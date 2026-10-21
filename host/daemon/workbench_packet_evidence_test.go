package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/evidence"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
	"github.com/sunholo-data/ailang-world/host/workbench"
)

// Forward real reads; controls alter one field after the genuine producer/store.
type rosterReadControl struct {
	readStore
	edit                          func(hashref.HashRef, *store.Object)
	readErr                       error
	missing                       hashref.HashRef
	entryEdit                     func(*store.LogEntry)
	entryErr                      error
	entryMissing                  bool
	waitAt                        hashref.HashRef
	deadlines                     []time.Time
	commitMissing, commitMismatch bool
}

func (r *rosterReadControl) GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		panic("roster read without deadline")
	}
	r.deadlines = append(r.deadlines, deadline)
	o, found, err := r.readStore.GetObject(ctx, ref)
	if err == nil && found && r.edit != nil {
		o.Payload = append([]byte(nil), o.Payload...)
		r.edit(ref, &o)
	}
	if ref == r.waitAt {
		<-ctx.Done()
	}
	if r.readErr != nil {
		return o, found, r.readErr
	}
	if ref == r.missing {
		return store.Object{}, false, nil
	}
	return o, found, err
}
func (r *rosterReadControl) GetLogEntry(ctx context.Context, index int64) (store.LogEntry, bool, error) {
	e, ok, err := r.readStore.GetLogEntry(ctx, index)
	if r.entryEdit != nil {
		r.entryEdit(&e)
	}
	if r.entryErr != nil {
		return e, ok, r.entryErr
	}
	if r.entryMissing {
		return e, false, nil
	}
	return e, ok, err
}
func (r *rosterReadControl) ObjectCommits(ctx context.Context, ref hashref.HashRef, after int64, limit int) ([]int64, error) {
	if r.commitMissing {
		return nil, nil
	}
	if r.commitMismatch {
		return []int64{0}, nil
	}
	return r.readStore.ObjectCommits(ctx, ref, after, limit)
}
func (r *rosterReadControl) ObjectsBySemanticID(context.Context, string, string, int) ([]store.Object, error) {
	panic("roster must never scan semantic IDs")
}

func flipRosterPayloadByte(o *store.Object) {
	// Flip one value byte while leaving the JSON structure intact.
	i := bytes.Index(o.Payload, []byte(`:"`)) + 2
	o.Payload[i] ^= 1
}

type rosterFixture struct {
	d                                             *Daemon
	packet                                        store.Object
	validated                                     coordinator.PacketResult
	input, record, effect, request, result, claim store.Object
	args                                          map[string]any
}

func rosterJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, e := json.Marshal(v)
	if e != nil {
		t.Fatal(e)
	}
	return b
}
func rosterObject(t *testing.T, d *Daemon, id string, payload []byte) store.Object {
	t.Helper()
	o := forgeObject(id, payload)
	o.Provenance = "demo-authored"
	commitObjects(t, d, "authored-claim", []store.Object{o}, o.Hash)
	return o
}
func rosterGet(t *testing.T, d *Daemon, ref hashref.HashRef) store.Object {
	t.Helper()
	o, ok, e := d.store.GetObject(boundedTestContext(t), ref)
	if e != nil || !ok {
		t.Fatalf("fixture read: %v %v", ok, e)
	}
	return o
}
func newRosterFixture(t *testing.T, change func(*rosterFixture)) *rosterFixture {
	t.Helper()
	d := newHandlerDaemon(t)
	return rosterFixtureOn(t, d, change)
}
func rosterFixtureOn(t *testing.T, d *Daemon, change func(*rosterFixture)) *rosterFixture {
	t.Helper()
	publishTools(t, d, hashref.SumSHA256([]byte("fake-interpreter")), "module host/capsule/main\n\nexport func main(input: string) -> string {\n input\n}\n",
		wsTool{ID: "hu.approve", Effect: broker.EffectHumanApprove}, wsTool{ID: "roster.test", Effect: broker.EffectAilangCLI})
	first, _ := producePacket(t, d, humanPlan("e1"), "evidence")
	rec := rosterGet(t, d, first.RecordRef)
	var wire forgedRecordV2
	if e := json.Unmarshal(rec.Payload, &wire); e != nil {
		t.Fatal(e)
	}
	effect := rosterGet(t, d, hashref.MustParse(wire.Effects[0]))
	er, e := broker.DecodeRecord(effect.Payload)
	if e != nil {
		t.Fatal(e)
	}
	claim := rosterObject(t, d, "world/demo/authored-claim/v1", []byte(`{"kind":"authored-claim","author":"demo-authored","text":"Demonstration commentary authored for this demo; not an AI review or proof."}`))
	// Both suites use a genuine coordinator/broker TESTED invocation.
	f := &rosterFixture{d: d, effect: effect, request: rosterGet(t, d, er.RequestRef), result: rosterGet(t, d, er.ResultRef), claim: claim}
	f.args = map[string]any{"goal": "demo goal", "candidates": []string{"a", "b"}, "roster": []any{
		map[string]any{"kind": "test-invocation", "ref": rosterTestInvocation(t, d, "").String()},
		map[string]any{"kind": "recorded-effect", "ref": effect.Hash.String()},
		map[string]any{"kind": "authored-claim", "ref": claim.Hash.String()},
		map[string]any{"kind": "proof-envelope", "ref": rec.Hash.String()}}}
	if change != nil {
		change(f)
	}
	call := humanCall(t, d, "ep1", "roster")
	call.Input = f.args
	c := humanCoordinator(t, d, testHumanBinder(d), humanPlan("e2"))
	res, e := c.Dispatch(boundedTestContext(t), call)
	if e != nil {
		t.Fatal(e)
	}
	f.packet = packetForEntry(t, d, res.EntryIndex)
	f.record = rosterGet(t, d, res.RecordRef)
	if e = json.Unmarshal(f.record.Payload, &wire); e != nil {
		t.Fatal(e)
	}
	f.input = rosterGet(t, d, hashref.MustParse(wire.Input))
	f.validated, e = coordinator.ValidatePacket(boundedTestContext(t), d.reads, f.packet)
	if e != nil || f.validated.Code != 0 {
		t.Fatalf("real packet validation: %+v %v", f.validated, e)
	}
	return f
}
func rosterPair(f *rosterFixture, i int) map[string]any {
	return f.args["roster"].([]any)[i].(map[string]any)
}
func replaceRosterObject(t *testing.T, f *rosterFixture, i int, o store.Object, edit func(*store.Object)) {
	t.Helper()
	o.Payload = append([]byte(nil), o.Payload...)
	edit(&o)
	o.Hash = hashref.SumSHA256(o.Payload)
	if e := f.d.store.PutObject(boundedTestContext(t), o); e != nil {
		t.Fatal(e)
	}
	rosterPair(f, i)["ref"] = o.Hash.String()
}
func padRosterInput(t *testing.T, f *rosterFixture, n int) {
	t.Helper()
	f.args["goal"] = "x"
	size := len(rosterJSON(t, f.args))
	f.args["goal"] = strings.Repeat("x", n-size+1)
	if len(rosterJSON(t, f.args)) != n {
		t.Fatal("input boundary fixture")
	}
}
func padRosterClaim(t *testing.T, f *rosterFixture, n int) {
	t.Helper()
	replaceRosterObject(t, f, 2, f.claim, func(o *store.Object) {
		var v map[string]any
		if e := json.Unmarshal(o.Payload, &v); e != nil {
			t.Fatal(e)
		}
		v["text"] = ""
		size := len(rosterJSON(t, v))
		v["text"] = strings.Repeat("x", n-size)
		o.Payload = rosterJSON(t, v)
	})
}

