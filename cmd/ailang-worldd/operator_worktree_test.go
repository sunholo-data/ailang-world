package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

const worktreeExample = "ailang-worldd service runner reconcile-worktree oa-0123456789abcdef0123456789abcdef"

type worktreeStartRig struct {
	operatorStartRig
	*operatorServiceFixture
	db  *sql.DB
	src operatorStartSource
}

func newWorktreeStartRig(t *testing.T, claude bool, ttl ...int64) *worktreeStartRig {
	t.Helper()
	return newWorktreeStartRigWithContext(t, claude, nil, ttl...)
}

func newWorktreeStartRigWithContext(t *testing.T, claude bool, ctx context.Context, ttl ...int64) *worktreeStartRig {
	t.Helper()
	return newWorktreeStartRigWithAgent(t, claude, ctx, "", ttl...)
}

func newWorktreeStartRigWithAgent(t *testing.T, claude bool, ctx context.Context, body string, ttl ...int64) *worktreeStartRig {
	t.Helper()
	return newWorktreeStartRigWithTask(t, claude, ctx, body, operatorStartTask, ttl...)
}

func newWorktreeStartRigWithTask(t *testing.T, claude bool, ctx context.Context, body, task string, ttl ...int64) *worktreeStartRig {
	t.Helper()
	return newWorktreeStartRigWithRemote(t, claude, ctx, body, task, "", ttl...)
}

func newWorktreeStartRigWithRemote(t *testing.T, claude bool, ctx context.Context, body, task, remote string, ttl ...int64) *worktreeStartRig {
	t.Helper()
	s := operatorService(t)
	s.mustInstall(t)
	s.writes = nil
	s.env.terminal.now = func() int64 { return 1012 }
	f := newPresetFixture(t, ctx)
	if remote != "" {
		presetGit(t, f.r.f, "remote", "set-url", "origin", remote)
	}
	// The runner checks the agent and session directory up to its login home.
	bin, err := os.MkdirTemp(s.config.Home, "agent-")
	if err != nil {
		t.Fatal(err)
	}
	startReadFakeAgentIn(t, f, bin)
	if body != "" {
		if err := os.WriteFile(filepath.Join(bin, "codex"), []byte(body), 0700); err != nil {
			t.Fatal(err)
		}
	}
	f.r.sessions = filepath.Join(s.config.Home, "sessions")
	r := operatorStartRig{f: f, bin: bin, creds: s.config.Credentials, hosts: t.TempDir()}
	r.f.r.f.db = s.config.World
	r.action = approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: fixtureActionID, Kind: "start", Episode: "ep",
		Store: r.creds, Options: approval.OperatorOptions{Preset: "p", Headless: true}, Task: task, Review: "fixture"}
	r.action.Argv = approval.OperatorArgv(r.action.Kind, r.action.Episode, r.action.Options)
	if claude {
		if err := os.Rename(filepath.Join(r.bin, "codex"), filepath.Join(r.bin, "claude")); err != nil {
			t.Fatal(err)
		}
	}
	r.editPreset(t, func(p *launchPreset) {
		p.Daemon, _ = json.Marshal(presetDaemon{Addr: defaultEpisodeAddr, DB: s.config.World, CredentialsDB: r.creds,
			WorkspaceRoot: r.f.r.f.ws, SessionsDir: r.f.r.sessions})
		if claude {
			p.Agent = "claude"
		}
		if remote != "" {
			p.Committer = "Fixture Operator <fixture@example.invalid>"
			p.Checks = p.Checks[:1]
		}
		if len(ttl) != 0 {
			p.TTLS = ttl[0]
		}
	})
	h := matchingHealth(t, r.f.r.f, r.f.r.profiles, r.f.profile)
	h.DBPath = s.config.World
	if body != "" {
		h.Build.ExecutableSHA256 = s.config.Executables[s.env.paths.Binary]
		s.health = h
	}
	stubDaemon(t, h, nil)
	if body != "" {
		episodeSelfBuild = func() daemon.BuildIdentity { return h.Build }
	}
	r.at = store.OperatorMintAttempt{Catalogue: 1}
	src := operatorStartSource{r.f.dir, r.f.r.profiles, r.f.r.f.ws, r.hosts, s.env.terminal.now}
	if ctx == nil {
		ctx = presetTestContext(t)
	}
	in, facts, err := operatorStartFacts(ctx, r.action, r.at, src.presets, src.profiles, src.workspace, src.hosts, src.clock, resolvePresetStart)
	if err != nil {
		t.Fatal(err)
	}
	r.in, r.action.Resolution = in, facts
	held := claimFixtureEpisodeAction(t, r.creds, r.action.Episode, func(a *approval.OperatorAction) {
		*a = r.action
	})
	if err := held.Close(); err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite", runnerDB(r.creds))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	err = runnerWrite(r.creds, func(c *store.Credentials) error {
		ctx := boundedTestContext(t)
		row, err := c.OperatorAction(ctx, r.action.ID, src.clock())
		if err != nil {
			return err
		}
		h, err := c.OperatorHeartbeat(ctx)
		r.at.CAS = store.OperatorCAS{ID: row.ID, Digest: [32]byte(row.Digest), Generation: row.Generation}
		r.at.Runner = h.Owner
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return &worktreeStartRig{r, s, db, src}
}

