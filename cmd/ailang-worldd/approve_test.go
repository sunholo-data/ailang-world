package main

// Row 201 A0 (design §2.5, §2.6a): `episode start --request` stores a mint request and opens no
// terminal; `approve` (M6b) mints it. The helpers here are shared by both.

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

// approveT0 is the fixed creation time every approve test uses (2027-01-15T08:00:00Z).
const approveT0 int64 = 1800000000

// seqRandom fills each call's buffer with the next bytes of 1, 2, 3, ... (mod 256), so the first
// request's id is mr-0102030405060708 and its challenge is 09..28.
func seqRandom() func([]byte) error {
	next := byte(1)
	return func(b []byte) error {
		for i := range b {
			b[i] = next
			next++
		}
		return nil
	}
}

// requestEnv is episode start --request's env: openTerminal records into *opened and fails,
// now returns now, random is seqRandom().
func requestEnv(opened *bool, now int64) sessionEnv {
	return sessionEnv{
		openTerminal: func() (io.ReadWriteCloser, error) { *opened = true; return nil, os.ErrNotExist },
		now:          func() int64 { return now },
		random:       seqRandom(),
	}
}

// clockSeq returns before for the first n calls and after from then on (D8 call order).
func clockSeq(n int, before, after int64) func() int64 {
	calls := 0
	return func() int64 {
		calls++
		if calls <= n {
			return before
		}
		return after
	}
}

type approveRig struct {
	r        episodeRun
	ep, id   string
	credPath string
	profile  *broker.ExecProfile
	task     string
}

// newApproveRig writes profile "proj", a task file, stubs a matching daemon, and runs
// `episode start <ep> --profile proj --task <task> --model claude-opus-5-5 --request` at createdAt
// with requestEnv. It fails the test unless that exits 0 with no terminal opened, and returns the
// rig with the request id read back from PendingMintRequests.
func newApproveRig(t *testing.T, ep string, createdAt int64, extra ...string) approveRig {
	t.Helper()
	return newApproveRigIn(t, newEpisodeRun(t), ep, createdAt, extra...)
}

