package main

// Row 189 M2/M3 (w-world-onboarding §3.1, AC2-AC4, D-189-2 A): episode start
// and session new --profile. The daemon check is driven through the
// fetchDaemonHealth and episodeSelfBuild seams here; episode_loopback_test.go
// drives the real HTTP fetch and the real /dev/tty.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

// testBuild is the build identity the stubbed daemon and this binary share.
var testBuild = daemon.BuildIdentity{VCSRevision: "0123456789abcdef0123456789abcdef01234567"}

// writeTestProfile writes <dir>/<project>.json, a minimal profile the
// daemon's own loader accepts; extra lands inside the object, so a different
// extra is different bytes (a different digest).
func writeTestProfile(t *testing.T, dir, project, extra string) *broker.ExecProfile {
	t.Helper()
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	data := []byte(fmt.Sprintf(`{"profile": "world/exec-profile/v1", "project": %q, "path": ["/usr/bin", "/bin"], "timeout_ms": 9000, `+
		`"commands": {"true": {"argv": ["sh", "-c", "exit 0"]}}%s}`, project, extra))
	if err := os.WriteFile(filepath.Join(dir, project+".json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	p, err := broker.ParseExecProfile(data)
	if err != nil {
		t.Fatalf("test profile does not parse: %v", err)
	}
	return p
}

// matchingHealth is the /v1/health of a daemon episode start accepts for f.
func matchingHealth(t *testing.T, f episodeFixture, profilesDir string, profiles ...*broker.ExecProfile) daemon.HealthResponse {
	t.Helper()
	ref, _, err := pinnedToolRef()
	if err != nil {
		t.Skipf("SKIP %s: %v", t.Name(), err)
	}
	cred, err := credentialsFor(f.db, "")
	if err != nil {
		t.Fatal(err)
	}
	h := daemon.HealthResponse{Status: "ok", Build: testBuild, ToolRef: ref, CredentialsDB: cred, ExecProfiles: []daemon.HealthExecProfile{}}
	for _, p := range profiles {
		h.ExecProfiles = append(h.ExecProfiles, daemon.HealthExecProfile{Project: p.Project, Digest: p.Digest,
			File: filepath.Join(profilesDir, p.Project+".json")})
	}
	return h
}

// stubDaemon answers the daemon check with h (or err) and makes this binary
// testBuild.
func stubDaemon(t *testing.T, h daemon.HealthResponse, err error) {
	t.Helper()
	oldFetch, oldSelf := fetchDaemonHealth, episodeSelfBuild
	t.Cleanup(func() { fetchDaemonHealth, episodeSelfBuild = oldFetch, oldSelf })
	fetchDaemonHealth = func(string) (daemon.HealthResponse, error) { return h, err }
	episodeSelfBuild = func() daemon.BuildIdentity { return testBuild }
}

type episodeRun struct {
	f        episodeFixture
	sessions string
	profiles string
}

func newEpisodeRun(t *testing.T) episodeRun {
	base := t.TempDir()
	return episodeRun{f: newEpisodeFixture(t), sessions: filepath.Join(base, "sessions"), profiles: filepath.Join(base, "profiles")}
}

func (r episodeRun) args(ep string, extra ...string) []string {
	return append([]string{ep, "--db", r.f.db, "--workspace-root", r.f.ws, "--repo", r.f.repo,
		"--sessions-dir", r.sessions, "--exec-profiles-dir", r.profiles}, extra...)
}

// spyEnv records whether the terminal was opened, then has none.
func spyEnv(opened *bool) sessionEnv {
	return sessionEnv{openTerminal: func() (io.ReadWriteCloser, error) { *opened = true; return nil, os.ErrNotExist },
		now: func() int64 { return 1000 }}
}

// assertNoSessionFiles: neither the sessions dir nor its files exist.
func (r episodeRun) assertNoSessionFiles(t *testing.T) {
	t.Helper()
	if _, err := os.Lstat(r.sessions); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("the sessions directory %s exists (err %v)", r.sessions, err)
	}
}

// AC4 and the happy path: worktree on world/<ep>, a bound credential, the
// session file and the MCP config at 0600 in a 0700 dir, the launch line —
// and the token in no other file and on neither stream.
func TestEpisodeStartProvisionsAndPrintsTheLaunchWithTheTokenInTwoFilesOnly(t *testing.T) {
	r := newEpisodeRun(t)
	p := writeTestProfile(t, r.profiles, "proj", "")
	stubDaemon(t, matchingHealth(t, r.f, r.profiles, p), nil)
	task := filepath.Join(t.TempDir(), "task.md")
	if err := os.WriteFile(task, []byte("do the thing\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	now := time.Now().Unix()
	term := &promptTerm{in: strings.NewReader("y\n")}
	env := sessionEnv{openTerminal: func() (io.ReadWriteCloser, error) { return term, nil }, now: func() int64 { return now }}
	var stdout, stderr bytes.Buffer
	code := runEpisodeStart(r.args("ep1", "--profile", "proj", "--task", task, "--model", "claude-opus-5-5"), &stdout, &stderr, env)
	if code != exitOK {
		t.Fatalf("episode start rc %d\nstdout %s\nstderr %s", code, stdout.String(), stderr.String())
	}
	sessionFile, mcpFile := filepath.Join(r.sessions, "ep1.session"), filepath.Join(r.sessions, "ep1.mcp.json")
	tok, err := os.ReadFile(sessionFile)
	if err != nil || !isHex64(string(tok)) {
		t.Fatalf("session file: %q %v", tok, err)
	}
	token := string(tok)
	for path, want := range map[string]os.FileMode{r.sessions: 0o700, sessionFile: 0o600, mcpFile: 0o600} {
		if info, err := os.Stat(path); err != nil || info.Mode().Perm() != want {
			t.Errorf("%s: mode %v (%v), want %v", path, info.Mode().Perm(), err, want)
		}
	}
	var mcp struct {
		MCPServers map[string]struct {
			Type, URL string
			Headers   map[string]string
		} `json:"mcpServers"`
	}
	raw, _ := os.ReadFile(mcpFile)
	if err := json.Unmarshal(raw, &mcp); err != nil || mcp.MCPServers["world"].URL != "http://"+defaultEpisodeAddr+"/mcp/" ||
		mcp.MCPServers["world"].Headers["Authorization"] != "Bearer "+token || mcp.MCPServers["world"].Type != "http" {
		t.Fatalf("MCP config %s: %v", raw, err)
	}

	// The credential is bound, and the y/N named the digest and the grant.
	rows := r.f.rows(t)
	if len(rows) != 1 || rows[0].Project != "proj" || rows[0].ProfileDigest != p.Digest {
		t.Fatalf("rows = %+v, want one bound to proj %s", rows, p.Digest)
	}
	for _, want := range []string{"exec profile proj " + p.Digest, "Workspace.Exec=worktree:50", "on new branch world/ep1"} {
		if !strings.Contains(term.out.String(), want) {
			t.Errorf("the y/N lacks %q: %s", want, term.out.String())
		}
	}
	branch := exec.Command(r.f.git, "-C", filepath.Join(r.f.ws, "ep1"), "rev-parse", "--abbrev-ref", "HEAD")
	branch.Env = worktree.Env(os.Environ())
	if out, err := branch.Output(); err != nil || strings.TrimSpace(string(out)) != "world/ep1" {
		t.Fatalf("worktree branch %q %v, want world/ep1", out, err)
	}
	for _, want := range []string{"claude --tools '' --strict-mcp-config --mcp-config " + mcpFile, "--allowedTools ",
		"mcp__world__workspace-exec", "mcp__world__ailang-read", "--model claude-opus-5-5", "task: " + task} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("the launch lacks %q:\n%s", want, stdout.String())
		}
	}

	// AC4: the token is on neither stream, not on the terminal, and in no
	// file of the worktree, the sessions dir or the store dir but the two.
	for name, s := range map[string]string{"stdout": stdout.String(), "stderr": stderr.String(), "terminal": term.out.String()} {
		if strings.Contains(s, token) {
			t.Errorf("the token is on %s", name)
		}
	}
	var holders []string
	for _, root := range []string{filepath.Join(r.f.ws, "ep1"), r.sessions, filepath.Dir(r.f.db)} {
		_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			if b, err := os.ReadFile(path); err == nil && bytes.Contains(b, []byte(token)) {
				holders = append(holders, path)
			}
			return nil
		})
	}
	if len(holders) != 2 || !contains(holders, sessionFile) || !contains(holders, mcpFile) {
		t.Fatalf("the token is in %v, want exactly the session file and the MCP config", holders)
	}
	// Instrument control: the grep finds a token where one is.
	if b, _ := os.ReadFile(sessionFile); !bytes.Contains(b, []byte(token)) {
		t.Fatal("instrument failure: the grep cannot see the token in the session file")
	}
}

