# Evaluator report, iteration 329, row 230 (r1)

VERDICT: PASS, 88/100 (Correctness 31/35, Test quality 21/25, Design conformance 24/25, Scope/evidence 12/15 -> 88)

Gates (eval worktree d6665670): gofmt -l host/broker = 0; go vet ok; go test -count=1 -v ./host/broker/ = 335 PASS, 0 FAIL, 30 SKIP, ok (base 333 + T10 + T11; matches executor M2 G4 and out/fullpkg-control 333 at base). TestExecProtectedCheckTimeout stayed green. T10 -count=500 under 32 `yes` (16 cores): ok. exec_probe.go sha256 96f33172... = executor M2 value; base prefix 390c74e841f26a12 reproduces.

BLOCKING: none.

NON-BLOCKING
1. Survivor: `time.After(execProbePollEvery)` -> `time.After(500*time.Millisecond)` stays green. Nothing pins that the constant is what drives the poll cadence, so T10's ceiling derivation (listen window less one poll) can drift from the real cadence. Not equivalent. Repro: sed the line in pollLoopback, run TestExecProbeReverseHelloWaitBounds.
2. Survivor: message literal `within 500 ms` with the argument removed stays green in tests (the T7 extension asserts "500 ms", equal to the literal at today's value). Equivalent at 500; killed only by a paired wait change. (Dropping the argument alone is also caught by go vet in test builds.)
3. T10 ceiling precision: the check is deadline <= stop+950ms, where stop is taken after the poll returns, so wait=951 is red by only ~0.9 ms margin (stop+950.99ms), and slack under load could let 951..~960 survive. wait=960, 1000, 60s are red now. No effect on a correct tree. Note the 950 ceiling reduces to exactly the design's number.
4. T10 test-arithmetic mutants (ceiling=listenWindow, ceiling-2*poll, 250->100 or 500 floor, arrival=before, drop sets==0, drop wrote==hello, after->before) are green alone. Each is equivalent at wait=500 and only distinguishable in combination with a production K-mutant. `conn.wrote` assertion is the only guard that the hello is written. Acceptable.
5. Poll cadence changes to 100 ms or 1 ms, and listen window 1100, are green: legitimate values inside the bounds (equivalent).

MUTATION TABLE (mine; focused regex on T10/T11/T7/Judge/ForeignListener)
| Mutant | Result |
| wait 5 ms (K1) | RED T10 x3: "seen {silent:true} ... want child/foreign, not silent" |
| wait 249 ms | RED same |
| wait 951 ms | RED "hello deadline at stop+950.99ms ... no later than stop+950ms" |
| wait 960 / 1000 ms | RED same ceiling assertion |
| poll 600 ms | RED ceiling (stop+400ms) |
| listen window 400 ms | RED ceiling (stop+350ms) |
| poll 100 / poll 1 ms / listen 1100 | GREEN (in-bounds, equivalent) |
| time.After literal 500 ms | GREEN (survivor 1) |
| Milliseconds() -> Microseconds() | RED T7 x4 ("within 500000 ms") |
| Milliseconds() -> Seconds | RED T7 x4 ("within 0 ms") |
| /10 on Milliseconds() | RED T7 ("within 50 ms") |
| message literal 500 + arg dropped | GREEN (survivor 2) |
| remove c.resume() (two forms) | RED T11 "handler is not net.createServer(c=>{c.resume();..." |
| test: ceiling=listenWindow / -2*poll | GREEN alone (equiv.) |
| test: floor 100 / 500 | GREEN alone; floor 600 RED |
| test: arrival=before, drop sets==0, drop wrote | GREEN alone |
| test: ceiling literal 950 | compile red (unused var), not run |
Planner/executor K1, K2/K2c, K3, K4, K5 rows re-run by hand: reproduced as stated in executor-M2-mutants.txt.

PER-MILESTONE NON-VACUITY (plan D-1, adjudicated by measurement)
- M2 (comments + T11): reverting M2's production hunk (git apply -R of 5ddc93d3 on exec_probe.go; comment-only, plus it removes none of the code) leaves T10/T11/T7 GREEN. Confirmed: behaviour-neutral; T11's kill is the c.resume() text mutant (RED, above), which is a structural pin the design labels S6.
- M1: reverting M1's hunk removes the constant; T10 then fails to COMPILE (undefined: execProbePollEvery) - red by build, not by assertion. With the constant kept and message literal 500 + time.After 50 ms restored, tests are GREEN. Confirmed D-1: at today's values no production revert reds an assertion; the real kills are the K-mutants (wait/poll/listen/format derivation), which I reproduced.

TIMING STABILITY: T10 uses a fake net.Conn (helloConn) with no sleeps or sockets; the floor compares the deadline the classifier set against `before+250ms` where deadline >= before+500ms, and the ceiling compares deadline (set at dial) with stop+950ms, so load only widens the margin. A correct tree cannot red from scheduling. The only wall-clock term, margin for mutants (NB3), degrades toward survival, not false red. 500 runs under load green.

SCOPE: production change is exec_probe.go only: constant execProbePollEvery + use, message derived from execProbeHelloWait.Milliseconds(), comment hunks (P1-P4), no logic change; tests only exec_probe_test.go; ci.yml gained exactly TestExecProbeReverseHelloWaitBounds and TestExecProbeChildListenerDrainsTheHello, and the per-test `--- PASS` guard loop is intact. Remaining diff is design/plan/evidence under design_docs.

DISPOSITIONS: items (1) F2 and (2) F3 addressed (T10; T11 structural, honestly labelled, refuted mechanism recorded in the code comment). (3) cross-profile: refuted at HEAD, invariant commented (P4); enforcement left to a successor row; acceptable per design but it is a comment, not a test. (4) R4 watch only. Evidence spot checks reproduced: base sha prefix; 333 PASS/30 SKIP control count; executor-M2 mutant outputs vs my reruns; 335 PASS now.

Deductions: survivor 1 (cadence not pinned, -2 test quality), item 3 left as a comment-only pin (-2 correctness), T11 is textual (-2 test quality), numbers not independently re-measured for V11/V12 srt latency (-3 evidence; not rerun, needs srt/load harness).

git status --porcelain after restore: 0 lines.
