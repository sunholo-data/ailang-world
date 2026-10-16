"""Row 182 (w-episode-outcome-capture): the harness's world/task-outcome/v1 payload and its record call.

The payload is checked against the REAL `ailang-worldd outcome record` validator (built from this
tree), so the Python and Go sides cannot drift: a valid payload passes validation and reaches the
commit step; a broken one is refused before anything is sent (exit 1).
"""
import json
import os
import subprocess
import sys
import tempfile
import types
import unittest

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import run  # noqa: E402

REPO = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..'))
SPEC = {'id': 'config_file_parser', 'caps': ['IO', 'FS'], 'expected_stdout': 'Loaded MyApp v1.0.0 on port 8080 with 3 features\n'}


def row(arm='world', cls='FAIL'):
    return {'agent': 'codex', 'arm': arm, 'run': 2, 'label': 'final-paired-r2', 'task': 'config_file_parser',
            'class': cls, 'error_category': 'logic_error', 'compile_ok': True, 'runtime_ok': True, 'stdout_ok': False,
            'grade_stdout_head': 'configuration file not found\n', 'grade_stderr_head': '',
            'solution_sha256': 'a' * 64, 'solution_state': 'present'}


DRV_WORLD = {'world': {'episode': 'fl-cx-r2', 'log_from': 259, 'log_to': 261}}
GRADED = {'tool_sha256': 'b' * 64}


class Payload(unittest.TestCase):
    def test_world_row(self):
        p = run.outcome_payload('fl-cx-r2', row(), SPEC, GRADED, DRV_WORLD, '/abs/ailang')
        self.assertEqual(p['entries'], {'from': 259, 'to': 261})
        self.assertEqual(p['grading']['argv'][0], 'run')  # the tool path is not part of the rule
        self.assertNotIn('/abs/ailang', p['grading']['argv'])
        self.assertEqual(p['grading']['argv'][-1], 'benchmark/solution.ail')
        self.assertIn('--caps', p['grading']['argv'])
        self.assertEqual(p['grading']['rule'], run.OUTCOME_RULE)
        self.assertEqual(p['grading']['expectedStdout'], SPEC['expected_stdout'])
        self.assertEqual(p['outcome']['class'], 'FAIL')
        self.assertEqual(p['outcome']['stdoutHead'], 'configuration file not found\n')
        self.assertEqual(p['episodeId'], 'fl-cx-r2')

    def test_shell_row_has_no_entries(self):
        p = run.outcome_payload('fl-cx-r2', row('shell', 'PASS'), SPEC, GRADED, {}, '/abs/ailang')
        self.assertIsNone(p['entries'])

    def test_world_row_without_range_is_refused(self):
        with self.assertRaises(ValueError):
            run.outcome_payload('fl-cx-r2', row(), SPEC, GRADED, {'world': {'log_from': 259, 'log_to': None}}, '/a')

    def test_input_files_are_named(self):
        spec = dict(SPEC, input_files={'b.txt': 'x', 'a.csv': 'y'})
        p = run.outcome_payload('fl-cx-r2', row(), spec, GRADED, DRV_WORLD, '/a')
        self.assertEqual(p['grading']['inputFiles'], ['a.csv', 'b.txt'])


