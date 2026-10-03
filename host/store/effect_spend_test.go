package store

import (
	"path/filepath"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

func spendIntent(ep, effect, scope string, cost int64) EffectIntent {
	in := effectIntentFixture(ep, 1)
	in.Effect, in.Scope, in.Cost = effect, scope, cost
	return in
}

func mustAppendEffect(t *testing.T, s *Store, in EffectIntent) {
	t.Helper()
	if _, _, err := s.AppendNextEffectIntent(boundedTestContext(t), in.EpisodeID, in); err != nil {
		t.Fatal(err)
	}
}

func assertSpend(t *testing.T, s *Store, ep string, want map[EffectKey]int64) {
	t.Helper()
	got, err := s.EffectSpend(boundedTestContext(t), ep)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("EffectSpend(%q) = %v, want %v", ep, got, want)
	}
	for k, v := range want {
		if got[k] != v {
			t.Fatalf("EffectSpend(%q) = %v, want %v", ep, got, want)
		}
	}
}

func TestEffectSpendTable(t *testing.T) {
	ref := func(n string) hashref.HashRef { return hashref.SumSHA256([]byte(n)) }
	cases := []struct {
		name  string
		setup func(t *testing.T, s *Store)
		ep    string
		want  map[EffectKey]int64
	}{
		{"empty episode", func(t *testing.T, s *Store) {}, "ep1", map[EffectKey]int64{}},
		{"same effect and scope summed", func(t *testing.T, s *Store) {
			mustAppendEffect(t, s, spendIntent("ep1", "FS.Write", "/a", 2))
			mustAppendEffect(t, s, spendIntent("ep1", "FS.Write", "/a", 3))
		}, "ep1", map[EffectKey]int64{{"FS.Write", "/a"}: 5}},
		{"different scopes and effects separated", func(t *testing.T, s *Store) {
			mustAppendEffect(t, s, spendIntent("ep1", "FS.Write", "/a", 2))
			mustAppendEffect(t, s, spendIntent("ep1", "FS.Write", "/b", 3))
			mustAppendEffect(t, s, spendIntent("ep1", "Proc.Run", "/a", 4))
		}, "ep1", map[EffectKey]int64{{"FS.Write", "/a"}: 2, {"FS.Write", "/b"}: 3, {"Proc.Run", "/a"}: 4}},
		{"other episodes excluded incl ep1 vs ep10 prefix trap", func(t *testing.T, s *Store) {
			mustAppendEffect(t, s, spendIntent("ep1", "FS.Write", "/a", 1))
			mustAppendEffect(t, s, spendIntent("ep10", "FS.Write", "/a", 100))
			mustAppendEffect(t, s, spendIntent("ep2", "FS.Write", "/a", 1000))
			mustAppendEffect(t, s, spendIntent("ep", "FS.Write", "/a", 10000))
		}, "ep1", map[EffectKey]int64{{"FS.Write", "/a"}: 1}},
		{"ep10 not polluted by ep1", func(t *testing.T, s *Store) {
			mustAppendEffect(t, s, spendIntent("ep1", "FS.Write", "/a", 1))
			mustAppendEffect(t, s, spendIntent("ep10", "FS.Write", "/a", 100))
		}, "ep10", map[EffectKey]int64{{"FS.Write", "/a"}: 100}},
		{"invocation intent excluded", func(t *testing.T, s *Store) {
			mustAppendEffect(t, s, spendIntent("ep1", "FS.Write", "/a", 1))
			if _, _, err := s.AppendIntent(boundedTestContext(t), "ep1", JournalIntent{
				InvocationID: "ep1", WorldRef: ref("w"), EntryHash: ref("e"), PrevEntryHash: ref("p"),
				TransitionFn: ref("f"), TransitionRef: ref("t"), Interpreter: ref("i"),
			}); err != nil {
				t.Fatal(err)
			}
		}, "ep1", map[EffectKey]int64{{"FS.Write", "/a"}: 1}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := openMem(t)
			tc.setup(t, s)
			assertSpend(t, s, tc.ep, tc.want)
		})
	}
}

func TestEffectSpendRequiresDeadline(t *testing.T) {
	s := openMem(t)
	if _, err := s.EffectSpend(nil, "ep1"); err == nil { //nolint:staticcheck // nil ctx is the case under test
		t.Fatal("nil ctx accepted")
	}
}

func TestEffectSpendSurvivesReopenAndReadOnly(t *testing.T) {
	path := filepath.Join(t.TempDir(), "world.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	mustAppendEffect(t, s, spendIntent("ep1", "FS.Write", "/a", 2))
	mustAppendEffect(t, s, spendIntent("ep1", "FS.Write", "/a", 3))
	mustAppendEffect(t, s, spendIntent("ep10", "FS.Write", "/a", 50))
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	want := map[EffectKey]int64{{"FS.Write", "/a"}: 5}
	s2, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	assertSpend(t, s2, "ep1", want)
	if err := s2.Close(); err != nil {
		t.Fatal(err)
	}
	ro, err := OpenReadOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	defer ro.Close()
	assertSpend(t, ro, "ep1", want)
}
