package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

// newEpisodeFixture is an existing store, a workspace root and a one-commit
// git repo to add worktrees from.
type episodeFixture struct {
	db, ws, repo, git string
}

func newEpisodeFixture(t *testing.T) episodeFixture {
	t.Helper()
	git, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git not on PATH")
	}
	base := t.TempDir()
	f := episodeFixture{db: dbPathOf(t), ws: filepath.Join(base, "ws"), repo: filepath.Join(base, "proj"), git: git}
	if err := os.MkdirAll(f.ws, 0o700); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{"init", "-q", f.repo},
		{"-C", f.repo, "-c", "user.email=se@example.invalid", "-c", "user.name=se", "commit", "-q", "--allow-empty", "-m", "init"},
	} {
		cmd := exec.Command(git, args...)
		cmd.Env = worktree.Env(os.Environ())
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	return f
}

func (f episodeFixture) args(ep string, extra ...string) []string {
	return append([]string{ep, "--db", f.db, "--workspace-root", f.ws, "--repo", f.repo}, extra...)
}

func (f episodeFixture) rows(t *testing.T) []store.SessionRow {
	t.Helper()
	return credRows(t, f.db)
}

func (f episodeFixture) worktrees(t *testing.T) int {
	t.Helper()
	cmd := exec.Command(f.git, "-C", f.repo, "worktree", "list", "--porcelain")
	cmd.Env = worktree.Env(os.Environ())
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	return strings.Count(string(out), "worktree ")
}

// assertNothingProvisioned: no worktree dir, no extra git worktree, no row.
func (f episodeFixture) assertNothingProvisioned(t *testing.T, ep string) {
	t.Helper()
	if _, err := os.Lstat(filepath.Join(f.ws, ep)); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("the worktree directory exists (err %v)", err)
	}
	if n := f.worktrees(t); n != 1 {
		t.Errorf("git lists %d worktrees, want 1 (the main one)", n)
	}
	if rows := f.rows(t); len(rows) != 0 {
		t.Errorf("the store holds %d session row(s), want 0", len(rows))
	}
}

func TestSessionNewProvisionsTheWorktreeAndTheSession(t *testing.T) {
	f := newEpisodeFixture(t)
	out := filepath.Join(t.TempDir(), "ep1.session")
	now := time.Now().Unix()
	var stdout, stderr bytes.Buffer
	code := runSessionNew(f.args("ep1", "--preset", "se-tools", "--ttl", "14400", "--out", out), &stdout, &stderr, testMintEnv("y\n", now))
	if code != exitOK {
		t.Fatalf("session new rc %d\n%s%s", code, stdout.String(), stderr.String())
	}
	if !isGitWorktree(filepath.Join(f.ws, "ep1")) || f.worktrees(t) != 2 {
		t.Fatalf("no git worktree at %s", filepath.Join(f.ws, "ep1"))
	}
	tok, err := os.ReadFile(out)
	if err != nil || !isHex64(string(tok)) {
		t.Fatalf("token file: %q %v", tok, err)
	}
	if info, _ := os.Stat(out); info.Mode().Perm() != 0o600 {
		t.Fatalf("token file mode %v", info.Mode().Perm())
	}
	if strings.Contains(stdout.String()+stderr.String(), string(tok)) {
		t.Fatal("the token was printed (MUT-TOKEN-ECHO)")
	}
	rows := f.rows(t)
	if len(rows) != 1 || rows[0].EpisodeID != "ep1" || rows[0].ExpiresAt != now+14400 {
		t.Fatalf("rows = %+v", rows)
	}
	if !strings.Contains(stderr.String(), "credential_id="+rows[0].CredentialID) {
		t.Fatalf("stderr lacks the credential_id: %s", stderr.String())
	}
	var grants []broker.Capability
	_ = json.Unmarshal([]byte(rows[0].GrantsJSON), &grants)
	if len(grants) != 10 {
		t.Fatalf("%d grants, want the ten se-tools grants", len(grants))
	}
	for _, g := range grants {
		if g.Scope != "worktree" || g.Budget != 50 || g.ExpiresAt != now+14400 {
			t.Fatalf("grant %+v, want scope worktree, budget 50, expiring with the session", g)
		}
	}

	// A second session for the same episode reuses the worktree: no --repo.
	stdout.Reset()
	stderr.Reset()
	code = runSessionNew([]string{"ep1", "--db", f.db, "--workspace-root", f.ws, "--grant", "Workspace.Read=worktree:3"}, &stdout, &stderr, testMintEnv("y\n", now))
	if code != exitOK || f.worktrees(t) != 2 || len(f.rows(t)) != 2 {
		t.Fatalf("reuse: rc %d worktrees %d rows %d\n%s", code, f.worktrees(t), len(f.rows(t)), stderr.String())
	}
}

