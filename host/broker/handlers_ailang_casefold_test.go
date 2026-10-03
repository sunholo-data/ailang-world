package broker

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"unicode"
)

// Row 134 V68 / R-SE-15: on macOS APFS (case-insensitive, the default)
// v0.51.0 matches fs_deny_write case-SENSITIVELY in both the FS effect and
// policy-tool; only `.git` is case-folded upstream. Under a lowercase-only
// list a running program and the write/edit ops wrote `.CLAUDE/settings.json`
// (= `.claude/settings.json`), `.GITMODULES` (= `.gitmodules`) and the rest.
// World closes it in the rendered policy alone: every fold variant of each
// protected directory as a literal `/**` prefix, one character-class pattern
// per protected file. The AC4.1 matrix (handlers_ailang_test.go) carries the
// case-variant rows through caseVariantDenyTargets.

// asciiCaseVariants lists every ASCII-case spelling of name, all-lowercase
// first, the last letter varying fastest.
func asciiCaseVariants(name string) []string {
	var letters []int
	for i, r := range name {
		if unicode.IsLetter(r) {
			letters = append(letters, i)
		}
	}
	n := len(letters)
	out := make([]string, 0, 1<<n)
	for mask := 0; mask < 1<<n; mask++ {
		b := []byte(name)
		for j, at := range letters {
			if mask&(1<<(n-1-j)) != 0 {
				b[at] = byte(unicode.ToUpper(rune(b[at])))
			}
		}
		out = append(out, string(b))
	}
	return out
}

// v051MatchDenyWrite mirrors effects.MatchDenyWrite at AILANG v0.51.0
// (internal/effects/fs_root.go:183-200), the matcher both the FS effect and
// policy-tool share. It lets the binary-independent tier sweep every variant;
// the real-binary tier below proves the binary agrees.
func v051MatchDenyWrite(patterns []string, rel string) string {
	rel = strings.TrimPrefix(filepath.ToSlash(rel), "./")
	base := rel
	if i := strings.LastIndexByte(rel, '/'); i >= 0 {
		base = rel[i+1:]
	}
	for _, pat := range patterns {
		pat = strings.TrimPrefix(filepath.ToSlash(pat), "./")
		if dir, ok := strings.CutSuffix(pat, "/**"); ok {
			if rel == dir || strings.HasPrefix(rel, dir+"/") {
				return pat
			}
			continue
		}
		if m, _ := path.Match(pat, rel); m {
			return pat
		}
		if m, _ := path.Match(pat, base); m {
			return pat
		}
	}
	return ""
}

// allFoldSpellings is every spelling of name that APFS resolves to it: ASCII
// case, plus ſ for each `s` (V68 measured no other non-ASCII fold onto these
// letters).
func allFoldSpellings(name string) []string {
	out := []string{""}
	for _, r := range name {
		alts := []rune{r}
		if unicode.IsLetter(r) {
			alts = append(alts, unicode.ToUpper(r))
		}
		if r == 's' {
			alts = append(alts, 'ſ')
		}
		var next []string
		for _, p := range out {
			for _, a := range alts {
				next = append(next, p+string(a))
			}
		}
		out = next
	}
	return out
}

