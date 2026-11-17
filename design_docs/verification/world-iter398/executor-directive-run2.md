MISSION-ROLE: executor

You are the sprint-executor for Ailang World mission iteration 398 (unattended; no human present).
This is **Run 2: sub-milestones M5r, M6e, in that order**. The plan is committed (`1a244fee`) and M4, M6p, M5g are applied AND COMMITTED (HEAD `55d7c197`); M5r and M6e are not applied
Start at M5r.

Worktree: `/Users/voightkampff/dev/sunholo-data/.wt-world-iter398-row243-v2b` (branch `sprint/world-iter398-row243-v2b`).
Plan: `design_docs/planned/w-world-address-bar-verbs-v2b-sprint-plan.md` in that worktree (committed; do not edit it).

## What to do
1. Read the plan's §0, §2, §3 (the acceptance protocol, the `accept.sh` block, and the sections for M5r, M6e),
   §5 and §7 with the Read tool. Do not read Appendix P by hand — §2's awk extracts each patch.
2. For each sub-milestone in order:
   a. Apply its patch EXACTLY with §2's bash heredoc (set `M=`; `--dry-run` first, then apply with `-E`, then the
      .orig/.rej check). Pre-existing committed `*.orig` files under `design_docs/verification/` are not yours; §2's
      check looks only in cmd host docs.
   b. Run its §3 Acceptance with §3's extract-and-run heredoc (set `M=`). It must end `ACCEPT_OK <M>`; keep its log
      under `$HOME/.ailang/state/world-iter398/exec-logs/`. Run it as ONE foreground Bash call with the tool's
      `timeout: 600000`. If a single call cannot finish within 600 s, do NOT background it: run the [base] lines and
      each step line of that sub-milestone's arm as separate foreground calls (each `> log 2>&1; rc=$?`, its own log
      under exec-logs/), and grep each log for `--- FAIL`, `--- SKIP` (in `-v` runs) and `^FAIL`.
   c. Snapshot with the plan's §7 command (set `M=`). Check the snapshot against §7's "Expected" line (STATUS lines,
      DELETED count, MANIFEST count); stop on a mismatch.
3. Do NOT commit. Do NOT check numstat (the controller commits each snapshot and checks numstat).

## Rules (stop conditions — stop and report, never work around)
- **No git writes** of any kind (`add`, `rm`, `mv`, `commit`, `stash`, `checkout`, `reset`, `restore`, `rebase`, `push`,
  `worktree`, `apply`). `/usr/bin/patch -p1 -E` per §2 is the only writer. No hand edits to any file, no `gofmt -w`,
  no `WORLD_UPDATE_VERB_GOLDEN=1`, no `WORLD_UPDATE_WHY_GOLDEN=1`.
- **EVERY Bash call that runs go build/vet/test MUST pass the Bash tool's `timeout: 600000` parameter and run in the
  FOREGROUND** (no `run_in_background`, no `&`, no `nohup`, no Monitor). The default is 120 s.
- Rig load varies; Go compiles take tens of seconds. Be patient in the foreground.
- Stop on: a `patch --dry-run` failure; `ACCEPT_RED`; a `--- FAIL`; a `--- SKIP` in a `-v` run; a `FAIL` line; any
  changed file outside §5's list; any reference to `~/.ssh`, `~/.ailang/world`, a real `/dev/tty`, a real key, a real
  browser profile or a real daemon (never dial `127.0.0.1:7644`).
- **Live Chrome exception (plan §7/§8 R6):** a live-Chrome `x` step that fails on a timeout under load is re-run ONCE,
  alone (`go test -count=1 -v -run '^<Name>$' <pkg>` with the same env as accept.sh), logged to its own file; if it
  passes alone, record both logs and continue with the remaining steps as separate calls. A second failure: stop.
- **Load-sensitive timing tests** in WHOLE-PACKAGE `p` runs only: if a test in this list fails —
  `TestDemoSeedPreflightRefusals`, `TestSeedRestartIgnoresASurvivingOldCopy`, `TestHostsLaunchersExecHost`,
  `TestWorktreeSnapshotRefusesAConfiguredFilter`, `TestMissionConfigDefaultFlowIntegration`,
  `TestExecProtectedCheckTimeout`, `TestMCPPostBudgetProductionConstants`, `TestBackoff`,
  `TestMarketingCaptureInventoryMutation`, `TestRegistryPhaseHasItsOwnClock`, `TestSessionSelfGit`,
  `TestServeStartsWhenTheProbeSpendsTheStartupClock`, `TestWorkspaceProtectedPolicyWiring`,
  `TestCLIRealSubprocessEpisode`, or any test that fails only on a subprocess / `ailang --version` / exec-profile /
  startup-clock deadline and is NOT a `TestWhy*`, `TestVerb*`, `TestSurfaceVerb*`, `TestAdmitWhy*` or other test this
  sprint adds or edits — re-run that ONE test alone with `go test -count=1 -v -run '^<Name>$' <pkg>`, logged to its
  own file; if it passes alone, record both logs and continue (then run the rest of that arm's steps as separate
  calls). Any other failure: stop.
- Your shell is zsh: wrap loops/globs/`$VAR` in `bash -c` or `bash <<'X2'`.
- Waiting discipline: there is nothing to background. Never end your turn announcing a wait, a pause or "once it
  completes"; end only when this run is complete or stopped.

## Final message
For each sub-milestone: patch applied (yes/no), the acceptance result (`ACCEPT_OK`/`ACCEPT_RED` and the log path, or
each split command with its rc and log path), the PASS/FAIL/SKIP counts you grepped, the snapshot path, its STATUS
line count, its DELETED count and its MANIFEST file count; any exception re-runs; any stop and why.
