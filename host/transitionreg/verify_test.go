package transitionreg

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/registry"
	"github.com/sunholo-data/ailang-world/host/store"
)

// testFakeRelease is the release string the default fake interpreter reports.
const testFakeRelease = "TEST-FAKE v9.9.9"

// fakeInterpreterScript writes the house-pattern shell-script fake
// interpreter (host/archive/archive_test.go): --version prints version;
// `check` exits with checkExit (0 accepts every source, 1 refuses all).
func fakeInterpreterScript(t *testing.T, version string, checkExit int) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "fake-ailang")
	script := "#!/bin/sh\n" +
		"if [ \"$1\" = \"--version\" ]; then\n" +
		"  echo \"" + version + "\"\n" +
		"  exit 0\n" +
		"fi\n" +
		"if [ \"$1\" = \"check\" ]; then\n" +
		"  echo \"Error: this fake refuses the check\" >&2\n" +
		"  exit " + strconv.Itoa(checkExit) + "\n" +
		"fi\n" +
		"exit 0\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

// publisherStore opens a FILE-BACKED temp store (the archive roots at
// <db>.artifacts), archives a fake interpreter reporting interpreterRelease,
// and bootstraps the epoch registry with bootstrapRelease ("" skips the
// bootstrap). This is the daemon's own startup shape: archive the interpreter,
// then registry.Bootstrap(ctx, s, release) (daemon.go:494-512).
func publisherStore(t *testing.T, interpreterRelease, bootstrapRelease string) (*store.Store, *archive.Archive, hashref.HashRef) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "world.db")
	s, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	arch := archive.New(dbPath)
	ref, err := arch.Archive(fakeInterpreterScript(t, interpreterRelease, 0))
	if err != nil {
		t.Fatalf("archive fake interpreter: %v", err)
	}
	if bootstrapRelease != "" {
		if _, _, err := registry.Bootstrap(context.Background(), s, bootstrapRelease); err != nil {
			t.Fatalf("bootstrap epoch registry: %v", err)
		}
	}
	return s, arch, ref
}

// TestEpochsForInterpreterMatchesBootstrapRelease measures the derivation
// itself: the manifest's verbatim --version output reduces to the release
// string the daemon bootstrapped (first non-blank line, trimmed), and that
// release names exactly the bootstrapped epoch.
func TestEpochsForInterpreterMatchesBootstrapRelease(t *testing.T) {
	s, arch, ref := publisherStore(t, testFakeRelease, testFakeRelease)
	epochs, release, err := NewPublisher(s, arch).EpochsForInterpreter(context.Background(), ref)
	if err != nil {
		t.Fatalf("EpochsForInterpreter: %v", err)
	}
	if release != testFakeRelease {
		t.Fatalf("derived release = %q, want %q (the daemon's reduction of the manifest version)", release, testFakeRelease)
	}
	if len(epochs) != 1 || epochs[0] != 1 {
		t.Fatalf("derived epochs = %v, want [1] (registry.Bootstrap created epoch 1)", epochs)
	}
}

// TestReleaseFromManifestMirrorsDaemonReduction pins the reduction rule the
// daemon applies at bootstrap (daemon.go releaseFromVersion): first
// non-blank line, trimmed; "unpinned" only when there is none.
func TestReleaseFromManifestMirrorsDaemonReduction(t *testing.T) {
	cases := []struct{ in, want string }{
		{"TEST-FAKE v9.9.9\n", "TEST-FAKE v9.9.9"},
		{"\n\n  AILANG v0.41.0  \nCommit: 24ee108\n", "AILANG v0.41.0"},
		{"   \n\t\n", "unpinned"},
	}
	for _, c := range cases {
		if got := releaseFromManifest(c.in); got != c.want {
			t.Errorf("releaseFromManifest(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestEnsureSourceLoadableStagesHermetically pins the check's honest scope
// directly: the source is staged in a scratch root and checked with the
// archived interpreter — the accepting fake passes, the refusing fake fails
// with its captured output in the error.
func TestEnsureSourceLoadableStagesHermetically(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "world.db")
	arch := archive.New(dbPath)
	accepting, err := arch.Archive(fakeInterpreterScript(t, testFakeRelease, 0))
	if err != nil {
		t.Fatal(err)
	}
	refusing, err := arch.Archive(fakeInterpreterScript(t, "TEST-REFUSING v1", 1))
	if err != nil {
		t.Fatal(err)
	}
	source := []byte("module world/check_entry\n\nexport func main() -> string { \"x\" }\n")
	if err := EnsureSourceLoadable(context.Background(), arch, accepting, source); err != nil {
		t.Fatalf("accepting interpreter refused a valid source: %v", err)
	}
	err = EnsureSourceLoadable(context.Background(), arch, refusing, source)
	var invalid *TransitionSourceInvalidError
	if !errors.As(err, &invalid) {
		t.Fatalf("refusing interpreter error = %v, want *TransitionSourceInvalidError", err)
	}
	if invalid.Output == "" || invalid.Ref.IsZero() {
		t.Fatalf("invalid error = %+v, want the interpreter's captured output and the source ref", invalid)
	}
}
