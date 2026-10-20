package broker

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// protectedGitPreflight refuses config-defined host commands before status.
func protectedGitPreflight(ctx context.Context, git, root string) error {
	filters, err := configuredGitFilters(ctx, git, root)
	if err != nil {
		return fmt.Errorf("protected Git filter query: %w", err)
	}
	if len(filters) > 0 {
		return fmt.Errorf("%s: %s (git status would run a configured command)",
			WorktreeFilterRefused, strings.Join(filters, ", "))
	}
	return nil
}

// protectedCheckTimeout is the provisional phase bound; M6-pre (V24) resizes it.
var protectedCheckTimeout = 2500 * time.Millisecond

type protectedCheck struct {
	intact  *bool
	changes []protectedChange
	err     error
}

// checkProtected reports integrity without changing the command's result.
func checkProtected(ctx context.Context, git, root, base string, canonical []string, ignored map[string]struct{}) (result protectedCheck) {
	ctx, cancel := context.WithTimeout(ctx, protectedCheckTimeout)
	defer cancel()
	if result.err = protectedGitPreflight(ctx, git, root); result.err != nil {
		return
	}
	commands := [][]string{
		{"status", "--porcelain=v2", "-z", "--branch", "--untracked-files=all", "--ignored=no", "--no-renames", "--ignore-submodules=all", "--"},
		strings.Fields(protectedIgnoredCommand),
		{"ls-files", "-v", "-z", "--"},
	}
	var streams [3][]byte
	for i, args := range commands {
		args = append(append([]string{"-c", "core.fsmonitor=false"}, args...), gitCandidatePathspecs(canonical)...)
		stderr := newBoundedSink(protectedCheckMaxOutput)
		out, err := runBounded(ctx, handlerBounds{execTimeout: protectedCheckTimeout, maxOutputBytes: protectedCheckMaxOutput},
			handlerCommand{path: git, args: args, dir: root, env: gitEnv(), stderr: stderr})
		if err != nil {
			result.err = fmt.Errorf("protected Git %s: %w (stderr %q)", commands[i][0], err, stderr.Bytes())
			return
		}
		streams[i] = out
	}
	head, changes, err := parseProtectedCandidates(canonical, streams[0], streams[1], streams[2], ignored)
	if err != nil || head == "" {
		result.err = fmt.Errorf("protected Git candidates: malformed output (HEAD %q): %v", head, err)
		return
	}
	if head != base {
		changes = append(changes, protectedChange{Path: "HEAD", Status: "head-moved"})
	}
	intact := len(changes) == 0
	result.intact, result.changes = &intact, changes
	return
}
