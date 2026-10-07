package archive

// Row 153 M3a: template promotion is racer-tolerant and EXDEV-safe.

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

func templateFixture(t *testing.T) (*Archive, hashref.HashRef, []byte) {
	t.Helper()
	return New(storeDBPath(t)), hashref.SumSHA256([]byte("interpreter")), []byte("module m\n")
}

// writeRunCache fakes a finished run's cache dir: compile/manifest.json plus
// one module file.
func writeRunCache(t *testing.T, dir, marker string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, "compile", "modules"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "compile", "manifest.json"), []byte(marker), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "compile", "modules", "m.bin"), []byte(marker), 0o600); err != nil {
		t.Fatal(err)
	}
}

func readTemplateManifest(t *testing.T, a *Archive, interp hashref.HashRef, src []byte) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(a.CapsuleTemplateDir(interp, src), "compile", "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

func noTemplateTmp(t *testing.T, a *Archive) {
	t.Helper()
	entries, err := os.ReadDir(a.CapsuleCacheRoot())
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".tmp-") {
			t.Fatalf("in-flight template tmp %q left in capsule-cache/", e.Name())
		}
	}
}

// AC3.1 (MUT-RACER-ERROR): a promotion that loses the rename to a racer that
// already promoted the same template is success: the loser's tmp is discarded
// and the winner's bytes are untouched.
func TestTemplateRenameRacerIsSuccess(t *testing.T) {
	a, interp, src := templateFixture(t)
	writeRunCache(t, a.CapsuleTemplateDir(interp, src), "winner") // the racer got there first
	if !a.CapsuleTemplateReady(interp, src) {
		t.Fatal("fixture: the winner's template is not ready")
	}

	tmp, err := a.NewCapsuleTemplateTmp()
	if err != nil {
		t.Fatal(err)
	}
	writeRunCache(t, tmp, "loser")
	if err := a.PromoteTemplate(tmp, interp, src); err != nil {
		t.Fatalf("PromoteTemplate lost to a racer: %v, want nil", err)
	}
	if _, err := os.Stat(tmp); !os.IsNotExist(err) {
		t.Fatalf("the loser's tmp %s survived (%v)", tmp, err)
	}
	if got := readTemplateManifest(t, a, interp, src); got != "winner" {
		t.Fatalf("the winner's manifest is %q, want it unchanged", got)
	}

	run := t.TempDir()
	writeRunCache(t, run, "loser2")
	if err := a.PromoteCopy(run, interp, src); err != nil {
		t.Fatalf("PromoteCopy lost to a racer: %v, want nil", err)
	}
	if got := readTemplateManifest(t, a, interp, src); got != "winner" {
		t.Fatalf("after PromoteCopy the manifest is %q, want the winner's", got)
	}
	noTemplateTmp(t, a)
}

// AC3.1 (MUT-CROSS-DEVICE-RENAME): a lazy promotion copies the run's cache
// from another filesystem. The rename seam is stubbed to fail EXDEV for any
// source outside capsule-cache/, as a rename from a run root on another
// device would (the plan's "Dir(old) != Dir(new)" stub also rejects the
// legitimate tmp -> <digest>/<hash> rename, which crosses directories on one
// device); PromoteCopy must stay inside capsule-cache/ and succeed.
func TestCapsulePromotionSurvivesCrossDevice(t *testing.T) {
	a, interp, src := templateFixture(t)
	old := renameDir
	t.Cleanup(func() { renameDir = old })
	cacheRoot := a.CapsuleCacheRoot() + string(filepath.Separator)
	renameDir = func(from, to string) error {
		if !strings.HasPrefix(from, cacheRoot) {
			return &os.LinkError{Op: "rename", Old: from, New: to, Err: syscall.EXDEV}
		}
		return old(from, to)
	}
	run := t.TempDir() // a separate root, as the capsule's run root is
	writeRunCache(t, run, "cold-run")
	if err := a.PromoteCopy(run, interp, src); err != nil {
		t.Fatalf("PromoteCopy across devices: %v", err)
	}
	if !a.CapsuleTemplateReady(interp, src) {
		t.Fatal("no ready template after a cross-device promotion")
	}
	if got := readTemplateManifest(t, a, interp, src); got != "cold-run" {
		t.Fatalf("template manifest = %q", got)
	}
	noTemplateTmp(t, a)

	// The EXDEV the stub models is a real, typed failure when it is NOT
	// avoided: a direct cross-directory promotion is a TemplateBuildError.
	tmp := t.TempDir()
	writeRunCache(t, tmp, "direct")
	a2, interp2, src2 := templateFixture(t)
	err := a2.PromoteTemplate(tmp, interp2, src2)
	var build *TemplateBuildError
	if !errors.As(err, &build) || !errors.Is(err, syscall.EXDEV) {
		t.Fatalf("direct cross-device PromoteTemplate = %v, want *TemplateBuildError wrapping EXDEV", err)
	}
}

