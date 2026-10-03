#!/usr/bin/env python3
"""AC3.2 / AC3.3 — exact statistics and the verdict matrix (§4.7).

Fixtures sit ON the thresholds: Δ = −2.00 pp exactly passes (T=50, MUT-DELTA-STRICT); a run range
of exactly 5.00 pp is eligible (T=20, MUT-RANGE-STRICT); one outlier task cannot move the median
overhead (MUT-MEAN); even and odd medians; all four verdicts.
"""
import os
import sys
import unittest
from fractions import Fraction

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import classify as C  # noqa: E402
import stats as S  # noqa: E402


def tasks(T):
    return [f't{i:02d}' for i in range(T)]


def rows(arm, T, passes_per_run, wall=lambda t, k: 1000, extra=None):
    """passes_per_run: {run: number of PASS tasks} — the first n tasks pass."""
    out = []
    for k, n in passes_per_run.items():
        for i, t in enumerate(tasks(T)):
            ok = i < n
            r = {'arm': arm, 'task': t, 'run': k, 'stdout_ok': ok,
                 'error_category': 'none' if ok else 'logic_error', 'wall_ms': wall(t, k)}
            if extra:
                r.update(extra(t, k) or {})
            r['class'] = C.classify(r)
            out.append(r)
    return out


class Median(unittest.TestCase):
    def test_odd_even(self):
        self.assertEqual(S.median([3, 1, 2]), 2)
        self.assertEqual(S.median([4, 1, 3, 2]), Fraction(5, 2))
        self.assertEqual(S.median([7]), 7)
        with self.assertRaises(S.StatsError):
            S.median([])


class Eligibility(unittest.TestCase):
    def test_range_exactly_5pp_is_eligible_T20(self):
        e = S.eligibility(rows('shell', 20, {1: 10, 2: 11, 3: 10}), tasks(20))
        self.assertEqual(e['range'], '1/20')
        self.assertTrue(e['eligible'], e)

    def test_range_over_5pp_is_not(self):
        e = S.eligibility(rows('shell', 20, {1: 10, 2: 12, 3: 10}), tasks(20))
        self.assertFalse(e['eligible'])

    def test_T23_one_task_apart_is_eligible_two_is_not(self):
        self.assertTrue(S.eligibility(rows('shell', 23, {1: 20, 2: 21, 3: 21}), tasks(23))['eligible'])
        self.assertFalse(S.eligibility(rows('shell', 23, {1: 19, 2: 21, 3: 21}), tasks(23))['eligible'])

    def test_needs_three_runs(self):
        e = S.eligibility(rows('shell', 23, {1: 20, 2: 20}), tasks(23))
        self.assertFalse(e['eligible'])
        self.assertIn('N=2 < 3', e['reasons'])

    def test_incomplete_run_refused(self):
        r = rows('shell', 23, {1: 20, 2: 20, 3: 20})[:-1]
        with self.assertRaises(S.StatsError):
            S.eligibility(r, tasks(23))


class PassRate(unittest.TestCase):
    def test_delta_exactly_minus_2pp_passes_T50(self):
        s = rows('shell', 50, {1: 40, 2: 40, 3: 40})
        w = rows('world', 50, {1: 39, 2: 39, 3: 39})  # ΣW = ΣS − 3 over 150 task-runs
        p = S.pass_rate_test(s, w, tasks(50))
        self.assertEqual(p['delta'], Fraction(-2, 100))
        self.assertTrue(p['passes'])

    def test_delta_below_fails(self):
        s = rows('shell', 50, {1: 40, 2: 40, 3: 40})
        w = rows('world', 50, {1: 39, 2: 39, 3: 38})
        self.assertFalse(S.pass_rate_test(s, w, tasks(50))['passes'])

    def test_T23_N3_world_may_lose_one(self):
        s = rows('shell', 23, {1: 21, 2: 21, 3: 21})
        self.assertTrue(S.pass_rate_test(s, rows('world', 23, {1: 21, 2: 20, 3: 21}), tasks(23))['passes'])
        self.assertFalse(S.pass_rate_test(s, rows('world', 23, {1: 20, 2: 20, 3: 21}), tasks(23))['passes'])

    def test_world_final_fault_blocks(self):
        s = rows('shell', 23, {1: 21, 2: 21, 3: 21})
        w = rows('world', 23, {1: 21, 2: 21, 3: 21},
                 extra=lambda t, k: {'stdout_ok': False, 'error_category': 'quota_exhausted',
                                     'cause_evidence': '429'} if (t, k) == ('t00', 1) else None)
        with self.assertRaises(S.StatsError):
            S.agent_tests(s, w, tasks(23))

    def test_rerun_attempt_replaces_the_faulted_one(self):
        s = rows('shell', 23, {1: 21, 2: 21, 3: 21})
        w = rows('world', 23, {1: 21, 2: 21, 3: 21},
                 extra=lambda t, k: {'stdout_ok': False, 'error_category': 'quota_exhausted',
                                     'cause_evidence': '429'} if (t, k) == ('t00', 1) else None)
        again = dict(next(r for r in w if (r['task'], r['run']) == ('t00', 1)))
        again.update(stdout_ok=True, error_category='none', attempt=2)
        again.pop('cause_evidence')
        again['class'] = C.classify(again)
        res = S.agent_tests(s, w + [again], tasks(23))
        self.assertEqual(res['pass_rate']['delta'], 0)


