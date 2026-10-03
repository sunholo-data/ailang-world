---
title: MCP and A2A
sidebar_position: 3
description: The two agent-facing endpoints — MCP over Streamable HTTP and A2A JSON-RPC — with authentication, naming, framing, limits and the full error mapping.
---

# MCP and A2A

World serves its transition registry to agents over two open protocols and invents no wire
format of its own. Both endpoints use the same session authority, the same grant-filtered view of
the registry, and the same invocation coordinator. The MCP handler is the released
`serveapi/protocol/mcphttp` package from the AILANG Go module (`v0.47.2`); World supplies the
session, the tool list and the invoker.

## Authentication

Both endpoints require:

```http
Authorization: Bearer <64-hex-credential>
```

The credential is a session minted with `ailang-worldd session mint`. No other header is read.

| Endpoint | Absent | Malformed | Unknown | Expired |
|---|---|---|---|---|
| `POST /mcp/` | HTTP `401`, text `session credential is absent: send Authorization: Bearer <session-credential>` | HTTP `401`, `malformed Authorization header: expected Bearer <64-hex-credential>` | HTTP `401`, `unknown session credential` | HTTP `401`, `session credential has expired` |
| `POST /a2a/` | JSON-RPC error `-32001` | JSON-RPC error `-32600` | `-32001` | `-32001` |
| `GET /.well-known/agent.json` | `401 SessionAbsent` (REST envelope) | `400 InvalidSession` | `401 SessionUnknown` | `401 SessionExpired` |

Authorization is checked first; an unauthorized caller learns nothing about the transport.

## What a session sees

A session sees exactly the transitions whose `access` its live grants satisfy. The A2A card's
`skills` and MCP's `tools/list` come from the same function, and admission to a call rebuilds
that view fresh for each request: having seen a tool earlier conveys no authority.

## MCP: `POST /mcp/`

### Request requirements

- `Content-Type: application/json`, else `400` `Content-Type must be application/json`.
- `Accept` listing both `application/json` and `text/event-stream`, else `400`
  `Accept must list both application/json and text/event-stream`.
- Optional `MCP-Protocol-Version`: `2025-11-25`, `2025-06-18` or `2025-03-26`. Without it the
  server assumes `2025-03-26`. Any other value is `400`.
- Body at most 4 MiB.

Methods: `initialize`, `ping`, `tools/list`, `tools/call`. Notifications and responses are
accepted and ignored (`202`). The server is stateless and advertises `tools.listChanged: false`.

### Responses are SSE

A successful reply is one Server-Sent Event:

```text
event: message
data: {"jsonrpc":"2.0","id":2,"result":{…}}
```

Read the JSON-RPC object from the `data:` line.

### Tool names

The MCP name is the transition ID with a reversible escape: `_` → `_u`, `.` → `_d`, `/` → `_s`.
So `tools.echo` is `tools_decho`, and the `se-tools` IDs such as `ailang-read` are their own
names (`-` needs no escaping). A name longer than 64 bytes refuses the whole tool surface rather
than dropping or truncating a tool. An input schema missing a top-level `type` gains
`"type":"object"`, with constraints preserved.

### `tools/call` results

The transition's output JSON object is returned twice, as upstream MCP does: as text and as
`structuredContent`.

```json
{"jsonrpc":"2.0","id":2,"result":{
  "content":[{"type":"text","text":"{…the output object…}"}],
  "structuredContent":{…the output object…}}}
```

For the coding tools the output object is the handler's fields plus `world`:

```json
{"ok":true,"content":"…","tool":"sha256:…","policy_digest":"…",
 "world":{"effects":[{"id":"e1","status":"ok","record":"sha256:…"}],"plan":"sha256:…"}}
```

`status` is `ok`, `failed`, `denied` or `skipped`. A plan refusal returns the plan's result plus
`"world":{"effects":[],"plan":"sha256:…"}`.

An MCP result does **not** carry the log entry index. Use the effect record to walk back to the
entry ([Provenance walks](../guides/provenance-walks.md)).

### Batches

JSON arrays are accepted only under `2025-03-26` (the default) and run sequentially; `2025-06-18`
and `2025-11-25` refuse batches with `400`. Each call item gets fresh admission and a new random
task ID, even when JSON-RPC IDs repeat. If a host call fails, the whole POST is answered with one
error envelope, and earlier items may already have committed.

### Retries and idempotency

MCP JSON-RPC IDs give **no** idempotency. Every `tools/call` mints a new random task ID, so a
retry is a new invocation and may repeat an effect. Every execution is recorded. Before retrying
a lost response, inspect the log.

### Limits

The aggregate POST deadline is 20 s, including credential resolution (3 s) and initial tool
admission (10 s); each invocation re-admits under its own 10 s read bound. Up to eight host
callbacks run at once, each bounded at 20 s; the subprocess limit is also eight.

## A2A: `POST /a2a/` and the agent card

