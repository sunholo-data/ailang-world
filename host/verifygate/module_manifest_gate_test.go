package verifygate

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

var liveStorejournalDigest = func() [sha256.Size]byte {
	raw, err := os.ReadFile(filepath.Join(repoRoot, "design_docs", "sketches", "storejournal.ail"))
	if err != nil {
		panic(fmt.Sprintf("read live storejournal baseline: %v", err))
	}
	return sha256.Sum256(raw)
}()

func copyGateFile(t *testing.T, root, rel string, mode os.FileMode) {
	t.Helper()
	if err := copyGateFileErr(root, rel, mode); err != nil {
		t.Fatal(err)
	}
}

func copyGateFileErr(root, rel string, mode os.FileMode) error {
	src := filepath.Join(repoRoot, filepath.FromSlash(rel))
	dst := filepath.Join(root, filepath.FromSlash(rel))
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open copy source %s: %v", rel, err)
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return copyFileForkLocked(in, dst, rel, mode)
}

func newIsolatedGateRoot(t *testing.T) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), "iso")
	if err := buildIsolatedGateRoot(root); err != nil {
		t.Fatal(err)
	}
	return root
}

// buildIsolatedGateRoot is newIsolatedGateRoot without a *testing.T, so the ONE shared pristine
// control (sharedPristineControl) can build its own root outside any single test's lifetime.
func buildIsolatedGateRoot(root string) error {
	if err := copyGateFileErr(root, "scripts/verify_ail.sh", 0o755); err != nil {
		return err
	}
	if err := copyGateFileErr(root, "scripts/testdata/ailang_release_observed.txt", 0o644); err != nil {
		return err
	}
	for _, pattern := range []string{"world/*.ail", "design_docs/sketches/*.ail", "packages/se-tools/se_tools/*.ail"} {
		matches, err := filepath.Glob(filepath.Join(repoRoot, filepath.FromSlash(pattern)))
		if err != nil {
			return err
		}
		if len(matches) == 0 {
			return fmt.Errorf("copy pattern %q matched zero files", pattern)
		}
		for _, src := range matches {
			rel, err := filepath.Rel(repoRoot, src)
			if err != nil {
				return err
			}
			if err := copyGateFileErr(root, filepath.ToSlash(rel), 0o644); err != nil {
				return err
			}
		}
	}
	files, ailFiles := 0, 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files++
			if filepath.Ext(path) == ".ail" {
				ailFiles++
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if files != 22 || ailFiles != 20 {
		return fmt.Errorf("isolated copy landed %d files / %d .ail files, want 22 / 20", files, ailFiles)
	}
	return nil
}

// isolatedTreeDigest is a content+mode digest of every regular file under root, keyed by its
// root-relative path in sorted (Walk) order. Two roots with equal digests are byte- and
// mode-identical gate inputs, which is what lets one pristine control stand for every root.
func isolatedTreeDigest(root string) (string, int, error) {
	h := sha256.New()
	n := 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(raw)
		fmt.Fprintf(h, "%s\x00%o\x00%x\n", filepath.ToSlash(rel), info.Mode().Perm(), sum)
		n++
		return nil
	})
	return fmt.Sprintf("%x", h.Sum(nil)), n, err
}

func runGateAt(t *testing.T, root string, env map[string]string) (int, string) {
	t.Helper()
	rc, out, err := runGateAtErr(root, env)
	if err != nil {
		t.Fatalf("start isolated verify gate: %v", err)
	}
	return rc, out
}

func runGateAtErr(root string, env map[string]string) (int, string, error) {
	cmd := exec.Command(filepath.Join(root, "scripts", "verify_ail.sh"))
	cmd.Dir = root
	blocked := map[string]bool{
		"AILANG_BIN": true, "WORLD_PKG_AILANG_BIN": true,
		"AILANG_SHIM_VERSION_LINE": true, "AILANG_SHIM_DELEGATE": true,
		"AILANG_Z3_PATH": true,
	}
	cmd.Env = make([]string, 0, len(os.Environ())+len(env))
	for _, item := range os.Environ() {
		if !blocked[strings.SplitN(item, "=", 2)[0]] {
			cmd.Env = append(cmd.Env, item)
		}
	}
	for key, value := range env {
		cmd.Env = append(cmd.Env, key+"="+value)
	}
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	err := cmd.Run()
	if err == nil {
		return 0, output.String(), nil
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), output.String(), nil
	}
	return -1, output.String(), err
}

