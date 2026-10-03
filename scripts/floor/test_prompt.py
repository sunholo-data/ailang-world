#!/usr/bin/env python3
"""AC1.4 — the shell and World prompts differ ONLY in the tools block and the solution path.

The diff is measured on the full task MESSAGE (teaching prompt + task) with difflib, and every
differing line on either side must lie inside the template's tool-specific sections
("## Available Tools", "## Verification Steps") or carry the arm's solution path. Separately, the
World task text must carry no shell command at all (MUT-PROMPT-LEAK) while the shell task text
must (the diff is not vacuous). Rig-only arms render all 23 pinned tasks with the real template
(at the pinned commit) and the real teaching prompt from $TOOL.
"""
import difflib
import os
import sys
import unittest

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import prompt  # noqa: E402

with open(os.path.join(HERE, 'testdata', 'agent_task_ailang.txt')) as _f:
    TEMPLATE = _f.read()
TEACHING = ('# AILANG teaching prompt (fixture)\nUse recursion.\n'
            '-- Run: ailang run --caps IO --entry main program.ail\n')  # the real one mentions it too
SHELL_PATH = '/work/floor/wt-cc/benchmark/solution.ail'
SHELL_CMDS = ('ailang run', 'ailang check', 'ailang test')
TOOL_SECTIONS = ('## Available Tools', '## Verification Steps')

SPECS = [
    {'id': 't_plain', 'description': 'Plain task', 'caps': ['IO'], 'task_prompt': 'Write a program in <LANG>.\n',
     'expected_stdout': 'ok\n'},
    {'id': 'pipeline', 'description': 'Read stdin', 'caps': ['IO'], 'prompt': 'Read numbers from stdin.\n',
     'expected_stdout': '2\n4\n', 'stdin': '1\n2\n'},
    {'id': 'api_call_json', 'description': 'POST', 'caps': ['Net', 'IO'], 'net_allow_localhost': True,
     'task_prompt': 'POST to {{MOCK_HTTP_URL}} in <LANG>.\n', 'expected_stdout': '200\n'},
]
MOCK = 'http://127.0.0.1:49152'


def sections(text):
    """Line index -> the '## ' heading it falls under."""
    cur, out = None, []
    for line in text.split('\n'):
        if line.startswith('## '):
            cur = line
        out.append(cur)
    return out


def render_pair(spec, template=TEMPLATE, teaching=TEACHING):
    kw = dict(template=template, deadline_s=600, mock_url=MOCK)
    s = prompt.render_task(spec, 'shell', solution_path=SHELL_PATH, **kw)
    w = prompt.render_task(spec, 'world', solution_path=prompt.WORLD_SOLUTION_PATH, **kw)
    return s, w, prompt.render_message(teaching, s), prompt.render_message(teaching, w)


def disallowed_diff(ms, mw):
    """Differing lines outside the tools block that do not carry a solution path."""
    a, b = ms.split('\n'), mw.split('\n')
    sa, sb = sections(ms), sections(mw)
    bad = []
    for op, i1, i2, j1, j2 in difflib.SequenceMatcher(None, a, b, autojunk=False).get_opcodes():
        if op == 'equal':
            continue
        for i in range(i1, i2):
            if not ((sa[i] or '').startswith(TOOL_SECTIONS) or SHELL_PATH in a[i]):
                bad.append(('shell', i, a[i]))
        for j in range(j1, j2):
            if not ((sb[j] or '').startswith(TOOL_SECTIONS) or prompt.WORLD_SOLUTION_PATH in b[j]):
                bad.append(('world', j, b[j]))
    return bad


class ArmDiff(unittest.TestCase):
    def test_only_tools_block_and_solution_path_differ(self):
        for spec in SPECS:
            with self.subTest(spec['id']):
                s, w, ms, mw = render_pair(spec)
                self.assertNotEqual(ms, mw)
                self.assertEqual(disallowed_diff(ms, mw), [])
                self.assertTrue(ms.startswith(TEACHING.rstrip('\n')) and mw.startswith(TEACHING.rstrip('\n')))

    def test_world_has_no_shell_command_and_shell_does(self):
        for spec in SPECS:
            with self.subTest(spec['id']):
                s, w, _, _ = render_pair(spec)
                for cmd in SHELL_CMDS:
                    self.assertNotIn(cmd, w)
                    self.assertIn(cmd, s)
                for name, _ in prompt.WORLD_TOOLS:
                    self.assertIn(f'`{name}`', w)
                    self.assertNotIn(f'`{name}`', s)
                self.assertIn(prompt.WORLD_RUN_NOTE, w)

    def test_solution_paths(self):
        s, w, _, _ = render_pair(SPECS[0])
        self.assertIn(f'Write your solution to: **{SHELL_PATH}**', s)
        self.assertIn(f'Write your solution to: **{prompt.WORLD_SOLUTION_PATH}**', w)
        self.assertNotIn(SHELL_PATH, w)
        with self.assertRaises(prompt.PromptError):
            prompt.render_task(SPECS[0], 'world', template=TEMPLATE, solution_path=SHELL_PATH, deadline_s=600)
        with self.assertRaises(prompt.PromptError):
            prompt.render_task(SPECS[0], 'shell', template=TEMPLATE, solution_path='benchmark/solution.ail',
                               deadline_s=600)

    def test_go_baseline_fill(self):
        s, w, _, _ = render_pair(SPECS[2])
        self.assertIn(f'POST to {MOCK} in AILANG.', s)
        self.assertIn(f'POST to {MOCK} in AILANG.', w)
        self.assertIn('POST\n\nPOST to', s)  # description + "\n\n" + task text
        self.assertIn('- Timeout: 600 seconds', s)
        self.assertIn('--caps Net,IO --net-allow-http --net-allow-localhost solution.ail', s)
        self.assertNotIn('{{', s.replace('{"path"', ''))
        self.assertIn('Read numbers from stdin.', render_pair(SPECS[1])[0])  # prompt: fallback

    def test_mock_token_needs_a_url(self):
        with self.assertRaises(prompt.PromptError):
            prompt.render_task(SPECS[2], 'shell', template=TEMPLATE, solution_path=SHELL_PATH, deadline_s=600)

    def test_unknown_placeholder_refused(self):
        with self.assertRaises(prompt.PromptError):
            render_pair(SPECS[0], template=TEMPLATE + '\n{{NEW_SLOT}}\n')


REPO = os.environ.get('FLOOR_AILANG_REPO')
TOOL = os.environ.get('FLOOR_TOOL_BIN', os.path.expanduser('~/.pinned-ailang-tools/v0.51.0/ailang'))


@unittest.skipUnless(REPO and os.path.isfile(TOOL), 'rig-only: FLOOR_AILANG_REPO and the pinned tool binary')
class PinnedRender(unittest.TestCase):
    def test_all_23_pinned_tasks(self):
        import corpus
        c = corpus.load_corpus(repo=REPO, commit='76a5aef65c81ac2fb4657fafaeff4a89e9faaf75')
        tmpl = prompt.load_template(REPO, c.commit)
        self.assertEqual(tmpl, TEMPLATE)  # the checked-in fixture is the pinned template
        teaching = prompt.load_teaching(TOOL)
        for t in c.tasks:
            with self.subTest(t.id):
                s, w, ms, mw = render_pair(t.spec, template=tmpl, teaching=teaching)
                self.assertEqual(disallowed_diff(ms, mw), [])
                for cmd in SHELL_CMDS:
                    self.assertNotIn(cmd, w)


if __name__ == '__main__':
    unittest.main()
