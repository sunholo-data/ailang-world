package appmanifest

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/sketchtest"
)

// The drift test binds manifest.go to design_docs/sketches/appmanifest.ail: it PARSES the
// sketch (host/sketchtest) and replays every inline test row of every exported law through
// the Go mirror of the same name. A law without a mirror, a mirror without a law, a moved
// bound on either side, or a row the mirror disagrees with turns it red.

var mirrors = map[string]sketchtest.Mirror{
	"schemaOk":  func(t testing.TB, a []any) any { return schemaOk(sketchtest.Str(t, a, 0)) },
	"nameLenOk": func(t testing.TB, a []any) any { return nameLenOk(int(sketchtest.Int(t, a, 0))) },
	"nameByteOk": func(t testing.TB, a []any) any {
		return nameByteOk(int(sketchtest.Int(t, a, 0)), sketchtest.Bool(t, a, 1))
	},
	"nameReserved":   func(t testing.TB, a []any) any { return nameReserved(sketchtest.Str(t, a, 0)) },
	"titleLenOk":     func(t testing.TB, a []any) any { return titleLenOk(int(sketchtest.Int(t, a, 0))) },
	"summaryLenOk":   func(t testing.TB, a []any) any { return summaryLenOk(int(sketchtest.Int(t, a, 0))) },
	"agentOk":        func(t testing.TB, a []any) any { return agentOk(sketchtest.Str(t, a, 0)) },
	"modeOk":         func(t testing.TB, a []any) any { return modeOk(sketchtest.Str(t, a, 0)) },
	"deadlineOk":     func(t testing.TB, a []any) any { return deadlineOk(sketchtest.Int(t, a, 0)) },
	"ttlOk":          func(t testing.TB, a []any) any { return ttlOk(sketchtest.Int(t, a, 0)) },
	"effectAllowed":  func(t testing.TB, a []any) any { return effectAllowed(sketchtest.Str(t, a, 0)) },
	"scopeOk":        func(t testing.TB, a []any) any { return scopeOk(sketchtest.Str(t, a, 0)) },
	"budgetOk":       func(t testing.TB, a []any) any { return budgetOk(sketchtest.Int(t, a, 0)) },
	"grantNarrows":   func(t testing.TB, a []any) any { return grantNarrows(sketchtest.Int(t, a, 0), sketchtest.Int(t, a, 1)) },
	"grantCountOk":   func(t testing.TB, a []any) any { return grantCountOk(int(sketchtest.Int(t, a, 0))) },
	"listCountOk":    func(t testing.TB, a []any) any { return listCountOk(int(sketchtest.Int(t, a, 0))) },
	"checkNameLenOk": func(t testing.TB, a []any) any { return checkNameLenOk(int(sketchtest.Int(t, a, 0))) },
	"checkTimeoutOk": func(t testing.TB, a []any) any { return checkTimeoutOk(sketchtest.Int(t, a, 0)) },
	"outputKindOk":   func(t testing.TB, a []any) any { return outputKindOk(sketchtest.Str(t, a, 0)) },
	"noDotDotSegment": func(t testing.TB, a []any) any {
		return noDotDotSegment(sketchtest.Str(t, a, 0))
	},
	"outputGlobOk": func(t testing.TB, a []any) any { return outputGlobOk(sketchtest.Str(t, a, 0)) },
	"viewOk":       func(t testing.TB, a []any) any { return viewOk(sketchtest.Str(t, a, 0)) },
	"viewCountOk":  func(t testing.TB, a []any) any { return viewCountOk(int(sketchtest.Int(t, a, 0))) },
}

func sketchPath() string {
	return filepath.Join("..", "..", "design_docs", "sketches", "appmanifest.ail")
}

func TestAppManifestSketchDrift(t *testing.T) {
	laws := sketchtest.Load(t, sketchPath(), "sketches/appmanifest")
	n := sketchtest.Replay(t, laws, mirrors, 122)
	t.Logf("evaluated %d sketch rows over %d laws", n, len(laws))
}

// TestAllowedEffectsAreTheSketchList keeps effectAllowed's Go set equal to the names the
// sketch's ensures lists, so neither side can grow an effect alone. allowed_effects_test.go
// binds the same set to the broker's effect constants.
func TestAllowedEffectsAreTheSketchList(t *testing.T) {
	src, err := os.ReadFile(sketchPath())
	if err != nil {
		t.Fatal(err)
	}
	s := string(src)
	start := strings.Index(s, "export func effectAllowed")
	end := strings.Index(s[start:], "\ntests [")
	if start < 0 || end < 0 {
		t.Fatal("effectAllowed law not found in the sketch")
	}
	var fromSketch []string
	for _, m := range regexp.MustCompile(`e == "([A-Za-z.]+)"`).FindAllStringSubmatch(s[start:start+end], -1) {
		fromSketch = append(fromSketch, m[1])
	}
	var fromGo []string
	for e := range allowedEffects {
		fromGo = append(fromGo, e)
	}
	sort.Strings(fromSketch)
	sort.Strings(fromGo)
	if !reflect.DeepEqual(fromSketch, fromGo) {
		t.Fatalf("effectAllowed: sketch ensures lists %v, Go allows %v", fromSketch, fromGo)
	}
}
