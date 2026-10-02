package main

// Row 134 M7 (AC7.1, coding-standards S7): QUICKSTART §9's `serve` and
// `session mint` command lines are bound to the real CLI — every flag the
// runbook names is documented by `serve --help` and accepted by the parser —
// so a renamed or mistyped flag reds here instead of at the operator's
// keyboard in the attended M6 smoke.

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

// section9Commands returns each `/tmp/ailang-worldd <verb...>` command line
// of QUICKSTART §9, backslash continuations joined, as argv after the binary.
func section9Commands(t *testing.T) [][]string {
	t.Helper()
	guide, err := os.ReadFile("../../docs/QUICKSTART.md")
	if err != nil {
		t.Fatal(err)
	}
	_, section, ok := strings.Cut(string(guide), "### 9. Software-engineering tools")
	if !ok {
		t.Fatal("QUICKSTART §9 is absent")
	}
	joined := strings.ReplaceAll(section, "\\\n", " ")
	var cmds [][]string
	for _, line := range strings.Split(joined, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "/tmp/ailang-worldd ") {
			continue
		}
		fields := strings.Fields(strings.TrimSuffix(strings.TrimSpace(line), "&"))
		cmds = append(cmds, fields[1:])
	}
	return cmds
}

var flagToken = regexp.MustCompile(`^--[a-z][a-z-]*$`)

func TestQuickstartSection9FlagsMatchTheCLI(t *testing.T) {
	var help, ignored bytes.Buffer
	if got := run([]string{"serve", "--help"}, &ignored, &help); got != exitUsage {
		t.Fatalf("serve --help exit = %d", got)
	}
	_, serveHelp, ok := strings.Cut(help.String(), "serve flags:")
	if !ok {
		t.Fatalf("serve --help has no serve flags block:\n%s", help.String())
	}
	var serves, mints int
	for _, argv := range section9Commands(t) {
		switch {
		case argv[0] == "serve":
			serves++
			for _, tok := range argv[1:] {
				if flagToken.MatchString(tok) && !regexp.MustCompile(`(?m)^\s+`+regexp.QuoteMeta(tok)+`\b`).MatchString(serveHelp) {
					t.Errorf("QUICKSTART §9 serve flag %s is not documented by serve --help", tok)
				}
			}
			// The parser accepts every documented flag: with a bad --bind
			// appended, serve stops at the bind check, after flag parsing.
			var out, errw bytes.Buffer
			got := run(append(append([]string{}, argv...), "--bind", "not-a-hostport"), &out, &errw)
			if got != exitUsage || !strings.Contains(errw.String(), "is not host:port") {
				t.Errorf("QUICKSTART §9 `serve %s` does not parse: exit %d, %s", strings.Join(argv[1:], " "), got, errw.String())
			}
		case len(argv) > 1 && argv[0] == "session" && argv[1] == "mint":
			mints++
			// Flags and every --grant spec parse; the tty fence (absent here)
			// is the first check after parsing.
			var out, errw bytes.Buffer
			env := sessionEnv{openTerminal: func() (io.ReadWriteCloser, error) { return nil, os.ErrNotExist }, now: func() int64 { return 0 }}
			got := runSessionMint(argv[2:], &out, &errw, env)
			if got != exitUsage || !strings.Contains(errw.String(), "no controlling terminal") {
				t.Errorf("QUICKSTART §9 `session mint %s` does not parse: exit %d, %s", strings.Join(argv[2:], " "), got, errw.String())
			}
			grants := 0
			for _, tok := range argv {
				if tok == "--grant" {
					grants++
				}
			}
			if grants != 6 {
				t.Errorf("QUICKSTART §9 mint names %d grants, want the six effect grants", grants)
			}
		}
	}
	// serve twice (pin-only bootstrap, then with the workspace tools) and one mint.
	if serves != 2 || mints != 1 {
		t.Fatalf("QUICKSTART §9 holds %d serve and %d mint commands, want 2 and 1", serves, mints)
	}
}

