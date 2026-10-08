package broker

import (
	"context"
	"errors"
	"fmt"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"reflect"
	"testing"
)

type approvalViewReader struct {
	objects                             map[hashref.HashRef]store.Object
	head                                hashref.HashRef
	registryReads, headReads, bodyReads int
	failRegistry, failObject            bool
}

func (r *approvalViewReader) GetRegistryHead(context.Context, string) (hashref.HashRef, bool, error) {
	r.registryReads++
	if r.failRegistry {
		return hashref.HashRef{}, false, errors.New("SECRET registry I/O")
	}
	return r.head, !r.head.IsZero(), nil
}
func (r *approvalViewReader) GetObject(_ context.Context, ref hashref.HashRef) (store.Object, bool, error) {
	if r.failObject {
		return store.Object{}, false, errors.New("SECRET object I/O")
	}
	o, ok := r.objects[ref]
	if o.SemanticID == ApprovalsV1 {
		r.headReads++
	} else {
		r.bodyReads++
	}
	return o, ok, nil
}
func (r *approvalViewReader) put(id string, wire any) hashref.HashRef {
	o := brokerObject(id, mustApprovalJSON(wire))
	r.objects[o.Hash] = o
	return o.Hash
}
func approvalViewFixture(n int, duplicate bool) *approvalViewReader {
	r := &approvalViewReader{objects: map[hashref.HashRef]store.Object{}}
	for i := 0; i < n; i++ {
		key := i
		if duplicate {
			key = i / 2
		}
		req := r.put(ApprovalRequestV1, approvalRequestWire{Effect: EffectHumanApprove, Scope: fmt.Sprintf("scope-%d", key), Requester: "requester", Cost: 3, Now: int64(key)})
		dec := r.put(ApprovalDecisionV1, approvalDecisionWire{RequestRef: req.String(), Decision: "approve", DecidedBy: fmt.Sprintf("operator-%d", i), Now: int64(i)})
		r.head = r.put(ApprovalsV1, approvalHeadWire{PreviousHead: r.head.String(), RequestRef: req.String(), DecisionRef: dec.String()})
	}
	return r
}
func TestFoldApprovals(t *testing.T) {
	heads := []ApprovalHead{{RequestRef: "C", DecisionRef: "C-approve"}, {RequestRef: "C"}, {RequestRef: "A", DecisionRef: "A-deny"}, {RequestRef: "B"}, {RequestRef: "A"}}
	got := FoldApprovals(heads)
	var refs, decisions []string
	for _, s := range got {
		refs = append(refs, s.RequestRef)
		decisions = append(decisions, s.DecisionRef)
	}
	if !reflect.DeepEqual(refs, []string{"C", "A", "B"}) || !reflect.DeepEqual(decisions, []string{"C-approve", "A-deny", ""}) || got[2].Status != "pending" {
		t.Fatalf("fold order/decisions: %+v", got)
	}
	heads = append([]ApprovalHead{{RequestRef: "A", DecisionRef: "A-new-approve"}}, heads...)
	got = FoldApprovals(heads)
	if got[0].RequestRef != "A" || got[0].DecisionRef != "A-new-approve" {
		t.Fatalf("newest decision lost: %+v", got)
	}

	// Producer-shaped heads independently pin [C approve, A deny, B pending].
	r := &approvalViewReader{objects: map[hashref.HashRef]store.Object{}}
	refsByName := map[string]hashref.HashRef{}
	for i, name := range []string{"A", "B", "C"} {
		refsByName[name] = r.put(ApprovalRequestV1, approvalRequestWire{Effect: EffectHumanApprove, Scope: name, Requester: "requester", Cost: 3, Now: int64(100 - i)})
	}
	appendHead := func(name, decision string) {
		dec := hashref.HashRef{}
		if decision != "" {
			dec = r.put(ApprovalDecisionV1, approvalDecisionWire{RequestRef: refsByName[name].String(), Decision: decision, DecidedBy: "operator", Now: 1})
		}
		r.head = r.put(ApprovalsV1, approvalHeadWire{PreviousHead: r.head.String(), RequestRef: refsByName[name].String(), DecisionRef: dec.String()})
	}
	appendHead("A", "")
	appendHead("B", "")
	appendHead("A", "deny")
	appendHead("C", "")
	appendHead("C", "approve")
	page, err := RecentApprovals(boundedTestContext(t), r, 40)
	if err != nil || len(page.Summaries) != 3 {
		t.Fatalf("resolved fold: %+v %v", page, err)
	}
	for i, want := range []struct{ name, status string }{{"C", "approve"}, {"A", "deny"}, {"B", "pending"}} {
		row := page.Summaries[i]
		if row.RequestRef != refsByName[want.name].String() || row.Status != want.status || row.Scope != want.name || row.Effect != EffectHumanApprove || row.Cost != 3 {
			t.Fatalf("resolved summary=%+v want %+v", row, want)
		}
	}
	appendHead("A", "approve")
	page, err = RecentApprovals(boundedTestContext(t), r, 40)
	if err != nil || page.Summaries[0].RequestRef != refsByName["A"].String() || page.Summaries[0].Status != "approve" {
		t.Fatalf("duplicate newest decision=%+v %v", page, err)
	}

}
func TestRecentApprovalsBounded(t *testing.T) {
	for _, duplicate := range []bool{false, true} {
		r := approvalViewFixture(200, duplicate)
		p, err := RecentApprovals(boundedTestContext(t), r, 40)
		if err != nil {
			t.Fatal(err)
		}
		if r.registryReads != 1 || r.headReads != 40 || len(p.Summaries) != 20 || !p.Truncated || r.bodyReads > 2*len(p.Summaries) || r.headReads+r.bodyReads > 80 {
			t.Fatalf("bounds heads=%d bodies=%d registry=%d page=%+v", r.headReads, r.bodyReads, r.registryReads, p)
		}
		if p.Summaries[0].DecidedBy != "operator-199" {
			t.Fatal("duplicate newest decision lost")
		}
	}
	r := &approvalViewReader{objects: map[hashref.HashRef]store.Object{}}
	p, err := RecentApprovals(boundedTestContext(t), r, 40)
	if err != nil || len(p.Summaries) != 0 || p.Truncated || r.registryReads != 1 || r.headReads+r.bodyReads != 0 {
		t.Fatal("no-head state incorrect")
	}
}
func TestRecentApprovalsMalformed(t *testing.T) {
	for _, kind := range []string{"semantic", "json", "missing-head", "missing-request", "missing-decision", "bad-ref", "cycle"} {
		t.Run(kind, func(t *testing.T) {
			r := approvalViewFixture(1, false)
			head := r.objects[r.head]
			var wire approvalHeadWire
			_ = decodeApprovalJSON(head.Payload, &wire)
			switch kind {
			case "semantic":
				head.SemanticID = "wrong"
				r.objects[r.head] = head
			case "json":
				head.Payload = []byte("{")
				r.objects[r.head] = head
			case "missing-head":
				delete(r.objects, r.head)
			case "missing-request":
				delete(r.objects, hashref.MustParse(wire.RequestRef))
			case "missing-decision":
				delete(r.objects, hashref.MustParse(wire.DecisionRef))
			case "bad-ref":
				wire.RequestRef = "bad"
				head.Payload = mustApprovalJSON(wire)
				r.objects[r.head] = head
			case "cycle":
				wire.PreviousHead = r.head.String()
				head.Payload = mustApprovalJSON(wire)
				r.objects[r.head] = head
			}
			p, err := RecentApprovals(boundedTestContext(t), r, 40)
			if !errors.Is(err, ErrApprovalChainMalformed) || p.MalformedRef == "" || r.headReads > 40 {
				t.Fatalf("malformed=%+v err=%v heads=%d", p, err, r.headReads)
			}
		})
	}
}
func TestRecentApprovalsStoreError(t *testing.T) {
	for _, registry := range []bool{true, false} {
		r := approvalViewFixture(1, false)
		r.failRegistry = registry
		r.failObject = !registry
		_, err := RecentApprovals(boundedTestContext(t), r, 40)
		if err == nil || errors.Is(err, ErrApprovalChainMalformed) {
			t.Fatalf("I/O classified as malformed: %v", err)
		}
	}
}
