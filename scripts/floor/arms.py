#!/usr/bin/env python3
"""Row 93 floor harness — the World arm's tool and grant set (design §4.3, §4.5; 2026-10-06).

M2 grows this module into the arm drivers (argv goldens, AC2.1). Today it pins the one thing the
rest of the harness must agree on: the World arm is EXACTLY the 8 AILANG se-tools.

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
