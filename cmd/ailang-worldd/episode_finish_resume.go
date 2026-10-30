package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
)

// finishMode reads only durable state. A pending commit takes precedence over
// uncommitted edits; only a completed commit can begin a new fix-up round.
func finishMode(message, episode, head string, dirty bool, record *daemon.EpisodeFinish) (string, string) {
	trailer := ""
	for _, line := range strings.Split(message, "\n") {
		if strings.HasPrefix(line, "World-Provenance: ") && strings.Contains(line, " episode="+episode+" ") {
			trailer = line
		}
	}
	if trailer == "" {
		return "fresh", ""
	}
	if record != nil && record.Commit == head {
		if record.Outcome == "done" {
			if dirty {
				return "fresh", trailer
			}
			return "finished", trailer
		}
		switch record.Outcome {
		case "push-failed", "pr-failed", "ambiguous":
			return "resume-dispatch", trailer
		}
	}
	return "resume", trailer
}

// finishChangeBase excludes earlier rounds on this branch, even those owned
// by another episode, while never walking a merged branch's history.
func finishChangeBase(query func(...string) (string, error), start *store.EpisodeStart) (string, error) {
	history, err := query("log", "--first-parent", "--format=%H%x00%B%x00", start.Base.SHA+"..HEAD")
	if err != nil {
		return "", err
	}
	parts := strings.Split(history, "\x00")
	for i := 0; i+1 < len(parts); i += 2 {
		for _, line := range strings.Split(parts[i+1], "\n") {
			if strings.HasPrefix(line, "World-Provenance:") {
				return strings.TrimSpace(parts[i]), nil
			}
		}
	}
	return start.Base.SHA, nil
}

func finishStatusCount(status string) int {
	parts := strings.Split(status, "\x00")
	n := 0
	for i := 0; i < len(parts); i++ {
		if parts[i] == "" {
			continue
		}
		n++
		if len(parts[i]) > 1 && (strings.ContainsAny(parts[i][:2], "RC")) {
			i++ // -z renames/copies have a second path field.
		}
	}
	return n
}

var finishBodyCheck = regexp.MustCompile(`outcome entry [0-9]+ \(([^)]+)\)$`)

func finishResumeBody(file, sha, trailer string, record *daemon.EpisodeFinish) (string, []finishCheckResult, []string, []string, error) {
	bad := func(reason string) (string, []finishCheckResult, []string, []string, error) {
		return "", nil, nil, nil, fmt.Errorf("the PR body %s %s", file, reason)
	}
	info, err := os.Lstat(file)
	if err != nil || !info.Mode().IsRegular() {
		return bad("is not for " + sha)
	}
	data, err := os.ReadFile(file)
	body := string(data)
	lines := strings.Split(strings.TrimRight(body, "\n\r \t"), "\n")
	if err != nil || lines[len(lines)-1] != trailer || !strings.Contains("\n"+body, "\nCommit: "+sha+"\n") || scrubMarkdown(body) != nil {
		return bad("is not for " + sha)
	}
	if record != nil && record.Commit == sha && finishSHA(data) != record.ReportSHA256 {
		return bad("changed since it was recorded")
	}
	var checks []finishCheckResult
	var unattributed, protected []string
	section := ""
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "Files read:") {
			section = line
		}
		if section == "## Checks" {
			if m := finishBodyCheck.FindStringSubmatch(line); m != nil {
				checks = append(checks, finishCheckResult{Ref: m[1]})
			}
		}
		if section == "### Not made through World" && strings.HasPrefix(line, "- ") {
			unattributed = append(unattributed, strings.TrimPrefix(line, "- "))
		}
		if strings.HasPrefix(line, "Protected paths changed: ") {
			protected = strings.Split(strings.TrimPrefix(line, "Protected paths changed: "), ", ")
		}
	}
	return body, checks, unattributed, protected, nil
}

func finishReusableApproval(record *daemon.EpisodeFinish, mode string, want broker.FinishApprovalScope, hasLabel bool, now int64) (broker.FinishApprovalScope, bool) {
	if mode != "resume-dispatch" || record == nil {
		return broker.FinishApprovalScope{}, false
	}
	old, err := broker.ParseFinishApprovalScope(record.Scope)
	if err != nil || old.ExpiresAt < now+300 || (!old.CreateLabel && !hasLabel) || record.Approval == "" {
		return old, false
	}
	comparable := old
	comparable.ExpiresAt, comparable.CreateLabel = want.ExpiresAt, want.CreateLabel
	return old, comparable == want
}

// finishDetach uses the existing bounded worktree.Output/runGit boundary.
// The worktree package is outside this milestone's allowed edit paths.
func finishDetach(query func(...string) (string, error), wt string, stderr io.Writer) {
	if _, err := query("-c", "core.hooksPath=/dev/null", "checkout", "--detach"); err != nil {
		fmt.Fprintf(stderr, "warning: git -C %s checkout --detach: %s\n", wt, strings.ReplaceAll(err.Error(), "\n", " "))
	}
}

func finishStep(step string) {
	if finishAfterStep != nil {
		finishAfterStep(step)
	}
}

func finishRevoke(ctx context.Context, credentials, ep string, stdout io.Writer) error {
	creds, err := store.OpenCredentials(credentials)
	if err != nil {
		return err
	}
	defer creds.Close()
	err = revokeEpisodeCredentials(ctx, creds, ep, stdout)
	if errors.Is(err, errNoEpisodeCredential) {
		fmt.Fprintln(stdout, "session already closed")
		return nil
	}
	return err
}
