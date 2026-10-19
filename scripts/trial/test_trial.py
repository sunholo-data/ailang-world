#!/usr/bin/env python3
"""Row 183 trial runner — hermetic: FAKE ``claude``/``codex`` on PATH, a fake World client, fake
``ailang-worldd`` runner; plus the REAL Go ``outcome record`` validator built from this tree (as
``scripts/floor/test_outcome.py`` does), so the trial payload and the verb cannot drift.
"""
import json
import os
import subprocess
import sys
import tempfile
import time
import unittest
from unittest import mock

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import trial  # noqa: E402
import arms  # noqa: E402  (scripts/floor, via trial's path insert)
import fakeworld  # noqa: E402

REPO = os.path.abspath(os.path.join(HERE, '..', '..'))
TOKEN = 'fedcba9876543210' * 4
EIGHT = [n.split('__', 2)[2] for n in arms.world_allowed_tools()]
NINE = EIGHT + ['workspace-exec']

FAKE_AGENT = r'''#!/usr/bin/env python3
import json, os, sys
argv = sys.argv[1:]
agent = os.path.basename(sys.argv[0])
spec = json.load(open(os.environ['FAKE_SPEC']))
facts = {'argv': argv, 'cwd': os.getcwd(), 'token_env': os.environ.get('WORLD_SESSION'),
         'registry_key': 'AILANG_REGISTRY_API_KEY' in os.environ}
if '--mcp-config' in argv:
    p = argv[argv.index('--mcp-config') + 1]
    # the World arm passes a 0600 file; the shell arm passes inline JSON (no servers)
    facts['mcp_mode'] = oct(os.stat(p).st_mode & 0o777) if os.path.exists(p) else 'inline:' + p
json.dump(facts, open(os.path.join(spec['marker_dir'], 'agent.facts'), 'w'))
open(os.path.join(spec['marker_dir'], 'agent.ran'), 'w').write('1')
for ev in spec['events']:
    print(json.dumps(ev), flush=True)
sys.exit(spec.get('rc', 0))
'''

CLAUDE_EVENTS = [
    {'type': 'system', 'subtype': 'init', 'tools': ['mcp__world__workspace-exec'], 'model': 'm'},
    {'type': 'assistant', 'message': {'content': [{'type': 'tool_use', 'name': 'mcp__world__workspace-exec', 'input': {}}]}},
    {'type': 'result', 'subtype': 'success', 'is_error': False, 'num_turns': 2, 'result': 'Added tests 1-4.'},
]
CODEX_EVENTS = [
    {'type': 'item.completed', 'item': {'id': 'i1', 'type': 'mcp_tool_call', 'server': 'world', 'tool': 'workspace-exec', 'status': 'completed'}},
    {'type': 'item.completed', 'item': {'id': 'i2', 'type': 'agent_message', 'text': 'Added tests 5-8.'}},
    {'type': 'turn.completed', 'usage': {}},
]

# The grader: records whether the agent had already run, prints what the test asks for.
GRADER = r'''import os, sys
marker, mode = sys.argv[1], sys.argv[2]
open(os.path.join(os.path.dirname(marker), 'grader.%s' % mode), 'w').write(
    'after-agent' if os.path.exists(marker) else 'before-agent')
if mode == 'multibyte':
    sys.stdout.write('é' * 1000)
    sys.stderr.buffer.write(b'\xff\xfe' * 500)
elif mode == 'sleep':
    import time; time.sleep(30)
elif mode == 'env':
    print('WORLD_SESSION' in os.environ, 'AILANG_REGISTRY_API_KEY' in os.environ)
sys.exit(0 if mode in ('pass', 'env') else 3)
'''


REPORT = {'episode': 'tps-a', 'entries': [11, 14], 'total_calls': 4,
          'calls_by_tool': {'ailang-read': 2, 'workspace-exec': 2},
          'plan_refusals': {'count': 1, 'by_tool': {'ailang-read': 1}, 'first': []},
          'handler_refusals': {'count': 0, 'by_tool': {}, 'first': []}}


class FakeClient:
    """head_index answers by whether the agent has run, so before/after is observable."""

    def __init__(self, marker, before=10, after=14, tools=NINE):
        self.marker, self.before, self.after, self.tools = marker, before, after, tools
        self.calls = []

    def health(self):
        return 200, {'status': 'ok', 'daemon_version': 'fake', 'db_path': '/rig/world.db'}

    def mcp(self, token, method, params=None, rpc_id=1):
        self.calls.append(('mcp', method, token == TOKEN))
        return {'result': {'tools': [{'name': t} for t in self.tools]}}

    def head_index(self):
        ran = os.path.exists(self.marker)
        self.calls.append(('head', ran))
        return self.after if ran else self.before