// sharedPristineControl is the ONE pristine-control gate run that every mutation arm in this file
// stands on (CI budget fix, 2026-10-03). Each arm used to re-run the FULL gate on its own fresh,
// unmutated copy before mutating it: 8 identical full runs (~16 s each locally, Leg 2b ~12 s of
// that) on byte-identical inputs, the largest single cost in host/verifygate and the reason the
// -race leg overran its 600 s budget on the 2-core runner once row 135 grew the se-tools tests.
// The claim each arm needs is unchanged -- "THIS root, unmutated, satisfies the control" -- and is
// now proven as two facts instead of one rerun: the shared run on a freshly built root carries the
// marker, AND the arm's own root is digest-identical (every file's path, mode and sha256) to that
// root before the arm mutates it. A root that differs in any byte is refused, never assumed equal.
var sharedPristine struct {
	once   sync.Once
	digest string
	files  int
	rc     int
	out    string
	err    error
}

func sharedPristineControl() (digest string, files, rc int, out string, err error) {
	s := &sharedPristine
	s.once.Do(func() {
		dir, err := os.MkdirTemp("", "verifygate-pristine-")
		if err != nil {
			s.err = err
			return
		}
		defer os.RemoveAll(dir)
		root := filepath.Join(dir, "iso")
		if s.err = buildIsolatedGateRoot(root); s.err != nil {
			return
		}
		if s.digest, s.files, s.err = isolatedTreeDigest(root); s.err != nil {
			return
		}
		s.rc, s.out, s.err = runGateAtErr(root, map[string]string{
			"AILANG_BIN": pinned, "WORLD_PKG_AILANG_BIN": pinned,
		})
	})
	return s.digest, s.files, s.rc, s.out, s.err
}

func requirePristineControl(t *testing.T, root string) string {
	t.Helper()
	requirePinned(t)
	out, digest, err := pristineControlCovers(root)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("pristine control observed (shared run, digest-identical root %s): %s", digest[:12], pristineMarker)
	return out
}

const pristineMarker = "✓ 16/16 required world/ identities verified across 20 module(s)"

// pristineControlCovers is requirePristineControl's contract as a predicate: it returns the shared
// control's output only when that run carries the marker AND root is digest-identical to the root
// it ran on. TestPristineControlRefusesANonIdenticalRoot points it at altered roots.
func pristineControlCovers(root string) (string, string, error) {
	digest, files, rc, out, err := sharedPristineControl()
	if err != nil {
		return "", "", fmt.Errorf("shared pristine control could not run: %v", err)
	}
	if files != 22 {
		return "", "", fmt.Errorf("shared pristine control digested %d files, want 22", files)
	}
	mine, mineFiles, err := isolatedTreeDigest(root)
	if err != nil {
		return "", "", err
	}
	if mine != digest || mineFiles != files {
		return "", "", fmt.Errorf("this arm's isolated root (%d files, digest %s) is NOT identical to the root the shared "+
			"pristine control ran on (%d files, digest %s) -- the control does not cover it", mineFiles, mine, files, digest)
	}
	if !strings.Contains(out, pristineMarker) {
		return "", "", fmt.Errorf("pristine isolated control missing %q (rc=%d)\n%s", pristineMarker, rc, out)
	}
	return out, digest, nil
}

