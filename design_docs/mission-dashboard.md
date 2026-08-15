# AILANG World — mission dashboard

*Snapshot, overwritten every iteration. History: charter STATUS + `world-mission-log.md`.*
**As of** 2026-08-15 (iter-87) · **dev** `bef0153` · **CI** green, SHA-addressed, `checks=2` both
`success`. Doc revision only; no code changed. `metered=$0.266188` of $5.

## ⚠ THE HEADLESS QUEUE IS FULLY BLOCKED — the two asks are the only unblockers

All 19 rows (18 + `4b`) complete or blocked. **1,2,3,4,4b,9,10,11,12,13,15,16** LANDED · **5**
upstream `ailang#498` still **OPEN** (re-checked `2026-08-04`) + unbuilt commit-boundary contract ·
**6** until 2–5 · **7** until 5's `P6.B` · **8** attended-only · **14** behind 18 · **17**/**18** on
one-word asks.

## Item 17 `w-validated-proven-evidence-boundary` — revised twice, re-PARKED

Designer `codex:gpt-5.6-sol` ×2, 566 → 711 lines; quorum rounds **3 and 4**, `absent_reviewers: []`
both. Round 4 is the **first reviewer flip to `pass`** in this item's history (`gemini-3-1-pro`);
`gpt5-6-sol` rejects. All three prescribed deliverables landed — MAC seam per Mark's Option B, the
V27 repair (`verify.results[].function`, never the bare int), and the `unauthenticated_report`
negative control.

**The spine: the reviewer's own verbatim alternative was the trap.** Round 3's `gpt5-6-sol` catch
was real — AC14 demanded a *daemon-owned* key at *first startup* in a tranche whose §8.2 excludes
`host/daemon`/`cmd/**`: unsatisfiable by construction. Its fix had two arms; the designer took arm 2
**verbatim** (the carve-out's whole safeguard) — and round 4 rejected the result, correctly: with no
production root the key is caller-supplied, so `NewValidator(key [32]byte, …)` (`:198`/`:211`) plus
the **free** `GradeOfValidated` (`:201`) lets any Go caller mint `ResolvedGradeProven`. **Arm 2 and
"authority boundary" are incompatible, and only applying arm 2 revealed it.**

**Measured, not forwarded (rule 3f):** gemini's premise objection — 4 record shapes, one `ai-check`,
pinned v0.30.0. Scalar record, `list[scalar]` and bare-ADT *param* **verify**; record with a **bare
ADT field** or `list[ADT]` **error**. Broader than the doc claimed, and `rc=0`/`check.passed=true`
throughout — **silent to the exit code**. Rows V31/V32; gemini passed. **Also corrected:** §3.6 said
`EXACT_TOTAL_VERIFIED=5`/`TESTS=20`; at HEAD **10**/**39** (`aaada20`, one iteration earlier).

## Second finding — `grep` cannot see `~~`

"What's next" came out wrong twice: item 16 (COMPLETE), then item 8 — by reading
`HEADLESS-ROUTABLE` out of **struck-through** prior-head text while the live head says
attended-only. ~19 rows keep dead heads inline *on purpose*. Process fix recorded; watch-item at 1.

## Loop

launchd, ~6h watchdog. Controller `opus` · designer `codex:gpt-5.6-sol` (rotation `claude` →
`codex`, written back). **TWO open asks** (17 and 18), both one-word.
