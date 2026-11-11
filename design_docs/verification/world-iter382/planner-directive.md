MISSION-ROLE: planner

You are the sprint-planner for Ailang World mission iteration 382 (unattended run; no human is present).
Write the sprint plan for **row 201 `w-world-approvals-by-key`, milestone A2, SPRINT 4** — M7 "the page,
the live browser and QUICKSTART", exactly the work the sprint-3 plan's **Appendix C** moved here by the
design's §9 fallback. No re-design: the A2 design r2 is ratified (quorum + carve-out) and NOT reopened.

## Where you work
- Sprint worktree (already created, branch `sprint/world-iter382-row201-a2s4`, base = origin/dev
  `27ce163eb1ae8af55151ea6548591c15a3126835`, which contains A2 sprints 1–3 merged, sprint 3 = `0ce117c5`):
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter382-row201-a2s4`
- Scratch (compile your stages, run mutants and browsers) OUTSIDE the repo:
  `$HOME/.ailang/state/world-iter382/plan-scratch/` (never under /tmp — location-sensitive tests go red there).
- Write exactly two files in the worktree, and DO NOT git add/commit them (the controller commits):
  - `design_docs/planned/w-world-approvals-by-key-a2-sprint4-plan.md`
  - `design_docs/planned/sprint_w-world-approvals-by-key-a2-s4.json` (same shape as `sprint_w-world-approvals-by-key-a2-s3.json`)
- Leave every code file in the worktree untouched (the executor applies your patches there). No git writes of any kind in the worktree.

## Read first
1. `design_docs/planned/w-world-approvals-by-key-a2-sprint3-plan.md` **Appendix C** (the sprint-4 scope, the
   inherited link format, XB1/XB2 seams, file names) and §1.2 C17/C18, §2 (patch protocol + awk
   extraction), §6 (gate script), §7 (executor notes), Appendix P's fence format. Reproduce this plan's
   format — it worked five sprints running (haiku executor, zero deviations, numstat equal at every step).
2. `design_docs/planned/w-world-approvals-by-key-a2.md` (r2): §2.7 (the JSON the page builds), §2.8
   (enrolment code grouping), §2.14 (the page, incl. the [r2] AbortController/240000 ms timer rule and
   V31 `blob:` save), §2.15 (generate.sh + the recording driver), §6 "Sprint 3" rows PG1/XB1/XB2 and the
   edited/unedited lists, §7 (repo-wide inventories), §8 conflict surface, §9, §10 (the attended run —
   NOT run by the loop; QUICKSTART §16 documents it).
3. The production link builder you must render against: `cmd/ailang-worldd/approver_passkey.go`
   (`ceremonyLink`, `writePasskeyLinks`) and the daemon's surface routes/CSP (`host/daemon`,
   `setSurfaceHeaders`, `surfaceCSP`, `TestSurfaceAppRoutes`, `TestFrozenV1RouteTableMatchesMux`).
4. `CLAUDE.md`, `design_docs/coding-standards.md` (binding), `docs/QUICKSTART.md` §15 (format for §16).

## Browser material (already preserved for you, outside /tmp)
- `$HOME/.ailang/state/world-iter382/plan-scratch/vt_probe.mjs` — sprint 3's measured XB2 virtual-time
  probe (`Emulation.setVirtualTimePolicy` advance budget 250000 → page time exactly 240000 ms in ~55 ms).
- `$HOME/.ailang/state/world-iter382/plan-scratch/w379-probes/browser/probe.mjs` (sha256 `01a17c0d…37cf`,
  the iteration-379 recording probe: WebAuthn.enable, addVirtualAuthenticator ctap2_1/internal/rk/UV,
  create + gets) and `.../w379-probes/blob/probe.mjs` (the V31 blob-save probe).
- Chrome: `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`; node: `/opt/homebrew/bin/node` (v26).
  **Temp profile, `--headless=new`, `--remote-debugging-port=0`, CDP virtual authenticator only** — never
  the operator's real profile, keychain, iCloud passkeys, or a platform authenticator.

## Scope
- **M7** as Appendix C lists: `host/surfacepage/passkey/{index.html,passkey.js}` + `host/surfacepage/passkey.go`
  (`go:embed`), served by `host/daemon` at `GET /surface/passkey` and `GET /surface/passkey.js` with
  `setSurfaceHeaders` and the **unchanged** `surfaceCSP`; no store, no request, no `fetch`/`XMLHttpRequest`/
  form `action`, `textContent` only; the [r2] AbortController + fixed 240000 ms timer, `signal` beside
  `publicKey`, button disabled while pending, a late result ignored, saving disabled after expiry; the
  §2.7 JSON built by the page; the grouped enrolment code; the `blob:` Save link. Tests: PG1
  `TestPasskeyPageRoutes`, the `TestSurfaceAppRoutes` edit, **XB1** `TestPasskeyPageInChrome`, **XB2**
  `TestPasskeyPageAbortsAPendingCeremony` (gated: `WORLD_CHROME` + `node` on PATH; `t.Skip("SKIP COUNTED: …")`
  unless `WORLD_REQUIRE_PASSKEY_BROWSER=1`, then `t.Fatal`), `testdata/passkey_cdp.mjs`, and
  `host/approval/testdata/webauthn/generate.sh` made real (today it is a stub that exits 1).
- **XB1 must build its links with the production code** (through the CLI or `ceremonyLink.link()` from a
  `cmd` test), and must carry the page's output back through the production verifier: the enrolment
  registration JSON the page saves must be accepted by the real `approver enroll --format webauthn` path
  (or the exact function it calls), and the assertion JSON must mint through the real
  `approve --webauthn-assertion` path (or `ApproveByAssertion` + the consume it uses) against a temp store.
  Rule 3k: if the product hands a human something to run or open, a test must exercise exactly that.
- **QUICKSTART §16** "Approve with a passkey" *(attended — pending first verbatim run)*. Written, not run.
- Optional carry-forwards (cheap, take them if they fit): the judge-381 early `--name` check
  (`approval.ApproverNamePattern` beside WN2, before any link prints; one branch + one CW1 arm) and pinning
  the next-step line's `--db`/`--credentials-db` order (one assertion).
- **Size cap: ≤ 2,000 changed lines**, measured `git diff --numstat <base> HEAD -- cmd host` (state this exact
  filter; note `docs/` and `design_docs/` are outside it). Split into sub-milestones of ≤ ~250 changed lines
  each where possible (M7a/M7b/…), each independently green and safe at commit (no route that serves a page
  whose tests are absent; no link pointing at a page that is not served).
- If XB1/XB2 cannot be made reliably green on this machine (≥ 3 consecutive green runs each with
  `WORLD_REQUIRE_PASSKEY_BROWSER=1`, timings recorded), say so with measurements; a flaky browser test is a
  reason to narrow, never to loosen its assertions. Note the machine is heavily loaded by other sessions
  (load 10–38) — give browser-test deadlines derived from measured runs with headroom (rule 3m).

## Retro directives (iterations 374–381; each was paid for)
1. **Ship compiled patches**: every line of code, test, page asset, `.mjs`, `generate.sh` and doc as a unified
   diff in Appendix P (```` ````diff <M> ```` fences, as sprint 3), generated LAST from compiled, vetted,
   gofmt-clean, green scratch stages; then extract them FROM THE PLAN FILE with §2's awk, apply to a fresh
   `git archive 27ce163e` export, and `diff -r` against your mutant tree — byte-identical or regenerate. The
   controller re-does this round trip independently. Embedded assets and `.mjs`/`.sh` files must survive
   byte-for-byte (trailing newlines, tabs; and the executable bit of `generate.sh` — say how it is set, since
   `patch` may not carry modes; if the mode matters, have a test or the gate check it).
2. **Expected numstat after each sub-milestone's commit** (measured after commit, `-- cmd host`).
3. **Mutants anchored to the DIFF**: one mutant per new production branch (PG1's method/path/CSP/content-type/
   asset branches; the page's origin/rp/op checks, abort, timer, late-result, save-after-expiry, button-disable
   and JSON-build branches; generate.sh's refusal paths), enumerated from the diff, not from the test plan;
   run each against the named test row AND the package (browser mutants with `WORLD_REQUIRE_PASSKEY_BROWSER=1`);
   record killed/survived; name equivalents with reasons; base controls where a test pins pre-existing
   behaviour. Write the mutant ledger to a file in scratch; every re-run (incl. a COMPILEERR re-run) appends
   its own line — a re-run in prose only is not evidence.
4. **Repo-wide inventories**: run every design §7 inventory by name with `-v` on your final scratch tree
   (`TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` (AC10 — confirm it skips `_test.go`),
   `TestDaemonDependencyAllowlist`, `TestFrozenV1RouteTableMatchesMux` (stays 13), `TestSurfaceAppRoutes`,
   `TestDaemonConstructorCallerInventory`, `TestWorldBoundaryDependencyAllowlist`,
   `TestEveryCommandDirectoryIsAProtectedGroup`, `TestProductionContextRoots`, `TestProductionGoSurface`,
   `TestNoNewDeadlineFreeStoreReads`, `TestCredentialsMethodsRequireADeadline`,
   `TestNoParallelWriteForkPackagesOutsideVerifygate`, `TestEveryVerbHasHelp`, the quickstart tests, and any
   tracked-binary / embedded-asset / docs-link checks a new `host/surfacepage/passkey/` dir or a QUICKSTART
   section would trip). Search for tests that inventory `host/surfacepage` or QUICKSTART sections.
5. The §6 gate script for the controller (one bash shell, `tee` every leg to
   `$HOME/.ailang/state/world-iter382/gates/`, never only through `tail`), incl. `-race` on touched packages,
   the inventories, XB1/XB2 with `WORLD_REQUIRE_PASSKEY_BROWSER=1`, `./scripts/verify_ail.sh`, and
   `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` after commit. Measure each gate at the base first. **Never
   start a gate under `nohup`**.
6. Executor notes as sprint 3 §7: no git writes; `patch -p1` is the only writer; runs of ~2 sub-milestones;
   cumulative snapshots to `$HOME/.ailang/state/world-iter382/snap/<M>/` with MANIFEST.sha256 (exclude the two
   plan files); stop conditions; numstat checks. The executor is a small model and must NOT run Chrome/node
   unless your §7 says exactly which command to run — prefer that the executor runs only non-browser
   acceptance and the controller runs XB1/XB2 in §6.
7. Acceptance criteria (§9-style), numbered, that the judge can score against.

## Hard constraints
- D-WORLD-109: test keys, temp stores, in-memory/test authenticators and a CDP **virtual** authenticator only.
  No real key, no `~/.ssh`, no `~/.ailang/world`, no real `/dev/tty`, no real browser profile, no mint or
  publish against a real store, no `tailscale serve`. Never run bare `ssh-keygen`. QUICKSTART §16 is written,
  tagged *attended — pending first verbatim run*, and NOT executed by you.
- `host/approval` stays pure and standard-library-only. Effects only at the store/CLI/daemon boundary.
  `surfaceCSP` is not changed.
- Do not touch: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/` (except reading),
  `cmd/world-publish/*`, `~/.ailang/state/mission-v1*`.
