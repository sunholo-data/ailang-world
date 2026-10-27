package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/toolprobe"
)

func hostsInitTest(t *testing.T, dir string, extra ...string) (int, string) {
	t.Helper()
	var out bytes.Buffer
	code := runHosts(append([]string{"init", "--out-dir", dir}, extra...), &out, &out)
	return code, out.String()
}

func TestHostsInitWritesEveryFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if err := os.MkdirAll(filepath.Join(home, ".codex"), 0700); err != nil {
		t.Fatal(err)
	}
	auth := filepath.Join(home, ".codex", "auth.json")
	// Unreadable credentials still only need a metadata lookup, never a read.
	if err := os.WriteFile(auth, []byte("secret-must-not-be-copied"), 0000); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(t.TempDir(), "hosts")
	if code, out := hostsInitTest(t, dir); code != 0 {
		t.Fatalf("%d: %s", code, out)
	}
	files := hostFiles(dir, "http://127.0.0.1:47330")
	// Independent wire-format assertions protect the measured host flags.
	wantSettings := "{\"permissions\":{\"deny\":[\"Bash\",\"Edit\",\"Write\",\"Read\",\"Glob\",\"Grep\",\"NotebookEdit\",\"WebFetch\",\"WebSearch\"],\"disableBypassPermissionsMode\":\"disable\"}}\n"
	var settings map[string]any
	if json.Unmarshal(files["claude-world.settings.json"], &settings) != nil {
		t.Fatal("invalid settings")
	}
	delete(settings, "statusLine")
	settingsBytes, _ := json.Marshal(settings)
	if string(settingsBytes)+"\n" != wantSettings {
		t.Fatal("settings bytes differ")
	}
	wantCodex := "sandbox_mode = \"read-only\"\napproval_policy = \"never\"\nweb_search = \"disabled\"\n[features]\nhooks = true\nshell_tool = false\nunified_exec = false\n[mcp_servers.world]\nurl = \"http://127.0.0.1:47330/mcp/\"\nbearer_token_env_var = \"WORLD_SESSION\"\nrequired = true\ndefault_tools_approval_mode = \"approve\"\n"
	if string(files["codex/config.toml"]) != wantCodex {
		t.Fatal("Codex bytes differ")
	}
	for name, want := range map[string]string{
		"claude-world.mcp.json": "{\"mcpServers\":{\"world\":{\"headers\":{\"Authorization\":\"Bearer ${WORLD_SESSION}\"},\"type\":\"http\",\"url\":\"http://127.0.0.1:47330/mcp/\"}}}\n",
		"pi-world-mcp.json":     "{\"mcpServers\":{\"world\":{\"auth\":\"bearer\",\"bearerTokenEnv\":\"WORLD_SESSION\",\"lifecycle\":\"eager\",\"url\":\"http://127.0.0.1:47330/mcp/\"}}}\n",
	} {
		if string(files[name]) != want {
			t.Errorf("%s bytes: %s", name, files[name])
		}
	}
	b, err := os.ReadFile(filepath.Join(dir, "hosts.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest hostsManifest
	if err := json.Unmarshal(b, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Schema != "world/hosts/v1" || manifest.Addr != "http://127.0.0.1:47330" || len(manifest.Files) != len(files) {
		t.Fatalf("manifest: %+v", manifest)
	}
	for name, want := range files {
		got, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil || !bytes.Equal(got, want) {
			t.Errorf("%s bytes: %q (%v)", name, got, err)
		}
		info, err := os.Stat(filepath.Join(dir, name))
		mode := os.FileMode(0600)
		if strings.HasPrefix(name, "bin/") {
			mode = 0700
		}
		if err != nil || info.Mode().Perm() != mode {
			t.Errorf("%s mode: %v, %v", name, info, err)
		}
		if manifest.Files[name] != "sha256:"+fmtHash(got) {
			t.Errorf("%s digest", name)
		}
	}
	for _, name := range []string{"", "bin", "codex"} {
		info, err := os.Stat(filepath.Join(dir, name))
		if err != nil || info.Mode().Perm() != 0700 {
			t.Errorf("%s directory mode", name)
		}
	}
	info, err := os.Stat(filepath.Join(dir, "hosts.json"))
	if err != nil || info.Mode().Perm() != 0600 {
		t.Fatal("manifest mode")
	}
	target, err := os.Readlink(filepath.Join(dir, "codex", "auth.json"))
	if err != nil || target != auth {
		t.Fatalf("auth link %q: %v", target, err)
	}
	entries := 0
	filepath.WalkDir(dir, func(_ string, d os.DirEntry, err error) error {
		if err != nil {
			t.Fatal(err)
		}
		if !d.IsDir() {
			entries++
		}
		return nil
	})
	if entries != len(files)+2 {
		t.Fatalf("wrote %d files, want %d including manifest and symlink", entries, len(files)+2)
	}
	if err := os.Remove(auth); err != nil {
		t.Fatal(err)
	}
	if code, out := hostsInitTest(t, filepath.Join(t.TempDir(), "no-auth")); code != 0 || !strings.Contains(out, "! Codex will ask you to log in") {
		t.Fatalf("%d: %s", code, out)
	}
}

func fmtHash(b []byte) string {
	const hex = "0123456789abcdef"
	sum := sha256.Sum256(b)
	out := make([]byte, 64)
	for i, v := range sum {
		out[2*i], out[2*i+1] = hex[v>>4], hex[v&15]
	}
	return string(out)
}

func TestHostsInitRefusesChangedFileWithoutReplace(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir := filepath.Join(t.TempDir(), "hosts with 'quote")
	if code, out := hostsInitTest(t, dir); code != 0 {
		t.Fatalf("%d: %s", code, out)
	}
	p := filepath.Join(dir, "codex", "config.toml")
	before, _ := os.Stat(p)
	if code, out := hostsInitTest(t, dir); code != 0 {
		t.Fatalf("no-op: %d %s", code, out)
	}
	after, _ := os.Stat(p)
	if !before.ModTime().Equal(after.ModTime()) {
		t.Fatal("unchanged file rewritten")
	}
	if err := os.WriteFile(p, []byte("changed"), 0600); err != nil {
		t.Fatal(err)
	}
	if code, out := hostsInitTest(t, dir); code == 0 || !strings.Contains(out, p) {
		t.Fatalf("refusal: %d %s", code, out)
	}
	if code, out := hostsInitTest(t, dir, "--replace"); code != 0 {
		t.Fatalf("replace: %d %s", code, out)
	}
	got, _ := os.ReadFile(p)
	if !bytes.Equal(got, hostFiles(dir, "http://127.0.0.1:47330")["codex/config.toml"]) {
		t.Fatal("replace did not restore config")
	}
}

func hostsLauncherFixture(t *testing.T) (string, string) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	writeTestPiAdapter(t, defaultPiMCPAdapter())
	root := t.TempDir()
	dir := filepath.Join(root, "host's configs")
	if code, out := hostsInitTest(t, dir); code != 0 {
		t.Fatalf("%d: %s", code, out)
	}
	bin := filepath.Join(root, "fake-bin")
	if err := os.Mkdir(bin, 0700); err != nil {
		t.Fatal(err)
	}
	for _, host := range []string{"claude", "codex", "pi"} {
		script := "#!/bin/sh\nprintf '%s\\n' \"$CODEX_HOME\" \"$@\"\n"
		if err := os.WriteFile(filepath.Join(bin, host), []byte(script), 0700); err != nil {
			t.Fatal(err)
		}
	}
	return dir, bin
}

func TestHostsLaunchersRefuseWithoutSession(t *testing.T) {
	dir, bin := hostsLauncherFixture(t)
	marker := filepath.Join(t.TempDir(), "host-ran")
	for _, host := range []string{"claude", "codex", "pi"} {
		if err := os.WriteFile(filepath.Join(bin, host), []byte("#!/bin/sh\ntouch "+shellQuote(marker)+"\n"), 0700); err != nil {
			t.Fatal(err)
		}
		ctx, cancel := budgetContext(5 * time.Second)
		// AC10: launcher fixtures run via the same driven, scrubbed toolprobe site.
		_, err := toolprobe.OutputEnv(ctx, filepath.Join(dir, "bin", host+"-world"), "", []string{"PATH=" + bin + ":/usr/bin:/bin"})
		cancel()
		var exit *exec.ExitError
		if !errors.As(err, &exit) || exit.ExitCode() != 3 {
			t.Fatalf("%s: %v", host, err)
		}
		for _, want := range []string{host + "-world: no World session", "why: World agents need a session you approved.", "next: start one with", "export WORLD_SESSION="} {
			if !strings.Contains(string(exit.Stderr), want) {
				t.Errorf("%s lacks %q: %s", host, want, exit.Stderr)
			}
		}
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatal("host ran without session")
		}
	}
}

