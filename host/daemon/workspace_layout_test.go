package daemon

// Row 141 (w-workspace-project-layouts) M1: the module-root sandbox, and M2:
// the read-only registry package cache. These
// run in CI with fakeToolBin (the fake reports fs_sandbox = its cwd and
// answers every dispatch with `fake:<cwd>`); the real-binary arms live in
// setools_e2e_test.go.

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"testing"
	"time"

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

// configureStubExec gives the daemon a Workspace.Exec profile after New, as
// workspace_exec_test.go does (row 152: exec must survive an AILANG refusal).
func configureStubExec(t *testing.T, d *Daemon, f wsFixture) {
	t.Helper()
	d.workspace.exec = &workspaceExec{profile: execProfile(t), sandbox: stubExecSandbox(t, f.stateDir)}
}

// requireOnlyExec fails unless the registry is exactly [Workspace.Exec]: the
// AILANG family is refused and the exec family is bound (row 152 AC1.5).
func requireOnlyExec(t *testing.T, reg broker.Registry) {
	t.Helper()
	if names := registryNames(reg); strings.Join(names, ",") != broker.EffectWorkspaceExec {
		t.Fatalf("registry = %v, want exactly [%s]", names, broker.EffectWorkspaceExec)
	}
}

// requireExecRunsIn drives one `where` call through the registry's bound exec
// handler and requires exit 0 with stdout == want. Membership alone is not
// enough: a refusal handler would also be a member. It reuses the caller's
// single registry() result, since another registry() call prints another
// operator line.
func requireExecRunsIn(t *testing.T, reg broker.Registry, want string) {
	t.Helper()
	out, err := execCall(t, reg[broker.EffectWorkspaceExec])
	if err != nil {
		t.Fatalf("Workspace.Exec on the refused episode: %v", err)
	}
	var res struct {
		ExitCode *int   `json:"exit_code"`
		Stdout   string `json:"stdout"`
	}
	if err := json.Unmarshal(out, &res); err != nil || res.ExitCode == nil || *res.ExitCode != 0 || strings.TrimSpace(res.Stdout) != want {
		t.Fatalf("Workspace.Exec = %s (%v), want exit_code 0 and stdout %s", out, err, want)
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
			configureStubExec(t, d, f)
			reg := d.workspace.registry("ep1")
			requireOnlyExec(t, reg)
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
			requireExecRunsIn(t, reg, filepath.Join(f.root, "ep1"))
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

// ---------------------------------------------------------------------------
// M2: the read-only registry package cache
// ---------------------------------------------------------------------------

// lockdown makes dir a read-only tree (0o444 files, 0o555 dirs; symlinks and
// special files are left alone, since chmod would follow a link), bottom-up so
// every directory is still traversable while its children are changed. It
// registers, right away, the cleanup that makes the tree removable again:
// t.TempDir's own cleanup was registered first, so this runs before it.
func lockdown(t *testing.T, dir string) {
	t.Helper()
	var paths []string
	_ = filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err == nil {
			paths = append(paths, p)
		}
		return nil
	})
	t.Cleanup(func() {
		_ = filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
			if d == nil {
				return nil
			}
			switch {
			case d.Type()&os.ModeSymlink != 0:
			case d.IsDir():
				_ = os.Chmod(p, 0o755)
			case d.Type().IsRegular():
				_ = os.Chmod(p, 0o644)
			}
			return nil
		})
	})
	for i := len(paths) - 1; i >= 0; i-- {
		fi, err := os.Lstat(paths[i])
		if err != nil || fi.Mode()&os.ModeSymlink != 0 {
			continue
		}
		switch {
		case fi.IsDir():
			if err := os.Chmod(paths[i], 0o555); err != nil {
				t.Fatal(err)
			}
		case fi.Mode().IsRegular():
			if err := os.Chmod(paths[i], 0o444); err != nil {
				t.Fatal(err)
			}
		}
	}
}

// readOnlyPackageCache builds dir as a snapshot of the packages (each
// "ns/name/ver"): an ailang.toml and one .ail file apiece, then lockdown.
func readOnlyPackageCache(t *testing.T, dir string, pkgs ...string) {
	t.Helper()
	buildPackageCache(t, dir, pkgs...)
	lockdown(t, dir)
}

func buildPackageCache(t *testing.T, dir string, pkgs ...string) {
	t.Helper()
	for _, pkg := range pkgs {
		parts := strings.Split(pkg, "/")
		mkdirs(t, filepath.Join(dir, pkg))
		writeFile(t, filepath.Join(dir, pkg, "ailang.toml"),
			fmt.Sprintf("[package]\nname = %q\nversion = %q\nedition = \"1\"\n", parts[0]+"/"+parts[1], parts[2]))
		writeFile(t, filepath.Join(dir, pkg, "greet.ail"), "module "+pkg+"/greet\n\nexport func greet() -> string {\n  \"hi\"\n}\n")
	}
}

