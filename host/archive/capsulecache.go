package archive

// The capsule compile-cache template (row 153 M3).
//
// Running a transition cold compiles it and the std modules it imports on
// every call, which is what made the plan phase miss a 2 s cap under load.
// The interpreter keeps its compiled output in the directory named by
// AILANG_CACHE_DIR; `check` at the capsule's own staging path populates it. So
// publication builds one template per (interpreter, source) pair, stored under
// the archive root, and every capsule run starts from a COPY of it (never the
// template itself: the interpreter rewrites manifest.json on a warm run, and a
// shared directory would let one run corrupt the next).
//
// This file is file I/O only (no exec.Command*): the subprocess-site gate in
// host/broker maps every exec site to a driver, and the one that builds a
// template is CheckSource in check.go.

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// CapsuleEntryPath is the path, relative to a run's root, at which a source is
// staged both by the template build (CheckSource) and by the capsule runner.
// The interpreter's cache key includes the module path, so check and run must
// stage the same one by construction.
const CapsuleEntryPath = "host/capsule/main.ail"

// CacheDirEnv names the interpreter's compile-cache directory variable.
const CacheDirEnv = "AILANG_CACHE_DIR"

const (
	capsuleCacheDir = "capsule-cache"
	tmpPrefix       = ".tmp-"
	manifestRel     = "compile/manifest.json"
)

// renameDir is the promotion seam: tests substitute a rename that fails with
// EXDEV across directories.
var renameDir = os.Rename

// TemplateBuildError reports that a template could not be built or promoted
// even though the interpreter accepted the source: host-side infrastructure
// (a full disk, a permission fault, an interpreter that does not honour
// AILANG_CACHE_DIR), not a verdict on the source. It is retryable.
type TemplateBuildError struct {
	Interpreter, Source hashref.HashRef
	Reason              string
	Err                 error
}

func (e *TemplateBuildError) Error() string {
	msg := fmt.Sprintf("capsule template for interpreter %s source %s: %s", e.Interpreter, e.Source, e.Reason)
	if e.Err != nil {
		msg += ": " + e.Err.Error()
	}
	return msg
}

func (e *TemplateBuildError) Unwrap() error { return e.Err }

// CapsuleCacheRoot is the directory holding every template and in-flight tmp.
func (a *Archive) CapsuleCacheRoot() string { return filepath.Join(a.root, capsuleCacheDir) }

// CapsuleTemplateDir is where the template for (interpreter, source) lives:
// <root>/capsule-cache/<interpreter digest>/<hex sha256 of the source>.
func (a *Archive) CapsuleTemplateDir(interpreter hashref.HashRef, source []byte) string {
	return filepath.Join(a.CapsuleCacheRoot(), interpreter.Digest(), hashref.SumSHA256(source).Digest())
}

// CapsuleTemplateReady reports whether the template exists and is complete: a
// promoted template always holds compile/manifest.json, written last by the
// interpreter and renamed into place whole.
func (a *Archive) CapsuleTemplateReady(interpreter hashref.HashRef, source []byte) bool {
	info, err := os.Stat(filepath.Join(a.CapsuleTemplateDir(interpreter, source), filepath.FromSlash(manifestRel)))
	return err == nil && info.Mode().IsRegular()
}

// NewCapsuleTemplateTmp makes a fresh private directory inside capsule-cache/,
// the same filesystem as every template, so a promotion is one rename.
func (a *Archive) NewCapsuleTemplateTmp() (string, error) {
	if err := os.MkdirAll(a.CapsuleCacheRoot(), 0o700); err != nil {
		return "", err
	}
	return os.MkdirTemp(a.CapsuleCacheRoot(), tmpPrefix+"*")
}

// PromoteTemplate renames a fully built tmp directory into its final place. A
// racer that got there first is success: the winner's template is equivalent
// (same interpreter, same source), so tmp is discarded and the winner kept.
// Any other failure removes tmp and is a *TemplateBuildError.
func (a *Archive) PromoteTemplate(tmp string, interpreter hashref.HashRef, source []byte) error {
	final := a.CapsuleTemplateDir(interpreter, source)
	fail := func(reason string, err error) error {
		_ = os.RemoveAll(tmp)
		return &TemplateBuildError{Interpreter: interpreter, Source: hashref.SumSHA256(source), Reason: reason, Err: err}
	}
	if err := os.MkdirAll(filepath.Dir(final), 0o700); err != nil {
		return fail("create template directory", err)
	}
	if err := renameDir(tmp, final); err != nil {
		if isRacer(err) {
			_ = os.RemoveAll(tmp)
			return nil
		}
		return fail("promote template", err)
	}
	return nil
}

// isRacer is a rename onto a template another writer already promoted.
func isRacer(err error) bool {
	return errors.Is(err, fs.ErrExist) || errors.Is(err, syscall.ENOTEMPTY) || errors.Is(err, syscall.EEXIST)
}

// PromoteCopy turns a finished cold run's cache directory into a template: it
// copies runCacheDir into a fresh tmp INSIDE capsule-cache/ (the run's root may
// be on another filesystem, where a rename would fail with EXDEV) and promotes
// that. A run directory without compile/manifest.json is not a template.
func (a *Archive) PromoteCopy(runCacheDir string, interpreter hashref.HashRef, source []byte) error {
	build := func(reason string, err error) error {
		return &TemplateBuildError{Interpreter: interpreter, Source: hashref.SumSHA256(source), Reason: reason, Err: err}
	}
	if info, err := os.Stat(filepath.Join(runCacheDir, filepath.FromSlash(manifestRel))); err != nil || !info.Mode().IsRegular() {
		return build("run left no compile/manifest.json", err)
	}
	tmp, err := a.NewCapsuleTemplateTmp()
	if err != nil {
		return build("create template tmp", err)
	}
	if err := CopyCacheTree(runCacheDir, tmp); err != nil {
		_ = os.RemoveAll(tmp)
		return build("copy run cache", err)
	}
	return a.PromoteTemplate(tmp, interpreter, source)
}

// CopyCacheTree deep-copies src into dst (which must exist). Only regular
// files and directories are copied: a symlink, device or socket in a cache is
// refused, so a template can never smuggle a path out of the run's root.
func CopyCacheTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		switch {
		case d.IsDir():
			if rel == "." {
				return nil
			}
			return os.MkdirAll(target, 0o700)
		case d.Type().IsRegular():
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			return os.WriteFile(target, data, 0o600)
		default:
			return fmt.Errorf("refusing to copy %s (%s): only regular files and directories", rel, d.Type())
		}
	})
}
