package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDoctorCodexStateSuffix(t *testing.T) {
	prefix := []byte("[mcp_servers.world]\nurl = \"http://fixture\"\n")
	for _, tc := range []struct {
		name, suffix string
		bad          int
	}{
		{"state", "\n[tui]\nscreen_reader_detection_done = true\n\n[hooks.state]\n\n[hooks.state.\"/Users/test/.ailang-world/hosts/codex/hooks.json:session_start:0:0\"]\ntrusted_hash = \"sha256:5ef8\"\n", 0},
		{"projects", "[projects.\"/x.y\"]\ntrust_level = \"trusted\"\n", 0},
		{"features", "[features]\nshell_tool = true\n", 3},
		{"mcp", "[mcp_servers.other]\n", 3},
		{"bare", "sandbox_mode = \"danger-full-access\"\n", 3},
		{"hooks", "[hooks]\nfoo = 1\n", 3},
		{"project-key", "[projects.\"/x\"]\nsandbox_mode = \"x\"\n", 4},
		{"array", "[[tui]]\n", 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, bad := codexStateSuffix(prefix, append(append([]byte{}, prefix...), []byte(tc.suffix)...))
			if bad != tc.bad {
				t.Fatalf("line %d, want %d", bad, tc.bad)
			}
		})
	}
	changed := append([]byte{}, prefix...)
	changed[0] = 'x'
	if _, bad := codexStateSuffix(prefix, changed); bad != 1 {
		t.Fatalf("prefix change: %d", bad)
	}
}

func TestDoctorHostsCodexStateLaunchPreflight(t *testing.T) {
	dir, bin := hostsLauncherFixture(t)
	t.Setenv("PATH", bin+":/usr/bin:/bin")
	path := filepath.Join(dir, "codex", "config.toml")
	prefix, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	suffix := "\n[tui]\nscreen_reader_detection_done = true\n[hooks.state]\n[hooks.state.\"/fixture/hooks.json:session_start:0:0\"]\ntrusted_hash = \"sha256:fixture\"\n[projects.\"/x.y\"]\ntrust_level = \"trusted\"\n"
	if err := os.WriteFile(path, append(prefix, []byte(suffix)...), 0600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	doctorHost(&doctorReport{w: &out}, "codex", dir)
	if !strings.Contains(out.String(), "World's config plus Codex state (8 lines)") || strings.Contains(out.String(), "digest differs") {
		t.Fatal(out.String())
	}
	if _, err := launchLauncher("codex", dir, defaultEpisodeAddr); err != nil {
		t.Fatal(err)
	}
}
