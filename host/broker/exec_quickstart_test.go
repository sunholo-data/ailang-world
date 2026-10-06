package broker

// Row 140 M3 (w-workspace-exec-toolchain-effect §6 AC3.5, coding-standards
// S7): QUICKSTART §10's profile script is run verbatim (python3, with its
// environment pointed at a fixture) and the three profiles it writes — Go,
// TypeScript, Python — load: every key known, every argv[0] resolved, the
// read_roots legal against the host rules. Its tools/call payload matches
// the Go profile's grammar.

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"
)

func quickstartSection10(t *testing.T) string {
	t.Helper()
	guide, err := os.ReadFile(filepath.Join("..", "..", "docs", "QUICKSTART.md"))
	if err != nil {
		t.Fatal(err)
	}
	_, section, ok := strings.Cut(string(guide), "### 10. Build and test non-AILANG projects")
	if !ok {
		t.Fatal("QUICKSTART §10 is absent")
	}
	return section
}

func TestQuickstartSection10ProfilesLoad(t *testing.T) {
	py, err := exec.LookPath("python3")
	if err != nil {
		t.Skipf("SKIP %s: python3 absent: the §10 profile script is python3", t.Name())
	}
	section := quickstartSection10(t)
	_, script, ok := strings.Cut(section, "python3 - <<'EOF'\n")
	script, _, ok2 := strings.Cut(script, "\nEOF\n")
	if !ok || !ok2 {
		t.Fatal("QUICKSTART §10 has no python3 profile script")
	}
	base, h := hostPathsFixture(t)
	bin := fakeToolDir(t, "go", "node")
	execDir := filepath.Join(base, "exec")
	gomod := filepath.Join(h.OperatorHome, "go", "pkg", "mod")
	pybin := filepath.Join(h.OperatorHome, ".local", "share", "uv", "python", "cpython-3.12.13-test", "bin", "python3.12")
	mustMkdir(t, filepath.Join(execDir, "profiles"), filepath.Dir(pybin))
	if err := os.WriteFile(pybin, []byte("#!/bin/sh\necho Python 3.12.13\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	// uv's layout (M4 finding F1): the venv links to a version-ALIAS dir, which links to the
	// real one. The script must make both read roots.
	aliasDir := filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(pybin))), "cpython-3.12-test")
	if err := os.Symlink(filepath.Dir(filepath.Dir(pybin)), aliasDir); err != nil {
		t.Fatal(err)
	}
	pylink := filepath.Join(aliasDir, "bin", "python3.12")
	cmd := exec.Command(py, "-")
	cmd.Stdin = strings.NewReader(script)
	cmd.Env = []string{"PATH=/usr/bin:/bin", "HOME=" + h.OperatorHome, "EXEC=" + execDir, "GODIR=" + bin, "NODEDIR=" + bin,
		"GOTC=go1.26.6", "GOMOD=" + gomod, "PYBIN=" + pybin, "PYLINK=" + pylink}
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("the §10 profile script failed: %v\n%s", err, out)
	}
	profiles := map[string]*ExecProfile{}
	for _, name := range []string{"ailang-compiler", "twilightgame", "sunholo-cli"} {
		data, err := os.ReadFile(filepath.Join(execDir, "profiles", name+".json"))
		if err != nil {
			t.Fatalf("the §10 script did not write %s: %v", name, err)
		}
		p, err := ParseExecProfile(data)
		if err != nil {
			t.Fatalf("§10 profile %s does not load: %v", name, err)
		}
		if err := CheckExecReadRoots(p.ReadRoots, h); err != nil {
			t.Fatalf("§10 profile %s: %v", name, err)
		}
		if p.Project != name {
			t.Fatalf("§10 profile %s names project %q", name, p.Project)
		}
		profiles[name] = p
	}
	// Both uv dirs are read roots: the alias the venv links to, and its realpath (F1).
	for _, root := range []string{aliasDir, filepath.Dir(filepath.Dir(pybin))} {
		if !slices.Contains(profiles["sunholo-cli"].ReadRoots, root) {
			t.Errorf("§10 sunholo-cli read_roots %q lack %s", profiles["sunholo-cli"].ReadRoots, root)
		}
	}
	// The probes: go version, node --version, the venv's interpreter.
	for name, want := range map[string][]string{"ailang-compiler": {"go", "version"}, "twilightgame": {"node", "--version"},
		"sunholo-cli": {pybin, "--version"}} {
		if got := profiles[name].Probe.Argv; !reflect.DeepEqual(got, want) {
			t.Errorf("§10 profile %s probe = %q, want %q", name, got, want)
		}
	}
	// pytest's -k is emitted as two items (form sep), after the profile's
	// sh -c '… "$@"' prefix, which receives the agent's args verbatim.
	py3 := profiles["sunholo-cli"].Commands["test-file"]
	norm, err := MatchExecArgs(py3.ExecCommand, []string{"-k", "not slow", "tests/test_cli.py"})
	if err != nil {
		t.Fatal(err)
	}
	if got := EmitExecArgs(py3.ExecCommand, norm); !reflect.DeepEqual(got, []string{"-k", "not slow", "tests/test_cli.py"}) {
		t.Fatalf("pytest args emitted as %q", got)
	}
	if !strings.HasSuffix(py3.Argv[2], `"$@"`) || py3.Argv[3] != "pytest" {
		t.Fatalf("the Python prefix does not hand the args to pytest as \"$@\": %q", py3.Argv)
	}
	// The §10 tools/call payload is a call the Go profile admits.
	m := regexp.MustCompile(`-d '([^']+)'`).FindStringSubmatch(section)
	if m == nil {
		t.Fatal("QUICKSTART §10 has no tools/call payload")
	}
	var call struct {
		Params struct {
			Name      string `json:"name"`
			Arguments struct {
				Command string   `json:"command"`
				Args    []string `json:"args"`
			} `json:"arguments"`
		} `json:"params"`
	}
	if err := json.Unmarshal([]byte(m[1]), &call); err != nil || call.Params.Name != "workspace-exec" {
		t.Fatalf("§10 payload %s: %v", m[1], err)
	}
	c, ok := profiles["ailang-compiler"].Commands[call.Params.Arguments.Command]
	if !ok {
		t.Fatalf("§10 calls command %q, which the Go profile lacks", call.Params.Arguments.Command)
	}
	if _, err := MatchExecArgs(c.ExecCommand, call.Params.Arguments.Args); err != nil {
		t.Fatalf("§10's call is refused by the Go profile's grammar: %v", err)
	}
}
