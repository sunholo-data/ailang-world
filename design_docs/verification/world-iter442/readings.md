# Iteration 442 — controller readings (row 243 V4 S5f4, activation)

Fire `world-1795730668-41608`, attempt 1. Controller `claude-opus-5-5` (session). Base `origin/dev` `ad1a482c` (= S5f3 merge `a50b0d72` + record #688).

## Gates 0–2
- Gate 0: kill switch absent; `gh` = `sunholo-voight-kampff`; billing CLEAN; `mission_directives.sh` on #635: 0 allowlisted directives (22 comments; LEGACY watermark file); step 6a: 0 crash and 0 transcript notices on #635/#574 (control 107:4); `mission-world` inbox 0 unread.
- Gate 1: local `dev` was 11 commits behind `origin/dev`, so the charter and log were read from origin. Base `ad1a482c…@2026-11-26T22:07:18Z`. Skill drift: the resolved symlink and the pin copy were each byte-identical to ailang `origin/dev` on all 14 files. Dev CI at `ad1a482c` was in progress at Gate 1 and finished **3/3 green at attempt 1** (run 45901932809), so D-WORLD-121's default B did not trigger row 267.
- Gate 2: clause map 1.0 clauses 1–7 MET. Ledger 1 OPEN: D-WORLD-121, default B. Pick = row 243 V4 S5f4, record 441's Next. No S5f4 plan, branch, PR or worktree existed. Controller census at HEAD: four start-path DECLARATIONS, zero callers. `M4.patch` (iteration 441's prototype) dry-run on a fresh export of HEAD: every hunk applied except 1 of 1 in `host/store/operator_start_recovery_test.go`.

