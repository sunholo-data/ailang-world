#!/usr/bin/env python3
"""Row 93 floor harness — one task-run's agent half (design §4.4 steps 2–4; AC2.2–AC2.4).

``run_agent`` spawns one agent in one arm and returns the driver's part of the result row:
  * World arm only: pre-flight (``world.preflight``; a failure is ``harness_setup``) and
    ``log_from`` before the spawn, ``log_to`` after it;
  * the exact argv (``arms.build_argv``), recorded REDACTED (``arms.redact_argv``);
  * a 600 s process-group deadline, ``t0``/``t1`` on the monotonic clock,
    ``finish_reason=timeout`` when the deadline kills it;
  * the raw transcript (Claude stream-json / codex ``--json``) to ``transcript_path``, parsed into
    ``tool_items`` (Claude ``tool_use`` names, codex ``item.type`` counts);
  * the ``void/<ep>`` cwd asserted still empty (World arm).

The bearer token never reaches the argv or the transcript path: the Claude World arm reads it from
a mode-0600 ``<private>/<ep>.mcp.json`` that is created just before the spawn and deleted in a
``finally``; the codex World arm gets it in the child's environment as ``WORLD_SESSION`` only. The
private dir must lie outside every worktree and void dir (refused otherwise). Grading (§4.4 step 5)
is the caller's: it reads the WORKTREE file (``world.solution_for_grading``).
"""
from __future__ import annotations

import collections
import json
import os
import signal
import subprocess
import time

import arms
import world

KILL_GRACE_S = 5.0


class DriverError(Exception):
    pass


def _within(path: str, base: str) -> bool:
    path, base = os.path.realpath(path), os.path.realpath(base)
    return path == base or path.startswith(base.rstrip(os.sep) + os.sep)


def check_private_dir(private_dir: str, *others: str) -> None:
    """The harness-private dir: absolute, mode 0700, outside (and not containing) every
    worktree/void dir the agents see."""
    if not os.path.isabs(private_dir) or not os.path.isdir(private_dir):
        raise DriverError(f'private dir must be an existing absolute dir: {private_dir}')
    mode = os.stat(private_dir).st_mode & 0o777
    if mode & 0o077:
        raise DriverError(f'private dir {private_dir} is mode {mode:o}; it must be 0700')
    for o in others:
        if o and (_within(private_dir, o) or _within(o, private_dir)):
            raise DriverError(f'private dir {private_dir} overlaps agent-visible dir {o}')


def write_private_mcp_config(private_dir: str, episode: str, token: str, addr: str) -> str:
    path = os.path.join(private_dir, f'{episode}.mcp.json')
    fd = os.open(path, os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600)
    with os.fdopen(fd, 'w') as f:
        json.dump(arms.claude_mcp_config(token, addr), f)
    return path


def spawn(argv: list, *, cwd: str, env: dict, transcript_path: str,
          deadline_s: float = arms.DEADLINE_S) -> dict:
    """Run argv in its own process group; stdout -> transcript, stderr -> transcript + '.stderr'."""
    with open(transcript_path, 'wb') as out, open(transcript_path + '.stderr', 'wb') as err:
        t0 = time.monotonic()
        p = subprocess.Popen(argv, cwd=cwd, env=env, stdin=subprocess.DEVNULL, stdout=out,
                             stderr=err, start_new_session=True)
        timed_out = False
        try:
            p.wait(timeout=deadline_s)
        except subprocess.TimeoutExpired:
            timed_out = True
            for sig, grace in ((signal.SIGTERM, KILL_GRACE_S), (signal.SIGKILL, None)):
                try:
                    os.killpg(p.pid, sig)
                except ProcessLookupError:
                    break
                try:
                    p.wait(timeout=grace)
                    break
                except subprocess.TimeoutExpired:
                    continue
            p.wait()
        t1 = time.monotonic()
    return {'rc': p.returncode, 't0': t0, 't1': t1, 'wall_ms': int(round((t1 - t0) * 1000)),
            'finish_reason': 'timeout' if timed_out else 'exit'}


