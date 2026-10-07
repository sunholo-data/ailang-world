package projection

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/canon"
	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// Row 123 (w-a2a-invocation-wire-coverage): the refusal mapping is proven on
// the WIRE, not only at dispatchError. Every failure below is produced by the
// REAL coordinator behind the REAL /a2a/ handler; only the coordinator's own
// construction seams (Store, Runner, BinderFor) inject the failure. The
// assertions read the recorded HTTP response, so a refactor of the
// projection.go error write (id, envelope, message, Connection header) fails
// here even when the dispatchError table still passes.

// wireRig is a published, invocable echo transition on a genesis world.
type wireRig struct {
	st    *store.Store
	tok   string
	world store.World
}

func newWireRig(t *testing.T) *wireRig { return newWireRigWith(t, false) }

// newWireRigWith is newWireRig; effectful keeps the descriptor's declared
// effect (alpha@world) so the coordinator takes the plan/finish path.
func newWireRigWith(t *testing.T, effectful bool) *wireRig {
	t.Helper()
	st := openStore(t)
	genesis := store.Object{Hash: hashref.SumSHA256([]byte("genesis-state")), InterfaceHash: hashref.SumSHA256([]byte("test/genesis")),
		SemanticID: "test/genesis", Provenance: "projection-wire-test", Payload: []byte("genesis-state")}
	entryHash := hashref.SumSHA256([]byte("genesis-entry"))
	world := store.World{Ref: hashref.SumSHA256([]byte("genesis-world")), Revision: 0, StateRoot: genesis.Hash, LogHead: entryHash}
	interp := hashref.SumSHA256([]byte("interpreter"))
	if err := st.Commit(boundedTestContext(t), store.Commit{Objects: []store.Object{genesis}, NextWorld: world,
		Entry: store.LogEntry{Header: store.LogHeader{EntryIndex: 0, SemanticsEpoch: 1, TransitionFn: hashref.SumSHA256([]byte("genesis-fn")),
			Interpreter: interp, PrevEntryHash: hashref.SumSHA256([]byte("genesis-prev")), WrittenBy: "projection-wire-test"},
			EntryHash: entryHash, TransitionRef: genesis.Hash}}); err != nil {
		t.Fatal(err)
	}
	source, err := canon.Source([]byte("module transitions/echo\nexport func main(input: string) -> string { input }\n"))
	if err != nil {
		t.Fatal(err)
	}
	src := store.Object{Hash: hashref.SumSHA256(source), InterfaceHash: hashref.SumSHA256([]byte("world/transition-source/v1")),
		SemanticID: "world/transition-source/v1", Provenance: "projection-wire-test", Payload: source}
	if err := st.PutObject(boundedTestContext(t), src); err != nil {
		t.Fatal(err)
	}
	d := descriptor("tools.echo", "alpha")
	d.TransitionFn, d.Interpreter = src.Hash, interp
	if !effectful {
		d.DeclaredEffects = nil
	}
	seedRegistry(t, st, d)
	return &wireRig{st: st, tok: mintToken(t, st, "ep-a", []broker.Capability{liveGrant("alpha")}), world: world}
}

// handler wires a REAL coordinator over the given seams into a REAL handler.
// A nil cs/binder means the production store / broker.OpenBinder.
func (r *wireRig) handler(t *testing.T, cs coordinator.Store, runner coordinator.Runner, binder coordinator.BinderFor, invokeWait time.Duration) *Handler {
	t.Helper()
	if cs == nil {
		cs = r.st
	}
	if binder == nil {
		binder = func(ep string, grants []broker.Capability) transitionreg.Binder {
			return broker.OpenBinder(r.st, ep, grants, nil)
		}
	}
	coord, err := coordinator.New(coordinator.Config{Store: cs, Runner: runner, Binder: binder,
		Now: func() int64 { return time.Now().Unix() }, MaxInput: 1 << 20, MaxOutput: 1 << 20})
	if err != nil {
		t.Fatal(err)
	}
	cfg := testConfig(r.st)
	cfg.Coordinator = coord
	cfg.InvokeWait = invokeWait
	return mustHandler(t, cfg)
}

