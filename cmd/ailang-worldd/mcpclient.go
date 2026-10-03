package main

// Row 138 M1 (design_docs/planned/w-worldd-developer-cli.md §3.1, §3.2): a
// built-in MCP client, so an operator or an agent calls a World tool with one
// command instead of hand-rolling curl, two Accept types, SSE parsing and a
// Bearer header.
//
// /mcp/ answers in three shapes that a naive client conflates (V9–V11):
//
//	200 text/event-stream   the JSON-RPC response, as SSE `data:` lines
//	200 application/json    the frozen -32603 host-failure envelope (NOT SSE)
//	401 text/plain          a session denial
//
// classifyMCP tells them apart by status and Content-Type only, never by
// sniffing the body. A host failure is ambiguous — it may follow a commit that
// did land — so it runs the disambiguation probe (head before/after) instead
// of inviting a blind retry.

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

// exitIntegrity is the row-138 integrity refusal: a broken provenance link
// (why) or a committed ok:false under call --strict (D-CLI-3).
const exitIntegrity = 3

// sessionEnvVar is the session resolver's fallback after --session.
const sessionEnvVar = "WORLD_SESSION"

// maxSessionFileBytes bounds a session file: the token is 64 bytes.
const maxSessionFileBytes = 256

// mcpCommandBudget bounds one tools/call command end to end: the head probe,
// the call (each request is itself capped by the client deadline) and the
// second probe.
const mcpCommandBudget = 2 * time.Minute

// resolveSession is the ONE session resolver (§3.1), shared by tools, call
// and commit. Precedence: --session, then $WORLD_SESSION, else an error. A
// 64-hex value is the token itself; anything else names a file of at most
// 256 bytes holding it. It warns (stderr) when a token is given raw on argv
// and when the file is group/other readable.
//
// The token never appears in any output this function writes, nor in any
// error it returns (MUT-TOKEN-ECHO).
func resolveSession(flagValue string, stderr io.Writer) (string, error) {
	value, source := flagValue, "--session"
	if value == "" {
		value, source = os.Getenv(sessionEnvVar), "$"+sessionEnvVar
	}
	if value == "" {
		return "", fmt.Errorf("no session credential: pass --session <file> or set %s (mint one with `ailang-worldd session mint ... --out <file>`)", sessionEnvVar)
	}
	if isHex64(value) {
		if source == "--session" {
			fmt.Fprintln(stderr, "ailang-worldd: warning: a raw session token on the command line is visible to other local processes; "+
				"prefer --session <file> (mode 0600) or "+sessionEnvVar)
		}
		return value, nil
	}
	info, err := os.Stat(value)
	if err != nil {
		return "", fmt.Errorf("%s is neither a 64-hex session token nor a readable file: %v", source, errors.Unwrap(err))
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("%s session file %s is not a regular file", source, value)
	}
	if info.Size() > maxSessionFileBytes {
		return "", fmt.Errorf("%s session file %s is %d bytes; a session file holds one 64-hex token (at most %d bytes)",
			source, value, info.Size(), maxSessionFileBytes)
	}
	if info.Mode().Perm()&0o077 != 0 {
		fmt.Fprintf(stderr, "ailang-worldd: warning: session file %s is readable by group/other (mode %04o); chmod 600 it\n",
			value, info.Mode().Perm())
	}
	raw, err := os.ReadFile(value)
	if err != nil {
		return "", fmt.Errorf("read %s session file %s: %v", source, value, errors.Unwrap(err))
	}
	token := strings.TrimSpace(string(raw))
	if !isHex64(token) {
		return "", fmt.Errorf("%s session file %s does not hold one 64-hex session token", source, value)
	}
	return token, nil
}

// ---------------------------------------------------------------------------
// wire classification
// ---------------------------------------------------------------------------

type mcpKind int

