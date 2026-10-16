#!/usr/bin/env python3
"""Row 183 (w-real-project-trials, r1 ratified attended 2026-10-16) — the trial runner.

One real-project episode through World, then its outcomes (design §3 D4–D6):

  ``run``      reads the trial spec, records ``log_from = head + 1``, spawns the agent with the World
               MCP tools INCLUDING ``workspace-exec`` (D4; the floor's argv allowlists only the 8
               AILANG tools, F6), records ``log_to = head`` after it exits, runs every grader argv
               UNSANDBOXED in the episode worktree (D5) and commits one ``world/task-outcome/v1``
               per grader with the rig's ``ailang-worldd outcome record`` under the episode's own
               session. One row JSON per episode lands in ``<evidence>/episodes/<ep>.json``.
  ``pr-body``  prints a DRAFT-PR body for one episode row: the ``ailang-world`` label note, the
               D-WORLD-60 ``World-Provenance:`` trailer (as ``ailang-worldd provenance`` printed it
               during ``run``) and the outcome refs. It never pushes and never calls ``gh``: the
               operator opens the draft (D6).

Reuse, not change: the spawn, transcript streaming, private MCP config and agent env come from
``scripts/floor/run_task.py``; the argv from ``scripts/floor/arms.build_argv``. Where those hard-code
floor specifics (the 8-tool pre-flight, the 8-tool ``--allowedTools``) this file wraps them; the
floor's code, config and digest are untouched (AC6).

Token discipline (as the floor's): the session token is read from ``<rig>/private/<ep>.session``
(mode 0600) and reaches only the 0600 MCP config file (Claude, deleted after the run) or the codex
child's ``WORLD_SESSION``; graders and ``ailang-worldd`` run with ``AILANG_REGISTRY_API_KEY`` and
``WORLD_SESSION`` stripped, ``ailang-worldd`` takes the session by FILE path. Nothing written to
evidence carries the token.
"""
from __future__ import annotations

import argparse
import datetime as dt
import hashlib
import json
import os
import re
import shlex
import shutil
import signal
import subprocess
import sys
import time

HERE = os.path.dirname(os.path.abspath(__file__))
FLOOR = os.path.join(os.path.dirname(HERE), 'floor')
sys.path.insert(0, FLOOR)
import arms  # noqa: E402
import run_task  # noqa: E402
import world  # noqa: E402

TRIAL_GRADER_VERSION = 'trial-grader/1'
EXEC_TOOL = 'workspace-exec'
EXEC_ALLOWED = f'mcp__{arms.MCP_SERVER}__{EXEC_TOOL}'
AGENTS = ('claude', 'codex')
DEFAULT_DEADLINE_S = 45 * 60  # D4: 45 min per episode
DEFAULT_GRADER_TIMEOUT_S = 900
WORLDD_TIMEOUT_S = 60
HEAD_BYTES = 400  # world/task-outcome/v1 bounds stdout/stderr heads in BYTES
SECRET_ENV = ('AILANG_REGISTRY_API_KEY', arms.TOKEN_ENV)
OUTCOME_RE = re.compile(r'object (sha256:[0-9a-f]{64}) entry (\d+)')
TRAILER_RE = re.compile(r'^World-Provenance: store=\S+ episode=\S+ entries=\d+-\d+$', re.M)
BARE_HEX64_RE = re.compile(r'(?<![0-9a-fA-F:])[0-9a-fA-F]{64}(?![0-9a-fA-F])')

CHANNEL_NOTE = (
    '---\n'
    'Channel note (trial operator): work only through the World MCP tools (`mcp__world__*`); there is '
    'no shell and no native file access, and paths are relative to the episode worktree. '
    '`workspace-exec` runs the operator\'s profiled commands by id: {commands}. Each call is sandboxed '
    '(no network, short timeout). When you are done, finish by summarising the change you made.\n')


class TrialError(Exception):
    """A refusal before anything runs (bad spec, bad rig layout, unreadable session)."""


# ---------------------------------------------------------------- the spec

def _abs(base: str, p: str) -> str:
    return p if os.path.isabs(p) else os.path.normpath(os.path.join(base, p))


