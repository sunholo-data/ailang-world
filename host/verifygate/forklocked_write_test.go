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
// production-of-tests caller. scanForkLockedWrites hardcodes this function's first
// argument (Args[0]) as the exempt open closure. Do not reorder the signature.
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
		return os.CreateTemp(dir, pattern)
	}, func(f *os.File) error {
		name = f.Name()
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

var bannedOSFuncs = map[string]bool{"WriteFile": true, "OpenFile": true, "Create": true, "CreateTemp": true, "NewFile": true, "CopyFS": true}
var bannedSyscallFuncs = map[string]bool{"Open": true, "Openat": true, "Creat": true}

// wrapperNames are the fork-locked entry points whose uses the floor counts
// (identity-matched against the package-scope objects, not by bare name).
var wrapperNames = []string{
	"writeFileForkLocked",
	"copyFileForkLocked",
	"createTempForkLocked",
	"createForkLocked",
	"forkLockedWrite",
}

const wrapperFile = "forklocked_write_test.go"

type scanReport struct {
	violations   []string
	filesSeen    int
	wrapperCalls int
}

// scanForkLockedWrites type-checks the package in dir (all *.go files, test and
// non-test) or, when files is non-nil, the given name->source map, and reports
// every use (call, function value, method value, initializer) of a banned
// os/syscall file-open function outside the two exemptions: the sole opener
// call returned by the open closure (Args[0]) of forkLockedDo, and the body of
// TestKernelRefusesExecOfWriterOpenFile; plus aliased/dot "os" imports and
// imports of io/ioutil and golang.org/x/sys/unix. Exemptions match by
// package-scope object identity, never by the name of the enclosing function.
func scanForkLockedWrites(dir string, files map[string][]byte) (scanReport, error) {
	var rep scanReport
	live := files == nil
	if live {
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
	info := &types.Info{Uses: map[*ast.Ident]types.Object{}, Defs: map[*ast.Ident]types.Object{}}
	conf := types.Config{Importer: importer.ForCompiler(fset, "source", nil)}
	pkg, err := conf.Check("verifygate", fset, parsed, info)
	if err != nil {
		return rep, fmt.Errorf("type-check: %w", err)
	}
	add := func(pos token.Pos, msg string) {
		p := fset.Position(pos)
		rep.violations = append(rep.violations, fmt.Sprintf("%s:%d: %s", p.Filename, p.Line, msg))
	}
	isFunc := func(obj types.Object, path string, set map[string]bool) bool {
		fo, ok := obj.(*types.Func)
		return ok && fo.Pkg() != nil && fo.Pkg().Path() == path && set[fo.Name()]
	}
	doObj := pkg.Scope().Lookup("forkLockedDo")
	kernelObj := pkg.Scope().Lookup("TestKernelRefusesExecOfWriterOpenFile")
	wrapperObjs := map[types.Object]bool{}
	for _, n := range wrapperNames {
		o := pkg.Scope().Lookup(n)
		if live {
			if o == nil {
				return rep, fmt.Errorf("wrapper %s is not declared in the package", n)
			}
			if filepath.Base(fset.Position(o.Pos()).Filename) != wrapperFile {
				return rep, fmt.Errorf("wrapper %s is not declared in %s", n, wrapperFile)
			}
		}
		if o != nil {
			wrapperObjs[o] = true
		}
	}

	type span struct {
		from, to token.Pos
		name     string
	}
	var decls []span
	var kernelBody *span
	exempt := map[*ast.Ident]bool{}
	for _, f := range parsed {
		for _, imp := range f.Imports {
			switch imp.Path.Value {
			case `"os"`:
				if imp.Name != nil {
					add(imp.Pos(), fmt.Sprintf("aliased or dot import of os (%s) evades the write scan", imp.Name.Name))
				}
			case `"io/ioutil"`, `"golang.org/x/sys/unix"`:
				add(imp.Pos(), fmt.Sprintf("banned import %s evades the write scan", imp.Path.Value))
			}
		}
		for _, d := range f.Decls {
			if fd, ok := d.(*ast.FuncDecl); ok {
				sp := span{fd.Pos(), fd.End(), fd.Name.Name}
				decls = append(decls, sp)
				if kernelObj != nil && info.Defs[fd.Name] == kernelObj {
					c := sp
					kernelBody = &c
				}
			}
		}
		// Exempt opener: the sole result call of the last statement of the
		// FuncLit that is Args[0] of a forkLockedDo call.
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || doObj == nil || len(call.Args) == 0 {
				return true
			}
			fid, ok := call.Fun.(*ast.Ident)
			if !ok || info.Uses[fid] != doObj {
				return true
			}
			lit, ok := call.Args[0].(*ast.FuncLit)
			if !ok || len(lit.Body.List) == 0 {
				return true
			}
			ret, ok := lit.Body.List[len(lit.Body.List)-1].(*ast.ReturnStmt)
			if !ok || len(ret.Results) != 1 {
				return true
			}
			opener, ok := ret.Results[0].(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := opener.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if fo, ok := info.Uses[sel.Sel].(*types.Func); ok && fo.Pkg() != nil && fo.Pkg().Path() == "os" &&
				(fo.Name() == "OpenFile" || fo.Name() == "CreateTemp") {
				exempt[sel.Sel] = true
			}
			return true
		})
	}
	for id, obj := range info.Uses {
		if wrapperObjs[obj] && filepath.Base(fset.Position(id.Pos()).Filename) != wrapperFile {
			rep.wrapperCalls++
		}
		if !isFunc(obj, "os", bannedOSFuncs) && !isFunc(obj, "syscall", bannedSyscallFuncs) {
			continue
		}
		if exempt[id] {
			continue
		}
		if kernelBody != nil && id.Pos() >= kernelBody.from && id.Pos() < kernelBody.to {
			continue
		}
		where := "package-level initializer"
		for _, d := range decls {
			if id.Pos() >= d.from && id.Pos() < d.to {
				where = "func " + d.name
			}
		}
		add(id.Pos(), fmt.Sprintf("raw %s.%s in %s: use a fork-locked wrapper", obj.Pkg().Name(), obj.Name(), where))
	}
	sort.Strings(rep.violations)
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

	// Known positives: one per banned function, plus evasion shapes. A want
	// entry is "file:line" or "file:line|text" (the reported line must contain text).
	const doDecl = "func forkLockedDo(open func() (*os.File, error), fill func(*os.File) error, hold func()) error { return nil }\n"
	positives := []struct {
		name, src string
		want      []string
		absent    []string // "file:line" prefixes that must NOT be reported
	}{
		{"WriteFile", "package p\n\nimport \"os\"\n\nfunc f() { _ = os.WriteFile(\"x\", nil, 0o644) }\n", []string{"fixture_WriteFile.go:5"}, nil},
		{"OpenFile", "package p\n\nimport \"os\"\n\nfunc f() { _, _ = os.OpenFile(\"x\", 0, 0) }\n", []string{"fixture_OpenFile.go:5"}, nil},
		{"Create", "package p\n\nimport \"os\"\n\nfunc f() { _, _ = os.Create(\"x\") }\n", []string{"fixture_Create.go:5"}, nil},
		{"CreateTemp", "package p\n\nimport \"os\"\n\nfunc f() { _, _ = os.CreateTemp(\"\", \"x\") }\n", []string{"fixture_CreateTemp.go:5"}, nil},
		{"AliasedImport", "package p\n\nimport osw \"os\"\n\nfunc f() { _ = osw.WriteFile(\"x\", nil, 0o644) }\n", []string{"fixture_AliasedImport.go:3", "fixture_AliasedImport.go:5"}, nil},
		{"FuncValue", "package p\n\nimport \"os\"\n\nfunc h() { wf := os.WriteFile; _ = wf(\"x\", nil, 0o755) }\n",
			[]string{"fixture_FuncValue.go:5|raw os.WriteFile in func h"}, nil},
		{"RawFd", "package p\n\nimport (\n\t\"os\"\n\t\"syscall\"\n)\n\nfunc h() {\n\tfd, _ := syscall.Open(\"x\", syscall.O_CREAT|syscall.O_WRONLY, 0o755)\n\tf := os.NewFile(uintptr(fd), \"x\")\n\tf.Close()\n}\n",
			[]string{"fixture_RawFd.go:9|raw syscall.Open in func h", "fixture_RawFd.go:10|raw os.NewFile in func h"}, nil},
		{"PkgInit", "package p\n\nimport \"os\"\n\nvar _ = os.WriteFile(\"x\", nil, 0o644)\n",
			[]string{"fixture_PkgInit.go:5|package-level initializer"}, nil},
		{"NamedWrapper", "package p\n\nimport \"os\"\n\nfunc forkLockedWrite() { _, _ = os.OpenFile(\"x\", 0, 0) }\n",
			[]string{"fixture_NamedWrapper.go:5|raw os.OpenFile in func forkLockedWrite"}, nil},
		{"Ioutil", "package p\n\nimport \"io/ioutil\"\n\nfunc h() { _ = ioutil.WriteFile(\"x\", nil, 0o644) }\n",
			[]string{"fixture_Ioutil.go:3|banned import"}, nil},
		{"RootWrite", "package p\n\nimport \"os\"\n\nfunc h(r *os.Root) { _ = r.WriteFile(\"x\", nil, 0o755) }\n",
			[]string{"fixture_RootWrite.go:5|raw os.WriteFile in func h"}, nil},
		// P-Smuggle: a raw write inside the open closure is reported; only the
		// returned opener stays exempt.
		{"Smuggle", "package p\n\nimport \"os\"\n\n" + doDecl + "\nfunc h() {\n\t_ = forkLockedDo(func() (*os.File, error) {\n\t\t_ = os.WriteFile(\"y\", nil, 0o644)\n\t\treturn os.OpenFile(\"x\", 0, 0)\n\t}, nil, nil)\n}\n",
			[]string{"fixture_Smuggle.go:9|raw os.WriteFile in func h"}, []string{"fixture_Smuggle.go:10"}},
		{"CopyFS", "package p\n\nimport \"os\"\n\nfunc h() { _ = os.CopyFS(\"dst\", os.DirFS(\"src\")) }\n",
			[]string{"fixture_CopyFS.go:5|raw os.CopyFS in func h"}, nil},
		// P-LocalDo: a local forkLockedDo is not the package-scope wrapper, so its
		// open closure gets no exemption.
		{"LocalDo", "package p\n\nimport \"os\"\n\n" + doDecl + "\nfunc h() {\n\tforkLockedDo := func(open func() (*os.File, error), fill func(*os.File) error, hold func()) error { return nil }\n\t_ = forkLockedDo(func() (*os.File, error) { return os.OpenFile(\"x\", 0, 0) }, nil, nil)\n}\n",
			[]string{"fixture_LocalDo.go:9|raw os.OpenFile in func h"}, nil},
		// P-CreateOpener: only os.OpenFile / os.CreateTemp are exempt openers.
		{"CreateOpener", "package p\n\nimport \"os\"\n\n" + doDecl + "\nfunc h() {\n\t_ = forkLockedDo(func() (*os.File, error) { return os.Create(\"x\") }, nil, nil)\n}\n",
			[]string{"fixture_CreateOpener.go:8|raw os.Create in func h"}, nil},
	}
	for _, p := range positives {
		fname := "fixture_" + p.name + ".go"
		rep, err := scanForkLockedWrites("", map[string][]byte{fname: []byte(p.src)})
		if err != nil {
			t.Fatalf("fixture %s: %v", p.name, err)
		}
		for _, want := range p.want {
			prefix, text, _ := strings.Cut(want, "|")
			found := false
			for _, v := range rep.violations {
				if strings.HasPrefix(v, prefix+":") && strings.Contains(v, text) {
					found = true
				}
			}
			if !found {
				t.Fatalf("fixture %s: %s not reported (got %v)", p.name, want, rep.violations)
			}
		}
		for _, absent := range p.absent {
			for _, v := range rep.violations {
				if strings.HasPrefix(v, absent+":") {
					t.Fatalf("fixture %s: exempt opener %s reported", p.name, v)
				}
			}
		}
	}

	// Known negatives: a wrapper-shaped call, os.ReadFile, and an exempt open closure are clean.
	negatives := map[string]string{
		"negative":    "package p\n\nimport \"os\"\n\nfunc writeFileForkLocked(p string) error { return nil }\n\nfunc f() {\n\t_ = writeFileForkLocked(\"x\")\n\t_, _ = os.ReadFile(\"x\")\n}\n",
		"OpenClosure": "package p\n\nimport \"os\"\n\n" + doDecl + "\nfunc h() {\n\t_ = forkLockedDo(func() (*os.File, error) { return os.OpenFile(\"x\", 0, 0) }, nil, nil)\n}\n",
	}
	for name, src := range negatives {
		rep, err := scanForkLockedWrites("", map[string][]byte{"fixture_" + name + ".go": []byte(src)})
		if err != nil {
			t.Fatalf("negative fixture %s: %v", name, err)
		}
		if len(rep.violations) != 0 {
			t.Fatalf("negative fixture %s was flagged: %v", name, rep.violations)
		}
	}
}

// ---- R2 tripwire: parallel tests in write+fork packages outside this scan ----

// loadGoDirs reads every *.go file under roots, grouped by directory. Keys are
// the directory paths with the leading "../" segments trimmed (host/broker).
func loadGoDirs(roots ...string) (map[string]map[string][]byte, error) {
	dirs := map[string]map[string][]byte{}
	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(d.Name(), ".go") {
				return nil
			}
			raw, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			dir := strings.TrimLeft(filepath.ToSlash(filepath.Dir(path)), "./")
			if dirs[dir] == nil {
				dirs[dir] = map[string][]byte{}
			}
			dirs[dir][d.Name()] = raw
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return dirs, nil
}

// parallelWriteForkViolations reports every t.Parallel() call in a _test.go of
// a directory (other than host/verifygate, which the full scan covers) whose
// .go files both write files and fork. It is syntactic: it matches the selector's
// X identifier name, so an aliased import evades it (R2').
func parallelWriteForkViolations(dirs map[string]map[string][]byte) (viol []string, nDirs, nWriteFork int, err error) {
	writeSel := map[string]map[string]bool{
		"os":      {"WriteFile": true, "OpenFile": true, "Create": true, "CreateTemp": true, "NewFile": true, "CopyFS": true},
		"syscall": {"Open": true, "Openat": true, "Creat": true},
	}
	forkSel := map[string]map[string]bool{
		"exec":    {"Command": true, "CommandContext": true},
		"os":      {"StartProcess": true},
		"syscall": {"ForkExec": true, "StartProcess": true, "Exec": true},
	}
	names := make([]string, 0, len(dirs))
	for d := range dirs {
		names = append(names, d)
	}
	sort.Strings(names)
	for _, dir := range names {
		nDirs++
		fset := token.NewFileSet()
		var writes, forks bool
		var parallel []string
		files := make([]string, 0, len(dirs[dir]))
		for n := range dirs[dir] {
			files = append(files, n)
		}
		sort.Strings(files)
		for _, n := range files {
			f, perr := parser.ParseFile(fset, n, dirs[dir][n], 0)
			if perr != nil {
				return nil, 0, 0, perr
			}
			isTest := strings.HasSuffix(n, "_test.go")
			ast.Inspect(f, func(node ast.Node) bool {
				switch x := node.(type) {
				case *ast.SelectorExpr:
					if id, ok := x.X.(*ast.Ident); ok {
						if writeSel[id.Name][x.Sel.Name] {
							writes = true
						}
						if forkSel[id.Name][x.Sel.Name] {
							forks = true
						}
					}
				case *ast.CallExpr:
					if sel, ok := x.Fun.(*ast.SelectorExpr); ok && isTest && sel.Sel.Name == "Parallel" && len(x.Args) == 0 {
						parallel = append(parallel, fmt.Sprintf("%s/%s:%d", dir, n, fset.Position(x.Pos()).Line))
					}
				}
				return true
			})
		}
		if writes && forks {
			nWriteFork++
			if dir != "host/verifygate" {
				for _, p := range parallel {
					viol = append(viol, p+": t.Parallel in a write+fork package outside the fork-locked scan")
				}
			}
		}
	}
	sort.Strings(viol)
	return viol, nDirs, nWriteFork, nil
}

func TestNoParallelWriteForkPackagesOutsideVerifygate(t *testing.T) {
	dirs, err := loadGoDirs("../../host", "../../cmd")
	if err != nil {
		t.Fatal(err)
	}
	viol, nDirs, nWriteFork, err := parallelWriteForkViolations(dirs)
	if err != nil {
		t.Fatal(err)
	}
	if len(viol) > 0 {
		t.Fatalf("parallel tests in write+fork packages outside the fork-locked scan:\n%s", strings.Join(viol, "\n"))
	}
	if nDirs < 26 {
		t.Fatalf("sweep saw %d Go dirs, want >= 26", nDirs)
	}
	if nWriteFork < 12 {
		t.Fatalf("sweep saw %d write+fork dirs, want >= 12", nWriteFork)
	}

	// Known positive: t.Parallel + a write + a fork must be named. Known
	// negative: the same without t.Parallel() must not.
	const body = "package zz\n\nimport (\n\t\"os\"\n\t\"os/exec\"\n\t\"testing\"\n)\n\nfunc TestX(t *testing.T) {\n%s\t_ = os.WriteFile(\"x\", nil, 0o755)\n\t_ = exec.Command(\"x\")\n}\n"
	pos := map[string]map[string][]byte{"zz": {"x_test.go": []byte(fmt.Sprintf(body, "\tt.Parallel()\n"))}}
	pv, _, _, err := parallelWriteForkViolations(pos)
	if err != nil {
		t.Fatal(err)
	}
	if len(pv) != 1 || !strings.HasPrefix(pv[0], "zz/x_test.go:10: t.Parallel in a write+fork package") {
		t.Fatalf("known positive not named: %v", pv)
	}
	neg := map[string]map[string][]byte{"zz": {"x_test.go": []byte(fmt.Sprintf(body, ""))}}
	nv, _, _, err := parallelWriteForkViolations(neg)
	if err != nil {
		t.Fatal(err)
	}
	if len(nv) != 0 {
		t.Fatalf("known negative flagged: %v", nv)
	}

	// Each fixture below is its own single-key map and its own call: the
	// function reports every callsite across all keys of one map.
	// os.CopyFS is a write (not host-prefixed, so only writeSel is exercised).
	const bodyCopyFS = "package zz\n\nimport (\n\t\"os\"\n\t\"os/exec\"\n\t\"testing\"\n)\n\nfunc TestX(t *testing.T) {\n\tt.Parallel()\n\t_ = os.CopyFS(\"dst\", os.DirFS(\"src\"))\n\t_ = exec.Command(\"x\")\n}\n"
	cv, _, _, err := parallelWriteForkViolations(map[string]map[string][]byte{"zzcopyfs": {"x_test.go": []byte(bodyCopyFS)}})
	if err != nil {
		t.Fatal(err)
	}
	if len(cv) != 1 || !strings.HasPrefix(cv[0], "zzcopyfs/x_test.go:10: t.Parallel in a write+fork package") {
		t.Fatalf("known CopyFS positive not named: %v", cv)
	}
	// The exemption is exactly host/verifygate: a sibling host package is named...
	hv, _, _, err := parallelWriteForkViolations(map[string]map[string][]byte{"host/zz": {"x_test.go": []byte(fmt.Sprintf(body, "\tt.Parallel()\n"))}})
	if err != nil {
		t.Fatal(err)
	}
	if len(hv) != 1 || !strings.HasPrefix(hv[0], "host/zz/x_test.go:10: t.Parallel in a write+fork package") {
		t.Fatalf("known host positive not named: %v", hv)
	}
	// ...and host/verifygate itself is not.
	ev, _, _, err := parallelWriteForkViolations(map[string]map[string][]byte{"host/verifygate": {"x_test.go": []byte(fmt.Sprintf(body, "\tt.Parallel()\n"))}})
	if err != nil {
		t.Fatal(err)
	}
	if len(ev) != 0 {
		t.Fatalf("exempt host/verifygate fixture flagged: %v", ev)
	}
}
