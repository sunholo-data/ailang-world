package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
)

// Row 189 M4: `profile init` (w-world-onboarding.md §3.2, AC5).

func writeExe(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatal(err)
	}
}

func writeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// goFixture is a Go repo whose PATH holds only a fake toolchain: its go
// answers `go env GOROOT GOMODCACHE GOVERSION` with a GOROOT whose bin/ holds
// that go and, with gofmt, a gofmt. GOROOT lies inside the module cache, as a
// toolchain go.mod switched to does. It returns the repo and the module cache.
func goFixture(t *testing.T, gofmt bool) (string, string) {
	t.Helper()
	base := t.TempDir()
	repo, mod := filepath.Join(base, "repo"), filepath.Join(base, "gomod")
	goroot := filepath.Join(mod, "golang.org", "toolchain@v0.0.1-go1.26.6.fixture")
	writeFile(t, filepath.Join(repo, "go.mod"), "module example.com/fixture\n\ngo 1.26\n")
	writeFile(t, filepath.Join(repo, "Makefile"), "VAR := x\n.PHONY: test\nbuild:\n\tgo build ./...\ntest: build\n\tgo test ./...\n")
	if err := os.MkdirAll(mod, 0o755); err != nil {
		t.Fatal(err)
	}
	writeExe(t, filepath.Join(goroot, "bin", "go"), "if [ \"$1\" = env ]; then printf '%s\\n%s\\n%s\\n' '"+goroot+"' '"+mod+
		"' go1.26.6; exit 0; fi\necho go version go1.26.6\n")
	if gofmt {
		writeExe(t, filepath.Join(goroot, "bin", "gofmt"), "exit 0\n")
	}
	t.Setenv("PATH", filepath.Join(goroot, "bin"))
	return repo, mod
}

// pyFixture is platform-source's layout: no marker at the root, the uv project
// in backend/, and backend/.venv/bin/python linking to uv's version-alias
// directory, which links to the real install (row 140 M4 finding F1). PATH
// holds nothing, so no uv, go or node is found.
type pyFix struct {
	repo, venvPython, realPython string
	roots                        []string
}

func pyFixture(t *testing.T) pyFix {
	t.Helper()
	base := t.TempDir()
	uv := filepath.Join(base, "uv-python")
	realDir := filepath.Join(uv, "cpython-3.12.13-test")
	writeExe(t, filepath.Join(realDir, "bin", "python3.12"), "echo Python 3.12.13\n")
	alias := filepath.Join(uv, "cpython-3.12-test")
	if err := os.Symlink("cpython-3.12.13-test", alias); err != nil {
		t.Fatal(err)
	}
	repo := filepath.Join(base, "repo")
	writeFile(t, filepath.Join(repo, "README.md"), "fixture\n")
	writeFile(t, filepath.Join(repo, "backend", "pyproject.toml"), "[project]\nname = \"fixture\"\n")
	writeFile(t, filepath.Join(repo, "backend", "uv.lock"), "version = 1\n")
	venvPython := filepath.Join(repo, "backend", ".venv", "bin", "python")
	if err := os.MkdirAll(filepath.Dir(venvPython), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(alias, "bin", "python3.12"), venvPython); err != nil {
		t.Fatal(err)
	}
	realPython, err := filepath.EvalSymlinks(venvPython)
	if err != nil {
		t.Fatal(err)
	}
	empty := filepath.Join(base, "empty-path")
	if err := os.MkdirAll(empty, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", empty)
	roots := []string{alias, filepath.Dir(filepath.Dir(realPython))}
	sort.Strings(roots)
	return pyFix{repo: repo, venvPython: venvPython, realPython: realPython, roots: roots}
}

func runProfileInitCLI(args ...string) (int, string, string) {
	var out, errw bytes.Buffer
	code := run(append([]string{"profile", "init"}, args...), &out, &errw)
	return code, out.String(), errw.String()
}

func readParsed(t *testing.T, path string) ([]byte, *broker.ExecProfile) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("the proposal was not written: %v", err)
	}
	p, err := broker.ParseExecProfile(data)
	if err != nil {
		t.Fatalf("the written proposal does not parse: %v\n%s", err, data)
	}
	return data, p
}

