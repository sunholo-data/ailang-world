package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDoctorCodexApprovalDrift(t *testing.T) {
	dir, bin := hostsLauncherFixture(t)
	t.Setenv("PATH", bin+":"+os.Getenv("PATH"))
	config := filepath.Join(dir, "codex", "config.toml")
	b, err := os.ReadFile(config)
	if err != nil {
		t.Fatal(err)
	}
	b = bytes.ReplaceAll(b, []byte("approval_policy = \"never\"\n"), nil)
	if err := os.WriteFile(config, b, 0600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	doctorHost(&doctorReport{w: &out}, "codex", dir)
	for _, want := range []string{"approval_policy", "drift", "hosts init", "--replace"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("missing %s: %s", want, out.String())
		}
	}
}
