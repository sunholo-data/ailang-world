MISSION-ROLE: planner

You are the sprint-planner for Ailang World mission iteration 384 (unattended run; no human is present).
Write the sprint plan for **row 201 `w-world-approvals-by-key`, milestone A3a, SPRINT 1** — "the record and
the minter": M1 credentials schema 7 `surface_approvals`, M2 the shared judge `host/authority/judge.go`,
M3 `approve <request-id> --await [--await-for <duration>]`, exactly as the A3 child design r2 §9 lists them.
No re-design: `design_docs/planned/w-world-approvals-by-key-a3.md` r2 is reviewed (quorum r1+r2 + controller
carve-out + independent audit VALID 84) and is NOT reopened. D-WORLD-114 (the doc's two amendments to the
attended parent: no Tauri half in A3a; a CLI minter re-judges and mints) is OPEN with default A = keep both,
which applies from this iteration; build as designed.

## Where you work
- Sprint worktree (already created, branch `sprint/world-iter384-row201-a3a-s1`, base = origin/dev
  `0a25c0d6c6c9af6cdf81da75270f41af41a44bad`, which contains A2 sprints 1–4 merged, sprint 4 = `88ff1565`):
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter384-a3a`
- The design's V32 base-drift trigger was run by the controller at `0a25c0d6`: all 12 V1 sha256 prefixes equal
  the doc's, `cmd/` and `host/` are byte-identical to the design's base `0140bd7f`, V8 holds (5
  `OpenCredentialsReadOnly` sites in host/daemon, 0 writer opens besides `UpgradeCredentials`), V18 `--- PASS` ×4,
  V15 `TestPasskeyPageInChrome` PASS. So no r3; the doc's line numbers are current.
- Scratch (compile your stages, run mutants) OUTSIDE the repo: `$HOME/.ailang/state/world-iter384/plan-scratch/`
  (never under /tmp — location-sensitive tests go red there).
- Write exactly two files in the worktree, and DO NOT git add/commit them (the controller commits):
  - `design_docs/planned/w-world-approvals-by-key-a3a-sprint1-plan.md`
  - `design_docs/planned/sprint_w-world-approvals-by-key-a3a-s1.json` (same shape as `sprint_w-world-approvals-by-key-a2-s4.json`)
- Leave every code file in the worktree untouched (the executor applies your patches there). No git writes of any kind in the worktree.

## Read first
1. `design_docs/planned/w-world-approvals-by-key-a3.md` r2 in full, especially §0 (scope, sprint boundaries,
   *Amendment posture*), §1 (what the code does today), §2.1, §2.2 (schema 7 and its S-guards), §2.5 (the shared
   judge; the CLI refusal TEXTS must stay byte-identical — every existing `TestApprove*`/`TestPasskey*` stays
   unedited), §2.6 (the minter), §3, §4, §5, §6 "Sprint 1" rows SA1–SA4, JG1–JG2, AW1–AW7 and the edited list,
   §7 (repo-wide inventories), §8 (conflict surface), §9 (sprint 1 gates and the fallback cut), §11, §14 (V-rows).
2. `design_docs/planned/w-world-approvals-by-key-a2-sprint4-plan.md` (and the a2 sprint-3 plan) for the plan
   FORMAT: §1.2 premise table, §2 patch protocol + awk extraction, §6 gate script, §7 executor notes,
   Appendix P fences, expected numstat per sub-milestone. Reproduce that format — it has worked six sprints
   running (haiku executor, zero deviations, numstat equal at every step).
3. The code you build on: `host/store/credentials.go`, `host/store/mint_requests.go`, `host/store/credentials_schema.sql`
   (or wherever the DDL lives), `host/authority/mint.go`, `cmd/ailang-worldd/approve.go`, `approve_passkey.go`,
   `approve_sig.go`, `session.go` (the `env.after` seam, V29), `session_new.go` (`provisionEpisode`),
   `host/approval/approvaltest` (test keys/passkeys).
4. `CLAUDE.md`, `design_docs/coding-standards.md` (binding: pure core, effects at the boundary, non-vacuous gates).

## Scope
- **M1 schema 7** (SA1–SA4), **M2 the shared judge** (JG1, JG2; the CLI wrappers rewritten to call it with
  byte-identical refusal texts), **M3 `approve --await`** (AW1–AW7, help line, +1 context root).
- **Size cap: ≤ 2,000 changed lines**, measured `git diff --numstat <base> HEAD -- cmd host` (state this exact
  filter; `docs/` and `design_docs/` are outside it). The design's own **fallback cut** applies: if the plan
  exceeds 2,000 lines, M3 moves to the head of sprint 2 — put it in an Appendix C as the next sprint's scope
  (no re-design), just as the a2 sprint-3 plan did. Split into sub-milestones of ≤ ~250 changed lines each where
  possible (M1a/M1b/…), each independently green and safe at commit (a table with no writer/reader outside tests
  is safe; a minter that mints before the judge is shared is not).
- Optional carry-forward if cheap: none required. Do not widen scope into sprint 2 (routes) or sprint 3 (page).

## Retro directives (iterations 374–383; each was paid for)
1. **Ship compiled patches**: every line of code and test as a unified diff in Appendix P (```` ````diff <M> ````
   fences), generated LAST from compiled, vetted, gofmt-clean, green scratch stages; then extract them FROM THE
   PLAN FILE with §2's awk, apply to a fresh `git archive 0a25c0d6` export, and `diff -r` against your mutant/work
   tree — byte-identical or regenerate. The controller re-does this round trip independently. **Ship the patches
   BEFORE running the full mutant set**, and checkpoint the plan file (write a first complete version, then
   append the mutant ledger) — iteration 382's planner took 218 min re-running mutants.
2. **Expected numstat after each sub-milestone's commit** (measured after commit, `-- cmd host`).
3. **Mutants anchored to the DIFF**: one mutant per new production branch (every S-guard arm, the upgrade step,
   the prune extension, each judge branch, each `--await` flag rule, the poll/expiry/consumed/refusal-once
   branches, AW7's "trust the row" mutant), enumerated from the diff, not from the test plan; run each against
   the named test row AND the package; record killed/survived; name equivalents with reasons; base controls where
   a test pins pre-existing behaviour. A mutant that does not compile is not a mutant (use `-vet=off` and a
   compiling edit; pair with a killed control). Write the ledger to a file in scratch; every re-run appends its
   own line — a re-run in prose only is not evidence.
4. **Repo-wide inventories**: run every design §7 inventory by name with `-v` on your final scratch tree
   (`TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`, `TestDaemonDependencyAllowlist`,
   `TestWorldBoundaryDependencyAllowlist`, `TestEveryCommandDirectoryIsAProtectedGroup`, `TestProductionContextRoots`,
   `TestProductionGoSurface`, `TestNoNewDeadlineFreeStoreReads`, `TestCredentialsMethodsRequireADeadline`,
   `TestFrozenV1RouteTableMatchesMux`, `TestApprovalImportsOnlyTheStandardLibrary`,
   `TestNoParallelWriteForkPackagesOutsideVerifygate`, `TestEveryVerbHasHelp`, the quickstart tests, any schema-pin
   test that enumerates credentials schema versions/tables, and any test that inventories `host/authority` files
   or imports). Search for tests that pin `credentialsSchemaVersion` or the table list before assuming none.
5. The §6 gate script for the controller (one bash shell, `tee` every leg to `$HOME/.ailang/state/world-iter384/gates/`,
   never only through `tail`), incl. `-race` on touched packages, the inventories, `./scripts/verify_ail.sh`, and
   `VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh` after commit. Measure each
   gate at the base first. **Never start a gate under `nohup`**.
6. Executor notes as the a2 plans' §7: no git writes; `patch -p1` is the only writer; runs of ~2–4 sub-milestones;
   cumulative snapshots to `$HOME/.ailang/state/world-iter384/snap/<M>/` with MANIFEST.sha256 (exclude the two plan
   files); stop conditions; numstat checks; exact acceptance commands per sub-milestone.
7. Acceptance criteria (§9-style), numbered, that the judge can score against.
8. **Every premise gets a §1.2 row**: any sentence saying "the code currently does X" cites the command and its
   observed output; an empty result needs a known-positive control in the same call. Where the design and HEAD
   disagree, measure, route minimally, record it in §1.4, and list every file you touch outside the design's §8.

## Hard constraints
- D-WORLD-109: test keys, temp stores, `approvaltest` passkeys/keys only. No real key, no `~/.ssh`, no
  `~/.ailang/world`, no real `/dev/tty`, no real browser profile, no mint or publish against a real store.
  Never run bare `ssh-keygen`.
- `host/approval` stays pure and standard-library-only (not touched by A3a). Effects only at the
  store/CLI/daemon boundary. `surfaceCSP` is not changed. The daemon still never keeps the credentials writer.
- Do not touch: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/` (except reading),
  `cmd/world-publish/*`, `~/.ailang/state/mission-v1*`.
