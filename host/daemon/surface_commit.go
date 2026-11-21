package daemon

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacefollow"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

const (
	surfaceCommitMessageBytes = 16 << 10
	surfaceCommitCacheEntries = 256
)

// readSurfaceCommits uses only the resolved published selection and a finish
// whose episode and operator provenance were already checked.
func (d *Daemon) readSurfaceCommits(ctx context.Context, stack surfacepage.FollowStack, finish *EpisodeFinish) []surfacepage.NodeCommit {
	repo, kind, selected, err := surface.ParseQualifiedID(stack.Front.ID)
	if err != nil || repo != stack.Repo || d.follow == nil {
		return nil
	}
	var id string
	var kinds []string
	switch {
	case kind == "commit" && stack.Front.Kind == surface.KindShared:
		id, kinds = selected, []string{"merge"}
	case kind == "commit" && stack.Front.Kind == surface.KindStep:
		id, kinds = selected, []string{"branch-tip"}
		if stack.Front.Run != nil && qualifiesSurfaceFinish(finish, stack.Front.Run.Episode, selected) {
			kinds = append(kinds, "episode-finish")
		}
	case kind == "episode" && stack.Front.Kind == surfacepage.KindOpen &&
		qualifiesSurfaceFinish(finish, selected, ""):
		id, kinds = finish.Commit, []string{"episode-finish"}
	default:
		return nil
	}
	for _, source := range d.follow.config.Repos {
		if source.ID == repo {
			msg, ok := d.readSurfaceCommitMessage(ctx, source, id)
			if !ok {
				return nil
			}
			msg.Kinds = kinds
			return []surfacepage.NodeCommit{msg}
		}
	}
	return nil
}

func qualifiesSurfaceFinish(finish *EpisodeFinish, episode, selected string) bool {
	if finish == nil || episode == "" || finish.Episode != episode || (selected != "" && selected != finish.Commit) {
		return false
	}
	switch finish.Outcome {
	case "done", "push-failed", "pr-failed", "ambiguous":
		return true
	}
	return false
}

func (d *Daemon) readSurfaceCommitMessage(ctx context.Context, source surfacefollow.Repo, id string) (surfacepage.NodeCommit, bool) {
	qualified, err := surface.QualifyID(source.ID, "commit", id)
	if err != nil {
		return surfacepage.NodeCommit{}, false
	}
	key := source.Repo + "\x00" + qualified
	d.follow.mu.Lock()
	cached, ok := d.follow.commitMessages[key]
	d.follow.mu.Unlock()
	if ok {
		return cached, true
	}
	if d.follow.commitReads.Add(1) > 2 {
		d.follow.commitReads.Add(-1)
		return surfacepage.NodeCommit{}, false
	}
	defer d.follow.commitReads.Add(-1)
	ctx, cancel := context.WithTimeout(ctx, surfacefollow.CallTimeout)
	defer cancel()
	raw, err := (surfacefollow.ReadGit{}).Git(ctx, []string{"-C", source.Repo, "--no-optional-locks",
		"--no-lazy-fetch", "--no-replace-objects", "cat-file", "commit", id})
	if err != nil {
		return surfacepage.NodeCommit{}, false
	}
	msg, ok := parseSurfaceCommitMessage(raw)
	if !ok {
		return surfacepage.NodeCommit{}, false
	}
	msg.ID = qualified
	d.follow.mu.Lock()
	if d.follow.commitMessages == nil || len(d.follow.commitMessages) >= surfaceCommitCacheEntries {
		d.follow.commitMessages = make(map[string]surfacepage.NodeCommit)
	}
	d.follow.commitMessages[key] = msg
	d.follow.mu.Unlock()
	return msg, true
}

func parseSurfaceCommitMessage(raw string) (surfacepage.NodeCommit, bool) {
	_, message, ok := strings.Cut(raw, "\n\n")
	if !ok || len(message) > surfaceCommitMessageBytes || !utf8.ValidString(message) {
		return surfacepage.NodeCommit{}, false
	}
	subject, body, _ := strings.Cut(message, "\n")
	// One blank line separates the subject from the body in a Git message.
	body = strings.TrimPrefix(body, "\n")
	cleanSubject, cleanBody := surface.SanitizeTaskBody(subject), surface.SanitizeTaskBody(body)
	return surfacepage.NodeCommit{Subject: strings.Clone(cleanSubject), Body: strings.Clone(cleanBody),
		Redacted: cleanSubject != subject || cleanBody != body}, true
}
