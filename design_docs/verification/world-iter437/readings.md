# Iteration 437 — readings (row 243 V4 S5e3)

All times are UTC. The fire started 2026-11-25 23:18 and these readings end 2026-11-26. Rig: Voights-Mac-Studio.local. Artifact times were read with `date -u -r`.

## Gate 0–2
- **Gate 0.** Kill switch absent (`~/.ailang/state/mission-world.disabled`). `gh` account is sunholo-voight-kampff. Billing CLEAN.
  - `mission_directives.sh` on #635: 0 allowlisted directives in 16 comments. The watermark is a LEGACY bare timestamp.
  - Step 6a, `mission_gate0_self_notices.sh`: rc 0, with 0 crash and 0 transcript notices on #635 and #574. Control issue 107 expected 4 and got 4. The `-prev` watermark is absent, so this was a single-watermark read.
  - The `mission-world` inbox has 0 unread. The `world` and `ailang-world` inboxes each hold one Sep/Oct unread already triaged by earlier fires.
- **Driver routing at fire start.**
  - Planner and executor: `codex:gpt-6.2-sol`. Evaluator: `sonnet`. Designer seed: `claude:claude-opus-5-5`.
  - `MISSION_ANTHROPIC_AVAILABLE=1`. Over-ration bucket: `google`. `MISSION_ROUTING_NOTE` empty.
  - Pinned driver `63b605cbc`, `AILANG_DRIVER_DRIFT=1`.
- **Gate 1.** Local `dev` was 5 behind `origin/dev` (at `90e97b1f`): the S5e2 plan `81ebfd7e`, M1 `3795994c`, the #672 merge, the iteration-436 record and the #673 merge. Fast-forwarded with no local changes. `dev` = `origin/dev` = `eaa202b1d23f601ce597f063034a061b5bf5dea3`. Base recorded `eaa202b1…@2026-11-25T23:18:25Z`.
- **Skill drift.** Checked the resolved symlink (the ailang main checkout) and the pin worktree `~/.ailang-driver-pin/world`, 14 files each. All 28 are byte-identical to `origin/dev`, with no uncommitted skill edit.
- **Dev CI at Gate 1.** `commits/eaa202b1/check-runs`: checks=3, all success.
- **Open PRs from the bot account.** #617 `ui/live-branch-heads` is attended work with no mission worktree. Left alone.
- **Clause map.** 1.0 clauses 1–7 are MET and none is UNMET. Row 243 serves D-WORLD-94. Record 436's Next = V4 S5e3.
- **Ledger.** `mission_decisions.sh --check` valid at 107 rows; `--open` empty.
- **Quorum** for `w-world-address-bar-verbs-v4`, via the checked reader (rc 0): r1 `…T10-13-56Z` was blocked (reject/pass/pass) and r2 `…T10-24-51Z` proceeded (pass ×3), with `absent_reviewers` = []. The design last changed in `8dc91744`.
- **Died-mid-flight traces.** None: no S5e3 plan, branch, PR or `world-iter437` worktree existed before this fire.
- **Carried premise.** Attended #671 (`90e97b1f`) overlaps the iteration-436 deferred publication patch.
  - `comm -12` of `git diff --name-only 5cac9316 eaa202b1 -- cmd host docs app` (42 files) against the patch's 14 headers gives 4 shared files: `cmd/ailang-worldd/episode.go`, `cmd/ailang-worldd/main.go`, `docs/QUICKSTART.md`, `host/store/credentials_test.go`.
  - Positive control: the 3 `operator_worktree*` test files of S5e2 appear in the changed list.
- **Claim** sent to `controlplane` before routing (`inbox_1795648762897_f859dc3d`).

## Gate 3
- **Resolver.**
  - Planner: `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`, and `derive-planner-lane.sh` gives the same answer.
  - Executor: `recipe codex:gpt-6.2-sol declared:provider-pin`.
  - Evaluator: `agent-tool sonnet declared:alias-pin`.