## Routing
- Resolver, followed verbatim: planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`, executor `recipe codex:gpt-6.2-sol declared:provider-pin`, evaluator `agent-tool sonnet declared:alias-pin`. Designer not spawned (reviewed design r1, quorum r2 PROCEED 3/3).
- The operator's standing request asked for every role through the Agent tool "exactly as the routing table specifies". The table returned `recipe` for the planner and executor, so the Agent tool is not a valid path for them, and they ran as `codex exec` recipes. The REQUIRED judge ran on the Agent tool: round 1, a fresh round-2 instance, and that instance resumed by name for a round-3 addendum. No role failed to spawn.
- Codex probe rc 0.

## Baseline (controller, pristine `ad1a482c`, outside the sandbox)
- `base/` (`-timeout 20m`, co-loaded with the planner): `store-plain` 1 (`TestConvertedStoreMethodsM5b/cutoff`, deadline under load), `store-race` 0, `authority-plain` 0, `authority-race` 0, `verifygate-plain` 1 (`TestMissionConfigDefaultFlowIntegration`, load); `cmd-plain` and `daemon-race` TIMED OUT at 20 m (instrument failure, discarded).
- `base2/cmd-plain` (`-timeout 45m`): complete, 696 s, 13 fails = the 11 row-271 real-binary tests + `TestOperatorFilePublicationAdmission` + `TestOperatorFilesReconcile` (git-fixture load).
- `base/worktree-{plain,race}`: 0 fails (recorded for round 2, which touched `host/worktree`).

## Planner
- `codex exec`, workspace-write, 23:21–00:45 (~84 min), **609,798 tok**, rc 0. Plan + JSON committed `48b4a9ec`.
- **Outcome A**: one milestone, the whole activation unit. The rebased prototype compiled at 917; M1 measured **1,391** (`-- cmd host docs app`; unfiltered product sum also 1,391). §1.4 premises found false: the old patch's one rejected hunk; the prototype read a catalogue generation the live loop never fills (M1 re-reads the durable heartbeat).
- Ledgers: 136 named diff-anchored kills + 1 equivalent (`D-OWNED-LIMIT`); all 5 carried judge-441 ids killed (`J-FROM-BINARY`, `J-BUDGET-EXPIRED`, `J-RUNNER-STOP-CANON`, `O-DROP-LOADED` in `TestOperatorActivationDrainRefusals`; `J-REVOKE-MISSINGSTAT` in `TestOperatorRevokeStatRefusal`).
- Race: the five new CLI tests measured 12.73 s at the lowest sample against the 12 s budget, with samples ranging to 58 s under load; daemon +0.
- Controller round trip `ctl_roundtrip.sh M1`: `diff -r` rc 0 on `cmd`, `host`, `docs`.

## Executor (round 1)
- `codex exec` (~6 min, **70,401 tok**), rc 0. M1 applied cleanly; `accept.sh M1 fast` `ACCEPT_OK` (125 PASS, 0 FAIL); snapshot 14/0/14, 14/14 cmp-equal to the planner's `st-M1`.
- `commit_ms.sh M1` attempt 1 REFUSED: the NEW `TestOperatorFirstManagedEpisodeProcess` failed at 24.99 s (`production tick failed: 0 operator start blocked`) in a run where existing git-fixture tests also timed out at 30 s. Alone, `-count=3`: 5.30 / 3.34 / 6.08 s PASS. Retry committed **M1 `1602bc5b`** (NUMSTAT 1,391; `VERIFYGATE_FAILSET_OK`). The judge then made this the blocking finding.

## Controller gates on `1602bc5b`
- gofmt, build, vet, Linux vet, test compile, `verify_ail`, caller census (four call sites in `operator_preparation.go`, one entry in `operator_runner.go`): rc 0.
- `host/store` plain/race 0 fails; `host/verbs` plain/race 0; whole `cmd/ailang-worldd` plain 11 fails, **0 new** vs `base2`.
- `accept.sh M1 race`: one red, existing `TestOperatorBorrowedWorktreeGitProcess/during_git` (git subprocess under load).

## Judge round 1
- Agent `sonnet`, foreground, ~65 min, **91,556 tok**: **FAIL 68** (25/13/16/8/6).
- **B1:** `TestOperatorFirstManagedEpisodeProcess` ran the real loop under a real `budgetContext(22 * time.Second)`. `GOMAXPROCS=1 -race -count=3` with 32 busy loops: PASS 27.99 s, PASS 22.36 s, FAIL 43.70 s with the controller's exact message. Controller had reproduced it independently (attempt 1 above).
- Survivors: J08 (post-prompt `before != after`), J09 (companion-unloaded check); J04, J10, J19 (five new tests only); J18 killed only by existing load-sensitive worktree tests. 11 sampled planner ids reproduced KILLED.
- No double-effect hole across the seven windows; lock order holds; nothing of S6 reachable.

## Round 2
- Codex (~10 min, **231,842 tok**), 250 lines over 15 files: named constant `operatorPreparationBudget` (22 s); `serviceEnv.operatorBudget` (zero in production) overrides the tick, its child budgets and the per-git cap via a context key and a new exported `host/worktree.WithCommandTimeout`; pins J04/J19 `TestOperatorStartMintRecoveryRemedies`, J08/J09 `TestOperatorAttendedStartDisposition`, J10 `TestOperatorFirstManagedEpisodeProcess`, J18 `TestOperatorPreparationRechecksHeartbeatAfterWorktree`.
- `commit_ms.sh R2` → **`c1f848de`** (NUMSTAT 1,601; `ACCEPT_OK`; `VERIFYGATE_FAILSET_OK`).
- Controller round-2 gates: every leg rc 0 except `accept-race` (existing `TestOperatorTerminalRevokeExclusion/durable_owner`, 4 s mint-fixture deadline under load; alone PASS 4.7 s per the judge); whole `cmd/ailang-worldd` 12 fails, **0 new**; `host/worktree` plain/race 0.
- **Judge round 2** (fresh Agent `sonnet`, ~37 min, **74,419 tok**): **PASS 85** (28/25/18/8/6), 0 blocking. B1 CLOSED (contended repro: PASS 79.24 / 76.59 / 28.88 s). J04/J08/J09/J10/J18/J19 all KILLED by the named tests. Non-blocking: the production budgets are pinned only by constant equality (its J21 tick 22 s → 1 s and J22 +1 h on the zero-override path both survive); `WithCommandTimeout` is an unclamped exported seam with no production caller; the override sets every child budget to one value; 1,601 > 1,400 cap.

## Round 3 (CI-found)
- PR #689 CI on `c1f848de`: both Go jobs red on the five new activation tests, each < 1 s, on ubuntu-24.04: `private session directory: path writable by others` (and `:175 … operator action conflicts with existing work`). Cause: `runnerPath` walks ancestors to the injected home and refuses `perm & 0022`; the fixture's sessions dir and fake agent sat outside the runner home, so the walk reached `/tmp` (1777). macOS's private temp root hid it.
- Controller reproduction on macOS: `TMPDIR=/private/tmp/w442tmp` (parent 1777) → both messages; without the override → PASS.
- Codex (test-only, 16 lines; `newWorktreeStartRigWithContext` puts the agent dir and sessions under `s.config.Home`) → **`f250e9f8`** (NUMSTAT 1,617). Commit gate: the seven activation tests PASS with the 1777-parent `TMPDIR`.

## Gate 3b
- PR #689 on `f250e9f8`: 3/3 at attempt 1 (race `cmd/ailang-worldd` 283.984 s, `host/daemon` 358.156 s, `host/verifygate` 325.666 s). `MERGEABLE CLEAN`; merged `--match-head-commit` → `c3504c56`.
- Merge `c3504c56` CI: **3/3 at attempt 1** (run 45931571443; race `cmd/ailang-worldd` 374.780 s, `host/daemon` 421.681 s — a slower runner, under the 480 s bound). **LANDED.**
- Round-3 executor tokens: not captured. Judge round-3 addendum: 83,213 tok (cumulative for the resumed instance).
