package daemon

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"html"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

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

// countingProofScanStore counts ObjectsBySemanticID scans of the proof ID only
// (PD13), so AC1's zero stays honest if another pane scans a different ID.
type countingProofScanStore struct {
	readStore
	calls *atomic.Int32
}

func (s countingProofScanStore) ObjectsBySemanticID(ctx context.Context, id, after string, limit int) ([]store.Object, error) {
	if id == evidence.ProofSemanticID {
		s.calls.Add(1)
	}
	return s.readStore.ObjectsBySemanticID(ctx, id, after, limit)
}

func TestWorkbenchGradeUnconfiguredServesTheConstantWithoutScanning(t *testing.T) {
	d := newHandlerDaemon(t)
	proof := workbenchTestObject("grade-unconfigured-proof", evidence.ProofSemanticID, evidence.InterfaceHashV1)
	if err := d.store.PutObject(boundedTestContext(t), proof); err != nil {
		t.Fatalf("PutObject proof: %v", err)
	}
	subject := storeGradeSubject(t, d, "grade-unconfigured-subject")
	var n atomic.Int32
	d.reads = countingProofScanStore{readStore: d.reads, calls: &n}
	for _, ref := range []hashref.HashRef{proof.Hash, subject} {
		code, body := gradePage(t, d, ref)
		if code != http.StatusOK {
			t.Fatalf("object page %s: status %d, want 200", ref, code)
		}
		if c := strings.Count(body, wantGradeUnavailable(objectGradeUnavailableReason)); c != 1 {
			t.Fatalf("object page %s: constant grade line count %d, want 1", ref, c)
		}
		if strings.Contains(body, "<span>PROVEN</span>") {
			t.Fatalf("object page %s rendered PROVEN without a validator", ref)
		}
	}
	if got := n.Load(); got != 0 {
		t.Fatalf("ObjectsBySemanticID(proof) calls = %d, want 0", got)
	}

	// Control: the same wrapper on a configured daemon counts the scan.
	cd, _ := newGradeDaemon(t, []string{gradeTestIdentity})
	var m atomic.Int32
	cd.reads = countingProofScanStore{readStore: cd.reads, calls: &m}
	gradePage(t, cd, storeGradeSubject(t, cd, "grade-control-subject"))
	if got := m.Load(); got != 1 {
		t.Fatalf("control: ObjectsBySemanticID(proof) calls = %d, want 1", got)
	}
}

func TestWorkbenchGradeProvenFromSubjectBoundEnvelope(t *testing.T) {
	d, rig := newGradeDaemon(t, []string{gradeTestIdentity})
	subject := storeGradeSubject(t, d, "grade-proven-subject")
	env := mintGradeEnvelope(t, rig.key, gradeReport(rig, subject))
	envRef := storeGradeEnvelope(t, d, env)
	if rig.version != "AILANG v0.41.0 " {
		t.Fatalf("instrument: producer-side version %q, want the discriminating %q", rig.version, "AILANG v0.41.0 ")
	}
	res := d.grader.ValidateProof(boundedTestContext(t), envRef, subject)
	sealed, ok := res.Validated()
	if !ok {
		t.Fatalf("premise: ValidateProof did not validate (err %v)", res.Err())
	}
	if _, ok := d.grader.Resolve(sealed).Proven(); !ok {
		t.Fatalf("premise: Resolve did not prove (err %v)", d.grader.Resolve(sealed).Err())
	}

	code, body := gradePage(t, d, subject)
	if code != http.StatusOK {
		t.Fatalf("subject page: status %d, want 200", code)
	}
	if c := strings.Count(body, "<p><span>PROVEN</span></p>"); c != 1 {
		t.Fatalf("subject page: PROVEN line count %d, want 1", c)
	}
	if strings.Contains(body, "GRADE UNAVAILABLE") {
		t.Fatalf("subject page still says GRADE UNAVAILABLE beside PROVEN")
	}
	// "verdict-" alone is vacuous: the inline stylesheet always carries
	// .verdict-pass/.verdict-fail. Assert the channel's own markup
	// (host/workbench/render_test.go:136 uses the same two markers).
	for _, marker := range []string{`class="verdict-`, `aria-label="test verdict`} {
		if strings.Contains(body, marker) {
			t.Fatalf("subject page rendered a verdict channel: %s", marker)
		}
	}
	for _, label := range []string{"<span>TESTED</span>", "<span>ATTESTED</span>", "<span>CLAIMED</span>"} {
		if strings.Contains(body, label) {
			t.Fatalf("subject page rendered unminted label %s", label)
		}
	}

	// Honesty control: the envelope's own page is bound elsewhere, never PROVEN.
	code, body = gradePage(t, d, envRef)
	if code != http.StatusOK {
		t.Fatalf("envelope page: status %d, want 200", code)
	}
	if c := strings.Count(body, wantGradeUnavailable(gradeReasonBoundElsewhere)); c != 1 {
		t.Fatalf("envelope page: bound-elsewhere line count %d, want 1", c)
	}
}

