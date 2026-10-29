package surface

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestFollowLegacyV1FixturesStillDecode(t *testing.T) {
	for _, name := range []string{goldenPath, rigGoldenPath} {
		b, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		ix, err := Decode(b)
		if err != nil {
			t.Fatal(err)
		}
		if err := Validate(ix); err != nil {
			t.Fatal(err)
		}
		if _, err := DecodeFollow(b); err == nil || !strings.Contains(err.Error(), "schema") {
			t.Fatalf("%s accepted or schema not named: %v", name, err)
		}
	}
}

func followFixture(t *testing.T) FollowIndex {
	t.Helper()
	ix := goldenIndex(t)
	return FollowIndex{Schema: FollowSchema, Generation: "g1", Repos: []FollowPartition{
		{ID: "world", Head: ix.Source.Head, BuiltAt: ix.BuiltAt, CheckedAt: ix.BuiltAt, Index: ix},
		{ID: "tools", Head: ix.Source.Head, BuiltAt: ix.BuiltAt, CheckedAt: ix.BuiltAt, Index: ix},
	}}
}

func TestFollowEnvelopeValidates(t *testing.T) {
	ix := followFixture(t)
	b, err := json.Marshal(ix)
	if err != nil {
		t.Fatal(err)
	}
	got, err := DecodeFollow(b)
	if err != nil || !reflect.DeepEqual(ix, got) {
		t.Fatalf("round trip: %v", err)
	}
	if _, err := Decode(b); err == nil {
		t.Fatal("legacy Decode accepted envelope")
	}
}

func TestFollowEnvelopeRejectsMalformedNamespaces(t *testing.T) {
	mutations := map[string]func(*FollowIndex){
		"slug":       func(x *FollowIndex) { x.Repos[0].ID = "Bad" },
		"duplicate":  func(x *FollowIndex) { x.Repos[1].ID = x.Repos[0].ID },
		"zero":       func(x *FollowIndex) { x.Repos = nil },
		"five":       func(x *FollowIndex) { x.Repos = append(x.Repos, x.Repos[0], x.Repos[0], x.Repos[0]) },
		"head":       func(x *FollowIndex) { x.Repos[0].Head = strings.Repeat("f", 40) },
		"v1":         func(x *FollowIndex) { x.Repos[0].Index.Schema = "bad" },
		"built":      func(x *FollowIndex) { x.Repos[0].BuiltAt = "yesterday" },
		"checked":    func(x *FollowIndex) { x.Repos[0].CheckedAt = "" },
		"failure":    func(x *FollowIndex) { x.Repos[0].LastFailure = &FollowFailure{Class: "source", Text: "/secret"} },
		"generation": func(x *FollowIndex) { x.Generation = "" },
		"sum": func(x *FollowIndex) {
			merges, bases := []Merge{}, map[string]string{}
			for i := 0; i < 1300; i++ {
				m := Merge{Hash: fmt.Sprintf("%040x", 3*i+3),
					Parents: []string{fmt.Sprintf("%040x", 3*i+1), fmt.Sprintf("%040x", 3*i+2)},
					At:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
				merges = append([]Merge{m}, merges...)
				bases[m.Hash] = fmt.Sprintf("%040x", 3*i+1)
			}
			ix, err := Build(BuildInput{Merges: merges, Bases: bases, BuiltAt: time.Now(), Rev: "dev"})
			if err != nil || Validate(ix) != nil {
				t.Fatalf("large valid partition: %v", err)
			}
			for i := range x.Repos {
				x.Repos[i].Index, x.Repos[i].Head = ix, ix.Source.Head
			}
		},
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			x := followFixture(t)
			mutate(&x)
			b, _ := json.Marshal(x)
			if _, err := DecodeFollow(b); err == nil {
				t.Fatal("malformed envelope accepted")
			}
		})
	}
	b, _ := json.Marshal(followFixture(t))
	for _, raw := range [][]byte{
		append([]byte(`{"unknown":true,`), b[1:]...),
		append([]byte(`{"overlay":[],`), b[1:]...),
		append(b, []byte("{}")...),
		[]byte(strings.Repeat(" ", MaxBytes+1)),
	} {
		if _, err := DecodeFollow(raw); err == nil {
			t.Fatal("unknown field, overlay, trailing or oversized data accepted")
		}
	}
}

func TestQualifiedIDs(t *testing.T) {
	hash := strings.Repeat("a", 40)
	for _, pair := range [][2]string{{"commit", hash}, {"commit", strings.Repeat("b", 64)}, {"episode", "w220a"}, {"brane", "row-220"}, {"brane", "pr:7"}} {
		a, err := QualifyID("world", pair[0], pair[1])
		if err != nil {
			t.Fatal(err)
		}
		r, k, l, err := ParseQualifiedID(a)
		if err != nil || r != "world" || k != pair[0] || l != pair[1] {
			t.Fatalf("round trip %s: %v", a, err)
		}
		b, _ := QualifyID("tools", k, l)
		if a == b {
			t.Fatal("repo collision")
		}
	}
	for _, v := range [][3]string{
		{"Bad", "episode", "ep"}, {"a", "unknown", "x"}, {"a", "brane", ""},
		{"a", "commit", "abc"}, {"a", "commit", strings.Repeat("A", 40)},
		{"a", "commit", hash + ":x"}, {"a", "episode", "ep:x"},
		{"", "brane", "x"}, {strings.Repeat("a", 33), "brane", "x"},
	} {
		if _, err := QualifyID(v[0], v[1], v[2]); err == nil {
			t.Fatalf("accepted %v", v)
		}
	}
	for _, s := range []string{"x", "a:episode", "a:episode:", "a:no:x", "A:brane:x", "a:commit:abc"} {
		if _, _, _, err := ParseQualifiedID(s); err == nil {
			t.Fatalf("accepted %s", s)
		}
	}
}
