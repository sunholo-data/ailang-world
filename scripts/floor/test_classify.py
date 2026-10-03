#!/usr/bin/env python3
"""AC3.1 — the pre-registered classification table, one row per category × arm (§4.6).

MUT-FAULT-AS-FAIL (a shell api_error classified FAIL) is killed by the table row and by the
eligibility arm: one evidenced shell api_error must make the agent INELIGIBLE, never cost it a
task.
"""
import os
import sys
import unittest

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import classify as C  # noqa: E402
import stats  # noqa: E402

EV = 'codex exit rc=1 after 180 s idle; no provider error in the transcript'

# (category, arm, extra fields, expected class)
TABLE = [
    ('compile_error', 'shell', {}, C.FAIL), ('compile_error', 'world', {}, C.FAIL),
    ('runtime_error', 'shell', {}, C.FAIL), ('runtime_error', 'world', {}, C.FAIL),
    ('logic_error', 'shell', {}, C.FAIL), ('logic_error', 'world', {}, C.FAIL),
    ('step_exhausted', 'shell', {}, C.FAIL), ('step_exhausted', 'world', {}, C.FAIL),
    ('non_agentic', 'shell', {}, C.FAIL), ('non_agentic', 'world', {}, C.FAIL),
    ('refused', 'shell', {}, C.FAIL), ('refused', 'world', {}, C.FAIL),
    ('output_format', 'shell', {}, C.FAIL), ('output_format', 'world', {}, C.FAIL),
    ('reasoning_stall', 'shell', {}, C.FAIL), ('reasoning_stall', 'world', {}, C.FAIL),
    ('timeout', 'shell', {'cause_evidence': 'deadline 600 s'}, C.HARNESS_FAULT),
    ('timeout', 'world', {}, C.FAIL),
    ('world_tool_error', 'world', {}, C.FAIL),
    ('api_error', 'shell', {'cause_evidence': EV}, C.HARNESS_FAULT),
    ('api_error', 'world', {'cause_evidence': EV}, C.FAIL),  # evidenced but untyped: World pays
    ('api_error', 'world', {}, C.FAIL),
    ('api_error', 'world', {'cause_evidence': 'HTTP 529 overloaded in transcript', 'cause_typed': True},
     C.HARNESS_FAULT),
    ('resource_limit', 'shell', {'cause_evidence': EV}, C.HARNESS_FAULT),
    ('resource_limit', 'world', {'cause_evidence': EV}, C.HARNESS_FAULT),
    ('quota_exhausted', 'shell', {'cause_evidence': EV}, C.HARNESS_FAULT),
    ('quota_exhausted', 'world', {'cause_evidence': EV}, C.HARNESS_FAULT),
    ('rate_limit', 'shell', {'cause_evidence': EV}, C.HARNESS_FAULT),
    ('rate_limit', 'world', {'cause_evidence': EV}, C.HARNESS_FAULT),
    ('wire_drift', 'shell', {'cause_evidence': EV}, C.HARNESS_FAULT),
    ('wire_drift', 'world', {'cause_evidence': EV}, C.HARNESS_FAULT),
    ('harness_setup', 'shell', {'cause_evidence': 'mock bind failed'}, C.HARNESS_FAULT),
    ('harness_setup', 'world', {'cause_evidence': 'pre-flight tools/list = 7'}, C.HARNESS_FAULT),
]
REFUSED = [
    ('verify_error', 'shell', {}), ('constraint_violation', 'world', {}), ('cost_killed', 'shell', {}),
    ('thrash_aborted', 'world', {}), ('policy_violation', 'shell', {}), ('none', 'shell', {}),
    ('world_tool_error', 'shell', {}),
    ('api_error', 'shell', {}),            # a harness fault must carry cause_evidence
    ('timeout', 'shell', {}),
    ('harness_setup', 'world', {}),
]


def row(cat, arm, ok=False, **extra):
    return dict({'arm': arm, 'stdout_ok': ok, 'error_category': cat}, **extra)


class Table(unittest.TestCase):
    def test_every_category_by_arm(self):
        for cat, arm, extra, want in TABLE:
            with self.subTest(cat=cat, arm=arm, extra=sorted(extra)):
                self.assertEqual(C.classify(row(cat, arm, **extra)), want)

    def test_pass_is_stdout_ok(self):
        for arm in ('shell', 'world'):
            for cat in ('none', 'timeout', 'api_error'):
                self.assertEqual(C.classify(row(cat, arm, ok=True)), C.PASS)

    def test_unlisted_refused(self):
        for cat, arm, extra in REFUSED:
            with self.subTest(cat=cat, arm=arm):
                with self.assertRaises(C.ClassificationError):
                    C.classify(row(cat, arm, **extra))


class FaultMakesIneligible(unittest.TestCase):
    def test_shell_api_error_is_ineligibility_not_a_lost_task(self):
        tasks = [f't{i}' for i in range(23)]
        rows = []
        for k in (1, 2, 3):
            for t in tasks:
                if (k, t) == (2, 't5'):
                    r = row('api_error', 'shell', cause_evidence=EV)
                else:
                    r = row('none', 'shell', ok=True)
                r.update(task=t, run=k, wall_ms=1000)
                r['class'] = C.classify(r)
                rows.append(r)
        e = stats.eligibility(rows, tasks)
        self.assertFalse(e['eligible'])
        self.assertEqual(e['harness_faults'], [['t5', 2]])
        # counted as a FAIL instead, the run range would be 1/23 and the agent eligible
        self.assertEqual(e['range'], '1/23')


if __name__ == '__main__':
    unittest.main()
