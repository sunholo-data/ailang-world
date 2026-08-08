package broker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// ---------------------------------------------------------------------------
// SM.C — probe-then-resolve reconciliation (AC13, AC14, AC15, AC16, AC16a,
// AC16b, AC17)
//
// SAFETY, restated for this file because it is the file that adds an
// IN-PROCESS HTTP client to host/broker:
//
//  1. every reconciliation here is driven through reconcileLoopback, which is
//     UNEXPORTED and, via validatePublishOrigin(allowLoopback=true), REFUSES a
//     non-loopback origin. TestReconcileRefusalSetWithAPassingPositiveControl
//     drives that refusal.
//  2. every origin any test hands it is an httptest server, which binds
//     127.0.0.1 only. TestEveryReconcileRequestWasLoopbackAndAGet enumerates
//     every request every fake bucket in this file ever saw.
//  3. the two arms that DO call the production entry point
//     ReconcileRegistryPublish refuse at R1/R2, which are evaluated before the
//     first fetch, and the assertion is on the refusal — no sample is taken.
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// MEASURED fixture documents
//
// Every byte pattern below is a MEASUREMENT, not an invention. The controller
// re-measured all three arms first-party on 2026-08-08; the design doc's row
// V-N is the original measurement.
// ---------------------------------------------------------------------------

// gcsNoSuchKey reproduces the bucket's 404 document:
//
//	GET https://storage.googleapis.com/ailang-registry/packages/world/core/0.1.0/metadata.json
//	-> 404, 217 bytes, `<?xml ...?><Error><Code>NoSuchKey</Code>...</Error>`
//
// It is reproduced to SHAPE (the element structure the decoder reads), not to
// the byte, because the <Details> element embeds the requested key and the
// literal length is therefore a function of the key, not of the server.
const gcsNoSuchKey = `<?xml version='1.0' encoding='UTF-8'?><Error>` +
	`<Code>NoSuchKey</Code><Message>The specified key does not exist.</Message>` +
	`<Details>No such object: ailang-registry/packages/world/core/0.1.0/metadata.json</Details>` +
	`</Error>`

// validatorPlainText404 is the MEASURED validator-origin 404 body:
//
//	GET https://registry.ailang.sunholo.com/packages/sunholo/auth/0.4.1/metadata.json
//	-> 404, 19 bytes, PLAIN TEXT `404 page not found` — NOT XML.
//
// It matters because it is the body a failover to the validator would produce,
// and it must never be mistaken for the bucket's absence document.
const validatorPlainText404 = "404 page not found\n"

// reconcileExpected are the three digests the durable publish request bound.
// They are PAIRWISE DISTINCT on purpose: "binds all three hashes" (AC14) is
// unfalsifiable if two of the three carry the same value, because a comparator
// that reads the wrong field would still agree.
var reconcileExpected = PublishHashes{
	TarballSHA256: "sha256:" + strings.Repeat("11", 32),
	ContentHash:   "sha256:" + strings.Repeat("22", 32),
	InterfaceHash: "sha256:" + strings.Repeat("33", 32),
}

// ---------------------------------------------------------------------------
// the loopback fake bucket
// ---------------------------------------------------------------------------

type bucketRequest struct{ Method, Path string }

type bucketResponse struct {
	status int
	body   []byte
	// hijack aborts the connection instead of answering. It is the only way to
	// make a TARGET fetch FAIL while the CONTROL on the same origin succeeds
	// (classifier branch C4); no status code can express that, and the two
	// URLs are structurally forced onto one origin.
	hijack bool
}

// fakeBucket is a loopback stand-in for the READ-ONLY public bucket. It serves
// exact object keys, records every request it ever sees (method AND path), and
// falls back to the measured GCS NoSuchKey document for an unknown key —
// exactly as the real bucket does.
type fakeBucket struct {
	server   *httptest.Server
	mu       sync.Mutex
	objects  map[string]bucketResponse
	log      []bucketRequest
	notFound bucketResponse
	onGet    func(path string)
	closed   bool
	url      string
}

// bucketAudit is the process-wide census every fake bucket in this file
// registers with. TestEveryReconcileRequestWasLoopbackAndAGet reads it, so the
// "read-only, loopback-only" claim is enumerated over the requests that
// actually happened rather than asserted per test.
var bucketAudit struct {
	mu       sync.Mutex
	origins  []string
	requests []bucketRequest
}

func newFakeBucket(t *testing.T) *fakeBucket {
	t.Helper()
	b := &fakeBucket{
		objects:  map[string]bucketResponse{},
		notFound: bucketResponse{status: http.StatusNotFound, body: []byte(gcsNoSuchKey)},
	}
	b.server = httptest.NewServer(http.HandlerFunc(b.serve))
	b.url = b.server.URL
	bucketAudit.mu.Lock()
	bucketAudit.origins = append(bucketAudit.origins, b.url)
	bucketAudit.mu.Unlock()
	t.Cleanup(b.close)
	return b
}

func (b *fakeBucket) serve(w http.ResponseWriter, r *http.Request) {
	b.mu.Lock()
	entry := bucketRequest{Method: r.Method, Path: r.URL.Path}
	b.log = append(b.log, entry)
	resp, ok := b.objects[r.URL.Path]
	if !ok {
		resp = b.notFound
	}
	hook := b.onGet
	b.mu.Unlock()

	bucketAudit.mu.Lock()
	bucketAudit.requests = append(bucketAudit.requests, entry)
	bucketAudit.mu.Unlock()

	_, _ = io.Copy(io.Discard, r.Body)
	if hook != nil {
		hook(r.URL.Path)
	}
	if resp.hijack {
		hijacker, ok := w.(http.Hijacker)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		conn, _, err := hijacker.Hijack()
		if err != nil {
			return
		}
		_ = conn.Close()
		return
	}
	w.WriteHeader(resp.status)
	_, _ = w.Write(resp.body)
}

func (b *fakeBucket) origin() string { return b.url }

func (b *fakeBucket) put(path string, resp bucketResponse) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.objects[path] = resp
}

func (b *fakeBucket) requests() []bucketRequest {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]bucketRequest(nil), b.log...)
}

func (b *fakeBucket) count() int { return len(b.requests()) }

func (b *fakeBucket) countPath(path string) int {
	n := 0
	for _, req := range b.requests() {
		if req.Path == path {
			n++
		}
	}
	return n
}

// close stops the server. It is idempotent so the "unreachable host" fixture
// can close it early and t.Cleanup can still run.
func (b *fakeBucket) close() {
	b.mu.Lock()
	already := b.closed
	b.closed = true
	b.mu.Unlock()
	if !already {
		b.server.Close()
	}
}

// bucketKey builds the object key INDEPENDENTLY of metadataObjectURL. Every
// fixture in this file seeds the bucket through it, and
// TestAbsenceControlSharesTargetKeySpace compares the two constructions. A
// single source would make that comparison a tautology and would let a change
// to the production key shape re-green the fixtures with it.
func bucketKey(vendor, name, version string) string {
	return "/packages/" + vendor + "/" + name + "/" + version + "/metadata.json"
}

func targetKey() string { return bucketKey(fixtureVendor, fixtureName, fixtureVersion) }

func controlKey() string {
	return bucketKey(DefaultProbeControlVendor, DefaultProbeControlName, DefaultProbeControlVersion)
}

// metadataDocument writes the served document's JSON key names out LITERALLY
// rather than marshalling a RegistryMetadata. That makes the fixture an
// independent second source of the wire names measured off the live object
// (`name`, `version`, `tarball_hash`, `content_hash`, `interface_hash`): if a
// struct tag is renamed, the reconciler stops recognising this fixture instead
// of the pair drifting together in silence.
func metadataDocument(name, version string, h PublishHashes) []byte {
	body, err := json.Marshal(map[string]any{
		"name":           name,
		"version":        version,
		"tarball_hash":   h.TarballSHA256,
		"content_hash":   h.ContentHash,
		"interface_hash": h.InterfaceHash,
	})
	if err != nil {
		panic("metadataDocument: " + err.Error())
	}
	return body
}

// firingControl seeds the same-pass known-positive control object: a 200 with a
// well-formed metadata document for the MEASURED control package
// sunholo/auth@0.4.1, under the target's own origin and key-space.
func firingControl(b *fakeBucket) {
	b.put(controlKey(), bucketResponse{
		status: http.StatusOK,
		body: metadataDocument("sunholo/auth", "0.4.1", PublishHashes{
			TarballSHA256: "sha256:" + strings.Repeat("aa", 32),
			ContentHash:   "sha256:" + strings.Repeat("bb", 32),
			InterfaceHash: "sha256:" + strings.Repeat("cc", 32),
		}),
	})
}

