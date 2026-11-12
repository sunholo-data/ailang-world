MISSION-ROLE: planner

You are the sprint-planner for Ailang World mission iteration 385 (unattended run; no human is present).
Write the sprint plan for **row 201 `w-world-approvals-by-key`, milestone A3a, SPRINT 2** — "the routes":
M4 the read routes, the origin set and the guards (AM1–AM3, RG1, RT1–RT3, RT7, CF1) and M5 the POST
(RT4–RT6, RT8–RT13, E2E1), exactly as the A3 child design r2 §6 "Sprint 2" and §9 list them.
No re-design: `design_docs/planned/w-world-approvals-by-key-a3.md` r2 is reviewed (quorum r1+r2 + controller
carve-out + independent audit VALID 84) and is NOT reopened. D-WORLD-114 (the doc's two amendments to the
attended parent) is still OPEN with default A = keep both, applied since iteration 384; build as designed.

## Where you work
- Sprint worktree (already created, branch `sprint/world-iter385-row201-a3a-s2`, base = origin/dev
  `a8c666554f284651de33d2f6a2c85a16f5b72c02`, which contains A3a sprint 1 merged as `b44b5f7d`: credentials
  schema 7 `surface_approvals` with `RecordSurfaceApproval`/`SurfaceApproval` in `host/store/surface_approvals.go`,
  the shared judge `host/authority/judge.go` (`JudgeAssertion`, `JudgeSSHSignature`), and the CLI minter
  `cmd/ailang-worldd/approve_await.go`):
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter385-a3a-s2`
- VERIFIED BY CONTROLLER at the base: production code has 0 hits for `requests.json` / `"/surface/approve`
  and 0 for `SurfaceOrigins|surface-origin` (control: `/surface/passkey` = 11 prod hits), so nothing of
  sprint 2 has landed.
- The design's line numbers were measured at `0140bd7f`; sprint 1 has since changed `cmd/` and `host/`. Re-measure
  every design line/premise you rely on at `a8c66655` and record drift in §1.2/§1.4. Read the landed sprint-1
  plan `design_docs/planned/w-world-approvals-by-key-a3a-sprint1-plan.md` §1.4 (P1–P6) — its premises are now facts
  of HEAD you build on.
- Scratch (compile your stages, run mutants) OUTSIDE the repo: `$HOME/.ailang/state/world-iter385/plan-scratch/`
  (never under /tmp — location-sensitive tests go red there).
- Write exactly two files in the worktree, and DO NOT git add/commit them (the controller commits):
  - `design_docs/planned/w-world-approvals-by-key-a3a-sprint2-plan.md`
  - `design_docs/planned/sprint_w-world-approvals-by-key-a3a-s2.json` (same shape as `sprint_w-world-approvals-by-key-a3a-s1.json`)
- Leave every code file in the worktree untouched (the executor applies your patches there). No git writes of any kind in the worktree.

