MISSION-ROLE: executor

You are the sprint-executor for Ailang World mission iteration 387 (unattended; no human present).
This is **Run A: sub-milestones M9a, M9b, M9c, in that order** (the plan's Run 1).

Worktree: `/Users/voightkampff/dev/sunholo-data/.wt-world-iter387-a3a-s4` (branch `sprint/world-iter387-row201-a3a-s4`).
Plan: `design_docs/planned/w-world-approvals-by-key-a3a-sprint4-plan.md` in that worktree (committed; do not edit it).

## What to do
1. Read the plan's §0, §2, §3 (the sections for M9a, M9b, M9c and the **[base]** lines), §5 and §7 with the Read tool.
   Do not read Appendix P by hand — §2's awk extracts each patch.
2. For each sub-milestone in order:
   a. Apply its patch EXACTLY with §2's bash heredoc (`--dry-run` first, then apply, then the .orig/.rej check).
   b. Run its §3 Acceptance in ONE bash shell (`bash <<'X' … X`), with `export AILANG_BIN=$HOME/.pinned-ailang/ailang`
      and `unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR` (a `[chrome]` leg sets the two
      Chrome variables on its own command line, exactly as §3 writes it), `tee`-ing every command's output to
      `$HOME/.ailang/state/world-iter387/exec-logs/<M>-<n>.log` and capturing each command's rc without a pipe
      (e.g. `cmd > log 2>&1; rc=$?`), then grep the logs for `--- FAIL`, `--- SKIP` (in `-v` runs) and `^FAIL`.
   c. Snapshot with §7's snapshot command (`M=<that sub-milestone>`), cumulative.
3. Do NOT commit. Do NOT check numstat (the controller does that after committing each snapshot).

## Rules (stop conditions — stop and report, never work around)
- **No git writes** of any kind (`add`, `commit`, `stash`, `checkout`, `reset`, `restore`, `rebase`, `push`, `worktree`,
  `apply`). `patch -p1` per §2 is the only writer. No hand edits to any file.
- **Never `nohup`, never background a test leg, never `&`.** Every acceptance command runs in the foreground.
  (Iteration 385's executor broke this once; do not.)
- A `patch --dry-run` failure; a `--- FAIL`; a `--- SKIP` in a §3 `-v` run; a `FAIL` line; any need to touch a file
  outside §5; any reference to `~/.ssh`, `~/.ailang/world` or a real `/dev/tty`; never run `ssh-keygen`; Chrome and
  node ONLY through the plan's `go test` commands (headless, temporary profile, virtual authenticator), never by hand.
- **One allowed exception for load-sensitive timing tests** in WHOLE-PACKAGE runs only: if a test in this list fails
  — `TestDemoSeedPreflightRefusals`, `TestSeedRestartIgnoresASurvivingOldCopy`, `TestHostsLaunchersExecHost`,
  `TestWorktreeSnapshotRefusesAConfiguredFilter`, `TestMissionConfigDefaultFlowIntegration/invalid_ledger`,
  `TestExecProtectedCheckTimeout/reserve`, `TestMissionConfigDefaultFlowIntegration/valid_ledger`, `TestPublishApprovalRefusalSetWithALandedPositiveControl`, `TestMCPPostBudgetProductionConstants/batch`, `TestBackoff` — re-run that ONE test alone with `-run '^<Name>$' -count=1 -v`, tee'd to its
  own log; if it passes alone, record both logs and continue. Any other failure: stop. A Chrome deadline test
  (`TestApproveCardNetworkDeadlines`, `TestApproveCardPasskeyDeadline`, `TestNeedsYouNotifiesOnlyWhenAskedOnce`) that fails on a timing window may be
  re-run ONCE alone the same way; record both logs.
- Long legs (whole `host/daemon` ~5–7 min; whole `cmd/ailang-worldd` ~3–5 min): give the Bash tool a timeout of
  600000 ms per call and split long legs across calls if needed.
- Your shell is zsh: wrap loops/globs/`$VAR` in `bash -c` or `bash <<'X'`.
- Waiting discipline: there is nothing to background. Never end your turn announcing a wait; end only when Run A is
  complete or stopped.

## Final message
For each sub-milestone: patch applied (yes/no), each acceptance command with its rc and log path, the PASS/FAIL/SKIP
counts you grepped, the snapshot path and its MANIFEST file count; any timing-exception re-runs; any stop and why.