// reconcileCfg is the base configuration every arm starts from. The control
// vendor/name/version are LEFT EMPTY so withDefaults supplies the production
// defaults — the fixtures therefore exercise the shipped control, not a
// test-chosen one.
func reconcileCfg(b *fakeBucket) ReconcileConfig {
	return ReconcileConfig{
		RegistryOrigin:        b.origin(),
		Vendor:                fixtureVendor,
		Name:                  fixtureName,
		Version:               fixtureVersion,
		Expected:              reconcileExpected,
		AbsentSamplesRequired: 1,
		MaxAttempts:           2,
		RequestTimeout:        10 * time.Second,
	}
}

func mustReconcile(t *testing.T, cfg ReconcileConfig) ReconcileReceipt {
	t.Helper()
	receipt, err := reconcileLoopback(context.Background(), cfg)
	if err != nil {
		t.Fatalf("reconcileLoopback: %v", err)
	}
	return receipt
}

// ---------------------------------------------------------------------------
// AC13 — recovery REPORTS an indeterminate publish and DISPATCHES NOTHING
// ---------------------------------------------------------------------------

// countingNetworkHandler is the AC13/AC17 instrument. It does exactly two
// observable things: it increments a counter, and it performs a metadata GET
// against a loopback bucket. Two INDEPENDENT witnesses of the same event, so a
// zero reading cannot be an artefact of one of them being unwired — and each
// test that reads zero from it ends by driving it once on purpose and watching
// BOTH witnesses move.
//
// It is deliberately NOT a *RegistryPublishHandler. A real publish handler
// refuses a malformed request long before its own dispatch counter moves, so a
// mutation that made recovery call it would leave every counter at zero and
// the guard would report success. This handler counts the CALL, which is the
// event AC13 is about.
type countingNetworkHandler struct {
	calls  atomic.Int64
	bucket *fakeBucket
	url    string
}

func newCountingNetworkHandler(t *testing.T) *countingNetworkHandler {
	t.Helper()
	bucket := newFakeBucket(t)
	firingControl(bucket)
	return &countingNetworkHandler{bucket: bucket, url: bucket.origin() + controlKey()}
}

func (h *countingNetworkHandler) Execute(
	ctx context.Context, _ EffectRequest, _ []byte,
) ([]byte, error) {
	h.calls.Add(1)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, resp.Body)
	return []byte(`{"schema":"counting-network-handler"}`), nil
}

func (h *countingNetworkHandler) String() string {
	return fmt.Sprintf("handlerCalls=%d bucketRequests=%d", h.calls.Load(), h.bucket.count())
}

// provePositiveControl drives the instrument once and requires BOTH witnesses
// to move. Without it every zero this handler reports is compatible with "the
// handler was never wired to anything".
func (h *countingNetworkHandler) provePositiveControl(t *testing.T) {
	t.Helper()
	beforeCalls, beforeRequests := h.calls.Load(), h.bucket.count()
	if _, err := h.Execute(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish,
	}, nil); err != nil {
		t.Fatalf("positive control: the instrument could not be driven at all: %v", err)
	}
	if got := h.calls.Load(); got != beforeCalls+1 {
		t.Fatalf("positive control: handler call counter %d -> %d, want +1", beforeCalls, got)
	}
	if got := h.bucket.count(); got != beforeRequests+1 {
		t.Fatalf("positive control: bucket request counter %d -> %d, want +1", beforeRequests, got)
	}
	t.Logf("positive control fired: %s", h)
}

func TestRecoveryReportsTheIndeterminatePublishWithoutDispatchingAnyHandler(t *testing.T) {
	// "reset": the validator ACCEPTS the request and then aborts the
	// connection, which is the genuinely ambiguous outcome SM.C exists to
	// resolve.
	validator := newFakeValidator(t, "reset")
	base := openPublishStore(t)
	fixture := newPublishFixture(t, validator.origin(), "ac13").
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

	session, _ := publishSession(t, base.Store, "ac13-first", publishGrant(fixture.scope), handler)
	_, _, err := session.Invoke(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish, Scope: fixture.scope, Cost: PublishCost, Now: 50,
	}, fixture.payload)
	var indeterminate *IndeterminateEffectError
	if !errors.As(err, &indeterminate) {
		t.Fatalf("AC13 setup: first invocation error = %T %v, want *IndeterminateEffectError", err, err)
	}
	postsBeforeRecovery := validator.count()

	reopened := base.reopen(t)

	// THE INSTRUMENT. It is handed to Recover through the production
	// `registries ...Registry` parameter, which exists for exactly this reason:
	// it makes the no-dispatch policy observable at the production boundary.
	instrument := newCountingNetworkHandler(t)
	findings, err := Recover(reopened.Store, Registry{
		EffectRegistryPublish: instrument,
		"FS.Write":            instrument,
	})
	if err != nil {
		t.Fatal(err)
	}

	// (a) recovery REPORTS the indeterminate publish.
	publishes := PendingPublishes(findings)
	if len(publishes) != 1 {
		t.Fatalf("AC13: PendingPublishes = %d, want exactly 1 (all findings %d)",
			len(publishes), len(findings))
	}
	reported := publishes[0]
	if reported.EffectIntent.Effect != EffectRegistryPublish {
		t.Errorf("AC13: reported effect = %q, want %q",
			reported.EffectIntent.Effect, EffectRegistryPublish)
	}
	if reported.EffectIntent.Scope != fixture.scope {
		t.Errorf("AC13: reported scope = %q, want %q", reported.EffectIntent.Scope, fixture.scope)
	}
	if reported.Err == nil || reported.Err.InvocationID == "" {
		t.Fatalf("AC13: reported finding carries no identifying error: %+v", reported)
	}
	if !strings.Contains(reported.Err.Error(), "indeterminate") {
		t.Errorf("AC13: reported error %q does not say the publish is indeterminate",
			reported.Err.Error())
	}
	// A report that names nothing is not a report: the finding must identify
	// WHICH irreversible attempt is unresolved.
	if reported.EffectIntent.RequestRef.IsZero() {
		t.Error("AC13: the reported finding names no durable request object")
	}

	// (b) and it dispatched NOTHING.
	if got := instrument.calls.Load(); got != 0 {
		t.Fatalf("AC13: recovery called a handler %d time(s), want 0 (%s)", got, instrument)
	}
	if got := instrument.bucket.count(); got != 0 {
		t.Fatalf("AC13: recovery caused %d network request(s), want 0 (%s)", got, instrument)
	}
	if got := handler.Dispatches(); got != 1 {
		t.Errorf("AC13: publisher dispatches = %d, want the 1 from the ambiguous attempt only", got)
	}
	if got := handler.CredentialLoads(); got != 1 {
		t.Errorf("AC13: credential loads = %d, want the 1 from the ambiguous attempt only", got)
	}
	if got := validator.count(); got != postsBeforeRecovery {
		t.Fatalf("AC13: recovery moved the POST counter %d -> %d", postsBeforeRecovery, got)
	}

	// (c) THE SAME-RUN POSITIVE CONTROL for the zero above.
	instrument.provePositiveControl(t)
}

// ---------------------------------------------------------------------------
// AC14 — exact public metadata resolves the indeterminate receipt to
//        `succeeded-reconciled`, binding ALL THREE hashes
// ---------------------------------------------------------------------------

func TestExactPublicMetadataResolvesToSucceededReconciled(t *testing.T) {
	bucket := newFakeBucket(t)
	firingControl(bucket)

	// MUT-SM-RECON-HASH anchor: the ONE expression that supplies the served
	// interface hash. Flipping it here must turn this success into a conflict.
	servedInterfaceHash := reconcileExpected.InterfaceHash // MUT-SM-RECON-HASH anchor

	bucket.put(targetKey(), bucketResponse{
		status: http.StatusOK,
		body: metadataDocument(fixturePackage, fixtureVersion, PublishHashes{
			TarballSHA256: reconcileExpected.TarballSHA256,
			ContentHash:   reconcileExpected.ContentHash,
			InterfaceHash: servedInterfaceHash,
		}),
	})

	receipt := mustReconcile(t, reconcileCfg(bucket))
	t.Logf("AC14 receipt: %s", receipt)
	if receipt.State != ReconcileSucceededReconciled {
		t.Fatalf("AC14: receipt state = %q, want %q (detail %q)",
			receipt.State, ReconcileSucceededReconciled, receipt.Detail)
	}
	if receipt.Served == nil {
		t.Fatal("AC14: a succeeded-reconciled receipt carries no served document")
	}

	// ALL THREE, named individually, against digests that are pairwise
	// distinct — so a comparator reading the wrong field cannot agree by luck.
	for _, arm := range []struct{ name, got, want string }{
		{"tarball", receipt.Served.TarballHash, reconcileExpected.TarballSHA256},
		{"content", receipt.Served.ContentHash, reconcileExpected.ContentHash},
		{"interface", receipt.Served.InterfaceHash, reconcileExpected.InterfaceHash},
	} {
		if arm.got != arm.want {
			t.Errorf("AC14: bound %s hash = %q, want %q", arm.name, arm.got, arm.want)
		}
	}
	if receipt.Served.Name != fixturePackage || receipt.Served.Version != fixtureVersion {
		t.Errorf("AC14: bound identity = %s@%s, want %s@%s",
			receipt.Served.Name, receipt.Served.Version, fixturePackage, fixtureVersion)
	}
	// Non-vacuity of "all three": the three expected digests must differ from
	// each other, or the loop above proves nothing.
	if reconcileExpected.TarballSHA256 == reconcileExpected.ContentHash ||
		reconcileExpected.ContentHash == reconcileExpected.InterfaceHash ||
		reconcileExpected.TarballSHA256 == reconcileExpected.InterfaceHash {
		t.Fatal("AC14: the three expected digests are not pairwise distinct; " +
			"the three-hash binding assertion above would be satisfiable by reading one field")
	}
	if got := len(receipt.Samples); got != 1 {
		t.Errorf("AC14: samples = %d, want 1 (present resolves on the first pass)", got)
	}
	if got := receipt.Samples[0].Verdict; got != ProbePresent {
		t.Errorf("AC14: sample verdict = %q, want %q", got, ProbePresent)
	}
	// The reconciler READ. It never wrote.
	for _, req := range bucket.requests() {
		if req.Method != http.MethodGet {
			t.Errorf("AC14: reconciliation issued a %s to %s; its only verb is GET",
				req.Method, req.Path)
		}
	}
}

