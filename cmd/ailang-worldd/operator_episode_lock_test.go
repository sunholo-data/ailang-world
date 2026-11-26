package main

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
)

const fixtureActionID = "oa-0123456789abcdef0123456789abcdef"

func claimFixtureEpisode(t *testing.T, creds, episode string) *os.File {
	return claimFixtureEpisodeAction(t, creds, episode, func(*approval.OperatorAction) {})
}

func claimFixtureEpisodeAction(t *testing.T, creds, episode string, prepare func(*approval.OperatorAction)) *os.File {
	t.Helper()
	canonical, err := store.CanonicalPath(creds)
	if err != nil {
		t.Fatal(err)
	}
	held, err := runnerKernelLock(operatorEpisodeLockPath(canonical, episode))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { held.Close() })
	if _, err := store.UpgradeCredentials(creds); err != nil {
		t.Fatal(err)
	}
	db, err := store.OpenCredentials(creds)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx := boundedTestContext(t)
	key := approvaltest.NewSigner(51)
	if err := db.EnrolApprover(ctx, store.ApproverRow{Name: "fixture", Format: "ssh", PublicKey: key.PublicKey().Blob, EnrolledAt: 900}); err != nil {
		t.Fatal(err)
	}
	a := approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: fixtureActionID, Kind: "start", Episode: episode,
		Store: canonical, Options: approval.OperatorOptions{Preset: "p", Headless: true}, Task: "fixture", Review: "fixture",
		Resolution: []approval.OperatorResolved{{Name: "agent", Value: "codex"}}}
	prepare(&a)
	a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
	req := store.OperatorRequest{Action: a, OperationID: strings.Repeat("1", 32), Challenge: sha256.Sum256([]byte(t.Name())), CreatedAt: 1000, ExpiresAt: 1600}
	row, err := db.RequestOperatorAction(ctx, req, 1000)
	if err != nil {
		t.Fatal(err)
	}
	cas := store.OperatorCAS{ID: row.ID, Digest: [32]byte(row.Digest), Generation: row.Generation}
	in := authority.OperatorJudgeInput{Action: a, Digest: cas.Digest, BytesDigest: sha256.Sum256(row.Action), Challenge: req.Challenge,
		CreatedAt: req.CreatedAt, ExpiresAt: req.ExpiresAt, Now: 1010}
	message, err := approval.OperatorActionMessage(cas.Digest, req.Challenge, req.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	artifact := key.Sign(approval.OperatorActionNamespace, message)
	proof, err := authority.JudgeOperatorAction(ctx, db, in, artifact)
	if err != nil {
		t.Fatal(err)
	}
	row, err = db.ConfirmOperatorAction(ctx, cas, artifact, proof, 1010)
	if err != nil {
		t.Fatal(err)
	}
	h := store.OperatorHeartbeat{Owner: store.OperatorOwner{Nonce: "test-runner"}, Protocol: store.OperatorRunnerProtocol, Build: "fixture", Readiness: "start", SeenAt: 1010}
	if h.Owner.Generation, err = db.PutOperatorHeartbeat(ctx, h); err != nil {
		t.Fatal(err)
	}
	cas.Generation = row.Generation
	if err := db.ClaimOperatorAction(ctx, cas, h.Owner, 1011); err != nil {
		t.Fatal(err)
	}
	// Bound the caller's contention test after the fixture is ready.
	fuse := time.AfterFunc(2*time.Second, func() { held.Close() })
	t.Cleanup(func() { fuse.Stop() })
	return held
}

func expectEpisodeRefusal(t *testing.T, code int, output string, before, creds string) {
	t.Helper()
	if code != exitUsage || !strings.Contains(output, "active action "+fixtureActionID+";") ||
		!strings.Contains(output, "service runner status") || serviceFileHash(creds) != before {
		t.Fatalf("expected exact owner refusal with unchanged credentials: %d %s", code, output)
	}
}

