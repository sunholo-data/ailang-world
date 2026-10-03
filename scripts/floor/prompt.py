#!/usr/bin/env python3
"""Row 93 floor harness — prompt rendering (design §4.2, V16, V34; AC1.4).

One task message per (agent, arm, task): the teaching prompt (``$TOOL prompt``, V34) followed
by the ailang agent task template ``internal/eval_harness/templates/agent_task_ailang.txt`` at
the pinned corpus commit, filled exactly as ``GenerateAgentPromptsWithSystemPrompt`` fills it
under the default (``baseline``) condition: no contract spec, no Z3 hints, no tool guidance,
the default verification steps, no extra tools. Both agents get the same message (codex has no
system-prompt flag; one mechanism keeps them comparable).

The two arms of one agent differ in EXACTLY two places (§4.2):
  (a) the TOOLS BLOCK — the template's two tool-specific slots: the ``## Available Tools``
      command list (+ ``{{EXTRA_TOOLS}}``) and the ``{{VERIFICATION_STEPS}}`` body, which in
      the Go default is itself a shell ``ailang run`` command. Shell keeps ``ailang run/check/
      test``; World lists its 8 tools and says plainly what ``ailang-run`` cannot do;
  (b) ``{{SOLUTION_PATH}}`` — the absolute worktree path (shell) vs ``benchmark/solution.ail``
      (World).
Everything else — including the teaching prompt, which itself mentions ``ailang run`` — is
byte-identical across arms.

Two deliberate departures from the Go agent path, identical in both arms:
  * the task text is ``task_prompt``, falling back to ``prompt`` (the Go agent path uses
    ``task_prompt`` only, so ``cli_args``/``pipeline``/``float_eq`` reach a Go agent with no task
    text at all; standard mode falls back exactly like this — spec.go PromptPartsForLanguage);
  * the shell arm's run command carries ``--net-allow-http --net-allow-localhost`` for a spec
    with ``net_allow_localhost`` (the Go template loader's M-EVAL-NETWORK-MOCK-FIXTURE rule), so
    the shell agent can reach the live mock (V8) when it tests ``api_call_json``.
"""
from __future__ import annotations

import hashlib
import os
import re
import subprocess

TEMPLATE_PATH = 'internal/eval_harness/templates/agent_task_ailang.txt'
MOCK_TOKEN = '{{MOCK_HTTP_URL}}'
WORLD_SOLUTION_PATH = 'benchmark/solution.ail'
ARMS = ('shell', 'world')

# Row 134's se-tools ids (V29), in transitions.json order.
WORLD_TOOLS = (
    ('ailang-read', 'read a file in your worktree (`path` relative to the worktree root)'),
    ('ailang-write', 'write a file (`path`, `content`)'),
    ('ailang-edit', 'replace text in a file (`path`, `old_text`, `new_text`)'),
    ('ailang-check', 'type-check a file (`path`)'),
    ('ailang-run', "run a file's `main` under World's policy (`path`, optional `args_json`)"),
    ('builtins-search', 'search the AILANG builtins (`query`, optional `module`)'),
    ('examples-search', 'search AILANG examples (`query`)'),
    ('ailang-cli', 'policy-allowlisted AILANG CLI operations (`op`, ...)'),
)
# P4 (§1): today's ailang-run has no stdin, no argv and only IO/FS. Row 135 (D-NF-3 = B) widens
# it; this sentence is a recorded tuning knob, applied before the prereg is committed.
WORLD_RUN_NOTE = ('`ailang-run` cannot pipe stdin, cannot pass command-line arguments and grants '
                  'no capabilities beyond IO and FS. If the task needs any of those, it cannot be '
                  'run end to end here: use `ailang-check` and reason carefully about the output. '
                  'Your solution is graded with the task\'s own stdin, arguments and capabilities.')

SHELL_TOOLS_TMPL = ('You have access to the `ailang` command:\n'
                    '- **Run code:** `ailang run --entry main --caps {caps} solution.ail`\n'
                    '- **Type check:** `ailang check solution.ail`\n'
                    '- **Run tests:** `ailang test solution.ail`\n')

# getDefaultVerificationSteps(caps), agent_prompt.go — verbatim
SHELL_STEPS_TMPL = ('1. **Run your solution:**\n'
                    '   ```\n'
                    '   ailang run --entry main --caps {caps} benchmark/solution.ail\n'
                    '   ```\n'
                    '\n'
                    '2. **Compare output carefully:**\n'
                    '   - Check that every line matches expected output\n'
                    '   - No extra spaces or formatting differences\n'
                    '\n'
                    '3. **If output doesn\'t match:**\n'
                    '   - Fix your code\n'
                    '   - Re-run to verify\n'
                    '   - Repeat until output matches exactly\n'
                    '\n'
                    '4. **Only finish once verified!**')

WORLD_STEPS_TMPL = ('1. **Run your solution:**\n'
                    '   call the `ailang-run` tool with `{{"path": "benchmark/solution.ail"}}`\n'
                    '\n'
                    '2. **Compare output carefully:**\n'
                    '   - Check that every line matches expected output\n'
                    '   - No extra spaces or formatting differences\n'
                    '\n'
                    '3. **If output doesn\'t match:**\n'
                    '   - Fix your code\n'
                    '   - Re-run to verify\n'
                    '   - Repeat until output matches exactly\n'
                    '\n'
                    '4. **Only finish once verified!**')

