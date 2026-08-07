package broker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/childenv"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/pkgproj"
	"github.com/sunholo-data/ailang-world/host/replay"
	"github.com/sunholo-data/ailang-world/host/store"
)

// ---------------------------------------------------------------------------
// SAFETY
//
// Nothing in this file may reach the public AILANG registry. The registry is
// immutable and a publisher cannot recall a version, so a single stray request
// is permanent. Three independent things make that structurally true rather
// than a matter of care:
//
//  1. every handler here is built by newLoopbackRegistryPublishHandler, which
//     is UNEXPORTED and REFUSES a non-loopback origin (see
//     TestLoopbackConstructorRefusesEveryNonLoopbackOrigin);
//  2. the fake validator is an httptest server, which binds 127.0.0.1 only,
//     and every request it ever sees is enumerated by
//     TestFakeValidatorSawOnlyLoopbackTraffic;
//  3. the publisher subprocess is NEVER the real `ailang` binary — it is this
//     test binary re-exec'd into TestRegistryPublishHelperProcess.
//
// The environment guard is TestRegistryEnvironmentPointsOnlyAtLoopback.
// ---------------------------------------------------------------------------

const (
	fixtureVendor  = "world"
	fixtureName    = "core"
	fixtureVersion = "0.1.0"
	fixturePackage = fixtureVendor + "/" + fixtureName

	// helperModeEnv and friends are read ONLY by the re-exec'd helper process
	// below. They are exported into the child by the generated shell script,
	// never by the handler, so the handler's minimal environment stays exactly
	// what childEnvironment builds.
	helperModeEnv      = "WORLD_FAKE_PUBLISHER_MODE"
	helperDumpEnv      = "WORLD_FAKE_PUBLISHER_DUMP"
	helperValidatorEnv = "WORLD_FAKE_PUBLISHER_DEAD_URL"

	// ac10Sentinel stands in for the credential in every test that must prove
	// a secret does not travel. It is deliberately not a plausible key.
	ac10Sentinel = "world-ac10-sentinel-value-not-a-real-key"
	// ac10Marker is the NON-secret companion injected the same way. Requiring
	// it to APPEAR is what proves the redaction scanner read the stream at
	// all; without it, "the sentinel is absent" is satisfied by an empty
	// string.
	ac10Marker = "world-ac10-nonsecret-marker"
)

// ---------------------------------------------------------------------------
// the loopback fake validator
// ---------------------------------------------------------------------------

// validatorRequest is one observed request. It records whether an API key
// header was present and how long it was — NEVER its value.
type validatorRequest struct {
	Method     string
	Path       string
	Host       string
	RemoteAddr string
	HasAPIKey  bool
	APIKeyLen  int
}

type fakeValidator struct {
	server  *httptest.Server
	mu      sync.Mutex
	log     []validatorRequest
	mode    string
	release chan struct{}
}

// newFakeValidator starts a loopback-only stand-in for the registry validator
// service. mode selects the arm: "ok" (200), "namespace" (403), "validation"
// (400 echoing the API key and the non-secret marker), "reset" (accept the
// request, then abort the connection — the genuinely ambiguous case, because
// the request body HAS arrived), and "hang" (never answer).
func newFakeValidator(t *testing.T, mode string) *fakeValidator {
	t.Helper()
	v := &fakeValidator{mode: mode, release: make(chan struct{})}
	v.server = httptest.NewServer(http.HandlerFunc(v.serve))
	// Cleanup is LIFO: registering Close first means it runs LAST, after the
	// release below has unblocked any request parked in the "hang" arm.
	t.Cleanup(v.server.Close)
	t.Cleanup(func() { close(v.release) })
	return v
}

func (v *fakeValidator) serve(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("X-API-Key")
	v.mu.Lock()
	v.log = append(v.log, validatorRequest{
		Method: r.Method, Path: r.URL.Path, Host: r.Host,
		RemoteAddr: r.RemoteAddr, HasAPIKey: key != "", APIKeyLen: len(key),
	})
	v.mu.Unlock()
	_, _ = io.Copy(io.Discard, r.Body)

	switch v.mode {
	case "namespace":
		w.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(w, "forbidden")
	case "validation":
		// MUT-SM-SECRET-ERROR's fixture: the server echoes the credential it
		// was sent, alongside a non-secret marker injected identically.
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "rejected marker="+ac10Marker+" key="+key)
	case "reset":
		hijacker, ok := w.(http.Hijacker)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		conn, _, err := hijacker.Hijack()
		if err != nil {
			return
		}
		// The request is already logged, so the body demonstrably left the
		// publisher. Killing the connection here is exactly Decision 3's
		// "reset after request body may have left process".
		_ = conn.Close()
	case "hang":
		<-v.release
	default:
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "ok")
	}
}

func (v *fakeValidator) origin() string { return v.server.URL }

func (v *fakeValidator) requests() []validatorRequest {
	v.mu.Lock()
	defer v.mu.Unlock()
	return append([]validatorRequest(nil), v.log...)
}

func (v *fakeValidator) count() int { return len(v.requests()) }

// ---------------------------------------------------------------------------
// the fake publisher: this test binary, re-exec'd
// ---------------------------------------------------------------------------

// TestRegistryPublishHelperProcess is not a test of this package: it is the
// body of the fake publisher subprocess, mirroring host/store's
// TestCrashHelperProcess idiom. It stands in for `ailang publish` and prints
// the v0.30.0 messages this handler classifies (e37b370:cmd/ailang/pkg_publish.go).
//
// The real binary is never invoked: an accidental live publish is unrecallable.
func TestRegistryPublishHelperProcess(t *testing.T) {
	mode := os.Getenv(helperModeEnv)
	if mode == "" {
		t.Skip("subprocess helper; runs only when re-exec'd with " + helperModeEnv)
	}
	if dump := os.Getenv(helperDumpEnv); dump != "" {
		if err := os.WriteFile(dump, []byte(strings.Join(os.Environ(), "\n")), 0o600); err != nil {
			fmt.Fprintf(os.Stderr, "publisher helper: write env dump: %v\n", err)
			os.Exit(9)
		}
	}

	post := func(target string) (int, string, error) {
		req, err := http.NewRequest(http.MethodPost, target+"/publish", bytes.NewReader([]byte("tarball")))
		if err != nil {
			return 0, "", err
		}
		req.Header.Set("Content-Type", "multipart/form-data")
		if key := os.Getenv(RegistryCredentialVariable); key != "" {
			req.Header.Set("X-API-Key", key)
		}
		resp, err := (&http.Client{}).Do(req)
		if err != nil {
			return 0, "", err
		}
		defer func() { _ = resp.Body.Close() }()
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, string(body), nil
	}

	validator := os.Getenv("AILANG_REGISTRY_VALIDATOR")
	switch mode {
	case "envdump":
		fmt.Printf("Published %s@%s\n", fixturePackage, fixtureVersion)
		os.Exit(0)

	case "dryrun":
		fmt.Printf("  Tarball: 1 bytes (sha256:%s...)\n", strings.Repeat("0", 17))
		fmt.Println("⚠ Dry run complete. Tarball ready but not uploaded.")
		os.Exit(0)

	case "transport":
		// Point at a loopback port nothing is listening on.
		_, _, err := post(os.Getenv(helperValidatorEnv))
		fmt.Printf("upload failed: %v\n", err)
		os.Exit(1)

	case "success", "namespace", "validation", "reset", "hang":
		status, body, err := post(validator)
		if err != nil {
			// v0.30.0 wraps every client.Do error exactly this way.
			fmt.Printf("upload failed: %v\n", err)
			os.Exit(1)
		}
		switch status {
		case http.StatusOK:
			fmt.Printf("Published %s@%s\n", fixturePackage, fixtureVersion)
			os.Exit(0)
		case http.StatusForbidden:
			fmt.Println("not authorized to publish to this namespace")
			os.Exit(1)
		case http.StatusBadRequest:
			fmt.Printf("validation failed:\n%s\n", body)
			os.Exit(1)
		default:
			fmt.Printf("publish failed (HTTP %d): %s\n", status, body)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "publisher helper: unknown mode %q\n", mode)
		os.Exit(8)
	}
}

