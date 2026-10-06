#!/usr/bin/env python3
"""M4 — run.py's guards and pure parts (design §6 M4 AC4.1/AC4.2, M5 AC5.1, §4.9, §4.10).

Mutants killed here (each turns the NAMED test red; README of world-floor-m4-2026-10-06 records
the runs):
  * MUT-DIRTY-LEDGER — ``require_clean_ledger`` accepts a dirty ledger
    -> ``LedgerGuard.test_refuses_dirty_ledger`` (+ ``test_iterate_refuses_before_any_work``).
  * MUT-FINAL-NOPREREG — ``check_final`` drops ``require_committed``
    -> ``FinalGuard.test_final_refuses_uncommitted_prereg``.
"""
import io
import json
import os
import shutil
import subprocess
import sys
import tempfile
import unittest
from contextlib import redirect_stderr, redirect_stdout
from types import SimpleNamespace
from unittest import mock

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import arms  # noqa: E402
import classify  # noqa: E402
import prompt  # noqa: E402
import run  # noqa: E402

with open(os.path.join(HERE, 'testdata', 'agent_task_ailang.txt')) as _f:
    TEMPLATE = _f.read()


def git(cwd, *args):
    subprocess.run(['git', '-c', 'user.name=t', '-c', 'user.email=t@t.invalid', *args], cwd=cwd,
                   check=True, capture_output=True)


class TempRepo:
    def __enter__(self):
        self.d = tempfile.mkdtemp(prefix='floor-run-test-')
        git(self.d, 'init', '-q')
        return self

    def __exit__(self, *exc):
        shutil.rmtree(self.d, ignore_errors=True)

    def path(self, name):
        return os.path.join(self.d, name)

    def write(self, name, text):
        with open(self.path(name), 'w') as f:
            f.write(text)
        return self.path(name)

    def commit(self, *names):
        git(self.d, 'add', '--', *names)
        git(self.d, 'commit', '-q', '-m', 'c')


def fake_task(i, tier='core'):
    return SimpleNamespace(id=i, sha256=('%064x' % abs(hash((i, tier))))[:64],
                           spec={'id': i, 'caps': ['IO'], 'expected_stdout': 'x\n'})


CORE_IDS = sorted(['pipeline', 'cli_args', 'api_call_json', 'prompt_injection', 'effect_tracking_io_fs',
                   'higher_order_functions'] + [f'core_{n:02d}' for n in range(17)])
SMOKE_IDS = sorted(['fizzbuzz'] + [f'smoke_{n:02d}' for n in range(22)])


class FakeProbes(run.Probes):
    def __init__(self, cli=None, teaching='TEACH'):
        super().__init__('/nonexistent')
        self.cli = cli or {'claude': '2.1.291 (Claude Code)', 'codex': 'codex-cli 0.159.2'}
        self.teach = teaching

    def tool_sha256(self, path):
        return 'f' * 64

    def teaching(self, path, compact):
        return self.teach + (' compact' if compact else '')

    def template(self, commit):
        return TEMPLATE

    def corpus(self, commit, tier, expected):
        ids = CORE_IDS if tier == 'core' else SMOKE_IDS
        assert len(ids) == expected == 23
        return SimpleNamespace(tasks=tuple(fake_task(i, tier) for i in ids))

    def cli_versions(self):
        return dict(self.cli)


def cfg(**knobs):
    c = run.load_config(run.CONFIG)
    c['knobs'] = dict(c['knobs'], **knobs)
    return c


class LedgerGuard(unittest.TestCase):  # AC4.2
    def test_refuses_dirty_ledger(self):  # MUT-DIRTY-LEDGER
        with TempRepo() as r:
            led = r.write('ledger.jsonl', '{"iteration": 0}\n')
            r.commit('ledger.jsonl')
            run.require_clean_ledger(led)  # committed and clean: accepted
            with open(led, 'a') as f:
                f.write('{"iteration": 1}\n')
            with self.assertRaisesRegex(run.RunRefused, 'AC4.2'):
                run.require_clean_ledger(led)

    def test_refuses_untracked_ledger(self):
        with TempRepo() as r:
            r.write('other', 'x')
            r.commit('other')
            led = r.write('ledger.jsonl', '{"iteration": 0}\n')
            with self.assertRaisesRegex(run.RunRefused, 'AC4.2'):
                run.require_clean_ledger(led)

    def test_absent_ledger_is_clean(self):  # the first iteration creates it
        with TempRepo() as r:
            r.write('other', 'x')
            r.commit('other')
            run.require_clean_ledger(r.path('ledger.jsonl'))

    def test_iterate_refuses_before_any_work(self):  # MUT-DIRTY-LEDGER, end to end through main()
        with TempRepo() as r:
            led = r.write('ledger.jsonl', '{"iteration": 0}\n')
            r.commit('ledger.jsonl')
            with open(led, 'a') as f:
                f.write('dirty\n')
            ev = os.path.join(r.d, 'evidence')
            err = io.StringIO()
            with mock.patch.object(run, 'LEDGER', led), mock.patch.object(run, 'load_config') as lc, \
                    redirect_stderr(err):
                rc = run.main(['iterate', '--iteration', '1', '--knob', 'k', '--side', 'harness',
                               '--change', 'c', '--evidence', ev])
            self.assertEqual(rc, 2)
            self.assertIn('AC4.2', err.getvalue())
            lc.assert_not_called()  # refused before reading the config, spawning or writing evidence
            self.assertFalse(os.path.exists(ev))


