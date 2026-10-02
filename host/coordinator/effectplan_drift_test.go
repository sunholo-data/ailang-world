package coordinator

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// The drift test binds effectplan.go to design_docs/sketches/effectplan.ail:
// it READS the sketch, parses every inline `tests [...]` row of every law,
// and evaluates each row through the Go mirror of the same name. A law added
// to the sketch with no Go mirror, a mirror with no sketch law, a changed
// boundary (16 bytes, 1 MiB, one effect) on either side, or a row whose
// expected value the Go mirror disagrees with, turns this test red. This is
// the effectbroker pattern (host/broker/decide_test.go TestSketchRows) with
// the rows read from the sketch instead of transcribed by hand.

type sketchTuple []any

type sketchRow struct {
	line int
	args []any
	want any
}

var sketchMirrors = map[string]func(t *testing.T, args []any) any{
	"effectCountOk":     func(t *testing.T, a []any) any { return effectCountOk(int(argInt(t, a, 0))) },
	"idLengthOk":        func(t *testing.T, a []any) any { return idLengthOk(argString(t, a, 0)) },
	"idByteOk":          func(t *testing.T, a []any) any { return idByteOk(int(argInt(t, a, 0))) },
	"idsDistinct":       func(t *testing.T, a []any) any { return idsDistinct(argString(t, a, 0), argString(t, a, 1)) },
	"payloadSizeOk":     func(t *testing.T, a []any) any { return payloadSizeOk(int(argInt(t, a, 0))) },
	"reservedOutputKey": func(t *testing.T, a []any) any { return reservedOutputKey(argString(t, a, 0)) },
	"finishKeyAllowed":  func(t *testing.T, a []any) any { return finishKeyAllowed(argString(t, a, 0)) },
	"resultPresenceOk": func(t *testing.T, a []any) any {
		return resultPresenceOk(int(argInt(t, a, 0)), argBool(t, a, 1))
	},
	"finishNeedsEffect": func(t *testing.T, a []any) any {
		return finishNeedsEffect(int(argInt(t, a, 0)), argBool(t, a, 1))
	},
	"requirementMatches": func(t *testing.T, a []any) any {
		return requirementMatches(argRequirement(t, a, 0), argEffect(t, a, 1))
	},
	"declaredContains": func(t *testing.T, a []any) any {
		list, ok := argAt(t, a, 0).([]any)
		if !ok {
			t.Fatalf("declaredContains arg 0 is %T, want a list", argAt(t, a, 0))
		}
		ds := make([]transitionreg.EffectRequirement, len(list))
		for i := range list {
			ds[i] = argRequirement(t, list, i)
		}
		return declaredContains(ds, argEffect(t, a, 1))
	},
	"effectLawfulAgainst": func(t *testing.T, a []any) any {
		return effectLawfulAgainst(argRequirement(t, a, 0), argEffect(t, a, 1))
	},
	"planShapeLawful": func(t *testing.T, a []any) any {
		return planShapeLawful(int(argInt(t, a, 0)), argBool(t, a, 1), argBool(t, a, 2))
	},
}

func TestEffectPlanSketchDrift(t *testing.T) {
	src, err := os.ReadFile(filepath.Join("..", "..", "design_docs", "sketches", "effectplan.ail"))
	if err != nil {
		t.Fatalf("read sketch: %v", err)
	}
	if !strings.HasPrefix(string(src), "module sketches/effectplan\n") {
		t.Fatal("sketch module header moved")
	}
	laws := parseSketchTests(t, string(src))

	var sketchNames, goNames []string
	for name := range laws {
		sketchNames = append(sketchNames, name)
	}
	for name := range sketchMirrors {
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
		if len(rows) < 2 {
			t.Fatalf("sketch law %s has %d test rows; a law needs a true and a false row", name, len(rows))
		}
		sawTrue, sawFalse := false, false
		for _, row := range rows {
			total++
			got := sketchMirrors[name](t, row.args)
			if got != row.want {
				t.Errorf("sketch line %d %s%v: Go mirror = %v, sketch expects %v", row.line, name, row.args, got, row.want)
			}
			sawTrue = sawTrue || row.want == true
			sawFalse = sawFalse || row.want == false
		}
		if !sawTrue || !sawFalse {
			t.Errorf("sketch law %s lacks a true or a false row", name)
		}
	}
	// Instrument-health floor: the parser found the rows the sketch carries
	// today (a parser that silently found nothing would pass the loop above).
	if total < 50 {
		t.Fatalf("parsed %d sketch rows; expected at least 50", total)
	}
	t.Logf("evaluated %d sketch rows over %d laws", total, len(sketchNames))
}

