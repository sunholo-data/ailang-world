package broker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"
)

// The drift test binds execargs.go to design_docs/sketches/execargs.ail
// (row 140 M1, w-workspace-exec-toolchain-effect §4.2-§4.3): it READS the
// sketch, parses every inline `tests [...]` row of every law, and evaluates
// each row through the Go mirror of the same name. A law added to the sketch
// with no Go mirror, a mirror with no sketch law, a changed boundary (256
// bytes, 10^6, max_args) on either side, or a row whose expected value the Go
// mirror disagrees with, turns this test red. This is the effectplan pattern
// (host/coordinator/effectplan_drift_test.go); the parser is repeated here
// because it lives in another package's test files, and it adds the string
// escapes the newline and NUL rows need.

type execSketchTuple []any

type execSketchRow struct {
	line int
	args []any
	want any
}

var execSketchMirrors = map[string]func(t *testing.T, a []any) any{
	"noDotDotSegment": func(t *testing.T, a []any) any { return execNoDotDotSegment(xString(t, a, 0)) },
	"pathOk":          func(t *testing.T, a []any) any { return execPathOk(xString(t, a, 0)) },
	"flagNameOf":      func(t *testing.T, a []any) any { return flagNameOf(xString(t, a, 0)) },
	"flagValueOf":     func(t *testing.T, a []any) any { return flagValueOf(xString(t, a, 0)) },
	"flagRuleMatches": func(t *testing.T, a []any) any {
		r := xRule(t, a, 0)
		return flagRuleMatches(r.name, xString(t, a, 1))
	},
	"consumesNext": func(t *testing.T, a []any) any { return consumesNext(xString(t, a, 0), xBool(t, a, 1)) },
	"digitByteOk":  func(t *testing.T, a []any) any { return digitByteOk(int(xInt(t, a, 0))) },
	"intValueOk":   func(t *testing.T, a []any) any { return intValueOk(int(xInt(t, a, 0))) },
	"regexValueOk": func(t *testing.T, a []any) any { return regexValueOk(xString(t, a, 0), int(xInt(t, a, 1))) },
	"flagValueOk": func(t *testing.T, a []any) any {
		return flagValueOk(xString(t, a, 0), xString(t, a, 1), int(xInt(t, a, 2)), xBool(t, a, 3), xBool(t, a, 4))
	},
	"pkgpatternOk": func(t *testing.T, a []any) any { return pkgpatternOk(xString(t, a, 0)) },
	"testfileOk":   func(t *testing.T, a []any) any { return testfileOk(xString(t, a, 0), xString(t, a, 1)) },
	"positionalOk": func(t *testing.T, a []any) any {
		return positionalOk(xString(t, a, 0), xString(t, a, 1), xBool(t, a, 2))
	},
	"dashDashOk":     func(t *testing.T, a []any) any { return dashDashOk(xString(t, a, 0), xBool(t, a, 1)) },
	"argCountOk":     func(t *testing.T, a []any) any { return argCountOk(int(xInt(t, a, 0)), int(xInt(t, a, 1))) },
	"intValueParses": func(t *testing.T, a []any) any { return intValueParses(xString(t, a, 0)) },
	"enumHas":        func(t *testing.T, a []any) any { return enumHas(xString(t, a, 0), xString(t, a, 1)) },
	"flagClassOf": func(t *testing.T, a []any) any {
		return flagClassOf(xFlags(t, xList(t, a, 0)), xString(t, a, 1))
	},
	"matchArgs": func(t *testing.T, a []any) any {
		args := []string{}
		for i := range xList(t, a, 1) {
			args = append(args, xString(t, xList(t, a, 1), i))
		}
		out, err := MatchExecArgs(xGrammar(t, a, 0), args)
		return renderExecMatch(t, out, err)
	},
}

