---
title: Connecting
sidebar_position: 2
description: Connect pi, Claude Code, codex, curl or an A2A client to an AILANG World daemon.
---

# Connecting

Your operator gives you two things: the daemon's address (by default `http://127.0.0.1:7644`) and a **session credential**, a 64-character hex string bound to one episode and a set of grants. Keep the credential out of files you commit, logs and transcripts. The examples below read it from an environment variable or a file, never inline.

| Surface | Endpoint | Auth |
|---|---|---|
| MCP (streamable HTTP) | `POST /mcp/` | `Authorization: Bearer <session>` |
| A2A agent card | `GET /.well-known/agent.json` | `Authorization: Bearer <session>` |
| A2A JSON-RPC | `POST /a2a/` | `Authorization: Bearer <session>` |
| Log and object reads | `GET /v1/log/…`, `GET /v1/objects/…` | see [Provenance](./provenance.md) |

A missing, malformed, unknown or expired credential is refused before anything runs. On `/mcp/` that is HTTP 401 with one of these messages:

```text
session credential is absent: send Authorization: Bearer <session-credential>
malformed Authorization header: expected Bearer <64-hex-credential>
unknown session credential
session credential has expired
```

## MCP over HTTP with curl

Send both `Content-Type: application/json` and `Accept: application/json, text/event-stream`. Successful replies are Server-Sent Events: the JSON-RPC response is the `data:` line of an `event: message`.

```bash
export WORLD_SESSION="$(cat /path/to/session-file)"
curl -s -H "Authorization: Bearer $WORLD_SESSION" \
  -H 'Content-Type: application/json' -H 'Accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' http://127.0.0.1:7644/mcp/
curl -s -H "Authorization: Bearer $WORLD_SESSION" \
  -H 'Content-Type: application/json' -H 'Accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"ailang-read","arguments":{"path":"hello.ail"}}}' \
  http://127.0.0.1:7644/mcp/
```

The reply to the second call, as recorded in the attended smoke (digests shortened):

```text
event: message
data: {"jsonrpc":"2.0","id":2,"result":{"content":[{"type":"text","text":"{…same object as structuredContent…}"}],"structuredContent":{"content":"module hello\n\nimport std/io (println)\n\nexport func main() -> () ! {IO} {\n  println(\"hello from ailang-run\")\n}\n","ok":true,"policy_digest":"87d52ed5…","tool":"sha256:e55ff71c…","world":{"effects":[{"id":"e1","record":"sha256:76404771…","status":"ok"}],"plan":"sha256:ce4340a1…"}}}}
```

The result is delivered twice: as one `text` content item holding the JSON, and as `structuredContent`. Read `structuredContent` when your client exposes it.

MCP facts worth knowing:

- **Tool names are the tool ids** (`ailang-read`, `builtins-search`, …). The MCP name encoding escapes only `_`, `.` and `/`, and these ids contain none of them.
- **`tools/list` is filtered by your grants.** A tool you hold no grant for is absent, and calling it anyway is a JSON-RPC error `unknown tool "<name>"` (code `-32602`).
- **The server is stateless.** No `initialize` handshake is required for `tools/list` or `tools/call`. `MCP-Protocol-Version` may be omitted, `2025-03-26`, `2025-06-18` or `2025-11-25`; only `2025-03-26` (or no header) accepts JSON-RPC batches.
- **Each call is a new invocation.** MCP JSON-RPC ids are not idempotency keys: a repeated `tools/call` runs again and is recorded again.
- **Deadlines.** One POST has 20 seconds in total; `ailang-run` programs are stopped after 8 seconds by the AILANG supervisor.

## pi (pi-mcp-adapter)

Measured with pi 0.85.1 and pi-mcp-adapter 2.32.1. `directTools: true` registers each World tool as its own pi tool, named `<server>_<tool>`, so server `world` gives `world_ailang-read` and so on.

```json title="world-mcp.json"
{"mcpServers": {"world": {"url": "http://127.0.0.1:7644/mcp/", "auth": "bearer",
  "bearerTokenEnv": "WORLD_SESSION", "directTools": true, "toolPrefix": "server"}}}
```

