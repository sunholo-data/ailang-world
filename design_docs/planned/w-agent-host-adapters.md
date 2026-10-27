# w-agent-host-adapters (queue row 223): Claude Code, Codex, pi and herdr into World with no flags — r1

**Status:** r0, written attended 2026-10-27 after Mark's *"let's start some work strands to do
this, yes. A herdr plugin is a good idea"* (roadmap
[`w-all-agent-work-through-world`](w-all-agent-work-through-world.md) §5, D-WORLD-94). Each fact
about a host below was checked on this Mac on 2026-10-27 or taken from the host's official
docs, which are listed in Sources. It is not quorum-reviewed: the pick-time designer
re-measures and runs the quorum as usual. **D-WORLD-97** (RESOLVED 2026-10-27, defaults: stages A-C per launch; machine-wide only on the scorecard, row 219 and rows 196/222) is the one ruling it needs.

**r1 (attended 2026-10-27):** World becomes a dependency of every agent Mark runs through
it, so it must be *"either transparent or clear for errors"* (Mark). §2.5 adds `preflight`, the
fallback and live state. §3 splits stage A into 223a and 223b.

## §0 What "in World" means for a host

An agent is **confined** when its only way to touch files, the shell and the network is World's
MCP tools, under a session Mark authorised. It is **observed** when the host keeps its own tools
and World sees only what passes through it. Every session records which host started it and
the level that host's configuration gives it. The app shows that level, and observed is never
drawn as confined.

## §1 What each host can enforce (measured 2026-10-27)

| host | how World's tools get in | how its own tools are switched off | durable machine-wide lever | level |
|---|---|---|---|---|
| **Claude Code** 2.x (CLI, VS Code and JetBrains extensions, Desktop's local Code tab) | an `http` MCP entry with `headersHelper` or `${VAR}` for the bearer token; `managed-mcp.json` gives World exclusive control | `permissions.deny` by tool name (`Bash`, `Edit`, `Write`, `Read`, `Glob`, `Grep`, `NotebookEdit`, `WebFetch`, `WebSearch`; the full list re-checked at pick time) removes them from the model's context; or `--tools ""`, `--strict-mcp-config` per launch | `/Library/Application Support/ClaudeCode/managed-settings.json` and `managed-mcp.json`: they beat user, project and CLI settings, and a deny at any level can't be re-allowed; plus `allowManagedPermissionRulesOnly`, `allowManagedMcpServersOnly`, `disableBypassPermissionsMode` | **confined** |
| **Codex CLI** 0.164.0 | `[mcp_servers.world]` with `url`, `bearer_token_env_var`, `required = true` and `default_tools_approval_mode = "approve"`: tested, a profile can set all of them | `sandbox_mode = "read-only"`, `[features] shell_tool = false, unified_exec = false`, `web_search = "disabled"`: profile-able; tested, `exec_command` leaves the prompt | `/etc/codex/requirements.toml` (enforced, most restrictive wins; `allowed_sandbox_modes`, `feature_requirements`, an MCP requirement), or the `com.openai.codex` MDM domain | **confined** |
| **pi** 0.85.1 | `pi-mcp-adapter` (`--mcp-config`; `auth: "bearer"`, `bearerTokenEnv`) | `--no-builtin-tools` (read, bash, edit, write); `--no-extensions` then `-e` only the herdr state extension | none found | **confined** per launch |
| **herdr** 0.9.3 | not a host: it runs the hosts above in panes | — | — | inherits its pane's agent |
| **Claude Desktop** (chat) | `claude_desktop_config.json` MCP servers | its own code-execution, file and connector tools can't be switched off (no documented control) | Desktop device policy (`com.anthropic.claudefordesktop`), which doesn't cover those tools | **observed** |
| **IDE agents** (Cursor, Copilot) | their MCP configs, where supported | not controllable | — | **observed** |

**Catches found by measuring:**
- **A Codex profile is layered, not isolated.** `codex -p world` layers `~/.codex/world.config.toml`
  on the user's `config.toml`, so the user's own MCP servers (including `computer-use`), plugins and
  hooks still load. The legacy `[profiles.<name>]` table is now a hard error. The adapter therefore
  uses its **own `CODEX_HOME`** (`~/.ailang-world/codex/`: a `config.toml` holding only World's
  settings, with `auth.json` linked from the real home), not a profile.
- **Claude Code's SessionStart hook can't prompt or block.** Approving a session can't happen
  there. It happens in the mint (a typed `y`, a herdr popup, row 201's key or row 224's
  allowance). The hook only attaches a session that is already approved.
- **Machine-wide policy is machine-wide.** Claude Code's managed files and Codex's
  `requirements.toml` would also confine the mission loop's agents (which run on this Mac) and
  attended sessions like this one, before either goes through World. They also allow only one
  MCP set for every session. That is why §3 starts per launch and goes machine-wide last.
- **An admin can remove the managed files.** It is policy, not a boundary. World OS (rows
  209–211) is the hard boundary.

## §2 The pieces

### §2.1 A session for the launch

Every adapter needs `WORLD_SESSION` for the pane or process it starts. The source moves as the
roadmap lands, and the adapters don't change:
1. **Now:** `episode start` (or `session new`), typed, writing `private/<ep>.session`.
2. **Row 217:** `episode start --preset <p>`, still typed, and short.
3. **Row 201:** a key tap.
4. **Row 224:** drawn from a standing allowance, with no prompt.

### §2.2 Launch configs (World-owned files under `~/.ailang-world/hosts/`, written by `ailang-worldd hosts init`)

- **Claude Code:** `claude-world.settings.json` (the deny list,
  `disableBypassPermissionsMode`, the herdr state hook kept) and `claude-world.mcp.json` (World's
  server, `headersHelper` printing `Authorization: Bearer $WORLD_SESSION`). The launcher is
  `claude --settings … --strict-mcp-config --mcp-config …`, run from an empty directory, as
  `episode start` already instructs.
- **Codex:** `CODEX_HOME=~/.ailang-world/codex codex …`. Its `config.toml` holds w188l's flag set
  as config. The CLI-only flags (`--skip-git-repo-check`, `--ephemeral`) stay in the launcher.
- **pi:** `pi --no-builtin-tools --no-extensions -e <herdr-state.ts> --mcp-config <world-mcp.json>`.
- **Claude Desktop:** an optional World server entry in `claude_desktop_config.json`, recorded as
  observed.
- **The level is recorded.** The launcher sends the host, its version and the digest of the
  config it used with the session's first call. A digest that doesn't match `hosts init`'s files
  records **observed**, never confined. This is the same idea as the exec profile digest bound at
  mint (D-WORLD-85).

### §2.3 The herdr plugin (`integrations/herdr/ailang-world/`, bash)

`herdr-plugin.toml`, using only what herdr 0.9.3 provides:
- **`[[actions]]` "New World agent"** (contexts: workspace, pane):
  1. picks the preset for the workspace's repo (row 217; until then, the rig's episode defaults);
  2. opens a **popup pane** (`[[panes]] placement = "popup"`) that runs the mint. Mark's `y` is
     typed there: the mint stays his act, and it no longer needs an outside terminal;
  3. splits a pane with `--env WORLD_SESSION=… --cwd <empty dir>` (herdr sets env only at pane
     creation) and starts the agent there with §2.2's launcher for Claude, Codex or pi.