`GET /.well-known/agent.json` returns the session's card (see the
[HTTP API](http-api.md#get-well-knownagentjson)). Skill IDs are the transition IDs, verbatim.

`POST /a2a/` accepts one JSON-RPC request (at most 1 MiB) with method `tasks/send`:

```bash
curl -s -H "Authorization: Bearer $(cat /tmp/se-session)" \
  -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tasks/send","params":{"id":"read-1","metadata":{"skill_id":"ailang-read"},"message":{"role":"user","parts":[{"type":"data","data":{"path":"hello.ail"}}]}}}' \
  http://127.0.0.1:7644/a2a/
```

- `params.metadata.skill_id` names the transition.
- `params.message.parts` must be exactly one part of `"type":"data"` whose `data` is the
  arguments object.
- `params.id` is the task ID: 1–128 characters from `A–Z a–z 0–9 . _ -`. It is an idempotency key
  within the episode: sending the same task ID again returns the committed result without running
  the transition again.
- Optional `params.metadata.transition_fn` pins the exact transition source hash; a mismatch is
  refused.

A success:

```json
{"jsonrpc":"2.0","id":1,"result":{
  "id":"read-1","status":{"state":"completed"},
  "artifacts":[{"parts":[{"type":"data","data":{…the output object…}}]}],
  "metadata":{"invocation_id":"a2a:ep1:read-1","world_ref":"sha256:…","entry_index":3}}}
```

`metadata.entry_index` is the log entry this call committed.

## Error mapping

### A2A (`/a2a/`)

Errors are JSON-RPC errors in an HTTP `200`.

| Code | Message | Cause |
|---|---|---|
| `-32001` | session absent / unknown / expired messages | Session denial |
| `-32600` | `invalid JSON-RPC request` / malformed header message | Bad envelope or `jsonrpc` not `2.0`; malformed `Authorization` |
| `-32601` | `method not found` | Method is not `tasks/send` |
| `-32602` | `invalid params` | Bad params, wrong part shape, bad `transition_fn` |
| `-32602` | `not authorized` | Skill unknown, not visible to this session, or access denied |
| `-32602` | `proposal does not match the registered transition` | `transition_fn` pin does not match |
| `-32602` | `task id already used in this session` | Task ID reused for a different call, or still in flight |
| `-32603` | `effects were requested but the invocation was not confirmed committed; effect records: sha256:… …` | **EffectsUnrecorded**: an effect ran and was recorded, then a post-effect step failed (head conflict, finish failure, commit error). Checked first, whatever the cause. |
| `-32603` | `transition declares an effect this daemon has no handler for` | Rule R8: tools not enabled for this daemon or episode |
| `-32603` | `no world is selected; commit a genesis world first` | No genesis commit yet |
| `-32603` | `transition does not implement the invocation calling convention` | Incompatible transition |
| `-32603` | `transition execution failed` | The capsule run failed |
| `-32603` | `transition output is not a JSON object` | Bad output |
| `-32603` | `invocation exceeded its deadline` | 20 s deadline; the connection is closed |
| `-32603` | `invocation outcome is not confirmed; resend the same task id` | Commit outcome uncertain |
| `-32603` | `invocation was not committed; send a new task id` | Known not committed |
| `-32603` | `world head moved during invocation; not committed; send a new task id` | Head conflict on a pure transition |
| `-32603` | `transition invocation is not available in this daemon` | No coordinator (`serve` without `--ailang-bin`), registry read failure, integrity failure; the cause is logged to stderr as an `a2a refusal` line |

### `EffectsUnrecorded` in detail

For an effectful call, everything after the effect runs on a detached 4 s budget. If any of it
fails, the coordinator returns `EffectsUnrecordedError` with the effect record refs. The effect's
intent, outcome and record are already durable and the budget was debited exactly once. The
call did **not** commit a log entry. Do not assume the effect did not happen: resolve each
listed record with `GET /v1/objects/{ref}?payload=true`. There is no automatic retry, because
re-planning against a new head would be a different invocation.

### MCP (`/mcp/`)

Protocol-level problems get JSON-RPC errors from the upstream handler, for example
`-32601 method not found: "…"`, `-32602 unknown tool "…"` (also what an unauthorized tool name
gets, since the session's surface does not contain it) and `-32602 invalid params: …`.

When the host's invocation fails for **any** reason — R8, an `EffectsUnrecorded` failure, a
deadline, a coordinator error — the upstream handler answers the **whole POST** with a fixed
envelope (plain JSON, not SSE):

```json
{"jsonrpc":"2.0","id":2,"error":{"code":-32603,"message":"host callback failed"}}
```

The message is `host callback timed out` for a deadline, `host callback canceled` for a
cancellation, or `host callback capacity exceeded` when all callback slots are busy. The
specific cause, including effect record refs, is **not** on the MCP wire; it is in the daemon's
stderr and in the store (the effect journal and records).

:::warning EffectsUnrecorded is only named over A2A
The effect-record refs of an `EffectsUnrecordedError` reach the caller over `/a2a/`. Over
`/mcp/` the caller sees only `host callback failed`. After such a failure, check the log and the
episode's effect records before retrying.
:::
