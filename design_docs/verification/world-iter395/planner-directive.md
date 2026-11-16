MISSION-ROLE: planner

You are the sprint-planner for Ailang World mission iteration 395 (unattended run; no human is present).
Write the sprint plan for **row 243 `w-world-address-bar-verbs`, phase V1 — the ONE sprint (M1–M4)**. No re-design:
`design_docs/planned/w-world-address-bar-verbs.md` r2 is reviewed (quorum r1+r2 blocked 3/3 on narrow defects,
controller narrow-refinement carve-out r2, independent sonnet audit CARVE-OUT VALID 84, iteration 394) and is NOT
reopened. Build V1 as designed. V1 is read-only by construction (§13): no POST, no store write, no process spawn, no
key path, no `.ail`, no `/v1` route.

## Scope (design §0.3, §2.1–§2.4, §6, §8, §9 — read them verbatim)
- **M1 `host/verbs`** — pure Go, no I/O: `Split`, `Quote`, `Lookup` + the verb table, `ParseInterleaved` (moved from
  `cmd/ailang-worldd/session.go`, identical body), `EpisodeReportFlags` + the bar admission; tests VB1–VB5 (and RK1's
  shadow set as the design places it).
- **M2 CLI repoint** — the call sites of `parseInterleaved` (design says 10; re-measure), `episode report` on the
  shared flags, the printer moved to `host/daemon`. The CLI's own tests stay unedited (design V16, VR3).
- **M3 `GET /surface/verb.json?q=<line>`** — one registration under `approvalRoute`, `d.episodeReport` extracted from
  the `/v1` handler; tests VR1–VR4, VS1. `/v1` stays 13 routes.
- **M4 the page** — `host/surfacepage/app/app.js` + `index.html`: `go(raw)` asks the route first; no silent fallback
  (only a schema-valid `outcome:"address"` navigates; every failure is an explicit refusal with the view unchanged);
  the bounded page wait (`VERB_ROUTE_BUDGET_MS`, AbortController, design r2); step 4b; tests VP1–VP5, XV1 (live
  Chrome, as the A3a/A4 plans ran theirs), QUICKSTART §19.
