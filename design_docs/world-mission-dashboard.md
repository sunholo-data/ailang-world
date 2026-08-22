# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS), `-status-archive.md`, `-log.md`.*
**Iteration 108** · 2026-08-22 · `dev` @ `189299b` · CI green (both jobs, SHA-addressed on the merge commit, `present=2 == expected=2`, 0 not-green)

> **A removal proves a check FIRES; only an addition proves it LOOKS.**
> Every mutation rule this loop owns is removal-shaped, so a gate can pass every drill in the
> rulebook and still be blind to the case it exists for. PE.F's manifest gate survived four
> removal-direction mutants; the arm that mattered was appending one real *passing* test to
> `host/evidence` — count 37→38, RED.

## Landed this iteration
- **Item 17 `w-validated-proven-evidence-boundary` is COMPLETE — all six milestones in.**
  `PE.F`: PR #82 → squash `189299b`.
- **AC8**: focused `host/evidence` named-manifest leg in `verify_go.sh` before the broad legs —
  terminal `Action=pass` top-level identities, **set equality**, `EXACT_EVIDENCE_TESTS=37` pinned
  from the shell, three anti-vacuity floors failing LOUDLY. Its isolated self-mutation gate in
  `host/verifygate` executes the **live** script, not a copy of its logic.
- **AC9**: the full **27-row** re-drill — every row with sha256 LANDED, `go build` rc=0 *before* any
  test result, a red set enumerated by **running** it, sole-killer vs one-of-N, `cp` restore.
- **Twelve §6 divergences recorded, not patched.** A refuted premise corrected: M26/M30 do **not**
  need the real store (V56).
- Judge `sonnet` **96/100 PASS**, zero blocking; two of its three non-blocking findings reproduced
  and **fixed in a round-3 commit** — a floor with no arm, and a record citing a deleted `.snap/` path.

## Next
- **Item 22 `w-daemon-lock-wait-not-deadline-bound`** (clause-2) is the queue head. Re-verify its
  predicate at pick time — a declared blocker is a claim. Row 14 becomes eligible now 17 is closed.

## Loop + routing
- Controller `opus` ×1 · executor `codex:gpt-5.6-sol` ×3 rounds (gate · 6 drill rows · the other 21)
  · judge `sonnet` ×1 in its **own** worktree.
- Round B measured its own bottleneck (`host/verifygate` ≈27 s/row); round C narrowed the
  enumeration **with a control proving verifygate green under a Go mutant** and tripled the pace.
- **Fable and the designer rotation unspent a NINTH consecutive iteration.** `metered=$0.00` of $5.
  No quorum (in-sprint continuation), no GPU.
- Gates need BOTH `AILANG_BIN=/tmp/ailang-v0300/ailang` and `GOTOOLCHAIN=go1.25.6`; `verify_go.sh`
  runs ~150 s and there is no `timeout` binary on this rig.

## Parked on Mark · quota posture
- **Nothing parked.** Decision ledger: 11 rows, **0 OPEN** — a ninth iteration with zero open asks.
- Billing tripwire CLEAN. Thread `#68` (25 comments, cap 80); rotation not due (next Monday-07:00
  **local** boundary is 2026-08-24).
