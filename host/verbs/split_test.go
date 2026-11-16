package verbs

// Row 243 V1 (design §6 VB1, VB2).

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

// VB1: Split's words and each named refusal, at and past each bound.
func TestSplit(t *testing.T) {
	for _, c := range []struct {
		line string
		want []string
	}{
		{"", []string{}},
		{"  episode   report  w1 ", []string{"episode", "report", "w1"}},
		{`why 'a b' "c d"`, []string{"why", "a b", "c d"}},
		{`'' x`, []string{"", "x"}},
		{`a'b'"c"d`, []string{"abcd"}},
		{`"q\"uo\\te\n"`, []string{`q"uo\te\n`}},
		{`'it'"'"'s'`, []string{"it's"}},
		{`a\b`, []string{`a\b`}},
		{`'$x' "a;b|c&d<e>f(g)h` + "`" + `"`, []string{"$x", "a;b|c&d<e>f(g)h`"}},
		{"é ü", []string{"é", "ü"}},
		{strings.Repeat("x", MaxLine), []string{strings.Repeat("x", MaxLine)}},
		{strings.TrimSpace(strings.Repeat("w ", MaxWords)), strings.Fields(strings.Repeat("w ", MaxWords))},
	} {
		got, err := Split(c.line)
		if err != nil || !reflect.DeepEqual(got, c.want) {
			t.Errorf("Split(%q) = %q, %v; want %q", c.line, got, err, c.want)
		}
	}
	for _, c := range []struct {
		line string
		want error
	}{
		{strings.Repeat("x", MaxLine+1), ErrTooLong},
		{strings.TrimSpace(strings.Repeat("w ", MaxWords+1)), ErrTooMany},
		{"a\nb", ErrNotText}, {"a\tb", ErrNotText}, {"a\x7fb", ErrNotText}, {"a\x00", ErrNotText}, {"a\x1f", ErrNotText},
		{"a \xff", ErrNotText}, {`'a`, ErrUnclosed}, {`"a`, ErrUnclosed}, {`"a\"`, ErrUnclosed}, {`a 'b' "c`, ErrUnclosed},
	} {
		if got, err := Split(c.line); err != c.want || got != nil {
			t.Errorf("Split(%.40q) = %q, %v; want %v", c.line, got, err, c.want)
		}
	}
	for _, ch := range "$`;|&<>()" { // the nine
		if _, err := Split("why a" + string(ch) + "b"); err != ErrShell {
			t.Errorf("unquoted %q: %v, want Shell", ch, err)
		}
		for _, q := range []string{"'a" + string(ch) + "b'", `"a` + string(ch) + `b"`} {
			if got, err := Split(q); err != nil || len(got) != 1 || got[0] != "a"+string(ch)+"b" {
				t.Errorf("quoted %s = %q, %v", q, got, err)
			}
		}
	}
	if !strings.HasPrefix(ErrTooLong.Error(), "TooLong") || !strings.HasPrefix(ErrShell.Error(), "Shell") {
		t.Error("a refusal's message does not lead with its name")
	}
}

// VB2: Quote is Split's inverse over generated argv, and a safe word stays bare.
func TestQuoteRoundTrip(t *testing.T) {
	alphabet := []rune("ab-_./:@=+%,'\"\\ $`;|&<>()*?~#!{}[]éx")
	rng := rand.New(rand.NewSource(243))
	check := func(argv []string) {
		t.Helper()
		line := Quote(argv)
		got, err := Split(line)
		if err != nil || !reflect.DeepEqual(got, argv) {
			t.Fatalf("Split(Quote(%q)) = Split(%s) = %q, %v", argv, line, got, err)
		}
	}
	for i := 0; i < 2000; i++ {
		argv := make([]string, 1+rng.Intn(6))
		for j := range argv {
			w := make([]rune, rng.Intn(8))
			for k := range w {
				w[k] = alphabet[rng.Intn(len(alphabet))]
			}
			argv[j] = string(w)
		}
		check(argv)
	}
	if got := Quote([]string{"ailang-worldd", "episode", "report", "w1", "--json"}); got != "ailang-worldd episode report w1 --json" {
		t.Errorf("safe words quoted: %s", got)
	}
	if got := Quote([]string{"a b", "it's", "$x", ""}); got != `'a b' 'it'"'"'s' '$x' ''` {
		t.Errorf("Quote = %s", got)
	}
}
