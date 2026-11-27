package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
)

func operatorActivationFixture(t *testing.T) (*worktreeStartRig, runnerConfig, store.OperatorHeartbeat, *os.File) {
	t.Helper()
	r := newWorktreeStartRig(t, false)
	c := r.config
	agent := operatorFileFact(r.action, "launch_binary")
	c.PresetNames = []string{r.action.Options.Preset}
	c.Executables[agent] = serviceFileHash(agent)
	p := r.env.runnerPaths()
	config, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	if err := r.env.atomic(p.Record, config, 0600); err != nil {
		t.Fatal(err)
	}
	if err := r.env.atomic(p.Plist, []byte(runnerPlist(p, c, config)), 0600); err != nil {
		t.Fatal(err)
	}
	r.env.controller = r.homeJob
	r.env.runnerProcesses = func() ([]worldProcess, error) {
		return []worldProcess{{os.Getpid(), "ailang-worldd"}}, nil
	}
	record, err := readServiceRecord(r.env.paths)
	if err != nil {
		t.Fatal(err)
	}
	record.Workspace = r.src.workspace
	body, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	if err := r.env.atomic(r.env.paths.Record, body, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(r.f.r.sessions, 0700); err != nil {
		t.Fatal(err)
	}
	r.exec(t, "UPDATE operator_actions SET state='confirmed', owner=NULL")
	build := c.Executables[r.env.paths.Binary]
	catalogue, err := json.Marshal(c.PresetNames)
	if err != nil {
		t.Fatal(err)
	}
	r.exec(t, "UPDATE operator_runner_state SET build=?, catalogue=?", build, catalogue)
	h, err := runnerHeartbeat(r.creds, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	held, err := runnerKernelLock(r.creds + ".operator-runner.lock")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { held.Close() })
	return r, c, h, held
}

func activationEvidence(t *testing.T, r *worktreeStartRig) [5]any {
	t.Helper()
	git := filepath.Join(r.f.r.f.repo, ".git")
	return [5]any{serviceFileHash(r.creds), provisionFiles(t, filepath.Join(git, "refs")),
		provisionFiles(t, filepath.Join(git, "worktrees")), provisionFiles(t, r.src.workspace), provisionFiles(t, r.f.r.sessions)}
}

func TestOperatorFirstManagedEpisodeProcess(t *testing.T) {
	r, c, h, runner := operatorActivationFixture(t)
	ctx := boundedTestContext(t)
	prepare := func(acquire func(string, string) (*store.OperatorEpisodeLock, error), stops <-chan os.Signal) error {
		cursor := ""
		return r.env.consumeOperatorPreparation(ctx, c, h, r.src, stops, runner, &cursor, acquire)
	}
	snapshot := func() any {
		return activationEvidence(t, r)
	}
	row := candidateRow(t, r.creds, r.src.clock())
	legacy := r.action
	legacy.Resolution = nil
	for _, fact := range r.action.Resolution {
		if fact.Name != "launch_binary_sha256" {
			legacy.Resolution = append(legacy.Resolution, fact)
		}
	}
	wire, err := approval.CanonicalOperatorAction(legacy)
	if err != nil {
		t.Fatal(err)
	}
	digest, err := approval.OperatorSignatureDigest(legacy)
	if err != nil {
		t.Fatal(err)
	}
	msg, err := approval.OperatorActionMessage(digest, [32]byte(row.Challenge), row.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	r.exec(t, "UPDATE operator_actions SET action=?, digest=?, artifact=?", wire, digest[:], approvaltest.NewSigner(51).Sign(approval.OperatorActionNamespace, msg))
	before := snapshot()
	var binding *approval.BindingError
	err = prepare(store.OpenOperatorEpisodeLock, nil)
	if !errors.As(err, &binding) || binding.Field != "launch_binary_sha256" || !reflect.DeepEqual(before, snapshot()) {
		t.Fatalf("old signed review: %v", err)
	}
	r.exec(t, "UPDATE operator_actions SET action=?, digest=?, artifact=?", row.Action, row.Digest, row.Artifact)

	if err := runnerWrite(r.creds, func(db *store.Credentials) error {
		return db.ClaimOperatorAction(ctx, operatorPreparationCAS(row), h.Owner, r.src.clock())
	}); err != nil {
		t.Fatal(err)
	}
	beforeRelease := activationEvidence(t, r)
	// The live loop only caches the lease generation, never the catalogue.
	loopHeartbeat := h
	loopHeartbeat.CatalogueGeneration = 0
	cursor := ""
	r.exec(t, "UPDATE operator_runner_state SET catalogue_generation=catalogue_generation+1")
	beforeDrift := snapshot()
	err = r.env.recoverOperatorStart(ctx, c, loopHeartbeat, r.src, &cursor, nil)
	if !errors.As(err, &binding) || binding.Field != "mint_catalogue" || !reflect.DeepEqual(beforeDrift, snapshot()) {
		t.Fatalf("untouched recovery crossed a changed catalogue: %v", err)
	}
	r.exec(t, "UPDATE operator_runner_state SET catalogue_generation=catalogue_generation-1")
	cursor = ""
	foreignHeartbeat := loopHeartbeat
	foreignHeartbeat.Owner.Generation++
	beforeForeign := snapshot()
	err = r.env.recoverOperatorStart(ctx, c, foreignHeartbeat, r.src, &cursor, nil)
	if !errors.As(err, &binding) || binding.Field != "mint_catalogue" || !reflect.DeepEqual(beforeForeign, snapshot()) {
		t.Fatalf("foreign recovery lease released the untouched claim: %v", err)
	}
	r.exec(t, "UPDATE operator_runner_state SET seen_at=?", r.src.clock()-10)
	cursor = ""
	beforeStale := snapshot()
	err = r.env.recoverOperatorStart(ctx, c, loopHeartbeat, r.src, &cursor, nil)
	if !errors.As(err, &binding) || binding.Field != "mint_catalogue" || !reflect.DeepEqual(beforeStale, snapshot()) {
		t.Fatalf("stale recovery lease released the untouched claim: %v", err)
	}
	r.exec(t, "UPDATE operator_runner_state SET seen_at=?", r.src.clock())
	cursor = ""
	if err := r.env.recoverOperatorStart(ctx, c, loopHeartbeat, r.src, &cursor, nil); err != nil {
		t.Fatalf("untouched claim with current signed catalogue: %v", err)
	}
	released := candidateRow(t, r.creds, r.src.clock())
	afterRelease := activationEvidence(t, r)
	if released.ID != row.ID || released.Owner != "" || released.Generation != row.Generation+2 ||
		!reflect.DeepEqual(beforeRelease[1:], afterRelease[1:]) {
		t.Fatalf("untouched release changed physical evidence or lost the confirmed action: %+v", released)
	}

	binary := operatorFileFact(r.action, "launch_binary")
	if err := os.Chmod(binary, 0600); err != nil {
		t.Fatal(err)
	}
	before = snapshot()
	err = prepare(store.OpenOperatorEpisodeLock, nil)
	if err == nil || !strings.Contains(err.Error(), "(F10)") || !reflect.DeepEqual(before, snapshot()) {
		t.Fatalf("non-executable unchanged bytes admitted preparation: %v", err)
	}
	if err := r.env.checkStartPair(c, r.src.workspace); !errors.Is(err, store.ErrOperatorConflict) {
		t.Fatalf("installed pair accepted a non-executable file: %v", err)
	}
	if err := os.Chmod(binary, 0700); err != nil {
		t.Fatal(err)
	}
	r.exec(t, "UPDATE operator_runner_state SET seen_at=?", r.src.clock()-10)
	before = snapshot()
	err = prepare(store.OpenOperatorEpisodeLock, nil)
	if !errors.Is(err, store.ErrOperatorRunner) || !reflect.DeepEqual(before, snapshot()) {
		t.Fatalf("stale heartbeat admitted preparation: %v", err)
	}
	r.exec(t, "UPDATE operator_runner_state SET seen_at=?", r.src.clock())
	before = snapshot()
	oldOwner := h.Owner
	oldOwner.Generation++
	if err := r.env.prepareOperatorStart(ctx, c, oldOwner, h.Build, r.src, nil, runner, store.OpenOperatorEpisodeLock); !errors.Is(err, store.ErrOperatorRunner) || !reflect.DeepEqual(before, snapshot()) {
		t.Fatalf("foreign runner admitted preparation: %v", err)
	}
	readProcesses := r.env.runnerProcesses
	r.env.runnerProcesses = func() ([]worldProcess, error) { return nil, nil }
	err = prepare(store.OpenOperatorEpisodeLock, nil)
	if err == nil || !strings.Contains(err.Error(), "did not include this process") || !reflect.DeepEqual(before, snapshot()) {
		t.Fatalf("unproved terminal retirement admitted preparation: %v", err)
	}
	r.env.runnerProcesses = readProcesses
	healthy := r.health.Build.ExecutableSHA256
	r.health.Build.ExecutableSHA256 = "sha256:unapproved"
	err = r.env.checkStartPair(c, r.src.workspace)
	if err == nil || !strings.Contains(err.Error(), "daemon build mismatch") || !reflect.DeepEqual(before, snapshot()) {
		t.Fatalf("drifted World build admitted preparation: %v", err)
	}
	r.health.Build.ExecutableSHA256 = healthy

	before = snapshot()
	paused := false
	err = prepare(func(path, ep string) (*store.OperatorEpisodeLock, error) {
		paused = true
		_, release := operatorProcessWitness(t, "episode", path)
		defer release()
		return store.OpenOperatorEpisodeLock(path, ep)
	}, nil)
	if !paused || !errors.Is(err, store.ErrOperatorConflict) || !reflect.DeepEqual(before, snapshot()) {
		t.Fatalf("terminal holder did not defeat runner before CAS: %v", err)
	}
	stopped := make(chan os.Signal, 1)
	stopped <- os.Interrupt
	err = prepare(func(path, ep string) (*store.OperatorEpisodeLock, error) {
		t.Fatal("stopped tick entered candidate selection")
		return nil, store.ErrOperatorConflict
	}, stopped)
	if !errors.Is(err, errOperatorDraining) || !reflect.DeepEqual(before, snapshot()) {
		t.Fatalf("stop admitted a claim: %v", err)
	}
	err = prepare(func(path, ep string) (*store.OperatorEpisodeLock, error) {
		stopped <- os.Interrupt
		return store.OpenOperatorEpisodeLock(path, ep)
	}, stopped)
	if !errors.Is(err, errOperatorDraining) || !reflect.DeepEqual(before, snapshot()) {
		t.Fatalf("stop after candidate selection admitted a claim: %v", err)
	}
	bin, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("WORLD_TEST_OPERATOR_WITNESS", "git-inherited")
	t.Setenv("WORLD_TEST_OPERATOR_WITNESS_PATH", runner.Name())
	t.Setenv("GORACE", strings.TrimSpace(os.Getenv("GORACE")+" atexit_sleep_ms=0"))
	hook := filepath.Join(r.f.r.f.repo, ".git", "hooks", "post-checkout")
	if err := os.WriteFile(hook, []byte("#!/bin/sh\nexec "+shellQuote(bin)+" -test.run '^TestOperatorProcessWitness$'\n"), 0700); err != nil {
		t.Fatal(err)
	}
	lifecycle, err := runnerLifecycleLock(filepath.Dir(r.env.paths.Record))
	if err != nil {
		t.Fatal(err)
	}
	err = prepare(store.OpenOperatorEpisodeLock, nil)
	lifecycle.Close()
	if !errors.Is(err, errOperatorDraining) || !reflect.DeepEqual(before, snapshot()) {
		t.Fatalf("claim crossed lifecycle exclusion: %v", err)
	}
	ticks := make(chan time.Time)
	close(ticks)
	var runLog bytes.Buffer
	if code := r.env.runOperatorPreparationLoop(&runLog, c, h, r.src, nil, ticks, runner); code != exitOK || runLog.Len() != 0 {
		t.Fatalf("production tick failed: %d %s", code, runLog.String())
	}
	closed, err := runnerLifecycleLock(filepath.Dir(r.env.paths.Record))
	if err != nil {
		t.Fatalf("claim retained lifecycle exclusion: %v", err)
	}
	closed.Close()
	if _, err := os.Stat(runner.Name() + ".git-observed"); err != nil {
		t.Fatalf("real Git did not observe the inherited process lock: %v", err)
	}
	s, _, err := readRunnerFiles(ctx, r.creds, row.ID)
	if err != nil || s.Phase != "observed" {
		t.Fatalf("file receipt: %s %v", s.Phase, err)
	}
	ro, err := store.OpenCredentialsReadOnly(runnerDB(r.creds))
	if err != nil {
		t.Fatal(err)
	}
	defer ro.Close()
	prepared, err := ro.ReadOperatorStartRecovery(ctx, row.ID)
	if err != nil || prepared.Window != store.OperatorStartPrepared || prepared.Owner != h.Owner.Nonce {
		t.Fatalf("prepared owner and evidence: %+v %v", prepared, err)
	}
	before = snapshot()
	var out bytes.Buffer
	if code := runSessionRevoke([]string{"--db", c.World, s.Session.CredentialID}, &out, &out); code != exitUsage ||
		!strings.Contains(out.String(), "active action "+row.ID) {
		t.Fatalf("terminal could cross prepared ownership: %d %s", code, &out)
	}
	report := daemon.RunEpisodeGC(ctx, daemon.EpisodeGCOptions{WorkspaceRoot: r.src.workspace, CredentialsPath: r.creds, Now: time.Unix(r.src.clock(), 0)})
	if len(report.Removed) != 0 || len(report.Skipped) != 1 || report.Skipped[0].Reason != "non-terminal operator step" ||
		!reflect.DeepEqual(before, snapshot()) {
		t.Fatalf("GC altered prepared owner or evidence: %+v", report)
	}
	if err := prepare(store.OpenOperatorEpisodeLock, nil); !errors.Is(err, store.ErrOperatorMissing) || !reflect.DeepEqual(before, snapshot()) {
		t.Fatalf("next tick prepared the same row twice: %v", err)
	}
	if _, err := os.Stat(filepath.Join(r.bin, "ran")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("spawn crossed the S6 boundary: %v", err)
	}
}

func TestOperatorCurrentOwnerProcessRecovery(t *testing.T) {
	r, c, h, runner := operatorActivationFixture(t)
	ctx := boundedTestContext(t)
	if err := r.env.prepareOperatorStart(ctx, c, h.Owner, h.Build, r.src, nil, runner, store.OpenOperatorEpisodeLock); err != nil {
		t.Fatal(err)
	}
	_, release := operatorProcessWitness(t, "owner", r.creds)
	observe := func() any {
		return activationEvidence(t, r)
	}
	recover := func() error { cursor := ""; return r.env.recoverOperatorStart(ctx, c, h, r.src, &cursor, nil) }
	for _, seen := range []int64{r.src.clock(), r.src.clock() - 11} {
		h.SeenAt = seen
		r.exec(t, "UPDATE operator_runner_state SET seen_at=?", seen)
		before := observe()
		if err := recover(); !errors.Is(err, store.ErrOperatorConflict) || !strings.Contains(err.Error(), "live episode owner") || !reflect.DeepEqual(before, observe()) {
			t.Fatalf("heartbeat displaced live owner: %v", err)
		}
	}
	release()
	before := observe()
	clock := r.src.clock
	r.src.clock = func() int64 {
		other, err := store.OpenOperatorEpisodeLock(r.creds, r.action.Episode)
		if other != nil {
			other.Close()
		}
		if !errors.Is(err, syscall.EWOULDBLOCK) {
			t.Fatalf("prepared readback lost episode exclusion: %v", err)
		}
		return clock()
	}
	if err := recover(); err != nil || !reflect.DeepEqual(before, observe()) {
		t.Fatalf("quiescent prepared recovery retried an effect: %v", err)
	}
	r.src.clock = clock
	s, _, err := readRunnerFiles(ctx, r.creds, r.action.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, reason := range []string{"worktree", "identity", "task", "metadata", "expiry", "revocation"} {
		t.Run(reason, func(t *testing.T) {
			switch reason {
			case "worktree":
				from, to := r.path(), r.path()+".moved"
				if err := os.Rename(from, to); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					if err := os.Rename(to, from); err != nil {
						t.Error(err)
					}
				})
			case "identity", "task":
				path := s.Mint.TokenPath
				if reason == "task" {
					path = r.in.Start.Task.Path
				}
				token, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					if err := os.WriteFile(path, token, 0600); err != nil {
						t.Error(err)
					}
				})
				if err := os.WriteFile(path, []byte("changed file fixture"), 0600); err != nil {
					t.Fatal(err)
				}
			case "metadata":
				r.exec(t, "UPDATE session_credentials SET project='unreviewed' WHERE credential_id=?", s.Session.CredentialID)
				t.Cleanup(func() {
					r.exec(t, "UPDATE session_credentials SET project=? WHERE credential_id=?", s.Session.Project, s.Session.CredentialID)
				})
			case "expiry":
				r.src.clock = func() int64 { return s.Session.ExpiresAt }
				t.Cleanup(func() { r.src.clock = func() int64 { return 1012 } })
			case "revocation":
				r.exec(t, "DELETE FROM session_credentials WHERE credential_id=?", s.Session.CredentialID)
				t.Cleanup(func() {
					if err := runnerWrite(r.creds, func(db *store.Credentials) error { return db.MintSession(ctx, s.Session) }); err != nil {
						t.Error(err)
					}
				})
			}
			before := observe()
			if err := recover(); !errors.Is(err, store.ErrOperatorConflict) || !reflect.DeepEqual(before, observe()) {
				t.Fatalf("recovery crossed %s evidence: %v", reason, err)
			}
		})
	}
	t.Run("partial_orphan", func(t *testing.T) {
		runner.Close()
		// Restore the pre-observation receipt boundary using the real file
		// reservation identity just produced by the runner.
		r.exec(t, `UPDATE operator_action_steps SET phase='reserved', receipt=NULL, version=version-1
			WHERE action_id=? AND kind='files'`, r.action.ID)
		before := activationEvidence(t, r)
		code, out := r.command("reconcile-start", r.action.ID)
		if code != exitUsage || !strings.Contains(out, store.ErrOperatorConflict.Error()) ||
			!reflect.DeepEqual(before, activationEvidence(t, r)) {
			t.Fatalf("existing token allowed orphan disposition: %d %s", code, out)
		}
		if err := os.Remove(s.Mint.TokenPath); err != nil {
			t.Fatal(err)
		}
		before = activationEvidence(t, r)
		code, out = r.command("reconcile-start", r.action.ID)
		after := activationEvidence(t, r)
		if code != exitOK || !strings.Contains(out, "existing evidence retained") || !reflect.DeepEqual(before[1:], after[1:]) {
			t.Fatalf("partial-file disposition: %d %s", code, out)
		}
		if err := runnerWrite(r.creds, func(db *store.Credentials) error {
			_, found, err := db.ResolveSession(ctx, s.Session.CredentialID)
			if found {
				t.Error("orphan credential was retained")
			}
			return err
		}); err != nil {
			t.Fatal(err)
		}
	})
}

