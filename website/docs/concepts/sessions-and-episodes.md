---
title: Sessions and episodes
sidebar_position: 4
description: Bearer session credentials, the episode they bind to, how they are minted and revoked, and what they gate.
---

# Sessions and episodes

## Episodes

An **episode** is the unit of work an agent session belongs to. It names:

- the budget ledger (spend is summed per episode, see
  [budgets](effects-capabilities-budgets.md#budgets));
- for the coding tools, the worktree `<workspace-root>/<episode>` the tools act in;
- part of every invocation ID: tool calls are recorded as `a2a:<episode>:<task-id>`.

For the coding tools the episode ID must match `^[a-z0-9][a-z0-9-]{0,63}$`, and
`<workspace-root>/<episode>` must be a real directory (not a symlink) that the operator created
beforehand, normally with `git worktree add`. Anything else gives the episode an empty handler
registry, so every effectful call is refused before anything runs.

## Session credentials

A **session** is a 64-hex bearer token bound to one episode, a set of grants and an expiry. It
is sent as a standard HTTP header:

```http
Authorization: Bearer <64-hex-credential>
```

No other header is read. The token is a session credential, never an API key.

The store keeps only a hash of the token (the `credential_id`). The raw token is printed or
written exactly once, at mint time.

### Minting

Minting is an attended act. It speaks directly to the store file, so the daemon must be
**stopped** (the store has one writer process at a time):

```bash
ailang-worldd session mint --db /path/world.db --episode ep1 \
  --grant Workspace.Read=worktree:50 --grant Workspace.Write=worktree:50 \
  --ttl 14400 --out /path/session
```

- `--grant EFFECT=SCOPE:BUDGET` is repeatable; at least one is required. The budget is the last
  `:`-separated segment, so a scope may contain `:`.
- `--ttl` is the lifetime in whole seconds (default 3600). Every grant expires with the
  session.
- `--out` writes the token to a file at mode `0600` and prints a one-line confirmation;
  without it the token is printed to stdout once.
- The `credential_id` is printed to stderr for later revocation.

`session mint` opens `/dev/tty` and asks a `y/N` confirmation there. Without a controlling
terminal it refuses. See [Attended steps](../guides/attended-steps.md).

### Revoking

```bash
ailang-worldd session revoke --db /path/world.db <credential_id>
```

Revocation deletes the mapping row, so the next request with that token is `SessionUnknown`.
It needs no terminal. Like mint, it opens the store for writing, so stop the daemon first.

Put `--db` **before** the `credential_id`. The command's own usage line shows
`session revoke <credential_id-hash> [--db <path>]`, but flags after the ID are not parsed and
that order fails with the usage message. Revoking an ID that does not exist still prints
`revoked session credential …` and exits 0.

## What a session gates

| Surface | Without a valid session |
|---|---|
| `POST /v1/commit` | `401 SessionAbsent` / `SessionUnknown` / `SessionExpired`, or `400 InvalidSession` for a malformed header |
| `GET /.well-known/agent.json` | Same REST error envelope and status mapping |
| `POST /a2a/` | JSON-RPC error `-32001` (absent, unknown, expired) or `-32600` (malformed) |
| `POST /mcp/` | HTTP `401` with a plain-text message |
| All other `GET /v1/*` routes and `/workbench` | **Open.** Reads are not session-gated today (residual R1 of the session-authority design) |

The commit route only checks that a valid session exists; it does not check grants. Grants
govern which transitions a session can see and call, and what their effects may do.

:::warning Reads are unauthenticated
Anything that can reach the loopback port can read every object, log entry and world in the
store. Keep the daemon on loopback (it refuses anything else) and treat the machine as the
trust boundary.
:::
