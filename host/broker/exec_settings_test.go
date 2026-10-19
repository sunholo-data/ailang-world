package broker

// Row 140 M2 AC2.1 (w-workspace-exec-toolchain-effect §4.4): the rendered srt
// settings' exact key set, a golden per platform for one profile, and the
// verify-at-render checks.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestProtectedMatrixModeContract(t *testing.T) {
	// python3 is a CI prerequisite (plan "Environment and honest gates"): a
	// silent skip here made the contract vacuous, so absence is a failure.
	python, err := exec.LookPath("python3")
	if err != nil {
		t.Fatalf("python3 is required for the protected self-test: %v", err)
	}
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate test source")
	}
	script := filepath.Join(filepath.Dir(source), "../..", "design_docs/verification/world-row140-m0/m0_matrix.py")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, python, script, "--protected-self-test")
	cmd.Dir = t.TempDir()
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("protected self-test failed: %v\n%s", err, out)
	}
	names := []string{"ok", "empty", "missing-client", "outside-root", "none-not-written", "host-changed",
		"wrote-unchanged", "invalid-outcome", "bad-digest", "root-equals-allowed", "symlink-escape", "bad-token",
		"invalid-rc", "snapshot-keys", "partial-token", "wrote-refused-token"}
	census := regexp.MustCompile(`(?m)^self-test: ([0-9]+) cases, ([0-9]+) passed$`).FindAllStringSubmatch(string(out), -1)
	if len(census) != 1 || census[0][1] != strconv.Itoa(len(names)) || census[0][2] != strconv.Itoa(len(names)) {
		t.Fatalf("missing or incorrect self-test census:\n%s", out)
	}
	for _, name := range names {
		if !regexp.MustCompile(`(?m)^\[PASS\] ` + regexp.QuoteMeta(name) + `$`).Match(out) {
			t.Errorf("missing passing case %s:\n%s", name, out)
		}
	}
	t.Logf("%s", out)
}

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
			// [M2-linux] V50: .git/hooks and .git/config precede .git, or srt's stubs for them
			// land inside the read-only .git bind and bwrap refuses to start.
			gitFirst := []string{filepath.Join(spec.Worktree, ".git", "hooks"), filepath.Join(spec.Worktree, ".git", "config"),
				filepath.Join(spec.Worktree, ".git")}
			if len(denyWrite) < 3 || !reflect.DeepEqual(denyWrite[:3], gitFirst) {
				t.Fatalf("denyWrite %v must open with %v (V50)", denyWrite, gitFirst)
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

func protectedSettingsSpec(t *testing.T, goos string) execSettingsSpec {
	t.Helper()
	spec := goldenSettingsSpec(goos)
	var err error
	spec.Protected, err = ParseProtectedPatterns([]string{"tests/**", "*_test.go"})
	if err != nil {
		t.Fatal(err)
	}
	return spec
}

func TestExecSettingsProtectedRendered(t *testing.T) {
	for _, goos := range []string{"darwin", "linux"} {
		t.Run(goos, func(t *testing.T) {
			data, err := renderExecSettings(protectedSettingsSpec(t, goos))
			if err != nil {
				t.Fatal(err)
			}
			var settings execSettings
			if err := json.Unmarshal(data, &settings); err != nil {
				t.Fatal(err)
			}
			want := []string{"/W/workspace/ep1/tests"}
			if goos == "darwin" {
				want = append([]string{"/W/workspace/ep1/**/*_test.go"}, want...)
			}
			legacyBytes, err := os.ReadFile(filepath.Join("testdata", "exec_settings_"+goos+".golden.json"))
			if err != nil {
				t.Fatal(err)
			}
			var legacy execSettings
			if err := json.Unmarshal(legacyBytes, &legacy); err != nil {
				t.Fatal(err)
			}
			want = append(legacy.Filesystem.DenyWrite, want...)
			got := settings.Filesystem.DenyWrite
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("protected denyWrite suffix = %q, want %q", got, want)
			}
		})
	}
}

// Full protected bytes, independent of the production renderer and deny lists.
const protectedSettingsGolden = `{
 "filesystem": {
  "denyRead": [
   "/Users/op",
   "/W/state",
   "/W/workspace"
  ],
  "allowRead": [
   "/W/workspace/ep1",
   "/W/state/exec-cache/ep1/ailang-compiler",
   "/Users/op/go/pkg/mod",
   "/opt/homebrew/Cellar/go/1.26.4/libexec"@SECCOMP@
  ],
  "allowWrite": [
   "/W/workspace/ep1",
   "/W/state/exec-cache/ep1/ailang-compiler"
  ],
  "denyWrite": [
   "/W/workspace/ep1/.git/hooks",
   "/W/workspace/ep1/.git/config",
   "/W/workspace/ep1/.git",
   "/W/workspace/ep1/.github",
   "/W/workspace/ep1/.pi",
   "/W/workspace/ep1/.claude",
   "/W/workspace/ep1/.ailang",
   "/W/workspace/ep1/.gitmodules",
   "/W/workspace/ep1/.gitattributes",
   "/tmp/claude",
@PRIVATE@   "/Users/op/.npm/_logs",
   "/Users/op/.claude/debug",
@BASENAME@   "/W/workspace/ep1/tests"
  ]
 },
 "network": {
  "allowedDomains": [],
  "deniedDomains": [],
  "allowLocalBinding": false,
  "allowAllUnixSockets": false
 },
 "enableWeakerNetworkIsolation": false,
 "allowAppleEvents": false
}
`

func checkProtectedSettingsGolden(t *testing.T, goos string) {
	t.Helper()
	seccomp, private, basename := "", "", ""
	if goos == "darwin" {
		private = "   \"/private/tmp/claude\",\n"
		basename = "   \"/W/workspace/ep1/**/*_test.go\",\n"
	} else {
		seccomp = ",\n   \"/W/state/archive/exec-sandbox/d/node_modules/@anthropic-ai/sandbox-runtime/vendor/seccomp/x64\""
	}
	want := strings.NewReplacer("@SECCOMP@", seccomp, "@PRIVATE@", private, "@BASENAME@", basename).Replace(protectedSettingsGolden)
	got, err := renderExecSettings(protectedSettingsSpec(t, goos))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatalf("protected %s golden bytes/order differ:\ngot:\n%s\nwant:\n%s", goos, got, want)
	}
}

func TestExecSettingsGoldenDarwin(t *testing.T) { checkProtectedSettingsGolden(t, "darwin") }
func TestExecSettingsGoldenLinux(t *testing.T)  { checkProtectedSettingsGolden(t, "linux") }

func TestExecSettingsNoProtectedIsByteIdentical(t *testing.T) {
	for _, goos := range []string{"darwin", "linux"} {
		for _, empty := range []bool{false, true} {
			t.Run(goos+"/empty="+strconv.FormatBool(empty), func(t *testing.T) {
				spec := goldenSettingsSpec(goos)
				if empty {
					spec.Protected = []string{}
				}
				want, err := os.ReadFile(filepath.Join("testdata", "exec_settings_"+goos+".golden.json"))
				if err != nil {
					t.Fatal(err)
				}
				got, err := renderExecSettings(spec)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(got, want) {
					t.Fatal("absent/empty protected changed legacy golden bytes")
				}
				path := filepath.Join(t.TempDir(), "settings.json")
				digest, err := writeExecSettings(path, got)
				if err != nil {
					t.Fatal(err)
				}
				if digest != "sha256:"+sha256Hex(want) {
					t.Fatal("absent/empty protected changed settings_digest")
				}
			})
		}
	}
}