type lockEntry struct{ Name, Version, Source string }

// writeLock writes an ailang.lock marshalled from a Go struct.
func writeLock(t *testing.T, path string, entries ...lockEntry) {
	t.Helper()
	type pkg struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Source      string `json:"source"`
		ContentHash string `json:"content_hash"`
	}
	doc := struct {
		Schema        string `json:"schema"`
		SchemaVersion string `json:"schema_version"`
		Packages      []pkg  `json:"packages"`
	}{Schema: "ailang.lock/v1", SchemaVersion: "1.0.0", Packages: []pkg{}}
	for _, e := range entries {
		doc.Packages = append(doc.Packages, pkg{e.Name, e.Version, e.Source, "sha256:" + strings.Repeat("0", 64)})
	}
	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, path, string(data))
}

// withUID sets the geteuid seam for the test (serial: it is a package var).
func withUID(t *testing.T, uid int) {
	t.Helper()
	old := geteuid
	geteuid = func() int { return uid }
	t.Cleanup(func() { geteuid = old })
}

// pkgcacheDir is where the fixtures place the snapshot: outside f.root
// (base/ws) and f.stateDir (base/state).
func (f wsFixture) pkgcacheDir() string { return filepath.Join(f.base, "pkgcache") }

func (f wsFixture) cfgWithCache(t *testing.T, dir string, log *bytes.Buffer) Config {
	cfg := Config{DBPath: f.db, WorkspaceRoot: f.root, WorkspacePackageCache: dir,
		ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease)}
	if log != nil {
		cfg.ErrorLog = log
	}
	return cfg
}

func (f wsFixture) registryLink() string {
	return filepath.Join(f.stateDir, "cache", "ep1", ".ailang", "cache", "registry")
}

