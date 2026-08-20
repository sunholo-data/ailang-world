package capsule

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

func pinnedBinary(t *testing.T) string {
	t.Helper()
	bin := os.Getenv("AILANG_BIN")
	if bin == "" {
		t.Skip("AILANG_BIN not set; capsule requires the pinned released ailang binary")
	}
	info, err := os.Stat(bin)
	if err != nil || !info.Mode().IsRegular() {
		t.Skipf("AILANG_BIN %q is not a usable executable: %v", bin, err)
	}
	return bin
}

type archivedFixture struct {
	archive *archive.Archive
	ref     hashref.HashRef
	path    string
}

func archiveExecutable(t *testing.T, path string) archivedFixture {
	t.Helper()
	a := archive.New(filepath.Join(t.TempDir(), "world.db"))
	ref, err := a.Archive(path)
	if err != nil {
		t.Fatalf("archive interpreter: %v", err)
	}
	resolved, err := a.Resolve(ref)
	if err != nil {
		t.Fatalf("resolve interpreter: %v", err)
	}
	return archivedFixture{archive: a, ref: ref, path: resolved}
}

func archivePinned(t *testing.T) archivedFixture {
	t.Helper()
	return archiveExecutable(t, pinnedBinary(t))
}

func source(body string) []byte {
	return []byte("module host/capsule/main\n\n" + body + "\n")
}

