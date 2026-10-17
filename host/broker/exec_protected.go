package broker

// Row 156 M1a (design_docs/planned/w-exec-protected-test-paths.md §4.1-§4.2):
// the operator's `protected` pattern set, as pure functions. ParseProtectedPatterns
// validates the grammar and returns the canonical list; canonicalProtected,
// protectedDigest, foldProtectedKey and matchProtected are the canonical form,
// its digest and the matcher. No I/O and no profile wiring here: the renderings
// for AILANG, srt and git are M1b and later.

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path"
	"path/filepath"
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
			"'**' is only valid as the final '/**' of a directory form; use the base-name form `%s` (World's write tools read '**/' as exactly one directory)",
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

// canonicalProtected NFC-normalizes each pattern, drops duplicates by fold key
// (keeping the first spelling seen) and sorts by fold-key bytes.
func canonicalProtected(patterns []string) []string {
	seen := make(map[string]bool, len(patterns))
	out := make([]string, 0, len(patterns))
	for _, p := range patterns {
		p = norm.NFC.String(p)
		k := foldProtectedKey(p)
		if seen[k] {
			continue
		}
		seen[k] = true
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
