package broker

import (
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// TestPacketBrokerExports pins the row-172 broker exports the decision-packet
// producer relies on: AuthorizingGrant (D3), ParsePendingApproval and
// DecodeApprovalRequest.
func TestPacketBrokerExports(t *testing.T) {
	req := EffectRequest{Effect: EffectHumanApprove, Scope: "demo", Cost: 1, Now: 100}
	expired := Capability{Effect: EffectHumanApprove, Scope: "demo", ExpiresAt: 100, Budget: 5}
	live1760 := Capability{Effect: EffectHumanApprove, Scope: "demo", ExpiresAt: 1760, Budget: 5}
	live1900 := Capability{Effect: EffectHumanApprove, Scope: "demo", ExpiresAt: 1900, Budget: 5}
	ledger := []Capability{expired, live1760, live1900}

	got, ok := AuthorizingGrant(ledger, req)
	if !ok || got != live1760 {
		t.Fatalf("AuthorizingGrant = %+v, %v; want the 1760 grant (first allowed in order is index 1)", got, ok)
	}
	for _, c := range ledger {
		if Decide(c, req).Allowed {
			if got != c {
				t.Fatalf("first Decide-allowed grant %+v differs from AuthorizingGrant %+v", c, got)
			}
			break
		}
	}
	if _, ok := AuthorizingGrant([]Capability{expired}, req); ok {
		t.Fatal("all-expired ledger: AuthorizingGrant reported an authorizing grant")
	}
	if _, ok := AuthorizingGrant(nil, req); ok {
		t.Fatal("empty ledger: AuthorizingGrant reported an authorizing grant")
	}

	ref := hashref.SumSHA256([]byte("request"))
	if back, err := ParsePendingApproval(pendingBytes(ref)); err != nil || back != ref {
		t.Fatalf("ParsePendingApproval(pending) = %v, %v; want %v", back, err, ref)
	}
	refusals := map[string]string{
		"status approve": `{"status":"approve","requestRef":"` + ref.String() + `"}`,
		"extra field":    `{"status":"pending","requestRef":"` + ref.String() + `","x":1}`,
		"malformed ref":  `{"status":"pending","requestRef":"sha256:zz"}`,
	}
	for name, bad := range refusals {
		if _, err := ParsePendingApproval([]byte(bad)); err == nil {
			t.Fatalf("ParsePendingApproval %s: no error for %s", name, bad)
		}
	}

	s := openTestStore(t)
	ctx := boundedTestContext(t)
	out, err := NewHumanHandler(s).Execute(ctx, req, []byte(`{"requester":"coord-test"}`))
	if err != nil {
		t.Fatalf("Human.Approve: %v", err)
	}
	requestRef, err := ParsePendingApproval(out)
	if err != nil {
		t.Fatalf("parse pending output %s: %v", out, err)
	}
	obj, found, err := s.GetObject(ctx, requestRef)
	if err != nil || !found {
		t.Fatalf("approval request object: found=%v err=%v", found, err)
	}
	decoded, err := DecodeApprovalRequest(obj.Payload)
	if err != nil {
		t.Fatalf("DecodeApprovalRequest: %v", err)
	}
	want := ApprovalRequest{Effect: EffectHumanApprove, Scope: "demo", Cost: 1, Requester: "coord-test", Now: 100}
	if decoded != want {
		t.Fatalf("decoded request = %+v, want %+v", decoded, want)
	}
	if _, err := DecodeApprovalRequest([]byte(`{"effect":"Human.Approve","scope":"demo","cost":1,"requester":"x","now":100,"extra":1}`)); err == nil {
		t.Fatal("DecodeApprovalRequest accepted an extra field")
	}
}
