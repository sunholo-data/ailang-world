---
title: Operating the daemon
sidebar_position: 1
description: Running ailang-worldd serve — flags, the single-writer store lock, the registry-credential guard, startup refusals, timeouts, logs and shutdown.
---

# Operating the daemon

`ailang-worldd serve` is one local process over one SQLite store file. It serves REST, MCP and
A2A on a loopback address.

## Flags

```bash
ailang-worldd serve --db <path> [--bind host:port] [--ailang-bin <path>]
                    [--workspace-root <dir> --tool-ailang-bin <path>]
                    [--examples-dir <dir>]
```

| Flag | Default | Meaning |
|---|---|---|
| `--db` | required | The world store database. Its parent directory must already exist. |
| `--bind` | `127.0.0.1:7644` | Loopback listen address. A non-loopback host is refused; there is no override. |
| `--ailang-bin` | none | The interpreter (`PIN`, v0.41.0) to archive and pin at startup. **Without it, transitions cannot be invoked**: the A2A and MCP surfaces answer `transition invocation is not available in this daemon`. |
| `--workspace-root` | none | Directory holding one worktree per episode, at `<dir>/<episode>`. Must not contain the store, its archive, its rendered policies or its tool cache. |
| `--tool-ailang-bin` | none | The AILANG binary the coding tools run (`TOOL`). Must be AILANG v0.52.1. Archived and hash-verified like `--ailang-bin`. |
| `--examples-dir` | `~/.ailang/examples` if it exists | The AILANG examples corpus `examples-search` reads, passed to the tool as `AILANG_EXAMPLES`. Must be outside `--workspace-root`. |

The coding tools are served only when **both** `--workspace-root` and `--tool-ailang-bin` are
set. With one but not the other, startup logs:

```text
ailang-worldd: workspace tools disabled: --workspace-root and --tool-ailang-bin must both be set; every effect-declaring transition is refused
```

`--addr` is a **client** flag. It is refused with `serve` (use `--bind`).

Exit codes: `0` ok, `1` usage or client error, `2` fatal startup.

## What startup does

In order, within a 9 s startup budget:

1. Refuses if `AILANG_REGISTRY_API_KEY` is set (see below).
2. Refuses a non-loopback `--bind` host (allowed: `127.0.0.1`, `::1`, `localhost`).
3. Resolves `--workspace-root` and `--examples-dir` and refuses a workspace root that contains
   the daemon's own state, or an examples corpus inside the workspace root.
4. Opens the store and takes the writer lock.
5. Archives `--tool-ailang-bin` (if given) and refuses any release but `AILANG v0.52.1`.
6. Archives `--ailang-bin` (if given) and builds the invocation coordinator.
7. Bootstraps the epoch registry (`world/epoch-registry/v1`) for the interpreter's release.
   A registry head naming different bytes is a fatal divergence, never silently rewritten.
8. Runs a bounded integrity scan of log and world rows.
9. Binds the socket and prints `ailang-worldd listening on http://127.0.0.1:7644`.

A startup failure prints `daemon startup failed at <stage>: <detail>: <cause>` and exits 2. The
stages are `config`, `store-open`, `archive`, `registry-bootstrap`, `bind-policy` and `listen`.

For each episode the first tool call constructs the episode's handler: it renders the policy at
`<db-dir>/policies/<episode>.toml` (mode `0600`), creates the tool cache at
`<db-dir>/cache/<episode>`, and runs one `policy-tool` summary to verify the binary honours the
policy. If that fails, the daemon logs
`ailang-worldd: workspace tools unavailable for episode "<ep>": …` and the episode gets no
tools.

## One writer per store

Each store file permits exactly one writer process. `Open` takes a non-waiting exclusive OS lock
on `<db>.writer.lock` before any SQLite write handle. A second writer fails immediately:

```text
ailang-worldd: daemon startup failed at store-open: another process already holds writer authority for this database (single-writer is enforced, not conventional): …
```

This covers every writer, not only a second daemon. **Stop the daemon before** running
`session mint`, `session revoke`, `world-publish transitions` or any other `world-publish`
verb that opens the same store. Read-only store users may coexist.

The operator sequence is therefore always: serve, stop, do the attended step, serve again.

## The registry-credential guard

```text
ailang-worldd: broker: AILANG_REGISTRY_API_KEY is set in the process environment: move it to a mode-0600 file outside the working tree and unset it (see design_docs/planned/w-self-mod-vertical.md Decision 4)
```

The check runs before flag parsing, so it applies to `serve` and to every client verb. The
message names the variable, never its value. Fix it with `unset AILANG_REGISTRY_API_KEY` in the
shell that starts the daemon. The publish credential belongs in a file named by
`world-publish --credential-file`, and nowhere else.

## Logs

- **stdout**: the `listening on` line, then integrity warnings only if the scan found a hole or
  could not finish (`integrity_hole …`, `integrity_scan_incomplete …`). A clean scan prints
  nothing.
- **stderr**: one line per internal error with its route (the HTTP client sees only
  `internal store failure`), A2A and MCP refusal causes
  (`ailang-worldd: a2a refusal: tasks/send <invocation-id>: "<cause>"` and
  `ailang-worldd: mcp refusal: tools/call <invocation-id>: "<cause>"`), and workspace-tool
  messages.

Run `serve` in a terminal you can see, or redirect both streams to a file. When it is
backgrounded, stop it by PID (the smoke script keeps a pidfile), not by shell job control from
another shell.

## Bounds

| Bound | Value |
|---|---|
| Store read per GET request | 10 s, then `503 Timeout` |
| Commit store work | 3 s, then `503 Timeout` or `503 CommitUncertain` |
| Credential lookup | 3 s |
| One transition invocation (MCP or A2A) | 20 s |
| Plan phase / finish phase | 2 s each |
| Tool handler | at most 10 s; the rendered policy's own `timeout_ms` is 8 s |
| Post-effect writes and commit | 4 s, detached from the caller's deadline |
| HTTP read / write timeout | 30 s / 30 s |
| Commit body | 8 MiB (`413 PayloadTooLarge`) |
| A2A request body | 1 MiB |
| MCP request body | 4 MiB |
| Log range and object pages | default 100, maximum 500 |
| Shutdown drain | 10 s |

## Health

```bash
ailang-worldd health
```

returns `status`, `daemon_version`, `db_path`, `interpreter_ref` and `interpreter_version`
(`interpreter_ref` is empty when `serve` ran without `--ailang-bin`).

## The workbench

`GET /workbench` is a read-only HTML renderer for operators: the log, worlds, objects and the
references between them. It accepts the query parameters `world`, `object`, `from`, `entry`,
`payload` (`0` or `1`), `refsAfter` and `commitsAfter`. It is not part of the frozen v1 API and
its HTML may change.

## Shutdown

`SIGINT` or `SIGTERM` drains in-flight requests (bounded at 10 s) and releases the writer lock.
A clean drain prints nothing.
