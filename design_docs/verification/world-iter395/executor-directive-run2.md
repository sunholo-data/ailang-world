MISSION-ROLE: executor

You are the sprint-executor for Ailang World mission iteration 395 (unattended; no human present).
This is **Run 2: sub-milestones M2, M3a, M3b, in that order.** M1a, M1b and M1c are applied and committed by the controller (`9e24bb2b`; the tree is clean). Start at M2.

Worktree: `/Users/voightkampff/dev/sunholo-data/.wt-world-iter395-row243-v1` (branch `sprint/world-iter395-row243-v1`).
Plan: `design_docs/planned/w-world-address-bar-verbs-v1-sprint-plan.md` in that worktree (committed; do not edit it).

## What to do
1. Read the plan's §0, §2, §3 (the `accept.sh` block and the sections for M2, M3a, M3b), §5 and §7 with the Read tool.
   Do not read Appendix P by hand — §2's awk extracts each patch.
2. For each sub-milestone in order:
   a. Apply its patch EXACTLY with §2's bash heredoc (set `M=`; `--dry-run` first, then apply, then the .orig/.rej
      check). Pre-existing committed `*.orig` files under `design_docs/verification/` are not yours; §2's check looks
      only in cmd host docs.
   b. Run its §3 Acceptance with §3's extract-and-run heredoc (set `M=`). It writes
      `$HOME/.ailang/state/world-iter395/exec-logs/accept_<M>.log` and must end `ACCEPT_OK <M>`. Run it as ONE
      foreground Bash call with the tool's `timeout: 600000`. If a single call cannot finish within 600 s, do NOT
      background it: run the [base] lines and each `step` line of that sub-milestone's arm as separate foreground
      calls (each `> log 2>&1; rc=$?`, its own log under exec-logs/), and grep each log for `--- FAIL`, `--- SKIP`
      (in `-v` runs) and `^FAIL`.
   c. Snapshot with the plan's §7 command, with ONE controller correction: in its MANIFEST line replace
      `find cmd host docs -type f` with `find $(for d in cmd host docs; do [ -d "$d" ] && echo "$d"; done) -type f`
      (§7 aborts under `set -euo pipefail` when the snapshot has no `docs/` directory — measured iteration 392).
      Stop on any `git status` code other than ` M` and `??`.
3. Do NOT commit. Do NOT check numstat (the controller commits each snapshot and checks numstat).

## Rules (stop conditions — stop and report, never work around)
- **No git writes** of any kind (`add`, `commit`, `stash`, `checkout`, `reset`, `restore`, `rebase`, `push`,
  `worktree`, `apply`). `patch -p1` per §2 is the only writer. No hand edits to any file, no `gofmt -w`, no
  `WORLD_UPDATE_VERB_GOLDEN=1`.
- **EVERY Bash call that runs go build/vet/test MUST pass the Bash tool's `timeout: 600000` parameter and run in the
  FOREGROUND** (no `run_in_background`, no `&`, no `nohup`, no Monitor). The default is 120 s: an executor once hit
  it, got pushed into a background shell, ended its turn on a wait and lost the shell.
- Rig load average is ~7–15 on 16 cores; Go compiles take tens of seconds. Be patient in the foreground.
- Stop on: a `patch --dry-run` failure; `ACCEPT_RED`; a `--- FAIL`; a `--- SKIP` in a `-v` run; a `FAIL` line; any
  need to touch a file outside §5; any reference to `~/.ssh`, `~/.ailang/world`, a real `/dev/tty`, a real key, a
  real browser profile or a real daemon (never dial `127.0.0.1:7644`).
- **One allowed exception for load-sensitive timing tests** in WHOLE-PACKAGE `p` runs only: if a test in this list
  fails — `TestDemoSeedPreflightRefusals`, `TestSeedRestartIgnoresASurvivingOldCopy`, `TestHostsLaunchersExecHost`,
  `TestWorktreeSnapshotRefusesAConfiguredFilter`, `TestMissionConfigDefaultFlowIntegration`,
  `TestExecProtectedCheckTimeout`, `TestMCPPostBudgetProductionConstants`, `TestBackoff`,
  `TestMarketingCaptureInventoryMutation`, `TestRegistryPhaseHasItsOwnClock`, `TestSessionSelfGit`,
  `TestServeStartsWhenTheProbeSpendsTheStartupClock`, `TestWorkspaceProtectedPolicyWiring`, or any test that fails
  only on a subprocess / `ailang --version` / exec-profile deadline in an UNTOUCHED package — re-run that ONE test
  alone with `go test -count=1 -v -run '^<Name>$' <pkg>`, logged to its own file; if it passes alone, record both
  logs and continue (then run the rest of that arm's steps as separate calls). Any other failure: stop.
- Your shell is zsh: wrap loops/globs/`$VAR` in `bash -c` or `bash <<'X2'`.
- Waiting discipline: there is nothing to background. Never end your turn announcing a wait, a pause or "once it
  completes"; end only when this run is complete or stopped.

## Final message
For each sub-milestone: patch applied (yes/no), the acceptance result (`ACCEPT_OK`/`ACCEPT_RED` and the log path, or
each split command with its rc and log path), the PASS/FAIL/SKIP counts you grepped, the snapshot path and its
MANIFEST file count; any timing-exception re-runs; any stop and why.