func commandNames(p *broker.ExecProfile) []string {
	var names []string
	for n := range p.Commands {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func dirNames(t *testing.T, dir string) []string {
	t.Helper()
	ents, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range ents {
		names = append(names, e.Name())
	}
	return names
}

func TestProfileInitGoRepo(t *testing.T) {
	repo, mod := goFixture(t, true)
	dir := filepath.Join(t.TempDir(), "profiles") // does not exist yet
	code, out, errs := runProfileInitCLI("--repo", repo, "--project", "go-fixture", "--profiles-dir", dir)
	if code != exitOK {
		t.Fatalf("exit %d, stderr %q", code, errs)
	}
	file := filepath.Join(dir, "go-fixture.proposed.json")
	data, p := readParsed(t, file)
	if got, want := commandNames(p), []string{"go-test", "go-vet", "gofmt-l"}; !reflect.DeepEqual(got, want) {
		t.Errorf("commands %v, want %v", got, want)
	}
	if p.Root != "." || p.Env["GOTOOLCHAIN"] != "go1.26.6" || p.Env["GOMODCACHE"] != mod || !reflect.DeepEqual(p.ReadRoots, []string{mod}) {
		t.Errorf("root %q, env %v, read_roots %v: want the values go env printed", p.Root, p.Env, p.ReadRoots)
	}
	if _, ok := p.Caches["GOCACHE"]; !ok || p.Probe.Argv[0] != "go" {
		t.Errorf("caches %v, probe %v: want GOCACHE and go version", p.Caches, p.Probe.Argv)
	}
	for _, name := range []string{"go-test", "go-vet", "gofmt-l"} {
		if !strings.Contains(out, "  "+name+" ") || !strings.Contains(out, "go.mod; the shape of ailang-world.json's "+name) {
			t.Errorf("stdout has no line explaining %s:\n%s", name, out)
		}
	}
	for _, want := range []string{"mv " + file + " " + filepath.Join(dir, "go-fixture.json"), "It is NOT active",
		"Makefile targets (listed only; no command is proposed from them): build, test", "serve validates and probes it at startup"} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout lacks %q:\n%s", want, out)
		}
	}
	if got := dirNames(t, dir); !reflect.DeepEqual(got, []string{"go-fixture.proposed.json"}) {
		t.Errorf("profiles dir holds %v; init must write only the proposal, never the active profile", got)
	}
	// Control: the parse that passed is not vacuous. Without read_roots the
	// same bytes fail on GOMODCACHE.
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	delete(m, "read_roots")
	mutant, _ := json.Marshal(m)
	if _, err := broker.ParseExecProfile(mutant); err == nil || !strings.Contains(err.Error(), "GOMODCACHE") {
		t.Errorf("control: the proposal without read_roots parsed (err %v)", err)
	}
}

func TestProfileInitPythonUvRepo(t *testing.T) {
	fx := pyFixture(t)
	dir := t.TempDir()
	code, out, errs := runProfileInitCLI("--repo", fx.repo, "--project", "py-fixture", "--profiles-dir", dir)
	if code != exitOK {
		t.Fatalf("exit %d, stderr %q", code, errs)
	}
	_, p := readParsed(t, filepath.Join(dir, "py-fixture.proposed.json"))
	if got, want := commandNames(p), []string{"pytest-file", "ruff-file"}; !reflect.DeepEqual(got, want) {
		t.Errorf("commands %v, want %v", got, want)
	}
	if p.Root != "backend" {
		t.Errorf("root %q, want backend (the only directory with a marker)", p.Root)
	}
	if !reflect.DeepEqual(p.ReadRoots, fx.roots) {
		t.Errorf("read_roots %v, want the alias and the real install %v", p.ReadRoots, fx.roots)
	}
	if !reflect.DeepEqual(p.Probe.Argv, []string{fx.realPython, "--version"}) {
		t.Errorf("probe %v, want the interpreter's realpath --version", p.Probe.Argv)
	}
	for _, name := range []string{"pytest-file", "ruff-file"} {
		if !strings.Contains(out, "backend/pyproject.toml + backend/uv.lock; the shape of platform-source.json's "+name) {
			t.Errorf("stdout has no line explaining %s:\n%s", name, out)
		}
	}
	// Control: the interpreter is resolved, never guessed. With no venv and
	// no uv on PATH, init refuses and writes nothing.
	if err := os.RemoveAll(filepath.Join(fx.repo, "backend", ".venv")); err != nil {
		t.Fatal(err)
	}
	dir2 := t.TempDir()
	code, _, errs = runProfileInitCLI("--repo", fx.repo, "--project", "py-fixture", "--profiles-dir", dir2)
	if code != exitUsage || !strings.Contains(errs, "no uv on PATH (run uv sync in") {
		t.Errorf("control: no interpreter: exit %d, stderr %q", code, errs)
	}
	if got := dirNames(t, dir2); len(got) != 0 {
		t.Errorf("control: wrote %v", got)
	}
}

