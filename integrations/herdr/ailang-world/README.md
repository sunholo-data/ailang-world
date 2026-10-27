# ailang-world: a herdr plugin that starts World agent panes

Row 223b ([`w-agent-host-adapters`](../../../design_docs/planned/w-agent-host-adapters.md) §2.3,
§2.5, stage A2). For herdr 0.9.3. POSIX `sh` and `jq`.

Website guide: [Run your agents through World](../../../website/docs/guides/agents-through-world.md).

- **New World agent:** you pick the host and the episode, doctor checks World, you type the
  mint's `y`, doctor checks the session and the host, and a **confined** agent starts in a new
  pane beside yours.
- **New plain agent:** the comparison arm. The host runs plainly in the repo, labelled
  `arm=plain`.
- **The fallback:** when a World check fails, you choose: see the fix again, open a plain pane
  labelled `world: fallback`, or quit. A fallback is never silent and never labelled confined.
- **Live state:** each World pane's sidebar row shows whether World is up and how long its
  session has left.
- **Key bindings:** one action shows three bindings and adds them to herdr's config after you
  type `y` (row 223c).
- **The World panel:** a pane listing this workspace's agent panes, their sessions, calls and
  budgets, with keys to start, inspect and revoke (row 223c).

The plugin uses only `ailang-worldd hosts init`'s launchers, `doctor --quick --json` and
`episode start` (row 223a). It never types, pipes or simulates the mint's `y`.

## Install

1. Write World's launchers once: `ailang-worldd hosts init --addr http://127.0.0.1:47330`.
2. The commands are executable in git (mode 755); if a copy lost that, `chmod 755 bin/* test/run.sh`.
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
   edits it only through "World: install key bindings", only between its own markers, and
   only after your `y`):

   ```toml
   [ui.sidebar.agents]
   rows = [["state_icon", "machine", "workspace", "tab"], ["agent"], ["$episode", "$world", "$session"]]
   ```

   `rows` is a list of rows, each a list of tokens (herdr 0.9.3's `herdr --default-config`);
   the first two rows are herdr's defaults, and the third adds the plugin's tokens.

## Running an action without keys

Every action runs from herdr's command palette, or from a shell:
`herdr plugin action invoke <id> --plugin ailang-world`, where `<id>` is `world-agent`,
`plain-agent`, `world-config`, `world-keys` or `world-panel`.

## Key bindings

A herdr manifest cannot declare keys. They live only in herdr's own config:
`$HERDR_CONFIG_FILE` when set, else `~/.config/herdr/config.toml`. The **World: install key
bindings** action (`world-keys`) opens a popup that runs `bin/world-keys`. That is the only
command in the plugin that writes herdr's config, and it writes only after you type `y`.

| key | action |
|---|---|
| `prefix+shift+a` | New World agent (`ailang-world.world-agent`) |
| `prefix+shift+l` | New plain agent (`ailang-world.plain-agent`) |
| `prefix+shift+o` | World panel (`ailang-world.world-panel`) |

- `bin/world-keys --check` prints the config path, the block and any conflicts, and checks
  that the result would be valid TOML. It changes nothing. It exits 0 when `--apply` could
  apply.
- `bin/world-keys --apply` (the default, and the action):
  1. prints the path and the exact block, then asks `Add these key bindings to <path>? [y/N]`.
     Any answer but `y` changes nothing.
  2. on `y`, it writes a one-time backup of the original to
     `$HERDR_PLUGIN_STATE_DIR/herdr-config.original.toml` (it never overwrites an existing
     backup; it is empty when there was no config);
  3. it writes the new file beside the config, checks it with Python's `tomllib`, and only
     then moves it into place, keeping the file's mode;
  4. it runs `herdr server reload-config`.

  The block sits between `# >>> ailang-world key bindings (bin/world-keys) >>>` and
  `# <<< ailang-world key bindings <<<`. A second apply replaces the block and never adds
  another. Nothing outside the markers is touched. A missing config is created holding just
  the block. A symbolic link refuses: set `HERDR_CONFIG_FILE` to its target.
- `bin/world-keys --uninstall` shows the block, asks `y` the same way, removes exactly the
  marker-to-marker lines, checks, moves into place and reloads. By hand outside herdr, set
  `HERDR_BIN_PATH` to herdr's binary first.
- **Conflicts:** a key bound anywhere outside the block, or one of herdr's defaults
  (`prefix+shift+r`, `n`, `g`, `w`, `d`, `t`, `x`, `p`), refuses before asking. To choose
  other keys, set the optional `WORLD_KEYS` in `world.env`, e.g.
  `WORLD_KEYS="agent=prefix+shift+a plain=prefix+shift+l panel=prefix+shift+o"`.

It needs `python3` 3.11 or later (`tomllib`), and refuses rather than write a config it could
not check.

## World panel

The **World panel** action (`world-panel`) opens `bin/world-panel` split beside your pane.
It redraws every 5 s and on each key:

- **the header:** World `ok`/`down` (`doctor --quick --json`), the daemon address, the
  workspace and its repo (from `WORLD_REPOS`), and whether the watcher is running;
- **one row per open pane** of this repo, from `panes.jsonl` and `arms.jsonl`: pane, episode,
  host, level (`confined`, `plain`, `fallback`), World state (`ok`, `down`, `expired`,
  `host-drift`), the session's time left, calls so far and the lowest budget left (from
  `episode report --json`);
- **the selected row:** the last lines of `provenance --episode <ep>`. Today that is the
  `World-Provenance` trailer with the episode's entry range, so the panel names the last
  entry index. Also the surface link: `$WORLD_ADDR/surface` (a surface place is a brane,
  which is a charter row's lane, not an episode);
- **the keys:** `j`/`k` or the arrows select; `n` New World agent; `l` New plain agent; `o` opens
  the surface (`open` or `xdg-open`) and prints the URL; `w` asks for an entry index and pages
  `ailang-worldd why <n>` with `less -R`; `r` asks `y`, then revokes the selected episode's
  session (`session revoke --db … --episode <ep>`); `q` quits.

`bin/world-panel --once` prints one frame. `--keys` reads keys one per line from stdin.
`--workspace <label>` overrides the workspace. Sessions are read only by path, by doctor.
The panel never prints a token.

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
- `world-watch.lock`, `watch/`: the watcher's lock and its last state per pane;
- `world-watch.err`: the stderr of a watcher started by an action;
- `herdr-config.original.toml`: `bin/world-keys`'s one-time backup of herdr's config.

## The watcher and events

`bin/world-watch` (startup) runs doctor for each open World pane every 60 s and updates its
tokens. A session under 60 minutes, or a failed check, shows one notification per pane per
state change, with the fix. A lock keeps one instance. herdr runs startup commands only when
it starts, so after `herdr plugin link` the New World agent, New plain agent and World panel
actions start the watcher, detached, when its lock shows it isn't running.

When a pane closes or exits, it is marked ended. For a World pane, a notification says its
worktree and session remain: finish the episode from the CLI (row 218 will make this one step).

## Tests

`sh test/run.sh`: fake `herdr` and `ailang-worldd`, no network, no real herdr. Every case
sets `HERDR_CONFIG_FILE` to a temporary file, so the tests never read or write your herdr
config.