// renderExecMatch is the sketch's scalar rendering of a match: "ok <the
// normalized args as a JSON array>" or "refused <index> <reason>".
func renderExecMatch(t *testing.T, out []string, err error) string {
	t.Helper()
	if err != nil {
		var bad *ExecArgError
		if !asExecArgError(err, &bad) {
			t.Fatalf("MatchExecArgs error %T %v is not an *ExecArgError", err, err)
		}
		return fmt.Sprintf("refused %d %s", bad.Index, bad.Reason)
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(out); err != nil {
		t.Fatal(err)
	}
	return "ok " + strings.TrimSuffix(buf.String(), "\n")
}

func asExecArgError(err error, target **ExecArgError) bool {
	e, ok := err.(*ExecArgError)
	if ok {
		*target = e
	}
	return ok
}

func TestExecArgsSketchDrift(t *testing.T) {
	src, err := os.ReadFile(filepath.Join("..", "..", "design_docs", "sketches", "execargs.ail"))
	if err != nil {
		t.Fatalf("read sketch: %v", err)
	}
	if !strings.HasPrefix(string(src), "module sketches/execargs\n") {
		t.Fatal("sketch module header moved")
	}
	laws := parseExecSketchTests(t, string(src))

	var sketchNames, goNames []string
	for name := range laws {
		sketchNames = append(sketchNames, name)
	}
	for name := range execSketchMirrors {
		goNames = append(goNames, name)
	}
	sort.Strings(sketchNames)
	sort.Strings(goNames)
	if !reflect.DeepEqual(sketchNames, goNames) {
		t.Fatalf("sketch laws with tests = %v\nGo mirrors            = %v\nevery sketch law needs a Go mirror and vice versa", sketchNames, goNames)
	}

	total := 0
	for _, name := range sketchNames {
		rows := laws[name]
		// A law needs at least two DIFFERENT expected values (for a predicate:
		// a true row and a false row), or a constant mirror would pass it.
		distinct := map[any]bool{}
		for _, row := range rows {
			total++
			got := execSketchMirrors[name](t, row.args)
			if got != row.want {
				t.Errorf("sketch line %d %s%q: Go mirror = %#v, sketch expects %#v", row.line, name, row.args, got, row.want)
			}
			distinct[row.want] = true
		}
		if len(distinct) < 2 {
			t.Errorf("sketch law %s has %d distinct expected value(s); it needs at least 2", name, len(distinct))
		}
	}
	// Instrument-health floor: the parser found the rows the sketch carries
	// today (a parser that silently found nothing would pass the loop above).
	if total < 130 {
		t.Fatalf("parsed %d sketch rows; expected at least 130", total)
	}
	t.Logf("evaluated %d sketch rows over %d laws", total, len(sketchNames))
}

var execExportFuncRE = regexp.MustCompile(`(?m)^export func ([A-Za-z0-9_]+)\(`)

// parseExecSketchTests returns, per exported function that carries inline
// tests, its rows. A function's tests block is the first `tests [` between its
// header and the next exported function.
func parseExecSketchTests(t *testing.T, src string) map[string][]execSketchRow {
	t.Helper()
	out := map[string][]execSketchRow{}
	heads := execExportFuncRE.FindAllStringSubmatchIndex(src, -1)
	for i, h := range heads {
		name := src[h[2]:h[3]]
		end := len(src)
		if i+1 < len(heads) {
			end = heads[i+1][0]
		}
		seg := src[h[0]:end]
		at := strings.Index(seg, "\ntests [")
		if at < 0 {
			continue
		}
		p := &execSketchParser{src: src, pos: h[0] + at + len("\ntests "), t: t}
		list, ok := p.value().([]any)
		if !ok {
			t.Fatalf("%s: tests is not a list", name)
		}
		for j, item := range list {
			pair, ok := item.(execSketchTuple)
			if !ok || len(pair) != 2 {
				t.Fatalf("%s row %d is not an (input, expected) pair", name, j)
			}
			args := []any{pair[0]}
			if tup, ok := pair[0].(execSketchTuple); ok {
				args = []any(tup)
			}
			out[name] = append(out[name], execSketchRow{line: p.rowLines[j], args: args, want: pair[1]})
		}
	}
	return out
}

// execSketchParser reads the literal subset the sketch's tests use: ints,
// strings (with the escapes \" \\ \n \t and \uXXXX), booleans, lists, records
// and tuples.
type execSketchParser struct {
	src      string
	pos      int
	depth    int
	rowLines []int // line of each element of the outermost list
	t        *testing.T
}

func (p *execSketchParser) skip() {
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

func (p *execSketchParser) line() int { return strings.Count(p.src[:p.pos], "\n") + 1 }

func (p *execSketchParser) fail(format string, args ...any) {
	p.t.Fatalf("sketch parse at line %d: %s", p.line(), fmt.Sprintf(format, args...))
}

func (p *execSketchParser) expect(c byte) {
	p.skip()
	if p.pos >= len(p.src) || p.src[p.pos] != c {
		p.fail("want %q", c)
	}
	p.pos++
}

func (p *execSketchParser) seq(closer byte, item func()) {
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

func (p *execSketchParser) str() string {
	p.pos++ // opening quote
	var b strings.Builder
	for {
		if p.pos >= len(p.src) {
			p.fail("unterminated string")
		}
		c := p.src[p.pos]
		switch {
		case c == '"':
			p.pos++
			return b.String()
		case c == '\n':
			p.fail("newline inside a string literal")
		case c != '\\':
			_, size := utf8.DecodeRuneInString(p.src[p.pos:])
			b.WriteString(p.src[p.pos : p.pos+size])
			p.pos += size
			continue
		}
		if p.pos+1 >= len(p.src) {
			p.fail("dangling escape")
		}
		switch e := p.src[p.pos+1]; e {
		case '"', '\\':
			b.WriteByte(e)
			p.pos += 2
		case 'n':
			b.WriteByte('\n')
			p.pos += 2
		case 't':
			b.WriteByte('\t')
			p.pos += 2
		case 'u':
			if p.pos+6 > len(p.src) {
				p.fail("short \\u escape")
			}
			n, err := strconv.ParseUint(p.src[p.pos+2:p.pos+6], 16, 32)
			if err != nil {
				p.fail("bad \\u escape %q", p.src[p.pos:p.pos+6])
			}
			b.WriteRune(rune(n))
			p.pos += 6
		default:
			p.fail("escape \\%c is outside the parsed subset", e)
		}
	}
}

func (p *execSketchParser) value() any {
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
				p.rowLines = append(p.rowLines, p.line())
			}
			list = append(list, p.value())
		})
		p.depth--
		return list
	case c == '(':
		p.pos++
		p.depth++
		tup := execSketchTuple{}
		p.seq(')', func() { tup = append(tup, p.value()) })
		p.depth--
		if len(tup) == 1 {
			return tup[0]
		}
		return tup
	case c == '{':
		p.pos++
		p.depth++
		rec := map[string]any{}
		p.seq('}', func() {
			p.skip()
			start := p.pos
			for p.pos < len(p.src) && (p.src[p.pos] == '_' || execIsAlnum(p.src[p.pos])) {
				p.pos++
			}
			key := p.src[start:p.pos]
			p.expect(':')
			rec[key] = p.value()
		})
		p.depth--
		return rec
	case c == '"':
		return p.str()
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
	default:
		p.fail("unsupported literal starting %q", c)
	}
	return nil
}

