# Evaluator report, iter-285 round 2, row 156 M4a + M4b (fix commit d4ff957f over 167ec9af)

Verdict: PASS, 93/100. No blocking findings. The r1 blocker is closed and the r1 flake is gone.

Score: correctness vs spec 38/40, test strength/mutation 22/25, gates green 15/15, scope 10/10, docs/residual honesty 8/10 = 93/100. Unaddressed r1 residuals and the order-dependence in finding 1 account for the deductions.

## Diff 167ec9af..d4ff957f (two Go files only, 31 insertions, 2 deletions)
- `protectedCheckTimeout` const -> var, still 2500 ms.
- `protectedRelaxBound(t)` sets it to 30 s with `t.Cleanup` restoring the old value.
- `TestProtectedChecker` rows relax the bound; new row `filter-refused`.
- `TestProtectedCheckerRootAndBound` relaxes every row except `shared-phase-bound`.

## Gates (pinned AILANG_BIN, WORLD_TOOL_AILANG_BIN unset)
Run while the machine was heavily loaded (single rows took 3-6 s versus 0.1-0.3 s in r1), which made a good stress.
- `go vet ./host/broker/`: ok.
- `TestProtectedGitFilterPreflight` 5/5, `TestProtectedChecker` 9/9 (includes `filter-refused`), `TestProtectedCheckerRootAndBound` 7/7: all PASS.
- `go test -race -count=1 -run '^TestProtected'`: ok.
- Stress: 8 concurrent `-race` test binaries, `-test.count=2` each, on the loaded machine: 8/8 PASS (r1: 1 of 6 failed on `RootAndBound/cwd`).

## Mutants (one at a time, bytes restored from a scratch copy each time)
| Mutant | Edit | Killer | Result |
|---|---|---|---|
| r1 survivor: checkProtected ignores preflight | `_ = protectedGitPreflight(ctx, git, root)` (exec_protected_check.go:~29) | `TestProtectedChecker/filter-refused` | KILLED: "filter not refused: {intact:... changes:[{a_test.go modified}] err:<nil>}" |
| MUT-CHECK-DOUBLE-BUDGET | legs use `context.Background()` | `RootAndBound/shared-phase-bound` | KILLED: "shared-phase-bound accepted" |
| New, relax helper: relax every row including shared-phase-bound | `if name != "shared-phase-bound"` -> `if true` | `RootAndBound/shared-phase-bound` | KILLED (accepted; the 30 s bound lets 4 x 0.9 s through) |
| New, relax helper: cleanup does not restore | `t.Cleanup(...)` -> `_ = old` | `RootAndBound/shared-phase-bound` | KILLED (relaxed value leaks from the earlier `cwd` row, so the bound row sees 30 s) |

All kills are assertion failures, not compile errors, panics or timeouts.

## Findings
Blocking: none.

Non-blocking:
1. The NO-RESTORE kill relies on subtest order (`shared-phase-bound` is second, after `cwd`). Reordering the list so the bound row runs first would let a non-restoring cleanup survive. Cheap hardening: assert `protectedCheckTimeout == 2500*time.Millisecond` at the start of the bound row. Reproduce: swap `cwd` and `shared-phase-bound` in the name list at exec_protected_check_test.go:~211 and apply the NO-RESTORE mutant.
2. The package var is mutated without a lock. Safe today, because no test uses `t.Parallel`, but a future `t.Parallel` under `-race` would trip it. A comment on `protectedRelaxBound` would help.
3. r1 residuals still open and still acceptable: (a) preflight output cap is 96 MiB (row-184 code); (b) rc 1 with stdout reads as "no filters" (dead check in `configuredGitFilters`); (c) `%q` of up to 8 MiB stderr in the error string (exec_protected_check.go:49). Carry (c) into the V24 resize.
