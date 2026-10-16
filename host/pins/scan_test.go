package pins

import (
	"go/scanner"
	"go/token"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"testing"
)

// releaseLit is an AILANG release literal. Every World release is v0.N.M, and
// a dev-build suffix ("v0.30.0-205-g…-dirty") still names its base release.
// No leading \b: a release glued to a word character ("ailang_v0.41.0") still counts
// (row 139 M1 judge r1, M9). Digits/dots on either side are excluded so "v0.41.0.1" and
// "1v0.41.0" are not misread as releases.
var releaseLit = regexp.MustCompile(`v0\.[0-9]+\.[0-9]+`)

// releaseBoundary reports whether s[i:j] is a whole release: no digit or dot touches it
// on either side. Nothing is consumed, so adjacent releases ("v0.41.0 v0.33.0") all count.
func releaseBoundary(s string, i, j int) bool {
	num := func(c byte) bool { return c == '.' || (c >= '0' && c <= '9') }
	return (i == 0 || !num(s[i-1])) && (j == len(s) || !num(s[j]))
}

// repoRoot is the module root, found from this file's own path.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

// productionFiles lists the TRACKED production files AC1 scans (git ls-files,
// so gitignored state such as .ailang/ and untracked scratch never counts):
//   - non-test Go under cmd/ and host/ (the whole production Go surface,
//     host/store TestProductionGoSurface);
//   - scripts/**/*.sh, *.py, *.json, except test_* files and testdata/;
//   - scripts/testdata/ailang_release_observed.txt (a production input of
//     verify_ail.sh, not a fixture);
//   - .github/workflows/*.yml.
func productionFiles(t *testing.T, root string) []string {
	t.Helper()
	out, err := exec.Command("git", "-C", root, "ls-files", "-z").Output()
	if err != nil {
		t.Fatalf("git ls-files: %v", err)
	}
	var files []string
	for _, f := range strings.Split(string(out), "\x00") {
		if isProductionFile(f) {
			files = append(files, f)
		}
	}
	sort.Strings(files)
	return files
}

func isProductionFile(f string) bool {
	base := path.Base(f)
	switch {
	case f == "":
		return false
	case f == "scripts/testdata/ailang_release_observed.txt":
		return true
	case strings.Contains("/"+f, "/testdata/"):
		return false
	case strings.HasSuffix(f, ".go"):
		return (strings.HasPrefix(f, "cmd/") || strings.HasPrefix(f, "host/")) && !strings.HasSuffix(f, "_test.go")
	case strings.HasPrefix(f, "scripts/"):
		if strings.HasPrefix(base, "test_") {
			return false
		}
		ext := path.Ext(f)
		return ext == ".sh" || ext == ".py" || ext == ".json"
	case strings.HasPrefix(f, ".github/workflows/"):
		return strings.HasSuffix(f, ".yml") || strings.HasSuffix(f, ".yaml")
	}
	return false
}

// releaseLiterals counts the release literals in one file's CODE. Go is
// tokenized and only string/char literals count (comments are history, not
// pins). Other files drop whole-line `#` comments; everything else counts,
// prose in a docstring included, because a sentence naming a pinned release
// goes stale when the pin moves.
func releaseLiterals(name string, src []byte) (map[string]int, error) {
	counts := map[string]int{}
	add := func(s string) {
		for _, ix := range releaseLit.FindAllStringIndex(s, -1) {
			if releaseBoundary(s, ix[0], ix[1]) {
				counts[s[ix[0]:ix[1]]]++
			}
		}
	}
	if strings.HasSuffix(name, ".go") {
		fset := token.NewFileSet()
		file := fset.AddFile(name, -1, len(src))
		var s scanner.Scanner
		var firstErr error
		s.Init(file, src, func(pos token.Position, msg string) {
			if firstErr == nil {
				firstErr = &scanError{pos.String() + ": " + msg}
			}
		}, 0)
		for {
			_, tok, lit := s.Scan()
			if tok == token.EOF {
				break
			}
			if tok == token.STRING || tok == token.CHAR {
				add(lit)
			}
		}
		return counts, firstErr
	}
	for _, line := range strings.Split(string(src), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		add(line)
	}
	return counts, nil
}

type scanError struct{ s string }

func (e *scanError) Error() string { return e.s }

// site is one file's EXPECTED release literals, by role: n occurrences of the
// release holding role (for observed: n of EACH observed release), or, for a
// role-free entry, n of the named release, which must still be in the manifest.
type site struct {
	role, release string
	n             int
	why           string
}

