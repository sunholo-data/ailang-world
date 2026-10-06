---
title: HTTP API
sidebar_position: 2
description: Every route ailang-worldd serves — method, authentication, request and response shapes, and error classes — taken from the route table and handlers.
---

# HTTP API

The daemon listens on loopback only (default `http://127.0.0.1:7644`). The route table is
`Daemon.Handler` in `host/daemon/daemon.go`. Routes use Go method patterns, so a wrong method
gets `405 Method Not Allowed` from the router.

| Method | Path | Auth | Purpose |
|---|---|---|---|
| `GET` | `/v1/health` | none | Daemon status and interpreter pin |
| `GET` | `/v1/head` | none | Selected world head (plain text) |
| `GET` | `/v1/worlds/{ref}` | none | One world row |
| `GET` | `/v1/objects/{ref}` | none | One object, payload optional |
| `GET` | `/v1/objects/by-semantic-id/{name...}` | none | Objects with a given semantic ID, paged |
| `GET` | `/v1/log/{index}` | none | One log entry |
| `GET` | `/v1/log` | none | A page of log entries |
| `GET` | `/v1/registry/{name...}` | none | A registry head |
| `POST` | `/v1/commit` | **session** | Commit a world transition (the only REST mutation) |
| `GET` | `/v1/receipts/{id}` | none | Receipt of a `rest:` commit |
| `GET` | `/workbench` | none | Read-only HTML operator renderer |
| `GET` | `/.well-known/agent.json` | **session** | Session-scoped A2A agent card |
| `POST` | `/a2a/` | **session** | A2A JSON-RPC (`tasks/send`) |
| `POST` | `/mcp/` | **session** | MCP Streamable HTTP (`initialize`, `ping`, `tools/list`, `tools/call`) |

The ten `/v1` routes are the frozen v1 machine table, pinned against a checked sketch by a test.
`/workbench`, the agent card, `/a2a/` and `/mcp/` are additive.

**Auth** means `Authorization: Bearer <64-hex-credential>` holding a live session (see
[Sessions and episodes](../concepts/sessions-and-episodes.md)). The `GET /v1/*` routes are not
authenticated today.

## Conventions

**Refs** are canonical `sha256:<64 hex>` text. A malformed ref is `400 BadRequest` naming the
field.

**Errors** on the `/v1` routes and the agent card use one envelope:

```json
{"error":{"class":"NotFound","message":"log entry not found"}}
```

| Class | Status | Meaning |
|---|---|---|
| `BadRequest` | 400 | The client's input is malformed; the message names what (not sanitized) |
| `InvalidSession` | 400 | Malformed `Authorization` header |
| `SessionAbsent`, `SessionUnknown`, `SessionExpired` | 401 | No, unknown or expired session |
| `NotFound` | 404 | The thing is not there (distinct from malformed) |
| `HeadConflict` | 409 | Stale `observedHead`; carries `observedHead` and `selectedHead` |
| `PayloadTooLarge` | 413 | Commit body over 8 MiB |
| `Internal` | 500 | Always the fixed message `internal store failure`; the cause goes to the daemon's stderr |
| `Timeout` | 503 | A read (10 s), commit (3 s) or credential lookup (3 s) deadline ended |
| `CommitUncertain` | 503 | The commit's outcome is unknown; reconcile before resending |
| `LookupIndexUnavailable` | 503 | The semantic-ID index is absent or incompatible |
| `ProjectionUnavailable` | 503 | Agent card: a registry read failed; retry |
| `ProjectionDeadlineExceeded` | 504 | Agent card: the bounded read deadline was hit |

**Paging.** `limit` defaults to 100 when absent or below 1, and is capped at 500.

## `GET /v1/health`

```json
{"status":"ok","daemon_version":"0.1.0","db_path":"/tmp/world-demo.db",
 "interpreter_ref":"sha256:…","interpreter_version":"AILANG v0.41.0\nCommit: 24ee108\n…"}
```

`interpreter_ref` and `interpreter_version` are empty when `serve` ran without `--ailang-bin`.

## `GET /v1/head`

`200 text/plain`: the selected head ref, for example `sha256:…`.
`404 NotFound` `no world head has been selected yet` before genesis.

## `GET /v1/worlds/{ref}`

```json
{"ref":"sha256:…","revision":0,"stateRoot":"sha256:…","logHead":"sha256:…"}
```

`404 NotFound` `world not found`.

## `GET /v1/objects/{ref}`

```json
{"hash":"sha256:…","interfaceHash":"sha256:…","semanticId":"world/effect-record/v1",
 "provenance":"…"}
```

Add `?payload=true` to include `"payload":"<base64>"`. Without it the key is absent, so object
reads stay bounded. `404 NotFound` `object not found`.

## `GET /v1/objects/by-semantic-id/{name...}`

The name may contain slashes. Query: `after=<ref>` resumes strictly after that hash;
`limit=N`.

