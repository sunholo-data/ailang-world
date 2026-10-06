#!/usr/bin/env python3
"""Row 93 floor harness — the phase/iteration runner (design §4.4, §4.9–§4.11; AC4.1, AC4.2, AC5.1).

    run.py iterate --iteration K --knob K --side S --change TEXT [--reuse-shell-from DIR]
        One smoke-tuning iteration (§4.10): canary per (agent, arm) (§4.9), then every smoke task
        N=1 per arm for both agents on the tuning models; grade, classify (§4.6), append JSONL rows
        under <evidence>/iter-K/runs/<agent>/<arm>/r1.jsonl, then append ONE ledger row (AC4.1).
    run.py digest [--mode smoke|final]       print the live config and its digest
    run.py prereg-draft [--out PATH]         write the FINAL preregistration draft (M5 commits it)
    run.py final --prereg PATH --phase shell|world --run K
        M5 only. REFUSES unless the prereg is committed (git ls-files + clean against HEAD) and its
        ``config_digest`` equals the live FINAL config's digest (AC5.1, MUT-FINAL-NOPREREG).

Refusals that guard the record:
  * AC4.2 — ``iterate`` and ``final`` refuse while the tuning ledger has uncommitted changes
    (``git status --porcelain`` on it is non-empty): every ledger row is committed before the next
    iteration runs, so a ledger row's ``commit`` names the code that produced it.
  * ``iterate`` also refuses while ``scripts/floor/`` is dirty (the config and harness under test
    must be the committed ones the ledger row cites).

The tuned config is ``scripts/floor/floor_config.json``. Its digest (``config_digest``) covers the
config AND what it renders: the argv template per (agent, arm), the teaching-prompt and template
digests, the World tools block text, the tool binary's sha256, the corpus commit and the YAML
hashes of the mode's tasks, the grader version and the CLI versions.

Cost guard (§4.9): Claude runs carry ``--max-budget-usd``; ``iterate`` stops before a Claude run
that could push the summed ``total_cost_usd`` past ``--spend-cap`` (ledger + this iteration), and
aborts on the first quota/rate-limit evidence. codex has no budget flag: the 600 s per-task
deadline plus ``--max-wall-s`` bound it.
"""
from __future__ import annotations

import argparse
import datetime as dt
import hashlib
import json
import os
import re
import shutil
import subprocess
import sys
import time

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import arms  # noqa: E402
import attest  # noqa: E402
import classify as classify_mod  # noqa: E402
import grade  # noqa: E402
import prompt  # noqa: E402
import stats  # noqa: E402

REPO = os.path.abspath(os.path.join(HERE, '..', '..'))
CONFIG = os.path.join(HERE, 'floor_config.json')
LEDGER = os.path.join(REPO, 'design_docs', 'verification', 'world-floor-tuning-ledger.jsonl')
EVIDENCE = os.path.join(REPO, 'design_docs', 'verification', 'world-floor-m4-2026-10-06')
DEFAULT_ROOT = os.path.expanduser('~/.ailang/state/floor-tune')
AILANG_REPO = os.environ.get('FLOOR_AILANG_REPO') or os.path.expanduser('~/dev/sunholo-data/ailang')
EPISODES = {'claude': 'fl-tune-cc', 'codex': 'fl-tune-cx'}
SOL = 'benchmark/solution.ail'
MOCK_PORT = 7655  # the one loopback pair the scratch daemon's ailang-run may reach (canary.serve_argv)
TEACHING_MODES = ('full', 'compact')
MODES = ('smoke', 'final')
PREREG_KIND = 'row93-floor-preregistration'
# Row 153 (the known product defect): a World tool call whose host callback outlives the daemon's
# plan-phase deadline. Counted, never tuned around (§4.6: World pays for its own faults).
ROW153_PATTERNS = ('host callback timed out', 'execute plan phase: context deadline exceeded',
                   'context deadline exceeded')
QUOTA_PATTERNS = (('quota_exhausted', re.compile(r'(?i)usage limit|quota|insufficient[_ ]credit|credit balance')),
                  ('rate_limit', re.compile(r'(?i)rate[_ -]?limit|\b429\b|too many requests')))
TYPED_API = re.compile(r'(?i)overloaded|\b5\d\d\b|internal server error|api_error|service unavailable')


class RunRefused(Exception):
    """A guard refused the run (AC4.2 dirty ledger, AC5.1 prereg, spend cap, …)."""


class QuotaAbort(Exception):
    """§4.9: the first quota_exhausted / rate_limit evidence aborts the phase."""


# ---------------------------------------------------------------- git guards (AC4.2, AC5.1)

def _git(args, cwd):
    return subprocess.run(['git', *args], cwd=cwd, capture_output=True, text=True)


def dirty(path: str) -> str:
    """``git status --porcelain`` for one path ('' = clean). Untracked counts as dirty."""
    d = path if os.path.isdir(path) else os.path.dirname(os.path.abspath(path))
    p = _git(['status', '--porcelain', '--untracked-files=all', '--', os.path.abspath(path)], d)
    if p.returncode != 0:
        raise RunRefused(f'git status failed for {path}: {p.stderr.strip()}')
    return p.stdout.strip()


def require_clean_ledger(ledger: str = LEDGER) -> None:
    """AC4.2: refuse to run while the tuning ledger has uncommitted changes."""
    st = dirty(ledger)
    if st:
        raise RunRefused(f'AC4.2: the tuning ledger {ledger} has uncommitted changes ({st!r}); '
                         'commit the previous iteration\'s ledger row before running another')


def require_committed(path: str) -> None:
    """The file is tracked and its working copy equals HEAD's blob (report.require_sealed's rule)."""
    path = os.path.abspath(path)
    if not os.path.isfile(path):
        raise RunRefused(f'{path} does not exist')
    d, base = os.path.dirname(path), os.path.basename(path)
    if _git(['ls-files', '--error-unmatch', '--', base], d).returncode != 0:
        raise RunRefused(f'{path} is not committed (git ls-files --error-unmatch)')
    if _git(['diff', '--quiet', 'HEAD', '--', base], d).returncode != 0:
        raise RunRefused(f'{path} differs from its committed blob (git diff --quiet HEAD)')