func wireBody(reqID, taskID string) string {
	return fmt.Sprintf(`{"jsonrpc":"2.0","id":%q,"method":"tasks/send","params":{"id":%q,"metadata":{"skill_id":"tools.echo"},"message":{"parts":[{"type":"data","data":{"x":1}}]}}}`, reqID, taskID)
}

// assertWireRefusal pins one refusal's bytes on the recorder: HTTP 200 JSON,
// the JSON-RPC id echoed verbatim, the exact constant code and message, no
// inner error detail, and the Connection header (close only when the
// invocation deadline ended the request).
func assertWireRefusal(t *testing.T, rec *recorded, reqID string, wantCode int, wantMsg string, wantClose bool, hidden ...string) {
	t.Helper()
	if rec.code != http.StatusOK {
		t.Fatalf("status=%d, want 200 (A2AError always writes 200); body=%s", rec.code, rec.body)
	}
	if ct := rec.header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type=%q, want application/json; body=%s", ct, rec.body)
	}
	code, msg, id := a2aErr(t, rec.body)
	if want := fmt.Sprintf("%q", reqID); string(id) != want {
		t.Fatalf("JSON-RPC id=%s, want the request id %s echoed; body=%s", id, want, rec.body)
	}
	for _, h := range append(hidden, "coordinator:", "store:", "transition registry:", "context deadline exceeded") {
		if bytes.Contains(rec.body, []byte(h)) {
			t.Fatalf("body interpolates inner error detail %q: %s", h, rec.body)
		}
	}
	if code != wantCode || msg != wantMsg {
		t.Fatalf("error=%d %q, want %d %q; body=%s", code, msg, wantCode, wantMsg, rec.body)
	}
	got := rec.header.Get("Connection")
	if wantClose && got != "close" {
		t.Fatalf("Connection=%q, want close: a timed-out invocation must release its transport", got)
	}
	if !wantClose && got != "" {
		t.Fatalf("Connection=%q, want absent: only an expired invocation deadline closes the transport", got)
	}
}

type recorded struct {
	code   int
	header http.Header
	body   []byte
}

func post(t *testing.T, h *Handler, auth, body string) *recorded {
	t.Helper()
	rec := postA2A(t, h, auth, body)
	return &recorded{code: rec.Code, header: rec.Header(), body: rec.Body.Bytes()}
}

// failStore injects one failure into the coordinator's store seam and counts
// that the coordinator reached it (so a refusal that happened earlier — e.g.
// at admission, which shares some messages — cannot pass as this row).
type failStore struct {
	*store.Store
	getObject error
	commit    error
	reached   *atomic.Int32
}

func (s failStore) GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error) {
	if s.getObject != nil {
		s.reached.Add(1)
		return store.Object{}, false, s.getObject
	}
	return s.Store.GetObject(ctx, ref)
}

func (s failStore) Commit(ctx context.Context, c store.Commit) error {
	if s.commit != nil {
		s.reached.Add(1)
		return s.commit
	}
	return s.Store.Commit(ctx, c)
}

type denyingBinder struct {
	label   string
	reached *atomic.Int32
}

func (b denyingBinder) Bind(broker.Manifest) (*broker.BoundInvoker, error) {
	b.reached.Add(1)
	return nil, &transitionreg.AccessDeniedError{ID: "tools.echo", Label: b.label}
}

// hookRunner runs hook (with the invocation ctx) and then returns stdout/err.
type hookRunner struct {
	hook   func(ctx context.Context) error
	stdout string
}

func (r hookRunner) RunContext(ctx context.Context, _ capsule.Entry) (capsule.Result, error) {
	if err := r.hook(ctx); err != nil {
		return capsule.Result{}, err
	}
	return capsule.Result{Stdout: []byte(r.stdout)}, nil
}

