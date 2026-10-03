package main

// Row 134 M7 (AC7.1, coding-standards S7): QUICKSTART §9's `serve` and
// `session new` (row 138 M5; `session mint` before it) command lines are
// bound to the real CLI — every flag the
// runbook names is documented by `serve --help` and accepted by the parser —
// so a renamed or mistyped flag reds here instead of at the operator's
// keyboard in the attended M6 smoke.

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
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
	if got := run([]string{"serve", "--help"}, &help, &ignored); got != exitOK {
		t.Fatalf("serve --help exit = %d", got)
	}
	_, serveHelp, ok := strings.Cut(help.String(), "serve flags:")
	if !ok {
		t.Fatalf("serve --help has no serve flags block:\n%s", help.String())
	}
	var serves, mints, commits int
	for i, argv := range section9Commands(t) {
		switch {
		case len(argv) > 2 && argv[0] == "--addr" && argv[2] == "commit":
			commits++
			// M6 measured POST /v1/commit as session-gated: the genesis commit
			// runs against the live daemon, after mint, with --session.
			if mints != 1 || serves != 2 {
				t.Errorf("QUICKSTART §9 commits (command %d) before the mint and the tools serve", i)
			}
			assertSection9CommitLine(t, argv)
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
		case len(argv) > 1 && argv[0] == "session" && argv[1] == "new":
			mints++
			// Row 138 M5: the runbook provisions the session with `session new`.
			// Its paths are swapped for a fixture (an existing store, a
			// workspace root holding an ep1 worktree), so validation passes and
			// the tty fence (absent here) is the first refusal: every flag and
			// the preset parse.
			fx := t.TempDir()
			db := filepath.Join(fx, "world.db")
			st, err := store.Open(db)
			if err != nil {
				t.Fatal(err)
			}
			_ = st.Close()
			ws := filepath.Join(fx, "ws")
			_ = os.MkdirAll(filepath.Join(ws, "ep1"), 0o700)
			_ = os.WriteFile(filepath.Join(ws, "ep1", ".git"), []byte("gitdir: /x/.git/worktrees/ep1\n"), 0o644)
			args := make([]string, 0, len(argv))
			for _, tok := range argv[2:] {
				switch tok {
				case "/tmp/se-world/world.db":
					tok = db
				case "/tmp/se-ws":
					tok = ws
				case "/tmp/se-session":
					tok = filepath.Join(fx, "session")
				}
				args = append(args, tok)
			}
			var out, errw bytes.Buffer
			env := sessionEnv{openTerminal: func() (io.ReadWriteCloser, error) { return nil, os.ErrNotExist }, now: func() int64 { return 0 }}
			got := runSessionNew(args, &out, &errw, env)
			if got != exitUsage || !strings.Contains(errw.String(), "no controlling terminal") {
				t.Errorf("QUICKSTART §9 `session new %s` does not parse: exit %d, %s", strings.Join(argv[2:], " "), got, errw.String())
			}
			if strings.Join(argv, " ") != "session new ep1 --db /tmp/se-world/world.db --workspace-root /tmp/se-ws --preset se-tools --ttl 14400 --out /tmp/se-session" {
				t.Errorf("QUICKSTART §9 session line changed: %q (the eight se-tools grants come from --preset se-tools)", strings.Join(argv, " "))
			}
		}
	}
	// serve twice (pin-only bootstrap, then with the workspace tools), one
	// session new, then the one session-gated genesis commit.
	if serves != 2 || mints != 1 || commits != 1 {
		t.Fatalf("QUICKSTART §9 holds %d serve, %d session new and %d commit commands, want 2, 1 and 1", serves, mints, commits)
	}
}

// assertSection9CommitLine runs §9's commit line, with its file and its
// `$(cat /tmp/se-session)` substitution replaced by fixtures and --addr by a
// recording server: it must POST the file to /v1/commit as the Bearer session.
func assertSection9CommitLine(t *testing.T, argv []string) {
	t.Helper()
	want := []string{"--addr", "http://127.0.0.1:7644", "commit", "--file", "/tmp/se-genesis.json", "--session", `"$(cat`, `/tmp/se-session)"`}
	if strings.Join(argv, " ") != strings.Join(want, " ") {
		t.Fatalf("QUICKSTART §9 commit line = %q, want %q", argv, want)
	}
	var method, path, auth, body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		method, path, auth, body = r.Method, r.URL.Path, r.Header.Get("Authorization"), string(data)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"ok":true}`)
	}))
	defer srv.Close()
	file := filepath.Join(t.TempDir(), "se-genesis.json")
	if err := os.WriteFile(file, []byte(`{"observedHead":""}`), 0o600); err != nil {
		t.Fatal(err)
	}
	token := strings.Repeat("ab", 32)
	var out, errw bytes.Buffer
	if got := run([]string{"--addr", srv.URL, "commit", "--file", file, "--session", token}, &out, &errw); got != exitOK {
		t.Fatalf("§9 commit line exit %d: %s", got, errw.String())
	}
	if method != http.MethodPost || path != "/v1/commit" || auth != "Bearer "+token || body != `{"observedHead":""}` {
		t.Fatalf("§9 commit sent %s %s auth=%q body=%q", method, path, auth, body)
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
	db, out := dbPathOf(t), filepath.Join(dir, "session")
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

// TestServeRunCapsFlags (row 135 AC3.1, CLI half): serve --help documents the
// three --run-* flags, the parser accepts them (stopping at a later bad
// --bind), and an empty name in --run-allow-caps is a usage error. The
// semantic refusals (unknown names, bare or non-loopback hosts) are the
// daemon's startup refusals (host/daemon TestRunCapsStartupRefusals).
func TestServeRunCapsFlags(t *testing.T) {
	var help, ignored bytes.Buffer
	run([]string{"serve", "--help"}, &help, &ignored)
	for _, flag := range []string{"--run-allow-caps", "--run-net-allow", "--run-net-allow-http"} {
		if !regexp.MustCompile(`(?m)^\s+` + regexp.QuoteMeta(flag) + `\b`).MatchString(help.String()) {
			t.Errorf("serve --help does not document %s", flag)
		}
	}
	var out, errw bytes.Buffer
	if got := run([]string{"serve", "--db", "/x/world.db", "--run-allow-caps", "Env,Net", "--run-net-allow", "127.0.0.1:7655",
		"--run-net-allow", "[::1]:7655", "--run-net-allow-http", "--bind", "not-a-hostport"}, &out, &errw); got != exitUsage ||
		!strings.Contains(errw.String(), "is not host:port") {
		t.Fatalf("serve with the --run-* flags: exit %d, %s", got, errw.String())
	}
	errw.Reset()
	if got := run([]string{"serve", "--db", "/x/world.db", "--run-allow-caps", "Env,,Net"}, &out, &errw); got != exitUsage ||
		!strings.Contains(errw.String(), "empty capability name") {
		t.Fatalf("serve --run-allow-caps Env,,Net: exit %d, %s", got, errw.String())
	}
}
