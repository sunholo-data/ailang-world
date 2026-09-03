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

	// The three fields below are populated only by landApproval. Until then
	// identity.ApprovalRef is a SYNTHETIC digest that names no object, which is
	// exactly right for the handler-only tests (they never reach the broker's
	// approval traversal) and exactly wrong for anything routed through
	// Session.Invoke.
	approvalScope      string
	approvalRequestRef hashref.HashRef
	approvalDecision   string
}

// ---------------------------------------------------------------------------
// SM.B2b — a FILE-BACKED store, because AC9a and AC9c must close and reopen it
//
// openTestStore uses ":memory:", which host/store documents as per-connection
// and physically unreachable from a second handle. A close/reopen criterion
// evaluated against it would be satisfied by an empty database — i.e. it would
// pass for the wrong reason. Every SM.B2b test therefore uses this instead.
// ---------------------------------------------------------------------------

type publishStore struct {
	*store.Store
	path   string
	closed bool
}

func openPublishStore(t *testing.T) *publishStore {
	t.Helper()
	return openPublishStoreAt(t, filepath.Join(t.TempDir(), "world.db"))
}

func openPublishStoreAt(t *testing.T, path string) *publishStore {
	t.Helper()
	s, err := store.Open(path)
	if err != nil {
		t.Fatalf("store.Open(%q): %v", path, err)
	}
	p := &publishStore{Store: s, path: path}
	t.Cleanup(func() { p.close(t) })
	return p
}

func (p *publishStore) close(t *testing.T) {
	t.Helper()
	if p.closed {
		return
	}
	p.closed = true
	if err := p.Store.Close(); err != nil {
		t.Fatalf("close store %q: %v", p.path, err)
	}
}

// reopen closes this handle and opens the SAME FILE again. The returned handle
// shares nothing but the bytes on disk: no cache, no connection, no in-memory
// budget. That is the whole content of AC9a's "new process".
func (p *publishStore) reopen(t *testing.T) *publishStore {
	t.Helper()
	p.close(t)
	return openPublishStoreAt(t, p.path)
}

// approvalTimes are the logical times of one attended approval. They are
// explicit because the broker never reads a wall clock and every ordering
// refusal in validatePublishApproval is stated in terms of them.
type approvalTimes struct{ request, decide, expires int64 }

func defaultApprovalTimes() approvalTimes {
	return approvalTimes{request: 10, decide: 11, expires: 100}
}

// landApproval mints the attended stamp for f through the LANDED attended
// path and returns a fixture whose ApprovalRef IS the content hash of the
// resulting immutable ApprovalDecisionV1.
//
// Every step is the landed surface, unmodified:
//
//	Session.Invoke(Human.Approve)  -> HumanHandler mints ApprovalRequestV1
//	DecideApproval(...)            -> the operator entry point mints ApprovalDecisionV1
//	Session.Invoke(Human.PollApproval) -> HumanHandler observes it back
//
// The poll is not decoration: AC9's non-vacuity requirement is that the
// positive control be an approval TRAVERSED through DecideApproval and
// EffectHumanPollApproval, so the traversal is performed here for every
// fixture and its result is checked to hash back to the same decision ref.
func (f publishFixture) landApproval(
	t *testing.T, base approvalStore, decision string, times approvalTimes,
) publishFixture {
	t.Helper()
	return f.landApprovalWithScope(t, base,
		PublishApprovalScopeFor(f.identity, f.hashes, times.expires), decision, times)
}

// landApprovalWithScope is landApproval with the canonical scope supplied by
// the caller instead of derived. It exists for exactly one reason: AC9's
// wrong-effect arm needs a stamp that is publish-grammatical in every respect
// EXCEPT its frozen `effect` term, and PublishApprovalScopeFor — correctly —
// cannot mint one. It reaches for FormatPublishApprovalScope, which is landed
// exported production code, so no codec is widened to build the fixture.
func (f publishFixture) landApprovalWithScope(
	t *testing.T, base approvalStore, approvalScope, decision string, times approvalTimes,
) publishFixture {
	t.Helper()
	return f.landApprovalWithScopeAndCost(t, base, approvalScope, decision, times, PublishCost)
}

// landApprovalWithScopeAndCost additionally lets the caller choose the cost the
// attended request is minted AT. It exists for AC9's wrong-cost arm: the landed
// HumanHandler copies req.Cost into approvalRequestWire.Cost verbatim, so a
// request priced at anything other than PublishCost is minted by invoking
// Human.Approve at that cost — no codec is touched to produce it.
func (f publishFixture) landApprovalWithScopeAndCost(
	t *testing.T, base approvalStore, approvalScope, decision string,
	times approvalTimes, cost int64,
) publishFixture {
	t.Helper()
	human := newHumanHandler(base)
	session := newSession(base, "attended-"+f.identity.Vendor+"-"+f.identity.Version,
		[]Capability{
			{Effect: EffectHumanApprove, Scope: approvalScope, ExpiresAt: times.expires, Budget: cost},
			{Effect: EffectHumanPollApproval, Scope: approvalScope, ExpiresAt: times.expires, Budget: 4},
		}, Registry{
			EffectHumanApprove: human, EffectHumanPollApproval: human,
		}, Live, nil)

	pending, _, err := session.Invoke(context.Background(), EffectRequest{
		Effect: EffectHumanApprove, Scope: approvalScope, Cost: cost, Now: times.request,
	}, mustApprovalJSON(approvalInputWire{Requester: "sm-b2b-fixture"}))
	if err != nil {
		t.Fatalf("landed Human.Approve: %v", err)
	}
	requestRef := decodePendingRef(t, pending)

	decisionRef := decideAndPollLandedApproval(t, base, approvalScope, requestRef, decision, times)

	landed := f
	landed.identity.ApprovalRef = decisionRef
	landed.approval.ApprovalRef = decisionRef
	landed.payload = EncodePublishPayload(landed.identity, landed.hashes)
	landed.approvalScope = approvalScope
	landed.approvalRequestRef = requestRef
	landed.approvalDecision = decision
	return landed
}

