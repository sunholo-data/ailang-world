---
title: Run your agents through World
sidebar_position: 2
description: Launch Claude Code, Codex and pi confined to World's tools, by hand or from herdr; what the agent is told; and what to do when a check fails.
---

# Run your agents through World

The goal (D-WORLD-94) is that all agent work runs through World: headless or interactive, on
any repo. This page covers what exists as of rows 223a–223d: launchers for Claude Code, Codex
and pi, a `doctor` preflight, a herdr plugin, and the session context each agent starts with.
The roadmap is `design_docs/planned/w-all-agent-work-through-world.md`.

## Why

An agent that runs through World leaves an entry for every effect it causes, and each entry
can be walked back with `ailang-worldd why`. How much World sees depends on the level:

- **Confined.** World's tools are the agent's only file and shell access, so every effect is
  a World entry. The launchers on this page start agents confined.
- **Observed.** The agent keeps its own built-in tools. World sees its MCP calls, not its
  built-ins. An observed agent is never labelled confined.
- **Off-machine actions** are never controlled. What a cloud service does on its own servers
  happens outside World.

Confinement is **per launch**: a launcher confines the agent it starts, and nothing else on the
machine. See [Tool confinement](../security/tool-confinement.md) for what each host's launch
turns off.

## Set up once

Write World's launch settings and launchers:

```bash
ailang-worldd hosts init
```

```text
usage: ailang-worldd hosts init [--out-dir <dir>] [--addr <url>] [--replace] [--pi-mcp-adapter <dir>]
```

- The files go to `~/.ailang-world/hosts` (`--out-dir`). Directories are 0700, configs 0600,
  launchers 0700.
- The launchers are `bin/claude-world`, `bin/codex-world` and `bin/pi-world`. Put that `bin`
  on your `PATH`, or call them by path.
- `--addr` is the World URL the launchers use (default `http://127.0.0.1:47330`).
- A file that differs from what `hosts init` would write is refused unless `--replace`.
  Unchanged files are left alone.
- pi needs `pi-mcp-adapter` (pi has no built-in MCP). `--pi-mcp-adapter` defaults to
  `~/.pi/agent/npm/node_modules/pi-mcp-adapter`. Without it, `pi-world` refuses to start.
- When `~/.codex/auth.json` is absent, it prints
  `! Codex will ask you to log in: ~/.codex/auth.json does not exist.`

Then check that World answers:

```bash
ailang-worldd --addr http://127.0.0.1:47330 doctor --quick --db <the store World serves>
```

`--quick` checks only the daemon, the session and the host. `--db` must name the store the
daemon serves, or the daemon check fails. Add `--session <file> --host claude` (or `codex`,
`pi`) to check a session and a launch config too.

The context block and Claude's status line call `ailang-worldd` and `jq`, so both must be on
the agent's `PATH`.

## Launch by hand

1. Start an episode. This is an attended step: you type `y` at the terminal.

   ```bash
   ailang-worldd episode start ep1 --repo ~/src/my-project --agent claude
   ```

   The token goes to `<sessions-dir>/ep1.session` (default `~/.ailang/world/sessions`).
   `episode start` prints launch lines for `claude` and `codex` only. For pi, mint with
   `--agent codex`.
2. Start the agent from an empty directory, not the worktree. The launcher reads the session
   from `WORLD_SESSION`:

   ```bash
   WORLD_SESSION="$(cat ~/.ailang/world/sessions/ep1.session)" claude-world
   WORLD_SESSION="$(cat ~/.ailang/world/sessions/ep1.session)" codex-world
   WORLD_SESSION="$(cat ~/.ailang/world/sessions/ep1.session)" pi-world
   ```

Without `WORLD_SESSION` a launcher refuses (exit 3):

```text
claude-world: no World session (WORLD_SESSION is unset).
  why: World agents need a session you approved.
  next: start one with `ailang-worldd episode start …`, then export WORLD_SESSION="$(cat <session file>)"
```

Arguments after the launcher's name pass to the host.

## In herdr

The herdr plugin is `integrations/herdr/ailang-world/` (herdr 0.9.3 or later). Its README has
the full detail; this is the short form.

1. **Link it.** `herdr plugin link <repo>/integrations/herdr/ailang-world`.
2. **Configure it.** Run **World: write config**. It writes `world.env` once in the plugin's
   herdr config directory and asks you to edit it: World's binary, address, store, workspace
   root, sessions and launchers; the mint's profile, lifetime and grants; one
   `<workspace label>=<repo path>[:<exec profile>]` line per herdr workspace in `WORLD_REPOS`;
   and the models.
3. **Key bindings.** Run **World: install key bindings**. It shows the block and asks `y`,
   keeps a one-time backup of your herdr config, checks the result is valid TOML, then reloads
   herdr. `bin/world-keys --uninstall` removes the block exactly.

   | key | action |
   |---|---|
   | `prefix+shift+a` | New World agent |
   | `prefix+shift+l` | New plain agent |
   | `prefix+shift+o` | World panel |

   Without keys, run any action from herdr's command palette, or from a shell:
   `herdr plugin action invoke <id> --plugin ailang-world`, where `<id>` is `world-agent`,
   `plain-agent`, `world-config`, `world-keys` or `world-panel`.