func execIsAlnum(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9'
}

func xAt(t *testing.T, a []any, i int) any {
	t.Helper()
	if i >= len(a) {
		t.Fatalf("sketch row has %d args, mirror reads arg %d", len(a), i)
	}
	return a[i]
}

func xInt(t *testing.T, a []any, i int) int64 {
	t.Helper()
	v, ok := xAt(t, a, i).(int64)
	if !ok {
		t.Fatalf("arg %d is %T, want int", i, a[i])
	}
	return v
}

func xString(t *testing.T, a []any, i int) string {
	t.Helper()
	v, ok := xAt(t, a, i).(string)
	if !ok {
		t.Fatalf("arg %d is %T, want string", i, a[i])
	}
	return v
}

func xBool(t *testing.T, a []any, i int) bool {
	t.Helper()
	v, ok := xAt(t, a, i).(bool)
	if !ok {
		t.Fatalf("arg %d is %T, want bool", i, a[i])
	}
	return v
}

func xList(t *testing.T, a []any, i int) []any {
	t.Helper()
	v, ok := xAt(t, a, i).([]any)
	if !ok {
		t.Fatalf("arg %d is %T, want a list", i, a[i])
	}
	return v
}

func xRecord(t *testing.T, a []any, i int, keys ...string) map[string]any {
	t.Helper()
	rec, ok := xAt(t, a, i).(map[string]any)
	if !ok {
		t.Fatalf("arg %d is %T, want a record", i, a[i])
	}
	var got []string
	for k := range rec {
		got = append(got, k)
	}
	sort.Strings(got)
	want := append([]string(nil), keys...)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("record fields %v, want %v", got, want)
	}
	return rec
}

type execRule struct{ name, cls string }

func xRule(t *testing.T, a []any, i int) execRule {
	t.Helper()
	r := xRecord(t, a, i, "name", "cls")
	f := []any{r["name"], r["cls"]}
	return execRule{name: xString(t, f, 0), cls: xString(t, f, 1)}
}

// xFlags turns the sketch's rule list into the profile's name -> class map;
// a repeated name would make the two shapes disagree, so it is refused.
func xFlags(t *testing.T, list []any) map[string]string {
	t.Helper()
	flags := map[string]string{}
	for i := range list {
		r := xRule(t, list, i)
		if _, dup := flags[r.name]; dup {
			t.Fatalf("sketch rule list repeats flag %q", r.name)
		}
		flags[r.name] = r.cls
	}
	return flags
}

func xGrammar(t *testing.T, a []any, i int) ExecCommand {
	t.Helper()
	g := xRecord(t, a, i, "flags", "positional", "suffixes", "maxArgs", "passthrough")
	f := []any{g["flags"], g["positional"], g["suffixes"], g["maxArgs"], g["passthrough"]}
	cmd := ExecCommand{
		Argv:       []string{"tool"},
		Flags:      xFlags(t, xList(t, f, 0)),
		Positional: xString(t, f, 1),
		MaxArgs:    int(xInt(t, f, 3)),
	}
	suffixes := xList(t, f, 2)
	for j := range suffixes {
		cmd.Suffixes = append(cmd.Suffixes, xString(t, suffixes, j))
	}
	if xBool(t, f, 4) {
		cmd.Passthrough = ExecPassthrough
	}
	return cmd
}
