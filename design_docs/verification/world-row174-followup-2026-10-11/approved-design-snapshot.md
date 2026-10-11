# w-exec-timeout-stimulus-fixture

**Status:** APPROVED for one bounded test-only follow-up, D-WORLD-71 YES (Mark, attended 2026-10-11); implementation/evaluation pending.
**Scope:** Only the body of `TestExecHandlerTimeoutKillsTheGroup` in `host/broker/handlers_exec_m2_test.go` plus necessary imports. No other tests/helpers or production edits.
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
   Use the existing `profileF` option to override this fixture's `sleep` command with
   fixed `/bin/sh` + helper/token/data-file paths. Keep script/profile/argv0 stable;
   select finite hold through a fixture-owned data file, not executable rewriting.
   Each child atomically publishes its own unique token and checked numeric PID.
2. Run three successful short/no-hold controls through that complete launch path;
   require both streams/text/bytes/SHA, both token-tagged child markers and natural
   exit/reap. Log maximum complete control duration `C`; do not bypass the trampoline.
   Use one standard `boundedTestContext` for all controls and the exercised arm via
   explicit `h.Execute(ctx,...)`; existing `f.call` creates fresh contexts and is not
   used here. Record each full result/child marker; remove old markers between calls.
3. Before building the exercised profile, derive `B=ceil_ms(C+2*pipeCloseGrace)` and
   `H=ceil_seconds(3*B+2*pipeCloseGrace)`; children hold `H`. Require
   `B+H+pipeCloseGrace < remaining ordinary 30 s caller deadline`; otherwise fail
   loudly, without clamp/skip/retry. The fixture profile alone takes `B` milliseconds.
   Also require `B/time.Millisecond <= ExecMaxTimeoutMS` (existing 9000); overflow-safe
   arithmetic and fit checks fail loudly. Rebuild/parse the exercised profile and
   handler through existing constructors if its timeout changes; no direct digest drift.
4. Wrap the existing kill seam only for the exercised call, delegate the real syscall,
   and first verify matching readiness tokens plus **both recorded child PIDs alive
   in the actual leader's process group** (`Getpgid`/signal 0), not a process-name census.
   The callback records synchronized observations/errors and always delegates the
   original real kill; assertions run on the calling test goroutine, not via Fatal
   in exec's watcher. Validate each PID through existing `proctest.Check/Signal`.
   Require real deadline kill, `timed_out`, no exit code, exact two-stream retention,
   then both specific children no longer live. Distinguish zombie/dead state from
   running children; the host owns the direct leader's reap, not an orphan's init reap.
   Existing `proctest.Alive/DeadWithin` count zombies as present: local test closures
   may observe bounded `ps -o stat= -p PID` with checked exit/status plus signal 0.
   Known absent or explicit zombie/dead states pass; unknown/error never means dead.
5. Cleanup targets only the proven fixture PIDs/group with real SIGKILL if necessary,
   waits boundedly for reap/reservation baseline, and restores the global seam after
   callbacks finish. A disabled group kill must never leave the finite children running.
   Cleanup first revalidates token/PID/PGID ownership; only proven fixture targets
   are signalled, never sentinel/ambient PIDs or a global pgrep/pkill pattern. All ps
   observations/polls use finite contexts within the caller guard; no new shared helper.
## Required non-vacuity and delivery

Isolated mutations: **group-kill-disabled** kills only the leader after proven live
children (specific-child-alive assertion red); **stdout/stderr deadline-stream-loss**
leave controls and real readiness/kill intact (each retention assertion red);
**deadline-disabled-natural** removes inner deadline and extends only mutant cleanup
deadline to caller deadline minus grace (finite natural exit/reap/no kill, timeout
assertion red). Restore exact bytes and green controls after each; bank actual assertion.
No custom context, capture preseed, omitted children, assertion weakening or production
helper/bounds edit. All caps, floor tasks/models/thresholds/N and shared harness freeze.
Approved next step: separate test-only commit, scoped plain/race controls and mutations,
fresh full normal Go gate, then new independent judge; rows remain open until accepted.
If accurate PID/group readiness or finite cleanup needs production support, stop this
proposal and route a normal new design rather than inventing a pass. All quorum triggers
false: no freeze item, no shared override, no KPI/banked schema change, no external
premise. Record quorum SKIPPED, never PASS. Why host rather than package: this is an
existing Go lifecycle fixture, not a semantic addition. Axioms A5/A7 +1, other ten 0,
net +2, no hard violation. The approved one follow-up does not authorize further rounds.

## Verification log / 2026-10-11 designer handoff

| Evidence read | Verified scope |
|---|---|
| CLAUDE, coding-standards, charter D-WORLD-71 / row174; design-doc-creator skill | One approved follow-up; mutations/gates; prior three FAIL68 banks retained |
| `handlers_exec_m2_test.go` fixture/profileF/call + exact timeout test | Existing fixed launch path; test-only override available; call creates fresh context |
| `handlers_exec.go` Execute/launch; `exec_profile.go` | Verification/cache before runBounded; normal trampoline; profile limit9000 unchanged |
| `handlers.go` runBounded/runHeadTail; `procbound.go` | Deadline/group kill and joined pipes/reap; reservation baseline and late-seam risk |
| `host/proctest/proctest.go` bodies | Checked signals reusable; Alive/DeadWithin cannot alone distinguish zombies |
| `handlers_parallel_guard_test.go` | Broker tests guarded against parallel calls; do not introduce parallelism |

No fresh control, mutant, full gate or judge result is represented by this designer
handoff. Evidence goes in `verification/world-row174-followup-2026-10-11/`; preserve
earlier reviews and FINAL unchanged. Runtime model attestation unavailable; native
fallback FLAG under D-WORLD-48, not a cross-vendor quorum. See designer receipt there.
