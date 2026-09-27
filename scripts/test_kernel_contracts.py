#!/usr/bin/env python3
"""Iteration 198: real compiler controls, scratch confined to this worktree."""
import copy
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import tempfile
from check_world_contract_quality import check, check_previous
from test_world_contract_history import controls as history_controls

ROOT = Path(__file__).resolve().parents[1]
SCRATCH = ROOT / '.scratch'
SCRATCH.mkdir(exist_ok=True)
BIN = os.environ['AILANG_BIN']
assert subprocess.check_output([BIN, '--version'], text=True).splitlines()[0] == 'AILANG v0.41.0'
FUNCTIONS = {'contracts': ['proposalMatchesWorld', 'verificationMatchesProposal', 'commitAllowed'],
             'transitions': ['plan', 'verify']}


def run(base, module, function, expected):
    cmd = [BIN, 'ai-check', f'world/{module}.ail']
    p = subprocess.run(cmd, cwd=base, capture_output=True, text=True, timeout=30)
    d = json.loads(p.stdout)
    assert d['check']['passed'], d
    results = {r['function']: r for r in d['verify']['results']}
    assert results[function]['status'] == expected, results
    assert any(r['status'] == 'verified' for n, r in results.items() if n != function), results
    print(f'{function}: {expected}; rc={p.returncode}; control=other identity verified',
          results[function].get('reason', ''))


with tempfile.TemporaryDirectory(prefix='kernel-contracts-', dir=SCRATCH) as td:
    base = Path(td)
    shutil.copytree(ROOT / 'world', base / 'world')
    backups = base / 'backups'
    shutil.copytree(base / 'world', backups)

    def restore(path):
        saved = backups / path.relative_to(base / 'world')
        shutil.copyfile(saved, path)
        assert hashlib.sha256(path.read_bytes()).digest() == hashlib.sha256(saved.read_bytes()).digest()

    for module, names in FUNCTIONS.items():
        path = base / 'world' / (module + '.ail')
        original = path.read_text()
        for fn in names:
            run(base, module, fn, 'verified')
            start = original.index('export func ' + fn + '(')
            a = original.index('ensures {', start)
            b = original.index('\n', a)
            path.write_text(original[:a] + 'ensures { false }' + original[b:])
            run(base, module, fn, 'counterexample')
            restore(path)
    # Actual implementation bugs, not merely false specifications.
    mutations = [
        ('contracts', 'proposalMatchesWorld', 'sameRef(w.stateRoot, p.inputWorld)', 'w.stateRoot.digest == p.inputWorld.digest'),
        ('contracts', 'verificationMatchesProposal', 'sameRef(p.proposalHash, v.proposalHash)', 'p.proposalHash.digest == v.proposalHash.digest'),
        ('contracts', 'commitAllowed', '    && v.accepted', '    && true'),
        ('transitions', 'plan', 'inputWorld: w.stateRoot', 'inputWorld: w.logHead'),
        ('transitions', 'verify', 'accepted: false', 'accepted: true'),
    ]
    for module, fn, old, new in mutations:
        path = base / 'world' / (module + '.ail'); original = path.read_text()
        assert old in original
        path.write_text(original.replace(old, new))
        run(base, module, fn, 'counterexample')
        restore(path)
    probes = [
        ('logepoch', 'renderRef', 'string', 'result == "${r.algo}:${r.digest}"', 'error'),
        ('logepoch', 'cacheKey', 'string', 'result == "${h.transitionFn.algo}:${h.transitionFn.digest}@${h.interpreter.algo}:${h.interpreter.digest}"', 'skipped'),
        ('transitions', 'commit', 'CommitResult', 'match result { Applied(next, t) => v.accepted && w.stateRoot == p.inputWorld && p.proposalHash == v.proposalHash && next.revision == w.revision + 1 && next.stateRoot == outputWorld && next.logHead == nextLogHead && t.proposalHash == p.proposalHash && t.inputWorld == p.inputWorld && t.outputWorld == outputWorld && t.transitionFn == p.transitionFn, Denied(reason) => not(v.accepted && w.stateRoot == p.inputWorld && p.proposalHash == v.proposalHash) }', 'skipped'),
    ]
    for module, fn, typ, law, verdict in probes:
        path = base / 'world' / (module + '.ail'); original = path.read_text()
        for expr in [law, 'false']:
            spec = '\n' + ('requires { w.revision >= 0 }\n' if fn == 'commit' else '') + 'ensures { ' + expr + ' }\n'
            path.write_text(re.sub(r'(export func '+fn+r'\([\s\S]*?\) -> '+typ+')', r'\1 ! {}'+spec, original, count=1))
            run(base, module, fn, 'counterexample' if fn == 'renderRef' and expr == 'false' else verdict)
        restore(path)
    path = base / 'world/logepoch.ail'; original = path.read_text()
    path.write_text(original.replace('module world/logepoch', 'module world/logepoch\nimport std/string (concat)').replace(
        'export func renderRef(r: HashRef) -> string',
        'export func renderRef(r: HashRef) -> string ! {}\nensures { result == concat([r.algo, ":", r.digest]) }'))
    run(base, 'logepoch', 'renderRef', 'error')
    restore(path)
    shutil.copytree(ROOT / 'packages/world-core', base / 'package')
    def quality():
        p = subprocess.run([BIN, 'pkg', 'quality', '--json', str(base / 'package')], capture_output=True, text=True, timeout=120)
        assert p.returncode == 0, p.stderr
        return json.loads(p.stdout)
    good = quality(); print(check(good))
    path = base / 'package/world/logepoch.ail'
    package_backup = base / 'package-logepoch.backup'
    shutil.copyfile(path, package_backup)
    path.write_text(path.read_text() + '\nexport func ratchetMutant(x: int) -> int ! {} { x }\n')
    bad = quality()
    assert bad['contracts']['uncontracted_exports'] == 4, bad
    try:
        check(bad)
    except AssertionError as e:
        print('added uncontracted export killed:', e)
    else:
        raise AssertionError('ratchet accepted added export')
    shutil.copyfile(package_backup, path)
    assert hashlib.sha256(path.read_bytes()).digest() == hashlib.sha256(package_backup.read_bytes()).digest()
    for label, mutate in [('missing field', lambda d: d['contracts'].pop('uncontracted_exports')),
                          ('zero proofs', lambda d: d['contracts'].update(verified=0)),
                          ('lower count needs ratchet', lambda d: d['contracts'].update(uncontracted_exports=2)),
                          ('contract errors', lambda d: d['contracts'].update(errors=1)),
                          ('contract counterexample', lambda d: d['contracts'].update(counterexample=1)),
                          ('contract skipped', lambda d: d['contracts'].update(skipped=1)),
                          ('compile failed', lambda d: d['compile'].update(ok=False)),
                          ('wrong package', lambda d: d.update(package='wrong/package')),
                          ('wrong schema', lambda d: d.update(schema='wrong/schema'))]:
        d = copy.deepcopy(good); mutate(d)
        try:
            check(d)
        except (AssertionError, KeyError) as e:
            print(label, 'killed:', e)
        else:
            raise AssertionError(label)
print('kernel contract controls PASSED')

check_previous('UNCONTRACTED_EXPORTS = 3')
try:
    check_previous('UNCONTRACTED_EXPORTS = 2')
except AssertionError as e:
    print('raised baseline killed:', e)
else:
    raise AssertionError('ratchet increase accepted')

history_controls()
