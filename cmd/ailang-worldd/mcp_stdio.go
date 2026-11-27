package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const mcpStdioHelp = `usage: ailang-worldd [--addr <url>] mcp-stdio [--addr <url>] [--session <file|token>]

Bridge newline-delimited JSON-RPC MCP messages on stdin/stdout to the World
daemon's Streamable HTTP /mcp/ endpoint using the resolved session token.
Tokens are read from --session or $WORLD_SESSION and are never logged or printed.
`

func runMCPStdio(addr string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("mcp-stdio", flag.ContinueOnError)
	fs.SetOutput(stderr)
	localAddr := fs.String("addr", addr, "World daemon base URL")
	sessionFlag := fs.String("session", "", "session token file or 64-hex token (default $WORLD_SESSION)")
	if code, done := parseVerbFlags(fs, args, mcpStdioHelp, stdout, stderr); done {
		return code
	}
	if fs.NArg() != 0 {
		fmt.Fprint(stderr, mcpStdioHelp)
		return exitUsage
	}
	token, err := resolveSession(*sessionFlag, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd mcp-stdio: %v\n", err)
		return exitUsage
	}
	base := strings.TrimRight(clientBaseURL(*localAddr), "/")
	endpoint := base + "/mcp/"
	httpClient := &http.Client{}
	if err := bridgeMCPStdio(stdin, stdout, httpClient, endpoint, token); err != nil {
		fmt.Fprintf(stderr, "ailang-worldd mcp-stdio: %s\n", sanitizeToken(err.Error(), token))
		return exitFatal
	}
	return exitOK
}

func sanitizeToken(s, token string) string {
	if token == "" {
		return s
	}
	return strings.ReplaceAll(s, token, "<redacted>")
}

const mcpStdioMaxConcurrent = 8

type mcpStdioEnvelope struct {
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

func canonicalMCPRequestID(raw json.RawMessage) string {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return ""
	}
	return string(trimmed)
}

type mcpStdioBridge struct {
	stdout   io.Writer
	client   *http.Client
	endpoint string
	token    string

	writeMu sync.Mutex

	sidMu        sync.RWMutex
	mcpSessionID string

	inFlightMu sync.Mutex
	inFlight   map[string]context.CancelFunc

	errMu    sync.Mutex
	writeErr error
}

func (b *mcpStdioBridge) sessionID() string {
	b.sidMu.RLock()
	defer b.sidMu.RUnlock()
	return b.mcpSessionID
}

func (b *mcpStdioBridge) setSessionID(sid string) {
	if sid == "" {
		return
	}
	b.sidMu.Lock()
	b.mcpSessionID = sid
	b.sidMu.Unlock()
}

func (b *mcpStdioBridge) registerCancel(idKey string, cancel context.CancelFunc) {
	if idKey == "" {
		return
	}
	b.inFlightMu.Lock()
	b.inFlight[idKey] = cancel
	b.inFlightMu.Unlock()
}

func (b *mcpStdioBridge) unregisterCancel(idKey string) {
	if idKey == "" {
		return
	}
	b.inFlightMu.Lock()
	delete(b.inFlight, idKey)
	b.inFlightMu.Unlock()
}

func (b *mcpStdioBridge) cancelRequest(idKey string) {
	if idKey == "" {
		return
	}
	b.inFlightMu.Lock()
	cancel := b.inFlight[idKey]
	b.inFlightMu.Unlock()
	if cancel != nil {
		cancel()
	}
}

func (b *mcpStdioBridge) recordWriteErr(err error) {
	if err == nil {
		return
	}
	b.errMu.Lock()
	if b.writeErr == nil {
		b.writeErr = err
	}
	b.errMu.Unlock()
}

func (b *mcpStdioBridge) getWriteErr() error {
	b.errMu.Lock()
	defer b.errMu.Unlock()
	return b.writeErr
}

func (b *mcpStdioBridge) writeRPCError(id json.RawMessage, code int, message string) {
	b.writeMu.Lock()
	writeStdioRPCError(b.stdout, id, code, message)
	b.writeMu.Unlock()
}

