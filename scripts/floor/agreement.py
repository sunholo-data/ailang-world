#!/usr/bin/env python3
"""Row 93 floor harness — AC1.2 grader agreement against banked agent rows (V18).

Re-grades banked ailang agent rows with the in-repo grader port and compares ``stdout_ok``
with the banked verdict. A row is eligible only when:
  * it is an ``ailang`` row of a pinned gate-tier task, under an ``agent/`` results dir;
  * its task's YAML is UNCHANGED since the row was banked: the row's timestamp is after the
    last commit touching ``benchmarks/<id>.yml`` at the pinned commit, AND the row's own
    ``expected_stdout``/``caps`` equal the pinned spec's (a second, content-level check);
  * the banked grade actually ran the code: rows whose agent run failed (Go's
    ``validateSolution`` returns all-false WITHOUT running the code when ``!result.Success``)
    carry an executor-level ``error_category``; they are excluded and counted, because their
    ``stdout_ok`` is not a grade of ``code``.

Banked rows were graded by whatever ailang binary (and the ailang repo's SOURCE stdlib, V37)
was current when they were banked; this grader uses the pinned v0.51.0 tool binary and its
built-in stdlib. Every disagreement is therefore re-run under each available older binary
(``--alt-bin``) and given a TYPED cause, most specific first:
  * ``banked_stdlib_skew_v37`` — the banked grade failed inside ``std/`` after a "stdlib version
    mismatch" warning: the Go grader ran from the ailang repo root and passed
    ``--stdlib-path <repo>/std`` (V37), pairing its binary with a source stdlib of another
    version. The banked verdict is the artefact; the port does not do this by construction;
  * ``banked_grade_timeout`` — the banked grade hit the 10 s grading timeout and the same code
    now completes in under half of it (a loaded rig at banking time);
  * ``version`` — an older binary reproduces the banked verdict (the language moved);
  * ``unexplained`` — otherwise (the grader is wrong until shown otherwise).
AC1.2 holds at 100% agreement, or at >= 99% with zero ``unexplained`` disagreements, and only if
every grade executed the pinned $TOOL digest.

``api_call_json`` solutions have the mock URL of their own run baked in (the harness
substituted ``{{MOCK_HTTP_URL}}`` with an ephemeral port); the re-grade re-creates the mock on
that exact port for the duration of the grade (V8).
"""
from __future__ import annotations

import argparse
import collections
import datetime as dt
import glob
import json
import os
import re
import subprocess
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import corpus as corpus_mod  # noqa: E402
import grade as grade_mod  # noqa: E402

GRADED_CATEGORIES = {'none', 'compile_error', 'runtime_error', 'logic_error'}
# The pinned $TOOL (v0.51.0 darwin/arm64, ~/.pinned-ailang-tools/v0.51.0/ailang) — D-WORLD-55 M0.
PINNED_TOOL_SHA256 = 'e55ff71c710c10395ebf949f6777ef9114bed8bba859e3ba719ed73c6d1b9c7f'


def attribute(rec: dict, banked: dict, regraded: dict) -> str:
    """A typed cause for one disagreement, most specific first; 'unexplained' otherwise."""
    bst = banked.get('stderr') or ''
    # V37: the banked grade ran from the ailang repo root, so `--stdlib-path <repo>/std` pointed
    # the (dev) binary at a SOURCE stdlib of another version; the stdlib itself failed to type.
    if (not banked['stdout_ok'] and regraded['stdout_ok'] and 'stdlib version mismatch' in bst
            and re.search(r'error in std/', bst)):
        return 'banked_stdlib_skew_v37'
    # The banked grade hit the 10 s grading timeout; the same code now completes well inside it.
    if (not banked['stdout_ok'] and regraded['stdout_ok'] and bst.strip().startswith('execution timed out')
            and regraded.get('elapsed_s', grade_mod.TIMEOUT_S) < grade_mod.TIMEOUT_S / 2):
        return 'banked_grade_timeout'
    if any(a['stdout_ok'] == bool(banked['stdout_ok']) for a in rec.get('alt', [])):
        return 'version'
    return 'unexplained'
_MOCK_URL_RE = re.compile(r'http://127\.0\.0\.1:(\d+)')


def yaml_last_change(repo: str, commit: str, file: str) -> dt.datetime:
    out = subprocess.run(['git', '-C', repo, 'log', '-1', '--format=%cI', commit, '--', f'benchmarks/{file}'],
                         capture_output=True, text=True, check=True).stdout.strip()
    return dt.datetime.fromisoformat(out)