// decideAndPollLandedApproval runs the second and third legs of the attended
// path — the landed operator entry point and the landed poll effect — and
// asserts that what comes back out of the poll hashes to the decision the
// publish payload will name. Without that check the traversal could be
// observing some other object and every "landed" claim below would be a
// coincidence.
func decideAndPollLandedApproval(
	t *testing.T, base approvalStore, approvalScope string,
	requestRef hashref.HashRef, decision string, times approvalTimes,
) hashref.HashRef {
	t.Helper()
	decisionRef, err := decideApproval(base, requestRef, decision, "attended-operator", times.decide)
	if err != nil {
		t.Fatalf("landed DecideApproval(%q): %v", decision, err)
	}
	human := newHumanHandler(base)
	pollSession := newSession(base, "attended-poll", []Capability{
		{Effect: EffectHumanPollApproval, Scope: approvalScope, ExpiresAt: times.expires, Budget: 4},
	}, Registry{EffectHumanPollApproval: human}, Live, nil)
	polled, _, err := pollSession.Invoke(context.Background(), EffectRequest{
		Effect: EffectHumanPollApproval, Scope: approvalScope, Cost: PublishCost, Now: times.decide,
	}, mustApprovalJSON(approvalInputWire{RequestRef: requestRef.String()}))
	if err != nil {
		t.Fatalf("landed Human.PollApproval: %v", err)
	}
	var observed observedDecisionWire
	if err := decodeApprovalJSON(polled, &observed); err != nil {
		t.Fatal(err)
	}
	if observed.Status != "decided" {
		t.Fatalf("landed poll status = %q, want \"decided\"", observed.Status)
	}
	if got := hashref.SumSHA256(observed.Decision); got != decisionRef {
		t.Fatalf("polled decision hashes to %s, want the decision ref %s", got, decisionRef)
	}
	return decisionRef
}

// landApprovalOverRawRequest mints the ApprovalRequestV1 DIRECTLY, then heads,
// decides and polls it through the landed surface.
//
// It exists for exactly one arm, and it is the only fixture in this file that
// does not go through the Human.Approve effect — NECESSARILY so. The guard it
// pins refuses a request that was not minted by Human.Approve, and
// HumanHandler.Execute writes `Effect: req.Effect` inside `case
// EffectHumanApprove:`, so the landed effect can only ever produce
// "Human.Approve". A fixture the guard is supposed to reject therefore cannot
// be produced by the path the guard exists to bless; the alternative to
// building it here is leaving the branch untested, which is what this arm is
// repairing.
//
// Everything it reaches for is landed production code — brokerObject,
// mustApprovalJSON, the unchanged approvalRequestWire, appendApprovalHead,
// decideApproval and the poll effect. No codec is widened and no wire type is
// added to build it.
func (f publishFixture) landApprovalOverRawRequest(
	t *testing.T, base approvalStore, wire approvalRequestWire, decision string, times approvalTimes,
) publishFixture {
	t.Helper()
	requestObj := brokerObject(ApprovalRequestV1, mustApprovalJSON(wire))
	if err := base.PutObject(requestObj); err != nil {
		t.Fatalf("put raw approval request: %v", err)
	}
	if err := appendApprovalHead(base, requestObj.Hash, hashref.HashRef{}); err != nil {
		t.Fatalf("head raw approval request: %v", err)
	}
	decisionRef := decideAndPollLandedApproval(
		t, base, wire.Scope, requestObj.Hash, decision, times)

	landed := f
	landed.identity.ApprovalRef = decisionRef
	landed.approval.ApprovalRef = decisionRef
	landed.payload = EncodePublishPayload(landed.identity, landed.hashes)
	landed.approvalScope = wire.Scope
	landed.approvalRequestRef = requestObj.Hash
	landed.approvalDecision = decision
	return landed
}

// putRawApprovalObject stores an approval-shaped object with caller-chosen
// payload bytes and returns its content hash. It is how AC9 reaches the SIX
// traversal-error branches of validatePublishApproval — an unparseable
// requestRef, a decision naming an absent or wrong-kind request, undecodable
// request or decision bytes — none of which the landed minting path can
// produce, because the landed minting path is the thing that produces
// well-formed objects. brokerObject is landed production code, so the fixture
// still costs no codec change.
func putRawApprovalObject(t *testing.T, base approvalStore, semanticID string, payload []byte) hashref.HashRef {
	t.Helper()
	obj := brokerObject(semanticID, payload)
	if err := base.PutObject(obj); err != nil {
		t.Fatalf("put raw %s object: %v", semanticID, err)
	}
	return obj.Hash
}

// withApprovalRef re-points a payload at a different approval decision, leaving
// every other field of the identity alone.
func (f publishFixture) withApprovalRef(ref hashref.HashRef) []byte {
	id := f.identity
	id.ApprovalRef = ref
	return EncodePublishPayload(id, f.hashes)
}

// publishGrant is the one-shot attended capability: exact scope, budget one.
func publishGrant(scope string) Capability {
	return Capability{Effect: EffectRegistryPublish, Scope: scope, ExpiresAt: 100, Budget: PublishCost}
}

// publishCounters is the observable AC8/AC9/AC9a/AC9b/AC9c are stated in: the
// number of POSTs the SHARED fake validator saw, the number of publisher
// subprocess dispatches, and the number of times the credential provider was
// consulted. They are read together so an arm can assert that a refusal moved
// NONE of them.
type publishCounters struct {
	posts, dispatches, credentialLoads int64
}

func readPublishCounters(v *fakeValidator, h *RegistryPublishHandler) publishCounters {
	return publishCounters{
		posts: int64(v.count()), dispatches: h.Dispatches(), credentialLoads: h.CredentialLoads(),
	}
}

func (c publishCounters) String() string {
	return fmt.Sprintf("POST=%d dispatches=%d credentialLoads=%d",
		c.posts, c.dispatches, c.credentialLoads)
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

func (s *publishRecordingStore) GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error) {
	return s.base.GetObject(ctx, ref)
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
	base := openPublishStore(t)
	// SM.B2b: the positive control now runs through validatePublishApproval, so
	// its approval must be a LANDED ApprovalDecisionV1 rather than a synthetic
	// digest. The denial arms below are unaffected — they are refused by the
	// capability decision, which is strictly earlier.
	fixture := newPublishFixture(t, validator.origin(), "ac7").
		landApproval(t, base, "approve", defaultApprovalTimes())
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
			if _, ok, getErr := recording.base.GetObject(context.Background(), ref); getErr != nil || !ok {
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
	// It runs on `base`, the store the attended approval actually landed in.
	session, recording := publishSession(t, base.Store, "ac7-allowed", validGrant, handler)
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
	base := openPublishStore(t)
	grant := publishGrant

	// ARM A — a DEFINITE handler failure: the validator answers 403, so the
	// attempt is over and its result is known.
	definiteValidator := newFakeValidator(t, "namespace")
	definiteFixture := newPublishFixture(t, definiteValidator.origin(), "ac11-definite").
		landApproval(t, base, "approve", defaultApprovalTimes())
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
		t, base.Store, "ac11-definite", grant(definiteFixture.scope), definiteHandler)
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
	ambiguousFixture := newPublishFixture(t, ambiguousValidator.origin(), "ac11-ambiguous").
		landApproval(t, base, "approve", defaultApprovalTimes())
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
		t, base.Store, "ac11-ambiguous", grant(ambiguousFixture.scope), ambiguousHandler)
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
	base := openPublishStore(t)
	fixture := newPublishFixture(t, validator.origin(), "ac11-timeout").
		landApproval(t, base, "approve", defaultApprovalTimes())
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
	session, recording := publishSession(t, base.Store, "ac11-timeout", publishGrant(fixture.scope), handler)
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
		t.Fatalf("%s", archive.AttributeFailure("archive.Archive", err))
	}
}

