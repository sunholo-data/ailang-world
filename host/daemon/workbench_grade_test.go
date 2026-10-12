package daemon

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"html"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/evidence"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// The discriminating version fixture (design V43): the producer's rule gives
// "AILANG v0.41.0 " and releaseFromVersion's TrimSpace gives "AILANG v0.41.0".
const gradeTestVersion = "AILANG v0.41.0 \r\nbuild abc\n"
const gradeTestIdentity = "world/types.ail"

// producerVersionLine is the TEST's own copy of proof_producer.go:328 — never
// the production proofCompilerVersionLine — so the two halves are independent.
func producerVersionLine(raw string) string {
	return strings.TrimRight(strings.SplitN(raw, "\n", 2)[0], "\r")
}

func wantGradeUnavailable(reason string) string {
	return "<p>GRADE UNAVAILABLE — " + html.EscapeString(reason) + "</p>"
}

type gradeRig struct {
	key          [32]byte
	compiler     hashref.HashRef // SumSHA256(fake interpreter bytes)
	version      string          // PRODUCER-side: producerVersionLine(stdout of the fake's own --version run)
	rawVersion   string          // that raw stdout
	bin, keyPath string
}

func writeProofKeyFile(t *testing.T, dir string, b []byte) string {
	t.Helper()
	p := filepath.Join(dir, "proof.key")
	if err := os.WriteFile(p, b, 0o600); err != nil {
		t.Fatalf("write proof key file: %v", err)
	}
	return p
}

// gradeDaemonConfig builds a configured Config over a fresh fake interpreter
// (gradeTestVersion) and a fresh 32-byte key file holding strings.Repeat("K", 32)
// (so a leak of the key bytes into an error is greppable), and fills the rig.
// The rig's version comes from RUNNING the fake.
func gradeDaemonConfig(t *testing.T, identities []string) (Config, gradeRig) {
	t.Helper()
	dir := t.TempDir()
	bin, binBytes := fakeInterpreter(t, dir, gradeTestVersion)
	raw, err := exec.CommandContext(boundedTestContext(t), bin, "--version").Output()
	if err != nil {
		t.Fatalf("run fake interpreter --version: %v", err)
	}
	keyPath := writeProofKeyFile(t, dir, []byte(strings.Repeat("K", 32)))
	var rig gradeRig
	copy(rig.key[:], strings.Repeat("K", 32))
	rig.compiler = hashref.SumSHA256(binBytes)
	rig.rawVersion = string(raw)
	rig.version = producerVersionLine(rig.rawVersion)
	rig.bin, rig.keyPath = bin, keyPath
	cfg := Config{DBPath: filepath.Join(dir, "world.db"), BindHost: DefaultBindHost, AilangBin: bin,
		ProofKeyFile: keyPath, ProofRequiredIdentities: identities, ErrorLog: testErrorLog(t)}
	return cfg, rig
}

func newGradeDaemon(t *testing.T, identities []string) (*Daemon, gradeRig) {
	t.Helper()
	cfg, rig := gradeDaemonConfig(t, identities)
	d, err := New(boundedTestContext(t), cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() {
		if err := d.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	})
	return d, rig
}

// mintGradeEnvelope builds an authenticated envelope from EXPORTED codecs only
// (PD6): the test cannot import evidence's unexported authority_test helpers.
func mintGradeEnvelope(t *testing.T, key [32]byte, report evidence.ProofReportV1) []byte {
	t.Helper()
	raw, err := evidence.EncodeProofReportV1(report)
	if err != nil {
		t.Fatalf("EncodeProofReportV1: %v", err)
	}
	m := hmac.New(sha256.New, key[:])
	_, _ = m.Write(raw)
	env, err := evidence.EncodeAuthenticatedEnvelope(evidence.AuthenticatedEnvelope{Report: raw, MAC: m.Sum(nil), MACValid: true})
	if err != nil {
		t.Fatalf("EncodeAuthenticatedEnvelope: %v", err)
	}
	return env
}

func gradeReport(rig gradeRig, subject hashref.HashRef) evidence.ProofReportV1 {
	return evidence.ProofReportV1{Schema: evidence.ProofReportSchemaV1, Subject: subject, Compiler: rig.compiler,
		CompilerVersion: rig.version, Verified: []string{gradeTestIdentity}, CheckPassed: true, ProofSucceeded: true}
}

