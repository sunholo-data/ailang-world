package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/evidence"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// demoProofRequired are the identities world/logepoch.ail verifies under the
// pinned v0.41.0 ai-check.
var demoProofRequired = []string{"sameRef", "servesEntry"}

// requireProofPin fails, never skips, when the v0.41.0 pin is absent.
func requireProofPin(t *testing.T) string {
	t.Helper()
	pin := os.Getenv("AILANG_BIN")
	if pin == "" {
		t.Fatal("AILANG_BIN unset: the proof adapters need the pinned v0.41.0 release; never skip")
	}
	requireBinaryRelease(t, pin, "v0.41.0")
	return pin
}

// demoProofRig is a real daemon over a fresh store with one genesis world. The
// store handle is closed before the daemon starts; from then on every read and
// write in a test goes through HTTP.
type demoProofRig struct {
	url, token string
}

func newDemoProofRig(t *testing.T, pin string) demoProofRig {
	t.Helper()
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	db := filepath.Join(base, "state", "world.db")
	if err := os.MkdirAll(filepath.Dir(db), 0o755); err != nil {
		t.Fatal(err)
	}
	pinBytes, err := os.ReadFile(pin)
	if err != nil {
		t.Fatal(err)
	}
	interp := hashref.SumSHA256(pinBytes)
	st, err := store.Open(db)
	if err != nil {
		t.Fatal(err)
	}
	obj := store.Object{Hash: hashref.SumSHA256([]byte("genesis-state")), InterfaceHash: hashref.SumSHA256([]byte("test/genesis")),
		SemanticID: "test/genesis", Provenance: "demo_proof_test", Payload: []byte("genesis-state")}
	entryHash := hashref.SumSHA256([]byte("genesis-entry"))
	w := store.World{Ref: hashref.SumSHA256([]byte("genesis-world")), Revision: 0, StateRoot: obj.Hash, LogHead: entryHash}
	if err := st.Commit(boundedTestContext(t), store.Commit{Objects: []store.Object{obj}, NextWorld: w, Entry: store.LogEntry{Header: store.LogHeader{
		EntryIndex: 0, SemanticsEpoch: 1, TransitionFn: hashref.SumSHA256([]byte("genesis-fn")), Interpreter: interp,
		PrevEntryHash: hashref.SumSHA256([]byte("genesis-prev")), WrittenBy: "demo_proof_test"}, EntryHash: entryHash, TransitionRef: obj.Hash}}); err != nil {
		t.Fatalf("genesis: %v", err)
	}
	now := time.Now().Unix()
	grants := []broker.Capability{{Effect: "world.apply", Scope: "world", ExpiresAt: now + 7200, Budget: 1000}}
	tok, _, _, err := authority.Mint(boundedTestContext(t), credsFor(t, db), "demo", grants, 3600, now, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	d, err := daemon.New(boundedTestContext(t), daemon.Config{DBPath: db, BindHost: daemon.DefaultBindHost, AilangBin: pin, ErrorLog: io.Discard})
	if err != nil {
		t.Fatalf("daemon.New: %v", err)
	}
	srv := httptest.NewServer(d.Handler())
	t.Cleanup(func() {
		srv.Close()
		_ = d.Close()
	})
	return demoProofRig{url: srv.URL, token: tok}
}

// client returns a REST client with the M3a 10 s request bound and the given transport.
func (r demoProofRig) client(rt http.RoundTripper) *client {
	return &client{base: r.url, timeout: demoProofHTTPTimeout, http: &http.Client{Transport: rt}}
}

// forwardRT forwards every request to the real daemon. A non-empty flip names
// the object path whose returned payload gets one byte flipped, inside an
// otherwise valid HTTP 200 with valid JSON.
type forwardRT struct{ flip string }

func (f forwardRT) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil || f.flip == "" || req.URL.Path != "/v1/objects/"+f.flip {
		return resp, err
	}
	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(body, &fields); err != nil {
		return nil, err
	}
	var payload []byte
	if err := json.Unmarshal(fields["payload"], &payload); err != nil || len(payload) == 0 {
		return nil, fmt.Errorf("forwardRT: no payload to flip: %v", err)
	}
	payload[len(payload)-1] ^= 0xff
	fields["payload"], _ = json.Marshal(payload)
	out, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}
	resp.Body = io.NopCloser(bytes.NewReader(out))
	resp.ContentLength = int64(len(out))
	return resp, nil
}

// recordingWriter counts and remembers every envelope the producer hands the
// writer, then forwards it to the inner writer.
type recordingWriter struct {
	inner evidence.ObjectWriter
	calls int
	last  store.Object
}

func (w *recordingWriter) PutObject(ctx context.Context, o store.Object) error {
	w.calls++
	w.last = o
	return w.inner.PutObject(ctx, o)
}
func (w *recordingWriter) BusyTimeout() time.Duration { return w.inner.BusyTimeout() }

func freshProofKey(t *testing.T) [32]byte {
	t.Helper()
	var k [32]byte
	if _, err := rand.Read(k[:]); err != nil {
		t.Fatal(err)
	}
	return k
}

func demoProofSource(t *testing.T) []byte {
	t.Helper()
	source, err := os.ReadFile(filepath.Join("..", "..", demoLogepochPath))
	if err != nil {
		t.Fatal(err)
	}
	return source
}

// putProofSource commits the production source through the REST writer, so the
// source is a carried object like every other demo object.
func putProofSource(t *testing.T, ctx context.Context, c *client, token, interp string, source []byte) hashref.HashRef {
	t.Helper()
	ref := hashref.SumSHA256(source)
	w := demoProofWriter{c: c, token: token, interpreter: interp}
	if err := w.PutObject(ctx, store.Object{Hash: ref, InterfaceHash: hashref.SumSHA256([]byte("world/source/ailang/v1")),
		SemanticID: "world/source/ailang/v1", Provenance: "demo-authored", Payload: source}); err != nil {
		t.Fatalf("commit source %s over REST: %v", ref, err)
	}
	return ref
}

