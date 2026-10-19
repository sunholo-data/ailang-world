package main

// Row 189 M4 (design_docs/planned/w-world-onboarding.md §3.2, D-WORLD-85:
// D-189-5 = A): `ailang-worldd profile init` proposes an exec profile for a
// repo from the repo's own toolchain files.
//
//  1. Refuse first: a --project outside the profile grammar, an output file
//     that exists, and an --out that is the active <profiles-dir>/<project>.json.
//  2. Detect: go.mod → go-test, go-vet, gofmt-l; pyproject.toml or uv.lock →
//     pytest-file, ruff-file; package.json → vitest-file or jest-file, when
//     its devDependencies name vitest or jest. A Makefile's targets are only
//     listed. The repo root wins; a repo with no marker at its root is read
//     one directory down (platform-source keeps its Python in backend/), and
//     that directory becomes the profile's root.
//  3. Resolve every path on this machine, never by guessing: go env GOROOT
//     GOMODCACHE GOVERSION in the root, the venv python's link chain (or uv
//     python find), node on PATH. A toolchain that does not resolve is
//     skipped with its reason; if none is left, init refuses.
//  4. Each command has the shape of the measured profile's command of the
//     same name (initShapes). The proposal must pass broker.ParseExecProfile
//     before anything is written, and then it is written once (O_EXCL).
//
// It never activates a profile: serve skips *.proposed.json, and init prints
// the mv that activates one. serve's startup probe stays the real check.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/toolprobe"
)

const profileHelp = `usage: ailang-worldd profile init --repo <dir> --project <name> [--out <file>] [--profiles-dir <dir>]

Propose an exec profile (` + broker.ExecProfileVersion + `) for a repo from its
own toolchain files, and write it to <profiles-dir>/<project>.proposed.json. It
never activates the profile: serve skips *.proposed.json. Review the file, then
run the mv it prints. serve validates and probes the profile at startup.

Detection (at the repo root first; with no marker there, one directory down,
and that directory becomes the profile's root):
  go.mod                     go-test, go-vet, gofmt-l
  pyproject.toml or uv.lock  pytest-file, ruff-file
  package.json               vitest-file or jest-file, whichever its
                             devDependencies name (neither: skipped)
  Makefile                   its targets are listed; no command is proposed
Each command has the shape (argv, flags, positional, max_args) of the measured
profile's command of the same name. Every path comes from this machine, never a
guess: go env GOROOT GOMODCACHE GOVERSION run in the root, the venv's python
(<root>/.venv/bin/python, else uv python find) and its link chain, and node on
PATH. A toolchain that does not resolve is skipped, and the reason is printed.

The proposal must pass the profile parser before it is written. If it does not,
nothing is written and the parser's error is printed. Refused: a --project
outside ` + broker.ExecProjectGrammar + `, an output file that already exists
(the file is written once, O_EXCL, and never overwritten), an --out naming the
active <profiles-dir>/<project>.json, and a repo with no detectable toolchain.

  --repo <dir>          the repository to read (required)
  --project <name>      the profile's project (required)
  --out <file>          default <profiles-dir>/<project>.proposed.json
  --profiles-dir <dir>  default ~/.ailang-world-exec/profiles
`

// initFlag is one flags entry: a class, or {"class", "form"} when it has a form.
type initFlag struct {
	Class string
	Form  string
}

func (f initFlag) MarshalJSON() ([]byte, error) {
	if f.Form == "" {
		return json.Marshal(f.Class)
	}
	return json.Marshal(struct {
		Class string `json:"class"`
		Form  string `json:"form"`
	}{f.Class, f.Form})
}

func (f *initFlag) UnmarshalJSON(b []byte) error {
	*f = initFlag{}
	if json.Unmarshal(b, &f.Class) == nil {
		return nil
	}
	var obj struct {
		Class string `json:"class"`
		Form  string `json:"form"`
	}
	if err := json.Unmarshal(b, &obj); err != nil {
		return err
	}
	f.Class, f.Form = obj.Class, obj.Form
	return nil
}

