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
