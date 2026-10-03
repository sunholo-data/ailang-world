package main

// Row 138 M1 (w-worldd-developer-cli §5 AC1.1, AC1.2, AC1.4, AC1.5) against
// fake /mcp/ servers that reproduce the three measured wire shapes (V9–V11)
// byte for byte. The real-daemon arms are in cliwalk_e2e_test.go.

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// fakeMCP serves /mcp/ with a fixed status, Content-Type and body, and
// /v1/head from heads (one answer per GET, the last repeated; "" = 404).
type fakeMCP struct {
	status int
	ctype  string
	body   string
	heads  []string

	mu       sync.Mutex
	headGets int
	auth     []string
	bodies   []string
}

func (f *fakeMCP) serve(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/mcp/":
			b, _ := io.ReadAll(r.Body)
			f.auth = append(f.auth, r.Header.Get("Authorization"))
			f.bodies = append(f.bodies, string(b))
			if r.Header.Get("Accept") != "application/json, text/event-stream" {
				http.Error(w, "bad accept", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", f.ctype)
			w.WriteHeader(f.status)
			_, _ = io.WriteString(w, f.body)
		case r.URL.Path == "/v1/head":
			h := ""
			if len(f.heads) > 0 {
				i := f.headGets
				if i >= len(f.heads) {
					i = len(f.heads) - 1
				}
				h = f.heads[i]
			}
			f.headGets++
			if h == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				_, _ = io.WriteString(w, `{"error":{"class":"NotFound","message":"no world head has been selected yet"}}`)
				return
			}
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, _ = io.WriteString(w, h)
		case strings.HasPrefix(r.URL.Path, "/v1/worlds/"):
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"ref":"x","revision":7,"stateRoot":"s","logHead":"l"}`)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

const fakeToken = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

// sessionFile writes the token to a 0600 file and returns its path.
func sessionFile(t *testing.T, token string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "session")
	if err := os.WriteFile(p, []byte(token+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func sseFrame(data string) string { return "event: message\ndata: " + data + "\n\n" }

// The measured success result (V9): the committed output is the text content.
const fakeOutput = `{"content":"hello from ep1\n","ok":true,"world":{"effects":[{"id":"e1","record":"sha256:1111111111111111111111111111111111111111111111111111111111111111","status":"ok"}],"plan":"sha256:2222222222222222222222222222222222222222222222222222222222222222"}}`

func fakeResult(id int, text string) string {
	tb, _ := json.Marshal(text)
	return `{"jsonrpc":"2.0","id":` + itoa(id) + `,"result":{"content":[{"type":"text","text":` + string(tb) + `}],"structuredContent":` + text + `}}`
}

func itoa(i int) string { b, _ := json.Marshal(i); return string(b) }

// TestMCPWireShapes is AC1.1: the httptest table over every wire shape —
// SSE single-line, multi-line and id-order; the 200 JSON -32603 envelope;
// the four 401 bodies with their hints; a 400. It kills MUT-SSE-FIRSTLINE
// (the multi-line row) and MUT-ENVELOPE-AS-RESULT (the envelope row).
func TestMCPWireShapes(t *testing.T) {
	result := fakeResult(1, fakeOutput)
	// The same JSON-RPC response split over several data: lines (SSE spec:
	// joined with "\n", which JSON treats as whitespace).
	split := strings.Replace(strings.Replace(result, `"id":1,`, "\"id\":1,\ndata: ", 1), `"result":`, "\ndata: \"result\":", 1)
	cases := []struct {
		name       string
		status     int
		ctype      string
		body       string
		kind       mcpKind
		wantCode   int
		wantStdout string // substring
		wantStderr string // substring
	}{
		{"sse single line", 200, "text/event-stream", sseFrame(result), mcpResult, exitOK, "plan       sha256:2222", ""},
		{"sse multi-line data", 200, "text/event-stream", "event: message\ndata: " + split + "\n\n", mcpResult, exitOK, "effect     e1 ok sha256:1111", ""},
		{"sse id order: the matching id wins", 200, "text/event-stream",
			sseFrame(`{"jsonrpc":"2.0","id":9,"error":{"code":-32601,"message":"not mine"}}`) + ": keepalive\n\n" + sseFrame(result),
			mcpResult, exitOK, "content:     hello from ep1\\n", ""},
		{"sse json-rpc error is a tool error", 200, "text/event-stream",
			sseFrame(`{"jsonrpc":"2.0","id":1,"error":{"code":-32602,"message":"unknown tool \"nope\""}}`), mcpToolError, exitUsage, "", "tool error -32602"},
		{"200 json -32603 envelope is a host failure, never a result", 200, "application/json",
			`{"jsonrpc":"2.0","id":1,"error":{"code":-32603,"message":"host callback failed"}}` + "\n", mcpHostFailure, exitUsage, "", "host failure -32603: host callback failed"},
		{"401 absent", 401, "text/plain; charset=utf-8", "session credential is absent: send Authorization: Bearer <session-credential>\n", mcpDenied, exitUsage, "", "fix: pass --session"},
		{"401 unknown", 401, "text/plain; charset=utf-8", "unknown session credential\n", mcpDenied, exitUsage, "", "does not know that credential"},
		{"401 malformed", 401, "text/plain; charset=utf-8", "malformed Authorization header: expected Bearer <64-hex-credential>\n", mcpDenied, exitUsage, "", "exactly one 64-hex token"},
		{"401 expired", 401, "text/plain; charset=utf-8", "session credential has expired\n", mcpDenied, exitUsage, "", "mint a new one"},
		{"400 transport is verbatim", 400, "application/json",
			`{"jsonrpc":"2.0","id":1,"error":{"code":-32600,"message":"Accept must list both application/json and text/event-stream"}}`,
			mcpTransport, exitUsage, "", "HTTP 400: {\"jsonrpc\""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := classifyMCP(c.status, c.ctype, []byte(c.body), "1"); got.Kind != c.kind {
				t.Fatalf("classifyMCP kind = %d, want %d (%+v)", got.Kind, c.kind, got)
			}
			f := &fakeMCP{status: c.status, ctype: c.ctype, body: c.body, heads: []string{"sha256:aa"}}
			srv := f.serve(t)
			code, stdout, stderr := runCLI(t, srv.URL, "call", "ailang-read", "--session", sessionFile(t, fakeToken), "--arg", "path=data.txt")
			if code != c.wantCode || !strings.Contains(stdout, c.wantStdout) || !strings.Contains(stderr, c.wantStderr) {
				t.Fatalf("call exit %d stdout %q stderr %q; want exit %d, stdout ∋ %q, stderr ∋ %q",
					code, stdout, stderr, c.wantCode, c.wantStdout, c.wantStderr)
			}
			if c.kind != mcpResult && stdout != "" {
				t.Fatalf("a non-result printed to stdout: %q", stdout)
			}
			if len(f.auth) != 1 || f.auth[0] != "Bearer "+fakeToken {
				t.Fatalf("the call did not present the session as a Bearer header: %v", len(f.auth))
			}
			if !strings.Contains(f.bodies[0], `"method":"tools/call"`) || !strings.Contains(f.bodies[0], `"arguments":{"path":"data.txt"}`) {
				t.Fatalf("request body = %s", f.bodies[0])
			}
		})
	}
}

// TestCallJSONOutIsByteExact is AC1.2: --json-out prints the committed output
// text exactly, plus one newline, so sha256(stdout − "\n") is the output ref.
func TestCallJSONOutIsByteExact(t *testing.T) {
	text := `{"content":"tab\there — ünïcode  and  spaces\n","ok":true,"world":{"effects":[],"plan":"sha256:2222222222222222222222222222222222222222222222222222222222222222"}}`
	f := &fakeMCP{status: 200, ctype: "text/event-stream", body: sseFrame(fakeResult(1, text)), heads: []string{"sha256:aa"}}
	srv := f.serve(t)
	code, stdout, stderr := runCLI(t, srv.URL, "call", "ailang-read", "--session", sessionFile(t, fakeToken), "--json", `{"path":"x"}`, "--json-out")
	if code != exitOK {
		t.Fatalf("exit %d: %s", code, stderr)
	}
	if stdout != text+"\n" {
		t.Fatalf("--json-out = %q, want the output bytes %q plus one newline", stdout, text)
	}
	got := sha256.Sum256([]byte(strings.TrimSuffix(stdout, "\n")))
	want := sha256.Sum256([]byte(text))
	if got != want {
		t.Fatalf("sha256(stdout - newline) = %s, want %s", hex.EncodeToString(got[:]), hex.EncodeToString(want[:]))
	}
}

// TestCallStrictRefusal is D-CLI-3: a committed refusal (ok:false) is exit 0,
// and exit 3 under --strict.
func TestCallStrictRefusal(t *testing.T) {
	text := `{"ok":false,"refused":"unknown argument limit","world":{"effects":[],"plan":"sha256:2222222222222222222222222222222222222222222222222222222222222222"}}`
	f := &fakeMCP{status: 200, ctype: "text/event-stream", body: sseFrame(fakeResult(1, text)), heads: []string{"sha256:aa"}}
	srv := f.serve(t)
	sess := sessionFile(t, fakeToken)
	if code, out, errOut := runCLI(t, srv.URL, "call", "ailang-read", "--session", sess); code != exitOK || !strings.Contains(out, "refused:") {
		t.Fatalf("committed refusal: exit %d %q %q, want 0", code, out, errOut)
	}
	if code, _, errOut := runCLI(t, srv.URL, "call", "ailang-read", "--session", sess, "--strict"); code != exitIntegrity {
		t.Fatalf("--strict committed refusal: exit %d %q, want %d", code, errOut, exitIntegrity)
	}
}

// TestCallHostFailureProbe is AC1.4's fake arm and the "a commit landed"
// branch: the -32603 envelope is disambiguated by the head before and after.
func TestCallHostFailureProbe(t *testing.T) {
	envelope := `{"jsonrpc":"2.0","id":1,"error":{"code":-32603,"message":"host callback failed"}}`
	cases := []struct {
		name  string
		heads []string
		want  string
	}{
		{"no head", []string{""}, "no world head"},
		{"head moved: a commit landed", []string{"sha256:aa", "sha256:bb"}, "a commit landed (entry 7) — do not retry; inspect it with: ailang-worldd why 7"},
		{"head unchanged: passthrough", []string{"sha256:aa"}, "host failure -32603: host callback failed\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := &fakeMCP{status: 200, ctype: "application/json", body: envelope, heads: c.heads}
			srv := f.serve(t)
			code, stdout, stderr := runCLI(t, srv.URL, "call", "ailang-read", "--session", sessionFile(t, fakeToken))
			if code != exitUsage || stdout != "" || !strings.Contains(stderr, c.want) {
				t.Fatalf("exit %d stdout %q stderr %q, want exit 1 and %q", code, stdout, stderr, c.want)
			}
			if c.name == "head unchanged: passthrough" && (strings.Contains(stderr, "landed") || strings.Contains(stderr, "no world head")) {
				t.Fatalf("an unchanged head was misread: %q", stderr)
			}
		})
	}
}

// TestSessionResolver covers §3.1's precedence and file rules.
func TestSessionResolver(t *testing.T) {
	dir := t.TempDir()
	good := sessionFile(t, fakeToken)
	loose := filepath.Join(dir, "loose")
	if err := os.WriteFile(loose, []byte(fakeToken), 0o644); err != nil {
		t.Fatal(err)
	}
	big := filepath.Join(dir, "big")
	if err := os.WriteFile(big, bytes.Repeat([]byte("a"), 300), 0o600); err != nil {
		t.Fatal(err)
	}
	junk := filepath.Join(dir, "junk")
	if err := os.WriteFile(junk, []byte("not-a-token"), 0o600); err != nil {
		t.Fatal(err)
	}
	other := strings.Repeat("f", 64)

	cases := []struct {
		name, flag, env string
		wantTok         string
		wantErr, warn   string
	}{
		{"flag file", good, "", fakeToken, "", ""},
		{"flag raw token warns", fakeToken, "", fakeToken, "", "visible to other local processes"},
		{"env token, no warning", "", other, other, "", ""},
		{"env file", "", good, fakeToken, "", ""},
		{"flag beats env", good, other, fakeToken, "", ""},
		{"neither", "", "", "", "no session credential", ""},
		{"loose file warns", loose, "", fakeToken, "", "readable by group/other"},
		{"oversized file", big, "", "", "at most 256 bytes", ""},
		{"junk file", junk, "", "", "does not hold one 64-hex session token", ""},
		{"missing file", filepath.Join(dir, "absent"), "", "", "neither a 64-hex session token nor a readable file", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv(sessionEnvVar, c.env)
			var stderr bytes.Buffer
			tok, err := resolveSession(c.flag, &stderr)
			if c.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("err = %v, want %q", err, c.wantErr)
				}
			} else if err != nil || tok != c.wantTok {
				t.Fatalf("token ok=%t err=%v", tok == c.wantTok, err)
			}
			if c.warn != "" && !strings.Contains(stderr.String(), c.warn) {
				t.Fatalf("stderr %q, want warning %q", stderr.String(), c.warn)
			}
			if c.warn == "" && stderr.Len() != 0 {
				t.Fatalf("unexpected warning %q", stderr.String())
			}
			for _, s := range []string{stderr.String(), errString(err)} {
				if strings.Contains(s, fakeToken) || strings.Contains(s, other) {
					t.Fatal("the resolver echoed the token")
				}
			}
		})
	}
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// TestTokenNeverPrinted is AC1.5 (kills MUT-TOKEN-ECHO): across every call
// and tools outcome — result, tool error, host failure, all four denials —
// and commit --session <file>, grep stdout and stderr for the token: 0
// matches. The token reaches the daemon only as the Bearer header.
func TestTokenNeverPrinted(t *testing.T) {
	bodies := []fakeMCP{
		{status: 200, ctype: "text/event-stream", body: sseFrame(fakeResult(1, fakeOutput))},
		{status: 200, ctype: "text/event-stream", body: sseFrame(`{"jsonrpc":"2.0","id":1,"result":{"tools":[{"name":"ailang-read","description":"Read a file","inputSchema":{"required":["path"]}}]}}`)},
		{status: 200, ctype: "text/event-stream", body: sseFrame(`{"jsonrpc":"2.0","id":1,"error":{"code":-32602,"message":"unknown tool"}}`)},
		{status: 200, ctype: "application/json", body: `{"jsonrpc":"2.0","id":1,"error":{"code":-32603,"message":"host callback failed"}}`},
		{status: 401, ctype: "text/plain", body: "unknown session credential"},
		{status: 401, ctype: "text/plain", body: "session credential has expired"},
	}
	file := sessionFile(t, fakeToken)
	matches := 0
	for i := range bodies {
		f := &bodies[i]
		f.heads = []string{"sha256:aa"}
		srv := f.serve(t)
		for _, sess := range []string{file, fakeToken} {
			for _, argv := range [][]string{
				{"call", "ailang-read", "--session", sess},
				{"call", "ailang-read", "--session", sess, "--json-out"},
				{"tools", "list", "--session", sess},
				{"tools", "list", "--session", sess, "--json"},
			} {
				_, stdout, stderr := runCLI(t, srv.URL, argv...)
				matches += strings.Count(stdout+stderr, fakeToken)
			}
			t.Setenv(sessionEnvVar, sess)
			_, stdout, stderr := runCLI(t, srv.URL, "call", "ailang-read")
			matches += strings.Count(stdout+stderr, fakeToken)
			t.Setenv(sessionEnvVar, "")
		}
		for _, a := range f.auth {
			if a != "Bearer "+fakeToken {
				t.Fatalf("a request carried Authorization other than the session (%d bytes)", len(a))
			}
		}
	}
	// commit --session <file>: the token is the Bearer header, never output.
	var auth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"selectedHead":"sha256:aa"}`)
	}))
	defer srv.Close()
	commit := filepath.Join(t.TempDir(), "c.json")
	if err := os.WriteFile(commit, []byte(`{"observedHead":""}`), 0o600); err != nil {
		t.Fatal(err)
	}
	code, stdout, stderr := runCLI(t, srv.URL, "commit", "--file", commit, "--session", file)
	if code != exitOK || auth != "Bearer "+fakeToken {
		t.Fatalf("commit --session <file>: exit %d, bearer ok=%t, %s", code, auth == "Bearer "+fakeToken, stderr)
	}
	matches += strings.Count(stdout+stderr, fakeToken)
	t.Setenv(sessionEnvVar, file)
	auth = ""
	if code, stdout, stderr = runCLI(t, srv.URL, "commit", "--file", commit); code != exitOK || auth != "Bearer "+fakeToken {
		t.Fatalf("commit with WORLD_SESSION: exit %d, bearer ok=%t", code, auth == "Bearer "+fakeToken)
	}
	matches += strings.Count(stdout+stderr, fakeToken)
	if matches != 0 {
		t.Fatalf("the session token appeared %d time(s) in CLI output", matches)
	}
}

