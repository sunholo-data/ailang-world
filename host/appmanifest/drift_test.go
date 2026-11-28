package appmanifest

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// The drift test binds manifest.go to design_docs/sketches/appmanifest.ail: it reads the
// sketch, parses every inline `tests [...]` row of every exported law, and evaluates each
// row through the Go mirror of the same name. A law without a mirror, a mirror without a
// law, a moved bound on either side, or a row the mirror disagrees with turns it red.
// This is the execargs pattern (host/broker/execargs_drift_test.go), cut down to the
// literal subset this sketch uses: ints, strings, booleans and tuples.

type sketchRow struct {
	line int
	args []any
	want any
}

var mirrors = map[string]func(t *testing.T, a []any) any{
	"schemaOk":       func(t *testing.T, a []any) any { return schemaOk(str(t, a, 0)) },
	"nameLenOk":      func(t *testing.T, a []any) any { return nameLenOk(int(num(t, a, 0))) },
	"nameByteOk":     func(t *testing.T, a []any) any { return nameByteOk(int(num(t, a, 0)), boolean(t, a, 1)) },
	"nameReserved":   func(t *testing.T, a []any) any { return nameReserved(str(t, a, 0)) },
	"titleLenOk":     func(t *testing.T, a []any) any { return titleLenOk(int(num(t, a, 0))) },
	"summaryLenOk":   func(t *testing.T, a []any) any { return summaryLenOk(int(num(t, a, 0))) },
	"agentOk":        func(t *testing.T, a []any) any { return agentOk(str(t, a, 0)) },
	"modeOk":         func(t *testing.T, a []any) any { return modeOk(str(t, a, 0)) },
	"deadlineOk":     func(t *testing.T, a []any) any { return deadlineOk(num(t, a, 0)) },
	"ttlOk":          func(t *testing.T, a []any) any { return ttlOk(num(t, a, 0)) },
	"effectAllowed":  func(t *testing.T, a []any) any { return effectAllowed(str(t, a, 0)) },
	"scopeOk":        func(t *testing.T, a []any) any { return scopeOk(str(t, a, 0)) },
	"budgetOk":       func(t *testing.T, a []any) any { return budgetOk(num(t, a, 0)) },
	"grantNarrows":   func(t *testing.T, a []any) any { return grantNarrows(num(t, a, 0), num(t, a, 1)) },
	"grantCountOk":   func(t *testing.T, a []any) any { return grantCountOk(int(num(t, a, 0))) },
	"listCountOk":    func(t *testing.T, a []any) any { return listCountOk(int(num(t, a, 0))) },
	"checkNameLenOk": func(t *testing.T, a []any) any { return checkNameLenOk(int(num(t, a, 0))) },
	"checkTimeoutOk": func(t *testing.T, a []any) any { return checkTimeoutOk(num(t, a, 0)) },
	"outputKindOk":   func(t *testing.T, a []any) any { return outputKindOk(str(t, a, 0)) },
	"noDotDotSegment": func(t *testing.T, a []any) any {
		return noDotDotSegment(str(t, a, 0))
	},
	"outputGlobOk": func(t *testing.T, a []any) any { return outputGlobOk(str(t, a, 0)) },
	"viewOk":       func(t *testing.T, a []any) any { return viewOk(str(t, a, 0)) },
	"viewCountOk":  func(t *testing.T, a []any) any { return viewCountOk(int(num(t, a, 0))) },
}

func sketchPath() string {
	return filepath.Join("..", "..", "design_docs", "sketches", "appmanifest.ail")
}

