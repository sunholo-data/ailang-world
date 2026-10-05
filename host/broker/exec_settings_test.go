package broker

// Row 140 M2 AC2.1 (w-workspace-exec-toolchain-effect §4.4): the rendered srt
// settings' exact key set, a golden per platform for one profile, and the
// verify-at-render checks.

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

var updateExecGolden = flag.Bool("update-exec-settings-golden", false, "rewrite testdata/exec_settings_*.golden.json")

func goldenSettingsSpec(goos string) execSettingsSpec {
	return execSettingsSpec{
		Worktree:      "/W/workspace/ep1",
		Cache:         "/W/state/exec-cache/ep1/ailang-compiler",
		OperatorHome:  "/Users/op",
		StateDir:      "/W/state",
		WorkspaceRoot: "/W/workspace",
		ReadRoots:     []string{"/Users/op/go/pkg/mod", "/opt/homebrew/Cellar/go/1.26.4/libexec", "/Users/op/go/pkg/mod"},
		SeccompDir:    "/W/state/archive/exec-sandbox/d/node_modules/@anthropic-ai/sandbox-runtime/vendor/seccomp/x64",
		GOOS:          goos,
	}
}

func keysOf(m map[string]any) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

func TestRenderExecSettingsExactKeySet(t *testing.T) {
	for _, goos := range []string{"darwin", "linux"} {
		t.Run(goos, func(t *testing.T) {
			data, err := renderExecSettings(goldenSettingsSpec(goos))
			if err != nil {
				t.Fatal(err)
			}
			var m map[string]any
			if err := json.Unmarshal(data, &m); err != nil {
				t.Fatal(err)
			}
			if got, want := keysOf(m), []string{"allowAppleEvents", "enableWeakerNetworkIsolation", "filesystem", "network"}; !reflect.DeepEqual(got, want) {
				t.Fatalf("top-level keys = %v, want %v", got, want)
			}
			fs := m["filesystem"].(map[string]any)
			if got, want := keysOf(fs), []string{"allowRead", "allowWrite", "denyRead", "denyWrite"}; !reflect.DeepEqual(got, want) {
				t.Fatalf("filesystem keys = %v, want %v", got, want)
			}
			net := m["network"].(map[string]any)
			if got, want := keysOf(net), []string{"allowAllUnixSockets", "allowLocalBinding", "allowedDomains", "deniedDomains"}; !reflect.DeepEqual(got, want) {
				t.Fatalf("network keys = %v, want %v", got, want)
			}
			for k, v := range map[string]any{"allowAppleEvents": m["allowAppleEvents"], "enableWeakerNetworkIsolation": m["enableWeakerNetworkIsolation"],
				"allowLocalBinding": net["allowLocalBinding"], "allowAllUnixSockets": net["allowAllUnixSockets"]} {
				if v != false {
					t.Fatalf("%s = %v, want an explicit false", k, v)
				}
			}
			for _, k := range []string{"allowedDomains", "deniedDomains"} {
				if l, ok := net[k].([]any); !ok || len(l) != 0 {
					t.Fatalf("%s = %v, want an explicit [] (D-140-2 = A)", k, net[k])
				}
			}
			spec := goldenSettingsSpec(goos)
			strs := func(k string) []string {
				var out []string
				for _, v := range fs[k].([]any) {
					out = append(out, v.(string))
				}
				return out
			}
			if got, want := strs("denyRead"), []string{spec.OperatorHome, spec.StateDir, spec.WorkspaceRoot}; !reflect.DeepEqual(got, want) {
				t.Fatalf("denyRead = %v, want %v (D-140-1 = A)", got, want)
			}
			if got, want := strs("allowWrite"), []string{spec.Worktree, spec.Cache}; !reflect.DeepEqual(got, want) {
				t.Fatalf("allowWrite = %v, want %v", got, want)
			}
			allowRead := strs("allowRead")
			if hasSeccomp := slicesContains(allowRead, spec.SeccompDir); hasSeccomp != (goos == "linux") {
				t.Fatalf("allowRead %v: the seccomp helper dir must be re-allowed on linux only [M0, V43/V44]", allowRead)
			}
			denyWrite := strs("denyWrite")
			for _, name := range []string{".git", ".github", ".pi", ".claude", ".ailang", ".gitmodules", ".gitattributes"} {
				if !slicesContains(denyWrite, filepath.Join(spec.Worktree, name)) {
					t.Fatalf("denyWrite %v lacks World's %s", denyWrite, name)
				}
			}
			for _, p := range []string{"/tmp/claude", spec.OperatorHome + "/.npm/_logs", spec.OperatorHome + "/.claude/debug"} {
				if !slicesContains(denyWrite, p) {
					t.Fatalf("denyWrite %v lacks srt's always-writable %s (V4)", denyWrite, p)
				}
			}
			if slicesContains(denyWrite, "/private/tmp/claude") != (goos == "darwin") {
				t.Fatalf("denyWrite %v: /private/tmp/claude is darwin's alias only", denyWrite)
			}
		})
	}
}

