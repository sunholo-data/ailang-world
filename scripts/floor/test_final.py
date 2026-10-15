#!/usr/bin/env python3
"""Row 93 M5 — the FINAL phase guards (§4.11 order, write-once rows), the paired order (D-WORLD-74),
the mint block (§4.5), the drift readout (informative) and the prereg task list report.py reads."""
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

    def test_old_sequential_phases_are_gone(self):
        # D-WORLD-74: a sequential World phase would put the arms in separate time windows again.
        for phase in ('world', 'drift'):
            with self.assertRaisesRegex(run.RunRefused, 'unknown phase'):
                run.check_final_phase(self.ev, phase, 1, 3, AGENTS)

    def test_paired_refused_before_the_seal(self):
        with self.assertRaisesRegex(run.RunRefused, '§4.11'):
            run.check_final_phase(self.ev, 'paired', 1, 3, AGENTS)

    def test_paired_refused_with_an_uncommitted_seal(self):
        for a in AGENTS:
            write_runs(self.ev, a, 'shell', {1: 21, 2: 21, 3: 21})
        report.seal(self.ev, AGENTS, T)
        with self.assertRaisesRegex(run.RunRefused, 'not committed'):
            run.check_final_phase(self.ev, 'paired', 1, 3, AGENTS)

    def test_paired_allowed_after_the_committed_seal(self):
        self.seal()
        for k in (1, 2, 3):
            run.check_final_phase(self.ev, 'paired', k, 3, AGENTS)
        with self.assertRaisesRegex(run.RunRefused, r'--run must be 1\.\.3'):
            run.check_final_phase(self.ev, 'paired', 4, 3, AGENTS)

    def test_rows_are_written_once(self):
        run.check_final_phase(self.ev, 'shell', 2, 3, ['claude', 'codex'])  # nothing exists yet
        os.makedirs(os.path.dirname(run.final_rows_path(self.ev, 'shell', 'claude', 2)))
        open(run.final_rows_path(self.ev, 'shell', 'claude', 2), 'w').close()
        with self.assertRaisesRegex(run.RunRefused, 'written once'):
            run.check_final_phase(self.ev, 'shell', 2, 3, ['claude', 'codex'])
        run.check_final_phase(self.ev, 'shell', 3, 3, ['claude', 'codex'])

    def test_paired_rows_are_written_once_per_arm(self):
        self.seal()
        for arm in ('shell', 'world'):
            p = run.final_rows_path(self.ev, 'paired', 'codex', 1, arm)
            os.makedirs(os.path.dirname(p), exist_ok=True)
            open(p, 'w').close()
            with self.assertRaisesRegex(run.RunRefused, 'written once'):
                run.check_final_phase(self.ev, 'paired', 1, 3, AGENTS)
            os.remove(p)
        run.check_final_phase(self.ev, 'paired', 1, 3, AGENTS)

    def test_rows_paths(self):
        self.assertEqual(run.final_rows_path('E', 'shell', 'codex', 2), os.path.join('E', 'runs', 'codex', 'shell', 'r2.jsonl'))
        self.assertEqual(run.final_rows_path('E', 'paired', 'codex', 2, 'world'),
                         os.path.join('E', 'paired', 'runs', 'codex', 'world', 'r2.jsonl'))
        self.assertEqual(run.final_rows_path('E', 'paired', 'claude', 1, 'shell'),
                         os.path.join('E', 'paired', 'runs', 'claude', 'shell', 'r1.jsonl'))
        with self.assertRaises(run.RunRefused):
            run.final_rows_path('E', 'shell', 'codex', 1, 'world')  # Phase 1 has no World rows


class PairedOrder(unittest.TestCase):
    """D-WORLD-74: both arms of a task run back to back, the order counterbalanced."""

    def test_each_task_runs_both_arms_once(self):
        for k in (1, 2, 3):
            for i in range(23):
                self.assertEqual(sorted(run.paired_order(i, k)), ['shell', 'world'])

    def test_each_run_is_balanced(self):
        for k in (1, 2, 3):
            firsts = [run.paired_order(i, k)[0] for i in range(23)]
            self.assertLessEqual(abs(firsts.count('world') - firsts.count('shell')), 1)

    def test_order_flips_between_consecutive_runs(self):
        for k in (1, 2):
            for i in range(23):
                self.assertEqual(run.paired_order(i, k), tuple(reversed(run.paired_order(i, k + 1))))

    def test_adjacent_tasks_alternate(self):
        for i in range(22):
            self.assertNotEqual(run.paired_order(i, 1)[0], run.paired_order(i + 1, 1)[0])


class IterateSchedule(unittest.TestCase):
    """The attended smoke's --paired: canaries first, then each task's arms in paired_order(i, 1)."""
    T = ['a', 'b', 'c']

    def test_unpaired_keeps_the_m4_order(self):
        s = run.iterate_schedule(['codex'], ['shell', 'world'], self.T, set(), False)
        self.assertEqual([(k, x, t) for k, _, x, t in s],
                         [('canary', 'shell', None)] + [('task', 'shell', t) for t in self.T]
                         + [('canary', 'world', None)] + [('task', 'world', t) for t in self.T])

    def test_paired_interleaves_in_the_final_order(self):
        s = run.iterate_schedule(['claude', 'codex'], ['shell', 'world'], self.T, set(), True)
        for ag in ('claude', 'codex'):
            mine = [(k, x, t) for k, a, x, t in s if a == ag]
            self.assertEqual(mine[:2], [('canary', 'shell', None), ('canary', 'world', None)])
            tasks = mine[2:]
            self.assertEqual(len(tasks), 2 * len(self.T))
            for i, t in enumerate(self.T):
                pair = tasks[2 * i:2 * i + 2]
                self.assertEqual([t, t], [p[2] for p in pair])
                self.assertEqual(run.paired_order(i, 1), tuple(p[1] for p in pair))

    def test_paired_refuses_one_arm_or_reuse(self):
        with self.assertRaises(run.RunRefused):
            run.iterate_schedule(['codex'], ['world'], self.T, set(), True)
        with self.assertRaises(run.RunRefused):
            run.iterate_schedule(['codex'], ['shell', 'world'], self.T, {'claude/shell'}, True)

    def test_reused_cell_is_not_scheduled(self):
        s = run.iterate_schedule(['codex'], ['shell', 'world'], self.T, {'codex/shell'}, False)
        self.assertNotIn('shell', {x for _, _, x, _ in s})


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
        write_runs(os.path.join(self.ev, 'paired'), 'a', 'shell', {1: 21, 2: 21, 3: 21}, wall=11_000)  # +10%
        write_runs(self.ev, 'b', 'shell', {1: 21, 2: 21, 3: 21}, wall=10_000)
        write_runs(os.path.join(self.ev, 'paired'), 'b', 'shell', {1: 21, 2: 21, 3: 21}, wall=11_100)  # +11%
        d = report.drift(self.ev, ['a', 'b', 'c'], T)
        self.assertFalse(d['a']['drift_confounded'])
        self.assertTrue(d['b']['drift_confounded'])
        self.assertEqual(d['a']['tasks'], 23)
        self.assertNotIn('c', d)  # no probe rows: nothing reported

    def test_negative_drift_counts_too(self):
        write_runs(self.ev, 'a', 'shell', {1: 21, 2: 21, 3: 21}, wall=10_000)
        write_runs(os.path.join(self.ev, 'paired'), 'a', 'shell', {1: 21}, wall=8_000)
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