// AC5.1: with no controlling terminal there is no worktree and no row.
// Kills MUT-FENCE-BYPASS.
func TestSessionNewWithoutATerminalProvisionsNothing(t *testing.T) {
	f := newEpisodeFixture(t)
	var stdout, stderr bytes.Buffer
	code := runSessionNew(f.args("ep1", "--preset", "se-tools"), &stdout, &stderr, noTerminalEnv(1000))
	if code != exitUsage || !strings.Contains(stderr.String(), "no controlling terminal") {
		t.Fatalf("rc %d stderr %q", code, stderr.String())
	}
	f.assertNothingProvisioned(t, "ep1")
}

// AC5.2: answering n provisions neither.
func TestSessionNewDeclinedProvisionsNothing(t *testing.T) {
	f := newEpisodeFixture(t)
	var stdout, stderr bytes.Buffer
	code := runSessionNew(f.args("ep1", "--preset", "se-tools"), &stdout, &stderr, testMintEnv("n\n", 1000))
	if code != exitUsage || !strings.Contains(stderr.String(), "aborted") {
		t.Fatalf("rc %d stderr %q", code, stderr.String())
	}
	f.assertNothingProvisioned(t, "ep1")
}

// AC5.3: a mint failure removes the worktree this run created.
func TestSessionNewMintFailureRollsBackTheWorktree(t *testing.T) {
	f := newEpisodeFixture(t)
	old := sessionMint
	defer func() { sessionMint = old }()
	sessionMint = func(context.Context, *store.Credentials, string, []broker.Capability, int64, int64, io.Writer) (string, string, store.SessionRow, error) {
		return "", "", store.SessionRow{}, errors.New("injected mint failure")
	}
	var stdout, stderr bytes.Buffer
	code := runSessionNew(f.args("ep1", "--preset", "se-tools"), &stdout, &stderr, testMintEnv("y\n", 1000))
	if code != exitFatal || !strings.Contains(stderr.String(), "removed the worktree") {
		t.Fatalf("rc %d stderr %q", code, stderr.String())
	}
	f.assertNothingProvisioned(t, "ep1")
}

// Row 189 (D-WORLD-85 C): with the WORLD writer lock held — what a serving
// daemon holds — session new provisions the worktree and the session, and
// never tells the operator to stop the daemon. Mutant: point the session verbs
// back at the world store's lock and this goes red with WriterAlreadyActive.
func TestSessionNewWithTheDaemonRunningProvisions(t *testing.T) {
	f := newEpisodeFixture(t)
	w, err := store.Open(f.db)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = w.Close() }()
	var stdout, stderr bytes.Buffer
	code := runSessionNew(f.args("ep1", "--preset", "se-tools"), &stdout, &stderr, testMintEnv("y\n", 1000))
	if code != exitOK || strings.Contains(stderr.String(), "stop the daemon") {
		t.Fatalf("session new beside a live world writer: rc %d stderr %q", code, stderr.String())
	}
	if !isGitWorktree(filepath.Join(f.ws, "ep1")) || len(f.rows(t)) != 1 {
		t.Fatalf("worktree %v rows %d", isGitWorktree(filepath.Join(f.ws, "ep1")), len(f.rows(t)))
	}
	// The world writer lock still protects everything else.
	if s, err := store.Open(f.db); err == nil {
		_ = s.Close()
		t.Fatal("a second world writer opened beside the first")
	}
}

// A held CREDENTIALS writer lock (another session verb mid-write) refuses
// before the worktree exists: the lock is real, it just is not the daemon's.
func TestSessionNewWithACredentialsWriterRunningProvisionsNothing(t *testing.T) {
	f := newEpisodeFixture(t)
	_ = credsFor(t, f.db)
	var stdout, stderr bytes.Buffer
	code := runSessionNew(f.args("ep1", "--preset", "se-tools"), &stdout, &stderr, testMintEnv("y\n", 1000))
	if code != exitFatal || !strings.Contains(stderr.String(), "another session new/mint/revoke") {
		t.Fatalf("rc %d stderr %q", code, stderr.String())
	}
	f.assertNothingProvisioned(t, "ep1")
}

