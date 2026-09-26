package coordinator

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/canon"
	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

const (
	now         = int64(100)
	accessScope = "world"
)

var grants = []broker.Capability{{Effect: "Invoke", Scope: accessScope, ExpiresAt: 1000, Budget: 10}}

// ---------- fixture ----------

type rig struct {
	t      *testing.T
	st     *store.Store
	interp hashref.HashRef
	runner Runner
	descs  []transitionreg.Descriptor
}

func newRig(t *testing.T, withRealInterpreter bool) *rig {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "world.db")
	st, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	r := &rig{t: t, st: st, interp: hashref.SumSHA256([]byte("fake-interpreter"))}
	if withRealInterpreter {
		bin := os.Getenv("AILANG_BIN")
		if bin == "" {
			t.Skip("AILANG_BIN not set; real-interpreter invocation needs the pinned ailang")
		}
		a := archive.New(dbPath)
		ref, err := a.Archive(bin)
		if err != nil {
			t.Fatalf("archive interpreter: %v", err)
		}
		r.interp = ref
		r.runner = capsule.New(a, capsule.Config{})
	}
	return r
}

func (r *rig) genesis() store.World {
	r.t.Helper()
	obj := store.Object{Hash: hashref.SumSHA256([]byte("genesis-state")), InterfaceHash: hashref.SumSHA256([]byte("test/genesis")),
		SemanticID: "test/genesis", Provenance: "coordinator-test", Payload: []byte("genesis-state")}
	entryHash := hashref.SumSHA256([]byte("genesis-entry"))
	w := store.World{Ref: hashref.SumSHA256([]byte("genesis-world")), Revision: 0, StateRoot: obj.Hash, LogHead: entryHash}
	if err := r.st.Commit(store.Commit{
		Objects: []store.Object{obj}, NextWorld: w,
		Entry: store.LogEntry{Header: store.LogHeader{
			EntryIndex: 0, SemanticsEpoch: 1, TransitionFn: hashref.SumSHA256([]byte("genesis-fn")),
			Interpreter: r.interp, PrevEntryHash: hashref.SumSHA256([]byte("genesis-prev")), WrittenBy: "coordinator-test",
		}, EntryHash: entryHash, TransitionRef: obj.Hash},
	}); err != nil {
		r.t.Fatalf("genesis commit: %v", err)
	}
	return w
}

// source stores canonical transition source under row 107's convention.
func (r *rig) source(raw string) hashref.HashRef {
	r.t.Helper()
	b, err := canon.Source([]byte(raw))
	if err != nil {
		r.t.Fatalf("canon: %v", err)
	}
	obj := store.Object{Hash: hashref.SumSHA256(b), InterfaceHash: hashref.SumSHA256([]byte("world/transition-source/v1")),
		SemanticID: "world/transition-source/v1", Provenance: "coordinator-test", Payload: b}
	if err := r.st.PutObject(obj); err != nil {
		r.t.Fatalf("put source: %v", err)
	}
	return obj.Hash
}

func (r *rig) describe(id string, fn hashref.HashRef, access string, declared ...transitionreg.EffectRequirement) {
	r.descs = append(r.descs, transitionreg.Descriptor{
		ID: id, TransitionFn: fn, Interpreter: r.interp, SemanticsEpoch: 1,
		InputSchema: []byte(`{}`), OutputSchema: []byte(`{}`),
		Access:          transitionreg.EffectRequirement{Effect: access, Scope: accessScope, Cost: 1},
		DeclaredEffects: declared, Title: id, Description: id,
	})
}

func (r *rig) publish() {
	r.t.Helper()
	sort.Slice(r.descs, func(i, j int) bool { return r.descs[i].ID < r.descs[j].ID })
	payload, err := transitionreg.EncodeRevision(transitionreg.Revision{
		SemanticID: transitionreg.SemanticIDV1, InterfaceHash: transitionreg.InterfaceHashV1, Revision: 1, Entries: r.descs,
	})
	if err != nil {
		r.t.Fatalf("encode revision: %v", err)
	}
	obj := store.Object{Hash: hashref.SumSHA256(payload), InterfaceHash: transitionreg.InterfaceHashV1,
		SemanticID: transitionreg.SemanticIDV1, Provenance: "coordinator-test", Payload: payload}
	if err := r.st.PutObject(obj); err != nil {
		r.t.Fatalf("put revision: %v", err)
	}
	if err := r.st.CompareAndSetRegistryHead(store.TransitionRegistryV1, hashref.HashRef{}, obj.Hash); err != nil {
		r.t.Fatalf("cas head: %v", err)
	}
}