// TestToolsList renders name, description and required arguments.
func TestToolsList(t *testing.T) {
	list := `{"jsonrpc":"2.0","id":1,"result":{"tools":[{"name":"ailang-read","description":"Read a file in the worktree","inputSchema":{"type":"object","required":["path"]}},{"name":"ailang-cli","description":"Run one CLI op","inputSchema":{"type":"object"}}]}}`
	f := &fakeMCP{status: 200, ctype: "text/event-stream", body: sseFrame(list)}
	srv := f.serve(t)
	sess := sessionFile(t, fakeToken)
	code, stdout, stderr := runCLI(t, srv.URL, "tools", "list", "--session", sess)
	if code != exitOK || !strings.Contains(stdout, "ailang-read") || !strings.Contains(stdout, "required: path") ||
		!strings.Contains(stdout, "required: (none)") || !strings.Contains(stdout, "2 tool(s)") {
		t.Fatalf("tools list exit %d stdout %q stderr %q", code, stdout, stderr)
	}
	code, stdout, _ = runCLI(t, srv.URL, "tools", "list", "--session", sess, "--json")
	var tools []map[string]any
	if code != exitOK || json.Unmarshal([]byte(stdout), &tools) != nil || len(tools) != 2 {
		t.Fatalf("tools list --json exit %d %q", code, stdout)
	}
	if !strings.Contains(f.bodies[0], `"method":"tools/list"`) {
		t.Fatalf("request = %s", f.bodies[0])
	}
}

// TestNewVerbsHelp: `<verb> --help` exits 0 with the verb's help on stdout,
// and the top-level usage names every new verb.
func TestNewVerbsHelp(t *testing.T) {
	for _, argv := range [][]string{
		{"tools", "--help"}, {"tools", "list", "--help"}, {"call", "--help"}, {"call", "ailang-read", "--help"},
		{"why", "--help"}, {"log", "tail", "--help"}, {"provenance", "--help"},
	} {
		var stdout, stderr bytes.Buffer
		if code := run(argv, &stdout, &stderr); code != exitOK || !strings.HasPrefix(stdout.String(), "usage: ailang-worldd") {
			t.Errorf("%v: exit %d stdout %q stderr %q", argv, code, stdout.String(), stderr.String())
		}
	}
	var stdout, stderr bytes.Buffer
	if code := run([]string{"help"}, &stdout, &stderr); code != exitOK {
		t.Fatal(code)
	}
	for _, verb := range []string{" tools list ", " call <tool> ", " why <index", " log tail ", " provenance ", "commit --file <commit.json> [--session", " session mint ", " session revoke "} {
		if !strings.Contains(stdout.String(), verb) {
			t.Errorf("usage does not list %q", verb)
		}
	}
}