// AC2 (MUT-FENCE-BYPASS kept): no controlling terminal -> no worktree, no
// row, no sessions dir, no file — after the daemon check passed.
func TestEpisodeStartWithoutATerminalProvisionsNothing(t *testing.T) {
	r := newEpisodeRun(t)
	p := writeTestProfile(t, r.profiles, "proj", "")
	stubDaemon(t, matchingHealth(t, r.f, r.profiles, p), nil)
	var stdout, stderr bytes.Buffer
	code := runEpisodeStart(r.args("ep1", "--profile", "proj"), &stdout, &stderr, noTerminalEnv(1000))
	if code != exitUsage || !strings.Contains(stderr.String(), "refusing: no controlling terminal") {
		t.Fatalf("rc %d stderr %q", code, stderr.String())
	}
	r.f.assertNothingProvisioned(t, "ep1")
	r.assertNoSessionFiles(t)
}

// AC3: a daemon that is not this build, not this tool pin, not this
// credentials store, or without the profile's current bytes is refused
// before the terminal opens, naming the fix. Control: the matching daemon
// gets as far as the terminal.
func TestEpisodeStartRefusesAMismatchedDaemonBeforeTheTerminal(t *testing.T) {
	r := newEpisodeRun(t)
	p := writeTestProfile(t, r.profiles, "proj", "")
	good := matchingHealth(t, r.f, r.profiles, p)
	mutate := func(m func(*daemon.HealthResponse)) daemon.HealthResponse {
		h := good
		h.ExecProfiles = append([]daemon.HealthExecProfile(nil), good.ExecProfiles...)
		m(&h)
		return h
	}
	for name, tc := range map[string]struct {
		h    daemon.HealthResponse
		want string
	}{
		"older build, no identity": {daemon.HealthResponse{Status: "ok", DBPath: r.f.db}, "Rebuild it from this checkout"},
		"other revision":           {mutate(func(h *daemon.HealthResponse) { h.Build.VCSRevision = "feedface" }), "Rebuild it from this checkout"},
		"dirty build":              {mutate(func(h *daemon.HealthResponse) { h.Build.VCSModified = true }), "Rebuild"},
		"other tool binary":        {mutate(func(h *daemon.HealthResponse) { h.ToolRef = "sha256:" + strings.Repeat("0", 64) }), "--tool-ailang-bin"},
		"no tool binary":           {mutate(func(h *daemon.HealthResponse) { h.ToolRef = "" }), "serves the tool binary none"},
		"other credentials store":  {mutate(func(h *daemon.HealthResponse) { h.CredentialsDB = "/elsewhere/credentials.db" }), "which it never reads"},
		"profile not loaded":       {mutate(func(h *daemon.HealthResponse) { h.ExecProfiles = nil }), "has not loaded exec profile proj"},
		"profile changed":          {mutate(func(h *daemon.HealthResponse) { h.ExecProfiles[0].Digest = "sha256:" + strings.Repeat("e", 64) }), "restart serve so it loads"},
		"profile left out": {mutate(func(h *daemon.HealthResponse) {
			h.ExecProfiles = nil
			h.ExecProfilesExcluded = []daemon.HealthExecExclusion{{Project: "proj", File: "/p/proj.json", Reason: "arm7-toolchain: exit 3"}}
		}), "arm7-toolchain: exit 3"},
	} {
		t.Run(name, func(t *testing.T) {
			stubDaemon(t, tc.h, nil)
			opened := false
			var stdout, stderr bytes.Buffer
			code := runEpisodeStart(r.args("ep1", "--profile", "proj"), &stdout, &stderr, spyEnv(&opened))
			if code != exitUsage || !strings.Contains(stderr.String(), tc.want) || !strings.Contains(stderr.String(), "refusing") {
				t.Fatalf("rc %d stderr %q, want a refusal naming %q", code, stderr.String(), tc.want)
			}
			if opened {
				t.Fatal("the daemon check ran after the terminal fence")
			}
			r.f.assertNothingProvisioned(t, "ep1")
			r.assertNoSessionFiles(t)
		})
	}
	// Control: the matching daemon passes the check and reaches the fence.
	stubDaemon(t, good, nil)
	opened := false
	var stdout, stderr bytes.Buffer
	if code := runEpisodeStart(r.args("ep1", "--profile", "proj"), &stdout, &stderr, spyEnv(&opened)); !opened || code != exitUsage {
		t.Fatalf("the matching daemon did not reach the terminal: rc %d %s", code, stderr.String())
	}
	// The build rule: the same executable bytes match, stamped or not (a
	// build in a git worktree carries no vcs.revision: the grader's case);
	// other bytes match only on one clean revision.
	exe := func(n string) string { return "sha256:" + strings.Repeat(n, 64) }
	for name, tc := range map[string]struct {
		self, daemon daemon.BuildIdentity
		match        bool
	}{
		"same bytes, no stamp":        {daemon.BuildIdentity{ExecutableSHA256: exe("a")}, daemon.BuildIdentity{ExecutableSHA256: exe("a")}, true},
		"other bytes, same clean rev": {daemon.BuildIdentity{VCSRevision: "abc", ExecutableSHA256: exe("a")}, daemon.BuildIdentity{VCSRevision: "abc", ExecutableSHA256: exe("b")}, true},
		"other bytes, no stamps":      {daemon.BuildIdentity{ExecutableSHA256: exe("a")}, daemon.BuildIdentity{ExecutableSHA256: exe("b")}, false},
		"other bytes, dirty revision": {daemon.BuildIdentity{VCSRevision: "abc", VCSModified: true, ExecutableSHA256: exe("a")}, daemon.BuildIdentity{VCSRevision: "abc", VCSModified: true, ExecutableSHA256: exe("b")}, false},
		"no identity here (go run)":   {daemon.BuildIdentity{}, testBuild, false},
	} {
		t.Run("build "+name, func(t *testing.T) {
			h := good
			h.Build = tc.daemon
			stubDaemon(t, h, nil)
			episodeSelfBuild = func() daemon.BuildIdentity { return tc.self }
			opened := false
			var stdout, stderr bytes.Buffer
			code := runEpisodeStart(r.args("ep1", "--profile", "proj"), &stdout, &stderr, spyEnv(&opened))
			if tc.match {
				if !opened || strings.Contains(stderr.String(), "Rebuild it from this checkout") {
					t.Fatalf("a matching build did not reach the terminal: rc %d %s", code, stderr.String())
				}
				return
			}
			if code != exitUsage || opened || !strings.Contains(stderr.String(), "Rebuild it from this checkout") {
				t.Fatalf("a mismatched build: rc %d opened %v %s, want the rebuild refusal", code, opened, stderr.String())
			}
			if !tc.self.Known() && !strings.Contains(stderr.String(), "only the same binary file matches") {
				t.Fatalf("the refusal does not say why: %s", stderr.String())
			}
			r.f.assertNothingProvisioned(t, "ep1")
			r.assertNoSessionFiles(t)
		})
	}
}

