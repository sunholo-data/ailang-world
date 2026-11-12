# Iteration 385 controller gates (row 201 A3a sprint 2)

All logs under `~/.ailang/state/world-iter385/` on the rig.

## Round trip
`git archive a8c66655` + the 8 patches extracted from the plan file with §2's awk (M4a 359, M4b 205, M4c 345,
M4d 515, M4e 142, M5a 391, M5b 174, M5c 86 patch lines) → `diff -rq` against the planner's `plan-scratch/work`
tree: empty for `cmd` and `host` (control: an appended byte made `diff` report the file).
After the last commit, the sprint worktree's `cmd` and `host` equal the round-trip tree.

## Commits (numstat `git diff --numstat a8c66655 HEAD -- cmd host`, ins+del, measured after each commit)
| Sub-milestone | Commit | Numstat | Plan §7 |
|---|---|---|---|
| plan | `0d88d820` | — | — |
| M4a | `fb4e2ae2` | 350 | 350 |
| M4b | `b9025b20` | 517 | 517 |
| M4c | `627d825e` | 793 | 793 |
| M4d | `f76c6bea` | 1,257 | 1,257 |
| M4e | `6a3f07ae` | 1,396 | 1,396 |
| M5a | `76aa7dcd` | 1,741 | 1,741 |
| M5b | `456fb26c` | 1,912 | 1,912 |
| M5c | `d52b3ae3` | 1,995 | 1,995 |

Every commit went through `commit_ms.sh` (gofmt, build, vet, the sub-milestone's named tests). The executor's M4e and M5c
snapshots equal the committed tree (`cmp` per file, empty).

## Plan §6 pre-commit gates (`gates.sh`, one bash shell, not under nohup)
rc 0, `PRE_COMMIT_GATES_GREEN 355s`: verify_ail, gofmt, build, vet, compile fence, the six touched packages, the
three `-race` legs, every inventory and all 20 new tests by name (`--- PASS`), the AW leg 9/9.

## After commit
`VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh` → rc 0: "0 other binary blobs
among 6778 tracked files", toolchain floor ✓, "all 64 required top-level evidence tests passed exactly once",
"go gate PASSED (plain)".

## Executor deviation (adjudicated by measurement)
Run B started M5a's whole-`host/daemon` leg under `nohup` (plan §7 forbids it). Its log reads
`ok host/daemon 250.244s` and its rc file reads `0`; the controller's §6 gate re-ran the whole package (not under nohup)
green. Procedural only; no result depends on it.

## Judge findings reproduced by the controller
- **J36 (survivor, non-blocking):** `host/daemon/surface_approve.go` `handleSurfaceRequestCard`, the
  `ActivePasskeys` error check replaced by `if false` → `go test ./host/daemon/ -run '^(TestApproval|TestApprovePOST|TestRequestsJSON|TestRequestCardJSON|TestDaemonCredentialWriterSites|TestSurfaceAppRoutes)'`
  rc 0 (survives). Read-only and fail-soft (the card shows no passkeys instead of a 500). Carried to sprint 3.
- **Record before the law (must-die):** the judge's K01 was killed by the writer inventory under `-failfast`.
  The controller re-ran it with the inventory EXCLUDED — a record inserted before `JudgeAssertion` →
  `go test ./host/daemon/ -run '^(TestApprovePOSTRefusalArms|TestApprovePOSTRecordsAVerifiedPasskey)$'` rc 1
  (R11-other-request, R11-unenrolled, R11-counter, R11-garbage fail on their 0-row check); unmutated, the same command rc 0.
- **Judge's numstat 1,596** is a measurement error: the stated unit (ins+del) reads 1,995, insertions alone 1,991.
  Under the 2,000 cap either way.

## CI
- PR #556 head `d52b3ae3`, run 43138995090 attempt 1: `go host build + test gate` red on
  `host/broker` `TestSeedRestartIgnoresASurvivingOldCopy/survivor` (`exec_seed_test.go:656`, "the old staging dir
  survived the restart: <nil>"); the sprint touches 0 `host/broker` files. Row 259's known flake, instance 4.
  `gh run rerun --failed` → attempt 2: 3/3 success.
- Merged `--match-head-commit d52b3ae3` → `d1027477` (02:36:43Z). No closing keywords in title or body.
