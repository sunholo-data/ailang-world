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
)

// ---------------------------------------------------------------------------
// `session new` END TO END, in a real child process with realSessionEnv —
// the production /dev/tty — rather than an injected terminal. The two tests
// differ only in the kernel's answer to open("/dev/tty"): none (setsid, the
// agent-harness shape) versus a pty the test answers on. In both, stdin is
// NOT the terminal.
// ---------------------------------------------------------------------------

const sessionPtyChildArgs = "AILANG_WORLDD_PTY_CHILD_ARGS"

// TestSessionPtyChildProcess is not a test of its own: it is the command, run
// in the child the tests below start.
func TestSessionPtyChildProcess(t *testing.T) {
	raw := os.Getenv(sessionPtyChildArgs)
	if raw == "" {
		t.Skip("the child process of the pty tests; it runs only when they start it")
	}
	var args []string
	if err := json.Unmarshal([]byte(raw), &args); err != nil {
		t.Fatal(err)
	}
	os.Exit(run(args, os.Stdout, os.Stderr))
}

func sessionPtyChild(t *testing.T, args []string) (*exec.Cmd, *ptytest.Buffer, *ptytest.Buffer) {
	t.Helper()
	encoded, err := json.Marshal(args)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestSessionPtyChildProcess$", "-test.count=1")
	cmd.Env = append(os.Environ(), sessionPtyChildArgs+"="+string(encoded))
	var stdout, stderr ptytest.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	return cmd, &stdout, &stderr
}

func childExitCode(t *testing.T, err error) int {
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

// No controlling terminal: refused before anything is provisioned — even with
// "y" piped to stdin — and the refusal ends with the fix.
func TestSessionNewWithoutAControllingTerminalRefusesEndToEnd(t *testing.T) {
	f := newEpisodeFixture(t)
	cmd, stdout, stderr := sessionPtyChild(t, append([]string{"session", "new"}, f.args("ep1", "--preset", "se-tools")...))
	cmd.Stdin = strings.NewReader("y\n")
	ptytest.WithoutControllingTerminal(cmd)
	code := childExitCode(t, cmd.Run())
	if code != exitUsage || !strings.Contains(stderr.String(), "refusing: no controlling terminal") {
		t.Fatalf("exit %d, want %d refusing for no controlling terminal\nstdout %q\nstderr %q", code, exitUsage, stdout.String(), stderr.String())
	}
	if !strings.HasSuffix(strings.TrimSpace(stderr.String()), sessionTTYFix) {
		t.Fatalf("the refusal does not end with the fix %q:\n%s", sessionTTYFix, stderr.String())
	}
	f.assertNothingProvisioned(t, "ep1")
}

// The pty as controlling terminal, stdin from /dev/null (the IDE-pane shape):
// the y/N appears on the pty, the answer typed there is read, and the worktree
// and the session are provisioned.
func TestSessionNewReadsTheAnswerFromThePtyWhenStdinIsDevNull(t *testing.T) {
	f := newEpisodeFixture(t)
	out := filepath.Join(t.TempDir(), "ep1.session")
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

	cmd, stdout, stderr := sessionPtyChild(t, append([]string{"session", "new"}, f.args("ep1", "--preset", "se-tools", "--out", out)...))
	cmd.Stdin = devNull
	ptytest.WithControllingTerminal(cmd, slave)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	_ = slave.Close()
	tr := ptytest.Record(master)
	waited := make(chan error, 1)
	go func() { waited <- cmd.Wait() }()
	if err := tr.WaitFor("[y/N]", 30*time.Second); err != nil {
		_ = cmd.Process.Kill()
		<-waited
		t.Fatalf("%v\nstderr %q", err, stderr.String())
	}
	if _, err := master.Write([]byte("y\n")); err != nil {
		t.Fatal(err)
	}
	var runErr error
	select {
	case runErr = <-waited:
	case <-time.After(60 * time.Second):
		_ = cmd.Process.Kill()
		t.Fatalf("the child did not finish after y was typed\npty %q\nstderr %q", tr.String(), stderr.String())
	}
	if code := childExitCode(t, runErr); code != exitOK {
		t.Fatalf("exit %d, want 0\nstdout %q\nstderr %q\npty %q", code, stdout.String(), stderr.String(), tr.String())
	}
	if !isGitWorktree(filepath.Join(f.ws, "ep1")) || len(f.rows(t)) != 1 {
		t.Fatalf("not provisioned: worktree %v, rows %d", isGitWorktree(filepath.Join(f.ws, "ep1")), len(f.rows(t)))
	}
	if tok, err := os.ReadFile(out); err != nil || !isHex64(string(tok)) {
		t.Fatalf("token file: %q %v", tok, err)
	}
}
