# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History: charter STATUS + `world-mission-log.md`.
Written: iteration 117, 2026-08-24.*

## Where we are

- **dev**: GREEN at `8f0037c`; latest landed milestone **`WB.F`** (6 of item 14's 11).
- **In flight**: item 14 `w-workbench-read-only`, milestone `WB.G` is NEXT (7 of 11).
- **This iteration**: resumed iter-116's `PARKED-ON-LANE` `WB.F`. Sonnet re-probed available
  after the 07:00 reset, evaluated `a96fd67` **96/100 PASS, 0 blocking**; landed via PR #88 →
  squash `8f0037c`, Gate 3b green on the merge commit.
- **NEXT**: `WB.G`, then queue rows 38, 37, 36, 35, 34, 32, 33, item 22, row 31.

## Routing · cost · parked

- Controller `claude:claude-opus-4-8`; executor (WB.F) was `codex:gpt-5.6-sol`; evaluator `sonnet`.
- Generator ≠ judge preserved. Metered spend **$0.00**. Fable/designer rotation unspent (13th iter).
- **Nothing is blocked on a human answer** (`scripts/mission_decisions.sh --open` empty; 0 OPEN).
- Row 5 remains blocked upstream on `ailang#764` (OPEN, 2 comments; V1's fix planned, not landed).
- Capacity note (not a park): the only evaluator lane is one weekly-bucketed Sonnet subscription;
  a late-week fresh evaluation can only wait for Monday. Routing-capacity signal, human-only fix.

## Bookkeeping

- Bookkeeping thread **rotated** this iteration (Mon 2026-08-24 past 07:00 local). New issue supersedes #68.
- Weekly external-issue sweep: 1 open issue (the bookkeeping thread itself), 0 orphans of 1.
