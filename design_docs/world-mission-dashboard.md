# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History: `world-mission.md` STATUS · `world-mission-status-archive.md` · `world-mission-log.md`.*

**Iteration 98 · 2026-08-20 · LANDED**

## Latest landing
- **Row 21 `w-archive-stderr-in-manifest`** — PR [#73](https://github.com/sunholo-data/ailang-world/pull/73) → squash `9fa2647`; Gate 3b green on the **merge** commit (`present=2 == expected=2`). Evaluator `sonnet` **97/100**, zero blocking.
- stderr no longer reaches the archive manifest: bounded stream-separated probes at both `CombinedOutput` sites + a **conditional** self-heal (zero process execs on a healthy artifact). 5 mutations, all killed; M4a/M4b a dual pair.

## Next picks (all unblocked, headless-routable)
1. **Row 20** `w-capsule-output-cap-load-flake` — re-scoped: measure the **margin**, not the outcome.
2. **Row 22** `w-daemon-lock-wait-not-deadline-bound` — a lock-contended read returns via `busy_timeout` (2 s), not the read deadline. Two independent reviewers cite it.
3. **Row 23** `w-store-deadline-free-residue-owner` — 11 deadline-free store reads have a ratchet and no owner.

## Parked on Mark
- **`D-WORLD-22`** (item 17) — **one word**: does tranche 1 absorb row 22's lock-wait bound (**A**), or does the claim weaken to what is proven (**B**)? Open since iteration 96; the **only** open decision (ledger valid, 9 rows).

## Routing
- Designer rotation seeded `claude:claude-fable-5`; codex quota-limited, gemini read-only under `CapRemoteSandbox` → rotation WRAPPING. **Unspent this iteration** (no new doc) for the first time in four.
- Planner `opus` (`derive-planner-lane.sh` → `opus fail-closed:env-pin`) · Executor `opus` · Evaluator `sonnet` (generator ≠ judge).
- `pi:deepseek` lane **SUSPENDED** by `D-WORLD-20`; chain is codex → opus.

## Cost / gates
- Metered **$0.00** of $5 (no quorum ran — doc arrived quorum-clean). Buckets: `opus` ×2, `sonnet` ×1. Billing tripwire **CLEAN**.
- `dev` green at `9fa2647` — both CI jobs `success`.
- Pinned: `ailang` v0.30.0 at `/tmp/ailang-v0300/ailang`; `GOTOOLCHAIN=go1.25.6` (ambient go1.26.4 is DENIED and reds the store canary).
- Bookkeeping issue **#68** (rotates Mondays 07:00 local; not due — created 2026-08-17T19:19:43Z).
