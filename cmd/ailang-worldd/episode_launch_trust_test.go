package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInteractiveLaunchTrustOverride(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "with.dot\"quote")
	if err := os.Mkdir(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	real, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := "projects={ \"" + strings.ReplaceAll(real, "\"", "\\\"") + "\" = { trust_level = \"trusted\" } }"
	for _, agent := range []string{"codex", "claude"} {
		args := interactiveLaunchArgv("host", agent, "", "session", "context", "prompt", dir, nil)
		count := 0
		for i, arg := range args {
			if strings.HasPrefix(arg, "projects=") {
				count++
				if arg != want || i == 0 || args[i-1] != "-c" {
					t.Fatalf("override = %q, want %q", arg, want)
				}
			}
		}
		expected := 1
		if agent == "claude" {
			expected = 0
		}
		if count != expected {
			t.Fatalf("%s: %d overrides in %q", agent, count, args)
		}
	}
}

func TestInteractiveLaunchTrustCleanFallback(t *testing.T) {
	args := interactiveLaunchArgv("host", "codex", "", "", "", "", "/absent/a/../b", nil)
	if args[len(args)-1] != `projects={ "/absent/b" = { trust_level = "trusted" } }` {
		t.Fatal(args)
	}
}

func TestTOMLBasicString(t *testing.T) {
	if got := tomlBasicString("a\\\"\n\t\x00\x7f😀"); got != "\"a\\\\\\\"\\u000a\\u0009\\u0000\\u007f😀\"" {
		t.Fatal(got)
	}
}