// AC3, no daemon: refused before the terminal, printing the serve line.
func TestEpisodeStartRefusesWhenNoDaemonServes(t *testing.T) {
	r := newEpisodeRun(t)
	stubDaemon(t, daemon.HealthResponse{}, &daemonUnreachableError{addr: defaultEpisodeAddr, err: errors.New("connection refused")})
	opened := false
	var stdout, stderr bytes.Buffer
	code := runEpisodeStart(r.args("ep1"), &stdout, &stderr, spyEnv(&opened))
	cred, _ := credentialsFor(r.f.db, "")
	for _, want := range []string{"no daemon is serving on " + defaultEpisodeAddr, "never starts one",
		" serve --db " + r.f.db, "--credentials-db " + cred, "--bind " + defaultEpisodeAddr, "--workspace-root ", "--tool-ailang-bin "} {
		if !strings.Contains(stderr.String(), want) {
			t.Errorf("refusal lacks %q:\n%s", want, stderr.String())
		}
	}
	if code != exitUsage || opened {
		t.Fatalf("rc %d, terminal opened %v", code, opened)
	}
	r.f.assertNothingProvisioned(t, "ep1")
	r.assertNoSessionFiles(t)
}

// Step 3: an existing session file (or MCP config) refuses before the
// terminal, and is left as it was. Control: removing it lets the run reach
// the terminal.
func TestEpisodeStartRefusesAnExistingSessionFileBeforeTheTerminal(t *testing.T) {
	r := newEpisodeRun(t)
	stubDaemon(t, matchingHealth(t, r.f, r.profiles), nil)
	if err := os.MkdirAll(r.sessions, 0o700); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"ep1.session", "ep1.mcp.json"} {
		path := filepath.Join(r.sessions, name)
		if err := os.WriteFile(path, []byte("earlier"), 0o600); err != nil {
			t.Fatal(err)
		}
		opened := false
		var stdout, stderr bytes.Buffer
		code := runEpisodeStart(r.args("ep1"), &stdout, &stderr, spyEnv(&opened))
		if code != exitUsage || opened || !strings.Contains(stderr.String(), path+" already exists") {
			t.Fatalf("%s: rc %d opened %v stderr %q", name, code, opened, stderr.String())
		}
		if b, _ := os.ReadFile(path); string(b) != "earlier" {
			t.Fatalf("%s was changed: %q", name, b)
		}
		r.f.assertNothingProvisioned(t, "ep1")
		_ = os.Remove(path)
	}
	opened := false
	var stdout, stderr bytes.Buffer
	_ = runEpisodeStart(r.args("ep1"), &stdout, &stderr, spyEnv(&opened))
	if !opened {
		t.Fatalf("with the files gone the run did not reach the terminal: %s", stderr.String())
	}
	// A sessions dir others can read is refused too.
	if err := os.Chmod(r.sessions, 0o755); err != nil {
		t.Fatal(err)
	}
	opened = false
	stderr.Reset()
	if code := runEpisodeStart(r.args("ep1"), &stdout, &stderr, spyEnv(&opened)); code != exitUsage || opened || !strings.Contains(stderr.String(), "chmod 700") {
		t.Fatalf("a 0755 sessions dir: rc %d opened %v %s", code, opened, stderr.String())
	}
}

