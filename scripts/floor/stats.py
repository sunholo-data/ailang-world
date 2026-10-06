#!/usr/bin/env python3
"""Row 93 floor harness — exact statistics and the verdict (design §4.7; AC3.2, AC3.3).

All gate arithmetic is EXACT (``fractions.Fraction``): the thresholds sit on lattice points
(Δ = −2.00 pp at T=50, range = 5.00 pp at T=20), where a float comparison can flip.

Notation (§4.7): T tasks, N runs per arm; p = 1 for a PASS task-run; r_{a,x,k} = (1/T)·Σ_t p;
R_{a,x} = (1/(T·N))·Σ_k Σ_t p.

  * Eligible(a) ⇔ N ≥ 3 ∧ zero HARNESS_FAULT rows in S ∧ max_k r_{a,S,k} − min_k r_{a,S,k} ≤ 0.05
  * Pass-rate test: Δ_a = R_{a,W} − R_{a,S} ≥ −0.02 (the POINT ESTIMATE; D-WORLD-55)
  * Overhead: o_t = median_k wall_W(t,k) / median_k wall_S(t,k) − 1; O_a = median_t o_t ≤ 0.25.
    Every non-harness-fault task-run counts; a timeout counts as the 600 s deadline.
  * Informative only, never the gate: ratio of pooled medians, ratio of total times, a paired
    (by task) bootstrap 95% CI for Δ, and per-task tables.

Verdict, in order: no agent eligible → PAUSED; any eligible agent fails either test → FAILS;
both eligible and both pass → HOLDS; otherwise → INCOMPLETE.

Rows are dicts with at least ``task, run, arm, class`` (PASS/FAIL/HARNESS_FAULT, from
classify.py), ``wall_ms`` and ``finish_reason``; a re-run task-run (World quota/rate-limit, or the
capacity rule in either arm, §4.6) carries ``attempt`` — the highest attempt is the task-run's outcome.
"""
from __future__ import annotations

import random
from fractions import Fraction

from classify import FAIL, HARNESS_FAULT, PASS

MIN_RUNS = 3
RANGE_MAX = Fraction(5, 100)
DELTA_MIN = Fraction(-2, 100)
OVERHEAD_MAX = Fraction(25, 100)
DEADLINE_MS = 600_000

HOLDS, FAILS, INCOMPLETE, PAUSED = 'HOLDS', 'FAILS', 'INCOMPLETE', 'PAUSED'


class StatsError(Exception):
    """The rows cannot support the computation (incomplete, inconsistent, or blocked)."""