// Redirect only the record ref returned by the carrying-entry read to an actual
// stored object derived from the coordinator's record. Hash validation still runs.
func rosterRecordControl(t *testing.T, f *rosterFixture, r *rosterReadControl, edit func(*store.Object)) {
	t.Helper()
	o := f.record
	o.Payload = append([]byte(nil), o.Payload...)
	edit(&o)
	o.Hash = hashref.SumSHA256(o.Payload)
	if err := f.d.store.PutObject(boundedTestContext(t), o); err != nil {
		t.Fatal(err)
	}
	r.entryEdit = func(e *store.LogEntry) {
		if e.Header.EntryIndex == f.validated.Entry {
			e.TransitionRef = o.Hash
		}
	}
}
func padRosterEffectObject(t *testing.T, f *rosterFixture, result bool, n int) {
	t.Helper()
	target := f.request
	if result {
		target = f.result
	}
	target.Payload = append(append([]byte(nil), target.Payload...), []byte(strings.Repeat(" ", n-len(target.Payload)))...)
	target.Hash = hashref.SumSHA256(target.Payload)
	if err := f.d.store.PutObject(boundedTestContext(t), target); err != nil {
		t.Fatal(err)
	}
	replaceRosterObject(t, f, 1, f.effect, func(o *store.Object) {
		er, err := broker.DecodeRecord(o.Payload)
		if err != nil {
			t.Fatal(err)
		}
		if result {
			er.ResultRef = target.Hash
		} else {
			er.RequestRef = target.Hash
		}
		o.Payload = broker.EncodeRecord(er)
	})
}

