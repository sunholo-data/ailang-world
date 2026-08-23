# Mission Dashboard — AILANG World

_Snapshot, overwritten each iteration. History lives in `world-mission.md` (STATUS),
`world-mission-status-archive.md` and `world-mission-log.md`._

**As of**: iteration 113 · 2026-08-23 · `dev` == `origin/dev` at `5fd6fb3`, CI green (`checks=2`)

## In flight
- **Queue row 14 `w-workbench-read-only` — [IN-SPRINT], 3 of 11 milestones done.**
  `WB.C` landed (PR #85 → `5fd6fb3`, evaluator `sonnet` 82/100, **1 blocking, fixed**).
  `GET /workbench` is now the ninth mux registration; `handleWorkbench` reads through the
  request-scoped deadline and sets the five security headers on every response, errors included.
  **AC5 closed** (route cardinality, base rc=1 → 0). **AC6 did NOT close** — it is rc=0 at base
  and is carried as a regression pin; see the STATUS stamp.
- **NEXT: `WB.D`** — the closed query grammar and every refusal branch (claims M2–M9, M13,
  M31, M32). `WB.C` shipped the minimal unknown-key guard only, by controller adjudication, and
  `WB.D` extends that one site rather than rewriting it.
- Plan: `design_docs/planned/w-workbench-read-only-sprint-plan.md` (tracked — the machine
  plan under `.ailang/` is gitignored and absent from sprint worktrees).

## Then
rows 35 (three populated `EntryView` fields rendered by nothing — new) → 34 (three shipped
template hunks pinned by nothing) → 32 (`host/capsule` load-dependent red) → 33 (`go test -run`
empty-selector AC census) → item 22 → row 31.

## Blocked
- **Row 5** — waits on `sunholo-data/ailang#764` (protocol-only module). Re-measured this
  iteration as a command: `OPEN`, 0 comments, `updatedAt` unchanged since 2026-08-17
  (control `#676` answers 3 comments / a different `updatedAt`; negative control errors).

## Loop + routing
- Cadence: launchd `dev.ailang.mission-world`, 6h hard timeout per iteration.
- Controller `claude-opus-5` · executor **`codex:gpt-5.6-sol`** → fallback `opus`
  (D-WORLD-20: DeepSeek link removed) · evaluator `sonnet` (generator≠judge).
- Designer rotation **unspent for 11 consecutive iterations** (no new doc authored).
- Verify profile `ailang-code`; both gates need `AILANG_BIN=/tmp/ailang-v0300/ailang`
  **and** `GOTOOLCHAIN=go1.25.6`.

## Human
- **Parked for Mark: none.** Decision ledger valid, 11 rows, **0 OPEN**.
- Bookkeeping issue **#68** (rotates Mondays 07:00 local; next boundary Mon 2026-08-24).

## Cost
metered **$0.00** of the $5/iteration ceiling. Quota this iteration: `opus` ×1, `codex` ×2,
`sonnet` ×1. Fable unspent.
