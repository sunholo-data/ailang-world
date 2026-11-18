package daemon

// Row 243 V3 sprint 5 (the sprint-4 judge's note 1, mutant J-NOW-C14): C14 records the finish
// approval at the handler's own clock — the world decision carries the injected d.approveNow, as
// C15's consumed_at does: one instant for the whole confirmation.

import (
	"encoding/json"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

func TestConfirmPOSTRecordsAtItsClock(t *testing.T) {
	a, row := newConfirmRig(t)
	at := apT0 + 7
	a.d.approveNow = func() int64 { return at }
	if rec := a.confirm(row.RequestID, confirmSSH(t, apSigner, approval.VerbNamespace, row)); rec.Code != 200 {
		t.Fatalf("%d %s", rec.Code, rec.Body.String())
	}
	got := a.readVerb(row.RequestID)
	var decision struct {
		Now int64 `json:"now"`
	}
	obj, ok, err := a.d.store.GetObject(boundedTestContext(t), hashref.MustParse(got.DecisionRef))
	if err != nil || !ok || json.Unmarshal(obj.Payload, &decision) != nil || decision.Now != at || got.ConsumedAt != at {
		t.Fatalf("decision at %d (%v %v), consumed at %d; the clock said %d", decision.Now, ok, err, got.ConsumedAt, at)
	}
}