var exportFuncRE = regexp.MustCompile(`(?m)^export func ([A-Za-z0-9_]+)\(`)

// parseSketchTests returns, per exported function that carries inline tests,
// its rows. A function's tests block is the first `tests [` between its
// header and the next exported function.
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
		seg := src[h[0]:end]
		at := strings.Index(seg, "\ntests [")
		if at < 0 {
			continue
		}
		p := &sketchParser{src: src, pos: h[0] + at + len("\ntests "), t: t}
		list, ok := p.value().([]any)
		if !ok {
			t.Fatalf("%s: tests is not a list", name)
		}
		lines := p.rowLines
		for j, item := range list {
			pair, ok := item.(sketchTuple)
			if !ok || len(pair) != 2 {
				t.Fatalf("%s row %d is not an (input, expected) pair", name, j)
			}
			args := []any{pair[0]}
			if tup, ok := pair[0].(sketchTuple); ok {
				args = []any(tup)
			}
			out[name] = append(out[name], sketchRow{line: lines[j], args: args, want: pair[1]})
		}
	}
	return out
}

// sketchParser reads the literal subset the sketch's tests use: ints,
// strings without escapes, booleans, lists, records and tuples.
type sketchParser struct {
	src      string
	pos      int
	depth    int
	rowLines []int // line of each element of the outermost list
	t        *testing.T
}

func (p *sketchParser) skip() {
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

func (p *sketchParser) line() int { return strings.Count(p.src[:p.pos], "\n") + 1 }

func (p *sketchParser) fail(format string, args ...any) {
	p.t.Fatalf("sketch parse at line %d: %s", p.line(), fmt.Sprintf(format, args...))
}

func (p *sketchParser) expect(c byte) {
	p.skip()
	if p.pos >= len(p.src) || p.src[p.pos] != c {
		p.fail("want %q", c)
	}
	p.pos++
}

func (p *sketchParser) seq(closer byte, item func()) {
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

func (p *sketchParser) value() any {
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
		tup := sketchTuple{}
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
			for p.pos < len(p.src) && (p.src[p.pos] == '_' || isAlnum(p.src[p.pos])) {
				p.pos++
			}
			key := p.src[start:p.pos]
			p.expect(':')
			rec[key] = p.value()
		})
		p.depth--
		return rec
	case c == '"':
		end := strings.IndexByte(p.src[p.pos+1:], '"')
		if end < 0 {
			p.fail("unterminated string")
		}
		s := p.src[p.pos+1 : p.pos+1+end]
		if strings.Contains(s, `\`) {
			p.fail("escapes are outside the parsed subset")
		}
		p.pos += end + 2
		return s
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

func isAlnum(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9'
}

func argAt(t *testing.T, a []any, i int) any {
	t.Helper()
	if i >= len(a) {
		t.Fatalf("sketch row has %d args, mirror reads arg %d", len(a), i)
	}
	return a[i]
}

func argInt(t *testing.T, a []any, i int) int64 {
	t.Helper()
	v, ok := argAt(t, a, i).(int64)
	if !ok {
		t.Fatalf("arg %d is %T, want int", i, a[i])
	}
	return v
}

func argString(t *testing.T, a []any, i int) string {
	t.Helper()
	v, ok := argAt(t, a, i).(string)
	if !ok {
		t.Fatalf("arg %d is %T, want string", i, a[i])
	}
	return v
}

func argBool(t *testing.T, a []any, i int) bool {
	t.Helper()
	v, ok := argAt(t, a, i).(bool)
	if !ok {
		t.Fatalf("arg %d is %T, want bool", i, a[i])
	}
	return v
}

func argRecord(t *testing.T, a []any, i int, keys ...string) map[string]any {
	t.Helper()
	rec, ok := argAt(t, a, i).(map[string]any)
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

func argRequirement(t *testing.T, a []any, i int) transitionreg.EffectRequirement {
	t.Helper()
	r := argRecord(t, a, i, "effect", "scope", "cost")
	f := []any{r["effect"], r["scope"], r["cost"]}
	return transitionreg.EffectRequirement{Effect: argString(t, f, 0), Scope: argString(t, f, 1), Cost: argInt(t, f, 2)}
}

func argEffect(t *testing.T, a []any, i int) lawEffect {
	t.Helper()
	r := argRecord(t, a, i, "id", "effect", "scope", "cost", "payloadBytes")
	f := []any{r["id"], r["effect"], r["scope"], r["cost"], r["payloadBytes"]}
	return lawEffect{
		ID: argString(t, f, 0), Effect: argString(t, f, 1), Scope: argString(t, f, 2),
		Cost: argInt(t, f, 3), PayloadBytes: int(argInt(t, f, 4)),
	}
}