// initCommand is one commands entry of a proposal.
type initCommand struct {
	Argv       []string            `json:"argv"`
	Flags      map[string]initFlag `json:"flags,omitempty"`
	Positional string              `json:"positional,omitempty"`
	Suffixes   []string            `json:"suffixes,omitempty"`
	MaxArgs    int                 `json:"max_args,omitempty"`
}

type initShape struct {
	cmd     initCommand
	mirrors string // the measured command this one copies
}

// initShapes are the measured profiles' commands, by name.
// TestProfileInitShapesMatchTheMeasuredProfiles binds them to the checked-in
// copies under testdata/profile-init/.
var initShapes = map[string]initShape{
	"go-test": {initCommand{Argv: []string{"go", "test", "-count=1"},
		Flags:      map[string]initFlag{"-run": {Class: broker.ExecClassRegex}, "-v": {Class: broker.ExecClassBool}, "-short": {Class: broker.ExecClassBool}},
		Positional: broker.ExecPositionalPkgPattern, MaxArgs: 8}, "ailang-world.json's go-test"},
	"go-vet": {initCommand{Argv: []string{"go", "vet"}, Positional: broker.ExecPositionalPkgPattern, MaxArgs: 4},
		"ailang-world.json's go-vet"},
	"gofmt-l": {initCommand{Argv: []string{"gofmt", "-l"}, Positional: broker.ExecPositionalRelpath, MaxArgs: 8},
		"ailang-world.json's gofmt-l"},
	"pytest-file": {initCommand{Argv: []string{"sh", "-c", `exec .venv/bin/python -m pytest -q -p no:cacheprovider "$@"`, "pytest"},
		Flags:      map[string]initFlag{"-k": {Class: broker.ExecClassRegex, Form: broker.ExecFormSep}, "-x": {Class: broker.ExecClassBool}},
		Positional: broker.ExecPositionalTestfile, Suffixes: []string{".py"}, MaxArgs: 6}, "platform-source.json's pytest-file"},
	"ruff-file": {initCommand{Argv: []string{"sh", "-c", `exec .venv/bin/ruff check "$@"`, "ruff"},
		Positional: broker.ExecPositionalRelpath, MaxArgs: 8}, "platform-source.json's ruff-file"},
	"vitest-file": {initCommand{Argv: []string{"node", "node_modules/vitest/vitest.mjs", "run"},
		Positional: broker.ExecPositionalTestfile, Suffixes: []string{".test.ts", ".test.tsx"}, MaxArgs: 4},
		"QUICKSTART section 10's twilightgame test-file (row 140 M4)"},
	"jest-file": {initCommand{Argv: []string{"node", "node_modules/jest/bin/jest.js"},
		Positional: broker.ExecPositionalTestfile, Suffixes: []string{".test.js", ".test.jsx", ".test.ts", ".test.tsx"}, MaxArgs: 4},
		"vitest-file's shape with jest's entry point (no jest profile is measured)"},
}

// initProposal is the written profile, in the measured profiles' key order.
type initProposal struct {
	Profile   string                 `json:"profile"`
	Project   string                 `json:"project"`
	Root      string                 `json:"root"`
	Path      []string               `json:"path"`
	Env       map[string]string      `json:"env,omitempty"`
	ReadRoots []string               `json:"read_roots,omitempty"`
	Caches    map[string]initCache   `json:"caches,omitempty"`
	TimeoutMS int                    `json:"timeout_ms"`
	Probe     []string               `json:"probe"`
	Commands  map[string]initCommand `json:"commands"`
}

type initCache struct {
	Seed string `json:"seed,omitempty"`
}

// Toolchain kinds, in priority order: below the repo root the first kind
// found picks the directory, and the first resolved kind gives the probe.
const (
	initKindGo     = "go"
	initKindPython = "python"
	initKindNode   = "node"
)

var initKindRank = map[string]int{initKindGo: 0, initKindPython: 1, initKindNode: 2}