// wantGradePremise checks the validator's own verdict on ref, so a page
// assertion is never made without the reason it depends on.
func wantGradePremise(t *testing.T, d *Daemon, ref, subject hashref.HashRef, want evidence.UnsupportedReason) {
	t.Helper()
	res := d.grader.ValidateProof(boundedTestContext(t), ref, subject)
	got, ok := res.Unsupported()
	if !ok || got != want {
		t.Fatalf("premise: ValidateProof(%s) reason = %q (unsupported %v, err %v), want %q", ref, got, ok, res.Err(), want)
	}
}

// wantGradePage requires the subject page to show exactly one grade line with
// the given sentence (escaped as the template renders it).
func wantGradePage(t *testing.T, d *Daemon, ref hashref.HashRef, sentence string) {
	t.Helper()
	code, body := gradePage(t, d, ref)
	if code != http.StatusOK {
		t.Fatalf("object page %s: status %d, want 200", ref, code)
	}
	if c := strings.Count(body, wantGradeUnavailable(sentence)); c != 1 {
		t.Fatalf("object page %s: line %q count %d, want 1", ref, sentence, c)
	}
}

// storeRawProofCandidate stores a proof-labelled object with an arbitrary
// interface and payload, for the structural-candidate reasons.
func storeRawProofCandidate(t *testing.T, d *Daemon, iface hashref.HashRef, payload []byte) hashref.HashRef {
	t.Helper()
	obj := store.Object{Hash: hashref.SumSHA256(payload), InterfaceHash: iface,
		SemanticID: evidence.ProofSemanticID, Provenance: "bounded-ai-check-producer", Payload: payload}
	if err := d.store.PutObject(boundedTestContext(t), obj); err != nil {
		t.Fatalf("PutObject candidate: %v", err)
	}
	return obj.Hash
}

func TestWorkbenchGradeReasonNoStoredEvidence(t *testing.T) {
	d, _ := newGradeDaemon(t, []string{gradeTestIdentity})
	subject := storeGradeSubject(t, d, "grade-no-evidence-subject")
	if rows, err := d.store.ObjectsBySemanticID(boundedTestContext(t), evidence.ProofSemanticID, "", 1); err != nil || len(rows) != 0 {
		t.Fatalf("premise: stored proof rows = %d (err %v), want 0", len(rows), err)
	}
	wantGradePage(t, d, subject, gradeReasonNoStoredEvidence)
}

func TestWorkbenchGradeReasonBoundToOtherSubjects(t *testing.T) {
	d, rig := newGradeDaemon(t, []string{gradeTestIdentity})
	subject := storeGradeSubject(t, d, "grade-bound-subject")
	other := hashref.SumSHA256([]byte("grade-other-subject"))
	envRef := storeGradeEnvelope(t, d, mintGradeEnvelope(t, rig.key, gradeReport(rig, other)))
	wantGradePremise(t, d, envRef, subject, evidence.UnsupportedSubjectMismatch)
	wantGradePage(t, d, subject, gradeReasonBoundElsewhere)
	_, body := gradePage(t, d, subject)
	if strings.Contains(body, html.EscapeString(gradeScanTruncation)) {
		t.Fatalf("a 1-candidate scan carried the truncation qualifier")
	}
}

func TestWorkbenchGradeReasonUnauthenticatedEnvelope(t *testing.T) {
	d, rig := newGradeDaemon(t, []string{gradeTestIdentity})
	subject := storeGradeSubject(t, d, "grade-unauth-subject")
	foreign := [32]byte{9, 9, 9}
	envRef := storeGradeEnvelope(t, d, mintGradeEnvelope(t, foreign, gradeReport(rig, subject)))
	wantGradePremise(t, d, envRef, subject, evidence.UnsupportedUnauthenticatedReport)
	wantGradePage(t, d, subject, gradeReasonUnauthenticated)
	_, body := gradePage(t, d, subject)
	if strings.Contains(body, html.EscapeString(gradeReasonBoundElsewhere)) {
		t.Fatalf("an unauthenticated envelope was reported as bound elsewhere")
	}
}