const (
	mcpResult      mcpKind = iota // 200 SSE carrying a JSON-RPC result
	mcpToolError                  // 200 SSE carrying a JSON-RPC error
	mcpHostFailure                // 200 application/json -32603 envelope
	mcpDenied                     // 401/403: session denial (text/plain)
	mcpTransport                  // 400 and anything else: verbatim
)

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// mcpOutcome is one classified /mcp/ answer.
type mcpOutcome struct {
	Kind   mcpKind
	Status int
	Result json.RawMessage // mcpResult
	Error  *rpcError       // mcpToolError, mcpHostFailure
	Body   string          // mcpDenied, mcpTransport
}

// sseEvents parses a text/event-stream body into the data payload of each
// event, per the SSE spec: an event ends at a blank line; each `data:` line
// contributes its value (one leading space stripped); several data lines of
// one event are joined with "\n". Comment lines (":") and other fields are
// ignored. An event still open at EOF is dispatched too.
func sseEvents(body []byte) []string {
	var events []string
	var data []string
	hasData := false
	flush := func() {
		if hasData {
			events = append(events, strings.Join(data, "\n"))
		}
		data, hasData = nil, false
	}
	sc := bufio.NewScanner(bytes.NewReader(body))
	sc.Buffer(make([]byte, 0, 64*1024), maxClientResponseBytes+1)
	for sc.Scan() {
		line := strings.TrimSuffix(sc.Text(), "\r")
		if line == "" {
			flush()
			continue
		}
		if strings.HasPrefix(line, ":") {
			continue
		}
		field, value, found := strings.Cut(line, ":")
		if !found {
			field, value = line, ""
		}
		value = strings.TrimPrefix(value, " ")
		if field == "data" {
			data = append(data, value)
			hasData = true
		}
	}
	flush()
	return events
}

// classifyMCP classifies one /mcp/ answer by status and Content-Type (§3.2).
// wantID is the compact JSON of the request id: among several SSE events the
// one whose JSON-RPC id matches is the answer.
func classifyMCP(status int, contentType string, body []byte, wantID string) mcpOutcome {
	media, _, _ := mime.ParseMediaType(contentType)
	switch {
	case status == http.StatusOK && media == "text/event-stream":
		for _, data := range sseEvents(body) {
			var msg struct {
				ID     json.RawMessage `json:"id"`
				Result json.RawMessage `json:"result"`
				Error  *rpcError       `json:"error"`
			}
			if json.Unmarshal([]byte(data), &msg) != nil {
				continue
			}
			var id bytes.Buffer
			if json.Compact(&id, msg.ID) != nil || id.String() != wantID {
				continue
			}
			if msg.Error != nil {
				return mcpOutcome{Kind: mcpToolError, Status: status, Error: msg.Error}
			}
			return mcpOutcome{Kind: mcpResult, Status: status, Result: msg.Result}
		}
		return mcpOutcome{Kind: mcpTransport, Status: status,
			Body: fmt.Sprintf("no SSE event answers request id %s", wantID)}
	case status == http.StatusOK && media == "application/json":
		var env struct {
			Error *rpcError `json:"error"`
		}
		if json.Unmarshal(body, &env) == nil && env.Error != nil {
			return mcpOutcome{Kind: mcpHostFailure, Status: status, Error: env.Error}
		}
		return mcpOutcome{Kind: mcpTransport, Status: status, Body: strings.TrimSpace(string(body))}
	case status == http.StatusUnauthorized || status == http.StatusForbidden:
		return mcpOutcome{Kind: mcpDenied, Status: status, Body: strings.TrimSpace(string(body))}
	default:
		return mcpOutcome{Kind: mcpTransport, Status: status, Body: strings.TrimSpace(string(body))}
	}
}

// denialHint is the one-line fix for a 401 body (projection.go messages).
func denialHint(body string) string {
	switch {
	case strings.Contains(body, "absent"):
		return "pass --session <file> or set " + sessionEnvVar
	case strings.Contains(body, "malformed"):
		return "the session file must hold exactly one 64-hex token"
	case strings.Contains(body, "expired"):
		return "the session has expired: mint a new one (`ailang-worldd session mint`, daemon stopped)"
	case strings.Contains(body, "unknown"):
		return "this daemon's store does not know that credential: check --addr and that the session was minted into this --db (it may have been revoked)"
	default:
		return "check the session credential"
	}
}