func TestHostsLaunchersExecHost(t *testing.T) {
	dir, bin := hostsLauncherFixture(t)
	for _, host := range []string{"claude", "codex", "pi"} {
		ctx, cancel := budgetContext(5 * time.Second)
		out, err := toolprobe.OutputEnv(ctx, filepath.Join(dir, "bin", host+"-world"), "", []string{"PATH=" + bin + ":/usr/bin:/bin", "WORLD_SESSION=test-session"}, "a b", "--flag")
		cancel()
		if err != nil {
			t.Fatalf("%s: %v", host, err)
		}
		want := []string{""}
		switch host {
		case "claude":
			want = append(want, "--plugin-dir", filepath.Join(dir, "claude-plugin"), "--settings", filepath.Join(dir, "claude-world.settings.json"), "--strict-mcp-config", "--mcp-config", filepath.Join(dir, "claude-world.mcp.json"))
		case "codex":
			want[0] = filepath.Join(dir, "codex")
			want = append(want, "-c", "developer_instructions=\"World context unavailable: daemon down; run ailang-worldd doctor and start World.\"")
		case "pi":
			want = append(want, "-e", filepath.Join(dir, "pi", "world-status.ts"), "--append-system-prompt", "World context unavailable: daemon down; run ailang-worldd doctor and start World.", "--no-builtin-tools", "--no-extensions", "-e", defaultPiMCPAdapter(), "--mcp-config", filepath.Join(dir, "pi-world-mcp.json"))
		}
		want = append(want, "a b", "--flag")
		got := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s argv/env: %#v, want %#v", host, got, want)
		}
	}
}

