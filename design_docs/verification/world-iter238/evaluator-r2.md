VERDICT: PASS 96

# Row 141 evaluator r2 (iteration 238), judging `2319d9d` against r1 `191d4e5`

- Evaluator worktree at `2319d9d`: `git diff --quiet` is clean and `git status --porcelain` is empty.
- Mutants were run on a `git archive 2319d9d` copy (`~/.ailang/state/world-iter238-row141/eval/mut`). It was restored, and `diff -r` against a fresh archive is empty.
- Raw outputs are in `~/.ailang/state/world-iter238-row141/eval/r2/`.

## What changed (re-derived)
`git diff --stat 191d4e5 2319d9d` lists 4 files: `executor-report.md` (+42), `docs/QUICKSTART.md` (8), `host/daemon/workspace_layout.go` (+23 lines, net +18) and `host/daemon/workspace_layout_test.go` (+96).

- `f680e28` against `191d4e5` changes `executor-report.md` only (+6).
- The coordinator's list is complete; nothing was missed.
- `host/broker` is byte-unchanged since r1.
- Over the whole row (3c76e99..2319d9d), the only exec-adjacent edits are:
  - `host/broker/execargs.go`: a one-line exported `ExecPathOk` wrapper with no behaviour change.
  - `cmd/ailang-worldd/quickstart_exec_test.go`: the one-line §10 cut.
- No MCP or srt file is touched. This supports the coordinator's reading that the `TestExecSrtMCPEndToEnd` CI red is the base-inherited row-153 class.

## Findings and survivors: status (my own mutants re-applied at 2319d9d)
| Item | Status | Evidence |
|---|---|---|
| N1 lock TOCTOU | CLOSED | Single `OpenFile(O_RDONLY\|O_NONBLOCK\|O_NOFOLLOW)` then `f.Stat()`. The window is gone. Dropping O_NONBLOCK → KILLED (FIFO test: "hung … no refusal within 2s"). Disabling the Fstat type check → KILLED (the FIFO now reads as EOF and is refused as "not valid JSON", which fails the message assertion). Dropping O_NOFOLLOW → KILLED (CoversLock/lock_is_a_symlink). |
| N2 / E8 1 MiB bound | CLOSED | KILLED by CoversLock/lock_over_1_MiB. |
| N3 / E10 group/other-writable | CLOSED | KILLED by 3 new StartupChecks rows. |
| N4 recipe `set -e` leak | CLOSED (by inspection) | The recipe is now `( set -euo pipefail … chmod -R a-w "$SNAP" )`. The guard's `exit 1` ends the pipeline subshell, pipefail and `-e` end the outer subshell, and the operator's shell keeps its options. My verbatim execution attempt was denied by the tool permission layer, so I did not re-run it. The executor reports running it verbatim. |
| N5 / E25 unset arm creates registry | CLOSED | KILLED by UnsetLeavesHomeAlone/missing_registry_stays_missing. |
| N6 moved-snapshot message | CLOSED | New symlink case with its own line, and the §11 line table is updated. A mutant that removes the link before refusing → KILLED (the link survives). Routing the case to the old default → KILLED (exact line). |
| N7 / E9 toml is a dir | CLOSED | KILLED by CoversLock/"ailang.toml of a locked package is a directory". |
| N7 / E16 deep toml counted | CLOSED | KILLED by StartupChecks/"ailang.toml one level too deep". |
| N8 SKIP count | CLOSED (informational) | Daemon SKIP is 24 here too; it depends on the environment. |
| Regressions spot-check | — | E1 (RemoveAll any dir) and E7 (bad JSON fail-open) are still KILLED. |

**Adjudication, E9/E16:** the executor and r1 are both right. My r1 mutants changed correct production code (`IsRegular` → `st != nil`; a 4-segment check widened to any depth), and the r1 suite failed to catch either change: that was a test gap, not a code defect. Adding tests and leaving the code alone was the correct fix, and both mutants are now killed.

## New production hunk: scan
- `O_NOFOLLOW` on a final-component symlink returns ELOOP on both Linux and darwin. The `errors.Is(err, syscall.ELOOP)` branch maps it to "is not a regular file". Intermediate components are not covered by O_NOFOLLOW, but the sandbox is already resolved symlink-free.
- `O_NONBLOCK` with `O_RDONLY` on a FIFO returns at once on both systems. On a regular file it has no effect. Go's poller registration falls back for regular files: epoll gives EPERM on Linux, and on darwin regular files are kindOpenFile and not pollable.
- Both flags are POSIX constants in `syscall` on linux and darwin, so they compile on both. CI is ubuntu.
- The FIFO test cannot hang CI. Its 2 s `select` timeout releases a blocked open by opening the FIFO with `O_RDWR|O_NONBLOCK` before `t.Fatal`. With the real code it returns in about 0.2 s. The buffer is read only after the channel receive (happens-before), so there is no data race; the `-race` daemon run is green.
- New symlink case ordering: a symlink pointing at the snapshot is still kept, and an empty dir is still replaced. A dangling or foreign symlink is refused and never deleted. No defect found.
- Minor: on a FIFO the open-error branch "cannot be read" is unreachable, and a Readlink error leaves `old` empty in the message. Cosmetic.

## Gates at 2319d9d (packages run one at a time, `-race -count=1 -v`)
- `go vet ./...`: rc 0.
- daemon: rc 0. 481 RUN, 457 PASS, 24 SKIP, 0 FAIL. That is 10 more RUN than r1 (the new rows plus the FIFO test).
- broker: rc 1. 627 RUN, 603 PASS, 23 SKIP, 1 FAIL. The single failure is `TestRunBoundedHeadTailTimeoutKeepsPartialOutput` (`partial stdout = ""`), the known base flake. Re-run alone it went **3/3 ok** (1.75 s each). The broker code is byte-unchanged since r1, where the same package was green, so this is not this diff.
- cmd: rc 0. 115 RUN, 112 PASS, 3 SKIP, 0 FAIL.
- RIG pair: rc 0. 2 RUN, 2 PASS, 0 SKIP.
- `verify_ail.sh` was not re-run: there are no `.ail` changes since r1, where it passed.

## Score (100)
| Area | Score |
|---|---|
| AC coverage | 30/30 |
| Tests and mutation | 24/25 (every r1 survivor killed) |
| Code | 14/15 |
| Security | 14/15 (the residual is risk R2, accepted by design) |
| Docs | 9/10 (N4 checked by inspection only) |
| Gates | 5/5 (broker red is the base flake, green alone 3/3) |
| **Total** | **96** |

There are no blocking findings.
