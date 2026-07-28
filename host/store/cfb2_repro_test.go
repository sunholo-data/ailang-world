// CF-B-2 is OPEN and UNFIXED; its fix is blocked on a human ratification
// packet. These characterization tests deliberately assert the current, wrong
// behavior and are expected to fail and be rewritten when the real fix lands.
package store

import (
	"strings"
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
			Ref:       hashref.SumSHA256([]byte("cf-b-2 world 1")),
			Revision:  genesis.Revision + 1,
			StateRoot: hashref.SumSHA256([]byte("cf-b-2 state 1")),
			LogHead:   entryHash,
		},
		Entry: LogEntry{
			Header: LogHeader{
				EntryIndex:     1,
				SemanticsEpoch: 1,
				TransitionFn:   hashref.SumSHA256([]byte("cf-b-2 transition fn")),
				Interpreter:    hashref.SumSHA256([]byte("cf-b-2 interpreter")),
				PrevEntryHash:  genesis.LogHead,
				WrittenBy:      "cf-b-2 characterization",
			},
			EntryHash:     entryHash,
			TransitionRef: body.Hash,
		},
	}
}

// TestCFB2CommitValidatesNoRefFieldOnWrite waits for Commit to reject every zero HashRef before writing.
func TestCFB2CommitValidatesNoRefFieldOnWrite(t *testing.T) {
	cases := []struct {
		name   string
		poison func(*Commit)
	}{
		{"transition fn", func(c *Commit) { c.Entry.Header.TransitionFn = hashref.HashRef{} }},
		{"interpreter", func(c *Commit) { c.Entry.Header.Interpreter = hashref.HashRef{} }},
		{"previous entry hash", func(c *Commit) { c.Entry.Header.PrevEntryHash = hashref.HashRef{} }},
		{"entry hash", func(c *Commit) { c.Entry.EntryHash = hashref.HashRef{} }},
		{"transition ref", func(c *Commit) { c.Entry.TransitionRef = hashref.HashRef{} }},
		{"world state root", func(c *Commit) { c.NextWorld.StateRoot = hashref.HashRef{} }},
		{"world log head", func(c *Commit) { c.NextWorld.LogHead = hashref.HashRef{} }},
		{"world ref", func(c *Commit) { c.NextWorld.Ref = hashref.HashRef{} }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := openMem(t)
			c := cfb2Commit(seedGenesis(t, s))
			tc.poison(&c)

			if err := s.Commit(c); err != nil {
				t.Fatalf("Commit rejected zero %s; current behavior is acceptance: %v", tc.name, err)
			}
		})
	}
}

// TestCFB2UnreadableLogEntryClass waits for Commit to reject references that make a stored log entry unreadable.
func TestCFB2UnreadableLogEntryClass(t *testing.T) {
	cases := []struct {
		name    string
		wantErr string
		poison  func(*Commit)
	}{
		{"transition fn", "transitionFn", func(c *Commit) { c.Entry.Header.TransitionFn = hashref.HashRef{} }},
		{"interpreter", "interpreter", func(c *Commit) { c.Entry.Header.Interpreter = hashref.HashRef{} }},
		{"previous entry hash", "prevEntryHash", func(c *Commit) { c.Entry.Header.PrevEntryHash = hashref.HashRef{} }},
		{"entry hash", "hash", func(c *Commit) { c.Entry.EntryHash = hashref.HashRef{} }},
		{"transition ref", "transitionRef", func(c *Commit) { c.Entry.TransitionRef = hashref.HashRef{} }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := openMem(t)
			c := cfb2Commit(seedGenesis(t, s))
			tc.poison(&c)
			if err := s.Commit(c); err != nil {
				t.Fatalf("Commit: %v", err)
			}

			_, ok, err := s.GetLogEntry(1)
			if err == nil || ok {
				t.Fatalf("GetLogEntry(1): ok=%v err=%v, want ok=false and non-nil error", ok, err)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("GetLogEntry(1) error %q does not mention %q", err, tc.wantErr)
			}
		})
	}
}

// TestCFB2UnloadableWorldClass waits for Commit to reject references that make the selected world unloadable.
func TestCFB2UnloadableWorldClass(t *testing.T) {
	cases := []struct {
		name    string
		wantErr string
		poison  func(*Commit)
	}{
		{"state root", "state root", func(c *Commit) { c.NextWorld.StateRoot = hashref.HashRef{} }},
		{"log head", "log head", func(c *Commit) { c.NextWorld.LogHead = hashref.HashRef{} }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := openMem(t)
			c := cfb2Commit(seedGenesis(t, s))
			tc.poison(&c)
			if err := s.Commit(c); err != nil {
				t.Fatalf("Commit: %v", err)
			}

			if _, ok, err := s.GetLogEntry(1); err != nil || !ok {
				t.Fatalf("GetLogEntry(1): ok=%v err=%v, want readable entry", ok, err)
			}
			_, ok, err := s.GetWorld(c.NextWorld.Ref)
			if err == nil || ok {
				t.Fatalf("GetWorld(next): ok=%v err=%v, want ok=false and non-nil error", ok, err)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("GetWorld(next) error %q does not mention %q", err, tc.wantErr)
			}
			head, ok, err := s.SelectedHead()
			if err != nil || !ok {
				t.Fatalf("SelectedHead: ok=%v err=%v", ok, err)
			}
			if head != c.NextWorld.Ref {
				t.Fatalf("SelectedHead = %q, want %q", head, c.NextWorld.Ref)
			}
		})
	}
}