// writePublisherScript generates the executable the handler launches. It is a
// shell wrapper rather than the test binary directly because the handler puts
// "publish" first in the argument vector, which would stop Go's flag parsing
// before -test.run and cause the child to run the whole package.
func writePublisherScript(t *testing.T, mode string, extra map[string]string) string {
	t.Helper()
	var body strings.Builder
	fmt.Fprintf(&body, "%s=%q\nexport %s\n", helperModeEnv, mode, helperModeEnv)
	names := make([]string, 0, len(extra))
	for name := range extra {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Fprintf(&body, "%s=%q\nexport %s\n", name, extra[name], name)
	}
	fmt.Fprintf(&body, "exec %q -test.run='^TestRegistryPublishHelperProcess$'\n", os.Args[0])

	path := filepath.Join(t.TempDir(), "fake-ailang")
	if err := os.WriteFile(path, []byte("#!/bin/sh\nset -eu\n"+body.String()), 0o700); err != nil {
		t.Fatal(err)
	}
	return path
}

// writeProbeScript generates a minimal executable that records the environment
// it was handed and then succeeds. It is the AC10(a) instrument, and its
// ability to SEE the variable is proven by a known-positive control before any
// absence it reports is believed.
func writeProbeScript(t *testing.T, dump string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "env-probe")
	body := "#!/bin/sh\n/usr/bin/env > " + shellQuote(dump) + "\n" +
		"echo 'ailang version v0.30.0 (AC10 probe)'\nexit 0\n"
	if err := os.WriteFile(path, []byte(body), 0o700); err != nil {
		t.Fatal(err)
	}
	return path
}

func shellQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'" }

func readEnvDump(t *testing.T, path string) []string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read env dump %q: %v", path, err)
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return nil
	}
	return lines
}

// ---------------------------------------------------------------------------
// fixtures
// ---------------------------------------------------------------------------

