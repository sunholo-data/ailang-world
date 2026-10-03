package verifygate

import (
	"errors"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

// Every test-side file write in this package goes through one of the wrappers
// below so that no fork in the test process can inherit a write fd (golang/go
// #22315: a forked child holds a copy of every open fd until its own execve,
// so a later execve of a just-written script fails with ETXTBSY on Linux).
// Every fork takes syscall.ForkLock for WRITING, so a writer holding RLock
// from open to close excludes all forks. See
// design_docs/planned/w-verifygate-etxtbsy.md.

// forkLockedProbe is a test-only seam. When set it runs on the goroutine that
// holds ForkLock.RLock, at stage "filled" (after fill/hold, before Close) and
// stage "closed" (after Close returned, RLock still held). It MUST NOT fork:
// a fork calls ForkLock.Lock and would deadlock against the caller's own RLock
// (R1). Only serial tests may set it.
var forkLockedProbe atomic.Pointer[func(stage string, f *os.File)]

// forkLockedDo holds syscall.ForkLock.RLock across open, fill, hold and Close.
// NOTHING inside the region may fork (exec.Command.Start, os.StartProcess, ...):
// a fork calls ForkLock.Lock and would deadlock against this goroutine's own
// RLock. fill must write bytes only. hold is a test-only seam, nil in every
// production-of-tests caller.
func forkLockedDo(open func() (*os.File, error), fill func(*os.File) error, hold func()) error {
	syscall.ForkLock.RLock()
	defer syscall.ForkLock.RUnlock()
	f, err := open()
	if err != nil {
		return err
	}
	if fill != nil {
		if err := fill(f); err != nil {
			f.Close()
			return err
		}
	}
	if hold != nil {
		hold()
	}
	if p := forkLockedProbe.Load(); p != nil {
		(*p)("filled", f)
	}
	err = f.Close()
	if p := forkLockedProbe.Load(); p != nil {
		(*p)("closed", f)
	}
	return err
}

func forkLockedWrite(path string, flag int, mode os.FileMode, fill func(*os.File) error, hold func()) error {
	return forkLockedDo(func() (*os.File, error) { return os.OpenFile(path, flag, mode) }, fill, hold)
}

// writeFileForkLocked has os.WriteFile semantics (it does not chmod an existing file).
func writeFileForkLocked(path string, data []byte, mode os.FileMode) error {
	return forkLockedWrite(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode, func(f *os.File) error {
		_, err := f.Write(data)
		return err
	}, nil)
}

// copyFileForkLocked copies src into the fresh path dst (O_EXCL). Error texts
// match the original copyGateFileErr.
func copyFileForkLocked(src *os.File, dst, rel string, mode os.FileMode) error {
	filled := false
	var copyErr error
	err := forkLockedWrite(dst, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode, func(f *os.File) error {
		filled = true
		_, copyErr = io.Copy(f, src)
		return copyErr
	}, nil)
	if err == nil {
		return nil
	}
	if !filled {
		return fmt.Errorf("create copy target %s: %v", rel, err)
	}
	if copyErr != nil {
		return fmt.Errorf("copy %s: %v", rel, copyErr)
	}
	return err
}

// createTempForkLocked replaces CreateTemp+WriteString+Close. On a fill or
// close error the file is removed and "" is returned.
func createTempForkLocked(dir, pattern string, data []byte) (string, error) {
	var name string
	err := forkLockedDo(func() (*os.File, error) {
		f, err := os.CreateTemp(dir, pattern)
		if f != nil {
			name = f.Name()
		}
		return f, err
	}, func(f *os.File) error {
		_, err := f.Write(data)
		return err
	}, nil)
	if err != nil {
		if name != "" {
			os.Remove(name)
		}
		return "", err
	}
	return name, nil
}

// createForkLocked is the os.Create shape: create/truncate, empty, closed.
func createForkLocked(path string) error {
	return forkLockedWrite(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o666, nil, nil)
}

func TestForkLockedWriteBlocksConcurrentFork(t *testing.T) {
	path := filepath.Join(t.TempDir(), "locked.txt")
	inside := make(chan struct{})
	release := make(chan struct{})
	writeDone := make(chan error, 1)
	go func() {
		writeDone <- forkLockedWrite(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644, nil, func() {
			close(inside)
			<-release
		})
	}()
	<-inside
	started := make(chan error, 1)
	cmd := exec.Command("/bin/sh", "-c", "exit 0")
	go func() { started <- cmd.Start() }()
	select {
	case <-started:
		close(release)
		t.Fatal("Start returned while a fork-locked write fd was open")
	case <-time.After(250 * time.Millisecond):
	}
	close(release)
	select {
	case err := <-started:
		if err != nil {
			t.Fatalf("Start failed after release: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Start did not return within 5 s of release")
	}
	if err := cmd.Wait(); err != nil {
		t.Fatal(err)
	}
	if err := <-writeDone; err != nil {
		t.Fatal(err)
	}
}

// TestForkLockedWrappersHoldLockThroughClose drives every wrapper with the
// after-Close probe: ForkLock must be read-held at "filled" and still at
// "closed" (after f.Close returned), and the stages must be exactly
// [filled closed] (a wrapper that bypasses forkLockedDo records none).
// Serial on purpose: it installs the package-global probe.
func TestForkLockedWrappersHoldLockThroughClose(t *testing.T) {
	root := t.TempDir()
	var stages []string
	var problems []string
	var wantLen int
	probe := func(stage string, f *os.File) {
		if !strings.HasPrefix(f.Name(), root) {
			return
		}
		stages = append(stages, stage)
		held := !syscall.ForkLock.TryLock()
		if !held {
			syscall.ForkLock.Unlock()
		}
		_, statErr := f.Stat()
		switch stage {
		case "filled":
			if !held {
				problems = append(problems, "ForkLock not read-held during fill")
			}
			if fi, err := f.Stat(); err != nil || fi.Size() != int64(wantLen) {
				problems = append(problems, fmt.Sprintf("file size at filled stage != %d (err %v)", wantLen, err))
			}
		case "closed":
			if !held {
				problems = append(problems, "ForkLock released before Close returned")
			}
			if !errors.Is(statErr, os.ErrClosed) {
				problems = append(problems, "probe at closed stage saw an open file")
			}
		}
	}
	data := []byte("#!/bin/sh\nexit 0\n")

	// src for copyFileForkLocked is written before the probe is installed.
	srcDir := t.TempDir()
	srcPath := filepath.Join(srcDir, "src.sh")
	if err := writeFileForkLocked(srcPath, data, 0o755); err != nil {
		t.Fatal(err)
	}
	src, err := os.Open(srcPath)
	if err != nil {
		t.Fatal(err)
	}
	defer src.Close()

	forkLockedProbe.Store(&probe)
	t.Cleanup(func() { forkLockedProbe.Store(nil) })

	type perm int
	const (
		exec755 perm = iota
		priv600
		noExec
	)
	cases := []struct {
		name    string
		mode    perm
		size    int
		content []byte
		run     func(dir string) (string, error)
	}{
		{"writeFileForkLocked", exec755, len(data), data, func(dir string) (string, error) {
			p := filepath.Join(dir, "w.sh")
			return p, writeFileForkLocked(p, data, 0o755)
		}},
		{"copyFileForkLocked", exec755, len(data), data, func(dir string) (string, error) {
			p := filepath.Join(dir, "c.sh")
			if _, err := src.Seek(0, io.SeekStart); err != nil {
				return p, err
			}
			return p, copyFileForkLocked(src, p, "rel", 0o755)
		}},
		{"createTempForkLocked", priv600, len(data), data, func(dir string) (string, error) {
			return createTempForkLocked(dir, "p*", data)
		}},
		{"createForkLocked", noExec, 0, nil, func(dir string) (string, error) {
			p := filepath.Join(dir, "e.txt")
			return p, createForkLocked(p)
		}},
		{"forkLockedWrite", exec755, len(data), data, func(dir string) (string, error) {
			p := filepath.Join(dir, "f.sh")
			return p, forkLockedWrite(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755, func(f *os.File) error {
				_, err := f.Write(data)
				return err
			}, nil)
		}},
	}
	for _, c := range cases {
		stages, problems, wantLen = nil, nil, c.size
		dir := filepath.Join(root, c.name)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		path, err := c.run(dir)
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}
		if got := strings.Join(stages, " "); got != "filled closed" {
			t.Errorf("%s: fork-lock probe stages = %v, want [filled closed]", c.name, stages)
		}
		for _, p := range problems {
			t.Errorf("%s: %s", c.name, p)
		}
		got, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read back: %v", c.name, err)
			continue
		}
		if string(got) != string(c.content) {
			t.Errorf("%s: content = %q, want %q", c.name, got, c.content)
		}
		fi, err := os.Stat(path)
		if err != nil {
			t.Errorf("%s: stat: %v", c.name, err)
			continue
		}
		pm := fi.Mode().Perm()
		switch c.mode {
		case exec755:
			if pm&0o111 != 0o111 {
				t.Errorf("%s: mode %v lacks exec bits", c.name, pm)
			}
		case priv600:
			if pm != 0o600 {
				t.Errorf("%s: mode %v, want 0600", c.name, pm)
			}
		case noExec:
			if pm&0o111 != 0 {
				t.Errorf("%s: mode %v has exec bits", c.name, pm)
			}
		}
	}
}

