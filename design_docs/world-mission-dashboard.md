# Mission Dashboard — AILANG World

_Snapshot, overwritten each iteration. History lives in `world-mission.md` (STATUS),
`world-mission-status-archive.md` and `world-mission-log.md`._

**As of**: iteration 112 · 2026-08-23 · `dev` == `origin/dev` at `75bc23f`, CI green (`checks=2`)

## In flight
- **Queue row 14 `w-workbench-read-only` — [IN-SPRINT], 2 of 11 milestones done.**
  `WB.B` landed (PR #84 → `75bc23f`, evaluator `sonnet` 91/100, zero blocking).
  `host/workbench` now renders: one parsed `html/template`, landmarks, escaping,
  local-only links, explicit UNAVAILABLE states, dual-channel verdict.
- **NEXT: `WB.C`** — ninth registration `GET /workbench`, `handleWorkbench` happy path,
  security headers, §3.5 comment. It closes **AC5 and AC6**, the first ACs any milestone closes.
- Plan: `design_docs/planned/w-workbench-read-only-sprint-plan.md` (tracked — the machine
  plan under `.ailang/` is gitignored and absent from sprint worktrees).

## Then
rows 34 (three shipped template hunks pinned by nothing — new) → 32 (`host/capsule` load-dependent
red) → 33 (`go test -run` empty-selector AC census) → item 22 → row 31.

## Blocked
- **Row 5** — waits on `sunholo-data/ailang#764` (protocol-only module). Re-measured this
  iteration as a command: `OPEN`, 0 comments, `updatedAt` unchanged since 2026-08-17. Still blocked.

## Loop + routing
- Cadence: launchd `dev.ailang.mission-world`, 6h hard timeout per iteration.
- Controller `claude-opus-5` · executor **`codex:gpt-5.6-sol`** → fallback `opus`
  (D-WORLD-20: DeepSeek link removed) · evaluator `sonnet` (generator≠judge).
- Designer rotation **unspent for 10 consecutive iterations** (no new doc authored).
- Verify profile `ailang-code`; both gates need `AILANG_BIN=/tmp/ailang-v0300/ailang`
  **and** `GOTOOLCHAIN=go1.25.6`.

## Human
- **Parked for Mark: none.** Decision ledger valid, 11 rows, **0 OPEN**.
- Bookkeeping issue **#68** (rotates Mondays 07:00 local; next boundary Mon 2026-08-24).

## Cost
metered **$0.00** of the $5/iteration ceiling. Quota this iteration: `opus` ×1, `codex` ×2,
`sonnet` ×1. Fable unspent.
