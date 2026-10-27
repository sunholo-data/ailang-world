# ailang-world: a herdr plugin that starts World agent panes

Row 223b ([`w-agent-host-adapters`](../../../design_docs/planned/w-agent-host-adapters.md) §2.3,
§2.5, stage A2). For herdr 0.9.3. POSIX `sh` and `jq`.

- **New World agent:** you pick the host and the episode, doctor checks World, you type the
  mint's `y`, doctor checks the session and the host, and a **confined** agent starts in a new
  pane beside yours.
- **New plain agent:** the comparison arm. The host runs plainly in the repo, labelled
  `arm=plain`.
- **The fallback:** when a World check fails, you choose: see the fix again, open a plain pane
  labelled `world: fallback`, or quit. A fallback is never silent and never labelled confined.
- **Live state:** each World pane's sidebar row shows whether World is up and how long its
  session has left.

The plugin uses only `ailang-worldd hosts init`'s launchers, `doctor --quick --json` and
`episode start` (row 223a). It never types, pipes or simulates the mint's `y`.

## Install

1. Write World's launchers once: `ailang-worldd hosts init --addr http://127.0.0.1:47330`.
2. Make the commands executable: `chmod 755 bin/* test/run.sh`.
3. Link the plugin into herdr (the operator does this, after review):
   `herdr plugin link /path/to/ailang-world/integrations/herdr/ailang-world`.
4. Write the config: run the **World: write config** action, or
   `bin/world-config-init --config-dir <the plugin's herdr config directory>`.
   It writes `world.env` once (0600), prints it and asks you to edit it:
   - `WORLDD`, `WORLD_ADDR`, `WORLD_DB`, `WORLD_WS`, `WORLD_SESSIONS`, `WORLD_HOSTS`: World's
     binary, daemon address, store, workspace root, session directory and launchers;
   - `WORLD_PROFILE`, `WORLD_TTL`, `WORLD_GRANTS`: the mint's default exec profile, lifetime
     and ten `--grant` flags (check the budgets against your rig);
   - `WORLD_REPOS`: one line per herdr workspace, `<workspace label>=<repo path>[:<exec profile>]`;
   - `WORLD_MODELS`: `claude=… codex=… pi=…`.

   Every command refuses, saying what, why and what next, while `world.env` is missing.
5. Show the labels in herdr's sidebar. Add the tokens to your own herdr config (the plugin
   never edits it):

   ```toml
   [ui.sidebar.agents]
   rows = ["$episode $world $session"]
   ```

   Check the `rows` shape against `herdr --default-config` before you add it.

## New World agent

The popup reads one answer per line: the host, the episode id, and only when a check fails,
`f`/`p`/`q`. It can also be run by hand:
`bin/world-agent-popup [--workspace <label>] [--pane <id>]`.

1. The repo comes from the workspace's label in `WORLD_REPOS`. With no entry, it says so and
   stops.
2. `doctor --quick --json` checks the daemon.
3. `episode start <ep> --repo … --agent <host> …` runs in the popup with its stdin, stdout and
   stderr. It asks `y/N` on the terminal, and **you** answer. No `--task`: the agent is
   interactive. For `pi`, the mint uses `--agent codex` (episode start prints launch lines for
   claude and codex only); the pane still runs `pi-world`.
4. `doctor --quick --json --session … --host … --hosts-dir …` checks the session and the
   launch config's digests. A host warning counts as a failure: that launch would be observed.
5. A pane is split to the right, in an empty directory under the plugin's state dir (never the
   worktree), and runs `WORLD_SESSION="$(cat '<session file>')" exec '<hosts>/bin/<host>-world'`.
   The token is read in the pane; it is never in an argv, an env flag or a log.

## Labels

| token | values | meaning |
|---|---|---|
| `episode` | `<ep>` | the World episode this pane works in |
| `level` | `confined` | the agent's only tools are World's, under a session you approved |
| `world` | `ok` | doctor's last check passed |
| | `down` | the daemon does not answer (or doctor could not run) |
| | `expired` | the session expired, was revoked, or its file is unreadable |
| | `host-drift` | the launch config differs from `hosts init`'s: a new launch would be observed |
| | `fallback` | World failed and you chose a plain pane: **not confined** |
| `session` | `<h>h<m>m left`, `expired`, `unknown` | the session's time left |
| `arm` | `plain` | a plain pane: the comparison arm, or a fallback |
| `pair_of` | `<ep>` | the World episode a plain arm is compared with |

A fallback pane also has the state label `world: fallback`.

## State (`$HERDR_PLUGIN_STATE_DIR`)

- `panes.jsonl`: World panes (`pane`, `episode`, `host`, `repo`, `started`, `level`, and
  `ended` once closed);
- `arms.jsonl`: plain and fallback panes;
- `fallbacks.jsonl`: each fallback with the failed check, its detail and its fix;
- `empty/<ep>/`: each agent's empty starting directory;
- `world-watch.lock`, `watch/`: the watcher's lock and its last state per pane.

## The watcher and events

`bin/world-watch` (startup) runs doctor for each open World pane every 60 s and updates its
tokens. A session under 60 minutes, or a failed check, shows one notification per pane per
state change, with the fix. A lock keeps one instance.

When a pane closes or exits, it is marked ended. For a World pane, a notification says its
worktree and session remain: finish the episode from the CLI (row 218 will make this one step).

## Tests

`sh test/run.sh`: fake `herdr` and `ailang-worldd`, no network, no real herdr.