func publishFixtureDir(t *testing.T) (string, pkgproj.Manifest) {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "world"), 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"ailang.toml":        "[package]\nname = \"" + fixturePackage + "\"\nversion = \"" + fixtureVersion + "\"\n",
		"world/types.ail":    "module world/types\n",
		"world/logepoch.ail": "module world/logepoch\n",
	}
	for rel, body := range files {
		if err := os.WriteFile(filepath.Join(dir, filepath.FromSlash(rel)), []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return dir, pkgproj.Manifest{
		Package: pkgproj.Package{Name: fixturePackage, Edition: "2026", AILANG: ">=0.30.0"},
		Exports: pkgproj.Exports{Modules: []string{"world/logepoch", "world/types"}},
		Effects: pkgproj.Effects{Max: []string{}},
	}
}

type publishFixture struct {
	dir      string
	manifest pkgproj.Manifest
	identity PublishIdentity
	hashes   PublishHashes
	approval PublishApproval
	scope    string
	payload  []byte
}

func newPublishFixture(t *testing.T, registryOrigin, approvalSeed string) publishFixture {
	t.Helper()
	dir, manifest := publishFixtureDir(t)
	hashes, err := recomputePublishHashes(dir, manifest)
	if err != nil {
		t.Fatal(err)
	}
	id := PublishIdentity{
		Vendor: fixtureVendor, Name: fixtureName, Version: fixtureVersion,
		RegistryOrigin:  registryOrigin,
		ManifestRef:     hashref.SumSHA256([]byte("manifest:" + approvalSeed)),
		ApprovalRef:     hashref.SumSHA256([]byte("approval:" + approvalSeed)),
		Exports:         manifest.Exports.Modules,
		Effects:         manifest.Effects.Max,
		CompilerVersion: "AILANG v0.30.0",
		CompilerSHA256:  "sha256:" + strings.Repeat("ab", 32),
	}
	return publishFixture{
		dir: dir, manifest: manifest, identity: id, hashes: hashes,
		approval: PublishApproval{
			ApprovalRef: id.ApprovalRef, Vendor: id.Vendor, Name: id.Name,
			Version: id.Version, RegistryOrigin: registryOrigin, Hashes: hashes,
		},
		scope:   PublishScope(registryOrigin, fixtureVendor, fixtureName, fixtureVersion),
		payload: EncodePublishPayload(id, hashes),
	}
}

func writeCredentialFile(t *testing.T, secret string) *FileRegistryCredentialProvider {
	t.Helper()
	// t.TempDir() is outside the working tree, which is exactly what the
	// provider requires.
	path := filepath.Join(t.TempDir(), "registry.key")
	if err := os.WriteFile(path, []byte(secret+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	provider, err := NewFileRegistryCredentialProvider(path, testRepoRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	return provider
}

func testRepoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate repository root")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

// publishRecordingStore is the objectStore double for publish sessions. It
// captures the minted effect invocation IDs so a test can ask the journal for
// the receipt STATE of a specific intent rather than inferring it.
type publishRecordingStore struct {
	base      *store.Store
	records   []store.Object
	effectIDs []string
}

func (s *publishRecordingStore) PutObject(obj store.Object) error {
	if err := s.base.PutObject(obj); err != nil {
		return err
	}
	if obj.SemanticID == EffectRecordV1 {
		s.records = append(s.records, obj)
	}
	return nil
}

func (s *publishRecordingStore) GetObject(ref hashref.HashRef) (store.Object, bool, error) {
	return s.base.GetObject(ref)
}

func (s *publishRecordingStore) AppendNextEffectIntent(
	episodeID string, intent store.EffectIntent,
) (string, int64, error) {
	id, ordinal, err := s.base.AppendNextEffectIntent(episodeID, intent)
	if err == nil {
		s.effectIDs = append(s.effectIDs, id)
	}
	return id, ordinal, err
}

func (s *publishRecordingStore) AppendClaimedEffectIntent(
	episodeID string, intent store.EffectIntent, approvalRef, requestRef hashref.HashRef,
) (string, int64, error) {
	id, ordinal, err := s.base.AppendClaimedEffectIntent(episodeID, intent, approvalRef, requestRef)
	if err == nil {
		s.effectIDs = append(s.effectIDs, id)
	}
	return id, ordinal, err
}

func (s *publishRecordingStore) AppendEffectOutcome(
	id string, outcome store.EffectOutcome,
) (int64, hashref.HashRef, error) {
	return s.base.AppendEffectOutcome(id, outcome)
}

func publishSession(
	t *testing.T, base *store.Store, episodeID string, grant Capability, handler Handler,
) (*Session, *publishRecordingStore) {
	t.Helper()
	recording := &publishRecordingStore{base: base}
	session := newSession(recording, episodeID, []Capability{grant},
		Registry{EffectRegistryPublish: handler}, Live, nil)
	return session, recording
}

func loopbackHandler(t *testing.T, cfg RegistryPublishConfig) *RegistryPublishHandler {
	t.Helper()
	handler, err := newLoopbackRegistryPublishHandler(cfg)
	if err != nil {
		t.Fatalf("newLoopbackRegistryPublishHandler: %v", err)
	}
	return handler
}

// ---------------------------------------------------------------------------
// AC7 — denial is persisted, dispatch never happens, and a valid grant does
// ---------------------------------------------------------------------------

func TestPublishDenialsPersistAndNeverDispatchWithLivePositiveControl(t *testing.T) {
	validator := newFakeValidator(t, "ok")
	fixture := newPublishFixture(t, validator.origin(), "ac7")
	handler := loopbackHandler(t, RegistryPublishConfig{
		PublisherPath:   writePublisherScript(t, "success", nil),
		PackageDir:      fixture.dir,
		Manifest:        fixture.manifest,
		RegistryOrigin:  validator.origin(),
		ValidatorOrigin: validator.origin(),
		Credential:      writeCredentialFile(t, ac10Sentinel),
		Approval:        fixture.approval,
		ExecTimeout:     20 * time.Second,
	})

	const now = 50
	validGrant := Capability{
		Effect: EffectRegistryPublish, Scope: fixture.scope, ExpiresAt: 100, Budget: 1,
	}
	denials := []struct {
		name  string
		grant Capability
		label string
	}{
		{"no Registry.Publish grant at all", Capability{
			Effect: "FS.Write", Scope: fixture.scope, ExpiresAt: 100, Budget: 1,
		}, LabelDeniedEffectName},
		{"grant scopes a different version", Capability{
			Effect:    EffectRegistryPublish,
			Scope:     PublishScope(validator.origin(), fixtureVendor, fixtureName, "0.1.1"),
			ExpiresAt: 100, Budget: 1,
		}, LabelDeniedScope},
		{"grant expired at the caller-supplied logical time", Capability{
			Effect: EffectRegistryPublish, Scope: fixture.scope, ExpiresAt: 10, Budget: 1,
		}, LabelDeniedExpired},
		{"grant budget is zero", Capability{
			Effect: EffectRegistryPublish, Scope: fixture.scope, ExpiresAt: 100, Budget: 0,
		}, LabelDeniedBudget},
	}

	for _, tc := range denials {
		t.Run(tc.name, func(t *testing.T) {
			session, recording := publishSession(t, openTestStore(t), "ac7-denial", tc.grant, handler)
			_, ref, err := session.Invoke(context.Background(), EffectRequest{
				Effect: EffectRegistryPublish, Scope: fixture.scope, Cost: PublishCost, Now: now,
			}, fixture.payload)

			var denied *DenialError
			if !errors.As(err, &denied) {
				t.Fatalf("Invoke error = %T %v, want *DenialError", err, err)
			}
			if denied.Decision.Label != tc.label {
				t.Errorf("denial label = %q, want %q", denied.Decision.Label, tc.label)
			}
			if ref.IsZero() {
				t.Fatal("denial returned a zero record ref: nothing was persisted")
			}
			if got := len(recording.records); got != 1 {
				t.Fatalf("persisted denial records = %d, want exactly 1", got)
			}
			rec, decodeErr := DecodeRecord(recording.records[0].Payload)
			if decodeErr != nil {
				t.Fatal(decodeErr)
			}
			if rec.Allowed || rec.Denial != tc.label || rec.Effect != EffectRegistryPublish {
				t.Errorf("persisted denial record = %#v", rec)
			}
			// The record must be READABLE FROM THE STORE, not merely observed
			// passing through the double.
			if _, ok, getErr := recording.base.GetObject(ref); getErr != nil || !ok {
				t.Fatalf("denial record %s absent from the store (ok=%v err=%v)", ref, ok, getErr)
			}
			if len(recording.effectIDs) != 0 {
				t.Errorf("a denial minted durable effect intents %v, want none", recording.effectIDs)
			}
		})
	}

	if got := handler.Dispatches(); got != 0 {
		t.Fatalf("publisher dispatches after %d denials = %d, want 0", len(denials), got)
	}
	if got := handler.CredentialLoads(); got != 0 {
		t.Fatalf("credential loads after %d denials = %d, want 0", len(denials), got)
	}
	if got := validator.count(); got != 0 {
		t.Fatalf("validator requests after %d denials = %d, want 0", len(denials), got)
	}

	// SAME-RUN POSITIVE CONTROL. A zero counter proves nothing unless the same
	// counter, the same handler and the same validator can be shown to move.
	session, recording := publishSession(t, openTestStore(t), "ac7-allowed", validGrant, handler)
	result, ref, err := session.Invoke(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish, Scope: fixture.scope, Cost: PublishCost, Now: now,
	}, fixture.payload)
	if err != nil {
		t.Fatalf("positive control Invoke: %v", err)
	}
	if ref.IsZero() || len(result) == 0 {
		t.Fatalf("positive control produced ref %s and %d result bytes", ref, len(result))
	}
	if got := handler.Dispatches(); got != 1 {
		t.Fatalf("publisher dispatches after the valid grant = %d, want exactly 1", got)
	}
	if got := handler.CredentialLoads(); got != 1 {
		t.Fatalf("credential loads after the valid grant = %d, want exactly 1", got)
	}
	if got := validator.count(); got != 1 {
		t.Fatalf("validator requests after the valid grant = %d, want exactly 1", got)
	}
	if got := validator.requests()[0]; !got.HasAPIKey || got.Path != "/publish" {
		t.Errorf("validator saw %+v, want a POST to /publish carrying an API key", got)
	}
	if got := session.grants[0].Budget; got != 0 {
		t.Errorf("budget after one publish = %d, want 0", got)
	}
	receipt := effectReceipt(t, recording, 0)
	if receipt.State != store.ReceiptResolved {
		t.Errorf("positive control receipt state = %q, want %q", receipt.State, store.ReceiptResolved)
	}
}

func effectReceipt(t *testing.T, recording *publishRecordingStore, index int) store.Receipt {
	t.Helper()
	if index >= len(recording.effectIDs) {
		t.Fatalf("no effect intent at index %d; minted %v", index, recording.effectIDs)
	}
	receipt, ok, err := recording.base.GetEffectReceipt(recording.effectIDs[index])
	if err != nil || !ok {
		t.Fatalf("GetEffectReceipt(%q) = ok %v, err %v", recording.effectIDs[index], ok, err)
	}
	return receipt
}

// ---------------------------------------------------------------------------
// AC11 — the two outcome arms must DIFFER in the same run
// ---------------------------------------------------------------------------

func TestDefiniteFailureResolvesWhileAmbiguityStaysIndeterminate(t *testing.T) {
	base := openTestStore(t)
	grant := func(scope string) Capability {
		return Capability{Effect: EffectRegistryPublish, Scope: scope, ExpiresAt: 100, Budget: 1}
	}

	// ARM A — a DEFINITE handler failure: the validator answers 403, so the
	// attempt is over and its result is known.
	definiteValidator := newFakeValidator(t, "namespace")
	definiteFixture := newPublishFixture(t, definiteValidator.origin(), "ac11-definite")
	definiteHandler := loopbackHandler(t, RegistryPublishConfig{
		PublisherPath:   writePublisherScript(t, "namespace", nil),
		PackageDir:      definiteFixture.dir,
		Manifest:        definiteFixture.manifest,
		RegistryOrigin:  definiteValidator.origin(),
		ValidatorOrigin: definiteValidator.origin(),
		Credential:      writeCredentialFile(t, ac10Sentinel),
		Approval:        definiteFixture.approval,
		ExecTimeout:     20 * time.Second,
	})
	definiteSession, definiteRecording := publishSession(
		t, base, "ac11-definite", grant(definiteFixture.scope), definiteHandler)
	_, definiteRef, definiteErr := definiteSession.Invoke(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish, Scope: definiteFixture.scope, Cost: PublishCost, Now: 50,
	}, definiteFixture.payload)

	var failed *EffectFailedError
	if !errors.As(definiteErr, &failed) {
		t.Fatalf("definite arm error = %T %v, want *EffectFailedError", definiteErr, definiteErr)
	}
	if definiteRef.IsZero() {
		t.Fatal("definite arm produced no effect record")
	}
	var wronglyIndeterminate *IndeterminateEffectError
	if errors.As(definiteErr, &wronglyIndeterminate) {
		t.Fatalf("definite arm was classified indeterminate: %v", definiteErr)
	}
	definiteReceipt := effectReceipt(t, definiteRecording, 0)

	// ARM B — a typed AMBIGUOUS transport result on the SAME store, in the
	// SAME run. The validator accepts the request (so it is logged: the body
	// demonstrably left the publisher) and then aborts the connection.
	ambiguousValidator := newFakeValidator(t, "reset")
	ambiguousFixture := newPublishFixture(t, ambiguousValidator.origin(), "ac11-ambiguous")
	ambiguousHandler := loopbackHandler(t, RegistryPublishConfig{
		PublisherPath:   writePublisherScript(t, "reset", nil),
		PackageDir:      ambiguousFixture.dir,
		Manifest:        ambiguousFixture.manifest,
		RegistryOrigin:  ambiguousValidator.origin(),
		ValidatorOrigin: ambiguousValidator.origin(),
		Credential:      writeCredentialFile(t, ac10Sentinel),
		Approval:        ambiguousFixture.approval,
		ExecTimeout:     20 * time.Second,
	})
	ambiguousSession, ambiguousRecording := publishSession(
		t, base, "ac11-ambiguous", grant(ambiguousFixture.scope), ambiguousHandler)
	_, ambiguousRef, ambiguousErr := ambiguousSession.Invoke(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish, Scope: ambiguousFixture.scope, Cost: PublishCost, Now: 50,
	}, ambiguousFixture.payload)

	var indeterminate *IndeterminateEffectError
	if !errors.As(ambiguousErr, &indeterminate) {
		t.Fatalf("ambiguous arm error = %T %v, want *IndeterminateEffectError", ambiguousErr, ambiguousErr)
	}
	if !ambiguousRef.IsZero() {
		t.Errorf("ambiguous arm produced effect record %s, want none", ambiguousRef)
	}
	if len(ambiguousRecording.records) != 0 {
		t.Errorf("ambiguous arm persisted %d effect records, want 0", len(ambiguousRecording.records))
	}
	if indeterminate.InvocationID == "" {
		t.Error("ambiguous arm did not return the effect invocation ID")
	}
	if got := ambiguousRecording.effectIDs; len(got) != 1 || got[0] != indeterminate.InvocationID {
		t.Errorf("returned invocation ID %q, durable intents %v", indeterminate.InvocationID, got)
	}
	if ambiguousValidator.count() != 1 {
		t.Errorf("ambiguous arm validator requests = %d, want 1 (the body must have left)",
			ambiguousValidator.count())
	}
	if got := ambiguousSession.grants[0].Budget; got != 0 {
		t.Errorf("budget after an ambiguous dispatch = %d, want 0 (the attempt is consumed)", got)
	}
	ambiguousReceipt := effectReceipt(t, ambiguousRecording, 0)

	// THE CRITERION: the two receipt STATE STRINGS must differ, in this run.
	if definiteReceipt.State != store.ReceiptResolved {
		t.Errorf("definite arm receipt state = %q, want %q", definiteReceipt.State, store.ReceiptResolved)
	}
	if definiteReceipt.EffectOutcome == nil || definiteReceipt.EffectOutcome.Status != "failed" {
		t.Errorf("definite arm outcome = %+v, want a recorded \"failed\" outcome", definiteReceipt.EffectOutcome)
	}
	if ambiguousReceipt.State != store.ReceiptIndeterminate {
		t.Errorf("ambiguous arm receipt state = %q, want %q",
			ambiguousReceipt.State, store.ReceiptIndeterminate)
	}
	if ambiguousReceipt.EffectOutcome != nil {
		t.Errorf("ambiguous arm appended outcome %+v, want none", ambiguousReceipt.EffectOutcome)
	}
	if definiteReceipt.State == ambiguousReceipt.State {
		t.Fatalf("both arms reported receipt state %q: the criterion is that they DIFFER",
			definiteReceipt.State)
	}
}

// TestHandlerTimeoutIsAmbiguousButAnExitCodeIsNot pins the other half of the
// narrowness: a subprocess bound expiring mid-request is ambiguous, while an
// ordinary non-zero exit from the same handler is not.
func TestHandlerTimeoutIsAmbiguousButAnExitCodeIsNot(t *testing.T) {
	validator := newFakeValidator(t, "hang")
	fixture := newPublishFixture(t, validator.origin(), "ac11-timeout")
	handler := loopbackHandler(t, RegistryPublishConfig{
		PublisherPath:   writePublisherScript(t, "hang", nil),
		PackageDir:      fixture.dir,
		Manifest:        fixture.manifest,
		RegistryOrigin:  validator.origin(),
		ValidatorOrigin: validator.origin(),
		Credential:      writeCredentialFile(t, ac10Sentinel),
		Approval:        fixture.approval,
		ExecTimeout:     700 * time.Millisecond,
	})
	session, recording := publishSession(t, openTestStore(t), "ac11-timeout", Capability{
		Effect: EffectRegistryPublish, Scope: fixture.scope, ExpiresAt: 100, Budget: 1,
	}, handler)
	_, _, err := session.Invoke(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish, Scope: fixture.scope, Cost: PublishCost, Now: 50,
	}, fixture.payload)
	var indeterminate *IndeterminateEffectError
	if !errors.As(err, &indeterminate) {
		t.Fatalf("timed-out publish error = %T %v, want *IndeterminateEffectError", err, err)
	}
	if got := effectReceipt(t, recording, 0).State; got != store.ReceiptIndeterminate {
		t.Errorf("timed-out publish receipt state = %q, want %q", got, store.ReceiptIndeterminate)
	}
}