func (r *worktreeStartRig) start(t *testing.T, step func(string) error) error {
	t.Helper()
	ctx, cancel := context.WithTimeout(presetTestContext(t), 5*time.Second)
	defer cancel()
	held := r.borrow(t)
	defer held.Close()
	return reserveOperatorStartWorktree(ctx, runnerDB(r.creds), held, r.action, r.at, r.src, step)
}

func (r *worktreeStartRig) borrow(t *testing.T) *store.OperatorEpisodeLock {
	t.Helper()
	held, err := store.OpenOperatorEpisodeLock(r.creds, r.action.Episode)
	if err != nil {
		t.Fatal(err)
	}
	return held
}

func (r *worktreeStartRig) exec(t *testing.T, q string, args ...any) {
	t.Helper()
	ctx := r.f.r.f.gitContext
	if ctx == nil {
		ctx = boundedTestContext(t)
	}
	if _, err := r.db.ExecContext(ctx, q, args...); err != nil {
		t.Fatal(err)
	}
}

func (r *worktreeStartRig) path() string { return filepath.Join(r.f.r.f.ws, r.action.Episode) }

func (r *worktreeStartRig) rows(t *testing.T, wantState, wantPhase, wantReceipt string, wantVersion int) {
	t.Helper()
	ctx := boundedTestContext(t)
	var state, owner string
	var generation, terminal, sessions, starts int
	q := `SELECT state, coalesce(owner,''), generation, coalesce(terminal_at,0),
		(SELECT count(*) FROM session_credentials), (SELECT count(*) FROM episode_starts) FROM operator_actions`
	err := r.db.QueryRowContext(ctx, q).Scan(&state, &owner, &generation, &terminal, &sessions, &starts)
	wantOwner, wantGeneration, wantTerminal := r.at.Runner.Nonce, r.at.CAS.Generation, 0
	switch wantState {
	case "running":
		wantGeneration++
	case "uncertain":
		wantGeneration += 2
	case "failed":
		wantOwner, wantGeneration, wantTerminal = "", wantGeneration+int64(wantVersion), 1012
	}
	if err != nil || state != wantState || owner != wantOwner || int64(generation) != wantGeneration ||
		terminal != wantTerminal || starts != 0 || sessions != 0 {
		t.Fatalf("action=%s owner=%q generation=%d terminal=%d starts=%d sessions=%d: %v", state, owner, generation, terminal, starts, sessions, err)
	}
	var phase, receipt, identity, stepOwner, actionID string
	var version, created int
	err = r.db.QueryRowContext(ctx, `SELECT phase, coalesce(receipt,''), version, identity, owner, action_id, created_at
		FROM operator_action_steps WHERE kind='worktree'`).Scan(&phase, &receipt, &version, &identity, &stepOwner, &actionID, &created)
	if wantPhase == "" {
		if !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("unexpected step: %v", err)
		}
	} else {
		var target operatorWorktreeTarget
		path, pathErr := filepath.EvalSymlinks(filepath.Dir(r.path()))
		repo, repoErr := filepath.EvalSymlinks(r.f.r.f.repo)
		if err != nil || phase != wantPhase || receipt != wantReceipt || version != wantVersion ||
			actionID != r.action.ID || stepOwner != r.at.Runner.Nonce || created != 1012 ||
			json.Unmarshal([]byte(identity), &target) != nil || pathErr != nil || repoErr != nil ||
			target != (operatorWorktreeTarget{r.action.ID, repo, filepath.Join(repo, ".git"), filepath.Join(path, r.action.Episode), "world/ep", r.f.base}) {
			t.Fatalf("step=%q receipt=%q version=%d target=%+v owner=%q action=%q created=%d: %v", phase, receipt, version, target, stepOwner, actionID, created, err)
		}
	}
	r.f.r.assertNoSessionFiles(t)
	if _, err := os.Lstat(filepath.Join(r.bin, "ran")); !errors.Is(err, os.ErrNotExist) ||
		len(r.writes) != 0 || r.ctl.boots != 0 || r.ctl.stops != 0 {
		t.Fatalf("unexpected launch, terminal service effect or file write: %v", err)
	}
}

