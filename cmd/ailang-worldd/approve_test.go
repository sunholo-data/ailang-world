package main

// Row 201 A0 (design §2.5, §2.6a): `episode start --request` stores a mint request and opens no
// terminal; `approve` (M6b) mints it. The helpers here are shared by both.

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
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
	// requestEnv's seqRandom restarts at 1 on every call, so a second request in one store would
	// draw the first one's id and collide on its primary key; the rig takes the production source.
	env := requestEnv(&opened, createdAt)
	env.random = realSessionEnv().random
	if code := runEpisodeStart(args, &stdout, &stderr, env); code != exitOK || opened {
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

// approveEnv is approve's env: openTerminal returns term (and records into *opened when non-nil),
// now is clock, after is time.After.
func approveEnv(term io.ReadWriteCloser, opened *bool, clock func() int64) sessionEnv {
	return sessionEnv{
		openTerminal: func() (io.ReadWriteCloser, error) {
			if opened != nil {
				*opened = true
			}
			return term, nil
		},
		now:   clock,
		after: time.After,
	}
}

// CL17: the unbounded callers (session new, episode start) never reach after: their y/N has no
// deadline (the plan's answerBy == 0 rule).
func TestExistingCallersReadUnbounded(t *testing.T) {
	after := func(time.Duration) <-chan time.Time {
		t.Error("after called on an unbounded caller")
		return nil
	}
	t.Run("session-new", func(t *testing.T) {
		f := newEpisodeFixture(t)
		env := sessionEnv{
			openTerminal: func() (io.ReadWriteCloser, error) { return &fakeTerm{r: strings.NewReader("y\n")}, nil },
			now:          func() int64 { return 1000 },
			after:        after,
		}
		var stdout, stderr bytes.Buffer
		if code := runSessionNew(f.args("ep1", "--grant", "Workspace.Read=worktree:1"), &stdout, &stderr, env); code != exitOK {
			t.Fatalf("session new rc %d\nstderr %s", code, stderr.String())
		}
		if n := len(f.rows(t)); n != 1 {
			t.Errorf("%d session rows, want 1", n)
		}
	})
	t.Run("episode-start", func(t *testing.T) {
		r := newEpisodeRun(t)
		p := writeTestProfile(t, r.profiles, "proj", "")
		stubDaemon(t, matchingHealth(t, r.f, r.profiles, p), nil)
		env := sessionEnv{
			openTerminal: func() (io.ReadWriteCloser, error) { return &fakeTerm{r: strings.NewReader("y\n")}, nil },
			now:          func() int64 { return 1000 },
			after:        after,
		}
		var stdout, stderr bytes.Buffer
		if code := runEpisodeStart(r.args("ep2", "--profile", "proj"), &stdout, &stderr, env); code != exitOK {
			t.Fatalf("episode start rc %d\nstderr %s", code, stderr.String())
		}
	})
}

// CL18: readTermLineWithin returns the line when it arrives first, and gives up when the timer
// fires first, without waiting on the blocked reader.
func TestReadTermLineWithin(t *testing.T) {
	t.Run("line-first", func(t *testing.T) {
		never := func(time.Duration) <-chan time.Time { return make(chan time.Time) }
		line, ok := readTermLineWithin(strings.NewReader("y\n"), time.Second, never)
		if !ok || line != "y" {
			t.Fatalf("(%q, %v), want (\"y\", true)", line, ok)
		}
	})
	t.Run("timer-first", func(t *testing.T) {
		pr, pw := io.Pipe()
		t.Cleanup(func() { _ = pw.Close() })
		fired := func(time.Duration) <-chan time.Time {
			c := make(chan time.Time)
			close(c)
			return c
		}
		type result struct {
			line string
			ok   bool
		}
		done := make(chan result, 1)
		go func() {
			line, ok := readTermLineWithin(pr, time.Second, fired)
			done <- result{line, ok}
		}()
		select {
		case got := <-done:
			if got.ok || got.line != "" {
				t.Fatalf("(%q, %v), want (\"\", false)", got.line, got.ok)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("readTermLineWithin did not return on the timer")
		}
	})
}

// rawCredentials opens the credentials store at path read-write, for tests that edit a stored row.
func rawCredentials(t *testing.T, path string) *sql.DB {
	t.Helper()
	raw, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

// atTime is a clock that always reads v.
func atTime(v int64) func() int64 { return func() int64 { return v } }

// CL4: --list shows the pending requests only: not the expired one, not the consumed one.
func TestApproveListShowsPending(t *testing.T) {
	r := newEpisodeRun(t)
	a := newApproveRigIn(t, r, "ep1", approveT0-700)
	b := newApproveRigIn(t, r, "ep2", approveT0)
	c := newApproveRigIn(t, r, "ep3", approveT0)
	var stdout, stderr bytes.Buffer
	term := &promptTerm{in: strings.NewReader("y\n")}
	if code := runApprove(c.args(), &stdout, &stderr, approveEnv(term, nil, atTime(approveT0+1))); code != exitOK {
		t.Fatalf("approve C: rc %d\nstderr %s", code, stderr.String())
	}
	listArgs := []string{"--list", "--db", r.f.db, "--credentials-db", c.credPath}
	var listOut, listErr bytes.Buffer
	if code := runApprove(listArgs, &listOut, &listErr, approveEnv(nil, nil, atTime(approveT0+2))); code != exitOK {
		t.Fatalf("--list: rc %d\nstderr %s", code, listErr.String())
	}
	out := listOut.String()
	if !strings.Contains(out, b.id) {
		t.Errorf("--list lacks pending %s:\n%s", b.id, out)
	}
	if strings.Contains(out, a.id) {
		t.Errorf("--list shows the expired %s:\n%s", a.id, out)
	}
	if strings.Contains(out, c.id) {
		t.Errorf("--list shows the consumed %s:\n%s", c.id, out)
	}
	bRow := b.row(t)
	if !strings.Contains(out, "episode ep2") || !strings.Contains(out, bRow.Digest[:16]) {
		t.Errorf("B's line lacks its episode or digest:\n%s", out)
	}
	listOut.Reset()
	if code := runApprove(listArgs, &listOut, &listErr, approveEnv(nil, nil, atTime(approveT0+1000))); code != exitOK {
		t.Fatalf("--list at +1000: rc %d", code)
	}
	if got := strings.TrimSpace(listOut.String()); got != "no pending mint requests" {
		t.Errorf("--list with nothing pending printed %q", got)
	}
}

// CL6: without a terminal, approve mints nothing and the request stays pending.
func TestApproveWithoutATerminalMintsNothing(t *testing.T) {
	a := newApproveRig(t, "ep1", approveT0)
	env := noTerminalEnv(approveT0 + 1)
	env.after = time.After
	var stdout, stderr bytes.Buffer
	if code := runApprove(a.args(), &stdout, &stderr, env); code != exitUsage {
		t.Fatalf("rc %d, want exitUsage\nstderr %s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), sessionTTYFix) {
		t.Errorf("stderr lacks the tty fix:\n%s", stderr.String())
	}
	if a.row(t).Consumed {
		t.Error("the request was consumed without a terminal")
	}
	if _, err := os.Lstat(a.worktree()); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("the worktree exists (err %v)", err)
	}
	if n := len(a.sessionRows(t)); n != 0 {
		t.Errorf("%d session rows, want 0", n)
	}
}

// CL7: a declined y/N leaves the request pending, and the prompt shows the whole request.
func TestApproveDeclinedLeavesTheRequestPending(t *testing.T) {
	a := newApproveRig(t, "ep1", approveT0)
	term := &promptTerm{in: strings.NewReader("n\n")}
	var stdout, stderr bytes.Buffer
	if code := runApprove(a.args(), &stdout, &stderr, approveEnv(term, nil, atTime(approveT0+1))); code != exitUsage {
		t.Fatalf("rc %d, want exitUsage\nstderr %s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "aborted (not confirmed)") {
		t.Errorf("stderr lacks 'aborted (not confirmed)': %s", stderr.String())
	}
	if a.row(t).Consumed {
		t.Error("a declined request was consumed")
	}
	if _, err := os.Lstat(a.worktree()); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("the worktree exists (err %v)", err)
	}
	if n := len(a.sessionRows(t)); n != 0 {
		t.Errorf("%d session rows, want 0", n)
	}
	if !strings.Contains(term.out.String(), "mint request "+a.id) {
		t.Errorf("the prompt lacks the request block:\n%s", term.out.String())
	}
}

// CL8: a used request is refused before any terminal opens, and mints once.
func TestApproveRefusesReplay(t *testing.T) {
	a := newApproveRig(t, "ep1", approveT0)
	var stdout, stderr bytes.Buffer
	term := &promptTerm{in: strings.NewReader("y\n")}
	if code := runApprove(a.args(), &stdout, &stderr, approveEnv(term, nil, atTime(approveT0+1))); code != exitOK {
		t.Fatalf("first approve: rc %d\nstderr %s", code, stderr.String())
	}
	if row := a.row(t); !row.Consumed || row.ApprovedBy != "tty" {
		t.Fatalf("after the mint: consumed %v by %q", row.Consumed, row.ApprovedBy)
	}
	var opened bool
	stdout.Reset()
	stderr.Reset()
	if code := runApprove(a.args(), &stdout, &stderr, approveEnv(nil, &opened, atTime(approveT0+2))); code != exitUsage {
		t.Fatalf("second approve: rc %d, want exitUsage\nstderr %s", code, stderr.String())
	}
	for _, want := range []string{"already used", "by tty"} {
		if !strings.Contains(stderr.String(), want) {
			t.Errorf("stderr lacks %q: %s", want, stderr.String())
		}
	}
	if opened {
		t.Error("a terminal opened for a used request")
	}
	if n := len(a.sessionRows(t)); n != 1 {
		t.Errorf("%d session rows, want 1", n)
	}
}

// CL12: approve refuses bad ids and misuse before any terminal; --help is approveHelp.
func TestApproveRefusesBadIDsAndMisuse(t *testing.T) {
	a := newApproveRig(t, "ep1", approveT0)
	run1 := func(args []string, env sessionEnv) (int, string) {
		var stdout, stderr bytes.Buffer
		code := runApprove(args, &stdout, &stderr, env)
		return code, stderr.String()
	}
	t.Run("grammar", func(t *testing.T) {
		code, msg := run1(append([]string{"mr-XYZ"}, a.args()[1:]...), approveEnv(nil, nil, atTime(approveT0+1)))
		if code != exitUsage || !strings.Contains(msg, "is not a mint request id") {
			t.Errorf("rc %d, stderr %q", code, msg)
		}
	})
	t.Run("unknown", func(t *testing.T) {
		code, msg := run1(append([]string{"mr-0000000000000000"}, a.args()[1:]...), approveEnv(nil, nil, atTime(approveT0+1)))
		if code != exitUsage || !strings.Contains(msg, "no mint request") {
			t.Errorf("rc %d, stderr %q", code, msg)
		}
	})
	t.Run("addr", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		code := run([]string{"--addr", "127.0.0.1:1", "approve", "mr-0000000000000000",
			"--db", a.r.f.db, "--credentials-db", a.credPath}, &stdout, &stderr)
		if code != exitUsage || !strings.Contains(stderr.String(), "not valid for 'approve'") {
			t.Errorf("rc %d, stderr %q", code, stderr.String())
		}
	})
	t.Run("no-after", func(t *testing.T) {
		var opened bool
		env := approveEnv(nil, &opened, atTime(approveT0+1))
		env.after = nil
		code, msg := run1(a.args(), env)
		if code != exitFatal || !strings.Contains(msg, "no prompt timer") || opened {
			t.Errorf("rc %d, opened %v, stderr %q", code, opened, msg)
		}
	})
	t.Run("other-store", func(t *testing.T) {
		dir2 := t.TempDir()
		world2 := filepath.Join(dir2, "world.db")
		s, err := store.Open(world2)
		if err != nil {
			t.Fatal(err)
		}
		_ = s.Close()
		var opened bool
		args := []string{a.id, "--db", world2, "--credentials-db", a.credPath}
		code, msg := run1(args, approveEnv(nil, &opened, atTime(approveT0+1)))
		if code != exitUsage || !strings.Contains(msg, "was made for") || opened {
			t.Errorf("rc %d, opened %v, stderr %q", code, opened, msg)
		}
	})
	t.Run("help", func(t *testing.T) {
		for _, argv := range [][]string{{"approve", "--help"}, {"help", "approve"}} {
			var stdout, stderr bytes.Buffer
			if code := run(argv, &stdout, &stderr); code != exitOK || stdout.String() != approveHelp {
				t.Errorf("%q: rc %d, stdout is not approveHelp (stderr %q)", argv, code, stderr.String())
			}
		}
	})
}

// CL9: an expired request is refused before any terminal opens, and a y that lands after the
// expiry (the in-transaction clock reading) mints nothing and removes the worktree the run made.
func TestApproveRefusesExpired(t *testing.T) {
	exp := approveT0 + 600
	t.Run("before-the-terminal", func(t *testing.T) {
		a := newApproveRig(t, "ep1", approveT0)
		var opened bool
		term := &promptTerm{in: strings.NewReader("y\n")}
		var stdout, stderr bytes.Buffer
		code := runApprove(a.args(), &stdout, &stderr, approveEnv(term, &opened, atTime(exp)))
		if code != exitUsage || !strings.Contains(stderr.String(), "expired at") || opened {
			t.Fatalf("rc %d, opened %v, want exitUsage, \"expired at\", not opened\nstderr %s", code, opened, stderr.String())
		}
	})
	t.Run("in-the-transaction", func(t *testing.T) {
		a := newApproveRig(t, "ep1", approveT0)
		term := &promptTerm{in: strings.NewReader("y\n")}
		var stdout, stderr bytes.Buffer
		code := runApprove(a.args(), &stdout, &stderr, approveEnv(term, nil, clockSeq(2, approveT0+1, approveT0+601)))
		if code != exitFatal || !strings.Contains(stderr.String(), "expired before it was consumed") {
			t.Fatalf("rc %d, want exitFatal and \"expired before it was consumed\"\nstderr %s", code, stderr.String())
		}
		if a.row(t).Consumed {
			t.Error("an expired request was consumed")
		}
		if n := len(a.sessionRows(t)); n != 0 {
			t.Errorf("%d session rows, want 0", n)
		}
		if _, err := os.Lstat(a.worktree()); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("the worktree exists (err %v)", err)
		}
		if !strings.Contains(stderr.String(), "removed the worktree") {
			t.Errorf("stderr lacks 'removed the worktree': %s", stderr.String())
		}
	})
}

// CL10: a world that changed since the request was made is refused before any terminal opens,
// and the request stays pending. Each arm is a fresh rig.
func TestApproveRefusesAStaleWorld(t *testing.T) {
	refusedStale := func(t *testing.T, a approveRig, want string) {
		t.Helper()
		var opened bool
		term := &promptTerm{in: strings.NewReader("y\n")}
		var stdout, stderr bytes.Buffer
		code := runApprove(a.args(), &stdout, &stderr, approveEnv(term, &opened, atTime(approveT0+1)))
		if code != exitUsage || !strings.Contains(stderr.String(), want) {
			t.Fatalf("rc %d, want exitUsage and %q\nstderr %s", code, want, stderr.String())
		}
		if opened {
			t.Error("the terminal opened for a stale request")
		}
		if a.row(t).Consumed {
			t.Error("the stale request was consumed")
		}
		if n := len(a.sessionRows(t)); n != 0 {
			t.Errorf("%d session rows, want 0", n)
		}
	}
	t.Run("profile-digest", func(t *testing.T) {
		r := newEpisodeRun(t)
		a := newApproveRigIn(t, r, "ep1", approveT0)
		p2 := writeTestProfile(t, r.profiles, "proj", `, "env": {"APPROVE_TEST": "1"}`)
		stubDaemon(t, matchingHealth(t, r.f, r.profiles, p2), nil)
		refusedStale(t, a, "profile.digest")
	})
	t.Run("daemon-build", func(t *testing.T) {
		r := newEpisodeRun(t)
		a := newApproveRigIn(t, r, "ep1", approveT0)
		otherBuild := daemon.BuildIdentity{VCSRevision: strings.Repeat("f", 40)}
		h := matchingHealth(t, r.f, r.profiles, a.profile)
		h.Build = otherBuild
		stubDaemon(t, h, nil)
		episodeSelfBuild = func() daemon.BuildIdentity { return otherBuild }
		refusedStale(t, a, "daemon.vcs_revision")
	})
	t.Run("session-file-exists", func(t *testing.T) {
		a := newApproveRig(t, "ep1", approveT0)
		if err := os.MkdirAll(a.r.sessions, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(a.r.sessions, "ep1.session"), []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
		refusedStale(t, a, "already exists")
	})
	t.Run("mcp-config-exists", func(t *testing.T) {
		a := newApproveRig(t, "ep1", approveT0)
		if err := os.MkdirAll(a.r.sessions, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(a.r.sessions, "ep1.mcp.json"), []byte("{}"), 0o600); err != nil {
			t.Fatal(err)
		}
		refusedStale(t, a, "already exists")
	})
	t.Run("worktree-exists", func(t *testing.T) {
		a := newApproveRig(t, "ep1", approveT0)
		if out, err := exec.Command("git", "-C", a.r.f.repo, "worktree", "add", "--detach", a.worktree()).CombinedOutput(); err != nil {
			t.Fatalf("git worktree add: %v\n%s", err, out)
		}
		refusedStale(t, a, "reuse")
	})
}

// CL11: a request whose stored bytes no longer match their digest is refused before any terminal
// opens. Re-digested, the same edited bytes mint exactly what they say: the grant the session
// row carries is the one the request holds.
func TestApproveRefusesATamperedRow(t *testing.T) {
	tamper := func(t *testing.T, a approveRig, redigest bool) []byte {
		t.Helper()
		row := a.row(t)
		if !bytes.Contains(row.Request, []byte(`"budget":50`)) {
			t.Fatalf("the request lacks \"budget\":50: %s", row.Request)
		}
		edited := bytes.Replace(row.Request, []byte(`"budget":50`), []byte(`"budget":999`), 1)
		digest := row.Digest
		if redigest {
			sum := sha256.Sum256(edited)
			digest = hex.EncodeToString(sum[:])
		}
		raw := rawCredentials(t, a.credPath)
		defer func() { _ = raw.Close() }()
		if _, err := raw.Exec(`UPDATE mint_requests SET request = ?, digest = ? WHERE request_id = ?`, edited, digest, a.id); err != nil {
			t.Fatal(err)
		}
		return edited
	}
	t.Run("digest-kept", func(t *testing.T) {
		a := newApproveRig(t, "ep1", approveT0)
		tamper(t, a, false)
		var opened bool
		term := &promptTerm{in: strings.NewReader("y\n")}
		var stdout, stderr bytes.Buffer
		code := runApprove(a.args(), &stdout, &stderr, approveEnv(term, &opened, atTime(approveT0+1)))
		if code != exitUsage || !strings.Contains(stderr.String(), "do not match its digest") || opened {
			t.Fatalf("rc %d, opened %v, want exitUsage, \"do not match its digest\", not opened\nstderr %s", code, opened, stderr.String())
		}
	})
	t.Run("redigested", func(t *testing.T) {
		a := newApproveRig(t, "ep1", approveT0)
		tamper(t, a, true)
		term := &promptTerm{in: strings.NewReader("y\n")}
		var stdout, stderr bytes.Buffer
		if code := runApprove(a.args(), &stdout, &stderr, approveEnv(term, nil, atTime(approveT0+1))); code != exitOK {
			t.Fatalf("rc %d\nstderr %s", code, stderr.String())
		}
		if !strings.Contains(term.out.String(), "Workspace.Read=worktree:999") {
			t.Errorf("the prompt lacks the edited grant:\n%s", term.out.String())
		}
		rows := a.sessionRows(t)
		if len(rows) != 1 {
			t.Fatalf("%d session rows, want 1", len(rows))
		}
		var grants []broker.Capability
		if err := json.Unmarshal([]byte(rows[0].GrantsJSON), &grants); err != nil {
			t.Fatal(err)
		}
		found := false
		for _, g := range grants {
			if g.Effect == broker.EffectWorkspaceRead && g.Budget == 999 {
				found = true
			}
		}
		if !found {
			t.Errorf("grants %s lack Workspace.Read with budget 999", rows[0].GrantsJSON)
		}
	})
}

// CL13: an unanswered y/N gives up when the request expires: the wait is exactly the remaining
// lifetime, and nothing is consumed or minted. Arm 2: a request expired before the prompt is
// refused before the prompt is written.
func TestApproveUnansweredPromptTimesOut(t *testing.T) {
	exp := approveT0 + 600
	t.Run("no-answer", func(t *testing.T) {
		a := newApproveRig(t, "ep1", approveT0)
		pr, pw := io.Pipe()
		t.Cleanup(func() { _ = pw.Close() })
		term := &promptTerm{in: pr}
		var waited time.Duration
		env := approveEnv(term, nil, atTime(exp-37))
		env.after = func(d time.Duration) <-chan time.Time {
			waited = d
			c := make(chan time.Time)
			close(c)
			return c
		}
		var stdout, stderr bytes.Buffer
		var code int
		done := make(chan struct{})
		go func() {
			code = runApprove(a.args(), &stdout, &stderr, env)
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("runApprove did not return on the timer")
		}
		if code != exitUsage || !strings.Contains(stderr.String(), "no answer within the request's remaining lifetime (37s)") ||
			!strings.Contains(stderr.String(), "still pending") {
			t.Fatalf("rc %d, want exitUsage, the 37s refusal and 'still pending'\nstderr %s", code, stderr.String())
		}
		if waited != 37*time.Second {
			t.Errorf("waited %v, want exactly 37s", waited)
		}
		if a.row(t).Consumed {
			t.Error("an unanswered request was consumed")
		}
		if _, err := os.Lstat(a.worktree()); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("the worktree exists (err %v)", err)
		}
		if n := len(a.sessionRows(t)); n != 0 {
			t.Errorf("%d session rows, want 0", n)
		}
	})
	t.Run("expired-before-the-prompt", func(t *testing.T) {
		a := newApproveRig(t, "ep1", approveT0)
		term := &promptTerm{in: strings.NewReader("y\n")}
		var stdout, stderr bytes.Buffer
		code := runApprove(a.args(), &stdout, &stderr, approveEnv(term, nil, clockSeq(1, exp-1, exp)))
		if code != exitUsage || !strings.Contains(stderr.String(), "expired before the prompt") {
			t.Fatalf("rc %d, want exitUsage and \"expired before the prompt\"\nstderr %s", code, stderr.String())
		}
		if term.out.Len() != 0 {
			t.Errorf("the prompt was written:\n%s", term.out.String())
		}
	})
}

// CL19: expiry during worktree preparation is refused: the worktree exists only after worktree.Add,
// so only a reading taken after preparation sees the expiry.
func TestApproveRefusesExpiryDuringWorktreePreparation(t *testing.T) {
	a := newApproveRig(t, "ep1", approveT0)
	wt := a.worktree()
	exp := approveT0 + 600
	clock := func() int64 {
		if _, err := os.Stat(wt); err == nil {
			return exp
		}
		return approveT0 + 1
	}
	term := &promptTerm{in: strings.NewReader("y\n")}
	var stdout, stderr bytes.Buffer
	code := runApprove(a.args(), &stdout, &stderr, approveEnv(term, nil, clock))
	if code != exitFatal || !strings.Contains(stderr.String(), "expired before it was consumed") {
		t.Fatalf("rc %d, want exitFatal and \"expired before it was consumed\"\nstderr %s", code, stderr.String())
	}
	if a.row(t).Consumed {
		t.Error("the request was consumed")
	}
	if n := len(a.sessionRows(t)); n != 0 {
		t.Errorf("%d session rows, want 0", n)
	}
	if _, err := os.Lstat(wt); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("the worktree still exists (err %v)", err)
	}
}