type capSource []broker.Capability

func (c capSource) CapabilitySnapshot(n int64) broker.CapabilitySnapshot {
	return broker.NewCapabilitySnapshot(c, n)
}

func (r *rig) coordinator(st Store, runner Runner) *Coordinator {
	r.t.Helper()
	c, err := New(Config{
		Store: st, Runner: runner, Now: func() int64 { return now }, MaxInput: 1 << 16, MaxOutput: 1 << 16,
		Binder: func(ep string, caps []broker.Capability) transitionreg.Binder {
			return broker.OpenBinder(r.st, ep, caps)
		},
	})
	if err != nil {
		r.t.Fatalf("New: %v", err)
	}
	return c
}

func (r *rig) call(skill, task string, caps []broker.Capability) Call {
	r.t.Helper()
	req, err := transitionreg.NewRequest(context.Background(), transitionreg.NewReader(r.st), capSource(caps), now)
	if err != nil {
		r.t.Fatalf("NewRequest: %v", err)
	}
	return Call{Request: req, EpisodeID: "ep1", Grants: caps, SkillID: skill, TaskID: task, Input: map[string]any{"msg": "hi"}}
}

// assertUntouched: the refusal left no durable mutation.
func (r *rig) assertUntouched(want hashref.HashRef, task string) {
	r.t.Helper()
	head, _, err := r.st.SelectedHead(context.Background())
	if err != nil || head != want {
		r.t.Fatalf("selected head = %v (%v), want unchanged %v", head, err, want)
	}
	if rc, ok, err := r.st.GetReceipt(InvocationID("ep1", task)); err != nil || (ok && rc.State != store.ReceiptNotStarted) {
		r.t.Fatalf("receipt for refused task = %+v ok=%v err=%v, want none", rc, ok, err)
	}
}

const echoSrc = "module transitions/echo\n\nexport func main(input: string) -> string {\n  \"{\\\"echo\\\":${input}}\"\n}\n"

// fakeRunner returns fixed stdout or a fixed error, optionally running a hook.
type fakeRunner struct {
	stdout string
	err    error
	hook   func(ctx context.Context)
}

func (f fakeRunner) RunContext(ctx context.Context, _ capsule.Entry) (capsule.Result, error) {
	if f.hook != nil {
		f.hook(ctx)
	}
	return capsule.Result{Stdout: []byte(f.stdout)}, f.err
}

// ---------- pure core: the kernel laws ----------

func lawPlan() (store.World, transitionreg.Descriptor, plan) {
	w := store.World{Ref: hashref.SumSHA256([]byte("w")), Revision: 7,
		StateRoot: hashref.SumSHA256([]byte("s")), LogHead: hashref.SumSHA256([]byte("l"))}
	d := transitionreg.Descriptor{ID: "x", TransitionFn: hashref.SumSHA256([]byte("fn")),
		Interpreter: hashref.SumSHA256([]byte("i")), SemanticsEpoch: 1}
	return w, d, planInvocation(w, "a2a:e:t", "e", d, []byte(`{"a":1}`), []byte(`{"b":2}`), 5)
}

// The record's skillId is the descriptor's ID — the registry key the card
// emits as skills[].id (projection.go:240).
func TestPlanRecordSkillID(t *testing.T) {
	_, d, p := lawPlan()
	var rec record
	if err := json.Unmarshal(p.Objects[2].Payload, &rec); err != nil || p.Objects[2].SemanticID != RecordV1 {
		t.Fatalf("record object: %v / %s", err, p.Objects[2].SemanticID)
	}
	if rec.SkillID != d.ID || rec.SkillID == "" {
		t.Fatalf("record skillId = %q, want descriptor ID %q", rec.SkillID, d.ID)
	}
}

func TestPlanLawRevision(t *testing.T) {
	w, _, p := lawPlan()
	if p.Commit.NextWorld.Revision != w.Revision+1 || p.Commit.Entry.Header.EntryIndex != w.Revision+1 {
		t.Fatalf("applyRevision: next revision %d / entry %d, want %d", p.Commit.NextWorld.Revision, p.Commit.Entry.Header.EntryIndex, w.Revision+1)
	}
}