func (r *worktreeStartRig) state(t *testing.T, wantState, wantPhase, wantReceipt string, wantVersion int, tree bool) {
	t.Helper()
	r.rows(t, wantState, wantPhase, wantReceipt, wantVersion)
	if tree {
		out, err := worktree.Output(boundedTestContext(t), r.f.r.f.git, []string{"-C", r.path(), "rev-parse", "HEAD", "--symbolic-full-name", "HEAD"})
		if err != nil || strings.TrimSpace(out) != r.f.base+"\nrefs/heads/world/ep" ||
			r.f.r.f.worktrees(t) != 2 || presetGit(t, r.f.r.f, "rev-parse", "refs/heads/world/ep") != r.f.base {
			t.Fatalf("registered worktree and pinned ref: %q %v", out, err)
		}
	} else {
		if _, err := os.Lstat(r.path()); !errors.Is(err, os.ErrNotExist) || r.f.r.f.worktrees(t) != 1 ||
			strings.Contains(presetGit(t, r.f.r.f, "for-each-ref", "--format=%(refname)", "refs/heads"), "refs/heads/world/ep") {
			t.Fatalf("unexpected worktree or branch: %v", err)
		}
	}
}

func (r *worktreeStartRig) run(t *testing.T, cmd string) (int, string) {
	t.Helper()
	argv := strings.Fields(cmd)
	if len(argv) != 5 || argv[0] != "ailang-worldd" || argv[1] != "service" {
		t.Fatalf("invalid worktree remedy %q", cmd)
	}
	return r.operatorServiceFixture.run(argv[2:]...)
}

func TestOperatorStartOnce(t *testing.T) {
	t.Run("worktree", testOperatorStartOnceWorktree)
	t.Run("launch", testOperatorStartOnceLaunch)
}