// TestDemoProofHTTPProducerRoundTrip mints the proof of the production source
// with the real v0.41.0 checker, through the HTTP reader and writer, and checks
// the committed envelope, the production validator, and the REST inventory.
func TestDemoProofHTTPProducerRoundTrip(t *testing.T) {
	pin := requireProofPin(t)
	rig := newDemoProofRig(t, pin)
	ctx := boundedTestContext(t)
	c := rig.client(forwardRT{})
	h := proofHealth(t, ctx, c)
	source := demoProofSource(t)
	srcRef := putProofSource(t, ctx, c, rig.token, h.InterpreterRef, source)
	key := freshProofKey(t)
	rec := &recordingWriter{inner: demoProofWriter{c: c, token: rig.token, interpreter: h.InterpreterRef}}
	prevWorld, prevEntry := proofHeadEntry(t, ctx, c)

	ref, err := demoProve(ctx, demoProofReader{c: c}, rec, h, key, srcRef, demoProofRequired)
	if err != nil || ref.IsZero() || rec.calls != 1 {
		t.Fatalf("the verified source must mint once over REST: ref=%s err=%v calls=%d", ref, err, rec.calls)
	}
	if rec.last.Hash != ref || rec.last.InterfaceHash != evidence.InterfaceHashV1 || rec.last.SemanticID != evidence.ProofSemanticID {
		t.Fatalf("the writer was handed %s/%s/%q, not the envelope identity", rec.last.Hash, rec.last.InterfaceHash, rec.last.SemanticID)
	}
	meta, payload, err := demoProofReader{c: c}.ReadObject(ctx, ref, evidence.MaxBytes)
	if err != nil || !bytes.Equal(payload, rec.last.Payload) || hashref.SumSHA256(payload) != ref ||
		meta.InterfaceHash != evidence.InterfaceHashV1 || meta.SemanticID != evidence.ProofSemanticID {
		t.Fatalf("the committed envelope must read back byte-identical with its identity: meta=%+v err=%v", meta, err)
	}
	checker, err := demoProofChecker(h)
	if err != nil {
		t.Fatal(err)
	}
	v, err := evidence.NewValidator(key, demoProofReader{c: c}, evidence.CompilerConfig{
		Compiler: checker.Compiler, CompilerVersion: checker.CompilerVersion, ObjectReadTimeout: demoProofObjectBound,
	}, demoProofRequired)
	if err != nil {
		t.Fatal(err)
	}
	sealed, ok := v.ValidateProof(ctx, ref, srcRef).Validated()
	if !ok {
		t.Fatal("the production validator must accept the REST-committed envelope for its exact subject")
	}
	if grade, proven := v.Resolve(sealed).Proven(); !proven || grade != evidence.ResolvedGradeProven {
		t.Fatalf("the exact subject must resolve PROVEN: grade=%v proven=%v", grade, proven)
	}
	other := hashref.SumSHA256([]byte("a different subject"))
	res := v.ValidateProof(ctx, ref, other)
	if _, ok := res.Validated(); ok {
		t.Fatal("a different subject must not validate")
	}
	if reason, _ := res.Unsupported(); reason != evidence.UnsupportedSubjectMismatch {
		t.Fatalf("a different subject must be refused as subject_mismatch, got %q (err %v)", reason, res.Err())
	}
	// Carrying-entry inventory over REST: the envelope object exists (read above)
	// and a committed log entry carries it as its transition object.
	page, err := worldReader{ctx: ctx, c: c}.page(0, 500)
	if err != nil {
		t.Fatal(err)
	}
	var carrying *wireEntry
	for i := range page {
		e := page[i]
		if e.TransitionRef == ref.String() && e.Header.Interpreter == h.InterpreterRef {
			carrying = &e
		}
	}
	if carrying == nil {
		t.Fatalf("no REST log entry carries envelope %s: the object exists without a commit", ref)
	}
	// The committed entry must extend the head the proof was minted onto.
	if carrying.Header.PrevEntryHash != prevWorld.LogHead {
		t.Fatalf("the carrying entry's prevEntryHash %s is not the observed head's log head %s", carrying.Header.PrevEntryHash, prevWorld.LogHead)
	}
	if carrying.Header.SemanticsEpoch != prevEntry.Header.SemanticsEpoch {
		t.Fatalf("the carrying entry's semanticsEpoch %d is not the head entry's epoch %d", carrying.Header.SemanticsEpoch, prevEntry.Header.SemanticsEpoch)
	}
	if carrying.Header.EntryIndex != prevWorld.Revision+1 {
		t.Fatalf("the carrying entry index %d is not the observed revision %d plus one", carrying.Header.EntryIndex, prevWorld.Revision)
	}
}

