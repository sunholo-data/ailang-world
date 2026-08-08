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
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
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

// granularityTrials is N in the BG.C AC1b calibration. It is the number of
// back-to-back stat -> write-new -> write-original -> stat pairs used to prove
// the nanosecond ModTime comparator can actually fire on a given filesystem
// before the ModTime half of the runtime backstop is believed. The threshold
// is deliberately NOT lowered below 20/20 on any filesystem: a probe that
// tolerates misses is a detector that can report a false negative.
const granularityTrials = 20

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

// rawWrite is the SINK. Tests swap this, never confinedWrite, so the CHECK is
// always exercised before any byte can be written.
var rawWrite = os.WriteFile

// confinedWrite is the single enforcement point every legacy os.WriteFile call
// in this file routes through. It resolves root and dst the same way insideRepo
// does (symlink-evaluated absolute paths; dst need not exist yet, its parent
// must), and if dst lies at or beneath root it returns an error SYNCHRONOUSLY,
// before rawWrite is invoked at all -- there is no write-then-check and no
// window. Otherwise it delegates to the swappable rawWrite sink.
func confinedWrite(root, dst string, data []byte, perm os.FileMode) error {
	inside, err := insideRepo(root, dst)
	if err != nil {
		return err
	}
	if inside {
		evalRoot, err := filepath.EvalSymlinks(root)
		if err != nil {
			return err
		}
		dir, base := filepath.Split(filepath.Clean(dst))
		evalDir, err := filepath.EvalSymlinks(filepath.Clean(dir))
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(evalRoot, filepath.Join(evalDir, base))
		if err != nil {
			return err
		}
		return fmt.Errorf("confined writer: destination inside repoRoot: %s", filepath.ToSlash(rel))
	}
	return rawWrite(dst, data, perm)
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
	for _, artifact := range []string{m.Mutant, m.OverlayJSON} {
		if _, err := os.Stat(artifact); err != nil {
			t.Fatalf("%s cannot arm %s: artifact %s is absent: %v", armBarrierEnv, arm, artifact, err)
		}
	}
	payload, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	// confinedWrite IS the AC4 enforcement point for the marker path: a marker
	// that resolves inside repoRoot is rejected synchronously, before a byte is
	// written. `git status --porcelain` reports untracked files, so an in-repo
	// marker would red the AC4 harness on its own artifact.
	if err := confinedWrite(root, absMarker, payload, 0o600); err != nil {
		t.Fatalf("%s marker %s rejected by confined writer (an in-repo marker would surface as untracked residue in `git status --porcelain`): %v", armBarrierEnv, absMarker, err)
	}
	t.Logf("BARRIER arm=%s marker=%s pid=%d target=%s mutant=%s overlay_json=%s timeout=%s", arm, absMarker, m.PID, m.Target, m.Mutant, m.OverlayJSON, armBarrierTimeout)
	deadline := time.Now().Add(armBarrierTimeout)
	for time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("%s barrier for arm %s timed out after %s without being killed: AC4 is fail-closed, so a timeout FAILS the criterion rather than passing it", armBarrierEnv, arm, armBarrierTimeout)
}

// fileObservables returns the five live-target observables BG.C's runtime
// backstop (AC1b) compares across a gate arm -- content sha256, size, mode,
// inode and nanosecond ModTime -- plus the st_dev of the filesystem the path
// lives on. The stat comes from os.Stat + fi.Sys().(*syscall.Stat_t); darwin's
// Stat_t.Dev is int32 and linux's is uint64 (and darwin's Ino is uint64 while
// on some platforms it differs), so Dev and Ino are converted with explicit
// uint64() -- the conversion compiles on both, CI is linux and this host is
// darwin, and no build tags are used. mtimeNs uses fi.ModTime().UnixNano(),
// which carries the filesystem's nanosecond resolution on both platforms
// (darwin Mtimespec / linux Mtim) -- that is exactly the resolution the 20-trial
// probe below measures.
func fileObservables(path string) (sha string, size int64, mode os.FileMode, inode uint64, mtimeNs int64, dev uint64, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", 0, 0, 0, 0, 0, err
	}
	sha = digest(data)
	fi, err := os.Stat(path)
	if err != nil {
		return "", 0, 0, 0, 0, 0, err
	}
	size = fi.Size()
	mode = fi.Mode()
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return "", 0, 0, 0, 0, 0, fmt.Errorf("file %s: syscall.Stat_t unavailable", path)
	}
	inode = uint64(st.Ino)
	mtimeNs = fi.ModTime().UnixNano()
	dev = uint64(st.Dev)
	return sha, size, mode, inode, mtimeNs, dev, nil
}