func TestOperatorTerminalExclusion(t *testing.T) {
	t.Run("nonwaiting", testOperatorEpisodeLockTransition)
	t.Run("plain", func(t *testing.T) {
		r := newEpisodeRun(t)
		health := matchingHealth(t, r.f, r.profiles)
		stubDaemon(t, health, nil)
		claimFixtureEpisode(t, health.CredentialsDB, "ep1")
		before, trees := serviceFileHash(health.CredentialsDB), r.f.worktrees(t)
		var out bytes.Buffer
		opened := false
		code := runEpisodeStart(r.args("ep1"), &out, &out, spyEnv(&opened))
		expectEpisodeRefusal(t, code, out.String(), before, health.CredentialsDB)
		r.assertNoSessionFiles(t)
		if opened || r.f.worktrees(t) != trees {
			t.Fatal("terminal opened or worktree changed")
		}
	})
	t.Run("preset", func(t *testing.T) {
		f := newPresetFixture(t)
		creds, _ := credentialsFor(f.r.f.db, "")
		claimFixtureEpisode(t, creds, "p-2610281352")
		before, trees := serviceFileHash(creds), f.r.f.worktrees(t)
		var out bytes.Buffer
		opened := false
		env := spyEnv(&opened)
		env.now = presetClock
		code := runEpisodeStart(f.args("--task", f.task), &out, &out, env)
		expectEpisodeRefusal(t, code, out.String(), before, creds)
		f.r.assertNoSessionFiles(t)
		if opened || f.r.f.worktrees(t) != trees {
			t.Fatal("terminal opened or worktree changed")
		}
	})
	t.Run("resume", func(t *testing.T) {
		f, ep, _ := resumeFixture(t)
		creds, _ := credentialsFor(f.r.f.db, "")
		held := claimFixtureEpisode(t, creds, ep)
		other, err := terminalEpisodeLock(creds, "other", nil)
		if err != nil || other == nil {
			t.Fatalf("durable owner of %s blocked another episode: %v", ep, err)
		}
		other.Close()
		if err := held.Close(); err != nil {
			t.Fatal(err)
		}
		before, trees := serviceFileHash(creds), f.r.f.worktrees(t)
		token := filepath.Join(f.r.sessions, ep+".session")
		tokenHash := serviceFileHash(token)
		now := presetClock() + 43201
		var out bytes.Buffer
		opened := false
		env := spyEnv(&opened)
		env.now = func() int64 { return now }
		code := runEpisodeStart(f.args("--resume", ep), &out, &out, env)
		expectEpisodeRefusal(t, code, out.String(), before, creds)
		archived, err := filepath.Glob(token + ".expired-*")
		if err != nil || len(archived) != 0 || opened || serviceFileHash(token) != tokenHash || f.r.f.worktrees(t) != trees {
			t.Fatal("expired approval or missing kernel holder permitted a resume effect")
		}
	})
	t.Run("launch", func(t *testing.T) {
		f := newLaunchFixture(t, "claude")
		claimFixtureEpisode(t, f.creds, "ep")
		before := serviceFileHash(f.creds)
		token := serviceFileHash(filepath.Join(f.sessions, "ep.session"))
		launchClient = func(string) *client { t.Fatal("launch reached report dispatch"); return nil }
		code, _, output := f.run()
		expectEpisodeRefusal(t, code, output, before, f.creds)
		if files, err := os.ReadDir(f.facts); err != nil || len(files) != 0 || token != serviceFileHash(filepath.Join(f.sessions, "ep.session")) {
			t.Fatal("launch modified session files or spawned")
		}
		if _, err := os.Lstat(filepath.Join(f.sessions, "ep.launch.lock")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("launch refusal left a PID file: %v", err)
		}
	})
	t.Run("finish", func(t *testing.T) {
		t.Setenv("HOME", t.TempDir())
		f := newFinishFixture(t)
		creds := f.health.CredentialsDB
		claimFixtureEpisode(t, creds, "ep1")
		before, status := serviceFileHash(creds), f.gitOutput("-C", f.wt, "status", "--porcelain")
		opened := false
		env := sessionEnv{openTerminal: func() (io.ReadWriteCloser, error) {
			if len(f.requests) != 0 {
				t.Fatal("health preceded F1")
			}
			opened = true
			return &fakeTerm{r: strings.NewReader("n\n")}, nil
		}}
		code, _, output := f.run(env)
		expectEpisodeRefusal(t, code, output, before, creds)
		if !opened || !reflect.DeepEqual(f.requests, []string{"health"}) || f.checks != 0 || f.outcomes != 0 ||
			strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != f.before || f.gitOutput("-C", f.wt, "status", "--porcelain") != status {
			t.Fatalf("F1 or finish effect order changed: %v", f.requests)
		}
	})
}

