MISSION-ROLE: executor

You are the sprint-executor for Ailang World mission iteration 385 (unattended; no human present).
This is **Run B: sub-milestones M5a, M5b, M5c, in that order** (the plan's Run 3; Run A's M4a–M4e are already applied AND committed by the controller, so the tree starts clean at M4e).

Worktree: `/Users/voightkampff/dev/sunholo-data/.wt-world-iter385-a3a-s2` (branch `sprint/world-iter385-row201-a3a-s2`).
Plan: `design_docs/planned/w-world-approvals-by-key-a3a-sprint2-plan.md` in that worktree.

## What to do
1. Read the plan's §0, §2, §3 (the sections for M5a, M5b, M5c and the **[base]** lines), §5 and §7 with the
   Read tool. Do not read Appendix P by hand — §2's awk extracts each patch.
2. For each sub-milestone in order:
   a. Apply its patch EXACTLY with §2's bash heredoc (`--dry-run` first, then apply, then the .orig/.rej check).
   b. Run its §3 Acceptance in ONE bash shell (`bash <<'EOF' … EOF`), with `export AILANG_BIN=$HOME/.pinned-ailang/ailang`
      and `unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR`, `tee`-ing every command's output
      to `$HOME/.ailang/state/world-iter385/exec-logs/<M>-<n>.log` and capturing each command's rc without a pipe
      (e.g. `cmd > log 2>&1; rc=$?`), then grep the logs for `--- FAIL`, `--- SKIP` (in `-v` runs) and `^FAIL`.
   c. Snapshot with §7's snapshot command (`M=<that sub-milestone>`), cumulative.
3. Do NOT commit. Do NOT check numstat (the controller does that after committing each snapshot).

## Rules (stop conditions — stop and report, never work around)
- **No git writes** of any kind (`add`, `commit`, `stash`, `checkout`, `reset`, `restore`, `rebase`, `push`, `worktree`,
  `apply`). `patch -p1` per §2 is the only writer. No hand edits to any file.
- A `patch --dry-run` failure; a `--- FAIL`; a `--- SKIP` in a §3 `-v` run; a `FAIL` line; any need to touch a file
  outside §5; any reference to `~/.ssh`, `~/.ailang/world` or a real `/dev/tty`; never run Chrome, node, `ssh-keygen`.
- **One allowed exception for load-sensitive timing tests** in WHOLE-PACKAGE runs only: if a test in this list fails
  — `TestDemoSeedPreflightRefusals`, `TestSeedRestartIgnoresASurvivingOldCopy`, `TestHostsLaunchersExecHost`,
  `TestWorktreeSnapshotRefusesAConfiguredFilter`, `TestMissionConfigDefaultFlowIntegration/invalid_ledger`,
  `TestExecProtectedCheckTimeout/reserve` — re-run that ONE test alone with `-run '^<Name>$' -count=1 -v`, tee'd to its
  own log; if it passes alone, record both logs and continue. Any other failure: stop.
- M5a's acceptance takes ~6–8 min (whole `host/daemon`) and M5c's ~4–5 min (whole `cmd/ailang-worldd`). The machine may be loaded.
  Give the Bash tool a timeout of 600000 ms per call and split long legs across calls if needed.
- Your shell is zsh: wrap loops/globs/`$VAR` in `bash -c` or `bash <<'EOF'`.
- Waiting discipline: if you background anything, poll its artifact (an rc file it writes) inside a bounded `date +%s`
  loop in the same turn. A Monitor is an event stream, not a wait; nothing will re-invoke you. Never end your turn
  announcing a wait; end only when Run B is complete or stopped.

## Final message
For each sub-milestone: patch applied (yes/no), each acceptance command with its rc and log path, the PASS/FAIL/SKIP
counts you grepped, the snapshot path and its MANIFEST file count; any timing-exception re-runs; any stop and why.