func TestProfileInitPackageJSON(t *testing.T) {
	for _, tc := range []struct {
		name, devDeps, want string
	}{
		{"vitest", `{"vitest": "^3.0.0"}`, "vitest-file"},
		{"jest", `{"jest": "^29.0.0"}`, "jest-file"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			base := t.TempDir()
			writeExe(t, filepath.Join(base, "bin", "node"), "echo v22.0.0\n")
			t.Setenv("PATH", filepath.Join(base, "bin"))
			repo := filepath.Join(base, "repo")
			writeFile(t, filepath.Join(repo, "package.json"), `{"name": "x", "devDependencies": `+tc.devDeps+`}`)
			dir := t.TempDir()
			code, out, errs := runProfileInitCLI("--repo", repo, "--project", "node-fixture", "--profiles-dir", dir)
			if code != exitOK {
				t.Fatalf("exit %d, stderr %q", code, errs)
			}
			_, p := readParsed(t, filepath.Join(dir, "node-fixture.proposed.json"))
			if got := commandNames(p); !reflect.DeepEqual(got, []string{tc.want}) {
				t.Errorf("commands %v, want [%s]", got, tc.want)
			}
			if !strings.Contains(out, tc.want+"  ") || !strings.Contains(out, "package.json; the shape of") {
				t.Errorf("stdout does not explain %s:\n%s", tc.want, out)
			}
		})
	}
}

type measuredCommands struct {
	Source   string                 `json:"source"`
	Commands map[string]initCommand `json:"commands"`
}

func loadMeasured(t *testing.T, name string) measuredCommands {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", "profile-init", name))
	if err != nil {
		t.Fatal(err)
	}
	var m measuredCommands
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	return m
}

// shapeDiff is "" when a and b are the same command shape.
func shapeDiff(a, b initCommand) string {
	ja, _ := json.Marshal(a)
	jb, _ := json.Marshal(b)
	if bytes.Equal(ja, jb) {
		return ""
	}
	return string(ja) + " != " + string(jb)
}

// The proposed commands have the measured profiles' shapes: initShapes and
// what init actually wrote for the Go and Python fixtures both equal the
// checked-in copies of the measured command entries.
func TestProfileInitShapesMatchTheMeasuredProfiles(t *testing.T) {
	goRepo, _ := goFixture(t, true)
	goDir := t.TempDir()
	if code, _, errs := runProfileInitCLI("--repo", goRepo, "--project", "go-shapes", "--profiles-dir", goDir); code != exitOK {
		t.Fatalf("go: exit %d, stderr %q", code, errs)
	}
	fx := pyFixture(t)
	pyDir := t.TempDir()
	if code, _, errs := runProfileInitCLI("--repo", fx.repo, "--project", "py-shapes", "--profiles-dir", pyDir); code != exitOK {
		t.Fatalf("python: exit %d, stderr %q", code, errs)
	}
	for _, tc := range []struct {
		file, written string
		want          int
	}{
		{"ailang-world.commands.json", filepath.Join(goDir, "go-shapes.proposed.json"), 3},
		{"platform-source.commands.json", filepath.Join(pyDir, "py-shapes.proposed.json"), 2},
	} {
		measured := loadMeasured(t, tc.file)
		if len(measured.Commands) != tc.want {
			t.Fatalf("%s holds %d commands, want %d", tc.file, len(measured.Commands), tc.want)
		}
		raw, err := os.ReadFile(tc.written)
		if err != nil {
			t.Fatal(err)
		}
		var written struct {
			Commands map[string]initCommand `json:"commands"`
		}
		if err := json.Unmarshal(raw, &written); err != nil {
			t.Fatal(err)
		}
		for name, want := range measured.Commands {
			if d := shapeDiff(initShapes[name].cmd, want); d != "" {
				t.Errorf("%s: initShapes[%s]: %s", tc.file, name, d)
			}
			if d := shapeDiff(written.Commands[name], want); d != "" {
				t.Errorf("%s: the written %s: %s", tc.file, name, d)
			}
		}
	}
	// Control: shapeDiff sees a one-field difference in each kind of field.
	base := initShapes["pytest-file"].cmd
	for _, mutate := range []func(c *initCommand){
		func(c *initCommand) { c.MaxArgs++ },
		func(c *initCommand) { c.Argv = append([]string{}, c.Argv[:3]...) },
		func(c *initCommand) { c.Flags = map[string]initFlag{"-k": {Class: "regex"}, "-x": {Class: "bool"}} },
		func(c *initCommand) { c.Positional = "relpath" },
	} {
		m := base
		mutate(&m)
		if shapeDiff(m, base) == "" {
			t.Errorf("control: shapeDiff missed a mutation: %+v", m)
		}
	}
}