func testOperatorEpisodeLockTransition(t *testing.T) {
	creds := filepath.Join(t.TempDir(), "credentials.db")
	canonical, err := store.CanonicalPath(creds)
	if err != nil {
		t.Fatal(err)
	}
	initial, err := terminalEpisodeLock(creds, "ep", nil)
	if err != nil || initial == nil {
		t.Fatalf("pre-management acquisition: %v %v", initial, err)
	}
	defer initial.Close()
	process, err := runnerKernelLock(canonical + ".operator-runner.lock")
	if err != nil {
		t.Fatal(err)
	}
	defer process.Close()
	// An independently opened fd represents the future paired process's
	// acquisition. Release on timeout so a blocking mutation cannot hang.
	results := make(chan error, 1)
	go func() {
		contender, err := runnerKernelLock(operatorEpisodeLockPath(canonical, "ep"))
		if contender != nil {
			contender.Close()
		}
		results <- err
	}()
	select {
	case err := <-results:
		if !errors.Is(err, syscall.EWOULDBLOCK) {
			t.Fatalf("same episode did not refuse: %v", err)
		}
	case <-time.After(time.Second):
		initial.Close()
		t.Fatal("episode acquisition waited on the first terminal holder")
	}
	if lock, err := terminalEpisodeLock(creds, "ep", nil); lock != nil || !errors.Is(err, store.ErrOperatorConflict) ||
		!strings.Contains(err.Error(), "episode lock held") {
		t.Fatalf("unidentified holder: %v %v", lock, err)
	}
	if _, err := os.Stat(creds); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("refusal created a database: %v", err)
	}
	for _, tc := range []struct{ creds, episode string }{{creds, "other"}, {filepath.Join(filepath.Dir(creds), "other.db"), "ep"}} {
		fd, err := terminalEpisodeLock(tc.creds, tc.episode, nil)
		if err != nil || fd == nil {
			t.Fatalf("independent store/episode denied: %v %v", fd, err)
		}
		fd.Close()
	}
	if _, err := terminalEpisodeLock(creds, "other", initial); err == nil {
		t.Fatal("transferred lock to another episode")
	}
	alias := filepath.Join(t.TempDir(), "alias")
	if err := os.Symlink(filepath.Dir(creds), alias); err != nil {
		t.Fatal(err)
	}
	if fd, err := terminalEpisodeLock(filepath.Join(alias, filepath.Base(creds)), "ep", initial); fd != nil || err != nil {
		t.Fatalf("canonical transfer: %v %v", fd, err)
	}
	if err := initial.Close(); err != nil {
		t.Fatal(err)
	}
	fd, err := runnerKernelLock(operatorEpisodeLockPath(canonical, "ep"))
	if err != nil {
		t.Fatal(err)
	}
	fd.Close()
}

func TestOperatorLaunchServedCredentials(t *testing.T) {
	f := newLaunchFixture(t, "claude")
	fetchDaemonHealth = func(string) (daemon.HealthResponse, error) {
		return daemon.HealthResponse{DBPath: f.db, CredentialsDB: filepath.Join(f.root, "other.db")}, nil
	}
	launchClient = func(string) *client { t.Fatal("mismatch reached report dispatch"); return nil }
	before := serviceFileHash(f.creds)
	code, _, output := f.run()
	if code != exitUsage || !strings.Contains(output, "credentials store differs") || serviceFileHash(f.creds) != before {
		t.Fatalf("credentials pair mismatch: %d %s", code, output)
	}
}