// TestA2ADispatchWire is AC-WIRE: representative typed refusals (R2/R3, R11,
// R14, R16 and the R4/R6 default) produced by the real coordinator and read
// back from the HTTP response bytes.
func TestA2ADispatchWire(t *testing.T) {
	ok := func(context.Context) error { return nil }

	t.Run("R13_integrity", func(t *testing.T) {
		r := newWireRig(t)
		h := r.handler(t, nil, hookRunner{stdout: `{}`, hook: ok}, nil, 2*time.Second)
		first := post(t, h, "Bearer "+r.tok, wireBody("first", "t1"))
		if bytes.Contains(first.body, []byte(`"error"`)) {
			t.Fatalf("commit: %s", first.body)
		}
		hits := 0
		h = r.handler(t, absentRecordStore{Store: r.st, hits: &hits}, unusedRunner{}, nil, 2*time.Second)
		var log bytes.Buffer
		h.errorLog = &log
		rec := post(t, h, "Bearer "+r.tok, wireBody("retry", "t1"))
		assertWireRefusal(t, rec, "retry", codeInternal, notAvailableMessage, false, "coordinator:", "a2a:ep-a:t1")
		if hits != 1 {
			t.Fatalf("fixture hits=%d", hits)
		}
		assertRefusalLog(t, &log, "a2a:ep-a:t1", "is absent")
	})
	t.Run("store_error_newline", func(t *testing.T) {
		r := newWireRig(t)
		var reached atomic.Int32
		secret := "SECRET\nsecond line"
		h := r.handler(t, failStore{Store: r.st, reached: &reached, getObject: errors.New(secret)}, unusedRunner{}, nil, 2*time.Second)
		var log bytes.Buffer
		h.errorLog = &log
		rec := post(t, h, "Bearer "+r.tok, wireBody("newline", "t1"))
		assertWireRefusal(t, rec, "newline", codeInternal, notAvailableMessage, false, "SECRET")
		if reached.Load() != 1 {
			t.Fatal("fixture did not fire")
		}
		assertRefusalLog(t, &log, "a2a:ep-a:t1", `SECRET\nsecond line`)
	})
	t.Run("R15", func(t *testing.T) {
		r := newWireRig(t)
		h := r.handler(t, intentOnlyStore{Store: r.st}, unusedRunner{}, nil, 2*time.Second)
		var log bytes.Buffer
		h.errorLog = &log
		rec := post(t, h, "Bearer "+r.tok, wireBody("spent", "t1"))
		assertWireRefusal(t, rec, "spent", codeInternal, "invocation was not committed; send a new task id", false)
		if log.Len() != 0 {
			t.Fatalf("R15 logged: %q", log.String())
		}
	})

	t.Run("R2_R3", func(t *testing.T) {
		r := newWireRig(t)
		var reached atomic.Int32
		const secret = "SECRET-R2R3-broker-label"
		binder := func(string, []broker.Capability) transitionreg.Binder {
			return denyingBinder{label: secret, reached: &reached}
		}
		h := r.handler(t, nil, unusedRunner{}, binder, 2*time.Second)
		rec := post(t, h, "Bearer "+r.tok, wireBody("wire-R2R3", "t-r2r3"))
		if reached.Load() != 1 {
			t.Fatalf("binder reached %d times, want 1: the refusal must come from dispatch, not admission", reached.Load())
		}
		assertWireRefusal(t, rec, "wire-R2R3", codeInvalidParams, msgNotAuthorized, false, secret)
	})

	// Exercise errors.Is deadline classification without waiting for a wall clock.
	t.Run("R11", func(t *testing.T) {
		r := newWireRig(t)
		var reached atomic.Int32
		const secret = "SECRET-R11-store-detail"
		h := r.handler(t, failStore{Store: r.st, reached: &reached, getObject: fmt.Errorf("%s: %w", secret, context.DeadlineExceeded)}, unusedRunner{}, nil, 20*time.Second)
		var log bytes.Buffer
		h.errorLog = &log
		rec := post(t, h, "Bearer "+r.tok, wireBody("wire-R11", "t-r11"))
		if reached.Load() != 1 {
			t.Fatalf("store reached %d times, want 1", reached.Load())
		}
		assertWireRefusal(t, rec, "wire-R11", codeInternal, "invocation exceeded its deadline", false, secret)
		if log.Len() != 0 {
			t.Fatalf("deadline logged: %q", log.String())
		}
	})

	t.Run("R11_transport", func(t *testing.T) {
		r := newWireRig(t)
		var reached atomic.Int32
		const secret = "SECRET-R11-capsule-detail"
		block := hookRunner{hook: func(ctx context.Context) error {
			reached.Add(1)
			select {
			case <-ctx.Done():
				return fmt.Errorf("%s: %w", secret, ctx.Err())
			case <-time.After(5 * time.Second): // a ctx-dropping mutant fails, never hangs
				return errors.New(secret + ": the invocation deadline never cancelled the capsule")
			}
		}}
		// The deadline also covers admission (SQLite), so it must be wide enough
		// that a loaded -race runner still reaches the runner before it expires.
		h := r.handler(t, nil, block, nil, 2*time.Second)
		var log bytes.Buffer
		h.errorLog = &log
		rec := post(t, h, "Bearer "+r.tok, wireBody("wire-R11", "t-r11"))
		if reached.Load() != 1 {
			t.Fatalf("runner reached %d times, want 1", reached.Load())
		}
		assertWireRefusal(t, rec, "wire-R11", codeInternal, "invocation exceeded its deadline", true, secret)
		if log.Len() != 0 {
			t.Fatalf("deadline logged: %q", log.String())
		}
	})

	t.Run("R14", func(t *testing.T) {
		// A REAL store conflict: the world head moves while the capsule runs,
		// so the store's compare-and-append refuses the planned commit.
		r := newWireRig(t)
		moved := store.World{Ref: hashref.SumSHA256([]byte("moved-world")), Revision: 1,
			StateRoot: r.world.StateRoot, LogHead: hashref.SumSHA256([]byte("moved-entry"))}
		moveHead := hookRunner{stdout: `{"ok":true}`, hook: func(ctx context.Context) error {
			if err := r.st.PutWorld(ctx, moved); err != nil {
				return err
			}
			return r.st.SelectHead(ctx, moved.Ref)
		}}
		h := r.handler(t, nil, moveHead, nil, 2*time.Second)
		rec := post(t, h, "Bearer "+r.tok, wireBody("wire-R14", "t-r14"))
		head, _, err := r.st.SelectedHead(boundedTestContext(t))
		if err != nil || head != moved.Ref {
			t.Fatalf("control: head=%v err=%v, want the moved world %v", head, err, moved.Ref)
		}
		assertWireRefusal(t, rec, "wire-R14", codeInternal, "world head moved during invocation; not committed; send a new task id", false,
			moved.Ref.String(), r.world.Ref.String(), "stale observed")
	})

	t.Run("R16", func(t *testing.T) {
		r := newWireRig(t)
		var reached atomic.Int32
		const secret = "SECRET-R16-driver-report"
		st := failStore{Store: r.st, reached: &reached, commit: &store.UncertainError{Op: "Commit", Cause: errors.New(secret)}}
		h := r.handler(t, st, hookRunner{stdout: `{"ok":true}`, hook: ok}, nil, 2*time.Second)
		rec := post(t, h, "Bearer "+r.tok, wireBody("wire-R16", "t-r16"))
		if reached.Load() != 1 {
			t.Fatalf("Commit reached %d times, want 1", reached.Load())
		}
		assertWireRefusal(t, rec, "wire-R16", codeInternal, "invocation outcome is not confirmed; resend the same task id", false, secret)
	})

	t.Run("R4_R6_default", func(t *testing.T) {
		r := newWireRig(t)
		var reached atomic.Int32
		const secret = "SECRET-R4R6-disk-io"
		st := failStore{Store: r.st, reached: &reached, getObject: errors.New(secret)}
		h := r.handler(t, st, unusedRunner{}, nil, 2*time.Second)
		var log bytes.Buffer
		h.errorLog = &log
		rec := post(t, h, "Bearer "+r.tok, wireBody("wire-R4R6", "t-r4r6"))
		if reached.Load() != 1 {
			t.Fatalf("GetObject reached %d times, want 1: the refusal must come from dispatch, not the nil-coordinator branch", reached.Load())
		}
		assertWireRefusal(t, rec, "wire-R4R6", codeInternal, notAvailableMessage, false, secret)
		assertRefusalLog(t, &log, "a2a:ep-a:t-r4r6", secret)
	})
}

