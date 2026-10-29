# Executor report: w-exec-probe-reset-holder-false-pass (iteration 318, row 232)

Worktree: `/Users/voightkampff/dev/sunholo-data/.wt-world-iter318-row232`, branch `sprint/world-iter318-row232`, base HEAD `12ddcf2a`. No git write operations were run. Snapshots: `.snap/M1/` (2 Go files) and `.snap/M2/` (all six scope files, cumulative).

Base sha256 (all six) matched the sprint JSON `base.sha256` values before the first edit.

## Milestone M1: classifier, dial seam, T8, T9, mutation drill

Files changed (`git diff --stat`):
- `host/broker/exec_probe.go`: 25 lines changed (import `syscall`; `dialErrIsSilent`; `pollLoopback` dial seam, rewritten doc comment and `else if dialErrIsSilent(err) { silent.Store(true) }` branch; call site `:566` now passes `nil`)
- `host/broker/exec_probe_test.go`: +93 (import `syscall`; `revDialErr`; `TestExecProbeReverseDialErrIsSilent`; `TestExecProbeReversePollRecordsConnectReset`)

Acceptance:
| Command | rc | Result |
|---|---|---|
| G1 `gofmt -l host/broker/` | 0 | empty |
| G2 `go vet ./host/broker/` | 0 | clean |
| G3 `go test -run '^$' ./host/broker/` | 0 | `ok ... [no tests to run]` (compile fence) |
| T8+T9 `-run '^TestExecProbeReverse(DialErrIsSilent\|PollRecordsConnectReset)$' -count=1 -v` | 0 | 2 top-level PASS, 14 subtest PASS, 0 FAIL, 0 SKIP |
| Focused `-run '^TestExecProbeReverse' -count=1 -v` | 0 | 8 top-level PASS, 0 FAIL |
| G4 `go build ./... && go vet ./host/broker/ && go test ./host/broker/ -count=1` | 1, then 0 | first run: FAIL only in `TestExecProtectedCheckTimeout/after-cancel` (see deviations); rerun: `ok github.com/sunholo-data/ailang-world/host/broker` 276 s |
| G5 `go test -race ./host/broker/ -run '^TestExecProbeReverse' -count=1` | 0 | 0 `DATA RACE`, `ok` 6.9 s |
| G7 `./scripts/verify_ail.sh` | 0 | `verify gate PASSED` (no `.ail` touched) |
| G9 `git diff --name-only` | - | exactly the two Go files |

Mutation drill (each mutant applied alone, only its named killer rows run, then restored from `/tmp/w318-exec/exec_probe.go.snap` and checked with `shasum -a 256 -c /tmp/w318-exec/m1.sha`; both files `OK` after every restore):

| Mutant | Killer rows run | Red/green observed | Restore sha256 |
|---|---|---|---|
| M22 MUT-DIALERR-DROPPED | T9 `connect-reset`, `refused-only` | `connect-reset` FAIL; control `refused-only` PASS | OK |
| M23 MUT-RESET-NOT-SILENT | T8 `op-reset`, `bare-reset`, `wrapped-reset`; T9 `connect-reset` | all 4 FAIL | OK |
| M24 MUT-RESET-NO-UNWRAP | T8 `op-reset`, `wrapped-reset`, `bare-reset` | `op-reset`, `wrapped-reset` FAIL; control `bare-reset` PASS | OK |
| M25 MUT-REFUSED-IS-SILENT | T8 `op-refused`, `real-refused`; T9 `refused-only`; T5 `clean`; all 8 T2 rows | all FAIL (8/8 T2 rows) | OK |
| M26 MUT-ANY-ERR-IS-SILENT | T8 `op-refused`, `op-einval`, `op-aborted`, `timeout`, `addrnotavail`, `nil`; T5 `clean` | 5 FAIL, control `nil` PASS; `clean` FAIL | OK |
| M27 MUT-EINVAL-IS-SILENT | T8 `op-einval`; T9 `einval-then-refused` | both FAIL | OK |
| M28 MUT-TIMEOUT-IS-SILENT | T8 `timeout`, `deadline-exceeded` | both FAIL | OK |
| M28b MUT-DEADLINE-IS-SILENT | T8 `deadline-exceeded`, `timeout` | `deadline-exceeded` FAIL; control `timeout` PASS | OK |
| M29 MUT-CONNECT-RESET-IS-FOREIGN | T9 `connect-reset` | FAIL | OK |

