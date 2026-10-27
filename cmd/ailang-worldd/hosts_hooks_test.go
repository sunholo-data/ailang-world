package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHostAwarenessHooks(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	for _, rel := range []string{".codex/herdr-agent-state.sh", ".pi/agent/extensions/herdr-agent-state.ts"} {
		p := filepath.Join(home, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("user-owned"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	dir := filepath.Join(t.TempDir(), "hosts")
	if code, out := hostsInitTest(t, dir); code != 0 {
		t.Fatalf("%d: %s", code, out)
	}
	files := hostFiles(dir, "http://127.0.0.1:47330")
	var hook struct {
		Hooks map[string][]struct {
			Hooks []struct {
				Command string `json:"command"`
			} `json:"hooks"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal(files["codex/hooks.json"], &hook); err != nil {
		t.Fatal(err)
	}
	if got := hook.Hooks["SessionStart"][0].Hooks[0].Command; got != "bash "+shellQuote(filepath.Join(home, ".codex/herdr-agent-state.sh"))+" session" {
		t.Fatal(got)
	}
	pi := string(files["bin/pi-world"])
	if !strings.Contains(pi, "-e "+shellQuote(filepath.Join(home, ".pi/agent/extensions/herdr-agent-state.ts"))) {
		t.Fatal(pi)
	}
	if err := json.Unmarshal(files["claude-plugin/hooks/hooks.json"], &hook); err != nil {
		t.Fatal(err)
	}
	if got := hook.Hooks["SessionStart"][0].Hooks[0].Command; got != shellQuote(filepath.Join(dir, "bin/world-context")) {
		t.Fatal(got)
	}
	for _, rel := range []string{".codex/herdr-agent-state.sh", ".pi/agent/extensions/herdr-agent-state.ts"} {
		b, _ := os.ReadFile(filepath.Join(home, rel))
		if string(b) != "user-owned" {
			t.Fatal("user hook changed")
		}
	}
}