## Read first
1. `design_docs/planned/w-world-approvals-by-key-a3.md` r2 in full, especially §0 (scope, sprint boundaries),
   §1, §2.1, §2.3 (the read routes), §2.4 (the POST, refusal arms R1–R14), §2.5, §2.7 only for which routes are
   sprint 3 (the card PAGE and approve.js are NOT this sprint — but RT1/RT7 name five routes; decide from §2/§6/§9
   precisely which routes this sprint registers and how RT1's "five routes" and `TestSurfaceAppRoutes` edits are
   satisfied without shipping sprint 3's page; record the reading in §1.4), §3, §4 (threats), §5, §6 "Sprint 2",
   §7 inventories, §8 conflict surface, §9 sprint-2 gates and the **fallback cut** (the SSHSIG arm: R11's
   `sshsig`, RT6, RT13 and AW6's route half move to sprint 3; the guards R2–R5, the origin set `O` and
   `--surface-origin`/CF1 are never cut), §11, §14 V-rows.
2. `design_docs/planned/w-world-approvals-by-key-a3a-sprint1-plan.md` for the plan FORMAT (it is the newest and
   it worked: §1.2 premise table, §1.4 false premises, §2 patch protocol + awk extraction, §6 gate script, §7
   executor notes, Appendix P fences, expected numstat per sub-milestone). Reproduce that format.
3. The code you build on: `host/daemon` (the mux, `Handler()`, `isProtected`, `surfaceCSP`, `ensureCredentialsStore`,
   the existing `/surface/*` routes and their tests incl. `TestSurfaceAppRoutes`, `TestFrozenV1RouteTableMatchesMux`),
   `host/surfacepage` (models), `host/store/surface_approvals.go`, `host/store/mint_requests.go`,
   `host/authority/judge.go`, `cmd/ailang-worldd/serve*.go` (flags, `Config`), `cmd/ailang-worldd/approve_await.go`,
   `host/approval/approvaltest`, and wherever `ParseCeremonyOrigin` lives.
4. `CLAUDE.md`, `design_docs/coding-standards.md` (binding: pure core, effects at the boundary, non-vacuous gates).

## Scope
- **M4** and **M5** as above. **Plus one carried test obligation from sprint 1's judge** (iteration 384,
  controller-reproduced): the minter judging at `p.row.CreatedAt` instead of the poll's `now`
  (`cmd/ailang-worldd/approve_await.go`, the two `JudgeAssertion/JudgeSSHSignature(..., now)` calls) leaves every
  await test green. Add an await analogue of `TestApproveByPasskeyJudgesTheApproverAtNow` (e.g. the approver is
  revoked, or enrolled, between the request's creation and the poll) that KILLS that mutation; prove it with the
  mutation applied (`sed 's/Artifact, now)/Artifact, p.row.CreatedAt)/'` changed 2 lines in 384). Put it in the
  earliest sub-milestone where it fits (it needs no route).
- **Size cap: ≤ 2,000 changed lines**, measured `git diff --numstat <base> HEAD -- cmd host` (state this exact
  filter). The design estimates ≈ 1,900 for sprint 2. If the plan exceeds 2,000, apply the design's fallback cut
  (SSHSIG arm to sprint 3) and record the moved scope in an Appendix C as sprint-3 scope (no re-design).
  Split into sub-milestones of ≤ ~250–350 changed lines each where possible (M4a/M4b/…/M5a/…), each
  independently green and safe at commit. **Safety ordering is the design's: guards and the Host/Origin checks
  before any route that writes; the POST must judge before it records (RT5's "recording before the law"
  mutant); the daemon never keeps the credentials writer (RT10, RT12).**

## Retro directives (iterations 374–384; each was paid for)
1. **Ship compiled patches**: every line of code and test as a unified diff in Appendix P (```` ````diff <M> ````
   fences), generated LAST from compiled, vetted, gofmt-clean, green scratch stages; then extract them FROM THE
   PLAN FILE with §2's awk, apply to a fresh `git archive a8c66655` export, and `diff -r` against your work
   tree — byte-identical or regenerate. The controller re-does this round trip independently. **Ship the patches
   BEFORE running the full mutant set, and checkpoint the plan file** (write a first complete version, then
   append the mutant ledger) — this took 69 min at 384 vs 218 at 382. Keep it.
2. **Expected numstat after each sub-milestone's commit** (measured after commit, `-- cmd host`).
3. **Mutants anchored to the DIFF**: one mutant per new production branch (every guard arm — Host, Origin incl.
   the set `O`, `Sec-Fetch-Site`, media type —, each refusal arm R1–R14, the judge-before-record order, the
   replay/consumed arms, writer-busy 503, the no-session/no-cookie refusal, the redaction of artifact/challenge in
   `requests.json`, the `--surface-origin` grammar check), enumerated from the diff, not from the test plan; run
   each against the named test row AND the package; record killed/survived; name equivalents with reasons; base
   controls where a test pins pre-existing behaviour (RT7's `/surface/model.json` control). A mutant that does not
   compile is not a mutant (use `-vet=off` and a compiling edit; pair with a killed control). Write the ledger to
   a file in scratch; every re-run appends its own line.
4. **Repo-wide inventories**: run every design §7 inventory by name with `-v` on your final scratch tree (incl.
   `TestFrozenV1RouteTableMatchesMux`, `TestSurfaceAppRoutes`, `TestDaemonDependencyAllowlist`,
   `TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots`, `TestProductionGoSurface`,
   `TestNoNewDeadlineFreeStoreReads`, `TestCredentialsMethodsRequireADeadline`, `TestEveryVerbHasHelp`, the
   quickstart tests, the new `TestDaemonCredentialWriterSites`, `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`,
   `TestApprovalImportsOnlyTheStandardLibrary`). Search for any test that enumerates daemon routes, serve flags or
   `Config` fields before assuming none.
5. The §6 gate script for the controller (one bash shell, `tee` every leg to `$HOME/.ailang/state/world-iter385/gates/`,
   never only through `tail`), incl. `-race` on touched packages, the inventories, `./scripts/verify_ail.sh`, and
   `VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh` after commit. Measure each
   gate at the base first. **Never start a gate under `nohup`**.
6. Executor notes as the sprint-1 plan's §7: no git writes; `patch -p1` is the only writer; runs of ~2–4
   sub-milestones; cumulative snapshots to `$HOME/.ailang/state/world-iter385/snap/<M>/` with MANIFEST.sha256
   (exclude the two plan files); stop conditions; numstat checks; exact acceptance commands per sub-milestone.
7. Acceptance criteria (§9-style), numbered, that the judge can score against.
8. **Every premise gets a §1.2 row**: any sentence saying "the code currently does X" cites the command and its
   observed output; an empty result needs a known-positive control in the same call. Where the design and HEAD
   disagree, measure, route minimally, record it in §1.4, and list every file you touch outside the design's §8.

## Hard constraints
- D-WORLD-109: test keys, temp stores, `approvaltest` passkeys/keys only. No real key, no `~/.ssh`, no
  `~/.ailang/world`, no real `/dev/tty`, no real browser profile, no mint or publish against a real store, no
  real network listener beyond `httptest`/loopback in tests. Never run bare `ssh-keygen`.
- `host/approval` stays pure and standard-library-only. Effects only at the store/CLI/daemon boundary.
  `surfaceCSP` is not changed. The daemon never KEEPS the credentials writer (the POST opens, records, closes).
  No CORS headers anywhere.
- Do not touch: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/` (except reading),
  `cmd/world-publish/*`, `~/.ailang/state/mission-v1*`, `/surface/passkey`, `passkey.js`, the A2 law.
- Export `AILANG_BIN=$HOME/.pinned-ailang/ailang` in every shell that runs Go tests.
- Your Bash tool's shell is zsh: run loops/globs/`$VAR` work under `bash -c` or `bash <<'EOF'`.
- Known load-sensitive / flaky tests (re-run once alone, do not chase): `TestSeedRestartIgnoresASurvivingOldCopy`
  (host/broker), `TestHostsLaunchersExecHost`, host/broker exec-probe timeouts incl.
  `TestExecProtectedCheckTimeout/reserve`, local `-race` of host/broker, host/daemon, host/verifygate,
  `TestDemoSeedPreflightRefusals`, `TestWorktreeSnapshotRefusesAConfiguredFilter`,
  `TestMissionConfigDefaultFlowIntegration/invalid_ledger`, `scripts/marketingcheck` supervisor tests,
  `TestMarketingCaptureFakeChromeSixShots/inventory`.
- **Time budget: aim to finish within ~100 minutes; hard stop at 150.** The iteration has a 6 h wall clock that
  started 00:05Z, and the executor, gates, judge, CI and record still follow you. Prefer a smaller,
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
