package broker

// Row 156 M1a: the protected grammar, canonical digest and matcher fixtures
// (exec_protected.go). Pure functions: no profile, no filesystem, no sockets.

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
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
		for _, bad := range []string{"a\x00.go", "a\x01.go", `a\b.go`, "\xff.go", "x.go ", " x.go", "tests/**\t", "\u00a0x.go", "x.go\u3000", "x.go\u2003"} {
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
	// (TESTS/** duplicates tests/**, *_test.go repeats); the byte-smallest spelling per fold key is kept;
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

// protectedFixtureRow is one V7 fixture (q4_policytool_denywrite.out.txt): a
// pattern set, a path and the expected match. The unit test and the pinned
// real-binary test share the table.
type protectedFixtureRow struct {
	name     string
	patterns []string
	path     string
	want     bool
}

var protectedFixtureRows = []protectedFixtureRow{
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

func TestProtectedMatchFixture(t *testing.T) {
	for _, tc := range protectedFixtureRows {
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

func TestProtectedRebase(t *testing.T) {
	t.Run("root-dot", func(t *testing.T) {
		in := []string{"tests/**", "*.test.ts"}
		for _, m := range []string{".", ""} {
			if got := rebaseForSandbox(in, m); !reflect.DeepEqual(got, in) {
				t.Errorf("module %q: got %q; want %q", m, got, in)
			}
		}
	})
	cases := []struct {
		name   string
		canon  []string
		module string
		want   []string
	}{
		{"under-sub", []string{"sub/tests/**"}, "sub", []string{"tests/**"}},
		{"ancestor", []string{"sub/**"}, "sub/pkg", []string{"*"}},
		{"equal", []string{"sub/**"}, "sub", []string{"*"}},
		{"disjoint", []string{"other/**"}, "sub", nil},
		{"sibling-prefix", []string{"subway/**"}, "sub", nil},
		{"case", []string{"sub/tests/**"}, "Sub", []string{"tests/**"}},
		{"basename", []string{"*.test.ts"}, "sub", []string{"*.test.ts"}},
		{"dedupe-star", []string{"sub/**", "sub/pkg/**"}, "sub/pkg", []string{"*"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := rebaseForSandbox(tc.canon, tc.module); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("rebaseForSandbox(%q, %q) = %q; want %q", tc.canon, tc.module, got, tc.want)
			}
		})
	}
	t.Run("fresh-slice", func(t *testing.T) {
		in := []string{"*.test.ts", "tests/**"}
		got := rebaseForSandbox(in, ".")
		if !reflect.DeepEqual(got, in) {
			t.Fatalf("got %q; want %q", got, in)
		}
		got[0], got[1] = "MUT", "MUT"
		if in[0] != "*.test.ts" || in[1] != "tests/**" {
			t.Fatalf("mutating the result changed the input: %q", in)
		}
	})
	t.Run("nil-empty", func(t *testing.T) {
		if got := rebaseForSandbox(nil, "sub"); got != nil {
			t.Errorf("nil input = %q; want nil", got)
		}
		if got := rebaseForSandbox([]string{}, "."); got != nil {
			t.Errorf("empty input = %q; want nil", got)
		}
		if got := rebaseForSandbox([]string{"other/**"}, "sub"); got != nil {
			t.Errorf("all-dropped input = %q; want nil", got)
		}
	})
}

func TestProtectedRenderers(t *testing.T) {
	canon, err := ParseProtectedPatterns([]string{"tests/**", "*.test.ts", "src/gen/**", "conftest.py"})
	if err != nil {
		t.Fatal(err)
	}
	t.Run("ailang", func(t *testing.T) {
		want := []string{"*.test.ts", "conftest.py", "src/gen/**", "tests/**"}
		if got := rebaseForSandbox(canon, "."); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q; want %q", got, want)
		}
	})
	t.Run("darwin", func(t *testing.T) {
		want := []string{"/w/ep1/**/*.test.ts", "/w/ep1/**/conftest.py", "/w/ep1/src/gen", "/w/ep1/tests"}
		if got := renderSrtProtected(canon, "/w/ep1", "darwin"); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q; want %q", got, want)
		}
	})
	t.Run("linux", func(t *testing.T) {
		want := []string{"/w/ep1/src/gen", "/w/ep1/tests"}
		got := renderSrtProtected(canon, "/w/ep1", "linux")
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q; want %q", got, want)
		}
		for _, e := range got {
			if strings.Contains(e, "**") || strings.HasSuffix(e, ".test.ts") || strings.HasSuffix(e, "conftest.py") {
				t.Errorf("linux rendered a glob or base name: %q", e)
			}
		}
	})
	t.Run("git", func(t *testing.T) {
		want := []string{
			":(top,glob,icase)**/*.test.ts",
			":(top,glob,icase)**/conftest.py",
			":(top,glob,icase)src/gen",
			":(top,glob,icase)src/gen/**",
			":(top,glob,icase)tests",
			":(top,glob,icase)tests/**",
		}
		got := gitPathspecs(canon)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %q; want %q", got, want)
		}
		for _, e := range got {
			if !strings.HasPrefix(e, ":(top,glob,icase)") {
				t.Errorf("pathspec %q lacks top/glob/icase magic", e)
			}
		}
	})
	t.Run("nil-empty", func(t *testing.T) {
		if rebaseForSandbox(nil, ".") != nil || renderSrtProtected(nil, "/w", "darwin") != nil ||
			renderSrtProtected([]string{}, "/w", "linux") != nil || gitPathspecs(nil) != nil {
			t.Error("empty input rendered a non-nil result")
		}
	})
}

