package main

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"errors"
	"io"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

func TestOperatorFinishF1(t *testing.T) {
	f := newFinishFixture(t)
	t.Setenv("WORLD_OPERATOR_CONFIG_SHA256", strings.Repeat("a", 64))
	t.Setenv("WORLD_OPERATOR_FINISH_PREP", fixtureActionID)
	head := strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "--path-format=absolute", "--git-path", "HEAD"))
	common := filepath.Join(f.root, "repo", ".git")
	snapshot := func() []any {
		return []any{
			serviceFileHash(head), serviceFileHash(filepath.Join(common, "packed-refs")),
			provisionFiles(t, filepath.Join(common, "refs")),
			serviceFileHash(filepath.Join(filepath.Dir(head), "index")),
			serviceFileHash(f.health.CredentialsDB),
			serviceFileHash(filepath.Join(f.sessions, "ep1.pr-body.md")),
		}
	}
	before := snapshot()
	oldGit := sessionGit
	t.Cleanup(func() { sessionGit = oldGit })
	calls := 0
	sessionGit = func() (string, error) { calls++; return f.git, nil }
	for _, flags := range [][]string{nil, {"--confirm-on-surface"}, {"--operator-prep", fixtureActionID}, {"--yes"}, {"--confirm-on-surface", "--dry-run"}} {
		rc, _, errOut := f.run(noTerminalEnv(1), flags...)
		if rc != exitUsage || !strings.Contains(errOut, "refusing: no controlling terminal") ||
			calls != 0 || len(f.requests) != 0 || len(f.seam.calls) != 0 || f.checks != 0 || f.outcomes != 0 ||
			len(f.finishRecords) != 0 || !reflect.DeepEqual(before, snapshot()) {
			t.Fatalf("F1 flags %q: rc %d, git %d, requests %v, dispatch %v, %s", flags, rc, calls, f.requests, f.seam.calls, errOut)
		}
	}
}

type operatorFinishTTY struct {
	io.Reader
	bytes.Buffer
}

func (t *operatorFinishTTY) Read(b []byte) (int, error) { return t.Reader.Read(b) }
func (t *operatorFinishTTY) Close() error               { return nil }

func TestOperatorFinishTTYParity(t *testing.T) {
	f := newFinishFixture(t)
	tty := &operatorFinishTTY{Reader: strings.NewReader("y\n")}
	env := testMintEnv("", 1)
	env.openTerminal = func() (io.ReadWriteCloser, error) { return tty, nil }
	rc, _, errOut := f.run(env, "--include-unattributed", "--title", "Reviewed fixture")
	if rc != exitOK || len(f.scopes) != 1 || len(f.finishRecords) != 1 {
		t.Fatalf("typed finish: rc %d scopes %d records %d %s", rc, len(f.scopes), len(f.finishRecords), errOut)
	}
	scope, err := broker.ParseFinishApprovalScope(f.scopes[0])
	if err != nil {
		t.Fatal(err)
	}
	want := broker.FinishPromptText(scope, f.seam.title, f.seam.body, nil, "") + "Open this pull request? [y/N] "
	if tty.String() != want {
		t.Fatalf("terminal prompt changed\ngot %q\nwant %q", tty.String(), want)
	}
}