func TestAppManifestSketchDrift(t *testing.T) {
	src, err := os.ReadFile(sketchPath())
	if err != nil {
		t.Fatalf("read sketch: %v", err)
	}
	if !strings.HasPrefix(string(src), "module sketches/appmanifest\n") {
		t.Fatal("sketch module header moved")
	}
	laws := parseSketchTests(t, string(src))

	var sketchNames, goNames []string
	for n := range laws {
		sketchNames = append(sketchNames, n)
	}
	for n := range mirrors {
		goNames = append(goNames, n)
	}
	sort.Strings(sketchNames)
	sort.Strings(goNames)
	if !reflect.DeepEqual(sketchNames, goNames) {
		t.Fatalf("sketch laws with tests = %v\nGo mirrors            = %v\nevery sketch law needs a Go mirror and vice versa", sketchNames, goNames)
	}

	total := 0
	for _, name := range sketchNames {
		distinct := map[any]bool{}
		for _, row := range laws[name] {
			total++
			if got := mirrors[name](t, row.args); got != row.want {
				t.Errorf("sketch line %d %s%v: Go mirror = %#v, sketch expects %#v", row.line, name, row.args, got, row.want)
			}
			distinct[row.want] = true
		}
		// A predicate needs a true row and a false row, or a constant mirror would pass.
		if len(distinct) < 2 {
			t.Errorf("sketch law %s has %d distinct expected value(s); it needs at least 2", name, len(distinct))
		}
	}
	// Instrument-health floor: the rows the sketch carries today.
	if total < 122 {
		t.Fatalf("parsed %d sketch rows; expected at least 122", total)
	}
	t.Logf("evaluated %d sketch rows over %d laws", total, len(sketchNames))
}

// TestAllowedEffectsAreTheSketchList keeps effectAllowed's Go set equal to the names the
// sketch's ensures lists, so neither side can grow an effect alone. allowed_effects_test.go
// binds the same set to the broker's effect constants.
func TestAllowedEffectsAreTheSketchList(t *testing.T) {
	src, err := os.ReadFile(sketchPath())
	if err != nil {
		t.Fatal(err)
	}
	s := string(src)
	start := strings.Index(s, "export func effectAllowed")
	end := strings.Index(s[start:], "\ntests [")
	if start < 0 || end < 0 {
		t.Fatal("effectAllowed law not found in the sketch")
	}
	var fromSketch []string
	for _, m := range regexp.MustCompile(`e == "([A-Za-z.]+)"`).FindAllStringSubmatch(s[start:start+end], -1) {
		fromSketch = append(fromSketch, m[1])
	}
	var fromGo []string
	for e := range allowedEffects {
		fromGo = append(fromGo, e)
	}
	sort.Strings(fromSketch)
	sort.Strings(fromGo)
	if !reflect.DeepEqual(fromSketch, fromGo) {
		t.Fatalf("effectAllowed: sketch ensures lists %v, Go allows %v", fromSketch, fromGo)
	}
}

var exportFuncRE = regexp.MustCompile(`(?m)^export func ([A-Za-z0-9_]+)\(`)

func parseSketchTests(t *testing.T, src string) map[string][]sketchRow {
	t.Helper()
	out := map[string][]sketchRow{}
	heads := exportFuncRE.FindAllStringSubmatchIndex(src, -1)
	for i, h := range heads {
		name := src[h[2]:h[3]]
		end := len(src)
		if i+1 < len(heads) {
			end = heads[i+1][0]
		}
		at := strings.Index(src[h[0]:end], "\ntests [")
		if at < 0 {
			continue
		}
		p := &parser{src: src, pos: h[0] + at + len("\ntests "), t: t}
		list, ok := p.value().([]any)
		if !ok {
			t.Fatalf("%s: tests is not a list", name)
		}
		for j, item := range list {
			pair, ok := item.(tuple)
			if !ok || len(pair) != 2 {
				t.Fatalf("%s row %d is not an (input, expected) pair", name, j)
			}
			args := []any{pair[0]}
			if tup, ok := pair[0].(tuple); ok {
				args = []any(tup)
			}
			out[name] = append(out[name], sketchRow{line: p.rowLines[j], args: args, want: pair[1]})
		}
	}
	return out
}

type tuple []any

type parser struct {
	src      string
	pos      int
	depth    int
	rowLines []int
	t        *testing.T
}