// approveFacts is everything an approve did outside the session row (CL5, CL16).
type approveFacts struct {
	WorktreeDir         string // mr.Worktree if it exists, else ""
	WorktreeIsGit       bool   // isGitWorktree(WorktreeDir)
	Branch              string // git branch --show-current in WorktreeDir
	SessionFile         string // mr.SessionFile if it exists, else ""
	SessionMode         os.FileMode
	TokenHash           string // hex(sha256(the session file's bytes))
	MCPFile             string // mr.MCPConfig if it exists, else ""
	MCPAddr, MCPToken   string // from the MCP config's world server URL and Authorization header
	Launch              string // approve's stdout
	ApprovedBy          string // the request row
	RowCredentialID     string // the request row
	SessionCredentialID string // the one session row
}

// collectApproveFacts reads what approve left on disk and in the two stores (CL5).
func collectApproveFacts(t *testing.T, a approveRig, mr approval.MintRequest, stdout string) approveFacts {
	t.Helper()
	f := approveFacts{Launch: stdout}
	if _, err := os.Stat(mr.Worktree); err == nil {
		f.WorktreeDir = mr.Worktree
		f.WorktreeIsGit = isGitWorktree(f.WorktreeDir)
		out, err := exec.Command("git", "-C", f.WorktreeDir, "branch", "--show-current").Output()
		if err != nil {
			t.Fatal(err)
		}
		f.Branch = strings.TrimSpace(string(out))
	}
	if info, err := os.Stat(mr.SessionFile); err == nil {
		f.SessionFile = mr.SessionFile
		f.SessionMode = info.Mode().Perm()
		b, err := os.ReadFile(mr.SessionFile)
		if err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(b)
		f.TokenHash = hex.EncodeToString(sum[:])
	}
	if mr.MCPConfig != "" {
		if _, err := os.Stat(mr.MCPConfig); err == nil {
			f.MCPFile = mr.MCPConfig
			b, err := os.ReadFile(mr.MCPConfig)
			if err != nil {
				t.Fatal(err)
			}
			var cfg struct {
				MCPServers map[string]struct {
					URL     string            `json:"url"`
					Headers map[string]string `json:"headers"`
				} `json:"mcpServers"`
			}
			if err := json.Unmarshal(b, &cfg); err != nil {
				t.Fatal(err)
			}
			srv := cfg.MCPServers[worldMCPServer]
			f.MCPAddr = strings.TrimSuffix(strings.TrimPrefix(srv.URL, "http://"), "/mcp/")
			f.MCPToken = strings.TrimPrefix(srv.Headers["Authorization"], "Bearer ")
		}
	}
	row := a.row(t)
	f.ApprovedBy = row.ApprovedBy
	f.RowCredentialID = row.CredentialID
	if rows := a.sessionRows(t); len(rows) == 1 {
		f.SessionCredentialID = rows[0].CredentialID
	}
	return f
}

