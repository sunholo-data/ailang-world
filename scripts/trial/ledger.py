#!/usr/bin/env python3
"""Episode evidence ledger (row 187): one row per trial/dogfood episode, from what the runner recorded.

Reads ``<root>/evidence/<project>/episodes/*.json`` (written by ``trial.py run``) and each episode's
transcript, and prints a Markdown table: wall time, World tool calls by kind, tokens, grader results,
outcome records and the PR. It reads only; nothing is written to World or the rig.

    python3 scripts/trial/ledger.py --root ~/.ailang/state/trial-2026-10-16 [--prs prs.json] [--json]

``--prs`` maps an episode id to its PR (e.g. ``{"wdocs": "ailang-world#297"}``). Token counts come
from the transcript's final usage record: Claude Code's ``result`` event (input, cache read and
write, output) or Codex's ``turn.completed``. A missing or unreadable transcript leaves them blank.

A World arm's tool calls and refusals come from World's own record when ``trial.py run`` banked it
(row 198): ``<evidence>/reports/<ep>.json``, the output of ``ailang-worldd episode report --json``
(``calls_source: world-report``). It is agent-independent, so Claude and Codex rows are counted the
same way. Without it (episodes run before row 198) they come from the transcript as before
(``calls_source: transcript``). Tokens, and every shell-arm count, always come from the transcript.

There is no comparison arm here: every row is a World episode. A row says what World cost and
delivered on real work; it does not say World beat a plain shell.
"""
from __future__ import annotations

import argparse
import datetime as dt
import glob
import json
import os
import sys

# World tool name suffix -> ledger column. Anything else is counted under "other".
KINDS = {'ailang-read': 'read', 'ailang-write': 'write', 'ailang-edit': 'edit', 'ailang-delete': 'write', 'ailang-run': 'run',
         'ailang-check': 'check', 'ailang-cli': 'cli', 'workspace-exec': 'exec',
         # the paired shell arm's native tools (Claude Code)
         'Read': 'read', 'Grep': 'read', 'Glob': 'read', 'LS': 'read', 'Write': 'write', 'Edit': 'edit',
         'MultiEdit': 'edit', 'NotebookEdit': 'edit', 'Bash': 'exec'}


def _tool_kind(name: str) -> str:
    base = name.rsplit('__', 1)[-1]
    return KINDS.get(base, 'other')


def tool_counts(tool_items: dict | None) -> dict:
    """Group a runner ``tool_items`` map (tool name -> calls) into ledger kinds."""
    out: dict = {}
    for name, n in (tool_items or {}).items():
        k = _tool_kind(name)
        out[k] = out.get(k, 0) + int(n)
    return out


def world_report(path: str, d: dict) -> dict | None:
    """The episode's banked ``episode report --json`` (row 198), or None: the row's recorded
    ``report.path``, else ``reports/<ep>.json`` beside the ``episodes/`` dir holding ``path``. Only a
    report that names this episode and carries ``calls_by_tool`` is used."""
    cands = []
    rec = d.get('report')
    if isinstance(rec, dict) and rec.get('path'):
        cands.append(rec['path'])
    cands.append(os.path.join(os.path.dirname(os.path.dirname(path)), 'reports', f"{d.get('episode')}.json"))
    for p in cands:
        try:
            with open(p, encoding='utf-8') as f:
                rep = json.load(f)
        except (OSError, ValueError):
            continue
        if isinstance(rep, dict) and rep.get('episode') == d.get('episode') and isinstance(rep.get('calls_by_tool'), dict):
            return rep
    return None


def report_refusals(rep: dict) -> dict:
    """Plan and handler refusal counts from a World report."""
    def count(k):
        v = rep.get(k)
        return int(v.get('count') or 0) if isinstance(v, dict) else 0
    return {'plan': count('plan_refusals'), 'handler': count('handler_refusals')}


