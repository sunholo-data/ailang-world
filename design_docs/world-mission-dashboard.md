# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History: charter STATUS + `world-mission-log.md`.
Written: iteration 118, 2026-08-24.*

## Where we are

- **dev**: GREEN at `bedc0d1`; latest landed milestone **`WB.G`** (7 of item 14's 11).
- **In flight**: item 14 `w-workbench-read-only`. `WB.G` closed **AC2** — the boundary gate now
  measures the workbench's transport-freedom instead of inheriting it.
- **This iteration**: executor `codex:gpt-5.6-sol` landed +129 test-only lines in one file; judge
  `sonnet` **97/100, 0 blocking**; PR #90 → squash `bedc0d1`, Gate 3b green on the merge commit.
- **NEXT**: `WB.H` (8 of 11 — mutation drill 1/4, discharges M14–M21), then WB.I, WB.J, WB.K,
  then queue rows 38, 37, 36, 35, 34, 32, 33, item 22, row 31.

## The find worth knowing

The sprint plan's task 3b required the workbench dependency closure to contain **zero** `net/url`
*and* — correctly, so eight zeros could not be vacuous — **exactly one** `html/template`.
`html/template` transitively imports `net/url`, so the anti-vacuity control is what made its own
sibling assertion false. The executor disclosed it rather than quietly dropping a check, and split
it: pin the transitive count at 1, plus an AST scan forbidding a *direct* import. Two more plan
defects surfaced the same way (an inverted callback polarity, and a `=== RUN` count that go1.25.6
refutes — the second one mine). All three were adjudicated by measurement; all three favoured the
executor. **Nothing in this loop baselines a plan's *task bodies* — only its acceptance commands.**

## Routing · cost · parked

- Controller `claude:claude-opus-5` (driver probed up from `opus-4-8` at 13:05 CEST).
- Executor `codex:gpt-5.6-sol` · evaluator `sonnet` · generator ≠ judge preserved.
- Metered spend **$0.00** of $5. Fable / designer rotation unspent (14th consecutive iteration).
- **Nothing is blocked on a human answer** (`scripts/mission_decisions.sh --open` empty; 0 OPEN).
- Row 5 remains blocked upstream on `ailang#764` (OPEN, 2 comments) — re-measured as a command
  this iteration, not transcribed.
- Standing capacity note (not a park): the only evaluator lane is one weekly-bucketed Sonnet
  subscription. A late-week iteration needing a fresh evaluation can only wait for Monday.
  Routing-capacity signal; the loop cannot widen its own rotation.

## Bookkeeping

- Bookkeeping thread **#89** (rotated by iteration 117 this morning; supersedes #68). Not due
  again until Mon 2026-09-01, or 80 comments.
- 0 `MarkEdmondson1234` directives since the watermark on #89 *or* the predecessor #68, with the
  allowlist script's control firing at 3 historical directives.