// TestEpisodeDenyWriteCoversEveryFoldSpelling sweeps, through the v0.51.0
// matcher, every spelling of every protected name at the root, as a directory
// prefix and (for files) nested; and checks the list over-protects nothing a
// legitimate project needs (MUT-CASEFOLD-DIRS, MUT-CASEFOLD-FILES,
// MUT-CASEFOLD-LONGS).
func TestEpisodeDenyWriteCoversEveryFoldSpelling(t *testing.T) {
	swept := 0
	for _, dir := range episodeDenyDirs {
		for _, v := range allFoldSpellings(dir) {
			for _, rel := range []string{v, v + "/x", v + "/a/b/settings.json"} {
				if v051MatchDenyWrite(episodeDenyWrite, rel) == "" {
					t.Errorf("%s is not denied", rel)
				}
				swept++
			}
		}
	}
	for _, file := range episodeDenyFiles {
		for _, v := range allFoldSpellings(file) {
			for _, rel := range []string{v, "sub/" + v, "a/b/" + v} {
				if v051MatchDenyWrite(episodeDenyWrite, rel) == "" {
					t.Errorf("%s is not denied", rel)
				}
				swept++
			}
		}
	}
	for _, ok := range []string{"src/main.ail", "github/x", ".githubx/y", ".claudex", "x/.claude/y", ".gitmodules.bak",
		".gitignore", "README.md", ".pix/y", ".ailangrc", "ailang/x"} {
		if pat := v051MatchDenyWrite(episodeDenyWrite, ok); pat != "" {
			t.Errorf("%s is denied by %q, want writable", ok, pat)
		}
	}
	t.Logf("%d fold spellings denied; %d patterns rendered", swept, len(episodeDenyWrite))
}

// caseVariantSeeds are the lowercase protected targets the case-variant rows
// aim at, seeded so that a write which got through would change bytes (a
// target whose parent did not exist would fail for an unrelated reason).
var caseVariantSeeds = map[string]string{
	".claude/settings.json":   "seed claude\n",
	".github/workflows/x.yml": "seed github\n",
	".pi/settings.json":       "seed pi\n",
	".ailang/cache/z":         "seed ailang\n",
	".gitmodules":             "seed gitmodules\n",
	".gitattributes":          "seed gitattributes\n",
	"sub/.gitmodules":         "seed nested gitmodules\n",
}

func (f *realFixture) seedCaseVariantTargets(t *testing.T) {
	t.Helper()
	for rel, content := range caseVariantSeeds {
		mustMkdir(t, filepath.Dir(filepath.Join(f.root, rel)))
		mustWrite(t, filepath.Join(f.root, rel), content)
	}
	mustMkdir(t, filepath.Join(f.root, ".pi", "fresh"))
}

// caseVariantDenyTargets is variant → the lowercase path APFS resolves it to.
// Every protected entry appears; `.pi/fresh/new` is an absent target (it must
// stay absent).
var caseVariantDenyTargets = [][2]string{
	{".CLAUDE/settings.json", ".claude/settings.json"},
	{".Claude/settings.json", ".claude/settings.json"},
	{".GITHUB/workflows/x.yml", ".github/workflows/x.yml"},
	{".GitHub/workflows/x.yml", ".github/workflows/x.yml"},
	{".AILANG/cache/z", ".ailang/cache/z"},
	{".aiLang/cache/z", ".ailang/cache/z"},
	{".Pi/settings.json", ".pi/settings.json"},
	{".PI/fresh/new", ".pi/fresh/new"},
	{".GITMODULES", ".gitmodules"},
	{".GitModuleſ", ".gitmodules"},
	{"sub/.GITMODULES", "sub/.gitmodules"},
	{".GitAttributes", ".gitattributes"},
	{".gitattributeſ", ".gitattributes"},
}

