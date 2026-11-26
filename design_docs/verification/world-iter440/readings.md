# Iteration 440 — controller readings (row 243 V4 S5f2)

Fire `world-1795694137-68847`, attempt 1. Controller `claude-opus-5-5` (session). Base `origin/dev` `b5f41458` (= S5f merge `886648ab` + attended #680 `0033ca11` + record #682).

## Gates 0–2
- Gate 0: kill switch absent; `gh` = `sunholo-voight-kampff`; billing CLEAN; `mission_directives.sh` on #635: 0 directives from the allowlist (19 comments); step 6a: 0 crash and 0 transcript notices on #635/#574 (control 107:4); `mission-world` inbox 0 unread.
- Gate 1: local `dev` was 3 commits behind `origin/dev` (the iteration-439 record merge), so the charter and log were read from origin. Base `b5f41458…@2026-11-26T11:56:44Z`. Skill drift: the resolved symlink and the pin copy were each byte-identical to ailang `origin/dev` on all 14 files. Dev CI 3/3 green at `b5f41458`. (`0033ca11`'s CI run was red on the race job.)
- Gate 2: clause map 1.0 clauses 1–7 MET. Ledger 108 rows, 1 OPEN: D-WORLD-121 (row 267 before S5f2?), default B "from the next fire". So the pick was row 243 S5f2, record 439's Next, as groomed. No S5f2 plan, branch, PR or worktree existed. Claim `inbox_1795694365193_283d84da`.

## Routing
- Resolver, followed verbatim: designer `recipe claude:claude-opus-5-5` (not spawned; reviewed design), planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin` (derive-planner-lane agrees), executor `recipe codex:gpt-6.2-sol declared:provider-pin`, evaluator `agent-tool sonnet declared:alias-pin`.
- Codex probe rc 0 from the sprint worktree. The planner and executor ran as `codex exec` recipes, because the resolver returned `recipe` for both; the Agent tool is not a valid path for a `recipe` role. The evaluator ran via the Agent tool.

## Baseline (controller, pristine `b5f41458`, outside the sandbox)
- First recording (`base/`), co-loaded with the planner:
  - `store-race`: 0 fails, 42 s. `cmd-plain`: the 11 row-271 reds, 591 s. Both valid.
  - `daemon-race`, `cmd-race` and `verify_go` each hit Go's default 10-minute timeout (`panic: test timed out after 10m0s`). Their fail sets are truncated, so they are INSTRUMENT FAILURES, not base sets.
- Re-record (`base2/`, `-timeout 20m`): `daemon-race` 6 fails, 916 s (the known SeTools reds plus load-sensitive tests); `worktree-plain`, `worktree-race` and `store-plain` 0 fails.
- Not re-recorded: whole `cmd/ailang-worldd` race and `verify_go`. CI arbitrated both, as at iteration 439.

## Planner
- `codex exec`, workspace-write, 13:01–14:09 (~68 min), **914,431 tok**, rc 0. Plan + JSON committed `4dd16433`.
- Scope: M1 73, M2 647, M3 1,292, M4p **1,400** cumulative (`-- cmd host docs app`). Activation moves to S5f3 in Appendix C.
- The S5f3 prototype measures 1,984 cumulative (+584 over M4p). Its 410-line core plus M1–M3 alone reaches 1,702, over the 1,650 safe-unit ceiling.
- §1.4 premises found false: the upgrade path refuses an installed or loaded companion; the prototype reproduced two durable states that all three attended remedies refused; the CLI GC override bypassed workspace binding; terminal revoke lacked episode exclusion; the heartbeat published preset names without an eligible catalogue.
- Controller round trip `ctl_roundtrip.sh M1 M2 M3 M4p`: `diff -r` rc 0 on `cmd`, `host` and `docs` for each milestone.

## Executor
- Run 1 (14:12–14:14, ~75k tok) applied M1 cleanly, then **stopped correctly** on a defect in the controller's own directive: the directive required an `ACCEPT_OK <M> fast` last line, a marker iteration 439's `accept.sh` printed and this plan's (`set -euo pipefail`, exit 0 = pass) does not. False-green rule (4): the directive's gate list was never baselined. The controller reset the worktree to `4dd16433`, corrected the directive, and re-ran.
- Run 2 (~4 min, **75,210 tok**) applied M1, M2, M3 and M4p. `accept.sh` exit 0 each time (12/32/56/60 PASS). Snapshots 3/0/3, 12/0/12, 17/0/17 and 18/0/18 were `cmp`-equal to the planner stages.
- Controller commits (`commit_ms.sh`): M1 `9094d9db`, M2 `1bd64507`, M3 `651077c2`, M4p `c4fb23f6` (NUMSTAT 73 → 647 → 1,292 → 1,400). The final tree matched the M4p manifest, 18/18 files.

## Controller gates on `c4fb23f6`
- Clean legs: build, vet, Linux vet, test compile, five inventories, `verify_ail`, GC probe, the caller census (four definitions, zero production start callers) and the lock census.
- `host/store` and `host/worktree` plain and race: 0 fails each.
- `gofmt` red only on untouched `host/daemon/daemon.go` (red at base too).
- **`cmd-plain` added `TestOperatorStartCollisions` and `TestOperatorFilePublicationAdmission`** (plus load-sensitive `TestWorldContextAndStatusline`). This was the judge's B1.
- `daemon-race` drift against `base2`: five added tests (`TestSeToolsMCPEndToEnd`, `TestMarketingCaptureInventoryMutation`, `TestSessionSelfGit`, `TestWorkspaceEpisodeModuleRootOverrides`, `TestWorkspaceModuleRootIsTheSandbox`). All five pass alone under `-race` at head and at base (`daemon-drift-alone.txt`), so they are load-sensitive, not attributable.
- Controller named-mutant replay outside the sandbox: all 102 `mutants.json` specs, each on its own milestone stage. **101 KILLED, 1 SURVIVED** (`D-LOCK-NO-NONBLOCK`, declared equivalent; the judge reproduced it in a whole-`host/store` run and agreed). The 70 whole-package comparisons the planner listed (44 `host/daemon`, 26 `cmd/ailang-worldd`, ~5–10 min each) were NOT run.

## Judge round 1
- Agent `sonnet`, foreground, ~70 min, **128,239 tok**, own worktree at `c4fb23f6`: **FAIL 67** (23/14/16/8/6).
- **B1:** M2 made the `r.start` / `publishOperatorFilesFixture` helpers borrow the lock before the function under test. The `held` arms then hit EAGAIN in the helper, and `missing_lock` re-created the sidecar, so three existing refusal arms were red or dead.
- Controller reproduction: the three arms FAIL at `c4fb23f6` and PASS at `b5f41458` (same command, same rig).
- Non-blocking findings:
  - N1: `J-RR-FORCE` (`--force` in `RemoveRegistered`) survives.
  - N2: `episode gc --credentials-db` for a non-default workspace skips every worktree.
  - N3: GC leaves persistent `.operator-episode-<ep>.lock` sidecars.
  - N4: tests pin error wording.
  - N5: `daemon.go` gofmt (pre-existing).
  - N6: daemon drift (load).
- 33 of the judge's own mutants, plus controls, all killed except `J-RR-FORCE`; 7 ledger kills re-confirmed.

## Round 2
- Codex executor fix (15:42–15:48, **104,364 tok**), test files only: `git diff --quiet -- ':!*_test.go'` held. +66/−14 = 80 lines.
- The `held`, `missing_lock` and `lock_held` arms now pass a non-holder description, a sidecar unlinked after borrow, and another episode's description. Each drives the production helper and kills a named `Verify` mutation.
- New `TestRemoveRegisteredRefusesLockedOrDirty` kills `--force` (`/dirty`) and `--force --force` (`/locked`).
- Commit `c7c08e5b`, NUMSTAT **1,476**: 76 over the 1,400 cap, all of it this judge-mandated test fix.
- Controller round-2 gates:
  - gofmt on touched files, build, vet and Linux vet: rc 0.
  - `host/worktree` plain: 0 fails.
  - `cmd-plain`: base 11 plus `TestPresetInitSeToolsModeAndModel` (`git init … timed out after 30s`) and `TestPresetStartReadRefusals` (load average ~18). Both pass alone (`alone-r2.txt`); the sprint touches neither file.
  - The B1 tests are absent from the head fail set. The changed cmd tests pass alone under `-race` (57.3 s).
- **Judge round 2:** Agent `sonnet` (fresh instance), ~29 min, **64,235 tok**: **PASS 77** (23/25/16/8/5), 0 blocking.
  - B1, N1 and N6 CLOSED; N2 and N3 OPEN, accepted; N4 and N5 not defects.
  - Its own mutants: R2-SAMEFILE, R2-NAMECHECK, R2-WT-ENTRY, R2-FILES-VERIFY and R2-FORCE are each killed by a B1 arm or the N1 test.

## Gate 3b
- PR #685 CI on `c7c08e5b`: **3/3 success at attempt 1**.
- PR race job: `cmd/ailang-worldd` 273.2 s (269.4 s on S5f's PR), `host/daemon` 360.3 s (357.1 s), `host/verifygate` 334.9 s, `host/broker` 280.1 s, `host/store` 63.9 s, `host/worktree` 1.8 s.
- `origin/dev` had moved to `91d1cabf` (attended #683, #684; no file overlap). The PR was `MERGEABLE CLEAN` and was merged `--match-head-commit` as `150d0d96`.
- Merge CI: see the log entry.
- **Merge CI on `150d0d96`** (run 45852708959): **3/3 success at attempt 1**. Race job: `cmd/ailang-worldd` 351.0 s, `host/daemon` 409.3 s, `host/verifygate` 322.5 s, `host/broker` 254.0 s. Slower than the PR runner and under the 480 s bound. **LANDED.**