func TestSessionNewValidatesBeforeTheFence(t *testing.T) {
	f := newEpisodeFixture(t)
	missing := filepath.Join(t.TempDir(), "world.db")
	for name, tc := range map[string]struct {
		args []string
		want string
	}{
		"bad-episode": {f.args("Bad_Ep", "--preset", "se-tools"), "does not match"},
		"no-grants":   {f.args("ep1"), "no grants"},
		"bad-preset":  {f.args("ep1", "--preset", "nope"), "unknown --preset"},
		"missing-db":  {[]string{"ep1", "--db", missing, "--workspace-root", f.ws, "--repo", f.repo, "--preset", "se-tools"}, "refusing to create one"},
		"no-repo":     {[]string{"ep1", "--db", f.db, "--workspace-root", f.ws, "--preset", "se-tools"}, "--repo is required"},
	} {
		t.Run(name, func(t *testing.T) {
			opened := false
			env := sessionEnv{openTerminal: func() (io.ReadWriteCloser, error) { opened = true; return nil, os.ErrNotExist }, now: func() int64 { return 0 }}
			var stdout, stderr bytes.Buffer
			if code := runSessionNew(tc.args, &stdout, &stderr, env); code != exitUsage || !strings.Contains(stderr.String(), tc.want) {
				t.Fatalf("rc %d stderr %q, want %q", code, stderr.String(), tc.want)
			}
			if opened {
				t.Fatal("validation ran after the fence")
			}
			if _, err := os.Stat(missing); !errors.Is(err, os.ErrNotExist) {
				t.Fatal("a missing --db was created")
			}
		})
	}
}

func TestSessionNewRefusesAWorkspaceRootHoldingTheStore(t *testing.T) {
	f := newEpisodeFixture(t)
	inner := filepath.Join(f.ws, "store", "world.db")
	_ = os.MkdirAll(filepath.Dir(inner), 0o700)
	st, err := store.Open(inner)
	if err != nil {
		t.Fatal(err)
	}
	_ = st.Close()
	var stdout, stderr bytes.Buffer
	code := runSessionNew([]string{"ep1", "--db", inner, "--workspace-root", f.ws, "--repo", f.repo, "--preset", "se-tools"}, &stdout, &stderr, testMintEnv("y\n", 1))
	if code != exitUsage || !strings.Contains(stderr.String(), "contains the store directory") {
		t.Fatalf("rc %d stderr %q", code, stderr.String())
	}
}

// AC5.4: the se-tools preset is exactly the union of the (effect, scope)
// pairs the se-tools transitions declare. Kills MUT-PRESET-DRIFT.
func TestSessionNewPresetEqualsDeclaredEffects(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "packages", "se-tools", "transitions.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest []struct {
		DeclaredEffects []struct {
			Effect string `json:"effect"`
			Scope  string `json:"scope"`
		} `json:"declaredEffects"`
	}
	if err := json.Unmarshal(raw, &manifest); err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{}
	for _, d := range manifest {
		for _, e := range d.DeclaredEffects {
			want[e.Effect+"@"+e.Scope] = true
		}
	}
	if len(manifest) == 0 || len(want) == 0 {
		t.Fatal("instrument failure: the manifest declares no effects")
	}
	got := map[string]bool{}
	for _, g := range sePresetGrants {
		if got[g.Effect+"@"+g.Scope] {
			t.Fatalf("preset repeats %s@%s", g.Effect, g.Scope)
		}
		got[g.Effect+"@"+g.Scope] = true
	}
	keys := func(m map[string]bool) string {
		var ks []string
		for k := range m {
			ks = append(ks, k)
		}
		sort.Strings(ks)
		return strings.Join(ks, " ")
	}
	if keys(got) != keys(want) {
		t.Fatalf("preset %s\nmanifest %s", keys(got), keys(want))
	}
	if len(got) != 10 {
		t.Fatalf("preset has %d grants, want 10", len(got))
	}
}

