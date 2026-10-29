package broker

// Row 156 M1a (design_docs/implemented/w-exec-protected-test-paths.md §4.1-§4.2):
// the operator's `protected` pattern set, as pure functions. ParseProtectedPatterns
// validates the grammar and returns the canonical list; canonicalProtected,
// protectedDigest, foldProtectedKey and matchProtected are the canonical form,
// its digest and the matcher. No I/O and no profile wiring here: the renderings
// for AILANG, srt and git (rebaseForSandbox, renderSrtProtected, gitCandidatePathspecs) are
// M1b; candidate parsing is M1c, with real-Git equivalence in M1d.

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

const (
	protectedMaxEntries = 64
	protectedMaxBytes   = 256
	protectedDomain     = "world/exec-protected/v1\n"
)

// Reason tags, one per refusal branch in ParseProtectedPatterns. Each error
// names the index of the offending entry and its tag.
const (
	protTagCount      = "entry-count"
	protTagLength     = "entry-length"
	protTagCharacters = "entry-characters"
	protTagMeta       = "entry-meta"
	protTagEverything = "entry-everything"
	protTagDoubleStar = "entry-doublestar"
	protTagSlash      = "entry-slash"
	protTagRelative   = "entry-relative"
	protTagGlob       = "entry-glob"
	protTagDenied     = "entry-denied"
)

