package store

import (
	"context"
	"testing"
	"time"
)

// TestCommitLandedReadsTheLogRowAtItsIndex pins INV-LOG's four observations
// (row 23 design §2.5) on a real file store. Each arm is a distinct refusal or
// acceptance branch of CommitLanded, so each has its own mutation.
func TestCommitLandedReadsTheLogRowAtItsIndex(t *testing.T) {
	ctx := func(t *testing.T) context.Context {
		c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		t.Cleanup(cancel)
		return c
	}
	t.Run("absent_row_is_not_committed", func(t *testing.T) {
		s := openFileStore(t)
		a := bareCommit(seedGenesis(t, s), 1, "A")
		if landed, err := s.CommitLanded(ctx(t), a); err != nil || landed {
			t.Fatalf("CommitLanded(A) with no row = %v, %v; want false, nil (not committed, resubmit allowed)", landed, err)
		}
	})
	t.Run("equal_row_is_entry_landed", func(t *testing.T) {
		s := openFileStore(t)
		a := bareCommit(seedGenesis(t, s), 1, "A")
		if err := s.Commit(ctx(t), a); err != nil {
			t.Fatalf("commit A: %v", err)
		}
		if landed, err := s.CommitLanded(ctx(t), a); err != nil || !landed {
			t.Fatalf("CommitLanded(A) after A landed = %v, %v; want true, nil", landed, err)
		}
	})
	t.Run("same_hash_different_header_is_not_committed", func(t *testing.T) {
		// A' carries A's EntryHash (the store never recomputes it) but a
		// different header field, so a hash-only comparison would misread it.
		s := openFileStore(t)
		a := bareCommit(seedGenesis(t, s), 1, "A")
		other := a
		other.Entry.Header.WrittenBy = "someone-else"
		if err := s.Commit(ctx(t), other); err != nil {
			t.Fatalf("commit A': %v", err)
		}
		if landed, err := s.CommitLanded(ctx(t), a); err != nil || landed {
			t.Fatalf("CommitLanded(A) over A' (same EntryHash, different WrittenBy) = %v, %v; want false, nil", landed, err)
		}
	})
	t.Run("read_error_is_unknown", func(t *testing.T) {
		s := openFileStore(t)
		a := bareCommit(seedGenesis(t, s), 1, "A")
		expired, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
		defer cancel()
		if landed, err := s.CommitLanded(expired, a); err == nil || landed {
			t.Fatalf("CommitLanded under an expired ctx = %v, %v; want false with an error (unknown)", landed, err)
		}
	})
}