func TestWorkbenchPacketRosterAnchoredBounds(t *testing.T) {
	cases := []struct {
		name, reason string
		change       func(*testing.T, *rosterFixture)
		read         func(*rosterFixture, hashref.HashRef, *store.Object)
		expired      bool
		prepare      func(*testing.T, *rosterFixture, *rosterReadControl)
		nonpacket    bool
		cancelParent bool
	}{
		{name: "real-attested-and-authored-claimed"},
		{name: "input-boundary", change: func(t *testing.T, f *rosterFixture) { padRosterInput(t, f, (64 << 10)) }},
		{name: "input-cap-plus-one", reason: "input size", change: func(t *testing.T, f *rosterFixture) { padRosterInput(t, f, (64<<10)+1) }},
		{name: "claim-boundary", change: func(t *testing.T, f *rosterFixture) { padRosterClaim(t, f, (4 << 10)) }},
		{name: "claim-cap-plus-one", reason: "claim size", change: func(t *testing.T, f *rosterFixture) { padRosterClaim(t, f, (4<<10)+1) }},
		{name: "read-boundary", change: func(t *testing.T, f *rosterFixture) {
			replaceRosterObject(t, f, 1, f.effect, func(o *store.Object) {
				o.Payload = append(o.Payload, []byte(strings.Repeat(" ", (1<<20)-len(o.Payload)))...)
			})
		}},
		{name: "read-cap-plus-one", reason: "effect size", change: func(t *testing.T, f *rosterFixture) {
			replaceRosterObject(t, f, 1, f.effect, func(o *store.Object) {
				o.Payload = append(o.Payload, []byte(strings.Repeat(" ", (1<<20)+1-len(o.Payload)))...)
			})
		}},
		{name: "result-boundary", change: func(t *testing.T, f *rosterFixture) { padRosterEffectObject(t, f, true, (1 << 20)) }},
		{name: "result-cap-plus-one", reason: "result size", change: func(t *testing.T, f *rosterFixture) { padRosterEffectObject(t, f, true, (1<<20)+1) }},
		{name: "request-boundary", change: func(t *testing.T, f *rosterFixture) { padRosterEffectObject(t, f, false, (1 << 20)) }},
		{name: "request-cap-plus-one", reason: "request size", change: func(t *testing.T, f *rosterFixture) { padRosterEffectObject(t, f, false, (1<<20)+1) }},
		{name: "no-roster-for-nonpacket", reason: "packet not validated", nonpacket: true},
		{name: "unrelated-input-ignored", change: func(t *testing.T, f *rosterFixture) {
			rosterObject(t, f.d, coordinator.InputV1, []byte(`{"grade":"PROVEN"}`))
		}},
		{name: "carrying-entry-missing", reason: "carrying entry missing", prepare: func(t *testing.T, f *rosterFixture, r *rosterReadControl) { r.entryMissing = true }},
		{name: "carrying-entry-read-error", reason: "carrying entry missing", prepare: func(t *testing.T, f *rosterFixture, r *rosterReadControl) {
			r.entryErr = errors.New("read unavailable")
		}},
		{name: "carrying-entry-deadline", reason: "read deadline", prepare: func(t *testing.T, f *rosterFixture, r *rosterReadControl) { r.entryErr = context.DeadlineExceeded }},
		{name: "record-missing", reason: "record missing", prepare: func(t *testing.T, f *rosterFixture, r *rosterReadControl) { r.missing = f.record.Hash }},
		{name: "input-missing", reason: "input missing", prepare: func(t *testing.T, f *rosterFixture, r *rosterReadControl) { r.missing = f.input.Hash }},
		{name: "effect-missing", reason: "effect missing", prepare: func(t *testing.T, f *rosterFixture, r *rosterReadControl) { r.missing = f.effect.Hash }},
		{name: "request-missing", reason: "request missing", prepare: func(t *testing.T, f *rosterFixture, r *rosterReadControl) { r.missing = f.request.Hash }},
		{name: "result-missing", reason: "result missing", prepare: func(t *testing.T, f *rosterFixture, r *rosterReadControl) { r.missing = f.result.Hash }},
		{name: "claim-missing", reason: "claim missing", prepare: func(t *testing.T, f *rosterFixture, r *rosterReadControl) { r.missing = f.claim.Hash }},
		{name: "record-schema", reason: "record schema", prepare: func(t *testing.T, f *rosterFixture, r *rosterReadControl) {
			rosterRecordControl(t, f, r, func(o *store.Object) { o.Payload[0] = '[' })
		}},
		{name: "record-input-ref", reason: "input ref", prepare: func(t *testing.T, f *rosterFixture, r *rosterReadControl) {
			rosterRecordControl(t, f, r, func(o *store.Object) {
				o.Payload = []byte(strings.Replace(string(o.Payload), `"input":"sha256:`, `"input":"bad:`, 1))
			})
		}},
		{name: "record-boundary", prepare: func(t *testing.T, f *rosterFixture, r *rosterReadControl) {
			rosterRecordControl(t, f, r, func(o *store.Object) {
				o.Payload = append(o.Payload, []byte(strings.Repeat(" ", (1<<20)-len(o.Payload)))...)
			})
		}},
		{name: "record-cap-plus-one", reason: "record size", prepare: func(t *testing.T, f *rosterFixture, r *rosterReadControl) {
			rosterRecordControl(t, f, r, func(o *store.Object) {
				o.Payload = append(o.Payload, []byte(strings.Repeat(" ", (1<<20)+1-len(o.Payload)))...)
			})
		}},
		{name: "record-kind", reason: "record kind", read: func(f *rosterFixture, ref hashref.HashRef, o *store.Object) {
			if ref == f.record.Hash {
				o.SemanticID = coordinator.RecordV1
			}
		}},
		{name: "request-kind", reason: "request kind", read: func(f *rosterFixture, ref hashref.HashRef, o *store.Object) {
			if ref == f.request.Hash {
				o.SemanticID = broker.EffectResultV1
			}
		}},
		{name: "result-kind", reason: "result kind", read: func(f *rosterFixture, ref hashref.HashRef, o *store.Object) {
			if ref == f.result.Hash {
				o.SemanticID = broker.EffectRequestV1
			}
		}},
		{name: "object-ref-mismatch", reason: "claim hash", read: func(f *rosterFixture, ref hashref.HashRef, o *store.Object) {
			if ref == f.claim.Hash {
				o.Hash = f.effect.Hash
			}
		}},

		{name: "count", reason: "roster count", change: func(t *testing.T, f *rosterFixture) { f.args["roster"] = f.args["roster"].([]any)[:3] }},
		{name: "duplicate-kind", reason: "duplicate kind", change: func(t *testing.T, f *rosterFixture) { rosterPair(f, 3)["kind"] = "authored-claim" }},
		{name: "duplicate-pair", reason: "duplicate kind", change: func(t *testing.T, f *rosterFixture) { f.args["roster"].([]any)[3] = rosterPair(f, 2) }},
		{name: "count-extra", reason: "roster count", change: func(t *testing.T, f *rosterFixture) {
			f.args["roster"] = append(f.args["roster"].([]any), rosterPair(f, 2))
		}},
		{name: "goal-missing", reason: "input schema", change: func(t *testing.T, f *rosterFixture) { delete(f.args, "goal") }},
		{name: "candidates-duplicate", reason: "input schema", change: func(t *testing.T, f *rosterFixture) { f.args["candidates"] = []string{"a", "a"} }},
		{name: "pair-ref-missing", reason: "ref syntax", change: func(t *testing.T, f *rosterFixture) { delete(rosterPair(f, 3), "ref") }},
		{name: "pair-kind-missing", reason: "unknown kind", change: func(t *testing.T, f *rosterFixture) { delete(rosterPair(f, 3), "kind") }},

		{name: "unknown-kind", reason: "unknown kind", change: func(t *testing.T, f *rosterFixture) { rosterPair(f, 3)["kind"] = "AiReview" }},
		{name: "caller-grade", reason: "input shape", change: func(t *testing.T, f *rosterFixture) { rosterPair(f, 1)["grade"] = "PROVEN" }},
		{name: "input-key-case", reason: "input shape", change: func(t *testing.T, f *rosterFixture) { f.args["Goal"] = f.args["goal"]; delete(f.args, "goal") }},
		{name: "pair-key-case", reason: "input shape", change: func(t *testing.T, f *rosterFixture) { p := rosterPair(f, 0); p["Kind"] = p["kind"]; delete(p, "kind") }},
		{name: "claim-key-case", reason: "claim shape", change: func(t *testing.T, f *rosterFixture) {
			replaceRosterObject(t, f, 2, f.claim, func(o *store.Object) {
				o.Payload = []byte(strings.Replace(string(o.Payload), `"author":`, `"Author":`, 1))
			})
		}},

		{name: "unknown-key", reason: "input shape", change: func(t *testing.T, f *rosterFixture) { f.args["extra"] = true }},
		{name: "schema", reason: "input schema", change: func(t *testing.T, f *rosterFixture) { f.args["candidates"] = []string{"a"} }},
		{name: "ref-syntax", reason: "ref syntax", change: func(t *testing.T, f *rosterFixture) { rosterPair(f, 3)["ref"] = "invalid" }},
		{name: "claim-other-declared-kind", reason: "claim kind", change: func(t *testing.T, f *rosterFixture) { rosterPair(f, 2)["ref"] = f.effect.Hash.String() }},
		{name: "claim-payload-kind", reason: "claim kind", change: func(t *testing.T, f *rosterFixture) {
			replaceRosterObject(t, f, 2, f.claim, func(o *store.Object) {
				o.Payload = []byte(strings.Replace(string(o.Payload), "authored-claim", "AiReview", 1))
			})
		}},
		{name: "claim-author", reason: "claim author", change: func(t *testing.T, f *rosterFixture) {
			replaceRosterObject(t, f, 2, f.claim, func(o *store.Object) {
				o.Payload = []byte(strings.Replace(string(o.Payload), "demo-authored", "AiReview", 1))
			})
		}},
		{name: "claim-producer", reason: "claim author", read: func(f *rosterFixture, r hashref.HashRef, o *store.Object) {
			if r == f.claim.Hash {
				o.Provenance = "AiReview"
			}
		}},
		{name: "claim-text-empty", reason: "claim shape", change: func(t *testing.T, f *rosterFixture) {
			replaceRosterObject(t, f, 2, f.claim, func(o *store.Object) {
				var v map[string]any
				if err := json.Unmarshal(o.Payload, &v); err != nil {
					t.Fatal(err)
				}
				v["text"] = ""
				o.Payload = rosterJSON(t, v)
			})
		}},
		{name: "claim-trailing-json", reason: "claim shape", change: func(t *testing.T, f *rosterFixture) {
			replaceRosterObject(t, f, 2, f.claim, func(o *store.Object) { o.Payload = append(o.Payload, []byte(" {}")...) })
		}},

		{name: "claim-shape", reason: "claim shape", change: func(t *testing.T, f *rosterFixture) {
			replaceRosterObject(t, f, 2, f.claim, func(o *store.Object) {
				o.Payload = []byte(strings.Replace(string(o.Payload), `"text":`, `"extra":`, 1))
			})
		}},
		{name: "record-consistency", reason: "effect inconsistent", change: func(t *testing.T, f *rosterFixture) {
			replaceRosterObject(t, f, 1, f.effect, func(o *store.Object) {
				er, e := broker.DecodeRecord(o.Payload)
				if e != nil {
					t.Fatal(e)
				}
				er.BudgetAfter++
				o.Payload = broker.EncodeRecord(er)
			})
		}},
		{name: "record-decode", reason: "effect decode", change: func(t *testing.T, f *rosterFixture) {
			replaceRosterObject(t, f, 1, f.effect, func(o *store.Object) {
				o.Payload = []byte(strings.Replace(string(o.Payload), `"requestRef":"sha256:`, `"requestRef":"bad:`, 1))
			})
		}},
		{name: "effect-kind", reason: "effect kind", change: func(t *testing.T, f *rosterFixture) { rosterPair(f, 1)["ref"] = f.claim.Hash.String() }},
		{name: "record-hash", reason: "effect hash", read: func(f *rosterFixture, r hashref.HashRef, o *store.Object) {
			if r == f.effect.Hash {
				flipRosterPayloadByte(o)
			}
		}},
		{name: "request-hash", reason: "request hash", read: func(f *rosterFixture, r hashref.HashRef, o *store.Object) {
			if r == f.request.Hash {
				flipRosterPayloadByte(o)
			}
		}},
		{name: "result-hash", reason: "result hash", read: func(f *rosterFixture, r hashref.HashRef, o *store.Object) {
			if r == f.result.Hash {
				flipRosterPayloadByte(o)
			}
		}},
		{name: "claim-hash", reason: "claim hash", read: func(f *rosterFixture, r hashref.HashRef, o *store.Object) {
			if r == f.claim.Hash {
				flipRosterPayloadByte(o)
			}
		}},
		{name: "input-hash", reason: "input hash", read: func(f *rosterFixture, r hashref.HashRef, o *store.Object) {
			if r == f.input.Hash {
				flipRosterPayloadByte(o)
			}
		}},
		{name: "carrying-record-hash", reason: "record hash", read: func(f *rosterFixture, r hashref.HashRef, o *store.Object) {
			if r == f.record.Hash {
				flipRosterPayloadByte(o)
			}
		}},
		{name: "input-semantic-id", reason: "input kind", read: func(f *rosterFixture, r hashref.HashRef, o *store.Object) {
			if r == f.input.Hash {
				o.SemanticID = coordinator.OutputV1
			}
		}},
		{name: "input-interface", reason: "input kind", read: func(f *rosterFixture, r hashref.HashRef, o *store.Object) {
			if r == f.input.Hash {
				o.InterfaceHash = hashref.SumSHA256([]byte("wrong"))
			}
		}},
		{name: "duplicate-json-key", reason: "input shape", change: func(t *testing.T, f *rosterFixture) {
			raw := rosterJSON(t, rosterPair(f, 0))
			raw = []byte(strings.Replace(string(raw), `"kind":`, `"kind":"test-invocation","kind":`, 1))
			f.args["roster"].([]any)[0] = json.RawMessage(raw)
		}},
		{name: "shared-deadline-overflow", reason: "read deadline", prepare: func(t *testing.T, f *rosterFixture, r *rosterReadControl) { r.waitAt = f.claim.Hash }},
		{name: "parent-deadline", reason: "read deadline", cancelParent: true},

		{name: "deadline", reason: "read deadline", expired: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Each row owns all its contexts and its daemon. Never share the parent's 30s budget.
			ctx := boundedTestContext(t)
			f := newRosterFixture(t, func(f *rosterFixture) {
				if tc.change != nil {
					tc.change(t, f)
				}
			})
			r := &rosterReadControl{readStore: f.d.reads}
			if tc.read != nil {
				r.edit = func(ref hashref.HashRef, o *store.Object) { tc.read(f, ref, o) }
			}
			if tc.expired {
				r.readErr = context.DeadlineExceeded
			}
			if tc.prepare != nil {
				tc.prepare(t, f, r)
			}
			if tc.nonpacket {
				var err error
				f.validated, err = coordinator.ValidatePacket(ctx, f.d.reads, f.input)
				if err != nil || f.validated.Code == 0 {
					t.Fatalf("nonpacket validation: %+v %v", f.validated, err)
				}
				projected, err := f.d.projectPacket(ctx, f.input)
				if err != nil || projected.view.Validated || len(projected.roster.rows) != 0 {
					t.Fatalf("nonpacket decorated: %+v %v", projected, err)
				}
			}
			if tc.cancelParent {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}
			got := packetRosterFromValidated(ctx, r, f.validated, nil)
			if tc.reason != "" {
				if got.unavailable != tc.reason {
					t.Fatalf("UNAVAILABLE guard %s: got %+v", tc.reason, got)
				}
				if len(got.rows) != 0 {
					t.Fatalf("refused roster carries labels: %+v", got)
				}
				return
			}
			if got.unavailable != "" || len(got.rows) != 4 {
				t.Fatalf("positive roster: %+v", got)
			}
			if tc.name == "real-attested-and-authored-claimed" {
				projected, err := f.d.projectPacket(ctx, f.packet)
				if err != nil || !projected.view.Validated || len(projected.roster.rows) != 4 {
					t.Fatalf("packet projection lost roster: %+v %v", projected, err)
				}
			}
			if got.rows[1].label != "ATTESTED" || got.rows[2].label != "CLAIMED" || got.rows[2].producer != "demo-authored" {
				t.Fatalf("real evidence labels: %+v", got)
			}
			if got.rows[0].label != "TESTED" {
				t.Fatalf("TESTED missing: %+v", got)
			}
			for _, i := range []int{3} {
				if got.rows[i].label != "" || got.rows[i].unavailable == "" {
					t.Fatalf("pending row gained grade: %+v", got.rows[i])
				}
			}
			if len(r.deadlines) == 0 {
				t.Fatal("no anchored reads")
			}
			for _, deadline := range r.deadlines {
				if !deadline.Equal(r.deadlines[0]) || time.Until(deadline) > (10*time.Second) {
					t.Fatal("reads do not share the 10s deadline")
				}
			}
		})
	}
}

