# Mission Dashboard — AILANG World

_Snapshot, overwritten every iteration. History lives in `world-mission.md` (STATUS),
`world-mission-status-archive.md` (rotated stamps) and `world-mission-log.md` (full entries)._

**As of** 2026-08-20 · **iteration 100** · controller `claude:claude-opus-5`

## State
- **Latest landing**: queue row 20 `w-capsule-output-cap-load-flake` — PR #74 → squash `912009d`,
  Gate 3b GREEN on the **merge** commit (`present=2 == expected=2`), evaluator `sonnet` 93/100,
  zero blocking.
- **dev**: green, `912009d`, both checks `success`. Local checkout == `origin/dev`.
- **Decision ledger**: **10 rows, 0 OPEN** — nothing is parked on Mark for the first time in this
  mission's history.

## In flight / next
1. **row 17** `w-validated-proven-evidence-boundary` — UNPARKED, no ask open. Owed: a bounded
   revision narrowing the claim + an assertion pinning `busy_timeout` < `ObjectReadTimeout`.
   Priced 4.75 d vs a 3–4 d guardrail; re-scope is part of the revision.
2. **row 22** `w-daemon-lock-wait-not-deadline-bound` · **row 23** `w-store-deadline-free-residue-owner`
   · **row 24** `w-host-subprocess-cleanup-boundary` · **row 25** `w-capsule-blocked-child-kill-coverage`
   (new this iteration) — all open, all headless-routable.

## Loop / routing
- Cadence: launchd, driver `tools/launchd/mission-control.sh` (FLEET-owned, `D-WORLD-DRIVER-1`).
- Controller `claude:claude-opus-5` · planner `opus` (`fail-closed:env-pin`) · executor
  `codex:gpt-5.6-sol` · evaluator `sonnet` · designer rotation currently `codex:gpt-5.6-sol`.
- `pi:deepseek` lane SUSPENDED by `D-WORLD-20`. Fable lane unspent 3 consecutive iterations.
- Verify profile `ailang-code`: `scripts/verify_ail.sh` + `scripts/verify_go.sh`.
  **Ambient Go is the DENIED go1.26.4 — always pin `GOTOOLCHAIN=go1.25.6`.**

## Cost
- Iteration 100 `metered=$0.00` of the $5 ceiling (no quorum bought — the ruling replaced round 3).
- Quota buckets: `opus` ×2, `codex` ×2, `sonnet` ×1.

## Parked on Mark
**None.** `D-WORLD-22` and `D-WORLD-23` both resolved 2026-08-20 by one attended one-word answer.