// ---------------------------------------------------------------------------
// the client
// ---------------------------------------------------------------------------

type mcpClient struct {
	c     *client
	token string
}

func (m mcpClient) post(ctx context.Context, method string, id int, params any) (mcpOutcome, error) {
	req := map[string]any{"jsonrpc": "2.0", "id": id, "method": method}
	if params != nil {
		req["params"] = params
	}
	body, err := json.Marshal(req)
	if err != nil {
		return mcpOutcome{}, err
	}
	status, ctype, data, err := m.c.doRequest(ctx, http.MethodPost, "/mcp/", bytes.NewReader(body), map[string]string{
		"Authorization": "Bearer " + m.token,
		"Accept":        "application/json, text/event-stream",
	})
	if err != nil {
		return mcpOutcome{}, err
	}
	return classifyMCP(status, ctype, data, fmt.Sprint(id)), nil
}

// headState reads GET /v1/head: the selected world ref, or present=false on
// the 404 "no world head has been selected yet".
func headState(ctx context.Context, c *client) (ref string, present bool, err error) {
	status, body, err := c.do(ctx, http.MethodGet, "/v1/head", nil)
	if err != nil {
		return "", false, err
	}
	switch status {
	case http.StatusOK:
		return strings.TrimSpace(string(body)), true, nil
	case http.StatusNotFound:
		return "", false, nil
	default:
		return "", false, fmt.Errorf("GET /v1/head returned HTTP %d: %s", status, strings.TrimSpace(string(body)))
	}
}

// worldRevision resolves a world ref to its revision (= the entry index of
// its log head) through GET /v1/worlds/{ref}.
func worldRevision(ctx context.Context, c *client, ref string) (int64, error) {
	var w struct {
		Revision int64 `json:"revision"`
	}
	status, body, err := c.do(ctx, http.MethodGet, "/v1/worlds/"+url.PathEscape(ref), nil)
	if err != nil {
		return 0, err
	}
	if status != http.StatusOK || json.Unmarshal(body, &w) != nil {
		return 0, fmt.Errorf("GET /v1/worlds/%s returned HTTP %d", ref, status)
	}
	return w.Revision, nil
}

// reportOutcome writes a non-result outcome and returns its exit code. For a
// host failure it runs the disambiguation probe (§3.2): no head after the
// call → commit genesis; the head moved → a commit landed, do not retry.
func reportOutcome(ctx context.Context, c *client, verb string, out mcpOutcome, before string, beforeOK bool, stderr io.Writer) int {
	switch out.Kind {
	case mcpToolError:
		fmt.Fprintf(stderr, "ailang-worldd %s: tool error %d: %s\n", verb, out.Error.Code, out.Error.Message)
		return exitUsage
	case mcpDenied:
		fmt.Fprintf(stderr, "ailang-worldd %s: session denied (HTTP %d): %s\n  fix: %s\n", verb, out.Status, out.Body, denialHint(out.Body))
		return exitUsage
	case mcpTransport:
		fmt.Fprintf(stderr, "ailang-worldd %s: /mcp/ returned HTTP %d: %s\n", verb, out.Status, out.Body)
		return exitUsage
	case mcpHostFailure:
		after, afterOK, err := headState(ctx, c)
		switch {
		case err != nil:
			fmt.Fprintf(stderr, "ailang-worldd %s: host failure %d: %s (and the head probe failed: %v)\n", verb, out.Error.Code, out.Error.Message, err)
		case !afterOK:
			fmt.Fprintf(stderr, "ailang-worldd %s: host failure %d: %s\n  no world head: this store has no selected world yet; commit a genesis world first (see QUICKSTART §9)\n",
				verb, out.Error.Code, out.Error.Message)
		case !beforeOK || after != before:
			n, werr := worldRevision(ctx, c, after)
			if werr != nil {
				fmt.Fprintf(stderr, "ailang-worldd %s: host failure %d: %s\n  a commit landed (head %s) — do not retry; inspect it with: ailang-worldd why head\n",
					verb, out.Error.Code, out.Error.Message, after)
			} else {
				fmt.Fprintf(stderr, "ailang-worldd %s: host failure %d: %s\n  a commit landed (entry %d) — do not retry; inspect it with: ailang-worldd why %d\n",
					verb, out.Error.Code, out.Error.Message, n, n)
			}
		default:
			fmt.Fprintf(stderr, "ailang-worldd %s: host failure %d: %s\n", verb, out.Error.Code, out.Error.Message)
		}
		return exitUsage
	}
	return exitOK
}

