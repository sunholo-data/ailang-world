#!/usr/bin/env python3
"""The World arm is EXACTLY the 8 AILANG se-tools — never row 140's ``workspace-exec``.

Checked against the se-tools manifest World actually serves (packages/se-tools/transitions.json),
so a 10th served tool, a renamed tool or a drifted effect turns this red and forces a decision.

Mutant killed here: MUT-EXEC-GRANT (World-arm grants include ``Workspace.Exec``).

M2 (AC2.1): the exact argv per (agent, arm) against the checked-in goldens
``testdata/argv_<agent>_<arm>.json`` AND against the design's own §4.3 shell text (shlex-split,
so a golden regenerated from a drifted builder still fails). Mutants killed here:
MUT-ALLOWEDTOOLS-ONLY (``--tools ""`` dropped) and MUT-CODEX-SHELL (``--disable shell_tool``
dropped); their transcript halves (AC2.2/AC2.3) are ``AttestFixtures`` below.
"""
import json
import os
import shlex
import sys
import unittest

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import arms  # noqa: E402
import attest  # noqa: E402
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


# The §4.3 table, verbatim from design_docs/planned/w-resident-agent-non-inferiority-floor-run.md
# (the World Claude arm's "…" expanded to the 8 names; codex World = "same flags, but …"; the Claude
# shell arm's --permission-mode is acceptEdits as amended by M2, V-M2-3; the M4 additions ACKed by
# Mark Edmondson attended 2026-10-06, D-WORLD-67: Claude isolation in both arms, --max-budget-usd
# 5.00, and the codex World arm's MCP approval override).
ISOLATION_JSON_BODY = ('"enabledPlugins":{"cc-plugin-agents-md@builtin":false,"cc-plugin-telemetry@builtin":false,"cc-plugin-plugin-authoring@builtin":false},"autoMemoryEnabled":false')
DESIGN_CLAUDE_SHELL = (
    'claude -p "$P" --model $M --output-format stream-json --verbose --no-session-persistence '
    '--disable-slash-commands --tools Bash,Read,Write,Edit,Glob,Grep --strict-mcp-config '
    "--mcp-config '{\"mcpServers\":{}}' --permission-mode acceptEdits --settings "
    "'{\"sandbox\":{\"enabled\":true,\"autoAllowBashIfSandboxed\":true,\"allowUnsandboxedCommands\":false},"
    + ISOLATION_JSON_BODY + "}' --setting-sources \"\" --max-budget-usd 5.00")
DESIGN_CLAUDE_WORLD = (
    'claude -p "$P" --model $M --output-format stream-json --verbose --no-session-persistence '
    '--disable-slash-commands --tools "" --strict-mcp-config --mcp-config <private>/<ep>.mcp.json '
    '--allowedTools ' + ','.join(f'mcp__world__{n}' for n in EIGHT)
    + " --setting-sources \"\" --settings '{" + ISOLATION_JSON_BODY + "}' --max-budget-usd 5.00")
DESIGN_CODEX_SHELL = (
    'codex exec --json --ephemeral --ignore-user-config --skip-git-repo-check -m $M '
    '-s workspace-write -C <worktree> --disable plugins --disable apps --disable browser_use '
    '--disable computer_use --disable image_generation --disable multi_agent '
    '-c web_search="disabled" "$P"')
DESIGN_CODEX_WORLD = (
    'codex exec --json --ephemeral --ignore-user-config --skip-git-repo-check -m $M '
    '-s read-only -C <void/ep> --disable plugins --disable apps --disable browser_use '
    '--disable computer_use --disable image_generation --disable multi_agent '
    '--disable shell_tool --disable unified_exec -c web_search="disabled" '
    '-c mcp_servers.world.url="http://127.0.0.1:7644/mcp/" '
    '-c mcp_servers.world.bearer_token_env_var="WORLD_SESSION" '
    '-c mcp_servers.world.required=true '
    '-c mcp_servers.world.default_tools_approval_mode="approve" "$P"')
DESIGN = {('claude', 'shell'): DESIGN_CLAUDE_SHELL, ('claude', 'world'): DESIGN_CLAUDE_WORLD,
          ('codex', 'shell'): DESIGN_CODEX_SHELL, ('codex', 'world'): DESIGN_CODEX_WORLD}


def golden(agent, arm):
    with open(os.path.join(HERE, 'testdata', f'argv_{agent}_{arm}.json')) as f:
        return json.load(f)


