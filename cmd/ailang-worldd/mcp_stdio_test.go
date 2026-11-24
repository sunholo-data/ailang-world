package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/toolprobe"
)

func TestMCPStdioBridgeSSEJSONSessionAndRedaction(t *testing.T) {
	const token = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	var seenAuth []string
	var seenSID []string
	call := 0
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/mcp/" {
			return &http.Response{StatusCode: http.StatusNotFound, Body: io.NopCloser(strings.NewReader("not found")), Header: http.Header{}}, nil
		}
		seenAuth = append(seenAuth, r.Header.Get("Authorization"))
		seenSID = append(seenSID, r.Header.Get("Mcp-Session-Id"))
		call++
		switch call {
		case 1:
			h := http.Header{"Mcp-Session-Id": {"sid-123"}, "Content-Type": {"text/event-stream"}}
			return &http.Response{StatusCode: http.StatusOK, Header: h, Body: io.NopCloser(strings.NewReader("event: message\ndata: {\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"ok\":true}}\n\n"))}, nil
		case 2:
			return &http.Response{StatusCode: http.StatusAccepted, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(""))}, nil
		case 3:
			h := http.Header{"Content-Type": {"application/json"}}
			return &http.Response{StatusCode: http.StatusOK, Header: h, Body: io.NopCloser(strings.NewReader("{\"jsonrpc\":\"2.0\",\"id\":2,\"result\":{\"tools\":[]}}\n"))}, nil
		default:
			return &http.Response{StatusCode: http.StatusUnauthorized, Header: http.Header{}, Body: io.NopCloser(strings.NewReader("denied token " + token))}, nil
		}
	})}

	input := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize"}`,
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
		`{"jsonrpc":"2.0","id":3,"method":"tools/call"}`,
	}, "\n") + "\n"

	var stdout bytes.Buffer
	if err := bridgeMCPStdio(strings.NewReader(input), &stdout, client, "http://127.0.0.1:7644/mcp/", token); err != nil {
		t.Fatalf("bridgeMCPStdio: %v", err)
	}
	if strings.Contains(stdout.String(), token) {
		t.Fatalf("token leaked in output: %s", stdout.String())
	}
	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3: %q", len(lines), stdout.String())
	}
	if lines[0] != `{"jsonrpc":"2.0","id":1,"result":{"ok":true}}` {
		t.Errorf("line 0 = %s", lines[0])
	}
	if lines[1] != `{"jsonrpc":"2.0","id":2,"result":{"tools":[]}}` {
		t.Errorf("line 1 = %s", lines[1])
	}
	var errResp struct {
		ID    int `json:"id"`
		Error struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(lines[2]), &errResp); err != nil {
		t.Fatalf("line 2 invalid JSON: %v (%s)", err, lines[2])
	}
	if errResp.ID != 3 || !strings.Contains(errResp.Error.Message, "HTTP 401") || !strings.Contains(errResp.Error.Message, "<redacted>") {
		t.Errorf("unexpected error response: %+v", errResp)
	}
	for i, auth := range seenAuth {
		if auth != "Bearer "+token {
			t.Errorf("call %d auth = %q", i+1, auth)
		}
	}
	if seenSID[0] != "" || seenSID[1] != "sid-123" || seenSID[2] != "sid-123" {
		t.Errorf("Mcp-Session-Id propagation = %#v", seenSID)
	}
}

func TestAgyHostInitLauncherGateAndStage(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root := t.TempDir()
	dir := filepath.Join(root, "hosts")
	if code, out := hostsInitTest(t, dir); code != 0 {
		t.Fatalf("%d: %s", code, out)
	}
	// Verify agy-world-gate allows World MCP tools and denies built-ins.
	gate := filepath.Join(dir, "bin", "agy-world-gate")
	for _, tc := range []struct {
		tool string
		want string
	}{
		{"mcp_world_ailang_read", `"decision":"allow"`},
		{"mcp__world__workspace-exec", `"decision":"allow"`},
		{"update_steps", `"decision":"allow"`},
		{"run_command", `"decision":"deny"`},
		{"write_file", `"decision":"deny"`},
		{"read_file", `"decision":"deny"`},
	} {
		wrapper := filepath.Join(t.TempDir(), "run-gate")
		payload := `{"toolCall":{"name":"` + tc.tool + `"}}`
		if err := os.WriteFile(wrapper, []byte("#!/bin/sh\nprintf '%s' "+shellQuote(payload)+" | "+shellQuote(gate)+"\n"), 0o700); err != nil {
			t.Fatal(err)
		}
		ctx, cancel := budgetContext(5 * time.Second)
		out, err := toolprobe.OutputEnv(ctx, wrapper, "", []string{"PATH=" + os.Getenv("PATH")})
		cancel()
		if err != nil {
			t.Fatalf("%s: %v", tc.tool, err)
		}
		if !strings.Contains(string(out), tc.want) {
			t.Errorf("gate(%s) = %s, want %s", tc.tool, out, tc.want)
		}
	}
	// Verify agy-world mounts .agents/plugins.json when started in an empty dir.
	fakeBin := filepath.Join(root, "bin")
	if err := os.MkdirAll(fakeBin, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fakeBin, "agy"), []byte("#!/bin/sh\nprintf '%s\\n' \"$@\"\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	emptyDir := filepath.Join(root, "empty-launch")
	if err := os.MkdirAll(emptyDir, 0o700); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := budgetContext(5 * time.Second)
	out, err := toolprobe.OutputEnv(ctx, filepath.Join(dir, "bin", "agy-world"), emptyDir,
		[]string{"PATH=" + fakeBin + ":/usr/bin:/bin", "WORLD_SESSION=test-session"}, "--model", "gemini")
	cancel()
	if err != nil {
		t.Fatalf("agy-world: %v", err)
	}
	if !strings.Contains(string(out), filepath.Join(dir, "agy-world.settings.json")) || !strings.Contains(string(out), "--model\ngemini") {
		t.Fatalf("unexpected agy-world argv: %s", out)
	}
	pluginsJSON, err := os.ReadFile(filepath.Join(emptyDir, ".agents", "plugins.json"))
	if err != nil || !strings.Contains(string(pluginsJSON), `"agy-plugin"`) {
		t.Fatalf("plugins.json: %v %s", err, pluginsJSON)
	}
	// Verify headless staging writes no secrets and valid plugin files.
	headlessDir := filepath.Join(root, "headless-launch")
	if err := os.MkdirAll(headlessDir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := stageAgyHeadlessPlugin(headlessDir, "127.0.0.1:7644", "Episode test context"); err != nil {
		t.Fatal(err)
	}
	mcpBytes, err := os.ReadFile(filepath.Join(headlessDir, ".agents", "plugins", "ailang-world", "mcp_config.json"))
	if err != nil || !strings.Contains(string(mcpBytes), `"mcp-stdio"`) || !strings.Contains(string(mcpBytes), "http://127.0.0.1:7644") {
		t.Fatalf("headless mcp_config.json: %v %s", err, mcpBytes)
	}
}
