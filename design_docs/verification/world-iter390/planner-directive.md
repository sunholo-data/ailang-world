MISSION-ROLE: planner

You are the sprint-planner for Ailang World mission iteration 390 (unattended run; no human is present).
Write the sprint plan for **row 201 `w-world-approvals-by-key`, milestone A4, SPRINT 2 — "the keys", headed by the
typed CLI moved from sprint 1**. No re-design: `design_docs/planned/w-world-approvals-by-key-a4.md` r2 is reviewed
(iteration 388: quorum r1+r2 blocked on narrow defects, controller carve-out, independent sonnet audit VALID 84) and
is NOT reopened. Sprint 1 LANDED in iteration 389 (PR #563, merge `338ff821`; judge sonnet PASS 88). D-WORLD-115
(the decision message with a derived challenge) is still OPEN; its default **A = keep it** is applied. If you find a
place where sprint 2 depends on the D-WORLD-115 answer (design §0 says only §2.3's challenge slot and sprint 2 change
under B), record it in §1.4 and build under A.

## Scope — read these verbatim first
- Design §0 "Sprint 2", §2.3, §2.7 (the keys: judge, evidence, counter, K0–K5), §2.8 (routes, handler order),
  §2.10 (CLI), §6 "Sprint 2" table, §7, §8, §9 "Sprint 2 = M4–M5" incl. its **fallback cut**.
- **Sprint-1 plan Appendix C** (`design_docs/planned/w-world-approvals-by-key-a4-sprint1-plan.md`, from line ~835):
  - **The head of this sprint is the typed CLI CD1–CD4** (`cmd/ailang-worldd/approve_decide.go`, ≈ 400 lines by the
    design's price), with everything Appendix C lists for it (id-shape dispatch beside the `mr-` check, `--deny`,
    `--await` refused with a request ref, finding the daemon for a request ref — sprint 2 decides this —, client-side
    hash checks, the printed fields, `y` at `/dev/tty` via `openTerminal`, `writeOperatorProof`, the POST with
    `Origin` and `typed: true`, help lines, `TestProductionContextRoots` +1, design Q2 `approve --list` decisions).
    The CLI is what makes sprint 1's `POST /surface/decide/{ref}` reachable by an operator; it is NOT cuttable
    in this sprint.
  - **N4 / K0's slot (S1-D9)** — the handler order is fixed; slot the key arms in exactly as Appendix C says, and
    edit sprint-1 code only at the places Appendix C names.
  - Carried: R4e, a concurrent-decide test for `humanMu` (named survivor `H-LOCK`), `TestDaemonCredentialWriterSites`
    1 → 2.
- **Sprint-1 judge's non-blocking test gaps, to close in this sprint (small, test-only):**
  - **J07**: the L1 reserved-scope effect-substring arm — mutant `strings.Contains` → `strings.HasSuffix` survived.
  - **J31/J32**: the daemon tests do not pin D1's request-ref regex anchor (`{64}$` → `{64}` in
    `host/daemon/surface_decide.go:29` survives: controller-reproduced in 389) or hex case; a regression would let
    `hashref.MustParse` panic on a 65-hex ref.
  - **J33**: the GET `/surface/decisions/{ref}` D4a mapping is untested.
  Evidence: `$HOME/.ailang/state/world-iter389/judge/` and `$HOME/.ailang/state/world-iter389/evaluator-r1.md`.
- **Size and the fallback**: CLI ≈ 400 + the design's sprint 2 ≈ 1,800 = ≈ 2,200, over the cap. Expect to apply the
  design's priced fallback: **the webauthn arm (RD10 passkey half, CD8–CD9, DM7, and the passkey halves of K0/KD7/
  KD8/RD13/RD14) moves to sprint 3**; SSHSIG ships alone. **Never cut**: CD1–CD4, R4e, RD13 (a) SSHSIG, the DDL
  CHECKs KD5/KD6, schema 8, the judge, the J-gap tests. Decide from the MEASURED size and record the decision. If the
  schema-8 table must carry passkey columns that only the cut arm uses, say how you keep the DDL forward-compatible so
  sprint 3 does not need schema 9 (or price schema 9 honestly).