// TestDemoProofSourceHashRefusal flips one byte of the source inside a valid
// HTTP 200 response. The producer's hash guard must refuse before any writer call.
func TestDemoProofSourceHashRefusal(t *testing.T) {
	pin := requireProofPin(t)
	rig := newDemoProofRig(t, pin)
	ctx := boundedTestContext(t)
	c := rig.client(forwardRT{})
	h := proofHealth(t, ctx, c)
	source := demoProofSource(t)
	srcRef := putProofSource(t, ctx, c, rig.token, h.InterpreterRef, source)
	flipped := rig.client(forwardRT{flip: srcRef.String()})
	rec := &recordingWriter{inner: demoProofWriter{c: c, token: rig.token, interpreter: h.InterpreterRef}}

	ref, err := demoProve(ctx, demoProofReader{c: flipped}, rec, h, freshProofKey(t), srcRef, demoProofRequired)
	if err == nil || !strings.Contains(err.Error(), "bytes hash to") {
		t.Fatalf("a flipped source byte must be refused by the hash guard: ref=%s err=%v", ref, err)
	}
	if !ref.IsZero() {
		t.Fatalf("a refused source must return the zero ref, got %s", ref)
	}
	if rec.calls != 0 {
		t.Fatalf("a refused source must never reach the envelope writer: calls=%d", rec.calls)
	}
}

// scriptRT answers every request from a function: the adapters' negative
// controls run against scripted REST replies, with no daemon.
type scriptRT func(*http.Request) (*http.Response, error)

func (f scriptRT) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func jsonReply(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}},
		Body: io.NopCloser(strings.NewReader(body)), ContentLength: int64(len(body))}
}

func scriptClient(rt http.RoundTripper) *client {
	return &client{base: "http://proof.invalid", timeout: demoProofHTTPTimeout, http: &http.Client{Transport: rt}}
}

// objectReply is a 200 object response carrying payload bytes.
func objectReply(t *testing.T, hash, iface string, payload []byte) *http.Response {
	t.Helper()
	raw, err := json.Marshal(map[string]any{"hash": hash, "interfaceHash": iface, "semanticId": "x", "provenance": "p", "payload": payload})
	if err != nil {
		t.Fatal(err)
	}
	return jsonReply(http.StatusOK, string(raw))
}

