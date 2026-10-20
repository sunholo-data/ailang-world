# Use World on your repo

Run a coding agent on your own repository through AILANG World. The agent works only through
World's tools, so every act it takes is in World's log and can be explained afterwards.

**Before you start:** you have a git repository, a coding agent (Claude Code or codex), and an
`ailang-worldd` built from this repo ([QUICKSTART §1](QUICKSTART.md#1-build-and-start)).

Each step says who does it: **you** (the human, at a terminal), **the agent**, or **either** (you,
a script or an agent). Placeholders: `<dir>` is your repo, `<project>` its profile name, `<ep>`
the episode (`^[a-z0-9][a-z0-9-]{0,63}$`), `<world.db>` the world store (`setup`'s default is
`~/.ailang/world/world.db`), `<workspace-root>` where worktrees live (default
`~/.ailang/world-ws`), and `<sessions-dir>` where tokens live (default `~/.ailang/world/sessions`).

## Once per machine

### 1. Install the plugin, the pins and the store (you)

```sh
claude plugin marketplace add sunholo-data/ailang-world
claude plugin install ailang-world@ailang-world-marketplace
```

The repo is public and its default branch is `dev`. The plugin gives your agents the `ailang-world`
skills and the print-only commands `/world-start`, `/world-watch` and `/world-why`.

Then install the two pinned AILANG binaries and the store directories, and check the machine:

```sh
ailang-worldd setup
ailang-worldd doctor
```

Finally publish World's coding tools into the store. This is a human act (you type the phrase it
shows), and the daemon must not be running yet. Run it from a checkout of this repo, with
`world-publish` built from `./cmd/world-publish`:

```sh
world-publish transitions --store <world.db> --manifest packages/se-tools/transitions.json --ailang-bin <interpreter-pin>
```

[QUICKSTART §0](QUICKSTART.md#0-attended-steps-the-only-commands-a-person-must-run) explains the
publish prompt and the snags (`AILANG_REGISTRY_API_KEY`, IDE terminals).

## Once per repo

### 2. Propose an exec profile (either), then activate it (you)

```sh
ailang-worldd profile init --repo <dir> --project <project>
```

It reads the repo's own files and needs no daemon and no terminal. The proposal holds the
profile's `root`, the toolchain `path`, `env` and `read_roots` resolved on this machine, a
`probe`, `caches`, and the commands the agent may run through `workspace-exec`: `go-test`,
`go-vet` and `gofmt-l` from `go.mod`; `pytest-file` and `ruff-file` from `pyproject.toml` or
`uv.lock`; `vitest-file` or `jest-file` from `package.json`. It prints why each command is there,
then writes `<project>.proposed.json` in `~/.ailang-world-exec/profiles/` and prints its digest
and an `mv`.

**You** review the file and run that `mv`. It renames the proposal to `<project>.json`, which
activates it. `serve` loads it and probes it at startup. A profile that fails the probe is left
out, named on `serve`'s stderr, and every other profile still serves. Adding a profile while
`serve` runs needs one restart of `serve`.

## Once per machine session

### 3. Start the daemon (either)

```sh
ailang-worldd serve --db <world.db> --bind 127.0.0.1:47330 --ailang-bin <interpreter-pin> --workspace-root <workspace-root> --tool-ailang-bin <tool-pin> --exec-sandbox <srt-node_modules>
```

- `--ailang-bin` and `--tool-ailang-bin` are the two binaries `setup` installed.
- `--exec-sandbox` names the `node_modules` holding `@anthropic-ai/sandbox-runtime`. With it,
  and with no `--exec-profile`, `serve` loads every `<project>.json` in
  `~/.ailang-world-exec/profiles/` (`--exec-profiles-dir` names another directory).
- If you run step 5 with no daemon serving, its refusal prints this line with your paths filled
  in.

Leave it running. It serves every episode. You never restart it for a new episode, and minting
works while it serves (row 189 M1–M3): an episode bound to a profile runs under that profile from
its first call.

## Per episode

### 4. Write the task (you, or an agent in your repo)

Put what the agent should do in a file, for example `<task-file>`: the change, the files to read
first, the checks you will run, and what the summary must say. The agent cannot ask you
questions mid-run, so put everything it needs in the file.

### 5. Start the episode (you)

```sh
ailang-worldd episode start <ep> --repo <dir> --profile <project> --task <task-file> --addr 127.0.0.1:47330
```

In the plugin, `/world-start <ep>` prints this line for your repo. If your store or workspace
root are not the defaults, add `--db <world.db>` and `--workspace-root <workspace-root>`.

It validates everything and checks the daemon first, then opens your terminal and asks a y/N.
The y/N names:
- the worktree, `<workspace-root>/<ep>`, created from `<dir>` on the new branch `world/<ep>`;
- every grant;
- the profile's project and the sha256 of its exact bytes.

**Type `y`.**

The default grants are the `se-tools` preset without `Workspace.Exec`, plus `Workspace.Exec`
because you passed `--profile`. That is `Workspace.Read`, `Workspace.Write`, `Ailang.Check`,
`Ailang.Run`, `Ailang.RunEnv`, `Ailang.RunNet`, `Ailang.Discover`, `Ailang.CLI`,
`Workspace.Exec` and `Workspace.Delete`. Each has scope `worktree` and 50 calls (`--budget`), and
the session lasts 3600 s (`--ttl`). `--preset` or `--grant` replace the default.

It refuses, before any worktree, credential or file exists, when:
- there is no controlling terminal (an agent, cron or CI ran it);
- the daemon on `--addr` is another build, serves another tool binary than this build pins, reads
  another credentials store, or has not loaded the profile's current bytes. The refusal names
  the fix;
- no daemon is serving. The refusal prints the `serve` line;
- `<sessions-dir>/<ep>.session` or `<sessions-dir>/<ep>.mcp.json` already exists.

On success it writes the token once to `<sessions-dir>/<ep>.session` (0600) and, for claude, a
private 0600 `<sessions-dir>/<ep>.mcp.json`. The token is never printed. It prints the worktree
and branch, the launch line, the task path, the file to delete after the run, and a
`provenance` line to watch with. It never runs the agent, commits or pushes.

### 6. Launch the agent (you)

Open an **empty directory, never the worktree**, and run the launch line `episode start` printed.
For claude it has this shape:

```sh
claude --tools '' --strict-mcp-config --mcp-config <sessions-dir>/<ep>.mcp.json --allowedTools <the World tools the grants cover>
```

Give it the task file. The agent reaches the worktree only through World's tools, so its record
is complete. With `--agent codex` in step 5, the launch line is a `codex exec` command that reads
the token from the session file through `WORLD_SESSION`, and no `.mcp.json` is written.

### 7. Watch (either)

```sh
ailang-worldd episode report <ep> --follow --addr 127.0.0.1:47330
ailang-worldd --addr 127.0.0.1:47330 provenance --episode <ep>
ailang-worldd --addr 127.0.0.1:47330 why <entry> --json
```

`episode report` summarises calls by tool, effect decisions, the budget left on each grant, and
refusals. `provenance` lists the episode's log entries. `why` explains one entry, or one MCP
call id (`mcp:<id>`).

### 8. Review (you)

Run your own tests, lint and review in the worktree, `<workspace-root>/<ep>`. Then read the
refusals the agent hit:

```sh
ailang-worldd episode report <ep> --json --addr 127.0.0.1:47330
```

Plan refusals and handler `ok:false` results are listed there. Check that the agent's summary
reports each of them.

### 9. Finish (you)

Commit and push from the worktree, which is on branch `world/<ep>`:

```sh
git -C <workspace-root>/<ep> add -A
git -C <workspace-root>/<ep> commit
git -C <workspace-root>/<ep> push -u origin world/<ep>
```

Delete the launch's MCP config, which holds the token, then revoke the session:

```sh
rm <sessions-dir>/<ep>.mcp.json
ailang-worldd session revoke --db <world.db> --episode <ep>
```

## What only you can do

Minting (`episode start`, `session new`) and publishing (`world-publish`) read a line you type at
a real terminal (`/dev/tty`). An agent cannot run them: it has no terminal, and the command
refuses. The plugin's skills and commands print these lines for you and never run them. See
[QUICKSTART §0](QUICKSTART.md#0-attended-steps-the-only-commands-a-person-must-run).

## Steps per episode

| | Human acts per episode |
|---|---|
| Before row 189 ([onboarding doc §1](../design_docs/planned/w-world-onboarding.md#1-why-measured)) | 11, including stopping, rebuilding and restarting the daemon |
| This page | 6: write the task (4), `episode start` and `y` (5), launch the agent (6), then commit and push, delete the `.mcp.json`, and revoke (9) |

Your review (step 8) is extra, as reading the graders was before. AC7 of the onboarding doc
measures this count on a real repo, in an attended run that follows this page.

## Troubleshooting

- **`refusing: no controlling terminal`**, or **`STOP fence=tty reason=no-controlling-terminal`**
  from `world-publish`: an agent, a script or a pipe ran an attended step. Run it yourself in a
  terminal window. An IDE terminal pane works.
- **`refusing: the daemon on … is build …`**: `serve` runs another build than your
  `ailang-worldd`. Rebuild `ailang-worldd` from your checkout, stop `serve`, start it again with
  the new binary, and rerun `episode start`. The same fix applies when the daemon serves another
  tool binary than this build pins: rerun `setup` first if the pin is missing.
- **`refusing: no daemon is serving on …`**: run the `serve` line it prints (step 3), then rerun.
- **The daemon has not loaded the profile, or the profile changed:** put `<project>.json` in the
  profiles directory, or keep its current bytes, and restart `serve` once.
- **`… .session already exists`** (or `.mcp.json`): an earlier run of that episode left it. If
  that episode is done, delete the file and revoke it (step 9), or pick another episode name.
- **`no exec profile configured`** from `workspace-exec`: `serve` runs with no exec profile
  (started without `--exec-sandbox`), or it loaded several and the episode is bound to none.
  Restart `serve` as in step 3, and start episodes with `--profile <project>`. Without
  `--profile`, `episode start` grants no `Workspace.Exec` at all.
