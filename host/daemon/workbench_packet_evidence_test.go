package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// Forward real reads; controls alter one field after the genuine producer/store.
type rosterReadControl struct {
	readStore
	edit         func(hashref.HashRef, *store.Object)
	readErr      error
	missing      hashref.HashRef
	entryEdit    func(*store.LogEntry)
	entryErr     error
	entryMissing bool
	waitAt       hashref.HashRef
	deadlines    []time.Time
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
	publishHumanTool(t, d)
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
	// M5a does not read TESTED or PROVEN. These genuine refs are placeholders for M5b.
	f := &rosterFixture{d: d, effect: effect, request: rosterGet(t, d, er.RequestRef), result: rosterGet(t, d, er.ResultRef), claim: claim}
	f.args = map[string]any{"goal": "demo goal", "candidates": []string{"a", "b"}, "roster": []any{
		map[string]any{"kind": "test-invocation", "ref": rec.Hash.String()},
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
	r.entryEdit = func(e *store.LogEntry) { e.TransitionRef = o.Hash }
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
			got := packetRosterFromValidated(ctx, r, f.validated)
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
			for _, i := range []int{0, 3} {
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