- **Codex probe.**
  - The first probe ran from `/tmp` and returned rc 1: "Not inside a trusted directory and --skip-git-repo-check was not specified". This was an instrument error (the wrong CWD), not a lane failure.
  - Re-run from the sprint worktree: rc 0, replied `ok`.
  - Codex usage was 51.0% used against 86.3% allowed.
- **Designer** not spawned: the design was already reviewed and no new doc was needed. Rotation pointer unchanged.
- **Sprint worktree.** `.wt-world-iter437-row243-v4s5e3` from `eaa202b1`. The `mission-base.sh snap` read equalled the Gate 1 base, so there was no drift. `mission-worktree.sh wait` reported ready.
- **Planner.** `codex exec`, `workspace-write`, 23:22–00:09, about 47 min, 430,131 tok. The controller committed the plan and JSON as `d34c358b`.
  - **M1 is the whole sprint: 1,650 lines** (1,625+/25−, 16 files). The directive allowed up to 1,650 for the minimum safe unit.
  - The pure rebase over #671 reproduced the inherited 1,259 lines (§1.4 F2). Hardening added 391.
  - The #671 allowance credential origin is refused before reservation, using the real `GrantAllowance`/`MintFromAllowance` producer (`TestOperatorFileAllowanceOrigin`).
  - Named mutants: 95 of 98 KILLED and 3 MASKED (`F-CREATE-NOFOLLOW`, `R-POST-PARENT`, `S-WORKTREE-PHASE`). The planner's own diff census lists 122 AUDIT lines with no independent mutation measurement.
  - The re-run GC probe still lists both a reserved and a published operator worktree as candidates at 73 h.
- **Controller round trip.** `ctl_roundtrip.sh M1` extracted `diff M1` (1,895 lines) onto a fresh `git archive eaa202b1`. `diff -r` against `st-M1` was rc 0 for cmd, host and docs.
- **Executor.** `codex exec`, sandboxed, 00:09:45–00:12:17, about 2.5 min, 77,895 tok.
  - Patch clean. `ACCEPT_OK M1 fast`: 320 PASS, 0 FAIL, 0 SKIP.
  - Snapshot 16/0/16, all 16 files `cmp`-equal to `st-M1`.
  - `commit_ms.sh M1` ran gofmt, build, vet and acceptance outside the sandbox, then committed `de608652` with NUMSTAT 1650.
- **Controller baseline.** `baseline-controller.sh` in a pristine detached worktree `.wt-world-iter437-base` at `eaa202b1`, about 00:14–00:37. Load ran 3.6→9.0, with fseventsd at 99% CPU at the start.
  - Plain and race `cmd/ailang-worldd` fail sets were each the 11 row-271 CLI tests exactly.
  - **`verify_go` came back with 60 failures, an INSTRUMENT FAILURE.** The plan's `gate-bin/git` PATH wrapper ran `exec "$WORLD_GATE_GIT"`. Tests that scrub their child env (host/broker and others) drop that variable, so every `git` they spawned died with `exec: : not found`. That accounts for 43 extra host/broker failures (44 in the contaminated set against 1 in the clean one); the CLI failures were identical in both runs.
  - **Controller fix (scratch only):** the wrapper now falls back with `${WORLD_GATE_GIT:-/opt/homebrew/bin/git}`. Under `env -i PATH=gate-bin:/usr/bin:/bin` the fixed wrapper gives rc 0 and the original (`git.planner-orig`) gives rc 127.
  - The contaminated files are kept as `*.contaminated-gitwrapper`.
  - The base `verify_go` set alone was re-recorded in the pristine base worktree, 00:37–00:48. It returned rc 1 with a firing parser control: **18 failures** = row 271's 14, plus `TestWorldContextAndStatusline`, `host/broker TestExecProtectedCheckTimeout`, `host/daemon TestSeToolsRunBenchmarkTasksEndToEnd` and `host/verifygate TestMissionConfigDefaultFlowIntegration`. All four are on the known load-sensitive list in the evaluator directive; two (`TestExecProtectedCheckTimeout`, `TestSeToolsRunBenchmarkTasksEndToEnd`) are also charter-listed.
