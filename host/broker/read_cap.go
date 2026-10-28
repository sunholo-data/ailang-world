package broker

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"
)

const defaultReadMaxChars = 40000

// capReadResponse withholds content only after the read has happened. Returning
// a normal response preserves its one charged, committed Workspace.Read effect.
func capReadResponse(resp, request map[string]json.RawMessage, cap int) map[string]json.RawMessage {
	var content, path string
	_ = json.Unmarshal(resp["content"], &content)
	if utf8.RuneCountInString(content) <= cap {
		return resp
	}
	_ = json.Unmarshal(request["path"], &path)
	offset := 1
	_ = json.Unmarshal(request["offset"], &offset)
	if offset < 1 {
		offset = 1
	}
	// Retain newline characters in the measurement, without inventing an
	// additional line for a trailing newline. Plan only the returned range.
	lines := strings.SplitAfter(content, "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	var next []string
	for i := 0; i < len(lines); {
		if utf8.RuneCountInString(lines[i]) > cap {
			next = append(next, fmt.Sprintf("line %d cannot be read whole yet; workspace-search will return a window of it (M2)", offset+i))
			i++
			continue
		}
		start, count := i, 0
		for i < len(lines) {
			n := utf8.RuneCountInString(lines[i])
			if count+n > cap {
				break
			}
			count += n
			i++
		}
		next = append(next, fmt.Sprintf("ailang-read {path: %q, offset: %d, limit: %d}", path, offset+start, i-start))
	}
	msg := fmt.Sprintf("refused: %s is %d characters over %d lines; the cap is %d\n  why: one read may not exceed max_chars (default 40,000) so the result fits the agent's context\n  next: %s", path, utf8.RuneCountInString(content), len(lines), cap, strings.Join(next, ", then "))
	refused, _ := json.Marshal(msg)
	return map[string]json.RawMessage{"ok": json.RawMessage("false"), "refused": refused}
}
