# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS), `-status-archive.md`, `-log.md`.*

**Iteration 103** · 2026-08-21 · `dev` @ `cbd17de` · CI green (both jobs, SHA-addressed, run confirmed)

> **This fire died at Gate 4 and its own retry landed the record.** The 02:24 fire merged `PE.A` and
> wrote the record into the main checkout, then died at 03:34 on `API Error: Connection lost
> mid-response` before `git add`; the driver re-fired at 03:35. Same iteration, finished by its retry
> — **not** an iteration 104. Every inherited claim re-derived first; one REFUTED (worktrees clean).

## In flight
- **Item 17 `w-validated-proven-evidence-boundary`** — `[IN-SPRINT]`. Six CI-green milestones
  `PE.A`–`PE.F` (4.70 d; 27/27 live mutation rows mapped; zero human asks).
- **`PE.A` LANDED** — PR #76 → `cbd17de`, Gate 3b green on the merge, judge `sonnet` 96/100 zero
  blocking. Kernel `ProofReceipt` arm, projection, golden, gate pins.
- **NEXT: `PE.B`** (0.83 d) — bounded one-snapshot store read + cached `BusyTimeout()`; also carries
  the DR-2 ratchet fix (`TestNoNewDeadlineFreeStoreReads` cannot see `ReadObject`). Then `PE.C` →
  `PE.D` → `PE.E` → `PE.F` in compile order; `PE.F` must be last.

## Queue after item 17
Rows **22**, **23** unblocked and headless-routable · **24**–**27** designed-pending. Row **27**
`w-interface-hash-does-not-cover-the-interface`: `interfaceHash` hashes only `ailang.toml`, so it does
not move when an exported ADT gains a constructor — found at iteration 81, re-broken by item 17's §8.3.

## Parked on Mark
**NONE.** Decision ledger: **11 rows, 0 OPEN** (`scripts/mission_decisions.sh --check`).

## Cross-mission
`mission-v1` accepted the "a Total is a claim about a column" proposal **in principle**. Its own finding
measured here first-party: **auto-close exposure is ZERO** (0 hits over 286 commit messages / 70 PR
records, both controls firing). **A zero is not a guard** — it is a habit, and this repo's one open
issue is `#68`, the loop's own human channel.

## Loop / routing / cost
Controller `claude:claude-opus-5` · planner `opus` (`fail-closed:env-pin`) · executor
`codex:gpt-5.6-sol` → `opus` (deepseek SUSPENDED, `D-WORLD-20`) · judge `sonnet` (generator≠judge).
Designer rotation and Fable both unspent (4th consecutive iteration). `metered=$0.00` of $5 · quota
`opus` ×2, `codex` ×2, `sonnet` ×1 · billing tripwire CLEAN.

## Gates
`./scripts/verify_ail.sh` · `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh` (the export is mandatory — it fails closed without it). Pinned `AILANG v0.30.0`.
