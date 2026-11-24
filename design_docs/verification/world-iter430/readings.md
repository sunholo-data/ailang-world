# Iteration 430 — readings (row 243 V4 S4b)

All times UTC, 2026-11-24, rig Voights-Mac-Studio.local.

## Gate 0–2
- Gate 0: kill switch absent; `gh` = sunholo-voight-kampff; billing CLEAN; 0 allowlisted directives on #635 (7 comments); step 6a rc 0 (0 crash, 0 transcript notices; `-prev` #574 watermark absent → single-watermark read).
- Gate 1: local `dev` `af7a610c` behind `origin/dev` `c42fe698` (iteration 429's record merge) — state read from origin. Skill drift: 0 of 14 files (resolved symlink target vs the driver clone's origin/dev).
- dev CI at `c42fe698`: `go race gate` FAILURE (run 45464352750 attempt 1): `panic: test timed out after 8m0s`, `FAIL host/daemon 480.105s`, running `TestWorkspacePackageCacheDigestIsStamped`; the same run's `host/verifygate` took 468.7 s (typical ~320 s). `gh run rerun --failed` → attempt 2 **success**. Docs-only tree; its code parent `ee8ef1c0` was 3/3 green.
- host/daemon race package time on the last 9 dev runs (CI job logs): 345.2, 341.6, 343.0, 342.0, **480.1 FAIL**, 341.9, 385.0, 342.8, 328.0 s. Failing runs coincide with a slow runner (verifygate 464.9 / 468.7 s vs ~320 s).
- Third recorded instance of the host/daemon race 8-minute package timeout: (1) dev `212a87ba` run 45425824922 (running `TestWorkspacePackageCacheRefusesUnprovisionedRegistry`); (2) PR #651 head `4fc443cd` run 45453319415 attempt 1 (`FAIL host/daemon 480.125s`, same test); (3) dev `c42fe698` run 45464352750 attempt 1. Filed as row 267 under D-WORLD-109 (2).
- Ledger valid at 107 rows, 0 OPEN. Quorum for `w-world-address-bar-verbs-v4`: r2 `…2026-11-18T10-24-51Z` PROCEED, `absent_reviewers` = [], verdicts pass/pass/pass (claude-sonnet-5@claude-p, gemini-4-argon, oc-glm-5-3). Design unchanged since `8dc91744`.
- Died-mid-flight traces: 0 PRs, 0 branches, 0 worktrees naming iter430 or S4b; `git log origin/dev --grep S4b` shows only S4's M4 and plan.

## Gate 3
- Resolver: planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin` (derive-planner-lane same); executor `recipe codex:gpt-6.2-sol declared:provider-pin`; evaluator `agent-tool sonnet declared:alias-pin`. Codex probe from the sprint worktree rc 0.
- Planner 19:21:54 → 20:02:48 (~41 min), 381,380 tok. Plan `a8bf0858`: M1 638 / M2 878 / M3 978 (`git diff -M --numstat c42fe698 HEAD -- cmd host docs app`); 72/72 named mutant kills; 55 whole-package rows UNMEASURED IN SANDBOX.
- Controller round trip (`ctl_roundtrip.sh M1 M2 M3` on a fresh `git archive c42fe698`): every `diff -r` of cmd/host/docs vs `st-M*` rc 0.
- Executor 20:03:52 → 20:06:08, 56,458 tok: M1–M3 `ACCEPT_OK` (PASS 44 / 85 / 112, 0 FAIL, 0 SKIP); snapshots byte-equal to the planner's stages; `commit_ms.sh` (acceptance re-run outside the sandbox) → `dcd65dff` 638, `1afe1ad9` 878, `40c21340` 978; committed tree = snapshot.
- False-green (5) pre-check: every new `os.WriteFile` in the diff targets `t.TempDir()`; production adds `<creds>.operator-episode-<ep>.lock` sidecars beside the served credentials store (documented in plan §1.3/F5).
- Gates on `40c21340` (plan §6, outside the sandbox, 1,243 s): every leg rc 0 except `verify_go` (rc 1, 542 s): `host/daemon` `TestSeToolsRunBenchmarkTasksEndToEnd` ("the mock saw [POST …, GET ] want exactly the one api_call_json POST") and `host/verifygate` `TestMissionConfigDefaultFlowIntegration/valid_ledger`. 0 changed lines in either package. Alone at head: daemon test 1 FAIL (during judge load) then 3/3 PASS; verifygate 1/1 PASS. At base `c42fe698`: daemon test 2/2 PASS. `pkg-cmd-ailang-worldd` 235 s, `race-cmd-ailang-worldd` 348 s. Observation of real LaunchAgents / `~/.ailang-world` / launchctl count: identical before and after.
- Evaluator Agent `sonnet`, foreground, own worktree `.wt-world-iter430-row243-v4s4b-eval`, ~46 min, 140,846 tok: **PASS 85**, 0 blocking. Report `evaluator-r1.md`; ledger 70 lines.
- Controller reproduction of judge finding 1: `grep -n terminalEpisodeLock` hits episode.go:364, episode_preset.go:282, episode_resume.go:175, episode_launch.go:289/456, episode_finish.go:128 (positive control) and nothing in `approve.go` / `session_new.go`, while `provisionEpisode(` is called at approve.go:398 and session_new.go:649.