func TestWorkspacePackageCacheStartupChecks(t *testing.T) {
	type row struct {
		name  string
		uid   int
		build func(t *testing.T, f wsFixture) string // returns the --workspace-package-cache value
		want  string                                 // "" = accepted
	}
	valid := func(t *testing.T, f wsFixture) string {
		readOnlyPackageCache(t, f.pkgcacheDir(), "acme/util/0.1.0")
		return f.pkgcacheDir()
	}
	rows := []row{
		{"accepted", 1000, valid, ""},
		{"not a dir", 1000, func(t *testing.T, f wsFixture) string {
			writeFile(t, f.pkgcacheDir(), "x")
			return f.pkgcacheDir()
		}, "not a directory"},
		{"inside the workspace root", 1000, func(t *testing.T, f wsFixture) string {
			readOnlyPackageCache(t, filepath.Join(f.root, "pkgs"), "acme/util/0.1.0")
			return filepath.Join(f.root, "pkgs")
		}, "inside --workspace-root"},
		{"inside the state dir", 1000, func(t *testing.T, f wsFixture) string {
			readOnlyPackageCache(t, filepath.Join(f.stateDir, "pkgs"), "acme/util/0.1.0")
			return filepath.Join(f.stateDir, "pkgs")
		}, "inside the state directory"},
		{"writable file", 1000, func(t *testing.T, f wsFixture) string {
			valid(t, f)
			if err := os.Chmod(filepath.Join(f.pkgcacheDir(), "acme/util/0.1.0/greet.ail"), 0o644); err != nil {
				t.Fatal(err)
			}
			return f.pkgcacheDir()
		}, "is writable"},
		{"writable inner dir", 1000, func(t *testing.T, f wsFixture) string {
			valid(t, f)
			if err := os.Chmod(filepath.Join(f.pkgcacheDir(), "acme/util"), 0o755); err != nil {
				t.Fatal(err)
			}
			return f.pkgcacheDir()
		}, "is writable"},
		{"DIR itself writable", 1000, func(t *testing.T, f wsFixture) string {
			valid(t, f)
			if err := os.Chmod(f.pkgcacheDir(), 0o755); err != nil {
				t.Fatal(err)
			}
			return f.pkgcacheDir()
		}, "is writable"},
		{"group-writable file", 1000, func(t *testing.T, f wsFixture) string {
			valid(t, f)
			if err := os.Chmod(filepath.Join(f.pkgcacheDir(), "acme/util/0.1.0/greet.ail"), 0o464); err != nil {
				t.Fatal(err)
			}
			return f.pkgcacheDir()
		}, "is writable"},
		{"other-writable file", 1000, func(t *testing.T, f wsFixture) string {
			valid(t, f)
			if err := os.Chmod(filepath.Join(f.pkgcacheDir(), "acme/util/0.1.0/greet.ail"), 0o446); err != nil {
				t.Fatal(err)
			}
			return f.pkgcacheDir()
		}, "is writable"},
		{"group-writable dir", 1000, func(t *testing.T, f wsFixture) string {
			valid(t, f)
			if err := os.Chmod(filepath.Join(f.pkgcacheDir(), "acme/util"), 0o575); err != nil {
				t.Fatal(err)
			}
			return f.pkgcacheDir()
		}, "is writable"},
		{"other-writable dir", 1000, func(t *testing.T, f wsFixture) string {
			valid(t, f)
			if err := os.Chmod(filepath.Join(f.pkgcacheDir(), "acme/util"), 0o557); err != nil {
				t.Fatal(err)
			}
			return f.pkgcacheDir()
		}, "is writable"},
		{"ailang.toml one level too deep", 1000, func(t *testing.T, f wsFixture) string {
			mkdirs(t, filepath.Join(f.pkgcacheDir(), "acme/util/0.1.0/sub"))
			writeFile(t, filepath.Join(f.pkgcacheDir(), "acme/util/0.1.0/sub/ailang.toml"), "[package]\n")
			lockdown(t, f.pkgcacheDir())
			return f.pkgcacheDir()
		}, "holds no package"},
		{"ailang.toml is a directory", 1000, func(t *testing.T, f wsFixture) string {
			mkdirs(t, filepath.Join(f.pkgcacheDir(), "acme/util/0.1.0/ailang.toml"))
			lockdown(t, f.pkgcacheDir())
			return f.pkgcacheDir()
		}, "holds no package"},
		{"symlink inside", 1000, func(t *testing.T, f wsFixture) string {
			buildPackageCache(t, f.pkgcacheDir(), "acme/util/0.1.0")
			sibling := filepath.Join(f.base, "sibling")
			buildPackageCache(t, sibling, "acme/other/0.1.0")
			if err := os.Symlink(sibling, filepath.Join(f.pkgcacheDir(), "link")); err != nil {
				t.Fatal(err)
			}
			lockdown(t, sibling)
			lockdown(t, f.pkgcacheDir())
			return f.pkgcacheDir()
		}, "is a symlink"},
		{"zero packages", 1000, func(t *testing.T, f wsFixture) string {
			mkdirs(t, filepath.Join(f.pkgcacheDir(), "acme", "util"))
			lockdown(t, f.pkgcacheDir())
			return f.pkgcacheDir()
		}, "holds no package"},
		{"uid 0", 0, valid, "uid 0"},
		{"a FIFO inside", 1000, func(t *testing.T, f wsFixture) string {
			buildPackageCache(t, f.pkgcacheDir(), "acme/util/0.1.0")
			if err := syscall.Mkfifo(filepath.Join(f.pkgcacheDir(), "pipe"), 0o444); err != nil {
				t.Fatal(err)
			}
			lockdown(t, f.pkgcacheDir())
			return f.pkgcacheDir()
		}, "neither a directory nor a regular file"},
	}
	for _, tc := range rows {
		t.Run(tc.name, func(t *testing.T) {
			f := newWSFixture(t)
			withUID(t, tc.uid)
			dir := tc.build(t, f)
			var log bytes.Buffer
			d, err := newWSDaemon(t, f.cfgWithCache(t, dir, &log))
			if tc.want == "" {
				if err != nil {
					t.Fatalf("New = %v, want accepted", err)
				}
				// Precondition (kimi r2): the accepted tree really holds no symlink.
				symlinks, total := 0, 0
				_ = filepath.WalkDir(dir, func(p string, de os.DirEntry, err error) error {
					if err == nil {
						total++
						if fi, err := os.Lstat(p); err == nil && fi.Mode()&os.ModeSymlink != 0 {
							symlinks++
						}
					}
					return nil
				})
				if symlinks != 0 || total == 0 {
					t.Fatalf("accepted tree: %d symlinks of %d entries, want 0 of >0", symlinks, total)
				}
				if want := "workspace package cache " + dir + ": 1 packages, read-only, sha256:"; !strings.Contains(log.String(), want) {
					t.Fatalf("startup log = %q, want it to contain %q", log.String(), want)
				}
				_, d1, _, err1 := resolvePackageCache(dir, f.root, f.stateDir)
				_, d2, _, err2 := resolvePackageCache(dir, f.root, f.stateDir)
				if err1 != nil || err2 != nil || d1 != d2 || d1 != d.workspace.packageCacheDigest {
					t.Fatalf("digests %q %q (%v %v), startup %q: want all equal", d1, d2, err1, err2, d.workspace.packageCacheDigest)
				}
				return
			}
			var startup *StartupError
			if !errors.As(err, &startup) || startup.Stage != StageConfig ||
				!strings.Contains(err.Error(), "--workspace-package-cache") || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("New = %v, want a %s StartupError naming --workspace-package-cache and %q", err, StageConfig, tc.want)
			}
		})
	}
	t.Run("no --workspace-root", func(t *testing.T) {
		f := newWSFixture(t)
		withUID(t, 1000)
		readOnlyPackageCache(t, f.pkgcacheDir(), "acme/util/0.1.0")
		_, err := newWSDaemon(t, Config{DBPath: f.db, WorkspacePackageCache: f.pkgcacheDir()})
		var startup *StartupError
		if !errors.As(err, &startup) || startup.Stage != StageConfig || !strings.Contains(err.Error(), "need --workspace-root") {
			t.Fatalf("New = %v, want a %s StartupError containing %q", err, StageConfig, "need --workspace-root")
		}
		if _, err := os.Stat(f.db); !os.IsNotExist(err) {
			t.Fatalf("the store was opened before the refusal (stat: %v)", err)
		}
	})
}