def parse_ts(s: str) -> dt.datetime:
    s = s.strip().replace('Z', '+00:00')
    m = re.match(r'(.*T\d\d:\d\d:\d\d)(\.\d+)?(.*)$', s)
    if m and m.group(2):
        s = m.group(1) + m.group(2)[:7] + m.group(3)
    return dt.datetime.fromisoformat(s)


def select_rows(results_dir: str, corpus, repo: str):
    tasks = corpus.by_id()
    since = {t.id: yaml_last_change(repo, corpus.commit, t.file) for t in corpus.tasks}
    stats = collections.Counter()
    rows = []
    paths = sorted(glob.glob(os.path.join(results_dir, '**', 'agent', '*_ailang_*.json'), recursive=True))
    for p in paths:
        base = os.path.basename(p)
        tid = next((t for t in tasks if base.startswith(t + '_ailang_')), None)
        if tid is None:
            continue
        stats['core_files'] += 1
        try:
            d = json.load(open(p))
        except (ValueError, OSError):
            stats['excluded_unreadable'] += 1
            continue
        if d.get('id') != tid or d.get('lang') != 'ailang' or 'stdout_ok' not in d or 'code' not in d:
            stats['excluded_shape'] += 1
            continue
        spec = tasks[tid].spec
        try:
            ts = parse_ts(d.get('timestamp', ''))
        except ValueError:
            stats['excluded_no_timestamp'] += 1
            continue
        if ts <= since[tid]:
            stats['excluded_yaml_changed_since'] += 1
            continue
        if d.get('expected_stdout') is not None and d['expected_stdout'] != spec.get('expected_stdout'):
            stats['excluded_expected_differs'] += 1
            continue
        if d.get('caps') is not None and list(d['caps']) != list(spec.get('caps') or []):
            stats['excluded_caps_differ'] += 1
            continue
        if d.get('error_category') not in GRADED_CATEGORIES:
            stats['excluded_not_graded_' + str(d.get('error_category'))] += 1
            continue
        stats['eligible'] += 1
        rows.append((p, d))
    return rows, stats, {k: v.isoformat() for k, v in since.items()}


def regrade(d: dict, spec: dict, tool_bin: str) -> dict:
    code = d.get('code')
    sol = code.encode() if code is not None else None
    if spec.get('net_allow_localhost') and code:
        ports = sorted(set(int(m) for m in _MOCK_URL_RE.findall(code)))
        if len(ports) == 1:
            try:
                with grade_mod.HTTPMock(ports[0]) as mock:
                    g = grade_mod.grade(sol, spec, tool_bin)
                    g['mock'] = {'port': ports[0], 'hits': len(mock.hits)}
                    return g
            except OSError as e:
                g = grade_mod.grade(sol, spec, tool_bin)
                g['mock'] = {'port': ports[0], 'bind_error': str(e)}
                return g
        g = grade_mod.grade(sol, spec, tool_bin)
        g['mock'] = {'ports_in_source': ports}
        return g
    return grade_mod.grade(sol, spec, tool_bin)


