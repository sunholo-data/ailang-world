package broker

// Row 140 M3 (w-workspace-exec-toolchain-effect §6 AC3.5, coding-standards
// S7): QUICKSTART §10's profile script is run verbatim (python3, with its
// environment pointed at a fixture) and the three profiles it writes — Go,
// TypeScript, Python — load: every key known, every argv[0] resolved, the
// read_roots legal against the host rules. Its tools/call payload matches
// the Go profile's grammar.

import (
	"bytes"
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
	section, _, _ = strings.Cut(section, "### 11.")
	return section
}

// quickstartProfileFiles executes the checked-in generator in an isolated fixture.
func quickstartProfileFiles(t *testing.T) (map[string][]byte, ExecHostPaths, string, string) {
	t.Helper()
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
	cmd := exec.CommandContext(boundedTestContext(t), py, "-")
	cmd.Stdin = strings.NewReader(script)
	cmd.Env = []string{"PATH=/usr/bin:/bin", "HOME=" + h.OperatorHome, "EXEC=" + execDir, "GODIR=" + bin, "NODEDIR=" + bin,
		"GOTC=go1.26.6", "GOMOD=" + gomod, "PYBIN=" + pybin, "PYLINK=" + pylink}
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("the §10 profile script failed: %v\n%s", err, out)
	}
	files := map[string][]byte{}
	for _, name := range []string{"ailang-compiler", "twilightgame", "sunholo-cli"} {
		data, err := os.ReadFile(filepath.Join(execDir, "profiles", name+".json"))
		if err != nil {
			t.Fatalf("the §10 script did not write %s: %v", name, err)
		}
		files[name] = data
	}
	return files, h, aliasDir, pybin
}