class Record(unittest.TestCase):
    L = {'bin': '/rig/bin', 'private': '/rig/private'}

    def test_success_parses_ref_and_entry(self):
        seen = {}

        def runner(argv, input, capture_output, env, timeout):
            seen.update(argv=argv, input=json.loads(input), env=env)
            ref = 'sha256:' + 'c' * 64
            return subprocess.CompletedProcess(argv, 0, f'recorded world/task-outcome/v1 codex/world/config_file_parser run 2: object {ref} entry 262\n'.encode(), b'')

        os.environ['AILANG_REGISTRY_API_KEY'] = 'must-not-leak'
        try:
            got = run.record_outcome(self.L, '127.0.0.1:7644', 'fl-cx-r2', row(), SPEC, GRADED, DRV_WORLD, '/a', runner=runner)
        finally:
            del os.environ['AILANG_REGISTRY_API_KEY']
        self.assertEqual(got, {'ref': 'sha256:' + 'c' * 64, 'entry': 262})
        self.assertEqual(seen['argv'][:5], ['/rig/bin/ailang-worldd', '--addr', 'http://127.0.0.1:7644', 'outcome', 'record'])
        self.assertIn('/rig/private/fl-cx-r2.session', seen['argv'])
        self.assertNotIn('AILANG_REGISTRY_API_KEY', seen['env'])
        self.assertEqual(seen['input']['task'], 'config_file_parser')

    def test_failure_is_an_error_never_a_raise(self):
        def bad(argv, **kw):
            return subprocess.CompletedProcess(argv, 2, b'', b'outcome writer: commit refused HTTP 401 Unauthorized; not retried')

        def boom(argv, **kw):
            raise OSError('no such file')

        def ok_rc_no_ref(argv, **kw):
            return subprocess.CompletedProcess(argv, 0, b'something else', b'')

        def ref_but_nonzero(argv, **kw):
            return subprocess.CompletedProcess(argv, 3, f'object sha256:{"c" * 64} entry 5'.encode(), b'uncertain')

        for runner in (bad, boom, ok_rc_no_ref, ref_but_nonzero):
            got = run.record_outcome(self.L, 'x', 'fl-cx-r2', row(), SPEC, GRADED, DRV_WORLD, '/a', runner=runner)
            self.assertIn('error', got)
        self.assertIn('401', run.record_outcome(self.L, 'x', 'fl-cx-r2', row(), SPEC, GRADED, DRV_WORLD, '/a', runner=bad)['error'])
        got = run.record_outcome(self.L, 'x', 'fl-cx-r2', row(), SPEC, GRADED, {'world': {}}, '/a', runner=bad)
        self.assertTrue(got['error'].startswith('payload:'))

    def test_error_count(self):
        rows = [{'outcome_record': {'error': 'x'}}, {'outcome_record': {'ref': 'r', 'entry': 1}}, {}]
        self.assertEqual(run.outcome_record_errors(rows), 1)


class Hook(unittest.TestCase):
    """maybe_record_outcome: the run_cell hook (AC5; MU-3, MU-5)."""

    def ctx(self, knob=True, client=True, token=True):
        return types.SimpleNamespace(cfg={'knobs': {'outcome_capture': knob}, 'tool_bin': '/abs/ailang'},
                                     L={'bin': '/rig/bin', 'private': '/rig/private'},
                                     client=types.SimpleNamespace(addr='127.0.0.1:7644') if client else None,
                                     tokens={'codex': 'tok'} if token else {})

    def test_records_both_arms_after_classify(self):
        seen = []

        def rec(L, addr, ep, r, spec, graded, drv, tool):
            seen.append((r['arm'], r.get('class')))
            return {'ref': 'sha256:' + 'c' * 64, 'entry': 9}

        for arm in ('shell', 'world'):
            r = row(arm, 'FAIL')
            run.maybe_record_outcome(self.ctx(), 'codex', 'fl-cx-r2', r, SPEC, GRADED, DRV_WORLD, recorder=rec)
            self.assertEqual(r['outcome_record'], {'ref': 'sha256:' + 'c' * 64, 'entry': 9})
        self.assertEqual(seen, [('shell', 'FAIL'), ('world', 'FAIL')])  # class is set before the record (MU-5)

    def test_off_or_no_client_or_no_token_records_nothing(self):
        def rec(*a):
            raise AssertionError('recorder must not be called')

        for c in (self.ctx(knob=False), self.ctx(client=False), self.ctx(token=False)):
            r = row()
            run.maybe_record_outcome(c, 'codex', 'fl-cx-r2', r, SPEC, GRADED, DRV_WORLD, recorder=rec)
            self.assertNotIn('outcome_record', r)
        c = self.ctx()
        del c.cfg['knobs']['outcome_capture']  # absent knob = off
        r = row()
        run.maybe_record_outcome(c, 'codex', 'fl-cx-r2', r, SPEC, GRADED, DRV_WORLD, recorder=rec)
        self.assertNotIn('outcome_record', r)

    def test_failure_never_changes_class_and_is_counted(self):
        def err(*a):
            return {'error': 'rc=2: refused'}

        def boom(*a):
            raise RuntimeError('unexpected')

        rows = []
        for rec in (err, boom):
            r = row('world', 'PASS')
            run.maybe_record_outcome(self.ctx(), 'codex', 'fl-cx-r2', r, SPEC, GRADED, DRV_WORLD, recorder=rec)
            self.assertEqual(r['class'], 'PASS')  # MU-3
            self.assertIn('error', r['outcome_record'])
            rows.append(r)
        self.assertEqual(run.outcome_record_errors(rows), 2)

    def test_run_cell_calls_the_hook_after_classify(self):
        import inspect
        src = inspect.getsource(run.run_cell)
        self.assertLess(src.index("row['class'] = classify_row(row)"), src.index('maybe_record_outcome('))