- **Gates** on `de608652`, using plan `gates.sh … final` with the fixed wrapper and `WORLD_GATE_OBSERVE=1`, 00:48–01:15, 1,622 s:
  - Every leg was rc 0 except `compare-verify_go` (below): gofmt, build, vet, linux-vet, compile, verify_ail, plain and race store, inventories, s5e2-race, accept, files-race 48 s, callers, and observation before/after/compare.
  - `pkg-cmd-ailang-worldd` (285 s) and `race-cmd-ailang-worldd` (456 s, against 405 s at base) both had fail sets equal to base.
  - `compare-verify_go` was rc 1 only because head had ONE FEWER failure: `host/broker TestExecProtectedCheckTimeout` passed. No new failure at head.
- **Evaluator.** Agent `sonnet`, background, 00:38–02:36, about 119 min, 133,078 tok, 56 tool uses. It worked in its own worktree `.wt-world-iter437-row243-v4s5e3-eval` at `de608652` and sequenced on `gates-final.rc`.
  - **PASS 79, 0 blocking.** Report in `evaluator-r1.md` here; ledger `~/.ailang/state/world-iter437/judge/mutants.tsv` (70 lines, 36 ids not in the plan).
- **Controller reproduction of judge finding 3**, in the judge worktree, restored clean afterwards. The edit dropped `now >= row.ExpiresAt ||` from `host/store/operator_files.go:191` (N31); `git diff --stat` showed 1 file and `go vet ./host/store` was ok.
  - `go test -count=1 ./host/store/` **rc 0**.
  - `go test ./cmd/ailang-worldd -run '^TestOperatorFile'` **rc 0**.
  - So it survives both. It is inert today: no production caller.

## Gate 3b
- **PR #674** head `de608652`, CI run 45731365295:
  - **Attempt 1** was red only on the race job: `FAIL host/daemon 480.093s`, which is row 267's signature in a package this sprint does not touch. The same log shows `cmd/ailang-worldd` race at **477.565 s**.
  - Dev race times for `cmd/ailang-worldd` on that workflow were 219.863 s at `d9c24a3b` and 301.979 s at `eaa202b1`. `host/daemon` was 354.9 s and 398.9 s.
  - **Attempt 2** (failed jobs re-run): 3/3 success. `cmd/ailang-worldd` race took **262.975 s** and `host/daemon` 353.079 s.
  - So attempt 1 was a slow runner (row 267 instance 11). The sprint's real local cost is about +51 s (456 against 405 s).
- **Merge.** `origin/dev` was unchanged at `eaa202b1` and the PR was `MERGEABLE CLEAN`. Merged with `--match-head-commit de608652` → `22bb7804`.
- **Merge CI** on `22bb7804`, run 45743770171. Tree-identical to `de608652`: `git diff --quiet de608652 22bb7804` rc 0.
  - **Attempt 1** was red on two jobs:
    - Build+test job: `--- FAIL: TestSeedRestartIgnoresASurvivingOldCopy/survivor`, `FAIL host/broker 234.926s`. This test is charter-listed as load-sensitive.
    - Race job: `FAIL host/daemon 480.125s` (row 267 instance 12). `cmd/ailang-worldd` race took **468.977 s** and `host/broker` 426.477 s.
    - `git diff --quiet eaa202b1 22bb7804 -- host/broker host/daemon` rc 0, so neither package changed.
  - **Attempt 2** (failed jobs re-run): **3/3 success**. `cmd/ailang-worldd` race took 259.477 s, `host/broker` 279.714 s and `host/daemon` 353.736 s. **LANDED.**

## Gate 4
- **Record worktree** `.wt-world-iter437-record` from `origin/dev` and fast-forwarded to `22bb7804`. Base `22bb7804…@2026-11-26T03:11:45Z`.
- **STATUS rotation.** Moved 434 to the archive with `edit_charter.py`, which asserts the charter length is unchanged (6,366 lines by `wc -l` before and after).
  - The charter has `^## STATUS 2026` = 3.
  - The archive carries 434, and control 433 is present.
- **Ledger** valid at 107 rows. The new stamp matches the canonical regex.
