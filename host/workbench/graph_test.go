package workbench

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

var updateWorkbenchGolden = flag.Bool("update-workbench-golden", false, "emit graph goldens from TempDir")

func graphFixture(n int) []GraphEntry {
	var entries []GraphEntry
	for i := 0; i < n; i++ {
		entries = append(entries, GraphEntry{Index: int64(9 - i), Href: fmt.Sprintf("?from=0&entry=%d", 9-i), Edges: []EdgeView{
			{Relation: "transitionFn", Target: "shared", Available: true, Href: "?object=shared"},
			{Relation: "interpreter", Target: fmt.Sprintf("interpreter-%d", i), Available: true, Href: fmt.Sprintf("?object=interpreter-%d", i)},
			{Relation: "transitionRef", Target: fmt.Sprintf("missing-%d", i), Missing: "not stored"},
		}})
	}
	return entries
}
func graphSection(t *testing.T, body string) string {
	t.Helper()
	start := strings.Index(body, `<section aria-label="world graph">`)
	if start < 0 {
		t.Fatal("graph section missing")
	}
	end := strings.Index(body[start:], "</section>")
	if end < 0 {
		t.Fatal("graph section unclosed")
	}
	return body[start : start+end+10]
}
func graphRender(t *testing.T, g GraphView) string {
	return graphSection(t, renderPage(t, Page{Graph: g}))
}
func TestGraphGolden(t *testing.T) {
	for _, f := range []struct {
		name string
		n    int
	}{{"graph_home.golden", 3}, {"graph_selected.golden", 1}} {
		t.Run(f.name, func(t *testing.T) {
			g := LayoutGraph(graphFixture(f.n))
			got := []byte(graphRender(t, g))
			if *updateWorkbenchGolden {
				path := filepath.Join(t.TempDir(), f.name)
				if err := os.WriteFile(path, got, 0600); err != nil {
					t.Fatal(err)
				}
				b, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				t.Logf("WORLD_GRAPH_GOLDEN %s %x %s", f.name, sha256.Sum256(b), base64.StdEncoding.EncodeToString(b))
				return
			}
			want, err := os.ReadFile(filepath.Join("testdata", f.name))
			if err != nil {
				t.Fatal(err)
			}
			if len(want) == 0 || !strings.Contains(string(want), "<svg") || strings.Count(string(want), "<rect") != len(g.Nodes) {
				t.Fatal("golden is empty or has invalid SVG/node census")
			}
			if string(want) != string(got) {
				t.Fatal("graph differs from golden")
			}
		})
	}
}
func TestGraphDeterministic(t *testing.T) {
	entries := graphFixture(3)
	base := LayoutGraph(entries)
	rendered := graphRender(t, base)
	want := []string{"entry 9", "entry 8", "entry 7", "shared", "interpreter-0", "missing-0", "interpreter-1", "missing-1", "interpreter-2", "missing-2"}
	for i := 0; i < 32; i++ {
		var copied []GraphEntry
		b, _ := json.Marshal(entries)
		if err := json.Unmarshal(b, &copied); err != nil {
			t.Fatal(err)
		}
		g := LayoutGraph(copied)
		var labels []string
		for _, n := range g.Nodes {
			labels = append(labels, n.Label)
		}
		if !reflect.DeepEqual(labels, want) {
			t.Fatalf("first-appearance node order=%v want %v", labels, want)
		}
		if !reflect.DeepEqual(base, g) || rendered != graphRender(t, g) {
			t.Fatal("graph is nondeterministic")
		}
	}
}
func TestGraphNodesDisjointAndInBounds(t *testing.T) {
	for count := 1; count <= 5; count++ {
		g := LayoutGraph(graphFixture(count))
		rows := 1 + 2*count
		if g.Height != 32+36*rows {
			t.Fatalf("height=%d want %d", g.Height, 32+36*rows)
		}
		for i, n := range g.Nodes {
			row := i
			if i >= count {
				row = i - count
			}
			x := 16
			if i >= count {
				x = 344
			}
			if n.X != x || n.Y != 16+36*row || n.X+280 > 640 || n.Y+28 > g.Height {
				t.Fatalf("geometry node=%+v", n)
			}
			for j, m := range g.Nodes {
				if j < i && n.X < m.X+280 && m.X < n.X+280 && n.Y < m.Y+28 && m.Y < n.Y+28 {
					t.Fatal("graph boxes intersect")
				}
			}
		}
		s := graphRender(t, g)
		if !strings.Contains(s, fmt.Sprintf(`viewBox="0 0 640 %d"`, g.Height)) || strings.Count(s, `width="280" height="28"`) != len(g.Nodes) {
			t.Fatal("rendered geometry differs")
		}
	}
}
func TestGraphUnavailableNodeHasNoLink(t *testing.T) {
	g := LayoutGraph(graphFixture(1))
	s := graphRender(t, g)
	available := 0
	for _, n := range g.Nodes {
		if n.Available {
			available++
		}
	}
	if strings.Count(s, "<a ") != available || !strings.Contains(s, "UNAVAILABLE: missing-0") {
		t.Fatal("available anchors or unavailable label incorrect")
	}
	if strings.Contains(s, `href="/workbench?object=missing-0"`) {
		t.Fatal("unstored target linked")
	}
}
func TestGraphCarriesNoGrade(t *testing.T) {
	s := graphRender(t, LayoutGraph(graphFixture(1)))
	if !strings.Contains(s, "<svg") {
		t.Fatal("SVG missing")
	}
	for _, label := range []string{"PROVEN", "TESTED", "ATTESTED", "CLAIMED"} {
		if strings.Contains(s, label) {
			t.Fatalf("graph carries grade %s", label)
		}
	}
}