**New World agent** opens a popup. You pick the host and the episode. `doctor` checks the
daemon, then `episode start` asks its `y/N` and **you** type the answer; the plugin never does.
`doctor` then checks the session and the host. A confined pane starts beside yours, from an
empty directory, with these labels:

| token | value |
|---|---|
| `episode` | the episode id |
| `level` | `confined` |
| `world` | `ok`, `down`, `expired` or `host-drift` |
| `session` | the time left, e.g. `1h42m left` |

**New plain agent** opens the host plainly in the repo, labelled `arm=plain`. It is the
comparison arm.

**The watcher** runs `doctor` for each open World pane every 60 s and keeps the labels live. A
session under 60 minutes, or a failed check, sends one notification per change.

**The World panel** opens in a split pane. It lists this workspace's World, plain and fallback
panes with their state, time left, calls and lowest budget left, and the selected pane's recent
entries. Keys: `n` New World agent, `l` New plain agent, `o` open the surface, `w` explain an
entry (`ailang-worldd why`), `r` revoke the selected session (asks `y`), `q` quit.

**The fallback.** When a check fails, the popup shows what failed and the fix, then offers
`[f]` show the fix again, `[p]` open a plain pane instead, or `[q]` quit. A plain pane opened
this way is labelled `world: fallback`, recorded in the plugin's `fallbacks.jsonl`, and
announced with a notification. It is never labelled confined.

## What the agent knows

Every World launch starts with a context block:

```text
You are working inside AILANG World (confined).
Episode ep1 · worktree <workspace-root>/ep1 · session 1h42m left · daemon http://127.0.0.1:47330
Your only file and shell access is World's MCP tools:
  read/edit/write files -> ailang-read / ailang-edit / ailang-write / ailang-delete (paths relative to the worktree)
  run commands         -> workspace-exec, the exec profile's commands only (list them with the tool)
  AILANG               -> ailang-check / ailang-run / ailang-cli; docs -> builtins-search / examples-search
Grants left: Workspace.Read 48/50, Workspace.Write 50/50, …
Every call is recorded. See it: http://127.0.0.1:47330/surface  ·  explain an entry: ailang-worldd why <entry>
Out of budget or need a new grant? Say so; the operator mints it. Don't try to work around World.
```

When World can't answer, the block is one line instead, such as
`World context unavailable: session invalid or expired; ask the operator to mint a new session.`

How each host gets it:

- **Claude Code:** a World plugin (`--plugin-dir`) whose session-start hook prints the block.
  The status line reads `World · <episode> · <h>h<m>m · ok` (or `expired`, `down`), refreshed
  every 60 s. If your own Claude settings have a status line, its output follows after ` │ `.
- **Codex:** the block as developer instructions, plus a fixed `AGENTS.md` with the tool map.
  When herdr's Codex state hook is installed, the launch runs it too.
- **pi:** the block appended to the system prompt, and a footer `World · <episode> · <h>h<m>m`.
  When herdr's pi state extension is installed, the launch loads it too.

The agent, or you, can ask World directly:

```bash
ailang-worldd --addr http://127.0.0.1:47330 whoami
```

```text
World · ep1
worktree: <workspace-root>/ep1
session: 1h42m left
daemon: http://127.0.0.1:47330
surface: http://127.0.0.1:47330/surface
  Workspace.Read worktree: 48/50 left
```

`whoami` reads the session from `WORLD_SESSION` or `--session <file>`, never from an argument.
`--json` prints the body of [`GET /v1/session/self`](../reference/http-api.md). It exits 10 when
the daemon doesn't answer and 11 when the session is invalid or expired.

## When something fails

Every failure says what failed, why, and the next command. `doctor` exits with the class of the
first failing check, in this order:

| exit | class | typical cause | next |
|---|---|---|---|
| 10 | daemon | nothing listens at `--addr`, or it serves another store | `ailang-worldd serve --db <store> …`, or pass the store's `--db` |
| 11 | session | the session file is missing, not mode 0600, expired or revoked | `ailang-worldd episode start …` (or `session new`) for a new session |
| 12 | host | the host is not on `PATH`, or pi's adapter is missing | install the host; for pi, `pi install npm:pi-mcp-adapter`, then `ailang-worldd hosts init --replace` |
| 13 | exec | an exec profile's tool moved or is missing | install it in the profile's path, then restart World |

A launch config that differs from what `hosts init` wrote is a host **warning**: doctor still
exits 0, but that launch would be observed, not confined. The fix is
`ailang-worldd hosts init --out-dir <dir> --replace`. The herdr plugin treats this warning as a
failure and labels the pane `host-drift`.

In herdr, a failed check offers the fallback (above). See
[Troubleshooting](troubleshooting.md) for each message.

## What's next

- **Row 219:** World as a supervised user service that survives a reboot and an upgrade.
- **Row 217:** a one-command start from a per-repo preset.
- **Row 218:** `episode finish`: checks, commit, push and PR in one step.
- **Row 224:** standing allowances: at most 10 sessions a day, 12 h each, for 30 days,
  revocable, and never for the mission loop.

Confinement stays per launch. It goes machine-wide only when the value scorecard (row 226)
holds over 10 episodes with zero unclear errors, the daemon is a supervised service (row 219),
and the mission loop and attended sessions run in World (rows 196 and 222). See the
[roadmap](../roadmap/index.md).