// TestOrdinaryHandlerFailuresKeepLandedResolvedBehaviour is the negative
// control for the narrow special case. It uses a NON-publish effect, so the
// only way it can red is if Session.Invoke's indeterminate arm stopped being
// keyed on the typed error — which is exactly MUT-SM-ALL-ERRORS-PENDING.
func TestOrdinaryHandlerFailuresKeepLandedResolvedBehaviour(t *testing.T) {
	base := openTestStore(t)
	recording := &publishRecordingStore{base: base}
	boom := errors.New("ordinary handler failure")
	session := newSession(recording, "narrowness", []Capability{{
		Effect: "probe", Scope: "s", ExpiresAt: 100, Budget: 3,
	}}, Registry{"probe": HandlerFunc(func(context.Context, EffectRequest, []byte) ([]byte, error) {
		return nil, boom
	})}, Live, nil)

	_, ref, err := session.Invoke(context.Background(), EffectRequest{
		Effect: "probe", Scope: "s", Cost: 1, Now: 1,
	}, []byte("payload"))
	var failed *EffectFailedError
	if !errors.As(err, &failed) {
		t.Fatalf("ordinary failure error = %T %v, want *EffectFailedError", err, err)
	}
	if !errors.Is(err, boom) {
		t.Errorf("ordinary failure lost its cause: %v", err)
	}
	if ref.IsZero() || len(recording.records) != 1 {
		t.Fatalf("ordinary failure produced ref %s and %d records, want a ref and 1 record",
			ref, len(recording.records))
	}
	receipt := effectReceipt(t, recording, 0)
	if receipt.State != store.ReceiptResolved {
		t.Fatalf("ordinary failure receipt state = %q, want %q", receipt.State, store.ReceiptResolved)
	}
	if receipt.EffectOutcome == nil || receipt.EffectOutcome.Status != "failed" {
		t.Fatalf("ordinary failure outcome = %+v, want \"failed\"", receipt.EffectOutcome)
	}
}