func TestPlanLawStateRoot(t *testing.T) {
	_, _, p := lawPlan()
	if p.Commit.NextWorld.StateRoot != hashref.SumSHA256([]byte(`{"b":2}`)) {
		t.Fatal("applyRevision: stateRoot is not the output object")
	}
}

func TestPlanLawLogHead(t *testing.T) {
	_, _, p := lawPlan()
	if p.Commit.NextWorld.LogHead != p.Commit.Entry.EntryHash || p.Commit.Entry.TransitionRef != p.Record {
		t.Fatal("applyRevision: logHead is not the new entry / entry does not reference the record")
	}
}

func TestPlanLawPrevEntry(t *testing.T) {
	w, _, p := lawPlan()
	if p.Commit.Entry.Header.PrevEntryHash != w.LogHead || p.Intent.PrevEntryHash != w.LogHead {
		t.Fatal("entry does not chain from the observed log head")
	}
}

func TestPlanLawObservedHead(t *testing.T) {
	w, _, p := lawPlan()
	if p.Commit.ObservedHead != w.Ref || p.Intent.ObservedHead != w.Ref || p.Intent.WorldRef != p.Commit.NextWorld.Ref {
		t.Fatal("proposalMatchesWorld: the commit/intent do not observe exactly the planned-against world")
	}
}

// ---------- real interpreter ----------

// AC-HAPPY.
func TestDispatchEchoRealInterpreter(t *testing.T) {
	r := newRig(t, true)
	w := r.genesis()
	r.describe("echo", r.source(echoSrc), "Invoke")
	r.publish()
	res, err := r.coordinator(r.st, r.runner).Dispatch(context.Background(), r.call("echo", "t1", grants))
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	if string(res.OutputBytes) != `{"echo":{"msg":"hi"}}` {
		t.Fatalf("output = %s", res.OutputBytes)
	}
	head, _, _ := r.st.SelectedHead(context.Background())
	if head != res.WorldRef || res.EntryIndex != w.Revision+1 {
		t.Fatalf("head %v entry %d, want %v / %d", head, res.EntryIndex, res.WorldRef, w.Revision+1)
	}
	entry, ok, err := r.st.GetLogEntry(context.Background(), res.EntryIndex)
	if err != nil || !ok || entry.Header.TransitionFn != r.descs[0].TransitionFn || entry.Header.Interpreter != r.interp {
		t.Fatalf("log entry %+v ok=%v err=%v", entry, ok, err)
	}
	rc, ok, err := r.st.GetReceipt(res.InvocationID)
	if err != nil || !ok || rc.State != store.ReceiptResolved {
		t.Fatalf("receipt %+v ok=%v err=%v, want resolved", rc, ok, err)
	}
}

// AC-REPLAY: the committed pins + input reproduce the committed output bytes.
func TestReplayCommittedInvocation(t *testing.T) {
	r := newRig(t, true)
	r.genesis()
	r.describe("echo", r.source(echoSrc), "Invoke")
	r.publish()
	res, err := r.coordinator(r.st, r.runner).Dispatch(context.Background(), r.call("echo", "t1", grants))
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	ctx := context.Background()
	recObj, ok, err := r.st.GetObject(ctx, res.RecordRef)
	if err != nil || !ok {
		t.Fatalf("record: ok=%v err=%v", ok, err)
	}
	var rec record
	if err := json.Unmarshal(recObj.Payload, &rec); err != nil {
		t.Fatalf("decode record: %v", err)
	}
	get := func(text string) []byte {
		ref := hashref.MustParse(text)
		o, ok, err := r.st.GetObject(ctx, ref)
		if err != nil || !ok {
			t.Fatalf("object %s: ok=%v err=%v", text, ok, err)
		}
		return o.Payload
	}
	out, err := r.runner.RunContext(ctx, capsule.Entry{
		Interpreter: hashref.MustParse(rec.Interpreter), Source: get(rec.TransitionFn), Args: mustJSON(string(get(rec.Input))),
	})
	if err != nil {
		t.Fatalf("re-execute: %v", err)
	}
	if got, want := string(out.Stdout), string(get(rec.Output))+"\n"; got != want {
		t.Fatalf("replay diverged: %q vs recorded %q", got, want)
	}
}