// newApproveRigIn is newApproveRig in an existing run, so several requests share one store
// (writeTestProfile rewrites the same bytes, so the digest is unchanged; stubDaemon is re-applied).
func newApproveRigIn(t *testing.T, r episodeRun, ep string, createdAt int64, extra ...string) approveRig {
	t.Helper()
	p := writeTestProfile(t, r.profiles, "proj", "")
	stubDaemon(t, matchingHealth(t, r.f, r.profiles, p), nil)
	task := filepath.Join(t.TempDir(), "task.md")
	if err := os.WriteFile(task, []byte("do the thing\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var opened bool
	var stdout, stderr bytes.Buffer
	args := r.args(ep, append([]string{"--profile", "proj", "--task", task, "--model", "claude-opus-5-5", "--request"}, extra...)...)
	if code := runEpisodeStart(args, &stdout, &stderr, requestEnv(&opened, createdAt)); code != exitOK || opened {
		t.Fatalf("episode start --request: rc %d, terminal opened %v\nstdout %s\nstderr %s", code, opened, stdout.String(), stderr.String())
	}
	credPath, err := credentialsFor(r.f.db, "")
	if err != nil {
		t.Fatal(err)
	}
	var id string
	for _, row := range pendingRows(t, r.f.db, createdAt) {
		if row.EpisodeID == ep {
			id = row.RequestID
		}
	}
	if id == "" {
		t.Fatalf("no pending request for %s", ep)
	}
	return approveRig{r: r, ep: ep, id: id, credPath: credPath, profile: p, task: task}
}

func (a approveRig) args(extra ...string) []string { // id, --db, db, --credentials-db, credPath, extra...
	return append([]string{a.id, "--db", a.r.f.db, "--credentials-db", a.credPath}, extra...)
}

func (a approveRig) row(t *testing.T) store.MintRequestRow {
	t.Helper()
	ro := openCredsReadOnly(t, a.r.f.db)
	defer func() { _ = ro.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row, ok, err := ro.MintRequest(ctx, a.id)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("no mint request %s", a.id)
	}
	return row
}

func (a approveRig) sessionRows(t *testing.T) []store.SessionRow {
	t.Helper()
	return credRows(t, a.r.f.db)
}

func (a approveRig) worktree() string { return filepath.Join(a.r.f.ws, a.ep) }

// openCredsReadOnly opens db's credentials store read-only (the caller closes it).
func openCredsReadOnly(t *testing.T, db string) *store.Credentials {
	t.Helper()
	path, err := credentialsFor(db, "")
	if err != nil {
		t.Fatal(err)
	}
	ro, err := store.OpenCredentialsReadOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	return ro
}

// pendingRows reads db's unconsumed, unexpired requests at now, read-only; no store means none.
func pendingRows(t *testing.T, db string, now int64) []store.MintRequestRow {
	t.Helper()
	path, err := credentialsFor(db, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	ro := openCredsReadOnly(t, db)
	defer func() { _ = ro.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := ro.PendingMintRequests(ctx, now, 0)
	if err != nil {
		t.Fatal(err)
	}
	return rows
}

// firstRequestID is the id episode start --request prints for the first request under seqRandom.
const firstRequestID = "mr-0102030405060708"

// seqChallenge is the challenge seqRandom draws second: bytes 9..40.
func seqChallenge() string {
	b := make([]byte, 32)
	for i := range b {
		b[i] = byte(9 + i)
	}
	return hex.EncodeToString(b)
}

// CL1: --request stores one pending request, opens no terminal, provisions nothing, and prints the
// id and the approve line.
func TestEpisodeStartRequestStoresAndOpensNoTerminal(t *testing.T) {
	r := newEpisodeRun(t)
	p := writeTestProfile(t, r.profiles, "proj", "")
	stubDaemon(t, matchingHealth(t, r.f, r.profiles, p), nil)
	task := filepath.Join(t.TempDir(), "task.md")
	if err := os.WriteFile(task, []byte("do the thing\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var opened bool
	var stdout, stderr bytes.Buffer
	args := r.args("ep1", "--profile", "proj", "--task", task, "--model", "claude-opus-5-5", "--request")
	if code := runEpisodeStart(args, &stdout, &stderr, requestEnv(&opened, approveT0)); code != exitOK {
		t.Fatalf("rc %d\nstderr %s", code, stderr.String())
	}
	if opened {
		t.Fatal("a terminal was opened")
	}
	rows := pendingRows(t, r.f.db, approveT0)
	if len(rows) != 1 {
		t.Fatalf("%d pending requests, want 1", len(rows))
	}
	row := rows[0]
	if row.RequestID != firstRequestID {
		t.Errorf("id %q, want %s", row.RequestID, firstRequestID)
	}
	if row.Challenge != seqChallenge() {
		t.Errorf("challenge %q, want %s", row.Challenge, seqChallenge())
	}
	if row.CreatedAt != approveT0 || row.ExpiresAt != approveT0+600 {
		t.Errorf("created %d expires %d, want %d and %d", row.CreatedAt, row.ExpiresAt, approveT0, approveT0+600)
	}
	sum := sha256.Sum256(row.Request)
	if row.Digest != hex.EncodeToString(sum[:]) {
		t.Errorf("digest %q is not sha256 of the request", row.Digest)
	}
	mr, err := approval.Decode(row.Request)
	if err != nil {
		t.Fatal(err)
	}
	wantGrants, err := episodeStartGrants("", nil, 50, true)
	if err != nil {
		t.Fatal(err)
	}
	absTask, _ := filepath.Abs(task)
	if mr.Episode != "ep1" || !reflect.DeepEqual(mr.Grants, fromCapabilities(wantGrants)) || mr.TTLSeconds != 3600 ||
		mr.Profile != (approval.ProfileRef{Project: "proj", Digest: p.Digest}) || mr.Agent != "claude" ||
		mr.Model != "claude-opus-5-5" || mr.Task != absTask ||
		mr.SessionFile != filepath.Join(r.sessions, "ep1.session") || mr.MCPConfig != filepath.Join(r.sessions, "ep1.mcp.json") ||
		mr.Daemon.Addr != defaultEpisodeAddr || mr.Daemon.VCSRevision != testBuild.VCSRevision {
		t.Errorf("request %+v", mr)
	}
	out := stdout.String()
	for _, want := range []string{firstRequestID, "ailang-worldd approve " + firstRequestID,
		"  grants    " + grantList(toCapabilities(mr.Grants))} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout lacks %q:\n%s", want, out)
		}
	}
	if _, err := os.Lstat(filepath.Join(r.f.ws, "ep1")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("the worktree exists (err %v)", err)
	}
	r.assertNoSessionFiles(t)
	if n := len(r.f.rows(t)); n != 0 {
		t.Errorf("%d session rows, want 0", n)
	}
}

// CL1b: a nil random source refuses before the store is touched.
func TestEpisodeStartRequestRefusesWithoutRandomness(t *testing.T) {
	r := newEpisodeRun(t)
	p := writeTestProfile(t, r.profiles, "proj", "")
	stubDaemon(t, matchingHealth(t, r.f, r.profiles, p), nil)
	var opened bool
	env := requestEnv(&opened, approveT0)
	env.random = nil
	var stdout, stderr bytes.Buffer
	args := r.args("ep1", "--profile", "proj", "--model", "claude-opus-5-5", "--request")
	if code := runEpisodeStart(args, &stdout, &stderr, env); code != exitFatal {
		t.Fatalf("rc %d, want exitFatal\nstderr %s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "no randomness source") {
		t.Errorf("stderr %q lacks 'no randomness source'", stderr.String())
	}
	if opened {
		t.Error("a terminal was opened")
	}
	if n := len(pendingRows(t, r.f.db, approveT0)); n != 0 {
		t.Errorf("%d pending requests, want 0", n)
	}
}

// CL2: --request refuses the modes it cannot carry, before any terminal or daemon; the control
// (--preset se-tools) is accepted.
func TestEpisodeStartRequestRefusesUnsupportedModes(t *testing.T) {
	cases := []struct {
		name string
		args func(r episodeRun) []string
		want string
	}{
		{"launch", func(r episodeRun) []string { return r.args("ep1", "--request", "--launch") },
			"--request cannot be combined with --launch"},
		{"headless", func(r episodeRun) []string { return r.args("ep1", "--request", "--headless") },
			"--request cannot be combined with --headless"},
		{"resume", func(r episodeRun) []string {
			return []string{"--request", "--resume", "ep1", "--preset", "p", "--db", r.f.db}
		}, "--request cannot be combined with --resume"},
		{"preset", func(r episodeRun) []string { return r.args("ep1", "--request", "--preset", "mypreset") },
			"--request cannot be combined with a launch preset"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := newEpisodeRun(t)
			p := writeTestProfile(t, r.profiles, "proj", "")
			stubDaemon(t, matchingHealth(t, r.f, r.profiles, p), nil)
			var opened bool
			var stdout, stderr bytes.Buffer
			if code := runEpisodeStart(c.args(r), &stdout, &stderr, requestEnv(&opened, approveT0)); code != exitUsage {
				t.Fatalf("rc %d, want exitUsage\nstderr %s", code, stderr.String())
			}
			if !strings.Contains(stderr.String(), c.want) {
				t.Errorf("stderr %q lacks %q", stderr.String(), c.want)
			}
			if opened {
				t.Error("a terminal was opened")
			}
			if n := len(pendingRows(t, r.f.db, approveT0)); n != 0 {
				t.Errorf("%d pending requests, want 0", n)
			}
		})
	}
	t.Run("control-se-tools", func(t *testing.T) {
		r := newEpisodeRun(t)
		p := writeTestProfile(t, r.profiles, "proj", "")
		stubDaemon(t, matchingHealth(t, r.f, r.profiles, p), nil)
		var opened bool
		var stdout, stderr bytes.Buffer
		if code := runEpisodeStart(r.args("ep1", "--request", "--preset", "se-tools"), &stdout, &stderr, requestEnv(&opened, approveT0)); code != exitOK {
			t.Fatalf("rc %d\nstderr %s", code, stderr.String())
		}
		if n := len(pendingRows(t, r.f.db, approveT0)); n != 1 {
			t.Errorf("%d pending requests, want 1", n)
		}
	})
}

// CL3: a v3 credentials store is refused with the restart line, and is not upgraded.
func TestEpisodeStartRequestRefusesAnOlderStore(t *testing.T) {
	r := newEpisodeRun(t)
	p := writeTestProfile(t, r.profiles, "proj", "")
	stubDaemon(t, matchingHealth(t, r.f, r.profiles, p), nil)
	credPath, err := credentialsFor(r.f.db, "")
	if err != nil {
		t.Fatal(err)
	}
	raw, err := sql.Open("sqlite", credPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = raw.Close() }()
	for _, stmt := range []string{
		`CREATE TABLE session_credentials (
    credential_id TEXT PRIMARY KEY,
    episode_id   TEXT NOT NULL CHECK (episode_id <> ''),
    grants_json  TEXT NOT NULL,
    expires_at   INTEGER NOT NULL,
    created_at   INTEGER NOT NULL,
    project TEXT CHECK (project IS NULL OR project <> ''),
    profile_digest TEXT CHECK ((profile_digest IS NULL) = (project IS NULL) AND (profile_digest IS NULL OR profile_digest <> ''))
)`,
		`CREATE TABLE episode_starts (
    episode TEXT PRIMARY KEY CHECK (episode <> ''),
    start_json TEXT NOT NULL
)`,
		`PRAGMA user_version = 3`,
	} {
		if _, err := raw.Exec(stmt); err != nil {
			t.Fatal(err)
		}
	}
	var opened bool
	var stdout, stderr bytes.Buffer
	args := r.args("ep1", "--grant", "Workspace.Read=worktree:1", "--request")
	if code := runEpisodeStart(args, &stdout, &stderr, requestEnv(&opened, approveT0)); code != exitUsage {
		t.Fatalf("rc %d, want exitUsage\nstderr %s", code, stderr.String())
	}
	for _, want := range []string{"mint requests need schema 4", "restart `ailang-worldd serve`"} {
		if !strings.Contains(stderr.String(), want) {
			t.Errorf("stderr lacks %q:\n%s", want, stderr.String())
		}
	}
	var version int
	if err := raw.QueryRow(`PRAGMA user_version`).Scan(&version); err != nil || version != 3 {
		t.Errorf("user_version %d (%v), want 3", version, err)
	}
	var n int
	if err := raw.QueryRow(`SELECT count(*) FROM sqlite_master WHERE name = 'mint_requests'`).Scan(&n); err != nil || n != 0 {
		t.Errorf("mint_requests tables %d (%v), want 0", n, err)
	}
}