class FakeWorldd:
    """Stands in for <rig>/bin/ailang-worldd: records every call, answers outcome/provenance."""

    def __init__(self, outcome_rc=0, raise_on=None, report_rc=0):
        self.calls = []
        self.outcome_rc, self.raise_on, self.report_rc = outcome_rc, raise_on, report_rc
        self.n = 0

    def __call__(self, argv, input=None, capture_output=True, env=None, timeout=None):
        verb = argv[3]
        self.calls.append({'argv': argv, 'verb': verb, 'payload': json.loads(input) if input else None, 'env': env})
        if self.raise_on == verb:
            raise OSError('no such file')
        if verb == 'outcome':
            self.n += 1
            if self.outcome_rc:
                return subprocess.CompletedProcess(argv, self.outcome_rc, b'', b'commit refused HTTP 401')
            return subprocess.CompletedProcess(argv, 0, (f'recorded world/task-outcome/v1 x: object sha256:'
                                                         f'{"%064x" % self.n} entry {100 + self.n}\n').encode(), b'')
        if verb == 'episode':  # episode report <ep> --json (row 198)
            if self.report_rc:
                return subprocess.CompletedProcess(argv, self.report_rc, b'', b'episode has no committed calls')
            return subprocess.CompletedProcess(argv, 0, (json.dumps(dict(REPORT, episode=argv[5])) + '\n').encode(), b'')
        ep = argv[argv.index('--episode') + 1]
        since = argv[argv.index('--since') + 1]
        return subprocess.CompletedProcess(argv, 0, f'World-Provenance: store=/rig/world.db episode={ep} entries={since}-14\n'.encode(), b'')


def git(cwd, *a):
    subprocess.run(['git', '-C', cwd, *a], check=True, capture_output=True,
                   env=dict(os.environ, GIT_AUTHOR_NAME='t', GIT_AUTHOR_EMAIL='t@x', GIT_COMMITTER_NAME='t',
                            GIT_COMMITTER_EMAIL='t@x'))