func TestOperatorEpisodeLockRefusals(t *testing.T) {
	for _, arm := range []string{"unsafe lock", "ownership read", "invalid action id", "unclaimed action"} {
		t.Run(arm, func(t *testing.T) {
			creds := filepath.Join(t.TempDir(), "credentials.db")
			canonical, err := store.CanonicalPath(creds)
			if err != nil {
				t.Fatal(err)
			}
			lockPath := operatorEpisodeLockPath(canonical, "ep")
			want, change := "could not read action ownership", ""
			switch arm {
			case "unsafe lock":
				want = "unsafe runner lock file"
				if err := os.WriteFile(lockPath, []byte("keep"), 0600); err != nil {
					t.Fatal(err)
				}
				if err := os.Chmod(lockPath, 0644); err != nil {
					t.Fatal(err)
				}
			case "ownership read":
				change = "DROP TABLE operator_actions"
			case "invalid action id":
				change = "PRAGMA foreign_keys=OFF; UPDATE operator_actions SET id='invalid'"
			case "unclaimed action":
				change = "UPDATE operator_actions SET state='confirmed', owner=NULL"
			}
			if change != "" {
				held := claimFixtureEpisode(t, creds, "ep")
				held.Close()
				credSQL(t, creds, change)
			}
			before, lockBefore := serviceFileHash(creds), serviceFileHash(lockPath)
			lock, err := terminalEpisodeLock(creds, "ep", nil)
			if lock != nil {
				lock.Close()
			}
			if arm == "unclaimed action" {
				if err != nil || lock == nil {
					t.Fatalf("confirmed action without an owner excludes terminal: %v", err)
				}
			} else if err == nil || lock != nil || !strings.Contains(err.Error(), want) {
				t.Fatalf("refusal: %v %v", lock, err)
			}
			if serviceFileHash(creds) != before || serviceFileHash(lockPath) != lockBefore {
				t.Fatal("store or lock bytes changed")
			}
			if arm != "unsafe lock" {
				fd, err := runnerKernelLock(lockPath)
				if err != nil {
					t.Fatalf("terminal failed to close its lock: %v", err)
				}
				fd.Close()
			}
		})
	}
}

func TestOperatorEpisodeLockExclusive(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lock")
	holder, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer holder.Close()
	if err := syscall.Flock(int(holder.Fd()), syscall.LOCK_SH|syscall.LOCK_NB); err != nil {
		t.Fatal(err)
	}
	result := make(chan error, 1)
	go func() {
		lock, err := runnerKernelLock(path)
		if lock != nil {
			lock.Close()
		}
		result <- err
	}()
	select {
	case err := <-result:
		if !errors.Is(err, syscall.EWOULDBLOCK) {
			t.Fatalf("runner lock did not exclude a shared holder: %v", err)
		}
	case <-time.After(time.Second):
		holder.Close()
		t.Fatal("runner lock waited on a shared holder")
	}
}

func TestOperatorLaunchClaimDuringPreflight(t *testing.T) {
	f := newLaunchFixture(t, "claude")
	var before string
	launchClient = func(string) *client {
		return memoryDaemon(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/v1/episodes/ep/report" {
				claimFixtureEpisode(t, f.creds, "ep")
				before = serviceFileHash(f.creds)
			}
			f.serve(w, r)
		})
	}
	code, _, output := f.run()
	expectEpisodeRefusal(t, code, output, before, f.creds)
	if facts, err := os.ReadDir(f.facts); err != nil || len(facts) != 0 {
		t.Fatal("claim during preflight allowed a child")
	}
	if _, err := os.Lstat(filepath.Join(f.sessions, "ep.launch.lock")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("claim during preflight left a PID file: %v", err)
	}
}
