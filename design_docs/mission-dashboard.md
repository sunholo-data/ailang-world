# AILANG World — mission dashboard

*Snapshot, overwritten every iteration. History: charter STATUS + `world-mission-log.md`.*
**As of** 2026-08-14 (iteration 84) · **dev** `323baf6` · **CI** green on the base, SHA-addressed.

## Just designed — and PARKED
**Item 17 `w-validated-proven-evidence-boundary`** — 566 lines, `323baf6`, designer
`codex:gpt-5.6-sol`, **28 verification rows**, decomposed **3.5d / 3.0d / 2.0d** (the row said
1.5–2d). Two quorum rounds, all four slots **present**: R1 both reject → designer adopted the
reviewer's own direction fix → R2 both reject on **new** grounds. `metered=$0.179422`.

## Parked on Mark — THREE one-word asks (14 and 15 unanswered since 2026-08-13)

1. **Item 17 — how does `ValidateProof` earn authority?** A content-addressed report is not an
   *authenticated* one. **A** = re-execute the pinned checker inside the validator, stored reports
   become non-authoritative cache. **B** = MAC/sign reports with a host-held key. **C** = ship as
   designed, record the forgery route as a declared limitation.
2. **Item 15, §7.3 freeze timing.** **A** = freeze the v1 `DecisionPacket` now. **B** = ratify only
   the `TimeoutPolicy` set. *The rejecting reviewer's `proposed_fix` IS option B.*
3. **Item 14.** **A** = context-aware store reads + request-scoped deadline. **B** = defer behind 18.

**Next**: item **18** (`w-daemon-read-cancellation`) — the only unparked, unblocked row left;
item 5 `P6.B` still blocked (`#498` Lane A landed, but item 5 needs a *public* seam).

## Loop

launchd, ~6h watchdog. Controller `opus` · designer `codex:gpt-5.6-sol` (probe rc=0; ran twice —
initial + the one prescribed revision) · planner/executor/evaluator did not fire. Cap $5.
**FLAGGED**: generator≠judge breaks when the rotation lands on codex — `gpt-5.6-sol` designed this
doc while `gpt5-6-sol` reviews it.

## Carry-forward

**A reviewer's fix can be runnable and still be a DOWNGRADE.** Iter-83 checked a remedy is
*implementable*; `gemini-3-1-pro`'s premise was right (`verify.verified` is an int) and its fix
(use a count) discards what `verify.results[].function` already supplies. Ask whether a fix
**reduces what the gate can observe**, not only whether it runs.
**The kernel's four "exports" are MODULES** — publishing `world/types` publishes every `Evidence`
constructor *and* `gradeOf`: a foreign `.ail` module minted `PROVEN` from a made-up digest (check
rc=0, test PASSING, control `IMP010` fired); same attack on the fix: `expected 4, got 1`.