// initDetected is one toolchain found in one directory of the repo.
type initDetected struct {
	kind     string
	dir      string   // "." or a child of the repo
	markers  []string // the repo-relative files that triggered it
	commands []string
}

// initResolved is what a detected toolchain contributes to the profile.
type initResolved struct {
	path      []string
	env       map[string]string
	readRoots []string
	caches    map[string]initCache
	probe     []string
	note      string // where the paths came from
}

// profileInitBudget bounds profile init's toolchain probes as a whole: at
// most go env and uv python find, each also under toolprobe's own 30 s.
const profileInitBudget = 75 * time.Second

// runProfile dispatches the `profile` verb family.
func runProfile(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, profileHelp)
		return exitUsage
	}
	switch args[0] {
	case "init":
		return runProfileInit(args[1:], stdout, stderr)
	case "--help", "-help", "-h", "help":
		fmt.Fprint(stdout, profileHelp)
		return exitOK
	default:
		fmt.Fprintf(stderr, "ailang-worldd profile: unknown subcommand %q\n\n", args[0])
		fmt.Fprint(stderr, profileHelp)
		return exitUsage
	}
}

func runProfileInit(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("ailang-worldd profile init", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	repo := fs.String("repo", "", "the repository to read")
	project := fs.String("project", "", "the profile's project")
	out := fs.String("out", "", "the proposal file (default <profiles-dir>/<project>.proposed.json)")
	profilesDir := fs.String("profiles-dir", defaultExecProfilesDir(), "the exec profiles directory")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprint(stdout, profileHelp)
			return exitOK
		}
		fmt.Fprint(stderr, profileHelp)
		return exitUsage
	}
	fail := func(format string, a ...any) int {
		fmt.Fprintf(stderr, "ailang-worldd profile init: refused: "+format+"\n", a...)
		return exitUsage
	}
	if fs.NArg() > 0 {
		return fail("unexpected argument %q", fs.Arg(0))
	}
	if *repo == "" || *project == "" {
		return fail("--repo and --project are required")
	}
	if !broker.ExecProjectOk(*project) {
		return fail("--project %q is outside the profile grammar %s", *project, broker.ExecProjectGrammar)
	}
	repoAbs, err := filepath.Abs(*repo)
	if err == nil {
		var info os.FileInfo
		if info, err = os.Stat(repoAbs); err == nil && !info.IsDir() {
			err = errors.New("not a directory")
		}
	}
	if err != nil {
		return fail("--repo %s: %v", *repo, err)
	}
	dirAbs, err := filepath.Abs(*profilesDir)
	if err != nil {
		return fail("--profiles-dir %s: %v", *profilesDir, err)
	}
	active := filepath.Join(dirAbs, *project+broker.ExecProfileSuffix)
	outPath := filepath.Join(dirAbs, *project+broker.ExecProposedSuffix)
	if *out != "" {
		if outPath, err = filepath.Abs(*out); err != nil {
			return fail("--out %s: %v", *out, err)
		}
	}
	if sameFileName(outPath, active) {
		return fail("--out %s is the active profile %s: profile init writes only a proposal, and the operator activates it (D-189-5)", outPath, active)
	}
	if _, err := os.Lstat(outPath); err == nil {
		return fail("%s exists: profile init writes a proposal once and never overwrites it", outPath)
	}

	ctx, cancel := context.WithTimeout(context.Background(), profileInitBudget)
	defer cancel()
	p, lines, err := proposeProfile(ctx, repoAbs, *project)
	if err != nil {
		return fail("%v", err)
	}
	data, err := marshalProposal(p)
	if err != nil {
		return fail("%v", err)
	}
	parsed, err := broker.ParseExecProfile(data)
	if err != nil {
		return fail("the proposal does not pass the profile parser, so nothing was written: %v", err)
	}
	if err := writeProposalOnce(outPath, data); err != nil {
		return fail("%v", err)
	}

	fmt.Fprintf(stdout, "proposed exec profile %q for %s (root %s):\n", p.Project, repoAbs, p.Root)
	for _, l := range lines {
		fmt.Fprintf(stdout, "  %s\n", l)
	}
	fmt.Fprintf(stdout, "wrote %s (%s). It is NOT active: serve skips *%s.\n", outPath, parsed.Digest, broker.ExecProposedSuffix)
	fmt.Fprintln(stdout, "review it, then activate it with:")
	fmt.Fprintf(stdout, "  mv %s %s\n", shellQuote(outPath), shellQuote(active))
	if _, err := os.Lstat(active); err == nil {
		fmt.Fprintf(stdout, "note: %s exists already; that mv replaces it.\n", active)
	}
	fmt.Fprintf(stdout, "serve validates and probes it at startup: started with --exec-sandbox, it loads every <project>.json in "+
		"--exec-profiles-dir (here %s), or pass the file alone with --exec-profile. A profile that fails to parse or probe "+
		"is left out and named in /v1/health.\n", dirAbs)
	return exitOK
}