// deviceOf returns st_dev for any path (file OR directory) via os.Stat. It
// exists because the granularity probe needs the st_dev of the arm's t.TempDir()
// directory and of repoRoot, which are directories and therefore cannot be
// read as files by fileObservables. Same uint64() portability discipline as
// fileObservables.
func deviceOf(path string) (uint64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, fmt.Errorf("path %s: syscall.Stat_t unavailable", path)
	}
	return uint64(st.Dev), nil
}

// granularityProbe measures the WORST CASE ModTime resolution of the filesystem
// the probe lives on. Under the arm's t.TempDir() it creates a probe file via
// confinedWrite (the single permitted write path, so the AST write-guard stays
// green) and runs granularityTrials back-to-back stat -> write-new ->
// write-original -> stat pairs with NO sleeps, counting how many trials the
// nanosecond ModTime actually changed. It records st_dev for the t.TempDir()
// path and for root so the transferability of the result can be asserted (W6).
func granularityProbe(t *testing.T, root, dir string) (fired int, tmpdirDev, repoDev uint64) {
	t.Helper()
	probe := filepath.Join(dir, "ac1b_mtime_probe")
	orig := []byte("ac1b mtime granularity probe -- original content\n")
	if err := confinedWrite(root, probe, orig, 0o600); err != nil {
		t.Fatal(err)
	}
	fired = 0
	for i := 0; i < granularityTrials; i++ {
		_, _, _, _, m1, _, err := fileObservables(probe)
		if err != nil {
			t.Fatal(err)
		}
		if err := confinedWrite(root, probe, []byte(fmt.Sprintf("ac1b probe trial %d different content\n", i)), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := confinedWrite(root, probe, orig, 0o600); err != nil {
			t.Fatal(err)
		}
		_, _, _, _, m2, _, err := fileObservables(probe)
		if err != nil {
			t.Fatal(err)
		}
		if m1 != m2 {
			fired++
		}
	}
	tmpdirDev, err := deviceOf(dir)
	if err != nil {
		t.Fatal(err)
	}
	repoDev, err = deviceOf(root)
	if err != nil {
		t.Fatal(err)
	}
	return fired, tmpdirDev, repoDev
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

	// BG.C AC1b runtime backstop -- BEFORE half. Taken immediately after absTarget
	// is known and the baseline is read, and BEFORE any overlay artifact is
	// written: this is the oldest honest snapshot of the live target on disk.
	beforeSHA, beforeSize, beforeMode, beforeInode, beforeMtimeNs, _, oerr := fileObservables(absTarget)
	if oerr != nil {
		t.Fatalf("AC1b arm=%s path=%s cannot stat live target before arm: %v", arm, rel, oerr)
	}
	if beforeSHA != baseSHA {
		t.Fatalf("AC1b arm=%s path=%s baseline read and live-target stat disagree (sha=%s vs stat=%s): the snapshot is not of the file the guard will protect", arm, rel, baseSHA, beforeSHA)
	}
	mutant := mutate(baseline)
	mutantSHA := digest(mutant)
	// Retained from the old harness: a mutation that never applied and a
	// mutation that failed to red are the same exit code.
	if mutantSHA == baseSHA {
		t.Fatalf("mutation did not apply to %s: sha256 remained %s", rel, baseSHA)
	}

	dir := t.TempDir()
	absMutant := filepath.Join(dir, "mutant__"+strings.ReplaceAll(rel, string(filepath.Separator), "__"))
	if err := confinedWrite(root, absMutant, mutant, 0o600); err != nil {
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
	if err := confinedWrite(root, absOverlayJSON, blob, 0o600); err != nil {
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

	// BG.C AC1b runtime backstop -- AFTER half. Runs on EVERY path, BEFORE the
	// guard-red fatals below: a backstop that only runs when the rest of the test
	// already passed could never report the interesting case (the guard failed to
	// red while the live file was nonetheless untouched).
	afterSHA, afterSize, afterMode, afterInode, afterMtimeNs, _, oerr := fileObservables(absTarget)
	if oerr != nil {
		t.Fatalf("AC1b arm=%s path=%s cannot stat live target after arm: %v", arm, rel, oerr)
	}

	// W3 -- four UNCONDITIONAL, filesystem-independent assertions. Inode is not
	// decoration: os.Rename is on the AST deny-list precisely because it is a
	// mutator, and a rename-based restore leaves content byte-identical and mtime
	// possibly unchanged while CHANGING the inode.
	if beforeSHA != afterSHA {
		t.Fatalf("AC1b arm=%s path=%s live-target sha256 changed: before=%s after=%s", arm, rel, beforeSHA, afterSHA)
	}
	if beforeSize != afterSize {
		t.Fatalf("AC1b arm=%s path=%s live-target size changed: before=%d after=%d", arm, rel, beforeSize, afterSize)
	}
	if beforeMode != afterMode {
		t.Fatalf("AC1b arm=%s path=%s live-target mode changed: before=%v after=%v", arm, rel, beforeMode, afterMode)
	}
	if beforeInode != afterInode {
		t.Fatalf("AC1b arm=%s path=%s live-target inode changed: before=%d after=%d (an os.Rename-based restore is exactly this signature)", arm, rel, beforeInode, afterInode)
	}

	// W7 -- no-write control: stat the same file twice with NO write between and
	// assert all five observables are equal. Without it, an always-unequal
	// comparator would make W3 red for the wrong reason and an always-equal one
	// would make it vacuously green.
	n1sha, n1size, n1mode, n1ino, n1mt, _, nerr := fileObservables(absTarget)
	if nerr != nil {
		t.Fatalf("AC1b no-write control: first stat failed: %v", nerr)
	}
	n2sha, n2size, n2mode, n2ino, n2mt, _, nerr := fileObservables(absTarget)
	if nerr != nil {
		t.Fatalf("AC1b no-write control: second stat failed: %v", nerr)
	}
	noWriteControl := n1sha == n2sha && n1size == n2size && n1mode == n2mode && n1ino == n2ino && n1mt == n2mt
	if !noWriteControl {
		t.Fatalf("AC1b no-write control failed: comparator reports a difference on a file no write touched (sha=%v size=%v mode=%v inode=%v mtime=%v)", n1sha == n2sha, n1size == n2size, n1mode == n2mode, n1ino == n2ino, n1mt == n2mt)
	}
	t.Logf("AC1b nowrite_control equal=%v", noWriteControl)

	// W4 -- the granularity probe, measured under the arm's t.TempDir() before the
	// ModTime comparison is believed.
	probeFired, tmpdirDev, repoDev := granularityProbe(t, root, dir)

	// W6 -- CONTROLLER ADDITION: the probe runs under t.TempDir() but the mtime
	// assertion it licenses is made about a file under repoRoot. If those are
	// different filesystems, a 20/20 probe on a fine-grained tmpfs would license an
	// mtime assertion on a coarse-grained repo volume -- a detector that cannot
	// detect, certified by a measurement of somewhere else. Checked BEFORE the
	// 20/20 gate: an untransferable probe is not a passing probe.
	if tmpdirDev != repoDev {
		t.Fatalf("AC1b granularity probe measured a different filesystem than the live target: probe dev=%d, live-target dev=%d; the probe result is not transferable -- see 10/OD-9", tmpdirDev, repoDev)
	}

	// W5 -- the 20/20 gate, with no third state and no silent state. fired < 20
	// FAILS the test; there is no logging-and-continuing branch. In CI a passing
	// test emits nothing, so a conditional backstop that logs and continues would
	// produce an outcome byte-identical to one that asserted -- "the check never
	// ran" wearing a green, which is this whole item's defect class. The threshold
	// is never lowered.
	if probeFired < granularityTrials {
		t.Fatalf("AC1b backstop is not armed on this filesystem: granularity probe fired %d/%d (tmpdir dev=%d, repo dev=%d); the ModTime half of the backstop cannot be relied on here -- see 10/OD-9", probeFired, granularityTrials, tmpdirDev, repoDev)
	}

	// The mtime comparator has now PROVEN it can fire (20/20) on the SAME
	// filesystem as the live target, so the nanosecond ModTime assertion is
	// licensed and believed.
	mtimeEqual := beforeMtimeNs == afterMtimeNs
	if !mtimeEqual {
		t.Fatalf("AC1b arm=%s path=%s live-target nanosecond ModTime changed: before=%d after=%d", arm, rel, beforeMtimeNs, afterMtimeNs)
	}

	// W8 -- evidence logs (visible only under -v, which is why every gate carries
	// -v). The probe X/20 and both st_dev values are the first measurement of
	// ModTime granularity on the CI filesystem in this mission's history.
	t.Logf("AC1b arm=%s path=%s probe_fired=%d/20 tmpdir_dev=%d repo_dev=%d", arm, rel, probeFired, tmpdirDev, repoDev)
	t.Logf("AC1b arm=%s path=%s sha256=%s size=%d mode=%v inode=%d mtime_ns_equal=%v", arm, rel, afterSHA, afterSize, afterMode, afterInode, mtimeEqual)

	if err == nil {
		t.Fatalf("mutation in %s passed boundary guard", rel)
	}
	if !strings.Contains(err.Error(), rel) {
		t.Fatalf("guard failure did not name exact path %s: %v", rel, err)
	}
	t.Logf("MUTATION path=%s baseline_sha256=%s mutant_sha256=%s guard=%q", rel, baseSHA, mutantSHA, err)
	return mutantSHA
}

// goArmMutate is the import-injecting mutation every protected Go arm applies,
// and ailArmMutate is the world arm's package-cache marker. Both are named
// rather than inlined so the gate and the recording-writer test below drive the
// IDENTICAL mutation: a recording test that mutated differently would be
// observing a write path the gate does not actually take.
func goArmMutate(t *testing.T, group goGroup) func([]byte) []byte {
	return func(baseline []byte) []byte {
		anchor := []byte("import (\n")
		if bytes.Count(baseline, anchor) != 1 {
			t.Fatalf("mutation anchor count for %s is %d, want 1", group.mutantFile, bytes.Count(baseline, anchor))
		}
		insert := []byte(fmt.Sprintf("import (\n\t_ %q // boundary mutation: compiling HTTP import\n", group.mutantImport))
		return bytes.Replace(baseline, anchor, insert, 1)
	}
}

func ailArmMutate(baseline []byte) []byte {
	return append(append([]byte(nil), baseline...), []byte("\n-- boundary mutation: package-cache .ailang/cache\n")...)
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
			mutateViaOverlay(t, root, group.mutantFile, group.name, goArmMutate(t, group), func(ov overlay) error {
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
		mutateViaOverlay(t, root, "world/types.ail", "world", ailArmMutate, func(ov overlay) error { return checkAILGroup(root, ov) })
	})

	// Green control: the broker is intentionally the network boundary. Its
	// non-empty closure must remain permitted by this inverse protected-group guard.
	//
	// UNTIL SM.C THIS CONTROL WAS VACUOUS AND THE CHARTER SAID SO. `len(deps) != 0`
	// is true of every Go package, so the control could only ever have proved that
	// `go list` ran — the positive half, that network code once it EXISTS in
	// host/broker is genuinely permitted there, was unproven because host/broker
	// had ZERO net/http in its production source and ZERO in its dependency
	// closure (measured at 2ef4a23: 0 files, 0/168 deps). SM.B2a did not discharge
	// it either: its POST happens inside the pinned `ailang` CHILD PROCESS, which
	// is not a Go dependency of anything here.
	//
	// SM.C's reconciler is the first IN-PROCESS net/http in host/broker, so the
	// control can finally be one that fails. Deleting the reconciler's import, or
	// moving the network boundary out of host/broker without moving this gate,
	// now reds here instead of passing silently.
	broker, err := goListDeps(root, "", "./host/broker/...")
	if err != nil {
		t.Fatal(err)
	}
	if len(broker) == 0 {
		t.Fatal("host/broker dependency enumeration is empty")
	}
	httpDeps := countDep(broker, "net/http")
	t.Logf("GREEN_CONTROL group=host/broker exact_count=%d net/http=%d result=PASS", len(broker), httpDeps)
	if httpDeps != 1 {
		t.Fatalf("host/broker closure contains %q %d time(s), want exactly 1: the network boundary "+
			"is supposed to LIVE here (host/broker/registry_reconcile.go's bounded read-only GET). "+
			"A zero here means either the reconciler lost its transport or the boundary moved out of "+
			"host/broker without this gate moving with it -- and the permissive half of AC12 would go "+
			"back to being unfalsifiable (closure=%d packages)", "net/http", httpDeps, len(broker))
	}
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

// TestWorldBoundaryRecordingWriter proves the write path of the mutation
// harness never touches repoRoot and always lands inside the acting arm's temp
// root. It swaps the SINK (rawWrite) -- never confinedWrite, so the confinement
// CHECK still runs in front of every write -- and drives THE REAL HARNESS,
// mutateViaOverlay, once per arm. A captured destination that resolves inside
// repoRoot, one that escapes the arm's temp root, or an arm that recorded ZERO
// destinations, must RED this test rather than pass vacuously.
//
// It drives mutateViaOverlay rather than replaying its write sequence inline,
// and that distinction is the whole value of the test. The first version of
// this test synthesised its own mutant and overlay paths and called
// confinedWrite directly; the controller MEASURED it and it PASSED 4/4 arms
// with mutateViaOverlay's own writes reverted to bare os.WriteFile -- i.e. it
// was blind to the harness whose write path it claimed to cover, and its
// exact-count assertion counted only the writes the test itself had just made.
// That is this sprint's own spine (a check shaped to itself tests itself, not
// the threat), so the assertion now counts the HARNESS's writes.
//
// The sink is a TEE, not a no-op: mutateViaOverlay validates its artifacts by
// reading them back off disk, so a sink that swallows writes would fail the
// harness for the wrong reason. It delegates to the CAPTURED original rather
// than to os.WriteFile, which would itself be an unpermitted call under the AST
// write-guard below.
func TestWorldBoundaryRecordingWriter(t *testing.T) {
	origRawWrite := rawWrite
	defer func() { rawWrite = origRawWrite }()

	root := repoRoot(t)

	// drive runs ONE real arm with the sink teed, and returns what the harness
	// actually wrote.
	drive := func(t *testing.T, rel, arm string, mutate func([]byte) []byte, check func(overlay) error) []string {
		t.Helper()
		var recorded []string
		rawWrite = func(dst string, data []byte, perm os.FileMode) error {
			recorded = append(recorded, dst)
			return origRawWrite(dst, data, perm)
		}
		defer func() { rawWrite = origRawWrite }()
		mutateViaOverlay(t, root, rel, arm, mutate, check)
		return recorded
	}

	assertArm := func(t *testing.T, arm string, recorded []string) {
		t.Helper()
		// (c) non-empty AND exact count. Zero destinations means the harness no
		// longer writes through the confined sink at all -- the exact regression
		// this test exists to catch -- so it REDS as vacuous rather than passing.
		if len(recorded) == 0 {
			t.Fatalf("%s: the harness recorded ZERO writes through the confined sink: either mutateViaOverlay bypasses rawWrite (regression) or this test no longer drives it (vacuous)", arm)
		}
		// Every arm writes through the confined sink: the harness's own two
		// artifacts (the mutant file and the overlay JSON), plus the BG.C AC1b
		// backstop's granularity probe under the arm's t.TempDir() -- ONE create
		// (the original content) and then TWO confined writes per trial (new
		// content, then the original written back) for granularityTrials trials.
		// All of them resolve beneath the arm's temp root, so the confinement
		// assertions below stay meaningful while the exact-count check stays
		// exact. Raising the count to include the probe is strict, not a
		// weakening: it verifies the BACKSTOP's own writes are confined too.
		const wantWrites = 2 + 1 + 2*granularityTrials
		if len(recorded) != wantWrites {
			t.Fatalf("%s: the harness recorded %d writes through the confined sink, want %d: %q", arm, len(recorded), wantWrites, recorded)
		}
		for _, dst := range recorded {
			// (b) no recorded destination may resolve beneath repoRoot.
			if inside, err := insideRepo(root, dst); err != nil || inside {
				t.Errorf("%s: recorded destination %s resolves beneath repoRoot (inside=%v, err=%v)", arm, dst, inside, err)
			}
			// (a) every recorded destination must resolve beneath a temp root,
			// never the repository -- asserted against the destination's own
			// prefix because t.TempDir() is owned by mutateViaOverlay.
			if inside, err := insideRepo(os.TempDir(), dst); err != nil || !inside {
				t.Errorf("%s: recorded destination %s does not resolve beneath the temp root %s (inside=%v, err=%v)", arm, dst, os.TempDir(), inside, err)
			}
		}
		sorted := append([]string(nil), recorded...)
		sort.Strings(sorted)
		t.Logf("RECORDED_WRITES arm=%s count=%d destinations=%q", arm, len(sorted), sorted)
	}

	for _, group := range protectedGoGroups {
		group := group
		t.Run("recording_"+strings.ReplaceAll(group.name, "/", "_"), func(t *testing.T) {
			recorded := drive(t, group.mutantFile, group.name, goArmMutate(t, group), func(ov overlay) error {
				_, err := checkGoGroup(root, group, ov)
				return err
			})
			assertArm(t, group.name, recorded)
		})
	}
	t.Run("recording_world", func(t *testing.T) {
		recorded := drive(t, "world/types.ail", "world", ailArmMutate, func(ov overlay) error { return checkAILGroup(root, ov) })
		assertArm(t, "world", recorded)
	})
}

// TestBoundaryASTWriteGuard walks the FULL AST of every .go file in
// host/boundary and REDS if any os.WriteFile / os.OpenFile / os.Create /
// os.Rename CALL appears anywhere except the single permitted site (the
// rawWrite initializer). It is deliberately an AST walk over identifiers and
// CallExprs, NOT a grep needle: iteration 54's textual self-guard in host/store
// was vacuous because its positive needle matched its own check line, and an
// AST walk cannot match its own source text at all.
//
// TWO STATED LIMITATIONS, recorded here so a green from this guard is never read
// as more than it is:
//
//  1. The deny-list matches a SELECTOR, `os.<Name>`, so it is bypassable by
//     import aliasing (`import w "os"; w.WriteFile(...)`), by reflection, and by
//     a write reached through a helper in another package. Resolving aliases
//     needs go/packages type information, which this walk deliberately does not
//     load. That residual escape is what mutation M7 and milestone BG.C exist to
//     cover with a RUNTIME backstop; this guard is the structural half only.
//  2. The guard covers WRITES only. A future check added inside checkGoGroup or
//     checkAILGroup that reads the disk DIRECTLY -- rather than through the
//     overlay read helper -- would silently not be exercised with the mutant, and
//     nothing here would notice: a checker reading past the thing it is meant to
//     inspect. That is a REVIEW RULE, not an assertion, and it is recorded as
//     such rather than as an acceptance criterion that cannot fail.
func TestBoundaryASTWriteGuard(t *testing.T) {
	// (iii) deny-list, enforced by length and by membership below.
	denyList := []string{"WriteFile", "OpenFile", "Create", "Rename"}
	if len(denyList) != 4 {
		t.Fatalf("AST deny-list has %d entries, want 4", len(denyList))
	}
	for _, name := range []string{"WriteFile", "OpenFile", "Create", "Rename"} {
		found := false
		for _, d := range denyList {
			if d == name {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("AST deny-list is missing %q (present: %q)", name, denyList)
		}
	}

	_, here, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate host/boundary directory")
	}
	dir := filepath.Dir(here)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var goFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".go") {
			goFiles = append(goFiles, e.Name())
		}
	}
	sort.Strings(goFiles)
	// (i) non-empty and equal to the explicit expected constant.
	const wantFileCount = 1
	if len(goFiles) != wantFileCount {
		t.Fatalf("host/boundary contains %d .go files, want %d: %q", len(goFiles), wantFileCount, goFiles)
	}
	t.Logf("AST_GUARD go_file_count=%d files=%q", len(goFiles), goFiles)

	permittedLine := -1
	permittedFound := false
	violations := 0
	for _, name := range goFiles {
		path := filepath.Join(dir, name)
		fset := token.NewFileSet()
		// UNTYPED nil src = read from disk (see (overlay).parseSrc: a typed nil
		// []byte is a non-nil interface and would parse as empty + EOF).
		tree, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		ast.Inspect(tree, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.CallExpr:
				sel, ok := node.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				pkg, ok := sel.X.(*ast.Ident)
				if !ok || pkg.Name != "os" {
					return true
				}
				for _, d := range denyList {
					if sel.Sel.Name == d {
						pos := fset.Position(node.Pos())
						t.Errorf("%s:%d: unpermitted os.%s CALL outside the rawWrite initializer", pos.Filename, pos.Line, d)
						violations++
						break
					}
				}
			case *ast.ValueSpec:
				// The SINGLE permitted site. `var rawWrite = os.WriteFile` is an
				// identifier REFERENCE, not a CallExpr, so the deny-list walk
				// above cannot and must not flag it. We assert separately that
				// the walker LOCATES this assignment and its line, so a walker
				// that finds as many permitted sites as violations is shown to
				// actually work before its green is believed.
				for i, id := range node.Names {
					if id.Name != "rawWrite" || i >= len(node.Values) {
						continue
					}
					sel, ok := node.Values[i].(*ast.SelectorExpr)
					if !ok {
						continue
					}
					pkg, ok := sel.X.(*ast.Ident)
					if !ok || pkg.Name != "os" || sel.Sel.Name != "WriteFile" {
						continue
					}
					permittedLine = fset.Position(node.Pos()).Line
					permittedFound = true
				}
			}
			return true
		})
	}

	// (ii) the walker must PROVE it found the permitted rawWrite initializer.
	if !permittedFound {
		t.Fatal("AST walker did not locate the permitted rawWrite = os.WriteFile initializer: a walker reporting as many permitted sites as violations is untested and its green is void")
	}
	t.Logf("AST_GUARD permitted_site=rawWrite line=%d", permittedLine)
	// (iv) enumeration: exact .go file count and the permitted call site's line.
	t.Logf("AST_GUARD enumeration go_file_count=%d permitted_call_line=%d", len(goFiles), permittedLine)

	if violations != 0 {
		t.Fatalf("AST write-guard found %d unpermitted os write call(s)", violations)
	}
}