// ---------------------------------------------------------------------------
// AC15 — an observed 409 resolves BY METADATA, never automatically
// ---------------------------------------------------------------------------

func TestObserved409ResolvesByMetadataNeverAutomatically(t *testing.T) {
	// Both arms observe the SAME 409. The only thing that differs is what the
	// bucket serves, which is the whole point: a 409 is evidence about nothing
	// until the metadata has been read.
	const observed409 = http.StatusConflict

	t.Run("mismatching_metadata_resolves_conflict", func(t *testing.T) {
		bucket := newFakeBucket(t)
		firingControl(bucket)
		bucket.put(targetKey(), bucketResponse{
			status: http.StatusOK,
			body: metadataDocument(fixturePackage, fixtureVersion, PublishHashes{
				TarballSHA256: reconcileExpected.TarballSHA256,
				ContentHash:   reconcileExpected.ContentHash,
				InterfaceHash: flipLastNibble(reconcileExpected.InterfaceHash),
			}),
		})
		cfg := reconcileCfg(bucket)
		cfg.ObservedPublishStatus = observed409

		receipt := mustReconcile(t, cfg)
		t.Logf("AC15 mismatching arm: %s", receipt)
		if receipt.State == ReconcileSucceededReconciled {
			t.Fatalf("AC15: a 409 with MISMATCHING metadata resolved to SUCCESS (detail %q); "+
				"someone else may hold this immutable version with different bytes", receipt.Detail)
		}
		if receipt.State != ReconcileConflict {
			t.Fatalf("AC15: receipt state = %q, want %q (detail %q)",
				receipt.State, ReconcileConflict, receipt.Detail)
		}
		if !strings.Contains(receipt.Detail, "interface") {
			t.Errorf("AC15: conflict detail %q does not name WHICH digest diverged", receipt.Detail)
		}
		if receipt.ObservedPublishStatus != observed409 {
			t.Errorf("AC15: receipt observedPublishStatus = %d, want the recorded %d",
				receipt.ObservedPublishStatus, observed409)
		}
	})

	t.Run("matching_metadata_negative_control", func(t *testing.T) {
		bucket := newFakeBucket(t)
		firingControl(bucket)
		bucket.put(targetKey(), bucketResponse{
			status: http.StatusOK,
			body:   metadataDocument(fixturePackage, fixtureVersion, reconcileExpected),
		})
		cfg := reconcileCfg(bucket)
		cfg.ObservedPublishStatus = observed409

		receipt := mustReconcile(t, cfg)
		t.Logf("AC15 matching arm (negative control): %s", receipt)
		if receipt.State != ReconcileSucceededReconciled {
			t.Fatalf("AC15 negative control: receipt state = %q, want %q (detail %q). "+
				"Without this arm, 'never success' would be satisfiable by a reconciler "+
				"that can only ever say conflict",
				receipt.State, ReconcileSucceededReconciled, receipt.Detail)
		}
	})
}

// ---------------------------------------------------------------------------
// AC16 — bounded repeated absence resolves `not-published`, without a POST
// ---------------------------------------------------------------------------

func TestBoundedAbsenceResolvesNotPublishedWithoutAPost(t *testing.T) {
	bucket := newFakeBucket(t)
	firingControl(bucket)
	// The target key is deliberately NOT seeded: the bucket answers it with the
	// measured NoSuchKey document, exactly as the real one does.

	cfg := reconcileCfg(bucket)
	cfg.AbsentSamplesRequired = 3
	cfg.MaxAttempts = 8

	receipt := mustReconcile(t, cfg)
	t.Logf("AC16 receipt: %s", receipt)
	if receipt.State != ReconcileNotPublished {
		t.Fatalf("AC16: receipt state = %q, want %q (detail %q)",
			receipt.State, ReconcileNotPublished, receipt.Detail)
	}

	// THE BOUND. Stated as an exact sample count rather than as wall clock:
	// the count is exact and noise-free, where a duration threshold on this rig
	// has noise the size of the signal it would be guarding.
	if got := receipt.AbsentSamples; got != cfg.AbsentSamplesRequired {
		t.Errorf("AC16: absent samples = %d, want exactly the required %d",
			got, cfg.AbsentSamplesRequired)
	}
	if got := len(receipt.Samples); got != cfg.AbsentSamplesRequired {
		t.Fatalf("AC16: the pass took %d samples for a %d-sample window; it did not stop "+
			"when the window closed (MaxAttempts was %d)",
			got, cfg.AbsentSamplesRequired, cfg.MaxAttempts)
	}
	if got := receipt.UninformativeSamples; got != 0 {
		t.Errorf("AC16: uninformative samples = %d, want 0 with a firing control", got)
	}
	if got := bucket.countPath(targetKey()); got != cfg.AbsentSamplesRequired {
		t.Errorf("AC16: target fetches = %d, want exactly %d", got, cfg.AbsentSamplesRequired)
	}
	if got := bucket.countPath(controlKey()); got != cfg.AbsentSamplesRequired {
		t.Errorf("AC16: control fetches = %d, want one per sample (%d)",
			got, cfg.AbsentSamplesRequired)
	}

	// WITHOUT A POST. The bucket is the only origin the reconciler was given,
	// and every request it saw must be a GET.
	for _, req := range bucket.requests() {
		if req.Method != http.MethodGet {
			t.Fatalf("AC16: reconciliation issued %s %s; resolving absence must not write",
				req.Method, req.Path)
		}
	}
	if got := bucket.count(); got != 2*cfg.AbsentSamplesRequired {
		t.Errorf("AC16: total bucket requests = %d, want %d (one control + one target per sample)",
			got, 2*cfg.AbsentSamplesRequired)
	}
}