// TestCFB2ZeroWorldRefWedgesTheStore waits for Commit to reject a zero world ref instead of permanently wedging writes.
func TestCFB2ZeroWorldRefWedgesTheStore(t *testing.T) {
	s := openMem(t)
	genesis := seedGenesis(t, s)
	c := cfb2Commit(genesis)
	c.NextWorld.Ref = hashref.HashRef{}
	if err := s.Commit(c); err != nil {
		t.Fatalf("poisoned Commit: %v", err)
	}

	if _, ok, err := s.GetLogEntry(1); err != nil || !ok {
		t.Fatalf("GetLogEntry(1): ok=%v err=%v, want readable entry", ok, err)
	}
	if _, ok, err := s.GetWorld(c.NextWorld.Ref); err != nil || !ok {
		t.Fatalf("GetWorld(zero ref): ok=%v err=%v, want readable world", ok, err)
	}
	_, ok, err := s.SelectedHead()
	if err == nil || ok {
		t.Fatalf("SelectedHead: ok=%v err=%v, want ok=false and non-nil error", ok, err)
	}
	if !strings.Contains(err.Error(), "selected head: hashref: empty hashref text") {
		t.Fatalf("SelectedHead error = %q, want empty-hashref error", err)
	}

	next := cfb2Commit(genesis)
	next.Entry.Header.EntryIndex = 2
	next.Entry.Header.PrevEntryHash = c.NextWorld.LogHead
	next.Entry.EntryHash = hashref.SumSHA256([]byte("cf-b-2 entry 2"))
	next.NextWorld.Ref = hashref.SumSHA256([]byte("cf-b-2 world 2"))
	next.NextWorld.Revision = 2
	next.NextWorld.LogHead = next.Entry.EntryHash
	err = s.Commit(next)
	if err == nil {
		t.Fatal("subsequent Commit succeeded; want wedged-store error")
	}
	if IsConflict(err) {
		t.Fatalf("subsequent Commit error is a ConflictError; want non-conflict wedge: %v", err)
	}
	if !strings.Contains(err.Error(), "selected head: hashref: empty hashref text") {
		t.Fatalf("subsequent Commit error = %q, want selected-head empty-hashref error", err)
	}
}

// TestCFB2PoisonedEntryLeavesPermanentHoleMidChain waits for Commit to reject a log chain containing an unreadable entry.
func TestCFB2PoisonedEntryLeavesPermanentHoleMidChain(t *testing.T) {
	s := openMem(t)
	genesis := seedGenesis(t, s)
	c1 := cfb2Commit(genesis)
	c1.Entry.Header.PrevEntryHash = hashref.HashRef{}
	if err := s.Commit(c1); err != nil {
		t.Fatalf("Commit(entry 1): %v", err)
	}
	if _, ok, err := s.GetLogEntry(1); err == nil || ok {
		t.Fatalf("GetLogEntry(1): ok=%v err=%v, want permanent unreadable hole", ok, err)
	}
	head, ok, err := s.SelectedHead()
	if err != nil || !ok {
		t.Fatalf("SelectedHead after entry 1: ok=%v err=%v", ok, err)
	}
	if head != c1.NextWorld.Ref {
		t.Fatalf("SelectedHead after entry 1 = %q, want %q", head, c1.NextWorld.Ref)
	}

	c2 := cfb2Commit(genesis)
	c2.ObservedHead = c1.NextWorld.Ref
	c2.Entry.Header.EntryIndex = 2
	c2.Entry.Header.PrevEntryHash = c1.NextWorld.LogHead
	c2.Entry.EntryHash = hashref.SumSHA256([]byte("cf-b-2 entry 2"))
	c2.NextWorld.Ref = hashref.SumSHA256([]byte("cf-b-2 world 2"))
	c2.NextWorld.Revision = 2
	c2.NextWorld.StateRoot = hashref.SumSHA256([]byte("cf-b-2 state 2"))
	c2.NextWorld.LogHead = c2.Entry.EntryHash
	if err := s.Commit(c2); err != nil {
		t.Fatalf("Commit(entry 2): %v", err)
	}
	if _, ok, err := s.GetLogEntry(2); err != nil || !ok {
		t.Fatalf("GetLogEntry(2): ok=%v err=%v, want readable entry after hole", ok, err)
	}
	if _, ok, err := s.GetLogEntry(1); err == nil || ok {
		t.Fatalf("GetLogEntry(1) after entry 2: ok=%v err=%v, want hole to remain", ok, err)
	}
}