// TestSessionMintGrantsExpireWithTheSession is the positive form of the row
// 134 M5b break-1 tripwire: `session mint` used to store each --grant with
// ExpiresAt 0, which the broker's liveness rule (now < ExpiresAt, Unix
// seconds) treats as expired, so a CLI-minted session resolved but saw zero
// tools. Each CLI-minted grant now expires at the session's own expiry
// (mint time + --ttl): live now and through the last second, denied:expired
// at the session expiry and after.
func TestSessionMintGrantsExpireWithTheSession(t *testing.T) {
	dir := t.TempDir()
	db, out := filepath.Join(dir, "world.db"), filepath.Join(dir, "session")
	var stdout, stderr bytes.Buffer
	now := time.Now().Unix()
	env := sessionEnv{openTerminal: func() (io.ReadWriteCloser, error) { return &fakeTerm{r: strings.NewReader("y\n")}, nil }, now: func() int64 { return now }}
	if got := runSessionMint([]string{"--db", db, "--episode", "ep1", "--grant", "Workspace.Read=worktree:5", "--grant", "Ailang.CLI=worktree:1", "--ttl", "3600", "--out", out}, &stdout, &stderr, env); got != exitOK {
		t.Fatalf("mint exit %d: %s", got, stderr.String())
	}
	token, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	st, err := store.Open(db)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = st.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := authority.New(st).ResolveContext(ctx, "Bearer "+string(token), now)
	if err != nil || res.Success == nil || len(res.Success.Caps) != 2 {
		t.Fatalf("resolve = %+v %v", res, err)
	}
	if res.Success.ExpiresAt != now+3600 {
		t.Fatalf("session expiry %d, want %d", res.Success.ExpiresAt, now+3600)
	}
	for _, c := range res.Success.Caps {
		if c.ExpiresAt != res.Success.ExpiresAt {
			t.Fatalf("grant %s expiry %d, want the session expiry %d", c.Effect, c.ExpiresAt, res.Success.ExpiresAt)
		}
	}
	req := broker.Requirement{Effect: "Workspace.Read", Scope: "worktree", Cost: 1}
	for _, at := range []int64{now, now + 3599} {
		if d := broker.Allows(broker.NewCapabilitySnapshot(res.Success.Caps, at), req); !d.Allowed {
			t.Fatalf("CLI-minted grant at t=now+%d: %+v, want allowed", at-now, d)
		}
	}
	for _, at := range []int64{now + 3600, now + 3601} {
		if d := broker.Allows(broker.NewCapabilitySnapshot(res.Success.Caps, at), req); d.Allowed || d.Label != broker.LabelDeniedExpired {
			t.Fatalf("CLI-minted grant at t=now+%d: %+v, want denied:expired", at-now, d)
		}
	}
}

// TestServeExamplesDirDefault (row 134 break-3 fix): an explicit
// --examples-dir wins; otherwise the operator's ~/.ailang/examples is the
// default when it is a directory; otherwise no corpus ("").
func TestServeExamplesDirDefault(t *testing.T) {
	home := t.TempDir()
	homeFn := func() (string, error) { return home, nil }
	if got := resolveExamplesDefault("", homeFn); got != "" {
		t.Fatalf("no ~/.ailang/examples: default = %q, want none", got)
	}
	corpus := filepath.Join(home, ".ailang", "examples")
	if err := os.MkdirAll(corpus, 0o755); err != nil {
		t.Fatal(err)
	}
	if got := resolveExamplesDefault("", homeFn); got != corpus {
		t.Fatalf("default = %q, want %q", got, corpus)
	}
	if got := resolveExamplesDefault("/explicit/corpus", homeFn); got != "/explicit/corpus" {
		t.Fatalf("explicit flag = %q, want it kept", got)
	}
	if got := resolveExamplesDefault("", func() (string, error) { return "", os.ErrNotExist }); got != "" {
		t.Fatalf("no home: default = %q, want none", got)
	}
}