def head_commit(cwd: str = REPO) -> str:
    return _git(['rev-parse', 'HEAD'], cwd).stdout.strip()


# ---------------------------------------------------------------- config + live digest

def load_config(path: str = CONFIG) -> dict:
    with open(path) as f:
        cfg = json.load(f)
    k = cfg.get('knobs') or {}
    if k.get('world_tools_wording') not in prompt.WORLD_WORDINGS:
        raise RunRefused(f'knobs.world_tools_wording must be one of {prompt.WORLD_WORDINGS}')
    if k.get('teaching') not in TEACHING_MODES:
        raise RunRefused(f'knobs.teaching must be one of {TEACHING_MODES}')
    if not isinstance(k.get('claude_isolation'), bool):
        raise RunRefused('knobs.claude_isolation must be a bool')
    for m in MODES:
        sec = cfg.get(m) or {}
        if set((sec.get('models') or {})) != set(arms.AGENTS):
            raise RunRefused(f'{m}.models must name exactly {arms.AGENTS}')
        if not (sec.get('claude_max_budget_usd') or 0) > 0:
            raise RunRefused(f'{m}.claude_max_budget_usd must be positive')
    # Forbidden knobs (§4.10) are fixed here, not in the config: N for FINAL is D-NF-4's 3, the gate
    # tier is core, the deadline is §4.3's 600 s.
    if cfg['final'].get('N') != 3 or cfg['final'].get('tier') != 'core':
        raise RunRefused('final.N must be 3 and final.tier core (D-NF-4; forbidden knobs, §4.10)')
    if cfg.get('deadline_s') != arms.DEADLINE_S:
        raise RunRefused(f'deadline_s must be {arms.DEADLINE_S} (§4.3)')
    cfg['_path'] = os.path.abspath(path)
    return cfg


def tool_bin(cfg: dict) -> str:
    return os.path.expanduser(cfg['tool_bin'])


def mode_section(cfg: dict, mode: str) -> dict:
    if mode not in MODES:
        raise RunRefused(f'unknown mode {mode!r}')
    return cfg[mode]


def _sha(b: bytes) -> str:
    return hashlib.sha256(b).hexdigest()


def canonical(obj) -> bytes:
    return json.dumps(obj, sort_keys=True, separators=(',', ':'), ensure_ascii=False).encode()


def argv_template(cfg: dict, mode: str, agent: str, arm: str) -> list:
    sec = mode_section(cfg, mode)
    k = cfg['knobs']
    cwd = '<worktree>' if arm == 'shell' else '<void/ep>'
    mcp = '<private>/<ep>.mcp.json' if (agent, arm) == ('claude', 'world') else None
    return arms.build_argv(agent, arm, model=sec['models'][agent], prompt_text='$P', cwd=cwd,
                           mcp_config_path=mcp, claude_isolation=k['claude_isolation'],
                           max_budget_usd=sec['claude_max_budget_usd'] if agent == 'claude' else None)


class Probes:
    """Everything the live config reads from outside the config file (injectable for tests)."""

    def __init__(self, ailang_repo: str = AILANG_REPO):
        self.ailang_repo = ailang_repo
        self._cache = {}

    def _memo(self, key, fn):
        if key not in self._cache:
            self._cache[key] = fn()
        return self._cache[key]

    def tool_sha256(self, path):
        return self._memo(('tool', path), lambda: grade.sha256_file(grade.resolve_tool(path)))

    def teaching(self, path, compact):
        return self._memo(('teach', path, compact), lambda: prompt.load_teaching(path, compact))

    def template(self, commit):
        return self._memo(('tmpl', commit), lambda: prompt.load_template(self.ailang_repo, commit))

    def corpus(self, commit, tier, expected):
        import corpus as corpus_mod
        return self._memo(('corpus', commit, tier, expected), lambda: corpus_mod.load_corpus(
            repo=self.ailang_repo, commit=commit, tier=tier, expected=expected))

    def cli_versions(self):
        def v(argv):
            try:
                return subprocess.run(argv, capture_output=True, text=True, timeout=60).stdout.strip()
            except (OSError, subprocess.TimeoutExpired) as e:
                return f'error: {e}'
        return self._memo('cli', lambda: {'claude': v(['claude', '--version']), 'codex': v(['codex', '--version'])})


def mode_tasks(cfg: dict, mode: str, probes: Probes) -> list:
    """The task specs of a mode, in id order. FINAL = all 23 core tasks (forbidden to change)."""
    sec = mode_section(cfg, mode)
    core = probes.corpus(cfg['corpus_commit'], 'core', 23)
    by = {t.id: t for t in core.tasks}
    if mode == 'final':
        return sorted(core.tasks, key=lambda t: t.id)
    ids = sec['tasks']
    missing = [i for i in ids if i not in by]
    if missing or len(set(ids)) != len(ids):
        raise RunRefused(f'smoke tasks must be distinct core ids; unknown {missing}')
    return [by[i] for i in sorted(ids)]


def canary_task(cfg: dict, mode: str, probes: Probes):
    sec = mode_section(cfg, mode)
    tier = sec['canary_tier']
    c = probes.corpus(cfg['corpus_commit'], tier, 23)
    by = {t.id: t for t in c.tasks}
    if sec['canary_task'] not in by:
        raise RunRefused(f'canary task {sec["canary_task"]!r} is not in tier {tier}')
    return by[sec['canary_task']]


def arm_digest(cfg: dict, mode: str, agent: str, arm: str, probes: Probes) -> str:
    """What one (agent, arm) cell depends on: shell rows are reusable across a World-only knob."""
    k = cfg['knobs']
    d = {'argv': argv_template(cfg, mode, agent, arm),
         'teaching_sha256': _sha(probes.teaching(tool_bin(cfg), k['teaching'] == 'compact').encode()),
         'template_sha256': _sha(probes.template(cfg['corpus_commit']).encode()),
         'tasks': [(t.id, t.sha256) for t in mode_tasks(cfg, mode, probes)],
         'tool_sha256': probes.tool_sha256(tool_bin(cfg)), 'deadline_s': cfg['deadline_s']}
    if arm == 'world':
        d['world_tools_wording'] = k['world_tools_wording']
        d['world_tools_block'] = prompt.world_tools_block(*reversed(prompt.world_wording(k['world_tools_wording'])[:2]))
    return _sha(canonical(d))