- **Labels:**
  - `herdr pane report-metadata <pane> --source ailang-world --token episode=<ep> --token level=confined --token calls=<n>` (tokens render as `$episode`/`$level` in `[ui.sidebar.agents]` rows);
  - `herdr workspace report-metadata` for the count of running episodes.
- **Events:**
  - `pane.exited` / `pane.closed` → offers `episode finish` (row 218), or records that the
    episode is still open;
  - `pane.agent_status_changed` → nothing is written to World. Agent states come from World's
    records (row 207), not the terminal, though herdr's own states keep working.
- **Blocked:** when World has an open decision or approval for the pane's episode, the plugin
  claims state authority (`herdr pane report-agent … --state blocked`) and shows a notification.
  It releases it when the decision is answered.

### §2.4 In the World app

The now card's Running list and each episode's front card show the host, the agent, the level
(confined or observed) and the herdr pane when there is one. Observed sessions use a different,
dashed mark and are never counted as confined.

### §2.5 Transparent or clear: preflight, fallback, live state (r1)

- **`ailang-worldd preflight --session <file> --host <claude|codex|pi>`** checks, in order:
  1. the daemon answers, and its build is the expected one;
  2. the session is valid, with its time and budget left;
  3. the host binary and the launch config, by digest;
  4. the exec tools the profile names, resolved now (friction 23).

  On success it prints one line. On failure it prints **what failed, why, and the next step as a
  command**, with its own exit code per class. It never prints "HTTP 0" or a bare exit code.
