#!/usr/bin/env python3
"""M2 driver — one task-run's agent half, hermetic: FAKE ``claude``/``codex`` on PATH and the
``fakeworld`` daemon. Checks the §4.4 steps 2–4 plumbing and the token discipline:

  * World pre-flight gates the spawn (a failure is ``harness_setup`` and nothing runs);
  * ``log_from``/``log_to`` bracket the run; ``void/`` is checked after it;
  * the 600 s deadline (here 1 s) kills the WHOLE process group, ``finish_reason=timeout``;
  * the Claude World arm's MCP config is mode 0600 in the private dir while the agent runs and
    is gone after; the codex World arm gets the token by env only; no shell arm ever sees it;
  * neither the recorded argv nor the transcript carries the token;
  * the shell arm's ``ailang`` is a symlink to ``$TOOL``.
"""
import json
import os
import stat
import sys
import tempfile
import time
import unittest
from unittest import mock

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import arms  # noqa: E402
import fakeworld  # noqa: E402
import run_task  # noqa: E402

TOKEN = '0123456789abcdef' * 4

FAKE_AGENT = r'''#!/usr/bin/env python3
import hashlib, json, os, shutil, subprocess, sys, time
argv = sys.argv[1:]
agent = os.path.basename(sys.argv[0])
spec = json.load(open(os.environ['FAKE_SPEC']))
open(os.path.join(spec['marker_dir'], agent + '.ran'), 'w').write('1')
facts = {'cwd': os.getcwd(), 'has_token_env': 'WORLD_SESSION' in os.environ,
         'token_env_sha': hashlib.sha256(os.environ.get('WORLD_SESSION', '').encode()).hexdigest(),
         'ailang': os.path.realpath(shutil.which('ailang') or '') if shutil.which('ailang') else None}
if '--mcp-config' in argv:
    p = argv[argv.index('--mcp-config') + 1]
    if os.path.isfile(p):
        facts['mcp_mode'] = oct(os.stat(p).st_mode & 0o777)
        cfg = json.load(open(p))
        h = cfg['mcpServers']['world']['headers']['Authorization']
        facts['mcp_auth_sha'] = hashlib.sha256(h.encode()).hexdigest()
json.dump(facts, open(os.path.join(spec['marker_dir'], agent + '.facts'), 'w'))
if spec.get('grandchild'):
    g = subprocess.Popen(['sleep', '30'])
    open(os.path.join(spec['marker_dir'], 'grandchild.pid'), 'w').write(str(g.pid))
if spec.get('sleep'):
    time.sleep(spec['sleep'])
if spec.get('touch_void'):
    open('native.txt', 'w').write('x')
for ev in spec['events']:
    print(json.dumps(ev), flush=True)
'''

CLAUDE_EVENTS = [
    {'type': 'system', 'subtype': 'init', 'tools': [f'mcp__world__{n}' for n in
                                                    ('ailang-read', 'ailang-write', 'ailang-edit', 'ailang-check',
                                                     'ailang-run', 'builtins-search', 'examples-search', 'ailang-cli')],
     'mcp_servers': [{'name': 'world', 'status': 'connected'}], 'model': 'm'},
    {'type': 'assistant', 'message': {'content': [{'type': 'tool_use', 'name': 'mcp__world__ailang-read', 'input': {}}]}},
    {'type': 'user', 'message': {'content': [{'type': 'tool_result', 'is_error': False, 'content': 'ok'}]}},
    {'type': 'result', 'subtype': 'success', 'is_error': False, 'total_cost_usd': 0.001, 'num_turns': 2},
]
CODEX_EVENTS = [
    {'type': 'thread.started', 'thread_id': 't'},
    {'type': 'item.started', 'item': {'id': 'i1', 'type': 'mcp_tool_call', 'server': 'world', 'tool': 'ailang-read', 'status': 'in_progress'}},
    {'type': 'item.completed', 'item': {'id': 'i1', 'type': 'mcp_tool_call', 'server': 'world', 'tool': 'ailang-read', 'status': 'completed'}},
    {'type': 'item.completed', 'item': {'id': 'i2', 'type': 'file_change', 'status': 'failed', 'changes': []}},
    {'type': 'item.completed', 'item': {'id': 'i3', 'type': 'agent_message', 'text': 'done'}},
    {'type': 'turn.completed', 'usage': {'input_tokens': 10, 'output_tokens': 2}},
]


