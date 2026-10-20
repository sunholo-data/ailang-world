Changed only the two authorized files.

Sections changed: header/totals; decisions; environment/gates; controller baseline; milestone table; M5b annotation; new M5e; held M5c/M5d; M6-pre dependency; mutation boundaries; JSON statuses, gates, baselines and coverage.

M0a–M1b remain unchanged. Consistency and diff checks passed.

```sh
jq -e . design_docs/planned/sprint_w-exec-protected-test-paths.json > /dev/null; echo jq_rc=$?
# jq_rc=0

grep -c 'D-WORLD-86' design_docs/planned/w-exec-protected-test-paths-sprint-plan.md
# 13
```