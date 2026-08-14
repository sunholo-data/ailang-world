# AILANG World — mission dashboard

*Snapshot, overwritten every iteration. History: charter STATUS + `world-mission-log.md`.*
**As of** 2026-08-14 (iter-84 + attended ratification) · **dev** `4787542` · **CI** green on the base.

## ALL THREE ASKS ANSWERED — ratified attended 2026-08-14, recorded in the queue rows

1. **Item 17 — how `ValidateProof` earns authority: OPTION B.** Authenticate reports with a host-held
   MAC/signing key; hash recomputation stays the integrity check it is, and the MAC supplies the
   *provenance* it never could. **A** rejected on cost — a compiler+solver on the critical path of
   every grade resolution, on a daemon whose store layer has no bounded-wait discipline yet. **C**
   rejected as a resting state. Key custody is bounded: single-host, no trust boundary crossed,
   worst-case loss = regenerate reports = exactly A's steady state.
2. **Item 15, §7.3: OPTION A.** Freeze the v1 `DecisionPacket` now; amendments become `/v2`, never
   in-place edits. B's merit needs item 7, parked behind item 5's upstream blocker with no ETA, and
   R1's defect (laws blind to `deadlineAt`) is fixed and proven, so B's stated ground is gone.
3. **Item 14: OPTION B.** Defer behind item 18. The defect is pre-existing on all seven GET routes —
   B *removes* it rather than extending it, satisfying the reviewer's catch, at nil ordering cost.

## What each item needs now

- **18 `w-daemon-read-cancellation`** — `[NEXT]`, needs a design doc. Now blocks item 14. **15** —
  unparked, **sprint-ready**, gated on nothing. **14** — unparked, blocked on 18; carries an
  independent finding, `Internal` branches pass `err.Error()` verbatim to an unauthenticated client.
- **17** — unparked, needs a *revision round*, not a new design: (i) the MAC seam; (ii) the **V27
  repair** — re-point §3.3/§3.4 at `verify.results[]` (per-identity `function`+`status`), never the
  `verify.verified` integer nor the reviewer's weaker `verifiedCount`; (iii) the negative control —
  hand-authored perfect report bytes must yield `unauthenticated_report`, not a seal.
- **5** still blocked (`#498` Lane A landed; item 5 needs a *public* seam, upstream has no `pkg/`).

## Loop · carry-forward

launchd, ~6h watchdog. Controller `opus` · rotation at `codex:gpt-5.6-sol`. Cap $5; iter-84 metered
`$0.179422`. **FLAGGED**: generator≠judge breaks when the rotation lands on codex — `gpt-5.6-sol`
designed item 17's doc while `gpt5-6-sol` reviewed it. Same collision at iter-82.

**Hash recomputation proves integrity, never provenance.** Every field a content-addressed report
asserts is a public value its author chose; the hash matches *because they authored it*. Authority
needs a secret or a re-execution — the choice between them is cost, not strength. (V27, still live:
a runnable reviewer fix can still be a DOWNGRADE — ask what it costs the gate's *observability*.)
