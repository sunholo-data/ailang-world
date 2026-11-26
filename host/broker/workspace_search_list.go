package broker

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

const (
	defaultSearchMaxResults = 100
	maxSearchMaxResults     = 500
	maxSearchPatternBytes   = 512
	searchWindowRunes       = 240
	searchBinaryProbeBytes  = 8192
	maxSearchFileBytes      = 32 << 20 // 32 MiB per file safety cap

	defaultListDepth      = 2
	maxListDepth          = 6
	defaultListMaxEntries = 300
	maxListMaxEntries     = 1000

	// responseEnvelopeHeadroom reserves room for the response envelope keys
	// (ok, total_matches/total, files_searched, truncated, next, tool,
	// policy_digest, package_cache, world) so the full JSON stays <= defaultReadMaxChars.
	responseEnvelopeHeadroom = 1024
)

var (
	searchAllowedKeys = map[string]bool{
		"op":          true,
		"pattern":     true,
		"regex":       true,
		"ignore_case": true,
		"path":        true,
		"include":     true,
		"max_results": true,
	}
	listAllowedKeys = map[string]bool{
		"op":          true,
		"path":        true,
		"depth":       true,
		"max_entries": true,
		"include":     true,
	}
)

// threePartRefusal formats a row-227 three-part refusal in policy-tool's
// ok:false shape (w-agent-search-and-read-cap §3.1).
func threePartRefusal(what, why, next string) map[string]json.RawMessage {
	msg := fmt.Sprintf("refused: %s\n  why: %s\n  next: %s", what, why, next)
	refused, _ := json.Marshal(msg)
	return map[string]json.RawMessage{
		"ok":      json.RawMessage("false"),
		"refused": refused,
	}
}

func checkAllowedKeys(fields map[string]json.RawMessage, allowed map[string]bool, toolName string) map[string]json.RawMessage {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if !allowed[k] {
			allowedList := make([]string, 0, len(allowed))
			for ak := range allowed {
				if ak != "op" {
					allowedList = append(allowedList, ak)
				}
			}
			sort.Strings(allowedList)
			return threePartRefusal(
				fmt.Sprintf("unknown argument %q", k),
				fmt.Sprintf("%s rejects keys outside its schema (including case-folded aliases)", toolName),
				fmt.Sprintf("retry %s using only: %s", toolName, strings.Join(allowedList, ", ")),
			)
		}
	}
	return nil
}

func parseOptionalString(fields map[string]json.RawMessage, key, toolName string) (string, bool, map[string]json.RawMessage) {
	raw, ok := fields[key]
	if !ok {
		return "", false, nil
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return "", false, threePartRefusal(
			fmt.Sprintf("argument %q must be a string", key),
			fmt.Sprintf("%s requires a JSON string for %q", toolName, key),
			fmt.Sprintf("pass a string value for %q", key),
		)
	}
	return s, true, nil
}

func parseOptionalBool(fields map[string]json.RawMessage, key, toolName string) (bool, map[string]json.RawMessage) {
	raw, ok := fields[key]
	if !ok {
		return false, nil
	}
	var b bool
	if err := json.Unmarshal(raw, &b); err != nil {
		return false, threePartRefusal(
			fmt.Sprintf("argument %q must be a boolean", key),
			fmt.Sprintf("%s requires true or false for %q", toolName, key),
			fmt.Sprintf("pass true or false for %q", key),
		)
	}
	return b, nil
}

