#!/usr/bin/env python3
"""AC2.4 — solution provenance from World's own log; §4.4 pre-flight (AC2.2's served half).

Driven against ``fakeworld.FakeWorld``, whose objects use the host's codecs (request bytes,
effect-record field order, invocation-record v2). The graded bytes come from
``world.check_task_run``, which reads the arm's WORKTREE file exactly as the grader does.

Mutant killed here: MUT-NATIVE-WRITE (the solution taken from ``void/`` instead of the worktree):
``test_world_written_solution_is_world_provenance`` plants a native write in void/ and must stay
``world``.
"""
import os
import sys
import tempfile
import unittest

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import fakeworld  # noqa: E402
import world  # noqa: E402

TOKEN = 'f' * 64
EP = 'fl-tune-cx'
SOL = 'benchmark/solution.ail'
BASE = b'module benchmark/solution\n// TODO\n'
GOOD = 'module benchmark/solution\nexport func main() -> () ! {IO} { println("ok") }\n'


class Fixture(unittest.TestCase):
    def setUp(self):
        self.fw = fakeworld.FakeWorld(TOKEN)
        self.fw.genesis()
        self.addr = self.fw.start()
        self.client = world.Client(self.addr)
        self.tmp = tempfile.TemporaryDirectory()
        self.wt = os.path.join(self.tmp.name, 'ws', EP)
        self.void = os.path.join(self.tmp.name, 'void', EP)
        os.makedirs(os.path.join(self.wt, 'benchmark'))
        os.makedirs(self.void)
        self.put(self.wt, BASE)

    def tearDown(self):
        self.fw.stop()
        self.tmp.cleanup()

    def put(self, root, data):
        os.makedirs(os.path.join(root, 'benchmark'), exist_ok=True)
        with open(os.path.join(root, SOL), 'wb') as f:
            f.write(data)

    def check(self, log_from):
        log_to = self.client.head_index()
        writes = world.workspace_writes(self.client, EP, log_from, log_to)
        p = world.check_task_run(self.client, episode=EP, log_from=log_from, log_to=log_to,
                                 worktree=self.wt, void_dir=self.void, base=BASE)
        self.assertEqual(p['world_writes'], len(writes))
        return p, writes


class Provenance(Fixture):
    def test_world_written_solution_is_world_provenance(self):
        log_from = self.client.head_index() + 1
        self.fw.write(EP, SOL, GOOD)
        self.put(self.wt, GOOD.encode())             # what World's handler wrote
        self.put(self.void, b'module x\n// native\n')  # a native write in the agent's cwd
        p, writes = self.check(log_from)
        self.assertEqual(len(writes), 1)
        self.assertEqual(p['solution_provenance'], 'world', p)
        self.assertFalse(p['native_write_detected'])
        self.assertEqual(p['last_world_write_index'], log_from)

    def test_native_overwrite_of_the_worktree_is_detected(self):
        log_from = self.client.head_index() + 1
        self.fw.write(EP, SOL, GOOD)
        self.put(self.wt, GOOD.encode() + b'// patched natively\n')
        p, _ = self.check(log_from)
        self.assertTrue(p['native_write_detected'])
        self.assertEqual(p['reason'], 'bytes_differ')

    def test_native_write_with_no_world_write_is_detected(self):
        log_from = self.client.head_index() + 1
        self.fw.read(EP, SOL)
        self.put(self.wt, GOOD.encode())
        p, writes = self.check(log_from)
        self.assertEqual(writes, [])
        self.assertTrue(p['native_write_detected'])

    def test_untouched_placeholder_is_world(self):
        log_from = self.client.head_index() + 1
        p, _ = self.check(log_from)
        self.assertEqual(p['solution_provenance'], 'world')

    def test_last_write_then_edits_fold_in_order(self):
        log_from = self.client.head_index() + 1
        self.fw.write(EP, SOL, 'module benchmark/solution\nA\n')
        self.fw.write(EP, './' + SOL, GOOD)  # path normalised
        self.fw.edit(EP, SOL, 'println("ok")', 'println("ok!")')
        self.fw.edit(EP, SOL, 'NOT THERE', 'x', outcome='refused')  # policy-tool said ok:false
        self.fw.write(EP, SOL, 'denied', outcome='denied')
        self.fw.write(EP, SOL, 'failed', outcome='failed')
        self.fw.write('fl-other', SOL, 'another episode')
        self.fw.write(EP, 'benchmark/other.ail', 'other file')
        self.put(self.wt, GOOD.replace('println("ok")', 'println("ok!")').encode())
        p, writes = self.check(log_from)
        self.assertEqual(len(writes), 7)  # the other episode's write is not this episode's
        self.assertEqual(p['solution_provenance'], 'world', p)

    def test_writes_before_log_from_do_not_count(self):
        self.fw.write(EP, SOL, GOOD)  # an earlier task-run
        log_from = self.client.head_index() + 1
        self.put(self.wt, GOOD.encode())  # the reset should have put BASE back
        p, writes = self.check(log_from)
        self.assertEqual(writes, [])
        self.assertTrue(p['native_write_detected'])

    def test_edit_on_a_natively_changed_file_is_detected(self):
        log_from = self.client.head_index() + 1
        self.fw.edit(EP, SOL, 'NATIVE', 'WORLD')  # applied ok, so the real file held NATIVE once
        self.put(self.wt, b'WORLD')
        p, _ = self.check(log_from)
        self.assertTrue(p['native_write_detected'])
        self.assertEqual(p['reason'], 'edit_base_mismatch')

    def test_log_paging_reaches_past_500_entries(self):
        log_from = self.client.head_index() + 1
        for _ in range(520):
            self.fw.read('fl-noise', 'x')
        self.fw.write(EP, SOL, GOOD)
        self.put(self.wt, GOOD.encode())
        p, writes = self.check(log_from)
        self.assertEqual(len(writes), 1)
        self.assertEqual(p['solution_provenance'], 'world')

    def test_request_bytes_roundtrip(self):
        b = fakeworld.request_bytes('Workspace.Write', 'worktree', 1, 1759700000, b'{"op":"write"}')
        self.assertEqual(b[:30], b'15:Workspace.Write8:worktree1:')
        self.assertEqual(world.parse_request_bytes(b),
                         ('Workspace.Write', 'worktree', 1, 1759700000, b'{"op":"write"}'))


class Preflight(Fixture):
    def test_green_preflight(self):
        r = world.preflight(self.client, TOKEN)
        self.assertEqual(len(r['tools']), 8)

    def test_ninth_tool_is_harness_setup(self):
        self.fw.tools.append('workspace-exec')
        with self.assertRaises(world.HarnessSetup) as cm:
            world.preflight(self.client, TOKEN)
        self.assertIn('workspace-exec', str(cm.exception))

    def test_missing_tool_is_harness_setup(self):
        self.fw.tools.remove('ailang-run')
        with self.assertRaises(world.HarnessSetup):
            world.preflight(self.client, TOKEN)

    def test_wrong_token_is_harness_setup_and_never_echoed(self):
        with self.assertRaises(world.HarnessSetup) as cm:
            world.preflight(self.client, 'e' * 64)
        self.assertNotIn('e' * 64, str(cm.exception))

    def test_unhealthy_daemon_is_harness_setup(self):
        self.fw.health_status = 503
        with self.assertRaises(world.HarnessSetup):
            world.preflight(self.client, TOKEN)

    def test_dead_daemon_is_harness_setup(self):
        with self.assertRaises(world.HarnessSetup):
            world.preflight(world.Client('127.0.0.1:1'), TOKEN)


if __name__ == '__main__':
    unittest.main()