func TestWorkbenchGradeReasonToolMismatch(t *testing.T) {
	t.Run("compiler", func(t *testing.T) {
		d, rig := newGradeDaemon(t, []string{gradeTestIdentity})
		subject := storeGradeSubject(t, d, "grade-tool-compiler")
		report := gradeReport(rig, subject)
		report.Compiler = hashref.SumSHA256([]byte("another compiler"))
		envRef := storeGradeEnvelope(t, d, mintGradeEnvelope(t, rig.key, report))
		wantGradePremise(t, d, envRef, subject, evidence.UnsupportedToolMismatch)
		wantGradePage(t, d, subject, gradeReasonToolMismatch)
	})
	t.Run("trimspace-version", func(t *testing.T) {
		d, rig := newGradeDaemon(t, []string{gradeTestIdentity})
		subject := storeGradeSubject(t, d, "grade-tool-trimspace")
		report := gradeReport(rig, subject)
		report.CompilerVersion = releaseFromVersion(rig.rawVersion)
		envRef := storeGradeEnvelope(t, d, mintGradeEnvelope(t, rig.key, report))
		wantGradePremise(t, d, envRef, subject, evidence.UnsupportedToolMismatch)
		wantGradePage(t, d, subject, gradeReasonToolMismatch)
	})
	t.Run("mismatch-sorted-first", func(t *testing.T) {
		d, rig := newGradeDaemon(t, []string{gradeTestIdentity})
		subject := storeGradeSubject(t, d, "grade-tool-sorted")
		report := gradeReport(rig, subject)
		report.Compiler = hashref.SumSHA256([]byte("another compiler"))
		toolRef := storeGradeEnvelope(t, d, mintGradeEnvelope(t, rig.key, report))
		var foreignRef hashref.HashRef
		found := false
		for i := 0; i < 64 && !found; i++ {
			foreignSubject := hashref.SumSHA256([]byte(fmt.Sprintf("grade-foreign-%d", i)))
			env := mintGradeEnvelope(t, rig.key, gradeReport(rig, foreignSubject))
			if hashref.SumSHA256(env).String() < toolRef.String() {
				foreignRef = storeGradeEnvelope(t, d, env)
				found = true
			}
		}
		if !found {
			t.Fatal("no foreign envelope sorts before the tool-mismatch candidate")
		}
		rows, err := d.store.ObjectsBySemanticID(boundedTestContext(t), evidence.ProofSemanticID, "", 10)
		if err != nil || len(rows) < 2 || rows[0].Hash != foreignRef {
			t.Fatalf("premise: the foreign envelope must scan first (err %v, rows %d)", err, len(rows))
		}
		wantGradePremise(t, d, foreignRef, subject, evidence.UnsupportedSubjectMismatch)
		wantGradePremise(t, d, toolRef, subject, evidence.UnsupportedToolMismatch)
		wantGradePage(t, d, subject, gradeReasonToolMismatch)
	})
}

func TestWorkbenchGradeReasonProofFailed(t *testing.T) {
	d, rig := newGradeDaemon(t, []string{gradeTestIdentity})
	subject := storeGradeSubject(t, d, "grade-proof-failed")
	report := gradeReport(rig, subject)
	report.CheckPassed = false
	envRef := storeGradeEnvelope(t, d, mintGradeEnvelope(t, rig.key, report))
	wantGradePremise(t, d, envRef, subject, evidence.UnsupportedProofFailed)
	wantGradePage(t, d, subject, gradeReasonProofFailed)
}

func TestWorkbenchGradeReasonProofIncomplete(t *testing.T) {
	d, rig := newGradeDaemon(t, []string{gradeTestIdentity})
	subject := storeGradeSubject(t, d, "grade-proof-incomplete")
	report := gradeReport(rig, subject)
	report.Verified = []string{"world/other.ail"}
	envRef := storeGradeEnvelope(t, d, mintGradeEnvelope(t, rig.key, report))
	wantGradePremise(t, d, envRef, subject, evidence.UnsupportedProofIncomplete)
	wantGradePage(t, d, subject, gradeReasonProofIncomplete)
}