// The default grants: se-tools without Workspace.Exec, plus it with
// --profile; the launch allows exactly the tools those grants cover.
func TestEpisodeStartDefaultGrantsAndAllowedTools(t *testing.T) {
	without, err := episodeStartGrants("", nil, 50, false)
	if err != nil || hasGrant(without, broker.EffectWorkspaceExec) || len(without) != len(sePresetGrants)-1 {
		t.Fatalf("grants without --profile = %v %v", without, err)
	}
	with, _ := episodeStartGrants("", nil, 50, true)
	if !hasGrant(with, broker.EffectWorkspaceExec) || len(with) != len(sePresetGrants) {
		t.Fatalf("grants with --profile = %v", with)
	}
	if tools := strings.Join(allowedTools(without), ","); strings.Contains(tools, "workspace-exec") || !strings.Contains(tools, "mcp__world__ailang-run") {
		t.Fatalf("allowed without exec: %s", tools)
	}
	if n := len(allowedTools(with)); n != len(worldToolEffects) {
		t.Fatalf("allowed with every grant: %d tools, want %d", n, len(worldToolEffects))
	}
	// Explicit grants are taken exactly, as session new takes them.
	only, _ := episodeStartGrants("", []string{"Workspace.Read=worktree:3"}, 50, true)
	if len(only) != 1 || strings.Join(allowedTools(only), ",") != "mcp__world__ailang-read" {
		t.Fatalf("explicit grants = %v, tools %v", only, allowedTools(only))
	}
}