// sameFileName reports a and b naming one file: equal, or equal once their
// directories are resolved (a symlinked profiles dir).
func sameFileName(a, b string) bool {
	if a == b {
		return true
	}
	return filepath.Base(a) == filepath.Base(b) && realpathOrSelf(filepath.Dir(a)) == realpathOrSelf(filepath.Dir(b))
}

func realpathOrSelf(p string) string {
	if real, err := filepath.EvalSymlinks(p); err == nil {
		return real
	}
	return p
}

func marshalProposal(p initProposal) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")
	if err := enc.Encode(p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// writeProposalOnce creates path (and its directory) and writes data, or
// refuses when path exists: O_EXCL, so a file that appears after the early
// check is never overwritten either.
func writeProposalOnce(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("%s exists: profile init writes a proposal once and never overwrites it", path)
		}
		return err
	}
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		_ = os.Remove(path)
		return fmt.Errorf("write %s: %w", path, err)
	}
	return f.Close()
}

// proposeProfile detects and resolves the repo's toolchains and assembles
// the proposal, with one line per command saying why it is there, then the
// notes (where each path came from, what was skipped, the Makefile targets).
func proposeProfile(ctx context.Context, repoAbs, project string) (initProposal, []string, error) {
	root, found, notes := detectInitRoot(repoAbs)
	rootAbs := filepath.Join(repoAbs, root)
	for _, dir := range uniqueStrings([]string{".", root}) {
		if targets := makefileTargets(filepath.Join(repoAbs, dir, "Makefile")); len(targets) > 0 {
			notes = append(notes, fmt.Sprintf("%s targets (listed only; no command is proposed from them): %s",
				filepath.Join(dir, "Makefile"), strings.Join(targets, ", ")))
		}
	}
	if len(found) == 0 {
		return initProposal{}, nil, fmt.Errorf("no detectable toolchain in %s: looked for go.mod, pyproject.toml, uv.lock, "+
			"and a package.json whose devDependencies name vitest or jest, at the root and one directory down%s",
			repoAbs, joinNotes(notes))
	}
	home, _ := os.UserHomeDir()
	p := initProposal{Profile: broker.ExecProfileVersion, Project: project, Root: root, Env: map[string]string{},
		Caches: map[string]initCache{}, TimeoutMS: broker.ExecMaxTimeoutMS, Commands: map[string]initCommand{}}
	var lines, skipped []string
	for _, d := range found {
		var r initResolved
		var err error
		switch d.kind {
		case initKindGo:
			r, err = resolveInitGo(ctx, rootAbs)
		case initKindPython:
			r, err = resolveInitPython(ctx, repoAbs, rootAbs)
		case initKindNode:
			r, err = resolveInitNode(home)
		}
		if err != nil {
			skipped = append(skipped, fmt.Sprintf("skipped %s (%s): %v", strings.Join(d.commands, ", "), strings.Join(d.markers, " + "), err))
			continue
		}
		for _, name := range d.commands {
			s := initShapes[name]
			p.Commands[name] = s.cmd
			lines = append(lines, fmt.Sprintf("%-12s %s; the shape of %s", name, strings.Join(d.markers, " + "), s.mirrors))
		}
		p.Path = append(p.Path, r.path...)
		p.ReadRoots = append(p.ReadRoots, r.readRoots...)
		for k, v := range r.env {
			p.Env[k] = v
		}
		for k, v := range r.caches {
			p.Caches[k] = v
		}
		if p.Probe == nil {
			p.Probe = r.probe
			notes = append(notes, fmt.Sprintf("probe %s: the %s toolchain's", strings.Join(r.probe, " "), d.kind))
		}
		notes = append(notes, r.note)
	}
	notes = append(notes, skipped...)
	if len(p.Commands) == 0 {
		return initProposal{}, nil, fmt.Errorf("no detected toolchain resolves on this machine%s", joinNotes(notes))
	}
	for _, d := range []string{"/usr/bin", "/bin"} {
		if info, err := os.Stat(d); err == nil && info.IsDir() {
			p.Path = append(p.Path, d)
		}
	}
	p.Path = uniqueStrings(p.Path)
	p.ReadRoots = uniqueStrings(p.ReadRoots)
	sort.Strings(p.ReadRoots)
	return p, append(lines, notes...), nil
}