// stallStore holds one durable operation until the invocation ctx ends — the
// held-sole-connection shape at the coordinator's store seam. release frees a
// stall whose ctx never ends (a ctx-dropping mutant) so the test fails instead
// of hanging.
type stallStore struct {
	*store.Store
	op        string
	uncertain bool
	armed     *atomic.Bool
	stalls    *atomic.Int32
	release   chan struct{}
}

func (s stallStore) stall(ctx context.Context, op string) (bool, error) {
	if !s.armed.Load() || s.op != op {
		return false, nil
	}
	s.stalls.Add(1)
	select {
	case <-ctx.Done():
	case <-s.release:
		return true, errors.New("stall released by the test bound")
	}
	if s.uncertain {
		return true, &store.UncertainError{Op: op, Cause: ctx.Err()}
	}
	return true, ctx.Err()
}

func (s stallStore) GetReceipt(ctx context.Context, id string) (store.Receipt, bool, error) {
	if stalled, err := s.stall(ctx, "GetReceipt"); stalled {
		return store.Receipt{}, false, err
	}
	return s.Store.GetReceipt(ctx, id)
}

func (s stallStore) AppendIntent(ctx context.Context, id string, intent store.JournalIntent) (int64, hashref.HashRef, error) {
	if stalled, err := s.stall(ctx, "AppendIntent"); stalled {
		return 0, hashref.HashRef{}, err
	}
	return s.Store.AppendIntent(ctx, id, intent)
}