func TestWorkspacePackageCacheLinkedIntoEpisodeHome(t *testing.T) {
	f := newWSFixture(t)
	withUID(t, 1000)
	readOnlyPackageCache(t, f.pkgcacheDir(), "acme/util/0.1.0")
	d := mustWSDaemon(t, f.cfgWithCache(t, f.pkgcacheDir(), nil))
	if len(d.workspace.registry("ep1")) == 0 {
		t.Fatal("registry(ep1) is empty")
	}
	link := f.registryLink()
	fi, err := os.Lstat(link)
	if err != nil || fi.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("Lstat(%s) = %v, %v, want a symlink", link, fi, err)
	}
	if target, err := os.Readlink(link); err != nil || target != d.workspace.packageCache || target != f.pkgcacheDir() {
		t.Fatalf("Readlink = %q, %v, want %s", target, err, f.pkgcacheDir())
	}
	env, err := os.ReadFile(filepath.Join(f.logDir, "env"))
	if err != nil {
		t.Fatal(err)
	}
	if want := "HOME=" + filepath.Join(f.stateDir, "cache", "ep1") + "\n"; !strings.Contains(string(env), want) {
		t.Fatalf("tool env lacks %q:\n%s", want, env)
	}
	got, _ := os.ReadFile(filepath.Join(f.stateDir, "policies", "ep1.toml"))
	if want, _ := broker.RenderEpisodePolicy(filepath.Join(f.root, "ep1")); string(got) != string(want) {
		t.Fatalf("policy bytes changed with the cache flag:\n%s\nwant\n%s", got, want)
	}
	// The keep arm: a symlink already pointing at the snapshot is left as is.
	if err := d.workspace.linkPackageCache("ep1", filepath.Join(f.stateDir, "cache", "ep1")); err != nil {
		t.Fatalf("second linkPackageCache = %v, want nil", err)
	}
	if target, _ := os.Readlink(link); target != f.pkgcacheDir() {
		t.Fatalf("link changed to %q", target)
	}
}