// Every stdout is a checked-in real report, or the single named mutation supplied
// by the table below. The fake only transports bytes; broker and coordinator
// construct all request/result/effect/record hashes and the carrying commit.
func rosterTestInvocation(t *testing.T, d *Daemon, stdout string) hashref.HashRef {
	t.Helper()
	if stdout == "" {
		raw, err := os.ReadFile("testdata/roster_test_report_pass.json")
		if err != nil {
			t.Fatal(err)
		}
		stdout = string(raw)
	}
	root, logs := canonicalTemp(t), canonicalTemp(t)
	out := string(rosterJSON(t, map[string]any{"ok": true, "exit_code": 0, "argv": []string{"test", "--json", "world/logepoch.ail"}, "stdout": stdout, "stderr": ""}))
	bin := fakeToolBinWith(t, logs, out)
	raw, err := os.ReadFile(bin)
	if err != nil {
		t.Fatal(err)
	}
	raw = bytes.ReplaceAll(raw, []byte(`"cli":["check"]`), []byte(`"cli":["check","test"]`))
	if err = os.WriteFile(bin, raw, 0700); err != nil {
		t.Fatal(err)
	}
	handler, err := broker.NewAilangToolHandler(boundedTestContext(t), broker.AilangToolConfig{Bin: bin, BinRef: hashref.SumSHA256(raw), Root: root, PolicyPath: filepath.Join(logs, "policy.json"), CacheDir: filepath.Join(logs, "cache")})
	if err != nil {
		t.Fatal(err)
	}
	grants := []broker.Capability{{Effect: broker.EffectAilangCLI, Scope: broker.WorkspaceScope, ExpiresAt: packetExpiry, Budget: 5}}
	req, err := transitionreg.NewRequest(boundedTestContext(t), transitionreg.NewReader(d.store), wsCaps(grants), packetNow)
	if err != nil {
		t.Fatal(err)
	}
	binder := func(episode string, caps []broker.Capability) transitionreg.Binder {
		return broker.OpenBinder(d.store, episode, caps, broker.Registry{broker.EffectAilangCLI: handler})
	}
	plan := `{"plan":"world/effect-plan/v1","effects":[{"id":"test","effect":"Ailang.CLI","scope":"worktree","cost":1,"payload":{"op":"test","path":"world/logepoch.ail","flags":{"json":""}}}],"finish":false,"result":null}`
	res, err := humanCoordinator(t, d, binder, plan).Dispatch(boundedTestContext(t), coordinator.Call{Surface: coordinator.SurfaceA2A, Request: req, EpisodeID: "ep1", Grants: grants, SkillID: "roster.test", TaskID: "roster-test", Input: map[string]any{}})
	if err != nil {
		t.Fatal(err)
	}
	return res.RecordRef
}

