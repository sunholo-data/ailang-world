package agentrun

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAgentrunInteractiveStdinIsTTY(t *testing.T) {
	agent, dir, env, tty := interactiveAgent(t, "IFS= read -r line || exit 23\nprintf 'read:%s\\n' \"$line\"\n")
	if _, err := tty.WriteString("terminal-input\n"); err != nil {
		t.Fatal(err)
	}
	if _, err := tty.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	result, err := RunInteractive(ctx, []string{agent}, dir, env, tty, 5*time.Second)
	data, readErr := os.ReadFile(tty.Name())
	if err != nil || readErr != nil || result.ExitCode != 0 || !strings.Contains(string(data), "read:terminal-input\n") {
		t.Fatal(result, err, readErr, string(data))
	}
}

func TestAgentrunInteractiveCancelReason(t *testing.T) {
	for _, kind := range []string{"canceled", "timeout"} {
		t.Run(kind, func(t *testing.T) {
			agent, dir, env, tty := interactiveAgent(t, "echo ready > ready\nexec /bin/sleep 30\n")
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			type outcome struct {
				result Result
				err    error
			}
			done := make(chan outcome, 1)
			go func() {
				result, err := RunInteractive(ctx, []string{agent}, dir, env, tty, 5*time.Second)
				done <- outcome{result, err}
			}()
			for deadline := time.Now().Add(4 * time.Second); ; time.Sleep(10 * time.Millisecond) {
				if _, err := os.Stat(filepath.Join(dir, "ready")); err == nil {
					break
				}
				if time.Now().After(deadline) {
					t.Fatal("child not ready")
				}
			}
			if kind == "canceled" {
				cancel()
			}
			select {
			case got := <-done:
				if got.err != nil || got.result.Reason != kind {
					t.Fatal(got)
				}
			case <-time.After(7 * time.Second):
				t.Fatal("wait did not finish")
			}
		})
	}
}