class HeadBytes(unittest.TestCase):
    def test_multibyte_head_fits_the_byte_bound(self):
        r = row()
        r['grade_stdout_head'] = 'é' * 400  # 800 UTF-8 bytes, 400 characters
        p = run.outcome_payload('fl-cx-r2', r, SPEC, GRADED, DRV_WORLD, '/a')
        self.assertLessEqual(len(p['outcome']['stdoutHead'].encode()), 400)
        self.assertEqual(p['outcome']['stdoutHead'], 'é' * 200)

    def test_surrogate_escape_is_bounded(self):
        r = row()
        r['grade_stdout_head'] = b'\xff\xfe'.decode('utf-8', errors='surrogateescape') * 200
        p = run.outcome_payload('fl-cx-r2', r, SPEC, GRADED, DRV_WORLD, '/a')
        self.assertLessEqual(len(p['outcome']['stdoutHead'].encode()), 400)


class KnobValidation(unittest.TestCase):
    def test_outcome_capture_must_be_bool(self):
        with open(run.CONFIG) as f:
            cfg = json.load(f)
        self.assertIs(cfg['knobs']['outcome_capture'], True)
        cfg['knobs']['outcome_capture'] = 'yes'
        with tempfile.NamedTemporaryFile('w', suffix='.json', delete=False) as f:
            json.dump(cfg, f)
        try:
            with self.assertRaises(run.RunRefused):
                run.load_config(f.name)
        finally:
            os.unlink(f.name)


class AgainstTheRealValidator(unittest.TestCase):
    """The Python payload passes the Go validator; a broken one is refused before any send."""

    @classmethod
    def setUpClass(cls):
        cls.tmp = tempfile.mkdtemp(prefix='outcome-validator-')
        cls.bin = os.path.join(cls.tmp, 'ailang-worldd')
        subprocess.run(['go', 'build', '-o', cls.bin, './cmd/ailang-worldd'], cwd=REPO, check=True, capture_output=True)
        cls.session = os.path.join(cls.tmp, 's.session')
        with open(cls.session, 'w') as f:
            f.write('1' * 64)
        os.chmod(cls.session, 0o600)

    def record(self, payload):
        env = {k: v for k, v in os.environ.items() if k not in ('AILANG_REGISTRY_API_KEY', 'WORLD_SESSION')}
        # Port 9 (discard) on loopback: nothing answers, so a payload that passes validation fails
        # at the commit step (transport: uncertain, exit 3) — never with the validator's exit 1.
        return subprocess.run([self.bin, '--addr', 'http://127.0.0.1:9', 'outcome', 'record', '--file', '-',
                               '--session', self.session], input=json.dumps(payload).encode(),
                              capture_output=True, env=env, timeout=60)

    def test_valid_world_and_shell_payloads_pass_validation(self):
        for r, drv in ((row(), DRV_WORLD), (row('shell', 'PASS'), {})):
            p = self.record(run.outcome_payload('fl-cx-r2', r, SPEC, GRADED, drv, '/abs/ailang'))
            self.assertNotEqual(p.returncode, 1, p.stderr.decode())
            self.assertNotIn(b'refused:', p.stderr)

    def test_multibyte_head_passes_the_real_validator(self):
        r = row()
        r['grade_stdout_head'] = 'é' * 400
        p = self.record(run.outcome_payload('fl-cx-r2', r, SPEC, GRADED, DRV_WORLD, '/abs/ailang'))
        self.assertNotEqual(p.returncode, 1, p.stderr.decode())

    def test_broken_payload_is_refused_before_send(self):
        bad = run.outcome_payload('fl-cx-r2', row(), SPEC, GRADED, DRV_WORLD, '/abs/ailang')
        bad['entries'] = None
        p = self.record(bad)
        self.assertEqual(p.returncode, 1)
        self.assertIn(b'refused:', p.stderr)


if __name__ == '__main__':
    unittest.main()