func TestOperatorFinishNoProof(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "creds.db")
	db, err := store.OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx := boundedTestContext(t)
	frozen := approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: fixtureActionID, Kind: "finish",
		Episode: "ep1", Store: path, Options: approval.OperatorOptions{Preset: "fixture", ConfirmOnSurface: true},
		Review: "finish fixture", Resolution: []approval.OperatorResolved{{Name: "manifest", Value: "fixture"}}}
	frozen.Argv = approval.OperatorArgv(frozen.Kind, frozen.Episode, frozen.Options)
	request := store.OperatorRequest{Action: frozen, OperationID: strings.Repeat("1", 32),
		Challenge: sha256.Sum256([]byte(t.Name())), CreatedAt: 1000, ExpiresAt: 1600}
	row, err := db.RequestOperatorAction(ctx, request, 1000)
	if err != nil {
		t.Fatal(err)
	}
	key := approvaltest.NewSigner(51)
	if err := db.EnrolApprover(ctx, store.ApproverRow{Name: "fixture", Format: "ssh", PublicKey: key.PublicKey().Blob, EnrolledAt: 900}); err != nil {
		t.Fatal(err)
	}
	cas := operatorPreparationCAS(row)
	refuse := func(action approval.OperatorAction, at int64, cas store.OperatorCAS, owner string, want error) {
		t.Helper()
		before := serviceFileHash(path)
		got, err := verifyOperatorFinishPrep(ctx, db, action, at, cas, owner)
		if !errors.Is(err, want) || got.claimed() || serviceFileHash(path) != before {
			t.Fatalf("want %v without state change, got %v", want, err)
		}
	}
	refuse(frozen, 1010, cas, "", approval.ErrNotVerified)
	if (operatorFinishPrep{}).claimed() {
		t.Fatal("zero-value capability opened the claimed gate")
	}

	msg, err := approval.OperatorActionMessage(cas.Digest, request.Challenge, request.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	sig := key.Sign(approval.OperatorActionNamespace, msg)
	proof, err := authority.JudgeOperatorAction(ctx, db, authority.OperatorJudgeInput{Action: frozen, Digest: cas.Digest,
		BytesDigest: sha256.Sum256(row.Action), Challenge: request.Challenge, CreatedAt: 1000, ExpiresAt: 1600, Now: 1010}, sig)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.ConfirmOperatorAction(ctx, cas, sig, proof, 1010); err != nil {
		t.Fatal(err)
	}
	writer, err := sql.Open("sqlite", runnerDB(path))
	if err != nil {
		t.Fatal(err)
	}
	defer writer.Close()
	setArtifact := func(artifact []byte) {
		t.Helper()
		if _, err := writer.ExecContext(ctx, "UPDATE operator_actions SET artifact=? WHERE id=?", artifact, cas.ID); err != nil {
			t.Fatal(err)
		}
	}
	setArtifact(approvaltest.NewSigner(52).Sign(approval.OperatorActionNamespace, msg))
	refuse(frozen, 1010, cas, "", approval.ErrNotVerified)
	setArtifact(sig)
	verified, err := verifyOperatorFinishPrep(ctx, db, frozen, 1010, cas, "")
	if err != nil || verified.claimed() {
		t.Fatalf("confirmed proof must precede claim: %v", err)
	}
	foreign := frozen
	foreign.Episode = "ep-other"
	foreign.Argv = approval.OperatorArgv(foreign.Kind, foreign.Episode, foreign.Options)
	refuse(foreign, 1010, cas, "", approval.ErrNotVerified)
	foreign = frozen
	foreign.Resolution = []approval.OperatorResolved{{Name: "manifest", Value: "different"}}
	refuse(foreign, 1010, cas, "", approval.ErrNotVerified)

	h := store.OperatorHeartbeat{Owner: store.OperatorOwner{Nonce: "test-runner"}, Protocol: store.OperatorRunnerProtocol, Build: "fixture", Readiness: "both", SeenAt: 1010}
	if h.Owner.Generation, err = db.PutOperatorHeartbeat(ctx, h); err != nil {
		t.Fatal(err)
	}
	if err := db.ClaimOperatorAction(ctx, cas, h.Owner, 1011); err != nil {
		t.Fatal(err)
	}
	refuse(frozen, 1011, cas, "", approval.ErrNotVerified)
	refuse(frozen, 1011, cas, h.Owner.Nonce, approval.ErrNotVerified)
	cas.Generation++
	refuse(frozen, 1011, cas, "other-runner", approval.ErrNotVerified)
	refuse(frozen, 1600, cas, h.Owner.Nonce, store.ErrOperatorExpired)
	verified, err = verifyOperatorFinishPrep(ctx, db, frozen, 1011, cas, h.Owner.Nonce)
	if err != nil || !verified.claimed() {
		t.Fatalf("live held proof failed: %v", err)
	}
	verified.authority = proof
	if verified.claimed() {
		t.Fatal("daemon's confirmation proof skipped independent runner recheck")
	}
}

func TestOperatorFinishProofScope(t *testing.T) {
	for _, tc := range []struct {
		kind, task string
		want       error
	}{
		{"start", "reviewed task", approval.ErrNotVerified},
		{"finish", "", errOperatorFinishSurface},
	} {
		t.Run(tc.kind, func(t *testing.T) {
			dir, err := filepath.EvalSymlinks(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(dir, "credentials.db")
			db, err := store.OpenCredentials(path)
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			ctx := boundedTestContext(t)
			const now = 1000
			key := approvaltest.NewSigner(51)
			if err := db.EnrolApprover(activationTestContext(t), store.ApproverRow{Name: "fixture", Format: "ssh", PublicKey: key.PublicKey().Blob, EnrolledAt: now - 1}); err != nil {
				t.Fatal(err)
			}
			a := approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: fixtureActionID, Kind: tc.kind,
				Episode: "ep1", Store: path, Task: tc.task, Options: approval.OperatorOptions{Preset: "fixture"},
				Review: "scope fixture", Resolution: []approval.OperatorResolved{{Name: "manifest", Value: "fixture"}}}
			a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
			req := store.OperatorRequest{Action: a, OperationID: strings.Repeat("2", 32),
				Challenge: sha256.Sum256([]byte(t.Name())), CreatedAt: now, ExpiresAt: now + 600}
			row, err := db.RequestOperatorAction(activationTestContext(t), req, now)
			if err != nil {
				t.Fatal(err)
			}
			cas := operatorPreparationCAS(row)
			msg, err := approval.OperatorActionMessage(cas.Digest, req.Challenge, req.ExpiresAt)
			if err != nil {
				t.Fatal(err)
			}
			sig := key.Sign(approval.OperatorActionNamespace, msg)
			proof, err := authority.JudgeOperatorAction(ctx, db, authority.OperatorJudgeInput{Action: a, Digest: cas.Digest,
				BytesDigest: sha256.Sum256(row.Action), Challenge: req.Challenge, CreatedAt: now, ExpiresAt: req.ExpiresAt, Now: now + 1}, sig)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := db.ConfirmOperatorAction(activationTestContext(t), cas, sig, proof, now+1); err != nil {
				t.Fatal(err)
			}
			before, err := db.OperatorAction(ctx, row.ID, now+1)
			if err != nil || before.State != "confirmed" {
				t.Fatalf("signed control: %s %v", before.State, err)
			}
			digest := serviceFileHash(path)
			got, err := verifyOperatorFinishPrep(ctx, db, a, now+1, cas, "")
			seen, readErr := db.OperatorAction(ctx, row.ID, now+1)
			if !errors.Is(err, tc.want) || got.authority != nil || readErr != nil ||
				!reflect.DeepEqual(before, seen) || digest != serviceFileHash(path) {
				t.Fatalf("scope refusal changed evidence: %v, row %s (%v)", err, seen.State, readErr)
			}
		})
	}
}
