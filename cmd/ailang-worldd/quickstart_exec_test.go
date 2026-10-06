package main

// Row 140 M3 (w-workspace-exec-toolchain-effect §4.6, §6 AC3.1/AC3.5,
// coding-standards S7): the four `serve --exec-*` flags reach the daemon's
// Config exactly as given, `serve --help` documents each, and QUICKSTART
// §10's serve line names only documented flags and parses.

import (
	"bytes"
	"context"
	"io"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/daemon"
)

// captureServeConfig swaps runDaemon for one that records the Config.
func captureServeConfig(t *testing.T) *daemon.Config {
	t.Helper()
	got := &daemon.Config{}
	prev := runDaemon
	runDaemon = func(_ context.Context, cfg daemon.Config, _ io.Writer) error {
		*got = cfg
		return nil
	}
	t.Cleanup(func() { runDaemon = prev })
	return got
}

// The wiring: every exec flag lands in the Config serve hands the daemon.
func TestServeExecFlagsReachTheDaemonConfig(t *testing.T) {
	got := captureServeConfig(t)
	var out, errw bytes.Buffer
	code := run([]string{"serve", "--db", "/x/world.db",
		"--exec-profile", "/p/a.json", "--exec-profile", "/p/b.json",
		"--exec-episode-project", "go1=a", "--exec-episode-project", "ts1=b",
		"--exec-sandbox", "/s/node_modules", "--exec-node", "/n/node", "--exec-max-output-bytes", "1048576"}, &out, &errw)
	if code != exitOK {
		t.Fatalf("serve exit %d: %s", code, errw.String())
	}
	if !reflect.DeepEqual(got.ExecProfiles, []string{"/p/a.json", "/p/b.json"}) ||
		!reflect.DeepEqual(got.ExecEpisodeProjects, []string{"go1=a", "ts1=b"}) ||
		got.ExecSandbox != "/s/node_modules" || got.ExecNode != "/n/node" || got.ExecMaxOutputBytes != 1<<20 {
		t.Fatalf("serve built Config %+v", *got)
	}
	// Unset, every exec field is zero: the daemon then binds the typed
	// "no exec profile configured" refusal.
	*got = daemon.Config{}
	if code := run([]string{"serve", "--db", "/x/world.db"}, &out, &errw); code != exitOK {
		t.Fatalf("serve exit %d", code)
	}
	if got.ExecProfiles != nil || got.ExecEpisodeProjects != nil || got.ExecSandbox != "" || got.ExecNode != "" || got.ExecMaxOutputBytes != 0 {
		t.Fatalf("serve with no exec flag built %+v", *got)
	}
	// A non-positive cap is a usage error, before the daemon.
	for _, bad := range []string{"0", "-1", "64MiB"} {
		errw.Reset()
		if code := run([]string{"serve", "--db", "/x/world.db", "--exec-max-output-bytes", bad}, &out, &errw); code != exitUsage ||
			!strings.Contains(errw.String(), "not a positive byte count") {
			t.Fatalf("--exec-max-output-bytes %s: exit %d, %s", bad, code, errw.String())
		}
	}
}

// section10ServeLines returns QUICKSTART §10's `/tmp/ailang-worldd serve`
// lines, continuations joined, as argv after the binary.
func section10ServeLines(t *testing.T) [][]string {
	t.Helper()
	guide, err := os.ReadFile("../../docs/QUICKSTART.md")
	if err != nil {
		t.Fatal(err)
	}
	_, section, ok := strings.Cut(string(guide), "### 10. Build and test non-AILANG projects")
	if !ok {
		t.Fatal("QUICKSTART §10 is absent")
	}
	section, _, _ = strings.Cut(section, "### 11. ")
	var cmds [][]string
	for _, line := range strings.Split(strings.ReplaceAll(section, "\\\n", " "), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "/tmp/ailang-worldd serve ") {
			cmds = append(cmds, strings.Fields(strings.TrimSuffix(line, "&"))[1:])
		}
	}
	return cmds
}

// AC3.5 (S7): every flag §10's serve line names is documented by serve
// --help, the line parses, and it uses all five exec flags.
func TestQuickstartSection10FlagsMatchTheCLI(t *testing.T) {
	var help, ignored bytes.Buffer
	if got := run([]string{"serve", "--help"}, &help, &ignored); got != exitOK {
		t.Fatalf("serve --help exit = %d", got)
	}
	_, serveHelp, ok := strings.Cut(help.String(), "serve flags:")
	if !ok {
		t.Fatalf("serve --help has no serve flags block:\n%s", help.String())
	}
	for _, flag := range []string{"--exec-profile", "--exec-episode-project", "--exec-sandbox", "--exec-node", "--exec-max-output-bytes"} {
		if !regexp.MustCompile(`(?m)^\s+` + regexp.QuoteMeta(flag) + `\b`).MatchString(serveHelp) {
			t.Errorf("serve --help does not document %s", flag)
		}
	}
	lines := section10ServeLines(t)
	if len(lines) != 1 {
		t.Fatalf("QUICKSTART §10 holds %d serve lines, want 1", len(lines))
	}
	argv := lines[0]
	seen := map[string]bool{}
	for _, tok := range argv[1:] {
		if flagToken.MatchString(tok) {
			seen[tok] = true
			if !regexp.MustCompile(`(?m)^\s+` + regexp.QuoteMeta(tok) + `\b`).MatchString(serveHelp) {
				t.Errorf("QUICKSTART §10 serve flag %s is not documented by serve --help", tok)
			}
		}
	}
	for _, flag := range []string{"--exec-profile", "--exec-episode-project", "--exec-sandbox", "--exec-node", "--exec-max-output-bytes"} {
		if !seen[flag] {
			t.Errorf("QUICKSTART §10's serve line does not show %s", flag)
		}
	}
	// It parses: with a bad --bind appended, serve stops at the bind check.
	var out, errw bytes.Buffer
	if got := run(append(append([]string{}, argv...), "--bind", "not-a-hostport"), &out, &errw); got != exitUsage ||
		!strings.Contains(errw.String(), "is not host:port") {
		t.Fatalf("QUICKSTART §10 `serve %s` does not parse: exit %d, %s", strings.Join(argv[1:], " "), got, errw.String())
	}
}
