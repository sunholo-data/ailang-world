package broker

// Row 243 V3 sprint 1 M3 (design §2.2): the finish prompt's text as a pure function of the approval
// scope, the title, the body, the protected paths and the committer, moved verbatim from
// cmd/ailang-worldd finishPrompt (which now prints FinishPromptText and its question), so the
// terminal and a stored verb request carry the same bytes.

import (
	"fmt"
	"strings"
)

// DisplayCommitter is how a prompt shows a pinned committer: "" is the repository's git config.
func DisplayCommitter(v string) string {
	if v == "" {
		return "(repository git config)"
	}
	return v
}

// FinishPromptText is the text `episode finish` prints before "Open this pull request? [y/N] ":
// the committer, the scope's repository, branch, sha and base, the title, the body's length and
// first six lines, the label, draft or ready, and the protected paths changed.
func FinishPromptText(s FinishApprovalScope, title, body string, protected []string, committer string) string {
	lines := strings.Split(body, "\n")
	preview := strings.Join(lines[:min(6, len(lines))], "\n")
	state := "draft"
	if !s.Draft {
		state = "ready"
	}
	label := "ailang-world"
	if s.CreateLabel {
		label += " (will be created)"
	}
	return fmt.Sprintf("committer: %s\n", DisplayCommitter(committer)) +
		fmt.Sprintf("repository: %s\nbranch: %s\nsha: %s\nbase: %s\ntitle: %s\nbody (%d bytes), first lines:\n%s\nlabel: %s\n%s\nprotected paths changed: %s\n",
			s.Repository, s.Head, s.SHA, s.Base, title, len(body), preview, label, state, strings.Join(protected, ", "))
}