func storeGradeEnvelope(t *testing.T, d *Daemon, env []byte) hashref.HashRef {
	t.Helper()
	obj := store.Object{Hash: hashref.SumSHA256(env), InterfaceHash: evidence.InterfaceHashV1,
		SemanticID: evidence.ProofSemanticID, Provenance: "bounded-ai-check-producer", Payload: env}
	if err := d.store.PutObject(boundedTestContext(t), obj); err != nil {
		t.Fatalf("PutObject envelope: %v", err)
	}
	return obj.Hash
}

func storeGradeSubject(t *testing.T, d *Daemon, label string) hashref.HashRef {
	t.Helper()
	obj := workbenchTestObject(label, "test/grade-subject", hashref.SumSHA256([]byte("interface-"+label)))
	if err := d.store.PutObject(boundedTestContext(t), obj); err != nil {
		t.Fatalf("PutObject subject %s: %v", label, err)
	}
	return obj.Hash
}

func gradePage(t *testing.T, d *Daemon, ref hashref.HashRef) (int, string) {
	t.Helper()
	rec := requestRecorder(t, d, http.MethodGet, "/workbench?object="+ref.String(), nil)
	return rec.Code, rec.Body.String()
}

// refuseGradeConfig requires New to refuse cfg at StageConfig with the exact
// detail, and returns the refusal for the arm's own Err checks.
func refuseGradeConfig(t *testing.T, cfg Config, detail string) *StartupError {
	t.Helper()
	d, err := New(boundedTestContext(t), cfg)
	if err == nil {
		if d != nil {
			_ = d.Close()
		}
		t.Fatalf("New accepted the config; want refusal %q", detail)
	}
	var se *StartupError
	if !errors.As(err, &se) {
		t.Fatalf("err = %v (%T), want *StartupError", err, err)
	}
	if se.Stage != StageConfig {
		t.Fatalf("Stage = %q, want %q", se.Stage, StageConfig)
	}
	if se.Detail != detail {
		t.Fatalf("Detail = %q, want %q", se.Detail, detail)
	}
	if strings.Contains(err.Error(), strings.Repeat("K", 32)) {
		t.Fatalf("startup error leaks key bytes: %v", err)
	}
	return se
}

