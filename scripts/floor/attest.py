#!/usr/bin/env python3
"""Row 93 floor harness — effective-tool attestation from the transcripts (AC2.2, AC2.3, AC2.6).

Each function returns a list of failure strings; ``[]`` means attested. They read the parsed
transcript (``run_task.parse_claude`` / ``parse_codex``), never the argv: the argv is what we asked
for, the transcript is what the agent actually had.
"""
from __future__ import annotations

import arms

WEB_TOOLS = ('WebFetch', 'WebSearch')
CODEX_WEB_ITEMS = ('web_search',)


def claude_init(init: dict | None, arm: str) -> list:
    """AC2.2 (+ AC2.6 for the shell arm): the stream-json ``system/init`` event's tool set."""
    if not init:
        return ['no system/init event in the transcript']
    tools = list(init.get('tools') or [])
    servers = [s.get('name') for s in (init.get('mcp_servers') or [])]
    bad = []
    if arm == 'shell':
        if 'Bash' not in tools:
            bad.append('shell arm: Bash is not in init tools')
        mcp = [t for t in tools if t.startswith('mcp__')]
        if mcp:
            bad.append(f'shell arm: MCP tools present {mcp}')
        web = [t for t in tools if t in WEB_TOOLS]
        if web:
            bad.append(f'shell arm: web tools present {web}')
        extra = sorted(set(tools) - set(arms.CLAUDE_SHELL_TOOLS))
        if extra:
            bad.append(f'shell arm: tools beyond {list(arms.CLAUDE_SHELL_TOOLS)}: {extra}')
        if servers:
            bad.append(f'shell arm: mcp_servers {servers} != []')
    else:
        want = sorted(arms.world_allowed_tools())
        if sorted(tools) != want:
            bad.append(f'world arm: init tools {sorted(tools)} != the 8 mcp__world__* {want}')
        if servers != [arms.MCP_SERVER]:
            bad.append(f'world arm: mcp_servers {servers} != [{arms.MCP_SERVER!r}]')
        st = [s.get('status') for s in (init.get('mcp_servers') or [])]
        if servers == [arms.MCP_SERVER] and st != ['connected']:
            bad.append(f'world arm: world server status {st} != connected')
    return bad


def claude_used_only(parsed: dict, arm: str) -> list:
    """No tool_use outside the arm's set (a canary's second line of defence)."""
    allowed = set(arms.CLAUDE_SHELL_TOOLS) if arm == 'shell' else set(arms.world_allowed_tools())
    used = sorted(set(parsed.get('tool_items') or {}) - allowed)
    return [f'{arm} arm: used tools outside the arm {used}'] if used else []


def codex_world(parsed: dict, void_listing: list) -> list:
    """AC2.3: no shell, no successful native file change, nothing in void/."""
    bad = []
    if parsed.get('command_execution'):
        bad.append(f"world arm: {parsed['command_execution']} command_execution item(s)")
    if parsed.get('file_change_ok'):
        bad.append(f"world arm: {parsed['file_change_ok']} successful file_change item(s)")
    if void_listing:
        bad.append(f'world arm: void dir not empty {void_listing}')
    bad += codex_no_web(parsed)
    return bad


def codex_no_web(parsed: dict) -> list:
    """AC2.6 (codex): no web-search item in the transcript."""
    web = {k: v for k, v in (parsed.get('tool_items') or {}).items() if k in CODEX_WEB_ITEMS}
    return [f'web items present {web}'] if web else []