func (b *mcpStdioBridge) handleLine(line []byte, envelope mcpStdioEnvelope, idKey string) error {
	hasID := idKey != ""
	ctx, cancel := budgetContext(mcpCommandBudget)
	b.registerCancel(idKey, cancel)
	defer func() {
		b.unregisterCancel(idKey)
		cancel()
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, b.endpoint, bytes.NewReader(line))
	if err != nil {
		if hasID {
			b.writeRPCError(envelope.ID, -32603, sanitizeToken(err.Error(), b.token))
		}
		return nil
	}
	req.Header.Set("Authorization", "Bearer "+b.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	if sid := b.sessionID(); sid != "" {
		req.Header.Set("Mcp-Session-Id", sid)
	}

	resp, err := b.client.Do(req)
	if err != nil {
		if errors.Is(ctx.Err(), context.Canceled) {
			return nil
		}
		if hasID {
			b.writeRPCError(envelope.ID, -32000, "world daemon unreachable: "+sanitizeToken(err.Error(), b.token))
		}
		return nil
	}
	b.setSessionID(resp.Header.Get("Mcp-Session-Id"))
	body, readErr := io.ReadAll(io.LimitReader(resp.Body, 16*1024*1024))
	_ = resp.Body.Close()
	if readErr != nil {
		if errors.Is(ctx.Err(), context.Canceled) {
			return nil
		}
		if hasID {
			b.writeRPCError(envelope.ID, -32000, "read world response: "+sanitizeToken(readErr.Error(), b.token))
		}
		return nil
	}
	if resp.StatusCode == http.StatusAccepted || (resp.StatusCode == http.StatusOK && len(bytes.TrimSpace(body)) == 0) {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		if hasID {
			msg := strings.TrimSpace(sanitizeToken(string(body), b.token))
			if msg == "" {
				msg = http.StatusText(resp.StatusCode)
			}
			b.writeRPCError(envelope.ID, -32000, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, msg))
		}
		return nil
	}
	mediaType, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if mediaType == "text/event-stream" {
		b.writeMu.Lock()
		err := writeSSEDataLines(b.stdout, body)
		b.writeMu.Unlock()
		if err != nil {
			b.recordWriteErr(err)
			return err
		}
		if envelope.Method == "tools/list" {
			markMCPStdioReady()
		}
		return nil
	}
	var compact bytes.Buffer
	if err := json.Compact(&compact, body); err != nil {
		if hasID {
			b.writeRPCError(envelope.ID, -32603, "invalid JSON from world daemon")
		}
		return nil
	}
	b.writeMu.Lock()
	_, err = fmt.Fprintf(b.stdout, "%s\n", compact.Bytes())
	b.writeMu.Unlock()
	if err != nil {
		b.recordWriteErr(err)
		return err
	}
	if envelope.Method == "tools/list" {
		markMCPStdioReady()
	}
	return nil
}

