# World iteration 449 — readings (row 243 V4 S9b)

All readings are first-party controller measurements unless marked otherwise. Paths under `~/.ailang/state/world-iter449/` are
rig-local scratch and are not published.

## Gate 0–2
- Kill switch absent (`mission-world.disabled`); `gh` = `sunholo-voight-kampff`; billing tripwire CLEAN.
- `mission_directives.sh --issue 635`: 0 allowlisted directives (36 comments; LEGACY watermark `2026-08-17T06:15:12Z`).
- Step 6a `mission_gate0_self_notices.sh` rc 0: 0 crash notices (control 107:4 ok; `#574` watermark absent → DEGRADED single-watermark
  read). 4 transcript-scan notices in the window: fires `world-1795730668-41608` (2 hits), `world-1795769251-75944` (1),
  `world-1795786215-30397` (3) — reported by iterations 446–448 — and one new, `world-1795842267-52962` (3 hits), iteration 448's fire.
  Key names are in the private controlplane notice only.
- `mission-world` inbox: 0 unread. `world` and `ailang-world` inboxes hold the same two stale messages iteration 448 recorded.
- Local `dev` `9a40fe68` was 2 behind `origin/dev` `77dca5ed` (iteration 448's record merge); state read from origin. No open
  own-account PRs. Dev CI on `77dca5ed`: `checks=3`, all `success`.
- Skill drift: the symlink-resolved copy differed from origin on `resources/gate-5-retro.md` only (the idle-fingerprint input list), as
  in 448; no CWD skill copy exists in this repo. Gate 5 read origin's version.
- Ledger valid, 109 rows, 1 OPEN (D-WORLD-122, not an S9b input; default B — fakes). Clause map: 1.0 clauses 1–7 MET, none UNMET.
- Pick = record-448 Next, row 243 V4 S9b. Quorum reader rc 0: V4 r2 `…T10-24-51Z` PROCEED 3/3, `absent_reviewers` `[]`.
- Verified at HEAD before routing: `git diff --stat c1869819 77dca5ed -- cmd host docs app` empty; the iteration-448 `M3.patch`
  dry-run applies rc 0 with no offset/fuzz/reject; `frozen.Kind != "finish"` at `operator_finish_preflight.go:23`.

## Routing
- Resolver: planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin` against the V4 doc (a doc-less call returned
  `agent-tool opus fail-closed:no-doc`; `derive-planner-lane.sh` on the doc = codex, followed); executor `recipe codex:gpt-6.2-sol
  declared:provider-pin`; evaluator `agent-tool sonnet declared:alias-pin`; no designer.
- Codex probe rc 0; codex 17.0% used / 25.9% allowed (10080 m window).
- Worktrees via `mission-worktree.sh` at `77dca5ed` (sprint, base) and at `f11d09c7` (judge, branch `eval/world-iter449-r1`).

## Baseline (pristine `77dca5ed`, outside the sandbox, isolated `AILANG_STATE_DIR`, during the planner run)
store/authority/agentrun plain+race and verbs: 0 fails. verifygate 1 (`TestMissionConfigDefaultFlowIntegration`); broker-plain 4
(`TestAilangRunCapsOnTheToolBinary`, `TestExecProtectedCheckTimeout`, `TestPolicyToolWorktreeDeltaOnlyForToolchainOps`,
`TestWorktreeSnapshotRefusesAConfiguredFilter`); cmd-plain 16 (the 11 row-271 real-binary tests and 5 load-class); daemon-plain 8;
daemon-race 10 (974 s). `ALL_DONE for 77dca5ed…`.

## Planner (codex `gpt-6.2-sol`, ~67 min, rc 0, 538,978 tok)
- Outcome **B**: rebased whole 2,461 lines; S9b = M1 401 + M2 → 1,218 cumulative; M3 = S9c activation, 1,243 incremental, compiled in
  Appendix P with five activation gates listed.
- Its own §7 block extracted and executed verbatim: `APPLY_ACCEPT_SNAPSHOT_OK` M1/M2/M3, `FINAL_APPLY_RC 0`.
- Plan-reported (not first-party): 53 named mutants killed incl. 19 judge-448 carries; census private entry/local commit 0 → 0 → 1,
  claim 1 → 1 → 2. Premises corrected: a phantom path in the 257-name cap, an unverified author header, a loose object accepted as
  check evidence without a matching World log entry, invalid SQL bindings in judge-448's `e-store-commit-once`, an empty workspace in the
  inherited boot catalogue (S9c).
- Controller round trip (`ctl_roundtrip.sh M1 M2`, fresh `git archive 77dca5ed`): `diff -r` rc 0 for cmd, host and docs at both stages.
- Destructive-write scan of `accept.sh`/`outside.sh`/`gates.sh`/`whole-mutants.sh`/`commit-guards.sh`/`snapshot.sh`: no `rm`, no
  real-path write; writes under scratch or `mktemp -d`.

## Outside-only legs on the planner's stage, before the executor
- `outside.sh st-M1 M1` rc 0: `ACCEPT_OK M1 outside`, `outside-race` (`RACE_BUDGET cmd/ailang-worldd n=6/6 added=3.84s cap=15s`),
  `outside-1777`.
- `outside.sh st-M2 M2` rc 1: `ACCEPT_OK M2 outside` and `outside-1777`; `outside-race` assertions all PASS but `RACE_BUDGET n=16/16
  added=111.61s cap=15s` (`TestOperatorCommitReconcile` 35.1 s, `TestOperatorFinishUnsandboxed` 19.3 s, `TestOperatorCommitFailed`
  10.8 s), load average ~5 on 16 cores with `corespotlightd` ~140% CPU and a long-running `godot` job. Re-runs alone 3/3: 6.52 s,
  6.43 s, 6.71 s, each `ACCEPT_OK M2 outside-race`. Attributed to rig load; handed to the judge (O1).

## Executor (codex `gpt-6.2-sol`, ~5 min, rc 0, 118,719 tok)
§7 block verbatim (`TOOLS=$STATE/exec-tools`): `ACCEPT_OK`/`APPLY_ACCEPT_SNAPSHOT_OK` M1 and M2 fast, `FINAL_APPLY_RC 0`; snapshots
5/0/5 and 10/0/10, `shasum -c` OK, all files `cmp`-equal to the planner's stages; all 17 extracted tools `cmp`-equal; no `.orig`/`.rej`.

## Commits (quiet rig; no controller sweep running — the iteration-448 process fix)
- **M1 `5717be11`** on `ACCEPT_OK M1 outside` (NUMSTAT 401), first attempt.
- **M2 `f11d09c7`** on `ACCEPT_OK M2 outside` and `COMMIT_GUARDS_OK` (whole verifygate + broker, `-timeout 20m`, no new fail vs the
  base sets), first attempt, NUMSTAT 1,218. Load 2.99 at start.

## Controller sweep on `f11d09c7` (outside the sandbox, before the judge ran anything)
gofmt, build, vet, linux vet, `go test -run '^$' ./...`, guard-store, guard-cli, guard-boundary, guard-daemon, `TestSubprocessSinkGate`,
`TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`, `callers.sh` (`CALLER_CENSUS_OK`: local-commit 0, attended-commit 1,
push/label/view-PR/create-PR/V3-put 1 each), `verify_ail.sh`, checkrun/worktree race, store race 42 s: all rc 0.

## Judge (Agent `sonnet`, ~55 min, 105,809 tok): PASS 78 (25/21/17/7/8), 0 blocking
- Reachability: no production caller of the private worker, the private entry or the local-commit seam; claim/push/PR/label/V3-put
  unchanged; nothing new exported. Proof scope: a signed start and a finish without confirm-on-surface are refused with a no-change
  read-back; the kind and surface checks are load-bearing (j01, j02 killed); the ID-vs-CAS check is redundant (j03 survived).
- 35 own mutants; all 19 judge-448 survivor edits the plan claims to pin are now killed, including `e-store-commit-once`. Survivors:
  j03, j08 (earlier check-receipt World-evidence recheck on resume), j12, j13 (`commit-failed` recorded without an input recheck first),
  j14 (resume from a reserved-but-unobserved commit), j16 (reconcile config re-read), j17, j19 (evidence SemanticID), j25. Must-live
  controls j34/j35 survived as intended.
- O1: not reproduced on a quiet rig; the new tests inject a 3-minute budget, an injected clock and a fake check runner, so load slows
  them but does not change their verdicts. Whole `cmd` race delta UNMEASURED. Hostile-repo scratch cases UNMEASURED beyond the repo's own
  tests and mutants.
- Landing: safe; second consecutive sprint with production code and no production caller, accepted as dead weight owned by S9c.
- **Controller reproduction of one survivor:** replacing `return w.reconcileCommit(j.CommitKey)` at `operator_finish_local.go:101` with an
  error compiles and every `TestOperatorFinish*`/`TestOperatorCommit*` test stays `ok` (29.7 s) — survives, as reported. Restored
  (`cmp` equal); judge worktree clean.

## Gate 3b
- PR #709 on `f11d09c7`: CI 3/3 at attempt 1 (run 46148073561) → MERGEABLE/CLEAN, `origin/dev` `77dca5ed` unmoved → merged
  `--match-head-commit` → **`d9635675`**. Merge CI 3/3 at attempt 1 (run 46149287254): race `cmd/ailang-worldd` 292.3 s, `host/daemon`
  366.1 s, `host/verifygate` 315.4 s, `host/broker` 274.8 s, `host/store` 66.0 s — all under the 900 s bound. **LANDED.**
