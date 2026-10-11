# w-exec-timeout-stimulus-fixture

**Status:** PROPOSED — attended approval required for one bounded follow-up beyond three failed reviews.
**Scope:** Only `TestExecHandlerTimeoutKillsTheGroup` in `host/broker/handlers_exec_m2_test.go`; no implementation authorized by this proposal.
## Evidence and mechanism

At `84d3f5a`, normal/plain Go failed this unchanged fixture: profile 800 ms,
reported timeout around 802 ms, empty stdout/stderr; race was not reached.
The revised row166 fixture passed. Root's isolated untouched-parent control: 20 PASS/0 FAIL, rc0, 17.433 s (`row166/baseline-exec-timeout-result.json`); preserve both results, with no cause inference or gate waiver.
No unique startup, capture, or scheduler cause is established by those observations.
The read test launches two `sleep 3011` children, then `echo started`; a timeout
can precede that stimulus, and a zero global `pgrep` census can precede any child start.
Code-read: `ExecHandler.Execute` verifies stack/settings, prepares cache and launches
stub node → trampoline `/bin/sh` → `env -i` → configured command → two children;
`runBounded` owns the real deadline/group kill and `Wait` joins capture/reap.
These tests use a pass-through srt stub, not real sandbox confinement.

## Proposed test-only repair

1. In this test's fixture only, replace `sleep 3011` with a unique temporary helper
   invoked through the same ExecHandler/profile/stub/full trampoline. Its two real
   child helpers each write a token/PID readiness marker before a finite hold;
   leader emits exact stdout and stderr markers after both child markers, then waits.
2. Run three successful short/no-hold controls through that complete launch path;
   require both streams/text/bytes/SHA, both token-tagged child markers and natural
   exit/reap. Log maximum complete control duration `C`; do not bypass the trampoline.
3. Before building the exercised profile, derive `B=ceil_ms(C+2*pipeCloseGrace)` and
   `H=ceil_seconds(3*B+2*pipeCloseGrace)`; children hold `H`. Require
   `B+H+pipeCloseGrace < remaining ordinary 30 s caller deadline`; otherwise fail
   loudly, without clamp/skip/retry. The fixture profile alone takes `B` milliseconds.
4. Wrap the existing kill seam only for the exercised call, delegate the real syscall,
   and first verify matching readiness tokens plus **both recorded child PIDs alive
   in the actual leader's process group** (`Getpgid`/signal 0), not a process-name census.
   Require real deadline kill, `timed_out`, no exit code, exact two-stream retention,
   then both specific children no longer live. Distinguish zombie/dead state from
   running children; the host owns the direct leader's reap, not an orphan's init reap.
5. Cleanup targets only the proven fixture PIDs/group with real SIGKILL if necessary,
   waits boundedly for reap/reservation baseline, and restores the global seam after
   callbacks finish. A disabled group kill must never leave the finite children running.
## Required non-vacuity and delivery

Isolated mutations: **group-kill-disabled** kills only the leader after proven live
children (specific-child-alive assertion red); **stdout/stderr deadline-stream-loss**
leave controls and real readiness/kill intact (each retention assertion red);
**deadline-disabled-natural** removes inner deadline and extends only mutant cleanup
deadline to caller deadline minus grace (finite natural exit/reap/no kill, timeout
assertion red). Restore exact bytes and green controls after each; bank actual assertion.
No custom context, capture preseed, omitted children, assertion weakening or production
helper/bounds edit. All caps, floor tasks/models/thresholds/N and shared harness freeze.
After approval: separate test-only commit, scoped plain/race controls and mutations,
fresh full normal Go gate, then new independent judge; rows remain open until accepted.
If accurate PID/group readiness or finite cleanup needs production support, stop this
proposal and route a normal new design rather than inventing a pass. All quorum triggers
false; test-only host machinery is not a semantic package addition. Axioms A5/A7 +1,
other ten 0, net +2 with no hard violation. This proposal does not waive the three-review limit.
