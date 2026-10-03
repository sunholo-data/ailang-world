#!/usr/bin/env python3
"""AC1.3 — the grader port: CompareOutput unit table, argv, workspace, stdin, binary pinning.

Hermetic arms drive a FAKE tool binary (a tiny interpreter for directive "solutions"), so the
run plumbing — stdin, argv, input files, cwd, timeout, compile/runtime split, the executed
binary's digest — is tested without ailang. The REAL-binary arms run the pinned tool
($FLOOR_TOOL_BIN, default ~/.pinned-ailang-tools/v0.51.0/ailang) on banked solutions; they skip
when it is absent unless FLOOR_REQUIRE_TOOL=1, which turns the absence red (CI sets it).

Mutants killed here: MUT-GRADE-LOOSE (near-miss), MUT-NO-STDIN (stdin arms), MUT-PATH-AILANG
(executed-binary digest + impostor on PATH).
"""
import json
import os
import shutil
import stat
import sys
import tempfile
import unittest
import urllib.request

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import grade  # noqa: E402

FAKE_TOOL = r'''#!/usr/bin/env python3
import os, sys, time
argv = sys.argv[1:]
src = open(os.path.join(os.getcwd(), 'benchmark', 'solution.ail')).read()
rest = argv[argv.index('--') + 1:] if '--' in argv else []
for line in src.splitlines()[1:]:
    cmd, _, arg = line.partition(' ')
    if cmd == 'PRINT': sys.stdout.write(arg.encode().decode('unicode_escape'))
    elif cmd == 'ECHO_STDIN': sys.stdout.write(sys.stdin.read())
    elif cmd == 'ARGS': print(' '.join(rest))
    elif cmd == 'ARGV': print(' '.join(argv))
    elif cmd == 'CAT': sys.stdout.write(open(arg).read())
    elif cmd == 'HAS_STD': print(os.path.exists('std'))
    elif cmd == 'STDERR': sys.stderr.write(arg + '\n')
    elif cmd == 'EXIT': sys.exit(int(arg))
    elif cmd == 'SLEEP': time.sleep(float(arg))
    elif cmd == 'WHOAMI': print('fake-tool')
'''
IMPOSTOR = '#!/bin/sh\necho PATH-AILANG\n'


def _write_exe(path, body):
    with open(path, 'w') as f:
        f.write(body)
    os.chmod(path, os.stat(path).st_mode | stat.S_IXUSR | stat.S_IXGRP | stat.S_IXOTH)


def sol(*lines):
    return ('module benchmark/solution\n' + '\n'.join(lines) + '\n').encode()


class CompareOutputTable(unittest.TestCase):
    """AC1.3: exact, JSON-equal and normalized pass; near-misses and extra lines fail."""

    PASS = [
        ('exact', 'Sum: 30\nDone', 'Sum: 30\nDone'),
        ('outer whitespace', '2\n4\n6\n', '\n2\n4\n6   \n\n'),
        ('json spacing+order', '{"a":1,"b":[1,2]}', '{ "b": [1, 2], "a": 1 }'),
        ('json int vs float', '{"n":1}', '{"n":1.0}'),
        ('json big ints collapse like float64', '12345678901234567890', '12345678901234567891'),
        ('normalized bool case', 'valid: True', 'valid: true'),
        ('normalized decimals', 'total 7.50\nx 16.0', 'total 7.5\nx 16'),
        ('normalized punct space', '[1, 2, 3]', '[1,2,3]'),
        ('normalized -0.0', 'v -0.0', 'v -0'),
        ('huge decimal via FormatFloat f', 'n 123456789012345678901234567890.0',
         'n 123456789012345680000000000000'),
        ('go trims \\v', 'a', 'a\v'),
    ]
    FAIL = [
        ('near-miss digit', '200', '2000'),
        ('near-miss value', 'Sum: 30', 'Sum: 31'),
        ('space not adjacent to punct is kept', '1 2', '12'),
        ('extra line', '2\n4', '2\n4\n6'),
        ('missing line', '2\n4\n6', '2\n4'),
        ('substring is not a match', 'ok', 'ok: CREDIT(50)\nok'),
        ('json bool is not a number', 'true', '1'),
        ('json false is not 0', '[false]', '[0]'),
        ('set vs list', '[1,2]', '{1,2}'),
        ('go does not trim \\x1c', 'a', 'a\x1c'),
        ('RE2 \\s has no \\v', '[1,2]', '[1,\v2]'),
        ('NaN is not JSON', '[NaN]', '[ NaN ,1]'),
    ]

    def test_pass(self):
        for name, e, a in self.PASS:
            with self.subTest(name):
                self.assertTrue(grade.compare_output(e, a), (e, a))

    def test_fail(self):
        for name, e, a in self.FAIL:
            with self.subTest(name):
                self.assertFalse(grade.compare_output(e, a), (e, a))

    def test_canon_decimal_matches_go_formatfloat(self):
        for s, want in [('7.50', '7.5'), ('16.0', '16'), ('3.14', '3.14'), ('0.1', '0.1'),
                        ('-0.0', '-0'), ('10000000000000000.0', '10000000000000000'),
                        ('0.000001', '0.000001'), ('1' * 400 + '.0', '1' * 400 + '.0')]:
            with self.subTest(s):
                self.assertEqual(grade.canon_decimal(s), want)