func driveCapsuleRun(t *testing.T, probe string) {
	t.Helper()
	a := archive.New(filepath.Join(t.TempDir(), "world.db"))
	ref, err := a.Archive(probe)
	if err != nil {
		t.Fatalf("%s", archive.AttributeFailure("archive.Archive", err))
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
		t.Fatalf("%s", archive.AttributeFailure("archive.Archive", err))
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
	base := openPublishStore(t)
	fixture := newPublishFixture(t, validator.origin(), "ac10b").
		landApproval(t, base, "approve", defaultApprovalTimes())
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
	session, recording := publishSession(t, base.Store, "ac10b", publishGrant(fixture.scope), handler)
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
	base := openPublishStore(t)
	fixture := newPublishFixture(t, validator.origin(), "cost").
		landApproval(t, base, "approve", defaultApprovalTimes())
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
	session, recording := publishSession(t, base.Store, "cost", publishGrant(fixture.scope), handler)
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

// ---------------------------------------------------------------------------
// SM.B2b — AC8, AC9, AC9a, AC9b, AC9c
//
// Everything below shares one shape, because the criteria share one claim:
// that single use of the attended stamp is enforced DURABLY, by the store, and
// not by the in-memory budget that happens to sit in front of it. So every
// reuse arm is presented with a session whose budget is FRESH — a session that
// would happily allow the publish if the claim were the only thing stopping it.
// ---------------------------------------------------------------------------

// landNonPublishApproval mints an approval whose scope is NOT publish-shaped,
// through the same landed path. It is the fixture for AC9's malformed arm and
// for approve_test.go's backwards-compatibility arms.
func landNonPublishApproval(
	t *testing.T, base approvalStore, scope string, cost int64, times approvalTimes,
) (requestRef, decisionRef hashref.HashRef) {
	t.Helper()
	human := newHumanHandler(base)
	session := newSession(base, "attended-nonpublish-"+scope, []Capability{
		{Effect: EffectHumanApprove, Scope: scope, ExpiresAt: times.expires, Budget: cost},
	}, Registry{EffectHumanApprove: human}, Live, nil)
	pending, _, err := session.Invoke(context.Background(), EffectRequest{
		Effect: EffectHumanApprove, Scope: scope, Cost: cost, Now: times.request,
	}, mustApprovalJSON(approvalInputWire{Requester: "sm-b2b-fixture"}))
	if err != nil {
		t.Fatalf("non-publish Human.Approve: %v", err)
	}
	requestRef = decodePendingRef(t, pending)
	decisionRef, err = decideApproval(base, requestRef, "approve", "attended-operator", times.decide)
	if err != nil {
		t.Fatalf("non-publish DecideApproval: %v", err)
	}
	return requestRef, decisionRef
}

// AC8 — the claim and the intent are persisted atomically BEFORE dispatch, the
// budget lands at zero, and a FRESH session with a FRESH budget re-presenting
// the same approval is refused before the credential is read and before POST.
func TestPublishClaimIsDurableBeforeDispatchAndSurvivesAFreshBudget(t *testing.T) {
	validator := newFakeValidator(t, "ok")
	base := openPublishStore(t)
	fixture := newPublishFixture(t, validator.origin(), "ac8").
		landApproval(t, base, "approve", defaultApprovalTimes())
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

	// The pre-dispatch observation. It runs INSIDE the handler chain, so
	// "before dispatch" is a position in the real control flow rather than an
	// inference from ordering in the source.
	type preDispatch struct {
		ran            bool
		effectIDs      []string
		receiptState   store.ReceiptState
		hasReceipt     bool
		claimReuseErr  error
		countersAtGate publishCounters
	}
	var observed preDispatch

	recording := &publishRecordingStore{base: base.Store}
	probe := HandlerFunc(func(ctx context.Context, req EffectRequest, payload []byte) ([]byte, error) {
		observed.ran = true
		observed.effectIDs = append([]string(nil), recording.effectIDs...)
		observed.countersAtGate = readPublishCounters(validator, handler)
		if len(recording.effectIDs) == 1 {
			receipt, ok, err := base.GetEffectReceipt(recording.effectIDs[0])
			observed.hasReceipt, observed.receiptState = ok, receipt.State
			if err != nil {
				observed.claimReuseErr = err
				return nil, err
			}
		}
		// Ask the STORE whether the approval is already claimed, using the same
		// production entry point a second session would use. A rollback leaves
		// nothing behind, so this probe cannot itself consume anything.
		probeRef := hashref.SumSHA256([]byte("ac8-claim-probe"))
		_, _, observed.claimReuseErr = base.AppendClaimedEffectIntent(
			"ac8-claim-probe",
			store.EffectIntent{
				EpisodeID: "ac8-claim-probe", Effect: EffectRegistryPublish,
				Scope: fixture.scope, Cost: PublishCost, RequestRef: probeRef, LogicalTime: 51,
			}, fixture.identity.ApprovalRef, probeRef)
		return handler.Execute(ctx, req, payload)
	})
	session := newSession(recording, "ac8", []Capability{publishGrant(fixture.scope)},
		Registry{EffectRegistryPublish: probe}, Live, nil)

	if _, ref, err := session.Invoke(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish, Scope: fixture.scope, Cost: PublishCost, Now: 50,
	}, fixture.payload); err != nil || ref.IsZero() {
		t.Fatalf("AC8 grant Invoke = ref %s, err %v; want a record and no error "+
			"(counters at refusal: %s)", ref, err, readPublishCounters(validator, handler))
	}

	if !observed.ran {
		t.Fatal("AC8: the pre-dispatch probe never ran; every observation below is void")
	}
	if len(observed.effectIDs) != 1 {
		t.Fatalf("AC8: durable effect intents at dispatch = %v, want exactly 1", observed.effectIDs)
	}
	if !observed.hasReceipt || observed.receiptState != store.ReceiptIndeterminate {
		t.Errorf("AC8: receipt at dispatch = %q (present %v), want %q — the intent must be durable "+
			"BEFORE the POST", observed.receiptState, observed.hasReceipt, store.ReceiptIndeterminate)
	}
	if !errors.Is(observed.claimReuseErr, store.ErrApprovalAlreadyConsumed) {
		t.Errorf("AC8: re-claiming the approval at dispatch time = %v, want ErrApprovalAlreadyConsumed "+
			"— the claim must be durable BEFORE the POST", observed.claimReuseErr)
	}
	if observed.countersAtGate.posts != 0 || observed.countersAtGate.dispatches != 0 {
		t.Errorf("AC8: counters at the pre-dispatch gate = %s, want POST=0 dispatches=0",
			observed.countersAtGate)
	}
	if got := session.grants[0].Budget; got != 0 {
		t.Errorf("AC8: budget after the grant = %d, want 0", got)
	}
	afterGrant := readPublishCounters(validator, handler)
	t.Logf("AC8 counters after the valid grant: %s", afterGrant)
	if afterGrant != (publishCounters{posts: 1, dispatches: 1, credentialLoads: 1}) {
		t.Fatalf("AC8: counters after the valid grant = %s, want POST=1 dispatches=1 credentialLoads=1",
			afterGrant)
	}

	// THE NON-VACUITY REQUIREMENT: a FRESH session with a FRESH budget. Without
	// the fresh budget this arm would be satisfied by the in-memory debit that
	// AC8 exists to say is not enough.
	freshSession, _ := publishSession(t, base.Store, "ac8-fresh", publishGrant(fixture.scope), handler)
	if got := freshSession.grants[0].Budget; got != PublishCost {
		t.Fatalf("AC8: the fresh session's budget = %d, want %d — a spent budget would refuse for "+
			"the WRONG reason", got, PublishCost)
	}
	_, _, reuseErr := freshSession.Invoke(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish, Scope: fixture.scope, Cost: PublishCost, Now: 51,
	}, fixture.payload)
	if !errors.Is(reuseErr, store.ErrApprovalAlreadyConsumed) {
		t.Fatalf("AC8: fresh-budget reuse error = %T %v, want store.ErrApprovalAlreadyConsumed",
			reuseErr, reuseErr)
	}
	var denied *DenialError
	if errors.As(reuseErr, &denied) {
		t.Fatalf("AC8: the reuse was refused by BUDGET (%s), not by the durable claim", denied.Decision.Label)
	}
	afterReuse := readPublishCounters(validator, handler)
	t.Logf("AC8 counters after the fresh-budget reuse: %s", afterReuse)
	if afterReuse != afterGrant {
		t.Fatalf("AC8: counters moved across the refused reuse: %s -> %s; the refusal must land "+
			"BEFORE credential load and POST", afterGrant, afterReuse)
	}
}

