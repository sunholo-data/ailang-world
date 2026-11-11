# Iteration 382 — controller gates (row 201 A2 sprint 4, head `e07b7bad`)

Base `27ce163e`. Sprint branch `sprint/world-iter382-row201-a2s4`, PR #552.

## Round trip (controller, independent)
`git archive 27ce163e` → the 7 Appendix-P patches extracted from the plan file with §2's awk → `patch -p1`:
`diff -rq` against the planner's `plan-scratch/work/` tree empty for `cmd`, `host` and `docs`.
After the last sub-milestone commit the sprint worktree's `cmd`/`host`/`docs` are also `diff -rq`-identical to it.

## Sub-milestone commits (commit_ms.sh; numstat `git diff --numstat 27ce163e HEAD -- cmd host`)
| M | commit | numstat | plan |
|---|---|---|---|
| M7a | e4734707 | 11 | 11 |
| M7b | a7f45088 | 195 | 195 |
| M7c | 30c5debe | 367 | 367 |
| M7d | c1483836 | 597 | 597 |
| M7e | 17c21dae | 924 | 924 |
| M7f | e8077ba6 | 1054 | 1054 |
| M7g | e07b7bad | 1054 | 1054 |

## Plan §6 pre-commit script (gates.sh) at e07b7bad
rc 0, `PRE_COMMIT_GATES_GREEN`: verify_ail, gofmt, build, vet, compile fence, touched packages (`-timeout 20m`),
`-race` legs, 29 named inventory/help/quickstart tests `--- PASS` (0 SKIP/FAIL), XB1+XB2 required
(`WORLD_REQUIRE_PASSKEY_BROWSER=1`, headless Chrome 157, temp profile, CDP virtual authenticator) three
consecutive runs: 9 s, 13 s, 15 s.

## VERIFY_GO_LEGS=plain verify_go.sh at e07b7bad
rc 1, run concurrently with the judge's mutant runs (load 10–15): three package-level `test timed out
after 10m0s` panics — host/broker 601 s, host/daemon 604 s, host/verifygate 600 s. Re-run alone with
`-timeout 20m`: host/daemon rc 0 (381 s), host/broker rc 0 (256 s), host/verifygate rc 1 on
`TestMissionConfigDefaultFlowIntegration/valid_ledger` (plan §6 known load-sensitive list), which passed
alone (`-run ^TestMissionConfigDefaultFlowIntegration$ -v`: valid_ledger and invalid_ledger PASS).
host/broker and host/verifygate import nothing this sprint touches; the planner measured verify_go red at
the base `27ce163e` in six untouched packages at load 40–66 (plan §1.1).

## Judge reproductions (controller, in the judge worktree after it finished)
- W-wn2a (`if false && !approval.ApproverNamePattern…`): `TestPasskeyEnrollFlagRules` rc 1 — KILLED.
- R-csp (`style-src 'unsafe-inline' blob:` in `surfaceCSP`): `TestPasskeyPageRoutes` rc 0 and host/daemon
  `-run 'CSP|Csp|SurfaceApp|SurfaceScript'` rc 0 — SURVIVED; the CSP literal is not pinned by value
  (pre-existing; this sprint does not change it). Logged, not filed (D-WORLD-109 (2)).