def load_spec(path: str) -> dict:
    with open(path) as f:
        spec = json.load(f)
    base = os.path.dirname(os.path.abspath(path))
    for k in ('project', 'episodes'):
        if not spec.get(k):
            raise TrialError(f'trial spec: missing {k!r}')
    if not re.fullmatch(r'[A-Za-z0-9._-]{1,64}', spec['project']):
        raise TrialError(f'trial spec: bad project name {spec["project"]!r}')
    if spec.get('repo'):
        spec['repo'] = _abs(base, spec['repo'])
    spec['_base'] = base
    return spec


def episode_spec(spec: dict, ep: str) -> dict:
    """The episode entry, validated, with paths resolved against the spec file's dir."""
    hits = [e for e in spec['episodes'] if e.get('id') == ep]
    if len(hits) != 1:
        raise TrialError(f'trial spec: episode {ep!r} found {len(hits)} times')
    e = dict(hits[0])
    if not re.fullmatch(r'[A-Za-z0-9._-]{1,64}', ep):
        raise TrialError(f'episode id {ep!r} is not a plain name')
    if e.get('agent') not in AGENTS:
        raise TrialError(f'episode {ep}: agent must be one of {AGENTS}')
    for k in ('model', 'task', 'task_text_file', 'worktree'):
        if not e.get(k):
            raise TrialError(f'episode {ep}: missing {k!r}')
    e['task_text_file'] = _abs(spec['_base'], e['task_text_file'])
    e['worktree'] = _abs(spec['_base'], e['worktree'])
    if e.get('void_dir'):
        e['void_dir'] = _abs(spec['_base'], e['void_dir'])
    e['deadline_s'] = float(e.get('deadline_s') or spec.get('deadline_s') or DEFAULT_DEADLINE_S)
    graders = e.get('graders') if e.get('graders') is not None else spec.get('graders')
    if not graders:
        raise TrialError(f'episode {ep}: no graders')
    names = set()
    out = []
    for g in graders:
        argv = g.get('argv')
        if not (isinstance(argv, list) and argv and all(isinstance(a, str) and a for a in argv)):
            raise TrialError(f'episode {ep}: grader argv must be a non-empty list of strings: {g!r}')
        name = g.get('name') or os.path.basename(argv[0])
        if name in names:
            raise TrialError(f'episode {ep}: duplicate grader name {name!r}')
        names.add(name)
        cwd = g.get('cwd') or '.'
        if os.path.isabs(cwd) or os.path.normpath(cwd).split(os.sep)[0] == '..':
            raise TrialError(f'episode {ep}: grader cwd {cwd!r} must be relative inside the worktree')
        out.append({'name': name, 'argv': list(argv), 'cwd': cwd,
                    'timeout_s': float(g.get('timeout_s') or DEFAULT_GRADER_TIMEOUT_S)})
    e['graders'] = out
    e['exec_commands'] = list(e.get('exec_commands') or spec.get('exec_commands') or [])
    e['label'] = f'trial-{spec["project"]}-{ep}'
    return e


def build_prompt(task_text: str, exec_commands: list) -> str:
    """D4: the task text VERBATIM, then the fixed channel note."""
    cmds = ', '.join(f'`{c}`' for c in exec_commands) if exec_commands else '(none profiled)'
    sep = '' if task_text.endswith('\n') else '\n'
    return task_text + sep + '\n' + CHANNEL_NOTE.format(commands=cmds)


# ---------------------------------------------------------------- the agent (D4)

def trial_allowed_tools() -> tuple:
    """The floor's 8 AILANG tools plus ``workspace-exec``."""
    return tuple(arms.world_allowed_tools()) + (EXEC_ALLOWED,)


def build_trial_argv(agent: str, *, model: str, prompt_text: str, cwd: str, mcp_config_path: str | None,
                     addr: str, max_budget_usd: float | None = None) -> list:
    """``arms.build_argv`` for the World arm, with the Claude ``--allowedTools`` widened by
    ``mcp__world__workspace-exec``. codex has no per-tool allowlist on the ``world`` server (the
    overrides set url, bearer env var, required and approval mode only): it sees every tool the
    session's grants serve, so a session minted with ``Workspace.Exec`` exposes ``workspace-exec``."""
    argv = arms.build_argv(agent, 'world', model=model, prompt_text=prompt_text, cwd=cwd,
                           mcp_config_path=mcp_config_path, addr=addr, max_budget_usd=max_budget_usd)
    if agent == 'claude':
        i = argv.index('--allowedTools')
        argv[i + 1] = ','.join(trial_allowed_tools())
    return argv