// TestNotPublishedDoesNotReAuthorizeTheConsumedApproval is AC16's second
// clause: resolving `not-published` says the artifact is absent; it does NOT
// hand back the attended authority the ambiguous attempt already burned.
func TestNotPublishedDoesNotReAuthorizeTheConsumedApproval(t *testing.T) {
	validator := newFakeValidator(t, "reset")
	base := openPublishStore(t)
	fixture := newPublishFixture(t, validator.origin(), "ac16-first").
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
	first, _ := publishSession(t, base.Store, "ac16-first", publishGrant(fixture.scope), handler)
	_, _, err := first.Invoke(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish, Scope: fixture.scope, Cost: PublishCost, Now: 50,
	}, fixture.payload)
	var indeterminate *IndeterminateEffectError
	if !errors.As(err, &indeterminate) {
		t.Fatalf("setup: first invocation error = %T %v, want *IndeterminateEffectError", err, err)
	}

	// Reconciliation says: absent.
	bucket := newFakeBucket(t)
	firingControl(bucket)
	cfg := reconcileCfg(bucket)
	cfg.AbsentSamplesRequired = 2
	cfg.MaxAttempts = 4
	receipt := mustReconcile(t, cfg)
	if receipt.State != ReconcileNotPublished {
		t.Fatalf("setup: receipt state = %q, want %q", receipt.State, ReconcileNotPublished)
	}

	postsAfterReconcile := validator.count()

	// ARM 1 — the SAME approval, a FRESH session and a FRESH budget. `absent`
	// did not give it back.
	second, secondRecording := publishSession(
		t, base.Store, "ac16-retry", publishGrant(fixture.scope), handler)
	if got := second.grants[0].Budget; got != PublishCost {
		t.Fatalf("the retry session's budget = %d, want a FRESH %d", got, PublishCost)
	}
	_, _, retryErr := second.Invoke(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish, Scope: fixture.scope, Cost: PublishCost, Now: 60,
	}, fixture.payload)
	if !errors.Is(retryErr, store.ErrApprovalAlreadyConsumed) {
		t.Fatalf("AC16: retry after `not-published` error = %T %v, want store.ErrApprovalAlreadyConsumed",
			retryErr, retryErr)
	}
	var denial *DenialError
	if errors.As(retryErr, &denial) {
		t.Fatalf("AC16: the retry was refused by BUDGET (%s), not by the burnt attended stamp",
			denial.Decision.Label)
	}
	if len(secondRecording.effectIDs) != 0 {
		t.Errorf("AC16: the refused retry minted durable intents %v, want none",
			secondRecording.effectIDs)
	}
	if got := validator.count(); got != postsAfterReconcile {
		t.Fatalf("AC16: the refused retry POSTed (%d -> %d)", postsAfterReconcile, got)
	}

	// ARM 2 — THE POSITIVE CONTROL, and the reason arm 1 is not vacuous. A
	// retry carrying a NEW attended approval passes the same claim boundary and
	// reaches dispatch. Without this, "the retry is refused" would be equally
	// satisfied by a store that refuses everything.
	fresh := newPublishFixture(t, validator.origin(), "ac16-second").
		landApproval(t, base, "approve", approvalTimes{request: 70, decide: 71, expires: 200})
	freshHandler := loopbackHandler(t, RegistryPublishConfig{
		PublisherPath:   writePublisherScript(t, "reset", nil),
		PackageDir:      fresh.dir,
		Manifest:        fresh.manifest,
		RegistryOrigin:  validator.origin(),
		ValidatorOrigin: validator.origin(),
		Credential:      writeCredentialFile(t, ac10Sentinel),
		Approval:        fresh.approval,
		ExecTimeout:     20 * time.Second,
	})
	if fresh.approval.ApprovalRef == fixture.approval.ApprovalRef {
		t.Fatal("AC16 positive control: the 'new' approval is the same object as the burnt one")
	}
	third, _ := publishSession(t, base.Store, "ac16-reapproved", publishGrant(fresh.scope), freshHandler)
	_, _, freshErr := third.Invoke(context.Background(), EffectRequest{
		Effect: EffectRegistryPublish, Scope: fresh.scope, Cost: PublishCost, Now: 80,
	}, fresh.payload)
	if errors.Is(freshErr, store.ErrApprovalAlreadyConsumed) {
		t.Fatalf("AC16 positive control: a NEW attended approval was reported already-consumed")
	}
	if !errors.As(freshErr, &indeterminate) {
		t.Fatalf("AC16 positive control: error = %T %v, want the fixture's ambiguous dispatch",
			freshErr, freshErr)
	}
	if got := validator.count(); got != postsAfterReconcile+1 {
		t.Fatalf("AC16 positive control: POSTs %d -> %d, want exactly one more; the new "+
			"attended approval must be what re-authorizes, and it must be able to",
			postsAfterReconcile, got)
	}
}

// ---------------------------------------------------------------------------
// AC16a — the same-pass control's validity, in three distinguishable arms
// ---------------------------------------------------------------------------

func TestProbeControlValidityHasThreeDistinguishableArms(t *testing.T) {
	// (i) and (ii) are arm (iii)'s CONTROLS. Arm (iii) is the one that can fail
	// today; without (i) and (ii) a reconciler that answered
	// `probe-unavailable` unconditionally would pass it.

	t.Run("i_control200_target404_counts_as_absent", func(t *testing.T) {
		bucket := newFakeBucket(t)
		firingControl(bucket)
		receipt := mustReconcile(t, reconcileCfg(bucket))
		t.Logf("AC16a arm (i): %s", receipt)
		if receipt.State != ReconcileNotPublished {
			t.Fatalf("arm (i) state = %q, want %q (detail %q)",
				receipt.State, ReconcileNotPublished, receipt.Detail)
		}
		if got := receipt.Samples[0].Verdict; got != ProbeAbsent {
			t.Fatalf("arm (i) verdict = %q, want %q", got, ProbeAbsent)
		}
		if got := receipt.Samples[0].ControlStatus; got != http.StatusOK {
			t.Errorf("arm (i) control status = %d, want 200", got)
		}
	})

	t.Run("ii_control200_target200_resolves_by_hash", func(t *testing.T) {
		bucket := newFakeBucket(t)
		firingControl(bucket)
		bucket.put(targetKey(), bucketResponse{
			status: http.StatusOK,
			body:   metadataDocument(fixturePackage, fixtureVersion, reconcileExpected),
		})
		receipt := mustReconcile(t, reconcileCfg(bucket))
		t.Logf("AC16a arm (ii): %s", receipt)
		if receipt.State != ReconcileSucceededReconciled {
			t.Fatalf("arm (ii) state = %q, want %q (detail %q)",
				receipt.State, ReconcileSucceededReconciled, receipt.Detail)
		}
		if got := receipt.Samples[0].Verdict; got != ProbePresent {
			t.Fatalf("arm (ii) verdict = %q, want %q", got, ProbePresent)
		}
	})

	// Arm (iii): the control does NOT return 200-with-JSON, while the target
	// looks EXACTLY like a clean absence. Every fixture must resolve
	// `probe-unavailable`, and the load-bearing assertion is the NEGATIVE one:
	// the receipt must not read `not-published`, because `not-published` is
	// what re-authorizes an irreversible POST.
	arm3 := []struct {
		name  string
		setUp func(t *testing.T, b *fakeBucket)
		// wholeOriginDown says the whole origin is unreachable. MEASURED
		// (iteration 65, evaluator finding, reproduced by the controller): the
		// branch that refuses here is **C1** — "same-pass control fetch failed"
		// — NOT C4. Closing the origin kills the CONTROL request too, and the
		// control is examined first, so C1 wins before C4 is ever reached.
		//
		// Consequence, stated plainly because it is easy to misread this
		// fixture as a guard: an origin-down row is a COVERAGE BYSTANDER for
		// the control-validity mutations. Verified first-party — under both
		// MUT-SM-PROBE-NO-CONTROL and MUT-SM-PROBE-CONTROL-ALWAYS-OK the two
		// sibling rows (wrong_registry_origin_empty_bucket, control_403) red
		// while this one PASSES, because with C1-C3 neutered it simply falls
		// through to C4's target-fetch failure, which is uninformative for a
		// different reason. The arm's mutation-catching duty is discharged
		// entirely by those two siblings; this row documents that a wholly
		// dead origin is also never an absence, and TestTargetTransportFailure-
		// IsNeverAbsence is what actually pins the C4 path.
		wholeOriginDown bool
	}{
		{
			// The single most likely real misconfiguration: $AILANG_REGISTRY
			// points at a bucket that is reachable and EMPTY. Everything 404s,
			// including the known-positive control — which is exactly the
			// signal that the instrument is not working.
			name: "wrong_registry_origin_empty_bucket",
			setUp: func(_ *testing.T, _ *fakeBucket) {
				// Deliberately seed NOTHING: the control 404s with the same
				// measured NoSuchKey document the target does.
			},
		},
		{
			name: "control_403",
			setUp: func(_ *testing.T, b *fakeBucket) {
				b.put(controlKey(), bucketResponse{
					status: http.StatusForbidden, body: []byte("forbidden"),
				})
			},
		},
		{
			name: "unreachable_host",
			setUp: func(_ *testing.T, b *fakeBucket) {
				firingControl(b)
				b.close()
			},
			wholeOriginDown: true,
		},
	}

	for _, tc := range arm3 {
		t.Run("iii_probe_unavailable/"+tc.name, func(t *testing.T) {
			bucket := newFakeBucket(t)
			tc.setUp(t, bucket)
			cfg := reconcileCfg(bucket)
			cfg.AbsentSamplesRequired = 2
			cfg.MaxAttempts = 3
			cfg.RequestTimeout = 2 * time.Second

			receipt := mustReconcile(t, cfg)
			t.Logf("AC16a arm (iii) %s: %s", tc.name, receipt)

			// THE ASSERTION AC16a IS ABOUT.
			if receipt.State == ReconcileNotPublished {
				t.Fatalf("arm (iii) %s resolved %q under a control that never fired "+
					"(samples: %v). An absence believed on a broken instrument "+
					"re-authorizes an irreversible POST", tc.name, receipt.State, receipt.Samples)
			}
			if receipt.State != ReconcileProbeUnavailable {
				t.Fatalf("arm (iii) %s state = %q, want %q (detail %q)",
					tc.name, receipt.State, ReconcileProbeUnavailable, receipt.Detail)
			}
			if got := receipt.AbsentSamples; got != 0 {
				t.Errorf("arm (iii) %s decremented the bounded window by %d, want 0",
					tc.name, got)
			}
			if got := receipt.UninformativeSamples; got != cfg.MaxAttempts {
				t.Errorf("arm (iii) %s uninformative samples = %d, want all %d",
					tc.name, got, cfg.MaxAttempts)
			}
			for i, sample := range receipt.Samples {
				if sample.Verdict != ProbeUninformative {
					t.Errorf("arm (iii) %s sample %d verdict = %q, want %q",
						tc.name, i, sample.Verdict, ProbeUninformative)
				}
			}
			// The receipt must say a HUMAN is required, and must be
			// distinguishable in text from arm (i)'s absence.
			if !strings.Contains(receipt.Detail, "human required") {
				t.Errorf("arm (iii) %s detail %q does not require a human", tc.name, receipt.Detail)
			}
			if !tc.wholeOriginDown {
				// When the origin is UP, the refusal must be attributed to the
				// CONTROL, not to something about the target.
				if !strings.Contains(receipt.Samples[0].Reason, "control") {
					t.Errorf("arm (iii) %s first sample reason = %q, want it to name the control",
						tc.name, receipt.Samples[0].Reason)
				}
				// And the target must have looked like a clean absence, or this
				// arm is not testing what it claims: a broken control is only
				// interesting when believing it would have said `not-published`.
				if got := receipt.Samples[0].TargetStatus; got != http.StatusNotFound {
					t.Fatalf("arm (iii) %s target status = %d, want 404: this fixture does "+
						"not arm the threat it exists to test", tc.name, got)
				}
			}
		})
	}
}

