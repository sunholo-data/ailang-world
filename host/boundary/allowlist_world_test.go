// Package boundary holds executable tests for dependencies that cross host
// package boundaries. It intentionally contains no production code.
package boundary

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

type goGroup struct {
	name, pattern, dir, mutantFile, mutantImport string
	// extraForbidden holds import prefixes rejected for THIS group only, on top
	// of forbiddenImportPrefixes. It exists because the loopback exception below
	// is true of exactly one group, and a single shared list silently granted it
	// to all three: measured at the SM.B2a boundary, bare "net/http" blank-imported
	// into host/store/store.go left this gate green (rc=0) while the
	// "net/http/httputil" arm correctly red-ed. Baseline net/http presence in each
	// group's dependency closure is host/store 0, host/replay 0,
	// cmd/ailang-worldd 1 — so store and replay were exempt from a transport they
	// never used, in the gate whose whole purpose is confining network to
	// host/broker.
	extraForbidden []string
}

var protectedGoGroups = []goGroup{
	{"host/store", "./host/store/...", "host/store", "host/store/store.go", "net/http/httputil", []string{"net/http"}},
	{"host/replay", "./host/replay/...", "host/replay", "host/replay/replay.go", "net/http/httputil", []string{"net/http"}},
	// ailang-worldd already uses net/http for its loopback-only daemon client.
	// Keep that inherited transport visible as a narrow exception while rejecting
	// an additional HTTP surface (and all direct registry/cloud imports). This is
	// the ONLY group the exception is true of, hence the empty extraForbidden.
	{"cmd/ailang-worldd", "./cmd/ailang-worldd/...", "cmd/ailang-worldd", "cmd/ailang-worldd/main.go", "net/http/httputil", nil},
}

var forbiddenImportPrefixes = []string{
	"cloud.google.com/",
	"github.com/Azure/",
	"github.com/aws/",
	"github.com/sunholo-data/ailang-world/host/registry",
	"net/http/httptest",
	"net/http/httputil",
}

var forbiddenAILMarkers = []string{".ailang/cache", "package-cache", "package_cache", "http://", "https://"}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate repository root")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

// overlaySpec is the on-disk shape `go list -overlay` consumes.
type overlaySpec struct {
	Replace map[string]string `json:"Replace"`
}

// overlay declares a file substitution to BOTH halves of this gate: the JSON
// document at jsonPath is handed to `go list -overlay` (the dependency-closure
// half), and replace is consulted by the read helper below (the import-scan
// half). The zero value means "no substitution", i.e. exactly the behaviour
// this file had before overlays existed.
//
// The two halves are named SEPARATELY on purpose. `go list -overlay` returns
// rc=0, the BASE closure and NO stderr when a Replace key matches no real file
// (measured), so an overlay can reach the read helper and never reach the
// toolchain. If a single string served both halves, disarming one would disarm
// the other and AC2's anti-vacuity mutation (M2 form (a)) could not red at AC2
// at all -- the arm would red earlier, at AC3, and the toolchain half would go
// untested. Keeping them separable is what makes AC2 falsifiable.
type overlay struct {
	jsonPath string            // path passed to `go list -overlay`; "" = no flag
	replace  map[string]string // cleaned absolute real path -> substitute path
}

// srcFor returns the bytes the checkers must read for path, or nil meaning
// "read it from disk" -- which is byte-for-byte the pre-overlay behaviour for
// every path the overlay does not replace.
func (o overlay) srcFor(path string) ([]byte, error) {
	sub, ok := o.replace[filepath.Clean(path)]
	if !ok {
		return nil, nil
	}
	return os.ReadFile(sub)
}

// parseSrc returns the value to hand parser.ParseFile as its src argument: the
// overlay's bytes when the path is replaced, and an UNTYPED nil otherwise.
//
// Returning srcFor's []byte directly would be a silent disaster, and it was
// OBSERVED as one while building BG.A: go/parser's readSource tests `src != nil`
// on an INTERFACE, so a typed nil []byte is a non-nil interface and is handed
// back as an EMPTY source. Every unreplaced file then parses as
// "expected 'package', found 'EOF'" -- and a checker that cannot read the tree
// is a checker that finds no forbidden imports. Funnelling the conversion
// through one function keeps the nil discipline in one place.
func (o overlay) parseSrc(path string) (any, error) {
	src, err := o.srcFor(path)
	if err != nil || src == nil {
		return nil, err
	}
	return src, nil
}

