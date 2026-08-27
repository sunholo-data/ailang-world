# Mission Dashboard — Ailang World

*Snapshot, overwritten each iteration. History lives in `world-mission.md` STATUS + `world-mission-log.md`.*

**Iteration 132 — 2026-08-27 — `P43` LANDED.**

## Latest
- PR [#99](https://github.com/sunholo-data/ailang-world/pull/99) → squash [`ecfb62d`](https://github.com/sunholo-data/ailang-world/commit/ecfb62d).
- Gate 3b **GREEN on the merge commit**: `present=2 == expected=2` (enumerated from `ci.yml`'s
  own job list; `ci.yml` is the only workflow), `notdone=0`, both `success`, `runs=1 event=push`,
  parent control `checks=2`, `mergeable` read first.
- Row 43 CLOSED. The floor-raise coupling map now lives in `scripts/verify_ail.sh`'s head block
  and `design_docs/coding-standards.md` §S8, bound by
  `TestFloorRaiseInventoryNamesEveryCoupledFile`.

## In flight / next
- Queue: rows **44**, **45**, **46**, **47**, **48**, **49**, **50**, **51** (new), then **39**.
- Row 45's gate is DISCHARGED — it was "gated on row 41 landing", and row 41 landed at iteration 129.

## Loop cadence + routing
- Designer: rotation `claude:claude-fable-5` → `pi:ollama/kimi-k3:cloud`. This iteration the pi
  entry **probed rc=0 and then failed the real run** on the provider's own
  `503 model 'kimi-k3' is temporarily overloaded`; fell to Fable, FLAGGED. Pointer unchanged.
- Planner `opus` (lane derived `opus fail-closed:env-pin`) · Executor `codex:gpt-5.6-sol`
  (ChatGPT subscription, `$0`) · Evaluator `sonnet`.
- Evaluator **failed round 1 at 62/100**, passed round 2 at **90/100** after a structural fix.

## Cost
- `metered=$0.2730` of the `$5` ceiling (two full-strength quorum rounds). All other lanes are
  quota buckets.

## Parked on Mark
- **Nothing.** Ledger: 13 rows, **0 OPEN**.