// AC9 — the seven refusal classes, each with its own sentinel, each measured to
// leave the credential provider and the POST counter untouched; and a positive
// control that is an approval traversed through the LANDED DecideApproval and
// EffectHumanPollApproval.
func TestPublishApprovalRefusalSetWithALandedPositiveControl(t *testing.T) {
	validator := newFakeValidator(t, "ok")
	base := openPublishStore(t)
	good := newPublishFixture(t, validator.origin(), "ac9").
		landApproval(t, base, "approve", defaultApprovalTimes())
	handler := loopbackHandler(t, RegistryPublishConfig{
		PublisherPath:   writePublisherScript(t, "success", nil),
		PackageDir:      good.dir,
		Manifest:        good.manifest,
		RegistryOrigin:  validator.origin(),
		ValidatorOrigin: validator.origin(),
		Credential:      writeCredentialFile(t, ac10Sentinel),
		Approval:        good.approval,
		ExecTimeout:     20 * time.Second,
	})

	// ARM: expired. Same bytes, same package, stamped to expire at 40; the
	// publish is presented at 50.
	expired := newPublishFixture(t, validator.origin(), "ac9-expired").
		landApproval(t, base, "approve", approvalTimes{request: 5, decide: 6, expires: 40})

	// ARM: wrong-scope. The stamp is minted for 0.1.1; the payload publishes
	// 0.1.0. Everything else about it is valid.
	otherVersion := good
	otherVersion.identity.Version = "0.1.1"
	otherVersion = otherVersion.landApproval(t, base, "approve",
		approvalTimes{request: 7, decide: 8, expires: 100})
	wrongScopeIdentity := good.identity
	wrongScopeIdentity.ApprovalRef = otherVersion.identity.ApprovalRef

	// ARM: wrong-hash. The stamp is the valid one; the payload's tarball digest
	// differs from it by exactly one nibble.
	flipped := good.hashes
	flipped.TarballSHA256 = flipLastNibble(flipped.TarballSHA256)

	// ARM: denied.
	denied := newPublishFixture(t, validator.origin(), "ac9").
		landApproval(t, base, "deny", approvalTimes{request: 12, decide: 13, expires: 100})

	// ARM: wrong-effect. A GENUINELY LANDED approval — real Human.Approve, real
	// DecideApproval, real Human.PollApproval — whose scope is well-formed
	// publish-approval grammar in every single respect except that its frozen
	// `effect` term names FS.Write.
	//
	// It exists because nothing else reaches the effect check. The "scope is
	// not publish grammar" arm below is refused by ParsePublishApprovalScope
	// for carrying no mark, which returns BEFORE scope.Effect is ever read — so
	// that arm pins the parser, not the term. Everything here is deliberately
	// identical to the good fixture (same package, version, manifest ref and
	// all three digests) so that the ONLY thing that can refuse it is the
	// effect term: with that one check disabled, this payload passes every
	// remaining comparison in validatePublishApproval.
	wrongEffect := good.landApprovalWithScope(t, base, FormatPublishApprovalScope(PublishApprovalScope{
		Publish:       PublishScope(validator.origin(), fixtureVendor, fixtureName, fixtureVersion),
		Effect:        "FS.Write",
		ManifestRef:   good.identity.ManifestRef.String(),
		TarballSHA256: good.hashes.TarballSHA256,
		ContentHash:   good.hashes.ContentHash,
		InterfaceHash: good.hashes.InterfaceHash,
		ExpiresAt:     100,
	}), "approve", approvalTimes{request: 16, decide: 17, expires: 100})

	// goodScopeAt renders the good fixture's canonical scope with a chosen
	// expiry term. Every arm below that is not ABOUT the scope uses the
	// expiry-100 form, so the field under test is the only difference.
	goodScopeAt := func(expiresAt int64) string {
		return PublishApprovalScopeFor(good.identity, good.hashes, expiresAt)
	}

	// ARM: wrong-request-effect. The ApprovalRequestV1 was not minted by
	// Human.Approve. See landApprovalOverRawRequest for why this is the one
	// fixture that cannot come out of the landed effect: the guard exists to
	// refuse exactly what that effect cannot produce.
	wrongRequestEffect := good.landApprovalOverRawRequest(t, base, approvalRequestWire{
		Effect: EffectHumanPollApproval, Scope: goodScopeAt(100), Cost: PublishCost,
		Requester: "sm-b2b-fixture", Now: 18,
	}, "approve", approvalTimes{request: 18, decide: 19, expires: 100})

	// ARM: wrong-cost. Minted through the landed effect AT cost 2, which
	// HumanHandler copies into the request wire verbatim.
	wrongCost := good.landApprovalWithScopeAndCost(t, base, goodScopeAt(100), "approve",
		approvalTimes{request: 22, decide: 23, expires: 100}, 2)

	// ARM: request-after-publish. The attended request is made at logical time
	// 60; the publish it is presented for is at 50. An approval requested AFTER
	// the thing it authorizes is not evidence of anything.
	//
	// This arm and the next share ErrPublishApprovalMalformed and are pinned by
	// their MESSAGE as well as their sentinel, deliberately: [request 60,
	// publish 50] is an EMPTY interval, so any decision time whatsoever also
	// violates the decision-range branch. There is no fixture that violates
	// request-after-publish alone, so the substring is what makes each arm
	// attributable to its own branch.
	requestAfterPublish := good.landApprovalWithScope(t, base, goodScopeAt(100), "approve",
		approvalTimes{request: 60, decide: 61, expires: 100})

	// ARM: decision-outside-range. Requested at 20, DECIDED at 5 — before the
	// request it decides.
	decisionBeforeRequest := good.landApprovalWithScope(t, base, goodScopeAt(100), "approve",
		approvalTimes{request: 20, decide: 5, expires: 100})

	// ARM: expiry-before-its-own-request. The stamp expires at 15 but was
	// requested at 20, so it was spent authority the moment it was minted. The
	// capability that minted it is unexpired (times.expires = 100), which is
	// what keeps this arm about the SCOPE's expiry term rather than about the
	// grant.
	expiryBeforeRequest := good.landApprovalWithScope(t, base, goodScopeAt(15), "approve",
		approvalTimes{request: 20, decide: 21, expires: 100})

	// ARM: grant-scope-mismatch. The stamp and the payload agree completely;
	// the CAPABILITY the effect was decided against describes a different
	// package version. This is the third of the three identities
	// validatePublishApproval keeps separate, and it is the one the existing
	// wrong-scope arm does NOT exercise (there, stamp and grant agree and the
	// payload differs).
	otherVersionScope := PublishScope(validator.origin(), fixtureVendor, fixtureName, "0.1.1")

	// ARMS: the six TRAVERSAL-ERROR branches. The content-addressed walk from
	// payload.approvalRef to the canonical scope has six ways to fail that are
	// not policy decisions but broken links, and each returns its own message.
	// None of them is reachable through the landed minting path — that path
	// exists to produce well-formed objects — so each is built with
	// putRawApprovalObject over landed production encoders.
	absentRequest := hashref.SumSHA256([]byte("ac9-no-such-approval-request"))
	decisionNamingAbsentRequest := putRawApprovalObject(t, base, ApprovalDecisionV1,
		mustApprovalJSON(approvalDecisionWire{
			RequestRef: absentRequest.String(), Decision: "approve",
			DecidedBy: "attended-operator", Now: 11,
		}))
	decisionNamingADecision := putRawApprovalObject(t, base, ApprovalDecisionV1,
		mustApprovalJSON(approvalDecisionWire{
			RequestRef: good.identity.ApprovalRef.String(), Decision: "approve",
			DecidedBy: "attended-operator", Now: 11,
		}))
	decisionWithUnparseableRequestRef := putRawApprovalObject(t, base, ApprovalDecisionV1,
		mustApprovalJSON(approvalDecisionWire{
			RequestRef: "not-a-hash-ref", Decision: "approve",
			DecidedBy: "attended-operator", Now: 11,
		}))
	// The two undecodable objects below MUST NOT share payload bytes. The store
	// is content-addressed and PutObject is INSERT OR IGNORE, so identical
	// payloads are one object and the second semantic ID is silently dropped —
	// which is exactly what happened the first time this fixture was written.
	undecodableDecision := putRawApprovalObject(t, base, ApprovalDecisionV1,
		[]byte(`{"unknownDecisionField":1}`))

	// A decision that correctly heads and decides a request whose BYTES the
	// landed codec cannot read. decideApproval checks the semantic ID and the
	// head chain, never the request payload, so this is decidable and only
	// falls over at the broker's own decode.
	undecodableRequest := putRawApprovalObject(t, base, ApprovalRequestV1,
		[]byte(`{"unknownRequestField":1}`))
	if err := appendApprovalHead(base, undecodableRequest, hashref.HashRef{}); err != nil {
		t.Fatal(err)
	}
	decisionOverUndecodableRequest, err := decideApproval(
		base, undecodableRequest, "approve", "attended-operator", 11)
	if err != nil {
		t.Fatalf("decide over an undecodable request: %v", err)
	}

	// ARM: malformed. Two shapes — a ref that names an object of the WRONG
	// semantic kind, and a landed approval whose scope is not publish-grammar.
	malformedIdentity := good.identity
	malformedIdentity.ApprovalRef = good.approvalRequestRef // an ApprovalRequestV1, not a decision
	_, nonPublishDecision := landNonPublishApproval(t, base, "release", PublishCost,
		approvalTimes{request: 14, decide: 15, expires: 100})
	nonPublishIdentity := good.identity
	nonPublishIdentity.ApprovalRef = nonPublishDecision

	// ARM: missing. A syntactically valid ref that names no object at all.
	missing := newPublishFixture(t, validator.origin(), "ac9-missing")

	negatives := []struct {
		name     string
		payload  []byte
		scope    string
		now      int64
		sentinel error
		contains string
	}{
		{"missing", missing.payload, missing.scope, 50, ErrPublishApprovalMissing, "names no object"},
		{"expired", expired.payload, expired.scope, 50, ErrPublishApprovalExpired, "expired at logical time 40"},
		{"wrong-scope", EncodePublishPayload(wrongScopeIdentity, good.hashes), good.scope, 50,
			ErrPublishApprovalScope, "version:0.1.1"},
		{"wrong-hash", EncodePublishPayload(good.identity, flipped), good.scope, 50,
			ErrPublishApprovalScope, "tarball hash"},
		{"denied", denied.payload, denied.scope, 50, ErrPublishApprovalDenied, `says "deny"`},
		{"wrong-effect (landed, publish-grammatical, effect term is FS.Write)",
			wrongEffect.payload, wrongEffect.scope, 50,
			ErrPublishApprovalScope, `approval stamps effect "FS.Write"`},
		{"wrong-request-effect (the request was not minted by Human.Approve)",
			wrongRequestEffect.payload, wrongRequestEffect.scope, 50,
			ErrPublishApprovalScope, `was minted by effect "Human.PollApproval"`},
		{"wrong-cost (the attended request is priced at 2)",
			wrongCost.payload, wrongCost.scope, 50,
			ErrPublishApprovalScope, "approves cost 2, want exactly 1"},
		{"grant-scope-mismatch (stamp and payload agree, the capability does not)",
			good.payload, otherVersionScope, 50,
			ErrPublishApprovalScope, "does not describe the payload package"},
		{"request-after-publish (requested at 60, published at 50)",
			requestAfterPublish.payload, requestAfterPublish.scope, 50,
			ErrPublishApprovalMalformed, "at logical time 60, after the publish at 50"},
		{"decision-outside-request-publish-range (decided at 5, requested at 20)",
			decisionBeforeRequest.payload, decisionBeforeRequest.scope, 50,
			ErrPublishApprovalMalformed, "at logical time 5, outside [request 20, publish 50]"},
		{"expiry-precedes-its-own-request (expires 15, requested at 20)",
			expiryBeforeRequest.payload, expiryBeforeRequest.scope, 50,
			ErrPublishApprovalMalformed, "approval expiry 15 precedes its own request time 20"},
		{"traversal (the publish payload itself does not decode)",
			[]byte(`{"schema":"world/not-a-publish-request/v1"}`), good.scope, 50,
			ErrPublishApprovalMalformed, "publish payload schema"},
		{"traversal (decision object bytes do not decode)",
			good.withApprovalRef(undecodableDecision), good.scope, 50,
			ErrPublishApprovalMalformed, "decision: broker: decode approval payload"},
		{"traversal (decision requestRef is not a hash ref)",
			good.withApprovalRef(decisionWithUnparseableRequestRef), good.scope, 50,
			ErrPublishApprovalMalformed, "decision requestRef"},
		{"traversal (decision names a request that does not exist)",
			good.withApprovalRef(decisionNamingAbsentRequest), good.scope, 50,
			ErrApprovalRequestNotFound, "references request " + absentRequest.String()},
		{"traversal (decision names an object that is not a request)",
			good.withApprovalRef(decisionNamingADecision), good.scope, 50,
			ErrPublishApprovalMalformed, `want "world/approval-request/v1"`},
		{"traversal (request object bytes do not decode)",
			good.withApprovalRef(decisionOverUndecodableRequest), good.scope, 50,
			ErrPublishApprovalMalformed, "request: broker: decode approval payload"},
		{"malformed (decision ref names an approval REQUEST)",
			EncodePublishPayload(malformedIdentity, good.hashes), good.scope, 50,
			ErrPublishApprovalMalformed, "want \"world/approval-decision/v1\""},
		{"malformed (landed approval whose scope is not publish grammar)",
			EncodePublishPayload(nonPublishIdentity, good.hashes), good.scope, 50,
			ErrPublishApprovalMalformed, "carries no \"#\" mark"},
	}

	for _, tc := range negatives {
		t.Run(tc.name, func(t *testing.T) {
			before := readPublishCounters(validator, handler)
			session, recording := publishSession(
				t, base.Store, "ac9-"+tc.name, publishGrant(tc.scope), handler)
			_, ref, err := session.Invoke(context.Background(), EffectRequest{
				Effect: EffectRegistryPublish, Scope: tc.scope, Cost: PublishCost, Now: tc.now,
			}, tc.payload)
			if !errors.Is(err, tc.sentinel) {
				t.Fatalf("%s error = %T %v, want %v", tc.name, err, err, tc.sentinel)
			}
			if !strings.Contains(err.Error(), tc.contains) {
				t.Errorf("%s refusal %q does not mention %q", tc.name, err, tc.contains)
			}
			if !ref.IsZero() {
				t.Errorf("%s produced effect record %s; a pre-claim refusal records nothing", tc.name, ref)
			}
			if len(recording.effectIDs) != 0 {
				t.Errorf("%s minted durable effect intents %v, want none", tc.name, recording.effectIDs)
			}
			after := readPublishCounters(validator, handler)
			t.Logf("AC9 %s: counters %s -> %s", tc.name, before, after)
			if after != before {
				t.Fatalf("%s moved a counter (%s -> %s): the refusal must land BEFORE credential "+
					"load and BEFORE POST", tc.name, before, after)
			}
			if after != (publishCounters{}) {
				t.Fatalf("%s: counters are %s but nothing has legitimately dispatched yet", tc.name, after)
			}
		})
	}

	// POSITIVE CONTROL. Seven zero-valued counters prove nothing unless the
	// SAME counters, the SAME handler and the SAME validator can be shown to
	// move — under an approval that reached here through the landed
	// Human.Approve -> DecideApproval -> Human.PollApproval traversal that
	// landApproval performs and hash-checks.
	positiveSession, _ := publishSession(t, base.Store, "ac9-positive", publishGrant(good.scope), handler)
	if _, ref, err := positiveSession.Invoke(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish, Scope: good.scope, Cost: PublishCost, Now: 50,
	}, good.payload); err != nil || ref.IsZero() {
		t.Fatalf("AC9 positive control = ref %s, err %v; want a record and no error "+
			"(counters at refusal: %s)", ref, err, readPublishCounters(validator, handler))
	}
	positive := readPublishCounters(validator, handler)
	t.Logf("AC9 positive control: counters %s", positive)
	if positive != (publishCounters{posts: 1, dispatches: 1, credentialLoads: 1}) {
		t.Fatalf("AC9 positive control counters = %s, want POST=1 dispatches=1 credentialLoads=1", positive)
	}

	// ARM: already-consumed — the seventh negative class, and the only one that
	// can be stated against a counter that has just been shown to move.
	t.Run("already-consumed", func(t *testing.T) {
		before := readPublishCounters(validator, handler)
		session, recording := publishSession(
			t, base.Store, "ac9-consumed", publishGrant(good.scope), handler)
		if got := session.grants[0].Budget; got != PublishCost {
			t.Fatalf("the reuse session's budget = %d, want %d", got, PublishCost)
		}
		_, _, err := session.Invoke(context.Background(), EffectRequest{
			Effect: EffectRegistryPublish, Scope: good.scope, Cost: PublishCost, Now: 51,
		}, good.payload)
		if !errors.Is(err, store.ErrApprovalAlreadyConsumed) {
			t.Fatalf("already-consumed error = %T %v, want store.ErrApprovalAlreadyConsumed", err, err)
		}
		var denial *DenialError
		if errors.As(err, &denial) {
			t.Fatalf("the reuse was refused by BUDGET (%s), not by the durable claim", denial.Decision.Label)
		}
		if len(recording.effectIDs) != 0 {
			t.Errorf("already-consumed minted durable effect intents %v, want none", recording.effectIDs)
		}
		after := readPublishCounters(validator, handler)
		t.Logf("AC9 already-consumed: counters %s -> %s", before, after)
		if after != before {
			t.Fatalf("already-consumed moved a counter (%s -> %s)", before, after)
		}
	})
}