// policyToolReply is the part of the policy-tool reply the pinned test reads.
type policyToolReply struct {
	Ok      bool   `json:"ok"`
	Refused string `json:"refused"`
}

// policyWriteOnce runs one write op through the real policy-tool with the given
// fs_deny_write, from cwd sandbox, under the V7 environment (design §4.3).
func policyWriteOnce(ctx context.Context, bin, home, root, sandbox string, deny []string, rel string) (policyToolReply, error) {
	var r policyToolReply
	if deny == nil {
		deny = []string{}
	}
	denyJSON, err := json.Marshal(deny)
	if err != nil {
		return r, err
	}
	sandboxJSON, err := json.Marshal(sandbox)
	if err != nil {
		return r, err
	}
	policy := fmt.Sprintf("security_mode = \"restricted\"\nallowed_caps = [\"FS\", \"IO\"]\nfs_sandbox = %s\ntimeout_ms = 8000\nfs_deny_write = %s\nentry = \"main\"\n\n[budgets]\nFS = 0\n", sandboxJSON, denyJSON)
	policyPath := filepath.Join(root, "policy.toml")
	if err := os.WriteFile(policyPath, []byte(policy), 0o600); err != nil {
		return r, err
	}
	req, err := json.Marshal(map[string]string{"op": "write", "path": rel, "content": "X\n"})
	if err != nil {
		return r, err
	}
	cmd := exec.CommandContext(ctx, bin, "policy-tool", "--policy", policyPath)
	cmd.Dir = sandbox
	cmd.Env = []string{"HOME=" + home, "PATH=/usr/bin:/bin", "LANG=C", "LC_ALL=C"}
	cmd.Stdin = strings.NewReader(string(req))
	out, err := cmd.Output()
	if err != nil {
		return r, fmt.Errorf("policy-tool: %w (stdout %q)", err, out)
	}
	if err := json.Unmarshal(out, &r); err != nil {
		return r, fmt.Errorf("policy-tool reply %q: %w", out, err)
	}
	return r, nil
}

// protectedWriteRefused classifies a reply: true only for a refusal by
// fs_deny_write. Any other refusal is an error, never a protected hit.
func protectedWriteRefused(r policyToolReply) (bool, error) {
	if r.Ok {
		return false, nil
	}
	if strings.Contains(r.Refused, "matches fs_deny_write") {
		return true, nil
	}
	return false, fmt.Errorf("refused for another reason: %q", r.Refused)
}

