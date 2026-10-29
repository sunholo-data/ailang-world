package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/sunholo-data/ailang-world/host/daemon"
)

const reportMarkdownMax = 60000

// scrubMarkdown implements trial.py's BARE_HEX64_RE without RE2 lookbehind.
func scrubMarkdown(s string) error {
	if strings.Contains(strings.ToLower(s), "bearer ") {
		return fmt.Errorf("refusing markdown: bearer credential shape")
	}
	hex := func(c byte) bool {
		return c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
	}
	for i := 0; i < len(s); {
		if !hex(s[i]) {
			i++
			continue
		}
		start := i
		for i < len(s) && hex(s[i]) {
			i++
		}
		if i-start == 64 && (start == 0 || s[start-1] != ':') {
			return fmt.Errorf("refusing markdown: bare 64-hex value")
		}
	}
	return nil
}

func reportTrailer(r daemon.EpisodeReport, db string) string {
	return fmt.Sprintf("World-Provenance: store=%s episode=%s entries=%d-%d", db, r.Episode, r.Entries[0], r.Entries[1])
}

type finishCheckResult struct {
	Name        string
	Unsandboxed bool
	Passed      bool
	Entry       int64
}

type reportAdditions struct {
	Checks       []finishCheckResult
	Commit       string
	Unattributed []string
	Finish       bool
}

// renderReportMarkdown buffers everything, scrubs before truncation, and keeps
// the provenance trailer intact as the last line.
func renderReportMarkdown(r daemon.EpisodeReport, db string, add reportAdditions) (string, error) {
	var b strings.Builder
	agent := "unknown agent"
	if r.Start != nil {
		agent = strings.TrimSpace(r.Start.Agent + " " + r.Start.Model)
	}
	fmt.Fprintf(&b, "Episode %s: %s working through AILANG World only.\n\n", r.Episode, agent)
	if add.Finish {
		b.WriteString("## Checks\n\n")
		for _, c := range add.Checks {
			verdict, class := "PASS", "sandboxed"
			if !c.Passed {
				verdict = "FAILED"
			}
			if c.Unsandboxed {
				class = "UNSANDBOXED"
			}
			fmt.Fprintf(&b, "- %s: %s (%s), outcome entry %d\n", c.Name, verdict, class, c.Entry)
		}
		fmt.Fprintf(&b, "\nCommit: %s\n\nD-WORLD-60: apply the ailang-world label when opening the PR.\n\n", add.Commit)
	}
	introEnd := b.Len()
	b.WriteString("## What the agent did\n\n")
	fmt.Fprintf(&b, "%d calls.\n\nCalls by tool:\n", r.TotalCalls)
	for _, kv := range sortedCounts(r.CallsByTool) {
		fmt.Fprintf(&b, "- %s: %d\n", kv.name, kv.n)
	}
	b.WriteString("\nEffects by decision:\n")
	for _, e := range r.Effects {
		fmt.Fprintf(&b, "- %s %s: %d\n", e.Effect, e.Decision, e.Count)
	}
	b.WriteString("\n")
	printRefusals(&b, "Plan refusals", r.PlanRefusals)
	printRefusals(&b, "Handler refusals", r.HandlerRefusals)
	b.WriteString("\nBudgets (used / left):\n")
	for _, v := range r.Budgets {
		fmt.Fprintf(&b, "- %s@%s: %d / %d\n", v.Effect, v.Scope, v.Used, v.Left)
	}
	activityEnd := b.Len()
	b.WriteString("\n## Files\n\nWritten:\n")
	for _, v := range r.Writes.Written {
		fmt.Fprintf(&b, "- %s (%s), %d calls, first entry %d\n", v.Path, v.Op, v.Count, v.FirstEntry)
	}
	b.WriteString("\nProduced:\n")
	for _, v := range r.Writes.Produced {
		fmt.Fprintf(&b, "- %s, %d calls, first entry %d\n", v.Path, v.Count, v.FirstEntry)
	}
	fmt.Fprintf(&b, "\nEvidence gaps: %d written and %d produced paths beyond the report cap; %d nonlocal paths unlisted.\n",
		r.Writes.WrittenTruncated, r.Writes.ProducedTruncated, r.Writes.Unlisted)
	fmt.Fprintf(&b, "Evidence entries beyond the cap: %d truncated deltas, %d delta errors.\n", r.Writes.DeltaTruncatedOmitted, r.Writes.DeltaErrorsOmitted)
	for _, e := range r.Writes.DeltasTruncated {
		fmt.Fprintf(&b, "- deltas_truncated: entry %d\n", e)
	}
	for _, e := range r.Writes.DeltaErrors {
		fmt.Fprintf(&b, "- delta_errors: entry %d: %s\n", e.Entry, e.Error)
	}
	if len(add.Unattributed) > 0 {
		b.WriteString("\n### Not made through World\n\n")
		for _, p := range add.Unattributed {
			fmt.Fprintf(&b, "- %s\n", p)
		}
	}
	fmt.Fprintf(&b, "\nFiles read: %d (%d read calls; %d outside the worktree).\n", len(r.Reads.Files)+r.Reads.Truncated, r.Reads.Count, r.Reads.Unlisted)
	tail := fmt.Sprintf("\nWorld entries %d-%d; walk any of them with ailang-worldd why <entry>.\n\n%s\n",
		r.Entries[0], r.Entries[1], reportTrailer(r, db))
	full := b.String() + tail
	if err := scrubMarkdown(full); err != nil {
		return "", err
	}
	if len(full) <= reportMarkdownMax {
		return full, nil
	}
	note := fmt.Sprintf("\n\nSection truncated; run ailang-worldd episode report %s for the full form.\n", r.Episode)
	room := (reportMarkdownMax - 3*len(note) - len(tail)) / 3
	if room < 0 {
		return "", fmt.Errorf("markdown provenance exceeds %d bytes", reportMarkdownMax)
	}
	body := b.String()
	sections := []string{body[:introEnd], body[introEnd:activityEnd], body[activityEnd:]}
	var capped strings.Builder
	for _, section := range sections {
		if len(section) > room {
			section = section[:room]
			for !utf8.ValidString(section) && len(section) > 0 {
				section = section[:len(section)-1]
			}
			section += note
		}
		capped.WriteString(section)
	}
	return capped.String() + tail, nil
}
