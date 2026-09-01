# Mission Dashboard — Ailang World

*Snapshot, overwritten each iteration. History: `world-mission.md` STATUS + `world-mission-log.md`.*
**As of** 2026-09-01, iteration **146** · `dev` = [`cb73cab`](https://github.com/sunholo-data/ailang-world/commit/cb73cab) · CI **GREEN** (2/2) · **54 of 61** rows closed

## This iteration (REFUTATION — no code routed; the design was corrected and then correctly blocked)
- **Row 50 unparked by `D-WORLD-29`, revised to it, then RE-PARKED on the new `D-WORLD-31`.**
  Provenance of the attended ruling checked first-party: the commit that flipped the row is authored
  `mark@aitanalabs.com`, not the fleet account — a human answer, not a self-resolution.
- **Rule A does close row 50**, measured: the silent arm (a second, indented `KNOWN_BAD`) goes
  **rc=0 PASS=2 → rc=1 `count=2, want 1`** at both consumers, using only the pre-existing message.
- **And rule A introduces a new false green.** With the only `KNOWN_BAD` inside a never-executed
  `if false; then` / `fi`: base **rc=1 `count=0, want 1`** → rule A **rc=0 PASS=2**, while bash
  leaves the variable **`<UNSET>`**. `run.sh` has **17** control-flow openers, so that is 17 live
  positions, all loud before and silent after. Declared residual 6 + unit arm 11.
- **`verify_go.sh` flakiness pinned first-party:** same pristine tree, two runs — **rc=1 with 1 FAIL,
  then rc=0 with 0 FAIL**. The doc's hardcoded base set matched neither; AC7 now measures at drill
  start. Corroborates row 58 from the opposite direction.
- **Skill drift FIRED after 13 clean iterations.** The running skill is **147 lines behind**
  `origin/dev` (V1's checkout is 22 commits stale, 6 touching `SKILL.md`, plus an uncommitted
  in-place edit). Read the delta, followed origin's rules. **V1's to fix — World cannot.**

## Next picks (ready, ungated)
1. **Row 52** — CI step-scoping fix; direction ratified by `D-WORLD-30`, doc banked.
2. **Row 54** — the frozen-core drift gate compares the driver copy to *itself*.
3. **Row 55** — the row-47 lever parser false-reds on three forms of valid YAML.
4. **Rows 56–61** — canary fence; approvals spine; `verify_go.sh` flake; grep-cannot-prove-live;
   the row-48 needle's inert-rename false red; the one-inserted-line fail-open.
Then row **39**.

## Loop cadence + routing
Controller `claude:claude-opus-5`. Designer **`pi:ollama/deepseek-v4-flash:0731-cloud`** ran this
iteration (typed verdict `ok`, 112 s, non-empty diff) — **rotation advanced**, so
`claude:claude-fable-5` is next. Planner `opus`, executor `codex:gpt-5.6-sol` (probed rc=0, unused),
evaluator `sonnet`.

## Parked on Mark — ONE open ask
- **`D-WORLD-31`** — ship rule A as ratified and accept the declared residual (recommended), or hold
  row 50 for the declarative-fixture migration? `D-WORLD-29` is **not** reopened; this is a new fact
  measured after the ruling. Default if unanswered: row 50 stays parked, queue advances to 52.

## Quota / spend posture
`metered=$0.15101` of the $5 ceiling — one full-strength quorum round. Billing tripwire **CLEAN**;
nested `claude` calls go through `claude-sub`, which strips the API keys by construction.