// pinnedSandbox returns a fresh symlink-resolved temp root and the sandbox
// <root>/<module>, with every parent of rel created so that a missing parent
// cannot pass for a protection hit.
func pinnedSandbox(t *testing.T, module, rel string) (string, string) {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	sandbox := filepath.Join(root, module)
	if err := os.MkdirAll(filepath.Dir(filepath.Join(sandbox, rel)), 0o755); err != nil {
		t.Fatal(err)
	}
	return root, sandbox
}

func TestProtectedMatchFixturePinned(t *testing.T) {
	bin := toolBinary(t)
	for _, row := range protectedFixtureRows {
		t.Run(row.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			canon, err := ParseProtectedPatterns(row.patterns)
			if err != nil {
				t.Fatal(err)
			}
			root, sandbox := pinnedSandbox(t, "", row.path)
			r, err := policyWriteOnce(ctx, bin, t.TempDir(), root, sandbox, rebaseForSandbox(canon, "."), row.path)
			if err != nil {
				t.Fatal(err)
			}
			refused, err := protectedWriteRefused(r)
			if err != nil {
				t.Fatal(err)
			}
			if want := matchProtected(canon, row.path); refused != want {
				t.Fatalf("policy-tool refused=%v for %q under %q; matchProtected=%v", refused, row.path, row.patterns, want)
			}
		})
	}
	for _, row := range protectedFixtureRows {
		t.Run("empty-deny-"+row.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			root, sandbox := pinnedSandbox(t, "", row.path)
			r, err := policyWriteOnce(ctx, bin, t.TempDir(), root, sandbox, nil, row.path)
			if err != nil {
				t.Fatal(err)
			}
			if !r.Ok {
				t.Fatalf("empty deny list refused %q: %q", row.path, r.Refused)
			}
		})
	}
	t.Run("module-root-dir", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		canon, err := ParseProtectedPatterns([]string{"sub/tests/**"})
		if err != nil {
			t.Fatal(err)
		}
		deny := rebaseForSandbox(canon, "sub")
		if !reflect.DeepEqual(deny, []string{"tests/**"}) {
			t.Fatalf("rebased = %q; want [tests/**]", deny)
		}
		root, sandbox := pinnedSandbox(t, "sub", "tests/x.py")
		r, err := policyWriteOnce(ctx, bin, t.TempDir(), root, sandbox, deny, "tests/x.py")
		if err != nil {
			t.Fatal(err)
		}
		if refused, err := protectedWriteRefused(r); err != nil || !refused {
			t.Fatalf("tests/x.py: refused=%v err=%v reply=%+v; want refused by fs_deny_write", refused, err, r)
		}
		r, err = policyWriteOnce(ctx, bin, t.TempDir(), root, sandbox, deny, "x.py")
		if err != nil {
			t.Fatal(err)
		}
		if !r.Ok {
			t.Fatalf("x.py: reply %+v; want ok (outside the protected dir)", r)
		}
	})
	t.Run("module-root-star", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		canon, err := ParseProtectedPatterns([]string{"sub/**"})
		if err != nil {
			t.Fatal(err)
		}
		deny := rebaseForSandbox(canon, "sub")
		if !reflect.DeepEqual(deny, []string{"*"}) {
			t.Fatalf("rebased = %q; want [*]", deny)
		}
		root, sandbox := pinnedSandbox(t, "sub", "deep/y.py")
		r, err := policyWriteOnce(ctx, bin, t.TempDir(), root, sandbox, deny, "deep/y.py")
		if err != nil {
			t.Fatal(err)
		}
		if refused, err := protectedWriteRefused(r); err != nil || !refused {
			t.Fatalf("deep/y.py: refused=%v err=%v reply=%+v; want refused by fs_deny_write", refused, err, r)
		}
	})
}
