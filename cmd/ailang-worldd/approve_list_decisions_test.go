package main

// Row 201 A4 sprint 4 (design §2.10, Q2): approve --list prints the pending decisions of the world
// store --db, read-only and offline (it never dials a daemon). Temp stores only (D-WORLD-109).

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// TestApproveListShowsDecisions: a Human.Approve request on the --db store's approvals chain whose
// one stored packet is not carried by a commit is listed with the law's refusal; the mint list is
// printed as before; a --db that cannot be read is one stderr line and exit 0.
func TestApproveListShowsDecisions(t *testing.T) {
	db, creds := elsewhereStores(t)
	putPendingRequest(t, creds)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	st, err := store.Open(db)
	if err != nil {
		t.Fatal(err)
	}
	out, err := broker.NewHumanHandler(st).Execute(ctx, broker.EffectRequest{Effect: broker.EffectHumanApprove, Scope: broker.WorkspaceScope, Cost: 1, Now: 1700},
		[]byte(`{"requester":"agent-list"}`))
	if err != nil {
		t.Fatal(err)
	}
	q, err := broker.ParsePendingApproval(out)
	if err != nil {
		t.Fatal(err)
	}
	payload := coordinator.EncodeDecisionPacket(coordinator.DecisionPacket{ProposalHash: hashref.SumSHA256([]byte("plan")), RequestRef: q,
		CreatedAt: 1700, DeadlineAt: 1760, Policy: "Cancel"})
	if err := st.PutObject(ctx, store.Object{Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte(coordinator.DecisionPacketV1)),
		SemanticID: coordinator.DecisionPacketV1, Provenance: "list-test", Payload: payload}); err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if code := runApprove([]string{"--list", "--db", db, "--credentials-db", creds}, &stdout, &stderr, approveEnv(nil, nil, atTime(approveT0+1))); code != exitOK {
		t.Fatalf("rc %d, stderr %q", code, stderr.String())
	}
	want := firstRequestID + "  episode ep1  digest abababababababab  expires 2027-01-15T08:10:00Z\n" +
		"pending decisions (Human.Approve requests):\n" +
		q.String() + "  " + broker.WorkspaceScope + "  requested by agent-list  refused: the decision packet carrying request " + q.String() + " is refused: "
	if got := stdout.String(); !strings.HasPrefix(got, want) || strings.Count(got, "\n") != 3 || stderr.Len() != 0 {
		t.Fatalf("stdout %q\nwant the prefix %q; stderr %q", got, want, stderr.String())
	}
	stdout.Reset()
	missing := filepath.Join(t.TempDir(), "missing-world.db")
	if code := runApprove([]string{"--list", "--db", missing, "--credentials-db", creds}, &stdout, &stderr, approveEnv(nil, nil, atTime(approveT0+1))); code != exitOK ||
		!strings.HasPrefix(stdout.String(), firstRequestID) || strings.Contains(stdout.String(), "decisions") ||
		!strings.Contains(stderr.String(), "ailang-worldd approve: decisions: cannot") {
		t.Fatalf("missing world store: rc %d, stdout %q, stderr %q", code, stdout.String(), stderr.String())
	}
}

// TestApproveListDecisionLines: the section's text for a decidable, a refused and a cut list.
func TestApproveListDecisionLines(t *testing.T) {
	var b bytes.Buffer
	printDecisions(&b, nil, false)
	if b.Len() != 0 {
		t.Fatalf("no decision printed %q", b.String())
	}
	printDecisions(&b, []surfacepage.PendingDecision{
		{Request: "sha256:aa", Packet: "sha256:pp", Scope: "workspace:w", Requester: "agent-1", Deadline: 1760},
		{Request: "sha256:bb", Scope: "registry:x", Requester: "agent-2", Reserved: "publish", Refusal: "publish approvals are minted by `world-publish approve` at a terminal (R4)"},
	}, true)
	want := "pending decisions (Human.Approve requests):\n" +
		"sha256:aa  workspace:w  requested by agent-1  deadline 1970-01-01T00:29:20Z  decide: ailang-worldd approve sha256:aa [--deny]\n" +
		"sha256:bb  registry:x  requested by agent-2  refused: publish approvals are minted by `world-publish approve` at a terminal (R4)\n" +
		"(more requests than listed: the approvals chain or the packet scan reached its bound)\n"
	if b.String() != want {
		t.Fatalf("got\n%s\nwant\n%s", b.String(), want)
	}
	b.Reset()
	printDecisions(&b, nil, true)
	if b.String() != "pending decisions (Human.Approve requests):\n(more requests than listed: the approvals chain or the packet scan reached its bound)\n" {
		t.Fatalf("cut and empty: %q", b.String())
	}
}
