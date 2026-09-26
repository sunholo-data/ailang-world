package archive

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// checkingInterpreter writes the house-pattern shell-script fake whose
// `check <file>` echoes the STAGED file's bytes (so the test sees exactly what
// CheckSource staged, and where) and exits checkExit. body distinguishes
// fakes so their archived hashes differ.
func checkingInterpreter(t *testing.T, body string, checkExit int) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "check-fake-"+body)
	script := "#!/bin/sh\n" +
		"if [ \"$1\" = \"--version\" ]; then echo 'CHECK-FAKE v1'; exit 0; fi\n" +
		"if [ \"$1\" = \"check\" ]; then cat \"$2\"; echo \"staged-as:$2\"; exit " + strconv.Itoa(checkExit) + "; fi\n" +
		"# body marker: " + body + "\n" +
		"exit 0\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

// TestCheckSourceReportsTheArchivedInterpretersVerdict pins CheckSource's
// contract with fakes: the canonical bytes are staged as entry.ail in a
// scratch root and handed to the ARCHIVED interpreter's `check`; exit 0 →
// Passed with the captured output, non-zero → !Passed with the captured
// output (a refusal is a verdict, not an infrastructure error); an
// unarchived interpreter ref is an error.
func TestCheckSourceReportsTheArchivedInterpretersVerdict(t *testing.T) {
	a := New(storeDBPath(t))
	accepting, err := a.Archive(checkingInterpreter(t, "accepting", 0))
	if err != nil {
		t.Fatal(err)
	}
	refusing, err := a.Archive(checkingInterpreter(t, "refusing", 1))
	if err != nil {
		t.Fatal(err)
	}
	source := []byte("module check/entry\n\nexport func main() -> string { \"x\" }\n")
	ok, err := a.CheckSource(context.Background(), accepting, source)
	if err != nil || !ok.Passed {
		t.Fatalf("accepting check = (%+v, %v), want Passed", ok, err)
	}
	if !strings.Contains(ok.Output, "export func main()") || !strings.Contains(ok.Output, "staged-as:entry.ail") {
		t.Fatalf("accepting output = %q, want the staged source bytes checked as entry.ail", ok.Output)
	}
	bad, err := a.CheckSource(context.Background(), refusing, source)
	if err != nil {
		t.Fatalf("a refusing interpreter is a verdict, not an error: %v", err)
	}
	if bad.Passed || !strings.Contains(bad.Output, "staged-as:entry.ail") {
		t.Fatalf("refusing check = %+v, want !Passed with the interpreter's captured output", bad)
	}
	if _, err := a.CheckSource(context.Background(), hashref.SumSHA256([]byte("never archived")), source); err == nil {
		t.Fatal("an interpreter that is not archived must be an error, never a verdict")
	}
}

// TestCheckSourceUnderPinnedInterpreter is the real-interpreter arm (V26):
// the pinned released binary, ARCHIVED, accepts a standalone module-declared
// source and refuses an import-bearing one (LDR001, hermetic scope) and
// garbage. Gated on AILANG_BIN like the CLI arm — never skip.
func TestCheckSourceUnderPinnedInterpreter(t *testing.T) {
	bin := os.Getenv("AILANG_BIN")
	if bin == "" {
		t.Fatal("AILANG_BIN unset: the real-interpreter check arm needs the pinned released interpreter; never skip")
	}
	a := New(storeDBPath(t))
	ref, err := a.Archive(bin)
	if err != nil {
		t.Fatalf("archive pinned interpreter: %v", err)
	}
	cases := []struct {
		name   string
		source string
		pass   bool
		want   string
	}{
		{"standalone module", "module check/entry\n\nexport func echo(x: string) -> string { x }\n", true, "No errors found"},
		{"import-bearing module", "module check/entry\n\nimport world/types (World)\n\nexport func apply(w: World) -> World {\n  w\n}\n", false, "LDR001"},
		{"garbage", "this is not AILANG source at all\n", false, "rror"},
	}
	for _, c := range cases {
		res, err := a.CheckSource(context.Background(), ref, []byte(c.source))
		if err != nil {
			t.Fatalf("%s: infrastructure error %v", c.name, err)
		}
		if res.Passed != c.pass || !strings.Contains(res.Output, c.want) {
			t.Fatalf("%s: check = %+v, want Passed=%v with %q in the output", c.name, res, c.pass, c.want)
		}
	}
}
