package broker

// Row 156 M1a: the protected grammar, canonical digest and matcher fixtures
// (exec_protected.go). Pure functions: no profile, no filesystem, no sockets.

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// wantProtectedRefusal asserts that in is refused and that the error names the
// entry index and the reason tag.
func wantProtectedRefusal(t *testing.T, in []string, index int, tag string) {
	t.Helper()
	got, err := ParseProtectedPatterns(in)
	if err == nil {
		t.Fatalf("ParseProtectedPatterns(%q) accepted as %q; want refusal [%s]", in, got, tag)
	}
	msg := err.Error()
	if !strings.Contains(msg, "protected["+strconv.Itoa(index)+"]") || !strings.Contains(msg, "["+tag+"]") {
		t.Fatalf("ParseProtectedPatterns(%q) error %q; want index %d and reason [%s]", in, msg, index, tag)
	}
}

func TestProtectedGrammar(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		for _, form := range []string{"tests/**", "*_test.go", "*.test.ts", "conftest.py", "[ab]_test.go", "a?c.go", "pkg/sub/**"} {
			got, err := ParseProtectedPatterns([]string{form})
			if err != nil {
				t.Errorf("ParseProtectedPatterns(%q): %v", form, err)
				continue
			}
			if !reflect.DeepEqual(got, []string{form}) {
				t.Errorf("ParseProtectedPatterns(%q) = %q; want the form unchanged", form, got)
			}
		}
		if got, err := ParseProtectedPatterns(nil); err != nil || len(got) != 0 {
			t.Errorf("nil input = %q, %v; want empty, nil", got, err)
		}
		if got, err := ParseProtectedPatterns([]string{}); err != nil || len(got) != 0 {
			t.Errorf("empty input = %q, %v; want empty, nil", got, err)
		}
	})
	t.Run("count-65", func(t *testing.T) {
		in := make([]string, 65)
		for i := range in {
			in[i] = "f" + strconv.Itoa(i) + ".go"
		}
		wantProtectedRefusal(t, in, 64, protTagCount)
		if _, err := ParseProtectedPatterns(in[:64]); err != nil {
			t.Errorf("64 entries refused: %v", err)
		}
	})
	t.Run("empty-and-257", func(t *testing.T) {
		wantProtectedRefusal(t, []string{"ok.go", ""}, 1, protTagLength)
		wantProtectedRefusal(t, []string{"ok.go", strings.Repeat("a", 254) + ".go"}, 1, protTagLength)
		if _, err := ParseProtectedPatterns([]string{strings.Repeat("a", 253) + ".go"}); err != nil {
			t.Errorf("256-byte entry refused: %v", err)
		}
	})
	t.Run("characters", func(t *testing.T) {
		for _, bad := range []string{"a\x00.go", "a\x01.go", `a\b.go`, "\xff.go", "x.go ", " x.go", "tests/**\t"} {
			wantProtectedRefusal(t, []string{"ok.go", bad}, 1, protTagCharacters)
		}
	})
	t.Run("relative", func(t *testing.T) {
		for _, bad := range []string{"/abs/**", "./tests/**", "a/../b/**", "a/./b/**", "a//b/**"} {
			wantProtectedRefusal(t, []string{"ok.go", bad}, 1, protTagRelative)
		}
	})
	t.Run("slash-glob", func(t *testing.T) {
		for _, bad := range []string{"src/*.ts", "tests/", "tests/*"} {
			wantProtectedRefusal(t, []string{"ok.go", bad}, 1, protTagSlash)
		}
	})
	t.Run("dir-glob", func(t *testing.T) {
		for _, bad := range []string{"te*/**", "a?/**", "[ab]/**"} {
			wantProtectedRefusal(t, []string{"ok.go", bad}, 1, protTagGlob)
		}
	})
	t.Run("bracket", func(t *testing.T) {
		wantProtectedRefusal(t, []string{"ok.go", "[a.go"}, 1, protTagGlob)
	})
	t.Run("everything", func(t *testing.T) {
		for _, bad := range []string{"*", "**", "./**"} {
			wantProtectedRefusal(t, []string{"ok.go", bad}, 1, protTagEverything)
		}
	})
	t.Run("unsupported-meta", func(t *testing.T) {
		for _, bad := range []string{"{a,b}.go", "!x.go"} {
			wantProtectedRefusal(t, []string{"ok.go", bad}, 1, protTagMeta)
		}
	})
	t.Run("already-denied", func(t *testing.T) {
		for _, bad := range []string{".git/**", ".GIT/**", ".github/workflows/**", ".claude/**", ".gitmodules/**", "sub/.git/**", "x/.GIT/y/**"} {
			wantProtectedRefusal(t, []string{"ok.go", bad}, 1, protTagDenied)
		}
	})
}