func TestKernelRefusesExecOfWriterOpenFile(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("ETXTBSY is a Linux kernel rule; darwin permits exec of a writer-open file (design V7)")
	}
	path := filepath.Join(t.TempDir(), "writer-open.sh")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("#!/bin/sh\nexit 0\n"); err != nil {
		f.Close()
		t.Fatal(err)
	}
	err = exec.Command(path).Start()
	if !errors.Is(err, syscall.ETXTBSY) {
		f.Close()
		t.Fatalf("expected ETXTBSY exec of a writer-open file, got %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	if err := exec.Command(path).Run(); err != nil {
		t.Fatalf("exec after Close failed: %v", err)
	}
}

// ---- mechanical guard: every file write in this package is fork-locked ----

var bannedWriteFuncs = map[string]bool{"WriteFile": true, "OpenFile": true, "Create": true, "CreateTemp": true}

// wrapperFuncs are the fork-locked entry points whose call sites the floor counts.
var wrapperFuncs = map[string]bool{
	"writeFileForkLocked":  true,
	"copyFileForkLocked":   true,
	"createTempForkLocked": true,
	"createForkLocked":     true,
	"forkLockedWrite":      true,
}

// rawWriteAllowlist is the set of enclosing FuncDecls that may call the banned
// os functions: the wrapper cores and the Linux kernel control (calls inside
// func literals attribute to the enclosing FuncDecl).
var rawWriteAllowlist = map[string]bool{
	"forkLockedWrite":                       true,
	"createTempForkLocked":                  true,
	"TestKernelRefusesExecOfWriterOpenFile": true,
}