func joinNotes(notes []string) string {
	if len(notes) == 0 {
		return ""
	}
	return " (" + strings.Join(notes, "; ") + ")"
}

// detectInitRoot picks the profile's root: the repo root when it holds any
// marker, else the child directory whose first toolchain ranks highest
// (ties: name order). Toolchains elsewhere are noted, never proposed: a
// profile has one root.
func detectInitRoot(repoAbs string) (string, []initDetected, []string) {
	found, notes := detectInitIn(repoAbs, ".")
	if len(found) > 0 {
		return ".", found, notes
	}
	ents, _ := os.ReadDir(repoAbs)
	root := "."
	var best []initDetected
	for _, e := range ents {
		name := e.Name()
		if !e.IsDir() || strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" {
			continue
		}
		f, n := detectInitIn(repoAbs, name)
		notes = append(notes, n...)
		if len(f) == 0 {
			continue
		}
		if best != nil && initKindRank[f[0].kind] >= initKindRank[best[0].kind] {
			notes = append(notes, otherRootNote(f))
			continue
		}
		if best != nil {
			notes = append(notes, otherRootNote(best))
		}
		root, best = name, f
	}
	return root, best, notes
}

func otherRootNote(f []initDetected) string {
	var markers []string
	for _, d := range f {
		markers = append(markers, d.markers...)
	}
	return fmt.Sprintf("also found %s (not proposed: a profile has one root)", strings.Join(markers, " + "))
}

// detectInitIn reads one directory's markers, in kind order.
func detectInitIn(repoAbs, dir string) ([]initDetected, []string) {
	var found []initDetected
	var notes []string
	has := func(name string) bool {
		info, err := os.Stat(filepath.Join(repoAbs, dir, name))
		return err == nil && info.Mode().IsRegular()
	}
	rel := func(name string) string { return filepath.Join(dir, name) }
	if has("go.mod") {
		found = append(found, initDetected{kind: initKindGo, dir: dir, markers: []string{rel("go.mod")},
			commands: []string{"go-test", "go-vet", "gofmt-l"}})
	}
	var py []string
	for _, m := range []string{"pyproject.toml", "uv.lock"} {
		if has(m) {
			py = append(py, rel(m))
		}
	}
	if len(py) > 0 {
		found = append(found, initDetected{kind: initKindPython, dir: dir, markers: py, commands: []string{"pytest-file", "ruff-file"}})
	}
	if has("package.json") {
		var pkg struct {
			DevDependencies map[string]json.RawMessage `json:"devDependencies"`
		}
		raw, err := os.ReadFile(filepath.Join(repoAbs, dir, "package.json"))
		if err == nil {
			err = json.Unmarshal(raw, &pkg)
		}
		var cmds []string
		for _, t := range []string{"vitest", "jest"} {
			if _, ok := pkg.DevDependencies[t]; ok {
				cmds = append(cmds, t+"-file")
			}
		}
		switch {
		case err != nil:
			notes = append(notes, fmt.Sprintf("%s does not parse (%v); no command proposed from it", rel("package.json"), err))
		case len(cmds) == 0:
			notes = append(notes, fmt.Sprintf("%s: its devDependencies name neither vitest nor jest; no command proposed from it", rel("package.json")))
		default:
			found = append(found, initDetected{kind: initKindNode, dir: dir, markers: []string{rel("package.json")}, commands: cmds})
		}
	}
	return found, notes
}