def trial_preflight(client, token: str) -> dict:
    """``/v1/health`` 200 and ``tools/list`` serving ``workspace-exec`` (the floor's pre-flight wants
    EXACTLY the 8 tools, so it cannot be reused here). Raises world.HarnessSetup."""
    status, h = client.health()
    if status != 200:
        raise world.HarnessSetup(f'pre-flight: /v1/health -> HTTP {status}')
    try:
        resp = client.mcp(token, 'tools/list')
    except (OSError, ValueError, world.WorldReadError) as e:
        raise world.HarnessSetup(f'pre-flight: tools/list failed: {e}') from None
    names = sorted(t.get('name') for t in ((resp.get('result') or {}).get('tools') or []))
    if EXEC_TOOL not in names:
        raise world.HarnessSetup(f'pre-flight: tools/list {names} has no {EXEC_TOOL} (session needs a Workspace.Exec grant)')
    return {'health': {k: (h or {}).get(k) for k in ('status', 'daemon_version', 'interpreter_ref', 'db_path')},
            'tools': names}


def run_trial_agent(e: dict, *, prompt_text: str, private_dir: str, void_dir: str, token: str,
                    transcript_path: str, client, addr: str, max_budget_usd: float | None = None) -> dict:
    """The trial's agent half: pre-flight, log_from, spawn (run_task.spawn), log_to, parse."""
    agent, ep = e['agent'], e['id']
    run_task.check_private_dir(private_dir, e['worktree'], void_dir)
    if os.listdir(void_dir):
        raise TrialError(f'void dir {void_dir} is not empty before the run')
    row = {'error_category': None, 'cause_evidence': None, 'log_from': None, 'log_to': None}
    try:
        row['preflight'] = trial_preflight(client, token)
        row['log_from'] = client.head_index() + 1
    except (world.HarnessSetup, world.WorldReadError, OSError) as err:
        row.update(error_category='harness_setup', cause_evidence=str(err))
        return row
    mcp_path = None
    try:
        if agent == 'claude':
            mcp_path = run_task.write_private_mcp_config(private_dir, ep, token, addr)
        argv = build_trial_argv(agent, model=e['model'], prompt_text=prompt_text, cwd=void_dir,
                                mcp_config_path=mcp_path, addr=addr, max_budget_usd=max_budget_usd)
        row['argv'] = arms.redact_argv(argv, prompt_text, (token,))
        env = run_task.agent_env(agent, 'world', private_dir=private_dir, tool_bin=None, token=token)
        row.update(run_task.spawn(argv, cwd=void_dir, env=env, transcript_path=transcript_path,
                                  deadline_s=e['deadline_s']))
    finally:
        if mcp_path and os.path.exists(mcp_path):
            os.unlink(mcp_path)
    row['mcp_config_deleted'] = mcp_path is None or not os.path.exists(mcp_path)
    if row.get('finish_reason') == 'timeout':
        row['error_category'] = 'timeout'
    try:
        row['log_to'] = client.head_index()
    except (world.WorldReadError, OSError) as err:
        row['log_to_error'] = str(err)
    row['void_listing'] = sorted(os.listdir(void_dir))
    with open(transcript_path, encoding='utf-8', errors='replace') as f:
        row['transcript'] = run_task.parse_transcript(agent, f.read())
    return row


# ---------------------------------------------------------------- graders + outcomes (D5)

def _head_bytes(b: bytes, limit: int = HEAD_BYTES) -> str:
    """At most ``limit`` UTF-8 bytes, cut on a character boundary; invalid bytes become U+FFFD
    first (as Go decodes them), so the bound holds on the re-encoded string."""
    s = (b or b'').decode('utf-8', errors='replace').encode('utf-8')
    return s[:limit].decode('utf-8', errors='ignore')


def child_env(base: dict | None = None) -> dict:
    env = dict(os.environ if base is None else base)
    for k in SECRET_ENV:
        env.pop(k, None)
    return env


def resolve_tool(argv0: str, cwd: str, env: dict) -> str | None:
    if os.sep in argv0:
        p = _abs(cwd, argv0)
        return os.path.realpath(p) if os.path.isfile(p) else None
    p = shutil.which(argv0, path=env.get('PATH'))
    return os.path.realpath(p) if p else None