// ---------------------------------------------------------------------------
// AC10(a) — every production subprocess site, enumerated and DRIVEN
// ---------------------------------------------------------------------------

type subprocessSite struct {
	File string // repo-relative, slash-separated
	Line int
	Call string
}

func (s subprocessSite) String() string {
	return fmt.Sprintf("%s:%d %s", s.File, s.Line, s.Call)
}

// enumerateSubprocessSites re-derives N BY COMMAND (here: by parsing the same
// trees a grep would scan) in the run that uses it. Nothing about the count is
// transcribed from a document.
func enumerateSubprocessSites(t *testing.T, root string) []subprocessSite {
	t.Helper()
	var sites []subprocessSite
	fset := token.NewFileSet()
	for _, top := range []string{"host", "cmd"} {
		walkErr := filepath.WalkDir(filepath.Join(root, top), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			file, parseErr := parser.ParseFile(fset, path, nil, 0)
			if parseErr != nil {
				return fmt.Errorf("parse %s: %w", path, parseErr)
			}
			rel, relErr := filepath.Rel(root, path)
			if relErr != nil {
				return relErr
			}
			ast.Inspect(file, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				pkg, ok := sel.X.(*ast.Ident)
				if !ok || pkg.Name != "exec" {
					return true
				}
				if sel.Sel.Name != "Command" && sel.Sel.Name != "CommandContext" {
					return true
				}
				sites = append(sites, subprocessSite{
					File: filepath.ToSlash(rel),
					Line: fset.Position(call.Lparen).Line,
					Call: "exec." + sel.Sel.Name,
				})
				return true
			})
			return nil
		})
		if walkErr != nil {
			t.Fatal(walkErr)
		}
	}
	sort.Slice(sites, func(i, j int) bool {
		if sites[i].File != sites[j].File {
			return sites[i].File < sites[j].File
		}
		return sites[i].Line < sites[j].Line
	})
	return sites
}

func TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential(t *testing.T) {
	// Replace the process credential with a sentinel FIRST. Two reasons: a
	// leak observed below is then provably the sentinel and not a real key,
	// and the assertion becomes independent of whether the operator's shell
	// happens to carry one.
	t.Setenv(RegistryCredentialVariable, ac10Sentinel)

	root := testRepoRoot(t)
	sites := enumerateSubprocessSites(t, root)
	if len(sites) == 0 {
		// A zero-length enumeration is indistinguishable from a clean tree.
		t.Fatal("AC10(a): enumerated 0 production subprocess launch sites; the instrument is broken")
	}
	byFile := map[string][]subprocessSite{}
	for _, site := range sites {
		byFile[site.File] = append(byFile[site.File], site)
	}
	files := make([]string, 0, len(byFile))
	for file := range byFile {
		files = append(files, file)
	}
	sort.Strings(files)

	t.Logf("AC10(a) re-derived N = %d production subprocess launch sites in %d files:",
		len(sites), len(files))
	for _, site := range sites {
		t.Logf("    %s", site)
	}

	// KNOWN-POSITIVE CONTROL FOR THE INSTRUMENT. Before any absence is
	// believed, show that this exact dump mechanism DOES report the variable
	// when a child really is handed it.
	controlDump := filepath.Join(t.TempDir(), "control.env")
	if _, err := runBounded(context.Background(), handlerBounds{execTimeout: 20 * time.Second},
		handlerCommand{
			path: writeProbeScript(t, controlDump),
			dir:  t.TempDir(),
			env:  []string{"PATH=/usr/bin:/bin", RegistryCredentialVariable + "=" + ac10Sentinel},
		}); err != nil {
		t.Fatalf("known-positive control probe: %v", err)
	}
	control := readEnvDump(t, controlDump)
	if !childenv.Has(control, RegistryCredentialVariable) {
		t.Fatalf("known-positive control did NOT observe %s: every absence this test reports is void",
			RegistryCredentialVariable)
	}

	drivers := map[string]func(t *testing.T, probe string){
		"host/broker/handlers.go": driveBrokerDryRunPublish,
		"host/archive/archive.go": driveArchiveVersionProbe,
		"host/capsule/capsule.go": driveCapsuleRun,
		"host/pkgproj/pkgproj.go": drivePkgprojCrossCheck,
		"host/replay/replay.go":   driveReplayEntry,
	}
	if len(drivers) != len(files) {
		t.Fatalf("AC10(a): %d files carry subprocess sites %v but %d have drivers; "+
			"every enumerated site must be DRIVEN", len(files), files, len(drivers))
	}
	for _, file := range files {
		driver, ok := drivers[file]
		if !ok {
			t.Fatalf("AC10(a): subprocess site file %q has no driver", file)
		}
		t.Run(file, func(t *testing.T) {
			t.Setenv(RegistryCredentialVariable, ac10Sentinel)
			dump := filepath.Join(t.TempDir(), "child.env")
			driver(t, writeProbeScript(t, dump))
			observed := readEnvDump(t, dump)
			if len(observed) == 0 {
				t.Fatalf("%s: the child wrote a zero-length environment; nothing was measured", file)
			}
			if childenv.Has(observed, RegistryCredentialVariable) {
				// Name the VARIABLE, never the value.
				t.Fatalf("%s: the child observed %s", file, RegistryCredentialVariable)
			}
		})
	}
}

func driveBrokerDryRunPublish(t *testing.T, probe string) {
	t.Helper()
	validator := newFakeValidator(t, "ok")
	fixture := newPublishFixture(t, validator.origin(), "ac10-dryrun")
	handler := loopbackHandler(t, RegistryPublishConfig{
		PublisherPath:   probe,
		PackageDir:      fixture.dir,
		Manifest:        fixture.manifest,
		RegistryOrigin:  validator.origin(),
		ValidatorOrigin: validator.origin(),
		Approval:        fixture.approval,
		DryRun:          true,
		ExecTimeout:     20 * time.Second,
	})
	// The probe prints a version banner, not v0.30.0's success line, so this
	// classifies as a definite failure. The subprocess still ran, which is
	// what this driver exists to cause.
	_, _ = handler.Execute(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish, Scope: fixture.scope, Cost: PublishCost, Now: 50,
	}, fixture.payload)
	if got := handler.Dispatches(); got != 1 {
		t.Fatalf("dry-run dispatches = %d, want 1", got)
	}
	if got := handler.CredentialLoads(); got != 0 {
		t.Fatalf("dry-run consulted the credential provider %d times, want 0", got)
	}
	if got := validator.count(); got != 0 {
		t.Fatalf("dry-run reached the validator %d times, want 0", got)
	}
}

func driveArchiveVersionProbe(t *testing.T, probe string) {
	t.Helper()
	a := archive.New(filepath.Join(t.TempDir(), "world.db"))
	if _, err := a.Archive(probe); err != nil {
		t.Fatalf("archive.Archive: %v", err)
	}
}