// AC-INCOMPAT (Residual 8(c)): checks under the pinned interpreter (V10) but
// cannot take an argument.
func TestDispatchIncompatibleRealInterpreter(t *testing.T) {
	r := newRig(t, true)
	w := r.genesis()
	r.describe("zero", r.source("module transitions/zero\n\nexport func main() -> string {\n  \"{}\"\n}\n"), "Invoke")
	r.publish()
	_, err := r.coordinator(r.st, r.runner).Dispatch(context.Background(), r.call("zero", "t1", grants))
	var inc *IncompatibleError
	if !errors.As(err, &inc) {
		t.Fatalf("err = %T %v, want *IncompatibleError", err, err)
	}
	r.assertUntouched(w.Ref, "t1")
}

// ---------- every refusal branch (fake runner unless stated) ----------

type corruptStore struct{ *store.Store }

func (c corruptStore) GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error) {
	o, ok, err := c.Store.GetObject(ctx, ref)
	if ok && o.SemanticID == "world/transition-source/v1" {
		o.Payload = append([]byte(nil), o.Payload...)
		o.Payload[0] ^= 1
	}
	return o, ok, err
}

type failingBinder struct{ err error }

func (b failingBinder) Bind(broker.Manifest) (*broker.BoundInvoker, error) { return nil, b.err }

