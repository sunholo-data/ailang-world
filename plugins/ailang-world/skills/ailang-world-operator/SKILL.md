---
name: ailang-world-operator
description: Use when someone wants to start, watch, grade or finish an AILANG World episode on a repo (ailang-worldd profile init, episode start, episode report, provenance, why, session revoke, world-publish). Prints the commands for the human to run; it never mints a session and never publishes.
---

# Operating a World episode

You print commands; the human runs the attended ones. Two acts are **human only**, because
they read a line typed at a real terminal (`/dev/tty`), which an agent cannot open:

- **minting** a session (`episode start`; also `session new` and `session mint`);
- **publishing** (`world-publish`).

Never run them, never pipe or feed input to them, and never run `session new` or
`session mint` yourself. Print the line, say who runs it, and stop.

The daemon address used below is the rig's, `127.0.0.1:47330`; use the one the operator's
`serve --bind` names. A bare host:port is accepted everywhere.

## Setup (once per repo)

Propose an exec profile from the repo's own files (no daemon, no terminal needed):

```sh
ailang-worldd profile init --repo <dir> --project <name>
```

It writes `~/.ailang-world-exec/profiles/<name>.proposed.json` and prints the `mv` that
activates it. **The operator reviews the file and runs that `mv` themselves** (D-189-5 A).
`serve` started with `--exec-sandbox` loads every active `<project>.json` in the profiles
directory and skips `*.proposed.json`:

```sh
ailang-worldd serve --db <world.db> --workspace-root <dir> --tool-ailang-bin <pin> --exec-sandbox <dir> --bind 127.0.0.1:47330
```

Missing pins or store: `ailang-worldd setup`, then `ailang-worldd doctor`.

## Start (attended: the human runs it)

```sh
ailang-worldd episode start <ep> --repo <dir> --profile <project> --task <file> --addr 127.0.0.1:47330
```

**The human runs this at a terminal and types `y`** at the y/N, which names the worktree
(branch `world/<ep>`), every grant and the profile's digest. It mints the session, writes the
token to `~/.ailang/world/sessions/<ep>.session` (0600) and prints a launch line. The agent
then starts from that launch line, with World's tools and no shell. You never run
`episode start`, and the token never appears in chat, the repo or this plugin.

## Watch

```sh
ailang-worldd episode report <ep> --follow --addr 127.0.0.1:47330
ailang-worldd --addr 127.0.0.1:47330 provenance --episode <ep>
ailang-worldd --addr 127.0.0.1:47330 why <entry|mcp:id> --json
ailang-worldd --addr 127.0.0.1:47330 log tail --follow
```

`episode report` summarises calls by tool, effect decisions, budget left and refusals;
`provenance` lists the episode's entries; `why` explains one entry or one MCP call id.

## Grade

In the ailang-world repo, grading is this repo's own harness: `scripts/trial/trial.py run`
(see its `--help`) runs the graders in the worktree and records one outcome per grader. In any
other repo, the operator runs their own checks (tests, lint, review) on the episode worktree.

## Finish

Once the worktree's work is committed by the operator, revoke the session:

```sh
ailang-worldd session revoke --db <world.db> --episode <ep>
```

Delete the launch's `~/.ailang/world/sessions/<ep>.mcp.json` after the run.

## Publish (human only, irreversible)

`world-publish` writes registries. `publish --live` is an **irreversible public write**. Both
it and its typed confirmation phrase are typed by the human only. Print the command and stop:

```sh
world-publish transitions --store <world.db> --manifest <manifest.json> --ailang-bin <pin>
world-publish publish --store <world.db> --registry-origin <origin> --publisher <pin> --credential-file <file> --approval-ref <ref> --now <n> --expires <e> --live
```