// readFile is srcFor for consumers that need the bytes themselves rather than
// the nil-means-disk sentinel go/parser accepts (the AILANG arm).
func (o overlay) readFile(path string) ([]byte, error) {
	src, err := o.srcFor(path)
	if err != nil || src != nil {
		return src, err
	}
	return os.ReadFile(path)
}

// countDep counts exact matches of want in a `go list -deps` closure.
func countDep(deps []string, want string) int {
	n := 0
	for _, d := range deps {
		if d == want {
			n++
		}
	}
	return n
}

// goListDeps deliberately mirrors host/broker's bounded go-list helper. Go
// dependency closures and world/*.ail sources are different enumerations: Go
// can enumerate the former, but cannot see AILANG modules at all.
//
// overlayJSON, when non-empty, is passed as -overlay=<path>. cmd.Dir stays
// root: the mutant is DECLARED to the toolchain rather than relocated on disk,
// so the package patterns and the repository layout are unchanged.
func goListDeps(root, overlayJSON string, patterns ...string) ([]string, error) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		return nil, fmt.Errorf("cannot locate the `go` toolchain on PATH: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	args := []string{"list", "-deps"}
	if overlayJSON != "" {
		args = append(args, "-overlay="+overlayJSON)
	}
	cmd := exec.CommandContext(ctx, goBin, append(args, patterns...)...)
	cmd.Dir = root
	var stdout, stderr strings.Builder
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s list -deps %s: %w\nstderr: %s", goBin, strings.Join(patterns, " "), err, stderr.String())
	}
	var deps []string
	s := bufio.NewScanner(strings.NewReader(stdout.String()))
	for s.Scan() {
		if line := strings.TrimSpace(s.Text()); line != "" {
			deps = append(deps, line)
		}
	}
	return deps, s.Err()
}

func enumerateAIL(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(filepath.Join(root, "world"), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".ailang" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) == ".ail" {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			files = append(files, filepath.ToSlash(rel))
		}
		return nil
	})
	sort.Strings(files)
	return files, err
}

func forbiddenImport(path string, extra ...string) bool {
	for _, prefix := range append(append([]string{}, forbiddenImportPrefixes...), extra...) {
		if path == strings.TrimSuffix(prefix, "/") || strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

// checkGoGroup returns the dependency closure it ACTUALLY CONSUMED alongside
// its verdict. AC2 asserts on that returned slice, not on a second `go list`
// invocation: two calls can disagree (different pattern, different cwd, a stale
// path variable) and an assertion on the second one would be green while the
// call that gates was running on the base closure. Threading the closure out
// costs zero extra `go list` invocations -- this function already computed and
// discarded it.
func checkGoGroup(root string, group goGroup, ov overlay) ([]string, error) {
	deps, err := goListDeps(root, ov.jsonPath, group.pattern)
	if err != nil {
		return nil, err
	}
	if len(deps) == 0 {
		return deps, fmt.Errorf("%s dependency enumeration is empty: guard would pass vacuously", group.name)
	}
	return deps, filepath.WalkDir(filepath.Join(root, group.dir), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// go list -deps (without -test) describes production dependencies, so the
		// attribution scan must use the same scope.
		if d.IsDir() || filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		// src is nil for every path the overlay does not replace, which is
		// exactly today's call; for a replaced path it is the mutant's bytes,
		// parsed while positions are still attributed to the REAL path, so the
		// RED message names the production file rather than a temp file.
		src, err := ov.parseSrc(path)
		if err != nil {
			return err
		}
		parsed, err := parser.ParseFile(token.NewFileSet(), path, src, parser.ImportsOnly)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		for _, spec := range parsed.Imports {
			imp, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				return fmt.Errorf("parse import in %s: %w", path, err)
			}
			if forbiddenImport(imp, group.extraForbidden...) {
				rel, _ := filepath.Rel(root, path)
				return fmt.Errorf("%s: forbidden registry/HTTP/cloud dependency %q", filepath.ToSlash(rel), imp)
			}
		}
		return nil
	})
}