- **Fallback:** when preflight fails, the plugin shows the message and offers:
  - fix it (the printed command);
  - or open a **plain pane** labelled `world: fallback`. The fallback is written as a friction
    entry with the preflight message, and counted by the scorecard (row 226). Work never stops,
    and a fallback is never hidden.
- **Live state:** each World pane's sidebar row shows `world: ok`, `session: 1 h left`, or
  `world: down`, refreshed by a light `preflight --quick`. An agent never finds out first.

## §3 Rollout (nothing machine-wide until everything is in World)

| stage | what | needs | changes for Mark |
|---|---|---|---|
| **A1 (223a) — built** | §2.2's launchers and `ailang-worldd hosts init`; §2.5's preflight checks are built into the existing read-only `doctor` (`--session`, `--host`, `--hosts-dir`, `--profile`, `--json`, `--quick`), with three-part errors and class exit codes | today's daemon | per-launch settings only; doctor names stopped daemons, expired sessions and moved exec tools; `/v1/health.exec_tools` reports startup paths |
| **A2 (223b) — built** (`integrations/herdr/ailang-world/`) | §2.3's herdr plugin: "New World agent" (typed mint in a popup, preflight, a confined pane, labelled) and "New plain agent" (the paired arm); the fallback and live state of §2.5 | 223a | run real work through World from herdr, with plain arms to compare |
| **B** | the popup uses presets | row 217 | the mint is one short line |
| **C** | no popup | rows 201 / 224 | no prompt under an allowance |
| **D** | machine-wide: Claude Code managed settings and `managed-mcp.json`, Codex `requirements.toml` | the loop through World (row 196) and attended sessions (row 222); **D-WORLD-97** | every Claude Code or Codex launch on the Mac is confined, whatever starts it |

## §4 Acceptance

1. A Claude Code pane started by the plugin has none of the denied built-ins in its tool list,
   only `mcp__world__*`, and a file write goes through World (it appears in provenance).
2. The same holds for Codex (no `exec_command`; the `computer-use` server from the user's home is
   absent) and for pi (no read, bash, edit or write).
3. A launch whose config digest differs from `hosts init`'s is recorded as **observed**, and the
   app draws it dashed.
4. The pane shows its episode and level in herdr's sidebar. Closing it offers `episode finish`.
   An open World decision marks the pane blocked and shows a notification.
5. Claude Desktop with World's server appears as observed. Nothing in the app calls it confined.
7. **(r1)** With the daemon stopped, an expired session, or a moved exec binary, `preflight`
   prints the three-part message for that case and exits with its class code. The plugin then
   offers the fallback, and choosing it writes a friction entry.
6. Stage A changes nothing machine-wide: an ordinary `claude` or `codex` launch outside the
   plugin behaves exactly as before. The test: the mission loop's next iteration is unaffected.

## §5 Open questions

1. **D-WORLD-97, stage D's timing:** go machine-wide only after rows 196 and 222 (the default),
   or earlier with an exception list for the loop's agents?
2. **The empty working directory:** Claude Code reads `CLAUDE.md` and project settings from its
   working directory. Should World's launcher supply the repo's `CLAUDE.md` as context through
   World, since the agent can't `Read` it?
3. **One herdr workspace per repo,** or per episode?

## Sources

- **Claude Code:**
  - managed settings and precedence:
    [managed-settings](https://code.claude.com/docs/en/managed-settings.md),
    [settings](https://code.claude.com/docs/en/settings.md),
    [permissions](https://code.claude.com/docs/en/permissions.md);
  - [managed MCP](https://code.claude.com/docs/en/managed-mcp.md) and
    [MCP](https://code.claude.com/docs/en/mcp.md);
  - [hooks](https://code.claude.com/docs/en/hooks.md);
  - [CLI reference](https://code.claude.com/docs/en/cli-reference.md);
  - [desktop](https://code.claude.com/docs/en/desktop.md).

  All read 2026-10-27.
- **Codex 0.164.0:** `codex --help` (profile v2: `$CODEX_HOME/<name>.config.toml`), a scratch
  `CODEX_HOME` test of the `world` settings (`mcp get`, `debug prompt-input`), and the binary's
  strings for `/etc/codex/requirements.toml`.
- **herdr 0.9.3:**
  - `herdr api schema --json` (protocol 22), `herdr --default-config`, the herdr-ops skill;
  - the installed agent-quota plugin (`src/herdr.rs`, metadata tokens) as a worked example.
- **pi 0.85.1:** `pi --help`; the `pi-mcp-adapter` README.