func (p *parser) fail(format string, args ...any) {
	p.t.Helper()
	p.t.Fatalf("sketch parse at line %d: "+format, append([]any{strings.Count(p.src[:p.pos], "\n") + 1}, args...)...)
}

func (p *parser) skip() {
	for p.pos < len(p.src) {
		switch {
		case strings.HasPrefix(p.src[p.pos:], "--"):
			for p.pos < len(p.src) && p.src[p.pos] != '\n' {
				p.pos++
			}
		case strings.ContainsRune(" \t\r\n", rune(p.src[p.pos])):
			p.pos++
		default:
			return
		}
	}
}

func (p *parser) seq(closer byte, item func()) {
	p.skip()
	if p.src[p.pos] == closer {
		p.pos++
		return
	}
	for {
		item()
		p.skip()
		switch p.src[p.pos] {
		case ',':
			p.pos++
		case closer:
			p.pos++
			return
		default:
			p.fail("want ',' or %q", closer)
		}
	}
}

func (p *parser) value() any {
	p.skip()
	if p.pos >= len(p.src) {
		p.fail("unexpected end")
	}
	switch c := p.src[p.pos]; {
	case c == '[':
		p.pos++
		outer := p.depth == 0
		p.depth++
		list := []any{}
		p.seq(']', func() {
			if outer {
				p.skip()
				p.rowLines = append(p.rowLines, strings.Count(p.src[:p.pos], "\n")+1)
			}
			list = append(list, p.value())
		})
		p.depth--
		return list
	case c == '(':
		p.pos++
		p.depth++
		tup := tuple{}
		p.seq(')', func() { tup = append(tup, p.value()) })
		p.depth--
		if len(tup) == 1 {
			return tup[0]
		}
		return tup
	case c == '"':
		p.pos++
		var b strings.Builder
		for {
			if p.pos >= len(p.src) || p.src[p.pos] == '\n' {
				p.fail("unterminated string")
			}
			ch := p.src[p.pos]
			if ch == '"' {
				p.pos++
				return b.String()
			}
			if ch == '\\' {
				if p.pos+1 >= len(p.src) || (p.src[p.pos+1] != '\\' && p.src[p.pos+1] != '"') {
					p.fail("escape outside the parsed subset")
				}
				b.WriteByte(p.src[p.pos+1])
				p.pos += 2
				continue
			}
			b.WriteByte(ch)
			p.pos++
		}
	case c == '-' || (c >= '0' && c <= '9'):
		start := p.pos
		p.pos++
		for p.pos < len(p.src) && p.src[p.pos] >= '0' && p.src[p.pos] <= '9' {
			p.pos++
		}
		n, err := strconv.ParseInt(p.src[start:p.pos], 10, 64)
		if err != nil {
			p.fail("bad int %q", p.src[start:p.pos])
		}
		return n
	case strings.HasPrefix(p.src[p.pos:], "true"):
		p.pos += 4
		return true
	case strings.HasPrefix(p.src[p.pos:], "false"):
		p.pos += 5
		return false
	}
	p.fail("unsupported literal starting %q", p.src[p.pos])
	return nil
}

func at(t *testing.T, a []any, i int) any {
	t.Helper()
	if i >= len(a) {
		t.Fatalf("sketch row has %d args, mirror reads arg %d", len(a), i)
	}
	return a[i]
}

func num(t *testing.T, a []any, i int) int64 {
	t.Helper()
	v, ok := at(t, a, i).(int64)
	if !ok {
		t.Fatalf("arg %d is %T, want int", i, a[i])
	}
	return v
}

func str(t *testing.T, a []any, i int) string {
	t.Helper()
	v, ok := at(t, a, i).(string)
	if !ok {
		t.Fatalf("arg %d is %T, want string", i, a[i])
	}
	return v
}

func boolean(t *testing.T, a []any, i int) bool {
	t.Helper()
	v, ok := at(t, a, i).(bool)
	if !ok {
		t.Fatalf("arg %d is %T, want bool", i, a[i])
	}
	return v
}