func (s stallStore) Commit(ctx context.Context, c store.Commit) error {
	if stalled, err := s.stall(ctx, "Commit"); stalled {
		return err
	}
	return s.Store.Commit(ctx, c)
}

func settledGoroutineCount(base int) int {
	deadline := time.Now().Add(2 * time.Second)
	n := runtime.NumGoroutine()
	for n > base && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
		n = runtime.NumGoroutine()
	}
	return n
}

// TestA2ADurableDeadlineWire is AC-REPEAT: the round-2 reviewer request
// (hold the durable path past the invocation deadline, stall GetReceipt,
// AppendIntent and Commit separately) proven at the HTTP request level with
// REPEATED calls on one handler: each request completes within the invoke
// bound, answers the exact R11 (definite) or R16 (uncertain) refusal with
// Connection: close, no blocked worker accumulates, nothing commits, and a
// normal call on the same handler then succeeds.
func TestA2ADurableDeadlineWire(t *testing.T) {
	// invokeWait also covers admission and the ops before the held one, so it
	// is wide enough for a loaded 2-core -race runner; bound only has to
	// separate "deadline honoured" from "hung" (a mutant that drops the
	// deadline blocks until the 5 s release), so it is generous.
	const invokeWait = 300 * time.Millisecond
	const bound = invokeWait + 3*time.Second
	const msgR11 = "invocation exceeded its deadline"
	const msgR16 = "invocation outcome is not confirmed; resend the same task id"
	for _, tc := range []struct {
		name, op  string
		uncertain bool
		wantMsg   string
	}{
		{"GetReceipt", "GetReceipt", false, msgR11},
		{"AppendIntent", "AppendIntent", false, msgR11},
		{"Commit", "Commit", false, msgR11},
		{"uncertain_AppendIntent", "AppendIntent", true, msgR16},
		{"uncertain_Commit", "Commit", true, msgR16},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := newWireRig(t)
			var armed atomic.Bool
			var stalls atomic.Int32
			armed.Store(true)
			st := stallStore{Store: r.st, op: tc.op, uncertain: tc.uncertain, armed: &armed, stalls: &stalls, release: make(chan struct{})}
			var once sync.Once
			releaseAll := func() { once.Do(func() { close(st.release) }) }
			t.Cleanup(releaseAll)
			h := r.handler(t, st, hookRunner{stdout: `{"ok":true}`, hook: func(context.Context) error { return nil }}, nil, invokeWait)
			base := runtime.NumGoroutine()
			for i := 0; i < 3; i++ {
				reqID := fmt.Sprintf("held-%s-%d", tc.name, i)
				done := make(chan *recorded, 1)
				start := time.Now()
				go func() { done <- post(t, h, "Bearer "+r.tok, wireBody(reqID, fmt.Sprintf("t%d", i))) }()
				var rec *recorded
				select {
				case rec = <-done:
				case <-time.After(bound):
					releaseAll()
					<-done
					t.Fatalf("call %d did not complete within %v: the invocation deadline did not bound the stalled %s", i, bound, tc.op)
				}
				if elapsed := time.Since(start); elapsed < invokeWait {
					t.Fatalf("call %d returned after %v, before the %v deadline: the stall was not held", i, elapsed, invokeWait)
				}
				assertWireRefusal(t, rec, reqID, codeInternal, tc.wantMsg, true)
			}
			if got := stalls.Load(); got != 3 {
				t.Fatalf("stalled %s reached %d times, want 3 (one per request)", tc.op, got)
			}
			// Slack 2 < 3 requests: one leaked worker per request is still caught,
			// while lazily started runtime/database helpers cannot flake it.
			if n := settledGoroutineCount(base); n > base+2 {
				t.Fatalf("goroutines %d > baseline %d (+2 slack) after three held requests: a blocked worker accumulated", n, base)
			}
			head, _, err := r.st.SelectedHead(boundedTestContext(t))
			if err != nil || head != r.world.Ref {
				t.Fatalf("head=%v err=%v, want genesis %v: a timed-out invocation committed", head, err, r.world.Ref)
			}
			// Positive control: the same handler, released, still invokes.
			armed.Store(false)
			rec := post(t, h, "Bearer "+r.tok, wireBody("held-control", "t-control"))
			if rec.code != http.StatusOK || !strings.Contains(string(rec.body), `"state":"completed"`) || !strings.Contains(string(rec.body), `"id":"held-control"`) {
				t.Fatalf("control invocation after held requests: %d %s", rec.code, rec.body)
			}
			if c := rec.header.Get("Connection"); c != "" {
				t.Fatalf("control success Connection=%q, want absent", c)
			}
		})
	}
}

