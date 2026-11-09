package daemon

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacefollow"
)

// Row 220 M3 end to end: real credentials and world stores, the daemon's
// episode source, and the served follow routes. Fake identities only.

const leakedPreset = `{"schema":"world/preset/v1","secret":"preset-bytes-do-not-serve"}`

func overlayStart(ep, repo string, now time.Time) *store.EpisodeStart {
	return &store.EpisodeStart{Schema: "world/episode-start/v1", Episode: ep, StartedAt: now.Add(-time.Hour).Format(time.RFC3339),
		Preset: store.StartPreset{Name: "attended", Path: "/home/operator/.ailang-world/presets/attended.json", SHA256: "sha256:p", Bytes: leakedPreset},
		Repo:   repo, Remote: store.StartRemote{Name: "origin", URL: "https://operator:remote-token@example.invalid/world.git"},
		Base: store.StartBase{Branch: "dev", SHA: strings.Repeat("3", 40)}, Branch: "trial/" + ep, Agent: "claude", Model: "opus",
		Task:   &store.StartTask{Path: "/home/operator/tasks/" + ep + ".md", SHA256: "sha256:t", Title: "Ask operator@example.invalid about " + ep},
		Checks: []store.StartCheck{{Name: "unit", Argv: []string{"go", "test", "./secret-check/..."}, Cwd: "/home/operator/dev/world"}},
		TTLS:   3600}
}

func putFinish(t *testing.T, d *Daemon, record EpisodeFinish, provenance string) {
	t.Helper()
	payload, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	ctx := boundedTestContext(t)
	obj := store.Object{Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte(EpisodeFinishV1)),
		SemanticID: EpisodeFinishV1, Provenance: provenance, Payload: payload}
	if err := d.store.PutObject(ctx, obj); err != nil {
		t.Fatal(err)
	}
	if err := d.store.SetRegistryHead(ctx, finishHead(record.Episode), obj.Hash); err != nil {
		t.Fatal(err)
	}
}