// AC9a — close and REOPEN the store, then a new session with a new budget. The
// SHARED fake's total POST count must remain EXACTLY ONE.
func TestConsumedApprovalStaysConsumedAcrossStoreCloseAndReopen(t *testing.T) {
	validator := newFakeValidator(t, "ok")
	base := openPublishStore(t)
	fixture := newPublishFixture(t, validator.origin(), "ac9a").
		landApproval(t, base, "approve", defaultApprovalTimes())
	// ONE handler and ONE validator span both processes, so the POST count
	// below is a total and not a per-session reading.
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

	first, _ := publishSession(t, base.Store, "ac9a-first", publishGrant(fixture.scope), handler)
	if _, _, err := first.Invoke(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish, Scope: fixture.scope, Cost: PublishCost, Now: 50,
	}, fixture.payload); err != nil {
		t.Fatalf("AC9a first invocation: %v", err)
	}
	afterFirst := readPublishCounters(validator, handler)
	t.Logf("AC9a arm 1 (first invocation, original store): %s", afterFirst)
	if afterFirst != (publishCounters{posts: 1, dispatches: 1, credentialLoads: 1}) {
		t.Fatalf("AC9a: counters after the first invocation = %s, want POST=1 dispatches=1 credentialLoads=1",
			afterFirst)
	}

	// CLOSE and REOPEN. The reopened handle shares no connection, no cache and
	// no budget with the first — only the bytes on disk.
	reopened := base.reopen(t)
	afterReopen := readPublishCounters(validator, handler)
	t.Logf("AC9a arm 2 (after close+reopen, before retry): %s", afterReopen)

	second, recording := publishSession(t, reopened.Store, "ac9a-second", publishGrant(fixture.scope), handler)
	if got := second.grants[0].Budget; got != PublishCost {
		t.Fatalf("AC9a: the reopened session's budget = %d, want a FRESH %d", got, PublishCost)
	}
	_, _, err := second.Invoke(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish, Scope: fixture.scope, Cost: PublishCost, Now: 60,
	}, fixture.payload)
	if !errors.Is(err, store.ErrApprovalAlreadyConsumed) {
		t.Fatalf("AC9a: reuse after reopen = %T %v, want store.ErrApprovalAlreadyConsumed", err, err)
	}
	var denial *DenialError
	if errors.As(err, &denial) {
		t.Fatalf("AC9a: the reuse was refused by BUDGET (%s), not by the durable claim", denial.Decision.Label)
	}
	if len(recording.effectIDs) != 0 {
		t.Errorf("AC9a: the reopened session minted durable intents %v, want none", recording.effectIDs)
	}
	afterSecond := readPublishCounters(validator, handler)
	t.Logf("AC9a arm 3 (after the refused retry): %s", afterSecond)
	if got := int64(validator.count()); got != 1 {
		t.Fatalf("AC9a: SHARED fake's total POST count = %d, want EXACTLY 1", got)
	}
	if afterSecond != afterFirst {
		t.Fatalf("AC9a: counters moved across close+reopen+retry: %s -> %s", afterFirst, afterSecond)
	}
}

