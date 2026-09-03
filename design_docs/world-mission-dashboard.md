# Mission Dashboard — AILANG World

_Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS), `world-mission-log.md`._

**As of:** 2026-09-03 · iter 153 · `dev` = [`12b8c87`](https://github.com/sunholo-data/ailang-world/commit/12b8c87) · CI green (3/3) on the merge commit · local verify gate green both legs

## Last iteration
**Row 58 LANDED** · **HARNESS** · PR [#115](https://github.com/sunholo-data/ailang-world/pull/115) → [`12b8c87`](https://github.com/sunholo-data/ailang-world/commit/12b8c87) · metered **$0.00** (controller-authored; no sub-agent spawned).

**The probe timeout is a CONCURRENCY defect, not a speed one — and row 58's headline was still wrong after one amendment.** The archived interpreter's `--version` probe costs **47–52 ms** warm and **1211–1294 ms** on a first exec at a fresh path. macOS serializes that assessment **globally**: 8 concurrent first execs return at 1289…**8871 ms**, 12 at 1255…**13691 ms**, linear at ~1.13 s each. `probeTimeout` is a **per-probe** bound, so it is crossed at **N ≥ 9** no matter how fast one probe is. That is the whole flake — green package-by-package, flaky under `go test ./...`, never red on Linux CI.

**Two candidate causes refuted by the same measurements:** CPU contention (iter-141's axis) moves it only to 1322/1373/1385 ms under 16 spinners; the Observatory cleanup over a 553 MB DB is paid by the ~50 ms warm arm. Negative control: 8 concurrent **warm** execs max at **348 ms**, so the linearity belongs to *first* exec.

**Disposition — a deliberate non-change.** `probeTimeout` was **not** raised: the cost scales with a dimension the probe cannot observe, so a constant only moves the cliff. The deliverable is attribution — `archive.EnvironmentFailure` / `AttributeFailure` label the deadline an ENVIRONMENT failure at 8 call sites. Nothing skipped, nothing suppressed: the test still FAILS and now names which of the two things went wrong. **6 mutants, 6 RED**; M6 (drop `probeVersion`'s `%w`) is killed by the shipped-path arm alone.

**Second defect, fixed here:** `host/archive` resolved the pin from a hardcoded `/tmp/ailang-v0300/ailang` — dead on **every** machine since iter-151 moved it — so `TestArchivePinnedInterpreter` was a silent SKIP with no red anywhere. Now reads `AILANG_BIN` like its sibling packages. SKIP → PASS, and it runs in CI for the first time.

## Goal distance
**Goal unmoved** (no product surface changed; loop-machinery work). Row census remains **carried, not measured** — row 72 tracks that. Row 57 tracking-only upstream; row 50 parked on `D-WORLD-31`.

## Next picks
1. **Row 59** `w-static-grep-cannot-prove-an-assertion-is-live` — an AC proved "load-bearing" by `grep -c`, which cannot tell a live assertion from a compiled-and-unreached one.
2. **Row 76** (NEW) `w-verify-go-driver-drift-gate-short-circuits-the-entire-local-go-gate` — `verify_go.sh` is rc=1 in 0.99 s at base on a **fleet-owned** drift red that fatals at `:224`, before `go build` at `:443`. Correct red, wrong ordering, no opt-out — so the gate every AC here names cannot answer "is this tree's Go code green?".
3. **Row 74** `w-the-personal-email-gate-...-does-not-exist` — build the gate the rulebook already claims exists. Then 60–66, 68–73, 75, then 39.

## Routing / cadence
Controller `claude:claude-opus-5`. Designer rotation pointer unchanged at `claude:claude-fable-5` (no designer ran). Verify profile `ailang-code`; pin **v0.30.0** at `~/.pinned-ailang/ailang` (PATH's `ailang` is `-dirty` and is never used for gates).

## Parked on Mark
**`D-WORLD-31`** — ONE WORD: ship rule A as ratified, or hold row 50 for the fixture migration? Re-asked unchanged; **no new ask this iteration**.

## Quota posture
metered **$0.00** of $5 this iteration. Billing tripwire CLEAN.
