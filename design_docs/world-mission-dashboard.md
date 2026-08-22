# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS), `-status-archive.md`, `-log.md`.*
**Iteration 107** · 2026-08-22 · `dev` @ `a87c723` · CI green (both jobs, SHA-addressed on the merge commit, 0 not-green) — after the record commit turned dev RED in this iteration's own test and was fixed forward

> **A stress control varies the axis you thought of; the false red lives on the one you didn't.**
> I ran AC16's timing 23/23 green across CPU contention and called it sound. The judge varied
> *parallelism* — `GOMAXPROCS=1` reds 10/10 on unmutated code, wearing the exact mutant signature it
> exists to detect. More runs of the same shape would never have found it.

## In flight
- **Item 17 `w-validated-proven-evidence-boundary`** `[IN-SPRINT]` — `PE.A`–`PE.F` (4.70 d), **five landed**.
- **`PE.E` LANDED** — PR #80 → `daf48a6`. Four real-store integration proofs, test-only, no new
  exported production symbol; no fake participates in any kill.
- Judge `sonnet` **66/100 FAIL round 1** (all four findings reproduced first-party, all four REAL) →
  **85/100 PASS round 2**, zero blocking → round-3 text-only commit taking two of five non-blocking.
- Three of the milestone's own checks could not fail; deleting the wrong one exposed a second defect it
  had absorbed. Re-drilled module-wide: M4/M22/M23/M24/M26/M30 all red, named arm in every set.

## Next
- **`PE.F`**, the last milestone: the persistent named-manifest gate in `verify_go.sh`
  (`REQUIRED_EVIDENCE_TESTS` + an exact count, terminal `Action=pass` only, anti-vacuity floor), its
  self-mutation test in `host/verifygate`, AC12's zero-diff assertion, and the full 27-row re-drill.
  Pinned last **without exception** — its count gate reds on any test landed after it, and PE.E added
  four, so it pins the OBSERVED count, never one transcribed from the plan.
- `PE.F` must also correct the plan: M26/M30 are listed as real-store-only kills; PE.D's fakes kill
  both in isolation. Row 14's predicate stays flipped-but-unpicked while 17 is IN-SPRINT (rule 1).

## Loop + routing
- Controller `opus` ×1 · executor `codex:gpt-5.6-sol` ×2 · judge `sonnet` ×2 rounds, in its **own**
  worktree so its mutation drills could not race the controller's gate runs.
- **Fable and the designer rotation unspent an EIGHTH consecutive iteration.** `metered=$0.00` of $5.
  No quorum purchased (in-sprint continuation), no GPU.
- Gates need BOTH `AILANG_BIN=/tmp/ailang-v0300/ailang` and `GOTOOLCHAIN=go1.25.6` — `verify_go.sh`
  fails closed without them by design, runs ~150 s, and there is no `timeout` binary on this rig.

## Parked on Mark · quota posture
- **Nothing parked.** Decision ledger: 11 rows, **0 OPEN** — zero open asks for the eighth iteration.
- Billing tripwire CLEAN. Thread `#68` (24 comments, cap 80); rotation not due (created after the
  Monday-07:00 **local** boundary).