const wrapperFile = "forklocked_write_test.go"

type scanReport struct {
	violations   []string
	filesSeen    int
	wrapperCalls int
}

// scanForkLockedWrites type-checks the package in dir (all *.go files, test and
// non-test) or, when files is non-nil, the given name->source map, and reports
// every banned os write call outside the allowlist plus aliased/dot "os" imports.
func scanForkLockedWrites(dir string, files map[string][]byte) (scanReport, error) {
	var rep scanReport
	if files == nil {
		files = map[string][]byte{}
		entries, err := os.ReadDir(dir)
		if err != nil {
			return rep, err
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
				continue
			}
			raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				return rep, err
			}
			files[e.Name()] = raw
		}
	}
	names := make([]string, 0, len(files))
	for n := range files {
		names = append(names, n)
	}
	sort.Strings(names)
	fset := token.NewFileSet()
	var parsed []*ast.File
	for _, n := range names {
		f, err := parser.ParseFile(fset, n, files[n], 0)
		if err != nil {
			return rep, err
		}
		parsed = append(parsed, f)
	}
	rep.filesSeen = len(parsed)
	info := &types.Info{Uses: map[*ast.Ident]types.Object{}}
	conf := types.Config{Importer: importer.ForCompiler(fset, "source", nil)}
	if _, err := conf.Check("verifygate", fset, parsed, info); err != nil {
		return rep, fmt.Errorf("type-check: %w", err)
	}
	add := func(pos token.Pos, msg string) {
		p := fset.Position(pos)
		rep.violations = append(rep.violations, fmt.Sprintf("%s:%d: %s", p.Filename, p.Line, msg))
	}
	for _, f := range parsed {
		fname := fset.Position(f.Pos()).Filename
		for _, imp := range f.Imports {
			if imp.Path.Value == `"os"` && imp.Name != nil {
				add(imp.Pos(), fmt.Sprintf("aliased or dot import of os (%s) evades the write scan", imp.Name.Name))
			}
		}
		visit := func(encl string, root ast.Node) {
			ast.Inspect(root, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				var id *ast.Ident
				switch fn := call.Fun.(type) {
				case *ast.Ident:
					id = fn
				case *ast.SelectorExpr:
					id = fn.Sel
				}
				if id == nil {
					return true
				}
				fo, ok := info.Uses[id].(*types.Func)
				if !ok || fo.Pkg() == nil {
					return true
				}
				if fo.Pkg().Path() == "os" && bannedWriteFuncs[fo.Name()] && !(encl != "" && rawWriteAllowlist[encl]) {
					where := "package-level initializer"
					if encl != "" {
						where = "func " + encl
					}
					add(call.Pos(), fmt.Sprintf("raw os.%s in %s: use a fork-locked wrapper", fo.Name(), where))
				}
				if wrapperFuncs[fo.Name()] && fname != wrapperFile {
					rep.wrapperCalls++
				}
				return true
			})
		}
		for _, d := range f.Decls {
			if fd, ok := d.(*ast.FuncDecl); ok {
				if fd.Body != nil {
					visit(fd.Name.Name, fd.Body)
				}
			} else {
				visit("", d)
			}
		}
	}
	return rep, nil
}