def flag_value(argv, flag):
    return argv[argv.index(flag) + 1]


def disabled(argv):
    return [argv[i + 1] for i, a in enumerate(argv) if a == '--disable']


class ArgvGoldens(unittest.TestCase):
    def test_builder_matches_checked_in_golden(self):
        for agent in arms.AGENTS:
            for arm in arms.ARMS:
                with self.subTest(agent=agent, arm=arm):
                    self.assertEqual(arms.golden_argv(agent, arm), golden(agent, arm))

    def test_golden_matches_the_design_text(self):
        for (agent, arm), text in DESIGN.items():
            with self.subTest(agent=agent, arm=arm):
                self.assertEqual(golden(agent, arm), shlex.split(text))

    def test_claude_world_removes_every_builtin(self):  # MUT-ALLOWEDTOOLS-ONLY
        argv = arms.golden_argv('claude', 'world')
        self.assertIn('--tools', argv)
        self.assertEqual(flag_value(argv, '--tools'), '')
        self.assertIn('--strict-mcp-config', argv)
        self.assertEqual(flag_value(argv, '--allowedTools').split(','), list(arms.world_allowed_tools()))

    def test_codex_world_has_no_shell(self):  # MUT-CODEX-SHELL
        argv = arms.golden_argv('codex', 'world')
        self.assertIn('shell_tool', disabled(argv))
        self.assertIn('unified_exec', disabled(argv))
        self.assertEqual(flag_value(argv, '-s'), 'read-only')
        self.assertNotIn('shell_tool', disabled(arms.golden_argv('codex', 'shell')))  # non-vacuous

    def test_arms_differ_only_where_the_design_says(self):
        cs, cw = arms.golden_argv('codex', 'shell'), arms.golden_argv('codex', 'world')
        self.assertEqual(set(disabled(cw)) - set(disabled(cs)), {'shell_tool', 'unified_exec'})
        self.assertEqual(flag_value(cs, '-m'), flag_value(cw, '-m'))
        self.assertIn('web_search=disabled', cs)
        self.assertIn('web_search=disabled', cw)
        shell_tools = flag_value(arms.golden_argv('claude', 'shell'), '--tools').split(',')
        self.assertFalse(set(shell_tools) & {'WebFetch', 'WebSearch'})

    def test_claude_shell_edits_confined_and_bash_sandboxed(self):  # AC2.5, V-M2-3
        argv = arms.golden_argv('claude', 'shell')
        self.assertEqual(flag_value(argv, '--permission-mode'), 'acceptEdits')  # never bypassPermissions
        sb = json.loads(flag_value(argv, '--settings'))['sandbox']
        self.assertEqual(sb, {'enabled': True, 'autoAllowBashIfSandboxed': True, 'allowUnsandboxedCommands': False})
        self.assertEqual(flag_value(argv, '--setting-sources'), '')  # isolation is canonical (D-WORLD-67)
        self.assertEqual(flag_value(arms.golden_argv('codex', 'shell'), '-s'), 'workspace-write')

    def test_codex_world_tools_are_approved(self):  # MUT-CODEX-NOAPPROVE (V-M4-3, D-WORLD-67)
        argv = arms.golden_argv('codex', 'world')
        cs = [argv[i + 1] for i, a in enumerate(argv) if a == '-c']
        self.assertIn('mcp_servers.world.default_tools_approval_mode=approve', cs)
        self.assertNotIn('mcp_servers.world.default_tools_approval_mode=approve',
                         arms.golden_argv('codex', 'shell'))  # the shell arm has no World server

    def test_claude_isolation_and_budget_identical_in_both_arms(self):  # D-WORLD-67
        sh, wo = arms.golden_argv('claude', 'shell'), arms.golden_argv('claude', 'world')
        for argv in (sh, wo):
            self.assertEqual(flag_value(argv, '--setting-sources'), '')
            self.assertEqual(flag_value(argv, '--max-budget-usd'), '5.00')
            s = json.loads(flag_value(argv, '--settings'))
            self.assertEqual({k: s[k] for k in arms.CLAUDE_ISOLATION_SETTINGS}, arms.CLAUDE_ISOLATION_SETTINGS)

    def test_model_is_a_parameter(self):
        for agent in arms.AGENTS:
            argv = arms.build_argv(agent, 'shell', model='m-x', prompt_text='p', cwd='/w')
            self.assertIn('m-x', argv)

    def test_world_arm_requires_its_config_and_shell_refuses_one(self):
        with self.assertRaises(arms.ArmError):
            arms.build_argv('claude', 'world', model='m', prompt_text='p', cwd='/v')
        with self.assertRaises(arms.ArmError):
            arms.build_argv('claude', 'shell', model='m', prompt_text='p', cwd='/w', mcp_config_path='/x')

    def test_no_token_in_any_argv_and_redaction(self):
        tok = 'a' * 64
        cfg = arms.claude_mcp_config(tok)
        self.assertEqual(cfg['mcpServers']['world']['headers']['Authorization'], 'Bearer ' + tok)
        for agent in arms.AGENTS:
            for arm in arms.ARMS:
                mcp = '/p/ep.mcp.json' if (agent, arm) == ('claude', 'world') else None
                argv = arms.build_argv(agent, arm, model='m', prompt_text='do it', cwd='/v', mcp_config_path=mcp)
                self.assertFalse(any(tok in a for a in argv))
                red = arms.redact_argv(argv, 'do it', (tok,))
                self.assertNotIn('do it', red)
                self.assertTrue(any(a.startswith('<prompt sha256:') for a in red))
        red = arms.redact_argv(['x', 'Authorization: Bearer ' + tok, 'y' + tok], 'p', (tok,))
        self.assertEqual(red, ['x', '<redacted>', '<redacted>'])