// ParseProtectedPatterns validates a raw `protected` list (design §4.1) and
// returns its canonical form (NFC, deduplicated by fold key, sorted). Nil or
// empty input means no protected set: an empty result and a nil error.
func ParseProtectedPatterns(raw []string) ([]string, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	if len(raw) > protectedMaxEntries {
		return nil, protRefusal(protectedMaxEntries, raw[protectedMaxEntries], protTagCount,
			"%d entries exceeds the limit of %d", len(raw), protectedMaxEntries)
	}
	out := make([]string, 0, len(raw))
	for i, in := range raw {
		p, err := parseProtectedEntry(i, in)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return canonicalProtected(out), nil
}

// parseProtectedEntry checks one entry and returns its NFC spelling.
func parseProtectedEntry(i int, in string) (string, error) {
	if len(in) == 0 || len(in) > protectedMaxBytes {
		return "", protRefusal(i, in, protTagLength,
			"%d bytes; each entry is 1 to %d bytes", len(in), protectedMaxBytes)
	}
	if !utf8.ValidString(in) {
		return "", protRefusal(i, in, protTagCharacters, "not valid UTF-8")
	}
	if first, _ := utf8.DecodeRuneInString(in); unicode.IsSpace(first) {
		return "", protRefusal(i, in, protTagCharacters, "leading whitespace; fileguard trims patterns, so refuse rather than disagree")
	}
	if last, _ := utf8.DecodeLastRuneInString(in); unicode.IsSpace(last) {
		return "", protRefusal(i, in, protTagCharacters, "trailing whitespace; fileguard trims patterns, so refuse rather than disagree")
	}
	for _, r := range in {
		if r == 0 || unicode.IsControl(r) {
			return "", protRefusal(i, in, protTagCharacters, "contains a control character")
		}
		if r == '\\' {
			return "", protRefusal(i, in, protTagCharacters, "contains a backslash; use '/'")
		}
	}
	p := norm.NFC.String(in)
	if strings.ContainsAny(p, "{}") || strings.HasPrefix(p, "!") {
		return "", protRefusal(i, p, protTagMeta,
			"braces and a leading '!' are not supported; use one pattern per name")
	}
	switch p {
	case "*", "**", "./**", "./":
		return "", protRefusal(i, p, protTagEverything,
			"covers the whole tree; a protected set must name specific paths")
	}
	if strings.Contains(p, "**") && !(strings.HasSuffix(p, "/**") && !strings.Contains(p[:len(p)-3], "**")) {
		suggest := path.Base(strings.TrimRight(p, "/"))
		if !strings.Contains(p, "/") {
			return "", protRefusal(i, p, protTagDoubleStar,
				"'**' is not a base-name pattern; use '*' in place of '**'")
		}
		if strings.Contains(suggest, "**") {
			suggest = "*"
		}
		return "", protRefusal(i, p, protTagDoubleStar,
			"'**' is only valid as the final '/**' of a directory form; use the base-name form `%s` (a protected set has exactly two forms, so every matcher reads it the same way)",
			suggest)
	}
	if strings.Contains(p, "/") {
		if !strings.HasSuffix(p, "/**") {
			return "", protRefusal(i, p, protTagSlash,
				"a pattern with '/' must be the directory form 'D/**'")
		}
		if err := checkProtectedDir(i, p, strings.TrimSuffix(p, "/**")); err != nil {
			return "", err
		}
		return p, nil
	}
	if _, err := path.Match(p, ""); err != nil {
		return "", protRefusal(i, p, protTagGlob, "not a valid path.Match pattern")
	}
	return p, nil
}

// checkProtectedDir refuses a directory form D whose D is not clean and
// relative, uses a glob metacharacter, or names or lies inside a deny name.
func checkProtectedDir(i int, p, d string) error {
	if d == "" || strings.HasPrefix(d, "/") {
		return protRefusal(i, p, protTagRelative, "must be relative to the worktree")
	}
	for _, seg := range strings.Split(d, "/") {
		if seg == "" || seg == "." || seg == ".." {
			return protRefusal(i, p, protTagRelative, "must be clean: no '.', '..' or empty segment")
		}
	}
	if path.Clean(d) != d {
		return protRefusal(i, p, protTagRelative, "must be clean (path.Clean(D) == D)")
	}
	if strings.ContainsAny(d, "*?[") {
		return protRefusal(i, p, protTagGlob, "the directory part takes no glob metacharacter")
	}
	if protectedDenied(d) {
		return protRefusal(i, p, protTagDenied,
			"names or lies inside a path World already denies (.git or a deny name)")
	}
	return nil
}

// protectedDenied reports whether dir equals, or lies under, one of World's
// existing deny names (execDenyNames), compared by fold key.
func protectedDenied(dir string) bool {
	// fileguard's UnderGitDir denies a .git segment at any depth.
	for _, seg := range strings.Split(dir, "/") {
		if foldProtectedKey(seg) == foldProtectedKey(".git") {
			return true
		}
	}
	k := foldProtectedKey(dir)
	for _, n := range execDenyNames() {
		nk := foldProtectedKey(n)
		if k == nk || strings.HasPrefix(k, nk+"/") {
			return true
		}
	}
	return false
}

// protRefusal names the entry index, the pattern and a reason tag.
func protRefusal(i int, p, tag, format string, args ...any) error {
	return fmt.Errorf("protected[%d] %q: %s [%s]", i, p, fmt.Sprintf(format, args...), tag)
}

// canonicalProtected NFC-normalizes each pattern, collapses spellings that share
// a fold key to the byte-smallest NFC spelling (so the result does not depend on
// input order: "TESTS/**" beats "tests/**") and sorts by fold-key bytes.
func canonicalProtected(patterns []string) []string {
	idx := make(map[string]int, len(patterns))
	out := make([]string, 0, len(patterns))
	for _, p := range patterns {
		p = norm.NFC.String(p)
		k := foldProtectedKey(p)
		if i, ok := idx[k]; ok {
			if p < out[i] {
				out[i] = p
			}
			continue
		}
		idx[k] = len(out)
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool {
		return foldProtectedKey(out[i]) < foldProtectedKey(out[j])
	})
	return out
}

// protectedDigest is "sha256:" + hex of sha256 over the domain line and the
// canonical list, one pattern per line, each newline-terminated. The empty set
// has no digest: it returns "" and callers print `none declared` instead.
func protectedDigest(canonical []string) string {
	if len(canonical) == 0 {
		return ""
	}
	sum := sha256.Sum256([]byte(protectedDomain + strings.Join(canonical, "\n") + "\n"))
	return "sha256:" + hex.EncodeToString(sum[:])
}

// foldProtectedKey is fileguard's fold (AILANG v0.52.1 internal/fileguard):
// NFC, then each rune mapped to the smallest member of its SimpleFold orbit.
// It is not strings.ToLower: 'ſ' (U+017F) folds with 's' and KELVIN SIGN
// (U+212A) with 'k'.
func foldProtectedKey(s string) string {
	s = norm.NFC.String(s)
	return strings.Map(foldProtectedRune, s)
}

func foldProtectedRune(r rune) rune {
	m := r
	for f := unicode.SimpleFold(r); f != r; f = unicode.SimpleFold(f) {
		if f < m {
			m = f
		}
	}
	return m
}

// matchProtected reports whether the worktree-relative path rel is protected
// by the canonical set. A directory form D/** hits D and everything under it;
// a base-name form hits any base name it matches, at any depth.
// MatchProtected tests a path against a canonical exec-profile protected set.
func MatchProtected(canonical []string, rel string) bool {
	return matchProtected(canonical, rel)
}

func matchProtected(canonical []string, rel string) bool {
	key := foldProtectedKey(strings.TrimPrefix(path.Clean(filepath.ToSlash(rel)), "./"))
	base := path.Base(key)
	for _, p := range canonical {
		if strings.HasSuffix(p, "/**") {
			d := foldProtectedKey(strings.TrimSuffix(p, "/**"))
			if key == d || strings.HasPrefix(key, d+"/") {
				return true
			}
			continue
		}
		if ok, _ := path.Match(foldProtectedKey(p), base); ok {
			return true
		}
	}
	return false
}

// RebaseProtectedForSandbox rebases patterns for an AILANG sandbox. The input
// must be ParseProtectedPatterns output.
func RebaseProtectedForSandbox(canonical []string, moduleRoot string) []string {
	return rebaseForSandbox(canonical, moduleRoot)
}

// rebaseForSandbox renders the canonical set as AILANG's fs_deny_write, relative
// to the sandbox at moduleRoot (design §4.2 "Re-basing for AILANG"). Segments
// compare under foldProtectedKey. A D/** strictly under the root loses the root
// segments; D equal to the root, or an ancestor of it, protects the whole
// sandbox and renders as '*'; D outside the root is dropped. A base-name G is
// unchanged. Order is kept and exact duplicates are dropped. The result never
// aliases the input and is nil when nothing is left.
func rebaseForSandbox(canonical []string, moduleRoot string) []string {
	var mSeg []string
	if m := path.Clean(moduleRoot); m != "." {
		mSeg = strings.Split(m, "/")
	}
	var out []string
	seen := make(map[string]bool, len(canonical))
	for _, p := range canonical {
		r, keep := rebaseProtected(p, mSeg)
		if keep && !seen[r] {
			seen[r] = true
			out = append(out, r)
		}
	}
	return out
}

// rebaseProtected re-bases one canonical entry onto the root segments mSeg.
func rebaseProtected(p string, mSeg []string) (string, bool) {
	if !strings.HasSuffix(p, "/**") {
		return p, true
	}
	dSeg := strings.Split(strings.TrimSuffix(p, "/**"), "/")
	n := len(mSeg)
	if len(dSeg) <= n {
		if foldSegmentsEqual(dSeg, mSeg[:len(dSeg)]) {
			return "*", true
		}
		return "", false
	}
	if !foldSegmentsEqual(mSeg, dSeg[:n]) {
		return "", false
	}
	return strings.Join(dSeg[n:], "/") + "/**", true
}

// foldSegmentsEqual compares two equal-length segment lists by fold key.
func foldSegmentsEqual(a, b []string) bool {
	for i := range a {
		if foldProtectedKey(a[i]) != foldProtectedKey(b[i]) {
			return false
		}
	}
	return true
}

// renderSrtProtected renders the canonical set as srt denyWrite entries for the
// worktree (design §4.2). Darwin gets the literal directory and '<wt>/**/G' for
// base names. Every other goos gets directory literals only, because srt skips
// globs on Linux (D-156-3 = A, V8). Order is canonical and the result is fresh.
func renderSrtProtected(canonical []string, worktree, goos string) []string {
	var out []string
	for _, p := range canonical {
		if strings.HasSuffix(p, "/**") {
			out = append(out, filepath.Join(worktree, strings.TrimSuffix(p, "/**")))
			continue
		}
		if goos == "darwin" {
			out = append(out, worktree+"/**/"+p)
		}
	}
	return out
}

// protectedWidenRunes is closed by TestProtectedFoldToASCIIClosed: these ASCII
// literals have non-ASCII equivalents under NFC + SimpleFold (design §4.2).
const protectedWidenRunes = "sSkK;`"
const protectedCheckMaxOutput = 8 << 20
const protectedIgnoredCommand = "ls-files -z -o -i --exclude-standard --"

// gitCandidatePathspecs over-selects; matchProtected is the authority. Preserve
// top/glob/icase and the bare-directory arm for file and symlink replacements.
func gitCandidatePathspecs(canonical []string) []string {
	var out []string
	for _, p := range canonical {
		directory := strings.HasSuffix(p, "/**")
		if directory {
			p = strings.TrimSuffix(p, "/**")
		}
		w := protectedStars.ReplaceAllString(protectedWiden.ReplaceAllString(p, "*"), "*")
		if directory {
			out = append(out, ":(top,glob,icase)"+w, ":(top,glob,icase)"+w+"/**")
		} else {
			out = append(out, ":(top,glob,icase)**/"+w)
		}
	}
	return out
}

// protectedChange retains raw spellings and distinct statuses (M6 reuses it).
type protectedChange struct {
	Path   string `json:"path"`
	Status string `json:"status"`
}

const protectedXYSub = `[.MADRCUT]{2} (N\.\.\.|S[.C][.M][.U]) `
const protectedMode = `[0-7]{6} `
const protectedHash = `([0-9a-f]{40}|[0-9a-f]{64})`

var protectedWiden = regexp.MustCompile(`\[[^]]*\]|[` + protectedWidenRunes + `?]|[^\x00-\x7f]`)
var protectedStars = regexp.MustCompile(`\*+`)
var protectedRecordFields = map[byte]*regexp.Regexp{
	'1': regexp.MustCompile(`^` + protectedXYSub + strings.Repeat(protectedMode, 3) + protectedHash + ` ` + protectedHash + `$`),
	'2': regexp.MustCompile(`^` + protectedXYSub + strings.Repeat(protectedMode, 3) + protectedHash + ` ` + protectedHash + ` [RC]([0-9]|[1-9][0-9]|100)$`),
	'u': regexp.MustCompile(`^` + protectedXYSub + strings.Repeat(protectedMode, 4) + protectedHash + ` ` + protectedHash + ` ` + protectedHash + `$`),
}
var protectedBranch = regexp.MustCompile(`^# branch\.(oid (` + protectedHash + `|\(initial\))|head [\s\S]+|upstream [\s\S]+|ab \+[0-9]+ -[0-9]+)$`)

// parseProtectedCandidates validates all streams; only ignored entries consult the snapshot.
func parseProtectedCandidates(canonical []string, statusOut, ignoredOut, flagsOut []byte, ignoredBaseline map[string]struct{}) (head string, changes []protectedChange, err error) {
	fail := func() (string, []protectedChange, error) {
		return "", nil, fmt.Errorf("malformed protected Git candidates")
	}
	streams := make([][]string, 3)
	for i, b := range [][]byte{statusOut, ignoredOut, flagsOut} {
		if len(b) > protectedCheckMaxOutput {
			return fail()
		}
		if len(b) == 0 {
			continue
		}
		if b[len(b)-1] != 0 {
			return fail()
		}
		streams[i] = strings.Split(string(b[:len(b)-1]), "\x00")
		for _, record := range streams[i] {
			if record == "" {
				return fail()
			}
		}
	}
	seen := map[protectedChange]bool{}
	add := func(p, status string) {
		c := protectedChange{p, status}
		if matchProtected(canonical, p) && !seen[c] {
			changes = append(changes, c)
			seen[c] = true
		}
	}
	for i := 0; i < len(streams[0]); i++ {
		r := streams[0][i]
		if len(r) < 3 || r[1] != ' ' {
			return fail()
		}
		switch r[0] {
		case '#':
			if !protectedBranch.MatchString(r) {
				return fail()
			}
			if strings.HasPrefix(r, "# branch.oid ") {
				if head != "" {
					return fail()
				}
				head = strings.TrimPrefix(r, "# branch.oid ")
			}
		case '?':
			add(r[2:], "untracked")
		case '1', '2', 'u':
			n := map[byte]int{'1': 9, '2': 10, 'u': 11}[r[0]]
			f := strings.SplitN(r, " ", n)
			if len(f) != n || f[n-1] == "" || !protectedRecordFields[r[0]].MatchString(strings.Join(f[1:n-1], " ")) {
				return fail()
			}
			if r[0] == '2' {
				if i+1 >= len(streams[0]) {
					return fail()
				}
				i++
				if f[8][0] == 'R' {
					add(streams[0][i], "deleted")
				}
			}
			status := "modified"
			if strings.Contains(f[1], "D") {
				status = "deleted"
			}
			add(f[n-1], status)
		default:
			return fail()
		}
	}
	for _, p := range streams[1] {
		if _, standing := ignoredBaseline[foldProtectedKey(p)]; !standing {
			add(p, "ignored")
		}
	}
	for _, r := range streams[2] {
		if len(r) < 3 || r[1] != ' ' || !strings.ContainsRune("HSMRCK?hsmrck", rune(r[0])) {
			return fail()
		}
		if r[0] == 'S' || (r[0] >= 'a' && r[0] <= 'z') {
			add(r[2:], "index-flag")
		}
	}
	return head, changes, nil
}
