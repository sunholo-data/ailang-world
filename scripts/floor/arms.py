#!/usr/bin/env python3
"""Row 93 floor harness — the arms: tool and grant set, and the exact argv (design §4.3, §4.5).

Two things live here. First, the World arm's tool and grant set (2026-10-06): the World arm is
EXACTLY the 8 AILANG se-tools. Second (M2, AC2.1), the exact argv per (agent, arm), checked against
the goldens in ``testdata/argv_<agent>_<arm>.json`` by ``test_arms.py``.

World serves 9 se-tools since row 140 (``workspace-exec``, effect ``Workspace.Exec``). The floor
benchmark is AILANG-only, so the World arm:
  * lists only ``prompt.WORLD_TOOL_NAMES`` in its prompt (prompt.py);
  * allows only their ``mcp__world__<id>`` names (Claude ``--allowedTools``; §4.4 pre-flight
    expects exactly these from ``tools/list``);
  * is minted with EXPLICIT ``--grant`` flags for exactly the effects those 8 tools declare, never
    ``--preset se-tools`` (which grants ``Workspace.Exec``). A session sees only the tools it holds
    a grant for (QUICKSTART §9), so leaving the grant out is what keeps ``workspace-exec`` off the
    served list, not just off the prompt.

The 8 tools declare 8 effects since row 135: ``ailang-run`` declares ``Ailang.RunEnv`` and
``Ailang.RunNet`` beside ``Ailang.Run`` (D-NF-3 = B widened it for this floor), and the two
searches share ``Ailang.Discover``, ``write``/``edit`` share ``Workspace.Write``.
"""
from __future__ import annotations

import hashlib
import json

import prompt

MCP_SERVER = 'world'
GRANT_SCOPE = 'worktree'
GRANT_BUDGET = 5000  # §4.5: each grant `worktree:5000`, --ttl 43200
SESSION_TTL_S = 43200

# Effect name per World-arm tool, as packages/se-tools/transitions.json declares it
# (test_arms.py checks this table against the manifest).
WORLD_TOOL_EFFECTS = {
    'ailang-read': ('Workspace.Read',),
    'ailang-write': ('Workspace.Write',),
    'ailang-edit': ('Workspace.Write',),
    'ailang-check': ('Ailang.Check',),
    'ailang-run': ('Ailang.Run', 'Ailang.RunEnv', 'Ailang.RunNet'),
    'builtins-search': ('Ailang.Discover',),
    'examples-search': ('Ailang.Discover',),
    'ailang-cli': ('Ailang.CLI',),
}


def world_allowed_tools() -> tuple[str, ...]:
    """The Claude ``--allowedTools`` names (and the pre-flight's expected ``tools/list``)."""
    return tuple(f'mcp__{MCP_SERVER}__{name}' for name in prompt.WORLD_TOOL_NAMES)


def world_grant_effects() -> tuple[str, ...]:
    """The effect names the World-arm session is granted, in first-declared order, deduplicated."""
    out: list[str] = []
    for name in prompt.WORLD_TOOL_NAMES:
        for eff in WORLD_TOOL_EFFECTS[name]:
            if eff not in out:
                out.append(eff)
    return tuple(out)


def world_grant_args(budget: int = GRANT_BUDGET) -> list[str]:
    """``session new`` grant flags: explicit ``--grant EFFECT=worktree:N``, never ``--preset``."""
    args: list[str] = []
    for eff in world_grant_effects():
        args += ['--grant', f'{eff}={GRANT_SCOPE}:{int(budget)}']
    return args


# ---------------------------------------------------------------- the argv (§4.3, AC2.1)

AGENTS = ('claude', 'codex')
ARMS = ('shell', 'world')
DEFAULT_ADDR = '127.0.0.1:7644'
DEADLINE_S = 600  # §4.3: every task-run, both arms
TOKEN_ENV = 'WORLD_SESSION'  # codex World arm: the bearer token reaches codex ONLY through this env var