EIGHT_MCP = [f'mcp__world__{n}' for n in EIGHT]
SHELL_INIT = {'tools': ['Bash', 'Edit', 'Glob', 'Grep', 'Read', 'Write'], 'mcp_servers': []}
WORLD_INIT = {'tools': list(EIGHT_MCP), 'mcp_servers': [{'name': 'world', 'status': 'connected'}]}


class AttestFixtures(unittest.TestCase):
    def test_claude_shell_init(self):
        self.assertEqual(attest.claude_init(SHELL_INIT, 'shell'), [])
        self.assertTrue(attest.claude_init(dict(SHELL_INIT, tools=SHELL_INIT['tools'] + ['WebFetch']), 'shell'))
        self.assertTrue(attest.claude_init(dict(SHELL_INIT, tools=SHELL_INIT['tools'] + ['mcp__x__y']), 'shell'))
        self.assertTrue(attest.claude_init(dict(SHELL_INIT, tools=['Read']), 'shell'))
        self.assertTrue(attest.claude_init(dict(SHELL_INIT, mcp_servers=[{'name': 'eparse'}]), 'shell'))
        self.assertTrue(attest.claude_init(None, 'shell'))

    def test_claude_world_init(self):  # MUT-ALLOWEDTOOLS-ONLY's transcript: built-ins still there
        self.assertEqual(attest.claude_init(WORLD_INIT, 'world'), [])
        self.assertTrue(attest.claude_init(dict(WORLD_INIT, tools=['Bash', 'Read'] + EIGHT_MCP), 'world'))
        self.assertTrue(attest.claude_init(dict(WORLD_INIT, tools=EIGHT_MCP[:7]), 'world'))
        self.assertTrue(attest.claude_init(dict(WORLD_INIT, tools=EIGHT_MCP + ['mcp__world__workspace-exec']), 'world'))
        self.assertTrue(attest.claude_init(dict(WORLD_INIT, mcp_servers=[{'name': 'world', 'status': 'failed'}]), 'world'))
        two = [{'name': 'world', 'status': 'connected'}, {'name': 'obscura', 'status': 'connected'}]
        self.assertTrue(attest.claude_init(dict(WORLD_INIT, mcp_servers=two), 'world'))

    def test_codex_world_items(self):  # MUT-CODEX-SHELL's transcript: a command_execution item
        clean = {'command_execution': 0, 'file_change_ok': 0,
                 'tool_items': {'agent_message': 1, 'mcp_tool_call': 1}}
        self.assertEqual(attest.codex_world(clean, []), [])
        self.assertTrue(attest.codex_world(dict(clean, command_execution=1), []))
        self.assertTrue(attest.codex_world(dict(clean, file_change_ok=1), []))
        self.assertTrue(attest.codex_world(clean, ['note.txt']))
        self.assertTrue(attest.codex_world(dict(clean, tool_items={'web_search': 1}), []))


if __name__ == '__main__':
    unittest.main()
