package daemon

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"unicode/utf8"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

const maxSurfaceTaskBytes = 64 << 10

// readFollowEpisodeStart uses the resolved front state, never the request's
// raw selector. Landed states need the follower's operator-backed finish
// association. An episode displayed as metadata does not establish that association.
func (d *Daemon) readFollowEpisodeStart(ctx context.Context, env surface.FollowIndex, stack surfacepage.FollowStack) *store.EpisodeStart {
	repo, kind, episode, err := surface.ParseQualifiedID(stack.Front.ID)
	if err != nil || repo != stack.Repo || env.Open == nil {
		return nil
	}
	switch {
	case kind == "episode" && stack.Front.Kind == surfacepage.KindOpen:
		if episode != stack.Front.Episode {
			return nil
		}
	case kind == "commit" && stack.Front.Kind == surface.KindStep:
		episode = ""
		for _, alias := range env.Open.Aliases {
			if alias.Target != stack.Front.ID {
				continue
			}
			aliasRepo, aliasKind, aliasEpisode, err := surface.ParseQualifiedID(alias.ID)
			if err != nil || aliasRepo != repo || aliasKind != "episode" || episode != "" {
				return nil
			}
			episode = aliasEpisode
		}
	default:
		return nil
	}
	if d.follow == nil || !surface.EpisodeGrammar.MatchString(episode) {
		return nil
	}
	var configured string
	for _, r := range d.follow.config.Repos {
		if r.ID == repo {
			configured = r.Repo
			break
		}
	}
	if configured == "" {
		return nil
	}
	start, err := d.readEpisodeStart(ctx, episode)
	if err != nil || start == nil || start.Episode != episode ||
		!filepath.IsAbs(start.Repo) || canonicalRepo(start.Repo) != configured {
		return nil
	}
	return start
}

// readSurfaceTask reads the service CLI's retained task copy. Custom preset
// session directories and start.Task.Path are not part of daemon configuration.
func (d *Daemon) readSurfaceTask(start *store.EpisodeStart) (body string, redacted bool) {
	if start == nil || start.Task == nil || d.cfg.DBPath == "" {
		return "", false
	}
	name := start.Episode + ".task.md"
	if !surface.EpisodeGrammar.MatchString(start.Episode) {
		return "", false
	}
	want, err := hashref.Parse(start.Task.SHA256)
	if err != nil {
		return "", false
	}
	dir := filepath.Join(filepath.Dir(d.cfg.DBPath), "sessions")
	root := openSurfaceDir(dir)
	if root == nil {
		return "", false
	}
	defer root.Close()
	data := readSurfaceFile(root, name, maxSurfaceTaskBytes)
	if data == nil || !utf8.Valid(data) || hashref.SumSHA256(data) != want {
		return "", false
	}
	body = surface.SanitizeTaskBody(string(data))
	return body, body != string(data)
}

func openSurfaceDir(dir string) *os.Root {
	info, err := os.Lstat(dir)
	if err != nil || !info.IsDir() {
		return nil
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil
	}
	opened, err := root.Stat(".")
	if err != nil || !os.SameFile(info, opened) {
		_ = root.Close()
		return nil
	}
	return root
}

func readSurfaceFile(root *os.Root, name string, limit int64) []byte {
	fileInfo, err := root.Lstat(name)
	if err != nil || !fileInfo.Mode().IsRegular() {
		return nil
	}
	f, err := root.OpenFile(name, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() || !os.SameFile(fileInfo, info) || info.Size() > limit {
		return nil
	}
	data, err := io.ReadAll(io.LimitReader(f, limit+1))
	if err != nil || int64(len(data)) > limit {
		return nil
	}
	return data
}