// TestDemoProofAdapterGuards names one negative control per guard in the HTTP
// adapters and the checker pin. Each control reaches its branch and asserts its
// specific refusal; each has a positive partner where a branch could be vacuous.
func TestDemoProofAdapterGuards(t *testing.T) {
	ctx := boundedTestContext(t)
	ref := hashref.SumSHA256([]byte("object under read"))
	other := hashref.SumSHA256([]byte("some other object"))
	iface := evidence.InterfaceHashV1.String()

	t.Run("reader-status-refused", func(t *testing.T) {
		c := scriptClient(scriptRT(func(*http.Request) (*http.Response, error) { return jsonReply(500, `{}`), nil }))
		_, payload, err := demoProofReader{c: c}.ReadObject(ctx, ref, 16)
		if err == nil || !strings.Contains(err.Error(), "returned HTTP 500") || payload != nil {
			t.Fatalf("a 500 object read must be refused by status: payload=%v err=%v", payload, err)
		}
	})
	t.Run("reader-absent-maps-to-nil", func(t *testing.T) {
		c := scriptClient(scriptRT(func(*http.Request) (*http.Response, error) { return jsonReply(404, `{}`), nil }))
		meta, payload, err := demoProofReader{c: c}.ReadObject(ctx, ref, 16)
		if err != nil || payload != nil || meta.SemanticID != "" {
			t.Fatalf("a 404 must read as absent (nil, nil): meta=%+v payload=%v err=%v", meta, payload, err)
		}
	})
	t.Run("reader-no-payload-refused", func(t *testing.T) {
		c := scriptClient(scriptRT(func(*http.Request) (*http.Response, error) {
			return jsonReply(http.StatusOK, fmt.Sprintf(`{"hash":%q,"interfaceHash":%q}`, ref.String(), iface)), nil
		}))
		var err error
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("a response without payload must be refused with %q; the reader panicked instead: %v", "carries no payload", r)
				}
			}()
			_, _, err = demoProofReader{c: c}.ReadObject(ctx, ref, 16)
		}()
		if err == nil || !strings.Contains(err.Error(), "carries no payload") {
			t.Fatalf("a response without payload must be refused: err=%v", err)
		}
	})
	t.Run("reader-malformed-body-refused", func(t *testing.T) {
		c := scriptClient(scriptRT(func(*http.Request) (*http.Response, error) {
			return jsonReply(http.StatusOK, "not json at all"), nil
		}))
		_, payload, err := demoProofReader{c: c}.ReadObject(ctx, ref, 16)
		if err == nil || !strings.Contains(err.Error(), "malformed body") || payload != nil {
			t.Fatalf("a 200 whose body is not JSON must be refused as malformed: payload=%v err=%v", payload, err)
		}
	})
	t.Run("reader-wrong-hash-refused", func(t *testing.T) {
		c := scriptClient(scriptRT(func(*http.Request) (*http.Response, error) {
			return objectReply(t, other.String(), iface, []byte("x")), nil
		}))
		_, _, err := demoProofReader{c: c}.ReadObject(ctx, ref, 16)
		if err == nil || !strings.Contains(err.Error(), "answered as") {
			t.Fatalf("an object answered under another hash must be refused: err=%v", err)
		}
	})
	t.Run("reader-oversize-refused", func(t *testing.T) {
		c := scriptClient(scriptRT(func(*http.Request) (*http.Response, error) {
			return objectReply(t, ref.String(), iface, []byte("12345")), nil
		}))
		_, payload, err := demoProofReader{c: c}.ReadObject(ctx, ref, 4)
		var tooLarge *store.ObjectTooLargeError
		if !errors.As(err, &tooLarge) || tooLarge.Size != 5 || tooLarge.MaxBytes != 4 || payload != nil {
			t.Fatalf("a payload past maxBytes must be refused as too large: payload=%v err=%v", payload, err)
		}
	})
	t.Run("reader-bad-interface-refused", func(t *testing.T) {
		c := scriptClient(scriptRT(func(*http.Request) (*http.Response, error) {
			return objectReply(t, ref.String(), "not-a-hash", []byte("x")), nil
		}))
		_, payload, err := demoProofReader{c: c}.ReadObject(ctx, ref, 16)
		if err == nil || !strings.Contains(err.Error(), "interface hash") || payload != nil {
			t.Fatalf("a malformed interface hash must be refused: payload=%v err=%v", payload, err)
		}
	})
	t.Run("busy-reader-zero", func(t *testing.T) {
		if d := (demoProofReader{}).BusyTimeout(); d != 0 {
			t.Fatalf("the REST reader has no lock-retry window: BusyTimeout=%s, want 0", d)
		}
	})
	t.Run("busy-writer-zero", func(t *testing.T) {
		if d := (demoProofWriter{}).BusyTimeout(); d != 0 {
			t.Fatalf("the REST writer has no lock-retry window: BusyTimeout=%s, want 0", d)
		}
	})
	t.Run("reader-valid-object-positive", func(t *testing.T) {
		c := scriptClient(scriptRT(func(*http.Request) (*http.Response, error) {
			return objectReply(t, ref.String(), iface, []byte("x")), nil
		}))
		meta, payload, err := demoProofReader{c: c}.ReadObject(ctx, ref, 16)
		if err != nil || string(payload) != "x" || meta.InterfaceHash != evidence.InterfaceHashV1 || meta.Provenance != "p" || meta.SemanticID != "x" {
			t.Fatalf("a well-formed object must read: meta=%+v payload=%q err=%v", meta, payload, err)
		}
	})

	// The writer scripts: a head, its world row and head entry, and one POST.
	head := hashref.SumSHA256([]byte("proof-head")).String()
	writerScript := func(t *testing.T, miss string, interp string, post func(demoCommit) (*http.Response, error), posts *int) *client {
		return scriptClient(scriptRT(func(req *http.Request) (*http.Response, error) {
			switch {
			case req.Method == http.MethodGet && req.URL.Path == "/v1/head":
				if miss == "head" {
					return jsonReply(http.StatusNotFound, `{}`), nil
				}
				if miss == "headerror" {
					return jsonReply(http.StatusInternalServerError, `{}`), nil
				}
				return jsonReply(http.StatusOK, head), nil
			case req.Method == http.MethodGet && strings.HasPrefix(req.URL.Path, "/v1/worlds/"):
				if miss == "world" {
					return jsonReply(http.StatusNotFound, `{}`), nil
				}
				return jsonReply(http.StatusOK, fmt.Sprintf(`{"ref":%q,"revision":3,"stateRoot":%q,"logHead":%q}`,
					head, hashref.SumSHA256([]byte("state")).String(), hashref.SumSHA256([]byte("log")).String())), nil
			case req.Method == http.MethodGet && req.URL.Path == "/v1/log/3":
				if miss == "entry" {
					return jsonReply(http.StatusNotFound, `{}`), nil
				}
				return jsonReply(http.StatusOK, fmt.Sprintf(`{"header":{"entryIndex":3,"semanticsEpoch":1,"transitionFn":%q,"interpreter":%q,"prevEntryHash":%q,"writtenBy":"t"},"entryHash":%q,"transitionRef":%q}`,
					hashref.SumSHA256([]byte("fn")).String(), interp, hashref.SumSHA256([]byte("prev")).String(),
					hashref.SumSHA256([]byte("log")).String(), hashref.SumSHA256([]byte("fn")).String())), nil
			case req.Method == http.MethodPost && req.URL.Path == "/v1/commit":
				*posts++
				var c demoCommit
				if err := json.NewDecoder(req.Body).Decode(&c); err != nil {
					t.Fatalf("the commit body must be the REST shape: %v", err)
				}
				return post(c)
			}
			t.Fatalf("unexpected REST call %s %s", req.Method, req.URL.Path)
			return nil, nil
		}))
	}
	envelope := store.Object{Hash: ref, InterfaceHash: evidence.InterfaceHashV1, SemanticID: evidence.ProofSemanticID, Provenance: "t", Payload: []byte("env")}
	pinned := hashref.SumSHA256([]byte("pinned interpreter")).String()
	okPost := func(c demoCommit) (*http.Response, error) {
		return jsonReply(http.StatusOK, fmt.Sprintf(`{"selectedHead":%q}`, c.NextWorld.Ref)), nil
	}

	t.Run("writer-no-head-refused", func(t *testing.T) {
		posts := 0
		err := demoProofWriter{c: writerScript(t, "head", pinned, okPost, &posts), interpreter: pinned}.PutObject(ctx, envelope)
		if err == nil || !strings.Contains(err.Error(), "no selected head") || posts != 0 {
			t.Fatalf("an absent head must refuse before any commit: posts=%d err=%v", posts, err)
		}
	})
	t.Run("writer-head-read-error-refused", func(t *testing.T) {
		posts := 0
		err := demoProofWriter{c: writerScript(t, "headerror", pinned, okPost, &posts), interpreter: pinned}.PutObject(ctx, envelope)
		if err == nil || !strings.Contains(err.Error(), "read head") || posts != 0 {
			t.Fatalf("a failed head read must refuse before any commit: posts=%d err=%v", posts, err)
		}
	})
	t.Run("writer-world-absent-refused", func(t *testing.T) {
		posts := 0
		err := demoProofWriter{c: writerScript(t, "world", pinned, okPost, &posts), interpreter: pinned}.PutObject(ctx, envelope)
		if err == nil || !strings.Contains(err.Error(), "has no world row") || posts != 0 {
			t.Fatalf("a head without its world row must refuse before any commit: posts=%d err=%v", posts, err)
		}
	})
	t.Run("writer-head-entry-absent-refused", func(t *testing.T) {
		posts := 0
		err := demoProofWriter{c: writerScript(t, "entry", pinned, okPost, &posts), interpreter: pinned}.PutObject(ctx, envelope)
		if err == nil || !strings.Contains(err.Error(), "is absent") || posts != 0 {
			t.Fatalf("a head entry that is absent must refuse before any commit: posts=%d err=%v", posts, err)
		}
	})
	t.Run("writer-foreign-interpreter-refused", func(t *testing.T) {
		posts := 0
		err := demoProofWriter{c: writerScript(t, "", pinned, okPost, &posts), interpreter: hashref.SumSHA256([]byte("other")).String()}.PutObject(ctx, envelope)
		if err == nil || !strings.Contains(err.Error(), "is not the pinned") || posts != 0 {
			t.Fatalf("a head entry under another interpreter must refuse before any commit: posts=%d err=%v", posts, err)
		}
	})
	t.Run("writer-transport-error-uncertain-no-retry", func(t *testing.T) {
		posts := 0
		err := demoProofWriter{c: writerScript(t, "", pinned, func(demoCommit) (*http.Response, error) {
			return nil, errors.New("connection reset")
		}, &posts), interpreter: pinned}.PutObject(ctx, envelope)
		if !store.IsUncertain(err) || posts != 1 {
			t.Fatalf("a commit whose outcome is unknown must be uncertain and never retried: posts=%d err=%v", posts, err)
		}
	})
	t.Run("writer-503-commit-uncertain", func(t *testing.T) {
		posts := 0
		err := demoProofWriter{c: writerScript(t, "", pinned, func(demoCommit) (*http.Response, error) {
			return jsonReply(http.StatusServiceUnavailable, `{"error":{"class":"CommitUncertain"}}`), nil
		}, &posts), interpreter: pinned}.PutObject(ctx, envelope)
		if !store.IsUncertain(err) || posts != 1 {
			t.Fatalf("a CommitUncertain answer must translate to store.UncertainError without retry: posts=%d err=%v", posts, err)
		}
	})
	t.Run("writer-503-other-class-refused", func(t *testing.T) {
		posts := 0
		err := demoProofWriter{c: writerScript(t, "", pinned, func(demoCommit) (*http.Response, error) {
			return jsonReply(http.StatusServiceUnavailable, `{"error":{"class":"Timeout"}}`), nil
		}, &posts), interpreter: pinned}.PutObject(ctx, envelope)
		if err == nil || store.IsUncertain(err) || !strings.Contains(err.Error(), "HTTP 503") || posts != 1 {
			t.Fatalf("a 503 of another class must be a plain refusal, not uncertain, and not retried: posts=%d err=%v", posts, err)
		}
	})
	t.Run("writer-409-refused-not-uncertain", func(t *testing.T) {
		posts := 0
		err := demoProofWriter{c: writerScript(t, "", pinned, func(demoCommit) (*http.Response, error) {
			return jsonReply(http.StatusConflict, `{"error":{"class":"Conflict"}}`), nil
		}, &posts), interpreter: pinned}.PutObject(ctx, envelope)
		if err == nil || store.IsUncertain(err) || !strings.Contains(err.Error(), "HTTP 409") || posts != 1 {
			t.Fatalf("a refused CAS must be a plain refusal, not uncertain, and not retried: posts=%d err=%v", posts, err)
		}
	})
	t.Run("writer-head-mismatch-uncertain", func(t *testing.T) {
		posts := 0
		err := demoProofWriter{c: writerScript(t, "", pinned, func(demoCommit) (*http.Response, error) {
			return jsonReply(http.StatusOK, fmt.Sprintf(`{"selectedHead":%q}`, hashref.SumSHA256([]byte("elsewhere")).String())), nil
		}, &posts), interpreter: pinned}.PutObject(ctx, envelope)
		if !store.IsUncertain(err) || posts != 1 {
			t.Fatalf("a 200 naming another head must be uncertain and not retried: posts=%d err=%v", posts, err)
		}
	})
	t.Run("writer-committed-positive", func(t *testing.T) {
		posts := 0
		err := demoProofWriter{c: writerScript(t, "", pinned, okPost, &posts), interpreter: pinned}.PutObject(ctx, envelope)
		if err != nil || posts != 1 {
			t.Fatalf("a commit answering its own head must succeed once: posts=%d err=%v", posts, err)
		}
	})

	pinBytes := []byte("pinned interpreter bytes")
	t.Run("checker-version-refused", func(t *testing.T) {
		_, err := demoProofChecker(daemon.HealthResponse{InterpreterVersion: "AILANG v0.99.0\n", InterpreterRef: hashref.SumSHA256(pinBytes).String()})
		if err == nil || !strings.Contains(err.Error(), "is not AILANG v0.41.0") {
			t.Fatalf("a non-pinned interpreter version must be refused: err=%v", err)
		}
	})
	t.Run("checker-ref-malformed-refused", func(t *testing.T) {
		_, err := demoProofChecker(daemon.HealthResponse{InterpreterVersion: "AILANG v0.41.0\n", InterpreterRef: "not-a-ref"})
		if err == nil || !strings.Contains(err.Error(), "interpreter ref") {
			t.Fatalf("a malformed interpreter ref must be refused: err=%v", err)
		}
	})
	t.Run("checker-archive-absent-refused", func(t *testing.T) {
		db := filepath.Join(t.TempDir(), "world.db")
		_, err := demoProofChecker(daemon.HealthResponse{DBPath: db, InterpreterVersion: "AILANG v0.41.0\n", InterpreterRef: hashref.SumSHA256(pinBytes).String()})
		if err == nil || !strings.Contains(err.Error(), "proof checker: archived interpreter:") {
			t.Fatalf("an absent archived interpreter must be refused: err=%v", err)
		}
	})
	t.Run("checker-readfile", func(t *testing.T) {
		db := filepath.Join(t.TempDir(), "world.db")
		ref := hashref.SumSHA256(pinBytes)
		dir := filepath.Join(db+".artifacts", "interpreters", ref.Algo(), ref.Digest())
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(dir, "ailang")
		if err := os.WriteFile(path, pinBytes, 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(path, 0); err != nil {
			t.Fatal(err)
		}
		defer os.Chmod(path, 0o600)
		_, err := demoProofChecker(daemon.HealthResponse{DBPath: db, InterpreterVersion: demoProofVersion, InterpreterRef: ref.String()})
		if err == nil || !strings.Contains(err.Error(), "proof checker: read archived interpreter:") {
			t.Fatalf("checker-readfile must name the archived interpreter read failure: %v", err)
		}
	})
	t.Run("checker-archive-bytes-refused", func(t *testing.T) {
		db := filepath.Join(t.TempDir(), "world.db")
		want := hashref.SumSHA256(pinBytes)
		dir := filepath.Join(db+".artifacts", "interpreters", want.Algo(), want.Digest())
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "ailang"), []byte("tampered bytes"), 0o755); err != nil {
			t.Fatal(err)
		}
		_, err := demoProofChecker(daemon.HealthResponse{DBPath: db, InterpreterVersion: "AILANG v0.41.0\n", InterpreterRef: want.String()})
		if err == nil || !strings.Contains(err.Error(), "archived pin bytes hash to") {
			t.Fatalf("archived bytes that miss the health ref must be refused: err=%v", err)
		}
	})
}