class Argv(unittest.TestCase):
    def test_v9_argv(self):
        t = '/abs/ailang'
        self.assertEqual(grade.build_argv(t, {'caps': ['IO']}),
                         [t, 'run', '--entry', 'main', '--quiet', '--relax-modules', '--caps', 'IO',
                          'benchmark/solution.ail'])
        self.assertEqual(grade.build_argv(t, {'caps': ['IO', 'FS', 'Env'], 'cli_args': ['numbers.txt']})[-3:],
                         ['benchmark/solution.ail', '--', 'numbers.txt'])
        a = grade.build_argv(t, {'caps': ['Net', 'IO'], 'net_allow_localhost': True})
        self.assertEqual(a[6:], ['--caps', 'Net,IO', '--net-allow-http', '--net-allow-localhost',
                                 'benchmark/solution.ail'])
        self.assertIn('--ai-stub', grade.build_argv(t, {'caps': ['ai', 'IO']}))
        self.assertNotIn('--stdlib-path', grade.build_argv(t, {'caps': ['IO']}))

    def test_tool_must_be_explicit(self):
        for bad in ('ailang', 'bin/ailang', '', '/nonexistent/ailang'):
            with self.subTest(bad):
                with self.assertRaises(grade.GradeSetupError):
                    grade.resolve_tool(bad)