func assertRefusalLog(t *testing.T, log *bytes.Buffer, id, detail string) {
	t.Helper()
	if bytes.Count(log.Bytes(), []byte("\n")) != 1 || !strings.Contains(log.String(), "ailang-worldd: a2a refusal: tasks/send "+id+": \"") || !strings.Contains(log.String(), detail) {
		t.Fatalf("operator log=%q want one line carrying %q", log.String(), detail)
	}
}

type absentRecordStore struct {
	*store.Store
	hits *int
}

func (s absentRecordStore) GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error) {
	o, ok, err := s.Store.GetObject(ctx, ref)
	if ok && o.SemanticID == coordinator.RecordV1 {
		*s.hits++
		return store.Object{}, false, nil
	}
	return o, ok, err
}

type intentOnlyStore struct{ *store.Store }

func (s intentOnlyStore) GetReceipt(ctx context.Context, id string) (store.Receipt, bool, error) {
	return store.Receipt{InvocationID: id, State: store.ReceiptIndeterminate}, true, nil
}

// phaseRunnerFn answers a convention-v2 phase from fn (the phase name and the
// invocation ctx), so a test can stall exactly one phase.
type phaseRunnerFn func(ctx context.Context, phase string) (string, error)

func (f phaseRunnerFn) RunContext(ctx context.Context, e capsule.Entry) (capsule.Result, error) {
	var arg string
	if err := json.Unmarshal(e.Args, &arg); err != nil {
		return capsule.Result{}, err
	}
	var in struct {
		Phase string `json:"phase"`
	}
	if err := json.Unmarshal([]byte(arg), &in); err != nil {
		return capsule.Result{}, err
	}
	out, err := f(ctx, in.Phase)
	return capsule.Result{Stdout: []byte(out + "\n")}, err
}

