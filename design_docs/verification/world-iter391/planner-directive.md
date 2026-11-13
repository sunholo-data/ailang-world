MISSION-ROLE: planner

You are the sprint-planner for Ailang World mission iteration 391 (unattended run; no human is present).
Write the sprint plan for **row 201 `w-world-approvals-by-key`, milestone A4, SPRINT 3 — the passkey arm and
CD6–CD7 moved from sprint 2, then "the surface"**. No re-design: `design_docs/planned/w-world-approvals-by-key-a4.md`
r2 is reviewed (iteration 388) and is NOT reopened. Sprint 1 LANDED in iteration 389 (PR #563, merge `338ff821`);
sprint 2 LANDED in iteration 390 (PR #566, merge `9ea7b45f`; judge sonnet PASS 88). D-WORLD-115 is still OPEN; its
default **A = keep it** is applied. If sprint 3 depends on the D-WORLD-115 answer, record it in §1.4 and build under A.

## Scope — read these verbatim first
- Design §0 "Sprint 3", §1.12 (planner obligation P1), §2.3, §2.7 (passkey halves of the key rules), §2.8 (routes,
  handler order, `requests.json` `decisions`, `GET /surface/decide/{ref}` + `decide.js`), §2.9/§2.10 (app.js
  *Needs you*, CLI incl. `approve --list` Q2), §6 "Sprint 3" table, §7, §8, §9 "Sprint 3 = M6–M7" incl. its
  **fallback cut**, §10 (attended exit run).
- **Sprint-2 plan Appendix C** (`design_docs/planned/w-world-approvals-by-key-a4-sprint2-plan.md`, line ~984) is
  the handoff. **The head of this sprint is, in order:**
  1. **Planner obligation P1 first**: measure A4's 91-byte decision message with `crypto.subtle` in Chrome against
     Go's goldens, as A3 V23 did. Use a temp Chrome profile only.
  2. **The webauthn arm** (design price ≈ 500): `DecideByAssertion` (DM7, the A2 verifier's defer-W15 form),
     `JudgeDecisionAssertion`, the POST's `webauthn` grammar member and arm (W15 + `usePasskey` in
     `RecordKeyDecision`'s transaction on the non-exact path; the exact-match path skips only the counter), KD2,
     the passkey halves of KD7, KD8, RD13 ((b) counter N, (c) counter 0) and RD14 (across a restart),
     `DecisionSubject.Passkeys` (sprint-1 P7), CD8 `--webauthn-link`, CD9 `--webauthn-assertion` on a request ref.
     `RecordKeyDecision`'s signature carries the passkey use through `ApproverUse.Passkey` (already a field).
  3. **CD6–CD7** (`--write-message` / `--signature --expiry` on a request ref): built and green in iteration 390's
     planner scratch, 63 changed lines measured, tagged `CD67-CUT` (`afbc129`) in
     `$HOME/.ailang/state/world-iter390/plan-scratch/work` (a git repo; read it with `git -C … show CD67-CUT`).
     Its sign seam predates `signDecision`, so rebase it on sprint 2's M5c as landed. Treat it as a claim:
     re-derive and re-test it, do not paste it.
  4. **Then the surface (M6–M7)**: `requests.json` gains `decisions` (DMD1–DMD3), the decision card
     `GET /surface/decide/{ref}` + `decide.js` (PD1, PD2), *Needs you* in `app.js`, live Chrome XD1–XD4 (gated as
     XB1, as A3), QUICKSTART §18, design Q2 `approve --list` decisions.
- **Sprint-2 judge's non-blocking test gaps, to close in this sprint (small, test-only unless a test shows a bug):**
  - **JM22**: `LandKeyDecision`'s `request_ref` filter is unpinned — controller-reproduced in 390: mutating
    `WHERE request_ref = ? AND decision_ref IS NULL` → `WHERE decision_ref IS NULL` in
    `host/store/key_decisions.go` (~line 143) leaves `go test ./host/store/ -run KeyDecision` rc 0.
  - **K0's approver and format operands** are unpinned at the store API (equivalent through the daemon).
  - **JM01–JM04**: the four `isApproverRace` operands' mapping to 422 is unpinned.
  - **JM27**: `LastUsedAt` pass-through is unpinned.
  - **JM46**: the CLI's default daemon address is unpinned.
  Evidence: `$HOME/.ailang/state/world-iter390/judge/` and `$HOME/.ailang/state/world-iter390/evaluator-r1.md`.
  Also the sprint-2 plan's named survivors `D-RACE`, `D-LAND-ERR`, `C-PROOF-ERR`, `C-SIGN-NIL` (its §4): close any
  whose guarded code you change.
- **Size and the fallback**: ≈ 1,500 (design) + ≈ 500 (webauthn) + 63 ≈ 2,060, over the cap. Expect your own cut.
  The design's priced fallback is: **the card page (M6 page half, XD1, XD2, XD4) becomes a fourth sprint**; the
  `decisions` listing, the *Needs you* line (pointing at `ailang-worldd approve <ref>`) and QUICKSTART stay; PD1 is
  never separated from the page. **Resolve the interaction honestly**: the design says the webauthn arm is needed
  for XD1 — if a passkey decision cannot be produced by an operator without the card page, say so and choose the
  cut that leaves every shipped arm reachable end to end by an operator (e.g. CD8/CD9 vs the card). Never ship a
  POST arm with no operator path to it unless you name it as such in §3.1 and §9. **Never cut**: the J-gap tests,
  P1's measurement, KD2, RD13 (b)/(c). Decide from the MEASURED size and record the decision (S3-D1).

## Where you work
- Sprint worktree (already created, branch `sprint/world-iter391-row201-a4-s3`, base = origin/dev
  `7004f5913d85cc1f307068a7f4c4e4e31aa190ea`, which contains A4 sprints 1 and 2):
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter391-a4-s3`
- VERIFIED BY CONTROLLER at the base (commands run in the worktree):
  - `host/store/credentials.go:44` `const credentialsSchemaVersion = 8`.
  - `host/sshsign/sshsign.go:50` `func SignIn(ctx, bin, namespace, key, msg, stderr)` exists (sprint 2 P1).
  - `grep -n -i webauthn host/daemon/surface_decide.go` → no hits (sprint 2 refuses the member as an unknown field).
  - `host/daemon/daemon.go:1138` routes `GET /surface/requests.json`; `host/surfacepage/app/` holds `app.js`,
    `app.css`, `index.html`; `docs/QUICKSTART.md` ends its numbered sections at §17 (line 1841).
  - Pinned ailang: `$HOME/.pinned-ailang/ailang` = AILANG v0.41.0.
  - Rig load average ≈ 10 on 16 cores. Other sessions run; expect some timing-sensitive reds in untouched packages;
    re-run once alone before believing a red.
- Read the landed sprint-2 plan for the plan FORMAT and reproduce it (§0 binding constraints, §1.2 premise table,
  §1.3 decisions, §1.4 false premises, §2 patch protocol + awk extraction, §3 sub-milestones with expected numstat,
  §3.1 reachability, §4 mutation table, §5 files, §6 gate script, §7 executor notes, §8 risks, §9 acceptance
  criteria, Appendix C, Appendix P fences). Sub-milestone ids continue: M6a, M6b, … M7x.
- The design's line numbers were measured at `b0fbcd3e`; re-measure every line/premise you rely on at `7004f591`
  and record drift in §1.2/§1.4.
- Scratch (compile your stages, run mutants, Chrome profiles) OUTSIDE the repo:
  `$HOME/.ailang/state/world-iter391/plan-scratch/` (never under /tmp — location-sensitive tests go red there).
- Write exactly two files in the worktree, and DO NOT git add/commit them (the controller commits):
  - `design_docs/planned/w-world-approvals-by-key-a4-sprint3-plan.md`
  - `design_docs/planned/sprint_w-world-approvals-by-key-a4-s3.json` (same shape as `sprint_w-world-approvals-by-key-a4-s2.json`)
- Leave every code/doc file in the worktree untouched (the executor applies your patches there). No git writes of any
  kind in the worktree.

## Read first
1. The design (sections above), the sprint-2 plan (format + Appendix C + §4 survivors), and the sprint-1/2 code you
   build on: `host/approval/decision*.go`, `host/broker/decision_law.go`, `host/daemon/surface_decide.go` and tests,
   `host/store/key_decisions.go`, `host/authority/` (the decision judge), `cmd/ailang-worldd/approve_decide.go`.
2. The A2/A3 passkey machinery you reuse: `host/approval/` (`ApproveByAssertion`, the assertion verifier, W15),
   `host/store/credentials.go` (`usePasskey`, `useApprover`, the writer and its deadline rule),
   `cmd/ailang-worldd/approve*.go` (`--webauthn-link`, `--webauthn-assertion`, `--write-message`, `--signature`
   in mint mode), the A3a surface (`host/daemon/surface_approve.go`, `host/surfacepage/app/app.js`, the
   `/surface/passkey` page and its live-Chrome tests, `TestPasskeyPage`, `TestApprovePage`), `host/approvaltest`.
3. `CLAUDE.md`, `design_docs/coding-standards.md` (binding: pure core, effects at the boundary, non-vacuous gates).

## Size rule
- **Cap: ≤ 2,000 changed lines**, measured `git diff --numstat <base> HEAD -- cmd host docs` (state this exact
  filter; include docs because QUICKSTART §18 is in scope — if you prefer `-- cmd host` plus a separate docs count,
  state both). Prefer a smaller, fully proven sprint over a full, half-proven one.
- Sub-milestones of ≤ ~250–350 changed lines each, each independently green and safe at commit.
- No new `/v1` route (`TestFrozenV1RouteTableMatchesMux` stays 13). New `/surface` routes only as the design names
  them (`TestSurfaceAppRoutes` +2 in sprint 3 per design §7). `surfaceCSP` unchanged unless the design says
  otherwise. No CORS headers.

## Retro directives (iterations 374–390; each was paid for)
1. **Ship compiled patches**: every line of code, test and doc as a unified diff in Appendix P
   (```` ````diff <M> ```` fences), generated LAST from compiled, vetted, gofmt-clean, green scratch stages; then
   extract them FROM THE PLAN FILE with §2's awk, apply to a fresh `git archive 7004f591` export, and `diff -r`
   against your work tree — byte-identical or regenerate. The controller re-does this round trip independently.
   **Ship the patches BEFORE running the full mutant set, and checkpoint the plan file** (write a first complete
   version, then append the mutant ledger).
2. **Expected numstat after each sub-milestone's commit** (cumulative, measured after commit, with the §size filter).
3. **Mutants anchored to the DIFF**: one mutant per new production branch, enumerated from the diff, not from the
   test plan — in particular DM7's steps, the assertion judge's steps, the POST format switch's webauthn arm, W15
   and the counter (exact-match path skips only the counter), KD2, `usePasskey` in the transaction, the
   `decisions` listing's filters (DMD1–DMD3), the card page's refusals (PD1/PD2), `decide.js`'s message bytes,
   *Needs you*'s escaping, CD6–CD9 flag rules and the 91-byte message written by `--write-message`. Run each
   against the named test row AND the package; record killed/survived; name equivalents with reasons. A mutant that
   does not compile is not a mutant (pair with a killed control). **Name the test that killed each mutant; a mutant
   killed only by an inventory/static test is not behavioural proof — re-run with that test excluded and record
   both.** Write the ledger to a file in scratch; every re-run appends its own line. **The plan's mutant table ids
   must be the same ids as in the ledger** (the judge will `comm` them).
4. **New branch reachability**: for every new refusal/branch, give a trace from the request's entry point (the POST
   handler, the GET routes, the CLI, the page script) to that branch for each input class, and say which test drives
   it end to end vs unit-only. A branch unreachable end to end because an earlier step refuses first must say so.
5. **Repo-wide inventories**: run every design §7 inventory by name with `-v` on your final scratch tree (incl.
   `TestFrozenV1RouteTableMatchesMux`, `TestSurfaceAppRoutes`, `TestDaemonDependencyAllowlist`,
   `TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots`, `TestProductionGoSurface`,
   `TestEveryVerbHasHelp`, the quickstart tests, `TestDaemonCredentialWriterSites`,
   `TestCredentialsMethodsRequireADeadline`, `TestSurfaceAppPure`, any embedded-asset/static-file inventory for
   `host/surfacepage`, `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`, import allowlists for
   `host/approval`). Search for any test that enumerates routes, exported symbols, files, embedded assets, CLI verbs,
   schema versions or QUICKSTART sections before assuming none.
6. **A directive whose diff adds a production context root must name
   `go test ./host/store/ -run '^TestProductionContextRoots$'` as a gate** (charter process fix, iter-292).
7. The §6 gate script for the controller (one bash shell; **every leg records its own rc and the script runs every
   leg — no `set -e`**; `tee` every leg to `$HOME/.ailang/state/world-iter391/gates/` (mkdir it first), never only
   through `tail`), incl. `-race` on touched Go packages, the inventories, the live-Chrome leg (with its env, copied
   verbatim from the design/plan, and a temp profile), `./scripts/verify_ail.sh`, and `VERIFY_GO_LEGS=plain
   AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh` after commit. Measure each gate at the base first.
   **Never start a gate under `nohup`.** The script must run under macOS `/bin/bash` 3.2: **no `declare -A`**.
8. Executor notes as the sprint-2 plan's §7: no git writes; `patch -p1` is the only writer; runs of ~2–4
   sub-milestones; cumulative snapshots to `$HOME/.ailang/state/world-iter391/snap/<M>/` with MANIFEST.sha256
   (exclude the two plan files); stop conditions; numstat checks; exact acceptance commands per sub-milestone (copy
   any env such as `WORLD_CHROME`, `WORLD_REQUIRE_PASSKEY_BROWSER`, `WORLD_EXEC_NODE` verbatim into each command);
   **an explicit "never `nohup`, never background a test leg" line; and "every Go Bash call passes
   `timeout: 600000`, foreground, no Monitor"**.
9. Acceptance criteria (§9-style), numbered, that the judge can score against. Include what an operator can do after
   this sprint and what the attended exit run (design §10) still waits for.
10. **Every premise gets a §1.2 row**: any sentence saying "the code currently does X" cites the command and its
   observed output; an empty result needs a known-positive control in the same call. Where the design and HEAD
   disagree, measure, route minimally, record it in §1.4, and list every file you touch outside the design's §8.
11. **A test row that drives a fake tool must have a sibling row fed the real tool's captured response** (charter
   process fix, iter-299; it caught sprint 2's P1) — relevant to Chrome's `crypto.subtle`/WebAuthn output (P1's
   goldens) and to `ssh-keygen -Y sign` in CD6–CD7: use an `approvaltest` key file / virtual authenticator, and
   check in any captured real output under `testdata/` with its generating command in a comment.
12. Close with this plan's own Appendix C: anything moved to a fourth sprint (with size), and its handoff.

## Hard constraints
- D-WORLD-109: test keys, temp stores, `approvaltest` keys only. No real key, no `~/.ssh`, no `~/.ailang/world`, no
  real `/dev/tty`, no real browser profile (temp `--user-data-dir` only), no mint or publish against a real store,
  no real network listener beyond `httptest`/loopback in tests. A test may run `ssh-keygen -Y sign` only on an
  `approvaltest` key file in a temp dir.
- `host/approval` stays pure and standard-library-only; the laws are pure (no I/O, no clock reads — `now` is an
  argument). Effects only at the store/CLI/daemon boundary. The daemon never KEEPS the credentials writer. No CORS
  headers anywhere. `decidedBy` is derived by the daemon, never accepted from a caller.
- Do not touch: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/` (except reading),
  `cmd/world-publish/*`, `~/.ailang/state/mission-v1*`, the A2/A3 mint law and its matrix tests, mints,
  `surface_approvals` — and anything the design's §8 marks unchanged. If the design's sprint 3 needs an edit to an
  A3a file (`approve.js`, `passkey.js`, the `/surface/passkey` route), it must be named in the design; otherwise
  record the need in §1.4 and route minimally.
- Every existing test the design lists as "Must stay green, unedited" stays unedited.
- Export `AILANG_BIN=$HOME/.pinned-ailang/ailang` in every shell that runs Go tests.
- **Every Bash call that runs Go builds or tests passes `timeout: 600000`** (the tool default is 120 s). Run long
  loops in the foreground in chunks under 10 min, or background with an rc-file you poll in a bounded loop.
- Your Bash tool's shell is zsh: run loops/globs/`$VAR` work under `bash -c` or `bash <<'EOF'`. macOS has no
  `timeout` command.
- Known load-sensitive / flaky tests (re-run once alone, do not chase): `TestSeedRestartIgnoresASurvivingOldCopy`
  (host/broker), `TestHostsLaunchersExecHost`, host/broker exec-probe timeouts incl.
  `TestExecProtectedCheckTimeout/reserve-after-before`, `TestPublishApprovalRefusalSetWithALandedPositiveControl`,
  local `-race` of host/broker, host/daemon, host/verifygate, `TestMCPPostBudgetProductionConstants/batch`,
  `TestBackoff` (host/surfacefollow), `TestDemoSeedPreflightRefusals`, `TestWorktreeSnapshotRefusesAConfiguredFilter`,
  `TestCommitInvocationReceiptIdentity` (race), `TestMissionConfigDefaultFlowIntegration/invalid_ledger`,
  `scripts/marketingcheck` supervisor tests, `TestMarketingCaptureInventoryMutation`, `TestSessionSelfGit`,
  `TestServeStartsWhenTheProbeSpendsTheStartupClock` (red at base), `TestRegistryPhaseHasItsOwnClock`,
  `TestWorkspaceProtectedPolicyWiring`, `TestConvertedStoreMethodsM5b`, host/pkgproj under load.
- **Time budget: aim to finish within ~100 minutes; hard stop at 130.** The iteration has a 6 h wall clock that
  started 01:31Z; the executor, gates, judge, CI and record still follow you. If you run short, ship a smaller cut
  fully proven (the J-gap tests + P1 + the webauthn arm + CD6–CD7, with the surface priced into Appendix C as a
  fourth sprint) rather than everything half-proven.

## Waiting discipline (binding on you)
If you background anything (a long test run, a mutant loop, Chrome), poll its ARTIFACT (an output/rc file it writes
on completion) inside a bounded `date +%s` loop in the same turn. A Monitor is an event stream, not a wait; nothing
will re-invoke you when it fires. Never end your turn announcing a wait, a pause, or "once it completes" — end only
when both files are written and the round trip has passed.

## Your final message
Report: the plan path, scope (what is in, anything moved and why, with measured numbers), the sub-milestone list with
expected cumulative numstat after each, mutant counts (total / killed / named equivalents / survivors) and the ledger
path, the round-trip result, base-gate results, P1's result, every premise found false (§1.4), how an operator
produces a passkey decision after this sprint, and anything the controller must decide.
