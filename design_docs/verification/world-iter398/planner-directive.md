MISSION-ROLE: planner

You are the sprint-planner for Ailang World mission iteration 398 (unattended run; no human is present).
Write the sprint plan for **row 243 `w-world-address-bar-verbs-v2`, sprint V2b — "the bar" (design M4–M6, plus §2.2
moved from V2a)**. No re-design: `design_docs/planned/w-world-address-bar-verbs-v2.md` r1 is reviewed (quorum r2
PROCEED 3/3, no absentee; independent audit PASS 84; iteration 396) and is NOT reopened. V2a (the seam) LANDED at
iteration 397 (PR #585, merge `2d530cbe`, judge PASS 90): the walk is `host/why` behind a `Getter`
(`why.Run`, `why.ParseTarget`, `why.Head`, `why.Elide`, `why.Sum`, `why.MaxResponseBytes`, …). V2b adds the caller:
`why <target>` answered in the World app's address bar, in-process, read-only.

## Scope (read these verbatim: design §0.3 A1–A5, §0.4, §2.2, §2.3, §2.4, §2.5, §2.6, §2.7, §3, §4, §5, §6 "V2b — the
bar", §7, §8 rows M4–M6, §9, §10, §13.2; and the V2a plan `design_docs/planned/w-world-address-bar-verbs-v2a-sprint-plan.md`
Appendix C "the V2b handoff")
- **§2.2 (moved here from V2a by V2a-D2, measured 111 lines by a compiled probe)**: `verbs.WhyFlags`/`WhyArgs`,
  `WhyLeadingTarget`, `WhyTrailingTarget`, `WhyScanDefault`/`WhyScanMax`/`WhyScanOK`, the CLI's use of them
  (`cmd/ailang-worldd/why_cli.go`), VA1 and `verbs.WhyScan* == why.Scan*`. The V2a planner's probe patch is at
  `$HOME/.ailang/state/world-iter397/plan-scratch/patches/v2b-probe-2.2.patch` (and `plan-scratch/v2b-probe/`) —
  a starting point you must re-measure against HEAD, not a deliverable. **The CLI's `why` output, exits and wording
  stay byte-identical** (design RV3; V2a's VC1 golden `cmd/ailang-worldd/testdata/` must pass unedited).
- **M4**: `verbs.AdmitWhy` (every argument form, §2.3; `-`, `--result`, `--expect-store` refused with the CLI line
  BEFORE `why.ParseTarget`), the `why` row, VA2–VA3, VB6–VB7.
- **M5**: the daemon's in-process `Getter` (§2.4: private mux over the seven `/v1` GET handlers, the 10 s read
  deadline, `why.MaxResponseBytes`), `verbWhy`, the route's `why` branch and outcome precedence (§2.5), the A2 cap
  (1 MiB); VR1–VR8, VS2.
- **M6**: VE1–VE3, the page (A2–A4, §2.6) incl. the V1 judge notes the design disposes into V2b (§9: A4 + VP7 empty
  submit supersedes a pending ask; VB6; VP8 third-submit supersession; VP10 `result` without `cli`; VB7 `MaxWords`/
  `MaxLine`; A3 + VP11 refusals-with-cli open the panel; Escape closes the panel in XV2), live Chrome XV2, and
  QUICKSTART §19's `why` addition.

## Inputs you must resolve BY MEASUREMENT (record each in a §1.3/§1.4 row with the command and its output)
1. **Audit input 4**: VS2/T1's pinned `d.*` set must include `whyGetter`, `whyMux` (or whatever you name them), the
   seven handler method values the private mux selects, and any `sync.Once` field. VS1's `ast.Inspect` + `DeepEqual`
   shape in `host/daemon/surface_verb_test.go` fails loudly on an incomplete set — run it. State that T1's read-only
   argument also rests on VR6/VR8.
2. **Audit input 5**: `host/verbs/lookup.go` sets `Reason: "why reads World, but its walk runs in the CLI. The app
   cannot run it yet."` (re-measure the line) and the route copies it to `a.Reason`. When `why` runs, the reason must
   be cleared or replaced (A1), with a test that fails on the stale text.
3. **Audit input 6**: the 1 MiB cap is on raw bytes; the JSON-encoded body can be larger (`\u00XX` escapes) and is
   still bounded — say so in the plan's RV2 text and the route's comment; test the bound you claim.
