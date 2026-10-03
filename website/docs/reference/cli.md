---
title: CLI reference
sidebar_position: 1
description: Every ailang-worldd and world-publish command and flag, taken from the binaries' own help output.
---

# CLI reference

Two programs, both built from this repository. Help text below is the binaries' own output,
captured from a build of `dev` on 2026-10-03.

## `ailang-worldd`

```text
ailang-worldd — AILANG World local daemon (loopback only)

Usage:
  ailang-worldd serve --db <path> [--bind host:port] [--ailang-bin <path>]
                      [--workspace-root <dir> --tool-ailang-bin <path>]
                      [--examples-dir <dir>]
  ailang-worldd [--addr <url>] health
  ailang-worldd [--addr <url>] head
  ailang-worldd [--addr <url>] world get <ref>
  ailang-worldd [--addr <url>] object get <ref> [--payload]
  ailang-worldd [--addr <url>] object find <semanticId> [--after <ref>] [--limit N]
  ailang-worldd [--addr <url>] log get <index>
  ailang-worldd [--addr <url>] log range --from N [--limit M]
  ailang-worldd [--addr <url>] registry get <name>
  ailang-worldd [--addr <url>] commit --file <commit.json>

Global client flag:
  --addr <url>   base URL of the daemon (default http://127.0.0.1:7644).
                 Applies to client verbs only; it is NOT a 'serve' flag.

serve flags:
  --db <path>          world store database (required)
  --bind host:port     loopback listen address (default 127.0.0.1:7644);
                       a non-loopback host is refused — there is no override
  --ailang-bin <path>  interpreter to archive and pin at startup (optional)
  --workspace-root <dir>
                       episode worktrees live at <dir>/<episode> (made by the
                       operator with git worktree add before session mint); it
                       must not contain the store, its archive, its rendered
                       policies or its tool cache — startup refuses otherwise
  --tool-ailang-bin <path>
                       AILANG binary the workspace tools run (must be
                       AILANG v0.51.0); archived and hash-verified like
                       --ailang-bin. The Workspace.*/Ailang.* tools are served
                       only when both this and --workspace-root are set
  --examples-dir <dir> AILANG examples corpus examples-search reads (passed
                       to the tool as AILANG_EXAMPLES; not built into the
                       binary). Default: ~/.ailang/examples when it exists,
                       else examples-search refuses "no examples corpus
                       configured". Must be outside --workspace-root

Exit codes: 0 ok, 1 usage or client error, 2 fatal startup.
```

`ailang-worldd help` prints the same text. There is no `--help` per verb except where Go's flag
parser prints a flag list (shown below).

:::note `session` is not in the usage text
`session mint` and `session revoke` exist but are missing from the top-level usage above.
`--addr` is refused with them, as with `serve`, because they act on `--db` directly, not on a
running daemon.
:::

Every invocation refuses to start (exit 2) if `AILANG_REGISTRY_API_KEY` is set in its
environment.

### `serve`

See the flags in the usage text above, and [Operating the daemon](../guides/operating-the-daemon.md).

### Client verbs

All client verbs send one bounded request (30 s client timeout) to `--addr` and print the
response body on success. On an HTTP error they print
`ailang-worldd: <path> returned HTTP <status> <Class>: <message>` to stderr (with
`(observedHead=… selectedHead=…)` for a conflict) and exit 1.

| Verb | Request |
|---|---|
| `health` | `GET /v1/health` |
| `head` | `GET /v1/head` (prints the ref as plain text) |
| `world get <ref>` | `GET /v1/worlds/{ref}` |
| `object get <ref> [--payload]` | `GET /v1/objects/{ref}`, with `?payload=true` when `--payload` is given |
| `object find <semanticId> [--after <ref>] [--limit N]` | `GET /v1/objects/by-semantic-id/{semanticId}` |
| `log get <index>` | `GET /v1/log/{index}` |
| `log range --from N [--limit M]` | `GET /v1/log?from=N&limit=M` |
| `registry get <name>` | `GET /v1/registry/{name}` |
| `commit --file <commit.json> [--session <token>]` | `POST /v1/commit` |

Flags come **after** the positional argument for `object get`, `object find`, and `log range`.

#### `commit`

```text
Usage of ailang-worldd commit:
  -file string
    	commit JSON file
  -session string
    	session credential (64-hex Bearer token) for the session-gated /v1/commit
```

`--session` is sent as `Authorization: Bearer <token>`. Without it the daemon answers
`401 SessionAbsent`. (The top-level usage line omits `--session`.)

### `session mint`

