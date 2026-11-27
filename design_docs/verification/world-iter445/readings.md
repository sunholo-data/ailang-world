# Iteration 445 — controller readings (row 243 V4 S7, inert operator requests)

Fire `world-1795786215-30397`, attempt 1. Controller `claude-opus-5-5` (session). Base `origin/dev` `a83d40a6` (= S6 merge `00b4ec55` + #694 `2b3ba7af` + record #696).

## Gates 0–2
- Gate 0:
  - Kill switch absent (`mission-world.disabled`); `gh` = `sunholo-voight-kampff`; billing CLEAN.
  - `mission_directives.sh`: 0 allowlisted directives on #635 (28 comments; LEGACY watermark file).
  - Step 6a rc 0: 0 crash notices (control 107:4). 2 transcript-scan notices in window: fire `world-1795730668-41608` (2 hits) and fire `world-1795769251-75944` (1 hit). Key names are only in the private controlplane notice.
- Gate 1:
  - Local `dev` `2b3ba7af` was behind `origin/dev` `a83d40a6`; read from origin. Base `a83d40a6…@2026-11-27T13:33:15Z`.
  - Skill drift: the resolved symlink is byte-identical to ailang `origin/dev` on 14/14 files.
  - Dev CI: `a83d40a6` in progress, parent `2b3ba7af` success.
- Gate 2:
  - Clause map: 1.0 clauses 1–7 MET, none UNMET. Ledger `--check` valid (108 rows); 1 OPEN at pick (D-WORLD-121), resolved by attended groom #697 (`216c4c42`) during the iteration.
  - V4 quorum: `mc_quorum_gate2` rc 0 (`w-world-address-bar-verbs-v4-2026-11-18T10-13-56Z.json`).
  - No S7 commit, branch, PR or worktree existed. `ConfirmOperatorAction` / `RequestOperatorAction` had 0 production callers (grep, excluding `_test.go`; the `OperatorAction` read census control hit).

## Routing (resolver followed verbatim)
- planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`; executor `recipe codex:gpt-6.2-sol declared:provider-pin`; evaluator `agent-tool sonnet declared:alias-pin`. Codex probe rc 0.
- Tokens: planner 995,656; executor 73,587 + 177,288 (r2); judge sonnet 150,245 (r1), 173,014 cumulative (r2); controller not reported.

## Baseline (pristine `a83d40a6`, outside the sandbox, concurrent with the planner)
- store-plain/race, authority-plain/race, worktree-plain, verbs-plain: 0 fails.
- verifygate-plain 1 (`TestMissionConfigDefaultFlowIntegration`); broker-plain 4; cmd-plain 17; daemon-race 16.

## Commit gates
- M1 `80af7f36`: `accept.sh M1 fast` `ACCEPT_OK`, NUMSTAT 173.
- M2 attempt 1: `TestOperatorAdmissionHTTP` `--- FAIL (11.02s)`, `surface_operator_test.go:332: context deadline exceeded`. Alone `-count=10`: ok (6.0 s). `GOMAXPROCS=2 -race -count=5` beside whole `host/store`: 4 PASS, 1 FAIL at `:357`.
- M2 attempt 2: guard wrapper new fails `TestKnownPositiveDelegates`, `TestReleaseChangeNotice` (verifygate), `TestExecProtectedCheckTimeout`, `TestExecProtectedDetectsCreation` (broker).
  - verifygate pair alone: FAIL at head AND base without `AILANG_STATE_DIR`; `ok` (229.7 s) at head with an isolated state/cache dir.
  - `TestExecProtectedCheckTimeout` alone: FAIL at head and at base. `TestExecProtectedDetectsCreation` alone `-count=2` at head: ok.
  - M2 `e27c93ef` committed on `accept.sh M2 outside` (`ACCEPT_OK`), NUMSTAT 1,399.

## Gate sweep on `e27c93ef`
- gofmt/build/vet/linux-vet/test-compile/verify_ail/callers/accept outside-race rc 0. New cmd tests 1777 and `-count=3` rc 0.
- `daemon-new-repeat` rc 1: `TestSurfaceVerbOutcomes` (×3). At base `ok` (0.99 s); at head FAIL (`:60` two rows, `:95` golden).
- store-plain/race, verbs, approval+surfacepage, authority, verifygate: 0 new. broker: 1 new (`TestExecProtectedCheckTimeout`, base-red alone).

## Round 2 `eb762a2e` (226 lines) — controller gates outside the sandbox
- gofmt/build/vet/linux-vet/test-compile/callers/accept outside + outside-race rc 0.
- Daemon selector `-count=3` and `-race -count=2` rc 0. F1 contention shape: 5 PASS, 0 FAIL. cmd new tests `-count=3` and 1777 rc 0. store-plain 0 new. surfacepage+verbs rc 0.
- Whole `host/daemon` plain (692 s): 7 fails, 2 not in the base set — `TestA2AInvokeDeadlineClosesSocket`, `TestRegistryPhaseHasItsOwnClock` — both `-count=3` alone `ok` at head (18.96 s) and base (18.75 s).
- Cumulative `git diff -M --numstat a83d40a6 eb762a2e -- cmd host docs app` = 1,559.

## Gate 3b
- PR #700 head `eb762a2e`: CI 3/3 success at attempt 1 (run 46044386173). Race: `cmd/ailang-worldd` 503.351 s, `host/daemon` 513.203 s, `host/verifygate` 486.634 s, `host/broker` 404.641 s.
- `mergeable` MERGEABLE/CLEAN re-read on that head; merged `--match-head-commit` → `6836c9a9`.
- Merge CI: 3/3 at attempt 1 (run 46046770928; race `cmd/ailang-worldd` 284.9 s, `host/daemon` 363.2 s, `host/verifygate` 324.9 s, `host/broker` 279.0 s). LANDED
