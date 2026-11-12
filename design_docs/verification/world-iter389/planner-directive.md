MISSION-ROLE: planner

You are the sprint-planner for Ailang World mission iteration 389 (unattended run; no human is present).
Write the sprint plan for **row 201 `w-world-approvals-by-key`, milestone A4, SPRINT 1 — "the law and the typed
root" (M1–M3)**. No re-design: `design_docs/planned/w-world-approvals-by-key-a4.md` r2 is reviewed (quorum r1+r2
blocked on narrow defects, controller carve-out, independent sonnet audit VALID 84, iteration 388) and is NOT
reopened. D-WORLD-115 (the doc's one amendment to the attended parent: the decision message with a derived
challenge) is still OPEN; its default **A = keep it** applies from this iteration (389). Build as designed. The
design says sprint 1 does not depend on the D-WORLD-115 answer; if you find a place where it does, record it in §1.4.

## Scope (design §0 "Sprint 1", §6 "Sprint 1" table, §9 "Sprint 1 = M1–M3" — read them verbatim)
- **M1 the message, the law and the discovery**: `host/approval/decision.go` (the 91-byte decision message, §2.3),
  `host/broker/decision_law.go` (the pure refusal ladder D1–D10 incl. D4a/D4b, §2.5), the complete request-wide
  packet discovery in `host/coordinator` (§2.5 D4a, `Discovery{Complete,…}`, DS1), the reserved-scope predicate
  (§2.6 L1, R4a). Tests DM1, DM2, DL1–DL12, DL4a, DL4b, DS1, R4a.
- **M2 the broker rules**: the broker refuses a second decision of the same request (§2.4, BR1/BR2), the publish
  consumer refuses `ssh:`/`webauthn:` decisions (§2.6 L2, R4b/R4c), the caller census XG1.
- **M3 the subject, the typed POST, the typed CLI**: `GET /surface/decisions/{ref}`, `POST /surface/decide/{ref}`
  typed arm only (behind the A3a guards), `ailang-worldd approve <sha256:request-ref> [--deny]` typed. Tests RD1–RD9,
  RD9b–RD9d, R4d, CD1–CD4.