// ---------------------------------------------------------------------------
// tools list
// ---------------------------------------------------------------------------

const toolsHelp = `usage: ailang-worldd [--addr <url>] tools list [--session <file|token>] [--json]

Lists the tools the session may call (MCP tools/list on /mcp/): name,
description and required arguments. --json prints the tools array verbatim.

The session is --session (a file holding the 64-hex token, preferred; or the
token itself, which warns) or $WORLD_SESSION.
`

func runTools(addr string, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" || args[0] == "help" {
		if len(args) == 0 {
			fmt.Fprint(stderr, toolsHelp)
			return exitUsage
		}
		fmt.Fprint(stdout, toolsHelp)
		return exitOK
	}
	if args[0] != "list" {
		fmt.Fprintf(stderr, "ailang-worldd tools: unknown subcommand %q\n%s", args[0], toolsHelp)
		return exitUsage
	}
	fs := flag.NewFlagSet("ailang-worldd tools list", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	session := fs.String("session", "", "session file or token (default $WORLD_SESSION)")
	asJSON := fs.Bool("json", false, "print the tools array as JSON")
	if code, done := parseVerbFlags(fs, args[1:], toolsHelp, stdout, stderr); done {
		return code
	}
	token, err := resolveSession(*session, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd tools list: %v\n", err)
		return exitUsage
	}
	ctx, cancel := budgetContext(mcpCommandBudget)
	defer cancel()
	c := newClient(addr)
	out, err := mcpClient{c: c, token: token}.post(ctx, "tools/list", 1, nil)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd tools list: %v\n", err)
		return exitUsage
	}
	if out.Kind != mcpResult {
		return reportOutcome(ctx, c, "tools list", out, "", true, stderr)
	}
	var res struct {
		Tools []struct {
			Name        string          `json:"name"`
			Description string          `json:"description"`
			InputSchema json.RawMessage `json:"inputSchema"`
		} `json:"tools"`
	}
	var raw struct {
		Tools json.RawMessage `json:"tools"`
	}
	if json.Unmarshal(out.Result, &res) != nil || json.Unmarshal(out.Result, &raw) != nil {
		fmt.Fprintf(stderr, "ailang-worldd tools list: malformed tools/list result\n")
		return exitUsage
	}
	if *asJSON {
		fmt.Fprintln(stdout, string(raw.Tools))
		return exitOK
	}
	for _, t := range res.Tools {
		var schema struct {
			Required []string `json:"required"`
		}
		_ = json.Unmarshal(t.InputSchema, &schema)
		req := "(none)"
		if len(schema.Required) > 0 {
			req = strings.Join(schema.Required, ", ")
		}
		fmt.Fprintf(stdout, "%-18s %s\n%-18s required: %s\n", t.Name, oneLine(t.Description, 100), "", req)
	}
	fmt.Fprintf(stdout, "%d tool(s)\n", len(res.Tools))
	return exitOK
}

// parseVerbFlags parses a verb's flags; --help prints help to stdout, rc 0.
// It reports done=true when the caller must return code.
func parseVerbFlags(fs *flag.FlagSet, args []string, help string, stdout, stderr io.Writer) (int, bool) {
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprint(stdout, help)
			return exitOK, true
		}
		fmt.Fprint(stderr, help)
		return exitUsage, true
	}
	return 0, false
}

// ---------------------------------------------------------------------------
// call
// ---------------------------------------------------------------------------