def live_config(cfg: dict, mode: str, probes: Probes) -> dict:
    sec = mode_section(cfg, mode)
    k = cfg['knobs']
    tasks = mode_tasks(cfg, mode, probes)
    teach = probes.teaching(tool_bin(cfg), k['teaching'] == 'compact')
    w_tools, w_note, _ = prompt.world_wording(k['world_tools_wording'])
    cfg_clean = {kk: v for kk, v in cfg.items() if not kk.startswith('_')}
    return {
        'mode': mode, 'config': cfg_clean,
        'agents': list(arms.AGENTS), 'models': dict(sec['models']), 'N': sec['N'],
        'argv_templates': {f'{a}/{x}': argv_template(cfg, mode, a, x) for a in arms.AGENTS for x in arms.ARMS},
        'prompt': {'template_sha256': _sha(probes.template(cfg['corpus_commit']).encode()),
                   'teaching': k['teaching'], 'teaching_sha256': _sha(teach.encode()),
                   'teaching_bytes': len(teach.encode()),
                   'world_tools_wording': k['world_tools_wording'],
                   'world_tools_block_sha256': _sha(prompt.world_tools_block(w_note, w_tools).encode())},
        'corpus': {'commit': cfg['corpus_commit'], 'tasks': {t.id: t.sha256 for t in tasks},
                   'canary': {canary_task(cfg, mode, probes).id: canary_task(cfg, mode, probes).sha256}},
        'tool_bin_sha256': probes.tool_sha256(tool_bin(cfg)), 'grader': grade.GRADER_VERSION,
        'cli_versions': probes.cli_versions(), 'deadline_s': cfg['deadline_s'],
        'world_grants': arms.world_grant_args(), 'world_tools': list(prompt.WORLD_TOOL_NAMES),
    }


def config_digest(live: dict) -> str:
    return _sha(canonical(live))


def ledger_head(ledger: str = LEDGER) -> dict:
    rows = read_ledger(ledger)
    raw = b''
    if os.path.exists(ledger):
        with open(ledger, 'rb') as f:
            raw = f.read()
    return {'rows': len(rows), 'sha256': _sha(raw), 'last_iteration': rows[-1]['iteration'] if rows else None}


def prereg_draft(cfg: dict, probes: Probes, ledger: str = LEDGER) -> dict:
    live = live_config(cfg, 'final', probes)
    return {
        'kind': PREREG_KIND, 'status': 'DRAFT — M5 commits it (attended); FINAL refuses otherwise',
        'config_digest': config_digest(live), 'live_config': live, 'tuning_ledger_head': ledger_head(ledger),
        'store': 'LIVE World store (D-NF-6); tuning used the scratch store ~/.ailang/state/floor-tune',
        'episodes': 'one per (agent, run): fl-cc-r<k>, fl-cx-r<k> (§4.5), 8 explicit grants each, --ttl 43200',
        'classification': 'design §4.6, scripts/floor/classify.py (pre-registered table; unlisted categories refused)',
        'statistics': {'T': 23, 'N': 3, 'eligible': 'N>=3 and zero HARNESS_FAULT in S and max_k r_S,k - min_k r_S,k <= 0.05',
                       'pass_rate': 'Delta = R_W - R_S >= -0.02 (point estimate; D-WORLD-55)',
                       'overhead': 'O = median_t (median_k wall_W / median_k wall_S - 1) <= 0.25; timeout counts 600 s',
                       'thresholds': {'range_max': str(stats.RANGE_MAX), 'delta_min': str(stats.DELTA_MIN),
                                      'overhead_max': str(stats.OVERHEAD_MAX), 'min_runs': stats.MIN_RUNS},
                       'informative_only': 'pooled-median ratio, total-time ratio, paired bootstrap 95% CI for Delta, per-task tables'},
        'run_protocol': ['Phase 1: shell arms, N runs per agent; commit eligibility.json (report.py seal)',
                         'Phase 2: World arms, N runs per agent',
                         'Drift probe: one extra shell run per agent (informative)',
                         'Verdict: report.py verdict (refuses unless eligibility.json is committed)'],
        'canary': 'one smoke-tier task per (agent, arm) before each phase (§4.9); abort on the first quota/rate-limit evidence',
    }


def check_final(prereg_path: str, cfg: dict, probes: Probes, ledger: str = LEDGER) -> dict:
    """AC5.1: the prereg is committed and its digest equals the live FINAL config's digest."""
    require_clean_ledger(ledger)
    require_committed(prereg_path)  # MUT-FINAL-NOPREREG removes this line
    with open(prereg_path) as f:
        pre = json.load(f)
    if pre.get('kind') != PREREG_KIND:
        raise RunRefused(f'{prereg_path} is not a {PREREG_KIND}')
    live = config_digest(live_config(cfg, 'final', probes))
    if pre.get('config_digest') != live:
        raise RunRefused(f'AC5.1: prereg config_digest {pre.get("config_digest")} != live FINAL config {live}')
    return pre


# ---------------------------------------------------------------- the ledger (AC4.1)

LEDGER_FIELDS = ('date', 'iteration', 'knob', 'side', 'change', 'before', 'after', 'commit')


def read_ledger(ledger: str = LEDGER) -> list:
    if not os.path.exists(ledger):
        return []
    with open(ledger) as f:
        return [json.loads(ln) for ln in f if ln.strip()]


