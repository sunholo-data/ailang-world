#!/usr/bin/env python3
"""Bank actual source for reviewer-requested premises, not spelling-only absences."""
from pathlib import Path

root = Path(__file__).resolve().parents[3]
base = root/'.design-scratch/base'
fresh = root/'.design-scratch/fresh'
def show(tree, name, first, last):
    print(f'SOURCE {tree.name}/{name}:{first}-{last}')
    for line, text in enumerate((tree/name).read_text().splitlines(), 1):
        if first <= line <= last: print(f'{line}: {text}')

show(base, 'host/daemon/daemon.go', 140, 154) # commitBudget declaration + purpose
show(base, 'host/daemon/daemon.go', 607, 652) # additive HealthResponse contract
show(base, 'host/daemon/daemon.go', 1271, 1289) # sole baseline production construction path
show(base, 'host/daemon/daemon_test.go', 205, 245) # existing constant pin, including B6
show(base, 'cmd/ailang-worldd/mcpclient.go', 659, 666) # cap helper's ellipsis bytes
show(fresh, 'cmd/ailang-worldd/doctor.go', 348, 368) # corrected local cap and advice
show(base, 'host/daemon/context_thread_test.go', 1, 43) # actual old cancellation assertions

for tree in [base, fresh]:
    names = ['host/daemon/startup_under_load_test.go','host/broker/exec_probe_out_of_time_test.go','cmd/ailang-worldd/doctor_exec_profiles_test.go','host/daemon/startup_callers_test.go']
    print(f'{tree.name} new test presence=' + repr({n:(tree/n).is_file() for n in names}))
    # This text count is only an instrument-health check. Race tests/mutations
    # carry the dynamic claims, not this search of our own test text.
    for name in names:
        if (tree/name).exists():
            text = (tree/name).read_text()
            print(f'{tree.name}/{name}: t.Parallel occurrences={text.count("t.Parallel(")}; t.Run control={text.count("t.Run(")}')

patch = root/'design_docs/planned/w-daemon-startup-under-load.patch'
paths = [line[6:] for line in patch.read_text().splitlines() if line.startswith('+++ b/')]
assert len(paths) == 11
assert sum(p.endswith('.go') for p in paths) == 10
assert all(p.endswith('.go') or p == 'docs/QUICKSTART.md' for p in paths)
for p in paths: print('PATCH_PATH=' + p)
print('patch scope: 10 Go files, QUICKSTART only; no .ail, frozen machinery, row255 record/handler, row257 seed files')
