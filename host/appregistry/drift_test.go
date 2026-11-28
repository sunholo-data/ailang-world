package appregistry

import (
	"path/filepath"
	"testing"

	"github.com/sunholo-data/ailang-world/host/sketchtest"
)

// TestAppRegistrySketchDrift parses design_docs/sketches/appregistry.ail and replays every
// inline row through the Go mirror of the same name (laws.go).
func TestAppRegistrySketchDrift(t *testing.T) {
	laws := sketchtest.Load(t, filepath.Join("..", "..", "design_docs", "sketches", "appregistry.ail"), "sketches/appregistry")
	n := sketchtest.Replay(t, laws, map[string]sketchtest.Mirror{
		"revisionParentOk": func(tb testing.TB, a []any) any {
			return revisionParentOk(sketchtest.Int(tb, a, 0), sketchtest.Bool(tb, a, 1))
		},
		"nextRevisionOk": func(tb testing.TB, a []any) any {
			return nextRevisionOk(sketchtest.Int(tb, a, 0), sketchtest.Int(tb, a, 1))
		},
		"nextVersion":    func(tb testing.TB, a []any) any { return nextVersion(sketchtest.Int(tb, a, 0)) },
		"versionOk":      func(tb testing.TB, a []any) any { return versionOk(sketchtest.Int(tb, a, 0)) },
		"publishChanges": func(tb testing.TB, a []any) any { return publishChanges(sketchtest.Bool(tb, a, 0)) },
		"retireOk":       func(tb testing.TB, a []any) any { return retireOk(sketchtest.Bool(tb, a, 0)) },
		"appCountOk":     func(tb testing.TB, a []any) any { return appCountOk(int(sketchtest.Int(tb, a, 0))) },
		"transitionCountOk": func(tb testing.TB, a []any) any {
			return transitionCountOk(int(sketchtest.Int(tb, a, 0)))
		},
		"manifestPinned": func(tb testing.TB, a []any) any { return manifestPinned(int(sketchtest.Int(tb, a, 0))) },
		"listingState": func(tb testing.TB, a []any) any {
			return listingState(sketchtest.Bool(tb, a, 0), sketchtest.Bool(tb, a, 1), sketchtest.Bool(tb, a, 2))
		},
		"launchable": func(tb testing.TB, a []any) any { return Launchable(sketchtest.Str(tb, a, 0)) },
	}, 42)
	t.Logf("evaluated %d sketch rows", n)
}
