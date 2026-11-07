#!/usr/bin/env python3
"""Fresh archive + independent Git root + effect assertions, never production apply."""
import datetime
import hashlib
import io
from pathlib import Path
import subprocess
import tarfile

root = Path(__file__).resolve().parents[3]
fresh = root/'.design-scratch/fresh'
patch = root/'design_docs/planned/w-daemon-startup-under-load.patch'
print('operation_deadline_utc=' + (datetime.datetime.now(datetime.timezone.utc)+datetime.timedelta(seconds=120)).isoformat(), flush=True)
def run(cmd, cwd=root):
    print('command=' + repr([str(x) for x in cmd]), flush=True)
    r = subprocess.run(cmd, cwd=cwd, capture_output=True, timeout=30, check=True)
    if r.stdout: print(r.stdout.decode(), end='', flush=True)
    if r.stderr: print(r.stderr.decode(), end='', flush=True)
    print('rc=0', flush=True)
    return r.stdout

baseline = run(['git','rev-parse','HEAD']).decode().strip()
print('baseline_sha=' + baseline, flush=True)
print('revised_patch_sha256=' + hashlib.sha256(patch.read_bytes()).hexdigest(), flush=True)
assert not fresh.exists()
fresh.mkdir(parents=True)
r = subprocess.run(['git','archive',baseline], cwd=root, capture_output=True, check=True, timeout=30)
with tarfile.open(fileobj=io.BytesIO(r.stdout)) as tar: tar.extractall(fresh, filter='data')
run(['git','init','-q'], fresh)
run(['git','apply','--check',str(patch)], fresh)
run(['git','apply',str(patch)], fresh)
for name in ['host/daemon/startup_under_load_test.go','host/broker/exec_probe_out_of_time_test.go','cmd/ailang-worldd/doctor_exec_profiles_test.go']:
    run(['test','-f',name], fresh)
    print('PRESENCE_ASSERTED=' + name, flush=True)
for line in patch.read_text().splitlines():
    if not line.startswith('+++ b/'): continue
    name = line[6:]
    got = (fresh/name).read_bytes()
    assert got == (root/'.design-scratch/proto'/name).read_bytes(), name
    print('applied_sha256=' + hashlib.sha256(got).hexdigest() + ' ' + name, flush=True)
run(['git','apply','--stat',str(patch)], fresh)
run(['git','apply','--numstat',str(patch)], fresh)
print('all patched files match the measured prototype byte-exact', flush=True)
