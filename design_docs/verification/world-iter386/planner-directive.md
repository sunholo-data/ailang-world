MISSION-ROLE: planner

You are the sprint-planner for Ailang World mission iteration 386 (unattended run; no human is present).
Write the sprint plan for **row 201 `w-world-approvals-by-key`, milestone A3a, SPRINT 3** — "the page and the app":
M6 the card page + `approve.js` (PG2, PG3, PG3b, the card's network deadlines), M7 *Needs you*, Decided and the
opt-in notification in `app.js`, M8 the live browser tests XA1–XA5 on the A2 Chrome driver (extended) and
QUICKSTART §17 *(attended — pending first verbatim run)* — exactly as the A3 child design r2 §2.7–§2.10, §6
"Sprint 3" and §9 list them, **plus two carried obligations** (below).
No re-design: `design_docs/planned/w-world-approvals-by-key-a3.md` r2 is reviewed (quorum r1+r2 + controller
carve-out + independent audit VALID 84) and is NOT reopened. D-WORLD-114 (the doc's two amendments to the
attended parent) is still OPEN with default A = keep both, applied since iteration 384; build as designed.

## Carried obligations (both MUST be in this sprint unless the size rule below moves them, with numbers)
1. **Appendix C of the sprint-2 plan** (`design_docs/planned/w-world-approvals-by-key-a3a-sprint2-plan.md`
   §"Appendix C"): R11's `sshsig` arm of `POST /surface/approve/{id}` (`authority.JudgeSSHSignature` on
   `*body.SSHSIG`, recorded as `ssh`/`ssh:<name>`; it REPLACES sprint 2's refusal at
   `host/daemon/surface_approve.go:255` `if body.SSHSIG != nil { // R11, sprint 2's fallback cut …` and its RT4
   arm `R11-sshsig-cut`), **RT6** `TestApprovePOSTRecordsAVerifiedSSHSIG`, **RT13**
   `TestApprovePOSTSSHSIGViaTailnetOriginWithNoPasskey`, and **AW6's route half** (an SSHSIG POSTed, then
   `approve --await` mints `ssh:<name>`, as E2E1 does for the passkey). Design §2.4/§6 are the spec.
2. **Judge survivor J36 from iteration 385** (controller-reproduced): in `handleSurfaceRequestCard`
   (`host/daemon`), replacing the `ActivePasskeys` error check with `if false` leaves every test green. Add a
   test that injects an `ActivePasskeys` error on the card route and asserts the design's answer (a 500 with the
   `Internal` refusal shape, not a card with no passkeys) — read design §2.3 for the exact status — and PROVE it
   kills that mutation (apply it, test red; revert, green). Put it in the earliest sub-milestone where it fits.

