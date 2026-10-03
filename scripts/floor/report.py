#!/usr/bin/env python3
"""Row 93 floor harness — sealed eligibility, verdict and report (design §4.7, §4.8, §4.11; AC3.4).

The order is the gate's integrity (V36: "run the eligibility precondition first and report it
separately … only then run the paired arms"):

  1. ``seal`` — eligibility is computed from the SHELL rows alone and written to
     ``eligibility.json``; the harness/attended operator then commits it.
  2. ``verdict`` — REFUSES to compute any World statistic unless ``eligibility.json`` exists, is
     tracked (``git ls-files --error-unmatch``) and its working copy equals the committed blob
     (``git diff --quiet HEAD --``). Only then are World rows read. The sealed eligibility must
     also equal a recomputation from the shell rows (a sealed file cannot be swapped).

Evidence layout (§4.8): ``<dir>/runs/<agent>/<arm>/r<k>.jsonl``, ``eligibility.json``,
``verdict.json``, ``summary.md``.
"""
from __future__ import annotations

import argparse
import glob
import json
import os
import subprocess
import sys
from fractions import Fraction

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import classify as classify_mod  # noqa: E402
import stats  # noqa: E402

ELIGIBILITY = 'eligibility.json'


class SealError(Exception):
    """World statistics requested before eligibility was sealed (committed)."""


def load_rows(evidence_dir: str, agent: str, arm: str) -> list:
    rows = []
    paths = sorted(glob.glob(os.path.join(evidence_dir, 'runs', agent, arm, 'r*.jsonl')))
    if not paths:
        raise stats.StatsError(f'no {arm} rows for {agent} under {evidence_dir}')
    for p in paths:
        with open(p) as f:
            for line in f:
                if line.strip():
                    r = json.loads(line)
                    r.setdefault('arm', arm)
                    if r['arm'] != arm:
                        raise stats.StatsError(f'{p}: row arm {r["arm"]} in the {arm} directory')
                    r['class'] = classify_mod.classify(r)
                    rows.append(r)
    return rows


def _jsonable(x):
    if isinstance(x, Fraction):
        return {'exact': str(x), 'value': float(x)}
    if isinstance(x, dict):
        return {k: _jsonable(v) for k, v in x.items()}
    if isinstance(x, (list, tuple)):
        return [_jsonable(v) for v in x]
    return x


def compute_eligibility(evidence_dir: str, agents, tasks) -> dict:
    return {'tasks': list(tasks),
            'agents': {a: stats.eligibility(load_rows(evidence_dir, a, 'shell'), tasks) for a in agents}}


def seal(evidence_dir: str, agents, tasks) -> str:
    path = os.path.join(evidence_dir, ELIGIBILITY)
    if os.path.exists(path):
        raise SealError(f'{path} already exists; eligibility is sealed once per attempt')
    e = compute_eligibility(evidence_dir, agents, tasks)
    with open(path, 'w') as f:
        json.dump(e, f, indent=2, sort_keys=True)
        f.write('\n')
    return path


def require_sealed(evidence_dir: str) -> dict:
    """AC3.4: eligibility.json exists, is git-tracked and is committed as it stands."""
    path = os.path.abspath(os.path.join(evidence_dir, ELIGIBILITY))
    if not os.path.isfile(path):
        raise SealError(f'{path} does not exist: seal eligibility before any World statistic')
    d, base = os.path.dirname(path), os.path.basename(path)
    tracked = subprocess.run(['git', '-C', d, 'ls-files', '--error-unmatch', '--', base], capture_output=True)
    if tracked.returncode != 0:
        raise SealError(f'{path} is not committed (git ls-files --error-unmatch rc={tracked.returncode})')
    clean = subprocess.run(['git', '-C', d, 'diff', '--quiet', 'HEAD', '--', base], capture_output=True)
    if clean.returncode != 0:
        raise SealError(f'{path} differs from its committed blob (git diff --quiet HEAD rc={clean.returncode})')
    with open(path) as f:
        return json.load(f)


def compute_verdict(evidence_dir: str, agents, tasks) -> dict:
    sealed = require_sealed(evidence_dir)  # BEFORE any World row is read (MUT-SEAL-ORDER)
    tasks = list(tasks)
    if sealed.get('tasks') != tasks or set(sealed.get('agents', {})) != set(agents):
        raise SealError('sealed eligibility covers different tasks/agents')
    recomputed = compute_eligibility(evidence_dir, agents, tasks)
    if recomputed != sealed:
        raise SealError('sealed eligibility differs from a recomputation over the shell rows')
    eligible = {a: bool(sealed['agents'][a]['eligible']) for a in agents}
    tests, info = {}, {}
    for a in agents:
        if not eligible[a]:
            continue
        s_rows = load_rows(evidence_dir, a, 'shell')
        w_rows = load_rows(evidence_dir, a, 'world')
        tests[a] = stats.agent_tests(s_rows, w_rows, tasks)
        info[a] = stats.informative(s_rows, w_rows, tasks)
    v = stats.verdict(eligible, tests)
    return {'verdict': v, 'eligible': eligible, 'tests': _jsonable(tests), 'informative': info,
            'thresholds': {'range_max': str(stats.RANGE_MAX), 'delta_min': str(stats.DELTA_MIN),
                           'overhead_max': str(stats.OVERHEAD_MAX), 'min_runs': stats.MIN_RUNS,
                           'deadline_ms': stats.DEADLINE_MS}}


def summary_md(result: dict) -> str:
    lines = [f'# Floor verdict: {result["verdict"]}', '',
             '| agent | eligible | Δ (point) | Δ test | O (median o_t) | O test |', '|---|---|---|---|---|---|']
    for a, e in result['eligible'].items():
        t = result['tests'].get(a)
        if t:
            pr, ov = t['pass_rate'], t['overhead']
            lines.append(f'| {a} | yes | {pr["delta"]["exact"]} ({pr["delta"]["value"]:+.4f}) | '
                         f'{"pass" if pr["passes"] else "FAIL"} | {ov["O"]["exact"]} ({ov["O"]["value"]:+.4f}) | '
                         f'{"pass" if t["overhead_passes"] else "FAIL"} |')
        else:
            lines.append(f'| {a} | no | — | — | — | — |')
    lines += ['', 'Informative only (never the gate): bootstrap 95% CI for Δ, pooled-median and '
              'total-time ratios — see verdict.json `informative`.', '']
    return '\n'.join(lines)


def main(argv=None) -> int:
    ap = argparse.ArgumentParser(description='Seal eligibility, then compute the verdict.')
    ap.add_argument('action', choices=['seal', 'verdict'])
    ap.add_argument('--evidence', required=True)
    ap.add_argument('--agent', action='append', required=True)
    ap.add_argument('--tasks-from', required=True, help='preregistration.json (corpus.yaml_sha256 keys)')
    a = ap.parse_args(argv)
    with open(a.tasks_from) as f:
        tasks = sorted(json.load(f)['corpus']['yaml_sha256'])
    if a.action == 'seal':
        print(seal(a.evidence, a.agent, tasks))
        return 0
    r = compute_verdict(a.evidence, a.agent, tasks)
    with open(os.path.join(a.evidence, 'verdict.json'), 'w') as f:
        json.dump(r, f, indent=2, sort_keys=True)
        f.write('\n')
    with open(os.path.join(a.evidence, 'summary.md'), 'w') as f:
        f.write(summary_md(r))
    print(r['verdict'])
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