// TestTargetTransportFailureIsNeverAbsence drives classifier branch C4
// end-to-end rather than as a unit: the CONTROL succeeds, and the TARGET's
// connection is aborted mid-response.
//
// It is the one failure mode arm (iii) of AC16a cannot reach. Both URLs are
// structurally forced onto one origin, so "the control works and the target
// does not" cannot be produced by pointing them at different hosts — only by a
// transport failure on one specific key. A reconciler that read a broken socket
// as silence, and silence as absence, would re-authorize an irreversible POST
// while its control was cheerfully green.
func TestTargetTransportFailureIsNeverAbsence(t *testing.T) {
	bucket := newFakeBucket(t)
	firingControl(bucket)
	bucket.put(targetKey(), bucketResponse{hijack: true})

	cfg := reconcileCfg(bucket)
	cfg.AbsentSamplesRequired = 2
	cfg.MaxAttempts = 3
	cfg.RequestTimeout = 2 * time.Second

	receipt := mustReconcile(t, cfg)
	t.Logf("C4 end-to-end: %s", receipt)
	if receipt.State == ReconcileNotPublished {
		t.Fatalf("a target whose connection was ABORTED resolved %q", receipt.State)
	}
	if receipt.State != ReconcileProbeUnavailable {
		t.Fatalf("state = %q, want %q (detail %q)",
			receipt.State, ReconcileProbeUnavailable, receipt.Detail)
	}
	if got := receipt.AbsentSamples; got != 0 {
		t.Errorf("absent samples = %d, want 0", got)
	}
	// The premise: the CONTROL fired in the same pass. Without this the arm
	// would be indistinguishable from a wholly unreachable origin, which is a
	// different (and already covered) fixture.
	if got := receipt.Samples[0].ControlStatus; got != http.StatusOK {
		t.Fatalf("the control did not fire (status %d); this arm is not testing "+
			"a TARGET-only transport failure", got)
	}
	if !strings.Contains(receipt.Samples[0].Reason, "target fetch failed") {
		t.Errorf("sample reason = %q, want it to blame the target fetch",
			receipt.Samples[0].Reason)
	}
}

// TestAbsenceControlSharesTargetKeySpace is the test registry_reconcile.go's
// metadataObjectURL comment names. It is the structural half of AC16a: a
// control that does not travel the target's own key-space proves nothing about
// the target's key-space.
//
// The failure it forbids is MEASURED, not hypothetical. On 2026-08-08 the
// validator origin answered 200 with 35457 bytes of well-formed JSON at
// /api/packages while answering 404 at
// /packages/{vendor}/{name}/{version}/metadata.json. So a control implemented
// as "fetch the registry index" — the natural choice, and the one a failover
// or a copy-paste produces — would FIRE against a misconfigured origin while
// the target 404s, the sample would be believed absent, and an irreversible
// POST would be re-authorized.
func TestAbsenceControlSharesTargetKeySpace(t *testing.T) {
	const origin = "https://storage.googleapis.com/ailang-registry"

	target, err := metadataObjectURL(origin, fixtureVendor, fixtureName, fixtureVersion)
	if err != nil {
		t.Fatal(err)
	}
	control, err := metadataObjectURL(origin,
		DefaultProbeControlVendor, DefaultProbeControlName, DefaultProbeControlVersion)
	if err != nil {
		t.Fatal(err)
	}

	// (a) Both are built from the SAME origin, in the SAME key-space, and the
	// production construction agrees with this file's independent one.
	if got, want := target, origin+bucketKey(fixtureVendor, fixtureName, fixtureVersion); got != want {
		t.Errorf("target URL = %q, want %q", got, want)
	}
	if got, want := control, origin+bucketKey(
		DefaultProbeControlVendor, DefaultProbeControlName, DefaultProbeControlVersion); got != want {
		t.Errorf("control URL = %q, want %q", got, want)
	}
	for _, arm := range []struct{ what, url string }{{"target", target}, {"control", control}} {
		if !strings.HasPrefix(arm.url, origin+"/packages/") {
			t.Errorf("%s URL %q does not live in the bucket's object key-space", arm.what, arm.url)
		}
		if !strings.HasSuffix(arm.url, "/metadata.json") {
			t.Errorf("%s URL %q is not a metadata object", arm.what, arm.url)
		}
	}
	// (b) The MEASURED index route is exactly what the control must NOT be.
	if strings.Contains(control, "/api/") {
		t.Errorf("the control URL %q is an index route; a 200 there says nothing about "+
			"the bucket key the target reads", control)
	}
	// (c) They differ ONLY in vendor/name/version.
	if control == target {
		t.Fatal("the control is the target itself")
	}
	if strings.TrimSuffix(target, bucketKey(fixtureVendor, fixtureName, fixtureVersion)) !=
		strings.TrimSuffix(control, bucketKey(
			DefaultProbeControlVendor, DefaultProbeControlName, DefaultProbeControlVersion)) {
		t.Errorf("control %q and target %q do not share a common origin prefix", control, target)
	}

	// (d) THE STRUCTURAL GUARANTEE. metadataObjectURL is the only URL builder
	// because ReconcileConfig carries no field a caller could use to hand
	// reconciliation a ready-made URL. Enumerate the struct and prove it.
	cfgType := reflect.TypeOf(ReconcileConfig{})
	if cfgType.NumField() == 0 {
		t.Fatal("ReconcileConfig enumerated zero fields: this guard is not reading anything")
	}
	var offenders []string
	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgType.Field(i)
		name := strings.ToLower(field.Name)
		if field.Type.Kind() != reflect.String {
			continue
		}
		if strings.Contains(name, "url") || strings.Contains(name, "uri") ||
			strings.Contains(name, "endpoint") {
			offenders = append(offenders, field.Name)
		}
	}
	if len(offenders) != 0 {
		t.Errorf("ReconcileConfig exposes URL-shaped field(s) %v; a caller could then aim the "+
			"control at a different key-space than the target", offenders)
	}
	t.Logf("ReconcileConfig enumerated %d fields, %d URL-shaped", cfgType.NumField(), len(offenders))

	// (e) The known-positive control for (d): the same enumeration, run over a
	// struct that DOES carry such a field, must find it. Without this the empty
	// `offenders` above is compatible with the loop never running.
	type controlShape struct {
		RegistryOrigin string
		ControlURL     string
	}
	found := 0
	probe := reflect.TypeOf(controlShape{})
	for i := 0; i < probe.NumField(); i++ {
		if strings.Contains(strings.ToLower(probe.Field(i).Name), "url") {
			found++
		}
	}
	if found != 1 {
		t.Fatalf("the URL-field detector found %d fields in a struct that has exactly 1; "+
			"the assertion in (d) is not measuring anything", found)
	}
}

// ---------------------------------------------------------------------------
// AC16b — a non-JSON error body is an ERROR PATH, not a parse crash, and the
//         measured NoSuchKey document is the ONLY body absence is believed on
// ---------------------------------------------------------------------------

