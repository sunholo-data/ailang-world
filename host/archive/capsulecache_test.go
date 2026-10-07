package archive

// Row 153 M3a: template promotion is racer-tolerant and EXDEV-safe.

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

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