func TestWorkbenchGradeReasonStructuralCandidate(t *testing.T) {
	t.Run("wrong-interface", func(t *testing.T) {
		d, rig := newGradeDaemon(t, []string{gradeTestIdentity})
		subject := storeGradeSubject(t, d, "grade-structural-iface")
		env := mintGradeEnvelope(t, rig.key, gradeReport(rig, subject))
		ref := storeRawProofCandidate(t, d, hashref.SumSHA256([]byte("not the envelope interface")), env)
		wantGradePremise(t, d, ref, subject, evidence.UnsupportedWrongInterface)
		wantGradePage(t, d, subject, gradeReasonStructuralCandidate)
	})
	t.Run("malformed", func(t *testing.T) {
		d, _ := newGradeDaemon(t, []string{gradeTestIdentity})
		subject := storeGradeSubject(t, d, "grade-structural-malformed")
		ref := storeRawProofCandidate(t, d, evidence.InterfaceHashV1, []byte("not an envelope"))
		wantGradePremise(t, d, ref, subject, evidence.UnsupportedMalformed)
		wantGradePage(t, d, subject, gradeReasonStructuralCandidate)
	})
	t.Run("oversize", func(t *testing.T) {
		d, _ := newGradeDaemon(t, []string{gradeTestIdentity})
		subject := storeGradeSubject(t, d, "grade-structural-oversize")
		ref := storeRawProofCandidate(t, d, evidence.InterfaceHashV1, bytes.Repeat([]byte("x"), evidence.MaxBytes+1))
		wantGradePremise(t, d, ref, subject, evidence.UnsupportedOversize)
		wantGradePage(t, d, subject, gradeReasonStructuralCandidate)
	})
}

// TestWorkbenchGradeReasonScanTruncationIsDisclosed pins the exactly-100
// boundary: 100 candidates are a full scan with no qualifier, and the 101st
// row makes the scan bounded and says so.
func TestWorkbenchGradeReasonScanTruncationIsDisclosed(t *testing.T) {
	d, rig := newGradeDaemon(t, []string{gradeTestIdentity})
	subject := storeGradeSubject(t, d, "grade-trunc-subject")
	for i := 0; i < gradeCandidateScanLimit; i++ {
		other := hashref.SumSHA256([]byte(fmt.Sprintf("grade-trunc-%03d", i)))
		storeGradeEnvelope(t, d, mintGradeEnvelope(t, rig.key, gradeReport(rig, other)))
	}
	if rows, err := d.store.ObjectsBySemanticID(boundedTestContext(t), evidence.ProofSemanticID, "", gradeCandidateScanLimit+1); err != nil || len(rows) != gradeCandidateScanLimit {
		t.Fatalf("premise: stored rows = %d (err %v), want %d", len(rows), err, gradeCandidateScanLimit)
	}
	wantGradePage(t, d, subject, gradeReasonBoundElsewhere)
	_, body := gradePage(t, d, subject)
	if strings.Contains(body, html.EscapeString(gradeScanTruncation)) {
		t.Fatalf("a 100-candidate scan carried the truncation qualifier")
	}

	other := hashref.SumSHA256([]byte(fmt.Sprintf("grade-trunc-%03d", gradeCandidateScanLimit)))
	storeGradeEnvelope(t, d, mintGradeEnvelope(t, rig.key, gradeReport(rig, other)))
	if rows, err := d.store.ObjectsBySemanticID(boundedTestContext(t), evidence.ProofSemanticID, "", gradeCandidateScanLimit+1); err != nil || len(rows) != gradeCandidateScanLimit+1 {
		t.Fatalf("premise: stored rows = %d (err %v), want %d", len(rows), err, gradeCandidateScanLimit+1)
	}
	wantGradePage(t, d, subject, gradeReasonBoundElsewhere+gradeScanTruncation)
}

// scriptedReader serves one scripted object read to a real *evidence.Validator.
type scriptedReader struct {
	meta    store.ObjectMeta
	payload []byte
	err     error
}

func (r scriptedReader) ReadObject(context.Context, hashref.HashRef, int64) (store.ObjectMeta, []byte, error) {
	return r.meta, r.payload, r.err
}

func (r scriptedReader) BusyTimeout() time.Duration { return 0 }

// scriptedGrader returns ONLY results a real *evidence.Validator produced, or the
// zero value (PD1). It composes real validators; it never fabricates a seal.
type scriptedGrader struct {
	v        *evidence.Validator // validates (built over a scriptedReader)
	ref      *hashref.HashRef    // when non-nil, replaces the candidate ref passed to v
	resolver proofValidator      // when non-nil, Resolve is delegated here (a DIFFERENT real validator = foreign seal)
	zero     bool                // when true, ValidateProof returns evidence.ValidationResult{} (no verdict)
}

func (g scriptedGrader) ValidateProof(ctx context.Context, reportRef, expectedSubject hashref.HashRef) evidence.ValidationResult {
	if g.zero {
		return evidence.ValidationResult{}
	}
	if g.ref != nil {
		reportRef = *g.ref
	}
	return g.v.ValidateProof(ctx, reportRef, expectedSubject)
}

func (g scriptedGrader) Resolve(sealed evidence.ValidatedEvidence) evidence.ResolutionResult {
	if g.resolver != nil {
		return g.resolver.Resolve(sealed)
	}
	return g.v.Resolve(sealed)
}

