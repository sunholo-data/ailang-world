#!/usr/bin/env python3
"""Serial production mutations; each has a bounded gate, diff and restore proof."""
import hashlib
import json
from pathlib import Path
import subprocess
import sys

root = Path(__file__).resolve().parents[3]
tree = root / '.design-scratch/proto'
v = Path(__file__).resolve().parent
runner = v / 'run.py'
D = 'host/daemon/daemon.go'
B = 'host/broker/exec_probe.go'
S = 'host/broker/exec_startup.go'
W = 'host/daemon/workspace.go'
DOC = 'cmd/ailang-worldd/doctor.go'
tests = '^(TestExecProbeOutOfTimeNamesTheClockNotAnArm|TestExecProbeGenuineRefusalRecordedAfterDeadline|TestDoctorLateRefusalObservationAndReasonCap|TestRegistryPhasePreservesParentCancellation|TestExecExclusionLogIsObservationOnly|TestHealthExecExclusionIsAdditive|TestDaemonConstructorCallerInventory|TestDaemonCallerInventoryAliasControl)$'
M = [
 ('MUT-LATE-NEVER-SPENT', B, 'spent := !pp.deadline.IsZero() && !time.Now().Before(pp.deadline)', 'spent := false && !pp.deadline.IsZero() && !time.Now().Before(pp.deadline)', 'TestExecProbeGenuineRefusalRecordedAfterDeadline'),
 ('MUT-SPENT-ANY', B, 'if !f.spent {\n\t\t\treturn false\n\t\t}\n\t}\n\treturn len(e.Failed) > 0', 'if f.spent {\n\t\t\treturn true\n\t\t}\n\t}\n\treturn false', 'TestExecProbeOutOfTimeNamesTheClockNotAnArm'),
 ('MUT-LATE-STILL-SERVED', S, '\t\tdelete(st.Profiles, p.Project)', '\t\tif !spent { delete(st.Profiles, p.Project) }', 'TestExecProbeGenuineRefusalRecordedAfterDeadline'),
 ('MUT-FLAG-NOT-CARRIED', S, 'OutOfTime: spent})', 'OutOfTime: false && spent})', 'TestExecProbeGenuineRefusalRecordedAfterDeadline'),
 ('MUT-CAUSAL-OBSERVATION', B, 'Every failed arm was recorded at or after the startup deadline; startup timeout may have contributed.', 'The host was too busy; nothing to fix in the profile.', 'TestExecProbeGenuineRefusalRecordedAfterDeadline'),
 ('MUT-LOG-CAUSE', W, 'x.Project, broker.ExecProbeDeadlineObservation)', 'x.Project, "the host was too busy; nothing to fix in the profile")', 'TestExecExclusionLogIsObservationOnly'),
 ('MUT-DOCTOR-CAUSE', DOC, 'x.Project, x.File, broker.ExecProbeDeadlineObservation,', 'x.Project, x.File, "the host was too busy; nothing to fix in the profile",', 'TestDoctorLateRefusalObservationAndReasonCap'),
 ('MUT-DOCTOR-NO-REVIEW', DOC, 'Review the failure reasons and restart World if appropriate', 'Restart World when the host is quieter', 'TestDoctorLateRefusalObservationAndReasonCap'),
 ('MUT-REASON-CAP-240-PREFIX', DOC, '240-len("…")', '240', 'TestDoctorLateRefusalObservationAndReasonCap'),
 ('MUT-HEALTH-NO-OMITEMPTY', D, 'json:"out_of_time,omitempty"', 'json:"out_of_time"', 'TestHealthExecExclusionIsAdditive'),
 ('MUT-PARENT-CANCEL-DETACHED', D, 'context.WithTimeout(phaseParent, registryBudget)', 'context.WithTimeout(context.WithoutCancel(phaseParent), registryBudget)', 'TestRegistryPhasePreservesParentCancellation'),
 ('MUT-CALLER-IMPORTED-ALIAS', 'cmd/ailang-worldd/startup_caller_mutant.go', None, 'package main\nimport renamed "github.com/sunholo-data/ailang-world/host/daemon"\nvar unintendedConstructor = renamed.New\n', 'TestDaemonConstructorCallerInventory'),
 ('MUT-CALLER-DIRECT', 'cmd/ailang-worldd/startup_caller_mutant.go', None, 'package main\nimport renamed "github.com/sunholo-data/ailang-world/host/daemon"\nfunc unreviewedStartup() { renamed.New(nil, renamed.Config{}) }\n', 'TestDaemonConstructorCallerInventory'),
 ('MUT-CALLER-LOCAL-ALIAS', 'host/daemon/startup_caller_mutant.go', None, 'package daemon\nvar unintendedConstructor = New\n', 'TestDaemonConstructorCallerInventory'),
]

def gate(label):
    log = v / (label + '.txt')
    r = subprocess.run([sys.executable, str(runner), '180', str(log), str(tree), 'go', 'test', '-count=1', '-timeout', '120s', '-v', '-run', tests, './host/daemon/', './host/broker/', './cmd/ailang-worldd/'], timeout=185)
    return r.returncode, log.read_text()

rc, output = gate('mutation-CONTROL')
assert rc == 0 and output.count('--- PASS: Test') == 8, 'control red/missing; mutants void'
summary = []
selection = set(sys.argv[1].split(",")) if len(sys.argv) > 1 else None
prefix = "caller-final-" if selection else ""
for label, name, old, new, expected in M:
    if selection and label not in selection: continue
    p = tree/name
    before = p.read_bytes() if p.exists() else None
    try:
        if old is None:
            assert before is None
            p.write_text(new)
        else:
            s = before.decode()
            count = s.count(old)
            assert count == (2 if label == 'MUT-REASON-CAP-240-PREFIX' else 1), (label, count)
            p.write_text(s.replace(old, new))
        after = p.read_bytes()
        assert before != after
        rc, output = gate(prefix + 'mutation-' + label)
        match = rc == 1 and '--- FAIL: ' + expected in output and '[build failed]' not in output and '[setup failed]' not in output
        summary.append(dict(mutant=label, file=name, before=hashlib.sha256(before).hexdigest() if before else 'absent', after=hashlib.sha256(after).hexdigest(), rc=rc, assertion=expected, match=match))
    finally:
        if before is None: p.unlink()
        else: p.write_bytes(before)
        assert (p.read_bytes() if p.exists() else None) == before
    (v/(prefix+'mutations.json')).write_text(json.dumps(summary, indent=2)+'\n')
assert all(row['match'] for row in summary), 'inspect mutation mismatches'
print(f'{len(summary)} named mutations killed; all files restored byte-exact')
