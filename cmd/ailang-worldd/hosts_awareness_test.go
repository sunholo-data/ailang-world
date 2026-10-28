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
	// The scripts call the ailang-worldd that ran hosts init; here, the fake below.
	prev := hostsWorlddPath
	hostsWorlddPath = func() string { return filepath.Join(bin, "ailang-worldd") }
	t.Cleanup(func() { hostsWorlddPath = prev })
	if code, out := hostsInitTest(t, dir, "--replace"); code != 0 {
		t.Fatalf("%d: %s", code, out)
	}
	fake := "#!/bin/sh\n[ \"$FAKE_CODE\" = 0 ] || exit \"$FAKE_CODE\"\nprintf '%s\\n' " + shellQuote(`{"episode":"w223d","worktree":"/work/w223d","secondsLeft":3661,"git":{"base":"abcdef0123456789","base_date":"2026-10-28T09:00:00Z","base_subject":"base fixture","behind":3,"ahead":0,"dirty":false},"addr":"http://127.0.0.1:47330","grants":[{"effect":"Workspace.Read","left":9,"budget":10}]}`) + "\n"
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
	for _, want := range []string{"inside AILANG World (confined)", "Base abcdef0 (2026-10-28T09:00:00Z, \"base fixture\") · dev is 3 commits ahead of it", "Episode w223d · worktree /work/w223d · session 1h1m", worldToolMap, "Grants left: Workspace.Read 9/10", worldBoundary} {
		if !strings.Contains(context, want) {
			t.Fatalf("missing %q: %s", want, context)
		}
	}
	for _, tc := range []struct {
		git  string
		want string
	}{
		{`{"error":"local git metadata timed out (2s budget)"}`, "Base unknown: local git metadata timed out (2s budget)"},
		{`{"base":"abcdef0123456789","base_date":"date","base_subject":"subject","behind":3,"dirty":true}`, "dev is 3 commits ahead of it · worktree has uncommitted changes"},
		{"null", ""},
	} {
		body := `{"episode":"w223d","worktree":"/work/w223d","secondsLeft":3661,"addr":"http://127.0.0.1:47330","grants":[],"git":` + tc.git + `}`
		script := "#!/bin/sh\nprintf '%s\n' " + shellQuote(body) + "\n"
		if err := os.WriteFile(filepath.Join(bin, "ailang-worldd"), []byte(script), 0700); err != nil {
			t.Fatal(err)
		}
		got := invoke("world-context", "0")
		if (tc.want == "" && strings.Contains(got, "Base ")) || (tc.want != "" && !strings.Contains(got, tc.want)) || strings.Contains(got, "never-cache-this-token") {
			t.Fatal(got)
		}
	}
	if err := os.WriteFile(filepath.Join(bin, "ailang-worldd"), []byte(fake), 0700); err != nil {
		t.Fatal(err)
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

// TestHostsAwarenessCallsTheInitBinary pins the 27 Oct live failure: the
// scripts called a bare ailang-worldd, which was not on the agent's PATH, and
// world-context said "daemon down". They call the binary that ran hosts init,
// and a missing one is named as missing.
func TestHostsAwarenessCallsTheInitBinary(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	writeTestPiAdapter(t, defaultPiMCPAdapter())
	missing := filepath.Join(t.TempDir(), "gone", "ailang-worldd")
	prev := hostsWorlddPath
	hostsWorlddPath = func() string { return missing }
	t.Cleanup(func() { hostsWorlddPath = prev })
	dir := filepath.Join(t.TempDir(), "hosts")
	if code, out := hostsInitTest(t, dir); code != 0 {
		t.Fatalf("%d: %s", code, out)
	}
	for _, name := range []string{"bin/world-context", "bin/world-statusline", "pi/world-status.ts"} {
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), missing) || strings.Contains(string(b), "$(ailang-worldd ") || strings.Contains(string(b), `execFile("ailang-worldd"`) {
			t.Fatalf("%s does not call the init binary %s:\n%s", name, missing, b)
		}
	}
	ctx, cancel := budgetContext(5 * time.Second)
	defer cancel()
	out, err := toolprobe.OutputEnv(ctx, filepath.Join(dir, "bin", "world-context"), "", []string{"PATH=" + os.Getenv("PATH"), "WORLD_SESSION=x"})
	if err != nil {
		t.Fatal(err)
	}
	if want := "no ailang-worldd at " + missing; !strings.Contains(string(out), want) || strings.Contains(string(out), "daemon down") {
		t.Fatalf("want %q, got %s", want, out)
	}
}