class FinalGuard(unittest.TestCase):  # AC5.1
    def setUp(self):
        self.cfg = cfg()
        self.probes = FakeProbes()

    def prereg(self, r, probes=None):
        pre = run.prereg_draft(self.cfg, probes or self.probes, ledger=r.path('ledger.jsonl'))
        return r.write('preregistration.json', json.dumps(pre, indent=2, sort_keys=True))

    def test_final_refuses_uncommitted_prereg(self):  # MUT-FINAL-NOPREREG
        with TempRepo() as r:
            r.write('ledger.jsonl', '')
            r.commit('ledger.jsonl')
            p = self.prereg(r)  # written, digest matches, NOT committed
            with self.assertRaisesRegex(run.RunRefused, 'not committed'):
                run.check_final(p, self.cfg, self.probes, ledger=r.path('ledger.jsonl'))

    def test_final_refuses_edited_prereg(self):
        with TempRepo() as r:
            r.write('ledger.jsonl', '')
            p = self.prereg(r)
            r.commit('ledger.jsonl', 'preregistration.json')
            with open(p, 'a') as f:
                f.write('\n')
            with self.assertRaisesRegex(run.RunRefused, 'differs from its committed blob'):
                run.check_final(p, self.cfg, self.probes, ledger=r.path('ledger.jsonl'))

    def test_final_refuses_digest_drift(self):
        with TempRepo() as r:
            r.write('ledger.jsonl', '')
            p = self.prereg(r)
            r.commit('ledger.jsonl', 'preregistration.json')
            drift = FakeProbes(cli={'claude': '2.1.300 (Claude Code)', 'codex': 'codex-cli 0.159.2'})
            with self.assertRaisesRegex(run.RunRefused, 'AC5.1'):
                run.check_final(p, self.cfg, drift, ledger=r.path('ledger.jsonl'))
            knob = cfg(teaching='compact')
            with self.assertRaisesRegex(run.RunRefused, 'AC5.1'):
                run.check_final(p, knob, self.probes, ledger=r.path('ledger.jsonl'))

    def test_final_accepts_committed_matching_prereg(self):  # non-vacuous: the guard can pass
        with TempRepo() as r:
            r.write('ledger.jsonl', '')
            p = self.prereg(r)
            r.commit('ledger.jsonl', 'preregistration.json')
            pre = run.check_final(p, self.cfg, self.probes, ledger=r.path('ledger.jsonl'))
            self.assertEqual(pre['kind'], run.PREREG_KIND)

    def test_prereg_is_final_shaped(self):
        pre = run.prereg_draft(self.cfg, self.probes, ledger='/nonexistent/ledger.jsonl')
        live = pre['live_config']
        self.assertEqual(live['models'], {'claude': 'claude-opus-5-5', 'codex': 'gpt-6.1-sol'})
        self.assertEqual(live['N'], 3)
        self.assertEqual(sorted(live['corpus']['tasks']), CORE_IDS)
        self.assertEqual(pre['config_digest'], run.config_digest(live))
        self.assertEqual(pre['statistics']['thresholds']['delta_min'], '-1/50')
        self.assertIn('--max-budget-usd', live['argv_templates']['claude/world'])