func TestDispatchRefuses(t *testing.T) {
	okRun := fakeRunner{stdout: `{"ok":true}` + "\n"}
	setup := func(t *testing.T) (*rig, store.World) {
		r := newRig(t, false)
		w := r.genesis()
		r.describe("plain", r.source(echoSrc), "Invoke")
		r.describe("fx", r.source(echoSrc+"\n"), "Invoke", transitionreg.EffectRequirement{Effect: "Net", Scope: "x", Cost: 1})
		r.describe("ghost", hashref.SumSHA256([]byte("never stored")), "Invoke")
		r.describe("locked", r.source(echoSrc+"\n\n"), "Admin")
		r.publish()
		return r, w
	}
	cases := []struct {
		name   string
		skill  string
		task   string
		runner Runner
		mutate func(*Call)
		store  func(*rig) Store
		check  func(error) bool
	}{
		{"R1_bad_task_id", "plain", "bad id!", okRun, nil, nil, func(e error) bool { var x *InvalidCallError; return errors.As(e, &x) }},
		{"R1_empty_episode", "plain", "t1", okRun, func(c *Call) { c.EpisodeID = "" }, nil, func(e error) bool { var x *InvalidCallError; return errors.As(e, &x) }},
		{"R1_input_over_cap", "plain", "t1", okRun, func(c *Call) { c.Input = map[string]any{"oversized": string(make([]byte, 1<<17))} }, nil, func(e error) bool { var x *InvalidCallError; return errors.As(e, &x) }},
		{"R1_unmarshalable_input", "plain", "t1", okRun, func(c *Call) { c.Input = map[string]any{"bad": make(chan int)} }, nil, func(e error) bool { var x *InvalidCallError; return errors.As(e, &x) }},
		{"R1_nil_input", "plain", "t1", okRun, func(c *Call) { c.Input = nil }, nil, func(e error) bool { var x *InvalidCallError; return errors.As(e, &x) }},
		{"R2_absent", "nosuch", "t1", okRun, nil, nil, func(e error) bool { var x *transitionreg.TransitionAbsentError; return errors.As(e, &x) }},
		{"R3_access_denied", "locked", "t1", okRun, nil, nil, func(e error) bool { var x *transitionreg.AccessDeniedError; return errors.As(e, &x) }},
		{"R5_pin_mismatch", "plain", "t1", okRun, func(c *Call) { h := hashref.SumSHA256([]byte("other")); c.PinnedFn = &h }, nil,
			func(e error) bool { var x *transitionreg.ProposalMismatchError; return errors.As(e, &x) }},
		{"R6_absent_source", "ghost", "t1", okRun, nil, nil, func(e error) bool { var x *SourceError; return errors.As(e, &x) && x.Kind == "absent" }},
		{"R6_corrupt_source", "plain", "t1", okRun, nil, func(r *rig) Store { return corruptStore{r.st} },
			func(e error) bool { var x *SourceError; return errors.As(e, &x) && x.Kind == "corrupt" }},
		{"R8_effects_declared", "fx", "t1", okRun, nil, nil, func(e error) bool { var x *EffectsUnsupportedError; return errors.As(e, &x) }},
		{"R9_incompatible", "plain", "t1", fakeRunner{err: &capsule.ExecError{Stderr: []byte("ARG_DECODE_MISMATCH: expected ()")}}, nil, nil,
			func(e error) bool { var x *IncompatibleError; return errors.As(e, &x) }},
		{"R10_exec_failed", "plain", "t1", fakeRunner{err: &capsule.ExecError{Stderr: []byte("panic")}}, nil, nil,
			func(e error) bool { var x *ExecutionError; return errors.As(e, &x) }},
		{"R12_output_not_object", "plain", "t1", fakeRunner{stdout: "[1,2]\n"}, nil, nil, func(e error) bool { var x *OutputError; return errors.As(e, &x) }},
		{"R12_output_over_cap", "plain", "t1", fakeRunner{stdout: `{"x":"` + strings.Repeat("x", 1<<17) + `"}`}, nil, nil, func(e error) bool { var x *OutputError; return errors.As(e, &x) }},
		{"R12_output_not_json", "plain", "t1", fakeRunner{stdout: "echo:hi\n"}, nil, nil, func(e error) bool { var x *OutputError; return errors.As(e, &x) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r, w := setup(t)
			call := r.call(tc.skill, tc.task, grants)
			if tc.mutate != nil {
				tc.mutate(&call)
			}
			var st Store = r.st
			if tc.store != nil {
				st = tc.store(r)
			}
			_, err := r.coordinator(st, tc.runner).Dispatch(context.Background(), call)
			if !tc.check(err) {
				t.Fatalf("err = %T %v", err, err)
			}
			r.assertUntouched(w.Ref, tc.task)
		})
	}

	t.Run("R4_bind_error", func(t *testing.T) {
		r, w := setup(t)
		marker := errors.New("binder failed")
		c, err := New(Config{Store: r.st, Runner: fakeRunner{hook: func(context.Context) { t.Error("executed after binder refusal") }}, Now: func() int64 { return now }, MaxInput: 1 << 16, MaxOutput: 1 << 16,
			Binder: func(string, []broker.Capability) transitionreg.Binder { return failingBinder{marker} }})
		if err != nil {
			t.Fatal(err)
		}
		_, err = c.Dispatch(context.Background(), r.call("plain", "t1", grants))
		if !errors.Is(err, marker) {
			t.Fatalf("binder error = %T %v", err, err)
		}
		r.assertUntouched(w.Ref, "t1")
	})

	t.Run("R7_no_world", func(t *testing.T) {
		r := newRig(t, false)
		r.describe("plain", r.source(echoSrc), "Invoke")
		r.publish()
		_, err := r.coordinator(r.st, okRun).Dispatch(context.Background(), r.call("plain", "t1", grants))
		var x *WorldAbsentError
		if !errors.As(err, &x) {
			t.Fatalf("err = %T %v, want *WorldAbsentError", err, err)
		}
	})

	t.Run("R11_cancelled_during_exec", func(t *testing.T) {
		r, w := setup(t)
		ctx, cancel := context.WithCancel(context.Background())
		runner := fakeRunner{err: &capsule.ExecError{Stderr: []byte("ARG_DECODE_MISMATCH")}, hook: func(context.Context) { cancel() }}
		_, err := r.coordinator(r.st, runner).Dispatch(ctx, r.call("plain", "t1", grants))
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("err = %T %v, want context.Canceled (not a classified exec failure)", err, err)
		}
		r.assertUntouched(w.Ref, "t1")
	})

	// A resent task id is answered from the journal: same result, nothing
	// re-executed, head unchanged.
	t.Run("R13_resend_reconciles_committed", func(t *testing.T) {
		r, _ := setup(t)
		first, err := r.coordinator(r.st, okRun).Dispatch(context.Background(), r.call("plain", "t1", grants))
		if err != nil {
			t.Fatalf("first Dispatch: %v", err)
		}
		mustNotRun := fakeRunner{hook: func(context.Context) { t.Error("resend re-executed the transition") }, stdout: `{"second":true}`}
		again, err := r.coordinator(r.st, mustNotRun).Dispatch(context.Background(), r.call("plain", "t1", grants))
		if err != nil {
			t.Fatalf("resend: %T %v, want the committed result", err, err)
		}
		if !again.Reconciled || again.WorldRef != first.WorldRef || again.EntryIndex != first.EntryIndex ||
			string(again.OutputBytes) != string(first.OutputBytes) || again.RecordRef != first.RecordRef {
			t.Fatalf("resend = %+v, want the first call's committed result %+v", again, first)
		}
		r.assertHead(first.WorldRef)
	})

	// An intent without an outcome (here: a head-moved conflict) is answered
	// NotCommitted, and the transition is not re-run.
	t.Run("R15_resend_not_committed", func(t *testing.T) {
		r, w := setup(t)
		mover := fakeRunner{stdout: `{"ok":true}`, hook: func(context.Context) {
			if _, err := r.coordinator(r.st, fakeRunner{stdout: `{"other":1}`}).Dispatch(context.Background(), r.call("plain", "other", grants)); err != nil {
				t.Errorf("concurrent writer: %v", err)
			}
		}}
		if _, err := r.coordinator(r.st, mover).Dispatch(context.Background(), r.call("plain", "t1", grants)); !store.IsConflict(err) {
			t.Fatalf("first: %v, want conflict", err)
		}
		_ = w
		mustNotRun := fakeRunner{hook: func(context.Context) { t.Error("resend re-executed the transition") }, stdout: `{}`}
		_, err := r.coordinator(r.st, mustNotRun).Dispatch(context.Background(), r.call("plain", "t1", grants))
		var x *NotCommittedError
		if !errors.As(err, &x) {
			t.Fatalf("resend err = %T %v, want *NotCommittedError", err, err)
		}
	})

	// An untyped Commit error is UNCONFIRMED, never "failed": here the commit
	// actually landed, and a resend reconciles to the committed result.
	t.Run("R16_unconfirmed_then_reconciled", func(t *testing.T) {
		r, w := setup(t)
		landedButErr := errStore{Store: r.st, afterCommit: errors.New("disk I/O error")}
		_, err := r.coordinator(landedButErr, okRun).Dispatch(context.Background(), r.call("plain", "t1", grants))
		var x *UnconfirmedError
		if !errors.As(err, &x) {
			t.Fatalf("err = %T %v, want *UnconfirmedError", err, err)
		}
		again, err := r.coordinator(r.st, okRun).Dispatch(context.Background(), r.call("plain", "t1", grants))
		if err != nil || !again.Reconciled || again.EntryIndex != w.Revision+1 {
			t.Fatalf("resend = %+v, %v; want the landed commit reconciled", again, err)
		}
	})

	t.Run("R14_head_moved", func(t *testing.T) {
		r, w := setup(t)
		mover := fakeRunner{stdout: `{"ok":true}`, hook: func(context.Context) {
			// Another writer advances the world while the capsule runs.
			other := r.coordinator(r.st, fakeRunner{stdout: `{"other":1}`})
			if _, err := other.Dispatch(context.Background(), r.call("plain", "other", grants)); err != nil {
				t.Errorf("concurrent writer: %v", err)
			}
		}}
		_, err := r.coordinator(r.st, mover).Dispatch(context.Background(), r.call("plain", "t1", grants))
		var unconfirmed *UnconfirmedError
		if !store.IsConflict(err) || errors.As(err, &unconfirmed) {
			t.Fatalf("err = %T %v, want a bare ConflictError (definitively not committed), not Unconfirmed", err, err)
		}
		head, _, _ := r.st.SelectedHead(context.Background())
		if head == w.Ref {
			t.Fatal("the concurrent writer's commit did not land (fixture broken)")
		}
	})
}

