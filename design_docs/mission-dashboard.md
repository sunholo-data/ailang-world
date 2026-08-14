# AILANG World — mission dashboard

*Snapshot, overwritten every iteration. History: charter STATUS + `world-mission-log.md`.*
**As of** 2026-08-14 (iter-85) · **dev** `aaada20` · **CI** green on the merge commit, SHA-addressed,
`present(2) == expected(2)`.

## Just landed — item 15 `w-decision-lifecycle-freeze`

PR [#67](https://github.com/sunholo-data/ailang-world/pull/67) → `aaada20`. Evaluator `sonnet`
**96/100, ZERO blocking**, reproducing 10 of 16 mutation arms in its own worktree. The v1
`DecisionPacket` is frozen: three types, **five Z3-proven laws**, 39 named tests; pins moved in the
same commit — **5→10 identities, 20→39 tests**. `metered=$0.00`. The landing is **six files, not
five**: `docs/SELF_MOD_PUBLISH.md` quotes `contentHash`/`tarballSHA256` verbatim in its attended-
operator approval table, so reprojecting red-lit `host/runbook` (rc=0 pristine, rc=1 pre-repair).

## The iteration's real subject: a SECOND dead slot, and the fix didn't hold

Iteration 85 **attempt 1** did the design work, planned, spawned the executor — and died ending its
turn to wait. `rc=0`, no watchdog, no charter row. Worse than the 2026-08-07/08 pair on two counts:
World's driver **does** carry `CLAUDE_CODE_PRINT_BG_WAIT_CEILING_MS=0` (`mission-control.sh:55`), so
this is a **post-fix** orphan; and that ceiling **suppresses the reap line the skill's own
attribution tell greps for**, so the prescribed instrument reads clean (it also over-counts: 3, one
being a log record quoting the count). Attribute by shape — `rc=0` far too fast, a transcript
ending in an announced wait, an orphan worktree holding real work.

## What each item needs now

- **18 `w-daemon-read-cancellation`** — `[NEXT]`, needs a design doc; blocks item 14.
- **17** — *revision round*: MAC seam, V27 repair (`verify.results[]`, never the int), neg control.
- **14** — blocked on 18 (`err.Error()` leaks to an unauthenticated client). **5** — still blocked.
- HUMAN-SURFACE §7 **points 1 (timeout half) and 3 CLOSED**; point 1's "binding" half open.

## Loop · carry-forward

launchd, ~6h watchdog. Controller `opus` · evaluator `sonnet` · rotation still at
`codex:gpt-5.6-sol` (designer did not fire). Cap $5; spent **$0.00**. **FLAGGED**: M2/M3/M4 were
driven by the controller, not the pinned codex executor — deliberate, because attempt 1 died at
that exact spawn; the evaluator judged the evidence not weakened. **An arm can be decoration** —
MU15: `x <= y-1` ≡ `x < y` over ints, so it still refuses the input it was written to admit and
kills nothing; the plan pre-registered exactly that. **Zero open asks.**