4. **V2a judge's non-blocking notes** (log 397): (1) six base lines no test pins, unchanged by the move —
   `--scan 0`, `why head` with no head, the `rest:` receipt state operand, the fallback-plan operand, the head-absent
   error, a mid-rune `Elide` cut. For each: pin it where V2b's bar or route reaches it (design VA2's `--scan`
   out-of-range row drives `W-SCANMAX`; VE1 runs a `rest:` receipt through the route — RS-REST), or name it a residual
   with the reason. (2) the `--json` trailing newline and the `--result` error text rest on VC1's golden alone and
   `TestWhyRunReplaysTheCLI` is partly circular on them — the bar's `why --json` answer must be pinned by a test that
   does not derive its expectation from `why.Run` itself. (4) VW5 cannot tell a closure that still calls `why.Head`
   from the real delegation — note whether V2b changes that; no action required if not.
5. **V2a's named survivors** `W-LOGPAGE`, `W-MAXRESP`, `W-SCANMAX` (V2a plan §4/Appendix C): each names a V2b
   behavioural kill (VR5 for `MaxResponseBytes`, VA2 for `--scan` bounds). Run those V2a mutants against your final
   tree and record whether V2b's tests now kill them.
6. Every design line you rely on was measured at `200ea5d5`; V2a then moved code. Re-measure every file:line the
   design cites for M4–M6 at HEAD `c09e0cdb` (e.g. `surface_verb.go`, `lookup.go`, `app.js` line numbers have
   likely shifted or not; say which).

## Where you work
- Sprint worktree (already created, branch `sprint/world-iter398-row243-v2b`, base = origin/dev
  `c09e0cdb3c8a3eb25c208dbacedddfb690fe8ca0`): `/Users/voightkampff/dev/sunholo-data/.wt-world-iter398-row243-v2b`
