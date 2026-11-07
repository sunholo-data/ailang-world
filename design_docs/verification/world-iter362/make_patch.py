#!/usr/bin/env python3
"""Read-only comparison of scratch trees; no Git index or commit required."""
import difflib
import hashlib
import json
from pathlib import Path

root = Path(__file__).resolve().parents[3]
base, proto = (root/'.design-scratch'/n for n in ['base','proto'])
names = [
    'cmd/ailang-worldd/doctor.go', 'cmd/ailang-worldd/doctor_exec_profiles_test.go',
    'cmd/ailang-worldd/doctor_hosts.go', 'docs/QUICKSTART.md',
    'host/broker/exec_probe.go', 'host/broker/exec_probe_out_of_time_test.go',
    'host/broker/exec_startup.go', 'host/daemon/daemon.go',
    'host/daemon/startup_under_load_test.go', 'host/daemon/startup_callers_test.go',
    'host/daemon/workspace.go',
]
parts, manifest = [], []
for name in sorted(names):
    p = proto/name
    old = (base/name).read_text() if (base/name).exists() else ''
    new = p.read_text()
    assert old != new
    header = 'diff --git a/'+name+' b/'+name+'\n'
    if not (base/name).exists(): header += 'new file mode 100644\n'
    parts.append(header + ''.join(difflib.unified_diff(old.splitlines(True), new.splitlines(True), fromfile='a/'+name if (base/name).exists() else '/dev/null', tofile='b/'+name)))
    manifest.append(dict(path=name, sha256=hashlib.sha256(p.read_bytes()).hexdigest(), lines=len(new.splitlines())))
patch = root/'design_docs/planned/w-daemon-startup-under-load.patch'
patch.write_text(''.join(parts))
(Path(__file__).parent/'prototype-manifest.json').write_text(json.dumps(manifest, indent=2)+'\n')
print('patch_sha256='+hashlib.sha256(patch.read_bytes()).hexdigest())
print('patch_lines='+str(len(patch.read_text().splitlines())))