def append_ledger(row: dict, ledger: str = LEDGER) -> None:
    missing = [k for k in LEDGER_FIELDS if k not in row]
    if missing:
        raise RunRefused(f'ledger row missing {missing}')
    prev = read_ledger(ledger)
    if prev and row['iteration'] <= prev[-1]['iteration']:
        raise RunRefused(f'iteration {row["iteration"]} is not after the ledger head {prev[-1]["iteration"]}')
    with open(ledger, 'a') as f:
        f.write(json.dumps(row, sort_keys=True, ensure_ascii=False) + '\n')


# ---------------------------------------------------------------- outcome typing (pure)

def quota_evidence(agent: str, parsed: dict, stderr: str) -> tuple | None:
    """(category, evidence) when the provider evidences quota/rate limiting, else None."""
    texts = []
    if agent == 'claude':
        rej = [s for s in parsed.get('rate_limit_status') or [] if s and s not in ('allowed', 'allowed_warning')]
        if rej:
            return 'rate_limit', f'rate_limit_event status {rej}'
        res = parsed.get('result') or {}
        if res.get('is_error'):
            texts.append(str(res.get('result') or ''))
    else:
        texts += [str(e) for e in parsed.get('errors') or []]
    texts.append(stderr[-4000:] if stderr else '')
    for cat, rx in QUOTA_PATTERNS:
        for t in texts:
            m = rx.search(t)
            if m:
                return cat, t[max(0, m.start() - 120):m.end() + 120].strip()
    return None


def agent_failure(agent: str, rc, parsed: dict, stderr: str) -> tuple | None:
    """(category, evidence, typed) when the AGENT process failed (not the solution), else None.

    Claude: a ``result`` event with ``is_error`` / a non-success subtype, or no result at all.
    codex: an ``error`` / ``turn.failed`` event, or a non-zero exit. Unexplained = ``api_error``
    (§4.6), typed only when the evidence names a provider fault."""
    q = quota_evidence(agent, parsed, stderr)
    if q:
        return q[0], q[1], True
    if agent == 'claude':
        res = parsed.get('result')
        if res is None:
            return 'api_error', f'no result event (rc={rc}); stderr: {stderr[-300:].strip()}', False
        sub = res.get('subtype') or ''
        if sub == 'error_max_turns':
            return 'step_exhausted', f'result subtype {sub}', True
        if sub == 'error_max_budget_usd':
            return 'cost_killed', f'result subtype {sub} (--max-budget-usd)', True
        if res.get('is_error') or (sub and sub != 'success'):
            ev = f'result subtype={sub} is_error={res.get("is_error")}: {str(res.get("result"))[:300]}'
            return 'api_error', ev, bool(TYPED_API.search(ev))
        return None
    errs = [str(e) for e in parsed.get('errors') or []]
    if errs:
        ev = '; '.join(errs)[:400]
        return 'api_error', ev, bool(TYPED_API.search(ev))
    if rc not in (0, None):
        return 'api_error', f'codex rc={rc}; stderr: {stderr[-300:].strip()}', False
    return None


def row153_hits(transcript: str, daemon_log: str) -> dict:
    tl = [ln.strip()[:400] for ln in transcript.splitlines() if any(p in ln for p in ROW153_PATTERNS)]
    dl = [ln.strip()[:400] for ln in daemon_log.splitlines() if any(p in ln for p in ROW153_PATTERNS)]
    tcount = sum(transcript.count(p) for p in ROW153_PATTERNS[:1])  # 'host callback timed out' occurrences
    return {'transcript_callback_timeouts': tcount, 'transcript_lines': len(tl),
            'daemon_deadline_lines': len(dl), 'samples': (tl[:2] + dl[:3])}


def world_tool_errors(agent: str, parsed: dict) -> int:
    if agent == 'claude':
        return int(parsed.get('tool_result_errors') or 0)
    return sum(1 for c in parsed.get('mcp_tool_calls') or [] if not c.endswith(':completed'))


def decide_category(*, driver_category, agent_fail, graded: dict | None) -> tuple:
    """(error_category, cause_evidence, cause_typed). The driver's (timeout/harness_setup) wins,
    then a failed agent process, then the grade's own category."""
    if driver_category:
        return driver_category, None, False
    if graded is not None and graded.get('stdout_ok'):
        return 'none', None, False
    if agent_fail:
        return agent_fail
    if graded is None:
        return 'harness_setup', 'no grade', True
    return grade.error_category(graded), None, False


def classify_row(row: dict) -> str:
    try:
        return classify_mod.classify(row)
    except classify_mod.ClassificationError as e:
        row['classification_error'] = str(e)
        return 'UNCLASSIFIED'


# ---------------------------------------------------------------- smoke summary (Δ and O)

def summarize(rows: list, tasks: list) -> dict:
    """Per agent: pass counts, Δ = R_W − R_S and O = median_t o_t over tasks whose both arms have a
    non-fault row (N=1, so o_t = wall_W/wall_S − 1). Smoke numbers are informative tuning signal."""
    out = {}
    T = len(tasks)
    for agent in arms.AGENTS:
        a = [r for r in rows if r['agent'] == agent]
        if not a:
            continue
        by = {(r['arm'], r['task']): r for r in a}
        ps = sum(1 for t in tasks if by.get(('shell', t), {}).get('class') == classify_mod.PASS)
        pw = sum(1 for t in tasks if by.get(('world', t), {}).get('class') == classify_mod.PASS)
        o_t = {}
        for t in tasks:
            s, w = by.get(('shell', t)), by.get(('world', t))
            if not s or not w or classify_mod.HARNESS_FAULT in (s['class'], w['class']) \
                    or 'UNCLASSIFIED' in (s['class'], w['class']):
                continue
            ws = stats.DEADLINE_MS if w.get('finish_reason') == 'timeout' else w['wall_ms']
            ss = stats.DEADLINE_MS if s.get('finish_reason') == 'timeout' else s['wall_ms']
            o_t[t] = ws / ss - 1
        faults = {x: sum(1 for r in a if r['arm'] == x and r['class'] == classify_mod.HARNESS_FAULT) for x in arms.ARMS}
        out[agent] = {'T': T, 'pass_S': ps, 'pass_W': pw, 'delta': round((pw - ps) / T, 4) if T else None,
                      'O': round(float(stats.median(o_t.values())), 4) if o_t else None,
                      'o_t': {k: round(v, 3) for k, v in sorted(o_t.items())},
                      'harness_faults': faults,
                      'unclassified': sum(1 for r in a if r['class'] == 'UNCLASSIFIED')}
    return out