- **Two calls the r2 audit LEFT TO YOU — decide each, record it as a §1.3 decision with the reason:**
  (1) how a *superseded* in-flight request (a second submit while the first is pending) is aborted and how its
  outcome differs from the budget-timeout message (the controller's own sentence "a second submit aborts the first"
  was withdrawn as beyond the reviewers' text — it is your call, not a given); (2) VP2's hang test: how the page's
  budget is exercised without a real 15 s wall (a controlled clock / injected budget) — and whether the 15 s budget is
  a tunable. Neither may weaken "no silent fallback".
- **Fallback cut (design §9):** if M1–M3 measure over 1,500 changed lines, M4's live-Chrome XV1 and the `approve`
  navigate arm move to V2's first sprint (the route answers `terminal` for `approve` meanwhile). Decide from the
  measured size and record it.

## Where you work
- Sprint worktree (already created, branch `sprint/world-iter395-row243-v1`, base = origin/dev
  `db2b7423327d5223cd1a37cf40e0794c2019026a`): `/Users/voightkampff/dev/sunholo-data/.wt-world-iter395-row243-v1`
- **The design was measured at `48c92a52`. Four attended PRs merged since (#577–#580)** and they touch the surface
  this sprint edits. VERIFIED BY CONTROLLER (`git diff --stat 48c92a52 db2b7423 -- cmd host docs`, 12 files,
  +272/−3): `host/daemon/daemon.go` (+1, a new route registration), `host/daemon/surface_home_summary.go` (new,
  `GET /surface/home-summary.json`), `host/surfacepage/home_summary.go` (new), `host/surfacepage/app/app.css`,
  `host/surfacepage/app_test.go` (+16), `host/surfacepage/approve/approve.js` (+34; `TestDecidePageAssets` banks its
  hash), `host/daemon/surface_follow_test.go`, `cmd/ailang-worldd/service.go`, `cmd/ailang-worldd/surface_follow_test.go`,
  `docs/QUICKSTART.md` (+3/−1). So: any route enumeration (`TestSurfaceAppRoutes` and friends), any asset-hash test
  over `app.js`/`index.html`, and the QUICKSTART section numbering/tests may have moved. Re-measure every design
  line/premise you rely on at `db2b7423` and record drift in §1.2/§1.4.
- Pinned ailang: `$HOME/.pinned-ailang/ailang` = AILANG v0.41.0 (the Go suite's `AILANG_BIN`). CI's verify gate
  installs `releases/latest` (v0.60.4). V1 changes no `.ail`, so `./scripts/verify_ail.sh` must give the base's
  result byte-for-byte.
- Rig load average ≈ 7–14 on 16 cores right now (it has been ≈ 200 before; it varies). Re-run a red in an untouched
  package once alone before believing it.
- Read the newest landed plans for the plan FORMAT: `design_docs/planned/w-world-approvals-by-key-a4-sprint4-plan.md`
  (the surface sprint — closest to this one: page JS, live Chrome, asset hashes), then sprint 1's (§1.2 premise table,
  §1.4 false premises, §2 patch protocol + awk extraction, §6 gate script, §7 executor notes, Appendix P fences,
  expected numstat per sub-milestone, Appendix C). Reproduce that format. The sprint JSON shape:
  `design_docs/planned/sprint_w-world-approvals-by-key-a4-s4.json`.
- Scratch (compile your stages, run mutants) OUTSIDE the repo: `$HOME/.ailang/state/world-iter395/plan-scratch/`
  (exists; never under /tmp — location-sensitive tests go red there).
- Write exactly two files in the worktree, and DO NOT git add/commit them (the controller commits):
  - `design_docs/planned/w-world-address-bar-verbs-v1-sprint-plan.md`
  - `design_docs/planned/sprint_w-world-address-bar-verbs-v1.json`
- Leave every other file in the worktree untouched (the executor applies your patches there). No git writes.

## Read first
1. `design_docs/planned/w-world-address-bar-verbs.md` r2 in full (809 lines; §0–§9 and §14 V-rows, §15 for why r1/r2
   are shaped as they are). Re-measure what you rely on, live.
2. The code you build on: `cmd/ailang-worldd/session.go` (`parseInterleaved`), every caller of it,
   `cmd/ailang-worldd/episode_report.go`, `host/daemon/episode_report.go`, `host/daemon/daemon.go` route registration
   and `approvalRoute`, `host/daemon/surface_*.go` (the surface guards a read route reuses), `host/surfacepage/app/`
   (`app.js`, `index.html`, the tests over them), the static `surface.js` (`parseAddress` — untouched by V1),
   `cmd/ailang-worldd/main.go`'s verb switch and help, the inventories in design §7.
3. `CLAUDE.md`, `design_docs/coding-standards.md` (binding: pure core, effects at the boundary, non-vacuous gates).

## Size rule
- **Cap: ≤ 2,000 changed lines**, measured `git diff --numstat <base> HEAD -- cmd host docs` (state this exact
  filter; the design sized with it). Design estimate ≈ 1,750. A move counts both sides. Prefer a smaller, fully proven
  sprint over a full, half-proven one.
- Sub-milestones of ≤ ~250–350 changed lines each, each independently green and safe at commit (your call, measured).

## Retro directives (iterations 374–394; each was paid for)
1. **Ship compiled patches**: every line of code, test and doc as a unified diff in Appendix P
   (```` ````diff <M> ```` fences), generated LAST from compiled, vetted, gofmt-clean, green scratch stages; then
   extract them FROM THE PLAN FILE with §2's awk, apply to a fresh `git archive db2b7423` export, and `diff -r`
   against your work tree — byte-identical or regenerate. The controller re-does this round trip independently.
   **Ship the patches BEFORE running the full mutant set, and checkpoint the plan file** (write a first complete
   version, then append the mutant ledger).
2. **Expected numstat after each sub-milestone's commit** (measured after commit, `-- cmd host docs`).
3. **Mutants anchored to the DIFF**: one mutant per new production branch, enumerated from the diff, not from the
   test plan — in particular each `Split` quoting/escape arm, `Quote`'s escape set, each `Lookup` outcome class
   (address / run / terminal / refused / unknown), the shadow-set membership, step 3b/4b (bare and unknown `episode`,
   three-word `world …`), the leading `ailang-worldd` strip, `EpisodeReportFlags` admission, the route's guards and
   its `q` bound, the report byte-identity, and in `app.js` the no-fallback arms (non-2xx, network error, schema-invalid
   body, budget timeout, superseded request) — each must refuse with the view unchanged. Run each against the named
   test row AND the package; record killed/survived; name equivalents with reasons. A mutant that does not compile is
   not a mutant (pair with a killed control). **Name the test that killed each mutant; a mutant killed only by an
   inventory/static/hash test is not behavioural proof — re-run with that test excluded and record both.** Write the
   ledger to a file in scratch; every re-run appends its own line. **The plan's mutant table ids must be the same ids
   as in the ledger.**
4. **New-branch reachability (iterations 388/394 retros)**: for every new step, path or refusal, trace each input
   class AND each failure/timeout class from the entry point (the address bar's `go(raw)`; the CLI's argv; the route's
   `q`) to its outcome, and say which test drives it end to end vs unit-only. Name any branch that is unreachable end
   to end because an earlier step refuses first. (Round 2 of this doc's quorum blocked on exactly two surfaces round
   1's fixes had created — a classification step with no head-match and a fetch with no deadline.)
5. **Repo-wide inventories**: run every design §7 inventory by name with `-v` on your final scratch tree (incl.
   `TestFrozenV1RouteTableMatchesMux`, `TestSurfaceAppRoutes`, `TestDaemonDependencyAllowlist`,
   `TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots`, `TestProductionGoSurface`,
   `TestEveryVerbHasHelp`, the quickstart tests, the page asset-hash tests). Search for any test that enumerates
   routes, exported symbols, packages, files, CLI verbs, page assets or QUICKSTART sections before assuming none — a
   new package `host/verbs` will meet every package allowlist.
6. **Parsers and fixtures**: the route returns `episode report`'s text "byte-identical to the CLI". The fixture that
   proves it must be GENERATED by running the real CLI path against a test store, never hand-typed, and a second
   fixture from a different real invocation (e.g. a report with and without the optional fields) must be pinned too.
7. The §6 gate script for the controller (one bash shell, `tee` every leg to `$HOME/.ailang/state/world-iter395/gates/`,
   never only through `tail`), incl. `go test -run '^$' ./...` as the test-file compile fence, `-race` on touched Go
   packages, the inventories, the node/Chrome legs as CI runs them (`WORLD_EXEC_NODE`), `./scripts/verify_ail.sh`, and
   `VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh` after commit. Measure each gate
   at the base first. **Never start a gate under `nohup`.**
8. Executor notes as the A4 plans' §7: no git writes; `patch -p1` is the only writer; runs of ~2–4 sub-milestones;
   cumulative snapshots to `$HOME/.ailang/state/world-iter395/snap/<M>/` with MANIFEST.sha256 (exclude the two plan
   files); stop conditions; numstat checks; exact acceptance commands per sub-milestone; **an explicit "never `nohup`,
   never background a test leg; pass a 600000 ms timeout to every long Bash call" line** (the iteration-389 executor
   lost a backgrounded acceptance shell to the 120 s default).
9. Acceptance criteria (§9-style), numbered, that the judge can score against. Include what an operator can do after
   this sprint (type `episode report <e>` into the app's bar and read the CLI's text; see the exact CLI line for a
   mutating verb) and what the attended exit run (design §10) still waits for.
10. **Every premise gets a §1.2 row**: any sentence saying "the code currently does X" cites the command and its
    observed output; an empty result needs a known-positive control in the same call. Where the design and HEAD
    disagree, measure, route minimally, record it in §1.4, and list every file you touch outside the design's §8.
11. Close with Appendix C: anything moved to V2 (with size) and the V2 handoff.

## Hard constraints
- D-WORLD-109: test stores, temp dirs and test keys only. No real key, no `~/.ssh`, no `~/.ailang/world`, no real
  `/dev/tty`, no real browser profile, no mint or publish, no real network listener beyond `httptest`/loopback in tests.
- `host/verbs` stays pure and standard-library-only (no I/O, no clock). Effects only at the CLI/daemon boundary. The
  route is read-only: no store write, no spawn. No CORS headers anywhere. `surfaceCSP` unchanged. No new `/v1` route.
  `worlddapi.ail` and every `.ail` untouched. The static page's own bar (`surface.js`) untouched.
- Do not touch: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/` (except reading),
  `cmd/world-publish/*`, `~/.ailang/state/mission-v1*`, row 201's laws/routes/cards (`approve.js`, `passkey.js`, the
  decide/approve handlers), the `approve` CLI.
- Every existing test the design lists as "unedited" stays unedited.
- Export `AILANG_BIN=$HOME/.pinned-ailang/ailang` in every shell that runs Go tests.
- Your Bash tool's shell is zsh: run loops/globs/`$VAR` work under `bash -c` or `bash <<'EOF'`. macOS has no
  `timeout` command. Pass `timeout: 600000` on any Bash call that may run long; the default is 120 s.
- Known load-sensitive / flaky tests (re-run once alone, do not chase): `TestSeedRestartIgnoresASurvivingOldCopy`
  (host/broker), `TestHostsLaunchersExecHost`, host/broker exec-probe timeouts incl.
  `TestExecProtectedCheckTimeout/reserve`, `TestPublishApprovalRefusalSetWithALandedPositiveControl`, local `-race` of
  host/broker, host/daemon, host/verifygate, `TestMCPPostBudgetProductionConstants/batch`, `TestBackoff`
  (host/surfacefollow), `TestDemoSeedPreflightRefusals`, `TestWorktreeSnapshotRefusesAConfiguredFilter`,
  `TestMissionConfigDefaultFlowIntegration` (host/verifygate), `scripts/marketingcheck` supervisor tests.
- **Time budget: aim to finish within ~100 minutes; hard stop at 150.** The iteration's 6 h wall clock started
  13:39Z, and the executor, gates, judge, CI and record still follow you. If you run short, ship the design's priced
  fallback cut fully proven rather than everything half-proven.

## Waiting discipline (binding on you)
If you background anything (a long test run, a mutant loop), poll its ARTIFACT (an output/rc file it writes on
completion) inside a bounded `date +%s` loop in the same turn. A Monitor is an event stream, not a wait; nothing will
re-invoke you when it fires. Never end your turn announcing a wait, a pause, or "once it completes" — end only when
both files are written and the round trip has passed.

## Your final message
Report: the plan path, scope (what is in, anything moved and why, with measured numbers), the sub-milestone list with
expected numstat after each, mutant counts (total / killed / named equivalents) and the ledger path, the round-trip
result, base-gate results, every premise found false (§1.4), your two §1.3 decisions (superseded request; VP2 clock),
and anything the controller must decide.
