package agentrun

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/childenv"
)

func interactiveAgent(t *testing.T, body string) (string, string, []string, *os.File) {
	t.Helper()
	dir := t.TempDir()
	agent := filepath.Join(dir, "agent")
	script := "#!/bin/sh\nif [ \"$WORLD_TEST_WARM\" = 1 ]; then exit 0; fi\n" + body
	if err := os.WriteFile(agent, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	env := childenv.Scrubbed(os.Environ())
	warmAgent(t, agent, dir, env)
	tty, err := os.CreateTemp(dir, "tty")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = tty.Close() })
	return agent, dir, env, tty
}

func TestAgentrunInteractiveGroupAndOutput(t *testing.T) {
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("WORLD_TEST_BINARY", binary)
	t.Setenv("WORLD_TEST_GROUP_CHILD", "1")
	agent, dir, env, tty := interactiveAgent(t, "exec \"$WORLD_TEST_BINARY\" -test.run=^TestAgentrunInteractiveGroupChild$\n")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := RunInteractive(ctx, []string{agent}, dir, env, tty, 5*time.Second)
	data, _ := os.ReadFile(tty.Name())
	want := fmt.Sprint(syscall.Getpgrp())
	if err != nil || r.ExitCode != 0 || r.Reason != "exit" || strings.Fields(string(data))[0] != want || !strings.Contains(string(data), "stdout\nstderr\n") {
		t.Fatal(r, err, string(data), want)
	}
}

// The child reports its own group without a process-table query (ps is denied
// by the confined exec profile). The normal suite does not enter child mode.
func TestAgentrunInteractiveGroupChild(t *testing.T) {
	if os.Getenv("WORLD_TEST_GROUP_CHILD") != "1" {
		return
	}
	fmt.Fprintf(os.Stdout, "%d\nstdout\n", syscall.Getpgrp())
	fmt.Fprintln(os.Stderr, "stderr")
	os.Exit(0)
}

func interactiveSignal(t *testing.T, sig syscall.Signal) {
	t.Helper()
	agent, dir, env, tty := interactiveAgent(t, `trap 'echo int > int-marker' INT
trap 'echo term > term-marker; exit 0' TERM
echo ready > ready
while [ ! -f finish ]; do /bin/sleep 0.05; done
`)
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	type outcome struct {
		r   Result
		err error
	}
	done := make(chan outcome, 1)
	go func() {
		r, err := RunInteractive(ctx, []string{agent}, dir, env, tty, 8*time.Second)
		done <- outcome{r, err}
	}()
	for deadline := time.Now().Add(5 * time.Second); ; {
		if _, err := os.Stat(filepath.Join(dir, "ready")); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("child did not become ready")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err := syscall.Kill(os.Getpid(), sig); err != nil {
		t.Fatal(err)
	}
	if sig == syscall.SIGINT {
		select {
		case got := <-done:
			t.Fatal("SIGINT ended wait", got)
		case <-time.After(100 * time.Millisecond):
		}
		if _, err := os.Stat(filepath.Join(dir, "int-marker")); !os.IsNotExist(err) {
			t.Fatal("SIGINT forwarded")
		}
		if err := os.WriteFile(filepath.Join(dir, "finish"), nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	select {
	case got := <-done:
		want := "exit"
		if sig == syscall.SIGTERM {
			want = "signal"
			if _, err := os.Stat(filepath.Join(dir, "term-marker")); err != nil {
				t.Fatal("SIGTERM not forwarded", err)
			}
		}
		if got.err != nil || got.r.Reason != want {
			t.Fatal(got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("child wait did not finish")
	}
}

func TestAgentrunInteractiveSIGINT(t *testing.T) { interactiveSignal(t, syscall.SIGINT) }

func TestAgentrunInteractiveSIGTERM(t *testing.T) { interactiveSignal(t, syscall.SIGTERM) }

func TestAgentrunInteractiveWallCap(t *testing.T) {
	agent, dir, env, tty := interactiveAgent(t, "exec /bin/sleep 30\n")
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	before := time.Now()
	r, err := RunInteractive(ctx, []string{agent}, dir, env, tty, 5*time.Second)
	if err != nil || r.Reason != "timeout" || time.Since(before) > 7*time.Second {
		t.Fatal(r, err, time.Since(before))
	}
}
