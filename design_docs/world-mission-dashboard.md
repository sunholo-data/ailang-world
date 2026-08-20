# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History: `world-mission.md` STATUS · `world-mission-status-archive.md` · `world-mission-log.md`.*

**Iteration 99 · 2026-08-20 · PARKED (no sprint ran)**

## This iteration
- **Row 20 `w-capsule-output-cap-load-flake`** — DESIGNED, then **parked** `needs-human-review` on the new **`D-WORLD-23`**. Doc `591c16d` + revision `40b7f19`; designer `codex:gpt-5.6-sol`; 2 quorum rounds, **both reviewers present in both**, `absent_reviewers` empty.
- **The row's own filed fix is refuted by measurement.** `verifyExecutable` reads+sha256s a **91.8 MB** binary at `capsule.go:134` — **before** `context.WithTimeout` at `:152`, so bounded by **nothing**; 65–78% of idle `Run`, ×6.8 under load. The child exec `ExecTimeout` bounds is **~15 ms against 5 s** (~330× margin). "Raise the wall clock" enlarges the largest margin.
- **New row 24 `w-host-subprocess-cleanup-boundary`** — the OVERFLOW kill is not process-group-wide while the CANCELLATION kill is, in **both** `host/capsule` and `host/broker`. Bound still holds (**3.005 s** vs a 3 s cap with a grandchild on stdout; control **11.29 ms**), so the cost is promptness, not termination.

## Next picks
1. **Row 23** `w-store-deadline-free-residue-owner` — needs **no new design doc**; cleanest headless route.
2. **Row 24** `w-host-subprocess-cleanup-boundary` — designed-pending, gated on nothing.
3. **Row 22** — stays behind `D-WORLD-22`. **Row 20** unparks on `D-WORLD-23`.

## Parked on Mark — TWO open, and one answer settles both
- **`D-WORLD-23`** (NEW, one word): when a quorum fix would fold **separately-owned** work into the tranche in front of it, does the tranche **always keep scope and weaken its claim** (**A**), or must the controller escalate each time (**B**)? **Arm A also resolves `D-WORLD-22`.**
- **`D-WORLD-22`** (item 17, since iter-96): absorb row 22's lock-wait bound (**A**) or weaken the claim (**B**)?
- Carve-out now foreclosed on the **scope axis 7 times**; ledger valid, 10 rows.

## Routing
- Designer rotation **advanced** `claude:claude-fable-5` → `codex:gpt-5.6-sol`; **Fable lane unspent a 2nd consecutive iteration**.
- Planner/executor `opus` · Evaluator `sonnet` (generator ≠ judge) — **none ran this iteration**.
- `pi:deepseek` lane **SUSPENDED** by `D-WORLD-20`; chain is codex → opus.

## Cost / gates
- Metered **$0.108951** of $5 (quorum r1 `$0.048867` + r2 `$0.060084`). Codex rode the **quota** lane. Billing tripwire **CLEAN**.
- `dev` green at `47e12cc` — both CI jobs `success`, run confirmed to exist.
- Pinned: `ailang` v0.30.0 at `/tmp/ailang-v0300/ailang`; `GOTOOLCHAIN=go1.25.6` (ambient go1.26.4 is DENIED). **`AILANG_BIN` unset makes `host/capsule` tests `t.Skip` silently — a bare `go test` there is false-green.**
- Bookkeeping issue **#68** (rotates Mondays 07:00 local; not due).