func testOperatorStartOnceWorktree(t *testing.T) {
	r := newWorktreeStartRig(t, false)
	if err := r.start(t, func(step string) error {
		if step == "before add" {
			r.state(t, "running", "reserved", "", 1, false)
			contender, err := os.Open(operatorEpisodeLockPath(r.creds, r.action.Episode))
			if err != nil {
				t.Fatal(err)
			}
			defer contender.Close()
			if err := syscall.Flock(int(contender.Fd()), syscall.LOCK_SH|syscall.LOCK_NB); !errors.Is(err, syscall.EWOULDBLOCK) {
				t.Fatalf("episode lock not held across reservation: %v", err)
			}
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	r.state(t, "running", "observed", "created", 2, true)
	before := observeProvision(t, r.f.r, r.creds)
	err := r.start(t, func(string) error { t.Fatal("replay crossed an effect boundary"); return nil })
	if err == nil || !strings.Contains(err.Error(), "branch already exists") || !reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
		t.Fatalf("replay changed state or wrong refusal: %v", err)
	}
}

func TestOperatorStartCollisions(t *testing.T) {
	for _, tc := range []struct{ name, why string }{
		{"branch", "branch already exists"}, {"registration", "registration already exists"},
		{"registration_foreign", "registration already exists"}, {"ref_ancestor", "registration already exists"}, {"ref_descendant", "registration already exists"},
		{"git_initial_read", "fixture worktree list refusal"},
		{"directory", "not a git worktree"}, {"symlink", "not a directory"},
		{"token", "ep.session already exists"}, {"task", "ep.task.md already exists"}, {"mcp", "ep.mcp.json already exists"},
		{"episode", "operator action conflicts"}, {"credential", "operator action conflicts"},
		{"stale", "mint_start"}, {"catalogue", "operator action conflicts"},
		{"runner", "operator action conflicts"}, {"revoked", "revoked"}, {"expires_at_reservation", "operator action expired"},
		{"heartbeat_stale", "operator action conflicts"}, {"heartbeat_future", "operator action conflicts"},
		{"heartbeat_protocol", "operator action conflicts"}, {"heartbeat_readiness", "operator action conflicts"},
		{"held", "episode lock pathname changed"}, {"missing_lock", "episode lock pathname changed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := newWorktreeStartRig(t, tc.name == "mcp")
			ctx := boundedTestContext(t)
			lockPath := operatorEpisodeLockPath(r.creds, r.action.Episode)
			held := r.borrow(t)
			defer held.Close()
			switch tc.name {
			case "branch":
				presetGit(t, r.f.r.f, "branch", "world/ep", r.f.base)
			case "registration_foreign":
				other := filepath.Join(t.TempDir(), "other")
				presetGit(t, r.f.r.f, "worktree", "add", "-b", "world/ep", other, r.f.base)
				presetGit(t, r.f.r.f, "update-ref", "-d", "refs/heads/world/ep")
			case "ref_ancestor", "ref_descendant":
				branch := map[string]string{"ref_ancestor": "world", "ref_descendant": "world/ep/sub"}[tc.name]
				presetGit(t, r.f.r.f, "branch", branch, r.f.base)
			case "git_initial_read":
				git := filepath.Join(r.bin, "git")
				if err := os.Remove(git); err != nil {
					t.Fatal(err)
				}
				script := "#!/bin/sh\nif [ \"$3 $4 $5 $6\" = 'worktree list --porcelain -z' ]; then echo 'fixture worktree list refusal' >&2; exit 1; fi\nexec \"" + r.f.r.f.git + "\" \"$@\"\n"
				if err := os.WriteFile(git, []byte(script), 0700); err != nil {
					t.Fatal(err)
				}
			case "registration":
				presetGit(t, r.f.r.f, "worktree", "add", "--detach", r.path(), r.f.base)
				if err := os.RemoveAll(r.path()); err != nil {
					t.Fatal(err)
				}
			case "directory":
				if err := os.Mkdir(r.path(), 0700); err != nil {
					t.Fatal(err)
				}
			case "symlink":
				if err := os.Symlink(t.TempDir(), r.path()); err != nil {
					t.Fatal(err)
				}
			case "token", "task", "mcp":
				if err := os.MkdirAll(r.f.r.sessions, 0700); err != nil {
					t.Fatal(err)
				}
				suffix := map[string]string{"token": ".session", "task": ".task.md", "mcp": ".mcp.json"}[tc.name]
				seedStartFixtureFile(t, filepath.Join(r.f.r.sessions, "ep"+suffix), []byte(r.action.Task))
			case "credential", "episode":
				err := runnerWrite(r.creds, func(c *store.Credentials) error {
					return c.MintSessionWithStart(ctx, store.SessionRow{CredentialID: serviceHash([]byte(r.action.Task)), EpisodeID: "ep",
						CreatedAt: 1010, ExpiresAt: 1600, GrantsJSON: "[]"}, &r.in.Start)
				})
				if err != nil {
					t.Fatal(err)
				}
				if tc.name == "episode" {
					r.exec(t, "DELETE FROM session_credentials")
				} else {
					r.exec(t, "DELETE FROM episode_starts")
				}
			case "stale":
				r.editPreset(t, func(p *launchPreset) { p.TTLS-- })
			case "catalogue":
				r.exec(t, "UPDATE operator_runner_state SET catalogue_generation=catalogue_generation+1")
			case "runner":
				r.at.Runner.Generation++
			case "expires_at_reservation":
				calls := 0
				r.src.clock = func() int64 {
					calls++
					if calls >= 4 {
						return 1600
					}
					return 1012
				}
			case "heartbeat_stale", "heartbeat_future", "heartbeat_protocol", "heartbeat_readiness":
				q := map[string]string{"heartbeat_stale": "seen_at=1002", "heartbeat_future": "seen_at=1013",
					"heartbeat_protocol": "protocol=2", "heartbeat_readiness": "readiness='finish'"}
				r.exec(t, "UPDATE operator_runner_state SET "+q[tc.name])
			case "revoked":
				if err := runnerWrite(r.creds, func(c *store.Credentials) error {
					return c.RevokeApprover(ctx, "fixture", 1011)
				}); err != nil {
					t.Fatal(err)
				}
			case "held":
				r.editPreset(t, func(p *launchPreset) { p.TTLS-- })
				// Keep the original inode linked while a different description holds the canonical lock.
				if err := os.Rename(lockPath, lockPath+".old"); err != nil {
					t.Fatal(err)
				}
				other := r.borrow(t)
				defer other.Close()
			case "missing_lock":
				if err := os.Remove(lockPath); err != nil {
					t.Fatal(err)
				}
			}
			before := observeProvision(t, r.f.r, r.creds)
			err := reserveOperatorStartWorktree(presetTestContext(t), runnerDB(r.creds), held, r.action, r.at, r.src, func(string) error { return nil })
			if err == nil || !strings.Contains(err.Error(), tc.why) || !reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) ||
				len(r.writes) != 0 || r.ctl.boots != 0 || r.ctl.stops != 0 {
				t.Fatalf("collision %s changed state or wrong refusal: %v", tc.name, err)
			}
			if tc.name == "held" || tc.name == "missing_lock" {
				if !errors.Is(err, store.ErrOperatorConflict) {
					t.Fatalf("expected episode lock conflict: %v", err)
				}
				r.rows(t, "claimed", "", "", 0)
			}
		})
	}
}