def sha256_file(path: str | None) -> str:
    if not path:
        return ''
    h = hashlib.sha256()
    with open(path, 'rb') as f:
        for chunk in iter(lambda: f.read(1 << 20), b''):
            h.update(chunk)
    return h.hexdigest()


def run_grader(g: dict, worktree: str, env: dict | None = None) -> dict:
    """One grader argv, unsandboxed, in ``<worktree>/<cwd>``, own process group, wall-capped."""
    env = child_env() if env is None else env
    cwd = os.path.normpath(os.path.join(worktree, g['cwd']))
    tool = resolve_tool(g['argv'][0], cwd, env)
    out = {'name': g['name'], 'argv': g['argv'], 'cwd': g['cwd'], 'timeout_s': g['timeout_s'],
           'tool_path': tool, 'tool_sha256': sha256_file(tool), 'rc': None, 'timed_out': False,
           'spawn_error': None}
    t0 = time.monotonic()
    stdout = stderr = b''
    try:
        p = subprocess.Popen(g['argv'], cwd=cwd, env=env, stdin=subprocess.DEVNULL, stdout=subprocess.PIPE,
                             stderr=subprocess.PIPE, start_new_session=True)
    except OSError as err:
        out['spawn_error'] = str(err)
        stderr = str(err).encode()
    else:
        try:
            stdout, stderr = p.communicate(timeout=g['timeout_s'])
        except subprocess.TimeoutExpired:
            out['timed_out'] = True
            for sig in (signal.SIGTERM, signal.SIGKILL):
                try:
                    os.killpg(p.pid, sig)
                except ProcessLookupError:
                    break
                try:
                    stdout, stderr = p.communicate(timeout=run_task.KILL_GRACE_S)
                    break
                except subprocess.TimeoutExpired:
                    continue
            if p.poll() is None:
                p.kill()
                p.wait()
        out['rc'] = p.returncode if not out['timed_out'] else None
    out['wall_ms'] = int(round((time.monotonic() - t0) * 1000))
    out['passed'] = out['rc'] == 0
    out['stdout_head'] = _head_bytes(stdout)
    out['stderr_head'] = _head_bytes(stderr)
    out['_stdout'], out['_stderr'] = stdout, stderr
    return out


def worktree_state(worktree: str) -> dict:
    """The agent's change, hashed: ``git diff HEAD --binary`` plus every untracked, unignored file
    (path and bytes). Taken after the agent and BEFORE any grader (graders may write artefacts)."""
    def git(*a):
        return subprocess.run(['git', '-C', worktree, *a], capture_output=True, timeout=120, check=True).stdout
    try:
        head = git('rev-parse', 'HEAD').strip().decode()
        diff = git('diff', 'HEAD', '--binary')
        untracked = sorted(p for p in git('ls-files', '--others', '--exclude-standard', '-z').split(b'\0') if p)
    except (OSError, subprocess.SubprocessError) as err:
        return {'sha256': '', 'state': 'unknown', 'error': str(err)[:300]}
    h = hashlib.sha256(diff)
    for rel in untracked:
        h.update(b'\0untracked\0' + rel + b'\0')
        try:
            h.update(bytes.fromhex(sha256_file(os.path.join(worktree, rel.decode('utf-8', 'surrogateescape')))))
        except OSError:
            h.update(b'<unreadable>')
    changed = bool(diff) or bool(untracked)
    return {'sha256': h.hexdigest() if changed else '', 'state': 'changed' if changed else 'unchanged',
            'base': head, 'diff_bytes': len(diff), 'untracked': len(untracked)}


def grading_rule(g: dict) -> str:
    where = '' if g['cwd'] in ('.', '') else f' (cwd {g["cwd"]})'
    return f'the operator runs {shlex.join(g["argv"])} unsandboxed in the episode worktree{where} after the agent stops'