// errStore wraps the real store; afterCommit makes Commit land and then
// report an error; onCommit runs just before the real Commit.
type errStore struct {
	*store.Store
	afterCommit error
	onCommit    func()
}

func (e errStore) Commit(c store.Commit) error {
	if e.onCommit != nil {
		e.onCommit()
	}
	if err := e.Store.Commit(c); err != nil {
		return err
	}
	return e.afterCommit
}

func (r *rig) assertHead(want hashref.HashRef) {
	r.t.Helper()
	head, _, err := r.st.SelectedHead(context.Background())
	if err != nil || head != want {
		r.t.Fatalf("head = %v (%v), want %v", head, err, want)
	}
}

// AC-BOUNDARY: both sides of the commit boundary.
func TestCommitBoundary(t *testing.T) {
	t.Run("cancel_before_commit", func(t *testing.T) {
		r := newRig(t, false)
		w := r.genesis()
		r.describe("plain", r.source(echoSrc), "Invoke")
		r.publish()
		ctx, cancel := context.WithCancel(context.Background())
		// The capsule SUCCEEDS; the caller cancels before the commit boundary.
		runner := fakeRunner{stdout: `{"ok":true}`, hook: func(context.Context) { cancel() }}
		_, err := r.coordinator(r.st, runner).Dispatch(ctx, r.call("plain", "t1", grants))
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("err = %T %v, want context.Canceled", err, err)
		}
		r.assertUntouched(w.Ref, "t1")
		if pending, err := r.st.PendingIntents(10); err != nil || len(pending) != 0 {
			t.Fatalf("pending intents = %v (%v), want none", pending, err)
		}
	})
	// The durable steps take no ctx; a deadline that expires while they run
	// must not turn a landed commit into a reported failure.
	t.Run("deadline_during_commit_reports_success", func(t *testing.T) {
		r := newRig(t, false)
		r.genesis()
		r.describe("plain", r.source(echoSrc), "Invoke")
		r.publish()
		ctx, cancel := context.WithCancel(context.Background())
		st := errStore{Store: r.st, onCommit: cancel}
		res, err := r.coordinator(st, fakeRunner{stdout: `{"ok":true}`}).Dispatch(ctx, r.call("plain", "t1", grants))
		if err != nil {
			t.Fatalf("Dispatch = %v after a landed commit; must report success", err)
		}
		r.assertHead(res.WorldRef)
	})
	t.Run("receipt", func(t *testing.T) {
		r := newRig(t, false)
		r.genesis()
		r.describe("plain", r.source(echoSrc), "Invoke")
		r.publish()
		res, err := r.coordinator(r.st, fakeRunner{stdout: `{"ok":true}`}).Dispatch(context.Background(), r.call("plain", "t1", grants))
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		rc, ok, err := r.st.GetReceipt(res.InvocationID)
		if err != nil || !ok || rc.State != store.ReceiptResolved || rc.Intent == nil || rc.Intent.WorldRef != res.WorldRef {
			t.Fatalf("receipt = %+v ok=%v err=%v, want exactly one resolved intent+outcome for the commit", rc, ok, err)
		}
	})
}