func bridgeMCPStdio(stdin io.Reader, stdout io.Writer, client *http.Client, endpoint, token string) error {
	scanner := bufio.NewScanner(stdin)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	b := &mcpStdioBridge{
		stdout:   stdout,
		client:   client,
		endpoint: endpoint,
		token:    token,
		inFlight: make(map[string]context.CancelFunc),
	}
	sem := make(chan struct{}, mcpStdioMaxConcurrent)
	var wg sync.WaitGroup
	for scanner.Scan() {
		if b.getWriteErr() != nil {
			break
		}
		trimmed := bytes.TrimSpace(scanner.Bytes())
		if len(trimmed) == 0 {
			continue
		}
		line := append([]byte(nil), trimmed...)
		var envelope mcpStdioEnvelope
		_ = json.Unmarshal(line, &envelope)
		if envelope.Method == "notifications/cancelled" {
			var p struct {
				RequestID json.RawMessage `json:"requestId"`
			}
			if json.Unmarshal(envelope.Params, &p) == nil {
				b.cancelRequest(canonicalMCPRequestID(p.RequestID))
			}
			continue
		}
		idKey := canonicalMCPRequestID(envelope.ID)
		if envelope.Method == "tools/call" {
			sem <- struct{}{}
			wg.Add(1)
			go func(payload []byte, env mcpStdioEnvelope, key string) {
				defer func() {
					<-sem
					wg.Done()
				}()
				_ = b.handleLine(payload, env, key)
			}(line, envelope, idKey)
			continue
		}
		if err := b.handleLine(line, envelope, idKey); err != nil {
			break
		}
	}
	wg.Wait()
	if err := b.getWriteErr(); err != nil {
		return err
	}
	if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

func markMCPStdioReady() {
	ready := os.Getenv("WORLD_AGY_MCP_READY")
	if ready == "" || !filepath.IsAbs(ready) {
		return
	}
	_ = os.WriteFile(ready, []byte("ready\n"), 0o600)
}

func writeSSEDataLines(w io.Writer, body []byte) error {
	sc := bufio.NewScanner(bytes.NewReader(body))
	sc.Buffer(make([]byte, 64*1024), 16*1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "" {
			continue
		}
		if _, err := fmt.Fprintf(w, "%s\n", payload); err != nil {
			return err
		}
	}
	return sc.Err()
}

func writeStdioRPCError(w io.Writer, id json.RawMessage, code int, message string) {
	resp := struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      json.RawMessage `json:"id"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}{
		JSONRPC: "2.0",
		ID:      id,
	}
	resp.Error.Code = code
	resp.Error.Message = message
	b, err := json.Marshal(resp)
	if err != nil {
		return
	}
	_, _ = fmt.Fprintf(w, "%s\n", b)
}

// stageAgyHeadlessPlugin writes an ephemeral Antigravity plugin into the empty
// launch directory so headless `agy -p` connects to World via `mcp-stdio` (reading
// $WORLD_SESSION from its environment, never writing the token to disk) and
// enforces tool confinement through its PreToolUse hook and AGENTS.md rule.
func stageAgyHeadlessPlugin(dir, addr, contextText string) error {
	base := addr
	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		base = "http://" + base
	}
	pluginDir := filepath.Join(dir, ".agents", "plugins", "ailang-world")
	rulesDir := filepath.Join(pluginDir, "rules")
	if err := os.MkdirAll(rulesDir, 0o700); err != nil {
		return err
	}
	write0600 := func(path string, data []byte) error {
		return os.WriteFile(path, data, 0o600)
	}
	if err := write0600(filepath.Join(pluginDir, "plugin.json"), []byte("{\"name\":\"ailang-world\"}\n")); err != nil {
		return err
	}
	mcpCfg, err := json.Marshal(map[string]any{"mcpServers": map[string]any{"world": map[string]any{
		"command": hostsWorlddPath(),
		"args":    []string{"mcp-stdio", "--addr", base},
	}}})
	if err != nil {
		return err
	}
	if err := write0600(filepath.Join(pluginDir, "mcp_config.json"), append(mcpCfg, '\n')); err != nil {
		return err
	}
	gatePath := filepath.Join(pluginDir, "gate.sh")
	gateScript := []byte(`#!/bin/sh
payload=$(cat)
name=$(printf '%s' "$payload" | jq -r '.toolCall.name // ""' 2>/dev/null)
server=$(printf '%s' "$payload" | jq -r '.toolCall.args.ServerName // .toolCall.args.server_name // ""' 2>/dev/null)
case "$name" in
  call_mcp_tool)
    case "$server" in
      ailang-world_world|world)
        printf '{"decision":"allow"}\n'
        exit 0
        ;;
    esac
    ;;
  read_file)
    path=$(printf '%s' "$payload" | jq -r '.toolCall.args.Path // .toolCall.args.path // ""' 2>/dev/null)
    case "$path" in
      *../*|*/../*|*/..) ;;
      */.gemini/antigravity-cli/brain/*/.system_generated/steps/*/output.txt)
        printf '{"decision":"allow"}\n'
        exit 0
        ;;
    esac
    ;;
  mcp_ailang-world_world_*|ailang-world_world_*|mcp_world_*|mcp__world__*|world_*|mcp:world/*|mcp:ailang-world_world/*|ailang-read|ailang-write|ailang-edit|ailang-delete|ailang-check|ailang-run|ailang-cli|workspace-exec|workspace-search|workspace-list|builtins-search|examples-search|update_steps|ask_question|AskQuestion)
    printf '{"decision":"allow"}\n'
    exit 0
    ;;
esac
printf '{"decision":"deny","reason":"Confined to AILANG World: use World MCP tools (ailang-read, ailang-edit, ailang-write, workspace-exec)."}\n'
`)
	if err := os.WriteFile(gatePath, gateScript, 0o700); err != nil {
		return err
	}
	readyPath := filepath.Join(pluginDir, "ready.sh")
	readyScript := []byte(`#!/bin/sh
payload=$(cat)
inv=$(printf '%s' "$payload" | jq -r '.invocationNum // 0' 2>/dev/null)
if [ "$inv" = "0" ] || [ "$inv" = "1" ]; then
  if [ -n "${WORLD_AGY_MCP_READY:-}" ] && [ ! -f "${WORLD_AGY_MCP_READY}.waited" ]; then
    i=0
    while [ "$i" -lt 50 ] && [ ! -f "$WORLD_AGY_MCP_READY" ]; do
      sleep 0.1
      i=$((i + 1))
    done
    if [ -f "$WORLD_AGY_MCP_READY" ]; then
      sleep 1.2
      : >"${WORLD_AGY_MCP_READY}.waited" 2>/dev/null || :
    fi
  fi
fi
printf '{}\n'
`)
	if err := os.WriteFile(readyPath, readyScript, 0o700); err != nil {
		return err
	}
	hooksCfg, err := json.Marshal(map[string]any{"ailang-world": map[string]any{
		"PreInvocation": []any{map[string]any{"type": "command", "command": shellQuote(readyPath), "timeout": 10}},
		"PreToolUse": []any{map[string]any{"matcher": "*", "hooks": []any{
			map[string]any{"type": "command", "command": shellQuote(gatePath), "timeout": 10},
		}}},
	}})
	if err != nil {
		return err
	}
	if err := write0600(filepath.Join(pluginDir, "hooks.json"), append(hooksCfg, '\n')); err != nil {
		return err
	}
	rules := contextText
	if rules == "" {
		rules = "You are working inside AILANG World (confined).\n" + worldToolMap + worldBoundary
	}
	rules += "\nWorld's MCP tools are served by MCP server `ailang-world_world` (if listed as lazy tools, call them via `call_mcp_tool` with `ServerName: \"ailang-world_world\"`).\n"
	if !strings.HasSuffix(rules, "\n") {
		rules += "\n"
	}
	return write0600(filepath.Join(rulesDir, "AGENTS.md"), []byte(rules))
}
