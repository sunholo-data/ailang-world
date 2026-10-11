# Designer receipt — approved exec timeout fixture follow-up

**Date:** 2026-10-11. **Role:** designer, native fallback FLAG under D-WORLD-48;
exact runtime model/token attestation unavailable. No provider calls or quorum run.

**Ready design:** `design_docs/planned/w-exec-timeout-stimulus-fixture.md`.
Current read head `5f6ebfd`, merging decision-only `948e115` and `origin/dev`;
charter D-WORLD-71 YES authorizes one bounded follow-up through all four roles.
Iteration250/three prior FAIL68 reviews stay banked, not reset or superseded by this receipt.

Reviewed current CLAUDE, charter authority/row174, coding standards and previously
loaded design-doc-creator skill; read exact exec fixture, profile constructor/limit,
ExecHandler Execute/trampoline, runBounded/runHeadTail/procbound and proctest bodies.
Only design/status/this receipt changed by designer; no product code or git writes.

Approved implementation surface: existing timeout test body + necessary imports only.
Local closures and a fixture-created helper script/data files are permitted; shared
helpers, other tests, production support/bounds, floor/harness/FINAL changes are not.

Pinned implementability details:

- `profileF` supplies the exact fixed `/bin/sh` helper command through the existing
  archived stub/full trampoline. Stable script/token/path arguments; hold mode is
  fixture data. Controls do not call an alternate direct shell entry point.
- Three full-path positive controls produce exact both-stream text/bytes/SHA and
  atomically published child-owned token/PID markers, then naturally exit/reap.
  Reuse one ordinary 30 s caller context via explicit `h.Execute`; do not use
  `f.call`'s per-call fresh contexts. Preserve first-control cost in maximum `C`.
- `B=ceil_ms(C+2g)`, `H=ceil_seconds(3B+2g)`, require `B+H+g<R`, `B<=9000ms`,
  overflow-safe arithmetic, no clamp/skip/extra retries. Reparse exercised profile
  and construct handler normally when timeout changes; no silent digest mismatch.
- At real kill: validate both unique tokens/PIDs, actual same leader PGID and
  existing/running children, then delegate original SIGKILL. Record observations
  with synchronization; test goroutine owns assertions, never callback Fatal.
- `proctest.Check/Signal` guard every positive PID. Alive/DeadWithin alone count
  zombies as present. Bounded `ps`/exit check plus signal 0 distinguishes absent,
  explicit zombie/dead and live states; observation error/unknown fails loudly.
- Cleanup revalidates ownership before any fixture-only kill, observes finite
  cleanup/reservation baseline and callback completion before global seam restore.
  No global name census, sentinel PID, ambient group or unbounded ps/poll wait.

Planner must pin four isolated mutations and their single reached assertions:
leader-only kill leaves proven live children (liveness red); deadline-only stdout
loss and stderr loss leave controls/readiness/real kill green (retention red);
coherent internal-deadline removal with mutant-only cleanup deadline at caller
deadline minus grace allows finite natural exit/reap/no kill (timeout red).
Restore exact bytes/green after each, then unchanged scoped plain/race, full normal
gates and fresh independent judge. A new defect/support need stops this follow-up;
there is no authorization for endless repairs or review-round resets.

Quorum: SKIPPED, all four attended triggers false (no freeze decision, shared override,
KPI/banked schema change or external premise). This is not a quorum PASS.
No efficacy, M1 acceptance or row/FINAL closure is claimed. User-typed publish/mint
and benchmark model availability remain distinct later gates for row174 smoke.

Material risks: loaded/cold controls may not fit both finite guards; missing child
markers must remain a failure; `ps` portability/error classification must be explicit;
zombie presence is not live execution; global seam restoration must follow actual
cleanup. If any cannot be discharged inside this test body, stop for a new design.