func writeTestPiAdapter(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"pi":{"extensions":["./index.ts"]}}`), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestPiAdapterManifestAndDoctor(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	adapter := filepath.Join(t.TempDir(), "adapter with 'quote")
	writeTestPiAdapter(t, adapter)
	dir := filepath.Join(t.TempDir(), "hosts")
	if code, out := hostsInitTest(t, dir, "--pi-mcp-adapter", adapter); code != 0 {
		t.Fatalf("%d: %s", code, out)
	}
	b, err := os.ReadFile(filepath.Join(dir, "hosts.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest hostsManifest
	if err := json.Unmarshal(b, &manifest); err != nil {
		t.Fatal(err)
	}
	pkg, _ := os.ReadFile(filepath.Join(adapter, "package.json"))
	if manifest.PiMCPAdapter != adapter || manifest.PiMCPAdapterSHA256 != "sha256:"+fmtHash(pkg) {
		t.Fatalf("adapter manifest: %+v", manifest)
	}
	var out bytes.Buffer
	r := &doctorReport{w: &out}
	doctorPiAdapter(r, dir)
	if r.fails != 0 {
		t.Fatal(out.String())
	}
	for _, changed := range []bool{true, false} {
		if changed {
			if err := os.WriteFile(filepath.Join(adapter, "package.json"), append(pkg, '\n'), 0600); err != nil {
				t.Fatal(err)
			}
		} else if err := os.Remove(filepath.Join(adapter, "package.json")); err != nil {
			t.Fatal(err)
		}
		out.Reset()
		r = &doctorReport{w: &out}
		doctorPiAdapter(r, dir)
		if r.classCode != 12 || r.fails != 1 || !strings.Contains(out.String(), piAdapterFix) {
			t.Fatalf("changed=%v: class %d: %s", changed, r.classCode, out.String())
		}
	}
}

func TestPiLauncherRefusesMissingAdapter(t *testing.T) {
	for _, installed := range []bool{false, true} {
		t.Run(map[bool]string{false: "missing at init", true: "removed after init"}[installed], func(t *testing.T) {
			dir, bin := hostsLauncherFixture(t)
			adapter := defaultPiMCPAdapter()
			if !installed {
				if err := os.Remove(filepath.Join(adapter, "package.json")); err != nil {
					t.Fatal(err)
				}
				if code, out := hostsInitTest(t, dir, "--replace"); code != 0 {
					t.Fatalf("%d: %s", code, out)
				}
			} else if err := os.Remove(filepath.Join(adapter, "package.json")); err != nil {
				t.Fatal(err)
			}
			ctx, cancel := budgetContext(5 * time.Second)
			defer cancel()
			_, err := toolprobe.OutputEnv(ctx, filepath.Join(dir, "bin", "pi-world"), "", []string{"PATH=" + bin + ":/usr/bin:/bin", "WORLD_SESSION=test-session"})
			var exit *exec.ExitError
			if !errors.As(err, &exit) || exit.ExitCode() != 3 {
				t.Fatalf("launcher: %v", err)
			}
			for _, want := range []string{piAdapterWhat(adapter), "why: " + piAdapterWhy, "next: " + piAdapterFix} {
				if !strings.Contains(string(exit.Stderr), want) {
					t.Errorf("missing %q: %s", want, exit.Stderr)
				}
			}
		})
	}
}