func TestWorkspacePackageCacheRefusesUnprovisionedRegistry(t *testing.T) {
	type plant struct {
		name  string
		setup func(t *testing.T, f wsFixture, link string)
		check func(t *testing.T, f wsFixture, link string)
	}
	plants := []plant{
		{"non-empty dir", func(t *testing.T, f wsFixture, link string) {
			mkdirs(t, filepath.Join(link, "sunholo", "x", "9.9.9"))
			writeFile(t, filepath.Join(link, "sunholo", "x", "9.9.9", "ailang.toml"), "planted\n")
		}, func(t *testing.T, f wsFixture, link string) {
			if b, err := os.ReadFile(filepath.Join(link, "sunholo", "x", "9.9.9", "ailang.toml")); err != nil || string(b) != "planted\n" {
				t.Fatalf("planted file after the refusal = %q, %v", b, err)
			}
		}},
		{"regular file", func(t *testing.T, f wsFixture, link string) { writeFile(t, link, "planted\n") },
			func(t *testing.T, f wsFixture, link string) {
				if b, err := os.ReadFile(link); err != nil || string(b) != "planted\n" {
					t.Fatalf("planted file after the refusal = %q, %v", b, err)
				}
			}},
		{"foreign symlink", func(t *testing.T, f wsFixture, link string) {
			if err := os.Symlink(f.outside, link); err != nil {
				t.Fatal(err)
			}
		}, func(t *testing.T, f wsFixture, link string) {
			if target, err := os.Readlink(link); err != nil || target != f.outside {
				t.Fatalf("foreign link after the refusal = %q, %v", target, err)
			}
		}},
	}
	for _, tc := range plants {
		t.Run(tc.name, func(t *testing.T) {
			f := newWSFixture(t)
			withUID(t, 1000)
			readOnlyPackageCache(t, f.pkgcacheDir(), "acme/util/0.1.0")
			var log bytes.Buffer
			d := mustWSDaemon(t, f.cfgWithCache(t, f.pkgcacheDir(), &log))
			link := f.registryLink()
			mkdirs(t, filepath.Dir(link))
			tc.setup(t, f, link)
			log.Reset() // drop the startup line: assert on what the episode prints
			configureStubExec(t, d, f)
			reg := d.workspace.registry("ep1")
			requireOnlyExec(t, reg)
			want := fmt.Sprintf("ailang-worldd: workspace package cache refused for episode \"ep1\": %s is not empty and was not provisioned by "+
				"--workspace-package-cache; clear it manually\n", link)
			if tc.name == "foreign symlink" {
				// N6: say what the link is, since the operator may have moved the snapshot.
				want = fmt.Sprintf("ailang-worldd: workspace package cache refused for episode \"ep1\": %s is a symlink to %s, "+
					"not to --workspace-package-cache %s; remove it if the snapshot moved\n", link, f.outside, f.pkgcacheDir())
			}
			if log.String() != want {
				t.Fatalf("operator log = %q, want exactly %q", log.String(), want)
			}
			if n := f.summaries(t); n != 0 {
				t.Fatalf("a refused episode ran the tool (%d summaries)", n)
			}
			tc.check(t, f, link)
			requireExecRunsIn(t, reg, filepath.Join(f.root, "ep1"))
		})
	}
	t.Run("empty dir is replaced", func(t *testing.T) {
		f := newWSFixture(t)
		withUID(t, 1000)
		readOnlyPackageCache(t, f.pkgcacheDir(), "acme/util/0.1.0")
		d := mustWSDaemon(t, f.cfgWithCache(t, f.pkgcacheDir(), nil))
		mkdirs(t, f.registryLink())
		if len(d.workspace.registry("ep1")) == 0 {
			t.Fatal("registry(ep1) is empty with an empty registry dir in place")
		}
		if target, err := os.Readlink(f.registryLink()); err != nil || target != f.pkgcacheDir() {
			t.Fatalf("Readlink = %q, %v, want the snapshot", target, err)
		}
	})
}

func TestWorkspacePackageCacheUnsetLeavesHomeAlone(t *testing.T) {
	for _, tc := range []struct {
		name  string
		setup func(t *testing.T, f wsFixture, link string)
		check func(t *testing.T, f wsFixture, link string)
	}{
		{"non-empty dir untouched", func(t *testing.T, f wsFixture, link string) {
			mkdirs(t, link)
			writeFile(t, filepath.Join(link, "keep.txt"), "kept\n")
		}, func(t *testing.T, f wsFixture, link string) {
			if b, err := os.ReadFile(filepath.Join(link, "keep.txt")); err != nil || string(b) != "kept\n" {
				t.Fatalf("keep.txt = %q, %v", b, err)
			}
		}},
		{"empty dir untouched", func(t *testing.T, f wsFixture, link string) { mkdirs(t, link) },
			func(t *testing.T, f wsFixture, link string) {
				if fi, err := os.Lstat(link); err != nil || !fi.IsDir() {
					t.Fatalf("empty registry dir after startup: %v, %v", fi, err)
				}
			}},
		{"missing registry stays missing", func(t *testing.T, f wsFixture, link string) {},
			func(t *testing.T, f wsFixture, link string) {
				if _, err := os.Lstat(link); !os.IsNotExist(err) {
					t.Fatalf("flag unset: the absent registry path was created (lstat: %v)", err)
				}
			}},
		{"stale symlink removed, target intact", func(t *testing.T, f wsFixture, link string) {
			target := filepath.Join(f.outside, "target")
			mkdirs(t, target)
			writeFile(t, filepath.Join(target, "keep.txt"), "kept\n")
			if err := os.Symlink(target, link); err != nil {
				t.Fatal(err)
			}
		}, func(t *testing.T, f wsFixture, link string) {
			if _, err := os.Lstat(link); !os.IsNotExist(err) {
				t.Fatalf("stale symlink still present (lstat: %v)", err)
			}
			if b, err := os.ReadFile(filepath.Join(f.outside, "target", "keep.txt")); err != nil || string(b) != "kept\n" {
				t.Fatalf("the symlink's target was touched: %q, %v", b, err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newWSFixture(t)
			d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease)})
			link := f.registryLink()
			mkdirs(t, filepath.Dir(link))
			tc.setup(t, f, link)
			if len(d.workspace.registry("ep1")) == 0 {
				t.Fatal("registry(ep1) is empty")
			}
			tc.check(t, f, link)
		})
	}
}