// AC5.6: list works while a writer holds the store, and shows hashes only.
func TestSessionListBesideALiveWriter(t *testing.T) {
	f := newEpisodeFixture(t)
	now := time.Now().Unix()
	var stdout, stderr bytes.Buffer
	if code := runSessionNew(f.args("ep1", "--preset", "se-tools"), &stdout, &stderr, testMintEnv("y\n", now)); code != exitOK {
		t.Fatalf("session new: %d %s", code, stderr.String())
	}
	token := strings.TrimSpace(stdout.String())
	if !isHex64(token) {
		t.Fatalf("stdout mode should print the token once, got %q", stdout.String())
	}
	w, err := store.Open(f.db)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = w.Close() }()
	var out, errw bytes.Buffer
	if code := runSessionList([]string{"--db", f.db}, &out, &errw, testMintEnv("", now)); code != exitOK {
		t.Fatalf("list beside a live writer: rc %d %s", code, errw.String())
	}
	rows := f.rows(t)
	for _, want := range []string{rows[0].CredentialID, "ep1", "live", "Workspace.Read=worktree:50", "1 session credential(s)"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("list lacks %q:\n%s", want, out.String())
		}
	}
	if strings.Contains(out.String(), token) {
		t.Fatal("list printed the token")
	}
	out.Reset()
	if code := runSessionList([]string{"--db", f.db, "--json", "--episode", "ep1"}, &out, &errw, testMintEnv("", now+999999)); code != exitOK {
		t.Fatal(errw.String())
	}
	var listing []sessionListing
	if err := json.Unmarshal(out.Bytes(), &listing); err != nil || len(listing) != 1 || listing[0].Live || len(listing[0].Grants) != len(sePresetGrants) {
		t.Fatalf("json listing %+v %v", listing, err)
	}
}

// Row 185: an existing episode directory that is not a git worktree is
// refused, and the refusal names --allow-non-git; nothing is minted and no
// terminal is opened. Kills MUT-ADOPT-NON-GIT (dropping the isGitWorktree
// arm) and MUT-REFUSAL-NO-FLAG (a refusal that does not name the flag).
func TestSessionNewRefusesANonGitEpisodeNamingTheFlag(t *testing.T) {
	f := newEpisodeFixture(t)
	dir := filepath.Join(f.ws, "ep1")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	opened := false
	env := sessionEnv{openTerminal: func() (io.ReadWriteCloser, error) { opened = true; return nil, os.ErrNotExist }, now: func() int64 { return 0 }}
	var stdout, stderr bytes.Buffer
	code := runSessionNew([]string{"ep1", "--db", f.db, "--workspace-root", f.ws, "--preset", "se-tools"}, &stdout, &stderr, env)
	if code != exitUsage || !strings.Contains(stderr.String(), "is not a git worktree") || !strings.Contains(stderr.String(), "--allow-non-git") {
		t.Fatalf("rc %d stderr %q, want a refusal naming --allow-non-git", code, stderr.String())
	}
	if opened {
		t.Fatal("the refusal ran after the terminal fence")
	}
	if rows := f.rows(t); len(rows) != 0 {
		t.Fatalf("%d session row(s) minted for a refused episode", len(rows))
	}
}