def median(values) -> Fraction:
    """The ordinary median: middle value (odd), mean of the two middle values (even)."""
    v = sorted(Fraction(x) for x in values)
    if not v:
        raise StatsError('median of no values')
    n = len(v)
    if n % 2:
        return v[n // 2]
    return (v[n // 2 - 1] + v[n // 2]) / 2


def final_rows(rows) -> list:
    """One row per (task, run): the highest ``attempt``."""
    best = {}
    for r in rows:
        key = (r['task'], r['run'])
        if key not in best or r.get('attempt', 1) > best[key].get('attempt', 1):
            best[key] = r
    return list(best.values())


def _grid(rows, tasks, arm: str) -> dict:
    """{run: {task: row}} with completeness checked: every run holds every task exactly once."""
    tasks = list(tasks)
    grid: dict = {}
    for r in rows:
        if r['arm'] != arm:
            raise StatsError(f'{arm} computation given a {r["arm"]} row')
        if r['class'] not in (PASS, FAIL, HARNESS_FAULT):
            raise StatsError(f'unclassified row {r}')
        grid.setdefault(r['run'], {})
        if r['task'] in grid[r['run']]:
            raise StatsError(f'duplicate {arm} task-run {r["task"]} r{r["run"]}')
        grid[r['run']][r['task']] = r
    for k, g in grid.items():
        if set(g) != set(tasks):
            missing = sorted(set(tasks) - set(g))
            extra = sorted(set(g) - set(tasks))
            raise StatsError(f'{arm} run {k} incomplete: missing {missing} extra {extra}')
    return grid


def run_rates(grid, T: int) -> dict:
    return {k: Fraction(sum(r['class'] == PASS for r in g.values()), T) for k, g in grid.items()}


def eligibility(shell_rows, tasks) -> dict:
    """Eligible(a) from ONE shell-phase attempt (all its rows, faults included)."""
    tasks = list(tasks)
    T = len(tasks)
    shell_rows = final_rows(shell_rows)  # §4.6 capacity rule: the last attempt is the outcome
    grid = _grid(shell_rows, tasks, 'shell')
    N = len(grid)
    faults = sorted((r['task'], r['run']) for g in grid.values() for r in g.values() if r['class'] == HARNESS_FAULT)
    rates = run_rates(grid, T)
    rng = (max(rates.values()) - min(rates.values())) if rates else None
    reasons = []
    if N < MIN_RUNS:
        reasons.append(f'N={N} < {MIN_RUNS}')
    if faults:
        reasons.append(f'{len(faults)} harness fault(s) in the shell arm')
    if rng is not None and not (rng <= RANGE_MAX):
        reasons.append(f'run range {rng} > {RANGE_MAX}')
    pooled = Fraction(sum(r['class'] == PASS for g in grid.values() for r in g.values()), T * N) if N else None
    return {'eligible': not reasons, 'N': N, 'T': T, 'reasons': reasons,
            'harness_faults': [list(f) for f in faults],
            'run_pass_counts': {str(k): int(v * T) for k, v in sorted(rates.items())},
            'run_rates': {str(k): str(v) for k, v in sorted(rates.items())},
            'range': str(rng) if rng is not None else None, 'R_S': str(pooled) if pooled is not None else None}


def _wall(r) -> int:
    if r.get('finish_reason') == 'timeout':
        return DEADLINE_MS
    w = r.get('wall_ms')
    if w is None or w < 0:
        raise StatsError(f'row without wall_ms: {r["task"]} r{r["run"]}')
    return int(w)


def overhead(shell_rows, world_rows, tasks) -> dict:
    """o_t per task and O_a = median_t o_t (harness-fault task-runs excluded; timeout = 600 s)."""
    per_task = {}
    for t in tasks:
        ws = [_wall(r) for r in world_rows if r['task'] == t and r['class'] != HARNESS_FAULT]
        ss = [_wall(r) for r in shell_rows if r['task'] == t and r['class'] != HARNESS_FAULT]
        if not ws or not ss:
            raise StatsError(f'task {t}: no non-fault task-runs in one arm (W={len(ws)}, S={len(ss)})')
        ms = median(ss)
        if ms <= 0:
            raise StatsError(f'task {t}: non-positive shell median wall time')
        per_task[t] = median(ws) / ms - 1
    return {'o_t': per_task, 'O': median(per_task.values())}


def pass_rate_test(shell_rows, world_rows, tasks) -> dict:
    tasks = list(tasks)
    T = len(tasks)
    gs = _grid(shell_rows, tasks, 'shell')
    gw = _grid(world_rows, tasks, 'world')
    if len(gs) != len(gw):
        raise StatsError(f'N differs between arms: S={len(gs)} W={len(gw)}')
    for g in gw.values():
        for r in g.values():
            if r['class'] == HARNESS_FAULT:
                raise StatsError(f'World task-run {r["task"]} r{r["run"]} ends in a harness fault: '
                                 'the row blocks on quota (§4.6), no verdict')
    N = len(gs)
    sum_s = sum(r['class'] == PASS for g in gs.values() for r in g.values())
    sum_w = sum(r['class'] == PASS for g in gw.values() for r in g.values())
    R_S = Fraction(sum_s, T * N)
    R_W = Fraction(sum_w, T * N)
    delta = R_W - R_S
    return {'sum_S': sum_s, 'sum_W': sum_w, 'R_S': R_S, 'R_W': R_W, 'delta': delta,
            'passes': delta >= DELTA_MIN}


def agent_tests(shell_rows, world_rows, tasks) -> dict:
    shell_rows = final_rows(shell_rows)
    world_rows = final_rows(world_rows)
    pr = pass_rate_test(shell_rows, world_rows, tasks)
    ov = overhead(shell_rows, world_rows, tasks)
    return {'pass_rate': pr, 'overhead': ov,
            'overhead_passes': ov['O'] <= OVERHEAD_MAX,
            'passes': bool(pr['passes'] and ov['O'] <= OVERHEAD_MAX)}


def verdict(eligible: dict, tests: dict) -> str:
    """eligible: {agent: bool}; tests: {agent: {'passes': bool}} for the eligible agents."""
    if not eligible:
        raise StatsError('no agents')
    elig = [a for a, e in eligible.items() if e]
    if not elig:
        return PAUSED
    for a in elig:
        if a not in tests:
            raise StatsError(f'eligible agent {a} has no test result')
    if any(not tests[a]['passes'] for a in elig):
        return FAILS
    if len(elig) == len(eligible):
        return HOLDS
    return INCOMPLETE


# ---------------------------------------------------------------- informative only

def informative(shell_rows, world_rows, tasks, *, seed: int = 93, B: int = 10000) -> dict:
    """Never the gate: pooled-median ratio, total-time ratio, paired bootstrap CI for Δ."""
    tasks = list(tasks)
    shell_rows = final_rows(shell_rows)
    world_rows = final_rows(world_rows)
    s_ok = [r for r in shell_rows if r['class'] != HARNESS_FAULT]
    w_ok = [r for r in world_rows if r['class'] != HARNESS_FAULT]
    pooled_ratio = median([_wall(r) for r in w_ok]) / median([_wall(r) for r in s_ok])
    total_ratio = Fraction(sum(_wall(r) for r in w_ok), sum(_wall(r) for r in s_ok))
    per_task = {}
    for t in tasks:
        ps = [r['class'] == PASS for r in shell_rows if r['task'] == t]
        pw = [r['class'] == PASS for r in world_rows if r['task'] == t]
        per_task[t] = {'pass_S': sum(ps), 'n_S': len(ps), 'pass_W': sum(pw), 'n_W': len(pw)}
    diffs = [Fraction(v['pass_W'], v['n_W']) - Fraction(v['pass_S'], v['n_S']) for v in per_task.values()]
    rnd = random.Random(seed)
    T = len(diffs)
    boots = sorted(float(sum(rnd.choice(diffs) for _ in range(T)) / T) for _ in range(B))
    lo, hi = boots[int(0.025 * B)], boots[int(0.975 * B) - 1]
    return {'note': 'INFORMATIVE ONLY — never the gate (D-WORLD-55: point estimate)',
            'pooled_median_ratio': float(pooled_ratio), 'total_time_ratio': float(total_ratio),
            'delta_bootstrap_ci95': [lo, hi], 'bootstrap': {'B': B, 'seed': seed, 'paired_by': 'task'},
            'per_task': per_task}