func TestNewRefusesMissingSeamsAndCaps(t *testing.T) {
	valid := Config{Store: (*store.Store)(nil), Runner: fakeRunner{}, Binder: func(string, []broker.Capability) transitionreg.Binder { return nil }, Now: func() int64 { return 1 }, MaxInput: 10, MaxOutput: 10}
	if _, err := New(Config{}); err == nil {
		t.Fatal("accepted missing seams")
	}
	checks := []struct {
		name   string
		change func(*Config)
	}{
		{"store", func(c *Config) { c.Store = nil }}, {"runner", func(c *Config) { c.Runner = nil }},
		{"binder", func(c *Config) { c.Binder = nil }}, {"clock", func(c *Config) { c.Now = nil }},
		{"input cap", func(c *Config) { c.MaxInput = 0 }}, {"output cap", func(c *Config) { c.MaxOutput = -1 }},
	}
	for _, tc := range checks {
		t.Run(tc.name, func(t *testing.T) {
			c := valid
			tc.change(&c)
			if _, err := New(c); err == nil {
				t.Fatal("accepted invalid config")
			}
		})
	}
	c, err := New(valid)
	if err != nil {
		t.Fatal(err)
	}
	if InvocationID("ep", "task") != "a2a:ep:task" {
		t.Fatal("wrong invocation namespace")
	}
	if validTaskID("bad id") || !validTaskID("a-1") {
		t.Fatal("task validation")
	}
	if _, _, err := c.parseOutput([]byte("[]")); err == nil {
		t.Fatal("accepted array output")
	}
	if _, _, err := c.parseOutput([]byte(`{"a":1}`)); err != nil {
		t.Fatal(err)
	}
	execErr := &capsule.ExecError{Stderr: []byte("ARG_DECODE_MISMATCH")}
	var incompatible *IncompatibleError
	if !errors.As(classifyExec(execErr), &incompatible) {
		t.Fatal("wrong execution classification")
	}
}