- Export `AILANG_BIN=$HOME/.pinned-ailang/ailang` in every shell that runs Go tests.
- Your Bash tool's shell is zsh: run loops/globs/`$VAR` work under `bash -c` or `bash <<'EOF'`.
- Known load-sensitive / flaky tests (re-run once alone, do not chase): `TestSeedRestartIgnoresASurvivingOldCopy`
  (host/broker), `TestHostsLaunchersExecHost`, host/broker exec-probe timeouts, local `-race` of host/broker,
  host/daemon, host/verifygate, `TestDemoSeedPreflightRefusals`, `TestWorktreeSnapshotRefusesAConfiguredFilter`,
  `TestMissionConfigDefaultFlowIntegration/invalid_ledger`, `scripts/marketingcheck` supervisor tests.
- Premises: any sentence in your plan that says "the code currently does X" needs a §1.2 row with the command
  and its observed output; an empty result needs a known-positive control in the same call. Where Appendix C,
  the design and HEAD disagree, measure, route minimally, record it in §1.4, and list every file you touch
  outside the design's §8 conflict surface.
- **Time budget: aim to finish within ~120 minutes.** The iteration has a 6 h wall clock and the executor,
  gates, judge and CI still follow you. Prefer a smaller, fully-proven sprint over a larger one.

## Waiting discipline (binding on you)
If you background anything (a long test run, a mutant loop, a browser run), poll its ARTIFACT (an output/rc
file it writes on completion) inside a bounded `date +%s` loop in the same turn. A Monitor is an event stream,
not a wait; nothing will re-invoke you when it fires. Never end your turn announcing a wait, a pause, or "once
it completes" — end only when both files are written and the round trip has passed.

## Your final message
Report: the plan path, scope (what is in, anything narrowed and why, with measured numbers), the sub-milestone
list with expected numstat after each, mutant counts (total / killed / named equivalents) and the ledger path,
the round-trip result, base-gate results, XB1/XB2 timings (≥ 3 runs each), every premise found false (§1.4),
which optional carry-forwards you took, and anything the controller must decide.