# Claude shell arm: the native coding tools, NO web tools ([REFINED] §4.3), Bash sandboxed by settings.
CLAUDE_SHELL_TOOLS = ('Bash', 'Read', 'Write', 'Edit', 'Glob', 'Grep')
CLAUDE_SANDBOX_SETTINGS = {'sandbox': {'enabled': True, 'autoAllowBashIfSandboxed': True,
                                       'allowUnsandboxedCommands': False}}
CLAUDE_EMPTY_MCP = {'mcpServers': {}}
CLAUDE_SHELL_PERMISSION_MODE = 'acceptEdits'  # §4.3 as amended by M2 (V-M2-3)
# [M4 knob (b), harness-side, BOTH Claude arms identically; V-M4-1] §4.3's "user MCP servers,
# plugins and skills are suppressed identically in both arms" — M2 measured that the user's
# plugins were still loaded and 2 SessionStart hooks still ran in both arms (V-M2-3).
# `--setting-sources ""` loads no user/project/local settings file (so no user plugins, hooks or
# enabled-plugin list); the built-in plugins that remain are switched off and auto-memory is off
# through the same `--settings` object the shell arm already uses for its sandbox. The flag set is
# identical in both arms; only the shell arm's settings object additionally carries the sandbox.
CLAUDE_BUILTIN_PLUGINS = ('cc-plugin-agents-md@builtin', 'cc-plugin-telemetry@builtin',
                          'cc-plugin-plugin-authoring@builtin')
CLAUDE_ISOLATION_SETTINGS = {'enabledPlugins': {p: False for p in CLAUDE_BUILTIN_PLUGINS},
                             'autoMemoryEnabled': False}
# codex, both arms: every non-coding feature off (V26, V28); web search off by config key (V-M2-1).
CODEX_DISABLED_FEATURES = ('plugins', 'apps', 'browser_use', 'computer_use', 'image_generation',
                           'multi_agent')
CODEX_WORLD_DISABLED_FEATURES = ('shell_tool', 'unified_exec')  # MUT-CODEX-SHELL drops shell_tool


class ArmError(Exception):
    pass


def _compact(obj) -> str:
    return json.dumps(obj, separators=(',', ':'))


def mcp_url(addr: str = DEFAULT_ADDR) -> str:
    return f'http://{addr}/mcp/'


def claude_mcp_config(token: str, addr: str = DEFAULT_ADDR) -> dict:
    """The Claude World arm's ``--mcp-config`` file body. It HOLDS THE RAW TOKEN: run_task writes it
    mode 0600 in the harness-private dir and deletes it after the run; it is never evidence."""
    return {'mcpServers': {MCP_SERVER: {'type': 'http', 'url': mcp_url(addr),
                                        'headers': {'Authorization': f'Bearer {token}'}}}}


def codex_world_overrides(addr: str = DEFAULT_ADDR) -> list[str]:
    """The ``-c`` overrides that give codex exactly the ``world`` server (token by env var name)."""
    # The design's shell text `-c key="v"` reaches codex unquoted (`key=v`): a value that is not
    # TOML is taken as a literal string (codex -c help), measured identical by `codex mcp list`.
    return ['-c', f'mcp_servers.{MCP_SERVER}.url={mcp_url(addr)}',
            '-c', f'mcp_servers.{MCP_SERVER}.bearer_token_env_var={TOKEN_ENV}',
            '-c', f'mcp_servers.{MCP_SERVER}.required=true']


