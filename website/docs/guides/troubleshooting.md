---
title: Troubleshooting
sidebar_position: 4
description: Real errors operators have hit running World, what each one means, and what to do.
---

# Troubleshooting

Every message below is quoted from the code or from a measured run.

## Startup and environment

### `AILANG_REGISTRY_API_KEY is set in the process environment`

```text
ailang-worldd: broker: AILANG_REGISTRY_API_KEY is set in the process environment: move it to a mode-0600 file outside the working tree and unset it (see design_docs/planned/w-self-mod-vertical.md Decision 4)
```

Every `ailang-worldd` command refuses (exit 2), client verbs included. Run
`unset AILANG_REGISTRY_API_KEY` in that shell. `world-publish` reads the key only from
`--credential-file`.

### `another process already holds writer authority for this database`

Another process has the store open for writing: usually your running daemon. `session mint`,
`session revoke` and `world-publish transitions` all need the daemon **stopped** first. The
lock file is `<db>.writer.lock`; do not delete it by hand while a process holds it.

### `resolve parent of "…": no such file or directory` / `STOP fence=store reason=unopenable`

The store refuses to create its parent directory. `mkdir -p` it first.

### `STOP fence=store reason=unopenable … has user_version 2 … binary requires 4; refusing to modify`

The store was written by an older binary and there is no migration. Use a fresh store file and
leave the old one untouched.

### `bind host "…" is not loopback`

`serve` accepts only `127.0.0.1`, `::1` or `localhost` in `--bind`. There is no override.

### `--addr is a client flag and is not valid for 'serve'`

Use `--bind` with `serve`. `--addr` (with a scheme, for example `http://127.0.0.1:7645`) goes
before a client verb: `ailang-worldd --addr http://127.0.0.1:7645 health`.

### `tool binary is "…"; only "AILANG v0.51.0" is measured for the workspace tools`

`--tool-ailang-bin` must be exactly AILANG v0.51.0. See [Install](../getting-started/install.md).

### `the workspace root is refused` / `… is inside --workspace-root`

The store, its archive, `policies/` or `cache/` directory, or the examples corpus lies inside
`--workspace-root`. An agent that can edit its own policy has no policy. Put the store in a
directory outside the workspace root.

## Attended steps

### `STOP fence=tty reason=stdin-is-not-the-controlling-terminal`

`world-publish` was run with stdin that is a character device but not `/dev/tty`, typically an
IDE pane or an agent harness. Append `< /dev/tty` and type the phrase yourself.

### `STOP fence=tty reason=no-controlling-terminal`

There is no terminal at all. A headless process cannot do this step; a human must.

### `refusing: no controlling terminal (…); minting a session credential requires one human act at a terminal`

Same, for `session mint`. Run it at a real terminal.

### `STOP fence=confirmation reason=mismatch`

The typed line is not the exact phrase. All three attended `world-publish` verbs, including
`transitions`, currently ask for `publish world/core@0.1.1 irreversibly`. See
[Attended steps](attended-steps.md#the-confirmation-phrase).

### `interpreter release "…" is nominated by NO epoch in world/epoch-registry/v1`

`world-publish transitions` was pointed at a store whose daemon has never run with this
interpreter. Start `serve --ailang-bin $PIN` once against the store, stop it, then publish.

### `… is not loadable under its pinned interpreter` / `LDR001`

The transition source failed `check` under the pinned interpreter. A common cause is importing
`world/*`: publishable modules must be self-contained and import only `std/*`.

### `session revoke: usage: session revoke <credential_id-hash> [--db <path>]`

Put `--db` before the credential ID: `session revoke --db /path/world.db <credential_id>`.

## Committing

### `401 SessionAbsent` on commit

```text
ailang-worldd: /v1/commit returned HTTP 401 SessionAbsent: a session credential is required: no Authorization Bearer header was present
```

`POST /v1/commit` is session-gated. Pass `--session "$(cat /path/session)"` to `commit`, or send
`Authorization: Bearer <token>`. This also means the genesis commit comes after `session mint`.
Older runbook text that commits genesis without a session no longer works.

The related classes: `400 InvalidSession` (malformed header), `401 SessionUnknown` (revoked or
never minted, or minted in a different store), `401 SessionExpired`.

### `409 HeadConflict`

The world moved since the head you observed. The body carries `observedHead` and
`selectedHead`; re-plan against the selected head.

### `503 CommitUncertain`

The commit budget ended after the durable step began. It landed whole or not at all. Do not
resend until you have read the log row at your entry's index (or the receipt, if you sent an
`invocationId`). The current head is not evidence either way.

### `500 internal store failure` on a commit you built yourself

A payload whose bytes do not hash to its stated `hash` is refused by the store, but the route
currently reports it as an internal failure rather than a client error (finding F-3 of the
value demonstration). Check every object `hash` is `sha256:` of the exact payload bytes. The
daemon's stderr has the real cause.

## MCP and the coding tools

### The `/mcp/` response is `event: message` / `data: …`

That is correct. `/mcp/` answers with SSE framing; the JSON-RPC object is on the `data:` line.
Send both `Content-Type: application/json` and
`Accept: application/json, text/event-stream`, or the request is refused with
`Accept must list both application/json and text/event-stream`.

### `tools/list` is empty, or a tool is missing

A session sees only tools whose `access` effect it holds a live grant for. Check the grants you
minted, that the session has not expired, and that the transitions are published in **this**
store (`ailang-worldd registry get world/transition-registry/v1`).

### `workspace tools disabled`

```text
ailang-worldd: workspace tools disabled: --workspace-root and --tool-ailang-bin must both be set; every effect-declaring transition is refused
```

Start `serve` with both flags. Until then every tool call is refused before any effect runs
(rule R8). Over A2A the refusal reads
`transition declares an effect this daemon has no handler for`; over MCP the whole request
answers `host callback failed`.

### `workspace tools unavailable for episode "…"`

The episode's handler could not be built: check that `<workspace-root>/<episode>` exists, is a
real directory and not a symlink, that the episode ID matches `^[a-z0-9][a-z0-9-]{0,63}$`, and
that the tool binary still verifies. The episode then gets an empty handler registry.

### `{"ok":false,"refused":"no examples corpus configured: …"}`

`examples-search` has no corpus. Run `$TOOL examples download` (fills `~/.ailang/examples`) or
start `serve` with `--examples-dir DIR`, outside the workspace root.

### `{"ok":false,"refused":"…"}` with `effects: []`

The tool's own plan refused the arguments: an unknown key, a missing required key, or a path
that is absolute, empty, contains `..` or starts with `-`. The refusal is committed.

### `denied:budget`

The episode has spent its grant. Budget persists per episode across restarts and across new
sessions for the same episode; use a new episode (and worktree) for a fresh budget.

### `transition invocation is not available in this daemon`

`serve` ran without `--ailang-bin`, so there is no invocation coordinator, or a registry read
failed. Check the daemon's stderr for an `a2a refusal` line with the cause.

### `effects were requested but the invocation was not confirmed committed; effect records: …`

Over A2A: the effect ran and was recorded, but the commit did not land (for example another
writer moved the head). The listed effect records are durable. Do not assume the effect did not
happen. See [MCP and A2A](../reference/mcp-and-a2a.md#error-mapping).

### pi or `claude` hangs when run from a script

Both read a non-TTY stdin as extra prompt input and wait for EOF. Give them `</dev/null`.

### pi shows the World tools only through a proxy

pi-mcp-adapter registers direct tools from its cache (`~/.pi/agent/mcp-cache.json`). Run one
interactive `pi --mcp-config …` session first to write it.
