package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/toolprobe"
)

func TestWorldContextAndStatusline(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if err := os.MkdirAll(filepath.Join(home, ".claude"), 0700); err != nil {
		t.Fatal(err)
	}
	// The command proves that the original stdin is forwarded.
	settings := `{"statusLine":{"command":"cat"}}`
	if err := os.WriteFile(filepath.Join(home, ".claude", "settings.json"), []byte(settings), 0600); err != nil {
		t.Fatal(err)
	}
	dir, bin := hostsLauncherFixture(t)
	// hostsLauncherFixture supplies a different HOME; generate again with our settings.
	t.Setenv("HOME", home)
	if code, out := hostsInitTest(t, dir, "--replace"); code != 0 {
		t.Fatalf("%d: %s", code, out)
	}
	fake := "#!/bin/sh\n[ \"$FAKE_CODE\" = 0 ] || exit \"$FAKE_CODE\"\nprintf '%s\\n' " + shellQuote(`{"episode":"w223d","worktree":"/work/w223d","secondsLeft":3661,"addr":"http://127.0.0.1:47330","grants":[{"effect":"Workspace.Read","left":9,"budget":10}]}`) + "\n"
	if err := os.WriteFile(filepath.Join(bin, "ailang-worldd"), []byte(fake), 0700); err != nil {
		t.Fatal(err)
	}
	invoke := func(name, code string) string {
		t.Helper()
		ctx, cancel := budgetContext(5 * time.Second)
		defer cancel()
		input := filepath.Join(t.TempDir(), "stdin")
		if err := os.WriteFile(input, []byte("quota-input"), 0600); err != nil {
			t.Fatal(err)
		}
		// toolprobe's output API has no stdin; a small fixture wrapper supplies it.
		wrapper := filepath.Join(t.TempDir(), "run")
		if err := os.WriteFile(wrapper, []byte("#!/bin/sh\nexec "+shellQuote(filepath.Join(dir, "bin", name))+" <"+shellQuote(input)+"\n"), 0700); err != nil {
			t.Fatal(err)
		}
		out, err := toolprobe.OutputEnv(ctx, wrapper, "", []string{
			"PATH=" + bin + ":" + os.Getenv("PATH"), "WORLD_SESSION=never-cache-this-token", "FAKE_CODE=" + code,
		})
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		return string(out)
	}
	context := invoke("world-context", "0")
	for _, want := range []string{"inside AILANG World (confined)", "Episode w223d · worktree /work/w223d · session 1h1m", worldToolMap, "Grants left: Workspace.Read 9/10", worldBoundary} {
		if !strings.Contains(context, want) {
			t.Fatalf("missing %q: %s", want, context)
		}
	}
	failure := invoke("world-context", "11")
	if strings.Count(failure, "\n") != 1 || !strings.Contains(failure, "operator") {
		t.Fatal(failure)
	}
	status := invoke("world-statusline", "0")
	if !strings.Contains(status, "World · w223d · 1h1m · ok │ quota-input") {
		t.Fatal(status)
	}
	// A cached success is reused even when the daemon subsequently fails.
	if got := invoke("world-statusline", "10"); got != status {
		t.Fatalf("cache: %s", got)
	}
	entries, err := os.ReadDir(filepath.Join(dir, "state"))
	if err != nil || len(entries) != 1 {
		t.Fatalf("cache entries: %v %v", entries, err)
	}
	for _, entry := range entries {
		p := filepath.Join(dir, "state", entry.Name())
		b, _ := os.ReadFile(p)
		info, _ := os.Stat(p)
		if bytes.Contains(b, []byte("never-cache-this-token")) || info.Mode().Perm() != 0600 {
			t.Fatal("unsafe cache")
		}
	}
	for _, command := range []string{"", "exit 1"} {
		settings := "{\"statusLine\":{\"command\":" + jsString(command) + "}}"
		if err := os.WriteFile(filepath.Join(home, ".claude", "settings.json"), []byte(settings), 0600); err != nil {
			t.Fatal(err)
		}
		if code, out := hostsInitTest(t, dir, "--replace"); code != 0 {
			t.Fatalf("%d: %s", code, out)
		}
		if got := invoke("world-statusline", "0"); strings.Contains(got, "│") || !strings.Contains(got, "World · w223d") {
			t.Fatal(got)
		}
	}
}
