# World iteration 450 — readings (row 243 V4 S9c)

All readings are first-party controller measurements unless marked otherwise. Paths under `~/.ailang/state/world-iter450/` are
rig-local scratch and are not published.

## Gate 0–2
- Kill switch absent (`mission-world.disabled`); `gh` = `sunholo-voight-kampff`; billing tripwire CLEAN.
- `mission_directives.sh --issue 635`: 0 allowlisted directives (37 comments; LEGACY watermark `2026-08-17T06:15:12Z`).
- Step 6a rc 0: 0 crash notices (control 107:4 ok; `#574` watermark absent → DEGRADED single-watermark read); 4 transcript-scan
  notices in the window, the same four fires already reported by iterations 446–449; no new one.
- Local `dev` `9a40fe68` 2 behind `origin/dev` `91befdab`; state read from origin. Dev CI on `91befdab`: `checks=3`, all `success`.
- Skill drift: symlink-resolved copy differs from origin on `resources/gate-5-retro.md` only (main ailang checkout behind origin);
  the pin's copies match origin; Gate 5 read origin's version.
- Ledger valid, 109 rows, 1 OPEN (D-WORLD-122). Clause map: 1.0 clauses 1–7 MET. Pick = record-449 Next, row 243 V4 S9c.
- `git diff --stat d9635675 91befdab -- cmd host docs app` → 0 lines; the iteration-449 `M3.patch` dry-run rc 0, 0 offset/fuzz/reject.

## Routing
- Resolver against the V4 doc: planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`; executor `recipe codex:gpt-6.2-sol
  declared:provider-pin`; evaluator `agent-tool sonnet declared:alias-pin`; no designer. Codex probe rc 0.
- Tokens: planner 494,248; executor 57,489; judge 97,258 (r1), re-rulings 106,952 and 111,761 as reported.

## Baseline (pristine `91befdab`, outside the sandbox, isolated state)
- 0 fails: store, authority, agentrun (plain+race), verbs, broker, surfacepage (plain+race), checkrun-race, worktree-race.
- verifygate 1 (`TestMissionConfigDefaultFlowIntegration`); cmd-plain 11 (row-271 e2e); daemon-plain 3; daemon-race 5
  (adds `TestMCPPostBudgetProductionConstants`, `TestSeToolsRunBenchmarkTasksEndToEnd`). Pristine skip-11 cmd race: ok, 415.5 s.

## Plan, stage, executor, commit
- Plan outcome B: M1 activation 1,399 lines (+1,335/−64); M2 (S9d CLI mirror) 57 lines compiled, not applied.
- Controller round trip: `diff -r` rc 0 for cmd, host and docs at M1 and M2.
- Outside on the planner's stage: focused plain/race/1777 `ACCEPT_OK`; `RACE_BUDGET cmd/ailang-worldd min=12.47s cap=15s` (one reading),
  `host/daemon min=0.60s cap=12s`; Chrome `TestOperatorStartWalkthroughInChrome` PASS 2.76 s.
- Executor: `FINAL_APPLY_RC 0`; snapshot 21/0/21, `shasum -c` clean, 21/21 `cmp`-equal to the stage.
- M1 `e69298f4`: `ACCEPT_OK M1 outside` + `COMMIT_GUARDS_OK`, NUMSTAT 1,399, first attempt.

## Controller sweep on `e69298f4`
- `gates.sh … M1 outside` 2,824 s: every leg OK; compares cmd/daemon/verifygate/broker/surfacepage plain+race all OK;
  `GATE_RESULT issue=0 pending=0`; head cmd race 425.7 s (vs 415.5 s pristine; two distinct runs).
- Named mutants (overlay, tree untouched; run after the judge's verdict): 124 ids → 80 KILLED, 41 SURVIVED, 1 CONTROL_LIVE, 2 SKIPPED (S9d).
  First invocation of the plan's `named-mutants.sh` crashed (`args[@]: unbound variable`, bash 3.2); re-run as its exact `mutate.py` call.
- e12 reproduction: removing the runner `Unattributed` refusal (`operator_finish_runner.go:278`) leaves every `TestOperatorFinish*` /
  `TestOperatorCommit*` ok (21.0 s); restored byte-identical.

## Judge (Agent sonnet)
- r1 PASS 84 (26/24/17/9/8), 0 blocking; re-ruling 1 PASS 78 and re-ruling 2 PASS 77, both non-blocking, on the 41 survivors (code
  trace, execution UNMEASURED); precondition: S10 must not start until a test-only follow-up kills them.

## Rebase
- PR #716 on `e69298f4`: 0 runs; `mergeable` CONFLICTING/DIRTY (attended #711–#713 on `docs/QUICKSTART.md`). Rebased onto `4bac6fe1`;
  product diff byte-identical outside QUICKSTART; QUICKSTART added lines identical; NUMSTAT 1,399. Re-check of `f45da2a9`: build, vet,
  linux vet, test-compile, js-syntax, `ACCEPT_OK M1 outside`, QUICKSTART/verb/parity/app tests, verbs/appmanifest/surfacepage, verify_ail —
  all rc 0.

## Gate 3b
- PR #716 on `f45da2a9` (rebased): CI 3/3 (event `pull_request`) → MERGEABLE/CLEAN, `origin/dev` `4bac6fe1` unmoved → merged `--match-head-commit` → **`9b11c617`**. merge CI 3/3 at attempt 1 (run 46169801018; race `cmd/ailang-worldd` 528.9 s, `host/daemon` 532.0 s, `host/verifygate` 534.5 s, `host/broker` 422.3 s, all under the 900 s bound — broker and verifygate, untouched by S9c, rose by similar amounts, i.e. runner variance). **LANDED.**
