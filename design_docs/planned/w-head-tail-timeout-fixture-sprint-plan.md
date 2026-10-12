# Row166 timeout fixture prerequisite — sprint plan

Ready for test-only implementation, attended2026-10-10; no repair claimed yet.
Design: [w-head-tail-timeout-fixture.md](w-head-tail-timeout-fixture.md).
Only product edit: `host/broker/handlers_capture_test.go`; no production/.ail edits.
Row174 plan, benchmark, mission, handler caps and shared harness stay unchanged.
Fresh available-lane review required; D-WORLD-48 fallback flagged if used, no quorum claimed.

## M1 — finite measured fixture (~70 LOC, 1–2 active hours)

Keep TestRunBoundedHeadTailTimeoutKeepsPartialOutput, real runBounded, `/bin/sh`,
process-group SIGKILL, real captures and standard context. No resettable deadline/mock error.
1. One executable prints PARTIAL-OUT/stdout and PARTIAL-ERR/stderr, then writes readiness,
   then optionally sleeps an environment-selected integer-second hold. Never preseed sinks.
2. Run exactly three no-sleep controls through that same executable with fresh captures,
   normal20s handler allowances and one finite boundedTestContext (30s caller deadline).
   Require successful exit, readiness and both exact Text/Bytes/SHA256, untruncated.
   Remove readiness between runs; log each complete control elapsed duration and maximum C.
3. With existing grace G=pipeCloseGrace=1s, set B=ceil_ms(C+2*G) and
   H=ceil_seconds(3*B+2*G). C already includes same-executable startup, both writes,
   natural exit and drain; two existing cleanup-grace units provide finite scheduling
   headroom without the unsupported tenfold multiplier. Immediately before exercised arm
   measure R=time.Until(callerDeadline).
   Require B>0, H>=3*B+2*G and B+H+G<R; check overflow/invalid deadline and fail loudly.
   Thus startup allowance B plus natural hold H plus cleanup grace fits remaining30s.
   No clamp, skip, retry or alternate measurement to rescue an insufficient budget.
4. For exercised arm wrap killGroup, inspect readiness BEFORE calling captured original
   syscall, record target/result/readiness with synchronization and always delegate real kill.
   No t.Parallel. Keep seam installed until child/Wait cleanup completes; then restore it.
5. Require real kill, readiness-before-kill, errors.Is(ErrHandlerTimeout),
   errors.As(*HandlerTimeoutError) with Timeout==B, and exact Text/Bytes/SHA/untruncated
   for BOTH streams. Log C/B/H/R and kill observation. Calibration is not causal proof.

Acceptance: three controls plus one actual deadline arm, finite guard and exact stream/error
assertions pass; production files byte-equal pre-repair snapshots. Existing baseline red log
is preserved. Why host/not package: Go subprocess-capture regression fixture only.
Registry reuse:none; existing runBounded/capture/killGroup mechanisms reused, no package capability.

## M2 — mutation, repetitions, normal gate and fresh review (0 product LOC)

Depends M1. Estimate2–4h elapsed, gated on full normal Go profile and review.
Create isolated mutant copies and keep exact original product/test SHA snapshots.

| Mutant | Exact required red control |
|---|---|
| stdout capture discarded | stdout Text/Bytes/SHA exact assertion |
| stderr capture discarded | stderr Text/Bytes/SHA exact assertion |
| either child printf removed (two arms) | that stream's exact assertion, never zero-output pass |
| readiness creation removed | readiness assertion; disclose if positive control stops before exercised arm |
| MUT-TIMEOUT-DISABLED-NATURAL below | successful natural exit/no kill contradicts required exact timeout-B/real-kill assertion |
| MUT-DEADLINE-STDOUT-LOSS: lose stdout only in exercised deadline arm | Three controls PASS; actual kill/readiness and exact timeout B reached; deadline stdout Text/Bytes/SHA assertion RED |
| MUT-DEADLINE-STDERR-LOSS: lose stderr only in exercised deadline arm | Three controls PASS; actual kill/readiness and exact timeout B reached; deadline stderr Text/Bytes/SHA assertion RED |

Timeout-only disabling leaves runHeadTail cleanupDeadline=B+2G active: it may return
ErrCleanupIncomplete before H elapses. Do not claim that is a natural successful return.
For the explicitly compound isolated MUT-TIMEOUT-DISABLED-NATURAL, replace inner WithTimeout
with cancellable parent context AND move only the mutant cleanup deadline to caller deadline
minus G. Real exec/Wait/sleep, parent deadline and capture paths remain active. Require child
natural exit and reap before caller guard, no real kill; exact timeout assertion must fail.
Never land either production mutation. On failure paths await actual child/reap within the
finite caller guard before restoring killGroup; no active waiter may race seam restoration.
If this cannot be discharged without changing production beyond those disclosed mutant lines,
stop with evidence rather than inventing a workaround or claiming a kill.

Eight distinct mutation arms total: original six plus two deadline-only stream-loss arms.
Original discard/printf mutations may fail calibration early; report that scope honestly, never
credit them as deadline-retention proof. The added two preserve calibration and must reach
the exercised real deadline/kill/readiness before their exact partial-stream assertion fails.
Use isolated mutant edits only; no new landed scope or harness change.
After EACH mutant prove expected red identity, exact-restored SHA values and pristine green.
Do not use git restore on unsaved deliverables. Bank all attempts, no retry-until-green.
Bounded pristine repetitions: focused plain count10 and focused race count5, verbose,
normal defaults (no special GOMAXPROCS/harness/cap changes); all15 must pass. Each invocation
has existing per-test finite context; use a600s outer Go-test guard and record actual count.
Full required profile unchanged: released AILANG_BIN/WORLD_PKG_AILANG_BIN v0.41.0,
WORLD_TOOL_AILANG_BIN v0.52.1; verify_go.sh covers build/plain/race. Root credits existing
verify_ail only with .ail byte-equivalence. Do not suppress any failure or relevant skip.
Fresh reviewer checks finite math, actual kill/readiness, both stream hashes, exact error B,
mutation honesty, natural reaping and seam restoration, and full normal gate evidence.

Evidence/progress are separate under
`design_docs/verification/world-row166-fixture-attended-2026-10-10/` and
`.ailang/state/sprints/sprint_w-head-tail-timeout-fixture.json`.
M1 implementation proof plus M2 PASS allows row174 normal gate to resume; it does not close
row174, publish/mint sessions, or permit FINAL. Root owns executor handoff and commits.

Bounded revision after blocking round2: only the measured-B formula/rationale changes.
Pre-revision records/mutations/counts are preserved in evidence `factor10-round2/` with
SHA index. Reported cold C=2.353960541s yields B=4.354s,H=16s,total21.354s<R27.636s;
reported full-plain C=0.756807791s yields B=2.757s,H=11s,total14.757s<R29.234s.
These are arithmetic projections, not passing execution proof. Same eight mutants,
plain count10, race count5, full normal Go gate and FRESH round3 review are mandatory
for the revised fixture. Previous greens and blocking reds remain historical evidence.