class Base(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        t = self.tmp.name
        self.fakebin = os.path.join(t, 'fakebin')
        os.makedirs(self.fakebin)
        for name in ('claude', 'codex'):
            p = os.path.join(self.fakebin, name)
            with open(p, 'w') as f:
                f.write(FAKE_AGENT)
            os.chmod(p, 0o755)
        self.rig = os.path.join(t, 'rig')
        self.private = os.path.join(self.rig, 'private')
        os.makedirs(self.private, mode=0o700)
        os.chmod(self.private, 0o700)
        self.markers = os.path.join(t, 'markers')
        os.makedirs(self.markers)
        self.marker = os.path.join(self.markers, 'agent.ran')
        self.wt = os.path.join(t, 'repo-wt')
        os.makedirs(os.path.join(self.wt, 'backend'))
        with open(os.path.join(self.wt, 'backend', 'a.py'), 'w') as f:
            f.write('x = 1\n')
        git(self.wt, 'init', '-q')
        git(self.wt, 'add', '.')
        git(self.wt, 'commit', '-qm', 'base')
        self.grader = os.path.join(t, 'grader.py')
        with open(self.grader, 'w') as f:
            f.write(GRADER)
        with open(os.path.join(t, 'task.md'), 'w') as f:
            f.write('# Issue 67\n\n`--deep` prints a stale skipped warning.\n')
        self.env = mock.patch.dict(os.environ, {'PATH': self.fakebin + os.pathsep + os.environ.get('PATH', ''),
                                                'WORLD_SESSION': 'leaked-from-parent',
                                                'AILANG_REGISTRY_API_KEY': 'must-not-leak'})
        self.env.start()

    def tearDown(self):
        self.env.stop()
        self.tmp.cleanup()

    def session(self, ep, mode=0o600, tok=TOKEN):
        p = os.path.join(self.private, f'{ep}.session')
        with open(p, 'w') as f:
            f.write(tok + '\n')
        os.chmod(p, mode)

    def agent_spec(self, events, rc=0):
        p = os.path.join(self.tmp.name, 'fake.json')
        with open(p, 'w') as f:
            json.dump({'marker_dir': self.markers, 'events': events, 'rc': rc}, f)
        os.environ['FAKE_SPEC'] = p

    def trial_spec(self, agent='claude', graders=None, **extra):
        g = graders or [{'name': 'pass', 'argv': [sys.executable, self.grader, self.marker, 'pass'], 'cwd': 'backend'},
                        {'name': 'fail', 'argv': [sys.executable, self.grader, self.marker, 'fail']}]
        spec = {'project': 'platform-source', 'repo': '.', 'exec_commands': ['adk-contract-file', 'ruff-file'],
                'episodes': [{'id': 'tps-a', 'agent': agent, 'model': 'm-test', 'task': 'adk2-phase0-1-4',
                              'task_text_file': 'task.md', 'worktree': 'repo-wt', 'deadline_s': 30, 'graders': g,
                              **extra}]}
        p = os.path.join(self.tmp.name, 'trial.json')
        with open(p, 'w') as f:
            json.dump(spec, f)
        return trial.load_spec(p)

    def run_ep(self, agent='claude', graders=None, client=None, worldd=None, events=None):
        if os.path.exists(self.marker):
            os.unlink(self.marker)
        self.agent_spec(events or (CLAUDE_EVENTS if agent == 'claude' else CODEX_EVENTS))
        self.session('tps-a')
        self.client = client or FakeClient(self.marker)
        self.worldd = worldd or FakeWorldd()
        self.evidence = os.path.join(self.rig, 'evidence', 'platform-source')
        return trial.run_episode(self.trial_spec(agent, graders), 'tps-a', rig=self.rig, addr='127.0.0.1:7644',
                                 evidence=self.evidence, client=self.client, runner=self.worldd)

    def facts(self):
        with open(os.path.join(self.markers, 'agent.facts')) as f:
            return json.load(f)


class Argv(unittest.TestCase):
    def test_claude_allows_workspace_exec_plus_the_eight(self):
        argv = trial.build_trial_argv('claude', model='m', prompt_text='P', cwd='/v', mcp_config_path='/p/x.mcp.json',
                                      addr='127.0.0.1:7644')
        allowed = argv[argv.index('--allowedTools') + 1].split(',')
        self.assertIn('mcp__world__workspace-exec', allowed)
        self.assertEqual(allowed, list(arms.world_allowed_tools()) + ['mcp__world__workspace-exec'])
        self.assertEqual(argv[argv.index('--tools') + 1], '')  # still no native tools
        self.assertNotIn('--max-budget-usd', argv)  # subscription (F1): no budget flag unless the spec sets one

    def test_floor_argv_is_unchanged(self):
        golden = arms.golden_argv('claude', 'world')
        self.assertNotIn('workspace-exec', golden[golden.index('--allowedTools') + 1])

    def test_codex_has_no_tool_allowlist_and_no_shell(self):
        argv = trial.build_trial_argv('codex', model='m', prompt_text='P', cwd='/v', mcp_config_path=None,
                                      addr='127.0.0.1:7644')
        self.assertFalse([a for a in argv if 'enabled_tools' in a or 'disabled_tools' in a])
        self.assertIn('shell_tool', argv)
        self.assertIn(f'mcp_servers.world.bearer_token_env_var={arms.TOKEN_ENV}', argv)

    def test_prompt_is_task_verbatim_plus_note(self):
        p = trial.build_prompt('Fix #67.\n', ['quick-check'])
        self.assertTrue(p.startswith('Fix #67.\n'))
        self.assertIn('`quick-check`', p)
        self.assertIn('summarising the change', p)


class Episode(Base):
    def test_claude_episode_end_to_end(self):
        row = self.run_ep('claude')
        f = self.facts()
        self.assertIn('mcp__world__workspace-exec', f['argv'][f['argv'].index('--allowedTools') + 1])
        self.assertEqual(f['mcp_mode'], '0o600')
        self.assertIsNone(f['token_env'])  # Claude: token only in the 0600 config file
        self.assertFalse(f['registry_key'])
        self.assertEqual(f['cwd'], os.path.realpath(os.path.join(self.rig, 'void', 'tps-a')))
        self.assertTrue(row['agent_run']['mcp_config_deleted'])
        self.assertEqual(row['agent_run']['rc'], 0)
        self.assertEqual(row['agent_summary'], 'Added tests 1-4.')
        # log range: head read before (agent not yet run) and after
        heads = [c for c in self.client.calls if c[0] == 'head']
        self.assertEqual(heads, [('head', False), ('head', True)])
        self.assertEqual(row['log_range'], {'from': 11, 'to': 14})
        # graders after the agent
        for mode in ('pass', 'fail'):
            with open(os.path.join(self.markers, f'grader.{mode}')) as fh:
                self.assertEqual(fh.read(), 'after-agent')
        # one outcome per grader, correct payload
        outs = [c for c in self.worldd.calls if c['verb'] == 'outcome']
        self.assertEqual(len(outs), 2)
        p0, p1 = outs[0]['payload'], outs[1]['payload']
        self.assertEqual((p0['outcome']['class'], p0['outcome']['errorCategory']), ('PASS', ''))
        self.assertEqual((p1['outcome']['class'], p1['outcome']['errorCategory']), ('FAIL', 'grader_failed'))
        for p in (p0, p1):
            self.assertEqual(p['arm'], 'world')
            self.assertEqual(p['entries'], {'from': 11, 'to': 14})
            self.assertEqual(p['task'], 'adk2-phase0-1-4')
            self.assertEqual(p['label'], 'trial-platform-source-tps-a')
            self.assertEqual(p['run'], 1)
            self.assertEqual(p['episodeId'], 'tps-a')
            self.assertEqual(p['grading']['graderVersion'], 'trial-grader/1')
            self.assertEqual(p['grading']['toolSha256'], trial.sha256_file(os.path.realpath(sys.executable)))
            self.assertTrue(p['grading']['rule'].startswith('the operator runs '))
            self.assertTrue(p['grading']['rule'].endswith(' unsandboxed in the episode worktree'
                                                          + (' (cwd backend)' if p is p0 else '') + ' after the agent stops'))
        self.assertEqual(p0['grading']['argv'], [sys.executable, self.grader, self.marker, 'pass'])
        for c in self.worldd.calls:
            self.assertEqual(c['argv'][:3], [os.path.join(self.rig, 'bin', 'ailang-worldd'), '--addr', 'http://127.0.0.1:7644'])
            self.assertNotIn('AILANG_REGISTRY_API_KEY', c['env'])
            self.assertNotIn('WORLD_SESSION', c['env'])
        self.assertIn(os.path.join(self.private, 'tps-a.session'), outs[0]['argv'])
        self.assertEqual(row['graders'][0]['outcome_record']['entry'], 101)
        self.assertEqual(row['outcome_record_errors'], 0)
        self.assertEqual(row['provenance']['trailer'], 'World-Provenance: store=/rig/world.db episode=tps-a entries=11-14')
        # the row file, the transcript path + hash
        with open(os.path.join(self.evidence, 'episodes', 'tps-a.json')) as fh:
            on_disk = json.load(fh)
        self.assertEqual(on_disk['transcript']['sha256'], trial.sha256_file(on_disk['transcript']['path']))
        self.assertEqual(on_disk['solution']['state'], 'unchanged')

    def test_codex_gets_token_by_env_only(self):
        row = self.run_ep('codex')
        f = self.facts()
        self.assertEqual(f['token_env'], TOKEN)
        self.assertNotIn('--mcp-config', f['argv'])
        self.assertEqual(row['agent_summary'], 'Added tests 5-8.')
        self.assertEqual(len([c for c in self.worldd.calls if c['verb'] == 'outcome']), 2)

    def test_world_episode_banks_the_report_with_the_row(self):
        row = self.run_ep('claude')
        reps = [c for c in self.worldd.calls if c['verb'] == 'episode']
        self.assertEqual(len(reps), 1)
        self.assertEqual(reps[0]['argv'][3:], ['episode', 'report', 'tps-a', '--json'])
        path = os.path.join(self.evidence, 'reports', 'tps-a.json')
        self.assertEqual(row['report'], {'path': path, 'sha256': trial.sha256_file(path), 'total_calls': 4})
        with open(path) as fh:
            self.assertEqual(json.load(fh)['calls_by_tool'], REPORT['calls_by_tool'])
        with open(os.path.join(self.evidence, 'episodes', 'tps-a.json')) as fh:
            self.assertEqual(json.load(fh)['report']['path'], path)  # in the same row write

    def test_report_failure_is_recorded_not_raised(self):
        row = self.run_ep('claude', worldd=FakeWorldd(report_rc=1))
        self.assertIn('episode has no committed calls', row['report']['error'])
        self.assertFalse(os.path.exists(os.path.join(self.evidence, 'reports', 'tps-a.json')))
        self.assertEqual(row['outcome_record_errors'], 0)

    def test_report_for_another_episode_is_refused(self):
        class Wrong(FakeWorldd):
            def __call__(self, argv, **kw):
                p = super().__call__(argv, **kw)
                if argv[3] == 'episode':
                    p.stdout = p.stdout.replace(b'"tps-a"', b'"other"')
                return p
        row = self.run_ep('claude', worldd=Wrong())
        self.assertEqual(row['report'], {'error': 'report does not describe this episode'})

    def run_shell(self, agent='claude'):
        if os.path.exists(self.marker):
            os.unlink(self.marker)
        self.agent_spec(CLAUDE_EVENTS if agent == 'claude' else CODEX_EVENTS)
        os.unlink(os.path.join(self.private, 'tps-a.session')) if os.path.exists(os.path.join(self.private, 'tps-a.session')) else None
        self.client, self.worldd = FakeClient(self.marker), FakeWorldd()
        self.evidence = os.path.join(self.rig, 'evidence', 'platform-source')
        spec = self.trial_spec(agent, arm='shell', pair_of='tps-w')
        return trial.run_episode(spec, 'tps-a', rig=self.rig, addr='127.0.0.1:7644', evidence=self.evidence,
                                 client=self.client, runner=self.worldd)

    def test_shell_arm_runs_in_the_worktree_with_no_world_contact(self):
        row = self.run_shell('claude')
        f = self.facts()
        self.assertEqual(f['cwd'], os.path.realpath(self.wt))
        self.assertEqual(f['mcp_mode'], 'inline:{"mcpServers":{}}')  # the floor's no-MCP shell argv
        self.assertEqual(f['argv'], arms.build_argv('claude', 'shell', model='m-test', prompt_text=f['argv'][f['argv'].index('-p') + 1],
                                                    cwd=self.wt, mcp_config_path=None, max_budget_usd=None)[1:])
        self.assertIsNone(f['token_env'])
        self.assertFalse(f['registry_key'])
        self.assertEqual(self.client.calls, [])   # no daemon contact at all
        self.assertEqual(self.worldd.calls, [])   # no outcome record, no provenance
        self.assertEqual((row['arm'], row['pair_of']), ('shell', 'tps-w'))
        self.assertIsNone(row['provenance'])
        self.assertIsNone(row['report'])          # no World report for the shell arm
        self.assertEqual(row['log_range'], {'from': None, 'to': None})
        self.assertEqual([g['name'] for g in row['graders']], ['pass', 'fail'])
        self.assertTrue(all('outcome_record' not in g for g in row['graders']))
        self.assertEqual(row['outcome_record_errors'], 0)
        prompt = f['argv'][f['argv'].index('-p') + 1]
        self.assertTrue(prompt.startswith('# Issue 67'))
        self.assertIn('normal coding tools', prompt)
        self.assertNotIn('workspace-exec', prompt)

    def test_codex_shell_arm_is_workspace_write_without_token(self):
        self.run_shell('codex')
        f = self.facts()
        self.assertIsNone(f['token_env'])
        self.assertEqual(f['argv'][f['argv'].index('-s') + 1], 'workspace-write')

    def test_world_arm_is_the_default_and_bad_arms_refuse(self):
        self.assertEqual(trial.episode_spec(self.trial_spec(), 'tps-a')['arm'], 'world')
        with self.assertRaises(trial.TrialError):
            trial.episode_spec(self.trial_spec(arm='docker'), 'tps-a')
        with self.assertRaises(trial.TrialError):
            trial.episode_spec(self.trial_spec(pair_of='../x'), 'tps-a')

    def test_no_token_anywhere_in_evidence(self):
        self.run_ep('claude')
        for d, _, files in os.walk(self.evidence):
            for n in files:
                with open(os.path.join(d, n), 'rb') as fh:
                    self.assertNotIn(TOKEN.encode(), fh.read(), n)

    def test_graders_get_no_secrets(self):
        row = self.run_ep('claude', graders=[{'name': 'env', 'argv': [sys.executable, self.grader, self.marker, 'env']}])
        self.assertEqual(row['graders'][0]['stdout_head'], 'False False\n')

    def test_preflight_without_workspace_exec_spawns_nothing(self):
        row = self.run_ep('claude', client=FakeClient(self.marker, tools=EIGHT))
        self.assertEqual(row['agent_run']['error_category'], 'harness_setup')
        self.assertFalse(os.path.exists(self.marker))
        self.assertEqual(row['graders'], [])
        self.assertEqual(self.worldd.calls, [])

    def test_preflight_against_fakeworld(self):
        fw = fakeworld.FakeWorld(TOKEN, tools=NINE)
        fw.genesis()
        addr = fw.start()
        try:
            got = trial.trial_preflight(trial.world.Client(addr), TOKEN)
            self.assertIn('workspace-exec', got['tools'])
            fw.tools = EIGHT
            with self.assertRaises(trial.world.HarnessSetup):
                trial.trial_preflight(trial.world.Client(addr), TOKEN)
        finally:
            fw.stop()

    def test_record_failure_is_counted_not_raised(self):
        row = self.run_ep('claude', worldd=FakeWorldd(outcome_rc=2))
        self.assertEqual(row['outcome_record_errors'], 2)
        self.assertIn('401', row['graders'][0]['outcome_record']['error'])
        row = self.run_ep('claude', worldd=FakeWorldd(raise_on='outcome'))
        self.assertEqual(row['outcome_record_errors'], 2)
        with mock.patch.object(trial, 'record_outcome', side_effect=RuntimeError('boom')):
            row = self.run_ep('claude')
        self.assertEqual(row['outcome_record_errors'], 2)
        self.assertIn('recorder: boom', row['graders'][1]['outcome_record']['error'])
        with mock.patch('sys.stdout'):
            self.assertEqual(trial.main(['pr-body', '--episode-row', row['_path']]), 0)  # still has a body

    def test_missing_log_to_is_a_counted_payload_error(self):
        class NoAfter(FakeClient):
            def head_index(self):
                if os.path.exists(self.marker):
                    raise trial.world.WorldReadError('daemon gone')
                return 10
        row = self.run_ep('claude', client=NoAfter(self.marker))
        self.assertEqual(row['outcome_record_errors'], 2)
        self.assertTrue(row['graders'][0]['outcome_record']['error'].startswith('payload:'))
        self.assertIn('error', row['provenance'])

    def test_byte_bounded_heads(self):
        row = self.run_ep('claude', graders=[{'name': 'mb', 'argv': [sys.executable, self.grader, self.marker, 'multibyte']}])
        g = row['graders'][0]
        self.assertEqual(g['stdout_head'], 'é' * 200)
        self.assertLessEqual(len(g['stderr_head'].encode()), 400)
        p = [c for c in self.worldd.calls if c['verb'] == 'outcome'][0]['payload']
        self.assertLessEqual(len(p['outcome']['stdoutHead'].encode()), 400)
        self.assertLessEqual(len(p['outcome']['stderrHead'].encode()), 400)
        with open(g['stdout_path'], 'rb') as fh:
            self.assertEqual(len(fh.read()), 2000)  # full output kept as evidence

    def test_grader_wall_cap(self):
        row = self.run_ep('claude', graders=[{'name': 'slow', 'timeout_s': 0.5,
                                              'argv': [sys.executable, self.grader, self.marker, 'sleep']}])
        g = row['graders'][0]
        self.assertTrue(g['timed_out'])
        self.assertIsNone(g['rc'])
        self.assertLess(g['wall_ms'], 15000)
        p = [c for c in self.worldd.calls if c['verb'] == 'outcome'][0]['payload']
        self.assertEqual((p['outcome']['class'], p['outcome']['errorCategory']), ('FAIL', 'grader_failed'))

    def test_grader_timeout_kills_the_whole_group(self):
        pidfile = os.path.join(self.tmp.name, 'grandchild.pid')
        g = trial.run_grader({'name': 'bg', 'cwd': '.', 'timeout_s': 0.5,
                              'argv': ['/bin/sh', '-c', f'sleep 60 & echo $! > {pidfile}; wait']},
                             self.tmp.name)
        self.assertTrue(g['timed_out'])
        with open(pidfile) as f:
            pid = int(f.read().strip())
        deadline = time.monotonic() + 5
        alive = True
        while alive and time.monotonic() < deadline:
            try:
                os.kill(pid, 0)
                time.sleep(0.1)
            except ProcessLookupError:
                alive = False
        self.assertFalse(alive, 'grandchild survived the grader timeout')

    def test_unreadable_tool_does_not_raise(self):
        tool = os.path.join(self.tmp.name, 'noread.sh')
        with open(tool, 'w') as f:
            f.write('#!/bin/sh\nexit 0\n')
        os.chmod(tool, 0o111)
        try:
            g = trial.run_grader({'name': 'nr', 'cwd': '.', 'timeout_s': 5, 'argv': [tool]}, self.tmp.name)
        finally:
            os.chmod(tool, 0o755)
        self.assertEqual(g['tool_sha256'], '')

    def test_grader_output_is_capped(self):
        old = trial.GRADER_KEEP_BYTES
        trial.GRADER_KEEP_BYTES = 1000
        try:
            g = trial.run_grader({'name': 'big', 'cwd': '.', 'timeout_s': 10,
                                  'argv': [sys.executable, '-c', 'import sys; sys.stdout.write("x" * 50000)']},
                                 self.tmp.name)
        finally:
            trial.GRADER_KEEP_BYTES = old
        self.assertEqual(len(g['_stdout']), 1000)
        self.assertEqual(g['stdout_bytes'], 50000)
        self.assertTrue(g['passed'])

    def test_missing_grader_tool_is_recorded_fail(self):
        row = self.run_ep('claude', graders=[{'name': 'nope', 'argv': ['./bin/does-not-exist', '--check']}])
        g = row['graders'][0]
        self.assertTrue(g['spawn_error'])
        self.assertEqual(g['tool_sha256'], '')
        self.assertEqual(row['outcome_record_errors'], 0)

    def test_session_mode_is_checked(self):
        self.agent_spec(CLAUDE_EVENTS)
        self.session('tps-a', mode=0o644)
        with self.assertRaises(trial.TrialError):
            trial.run_episode(self.trial_spec(), 'tps-a', rig=self.rig, addr='x', evidence=self.tmp.name,
                              client=FakeClient(self.marker), runner=FakeWorldd())
        self.assertFalse(os.path.exists(self.marker))


class WorktreeState(Base):
    def test_change_and_untracked_are_hashed(self):
        self.assertEqual(trial.worktree_state(self.wt)['state'], 'unchanged')
        with open(os.path.join(self.wt, 'backend', 'a.py'), 'w') as f:
            f.write('x = 2\n')
        s1 = trial.worktree_state(self.wt)
        self.assertEqual(s1['state'], 'changed')
        self.assertRegex(s1['sha256'], '^[0-9a-f]{64}$')
        with open(os.path.join(self.wt, 'backend', 'test_new.py'), 'w') as f:
            f.write('def test(): pass\n')
        s2 = trial.worktree_state(self.wt)
        self.assertNotEqual(s1['sha256'], s2['sha256'])
        self.assertEqual(s2['untracked'], 1)

    def test_not_a_repo_is_unknown(self):
        self.assertEqual(trial.worktree_state(self.markers)['state'], 'unknown')


class SpecValidation(Base):
    def test_bad_specs_are_refused(self):
        spec = self.trial_spec()
        for bad in ({'agent': 'gpt'}, {'graders': [{'argv': []}]}, {'graders': [{'argv': ['x'], 'cwd': '../out'}]},
                    {'graders': [{'argv': ['x'], 'cwd': '/abs'}]}, {'graders': [{'argv': ['x']}, {'argv': ['x']}]}):
            s = dict(spec, episodes=[dict(spec['episodes'][0], **bad)])
            with self.assertRaises(trial.TrialError, msg=bad):
                trial.episode_spec(s, 'tps-a')
        with self.assertRaises(trial.TrialError):
            trial.episode_spec(spec, 'nope')

    def test_defaults(self):
        spec = self.trial_spec()
        e = dict(spec['episodes'][0])
        del e['deadline_s']
        e['graders'] = [{'argv': ['make', 'adk-conformance'], 'cwd': 'backend'}]
        got = trial.episode_spec(dict(spec, episodes=[e]), 'tps-a')
        self.assertEqual(got['deadline_s'], 2700)
        self.assertEqual(got['graders'][0]['timeout_s'], 900)
        self.assertEqual(got['graders'][0]['name'], 'make')
        self.assertEqual(got['worktree'], self.wt)


class PrBody(Base):
    def test_body_has_trailer_refs_and_no_secrets(self):
        row = self.run_ep('claude')
        with mock.patch('subprocess.run', side_effect=AssertionError('pr-body must not spawn')), \
                mock.patch('subprocess.Popen', side_effect=AssertionError('pr-body must not spawn')):
            body = trial.pr_body(self._load(row['_path']), evidence_link='design_docs/verification/world-trial-x/')
        self.assertIn('\nWorld-Provenance: store=/rig/world.db episode=tps-a entries=11-14\n', body)
        self.assertTrue(body.rstrip('\n').endswith('entries=11-14'))  # trailer is the last line
        self.assertIn('`ailang-world`', body)
        self.assertIn('outcome `sha256:' + '%064x' % 1 + '` (entry 101)', body)
        self.assertIn('**PASS**', body)
        self.assertIn('**FAIL rc=3**', body)
        self.assertIn('Added tests 1-4.', body)
        self.assertIn('Evidence: design_docs/verification/world-trial-x/', body)
        self.assertNotIn(TOKEN, body)
        self.assertNotIn('Bearer', body)

    def test_no_trailer_refuses(self):
        with self.assertRaises(trial.TrialError):
            trial.pr_body({'episode': 'e', 'provenance': {'error': 'rc=1'}})

    def test_token_shaped_text_refuses(self):
        row = self.run_ep('claude')
        row = self._load(row['_path'])
        row['agent_summary'] = 'my token is ' + TOKEN
        with self.assertRaises(trial.TrialError):
            trial.pr_body(row)
        with mock.patch('sys.stdout'), mock.patch('sys.stderr'):
            self.assertEqual(trial.main(['pr-body', '--episode-row', self._dump(row)]), 2)

    def _load(self, p):
        with open(p) as f:
            return json.load(f)

    def _dump(self, row):
        p = os.path.join(self.tmp.name, 'row.json')
        with open(p, 'w') as f:
            json.dump(row, f)
        return p


class AgainstTheRealValidator(Base):
    """The trial payload passes the Go validator (rc != 1: it reaches the commit step and fails on
    the unreachable address); a broken one is refused before any send (rc 1)."""

    @classmethod
    def setUpClass(cls):
        cls.btmp = tempfile.mkdtemp(prefix='trial-validator-')
        cls.bin = os.path.join(cls.btmp, 'ailang-worldd')
        subprocess.run(['go', 'build', '-o', cls.bin, './cmd/ailang-worldd'], cwd=REPO, check=True, capture_output=True)

    def setUp(self):
        super().setUp()
        os.makedirs(os.path.join(self.rig, 'bin'))
        os.symlink(self.bin, os.path.join(self.rig, 'bin', 'ailang-worldd'))
        self.session('tps-a')

    def payloads(self):
        captured = []

        def capture(argv, input=None, **kw):
            if argv[3] == 'outcome':
                captured.append(json.loads(input))
            return subprocess.CompletedProcess(argv, 2, b'', b'')
        self.run_ep('claude', worldd=capture, graders=[
            {'name': 'pass', 'argv': [sys.executable, self.grader, self.marker, 'pass']},
            {'name': 'mb', 'argv': [sys.executable, self.grader, self.marker, 'multibyte']},
            {'name': 'nope', 'argv': ['./bin/does-not-exist']}])
        return captured

    def real(self, payload):
        # Port 9 (discard) on loopback: nothing answers, so a valid payload fails at commit (exit 3).
        return trial.record_outcome(self.rig, '127.0.0.1:9', 'tps-a', payload)

    def test_valid_payloads_pass_validation(self):
        ps = self.payloads()
        self.assertEqual(len(ps), 3)
        for p in ps:
            got = self.real(p)
            self.assertIn('error', got)
            self.assertFalse(got['error'].startswith('rc=1:'), got['error'])
            self.assertNotIn('refused:', got['error'])

    def test_empty_episode_range_passes(self):
        p = self.payloads()[0]
        p['entries'] = {'from': 15, 'to': 14}  # a World episode that committed nothing
        self.assertFalse(self.real(p)['error'].startswith('rc=1:'))

    def test_broken_payloads_are_refused(self):
        base = self.payloads()[0]
        for mut in (lambda p: p.update(entries=None), lambda p: p['outcome'].update(stdoutHead='x' * 401),
                    lambda p: p['grading'].update(argv=[]), lambda p: p.update(extra=1)):
            p = json.loads(json.dumps(base))
            mut(p)
            got = self.real(p)
            self.assertTrue(got['error'].startswith('rc=1:'), got['error'])
            self.assertIn('refused:', got['error'])


if __name__ == '__main__':
    unittest.main()
