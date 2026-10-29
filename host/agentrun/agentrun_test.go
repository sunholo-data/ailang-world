package agentrun

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/childenv"
)

func TestAgentrunWallCapKillsDescendants(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "bin")
	home := filepath.Join(dir, "home")
	void := filepath.Join(dir, "void")
	for _, path := range []string{bin, home, void} {
		if err := os.Mkdir(path, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("HOME", home)
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	heartbeat := filepath.Join(dir, "heartbeat")
	agent := filepath.Join(bin, "codex")
	// The descendant holds inherited descriptors and writes repeatedly; killing
	// only the direct child leaves the heartbeat growing after Run returns.
	script := "#!/bin/sh\n(while :; do printf x >> \"$HEARTBEAT\"; /bin/sleep 0.05; done) &\nwait\n"
	if err := os.WriteFile(agent, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	out, err := os.CreateTemp(dir, "out")
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()
	errout, err := os.CreateTemp(dir, "err")
	if err != nil {
		t.Fatal(err)
	}
	defer errout.Close()
	env := []string{}
	for _, kv := range childenv.Scrubbed(os.Environ()) {
		if !strings.HasPrefix(kv, "WORLD_SESSION=") {
			env = append(env, kv)
		}
	}
	env = append(env, "HEARTBEAT="+heartbeat)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	before := time.Now()
	r, err := Run(ctx, []string{agent}, void, env, out, errout, time.Second)
	if err != nil || r.Reason != "timeout" || time.Since(before) > 3*time.Second {
		t.Fatalf("%+v %v", r, err)
	}
	first, err := os.ReadFile(heartbeat)
	if err != nil || len(first) == 0 {
		t.Fatal("child never ran", err)
	}
	time.Sleep(200 * time.Millisecond)
	second, _ := os.ReadFile(heartbeat)
	if string(first) != string(second) {
		t.Fatal("descendant survived wall cap")
	}
}

func TestAgentrunForwardsSIGTERM(t *testing.T) {
	root := t.TempDir()
	bin, void := filepath.Join(root, "bin"), filepath.Join(root, "void")
	for _, dir := range []string{bin, void} {
		if err := os.Mkdir(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("HOME", t.TempDir())
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	ready := filepath.Join(root, "ready")
	agent := filepath.Join(bin, "claude")
	script := "#!/bin/sh\ntrap 'printf received; exit 0' TERM\nprintf ready > \"$READY\"\nwhile :; do /bin/sleep 1; done\n"
	if err := os.WriteFile(agent, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	out, err := os.CreateTemp(root, "stdout")
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()
	errout, err := os.CreateTemp(root, "stderr")
	if err != nil {
		t.Fatal(err)
	}
	defer errout.Close()
	env := append(childenv.Scrubbed([]string{"PATH=" + os.Getenv("PATH"), "HOME=" + os.Getenv("HOME")}), "READY="+ready)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	type outcome struct {
		result Result
		err    error
	}
	done := make(chan outcome, 1)
	go func() {
		r, err := Run(ctx, []string{agent}, void, env, out, errout, 3*time.Second)
		done <- outcome{r, err}
	}()
	deadline := time.Now().Add(2 * time.Second)
	for {
		if _, err := os.Stat(ready); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("fake agent did not become ready")
		}
		time.Sleep(10 * time.Millisecond)
	}
	// Notify is installed before Start, so the ready file proves this signal
	// reaches Run's handler and can be forwarded safely to the child's group.
	if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	got := <-done
	if got.err != nil || got.result.Reason != "signal" || got.result.ExitCode != 0 {
		t.Fatalf("%+v %v", got.result, got.err)
	}
	_, _ = out.Seek(0, 0)
	data, err := os.ReadFile(out.Name())
	if err != nil || string(data) != "received" {
		t.Fatal("signal did not reach child", string(data), err)
	}
}