func parseOptionalIntRange(fields map[string]json.RawMessage, key string, def, minVal, maxVal int, toolName string) (int, map[string]json.RawMessage) {
	raw, ok := fields[key]
	if !ok {
		return def, nil
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err != nil || float64(int(f)) != f {
		return 0, threePartRefusal(
			fmt.Sprintf("argument %q must be an integer", key),
			fmt.Sprintf("%s requires an integer between %d and %d for %q", toolName, minVal, maxVal, key),
			fmt.Sprintf("pass an integer in %d..%d for %q", minVal, maxVal, key),
		)
	}
	n := int(f)
	if n < minVal || n > maxVal {
		return 0, threePartRefusal(
			fmt.Sprintf("argument %q=%d is out of range (%d..%d)", key, n, minVal, maxVal),
			fmt.Sprintf("%s bounds %q to %d..%d to keep results bounded", toolName, key, minVal, maxVal),
			fmt.Sprintf("pass an integer in %d..%d for %q", minVal, maxVal, key),
		)
	}
	return n, nil
}

func parseIncludeGlob(fields map[string]json.RawMessage, toolName string) (string, map[string]json.RawMessage) {
	inc, present, ref := parseOptionalString(fields, "include", toolName)
	if ref != nil || !present {
		return "", ref
	}
	if inc == "" || strings.ContainsRune(inc, 0) || strings.Contains(inc, "/") {
		return "", threePartRefusal(
			fmt.Sprintf("invalid include glob %q", inc),
			"include must be a non-empty basename glob without '/' or NUL",
			`pass a basename glob such as "*.go" or "*.md"`,
		)
	}
	if _, err := path.Match(inc, "probe"); err != nil {
		return "", threePartRefusal(
			fmt.Sprintf("invalid include glob %q: %v", inc, err),
			"include must be a valid path.Match glob pattern",
			`fix the glob syntax, e.g. "*.go" or "*.md*"`,
		)
	}
	return inc, nil
}

func isGitSegment(name string) bool {
	return strings.EqualFold(name, ".git")
}

// resolveWorktreeTarget validates relPath and resolves it strictly inside
// h.root without following any symlink along the path. It returns the absolute
// path, the clean slash-relative path from h.root ("." for the root), and the
// os.FileInfo from Lstat.
func (h *AilangToolHandler) resolveWorktreeTarget(relPath, toolName string) (string, string, os.FileInfo, map[string]json.RawMessage) {
	if relPath == "" {
		relPath = "."
	}
	if strings.ContainsRune(relPath, 0) || strings.HasPrefix(relPath, "/") || strings.HasPrefix(relPath, `\`) || filepath.IsAbs(relPath) || strings.HasPrefix(relPath, "-") {
		return "", "", nil, threePartRefusal(
			fmt.Sprintf("path %q escapes or violates the worktree sandbox", relPath),
			"paths must be relative to the worktree root with no leading '/', '\\', '-' or NUL",
			fmt.Sprintf("retry %s with a relative path inside the worktree (default \".\")", toolName),
		)
	}
	// Check raw segments before Clean so "a/../../b" or "a/.." never slips through.
	parts := strings.Split(filepath.ToSlash(relPath), "/")
	for _, seg := range parts {
		if seg == ".." {
			return "", "", nil, threePartRefusal(
				fmt.Sprintf("path %q contains \"..\" and escapes the worktree sandbox", relPath),
				"\"..\" segments are forbidden so reads stay inside the episode worktree",
				fmt.Sprintf("retry %s with a relative path without \"..\" segments", toolName),
			)
		}
		if isGitSegment(seg) {
			return "", "", nil, threePartRefusal(
				fmt.Sprintf("path %q touches .git metadata", relPath),
				".git is excluded from workspace search and list",
				fmt.Sprintf("retry %s on a non-.git path inside the worktree", toolName),
			)
		}
	}

	cleaned := path.Clean(filepath.ToSlash(relPath))
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return "", "", nil, threePartRefusal(
			fmt.Sprintf("path %q escapes the worktree sandbox", relPath),
			"target resolves outside the episode worktree",
			fmt.Sprintf("retry %s with a path inside the worktree", toolName),
		)
	}
	if cleaned == "." {
		info, err := os.Lstat(h.root)
		if err != nil {
			return "", "", nil, threePartRefusal(
				fmt.Sprintf("cannot stat worktree root: %v", err),
				"the worktree root is not accessible",
				"check the episode worktree state",
			)
		}
		return h.root, ".", info, nil
	}

	// Walk every component from h.root with Lstat so no symlink is ever followed.
	cur := h.root
	for _, seg := range strings.Split(cleaned, "/") {
		if seg == "" || seg == "." {
			continue
		}
		cur = filepath.Join(cur, seg)
		info, err := os.Lstat(cur)
		if err != nil {
			if os.IsNotExist(err) {
				return "", "", nil, threePartRefusal(
					fmt.Sprintf("path %q does not exist", relPath),
					"no file or directory exists at that worktree path",
					fmt.Sprintf("call workspace-list {path: \".\"} to inspect existing paths"),
				)
			}
			return "", "", nil, threePartRefusal(
				fmt.Sprintf("cannot access %q: %v", relPath, err),
				"lstat failed on the target path",
				fmt.Sprintf("retry %s on an accessible worktree path", toolName),
			)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", "", nil, threePartRefusal(
				fmt.Sprintf("path %q is or traverses a symlink", relPath),
				fmt.Sprintf("%s never follows symlinks so reads stay confined to the worktree", toolName),
				fmt.Sprintf("call %s on the containing directory or a real worktree path", toolName),
			)
		}
	}

	resolved, err := filepath.EvalSymlinks(cur)
	if err != nil || !within(resolved, h.root) {
		return "", "", nil, threePartRefusal(
			fmt.Sprintf("path %q resolves outside the worktree sandbox", relPath),
			"symlink-resolved path is not inside the episode worktree",
			fmt.Sprintf("retry %s with a path inside the worktree", toolName),
		)
	}
	info, err := os.Lstat(cur)
	if err != nil {
		return "", "", nil, threePartRefusal(
			fmt.Sprintf("cannot stat %q: %v", relPath, err),
			"the target path could not be stated",
			fmt.Sprintf("retry %s with a valid worktree path", toolName),
		)
	}
	return cur, cleaned, info, nil
}

type searchMatch struct {
	Path string `json:"path"`
	Line int    `json:"line"`
	Col  int    `json:"col"`
	Text string `json:"text"`
}

func (h *AilangToolHandler) executeSearch(ctx context.Context, fields map[string]json.RawMessage) ([]byte, error) {
	if ref := checkAllowedKeys(fields, searchAllowedKeys, "workspace-search"); ref != nil {
		return h.stamp(ref)
	}
	pattern, present, ref := parseOptionalString(fields, "pattern", "workspace-search")
	if ref != nil {
		return h.stamp(ref)
	}
	if !present || len(pattern) == 0 || len(pattern) > maxSearchPatternBytes || strings.ContainsRune(pattern, 0) {
		return h.stamp(threePartRefusal(
			"pattern must be 1 to 512 bytes with no NUL",
			"workspace-search requires a non-empty bounded search pattern",
			`pass pattern as a 1..512 byte string, e.g. workspace-search {pattern: "func main"}`,
		))
	}
	useRegex, ref := parseOptionalBool(fields, "regex", "workspace-search")
	if ref != nil {
		return h.stamp(ref)
	}
	ignoreCase, ref := parseOptionalBool(fields, "ignore_case", "workspace-search")
	if ref != nil {
		return h.stamp(ref)
	}
	relPath, _, ref := parseOptionalString(fields, "path", "workspace-search")
	if ref != nil {
		return h.stamp(ref)
	}
	include, ref := parseIncludeGlob(fields, "workspace-search")
	if ref != nil {
		return h.stamp(ref)
	}
	maxResults, ref := parseOptionalIntRange(fields, "max_results", defaultSearchMaxResults, 1, maxSearchMaxResults, "workspace-search")
	if ref != nil {
		return h.stamp(ref)
	}

	var matcher func(line string) (int, int, bool)
	if useRegex || ignoreCase {
		expr := pattern
		if !useRegex {
			expr = regexp.QuoteMeta(pattern)
		}
		if ignoreCase {
			expr = "(?i)" + expr
		}
		re, err := regexp.Compile(expr)
		if err != nil {
			return h.stamp(threePartRefusal(
				fmt.Sprintf("invalid regex %q: %v", pattern, err),
				"workspace-search uses Go RE2 syntax when regex is true",
				"fix the RE2 expression or set regex: false for a fixed-string search",
			))
		}
		matcher = func(line string) (int, int, bool) {
			loc := re.FindStringIndex(line)
			if loc == nil {
				return 0, 0, false
			}
			return loc[0], loc[1], true
		}
	} else {
		matcher = func(line string) (int, int, bool) {
			idx := strings.Index(line, pattern)
			if idx < 0 {
				return 0, 0, false
			}
			return idx, idx + len(pattern), true
		}
	}

	fullPath, cleanRel, info, ref := h.resolveWorktreeTarget(relPath, "workspace-search")
	if ref != nil {
		return h.stamp(ref)
	}

	matches := make([]searchMatch, 0)
	totalMatches := 0
	filesSearched := 0
	truncated := false
	usedRunes := 0
	runeBudget := defaultReadMaxChars - responseEnvelopeHeadroom

	scanFile := func(absPath, slashRel string, fi os.FileInfo) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if !fi.Mode().IsRegular() || fi.Size() > maxSearchFileBytes {
			return nil
		}
		if include != "" {
			ok, _ := path.Match(include, path.Base(slashRel))
			if !ok {
				return nil
			}
		}
		f, err := os.Open(absPath)
		if err != nil {
			return nil
		}
		defer func() { _ = f.Close() }()

		// Probe the leading bytes for NUL to skip binary files.
		var probe [searchBinaryProbeBytes]byte
		n, readErr := f.Read(probe[:])
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			return nil
		}
		if bytes.IndexByte(probe[:n], 0) >= 0 {
			return nil
		}
		filesSearched++

		reader := bufio.NewReader(io.MultiReader(bytes.NewReader(probe[:n]), f))
		lineNum := 0
		for {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			rawLine, err := reader.ReadString('\n')
			if err != nil && !errors.Is(err, io.EOF) {
				break
			}
			if rawLine == "" && errors.Is(err, io.EOF) {
				break
			}
			lineNum++
			line := strings.TrimRight(rawLine, "\r\n")
			if start, end, matched := matcher(line); matched {
				totalMatches++
				if len(matches) >= maxResults {
					truncated = true
					return errSearchStop
				}
				col := utf8.RuneCountInString(line[:start]) + 1
				item := searchMatch{
					Path: slashRel,
					Line: lineNum,
					Col:  col,
					Text: windowLine(line, start, end, searchWindowRunes),
				}
				encoded, _ := json.Marshal(item)
				itemRunes := utf8.RuneCount(encoded) + 1 // comma
				if usedRunes+itemRunes > runeBudget {
					truncated = true
					return errSearchStop
				}
				usedRunes += itemRunes
				matches = append(matches, item)
			}
			if errors.Is(err, io.EOF) {
				break
			}
		}
		return nil
	}

	if !info.IsDir() {
		err := scanFile(fullPath, cleanRel, info)
		if err != nil && !errors.Is(err, errSearchStop) {
			return h.stamp(threePartRefusal(
				fmt.Sprintf("search interrupted: %v", err),
				"the search deadline or context expired before completing",
				"narrow path or include and retry workspace-search",
			))
		}
	} else {
		var walk func(absDir, relDir string) error
		walk = func(absDir, relDir string) error {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			entries, err := os.ReadDir(absDir)
			if err != nil {
				return nil
			}
			for _, ent := range entries {
				name := ent.Name()
				if isGitSegment(name) {
					continue
				}
				// Never follow symlinks.
				if ent.Type()&os.ModeSymlink != 0 {
					continue
				}
				childAbs := filepath.Join(absDir, name)
				childRel := name
				if relDir != "." {
					childRel = relDir + "/" + name
				}
				if ent.IsDir() {
					if err := walk(childAbs, childRel); err != nil {
						return err
					}
					continue
				}
				fi, err := ent.Info()
				if err != nil || !fi.Mode().IsRegular() {
					continue
				}
				if err := scanFile(childAbs, childRel, fi); err != nil {
					return err
				}
			}
			return nil
		}
		err := walk(fullPath, cleanRel)
		if err != nil && !errors.Is(err, errSearchStop) {
			return h.stamp(threePartRefusal(
				fmt.Sprintf("search interrupted: %v", err),
				"the search deadline or context expired before completing",
				"narrow path or include and retry workspace-search",
			))
		}
	}

	resp := map[string]json.RawMessage{
		"ok": json.RawMessage("true"),
	}
	resp["matches"], _ = json.Marshal(matches)
	resp["total_matches"], _ = json.Marshal(totalMatches)
	resp["files_searched"], _ = json.Marshal(filesSearched)
	resp["truncated"], _ = json.Marshal(truncated)
	if truncated {
		resp["next"], _ = json.Marshal("narrow path or include")
	} else if totalMatches == 0 && !useRegex && strings.ContainsAny(pattern, "|*?+^$[]()\\") {
		resp["hint"], _ = json.Marshal("literal search found 0 matches; pass regex: true for Go RE2 regex matching (e.g. a|b)")
	}
	return h.stamp(resp)
}

var errSearchStop = errors.New("search cap reached")

// windowLine returns line cut to at most maxRunes runes centered around the
// byte span [startByte, endByte], placing '…' at any cut end.
func windowLine(line string, startByte, endByte, maxRunes int) string {
	runes := []rune(line)
	n := len(runes)
	if n <= maxRunes {
		return line
	}
	mStart := utf8.RuneCountInString(line[:startByte])
	mEnd := utf8.RuneCountInString(line[:endByte])
	mid := (mStart + mEnd) / 2
	winStart := mid - maxRunes/2
	if winStart > mStart {
		winStart = mStart
	}
	if winStart < 0 {
		winStart = 0
	}

	needLeft := winStart > 0
	avail := maxRunes
	if needLeft {
		avail--
	}
	if winStart+avail >= n {
		// Fits to the end of the line: shift window left to fill maxRunes.
		winStart = n - maxRunes
		if winStart < 0 {
			winStart = 0
		}
		if winStart > 0 {
			return "…" + string(runes[winStart+1:])
		}
		return string(runes)
	}
	// Cut at the right end too.
	take := avail - 1
	if take < 0 {
		take = 0
	}
	if needLeft {
		return "…" + string(runes[winStart:winStart+take]) + "…"
	}
	return string(runes[winStart:winStart+take]) + "…"
}

type listEntry struct {
	Path string `json:"path"`
	Type string `json:"type"`
	Size int64  `json:"size"`
}

func (h *AilangToolHandler) executeList(ctx context.Context, fields map[string]json.RawMessage) ([]byte, error) {
	if ref := checkAllowedKeys(fields, listAllowedKeys, "workspace-list"); ref != nil {
		return h.stamp(ref)
	}
	relPath, _, ref := parseOptionalString(fields, "path", "workspace-list")
	if ref != nil {
		return h.stamp(ref)
	}
	depth, ref := parseOptionalIntRange(fields, "depth", defaultListDepth, 1, maxListDepth, "workspace-list")
	if ref != nil {
		return h.stamp(ref)
	}
	maxEntries, ref := parseOptionalIntRange(fields, "max_entries", defaultListMaxEntries, 1, maxListMaxEntries, "workspace-list")
	if ref != nil {
		return h.stamp(ref)
	}
	include, ref := parseIncludeGlob(fields, "workspace-list")
	if ref != nil {
		return h.stamp(ref)
	}

	fullPath, cleanRel, info, ref := h.resolveWorktreeTarget(relPath, "workspace-list")
	if ref != nil {
		return h.stamp(ref)
	}

	entries := make([]listEntry, 0)
	truncated := false
	usedRunes := 0
	runeBudget := defaultReadMaxChars - responseEnvelopeHeadroom

	addEntry := func(e listEntry) bool {
		if len(entries) >= maxEntries {
			truncated = true
			return false
		}
		encoded, _ := json.Marshal(e)
		itemRunes := utf8.RuneCount(encoded) + 1
		if usedRunes+itemRunes > runeBudget {
			truncated = true
			return false
		}
		usedRunes += itemRunes
		entries = append(entries, e)
		return true
	}

	if !info.IsDir() {
		kind := "file"
		if info.Mode()&os.ModeSymlink != 0 {
			kind = "symlink"
		}
		if include == "" || func() bool { ok, _ := path.Match(include, path.Base(cleanRel)); return ok }() {
			addEntry(listEntry{Path: cleanRel, Type: kind, Size: info.Size()})
		}
	} else {
		var walk func(absDir, relDir string, curDepth int) error
		walk = func(absDir, relDir string, curDepth int) error {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			dirEntries, err := os.ReadDir(absDir)
			if err != nil {
				return nil
			}
			for _, ent := range dirEntries {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				name := ent.Name()
				if isGitSegment(name) {
					continue
				}
				childRel := name
				if relDir != "." {
					childRel = relDir + "/" + name
				}
				fi, err := ent.Info()
				if err != nil {
					continue
				}
				var kind string
				var size int64
				switch {
				case ent.Type()&os.ModeSymlink != 0:
					kind = "symlink"
					size = fi.Size()
				case ent.IsDir():
					kind = "dir"
					size = 0
				case fi.Mode().IsRegular():
					kind = "file"
					size = fi.Size()
				default:
					continue
				}

				matchesInclude := true
				if include != "" {
					matchesInclude, _ = path.Match(include, name)
				}
				if matchesInclude {
					if !addEntry(listEntry{Path: childRel, Type: kind, Size: size}) {
						return errSearchStop
					}
				}
				// Descend into real directories only, never symlinks, up to depth.
				if kind == "dir" && curDepth < depth {
					if err := walk(filepath.Join(absDir, name), childRel, curDepth+1); err != nil {
						return err
					}
				}
			}
			return nil
		}
		err := walk(fullPath, cleanRel, 1)
		if err != nil && !errors.Is(err, errSearchStop) {
			return h.stamp(threePartRefusal(
				fmt.Sprintf("list interrupted: %v", err),
				"the list deadline or context expired before completing",
				"narrow path, depth or include and retry workspace-list",
			))
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})

	resp := map[string]json.RawMessage{
		"ok": json.RawMessage("true"),
	}
	resp["entries"], _ = json.Marshal(entries)
	resp["total"], _ = json.Marshal(len(entries))
	resp["truncated"], _ = json.Marshal(truncated)
	if truncated {
		resp["next"], _ = json.Marshal("narrow path, depth or include")
	}
	return h.stamp(resp)
}
