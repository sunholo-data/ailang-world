package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"strings"
	"testing"
)

func recordWorkbenchApproval(t *testing.T, d *Daemon, n int) hashref.HashRef {
	t.Helper()
	out, err := broker.NewHumanHandler(d.store).Execute(boundedTestContext(t), broker.EffectRequest{Effect: broker.EffectHumanApprove, Scope: fmt.Sprintf("scope-%d", n), Cost: 3, Now: int64(n)}, []byte(`{"requester":"workbench-requester"}`))
	if err != nil {
		t.Fatal(err)
	}
	var result struct {
		RequestRef string `json:"requestRef"`
	}
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatal(err)
	}
	return hashref.MustParse(result.RequestRef)
}
func TestWorkbenchDecisionsPane(t *testing.T) {
	d := newHandlerDaemon(t)
	s := workbenchNamedSection(t, workbenchLiveBody(t, d, "/workbench"), "decisions")
	if !strings.Contains(s, "No approval requests recorded") {
		t.Fatal("empty decisions missing")
	}
	workbenchLiveREST(t, d, []int64{0})
	ref := recordWorkbenchApproval(t, d, 0)
	if _, err := broker.DecideApproval(boundedTestContext(t), d.store, ref, "deny", "workbench-operator", 1); err != nil {
		t.Fatal(err)
	}
	s = workbenchNamedSection(t, workbenchLiveBody(t, d, "/workbench"), "decisions")
	for _, want := range []string{`href="/workbench?object=` + ref.String() + `"`, broker.EffectHumanApprove, "scope-0", "workbench-requester", "cost: 3", "deny by workbench-operator", "as of entry 0"} {
		if !strings.Contains(s, want) {
			t.Fatalf("decisions missing %q: %s", want, s)
		}
	}
	for i := 1; i <= 21; i++ {
		recordWorkbenchApproval(t, d, i)
	}
	s = workbenchNamedSection(t, workbenchLiveBody(t, d, "/workbench"), "decisions")
	if !strings.Contains(s, "Showing recent approvals; more recorded") || strings.Count(s, "cost: 3") != 20 {
		t.Fatal("truncation/count missing")
	}
}
func TestRecentApprovalsMalformed(t *testing.T) {
	for _, kind := range []string{"semantic", "json", "missing-request"} {
		t.Run(kind, func(t *testing.T) {
			d := newHandlerDaemon(t)
			payload := []byte(`{"previousHead":"","requestRef":"` + hashref.SumSHA256([]byte("absent request")).String() + `","decisionRef":""}`)
			id := broker.ApprovalsV1
			if kind == "semantic" {
				id = "wrong"
			}
			if kind == "json" {
				payload = []byte("{")
			}
			o := store.Object{Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte("approval-interface")), SemanticID: id, Payload: payload}
			if err := d.store.PutObject(boundedTestContext(t), o); err != nil {
				t.Fatal(err)
			}
			if err := d.store.SetRegistryHead(boundedTestContext(t), broker.ApprovalsV1, o.Hash); err != nil {
				t.Fatal(err)
			}
			r := requestRecorder(t, d, "GET", "/workbench", nil)
			if r.Code != 200 {
				t.Fatalf("malformed status=%d want 200", r.Code)
			}
			s := workbenchNamedSection(t, r.Body.String(), "decisions")
			offending := o.Hash.String()
			if kind == "missing-request" {
				offending = hashref.SumSHA256([]byte("absent request")).String()
			}
			if !strings.Contains(s, "UNAVAILABLE: approvals chain is malformed at "+offending) {
				t.Fatalf("malformed pane missing ref: %s", s)
			}
		})
	}
}

type workbenchApprovalFailReads struct{ readStore }

func (r workbenchApprovalFailReads) GetRegistryHead(ctx context.Context, name string) (hashref.HashRef, bool, error) {
	if name == broker.ApprovalsV1 {
		return hashref.HashRef{}, false, errors.New("SECRET approval I/O")
	}
	return r.readStore.GetRegistryHead(ctx, name)
}
func TestWorkbenchDecisionsStoreError(t *testing.T) {
	d := newHandlerDaemon(t)
	d.reads = workbenchApprovalFailReads{d.reads}
	r := requestRecorder(t, d, "GET", "/workbench", nil)
	if r.Code != 500 || !strings.Contains(r.Body.String(), workbenchInternalStoreFailureMessage) || strings.Contains(r.Body.String(), "SECRET") {
		t.Fatal("approvals I/O must be sanitized 500")
	}
}

// Companion keeps the plan's daemon mutation command non-vacuous.
func TestWorkbenchDecisionsReadOnly(t *testing.T) {
	d := newHandlerDaemon(t)
	recordWorkbenchApproval(t, d, 0)
	s := workbenchNamedSection(t, workbenchLiveBody(t, d, "/workbench"), "decisions")
	if !strings.Contains(s, "pending") {
		t.Fatal("decision rows missing")
	}
	for _, bad := range []string{"<form", "<button", "method="} {
		if strings.Contains(s, bad) {
			t.Fatalf("action in decisions: %s", bad)
		}
	}
}