// allowedSites is the complete set of production files that may carry a
// release literal. Go has NONE: every Go site derives from this package. A
// file absent here must carry zero literals; a file present must carry
// EXACTLY these counts, so a new literal in an explained file reds too.
var allowedSites = map[string][]site{
	".github/workflows/ci.yml": {
		{role: RoleFrozenCompiler, n: 4, why: "the WORLD_PKG_AILANG_BIN install (world-package leg): step name, download URL, version assert and its error"},
		{role: RoleGatePrimary, n: 8, why: "the AILANG_BIN installs in go-verify and go-race-verify (row 181): per job step name, download URL, version assert and its error"},
		{role: RoleTool, n: 8, why: "the floor grader's tool-binary install: step name, URL, install paths, version assert and its error"},
	},
	"scripts/verify_world_package.sh": {
		{role: RoleFrozenCompiler, n: 2, why: "step 9 compiler version assert and its error message"},
	},
	"scripts/verify_go.sh": {
		{role: RoleGatePrimary, n: 4, why: "AILANG_BIN exact-token assert, its hint and its error message (which names the release twice)"},
	},
	"scripts/verify_ail.sh": {
		{release: "v0.30.0", n: 2, why: "release-verdict instrument control tokens (any released token works; this one is an observed release)"},
	},
	"scripts/testdata/ailang_release_observed.txt": {
		{role: RoleObserved, n: 1, why: "the observed set itself; equality with the manifest is TestObservedReleasesFileMatchesManifest"},
	},
	"scripts/world_package_ready_packet.golden.json": {
		{role: RoleFrozenCompiler, n: 1, why: "the ready packet's compilerVersion"},
	},
	"scripts/floor/floor_config.json": {
		{role: RoleTool, n: 1, why: "floor tool_bin path"},
	},
	"scripts/floor/grade.py": {
		{role: RoleTool, n: 2, why: "DEFAULT_TOOL_BIN and the docstring naming the tool binary"},
	},
	"scripts/floor/agreement.py": {
		{role: RoleTool, n: 1, why: "docstring naming the pinned tool binary"},
	},
}

// expectedCounts resolves a file's allowlist against manifest m.
func expectedCounts(t *testing.T, m Manifest, sites []site) map[string]int {
	t.Helper()
	want := map[string]int{}
	for _, s := range sites {
		switch {
		case s.release != "":
			if _, ok := m.Lookup(s.release); !ok {
				t.Errorf("allowlisted release %s (%s) is not in the manifest", s.release, s.why)
			}
			want[s.release] += s.n
		case s.role == RoleObserved:
			for _, r := range m.Holders(RoleObserved) {
				want[r.Release] += s.n
			}
		default:
			r, err := m.Holder(s.role)
			if err != nil {
				t.Fatal(err)
			}
			want[r.Release] += s.n
		}
	}
	return want
}

// checkLiterals compares a file's measured literals with the manifest and the
// allowlist, returning one finding per unexplained discrepancy.
func checkLiterals(t *testing.T, m Manifest, name string, got map[string]int) []string {
	t.Helper()
	var findings []string
	for rel := range got {
		if _, ok := m.Lookup(rel); !ok {
			findings = append(findings, name+": release literal "+rel+" is not in pins/ailang-releases.json")
		}
	}
	want := expectedCounts(t, m, allowedSites[name])
	keys := map[string]bool{}
	for k := range got {
		keys[k] = true
	}
	for k := range want {
		keys[k] = true
	}
	var sorted []string
	for k := range keys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	for _, rel := range sorted {
		if got[rel] != want[rel] {
			findings = append(findings, name+": "+rel+" appears "+itoa(got[rel])+" time(s), the manifest and allowlist explain "+itoa(want[rel]))
		}
	}
	return findings
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for ; n > 0; n /= 10 {
		b = append([]byte{byte('0' + n%10)}, b...)
	}
	return string(b)
}

// AC1 (row 139): every AILANG release literal in production files is
// explained by the manifest — a release it lists, at a site this test names,
// with the count the site's role resolves to. Moving a role in the manifest
// therefore reds every literal site still naming the old release.
func TestReleaseLiteralsExplainedByManifest(t *testing.T) {
	root := repoRoot(t)
	m := Current()
	files := productionFiles(t, root)
	if len(files) < 150 {
		t.Fatalf("scan saw %d production files, want >= 150 (is git ls-files working?)", len(files))
	}
	seen := map[string]bool{}
	goFiles, hits := 0, 0
	for _, f := range files {
		src, err := os.ReadFile(filepath.Join(root, f))
		if err != nil {
			t.Fatal(err)
		}
		got, err := releaseLiterals(f, src)
		if err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		if strings.HasSuffix(f, ".go") {
			goFiles++
		}
		seen[f] = true
		for _, n := range got {
			hits += n
		}
		for _, finding := range checkLiterals(t, m, f, got) {
			t.Error(finding)
		}
	}
	for f := range allowedSites {
		if !seen[f] {
			t.Errorf("allowlisted site %s is not a scanned production file", f)
		}
	}
	// Non-vacuity: the V7 sites the scan must reach were reached, and it
	// tokenized Go rather than silently skipping it.
	for _, f := range []string{"cmd/ailang-worldd/pins.go", "cmd/ailang-worldd/demo_proof.go", "cmd/ailang-worldd/demo_seed.go",
		"cmd/world-publish/fences.go", "host/daemon/workspace.go", "scripts/verify_world_package.sh",
		".github/workflows/ci.yml", "scripts/testdata/ailang_release_observed.txt"} {
		if !seen[f] {
			t.Errorf("scan missed V7 site %s", f)
		}
	}
	if goFiles < 100 || hits == 0 {
		t.Fatalf("scan saw %d Go files and %d literals; the instrument is not looking", goFiles, hits)
	}
	t.Logf("AC1 scan: %d production files (%d Go), %d release literals", len(files), goFiles, hits)
}