def codex_tool_items(path: str | None) -> dict:
    """Codex's runner ``tool_items`` counts item TYPES (``mcp_tool_call``), so count tools from the
    transcript's completed ``mcp_tool_call`` items instead (``item.tool``). Only the fallback for a
    World row without a banked report uses this."""
    out: dict = {}
    if not path or not os.path.isfile(path):
        return out
    with open(path, encoding='utf-8', errors='replace') as f:
        for line in f:
            try:
                e = json.loads(line)
            except ValueError:
                continue
            it = e.get('item') if e.get('type') == 'item.completed' else None
            if isinstance(it, dict) and it.get('type') == 'mcp_tool_call' and it.get('tool'):
                out[it['tool']] = out.get(it['tool'], 0) + 1
    return out


def transcript_usage(path: str | None) -> dict:
    """Final token usage from a Claude Code (``result``) or Codex (``turn.completed``) transcript."""
    if not path or not os.path.isfile(path):
        return {}
    last = None
    with open(path, encoding='utf-8', errors='replace') as f:
        for line in f:
            try:
                e = json.loads(line)
            except ValueError:
                continue
            if e.get('type') in ('result', 'turn.completed') and isinstance(e.get('usage'), dict):
                last = e
    if last is None:
        return _summed_claude_usage(path)
    u = last['usage']
    return {
        'input': int(u.get('input_tokens') or 0),
        'cache_read': int(u.get('cache_read_input_tokens') or u.get('cached_input_tokens') or 0),
        'cache_write': int(u.get('cache_creation_input_tokens') or u.get('cache_write_input_tokens') or 0),
        'output': int(u.get('output_tokens') or 0),
        'cost_usd_estimate': last.get('total_cost_usd'),
    }


def _summed_claude_usage(path: str) -> dict:
    """A Claude Code transcript killed before its ``result`` event: sum the per-message usage of
    its assistant events (each message id counted once). Marked ``partial``."""
    seen, tot = set(), {'input': 0, 'cache_read': 0, 'cache_write': 0, 'output': 0}
    with open(path, encoding='utf-8', errors='replace') as f:
        for line in f:
            try:
                e = json.loads(line)
            except ValueError:
                continue
            m = e.get('message') if e.get('type') == 'assistant' else None
            if not isinstance(m, dict) or not isinstance(m.get('usage'), dict) or m.get('id') in seen:
                continue
            seen.add(m.get('id'))
            u = m['usage']
            tot['input'] += int(u.get('input_tokens') or 0)
            tot['cache_read'] += int(u.get('cache_read_input_tokens') or 0)
            tot['cache_write'] += int(u.get('cache_creation_input_tokens') or 0)
            tot['output'] += int(u.get('output_tokens') or 0)
    if not seen:
        return {}
    return dict(tot, cost_usd_estimate=None, partial=True)


def _minutes(started: str | None, finished: str | None, wall_ms) -> float | None:
    if isinstance(wall_ms, (int, float)) and wall_ms > 0:
        return round(wall_ms / 60000, 1)
    try:
        a, b = dt.datetime.fromisoformat(started), dt.datetime.fromisoformat(finished)
        return round((b - a).total_seconds() / 60, 1)
    except (TypeError, ValueError):
        return None


