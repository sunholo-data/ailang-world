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
- part of every invocation ID: tool calls are recorded as `mcp:<episode>:<task-id>` (MCP) or `a2a:<episode>:<task-id>` (A2A).

For the coding tools the episode ID must match `^[a-z0-9][a-z0-9-]{0,63}$`, and
`<workspace-root>/<episode>` must be a real directory (not a symlink) that the operator created
beforehand, normally a git worktree made by `session new` (below) or `git worktree add`.
Anything else gives the episode an empty handler registry, so every effectful call is refused
before anything runs.

An episode directory that is not a git worktree can still be adopted, with
`session new --allow-non-git`, as a **degraded** episode: it has no own branch, no branch diff
or provenance trailer and no reset to a base, and its `worktree_delta` records use a bounded
file walk instead of `git status`. `session list` marks it `DEGRADED`, and `/v1/health` lists it
under `degraded_episodes`.

## Session credentials

A **session** is a 64-hex bearer token bound to one episode, a set of grants and an expiry. It
is sent as a standard HTTP header:

```http
Authorization: Bearer <64-hex-credential>
```

No other header is read. The token is a session credential, never an API key.

The store keeps only a hash of the token (the `credential_id`). The raw token is printed or
written exactly once, at mint time.

### Provisioning an episode: `session new`

`session new` is the usual path (QUICKSTART §9). In one attended act it provisions the episode's
worktree (from `--repo`, detached or on `--branch`, unless it already exists) and mints a
session for it. It speaks directly to the store file, so the daemon must be **stopped** (the
store has one writer process at a time):

```bash
ailang-worldd session new ep1 --db /path/world.db --workspace-root /path/ws --repo /path/repo \
  --preset se-tools --ttl 14400 --out /path/session
```

- `--preset se-tools` gives one grant per effect the coding tools declare: ten grants, each
  scope `worktree` with `--budget` calls (default 50). `--grant EFFECT=SCOPE:BUDGET` adds a grant
  or overrides a preset one.
- The y/N prompt names the worktree and every grant. If the mint fails, the worktree this run
  created is removed.
- `--allow-non-git` adopts an existing non-git directory as a degraded episode (see above).

`session list --db /path/world.db` shows the credentials in a store (the `credential_id`, never
the token), live or expired, with their grants; it reads beside a running daemon.

### Minting

`session mint` mints a session with explicit grants for an episode whose directory already
exists. It is attended too, and also needs the daemon **stopped**:

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

`session new` and `session mint` open `/dev/tty` and ask a `y/N` confirmation there. Without a
controlling terminal they refuse. See [Attended steps](../guides/attended-steps.md).

### Revoking

```bash
ailang-worldd session revoke --db /path/world.db <credential_id>
ailang-worldd session revoke --db /path/world.db --episode ep1
```

Revocation deletes the mapping row, so the next request with that token is `SessionUnknown`.
`--episode` revokes every credential of that episode. It needs no terminal. Like mint, it opens
the store for writing, so stop the daemon first.

Flags and the `credential_id` may come in either order. An ID (or episode) with no credential
in the store is refused with exit 1 (`refusing: no session credential … nothing was revoked`),
and so is a `--db` that does not exist; it is never created.

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
