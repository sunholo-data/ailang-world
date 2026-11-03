# Evaluator report, row 241 (iter 339), round 1

Subject: e2d418d0 (code base dcb2366e). Instrument: `go test ./host/broker -run '^(TestExecProbeReversePollCadence|TestExecProbeReverseHelloWaitBounds)$' -count=1` with AILANG_BIN pinned; script /tmp/j339.sh.

## Score: 88/100 -- PASS
Correctness 31/35, test quality 17/20, scope/process 14/15, evidence 13/15, maintainability 13/15.

## BLOCKING
None.

## NON-BLOCKING
1. T12's count bound (>= 6 dials of ~21) is loose by design: the poll interval mutants 60/150/200 ms survive (m5, m4, m2). Only intervals about >=4x too slow, or <50 ms, are killed. Row survivor 1 is the 500 ms literal, which is killed; the looseness is documented in the test comment.
2. T10 survivor 3 is only partly closed. `setAt` is read after the classifier's `time.Now()`, so a hello wait of 950+delta ms (delta = classifier-to-SetDeadline gap, microseconds) survives. Ceiling is exact in practice (960 ms and 1000 ms killed, m6/m7), but a 951 ms mutant is a just-over case that can pass if delta exceeds 1 ms under a stall. A correct tree still cannot red: deadline - setAt <= 500 ms < 950 ms.
3. The floor check in T12 has a 1 ms slack mutant that survives (m10: `gap < poll-1ms`). Equivalent in practice (timers never fire early); no consequence.
4. Evidence lives mostly in designer/controller scratch files; the planner's own sandbox run was uninformative (documented in the plan).

## Mutants run (my own)
- control: ok
- m1 pollLoopback `time.After(500ms)` (the row's mutant): KILLED (dials=3, fastest 500.5 ms)
- m3 25 ms: KILLED (39 dials, fastest 25.4 ms, floor clause)
- m2 200 ms, m4 150 ms, m5 60 ms: SURVIVE (below the 4x count bound; accepted by design)
- m12 execProbePollEvery const 50 -> 100: survives; equivalent because T10's ceiling derives from the constant and T12 reads it
- m6 helloWait 960 ms, m7 1000 ms: KILLED (set+960 > set+950)
- m8 helloWait 200 ms: KILLED (the 250 ms floor)
- m9 T10 ceiling +1 s (test-line mutant): SURVIVES (expected, no 1 s mutant exists... mutating the test's own bound is not a tree defect)
- m10 T12 floor -1 ms: survives (non-blocking 3)
- m11 count bound -> 1: by hand with sed, test still passes (control tree dials ~21), so the count clause is only discriminating via m1-type mutants through the 6 bound; not killed by lowering it alone, as expected

All mutations restored: `git diff --quiet` rc 0 and git status clean.

## Flakiness
`-count=100` of both tests, no synthetic load: ok (103.8 s). A correct-tree red needs the poll goroutine starved while the reference loop is not (about 4x), or a dial gap under 50 ms (impossible for monotonic timers).

## Scope
`exec_probe.go` byte-identical to dcb2366e (git diff --quiet rc 0). Only exec_probe_test.go and ci.yml (the name `TestExecProbeReversePollCadence` appended in the same line group as the other TestExecProbeReverse* names, in the Row 140 Workspace.Exec real-srt verbose step, which is where that list lives) changed, plus design/evidence docs.

## Evidence spot-checks
- 336 PASS / 30 SKIP / 0 FAIL: controller-A2-package-census.txt line 369 (`census: 336 PASS, 30 SKIP, 0 FAIL`) matches doc V15/plan line 76.
- Heavy load minimum dials 19 (gate 6), fastest 50.02 ms: controller-A4-heavy.txt (CADENCE line); load average 170 at the end.
- Mutation round EXPECTATION MATCH and restored(cmp ok): controller-A3-mutation.txt tail; gates summary all rc=0.
- Subject sha in controller-gates-summary.txt is 51622b82; the sprint head e2d418d0 adds only the evidence commit, consistent.

## Not verified
Full-suite and verify_ail (controller-owned, read only); the 61/61 mutation-row count (read EXPECTATION MATCH only); 160-yes load not re-run (forbidden); CI on the real linux runner.