All nine mutants compiled (no `build failed` or `imported and not used`). Raw logs: `/tmp/w318-exec/drill/*.txt`, summary `/tmp/w318-exec/drill-run.txt`.

## Milestone M2: CI, QUICKSTART, doc amendments, AC-STABLE, gates

Files changed (cumulative from base):
- `.github/workflows/ci.yml`: 1 line (step 8 `tests=` list at :247, + T8 and T9 names)
- `docs/QUICKSTART.md`: 1 line replaced by 2 (§10 arm 6, :1153)
- `design_docs/planned/w-exec-probe-reset-holder-false-pass.md`: four planned amendments (§2 timeout row; §6 T9 bullet; §6 AC-LINUX restated; §9 V16 row after V15), plus the D1 edit (§5/§6 M28 row killer and new M28b row)
- `design_docs/implemented/w-exec-probe-reverse-leg-port-race.md`: 2 notes (header `Amended by` line; §4.2 indented note)

Acceptance:
| Command | rc | Result |
|---|---|---|
| G1 gofmt | 0 | empty |
| G2 vet | 0 | clean |
| G3 compile fence (`-count=1`) | 0 | `ok ... [no tests to run]` |
| G8 `actionlint .github/workflows/ci.yml` | 0 | no output |
| YAML load (ruby) | 0 | `yaml ok` |
| `grep -c 'TestExecProbeReverseDialErrIsSilent TestExecProbeReversePollRecordsConnectReset"' ci.yml` | 0 | 1 |
| `-run '^TestQuickstartSection10ProfilesLoad$' -count=1 -v` | 0 | `--- PASS` |
| Step-8 emulation (3 names, `-p 1 -timeout 600s -v`) | 0 | PASS for all 3 top-level names; `reset-holder` subtest PASS line count 1 |
| G6 AC-STABLE `-run '^TestExecProbeReverse' -count=20 -v` | 0 | each of the 8 names `--- PASS` 20 times; 0 `--- FAIL`; `reset-holder` subtest PASS 20; 146 s |
| G4 full package (second run) | 0 | `ok github.com/sunholo-data/ailang-world/host/broker` 179 s |
| G5 race | 0 | 0 `DATA RACE`, `ok` 6.2 s |
| G7 verify_ail | 0 | `verify gate PASSED` |
| G9 `git diff --name-only` | - | exactly the six scope files |

Full `-count=20` log: `/tmp/w318-exec-c20.txt`.

## Deviations and observations

1. **G4 first run, M1 (unrelated test).** `TestExecProtectedCheckTimeout/after-cancel` failed (`exec_protected_record_test.go:719`, "the call was not cancelled during the after-phase"). Re-run alone 3 times: PASS 3/3 (`/tmp/w318-exec-protected-{1,2,3}.txt`). A full G4 re-run then passed. Reading: load-dependent flake, unrelated to `exec_probe`. Not touched.
2. **G4 at M2, unrelated test.** `TestWorktreeSnapshotRefusesAConfiguredFilter` failed (`worktree_delta_test.go:432`, the clean-filter control), the known flake the plan describes. Re-run alone 3 times: PASS 3/3. A full G4 re-run then passed (rc 0). Nothing in `worktree_delta*` was touched.
3. **Mutants applied by script, not the Edit tool.** `/tmp/w318-exec/scripts/mut.py` replaces the exact anchor text (asserted to occur exactly once) with each mutant's text. Restores were `cp` from the snapshot, checked by sha256. Same edits as the plan.
4. **Baseline focused run skipped.** I checked the six base sha256 values, but I did not run the plan's §3 focused base run (6 PASS) before editing. The post-M1 focused run is 8 PASS.
5. **V16 placement.** My first insert put V16 before V15. I moved it directly after the V15 row to match the plan.
6. **Design doc killer matrix not updated.** The design's §6 summary table still lists `M23…M28` for T8 and does not list M28b. The executor instructions allowed only the M28 row and the new M28b row, so I left the summary alone. The controller may want it updated.
7. **G3 caching.** The first G3 run at M2 was served from cache (`(cached)`). I re-ran with `-count=1` (rc 0).
8. **T7 under load (optional).** Not run (the plan leaves it to the executor or the controller). AC-LINUX-LOAD is the controller's.
9. **Base HEAD.** The sprint JSON lists base head `e1110207`. The worktree HEAD is `12ddcf2a` (the plan commit). The six base sha256 values match the JSON.
10. **New untracked paths.** `.snap/` and this report are new untracked paths in the worktree. No git writes were made.