// Rebuild one edge at a time, retaining genuine stored objects and real hashes.
func rosterTestChain(t *testing.T, f *rosterFixture, target string, edit func(*store.Object)) hashref.HashRef {
	t.Helper()
	rec := rosterGet(t, f.d, hashref.MustParse(rosterPair(f, 0)["ref"].(string)))
	var wire packetTestRecord
	if err := json.Unmarshal(rec.Payload, &wire); err != nil {
		t.Fatal(err)
	}
	effect := rosterGet(t, f.d, hashref.MustParse(wire.Effects[0]))
	er, err := broker.DecodeRecord(effect.Payload)
	if err != nil {
		t.Fatal(err)
	}
	put := func(o store.Object) store.Object {
		o.Hash = hashref.SumSHA256(o.Payload)
		if err := f.d.store.PutObject(boundedTestContext(t), o); err != nil {
			t.Fatal(err)
		}
		return o
	}
	if target == "result" || target == "request" {
		ref := er.ResultRef
		if target == "request" {
			ref = er.RequestRef
		}
		obj := rosterGet(t, f.d, ref)
		edit(&obj)
		obj = put(obj)
		if target == "request" {
			er.RequestRef = obj.Hash
		} else {
			er.ResultRef = obj.Hash
		}
		effect.Payload = broker.EncodeRecord(er)
	} else if target == "effect" {
		edit(&effect)
	}
	if target != "record" {
		effect = put(effect)
		wire.Effects = []string{effect.Hash.String()}
		rec.Payload = rosterJSON(t, wire)
	} else {
		edit(&rec)
	}
	rec = put(rec)
	// A real store commit makes the replacement record its TransitionRef.
	commitObjects(t, f.d, "test-chain-control", []store.Object{rec}, rec.Hash)
	rosterPair(f, 0)["ref"] = rec.Hash.String()
	return rec.Hash
}

