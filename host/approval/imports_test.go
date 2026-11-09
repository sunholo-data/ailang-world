package approval

import (
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"
	"testing"
)

// approvalForbidden are the packages host/approval may never import (coding
// standard S2: no I/O, no randomness).
var approvalForbidden = map[string]bool{
	"os": true, "os/exec": true, "os/signal": true, "net": true, "net/http": true,
	"syscall": true, "io/ioutil": true, "crypto/rand": true, "math/rand": true,
	"math/rand/v2": true, "unsafe": true,
}

// approvalImportViolation reports whether an import path is not standard
// library (a dot in its first segment) or is forbidden.
func approvalImportViolation(p string) bool {
	first := strings.SplitN(p, "/", 2)[0]
	return strings.Contains(first, ".") || approvalForbidden[p]
}

func TestApprovalImportsOnlyTheStandardLibrary(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	files, imports := 0, 0
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		files++
		f, err := parser.ParseFile(fset, name, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatal(err)
		}
		for _, im := range f.Imports {
			imports++
			p, err := strconv.Unquote(im.Path.Value)
			if err != nil {
				t.Fatal(err)
			}
			if approvalImportViolation(p) {
				t.Errorf("%s imports %q, which is not allowed in host/approval", name, p)
			}
		}
	}
	if files < 1 || imports < 5 {
		t.Fatalf("guard saw %d production files and %d imports; it must see at least 1 and 5", files, imports)
	}

	// Positive control: the same check must catch a non-stdlib import.
	src := "package x; import _ \"github.com/sunholo-data/ailang-world/host/store\""
	f, err := parser.ParseFile(fset, "x.go", src, parser.ImportsOnly)
	if err != nil {
		t.Fatal(err)
	}
	violations := 0
	for _, im := range f.Imports {
		p, err := strconv.Unquote(im.Path.Value)
		if err != nil {
			t.Fatal(err)
		}
		if approvalImportViolation(p) {
			violations++
		}
	}
	if violations != 1 {
		t.Fatalf("positive control found %d violations, want 1", violations)
	}
}
