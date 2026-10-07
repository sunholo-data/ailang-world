package capsule

// Row 153 r2 (judge N5): only a SUCCESSFUL cold run is promoted to a template.
// The interpreter is a shell fake that writes compile/manifest.json under
// AILANG_CACHE_DIR (as the real one does) and then exits with a chosen code,
// so the only variable between the two arms is the run's exit status.

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func cacheWritingInterpreter(t *testing.T, exit int) archivedFixture {
	t.Helper()
	path := filepath.Join(t.TempDir(), "cache-writer-"+strconv.Itoa(exit))
	script := "#!/bin/sh\n" +
		"if [ \"$1\" = \"--version\" ]; then echo 'CACHE-WRITER v1'; exit 0; fi\n" +
		"mkdir -p \"$AILANG_CACHE_DIR/compile\" && printf '{}' > \"$AILANG_CACHE_DIR/compile/manifest.json\"\n" +
		"echo '\"ok\"'\n" +
		"exit " + strconv.Itoa(exit) + "\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return archiveExecutable(t, path)
}

func TestFailedColdRunIsNotPromoted(t *testing.T) {
	src := source(`export func main() -> string { "x" }`)
	for _, tc := range []struct {
		name     string
		exit     int
		wantTmpl bool
	}{
		{"success promotes (control: the fake does write a manifest)", 0, true},
		{"failure does not promote", 3, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fx := cacheWritingInterpreter(t, tc.exit)
			runner := New(fx.archive, Config{Log: &syncBuf{}})
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			_, err := runner.RunContext(ctx, Entry{Interpreter: fx.ref, Source: src})
			if (err != nil) != (tc.exit != 0) {
				t.Fatalf("RunContext err = %v for exit %d", err, tc.exit)
			}
			if runner.ColdRuns() != 1 {
				t.Fatalf("ColdRuns = %d, want 1", runner.ColdRuns())
			}
			if got := fx.archive.CapsuleTemplateReady(fx.ref, src); got != tc.wantTmpl {
				t.Fatalf("template ready = %v after exit %d, want %v", got, tc.exit, tc.wantTmpl)
			}
		})
	}
}
