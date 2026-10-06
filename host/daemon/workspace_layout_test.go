package daemon

// Row 141 (w-workspace-project-layouts) M1: the module-root sandbox. These
// run in CI with fakeToolBin (the fake reports fs_sandbox = its cwd and
// answers every dispatch with `fake:<cwd>`); the real-binary arms live in
// setools_e2e_test.go.

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
)

// wsRead runs one Workspace.Read of path through the episode's registry.
func wsRead(t *testing.T, d *Daemon, episode, path string) map[string]any {
	t.Helper()
	h := d.workspace.registry(episode)[broker.EffectWorkspaceRead]
	if h == nil {
		t.Fatalf("registry(%s) has no Workspace.Read handler", episode)
	}
	body, _ := json.Marshal(map[string]string{"op": "read", "path": path})
	out, err := h.Execute(boundedTestContext(t), broker.EffectRequest{Effect: broker.EffectWorkspaceRead, Scope: broker.WorkspaceScope, Cost: 1}, body)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var resp map[string]any
	if err := json.Unmarshal(out, &resp); err != nil {
		t.Fatal(err)
	}
	return resp
}

func mkdirs(t *testing.T, dirs ...string) {
	t.Helper()
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
}

func TestWorkspaceModuleRootGrammar(t *testing.T) {
	type row struct {
		name string
		cfg  func(cfg *Config)
		want string // "" = accepted, else a substring of the refusal
	}
	def := func(rel string) func(*Config) { return func(c *Config) { c.WorkspaceModuleRoot = rel } }
	pair := func(p ...string) func(*Config) { return func(c *Config) { c.WorkspaceEpisodeModuleRoots = p } }
	rows := []row{
		{"tools", def("tools"), ""},
		{"tools/sub", def("tools/sub"), ""},
		{". as the default", def("."), ""},
		{"ep1=.", pair("ep1=."), ""},
		{"ep1= (empty REL)", pair("ep1="), "--workspace-episode-module-root"},
		{"/abs", def("/abs"), "--workspace-module-root"},
		{"..", def(".."), "--workspace-module-root"},
		{"../x", def("../x"), "--workspace-module-root"},
		{"tools/../../x", def("tools/../../x"), "--workspace-module-root"},
		{"-x", def("-x"), "--workspace-module-root"},
		{"tools/", def("tools/"), "is not clean"},
		{`a"b`, def(`a"b`), "--workspace-module-root"},
		{"EP outside the grammar", pair("EP1=tools"), "outside the episode grammar"},
		{"no =", pair("ep1"), "is not EP=REL"},
		{"mapped twice", pair("ep1=tools", "ep1=tools"), "a second time"},
	}
	for _, tc := range rows {
		t.Run(tc.name, func(t *testing.T) {
			f := newWSFixture(t)
			mkdirs(t, filepath.Join(f.root, "ep1", "tools", "sub"))
			cfg := Config{DBPath: f.db, WorkspaceRoot: f.root, ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease)}
			tc.cfg(&cfg)
			_, err := newWSDaemon(t, cfg)
			if tc.want == "" {
				if err != nil {
					t.Fatalf("New = %v, want accepted", err)
				}
				return
			}
			var startup *StartupError
			if !errors.As(err, &startup) || startup.Stage != StageConfig || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("New = %v, want a %s StartupError containing %q", err, StageConfig, tc.want)
			}
			if _, err := os.Stat(f.db); !os.IsNotExist(err) {
				t.Fatalf("a refused module root left the store open (stat %s: %v)", f.db, err)
			}
		})
	}
}

func TestWorkspaceModuleRootIsTheSandbox(t *testing.T) {
	f := newWSFixture(t)
	mkdirs(t, filepath.Join(f.root, "ep1", "tools"))
	d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, WorkspaceModuleRoot: "tools",
		ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease)})
	reg := d.workspace.registry("ep1")
	sandbox := filepath.Join(f.root, "ep1", "tools")
	got, err := os.ReadFile(filepath.Join(f.stateDir, "policies", "ep1.toml"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := broker.RenderEpisodePolicy(sandbox)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatalf("rendered policy =\n%s\nwant (fs_sandbox = %s)\n%s", got, sandbox, want)
	}
	if gotNames := registryNames(reg); strings.Join(gotNames, ",") != strings.Join(wantWorkspaceEffects, ",") {
		t.Fatalf("registry = %v, want %v", gotNames, wantWorkspaceEffects)
	}
	if resp := wsRead(t, d, "ep1", "x"); resp["content"] != "fake:"+sandbox {
		t.Fatalf("read ran in %v, want the module root %s", resp["content"], sandbox)
	}
	if n := f.summaries(t); n != 1 {
		t.Fatalf("summaries = %d, want 1", n)
	}
}