func TestErrorBodyClassificationDistinguishesNoSuchKeyFromGarbage(t *testing.T) {
	// The two arms DIFFERING is the whole point. A classifier that redded both
	// (or greened both) could not tell an absence from a broken response, which
	// is the failure this criterion exists to forbid.
	arms := []struct {
		name      string
		body      []byte
		wantState string
		want      string
	}{
		{
			name:      "measured_gcs_nosuchkey_xml",
			body:      []byte(gcsNoSuchKey),
			wantState: ReconcileNotPublished,
			want:      ProbeAbsent,
		},
		{
			name:      "truncated_nosuchkey_xml",
			body:      []byte(gcsNoSuchKey)[:40],
			wantState: ReconcileProbeUnavailable,
			want:      ProbeUninformative,
		},
		{
			name:      "garbage_bytes",
			body:      []byte("\x00\x01\x02not xml at all\xff"),
			wantState: ReconcileProbeUnavailable,
			want:      ProbeUninformative,
		},
		{
			name:      "validator_plain_text_404",
			body:      []byte(validatorPlainText404),
			wantState: ReconcileProbeUnavailable,
			want:      ProbeUninformative,
		},
		{
			name:      "empty_body",
			body:      nil,
			wantState: ReconcileProbeUnavailable,
			want:      ProbeUninformative,
		},
	}

	seen := map[string]bool{}
	for _, arm := range arms {
		t.Run(arm.name, func(t *testing.T) {
			bucket := newFakeBucket(t)
			firingControl(bucket)
			bucket.put(targetKey(), bucketResponse{
				status: http.StatusNotFound, body: arm.body,
			})
			receipt := mustReconcile(t, reconcileCfg(bucket))
			t.Logf("AC16b %s: %s", arm.name, receipt)

			// The control fires in EVERY arm, so the only thing that varies is
			// the target's body.
			if got := receipt.Samples[0].ControlStatus; got != http.StatusOK {
				t.Fatalf("AC16b %s: the control did not fire (status %d); this arm is not "+
					"testing the body classifier", arm.name, got)
			}
			if got := receipt.Samples[0].Verdict; got != arm.want {
				t.Fatalf("AC16b %s: verdict = %q, want %q (reason %q)",
					arm.name, got, arm.want, receipt.Samples[0].Reason)
			}
			if receipt.State != arm.wantState {
				t.Fatalf("AC16b %s: state = %q, want %q (detail %q)",
					arm.name, receipt.State, arm.wantState, receipt.Detail)
			}
			if arm.want == ProbeUninformative && receipt.State == ReconcileNotPublished {
				t.Fatalf("AC16b %s: a body that is not the measured absence document "+
					"resolved `not-published`", arm.name)
			}
		})
		seen[arm.wantState] = true
	}
	// The arms must actually DIFFER.
	if len(seen) < 2 {
		t.Fatalf("AC16b: every arm expected the same state %v; the fixture cannot tell an "+
			"absence from a broken response", seen)
	}
}

// ---------------------------------------------------------------------------
// AC17 — replay returns the RECORDED publish result and performs ZERO network
//        calls
// ---------------------------------------------------------------------------

func TestReplayReturnsTheRecordedPublishResultWithZeroNetworkCalls(t *testing.T) {
	validator := newFakeValidator(t, "ok")
	base := openPublishStore(t)
	fixture := newPublishFixture(t, validator.origin(), "ac17").
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
	request := EffectRequest{
		Effect: EffectRegistryPublish, Scope: fixture.scope, Cost: PublishCost, Now: now,
	}

	live, _ := publishSession(t, base.Store, "ac17-live", publishGrant(fixture.scope), handler)
	liveResult, liveRef, err := live.Invoke(context.Background(), request, fixture.payload)
	if err != nil {
		t.Fatalf("AC17 setup: live publish: %v", err)
	}
	if liveRef.IsZero() || len(liveResult) == 0 {
		t.Fatalf("AC17 setup: live publish produced ref %s and %d bytes", liveRef, len(liveResult))
	}
	livePosts, liveDispatches := validator.count(), handler.Dispatches()
	if livePosts != 1 || liveDispatches != 1 {
		t.Fatalf("AC17 setup: live POSTs=%d dispatches=%d, want 1 and 1", livePosts, liveDispatches)
	}
	// The recorded result must actually BE a publish result, or "replay returns
	// the recorded publish result" is a claim about opaque bytes.
	liveStatus := publishResultStatus(t, liveResult)
	if liveStatus != PublishStatusSucceeded {
		t.Fatalf("AC17 setup: live result status = %q, want %q", liveStatus, PublishStatusSucceeded)
	}

	// THE INSTRUMENT: a counting handler that ALSO reaches the network. It is
	// registered in the replay session's registry under the same effect name a
	// live session would use, so a replay that consults its registry is
	// observable twice over.
	instrument := newCountingNetworkHandler(t)
	replaySession := NewReplaySession(base.Store,
		[]Capability{publishGrant(fixture.scope)},
		Registry{EffectRegistryPublish: instrument},
		[]hashref.HashRef{liveRef})

	replayed, replayRef, err := replaySession.Invoke(context.Background(), request, fixture.payload)
	if err != nil {
		t.Fatalf("AC17: replay Invoke: %v", err)
	}

	// (a) it returned THE RECORDED RESULT.
	if replayRef != liveRef {
		t.Errorf("AC17: replay record ref = %s, want the recorded %s", replayRef, liveRef)
	}
	if string(replayed) != string(liveResult) {
		t.Fatalf("AC17: replayed %d bytes %q, recorded %d bytes %q",
			len(replayed), replayed, len(liveResult), liveResult)
	}
	if got := publishResultStatus(t, replayed); got != PublishStatusSucceeded {
		t.Errorf("AC17: replayed result status = %q, want %q", got, PublishStatusSucceeded)
	}

	// (b) it performed ZERO network calls.
	if got := instrument.calls.Load(); got != 0 {
		t.Fatalf("AC17: replay consulted a handler %d time(s), want 0 (%s)", got, instrument)
	}
	if got := instrument.bucket.count(); got != 0 {
		t.Fatalf("AC17: replay caused %d network request(s), want 0 (%s)", got, instrument)
	}
	if got := validator.count(); got != livePosts {
		t.Fatalf("AC17: replay moved the POST counter %d -> %d", livePosts, got)
	}
	if got := handler.Dispatches(); got != liveDispatches {
		t.Fatalf("AC17: replay moved the dispatch counter %d -> %d", liveDispatches, got)
	}
	if got := handler.CredentialLoads(); got != 1 {
		t.Errorf("AC17: credential loads = %d, want the 1 from the live publish only", got)
	}

	// (c) THE SAME-RUN POSITIVE CONTROL for the zeros above.
	instrument.provePositiveControl(t)
}