func TestOperatorAttendedStartDisposition(t *testing.T) {
	r, c, h, runner := operatorActivationFixture(t)
	ctx := boundedTestContext(t)
	r.exec(t, `CREATE TRIGGER interrupted_mint BEFORE INSERT ON session_credentials
		BEGIN SELECT RAISE(ABORT,'fixture mint rollback'); END`)
	err := r.env.prepareOperatorStart(ctx, c, h.Owner, h.Build, r.src, nil, runner, store.OpenOperatorEpisodeLock)
	if err == nil || !strings.Contains(err.Error(), "fixture mint rollback") {
		t.Fatalf("mint rollback: %v", err)
	}
	ro, err := store.OpenCredentialsReadOnly(runnerDB(r.creds))
	if err != nil {
		t.Fatal(err)
	}
	defer ro.Close()
	id := r.action.ID
	row, err := ro.ReadOperatorStartRecovery(ctx, id)
	if err != nil || row.Window != store.OperatorStartMintPending {
		t.Fatalf("durable rollback window: %s %v", row.Window, err)
	}
	before := activationEvidence(t, r)
	code, output := r.command("reconcile-start", "untyped")
	if code != exitUsage || !strings.Contains(output, "expected reconcile-start ACTION-ID") ||
		!reflect.DeepEqual(before, activationEvidence(t, r)) {
		t.Fatalf("invalid action identifier reached reconciliation: %d %s", code, output)
	}
	code, output = r.command("reconcile-start", id)
	if code != exitUsage || !strings.Contains(output, "runner lock") || !reflect.DeepEqual(before, activationEvidence(t, r)) {
		t.Fatalf("attended reconcile crossed live process: %d %s", code, output)
	}
	runner.Close()
	// A file finalizer must not mask a missing explicit close.
	defer debug.SetGCPercent(debug.SetGCPercent(-1))
	terminal := r.env.terminal.openTerminal
	r.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
		return &fakeTerm{r: strings.NewReader("n\n")}, nil
	}
	code, output = r.command("reconcile-start", id)
	if code != exitUsage || !reflect.DeepEqual(before, activationEvidence(t, r)) {
		t.Fatalf("attended disposition crossed typed refusal: %d %s", code, output)
	}
	r.env.terminal.openTerminal = terminal
	code, output = r.command("reconcile-start", id)
	after := activationEvidence(t, r)
	if code != exitOK || !strings.Contains(output, "existing evidence retained") {
		t.Fatalf("attended disposition: %d %s", code, output)
	}
	before[0] = after[0]
	if !reflect.DeepEqual(before, after) {
		t.Fatal("disposition modified the retained worktree, branch or private files")
	}
	stopped, err := ro.OperatorAction(ctx, id, 0)
	if err != nil || stopped.State != "failed" || stopped.Owner != "" {
		t.Fatalf("rollback disposition: %s %q %v", stopped.State, stopped.Owner, err)
	}
}