class FakeToolGrading(unittest.TestCase):
    def setUp(self):
        self.td = tempfile.mkdtemp(prefix='floor-test-')
        self.tool = os.path.join(self.td, 'tool', 'ailang-v0.51.0')
        os.makedirs(os.path.dirname(self.tool))
        _write_exe(self.tool, FAKE_TOOL)
        # an impostor `ailang` FIRST on PATH: a grader that drifts to PATH runs it (MUT-PATH-AILANG)
        self.pathdir = os.path.join(self.td, 'pathbin')
        os.makedirs(self.pathdir)
        _write_exe(os.path.join(self.pathdir, 'ailang'), IMPOSTOR)
        self.old_path = os.environ['PATH']
        os.environ['PATH'] = self.pathdir + os.pathsep + self.old_path

    def tearDown(self):
        os.environ['PATH'] = self.old_path
        shutil.rmtree(self.td, ignore_errors=True)

    def g(self, solution, spec, **kw):
        spec = dict({'caps': ['IO']}, **spec)
        return grade.grade(solution, spec, self.tool, **kw)

    def test_pins_the_explicit_binary_not_path(self):
        r = self.g(sol('WHOAMI'), {'expected_stdout': 'fake-tool\n'})
        self.assertTrue(r['stdout_ok'], r)
        self.assertEqual(r['executed_bin'], self.tool)
        self.assertEqual(r['tool_sha256'], grade.sha256_file(self.tool))

    def test_stdin_is_piped(self):
        r = self.g(sol('ECHO_STDIN'), {'stdin': '1\n2\n3\n', 'expected_stdout': '1\n2\n3\n'})
        self.assertTrue(r['stdout_ok'], r)

    def test_wrong_stdin_fails(self):
        r = self.g(sol('ECHO_STDIN'), {'stdin': '1\n2\n', 'expected_stdout': '1\n2\n3\n'})
        self.assertTrue(r['runtime_ok'])
        self.assertFalse(r['stdout_ok'])

    def test_no_stdin_is_devnull(self):
        r = self.g(sol('ECHO_STDIN', 'PRINT end'), {'expected_stdout': 'end'})
        self.assertTrue(r['stdout_ok'], r)

    def test_input_files_and_cli_args(self):
        spec = {'input_files': {'numbers.txt': '1\n2\n', 'sub/x.txt': 'X'}, 'cli_args': ['numbers.txt'],
                'caps': ['IO', 'FS', 'Env'], 'expected_stdout': 'numbers.txt\n1\n2\nX'}
        r = self.g(sol('ARGS', 'CAT numbers.txt', 'CAT sub/x.txt'), spec)
        self.assertTrue(r['stdout_ok'], r)
        self.assertEqual(r['argv'][-2:], ['--', 'numbers.txt'])

    def test_fresh_workspace_has_no_std(self):
        r = self.g(sol('HAS_STD'), {'expected_stdout': 'False'})
        self.assertTrue(r['stdout_ok'], r)

    def test_compile_vs_runtime_split(self):
        r = self.g(sol('STDERR Error: type error in benchmark/solution', 'EXIT 1'), {'expected_stdout': 'x'})
        self.assertEqual((r['compile_ok'], r['runtime_ok'], r['stdout_ok']), (False, False, False))
        self.assertEqual(grade.error_category(r), 'compile_error')
        r = self.g(sol('PRINT x', 'STDERR boom', 'EXIT 3'), {'expected_stdout': 'x'})
        self.assertEqual((r['compile_ok'], r['runtime_ok'], r['stdout_ok']), (True, False, False))
        self.assertEqual(grade.error_category(r), 'runtime_error')
        r = self.g(sol('PRINT y'), {'expected_stdout': 'x'})
        self.assertEqual(grade.error_category(r), 'logic_error')

    def test_timeout(self):
        r = self.g(sol('PRINT x', 'SLEEP 5'), {'expected_stdout': 'x'}, timeout=0.5)
        self.assertTrue(r['timed_out'])
        self.assertEqual((r['compile_ok'], r['runtime_ok'], r['stdout_ok']), (True, False, False))

    def test_missing_empty_placeholder_never_run_and_are_logic_errors(self):
        for s, state in ((None, 'missing'), (b'', 'empty'), (grade.GO_PLACEHOLDER.encode(), 'placeholder')):
            with self.subTest(state):
                r = self.g(s, {'expected_stdout': 'x'})
                self.assertEqual(r['solution_state'], state)
                self.assertIsNone(r['argv'])
                self.assertEqual(grade.error_category(r), 'logic_error')

    def test_unported_modes_refused(self):
        for extra in ({'grading': 'quine'}, {'grade_entrypoint': 'x'}):
            with self.assertRaises(grade.GradeSetupError):
                self.g(sol('PRINT x'), dict({'expected_stdout': 'x'}, **extra))


class HTTPMockShape(unittest.TestCase):
    def test_echo_like_go_mock(self):
        with grade.HTTPMock() as m:
            req = urllib.request.Request(m.url + '/post', data=b'{"count":42,"message":"hi"}', method='POST',
                                         headers={'x-test-header': 'value123', 'Content-Type': 'application/json'})
            with urllib.request.urlopen(req, timeout=5) as resp:
                self.assertEqual(resp.status, 200)
                body = resp.read()
        self.assertTrue(body.endswith(b'\n'))
        d = json.loads(body)
        self.assertEqual(list(d), ['headers', 'json', 'method', 'url'])
        self.assertEqual(d['json'], {'count': 42, 'message': 'hi'})
        self.assertEqual(d['headers']['X-Test-Header'], 'value123')
        self.assertNotIn('Host', d['headers'])
        self.assertEqual((d['method'], d['url']), ('POST', '/post'))
        self.assertEqual(len(m.hits), 1)