func TestProfileInitRefusals(t *testing.T) {
	goRepo, _ := goFixture(t, true)
	bare := filepath.Join(t.TempDir(), "bare")
	writeFile(t, filepath.Join(bare, "README.md"), "nothing to build\n")
	writeFile(t, filepath.Join(bare, "Makefile"), "all:\n\techo hi\n")
	writeFile(t, filepath.Join(bare, "web", "package.json"), `{"devDependencies": {"eslint": "^9"}}`)
	for _, tc := range []struct {
		name    string
		repo    string
		project string
		prepare func(t *testing.T, dir string) []string // extra args
		want    string
		code    int
	}{
		{name: "control", repo: goRepo, project: "fine", code: exitOK},
		{name: "project-grammar", repo: goRepo, project: "Bad_Name", want: "is outside the profile grammar", code: exitUsage},
		{name: "project-too-long", repo: goRepo, project: strings.Repeat("a", 65), want: "is outside the profile grammar", code: exitUsage},
		{name: "no-toolchain", repo: bare, project: "bare", code: exitUsage,
			want: "no detectable toolchain in " + bare},
		{name: "output-exists", repo: goRepo, project: "fine", code: exitUsage, want: "exists: profile init writes a proposal once",
			prepare: func(t *testing.T, dir string) []string {
				writeFile(t, filepath.Join(dir, "fine.proposed.json"), "keep\n")
				return nil
			}},
		{name: "out-is-active", repo: goRepo, project: "fine", code: exitUsage, want: "is the active profile",
			prepare: func(t *testing.T, dir string) []string { return []string{"--out", filepath.Join(dir, "fine.json")} }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			args := []string{"--repo", tc.repo, "--project", tc.project, "--profiles-dir", dir}
			if tc.prepare != nil {
				args = append(args, tc.prepare(t, dir)...)
			}
			before := dirNames(t, dir)
			code, _, errs := runProfileInitCLI(args...)
			if code != tc.code || !strings.Contains(errs, tc.want) {
				t.Fatalf("exit %d (want %d), stderr %q (want %q)", code, tc.code, errs, tc.want)
			}
			after := dirNames(t, dir)
			if tc.code == exitOK {
				if !reflect.DeepEqual(after, []string{"fine.proposed.json"}) {
					t.Errorf("control wrote %v", after)
				}
				return
			}
			if !reflect.DeepEqual(after, before) {
				t.Errorf("a refusal changed the profiles dir: %v -> %v", before, after)
			}
			if tc.name == "output-exists" {
				if b, _ := os.ReadFile(filepath.Join(dir, "fine.proposed.json")); string(b) != "keep\n" {
					t.Errorf("the existing file was overwritten: %q", b)
				}
			}
			if tc.name == "no-toolchain" && !strings.Contains(errs, "web/package.json: its devDependencies name neither vitest nor jest") {
				t.Errorf("no-toolchain does not say why web/package.json was skipped: %q", errs)
			}
		})
	}
}