AVAILABLE_TOOLS_HEADING = '## Available Tools\n\n'
NEXT_HEADING = '\n## Common Pitfalls'


class PromptError(Exception):
    pass


def sha256(s: str) -> str:
    return hashlib.sha256(s.encode()).hexdigest()


def load_template(repo: str, commit: str) -> str:
    p = subprocess.run(['git', '-C', repo, 'show', f'{commit}:{TEMPLATE_PATH}'], capture_output=True)
    if p.returncode != 0:
        raise PromptError(f'cannot read {TEMPLATE_PATH} at {commit}: {p.stderr.decode(errors="replace")}')
    return p.stdout.decode('utf-8')


def load_teaching(tool_bin: str, compact: bool = False) -> str:
    """The teaching prompt from the pinned tool binary (V34)."""
    if not os.path.isabs(tool_bin):
        raise PromptError('teaching prompt must come from an explicit tool binary path')
    argv = [tool_bin, 'prompt'] + (['--compact'] if compact else [])
    p = subprocess.run(argv, capture_output=True)
    if p.returncode != 0 or not p.stdout.strip():
        raise PromptError(f'{argv} failed: rc={p.returncode}')
    return p.stdout.decode('utf-8')


def world_tools_block(note: str = WORLD_RUN_NOTE, tools=WORLD_TOOLS) -> str:
    lines = ['You have NO shell and NO built-in file tools. You work only through these World '
             'tools (MCP server `world`); every path is relative to your worktree root:']
    lines += [f'- `{name}` — {desc}' for name, desc in tools]
    lines.append('')
    lines.append(note)
    return '\n'.join(lines) + '\n'


def _caps_value(spec: dict) -> str:
    v = ','.join(spec.get('caps') or [])
    if spec.get('net_allow_localhost'):
        v += ' --net-allow-http --net-allow-localhost'
    return v


def render_task(spec: dict, arm: str, *, template: str, solution_path: str, deadline_s: int,
                mock_url: str | None = None, world_note: str = WORLD_RUN_NOTE) -> str:
    if arm not in ARMS:
        raise PromptError(f'unknown arm {arm!r}')
    if arm == 'world' and solution_path != WORLD_SOLUTION_PATH:
        raise PromptError(f'the World arm solution path is {WORLD_SOLUTION_PATH!r}')
    if arm == 'shell' and not os.path.isabs(solution_path):
        raise PromptError('the shell arm solution path must be absolute')
    a = template.find(AVAILABLE_TOOLS_HEADING)
    b = template.find(NEXT_HEADING, a)
    if a < 0 or b < 0 or template.count(AVAILABLE_TOOLS_HEADING) != 1:
        raise PromptError('template has no unique "## Available Tools" … "## Common Pitfalls" section')
    if '{{VERIFICATION_STEPS}}' not in template:
        raise PromptError('template has no {{VERIFICATION_STEPS}} slot')
    caps = ','.join(spec.get('caps') or [])
    if arm == 'shell':
        tools = SHELL_TOOLS_TMPL.format(caps=_caps_value(spec))
        steps = SHELL_STEPS_TMPL.format(caps=caps)
    else:
        tools = world_tools_block(world_note)
        steps = WORLD_STEPS_TMPL.format()
    t = template[:a + len(AVAILABLE_TOOLS_HEADING)] + tools + template[b:]
    task_text = spec.get('task_prompt') or spec.get('prompt') or ''
    if MOCK_TOKEN in task_text:
        if not mock_url:
            raise PromptError(f'task {spec.get("id")!r} needs a live mock URL')
        task_text = task_text.replace(MOCK_TOKEN, mock_url)
    t = t.replace('{{DESCRIPTION}}', (spec.get('description') or '') + '\n\n' + task_text)
    t = t.replace('{{EXPECTED_OUTPUT}}', spec.get('expected_stdout') or '')
    t = t.replace('{{CAPS}}', caps)
    t = t.replace('{{TIMEOUT}}', str(int(deadline_s)))
    t = t.replace('{{SOLUTION_PATH}}', solution_path)
    t = t.replace('{{PYTHON_VERSION}}', '')
    t = t.replace('{{CONTRACT_SPEC}}', '')
    t = t.replace('{{Z3_HINTS}}', '')
    t = t.replace('{{TOOL_GUIDANCE}}', '')
    t = t.replace('{{VERIFICATION_STEPS}}', steps)
    t = t.replace('{{EXTRA_TOOLS}}', '')
    t = t.replace('<LANG>', 'AILANG')
    left = re.findall(r'\{\{[A-Z_]+\}\}', t)
    if left:
        raise PromptError(f'unfilled placeholder(s) in rendered prompt: {sorted(set(left))}')
    return t


def render_message(teaching: str, task_text: str) -> str:
    """The single task message both agents receive: teaching prompt, then the task (§4.2)."""
    if not teaching.strip():
        raise PromptError('empty teaching prompt')
    return teaching.rstrip('\n') + '\n\n---\n\n' + task_text
