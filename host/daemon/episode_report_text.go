package daemon

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// WriteEpisodeReportText renders the compact human table: episode report's text at a terminal and
// in the app's bar (row 243 V1, moved from cmd/ailang-worldd so both print one function's bytes).
func WriteEpisodeReportText(w io.Writer, r EpisodeReport) {
	if r.Start != nil {
		data, _ := json.MarshalIndent(r.Start, "  ", "  ")
		fmt.Fprintf(w, "start:\n  %s\n", data)
	}
	fmt.Fprintf(w, "episode %s: %d calls, entries %d-%d, %s .. %s\n",
		r.Episode, r.TotalCalls, r.Entries[0], r.Entries[1], r.FirstAt, r.LastAt)
	fmt.Fprintln(w, "calls by tool:")
	for _, kv := range SortedCounts(r.CallsByTool) {
		fmt.Fprintf(w, "  %-18s %6d\n", kv.Name, kv.N)
	}
	fmt.Fprintln(w, readsLine(r.Reads))
	fmt.Fprintln(w, "effects:")
	for _, e := range r.Effects {
		fmt.Fprintf(w, "  %-22s %-16s %6d\n", e.Effect, e.Decision, e.Count)
	}
	fmt.Fprintln(w, "budgets (minted / used / left):")
	for _, b := range r.Budgets {
		fmt.Fprintf(w, "  %-22s %6d %6d %6d\n", b.Effect+"@"+b.Scope, b.Minted, b.Used, b.Left)
	}
	WriteReportRefusals(w, "plan refusals (refused before any effect)", r.PlanRefusals)
	WriteReportRefusals(w, "handler refusals (effect allowed, handler ok:false)", r.HandlerRefusals)
	if r.Uncommitted > 0 {
		fmt.Fprintf(w, "uncommitted intents: %d\n", r.Uncommitted)
	}
	for _, e := range r.ExecCache {
		fmt.Fprintln(w, execCacheLine(e))
	}
}

// execCacheLine is one exec_cache entry of the report table (row 257).
func execCacheLine(e ReportExecCache) string {
	line := fmt.Sprintf("exec cache   %s %s %s", e.Project, e.Cache, e.State)
	if e.Attempts > 0 {
		line += fmt.Sprintf(" at %s (attempt %d of 3)", e.By, e.Attempts)
	}
	if e.State == "ready" && e.Attempts > 0 {
		line += fmt.Sprintf(" in %dms", e.TookMS)
	}
	if e.Failure != "" {
		line += ": " + e.Failure
	}
	return line
}

// readsLine is the one-line reads summary; --json carries the file list.
func readsLine(r ReportReads) string {
	line := fmt.Sprintf("reads: %d calls, %d files", r.Count, len(r.Files))
	if r.Truncated > 0 {
		line += fmt.Sprintf(", %d more not listed", r.Truncated)
	}
	if r.Unlisted > 0 {
		line += fmt.Sprintf(", %d outside the worktree", r.Unlisted)
	}
	return line
}

// WriteReportRefusals is one refusals line and its first entries.
func WriteReportRefusals(w io.Writer, title string, r ReportRefusals) {
	parts := make([]string, 0, len(r.ByTool))
	for _, kv := range SortedCounts(r.ByTool) {
		parts = append(parts, fmt.Sprintf("%s %d", kv.Name, kv.N))
	}
	line := fmt.Sprintf("%s: %d", title, r.Count)
	if len(parts) > 0 {
		line += " (" + strings.Join(parts, ", ") + ")"
	}
	fmt.Fprintln(w, line)
	for _, f := range r.First {
		fmt.Fprintf(w, "  #%d %s: %s\n", f.Entry, f.Tool, f.Reason)
	}
}

// NamedCount is one name and its count.
type NamedCount struct {
	Name string
	N    int
}

// SortedCounts orders a count map by count, then name.
func SortedCounts(m map[string]int) []NamedCount {
	out := make([]NamedCount, 0, len(m))
	for k, n := range m {
		out = append(out, NamedCount{k, n})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].N != out[j].N {
			return out[i].N > out[j].N
		}
		return out[i].Name < out[j].Name
	})
	return out
}