```json
{"items":[{"hash":"sha256:…","interfaceHash":"sha256:…","semanticId":"world/invocation-record/v2","provenance":"coordinator:mcp"}]}
```

Ascending hash order, never payloads. A name no object carries is `200` with `"items":[]`, never
404. `400` for an empty name, a malformed `after` or a non-integer `limit`;
`503 LookupIndexUnavailable` if the index is missing.

## `GET /v1/log/{index}`

```json
{"header":{"entryIndex":1,"semanticsEpoch":1,"transitionFn":"sha256:…",
           "interpreter":"sha256:…","prevEntryHash":"sha256:…","writtenBy":"coordinator:mcp"},
 "entryHash":"sha256:…","transitionRef":"sha256:…"}
```

`400` `log index must be a non-negative integer`; `404 NotFound` `log entry not found`.

## `GET /v1/log?from=N&limit=M`

```json
{"items":[ {…log entry…}, … ]}
```

`from` defaults to 0 and must be a non-negative integer. The page stops at the first absent index,
so `from` past the end returns `"items":[]`. The whole page shares one 10 s read deadline.

## `GET /v1/registry/{name...}`

```json
{"name":"world/epoch-registry/v1","head":"sha256:…"}
```

Names in use: `world/epoch-registry/v1` (bootstrapped at startup),
`world/transition-registry/v1` (after a transitions publish), `world/approvals/v1`.
`404 NotFound` `registry head not found`.

## `POST /v1/commit`

**Requires a session.** Body: one JSON object; unknown fields and trailing values are refused.

```json
{
  "invocationId": "rest:<your-id>",
  "observedHead": "",
  "objects": [
    {"hash": "sha256:…", "interfaceHash": "sha256:…", "semanticId": "…",
     "provenance": "…", "payload": "<base64>"}
  ],
  "nextWorld": {"ref": "sha256:…", "revision": 0, "stateRoot": "sha256:…", "logHead": "sha256:…"},
  "entry": {
    "header": {"entryIndex": 0, "semanticsEpoch": 1, "transitionFn": "sha256:…",
               "interpreter": "sha256:…", "prevEntryHash": "sha256:…", "writtenBy": "…"},
    "entryHash": "sha256:…",
    "transitionRef": "sha256:…"
  }
}
```

- `invocationId` is optional and must start with `rest:`. With it, the commit's intent is
  journaled first and a receipt is available; resending the same commit with the same ID is
  idempotent, and reusing the ID for a different commit is `400`.
- `observedHead` is empty only for genesis.
- Each object's `hash` must be `sha256:` of its payload bytes; the store verifies it.

Success: `200 {"selectedHead":"sha256:…"}`.

| Failure | Response |
|---|---|
| No or bad session | `401`/`400` session classes |
| Body over 8,388,608 bytes | `413 PayloadTooLarge` |
| Invalid JSON, unknown field, bad ref | `400 BadRequest` naming the problem |
| Stale `observedHead` | `409 HeadConflict` with `observedHead`, `selectedHead` |
| Budget ended before the durable step | `503 Timeout`: not committed; safe to resend |
| Budget ended after the durable step began | `503 CommitUncertain`: reconcile via `GET /v1/receipts/{id}` (with an `invocationId`) or `GET /v1/log/{entryIndex}` before any resend |
| Payload/hash mismatch, other store failures | `500 Internal` |

## `GET /v1/receipts/{id}`

Only IDs starting with `rest:` (otherwise `400` `receipt id must be rest:<id>`).

```json
{"invocationId":"rest:demo-1","state":"resolved","resultRef":"sha256:…"}
```

`state` is `resolved` (it landed; `resultRef` is the committed `nextWorld.ref`), `not-started`
or `indeterminate`. Treat `indeterminate` as not landed and fail-closed.

## `GET /workbench`

`200 text/html` with a strict Content-Security-Policy and `Cache-Control: no-store`. Query
parameters: `world`, `object`, `from`, `entry`, `payload` (`0` or `1`), `refsAfter`,
`commitsAfter`; an unknown or repeated parameter is `400`. Not part of the frozen API; its
HTML may change.

## `GET /.well-known/agent.json`

**Requires a session.** Returns the A2A agent card, with `skills` listing exactly the
transitions this session's grants admit, in registry byte order:

```json
{"name":"ailang-worldd","description":"…","url":"http://127.0.0.1:7644","version":"0.1.0",
 "capabilities":{"streaming":false,"pushNotifications":false,"stateTransitionHistory":false},
 "defaultInputModes":["application/json"],"defaultOutputModes":["application/json"],
 "skills":[{"id":"ailang-read","name":"Read (sandboxed)","description":"…","tags":[],"examples":[]}]}
```

No registry head yet means `200` with `"skills":[]`. Session denials use the REST error
envelope above; a failed registry read is `503 ProjectionUnavailable`, a deadline
`504 ProjectionDeadlineExceeded`.

## `POST /a2a/` and `POST /mcp/`

See [MCP and A2A](mcp-and-a2a.md).
