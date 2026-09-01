# Mission Dashboard — Ailang World

*Snapshot, overwritten each iteration. History: `world-mission.md` STATUS + `world-mission-log.md`.*
**As of** 2026-09-02, iteration **147** · `dev` = [`a1744d3`](https://github.com/sunholo-data/ailang-world/commit/a1744d3) · CI **GREEN** (2/2 on the merge commit) · **55 of 63** rows closed

## This iteration (PRODUCT — row 52 LANDED via PR [#110](https://github.com/sunholo-data/ailang-world/pull/110))
- **The row-44 wiring test's step locator was wrong in BOTH directions, and both were live at HEAD.**
  ARM A false positive: it blamed an *unrelated* step's `continue-on-error` flag (`rc=1 @ ci.yml:164`).
  ARM B fail-open: `rc=0 --- PASS` over a live forbidden flag on the guarded step. Re-derived
  first-party, not inherited.
- **Direction was human-ratified (`D-WORLD-30`) and hardened at quorum round 3.** The locator now
  anchors on the SHALLOWEST enclosing `steps:`, pins `expectedStepCol = 6` loudly, and refuses on
  containment and identity. Round 3 blocked at full strength; all three objections were MEASURED
  rather than forwarded and closed under the narrow-refinement carve-out.
- **Counterfactual, the number that justifies the sprint:** the judge reverted only the locator hunk
  and re-ran the arms — **7 of 14 are newly load-bearing** (both live defects, a third false
  positive, a wholly new identity invariant, and three re-indent arms the old scan absorbed at rc=0).
- **Two of this iteration's own claims were refuted by its own lanes.** The planner measured that the
  `expectedStepCol` pin I added at round 3 is *unreachable* on the arm the doc named (AC8 could not
  fail) and replaced it; the judge then found two false-positive shapes nobody had declared.

## Loop / routing
controller `claude:claude-opus-5` · designer `claude:claude-fable-5` (rotation; advanced) ·
planner `opus` (`opus fail-closed:env-pin`) · executor `codex:gpt-5.6-sol` · evaluator `sonnet`
(**93/100 PASS, zero blocking**; generator≠judge three ways). `metered=$0.18372` of $5.

## Next (banked, all gated on nothing)
`54` driver copy stale · `55` dispatch-lever parser false-reds · `56` canary fence blind to a
skipped canary · `57` approvals spine green-under-the-row · `58` verify_go.sh flaky at base ·
`59` static grep cannot prove an assertion live · `60` P1 needle reds on an inert rename ·
`61` P1 gate fails open on one inserted line · **`62`/`63` new, from this iteration's judge** ·
then `39` session authority.

## Parked on Mark
- **`D-WORLD-31` (row 50)** — one word. Ship the ratified rule A as-is (residual declared and
  pinned), or hold row 50 for the declarative-fixture migration `gpt5-6-sol` asked for?
  Default if unanswered: row 50 stays parked, the queue advances.

## Quota posture
Anthropic available (`MISSION_ANTHROPIC_AVAILABLE=1`); billing tripwire CLEAN; Fable diet spent on
ONE design doc (authoring run only — the round-3 fixes were a controller carve-out, not a designer run).