func TestProtectedGrammarRefusesDoubleStarPrefix(t *testing.T) {
	// Slash-free '**' has no directory to suffix; its own fix text is asserted.
	for _, bad := range []string{"**x.go", "a**b.go"} {
		wantProtectedRefusal(t, []string{"ok.go", bad}, 1, protTagDoubleStar)
		if _, err := ParseProtectedPatterns([]string{bad}); err == nil || !strings.Contains(err.Error(), "use '*' in place of '**'") {
			t.Errorf("%q: error %v lacks its fix text", bad, err)
		}
	}
	for _, tc := range []struct{ in, suggest string }{
		{"**/x.test.ts", "x.test.ts"},
		{"src/**/x.go", "x.go"},
	} {
		_, err := ParseProtectedPatterns([]string{tc.in})
		if err == nil {
			t.Fatalf("ParseProtectedPatterns(%q) accepted; want refusal", tc.in)
		}
		msg := err.Error()
		if !strings.Contains(msg, "use the base-name form") {
			t.Errorf("%q: error %q lacks the fix text", tc.in, msg)
		}
		if !strings.Contains(msg, "`"+tc.suggest+"`") {
			t.Errorf("%q: error %q does not name %q", tc.in, msg, tc.suggest)
		}
	}
}

func TestProtectedDigestGolden(t *testing.T) {
	// Hand-derived canonical list: fold keys are TESTS/**, *.TEST.TS, *_TEST.GO
	// (TESTS/** duplicates tests/**, *_test.go repeats); first spellings kept;
	// sorted by fold-key bytes ('*' 0x2A < 'T' 0x54, '.' 0x2E < '_' 0x5F).
	// The fold key TESTS/** has spellings "tests/**" and "TESTS/**"; the
	// byte-smallest, "TESTS/**" ('T' 0x54 < 't' 0x74), is kept.
	// Digest, computed outside Go:
	//   printf 'world/exec-protected/v1\n*.test.ts\n*_test.go\nTESTS/**\n' | shasum -a 256
	const want = "sha256:bbe9239ed936ef97f569a718024b72ca42d652ba2025a2006514badb18784965"
	wantCanon := []string{"*.test.ts", "*_test.go", "TESTS/**"}
	in := []string{"tests/**", "*_test.go", "TESTS/**", "*.test.ts", "*_test.go"}
	got, err := ParseProtectedPatterns(in)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, wantCanon) {
		t.Fatalf("canonical = %q; want %q", got, wantCanon)
	}
	if d := protectedDigest(got); d != want {
		t.Fatalf("digest = %s; want %s", d, want)
	}
	shuffled := []string{"*.test.ts", "tests/**", "TESTS/**", "*_test.go"}
	got2, err := ParseProtectedPatterns(shuffled)
	if err != nil {
		t.Fatal(err)
	}
	if d := protectedDigest(got2); d != want {
		t.Fatalf("shuffled digest = %s; want %s (order-independent)", d, want)
	}
	// Spelling order must not matter: both orders give the same list and digest.
	for _, pair := range [][]string{{"tests/**", "TESTS/**"}, {"TESTS/**", "tests/**"}} {
		c, err := ParseProtectedPatterns(pair)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(c, []string{"TESTS/**"}) || protectedDigest(c) != protectedDigest([]string{"TESTS/**"}) {
			t.Fatalf("%q canonical = %q; want [TESTS/**] with a shared digest", pair, c)
		}
	}
	if d := protectedDigest(nil); d != "" {
		t.Fatalf("empty digest = %q; want \"\" (none declared)", d)
	}
}

func TestProtectedMatchFixture(t *testing.T) {
	cases := []struct {
		name     string
		patterns []string
		path     string
		want     bool
	}{
		{"pkg-A_TEST.GO", []string{"*_test.go"}, "pkg/A_TEST.GO", true},
		{"root", []string{"*_test.go"}, "c_test.go", true},
		{"deep", []string{"*_test.go"}, "a/b/c/d_test.go", true},
		{"dir-root", []string{"tests/**"}, "tests", true},
		{"dir-deep", []string{"tests/**"}, "tests/a/b.py", true},
		{"long-s", []string{"tests/**"}, "teſts/x.py", true},
		{"kelvin", []string{"*k_test.go"}, "pkg/K_test.go", true},
		{"non-hit-sibling", []string{"tests/**"}, "tests2/x.py", false},
		{"non-hit-nested-dir", []string{"tests/**"}, "cli/tests/x.py", false},
		{"non-hit-base", []string{"*_test.go"}, "pkg/a.go", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			canon, err := ParseProtectedPatterns(tc.patterns)
			if err != nil {
				t.Fatal(err)
			}
			if got := matchProtected(canon, tc.path); got != tc.want {
				t.Fatalf("matchProtected(%q, %q) = %v; want %v", tc.patterns, tc.path, got, tc.want)
			}
		})
	}
	t.Run("clean-path", func(t *testing.T) {
		canon, err := ParseProtectedPatterns([]string{"tests/**"})
		if err != nil {
			t.Fatal(err)
		}
		for _, p := range []string{"./tests/x.py", "a/../tests/x.py"} {
			if !matchProtected(canon, p) {
				t.Errorf("tests/** did not match %q", p)
			}
		}
	})
	t.Run("nfc", func(t *testing.T) {
		nfc := "café_test.go"
		nfd := "café_test.go"
		canon, err := ParseProtectedPatterns([]string{nfc})
		if err != nil {
			t.Fatal(err)
		}
		if !matchProtected(canon, nfd) {
			t.Errorf("NFC pattern did not match NFD path %q", nfd)
		}
		canonD, err := ParseProtectedPatterns([]string{nfd})
		if err != nil {
			t.Fatal(err)
		}
		if !matchProtected(canonD, nfc) {
			t.Errorf("NFD pattern did not match NFC path %q", nfc)
		}
	})
}
