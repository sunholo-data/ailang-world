// Package sketchtest reads the inline `tests [...]` rows of an AILANG sketch so a Go drift
// test can replay each row through the Go mirror of the same law (coding-standards S3,
// D-WORLD-127: a law's drift test PARSES its .ail source, it does not restate rows).
//
// It parses the literal subset scalar laws use: ints, strings (with \" and \\ escapes),
// booleans and tuples. Sketches whose rows hold records or lists keep their own richer
// parser (host/broker/execargs_drift_test.go). It is a test helper: callers pass their
// testing.TB, and every parse failure is fatal.
package sketchtest

import (
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// Row is one (input, expected) pair: the input's arguments and the expected value.
type Row struct {
	Line int
	Args []any
	Want any
}

// Mirror evaluates one row's arguments through a law's Go mirror.
type Mirror func(tb testing.TB, args []any) any

var exportFuncRE = regexp.MustCompile(`(?m)^export func ([A-Za-z0-9_]+)\(`)

// Load reads the sketch at path, checks its module header, and parses every exported
// law's inline tests.
func Load(tb testing.TB, path, module string) map[string][]Row {
	tb.Helper()
	src, err := os.ReadFile(path)
	if err != nil {
		tb.Fatalf("read sketch %s: %v", path, err)
	}
	if !strings.HasPrefix(string(src), "module "+module+"\n") {
		tb.Fatalf("sketch %s does not start with module %s", path, module)
	}
	return Parse(tb, string(src))
}

// Parse returns each exported function's rows, keyed by the function's name.
func Parse(tb testing.TB, src string) map[string][]Row {
	tb.Helper()
	out := map[string][]Row{}
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
		p := &parser{src: src, pos: h[0] + at + len("\ntests "), tb: tb}
		list, ok := p.value().([]any)
		if !ok {
			tb.Fatalf("%s: tests is not a list", name)
		}
		for j, item := range list {
			pair, ok := item.(tuple)
			if !ok || len(pair) != 2 {
				tb.Fatalf("%s row %d is not an (input, expected) pair", name, j)
			}
			args := []any{pair[0]}
			if tup, ok := pair[0].(tuple); ok {
				args = []any(tup)
			}
			out[name] = append(out[name], Row{Line: p.rowLines[j], Args: args, Want: pair[1]})
		}
	}
	return out
}

// Replay requires the sketch's laws-with-tests and the mirrors to be the same set, runs
// every row through its mirror, requires each law to have at least two distinct expected
// values (so a constant mirror cannot pass), and requires at least minRows rows in all
// (an instrument-health floor against a parser that silently finds nothing). It returns
// the number of rows evaluated.
func Replay(tb testing.TB, laws map[string][]Row, mirrors map[string]Mirror, minRows int) int {
	tb.Helper()
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
		tb.Fatalf("sketch laws with tests = %v\nGo mirrors            = %v\nevery sketch law needs a Go mirror and vice versa", sketchNames, goNames)
	}
	total := 0
	for _, name := range sketchNames {
		distinct := map[any]bool{}
		for _, row := range laws[name] {
			total++
			if got := mirrors[name](tb, row.Args); got != row.Want {
				tb.Errorf("sketch line %d %s%v: Go mirror = %#v, sketch expects %#v", row.Line, name, row.Args, got, row.Want)
			}
			distinct[row.Want] = true
		}
		if len(distinct) < 2 {
			tb.Errorf("sketch law %s has %d distinct expected value(s); it needs at least 2", name, len(distinct))
		}
	}
	if total < minRows {
		tb.Fatalf("parsed %d sketch rows; expected at least %d", total, minRows)
	}
	return total
}

// Int, Str and Bool read argument i of a row with its expected type.
func Int(tb testing.TB, a []any, i int) int64 {
	tb.Helper()
	v, ok := at(tb, a, i).(int64)
	if !ok {
		tb.Fatalf("arg %d is %T, want int", i, a[i])
	}
	return v
}

func Str(tb testing.TB, a []any, i int) string {
	tb.Helper()
	v, ok := at(tb, a, i).(string)
	if !ok {
		tb.Fatalf("arg %d is %T, want string", i, a[i])
	}
	return v
}

func Bool(tb testing.TB, a []any, i int) bool {
	tb.Helper()
	v, ok := at(tb, a, i).(bool)
	if !ok {
		tb.Fatalf("arg %d is %T, want bool", i, a[i])
	}
	return v
}

func at(tb testing.TB, a []any, i int) any {
	tb.Helper()
	if i >= len(a) {
		tb.Fatalf("sketch row has %d args, mirror reads arg %d", len(a), i)
	}
	return a[i]
}

type tuple []any

type parser struct {
	src      string
	pos      int
	depth    int
	rowLines []int
	tb       testing.TB
}

func (p *parser) fail(format string, args ...any) {
	p.tb.Helper()
	p.tb.Fatalf("sketch parse at line %d: "+format, append([]any{strings.Count(p.src[:p.pos], "\n") + 1}, args...)...)
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
