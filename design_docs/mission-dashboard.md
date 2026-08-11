# AILANG World — mission dashboard

*Snapshot, overwritten every iteration. History lives in `world-mission.md` (STATUS) and
`world-mission-log.md`. Last written: iteration 71, 2026-08-11.*

## Now

- **dev**: `93df1ec` — CI **green both jobs**, SHA-addressed, `checks=2` = expected 2,
  0-incident window (so the green is attributable).
- **Last landed**: **`TR.A1`** (PR #59) — item 11's first executable milestone. Store CAS
  (`CompareAndSetRegistryHead`), the TR-CJSON-1 canonical codec, descriptors/revisions and
  literal goldens. **AC1/AC4/AC10 activated** (zero-arms removed: exactly 3/2/1).
  Evaluator `sonnet` **94/100, zero blocking**; its one non-blocking finding was APPLIED.

## Parked on Mark

**Nothing.** Zero open asks. Owed by the **shared driver** (frozen core — World cannot apply):
`ailang#611` (real per-role executor chain) and the World driver sync (missing `pi:*` pre-flight).

## Next

**`TR.A2`** — `Reader`, eager `Snapshot`, cache, `BuildNext`/`Publish`, closing AC2/AC3
(the plan's T5–T8; scoped and measured, needs no re-design). Then `TR.B`, then `TR.C`.
`TR.C` is the binding gate; P6.B's prerequisite is satisfied only when it is green.
Item 8's `SM.D` is **attended-only** and never routes headless.

## Loop

- Cadence: launchd, every ~6h, headless. Bookkeeping issue **#53** (rotates Mondays 07:00 local).
- Routing: controller/planner `opus` · designer rotation (last `codex:gpt-5.6-sol`)
  · executor `codex:gpt-5.6-sol` · evaluator `sonnet`. `pi` is **BARRED** from publish milestones.
- `derive-planner-lane.sh` is absent here: the lane fails closed to opus loudly every fire.
- Spend: `metered=$0.00` this iteration (all quota buckets); ceiling $5/iteration.

## Standing hazards

- Export `AILANG_BIN=/tmp/ailang-v0300/ailang` **and** `GOTOOLCHAIN=go1.25.6` before any gate.
- **`rg` is NOT a binary here** — a harness-injected shell function, absent under `env -i` and in
  CI. Never put it in a committed command; the repo uses `grep` throughout.
- **A refusal test that asserts only *that* an error occurred pins no branch.** `DecodeRevision`
  re-encodes canonically, so every guard has a backstop that refuses the same input — two guards
  were provably unpinned until each case pinned its own message (iter-71).
- **`verify_ail.sh` never compares the module count against 11** — only against 0 (`:233`); the
  total is printed (`:243`), not asserted. Follow-up item 12.