def outcome_payload(e: dict, log_from: int | None, log_to: int | None, gr: dict, sol: dict) -> dict:
    """One world/task-outcome/v1 for one grader (pure; the verb validates it again)."""
    if log_from is None or log_to is None:
        raise ValueError('world episode without a log range')
    ok = bool(gr['passed'])
    return {
        'agent': e['agent'], 'arm': 'world', 'episodeId': e['id'], 'run': 1, 'label': e['label'], 'task': e['task'],
        'entries': {'from': int(log_from), 'to': int(log_to)},
        'outcome': {'class': 'PASS' if ok else 'FAIL', 'errorCategory': '' if ok else 'grader_failed',
                    # A project grader has no compile/run/stdout split: all three carry its verdict.
                    'compileOk': ok, 'runtimeOk': ok, 'stdoutOk': ok,
                    'stdoutHead': _head_bytes(gr['stdout_head'].encode()),
                    'stderrHead': _head_bytes(gr['stderr_head'].encode()),
                    'solutionSha256': sol.get('sha256') or '', 'solutionState': sol.get('state') or ''},
        'grading': {'rule': grading_rule(gr), 'argv': list(gr['argv']), 'inputFiles': [], 'expectedStdout': '',
                    'graderVersion': TRIAL_GRADER_VERSION, 'toolSha256': gr.get('tool_sha256') or ''},
    }


def worldd_argv(rig: str, addr: str, *verb) -> list:
    return [os.path.join(rig, 'bin', 'ailang-worldd'), '--addr', f'http://{addr}', *verb]


def record_outcome(rig: str, addr: str, ep: str, payload: dict, runner=subprocess.run) -> dict:
    """Commit one outcome under the episode's session FILE. {'ref','entry'} or {'error'}; never raises."""
    argv = worldd_argv(rig, addr, 'outcome', 'record', '--file', '-',
                       '--session', os.path.join(rig, 'private', f'{ep}.session'))
    try:
        p = runner(argv, input=json.dumps(payload).encode(), capture_output=True, env=child_env(),
                   timeout=WORLDD_TIMEOUT_S)
    except (OSError, subprocess.SubprocessError) as err:
        return {'error': f'spawn: {err}'}
    m = OUTCOME_RE.search(p.stdout.decode(errors='replace'))
    if p.returncode != 0 or not m:
        return {'error': f'rc={p.returncode}: {p.stderr.decode(errors="replace").strip()[:300]}'}
    return {'ref': m.group(1), 'entry': int(m.group(2))}


def provenance_trailer(rig: str, addr: str, ep: str, log_from: int, log_to: int, runner=subprocess.run) -> dict:
    """The D-WORLD-60 trailer exactly as ``ailang-worldd provenance`` prints it. Never raises."""
    if log_from is None or log_to is None or log_to < log_from:
        return {'error': 'no World entries in the episode range'}
    scan = min(max(500, log_to - log_from + 1), 5000)
    argv = worldd_argv(rig, addr, 'provenance', '--since', str(log_from), '--episode', ep, '--scan', str(scan))
    try:
        p = runner(argv, capture_output=True, env=child_env(), timeout=WORLDD_TIMEOUT_S)
    except (OSError, subprocess.SubprocessError) as err:
        return {'error': f'spawn: {err}'}
    m = TRAILER_RE.search(p.stdout.decode(errors='replace'))
    if p.returncode != 0 or not m:
        return {'error': f'rc={p.returncode}: {p.stderr.decode(errors="replace").strip()[:300]}'}
    return {'trailer': m.group(0)}


def outcome_record_errors(row: dict) -> int:
    return sum(1 for g in row.get('graders') or [] if (g.get('outcome_record') or {}).get('error'))


# ---------------------------------------------------------------- one episode

def read_token(path: str) -> str:
    st = os.stat(path)
    if st.st_mode & 0o077:
        raise TrialError(f'{path} is not mode 0600')
    with open(path) as f:
        tok = f.read().strip()
    if not re.fullmatch(r'[0-9a-f]{64}', tok):
        raise TrialError(f'{path} does not hold a 64-hex session token')
    return tok


def _sha_file_or_none(p: str) -> str | None:
    return sha256_file(p) if os.path.isfile(p) else None


