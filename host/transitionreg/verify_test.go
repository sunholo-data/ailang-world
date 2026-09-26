package transitionreg

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/sunholo-data/ailang-world/host/archive"
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
