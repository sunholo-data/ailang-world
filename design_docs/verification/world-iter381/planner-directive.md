MISSION-ROLE: planner

You are the sprint-planner for Ailang World mission iteration 381 (unattended run; no human is present).
Write the sprint plan for **row 201 `w-world-approvals-by-key`, milestone A2, SPRINT 3** — "the operator paths and the page" (M6 the CLI, M7 the ceremony page + live browser test + QUICKSTART §16).

## Where you work
- Sprint worktree (already created, branch `sprint/world-iter381-row201-a2s3`, base = origin/dev
  `ab4b2402116243c0e0fab9ed5cee8ec80a660a4e`, which contains A2 sprint 1 (`d5d90e49`) and sprint 2 (`3e55451f`) merged):
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter381-row201-a2s3`
- Scratch (compile your stages, run mutants) OUTSIDE the repo: `$HOME/.ailang/state/world-iter381/plan-scratch/`
  (never under /tmp — location-sensitive tests go red there).
- Write exactly two files in the worktree, and DO NOT git add/commit them (the controller commits):
  - `design_docs/planned/w-world-approvals-by-key-a2-sprint3-plan.md`
  - `design_docs/planned/sprint_w-world-approvals-by-key-a2-s3.json` (same shape as `sprint_w-world-approvals-by-key-a2-s2.json`)
- Leave every code file in the worktree untouched (the executor applies your patches there). No git writes of any kind in the worktree.

## Read first
1. `design_docs/planned/w-world-approvals-by-key-a2.md` (r2, ratified via quorum + carve-out; NOT reopened): §0, §2.7, §2.8 (enrolment code), §2.12 (+2.12a, 2.12b), §2.13, §2.14 (the page, incl. the [r2] AbortController/240000 ms timer rule), §2.15 (generate.sh), §3, §4, §6 "Sprint 3" (CW1–CW9, PG1, XB1, XB2, edited/unedited lists), §7 (repo-wide inventories), §8, §9 (**sprint 3 = M6–M7, and its fallback cut: M7 becomes its own fourth sprint**), §10 (the attended run — NOT run by the loop), §11–§13.
2. `design_docs/planned/w-world-approvals-by-key-a2-sprint2-plan.md` — the format to reproduce (it worked: a haiku executor applied 8 compiled patches across four runs with zero deviations, numstat equal at every step). Especially §1 (premises table with commands), §1.5 (size and cut), §2 (patch protocol and the awk extraction), §6 (gate script), §7 executor notes, Appendix P, and any Appendix C that names sprint-3 work.
3. `design_docs/planned/w-world-approvals-by-key-a1b-sprint2-plan.md`/`-sprint3-plan.md`/`-sprint4-plan.md` only as needed for the A1b CLI rigs (`spyEnv`, `approveRig`, `noTerminalEnv`, `stubDaemon`, `confirmApproverAct`, `readRegularFileCapped`).
4. `CLAUDE.md`, `design_docs/coding-standards.md` (binding).

## Scope, order and the cut
- **M6 first** (the CLI: §2.12 enrol with WN1–WN11 and the §2.12a behaviour-preserving terminal refactor, §2.12b list/revoke labels L4/L5, §2.13 `--webauthn-link`/`--webauthn-assertion` with WF1–WF6, F1 widened to five modes, CW1–CW9 and the named edited tests). Then **M7** (PG1 routes + embedded page, XB1, XB2, `generate.sh`, QUICKSTART §16 *(attended — pending)*).
- **Size cap: ≤ 2,000 changed lines**, measured as `git diff --numstat <base> HEAD -- cmd host` (state this exact path filter in the plan). Sprint 2's real sizes ran well above the design's estimates (M3b 890 vs 450), so MEASURE M6 from compiled stages before deciding. The design's §9 sprint-3 fallback (M7 → its own sprint 4) is a measurement, not a preference: apply it if M6+M7 measured > 2,000, or if XB1/XB2 cannot be made reliably green on this machine within the planning run. If you cut, list the moved items in an Appendix C for sprint 4 exactly as sprint 1 did, and make sure every committed state is safe (e.g. no route that serves a page whose tests are absent). Never exceed 2,000.
- Split into sub-milestones of ≤ ~250 changed lines each where possible (M6a/M6b/…), each independently green, the §2.12a refactor as its own first sub-milestone so CA4/CA9/CA10 are its unedited oracle.

## M7 and the browser (if M7 stays in this sprint)
- Chrome is at `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`; node is `/opt/homebrew/bin/node` (v26). Use a **temp profile, `--headless=new`, `--remote-debugging-port=0`, and a CDP virtual authenticator only** — never the operator's real browser profile, keychain, iCloud passkeys, or a platform authenticator. The design's iteration-379 recording driver/probe (see §2.15 and `host/approval/testdata/webauthn/SOURCE.txt`) is your starting point; find it in the repo or in sprint 1's plan before writing a new one.
- XB1/XB2 are gated (`WORLD_CHROME` + `node` on PATH; `t.Skip("SKIP COUNTED: …")` unless `WORLD_REQUIRE_PASSKEY_BROWSER=1`, then `t.Fatal`). Run both with `WORLD_REQUIRE_PASSKEY_BROWSER=1` at least 3× in a row in your scratch tree and record the timings; a flaky browser test is a reason to cut M7, not to loosen its assertions.
- XB2's virtual-time / injected-clock seam: choose one, measure that it actually advances page time to 240000 ms (a test that waits 4 real minutes is not acceptable), and record the measurement in §1.2.
- The page makes NO network request (no `fetch`, `XMLHttpRequest`, form `action`) and uses `textContent` only; `surfaceCSP` is not changed (§2.14 V31).

## Carry-forwards (log 380 / judge 380)
- **Consider refusing an all-zero challenge** in the enrol path's random-challenge step (WN5 area): `env.random` succeeding but returning 32 zero bytes. Measure what the registration law does with a zero challenge today (sprint 2 measured: accepted by the law, left to the caller). Decide (typed refusal → `exitFatal` + a test arm, or a documented decision in §1.3) and say why.
- Sprint 2's judge noted `ApproveByAssertion` does not guard a zero-valued `Passkey.Key`; the controller found `COSEKey.verify`'s default arm returns `ErrCOSEAlg` (fails closed). No action unless your M6 code path changes that.
- **"The design says this sprint touches no file in X" is a claim to measure** (sprint 2: schema 6 reached `cmd/` test fixtures). List in §1.4 every file outside the design's §8 conflict surface that your patches touch, and why.

## Retro directives (iterations 374–380; each was paid for)
1. **Ship compiled patches**: every line of code, test, page asset, driver script and doc as a unified diff in Appendix P (```` ````diff <M> ```` fences, as sprint 2), generated LAST from compiled, vetted, gofmt-clean, green scratch stages; then extract them FROM THE PLAN FILE with §2's awk, apply to a fresh export of the base (`git archive ab4b2402`), and `diff -r` against your mutant tree — byte-identical or regenerate. The controller will re-do this round trip independently. Embedded assets (`go:embed`) and any `.mjs` must survive the round trip byte-for-byte (watch trailing newlines and tabs).
2. **Expected numstat after each sub-milestone's commit** (measured after commit, `-- cmd host`), as sprint 2 §7.
3. **Mutants anchored to the DIFF**: one mutant per new production branch/refusal (WN1–WN11, WF1–WF6, F1's five modes, L4/L5, PG1's method/CSP/asset branches, the page's origin/rp checks and abort/timer/save branches for XB1/XB2), enumerated from the diff, not from the test plan; run each against the named test row AND the package; record killed/survived; name equivalents with reasons. Base controls where a test pins pre-existing behaviour. Write the mutant ledger to a file in scratch; every re-run (including a COMPILEERR re-run) appends its own line to that file — a re-run that exists only in prose is not evidence (380's judge).
4. **Flag/mode coverage**: for every refusal that tests a combination of independent conditions (F1's exclusive modes: every pair; WN1's pairings: each; CW5's code-shape arms), give an arm for each, so `!a && !b`-style mutants die.
5. **Repo-wide inventories**: run every design §7 inventory by name with `-v` against your final scratch tree (`TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` — note M7 adds a Chrome/node subprocess in a TEST, check whether the AC10 census covers `_test.go` files; `TestDaemonDependencyAllowlist`, `TestFrozenV1RouteTableMatchesMux`, `TestSurfaceAppRoutes`, `TestDaemonConstructorCallerInventory`, `TestWorldBoundaryDependencyAllowlist`, `TestEveryCommandDirectoryIsAProtectedGroup`, `TestProductionContextRoots`, `TestProductionGoSurface`, `TestNoNewDeadlineFreeStoreReads`, `TestCredentialsMethodsRequireADeadline`, `TestNoParallelWriteForkPackagesOutsideVerifygate`, `TestEveryVerbHasHelp`, the quickstart tests). A new context root must appear where `TestProductionContextRoots` demands.
6. The §6 gate script for the controller (one bash shell, `tee` every leg to `$HOME/.ailang/state/world-iter381/gates/`, never only through `tail`), including `-race` on the touched packages (`cmd/ailang-worldd` touched tests, `host/daemon` touched tests), the inventories, XB1/XB2 with `WORLD_REQUIRE_PASSKEY_BROWSER=1` if M7 is in, and `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` after commit. Measure each gate green at the base first (a gate red at base measures the repo). **Never start a gate under `nohup`** (it reds `scripts/marketingcheck`'s `shell-hup` test).
7. Executor notes as sprint 2 §7: no git writes; `patch -p1` is the only writer; runs of ~2 sub-milestones; cumulative snapshots to `$HOME/.ailang/state/world-iter381/snap/<M>/` with MANIFEST.sha256 (exclude the two plan files); stop conditions; numstat checks. Paths must name iter381 and the worktree above.
8. Acceptance criteria (§9-style), numbered, that the judge can score against.

## Hard constraints
- D-WORLD-109: test keys, temp stores, in-memory/test authenticators and a CDP **virtual** authenticator only. No real key, no `~/.ssh`, no `~/.ailang/world`, no real `/dev/tty`, no real browser profile, no mint or publish, no `tailscale serve`. Never run bare `ssh-keygen`. QUICKSTART §16 is written, tagged *attended — pending first verbatim run*, and NOT executed by you.
- `host/approval` stays pure and standard-library-only. Effects only at the store/CLI/daemon boundary.
- Do not touch: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/` (except reading), `cmd/world-publish/*`, `~/.ailang/state/mission-v1*`. The typed-`y` path's messages and exit codes are preserved (§2.12a).
- Export `AILANG_BIN=$HOME/.pinned-ailang/ailang` in every shell that runs Go tests (`host/daemon` tests red by design without it).
- The controller shell is zsh: run loops/globs/`$VAR` work under `bash -c` or `bash <<'EOF'`.
- Known load-sensitive / flaky tests (re-run once alone, do not chase): `TestSeedRestartIgnoresASurvivingOldCopy/survivor` (host/broker, row 259), local `-race` of host/broker, host/daemon, host/verifygate, `TestWorktreeSnapshotRefusesAConfiguredFilter`, `TestMissionConfigDefaultFlowIntegration/invalid_ledger`.
- Premises: any sentence in your plan that says "the code currently does X" needs a §1.2 row with the command and its observed output; an empty result needs a known-positive control in the same call. Where the design and HEAD disagree, measure, route minimally, and record it in §1.4.

## Waiting discipline (binding on you)
If you background anything (a long test run, a mutant loop, a browser run), poll its ARTIFACT (an output/rc file it writes on completion) inside a bounded `date +%s` loop in the same turn. A Monitor is an event stream, not a wait; nothing will re-invoke you when it fires. Never end your turn announcing a wait, a pause, or "once it completes" — end only when both files are written and the round-trip has passed.

## Your final message
Report: the plan path, the cut you chose (what is in this sprint, what moved to sprint 4 and why, with measured numbers), the sub-milestone list with expected numstat after each, mutant counts (total / killed / named equivalents) and the ledger path, the round-trip result, the base-gate results, the browser-test timings (if M7 is in), every premise you found false (§1.4), the all-zero-challenge decision, and anything the controller must decide.
