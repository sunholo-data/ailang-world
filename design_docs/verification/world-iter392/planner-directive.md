MISSION-ROLE: planner

You are the sprint-planner for Ailang World mission iteration 392 (unattended run; no human is present).
Write the sprint plan for **row 201 `w-world-approvals-by-key`, milestone A4, SPRINT 4 — "the surface"**, the
part the sprint-3 planner moved out (its decision S3-D1). No re-design: `design_docs/planned/w-world-approvals-by-key-a4.md`
r2 is reviewed (iteration 388) and is NOT reopened. Sprints 1–3 LANDED (iterations 389/390/391; PRs #563, #566,
#569; sprint 3 merge `e6761873`, judge sonnet PASS 91). D-WORLD-115 is still OPEN; its default **A = keep it** is
applied. If anything in sprint 4 depends on the D-WORLD-115 answer, record it in §1.4 and build under A.
This should be the LAST A4 sprint: after it, A4 is code-complete except Mark's attended exit run (design §10).

## Scope — read these verbatim first
- Design §0 "Sprint 3", §1.12 (P1), §2.8 (routes: `requests.json` `decisions`, `GET /surface/decide/{ref}` +
  `decide.js`, handler order), §2.9 (*Needs you*, what the agent observes), §2.10 (CLI incl. `approve --list`, Q2),
  §6 "Sprint 3: the surface" table (DMD1–DMD3, PD1, PD2, XD1–XD4 …), §7 inventories, §8 conflict surface,
  §9 "Sprint 3" gates, §10 (attended exit run), §12 Q2.
- **Sprint-3 plan Appendix C** (`design_docs/planned/w-world-approvals-by-key-a4-sprint3-plan.md`, line 867) is the
  handoff, priced at ≈ 1,250:
  1. `requests.json` gains `decisions` + `decisionsTruncated` (DMD1–DMD3; ≈ 120 Go + 200 tests).
  2. *Needs you* in `app.js` linking each pending decision to its card; the reserved text without a link; the
     notification for a new decision id (XD3; ≈ 80 JS + 120 tests).
  3. The decision card `GET /surface/decide/{ref}` + `GET /surface/decide.js` (`DECISION_FIELDS`, the hash and
     `requestRef` checks, the passkey ceremony with `challenge = DecisionMessageHash` built with **`BigInt`** for the
     expiry (sprint-3 P1), **Save the message** + the signing line + a paste box; PD1, PD2, XD1, XD2, XD4; ≈ 300 page
     + 60 Go + 300 tests). The card POSTs to sprint 3's `POST /surface/decide/{ref}` arms UNCHANGED (the `webauthn`
     body = the page's envelope; the `sshsig` body from the paste box). XD1 can follow
     `TestApproveDecideWebAuthnLinkInChrome`'s driver use; `decide.js` repeats A3's encoding helpers (design §2.8) and
     P1's construction (sprint-3 C4's probe page: `$HOME/.ailang/state/world-iter391/plan-scratch/p1/srv/main.go`).
  4. `DecisionSubject.Passkeys` (sprint-1 P7, ≈ 30) — VERIFIED BY CONTROLLER: `host/surfacepage/decide_model.go`
     `type DecisionSubject struct` has no `Passkeys` field at `799660ac`.
  5. `approve --list` shows pending decisions (design Q2, ≈ 80).
  6. `TestSurfaceAppRoutes` +2; QUICKSTART §18's card paragraph (≈ 20) — §18 exists (`docs/QUICKSTART.md:1861`,
     "Decide an agent's approval request").
  7. **The deny verb on screen**: the sprint-3 judge named that the A2 passkey page (`/surface/passkey`, reached by
     `approve <request-ref> --deny --webauthn-link`) always reads "Approve request …", even for a deny (the signature
     binds the verb; a mismatched post is 422). The card must show its own verb correctly. Decide, and record as a
     plan decision, whether the CLI's `--webauthn-link` path also gets correct verb text in this sprint (that would
     edit an A2/A3a file — `passkey.js` or the `/surface/passkey` route — which the design does not name: measure
     the size, route minimally, list the file in §1.4 and §5, and keep the A2 mint page's behaviour and its tests
     unchanged for mints) or is named as a residual with its reason.
- **Sprint-3 named survivors carried** (its §4): `D-LAND-ERR`, `C-PROOF-ERR`, `C-SIGN-NIL`, `C-LINK-ERR`,
  `KA-RACE-TEXT` — close any whose guarded code you change; otherwise keep them named.
