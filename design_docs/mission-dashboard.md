# AILANG World — mission dashboard

*Snapshot, overwritten every iteration. History lives in `world-mission.md` (STATUS) and
`world-mission-log.md`. Last written: iteration 69, 2026-08-11.*

## Now

- **dev**: `32b086c` — CI **green both jobs**, SHA-addressed, `checks=2` = expected 2,
  0-incident window (so the green is attributable).
- **Last landed**: `VL.B` (PR #57) — Z3 installed in **both** CI jobs; `host/verifygate`'s
  accept-arms now assert `verify gate PASSED`. Evaluator `sonnet` **91/100, zero blocking**.
- **Item 9 `w-verify-binary-lockfile`**: COMPLETE. `9/OD-11` **ratified and discharged**.

## Parked on Mark

**Nothing.** Zero open asks. `9/OD-10` and `9/OD-11` are both ratified and closed.

Owed by the **shared driver** (frozen core — World cannot apply): `ailang#611` (real
per-role executor chain) and the World driver sync (missing `pi:*` pre-flight loop).

## Next

Item 5 `w-mcp-projection`'s single remaining prerequisite — the **transition registry**,
still ABSENT at HEAD. Either write it, or re-scope `P6.B` around its absence.
Item 8's `SM.D` is **attended-only** and never routes headless.

## Loop

- Cadence: launchd, every ~6h, headless. Bookkeeping issue **#53** (rotates Mondays 07:00 local).
- Routing: controller/planner `opus` · executor chain `codex:gpt-5.6-sol → pi:deepseek-v4-flash → opus`
  · evaluator `sonnet` (generator≠judge). `pi` is **BARRED** from publish-capable milestones.
- `derive-planner-lane.sh` is absent here: the lane fails closed to opus loudly every fire.
- Spend: `metered=$0.00` this iteration; budget ceiling $5/iteration. All lanes quota buckets.

## Standing hazards

- Export `AILANG_BIN=/tmp/ailang-v0300/ailang` **and** `GOTOOLCHAIN=go1.25.6` before any gate —
  `verify_go.sh` is rc=1 at BASE without the toolchain pin.
- The released `ailang` shells out to z3 via **hardcoded absolute paths**; PATH cannot hide it.
  The only faithful solverless control is `AILANG_Z3_PATH=<nonexistent>`.
