#!/usr/bin/env python3
"""PUB011-only ratchet. Input is pinned `pkg quality --json` stdout."""
import ast
import json
import os
import subprocess
import sys

# Lower this in the same change that contracts an exempt export; never raise it.
UNCONTRACTED_EXPORTS = 3


def require(condition, message):
    if not condition:
        raise AssertionError(message)


def check(data):
    require(data['schema'] == 'ailang.package-quality/v1', 'unknown quality schema')
    require(data['package'] == 'world/core', 'wrong package')
    require(data['compile']['ok'] is True, 'package did not compile')
    c = data['contracts']
    n = c['uncontracted_exports']
    require(type(n) is int and n >= 0, 'invalid uncontracted_exports')
    require(n == UNCONTRACTED_EXPORTS, (
        f'PUB011: uncontracted_exports={n}, ratchet={UNCONTRACTED_EXPORTS}; '
        'increases forbidden; decreases require lowering the ratchet'))
    require(c['verified'] > 0, 'zero proven contracts')
    require(c['errors'] == c['counterexample'] == c['skipped'] == 0, 'unproven contracts')
    return f'PUB011 ratchet passed: uncontracted_exports={n}'


def check_previous(source):
    tree = ast.parse(source)
    values = [n.value.value for n in tree.body if isinstance(n, ast.Assign)
              and any(isinstance(t, ast.Name) and t.id == 'UNCONTRACTED_EXPORTS' for t in n.targets)
              and isinstance(n.value, ast.Constant)]
    require(len(values) == 1 and type(values[0]) is int, 'invalid previous ratchet')
    require(0 <= UNCONTRACTED_EXPORTS <= values[0], 'ratchet increase forbidden')
    print(f'ratchet monotonic: {values[0]} -> {UNCONTRACTED_EXPORTS}')


def resolve_base(event_name, env):
    if event_name == 'pull_request':
        base = env.get('RATCHET_PR_BASE', '')
    elif event_name == 'push':
        base = env.get('RATCHET_PUSH_BEFORE', '')
    elif event_name == 'workflow_dispatch':
        print('contract quality gate: workflow_dispatch fallback to HEAD^', flush=True)
        return 'HEAD^'
    else:
        raise ValueError(f'unsupported ratchet event {event_name}')
    require(bool(base), f'missing ratchet base for {event_name}')
    return base


def check_base(ref):
    try:
        base = subprocess.check_output(
            ['git', 'rev-parse', '--verify', ref + '^{commit}'],
            text=True, stderr=subprocess.PIPE).strip()
    except subprocess.CalledProcessError as e:
        raise ValueError(f'cannot resolve ratchet base {ref}') from e
    path = 'scripts/check_world_contract_quality.py'
    names = subprocess.check_output(['git', 'ls-tree', '--name-only', base, '--', path], text=True)
    # Bootstrap only when the verified base commit predates this gate.
    source = subprocess.check_output(['git', 'show', base + ':' + path], text=True) if names.strip() else 'UNCONTRACTED_EXPORTS = 3'
    check_previous(source)


if __name__ == '__main__':
    try:
        if len(sys.argv) == 3 and sys.argv[1] == '--base':
            check_base(sys.argv[2])
        elif len(sys.argv) == 3 and sys.argv[1] == '--event':
            check_base(resolve_base(sys.argv[2], os.environ))
        else:
            require(len(sys.argv) == 1, 'usage: check_world_contract_quality.py [--base COMMIT | --event EVENT]')
            print(check(json.load(sys.stdin)))
    except (AssertionError, KeyError, ValueError, TypeError, subprocess.CalledProcessError) as e:
        sys.exit(f'contract quality gate: {e}')