// approveEffectProblems returns one entry per mismatch, each starting with its name (CL16).
func approveEffectProblems(mr approval.MintRequest, f approveFacts) []string {
	var out []string
	add := func(name, format string, a ...any) { out = append(out, name+": "+fmt.Sprintf(format, a...)) }
	if f.WorktreeDir != mr.Worktree {
		add("worktree", "the worktree is %q, the request needs %q", f.WorktreeDir, mr.Worktree)
	}
	if !f.WorktreeIsGit {
		add("worktree-git", "%q is not a git worktree", f.WorktreeDir)
	}
	if f.Branch != mr.Branch {
		add("branch", "the branch is %q, the request needs %q", f.Branch, mr.Branch)
	}
	if f.SessionFile != mr.SessionFile {
		add("session_file", "the session file is %q, the request needs %q", f.SessionFile, mr.SessionFile)
	}
	if f.SessionMode != 0o600 {
		add("session_mode", "the session file is %v, want 0600", f.SessionMode)
	}
	if f.TokenHash != f.SessionCredentialID {
		add("token_hash", "the session file's hash is not the session credential")
	}
	if mr.Agent == "codex" {
		if f.MCPFile != "" {
			add("mcp_for_codex", "a codex launch wrote an MCP config %q", f.MCPFile)
		}
	} else if f.MCPFile != mr.MCPConfig || f.MCPAddr != mr.Daemon.Addr || sha256Hex([]byte(f.MCPToken)) != f.SessionCredentialID {
		add("mcp", "the MCP config is %q at %q (token hash mismatch), the request needs %q at %q", f.MCPFile, f.MCPAddr, mr.MCPConfig, mr.Daemon.Addr)
	}
	if mr.Agent == "codex" {
		if !strings.Contains(f.Launch, " codex exec ") {
			add("agent", "the launch line is not a codex exec")
		}
	} else if !strings.Contains(f.Launch, "  claude --tools") {
		add("agent", "the launch line is not a claude --tools launch")
	}
	if mr.Model != "" {
		want := "--model " + shellQuote(mr.Model)
		if mr.Agent == "codex" {
			want = "-m " + shellQuote(mr.Model)
		}
		if !strings.Contains(f.Launch, want) {
			add("model", "the launch line lacks %q", want)
		}
	}
	if mr.Task != "" && !strings.Contains(f.Launch, "task: "+mr.Task) {
		add("task", "the launch line lacks %q", "task: "+mr.Task)
	}
	if f.ApprovedBy != "tty" {
		add("approved_by", "approved by %q, want tty", f.ApprovedBy)
	}
	if f.RowCredentialID != f.SessionCredentialID {
		add("credential_id", "the request row's credential %q is not the session row's %q", f.RowCredentialID, f.SessionCredentialID)
	}
	return out
}