// checkAILGroup takes the same read indirection as checkGoGroup. The AILANG arm
// never used `go list`, so ONLY its read moves: a replaced world/*.ail resolves
// to the mutant bytes, every other file to os.ReadFile exactly as before.
func checkAILGroup(root string, ov overlay) error {
	files, err := enumerateAIL(root)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("world AILANG enumeration is empty: guard would pass vacuously")
	}
	for _, rel := range files {
		data, err := ov.readFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			return err
		}
		lower := strings.ToLower(string(data))
		for _, marker := range forbiddenAILMarkers {
			if strings.Contains(lower, marker) {
				return fmt.Errorf("%s: forbidden registry/package-cache marker %q", rel, marker)
			}
		}
	}
	return nil
}

func digest(data []byte) string { return fmt.Sprintf("%x", sha256.Sum256(data)) }

// armBarrierEnv names the AC4 kill-harness barrier. UNSET -- the normal case,
// including CI -- the barrier is a NO-OP costing one os.Getenv per arm.
//
// Set, its value is "<arm>=<marker path>". When the running arm's name matches
// <arm>, the arm writes a JSON ready marker to <marker path> AFTER the mutant
// and the overlay JSON have been written and validated but BEFORE the checker
// consumes them, then blocks until the process is killed or armBarrierTimeout
// elapses. That window is exactly the interval in which the OLD harness held a
// forbidden import in a production source file, so it is the interval an
// external SIGKILL harness must be able to hit deterministically instead of by
// racing.
//
// The marker path MUST resolve OUTSIDE repoRoot and is rejected otherwise:
// `git status --porcelain` reports untracked files, so an in-repo marker would
// make the harness red on its own instrument and send a reviewer hunting a
// phantom residue.
//
// A timeout is a test FAILURE, never a pass. AC4 is fail-closed: a missing
// marker, a timeout, an absent artifact or an early exit must FAIL the
// criterion, because "the kill never landed while the arm was armed" and "the
// arm left no residue" are otherwise the same green.
const armBarrierEnv = "WORLD_BOUNDARY_ARM_BARRIER"

// armBarrierTimeout bounds the barrier so a mis-driven harness cannot hang CI.
const armBarrierTimeout = 120 * time.Second

// armBarrierMarker is the JSON an armed arm publishes. The harness verifies
// Mutant and OverlayJSON exist and that the overlay maps Target -> Mutant
// before it sends the kill, so a kill can never be recorded against an arm that
// was not actually armed.
type armBarrierMarker struct {
	Arm         string `json:"arm"`
	Target      string `json:"target"`
	Mutant      string `json:"mutant"`
	OverlayJSON string `json:"overlay_json"`
	PID         int    `json:"pid"`
}

// insideRepo reports whether path lies at or beneath root, comparing
// symlink-evaluated absolute paths. On darwin /tmp and /var are symlinks
// (/var -> private/var), so a naive string-prefix test on an un-evaluated
// t.TempDir() path mis-classifies. path need not exist yet; its parent
// directory must.
func insideRepo(root, path string) (bool, error) {
	evalRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return false, err
	}
	dir, base := filepath.Split(filepath.Clean(path))
	evalDir, err := filepath.EvalSymlinks(filepath.Clean(dir))
	if err != nil {
		return false, err
	}
	rel, err := filepath.Rel(evalRoot, filepath.Join(evalDir, base))
	if err != nil {
		return false, err
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)), nil
}

