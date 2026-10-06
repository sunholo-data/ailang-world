#!/usr/bin/env python3
"""M2 — the attended mint block is exactly §4.5's: two ``session new`` commands, the 8 explicit
grants from ``arms.world_grant_args`` (never ``--preset se-tools``), ``--ttl 43200``, ``--out``
under the harness-private dir; the scratch store, never the live one."""
import os
import shlex
import sys
import unittest

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import arms  # noqa: E402
import canary  # noqa: E402


class MintBlock(unittest.TestCase):
    def setUp(self):
        self.L = canary.layout('/scratch/floor-tune')
        self.block = canary.mint_block(self.L)
        self.mints = [shlex.split(ln) for ln in self.block.replace('\\\n', ' ').splitlines()
                      if 'session new' in ln and not ln.lstrip().startswith('#')]

    def test_two_mints_with_the_eight_explicit_grants(self):
        self.assertEqual([m[m.index('new') + 1] for m in self.mints], ['fl-tune-cc', 'fl-tune-cx'])
        for m in self.mints:
            grants = [m[i + 1] for i, a in enumerate(m) if a == '--grant']
            self.assertEqual(grants, arms.world_grant_args()[1::2])
            self.assertEqual(len(grants), 8)
            self.assertNotIn('--preset', m)
            self.assertFalse(any('Workspace.Exec' in a for a in m))
            self.assertEqual(m[m.index('--ttl') + 1], '43200')
            self.assertTrue(m[m.index('--out') + 1].startswith('$FT/private/'))
            self.assertEqual(m[m.index('--db') + 1], '$FT/store/world.db')

    def test_scratch_store_and_stopped_daemon(self):
        self.assertIn('FT=/scratch/floor-tune', self.block)
        self.assertNotIn('.ailang/world/world.db', self.block)  # never the live store
        self.assertIn('curl -sf http://127.0.0.1:7644/v1/health', self.block)
        self.assertIn('world-publish transitions', self.block)
        self.assertIn('canary.py world --root $FT', self.block)


if __name__ == '__main__':
    unittest.main()