class Driver(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        t = self.tmp.name
        self.bin = os.path.join(t, 'fakebin')
        os.makedirs(self.bin)
        for name in ('claude', 'codex'):
            p = os.path.join(self.bin, name)
            with open(p, 'w') as f:
                f.write(FAKE_AGENT)
            os.chmod(p, 0o755)
        self.tool = os.path.join(t, 'tool', 'ailang')
        os.makedirs(os.path.dirname(self.tool))
        with open(self.tool, 'w') as f:
            f.write('#!/bin/sh\n')
        os.chmod(self.tool, 0o755)
        self.private = os.path.join(t, 'private')
        os.makedirs(self.private, mode=0o700)
        os.chmod(self.private, 0o700)
        self.wt = os.path.join(t, 'ws', 'fl-x')
        self.void = os.path.join(t, 'void', 'fl-x')
        self.markers = os.path.join(t, 'markers')
        for d in (self.wt, self.void, self.markers):
            os.makedirs(d)
        self.fw = fakeworld.FakeWorld(TOKEN)
        self.fw.genesis()
        self.addr = self.fw.start()
        self.env = mock.patch.dict(os.environ, {'PATH': self.bin + os.pathsep + os.environ.get('PATH', ''),
                                                'WORLD_SESSION': 'leaked-from-parent'})
        self.env.start()

    def tearDown(self):
        self.env.stop()
        self.fw.stop()
        self.tmp.cleanup()

    def spec(self, **kw):
        s = dict(marker_dir=self.markers, **kw)
        p = os.path.join(self.tmp.name, 'spec.json')
        with open(p, 'w') as f:
            json.dump(s, f)
        os.environ['FAKE_SPEC'] = p

    def facts(self, agent):
        with open(os.path.join(self.markers, agent + '.facts')) as f:
            return json.load(f)

    def run_agent(self, agent, arm, **kw):
        tr = os.path.join(self.tmp.name, f'{agent}-{arm}.jsonl')
        extra = dict(void_dir=self.void, episode='fl-x', token=TOKEN) if arm == 'world' else {}
        extra.update(kw)
        row = run_task.run_agent(agent, arm, model='m-test', prompt_text='PROMPT', worktree=self.wt,
                                 private_dir=self.private, transcript_path=tr, tool_bin=self.tool,
                                 addr=self.addr, **extra)
        text = open(tr).read() if os.path.exists(tr) else None
        return row, text

    def test_claude_world_run(self):
        self.spec(events=CLAUDE_EVENTS)
        row, text = self.run_agent('claude', 'world')
        self.assertIsNone(row['error_category'])
        self.assertEqual(row['finish_reason'], 'exit')
        f = self.facts('claude')
        self.assertEqual(f['cwd'], os.path.realpath(self.void))
        self.assertEqual(f['mcp_mode'], '0o600')
        self.assertEqual(f['mcp_auth_sha'], __import__('hashlib').sha256(f'Bearer {TOKEN}'.encode()).hexdigest())
        self.assertFalse(f['has_token_env'])  # neither the token nor the parent's leak
        self.assertEqual(os.listdir(self.private), [])  # config deleted
        self.assertTrue(row['mcp_config_deleted'])
        self.assertEqual(row['world']['log_from'], 1)
        self.assertEqual(row['world']['log_to'], 0)  # nothing committed during the fake run
        self.assertTrue(row['void_empty'])
        self.assertEqual(row['tool_items'], {'mcp__world__ailang-read': 1})
        self.assertEqual(len(row['transcript']['init']['tools']), 8)
        self.assertNotIn(TOKEN, json.dumps(row))
        self.assertNotIn(TOKEN, text)
        self.assertIn('<prompt sha256:', ' '.join(row['argv']))
        self.assertGreaterEqual(row['wall_ms'], 0)

    def test_codex_world_run_token_by_env_only(self):
        self.spec(events=CODEX_EVENTS)
        row, text = self.run_agent('codex', 'world')
        f = self.facts('codex')
        self.assertTrue(f['has_token_env'])
        self.assertEqual(f['token_env_sha'], __import__('hashlib').sha256(TOKEN.encode()).hexdigest())
        self.assertNotIn(TOKEN, json.dumps(row))
        self.assertEqual(row['tool_items'], {'mcp_tool_call': 1, 'file_change': 1, 'agent_message': 1})
        self.assertEqual(row['transcript']['file_change_ok'], 0)
        self.assertEqual(row['transcript']['command_execution'], 0)
        self.assertIn('read-only', row['argv'])

    def test_shell_arm_never_sees_the_token_and_gets_the_tool_ailang(self):
        for agent in arms.AGENTS:
            with self.subTest(agent=agent):
                self.spec(events=CLAUDE_EVENTS if agent == 'claude' else CODEX_EVENTS)
                row, _ = self.run_agent(agent, 'shell')
                f = self.facts(agent)
                self.assertFalse(f['has_token_env'])
                self.assertEqual(f['ailang'], os.path.realpath(self.tool))
                self.assertEqual(f['cwd'], os.path.realpath(self.wt))
                self.assertIsNone(row['world'])

    def test_preflight_failure_is_harness_setup_and_nothing_spawns(self):
        self.spec(events=CLAUDE_EVENTS)
        self.fw.tools.append('workspace-exec')
        row, _ = self.run_agent('claude', 'world')
        self.assertEqual(row['error_category'], 'harness_setup')
        self.assertIn('workspace-exec', row['cause_evidence'])
        self.assertFalse(os.path.exists(os.path.join(self.markers, 'claude.ran')))

    def test_deadline_kills_the_process_group(self):
        self.spec(events=CODEX_EVENTS, sleep=30, grandchild=True)
        t = time.monotonic()
        row, _ = self.run_agent('codex', 'shell', deadline_s=1.0)
        self.assertLess(time.monotonic() - t, 15)
        self.assertEqual(row['finish_reason'], 'timeout')
        self.assertEqual(row['error_category'], 'timeout')
        with open(os.path.join(self.markers, 'grandchild.pid')) as f:
            pid = int(f.read())
        time.sleep(0.2)
        try:
            os.kill(pid, 0)
            alive = True
        except ProcessLookupError:
            alive = False
        except PermissionError:
            alive = True
        if alive:  # a zombie of an exited group is reaped by init; check it is not still running
            import subprocess
            st = subprocess.run(['ps', '-o', 'stat=', '-p', str(pid)], capture_output=True, text=True).stdout.strip()
            alive = bool(st) and not st.startswith('Z')
        self.assertFalse(alive, 'the grandchild survived the deadline')

    def test_void_write_is_reported(self):
        self.spec(events=CODEX_EVENTS, touch_void=True)
        row, _ = self.run_agent('codex', 'world')
        self.assertFalse(row['void_empty'])
        self.assertEqual(row['void_listing'], ['native.txt'])

    def test_private_dir_rules(self):
        self.spec(events=CLAUDE_EVENTS)
        os.chmod(self.private, 0o755)
        with self.assertRaises(run_task.DriverError):
            self.run_agent('claude', 'world')
        os.chmod(self.private, 0o700)
        inside = os.path.join(self.wt, 'priv')
        os.makedirs(inside, mode=0o700)
        os.chmod(inside, 0o700)
        with self.assertRaises(run_task.DriverError):
            run_task.run_agent('claude', 'shell', model='m', prompt_text='p', worktree=self.wt,
                               private_dir=inside, transcript_path=os.path.join(self.tmp.name, 'x.jsonl'))
        self.assertEqual(stat.S_IMODE(os.stat(self.private).st_mode), 0o700)


if __name__ == '__main__':
    unittest.main()