// P11: the keep set names every digest GC must leave (the daemon pin and each
// head descriptor's interpreter); unlisted digests, and only digest-named
// directories, go.
func TestPruneCapsuleCacheKeepsListedDigests(t *testing.T) {
	a := New(storeDBPath(t))
	keepA, keepB, drop := strings.Repeat("a", 64), strings.Repeat("b", 64), strings.Repeat("c", 64)
	for _, name := range []string{keepA, keepB, drop, "not-a-digest"} {
		if err := os.MkdirAll(filepath.Join(a.CapsuleCacheRoot(), name, "src"), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	removed, err := a.PruneCapsuleCache(map[string]bool{keepA: true, keepB: true}, StaleTemplateTmpAge, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if len(removed) != 1 || removed[0] != drop {
		t.Fatalf("removed = %v, want only %s", removed, drop)
	}
	for _, name := range []string{keepA, keepB, "not-a-digest"} {
		if _, err := os.Stat(filepath.Join(a.CapsuleCacheRoot(), name)); err != nil {
			t.Errorf("%s was pruned: %v", name, err)
		}
	}
}

// Judge N4 (OWN-4): CopyCacheTree refuses a symlink (here a link out of the
// tree to a real file), so a template cannot smuggle a path out of the run's
// root; PromoteCopy therefore promotes nothing from such a run.
func TestCopyCacheTreeRefusesSymlinks(t *testing.T) {
	outside := filepath.Join(t.TempDir(), "secret")
	if err := os.WriteFile(outside, []byte("secret"), 0o600); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(t.TempDir(), "run")
	writeRunCache(t, src, "run")
	if err := os.Symlink(outside, filepath.Join(src, "compile", "modules", "link.bin")); err != nil {
		t.Skipf("symlinks unsupported: %v", err)
	}

	dst := t.TempDir()
	err := CopyCacheTree(src, dst)
	if err == nil || !strings.Contains(err.Error(), "refusing to copy") {
		t.Fatalf("CopyCacheTree over a symlink = %v, want a refusal", err)
	}
	if _, serr := os.Stat(filepath.Join(dst, "compile", "modules", "link.bin")); serr == nil {
		t.Fatal("the symlink's target was copied into the destination")
	}

	a, interp, source := templateFixture(t)
	if perr := a.PromoteCopy(src, interp, source); perr == nil {
		t.Fatal("PromoteCopy of a run holding a symlink succeeded")
	}
	if a.CapsuleTemplateReady(interp, source) {
		t.Fatal("a template was promoted from a run holding a symlink")
	}
	noTemplateTmp(t, a)
}

// Judge N6 (OWN-7): a template counts as ready only if compile/manifest.json
// exists as a regular file; a bare directory (or one without the manifest) is
// not a template.
func TestCapsuleTemplateReadyNeedsTheManifest(t *testing.T) {
	a, interp, source := templateFixture(t)
	if a.CapsuleTemplateReady(interp, source) {
		t.Fatal("ready with no template directory")
	}
	dir := a.CapsuleTemplateDir(interp, source)
	if err := os.MkdirAll(filepath.Join(dir, "compile", "modules"), 0o700); err != nil {
		t.Fatal(err)
	}
	if a.CapsuleTemplateReady(interp, source) {
		t.Fatal("ready for a directory without compile/manifest.json")
	}
	if err := os.Mkdir(filepath.Join(dir, "compile", "manifest.json"), 0o700); err != nil {
		t.Fatal(err)
	}
	if a.CapsuleTemplateReady(interp, source) {
		t.Fatal("ready when manifest.json is a directory")
	}
	if err := os.Remove(filepath.Join(dir, "compile", "manifest.json")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "compile", "manifest.json"), []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
	if !a.CapsuleTemplateReady(interp, source) {
		t.Fatal("not ready with a regular compile/manifest.json")
	}
}
