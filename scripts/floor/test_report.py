#!/usr/bin/env python3
"""AC3.4 — report.py refuses ANY World statistic until eligibility.json is sealed (committed).

Each refusal arm also asserts that no World statistic was computed and no World row was read
before the refusal (MUT-SEAL-ORDER: computing World stats first and checking the seal after
would still raise, so the raise alone is not the assertion — the call log is).
"""
import json
import os
import shutil
import subprocess
import sys
import tempfile
import unittest

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import report  # noqa: E402
import stats  # noqa: E402

T = [f't{i:02d}' for i in range(23)]
AGENTS = ['claude-code', 'codex']
GIT = ['git', '-c', 'user.name=floor-test', '-c', 'user.email=floor-test@example.com', '-c', 'commit.gpgsign=false']


def write_runs(root, agent, arm, passes, wall=10_000):
    d = os.path.join(root, 'runs', agent, arm)
    os.makedirs(d, exist_ok=True)
    for k, n in passes.items():
        with open(os.path.join(d, f'r{k}.jsonl'), 'w') as f:
            for i, t in enumerate(T):
                ok = i < n
                f.write(json.dumps({'arm': arm, 'task': t, 'run': k, 'stdout_ok': ok,
                                    'error_category': 'none' if ok else 'logic_error', 'wall_ms': wall}) + '\n')


class Seal(unittest.TestCase):
    def setUp(self):
        self.repo = tempfile.mkdtemp(prefix='floor-report-')
        subprocess.run(GIT + ['init', '-q', self.repo], check=True)
        self.ev = os.path.join(self.repo, 'world-floor-test')
        for a in AGENTS:
            write_runs(self.ev, a, 'shell', {1: 21, 2: 21, 3: 22})
            write_runs(self.ev, a, 'world', {1: 21, 2: 21, 3: 21}, wall=11_000)
        self.calls = []
        self._orig = (stats.agent_tests, report.load_rows)

        def spy_tests(*a, **k):
            self.calls.append('agent_tests')
            return self._orig[0](*a, **k)

        def spy_load(ev, agent, arm):
            self.calls.append(f'load:{agent}:{arm}')
            return self._orig[1](ev, agent, arm)
        stats.agent_tests = spy_tests
        report.load_rows = spy_load

    def tearDown(self):
        stats.agent_tests, report.load_rows = self._orig
        shutil.rmtree(self.repo, ignore_errors=True)

    def git(self, *args):
        subprocess.run(GIT + ['-C', self.repo, *args], check=True, capture_output=True)

    def world_touched(self):
        return [c for c in self.calls if c == 'agent_tests' or c.endswith(':world')]

    def refuse(self):
        with self.assertRaises(report.SealError):
            report.compute_verdict(self.ev, AGENTS, T)
        self.assertEqual(self.world_touched(), [])

    def test_absent(self):
        self.refuse()

    def test_untracked(self):
        report.seal(self.ev, AGENTS, T)
        self.refuse()

    def test_staged_but_not_committed(self):
        report.seal(self.ev, AGENTS, T)
        self.git('add', '.')
        self.refuse()

    def test_committed_then_edited(self):
        p = report.seal(self.ev, AGENTS, T)
        self.git('add', '.')
        self.git('commit', '-qm', 'seal')
        with open(p) as f:
            d = json.load(f)
        d['agents']['codex']['eligible'] = False
        with open(p, 'w') as f:
            json.dump(d, f)
        self.refuse()

    def test_committed_but_swapped(self):
        p = report.seal(self.ev, AGENTS, T)
        with open(p) as f:
            d = json.load(f)
        d['agents']['codex']['eligible'] = False
        with open(p, 'w') as f:
            json.dump(d, f, indent=2, sort_keys=True)
        self.git('add', '.')
        self.git('commit', '-qm', 'seal (tampered)')
        self.refuse()

    def test_sealed_and_committed_computes(self):
        report.seal(self.ev, AGENTS, T)
        self.git('add', '.')
        self.git('commit', '-qm', 'seal eligibility first')
        r = report.compute_verdict(self.ev, AGENTS, T)
        self.assertEqual(r['verdict'], stats.HOLDS)  # Δ = −1/69 ≥ −2/100; O = 0.10
        self.assertEqual(r['tests']['codex']['pass_rate']['delta']['exact'], '-1/69')
        self.assertIn('agent_tests', self.calls)
        self.assertIn('Floor verdict: HOLDS', report.summary_md(r))

    def test_seal_is_once(self):
        report.seal(self.ev, AGENTS, T)
        with self.assertRaises(report.SealError):
            report.seal(self.ev, AGENTS, T)

    def test_seal_reads_no_world_rows(self):
        report.seal(self.ev, AGENTS, T)
        self.assertEqual(self.world_touched(), [])


if __name__ == '__main__':
    unittest.main()