def run_episode(spec: dict, ep: str, *, rig: str, addr: str, evidence: str, client=None,
                runner=subprocess.run, grader_env: dict | None = None) -> dict:
    """D4 + D5 for one episode; writes and returns ``<evidence>/episodes/<ep>.json``."""
    e = episode_spec(spec, ep)
    private = os.path.join(rig, 'private')
    void_dir = e.get('void_dir') or os.path.join(rig, 'void', ep)
    os.makedirs(void_dir, exist_ok=True)
    token = read_token(os.path.join(private, f'{ep}.session'))
    for d in ('episodes', 'transcripts', 'graders'):
        os.makedirs(os.path.join(evidence, d), exist_ok=True)
    with open(e['task_text_file'], encoding='utf-8') as f:
        task_text = f.read()
    prompt_text = build_prompt(task_text, e['exec_commands'])
    transcript = os.path.join(evidence, 'transcripts', f'{ep}.jsonl')
    client = client or world.Client(addr)
    row = {'schema': 'trial-episode/1', 'project': spec['project'], 'episode': ep, 'agent': e['agent'],
           'model': e['model'], 'task': e['task'], 'label': e['label'], 'worktree': e['worktree'],
           'started_at': dt.datetime.now(dt.timezone.utc).isoformat(timespec='seconds'),
           'prompt_sha256': hashlib.sha256(prompt_text.encode()).hexdigest(),
           'task_text_sha256': hashlib.sha256(task_text.encode()).hexdigest(),
           'exec_commands': e['exec_commands'], 'deadline_s': e['deadline_s']}
    drv = run_trial_agent(e, prompt_text=prompt_text, private_dir=private, void_dir=void_dir, token=token,
                          transcript_path=transcript, client=client, addr=addr,
                          max_budget_usd=e.get('max_budget_usd'))
    tr = drv.pop('transcript', None) or {}
    row['agent_run'] = {k: drv.get(k) for k in ('rc', 'wall_ms', 'finish_reason', 'error_category',
                                                 'cause_evidence', 'argv', 'preflight', 'mcp_config_deleted',
                                                 'void_listing', 'log_to_error', 'stdout_reader_done')}
    row['agent_run']['tool_items'] = tr.get('tool_items')
    row['agent_summary'] = ((tr.get('result') or {}).get('result') if e['agent'] == 'claude'
                            else tr.get('final_message'))
    row['transcript'] = {'path': transcript, 'sha256': _sha_file_or_none(transcript)}
    row['log_range'] = {'from': drv.get('log_from'), 'to': drv.get('log_to')}
    row['graders'] = []
    if drv.get('error_category') == 'harness_setup':
        row['outcome_record_errors'] = 0
        return _write_row(evidence, ep, row)
    sol = worktree_state(e['worktree'])
    row['solution'] = sol
    for g in e['graders']:
        gr = run_grader(g, e['worktree'], grader_env)
        base = os.path.join(evidence, 'graders', f'{ep}.{re.sub(r"[^A-Za-z0-9._-]", "_", g["name"])}')
        for stream in ('stdout', 'stderr'):
            with open(f'{base}.{stream}', 'wb') as f:
                f.write(gr.pop('_' + stream))
            gr[f'{stream}_path'] = f'{base}.{stream}'
        try:
            payload = outcome_payload(e, drv.get('log_from'), drv.get('log_to'), gr, sol)
        except (ValueError, KeyError, TypeError) as err:
            gr['outcome_record'] = {'error': f'payload: {err}'}
        else:
            try:
                gr['outcome_record'] = record_outcome(rig, addr, ep, payload, runner=runner)
            except Exception as err:  # noqa: BLE001 — a record must never take the episode down
                gr['outcome_record'] = {'error': f'recorder: {err}'}
        row['graders'].append(gr)
    row['outcome_record_errors'] = outcome_record_errors(row)
    row['provenance'] = provenance_trailer(rig, addr, ep, drv.get('log_from'), drv.get('log_to'), runner=runner)
    return _write_row(evidence, ep, row)


def _write_row(evidence: str, ep: str, row: dict) -> dict:
    row['finished_at'] = dt.datetime.now(dt.timezone.utc).isoformat(timespec='seconds')
    path = os.path.join(evidence, 'episodes', f'{ep}.json')
    with open(path, 'w') as f:
        json.dump(row, f, indent=2, sort_keys=True, ensure_ascii=False, default=str)
        f.write('\n')
    row['_path'] = path
    return row


# ---------------------------------------------------------------- the draft PR body (D6)