func driveCapsuleRun(t *testing.T, probe string) {
	t.Helper()
	a := archive.New(filepath.Join(t.TempDir(), "world.db"))
	ref, err := a.Archive(probe)
	if err != nil {
		t.Fatalf("archive.Archive: %v", err)
	}
	// The probe is not an interpreter, so Run's output will not be a valid
	// transition result. The subprocess launch is the measurement.
	_, _ = capsule.New(a, capsule.Config{ExecTimeout: 20 * time.Second}).
		Run(capsule.Entry{Interpreter: ref, Source: []byte("module world/transition\n")})
}

func drivePkgprojCrossCheck(t *testing.T, probe string) {
	t.Helper()
	dir, manifest := publishFixtureDir(t)
	// CrossCheck fails on the probe's output, which is expected: the site it
	// launches is what is under test.
	_, _ = pkgproj.CrossCheck(dir, manifest, probe)
}

func driveReplayEntry(t *testing.T, probe string) {
	t.Helper()
	base := openTestStore(t)
	a := archive.New(filepath.Join(t.TempDir(), "world.db"))
	ref, err := a.Archive(probe)
	if err != nil {
		t.Fatalf("archive.Archive: %v", err)
	}
	source, err := replay.SourceObject([]byte("module world/transition\n"),
		"world/source/v1", "host/broker/registry_publish_test.go")
	if err != nil {
		t.Fatal(err)
	}
	if err := base.PutObject(source); err != nil {
		t.Fatal(err)
	}
	// The recorded result deliberately does not match the probe's banner, so
	// ReplayEntry returns a divergence — AFTER step 4 has launched the child.
	_, _ = replay.NewEngine(base, a).ReplayEntry(replay.Episode{
		WorldLibDir: t.TempDir(),
	}, 0, replay.EpisodeEntry{
		TransitionFn:      source.Hash,
		Interpreter:       ref,
		RecordedResult:    []byte("unused"),
		RecordedWorldHash: hashref.SumSHA256([]byte("unused")),
	})
}

// ---------------------------------------------------------------------------
// AC10(b) — redaction, paired with a non-secret marker injected identically
// ---------------------------------------------------------------------------

func TestPublisherErrorRedactsTheSecretAndKeepsTheMarker(t *testing.T) {
	validator := newFakeValidator(t, "validation")
	fixture := newPublishFixture(t, validator.origin(), "ac10b")
	handler := loopbackHandler(t, RegistryPublishConfig{
		PublisherPath:   writePublisherScript(t, "validation", nil),
		PackageDir:      fixture.dir,
		Manifest:        fixture.manifest,
		RegistryOrigin:  validator.origin(),
		ValidatorOrigin: validator.origin(),
		Credential:      writeCredentialFile(t, ac10Sentinel),
		Approval:        fixture.approval,
		ExecTimeout:     20 * time.Second,
	})
	session, recording := publishSession(t, openTestStore(t), "ac10b", Capability{
		Effect: EffectRegistryPublish, Scope: fixture.scope, ExpiresAt: 100, Budget: 1,
	}, handler)
	_, _, err := session.Invoke(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish, Scope: fixture.scope, Cost: PublishCost, Now: 50,
	}, fixture.payload)
	if err == nil {
		t.Fatal("expected the 400 arm to fail")
	}

	// The scanner is proven to be reading the stream by requiring the
	// NON-SECRET marker — injected through the identical path — to survive.
	// Without this pairing, "no sentinel" is satisfied by an empty message.
	var dispatch *PublishDispatchError
	if !errors.As(err, &dispatch) {
		t.Fatalf("error = %T %v, want a wrapped *PublishDispatchError", err, err)
	}
	streams := map[string]string{
		"error chain":      err.Error(),
		"publisher detail": dispatch.Detail,
	}
	for i, obj := range recording.records {
		streams[fmt.Sprintf("record[%d]", i)] = string(obj.Payload)
	}
	if !strings.Contains(dispatch.Detail, ac10Marker) {
		t.Fatalf("the non-secret marker %q did not survive to the publisher detail: the "+
			"redaction assertion below would pass vacuously. detail = %q",
			ac10Marker, dispatch.Detail)
	}
	if !strings.Contains(dispatch.Detail, redactedMarker) {
		t.Errorf("the publisher detail does not show the redaction marker %q: %q",
			redactedMarker, dispatch.Detail)
	}
	for name, text := range streams {
		if strings.Contains(text, ac10Sentinel) {
			t.Errorf("%s leaked the credential", name)
		}
	}

	// And the fake validator did receive the secret — otherwise there was
	// nothing to redact and the whole test is vacuous.
	if got := validator.requests(); len(got) != 1 || !got[0].HasAPIKey ||
		got[0].APIKeyLen != len(ac10Sentinel) {
		t.Fatalf("validator observed %+v, want exactly one request carrying a %d-byte key",
			got, len(ac10Sentinel))
	}
}

func TestRedactSecretLeavesAnEmptySecretAlone(t *testing.T) {
	const text = "abc"
	if got := redactSecret(text, nil); got != text {
		t.Errorf("redactSecret(%q, nil) = %q, want unchanged", text, got)
	}
	if got := redactSecret("x"+ac10Sentinel+"y", []byte(ac10Sentinel)); strings.Contains(got, ac10Sentinel) {
		t.Errorf("redactSecret did not remove the secret: %q", got)
	}
}

// ---------------------------------------------------------------------------
// the refusal set, the cost law, and the frozen payload
// ---------------------------------------------------------------------------

func TestPublishCostLawIsExactlyOne(t *testing.T) {
	// MUT-SM-COST-ZERO's anchor. The unit is one irreversible public write, so
	// a zero cost would make an attended budget of one authorize unboundedly
	// many of them.
	if PublishCost != 1 {
		t.Fatalf("PublishCost = %d, want exactly 1", PublishCost)
	}
	validator := newFakeValidator(t, "ok")
	fixture := newPublishFixture(t, validator.origin(), "cost")
	handler := loopbackHandler(t, RegistryPublishConfig{
		PublisherPath:   writePublisherScript(t, "success", nil),
		PackageDir:      fixture.dir,
		Manifest:        fixture.manifest,
		RegistryOrigin:  validator.origin(),
		ValidatorOrigin: validator.origin(),
		Credential:      writeCredentialFile(t, ac10Sentinel),
		Approval:        fixture.approval,
		ExecTimeout:     20 * time.Second,
	})
	for _, cost := range []int64{0, 2} {
		_, err := handler.Execute(context.Background(), EffectRequest{
			Effect: EffectRegistryPublish, Scope: fixture.scope, Cost: cost, Now: 50,
		}, fixture.payload)
		var refusal *PublishRefusalError
		if !errors.As(err, &refusal) {
			t.Fatalf("cost %d error = %T %v, want *PublishRefusalError", cost, err, err)
		}
	}
	if got := handler.Dispatches(); got != 0 {
		t.Fatalf("a refused cost still dispatched %d times", got)
	}

	// The record must carry cost 1 and debit exactly one unit.
	session, recording := publishSession(t, openTestStore(t), "cost", Capability{
		Effect: EffectRegistryPublish, Scope: fixture.scope, ExpiresAt: 100, Budget: 1,
	}, handler)
	if _, _, err := session.Invoke(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish, Scope: fixture.scope, Cost: PublishCost, Now: 50,
	}, fixture.payload); err != nil {
		t.Fatal(err)
	}
	rec, err := DecodeRecord(recording.records[0].Payload)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Cost != 1 || rec.BudgetBefore != 1 || rec.BudgetAfter != 0 {
		t.Fatalf("publish record cost/budget = %d/%d->%d, want 1/1->0",
			rec.Cost, rec.BudgetBefore, rec.BudgetAfter)
	}
}

