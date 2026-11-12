# Iteration 386 controller gates (before the judge, not concurrent with it)

Worktree `.wt-world-iter386-a3a-s3`, head `2936964e`, base `82fb9afa`.

## Controller round trip
`git archive 82fb9afa` → awk-extract the 9 Appendix P patches from the committed plan → `patch -p1` in order → `diff -rq` vs the planner's `plan-scratch/work` for `cmd`, `host`, `docs`: empty. Appended-byte control detected. After all 9 commits, the sprint tree equals the round-trip tree (`diff -rq` empty).

## Numstat after each commit (`git diff --numstat 82fb9afa HEAD -- cmd host`, ins+del)

| Sub-milestone | Commit | Numstat | Plan |
|---|---|---|---|
| M6a | `854c0a82` | 29 | 29 |
| M5d | `3a023498` | 280 | 280 |
| M6b | `fea1889a` | 788 | 788 |
| M6c | `693f01c7` | 943 | 943 |
| M7a | `34aaa64b` | 976 | 976 |
| M8a | `9424444d` | 1,383 | 1,383 |
| M8b | `4d161376` | 1,502 | 1,502 |
| M8c | `ea34fd73` | 1,788 | 1,788 |
| M8d | `2936964e` | 1,843 | 1,843 |

Measured: 29, 280, 788, 943, 976, 1,383, 1,502, 1,788, 1,843 — equal to plan §7 at every step. `-- docs` 18.

## Pre-commit gate script (plan §6)
`bash gates.sh` rc 0: `PRE_COMMIT_GATES_GREEN 470s`.
    ok  	github.com/sunholo-data/ailang-world/host/surfacepage	0.757s
    ok  	github.com/sunholo-data/ailang-world/host/store	12.758s
    ok  	github.com/sunholo-data/ailang-world/host/authority	1.065s
    ok  	github.com/sunholo-data/ailang-world/host/approval	1.264s
    ok  	github.com/sunholo-data/ailang-world/host/approval/approvaltest	1.283s
    ok  	github.com/sunholo-data/ailang-world/cmd/ailang-worldd	120.023s
    ok  	github.com/sunholo-data/ailang-world/host/daemon	278.830s

Chrome leg (`gates/chrome.log`): 8 top-level `--- PASS`, 0 FAIL/SKIP.
Inventories (`gates/inv.log`): 57 top-level `--- PASS`.

## After commit
`VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh` rc 0: `✓ all 64 required top-level evidence tests passed exactly once`.

## Judge survivor reproduced by the controller
A02 (`host/surfacepage/app/app.js`: `${esc(q.agent)}` → `${q.agent}` in the Needs-you row): `WORLD_EXEC_NODE=$(command -v node) go test -run '^TestSurfaceAppPure$' ./host/surfacepage/` rc 0 with the mutation (`ctl-A02-survivor.log`). Unmutated, the test runs (not skipped) with and without `WORLD_EXEC_NODE`. Test gap only — the code escapes. Carried to sprint 4.