// AC9b — two concurrent sessions with distinct FRESH budgets race the same
// approval behind a START BARRIER. Run this under -race.
func TestTwoSessionsRacingOneApprovalDispatchExactlyOnce(t *testing.T) {
	validator := newFakeValidator(t, "ok")
	base := openPublishStore(t)
	fixture := newPublishFixture(t, validator.origin(), "ac9b").
		landApproval(t, base, "approve", defaultApprovalTimes())
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

	const racers = 2
	type outcome struct {
		err       error
		effectIDs []string
	}
	outcomes := make([]outcome, racers)
	// THE START BARRIER. Without it this is a sequential test wearing a race
	// test's name: both goroutines must be parked at the same instruction
	// before either is allowed to touch the store.
	var ready, done sync.WaitGroup
	start := make(chan struct{})
	ready.Add(racers)
	done.Add(racers)
	for i := 0; i < racers; i++ {
		go func(i int) {
			defer done.Done()
			// Each racer gets its OWN session over the SAME store, with its own
			// FRESH budget of one.
			recording := &publishRecordingStore{base: base.Store}
			session := newSession(recording, fmt.Sprintf("ac9b-%d", i),
				[]Capability{publishGrant(fixture.scope)},
				Registry{EffectRegistryPublish: handler}, Live, nil)
			ready.Done()
			<-start
			_, _, err := session.Invoke(context.Background(), EffectRequest{
				Effect: EffectRegistryPublish, Scope: fixture.scope, Cost: PublishCost, Now: 50,
			}, fixture.payload)
			outcomes[i] = outcome{err: err, effectIDs: recording.effectIDs}
		}(i)
	}
	ready.Wait()
	close(start)
	done.Wait()

	winners, losers := 0, 0
	for i, got := range outcomes {
		t.Logf("AC9b racer %d: err=%v durableIntents=%v", i, got.err, got.effectIDs)
		switch {
		case got.err == nil:
			winners++
			if len(got.effectIDs) != 1 {
				t.Errorf("AC9b racer %d won but minted durable intents %v, want exactly 1",
					i, got.effectIDs)
			}
		case errors.Is(got.err, store.ErrApprovalAlreadyConsumed):
			losers++
			if len(got.effectIDs) != 0 {
				t.Errorf("AC9b racer %d lost but minted durable intents %v, want none",
					i, got.effectIDs)
			}
			var denial *DenialError
			if errors.As(got.err, &denial) {
				t.Errorf("AC9b racer %d was refused by BUDGET (%s), not by the durable claim",
					i, denial.Decision.Label)
			}
		default:
			t.Errorf("AC9b racer %d: unexpected error %T %v", i, got.err, got.err)
		}
	}
	if winners != 1 || losers != racers-1 {
		t.Fatalf("AC9b: %d winners and %d losers, want exactly 1 and %d", winners, losers, racers-1)
	}
	counters := readPublishCounters(validator, handler)
	t.Logf("AC9b arm (both racers complete): %s", counters)
	if counters != (publishCounters{posts: 1, dispatches: 1, credentialLoads: 1}) {
		t.Fatalf("AC9b: counters = %s, want POST=1 dispatches=1 credentialLoads=1 — the loser must "+
			"be refused BEFORE credential load", counters)
	}
}