func TestLoopbackConstructorRefusesEveryNonLoopbackOrigin(t *testing.T) {
	validator := newFakeValidator(t, "ok")
	fixture := newPublishFixture(t, validator.origin(), "origins")
	base := RegistryPublishConfig{
		PublisherPath:   writePublisherScript(t, "success", nil),
		PackageDir:      fixture.dir,
		Manifest:        fixture.manifest,
		RegistryOrigin:  validator.origin(),
		ValidatorOrigin: validator.origin(),
		Credential:      writeCredentialFile(t, ac10Sentinel),
		Approval:        fixture.approval,
	}
	// The unmodified configuration MUST construct, or every refusal below is
	// satisfied by an unrelated defect.
	if _, err := newLoopbackRegistryPublishHandler(base); err != nil {
		t.Fatalf("control configuration was refused: %v", err)
	}

	refusals := []struct {
		name     string
		mutate   func(*RegistryPublishConfig)
		contains string
	}{
		{"wildcard origin", func(c *RegistryPublishConfig) {
			c.ValidatorOrigin = "https://*.sunholo.com"
		}, "wildcard"},
		{"non-https live origin", func(c *RegistryPublishConfig) {
			c.ValidatorOrigin = "http://registry.ailang.sunholo.com"
		}, "loopback"},
		{"public origin through the test constructor", func(c *RegistryPublishConfig) {
			c.ValidatorOrigin = ApprovedValidatorOrigin
		}, "non-loopback"},
		{"credentials embedded in the URL", func(c *RegistryPublishConfig) {
			u, _ := url.Parse(validator.origin())
			c.ValidatorOrigin = u.Scheme + "://user:secret@" + u.Host
		}, "credentials in the URL"},
		{"query string on the origin", func(c *RegistryPublishConfig) {
			c.ValidatorOrigin = validator.origin() + "?token=x"
		}, "query or fragment"},
		{"--allow-dotted-tool-names", func(c *RegistryPublishConfig) {
			c.ExtraArgs = []string{"--allow-dotted-tool-names"}
		}, "--allow-dotted-tool-names"},
		{"a live handler without a credential provider", func(c *RegistryPublishConfig) {
			c.Credential = nil
		}, "credential provider"},
		{"a dry-run handler holding a credential provider", func(c *RegistryPublishConfig) {
			c.DryRun = true
		}, "no credential provider"},
	}
	for _, tc := range refusals {
		t.Run(tc.name, func(t *testing.T) {
			cfg := base
			tc.mutate(&cfg)
			_, err := newLoopbackRegistryPublishHandler(cfg)
			if err == nil {
				t.Fatalf("%s was accepted", tc.name)
			}
			if !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("refusal %q does not mention %q", err, tc.contains)
			}
		})
	}

	// The credentials-in-URL refusal must never echo the userinfo it refused.
	u, _ := url.Parse(validator.origin())
	cfg := base
	cfg.ValidatorOrigin = u.Scheme + "://user:s3cr3t-userinfo@" + u.Host
	_, err := newLoopbackRegistryPublishHandler(cfg)
	if err == nil || strings.Contains(err.Error(), "s3cr3t-userinfo") {
		t.Errorf("credentials-in-URL refusal echoed the userinfo: %v", err)
	}
}

func TestPublishHandlerRefusesAlternatePackageVersionAndHashes(t *testing.T) {
	validator := newFakeValidator(t, "ok")
	fixture := newPublishFixture(t, validator.origin(), "identity")
	handler := loopbackHandler(t, RegistryPublishConfig{
		PublisherPath:   writePublisherScript(t, "success", nil),
		PackageDir:      fixture.dir,
		Manifest:        fixture.manifest,
		RegistryOrigin:  validator.origin(),
		ValidatorOrigin: validator.origin(),
		Credential:      writeCredentialFile(t, ac10Sentinel),
		Approval:        fixture.approval,
		ExecTimeout:     20 * time.Second,
	})
	alternateVersion := fixture.identity
	alternateVersion.Version = "0.1.1"
	alternateName := fixture.identity
	alternateName.Name = "other"
	flippedHashes := fixture.hashes
	flippedHashes.TarballSHA256 = flipLastNibble(flippedHashes.TarballSHA256)

	cases := []struct {
		name     string
		scope    string
		payload  []byte
		contains string
	}{
		{"alternate version in the payload",
			PublishScope(validator.origin(), fixtureVendor, fixtureName, "0.1.1"),
			EncodePublishPayload(alternateVersion, fixture.hashes), "approval stamps"},
		{"alternate package name in the payload",
			PublishScope(validator.origin(), fixtureVendor, "other", fixtureVersion),
			EncodePublishPayload(alternateName, fixture.hashes), "approval stamps"},
		{"scope that does not describe the payload", "registry:x/package:y/version:z",
			fixture.payload, "does not describe the payload"},
		{"a flipped tarball nibble", fixture.scope,
			EncodePublishPayload(fixture.identity, flippedHashes), "tarball hash"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := handler.Execute(context.Background(), EffectRequest{
				Effect: EffectRegistryPublish, Scope: tc.scope, Cost: PublishCost, Now: 50,
			}, tc.payload)
			if err == nil {
				t.Fatalf("%s was accepted", tc.name)
			}
			if !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("refusal %q does not mention %q", err, tc.contains)
			}
		})
	}
	if got := handler.Dispatches(); got != 0 {
		t.Fatalf("identity refusals dispatched %d times, want 0", got)
	}
	if got := handler.CredentialLoads(); got != 0 {
		t.Fatalf("identity refusals loaded the credential %d times, want 0", got)
	}
}

func flipLastNibble(hash string) string {
	if hash == "" {
		return hash
	}
	last := hash[len(hash)-1]
	if last == '0' {
		return hash[:len(hash)-1] + "1"
	}
	return hash[:len(hash)-1] + "0"
}

func TestPublishPayloadFieldOrderIsFrozen(t *testing.T) {
	fixture := newPublishFixture(t, "http://127.0.0.1:1", "field-order")
	dec := json.NewDecoder(bytes.NewReader(fixture.payload))
	if _, err := dec.Token(); err != nil { // opening brace
		t.Fatal(err)
	}
	var keys []string
	for dec.More() {
		key, err := dec.Token()
		if err != nil {
			t.Fatal(err)
		}
		name, ok := key.(string)
		if !ok {
			t.Fatalf("object key %v is not a string", key)
		}
		keys = append(keys, name)
		var discard json.RawMessage
		if err := dec.Decode(&discard); err != nil {
			t.Fatal(err)
		}
	}
	if len(keys) == 0 {
		t.Fatal("decoded zero payload keys; the instrument is broken")
	}
	if !sameStrings(keys, publishPayloadFields) {
		t.Fatalf("payload field order = %v, want %v", keys, publishPayloadFields)
	}
}

// ---------------------------------------------------------------------------
// the credential provider and the ambient-startup refusal
// ---------------------------------------------------------------------------

