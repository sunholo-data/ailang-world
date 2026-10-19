#!/usr/bin/env python3
"""Measure Unicode Git pathspec misses; exit 0 means valid evidence, not equivalence."""
import base64
import datetime
import json
import os
from pathlib import Path
import platform
import subprocess
import shutil
import tempfile

ROOT = Path(__file__).resolve().parents[3]
OUT = Path(__file__).with_name('result.json')
# Source literals, checked manually against gitPathspecs in exec_protected.go:327-337.
DIRECTORY = [':(top,glob,icase)tests', ':(top,glob,icase)tests/**']
BASENAME = [':(top,glob,icase)**/*k_test.go']
ROWS = [('ascii-directory', 'tests/x.py', 'tests/**', DIRECTORY),
        ('long-s', 'teſts/x.py', 'tests/**', DIRECTORY),
        ('ascii-k', 'pkg/k_test.go', '*k_test.go', BASENAME),
        ('kelvin', 'pkg/K_test.go', '*k_test.go', BASENAME)]


def blob(data):
    return {'utf8': data.decode('utf-8', 'backslashreplace'),
            'base64': base64.b64encode(data).decode('ascii'),
            'nul_fields_base64': [base64.b64encode(p).decode('ascii')
                                  for p in data.split(b'\0') if p]}


def main():
    OUT.unlink(missing_ok=True)
    go_binary = shutil.which('go')
    if not go_binary:
        raise RuntimeError('go is required')
    go_env = json.loads(subprocess.run([go_binary, 'env', '-json', 'GOMODCACHE', 'GOCACHE'], cwd=ROOT, capture_output=True, check=True, timeout=10).stdout)
    result = {'measurement_valid': False, 'equivalence': False,
              'M1c': 'incomplete', 'platform': platform.platform(),
              'utc': datetime.datetime.now(datetime.timezone.utc).isoformat(),
              'pathspec_origin': 'source literal used; manually checked against current gitPathspecs, host/broker/exec_protected.go:327-337; no source parser',
              'timeout_seconds_git': 10, 'timeout_seconds_go_test': 120, 'rows': []}
    with tempfile.TemporaryDirectory(prefix='.iter281-evidence-', dir=ROOT) as scratch:
        home = Path(scratch) / 'home'
        home.mkdir()
        env = {k: v for k, v in os.environ.items()
               if not k.startswith('GIT_') and k != 'WORLD_TOOL_AILANG_BIN'}
        env.update(HOME=str(home), XDG_CONFIG_HOME=str(home / '.config'),
                   TMPDIR=scratch, GOTMPDIR=scratch, GOCACHE=go_env['GOCACHE'],
                   GIT_CONFIG_NOSYSTEM='1', GIT_CONFIG_GLOBAL=os.devnull,
                   GIT_OPTIONAL_LOCKS='0', LANG='C', LC_ALL='C')
        # Use the already provisioned module cache read-only; no dependency downloads.
        env['GOPROXY'] = 'off'
        env['GOMODCACHE'] = go_env['GOMODCACHE']

        def run(argv, cwd, timeout=10):
            p = subprocess.run(argv, cwd=cwd, env=env, capture_output=True, timeout=timeout)
            item = {'argv': argv, 'cwd': str(cwd), 'rc': p.returncode, 'timeout_seconds': timeout,
                    'stdout': blob(p.stdout), 'stderr': blob(p.stderr)}
            if p.returncode:
                raise RuntimeError(json.dumps(item))
            return item, p.stdout

        result['git_version'] = run(['git', '--version'], ROOT)[0]
        result['base'] = run(['git', 'rev-parse', 'HEAD'], ROOT)[0]
        result['environment'] = {k: env[k] for k in ['HOME', 'TMPDIR', 'GOTMPDIR', 'GIT_CONFIG_GLOBAL', 'GIT_CONFIG_NOSYSTEM', 'GIT_OPTIONAL_LOCKS', 'LANG', 'LC_ALL', 'GOCACHE', 'GOMODCACHE', 'GOPROXY']}
        for name, path, pattern, specs in ROWS:
            repo = Path(scratch) / name
            repo.mkdir()
            row = {'name': name, 'path': path, 'path_utf8_base64': base64.b64encode(path.encode()).decode(),
                   'pattern': pattern, 'pathspecs': specs, 'mirror_expected': True, 'setup': []}
            def git(*args):
                return run(['git', '-c', 'core.fsmonitor=false', *args], repo)
            for args in [('init', '-q'), ('config', 'core.ignorecase', 'false'),
                         ('config', 'user.name', 'Evidence'), ('config', 'user.email', 'evidence@example.invalid')]:
                row['setup'].append(git(*args)[0])
            target = repo / path
            target.parent.mkdir(parents=True, exist_ok=True)
            target.write_bytes(b'before\n')
            row['setup'].append(git('add', '--', path)[0])
            row['setup'].append(git('-c', 'core.hooksPath=/dev/null', 'commit', '-qm', 'fixture')[0])
            target.write_bytes(b'after\n')
            outputs = {}
            for label, args in [('unfiltered_ls_files', ['ls-files', '-z']),
                                ('unfiltered_status', ['status', '--porcelain=v1', '-z', '--untracked-files=all', '--ignored=matching']),
                                ('selected_ls_files', ['ls-files', '-z', '--', *specs]),
                                ('selected_status', ['status', '--porcelain=v1', '-z', '--untracked-files=all', '--ignored=matching', '--', *specs])]:
                row[label], outputs[label] = git(*args)
            encoded = path.encode()
            if outputs['unfiltered_ls_files'] != encoded + b'\0' or outputs['unfiltered_status'] != b' M ' + encoded + b'\0':
                raise RuntimeError('instrument failure: missing exact unfiltered path/status control: ' + name)
            row['git_selected'] = (outputs['selected_ls_files'] == encoded + b'\0' and
                                   outputs['selected_status'] == b' M ' + encoded + b'\0')
            if name.startswith('ascii') and not row['git_selected']:
                raise RuntimeError('instrument failure: ASCII selected positive control missing: ' + name)
            row['equivalence'] = row['git_selected'] == row['mirror_expected']
            result['rows'].append(row)
        result['mirror_test'], output = run([go_binary, 'test', '-count=1', './host/broker/...', '-run',
            '^TestProtectedMatchFixture$/(long-s|kelvin)$', '-timeout', '110s', '-v'], ROOT, timeout=120)
        for name in ['long-s', 'kelvin']:
            if ('--- PASS: TestProtectedMatchFixture/' + name).encode() not in output:
                raise RuntimeError('instrument failure: mirror named PASS missing: ' + name)
        result['equivalence'] = all(row['equivalence'] for row in result['rows'])
        result['measurement_valid'] = True
    OUT.write_text(json.dumps(result, ensure_ascii=True, indent=2) + '\n')
    print(json.dumps({'measurement_valid': True, 'equivalence': result['equivalence'],
                      'M1c': 'incomplete', 'output': str(OUT)}))


if __name__ == '__main__':
    main()