class ConfigRules(unittest.TestCase):
    def test_forbidden_knobs_refused(self):
        base = run.load_config(run.CONFIG)
        base.pop('_path')
        for mut in (lambda c: c['final'].update(N=5), lambda c: c['final'].update(tier='smoke'),
                    lambda c: c.update(deadline_s=300), lambda c: c['knobs'].update(teaching='tiny'),
                    lambda c: c['knobs'].update(world_tools_wording='x'),
                    lambda c: c['smoke']['models'].pop('codex')):
            c = json.loads(json.dumps(base))
            mut(c)
            with tempfile.NamedTemporaryFile('w', suffix='.json', delete=False) as f:
                json.dump(c, f)
            try:
                with self.assertRaises(run.RunRefused):
                    run.load_config(f.name)
            finally:
                os.unlink(f.name)

    def test_smoke_subset_covers_the_p4_tasks_and_is_core(self):
        c = run.load_config(run.CONFIG)
        ids = c['smoke']['tasks']
        for t in ('pipeline', 'cli_args', 'api_call_json'):
            self.assertIn(t, ids)
        self.assertEqual([t.id for t in run.mode_tasks(c, 'smoke', FakeProbes())], sorted(ids))
        self.assertEqual(c['smoke']['models'], {'claude': 'claude-haiku-4-5', 'codex': 'gpt-6-luna'})

    def test_shell_digest_ignores_world_wording_world_digest_does_not(self):
        p = FakeProbes()
        a, b = cfg(world_tools_wording='m1', claude_isolation=False), cfg(world_tools_wording='row135', claude_isolation=False)
        for ag in arms.AGENTS:
            self.assertEqual(run.arm_digest(a, 'smoke', ag, 'shell', p), run.arm_digest(b, 'smoke', ag, 'shell', p))
            self.assertNotEqual(run.arm_digest(a, 'smoke', ag, 'world', p), run.arm_digest(b, 'smoke', ag, 'world', p))
        c = cfg(claude_isolation=True)
        self.assertNotEqual(run.arm_digest(a, 'smoke', 'claude', 'shell', p), run.arm_digest(c, 'smoke', 'claude', 'shell', p))
        self.assertEqual(run.arm_digest(a, 'smoke', 'codex', 'shell', p), run.arm_digest(c, 'smoke', 'codex', 'shell', p))


class Isolation(unittest.TestCase):  # knob (b): identical in both Claude arms
    def test_identical_additions_in_both_claude_arms(self):
        sh = arms.build_argv('claude', 'shell', model='m', prompt_text='p', cwd='/w', claude_isolation=True)
        wo = arms.build_argv('claude', 'world', model='m', prompt_text='p', cwd='/v', mcp_config_path='/x.json',
                             claude_isolation=True)
        for argv in (sh, wo):
            self.assertEqual(run_flag(argv, '--setting-sources'), '')
            s = json.loads(run_flag(argv, '--settings'))
            self.assertEqual({k: s[k] for k in ('enabledPlugins', 'autoMemoryEnabled')}, arms.CLAUDE_ISOLATION_SETTINGS)
        self.assertEqual(json.loads(run_flag(sh, '--settings'))['sandbox'], arms.CLAUDE_SANDBOX_SETTINGS['sandbox'])
        self.assertNotIn('sandbox', json.loads(run_flag(wo, '--settings')))
        base_sh = arms.build_argv('claude', 'shell', model='m', prompt_text='p', cwd='/w')
        self.assertNotIn('--setting-sources', base_sh)  # off by default: the goldens stay the §4.3 text

    def test_codex_unaffected_and_budget_claude_only(self):
        for arm in arms.ARMS:
            mcp = None
            a = arms.build_argv('codex', arm, model='m', prompt_text='p', cwd='/w', mcp_config_path=mcp)
            b = arms.build_argv('codex', arm, model='m', prompt_text='p', cwd='/w', mcp_config_path=mcp,
                                claude_isolation=True, max_budget_usd=1.0)
            self.assertEqual(a, b)
        c = arms.build_argv('claude', 'shell', model='m', prompt_text='p', cwd='/w', max_budget_usd=1)
        self.assertEqual(run_flag(c, '--max-budget-usd'), '1.00')
        with self.assertRaises(arms.ArmError):
            arms.build_argv('claude', 'shell', model='m', prompt_text='p', cwd='/w', max_budget_usd=0)