func TestVerifygateTestWritesAreForkLocked(t *testing.T) {
	live, err := scanForkLockedWrites(".", nil)
	if err != nil {
		t.Fatalf("live scan: %v", err)
	}
	// Violations first, so a raw write is named at its file:line rather than masked by the floors.
	if len(live.violations) > 0 {
		t.Fatalf("raw (non-fork-locked) file writes in verifygate tests:\n%s", strings.Join(live.violations, "\n"))
	}
	if live.filesSeen < 9 {
		t.Fatalf("scan enumerated %d files, want >= 9", live.filesSeen)
	}
	if live.wrapperCalls < 21 {
		t.Fatalf("scan found %d wrapper call sites, want >= 21", live.wrapperCalls)
	}

	// Known positives: one per banned function, plus an aliased import.
	positives := []struct {
		name, src string
		want      []string
	}{
		{"WriteFile", "package p\n\nimport \"os\"\n\nfunc f() { _ = os.WriteFile(\"x\", nil, 0o644) }\n", []string{"fixture_WriteFile.go:5"}},
		{"OpenFile", "package p\n\nimport \"os\"\n\nfunc f() { _, _ = os.OpenFile(\"x\", 0, 0) }\n", []string{"fixture_OpenFile.go:5"}},
		{"Create", "package p\n\nimport \"os\"\n\nfunc f() { _, _ = os.Create(\"x\") }\n", []string{"fixture_Create.go:5"}},
		{"CreateTemp", "package p\n\nimport \"os\"\n\nfunc f() { _, _ = os.CreateTemp(\"\", \"x\") }\n", []string{"fixture_CreateTemp.go:5"}},
		{"AliasedImport", "package p\n\nimport osw \"os\"\n\nfunc f() { _ = osw.WriteFile(\"x\", nil, 0o644) }\n", []string{"fixture_AliasedImport.go:3", "fixture_AliasedImport.go:5"}},
	}
	for _, p := range positives {
		fname := "fixture_" + p.name + ".go"
		rep, err := scanForkLockedWrites("", map[string][]byte{fname: []byte(p.src)})
		if err != nil {
			t.Fatalf("fixture %s: %v", p.name, err)
		}
		for _, want := range p.want {
			found := false
			for _, v := range rep.violations {
				if strings.HasPrefix(v, want+":") {
					found = true
				}
			}
			if !found {
				t.Fatalf("fixture %s: %s not reported (got %v)", p.name, want, rep.violations)
			}
		}
	}

	// Known negative: a wrapper-shaped call and os.ReadFile are clean.
	neg := "package p\n\nimport \"os\"\n\nfunc writeFileForkLocked(p string) error { return nil }\n\nfunc f() {\n\t_ = writeFileForkLocked(\"x\")\n\t_, _ = os.ReadFile(\"x\")\n}\n"
	rep, err := scanForkLockedWrites("", map[string][]byte{"fixture_negative.go": []byte(neg)})
	if err != nil {
		t.Fatalf("negative fixture: %v", err)
	}
	if len(rep.violations) != 0 {
		t.Fatalf("negative fixture was flagged: %v", rep.violations)
	}
}