# ---------------------------------------------------------------- the live loop (rig only)

def _layout(root: str) -> dict:
    import canary
    L = canary.layout(root)
    L['m4wt'] = os.path.join(root, 'm4wt')
    L['m4'] = os.path.join(root, 'm4')
    return L


def _git_fixture(args, cwd):
    subprocess.run(['git', '-c', 'user.name=floor-harness', '-c', 'user.email=floor-harness@floor.invalid',
                    *args], cwd=cwd, check=True, capture_output=True)


def task_base(L: dict, task, cache: dict) -> str:
    """A commit in the fixture repo holding the placeholder + the task's input_files (§4.4 step 1)."""
    if task.id in cache:
        return cache[task.id]
    fx = L['fixture']
    root = _git(['rev-list', '--max-parents=0', 'HEAD'], fx).stdout.split()[0]
    files = task.spec.get('input_files') or {}
    tag = f'floor-task-{task.id}-{_sha(canonical(files))[:12]}'
    have = _git(['rev-parse', '--verify', '-q', f'refs/tags/{tag}^{{commit}}'], fx).stdout.strip()
    if have:
        cache[task.id] = have
        return have
    tmp = os.path.join(L['m4'], 'basebuild')
    if os.path.isdir(tmp):
        _git_fixture(['worktree', 'remove', '--force', tmp], fx)
    _git_fixture(['worktree', 'add', '-q', '--detach', tmp, root], fx)
    for name, content in files.items():
        fp = os.path.join(tmp, name)
        os.makedirs(os.path.dirname(fp), exist_ok=True)
        with open(fp, 'w') as f:
            f.write(content or '')
    _git_fixture(['add', '-A'], tmp)
    _git_fixture(['commit', '-q', '--allow-empty', '-m', f'floor base: {task.id}'], tmp)
    sha = _git(['rev-parse', 'HEAD'], tmp).stdout.strip()
    _git_fixture(['tag', tag, sha], fx)
    _git_fixture(['worktree', 'remove', '--force', tmp], fx)
    cache[task.id] = sha
    return sha


def reset_worktree(wt: str, sha: str) -> bytes | None:
    _git_fixture(['checkout', '-q', '-f', '--detach', sha], wt)
    _git_fixture(['clean', '-q', '-ffdx'], wt)
    p = os.path.join(wt, SOL)
    return open(p, 'rb').read() if os.path.exists(p) else None


def ensure_shell_worktree(L: dict, agent: str) -> str:
    wt = os.path.join(L['m4wt'], agent)
    if not os.path.isdir(wt):
        os.makedirs(L['m4wt'], exist_ok=True)
        _git_fixture(['worktree', 'add', '-q', '--detach', wt], L['fixture'])
    return wt


def _read_from(path: str, offset: int) -> str:
    try:
        with open(path, 'rb') as f:
            f.seek(offset)
            return f.read().decode('utf-8', errors='replace')
    except FileNotFoundError:
        return ''


def _size(path: str) -> int:
    try:
        return os.path.getsize(path)
    except OSError:
        return 0


class Ctx:
    def __init__(self, cfg, mode, probes, L, out_dir, transcripts_dir, *, client=None, tokens=None,
                 spend_cap=None, spent_before=0.0):
        self.cfg, self.mode, self.probes, self.L = cfg, mode, probes, L
        self.sec = mode_section(cfg, mode)
        self.out_dir, self.tdir = out_dir, transcripts_dir
        self.client, self.tokens = client, tokens or {}
        self.template = probes.template(cfg['corpus_commit'])
        self.teaching = probes.teaching(tool_bin(cfg), cfg['knobs']['teaching'] == 'compact')
        self.base_cache = {}
        self.spend_cap, self.spent = spend_cap, spent_before
        self.spend_this = 0.0