type okEffect struct{}

func (okEffect) Execute(context.Context, broker.EffectRequest, []byte) ([]byte, error) {
	return []byte(`{"ok":true}`), nil
}

// effectfulPlan is a one-effect plan on the wire rig's alpha@world grant.
func effectfulPlan(finish bool) string {
	return fmt.Sprintf(`{"plan":"world/effect-plan/v1","effects":[{"id":"e1","effect":"alpha","scope":"world","cost":1,"payload":{}}],"finish":%t,"result":null}`, finish)
}

func effectfulHandler(t *testing.T, r *wireRig, runner coordinator.Runner, invokeWait time.Duration) *Handler {
	t.Helper()
	return r.handler(t, nil, runner, func(ep string, grants []broker.Capability) transitionreg.Binder {
		return broker.OpenBinder(r.st, ep, grants, broker.Registry{"alpha": okEffect{}})
	}, invokeWait)
}

// TestA2APlanTimeoutIsTypedAndNamesThePhase is row 153 AC1.3 on the wire: a
// plan that blocks until ITS budget ends, with the caller still live, answers
// -32603 with the plan-phase message (not the generic deadline text), and the
// connection is not closed (the invocation deadline did not expire).
// Waits the real PlanPhaseBudget (P1).
func TestA2APlanTimeoutIsTypedAndNamesThePhase(t *testing.T) {
	r := newWireRigWith(t, true)
	var plans atomic.Int32
	runner := phaseRunnerFn(func(ctx context.Context, phase string) (string, error) {
		if phase != "plan" {
			return "", errors.New("unexpected phase " + phase)
		}
		plans.Add(1)
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(coordinator.PlanPhaseBudget + 5*time.Second): // a budget-less mutant fails, never hangs
			return "", errors.New("the plan budget never cancelled the capsule")
		}
	})
	h := effectfulHandler(t, r, runner, 20*time.Second)
	var log bytes.Buffer
	h.errorLog = &log
	rec := post(t, h, "Bearer "+r.tok, wireBody("wire-plan-timeout", "t-plan"))
	if plans.Load() != 1 {
		t.Fatalf("plan ran %d times, want 1", plans.Load())
	}
	want := PhaseTimeoutPrefix + " " + coordinator.PlanPhaseBudget.String() + " budget; nothing ran or was committed; resend the same task id"
	assertWireRefusal(t, rec, "wire-plan-timeout", codeInternal, want, false)
}

// TestA2AFinishOverrunStaysEffectsUnrecorded is row 153 AC1.7: a finish
// phase that overruns AFTER the effect ran is an effects-unrecorded answer
// (the effect record is durable), never the plan-timeout "resend" message.
// Waits the real FinishPhaseBudget (P1).
func TestA2AFinishOverrunStaysEffectsUnrecorded(t *testing.T) {
	r := newWireRigWith(t, true)
	runner := phaseRunnerFn(func(ctx context.Context, phase string) (string, error) {
		if phase == "plan" {
			return effectfulPlan(true), nil
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(coordinator.FinishPhaseBudget + 5*time.Second):
			return "", errors.New("the finish budget never cancelled the capsule")
		}
	})
	h := effectfulHandler(t, r, runner, 20*time.Second)
	rec := post(t, h, "Bearer "+r.tok, wireBody("wire-finish-overrun", "t-finish"))
	code, msg, _ := a2aErr(t, rec.body)
	if code != codeInternal || !strings.HasPrefix(msg, EffectsUnrecordedPrefix) {
		t.Fatalf("answer = %d %q, want %d with prefix %q; body=%s", code, msg, codeInternal, EffectsUnrecordedPrefix, rec.body)
	}
	if strings.HasPrefix(msg, PhaseTimeoutPrefix) || msg == "invocation exceeded its deadline" {
		t.Fatalf("finish overrun answered %q: it must stay effects-unrecorded", msg)
	}
}
