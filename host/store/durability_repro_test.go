// CF-B-2 is closed by ARM V1: every reference persisted by Commit is validated
// before a transaction begins. These tests retain the three original damage
// classes and prove rejection leaves both rows and the selected head untouched.
package store

import (
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

func cfb2Commit(genesis World) Commit {
	body := obj("cf-b-2 transition body", "transition/cf-b-2")
	entryHash := hashref.SumSHA256([]byte("cf-b-2 entry 1"))
	return Commit{
		ObservedHead: genesis.Ref,
		Objects:      []Object{body},
		NextWorld: World{
			Ref: hashref.SumSHA256([]byte("cf-b-2 world 1")), Revision: genesis.Revision + 1,
			StateRoot: hashref.SumSHA256([]byte("cf-b-2 state 1")), LogHead: entryHash,
		},
		Entry: LogEntry{
			Header: LogHeader{
				EntryIndex: 1, SemanticsEpoch: 1,
				TransitionFn:  hashref.SumSHA256([]byte("cf-b-2 transition fn")),
				Interpreter:   hashref.SumSHA256([]byte("cf-b-2 interpreter")),
				PrevEntryHash: genesis.LogHead, WrittenBy: "cf-b-2 characterization",
			},
			EntryHash: entryHash, TransitionRef: body.Hash,
		},
	}
}

type storeState struct {
	counts [4]int
	head   string
	ok     bool
}

func snapshotStore(t *testing.T, s *Store) storeState {
	t.Helper()
	var got storeState
	for i, table := range []string{"objects", "worlds", "log_entries", "store_heads"} {
		if err := s.db.QueryRow("SELECT count(*) FROM " + table).Scan(&got.counts[i]); err != nil {
			t.Fatal(err)
		}
	}
	head, ok, err := s.SelectedHead()
	if err != nil {
		t.Fatal(err)
	}
	got.head, got.ok = head.String(), ok
	return got
}

func assertRejectedUntouched(t *testing.T, s *Store, before storeState, c Commit, field string) {
	t.Helper()
	err := s.Commit(c)
	var invalid *InvalidRefError
	if !errors.As(err, &invalid) {
		t.Fatalf("Commit error = %T %v, want *InvalidRefError", err, err)
	}
	if invalid.Op != "Commit" || invalid.Field != field {
		t.Fatalf("InvalidRefError = %#v, want Commit/%s", invalid, field)
	}
	if after := snapshotStore(t, s); after != before {
		t.Fatalf("store changed after rejection: before=%+v after=%+v", before, after)
	}
}

// CLASS 1: these fields formerly made the appended log entry unreadable.
func TestCFB2UnreadableLogEntryClass(t *testing.T) {
	cases := []struct {
		name, field string
		poison      func(*Commit)
	}{
		{"transition fn", "Entry.Header.TransitionFn", func(c *Commit) { c.Entry.Header.TransitionFn = hashref.HashRef{} }},
		{"interpreter", "Entry.Header.Interpreter", func(c *Commit) { c.Entry.Header.Interpreter = hashref.HashRef{} }},
		{"previous entry hash", "Entry.Header.PrevEntryHash", func(c *Commit) { c.Entry.Header.PrevEntryHash = hashref.HashRef{} }},
		{"entry hash", "Entry.EntryHash", func(c *Commit) { c.Entry.EntryHash = hashref.HashRef{} }},
		{"transition ref", "Entry.TransitionRef", func(c *Commit) { c.Entry.TransitionRef = hashref.HashRef{} }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := openMem(t)
			c := cfb2Commit(seedGenesis(t, s))
			before := snapshotStore(t, s)
			tc.poison(&c)
			assertRejectedUntouched(t, s, before, c, tc.field)
		})
	}
}

// CLASS 2: these fields formerly made the selected world unloadable.
func TestCFB2UnloadableWorldClass(t *testing.T) {
	cases := []struct {
		name, field string
		poison      func(*Commit)
	}{
		{"state root", "NextWorld.StateRoot", func(c *Commit) { c.NextWorld.StateRoot = hashref.HashRef{} }},
		{"log head", "NextWorld.LogHead", func(c *Commit) { c.NextWorld.LogHead = hashref.HashRef{} }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := openMem(t)
			c := cfb2Commit(seedGenesis(t, s))
			before := snapshotStore(t, s)
			tc.poison(&c)
			assertRejectedUntouched(t, s, before, c, tc.field)
		})
	}
}

// CLASS 3: the old assertion that a zero-ref world remained readable was the
// wedge's premise, not lenience: the poison was invisible to GetWorld, so no
// read-side repair could reach the selected-head failure or unblock later writes.
func TestCFB2ZeroWorldRefWedgeRejected(t *testing.T) {
	s := openMem(t)
	genesis := seedGenesis(t, s)
	c := cfb2Commit(genesis)
	before := snapshotStore(t, s)
	c.NextWorld.Ref = hashref.HashRef{}
	assertRejectedUntouched(t, s, before, c, "NextWorld.Ref")

	// AC1's CLASS 3 half: rejection alone is NOT the property under test. What
	// made CLASS 3 the worst of the eight is that the poison WEDGED the store —
	// SelectedHead() errored with a non-ConflictError and every later Commit
	// inherited that error, so a caller's standard re-plan-on-conflict path never
	// fired and the store could never accept another write. Asserting the bad
	// commit was refused says nothing about that. Only a subsequent VALID commit
	// succeeding proves the store is still live after a refusal.
	if _, _, err := s.SelectedHead(); err != nil {
		t.Fatalf("SelectedHead errored after a REFUSED commit: %v — that is the CLASS 3 wedge this fix exists to prevent", err)
	}
	good := cfb2Commit(genesis)
	if err := s.Commit(good); err != nil {
		t.Fatalf("a valid Commit after a refused one failed: %v — refusal must leave the store able to accept writes", err)
	}
	ref, ok, err := s.SelectedHead()
	if err != nil || !ok {
		t.Fatalf("SelectedHead after the valid commit: err=%v ok=%v", err, ok)
	}
	if ref.String() != good.NextWorld.Ref.String() {
		t.Fatalf("selected head = %q, want the valid commit's world %q", ref.String(), good.NextWorld.Ref.String())
	}
}