const callHelp = `usage: ailang-worldd [--addr <url>] call <tool> [--session <file|token>]
           [--arg k=v]... [--arg-json k=<json>]... | --json '<obj>'|@file|-
           [--json-out] [--strict]

Calls one tool (MCP tools/call on /mcp/). Arguments come from repeated --arg
(a string value) and --arg-json (any JSON value), or from one --json object
(inline, @file, or - for stdin); the two forms do not mix.

Output: the result's fields (long values elided with their byte counts), then
the world block: the plan ref and, per effect, its id, status and record ref.
--json-out prints the committed output bytes exactly (one trailing newline),
so sha256 of stdout minus that newline is the output ref, and the output can
be piped to 'ailang-worldd why -'.

Exit: 0 committed (including a committed refusal, ok:false); 3 with --strict
when the committed result has ok:false; 1 on a tool error, a session denial
or a host failure. A host failure is probed: "no world head" means commit a
genesis first; "a commit landed (entry N)" means do not retry — run why N.
`

type stringsFlag []string

func (s *stringsFlag) String() string     { return strings.Join(*s, ",") }
func (s *stringsFlag) Set(v string) error { *s = append(*s, v); return nil }

func runCall(addr string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	tool := ""
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		tool, args = args[0], args[1:]
	}
	fs := flag.NewFlagSet("ailang-worldd call", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	session := fs.String("session", "", "session file or token (default $WORLD_SESSION)")
	var kv, kvJSON stringsFlag
	fs.Var(&kv, "arg", "k=v string argument, repeatable")
	fs.Var(&kvJSON, "arg-json", "k=<json> argument, repeatable")
	whole := fs.String("json", "", "the whole arguments object: inline JSON, @file, or - for stdin")
	jsonOut := fs.Bool("json-out", false, "print the committed output bytes exactly")
	strict := fs.Bool("strict", false, "exit 3 when the committed result has ok:false")
	if code, done := parseVerbFlags(fs, args, callHelp, stdout, stderr); done {
		return code
	}
	if tool == "" && fs.NArg() > 0 {
		tool = fs.Arg(0)
		if fs.NArg() > 1 {
			fmt.Fprintf(stderr, "ailang-worldd call: unexpected argument %q\n", fs.Arg(1))
			return exitUsage
		}
	} else if fs.NArg() > 0 {
		fmt.Fprintf(stderr, "ailang-worldd call: unexpected argument %q (flags come after the tool name)\n", fs.Arg(0))
		return exitUsage
	}
	if tool == "" {
		fmt.Fprint(stderr, callHelp)
		return exitUsage
	}
	arguments, err := callArguments(kv, kvJSON, *whole, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd call: %v\n", err)
		return exitUsage
	}
	token, err := resolveSession(*session, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd call: %v\n", err)
		return exitUsage
	}
	ctx, cancel := budgetContext(mcpCommandBudget)
	defer cancel()
	c := newClient(addr)
	before, beforeOK, err := headState(ctx, c)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd call: %v\n", err)
		return exitUsage
	}
	out, err := mcpClient{c: c, token: token}.post(ctx, "tools/call", 1,
		map[string]any{"name": tool, "arguments": arguments})
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd call: %v\n", err)
		return exitUsage
	}
	if out.Kind != mcpResult {
		return reportOutcome(ctx, c, "call", out, before, beforeOK, stderr)
	}
	var res struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		IsError bool `json:"isError"`
	}
	if json.Unmarshal(out.Result, &res) != nil || len(res.Content) == 0 || res.Content[0].Type != "text" {
		fmt.Fprintf(stderr, "ailang-worldd call: malformed tools/call result: %s\n", oneLine(string(out.Result), 200))
		return exitUsage
	}
	text := res.Content[0].Text
	if res.IsError {
		fmt.Fprintf(stderr, "ailang-worldd call: tool reported an error: %s\n", oneLine(text, 400))
		return exitUsage
	}
	var fields map[string]json.RawMessage
	_ = json.Unmarshal([]byte(text), &fields)
	if *jsonOut {
		io.WriteString(stdout, text+"\n")
	} else {
		printCallResult(stdout, text, fields)
	}
	if *strict && string(fields["ok"]) == "false" {
		fmt.Fprintln(stderr, "ailang-worldd call: committed with ok:false (--strict)")
		return exitIntegrity
	}
	return exitOK
}

