package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/checkrun"
	"github.com/sunholo-data/ailang-world/host/store"
)

func operatorLocalFinishFixture(t *testing.T, change func(*store.OperatorFinishReview), confirmed ...bool) (*finishFixture, operatorFinishWork, *int64) {
	t.Helper()
	f := newFinishFixture(t)
	g, snap := operatorFinishFixtureSnapshot(t, f)
	return operatorLocalFinishOn(t, f, g, snap, change, confirmed...)
}

func operatorFinishFixtureSnapshot(t *testing.T, f *finishFixture) (operatorFinishGit, store.OperatorFinishTree) {
	t.Helper()
	g := operatorFinishGitFixture(t, f)
	snap, _, err := g.snapshot(activationTestContext(t), f.rep.Start)
	if err != nil {
		t.Fatal(err)
	}
	return g, snap
}

func operatorLocalFinishOn(t *testing.T, f *finishFixture, g operatorFinishGit, snap store.OperatorFinishTree, change func(*store.OperatorFinishReview), confirmed ...bool) (*finishFixture, operatorFinishWork, *int64) {
	t.Helper()
	ctx := activationTestContext(t)
	path, err := store.CanonicalPath(f.health.CredentialsDB)
	if err != nil {
		t.Fatal(err)
	}
	db, err := store.OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	key := approvaltest.NewSigner(51)
	if err := db.EnrolApprover(activationTestContext(t), store.ApproverRow{Name: "fixture", Format: "ssh", PublicKey: key.PublicKey().Blob, EnrolledAt: 900}); err != nil {
		t.Fatal(err)
	}
	m := store.OperatorFinishReview{Store: path, Episode: "ep1", Preset: "fixture", Worktree: g.dir,
		Remote: f.rep.Start.Remote, Checks: []store.StartCheck{{Name: "short", Argv: []string{"echo", "fixture"}, Exec: "echo", Cwd: ".", Sandbox: "sandboxed", TimeoutS: 15}},
		Commit: store.OperatorLocalCommit{Snapshot: snap, Message: "Fixture finish\n\nWorld-Provenance: episode=ep1\n", Committer: "Fixture Operator <fixture@example.invalid>"}}
	if change != nil {
		change(&m)
	}
	a, err := m.Action(1, false)
	if err != nil {
		t.Fatal(err)
	}
	a.ID = fixtureActionID
	req := store.OperatorRequest{Action: a, OperationID: strings.Repeat("1", 32), Challenge: sha256.Sum256([]byte(t.Name())), CreatedAt: 1000, ExpiresAt: 1600}
	row, err := db.RequestOperatorAction(activationTestContext(t), req, 1000)
	if err != nil {
		t.Fatal(err)
	}
	cas := operatorPreparationCAS(row)
	msg, err := approval.OperatorActionMessage(cas.Digest, req.Challenge, req.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	sig := key.Sign(approval.OperatorActionNamespace, msg)
	in := authority.OperatorJudgeInput{Action: a, Digest: cas.Digest, BytesDigest: sha256.Sum256(row.Action),
		Challenge: req.Challenge, CreatedAt: 1000, ExpiresAt: 1600, Now: 1001}
	proof, err := authority.JudgeOperatorAction(ctx, db, in, sig)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.ConfirmOperatorAction(activationTestContext(t), cas, sig, proof, 1001); err != nil {
		t.Fatal(err)
	}
	h := store.OperatorHeartbeat{Owner: store.OperatorOwner{Nonce: "finish-runner"}, Protocol: store.OperatorRunnerProtocol, Build: "fixture", Readiness: "both", SeenAt: 1001}
	h.Owner.Generation, err = db.PutOperatorHeartbeat(activationTestContext(t), h)
	if err != nil {
		t.Fatal(err)
	}
	owner := ""
	if len(confirmed) == 0 {
		if err := db.ClaimOperatorAction(activationTestContext(t), cas, h.Owner, 1001); err != nil {
			t.Fatal(err)
		}
		cas.Generation++
		owner = h.Owner.Nonce
	}
	ro, err := store.OpenCredentialsReadOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ro.Close() })
	now := int64(1001)
	p, err := verifyOperatorFinishPrep(ctx, ro, a, now, cas, owner)
	if err != nil {
		t.Fatal(err)
	}
	w := operatorFinishWork{prep: p, generation: 1, manifest: m, report: f.rep, git: g, spool: f.root,
		read: ro, now: func() int64 { return now }, budget: (serviceEnv{operatorBudget: activationTestBudget}).operatorContext}
	w.fresh = func(ctx context.Context) error {
		_, err := verifyOperatorFinishPrep(ctx, ro, a, now, cas, owner)
		return err
	}
	w.run = func(context.Context, store.StartCheck) (checkrun.Result, error) {
		if j := operatorFinishJournal(t, w); !j.OpenCheck || j.CommitPhase != "" {
			t.Fatalf("check ran without its reservation: %+v", j)
		}
		f.checks++
		return checkrun.Result{ExitCode: 0}, nil
	}
	w.record = func(context.Context, demoObject) (int64, error) {
		f.outcomes++
		return int64(f.outcomes), nil
	}
	w.evidence = func(context.Context, store.OperatorFinishCheck) bool { return f.outcomes > 0 }
	w.commit = func(ctx context.Context, spec store.OperatorLocalCommit, b [][]byte, tree string) error {
		return g.commit(ctx, f.root, spec, b, tree, f.rep.Start)
	}
	t.Cleanup(func() {
		row, err := ro.OperatorAction(activationTestContext(t), p.row.ID, now)
		if err != nil || row.V3ID != "" || len(f.seam.calls) != 0 || len(f.scopes) != 0 || len(f.finishRecords) != 0 || len(f.requests) != 0 {
			t.Errorf("private preparation leaked into dispatch: V3=%q, %v, %v", row.V3ID, f.seam.calls, err)
		}
	})
	return f, w, &now
}