func TestOperatorWorktreeWindow(t *testing.T) {
	md, err := os.ReadFile("../../docs/QUICKSTART.md")
	if err != nil || !strings.Contains(string(md), "\n"+worktreeExample+"\n") {
		t.Fatalf("missing executable QUICKSTART example: %v", err)
	}
	for _, window := range []string{"before reservation", "before add", "before readback", "before observe", "git refused", "success then absent", "path only", "ref changed", "head changed", "readback error", "observe error", "add error after checkout", "redirected checkout", "registration locked", "ref read skew"} {
		t.Run(window, func(t *testing.T) {
			r := newWorktreeStartRig(t, false)
			initial := observeProvision(t, r.f.r, r.creds)
			var atFailure provisionSnapshot
			injected := errors.New("injected interruption")
			gitPath := filepath.Join(r.bin, "git")
			err := r.start(t, func(step string) error {
				switch {
				case window == "add error after checkout" && step == "before add":
					hook := filepath.Join(r.f.r.f.repo, ".git", "hooks", "post-checkout")
					if err := os.WriteFile(hook, []byte("#!/bin/sh\nexit 1\n"), 0700); err != nil {
						t.Fatal(err)
					}
				case window == "add error after checkout" && step == "before readback":
					r.state(t, "running", "reserved", "", 1, true)
					atFailure = observeProvision(t, r.f.r, r.creds)
				case window == "success then absent" && step == "before readback":
					r.state(t, "running", "reserved", "", 1, true)
					presetGit(t, r.f.r.f, "worktree", "remove", r.path())
					presetGit(t, r.f.r.f, "branch", "-D", "world/ep")
					atFailure = observeProvision(t, r.f.r, r.creds)
				case window == "git refused" && step == "before add":
					dir := filepath.Join(r.f.r.f.repo, ".git", "refs", "heads", "world")
					if err := os.Mkdir(dir, 0700); err != nil {
						t.Fatal(err)
					}
					seedStartFixtureFile(t, filepath.Join(dir, "ep.lock"), []byte(r.action.Task))
				case (window == "path only" || window == "ref changed" || window == "head changed" || (window == "redirected checkout" || window == "registration locked")) && step == "before readback":
					switch window {
					case "path only":
						presetGit(t, r.f.r.f, "worktree", "remove", "--force", r.path())
						presetGit(t, r.f.r.f, "branch", "-D", "world/ep")
						if err := os.Mkdir(r.path(), 0700); err != nil {
							t.Fatal(err)
						}
					case "ref changed":
						presetGit(t, r.f.r.f, "update-ref", "refs/heads/world/ep", "feature")
					case "registration locked":
						presetGit(t, r.f.r.f, "worktree", "lock", r.path())
					case "redirected checkout":
						decoy := filepath.Join(t.TempDir(), "decoy")
						presetGit(t, r.f.r.f, "worktree", "add", "--detach", decoy, r.f.base)
						seedStartFixtureFile(t, filepath.Join(r.path(), ".git"), readStartFixtureFile(t, filepath.Join(decoy, ".git")))
					case "head changed":
						git := r.f.r.f
						git.repo = r.path()
						presetGit(t, git, "checkout", "--detach", "feature")
					}
					atFailure = observeProvision(t, r.f.r, r.creds)
				case window == "ref read skew" && step == "before readback":
					if err := os.Remove(gitPath); err != nil {
						t.Fatal(err)
					}
					script := "#!/bin/sh\nreal=" + r.f.r.f.git + "\n" +
						"if [ \"$3\" = for-each-ref ]; then\n" +
						"  \"$real\" -C \"$2\" update-ref refs/heads/world/ep feature || exit 1\n" +
						"  \"$real\" \"$@\" || exit 1\n" +
						"  \"$real\" -C \"$2\" update-ref refs/heads/world/ep " + r.f.base + "\n" +
						"else exec \"$real\" \"$@\"; fi\n"
					if err := os.WriteFile(gitPath, []byte(script), 0700); err != nil {
						t.Fatal(err)
					}
					atFailure = observeProvision(t, r.f.r, r.creds)
				case window == "ref read skew" && step == "before observe":
					if err := os.Remove(gitPath); err != nil {
						t.Fatal(err)
					}
					if err := os.Symlink(r.f.r.f.git, gitPath); err != nil {
						t.Fatal(err)
					}
				case window == "readback error" && step == "before readback":
					if err := os.Rename(gitPath, gitPath+".hidden"); err != nil {
						t.Fatal(err)
					}
					atFailure = observeProvision(t, r.f.r, r.creds)
				case window == "readback error" && step == "before observe":
					if err := os.Rename(gitPath+".hidden", gitPath); err != nil {
						t.Fatal(err)
					}
				case window == "observe error" && step == "before observe":
					r.exec(t, `CREATE TRIGGER fail_observe BEFORE UPDATE OF phase ON operator_action_steps
						BEGIN SELECT RAISE(ABORT,'fixture observation refused'); END`)
				case step == window:
					return injected
				}
				return nil
			})
			switch window {
			case "before reservation":
				r.state(t, "claimed", "", "", 0, false)
				if !errors.Is(err, injected) || !reflect.DeepEqual(initial, observeProvision(t, r.f.r, r.creds)) {
					t.Fatalf("failure before reservation changed state or class: %v", err)
				}
				if code, out := r.run(t, worktreeExample); code != exitUsage || !strings.Contains(out, "operator action conflicts") ||
					!reflect.DeepEqual(initial, observeProvision(t, r.f.r, r.creds)) {
					t.Fatalf("command accepted an action without reservation: %d %s", code, out)
				}
				return
			case "git refused":
				r.state(t, "failed", "observed", "stopped: absent", 2, false)
				if err == nil || err.Error() != "worktree add failed; worktree and branch absent" {
					t.Fatal(err)
				}
				return
			case "observe error":
				r.state(t, "running", "reserved", "", 1, true)
				if err == nil || !strings.Contains(err.Error(), "fixture observation refused") {
					t.Fatal(err)
				}
				r.exec(t, "DROP TRIGGER fail_observe")
			case "success then absent":
				r.state(t, "uncertain", "reserved", "", 2, false)
				got := observeProvision(t, r.f.r, r.creds)
				got.credentials = atFailure.credentials
				if !reflect.DeepEqual(atFailure, got) || err == nil || err.Error() != "worktree outcome uncertain; "+worktreeExample {
					t.Fatalf("successful add with absent readback lost uncertainty or evidence: %v", err)
				}
			case "path only", "ref changed", "head changed", "readback error", "add error after checkout", "redirected checkout", "registration locked", "ref read skew":
				r.rows(t, "uncertain", "reserved", "", 2)
				got := observeProvision(t, r.f.r, r.creds)
				got.credentials = atFailure.credentials
				if !reflect.DeepEqual(atFailure, got) || err == nil || !strings.Contains(err.Error(), "worktree outcome uncertain") {
					t.Fatalf("partial result was changed or misclassified: %v", err)
				}
			default:
				r.state(t, "running", "reserved", "", 1, window != "before add")
				if !errors.Is(err, injected) {
					t.Fatalf("lost failure class: %v", err)
				}
			}
			_, remedy, ok := strings.Cut(err.Error(), "; ")
			if !ok || remedy != worktreeExample {
				t.Fatalf("incorrect remedy: %v", err)
			}
			before := observeProvision(t, r.f.r, r.creds)
			for _, answer := range []string{"no tty", "n\n", "y\n"} {
				r.env.terminal = testMintEnv(answer, 1012)
				if answer == "no tty" {
					r.env.terminal = noTerminalEnv(1012)
				}
				code, msg := r.run(t, remedy)
				if answer != "y\n" {
					_, retry, ok := strings.Cut(msg, "Retry at a terminal: ")
					if code != exitUsage || !ok || strings.TrimSuffix(retry, "\n") != remedy ||
						!reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
						t.Fatalf("unattended recovery modified state: %d %s", code, msg)
					}
					remedy = strings.TrimSuffix(retry, "\n")
					continue
				}
				if code != exitOK || !strings.Contains(msg, "Worktree reservation disposed for "+r.action.ID) {
					t.Fatalf("%d %s", code, msg)
				}
			}
			receipt, version := "stopped: retained", 2
			switch window {
			case "success then absent":
				receipt, version = "stopped: absent", 3
			case "before add":
				receipt = "stopped: absent"
			case "path only", "ref changed", "head changed", "readback error", "add error after checkout", "redirected checkout", "registration locked", "ref read skew":
				version = 3
			}
			r.rows(t, "failed", "observed", receipt, version)
			got := observeProvision(t, r.f.r, r.creds)
			got.credentials = before.credentials
			if !reflect.DeepEqual(before, got) {
				t.Fatal("recovery modified files, worktree registration or branch")
			}
			before = observeProvision(t, r.f.r, r.creds)
			if code, msg := r.run(t, remedy); code != exitUsage || !strings.Contains(msg, "operator action conflicts") ||
				!reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
				t.Fatalf("second disposition changed the state: %d %s", code, msg)
			}
		})
	}
}