- VERIFIED BY CONTROLLER: `git diff --stat 2d530cbe c09e0cdb -- cmd host docs` is EMPTY (only the record changed since
  V2a's merge).
- Pinned ailang: `$HOME/.pinned-ailang/ailang` = AILANG v0.41.0 (the Go suite's `AILANG_BIN`). V2b changes no `.ail`.
- **Real-stack why tests skip in CI** (`WORLD_TOOL_AILANG_BIN` unset). Any `why` change must also run
  `TestDemoSeedWhyNamedLinks/real-stack` and `TestCLIAgainstSeToolsDaemon` with
  `WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang` — put them in the §6 gate script as a named leg
  (`--- PASS`, not SKIP), measured at the base first.
- Read the newest landed plans for the FORMAT and reproduce it: `design_docs/planned/w-world-address-bar-verbs-v2a-sprint-plan.md`
  (iteration 397) and `design_docs/planned/w-world-address-bar-verbs-v1-sprint-plan.md` (iteration 395 — the one
  with page tests, live Chrome XV1 and QUICKSTART §19; reuse its Chrome/page-test mechanism and its env gates
  e.g. `WORLD_REQUIRE_PASSKEY_BROWSER`/`WORLD_CHROME`, with a temp profile only) and their JSONs.
- Scratch (compile your stages, run mutants) OUTSIDE the repo: `$HOME/.ailang/state/world-iter398/plan-scratch/`
  (exists; never under /tmp — location-sensitive tests go red there).
- Write exactly two files in the worktree, and DO NOT git add/commit them (the controller commits):
  - `design_docs/planned/w-world-address-bar-verbs-v2b-sprint-plan.md`
  - `design_docs/planned/sprint_w-world-address-bar-verbs-v2b.json`
- Leave every other file in the worktree untouched (the executor applies your patches there). No git writes.

## Patches and the snapshot contract (as V2a)
- The executor's only writer is `/usr/bin/patch -p1 -E` (macOS; without `-E` a deleted file stays as 0 bytes).
- Executor snapshots: cumulative, `$HOME/.ailang/state/world-iter398/snap/<M>/` holds every file added or changed
  since base (repo-relative paths), a `DELETED` file (one repo path per line, may be empty), a `STATUS` file, and
  `MANIFEST.sha256` (exclude the two plan files, STATUS, DELETED). The controller's commit script copies the snapshot
  over the worktree and `git rm`s each `DELETED` path. Give the exact snapshot commands in §7, and make the `find`
  robust to a sub-milestone with no `docs/` or no `host/…` dir:
  `find $(for d in cmd host docs; do [ -d "$d" ] && echo "$d"; done) …` — and RUN your snapshot line on such a
  sub-milestone in scratch before shipping it.

## Size rule
- **Cap: ≤ 1,400 changed lines** for V2b, measured `git diff -M --numstat <base> HEAD -- cmd host docs` (state this
  exact filter). Budget ≈ 1,041 (design ≈ 930 + §2.2's measured 111; × 1.09 ≈ 1,135). If over, propose a cut that
  leaves the bar's `why` FULLY proven end to end (route + page + one live Chrome drill) and moves the rest (e.g. a V1
  judge-note test that does not involve `why`) to a named V2c with its measured size; record it as a §1.3 decision.
  Prefer a smaller, fully proven sprint over a full, half-proven one.
- Sub-milestones of ≤ ~250–400 changed lines each, each independently green and safe at commit (your call, measured).

## Retro directives (iterations 374–397; each was paid for)
1. **Ship compiled patches**: every line of code, test and doc as a unified diff in Appendix P
   (```` ````diff <M> ```` fences), generated LAST from compiled, vetted, gofmt-clean, green scratch stages; then
   extract them FROM THE PLAN FILE with §2's awk, apply to a fresh `git archive c09e0cdb` export, and `diff -r`
   against your work tree — byte-identical or regenerate. The controller re-does this round trip independently.
   **Ship the patches BEFORE running the full mutant set, and checkpoint the plan file** (write a first complete
   version, then append the mutant ledger).
2. **Expected numstat after each sub-milestone's commit** (measured after commit, `-- cmd host docs`, renames on).
3. **Mutants anchored to the DIFF**: one mutant per new or changed production branch, enumerated from the diff, not
   from the test plan — in particular: every `AdmitWhy` arm (each refused form and its CLI line, the accepted forms),
   the getter's mux selection, its deadline, its cap and its error classes (404 vs read failure vs over-cap), the
   route's outcome precedence (§2.5 a–d), the stale-reason clearing, the A2 cap's comparison, every page branch
   (A2–A4: supersession, the panel for refusals with `cli`, `result` requires `cli`), every exported constant, the
   `WhyFlags` registrations and defaults. Run each against the named test row AND the package; record
   killed/survived; name equivalents with reasons. A mutant that does not compile is not a mutant (pair with a killed
   control). **Name the test that killed each mutant; a mutant killed only by an inventory/static/AST test is not
   behavioural proof — re-run with that test excluded and record both.** Write the ledger to a file in scratch; every
   re-run appends its own line. **The plan's mutant table ids must be the ledger's.** JS page mutants count too (as
   V1's plan did).
4. **Trace every input AND every failure/timeout class through each new branch** (394/396 retro): §2.7's table is the
   design's; for each row name the test that drives it end to end at the route and, where visible, at the page.
5. **Repo-wide inventories**: run every design §7 inventory and any test that enumerates packages, imports, exported
   symbols, routes, files, CLI verbs or context roots, by name with `-v`, on your final scratch tree (incl.
   `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots`,
   `TestProductionGoSurface`, `TestEveryVerbHasHelp`, `TestFrozenV1RouteTableMatchesMux`, `TestSurfaceAppPure`,
   VS1/VS2 in `host/daemon/surface_verb_test.go`, the page golden `host/surfacepage/testdata/verb-route.golden.json`).
   Search for any such test before assuming none. If the diff adds a production context root,
   `go test ./host/store/ -run '^TestProductionContextRoots$'` is a named gate (charter Repo Profile).
6. **Fixtures from the real tool**: every golden or recorded body (the page golden's `why` row, VR4's transport-parity
   recording) is GENERATED by a named command against the real route/daemon, never typed. Say the command in the plan.
7. The §6 gate script for the controller (one bash shell, `tee` every leg to
   `$HOME/.ailang/state/world-iter398/gates/<leg>.txt` — `.gitignore` has `*.out`, so `.txt` — never only through
   `tail`), incl. `go test -run '^$' ./...` as the test-file compile fence, `-race` on touched Go packages, the
   inventories, the real-stack why leg, the live-Chrome leg (as V1's), `./scripts/verify_ail.sh`, and
   `VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh` after commit. Measure each gate
   at the base first. Honour `GATES_SKIP_VERIFY_GO=1` (the executor sets it; the controller runs `verify_go`). **Never
   start a gate under `nohup`.** Write `ALL_DONE.<sha>` with `FAILS: …` at the end.
8. Executor notes as the V2a plan's §7: no git writes; `patch -p1 -E` is the only writer; runs of ~2–4
   sub-milestones; cumulative snapshots per the contract above with MANIFEST.sha256; stop conditions; numstat checks;
   exact acceptance commands per sub-milestone (each fitting a 600 s foreground Bash call — split if not); **an
   explicit "never `nohup`, never background a test leg; pass a 600000 ms timeout to every long Bash call" line**.
9. Acceptance criteria (§9-style), numbered, that the judge can score against — including CLI behaviour-neutrality
   (VC1 unedited and green), the read-only argument (VS1/VS2, no new `/v1` route, no POST, no store write), and the
   exact list of pre-existing test files changed with the reason for each.
10. **Every premise gets a §1.2 row**: any sentence saying "the code currently does X" cites the command and its
    observed output; an empty result needs a known-positive control in the same call. Where the design and HEAD
    disagree, measure, route minimally, record it in §1.4, and list every file you touch outside the design's §7.
    **For every type, field or function your patches name, cite where it exists at HEAD or mark it new.**
11. Close with Appendix C: anything moved out (with size), residuals carried (RV1–RV5, RS-FOCUS, RS-REST, the V2a
    judge's unpinned lines you did not pin), and what row 243 V3 (the verb-approval law) inherits.

## Hard constraints
- D-WORLD-109: test stores, temp dirs and test keys only. No real key, no `~/.ssh`, no `~/.ailang/world`, no real
  `/dev/tty`, no real browser profile (temp Chrome profile only), no mint or publish, no real network listener beyond
  `httptest`/loopback in tests.
- Read-only: no new `/v1` route, no POST, no store write, no process spawn from the daemon, no `.ail`. The bar does
  not run `why --result`, `why -` or `--expect-store` (refused with the CLI line). The CLI's own output is unchanged.
- `host/why` keeps effects out of its core (VS2a). The getter's effects live in `host/daemon`.
- Do not touch: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/` (except reading),
  `cmd/world-publish/*`, `~/.ailang/state/mission-v1*`, row 201's approval-card files beyond what A3/A4 need (list
  any you touch, with the reason).
- Export `AILANG_BIN=$HOME/.pinned-ailang/ailang` in every shell that runs Go tests.
- Your Bash tool's shell is zsh: run loops/globs/`$VAR` work under `bash -c` or `bash <<'EOF'`. macOS has no
  `timeout` command. Pass `timeout: 600000` on any Bash call that may run long; the default is 120 s.
- Known load-sensitive / flaky tests (re-run once alone, do not chase): `TestSeedRestartIgnoresASurvivingOldCopy`
  (host/broker), `TestHostsLaunchersExecHost`, host/broker exec-probe timeouts incl.
  `TestExecProtectedCheckTimeout/reserve`, `TestPublishApprovalRefusalSetWithALandedPositiveControl`, local `-race` of
  host/broker, host/daemon, host/verifygate, `TestMCPPostBudgetProductionConstants/batch`, `TestBackoff`
  (host/surfacefollow), `TestDemoSeedPreflightRefusals`, `TestWorktreeSnapshotRefusesAConfiguredFilter`,
  `TestMissionConfigDefaultFlowIntegration` (host/verifygate), `TestCLIRealSubprocessEpisode` (cmd, under load),
  `TestCapsuleTemplateIsCopiedNotShared`, `scripts/marketingcheck` supervisor tests,
  `TestMarketingChromeSixShots` (red without `WORLD_TOOL_AILANG_BIN`).
- **Time budget: aim to finish within ~90 minutes; hard stop at 130.** The iteration's 6 h wall clock started
  01:39Z, and the executor, gates, judge, CI and record still follow you. If you run short, take a cut fully proven
  rather than everything half-proven.

## Waiting discipline (binding on you)
If you background anything (a long test run, a mutant loop), poll its ARTIFACT (an output/rc file it writes on
completion) inside a bounded `date +%s` loop in the same turn. A Monitor is an event stream, not a wait; nothing will
re-invoke you when it fires. Never end your turn announcing a wait, a pause, or "once it completes" — end only when
both files are written and the round trip has passed.

## Your final message
Report: the plan path, scope (what is in, anything moved out and why, with measured numbers), the sub-milestone list
with expected numstat after each, mutant counts (total / killed / named equivalents) and the ledger path, the
round-trip result, base-gate results (incl. the real-stack and Chrome legs), every premise found false (§1.4), how you
resolved inputs 1–6, the exact pre-existing test files changed, and anything the controller must decide.
