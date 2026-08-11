# AILANG World — mission dashboard

*Snapshot, overwritten every iteration. History lives in `world-mission.md` (STATUS) and
`world-mission-log.md`. Last written: iteration 70, 2026-08-11.*

## Now

- **dev**: `11fb1fd` — CI **green both jobs**, SHA-addressed, `checks=2` = expected 2,
  0-incident window (so the green is attributable).
- **Last landed**: the **`w-transition-registry` design doc** (PR #58) — item 5's single
  remaining blocker now has a quorum-reviewed design: 3 milestones, 28 verification rows,
  11 acceptance criteria, 44 named mutations, 3.5 World days. Item 5's blocker is now
  **designed, not absent**.

## Parked on Mark

**Nothing.** Zero open asks. Owed by the **shared driver** (frozen core — World cannot apply):
`ailang#611` (real per-role executor chain) and the World driver sync (missing `pi:*` pre-flight).

## Next

**Sprint-plan `w-transition-registry`** (`TR.A` → `TR.B` → `TR.C`), then execute `TR.A`.
`TR.C` is the binding gate; P6.B's prerequisite is satisfied only when it is green.
Item 8's `SM.D` is **attended-only** and never routes headless.

## Loop

- Cadence: launchd, every ~6h, headless. Bookkeeping issue **#53** (rotates Mondays 07:00 local).
- Routing: controller/planner `opus` · designer **rotation** (last used `codex:gpt-5.6-sol`)
  · executor `codex:gpt-5.6-sol` · evaluator `sonnet`. `pi` is **BARRED** from publish milestones.
- `derive-planner-lane.sh` is absent here: the lane fails closed to opus loudly every fire.
- Spend: `metered=$0.21259` this iteration (3 reviewer calls); ceiling $5/iteration.

## Standing hazards

- Export `AILANG_BIN=/tmp/ailang-v0300/ailang` **and** `GOTOOLCHAIN=go1.25.6` before any gate.
- **`rg` is NOT a binary here** — it is a shell function injected by the agent harness, absent
  under `env -i` and in CI. Never put it in a committed command; the repo uses `grep` throughout.
- A quorum reviewer can go **ABSENT ON BUDGET** and the verdict still reads `proceed` at N−1.
  Re-run the absent reviewer at a raised cap before banking that pass — iter-70 did, and it
  flipped to REJECT.