// scriptedFixture is a configured daemon with one stored placeholder candidate,
// so the object page's proof scan is non-empty and the scripted grader runs.
type scriptedFixture struct {
	d       *Daemon
	rig     gradeRig
	subject hashref.HashRef
	cand    hashref.HashRef
}

func scriptedRig(t *testing.T) scriptedFixture {
	t.Helper()
	d, rig := newGradeDaemon(t, []string{gradeTestIdentity})
	subject := storeGradeSubject(t, d, "grade-scripted-subject")
	cand := workbenchTestObject("grade-scripted-candidate", evidence.ProofSemanticID, evidence.InterfaceHashV1)
	if err := d.store.PutObject(boundedTestContext(t), cand); err != nil {
		t.Fatalf("PutObject candidate: %v", err)
	}
	return scriptedFixture{d: d, rig: rig, subject: subject, cand: cand.Hash}
}

func newScriptedV(t *testing.T, rig gradeRig, r scriptedReader, key [32]byte) *evidence.Validator {
	t.Helper()
	v, err := evidence.NewValidator(key, r, evidence.CompilerConfig{Compiler: rig.compiler,
		CompilerVersion: rig.version, ObjectReadTimeout: 10 * time.Second}, []string{gradeTestIdentity})
	if err != nil {
		t.Fatalf("NewValidator: %v", err)
	}
	return v
}

const (
	premiseUnsupported = "unsupported" // the validator reports this UnsupportedReason
	premiseValidated   = "validated"   // the validator seals the candidate
	premiseNone        = "none"        // the zero result: no verdict at all
)

// gradeSeamArm is one row of the seam table (C5 AC4): a grader, the premise the
// test asserts on it directly, and the exact page it must render.
type gradeSeamArm struct {
	name     string
	grader   proofValidator
	premise  string
	reason   evidence.UnsupportedReason
	status   int
	sentence string // the unavailable sentence ("" for PROVEN and for 500)
	proven   bool
}