func TestWorkbenchPacketRosterTestedAndProven(t *testing.T) {
	type control struct {
		name, reason string
		mutate       func(*testing.T, *rosterFixture, *gradeRig, *rosterReadControl)
		report       func(map[string]any)
		proof        bool
	}
	cases := []control{
		{name: "four-real-rows"},
		{name: "carrying-missing", reason: "test carrying commit missing", mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) { r.commitMissing = true }},
		{name: "carrying-mismatch", reason: "test carrying commit mismatch", mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) { r.commitMismatch = true }},
		{name: "effect-count", reason: "test effect count", mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			rosterTestChain(t, f, "record", func(o *store.Object) {
				var v packetTestRecord
				_ = json.Unmarshal(o.Payload, &v)
				v.Effects = append(v.Effects, v.Effects[0])
				o.Payload = rosterJSON(t, v)
			})
		}},
		{name: "effect-not-cli", reason: "test effect not Ailang.CLI", mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			rosterTestChain(t, f, "effect", func(o *store.Object) {
				er, _ := broker.DecodeRecord(o.Payload)
				er.Effect = broker.EffectAilangCheck
				o.Payload = broker.EncodeRecord(er)
			})
		}},
		{name: "op-not-test", reason: "test op", mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			rosterTestChain(t, f, "request", func(o *store.Object) {
				o.Payload = bytes.Replace(o.Payload, []byte(`"op":"test"`), []byte(`"op":"check"`), 1)
			})
		}},
	}
	for _, c := range []struct {
		name, key string
		value     any
		reason    string
	}{
		{"exit-code", "exit_code", 1, "test exit code"}, {"not-ok", "ok", false, "test not ok"}, {"argv", "argv", []string{"check"}, "test argv"},
		// Prose/null are single whole-stdout substitutions of the checked-in pass fixture.
		{"prose", "stdout", "ordinary prose", "test report JSON"}, {"null", "stdout", "null", "test report null"},
	} {
		c := c
		cases = append(cases, control{name: c.name, reason: c.reason, mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			rosterTestChain(t, f, "result", func(o *store.Object) {
				var v map[string]any
				_ = json.Unmarshal(o.Payload, &v)
				v[c.key] = c.value
				o.Payload = rosterJSON(t, v)
			})
		}})
	}
	for _, c := range []struct {
		name, key string
		value     any
		reason    string
	}{
		{"success-false", "success", false, "test success"}, {"failed", "failed_tests", 1, "test failed"}, {"skipped", "skipped_tests", 1, "test skipped"}, {"vacuous", "vacuous_skips", 1, "test vacuous"}, {"passed-zero", "passed_tests", 0, "test zero passes"}, {"passed-mismatch", "passed_tests", 11, "test passed count"},
	} {
		c := c
		cases = append(cases, control{name: c.name, reason: c.reason, report: func(v map[string]any) { v[c.key] = c.value }})
	}
	cases = append(cases,
		control{name: "named-over-cap", reason: "test named count", report: func(v map[string]any) {
			tests := v["tests"].([]any)
			for len(tests) < 255 { // 255 tests + the report's 2 properties = 257 named = cap+1

				tests = append(tests, tests[0])
			}
			v["tests"] = tests
		}},
		control{name: "named-empty", reason: "test named count", report: func(v map[string]any) { v["tests"] = []any{}; v["properties"] = []any{} }},
		control{name: "non-pass", reason: "test status", report: func(v map[string]any) { v["tests"].([]any)[0].(map[string]any)["status"] = "fail" }},
		control{name: "empty-name", reason: "test empty name", report: func(v map[string]any) { v["tests"].([]any)[0].(map[string]any)["name"] = "" }},
		control{name: "missing-counter", reason: "test report counters", report: func(v map[string]any) { delete(v, "failed_tests") }},
		control{name: "named-schema", reason: "test report JSON", report: func(v map[string]any) { v["tests"] = 7 }},
		control{name: "vacuous-real", reason: "test success", mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			raw, e := os.ReadFile("testdata/roster_test_report_vacuous.json")
			if e != nil {
				t.Fatal(e)
			}
			rosterTestChain(t, f, "result", func(o *store.Object) {
				var v map[string]any
				_ = json.Unmarshal(o.Payload, &v)
				v["stdout"] = string(raw)
				o.Payload = rosterJSON(t, v)
			})
		}},
		control{name: "no-grader", reason: objectGradeUnavailableReason, proof: true, mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) { f.d.grader = nil }},
		control{name: "wrong-key", reason: gradeReasonUnauthenticated, proof: true, mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			g.key[0] ^= 1
			ref := storeGradeEnvelope(t, f.d, mintGradeEnvelope(t, g.key, gradeReport(*g, rosterProofSubject(t, f))))
			rosterPair(f, 3)["ref"] = ref.String()
		}},
		control{name: "edited-subject", reason: gradeReasonUnauthenticated, proof: true, mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			obj := rosterGet(t, f.d, hashref.MustParse(rosterPair(f, 3)["ref"].(string)))
			env, e := evidence.DecodeAuthenticatedEnvelope(obj.Payload)
			if e != nil {
				t.Fatal(e)
			}
			rep, _ := evidence.DecodeProofReportV1(env.Report)
			other := rosterObject(t, f.d, "world/source/ailang/v1", []byte("different source"))
			rep.Subject = other.Hash
			env.Report, _ = evidence.EncodeProofReportV1(rep)
			raw, e := evidence.EncodeAuthenticatedEnvelope(env)
			if e != nil {
				t.Fatal(e)
			}
			rosterPair(f, 3)["ref"] = storeGradeEnvelope(t, f.d, raw).String()
		}},
		control{name: "subject-kind", reason: "proof subject kind", proof: true, mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			rep := gradeReport(*g, f.claim.Hash)
			rosterPair(f, 3)["ref"] = storeGradeEnvelope(t, f.d, mintGradeEnvelope(t, g.key, rep)).String()
		}},
		control{name: "promote-claim", reason: "unknown kind", mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			rosterPair(f, 2)["kind"] = "PROVEN"
		}},
	)
	for _, c := range []struct {
		name, reason string
		edit         func(*evidence.ProofReportV1)
	}{
		{"tool-mismatch", gradeReasonToolMismatch, func(p *evidence.ProofReportV1) { p.CompilerVersion = "other" }},
		{"proof-failed", gradeReasonProofFailed, func(p *evidence.ProofReportV1) { p.CheckPassed = false }},
		{"proof-incomplete", gradeReasonProofIncomplete, func(p *evidence.ProofReportV1) { p.Verified = []string{"different"} }},
	} {
		c := c
		cases = append(cases, control{name: c.name, reason: c.reason, proof: true, mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			rep := gradeReport(*g, rosterProofSubject(t, f))
			c.edit(&rep)
			rosterPair(f, 3)["ref"] = storeGradeEnvelope(t, f.d, mintGradeEnvelope(t, g.key, rep)).String()
		}})
	}

	cases = append(cases,
		control{name: "stdout-over-cap", reason: "test stdout size", mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			rosterTestChain(t, f, "result", func(o *store.Object) {
				var v map[string]any
				_ = json.Unmarshal(o.Payload, &v)
				raw := v["stdout"].(string)
				v["stdout"] = raw + strings.Repeat(" ", packetRosterStdoutLimit+1-len(raw))
				o.Payload = rosterJSON(t, v)
			})
		}},
		control{name: "subject-mismatch", reason: gradeReasonNoneValidated, proof: true, mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			f.d.grader = rosterMismatchedSubject{proofValidator: f.d.grader, subject: f.claim.Hash}
		}},
		control{name: "remac-edited-subject", reason: gradeReasonUnauthenticated, proof: true, mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			other := rosterObject(t, f.d, "world/source/ailang/v1", []byte("other source remac"))
			rep := gradeReport(*g, other.Hash)
			g.key[0] ^= 1
			rosterPair(f, 3)["ref"] = storeGradeEnvelope(t, f.d, mintGradeEnvelope(t, g.key, rep)).String()
		}},
		control{name: "record-schema", reason: "test record schema", mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			rosterTestChain(t, f, "record", func(o *store.Object) { o.Payload = append(o.Payload, []byte(" null")...) })
		}},
		control{name: "effect-ref", reason: "test effect ref", mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			rosterTestChain(t, f, "record", func(o *store.Object) {
				var v packetTestRecord
				_ = json.Unmarshal(o.Payload, &v)
				v.Effects = []string{"bad-ref"}
				o.Payload = rosterJSON(t, v)
			})
		}},
		control{name: "effect-inconsistent", reason: "test effect inconsistent", mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			rosterTestChain(t, f, "effect", func(o *store.Object) {
				er, _ := broker.DecodeRecord(o.Payload)
				er.BudgetAfter++
				o.Payload = broker.EncodeRecord(er)
			})
		}},
		control{name: "request-schema", reason: "test request schema", mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			rosterTestChain(t, f, "request", func(o *store.Object) { o.Payload = append(o.Payload, []byte(" null")...) })
		}},
		control{name: "result-schema", reason: "test result schema", mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			rosterTestChain(t, f, "result", func(o *store.Object) { o.Payload = append(o.Payload, []byte(" null")...) })
		}},
		control{name: "stdout-schema", reason: "test stdout schema", mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			rosterTestChain(t, f, "result", func(o *store.Object) {
				var v map[string]any
				_ = json.Unmarshal(o.Payload, &v)
				v["stdout"] = 7
				o.Payload = rosterJSON(t, v)
			})
		}},
		control{name: "proof-schema", reason: gradeReasonStructuralCandidate, proof: true, mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			rosterPair(f, 3)["ref"] = storeGradeEnvelope(t, f.d, []byte("null")).String()
		}},
		control{name: "resolve-refused", reason: gradeReasonNoneValidated, proof: true, mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			f.d.grader = rosterForeignResolver{proofValidator: f.d.grader}
		}},
		control{name: "proof-deadline", reason: "read deadline", mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
			r.waitAt = rosterProofSubject(t, f)
		}},
	)

	for _, role := range []string{"test record", "test effect", "test request", "test result", "proof envelope", "proof subject"} {
		for _, failure := range []string{"missing", "hash", "kind", "size"} {
			role, failure := role, failure
			cases = append(cases, control{name: strings.ReplaceAll(role, " ", "-") + "-" + failure, reason: role + " " + failure, proof: strings.HasPrefix(role, "proof"), mutate: func(t *testing.T, f *rosterFixture, g *gradeRig, r *rosterReadControl) {
				ref := hashref.MustParse(rosterPair(f, 0)["ref"].(string))
				if strings.HasPrefix(role, "proof") {
					ref = hashref.MustParse(rosterPair(f, 3)["ref"].(string))
					if role == "proof subject" {
						ref = rosterProofSubject(t, f)
					}
				} else if role != "test record" {
					rec := rosterGet(t, f.d, ref)
					var v packetTestRecord
					_ = json.Unmarshal(rec.Payload, &v)
					ref = hashref.MustParse(v.Effects[0])
					if role != "test effect" {
						er, _ := broker.DecodeRecord(rosterGet(t, f.d, ref).Payload)
						ref = er.RequestRef
						if role == "test result" {
							ref = er.ResultRef
						}
					}
				}
				if failure == "missing" {
					r.missing = ref
					return
				}
				r.edit = func(got hashref.HashRef, o *store.Object) {
					if got != ref {
						return
					}
					switch failure {
					case "hash":
						o.Payload = append(o.Payload, ' ')
					case "kind":
						o.SemanticID = "other"
					case "size":
						limit := packetRosterReadLimit
						if role == "proof envelope" {
							limit = evidence.MaxBytes
						}
						o.Payload = bytes.Repeat([]byte{' '}, limit+1)
					}
				}
			}})
		}
	}
	cases = append(cases, control{name: "named-boundary", report: func(v map[string]any) {
		tests := v["tests"].([]any)
		for len(tests) < 254 {
			tests = append(tests, tests[0])
		}
		v["tests"] = tests
		v["passed_tests"] = 256
	}})
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := boundedTestContext(t)
			d, rig := newGradeDaemon(t, []string{gradeTestIdentity})
			r := &rosterReadControl{readStore: d.reads}
			f := rosterFixtureOn(t, d, func(f *rosterFixture) {
				subject := rosterObject(t, d, "world/source/ailang/v1", []byte("module roster/source"))
				rosterPair(f, 3)["ref"] = storeGradeEnvelope(t, d, mintGradeEnvelope(t, rig.key, gradeReport(rig, subject.Hash))).String()
				if tc.report != nil {
					rosterTestChain(t, f, "result", func(o *store.Object) {
						var out map[string]any
						_ = json.Unmarshal(o.Payload, &out)
						var rep map[string]any
						_ = json.Unmarshal([]byte(out["stdout"].(string)), &rep)
						tc.report(rep)
						out["stdout"] = string(rosterJSON(t, rep))
						o.Payload = rosterJSON(t, out)
					})
				}
				if tc.mutate != nil {
					tc.mutate(t, f, &rig, r)
				}
			})
			var got packetRoster
			if tc.name == "stdout-over-cap" {
				// The outer 1MiB result cap necessarily fires first in the full projection.
				// Isolate the decoded-stdout guard with genuine stored hashes and a read seam.
				read := func(ref hashref.HashRef, role, id string, limit int) (store.Object, string) {
					o := rosterGet(t, d, ref)
					if o.Hash != hashref.SumSHA256(o.Payload) || o.SemanticID != id {
						t.Fatal("invalid isolated read")
					}
					return o, ""
				}
				_, reason := packetTestedRow(ctx, r, read, hashref.MustParse(rosterPair(f, 0)["ref"].(string)))
				got.unavailable = reason
			} else {
				got = packetRosterFromValidated(ctx, r, f.validated, d.grader)
			}
			if tc.reason != "" {
				reason := got.unavailable
				if tc.proof && len(got.rows) == 4 {
					reason = got.rows[3].unavailable
					if got.rows[3].label != "" {
						t.Fatal("refused proof promoted")
					}
				}
				if reason != tc.reason {
					t.Fatalf("refusal guard %s: got %+v", tc.reason, got)
				}
				// A stored-report refusal reaches the page through the real projection.
				if tc.report != nil && !tc.proof && tc.name != "stdout-over-cap" {
					projected, err := d.projectPacket(ctx, f.packet)
					if err != nil || projected.view.RosterUnavailable != tc.reason || len(projected.view.Roster) != 0 {
						t.Fatalf("refused projection %s: %+v %v", tc.reason, projected, err)
					}
				}
				return
			}
			if got.unavailable != "" || len(got.rows) != 4 {
				t.Fatalf("four rows: %+v", got)
			}
			for i, label := range []string{"TESTED", "ATTESTED", "CLAIMED", "PROVEN"} {
				if got.rows[i].label != label || got.rows[i].producer == "" {
					t.Fatalf("resolved row %d: %+v", i, got.rows[i])
				}
			}
			wantPasses := 10
			if tc.name == "named-boundary" {
				wantPasses = 256
			}
			if got.rows[0].passed != wantPasses || got.rows[0].verdict != "pass" || got.rows[3].subject == "" || got.rows[3].version != rig.version || got.rows[3].compiler != rig.compiler.String() {
				t.Fatalf("verdict/pins: %+v", got)
			}
			// Producers are the stored objects' own provenance, not merely non-empty.
			wantProducers := rosterWantProducers(t, f)
			for i := range got.rows {
				if got.rows[i].producer != wantProducers[i] {
					t.Fatalf("row %d producer %q, want %q", i, got.rows[i].producer, wantProducers[i])
				}
			}
			projected, err := d.projectPacket(ctx, f.packet)
			if err != nil || len(projected.view.Roster) != 4 {
				t.Fatalf("render projection: %+v %v", projected, err)
			}
			rosterAssertView(t, projected.view, got)
			code, body := gradePage(t, d, f.packet.Hash)
			if code != http.StatusOK {
				t.Fatalf("packet page: status %d", code)
			}
			for _, want := range []string{
				`aria-label="evidence roster"`,
				"test-invocation: TESTED · producer/author: " + wantProducers[0] + " · report:",
				"recorded-effect: ATTESTED · producer/author: " + wantProducers[1] + " · report:",
				"authored-claim: CLAIMED · producer/author: demo-authored · demo-authored · report:",
				"proof-envelope: PROVEN · producer/author: " + wantProducers[3] + " · report:",
				fmt.Sprintf("verdict pass · %d named passes", wantPasses),
				`subject: <a href="/workbench?object=` + got.rows[3].subject + `">`,
				"version " + rig.version,
			} {
				if !strings.Contains(body, want) {
					t.Fatalf("packet page lacks %q", want)
				}
			}
		})
	}
}
func rosterProofSubject(t *testing.T, f *rosterFixture) hashref.HashRef {
	t.Helper()
	obj := rosterGet(t, f.d, hashref.MustParse(rosterPair(f, 3)["ref"].(string)))
	env, e := evidence.DecodeAuthenticatedEnvelope(obj.Payload)
	if e != nil {
		t.Fatal(e)
	}
	rep, _ := evidence.DecodeProofReportV1(env.Report)
	return rep.Subject
}

