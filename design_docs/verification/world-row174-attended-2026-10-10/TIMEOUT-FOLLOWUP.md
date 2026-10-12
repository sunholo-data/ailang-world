# Pre-existing partial-output timeout failure

First full candidate normal Go profile failed only
`TestRunBoundedHeadTailTimeoutKeepsPartialOutput` in `host/broker`: at the
300 ms deadline, partial stdout was empty rather than `PARTIAL-OUT`.
See `verify-go-first-red.log`. The race leg was not reached in that attempt.

The exact unchanged test on untouched parent `5ca1115`, `-count=20`, produced
13 PASS and 7 FAIL with the same empty-stdout assertion. See
`baseline-partial-output-control.log`. Neither the broker implementation nor
that test differs between parent and candidate. The candidate changes eight
package description strings and adds metadata verification tests only.

This establishes a pre-existing timing-sensitive failure; it does not establish
whether startup scheduling, stream capture or another cause is responsible.
Host load observed later was `{5.38 6.34 6.95}`; this single observation does
not establish the failure's cause. No timeout, assertion or lifecycle was changed.

The unchanged full normal profile rerun passed the plain leg, then failed the
same test in the race leg. All other race packages passed. See `verify-go.log`.
No waiver is applied; M1 round 1 remains FAIL (68/100).

The mission already records this exact fixture as row 166,
`w-broker-head-tail-timeout-flake`. Its separate test-only repair is designed in
`../../planned/w-head-tail-timeout-fixture.md` and follows the normal role review
path because it now blocks row 174's normal gate. Row 126 names a different test.
Production bounds, benchmark settings and shared harness remain unchanged.
All failed attempts are retained. Calibration and readiness checks do not prove
one unique historical failure cause. Repair verification and fresh review are
pending; a later pass must not rewrite these observations.

## First repair rejected in round 2

The first measured fixture at product commit `38c61e4` used `B=10*C+2grace`.
Its focused plain10/race5 controls and eight isolated mutants passed. Fresh judge
pristine control then failed the finite-fit guard: C=2.353960541s, B=25.54s,
H=79s, R=27.636286083s. The full normal Go profile also failed the new guard
(C=756.807791ms, B=9.569s, H=31s, R=29.234203208s); rc=1, race NOT_REACHED.
This new fixture failure is not attributed to the untouched parent or waived.
A later standard-cache control passed; that does not erase the pristine red.
Round2 remains FAIL68. See separate row166 `normal-go-result-round2.json`,
`judge-round2-cold-guard-red.log` and preserved `factor10-round2/` evidence.

A bounded approved design revision uses `B=ceil_ms(C+2grace)`: C measures the
complete same-child control, already including startup and output drain.
The original tenfold multiplier had no verified capture-contract rationale.
The revised exercised deadline is tighter; both-output, real-kill, child-readiness,
cleanup, finite H and ordinary30s caller assertions remain. Arithmetic fits the two
reported C/R observations, but passing proof is pending new eight mutants,
plain10/race5, exact full Go profile and fresh final round3 review.

## Final repair review / separate gate blocker

At `84d3f5a`, revised C+2grace fixture current plain10/race5 and eight mutations
passed, and fresh judge3 independently verified a deadline-only loss mutant.
The final normal Go profile passed the repaired fixture but failed the unchanged
`TestExecHandlerTimeoutKillsTheGroup`: 800 ms timeout reported around802ms, empty
stdout/stderr instead of the expected started marker. rc1; race NOT_REACHED.
All five recorded product hashes remained exact. See separate row166
`normal-go-final-result.json` and `verify-go-row166-final.log`. Parent isolated
count20 passed all20 (rc0); this does not establish a unique cause or waive the red.
Round3 remains FAIL68 and human review is required by shared Gate3's max3 policy.
The test-only full-trampoline/readiness proposal is prepared, unimplemented,
in `../../planned/w-exec-timeout-stimulus-fixture.md`; approval is pending.