// gradeSeamArms is the 15-row table, shared by the mapping, the verdict sweep
// and the sanitized checks. Every grader is a real validator or the zero value.
func gradeSeamArms(t *testing.T, f scriptedFixture) []gradeSeamArm {
	t.Helper()
	rig, subject := f.rig, f.subject
	good := mintGradeEnvelope(t, rig.key, gradeReport(rig, subject))
	goodRef := hashref.SumSHA256(good)
	goodMeta := store.ObjectMeta{InterfaceHash: evidence.InterfaceHashV1, SemanticID: evidence.ProofSemanticID}
	ptr := func(h hashref.HashRef) *hashref.HashRef { return &h }
	nv := func(r scriptedReader) *evidence.Validator { return newScriptedV(t, rig, r, rig.key) }
	withReport := func(mut func(*evidence.ProofReportV1)) []byte {
		r := gradeReport(rig, subject)
		mut(&r)
		return mintGradeEnvelope(t, rig.key, r)
	}
	otherEnv := withReport(func(r *evidence.ProofReportV1) { r.Subject = hashref.SumSHA256([]byte("grade-scripted-other")) })
	toolEnv := withReport(func(r *evidence.ProofReportV1) { r.Compiler = hashref.SumSHA256([]byte("another compiler")) })
	failEnv := withReport(func(r *evidence.ProofReportV1) { r.Errors = 1 })
	incEnv := withReport(func(r *evidence.ProofReportV1) { r.Verified = []string{"world/other.ail"} })

	return []gradeSeamArm{
		{name: "invalid_ref", premise: premiseUnsupported, reason: evidence.UnsupportedInvalidRef, status: 200, sentence: gradeReasonStructuralCandidate,
			grader: scriptedGrader{v: nv(scriptedReader{meta: goodMeta, payload: good}), ref: &hashref.HashRef{}}},
		{name: "missing", premise: premiseUnsupported, reason: evidence.UnsupportedMissing, status: 200, sentence: gradeReasonStructuralCandidate,
			grader: scriptedGrader{v: nv(scriptedReader{meta: goodMeta}), ref: &goodRef}},
		{name: "oversize", premise: premiseUnsupported, reason: evidence.UnsupportedOversize, status: 200, sentence: gradeReasonStructuralCandidate,
			grader: scriptedGrader{v: nv(scriptedReader{err: &store.ObjectTooLargeError{Size: evidence.MaxBytes + 1, MaxBytes: evidence.MaxBytes}})}},
		{name: "hash_mismatch", premise: premiseUnsupported, reason: evidence.UnsupportedHashMismatch, status: 200, sentence: gradeReasonStructuralCandidate,
			grader: scriptedGrader{v: nv(scriptedReader{meta: goodMeta, payload: good})}},
		{name: "wrong_semantic_id", premise: premiseUnsupported, reason: evidence.UnsupportedWrongSemanticID, status: 200, sentence: gradeReasonStructuralCandidate,
			grader: scriptedGrader{v: nv(scriptedReader{meta: store.ObjectMeta{InterfaceHash: evidence.InterfaceHashV1, SemanticID: "other/id"}, payload: good}), ref: &goodRef}},
		{name: "wrong_interface", premise: premiseUnsupported, reason: evidence.UnsupportedWrongInterface, status: 200, sentence: gradeReasonStructuralCandidate,
			grader: scriptedGrader{v: nv(scriptedReader{meta: store.ObjectMeta{InterfaceHash: hashref.SumSHA256([]byte("x")), SemanticID: evidence.ProofSemanticID}, payload: good}), ref: &goodRef}},
		{name: "malformed", premise: premiseUnsupported, reason: evidence.UnsupportedMalformed, status: 200, sentence: gradeReasonStructuralCandidate,
			grader: scriptedGrader{v: nv(scriptedReader{meta: goodMeta, payload: []byte("{}")}), ref: ptr(hashref.SumSHA256([]byte("{}")))}},
		{name: "unauthenticated_report", premise: premiseUnsupported, reason: evidence.UnsupportedUnauthenticatedReport, status: 200, sentence: gradeReasonUnauthenticated,
			grader: scriptedGrader{v: newScriptedV(t, rig, scriptedReader{meta: goodMeta, payload: good}, [32]byte{7}), ref: &goodRef}},
		{name: "subject_mismatch", premise: premiseUnsupported, reason: evidence.UnsupportedSubjectMismatch, status: 200, sentence: gradeReasonBoundElsewhere,
			grader: scriptedGrader{v: nv(scriptedReader{meta: goodMeta, payload: otherEnv}), ref: ptr(hashref.SumSHA256(otherEnv))}},
		{name: "tool_mismatch", premise: premiseUnsupported, reason: evidence.UnsupportedToolMismatch, status: 200, sentence: gradeReasonToolMismatch,
			grader: scriptedGrader{v: nv(scriptedReader{meta: goodMeta, payload: toolEnv}), ref: ptr(hashref.SumSHA256(toolEnv))}},
		{name: "proof_failed", premise: premiseUnsupported, reason: evidence.UnsupportedProofFailed, status: 200, sentence: gradeReasonProofFailed,
			grader: scriptedGrader{v: nv(scriptedReader{meta: goodMeta, payload: failEnv}), ref: ptr(hashref.SumSHA256(failEnv))}},
		{name: "proof_incomplete", premise: premiseUnsupported, reason: evidence.UnsupportedProofIncomplete, status: 200, sentence: gradeReasonProofIncomplete,
			grader: scriptedGrader{v: nv(scriptedReader{meta: goodMeta, payload: incEnv}), ref: ptr(hashref.SumSHA256(incEnv))}},
		{name: "validated", premise: premiseValidated, status: 200, proven: true,
			grader: scriptedGrader{v: nv(scriptedReader{meta: goodMeta, payload: good}), ref: &goodRef}},
		{name: "resolve_refused", premise: premiseValidated, status: 500,
			grader: scriptedGrader{v: nv(scriptedReader{meta: goodMeta, payload: good}), ref: &goodRef,
				resolver: newScriptedV(t, rig, scriptedReader{}, rig.key)}},
		{name: "no_verdict", premise: premiseNone, status: 500,
			grader: scriptedGrader{zero: true}},
	}
}

// checkSeamPremise asserts the arm's validator outcome directly, so no page
// assertion is made without the premise that produces it.
func checkSeamPremise(t *testing.T, a gradeSeamArm, cand, subject hashref.HashRef) {
	t.Helper()
	res := a.grader.ValidateProof(boundedTestContext(t), cand, subject)
	switch a.premise {
	case premiseUnsupported:
		got, ok := res.Unsupported()
		if !ok || got != a.reason {
			t.Fatalf("%s: premise reason = %q (unsupported %v, err %v), want %q", a.name, got, ok, res.Err(), a.reason)
		}
	case premiseValidated:
		if _, ok := res.Validated(); !ok {
			t.Fatalf("%s: premise: not validated (err %v)", a.name, res.Err())
		}
	case premiseNone:
		_, validated := res.Validated()
		_, unsupported := res.Unsupported()
		if validated || unsupported || res.Err() != nil {
			t.Fatalf("%s: premise: want the zero result, got validated %v unsupported %v err %v", a.name, validated, unsupported, res.Err())
		}
	}
}

