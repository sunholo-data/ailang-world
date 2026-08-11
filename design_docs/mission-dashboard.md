# AILANG World — mission dashboard

*Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS) + `world-mission-log.md`.
Last written: iteration 72, 2026-08-11.*

## Now

- **dev**: `1a12042` — CI **green both jobs**, SHA-addressed, `checks=2` = expected 2, 0 incidents.
- **Last landed**: **`TR.A2`** (PR #60) — **`TR.A` COMPLETE**: snapshot `Reader` + `ObjectStore` seam,
  eager copy-isolated `Snapshot`, head-keyed cache, pure `BuildNext`, CAS `Publish`. **AC2/AC3
  activated** (exactly 3/4). Evaluator `sonnet` **86/100**; both blocking findings reproduced and
  **FIXED in-PR**, not carried.

## Parked on Mark

**Nothing.** Zero open asks. Owed by the **shared driver** (frozen core — World cannot apply):
`ailang#611` (real per-role executor chain) and the World driver sync (missing `pi:*` pre-flight).

## Next

**`TR.B`** (capability snapshot + declared-effect confinement, AC5/6/7), then **`TR.C`**, the binding
gate — P6.B's prerequisite is satisfied only when `TR.C` is green. `SM.D` (item 8) is attended-only;
item 12 is small; items 13/14/15 (UI programme) were filed attended.

## Loop

- launchd, ~6h, headless. Issue **#53** (rotates Mondays 07:00 **local**).
- controller/planner `opus` · designer rotation (last `codex:gpt-5.6-sol`) · executor
  `codex:gpt-5.6-sol` · evaluator `sonnet`. `pi` **BARRED** from publish milestones.
- `derive-planner-lane.sh` absent → lane fails closed to opus, loudly. `metered=$0.00`; cap $5.

## Standing hazards

- **`rg` is NOT a binary here** — a harness-injected shell function, absent in CI. Use `grep`.
- **A refusal test asserting only *that* an error occurred pins no branch**: `DecodeRevision`
  re-encodes canonically, so every guard has a backstop. Pin the **measured message** per branch.
- **A rule-3j audit anchored to a DECISION LIST cannot contain branches the sprint itself writes** —
  3 shipped uncovered here. Enumerate the branches the diff ADDS, not the ones the doc froze.
- **A green `go test` is not a green `go vet`** — `copylocks` is outside both, invisible to CI.
- **`verify_ail.sh` never asserts the module count against 11** — only against 0. Item 12.