func TestProofKeyWiringRefusesEveryHalfConfiguration(t *testing.T) {
	t.Run("control", func(t *testing.T) {
		cfg, _ := gradeDaemonConfig(t, []string{gradeTestIdentity})
		d, err := New(boundedTestContext(t), cfg)
		if err != nil {
			t.Fatalf("the unedited configuration was refused: %v", err)
		}
		if err := d.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	})
	t.Run("missing-file", func(t *testing.T) {
		cfg, _ := gradeDaemonConfig(t, []string{gradeTestIdentity})
		cfg.ProofKeyFile = filepath.Join(t.TempDir(), "absent.key")
		se := refuseGradeConfig(t, cfg, gradeConfigKeyRefused)
		if !errors.Is(se.Err, os.ErrNotExist) {
			t.Fatalf("cause = %v, want os.ErrNotExist", se.Err)
		}
	})
	t.Run("directory", func(t *testing.T) {
		cfg, _ := gradeDaemonConfig(t, []string{gradeTestIdentity})
		cfg.ProofKeyFile = t.TempDir()
		se := refuseGradeConfig(t, cfg, gradeConfigKeyRefused)
		if !strings.Contains(se.Error(), "is not a regular file") {
			t.Fatalf("error = %v, want it to say %q", se, "is not a regular file")
		}
	})
	t.Run("31-bytes", func(t *testing.T) {
		cfg, _ := gradeDaemonConfig(t, []string{gradeTestIdentity})
		cfg.ProofKeyFile = writeProofKeyFile(t, t.TempDir(), []byte(strings.Repeat("K", 31)))
		se := refuseGradeConfig(t, cfg, gradeConfigKeyRefused)
		if !strings.Contains(se.Error(), "holds 31 bytes; a proof key is exactly 32") {
			t.Fatalf("error = %v, want the 31-byte refusal", se)
		}
	})
	t.Run("33-bytes", func(t *testing.T) {
		cfg, _ := gradeDaemonConfig(t, []string{gradeTestIdentity})
		cfg.ProofKeyFile = writeProofKeyFile(t, t.TempDir(), []byte(strings.Repeat("K", 33)))
		se := refuseGradeConfig(t, cfg, gradeConfigKeyRefused)
		if !strings.Contains(se.Error(), "holds more than 32 bytes") {
			t.Fatalf("error = %v, want the 33-byte refusal", se)
		}
	})
	t.Run("no-ailang-bin", func(t *testing.T) {
		cfg, _ := gradeDaemonConfig(t, []string{gradeTestIdentity})
		cfg.AilangBin = ""
		se := refuseGradeConfig(t, cfg, gradeConfigNoInterpreter)
		if se.Err != nil {
			t.Fatalf("cause = %v, want nil", se.Err)
		}
	})
	t.Run("no-identities", func(t *testing.T) {
		cfg, _ := gradeDaemonConfig(t, nil)
		se := refuseGradeConfig(t, cfg, gradeConfigValidatorRefused)
		if !errors.Is(se.Err, evidence.ErrInvalidValidatorConfig) {
			t.Fatalf("cause = %v, want evidence.ErrInvalidValidatorConfig", se.Err)
		}
		if !strings.Contains(se.Error(), "required identities are empty") {
			t.Fatalf("error = %v, want %q", se, "required identities are empty")
		}
	})
	t.Run("identities-without-key", func(t *testing.T) {
		cfg, _ := gradeDaemonConfig(t, []string{gradeTestIdentity})
		cfg.ProofKeyFile = ""
		refuseGradeConfig(t, cfg, gradeConfigIdentitiesWithoutKey)
	})
}

func TestProofKeyWiringConstructsTheValidator(t *testing.T) {
	d, _ := newGradeDaemon(t, []string{gradeTestIdentity})
	if v, ok := d.grader.(*evidence.Validator); !ok || v == nil {
		t.Fatalf("configured d.grader = %T, want a constructed *evidence.Validator", d.grader)
	}
	u := newHandlerDaemon(t)
	if u.grader != nil {
		t.Fatalf("unconfigured d.grader = %T, want nil", u.grader)
	}
}

func TestWorkbenchGradeCompilerVersionDerivationsAgree(t *testing.T) {
	d, rig := newGradeDaemon(t, []string{gradeTestIdentity})
	if rig.rawVersion != gradeTestVersion {
		t.Fatalf("instrument: fake --version printed %q, want the fixture %q", rig.rawVersion, gradeTestVersion)
	}
	if d.interpreterVersion != rig.rawVersion {
		t.Fatalf("instrument: d.interpreterVersion = %q, want the raw %q", d.interpreterVersion, rig.rawVersion)
	}
	p := producerVersionLine(rig.rawVersion)
	if p != "AILANG v0.41.0 " {
		t.Fatalf("producer-side line = %q, want %q", p, "AILANG v0.41.0 ")
	}
	dv := proofCompilerVersionLine(d.interpreterVersion)
	if dv != p {
		t.Fatalf("daemon-side %q != producer-side %q", dv, p)
	}
	if r := releaseFromVersion(d.interpreterVersion); r != "AILANG v0.41.0" || r == p {
		t.Fatalf("discrimination control: releaseFromVersion = %q, want %q and != %q", r, "AILANG v0.41.0", p)
	}
	t.Run("pinned-binary", func(t *testing.T) {
		bin := os.Getenv("AILANG_BIN")
		if bin == "" {
			t.Skip("observation only")
		}
		out, err := exec.CommandContext(boundedTestContext(t), bin, "--version").Output()
		if err != nil {
			t.Fatalf("run pinned --version: %v", err)
		}
		got := proofCompilerVersionLine(string(out))
		if got != "AILANG v0.41.0" || producerVersionLine(string(out)) != got {
			t.Fatalf("pinned binary: daemon-side %q, producer-side %q, want both %q", got, producerVersionLine(string(out)), "AILANG v0.41.0")
		}
	})
}