- **Audit note N4 (iteration 388's carve-out audit, MUST be folded in)**: the design's K0 (idempotent retry on the
  exact artifact of an unlanded key row, §2.7) is a sprint-2 concern, but its place in the Judge and handler-order
  bullets of §2.8 must be consistent. If any sprint-1 code (the handler order of the POST, the law's inputs) fixes an
  order that sprint 2's K0 must slot into, write that order down in the plan so sprint 2 does not have to reorder
  sprint-1 code; record it in §1.3 as a decision.
- **Fallback cut (design §9):** CD1–CD4 (the typed CLI) move to sprint 2's head if the sprint exceeds the cap. The
  discovery, D4a and RD9b–RD9d are never cut. M1 and M2 are never cut. Decide from the measured size and record it.

## Where you work
- Sprint worktree (already created, branch `sprint/world-iter389-row201-a4-s1`, base = origin/dev
  `9a6b959deed573d7966863a1f166f3a8c5700f2e`, which contains A3a code-complete and the A4 design):
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter389-a4-s1`
- VERIFIED BY CONTROLLER at the base (commands run in the worktree):
  - `grep -rl "DecideApproval(" cmd host --include="*.go"` → `host/broker/approve.go` (the only non-test file; it is
    the definition) + 5 `_test.go` files (`context_thread_test.go`, `episode_test.go`, `registry_publish_test.go`,
    `workbench_decisions_test.go`, `workbench_still_test.go`). So there is no production caller of the exported
    decide path today (design V3).
  - `host/coordinator/decisionpacket_validate.go:34` `const packetScanLimit = 100`.
  - `host/broker/approve.go:517` `func validatePublishApproval(ctx context.Context,`.
  - Pinned ailang: `$HOME/.pinned-ailang/ailang` = AILANG v0.41.0.
  - **The rig is very heavily loaded (load average ≈ 210 on 16 cores from other sessions).** Expect slow Go builds
    and timing-sensitive test reds in untouched packages; re-run once alone before believing a red.
- Read the landed A3a sprint plans for the plan FORMAT, newest first:
  `design_docs/planned/w-world-approvals-by-key-a3a-sprint4-plan.md`, then the sprint-3 one (§1.2 premise table, §1.4
  false premises, §2 patch protocol + awk extraction, §6 gate script, §7 executor notes, Appendix P fences, expected
  numstat per sub-milestone, Appendix C). Reproduce that format.
- The design's line numbers were measured at `b0fbcd3e`; re-measure every line/premise you rely on at `9a6b959d`
  and record drift in §1.2/§1.4.
- Scratch (compile your stages, run mutants) OUTSIDE the repo:
  `$HOME/.ailang/state/world-iter389/plan-scratch/` (never under /tmp — location-sensitive tests go red there).
- Write exactly two files in the worktree, and DO NOT git add/commit them (the controller commits):
  - `design_docs/planned/w-world-approvals-by-key-a4-sprint1-plan.md`
  - `design_docs/planned/sprint_w-world-approvals-by-key-a4-s1.json` (same shape as `sprint_w-world-approvals-by-key-a3a-s4.json`)
- Leave every code/doc file in the worktree untouched (the executor applies your patches there). No git writes of any
  kind in the worktree.

## Read first
1. `design_docs/planned/w-world-approvals-by-key-a4.md` r2 in full — §0, §1 (what the code does today), §2.1–§2.6,
   §2.8 (routes, handler order), §2.10 (CLI), §3, §4, §5, §6 Sprint 1, §7, §8, §9, §14 V-rows, §15.
   Re-measure what you rely on, live.
2. The parent `design_docs/planned/w-world-approvals-by-key.md` §2.2 and the A3 child
   `design_docs/planned/w-world-approvals-by-key-a3.md` (the A3a guards R0–R6 that the new POST reuses).
3. The code you build on: `host/broker/approve.go` (`DecideApproval`, `decideApproval`, `validatePublishApproval`,
   `humanMu`), `host/coordinator/decisionpacket*.go`, `host/approval/` (the mint message `signed.go`/`request.go`
   for the encoding convention), `host/daemon/surface_approve*.go` and the A3a guard helpers, `host/daemon/daemon.go`
   route registration, `host/daemon/episode_finish.go` (`writeOperatorProof`, `RecordOperatorApproval`),
   `cmd/ailang-worldd/approve*.go` (the id-shape dispatch; F1 flag rules), the inventories in §7.
4. `CLAUDE.md`, `design_docs/coding-standards.md` (binding: pure core, effects at the boundary, non-vacuous gates).

## Size rule
- **Cap: ≤ 2,000 changed lines**, measured `git diff --numstat <base> HEAD -- cmd host` (state this exact filter).
  The design estimates ≈ 1,950 for sprint 1, i.e. at the cap — measure your real size and apply the design's priced
  fallback cut (CD1–CD4 to sprint 2) if needed. Prefer a smaller, fully proven sprint.
- Sub-milestones of ≤ ~250–350 changed lines each, each independently green and safe at commit (e.g. M1a message,
  M1b law, M1c discovery, M1d reserved scope, M2a/M2b, M3a subject route, M3b POST guards+law arms, M3c typed arm,
  M3d CLI — your call, measured).
- No new `/v1` route (`TestFrozenV1RouteTableMatchesMux` stays 13). `surfaceCSP` unchanged. No CORS headers.

## Retro directives (iterations 374–388; each was paid for)
1. **Ship compiled patches**: every line of code, test and doc as a unified diff in Appendix P
   (```` ````diff <M> ```` fences), generated LAST from compiled, vetted, gofmt-clean, green scratch stages; then
   extract them FROM THE PLAN FILE with §2's awk, apply to a fresh `git archive 9a6b959d` export, and `diff -r`
   against your work tree — byte-identical or regenerate. The controller re-does this round trip independently.
   **Ship the patches BEFORE running the full mutant set, and checkpoint the plan file** (write a first complete
   version, then append the mutant ledger).
2. **Expected numstat after each sub-milestone's commit** (measured after commit, `-- cmd host`).
3. **Mutants anchored to the DIFF**: one mutant per new production branch, enumerated from the diff, not from the
   test plan — in particular **each Dn branch of the law deleted** (design 3j: the unit is the branch), D4a reading
   `Carrying` instead of `Complete`, the discovery's sentinel `>=`/`>` and cursor, the decision message field order and
   endianness, the reserved-scope prefix and parse arms, the no-redecide check, L2, every POST guard, the typed proof
   consumption and verb binding, the CLI's subject-hash check. Run each against the named test row AND the package;
   record killed/survived; name equivalents with reasons. A mutant that does not compile is not a mutant (pair with a
   killed control). **Name the test that killed each mutant; a mutant killed only by an inventory/static test is not
   behavioural proof — re-run with that test excluded and record both.** Write the ledger to a file in scratch; every
   re-run appends its own line. **The plan's mutant table ids must be the same ids as in the ledger** (the iteration
   387 judge could not reconcile renamed rows).
4. **New branch reachability (iteration 388 retro)**: for every new refusal/branch, give a trace from the request's
   entry point (the POST handler, the CLI) to that branch for each input class, and say which test drives it end to
   end vs unit-only. A branch that is unreachable end to end because an earlier step refuses first must say so (the
   iteration-388 audit found exactly this for a revoked-approver arm).
5. **Repo-wide inventories**: run every design §7 inventory by name with `-v` on your final scratch tree (incl.
   `TestFrozenV1RouteTableMatchesMux`, `TestSurfaceAppRoutes`, `TestDaemonDependencyAllowlist`,
   `TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots`, `TestProductionGoSurface`,
   `TestEveryVerbHasHelp`, the quickstart tests, `TestDaemonCredentialWriterSites`,
   `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`, import allowlists for `host/approval`). Search for
   any test that enumerates routes, exported symbols, files, CLI verbs or QUICKSTART sections before assuming none.
6. The §6 gate script for the controller (one bash shell, `tee` every leg to
   `$HOME/.ailang/state/world-iter389/gates/`, never only through `tail`), incl. `-race` on touched Go packages, the
   inventories, `./scripts/verify_ail.sh`, and `VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang
   ./scripts/verify_go.sh` after commit. Measure each gate at the base first. **Never start a gate under `nohup`**.
7. Executor notes as the A3a plans' §7: no git writes; `patch -p1` is the only writer; runs of ~2–4 sub-milestones;
   cumulative snapshots to `$HOME/.ailang/state/world-iter389/snap/<M>/` with MANIFEST.sha256 (exclude the two plan
   files); stop conditions; numstat checks; exact acceptance commands per sub-milestone; **an explicit "never
   `nohup`, never background a test leg" line**.
8. Acceptance criteria (§9-style), numbered, that the judge can score against. Include what an operator can do after
   this sprint (approve/deny a packet at a terminal, if CD1–CD4 ship) and what the attended exit run (design §10)
   still waits for.
9. **Every premise gets a §1.2 row**: any sentence saying "the code currently does X" cites the command and its
   observed output; an empty result needs a known-positive control in the same call. Where the design and HEAD
   disagree, measure, route minimally, record it in §1.4, and list every file you touch outside the design's §8.
10. Close with this plan's own Appendix C: anything moved to sprint 2 (with size), and the sprint-2 handoff (K0's
    slot in the handler order, N4).

## Hard constraints
- D-WORLD-109: test keys, temp stores, `approvaltest` keys only. No real key, no `~/.ssh`, no `~/.ailang/world`, no
  real `/dev/tty`, no real browser profile, no mint or publish against a real store, no real network listener beyond
  `httptest`/loopback in tests. Never run bare `ssh-keygen`.
- `host/approval` stays pure and standard-library-only; the law in `host/broker/decision_law.go` is pure (no I/O, no
  clock reads — `now` is an argument). Effects only at the store/CLI/daemon boundary. The daemon never KEEPS the
  credentials writer. No CORS headers anywhere. `decidedBy` is derived by the daemon, never accepted from a caller.
- Do not touch: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/` (except reading),
  `cmd/world-publish/*` (the publish CLI; the L2 change is in the broker consumer), `~/.ailang/state/mission-v1*`,
  the `/surface/passkey` route, `passkey.js`, `approve.js`, the A2/A3a law, mints, `surface_approvals`.
- Every existing test the design lists as "Must stay green, unedited" stays unedited.
- Export `AILANG_BIN=$HOME/.pinned-ailang/ailang` in every shell that runs Go tests.
- Your Bash tool's shell is zsh: run loops/globs/`$VAR` work under `bash -c` or `bash <<'EOF'`. macOS has no
  `timeout` command.
- Known load-sensitive / flaky tests (re-run once alone, do not chase): `TestSeedRestartIgnoresASurvivingOldCopy`
  (host/broker), `TestHostsLaunchersExecHost`, host/broker exec-probe timeouts incl.
  `TestExecProtectedCheckTimeout/reserve`, `TestPublishApprovalRefusalSetWithALandedPositiveControl` (host/broker),
  local `-race` of host/broker, host/daemon, host/verifygate (incl. `TestSeToolsGateArmsAreLoadBearing/plan-bytes`),
  `TestMCPPostBudgetProductionConstants/batch` (host/daemon), `TestBackoff` (host/surfacefollow),
  `TestDemoSeedPreflightRefusals`, `TestWorktreeSnapshotRefusesAConfiguredFilter`,
  `TestMissionConfigDefaultFlowIntegration/invalid_ledger`, `scripts/marketingcheck` supervisor tests,
  `TestMarketingCaptureFakeChromeSixShots/inventory`.
- **Time budget: aim to finish within ~100 minutes; hard stop at 150.** The iteration has a 6 h wall clock that
  started 16:25Z, and the executor, gates, judge, CI and record still follow you. If you run short, ship the
  priced fallback cut with M1–M3 minus CD1–CD4 fully proven rather than everything half-proven.

## Waiting discipline (binding on you)
If you background anything (a long test run, a mutant loop), poll its ARTIFACT (an output/rc file it writes on
completion) inside a bounded `date +%s` loop in the same turn. A Monitor is an event stream, not a wait; nothing will
re-invoke you when it fires. Never end your turn announcing a wait, a pause, or "once it completes" — end only when
both files are written and the round trip has passed.

## Your final message
Report: the plan path, scope (what is in, anything moved and why, with measured numbers), the sub-milestone list with
expected numstat after each, mutant counts (total / killed / named equivalents) and the ledger path, the round-trip
result, base-gate results, every premise found false (§1.4), the N4/K0 handler-order decision, and anything the
controller must decide.
