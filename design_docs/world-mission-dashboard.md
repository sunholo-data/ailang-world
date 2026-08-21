# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS), `-status-archive.md`, `-log.md`.*

**Iteration 104** · 2026-08-21 · `dev` @ `3ddacae` · CI green (both jobs, SHA-addressed, run confirmed)

> **The judge scored 91/100, zero blocking — and two of its NON-blocking findings were this
> milestone's own anti-vacuity contract failing.** Both reproduced first-party, both fixed; the
> first remedy for the second was *itself* vacuous and its own drill caught it.

## In flight
- **Item 17 `w-validated-proven-evidence-boundary`** `[IN-SPRINT]` — six milestones `PE.A`–`PE.F`
  (4.70 d), **two landed**.
- **`PE.B` LANDED** — PR #77 → `3ddacae`, Gate 3b green on the merge (`present=2 == expected=2`,
  run confirmed). `ReadObject` (one snapshot for probe and payload), cached `BusyTimeout()`, B5.
- **NEXT: `PE.C`** (0.80 d) — `host/evidence` codecs, byte caps, nesting-depth pin. Then `PE.D` →
  `PE.E` → `PE.F` in compile order; `PE.F` last (its `EXACT_EVIDENCE_TESTS` pin forces it).

## The finding worth carrying
`busyTimeoutFromParams` reported the **last** `busy_timeout` pragma in a DSN; the pinned driver
applies the **first** (measured both directions against a live `PRAGMA` readback). Reachable, since
`withBusyTimeout` deliberately never overrides a caller's value — and under-reporting is the unsafe
direction for AC18/AC22's `ObjectReadTimeout > BusyTimeout()` pin. Fixed to first-wins, pinned
against the readback rather than against a comment.

## Queue after item 17
Rows **22**, **23** headless-routable · **24**–**27** designed-pending · **28**/**29** new, both
from the judge and both declined as silent edits: `ReadObject`'s six unpinned refusal branches
(sibling `GetObject`'s analogous branches *do* die), and its absent branch distinguishable from a
zero-length payload only by an invariant the *writers* enforce. Row **5** stays blocked —
`sunholo-data/ailang#764` re-measured OPEN today, 0 comments, untouched since 2026-08-17.

## Parked on Mark
**NONE.** Decision ledger: **11 rows, 0 OPEN** (`scripts/mission_decisions.sh --check`).

## Loop / routing / cost
Controller `claude:claude-opus-5` · no planner/designer (plan existed) · executor
`codex:gpt-5.6-sol` (probe rc=0) · judge `sonnet` in its OWN worktree (generator≠judge). Fable
unspent a **5th** iteration. `metered=$0.00` of $5 · quota `opus` ×1 / `codex` ×2 / `sonnet` ×1.

## Gates
`AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` ·
`AILANG_BIN=… GOTOOLCHAIN=go1.25.6 ./scripts/verify_go.sh` — **both exports mandatory**; the go gate
fails closed without `AILANG_BIN` and refuses the rig's default go1.26.4 outright. Pinned `v0.30.0`.
`TestHandlerTimeoutKillsTheWholeProcessGroup` (`host/broker`) is load-flaky — charter-recorded 2/5.
