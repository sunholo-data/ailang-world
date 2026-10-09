#!/usr/bin/env python3
"""Row 93 M5 — the FINAL phase guards (§4.11 order, write-once rows), the mint block (§4.5), the
drift probe (§4.11, informative) and the prereg task list report.py reads."""
import glob
import json
import os
import shutil
import subprocess
import sys
import tempfile
import unittest

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import arms  # noqa: E402
import report  # noqa: E402
import run  # noqa: E402
from test_report import AGENTS, GIT, T, write_runs  # noqa: E402


class PhaseGuards(unittest.TestCase):
    def setUp(self):
        self.repo = tempfile.mkdtemp(prefix='floor-final-')
        subprocess.run(GIT + ['init', '-q', self.repo], check=True)
        self.ev = os.path.join(self.repo, 'ev')
        os.makedirs(self.ev)

    def tearDown(self):
        shutil.rmtree(self.repo, ignore_errors=True)

    def seal(self):
        for a in AGENTS:
            write_runs(self.ev, a, 'shell', {1: 21, 2: 21, 3: 21})
        report.seal(self.ev, AGENTS, T)
        subprocess.run(GIT + ['-C', self.repo, 'add', '-A'], check=True)
        subprocess.run(GIT + ['-C', self.repo, 'commit', '-q', '-m', 'seal'], check=True)

    def test_shell_runs_need_no_seal(self):
        run.check_final_phase(self.ev, 'shell', 1, 3, AGENTS)

    def test_run_out_of_range_refused(self):
        for k in (0, 4):
            with self.assertRaisesRegex(run.RunRefused, r'--run must be 1\.\.3'):
                run.check_final_phase(self.ev, 'shell', k, 3, AGENTS)

    def test_world_and_drift_refused_before_the_seal(self):
        for phase, k in (('world', 1), ('drift', 4)):
            with self.assertRaisesRegex(run.RunRefused, '§4.11'):
                run.check_final_phase(self.ev, phase, k, 3, AGENTS)

    def test_world_refused_with_an_uncommitted_seal(self):
        for a in AGENTS:
            write_runs(self.ev, a, 'shell', {1: 21, 2: 21, 3: 21})
        report.seal(self.ev, AGENTS, T)
        with self.assertRaisesRegex(run.RunRefused, 'not committed'):
            run.check_final_phase(self.ev, 'world', 1, 3, AGENTS)

    def test_world_allowed_after_the_committed_seal(self):
        self.seal()
        run.check_final_phase(self.ev, 'world', 1, 3, AGENTS)
        run.check_final_phase(self.ev, 'drift', 4, 3, AGENTS)

    def test_rows_are_written_once(self):
        run.check_final_phase(self.ev, 'shell', 2, 3, ['claude', 'codex'])  # nothing exists yet
        os.makedirs(os.path.dirname(run.final_rows_path(self.ev, 'shell', 'claude', 2)))
        open(run.final_rows_path(self.ev, 'shell', 'claude', 2), 'w').close()
        with self.assertRaisesRegex(run.RunRefused, 'written once'):
            run.check_final_phase(self.ev, 'shell', 2, 3, ['claude', 'codex'])
        run.check_final_phase(self.ev, 'shell', 3, 3, ['claude', 'codex'])

    def test_rows_paths(self):
        self.assertEqual(run.final_rows_path('E', 'world', 'codex', 2), os.path.join('E', 'runs', 'codex', 'world', 'r2.jsonl'))
        self.assertEqual(run.final_rows_path('E', 'drift', 'claude', 4),
                         os.path.join('E', 'drift', 'runs', 'claude', 'shell', 'r4.jsonl'))


class MintBlock(unittest.TestCase):
    def test_two_n_sessions_with_explicit_grants(self):
        L = {'root': '/R'}
        b = run.final_mint_block(L, 3)
        news = [ln for ln in b.splitlines() if ' session new ' in ln]
        self.assertEqual([ln.split()[3] for ln in news],
                         ['fl-cc-r1', 'fl-cx-r1', 'fl-cc-r2', 'fl-cx-r2', 'fl-cc-r3', 'fl-cx-r3'])
        self.assertNotIn('--preset', '\n'.join(ln for ln in b.splitlines() if not ln.startswith('#')))
        self.assertEqual(b.count(' '.join(arms.world_grant_args())), 6)
        self.assertEqual(b.count(f'--ttl {arms.SESSION_TTL_S}'), 6)
        self.assertEqual(b.count('world-publish transitions'), 1)
        for k in (1, 2, 3):
            for ep in run.final_episodes(k).values():
                self.assertIn(f'--out $FT/private/{ep}.session', b)


class Drift(unittest.TestCase):
    def setUp(self):
        self.ev = tempfile.mkdtemp(prefix='floor-drift-')

    def tearDown(self):
        shutil.rmtree(self.ev, ignore_errors=True)

    def test_flags_only_beyond_ten_percent(self):
        write_runs(self.ev, 'a', 'shell', {1: 21, 2: 21, 3: 21}, wall=10_000)
        write_runs(os.path.join(self.ev, 'drift'), 'a', 'shell', {4: 21}, wall=11_000)  # +10%: not beyond
        write_runs(self.ev, 'b', 'shell', {1: 21, 2: 21, 3: 21}, wall=10_000)
        write_runs(os.path.join(self.ev, 'drift'), 'b', 'shell', {4: 21}, wall=11_100)  # +11%
        d = report.drift(self.ev, ['a', 'b', 'c'], T)
        self.assertFalse(d['a']['drift_confounded'])
        self.assertTrue(d['b']['drift_confounded'])
        self.assertEqual(d['a']['tasks'], 23)
        self.assertNotIn('c', d)  # no probe rows: nothing reported

    def test_negative_drift_counts_too(self):
        write_runs(self.ev, 'a', 'shell', {1: 21, 2: 21, 3: 21}, wall=10_000)
        write_runs(os.path.join(self.ev, 'drift'), 'a', 'shell', {4: 21}, wall=8_000)
        self.assertTrue(report.drift(self.ev, ['a'], T)['a']['drift_confounded'])


class PreregTasks(unittest.TestCase):
    def test_reads_the_draft_this_harness_writes(self):
        drafts = glob.glob(os.path.join(run.REPO, 'design_docs', 'verification', 'world-floor-*', 'preregistration*.json'))
        self.assertTrue(drafts)
        for p in drafts:
            with open(p) as f:
                pre = json.load(f)
            tasks = report.prereg_tasks(pre)
            self.assertEqual(len(tasks), 23, p)
            self.assertEqual(sorted(tasks), sorted(pre['live_config']['corpus']['tasks']))

    def test_refuses_an_empty_prereg(self):
        with self.assertRaises(SystemExit):
            report.prereg_tasks({'live_config': {'corpus': {}}})


if __name__ == '__main__':
    unittest.main()