// publishResultStatus reads the status field back out of a canonical publish
// result object. It decodes the OBJECT rather than re-deriving the status from
// the inputs, so it verifies the artifact and not the test's arithmetic.
func publishResultStatus(t *testing.T, payload []byte) string {
	t.Helper()
	var wire struct {
		Schema string `json:"schema"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(payload, &wire); err != nil {
		t.Fatalf("decode publish result %q: %v", payload, err)
	}
	if wire.Schema != PublishResultSchema {
		t.Fatalf("publish result schema = %q, want %q", wire.Schema, PublishResultSchema)
	}
	return wire.Status
}

// ---------------------------------------------------------------------------
// Refusal-branch coverage: R1-R9, U1-U3, C1-C8, P1-P3
//
// Each refusal in registry_reconcile.go is a one-way door placed there because
// its failure mode is a permanent duplicate public artifact. A door nobody has
// ever walked into is a door nobody knows is there.
// ---------------------------------------------------------------------------

func TestReconcileURLSegmentRefusalsU1U3(t *testing.T) {
	const origin = "http://127.0.0.1:1"
	// U3's positive control runs FIRST: an ordinary segment must build.
	if _, err := metadataObjectURL(origin, "world", "core", "0.1.0"); err != nil {
		t.Fatalf("positive control: a plain segment was refused: %v", err)
	}

	cases := []struct {
		branch, vendor, name, version, want string
	}{
		{"U1/vendor", "", "core", "0.1.0", "vendor is empty"},
		{"U1/name", "world", "", "0.1.0", "name is empty"},
		{"U1/version", "world", "core", "", "version is empty"},
		{"U2/dot", ".", "core", "0.1.0", "path traversal segment"},
		{"U2/dotdot", "world", "..", "0.1.0", "path traversal segment"},
		{"U3/slash", "world", "core", "0.1.0/../..", "not a single safe path segment"},
		{"U3/space", "wo rld", "core", "0.1.0", "not a single safe path segment"},
		{"U3/percent", "world", "co%2fre", "0.1.0", "not a single safe path segment"},
	}
	for _, tc := range cases {
		t.Run(tc.branch, func(t *testing.T) {
			got, err := metadataObjectURL(origin, tc.vendor, tc.name, tc.version)
			if err == nil {
				t.Fatalf("%s: built %q, want a refusal", tc.branch, got)
			}
			var refusal *PublishRefusalError
			if !errors.As(err, &refusal) {
				t.Fatalf("%s: error = %T %v, want *PublishRefusalError", tc.branch, err, err)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("%s: error %q does not contain %q", tc.branch, err.Error(), tc.want)
			}
		})
	}
}

func TestReconcileRefusalSetWithAPassingPositiveControl(t *testing.T) {
	// The positive control: the SAME entry point, with everything valid, takes
	// a sample and returns a receipt. Every refusal below is therefore a
	// refusal of the input, not of the call.
	control := newFakeBucket(t)
	firingControl(control)
	if got := mustReconcile(t, reconcileCfg(control)); got.State != ReconcileNotPublished {
		t.Fatalf("positive control: state = %q, want %q", got.State, ReconcileNotPublished)
	}

	// A dead loopback origin: R3-R7 all refuse BEFORE the first fetch, so
	// nothing here can reach a network even if a refusal regressed.
	dead := "http://127.0.0.1:1"
	valid := func() ReconcileConfig {
		return ReconcileConfig{
			RegistryOrigin: dead,
			Vendor:         fixtureVendor, Name: fixtureName, Version: fixtureVersion,
			Expected:              reconcileExpected,
			AbsentSamplesRequired: 1, MaxAttempts: 2,
			RequestTimeout: time.Second,
		}
	}

	loopbackCases := []struct {
		branch string
		mutate func(cfg *ReconcileConfig)
		want   string
	}{
		{"R1/non-loopback", func(c *ReconcileConfig) {
			c.RegistryOrigin = "https://storage.googleapis.com/ailang-registry"
		}, "refuses non-loopback"},
		{"R1/trailing-slash", func(c *ReconcileConfig) {
			c.RegistryOrigin = dead + "/"
		}, "must not end in a slash"},
		{"R1/userinfo", func(c *ReconcileConfig) {
			c.RegistryOrigin = "http://user:pass@127.0.0.1:1"
		}, "embeds credentials in the URL"},
		{"R1/query", func(c *ReconcileConfig) {
			c.RegistryOrigin = dead + "?a=b"
		}, "carries a query or fragment"},
		{"R1/wildcard", func(c *ReconcileConfig) {
			c.RegistryOrigin = "http://*.127.0.0.1"
		}, "is a wildcard"},
		{"R1/empty", func(c *ReconcileConfig) { c.RegistryOrigin = "" }, "is empty"},
		{"R3/target", func(c *ReconcileConfig) { c.Version = ".." }, "metadata object version"},
		{"R4/control", func(c *ReconcileConfig) {
			c.ControlVendor, c.ControlName, c.ControlVersion = "..", "auth", "0.4.1"
		}, "metadata object vendor"},
		{"R5/control-is-target", func(c *ReconcileConfig) {
			c.ControlVendor, c.ControlName, c.ControlVersion =
				fixtureVendor, fixtureName, fixtureVersion
		}, "is the probe target itself"},
		{"R6/tarball", func(c *ReconcileConfig) {
			c.Expected.TarballSHA256 = ""
		}, "requires all three expected digests"},
		{"R6/content", func(c *ReconcileConfig) {
			c.Expected.ContentHash = ""
		}, "requires all three expected digests"},
		{"R6/interface", func(c *ReconcileConfig) {
			c.Expected.InterfaceHash = ""
		}, "requires all three expected digests"},
		{"R7/zero-window", func(c *ReconcileConfig) {
			c.AbsentSamplesRequired, c.MaxAttempts = -1, 4
		}, "absence window is unsatisfiable"},
		{"R7/attempts-below-window", func(c *ReconcileConfig) {
			c.AbsentSamplesRequired, c.MaxAttempts = 3, 2
		}, "absence window is unsatisfiable"},
	}
	for _, tc := range loopbackCases {
		t.Run(tc.branch, func(t *testing.T) {
			cfg := valid()
			tc.mutate(&cfg)
			receipt, err := reconcileLoopback(context.Background(), cfg)
			if err == nil {
				t.Fatalf("%s: returned receipt %s, want a refusal", tc.branch, receipt)
			}
			var refusal *PublishRefusalError
			if !errors.As(err, &refusal) {
				t.Fatalf("%s: error = %T %v, want *PublishRefusalError", tc.branch, err, err)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("%s: error %q does not contain %q", tc.branch, err.Error(), tc.want)
			}
		})
	}

	// R1 and R2 through the PRODUCTION entry point. Both refuse strictly
	// before the first fetch, so no request leaves this process — which is
	// asserted below by the fact that a refusal, not a receipt, comes back.
	productionCases := []struct{ branch, origin, want string }{
		{"R1/production-requires-https", "http://storage.googleapis.com/ailang-registry",
			"must be https for a live origin"},
		{"R1/production-refuses-loopback", "https://127.0.0.1:1",
			"production constructor refuses a loopback"},
		{"R2/validator-origin", ApprovedValidatorOrigin,
			"reconciliation reads the bucket, not the validator service"},
	}
	for _, tc := range productionCases {
		t.Run(tc.branch, func(t *testing.T) {
			cfg := valid()
			cfg.RegistryOrigin = tc.origin
			receipt, err := ReconcileRegistryPublish(context.Background(), cfg)
			if err == nil {
				t.Fatalf("%s: returned receipt %s, want a refusal", tc.branch, receipt)
			}
			if len(receipt.Samples) != 0 {
				t.Fatalf("%s: the refusal took %d sample(s); it must refuse before any fetch",
					tc.branch, len(receipt.Samples))
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("%s: error %q does not contain %q", tc.branch, err.Error(), tc.want)
			}
		})
	}
}

// TestReconcileCancellationReportsTheRefusalItReached drives R8. The context is
// cancelled BY THE BUCKET, from inside the first target fetch, so the arm is
// deterministic rather than a race against a sleep.
func TestReconcileCancellationReportsTheRefusalItReached(t *testing.T) {
	bucket := newFakeBucket(t)
	firingControl(bucket)
	// A 404 body that is NOT the absence document, so every sample is
	// uninformative and the loop would otherwise run to MaxAttempts.
	bucket.put(targetKey(), bucketResponse{
		status: http.StatusNotFound, body: []byte(validatorPlainText404),
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	target := targetKey()
	bucket.onGet = func(path string) {
		if path == target {
			cancel()
		}
	}

	cfg := reconcileCfg(bucket)
	cfg.AbsentSamplesRequired = 1
	cfg.MaxAttempts = 6
	// Large enough that ctx.Done() wins the select deterministically.
	cfg.SampleInterval = 30 * time.Second

	start := time.Now()
	receipt, err := reconcileLoopback(ctx, cfg)
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("R8: cancellation returned an error rather than a receipt: %v", err)
	}
	t.Logf("R8 receipt after %s: %s", elapsed, receipt)
	if receipt.State != ReconcileProbeUnavailable {
		t.Fatalf("R8: a cancelled pass resolved %q; it must never upgrade to a resolution",
			receipt.State)
	}
	if !strings.Contains(receipt.Detail, "cancelled") {
		t.Errorf("R8: detail %q does not say the pass was cancelled", receipt.Detail)
	}
	if got := len(receipt.Samples); got != 1 {
		t.Errorf("R8: samples = %d, want the 1 taken before cancellation", got)
	}
	// It stopped at the FIRST interval rather than sleeping through it. Stated
	// against a bound 1000x smaller than the interval, so it is not a
	// noise-sized threshold.
	if elapsed >= cfg.SampleInterval {
		t.Errorf("R8: the pass took %s, i.e. it slept through the %s interval",
			elapsed, cfg.SampleInterval)
	}
}

// TestClassifyProbeSampleCoversEveryBranch drives all eight classifier
// branches directly. Six of them decline to call the sample informative; the
// table asserts each one's verdict AND a distinguishing fragment of its reason,
// so two branches cannot silently collapse into one.
func TestClassifyProbeSampleCoversEveryBranch(t *testing.T) {
	okJSON := metadataDocument("sunholo/auth", "0.4.1", reconcileExpected)
	cases := []struct {
		branch          string
		control, target probeResult
		wantVerdict     string
		wantReason      string
	}{
		{"C1/control-error",
			probeResult{Err: errors.New("dial tcp: connection refused")},
			probeResult{Status: 404, Body: []byte(gcsNoSuchKey)},
			ProbeUninformative, "same-pass control fetch failed"},
		{"C2/control-non-200",
			probeResult{Status: 403, Body: []byte("forbidden")},
			probeResult{Status: 404, Body: []byte(gcsNoSuchKey)},
			ProbeUninformative, "same-pass control returned 403"},
		{"C3/control-not-json",
			probeResult{Status: 200, Body: []byte("<html>captive portal</html>")},
			probeResult{Status: 404, Body: []byte(gcsNoSuchKey)},
			ProbeUninformative, "control body is not a well-formed JSON object"},
		{"C3/control-empty-json-object",
			probeResult{Status: 200, Body: []byte("{}")},
			probeResult{Status: 404, Body: []byte(gcsNoSuchKey)},
			ProbeUninformative, "control body is not a well-formed JSON object"},
		{"C4/target-error",
			probeResult{Status: 200, Body: okJSON},
			probeResult{Err: errors.New("unexpected EOF")},
			ProbeUninformative, "target fetch failed"},
		{"C5/absent",
			probeResult{Status: 200, Body: okJSON},
			probeResult{Status: 404, Body: []byte(gcsNoSuchKey)},
			ProbeAbsent, "measured GCS NoSuchKey document and a firing control"},
		{"C6/unexplained-404",
			probeResult{Status: 200, Body: okJSON},
			probeResult{Status: 404, Body: []byte(validatorPlainText404)},
			ProbeUninformative, "404 body is not the measured GCS NoSuchKey document"},
		{"C7/target-200-not-json",
			probeResult{Status: 200, Body: okJSON},
			probeResult{Status: 200, Body: []byte("not json")},
			ProbeUninformative, "target 200 body is not a well-formed JSON object"},
		{"C7/present",
			probeResult{Status: 200, Body: okJSON},
			probeResult{Status: 200, Body: okJSON},
			ProbePresent, "target 200 with a well-formed metadata document"},
		{"C8/target-500",
			probeResult{Status: 200, Body: okJSON},
			probeResult{Status: 500, Body: []byte("boom")},
			ProbeUninformative, "target returned 500"},
		{"C8/target-403",
			probeResult{Status: 200, Body: okJSON},
			probeResult{Status: 403, Body: []byte("nope")},
			ProbeUninformative, "target returned 403"},
	}
	reasons := map[string]bool{}
	for _, tc := range cases {
		t.Run(tc.branch, func(t *testing.T) {
			verdict, reason := classifyProbeSample(tc.control, tc.target)
			if verdict != tc.wantVerdict {
				t.Fatalf("%s: verdict = %q (%q), want %q", tc.branch, verdict, reason, tc.wantVerdict)
			}
			if !strings.Contains(reason, tc.wantReason) {
				t.Errorf("%s: reason %q does not contain %q", tc.branch, reason, tc.wantReason)
			}
			reasons[reason] = true
		})
	}
	// Distinctness: 11 rows across 8 branches must not collapse to one message.
	if len(reasons) < 8 {
		t.Errorf("the classifier produced %d distinct reasons across %d rows; branches are "+
			"indistinguishable in the receipt", len(reasons), len(cases))
	}
	// The absence branch is the only path to ProbeAbsent.
	absent := 0
	for _, tc := range cases {
		if v, _ := classifyProbeSample(tc.control, tc.target); v == ProbeAbsent {
			absent++
		}
	}
	if absent != 1 {
		t.Errorf("%d of %d rows reached %q; exactly one (C5) should",
			absent, len(cases), ProbeAbsent)
	}
}

// TestResolvePresentRefusalsP1P3 drives the three branches of resolvePresent.
// P1 is the interesting one: the classifier has already proved the body parses
// as a JSON OBJECT, so the only way the second decoder can fail is a member of
// the wrong TYPE — which a registry could serve tomorrow.
func TestResolvePresentRefusalsP1P3(t *testing.T) {
	base := ReconcileReceipt{State: ReconcileProbeUnavailable}
	cfg := ReconcileConfig{
		Vendor: fixtureVendor, Name: fixtureName, Version: fixtureVersion,
		Expected: reconcileExpected,
	}

	t.Run("P1/typed-decode-failure", func(t *testing.T) {
		body := []byte(`{"name":123,"version":"0.1.0"}`)
		// The premise: the classifier WOULD have accepted this body.
		if !wellFormedJSONObject(body) {
			t.Fatal("P1 premise broken: the classifier would have rejected this body, so " +
				"resolvePresent could never see it")
		}
		got, err := resolvePresent(base, body, cfg)
		if err != nil {
			t.Fatal(err)
		}
		if got.State != ReconcileProbeUnavailable {
			t.Fatalf("P1: state = %q, want %q", got.State, ReconcileProbeUnavailable)
		}
		if !strings.Contains(got.Detail, "did not decode") {
			t.Errorf("P1: detail %q does not say the document failed to decode", got.Detail)
		}
	})

	t.Run("P2/different-package", func(t *testing.T) {
		got, err := resolvePresent(base,
			metadataDocument("sunholo/auth", fixtureVersion, reconcileExpected), cfg)
		if err != nil {
			t.Fatal(err)
		}
		if got.State != ReconcileConflict {
			t.Fatalf("P2: state = %q, want %q", got.State, ReconcileConflict)
		}
		if !strings.Contains(got.Detail, "sunholo/auth") {
			t.Errorf("P2: detail %q does not name the document it read", got.Detail)
		}
	})

	t.Run("P2/different-version", func(t *testing.T) {
		got, err := resolvePresent(base,
			metadataDocument(fixturePackage, "0.9.9", reconcileExpected), cfg)
		if err != nil {
			t.Fatal(err)
		}
		if got.State != ReconcileConflict {
			t.Fatalf("P2: state = %q, want %q", got.State, ReconcileConflict)
		}
	})

	for _, arm := range []struct {
		name   string
		served PublishHashes
		want   string
	}{
		{"P3/tarball", PublishHashes{
			flipLastNibble(reconcileExpected.TarballSHA256),
			reconcileExpected.ContentHash, reconcileExpected.InterfaceHash}, "tarball"},
		{"P3/content", PublishHashes{
			reconcileExpected.TarballSHA256,
			flipLastNibble(reconcileExpected.ContentHash), reconcileExpected.InterfaceHash}, "content"},
		{"P3/interface", PublishHashes{
			reconcileExpected.TarballSHA256, reconcileExpected.ContentHash,
			flipLastNibble(reconcileExpected.InterfaceHash)}, "interface"},
	} {
		t.Run(arm.name, func(t *testing.T) {
			got, err := resolvePresent(base,
				metadataDocument(fixturePackage, fixtureVersion, arm.served), cfg)
			if err != nil {
				t.Fatal(err)
			}
			if got.State != ReconcileConflict {
				t.Fatalf("%s: state = %q, want %q", arm.name, got.State, ReconcileConflict)
			}
			if !strings.Contains(got.Detail, arm.want) {
				t.Errorf("%s: detail %q does not name the %s arm", arm.name, got.Detail, arm.want)
			}
		})
	}

	// The same-run positive control: unmodified digests resolve success, so the
	// three conflicts above are conflicts about the digests and not about the
	// comparator being unable to agree with anything.
	t.Run("P3/positive-control", func(t *testing.T) {
		got, err := resolvePresent(base,
			metadataDocument(fixturePackage, fixtureVersion, reconcileExpected), cfg)
		if err != nil {
			t.Fatal(err)
		}
		if got.State != ReconcileSucceededReconciled {
			t.Fatalf("P3 positive control: state = %q, want %q (detail %q)",
				got.State, ReconcileSucceededReconciled, got.Detail)
		}
	})
}

// TestGCSNoSuchKeyDecoderRejectsImpersonations pins isGCSNoSuchKey as a DECODE
// rather than a substring match. The last row is the one that matters: a body
// that CONTAINS the string but is not the document.
func TestGCSNoSuchKeyDecoderRejectsImpersonations(t *testing.T) {
	if !isGCSNoSuchKey([]byte(gcsNoSuchKey)) {
		t.Fatal("positive control: the measured NoSuchKey document was not recognised")
	}
	rejects := map[string][]byte{
		"truncated":           []byte(gcsNoSuchKey)[:40],
		"empty":               nil,
		"plain text":          []byte(validatorPlainText404),
		"json":                []byte(`{"Code":"NoSuchKey"}`),
		"wrong code":          []byte(`<Error><Code>AccessDenied</Code></Error>`),
		"wrong root element":  []byte(`<Result><Code>NoSuchKey</Code></Result>`),
		"substring in a blob": []byte("garbage NoSuchKey garbage"),
		"nested but not root": []byte(`<Wrapper><Error><Code>NoSuchKey</Code></Error></Wrapper>`),
	}
	for name, body := range rejects {
		if isGCSNoSuchKey(body) {
			t.Errorf("%s: %q was accepted as the measured absence document", name, body)
		}
	}
}

// TestEveryReconcileRequestWasLoopbackAndAGet enumerates every request every
// fake bucket in this file served, across the whole package run. It is the
// file-level safety enumeration: "read-only and loopback-only" measured over
// the requests that actually happened, not asserted per test.
//
// It is ordered last by name deliberately — `go test` runs tests in source
// order within a file, and this file's remaining tests all precede it.
func TestEveryReconcileRequestWasLoopbackAndAGet(t *testing.T) {
	bucketAudit.mu.Lock()
	origins := append([]string(nil), bucketAudit.origins...)
	requests := append([]bucketRequest(nil), bucketAudit.requests...)
	bucketAudit.mu.Unlock()

	if len(origins) == 0 || len(requests) == 0 {
		t.Skipf("no fake bucket traffic in this run (origins=%d requests=%d); "+
			"this enumeration is only meaningful after the reconcile tests have run",
			len(origins), len(requests))
	}
	for _, origin := range origins {
		if !isLoopbackHost(mustHost(t, origin)) {
			t.Fatalf("a fake bucket bound %q, which is not loopback", origin)
		}
	}
	for _, req := range requests {
		if req.Method != http.MethodGet {
			t.Errorf("a reconciliation issued %s %s; its only verb is GET", req.Method, req.Path)
		}
		if !strings.HasPrefix(req.Path, "/packages/") {
			t.Errorf("a reconciliation fetched %q, outside the bucket object key-space", req.Path)
		}
	}
	t.Logf("enumerated %d requests across %d loopback fake buckets", len(requests), len(origins))
}