def _json_lines(text: str):
    for ln in text.splitlines():
        ln = ln.strip()
        if not ln.startswith('{'):
            continue
        try:
            yield json.loads(ln)
        except ValueError:
            continue


def parse_claude(text: str) -> dict:
    """Claude stream-json: the ``system/init`` event, ``tool_use`` names, the ``result`` event."""
    init, result = None, None
    tools = collections.Counter()
    tool_errors = 0
    denied, hooks, rate = [], [], []
    for ev in _json_lines(text):
        t = ev.get('type')
        if t == 'system' and ev.get('subtype') == 'init' and init is None:
            init = {k: ev.get(k) for k in ('tools', 'mcp_servers', 'model', 'permissionMode', 'cwd',
                                          'claude_code_version', 'apiKeySource', 'skills',
                                          'slash_commands')}
            init['plugins'] = [p.get('name') for p in ev.get('plugins') or []]
        elif t == 'system' and ev.get('subtype') == 'permission_denied':
            denied.append({'tool': ev.get('tool_name'), 'reason': ev.get('decision_reason')})
        elif t == 'system' and ev.get('subtype') == 'hook_response':
            hooks.append({'hook': ev.get('hook_name'), 'outcome': ev.get('outcome'),
                          'output_bytes': len(ev.get('output') or '')})
        elif t == 'rate_limit_event':
            rate.append((ev.get('rate_limit_info') or {}).get('status'))
        elif t == 'assistant':
            for c in (ev.get('message') or {}).get('content') or []:
                if c.get('type') == 'tool_use':
                    tools[c.get('name')] += 1
        elif t == 'user':
            for c in (ev.get('message') or {}).get('content') or []:
                if isinstance(c, dict) and c.get('type') == 'tool_result' and c.get('is_error'):
                    tool_errors += 1
        elif t == 'result':
            result = {k: ev.get(k) for k in ('subtype', 'is_error', 'num_turns', 'total_cost_usd',
                                            'duration_ms', 'usage', 'result')}
    return {'init': init, 'tool_items': dict(tools), 'tool_result_errors': tool_errors,
            'permission_denied': denied, 'hooks': hooks, 'rate_limit_status': rate, 'result': result}


def parse_codex(text: str) -> dict:
    """codex ``exec --json``: the last state of every item by id, counted by ``item.type``."""
    items: dict = {}
    usage, errors = None, []
    for ev in _json_lines(text):
        t = ev.get('type') or ''
        if t.startswith('item.'):
            it = ev.get('item') or {}
            items[it.get('id') or f'_{len(items)}'] = it
        elif t == 'turn.completed':
            usage = ev.get('usage')
        elif t in ('error', 'turn.failed'):
            errors.append(ev.get('message') or (ev.get('error') or {}).get('message') or t)
    counts = collections.Counter(it.get('type') for it in items.values())
    ok_file_changes = sum(1 for it in items.values()
                          if it.get('type') == 'file_change' and it.get('status') == 'completed')
    mcp_calls = [f"{it.get('server')}/{it.get('tool')}:{it.get('status')}" for it in items.values()
                 if it.get('type') == 'mcp_tool_call']
    final = [it.get('text') for it in items.values() if it.get('type') == 'agent_message']
    return {'tool_items': dict(counts), 'file_change_ok': ok_file_changes,
            'command_execution': counts.get('command_execution', 0), 'mcp_tool_calls': mcp_calls,
            'usage': usage, 'errors': errors, 'final_message': final[-1] if final else None}


def parse_transcript(agent: str, text: str) -> dict:
    return parse_claude(text) if agent == 'claude' else parse_codex(text)