# ---------------------------------------------------------------- real pinned binary

TOOL = os.environ.get('FLOOR_TOOL_BIN', grade.DEFAULT_TOOL_BIN)
REQUIRE = os.environ.get('FLOOR_REQUIRE_TOOL') == '1'
HAVE_TOOL = os.path.isfile(TOOL) and os.access(TOOL, os.X_OK)
TD = os.path.join(HERE, 'testdata')
# the three V8 tasks as pinned at 76a5aef (test_corpus checks these against the real YAMLs)
SPEC_PIPELINE = {'id': 'pipeline', 'caps': ['IO'], 'expected_stdout': '2\n4\n6\n8\n10\n', 'stdin': '1\n2\n3\n4\n5\n'}
SPEC_CLI_ARGS = {'id': 'cli_args', 'caps': ['IO', 'FS', 'Env'], 'expected_stdout': '15\n',
                 'input_files': {'numbers.txt': '1\n2\n3\n4\n5\n'}, 'cli_args': ['numbers.txt']}
SPEC_API = {'id': 'api_call_json', 'caps': ['Net', 'IO'], 'expected_stdout': '200\n', 'net_allow_localhost': True}


class RealBinary(unittest.TestCase):
    def setUp(self):
        if not HAVE_TOOL:
            if REQUIRE:
                self.fail(f'FLOOR_REQUIRE_TOOL=1 but the pinned tool binary is absent: {TOOL}')
            self.skipTest(f'pinned tool binary absent: {TOOL}')

    def read(self, name):
        with open(os.path.join(TD, name), 'rb') as f:
            return f.read()

    def test_pipeline_stdin(self):
        s = self.read('pipeline_solution.ail.txt')
        self.assertTrue(grade.grade(s, SPEC_PIPELINE, TOOL)['stdout_ok'])
        wrong = dict(SPEC_PIPELINE, stdin='1\n2\n3\n4\n6\n')
        r = grade.grade(s, wrong, TOOL)
        self.assertTrue(r['runtime_ok'], r)
        self.assertFalse(r['stdout_ok'])

    def test_cli_args_file_and_argv(self):
        s = self.read('cli_args_solution.ail.txt')
        self.assertTrue(grade.grade(s, SPEC_CLI_ARGS, TOOL)['stdout_ok'])
        self.assertFalse(grade.grade(s, dict(SPEC_CLI_ARGS, input_files={'numbers.txt': '1\n2\n'}), TOOL)['stdout_ok'])

    def test_api_call_json_needs_the_live_mock(self):
        tmpl = self.read('api_call_json_solution.ail.txt').decode()
        with grade.HTTPMock() as m:
            r = grade.grade(tmpl.replace('{{MOCK_HTTP_URL}}', m.url).encode(), SPEC_API, TOOL)
            self.assertTrue(r['stdout_ok'], r)
            self.assertEqual(len(m.hits), 1)
            self.assertEqual(m.hits[0]['headers'].get('X-Test-Header'), 'value123')
            port = m.server.server_address[1]
        r = grade.grade(tmpl.replace('{{MOCK_HTTP_URL}}', f'http://127.0.0.1:{port}').encode(), SPEC_API, TOOL)
        self.assertFalse(r['stdout_ok'])  # mock down -> no 200

    def test_records_pinned_digest(self):
        r = grade.grade(self.read('pipeline_solution.ail.txt'), SPEC_PIPELINE, TOOL)
        self.assertEqual(r['executed_bin'], TOOL)
        self.assertEqual(r['tool_sha256'], grade.sha256_file(TOOL))


if __name__ == '__main__':
    unittest.main()
