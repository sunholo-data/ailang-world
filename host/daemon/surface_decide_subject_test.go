package daemon

// Row 201 A4 sprint 1 (design §2.8, §6 RD3): GET /surface/decisions/{ref} is the subject as bytes.

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// RD3: the bytes decode to exactly the stored payloads (the client re-hashes them); no message or
// hash field; the window; refusal/reserved/decision set by the law; D1-D4 answer the refusal.
func TestDecisionSubjectJSON(t *testing.T) {
	r := newDecideRig(t)
	q, p := r.packet("subject", broker.WorkspaceScope, "agent")
	rec := r.serve("GET", "/surface/decisions/"+q.String(), apLoopback, nil)
	if rec.Code != 200 {
		t.Fatalf("subject = %d %s", rec.Code, rec.Body.String())
	}
	if err := approvalHeaders(rec); err != nil {
		t.Fatal(err)
	}
	var keys map[string]json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &keys); err != nil {
		t.Fatal(err)
	}
	var names []string
	for k := range keys {
		names = append(names, k)
	}
	sort.Strings(names)
	if want := []string{"decidedBy", "decision", "entry", "now", "packet", "packetBytes", "passkeys", "proposal", "refusal", "request", "requestBytes", "reserved", "ssh", "v", "window"}; !reflect.DeepEqual(names, want) {
		t.Fatalf("subject keys %v, want %v (no message, no hash)", names, want)
	}
	var s surfacepage.DecisionSubject
	if err := json.Unmarshal(rec.Body.Bytes(), &s); err != nil {
		t.Fatal(err)
	}
	reqBytes, err1 := base64.RawURLEncoding.DecodeString(s.RequestBytes)
	pktBytes, err2 := base64.RawURLEncoding.DecodeString(s.PacketBytes)
	if err1 != nil || err2 != nil || sha256Of(reqBytes) != q.String() || sha256Of(pktBytes) != p.Hash.String() || s.Request != q.String() || s.Packet != p.Hash.String() {
		t.Fatalf("subject bytes do not hash to their refs: %+v", s)
	}
	stored, ok, err := r.d.store.GetObject(boundedTestContext(t), q)
	if err != nil || !ok || string(stored.Payload) != string(reqBytes) || string(p.Payload) != string(pktBytes) {
		t.Fatal("subject bytes are not the stored payloads")
	}
	packet, err := coordinator.DecodeDecisionPacket(pktBytes)
	if err != nil || packet.RequestRef != q || s.Proposal != packet.ProposalHash.String() || s.Entry < 1 {
		t.Fatalf("packet %+v, %v; subject %+v", packet, err, s)
	}
	if s.V != 1 || s.Now != decNow || s.Window != (surfacepage.DecisionWindow{Min: decNow + 1, Max: packetNow + 60}) ||
		s.Refusal != "" || s.Reserved != "" || s.Decision != "" || s.DecidedBy != "" || s.SSH.Namespace != "world-decision@sunholo" {
		t.Fatalf("subject %+v", s)
	}
	// A wide deadline: window.max is now + 600.
	r.d.approveNow = func() int64 { return packetNow - 600 }
	if err := json.Unmarshal(r.serve("GET", "/surface/decisions/"+q.String(), apLoopback, nil).Body.Bytes(), &s); err != nil || s.Window.Max != packetNow {
		t.Fatalf("window %+v, %v", s.Window, err)
	}
	r.d.approveNow = func() int64 { return decNow }
	// After a decision: the decision, and the law's D6 refusal as display.
	if _, err := broker.DecideApproval(boundedTestContext(t), r.d.store, q, "deny", "operator:terminal", decNow); err != nil {
		t.Fatal(err)
	}
	s = surfacepage.DecisionSubject{}
	if err := json.Unmarshal(r.serve("GET", "/surface/decisions/"+q.String(), apLoopback, nil).Body.Bytes(), &s); err != nil ||
		s.Decision != "deny" || s.DecidedBy != "operator:terminal" || !strings.Contains(s.Refusal, "was already decided deny by operator:terminal") {
		t.Fatalf("decided subject %+v, %v", s, err)
	}
	// Reserved: a finish-shaped request shows its owner and the D5 refusal.
	fq, _ := r.packet("subject-finish", "github-pr:o/r", "agent")
	s = surfacepage.DecisionSubject{}
	if err := json.Unmarshal(r.serve("GET", "/surface/decisions/"+fq.String(), apLoopback, nil).Body.Bytes(), &s); err != nil ||
		s.Reserved != "finish" || !strings.Contains(s.Refusal, "ailang-worldd episode finish") {
		t.Fatalf("reserved subject %+v, %v", s, err)
	}
	// D7 at the deadline is display too.
	r.d.approveNow = func() int64 { return packetNow + 60 }
	s = surfacepage.DecisionSubject{}
	if err := json.Unmarshal(r.serve("GET", "/surface/decisions/"+fq.String(), apLoopback, nil).Body.Bytes(), &s); err != nil || !strings.Contains(s.Refusal, "ailang-worldd episode finish") {
		t.Fatalf("reserved before the deadline: %+v", s)
	}
	r.d.approveNow = func() int64 { return decNow }
	// D1-D4: no subject, the law's refusal.
	pending := approvalRequestInChain(t, r.d, broker.WorkspaceScope, "no-packet", packetNow)
	for _, tc := range []struct {
		ref, class, reason string
		status             int
	}{
		{"sha256:XYZ", "BadRequest", "is not a request ref", 400},
		{"sha256:" + strings.Repeat("0", 64), "NotFound", "no such approval request", 404},
		{p.Hash.String(), "NotFound", "no such approval request", 404},
		{pending.String(), "Conflict", "no committed decision packet carries request", 409},
	} {
		got := requireRefusal(t, tc.ref, r.serve("GET", "/surface/decisions/"+tc.ref, apLoopback, nil), tc.status, tc.class)
		if !strings.Contains(got.Reason, tc.reason) {
			t.Errorf("%s: %q lacks %q", tc.ref, got.Reason, tc.reason)
		}
	}
}
