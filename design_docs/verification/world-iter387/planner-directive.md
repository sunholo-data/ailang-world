MISSION-ROLE: planner

You are the sprint-planner for Ailang World mission iteration 387 (unattended run; no human is present).
Write the sprint plan for **row 201 `w-world-approvals-by-key`, milestone A3a, SPRINT 4** — the last A3a sprint:
the scope the sprint-3 plan's fallback cut moved here, plus the iteration-386 judge's test gap. No re-design:
`design_docs/planned/w-world-approvals-by-key-a3.md` r2 is reviewed (quorum r1+r2 + controller carve-out +
independent audit VALID 84) and is NOT reopened. D-WORLD-114 (the doc's two amendments to the attended parent) is
still OPEN with default A = keep both, applied since iteration 384; build as designed. The browser notification in
this sprint is the WEB notification of design §2.9 (A3a), not the native/Tauri one D-WORLD-114 (a) defers.

## Scope (from `design_docs/planned/w-world-approvals-by-key-a3a-sprint3-plan.md` §"Appendix C" — read it verbatim)
1. **The opt-in browser notification (design §2.9, XA4)**: the "Notify me" control calling
   `Notification.requestPermission()` only from its click; on a refresh with a pending id not seen before in the tab
   and permission `granted`, one `new Notification("World needs you", {body: "Approve mint request for <episode>
   (<n> grants), expires <hh:mm>", tag: <id>})` whose click focuses the tab and opens `/surface/approve/<id>`;
   **XA4** `TestNeedsYouNotifiesOnlyWhenAskedOnce` with the Chrome driver's `grantNotifications` command and the
   `notify` recorder injection (design §6; sprint-3 plan §1.4 P3 moved these here).
2. **The Decided list** under Needs you ("You · Minted · <episode> · <approvedBy> · <hh:mm>", the last few, a
   "show all" toggling the 20 of `requests.json`'s `decided`) and **XA3's Decided half** ("after the mint, Decided
   shows `webauthn:<name>`").
3. **The `N` key** jumping to the first Needs-you link. Sprint-3 plan §1.4 P7 / §1.2 C15 found the `N` key lives in
   `host/surfacepage/surface.js` (the static pages' `KEYS` table and `case "n": case "N": return {ui: "needs"}` at
   ~line 156), and `app.js` has none. Decide from the design §2.8 and the code where the app's `N` belongs, measure,
   and record the decision in §1.3.
4. **Judge survivor A02/A03/A05 from iteration 386 (controller-reproduced, MUST be closed)**: in `needsYouHTML`
   (`host/surfacepage/app/app.js` ~line 426), removing `esc()` on `q.agent` leaves `TestSurfaceAppPure` green
   (rc 0, with `WORLD_EXEC_NODE` set). The same holds for the `id` in the href/`id=` attribute and `q.recorded`.
   Add hostile-value cases (e.g. `<img src=x onerror=…>`, a `"` breaking the attribute) to the pure test so that
   removing EACH of those `esc()` calls fails it; prove each kill (apply mutant → red, revert → green) and name the
   killing test. Do the same for every new `esc()` the Decided list adds.
5. **Optional (only if it fits under the cap)**: one live-Chrome case that renders a hostile signed-field value on
   the card page (judge survivor C03/C04: no Chrome test renders a hostile value) and asserts it appears as text.
   If it does not fit, list it in this plan's Appendix C with a size estimate.

## Where you work
- Sprint worktree (already created, branch `sprint/world-iter387-row201-a3a-s4`, base = origin/dev
  `bed9d135fe38a1c8f7e2784414aaac348b76879f`, which contains A3a sprints 1–3 merged (`b44b5f7d`, `d1027477`,
  `70823879`)): `/Users/voightkampff/dev/sunholo-data/.wt-world-iter387-a3a-s4`
- VERIFIED BY CONTROLLER at the base (commands run in the worktree):
  - `grep -c Notification host/surfacepage/app/app.js` → 0. `grep -rln grantNotifications cmd host` → no match;
    control in the same call: `getRecorder|fetchRewrite` match `cmd/ailang-worldd/approve_page_driver_test.go`,
    `approve_page_swap_test.go`, `approve_page_chrome_test.go`, `testdata/passkey_cdp.mjs`.
  - `needsYouHTML(r, now)` is at `host/surfacepage/app/app.js:426` and is exported via `module.exports` (line 435);
    it renders `r.pending` only. `requests.json`'s `decided` exists: `host/surfacepage/approve_model.go:68`
    `Decided []DecidedRequest json:"decided"`, filled in `host/daemon/surface_approve.go:99–105`.
  - `N`: `host/surfacepage/surface.js:156` `case "n": case "N": return {ui: "needs"};` and the help rows at 257/259.
  - Chrome is at `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`. Pinned ailang:
    `$HOME/.pinned-ailang/ailang` = AILANG v0.41.0.
- Read the landed sprint-3 plan (§1.2, §1.4 P1–P8, §6, §7, Appendix C) — its premises are facts of HEAD you build
  on. The design's line numbers were measured at `0140bd7f`; re-measure every line/premise you rely on at
  `bed9d135` and record drift in §1.2/§1.4.
- Scratch (compile your stages, run mutants, Chrome profiles) OUTSIDE the repo:
  `$HOME/.ailang/state/world-iter387/plan-scratch/` (never under /tmp — location-sensitive tests go red there).
- Write exactly two files in the worktree, and DO NOT git add/commit them (the controller commits):
  - `design_docs/planned/w-world-approvals-by-key-a3a-sprint4-plan.md`
  - `design_docs/planned/sprint_w-world-approvals-by-key-a3a-s4.json` (same shape as `sprint_w-world-approvals-by-key-a3a-s3.json`)
- Leave every code/doc file in the worktree untouched (the executor applies your patches there). No git writes of any kind in the worktree.

## Read first
1. `design_docs/planned/w-world-approvals-by-key-a3.md` r2, especially §2.8 (Needs you / Decided), §2.9 (the
   notification), §2.10, §3, §6 "Sprint 3" (XA3, XA4, the driver commands), §7 inventories, §8 conflict surface,
   §9 gates and the exit run, §14 V-rows on browser behaviour (Notification permission, CDP
   `Browser.grantPermissions`) — re-measure what you rely on, live.
2. `design_docs/planned/w-world-approvals-by-key-a3a-sprint3-plan.md` for the plan FORMAT (newest; it worked:
   §1.2 premise table, §1.4 false premises, §2 patch protocol + awk extraction, §6 gate script, §7 executor notes,
   Appendix P fences, expected numstat per sub-milestone, Appendix C). Reproduce that format.
3. The code you build on: `host/surfacepage/app/app.js` (+ `app.css`, `index.html`), `host/surfacepage/surface.js`,
   `host/surfacepage/app_test.go` (`TestSurfaceAppPure` and how it runs node via `WORLD_EXEC_NODE`),
   `host/surfacepage/approve_model.go`, `host/daemon/surface_approve.go`, the Chrome harness in `cmd/ailang-worldd`
   (`approve_page_chrome_test.go`, `approve_page_driver_test.go`, `testdata/passkey_cdp.mjs`), `docs/QUICKSTART.md`
   §17 and the quickstart tests that pin it (the A3a exit run waits on this sprint's notification step).
4. `CLAUDE.md`, `design_docs/coding-standards.md` (binding: pure core, effects at the boundary, non-vacuous gates).

## Size rule
- **Cap: ≤ 2,000 changed lines**, measured `git diff --numstat <base> HEAD -- cmd host` (state this exact filter;
  `.js`/`.html`/`.mjs`/`.css` under `host/`/`cmd/` count; `docs/QUICKSTART.md` does not, but list its lines). The
  sprint-3 plan estimated the moved part at ≈ 300; this should be a small sprint. Prefer a smaller, fully proven
  sprint.
- Sub-milestones of ≤ ~250–350 changed lines each, each independently green and safe at commit.
- `surfaceCSP` is unchanged unless the design says otherwise (check: a Notification needs no CSP change; say so
  with evidence). No new daemon route unless the design names one.

## Retro directives (iterations 374–386; each was paid for)
1. **Ship compiled patches**: every line of code, test, JS and doc as a unified diff in Appendix P
   (```` ````diff <M> ```` fences), generated LAST from compiled, vetted, gofmt-clean, green scratch stages; then
   extract them FROM THE PLAN FILE with §2's awk, apply to a fresh `git archive bed9d135` export, and `diff -r`
   against your work tree — byte-identical or regenerate. The controller re-does this round trip independently.
   **Ship the patches BEFORE running the full mutant set, and checkpoint the plan file** (write a first complete
   version, then append the mutant ledger).
2. **Expected numstat after each sub-milestone's commit** (measured after commit, `-- cmd host`).
3. **Mutants anchored to the DIFF**: one mutant per new production branch (JS and Go: the click-only permission
   request, the granted check, the seen-id set (once per id per tab), the notification tag/body, the click handler
   (focus + open card), the Decided rendering + each `esc()`, "show all" toggle, the `N` key handler and its target,
   and each existing Needs-you `esc()` the A02 test now pins), enumerated from the diff, not from the test plan; run
   each against the named test row AND the package; record killed/survived; name equivalents with reasons. A mutant
   that does not compile/parse is not a mutant (pair with a killed control). **Name the test that killed each
   mutant; a mutant killed only by an inventory/static banned-token test is not behavioural proof — re-run with that
   test excluded and record both.** Copy the plan's acceptance environment (e.g. `WORLD_EXEC_NODE`, `WORLD_CHROME`,
   `WORLD_REQUIRE_PASSKEY_BROWSER`) verbatim into every mutant command. Write the ledger to a file in scratch; every
   re-run appends its own line.
4. **Live Chrome**: the XA tests run on this machine (`WORLD_CHROME=... WORLD_REQUIRE_PASSKEY_BROWSER=1`, design
   §9). Run them yourself on your final scratch tree and record `--- PASS` per subtest; state which tests skip in CI
   (no Chrome step) and what in CI still pins the behaviour (the pure node tests). Chrome profiles go in scratch,
   temporary, headless; never a real profile. Measure live whether headless Chrome supports `Notification` +
   `Browser.grantPermissions` for the test origin before designing XA4 around it; if a premise fails, record it in
   §1.4 and route minimally (e.g. the `notify` recorder injection).
5. **Repo-wide inventories**: run every design §7 inventory by name with `-v` on your final scratch tree (incl.
   `TestFrozenV1RouteTableMatchesMux`, `TestSurfaceAppRoutes`, `TestDaemonDependencyAllowlist`,
   `TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots`, `TestProductionGoSurface`,
   `TestEveryVerbHasHelp`, the quickstart tests, `TestDaemonCredentialWriterSites`,
   `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`, PG2's static banned-token test). Search for any
   test that enumerates embedded files, keys tables or QUICKSTART sections before assuming none.
6. The §6 gate script for the controller (one bash shell, `tee` every leg to `$HOME/.ailang/state/world-iter387/gates/`,
   never only through `tail`), incl. `-race` on touched Go packages, the inventories, the live-Chrome leg,
   `./scripts/verify_ail.sh`, and `VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh`
   after commit. Measure each gate at the base first. **Never start a gate under `nohup`**.
7. Executor notes as the sprint-3 plan's §7: no git writes; `patch -p1` is the only writer; runs of ~2–4
   sub-milestones; cumulative snapshots to `$HOME/.ailang/state/world-iter387/snap/<M>/` with MANIFEST.sha256
   (exclude the two plan files); stop conditions; numstat checks; exact acceptance commands per sub-milestone;
   **an explicit "never `nohup`, never background a test leg" line**.
8. Acceptance criteria (§9-style), numbered, that the judge can score against. Include: A3a is then code-complete
   except Mark's attended exit run (QUICKSTART §17); say whether §17 needs a text change for the notification step
   and ship it if so.
9. **Every premise gets a §1.2 row**: any sentence saying "the code currently does X" cites the command and its
   observed output; an empty result needs a known-positive control in the same call. Where the design and HEAD
   disagree, measure, route minimally, record it in §1.4, and list every file you touch outside the design's §8.
10. Close with this plan's own Appendix C: anything left for A4 or later, and the attended exit-run steps.

## Hard constraints
- D-WORLD-109: test keys, temp stores, `approvaltest` passkeys/keys and virtual authenticators only. No real key,
  no `~/.ssh`, no `~/.ailang/world`, no real `/dev/tty`, no real browser profile, no mint or publish against a real
  store, no real network listener beyond `httptest`/loopback in tests. Never run bare `ssh-keygen`. Never trigger a
  real desktop notification outside a headless test profile.
- `host/approval` stays pure and standard-library-only. Effects only at the store/CLI/daemon boundary. `surfaceCSP`
  is not changed. The daemon never KEEPS the credentials writer. No CORS headers anywhere. `app.js`/`approve.js`: no
  `innerHTML` with unescaped data, no `eval(`, no `document.cookie`; follow the file's existing `esc()` convention.
- Do not touch: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/` (except reading),
  `cmd/world-publish/*`, `~/.ailang/state/mission-v1*`, the `/surface/passkey` route, `passkey.js`, the A2 law.
- Export `AILANG_BIN=$HOME/.pinned-ailang/ailang` in every shell that runs Go tests.
- Your Bash tool's shell is zsh: run loops/globs/`$VAR` work under `bash -c` or `bash <<'EOF'`.
- Known load-sensitive / flaky tests (re-run once alone, do not chase): `TestSeedRestartIgnoresASurvivingOldCopy`
  (host/broker), `TestHostsLaunchersExecHost`, host/broker exec-probe timeouts incl.
  `TestExecProtectedCheckTimeout/reserve`, local `-race` of host/broker, host/daemon, host/verifygate (incl.
  `TestSeToolsGateArmsAreLoadBearing/plan-bytes`), `TestMCPPostBudgetProductionConstants/batch` (host/daemon),
  `TestBackoff` (host/surfacefollow), `TestDemoSeedPreflightRefusals`, `TestWorktreeSnapshotRefusesAConfiguredFilter`,
  `TestMissionConfigDefaultFlowIntegration/invalid_ledger`, `scripts/marketingcheck` supervisor tests,
  `TestMarketingCaptureFakeChromeSixShots/inventory`.
- **Time budget: aim to finish within ~75 minutes (this is a small sprint); hard stop at 120.** The iteration has a
  6 h wall clock that started 08:21Z, and the executor, gates, judge, CI and record still follow you.

## Waiting discipline (binding on you)
If you background anything (a long test run, a mutant loop, Chrome), poll its ARTIFACT (an output/rc file it writes
on completion) inside a bounded `date +%s` loop in the same turn. A Monitor is an event stream, not a wait; nothing
will re-invoke you when it fires. Never end your turn announcing a wait, a pause, or "once it completes" — end only
when both files are written and the round trip has passed.

## Your final message
Report: the plan path, scope (what is in, anything moved and why, with measured numbers), the sub-milestone list with
expected numstat after each, mutant counts (total / killed / named equivalents) and the ledger path, the A02/A03/A05
kill proofs, the live-Chrome results per subtest, the round-trip result, base-gate results, every premise found false
(§1.4), and anything the controller must decide.
