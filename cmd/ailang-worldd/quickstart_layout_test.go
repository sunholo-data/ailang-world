package main

// Row 141 M3 (w-workspace-project-layouts AC3.1, S7): QUICKSTART §11 and the
// three project-layout flags are bound to `serve --help` and to the Config
// serve hands the daemon.

import (
	"bytes"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/daemon"
)

var layoutFlags = []string{"--workspace-module-root", "--workspace-episode-module-root", "--workspace-package-cache"}

// section11ServeLines returns QUICKSTART §11's `/tmp/ailang-worldd serve`
// lines, continuations joined, as argv after the binary.
func section11ServeLines(t *testing.T) [][]string {
	t.Helper()
	guide, err := os.ReadFile("../../docs/QUICKSTART.md")
	if err != nil {
		t.Fatal(err)
	}
	_, section, ok := strings.Cut(string(guide), "### 11. Projects with a module root and registry packages")
	if !ok {
		t.Fatal("QUICKSTART §11 is absent")
	}
	var cmds [][]string
	for _, line := range strings.Split(strings.ReplaceAll(section, "\\\n", " "), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "/tmp/ailang-worldd serve ") {
			cmds = append(cmds, strings.Fields(strings.TrimSuffix(line, "&"))[1:])
		}
	}
	return cmds
}

func TestQuickstartSection11FlagsMatchTheCLI(t *testing.T) {
	var help, ignored bytes.Buffer
	if got := run([]string{"serve", "--help"}, &help, &ignored); got != exitOK {
		t.Fatalf("serve --help exit = %d", got)
	}
	_, serveHelp, ok := strings.Cut(help.String(), "serve flags:")
	if !ok {
		t.Fatalf("serve --help has no serve flags block:\n%s", help.String())
	}
	documented := func(flag string) bool {
		return regexp.MustCompile(`(?m)^\s+` + regexp.QuoteMeta(flag) + `\b`).MatchString(serveHelp)
	}
	for _, flag := range layoutFlags {
		if !documented(flag) {
			t.Errorf("serve --help does not document %s", flag)
		}
	}
	lines := section11ServeLines(t)
	if len(lines) != 1 {
		t.Fatalf("QUICKSTART §11 holds %d serve lines, want 1", len(lines))
	}
	argv := lines[0]
	seen := map[string]bool{}
	for _, tok := range argv[1:] {
		if flagToken.MatchString(tok) {
			seen[tok] = true
			if !documented(tok) {
				t.Errorf("QUICKSTART §11 serve flag %s is not documented by serve --help", tok)
			}
		}
	}
	for _, flag := range layoutFlags {
		if !seen[flag] {
			t.Errorf("QUICKSTART §11's serve line does not show %s", flag)
		}
	}
	// It parses: with a bad --bind appended, serve stops at the bind check.
	var out, errw bytes.Buffer
	if got := run(append(append([]string{}, argv...), "--bind", "not-a-hostport"), &out, &errw); got != exitUsage ||
		!strings.Contains(errw.String(), "is not host:port") {
		t.Fatalf("QUICKSTART §11 `serve %s` does not parse: exit %d, %s", strings.Join(argv[1:], " "), got, errw.String())
	}
	checkLayoutWiring(t)
}

// The wiring: the three flags land in the Config serve hands the daemon, and
// unset they leave it zero (so the daemon's behaviour is unchanged).
func checkLayoutWiring(t *testing.T) {
	t.Helper()
	got := captureServeConfig(t)
	var out, errw bytes.Buffer
	code := run([]string{"serve", "--db", "/x/world.db", "--workspace-root", "/w",
		"--workspace-module-root", "tools",
		"--workspace-episode-module-root", "a=.", "--workspace-episode-module-root", "b=x",
		"--workspace-package-cache", "/p"}, &out, &errw)
	if code != exitOK {
		t.Fatalf("serve exit %d: %s", code, errw.String())
	}
	if got.WorkspaceModuleRoot != "tools" || !reflect.DeepEqual(got.WorkspaceEpisodeModuleRoots, []string{"a=.", "b=x"}) ||
		got.WorkspacePackageCache != "/p" {
		t.Fatalf("serve built Config %+v", *got)
	}
	*got = daemon.Config{}
	if code := run([]string{"serve", "--db", "/x/world.db"}, &out, &errw); code != exitOK {
		t.Fatalf("serve exit %d", code)
	}
	if got.WorkspaceModuleRoot != "" || got.WorkspaceEpisodeModuleRoots != nil || got.WorkspacePackageCache != "" {
		t.Fatalf("serve with no layout flag built %+v", *got)
	}
	// An empty --workspace-module-root is a usage error (it would read as
	// "unset"): it stops at flag parsing, before the daemon.
	errw.Reset()
	if code := run([]string{"serve", "--db", "/x/world.db", "--workspace-module-root", ""}, &out, &errw); code != exitUsage ||
		!strings.Contains(errw.String(), "is empty (use . for the worktree root)") {
		t.Fatalf("--workspace-module-root \"\": exit %d, %s", code, errw.String())
	}
}
