# Iteration 444 — controller readings (row 243 V4 S6, launch)

Fire `world-1795769251-75944`, attempt 1. Controller `claude-opus-5-5` (session). Base `origin/dev` `b2b0ebab` (= row 267 merge `6447e776` + record #692).

## Gates 0–2
- Gate 0:
  - Kill switch absent (`mission-world.disabled`); `gh` = `sunholo-voight-kampff`; billing CLEAN.
  - `mission_directives.sh` found 0 allowlisted directives on #635 (26 comments) and on #574 (44 comments).
  - Step 6a's default 5 s `gh` timeout returned rc 2 (`gh timed out after 5s for issue #635`; a direct `gh api` read took 1.2 s). That is an instrument failure, not a verdict. Re-run with `--gh-timeout 30`: rc 0, 0 crash notices (control 107:4), 1 transcript-scan notice (fire `world-1795730668-41608`, 2 hits; key names only in the private controlplane notice).
- Gate 1:
  - Local `dev` == `origin/dev` == `b2b0ebab`. Base `b2b0ebab…@2026-11-27T08:49:51Z`.
  - Skill drift: resolved symlink and pin copy byte-identical to ailang `origin/dev` on 14/14 files. The positive control (the commit before the last `gate-3-route.md` change) differs, so the instrument fires.
  - Dev CI on `b2b0ebab` 3/3 success (run 45967606138).
- Gate 2:
  - Clause map: 1.0 clauses 1–7 MET, none UNMET.
  - Ledger: 1 OPEN, D-WORLD-121, moot since row 267 landed. The pick is row 243 V4 S6, the record-443 Next.
  - V4 design quorum: `mc_quorum_gate2` rc 0; r2 PROCEED 3/3 (`claude-sonnet-5@claude-p`, `gemini-4-argon`, `oc-glm-5-3`), no absent seat.
  - No S6 commit, PR or orphan worktree exists.

## Routing (resolver followed verbatim)
- Designer: not needed (doc reviewed). Planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`. Executor `recipe codex:gpt-6.2-sol declared:provider-pin`. Evaluator `agent-tool sonnet declared:alias-pin`.
- Codex probe rc 0.
- The operator's standing request asked for every role through the Agent tool "exactly as the routing table specifies". The table returned `recipe` for the planner and executor, so they ran as `codex exec` recipes. The REQUIRED judge ran on the Agent tool. No role failed to spawn.

## Baseline (controller, pristine `b2b0ebab`, outside the sandbox, co-loaded with the planner)
- store-plain 0, store-race 0, authority-plain 0, authority-race 0, worktree-plain 0 fails.
- verifygate-plain 1 (`TestMissionConfigDefaultFlowIntegration`).
- cmd-plain 12 (11 row-271 + `TestOperatorWorktreeDispositionCAS`), 1,090 s.
- daemon-race 9 (SeTools e2e), 1,624 s.

## Planner
- `codex exec`, ~87 min (10:02 → 11:29 local), **656,730 tok**, rc 0.
- Outcome A: two milestones, M1 512, M2 1,650 cumulative (the safe-unit exception, exactly at its ceiling).
- Controller-measured: `git diff --no-index --numstat` over a fresh `b2b0ebab` archive against the round-trip export gives cmd 1,054, host 582, docs 14 = **1,650**, unfiltered 1,650. Controller round trip `ctl_roundtrip.sh M1 M2`: `diff -r` rc 0 on cmd, host, docs for both.
- Premise false in the controller's own directive: the E3 raw path `~/.ailang/state/world-iter429/e3.out` no longer exists. The tracked reading is `design_docs/verification/world-iter429/readings.md`. The planner reported it correctly as absent.

## Executor
- `codex exec`, **76,995 tok**, rc 0. M1 `ACCEPT_OK M1 fast` (29 PASS). M2 `ACCEPT_OK M2 fast` (40 PASS). Snapshots 20/0/20 and 34/0/34, cmp-equal to the planner stages. No residue.
- `commit_ms.sh` gave M1 `6398a4a2` (NUMSTAT 512) and M2 `808ce248` (NUMSTAT 1,650). M2's commit gate included whole `host/verifygate`: `ok` 509.990 s.

## Controller gates (outside the sandbox, head `808ce248`)
- gofmt, build, vet, linux-vet, test-compile (`go test -run '^$' ./...`), `verify_ail.sh`: all rc 0.
- New tests with `TMPDIR` under a 1777 parent: 0 fails, 145 s. New tests `-count=3`: 0 fails.
- store plain/race, verbs, worktree race, agentrun race, authority: 0 new against base.
- `accept.sh M2 race`: first run rc 1 on `TestOperatorStartOnce/worktree` (`git checkout -q -b feature: timed out after 30s: context deadline exceeded`, the git-fixture load class; `/launch` passed). Re-run at load ~4: `ACCEPT_OK M2 race`.
- `callers` leg rc 2: the controller passed extra arguments to the planner's `callers.sh`, an instrument error. The judge ran it correctly: `launchPreparedOperatorStart` has exactly 1 production caller (`operator_runner.go:118`).
- Planner's `outside.sh` (stage trees, byte-identical to the commits):
  - m1-cmd-plain new vs base: `TestDemoSeedCommitFailures`, `TestDemoSeedPreflightRefusals`, `TestOperatorWorktreeWindow`.
  - m2-cmd-plain new vs base: `TestDemoSeedPreflightRefusals`, `TestOperatorFilePublicationAdmission`.
  - All four pass alone at head (126.6 s, load ~4), so this is load with no diff-attributable fail.
  - The whole-package `cmd/ailang-worldd` race legs TRUNCATED at the 15 m bound for base, M1 and M2 alike under co-load, so locally they are uninformative. CI arbitrates, under its 900 s per-package bound.
- Whole-package mutants: the controller replayed the 9 non-cmd ids of the planner's 25. All KILLED, each with a named killer (`ledger-whole-noncmd.jsonl`).

## Judge r1 (Agent `sonnet`, cross-vendor; ~57 min, 161,452 tok, 48 tool calls)
- **PASS 73** (25/17/16/8/7), 0 blocking.
- Launch safety: no second child is possible, restart adoption cannot adopt a live process, the pre-exec rechecks are present, and scrubbing works on real child output.
- Survivors (test gaps, code traced correct):
  - T1: the real worker/supervisor path is never executed by a test. `operatorSupervisor` seam; J16 no ExtraFiles, J17 supervisor cap +1000 h.
  - T2: pre-exec refusal arms unpinned. J09 signed-binary hash, J10 launch-dir recheck, J22 `agentrun` ignores the `Before` error.
  - T3: J25, no group kill when `Started` errors.
  - Stuck launch rows have no settle path (episode fenced until a manual repair).
  - The child inherits the lock fds 3/4.
- Controller reproductions:
  - J09: `if false && serviceFileHash(in.binary) != in.digest` in the eval worktree, `go vet` ok, launch selector `ok` 5.574 s. **SURVIVES**.
  - T1: `grep -rn 'runOperatorLaunchWorker\|superviseOperatorLaunch' --include='*_test.go'` = 0, with production references at `operator_launch.go:212,215,268` and `main.go:396` as the control.
- Production reachability: `store.Credentials.ConfirmOperatorAction` (`host/store/operator_cas.go:42`) has no non-test caller (the `approval.ConfirmOperatorAction` call in `host/authority/judge_operator.go:66` is the firing control). So no production path creates a confirmed start, and S6's launch is dormant until S7 adds the confirmation routes.
- Decision: land on PASS 73 and carry T1–T3 as a precondition of S7. A round 2 (codex test-only + judge r2 + two CI cycles) did not fit the remaining slot (3h07m of 6h used at the decision).

## Gate 3b round 1 — PR #695 on `808ce248`: RED
- Run 45994531303: `go host build + test gate` and `go race gate` both failed, on one test only: `host/broker TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`. The output read "14 files carry subprocess sites … but 13 have drivers". The undriven site was `cmd/ailang-worldd/operator_launch.go:217`, `exec.Command(e.paths.Binary, "operator-launch-worker", …)`, and it set no `cmd.Env`, so the worker inherited the runner's whole environment.
- Attribution: dev `b2b0ebab` was 3/3 green; the guard is a whole-repo AST scan.
- The controller's touched-package sweep never ran `host/broker`. This is rule 3g, the third iteration running where a whole-repo guard outside the touched packages caught the sprint (441 `TestSessionRootBudgets`/`TestSubprocessSinkGate`, 442 the Linux fixture, 444 this guard).
- The red independently confirmed the judge's T1.

## Round 2 (codex `gpt-6.2-sol`, 155,961 tok, rc 0; 250 lines) → `b40e2f48`
- What changed:
  - The worker spawn moved to `host/agentrun.StartWorker`, with an explicit `cmd.Env = childenv.Scrubbed(os.Environ())` and Setpgid. `cmd/ailang-worldd` now has 0 `exec.Command` sites.
  - `driveAgentrunWorker` was added to the guard's agentrun entry.
  - New tests: `TestOperatorLaunchDispatchRechecks/{binary_digest,nonempty_directory}` (J09/J10), `TestRunObservedBeforeRefusalPreventsStart` (J22), `TestRunObservedStartedRefusalKillsGroup` (J25) and `TestStartWorkerEnvironmentAndFiles` (ENV).
- Inside the sandbox the full broker guard was UNINFORMATIVE (loopback bind denied).
- Controller re-run outside the sandbox:
  - full broker guard `ok` 24.2 s;
  - whole `host/agentrun` `ok` 16.1 s;
  - launch selector + dispatch rechecks `ok` 10.4 s;
  - build, vet and linux vet clean.
- Sprint total: `git diff -M --numstat b2b0ebab b40e2f48 -- cmd host docs app` = **1,888**, past the 1,650 safe-unit exception (FLAG).

## Judge r2 (the r1 judge resumed by name; 179,347 tok cumulative)
- **PASS 80** (26/24/17/9/4), 0 blocking.
- J09, J10, J22, J25 and ENV each re-run: each FAILS its named test. Un-mutated, the new tests take 0.00–0.01 s. They pass at `-count=3`, under `-race` and under `GOMAXPROCS=1`.
- r1 findings:
  - T2 closed. T3 closed at the agentrun level.
  - T1 partly closed: `StartWorker` is driven, but `runOperatorLaunchWorker` and the supervisor cap (J17) are still unexecuted.
  - Still open: J06, stuck launch rows (no attended settle verb) and the child inheriting lock fds 3/4.
- Scope scored 4/10 for the 1,888-line total.

## Gate 3b round 2
- PR #695 on `b40e2f48`: **3/3 at attempt 1** (run 45999595851). Race: `cmd/ailang-worldd` 287.863 s, `host/daemon` 356.323 s, `host/verifygate` 325.488 s, `host/broker` 278.316 s, `host/store` 66.761 s.
- `MERGEABLE CLEAN`; merged `--match-head-commit` → **`00b4ec55`**.
- Merge `00b4ec55` CI: **3/3 at attempt 1** (run 46001578504). Race: `cmd/ailang-worldd` 291.041 s, `host/daemon` 358.622 s, `host/verifygate` 323.035 s, `host/broker` 275.474 s. **LANDED.**