// Row 185: with --allow-non-git the directory is adopted (no git worktree is
// made), the prompt and stderr say DEGRADED, and session list marks it —
// in text and as "degraded": true in --json — while a git episode in the same
// store stays unmarked. Kills MUT-DEGRADED-SILENT (no marker on adoption).
func TestSessionNewAllowNonGitMintsADegradedEpisode(t *testing.T) {
	f := newEpisodeFixture(t)
	dir := filepath.Join(f.ws, "plain")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	now := time.Now().Unix()
	term := &promptTerm{in: strings.NewReader("y\n")}
	env := sessionEnv{openTerminal: func() (io.ReadWriteCloser, error) { return term, nil }, now: func() int64 { return now }}
	var stdout, stderr bytes.Buffer
	code := runSessionNew([]string{"plain", "--db", f.db, "--workspace-root", f.ws, "--grant", "Workspace.Read=worktree:3", "--allow-non-git"},
		&stdout, &stderr, env)
	if code != exitOK {
		t.Fatalf("--allow-non-git: rc %d\n%s", code, stderr.String())
	}
	if !strings.Contains(term.out.String(), "NON-GIT") || !strings.Contains(term.out.String(), "DEGRADED") {
		t.Fatalf("the prompt does not say the episode is degraded: %q", term.out.String())
	}
	if !strings.Contains(stderr.String(), "episode plain is DEGRADED") {
		t.Fatalf("stderr lacks the degraded notice: %s", stderr.String())
	}
	if isGitWorktree(dir) || f.worktrees(t) != 1 {
		t.Fatal("--allow-non-git turned the directory into a git worktree")
	}
	if rows := f.rows(t); len(rows) != 1 || rows[0].EpisodeID != "plain" {
		t.Fatalf("rows = %+v", rows)
	}

	// A git episode beside it.
	stdout.Reset()
	stderr.Reset()
	if code := runSessionNew(f.args("ep1", "--grant", "Workspace.Read=worktree:3"), &stdout, &stderr, testMintEnv("y\n", now)); code != exitOK {
		t.Fatalf("git episode: rc %d %s", code, stderr.String())
	}
	if strings.Contains(stderr.String(), "DEGRADED") {
		t.Fatalf("a git episode was reported degraded: %s", stderr.String())
	}

	var out, errw bytes.Buffer
	if code := runSessionList([]string{"--db", f.db, "--workspace-root", f.ws}, &out, &errw, testMintEnv("", now)); code != exitOK {
		t.Fatalf("list: rc %d %s", code, errw.String())
	}
	for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		switch {
		case strings.Contains(line, " plain "):
			if !strings.Contains(line, "DEGRADED (not a git worktree)") {
				t.Errorf("the degraded episode is unmarked: %q", line)
			}
		case strings.Contains(line, " ep1 "):
			if strings.Contains(line, "DEGRADED") {
				t.Errorf("the git episode is marked degraded: %q", line)
			}
		}
	}
	out.Reset()
	if code := runSessionList([]string{"--db", f.db, "--workspace-root", f.ws, "--json"}, &out, &errw, testMintEnv("", now)); code != exitOK {
		t.Fatal(errw.String())
	}
	var listing []sessionListing
	if err := json.Unmarshal(out.Bytes(), &listing); err != nil || len(listing) != 2 {
		t.Fatalf("json listing %s %v", out.String(), err)
	}
	for _, l := range listing {
		if l.Degraded != (l.Episode == "plain") {
			t.Errorf("episode %s degraded=%v", l.Episode, l.Degraded)
		}
	}
}

// Row 185: the flag changes nothing for a git episode — a new worktree is
// still created from --repo and nothing is marked degraded.
func TestSessionNewAllowNonGitLeavesGitEpisodesUnchanged(t *testing.T) {
	f := newEpisodeFixture(t)
	var stdout, stderr bytes.Buffer
	code := runSessionNew(f.args("ep1", "--grant", "Workspace.Read=worktree:1", "--allow-non-git"), &stdout, &stderr, testMintEnv("y\n", 1000))
	if code != exitOK || !isGitWorktree(filepath.Join(f.ws, "ep1")) || f.worktrees(t) != 2 {
		t.Fatalf("rc %d worktrees %d\n%s", code, f.worktrees(t), stderr.String())
	}
	if strings.Contains(stderr.String(), "DEGRADED") || episodeDegraded(f.ws, "ep1") {
		t.Fatalf("a git episode is degraded: %s", stderr.String())
	}
}

// promptTerm is fakeTerm that keeps the prompt: reads come from in, writes
// land in out.
type promptTerm struct {
	in  io.Reader
	out bytes.Buffer
}

func (t *promptTerm) Read(p []byte) (int, error)  { return t.in.Read(p) }
func (t *promptTerm) Write(p []byte) (int, error) { return t.out.Write(p) }
func (t *promptTerm) Close() error                { return nil }

func TestSessionRevokeByEpisode(t *testing.T) {
	f := newEpisodeFixture(t)
	for i := 0; i < 2; i++ {
		var stdout, stderr bytes.Buffer
		if code := runSessionNew(f.args("ep1", "--grant", "Workspace.Read=worktree:1"), &stdout, &stderr, testMintEnv("y\n", 1000)); code != exitOK {
			t.Fatalf("new %d: %s", i, stderr.String())
		}
		f.repo = "" // the second run reuses the worktree
	}
	var out, errw bytes.Buffer
	if code := runSessionRevoke([]string{"--episode", "ep1", "--db", f.db}, &out, &errw); code != exitOK || strings.Count(out.String(), "revoked") != 2 {
		t.Fatalf("revoke --episode: %d %q %q", code, out.String(), errw.String())
	}
	if rows := f.rows(t); len(rows) != 0 {
		t.Fatalf("%d rows left", len(rows))
	}
	out.Reset()
	errw.Reset()
	if code := runSessionRevoke([]string{"--episode", "ep1", "--db", f.db}, &out, &errw); code != exitUsage {
		t.Fatalf("revoke --episode of an episode with no credentials: %d", code)
	}
}
