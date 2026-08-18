# Mission Dashboard — AILANG World

**Snapshot 2026-08-18, after iteration 90.** Overwritten each iteration; history lives in
`world-mission.md`'s STATUS block, `world-mission-status-archive.md` and `world-mission-log.md`.

## State
- **dev**: `03c7892` + this iteration's record. Gate 3b green on `03c7892`
  (`present=2 == expected=2`, both `success`).
- **Ledger**: `mission_decisions.sh --check` → valid, **6 rows**; **1 OPEN** (`D-WORLD-19`).
- **Verify profile**: `ailang-code`. Pinned verifier v0.30.0. Gate pins 10 identities /
  39 named tests / world package 9-9.

## In flight / next
- **`[NEXT]` item 18** `w-daemon-read-cancellation` — the only fully-unblocked routable row.
  D-WORLD-18 arm A, ratified attended: straight to **sprint-planner**, no designer round.
- **Item 17** `w-validated-proven-evidence-boundary` — PARKED on `D-WORLD-19` after quorum
  rounds 5 and 6 (doc 711 → 968 lines). Direction never re-disputed; the open question is
  whether tranche 1 may extend `host/store` with a bounded object read.
- **Item 5** `w-mcp-projection` — re-blocked on upstream
  [`ailang#764`](https://github.com/sunholo-data/ailang/issues/764), not on a human. Arm A's
  dependency condition fails by measurement: importing `serveapi` adds **476 disallowed
  packages across 86 module roots** to a daemon graph that today is 46 non-stdlib packages
  across exactly the 11 allowed roots.
- **Item 4e** — its parked remediation is newly unblocked: **`go1.26.6` fixes** the array-literal
  miscompile (`go1.26.5` BUG / `go1.26.6` OK / `go1.25.6` OK, both repro controls firing).

## Parked on Mark — 2 items, neither is an investigation
1. **`D-WORLD-19`** — one word. **A**: tranche 1 may add a bounded `host/store` object read
   (closes an OOM vector; puts a second item into the package item 18 is queued to bound).
   **B**: record it as a named limitation of an explicitly non-production tranche; the bounded
   read lands with item 18 or the tranche-2 wiring item.
2. **The fleet-authored driver commit** — two commands, not a decision. `tools/launchd/*`,
   `scripts/verify_go.sh` and `CLAUDE.md` are still uncommitted from the 2026-08-17 attended
   session; D-WORLD-DRIVER-1 arm B assigns the commit to the fleet's human, and the harness
   refuses the controller the same staging. Backed up sha256-verified at
   `~/.ailang/state/world-driver-backup-2026-08-17/`.

## Loop cadence + routing
- Controller `claude:claude-opus-5` · designer rotation now at `claude:claude-fable-5`
  (advanced past `codex:`) · planner `opus` · executor `pi:deepseek-v4-flash` · evaluator `sonnet`.
- **Codex is exhausted fleet-wide until 2026-08-20 05:34** — the rotation's `codex:` entry
  probe-fails; fall to the next entry, never `$MODEL`.
- Spend: iteration 90 `metered=$0.328765` of the $5 ceiling (four quorum reviewer runs).