func operatorFinishJournal(t *testing.T, w operatorFinishWork) store.OperatorFinishJournal {
	t.Helper()
	j, err := w.read.ReadOperatorFinishJournal(activationTestContext(t), operatorPreparationCAS(w.prep.row))
	if err != nil {
		t.Fatal(err)
	}
	return j
}

func TestOperatorFinishPrep(t *testing.T) {
	f, w, _ := operatorLocalFinishFixture(t, nil)
	baseline := operatorFinishJournal(t, w)
	unheld := w
	unheld.prep.row.Owner = ""
	if err := runEpisodeFinishWithPrep(activationTestContext(t), unheld); !errors.Is(err, approval.ErrNotVerified) ||
		!reflect.DeepEqual(baseline, operatorFinishJournal(t, w)) || f.checks != 0 {
		t.Fatalf("unclaimed proof crossed the entry: %v", err)
	}
	original := w.manifest.Commit.Message
	w.manifest.Commit.Message = "different local commit"
	if err := runEpisodeFinishWithPrep(activationTestContext(t), w); !errors.Is(err, errOperatorFinishReview) ||
		!reflect.DeepEqual(baseline, operatorFinishJournal(t, w)) || f.checks != 0 {
		t.Fatalf("other signed intent was admitted: %v", err)
	}
	w.manifest.Commit.Message = original
	if !w.boundManifest(w.generation) {
		t.Fatal("signed manifest must bind original inputs")
	}
	for _, kind := range []string{"protected", "unattributed"} {
		t.Run(kind, func(t *testing.T) {
			rig := *f
			rig.t = t
			rig.health.CredentialsDB = filepath.Join(t.TempDir(), "credentials.db")
			_, blocked, _ := operatorLocalFinishOn(t, &rig, w.git, w.manifest.Commit.Snapshot, func(m *store.OperatorFinishReview) {
				if kind == "protected" {
					m.Protected = []string{"a.go"}
				} else {
					m.Unattributed = []string{"b.go"}
				}
			})
			before := operatorFinishJournal(t, blocked)
			if err := runEpisodeFinishWithPrep(activationTestContext(t), blocked); !errors.Is(err, errOperatorFinishFence) ||
				!reflect.DeepEqual(before, operatorFinishJournal(t, blocked)) || rig.checks+rig.outcomes != 0 {
				t.Fatalf("unsupported %s change reached execution: %v", kind, err)
			}
		})
	}
}

func TestOperatorCommitReconcile(t *testing.T) {
	f, w, _ := operatorLocalFinishFixture(t, nil)
	sentinel := filepath.Join(f.root, "unexpected-hook")
	for _, name := range []string{"pre-commit", "reference-transaction"} {
		hook := filepath.Join(w.git.common, "hooks", name)
		if err := os.WriteFile(hook, []byte("#!/bin/sh\ntouch "+shellQuote(sentinel)+"\n"), 0700); err != nil {
			t.Fatal(err)
		}
	}
	original := w.commit
	w.commit = func(ctx context.Context, spec store.OperatorLocalCommit, b [][]byte, tree string) error {
		j := operatorFinishJournal(t, w)
		if j.CommitPhase != "reserved" || !strings.HasSuffix(j.CommitKey.Identity, ":"+tree) || len(j.Checks) != 1 || j.Checks[0].Result != "passed" {
			t.Fatalf("commit without exact reserved intent: %+v", j)
		}
		if err := original(ctx, spec, b, tree); err != nil {
			t.Fatalf("fixture commit did not reach its lost-reply boundary: %v", err)
		}
		return io.ErrUnexpectedEOF // bounded commit advanced HEAD, reply was lost
	}
	ctx := activationTestContext(t)
	if err := runEpisodeFinishWithPrep(ctx, w); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(sentinel); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("a repository hook executed or sentinel unreadable: %v", err)
	}
	seen := operatorFinishJournal(t, w)
	head := strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD"))
	if head == f.before || seen.Commit != (store.OperatorFinishLocalCommit{State: store.OperatorLocalCommitDone, SHA: head}) {
		t.Fatalf("hand-off %+v, head %s", seen.Commit, head)
	}
	w.commit = func(context.Context, store.OperatorLocalCommit, [][]byte, string) error {
		t.Fatal("repeated an observed commit")
		return nil
	}
	if err := runEpisodeFinishWithPrep(activationTestContext(t), w); err != nil || f.checks != 1 || f.outcomes != 1 {
		t.Fatalf("resume: %v, checks %d, outcomes %d", err, f.checks, f.outcomes)
	}
}