- **Sprint-3 judge's non-blocking notes** (evidence: `$HOME/.ailang/state/world-iter391/evaluator-r1.md`,
  `$HOME/.ailang/state/world-iter391/judge/`): `hErr` placement unpinned; the packet the CLI posts with
  `--signature`/`--webauthn-assertion` is not asserted; the two file-size read caps unpinned. Small test-only
  closes are welcome if they fit; not required.
- **Size**: ≈ 1,250 priced. Cap 2,000. If your measured size exceeds the cap, cut the least operator-visible part
  and say why (S4-D1) — but the design's rule stands: **PD1 is never separated from the page**, and no POST arm or
  page ships without an operator path to it.

## Where you work
- Sprint worktree (already created, branch `sprint/world-iter392-row201-a4-s4`, base = origin/dev
  `799660ac4521919e2efe9c864fa8f5a0d7666ac8`, which contains A4 sprints 1–3 and the 391 record):
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter392-a4-s4`
- VERIFIED BY CONTROLLER at the base (commands run in the worktree):
  - `host/daemon/daemon.go:1138` `GET /surface/requests.json`; `:1143` `GET /surface/decisions/{ref}`;
    `:1144` `POST /surface/decide/{ref}`. `func TestSurfaceAppRoutes` at `host/daemon/surface_app_test.go:23`.
  - `host/surfacepage/` holds `app/` (`app.js`, `app.css`, `index.html`), `approve/`, `passkey/`, `decide_model.go`,
    `surface.js`.
  - `host/store/credentials.go:44` `const credentialsSchemaVersion = 8`.
  - Pinned ailang: `$HOME/.pinned-ailang/ailang` = AILANG v0.41.0.
  - Rig load average ≈ 10–11 on 16 cores. Other sessions run; expect some timing-sensitive reds in untouched
    packages; re-run once alone before believing a red.
- Read the landed sprint-3 plan for the plan FORMAT and reproduce it (§0 binding constraints, §1.1 baseline, §1.2
  premise table, §1.3 decisions, §1.4 false premises, §1.5 size, §2 patch protocol + awk extraction, §3
  sub-milestones with expected numstat, §3.1 reachability, §4 mutation table, §5 files, §6 gate script, §7 executor
  notes, §8 risks, §9 acceptance criteria, Appendix C, Appendix P fences). Sub-milestone ids continue: M7a, M7b, …
- The design's line numbers were measured at `b0fbcd3e`; re-measure every line/premise you rely on at `799660ac`
  and record drift in §1.2/§1.4.
- Scratch (compile your stages, run mutants, Chrome profiles) OUTSIDE the repo:
  `$HOME/.ailang/state/world-iter392/plan-scratch/` (never under /tmp — location-sensitive tests go red there).
- Write exactly two files in the worktree, and DO NOT git add/commit them (the controller commits):
  - `design_docs/planned/w-world-approvals-by-key-a4-sprint4-plan.md`
  - `design_docs/planned/sprint_w-world-approvals-by-key-a4-s4.json` (same shape as `sprint_w-world-approvals-by-key-a4-s3.json`)
- Leave every code/doc file in the worktree untouched (the executor applies your patches there). No git writes of any
  kind in the worktree.

## Read first
1. The design (sections above), the sprint-3 plan (format + Appendix C + §4 survivors + §6 gate script), and the code
   you build on: `host/daemon/surface_decide.go` (+ tests), `host/daemon/surface_approve.go` (A3a's requests.json and
   card route — your model for the decision listing and the card), `host/surfacepage/decide_model.go`,
   `host/surfacepage/approve*.go` + `approve/` (A3a's card page and `approve.js` — the model for `decide.js`),
   `host/surfacepage/passkey*.go` + `passkey/`, `host/surfacepage/app/app.js` (*Needs you*, Decided, the
   notification, the `N` key lives in `surface.js`), `cmd/ailang-worldd/approve*.go` (incl. `approve --list`,
   `approve_decide.go`), `host/approval/decision*.go`, `host/broker/decision_law.go`.
2. The live-Chrome tests you extend: `TestApproveDecideWebAuthnLinkInChrome`, `TestApprovePage*`, `TestPasskeyPage*`,
   the XA1–XA5 tests and their env gates.
3. `CLAUDE.md`, `design_docs/coding-standards.md` (binding: pure core, effects at the boundary, non-vacuous gates).

## Size rule
- **Cap: ≤ 2,000 changed lines**, measured `git diff --numstat <base> HEAD -- cmd host docs` (state this exact
  filter; also give `-- cmd host` separately). Prefer a smaller, fully proven sprint over a full, half-proven one.
- Sub-milestones of ≤ ~250–350 changed lines each, each independently green and safe at commit.
- No new `/v1` route (`TestFrozenV1RouteTableMatchesMux` stays 13). New `/surface` routes only as the design names
  them (`TestSurfaceAppRoutes` +2). `surfaceCSP` unchanged unless the design says otherwise. No CORS headers.

## Retro directives (iterations 374–391; each was paid for)
1. **Ship compiled patches**: every line of code, test and doc as a unified diff in Appendix P
   (```` ````diff <M> ```` fences), generated LAST from compiled, vetted, gofmt-clean, green scratch stages; then
   extract them FROM THE PLAN FILE with §2's awk, apply to a fresh `git archive 799660ac` export, and `diff -r`
   against your work tree — byte-identical or regenerate. The controller re-does this round trip independently.
   **Ship the patches BEFORE running the full mutant set, and checkpoint the plan file** (write a first complete
   version, then append the mutant ledger).
2. **Expected numstat after each sub-milestone's commit** (cumulative, measured after commit, with the §size filter).
3. **Mutants anchored to the DIFF**: one mutant per new production branch, enumerated from the diff, not from the
   test plan — in particular the `decisions` listing's filters and truncation (DMD1–DMD3), the card route's
   refusals (PD1/PD2), `decide.js`'s message bytes (verb, packet digest, request digest, the `BigInt` expiry), its
   hash and `requestRef` checks, its POST envelope, *Needs you*'s escaping and link, the notification's new-id test,
   `DecisionSubject.Passkeys`'s filter, `approve --list`'s decision rows, and the verb text. Run each against the
   named test row AND **the whole package** (a SURVIVED verdict requires a whole-package run — a narrow `-run`
   turned a killed mutant into a false survivor in 391); record killed/survived; name equivalents with reasons. A
   mutant that does not compile is not a mutant (pair with a killed control). **Name the test that killed each
   mutant; a mutant killed only by an inventory/static test is not behavioural proof — re-run with that test
   excluded and record both.** Write the ledger to a file in scratch; every re-run appends its own line. **The plan's
   mutant table ids must be the same ids as in the ledger** (the judge will `comm` them).
4. **New branch reachability**: for every new refusal/branch, give a trace from the request's entry point (the GET
   routes, the POST, the CLI, the page scripts) to that branch for each input class, and say which test drives it
   end to end vs unit-only. A branch unreachable end to end because an earlier step refuses first must say so.
5. **Repo-wide inventories**: run every design §7 inventory by name with `-v` on your final scratch tree (incl.
   `TestFrozenV1RouteTableMatchesMux`, `TestSurfaceAppRoutes`, `TestDaemonDependencyAllowlist`,
   `TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots`, `TestProductionGoSurface`,
   `TestEveryVerbHasHelp`, the quickstart tests, `TestDaemonCredentialWriterSites`,
   `TestCredentialsMethodsRequireADeadline`, `TestSurfaceAppPure`, any embedded-asset/static-file inventory for
   `host/surfacepage` (a new `decide/` asset dir WILL hit one), `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`,
   import allowlists for `host/approval`). Search for any test that enumerates routes, exported symbols, files,
   embedded assets, CLI verbs, CSP, schema versions or QUICKSTART sections before assuming none.
6. **A directive whose diff adds a production context root must name
   `go test ./host/store/ -run '^TestProductionContextRoots$'` as a gate** (charter process fix, iter-292).
7. The §6 gate script for the controller (one bash shell; **every leg records its own rc and the script runs every
   leg — no `set -e`**; `tee` every leg to `$HOME/.ailang/state/world-iter392/gates/` (mkdir it first), never only
   through `tail`), incl. `-race` on touched Go packages, the inventories, the live-Chrome leg (with its env, copied
   verbatim from the design/plan, and a temp profile), `./scripts/verify_ail.sh`, and `VERIFY_GO_LEGS=plain
   AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh` after commit. Measure each gate at the base first.
   **Never start a gate under `nohup`.** The script must run under macOS `/bin/bash` 3.2: **no `declare -A`**. The
   script writes `gates/ALL_DONE.<sha>` when every leg has finished (the judge waits on it).
8. Executor notes as the sprint-3 plan's §7, with ONE change: **snapshots copy only the files `git status
   --porcelain` reports changed (plus their deletions), NOT an `rsync` of whole directories** — the controller's
   `commit_ms.sh` commits every file present in a snapshot (391 had to override §7's rsync wording). No git writes;
   `patch -p1` is the only writer; runs of ~2–4 sub-milestones; cumulative snapshots to
   `$HOME/.ailang/state/world-iter392/snap/<M>/` with MANIFEST.sha256 (exclude the two plan files); stop
   conditions; numstat checks; exact acceptance commands per sub-milestone (copy any env such as `WORLD_CHROME`,
   `WORLD_REQUIRE_PASSKEY_BROWSER`, `WORLD_EXEC_NODE` verbatim into each command); **an explicit "never `nohup`,
   never background a test leg" line; and "every Go Bash call passes `timeout: 600000`, foreground, no Monitor"**.
9. Acceptance criteria (§9-style), numbered, that the judge can score against. Include what an operator can do after
   this sprint (open *Needs you*, open the card, decide with a passkey or an SSH signature, see it in Decided /
   `approve --list`) and what the attended exit run (design §10) still waits for.
10. **Every premise gets a §1.2 row**: any sentence saying "the code currently does X" cites the command and its
   observed output; an empty result needs a known-positive control in the same call. Where the design and HEAD
   disagree, measure, route minimally, record it in §1.4, and list every file you touch outside the design's §8.
11. **A test row that drives a fake tool must have a sibling row fed the real tool's captured response** (charter
   process fix, iter-299) — relevant to Chrome's WebAuthn assertion through the card (XD1 with a virtual
   authenticator) and to `ssh-keygen -Y sign` for the card's paste box: use an `approvaltest` key file / virtual
   authenticator, and check in any captured real output under `testdata/` with its generating command in a comment.
12. **No real daemon (D-WORLD-109, the 391 P7 slip)**: the sprint-3 planner's first JM46 test draft wrote an
   incomplete temp service record, the CLI fell back to `127.0.0.1:7644`, and one GET reached the machine's real
   World service. Every test that touches the CLI's service record (e.g. `approve --list`) must first PROVE its temp
   record resolves to the fake daemon (as `approve_decide_jgaps_test.go:31` does) before any dial, and your own
   scratch runs must never reach `127.0.0.1:7644`.
13. Close with this plan's own Appendix C: anything moved out (with size) and its handoff, plus the residuals the
   attended exit run must cover.

## Hard constraints
- D-WORLD-109: test keys, temp stores, `approvaltest` keys only. No real key, no `~/.ssh`, no `~/.ailang/world`, no
  real `/dev/tty`, no real browser profile (temp `--user-data-dir` only), no mint or publish against a real store,
  no real network listener beyond `httptest`/loopback in tests. A test may run `ssh-keygen -Y sign` only on an
  `approvaltest` key file in a temp dir.
- `host/approval` stays pure and standard-library-only; the laws are pure (no I/O, no clock reads — `now` is an
  argument). Effects only at the store/CLI/daemon boundary. The daemon never KEEPS the credentials writer. No CORS
  headers anywhere. `decidedBy` is derived by the daemon, never accepted from a caller. The card page computes
  nothing the POST trusts: the POST re-runs the law.
- Do not touch: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/` (except reading),
  `cmd/world-publish/*`, `~/.ailang/state/mission-v1*`, the A2/A3 mint law and its matrix tests, mints,
  `surface_approvals`, the sprint 1–3 laws and POST arms (reuse them unchanged) — and anything the design's §8 marks
  unchanged. If an A2/A3a file must change (see scope item 7), name it in §1.4 and §5 and keep its existing tests
  green and unedited unless you justify each edit.
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
  `TestWorkspaceProtectedPolicyWiring`, `TestConvertedStoreMethodsM5b`, host/pkgproj under load,
  `TestPasskeyPageAbortsAPendingCeremony/get-late` (Chrome, 240000.1 vs 240000 ms under load — 391).
- **Time budget: aim to finish within ~100 minutes; hard stop at 130.** The iteration has a 6 h wall clock that
  started 05:31Z; the executor, gates, judge, CI and record still follow you. If you run short, ship a smaller cut
  fully proven (e.g. the listing + *Needs you* + `Passkeys` + `approve --list` + the card with PD1 and XD1, with
  the rest priced into Appendix C) rather than everything half-proven.

## Waiting discipline (binding on you)
If you background anything (a long test run, a mutant loop, Chrome), poll its ARTIFACT (an output/rc file it writes
on completion) inside a bounded `date +%s` loop in the same turn. A Monitor is an event stream, not a wait; nothing
will re-invoke you when it fires. Never end your turn announcing a wait, a pause, or "once it completes" — end only
when both files are written and the round trip has passed.

## Your final message
Report: the plan path, scope (what is in, anything moved and why, with measured numbers), the sub-milestone list with
expected cumulative numstat after each, mutant counts (total / killed / named equivalents / survivors) and the ledger
path, the round-trip result, base-gate results, every premise found false (§1.4), the verb-text decision (scope item
7), how an operator decides a request from the surface after this sprint, and anything the controller must decide.