def run(repo: str, commit: str, results_dir: str, tool_bin: str, alt_bins: list, out_path: str | None,
        limit: int | None = None, only: list | None = None) -> dict:
    c = corpus_mod.load_corpus(repo=repo, commit=commit)
    tasks = c.by_id()
    rows, stats, since = select_rows(results_dir, c, repo)
    if only:
        rows = [(p, d) for p, d in rows if d['id'] in only]
    if limit:
        rows = rows[:limit]
    per_task = collections.defaultdict(lambda: collections.Counter())
    disagreements = []
    records = []
    for p, d in rows:
        spec = tasks[d['id']].spec
        g = regrade(d, spec, tool_bin)
        agree = bool(g['stdout_ok']) == bool(d['stdout_ok'])
        per_task[d['id']]['n'] += 1
        per_task[d['id']]['agree'] += int(agree)
        rec = {'path': os.path.relpath(p, repo), 'id': d['id'], 'timestamp': d.get('timestamp'),
               'model': d.get('model'), 'banked_stdout_ok': bool(d['stdout_ok']),
               'banked_error_category': d.get('error_category'), 'stdout_ok': bool(g['stdout_ok']),
               'compile_ok': g['compile_ok'], 'runtime_ok': g['runtime_ok'], 'agree': agree,
               'solution_state': g['solution_state']}
        rec['tool_sha256'] = g['tool_sha256']
        if 'mock' in g:
            rec['mock'] = g['mock']
        if not agree:
            rec['stderr'] = g['stderr'][:600]
            rec['stdout'] = g['stdout'][:300]
            rec['banked_stdout'] = (d.get('stdout') or '')[:300]
            rec['banked_stderr'] = (d.get('stderr') or '')[:600]
            alt = []
            for ab in alt_bins:
                ga = regrade(d, spec, ab)
                alt.append({'bin': ab, 'sha256': ga['tool_sha256'], 'stdout_ok': bool(ga['stdout_ok']),
                            'stderr': ga['stderr'][:300]})
            rec['alt'] = alt
            rec['elapsed_s'] = g.get('elapsed_s')
            rec['cause'] = attribute(rec, d, g)
            disagreements.append(rec)
        records.append(rec)
    n = len(records)
    agree_n = sum(r['agree'] for r in records)
    summary = {
        'corpus_commit': c.commit, 'corpus_digest': c.digest(), 'tool_bin': tool_bin,
        'tool_sha256': grade_mod.sha256_file(tool_bin), 'grader': grade_mod.GRADER_VERSION,
        'alt_bins': [{'bin': b, 'sha256': grade_mod.sha256_file(b)} for b in alt_bins],
        'selection': dict(stats), 'yaml_last_change': since,
        'regraded': n, 'agree': agree_n, 'agreement': (agree_n / n) if n else None,
        'disagreements': len(disagreements),
        'disagreement_causes': dict(collections.Counter(r['cause'] for r in disagreements)),
        'agreement_excluding_version_caused': (
            agree_n / (n - sum(r['cause'] == 'version' for r in disagreements))
            if n - sum(r['cause'] == 'version' for r in disagreements) else None),
        'per_task': {k: dict(v) for k, v in sorted(per_task.items())},
        'banked_pass_rows': sum(r['banked_stdout_ok'] for r in records),
        'executed_tool_sha256': sorted(set(r['tool_sha256'] for r in records if r['tool_sha256'])),
        'api_call_json_mock': dict(collections.Counter(
            ('bound+hit' if r['mock'].get('hits') else 'bound,no-hit' if 'hits' in r['mock']
             else 'bind_error' if 'bind_error' in r['mock'] else 'no-single-port')
            for r in records if 'mock' in r)),
        'banked_fail_rows': sum(not r['banked_stdout_ok'] for r in records),
    }
    if out_path:
        os.makedirs(os.path.dirname(out_path), exist_ok=True)
        with open(out_path, 'w') as f:
            json.dump({'summary': summary, 'disagreements': disagreements}, f, indent=2, sort_keys=True)
            f.write('\n')
        with open(out_path.replace('.json', '.rows.jsonl'), 'w') as f:
            for r in records:
                f.write(json.dumps({k: r[k] for k in ('path', 'id', 'banked_stdout_ok', 'stdout_ok', 'agree')},
                                   sort_keys=True) + '\n')
    return {'summary': summary, 'disagreements': disagreements}


def main(argv=None) -> int:
    ap = argparse.ArgumentParser(description=__doc__.split('\n')[0])
    ap.add_argument('--repo', required=True)
    ap.add_argument('--commit', required=True)
    ap.add_argument('--results-dir')
    ap.add_argument('--tool-bin', default=grade_mod.DEFAULT_TOOL_BIN)
    ap.add_argument('--alt-bin', action='append', default=[])
    ap.add_argument('--out')
    ap.add_argument('--limit', type=int)
    ap.add_argument('--min-rows', type=int, default=200)
    ap.add_argument('--only', action='append', help='restrict to these task ids (mutation proofs)')
    ap.add_argument('--expect-tool-sha256', default=PINNED_TOOL_SHA256)
    a = ap.parse_args(argv)
    res = run(a.repo, a.commit, a.results_dir or os.path.join(a.repo, 'eval_results'), a.tool_bin,
              a.alt_bin, a.out, a.limit, a.only)
    s = res['summary']
    print(json.dumps(s, indent=2, sort_keys=True))
    need = set(a.only) if a.only else {'pipeline', 'cli_args', 'api_call_json'}
    digest_ok = s['executed_tool_sha256'] == [a.expect_tool_sha256]
    if not digest_ok:
        print(f'executed binary digest(s) {s["executed_tool_sha256"]} != pinned $TOOL {a.expect_tool_sha256}')
    ok = (digest_ok and s['regraded'] >= a.min_rows and need <= set(s['per_task'])
          and (s['disagreements'] == 0 or (s['disagreement_causes'].get('unexplained', 0) == 0
                                           and s['agreement'] >= 0.99)))
    print('AC1.2', 'PASS' if ok else 'NOT MET (see disagreements)')
    return 0 if ok else 1


if __name__ == '__main__':
    raise SystemExit(main())