```text
Usage of ailang-worldd session mint:
  -db string
    	world store database (required)
  -episode string
    	episode the session binds to (required)
  -grant value
    	EFFECT=SCOPE:BUDGET, repeatable, at least one
  -out string
    	write the credential ONCE to <path> at mode 0600 (default: print to stdout exactly once)
  -ttl int
    	lifetime in whole seconds (default 3600) (default 3600)
```

Attended: opens `/dev/tty` and asks `Confirm mint for episode … [y/N]` there; refuses without
a controlling terminal. The daemon must be stopped (single writer). Each grant's expiry is set to
the session's (`now + ttl`). With `--out` it prints
`minted session credential for episode <ep>: <n> grant(s), expires epoch <t>, written to <path>`;
it always prints `credential_id=<hash>` to stderr. See
[Sessions and episodes](../concepts/sessions-and-episodes.md).

### `session revoke`

```text
Usage of ailang-worldd session revoke:
  -db string
    	world store database (required)
```

```bash
ailang-worldd session revoke --db <path> <credential_id>
```

`--db` must come **before** the ID; the command's own usage message
(`session revoke <credential_id-hash> [--db <path>]`) shows the other order, which fails. The ID
is the 64-hex `credential_id` printed at mint. Not attended; the daemon must be stopped.

## `world-publish`

```text
world-publish — the attended entrypoint for an IRREVERSIBLE public write

  world-publish packet     [--package-dir D] [--golden G]
  world-publish approve    --store S [--registry-origin O] --now N --expires E
  world-publish publish    --store S --registry-origin O --publisher P \
                           --credential-file C --approval-ref R --now N --expires E (--live | --dry-run)
  world-publish reconcile  --store S [--registry-origin O] [--probe]
  world-publish transitions --store S --manifest M (--ailang-bin B | --interpreter-ref R)  [local registry write]

Exit codes: 0 done · 1 failed · 2 usage · 3 STOP (a fence refused; nothing happened)
```

This usage is printed when the program is run with no verb or an unknown one
(`world-publish --help` prints `unknown verb "--help"` followed by it).

### Verbs

| Verb | Kind | What it does |
|---|---|---|
| `packet` | Read-only, headless-permitted | Recomputes the `world/core` ready-packet from `--package-dir` and compares every field with `--golden`. `STOP fence=packet reason=drift` on any difference. |
| `approve` | **Attended** | Mints a one-shot `ApprovalDecisionV1` for the current packet and prints its ref. |
| `publish` | **Attended**; irreversible with `--live` | Spends the approval, exactly once. Exactly one of `--live` or `--dry-run`. `--dry-run` runs every fence and sends nothing. |
| `reconcile` | Read-only, headless-permitted | Lists publish intents with no outcome; `--probe` issues read-only metadata GETs to resolve them. `--live` is refused. |
| `transitions` | **Attended**; local registry write | Publishes a descriptor manifest into the store's transition registry. `--live` is refused. See [Transitions manifest](transitions-manifest.md). |

### Flags

All verbs share one flag set. Each verb uses the subset shown in the usage above. `-h` on any
verb prints this list (and exits 2):

```text
  -ailang-bin string
    	transitions: interpreter binary to archive now, as the daemon does at startup
  -approval-ref world-publish approve
    	the ApprovalDecisionV1 hashref minted by world-publish approve
  -credential-file string
    	file outside the working tree holding the registry API key
  -decided-by string
    	who is granting the approval
  -dry-run
    	rehearse the publish: every fence, no request
  -episode string
    	durable episode ID (default "attended-publish")
  -expires int
    	logical time the approval expires at
  -golden string
    	the committed ready-packet golden (default "scripts/world_package_ready_packet.golden.json")
  -interpreter-ref string
    	transitions: archived interpreter HashRef pinned into every descriptor
  -live
    	PERFORM THE IRREVERSIBLE PUBLIC WRITE
  -manifest string
    	transitions: descriptor manifest file (JSON array)
  -now int
    	logical time of this act
  -package-dir string
    	the projected package directory (default "packages/world-core")
  -probe
    	reconcile: issue the read-only metadata GETs
  -publisher string
    	path to the pinned released ailang binary
  -registry-origin string
    	the read-only public bucket origin
  -requester string
    	who is requesting the approval
  -store string
    	path to the world database this publish is recorded in
```

There is deliberately no flag that sets the registry validator origin; it is a compiled
constant. A test compares the flag set to a frozen list, so adding a flag fails the build.

### STOP lines

A refusal prints `STOP fence=<name>[ reason=<reason>]` and exits 3. Fence names: `mode`, `ci`,
`tty`, `confirmation`, `store`, `approval`, `credential`, `packet`, `handler`. See
[Attended steps](../guides/attended-steps.md#how-the-fences-work).

Build `world-publish` to a binary for attended use rather than using `go run`: `go run` exits 1
for a child that exited 3, which hides the STOP contract.
