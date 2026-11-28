# World iteration 448 — readings (row 243 V4 S9)

All readings are first-party controller measurements unless marked otherwise. Paths under `~/.ailang/state/world-iter448/` are
rig-local scratch and are not published.

## Gate 0–2
- Kill switch absent (`mission-world.disabled`); `gh` = `sunholo-voight-kampff`; billing tripwire CLEAN.
- `mission_directives.sh --issue 635`: 0 allowlisted directives (34 comments; LEGACY watermark `2026-08-17T06:15:12Z`). The
  mission-scoped watermark `mission-world-last-seen` = `2026-10-27T09:03:42Z`; the older of the two was used, so nothing is missed.
- Step 6a `mission_gate0_self_notices.sh` rc 0: 0 crash notices (control 107:4 ok; `#574` watermark absent → DEGRADED single-watermark
  read). 3 transcript-scan notices in the window, fires `world-1795730668-41608` (2 hits), `world-1795769251-75944` (1 hit) and
  `world-1795786215-30397` (3 hits) — the same three iterations 446 and 447 reported; no new one.
- `mission-world` inbox: 0 unread. `ailang-world` and `world` inboxes hold one stale message each (2026-10-17, 2026-09-14), both
  previously triaged.
- Local `dev` = `origin/dev` = `9a40fe68`. No open own-account PRs. Dev CI on `9a40fe68`: `checks=3`, all `success`.
- Skill drift: the pin worktree matched origin on all skill files; the symlink-resolved copy (V1 main checkout) differed on
  `resources/gate-5-retro.md` only (origin's idle-fingerprint input list adds "open nightly-rig issues and their failure reports"). Gate 5
  read origin's version.
- Ledger valid, 109 rows, 1 OPEN (D-WORLD-122, Mark's Sol preset for E5 — not an S9 input). Clause map: 1.0 clauses 1–7 MET, none UNMET.
- Pick = record-447 Next, row 243 V4 S9. Design §8: S9 waits for "119 A/B and 120 A/B, E2, S8" — D-WORLD-119 and D-WORLD-120 RESOLVED (A),
  E2 read-level PASS (W33), S8 landed (`85ed2234`). No S9 plan, branch, PR or worktree existed. `git diff --stat 85ed2234 9a40fe68 -- cmd
  host docs app` empty. Quorum reader rc 0: V4 r2 `…T10-24-51Z` PROCEED 3/3, `absent_reviewers` `[]`.
- Verified at HEAD before routing: `runEpisodeFinishWith` at `cmd/ailang-worldd/episode_finish.go:63`, the F1 comment at `:92`.

## Routing
- Resolver: planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin` (`derive-planner-lane.sh` the same); executor
  `recipe codex:gpt-6.2-sol declared:provider-pin`; evaluator `agent-tool sonnet declared:alias-pin`; no designer (V4 reviewed).
- Codex probe rc 0. Codex usage 15.0% used / 25.9% allowed (10080 m window).
- Worktrees via `mission-worktree.sh` at `9a40fe68` (sprint, base) and at `c9cabdfd` (judge, branch `eval/world-iter448-r1`). A first
  judge-worktree call passed `--detach` where the helper takes a branch; it failed `worktree_add_rc_255` and left nothing (checked).

## Baseline (pristine `9a40fe68`, outside the sandbox, isolated `AILANG_STATE_DIR`, during the planner run)
store/authority/agentrun plain+race, verifygate and verbs: 0 fails. broker-plain 6 (`TestAilangRunLeavesNoCacheInTheWorktreeDelta`,
`TestExecHandlerWorktreeDelta`, `TestExecProtectedBaseline`, `TestExecProtectedCheckTimeout`, `TestExecProtectedIntact`,
`TestPolicyToolWorktreeDeltaOnlyForToolchainOps`); cmd-plain 13 (the 11 row-271 real-binary tests, `TestOperatorFilesReconcile`,
`TestOperatorStartCollisions`); daemon-plain 4; daemon-race 14. Base-extra (planner `gates.sh … outside-base`): checkrun, worktree and
surfacepage plain+race rc 0; whole `cmd -race` timed out (rc 124, the known can't-finish-locally leg).

## Planner (codex `gpt-6.2-sol`, ~96 min, rc 0, 819,311 tok)
- Outcome **B**: compiled whole 2,799 lines; S9 = M1 621 + M2 → 1,399 cumulative (cap 1,400); M3 = S9b, 1,400 incremental, compiled in
  Appendix P with its activation gates listed.
- Its own §7 block extracted and executed verbatim on a fresh export: `APPLY_ACCEPT_SNAPSHOT_OK` M1/M2/M3 fast, `FINAL_APPLY_RC 0`.
- Plan-reported (not first-party): named kills 21/12/12; caller census private entry and commit seam 0 → 0 → 1 (M3), claim 1 → 1 → 2.
- Controller round trip (`ctl_roundtrip.sh M1 M2`, fresh `git archive 9a40fe68`): `diff -r` rc 0 for cmd, host and docs at both stages.
- Destructive-write scan of `accept.sh`/`outside.sh`: every write under scratch or `t.TempDir()`.

## Outside-only legs on the planner's stage, before the executor (the iteration-447 process fix, first use)
- `outside.sh st-M2 M2` rc 0: `ACCEPT_OK M2 outside`, `outside-race` (`RACE_BUDGET cmd/ailang-worldd n=7 added=8.35s cap=15s`),
  `outside-1777`.
- `outside.sh st-M1 M1` rc 1: one fail, `TestOperatorFinishAuthInstall` in the TMPDIR-under-1777 leg (`start after finish auth failure 1:
  runner refused / why: could not bootstrap the companion`, 3.30 s). Re-runs: alone under 1777 3/3 PASS (0.27–0.29 s) on M1 and 3/3 on
  M2; alone plain 2/2; the whole M1 1777 leg 2/2 PASS (2.30 s, 2.42 s). Not reproduced; recorded as a flake candidate and handed to the
  judge, who attributed it (below). The executor was started on that basis.

## Executor (codex `gpt-6.2-sol`, ~9 min, rc 0, 133,203 tok)
§7 block verbatim (`TOOLS=$STATE/exec-tools`): `ACCEPT_OK`/`APPLY_ACCEPT_SNAPSHOT_OK` M1 and M2 fast; snapshots 9/0/9 and 15/0/15,
`shasum -c` OK, all files `cmp`-equal to the planner's stages; tools `cmp`-equal to the planner's; no `.orig`/`.rej`.

## Commits
- **M1 `9a04073e`** on `ACCEPT_OK M1 outside` (NUMSTAT 621).
- M2 attempt 1 REFUSED `INSTRUMENT FAILURE cmd-finish-stems truncated` (`TestEpisodeFinishDetachedPushFailedRefuses` running at the 240 s
  bound) while the controller's own whole `cmd -race` base leg ran concurrently: contention the controller caused. The controller stopped
  its own sweep driver (its PID only).
- M2 attempt 2, quiet rig: `ACCEPT_OK M2 outside`; `commit-guards.sh` REFUSED `INSTRUMENT FAILURE: truncated commit guard` — it runs whole
  `host/verifygate host/broker` with no `-timeout`, so broker hit Go's default 10 m (pristine broker 730 s here). A plan-instrument defect.
  Verifygate failed `TestMissionConfigDefaultFlowIntegration` (known load-sensitive) and `TestProtectedGitCandidateEquivalence`;
  `go list -deps ./host/verifygate` contains none of `host/worktree`, `host/checkrun`, `host/store` — S9 cannot reach it.
- Whole broker at the M2 stage with `-timeout 20m` (574.5 s): 5 fails, 3 outside the base 6
  (`TestAilangRunWorktreeDeltaOnTheFakeTool`, `TestExecNoProtectedIsByteIdentical`,
  `TestPublishDenialsPersistAndNeverDispatchWithLivePositiveControl`); all 3 PASS alone 2/2 at head AND 2/2 at base. Broker imports
  `host/store`; the set churns run to run (base 6 vs head 5, 2 shared). Load.
- **M2 `c9cabdfd`** committed on `ACCEPT_OK M2 outside` (the 445–447 precedent), NUMSTAT 1,399.

## Controller sweep on `c9cabdfd` (outside the sandbox)
gofmt, build, vet, linux vet, `go test -run '^$' ./...`, guard-store, guard-cli (`TestCLICredentialWriterSites`, `TestSessionRootBudgets`,
`TestOperatorProductionDeadlines`), guard-boundary, guard-daemon, `TestSubprocessSinkGate`,
`TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`, `callers.sh` (`CALLER_CENSUS_OK`; push/label/view-PR/create-PR/V3-put
1 each), `verify_ail.sh` (PASSED, 16 identities), checkrun/worktree race, store race 44.6 s (base 43 s): all rc 0.

## Judge (Agent `sonnet`, ~66 min, 120,903 tok): PASS 74 (25/18/16/6/9), 0 blocking
- Reachability: no production caller of the private entry or the local-commit primitive; push/PR/label unchanged. F1 now precedes flag
  parsing (only exact `--help`/`-h` forms exempt); a bad flag with no tty now gets the F1 refusal — stricter; no existing assertion weakened.
- 66 own mutants + 9 planner ids; the 8 planner mutants re-run all died. Non-blocking survivors: ~290 M2 lines with no test or caller (manifest
  file, store journal/review types, `NextOperatorFinish`, `OperatorCompletedLaunches`, `OperatorFinishLaunch`); the install pinning-skip
  scope (skipping `git` too, or every executable, survives); the private-entry `Kind=="finish"` and `ConfirmOnSurface` checks; a constant
  digest (untracked content hashing unasserted), walk/name/byte bounds, `config.worktree`/alternates fence, staged and race checks; commit
  branch check, post-commit recheck, update-ref CAS and author line.
- O1 (AuthInstall flake): not reproduced (25/25 under 32 CPU hogs); the existing `TestOperatorServiceStartup` failed 8/25 under the same
  load with the identical message. Cause: a 1 s real-time heartbeat-write budget inside `runnerTick` shared by every `start` test; the new
  test inherits it.
- O3 (store race): base 109.2 s vs head 40.3 s, once each — no S9 delta.
- **Controller reproduction of one survivor:** removing `frozen.Kind != "finish" || ` at `operator_finish_preflight.go:23` (applied, `cmp`
  confirmed) compiles and `TestOperatorFinishNoProof|TestOperatorFinishF1|TestOperatorFinishTTYParity` stay `ok` — survives, as reported.
  Restored; worktree clean.