func runControl(t *testing.T, execPath string, src []byte, caps string, sandbox bool) ([]byte, []byte, error) {
	t.Helper()
	root := t.TempDir()
	srcPath := filepath.Join(root, filepath.FromSlash(entryModulePath))
	if err := os.MkdirAll(filepath.Dir(srcPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(srcPath, src, 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(execPath,
		"run", "--quiet", "--caps", caps, "--entry", entryFn, entryModulePath)
	cmd.Dir = root
	if sandbox {
		cmd.Env = []string{"AILANG_FS_SANDBOX=" + root}
	} else {
		cmd.Env = []string{}
	}
	var stdout, stderr strings.Builder
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	return []byte(stdout.String()), []byte(stderr.String()), err
}

func TestF1PinnedInterpreterHashMismatchRefusedBeforeExec(t *testing.T) {
	fixture := archivePinned(t)
	info, err := os.Stat(fixture.path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(fixture.path, 0o755); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(fixture.path)
	if err != nil {
		t.Fatal(err)
	}
	data[len(data)-1] ^= 1
	if err := os.WriteFile(fixture.path, data, info.Mode()); err != nil {
		t.Fatal(err)
	}

	_, err = New(fixture.archive, Config{}).Run(Entry{
		Interpreter: fixture.ref,
		Source:      source(`export func main() -> string { "must-not-execute" }`),
	})
	var mismatch *HashMismatchError
	if !errors.As(err, &mismatch) {
		t.Fatalf("error = %T %v, want *HashMismatchError", err, err)
	}
}

func TestF2DefaultDenyCapabilitiesWithIOControl(t *testing.T) {
	fixture := archivePinned(t)
	src := source(`import std/io (println)

export func main() -> unit ! {IO} {
  println("io-control")
}`)
	_, err := New(fixture.archive, Config{}).Run(Entry{
		Interpreter: fixture.ref,
		Source:      src,
	})
	var execErr *ExecError
	if !errors.As(err, &execErr) {
		t.Fatalf("denied run error = %T %v, want *ExecError", err, err)
	}
	if !strings.Contains(string(execErr.Stderr), "effect 'IO' requires capability") {
		t.Fatalf("denied stderr = %q, want IO capability denial", execErr.Stderr)
	}

	stdout, stderr, err := runControl(t, fixture.path, src, "IO", true)
	if err != nil {
		t.Fatalf("--caps IO control failed: %v (stderr %q)", err, stderr)
	}
	if string(stdout) != "io-control\n" {
		t.Fatalf("--caps IO stdout = %q", stdout)
	}
}

func TestF3FilesystemSandboxEscapeWithUnsandboxedControl(t *testing.T) {
	fixture := archivePinned(t)
	outsideDir := t.TempDir()
	outsidePath := filepath.Join(outsideDir, "outside.txt")
	if err := os.WriteFile(outsidePath, []byte("outside-content"), 0o644); err != nil {
		t.Fatal(err)
	}
	src := source(fmt.Sprintf(`import std/fs (readFile)

export func main() -> string ! {FS} {
  readFile(%q)
}`, outsidePath))

	// Isolate F3 from F2 while retaining the shipped staging and child-env path.
	// caps is deliberately unexported; every external Runner stays default-deny.
	runner := New(fixture.archive, Config{})
	runner.caps = "FS"
	result, jailedErr := runner.Run(Entry{Interpreter: fixture.ref, Source: src})
	var execErr *ExecError
	if !errors.As(jailedErr, &execErr) {
		t.Fatalf("sandboxed error = %T %v, want *ExecError", jailedErr, jailedErr)
	}
	if !strings.Contains(string(result.Stderr), "escapes sandbox") {
		t.Fatalf("sandboxed stderr = %q, want escape refusal", result.Stderr)
	}
	stdout, stderr, err := runControl(t, fixture.path, src, "FS", false)
	if err != nil {
		t.Fatalf("unsandboxed control failed: %v (stderr %q)", err, stderr)
	}
	if string(stdout) != "outside-content\n" {
		t.Fatalf("unsandboxed stdout = %q", stdout)
	}
}

func TestF4ParentEnvironmentMarkerInvisible(t *testing.T) {
	const marker = "CAPSULE_F4_PARENT_MARKER"
	t.Setenv(marker, "visible")
	dir := t.TempDir()
	probe := filepath.Join(dir, "ailang-env-probe")
	script := fmt.Sprintf(`#!/bin/sh
if [ "$1" = "--version" ]; then
  echo "AILANG v0.30.0 env probe"
  exit 0
fi
if [ -n "$%s" ]; then
  echo visible
else
  echo hidden
fi
`, marker)
	if err := os.WriteFile(probe, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	fixture := archiveExecutable(t, probe)
	result, err := New(fixture.archive, Config{}).Run(Entry{
		Interpreter: fixture.ref,
		Source:      source(`export func main() -> string { "unused" }`),
	})
	if err != nil {
		t.Fatalf("environment probe: %v", err)
	}
	if string(result.Stdout) != "hidden\n" {
		t.Fatalf("probe observed parent marker: stdout %q", result.Stdout)
	}
}

func TestF5WallClockTimeoutHasElapsedBound(t *testing.T) {
	fixture := archivePinned(t)
	src := source(`func fib(n: int) -> int {
  if n < 2 then n else fib(n - 1) + fib(n - 2)
}

export func main() -> int {
  fib(28)
}`)
	const limit = 40 * time.Millisecond
	start := time.Now()
	_, err := New(fixture.archive, Config{ExecTimeout: limit}).Run(Entry{
		Interpreter: fixture.ref,
		Source:      src,
	})
	elapsed := time.Since(start)
	var timeout *TimeoutError
	if !errors.As(err, &timeout) {
		t.Fatalf("error after %s = %T %v, want *TimeoutError", elapsed, err, err)
	}
	if elapsed > 2*time.Second {
		t.Fatalf("timeout returned after %s, want <= 2s", elapsed)
	}
	if elapsed < limit/2 {
		t.Fatalf("timeout returned too early after %s, injected limit %s", elapsed, limit)
	}
}

func TestF6OutputCapReturnsStructuredOverflow(t *testing.T) {
	fixture := archivePinned(t)
	src := source(`func repeat(n: int) -> string {
  if n == 0 then "" else "0123456789abcdef${repeat(n - 1)}"
}

export func main() -> string {
  repeat(32)
}`)
	const limit = int64(64)
	control, err := New(fixture.archive, Config{MaxOutputBytes: 1024}).Run(Entry{
		Interpreter: fixture.ref,
		Source:      src,
	})
	if err != nil {
		t.Fatalf("large-cap control failed: %v (stdout %q, stderr %q)",
			err, control.Stdout, control.Stderr)
	}
	if len(control.Stdout) != 513 {
		t.Fatalf("large-cap control produced %d bytes, want 513", len(control.Stdout))
	}
	result, err := New(fixture.archive, Config{MaxOutputBytes: limit}).Run(Entry{
		Interpreter: fixture.ref,
		Source:      src,
	})
	var overflow *OutputLimitError
	if !errors.As(err, &overflow) {
		t.Fatalf("error = %T %v, want *OutputLimitError (stdout %q, stderr %q)",
			err, err, result.Stdout, result.Stderr)
	}
	if int64(len(result.Stdout))+int64(len(result.Stderr)) > 2*limit {
		t.Fatalf("captured %d bytes under %d-byte cap",
			len(result.Stdout)+len(result.Stderr), limit)
	}
}

var errWaitedWithoutKill = errors.New("fake child waited without kill")

type fakeChild struct {
	mu          sync.Mutex
	killCount   int
	killed      bool
	waitCount   int
	waitEntered chan struct{}
	waitRelease chan struct{}
}

func (c *fakeChild) Kill() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.killCount++
	c.killed = true
	return nil
}

func (c *fakeChild) Wait() error {
	c.mu.Lock()
	c.waitCount++
	killed := c.killed
	c.mu.Unlock()
	if c.waitEntered != nil {
		close(c.waitEntered)
	}
	if c.waitRelease != nil {
		<-c.waitRelease
	}
	if !killed {
		return errWaitedWithoutKill
	}
	return nil
}

func (c *fakeChild) state() (killCount, waitCount int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.killCount, c.waitCount
}

func TestOutputCollectionOverflowKillsAndOutranksDeadline(t *testing.T) {
	const limit = int64(16)
	want := bytes.Repeat([]byte("x"), int(limit))
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	cancel()
	if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Fatalf("context error = %v, want DeadlineExceeded", ctx.Err())
	}
	child := &fakeChild{}
	result, runErr, err := collectOutput(ctx, bytes.NewReader(append(bytes.Clone(want), 'x')), bytes.NewReader(nil), limit, time.Second, child)
	if !bytes.Equal(result.Stdout, want) {
		t.Fatalf("stdout = %q, want %q", result.Stdout, want)
	}
	killCount, waitCount := child.state()
	if killCount != 1 || waitCount != 1 {
		t.Fatalf("killCount = %d, waitCount = %d; want 1, 1", killCount, waitCount)
	}
	if errors.Is(runErr, errWaitedWithoutKill) {
		t.Fatalf("Wait did not observe Kill: %v", runErr)
	}
	var overflow *OutputLimitError
	if !errors.As(err, &overflow) {
		t.Fatalf("error = %T %v, want *OutputLimitError", err, err)
	}
	var timeout *TimeoutError
	if errors.As(err, &timeout) {
		t.Fatalf("error = %T %v, do not want *TimeoutError", err, err)
	}
}

func TestOutputCollectionAtLimitDoesNotKill(t *testing.T) {
	const limit = int64(16)
	want := bytes.Repeat([]byte("x"), int(limit))
	child := &fakeChild{}
	result, runErr, err := collectOutput(context.Background(), bytes.NewReader(want), bytes.NewReader(nil), limit, time.Second, child)
	if err != nil {
		t.Fatalf("collection error = %v, want nil", err)
	}
	if !bytes.Equal(result.Stdout, want) {
		t.Fatalf("stdout = %q, want %q", result.Stdout, want)
	}
	killCount, waitCount := child.state()
	if killCount != 0 || waitCount != 1 {
		t.Fatalf("killCount = %d, waitCount = %d; want 0, 1", killCount, waitCount)
	}
	if !errors.Is(runErr, errWaitedWithoutKill) {
		t.Fatalf("Wait result = %v, want %v", runErr, errWaitedWithoutKill)
	}
}

func TestOutputCollectionTwoOverflowsKillOnce(t *testing.T) {
	const limit = int64(16)
	overLimit := bytes.Repeat([]byte("x"), int(limit+1))
	child := &fakeChild{}
	_, runErr, err := collectOutput(context.Background(), bytes.NewReader(overLimit), bytes.NewReader(overLimit), limit, time.Second, child)
	if runErr != nil {
		t.Fatalf("Wait result = %v, want nil", runErr)
	}
	var overflow *OutputLimitError
	if !errors.As(err, &overflow) {
		t.Fatalf("error = %T %v, want *OutputLimitError", err, err)
	}
	killCount, _ := child.state()
	if killCount != 1 {
		t.Fatalf("killCount = %d, want 1", killCount)
	}
}

type blockingReader struct {
	entered chan<- struct{}
	release <-chan struct{}
	once    sync.Once
}

func (r *blockingReader) Read([]byte) (int, error) {
	r.once.Do(func() { r.entered <- struct{}{} })
	<-r.release
	return 0, io.EOF
}

func TestOutputCollectionCallerReleaseUnblocksReadersAndWait(t *testing.T) {
	const watchdog = 10 * time.Second
	readerEntered := make(chan struct{}, 2)
	readerRelease := make(chan struct{})
	waitEntered := make(chan struct{})
	waitRelease := make(chan struct{})
	child := &fakeChild{waitEntered: waitEntered, waitRelease: waitRelease}
	type outcome struct {
		result Result
		runErr error
		err    error
	}
	done := make(chan outcome, 1)
	go func() {
		result, runErr, err := collectOutput(
			context.Background(),
			&blockingReader{entered: readerEntered, release: readerRelease},
			&blockingReader{entered: readerEntered, release: readerRelease},
			16, time.Second, child,
		)
		done <- outcome{result: result, runErr: runErr, err: err}
	}()
	for i := 0; i < 2; i++ {
		select {
		case <-readerEntered:
		case <-time.After(watchdog):
			t.Fatalf("reader-entry phase %d never happened", i+1)
		}
	}
	close(readerRelease)
	select {
	case <-waitEntered:
	case <-time.After(watchdog):
		t.Fatal("wait-entry phase never happened")
	}
	close(waitRelease)
	select {
	case got := <-done:
		if got.err != nil {
			t.Fatalf("collection error = %v, want nil", got.err)
		}
		if !errors.Is(got.runErr, errWaitedWithoutKill) {
			t.Fatalf("Wait result = %v, want %v", got.runErr, errWaitedWithoutKill)
		}
		if len(got.result.Stdout) != 0 || len(got.result.Stderr) != 0 {
			t.Fatalf("result = %#v, want empty output", got.result)
		}
	case <-time.After(watchdog):
		t.Fatal("helper-return phase never happened")
	}
}

// F6 at production shape. The shipped F6 fixture emits 513 bytes — below one OS
// pipe buffer — so its child always exits on its own and the case where nothing
// drains the pipe is never exercised. Here the untrusted transition returns more
// than 64 KiB, which needs no capability: the interpreter prints the entry's
// return value even under --caps "". If the capsule does not kill an overflowing
// child, cmd.Wait() blocks until the wall clock and the caller is handed a
// *TimeoutError for what is an overflow — F6 silently degrading into F5.
func TestF6OutputCapKillsChildBeyondOnePipeBuffer(t *testing.T) {
	fixture := archivePinned(t)
	src := source(`func dbl(s: string, n: int) -> string {
  if n == 0 then s else dbl("${s}${s}", n - 1)
}

export func main() -> string {
  dbl("0123456789abcdef", 13)
}`)
	const limit = int64(1024)
	const clock = 5 * time.Second
	start := time.Now()
	result, err := New(fixture.archive, Config{
		ExecTimeout: clock, MaxOutputBytes: limit,
	}).Run(Entry{Interpreter: fixture.ref, Source: src})
	elapsed := time.Since(start)

	var overflow *OutputLimitError
	if !errors.As(err, &overflow) {
		t.Fatalf("error after %s = %T %v, want *OutputLimitError", elapsed, err, err)
	}
	var timeout *TimeoutError
	if errors.As(err, &timeout) {
		t.Fatalf("overflow was reported as a wall-clock timeout after %s", elapsed)
	}
	// The whole point: the child is killed at once rather than run to the clock.
	if elapsed >= clock {
		t.Fatalf("overflow took %s, i.e. the full %s bound — the child was not killed",
			elapsed, clock)
	}
	if int64(len(result.Stdout)) > limit {
		t.Fatalf("captured %d bytes under a %d-byte cap", len(result.Stdout), limit)
	}
}