class CodexMcpApproval(unittest.TestCase):  # knob (c)
    def test_default_is_the_golden_and_approve_adds_one_override(self):
        d = arms.build_argv('codex', 'world', model='m', prompt_text='p', cwd='/v')
        a = arms.build_argv('codex', 'world', model='m', prompt_text='p', cwd='/v', codex_mcp_approval='approve')
        self.assertEqual(d, arms.build_argv('codex', 'world', model='m', prompt_text='p', cwd='/v',
                                            codex_mcp_approval='default'))
        added = [x for x in a if x not in d]
        self.assertEqual(added, ['mcp_servers.world.default_tools_approval_mode=approve'])
        self.assertEqual(a[a.index(added[0]) - 1], '-c')
        self.assertIn('shell_tool', a)  # still no shell
        sh = arms.build_argv('codex', 'shell', model='m', prompt_text='p', cwd='/w', codex_mcp_approval='approve')
        self.assertEqual(sh, arms.build_argv('codex', 'shell', model='m', prompt_text='p', cwd='/w'))
        with self.assertRaises(arms.ArmError):
            arms.build_argv('codex', 'world', model='m', prompt_text='p', cwd='/v', codex_mcp_approval='auto')

    def test_only_the_codex_world_cell_digest_moves(self):
        p = FakeProbes()
        a, b = cfg(codex_mcp_approval='default'), cfg(codex_mcp_approval='approve')
        for ag in arms.AGENTS:
            for x in arms.ARMS:
                same = run.arm_digest(a, 'smoke', ag, x, p) == run.arm_digest(b, 'smoke', ag, x, p)
                self.assertEqual(same, (ag, x) != ('codex', 'world'), (ag, x))


def run_flag(argv, flag):
    return argv[argv.index(flag) + 1]


class WorldWording(unittest.TestCase):  # knob (a)
    SPEC = {'id': 'cli_args', 'description': 'd', 'caps': ['IO', 'FS', 'Env'], 'task_prompt': 'Sum <LANG>.\n',
            'expected_stdout': '1\n'}

    def render(self, wording):
        return prompt.render_task(self.SPEC, 'world', template=TEMPLATE, solution_path='benchmark/solution.ail',
                                  deadline_s=600, world_wording_id=wording)

    def test_m1_reproduces_the_stale_note_and_row135_tells_the_truth(self):
        m1, new = self.render('m1'), self.render('row135')
        self.assertIn('cannot pipe stdin', m1)
        self.assertNotIn('cannot pipe stdin', new)
        for field in ('`stdin`', '`argv`', '`caps`'):
            self.assertIn(field, new)
        self.assertIn('"caps": ["IO", "FS", "Env"]', new)
        self.assertNotIn('"caps"', m1)
        for w in (m1, new):
            self.assertNotIn('ailang run', w)  # MUT-PROMPT-LEAK holds for both wordings

    def test_wording_change_stays_inside_the_tools_block(self):
        import test_prompt
        s = prompt.render_task(self.SPEC, 'shell', template=TEMPLATE, solution_path=test_prompt.SHELL_PATH,
                               deadline_s=600)
        for wording in prompt.WORLD_WORDINGS:
            w = self.render(wording)
            self.assertEqual(test_prompt.disallowed_diff(prompt.render_message('T', s),
                                                         prompt.render_message('T', w)), [])

    def test_unknown_wording_refused(self):
        with self.assertRaises(prompt.PromptError):
            self.render('nope')