## Where you work
- Sprint worktree (already created, branch `sprint/world-iter390-row201-a4-s2`, base = origin/dev
  `8338eebb173290168f07e6b426d461c0122f67a8`, which contains A4 sprint 1):
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter390-a4-s2`
- VERIFIED BY CONTROLLER at the base (commands run in the worktree):
  - `host/store/credentials.go:43` `const credentialsSchemaVersion = 7`.
  - `host/authority/judge.go` exists (the A2 judge whose step names JD1–JD3 must agree with).
  - `cmd/ailang-worldd/approve.go:199` refuses a non-`mr-` id: `%q is not a mint request id (mr- and 16 lowercase hex)`.
  - `host/daemon/surface_decide.go` exists with the typed arm only (line ~179: "The grammar admits typed only in
    sprint 1").
  - Pinned ailang: `$HOME/.pinned-ailang/ailang` = AILANG v0.41.0.
  - Rig load average ≈ 30 on 16 cores (lighter than iteration 389's 210, but other sessions run). Expect some
    timing-sensitive test reds in untouched packages; re-run once alone before believing a red.
- Read the landed sprint-1 plan for the plan FORMAT and reproduce it (§1.2 premise table, §1.3 decisions, §1.4 false
  premises, §2 patch protocol + awk extraction, §3 sub-milestones with expected numstat, §3.1 reachability, §4
  mutation table, §5 files, §6 gate script, §7 executor notes, §8 risks, §9 acceptance criteria, Appendix C,
  Appendix P fences).
- The design's line numbers were measured at `b0fbcd3e`; re-measure every line/premise you rely on at `8338eebb`
  and record drift in §1.2/§1.4.
- Scratch (compile your stages, run mutants) OUTSIDE the repo:
  `$HOME/.ailang/state/world-iter390/plan-scratch/` (never under /tmp — location-sensitive tests go red there).
- Write exactly two files in the worktree, and DO NOT git add/commit them (the controller commits):
  - `design_docs/planned/w-world-approvals-by-key-a4-sprint2-plan.md`
  - `design_docs/planned/sprint_w-world-approvals-by-key-a4-s2.json` (same shape as `sprint_w-world-approvals-by-key-a4-s1.json`)
- Leave every code/doc file in the worktree untouched (the executor applies your patches there). No git writes of any
  kind in the worktree.

## Read first
1. The design (sections above), the sprint-1 plan (format + Appendix C), and the sprint-1 code you build on:
   `host/approval/decision.go`, `host/broker/decision_law.go`, `host/broker/approve.go`, `host/daemon/surface_decide.go`
   and its tests, `host/coordinator` discovery.
2. The A2/A3 key machinery you reuse: `host/approval/` (`ApproveBySignature`, `ApproveByAssertion`, `Binds`, the
   SSHSIG/assertion verifiers), `host/authority/judge.go`, `host/store/credentials.go` (schema steps v1→v7, the A2/A3
   upgrade pattern, `usePasskey`, `useApprover`, the writer and its deadline rule), `credentials_schema.sql`,
   `cmd/ailang-worldd/approve*.go` (F1 flag rules, `openTerminal`, `writeOperatorProof`, mint-mode flags
   `--sign-with`, `--signature`, `--write-message`, `--webauthn-assertion`), `host/approvaltest`.
3. `CLAUDE.md`, `design_docs/coding-standards.md` (binding: pure core, effects at the boundary, non-vacuous gates).

## Size rule
- **Cap: ≤ 2,000 changed lines**, measured `git diff --numstat <base> HEAD -- cmd host` (state this exact filter).
  Prefer a smaller, fully proven sprint over a full, half-proven one.
- Sub-milestones of ≤ ~250–350 changed lines each, each independently green and safe at commit.
- No new `/v1` route (`TestFrozenV1RouteTableMatchesMux` stays 13). `surfaceCSP` unchanged. No CORS headers.

## Retro directives (iterations 374–389; each was paid for)
1. **Ship compiled patches**: every line of code, test and doc as a unified diff in Appendix P
   (```` ````diff <M> ```` fences), generated LAST from compiled, vetted, gofmt-clean, green scratch stages; then
   extract them FROM THE PLAN FILE with §2's awk, apply to a fresh `git archive 8338eebb` export, and `diff -r`
   against your work tree — byte-identical or regenerate. The controller re-does this round trip independently.
   **Ship the patches BEFORE running the full mutant set, and checkpoint the plan file** (write a first complete
   version, then append the mutant ledger).
2. **Expected numstat after each sub-milestone's commit** (cumulative, measured after commit, `-- cmd host`).
3. **Mutants anchored to the DIFF**: one mutant per new production branch, enumerated from the diff, not from the
   test plan — in particular each judge step, each K-rule (K0 exact match on artifact AND signed fields, K1–K5), the
   SQL counter check, both DDL CHECK clauses, the schema upgrade step, the POST format switch, `decidedBy`
   derivation, the land call, the CLI's id-shape dispatch, subject-hash checks, `--deny`/`--await` flag rules, the
   SSHSIG namespace and the 91-byte message written by `--write-message`. Run each against the named test row AND the
   package; record killed/survived; name equivalents with reasons. A mutant that does not compile is not a mutant
   (pair with a killed control). **Name the test that killed each mutant; a mutant killed only by an inventory/static
   test is not behavioural proof — re-run with that test excluded and record both.** Write the ledger to a file in
   scratch; every re-run appends its own line. **The plan's mutant table ids must be the same ids as in the ledger.**
4. **New branch reachability**: for every new refusal/branch, give a trace from the request's entry point (the POST
   handler, the CLI) to that branch for each input class, and say which test drives it end to end vs unit-only. A
   branch unreachable end to end because an earlier step refuses first must say so.
5. **Repo-wide inventories**: run every design §7 inventory by name with `-v` on your final scratch tree (incl.
   `TestFrozenV1RouteTableMatchesMux`, `TestSurfaceAppRoutes`, `TestDaemonDependencyAllowlist`,
   `TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots`, `TestProductionGoSurface`,
   `TestEveryVerbHasHelp`, the quickstart tests, `TestDaemonCredentialWriterSites`,
   `TestCredentialsMethodsRequireADeadline`,
   `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`, import allowlists for `host/approval`). Search for
   any test that enumerates routes, exported symbols, files, CLI verbs, schema versions or QUICKSTART sections
   before assuming none.
6. **A directive whose diff adds a production context root must name
   `go test ./host/store/ -run '^TestProductionContextRoots$'` as a gate** (charter process fix, iter-292); a new
   root is pinned in `host/store/context_roots_test.go` with its bound in the comment.
7. The §6 gate script for the controller (one bash shell; **every leg records its own rc and the script runs every
   leg — no `set -e` stopping at the first red** (iteration-389 retro); `tee` every leg to
   `$HOME/.ailang/state/world-iter390/gates/` (mkdir it first), never only through `tail`), incl. `-race` on touched
   Go packages, the inventories, `./scripts/verify_ail.sh`, and `VERIFY_GO_LEGS=plain
   AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh` after commit. Measure each gate at the base first.
   **Never start a gate under `nohup`.** The script must run under macOS `/bin/bash` 3.2: **no `declare -A`**.
8. Executor notes as the sprint-1 plan's §7: no git writes; `patch -p1` is the only writer; runs of ~2–4
   sub-milestones; cumulative snapshots to `$HOME/.ailang/state/world-iter390/snap/<M>/` with MANIFEST.sha256
   (exclude the two plan files); stop conditions; numstat checks; exact acceptance commands per sub-milestone;
   **an explicit "never `nohup`, never background a test leg" line; and "every Go Bash call passes `timeout: 600000`,
   foreground, no Monitor"** (iteration 389: an executor hit the Bash tool's 120 s default, backgrounded, ended its
   turn waiting, and the shell was lost).
9. Acceptance criteria (§9-style), numbered, that the judge can score against. Include what an operator can do after
   this sprint (approve/deny a packet at a terminal by typed `y`; by an SSH key) and what the attended exit run
   (design §10) still waits for.
10. **Every premise gets a §1.2 row**: any sentence saying "the code currently does X" cites the command and its
   observed output; an empty result needs a known-positive control in the same call. Where the design and HEAD
   disagree, measure, route minimally, record it in §1.4, and list every file you touch outside the design's §8.
11. **A test row that drives a fake tool must have a sibling row fed the real tool's captured response** (charter
   process fix, iter-299) — relevant to `ssh-keygen -Y sign` in CD5–CD7: use an `approvaltest` key file, and check
   in any captured real output under `testdata/` with its generating command in a comment.
12. Close with this plan's own Appendix C: anything moved to sprint 3 (with size), and the sprint-3 handoff.

## Hard constraints
- D-WORLD-109: test keys, temp stores, `approvaltest` keys only. No real key, no `~/.ssh`, no `~/.ailang/world`, no
  real `/dev/tty`, no real browser profile, no mint or publish against a real store, no real network listener beyond
  `httptest`/loopback in tests. Never run bare `ssh-keygen` against a real key; a test may run `ssh-keygen -Y sign`
  only on an `approvaltest` key file in a temp dir.
- `host/approval` stays pure and standard-library-only; the laws are pure (no I/O, no clock reads — `now` is an
  argument). Effects only at the store/CLI/daemon boundary. The daemon never KEEPS the credentials writer. No CORS
  headers anywhere. `decidedBy` is derived by the daemon, never accepted from a caller.
- Do not touch: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/` (except reading),
  `cmd/world-publish/*`, `~/.ailang/state/mission-v1*`, the `/surface/passkey` route, `passkey.js`, `approve.js`,
  the A2/A3a law and its matrix tests, mints, `surface_approvals`.
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
  `TestMissionConfigDefaultFlowIntegration/invalid_ledger`, `scripts/marketingcheck` supervisor tests,
  `TestMarketingCaptureInventoryMutation`, `TestSessionSelfGit`, `TestServeStartsWhenTheProbeSpendsTheStartupClock`
  (red at base), `TestRegistryPhaseHasItsOwnClock`, `TestWorkspaceProtectedPolicyWiring`,
  `TestConvertedStoreMethodsM5b`, host/pkgproj under load.
- **Time budget: aim to finish within ~100 minutes; hard stop at 140.** The iteration has a 6 h wall clock that
  started 21:32Z; the executor, gates, judge, CI and record still follow you. If you run short, ship a smaller cut
  fully proven (CLI + J-gaps + schema 8 + judge + SSHSIG arm) rather than everything half-proven.

## Waiting discipline (binding on you)
If you background anything (a long test run, a mutant loop), poll its ARTIFACT (an output/rc file it writes on
completion) inside a bounded `date +%s` loop in the same turn. A Monitor is an event stream, not a wait; nothing will
re-invoke you when it fires. Never end your turn announcing a wait, a pause, or "once it completes" — end only when
both files are written and the round trip has passed.

## Your final message
Report: the plan path, scope (what is in, anything moved and why, with measured numbers), the sub-milestone list with
expected cumulative numstat after each, mutant counts (total / killed / named equivalents / survivors) and the ledger
path, the round-trip result, base-gate results, every premise found false (§1.4), how the CLI finds the daemon for a
request ref, the schema-8 forward-compatibility decision, and anything the controller must decide.