// TestWorkbenchGradeReasonMapsEveryUnsupportedReason is the AC4 table: every
// validator outcome renders its exact sentence, no other sentence appears, and
// the unknown-token arm runs on the pure production function.
func TestWorkbenchGradeReasonMapsEveryUnsupportedReason(t *testing.T) {
	f := scriptedRig(t)
	arms := gradeSeamArms(t, f)
	for _, a := range arms {
		// subject_mismatch is not one of the 11 mapped tokens: the scan loop
		// turns it into gradeReasonBoundElsewhere, so the map is checked on the rest.
		if a.premise == premiseUnsupported && a.reason != evidence.UnsupportedSubjectMismatch && gradeReasonSentence(a.reason) != a.sentence {
			t.Fatalf("%s: gradeReasonSentence(%q) = %q, want %q", a.name, a.reason, gradeReasonSentence(a.reason), a.sentence)
		}
	}
	if got := gradeReasonSentence(evidence.UnsupportedReason("future_reason")); got != gradeReasonNoneValidated {
		t.Fatalf("unknown token: gradeReasonSentence = %q, want %q", got, gradeReasonNoneValidated)
	}
	known := []string{gradeReasonNoStoredEvidence, gradeReasonBoundElsewhere, gradeReasonToolMismatch,
		gradeReasonProofFailed, gradeReasonProofIncomplete, gradeReasonUnauthenticated,
		gradeReasonStructuralCandidate, gradeReasonNoneValidated, objectGradeUnavailableReason}
	for _, a := range arms {
		f.d.grader = a.grader
		checkSeamPremise(t, a, f.cand, f.subject)
		code, body := gradePage(t, f.d, f.subject)
		if code != a.status {
			t.Fatalf("%s: status %d, want %d", a.name, code, a.status)
		}
		switch {
		case a.status == http.StatusInternalServerError:
			if !strings.Contains(body, ">Internal<") {
				t.Fatalf("%s: 500 body lacks the sanitized Internal marker", a.name)
			}
		case a.proven:
			if c := strings.Count(body, "<p><span>PROVEN</span></p>"); c != 1 {
				t.Fatalf("%s: PROVEN line count %d, want 1", a.name, c)
			}
		default:
			if c := strings.Count(body, wantGradeUnavailable(a.sentence)); c != 1 {
				t.Fatalf("%s: line %q count %d, want 1", a.name, a.sentence, c)
			}
		}
		for _, s := range known {
			if s != a.sentence && strings.Contains(body, wantGradeUnavailable(s)) {
				t.Fatalf("%s: unexpected grade sentence %q", a.name, s)
			}
		}
	}
}

// proofScanFailingStore fails the proof-ID scan with err and forwards every
// other read (the wrapper pattern of workbench_test.go's failing stores).
type proofScanFailingStore struct {
	readStore
	err error
}

func (s proofScanFailingStore) ObjectsBySemanticID(ctx context.Context, id, after string, limit int) ([]store.Object, error) {
	if id == evidence.ProofSemanticID {
		return nil, s.err
	}
	return s.readStore.ObjectsBySemanticID(ctx, id, after, limit)
}

// TestWorkbenchGradeScanFailureIsSanitized: an operational failure of the proof
// scan is a sanitized 500 or 503, never a reason sentence (AC5, M9).
func TestWorkbenchGradeScanFailureIsSanitized(t *testing.T) {
	cases := []struct {
		name   string
		err    error
		status int
		marker string
	}{
		{"internal", errSentinelInternal, http.StatusInternalServerError, ">Internal<"},
		{"deadline", context.DeadlineExceeded, http.StatusServiceUnavailable, ">Timeout<"},
		{"reference-index", &store.ReferenceIndexUnavailableError{}, http.StatusServiceUnavailable, ">ReferenceIndexUnavailable<"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d, _ := newGradeDaemon(t, []string{gradeTestIdentity})
			subject := storeGradeSubject(t, d, "grade-scan-"+c.name)
			buf := &bytes.Buffer{}
			d.errLog = buf
			d.reads = proofScanFailingStore{readStore: d.reads, err: c.err}
			code, body := gradePage(t, d, subject)
			if code != c.status || !strings.Contains(body, c.marker) {
				t.Fatalf("status %d body %q; want %d with %s", code, body, c.status, c.marker)
			}
			if strings.Contains(body, "GRADE UNAVAILABLE") {
				t.Fatalf("an operational failure rendered a grade reason")
			}
			if strings.Contains(body, internalDetailSentinel) {
				t.Fatalf("the response body leaks the internal detail")
			}
			if c.name == "internal" && !strings.Contains(buf.String(), internalDetailSentinel) {
				t.Fatalf("the operator log lacks the internal detail")
			}
		})
	}
	t.Run("control", func(t *testing.T) {
		d, _ := newGradeDaemon(t, []string{gradeTestIdentity})
		subject := storeGradeSubject(t, d, "grade-scan-control")
		d.reads = proofScanFailingStore{readStore: d.reads, err: nil}
		if code, _ := gradePage(t, d, subject); code != http.StatusOK {
			t.Fatalf("control: status %d, want 200", code)
		}
	})
}