class Outcomes(unittest.TestCase):
    def test_claude_success_is_no_agent_failure(self):
        self.assertIsNone(run.agent_failure('claude', 0, {'result': {'subtype': 'success', 'is_error': False}}, ''))

    def test_claude_missing_result_is_untyped_api_error(self):
        cat, ev, typed = run.agent_failure('claude', 1, {'result': None}, 'boom')
        self.assertEqual((cat, typed), ('api_error', False))
        self.assertIn('boom', ev)

    def test_claude_budget_and_turns(self):
        self.assertEqual(run.agent_failure('claude', 1, {'result': {'subtype': 'error_max_budget_usd'}}, '')[0], 'cost_killed')
        self.assertEqual(run.agent_failure('claude', 1, {'result': {'subtype': 'error_max_turns'}}, '')[0], 'step_exhausted')

    def test_rate_limit_and_quota_evidence(self):
        self.assertEqual(run.quota_evidence('claude', {'rate_limit_status': ['allowed', 'rejected']}, '')[0], 'rate_limit')
        self.assertIsNone(run.quota_evidence('claude', {'rate_limit_status': ['allowed', 'allowed_warning']}, ''))
        self.assertEqual(run.quota_evidence('codex', {'errors': ["You've hit your usage limit"]}, '')[0], 'quota_exhausted')
        self.assertEqual(run.quota_evidence('codex', {'errors': ['HTTP 429 Too Many Requests']}, '')[0], 'rate_limit')
        self.assertEqual(run.agent_failure('codex', 1, {'errors': ['429 Too Many Requests']}, '')[0], 'rate_limit')

    def test_codex_errors(self):
        cat, _, typed = run.agent_failure('codex', 1, {'errors': ['stream disconnected before completion']}, '')
        self.assertEqual((cat, typed), ('api_error', False))
        self.assertTrue(run.agent_failure('codex', 1, {'errors': ['503 Service Unavailable']}, '')[2])
        # measured in M4 iteration 5 (codex shell, gpt-6-luna): a provider capacity refusal is typed
        cap = 'Selected model is at capacity. Please try a different model.'
        self.assertEqual(run.agent_failure('codex', 1, {'errors': [cap]}, ''), ('api_error', cap, True))
        self.assertEqual(run.agent_failure('codex', 2, {'errors': []}, 'x')[0], 'api_error')
        self.assertIsNone(run.agent_failure('codex', 0, {'errors': []}, ''))

    def test_decide_category_precedence(self):
        g_pass = {'stdout_ok': True, 'compile_ok': True, 'runtime_ok': True, 'solution_state': 'present'}
        g_fail = {'stdout_ok': False, 'compile_ok': False, 'runtime_ok': False, 'solution_state': 'present'}
        self.assertEqual(run.decide_category(driver_category='timeout', agent_fail=None, graded=g_pass)[0], 'timeout')
        self.assertEqual(run.decide_category(driver_category=None, agent_fail=('api_error', 'e', False), graded=g_pass)[0], 'none')
        self.assertEqual(run.decide_category(driver_category=None, agent_fail=('api_error', 'e', False), graded=g_fail)[0], 'api_error')
        self.assertEqual(run.decide_category(driver_category=None, agent_fail=None, graded=g_fail)[0], 'compile_error')

    def test_world_api_error_untyped_is_fail_shell_is_fault(self):
        base = {'stdout_ok': False, 'error_category': 'api_error', 'cause_evidence': 'x', 'cause_typed': False}
        self.assertEqual(run.classify_row(dict(base, arm='world')), classify.FAIL)
        self.assertEqual(run.classify_row(dict(base, arm='shell')), classify.HARNESS_FAULT)
        r = dict(base, arm='shell', error_category='cost_killed')
        self.assertEqual(run.classify_row(r), 'UNCLASSIFIED')
        self.assertIn('pre-registered', r['classification_error'])

    def test_row153_counts(self):
        tr = '{"x":"host callback timed out after 10s"}\n{"y":"ok"}\n{"z":"host callback timed out"}\n'
        log = 'a\nexecute plan phase: context deadline exceeded\nb\n'
        h = run.row153_hits(tr, log)
        self.assertEqual(h['transcript_callback_timeouts'], 2)
        self.assertEqual(h['daemon_deadline_lines'], 1)
        self.assertEqual(run.row153_hits('{"ok":1}', '')['transcript_callback_timeouts'], 0)


class Summary(unittest.TestCase):
    def row(self, agent, arm, task, cls, wall):
        return {'agent': agent, 'arm': arm, 'task': task, 'class': cls, 'wall_ms': wall, 'finish_reason': 'exit'}

    def test_delta_and_median_overhead(self):
        P, F, H = classify.PASS, classify.FAIL, classify.HARNESS_FAULT
        rows = [self.row('claude', 'shell', 'a', P, 100), self.row('claude', 'world', 'a', P, 150),
                self.row('claude', 'shell', 'b', P, 100), self.row('claude', 'world', 'b', F, 300),
                self.row('claude', 'shell', 'c', H, 100), self.row('claude', 'world', 'c', P, 100)]
        s = run.summarize(rows, ['a', 'b', 'c'])['claude']
        self.assertEqual((s['pass_S'], s['pass_W']), (2, 2))
        self.assertEqual(s['delta'], 0.0)
        self.assertEqual(s['o_t'], {'a': 0.5, 'b': 2.0})  # c excluded: a harness fault
        self.assertEqual(s['O'], 1.25)
        self.assertEqual(s['harness_faults'], {'shell': 1, 'world': 0})


class Ledger(unittest.TestCase):
    def test_append_requires_fields_and_order(self):
        with tempfile.TemporaryDirectory() as d:
            led = os.path.join(d, 'l.jsonl')
            row = {k: None for k in run.LEDGER_FIELDS}
            row['iteration'] = 0
            run.append_ledger(row, led)
            with self.assertRaises(run.RunRefused):
                run.append_ledger(dict(row, iteration=0), led)
            bad = dict(row, iteration=1)
            bad.pop('commit')
            with self.assertRaises(run.RunRefused):
                run.append_ledger(bad, led)
            self.assertEqual(len(run.read_ledger(led)), 1)

    def test_committed_ledger_rows_carry_the_ac41_fields(self):
        for r in run.read_ledger(run.LEDGER):
            for k in run.LEDGER_FIELDS:
                self.assertIn(k, r)


if __name__ == '__main__':
    unittest.main()