func armBarrier(t *testing.T, root, arm string, m armBarrierMarker) {
	t.Helper()
	spec := os.Getenv(armBarrierEnv)
	if spec == "" {
		return // no-op: the only cost on a normal run is this one Getenv
	}
	name, marker, ok := strings.Cut(spec, "=")
	if !ok || name == "" || marker == "" {
		t.Fatalf("%s=%q is malformed: want \"<arm>=<marker path outside the repository>\"", armBarrierEnv, spec)
	}
	if name != arm {
		return
	}
	absMarker, err := filepath.Abs(marker)
	if err != nil {
		t.Fatal(err)
	}
	inside, err := insideRepo(root, absMarker)
	if err != nil {
		t.Fatalf("%s marker %s cannot be resolved against repoRoot %s: %v", armBarrierEnv, absMarker, root, err)
	}
	if inside {
		t.Fatalf("%s marker %s resolves inside repoRoot %s: `git status --porcelain` reports untracked files, so an in-repo marker would red the AC4 harness on its own artifact", armBarrierEnv, absMarker, root)
	}
	for _, artifact := range []string{m.Mutant, m.OverlayJSON} {
		if _, err := os.Stat(artifact); err != nil {
			t.Fatalf("%s cannot arm %s: artifact %s is absent: %v", armBarrierEnv, arm, artifact, err)
		}
	}
	payload, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absMarker, payload, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Logf("BARRIER arm=%s marker=%s pid=%d target=%s mutant=%s overlay_json=%s timeout=%s", arm, absMarker, m.PID, m.Target, m.Mutant, m.OverlayJSON, armBarrierTimeout)
	deadline := time.Now().Add(armBarrierTimeout)
	for time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("%s barrier for arm %s timed out after %s without being killed: AC4 is fail-closed, so a timeout FAILS the criterion rather than passing it", armBarrierEnv, arm, armBarrierTimeout)
}

// mutateViaOverlay proves the guard has teeth WITHOUT writing a byte inside the
// repository. The mutant and the overlay JSON are written to t.TempDir(); the
// substitution is DECLARED to the toolchain (`go list -overlay`) and to the
// import scan (the overlay read helper). Nothing under repoRoot is written, so
// there is nothing to restore -- and the restore logic is therefore DELETED
// rather than made safer. That is the whole repair: the old harness wrote a
// forbidden import into a production source file and undid it with a `defer`,
// and `defer` does not run on SIGKILL, so an interrupted run left the tree
// poisoned permanently with a residue invisible to `go build` but fatal to this
// gate.
//
// It returns the mutant's sha256 so callers can record the arm.
func mutateViaOverlay(t *testing.T, root, rel, arm string, mutate func([]byte) []byte, check func(ov overlay) error) string {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	absTarget, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	baseline, err := os.ReadFile(absTarget)
	if err != nil {
		t.Fatal(err)
	}
	baseSHA := digest(baseline)
	mutant := mutate(baseline)
	mutantSHA := digest(mutant)
	// Retained from the old harness: a mutation that never applied and a
	// mutation that failed to red are the same exit code.
	if mutantSHA == baseSHA {
		t.Fatalf("mutation did not apply to %s: sha256 remained %s", rel, baseSHA)
	}

	dir := t.TempDir()
	absMutant := filepath.Join(dir, "mutant__"+strings.ReplaceAll(rel, string(filepath.Separator), "__"))
	if err := os.WriteFile(absMutant, mutant, 0o600); err != nil {
		t.Fatal(err)
	}
	// Absolute paths on BOTH sides. Relative keys also apply, but absolute is
	// cmd.Dir-independent, which is what makes the gate insensitive to where
	// the checkout lives.
	blob, err := json.Marshal(overlaySpec{Replace: map[string]string{absTarget: absMutant}})
	if err != nil {
		t.Fatal(err)
	}
	absOverlayJSON := filepath.Join(dir, "overlay.json")
	if err := os.WriteFile(absOverlayJSON, blob, 0o600); err != nil {
		t.Fatal(err)
	}
	// Validate the artifacts before anything consumes them: `go list -overlay`
	// SILENTLY IGNORES a Replace key that matches no real file (rc=0, base
	// closure, no stderr), so a half-written overlay must not reach the checker
	// looking like a working one.
	written, err := os.ReadFile(absMutant)
	if err != nil {
		t.Fatal(err)
	}
	if digest(written) != mutantSHA {
		t.Fatalf("overlay mutant for %s did not round-trip: wrote %s, read back %s", rel, mutantSHA, digest(written))
	}
	if _, err := os.Stat(absTarget); err != nil {
		t.Fatalf("overlay Replace key %s does not name a real file: go list would silently ignore it: %v", absTarget, err)
	}
	ov := overlay{jsonPath: absOverlayJSON, replace: map[string]string{filepath.Clean(absTarget): absMutant}}
	t.Logf("OVERLAY path=%s arm=%s baseline_sha256=%s mutant_sha256=%s mutant=%s overlay_json=%s", rel, arm, baseSHA, mutantSHA, absMutant, absOverlayJSON)

	// AC4: the barrier sits after the artifacts are written and validated and
	// before the checker consumes them. No-op unless armBarrierEnv is set.
	armBarrier(t, root, arm, armBarrierMarker{
		Arm:         arm,
		Target:      absTarget,
		Mutant:      absMutant,
		OverlayJSON: absOverlayJSON,
		PID:         os.Getpid(),
	})

	err = check(ov)
	if err == nil {
		t.Fatalf("mutation in %s passed boundary guard", rel)
	}
	if !strings.Contains(err.Error(), rel) {
		t.Fatalf("guard failure did not name exact path %s: %v", rel, err)
	}
	t.Logf("MUTATION path=%s baseline_sha256=%s mutant_sha256=%s guard=%q", rel, baseSHA, mutantSHA, err)
	return mutantSHA
}