def run_cell(ctx: Ctx, agent: str, arm: str, task, *, run_k: int, rows_path: str, label: str) -> dict:
    """§4.4 for one task-run: reset, (pre-flight), spawn, grade, classify, append one row."""
    import run_task
    import world
    cfg, L = ctx.cfg, ctx.L
    k = cfg['knobs']
    if agent == 'claude' and ctx.spend_cap is not None:
        if ctx.spent + ctx.spend_this + ctx.sec['claude_max_budget_usd'] > ctx.spend_cap:
            raise RunRefused(f'spend cap: ${ctx.spent + ctx.spend_this:.2f} spent + '
                             f'${ctx.sec["claude_max_budget_usd"]:.2f} could exceed ${ctx.spend_cap:.2f}')
    ep = EPISODES[agent]
    wt = ensure_shell_worktree(L, agent) if arm == 'shell' else os.path.join(L['ws'], ep)
    void = os.path.join(L['void'], ep) if arm == 'world' else None
    base = reset_worktree(wt, task_base(L, task, ctx.base_cache))
    spec = task.spec
    sol_path = os.path.join(wt, SOL) if arm == 'shell' else prompt.WORLD_SOLUTION_PATH
    mock = grade.HTTPMock(port=MOCK_PORT) if spec.get('net_allow_localhost') else None
    tr = os.path.join(ctx.tdir, f'{label}-{agent}-{arm}-{task.id}.jsonl')
    os.makedirs(ctx.tdir, exist_ok=True)
    log_off = _size(L['serve_log'])
    t_start = time.time()
    if mock:
        mock.__enter__()
    try:
        task_text = prompt.render_task(spec, arm, template=ctx.template, solution_path=sol_path,
                                       deadline_s=cfg['deadline_s'],
                                       mock_url=mock.url if mock else None,
                                       world_wording_id=k['world_tools_wording'])
        msg = prompt.render_message(ctx.teaching, task_text)
        drv = run_task.run_agent(
            agent, arm, model=ctx.sec['models'][agent], prompt_text=msg, worktree=wt,
            private_dir=L['private'], transcript_path=tr, void_dir=void,
            episode=ep if arm == 'world' else None, token=ctx.tokens.get(agent) if arm == 'world' else None,
            tool_bin=tool_bin(cfg), deadline_s=cfg['deadline_s'], client=ctx.client,
            claude_isolation=k['claude_isolation'],
            max_budget_usd=ctx.sec['claude_max_budget_usd'] if agent == 'claude' else None)
        graded = None
        if drv.get('error_category') != 'harness_setup':
            try:
                graded = grade.grade(world.solution_for_grading(wt), spec, tool_bin(cfg))
            except grade.GradeSetupError as e:
                drv.update(error_category='harness_setup', cause_evidence=f'grader: {e}')
    finally:
        if mock:
            mock.__exit__(None, None, None)
    daemon_log = _read_from(L['serve_log'], log_off) if arm == 'world' else ''
    parsed = drv.get('transcript') or {}
    stderr = _read_from(tr + '.stderr', 0)
    transcript = _read_from(tr, 0)
    afail = None if drv.get('error_category') else agent_failure(agent, drv.get('rc'), parsed, stderr)
    cat, ev, typed = decide_category(driver_category=drv.get('error_category'), agent_fail=afail, graded=graded)
    if drv.get('cause_evidence'):
        ev = drv['cause_evidence']
    if cat == 'timeout':
        ev = ev or f'killed at the {cfg["deadline_s"]} s deadline'
    res = parsed.get('result') or {}
    row = {
        'id': task.id, 'task': task.id, 'lang': 'ailang', 'model': ctx.sec['models'][agent],
        'executor': agent, 'agent': agent, 'arm': arm, 'run': run_k, 'trial': run_k, 'label': label,
        'compile_ok': bool(graded and graded['compile_ok']), 'runtime_ok': bool(graded and graded['runtime_ok']),
        'stdout_ok': bool(graded and graded['stdout_ok']), 'error_category': cat, 'cause_evidence': ev,
        'cause_typed': typed, 'finish_reason': drv.get('finish_reason'), 'rc': drv.get('rc'),
        'wall_ms': drv.get('wall_ms'), 'duration_ms': drv.get('wall_ms'),
        'cost_usd': res.get('total_cost_usd') if agent == 'claude' else None,
        'cost_provenance': 'claude result.total_cost_usd' if agent == 'claude' else 'codex: tokens only',
        'input_tokens': ((res.get('usage') or {}).get('input_tokens') if agent == 'claude'
                         else (parsed.get('usage') or {}).get('input_tokens')),
        'output_tokens': ((res.get('usage') or {}).get('output_tokens') if agent == 'claude'
                          else (parsed.get('usage') or {}).get('output_tokens')),
        'num_turns': res.get('num_turns'), 'tool_items': drv.get('tool_items'),
        'world_tool_errors': world_tool_errors(agent, parsed) if arm == 'world' else None,
        'argv': drv.get('argv'), 'world': drv.get('world'), 'preflight': drv.get('preflight'),
        'void_empty': drv.get('void_empty'), 'mcp_config_deleted': drv.get('mcp_config_deleted'),
        'solution_sha256': (graded or {}).get('solution_sha256'), 'solution_state': (graded or {}).get('solution_state'),
        'grade_stdout_head': ((graded or {}).get('stdout') or '')[:400],
        'grade_stderr_head': ((graded or {}).get('stderr') or '')[:400],
        'digests': {'prompt': _sha(msg.encode()), 'teaching': _sha(ctx.teaching.encode()),
                    'yaml': task.sha256, 'tool_bin': ctx.probes.tool_sha256(tool_bin(cfg)), 'grader': grade.GRADER_VERSION},
        'transcript_sha256': _sha(transcript.encode('utf-8', errors='surrogateescape')),
        'transcript_bytes': len(transcript.encode('utf-8', errors='surrogateescape')),
        'transcript_file': os.path.basename(tr), 'started_at': dt.datetime.fromtimestamp(t_start).isoformat(timespec='seconds'),
    }
    # attestation of the effective tool set, every task-run (AC2.2/AC2.3/AC2.6)
    att = []
    if drv.get('error_category') != 'harness_setup':
        if agent == 'claude':
            att += attest.claude_init(parsed.get('init'), arm) + attest.claude_used_only(parsed, arm)
            init = parsed.get('init') or {}
            row['init_plugins'] = init.get('plugins')
            row['hooks'] = len(parsed.get('hooks') or [])
            if k['claude_isolation'] and (init.get('plugins') or parsed.get('hooks')):
                att.append(f'isolation: plugins {init.get("plugins")} hooks {len(parsed.get("hooks") or [])}')
        elif arm == 'world':
            att += attest.codex_world(parsed, drv.get('void_listing') or [])
        else:
            att += attest.codex_no_web(parsed)
    row['attestation'] = att
    if arm == 'world' and drv.get('world') and drv['world'].get('log_to') is not None:
        try:
            prov = world.check_task_run(ctx.client, episode=ep, log_from=drv['world']['log_from'],
                                        log_to=drv['world']['log_to'], worktree=wt, void_dir=void, base=base)
        except (world.WorldReadError, OSError, ValueError) as e:
            prov = {'error': str(e), 'native_write_detected': None}
        row['provenance'] = prov
    if arm == 'world':
        row['row153'] = row153_hits(transcript, daemon_log)
    row['class'] = classify_row(row)
    if agent == 'claude' and row['cost_usd']:
        ctx.spend_this += float(row['cost_usd'])
    os.makedirs(os.path.dirname(rows_path), exist_ok=True)
    with open(rows_path, 'a') as f:
        f.write(json.dumps(row, sort_keys=True, ensure_ascii=False, default=str) + '\n')
    q = quota_evidence(agent, parsed, stderr)
    if q:
        raise QuotaAbort(f'{agent}/{arm}/{task.id}: {q[0]}: {q[1]}')
    return row


