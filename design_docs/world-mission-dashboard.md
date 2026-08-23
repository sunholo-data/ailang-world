# Mission Dashboard — Ailang World

*Snapshot. Overwritten every iteration. History lives in the charter STATUS block and
`world-mission-log.md` — never here.*

**Last iteration**: 114 · 2026-08-23 · `WB.D` LANDED

## Where we are
- **Queue item 14** `w-workbench-read-only` **[IN-SPRINT]** — **4 of 11 milestones landed**
  (`WB.A` `83f1973` · `WB.B` `75bc23f` · `WB.C` `5fd6fb3` · `WB.D` `e50fbea`).
- **dev is GREEN** at `e50fbea`: `checks=2`, both `success`, 0 not-green, run existence asserted.
- **NEXT**: `WB.E` — payload opt-in, 64 KiB cap, 100-entry timeline cap (claims M10–M12).
- Then: row 37 → row 36 → row 35 → row 34 → row 32 → row 33 → item 22 → row 31.

## This iteration in one line
`WB.D` enforced §2.4's closed query grammar and every refusal branch (14 pinned subtests). The
finding: **`M9`'s mutation row names one mutant for a pattern with NINE call sites, and 7 of the 9
have no killer anywhere in the package** — filed as queue row 37, not absorbed.

## Loop cadence + routing
- Controller `claude:claude-opus-5` · executor `codex:gpt-5.6-sol` (probe rc=0, first try)
  · evaluator `sonnet` 82/100 PASS · planner `opus`.
- Designer rotation **unspent a 12th consecutive iteration** (no new doc needed since iter-102).
- `metered=$0.00` of the $5 ceiling. Quota: `opus` ×1, `codex` ×2, `sonnet` ×1.

## Parked on Mark
- **NONE.** Decision ledger: 11 rows, **0 OPEN**. Zero open asks this iteration.

## Blocked / external
- Queue **row 5** stays BLOCKED on upstream `sunholo-data/ailang#764` (re-measured as a command
  this iteration: `state=open comments=0 updatedAt=2026-08-17T23:34:55Z`, unchanged).

## Bookkeeping
- Issue **#68** (week of 2026-08-17), 33 comments, cap 80. Rotation next due **Mon 2026-08-24
  07:00 local**. Watermark `2026-08-23T07:50:00Z`; 0 directives since.