func TestOperatorCommitFailed(t *testing.T) {
	f, w, _ := operatorLocalFinishFixture(t, nil)
	n := 0
	w.commit = func(context.Context, store.OperatorLocalCommit, [][]byte, string) error {
		n++
		if got := operatorFinishJournal(t, w); got.CommitPhase != "reserved" {
			t.Fatal("commit attempt preceded reservation")
		}
		return errOperatorCommitNoRef
	}
	for range 2 {
		if err := runEpisodeFinishWithPrep(activationTestContext(t), w); !errors.Is(err, errOperatorCommitFailed) {
			t.Fatal(err)
		}
	}
	j := operatorFinishJournal(t, w)
	if n != 1 || len(j.Checks) != 1 || j.Commit.State != "commit-failed" || j.Commit.SHA != "" ||
		strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != f.before {
		t.Fatalf("failure lost evidence or retried: %+v attempts=%d", j, n)
	}
}

func TestOperatorFinishChecks(t *testing.T) {
	shared := newFinishFixture(t)
	git, snap := operatorFinishFixtureSnapshot(t, shared)
	for _, name := range []string{"failure", "abort", "timeout", "unknown-outcome", "changed-input"} {
		t.Run(name, func(t *testing.T) {
			rig := *shared
			rig.t = t
			rig.health.CredentialsDB = filepath.Join(t.TempDir(), "credentials.db")
			f, w, _ := operatorLocalFinishOn(t, &rig, git, snap, func(m *store.OperatorFinishReview) {
				if name == "changed-input" {
					next := m.Checks[0]
					next.Name = "after-drift"
					m.Checks = append(m.Checks, next)
				}
			})
			w.run = func(context.Context, store.StartCheck) (checkrun.Result, error) {
				f.checks++
				switch name {
				case "failure":
					return checkrun.Result{ExitCode: 1}, nil
				case "abort":
					return checkrun.Result{}, context.Canceled
				case "timeout":
					return checkrun.Result{}, context.DeadlineExceeded
				default:
					f.write("b.go", "different bytes")
					return checkrun.Result{ExitCode: 0}, nil
				}
			}
			if name == "unknown-outcome" {
				w.run = func(context.Context, store.StartCheck) (checkrun.Result, error) {
					f.checks++
					return checkrun.Result{ExitCode: 0}, nil
				}
				w.record = func(context.Context, demoObject) (int64, error) { return 0, nil }
			}
			want := map[string]error{"unknown-outcome": errOperatorFinishUncertain, "failure": errOperatorFinishCheck, "abort": errOperatorFinishUncertain,
				"timeout": errOperatorFinishUncertain, "changed-input": errOperatorFinishReview}[name]
			for range 2 {
				if err := runEpisodeFinishWithPrep(activationTestContext(t), w); !errors.Is(err, want) {
					t.Fatal(err)
				}
			}
			j := operatorFinishJournal(t, w)
			lost := name == "abort" || name == "timeout" || name == "unknown-outcome"
			if f.checks != 1 || j.OpenCheck != lost || j.CommitPhase != "" || strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != f.before {
				t.Fatalf("check state changed after refusal: %d %+v", f.checks, j)
			}
		})
	}
}

func TestOperatorFinishUnsandboxed(t *testing.T) {
	f, w, _ := operatorLocalFinishFixture(t, func(m *store.OperatorFinishReview) { m.Checks[0].Unsandboxed = true })
	for _, resume := range []bool{false, true} {
		if resume {
			if err := w.reserve(w.step("checks", "00")); err != nil {
				t.Fatal(err)
			}
		}
		before := operatorFinishJournal(t, w)
		if err := runEpisodeFinishWithPrep(activationTestContext(t), w); !errors.Is(err, errOperatorFinishTerminal) ||
			f.checks != 0 || !reflect.DeepEqual(before, operatorFinishJournal(t, w)) {
			t.Fatalf("unsandboxed resume=%v: %v", resume, err)
		}
	}
}

