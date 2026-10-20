// Package plugindoc finds the command lines the ailang-world Claude Code
// plugin prints (plugins/ailang-world/, row 189 M5) so each CLI's
// TestPluginCommandsMatchTheCLIHelp can check them against its own help
// (w-world-onboarding.md AC6: CI fails if a command the skill names is
// missing from the CLI's help).
//
// A command line is a line in a fenced block, or an inline code span, that
// starts with "<prog> " once a leading "$ " and VAR=value assignments are
// stripped. Fenced lines ending in a backslash are joined with the next, and a
// trailing " # comment" is dropped.
package plugindoc

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Line is one command line: the file it came from, its 1-based line number
// and the command text, starting with the program name.
type Line struct {
	File string
	Line int
	Text string
}

var (
	inlineCode = regexp.MustCompile("`([^`]+)`")
	envAssign  = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*=\S*\s+`)
)

// Walk returns the command lines for prog in every .md file under root, in
// path order.
func Walk(root, prog string) ([]Line, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".md") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	var out []Line
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		out = append(out, Lines(f, string(raw), prog)...)
	}
	return out, nil
}

// Lines returns the command lines for prog in one markdown text.
func Lines(file, md, prog string) []Line {
	var out []Line
	add := func(n int, text string) {
		if cmd, ok := commandText(text, prog); ok {
			out = append(out, Line{File: file, Line: n, Text: cmd})
		}
	}
	inFence := false
	joined, start := "", 0
	for i, l := range strings.Split(md, "\n") {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			if joined != "" {
				add(start, joined)
				joined = ""
			}
			inFence = !inFence
			continue
		}
		if !inFence {
			for _, m := range inlineCode.FindAllStringSubmatch(l, -1) {
				add(i+1, m[1])
			}
			continue
		}
		if joined == "" {
			start = i + 1
		}
		if cont, ok := strings.CutSuffix(trimmed, "\\"); ok {
			joined += cont + " "
			continue
		}
		add(start, joined+trimmed)
		joined = ""
	}
	return out
}

// commandText strips a leading "$ ", environment assignments and a trailing
// comment, and reports whether what is left starts with "<prog> ".
func commandText(text, prog string) (string, bool) {
	t := strings.TrimSpace(text)
	t = strings.TrimSpace(strings.TrimPrefix(t, "$ "))
	for {
		loc := envAssign.FindStringIndex(t)
		if loc == nil {
			break
		}
		t = t[loc[1]:]
	}
	if i := strings.Index(t, " #"); i >= 0 {
		t = strings.TrimSpace(t[:i])
	}
	if !strings.HasPrefix(t, prog+" ") {
		return "", false
	}
	return strings.Join(strings.Fields(t), " "), true
}

// Flags returns the --flag names among tokens, each without an "=value".
// A bare "--" is not a flag.
func Flags(tokens []string) []string {
	var out []string
	for _, tok := range tokens {
		if !strings.HasPrefix(tok, "--") || tok == "--" {
			continue
		}
		name, _, _ := strings.Cut(tok, "=")
		out = append(out, name)
	}
	return out
}

// FlagIn reports whether help names the flag as a whole word: "--follow" is
// not found in a help that only names "--follow-up".
func FlagIn(help, flag string) bool {
	re := regexp.MustCompile(`(^|[^A-Za-z0-9-])` + regexp.QuoteMeta(flag) + `([^A-Za-z0-9-]|$)`)
	return re.MatchString(help)
}

// VerbPath is the longest of paths that tokens start with, or nil.
func VerbPath(tokens []string, paths [][]string) []string {
	var best []string
	for _, p := range paths {
		if len(p) <= len(best) || len(p) > len(tokens) {
			continue
		}
		match := true
		for i := range p {
			if tokens[i] != p[i] {
				match = false
				break
			}
		}
		if match {
			best = p
		}
	}
	return best
}

var verbWord = regexp.MustCompile(`^[a-z][a-z-]*$`)

// Verb names what tokens ask for when no usage path matches: the leading
// lowercase words, at most two ("episode frobnicate").
func Verb(tokens []string) string {
	var words []string
	for _, w := range tokens {
		if len(words) == 2 || !verbWord.MatchString(w) {
			break
		}
		words = append(words, w)
	}
	return strings.Join(words, " ")
}

// UsagePaths returns the verb path of every usage line that starts with
// prefix: its leading lowercase words, after any of skip are removed (e.g.
// ailang-worldd's "[--addr <url>] "). "  ailang-worldd session new <episode>"
// gives [session new].
func UsagePaths(usage, prefix string, skip ...string) [][]string {
	var out [][]string
	for _, l := range strings.Split(usage, "\n") {
		rest, ok := strings.CutPrefix(l, prefix)
		if !ok {
			continue
		}
		for _, s := range skip {
			rest = strings.TrimPrefix(rest, s)
		}
		var path []string
		for _, w := range strings.Fields(rest) {
			if !verbWord.MatchString(w) {
				break
			}
			path = append(path, w)
		}
		if len(path) > 0 {
			out = append(out, path)
		}
	}
	return out
}