func TestQuickstartSection10ProfilesLoad(t *testing.T) {
	files, h, aliasDir, pybin := quickstartProfileFiles(t)
	section := quickstartSection10(t)
	profiles := map[string]*ExecProfile{}
	for name, data := range files {
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
	for _, row := range []struct {
		name, project string
		want          []string
	}{
		{"protected-go", "ailang-compiler", []string{"*.golden", "*_test.go", "tests/**"}},
		{"protected-ts", "twilightgame", []string{"*.test.ts", "*.test.tsx", "tests/**", "vitest.config.ts"}},
		{"protected-python", "sunholo-cli", []string{"cli/tests/**", "conftest.py", "pyproject.toml"}},
	} {
		t.Run(row.name, func(t *testing.T) {
			p := profiles[row.project]
			if !reflect.DeepEqual(p.Protected, row.want) {
				t.Fatalf("%s protected = %q, want %q", row.project, p.Protected, row.want)
			}
			if p.ProtectedDigest == "" || p.ProtectedDigest != protectedDigest(row.want) {
				t.Fatalf("%s protected digest = %q, want %q", row.project, p.ProtectedDigest, protectedDigest(row.want))
			}
		})
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

// quickstartProtectedBlock scopes parsing so unrelated tables/examples cannot satisfy a guard.
func quickstartProtectedBlock(t *testing.T, kind string) string {
	t.Helper()
	_, block, ok := strings.Cut(quickstartSection10(t), "<!-- exec-protected-"+kind+":begin -->")
	block, _, end := strings.Cut(block, "<!-- exec-protected-"+kind+":end -->")
	if !ok || !end {
		t.Fatalf("protected %s block is absent", kind)
	}
	return block
}

func quickstartProtectedRows(t *testing.T, kind string) [][]string {
	t.Helper()
	var rows [][]string
	for _, line := range strings.Split(quickstartProtectedBlock(t, kind), "\n") {
		if !strings.HasPrefix(line, "| ") {
			continue
		}
		cells := strings.Split(strings.Trim(line, "|"), "|")
		for i := range cells {
			cells[i] = strings.TrimSpace(cells[i])
		}
		rows = append(rows, cells)
	}
	if len(rows) < 2 {
		t.Fatalf("protected %s table has no data rows", kind)
	}
	return rows[1:] // header; the separator is not a data row
}

func TestQuickstartProtectedContract(t *testing.T) {
	t.Run("patterns", func(t *testing.T) {
		counts := map[string]int{}
		for _, row := range quickstartProtectedRows(t, "patterns") {
			if len(row) != 3 || (row[1] != "accepted" && row[1] != "refused") {
				t.Fatalf("malformed pattern row: %q", row)
			}
			// Remove only code delimiters: whitespace inside them must reach the parser.
			pattern := strings.Trim(row[0], "`")
			_, err := ParseProtectedPatterns([]string{pattern})
			if (err == nil) != (row[1] == "accepted") {
				t.Errorf("pattern %q documented %s, parser error = %v", pattern, row[1], err)
			}
			counts[row[1]]++
		}
		if counts["accepted"] < 1 || counts["refused"] < 1 {
			t.Fatalf("pattern table needs accepted and refused rows: %v", counts)
		}
	})
	t.Run("enforcement", func(t *testing.T) {
		root := t.TempDir()
		seen := map[string]bool{}
		for _, row := range quickstartProtectedRows(t, "os") {
			if len(row) != 4 || (row[0] != "linux" && row[0] != "darwin") || seen[row[0]] {
				t.Fatalf("malformed OS row: %q", row)
			}
			seen[row[0]] = true
			if got := protectedEnforced(row[0]); got != row[1] {
				t.Errorf("%s enforced documented %q, executable = %q", row[0], row[1], got)
			}
			want := []string{filepath.Join(root, "tests")}
			baseClaim := "detected only"
			if row[0] == "darwin" {
				want = append(want, root+"/**/*.test.ts")
				baseClaim = "refused"
			}
			if row[2] != "refused" || row[3] != baseClaim {
				t.Errorf("%s enforcement cells = %q, want refused / %s", row[0], row[2:], baseClaim)
			}
			if got := renderSrtProtected([]string{"tests/**", "*.test.ts"}, root, row[0]); !reflect.DeepEqual(got, want) {
				t.Errorf("%s rendered %q, want %q", row[0], got, want)
			}
		}
		if !seen["linux"] || !seen["darwin"] {
			t.Fatal("OS table must cover linux and darwin")
		}
	})
	t.Run("record", func(t *testing.T) {
		block := quickstartProtectedBlock(t, "record")
		_, block, ok := strings.Cut(block, "```json\n")
		block, _, end := strings.Cut(block, "\n```")
		var envelope map[string]json.RawMessage
		if !ok || !end || json.Unmarshal([]byte(block), &envelope) != nil || len(envelope) != 1 {
			t.Fatal("protected record example must be one JSON result block")
		}
		var record protectedRecord
		decoder := json.NewDecoder(bytes.NewReader(envelope["protected"]))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&record); err != nil {
			t.Fatalf("protected record example: %v", err)
		}
		var keys map[string]json.RawMessage
		if err := json.Unmarshal(envelope["protected"], &keys); err != nil {
			t.Fatal(err)
		}
		typ := reflect.TypeOf(record)
		if len(keys) != typ.NumField() {
			t.Fatalf("record has %d keys, struct has %d fields", len(keys), typ.NumField())
		}
		for i := 0; i < typ.NumField(); i++ {
			key := strings.Split(typ.Field(i).Tag.Get("json"), ",")[0]
			if _, ok := keys[key]; !ok {
				t.Errorf("record example lacks %q", key)
			}
		}
		if record.IntactAfter == nil || *record.IntactAfter || record.Enforced != protectedEnforced("darwin") {
			t.Fatalf("record must show a reported change: %+v", record)
		}
	})
	t.Run("no-default", func(t *testing.T) {
		files, _, _, _ := quickstartProfileFiles(t)
		var profile map[string]json.RawMessage
		if err := json.Unmarshal(files["ailang-compiler"], &profile); err != nil {
			t.Fatal(err)
		}
		delete(profile, "protected")
		data, err := json.Marshal(profile)
		if err != nil {
			t.Fatal(err)
		}
		p, err := ParseExecProfile(data)
		if err != nil || p == nil {
			t.Fatalf("profile without protected: %v", err)
		}
		if len(p.Protected) != 0 || p.ProtectedDigest != "" {
			t.Fatalf("absent protected must have no default: %q / %q", p.Protected, p.ProtectedDigest)
		}
		if !strings.Contains(quickstartSection10(t), "There is no default: absent `protected` or `[]` means no protected set, no `protected` key\nin the exec result, and `why` prints `protected: none declared`.") {
			t.Fatal("no-default sentence must state absent/empty behavior and protected: none declared")
		}
	})
	t.Run("refusal-message", func(t *testing.T) {
		_, err := ParseProtectedPatterns([]string{"**/x_test.go"})
		if err == nil || !strings.Contains(err.Error(), "use the base-name form `x_test.go`") || strings.Contains(err.Error(), "exactly one directory") {
			t.Fatalf("double-star refusal must name x_test.go without the old one-directory claim: %v", err)
		}
	})
}
