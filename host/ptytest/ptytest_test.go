//go:build darwin || linux

package ptytest

import (
	"os/exec"
	"strings"
	"testing"
	"time"
)

// Both shapes, measured with /bin/sh: with the pty as controlling terminal the
// child's /dev/tty is the pty (what it writes there reaches the master, and
// what the master writes is what it reads); without one, opening /dev/tty
// fails. Each arm is the other's control.
func TestTheTwoShapesOfControllingTerminal(t *testing.T) {
	master, slave, err := Open()
	if err != nil {
		t.Fatalf("open pty: %v", err)
	}
	defer func() { _ = master.Close() }()
	cmd := exec.Command("/bin/sh", "-c", `printf 'prompt> ' >/dev/tty; read x </dev/tty; echo "got:$x"`)
	var out Buffer
	cmd.Stdout = &out
	WithControllingTerminal(cmd, slave)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	_ = slave.Close()
	tr := Record(master)
	if err := tr.WaitFor("prompt> ", 10*time.Second); err != nil {
		t.Fatal(err)
	}
	if _, err := master.Write([]byte("typed\n")); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		t.Fatalf("child with a controlling terminal: %v (stdout %q, pty %q)", err, out.String(), tr.String())
	}
	if strings.TrimSpace(out.String()) != "got:typed" {
		t.Fatalf("stdout %q, want the line typed at the pty", out.String())
	}

	none := exec.Command("/bin/sh", "-c", `: </dev/tty`)
	WithoutControllingTerminal(none)
	if err := none.Run(); err == nil {
		t.Fatal("a child with no controlling terminal opened /dev/tty")
	}
}