func TestOperatorFinishPrepExpiry(t *testing.T) {
	f, w, now := operatorLocalFinishFixture(t, nil)
	w.run = func(context.Context, store.StartCheck) (checkrun.Result, error) {
		f.checks++
		*now = w.prep.row.ExpiresAt
		return checkrun.Result{ExitCode: 0}, nil
	}
	if err := runEpisodeFinishWithPrep(activationTestContext(t), w); !errors.Is(err, store.ErrOperatorExpired) {
		t.Fatal(err)
	}
	j := operatorFinishJournal(t, w)
	if len(j.Checks) != 1 || j.CommitPhase != "" || f.checks != 1 {
		t.Fatalf("expiry lost check receipt or advanced commit: %+v", j)
	}
}

func TestOperatorFinishStale(t *testing.T) {
	f, w, _ := operatorLocalFinishFixture(t, nil)
	f.write("b.go", "changed untracked contents")
	before := serviceFileHash(w.manifest.Store)
	if err := runEpisodeFinishWithPrep(activationTestContext(t), w); !errors.Is(err, errOperatorFinishReview) ||
		serviceFileHash(w.manifest.Store) != before || f.checks != 0 || f.outcomes != 0 {
		t.Fatalf("untracked drift crossed preflight: %v", err)
	}
	f.write("b.go", "hand\n")
	f.gitOutput("-C", f.wt, "commit", "--allow-empty", "-qm", "foreign HEAD")
	if err := runEpisodeFinishWithPrep(activationTestContext(t), w); !errors.Is(err, errOperatorFinishReview) ||
		serviceFileHash(w.manifest.Store) != before || f.checks != 0 || f.outcomes != 0 {
		t.Fatalf("HEAD drift crossed preflight: %v", err)
	}
}

func TestOperatorFinishJournalOnce(t *testing.T) {
	_, w, now := operatorLocalFinishFixture(t, nil, true)
	ctx := activationTestContext(t)
	next, err := w.read.NextOperatorFinish(ctx, *now)
	if err != nil || next.ID != w.prep.row.ID || next.State != "confirmed" {
		t.Fatalf("confirmed candidate: %s %s %v", next.ID, next.State, err)
	}
	c := runnerConfig{Credentials: w.manifest.Store, FinishRemotes: []string{w.manifest.Remote.URL}}
	a, m, generation, err := finishManifestForRow(c, next)
	if err != nil || !reflect.DeepEqual(m, w.manifest) || a.ID != next.ID || generation != 1 {
		t.Fatalf("decode exact frozen intent: %d (%v)", generation, err)
	}
	c.FinishRemotes = nil
	if _, _, _, err := finishManifestForRow(c, next); !errors.Is(err, approval.ErrNotVerified) {
		t.Fatalf("unpinned remote: %v", err)
	}
	if _, err := w.read.NextOperatorFinish(ctx, next.ExpiresAt); !errors.Is(err, store.ErrOperatorMissing) {
		t.Fatalf("expired confirmed candidate: %v", err)
	}
	h, err := w.read.OperatorHeartbeat(ctx)
	if err != nil {
		t.Fatal(err)
	}
	cas := operatorPreparationCAS(next)
	if err := runnerWrite(w.manifest.Store, func(db *store.Credentials) error {
		write, stop := w.budget(time.Second)
		defer stop()
		return db.ClaimOperatorAction(write, cas, h.Owner, *now)
	}); err != nil {
		t.Fatal(err)
	}
	w.prep.row, err = w.read.NextOperatorFinish(ctx, *now)
	if err != nil || w.prep.row.Owner != h.Owner.Nonce || w.prep.row.Generation != cas.Generation+1 {
		t.Fatalf("claim candidate: %s %v", w.prep.row.State, err)
	}
	key := w.step("commit", strings.Repeat("a", 40))
	if err := w.reserve(key); err != nil {
		t.Fatal(err)
	}
	other := w.step("commit", strings.Repeat("b", 40))
	for _, done := range []bool{false, true} {
		if done {
			if err := w.observe(key, store.OperatorFinishLocalCommit{State: store.OperatorLocalCommitDone, SHA: strings.Repeat("c", 40)}); err != nil {
				t.Fatal(err)
			}
		}
		before := operatorFinishJournal(t, w)
		if err := w.reserve(other); !errors.Is(err, store.ErrOperatorConflict) {
			t.Fatalf("second distinct commit intent after observed=%t: %v", done, err)
		}
		if after := operatorFinishJournal(t, w); after.CommitKey.Identity != key.Identity || !reflect.DeepEqual(before, after) {
			t.Fatalf("commit reservation changed after refusal: %+v", after)
		}
	}
	if _, err := w.read.NextOperatorFinish(ctx, *now); !errors.Is(err, store.ErrOperatorMissing) {
		t.Fatalf("completed candidate selected again: %v", err)
	}
}