func TestWorkspaceModuleRootSymlinkOrMissingIsR8(t *testing.T) {
	for _, tc := range []struct {
		name  string
		setup func(t *testing.T, f wsFixture)
	}{
		{"symlink to a real dir", func(t *testing.T, f wsFixture) {
			mkdirs(t, filepath.Join(f.root, "ep1", "real"))
			if err := os.Symlink(filepath.Join(f.root, "ep1", "real"), filepath.Join(f.root, "ep1", "tools")); err != nil {
				t.Fatal(err)
			}
		}},
		{"symlink to a sibling episode's tools", func(t *testing.T, f wsFixture) {
			mkdirs(t, filepath.Join(f.root, "ep2", "tools"))
			if err := os.Symlink("../ep2/tools", filepath.Join(f.root, "ep1", "tools")); err != nil {
				t.Fatal(err)
			}
		}},
		{"regular file", func(t *testing.T, f wsFixture) { writeFile(t, filepath.Join(f.root, "ep1", "tools"), "x") }},
		{"missing", func(t *testing.T, f wsFixture) {}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newWSFixture(t)
			tc.setup(t, f)
			var log bytes.Buffer
			d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, WorkspaceModuleRoot: "tools", ErrorLog: &log,
				ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease)})
			if reg := d.workspace.registry("ep1"); len(reg) != 0 {
				t.Fatalf("registry = %v, want empty", registryNames(reg))
			}
			if f.summaries(t) != 0 || f.dispatches(t) != 0 {
				t.Fatalf("a refused module root ran the tool (summaries %d, dispatches %d)", f.summaries(t), f.dispatches(t))
			}
			if n := strings.Count(log.String(), "\n"); n != 1 || !strings.Contains(log.String(), `workspace tools unavailable for episode "ep1": module root "tools":`) {
				t.Fatalf("operator log (%d lines) = %q", n, log.String())
			}
			if tc.name == "missing" {
				if _, err := os.Lstat(filepath.Join(f.root, "ep1", "tools")); !os.IsNotExist(err) {
					t.Fatalf("the missing module root was created (lstat: %v)", err)
				}
			}
		})
	}
}

func TestWorkspaceEpisodeModuleRootOverrides(t *testing.T) {
	f := newWSFixture(t)
	mkdirs(t, filepath.Join(f.root, "ep1", "tools"))
	d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, WorkspaceModuleRoot: "tools",
		WorkspaceEpisodeModuleRoots: []string{"ep2=."}, ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease)})
	if resp := wsRead(t, d, "ep1", "x"); resp["content"] != "fake:"+filepath.Join(f.root, "ep1", "tools") {
		t.Fatalf("ep1 read = %v", resp)
	}
	if resp := wsRead(t, d, "ep2", "x"); resp["content"] != "fake:"+filepath.Join(f.root, "ep2") {
		t.Fatalf("ep2 read = %v, want the worktree root (override ep2=.)", resp)
	}
}

func TestWorkspaceModuleRootNeedsWorkspaceRoot(t *testing.T) {
	for name, cfg := range map[string]Config{
		"default":  {WorkspaceModuleRoot: "tools"},
		"episodes": {WorkspaceEpisodeModuleRoots: []string{"ep1=tools"}},
	} {
		t.Run(name, func(t *testing.T) {
			f := newWSFixture(t)
			cfg.DBPath = f.db
			_, err := newWSDaemon(t, cfg)
			var startup *StartupError
			if !errors.As(err, &startup) || startup.Stage != StageConfig || !strings.Contains(err.Error(), "need --workspace-root") {
				t.Fatalf("New = %v, want a %s StartupError containing %q", err, StageConfig, "need --workspace-root")
			}
			if _, err := os.Stat(f.db); !os.IsNotExist(err) {
				t.Fatalf("the store was opened before the refusal (stat: %v)", err)
			}
		})
	}
	// Control: the same flag WITH --workspace-root starts and opens the store,
	// so the absence above is the instrument seeing a refusal, not a blind spot.
	f := newWSFixture(t)
	mkdirs(t, filepath.Join(f.root, "ep1", "tools"))
	mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, WorkspaceModuleRoot: "tools",
		ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease)})
	if _, err := os.Stat(f.db); err != nil {
		t.Fatalf("control: the store is absent after a successful start: %v", err)
	}
}
