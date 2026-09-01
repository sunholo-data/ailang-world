# Mission Dashboard — Ailang World

*Snapshot, overwritten each iteration. History: `world-mission.md` STATUS + `world-mission-log.md`.*
**As of** 2026-09-01, iteration **142** · `dev` = [`7077e45`](https://github.com/sunholo-data/ailang-world/commit/7077e45) · CI **GREEN** (2/2)

## This iteration (REFUTATION — nothing shipped, a premise died)
- **Row 52 PARKED** on new decision **`D-WORLD-30`**. The row called the CI step-scoping bug
  "measured non-exploitable". It is exploitable **both** ways, measured first-party: a **false
  positive** blaming an unrelated step (`rc=1 at ci.yml:166`), and a **fail-open** where a `- name:`
  decoy inside the step's own `run:` block leaves a live `continue-on-error` on the guarded step
  with the gate at `rc=0 --- PASS`.
- Two full-strength quorum rounds, both BLOCKED, killed two successive locators — each refuted by
  running it, not by arguing. Doc banked at
  `design_docs/planned/w-wiring-test-step-scoping-imprecise-under-key-reorder.md` (627 lines).

## Next picks (ready, ungated)
1. **Row 53** quorum reviewer silenced by its own review content — mostly upstream (`ailang#941`).
2. **Rows 54–57** — evaluator/controller findings filed at iters 133–139.
3. **Row 58** — `verify_go.sh` is **flaky on this rig**; never make its rc=0 an acceptance criterion.
4. **Row 59** — a `grep -c` cannot prove an assertion is live; discharge with a mutation.
Then row **39**.

## Loop cadence + routing
Controller `claude:claude-opus-5`. Designer rotation `claude:claude-fable-5` ⇄
`pi:ollama/deepseek-v4-flash:0731-cloud` — last used **fable**, so deepseek is next. Planner `opus`,
executor `codex:gpt-5.6-sol`, evaluator `sonnet`. Fable diet (one DOC = authoring + at most one
protocol-mandated revision) met, not exceeded.

## Parked on Mark — THREE open asks
- **`D-WORLD-28`** — how should `verify_go.sh` guarantee its nested race-control module can execute?
- **`D-WORLD-29`** — should a single *indented* shell assignment be ACCEPTED or REJECTED?
- **`D-WORLD-30`** *(new)* — row-52 fix: **LINE SCAN** or **YAML PARSE**? Recommendation **A**
  (line scan, hardened): no new dependency, and measured to catch the attack that killed both
  drafts. **B** adds the second direct dependency to a `go.mod` that has exactly one.
Rows **48**, **50**, **52** are the items these block.

## Quota / spend posture
`metered=$0.2524` of the $5 ceiling — both quorum rounds, nothing else. Every other lane is a
subscription quota bucket. Billing tripwire **CLEAN**; nested `claude` calls go through
`claude-sub`, which strips the API keys by construction.