// TestWorkbenchGradeValidationReadErrorIsSanitized: a validator read error is a
// sanitized 500 with the detail kept in the operator log only (AC5, M18).
func TestWorkbenchGradeValidationReadErrorIsSanitized(t *testing.T) {
	f := scriptedRig(t)
	buf := &bytes.Buffer{}
	f.d.errLog = buf
	f.d.grader = scriptedGrader{v: newScriptedV(t, f.rig, scriptedReader{err: errSentinelInternal}, f.rig.key)}
	if res := f.d.grader.ValidateProof(boundedTestContext(t), f.cand, f.subject); res.Err() == nil {
		t.Fatalf("premise: the read error did not reach the validator result")
	}
	code, body := gradePage(t, f.d, f.subject)
	if code != http.StatusInternalServerError || !strings.Contains(body, ">Internal<") {
		t.Fatalf("status %d body %q; want 500 with >Internal<", code, body)
	}
	if strings.Contains(body, internalDetailSentinel) {
		t.Fatalf("the response body leaks the read error detail")
	}
	if !strings.Contains(buf.String(), internalDetailSentinel) {
		t.Fatalf("the operator log lacks the read error detail")
	}
}

// TestWorkbenchGradeValidatorRefusalIsSanitized: a resolver that refuses the
// validator's own seal is a sanitized 500, never PROVEN (AC5, M5).
func TestWorkbenchGradeValidatorRefusalIsSanitized(t *testing.T) {
	f := scriptedRig(t)
	var refusal gradeSeamArm
	for _, a := range gradeSeamArms(t, f) {
		if a.name == "resolve_refused" {
			refusal = a
		}
	}
	f.d.grader = refusal.grader
	res := refusal.grader.ValidateProof(boundedTestContext(t), f.cand, f.subject)
	sealed, ok := res.Validated()
	if !ok {
		t.Fatalf("premise: not validated (err %v)", res.Err())
	}
	foreign := newScriptedV(t, f.rig, scriptedReader{}, f.rig.key)
	if err := foreign.Resolve(sealed).Err(); !errors.Is(err, evidence.ErrForeignSeal) {
		t.Fatalf("premise: foreign Resolve err = %v, want evidence.ErrForeignSeal", err)
	}
	buf := &bytes.Buffer{}
	f.d.errLog = buf
	code, body := gradePage(t, f.d, f.subject)
	if code != http.StatusInternalServerError || !strings.Contains(body, ">Internal<") {
		t.Fatalf("status %d; want 500 with >Internal<", code)
	}
	if strings.Contains(body, "<span>PROVEN</span>") {
		t.Fatalf("a refused seal rendered PROVEN")
	}
	if !strings.Contains(buf.String(), "refused the validator's own seal") {
		t.Fatalf("the operator log lacks the refusal detail")
	}
}

// TestWorkbenchGradeNeverRendersVerdictOrUnmintedLabels: across every validated
// or unavailable page, no verdict channel or unminted label is rendered, and
// PROVEN appears exactly once and only on the validated arm (AC6). The markers
// are the channel's markup, not the bare "verdict-" (see deviation 3): the
// inline stylesheet always carries .verdict-pass and .verdict-fail.
func TestWorkbenchGradeNeverRendersVerdictOrUnmintedLabels(t *testing.T) {
	f := scriptedRig(t)
	for _, a := range gradeSeamArms(t, f) {
		if a.status != http.StatusOK {
			continue
		}
		f.d.grader = a.grader
		code, body := gradePage(t, f.d, f.subject)
		if code != http.StatusOK {
			t.Fatalf("%s: status %d, want 200", a.name, code)
		}
		for _, m := range []string{`class="verdict-`, `aria-label="test verdict`, "<span>TESTED</span>", "<span>ATTESTED</span>", "<span>CLAIMED</span>"} {
			if strings.Contains(body, m) {
				t.Fatalf("%s: page renders %s", a.name, m)
			}
		}
		if a.proven {
			if c := strings.Count(body, "<p><span>PROVEN</span></p>"); c != 1 {
				t.Fatalf("%s: PROVEN line count %d, want 1", a.name, c)
			}
		} else if strings.Contains(body, "<span>PROVEN</span>") {
			t.Fatalf("%s: unavailable page renders PROVEN", a.name)
		}
	}
}
