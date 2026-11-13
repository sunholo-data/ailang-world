package surfacepage

// Row 201 A4 sprint 4 (design §6 DMD1–DMD3): the decisions fold, without a daemon or a store.

import (
	"reflect"
	"testing"
)

// DMD1: pending Human.Approve requests only; decided, other effects and unpacketed requests are
// omitted; one packet gives the packet and its deadline, the reserved owner and refusal pass through.
func TestDecideModelListsPendingDecisionsOnly(t *testing.T) {
	cands := []DecisionCandidate{
		{Request: "sha256:01", Effect: "Human.Approve", Scope: "workspace:w", Requester: "agent-1", Cost: 3, Packets: []string{"sha256:p1"}, Deadline: 1760},
		{Request: "sha256:02", Effect: "Human.Approve", Scope: "workspace:w", Requester: "agent-2", Decided: true, Packets: []string{"sha256:p2"}, Deadline: 1760},
		{Request: "sha256:03", Effect: "Registry.Publish", Scope: "registry:x", Requester: "agent-3", Packets: []string{"sha256:p3"}, Deadline: 1760},
		{Request: "sha256:04", Effect: "Human.Approve", Scope: "workspace:w", Requester: "agent-4"},
		{Request: "sha256:05", Effect: "Human.Approve", Scope: "registry:pkg#effect=Registry.Publish", Requester: "agent-5", Packets: []string{"sha256:p5"},
			Deadline: 1770, Reserved: "publish", Refusal: "publish approvals are minted by `world-publish approve`"},
	}
	rows, truncated := FoldDecisions(cands, false, false)
	want := []PendingDecision{
		{Request: "sha256:01", Packet: "sha256:p1", Effect: "Human.Approve", Scope: "workspace:w", Cost: 3, Requester: "agent-1", Deadline: 1760},
		{Request: "sha256:05", Packet: "sha256:p5", Effect: "Human.Approve", Scope: "registry:pkg#effect=Registry.Publish", Requester: "agent-5", Deadline: 1770,
			Reserved: "publish", Refusal: "publish approvals are minted by `world-publish approve`"},
	}
	if !reflect.DeepEqual(rows, want) || truncated {
		t.Fatalf("FoldDecisions = %+v, %v; want %+v, false", rows, truncated, want)
	}
}

// DMD2: a request two packets carry is listed with its refusal and no packet or deadline.
func TestDecideModelListsATwoPacketRequestRefused(t *testing.T) {
	rows, _ := FoldDecisions([]DecisionCandidate{{Request: "sha256:06", Effect: "Human.Approve", Scope: "s", Requester: "a",
		Packets: []string{"sha256:pa", "sha256:pb"}, Deadline: 1760, Refusal: "request sha256:06 is carried by 2 decision packets"}}, false, false)
	want := []PendingDecision{{Request: "sha256:06", Effect: "Human.Approve", Scope: "s", Requester: "a", Refusal: "request sha256:06 is carried by 2 decision packets"}}
	if !reflect.DeepEqual(rows, want) {
		t.Fatalf("two packets: %+v, want %+v", rows, want)
	}
}

// DMD3: truncated is either bound: the approvals read's or the packet scan's.
func TestDecideModelTruncation(t *testing.T) {
	for _, c := range []struct{ chain, scan, want bool }{{false, false, false}, {true, false, true}, {false, true, true}, {true, true, true}} {
		if _, got := FoldDecisions(nil, c.chain, c.scan); got != c.want {
			t.Errorf("chain %v scan %v: truncated %v, want %v", c.chain, c.scan, got, c.want)
		}
	}
}
