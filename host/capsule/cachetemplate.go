package capsule

// Row 153 M3: every run starts from a copy of the (interpreter, source)
// compile-cache template built at publication (host/archive/capsulecache.go).
// This file does file I/O and logging only, never exec (the subprocess-site
// gate maps every exec to a driver; the one launch stays in capsule.go).

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

// coldState is the per-Runner cold-run bookkeeping.
type coldState struct {
	log  io.Writer
	cold atomic.Int64
}

// ColdRuns counts the runs that started without a template and so compiled
// from scratch. A rig asserts it, because a cold run is the slow path row 153
// exists to remove.
func (r *Runner) ColdRuns() int64 { return r.cs.cold.Load() }

// seedCache prepares the run's private compile-cache directory under root and
// reports whether it was seeded from the template. A missing or unusable
// template is a cold run: it is counted and written to the operator log, never
// silent, and the run proceeds (the interpreter compiles and fills cacheDir).
func (r *Runner) seedCache(root string, interpreter hashref.HashRef, source []byte) (cacheDir string, cold bool, err error) {
	cacheDir = filepath.Join(root, ".capsule-cache")
	if err := os.MkdirAll(cacheDir, 0o700); err != nil {
		return "", false, fmt.Errorf("capsule: create cache directory: %w", err)
	}
	if r.archive.CapsuleTemplateReady(interpreter, source) {
		if cerr := archive.CopyCacheTree(r.archive.CapsuleTemplateDir(interpreter, source), cacheDir); cerr == nil {
			return cacheDir, false, nil
		}
		// An unreadable template: start from an empty cache, as if missing.
		_ = os.RemoveAll(cacheDir)
		if err := os.MkdirAll(cacheDir, 0o700); err != nil {
			return "", false, fmt.Errorf("capsule: create cache directory: %w", err)
		}
	}
	r.cs.cold.Add(1)
	fmt.Fprintf(r.cs.log, "capsule: cold compile %s/%s (template missing)\n", interpreter.Digest(), hashref.SumSHA256(source).Digest())
	return cacheDir, true, nil
}

// promoteCold turns a successful cold run's cache into the template, so the
// next run is warm. A failure is logged, never returned: the run succeeded.
func (r *Runner) promoteCold(cacheDir string, interpreter hashref.HashRef, source []byte) {
	if err := r.archive.PromoteCopy(cacheDir, interpreter, source); err != nil {
		fmt.Fprintf(r.cs.log, "capsule: promote template %s/%s: %v\n", interpreter.Digest(), hashref.SumSHA256(source).Digest(), err)
	}
}