def pr_body(row: dict, evidence_link: str | None = None) -> str:
    """The draft-PR body. Refuses (TrialError) without a trailer, or if anything token-shaped
    (a bare 64-hex run, ``Bearer``) would be printed."""
    trailer = (row.get('provenance') or {}).get('trailer')
    if not trailer or not TRAILER_RE.fullmatch(trailer):
        raise TrialError(f'episode {row.get("episode")}: no World-Provenance trailer recorded '
                         f'({(row.get("provenance") or {}).get("error") or "absent"})')
    lines = [f'Trial episode `{row["episode"]}` ({row["project"]}, task `{row["task"]}`): changes made by '
             f'{row["agent"]} `{row["model"]}` working through AILANG World only (row 183, w-real-project-trials).',
             '', 'Label: `ailang-world` (D-WORLD-60). Draft: the maintainer reviews and merges.', '']
    summary = (row.get('agent_summary') or '').strip()
    if summary:
        lines += ['## Agent summary', '', summary[:4000], '']
    lines += ['## Operator grading (unsandboxed, after the agent stopped)', '']
    for g in row.get('graders') or []:
        rec = g.get('outcome_record') or {}
        ref = f'outcome `{rec["ref"]}` (entry {rec["entry"]})' if rec.get('ref') else f'outcome NOT recorded: {rec.get("error")}'
        verdict = 'PASS' if g.get('passed') else ('TIMEOUT' if g.get('timed_out') else f'FAIL rc={g.get("rc")}')
        lines.append(f'- `{shlex.join(g["argv"])}` (cwd `{g["cwd"]}`): **{verdict}**; {ref}')
    lr = row.get('log_range') or {}
    tsha = (row.get('transcript') or {}).get('sha256')
    lines += ['', f'World entries {lr.get("from")}-{lr.get("to")}; walk any of them with `ailang-worldd why <entry>`.']
    if tsha:
        lines.append(f'Agent transcript kept on the trial rig, sha256:{tsha}.')
    if evidence_link:
        lines.append(f'Evidence: {evidence_link}')
    lines += ['', trailer]
    body = '\n'.join(lines) + '\n'
    if BARE_HEX64_RE.search(body) or 'bearer ' in body.lower():
        raise TrialError('pr-body: refusing to print a token-shaped value')
    return body


# ---------------------------------------------------------------- CLI

def main(argv=None) -> int:
    ap = argparse.ArgumentParser(prog='trial.py', description=__doc__.split('\n\n')[0])
    sub = ap.add_subparsers(dest='cmd', required=True)
    r = sub.add_parser('run', help='run one episode: agent, graders, outcomes')
    r.add_argument('--trial', required=True)
    r.add_argument('--episode', required=True)
    r.add_argument('--root', required=True, help='the trial rig (bin/ailang-worldd, private/<ep>.session)')
    r.add_argument('--addr', default=arms.DEFAULT_ADDR, help='daemon host:port')
    r.add_argument('--evidence', help='evidence dir (default <root>/evidence/<project>)')
    p = sub.add_parser('pr-body', help='print the draft-PR body for one episode row (never pushes)')
    p.add_argument('--episode-row', required=True)
    p.add_argument('--evidence-link', help='the banked evidence dir/URL to link (D-WORLD-60 (3))')
    a = ap.parse_args(argv)
    try:
        if a.cmd == 'run':
            spec = load_spec(a.trial)
            rig = os.path.abspath(a.root)
            evidence = os.path.abspath(a.evidence or os.path.join(rig, 'evidence', spec['project']))
            row = run_episode(spec, a.episode, rig=rig, addr=a.addr, evidence=evidence)
            print(json.dumps({'episode_row': row['_path'], 'agent_rc': row['agent_run'].get('rc'),
                              'error_category': row['agent_run'].get('error_category'),
                              'log_range': row['log_range'],
                              'graders': {g['name']: g['passed'] for g in row['graders']},
                              'outcome_record_errors': row['outcome_record_errors'],
                              'provenance': row.get('provenance')}, indent=2))
            if row['agent_run'].get('error_category') == 'harness_setup':
                return 2
            return 1 if row['outcome_record_errors'] else 0
        with open(a.episode_row) as f:
            sys.stdout.write(pr_body(json.load(f), a.evidence_link))
        return 0
    except TrialError as err:
        print(f'trial.py: {err}', file=sys.stderr)
        return 2


if __name__ == '__main__':
    sys.exit(main())