// Positive controls: the same scanner and checker, fed a literal that is not
// explained, must report it — in Go (string, not comment), in shell, and as a
// real production file mutated in memory.
func TestReleaseLiteralScanCatchesInjectedLiterals(t *testing.T) {
	m := Current()
	root := repoRoot(t)
	goSrc := []byte("package x\n// a comment naming v0.41.0 is history\nconst a = \"AILANG v0.41.0\"\nvar b = `raw v0.52.1 text`\nvar c = \"ailang_v0.41.0\"\nvar d = \"v0.30.0 v0.33.0 v0.30.0\"\nvar e = \"v0.41.0.1\"\n")
	got, err := releaseLiterals("host/x/x.go", goSrc)
	if err != nil {
		t.Fatal(err)
	}
	if got["v0.41.0"] != 2 || got["v0.52.1"] != 1 || got["v0.30.0"] != 2 || got["v0.33.0"] != 1 || len(got) != 4 {
		t.Fatalf("Go tokenizer counts %v, want v0.41.0 x2 (one glued to a word, the comment and v0.41.0.1 ignored), v0.52.1 x1, adjacent v0.30.0 x2 + v0.33.0 x1", got)
	}
	if f := checkLiterals(t, m, "host/x/x.go", got); len(f) < 2 {
		t.Fatalf("Go literals in an unlisted file gave %d findings, want 2: %v", len(f), f)
	}
	shSrc := []byte("# v0.41.0 in a comment\nbin=\"$HOME/.tools/v0.99.0/ailang\"\n")
	got, _ = releaseLiterals("scripts/x.sh", shSrc)
	f := checkLiterals(t, m, "scripts/x.sh", got)
	if len(f) != 2 || !strings.Contains(f[0], "not in pins/ailang-releases.json") {
		t.Fatalf("unknown release v0.99.0 in shell gave %v, want a not-in-manifest finding and a count finding", f)
	}
	// A real V7 site, mutated: re-introducing the pre-M1 literal reds.
	for name, inject := range map[string]string{
		"cmd/ailang-worldd/demo_proof.go": "\nconst reintroduced = \"AILANG v0.41.0\"\n",
		"scripts/verify_world_package.sh": "\necho 'AILANG v0.41.0'\n",
		".github/workflows/ci.yml":        "\n          echo v0.52.1\n",
	} {
		src, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		clean, _ := releaseLiterals(name, src)
		if f := checkLiterals(t, m, name, clean); len(f) != 0 {
			t.Fatalf("control precondition: %s is not clean: %v", name, f)
		}
		dirty, _ := releaseLiterals(name, append(src, inject...))
		if f := checkLiterals(t, m, name, dirty); len(f) == 0 {
			t.Errorf("injected literal in %s was not caught", name)
		}
	}
	// A manifest that moves the tool role reds the sites still naming v0.52.1.
	moved := m
	moved.Releases = append([]Release(nil), m.Releases...)
	for i := range moved.Releases {
		var roles []Role
		for _, r := range moved.Releases[i].Roles {
			if r.Role != RoleTool {
				roles = append(roles, r)
			}
		}
		if moved.Releases[i].Release == "v0.41.0" {
			roles = append(roles, Role{Role: RoleTool, Since: "2026-10-16", Evidence: "control"})
		}
		moved.Releases[i].Roles = roles
	}
	src, _ := os.ReadFile(filepath.Join(root, ".github/workflows/ci.yml"))
	got, _ = releaseLiterals(".github/workflows/ci.yml", src)
	if f := checkLiterals(t, moved, ".github/workflows/ci.yml", got); len(f) == 0 {
		t.Error("moving the tool role in the manifest did not red ci.yml's tool install")
	}
}

func TestIsProductionFile(t *testing.T) {
	for f, want := range map[string]bool{
		"host/daemon/workspace.go":                     true,
		"host/daemon/workspace_test.go":                false,
		"host/pkgproj/testdata/gen.go":                 false,
		"design_docs/verification/x/main.go":           false,
		"scripts/verify_ail.sh":                        true,
		"scripts/test_queue_census.sh":                 false,
		"scripts/floor/grade.py":                       true,
		"scripts/floor/test_grade.py":                  false,
		"scripts/testdata/upstream_release_tags.txt":   false,
		"scripts/testdata/ailang_release_observed.txt": true,
		".github/workflows/ci.yml":                     true,
		"tools/launchd/driver.sh":                      false,
		".ailang/state/x.go":                           false,
	} {
		if got := isProductionFile(f); got != want {
			t.Errorf("isProductionFile(%q) = %v, want %v", f, got, want)
		}
	}
}