// TestConfinedWriteRejectsInsideRepo is both the B4 rejection-path test and the
// M3 threat-shaped mutation. It routes today's live-path mutation at
// host/store/store.go through the writer, expects the exact mandated rejection
// message, and proves the rejection happens before a single byte can be written
// by asserting the target's sha256 is UNCHANGED. The green is carried by an
// EXERCISED rejection, never by the absence of observations.
func TestConfinedWriteRejectsInsideRepo(t *testing.T) {
	root := repoRoot(t)
	target := filepath.Join(root, "host", "store", "store.go")
	before, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	beforeSHA := digest(before)
	mutant := append(append([]byte(nil), before...), []byte("\n// M3 confined-writer rejection probe\n")...)

	// M3 threat-shaped mutation routed through the writer.
	err = confinedWrite(root, target, mutant, 0o600)
	if err == nil {
		t.Fatal("confinedWrite accepted a destination inside repoRoot: the rejection path is not exercised")
	}
	if !strings.Contains(err.Error(), "confined writer: destination inside repoRoot") {
		t.Fatalf("confinedWrite rejection message %q does not carry the mandated marker", err.Error())
	}
	predicted := "confined writer: destination inside repoRoot: host/store/store.go"
	if err.Error() != predicted {
		t.Fatalf("M3 predicted message %q but observed %q", predicted, err.Error())
	}
	t.Logf("M3 PREDICTED=%q OBSERVED=%q", predicted, err.Error())

	after, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if digest(after) != beforeSHA {
		t.Fatalf("confinedWrite rejected but %s changed sha256 before=%s after=%s", target, beforeSHA, digest(after))
	}
	t.Logf("M3 target=host/store/store.go sha256_after_rejection=%s UNCHANGED", beforeSHA)
}