class Overhead(unittest.TestCase):
    def test_median_of_per_task_medians_odd_N(self):
        T = 5
        sw = {'t00': 1000, 't01': 1000, 't02': 1000, 't03': 1000, 't04': 1000}
        ww = {'t00': 1000, 't01': 1000, 't02': 1100, 't03': 1100, 't04': 11000}  # one outlier task
        s = rows('shell', T, {1: 5, 2: 5, 3: 5}, wall=lambda t, k: sw[t])
        w = rows('world', T, {1: 5, 2: 5, 3: 5}, wall=lambda t, k: ww[t])
        o = S.overhead(s, w, tasks(T))
        self.assertEqual(o['o_t']['t04'], 10)
        self.assertEqual(o['O'], Fraction(1, 10))  # the mean would be 2.04 and FAIL
        self.assertTrue(S.agent_tests(s, w, tasks(T))['overhead_passes'])

    def test_even_N_per_task_median_and_even_T(self):
        T = 4
        s = rows('shell', T, {1: 4, 2: 4}, wall=lambda t, k: 1000 if k == 1 else 3000)  # median 2000
        w = rows('world', T, {1: 4, 2: 4}, wall=lambda t, k: 2000 if k == 1 else 3000)  # median 2500
        o = S.overhead(s, w, tasks(T))
        self.assertEqual(set(o['o_t'].values()), {Fraction(1, 4)})
        self.assertEqual(o['O'], Fraction(1, 4))  # exactly 0.25 passes (≤)

    def test_timeout_counts_600s_and_faults_are_excluded(self):
        T = 3
        s = rows('shell', T, {1: 3, 2: 3, 3: 3}, wall=lambda t, k: 300_000)
        w = rows('world', T, {1: 3, 2: 3, 3: 3}, wall=lambda t, k: 300_000,
                 extra=lambda t, k: {'stdout_ok': False, 'error_category': 'timeout', 'finish_reason': 'timeout',
                                     'wall_ms': 900_000} if t == 't00' and k in (1, 2) else None)
        o = S.overhead(s, w, tasks(T))
        self.assertEqual(o['o_t']['t00'], 1)  # median(600,600,300) s / 300 s − 1
        s2 = rows('shell', T, {1: 3, 2: 3, 3: 3}, wall=lambda t, k: 300_000,
                  extra=lambda t, k: {'stdout_ok': False, 'error_category': 'api_error', 'cause_evidence': 'x',
                                      'wall_ms': 1} if (t, k) == ('t01', 1) else None)
        self.assertEqual(S.overhead(s2, rows('world', T, {1: 3, 2: 3, 3: 3}, wall=lambda t, k: 300_000),
                                    tasks(T))['o_t']['t01'], 0)


class Verdict(unittest.TestCase):
    P, F = {'passes': True}, {'passes': False}

    def test_matrix(self):
        V = S.verdict
        self.assertEqual(V({'cc': False, 'cx': False}, {}), S.PAUSED)
        self.assertEqual(V({'cc': True, 'cx': True}, {'cc': self.P, 'cx': self.F}), S.FAILS)
        self.assertEqual(V({'cc': True, 'cx': False}, {'cc': self.F}), S.FAILS)
        self.assertEqual(V({'cc': True, 'cx': True}, {'cc': self.P, 'cx': self.P}), S.HOLDS)
        self.assertEqual(V({'cc': True, 'cx': False}, {'cc': self.P}), S.INCOMPLETE)
        self.assertEqual(V({'cc': False, 'cx': True}, {'cx': self.P}), S.INCOMPLETE)
        with self.assertRaises(S.StatsError):
            V({'cc': True}, {})

    def test_end_to_end_fails_on_overhead_alone(self):
        T = 23
        s = rows('shell', T, {1: 21, 2: 21, 3: 21}, wall=lambda t, k: 10_000)
        w = rows('world', T, {1: 21, 2: 21, 3: 21}, wall=lambda t, k: 12_600)  # o_t = 0.26
        t = S.agent_tests(s, w, tasks(T))
        self.assertTrue(t['pass_rate']['passes'])
        self.assertFalse(t['overhead_passes'])
        self.assertEqual(S.verdict({'cc': True, 'cx': True}, {'cc': t, 'cx': t}), S.FAILS)

    def test_informative_is_labelled_and_deterministic(self):
        T = 23
        s = rows('shell', T, {1: 21, 2: 21, 3: 21})
        w = rows('world', T, {1: 20, 2: 21, 3: 21})
        a, b = S.informative(s, w, tasks(T)), S.informative(s, w, tasks(T))
        self.assertEqual(a, b)
        self.assertIn('never the gate', a['note'])
        lo, hi = a['delta_bootstrap_ci95']
        self.assertLessEqual(lo, hi)


if __name__ == '__main__':
    unittest.main()
