package sketchtest

import "testing"

func TestParseReadsScalarRows(t *testing.T) {
	src := `module x/y

-- a comment
export func f(a: int, b: bool) -> bool ! {}
ensures { result == b }
tests [
  ((1, true), true),   -- trailing comment
  ((-2, false), false)
]
{ b }

export func g(s: string) -> string ! {}
tests [("a\"b", "c\\d"), ("", "x")]
{ s }

export func noTests(n: int) -> int ! {} { n }
`
	laws := Parse(t, src)
	if len(laws) != 2 || len(laws["f"]) != 2 || len(laws["g"]) != 2 {
		t.Fatalf("laws = %+v", laws)
	}
	r := laws["f"][1]
	if r.Line != 8 || r.Args[0] != int64(-2) || r.Args[1] != false || r.Want != false {
		t.Fatalf("f row 2 = %+v", r)
	}
	if g := laws["g"][0]; g.Args[0] != `a"b` || g.Want != `c\d` {
		t.Fatalf("g row 1 = %+v", g)
	}
	n := Replay(t, laws, map[string]Mirror{
		"f": func(tb testing.TB, a []any) any { return Bool(tb, a, 1) },
		"g": func(tb testing.TB, a []any) any {
			if Str(tb, a, 0) == "" {
				return "x"
			}
			return `c\d`
		},
	}, 4)
	if n != 4 {
		t.Fatalf("replayed %d rows, want 4", n)
	}
}