Direct tools register from the adapter's metadata cache (`~/.pi/agent/mcp-cache.json`). Warm it once, then run with every built-in tool disabled so World's tools are the only ones available:

```bash
export WORLD_SESSION="$(cat /path/to/session-file)"
pi --mcp-config world-mcp.json --no-builtin-tools -p "Reply with the single word ok." </dev/null >/dev/null
pi --mcp-config world-mcp.json --no-builtin-tools \
  --tools world_ailang-read,world_ailang-write,world_ailang-edit,world_ailang-check,world_ailang-run,world_builtins-search,world_examples-search,world_ailang-cli \
  -p "Read hello.ail, change its message to 'hello from pi', check it, then run it and report its stdout." </dev/null
```

`</dev/null` matters: pi reads a non-TTY stdin as extra prompt input and waits for end of file. In an interactive `pi --mcp-config world-mcp.json` session, `/mcp tools` lists the eight `world_*` names.

## Claude Code

Let `claude mcp add` write the config. With `-s project` it writes `.mcp.json` in the current directory, and that file holds the raw bearer token, so do this in a scratch directory and delete it afterwards.

```bash
mkdir -p /tmp/world-claude && cd /tmp/world-claude
claude mcp add --transport http world http://127.0.0.1:7644/mcp/ \
  --header "Authorization: Bearer $WORLD_SESSION" -s project
```

Then run with the built-in tools disabled. `--allowedTools` alone leaves the built-in tools in place; `--tools ""` removes them, and `--strict-mcp-config` ignores every other MCP config:

```bash
claude -p "Read hello.ail, change its message to 'hello from claude', check it, then run it and report its stdout." \
  --tools "" --strict-mcp-config --mcp-config /tmp/world-claude/.mcp.json \
  --allowedTools mcp__world__ailang-read,mcp__world__ailang-write,mcp__world__ailang-edit,mcp__world__ailang-check,mcp__world__ailang-run,mcp__world__builtins-search,mcp__world__examples-search,mcp__world__ailang-cli \
  </dev/null
```

Claude Code names the tools `mcp__world__<id>`. Confirm them with `/mcp` in an interactive `claude --strict-mcp-config --mcp-config /tmp/world-claude/.mcp.json` session.

## codex

codex was not part of the attended smoke. The registration below was measured against codex-cli 0.159.2 (`codex mcp add --help`); the end-to-end run is not yet verified.

```bash
codex mcp add world --url http://127.0.0.1:7644/mcp/ --bearer-token-env-var WORLD_SESSION
codex mcp list
```

That writes this block to `~/.codex/config.toml` (or `$CODEX_HOME/config.toml`). The token itself is read from `WORLD_SESSION` at run time and never written to the file:

```toml
[mcp_servers.world]
url = "http://127.0.0.1:7644/mcp/"
bearer_token_env_var = "WORLD_SESSION"
```

codex keeps its own shell and file tools alongside MCP servers. Removing or sandboxing them is the operator's job, through codex's own configuration. If you are a codex agent working on a World episode, use the World tools for every read, write and run in the worktree, so that each action is recorded.

## A2A

The agent card lists your session's tools as skills, with the tool id as the skill `id`:

```bash
curl -s -H "Authorization: Bearer $WORLD_SESSION" http://127.0.0.1:7644/.well-known/agent.json
```

Invoke a tool with `tasks/send`: `metadata.skill_id` names the tool, and the message carries exactly one `data` part holding the arguments object.

```bash
curl -s -H "Authorization: Bearer $WORLD_SESSION" -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tasks/send","params":{"id":"read-hello-1","metadata":{"skill_id":"ailang-read"},"message":{"role":"user","parts":[{"type":"data","data":{"path":"hello.ail"}}]}}}' \
  http://127.0.0.1:7644/a2a/
```

A completed task carries the tool result in `artifacts[0].parts[0].data`, plus `metadata.invocation_id`, `metadata.world_ref` and `metadata.entry_index`.

Unlike MCP, **the A2A task id is an idempotency key within your session**. Sending the same `params.id` again returns the committed result without running the tool again. That makes A2A the safer surface when you must retry a write or a run after a lost response. A skill id outside your grants is refused with `not authorized` (code `-32602`).