func TestSurfaceFollowOpenEpisodes(t *testing.T) {
	cfg, root := followFixture(t, 1)
	d := newFollowDaemon(t, cfg)
	repo := d.follow.config.Repos[0].Repo
	// A start may record a symlinked spelling of the configured checkout.
	link := filepath.Join(root, "repo-link")
	if err := os.Symlink(repo, link); err != nil {
		t.Fatal(err)
	}
	elsewhere := filepath.Join(root, "unconfigured")
	if err := os.Mkdir(elsewhere, 0700); err != nil {
		t.Fatal(err)
	}
	now := followClock{}.Now()
	w, err := store.OpenCredentials(d.credentialsPath)
	if err != nil {
		t.Fatal(err)
	}
	ctx := boundedTestContext(t)
	mint := func(id, ep string, created, expires int64, start *store.EpisodeStart) {
		t.Helper()
		row := store.SessionRow{CredentialID: id, EpisodeID: ep, GrantsJSON: "[]", CreatedAt: created, ExpiresAt: expires}
		if start != nil {
			err = w.MintSessionWithStart(ctx, row, start)
		} else {
			err = w.MintSession(ctx, row)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	// More rows than the page's 500-row bound: the join must read them all.
	for i := 0; i < store.MaxListSessions+5; i++ {
		mint(fmt.Sprintf("old-%04d", i), fmt.Sprintf("old-%04d", i), now.Unix()-90000+int64(i), now.Unix()-86400, nil)
	}
	live := now.Unix() + 3600
	mint("cred-w1-a", "w1", now.Unix()-60, live, overlayStart("w1", link, now))
	mint("cred-w1-b", "w1", now.Unix()-30, live, nil)
	mint("cred-w2", "w2", now.Unix()-20, live, overlayStart("w2", elsewhere, now))
	mint("cred-legacy", "legacy", now.Unix()-10, live, nil)
	mint("cred-w4", "w4", now.Unix()-5, live, overlayStart("w4", repo, now))
	mint("cred-w5", "w5", now.Unix()-4, live, overlayStart("w5", repo, now))
	mint("cred-late", "late", now.Unix()-1, live, overlayStart("late", repo, now))
	// Finish records, then revoke: the operator's record makes w4 PR open;
	// an object under the same head from any other writer is not a finish.
	putFinish(t, d, EpisodeFinish{Episode: "w4", Branch: "trial/w4", Commit: strings.Repeat("9", 40), Remote: "origin",
		PR: FinishPR{URL: "https://example.invalid/pull/12", Number: 12}, Outcome: "done"}, finishProvenance)
	putFinish(t, d, EpisodeFinish{Episode: "w5", Branch: "trial/w5", Commit: strings.Repeat("8", 40),
		PR: FinishPR{Number: 13}, Outcome: "done"}, "session:forged")
	for _, id := range []string{"cred-w4", "cred-w5"} {
		if err := w.RevokeSession(ctx, id); err != nil {
			t.Fatal(err)
		}
	}
	_ = w.Close()

	startFollowTest(t, d, &followGit{})
	waitFollow(t, func() bool {
		s := d.follow.publisher.Snapshot()
		return len(s.Envelope.Repos) == 1 && s.Envelope.Open != nil
	})
	body := followRequest(d, "/surface/index.json")
	if body.Code != 200 {
		t.Fatalf("%d %s", body.Code, body.Body)
	}
	ix, err := surface.DecodeFollow(body.Body.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	o := ix.Open
	got := map[string]surface.OpenEpisode{}
	for _, e := range o.Episodes {
		got[e.Episode] = e
	}
	if o.Status != surface.OpenComplete || len(o.Episodes) != 4 || o.Unconfigured != 1 {
		t.Fatalf("%s %q episodes=%d unconfigured=%d", o.Status, o.Reason, len(o.Episodes), o.Unconfigured)
	}
	if e := got["w1"]; e.ID != "repo0:episode:w1" || e.State != surface.OpenRunning || e.Anchor != "repo0:commit:"+strings.Repeat("3", 40) ||
		e.Agent != "claude" || e.Model != "opus" || e.Branch != "trial/w1" || e.Task != "Ask [email] about w1" {
		t.Fatalf("w1: %+v", e)
	}
	if e := got["w4"]; e.State != surface.OpenPR || e.PR != 12 || e.SessionLive {
		t.Fatalf("w4: %+v", e)
	}
	if e := got["w5"]; e.State != surface.OpenEnded || e.PR != 0 || e.EndedAt != now.Format(time.RFC3339) {
		t.Fatalf("forged finish trusted: %+v", e)
	}
	if e := got["late"]; e.State != surface.OpenRunning {
		t.Fatalf("row past the 500-row page not joined: %+v", e)
	}
	want := `[{legacy start record unavailable} {w2 repo not configured}]`
	if fmt.Sprint(o.Running) != want {
		t.Fatalf("running %v", o.Running)
	}
	status := followRequest(d, "/surface/follow-status.json")
	var doc surfacefollow.StatusDocument
	if json.Unmarshal(status.Body.Bytes(), &doc) != nil || doc.Open == nil || doc.Open.Status != surface.OpenComplete ||
		doc.Open.Episodes != 4 || doc.Open.Running != 2 || doc.Generation != ix.Generation {
		t.Fatalf("status %s", status.Body)
	}
	for _, out := range []string{body.Body.String(), status.Body.String()} {
		for _, secret := range []string{"preset-bytes-do-not-serve", "presets/attended.json", "remote-token", "example.invalid",
			"tasks/w1.md", "secret-check", "operator@", root, repo, link, elsewhere, "cred-w1-a", "/home/operator"} {
			if strings.Contains(out, secret) {
				t.Fatalf("leaked %q: %s", secret, out)
			}
		}
	}
}

func TestFollowEpisodesSourceFailsClosed(t *testing.T) {
	src := newFollowEpisodes(filepath.Join(t.TempDir(), "missing", store.CredentialsFileName), nil)
	if _, err := src.ReadEpisodes(boundedTestContext(t), time.Now(), nil); err == nil {
		t.Fatal("a missing credentials store read as zero running")
	}
}