func TestWorldBoundaryDependencyAllowlist(t *testing.T) {
	root := repoRoot(t)
	// Enumerate and print every protected group before any mutation. Go groups
	// are dependency closures; world is the distinct AILANG source enumeration.
	// baselineDeps holds the NON-overlay closure per group; it is the negative
	// half of AC2 and it is computed here, once, before any mutation exists.
	baselineDeps := make(map[string][]string, len(protectedGoGroups))
	for _, group := range protectedGoGroups {
		deps, err := goListDeps(root, "", group.pattern)
		if err != nil {
			t.Fatal(err)
		}
		if len(deps) == 0 {
			t.Fatalf("%s dependency enumeration is empty", group.name)
		}
		baseCount := countDep(deps, group.mutantImport)
		t.Logf("ENUMERATION group=%s exact_count=%d baseline_%s=%d dependencies=%q", group.name, len(deps), group.mutantImport, baseCount, deps)
		// AC2's negative half, asserted before any arm runs: if the mutant
		// import were ALREADY in the baseline closure, "the overlay closure
		// contains it" would be satisfied by a dead overlay too.
		if baseCount != 0 {
			t.Fatalf("baseline (non-overlay) closure for %s already contains %q %d time(s): the overlay assertion in the arms below would not discriminate", group.name, group.mutantImport, baseCount)
		}
		baselineDeps[group.name] = deps
		if _, err := checkGoGroup(root, group, overlay{}); err != nil {
			t.Fatal(err)
		}
	}
	ailFiles, err := enumerateAIL(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(ailFiles) == 0 {
		t.Fatal("world AILANG enumeration is empty")
	}
	t.Logf("ENUMERATION group=world exact_count=%d files=%q", len(ailFiles), ailFiles)
	if err := checkAILGroup(root, overlay{}); err != nil {
		t.Fatal(err)
	}

	for _, group := range protectedGoGroups {
		group := group
		t.Run("mutation_"+strings.ReplaceAll(group.name, "/", "_"), func(t *testing.T) {
			// overlayDeps is the closure checkGoGroup ITSELF CONSUMED. AC2 is
			// asserted on this slice and on nothing else: a second, separate
			// goListDeps call would prove the overlay works for THAT call and
			// say nothing about the call that gates.
			var overlayDeps []string
			mutateViaOverlay(t, root, group.mutantFile, group.name, func(baseline []byte) []byte {
				anchor := []byte("import (\n")
				if bytes.Count(baseline, anchor) != 1 {
					t.Fatalf("mutation anchor count for %s is %d, want 1", group.mutantFile, bytes.Count(baseline, anchor))
				}
				insert := []byte(fmt.Sprintf("import (\n\t_ %q // boundary mutation: compiling HTTP import\n", group.mutantImport))
				return bytes.Replace(baseline, anchor, insert, 1)
			}, func(ov overlay) error {
				var err error
				overlayDeps, err = checkGoGroup(root, group, ov)
				return err
			})

			// AC2 -- the overlay demonstrably reached `go list`. Both halves,
			// both counts printed. This is not decoration: `go list -overlay`
			// returns rc=0, the BASE closure and NO stderr when a Replace key
			// matches no real file, so without this assertion the toolchain half
			// of the gate could be dead while the import scan (which reads
			// through our own helper) kept producing a convincing RED.
			base := baselineDeps[group.name]
			baseCount := countDep(base, group.mutantImport)
			ovCount := countDep(overlayDeps, group.mutantImport)
			t.Logf("AC2 group=%s baseline_closure=%d baseline_%s=%d overlay_closure=%d overlay_%s=%d", group.name, len(base), group.mutantImport, baseCount, len(overlayDeps), group.mutantImport, ovCount)
			if len(base) == 0 || len(overlayDeps) == 0 {
				t.Fatalf("AC2 closure enumeration for %s is empty (baseline=%d, overlay=%d): the assertion below would pass vacuously", group.name, len(base), len(overlayDeps))
			}
			if ovCount == 0 {
				t.Fatalf("overlay closure does not contain %q: `go list` never saw the mutant (overlay closure=%d packages, baseline closure=%d) -- the toolchain half of the gate is dead", group.mutantImport, len(overlayDeps), len(base))
			}
			if baseCount != 0 {
				t.Fatalf("baseline closure for %s already contains %q %d time(s): the overlay assertion is not a discriminator", group.name, group.mutantImport, baseCount)
			}
		})
	}
	t.Run("mutation_world", func(t *testing.T) {
		mutateViaOverlay(t, root, "world/types.ail", "world", func(baseline []byte) []byte {
			return append(append([]byte(nil), baseline...), []byte("\n-- boundary mutation: package-cache .ailang/cache\n")...)
		}, func(ov overlay) error { return checkAILGroup(root, ov) })
	})

	// Green control: the broker is intentionally the network boundary. Its
	// non-empty closure must remain permitted by this inverse protected-group guard.
	broker, err := goListDeps(root, "", "./host/broker/...")
	if err != nil {
		t.Fatal(err)
	}
	if len(broker) == 0 {
		t.Fatal("host/broker dependency enumeration is empty")
	}
	t.Logf("GREEN_CONTROL group=host/broker exact_count=%d result=PASS", len(broker))
}

func TestWorldBoundaryNullCases(t *testing.T) {
	if err := func(files []string) error {
		if len(files) == 0 {
			return fmt.Errorf("empty")
		}
		return nil
	}(nil); err == nil {
		t.Fatal("empty AILANG enumeration passed vacuously")
	}
	if forbiddenImport("fmt") {
		t.Fatal("stdlib control unexpectedly forbidden")
	}
	if !forbiddenImport("net/http/httputil") {
		t.Fatal("HTTP mutation control unexpectedly permitted")
	}
}

// TestBareNetHTTPExemptionIsPerGroup pins the asymmetry that the single shared
// forbiddenImportPrefixes list used to erase. The loopback-daemon exception is
// true of cmd/ailang-worldd and of nothing else, so bare "net/http" must be
// rejected for host/store and host/replay while remaining permitted for
// cmd/ailang-worldd. Collapsing extraForbidden back into one global list makes
// this test fail rather than silently re-granting the exemption.
func TestBareNetHTTPExemptionIsPerGroup(t *testing.T) {
	byName := make(map[string]goGroup, len(protectedGoGroups))
	for _, g := range protectedGoGroups {
		byName[g.name] = g
	}
	if len(byName) != 3 {
		t.Fatalf("protected group enumeration is %d, want 3: the guard below would not cover what it claims", len(byName))
	}

	for _, want := range []struct {
		group     string
		forbidden bool
	}{
		{"host/store", true},
		{"host/replay", true},
		{"cmd/ailang-worldd", false}, // documented loopback-IPC exception
	} {
		g, ok := byName[want.group]
		if !ok {
			t.Fatalf("protected group %q absent: enumeration no longer covers it", want.group)
		}
		if got := forbiddenImport("net/http", g.extraForbidden...); got != want.forbidden {
			t.Errorf("forbiddenImport(\"net/http\") for %s = %v, want %v", want.group, got, want.forbidden)
		}
		// net/http/httputil stays forbidden everywhere, per the shared list.
		if !forbiddenImport("net/http/httputil", g.extraForbidden...) {
			t.Errorf("net/http/httputil unexpectedly permitted for %s", want.group)
		}
	}
}
