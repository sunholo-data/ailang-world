MISSION-ROLE: executor

You are the sprint-executor for Ailang World mission iteration 389 (unattended; no human present).
This is **Run A: sub-milestones M1a, M1b, M1c, M2, in that order** (the plan's Run 1).

Worktree: `/Users/voightkampff/dev/sunholo-data/.wt-world-iter389-a4-s1` (branch `sprint/world-iter389-row201-a4-s1`).
Plan: `design_docs/planned/w-world-approvals-by-key-a4-sprint1-plan.md` in that worktree (committed; do not edit it).

## What to do
1. Read the plan's §0, §2, §3 (the sections for M1a, M1b, M1c, M2 and the **[base]** lines), §5 and §7 with the Read tool.
   Do not read Appendix P by hand — §2's awk extracts each patch.
2. For each sub-milestone in order:
   a. Apply its patch EXACTLY with §2's bash heredoc (`--dry-run` first, then apply, then the .orig/.rej check).
   b. Run its §3 Acceptance in ONE bash shell (`bash <<'X' … X`), with `export AILANG_BIN=$HOME/.pinned-ailang/ailang`
      and `unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR` (this sprint has no browser
      leg), `tee`-ing every command's output to
      `$HOME/.ailang/state/world-iter389/exec-logs/<M>-<n>.log` and capturing each command's rc without a pipe
      (e.g. `cmd > log 2>&1; rc=$?`), then grep the logs for `--- FAIL`, `--- SKIP` (in `-v` runs) and `^FAIL`.
   c. Snapshot with §7's snapshot command (`M=<that sub-milestone>`), cumulative.
3. Do NOT commit. Do NOT check numstat (the controller does that after committing each snapshot).

## Rules (stop conditions — stop and report, never work around)
- **No git writes** of any kind (`add`, `commit`, `stash`, `checkout`, `reset`, `restore`, `rebase`, `push`, `worktree`,
  `apply`). `patch -p1` per §2 is the only writer. No hand edits to any file.
- **Never `nohup`, never background a test leg, never `&`.** Every acceptance command runs in the foreground.
  (Iteration 385's executor broke this once; do not.)
- The rig is very heavily loaded (load average ~200 on 16 cores). Go compiles are slow; be patient in the foreground.
- A `patch --dry-run` failure; a `--- FAIL`; a `--- SKIP` in a §3 `-v` run; a `FAIL` line; any need to touch a file
  outside §5; any reference to `~/.ssh`, `~/.ailang/world` or a real `/dev/tty`; never run `ssh-keygen`.
- **One allowed exception for load-sensitive timing tests** in WHOLE-PACKAGE runs only: if a test in this list fails
  — `TestDemoSeedPreflightRefusals`, `TestSeedRestartIgnoresASurvivingOldCopy`, `TestHostsLaunchersExecHost`,
  `TestWorktreeSnapshotRefusesAConfiguredFilter`, `TestMissionConfigDefaultFlowIntegration/invalid_ledger`,
  `TestExecProtectedCheckTimeout/reserve`, `TestMissionConfigDefaultFlowIntegration/valid_ledger`, `TestPublishApprovalRefusalSetWithALandedPositiveControl`, `TestMCPPostBudgetProductionConstants/batch`, `TestBackoff`, `TestCapsuleWarmEqualsColdForSeTools`,
  `TestMarketingCaptureInventoryMutation`, `TestRegistryPhaseHasItsOwnClock` — re-run that ONE test alone with `-run '^<Name>$' -count=1 -v`, tee'd to its
  own log; if it passes alone, record both logs and continue. Any other failure: stop.
  `TestServeStartsWhenTheProbeSpendsTheStartupClock` is RED AT THE BASE under this rig's load (the planner measured it failing alone at `9a6b959d`): if it fails, re-run it alone once, record both logs and continue — it is not a stop.
- Long legs (whole `host/daemon` ~5–10 min under the current rig load of ~200): give the Bash tool a timeout of
  600000 ms per call and split long legs across calls if needed.
- Your shell is zsh: wrap loops/globs/`$VAR` in `bash -c` or `bash <<'X'`.
- Waiting discipline: there is nothing to background. Never end your turn announcing a wait; end only when Run A is
  complete or stopped.

## Final message
For each sub-milestone: patch applied (yes/no), each acceptance command with its rc and log path, the PASS/FAIL/SKIP
counts you grepped, the snapshot path and its MANIFEST file count; any timing-exception re-runs; any stop and why.