func TestCredentialProviderRequiresModeAndLocation(t *testing.T) {
	root := testRepoRoot(t)
	outside := t.TempDir()

	good := filepath.Join(outside, "ok.key")
	if err := os.WriteFile(good, []byte(ac10Sentinel), 0o600); err != nil {
		t.Fatal(err)
	}
	provider, err := NewFileRegistryCredentialProvider(good, root)
	if err != nil {
		t.Fatalf("control credential file was refused: %v", err)
	}
	secret, err := provider.Credential()
	if err != nil || string(secret) != ac10Sentinel {
		t.Fatalf("Credential() = %d bytes, err %v, want the fixture secret", len(secret), err)
	}
	// The provider object must not carry the bytes.
	if strings.Contains(fmt.Sprintf("%+v", *provider), ac10Sentinel) {
		t.Error("the provider struct carries the secret")
	}

	loose := filepath.Join(outside, "loose.key")
	if err := os.WriteFile(loose, []byte(ac10Sentinel), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := NewFileRegistryCredentialProvider(loose, root); err == nil {
		t.Error("a mode-0644 credential file was accepted")
	}

	inside := filepath.Join(root, "credential-probe.key")
	if err := os.WriteFile(inside, []byte(ac10Sentinel), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(inside) })
	if _, err := NewFileRegistryCredentialProvider(inside, root); err == nil {
		t.Error("a credential file inside the working tree was accepted")
	}

	if _, err := NewFileRegistryCredentialProvider("relative.key", root); err == nil {
		t.Error("a relative credential path was accepted")
	}
}

func TestAmbientCredentialRefusalHasABothWaysControl(t *testing.T) {
	// MUT-SM-AMBIENT-KEY and its mandatory negative control: the SAME call,
	// differing only in the environment. Identical outcomes in both arms would
	// mean the check never read the environment.
	withKey := []string{"PATH=/usr/bin", RegistryCredentialVariable + "=" + ac10Sentinel}
	withoutKey := []string{"PATH=/usr/bin"}
	emptyKey := []string{"PATH=/usr/bin", RegistryCredentialVariable + "="}

	err := AssertNoAmbientRegistryCredential(withKey)
	if !errors.Is(err, ErrAmbientRegistryCredential) {
		t.Fatalf("ambient arm error = %v, want ErrAmbientRegistryCredential", err)
	}
	if strings.Contains(err.Error(), ac10Sentinel) {
		t.Fatal("the ambient-credential refusal printed the credential")
	}
	if !strings.Contains(err.Error(), RegistryCredentialVariable) {
		t.Errorf("the refusal %q does not name the variable", err)
	}
	if got := AssertNoAmbientRegistryCredential(withoutKey); got != nil {
		t.Fatalf("NEGATIVE CONTROL failed: absent variable produced %v", got)
	}
	if got := AssertNoAmbientRegistryCredential(emptyKey); got != nil {
		t.Fatalf("an empty assignment is not ambient authority, got %v", got)
	}
}

func TestProductionConstructorRefusesAnAmbientCredential(t *testing.T) {
	t.Setenv(RegistryCredentialVariable, ac10Sentinel)
	fixture := newPublishFixture(t, "https://storage.googleapis.com/ailang-registry", "prod")
	cfg := RegistryPublishConfig{
		PublisherPath:   filepath.Join(t.TempDir(), "ailang"),
		PackageDir:      fixture.dir,
		Manifest:        fixture.manifest,
		RegistryOrigin:  "https://storage.googleapis.com/ailang-registry",
		ValidatorOrigin: ApprovedValidatorOrigin,
		Credential:      writeCredentialFile(t, ac10Sentinel),
		Approval:        fixture.approval,
	}
	if _, err := NewRegistryPublishHandler(cfg); !errors.Is(err, ErrAmbientRegistryCredential) {
		t.Fatalf("ambient arm error = %v, want ErrAmbientRegistryCredential", err)
	}
	// NEGATIVE CONTROL: the same construction with the variable absent must
	// get PAST the ambient gate. It still fails on the missing publisher
	// binary, which is a DIFFERENT, named error — the two arms differ.
	t.Setenv(RegistryCredentialVariable, "")
	_, err := NewRegistryPublishHandler(cfg)
	if errors.Is(err, ErrAmbientRegistryCredential) {
		t.Fatal("NEGATIVE CONTROL failed: the ambient gate fired with the variable absent")
	}
	if err != nil {
		t.Logf("negative control passed the ambient gate and stopped later: %v", err)
	}
}

func TestProductionConstructorAcceptsOnlyTheApprovedValidatorOrigin(t *testing.T) {
	t.Setenv(RegistryCredentialVariable, "")
	validator := newFakeValidator(t, "ok")
	fixture := newPublishFixture(t, validator.origin(), "prod-origin")
	_, err := NewRegistryPublishHandler(RegistryPublishConfig{
		PublisherPath:   filepath.Join(t.TempDir(), "ailang"),
		PackageDir:      fixture.dir,
		Manifest:        fixture.manifest,
		RegistryOrigin:  validator.origin(),
		ValidatorOrigin: validator.origin(),
		Credential:      writeCredentialFile(t, ac10Sentinel),
		Approval:        fixture.approval,
	})
	if err == nil || !strings.Contains(err.Error(), ApprovedValidatorOrigin) {
		t.Fatalf("production constructor accepted a loopback validator: %v", err)
	}
}

// ---------------------------------------------------------------------------
// network evidence
// ---------------------------------------------------------------------------

func TestRegistryEnvironmentPointsOnlyAtLoopback(t *testing.T) {
	for _, name := range []string{"AILANG_REGISTRY", "AILANG_REGISTRY_VALIDATOR", "AILANG_REGISTRY_API"} {
		value, present := os.LookupEnv(name)
		if !present || value == "" {
			t.Logf("%s is unset for this run", name)
			continue
		}
		parsed, err := url.Parse(value)
		if err != nil || !isLoopbackHost(parsed.Hostname()) {
			t.Fatalf("%s=%q is not loopback: this test run could reach a real registry", name, value)
		}
		t.Logf("%s points at loopback host %q", name, parsed.Hostname())
	}
}

func TestFakeValidatorSawOnlyLoopbackTraffic(t *testing.T) {
	validator := newFakeValidator(t, "ok")
	if !isLoopbackHost(mustHost(t, validator.origin())) {
		t.Fatalf("the fake validator bound %q, which is not loopback", validator.origin())
	}
	fixture := newPublishFixture(t, validator.origin(), "loopback-log")
	handler := loopbackHandler(t, RegistryPublishConfig{
		PublisherPath:   writePublisherScript(t, "success", nil),
		PackageDir:      fixture.dir,
		Manifest:        fixture.manifest,
		RegistryOrigin:  validator.origin(),
		ValidatorOrigin: validator.origin(),
		Credential:      writeCredentialFile(t, ac10Sentinel),
		Approval:        fixture.approval,
		ExecTimeout:     20 * time.Second,
	})
	if _, err := handler.Execute(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish, Scope: fixture.scope, Cost: PublishCost, Now: 50,
	}, fixture.payload); err != nil {
		t.Fatal(err)
	}
	log := validator.requests()
	if len(log) == 0 {
		t.Fatal("the validator log is empty; a request log that never fills proves nothing")
	}
	for i, req := range log {
		t.Logf("validator request %d: %s %s host=%s remote=%s apiKey=%v(%d bytes)",
			i, req.Method, req.Path, req.Host, req.RemoteAddr, req.HasAPIKey, req.APIKeyLen)
		if !isLoopbackHost(mustHost(t, "http://"+req.Host)) {
			t.Errorf("validator request %d carried non-loopback Host %q", i, req.Host)
		}
	}
}

func mustHost(t *testing.T, raw string) string {
	t.Helper()
	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("parse %q: %v", raw, err)
	}
	return parsed.Hostname()
}
