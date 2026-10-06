//go:build darwin || linux

package main

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/ptytest"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// ---------------------------------------------------------------------------
// THE FENCE END TO END, in a real child process with the PRODUCTION probe.
//
// Every AC21 row injects a ttyProbe. These two tests do not: the child runs
// run() with probeControllingTerminal, and the only thing that differs between
// them is the kernel's answer to open("/dev/tty") — none at all (setsid with no
// TIOCSCTTY: an agent harness, cron, CI) versus a pty the test types into.
// In BOTH, stdin is not the controlling terminal, which is the IDE-pane shape
// that used to STOP every attended operator with
// `reason=stdin-is-not-the-controlling-terminal`.
// ---------------------------------------------------------------------------

// ptyChildArgs carries the child's argv as JSON. Its presence is what turns
// TestPtyChildProcess from a skip into the command.
const ptyChildArgs = "WORLD_PUBLISH_PTY_CHILD_ARGS"

// TestPtyChildProcess is not a test of its own: it is the command, run in the
// child the tests below start, with the production environment.
func TestPtyChildProcess(t *testing.T) {
	raw := os.Getenv(ptyChildArgs)
	if raw == "" {
		t.Skip("the child process of the pty tests; it runs only when they start it")
	}
	var args []string
	if err := json.Unmarshal([]byte(raw), &args); err != nil {
		t.Fatal(err)
	}
	os.Exit(run(args, os.Stdin, os.Stdout, os.Stderr, environment{getenv: os.Getenv, probe: probeControllingTerminal}))
}

// ptyChild builds the child command. The CI tripwire variables are removed
// from its environment, because the operator these tests stand in for is not
// in CI; the fence under test is the terminal, not the tripwire.
func ptyChild(t *testing.T, args []string) (*exec.Cmd, *ptytest.Buffer, *ptytest.Buffer) {
	t.Helper()
	encoded, err := json.Marshal(args)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestPtyChildProcess$", "-test.count=1")
	for _, kv := range os.Environ() {
		name, _, _ := strings.Cut(kv, "=")
		if name == "CI" || name == "GITHUB_ACTIONS" || name == ptyChildArgs {
			continue
		}
		cmd.Env = append(cmd.Env, kv)
	}
	cmd.Env = append(cmd.Env, ptyChildArgs+"="+string(encoded))
	var stdout, stderr ptytest.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	return cmd, &stdout, &stderr
}

// transitionsFixture is the in-process happy path's setup: a bootstrapped
// store, a one-descriptor manifest and the shell-script fake interpreter.
func transitionsFixture(t *testing.T) map[string]string {
	t.Helper()
	storePath := filepath.Join(t.TempDir(), "world.db")
	srcPath := filepath.Join(t.TempDir(), "echo.ail")
	if err := os.WriteFile(srcPath, []byte("module world/echo\n\nexport func apply(w: World) -> World ! {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	bootstrapEpochRegistry(t, storePath, "test-interpreter-version")
	return map[string]string{"store": storePath, "manifest": echoManifest(t, srcPath, false), "ailang-bin": fakeInterpreter(t)}
}

func registryHeadExists(t *testing.T, storePath string) bool {
	t.Helper()
	db, err := store.Open(storePath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	_, _, ok, err := transitionreg.NewReader(db).CurrentRevision(boundedTestContext(t))
	if err != nil {
		t.Fatal(err)
	}
	return ok
}

func exitCodeOf(t *testing.T, err error) int {
	t.Helper()
	if err == nil {
		return 0
	}
	var exit *exec.ExitError
	if !errors.As(err, &exit) {
		t.Fatalf("the child did not run: %v", err)
	}
	return exit.ExitCode()
}

// With NO controlling terminal the command refuses before the write — even
// with the exact phrase piped to stdin, which is what an agent would do — and
// the refusal ends with the fix.
func TestTransitionsWithoutAControllingTerminalRefusesAndWritesNothing(t *testing.T) {
	flags := transitionsFixture(t)
	cmd, stdout, stderr := ptyChild(t, invocation{verb: "transitions", flags: flags}.args(t))
	cmd.Stdin = strings.NewReader(expectedTransitionsPhrase(t, flags) + "\n")
	ptytest.WithoutControllingTerminal(cmd)
	code := exitCodeOf(t, cmd.Run())

	if code != exitStop || stopLine(stderr.String()) != "STOP fence=tty reason=no-controlling-terminal" {
		t.Fatalf("no controlling terminal: exit %d, want %d with the no-controlling-terminal STOP\nstdout %q\nstderr %q",
			code, exitStop, stdout.String(), stderr.String())
	}
	if !strings.HasSuffix(strings.TrimSpace(stderr.String()), ttyFix) {
		t.Fatalf("the refusal does not end with the fix %q:\n%s", ttyFix, stderr.String())
	}
	if registryHeadExists(t, flags["store"]) {
		t.Fatal("a refused run wrote a transition-registry head")
	}
}

// With the pty as the controlling terminal and stdin redirected from
// /dev/null, the prompt appears on the pty, the phrase typed there is read,
// and the write happens. This is the IDE-pane case; it is exactly as strong as
// QUICKSTART's former `< /dev/tty` advice, because the line comes from the
// same device either way.
func TestTransitionsReadsTheConfirmationFromThePtyWhenStdinIsDevNull(t *testing.T) {
	flags := transitionsFixture(t)
	phrase := expectedTransitionsPhrase(t, flags)
	master, slave, err := ptytest.Open()
	if err != nil {
		t.Fatalf("open pty: %v", err)
	}
	defer func() { _ = master.Close() }()
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = devNull.Close() }()

	cmd, stdout, stderr := ptyChild(t, invocation{verb: "transitions", flags: flags}.args(t))
	cmd.Stdin = devNull
	ptytest.WithControllingTerminal(cmd, slave)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	_ = slave.Close()
	tr := ptytest.Record(master)
	waited := make(chan error, 1)
	go func() { waited <- cmd.Wait() }()

	if err := tr.WaitFor(phrase, 30*time.Second); err != nil {
		_ = cmd.Process.Kill()
		<-waited
		t.Fatalf("%v\nstdout %q\nstderr %q", err, stdout.String(), stderr.String())
	}
	if _, err := master.Write([]byte(phrase + "\n")); err != nil {
		t.Fatal(err)
	}
	var runErr error
	select {
	case runErr = <-waited:
	case <-time.After(60 * time.Second):
		_ = cmd.Process.Kill()
		t.Fatalf("the child did not finish after the phrase was typed\npty %q\nstderr %q", tr.String(), stderr.String())
	}
	if code := exitCodeOf(t, runErr); code != exitOK {
		t.Fatalf("exit %d, want 0\nstdout %q\nstderr %q\npty %q", code, stdout.String(), stderr.String(), tr.String())
	}
	if !strings.Contains(stdout.String(), "published transition registry revision 1") {
		t.Fatalf("stdout %q, want the revision-1 publication line", stdout.String())
	}
	if strings.Contains(stdout.String(), phrase) {
		t.Fatalf("the prompt went to stdout %q; it belongs on the terminal the line is read from", stdout.String())
	}
	if !registryHeadExists(t, flags["store"]) {
		t.Fatal("instrument failure: the confirmed run left no registry head")
	}
}