// resolveInitGo asks the repo's own go for the toolchain: the go that runs
// in the root selects the toolchain go.mod names, so GOROOT/bin is that
// toolchain's go and gofmt.
func resolveInitGo(ctx context.Context, rootAbs string) (initResolved, error) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		return initResolved{}, fmt.Errorf("no go on PATH: %v", err)
	}
	out, err := toolprobe.Output(ctx, goBin, rootAbs, "env", "GOROOT", "GOMODCACHE", "GOVERSION")
	if err != nil {
		return initResolved{}, fmt.Errorf("%s env GOROOT GOMODCACHE GOVERSION in %s: %v", goBin, rootAbs, err)
	}
	vals := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(vals) != 3 {
		return initResolved{}, fmt.Errorf("%s env printed %q, not three lines", goBin, out)
	}
	goroot, modcache, version := strings.TrimSpace(vals[0]), strings.TrimSpace(vals[1]), strings.TrimSpace(vals[2])
	if !filepath.IsAbs(goroot) || !filepath.IsAbs(modcache) || !strings.HasPrefix(version, "go") {
		return initResolved{}, fmt.Errorf("%s env printed GOROOT %q, GOMODCACHE %q, GOVERSION %q", goBin, goroot, modcache, version)
	}
	goroot, modcache = filepath.Clean(goroot), filepath.Clean(modcache)
	r := initResolved{path: []string{filepath.Join(goroot, "bin")},
		env: map[string]string{"GOFLAGS": "-mod=readonly", "GOPROXY": "off", "GOTOOLCHAIN": version,
			"GOTELEMETRY": "off", "GOMODCACHE": modcache},
		readRoots: []string{modcache}, caches: map[string]initCache{"GOCACHE": {}}, probe: []string{"go", "version"},
		note: fmt.Sprintf("go: go env in the root gave GOROOT %s, GOMODCACHE %s, GOVERSION %s; GOCACHE starts empty "+
			"(give it a seed to warm the first build)", goroot, modcache, version)}
	if !pathInsideOrEqual(goroot, modcache) {
		r.readRoots = append(r.readRoots, goroot)
	}
	return r, nil
}

// resolveInitPython finds the interpreter the commands' .venv/bin/python
// runs and makes every directory on its link chain a read root: the venv
// links to uv's version-alias directory, which links to the real one, and
// both must be readable (row 140 M4 finding F1).
func resolveInitPython(ctx context.Context, repoAbs, rootAbs string) (initResolved, error) {
	start, source := filepath.Join(rootAbs, ".venv", "bin", "python"), ""
	if _, err := os.Lstat(start); err == nil {
		source = filepath.Join(initRel(repoAbs, rootAbs), ".venv", "bin", "python")
	} else if uv, err := exec.LookPath("uv"); err == nil {
		out, err := toolprobe.Output(ctx, uv, rootAbs, "python", "find")
		if err != nil {
			return initResolved{}, fmt.Errorf("no %s, and %s python find failed: %v (run uv sync in %s)", start, uv, err, rootAbs)
		}
		start, source = strings.TrimSpace(string(out)), "uv python find"
		if !filepath.IsAbs(start) {
			return initResolved{}, fmt.Errorf("uv python find printed %q, not an absolute path", start)
		}
	} else {
		return initResolved{}, fmt.Errorf("no %s and no uv on PATH (run uv sync in %s)", start, rootAbs)
	}
	roots, real, err := interpreterRoots(start, repoAbs)
	if err != nil {
		return initResolved{}, err
	}
	note := fmt.Sprintf("python: %s resolves to %s", source, real)
	if len(roots) > 0 {
		note += "; read_roots are the install directories on its link chain: " + strings.Join(roots, ", ")
	}
	if _, err := os.Lstat(filepath.Join(rootAbs, ".venv", "bin", "python")); err != nil {
		note += fmt.Sprintf("; the commands run .venv/bin/python, so run uv sync in %s before an episode", initRel(repoAbs, rootAbs))
	}
	return initResolved{readRoots: roots, probe: []string{real, "--version"}, note: note}, nil
}