// assertCaseVariantDenyRows is the AC4.1 case-variant rows: for each variant,
// the write op, the edit op and a running program (wr.ail) are each refused,
// and the lowercase target APFS resolves the variant to is byte-unchanged (or
// still absent). It returns the row count.
func (f *realFixture) assertCaseVariantDenyRows(t *testing.T) int {
	t.Helper()
	f.seedCaseVariantTargets(t)
	rows := 0
	for _, tc := range caseVariantDenyTargets {
		variant, canonical := tc[0], tc[1]
		want, existed := caseVariantSeeds[canonical]
		// Each probe starts from the seed, so each row measures one route.
		reset := func() {
			full := filepath.Join(f.root, canonical)
			if existed {
				mustWrite(t, full, want)
			} else if err := os.Remove(full); err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
		}
		check := func(how string, admitted any) {
			got, exists := f.snapshot(canonical)
			if exists != existed || got != want {
				t.Errorf("%s %s: admitted=%v and wrote through to %s (existed=%v exists=%v %q -> %q)",
					how, variant, admitted, canonical, existed, exists, want, got)
			}
		}
		resp := f.op(t, EffectWorkspaceWrite, map[string]any{"op": "write", "path": variant, "content": "pwned"})
		if resp["ok"] != false || !strings.Contains(fmt.Sprint(resp["refused"]), "fs_deny_write") {
			t.Errorf("op write %s = %v, want refused by fs_deny_write", variant, resp)
		}
		check("op write", resp["ok"])
		reset()
		resp = f.op(t, EffectWorkspaceWrite, map[string]any{"op": "edit", "path": variant, "old_text": "seed", "new_text": "pwned"})
		if resp["ok"] != false || !strings.Contains(fmt.Sprint(resp["refused"]), "fs_deny_write") {
			t.Errorf("op edit %s = %v, want refused by fs_deny_write", variant, resp)
		}
		check("op edit", resp["ok"])
		reset()
		arg, _ := json.Marshal(variant)
		res := f.run(t, f.h, "wr.ail", string(arg))
		if res.ExitCode != 1 || !res.Admitted || !strings.Contains(res.Stderr, "E_FS_PROTECTED") || strings.Contains(res.Stdout, "wrote") {
			t.Errorf("program write %s = %+v, want rc 1 E_FS_PROTECTED", variant, res)
		}
		check("program write", res.Stdout)
		reset()
		rows += 3
	}
	// No variant spelling was created as a distinct entry either (on a
	// case-sensitive volume the writes above would otherwise leave one).
	entries, err := os.ReadDir(f.root)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Name() != strings.ToLower(e.Name()) && strings.HasPrefix(e.Name(), ".") && !strings.EqualFold(e.Name(), ".git") {
			t.Errorf("a case-variant entry %q was created in the worktree", e.Name())
		}
	}
	return rows
}

// TestEpisodeDenyWriteFoldSweepOnTheToolBinary proves the real binary
// applies the rendered patterns as the mirror does: the write op on EVERY
// rendered directory variant and every file spelling is refused, the policy
// still loads and summarises, and its size is recorded (V68).
func TestEpisodeDenyWriteFoldSweepOnTheToolBinary(t *testing.T) {
	f := newRealFixture(t, "init")
	f.seedCaseVariantTargets(t)
	info, err := os.Stat(f.policyPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("rendered policy: %d bytes, %d fs_deny_write patterns", info.Size(), len(episodeDenyWrite))
	var targets []string
	for _, dir := range episodeDenyDirs {
		for _, v := range allFoldSpellings(dir) {
			targets = append(targets, v+"/swept")
		}
	}
	for _, file := range episodeDenyFiles {
		// A cross-section: the 2^n file spellings are swept through the
		// mirror; here each class position is exercised in both cases and ſ.
		up := strings.ToUpper(file)
		targets = append(targets, file, up, strings.ReplaceAll(file, "s", "ſ"),
			strings.ReplaceAll(up, "S", "ſ"), "sub/"+up, "a/b/"+strings.ReplaceAll(file, "s", "ſ"))
	}
	for _, target := range targets {
		resp := f.op(t, EffectWorkspaceWrite, map[string]any{"op": "write", "path": target, "content": "pwned"})
		if resp["ok"] != false || !strings.Contains(fmt.Sprint(resp["refused"]), "fs_deny_write") {
			t.Errorf("op write %s = %v, want refused by fs_deny_write", target, resp)
		}
	}
	for rel, want := range caseVariantSeeds {
		if got, _ := f.snapshot(rel); got != want {
			t.Errorf("%s changed: %q", rel, got)
		}
	}
	t.Logf("%d spellings refused by the tool binary", len(targets))
}