def episode_row(path: str, prs: dict) -> dict:
    d = json.load(open(path, encoding='utf-8'))
    run = d.get('agent_run') if isinstance(d.get('agent_run'), dict) else {}
    graders = d.get('graders') or []
    transcript = (d.get('transcript') or {}).get('path') if isinstance(d.get('transcript'), dict) else d.get('transcript')
    arm = d.get('arm') or 'world'
    rep = world_report(path, d) if arm == 'world' else None
    if rep is not None:
        tools, refusals, source = tool_counts(rep['calls_by_tool']), report_refusals(rep), 'world-report'
    else:
        items = run.get('tool_items')
        if 'mcp_tool_call' in (items or {}):  # Codex counts item types: the pre-row-198 fallback
            items = codex_tool_items(transcript)
        tools, refusals, source = tool_counts(items), None, 'transcript'
    return {
        'episode': d.get('episode'),
        'date': (d.get('started_at') or '')[:10],
        'project': d.get('project'),
        'agent': f"{d.get('agent')} {d.get('model')}",
        'arm': arm,
        'pair_of': d.get('pair_of'),
        'finish': (run.get('error_category') or run.get('finish_reason') or ''),
        'task': d.get('task'),
        'agent_rc': run.get('rc'),
        'wall_min': _minutes(d.get('started_at'), d.get('finished_at'), run.get('wall_ms')),
        'tools': tools,
        'calls_source': source,
        'refusals': refusals,
        'tokens': transcript_usage(transcript),
        'graders_passed': sum(1 for g in graders if g.get('passed')),
        'graders_total': len(graders),
        'outcomes_recorded': sum(1 for g in graders if g.get('outcome_record')),
        'outcome_record_errors': d.get('outcome_record_errors'),
        'diff_bytes': (d.get('solution') or {}).get('diff_bytes'),
        'pr': prs.get(d.get('episode'), ''),
    }


def _k(n) -> str:
    return '' if n is None else (f'{n / 1e6:.1f}M' if n >= 1e6 else f'{n / 1e3:.0f}k' if n >= 1e3 else str(n))


def markdown(rows: list) -> str:
    head = ('| episode | arm | date | project | agent | task | wall (min) | calls (read/edit+write/run/exec/other) '
            '| tokens in+cache / out | graders | outcomes | PR |')
    lines = [head, '|' + '---|' * 12]
    for r in rows:
        t, u = r['tools'], r['tokens']
        calls = sum(t.values())
        split = f"{t.get('read', 0)}/{t.get('edit', 0) + t.get('write', 0)}/{t.get('run', 0)}/{t.get('exec', 0)}/" \
                f"{t.get('check', 0) + t.get('cli', 0) + t.get('other', 0)}"
        if r.get('refusals'):
            split += f"; refused {r['refusals']['plan']} plan + {r['refusals']['handler']} handler"
        tok = f"{_k(u.get('input', 0) + u.get('cache_read', 0) + u.get('cache_write', 0))} / {_k(u.get('output'))}" if u else ''
        if u.get('partial'):
            tok += ' (partial)'
        arm = r['arm'] + (f" (pair of {r['pair_of']})" if r.get('pair_of') else '')
        wall = f"{r['wall_min']}" + (' **timeout**' if r.get('finish') == 'timeout' else '')
        lines.append(f"| {r['episode']} | {arm} | {r['date']} | {r['project']} | {r['agent']} | {r['task']} | {wall} "
                     f"| {calls} ({split}) | {tok} | {r['graders_passed']}/{r['graders_total']} "
                     f"| {r['outcomes_recorded']} | {r['pr']} |")
    return '\n'.join(lines) + '\n'


def main(argv=None) -> int:
    ap = argparse.ArgumentParser(description=__doc__.split('\n\n')[0])
    ap.add_argument('--root', required=True, action='append', help='rig root (repeatable)')
    ap.add_argument('--prs', help='JSON map episode -> PR')
    ap.add_argument('--json', action='store_true', help='print rows as JSON')
    a = ap.parse_args(argv)
    prs = json.load(open(a.prs, encoding='utf-8')) if a.prs else {}
    rows = []
    for root in a.root:
        base = os.path.join(os.path.expanduser(root), 'evidence')
        # <project>/episodes/ (current runner) and episodes/ (the first trial's layout)
        for p in glob.glob(os.path.join(base, '*', 'episodes', '*.json')) + glob.glob(os.path.join(base, 'episodes', '*.json')):
            rows.append(episode_row(p, prs))
    rows.sort(key=lambda r: (r['date'], r['episode'] or ''))
    sys.stdout.write(json.dumps(rows, indent=1) + '\n' if a.json else markdown(rows))
    return 0


if __name__ == '__main__':
    sys.exit(main())
