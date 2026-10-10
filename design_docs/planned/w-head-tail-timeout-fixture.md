# w-head-tail-timeout-fixture

**Status:** Planned prerequisite to row174 M1 normal gate; attended 2026-10-10.
**Scope:** Existing row166 hygiene item, test fixture only; ~0.25 day including mutations.
**Authority:** Continue necessary local gate repair; no mission-state or benchmark changes.

## Problem and verified mechanism

`TestRunBoundedHeadTailTimeoutKeepsPartialOutput` starts a shell that prints
`PARTIAL-OUT`/`PARTIAL-ERR` then sleeps 30 seconds, with a 300 ms execution budget.
The untouched parent `5ca1115` produced 13 PASS/7 FAIL at count20 with empty stdout;
candidate normal/plain and unchanged retry/race failed the same assertion.
These are [banked baseline observations](../verification/world-row174-attended-2026-10-10/baseline-partial-output-control.log), not a new candidate regression.

Read `handlers.go` `runBounded`/`runHeadTail`, `procbound.go`, and the exact test body:
`context.WithTimeout` precedes command construction, reservation and `cmd.Start`.
`Admit` is nonblocking, so a queued reservation wait is not the cause.
Both sinks attach before Start; `cmd.Wait` joins exec's real pipe-copy goroutines.
An expired context before Start is deliberately classified as handler timeout.
Thus timeout alone does not prove the child reached its output stimulus.
The current evidence does not identify a unique cause (startup delay vs another path).

## Smallest proposed repair

Keep the same `runBounded` entry point, real `/bin/sh` child, capture strategy,
standard context deadline, process-group SIGKILL and exact two-stream assertions.
Change only the fixture so its setup is measured and its child state is observed:

1. Use one executable that prints both exact markers, writes a readiness file only
   **after both printf operations**, then optionally sleeps for an environment-selected
   duration. Readiness is independent side-channel evidence, not captured-output preseed.
2. Run three no-sleep positive controls through **the same `runBounded` executable**
   with fresh real captures and finite existing 20 s test allowances. Require exact
   stdout/stderr, byte counts and hashes, successful exit and the readiness file.
   Log each total control duration and its maximum `C`; remove readiness between runs.
3. Derive the exercised timeout as `B = C + 2*pipeCloseGrace`, rounded upward to a
   whole millisecond. `C` already measures complete same-executable startup, both
   writes, natural exit and drain; adding two existing cleanup-grace units supplies
   finite scheduling headroom without multiplying all that work by an unsupported
   tenfold factor. This tighter fixture deadline changes no production handler bound.
4. Derive child hold time `H` as the whole-second ceiling of `3*B + 2*pipeCloseGrace`.
   Require `B + H + pipeCloseGrace < remaining outer deadline`: include startup,
   finite hold and reap headroom inside the ordinary 30 s caller context. If calibration cannot fit, fail
   loudly with measured values; never clamp, skip, or retry into a pass.
5. For the exercised arm only, wrap the existing `killGroup` test seam, delegate to
   the original real syscall, and record whether readiness existed at the kill.
   Restore the seam only after bounded cleanup/reap; synchronize observations. Do not parallelize this test.
   Require at least one real kill, readiness before that kill, `ErrHandlerTimeout`
   with the exercised `B`, and exact retained text/bytes/hash on **both** streams.

This asserts the child's state at cancellation instead of inferring it from 300 ms.
Readiness missing at kill is an explicit fixture failure, not a weakened partial-output pass.
Calibration executes a positive-control branch openly; the exercised branch still waits
for the actual `runBounded` deadline. No custom/resettable context or fake timeout error.

## Acceptance and mutations

- **AC1:** Three real controls and one real timeout arm complete; both output markers,
  byte counts and SHA-256 values match. Mutate either sink to discard output: its
  exact stream assertion must fail (including zero-output null case).
  Additionally run separate stdout-only and stderr-only deadline-arm stream-loss mutants: no-sleep controls stay green, real deadline kill/readiness is reached, and that stream's exercised retention assertion turns red; restore each to green.
- **AC2:** Exercised child reaches readiness before real deadline kill; returned error
  is timeout with `B`. Use coherent temporary **MUT-TIMEOUT-DISABLED-NATURAL**: replace the
  internal `WithTimeout` with standard `WithCancel`, and extend only the mutant's
  `runHeadTail` cleanup deadline to the existing caller deadline minus one
  `pipeCloseGrace` (require that deadline present; `B+H+grace<R` places reap before it).
  The finite child exits with no kill, real Wait joins streams/reaps, and timeout/kill assertion
  must fail before the outer guard. Both temporary product edits restore exactly.
- **AC3:** Delete either child printf while retaining readiness: that stream assertion
  fails. Remove readiness creation: readiness-before-kill assertion fails.
- **AC4:** Every exercised arm records pre-run `procbound.Outstanding` and verifies
  reservations return to that baseline before restoring the global kill seam.
  Mutation runs are isolated focused invocations; no concurrent broker tests.
  Cleanup uses a finite context/guard and the original real group kill if needed,
  waits for reap, and fails loudly on incomplete cleanup. Restore exact product/test
  bytes after every mutant and rerun green; preserve
  baseline red logs. Normal focused/race controls and complete required Go gate pass.

Only `host/broker/handlers_capture_test.go` changes (~35–60 LOC); production files,
3/10/20/30-second handler/transport caps, benchmark tasks/models/thresholds/N and
shared harness remain unchanged. Why host rather than package: this is a Go process
capture regression fixture, not a World semantic feature; pure kernel remains untouched.
Verification still pending: executor must bank measured `C/B/H`, readiness-at-kill,
all mutations, and full gate before declaring repaired. Calibration is not causal proof.
**Round2 formula correction (2026-10-10):** retain both controller-reported pristine
guard-red logs: cold-control `C=2.353960541s` and normal/plain `C=756.807791ms`, each
after all three controls passed. Old `10*C+2g` failed the finite-fit guard before
the exercised arm. Revised rounding gives respectively `B=4.354s,H=16s,B+H+g=21.354s`
against `R=27.636s`, and `B=2.757s,H=11s,B+H+g=14.757s` against `R=29.234s`.
These arithmetic fits are not passing executions or a unique startup-cause claim;
later standard-cache green does not erase either red. Fresh round3 evaluation,
all eight isolated mutants with exact restored greens, repetitions and full gates
remain required; readiness, retention, actual timeout and cleanup assertions stay fixed.
Single-site deadline removal is **not** the AC2 mutant: unchanged cleanup bookkeeping
can return `ErrCleanupIncomplete` at `B+2*pipeCloseGrace`, after which deferred cancel
can invoke a real kill and reap later. Never claim that runner naturally waits for exit.
MUT-TIMEOUT-DISABLED-NATURAL isolates execution-deadline enforcement from that second cleanup
timer; the widened finite wait is mutation-only, never a landed bound change.
Related: row166; coding-standards S6; mission rule3m; `cleanup_test.go` existing real-kill seams.
All attended quorum triggers false: no freeze decision, shared override, KPI/schema change
or external premise. Axioms A1–A12: `[0,0,0,0,+1,0,+1,0,0,0,0,0]`, net +2;
bounded state-observed controls improve verification/machine evidence, no hard violation.
