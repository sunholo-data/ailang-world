# AILANG World — mission dashboard

*Snapshot, overwritten every Gate 4. History lives in `world-mission.md` STATUS + `world-mission-log.md`.*

**As of** 2026-08-05, iteration 50 · dev @ `de80792` + this iteration's doc commit · CI green

## In flight

- **Item 4f `w-bench-load-confound` — COMPLETE.** `BC.A′` `0b72019` (iter-48) · `BC.B′` `d357474`
  (iter-49) · controller pass **`C2b`** (iter-50) discharged the last five mutations. Census
  **16 of 16**, re-derived by command. Doc → `design_docs/implemented/`.
- **`[NEXT]`: item 8 `w-self-mod-vertical`** (clause-7) — unparked; its gate ("until 4 lands")
  cleared at iter-35 and Mark re-scoped it attended on 2026-08-04. **Needs a design doc** →
  design-doc-creator. Its VERIFY-FIRST clause is binding at pick: live-repro `ailang publish`
  auth + vendor-registration mechanics against the pinned binary before writing milestones.
- **Item 5 `w-mcp-projection` — still BLOCKED, but on ONE prerequisite, not three.** `#498` seam
  cleared prereq 1; 4b's receipt law cleared prereq 3; the **transition registry is still absent**
  (measured at HEAD with a firing known-positive control). Broker `Session` API does exist.

## Latest

- Iter-50 headline: `MUT-PAIR-INLINE-BUILD`'s **rule** fired (R1 both sections + R2's orphan
  cascade) while its stated **secondary observable did not** — legs moved 7→8 s because a warm Go
  build cache prices the whole compile at 1–2 s. *A secondary observable a cache can erase is
  worse than none: it gives a reviewer a reason to stop looking.* Struck in the doc.
- Also: the "12 mutations" figure was a transcription (really 13 for `BC.B′`, 16 item-wide), and
  iter-49's P3/P4 citation repair had been written only into the row that found it.

## Loop

- Cadence: launchd, `mission-world`. Controller `claude-opus-5`.
- Routing: executor `codex:gpt-5.6-sol` · evaluator `sonnet` · planner `opus` · designer rotation
  pointer `claude:claude-fable-5` (not fired since iter-47). Iter-50 fired **none** of them — an
  evidence-only pass, and one the sandbox cannot run (loopback binds).
- Verify profile `ailang-code`; AILANG pinned **v0.30.0** at `/tmp/ailang-v0300/ailang`; Go pinned
  `GOTOOLCHAIN=go1.25.6` (except AC6's fixture session, which requires `auto`).
- Bookkeeping issue **#32** (week of 2026-08-03).

## Cost

- Iteration 50 `metered=$0.00` — every role on a quota bucket. Budget $5/iteration.

## Parked on Mark

**Nothing.** Zero open asks. Next free OD number: **`OD-9`**.
