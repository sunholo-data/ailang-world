package store

import (
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

func TestValidateRefTextRejectsNonCanonicalText(t *testing.T) {
	// HashRef's representation is private, so zero is the only non-canonical
	// HashRef constructible here. Testing the string primitive keeps the full
	// hashref.Parse canonicality claim non-vacuous.
	for _, text := range []string{"", "sha256:", "abc", "SHA256:ab", "sha256:ZZ", "a:b:c", "md5:ab"} {
		t.Run(text, func(t *testing.T) {
			err := validateRefText("test", "ref", text)
			var invalid *InvalidRefError
			if !errors.As(err, &invalid) || invalid.Op != "test" || invalid.Field != "ref" || invalid.Text != text {
				t.Fatalf("got %#v, want structured invalid-ref error", err)
			}
			if invalid.Err == nil || !IsInvalidRef(err) {
				t.Fatalf("error does not preserve cause/type: %#v", invalid)
			}
		})
	}
}

func TestEveryPersistedRefFieldIsValidated(t *testing.T) {
	type testCase struct {
		name, op, field string
		run             func(*Store, World) error
	}
	zero := hashref.HashRef{}
	commitCase := func(name, field string, poison func(*Commit)) testCase {
		return testCase{name, "Commit", field, func(s *Store, genesis World) error {
			c := cfb2Commit(genesis)
			poison(&c)
			return s.Commit(c)
		}}
	}
	cases := []testCase{
		commitCase("Commit TransitionFn", "Entry.Header.TransitionFn", func(c *Commit) { c.Entry.Header.TransitionFn = zero }),
		commitCase("Commit Interpreter", "Entry.Header.Interpreter", func(c *Commit) { c.Entry.Header.Interpreter = zero }),
		commitCase("Commit PrevEntryHash", "Entry.Header.PrevEntryHash", func(c *Commit) { c.Entry.Header.PrevEntryHash = zero }),
		commitCase("Commit EntryHash", "Entry.EntryHash", func(c *Commit) { c.Entry.EntryHash = zero }),
		commitCase("Commit TransitionRef", "Entry.TransitionRef", func(c *Commit) { c.Entry.TransitionRef = zero }),
		commitCase("Commit StateRoot", "NextWorld.StateRoot", func(c *Commit) { c.NextWorld.StateRoot = zero }),
		commitCase("Commit LogHead", "NextWorld.LogHead", func(c *Commit) { c.NextWorld.LogHead = zero }),
		commitCase("Commit World Ref", "NextWorld.Ref", func(c *Commit) { c.NextWorld.Ref = zero }),
		commitCase("Commit Object Hash", "Objects.Hash", func(c *Commit) { c.Objects[0].Hash = zero }),
		commitCase("Commit Object InterfaceHash", "Objects.InterfaceHash", func(c *Commit) { c.Objects[0].InterfaceHash = zero }),
		{"PutObject Hash", "PutObject", "Hash", func(s *Store, _ World) error {
			o := obj("validate", "validate")
			o.Hash = zero
			return s.PutObject(o)
		}},
		{"PutObject InterfaceHash", "PutObject", "InterfaceHash", func(s *Store, _ World) error {
			o := obj("validate", "validate")
			o.InterfaceHash = zero
			return s.PutObject(o)
		}},
		{"PutWorld Ref", "PutWorld", "Ref", func(s *Store, g World) error { g.Ref = zero; return s.PutWorld(g) }},
		{"PutWorld StateRoot", "PutWorld", "StateRoot", func(s *Store, g World) error { g.StateRoot = zero; return s.PutWorld(g) }},
		{"PutWorld LogHead", "PutWorld", "LogHead", func(s *Store, g World) error { g.LogHead = zero; return s.PutWorld(g) }},
		{"SetRegistryHead", "SetRegistryHead", "objectRef", func(s *Store, _ World) error { return s.SetRegistryHead("bad", zero) }},
		{"SelectHead", "SelectHead", "ref", func(s *Store, _ World) error { return s.SelectHead(zero) }},
		{"PutVerifyResult TransitionFn", "PutVerifyResult", "TransitionFn", func(s *Store, _ World) error {
			return s.PutVerifyResult(VerifyResult{TransitionFn: zero, Interpreter: hashref.SumSHA256([]byte("i"))})
		}},
		{"PutVerifyResult Interpreter", "PutVerifyResult", "Interpreter", func(s *Store, _ World) error {
			return s.PutVerifyResult(VerifyResult{TransitionFn: hashref.SumSHA256([]byte("f")), Interpreter: zero})
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := openMem(t)
			genesis := seedGenesis(t, s)
			before := snapshotStore(t, s)
			err := tc.run(s, genesis)
			var invalid *InvalidRefError
			if !errors.As(err, &invalid) || invalid.Op != tc.op || invalid.Field != tc.field {
				t.Fatalf("got %#v, want %s/%s InvalidRefError", err, tc.op, tc.field)
			}
			if after := snapshotStore(t, s); after != before {
				t.Fatalf("store changed: before=%+v after=%+v", before, after)
			}
		})
	}
}

func TestGenesisObservedHeadRemainsZeroLegal(t *testing.T) {
	s := openMem(t)
	c := cfb2Commit(World{})
	c.ObservedHead = hashref.HashRef{}
	c.Entry.Header.EntryIndex = 0
	c.Entry.Header.PrevEntryHash = hashref.SumSHA256([]byte("genesis previous"))
	if err := s.Commit(c); err != nil {
		t.Fatalf("genesis Commit rejected zero ObservedHead: %v", err)
	}
}