// worldToolEffects is the manifest's (tool, declared effects) table.
func TestEpisodeToolEffectsMatchTheManifest(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "packages", "se-tools", "transitions.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest []struct {
		ID              string `json:"id"`
		DeclaredEffects []struct {
			Effect string `json:"effect"`
		} `json:"declaredEffects"`
	}
	if err := json.Unmarshal(raw, &manifest); err != nil || len(manifest) == 0 {
		t.Fatalf("instrument failure: %v, %d transitions", err, len(manifest))
	}
	want := map[string]string{}
	for _, d := range manifest {
		var effs []string
		for _, e := range d.DeclaredEffects {
			effs = append(effs, e.Effect)
		}
		want[d.ID] = strings.Join(effs, ",")
	}
	got := map[string]string{}
	for _, te := range worldToolEffects {
		got[te.tool] = strings.Join(te.effects, ",")
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("tool table %v\nmanifest   %v", got, want)
	}
}

// D-189-2 A at the mint: session new --profile names the project and its
// digest on the y/N line, records both, and session list shows them.
// Control: a session minted without --profile is unbound everywhere.
func TestSessionNewProfileBindsAndNamesTheDigest(t *testing.T) {
	r := newEpisodeRun(t)
	p := writeTestProfile(t, r.profiles, "proj", "")
	now := time.Now().Unix()
	term := &promptTerm{in: strings.NewReader("y\n")}
	env := sessionEnv{openTerminal: func() (io.ReadWriteCloser, error) { return term, nil }, now: func() int64 { return now }}
	var stdout, stderr bytes.Buffer
	code := runSessionNew(r.f.args("ep1", "--grant", "Workspace.Exec=worktree:5", "--profile", "proj", "--exec-profiles-dir", r.profiles), &stdout, &stderr, env)
	if code != exitOK {
		t.Fatalf("rc %d %s", code, stderr.String())
	}
	if !strings.Contains(term.out.String(), "bound to exec profile proj "+p.Digest) {
		t.Fatalf("the y/N does not name the project and digest: %s", term.out.String())
	}
	// Control: an unbound mint for another episode.
	plain := &promptTerm{in: strings.NewReader("y\n")}
	env.openTerminal = func() (io.ReadWriteCloser, error) { return plain, nil }
	if code := runSessionNew(r.f.args("ep2", "--grant", "Workspace.Read=worktree:5"), &stdout, &stderr, env); code != exitOK {
		t.Fatalf("unbound: rc %d %s", code, stderr.String())
	}
	if strings.Contains(plain.out.String(), "exec profile") {
		t.Fatalf("an unbound y/N names a profile: %s", plain.out.String())
	}
	for _, row := range r.f.rows(t) {
		bound := row.Project == "proj" && row.ProfileDigest == p.Digest
		if bound != (row.EpisodeID == "ep1") || (row.EpisodeID == "ep2" && (row.Project != "" || row.ProfileDigest != "")) {
			t.Fatalf("row %+v", row)
		}
	}
	var out, errw bytes.Buffer
	if code := runSessionList([]string{"--db", r.f.db}, &out, &errw, testMintEnv("", now)); code != exitOK {
		t.Fatal(errw.String())
	}
	for _, line := range strings.Split(out.String(), "\n") {
		switch {
		case strings.Contains(line, " ep1 "):
			if !strings.Contains(line, "profile proj "+p.Digest) {
				t.Errorf("list does not show the binding: %q", line)
			}
		case strings.Contains(line, " ep2 "):
			if strings.Contains(line, "profile") {
				t.Errorf("list shows a binding for an unbound credential: %q", line)
			}
		}
	}
	out.Reset()
	if code := runSessionList([]string{"--db", r.f.db, "--json", "--episode", "ep1"}, &out, &errw, testMintEnv("", now)); code != exitOK {
		t.Fatal(errw.String())
	}
	var listing []sessionListing
	if err := json.Unmarshal(out.Bytes(), &listing); err != nil || len(listing) != 1 || listing[0].Project != "proj" || listing[0].ProfileDigest != p.Digest {
		t.Fatalf("json listing %s %v", out.String(), err)
	}
}

