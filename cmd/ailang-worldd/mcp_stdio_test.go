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

	readyFile := filepath.Join(t.TempDir(), ".world-mcp-ready")
	t.Setenv("WORLD_AGY_MCP_READY", readyFile)
	var stdout bytes.Buffer
	if err := bridgeMCPStdio(strings.NewReader(input), &stdout, client, "http://127.0.0.1:7644/mcp/", token); err != nil {
		t.Fatalf("bridgeMCPStdio: %v", err)
	}
	if _, err := os.Stat(readyFile); err != nil {
		t.Fatalf("WORLD_AGY_MCP_READY was not written after tools/list: %v", err)
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
	// Verify agy-world-gate allows World MCP tools (including call_mcp_tool targeting
	// ailang-world_world / world) and denies built-ins and non-World MCP servers.
	gate := filepath.Join(dir, "bin", "agy-world-gate")
	for _, tc := range []struct {
		label   string
		payload string
		want    string
	}{
		{"call_mcp_tool world plugin", `{"toolCall":{"name":"call_mcp_tool","args":{"ServerName":"ailang-world_world","ToolName":"ailang-read"}}}`, `"decision":"allow"`},
		{"call_mcp_tool world", `{"toolCall":{"name":"call_mcp_tool","args":{"ServerName":"world","ToolName":"ailang-read"}}}`, `"decision":"allow"`},
		{"call_mcp_tool foreign server", `{"toolCall":{"name":"call_mcp_tool","args":{"ServerName":"other-mcp","ToolName":"shell"}}}`, `"decision":"deny"`},
		{"mcp_world_ailang_read", `{"toolCall":{"name":"mcp_world_ailang_read"}}`, `"decision":"allow"`},
		{"mcp__world__workspace-exec", `{"toolCall":{"name":"mcp__world__workspace-exec"}}`, `"decision":"allow"`},
		{"update_steps", `{"toolCall":{"name":"update_steps"}}`, `"decision":"allow"`},
		{"run_command", `{"toolCall":{"name":"run_command"}}`, `"decision":"deny"`},
		{"write_file", `{"toolCall":{"name":"write_file"}}`, `"decision":"deny"`},
		{"read_file", `{"toolCall":{"name":"read_file"}}`, `"decision":"deny"`},
	} {
		wrapper := filepath.Join(t.TempDir(), "run-gate")
		if err := os.WriteFile(wrapper, []byte("#!/bin/sh\nprintf '%s' "+shellQuote(tc.payload)+" | "+shellQuote(gate)+"\n"), 0o700); err != nil {
			t.Fatal(err)
		}
		ctx, cancel := budgetContext(5 * time.Second)
		out, err := toolprobe.OutputEnv(ctx, wrapper, "", []string{"PATH=" + os.Getenv("PATH")})
		cancel()
		if err != nil {
			t.Fatalf("%s: %v", tc.label, err)
		}
		if !strings.Contains(string(out), tc.want) {
			t.Errorf("gate(%s) = %s, want %s", tc.label, out, tc.want)
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
	if !strings.Contains(string(out), `"trustedWorkspaces"`) || !strings.Contains(string(out), `"mcp(world/*)"`) || !strings.Contains(string(out), "--model\ngemini") {
		t.Fatalf("unexpected agy-world argv: %s", out)
	}
	pluginMCP, err := os.ReadFile(filepath.Join(emptyDir, ".agents", "plugins", "ailang-world", "mcp_config.json"))
	if err != nil || !strings.Contains(string(pluginMCP), `"mcp-stdio"`) {
		t.Fatalf("staged plugin mcp_config.json: %v %s", err, pluginMCP)
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

func TestAgyInteractiveTranscriptCaptureAndGuards(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	hostsDir := filepath.Join(root, "hosts")
	if code, out := hostsInitTest(t, hostsDir); code != 0 {
		t.Fatalf("%d: %s", code, out)
	}
	launchDir := filepath.Join(root, "launch-ep")
	if err := os.MkdirAll(filepath.Join(launchDir, ".agents"), 0o700); err != nil {
		t.Fatal(err)
	}
	brainFile := filepath.Join(home, ".gemini", "antigravity-cli", "brain", "conv-123", ".system_generated", "logs", "transcript.jsonl")
	if err := os.MkdirAll(filepath.Dir(brainFile), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(brainFile, []byte("{\"type\":\"USER_INPUT\"}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	// Run agy-world-session-start hook and confirm it records .world-session-meta.json.
	hook := filepath.Join(hostsDir, "bin", "agy-world-session-start")
	payload, _ := json.Marshal(map[string]string{
		"conversationId": "conv-123",
		"transcriptPath": brainFile,
		"cwd":            launchDir,
	})
	wrapper := filepath.Join(t.TempDir(), "run-start")
	if err := os.WriteFile(wrapper, []byte("#!/bin/sh\nprintf '%s' "+shellQuote(string(payload))+" | "+shellQuote(hook)+"\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := budgetContext(5 * time.Second)
	out, err := toolprobe.OutputEnv(ctx, wrapper, "", []string{"PATH=" + os.Getenv("PATH"), "HOME=" + home})
	cancel()
	if err != nil {
		t.Fatalf("session-start hook: %v", err)
	}
	if !strings.Contains(string(out), `"injectSteps"`) {
		t.Fatalf("unexpected hook output: %s", out)
	}
	env := []string{"HOME=" + home}
	got := agyTranscriptPath(launchDir, env)
	if got.path != brainFile || got.id != "conv-123" || got.reason != "" {
		t.Fatalf("agyTranscriptPath = %+v", got)
	}
	// Refuse a transcript path outside agy's brain directory.
	outside := filepath.Join(root, "outside.jsonl")
	if err := os.WriteFile(outside, []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	badMeta, _ := json.Marshal(map[string]string{"conversationId": "conv-123", "transcriptPath": outside})
	if err := os.WriteFile(filepath.Join(launchDir, ".agents", ".world-session-meta.json"), badMeta, 0o600); err != nil {
		t.Fatal(err)
	}
	if refused := agyTranscriptPath(launchDir, env); refused.path != "" || !strings.Contains(refused.reason, "outside agy brain root") {
		t.Fatalf("expected outside path refusal, got %+v", refused)
	}
	// Refuse a symlink transcript inside agy's brain directory.
	linkPath := filepath.Join(home, ".gemini", "antigravity-cli", "brain", "conv-123", "link.jsonl")
	if err := os.Symlink(brainFile, linkPath); err != nil {
		t.Fatal(err)
	}
	linkMeta, _ := json.Marshal(map[string]string{"conversationId": "conv-123", "transcriptPath": linkPath})
	if err := os.WriteFile(filepath.Join(launchDir, ".agents", ".world-session-meta.json"), linkMeta, 0o600); err != nil {
		t.Fatal(err)
	}
	if refused := agyTranscriptPath(launchDir, env); refused.path != "" || !strings.Contains(refused.reason, "not a regular file") {
		t.Fatalf("expected symlink refusal, got %+v", refused)
	}
}