// AC9c — a dispatched-but-ambiguous attempt BURNS the approval. After a store
// reopen a fresh session/budget retrying the same stamp receives
// ErrApprovalAlreadyConsumed (not budget denial), the total POST count stays
// one, and recovery reports the pending publish without launching anything.
func TestIndeterminatePublishBurnsTheApprovalAndRecoveryStaysReadOnly(t *testing.T) {
	// "reset": the validator ACCEPTS the request (so the body demonstrably left
	// the publisher) and then aborts the connection. That is the genuinely
	// ambiguous case, and the reason burning the stamp is the safe choice.
	validator := newFakeValidator(t, "reset")
	base := openPublishStore(t)
	fixture := newPublishFixture(t, validator.origin(), "ac9c").
		landApproval(t, base, "approve", defaultApprovalTimes())
	handler := loopbackHandler(t, RegistryPublishConfig{
		PublisherPath:   writePublisherScript(t, "reset", nil),
		PackageDir:      fixture.dir,
		Manifest:        fixture.manifest,
		RegistryOrigin:  validator.origin(),
		ValidatorOrigin: validator.origin(),
		Credential:      writeCredentialFile(t, ac10Sentinel),
		Approval:        fixture.approval,
		ExecTimeout:     20 * time.Second,
	})

	first, recording := publishSession(t, base.Store, "ac9c-first", publishGrant(fixture.scope), handler)
	_, ref, err := first.Invoke(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish, Scope: fixture.scope, Cost: PublishCost, Now: 50,
	}, fixture.payload)
	var indeterminate *IndeterminateEffectError
	if !errors.As(err, &indeterminate) {
		t.Fatalf("AC9c first invocation error = %T %v, want *IndeterminateEffectError", err, err)
	}
	if !ref.IsZero() {
		t.Errorf("AC9c: the ambiguous arm produced effect record %s, want none", ref)
	}
	afterFirst := readPublishCounters(validator, handler)
	t.Logf("AC9c arm 1 (first invocation, typed indeterminate): %s", afterFirst)
	if afterFirst != (publishCounters{posts: 1, dispatches: 1, credentialLoads: 1}) {
		t.Fatalf("AC9c: counters after the ambiguous dispatch = %s, want POST=1 dispatches=1 "+
			"credentialLoads=1", afterFirst)
	}
	if got := effectReceipt(t, recording, 0).State; got != store.ReceiptIndeterminate {
		t.Fatalf("AC9c: receipt state = %q, want %q", got, store.ReceiptIndeterminate)
	}

	reopened := base.reopen(t)

	// RECOVERY IS READ-ONLY. It is handed the real publish registry, and the
	// counters must not move: a surface that can report an unresolved
	// irreversible attempt must not be able to launch a second one.
	findings, err := Recover(reopened.Store, Registry{EffectRegistryPublish: handler})
	if err != nil {
		t.Fatal(err)
	}
	publishes := PendingPublishes(findings)
	if len(publishes) != 1 {
		t.Fatalf("AC9c: PendingPublishes = %d findings, want exactly 1 (all findings: %d)",
			len(publishes), len(findings))
	}
	if got := publishes[0].EffectIntent.Effect; got != EffectRegistryPublish {
		t.Errorf("AC9c: pending publish effect = %q, want %q", got, EffectRegistryPublish)
	}
	if got := publishes[0].EffectIntent.Scope; got != fixture.scope {
		t.Errorf("AC9c: pending publish scope = %q, want %q", got, fixture.scope)
	}
	// The selector must actually SELECT: a filter that returns everything would
	// pass the assertion above whenever the store holds exactly one finding.
	if got := PendingPublishes(nil); got != nil {
		t.Errorf("AC9c: PendingPublishes(nil) = %v, want nil", got)
	}
	if got := PendingPublishes([]IndeterminateEffect{{
		EffectIntent: store.EffectIntent{Effect: "FS.Write"},
	}}); len(got) != 0 {
		t.Errorf("AC9c: PendingPublishes selected a non-publish finding %v", got)
	}
	afterRecovery := readPublishCounters(validator, handler)
	t.Logf("AC9c arm 2 (after read-only recovery over the reopened store): %s", afterRecovery)
	if afterRecovery != afterFirst {
		t.Fatalf("AC9c: recovery moved a counter (%s -> %s); recovery must never dispatch",
			afterFirst, afterRecovery)
	}

	// The retry, with a FRESH session and a FRESH budget.
	second, secondRecording := publishSession(
		t, reopened.Store, "ac9c-second", publishGrant(fixture.scope), handler)
	if got := second.grants[0].Budget; got != PublishCost {
		t.Fatalf("AC9c: the retry session's budget = %d, want a FRESH %d", got, PublishCost)
	}
	_, _, retryErr := second.Invoke(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish, Scope: fixture.scope, Cost: PublishCost, Now: 60,
	}, fixture.payload)
	if !errors.Is(retryErr, store.ErrApprovalAlreadyConsumed) {
		t.Fatalf("AC9c: retry error = %T %v, want store.ErrApprovalAlreadyConsumed", retryErr, retryErr)
	}
	// THE DISTINCTION AC9c EXISTS FOR: budget denial would also leave the POST
	// counter reading 1 while proving nothing about the claim.
	var denial *DenialError
	if errors.As(retryErr, &denial) {
		t.Fatalf("AC9c: the retry was refused by BUDGET (%s), not by the durable claim",
			denial.Decision.Label)
	}
	if len(secondRecording.effectIDs) != 0 {
		t.Errorf("AC9c: the retry minted durable intents %v, want none", secondRecording.effectIDs)
	}
	afterRetry := readPublishCounters(validator, handler)
	t.Logf("AC9c arm 3 (after the refused fresh-budget retry): %s", afterRetry)
	if got := int64(validator.count()); got != 1 {
		t.Fatalf("AC9c: total POST count = %d, want EXACTLY 1", got)
	}
	if afterRetry != afterFirst {
		t.Fatalf("AC9c: counters moved across reopen+recovery+retry: %s -> %s", afterFirst, afterRetry)
	}
}