// CL5: a `y` mints exactly the stored request: the session row matches it, and every other effect
// (the worktree, the session file, the MCP config, the launch line) is what the request says.
func TestApproveMintsExactlyTheRequest(t *testing.T) {
	a := newApproveRig(t, "ep1", approveT0)
	term := &promptTerm{in: strings.NewReader("y\n")}
	var stdout, stderr bytes.Buffer
	if code := runApprove(a.args(), &stdout, &stderr, approveEnv(term, nil, atTime(approveT0+1))); code != exitOK {
		t.Fatalf("rc %d\nstderr %s", code, stderr.String())
	}
	row := a.row(t)
	mr, err := approval.Decode(row.Request)
	if err != nil {
		t.Fatal(err)
	}
	rows := a.sessionRows(t)
	if len(rows) != 1 {
		t.Fatalf("%d session rows, want 1", len(rows))
	}
	sess := rows[0]
	sessionBytes, err := os.ReadFile(mr.SessionFile)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(sessionBytes)
	facts := approval.MintFacts{CredentialID: sess.CredentialID, EpisodeID: sess.EpisodeID, GrantsJSON: sess.GrantsJSON,
		ExpiresAt: sess.ExpiresAt, CreatedAt: sess.CreatedAt, Project: sess.Project, ProfileDigest: sess.ProfileDigest,
		TokenHash: hex.EncodeToString(sum[:]), RequestCreatedAt: approveT0, ConsumedAt: row.ConsumedAt}
	if err := approval.MintMatches(mr, facts); err != nil {
		t.Errorf("the minted row is not the request: %v", err)
	}
	if probs := approveEffectProblems(mr, collectApproveFacts(t, a, mr, stdout.String())); len(probs) != 0 {
		t.Errorf("the effects differ from the request:\n%s", strings.Join(probs, "\n"))
	}
}

