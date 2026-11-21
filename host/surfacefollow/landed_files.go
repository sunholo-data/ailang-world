package surfacefollow

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

// Bound retained paths separately from the Git output and manifest caps.
// Account for worst-case JSON escaping of each byte of a path.
const (
	landedFilesBudget = 128 << 10
	landedFilesRetry  = time.Minute
)

// withLandedFiles runs one bounded object read from the partition's pinned
// merge IDs. Files are optional, so a failed read leaves the graph available.
func (c *Collector) withLandedFiles(ctx context.Context, r Repo, p surface.FollowPartition) surface.FollowPartition {
	if deadline, ok := ctx.Deadline(); ok && time.Until(deadline) < CallTimeout+100*time.Millisecond {
		return p
	}
	ids := make([]string, 0, len(p.Index.States)/2)
	for i := len(p.Index.States) - 1; i >= 0; i-- {
		if s := p.Index.States[i]; s.Kind == surface.KindShared {
			ids = append(ids, s.ID)
		}
	}
	p.LandedFilesChecked = true
	if len(ids) == 0 {
		return p
	}
	args := []string{"--no-lazy-fetch", "--no-replace-objects", "show", "--format=%H", "--name-status", "-z", "--diff-merges=first-parent",
		"--no-renames", "--no-ext-diff", "--no-textconv", "--no-notes", "--no-show-signature",
		"--no-color", "--no-relative", "--ignore-submodules=none"}
	args = append(args, ids...)
	args = append(args, "--")
	out, err := c.call(ctx, r, false, args...)
	switch {
	case errors.Is(err, errBudget), errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
		p.LandedFilesChecked = false
	case err == nil:
		p.Index.LandedFiles = parseLandedFiles(out, ids)
		if len(p.Index.LandedFiles) > 0 {
			if data, err := surface.Marshal(p.Index); err != nil || len(data) > surface.MaxBytes {
				p.Index.LandedFiles = nil
			}
		}
	}
	return p
}

// Git terminates every field with NUL and prefixes the first status of a
// commit with LF. Treat the following path as opaque, even if it looks like
// a status or a commit hash. Require the full expected set before counting.
func parseLandedFiles(out string, ids []string) map[string]surface.FileChanges {
	if len(out) > worktree.MaxStdout {
		return nil
	}
	remaining := make(map[string]bool, len(ids))
	for _, id := range ids {
		remaining[id] = true
	}
	results := make(map[string]surface.FileChanges, len(ids))
	budget := landedFilesBudget
	var id string
	for out != "" {
		token, rest, ok := strings.Cut(out, "\x00")
		if !ok {
			return nil
		}
		out = rest
		if remaining[token] {
			id = token
			results[id] = surface.FileChanges{}
			delete(remaining, id)
			continue
		}
		status := strings.TrimPrefix(token, "\n")
		switch status {
		case "A", "M", "D":
		case "T":
			status = "M" // same display convention as episode files
		default:
			return nil
		}
		path, rest, ok := strings.Cut(out, "\x00")
		if id == "" || !ok || !filepath.IsLocal(path) {
			return nil
		}
		out = rest
		files := results[id]
		cost := 6*len(path) + 32
		if len(files.Changes) < surface.MaxChanges && cost <= budget && len(path) <= MaxPreviewPath && utf8.ValidString(path) {
			files.Changes = append(files.Changes, surface.Change{Path: strings.Clone(path), Status: status})
			budget -= cost
		} else {
			files.Truncated++
		}
		results[id] = files
	}
	if len(remaining) != 0 {
		return nil
	}
	return results
}