## Where you work
- Sprint worktree (already created, branch `sprint/world-iter386-row201-a3a-s3`, base = origin/dev
  `82fb9afa9086eed2f781bb0d4fe54d4bae3b70d7`, which contains A3a sprints 1 and 2 merged (`b44b5f7d`, `d1027477`)):
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter386-a3a-s3`
- VERIFIED BY CONTROLLER at the base: production Go has 0 hits for `approve.js` and 0 for a
  `/surface/approve/` page route; the sprint-2 SSHSIG refusal is at `host/daemon/surface_approve.go:255`
  (controls: `/surface/passkey` = 11 prod hits, `requests.json` = 5). `host/surfacepage/app/` holds
  `app.js`, `app.css`, `index.html`; `host/surfacepage/passkey/` holds `index.html`, `passkey.js`.
  Chrome is present at `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`.
- The design's line numbers were measured at `0140bd7f`; sprints 1–2 have since changed `cmd/` and `host/`.
  Re-measure every design line/premise you rely on at `82fb9afa` and record drift in §1.2/§1.4. Read the landed
  sprint-2 plan §1.4 (P1–P7) and Appendix C — its premises are now facts of HEAD you build on.
- Scratch (compile your stages, run mutants, Chrome profiles) OUTSIDE the repo:
  `$HOME/.ailang/state/world-iter386/plan-scratch/` (never under /tmp — location-sensitive tests go red there).
- Write exactly two files in the worktree, and DO NOT git add/commit them (the controller commits):
  - `design_docs/planned/w-world-approvals-by-key-a3a-sprint3-plan.md`
  - `design_docs/planned/sprint_w-world-approvals-by-key-a3a-s3.json` (same shape as `sprint_w-world-approvals-by-key-a3a-s2.json`)
- Leave every code/doc file in the worktree untouched (the executor applies your patches there). No git writes of any kind in the worktree.

## Read first
1. `design_docs/planned/w-world-approvals-by-key-a3.md` r2 in full, especially §0, §2.3, §2.4 (the SSHSIG arm),
   §2.7 (the card page, `FIELDS`, digest check, network deadlines/retries [r2]), §2.8, §2.9, §2.10, §3, §4, §5,
   §6 "Sprint 3" (PG2, PG3, PG3b, XA1–XA5, the driver edits `fetchRewrite`/`grantNotifications`/`click`/
   `notify`/`getRecorder`), §7 inventories, §8 conflict surface, §9 sprint-3 gates and the **fallback cut**, §11,
   §14 V-rows (V24–V28 are the browser measurements you rely on — re-measure what you can).
2. `design_docs/planned/w-world-approvals-by-key-a3a-sprint2-plan.md` for the plan FORMAT (newest; it worked:
   §1.2 premise table, §1.4 false premises, §2 patch protocol + awk extraction, §6 gate script, §7 executor
   notes, Appendix P fences, expected numstat per sub-milestone, Appendix C). Reproduce that format.
3. The code you build on: `host/daemon` (`surface_approve.go`, the read routes, `approvalRoute`, `isProtected`,
   `surfaceCSP`, `TestSurfaceAppRoutes`, RT1's `approvalRoutes`, RT7's targets), `host/surfacepage`
   (`approve_model.go`, `app/app.js`, `passkey/passkey.js`, the embed files), the A2 Chrome harness in
   `cmd/ailang-worldd` (`TestPasskeyPageInChrome`, `passkey_cdp.mjs` or wherever the CDP driver lives),
   `host/authority/judge.go`, `cmd/ailang-worldd/approve_await.go`, `host/approval/approvaltest`, `docs/QUICKSTART.md`
   (§16 and the quickstart tests that pin it).
4. `CLAUDE.md`, `design_docs/coding-standards.md` (binding: pure core, effects at the boundary, non-vacuous gates).

## Size rule (decide with measured numbers)
- **Cap: ≤ 2,000 changed lines**, measured `git diff --numstat <base> HEAD -- cmd host` (state this exact
  filter; `.js`/`.html`/`.mjs` under `host/`/`cmd/` count; `docs/QUICKSTART.md` does not, but list its lines).
- The design estimates ≈ 2,200 for sprint 3 BEFORE the carried Appendix C, so expect the design's own fallback
  cut: **M7's notification and Decided list (XA4, half of XA3) move to a fourth sprint**; the card, Approve and the
  Needs-you link (XA1, XA2, XA3's list half) stay. **Never cut:** PG2, PG3, PG3b, XA1, XA2 incl. arm (c), XA5,
  the J36 test.
- If the plan still exceeds 2,000 after that cut, the next move is Appendix C's SSHSIG arm (RT6, RT13, AW6's
  route half — the 422 refusal stays in place, so nothing is silently accepted) to sprint 4 as well. Record every
  moved item in this plan's Appendix C as sprint-4 scope (no re-design), with the measured line counts that forced
  it. Prefer a smaller, fully-proven sprint over a larger one.
- Split into sub-milestones of ≤ ~250–350 changed lines each where possible (M6a/M6b/…/M8a/…), each
  independently green and safe at commit. **Safety ordering is the design's:** routes register through
  `approvalRoute` (the Host check); `surfaceCSP` unchanged; the page verifies the digest itself and never uses a
  server hash; the SSHSIG arm judges before it records, exactly as the passkey arm does.

## Retro directives (iterations 374–385; each was paid for)
1. **Ship compiled patches**: every line of code, test, JS and doc as a unified diff in Appendix P
   (```` ````diff <M> ```` fences), generated LAST from compiled, vetted, gofmt-clean, green scratch stages; then
   extract them FROM THE PLAN FILE with §2's awk, apply to a fresh `git archive 82fb9afa` export, and `diff -r`
   against your work tree — byte-identical or regenerate. The controller re-does this round trip independently.
   **Ship the patches BEFORE running the full mutant set, and checkpoint the plan file** (write a first complete
   version, then append the mutant ledger) — 384 took 69 min and 385 77 min this way.
2. **Expected numstat after each sub-milestone's commit** (measured after commit, `-- cmd host`).
3. **Mutants anchored to the DIFF**: one mutant per new production branch (Go AND JS: each route registration,
   the digest check, the FIELDS walk, unknown-field refusal, each deadline/AbortController, the retry predicate
   (503 only, ≤3, 1 s), unknown-outcome messaging, the SSHSIG arm's judge-before-record and each of its refusal
   arms, the Needs-you list wiring, the J36 error check), enumerated from the diff, not from the test plan; run
   each against the named test row AND the package; record killed/survived; name equivalents with reasons; base
   controls where a test pins pre-existing behaviour. A mutant that does not compile is not a mutant (use
   `-vet=off` and a compiling edit; pair with a killed control). **When a mutant is killed, name the test that
   killed it; a mutant killed only by a call-site INVENTORY (e.g. `TestDaemonCredentialWriterSites`) under
   `-failfast` is not behavioural proof — re-run excluding inventories and record both** (iteration 385's judge
   lesson). Write the ledger to a file in scratch; every re-run appends its own line.
4. **Live Chrome**: the XA tests run on this machine (`WORLD_CHROME=... WORLD_REQUIRE_PASSKEY_BROWSER=1`, design
   §9). Run them yourself on your final scratch tree and record `--- PASS` per subtest; CI has no Chrome step, so
   state in the plan which XA tests skip in CI and what (if anything) in CI still pins the page (PG2/PG3/PG3b).
   Chrome profiles go in scratch, temporary, headless; never a real profile.
5. **Repo-wide inventories**: run every design §7 inventory by name with `-v` on your final scratch tree (incl.
   `TestFrozenV1RouteTableMatchesMux`, `TestSurfaceAppRoutes`, `TestDaemonDependencyAllowlist`,
   `TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots`, `TestProductionGoSurface`,
   `TestNoNewDeadlineFreeStoreReads`, `TestCredentialsMethodsRequireADeadline`, `TestEveryVerbHasHelp`, the
   quickstart tests, `TestDaemonCredentialWriterSites`, `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`,
   `TestApprovalImportsOnlyTheStandardLibrary`). Search for any test that enumerates daemon routes, embedded
   files, serve flags or QUICKSTART sections before assuming none.
6. The §6 gate script for the controller (one bash shell, `tee` every leg to `$HOME/.ailang/state/world-iter386/gates/`,
   never only through `tail`), incl. `-race` on touched packages, the inventories, the live-Chrome leg,
   `./scripts/verify_ail.sh`, and `VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh`
   after commit. Measure each gate at the base first. **Never start a gate under `nohup`**.
7. Executor notes as the sprint-2 plan's §7: no git writes; `patch -p1` is the only writer; runs of ~2–4
   sub-milestones; cumulative snapshots to `$HOME/.ailang/state/world-iter386/snap/<M>/` with MANIFEST.sha256
   (exclude the two plan files); stop conditions; numstat checks; exact acceptance commands per sub-milestone;
   **an explicit "never `nohup`, never background a test leg" line** (385's executor broke it once).
8. Acceptance criteria (§9-style), numbered, that the judge can score against.
9. **Every premise gets a §1.2 row**: any sentence saying "the code currently does X" cites the command and its
   observed output; an empty result needs a known-positive control in the same call. Where the design and HEAD
   disagree, measure, route minimally, record it in §1.4, and list every file you touch outside the design's §8.

## Hard constraints
- D-WORLD-109: test keys, temp stores, `approvaltest` passkeys/keys and virtual authenticators only. No real key,
  no `~/.ssh`, no `~/.ailang/world`, no real `/dev/tty`, no real browser profile, no mint or publish against a
  real store, no real network listener beyond `httptest`/loopback in tests. Never run bare `ssh-keygen`.
- `host/approval` stays pure and standard-library-only. Effects only at the store/CLI/daemon boundary.
  `surfaceCSP` is not changed. The daemon never KEEPS the credentials writer. No CORS headers anywhere.
  `approve.js`: no `innerHTML`, `eval(`, `document.cookie`.
- Do not touch: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/` (except reading),
  `cmd/world-publish/*`, `~/.ailang/state/mission-v1*`, `/surface/passkey` route, `passkey.js`, the A2 law.
- Export `AILANG_BIN=$HOME/.pinned-ailang/ailang` in every shell that runs Go tests.
- Your Bash tool's shell is zsh: run loops/globs/`$VAR` work under `bash -c` or `bash <<'EOF'`.
- Known load-sensitive / flaky tests (re-run once alone, do not chase): `TestSeedRestartIgnoresASurvivingOldCopy`
  (host/broker), `TestHostsLaunchersExecHost`, host/broker exec-probe timeouts incl.
  `TestExecProtectedCheckTimeout/reserve`, local `-race` of host/broker, host/daemon, host/verifygate (incl.
  `TestSeToolsGateArmsAreLoadBearing/plan-bytes`), `TestDemoSeedPreflightRefusals`,
  `TestWorktreeSnapshotRefusesAConfiguredFilter`, `TestMissionConfigDefaultFlowIntegration/invalid_ledger`,
  `scripts/marketingcheck` supervisor tests, `TestMarketingCaptureFakeChromeSixShots/inventory`.
- **Time budget: aim to finish within ~100 minutes; hard stop at 150.** The iteration has a 6 h wall clock that
  started 04:06Z, and the executor, gates, judge, CI and record still follow you.

## Waiting discipline (binding on you)
If you background anything (a long test run, a mutant loop, Chrome), poll its ARTIFACT (an output/rc file it
writes on completion) inside a bounded `date +%s` loop in the same turn. A Monitor is an event stream, not a wait;
nothing will re-invoke you when it fires. Never end your turn announcing a wait, a pause, or "once it completes" —
end only when both files are written and the round trip has passed.

## Your final message
Report: the plan path, scope (what is in, anything moved by the fallback cut and why, with measured numbers), the
sub-milestone list with expected numstat after each, mutant counts (total / killed / named equivalents) and the
ledger path, the J36 kill proof, the live-Chrome results per subtest, the round-trip result, base-gate results,
every premise found false (§1.4), and anything the controller must decide.
