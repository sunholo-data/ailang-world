# Iteration 376 — controller gates (row 201 A1b sprint 2)

Sprint head `3ea0bea9` (branch `sprint/world-iter376-row201-a1b-s2`, base `115f5a2b`). AILANG_BIN = `~/.pinned-ailang/ailang` (AILANG v0.41.0). Script: `gates.sh`.

## Sweep 2 (the verdict), all rc 0

| Gate | Result |
|---|---|
| `./scripts/verify_ail.sh` | rc 0 (40 s) |
| `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` | rc 0 (415 s). 0 other binary blobs among 6,618 tracked files; 64 evidence tests passed exactly once (`verify_go-3ea0bea9-tail.txt`). |
| `go build ./... && go vet ./...` | rc 0 |
| `TestProductionContextRoots` | rc 0 |
| `TestDaemonDependencyAllowlist` | rc 0 |
| `go test -race` on store, approval, authority, and the approver tests in `cmd/ailang-worldd` | rc 0 |
| `go test ./cmd/ailang-worldd/ -count=1` | rc 0 (144 s) |

Also run locally: `scripts/check_no_personal_email.sh` rc 0; `scripts/queue_census.sh` with the CI controls rc 0.

## Sweep 1 (not a verdict)

Sweep 1 ran while the judge's own `go test` of `cmd/ailang-worldd` (291 s) was running on the same rig (load average about 12 on 16 cores). `verify_go` returned rc 1 with five failures in packages the sprint does not touch: `host/broker` (`TestExecProbeNet6ForeignDialIsRetriedNotRefused`, `TestPublishClaimIsDurableBeforeDispatchAndSurvivesAFreshBudget`), `host/daemon` (`TestWorkspaceExecEpisodeSelectsItsProfile`), `host/surfacepage` (`TestSurfaceAppFollow`, 30.04 s) and `host/workbench` (`TestLiveScriptPure`, 30.01 s).

- All five passed when re-run alone at the same head: broker 7.2 s, daemon 1.6 s, surfacepage 0.5 s, workbench 0.5 s.
- All five passed in sweep 2.
- `TestExecProbeNet6ForeignDialIsRetriedNotRefused` is the **second** recorded instance (the first was the iteration 375 judge's run).
- `TestPublishClaimIsDurableBeforeDispatchAndSurvivesAFreshBudget` failed on a counter (`POST=2`, want 1), not a timeout. This is its first recorded instance. It is logged per D-WORLD-109 and not filed.

## Per-milestone reconstruction

Each sub-milestone was restored from `~/.ailang/state/world-iter376/snap/<M>/` and run through gofmt, build, vet and its acceptance tests (`commit_ms.sh`), then committed. After each commit the tree was compared against the executor's snapshot with `cmp`: 0 mismatches.

| Sub-milestone | Commit |
|---|---|
| plan | `c5f00048` |
| M3a | `df9d6280` |
| M3b | `4efba493` |
| M3c | `8437ff2b` |
| M3d | `a80592fe` |
| M3e | `1aee825f` |
| M3f | `636e7761` |
| M3g | `cad109f9` |
| size-gate note | `3ea0bea9` |

Size gate (plan §7) after M3g: `git diff --numstat 115f5a2b HEAD -- cmd host` = **1,896**, over the 1,683 threshold. So Run 5 (M4a/M4b) did not run.

## Judge finding reproduced by the controller

This was run in the judge's worktree; the file was restored afterwards and `git diff --quiet` passed.

- **F-CREDS survives.** In `approverFlags`, the `--credentials-db` value was replaced with an always-empty string, via `strings.TrimPrefix(*credsFlag, *credsFlag)` so the code still compiles. The approver tests plus the help tests still pass (`ok`).
  - A first attempt that returned a literal `""` failed to build (an unused variable), so it was not a mutant. It was redone.
  - Control: changing `approverConfirmBudget` from 120 s to 121 s makes `TestApproverEnrollUnansweredPromptTimesOut` FAIL.