func TestOperatorWorktreeDispositionCAS(t *testing.T) {
	for _, change := range []string{"action_owner", "action_generation", "action_digest", "action_kind", "action_state", "step_identity", "step_phase", "step_owner", "step_version", "later_effect", "duplicate_step", "git_changed", "lock_held"} {
		t.Run(change, func(t *testing.T) {
			r := newWorktreeStartRig(t, false)
			if err := r.start(t, func(step string) error {
				if step == "before add" {
					return errors.New("crash")
				}
				return nil
			}); err == nil {
				t.Fatal("expected interrupted reservation")
			}
			before := observeProvision(t, r.f.r, r.creds)
			why := "operator action conflicts"
			if change == "lock_held" {
				fd, err := openRunnerLock(operatorEpisodeLockPath(r.creds, "ep"), 0)
				if err != nil {
					t.Fatal(err)
				}
				defer fd.Close()
				defer time.AfterFunc(3*time.Second, func() { fd.Close() }).Stop()
				why = "episode recovery lock held"
			}
			if change == "git_changed" {
				why = "worktree evidence changed after review"
			}
			r.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
				switch change {
				case "action_owner", "action_generation", "action_digest", "action_kind", "action_state":
					edit := map[string]string{"action_owner": "owner='replacement'", "action_generation": "generation=generation+1",
						"action_digest": "digest=zeroblob(32)", "action_kind": "kind='finish'", "action_state": "state='claimed'"}
					// Change only the action after the step update, inside the disposition transaction.
					r.exec(t, `CREATE TRIGGER stale_action AFTER UPDATE OF phase ON operator_action_steps
						BEGIN UPDATE operator_actions SET `+edit[change]+`; END`)
				case "step_phase":
					r.exec(t, "UPDATE operator_action_steps SET phase='observed',receipt=? WHERE kind='worktree'", []byte("created"))
				case "step_identity":
					r.exec(t, "UPDATE operator_action_steps SET identity=identity||' ' WHERE kind='worktree'")
				case "step_owner":
					r.exec(t, "UPDATE operator_action_steps SET owner='replacement' WHERE kind='worktree'")
				case "step_version":
					r.exec(t, "UPDATE operator_action_steps SET version=version+1 WHERE kind='worktree'")
				case "later_effect":
					r.exec(t, `INSERT INTO operator_action_steps(kind,identity,action_id,phase,version,owner,created_at)
							SELECT 'mint', id, id, 'reserved', 1, owner, 1012 FROM operator_actions`)
				case "duplicate_step":
					r.exec(t, `INSERT INTO operator_action_steps SELECT kind, identity||'another', action_id, phase, version, owner, created_at, receipt
						FROM operator_action_steps WHERE kind='worktree'`)
				case "git_changed":
					presetGit(t, r.f.r.f, "branch", "world/ep", r.f.base)
				default:
					t.Fatal("prompt opened despite lock contention")
				}
				before = observeProvision(t, r.f.r, r.creds)
				return testMintEnv("y\n", 1012).openTerminal()
			}
			if code, msg := r.run(t, worktreeExample); code != exitUsage || !strings.Contains(msg, why) ||
				!reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
				t.Fatalf("changed evidence was accepted: %d %s", code, msg)
			}
		})
	}
}