// Both seams exercise real validator results; neither can mint authority.
type rosterMismatchedSubject struct {
	proofValidator
	subject hashref.HashRef
}

func (r rosterMismatchedSubject) ValidateProof(ctx context.Context, ref, _ hashref.HashRef) evidence.ValidationResult {
	return r.proofValidator.ValidateProof(ctx, ref, r.subject)
}

type rosterForeignResolver struct{ proofValidator }

func (r rosterForeignResolver) Resolve(evidence.ValidatedEvidence) evidence.ResolutionResult {
	return r.proofValidator.Resolve(evidence.ValidatedEvidence{})
}

// rosterWantProducers reads each roster row's producer from the stored object
// that row is projected from: the test result, the recorded effect, the authored
// claim and the proof envelope.
func rosterWantProducers(t *testing.T, f *rosterFixture) [4]string {
	t.Helper()
	rec := rosterGet(t, f.d, hashref.MustParse(rosterPair(f, 0)["ref"].(string)))
	var wire packetTestRecord
	if err := json.Unmarshal(rec.Payload, &wire); err != nil {
		t.Fatal(err)
	}
	er, err := broker.DecodeRecord(rosterGet(t, f.d, hashref.MustParse(wire.Effects[0])).Payload)
	if err != nil {
		t.Fatal(err)
	}
	want := [4]string{
		rosterGet(t, f.d, er.ResultRef).Provenance,
		rosterGet(t, f.d, hashref.MustParse(rosterPair(f, 1)["ref"].(string))).Provenance,
		"demo-authored",
		rosterGet(t, f.d, hashref.MustParse(rosterPair(f, 3)["ref"].(string))).Provenance,
	}
	for i, p := range want {
		if p == "" || p == "unavailable" {
			t.Fatalf("instrument: row %d's stored provenance %q cannot discriminate a dropped producer", i, p)
		}
	}
	return want
}

// rosterAssertView requires the daemon's workbench rows to carry every resolved
// field, the CLAIMED caption on the authored claim alone, and links built from
// each row's own refs.
func rosterAssertView(t *testing.T, v *workbench.PacketView, got packetRoster) {
	t.Helper()
	if v.RosterUnavailable != "" || len(v.Roster) != len(got.rows) {
		t.Fatalf("view roster: unavailable %q, %d rows", v.RosterUnavailable, len(v.Roster))
	}
	for i, row := range got.rows {
		caption := ""
		if row.kind == "authored-claim" {
			caption = "demo-authored"
		}
		want := workbench.PacketEvidenceView{Kind: row.kind, Ref: row.ref, Href: "?object=" + row.ref, Label: row.label, Producer: row.producer,
			Caption: caption, Unavailable: row.unavailable, Subject: row.subject, SubjectHref: "?object=" + row.subject,
			Compiler: row.compiler, CompilerHref: "?object=" + row.compiler, Version: row.version, Verdict: row.verdict, Passed: row.passed}
		if v.Roster[i] != want {
			t.Fatalf("view row %d:\n got %+v\nwant %+v", i, v.Roster[i], want)
		}
	}
}
