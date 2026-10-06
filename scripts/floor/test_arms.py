#!/usr/bin/env python3
"""The World arm is EXACTLY the 8 AILANG se-tools — never row 140's ``workspace-exec``.

Checked against the se-tools manifest World actually serves (packages/se-tools/transitions.json),
so a 10th served tool, a renamed tool or a drifted effect turns this red and forces a decision.

Mutant killed here: MUT-EXEC-GRANT (World-arm grants include ``Workspace.Exec``).
"""
import json
import os
import sys
import unittest

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import arms  # noqa: E402
import prompt  # noqa: E402

MANIFEST = os.path.join(HERE, '..', '..', 'packages', 'se-tools', 'transitions.json')
EIGHT = ('ailang-read', 'ailang-write', 'ailang-edit', 'ailang-check', 'ailang-run',
         'builtins-search', 'examples-search', 'ailang-cli')
EXEC_TOOL, EXEC_EFFECT = 'workspace-exec', 'Workspace.Exec'


def served():
    with open(MANIFEST) as f:
        return {t['id']: tuple(e['effect'] for e in t['declaredEffects']) for t in json.load(f)}


class WorldArmToolSet(unittest.TestCase):
    def test_exactly_the_eight_ailang_tools(self):
        self.assertEqual(prompt.WORLD_TOOL_NAMES, EIGHT)
        self.assertEqual(arms.world_allowed_tools(), tuple(f'mcp__world__{n}' for n in EIGHT))

    def test_workspace_exec_absent_from_allowed_tools_and_grants(self):
        self.assertNotIn(EXEC_TOOL, prompt.WORLD_TOOL_NAMES)
        self.assertIn(EXEC_TOOL, prompt.WORLD_EXCLUDED_TOOLS)
        self.assertNotIn(f'mcp__world__{EXEC_TOOL}', arms.world_allowed_tools())
        self.assertFalse(any(EXEC_TOOL in t for t in arms.world_allowed_tools()))
        self.assertNotIn(EXEC_EFFECT, arms.world_grant_effects())
        args = arms.world_grant_args()
        self.assertFalse(any(EXEC_EFFECT in a for a in args), args)
        self.assertNotIn('--preset', args)  # --preset se-tools grants Workspace.Exec

    def test_grants_are_exactly_the_eight_tools_declared_effects(self):
        man = served()
        declared = []
        for name in EIGHT:
            self.assertIn(name, man, f'{name} is not served')
            self.assertEqual(arms.WORLD_TOOL_EFFECTS[name], man[name], name)
            declared += [e for e in man[name] if e not in declared]
        self.assertEqual(arms.world_grant_effects(), tuple(declared))
        self.assertEqual(len(declared), 8)
        self.assertEqual(arms.world_grant_args(),
                         [x for e in declared for x in ('--grant', f'{e}=worktree:5000')])

    def test_the_only_served_tool_left_out_is_workspace_exec(self):
        man = served()
        self.assertEqual(set(man) - set(EIGHT), {EXEC_TOOL})  # non-vacuous: it IS served
        self.assertEqual(man[EXEC_TOOL], (EXEC_EFFECT,))

    def test_world_prompt_never_mentions_exec(self):
        block = prompt.world_tools_block()
        for name in EIGHT:
            self.assertIn(f'`{name}`', block)
        self.assertNotIn(EXEC_TOOL, block)
        self.assertNotIn('Exec', block)


if __name__ == '__main__':
    unittest.main()