- Export `AILANG_BIN=$HOME/.pinned-ailang/ailang` in every shell that runs Go tests.
- Your Bash tool's shell is zsh: run loops/globs/`$VAR` work under `bash -c` or `bash <<'EOF'`.
- Known load-sensitive / flaky tests (re-run once alone, do not chase): `TestSeedRestartIgnoresASurvivingOldCopy`
  (host/broker), `TestHostsLaunchersExecHost`, host/broker exec-probe timeouts, local `-race` of host/broker,
  host/daemon, host/verifygate, `TestDemoSeedPreflightRefusals`, `TestWorktreeSnapshotRefusesAConfiguredFilter`,
  `TestMissionConfigDefaultFlowIntegration/invalid_ledger`, `scripts/marketingcheck` supervisor tests.
- **Time budget: aim to finish within ~100 minutes; hard stop at 150.** The iteration has a 6 h wall clock that
  started 20:05Z, and the executor, gates, judge, CI and record still follow you. Prefer a smaller,
  fully-proven sprint (use the fallback cut) over a larger one.

## Waiting discipline (binding on you)
If you background anything (a long test run, a mutant loop), poll its ARTIFACT (an output/rc file it writes on
completion) inside a bounded `date +%s` loop in the same turn. A Monitor is an event stream, not a wait; nothing
will re-invoke you when it fires. Never end your turn announcing a wait, a pause, or "once it completes" — end
only when both files are written and the round trip has passed.

## Your final message
Report: the plan path, scope (what is in, anything moved by the fallback cut and why, with measured numbers), the
sub-milestone list with expected numstat after each, mutant counts (total / killed / named equivalents) and the
ledger path, the round-trip result, base-gate results, every premise found false (§1.4), and anything the
controller must decide.
