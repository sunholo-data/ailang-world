package daemon

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

func TestSurfaceTaskExactSelectionAndVisibility(t *testing.T) {
	cfg, _ := followFixture(t, 2)
	d := newFollowDaemon(t, cfg)
	snapshot := followSnapshotFor(t, "task", []string{"repo0", "repo1"}, nil)
	injectFollow(d, snapshot)
	now := time.Now()
	const task = "Retained prompt\n\n  • Run local checks\n  • Ask operator@example.invalid\n  • Review /Users/operator/private/config\n  • Remove ghp_abcdefghijklmnopqrstuv\n"
	const want = "Retained prompt\n\n  • Run local checks\n  • Ask [email]\n  • Review [path]\n  • Remove [redacted]\n"
	path := writeSurfaceTask(t, d, "wa1", []byte(task))
	start := overlayStart("wa1", d.follow.config.Repos[0].Repo, now)
	start.Task.SHA256 = hashref.SumSHA256([]byte(task)).String()
	mintSurfaceTaskStart(t, d, start)
	wrongRepo := overlayStart("wb1", d.follow.config.Repos[0].Repo, now)
	wrongRepo.Task.SHA256 = start.Task.SHA256
	mintSurfaceTaskStart(t, d, wrongRepo)
	writeSurfaceTask(t, d, "wb1", []byte(task))

	stack := func(repo, lane, id string) surfacepage.StackFront {
		t.Helper()
		q := url.Values{"brane": {repo + ":brane:" + lane}, "at": {repo + ":" + id}}
		w := followRequest(d, "/surface/stack.json?"+q.Encode())
		var v surfacepage.FollowStack
		if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &v) != nil {
			t.Fatalf("stack: %d %s", w.Code, w.Body)
		}
		if strings.Contains(w.Body.String(), start.Task.Path) || strings.Contains(w.Body.String(), path) {
			t.Fatal("task path surfaced")
		}
		return v.Front
	}
	if f := stack("repo0", "shared", "episode:wa1"); f.TaskBody != want || !f.TaskRedacted {
		t.Fatalf("open body = %q (redacted %v)", f.TaskBody, f.TaskRedacted)
	}
	if f := stack("repo1", "shared", "episode:wb1"); f.TaskBody != "" {
		t.Fatalf("wrong repo body = %q", f.TaskBody)
	}
	for _, path := range []string{"/surface/model.json", "/surface/index.json", "/surface/follow-status.json", "/surface/static"} {
		if w := followRequest(d, path); strings.Contains(w.Body.String(), "Retained prompt") || strings.Contains(w.Body.String(), "taskBody") {
			t.Fatalf("%s exposes retained task body", path)
		}
	}
	var landed surface.State
	for _, s := range snapshot.Envelope.Repos[0].Index.States {
		if s.Kind == surface.KindStep && s.Brane != "" {
			landed = s
			break
		}
	}
	if landed.ID == "" {
		t.Fatal("fixture has no landed brane step")
	}
	if f := stack("repo0", landed.Brane, "commit:"+landed.ID); f.TaskBody != "" {
		t.Fatalf("unassociated step body = %q", f.TaskBody)
	}
	snapshot.Envelope.Open.Episodes = snapshot.Envelope.Open.Episodes[1:]
	snapshot.Envelope.Open.Aliases = append(snapshot.Envelope.Open.Aliases, surface.EpisodeAlias{ID: "repo0:episode:wa1", Target: "repo0:commit:" + landed.ID})
	if err := surface.ValidateFollow(snapshot.Envelope); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"commit:" + landed.ID, "episode:wa1"} {
		if f := stack("repo0", landed.Brane, id); f.ID != "repo0:commit:"+landed.ID || f.TaskBody != want || !f.TaskRedacted {
			t.Fatalf("%s body = %q, state %s", id, f.TaskBody, f.ID)
		}
	}
	if f := stack("repo0", "shared", "commit:"+landed.Merge); f.TaskBody != "" {
		t.Fatalf("shared merge body = %q", f.TaskBody)
	}
	if f := stack("repo1", landed.Brane, "commit:"+landed.ID); f.TaskBody != "" {
		t.Fatalf("colliding commit body = %q", f.TaskBody)
	}
}

func TestSurfaceTaskRefusesUnverifiedCopy(t *testing.T) {
	good := []byte("Exact prompt\n\n  Keep spacing.\n")
	for _, tc := range []struct {
		name string
		data []byte
		edit func(*testing.T, string, *store.EpisodeStart)
	}{
		{name: "valid", data: good},
		{name: "digest mismatch", data: good, edit: func(t *testing.T, _ string, s *store.EpisodeStart) {
			s.Task.SHA256 = hashref.SumSHA256([]byte("different")).String()
		}},
		{name: "missing", data: good, edit: func(t *testing.T, p string, s *store.EpisodeStart) {
			if err := os.Rename(p, p+".source"); err != nil {
				t.Fatal(err)
			}
			s.Task.Path = p + ".source"
		}},
		{name: "symlink", data: good, edit: func(t *testing.T, p string, _ *store.EpisodeStart) {
			if err := os.Rename(p, p+".source"); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(filepath.Base(p)+".source", p); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "old record without task", data: good, edit: func(_ *testing.T, _ string, s *store.EpisodeStart) { s.Task = nil }},
		{name: "old task without digest", data: good, edit: func(_ *testing.T, _ string, s *store.EpisodeStart) { s.Task.SHA256 = "" }},
		{name: "invalid UTF-8", data: []byte{0xff, 0xfe}},
		{name: "oversize", data: make([]byte, maxSurfaceTaskBytes+1)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := &Daemon{cfg: Config{DBPath: filepath.Join(t.TempDir(), "world.db")}}
			s := &store.EpisodeStart{Episode: "e1", Task: &store.StartTask{SHA256: hashref.SumSHA256(tc.data).String()}}
			path := writeSurfaceTask(t, d, s.Episode, tc.data)
			if tc.edit != nil {
				tc.edit(t, path, s)
			}
			want := ""
			if tc.name == "valid" {
				want = string(good)
			}
			if got, redacted := d.readSurfaceTask(s); got != want || redacted {
				t.Fatalf("body length = %d, want %d (redacted %v)", len(got), len(want), redacted)
			}
		})
	}
}

func writeSurfaceTask(t *testing.T, d *Daemon, episode string, bytes []byte) string {
	t.Helper()
	dir := filepath.Join(filepath.Dir(d.cfg.DBPath), "sessions")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, episode+".task.md")
	if err := os.WriteFile(path, bytes, 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func mintSurfaceTaskStart(t *testing.T, d *Daemon, start *store.EpisodeStart) {
	t.Helper()
	credentials, err := store.OpenCredentials(d.credentialsPath)
	if err != nil {
		t.Fatal(err)
	}
	defer credentials.Close()
	now := time.Now().Unix()
	row := store.SessionRow{CredentialID: "cred-" + start.Episode, EpisodeID: start.Episode, GrantsJSON: "[]", CreatedAt: now, ExpiresAt: now + 3600}
	if err := credentials.MintSessionWithStart(boundedTestContext(t), row, start); err != nil {
		t.Fatal(err)
	}
}
