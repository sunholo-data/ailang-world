---
title: Attended steps
sidebar_position: 2
description: Which operations need a human at a real terminal, why, how the TTY fence works, and how to run them from an IDE or agent harness.
---

# Attended steps

Some operations grant authority or write something that cannot be taken back. World fences
them to a human at a controlling terminal. An autonomous process — CI, a script, an agent — is
structurally unable to complete them. That is the design, not a defect to route around.

| Operation | Command | Fence | Why it is attended |
|---|---|---|---|
| Provision an episode | `ailang-worldd session new` | Opens `/dev/tty`, `y/N` confirmation (naming the worktree and every grant) read from it, before any side effect | Provisioning a session credential is one human act; otherwise whoever runs the script gets authority |
| Mint a session | `ailang-worldd session mint` | Opens `/dev/tty`, `y/N` confirmation read from it | Minting a credential is one human act; otherwise whoever runs the script gets authority |
| Publish transitions | `world-publish transitions` | CI tripwire, controlling terminal, typed phrase | The registry write path is the operator's hands, never an agent's |
| Open a new semantics epoch | `world-publish epoch nominate --new-epoch` | CI tripwire, controlling terminal, typed phrase naming the epoch, release and store | A release whose semantics moved changes how history replays; that is a human's decision |
| Mint a package-publish approval | `world-publish approve` | CI tripwire, controlling terminal, typed phrase | An approval for an irreversible act must not be self-minted by a loop |
| Publish `world/core` | `world-publish publish --live` | All of the above, plus a one-shot approval, golden packet and credential file | The public registry is immutable |

Not attended: `session list`, `session revoke`, `world-publish packet`,
`world-publish reconcile`, `world-publish epoch nominate --candidate` (headless-permitted, per
D-WORLD-78: it adds a release whose replay goldens are byte-identical to the current epoch) and
every `ailang-worldd` client verb.

## How the fences work

### `session new` and `session mint`

Both open `/dev/tty` read-write, write the prompt there and read the answer from there. If
`/dev/tty` cannot be opened they refuse, `session new` before it creates a worktree or a row:

```text
ailang-worldd session new: refusing: no controlling terminal (open /dev/tty: device not configured); provisioning a session credential requires one human act at a terminal. fix: run it yourself from a terminal window (an IDE terminal pane works); an agent cannot run this step
ailang-worldd session mint: refusing: no controlling terminal (open /dev/tty: device not configured); minting a session credential requires one human act at a terminal. fix: run it yourself from a terminal window (an IDE terminal pane works); an agent cannot run this step
```

Anything but `y` or `yes` aborts (`aborted (not confirmed)`).

### `world-publish`

Every refusal prints `STOP fence=<name> reason=<reason>` on stderr and exits **3**: nothing
happened and nothing will. The attended verbs (`approve`, `publish`, `transitions`,
`epoch nominate --new-epoch`) check, in order:

1. **`ci`** — refuses when `CI` or `GITHUB_ACTIONS` is set. This is a declared tripwire, not the
   fence: `env -u CI` defeats it.
2. **`tty`** — the load-bearing fence. `reason=no-controlling-terminal`: `/dev/tty` does not
   open, which is what an agent harness, `setsid`, cron or CI sees. It is the only `tty`
   refusal.
3. **`confirmation`** — prints the prompt and reads one line **from the controlling terminal**,
   then compares it with the exact phrase (`reason=eof` when the terminal closes first,
   `reason=mismatch` for anything else). When stdin is that terminal, it is read; otherwise
   the command reads the `/dev/tty` it opened. A redirected stdin (`echo phrase |`,
   `< /dev/null`) is never read as the confirmation.

Every refusal's last line is the fix.

The store is opened **before** the fence, so a mistyped path, or a daemon still holding the
writer lock, is reported before you are asked to type anything. `transitions` also reads its
manifest before the fence, so that the phrase can name the write.

### The confirmation phrase

Each verb's phrase names the write it confirms. `approve` and `publish` name the current
`world/core` candidate:

```text
Type exactly, to proceed with an IRREVERSIBLE public write:
  publish world/core@0.1.1 irreversibly
>
```

`transitions` is a **local** registry write, not a public one, and a later revision can
supersede it. Its phrase names the descriptor count and the store; for the se-tools manifest:

```text
Type exactly, to write a LOCAL transition-registry revision (10 descriptor(s) from packages/se-tools/transitions.json):
  publish 10 transitions to /tmp/se-world/world.db
>
```

`epoch nominate --new-epoch` asks for `open epoch <N> for <release> in <store>`: the epoch it
opens, the `--release` and the `--store`. Before the fence it has
already checked the release against the compiled-in manifest `pins/ailang-releases.json`, the
`--gate-verdict` and the `--binary` digest, so a refused nomination never asks you to type.

The `tty` and `ci` fences in front of them are the same as for the package publish. Type the
phrase exactly as printed.

## Running attended steps from an IDE pane or agent harness

Embedded terminals often give the process a stdin that is not the controlling terminal. That
no longer matters: `world-publish`, `session new` and `session mint` all open `/dev/tty`
themselves and prompt and read there, so an IDE terminal pane works with the plain command:

```bash
/tmp/world-publish transitions --store /tmp/se-world/world.db \
  --manifest packages/se-tools/transitions.json --ailang-bin $PIN
```

Until 2026-10-06 the publish refused such a stdin
(`reason=stdin-is-not-the-controlling-terminal`) and the advice was to append `< /dev/tty`.
Reading `/dev/tty` directly is exactly as strong: the line comes from the same device either
way, and you still type it. `< /dev/tty` still works. If `/dev/tty` itself does not open (a
true headless process: an agent harness, cron, CI), there is no way through, by design.

## The order of operations

Both `session new` (or mint) and publish need single-writer authority over the store, so the
daemon must be stopped. A typical attended session for the coding tools:

1. `serve` once with `--ailang-bin` to bootstrap the epoch registry, then stop it.
2. **Publish** the transitions (attended).
3. **Provision** the episode and its session with `session new` (attended).
4. `serve` with the tools enabled.
5. Commit genesis with `--session`, then hand the session to agents.

`tools/attended/se_smoke.sh` follows this order and refuses `publish` or `mint` while its daemon
is running.

## The `world/core` package publish

Publishing `world/core` to the public AILANG registry is the one irreversible effect in the
repository. Its runbook is `docs/SELF_MOD_PUBLISH.md`; the helper is
`tools/attended/self_mod_publish.sh`, which keeps each human act a separate subcommand
(`preflight`, `approve`, `dry-run`, `live`). In outline:

1. **Readiness (automated, no public write).** Build the projection
   (`./scripts/build_world_package.sh`), run the readiness gate
   (`./scripts/verify_world_package.sh`, which needs the v0.41.0 compiler via
   `WORLD_PKG_AILANG_BIN`), and run the full repo gate.
2. **Review** the golden ready-packet and confirm no drift: `world-publish packet` recomputes it
   from `packages/world-core/` and compares all fields (read-only; `STOP fence=packet reason=drift`
   on any difference).
3. **Approve** (attended): `world-publish approve --store … --registry-origin … --now N --expires E`
   mints a one-shot `ApprovalDecisionV1` and prints its ref.
4. **Rehearse**: `world-publish publish --dry-run …` runs every fence and sends nothing.
5. **Publish** (attended, irreversible): `world-publish publish --live … --approval-ref <ref>`.
   Outcomes: `PUBLISHED` (exit 0), `FAILED` (exit 1, known and recorded), or
   `INDETERMINATE` (exit 1): **do not retry**.
6. **Reconcile** an indeterminate attempt: `world-publish reconcile --store …`, adding `--probe`
   for read-only metadata GETs. It resolves to `succeeded-reconciled`, `conflict`,
   `not-published` (a retry then needs a new attended approval) or `probe-unavailable` (stop; a
   human decides).

Minting and spending are two invocations on purpose, so two human acts never collapse into one
keystroke. The approval is single-use durably (a primary key in the store), even across
restarts. The credential is read only from a mode-`0600` file outside the working tree
(`--credential-file`). Build `world-publish` into a fresh temp directory for this procedure; the
helper always rebuilds and never trusts a prebuilt binary, and `go run` is avoided because it
does not propagate exit code 3.

`world/core@0.1.1` was published on 2026-10-02 from a fresh store
(`~/.ailang/world/world-0.1.1.db`). The legacy store that recorded 0.1.0 is schema version 2,
which the current binary refuses to open; it must never be modified or deleted.