// A profile that is missing, does not parse or names another project, an
// episode already bound to other bytes, and no terminal at all: each refuses
// with nothing provisioned.
func TestSessionNewProfileRefusals(t *testing.T) {
	r := newEpisodeRun(t)
	writeTestProfile(t, r.profiles, "proj", "")
	if err := os.WriteFile(filepath.Join(r.profiles, "broken.json"), []byte(`{"profile": `), 0o600); err != nil {
		t.Fatal(err)
	}
	other := writeTestProfile(t, t.TempDir(), "other", "")
	raw, _ := json.Marshal(map[string]any{"profile": broker.ExecProfileVersion, "project": other.Project, "path": []string{"/bin"},
		"timeout_ms": 9000, "commands": map[string]any{"true": map[string]any{"argv": []string{"sh"}}}})
	if err := os.WriteFile(filepath.Join(r.profiles, "misnamed.json"), raw, 0o600); err != nil {
		t.Fatal(err)
	}
	for name, tc := range map[string]struct{ profile, want string }{
		"missing":  {"nope", "nope.json"},
		"broken":   {"broken", "broken.json"},
		"misnamed": {"misnamed", `declares project "other"`},
	} {
		t.Run(name, func(t *testing.T) {
			opened := false
			var stdout, stderr bytes.Buffer
			code := runSessionNew(r.f.args("ep1", "--grant", "Workspace.Exec=worktree:5", "--profile", tc.profile, "--exec-profiles-dir", r.profiles), &stdout, &stderr, spyEnv(&opened))
			if code != exitUsage || opened || !strings.Contains(stderr.String(), tc.want) {
				t.Fatalf("rc %d opened %v stderr %q, want %q", code, opened, stderr.String(), tc.want)
			}
			r.f.assertNothingProvisioned(t, "ep1")
		})
	}
	// No terminal with a good profile: the fence refuses (the grader's case).
	var stdout, stderr bytes.Buffer
	if code := runSessionNew(r.f.args("ep1", "--grant", "Workspace.Exec=worktree:5", "--profile", "proj", "--exec-profiles-dir", r.profiles), &stdout, &stderr, noTerminalEnv(1000)); code != exitUsage ||
		!strings.Contains(stderr.String(), "no controlling terminal") {
		t.Fatalf("no terminal: rc %d %s", code, stderr.String())
	}
	r.f.assertNothingProvisioned(t, "ep1")

	// One episode, one profile: once ep1 is bound, other bytes are refused
	// before the fence until the old credential is revoked.
	now := time.Now().Unix()
	if code := runSessionNew(r.f.args("ep1", "--grant", "Workspace.Exec=worktree:5", "--profile", "proj", "--exec-profiles-dir", r.profiles), &stdout, &stderr, testMintEnv("y\n", now)); code != exitOK {
		t.Fatalf("first bound mint: %s", stderr.String())
	}
	writeTestProfile(t, r.profiles, "proj", `, "root": "."`)
	opened := false
	stderr.Reset()
	code := runSessionNew([]string{"ep1", "--db", r.f.db, "--workspace-root", r.f.ws, "--grant", "Workspace.Exec=worktree:5", "--profile", "proj", "--exec-profiles-dir", r.profiles}, &stdout, &stderr, spyEnv(&opened))
	if code != exitUsage || opened || !strings.Contains(stderr.String(), "revoke it first") {
		t.Fatalf("re-mint over a live binding: rc %d opened %v %s", code, opened, stderr.String())
	}
}