// fakeToolBinWith is fakeToolBin with two changes: a policy-tool dispatch
// prints dispatchJSON (when non-empty), and `run` answers outcome shape (a) of
// composeRunOutcome (an admitted run, rc 0), so Ailang.Run reaches executeRun's
// stamp without a real binary.
func fakeToolBinWith(t *testing.T, logDir, dispatchJSON string) string {
	t.Helper()
	dispatch := `printf '{"ok":true,"content":"fake:%s"}' "$(pwd -P)"`
	if dispatchJSON != "" {
		dispatch = fmt.Sprintf("printf '%%s' '%s'", dispatchJSON)
	}
	script := fmt.Sprintf(`#!/bin/sh
log=%q
if [ "$1" = "--version" ]; then echo %q; exit 0; fi
if [ "$1" = "policy-tool" ]; then
  req=$(cat)
  case "$req" in
    *'"summary"'*)
      echo x >> "$log/summaries"
      env > "$log/env"
      printf '{"ok":true,"summary":{"security_mode":"restricted","policy_digest":"fakedigest","fs_sandbox":"%%s","cli":["check"]}}' "$(pwd -P)"
      exit 0 ;;
  esac
  echo x >> "$log/dispatches"
  %s
  exit 0
fi
if [ "$1" = "run" ]; then echo 'policy: {"ok":true,"decision":{"ok":true}}' >&2; exit 0; fi
echo x >> "$log/dispatches"
exit 0
`, logDir, ToolBinaryRelease, dispatch)
	bin := filepath.Join(canonicalTemp(t), "ailang")
	if err := os.WriteFile(bin, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	return bin
}

// ownDigest recomputes the snapshot digest with its own recursive ReadDir
// walker and its own sort (the P-D5 byte form), sharing nothing with
// resolvePackageCache.
func ownDigest(t *testing.T, dir string) string {
	t.Helper()
	var records []string
	var walk func(rel string)
	walk = func(rel string) {
		entries, err := os.ReadDir(filepath.Join(dir, rel))
		if err != nil {
			t.Fatal(err)
		}
		for _, e := range entries {
			r := e.Name()
			if rel != "" {
				r = rel + "/" + e.Name()
			}
			if e.IsDir() {
				records = append(records, r+"\x00d")
				walk(r)
				continue
			}
			b, err := os.ReadFile(filepath.Join(dir, r))
			if err != nil {
				t.Fatal(err)
			}
			sum := sha256.Sum256(b)
			records = append(records, r+"\x00f\x00"+hex.EncodeToString(sum[:]))
		}
	}
	walk("")
	sort.Strings(records)
	all := strings.Join(records, "\n") + "\n"
	sum := sha256.Sum256([]byte(all))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func execJSON(t *testing.T, d *Daemon, episode, effect, payload string) (map[string]any, error) {
	t.Helper()
	h := d.workspace.registry(episode)[effect]
	if h == nil {
		t.Fatalf("registry(%s) has no %s handler", episode, effect)
	}
	out, err := h.Execute(boundedTestContext(t), broker.EffectRequest{Effect: effect, Scope: broker.WorkspaceScope, Cost: 1}, []byte(payload))
	if err != nil {
		return nil, err
	}
	var resp map[string]any
	if err := json.Unmarshal(out, &resp); err != nil {
		t.Fatalf("%s: %v", out, err)
	}
	return resp, nil
}

func TestWorkspacePackageCacheDigestIsStamped(t *testing.T) {
	const readPayload, runPayload = `{"op":"read","path":"x"}`, `{"path":"r.ail"}`
	f := newWSFixture(t)
	withUID(t, 1000)
	// util and util-x: walk order differs from byte order of the full relpath.
	readOnlyPackageCache(t, f.pkgcacheDir(), "acme/util/0.1.0", "acme/util-x/0.1.0")
	cfg := Config{DBPath: f.db, WorkspaceRoot: f.root, WorkspacePackageCache: f.pkgcacheDir(),
		ToolAilangBin: fakeToolBinWith(t, f.logDir, "")}
	d := mustWSDaemon(t, cfg)
	want := ownDigest(t, f.pkgcacheDir())
	if d.workspace.packageCacheDigest != want {
		t.Fatalf("startup digest = %q, the test's own recomputation = %q", d.workspace.packageCacheDigest, want)
	}
	// (a) a policy-tool result names the snapshot.
	resp, err := execJSON(t, d, "ep1", broker.EffectWorkspaceRead, readPayload)
	if err != nil || resp["package_cache"] != want {
		t.Fatalf("Workspace.Read package_cache = %v (%v), want %s", resp["package_cache"], err, want)
	}
	// (c) and so does an ailang-run result.
	resp, err = execJSON(t, d, "ep1", broker.EffectAilangRun, runPayload)
	if err != nil || resp["admitted"] != true || resp["package_cache"] != want {
		t.Fatalf("Ailang.Run result = %v (%v), want admitted with package_cache %s", resp, err, want)
	}

	// (b) and (c, unset arm): with no snapshot the key is absent.
	f2 := newWSFixture(t)
	d2 := mustWSDaemon(t, Config{DBPath: f2.db, WorkspaceRoot: f2.root, ToolAilangBin: fakeToolBinWith(t, f2.logDir, "")})
	for effect, payload := range map[string]string{broker.EffectWorkspaceRead: readPayload, broker.EffectAilangRun: runPayload} {
		resp, err := execJSON(t, d2, "ep1", effect, payload)
		if err != nil {
			t.Fatalf("%s: %v", effect, err)
		}
		if _, present := resp["package_cache"]; present {
			t.Fatalf("%s result carries package_cache with the flag unset: %v", effect, resp)
		}
	}

	// (d) the key is reserved: a tool response carrying it is refused even
	// with no snapshot, so it can never pass for World provenance.
	f3 := newWSFixture(t)
	d3 := mustWSDaemon(t, Config{DBPath: f3.db, WorkspaceRoot: f3.root,
		ToolAilangBin: fakeToolBinWith(t, f3.logDir, `{"ok":true,"package_cache":"forged"}`)})
	if _, err := execJSON(t, d3, "ep1", broker.EffectWorkspaceRead, readPayload); err == nil || !strings.Contains(err.Error(), "package_cache") {
		t.Fatalf("a forged package_cache key = %v, want an error naming the reserved key", err)
	}
}

func TestWorkspacePackageCacheCoversLock(t *testing.T) {
	util := lockEntry{"acme/util", "0.1.0", "registry"}
	for _, tc := range []struct {
		name  string
		setup func(t *testing.T, f wsFixture, lock string)
		want  string // "" = the registry forms; else the exact operator line (with %DIR%) or a prefix ending in "lock"
		exact bool
	}{
		{"covered", func(t *testing.T, f wsFixture, lock string) { writeLock(t, lock, util) }, "", true},
		{"one registry entry absent", func(t *testing.T, f wsFixture, lock string) {
			writeLock(t, lock, util, lockEntry{"sunholo/oauth", "0.1.0", "registry"})
		}, `ailang-worldd: workspace package cache does not cover episode "ep1": lock requires sunholo/oauth@0.1.0, absent from %DIR%; ` +
			`rebuild the snapshot per QUICKSTART §11` + "\n", true},
		{"path-source entry is not checked", func(t *testing.T, f wsFixture, lock string) {
			writeLock(t, lock, util, lockEntry{"local/thing", "0.1.0", "path"})
		}, "", true},
		{"no lock", func(t *testing.T, f wsFixture, lock string) {}, "", true},
		{"invalid JSON", func(t *testing.T, f wsFixture, lock string) { writeFile(t, lock, "{") },
			`does not cover episode "ep1": lock`, false},
		{"traversal in the name", func(t *testing.T, f wsFixture, lock string) {
			// The traversal WOULD succeed: DIR/../esc/1.0.0/ailang.toml exists.
			mkdirs(t, filepath.Join(f.base, "esc", "1.0.0"))
			writeFile(t, filepath.Join(f.base, "esc", "1.0.0", "ailang.toml"), "[package]\n")
			writeLock(t, lock, lockEntry{"../esc", "1.0.0", "registry"})
		}, `does not cover episode "ep1": lock requires ../esc@1.0.0`, false},
		{"lock over 1 MiB", func(t *testing.T, f wsFixture, lock string) {
			// Valid JSON that WOULD be covered, so only the size bound refuses it.
			writeFile(t, lock, `{"packages":[],"pad":"`+strings.Repeat("a", 1<<20)+`"}`)
		}, `is larger than 1 MiB; rebuild the snapshot per QUICKSTART §11`, false},
		{"ailang.toml of a locked package is a directory", func(t *testing.T, f wsFixture, lock string) {
			mkdirs(t, filepath.Join(f.pkgcacheDir(), "acme", "dirpkg", "0.1.0", "ailang.toml"))
			writeLock(t, lock, util, lockEntry{"acme/dirpkg", "0.1.0", "registry"})
		}, `lock requires acme/dirpkg@0.1.0`, false},
		{"lock is a symlink", func(t *testing.T, f wsFixture, lock string) {
			real := filepath.Join(f.outside, "real.lock")
			writeLock(t, real, util)
			if err := os.Symlink(real, lock); err != nil {
				t.Fatal(err)
			}
		}, `does not cover episode "ep1": lock`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newWSFixture(t)
			withUID(t, 1000)
			buildPackageCache(t, f.pkgcacheDir(), "acme/util/0.1.0")
			mkdirs(t, filepath.Join(f.root, "ep1", "tools"))
			tc.setup(t, f, filepath.Join(f.root, "ep1", "tools", "ailang.lock"))
			lockdown(t, f.pkgcacheDir())
			var log bytes.Buffer
			cfg := f.cfgWithCache(t, f.pkgcacheDir(), &log)
			cfg.WorkspaceModuleRoot = "tools"
			d := mustWSDaemon(t, cfg)
			if tc.want != "" {
				configureStubExec(t, d, f)
			}
			log.Reset()
			reg := d.workspace.registry("ep1")
			if tc.want == "" {
				if len(reg) == 0 || log.Len() != 0 {
					t.Fatalf("registry = %v, log %q, want a registry and a silent log", registryNames(reg), log.String())
				}
				return
			}
			requireOnlyExec(t, reg)
			if f.summaries(t) != 0 {
				t.Fatalf("registry = %v, summaries %d, want only Workspace.Exec and 0", registryNames(reg), f.summaries(t))
			}
			want := strings.ReplaceAll(tc.want, "%DIR%", f.pkgcacheDir())
			if n := strings.Count(log.String(), "\n"); n != 1 || (tc.exact && log.String() != want) || (!tc.exact && !strings.Contains(log.String(), want)) {
				t.Fatalf("operator log (%d lines) = %q, want %q", n, log.String(), want)
			}
			// Exec runs in the worktree, not in the module root (<root>/ep1/tools).
			requireExecRunsIn(t, reg, filepath.Join(f.root, "ep1"))
		})
	}
}

// TestWorkspacePackageCacheLockFIFOIsRefusedPromptly (N1): ailang.lock is
// agent-writable, and opening a FIFO for reading blocks until a writer
// appears. The lock read must refuse it promptly instead of hanging episode
// construction while the handler mutex is held.
func TestWorkspacePackageCacheLockFIFOIsRefusedPromptly(t *testing.T) {
	f := newWSFixture(t)
	withUID(t, 1000)
	readOnlyPackageCache(t, f.pkgcacheDir(), "acme/util/0.1.0")
	mkdirs(t, filepath.Join(f.root, "ep1", "tools"))
	lock := filepath.Join(f.root, "ep1", "tools", "ailang.lock")
	if err := syscall.Mkfifo(lock, 0o644); err != nil {
		t.Fatal(err)
	}
	var log bytes.Buffer
	cfg := f.cfgWithCache(t, f.pkgcacheDir(), &log)
	cfg.WorkspaceModuleRoot = "tools"
	d := mustWSDaemon(t, cfg)
	configureStubExec(t, d, f)
	log.Reset()
	done := make(chan broker.Registry, 1)
	go func() { done <- d.workspace.registry("ep1") }()
	select {
	case reg := <-done:
		if names := registryNames(reg); strings.Join(names, ",") != broker.EffectWorkspaceExec ||
			!strings.Contains(log.String(), `does not cover episode "ep1": lock `) || !strings.Contains(log.String(), "is not a regular file") {
			t.Fatalf("registry %v, log %q, want only Workspace.Exec and a not-a-regular-file refusal", names, log.String())
		}
		requireExecRunsIn(t, reg, filepath.Join(f.root, "ep1"))
	case <-time.After(2 * time.Second):
		// Release the blocked open so the goroutine and the daemon's Close can finish.
		if w, err := os.OpenFile(lock, os.O_RDWR|syscall.O_NONBLOCK, 0); err == nil {
			_ = w.Close()
		}
		t.Fatal("a FIFO as ailang.lock hung registry construction (no refusal within 2s)")
	}
}