// CL16: approveEffectProblems passes a clean pair and names each fact that differs (pure).
func TestApproveEffectsCheckCatchesEachMismatch(t *testing.T) {
	const tok = "tok"
	credID := sha256Hex([]byte(tok))
	claude := approval.MintRequest{Episode: "ep1", Worktree: "/ws/ep1", Branch: "ailang/ep1", SessionFile: "/s/ep1.session",
		MCPConfig: "/s/ep1.mcp.json", Agent: "claude", Model: "claude-opus-5-5", Task: "/t/task.md",
		Daemon: approval.DaemonRef{Addr: "127.0.0.1:7777"}}
	cleanClaude := approveFacts{WorktreeDir: "/ws/ep1", WorktreeIsGit: true, Branch: "ailang/ep1", SessionFile: "/s/ep1.session",
		SessionMode: 0o600, TokenHash: credID, MCPFile: "/s/ep1.mcp.json", MCPAddr: "127.0.0.1:7777", MCPToken: tok,
		Launch: "  claude --tools '' --model claude-opus-5-5\ntask: /t/task.md\n", ApprovedBy: "tty",
		RowCredentialID: credID, SessionCredentialID: credID}
	if probs := approveEffectProblems(claude, cleanClaude); len(probs) != 0 {
		t.Fatalf("a clean claude pair has problems: %v", probs)
	}
	codex := approval.MintRequest{Episode: "ep1", Worktree: "/ws/ep1", Branch: "ailang/ep1", SessionFile: "/s/ep1.session",
		Agent: "codex", Model: "gpt-5", Task: "/t/task.md", Daemon: approval.DaemonRef{Addr: "127.0.0.1:7777"}}
	cleanCodex := approveFacts{WorktreeDir: "/ws/ep1", WorktreeIsGit: true, Branch: "ailang/ep1", SessionFile: "/s/ep1.session",
		SessionMode: 0o600, TokenHash: credID, Launch: "  WORLD_SESSION=\"$(cat /s/ep1.session)\" codex exec --ignore-user-config -m gpt-5\ntask: /t/task.md\n",
		ApprovedBy: "tty", RowCredentialID: credID, SessionCredentialID: credID}
	if probs := approveEffectProblems(codex, cleanCodex); len(probs) != 0 {
		t.Fatalf("a clean codex pair has problems: %v", probs)
	}
	type arm struct {
		name   string
		mr     approval.MintRequest
		mutate func(*approveFacts)
	}
	arms := []arm{
		{"worktree", claude, func(f *approveFacts) { f.WorktreeDir = "/elsewhere" }},
		{"worktree-git", claude, func(f *approveFacts) { f.WorktreeIsGit = false }},
		{"branch", claude, func(f *approveFacts) { f.Branch = "other" }},
		{"session_file", claude, func(f *approveFacts) { f.SessionFile = "/s/other.session" }},
		{"session_mode", claude, func(f *approveFacts) { f.SessionMode = 0o644 }},
		{"token_hash", claude, func(f *approveFacts) { f.TokenHash = "x" }},
		{"mcp", claude, func(f *approveFacts) { f.MCPAddr = "127.0.0.1:1" }},
		{"mcp_for_codex", codex, func(f *approveFacts) { f.MCPFile = "/s/ep1.mcp.json" }},
		{"agent", claude, func(f *approveFacts) { f.Launch = "  claude --model claude-opus-5-5\ntask: /t/task.md\n" }},
		{"model", claude, func(f *approveFacts) { f.Launch = "  claude --tools ''\ntask: /t/task.md\n" }},
		{"task", claude, func(f *approveFacts) { f.Launch = "  claude --tools '' --model claude-opus-5-5\n" }},
		{"approved_by", claude, func(f *approveFacts) { f.ApprovedBy = "ssh:x" }},
		{"credential_id", claude, func(f *approveFacts) { f.RowCredentialID = "other" }},
	}
	for _, ar := range arms {
		f := cleanClaude
		if ar.mr.Agent == "codex" {
			f = cleanCodex
		}
		ar.mutate(&f)
		probs := approveEffectProblems(ar.mr, f)
		named := false
		for _, p := range probs {
			if strings.HasPrefix(p, ar.name+":") {
				named = true
			}
		}
		if !named {
			t.Errorf("arm %s: problems %v do not name it", ar.name, probs)
		}
	}
}
