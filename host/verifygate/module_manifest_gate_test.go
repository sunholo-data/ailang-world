package verifygate

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	src := filepath.Join(repoRoot, filepath.FromSlash(rel))
	dst := filepath.Join(root, filepath.FromSlash(rel))
	in, err := os.Open(src)
	if err != nil {
		t.Fatalf("open copy source %s: %v", rel, err)
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		t.Fatal(err)
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
	if err != nil {
		t.Fatalf("create copy target %s: %v", rel, err)
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		t.Fatalf("copy %s: %v", rel, err)
	}
	if err := out.Close(); err != nil {
		t.Fatal(err)
	}
}

func newIsolatedGateRoot(t *testing.T) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), "iso")
	copyGateFile(t, root, "scripts/verify_ail.sh", 0o755)
	copyGateFile(t, root, "scripts/testdata/ailang_release_observed.txt", 0o644)
	for _, pattern := range []string{"world/*.ail", "design_docs/sketches/*.ail"} {
		matches, err := filepath.Glob(filepath.Join(repoRoot, filepath.FromSlash(pattern)))
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) == 0 {
			t.Fatalf("copy pattern %q matched zero files", pattern)
		}
		for _, src := range matches {
			rel, err := filepath.Rel(repoRoot, src)
			if err != nil {
				t.Fatal(err)
			}
			copyGateFile(t, root, filepath.ToSlash(rel), 0o644)
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
		t.Fatal(err)
	}
	if files != 13 || ailFiles != 11 {
		t.Fatalf("isolated copy landed %d files / %d .ail files, want 13 / 11", files, ailFiles)
	}
	return root
}

func runGateAt(t *testing.T, root string, env map[string]string) (int, string) {
	t.Helper()
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
		return 0, output.String()
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), output.String()
	}
	t.Fatalf("start isolated verify gate: %v", err)
	return -1, output.String()
}

func requirePristineControl(t *testing.T, root string) string {
	t.Helper()
	requirePinned(t)
	rc, out := runGateAt(t, root, map[string]string{
		"AILANG_BIN": pinned, "WORLD_PKG_AILANG_BIN": pinned,
	})
	const marker = "✓ 10/10 required world/ identities verified across 11 module(s)"
	if !strings.Contains(out, marker) {
		t.Fatalf("pristine isolated control missing %q (rc=%d)\n%s", marker, rc, out)
	}
	t.Logf("pristine control observed: %s", marker)
	return out
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
	if err := os.WriteFile(path, []byte(mutant), 0o755); err != nil {
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
	root := newIsolatedGateRoot(t)
	control := requirePristineControl(t, root)
	if got := strings.Count(control, "\n   ai-check "); got != 11 {
		t.Fatalf("pristine control emitted %d ai-check lines, want 11", got)
	}
	probe := filepath.Join(root, "world", "_stray_manifest_probe.ail")
	const source = "module world/_stray_manifest_probe\n\nexport func strayId(x: int) -> int = x\n"
	if err := os.WriteFile(probe, []byte(source), 0o644); err != nil {
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
	root := newIsolatedGateRoot(t)
	requirePristineControl(t, root)
	probe := filepath.Join(root, "world", "SNEAKY.AIL")
	const source = "module world/SNEAKY\n\nexport func sneakyId(x: int) -> int = x\n"
	if err := os.WriteFile(probe, []byte(source), 0o644); err != nil {
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

func TestModuleManifestRejectsDeletedModule(t *testing.T) {
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
	root := newIsolatedGateRoot(t)
	requirePristineControl(t, root)
	const old = `LEG1_MODULES=(
  design_docs/sketches/effectbroker.ail
  design_docs/sketches/logepoch.ail
  design_docs/sketches/storejournal.ail
  design_docs/sketches/transitions.ail
  design_docs/sketches/worlddapi.ail
  design_docs/sketches/worldkernel.ail
  design_docs/sketches/worldtypes.ail
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