def read_token(path: str) -> str:
    st = os.stat(path)
    if st.st_mode & 0o077:
        raise RunRefused(f'{path} is not mode 0600')
    with open(path) as f:
        return f.read().strip()


def spent_so_far(ledger: str = LEDGER) -> float:
    return round(sum(float(r.get('spend_usd_claude') or 0) for r in read_ledger(ledger)), 4)


def cmd_iterate(a) -> int:
    import canary
    import tokens as tokens_mod
    require_clean_ledger(LEDGER)  # AC4.2
    st = dirty(HERE)
    if st and not a.allow_dirty_harness:
        raise RunRefused(f'scripts/floor is dirty ({st[:200]!r}); commit the config/harness under test first')
    cfg = load_config(a.config)
    probes = Probes()
    mode = 'smoke'
    tasks = mode_tasks(cfg, mode, probes)
    can = canary_task(cfg, mode, probes)
    agents = a.agents.split(',')
    arms_run = a.arms.split(',')
    L = _layout(a.root)
    it_dir = os.path.join(a.evidence, f'iter-{a.iteration}')
    if os.path.exists(it_dir):
        raise RunRefused(f'{it_dir} exists: an iteration is run once')
    tdir = os.path.join(L['m4'], 'transcripts', f'iter-{a.iteration}')
    os.makedirs(it_dir)
    live = live_config(cfg, mode, probes)
    digests = {f'{ag}/{x}': arm_digest(cfg, mode, ag, x, probes) for ag in arms.AGENTS for x in arms.ARMS}
    reused = {}
    if a.reuse_shell_from:
        src = os.path.join(a.evidence, a.reuse_shell_from)
        with open(os.path.join(src, 'iteration.json')) as f:
            prev = json.load(f)
        for ag in agents:
            key = f'{ag}/shell'
            if prev['arm_digests'].get(key) != digests[key]:
                raise RunRefused(f'cannot reuse {key} rows from {a.reuse_shell_from}: arm digest changed')
            s = os.path.join(src, 'runs', ag, 'shell', 'r1.jsonl')
            d = os.path.join(it_dir, 'runs', ag, 'shell', 'r1.jsonl')
            os.makedirs(os.path.dirname(d), exist_ok=True)
            shutil.copyfile(s, d)
            reused[key] = a.reuse_shell_from
    toks, client, daemon = {}, None, None
    if 'world' in arms_run:
        import world
        sess = {ag: os.path.join(L['private'], f'{EPISODES[ag]}.session') for ag in agents}
        toks = {ag: read_token(p) for ag, p in sess.items()}
        daemon = canary.start_daemon(L, a.addr)
        client = world.Client(a.addr)
    spent_before = spent_so_far(LEDGER) + a.prior_spend
    ctx = Ctx(cfg, mode, probes, L, it_dir, tdir, client=client, tokens=toks, spend_cap=a.spend_cap,
              spent_before=spent_before)
    t0 = time.time()
    status, abort = 'complete', None
    canaries = []
    try:
        if client is not None:
            canary.commit_genesis_if_needed(L, client, toks[agents[0]])
        for ag in agents:
            for x in arms_run:
                if f'{ag}/{x}' in reused:
                    continue
                if time.time() - t0 > a.max_wall_s:
                    raise RunRefused(f'wall-clock ceiling {a.max_wall_s} s reached')
                c = run_cell(ctx, ag, x, can, run_k=1, rows_path=os.path.join(it_dir, 'canary.jsonl'), label='canary')
                canaries.append({'agent': ag, 'arm': x, 'task': can.id, 'class': c['class'],
                                 'attestation': c['attestation'], 'wall_ms': c['wall_ms']})
                print(f'canary {ag}/{x}: {c["class"]} att={c["attestation"]} {c["wall_ms"]} ms', flush=True)
                if c['attestation'] or c['error_category'] == 'harness_setup':
                    raise RunRefused(f'canary {ag}/{x} failed attestation/setup: {c["attestation"]} {c["cause_evidence"]}')
                for t in tasks:
                    if time.time() - t0 > a.max_wall_s:
                        raise RunRefused(f'wall-clock ceiling {a.max_wall_s} s reached')
                    r = run_cell(ctx, ag, x, t, run_k=1,
                                 rows_path=os.path.join(it_dir, 'runs', ag, x, 'r1.jsonl'), label=f'iter{a.iteration}')
                    print(f'{ag}/{x}/{t.id}: {r["class"]} {r["error_category"]} {r["wall_ms"]} ms '
                          f'${r["cost_usd"]} att={r["attestation"]} '
                          f'prov={(r.get("provenance") or {}).get("solution_provenance")} '
                          f'r153={(r.get("row153") or {}).get("transcript_callback_timeouts")}', flush=True)
    except (QuotaAbort, RunRefused) as e:
        status, abort = 'aborted', str(e)
        print(f'ABORT: {e}', flush=True)
    finally:
        if daemon is not None:
            canary.stop_daemon(daemon)
    wall_s = round(time.time() - t0, 1)
    rows = []
    for ag in agents:
        for x in arm_list():
            p = os.path.join(it_dir, 'runs', ag, x, 'r1.jsonl')
            if os.path.exists(p):
                rows += [json.loads(ln) for ln in open(p) if ln.strip()]
    summ = summarize(rows, [t.id for t in tasks])
    r153 = {}
    for r in rows:
        if r['arm'] == 'world':
            h = r.get('row153') or {}
            d = r153.setdefault(r['agent'], {'task_runs': 0, 'runs_with_callback_timeout': 0,
                                            'callback_timeouts': 0, 'daemon_deadline_lines': 0})
            d['task_runs'] += 1
            d['callback_timeouts'] += h.get('transcript_callback_timeouts', 0)
            d['runs_with_callback_timeout'] += 1 if (h.get('transcript_callback_timeouts') or h.get('daemon_deadline_lines')) else 0
            d['daemon_deadline_lines'] += h.get('daemon_deadline_lines', 0)
    native = [f'{r["agent"]}/{r["task"]}' for r in rows if r['arm'] == 'world' and (r.get('provenance') or {}).get('native_write_detected')]
    world_writes = {f'{r["agent"]}/{r["task"]}': (r.get('provenance') or {}).get('world_writes') for r in rows if r['arm'] == 'world'}
    it = {'iteration': a.iteration, 'knob': a.knob, 'side': a.side, 'change': a.change, 'status': status,
          'abort': abort, 'commit': head_commit(), 'config_digest_smoke': config_digest(live),
          'arm_digests': digests, 'shell_rows_reused_from': reused, 'tasks': [t.id for t in tasks],
          'canary': canaries, 'summary': summ, 'row153': r153, 'native_write_detected': native,
          'world_writes_per_run': world_writes, 'spend_usd_claude': round(ctx.spend_this, 4),
          'codex_tokens': sum((r.get('input_tokens') or 0) + (r.get('output_tokens') or 0)
                              for r in rows + read_jsonl(os.path.join(it_dir, 'canary.jsonl'))
                              if r['agent'] == 'codex' and f'{r["agent"]}/{r["arm"]}' not in reused),
          'wall_s': wall_s, 'transcripts': os.path.relpath(tdir, os.path.expanduser('~'))}
    with open(os.path.join(it_dir, 'iteration.json'), 'w') as f:
        json.dump(it, f, indent=2, sort_keys=True)
        f.write('\n')
    tg = tokens_mod.scan_dir(it_dir, tokens_mod.load_tokens(
        [os.path.join(L['private'], f'{EPISODES[ag]}.session') for ag in arms.AGENTS
         if os.path.exists(os.path.join(L['private'], f'{EPISODES[ag]}.session'))]))
    if tg:
        print(f'AC2.7 FAIL: token findings in {it_dir}: {tg}', flush=True)
        return 3
    prev = read_ledger(LEDGER)
    before = prev[-1]['after'] if prev else None
    after = {ag: {'delta': s['delta'], 'O': s['O'], 'pass_S': s['pass_S'], 'pass_W': s['pass_W'], 'T': s['T']}
             for ag, s in summ.items()}
    append_ledger({'date': dt.date.today().isoformat(), 'iteration': a.iteration, 'knob': a.knob,
                   'side': a.side, 'change': a.change, 'before': before, 'after': after,
                   'commit': it['commit'], 'config_digest_smoke': it['config_digest_smoke'],
                   'status': status, 'evidence': os.path.relpath(it_dir, REPO),
                   'spend_usd_claude': it['spend_usd_claude'], 'codex_tokens': it['codex_tokens'],
                   'wall_s': wall_s, 'row153': r153, 'shell_rows_reused_from': reused}, LEDGER)
    print(json.dumps({'summary': summ, 'row153': r153, 'spend': it['spend_usd_claude'], 'wall_s': wall_s,
                      'status': status}, indent=1))
    return 0 if status == 'complete' else 4