def build_argv(agent: str, arm: str, *, model: str, prompt_text: str, cwd: str,
               mcp_config_path: str | None = None, addr: str = DEFAULT_ADDR,
               claude_isolation: bool = False, max_budget_usd: float | None = None) -> list[str]:
    """The exact argv of §4.3 for one (agent, arm). ``cwd`` is the worktree (shell) or the empty
    ``void/<ep>`` (World); Claude takes it as the process cwd, codex also as ``-C``.

    The two M4 additions are Claude-only, keyword-only and applied IDENTICALLY to both arms, so the
    defaults reproduce the §4.3 goldens exactly: ``claude_isolation`` (knob (b), V-M4-1) appends
    ``--setting-sources ""`` and merges ``CLAUDE_ISOLATION_SETTINGS`` into the arm's ``--settings``;
    ``max_budget_usd`` (§4.9) appends ``--max-budget-usd``. codex ignores both (it has no budget
    flag; ``--ignore-user-config`` already drops its user config)."""
    if agent not in AGENTS or arm not in ARMS:
        raise ArmError(f'unknown agent/arm {agent!r}/{arm!r}')
    if not model or not prompt_text or not cwd:
        raise ArmError('model, prompt and cwd are required')
    if agent == 'claude':
        argv = ['claude', '-p', prompt_text, '--model', model, '--output-format', 'stream-json',
                '--verbose', '--no-session-persistence', '--disable-slash-commands']
        tail = []
        if claude_isolation:
            tail += ['--setting-sources', '']
            if arm == 'world':
                tail += ['--settings', _compact(CLAUDE_ISOLATION_SETTINGS)]
        if max_budget_usd is not None:
            if not max_budget_usd > 0:
                raise ArmError(f'max_budget_usd must be positive, got {max_budget_usd!r}')
            tail += ['--max-budget-usd', f'{float(max_budget_usd):.2f}']
        if arm == 'shell':
            if mcp_config_path is not None:
                raise ArmError('the Claude shell arm takes no MCP config file')
            return argv + ['--tools', ','.join(CLAUDE_SHELL_TOOLS), '--strict-mcp-config',
                           '--mcp-config', _compact(CLAUDE_EMPTY_MCP),
                           # [M2, V-M2-3] acceptEdits, not the drafted bypassPermissions: the
                           # sandbox confines Bash only; under bypassPermissions the Write tool
                           # wrote outside the worktree (AC2.5 canary r1). acceptEdits refuses an
                           # edit outside the cwd in -p mode and keeps sandboxed Bash auto-allowed.
                           '--permission-mode', CLAUDE_SHELL_PERMISSION_MODE,
                           '--settings', _compact(dict(CLAUDE_SANDBOX_SETTINGS,
                                                       **(CLAUDE_ISOLATION_SETTINGS if claude_isolation else {})))
                           ] + tail
        if not mcp_config_path:
            raise ArmError('the Claude World arm needs its private MCP config path')
        return argv + ['--tools', '', '--strict-mcp-config', '--mcp-config', mcp_config_path,
                       '--allowedTools', ','.join(world_allowed_tools())] + tail
    argv = ['codex', 'exec', '--json', '--ephemeral', '--ignore-user-config', '--skip-git-repo-check',
            '-m', model, '-s', 'workspace-write' if arm == 'shell' else 'read-only', '-C', cwd]
    for feat in CODEX_DISABLED_FEATURES:
        argv += ['--disable', feat]
    if arm == 'world':
        for feat in CODEX_WORLD_DISABLED_FEATURES:
            argv += ['--disable', feat]
    argv += ['-c', 'web_search=disabled']  # a real key: `web_search=bogus` is an enum error (V-M2-1)
    if arm == 'world':
        argv += codex_world_overrides(addr)
    return argv + [prompt_text]


def redact_argv(argv: list[str], prompt_text: str, secrets: tuple[str, ...] = ()) -> list[str]:
    """The argv as it may be RECORDED: the prompt by digest, any secret or Bearer value redacted.
    (The argv carries no token by construction — the token is in the 0600 config file or the env —
    this is the belt to that brace.)"""
    out = []
    for a in argv:
        if a == prompt_text:
            out.append(f'<prompt sha256:{hashlib.sha256(prompt_text.encode()).hexdigest()}>')
        elif any(s and s in a for s in secrets) or 'bearer ' in a.lower():
            out.append('<redacted>')
        else:
            out.append(a)
    return out


def golden_argv(agent: str, arm: str) -> list[str]:
    """The argv with the §4.3 placeholders, as the goldens store it."""
    cwd = '<worktree>' if arm == 'shell' else '<void/ep>'
    mcp = '<private>/<ep>.mcp.json' if (agent, arm) == ('claude', 'world') else None
    return build_argv(agent, arm, model='$M', prompt_text='$P', cwd=cwd, mcp_config_path=mcp)