// TestPristineControlRefusesANonIdenticalRoot is the non-vacuity arm for the shared pristine
// control: a fresh root is covered, and each single-fact divergence (one content byte, one mode
// bit, one extra file) is refused. Without it, a digest that ignored content or mode would let one
// control run silently stand for a root it never saw.
func TestPristineControlRefusesANonIdenticalRoot(t *testing.T) {
	// Parallel-safe: this arm runs the gate only in its OWN t.TempDir() copy (own .ailang compile
	// cache), never on the live tree; see sharedPristineControl for the CI budget rationale.
	t.Parallel()
	requirePinned(t)
	if _, _, err := pristineControlCovers(newIsolatedGateRoot(t)); err != nil {
		t.Fatalf("positive control: a fresh isolated root is not covered: %v", err)
	}
	for _, arm := range []struct {
		name  string
		alter func(t *testing.T, root string)
	}{
		{"content-byte", func(t *testing.T, root string) {
			p := filepath.Join(root, "world", "types.ail")
			raw, err := os.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}
			if err := writeFileForkLocked(p, append(raw, ' '), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{"mode-bit", func(t *testing.T, root string) {
			if err := os.Chmod(filepath.Join(root, "scripts", "verify_ail.sh"), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{"extra-file", func(t *testing.T, root string) {
			if err := writeFileForkLocked(filepath.Join(root, "world", "extra.txt"), nil, 0o644); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(arm.name, func(t *testing.T) {
			root := newIsolatedGateRoot(t)
			arm.alter(t, root)
			if _, _, err := pristineControlCovers(root); err == nil || !strings.Contains(err.Error(), "NOT identical") {
				t.Fatalf("altered root (%s) was covered by the shared control: err=%v", arm.name, err)
			}
		})
	}
}

func mutateCopiedScript(t *testing.T, root, old, replacement string) {
	t.Helper()
	path := filepath.Join(root, "scripts", "verify_ail.sh")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(string(raw), old); n != 1 {
		t.Fatalf("copied-script mutation anchor count=%d, want 1 for %q", n, old)
	}
	mutant := strings.Replace(string(raw), old, replacement, 1)
	if err := writeFileForkLocked(path, []byte(mutant), 0o755); err != nil {
		t.Fatal(err)
	}
}

func requireLiveTreeUntouched(t *testing.T) {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(repoRoot, "world", "_stray*"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 0 {
		t.Fatalf("committed arm wrote live stray files: %v", matches)
	}
	raw, err := os.ReadFile(filepath.Join(repoRoot, "design_docs", "sketches", "storejournal.ail"))
	if err != nil {
		t.Fatalf("read live storejournal after arm: %v", err)
	}
	if got := sha256.Sum256(raw); got != liveStorejournalDigest {
		t.Fatalf("committed arm changed live storejournal: got sha256 %x, baseline %x", got, liveStorejournalDigest)
	}
}

func TestModuleManifestRejectsStrayModule(t *testing.T) {
	// Parallel-safe: this arm runs the gate only in its OWN t.TempDir() copy (own .ailang compile
	// cache), never on the live tree; see sharedPristineControl for the CI budget rationale.
	t.Parallel()
	root := newIsolatedGateRoot(t)
	control := requirePristineControl(t, root)
	if got := strings.Count(control, "\n   ai-check "); got != 20 {
		t.Fatalf("pristine control emitted %d ai-check lines, want 20", got)
	}
	probe := filepath.Join(root, "world", "_stray_manifest_probe.ail")
	const source = "module world/_stray_manifest_probe\n\nexport func strayId(x: int) -> int = x\n"
	if err := writeFileForkLocked(probe, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	rc, out := runGateAt(t, root, map[string]string{
		"AILANG_BIN": pinned, "WORLD_PKG_AILANG_BIN": pinned,
	})
	if rc != 1 {
		t.Fatalf("stray module: want rc=1, got %d\n%s", rc, out)
	}
	for _, want := range []string{
		"LEG1_MODULES", "+world/_stray_manifest_probe.ail",
		"expected: LEG1_MODULES", "actual:   .ail files swept",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("stray refusal omits %q\n%s", want, out)
		}
	}
	if got := strings.Count(out, "\n   ai-check "); got != 0 {
		t.Fatalf("stray was ai-checked before refusal: count=%d\n%s", got, out)
	}
	if strings.Contains(out, "verify gate PASSED") {
		t.Fatalf("stray refusal printed terminal success\n%s", out)
	}
	requireLiveTreeUntouched(t)
}

// TestModuleManifestRejectsCaseVariantExtension pins the spelling that defeated the first cut of
// this gate. The enumeration used `find -name '*.ail'`, which is case-SENSITIVE, so a module named
// `SNEAKY.AIL` never entered the swept set — and a set the enumeration cannot see is a set the
// compare cannot refuse. Measured on the landed-then-repaired tree: with `-name`, the gate exits 0
// and prints "swept .ail module set equals the LEG1_MODULES allowlist (11 modules)" while the file
// sits in world/; with `-iname` it is refused before any ai-check runs.
//
// The assertion deliberately reads the offending path out of the gate's OWN diff output rather than
// reconstructing the expected text, so a gate that refuses for some unrelated reason cannot pass it.
func TestModuleManifestRejectsCaseVariantExtension(t *testing.T) {
	// Parallel-safe: this arm runs the gate only in its OWN t.TempDir() copy (own .ailang compile
	// cache), never on the live tree; see sharedPristineControl for the CI budget rationale.
	t.Parallel()
	root := newIsolatedGateRoot(t)
	requirePristineControl(t, root)
	probe := filepath.Join(root, "world", "SNEAKY.AIL")
	const source = "module world/SNEAKY\n\nexport func sneakyId(x: int) -> int = x\n"
	if err := writeFileForkLocked(probe, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	rc, out := runGateAt(t, root, map[string]string{
		"AILANG_BIN": pinned, "WORLD_PKG_AILANG_BIN": pinned,
	})
	if rc != 1 {
		t.Fatalf("case-variant extension: want rc=1, got %d\n%s", rc, out)
	}
	if !strings.Contains(out, "+world/SNEAKY.AIL") {
		t.Fatalf("case-variant refusal omits the offending path from the manifest diff\n%s", out)
	}
	if strings.Contains(out, "verify gate PASSED") {
		t.Fatalf("case-variant refusal printed terminal success\n%s", out)
	}
	if got := strings.Count(out, "\n   ai-check "); got != 0 {
		t.Fatalf("case-variant was ai-checked before refusal: count=%d\n%s", got, out)
	}
	requireLiveTreeUntouched(t)
}

// TestEffectPlanSketchIdentityIsRequired is row 134 AC1.2's mutation: the
// effect-plan sketch's proofs are gated BY NAME, so dropping one contract (here
// L4's payloadSizeOk ensures) leaves the module compiling but reds the gate on
// the vanished identity — it is not absorbed by the world/-only total.
func TestEffectPlanSketchIdentityIsRequired(t *testing.T) {
	// Parallel-safe: this arm runs the gate only in its OWN t.TempDir() copy (own .ailang compile
	// cache), never on the live tree; see sharedPristineControl for the CI budget rationale.
	t.Parallel()
	root := newIsolatedGateRoot(t)
	requirePristineControl(t, root)
	target := filepath.Join(root, "design_docs", "sketches", "effectplan.ail")
	raw, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	const contract = "ensures { result == (n >= 0 && n <= 1048576) }\n"
	if n := strings.Count(string(raw), contract); n != 1 {
		t.Fatalf("payloadSizeOk contract anchor count=%d, want 1", n)
	}
	if err := writeFileForkLocked(target, []byte(strings.Replace(string(raw), contract, "", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	rc, out := runGateAt(t, root, map[string]string{
		"AILANG_BIN": pinned, "WORLD_PKG_AILANG_BIN": pinned,
	})
	if rc != 1 || !strings.Contains(out, "payloadSizeOk) MISSING from verify.results[]") {
		t.Fatalf("dropped effectplan contract not refused by name: rc=%d\n%s", rc, out)
	}
	if strings.Contains(out, "verify gate PASSED") {
		t.Fatalf("dropped effectplan contract printed terminal success\n%s", out)
	}
	requireLiveTreeUntouched(t)
}

// TestSeToolsGateArmsAreLoadBearing is row 134 M5's pair of gate mutations on
// the se-tools package: (contract) dropping the L-ARGS key law's ensures from
// se_tools/read.ail leaves the module compiling but reds Leg 1 on the vanished
// identity by name; (plan bytes) changing the cost in read's pinned admitted-plan
// expectation reds Leg 2b on main_test_1 by name. Neither is absorbed by the
// world/-only total.
func TestSeToolsGateArmsAreLoadBearing(t *testing.T) {
	// Parallel-safe: this arm runs the gate only in its OWN t.TempDir() copy (own .ailang compile
	// cache), never on the live tree; see sharedPristineControl for the CI budget rationale.
	t.Parallel()
	arms := []struct {
		name, old, replacement, want string
	}{
		{"contract", "ensures { result == (k == \"path\") }\n", "",
			"argKeyAllowed) MISSING from verify.results[]"},
		{"plan-bytes", `cost\":1,\"payload`, `cost\":0,\"payload`,
			"se_tools/read.ail required named tests missing/failing: main_test_1="},
	}
	for _, arm := range arms {
		t.Run(arm.name, func(t *testing.T) {
			t.Parallel()
			root := newIsolatedGateRoot(t)
			requirePristineControl(t, root)
			target := filepath.Join(root, "packages", "se-tools", "se_tools", "read.ail")
			raw, err := os.ReadFile(target)
			if err != nil {
				t.Fatal(err)
			}
			if n := strings.Count(string(raw), arm.old); n != 1 {
				t.Fatalf("mutation anchor count=%d, want 1 for %q", n, arm.old)
			}
			if err := writeFileForkLocked(target, []byte(strings.Replace(string(raw), arm.old, arm.replacement, 1)), 0o644); err != nil {
				t.Fatal(err)
			}
			rc, out := runGateAt(t, root, map[string]string{
				"AILANG_BIN": pinned, "WORLD_PKG_AILANG_BIN": pinned,
			})
			if rc != 1 || !strings.Contains(out, arm.want) {
				t.Fatalf("se-tools %s mutant not refused by name (want %q): rc=%d\n%s", arm.name, arm.want, rc, out)
			}
			if strings.Contains(out, "verify gate PASSED") {
				t.Fatalf("se-tools %s mutant printed terminal success\n%s", arm.name, out)
			}
			requireLiveTreeUntouched(t)
		})
	}
}

func TestModuleManifestRejectsDeletedModule(t *testing.T) {
	// Parallel-safe: this arm runs the gate only in its OWN t.TempDir() copy (own .ailang compile
	// cache), never on the live tree; see sharedPristineControl for the CI budget rationale.
	t.Parallel()
	root := newIsolatedGateRoot(t)
	requirePristineControl(t, root)
	target := filepath.Join(root, "design_docs", "sketches", "storejournal.ail")
	if err := os.Remove(target); err != nil {
		t.Fatal(err)
	}
	rc, out := runGateAt(t, root, map[string]string{
		"AILANG_BIN": pinned, "WORLD_PKG_AILANG_BIN": pinned,
	})
	if rc != 1 || !strings.Contains(out, "-design_docs/sketches/storejournal.ail") {
		t.Fatalf("deleted leaf not refused by named manifest diff: rc=%d\n%s", rc, out)
	}
	requireLiveTreeUntouched(t)
}

func TestModuleManifestEmptyAllowlistFailsLoudly(t *testing.T) {
	// Parallel-safe: this arm runs the gate only in its OWN t.TempDir() copy (own .ailang compile
	// cache), never on the live tree; see sharedPristineControl for the CI budget rationale.
	t.Parallel()
	root := newIsolatedGateRoot(t)
	requirePristineControl(t, root)
	const old = `LEG1_MODULES=(
  design_docs/sketches/effectbroker.ail
  design_docs/sketches/effectplan.ail
  design_docs/sketches/logepoch.ail
  design_docs/sketches/storejournal.ail
  design_docs/sketches/transitions.ail
  design_docs/sketches/worlddapi.ail
  design_docs/sketches/worldkernel.ail
  design_docs/sketches/worldtypes.ail
  packages/se-tools/se_tools/builtins_search.ail
  packages/se-tools/se_tools/check.ail
  packages/se-tools/se_tools/cli.ail
  packages/se-tools/se_tools/edit.ail
  packages/se-tools/se_tools/examples_search.ail
  packages/se-tools/se_tools/read.ail
  packages/se-tools/se_tools/run.ail
  packages/se-tools/se_tools/write.ail
  world/contracts.ail
  world/logepoch.ail
  world/transitions.ail
  world/types.ail
)`
	mutateCopiedScript(t, root, old, "LEG1_MODULES=()")
	rc, out := runGateAt(t, root, map[string]string{
		"AILANG_BIN": pinned, "WORLD_PKG_AILANG_BIN": pinned,
	})
	if rc != 1 || !strings.Contains(out, "LEG1_MODULES allowlist is empty") {
		t.Fatalf("empty allowlist did not fail through its named guard: rc=%d\n%s", rc, out)
	}
	if strings.Contains(out, "unbound variable") {
		t.Fatalf("empty allowlist aborted before its own guard\n%s", out)
	}
	requireLiveTreeUntouched(t)
}

func TestModuleManifestEmptyEnumerationFailsLoudly(t *testing.T) {
	// Parallel-safe: this arm runs the gate only in its OWN t.TempDir() copy (own .ailang compile
	// cache), never on the live tree; see sharedPristineControl for the CI budget rationale.
	t.Parallel()
	root := newIsolatedGateRoot(t)
	requirePristineControl(t, root)
	const old = `mods+=("${f#./}")            # repo-relative path (manifest key), normalized (gemini catch)`
	const replacement = `mods+=("${f#./}"); mods=() # repo-relative path (manifest key), normalized (gemini catch)`
	mutateCopiedScript(t, root, old, replacement)
	rc, out := runGateAt(t, root, map[string]string{
		"AILANG_BIN": pinned, "WORLD_PKG_AILANG_BIN": pinned,
	})
	if rc != 1 || !strings.Contains(out, "swept .ail enumeration was empty") {
		t.Fatalf("empty enumeration did not fail through its named guard: rc=%d\n%s", rc, out)
	}
	requireLiveTreeUntouched(t)
}
