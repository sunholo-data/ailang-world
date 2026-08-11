# AILANG World — mission dashboard

*Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS) + `world-mission-log.md`.
Last written: iteration 71, 2026-08-11.*

## Now

- **dev**: `93df1ec` — CI **green both jobs**, SHA-addressed, `checks=2` = expected 2,
  0-incident window (so the green is attributable).
- **Last landed**: **`TR.A1`** (PR #59) — item 11's first executable milestone: store CAS, the
  TR-CJSON-1 canonical codec, descriptors/revisions, literal goldens. **AC1/AC4/AC10 activated**
  (zero-arms removed: exactly 3/2/1). Evaluator `sonnet` **94/100, zero blocking**; its one
  non-blocking finding APPLIED.

## Parked on Mark

**Nothing.** Zero open asks. Owed by the **shared driver** (frozen core — World cannot apply):
`ailang#611` (real per-role executor chain) and the World driver sync (missing `pi:*` pre-flight).

## Next

**`TR.A2`** — `Reader`, eager `Snapshot`, cache, `BuildNext`/`Publish`, closing AC2/AC3 (the plan's
T5–T8; measured and scoped, needs no re-design). Then `TR.B`, then `TR.C`, the binding gate —
P6.B's prerequisite is satisfied only when it is green. Item 8's `SM.D` is attended-only and never
routes headless. Item 12 (`verify_ail.sh` module pin) is small; may fold into a `TR.*` docs task.

## Loop

- launchd, ~6h, headless. Issue **#53** (rotates Mondays 07:00 **local**).
- controller/planner `opus` · designer rotation (last `codex:gpt-5.6-sol`) · executor
  `codex:gpt-5.6-sol` · evaluator `sonnet`. `pi` **BARRED** from publish milestones.
- `derive-planner-lane.sh` absent → lane fails closed to opus, loudly. `metered=$0.00`; cap $5.

## Standing hazards

- Export `AILANG_BIN=/tmp/ailang-v0300/ailang` **and** `GOTOOLCHAIN=go1.25.6` before any gate.
- **`rg` is NOT a binary here** — a harness-injected shell function, absent in CI. Use `grep`.
- **A refusal test asserting only *that* an error occurred pins no branch**: `DecodeRevision`
  re-encodes canonically, so every guard has a backstop (2 unpinned until each pinned its message).
- **`verify_ail.sh` never asserts the module count against 11** — only against 0. Item 12.