def agent_env(agent: str, arm: str, *, private_dir: str, tool_bin: str | None,
              token: str | None, base_env: dict | None = None) -> dict:
    env = dict(os.environ if base_env is None else base_env)
    env.pop(arms.TOKEN_ENV, None)
    env.pop('AILANG_REGISTRY_API_KEY', None)
    if arm == 'shell' and tool_bin:
        # §4.3: `ailang` on the shell arm's PATH is a symlink to $TOOL.
        bindir = os.path.join(private_dir, 'bin')
        os.makedirs(bindir, exist_ok=True)
        link = os.path.join(bindir, 'ailang')
        if os.path.islink(link) and os.readlink(link) != tool_bin:
            os.unlink(link)
        if not os.path.islink(link):
            os.symlink(tool_bin, link)
        env['PATH'] = bindir + os.pathsep + env.get('PATH', '')
    if agent == 'codex' and arm == 'world':
        if not token:
            raise DriverError('the codex World arm needs the session token')
        env[arms.TOKEN_ENV] = token
    return env


def run_agent(agent: str, arm: str, *, model: str, prompt_text: str, worktree: str,
              private_dir: str, transcript_path: str, void_dir: str | None = None,
              episode: str | None = None, token: str | None = None, tool_bin: str | None = None,
              addr: str = arms.DEFAULT_ADDR, deadline_s: float = arms.DEADLINE_S,
              client: world.Client | None = None, claude_isolation: bool = False,
              max_budget_usd: float | None = None) -> dict:
    """§4.4 steps 2–4 for one (agent, arm). Returns the driver fields of the row."""
    if arm == 'world' and not (void_dir and episode and token):
        raise DriverError('the World arm needs void_dir, episode and token')
    check_private_dir(private_dir, worktree, void_dir or '')
    row = {'agent': agent, 'arm': arm, 'model': model, 'episode': episode,
           'error_category': None, 'cause_evidence': None, 'world': None}
    secrets = (token,) if token else ()
    cwd = worktree if arm == 'shell' else void_dir
    mcp_path = None
    if arm == 'world':
        client = client or world.Client(addr)
        if os.listdir(void_dir):
            raise DriverError(f'void dir {void_dir} is not empty before the run')
        try:
            row['preflight'] = world.preflight(client, token)
            log_from = client.head_index() + 1
        except (world.HarnessSetup, world.WorldReadError, OSError) as e:
            row.update(error_category='harness_setup', cause_evidence=str(e))
            return row
        row['world'] = {'episode': episode, 'log_from': log_from, 'log_to': None}
    try:
        if agent == 'claude' and arm == 'world':
            mcp_path = write_private_mcp_config(private_dir, episode, token, addr)
        argv = arms.build_argv(agent, arm, model=model, prompt_text=prompt_text, cwd=cwd,
                               mcp_config_path=mcp_path, addr=addr, claude_isolation=claude_isolation,
                               max_budget_usd=max_budget_usd)
        row['argv'] = arms.redact_argv(argv, prompt_text, secrets)
        env = agent_env(agent, arm, private_dir=private_dir, tool_bin=tool_bin, token=token)
        row.update(spawn(argv, cwd=cwd, env=env, transcript_path=transcript_path,
                         deadline_s=deadline_s))
    finally:
        if mcp_path and os.path.exists(mcp_path):
            os.unlink(mcp_path)
    row['mcp_config_deleted'] = mcp_path is None or not os.path.exists(mcp_path)
    if row['finish_reason'] == 'timeout':
        row['error_category'] = 'timeout'
    if arm == 'world':
        try:
            row['world']['log_to'] = client.head_index()
        except (world.WorldReadError, OSError) as e:
            row['world']['log_to_error'] = str(e)
        row['void_listing'] = sorted(os.listdir(void_dir))
        row['void_empty'] = not row['void_listing']
    with open(transcript_path, encoding='utf-8', errors='replace') as f:
        parsed = parse_transcript(agent, f.read())
    row['transcript'] = parsed
    row['tool_items'] = parsed['tool_items']
    return row
