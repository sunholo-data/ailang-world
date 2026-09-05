# Mission Dashboard — Ailang World

**Snapshot: 2026-09-05, after iteration 156.** Overwritten every iteration; history lives in
`world-mission.md` (STATUS), `world-mission-status-archive.md`, `world-mission-log.md`.

## State
- `dev` == `origin/dev` == `3417088`. CI **GREEN 3/3** on that merge commit.
- Gate = **two legs**: `AILANG_BIN=~/.pinned-ailang/ailang ./scripts/verify_ail.sh` and
  `go build ./... && go test ./... -count=1`. Pin **AILANG v0.30.0 / `e37b370`**.
  `verify_go.sh` is **rc=1 at base** on a FLEET-OWNED drift arm (row 76) — use the two legs.
- Bookkeeping issue **#107** (prev #89). Decision ledger **18 rows, ZERO OPEN**.

## Just landed (iteration 156)
- **Row 60 LANDED** — PR #117 → squash `3417088`. The P1 needles now bind each operand's
  **derivation**, not its spelling, so a semantically inert rename no longer reds CI.
- The row named **one** identifier-pinned needle; there were **three**, each invisible until the
  one before it was relaxed. The predecessor was also **fail-open** on the inversion it is
  documented as the sole killer for (reachable by swapping the assignments; old count still 1).
- The green control is now a test over the file's text, not a drill row — a needle added later is
  covered without anyone widening the arm.
- **Controller-authored; no independent judge ran** (generator == judge). Compensating: 8
  file-level mutants + a 7-conjunct sensitivity drill, all landed/restored by sha256.

## Next picks
**Row 61** (head), then **62–66**, **68–78**, **81**, **82**, then **39**. Rows **79/80** (Astra)
are `[PARKED — DESIGN REVIEW]`; they need their own quorum and do not reorder release work.

## Parked on Mark
**Nothing.** Zero open decisions.

## Loop / routing / quota
Controller `claude:claude-opus-5` · planner `opus` · executor `codex:gpt-5.6-sol` · evaluator
`sonnet` · designer = rotation (attended 2026-09-05: fable-5-1 → `codex:gpt-6-astra` → pi/deepseek;
astra ADDED, fable not displaced). Iteration 156 spent **`metered=$0.00`** of $5.

## Known local gaps (not defects in this repo's code)
`mission_directives.sh`, `mission-heartbeat.sh`, `resolve-role-spawn.sh`, `mission_pi_run.sh` are
**absent here** — reached by ABSOLUTE PATH into the V1 checkout (row 69).
`tools/launchd/mission-control.sh.tmp.astra` is an untracked **fleet** artifact (frozen core),
left alone deliberately, second iteration running.
