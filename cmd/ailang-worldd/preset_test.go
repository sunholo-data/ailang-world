package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

func presetTestContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	t.Cleanup(cancel)
	return ctx
}
func presetGit(t *testing.T, f episodeFixture, args ...string) string {
	t.Helper()
	s, err := worktree.Output(presetTestContext(t), f.git, append([]string{"-C", f.repo}, args...))
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(s)
}

type presetFixture struct {
	r                     episodeRun
	dir, file, task, base string
	data                  []byte
	profile               *broker.ExecProfile
}

func newPresetFixture(t *testing.T) presetFixture {
	t.Helper()
	r := newEpisodeRun(t)
	presetGit(t, r.f, "branch", "dev")
	base := presetGit(t, r.f, "rev-parse", "dev")
	presetGit(t, r.f, "checkout", "-q", "-b", "feature")
	presetGit(t, r.f, "-c", "user.email=test@example.invalid", "-c", "user.name=test", "commit", "-q", "--allow-empty", "-m", "feature")
	presetGit(t, r.f, "remote", "add", "origin", "https://example.invalid/repo")
	profile := writeTestProfile(t, r.profiles, "proj", "")
	dir := t.TempDir()
	cred, err := credentialsFor(r.f.db, "")
	if err != nil {
		t.Fatal(err)
	}
	p := launchPreset{Schema: "world/preset/v1", Name: "p", Repo: r.f.repo, Remote: "origin", Base: "dev", Profile: "proj", Agent: "codex", Model: "fixture-model", TTLS: 43200,
		Grants: []presetGrant{{Set: "se-tools", Budget: 50}},
		Checks: []json.RawMessage{json.RawMessage(`{"name":"short","argv":["sh","-c","exit 0"],"timeout_s":9}`), json.RawMessage(`{"name":"suite","argv":["go","test","./..."],"unsandboxed":true}`)}}
	p.Daemon, _ = json.Marshal(presetDaemon{Addr: defaultEpisodeAddr, DB: r.f.db, CredentialsDB: cred, WorkspaceRoot: r.f.ws, SessionsDir: r.sessions})
	data, _ := json.MarshalIndent(p, "", "  ")
	file := filepath.Join(dir, "p.json")
	if err := os.WriteFile(file, data, 0o600); err != nil {
		t.Fatal(err)
	}
	task := filepath.Join(t.TempDir(), "t.md")
	if err := os.WriteFile(task, []byte("# Fix fixture\nExact task bytes.\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	h := matchingHealth(t, r.f, r.profiles, profile)
	h.DBPath = r.f.db
	stubDaemon(t, h, nil)
	return presetFixture{r: r, dir: dir, file: file, task: task, base: base, data: data, profile: profile}
}
func (f presetFixture) args(extra ...string) []string {
	return append([]string{"--preset", "p", "--presets-dir", f.dir, "--exec-profiles-dir", f.r.profiles}, extra...)
}
func presetClock() int64 { return time.Date(2026, 10, 28, 13, 52, 0, 0, time.Local).Unix() }
func presetEnv(term *promptTerm) sessionEnv {
	return sessionEnv{openTerminal: func() (io.ReadWriteCloser, error) { return term, nil }, now: presetClock}
}
func readPresetStart(t *testing.T, f presetFixture, ep string) *store.EpisodeStart {
	t.Helper()
	cred, _ := credentialsFor(f.r.f.db, "")
	c, err := store.OpenCredentialsReadOnly(cred)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	start, err := c.EpisodeStart(presetTestContext(t), ep)
	if err != nil {
		t.Fatal(err)
	}
	return start
}

// AC1/2/6/7/13: the local dev base survives HEAD on feature; approved bytes
// remain pinned after edits and revoke. No model or non-loopback calls.
func TestPresetStartLocalBaseAndPinnedHistory(t *testing.T) {
	f := newPresetFixture(t)
	term := &promptTerm{in: strings.NewReader("y\n")}
	var out, errOut bytes.Buffer
	if code := runEpisodeStart(f.args("--task", f.task), &out, &errOut, presetEnv(term)); code != 0 {
		t.Fatalf("%d: %s", code, errOut.String())
	}
	ep := "p-2610281352"
	dir := filepath.Join(f.r.f.ws, ep)
	head, err := worktree.Output(presetTestContext(t), f.r.f.git, []string{"-C", dir, "rev-parse", "HEAD"})
	if err != nil || strings.TrimSpace(head) != f.base {
		t.Fatalf("HEAD %s, want dev %s: %v", head, f.base, err)
	}
	branch, _ := worktree.Output(presetTestContext(t), f.r.f.git, []string{"-C", dir, "branch", "--show-current"})
	if strings.TrimSpace(branch) != "world/"+ep {
		t.Fatal(branch)
	}
	start := readPresetStart(t, f, ep)
	if start == nil || start.Base.SHA != f.base || start.Preset.Bytes != string(f.data) || start.Preset.SHA256 != serviceHash(f.data) || start.Checks[0].Exec != "true" || start.Checks[1].TimeoutS != 900 {
		t.Fatalf("%+v", start)
	}
	for _, needle := range []string{start.Preset.SHA256, f.base, "short (sandboxed: true)", "suite (UNSANDBOXED", "never fetched"} {
		if !strings.Contains(term.out.String(), needle) {
			t.Errorf("prompt missing %s: %s", needle, term.out.String())
		}
	}
	rows := f.r.f.rows(t)
	var grants []broker.Capability
	_ = json.Unmarshal([]byte(rows[0].GrantsJSON), &grants)
	want, _ := sessionNewGrants("se-tools", nil, 50)
	if len(grants) != len(want) {
		t.Fatal(grants)
	}
	for i, g := range grants {
		if g.Effect != want[i].Effect || g.Scope != want[i].Scope || g.Budget != want[i].Budget || !strings.Contains(term.out.String(), g.Effect+"=") {
			t.Fatal("grant set/prompt drift", g)
		}
	}
	for _, suffix := range []string{".session", ".task.md"} {
		path := filepath.Join(f.r.sessions, ep+suffix)
		info, err := os.Stat(path)
		if err != nil || info.Mode().Perm() != 0o600 {
			t.Fatalf("%s: %v", path, err)
		}
	}
	copied, _ := os.ReadFile(filepath.Join(f.r.sessions, ep+".task.md"))
	task, _ := os.ReadFile(f.task)
	if !bytes.Equal(copied, task) {
		t.Fatal("task bytes drift")
	}
	_ = os.WriteFile(f.file, []byte("{}"), 0o600)
	if got := readPresetStart(t, f, ep); !reflect.DeepEqual(start, got) {
		t.Fatal("preset edit changed history")
	}
	if code := runSessionRevoke([]string{"--db", f.r.f.db, "--episode", ep}, &out, &errOut); code != 0 {
		t.Fatalf("revoke: %d %s", code, errOut.String())
	}
	if got := readPresetStart(t, f, ep); !reflect.DeepEqual(start, got) {
		t.Fatal("revoke deleted history")
	}
	row, _ := json.Marshal(start)
	t.Logf("AC1 prompt:\n%s\nstart row:\n%s", term.out.String(), row)
}

// AC2/3: decline and missing terminal leave all provisioned state untouched.
func TestPresetDeclineAndNoTerminal(t *testing.T) {
	for _, noTTY := range []bool{false, true} {
		t.Run(fmt.Sprint(noTTY), func(t *testing.T) {
			f := newPresetFixture(t)
			term := &promptTerm{in: strings.NewReader("n\n")}
			env := presetEnv(term)
			if noTTY {
				env.openTerminal = func() (io.ReadWriteCloser, error) { return nil, os.ErrNotExist }
			}
			var out, errOut bytes.Buffer
			code := runEpisodeStart(f.args("--agent", "claude", "--model", "override", "--ttl", "12", "--branch", "override-branch"), &out, &errOut, env)
			if noTTY && code == 0 {
				t.Fatal("missing terminal accepted")
			}
			f.r.f.assertNothingProvisioned(t, "p-2610281352")
			f.r.assertNoSessionFiles(t)
			if !noTTY {
				for _, s := range []string{"override  agent", "override  model", "override  ttl", "override  branch"} {
					if !strings.Contains(term.out.String(), s) {
						t.Fatal(term.out.String())
					}
				}
			}
		})
	}
	var out, errOut bytes.Buffer
	runEpisodeStart([]string{"--help"}, &out, &errOut, sessionEnv{})
	for _, bad := range []string{"--yes", "--launch", "--resume"} {
		if strings.Contains(out.String(), bad) {
			t.Fatal("unexpected flag", bad)
		}
	}
}

// AC4/5: every unsafe layout or bad check refuses before daemon and fence.
func TestPresetRefusalsBeforeDaemon(t *testing.T) {
	cases := []string{"symlink", "writable", "dirwritable", "longname", "inside", "unknown", "name", "reserved", "slow", "make", "duplicate", "absolute", "parent", "missing", "grant"}
	for _, tc := range cases {
		t.Run(tc, func(t *testing.T) {
			f := newPresetFixture(t)
			name, dir := "p", f.dir
			var obj map[string]any
			_ = json.Unmarshal(f.data, &obj)
			switch tc {
			case "symlink":
				target := filepath.Join(t.TempDir(), "target")
				_ = os.Rename(f.file, target)
				_ = os.Symlink(target, f.file)
			case "writable":
				_ = os.Chmod(f.file, 0o660)
			case "dirwritable":
				_ = os.Chmod(f.dir, 0o770)
			case "longname":
				// 52 + "-yymmddhhmm" + "-9" is 65, one past the episode grammar.
				name = strings.Repeat("a", 52)
			case "inside":
				dir = f.r.f.ws
				_ = os.WriteFile(filepath.Join(dir, "p.json"), f.data, 0o600)
			case "unknown":
				obj["surprise"] = true
			case "name":
				obj["name"] = "wrong"
			case "reserved":
				name = "se-tools"
				_ = os.WriteFile(filepath.Join(dir, name+".json"), f.data, 0o600)
			case "slow":
				obj["checks"] = []any{map[string]any{"argv": []string{"go", "test", "./..."}, "timeout_s": 900}}
			case "make":
				obj["checks"] = []any{map[string]any{"argv": []string{"make", "test"}}}
			case "duplicate":
				obj["checks"] = []any{map[string]any{"argv": []string{"x"}}, map[string]any{"argv": []string{"x"}}}
			case "absolute":
				obj["checks"] = []any{map[string]any{"argv": []string{"x"}, "cwd": "/abs"}}
			case "parent":
				obj["checks"] = []any{map[string]any{"argv": []string{"x"}, "cwd": ".."}}
			case "missing":
				name = "absent"
			}
			switch tc {
			case "unknown", "name", "slow", "make", "duplicate", "absolute", "parent":
				data, _ := json.Marshal(obj)
				_ = os.WriteFile(f.file, data, 0o600)
			}
			fetchDaemonHealth = func(string) (daemon.HealthResponse, error) {
				t.Fatal("daemon contacted")
				return daemon.HealthResponse{}, nil
			}
			var out, errOut bytes.Buffer
			if tc == "reserved" {
				if _, err := loadLaunchPreset(presetTestContext(t), name, dir, f.r.profiles, f.r.f.ws); err == nil {
					t.Fatal("reserved accepted")
				}
				return
			}
			args := []string{"--preset", name, "--presets-dir", dir, "--exec-profiles-dir", f.r.profiles}
			if tc == "grant" {
				args = append(args, "--grant", "Workspace.Read=worktree:1")
			}
			if code := runEpisodeStart(args, &out, &errOut, sessionEnv{}); code == 0 {
				t.Fatal("accepted", tc)
			}
			want := map[string]string{"dirwritable": "presets directory", "longname": "at most 51 characters", "name": `the file says name "wrong"`}[tc]
			if want != "" && !strings.Contains(errOut.String(), want) {
				t.Fatalf("%s refusal lacks %q:\n%s", tc, want, errOut.String())
			}
			if tc == "slow" || tc == "make" {
				errOut.Reset()
				if code := runPreset([]string{"check", name, "--presets-dir", dir, "--exec-profiles-dir", f.r.profiles, "--no-daemon"}, &out, &errOut); code == 0 {
					t.Fatal("check accepted")
				}
				t.Logf("refused preset check:\n%s", errOut.String())
			}
		})
	}
}

func TestPresetCheckMatcherUsesProfileGrammarAndTimeout(t *testing.T) {
	p := &broker.ExecProfile{Project: "fixture", Root: ".", TimeoutMS: 9000, Commands: map[string]broker.ExecProfileCommand{
		"go-test": {ExecCommand: broker.ExecCommand{Argv: []string{"go", "test"}, Flags: map[string]string{"-run": broker.ExecClassRegex}, Positional: broker.ExecPositionalPkgPattern, MaxArgs: 16}},
	}}
	for _, tc := range []struct {
		name    string
		argv    []string
		cwd     string
		timeout float64
		valid   bool
	}{
		{"valid", []string{"go", "test", "-run", "TestOne", "./..."}, ".", 9, true},
		{"slow", []string{"go", "test", "./..."}, ".", 900, false},
		{"grammar", []string{"go", "test", "-bad"}, ".", 9, false},
		{"cwd", []string{"go", "test", "./..."}, "backend", 9, false},
		{"make", []string{"make", "test"}, ".", 9, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := []store.StartCheck{{Name: tc.name, Argv: tc.argv, Cwd: tc.cwd, TimeoutS: tc.timeout}}
			err := pinPresetChecks(c, p)
			if (err == nil) != tc.valid {
				t.Fatalf("%+v %v", c, err)
			}
			if tc.valid && (c[0].Exec != "go-test" || !reflect.DeepEqual(c[0].Args, []string{"-run=TestOne", "./..."})) {
				t.Fatalf("unpinned/unnormalized: %+v", c)
			}
		})
	}
}

func TestPresetSharedGraderVectors(t *testing.T) {
	data, err := os.ReadFile("../../scripts/trial/testdata/grader_spec_vectors.json")
	if err != nil {
		t.Fatal(err)
	}
	var vectors []struct {
		Name     string             `json:"name"`
		Graders  []json.RawMessage  `json:"graders"`
		Error    string             `json:"error"`
		Expected []store.StartCheck `json:"expected"`
	}
	if err := json.Unmarshal(data, &vectors); err != nil {
		t.Fatal(err)
	}
	if len(vectors) == 0 {
		t.Fatal("no vectors")
	}
	for _, v := range vectors {
		t.Run(v.Name, func(t *testing.T) {
			got, err := validateGraderChecks(v.Graders, "tps-a")
			if v.Error != "" {
				if err == nil || !strings.Contains(err.Error(), v.Error) {
					t.Fatalf("%v; want %s", err, v.Error)
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				for i := range got {
					got[i].Sandbox = ""
				}
				if !reflect.DeepEqual(got, v.Expected) {
					t.Fatalf("%+v != %+v", got, v.Expected)
				}
			}
		})
	}
}

func TestPresetDerivedCollisionsAndExistingBranch(t *testing.T) {
	f := newPresetFixture(t)
	privateSessions := t.TempDir()
	if err := os.WriteFile(filepath.Join(privateSessions, "p-2610281352.session"), []byte("occupied"), 0o600); err != nil {
		t.Fatal(err)
	}
	if id, err := derivePresetEpisode("p", f.r.f.ws, privateSessions, presetClock()); err != nil || id != "p-2610281352-2" {
		t.Fatalf("session collision: %s %v", id, err)
	}
	for n := 1; n <= 9; n++ {
		id, err := derivePresetEpisode("p", f.r.f.ws, f.r.sessions, presetClock())
		want := "p-2610281352"
		if n > 1 {
			want += fmt.Sprintf("-%d", n)
		}
		if err != nil || id != want {
			t.Fatalf("%s %v, want %s", id, err, want)
		}
		_ = os.Mkdir(filepath.Join(f.r.f.ws, id), 0o700)
	}
	if _, err := derivePresetEpisode("p", f.r.f.ws, f.r.sessions, presetClock()); err == nil || !strings.Contains(err.Error(), "F5") {
		t.Fatal(err)
	}
	presetGit(t, f.r.f, "branch", "existing", "dev")
	term := &promptTerm{in: strings.NewReader("y\n")}
	var out, errOut bytes.Buffer
	if code := runEpisodeStart(append(f.args("--branch", "existing"), "explicit"), &out, &errOut, presetEnv(term)); code != 0 {
		t.Fatalf("%d: %s", code, errOut.String())
	}
	if !strings.Contains(term.out.String(), "existing branch existing") {
		t.Fatal(term.out.String())
	}
	if !readPresetStart(t, f, "explicit").BranchExisted {
		t.Fatal("existing branch not recorded")
	}
	// The same branch cannot be checked out twice.
	if code := runEpisodeStart(append(f.args("--branch", "existing"), "another"), &out, &errOut, sessionEnv{}); code == 0 {
		t.Fatal("busy branch accepted")
	}
}

func TestPresetVerbsAndService(t *testing.T) {
	f := newPresetFixture(t)
	var out, errOut bytes.Buffer
	for _, verb := range []string{"check", "show", "list"} {
		args := []string{verb, "--presets-dir", f.dir, "--exec-profiles-dir", f.r.profiles, "--no-daemon"}
		if verb != "list" {
			args = append(args, "p")
		}
		if verb == "show" {
			args = append(args, "--json")
		}
		if code := runPreset(args, &out, &errOut); code != 0 {
			t.Fatalf("%s: %d %s", verb, code, errOut.String())
		}
		if verb == "check" {
			t.Logf("valid preset check:\n%s", out.String())
		}
		out.Reset()
	}
	home := t.TempDir()
	t.Setenv("HOME", home)
	paths := homeServicePaths(home)
	r := serviceRecord{Schema: "world/service/v1", Label: worldServiceLabel, Plist: paths.Plist, Binary: paths.Binary,
		Address: "127.0.0.1:7644", Store: f.r.f.db, Credentials: filepath.Join(filepath.Dir(f.r.f.db), "credentials.db"), Workspace: f.r.f.ws}
	data, _ := json.Marshal(r)
	_ = os.MkdirAll(filepath.Dir(paths.Record), 0o700)
	_ = os.WriteFile(paths.Record, data, 0o600)
	d, err := resolvePresetDaemon(json.RawMessage("\"service\""))
	if err != nil || d.Addr != r.Address || d.DB != r.Store || d.CredentialsDB != r.Credentials {
		t.Fatalf("%+v %v", d, err)
	}
	p, err := loadLaunchPreset(presetTestContext(t), "p", f.dir, f.r.profiles, f.r.f.ws)
	if err != nil {
		t.Fatal(err)
	}
	h := matchingHealth(t, f.r.f, f.r.profiles, f.profile)
	h.DBPath = filepath.Join(t.TempDir(), "other.db")
	if err := checkPresetDaemon(h, &mintRequest{}, p, d.Addr); err == nil || !strings.Contains(err.Error(), "F6") {
		t.Fatal(err)
	}
	newDir := filepath.Join(t.TempDir(), "presets")
	args := []string{"init", "new", "--presets-dir", newDir, "--repo", f.r.f.repo, "--base", "dev"}
	if code := runPreset(args, &out, &errOut); code != 0 {
		t.Fatal(errOut.String())
	}
	for path, mode := range map[string]os.FileMode{newDir: 0o700, filepath.Join(newDir, "new.json"): 0o600} {
		info, err := os.Stat(path)
		if err != nil || info.Mode().Perm() != mode {
			t.Fatalf("%s %v", path, err)
		}
	}
	if code := runPreset(args, &out, &errOut); code == 0 {
		t.Fatal("init overwrote")
	}
}