def arm_list():
    return arms.ARMS


def read_jsonl(path: str) -> list:
    if not os.path.exists(path):
        return []
    with open(path) as f:
        return [json.loads(ln) for ln in f if ln.strip()]


def cmd_digest(a) -> int:
    cfg = load_config(a.config)
    live = live_config(cfg, a.mode, Probes())
    print(json.dumps({'config_digest': config_digest(live), 'live_config': live}, indent=1, sort_keys=True))
    return 0


def cmd_prereg_draft(a) -> int:
    cfg = load_config(a.config)
    pre = prereg_draft(cfg, Probes())
    with open(a.out, 'w') as f:
        json.dump(pre, f, indent=2, sort_keys=True)
        f.write('\n')
    print(f'{a.out}: config_digest {pre["config_digest"]}')
    return 0


def cmd_final(a) -> int:
    cfg = load_config(a.config)
    check_final(a.prereg, cfg, Probes())  # AC5.1 — before any spawn
    raise RunRefused('FINAL phases are run in the attended M5 session (design §6 M5); the prereg '
                     'check passed, and this M4 build stops here by design')


def main(argv=None) -> int:
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    sub = ap.add_subparsers(dest='cmd', required=True)
    p = sub.add_parser('iterate')
    p.add_argument('--iteration', type=int, required=True)
    p.add_argument('--knob', required=True)
    p.add_argument('--side', required=True, choices=['baseline', 'harness', 'world'])
    p.add_argument('--change', required=True)
    p.add_argument('--agents', default='claude,codex')
    p.add_argument('--arms', default='shell,world')
    p.add_argument('--reuse-shell-from', default=None)
    p.add_argument('--root', default=DEFAULT_ROOT)
    p.add_argument('--addr', default=arms.DEFAULT_ADDR)
    p.add_argument('--evidence', default=EVIDENCE)
    p.add_argument('--spend-cap', type=float, default=15.0)
    p.add_argument('--prior-spend', type=float, default=0.0, help='Claude spend outside the ledger (probes)')
    p.add_argument('--max-wall-s', type=float, default=3 * 3600)
    p.add_argument('--allow-dirty-harness', action='store_true')
    p.set_defaults(fn=cmd_iterate)
    p = sub.add_parser('digest')
    p.add_argument('--mode', default='smoke', choices=MODES)
    p.set_defaults(fn=cmd_digest)
    p = sub.add_parser('prereg-draft')
    p.add_argument('--out', default=os.path.join(EVIDENCE, 'preregistration.draft.json'))
    p.set_defaults(fn=cmd_prereg_draft)
    p = sub.add_parser('final')
    p.add_argument('--prereg', required=True)
    p.add_argument('--phase', required=True, choices=['shell', 'world'])
    p.add_argument('--run', type=int, required=True)
    p.set_defaults(fn=cmd_final)
    for sp in sub.choices.values():
        sp.add_argument('--config', default=CONFIG)
    a = ap.parse_args(argv)
    if hasattr(a, 'root'):
        a.root = os.path.abspath(os.path.expanduser(a.root))
    try:
        return a.fn(a)
    except RunRefused as e:
        print(f'REFUSED: {e}', file=sys.stderr)
        return 2


if __name__ == '__main__':
    raise SystemExit(main())
