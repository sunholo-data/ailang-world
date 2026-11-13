MISSION-ROLE: executor

You are the sprint-executor for Ailang World mission iteration 390 (unattended; no human present).
This is **run2: sub-milestones M4d, M4e, M4f, M4g, in that order** (the plan's Run 2). Run 1 (M4a–M4c) has been applied, accepted and committed by the controller; HEAD contains it.

Worktree: `/Users/voightkampff/dev/sunholo-data/.wt-world-iter390-a4-s2` (branch `sprint/world-iter390-row201-a4-s2`).
Plan: `design_docs/planned/w-world-approvals-by-key-a4-sprint2-plan.md` in that worktree (committed; do not edit it).

## What to do
1. Read the plan's §0, §2, §3 (the **[base]** lines and the sections for M4d, M4e, M4f, M4g), §5 and §7 with the Read tool.
   Do not read Appendix P by hand — §2's awk extracts each patch.
2. For each sub-milestone in order:
   a. Apply its patch EXACTLY with §2's bash heredoc (`--dry-run` first, then apply, then the .orig/.rej check).
   b. Run its §3 Acceptance in bash (`bash <<'X' … X`), with `export AILANG_BIN=$HOME/.pinned-ailang/ailang`
      and `unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR` (no browser leg),
      writing every command's output to `$HOME/.ailang/state/world-iter390/exec-logs/<M>-<n>.log` and capturing
      each command's rc without a pipe (e.g. `cmd > log 2>&1; rc=$?`), then grep the logs for `--- FAIL`,
      `--- SKIP` (in `-v` runs) and `^FAIL`.
   c. Snapshot with §7's snapshot command (`M=<that sub-milestone>`), cumulative.
3. Do NOT commit. Do NOT check numstat (the controller does that after committing each snapshot).

## Rules (stop conditions — stop and report, never work around)
- **No git writes** of any kind (`add`, `commit`, `stash`, `checkout`, `reset`, `restore`, `rebase`, `push`,
  `worktree`, `apply`). `patch -p1` per §2 is the only writer. No hand edits to any file.
- **EVERY Bash call that runs go build/vet/test MUST pass the Bash tool's `timeout: 600000` parameter and run in the
  FOREGROUND** (no `run_in_background`, no `&`, no `nohup`, no Monitor). The tool's default timeout is 120 s:
  iteration 389's executor hit it, got pushed into a background shell, ended its turn on a wait and lost the shell.
  Split long legs across separate calls (the fence; each `-v` run; each whole-package run), each its own
  foreground call with its own log and rc. If one call cannot hold a whole-package leg in 10 min, split it by `-run`.
- Rig load average is ~30 on 16 cores; Go compiles are slowish. Be patient in the foreground.
- Stop on: a `patch --dry-run` failure; a `--- FAIL`; a `--- SKIP` in a §3 `-v` run; a `FAIL` line; any need to
  touch a file outside §5; any reference to `~/.ssh`, `~/.ailang/world`, a real `/dev/tty` or a real key. The
  only `ssh-keygen` allowed is the one the plan's tests run on `approvaltest` key files in temp dirs.
- **One allowed exception for load-sensitive timing tests** in WHOLE-PACKAGE runs only: if a test in this list fails
  — `TestDemoSeedPreflightRefusals`, `TestSeedRestartIgnoresASurvivingOldCopy`, `TestHostsLaunchersExecHost`,
  `TestWorktreeSnapshotRefusesAConfiguredFilter`, `TestMissionConfigDefaultFlowIntegration`,
  `TestExecProtectedCheckTimeout`, `TestPublishApprovalRefusalSetWithALandedPositiveControl`,
  `TestMCPPostBudgetProductionConstants`, `TestBackoff`, `TestCapsuleWarmEqualsColdForSeTools`,
  `TestMarketingCaptureInventoryMutation`, `TestRegistryPhaseHasItsOwnClock`, `TestSessionSelfGit`,
  `TestServeStartsWhenTheProbeSpendsTheStartupClock`, `TestWorkspaceProtectedPolicyWiring`,
  `TestConvertedStoreMethodsM5b` — re-run that ONE test alone with `-run '^<Name>$' -count=1 -v`, logged to its
  own file; if it passes alone, record both logs and continue (`TestServeStartsWhenTheProbeSpendsTheStartupClock`
  may stay red alone — it is red at base under load; record it and continue). Any other failure: stop.
- Your shell is zsh: wrap loops/globs/`$VAR` in `bash -c` or `bash <<'X'`.
- Waiting discipline: there is nothing to background. Never end your turn announcing a wait, a pause or "once it
  completes"; end only when this run is complete or stopped.

## Final message
For each sub-milestone: patch applied (yes/no), each acceptance command with its rc and log path, the PASS/FAIL/SKIP
counts you grepped, the snapshot path and its MANIFEST file count; any timing-exception re-runs; any stop and why.
