# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS), `-status-archive.md`, `-log.md`.*
**Iteration 106** · 2026-08-21 · `dev` @ `d1b7eae` · CI green (both jobs, SHA-addressed on the merge commit, run confirmed `event=push`)

> **A removal proves a check FIRES; only an ADDITION proves it LOOKS.** The frozen-surface gate killed
> the one mutation spelling the design doc names and was blind to three others — a pointer result, a
> method, and a type alias — each minting PROVEN from a raw hash with no seal, whole package green.

## In flight
- **Item 17 `w-validated-proven-evidence-boundary`** `[IN-SPRINT]` — `PE.A`–`PE.F` (4.70 d), **four landed**.
- **`PE.D` LANDED** — PR #79 → `d1b7eae`. The validator, the sealed mint authority (an unexported
  pointer to a per-instance non-zero-size allocation), the resolved grade, and four separately-pinned
  constructor refusals. `Resolve` runs mint-validity strictly before binding, because merged they
  would compare the zero-zero pair EQUAL and the forge would resolve.
- Judge `sonnet` **62/100 FAIL round 1** (two blocking, both real, both reproduced first-party) →
  repaired → **95/100 PASS round 2**, zero blocking, judge aimed at the repair.

## Next
- **`PE.E`** — real-store integration proofs, 0.85 d, test-only ~700 lines. No fake participates in
  any kill by construction. The plan flags it as the milestone most needing the out-of-sandbox re-run:
  its verdicts are wall-clock classified, not socket-bound, so a loaded sandbox can fake the mutant
  signature. Then `PE.F` last, forced by its own `EXACT_EVIDENCE_TESTS` pin.
- Row 14's predicate has been flipped for thirteen iterations (blocked only on item 18, complete since
  iter-93). It stays unpicked while item 17 is IN-SPRINT with an explicit NEXT — standing rule 1.

## Loop + routing
- Controller `opus` ×1 · executor `codex:gpt-5.6-sol` ×2 (probe + run) · judge `sonnet` ×2 (two rounds).
- **Fable and the designer rotation unspent a SEVENTH consecutive iteration** — no new doc needed.
- `metered=$0.00` of the $5 ceiling. No quorum purchased (in-sprint continuation), no GPU.
- Gates need BOTH `AILANG_BIN=/tmp/ailang-v0300/ailang` and `GOTOOLCHAIN=go1.25.6`; `verify_go.sh`
  fails closed without them, which is deliberate — a bare `go test` reports `ok` with the
  load-bearing assertions silently skipped.

## Parked on Mark
- **Nothing.** Decision ledger: 11 rows, **0 OPEN**. Zero open asks for the seventh iteration running.

## Quota posture
- Billing tripwire CLEAN (no `ANTHROPIC_API_KEY` / `ANTHROPIC_AUTH_TOKEN` in the loop's shells).
- Bookkeeping thread `#68` (23 comments, cap 80); weekly rotation not due.