// O_EXCL: a second run refuses and leaves the first proposal byte for byte,
// and the write itself refuses a file that appeared after the early check.
func TestProfileInitWritesOnceExclusively(t *testing.T) {
	repo, _ := goFixture(t, true)
	dir := t.TempDir()
	out := filepath.Join(dir, "elsewhere", "x.proposed.json")
	if code, _, errs := runProfileInitCLI("--repo", repo, "--project", "once", "--profiles-dir", dir, "--out", out); code != exitOK {
		t.Fatalf("first run: exit %d, stderr %q", code, errs)
	}
	first, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	code, _, errs := runProfileInitCLI("--repo", repo, "--project", "once", "--profiles-dir", dir, "--out", out)
	if code != exitUsage || !strings.Contains(errs, "never overwrites") {
		t.Errorf("second run: exit %d, stderr %q", code, errs)
	}
	if again, _ := os.ReadFile(out); !bytes.Equal(again, first) {
		t.Error("the second run changed the proposal")
	}
	late := filepath.Join(dir, "late.proposed.json")
	writeFile(t, late, "appeared late\n")
	if err := writeProposalOnce(late, []byte("{}\n")); err == nil || !strings.Contains(err.Error(), "never overwrites") {
		t.Errorf("writeProposalOnce over an existing file: %v", err)
	}
	if b, _ := os.ReadFile(late); string(b) != "appeared late\n" {
		t.Errorf("writeProposalOnce overwrote: %q", b)
	}
	// Control: a fresh path is written.
	fresh := filepath.Join(dir, "fresh.proposed.json")
	if err := writeProposalOnce(fresh, []byte("{}\n")); err != nil {
		t.Errorf("control: %v", err)
	}
}

// A command whose argv[0] is not in its toolchain's directories is refused
// and nothing is written, not even the directory. The parser alone would
// accept it wherever /usr/bin holds a gofmt: the GitHub Linux runner does,
// and this test passed on macOS but failed there (CI on #324).
func TestProfileInitWritesNothingWhenACommandIsNotInItsToolchain(t *testing.T) {
	repo, _ := goFixture(t, false) // GOROOT/bin has no gofmt
	dir := filepath.Join(t.TempDir(), "profiles")
	code, _, errs := runProfileInitCLI("--repo", repo, "--project", "nofmt", "--profiles-dir", dir)
	if code != exitUsage || !strings.Contains(errs, "commands.gofmt-l.argv: gofmt is not in the go toolchain's") ||
		!strings.Contains(errs, "so nothing was written") {
		t.Fatalf("exit %d, stderr %q", code, errs)
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Errorf("the profiles dir exists after a parse refusal (%v): %v", err, dirNames(t, dir))
	}
	// Control: with a gofmt the same run writes.
	goroot := filepath.Dir(strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))[0])
	writeExe(t, filepath.Join(goroot, "bin", "gofmt"), "exit 0\n")
	if code, _, errs := runProfileInitCLI("--repo", repo, "--project", "nofmt", "--profiles-dir", dir); code != exitOK {
		t.Errorf("control: exit %d, stderr %q", code, errs)
	}
}

// S7: QUICKSTART's profile init line uses only flags profile init documents,
// and it lies outside §0's attended-steps box (§0 lists only attended acts).
func TestProfileInitQuickstartLineMatchesTheCLI(t *testing.T) {
	raw, err := os.ReadFile("../../docs/QUICKSTART.md")
	if err != nil {
		t.Fatal(err)
	}
	guide := string(raw)
	_, box, _ := strings.Cut(guide, "## 0. Attended steps")
	box, _, _ = strings.Cut(box, "## 1. ")
	var lines []string
	for _, l := range quickstartCommandLines(guide) {
		if strings.HasPrefix(l, "/tmp/ailang-worldd profile init ") {
			lines = append(lines, l)
		}
	}
	if len(lines) != 1 {
		t.Fatalf("QUICKSTART holds %d profile init lines, want 1: %q", len(lines), lines)
	}
	if strings.Contains(box, "profile init") {
		t.Error("QUICKSTART §0 names profile init; it is not an attended step")
	}
	for _, tok := range strings.Fields(lines[0]) {
		if strings.HasPrefix(tok, "--") && !strings.Contains(profileHelp, "  "+tok+" ") {
			t.Errorf("QUICKSTART's profile init flag %s is not in profile init --help", tok)
		}
	}
}

func TestProfileInitHelp(t *testing.T) {
	for _, argv := range [][]string{{"profile", "init", "--help"}, {"help", "profile", "init"}, {"profile", "--help"}} {
		var out, errw bytes.Buffer
		if code := run(argv, &out, &errw); code != exitOK || !strings.Contains(out.String(), "ailang-worldd profile init --repo") {
			t.Errorf("%q: exit %d, stdout %q, stderr %q", argv, code, out.String(), errw.String())
		}
	}
}