func slicesContains(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

func TestRenderExecSettingsGolden(t *testing.T) {
	for _, goos := range []string{"darwin", "linux"} {
		t.Run(goos, func(t *testing.T) {
			data, err := renderExecSettings(goldenSettingsSpec(goos))
			if err != nil {
				t.Fatal(err)
			}
			golden := filepath.Join("testdata", "exec_settings_"+goos+".golden.json")
			if *updateExecGolden {
				if err := os.WriteFile(golden, data, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			want, err := os.ReadFile(golden)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(data, want) {
				t.Fatalf("rendered settings differ from %s:\n%s\nwant:\n%s", golden, data, want)
			}
		})
	}
}

func TestRenderExecSettingsRefusesUncleanPaths(t *testing.T) {
	spec := goldenSettingsSpec("linux")
	spec.Worktree = "/W/workspace/ep1/../ep2"
	if _, err := renderExecSettings(spec); err == nil {
		t.Fatal("an unclean worktree path was rendered")
	}
	spec = goldenSettingsSpec("linux")
	spec.ReadRoots = []string{"relative"}
	if _, err := renderExecSettings(spec); err == nil {
		t.Fatal("a relative read root was rendered")
	}
}

// Verify-at-render (§4.4): World re-reads the file, byte-compares it with
// its own canonical rendering and refuses a file it did not write.
func TestWriteAndVerifyExecSettings(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ep1.proj.srt.json")
	data, err := renderExecSettings(goldenSettingsSpec(hostGOOS()))
	if err != nil {
		t.Fatal(err)
	}
	digest, err := writeExecSettings(path, data)
	if err != nil {
		t.Fatalf("writeExecSettings: %v", err)
	}
	if want := "sha256:" + sha256Hex(data); digest != want {
		t.Fatalf("settings_digest = %q, want %q", digest, want)
	}
	info, err := os.Lstat(path)
	if err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("settings file mode = %v (%v), want 0600", info.Mode(), err)
	}
	if err := verifyExecSettings(path, data); err != nil {
		t.Fatalf("verifyExecSettings on World's own file: %v", err)
	}

	// Not World's bytes.
	if err := os.WriteFile(path, append(append([]byte(nil), data...), ' '), 0o600); err != nil {
		t.Fatal(err)
	}
	var unverified *ExecSettingsError
	if err := verifyExecSettings(path, data); !errors.As(err, &unverified) || !strings.Contains(err.Error(), path) {
		t.Fatalf("a changed file verified: %v", err)
	}
	// World's bytes behind a symlink.
	other := filepath.Join(dir, "other.json")
	if err := os.WriteFile(other, data, 0o600); err != nil {
		t.Fatal(err)
	}
	_ = os.Remove(path)
	if err := os.Symlink(other, path); err != nil {
		t.Fatal(err)
	}
	if err := verifyExecSettings(path, data); err == nil {
		t.Fatal("a symlink to the right bytes verified")
	}
	// World's bytes at a widened mode.
	_ = os.Remove(path)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := verifyExecSettings(path, data); err == nil {
		t.Fatal("a 0644 settings file verified")
	}
	_ = os.Remove(path)
	if err := verifyExecSettings(path, data); err == nil {
		t.Fatal("a missing settings file verified")
	}
}
