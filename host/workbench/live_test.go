package workbench

import (
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// Reflection keeps these feature oracles executable against the pre-view API.
func livePage(t *testing.T, cursor int64) Page {
	t.Helper()
	p := Page{Title: "Live"}
	v := reflect.ValueOf(&p).Elem().FieldByName("Live")
	if !v.IsValid() {
		t.Fatal("Page.Live is missing")
	}
	v.FieldByName("Cursor").SetInt(cursor)
	v.FieldByName("Head").SetString("head-seven")
	if cursor >= 0 {
		v.FieldByName("Recent").Set(reflect.ValueOf([]EntryView{{EntryIndex: 7, EntryHash: "hash-seven", SelectHref: "?from=0&entry=7"}, {EntryIndex: 6, EntryHash: "hash-six", SelectHref: "?from=0&entry=6"}}))
	}
	return p
}
func liveSection(t *testing.T, body string) string {
	t.Helper()
	start := strings.Index(body, `<section aria-label="live"`)
	if start < 0 {
		t.Fatal("live section missing")
	}
	end := strings.Index(body[start:], "</section>")
	if end < 0 {
		t.Fatal("live section unterminated")
	}
	return body[start : start+end+len("</section>")]
}
func TestTokensBothThemes(t *testing.T) {
	body := renderPage(t, Page{})
	style := regexp.MustCompile(`(?s)<style>(.*?)</style>`).FindStringSubmatch(body)
	if len(style) != 2 {
		t.Fatal("style block missing")
	}
	roots := regexp.MustCompile(`:root\s*\{([^}]+)\}`).FindAllStringSubmatch(style[1], -1)
	if len(roots) != 2 || !strings.Contains(style[1], "@media (prefers-color-scheme: dark)") {
		t.Fatal("light and dark token roots missing")
	}
	want := []string{"--bg", "--card", "--ink", "--mut", "--line", "--acc", "--accbg", "--ok", "--okbg", "--warn", "--warnbg"}
	sort.Strings(want)
	for i, root := range roots {
		var names []string
		for _, m := range regexp.MustCompile(`(--[a-z]+)\s*:`).FindAllStringSubmatch(root[1], -1) {
			names = append(names, m[1])
		}
		sort.Strings(names)
		if len(names) == 0 || !reflect.DeepEqual(names, want) {
			t.Errorf("theme %d tokens=%v want %v", i, names, want)
		}
		bg := []string{"#faf9f5", "#1f1e1b"}[i]
		if !regexp.MustCompile(`--bg\s*:\s*` + bg + `\s*[;}]?`).MatchString(root[1]) {
			t.Errorf("theme %d background missing %s", i, bg)
		}
	}
	if strings.Contains(body, "<link") {
		t.Fatal("external stylesheet")
	}
}
func TestRenderLiveRegion(t *testing.T) {
	for _, cursor := range []int64{7, -1} {
		t.Run(strconv.FormatInt(cursor, 10), func(t *testing.T) {
			s := liveSection(t, renderPage(t, livePage(t, cursor)))
			if cursor == 7 {
				if !strings.Contains(s, `data-live-cursor="7"`) || !strings.Contains(s, "Log at entry 7 (head head-seven) when this page was rendered") {
					t.Errorf("live position missing: %s", s)
				}
				a, b := strings.Index(s, "select entry 7"), strings.Index(s, "select entry 6")
				if a < 0 || b <= a {
					t.Errorf("newest-first rows missing: %s", s)
				}
				if !strings.Contains(s, "hash-seven") || !strings.Contains(s, "hash-six") {
					t.Error("row hashes missing")
				}
			} else if !strings.Contains(s, `data-live-cursor="-1"`) || !strings.Contains(s, "The log is empty") {
				t.Errorf("empty-log assertion missing: %s", s)
			}
			if strings.Contains(s, "<h3>entry ") {
				t.Error("live rows use timeline headings")
			}
		})
	}
}
func TestRenderQuietFooter(t *testing.T) {
	for _, cursor := range []int64{7, -1} {
		body := renderPage(t, livePage(t, cursor))
		start := strings.Index(body, `<footer aria-label="live status"`)
		if start < strings.Index(body, "</main>") || start < 0 {
			t.Fatal("footer missing or inside main")
		}
		footer := body[start:]
		if !strings.Contains(footer, `role="status"`) || !strings.Contains(footer, `<span data-live-status>Live updates: off. Reload to refresh.</span>`) {
			t.Error("no-JS status missing")
		}
		if strings.Contains(footer, "Log at entry") || strings.Contains(footer, "The log is empty") {
			t.Error("log position belongs in swapped live section")
		}
	}
}
func TestProvenanceWalkStaysLast(t *testing.T) {
	body := renderPage(t, Page{})
	walk := strings.Index(body, `<section aria-label="provenance walk">`)
	for _, label := range []string{"live", "world graph", "decisions"} {
		needle := `<section aria-label="` + label + `"`
		if strings.Count(body, needle) != 1 || strings.Index(body, needle) >= walk {
			t.Errorf("%s must occur once before provenance walk", label)
		}
	}
	close := strings.Index(body, "</section>\n</main>")
	if strings.Count(body, "</section>\n</main>") != 1 || close < walk || strings.Contains(body[walk:close], "<section") && strings.Count(body[walk:close], "<section") != 1 {
		t.Error("provenance walk must close main")
	}
}
func TestWorkbenchNewViewTypesInCensus(t *testing.T) {
	listed := map[reflect.Type]bool{}
	for _, typ := range workbenchViewTypes {
		listed[typ] = true
	}
	visited := map[reflect.Type]bool{}
	var visit func(reflect.Type)
	visit = func(typ reflect.Type) {
		for typ.Kind() == reflect.Pointer || typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array {
			typ = typ.Elem()
		}
		if typ.Kind() != reflect.Struct || visited[typ] {
			return
		}
		visited[typ] = true
		if !listed[typ] {
			t.Errorf("new view type %s missing from census", typ.Name())
		}
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			if f.IsExported() {
				visit(f.Type)
			}
		}
	}
	page := reflect.TypeOf(Page{})
	for _, name := range []string{"Live", "Graph", "Decisions"} {
		if f, ok := page.FieldByName(name); ok {
			visit(f.Type)
		}
	}
	if len(visited) == 0 {
		t.Fatal("no new view types found")
	}
}
