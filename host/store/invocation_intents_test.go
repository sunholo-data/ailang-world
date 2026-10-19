package store

import (
	"errors"
	"testing"
)

// Row 198: InvocationIntents reads exactly one invocation-id key range, in
// seq order, with the outcome the commit appended.
func TestInvocationIntentsReadsOneKeyRangeWithOutcomes(t *testing.T) {
	s := openMem(t)
	ctx := boundedTestContext(t)
	c := journalCommitFixture(t, s, "mcp:ep:1")
	if _, _, err := s.AppendIntent(ctx, c.InvocationID, testCommitIntent(c.InvocationID, c)); err != nil {
		t.Fatal(err)
	}
	if err := s.Commit(ctx, c); err != nil {
		t.Fatal(err)
	}
	// Neighbours of the range: a longer episode, a sibling, the other
	// surface, and an uncommitted intent inside the range.
	for _, id := range []string{"mcp:ep2:1", "mcp:eq:1", "a2a:ep:1", "mcp:ep:2"} {
		if _, _, err := s.AppendIntent(ctx, id, testCommitIntent(id, c)); err != nil {
			t.Fatal(err)
		}
	}
	got, err := s.InvocationIntents(ctx, "mcp:ep:", 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Intent.InvocationID != "mcp:ep:1" || got[1].Intent.InvocationID != "mcp:ep:2" {
		t.Fatalf("intents = %+v, want mcp:ep:1 then mcp:ep:2", got)
	}
	if got[0].Outcome == nil || got[0].Outcome.ResultRef != c.NextWorld.Ref || got[1].Outcome != nil {
		t.Errorf("outcomes = %+v / %+v, want the committed world, then none", got[0].Outcome, got[1].Outcome)
	}
	if got[0].Seq >= got[1].Seq || got[0].Intent.LogicalTime != 42 {
		t.Errorf("seq/time = %d,%d/%d", got[0].Seq, got[1].Seq, got[0].Intent.LogicalTime)
	}
	page, err := s.InvocationIntents(ctx, "mcp:ep:", got[0].Seq, 1)
	if err != nil || len(page) != 1 || page[0].Intent.InvocationID != "mcp:ep:2" {
		t.Errorf("cursor page = %+v, %v", page, err)
	}
	for _, bad := range []string{"mcp:ep", "effect:ep:", ""} {
		var prefixErr *InvalidIntentPrefixError
		if _, err := s.InvocationIntents(ctx, bad, 0, 10); !errors.As(err, &prefixErr) {
			t.Errorf("prefix %q: err = %v, want InvalidIntentPrefixError", bad, err)
		}
	}
	for _, limit := range []int{0, MaxInvocationIntentsPage + 1} {
		var limitErr *InvalidLimitError
		if _, err := s.InvocationIntents(ctx, "mcp:ep:", 0, limit); !errors.As(err, &limitErr) {
			t.Errorf("limit %d: err = %v, want InvalidLimitError", limit, err)
		}
	}
}
