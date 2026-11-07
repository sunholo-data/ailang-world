#!/usr/bin/env python3
"""Bound each evidence command, preserve its rc/output, and observe gate load.

Usage: python3 run.py SECONDS LOG CWD COMMAND [ARGS...]
No TMPDIR/GOTMPDIR/HOME overrides; no credential environment output.
"""
import datetime
import os
from pathlib import Path
import shlex
import signal
import subprocess
import sys

seconds, log, cwd, *cmd = sys.argv[1:]
seconds = int(seconds)
assert 0 < seconds <= 1800
env = dict(os.environ)
env['PATH'] = '/opt/homebrew/bin:' + env['PATH']
env['AILANG_BIN'] = os.path.expanduser('~/.pinned-ailang/ailang')
now = datetime.datetime.now(datetime.timezone.utc)
with Path(log).open('w') as out:
    def say(s):
        out.write(s + '\n')
        out.flush()
    say('start_utc=' + now.isoformat())
    say('deadline_utc=' + (now + datetime.timedelta(seconds=seconds)).isoformat())
    say('cwd=' + str(Path(cwd).resolve()))
    say('command=' + shlex.join(cmd))
    say('PATH=/opt/homebrew/bin:$PATH AILANG_BIN=$HOME/.pinned-ailang/ailang; TMPDIR/GOTMPDIR untouched')
    for obs in [['uptime'], ['go', 'version']]:
        r = subprocess.run(obs, env=env, text=True, capture_output=True, timeout=10)
        say(r.stdout + r.stderr)
    p = subprocess.Popen(cmd, cwd=cwd, env=env, stdout=out, stderr=subprocess.STDOUT, start_new_session=True)
    try:
        rc = p.wait(timeout=seconds-25)
    except subprocess.TimeoutExpired:
        os.killpg(p.pid, signal.SIGKILL)
        p.wait(timeout=10)
        rc = 124
        say('TIMEOUT: output VOID as a completed gate')
    say('rc=' + str(rc))
    say('end_utc=' + datetime.datetime.now(datetime.timezone.utc).isoformat())
    r = subprocess.run(['uptime'], text=True, capture_output=True, timeout=10)
    say(r.stdout + r.stderr)
print(str(log) + ': rc=' + str(rc), flush=True)
sys.exit(rc)