// interpreterRoots walks start's symlink chain. Each hop's install
// directory (the parent of its bin/) outside the repo is a read root, and so
// is the realpath's. It returns the roots, sorted, and the realpath.
func interpreterRoots(start, repoAbs string) ([]string, string, error) {
	repoReal := realpathOrSelf(repoAbs)
	var roots []string
	add := func(p string) {
		dir := filepath.Dir(p)
		if filepath.Base(dir) == "bin" {
			dir = filepath.Dir(dir)
		}
		if dir == "/" || pathInsideOrEqual(dir, repoAbs) || pathInsideOrEqual(dir, repoReal) {
			return
		}
		roots = append(roots, dir)
	}
	p := filepath.Clean(start)
	for hops := 0; ; hops++ {
		if hops > 40 {
			return nil, "", fmt.Errorf("%s: too many symlinks", start)
		}
		add(p)
		info, err := os.Lstat(p)
		if err != nil {
			return nil, "", fmt.Errorf("the interpreter %s does not resolve: %v", start, err)
		}
		if info.Mode()&os.ModeSymlink == 0 {
			break
		}
		target, err := os.Readlink(p)
		if err != nil {
			return nil, "", err
		}
		if !filepath.IsAbs(target) {
			target = filepath.Join(filepath.Dir(p), target)
		}
		p = filepath.Clean(target)
	}
	real, err := filepath.EvalSymlinks(start)
	if err != nil {
		return nil, "", fmt.Errorf("the interpreter %s does not resolve: %v", start, err)
	}
	add(real)
	roots = uniqueStrings(roots)
	sort.Strings(roots)
	return roots, real, nil
}

// resolveInitNode puts node's directory on the path; a node under the
// operator's HOME (nvm) also needs its install directory as a read root.
func resolveInitNode(home string) (initResolved, error) {
	node, err := exec.LookPath("node")
	if err != nil {
		return initResolved{}, fmt.Errorf("no node on PATH: %v", err)
	}
	abs, err := filepath.Abs(node)
	if err != nil {
		return initResolved{}, err
	}
	r := initResolved{path: []string{filepath.Dir(abs)}, probe: []string{"node", "--version"},
		note: fmt.Sprintf("node: %s on PATH", abs)}
	if real, err := filepath.EvalSymlinks(abs); err == nil && home != "" && pathInsideOrEqual(real, realpathOrSelf(home)) {
		dir := filepath.Dir(filepath.Dir(real))
		r.readRoots = []string{dir}
		r.note += fmt.Sprintf(" resolves under HOME to %s, so %s is a read root", real, dir)
	}
	return r, nil
}

var makeTargetLine = regexp.MustCompile(`^([A-Za-z0-9][A-Za-z0-9_./-]*)\s*:([^=]|$)`)

// makefileTargets lists a Makefile's explicit targets in file order.
func makefileTargets(path string) []string {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var out []string
	for _, line := range strings.Split(string(raw), "\n") {
		if m := makeTargetLine.FindStringSubmatch(line); m != nil {
			out = append(out, m[1])
		}
	}
	return uniqueStrings(out)
}

func initRel(base, p string) string {
	if r, err := filepath.Rel(base, p); err == nil {
		return r
	}
	return p
}

func pathInsideOrEqual(p, dir string) bool {
	return p == dir || strings.HasPrefix(p, strings.TrimSuffix(dir, "/")+"/")
}

// uniqueStrings keeps the first of each string, in order.
func uniqueStrings(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}
