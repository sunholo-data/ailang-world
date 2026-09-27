#!/usr/bin/env python3
"""Read-only history controls; no synthetic commits or external writes."""
from contextlib import redirect_stdout
from io import StringIO
from pathlib import Path
import subprocess
import sys

from check_world_contract_quality import resolve_base

ROOT = Path(__file__).resolve().parents[1]
CHECKER = ROOT / 'scripts/check_world_contract_quality.py'


def controls():
    # Simulate the workflow's extraction, with competing fields to catch mixups.
    event = {'pull_request': {'base': {'sha': 'a' * 40}}, 'before': 'b' * 40}
    env = {'RATCHET_PR_BASE': event['pull_request']['base']['sha'],
           'RATCHET_PUSH_BEFORE': event['before']}
    for name, expected in [('pull_request', 'a' * 40), ('push', 'b' * 40),
                           ('workflow_dispatch', 'HEAD^')]:
        out = StringIO()
        with redirect_stdout(out):
            actual = resolve_base(name, env)
        assert actual == expected, (name, actual)
        if name == 'workflow_dispatch':
            assert out.getvalue() == 'contract quality gate: workflow_dispatch fallback to HEAD^\n'
        else:
            assert out.getvalue() == '', out.getvalue()
        print(f'{name}: base={actual}; PASS')
        if out.getvalue():
            print(out.getvalue(), end='')

    workflow = (ROOT / '.github/workflows/ci.yml').read_text()
    for line in ['RATCHET_EVENT: ${{ github.event_name }}',
                 'RATCHET_PR_BASE: ${{ github.event.pull_request.base.sha }}',
                 'RATCHET_PUSH_BEFORE: ${{ github.event.before }}',
                 'run: python3 scripts/check_world_contract_quality.py --event "$RATCHET_EVENT"']:
        assert line in workflow, line
    print('CI event/base wiring: PASS')

    for ref in ['0' * 40, 'nosuchref']:
        p = subprocess.run([sys.executable, str(CHECKER), '--base', ref],
                           cwd=ROOT, capture_output=True, text=True)
        expected = f'contract quality gate: cannot resolve ratchet base {ref}'
        assert p.returncode == 1, p
        assert p.stderr.strip() == expected, p.stderr
        assert not p.stdout and 'Traceback' not in p.stderr, p
        print(f'invalid base: rc={p.returncode}; {p.stderr.strip()}; no traceback; PASS')

    # An actual root commit predates the checker; assert that precondition.
    base = subprocess.check_output(['git', 'rev-list', '--max-parents=0', 'HEAD'],
                                   cwd=ROOT, text=True).splitlines()[0]
    names = subprocess.check_output(['git', 'ls-tree', '--name-only', base, '--',
                                     'scripts/check_world_contract_quality.py'],
                                    cwd=ROOT, text=True)
    assert not names.strip(), names
    p = subprocess.run([sys.executable, str(CHECKER), '--base', base],
                       cwd=ROOT, capture_output=True, text=True)
    assert p.returncode == 0 and not p.stderr, p
    assert p.stdout == 'ratchet monotonic: 3 -> 3\n', p.stdout
    print(f'pre-checker base: {p.stdout.strip()}; rc=0; PASS')
    print('world contract history controls PASSED')


if __name__ == '__main__':
    controls()