// callArguments builds the tools/call arguments object.
func callArguments(kv, kvJSON []string, whole string, stdin io.Reader) (json.RawMessage, error) {
	if whole != "" {
		if len(kv) > 0 || len(kvJSON) > 0 {
			return nil, errors.New("--json does not mix with --arg/--arg-json")
		}
		var raw []byte
		switch {
		case whole == "-":
			b, err := io.ReadAll(io.LimitReader(stdin, maxClientCommitBytes+1))
			if err != nil {
				return nil, fmt.Errorf("read stdin: %v", err)
			}
			raw = b
		case strings.HasPrefix(whole, "@"):
			b, err := os.ReadFile(whole[1:])
			if err != nil {
				return nil, fmt.Errorf("read %s: %v", whole[1:], errors.Unwrap(err))
			}
			raw = b
		default:
			raw = []byte(whole)
		}
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(raw, &obj); err != nil || obj == nil {
			return nil, errors.New("--json must be one JSON object")
		}
		return json.RawMessage(bytes.TrimSpace(raw)), nil
	}
	obj := map[string]json.RawMessage{}
	for _, spec := range kv {
		k, v, ok := strings.Cut(spec, "=")
		if !ok || k == "" {
			return nil, fmt.Errorf("--arg %q: want k=v", spec)
		}
		b, _ := json.Marshal(v)
		obj[k] = b
	}
	for _, spec := range kvJSON {
		k, v, ok := strings.Cut(spec, "=")
		if !ok || k == "" {
			return nil, fmt.Errorf("--arg-json %q: want k=<json>", spec)
		}
		if !json.Valid([]byte(v)) {
			return nil, fmt.Errorf("--arg-json %s: value is not valid JSON", k)
		}
		obj[k] = json.RawMessage(v)
	}
	return json.Marshal(obj)
}

const elideAt = 160

// printCallResult renders a result: each top-level field but world (long
// values elided with byte counts), then the world block.
func printCallResult(w io.Writer, text string, fields map[string]json.RawMessage) {
	if fields == nil {
		fmt.Fprintf(w, "result: %s\n", elide(text))
		return
	}
	keys := make([]string, 0, len(fields))
	for k := range fields {
		if k != "world" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		var s string
		if json.Unmarshal(fields[k], &s) == nil {
			fmt.Fprintf(w, "%-12s %s\n", k+":", elide(s))
		} else {
			fmt.Fprintf(w, "%-12s %s\n", k+":", elide(string(fields[k])))
		}
	}
	sum := sha256.Sum256([]byte(text))
	fmt.Fprintf(w, "output:      sha256:%s (%d bytes)\n", hex.EncodeToString(sum[:]), len(text))
	var wb struct {
		Plan    string `json:"plan"`
		Effects []struct {
			ID     string  `json:"id"`
			Status string  `json:"status"`
			Record *string `json:"record"`
		} `json:"effects"`
	}
	if raw, ok := fields["world"]; !ok || json.Unmarshal(raw, &wb) != nil {
		fmt.Fprintln(w, "world:       (no world block: a pure transition)")
		return
	}
	fmt.Fprintf(w, "world:\n  plan       %s\n", wb.Plan)
	if len(wb.Effects) == 0 {
		fmt.Fprintln(w, "  effects    (none: the plan answered without effects)")
	}
	for _, e := range wb.Effects {
		rec := "-"
		if e.Record != nil {
			rec = *e.Record
		}
		fmt.Fprintf(w, "  effect     %s %s %s\n", e.ID, e.Status, rec)
	}
}

// elide shortens a value to elideAt bytes plus its byte count, on one line.
func elide(s string) string {
	one := strings.ReplaceAll(s, "\n", `\n`)
	if len(one) <= elideAt {
		return one
	}
	return fmt.Sprintf("%s… (%d bytes)", strings.ToValidUTF8(one[:elideAt], ""), len(s))
}

func oneLine(s string, n int) string {
	s = strings.Join(strings.Fields(s), " ")
	if len(s) > n {
		return strings.ToValidUTF8(s[:n], "") + "…"
	}
	return s
}