// proofHeadEntry reads the selected head, its world row and its head entry over REST.
func proofHeadEntry(t *testing.T, ctx context.Context, c *client) (wireWorld, wireEntry) {
	t.Helper()
	ref, has, err := headState(ctx, c)
	if err != nil || !has {
		t.Fatalf("read head: has=%v err=%v", has, err)
	}
	r := worldReader{ctx: ctx, c: c}
	world, found, err := r.world(ref)
	if err != nil || !found {
		t.Fatalf("read head world %s: found=%v err=%v", ref, found, err)
	}
	entry, found, err := r.entry(world.Revision)
	if err != nil || !found {
		t.Fatalf("read head entry %d: found=%v err=%v", world.Revision, found, err)
	}
	return world, entry
}

// proofHealth reads the daemon's health through the REST client.
func proofHealth(t *testing.T, ctx context.Context, c *client) daemon.HealthResponse {
	t.Helper()
	var h daemon.HealthResponse
	if err := demoGet(worldReader{ctx: ctx, c: c}, "/v1/health", &h); err != nil {
		t.Fatal(err)
	}
	return h
}

// A standalone no-contract module is intentional here: the negative control
// must exercise a real empty verified set, rather than a synthetic JSON report.
func demoNoContractSource(t *testing.T, pin string) []byte {
	t.Helper()
	path := filepath.Join(t.TempDir(), "plain.ail")
	source := []byte("module plain\nexport func identity(x: int) -> int { x }\n")
	if err := os.WriteFile(path, source, 0o600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(boundedTestContext(t), 10*time.Second)
	defer cancel()
	if out, err := exec.CommandContext(ctx, pin, "check", path).CombinedOutput(); err != nil {
		t.Fatalf("standalone no-contract control must check under the real pin: %v %s", err, out)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func proofValidator(t *testing.T, c *client, h daemon.HealthResponse, key [32]byte) *evidence.Validator {
	t.Helper()
	checker, err := demoProofChecker(h)
	if err != nil {
		t.Fatal(err)
	}
	v, err := evidence.NewValidator(key, demoProofReader{c: c}, evidence.CompilerConfig{
		Compiler: checker.Compiler, CompilerVersion: checker.CompilerVersion, ObjectReadTimeout: demoProofObjectBound,
	}, demoProofRequired)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func requireDemoProven(t *testing.T, ctx context.Context, v *evidence.Validator, ref, subject hashref.HashRef) {
	t.Helper()
	sealed, valid := v.ValidateProof(ctx, ref, subject).Validated()
	grade, proven := v.Resolve(sealed).Proven()
	if !valid || !proven || grade != evidence.ResolvedGradeProven {
		t.Fatalf("same-key positive arm must Resolve PROVEN for the exact subject: valid=%v proven=%v grade=%v", valid, proven, grade)
	}
}

func TestDemoProofWrongKeyCannotProve(t *testing.T) {
	pin := requireProofPin(t)
	rig := newDemoProofRig(t, pin)
	ctx := boundedTestContext(t)
	c := rig.client(forwardRT{})
	h := proofHealth(t, ctx, c)
	subject := putProofSource(t, ctx, c, rig.token, h.InterpreterRef, demoProofSource(t))
	key, different := freshProofKey(t), freshProofKey(t)
	if key == different {
		t.Fatal("independent random keys unexpectedly equal")
	}
	rec := &recordingWriter{inner: demoProofWriter{c: c, token: rig.token, interpreter: h.InterpreterRef}}
	ref, err := demoProve(ctx, demoProofReader{c: c}, rec, h, key, subject, demoProofRequired)
	if err != nil || ref.IsZero() || rec.calls != 1 {
		t.Fatalf("real HTTP mint must succeed once: ref=%s calls=%d err=%v", ref, rec.calls, err)
	}
	_, payload, err := (demoProofReader{c: c}).ReadObject(ctx, ref, evidence.MaxBytes)
	if err != nil || !bytes.Equal(payload, rec.last.Payload) {
		t.Fatalf("mint must be committed over REST: %v", err)
	}
	wrong := proofValidator(t, c, h, different)
	result := wrong.ValidateProof(ctx, ref, subject)
	reason, unsupported := result.Unsupported()
	sealed, valid := result.Validated()
	_, proven := wrong.Resolve(sealed).Proven()
	if !unsupported || reason != evidence.UnsupportedUnauthenticatedReport || valid || proven {
		t.Fatalf("wrong key must remain UNAVAILABLE, never PROVEN: reason=%s valid=%v proven=%v", reason, valid, proven)
	}
	same := proofValidator(t, c, h, key) // M3b-key: replace key with different here.
	requireDemoProven(t, ctx, same, ref, subject)
}

func TestDemoProofNoVerifiedSetRefuses(t *testing.T) {
	pin := requireProofPin(t)
	rig := newDemoProofRig(t, pin)
	ctx := boundedTestContext(t)
	c := rig.client(forwardRT{})
	h := proofHealth(t, ctx, c)
	plain := demoNoContractSource(t, pin)
	bad := putProofSource(t, ctx, c, rig.token, h.InterpreterRef, plain)
	before, _ := proofHeadEntry(t, ctx, c)
	rec := &recordingWriter{inner: demoProofWriter{c: c, token: rig.token, interpreter: h.InterpreterRef}}
	key := freshProofKey(t)
	ref, err := demoProve(ctx, demoProofReader{c: c}, rec, h, key, bad, demoProofRequired)
	if !errors.Is(err, evidence.ErrVacuousCheck) || !ref.IsZero() || rec.calls != 0 {
		t.Fatalf("no verified set must refuse with ErrVacuousCheck, zero ref and zero writes: ref=%s calls=%d err=%v", ref, rec.calls, err)
	}
	after, _ := proofHeadEntry(t, ctx, c)
	if after != before {
		t.Fatal("vacuous proof changed committed REST history")
	}
	source := demoProofSource(t) // M3b-null-verified: replace with plain (consistent hash).
	good := putProofSource(t, ctx, c, rig.token, h.InterpreterRef, source)
	ref, err = demoProve(ctx, demoProofReader{c: c}, rec, h, key, good, demoProofRequired)
	if err != nil || ref.IsZero() || rec.calls != 1 {
		t.Fatalf("positive proof-control must mint required identities sameRef and servesEntry: ref=%s calls=%d err=%v", ref, rec.calls, err)
	}
	requireDemoProven(t, ctx, proofValidator(t, c, h, key), ref, good)
}

// uncertainProofRT counts actual POST attempts. The lost-response arm waits for
// the real commit to succeed before discarding the acknowledgement.
type uncertainProofRT struct {
	mode   string
	calls  int
	object store.Object
}

func (rt *uncertainProofRT) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method != http.MethodPost || req.URL.Path != "/v1/commit" {
		return http.DefaultTransport.RoundTrip(req)
	}
	rt.calls++
	raw, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewReader(raw))
	var commit demoCommit
	if err := json.Unmarshal(raw, &commit); err != nil {
		return nil, err
	}
	o := commit.Objects[0]
	rt.object = store.Object{Hash: hashref.MustParse(o.Hash), Payload: o.Payload}
	if rt.mode == "safe-refusal" {
		return jsonReply(http.StatusForbidden, `{"error":{"class":"Denied"}}`), nil
	}
	if rt.mode == "uncertain-before-send" {
		return nil, errors.New("test transport lost before send")
	}
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("real commit did not succeed: HTTP %d %v", resp.StatusCode, err)
	}
	return nil, errors.New("test transport lost acknowledgement after real commit")
}

func TestDemoProofUncertainCommitStops(t *testing.T) {
	pin := requireProofPin(t)
	for _, mode := range []string{"safe-refusal", "uncertain-before-send", "lost-response-after-commit"} {
		t.Run(mode, func(t *testing.T) {
			rig := newDemoProofRig(t, pin)
			ctx := boundedTestContext(t)
			c := rig.client(forwardRT{})
			h := proofHealth(t, ctx, c)
			subject := putProofSource(t, ctx, c, rig.token, h.InterpreterRef, demoProofSource(t))
			key := freshProofKey(t)
			keyPath := filepath.Join(t.TempDir(), "operator.key")
			if err := os.WriteFile(keyPath, key[:], 0o600); err != nil {
				t.Fatal(err)
			}
			before, _ := proofHeadEntry(t, ctx, c)
			rt := &uncertainProofRT{mode: mode}
			ref, err := demoSeedProof(ctx, rig.client(rt), rig.token, h, keyPath, subject)
			if rt.calls != 1 {
				t.Fatalf("write attempt count must equal one: got %d", rt.calls)
			}
			if err == nil || !ref.IsZero() || !strings.Contains(err.Error(), "S4 refused:") || !strings.Contains(err.Error(), "inspect history") || !strings.Contains(err.Error(), "no success manifest written") {
				t.Fatalf("S4 refusal must stop without success or reseed: ref=%s err=%v", ref, err)
			}
			if store.IsUncertain(err) != (mode != "safe-refusal") {
				t.Fatalf("uncertain outcome must survive demo error wrapping: %v", err)
			}
			assertDemoProofKeyPrivate(t, map[string]string{"returned-error": err.Error()}, key[:])
			after, _ := proofHeadEntry(t, ctx, c)
			_, payload, readErr := (demoProofReader{c: c}).ReadObject(ctx, rt.object.Hash, evidence.MaxBytes)
			if readErr != nil {
				t.Fatal(readErr)
			}
			if mode == "lost-response-after-commit" {
				if after.Revision != before.Revision+1 || after.LogHead == before.LogHead || !bytes.Equal(payload, rt.object.Payload) {
					t.Fatal("lost acknowledgement must leave exactly one real committed envelope readable over REST")
				}
			} else if after != before || payload != nil {
				t.Fatal("uncommitted refusal changed REST history or stored an envelope")
			}
		})
	}
}

type blockingProofReader struct{ seen chan context.Context }

func (*blockingProofReader) BusyTimeout() time.Duration { return 0 }
func (r *blockingProofReader) ReadObject(ctx context.Context, _ hashref.HashRef, _ int64) (store.ObjectMeta, []byte, error) {
	r.seen <- ctx
	<-ctx.Done()
	return store.ObjectMeta{}, nil, ctx.Err()
}

func TestDemoProofBounds(t *testing.T) {
	c := demoProofClient("http://proof.invalid")
	reader, writer := demoProofReader{c: c}, demoProofWriter{c: c}
	if reader.c.timeout != demoProofHTTPTimeout || writer.c.http.Timeout != demoProofHTTPTimeout || demoProofHTTPTimeout >= demoProofObjectBound {
		t.Fatal("production client/adapters must carry the 10s HTTP bound below the 12s object bound")
	}
	whole, cancelWhole := demoProofContext(context.Background())
	defer cancelWhole()
	if deadline, ok := whole.Deadline(); !ok || time.Until(deadline) > demoProofWholeBound || time.Until(deadline) < demoProofWholeBound-time.Second {
		t.Fatal("whole proof context must carry the 90s deadline")
	}
	pin := requireProofPin(t)
	rig := newDemoProofRig(t, pin)
	h := proofHealth(t, boundedTestContext(t), rig.client(forwardRT{}))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	readerBlock := &blockingProofReader{seen: make(chan context.Context, 1)}
	rec := &recordingWriter{inner: demoProofWriter{}}
	done := make(chan error, 1)
	key := freshProofKey(t)
	go func() {
		_, err := demoProve(ctx, readerBlock, rec, h, key, hashref.SumSHA256([]byte("source")), demoProofRequired)
		done <- err
	}()
	select {
	case passed := <-readerBlock.seen:
		if d, ok := passed.Deadline(); !ok || time.Until(d) > demoProofWholeBound {
			t.Fatal("demoProve reader context must have a deadline no later than the whole bound")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("producer never reached blocking reader")
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) || rec.calls != 0 {
			t.Fatalf("cancelling demoProve must stop producer with context error before write: calls=%d err=%v", rec.calls, err)
		}
	case <-time.After(time.Second):
		t.Fatal("cancelled producer did not stop")
	}
}

// keyReadFault returns usable metadata/data alongside an error so removing the
// corresponding guard reaches its named assertion, rather than panicking.
type keyReadFault struct {
	info             os.FileInfo
	statErr, readErr error
	data             []byte
}

func (f keyReadFault) Stat() (os.FileInfo, error) { return f.info, f.statErr }
func (f keyReadFault) Read(p []byte) (int, error) { return copy(p, f.data), f.readErr }

func TestDemoProofKeyIORefusals(t *testing.T) {
	key := freshProofKey(t)
	path := filepath.Join(t.TempDir(), "proof.key")
	if err := os.WriteFile(path, key[:], 0o600); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, arm := range []string{"stat", "read"} {
		t.Run(arm, func(t *testing.T) {
			f := keyReadFault{info: info, data: key[:], readErr: io.EOF}
			if arm == "stat" {
				f.statErr = errors.New("stat control")
			} else {
				f.readErr = errors.New("read control")
			}
			got, err := demoReadProofKey(f)
			if err == nil || !strings.Contains(err.Error(), "proof key: "+arm+":") || got != [32]byte{} {
				t.Fatalf("key %s error must refuse without bytes: err=%v zero=%v", arm, err, got == [32]byte{})
			}
		})
	}
}

// assertDemoProofKeyPrivate checks each surface separately, including encodings
// of both the operator key and malformed key-file contents.
func assertDemoProofKeyPrivate(t *testing.T, surfaces map[string]string, keys ...[]byte) {
	t.Helper()
	for _, key := range keys {
		if len(key) == 0 {
			continue
		}
		forms := map[string]string{
			"raw":            string(key),
			"hex-lower":      hex.EncodeToString(key),
			"hex-upper":      strings.ToUpper(hex.EncodeToString(key)),
			"base64-std":     base64.StdEncoding.EncodeToString(key),
			"base64-raw-url": base64.RawURLEncoding.EncodeToString(key),
		}
		for surface, text := range surfaces {
			for encoding, secret := range forms {
				if strings.Contains(text, secret) {
					t.Fatalf("seed-key-disclosure: %s contains %s key", surface, encoding)
				}
			}
		}
	}
}

func TestDemoProofKeyDisclosure(t *testing.T) {
	key := freshProofKey(t)
	for _, size := range []int{31, 33} {
		t.Run(fmt.Sprint(size), func(t *testing.T) {
			raw := append(append([]byte(nil), key[:]...), key[0])[:size]
			path := filepath.Join(t.TempDir(), "operator.key")
			if err := os.WriteFile(path, raw, 0o600); err != nil {
				t.Fatal(err)
			}
			// Invalid lengths must stop before a reader/writer can use this nil client.
			_, err := demoSeedProof(context.Background(), nil, "", daemon.HealthResponse{}, path, hashref.HashRef{})
			if err == nil {
				t.Fatal("seed-key-disclosure: malformed key must return an error")
			}
			assertDemoProofKeyPrivate(t, map[string]string{"returned-error": err.Error()}, key[:], raw)
		})
	}
}
