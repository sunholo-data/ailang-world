# Sprint plan — w-world-approvals-by-key-a2, sprint 4 (iteration 382, row 201)

> **THE SPRINT — M7, "the page, the live browser and QUICKSTART", whole, as the sprint-3 plan's
> Appendix C moved it here.** The passkey ceremony page (`host/surfacepage/passkey/{index.html,
> passkey.js}`, embedded by `host/surfacepage/passkey.go`) served by the production daemon at `GET
> /surface/passkey` and `GET /surface/passkey.js` under `setSurfaceHeaders` and the **unchanged**
> `surfaceCSP`; PG1 `TestPasskeyPageRoutes`; the `TestSurfaceAppRoutes` edit; the Node/DevTools driver
> `cmd/ailang-worldd/testdata/passkey_cdp.mjs`; **XB1** `TestPasskeyPageInChrome` and **XB2**
> `TestPasskeyPageAbortsAPendingCeremony`; `host/approval/testdata/webauthn/generate.sh` made real; and
> QUICKSTART §16 *(attended — pending first verbatim run)*. Both optional carry-forwards are taken
> (M7a: the early `--name` grammar check WN2a, and the next-step line's `--db`/`--credentials-db` order).
> - **Measured: 1,054 changed lines under `cmd/ host/`** — `git diff --numstat <base> HEAD -- cmd
>   host`, insertions + deletions, after a commit of each of the seven scratch stages (206 production: the page 166, the embed 18, the routes and handlers 19, WN2a 3; 848 test, driver and generator).
>   `docs/QUICKSTART.md` (§16, M7g) and `design_docs/` are outside that filter. Margin to the 2,000 cap:
>   946. Nothing is narrowed.
> - **XB1 runs the product's own path end to end (rule 3k).** It enrols through `approver enroll
>   --format webauthn` with the link **that run printed**, opened in Chrome against the page **the
>   production daemon handler serves**; the page saves the registration through its `blob:` Save link
>   into the very path `--registration` named, and the terminal types **the code the page showed**. It
>   then approves through `approve --webauthn-link` → the printed link → the page's saved assertion →
>   `approve --webauthn-assertion`, which mints (`approved_by = webauthn:mark-mac`, the counter
>   advanced). So the page, the CLI's links and the production verifier cannot drift apart unseen.
> - **XB1 and XB2 are green on this machine, 3 consecutive required runs in the §6 script on the final code (XB2 3.3–4.4 s, XB1 3.1–3.8 s, load 7–12), after 7 consecutive green on the code before C34's one-line `syncBuffer` change and 6 more on the code before C23/C24**, each with `WORLD_REQUIRE_PASSKEY_BROWSER=1`
>   (§1.2 C19). XB2 drives the page's production 240000 ms timer in CDP virtual time (paused from before the
>   click): the injected operation sees its signal abort at exactly 240000 ms of page time, 250 s of
>   page time costing 26–98 ms of real time.
> - **Every line of code, test, page asset, `.mjs`, `generate.sh` and doc is a unified diff in Appendix
>   P**, generated LAST from scratch stages that were compiled, vetted, gofmt-clean and green; every
>   mutant ran on the final tree (69 mutants: 64 KILLED by the named test and the suite, 5 named equivalents SURVIVED as predicted). The patches were then extracted **from this file** with
>   §2's `awk`, applied to a fresh `git archive 27ce163e` export, each sub-milestone's acceptance run, and
>   the result compared byte-for-byte with the mutant tree (§1.2 C30–C31). The executor applies each
>   patch and runs its acceptance; it writes no code by hand and runs no browser.

- **Design**: [w-world-approvals-by-key-a2.md](w-world-approvals-by-key-a2.md), **r2** (quorum r2 +
  controller carve-out; **not reopened**). Where this plan and the design differ, this plan was
  measured against `27ce163e` and wins; each difference is a decision (§1.3) or a premise found false
  (§1.4).
- **Starting ground**: sprint 1 (M0–M3a) PR #546 (merge `d5d90e49`); sprint 2 (M3b, M4, M5) PR #548
  (merge `3e55451f`); sprint 3 (M6, the CLI) PR #550 (merge `0ce117c5`), recorded in `f89337af`. Sprint
  3's Appendix C is this sprint's scope; judge 381's early `--name` note is folded in (M7a).
- **Planner base**: `27ce163eb1ae8af55151ea6548591c15a3126835` (= origin/dev). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter382-row201-a2s4`, branch
  `sprint/world-iter382-row201-a2s4`. Planner scratch (outside the repo, never under `/tmp`):
  `$HOME/.ailang/state/world-iter382/plan-scratch/` (`base/` the pristine export, `work/` the tree the
  mutants ran on, `stagebuild/` one commit per sub-milestone, `patches/`, `tools/`, `logs/`,
  `roundtrip/`, `gen2/` a `generate.sh` recording).
- **Executor**: as routed by the controller (a small model). **Four runs**: Run 1 (M7a, M7b), Run 2
  (M7c, M7d), Run 3 (M7e), Run 4 (M7f, M7g). Each sub-milestone is one `patch -p1` (§2) plus its
  acceptance (§3). **Apply, never retype. The executor runs no Chrome and no node**: XB1, XB2 and
  `generate.sh` run in the controller's §6 gates.
- **Touched**: `cmd/ailang-worldd/{approver_passkey.go, approver_passkey_test.go,
  approve_passkey_test.go}` (M7a); new `host/surfacepage/{passkey.go, passkey/index.html,
  passkey/passkey.js}` (M7b); new `cmd/ailang-worldd/testdata/passkey_cdp.mjs` (M7c); new
  `cmd/ailang-worldd/passkey_page_driver_test.go` (M7d); `host/daemon/{daemon.go, surface_page.go,
  surface_app_test.go}`, new `host/daemon/passkey_page_test.go`, new
  `cmd/ailang-worldd/{passkey_page_chrome_test.go, passkey_page_abort_test.go}` (M7e); new
  `cmd/ailang-worldd/passkey_fixture_record_test.go`, `host/approval/testdata/webauthn/{generate.sh,
  SOURCE.txt}` (M7f); `docs/QUICKSTART.md` (M7g).
- **Not touched**: `surfaceCSP` and `setSurfaceHeaders` (read only), every other `host/daemon` and
  `host/surfacepage` file, every `host/approval` Go file (the package stays pure and
  standard-library-only), the committed fixtures `chrome-es256.json`/`chrome-eddsa.json` (byte-for-byte:
  XW0 stays green unedited), `host/store/**`, `host/authority/**`, every other `cmd/ailang-worldd` file,
  `.gitattributes`, `world/*.ail`, `packages/**`, `scripts/**`, `tools/launchd/*`, `cmd/world-publish/*`,
  any skill, any `/v1` route, `~/.ssh`, `~/.ailang/world`, `~/.ailang/state/mission-v1*`.

---

## 0. Binding constraints (each is restated where it governs)

1. **Test keys, temp stores and a CDP virtual authenticator only (D-WORLD-109 = A).** Every browser run
   (XB1, XB2, `generate.sh`, the planner's measurements) is Chrome `--headless=new` with a temporary
   profile under `t.TempDir()` (or the planner's scratch), `--remote-debugging-port=0`, and
   `WebAuthn.addVirtualAuthenticator` (`ctap2_1`, `internal`, resident key, UV): never the operator's
   profile, keychain, iCloud passkeys or a platform authenticator. Every store is under
   `t.TempDir()`; every terminal is a test double (`ceremonyTerm`). Nothing opens `~/.ailang/world`,
   `~/.ssh`, a real credentials store or `/dev/tty`; nothing mints into a real store, publishes or runs
   `tailscale serve`; no `ssh-keygen`. QUICKSTART §16 is written and tagged *attended — pending first
   verbatim run*; nothing in this sprint runs it.
2. **Pure core, effects at the boundary (coding standards S1–S3).** No `host/approval` Go file
   changes. The page holds no data and makes no request (`fetch`, `XMLHttpRequest`, form `action`:
   none; text with `textContent` only); the daemon routes read no store and no request (PG1 serves them
   with the world store closed). `surfaceCSP` is not changed. No new Go module, package, production
   exec site, `/v1` route, context root or `.ail`.
3. **One mutant per new production branch, anchored to the diff** (§4): every condition, operand,
   bound, constant and copied field of Appendix P's production hunks (the page's script and HTML, the
   embed, the two routes and handlers, WN2a, `generate.sh`'s refusals), each run against its single
   named test **and** its package suite (browser mutants with `WORLD_REQUIRE_PASSKEY_BROWSER=1`), named
   equivalents with their reasons. Every run appends a line to `plan-scratch/logs/mutant_ledger.tsv`.
4. **Gated browser tests fail loudly when required.** XB1 and XB2 skip with `SKIP COUNTED: …` when
   `WORLD_CHROME` is unset or `node` is not on `PATH`, and **fail** instead when
   `WORLD_REQUIRE_PASSKEY_BROWSER=1`. §6 runs them required.
5. **No `t.Parallel()`** in any new test (`TestNoParallelWriteForkPackagesOutsideVerifygate`).
6. **Bisectable and safe at every commit.** Each sub-milestone ends gofmt-clean, `go build ./...` and
   `go vet` green, the `_test.go` compile fence (`go test -run '^$' ./...`) green, its own tests and
   every earlier one green. **No route serves the page before its tests land** (M7b embeds it unserved;
   M7e adds the routes together with PG1, XB1 and XB2), and no link points at a page that is not served
   (the CLI's links have pointed at `/surface/passkey` since sprint 3; it is a 404 until M7e).
7. **No git writes by the executor** (§7). The controller commits from snapshots.

---

## 1. Pre-flight

### 1.1 Baseline at `27ce163e` (measured by the planner, 2026-11-11)

Every command ran with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` (`AILANG v0.41.0`), `go1.26.6
darwin/arm64`, in a `git archive 27ce163e` export at `$HOME/.ailang/state/world-iter382/plan-scratch/base`
(a throwaway git repo; `diff -r cmd host` against the worktree: identical, `BASE_EQUALS_WORKTREE`; the
worktree itself was only read). The pre-commit legs are `tools/gates.sh <tree> <logs> base` = §6
without the four new tests (they do not exist at the base). Logs: `plan-scratch/logs/baseline/`.

| Command | rc | Time | Result |
|---|---|---|---|
| `git -C <worktree> status --porcelain \| wc -l` | 0 | — | `0` (clean at `27ce163e`) |
| `./scripts/verify_ail.sh` | **0** | — | `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` |
| `gofmt -l cmd host`; `go build ./...`; `go vet ./...`; `go test -run '^$' ./...` | **0** | — | empty / clean / no FAIL line |
| `go test -count=1 ./host/surfacepage/... ./host/approval/... ./cmd/ailang-worldd/ ./host/daemon/` | **0** | — | all `ok` (`cmd/ailang-worldd` 159.8 s, `host/daemon` 325.6 s) |
| `-race` legs (§6, the base's tests of the same names) | **0** | — | `ok`, no race report |
| §6's 26 named inventory, reflective, help and quickstart tests, `-v` | **0** | — | 26 `--- PASS`, 0 SKIP, 0 FAIL; `TestSubprocessSinkGate` (added to §6 later, C34) run alone: `ok` (0.6 s) |
| the whole pre-commit script | **0** | **475 s** | `PRE_COMMIT_GATES_GREEN 475s` (load average 18–47 during the run) |
| `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` | **1** | ~12 min | **red at the base, measuring the machine**: run while the planner's six mutant workers and other sessions held the load at 40–66, it timed out or missed subprocess deadlines in `cmd/ailang-worldd` (`TestCLIRealSubprocessEpisode`: daemon announcement timed out; a 10 m package timeout), `host/broker` (`TestExecProbeNet6*`, `TestExecProtectedCheckTimeout`), `host/coordinator` (`TestCommittedOutputCarriesTheWorktreeDelta`, a 5 s exec), `host/daemon`, `host/evidence` and `host/verifygate`; every static leg green (`✓ 0 other binary blobs among 6716 tracked files`, toolchain floor, `✓ all 64 required top-level evidence tests passed exactly once`) — `logs/baseline/verify_go.log`. Run the after-commit leg on a quiet machine and re-run a red untouched package alone once (§6) |

**A trap measured at the base (C8):** the first base run exported `WORLD_CHROME` for the whole
script and went red in `cmd/ailang-worldd` on `TestMarketingChromeSixShots` ("AILANG_BIN and
WORLD_TOOL_AILANG_BIN are required"): `WORLD_CHROME` arms that test too. §6 therefore sets
`WORLD_CHROME` **only on the XB leg**, never exported (`logs/baseline_try1_worldchrome_exported/`).

**Rule (carried):** a red at the base measures the repo, not the sprint. The controller re-runs a red
untouched package alone once and does not chase it. **Never start a gate under `nohup`.**

### 1.2 Codebase claims checked (at `27ce163e` unless "work" — the planner's final scratch tree)

| # | Claim | Command → observed output |
|---|---|---|
| C1 | The base | `git rev-parse HEAD origin/dev` (worktree) → `27ce163eb1ae8af55151ea6548591c15a3126835` twice; `git status --porcelain \| wc -l` → `0` |
| C2 | **The page is not served yet** | `grep -n 'surface/passkey' host/daemon/daemon.go` → **empty**; control in the same call: `grep -n '"GET /surface/app.js"' host/daemon/daemon.go` → `1111:` |
| C3 | The surface headers and CSP the page must use | `surface_page.go:19 const surfaceCSP = "default-src 'none'; style-src 'unsafe-inline'; script-src 'self'; connect-src 'self'; font-src 'self'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'"`; `:21 func setSurfaceHeaders` = `setWorkbenchHeaders` (`workbench.go:57`: `text/html; charset=utf-8`, `no-store`, `nosniff`, `no-referrer`) + `surfaceCSP` |
| C4 | The app's assets use another header set | `surface_page.go:155`, `:177` (`writeSurfaceAsset`): `Content-Security-Policy: default-src 'none'`, `Cache-Control: no-cache`. The design says the passkey routes use `setSurfaceHeaders`; S4-D3 keeps that for both (the script then overrides only `Content-Type`) |
| C5 | **`ceremonyLink` lives in `package main`** (`cmd/ailang-worldd/approver_passkey.go:1 package main`, `:36 type ceremonyLink struct`) | a `host/daemon` test cannot import it, so XB1 there could not build its links with the production code (§1.4 P1) |
| C6 | A `cmd/ailang-worldd` test can serve the production daemon handler | `demo_proof_test.go:89 d, err := daemon.New(boundedTestContext(t), daemon.Config{…})`; on work `passkeyPageOrigin` serves `d.Handler()` with `httptest.NewServer` |
| C7 | `generate.sh` is the sprint-1 stub and is not pinned | `cat host/approval/testdata/webauthn/generate.sh` → prints "regeneration lands in row 201 A2 sprint 3 (M7)" and `exit 1`; mode `-rw-r--r--` (run as `sh generate.sh`); `grep -n webauthn host/approval/approver_test.go` pins only `testdata/approvers` (its `generate.sh` sha `54d00a3e…`), and XW0 (`webauthn_fixture_test.go:88`) pins only the two JSON files |
| C8 | `WORLD_CHROME` arms another `cmd/ailang-worldd` test | `marketing_chrome_test.go:52 chrome := os.Getenv("WORLD_CHROME")`; base run 1 with it exported → `--- FAIL: TestMarketingChromeSixShots … AILANG_BIN and WORLD_TOOL_AILANG_BIN are required` (§1.1) |
| C9 | The AC10 exec-site census skips `_test.go` | `registry_publish_test.go:1094 if d.IsDir() \|\| !strings.HasSuffix(path, ".go") \|\| strings.HasSuffix(path, "_test.go") {` — the Node launch in `passkey_page_driver_test.go` adds no production exec site; on work `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` `--- PASS` (C25) |
| C10 | The name grammar the early check uses | `host/approval/approver.go:25 var ApproverNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,47}$`)`, `:39 ErrApproverName` |
| C11 | Instruments | `Google Chrome --version` → `Google Chrome 157.0.8103.50`; `node --version` → `v26.13.0` (`/opt/homebrew/bin/node`); `go version` → `go1.26.6 darwin/arm64`; `$HOME/.pinned-ailang/ailang --version` → `AILANG v0.41.0` |
| C12 | **Chrome writes `DevToolsActivePort` in two steps** (work) | the first `generate.sh` run: `DOMException [SyntaxError]: TypeError: Invalid URL … new WebSocket` (the file existed but was empty); the driver now waits for `^\d+\n/devtools/browser/\S+` (`logs/generate1.log` before the fix; every later run green) |
| C13 | **A navigation that changes only the fragment fires no load event** (work) | the first XB1 run: `driver open: timed out waiting for the page load` in `links-for-another-page` (same path, new `#…`); the driver now loads `about:blank` before each link |
| C14 | **A Save link below the fold is not clicked** (work) | the first XB1 run: `driver ceremony: timed out waiting for the download`; the driver now `scrollIntoView({block:'center'})`s before reading the box |
| C15 | **XB2's exact deadline, measured — and a premise of the first draft found false under load** (work) | first draft: virtual time started only after the click, so the injected operation's clock started a moment after the page's `setTimeout`; quiet: `performance.now()` deltas `239999.6`–`239999.9`, `Date.now()` `240000`/`239999`; **at load 36–80: `239998`** (`logs/xb/run4.log`, a red run). Fixed in the driver, not the assertion: `Emulation.setVirtualTimePolicy({policy: "pause"})` **before** the click, so the page's timer and the operation start at one virtual instant. Then `Date.now()` delta = **exactly `240000`** in every run, `performance.now()` `239999.9`/`240000.1` (one 0.1 ms coarsening quantum; `logs/xb2_try*.log`). XB2 requires `Date.now()` = 240000 and `performance.now()` within 0.15 ms of it; MUT-JS-TIMEOUT moves it by 1000 ms. The late result arrives at `Date.now()` exactly `245000` |
| C16 | **`generate.sh` records fixtures the production law accepts** (work) | `sh generate.sh $S/gen2` → rc 0 (`logs/generate2.log`, ~22 s); copied over the committed fixtures in a throwaway tree, `go test ./host/approval/` → every law test green (XW1 layout, XW2: both recorded assertions approve the fixture request, `assertion` after `assertion2` refused) and **only the re-pins red**: `TestPasskeyFixturesAreTheGeneratedBytes` (the two SHA-256s), `TestEnrolmentCodeGolden`, `TestRenderPasskeyEnrolmentGolden`, and `TestPasskeyEnvelope`'s `!HasUserHandle` (the page asks for a resident key, so its assertions carry `userHandle`; iteration 379's probe did not) — `logs/genverify2.log` |
| C17 | The recorder enrols as `mark-mac` (work) | a first recording as `world-approver` made assertions whose `userHandle` the fixture tests (which enrol `mark-mac`) refuse: `approval: the assertion's userHandle is not the approver's name` (`logs/genverify.log`); recorded as `mark-mac`, C16 |
| C18 | **Chrome sometimes adds `other_keys_can_be_added_here` to `clientDataJSON`** (work) | the first recording's ES256 registration: `TestTestPasskeyMatchesChromeLayout … Chrome's type,challenge,origin,crossOrigin,other_keys_can_be_added_here` (`logs/genverify.log`); the second recording had none (`grep -c other_keys gen2/*.json` → 0, 0). Chrome injects it at random; the law accepts it (XB1 green either way). Whoever re-records with `generate.sh` re-runs it until XW1 is green or re-pins XW1 (Appendix C) |
| C19 | **XB1 and XB2, consecutive required runs** (the round-trip tree, `WORLD_REQUIRE_PASSKEY_BROWSER=1`) | `WORLD_REQUIRE_PASSKEY_BROWSER=1`, `-count=1`, both tests per run. **Final code, §6 script on work** (`logs/final_gates/xb_{1,2,3}.log`, `xb_times.log`): run 1 XB2 4.36 s / XB1 3.80 s (9 s wall), run 2 3.41 / 3.18 (7 s), run 3 3.25 / 3.13 (8 s), load 7–12 — 3 of 3 green, 7 subtests `--- PASS` each. **Code before C34** (differs only in the driver's stderr buffer type): round-trip tree 4.22/4.40, 5.53/7.64, 4.00/4.84, 4.40/6.10 s at load 34–38 (`logs/xbfinal/summary.txt`) and the first §6 run 3.48/3.10, 3.30/3.15, 3.95/4.71 s (`logs/final_gates_try1_pre_syncbuffer/`) — 7 of 7 green. **Code before C23/C24** (same page): 6 of 6 green at load 38–80, XB2 4.97–18.46 s, XB1 6.43–12.52 s (`logs/xb_c{1..6}.log`). 250 s of page time costs 26–98 ms. Every red run (C12, C13, C14, C15 `logs/xb/run4.log`, C21, C23) was fixed in the driver or test before the final code. Slowest test measured: 18.5 s at load 80, against `passkeyBrowserBudget` 120 s |
| C21 | **A stray `downloads.html` download under load** (work) | 3 of ~16 XB runs at load ≥ 60 found `downloads.html`, `downloads.html.crdownload` or `downloads.html (1).crdownload` in an XB2 subtest's download directory, with **no** `Browser.download*` event inside that ceremony's window and no click on Save (Save was hidden: `click` returned false). Not the page's save (that is `world-*.json`, and every page save in XB1 has its events). Unexplained; routed by counting only `world-*.json` files as saves and logging anything else (`worldFiles`), while still requiring zero download events in the window (`logs/xb2_b5.log`, `xb_c1.log`) |
| C22 | **The suite of a page-asset mutant** | the page's bytes reach a test only through the daemon handler: `grep -rn --include='*.go' 'surfacepage.PasskeyPage\|surfacepage.PasskeyScript' cmd host` → only `host/daemon/surface_page.go:177`, `:184`; the tests that serve or open the page are PG1 (`host/daemon`) and, in `cmd/ailang-worldd`, only `passkey_page_*_test.go` and `passkey_fixture_record_test.go` (`grep -rln --include='*_test.go' 'passkeyPageOrigin\|startPasskeyBrowser'`). So a JS/HTML mutant's suite is `cmd/ailang-worldd`'s **passkey suite**, `-run '^(TestPasskey\|TestApproveByPasskey\|TestApprovePasskey)'` (CW1–CW9, XB1, XB2, the recorder, the generator), with `WORLD_REQUIRE_PASSKEY_BROWSER=1`; a Go or `generate.sh` mutant's suite is its whole package. Measured why: a whole `cmd/ailang-worldd` run took **~15 min** per mutant at load 40–100 (6 in parallel, run2), and at that load it is red at the base on timing tests (`TestHostsLaunchersExecHost`, `TestWorldContextAndStatusline`, …; §1.1, C27), so a whole-package "suite" verdict for a page mutant measured the machine, not the mutant (run1's three equivalents "KILLED" by the suite, `logs/mutants_run1.err`) |
| C23 | **XB1 cannot observe "pending"** (work, run4) | a virtual authenticator answers within the driver's next round trip: under load the page was already `done` when read after the click (`while pending: state "done"`, run4 MUT-JS-ASR-USERHANDLE/MUT-JS-RP-DEFAULT, `logs/mutants_run4.err`). XB1 no longer asserts it; XB2 does, with an operation that stays pending |
| C24 | **`WebAuthn.setUserVerified(false)` does not tell `required` from `preferred`** (work, run4) | MUT-JS-UV-CREATE and MUT-JS-UV-GET SURVIVED XB1 when the UV-false variant was an authenticator *with* UV capability that fails it: Chrome refuses both. Routed: the UV-false variant has **no** UV capability (`hasUserVerification: false`), and the get arm moves the enrolled credential onto such an authenticator (`WebAuthn.getCredentials` → `removeVirtualAuthenticator` → `addVirtualAuthenticator` → `addCredential`), then back for the retry |
| C32 | **Another session trimmed the shared Go build cache mid-run** (run3) | 57 of 69 mutants `NOCOMPILE`: `could not import sync (open …/Library/Caches/go-build/0a/0a3891e5…-d: no such file or directory)` (`logs/mutants_run3.err`); `go vet` of the same packages was clean a minute later. The mutant runs from run4 on use a private `GOCACHE` under `plan-scratch/gocache` |
| C33 | **The whole `host/daemon` package exceeds `go test`'s default 10 min under load** | the round trip's M7e acceptance at load 30–60: `panic: test timed out after 10m0s` after `TestRegistryPhaseHasItsOwnClock` ("`node --version`: broker: handler subprocess timed out after 5s", a known load-sensitive test) — `logs/rt_try3_m7e_timeout/M7e.daemon.log`; the same package is 313–326 s at load ≤ 20 (§1.1, C26). §3 and §6 pass `-timeout 20m` to whole-package runs: a bound on the run, not an assertion |
| C34 | **A repo-wide inventory the first gate list missed: `TestSubprocessSinkGate`** (`host/verifygate`) | the first after-commit `verify_go.sh` on the round-trip tree: `cmd/ailang-worldd/passkey_page_driver_test.go:78: subprocess sink "stderr" is a bare bytes.Buffer …; use the syncBuffer pattern` (`logs/final_verify_go_try1.log`). Fixed (the driver's stderr is `cli_test.go`'s existing `syncBuffer`) and **added to §6's inventory list**; `--- PASS` at the base and on the final tree (`violations=0 filesParsed=721 startSites=27`). Every cmd/ailang-worldd mutant was re-run after the fix (run6) |
| C20 | XB1 and XB2 skip, counted, without Chrome (work) | `env -u WORLD_CHROME go test ./cmd/ailang-worldd -run '^(TestPasskeyPageInChrome\|TestPasskeyPageAbortsAPendingCeremony)$' -v` → two `--- SKIP` with `SKIP COUNTED: the passkey page browser test needs WORLD_CHROME ("") and node on PATH (<nil>)` (M7e's acceptance asserts both) |
| C25 | **Every repo-wide inventory, on the final tree** | the §6 script (`tools/sprint_gates.sh` = §6 verbatim) on `plan-scratch/work` (byte-identical to the round-trip tree) → `PRE_COMMIT_GATES_GREEN` in ~7 min (load 7–12): **29 named tests `--- PASS`, 0 SKIP/FAIL** (`logs/final_gates/inv.log`): design §7's inventories, the reflective, help and quickstart tests, PG1, the generator's refusals and `TestSubprocessSinkGate` (C34). `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` (the Node launch is in a `_test.go`), `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots` (no new root), `TestFrozenV1RouteTableMatchesMux` (13), `TestNoParallelWriteForkPackagesOutsideVerifygate`, `TestApprovalImportsOnlyTheStandardLibrary` unchanged; XW0 green (fixtures unchanged) |
| C26 | Whole packages on the final tree | same run: `go test -count=1 -timeout 20m ./host/surfacepage/... ./host/approval/... ./cmd/ailang-worldd/ ./host/daemon/` → all `ok` (`cmd/ailang-worldd` 156.5 s, `host/daemon` 301.9 s; `logs/final_gates/pkgs.log`); every earlier test green with only the three named edits. The run before it, at load 73, was red only on `TestWorldContextAndStatusline` (`world-context: signal: killed`, a 5 s subprocess; `logs/final_gates_try2_load73/`), green when re-run |
| C27 | **The whole repository on the final tree** (committed) | `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` on the round-trip tree at M7g (committed, load 7–15) → **rc 0**, `✓ go gate PASSED (plain)`: `✓ 0 other binary blobs among 6725 tracked files`, toolchain floor, `✓ all 64 required top-level evidence tests passed exactly once`, `go test ./...` green (`logs/final_verify_go.log`, ~8 min). The run before it (pre-C34) was red on `TestSubprocessSinkGate` — **this sprint's defect, fixed** (C34) — and on the known `TestWorktreeSnapshotRefusesAConfiguredFilter` (`logs/final_verify_go_try1.log`) |
| C28 | `go vet ./...`, gofmt and `-race` on the final tree | same run: `go vet ./...` empty; `gofmt -l cmd host` empty; `-race` `host/surfacepage` 3.1 s, `host/daemon` (PG1, `TestSurfaceAppRoutes`, `TestSurfaceScriptRoute`) 1.8 s, the `cmd/ailang-worldd` passkey/approver subset 8.7 s → `ok`, no race report (`logs/final_gates/race.log`); `./scripts/verify_ail.sh` → `✓ verify gate PASSED` |
| C29 | Every mutant, on the final tree | 69 rows, the last ledger line of each: **64 KILLED (named and suite), 5 named equivalents SURVIVED (named and suite)**, 69 of 69 as expected; 251 ledger lines in all (§4) |
| C30 | **The patches reproduce the mutant tree** | a fresh `git archive 27ce163e` export + `patch -p1` of Appendix P's seven patches **extracted from this file with §2's `awk`** (each `cmp`-equal to the generated patch; `--dry-run` first; no `.orig`/`.rej`) → `diff -r cmd host` and `cmp docs/QUICKSTART.md` against the tree §4's mutants ran on (`plan-scratch/work`): **identical** (`ROUNDTRIP_IDENTICAL_TO_MUTANT_TREE`, `tools/roundtrip.sh`, `logs/roundtrip.log`); nothing uncommitted after M7g; `generate.sh` stays `-rw-r--r--` |
| C31 | **Every sub-milestone is bisectable** | the same fresh export, after **each** patch: `tools/accept.sh <M>` = §3 verbatim → `ACCEPT <M> rc=0` for all seven (`logs/rt/`); the cumulative numstat after a commit of each = §7: **11, 195, 367, 597, 924, 1054, 1054** (`NUMSTAT` lines of `logs/roundtrip.log`). Two load reds on the way, both untouched timing tests, each green re-run alone: M7e's whole `host/daemon` hit the default 10 min timeout at load 30–60 (`logs/rt_try3_m7e_timeout/`; hence §3's `-timeout 20m`, C33), and M7g's whole `cmd/ailang-worldd` failed `TestWorldContextAndStatusline` once (its `1h0m`/`1h1m` minute-boundary cache line; green 3/3 alone on the round-trip tree and 2/2 at the base; `logs/rt/M7g.accept.try1_*.out`), then M7g's acceptance re-ran green |

### 1.3 Planner decisions where the design leaves a choice

- **S4-D1 — XB1 and XB2 live in `cmd/ailang-worldd`, not `host/daemon`** (§1.4 P1). The links come
  from the CLI itself (XB1) or `ceremonyLink.link()` (XB2, the bad-link arms), and the files go back
  through `runApprover`/`runApprove`; the page is still served by the production daemon (`daemon.New`
  → `Handler()` → `httptest.NewServer`, origin `http://localhost:<port>`). The driver is therefore
  `cmd/ailang-worldd/testdata/passkey_cdp.mjs`.
- **S4-D2 — One long-lived driver per test, JSON lines over stdin/stdout.** The virtual authenticator
  belongs to one page target, and the approve step needs the credential the enrol step created, so one
  Chrome serves a whole test: commands `tab` (a fresh target, optionally a virtual authenticator
  `uv`/`no-uv`, optionally an injected operation), `open`, `ceremony` (a trusted mouse click on the
  button; virtual time or real; the Save click and the download events), `uv` (moves the credentials,
  keys and counters, to a new virtual authenticator with or without user verification: C24) and
  `advance` (virtual time). `no-uv` is an authenticator with **no** user-verification capability. EOF on stdin kills Chrome and removes the profile; the Go cleanup
  kills Chrome's pid only if the driver died. `passkeyBrowserBudget` = 120 s per test, against 4–19 s
  measured at load 20–80 (C19; rule 3m: ≥ ×6 headroom). No `t.Fatal` from the terminal's read goroutine:
  the terminal double uses `b.call` and `t.Error`.
- **S4-D3 — Both routes use `setSurfaceHeaders`** (design §2.14 "with `setSurfaceHeaders` and
  `surfaceCSP`"), the script overriding only `Content-Type` to `text/javascript; charset=utf-8`; so both
  are `no-store`, `nosniff`, `no-referrer`, under `surfaceCSP`. Not `writeSurfaceAsset` (C4).
- **S4-D4 — The page's states are `data-state` on `<body>`**: `loading` → `ready` | `refused`; on the
  click `pending` → `done` | `failed` | `expired`. The driver and XB1/XB2 read it (no test hook in the
  page). The button is disabled while pending, re-enabled after `failed` or `expired` (retry on the same
  page), and stays disabled after `done`.
- **S4-D5 — The [r2] timer, structurally.** `setTimeout(240000)` starts immediately before the
  `create`/`get` call that carries `signal: controller.signal`; on expiry it sets `expired`, aborts,
  sets `expired` state and the timeout text. `clearTimeout` runs in `finally` around the pending
  operation only, so nothing after the result can expire. `if (expired) { re-enable; return }` is the
  one branch that both ignores a late result and keeps saving disabled after expiry: the Save link is
  hidden from load and shown only at `done`, which an expired attempt cannot reach (MUT-JS-LATE,
  MUT-HTML-SAVE-HIDDEN, XB2's late arms).
- **S4-D6 — XB2 has four arms**: `create`/`get` × an injected operation that rejects with
  `signal.reason` on abort (`abort`) or ignores the signal and resolves with a fake credential at 245000
  ms (`late`). Virtual time is paused from before the click, then advanced 250 s at once. Each asserts:
  pending + button disabled after the click; the operation got a signal, saw it abort at exactly 240000
  ms of `Date.now()` (C15); `expired`, the exact timeout text with the error flag; no
  textarea, no code, Save hidden with no `href`, no download event, no file; the button re-enabled; no
  CSP violation or page error; and for `late`, the fake result really arrived (at 245000) and was ignored.
- **S4-D7 — XB1's three subtests**: `enrol-then-approve` (above, plus a `get` on an authenticator
  without user verification first — `failed`, nothing saved, the retry on the same page succeeds — and, after the approval's
  Save, 250 s of virtual time with the page still `done`: the finished attempt's timer was cleared,
  MUT-JS-CLEAR); `user-verification-required` (a UV-false authenticator: `create` fails, nothing saved,
  the CLI refuses at WN8, nothing enrolled — design's "`userVerification: preferred` with a UV-false
  authenticator" mutant); `links-for-another-page` (a control link → `ready`; version 2, op `sign`,
  another origin, another rp, not base64url, not JSON → `refused` with the button hidden).
- **S4-D8 — `generate.sh` records through the page.** It runs `TestPasskeyPageRecordsTheFixtures`
  (skips unless `WORLD_PASSKEY_RECORD_DIR` is set) into an **empty directory it is given**, never over
  the committed fixtures: ES256 by the page as-is, EdDSA by an injected filter that leaves only `-8` in
  `pubKeyCredParams`; create + two gets over the fixture message hash `29181d22…fef5`, as `mark-mac`,
  in the committed layout. Copying its output over the fixtures is a re-pin (C16, Appendix C). It
  refuses (exit 2, nothing written, nothing run) a call without exactly one existing empty directory or
  without `WORLD_CHROME`; `TestPasskeyFixtureGeneratorRefuses` runs those refusals with no browser and
  `PATH=/usr/bin:/bin` (no `go`, no `node`), so a refusal that falls through cannot recurse.
  **Mode**: `generate.sh` stays `0644` and is run as `sh generate.sh` (as its usage line and the
  refusal test do): `patch` need not carry a mode, and nothing depends on the executable bit.
- **S4-D9 — The page builds the §2.7 JSON itself** (`id`, `rawId`, `type`, `response` with the two or
  four byte fields and `userHandle` or `null`), never `toJSON()`; `JSON.stringify` keeps §2.7's member
  order; the Save blob is that JSON plus `\n`, named `world-registration-<name>.json` /
  `world-assertion-<request>.json` (sprint 3's names).
- **S4-D10 — The page refuses a link for another page** before showing any button: a fragment that is
  not base64url or not JSON, `v` ≠ 1, `op` not `create`/`get`, `origin` ≠ `location.origin`, `rp` ≠
  `location.hostname`. Everything else in a link is the CLI's (WebAuthn itself refuses a bad challenge
  or credential ID, and the CLI judges the result).
- **S4-D11 — WN2a (carry-forward, judge 381; S3-D13 R6)**: right after WN2, a `--name` that fails
  `approval.ApproverNamePattern` refuses with `exitUsage`, "`--name "Mark": approval: approver name must
  match …`", before the stores, the terminal or the link (CW1 arm `WN2a-name-grammar`). The law's R1
  still judges the name (unchanged). **And** CW7 now pins the whole next-step line,
  `ailang-worldd approve <id> --webauthn-assertion <the saved file> --db <the request's world store>
  --credentials-db <its credentials store>` (from the request's own `WorldDB`/`CredentialsDB`).
- **S4-D12 — Sub-milestones** (§1.5): M7a carry-forwards; M7b the page and its embed, **unserved**;
  M7c the driver; M7d the Go harness; **M7e the routes with PG1, XB1 and XB2 together** (no commit
  serves the page without its tests); M7f the generator; M7g QUICKSTART §16. M7c and M7d are inert test
  support (nothing calls them until M7e); M7e (327) is above ~250 because the routes may not land
  without both browser tests.

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

- **P1 — Design §6/§8: XB1/XB2 in `host/daemon` (`passkey_page_*_test.go`, driver at
  `host/daemon/testdata/passkey_cdp.mjs`); the directive and Appendix C: "XB1 should build its links
  with the production code".** Both cannot hold: `ceremonyLink` is in `package main` (C5). Routed: the
  browser tests and the driver live in `cmd/ailang-worldd` (S4-D1). **Files outside design §8's
  sprint-3/4 conflict surface**: `cmd/ailang-worldd/{passkey_page_driver_test.go,
  passkey_page_chrome_test.go, passkey_page_abort_test.go, passkey_fixture_record_test.go,
  testdata/passkey_cdp.mjs}` (instead of the `host/daemon` names), `host/daemon/passkey_page_test.go`
  and `surface_app_test.go` (design §6 names PG1 and the edit, §8 does not list them),
  `cmd/ailang-worldd/{approver_passkey_test.go, approve_passkey_test.go}` (M7a's two assertions), and
  `host/approval/testdata/webauthn/SOURCE.txt` (its last line named "sprint 3" for the generator).
- **P2 — Design §2.15 / the stub: "generate.sh … runs the M7 Node/DevTools driver against Chrome's
  virtual authenticator" and "re-records the §2.15 fixtures".** Re-recording cannot reproduce the
  committed bytes (fresh keys, random ES256 signatures; C16) and the page asks for a resident key, so
  its assertions carry a `userHandle` iteration 379's did not. Routed: `generate.sh` writes into an
  empty directory and never over the fixtures (S4-D8); the committed fixtures and XW0 are unchanged;
  what a re-pin involves is measured (C16) and written in `generate.sh` and Appendix C.
- **P3 — Design §2.14 lists `alg -7, -8`; the page offers both, but Chrome's virtual authenticator
  always takes ES256**, so EdDSA through the page is not observable by XB1 (MUT-JS-ALG-NO-EDDSA, a
  named equivalent). XB1 asserts the enrolled algorithm is ES256 (the page's first preference); the
  EdDSA path through a real page is `generate.sh`'s (filtered) recording, and the CLI's EdDSA
  handling is CW2's (sprint 3).
- **P4 — Sprint 3 Appendix C: "the probes are in `/tmp` and may be gone".** They were preserved
  outside `/tmp` for this sprint (`plan-scratch/w379-probes/`, `vt_probe.mjs`) and generalised into
  the driver; none is committed.
- **P5 — The planner directive's "XB2 … `t.Skip("SKIP COUNTED: …")` unless
  `WORLD_REQUIRE_PASSKEY_BROWSER=1`"** holds as written; note that every `go test ./cmd/ailang-worldd`
  without `WORLD_CHROME` (CI, the executor) now shows two counted skips and one recorder skip. CI does
  not run the browser tests (design §9: Chrome and Node on the runner are UNVERIFIED); §6 does.

### 1.5 The size

Lines are insertions + deletions under `cmd/ host/`: **`git diff --numstat <base> HEAD -- cmd host`**
(`docs/QUICKSTART.md` and `design_docs/` are outside the filter). **Every row is measured** from the
scratch stages, each committed in a throwaway repo (`plan-scratch/stagebuild`, tags `M7a` … `M7g`;
`tools/stages.sh`, `tools/build_patches.sh`).

| Sub-milestone | Production | Test / driver / generator | Step | Cumulative vs `27ce163e` (after its commit) | Run |
|---|---|---|---|---|---|
| M7a — WN2a; CW1 arm; CW7's next-step line | 3 | 8 | **11** | 11 | 1 |
| M7b — `index.html`, `passkey.js`, `passkey.go` (embedded, not served) | 184 | 0 | **184** | 195 | 1 |
| M7c — `testdata/passkey_cdp.mjs` (inert) | 0 | 172 | **172** | 367 | 2 |
| M7d — `passkey_page_driver_test.go` (inert) | 0 | 230 | **230** | 597 | 2 |
| M7e — 2 routes + 2 handlers; PG1; `TestSurfaceAppRoutes`; XB1; XB2 | 19 | 308 | **327** | 924 | 3 |
| M7f — `generate.sh`, `SOURCE.txt`; the recorder; the refusals | 0 | 130 | **130** | 1,054 | 4 |
| M7g — QUICKSTART §16 (`docs/`, outside the filter) | — | — | **0** | 1,054 | 4 |
| **Sprint 4** | **206** | **848** | **1,054 ≤ 2,000** (margin 946) | | |

The design estimated M7 at ≈ 530 under the filter (230 page + 300 test/driver); it measured **1,054** (2.0×): the page as estimated (184 + 19), the tests and driver 2.8× (848 vs 300) — the driver alone is 172, the harness 230, XB1 171, and `generate.sh` with its recorder and refusal test 130, which the estimate counted as "the generator". Nothing was narrowed: the whole of Appendix C fits with 946 to spare.

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of
four backticks, in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; a new file's patch has an
empty "a" side. Apply it **only** this way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF'
set -euo pipefail
M=M7a   # <- the sub-milestone, exactly as in the Appendix P heading: M7a M7b M7c M7d M7e M7f M7g
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter382-row201-a2s4
PLAN="$W/design_docs/planned/w-world-approvals-by-key-a2-sprint4-plan.md"
T="$HOME/.ailang/state/world-iter382/exec-patches"; mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
patch -p1 --dry-run < "$T/$M.patch"
patch -p1 < "$T/$M.patch"
if find cmd host docs \( -name '*.orig' -o -name '*.rej' \) | grep .; then echo "STOP: patch left .orig/.rej"; exit 1; fi
git status --short
EOF
```

- If `--dry-run` fails, **stop and report** (no hand edits, no fuzz options). It can only fail if an
  earlier sub-milestone was skipped or changed.
- Then run that sub-milestone's **Acceptance** (§3) in one bash shell with
  `export AILANG_BIN=$HOME/.pinned-ailang/ailang` and **`unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER
  WORLD_PASSKEY_RECORD_DIR`**, each command's output `tee`'d to a log under
  `$HOME/.ailang/state/world-iter382/exec-logs/` (never only through `tail`).
- Then snapshot (§7). The controller commits and compares the numstat with §7's expected value.
- `generate.sh`'s mode: the patch modifies the existing `0644` file in place; nothing needs or sets
  the executable bit (S4-D8).

---

## 3. Milestones

Every sub-milestone's acceptance starts with the same three lines (abbreviated **[base]** below):

```
test -z "$(gofmt -l cmd host)"
go build ./... && go vet ./host/surfacepage/ ./host/daemon/ ./host/approval/... ./cmd/ailang-worldd/
go test -run '^$' ./...
```

Each `-v` run below must show `--- PASS` for every test it names and **no `--- FAIL`**, and **no
`--- SKIP`** except where the line says which skips it expects. The planner's `tools/accept.sh <M>` is
exactly these lines (§1.2 C31).

### M7a — carry-forwards: WN2a and the next-step line · 3 prod + 8 test · Run 1

**Patch** `M7a`:
- `cmd/ailang-worldd/approver_passkey.go`: **WN2a** after WN2 (S4-D11).
- `cmd/ailang-worldd/approver_passkey_test.go`: CW1 arm `WN2a-name-grammar` (`--name Mark` → the
  refusal, `exitUsage`, no terminal, no link, nothing written).
- `cmd/ailang-worldd/approve_passkey_test.go`: CW7 pins the whole next-step line from the request's
  own stores.

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestPasskeyEnrollFlagRules|TestApproveByPasskey)$' -count=1 -v      # no --- SKIP
go test ./cmd/ailang-worldd/ -count=1 -timeout 20m                                                      # ~3-5 min
```

**Expected cumulative numstat after the controller's commit of M7a: 11**.

### M7b — the page and its embed, not yet served · 184 prod · Run 1

**Patch** `M7b`: new `host/surfacepage/passkey/index.html` (the page: `data-state`, the button, the
status line, the code field, the result textarea, the hidden Save link, one `<script
src="/surface/passkey.js">`), new `host/surfacepage/passkey/passkey.js` (S4-D4–S4-D10: the fragment
checks, create/get with `signal`, the 240000 ms timer, the enrolment code — design §2.8, base32 of
the first 16 bytes of SHA-256(label ‖ 0x00 ‖ SHA-256(clientDataJSON) ‖ SHA-256(attestationObject)) in
five groups — the §2.7 JSON, the `blob:` Save link), new `host/surfacepage/passkey.go` (`PasskeyPage`,
`PasskeyScript`, `go:embed`). No route serves them yet: `/surface/passkey` is still a 404.

**Acceptance**
```
[base]
go test ./host/surfacepage/ -count=1
go test ./host/daemon/ -run '^(TestSurfaceAppRoutes|TestFrozenV1RouteTableMatchesMux)$' -count=1 -v   # no --- SKIP
```

**Expected cumulative numstat after the controller's commit of M7b: 195**.

### M7c — the browser driver · 172 test support · Run 2

**Patch** `M7c`: new `cmd/ailang-worldd/testdata/passkey_cdp.mjs` (S4-D2). Inert: nothing runs it
until M7e. **The executor does not run it** (no node, no Chrome).

**Acceptance**
```
[base]
test "$(head -c 2 cmd/ailang-worldd/testdata/passkey_cdp.mjs)" = "//"
grep -c 'WebAuthn.addVirtualAuthenticator' cmd/ailang-worldd/testdata/passkey_cdp.mjs                  # 1
```

**Expected cumulative numstat after the controller's commit of M7c: 367**.

### M7d — the Go harness · 230 test · Run 2

**Patch** `M7d`: new `cmd/ailang-worldd/passkey_page_driver_test.go`: `startPasskeyBrowser` (the
gate, S4-D2), `call`/`do`, `passkeyPageOrigin` (the production daemon handler on `localhost`),
`printedLink`, `saved` (a completed `blob:` download of the page's JSON plus `\n`, no CSP violation, no
page error), `worldFiles` (C21), `pageView`, `pageCeremony`, `opString`. Inert until M7e.

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApprove|TestApprover|TestPasskey|TestEveryVerbHasHelp)' -count=1
```

**Expected cumulative numstat after the controller's commit of M7d: 597**.

### M7e — the routes, PG1, XB1 and XB2 · 19 prod + 308 test · Run 3

**Patch** `M7e`:
- `host/daemon/daemon.go`: `GET /surface/passkey`, `GET /surface/passkey.js` (+ the routes comment).
- `host/daemon/surface_page.go`: `handleSurfacePasskey`, `handleSurfacePasskeyScript` (S4-D3).
- new `host/daemon/passkey_page_test.go`: **PG1** `TestPasskeyPageRoutes` (both routes 200 with the
  bytes of `host/surfacepage/passkey/*` read from disk, the five headers, `surfaceCSP`; POST 405; served
  with the world store closed; the page's one script, no handler attribute, form, image, link, iframe or
  URL; the labelled instrument-health controls: no `fetch(`, `XMLHttpRequest`, `innerHTML`, …).
- `host/daemon/surface_app_test.go`: `TestSurfaceAppRoutes`'s comment, and four more 404 names
  (`/surface/passkey/`, `/surface/passkey/index.html`, `/surface/passkey/passkey.js`,
  `/surface/passkey.html`): the two routes open no subtree.
- new `cmd/ailang-worldd/passkey_page_chrome_test.go`: **XB1** `TestPasskeyPageInChrome` (S4-D7).
- new `cmd/ailang-worldd/passkey_page_abort_test.go`: **XB2** `TestPasskeyPageAbortsAPendingCeremony`
  (S4-D6).

**Acceptance** (the executor has no Chrome: XB1/XB2 must SKIP, counted; §6 runs them required)
```
[base]
go test ./host/daemon/ -run '^(TestPasskeyPageRoutes|TestSurfaceAppRoutes|TestFrozenV1RouteTableMatchesMux|TestDaemonDependencyAllowlist)$' -count=1 -v   # no --- SKIP
go test ./cmd/ailang-worldd/ -run '^(TestPasskeyPageInChrome|TestPasskeyPageAbortsAPendingCeremony)$' -count=1 -v
#   exactly two "--- SKIP: TestPasskeyPage…" lines, each with "SKIP COUNTED: the passkey page browser test needs WORLD_CHROME"
go test ./host/daemon/ -count=1 -timeout 20m                                                              # ~5-8 min; the default 10 min timed out at load 30-60 (§1.2 C33)
```

**Expected cumulative numstat after the controller's commit of M7e: 924**.

### M7f — `generate.sh` made real · 130 test + generator · Run 4

**Patch** `M7f`: new `cmd/ailang-worldd/passkey_fixture_record_test.go`:
`TestPasskeyPageRecordsTheFixtures` (S4-D8; skips unless `WORLD_PASSKEY_RECORD_DIR`) and
`TestPasskeyFixtureGeneratorRefuses` (five arms: no directory, two, a missing one, a non-empty one, no
`WORLD_CHROME` → exit 2, the exact message, nothing written); `host/approval/testdata/webauthn/
generate.sh` (the real generator); `SOURCE.txt`'s last line. The committed fixtures do not change.

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestPasskeyFixtureGeneratorRefuses|TestPasskeyPageRecordsTheFixtures)$' -count=1 -v
#   TestPasskeyFixtureGeneratorRefuses --- PASS; exactly one --- SKIP: TestPasskeyPageRecordsTheFixtures
go test ./host/approval/ -run '^(TestPasskeyFixturesAreTheGeneratedBytes|TestApproveByAssertionAcceptsChrome)$' -count=1 -v   # no --- SKIP
```

**Expected cumulative numstat after the controller's commit of M7f: 1054**.

### M7g — QUICKSTART §16 *(attended — pending first verbatim run)* · docs only · Run 4

**Patch** `M7g`: `docs/QUICKSTART.md` gains §16 "Approve with a passkey": the §0 scratch store served
at `127.0.0.1:7644` and opened as `http://localhost:7644`; enrol (the refusals before any link; the
link; Create, Save; the 26-character code in five groups, typed or pasted, a fingerprint not a
secret; the y/N; the 5-minute code and 2-minute y/N bounds and the page's 4-minute timer); `episode
start … --request`; `--webauthn-link` (request, expiry, hash; compare with the page); Approve, Save;
`--webauthn-assertion` (mints, `approved_by=webauthn:demo-mac`); the same file again refused; `approve
--list --decided`; `approver revoke`; the refusals. **Not executed** (D-WORLD-109; design §10).

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestQuickstartAttendedStepsBoxMatchesSection9|TestQuickstartSection9FlagsMatchTheCLI|TestQuickstartSection10FlagsMatchTheCLI|TestQuickstartSection11FlagsMatchTheCLI|TestQuickstartServiceSectionMatchesTheCLIHelp|TestEveryVerbHasHelp)$' -count=1 -v   # no --- SKIP
go test ./cmd/ailang-worldd/ -count=1 -timeout 20m                                                        # ~3-5 min
```

**Expected cumulative numstat after the controller's commit of M7g: 1054** (docs are outside the filter).

---

## 4. Mutation table (69 rows run on the planner's final tree: 64 KILLED, 5 named equivalents SURVIVED as predicted)

**Derivation.** Anchored to Appendix P's production hunks — `approver_passkey.go` (WN2a), `host/surfacepage/passkey/passkey.js` and `index.html` (every condition, operand, bound, copied field, state change and the timer), `host/surfacepage/passkey.go` (both embeds), `host/daemon/daemon.go` and `surface_page.go` (each route, its method, each handler's headers, content type and body) and `generate.sh` (each refusal operand) — plus design §6's named mutants (the code over one input or truncated, the Save link removed, `userVerification: "preferred"` with a UV-false authenticator, the page using its own challenge — MUT-JS-B64-STRICT and the challenge's copy are covered by XB1's CLI refusal; no `signal`, another timer constant, saving left enabled after expiry, a late result accepted, the button left enabled while pending). **Procedure**: `plan-scratch/tools/mutate.py work tools/mutants.json logs/mutant_ledger.tsv <run> --workers 6` (outside the repo; `tools/mk_mutants.py` generates the table and asserts every anchor occurs exactly once in the final tree): for each row, in its own copy of the tree, apply **one** edit, assert the file changed, `go vet` the package, run **the single named test** (`go test <pkg> -run '^<test>$' -count=1`) **and then its package suite** — for a Go or `generate.sh` mutant the whole package (`go test <pkg> -count=1`), for a page-asset (JS/HTML) mutant `cmd/ailang-worldd`'s passkey suite `-run '^(TestPasskey|TestApproveByPasskey|TestApprovePasskey)'` (§1.2 C22) — every browser run with `WORLD_CHROME`, `WORLD_REQUIRE_PASSKEY_BROWSER=1` and `-skip '^TestMarketingChromeSixShots$'` (§1.1 C8), and a private `GOCACHE` (C32), restore, assert the restore is byte-identical, and **append one line to the ledger** `plan-scratch/logs/mutant_ledger.tsv` (UTC time, run label, id, stage, named, suite, expected, MATCH/MISMATCH). A row is KILLED only when both fail. The ledger keeps every run: `run1-final` and `run2-final` (stopped: the tests changed after them, C15, C21), `run3-final` (57 NOCOMPILE from the trimmed cache, C32), `run4-rerun` (stopped: C23, C24), and **`run5-final` + `run6-final-cmd`, the runs on the final tree** (the notes below say which line the table shows for each id). **Never commit or snapshot a mutant.** In the table tabs are dropped, a newline prints as ⏎ and a backquote as '.

| Id | Stage | Mutant (single edit) | Named test | Named | Suite |
|---|---|---|---|---|---|
| MUT-WN2A | M7a | approver_passkey.go: if !approval.ApproverNamePattern.MatchString(name) { → if false { | `TestPasskeyEnrollFlagRules` | KILLED | KILLED |
| MUT-JS-TIMEOUT | M7b | passkey.js: const TIMEOUT_MS = 240000; → const TIMEOUT_MS = 239000; | `TestPasskeyPageAbortsAPendingCeremony` | KILLED | KILLED |
| MUT-JS-B64-STRICT | M7b | passkey.js: if (!/^[A-Za-z0-9_-]*$/.test(s)) throw → if (false) throw | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-CODE-ONE-INPUT | M7b | passkey.js: msg.set(await sha256(clientDataJSON), label.length); → msg.set(await sha256(attestationObject), label.length); | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-CODE-NO-ATT | M7b | passkey.js:     msg.set(await sha256(attestationObject), label.length + 32); → (deleted) | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-CODE-LABEL | M7b | passkey.js: 'world-enrolment-code/v1\0' → 'world-enrolment-code/v1' | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-CODE-TRUNC | M7b | passkey.js: .slice(0, 16); → .slice(0, 15); | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-CODE-LAST | M7b | passkey.js:     out += alphabet[(bits << (5 - n)) & 31]; → (deleted) | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-CODE-GROUP | M7b | passkey.js: out.slice(20)] → out.slice(20, 25)] | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-REG-TYPE | M7b | passkey.js: type: c.type,⏎    response: { clientDataJSON: b64u(c.response.clientDataJSON), attestationObject → type: 'x',⏎    response: { clientDataJSON: b64u(c.response.clientDataJSON), attestationObject | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-REG-ATT | M7b | passkey.js: attestationObject: b64u(c.response.attestationObject) } }); → attestationObject: b64u(c.response.clientDataJSON) } }); | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-ASR-AUTHDATA | M7b | passkey.js: authenticatorData: b64u(c.response.authenticatorData), → authenticatorData: b64u(c.response.clientDataJSON), | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-ASR-SIG | M7b | passkey.js: signature: b64u(c.response.signature) → signature: b64u(c.response.authenticatorData) | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-ASR-USERHANDLE | M7b | passkey.js: userHandle: c.response.userHandle ? b64u(c.response.userHandle) : null → userHandle: null | `TestPasskeyPageInChrome` | SURVIVED | SURVIVED |
| MUT-JS-RP-DEFAULT | M7b | passkey.js: rp: { id: f.rp, name: 'AILANG World' } → rp: { name: 'AILANG World' } | `TestPasskeyPageInChrome` | SURVIVED | SURVIVED |
| MUT-JS-RPID-DEFAULT | M7b | passkey.js: { challenge, rpId: f.rp, → { challenge, | `TestPasskeyPageInChrome` | SURVIVED | SURVIVED |
| MUT-JS-USER-ID | M7b | passkey.js: user: { id: new TextEncoder().encode(f.name), → user: { id: new TextEncoder().encode('x'), | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-ALG-NO-ES256 | M7b | passkey.js: [{ type: 'public-key', alg: -7 }, { type: 'public-key', alg: -8 }] → [{ type: 'public-key', alg: -8 }] | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-ALG-NO-EDDSA | M7b | passkey.js: [{ type: 'public-key', alg: -7 }, { type: 'public-key', alg: -8 }] → [{ type: 'public-key', alg: -7 }] | `TestPasskeyPageInChrome` | SURVIVED | SURVIVED |
| MUT-JS-UV-CREATE | M7b | passkey.js: authenticatorSelection: { userVerification: 'required' → authenticatorSelection: { userVerification: 'preferred' | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-UV-GET | M7b | passkey.js: userVerification: 'required', timeout: TIMEOUT_MS }, signal }); → userVerification: 'preferred', timeout: TIMEOUT_MS }, signal }); | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-ATTESTATION | M7b | passkey.js: attestation: 'none' → attestation: 'direct' | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-ALLOW-EMPTY | M7b | passkey.js: allowCredentials: f.allow.map((id) => ({ type: 'public-key', id: unb64u(id) })), → allowCredentials: [], | `TestPasskeyPageInChrome` | SURVIVED | SURVIVED |
| MUT-JS-SIGNAL-CREATE | M7b | passkey.js: attestation: 'none', timeout: TIMEOUT_MS }, signal }); → attestation: 'none', timeout: TIMEOUT_MS } }); | `TestPasskeyPageAbortsAPendingCeremony` | KILLED | KILLED |
| MUT-JS-SIGNAL-GET | M7b | passkey.js: userVerification: 'required', timeout: TIMEOUT_MS }, signal }); → userVerification: 'required', timeout: TIMEOUT_MS } }); | `TestPasskeyPageAbortsAPendingCeremony` | KILLED | KILLED |
| MUT-JS-CALL-OP | M7b | passkey.js: if (f.op === 'create') {⏎      return navigator → if (true) {⏎      return navigator | `TestPasskeyPageAbortsAPendingCeremony` | KILLED | KILLED |
| MUT-JS-BTN-PENDING | M7b | passkey.js:     go.disabled = true;⏎    setState('pending'); →     setState('pending'); | `TestPasskeyPageAbortsAPendingCeremony` | KILLED | KILLED |
| MUT-JS-EXPIRED-FLAG | M7b | passkey.js:       expired = true; → (deleted) | `TestPasskeyPageAbortsAPendingCeremony` | KILLED | KILLED |
| MUT-JS-ABORT | M7b | passkey.js:       controller.abort(); → (deleted) | `TestPasskeyPageAbortsAPendingCeremony` | KILLED | KILLED |
| MUT-JS-STATE-EXPIRED | M7b | passkey.js:       setState('expired'); → (deleted) | `TestPasskeyPageAbortsAPendingCeremony` | KILLED | KILLED |
| MUT-JS-SAY-TIMEOUT | M7b | passkey.js:       say('Timed out: →       void ('Timed out: | `TestPasskeyPageAbortsAPendingCeremony` | KILLED | KILLED |
| MUT-JS-CLEAR | M7b | passkey.js:       clearTimeout(timer); → (deleted) | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-LATE | M7b | passkey.js:     if (expired) { // [r2] →     if (false) { // [r2] | `TestPasskeyPageAbortsAPendingCeremony` | KILLED | KILLED |
| MUT-JS-RETRY-EXPIRED | M7b | passkey.js:       go.disabled = false;⏎      return;⏎    }⏎    if (failure !== null) { →       return;⏎    }⏎    if (failure !== null) { | `TestPasskeyPageAbortsAPendingCeremony` | KILLED | KILLED |
| MUT-JS-FAIL | M7b | passkey.js:     if (failure !== null) { →     if (false) { | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-FAIL-STATE | M7b | passkey.js:       setState('failed'); → (deleted) | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-FAIL-SAY | M7b | passkey.js:       say('The browser did not make →       void ('The browser did not make | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-RETRY-FAILED | M7b | passkey.js:       go.disabled = false;⏎      return;⏎    }⏎    const create →       return;⏎    }⏎    const create | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-CREATE-OP | M7b | passkey.js: const create = f.op === 'create'; → const create = true; | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-RESULT | M7b | passkey.js: $('result').value = json; → $('result').value = ''; | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-SAVE-NEWLINE | M7b | passkey.js: new Blob([json + '\n'] → new Blob([json] | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-SAVE-NAME-REG | M7b | passkey.js: 'world-registration-${f.name}.json' → 'world-registration.json' | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-SAVE-NAME-ASR | M7b | passkey.js: 'world-assertion-${f.request}.json' → 'world-assertion.json' | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-SAVE-SHOW | M7b | passkey.js:     save.hidden = false; → (deleted) | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-DONE | M7b | passkey.js:     setState('done'); → (deleted) | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-V | M7b | passkey.js: if (f.v !== 1 \|\| (f.op !== 'create' && f.op !== 'get')) → if ((f.op !== 'create' && f.op !== 'get')) | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-OP | M7b | passkey.js: if (f.v !== 1 \|\| (f.op !== 'create' && f.op !== 'get')) → if (f.v !== 1) | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-ORIGIN | M7b | passkey.js: if (f.origin !== location.origin \|\| f.rp !== location.hostname) { → if (f.rp !== location.hostname) { | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-RP | M7b | passkey.js: if (f.origin !== location.origin \|\| f.rp !== location.hostname) { → if (f.origin !== location.origin) { | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-WHAT-OP | M7b | passkey.js: $('what').textContent = f.op === 'create' → $('what').textContent = true | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-READY-SHOW | M7b | passkey.js:   go.hidden = false;⏎  go.disabled = false;⏎  setState('ready'); →   go.disabled = false;⏎  setState('ready'); | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-JS-READY-ENABLE | M7b | passkey.js:   go.hidden = false;⏎  go.disabled = false;⏎  setState('ready'); →   go.hidden = false;⏎  setState('ready'); | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-HTML-SCRIPT | M7b | index.html: <script src="/surface/passkey.js"></script> → (deleted) | `TestPasskeyPageInChrome` | KILLED | KILLED |
| MUT-HTML-SAVE-HIDDEN | M7b | index.html: <a id="save" hidden>Save</a> → <a id="save">Save</a> | `TestPasskeyPageAbortsAPendingCeremony` | KILLED | KILLED |
| MUT-EMBED-PAGE | M7b | passkey.go: //go:embed passkey/index.html → //go:embed passkey/passkey.js | `TestPasskeyPageRoutes` | KILLED | KILLED |
| MUT-EMBED-SCRIPT | M7b | passkey.go: //go:embed passkey/passkey.js → //go:embed passkey/index.html | `TestPasskeyPageRoutes` | KILLED | KILLED |
| MUT-ROUTE-PAGE | M7e | daemon.go: mux.HandleFunc("GET /surface/passkey", d.handleSurfacePasskey) → (deleted) | `TestPasskeyPageRoutes` | KILLED | KILLED |
| MUT-ROUTE-SCRIPT | M7e | daemon.go: mux.HandleFunc("GET /surface/passkey.js", d.handleSurfacePasskeyScript) → (deleted) | `TestPasskeyPageRoutes` | KILLED | KILLED |
| MUT-ROUTE-PAGE-METHOD | M7e | daemon.go: "GET /surface/passkey", d.handleSurfacePasskey) → "/surface/passkey", d.handleSurfacePasskey) | `TestPasskeyPageRoutes` | KILLED | KILLED |
| MUT-ROUTE-SCRIPT-METHOD | M7e | daemon.go: "GET /surface/passkey.js", d.handleSurfacePasskeyScript) → "/surface/passkey.js", d.handleSurfacePasskeyScript) | `TestPasskeyPageRoutes` | KILLED | KILLED |
| MUT-PAGE-HEADERS | M7e | surface_page.go: setSurfaceHeaders(w)⏎_, _ = w.Write(surfacepage.PasskeyPage) → _, _ = w.Write(surfacepage.PasskeyPage) | `TestPasskeyPageRoutes` | KILLED | KILLED |
| MUT-SCRIPT-HEADERS | M7e | surface_page.go: setSurfaceHeaders(w)⏎w.Header().Set("Content-Type", "text/javascript; charset=utf-8") → w.Header().Set("Content-Type", "text/javascript; charset=utf-8") | `TestPasskeyPageRoutes` | KILLED | KILLED |
| MUT-SCRIPT-CONTENT-TYPE | M7e | surface_page.go: w.Header().Set("Content-Type", "text/javascript; charset=utf-8")⏎_, _ = w.Write(surfacepage.PasskeyScript) → _, _ = w.Write(surfacepage.PasskeyScript) | `TestPasskeyPageRoutes` | KILLED | KILLED |
| MUT-PAGE-BODY | M7e | surface_page.go: _, _ = w.Write(surfacepage.PasskeyPage) → _, _ = w.Write(surfacepage.PasskeyScript) | `TestPasskeyPageRoutes` | KILLED | KILLED |
| MUT-SCRIPT-BODY | M7e | surface_page.go: _, _ = w.Write(surfacepage.PasskeyScript) → _, _ = w.Write(surfacepage.PasskeyPage) | `TestPasskeyPageRoutes` | KILLED | KILLED |
| MUT-GEN-ARGC | M7f | generate.sh: if [ "$#" -ne 1 ]; then → if false; then | `TestPasskeyFixtureGeneratorRefuses` | KILLED | KILLED |
| MUT-GEN-NOT-DIR | M7f | generate.sh: if [ ! -d "$1" ] \|\| [ -n "$(ls -A "$1")" ]; then → if [ -n "$(ls -A "$1")" ]; then | `TestPasskeyFixtureGeneratorRefuses` | KILLED | KILLED |
| MUT-GEN-NOT-EMPTY | M7f | generate.sh: if [ ! -d "$1" ] \|\| [ -n "$(ls -A "$1")" ]; then → if [ ! -d "$1" ]; then | `TestPasskeyFixtureGeneratorRefuses` | KILLED | KILLED |
| MUT-GEN-CHROME | M7f | generate.sh: if [ -z "${WORLD_CHROME:-}" ]; then → if false; then | `TestPasskeyFixtureGeneratorRefuses` | KILLED | KILLED |

**Named equivalents** (each SURVIVED by the named test and the suite, as predicted):
- **MUT-JS-ASR-USERHANDLE** — the assertion law requires a present userHandle to equal the name and accepts an absent one (W-userHandle is optional in L3 §7.2); the page forwards what the authenticator returned.
- **MUT-JS-RP-DEFAULT** — an omitted rp.id defaults to the page's hostname, which the page has already required to equal f.rp (the rp check, MUT-JS-RP).
- **MUT-JS-RPID-DEFAULT** — an omitted rpId defaults to the page's hostname, which the page has already required to equal f.rp.
- **MUT-JS-ALG-NO-EDDSA** — -8 is offered for an authenticator without ES256; Chrome's virtual authenticator takes the first, ES256, so the second entry is not observable by XB1 (generate.sh's EdDSA recording filters the list to -8 and is its check, attended).
- **MUT-JS-ALLOW-EMPTY** — the virtual authenticator holds one discoverable credential, which an empty allow list also selects; the CLI judges the credential ID (WF5), so a wrong one is refused there.

**Ledger.** `plan-scratch/logs/mutant_ledger.tsv`, 193 lines: run1-final 19 and run2-final 7 (stopped when XB2's deadline measurement and the download-directory check changed, C15, C21), run3-final 69 (57 NOCOMPILE: the shared Go cache was trimmed by another session mid-run, C32), run4-rerun 29 (stopped: it exposed XB1's pending race, C23, and that `setUserVerified(false)` could not kill the two `userVerification: "preferred"` mutants, C24 — both fixed in the tests, not the page), **run5-final 69, all as expected**, and **run6-final-cmd 58**: every `cmd/ailang-worldd` mutant re-run after the harness's stderr became a `syncBuffer` (C34). The table above is the final tree's verdict for each id: run6's line for the 58 `cmd/ailang-worldd` mutants, run5's for the 11 `host/daemon` ones (no file of `host/daemon`, `host/surfacepage` or their tests changed after run5 started).

---


## 5. Files (final `git status --short` after M7g)

| File | Sub-milestone | Change |
|---|---|---|
| `cmd/ailang-worldd/approver_passkey.go` | M7a | WN2a (one branch) |
| `cmd/ailang-worldd/approver_passkey_test.go`, `approve_passkey_test.go` | M7a | one CW1 arm; CW7's next-step line pinned |
| `host/surfacepage/passkey/index.html`, `passkey/passkey.js`, `passkey.go` | M7b | new: the page, its script, the embed |
| `cmd/ailang-worldd/testdata/passkey_cdp.mjs` | M7c | new: the DevTools driver |
| `cmd/ailang-worldd/passkey_page_driver_test.go` | M7d | new: the Go harness |
| `host/daemon/daemon.go`, `surface_page.go` | M7e | two routes, two handlers |
| `host/daemon/passkey_page_test.go`, `surface_app_test.go` | M7e | new PG1; four 404 names |
| `cmd/ailang-worldd/passkey_page_chrome_test.go`, `passkey_page_abort_test.go` | M7e | new: XB1, XB2 |
| `cmd/ailang-worldd/passkey_fixture_record_test.go` | M7f | new: the recorder, the generator's refusals |
| `host/approval/testdata/webauthn/generate.sh`, `SOURCE.txt` | M7f | the real generator; one line |
| `docs/QUICKSTART.md` | M7g | §16 |
| this plan and `sprint_w-world-approvals-by-key-a2-s4.json` | planner | |

---

## 6. Sprint gates (controller, after the last sub-milestone; each measured at the base, §1.1)

In **one** `bash` shell, from the worktree, **not under `nohup`**. Every leg `tee`s its full output to
a log; a verdict is read from its log, never from `tail` alone. **`WORLD_CHROME` is set on the XB leg
only, never exported** (§1.1 C8). The XB leg runs the three browser-facing tests **three times in a
row**, required:

```
bash <<'EOF'
set -euo pipefail
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export PATH="/opt/homebrew/bin:$PATH"     # node v26 for the XB leg
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter382-row201-a2s4; cd "$W"
L=$HOME/.ailang/state/world-iter382/gates; mkdir -p "$L"
CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
./scripts/verify_ail.sh 2>&1 | tee "$L/verify_ail.log"                      # no .ail change: exactly green
out=$(gofmt -l cmd host); echo "$out" | tee "$L/gofmt.log"; test -z "$out"
go build ./... 2>&1 | tee "$L/build.log"
go vet ./... 2>&1 | tee "$L/vet.log"
go test -run '^$' ./... 2>&1 | tee "$L/fence.log"                           # compile fence for every _test.go
go test -count=1 -timeout 20m ./host/surfacepage/... ./host/approval/... ./cmd/ailang-worldd/ ./host/daemon/ 2>&1 | tee "$L/pkgs.log"
go test -race -count=1 ./host/surfacepage/ 2>&1 | tee "$L/race.log"
go test -race -count=1 ./host/daemon/ -run '^(TestPasskeyPageRoutes|TestSurfaceAppRoutes|TestSurfaceScriptRoute)$' 2>&1 | tee -a "$L/race.log"
go test -race -count=1 ./cmd/ailang-worldd/ -run '^(TestPasskey|TestApproveByPasskey|TestApprovePasskey|TestApprover|TestEveryVerbHasHelp)' 2>&1 | tee -a "$L/race.log"
# the repo-wide inventories (design §7), the reflective, help and quickstart tests, by name, -v:
go test -v -count=1 ./host/broker/ -run '^TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential$' 2>&1 | tee "$L/inv.log"
go test -v -count=1 ./host/daemon/ -run '^(TestDaemonDependencyAllowlist|TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestDaemonConstructorCallerInventory|TestQuickstartDecisionPacketSection|TestPasskeyPageRoutes)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/boundary/ -run '^(TestWorldBoundaryDependencyAllowlist|TestBareNetHTTPExemptionIsPerGroup|TestEveryCommandDirectoryIsAProtectedGroup)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/store/ -run '^(TestProductionContextRoots|TestProductionGoSurface|TestNoNewDeadlineFreeStoreReads|TestReadGettersHonorContext|TestCredentialsMethodsRequireADeadline)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/verifygate/ -run '^(TestNoParallelWriteForkPackagesOutsideVerifygate|TestSubprocessSinkGate)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/approval/... -run '^(TestApprovalImportsOnlyTheStandardLibrary|TestApprovalTestIsTestOnly|TestPasskeyFixturesAreTheGeneratedBytes)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./cmd/ailang-worldd/ -run '^(TestEveryVerbHasHelp|TestQuickstartAttendedStepsBoxMatchesSection9|TestQuickstartSection9FlagsMatchTheCLI|TestQuickstartSection10FlagsMatchTheCLI|TestQuickstartSection11FlagsMatchTheCLI|TestQuickstartServiceSectionMatchesTheCLIHelp|TestPasskeyFixtureGeneratorRefuses)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/broker/ -run '^(TestQuickstartSection10ProfilesLoad|TestQuickstartProtectedContract)$' 2>&1 | tee -a "$L/inv.log"
for t in TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux \
  TestSurfaceAppRoutes TestDaemonConstructorCallerInventory TestQuickstartDecisionPacketSection TestPasskeyPageRoutes \
  TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup \
  TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext \
  TestCredentialsMethodsRequireADeadline TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestApprovalImportsOnlyTheStandardLibrary \
  TestApprovalTestIsTestOnly TestPasskeyFixturesAreTheGeneratedBytes TestEveryVerbHasHelp TestQuickstartAttendedStepsBoxMatchesSection9 \
  TestQuickstartSection9FlagsMatchTheCLI TestQuickstartSection10FlagsMatchTheCLI TestQuickstartSection11FlagsMatchTheCLI \
  TestQuickstartServiceSectionMatchesTheCLIHelp TestPasskeyFixtureGeneratorRefuses TestQuickstartSection10ProfilesLoad TestQuickstartProtectedContract; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || { echo "$t did not PASS"; exit 1; }; done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then exit 1; fi
# XB1 and XB2, required, three consecutive runs (Chrome headless, temp profile, virtual authenticator):
for i in 1 2 3; do
  s=$(date +%s)
  WORLD_CHROME="$CHROME" WORLD_REQUIRE_PASSKEY_BROWSER=1 go test -v -count=1 ./cmd/ailang-worldd/ -run '^(TestPasskeyPageInChrome|TestPasskeyPageAbortsAPendingCeremony)$' 2>&1 | tee "$L/xb_$i.log"
  echo "XB run $i: $(( $(date +%s) - s )) s" | tee -a "$L/xb_times.log"
  for t in TestPasskeyPageInChrome TestPasskeyPageAbortsAPendingCeremony; do grep -qE -- "^--- PASS: $t " "$L/xb_$i.log" || { echo "$t did not PASS (run $i)"; exit 1; }; done
  if grep -E '^\s*--- (SKIP|FAIL)' "$L/xb_$i.log"; then exit 1; fi
  test "$(grep -cE '^    --- PASS: TestPasskeyPage' "$L/xb_$i.log")" = 7                   # 3 XB1 + 4 XB2 subtests
done
echo PRE_COMMIT_GATES_GREEN | tee "$L/pre_commit.ok"
EOF
# AFTER the controller commits (the tracked-binary gate reads committed blobs; this leg runs go test ./...):
bash -c 'set -o pipefail; export AILANG_BIN=$HOME/.pinned-ailang/ailang; cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter382-row201-a2s4 && VERIFY_GO_LEGS=plain ./scripts/verify_go.sh 2>&1 | tee $HOME/.ailang/state/world-iter382/gates/verify_go.log; echo "verify_go rc=$?" | tee -a $HOME/.ailang/state/world-iter382/gates/verify_go.log'
git -C /Users/voightkampff/dev/sunholo-data/.wt-world-iter382-row201-a2s4 status --short   # only the two plan files, until the controller commits them
```

**Optional (the controller's choice; not a gate):** `generate.sh` against the committed driver, into
an empty directory outside `/tmp`, then the law over its output in a throwaway copy (§1.2 C16):
`mkdir -p $HOME/.ailang/state/world-iter382/gen && WORLD_CHROME="$CHROME" sh
host/approval/testdata/webauthn/generate.sh $HOME/.ailang/state/world-iter382/gen` (~20 s). Never copy
its output over the committed fixtures in this sprint.

**At the base**: every pre-commit leg but the new tests is green at `27ce163e` in 475 s (§1.1; 26
named tests `--- PASS`). On the final tree: `PRE_COMMIT_GATES_GREEN` in ~7 min at load 7–12 (29 named tests `--- PASS`, XB ×3 green, 7–9 s each; §1.2 C19, C25, C26, C28), and the after-commit `verify_go.sh` **green, rc 0** (C27). At load 56–80 the same script was red only on `TestWorldContextAndStatusline` (C26), and the base's `verify_go.sh` at load 40–66 was red in six packages (§1.1): **run the gates on a quiet machine, and re-run a red untouched package alone once.**

**Known load- or environment-sensitive (re-run once, alone; do not chase)**:
`TestMissionConfigDefaultFlowIntegration/invalid_ledger` (and `/valid_ledger` under heavy load),
`TestSeedRestartIgnoresASurvivingOldCopy/survivor`, `host/broker`'s exec-probe timeouts and its local
`-race`, the local `-race` of `host/daemon` and `host/verifygate`, `TestHostsLaunchersExecHost`,
`TestDemoSeedPreflightRefusals`, `TestWorktreeSnapshotRefusesAConfiguredFilter`,
`scripts/marketingcheck`'s supervisor tests, `host/daemon`'s `TestAGUISlowReaderCutIsResumable`,
`TestSessionSelfGit`, `TestServeStartsWhenTheProbeSpendsTheStartupClock`, `TestRegistryPhaseHasItsOwnClock`,
`cmd/ailang-worldd`'s `TestWorldContextAndStatusline` (a 5 s subprocess and a minute-boundary cache line: red 3 times in the planner's runs at load 56–80, every time green alone, on trees that differ only in files it does not read), `TestCLIRealSubprocessEpisode`,
and anything run under `nohup`. Any other red, any race report, any red in `cmd/ailang-worldd`,
`host/surfacepage`, `host/approval` or PG1/`TestSurfaceAppRoutes`, or **any red XB run**, is a
regression: a flaky browser test is a reason to stop and report, never to loosen it. CI's
`go-race-verify` is the race verdict; CI does not run XB1/XB2 (design §9), so the XB leg above is their
only gate.

**Must stay green, unedited**: every a0a1, A1b, sprint-1, sprint-2 and sprint-3 test, except the named
edits (M7a's CW1 arm and CW7 line, M7e's `TestSurfaceAppRoutes` 404 list). In particular XW0
(`TestPasskeyFixturesAreTheGeneratedBytes`: the fixtures are unchanged), CW1–CW9, `TestSurfaceAppRoutes`'s
other subtests, `TestFrozenV1RouteTableMatchesMux` (13 routes), and every §7 inventory.

---

## 7. Executor directive notes

- **No git writes.** Never `git add`, `commit`, `stash`, `checkout`, `reset`, `restore`, `rebase`, `push`,
  `worktree` or `apply`. `patch -p1` (§2) is the only writer; read-only `git status`/`git diff` are fine.
- **No browser.** Never run Chrome, `node`, `generate.sh`, or any test with `WORLD_CHROME` set. M7c's
  acceptance only reads the file. M7e's acceptance expects XB1/XB2 to **skip, counted** (the gate
  message); the controller runs them required in §6.
- **Runs.** Run 1 = M7a, M7b; Run 2 = M7c, M7d; Run 3 = M7e; Run 4 = M7f, M7g. Each run reads §0, §2
  and its own §3 sections, then applies its patches in order.
- **Expected numstat after each sub-milestone, measured AFTER the controller's commit of it** (an
  uncommitted tree hides the new untracked files from numstat):
  ```
  bash -c 'cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter382-row201-a2s4 && git diff --numstat 27ce163eb1ae8af55151ea6548591c15a3126835 HEAD -- cmd host | awk "{s+=\$1+\$2} END {print s}"'
  ```
  After the controller's commit of M7a **11**, M7b **195**, M7c **367**, M7d **597**,
  M7e **924**, M7f **1054**, M7g **1054**. **Any other value means a patch was not applied as
  written: stop and report.** (M7g changes only `docs/`, outside the filter.)
- **Snapshot after every sub-milestone, cumulatively**:
  ```
  bash -c 'set -euo pipefail; M=M7a; W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter382-row201-a2s4
    D=$HOME/.ailang/state/world-iter382/snap/$M; rm -rf "$D"; mkdir -p "$D"; cd "$W"
    git status --porcelain --untracked-files=all | while read -r st f; do
      case "$f" in design_docs/planned/w-world-approvals-by-key-a2-sprint4-plan.md|design_docs/planned/sprint_w-world-approvals-by-key-a2-s4.json) continue;; esac
      mkdir -p "$D/$(dirname "$f")"; cp -p "$f" "$D/$f"; done
    (cd "$D" && find . -type f ! -name MANIFEST.sha256 | sort | xargs shasum -a 256) > "$D/MANIFEST.sha256"
    cat "$D/MANIFEST.sha256"'
  ```
  (`M` = `M7a` … `M7g`.) No file is deleted, so the porcelain status has only `M` and `??`; stop on any
  other status. The controller rebuilds one commit per sub-milestone (iteration 377's `commit_ms.sh`
  with `WT` and the snapshot root pointed here).
- **Stop conditions.** A `patch --dry-run` failure; a numstat other than the expected value; a red
  earlier-milestone test; a `--- FAIL`, or a `--- SKIP` other than the ones §3 names; any need to touch
  a file outside §5; running Chrome, node, `ssh-keygen` or any browser; any reference to `~/.ssh`,
  `~/.ailang/world` or a real `/dev/tty`. Stop and report; do not work around.
- **Timing.** M7a's and M7g's whole `cmd/ailang-worldd` ~3–4 min under load; M7e's whole `host/daemon`
  ~5–8 min; every other acceptance under 1 min; the compile fence ~25 s.
- **Shell.** Wrap loops, globs and `$VAR` expansions in `bash -c` or a `bash <<'EOF'` heredoc (zsh does
  not word-split). Export `AILANG_BIN` in the same shell as the Go tests. `tee` every gate's output to a
  file; never pipe a gate's only copy through `tail`. **Never start a gate under `nohup`.** Poll any
  background command's artifact in a bounded `date +%s` loop in the same turn; never end a turn
  announcing a wait.

---

## 8. Risks

- **R1 — Browser flakiness under load.** XB1/XB2 depend on Chrome, node and a loaded machine. Measured
  4–19 s per test on the planner's machine at load 20–80 (C19); `passkeyBrowserBudget` is 120 s and the
  driver's own waits are 15–30 s. The four races the planner hit (C12–C15) are fixed in the driver,
  not by loosening an assertion; the stray download (C21) is logged and not counted as a save. A red XB run in §6 stops the sprint (§6).
- **R2 — CI never runs XB1/XB2** (design §9). The committed state is checked by PG1 and the counted
  skips in CI; the browser behaviour only by §6 on the controller's machine. A future change to the page
  must re-run §6's XB leg.
- **R3 — Safari and Firefox are UNVERIFIED** (design R3): the page uses only `navigator.credentials`,
  `AbortController`, `crypto.subtle.digest`, `Blob`/`URL.createObjectURL` and `<a download>`; the attended
  run (design §10, QUICKSTART §16) is the first non-Chrome browser.
- **R4 — A re-recording with `generate.sh` is a re-pin, not a drop-in** (C16, C18): XW0's two SHA-256s,
  the two golden values, `TestPasskeyEnvelope`'s `!HasUserHandle` and, if Chrome adds
  `other_keys_can_be_added_here`, XW1. Nothing in this sprint re-records.
- **R5 — The page trusts nothing it shows**: its text is `textContent`; it acts only on its own origin's
  links; it sends nothing; the CLI judges every byte it saves (WN8–WN11, WF2–WF6). A malicious link at
  the right origin can only ask the operator's authenticator to sign a challenge the operator then
  carries to a CLI that refuses anything but the pending request's own hash.
- **R6 — Haiku and patches.** The executor transcribes nothing; a skipped patch fails `--dry-run` or the
  numstat; a hand "fix" after a red test is forbidden by §7.

---

## 9. Acceptance criteria (the judge scores against these)

1. **The routes**: `GET /surface/passkey` and `/surface/passkey.js` serve exactly
   `host/surfacepage/passkey/{index.html,passkey.js}` with `surfaceCSP` (unchanged), `no-store`,
   `nosniff`, `no-referrer` and the right content types; POST is 405; they answer with the world store
   closed; no other name under them is served (PG1, `TestSurfaceAppRoutes`); `TestFrozenV1RouteTableMatchesMux`
   stays at 13.
2. **The page makes no request and writes no markup**: no `fetch`, `XMLHttpRequest`, form, `innerHTML`
   (PG1's labelled controls); XB1 and XB2 record **zero** CSP violations and page errors under the
   production `surfaceCSP`.
3. **XB1 — the product's path end to end**: the link `approver enroll --format webauthn` printed, opened
   in Chrome with a virtual authenticator, creates a passkey; the page's grouped code equals Go's
   `EnrolmentCode` of the saved file; the `blob:` download completes into the `--registration` path,
   byte-equal to the page's JSON plus `\n`; the CLI enrols it with that code (`exitOK`, ES256, the
   origin and rp). The link `approve --webauthn-link` printed shows the request, expiry and hash; a get on an
   authenticator without user verification saves nothing; the retry saves the assertion, and `approve
   --webauthn-assertion` mints with no terminal (`approved_by = webauthn:mark-mac`, counter advanced).
4. **XB1 — refusals**: a create on an authenticator without user verification fails on the page, nothing is saved, nothing is
   enrolled; links for another origin, another rp, another version or op, not base64url or not JSON are
   refused before any button; the control link is `ready`.
5. **XB2 — the [r2] timer**: for create and get, an operation that never answers is aborted through the
   `signal` the page passed, at exactly 240000 ms of page time (`Date.now()`), the page shows the
   explicit timeout error, nothing can be saved (no textarea, code, Save, download or file), a result
   arriving later (245000 ms) is ignored, the button was disabled while pending and is enabled after;
   a finished attempt's timer never fires (XB1's 250 s advance).
6. **Gating**: without `WORLD_CHROME`/node XB1 and XB2 skip with `SKIP COUNTED`; with
   `WORLD_REQUIRE_PASSKEY_BROWSER=1` they fail instead; §6 runs them required, three consecutive green
   runs with timings.
7. **`generate.sh`** records, through the page and a virtual authenticator, fixtures the production law
   accepts (C16), only into an empty directory it is given; it refuses (exit 2, nothing written, nothing
   run) every other call and a missing `WORLD_CHROME` (`TestPasskeyFixtureGeneratorRefuses`); the
   committed fixtures and XW0 are unchanged.
8. **WN2a**: a `--name` outside the grammar refuses before any store, terminal or link (CW1); CW7 pins
   the next-step line's `--db … --credentials-db …` from the request's own stores.
9. **QUICKSTART §16** exists, is tagged *(attended — pending first verbatim run)*, shows the
   26-character code in five groups, says it may be typed or pasted and is a fingerprint not a secret,
   and walks serve → enrol → request → link → assertion → decided list → revoke → a refused replay; it
   was not executed.
10. Each §4 mutant, applied alone, makes its named test **and** its package suite fail; the named
    equivalents survive for the stated reasons; the ledger has a line per run.
11. Every §6 gate passes, including every inventory, reflective, help and quickstart test by name, the
    `-race` legs, the XB leg ×3, and `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` after the commit
    (modulo §6's re-run rule); `git status --short` lists only §5.
12. The numstat after each sub-milestone's commit equals §7's value; the sprint total is ≤ 2,000.
13. `host/approval` stays pure and standard-library-only (no Go file changed); no new production exec
    site, module, package, `/v1` route or context root.
14. No real key, browser profile, authenticator, store, `/dev/tty`, `~/.ssh` access, mint, publish or
    `tailscale serve` anywhere in the sprint (D-WORLD-109).

---

## Appendix C — after this sprint

- **A2's attended exit run** (design §10, Mark, NOT the loop) is now possible: QUICKSTART §16 is its
  script; it removes the *attended — pending* tag when run verbatim, and records the public values
  design §10 step 6 lists (browser, provider, `fmt`, algorithm, BE/BS, counter, canonical CBOR).
- **Re-recording the fixtures** (only if wanted): `WORLD_CHROME=… sh
  host/approval/testdata/webauthn/generate.sh <empty dir>`; copying the output over
  `chrome-es256.json`/`chrome-eddsa.json` must re-pin in the same commit XW0's two SHA-256s,
  `TestEnrolmentCodeGolden`, `TestRenderPasskeyEnrolmentGolden`, `TestPasskeyEnvelope`'s
  `!HasUserHandle` (page recordings carry `userHandle`), and update `SOURCE.txt`; re-run it if Chrome
  injected `other_keys_can_be_added_here` and XW1 goes red (C16–C18).
- **A3** (design §2.14) reuses this page and script unchanged: the needs-you card, the Approve POST, the
  move into the surface app.

---

## Appendix P — the patches (apply with §2; generated from the compiled scratch stages; do not edit)

### Patch M7a — carry-forwards: WN2a and the next-step line

````diff M7a
--- a/cmd/ailang-worldd/approve_passkey_test.go
+++ b/cmd/ailang-worldd/approve_passkey_test.go
@@ -135,6 +135,10 @@
 		}
 		return l
 	}
+	req, err := approval.Decode(row.Request) // the next step names the request's own stores, --db first
+	if err != nil {
+		t.Fatal(err)
+	}
 	want := []ceremonyLink{get(pkRPID, pkOrigin, other.CredentialID(), passkeySigner.CredentialID()), get("x.ts.net", "https://x.ts.net", third.CredentialID())}
 	out := stdout.String()
 	if n := strings.Count(out, "/surface/passkey#"); n != 2 {
@@ -145,7 +149,7 @@
 		t.Fatalf("links %+v, want %+v", got, want)
 	}
 	for _, w := range []string{"approve request " + a.id + " (expires " + approverTime(row.ExpiresAt) + ", message hash " + hex.EncodeToString(h[:8]) + ")",
-		"\n  " + pkOrigin + "/surface/passkey#", "\n  https://x.ts.net/surface/passkey#", "ailang-worldd approve " + a.id + " --webauthn-assertion <the saved file> --db "} {
+		"\n  " + pkOrigin + "/surface/passkey#", "\n  https://x.ts.net/surface/passkey#", "\n  ailang-worldd approve " + a.id + " --webauthn-assertion <the saved file> --db " + shellQuote(req.WorldDB) + " --credentials-db " + shellQuote(req.CredentialsDB) + "\n"} {
 		if !strings.Contains(out, w) {
 			t.Fatalf("stdout lacks %q:\n%s", w, out)
 		}
--- a/cmd/ailang-worldd/approver_passkey.go
+++ b/cmd/ailang-worldd/approver_passkey.go
@@ -59,6 +59,9 @@
 	if name == "" || origin == "" || regPath == "" { // WN2
 		return approverRefuse(stderr, prefix, "--name, --origin and --registration are required with --format webauthn")
 	}
+	if !approval.ApproverNamePattern.MatchString(name) { // WN2a: before any link; the law's R1 judges it again
+		return approverRefuse(stderr, prefix, "--name %q: %v", name, approval.ErrApproverName)
+	}
 	credPath, refusal := approverStores(db, creds) // N3
 	if refusal != "" {
 		return approverRefuse(stderr, prefix, "%s", refusal)
--- a/cmd/ailang-worldd/approver_passkey_test.go
+++ b/cmd/ailang-worldd/approver_passkey_test.go
@@ -121,6 +121,8 @@
 			want: []string{"--name, --origin and --registration are required with --format webauthn"}},
 		{name: "WN2-registration", args: func(t *testing.T, r enrolRig) []string { return dropFlag(r.args(), "--registration") },
 			want: []string{"--name, --origin and --registration are required with --format webauthn"}},
+		{name: "WN2a-name-grammar", args: func(t *testing.T, r enrolRig) []string { return replaceFlag(r.args(), "--name", "Mark") },
+			want: []string{`--name "Mark": approval: approver name must match`}},
 		{name: "N3-no-world-store", args: func(t *testing.T, r enrolRig) []string {
 			return replaceFlag(r.args(), "--db", filepath.Join(t.TempDir(), "absent.db"))
 		}, want: []string{"does not exist; refusing to create one"}},
````

### Patch M7b — the page and its embed, not yet served

````diff M7b
--- a/host/surfacepage/passkey.go
+++ b/host/surfacepage/passkey.go
@@ -0,0 +1,18 @@
+package surfacepage
+
+import _ "embed"
+
+// Row 201 A2 (design_docs/planned/w-world-approvals-by-key-a2.md §2.14): the passkey ceremony
+// page, served at GET /surface/passkey and GET /surface/passkey.js. It reads its link's fragment,
+// runs one WebAuthn create or get, and offers the result as a file; it makes no request and holds
+// no data.
+
+// PasskeyPage is GET /surface/passkey.
+//
+//go:embed passkey/index.html
+var PasskeyPage []byte
+
+// PasskeyScript is GET /surface/passkey.js.
+//
+//go:embed passkey/passkey.js
+var PasskeyScript []byte
--- a/host/surfacepage/passkey/index.html
+++ b/host/surfacepage/passkey/index.html
@@ -0,0 +1,28 @@
+<!doctype html>
+<html lang="en">
+<head>
+<meta charset="utf-8">
+<meta name="viewport" content="width=device-width, initial-scale=1">
+<title>AILANG World passkey</title>
+<style>
+body { font: 16px/1.5 system-ui, sans-serif; max-width: 42rem; margin: 2rem auto; padding: 0 1rem; color: #1d1d1f; }
+h1 { font-size: 1.4rem; }
+#what, #status { white-space: pre-wrap; }
+#status[data-error] { color: #b00020; font-weight: 600; }
+#code { font: 600 1.3rem ui-monospace, monospace; letter-spacing: 0.05em; width: 100%; }
+#result { font: 12px ui-monospace, monospace; width: 100%; height: 9rem; }
+button, #save { font-size: 1rem; padding: 0.4rem 1rem; }
+[hidden] { display: none !important; }
+</style>
+</head>
+<body data-state="loading">
+<h1>AILANG World passkey</h1>
+<p id="what">Reading the link…</p>
+<p><button id="go" type="button" hidden disabled></button></p>
+<p id="status"></p>
+<p id="code-box" hidden><label for="code">Enrolment code (type or paste it at the terminal; a fingerprint of the file, not a secret)</label><br><input id="code" readonly></p>
+<p id="result-box" hidden><label for="result">The file to save</label><br><textarea id="result" readonly></textarea></p>
+<p><a id="save" hidden>Save</a></p>
+<script src="/surface/passkey.js"></script>
+</body>
+</html>
--- a/host/surfacepage/passkey/passkey.js
+++ b/host/surfacepage/passkey/passkey.js
@@ -0,0 +1,138 @@
+'use strict';
+// AILANG World's passkey ceremony page (row 201 A2, design §2.14). It reads the link's fragment,
+// runs one WebAuthn create or get in this browser, and offers the result as a file to save. It
+// makes no network request and holds no data: the CLI judges the file. Text is set with
+// textContent only.
+(() => {
+  const TIMEOUT_MS = 240000; // [r2] the application's own bound: WebAuthn's timeout is only a hint
+  const $ = (id) => document.getElementById(id);
+  const setState = (s) => { document.body.dataset.state = s; };
+  const say = (text, error) => {
+    $('status').textContent = text;
+    $('status').toggleAttribute('data-error', error);
+  };
+
+  const b64u = (buf) => {
+    let s = '';
+    for (const b of new Uint8Array(buf)) s += String.fromCharCode(b);
+    return btoa(s).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
+  };
+  const unb64u = (s) => {
+    if (!/^[A-Za-z0-9_-]*$/.test(s)) throw new Error('not base64url');
+    return Uint8Array.from(atob(s.replace(/-/g, '+').replace(/_/g, '/')), (c) => c.charCodeAt(0));
+  };
+
+  // The enrolment code (design §2.8): base32 of the first 16 bytes of SHA-256("world-enrolment-code/v1"
+  // ‖ 0x00 ‖ SHA-256(clientDataJSON) ‖ SHA-256(attestationObject)), shown in five groups.
+  const sha256 = async (bytes) => new Uint8Array(await crypto.subtle.digest('SHA-256', bytes));
+  async function enrolmentCode(clientDataJSON, attestationObject) {
+    const label = new TextEncoder().encode('world-enrolment-code/v1\0');
+    const msg = new Uint8Array(label.length + 64);
+    msg.set(label);
+    msg.set(await sha256(clientDataJSON), label.length);
+    msg.set(await sha256(attestationObject), label.length + 32);
+    const h = (await sha256(msg)).slice(0, 16);
+    const alphabet = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ234567';
+    let bits = 0, n = 0, out = '';
+    for (const b of h) {
+      bits = (bits << 8) | b;
+      n += 8;
+      while (n >= 5) { out += alphabet[(bits >>> (n - 5)) & 31]; n -= 5; }
+    }
+    out += alphabet[(bits << (5 - n)) & 31];
+    return [out.slice(0, 5), out.slice(5, 10), out.slice(10, 15), out.slice(15, 20), out.slice(20)].join('-');
+  }
+
+  // The §2.7 files, built here rather than by toJSON().
+  const registrationJSON = (c) => ({ id: c.id, rawId: b64u(c.rawId), type: c.type,
+    response: { clientDataJSON: b64u(c.response.clientDataJSON), attestationObject: b64u(c.response.attestationObject) } });
+  const assertionJSON = (c) => ({ id: c.id, rawId: b64u(c.rawId), type: c.type,
+    response: { clientDataJSON: b64u(c.response.clientDataJSON), authenticatorData: b64u(c.response.authenticatorData),
+      signature: b64u(c.response.signature), userHandle: c.response.userHandle ? b64u(c.response.userHandle) : null } });
+
+  function call(f, signal) {
+    const challenge = unb64u(f.challenge);
+    if (f.op === 'create') {
+      return navigator.credentials.create({ publicKey: { challenge, rp: { id: f.rp, name: 'AILANG World' },
+        user: { id: new TextEncoder().encode(f.name), name: f.name, displayName: f.name },
+        pubKeyCredParams: [{ type: 'public-key', alg: -7 }, { type: 'public-key', alg: -8 }],
+        authenticatorSelection: { userVerification: 'required', residentKey: 'preferred' },
+        attestation: 'none', timeout: TIMEOUT_MS }, signal });
+    }
+    return navigator.credentials.get({ publicKey: { challenge, rpId: f.rp,
+      allowCredentials: f.allow.map((id) => ({ type: 'public-key', id: unb64u(id) })),
+      userVerification: 'required', timeout: TIMEOUT_MS }, signal });
+  }
+
+  async function ceremony(f) {
+    const go = $('go');
+    go.disabled = true;
+    setState('pending');
+    say('Waiting for your passkey…', false);
+    const controller = new AbortController();
+    let expired = false;
+    const timer = setTimeout(() => {
+      expired = true;
+      controller.abort();
+      setState('expired');
+      say(`Timed out: no passkey answer within ${TIMEOUT_MS / 60000} minutes. Nothing can be saved from this attempt; try again.`, true);
+    }, TIMEOUT_MS);
+    let cred = null, failure = null;
+    try {
+      cred = await call(f, controller.signal);
+    } catch (e) {
+      failure = e;
+    } finally {
+      clearTimeout(timer);
+    }
+    if (expired) { // [r2] a result after the timer is ignored: the attempt is over, nothing is saved
+      go.disabled = false;
+      return;
+    }
+    if (failure !== null) {
+      setState('failed');
+      say(`The browser did not make a passkey response (${failure.name}: ${failure.message}). Nothing was saved; try again.`, true);
+      go.disabled = false;
+      return;
+    }
+    const create = f.op === 'create';
+    const json = JSON.stringify(create ? registrationJSON(cred) : assertionJSON(cred));
+    if (create) {
+      $('code').value = await enrolmentCode(cred.response.clientDataJSON, cred.response.attestationObject);
+      $('code-box').hidden = false;
+    }
+    $('result').value = json;
+    $('result-box').hidden = false;
+    const save = $('save');
+    save.href = URL.createObjectURL(new Blob([json + '\n'], { type: 'application/json' }));
+    save.download = create ? `world-registration-${f.name}.json` : `world-assertion-${f.request}.json`;
+    save.hidden = false;
+    setState('done');
+    say(create ? 'Passkey created. Save the file, then type the enrolment code at the terminal.'
+      : 'Signed by your passkey. Save the file, then run the approve command the terminal printed.', false);
+  }
+
+  let f;
+  try {
+    f = JSON.parse(new TextDecoder().decode(unb64u(location.hash.slice(1))));
+    if (f.v !== 1 || (f.op !== 'create' && f.op !== 'get')) throw new Error('not a version 1 create or get');
+  } catch (e) {
+    setState('refused');
+    $('what').textContent = `This is not a World passkey link (${e.message}). Open the whole link the terminal printed.`;
+    return;
+  }
+  if (f.origin !== location.origin || f.rp !== location.hostname) {
+    setState('refused');
+    $('what').textContent = `This link is for ${f.origin} (passkey site ${f.rp}): open it there, not at ${location.origin}.`;
+    return;
+  }
+  $('what').textContent = f.op === 'create'
+    ? `Create a passkey for approver ${f.name} at ${f.rp}`
+    : `Approve request ${f.request}, expires ${new Date(f.expires * 1000).toLocaleString()}, message hash ${f.hash16}`;
+  const go = $('go');
+  go.textContent = f.op === 'create' ? 'Create passkey' : 'Approve with passkey';
+  go.addEventListener('click', () => { ceremony(f); });
+  go.hidden = false;
+  go.disabled = false;
+  setState('ready');
+})();
````

### Patch M7c — the browser driver

````diff M7c
--- a/cmd/ailang-worldd/testdata/passkey_cdp.mjs
+++ b/cmd/ailang-worldd/testdata/passkey_cdp.mjs
@@ -0,0 +1,172 @@
+// The passkey page's browser driver (row 201 A2 M7; XB1, XB2 and host/approval/testdata/webauthn/
+// generate.sh). Test-only: Chrome headless with a temporary profile under argv[2], driven over the
+// DevTools protocol, with a CDP VIRTUAL authenticator, a test key (D-WORLD-109): never a real
+// profile, keychain, passkey or platform authenticator. Commands arrive as JSON lines on stdin;
+// each answer is one JSON line on stdout. EOF on stdin kills Chrome and removes the profile.
+import { spawn } from 'node:child_process';
+import fs from 'node:fs';
+import path from 'node:path';
+import readline from 'node:readline';
+
+const dir = process.argv[2];
+const profile = path.join(dir, 'profile');
+fs.mkdirSync(profile);
+const chrome = spawn(process.env.WORLD_CHROME, ['--headless=new', '--remote-debugging-port=0', '--user-data-dir=' + profile,
+  '--no-first-run', '--no-default-browser-check', '--disable-background-networking', 'about:blank'], { stdio: 'ignore' });
+const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
+const until = async (ok, ms, what) => {
+  for (const end = Date.now() + ms; Date.now() < end; await sleep(25)) if (await ok()) return;
+  throw new Error('timed out waiting for ' + what);
+};
+let port;
+await until(() => { // Chrome writes the file in two steps: wait for both lines
+  try { port = fs.readFileSync(path.join(profile, 'DevToolsActivePort'), 'utf8'); } catch { return false; }
+  return /^\d+\n\/devtools\/browser\/\S+/.test(port);
+}, 30000, 'DevToolsActivePort');
+const [dport, dpath] = port.trim().split('\n');
+const ws = new WebSocket(`ws://127.0.0.1:${dport}${dpath}`);
+await new Promise((r, j) => { ws.onopen = r; ws.onerror = j; });
+let seq = 0, session = null, target = null, authenticatorId = null;
+const replies = new Map(), events = [];
+ws.onmessage = (m) => {
+  const msg = JSON.parse(m.data);
+  if (msg.id && replies.has(msg.id)) { replies.get(msg.id)(msg); replies.delete(msg.id); } else if (msg.method) events.push(msg);
+};
+const send = async (method, params = {}, sessionId = session) => {
+  const msg = await new Promise((r) => { const id = ++seq; replies.set(id, r); ws.send(JSON.stringify({ id, method, params, sessionId: sessionId || undefined })); });
+  if (msg.error) throw new Error(`${method}: ${msg.error.message}`);
+  return msg.result;
+};
+const evaluate = async (expression) => {
+  const r = await send('Runtime.evaluate', { expression: `JSON.stringify(${expression})`, returnByValue: true, awaitPromise: true });
+  if (r.exceptionDetails) throw new Error(r.exceptionDetails.text);
+  return JSON.parse(r.result.value);
+};
+const click = async (id) => {
+  const box = await evaluate(`(() => { const e = document.getElementById('${id}'); e.scrollIntoView({ block: 'center' }); const r = e.getBoundingClientRect(); return { x: r.x + r.width / 2, y: r.y + r.height / 2, w: r.width }; })()`);
+  if (box.w === 0) return false;
+  for (const type of ['mousePressed', 'mouseReleased']) await send('Input.dispatchMouseEvent', { type, x: box.x, y: box.y, button: 'left', clickCount: 1 });
+  return true;
+};
+
+// The injected credentials operations of XB2 (pending until aborted, or a late result) and of
+// generate.sh (one algorithm). Each records what the page passed in window.__op.
+const injections = {
+  abort: `for (const m of ['create', 'get']) navigator.credentials[m] = (o) => new Promise((_, reject) => {
+    const op = window.__op = { method: m, signal: !!(o && o.signal), t0: performance.now(), d0: Date.now() };
+    if (op.signal) o.signal.addEventListener('abort', () => { op.abortedAfter = performance.now() - op.t0; op.abortedAfterDate = Date.now() - op.d0; op.aborted = o.signal.aborted; reject(o.signal.reason); });
+  });`,
+  late: `for (const m of ['create', 'get']) navigator.credentials[m] = (o) => new Promise((resolve) => {
+    const op = window.__op = { method: m, signal: !!(o && o.signal), t0: performance.now(), d0: Date.now() };
+    if (op.signal) o.signal.addEventListener('abort', () => { op.abortedAfter = performance.now() - op.t0; op.abortedAfterDate = Date.now() - op.d0; op.aborted = o.signal.aborted; });
+    const b = (n) => new Uint8Array(n).buffer;
+    setTimeout(() => { op.resolvedAfter = performance.now() - op.t0; op.resolvedAfterDate = Date.now() - op.d0; resolve({ id: 'AA', rawId: b(1), type: 'public-key',
+      response: { clientDataJSON: b(2), attestationObject: b(3), authenticatorData: b(37), signature: b(8), userHandle: null } }); }, 245000);
+  });`,
+  eddsa: `{ const create = navigator.credentials.create.bind(navigator.credentials);
+    navigator.credentials.create = (o) => { o.publicKey.pubKeyCredParams = o.publicKey.pubKeyCredParams.filter((p) => p.alg === -8); return create(o); }; }`,
+};
+const watch = `window.__cspv = []; document.addEventListener('securitypolicyviolation', (e) => window.__cspv.push(e.violatedDirective + ' ' + e.blockedURI));`;
+
+// A virtual authenticator that can verify the user (uv) or has no user verification at all.
+const addAuthenticator = async (uv) => {
+  ({ authenticatorId } = await send('WebAuthn.addVirtualAuthenticator', { options: { protocol: 'ctap2', ctap2Version: 'ctap2_1', transport: 'internal',
+    hasResidentKey: true, hasUserVerification: uv, isUserVerified: uv, automaticPresenceSimulation: true } }));
+};
+
+async function tab({ authenticator = null, inject = null }) {
+  if (target) await send('Target.closeTarget', { targetId: target }, null);
+  const { targetId } = await send('Target.createTarget', { url: 'about:blank' }, null);
+  const { sessionId } = await send('Target.attachToTarget', { targetId, flatten: true }, null);
+  session = sessionId;
+  target = targetId;
+  for (const d of ['Page', 'Runtime', 'Log']) await send(d + '.enable');
+  await send('Page.addScriptToEvaluateOnNewDocument', { source: watch + (inject ? injections[inject] : '') });
+  if (authenticator) {
+    await send('WebAuthn.enable', { enableUI: false });
+    await addAuthenticator(authenticator === 'uv');
+  }
+  return {};
+}
+
+const pageState = () => evaluate(`(() => { const $ = (id) => document.getElementById(id); return {
+  state: document.body.dataset.state, what: $('what').textContent, status: $('status').textContent, statusError: $('status').hasAttribute('data-error'),
+  button: $('go').textContent, buttonHidden: $('go').hidden, buttonDisabled: $('go').disabled, code: $('code').value, result: $('result').value,
+  saveHidden: $('save').hidden, saveHref: $('save').getAttribute('href') || '', saveName: $('save').download,
+  violations: window.__cspv || [], op: window.__op || null }; })()`);
+
+async function open({ url }) {
+  for (const u of ['about:blank', url]) { // a new document each time, even when only the fragment differs
+    const from = events.length;
+    await send('Page.navigate', { url: u });
+    await until(() => events.slice(from).some((e) => e.method === 'Page.loadEventFired'), 15000, 'the page load');
+  }
+  await until(async () => (await evaluate(`document.body && document.body.dataset.state`)) !== 'loading', 5000, 'the page script');
+  return pageState();
+}
+
+async function ceremony({ virtualBudget = 0, download = null, settleMs = 30000 }) {
+  const from = events.length;
+  if (download) await send('Browser.setDownloadBehavior', { behavior: 'allow', downloadPath: download, eventsEnabled: true }, null);
+  // In virtual time the page's clock is paused from before the click, so the page's timer and the
+  // operation start at one instant, however loaded the machine is.
+  if (virtualBudget > 0) await send('Emulation.setVirtualTimePolicy', { policy: 'pause' });
+  if (!(await click('go'))) throw new Error('the ceremony button is not visible');
+  const pending = await pageState();
+  const t0 = Date.now();
+  if (virtualBudget > 0) {
+    await send('Emulation.setVirtualTimePolicy', { policy: 'advance', budget: virtualBudget });
+    await until(() => events.slice(from).some((e) => e.method === 'Emulation.virtualTimeBudgetExpired'), settleMs, 'the virtual time budget');
+  } else {
+    await until(async () => (await evaluate(`document.body.dataset.state`)) !== 'pending', settleMs, 'the ceremony');
+  }
+  const settled = await pageState();
+  const realMs = Date.now() - t0;
+  let saveClicked = false;
+  if (download) {
+    saveClicked = await click('save');
+    const done = () => events.slice(from).some((e) => e.method === 'Browser.downloadProgress' && e.params.state !== 'inProgress');
+    if (saveClicked) await until(done, settleMs, 'the download');
+    else await sleep(500);
+  }
+  const downloads = events.slice(from).filter((e) => e.method.startsWith('Browser.download'))
+    .map((e) => ({ event: e.method, state: e.params.state || '', name: e.params.suggestedFilename || '', url: e.params.url || '' }));
+  const errors = events.slice(from).filter((e) => e.method === 'Runtime.exceptionThrown' || (e.method === 'Log.entryAdded' && e.params.entry.level === 'error'))
+    .map((e) => e.params.entry ? e.params.entry.text : e.params.exceptionDetails.text);
+  return { pending, settled, after: await pageState(), saveClicked, downloads, errors, realMs };
+}
+
+// uv moves the authenticator's credentials (keys and counters) to a new virtual authenticator that
+// can verify the user, or has no user verification at all.
+async function uv({ verified }) {
+  const { credentials } = await send('WebAuthn.getCredentials', { authenticatorId });
+  await send('WebAuthn.removeVirtualAuthenticator', { authenticatorId });
+  await addAuthenticator(verified);
+  for (const credential of credentials) await send('WebAuthn.addCredential', { authenticatorId, credential });
+  return { moved: credentials.length };
+}
+
+// advance runs budget ms of page time at once (virtual time) and returns the page.
+async function advance({ budget }) {
+  const from = events.length;
+  await send('Emulation.setVirtualTimePolicy', { policy: 'advance', budget });
+  await until(() => events.slice(from).some((e) => e.method === 'Emulation.virtualTimeBudgetExpired'), 30000, 'the virtual time budget');
+  return pageState();
+}
+
+const product = (await send('Browser.getVersion', {}, null)).product;
+console.log(JSON.stringify({ ready: true, chrome: product, pid: chrome.pid }));
+const commands = { tab, open, ceremony, uv, advance };
+for await (const line of readline.createInterface({ input: process.stdin })) {
+  const { cmd, ...args } = JSON.parse(line);
+  try {
+    console.log(JSON.stringify({ ok: await commands[cmd](args) }));
+  } catch (e) {
+    console.log(JSON.stringify({ error: String(e && e.message || e) }));
+  }
+}
+ws.close();
+chrome.kill('SIGKILL');
+await sleep(300);
+fs.rmSync(profile, { recursive: true, force: true });
+process.exit(0);
````

### Patch M7d — the Go harness

````diff M7d
--- a/cmd/ailang-worldd/passkey_page_driver_test.go
+++ b/cmd/ailang-worldd/passkey_page_driver_test.go
@@ -0,0 +1,230 @@
+package main
+
+// Row 201 A2 sprint 4 M7 (design_docs/planned/w-world-approvals-by-key-a2.md §2.14, §6 XB1, XB2):
+// the harness of the passkey page's browser tests. testdata/passkey_cdp.mjs drives Chrome headless
+// with a temporary profile and a CDP virtual authenticator, a test key (D-WORLD-109): never a real
+// profile, keychain, passkey or platform authenticator. The page is served by the production daemon
+// handler. Gated: WORLD_CHROME and node on PATH, else a counted skip, which
+// WORLD_REQUIRE_PASSKEY_BROWSER=1 turns into a failure.
+
+import (
+	"bufio"
+	"context"
+	"encoding/json"
+	"fmt"
+	"io"
+	"net"
+	"net/http/httptest"
+	"os"
+	"os/exec"
+	"path/filepath"
+	"regexp"
+	"strings"
+	"syscall"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/daemon"
+)
+
+// passkeyBrowserBudget bounds one driver (Chrome's start, every ceremony and the reap). Measured
+// (plan §1.2 C19): one XB1 or XB2 test takes 4-19 s on the planner's machine at load 20-80.
+const passkeyBrowserBudget = 120 * time.Second
+
+// pageView is what the page shows (the driver reads it after each step).
+type pageView struct {
+	State, What, Status, Button, Code, Result, SaveHref, SaveName string
+	StatusError, ButtonHidden, ButtonDisabled, SaveHidden         bool
+	Violations                                                    []string
+	Op                                                            *struct {
+		Method                                                           string
+		Signal, Aborted                                                  bool
+		AbortedAfter, AbortedAfterDate, ResolvedAfter, ResolvedAfterDate *float64
+	}
+}
+
+// pageCeremony is one click on the ceremony button: the page while pending, when settled, and
+// after the Save click; the download events; the page's errors.
+type pageCeremony struct {
+	Pending, Settled, After pageView
+	SaveClicked             bool
+	Downloads               []struct{ Event, State, Name, URL string }
+	Errors                  []string
+	RealMs                  int64
+}
+
+type passkeyBrowser struct {
+	in     io.WriteCloser
+	out    *bufio.Reader
+	chrome string
+}
+
+// startPasskeyBrowser starts the driver, or skips (counted) when Chrome or node is missing.
+func startPasskeyBrowser(t *testing.T) *passkeyBrowser {
+	t.Helper()
+	chrome := os.Getenv("WORLD_CHROME")
+	node, err := exec.LookPath("node")
+	if chrome == "" || err != nil {
+		msg := fmt.Sprintf("SKIP COUNTED: the passkey page browser test needs WORLD_CHROME (%q) and node on PATH (%v)", chrome, err)
+		if os.Getenv("WORLD_REQUIRE_PASSKEY_BROWSER") == "1" {
+			t.Fatal(msg)
+		}
+		t.Skip(msg)
+	}
+	ctx, cancel := context.WithTimeout(context.Background(), passkeyBrowserBudget)
+	cmd := exec.CommandContext(ctx, node, filepath.Join("testdata", "passkey_cdp.mjs"), t.TempDir())
+	var stderr syncBuffer // cli_test.go's: os/exec's copier writes until Wait returns (TestSubprocessSinkGate)
+	cmd.Stderr = &stderr
+	in, err := cmd.StdinPipe()
+	if err != nil {
+		t.Fatal(err)
+	}
+	out, err := cmd.StdoutPipe()
+	if err != nil {
+		t.Fatal(err)
+	}
+	if err := cmd.Start(); err != nil {
+		cancel()
+		t.Fatal(err)
+	}
+	b := &passkeyBrowser{in: in, out: bufio.NewReader(out)}
+	var ready struct {
+		Ready  bool
+		Chrome string
+		Pid    int
+	}
+	t.Cleanup(func() {
+		_ = in.Close()
+		werr := cmd.Wait()
+		cancel()
+		if werr != nil && ready.Pid > 0 { // the driver did not reap Chrome itself
+			_ = syscall.Kill(ready.Pid, syscall.SIGKILL)
+		}
+		if t.Failed() {
+			t.Logf("driver: %v\n%s", werr, stderr.String())
+		}
+	})
+	if err := b.read(&ready); err != nil || !ready.Ready {
+		t.Fatalf("the driver did not start: %v", err)
+	}
+	b.chrome = ready.Chrome
+	return b
+}
+
+func (b *passkeyBrowser) read(v any) error {
+	line, err := b.out.ReadBytes('\n')
+	if err != nil {
+		return err
+	}
+	return json.Unmarshal(line, v)
+}
+
+// call sends one driver command and decodes its answer into out. It never calls t.Fatal, so a
+// terminal's read goroutine may use it.
+func (b *passkeyBrowser) call(out any, cmd string, args map[string]any) error {
+	if args == nil {
+		args = map[string]any{}
+	}
+	args["cmd"] = cmd
+	line, _ := json.Marshal(args)
+	if _, err := b.in.Write(append(line, '\n')); err != nil {
+		return err
+	}
+	var reply struct {
+		Ok    json.RawMessage
+		Error string
+	}
+	if err := b.read(&reply); err != nil {
+		return err
+	}
+	if reply.Error != "" {
+		return fmt.Errorf("driver %s: %s", cmd, reply.Error)
+	}
+	if out == nil {
+		return nil
+	}
+	return json.Unmarshal(reply.Ok, out)
+}
+
+func (b *passkeyBrowser) do(t *testing.T, out any, cmd string, args map[string]any) {
+	t.Helper()
+	if err := b.call(out, cmd, args); err != nil {
+		t.Fatal(err)
+	}
+}
+
+// passkeyPageOrigin serves the production daemon handler on localhost and returns its origin.
+func passkeyPageOrigin(t *testing.T) string {
+	t.Helper()
+	d, err := daemon.New(boundedTestContext(t), daemon.Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: daemon.DefaultBindHost, ErrorLog: io.Discard})
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { _ = d.Close() })
+	srv := httptest.NewServer(d.Handler())
+	t.Cleanup(srv.Close)
+	_, port, _ := net.SplitHostPort(srv.Listener.Addr().String())
+	return "http://localhost:" + port
+}
+
+var printedLinkRE = regexp.MustCompile(`(?m)^  (http\S+/surface/passkey#\S+)$`)
+
+// printedLink is the one ceremony link a CLI run printed.
+func printedLink(out string) (string, error) {
+	m := printedLinkRE.FindAllStringSubmatch(out, -1)
+	if len(m) != 1 {
+		return "", fmt.Errorf("%d ceremony links in %q, want 1", len(m), out)
+	}
+	return m[0][1], nil
+}
+
+// saved checks that c saved exactly the page's JSON plus a newline as dir/name through a blob:
+// download of origin, with no CSP violation or page error, and returns the file.
+func saved(t *testing.T, c pageCeremony, origin, dir, name string) []byte {
+	t.Helper()
+	if c.After.State != "done" || c.After.StatusError || !c.SaveClicked || c.After.SaveName != name {
+		t.Fatalf("the page did not finish: state %q, status %q, save %q clicked %v", c.After.State, c.After.Status, c.After.SaveName, c.SaveClicked)
+	}
+	var began, completed bool
+	for _, d := range c.Downloads {
+		began = began || (d.Event == "Browser.downloadWillBegin" && d.Name == name && strings.HasPrefix(d.URL, "blob:"+origin+"/"))
+		completed = completed || (d.Event == "Browser.downloadProgress" && d.State == "completed")
+	}
+	if !began || !completed {
+		t.Fatalf("no completed blob: download of %s: %+v", name, c.Downloads)
+	}
+	if len(c.After.Violations) != 0 || len(c.Errors) != 0 {
+		t.Fatalf("CSP violations %q, page errors %q", c.After.Violations, c.Errors)
+	}
+	file, err := os.ReadFile(filepath.Join(dir, name))
+	if err != nil {
+		t.Fatal(err)
+	}
+	if string(file) != c.After.Result+"\n" {
+		t.Fatalf("the saved file is not the page's JSON plus a newline:\n%s\n%s", file, c.After.Result)
+	}
+	return file
+}
+
+// opString shows the injected operation's record.
+func opString(o any) string { b, _ := json.Marshal(o); return string(b) }
+
+// worldFiles is the world-*.json files in dir: what the page's Save link can write. Anything else
+// Chrome leaves there is logged, not counted (plan §1.2 C21: a stray "downloads.html" download,
+// outside every ceremony's events, seen under load).
+func worldFiles(t *testing.T, dir string) []string {
+	t.Helper()
+	entries, err := os.ReadDir(dir)
+	if err != nil {
+		t.Fatal(err)
+	}
+	var names []string
+	for _, e := range entries {
+		if strings.HasPrefix(e.Name(), "world-") && strings.HasSuffix(e.Name(), ".json") {
+			names = append(names, e.Name())
+		} else {
+			t.Logf("not a page save, not counted: %s", e.Name())
+		}
+	}
+	return names
+}
````

### Patch M7e — the routes, PG1, XB1 and XB2

````diff M7e
--- a/cmd/ailang-worldd/passkey_page_abort_test.go
+++ b/cmd/ailang-worldd/passkey_page_abort_test.go
@@ -0,0 +1,63 @@
+package main
+
+// Row 201 A2 sprint 4 M7, XB2 (design §2.14 [r2], §6): the page's own 240000 ms timer, in virtual
+// time, against an injected credentials operation.
+
+import (
+	"encoding/base64"
+	"math"
+	"testing"
+)
+
+// XB2 [r2]: a ceremony the authenticator never answers is aborted by the page's own 240000 ms
+// timer, driven in virtual time (paused from before the click, so the timer and the operation start at
+// one instant): the injected operation saw the signal abort at exactly 240000 ms of page time
+// (Date.now(); performance.now() agrees within one 0.1 ms quantum of its coarsening), the page shows the timeout,
+// nothing can be saved, a result arriving after the timer (at 245000 ms) is ignored, and the button
+// was disabled while the operation was pending.
+func TestPasskeyPageAbortsAPendingCeremony(t *testing.T) {
+	b := startPasskeyBrowser(t)
+	origin := passkeyPageOrigin(t)
+	b64 := base64.RawURLEncoding.EncodeToString
+	links := map[string]string{
+		"create": ceremonyLink{V: 1, Op: "create", RP: "localhost", Origin: origin, Name: "mark-mac", Challenge: b64(make([]byte, 32))}.link(),
+		"get": ceremonyLink{V: 1, Op: "get", RP: "localhost", Origin: origin, Request: "mr-0123456789abcdef", Expires: approveT0 + 600,
+			Hash16: "0123456789abcdef", Challenge: b64(make([]byte, 32)), Allow: []string{b64(make([]byte, 32))}}.link(),
+	}
+	for _, op := range []string{"create", "get"} {
+		for _, inject := range []string{"abort", "late"} {
+			t.Run(op+"-"+inject, func(t *testing.T) {
+				b.do(t, nil, "tab", map[string]any{"inject": inject})
+				var view pageView
+				b.do(t, &view, "open", map[string]any{"url": links[op]})
+				if view.State != "ready" {
+					t.Fatalf("the page is %q before the click", view.State)
+				}
+				dl := t.TempDir()
+				var c pageCeremony
+				b.do(t, &c, "ceremony", map[string]any{"virtualBudget": 250000, "download": dl})
+				t.Logf("%s-%s: 250 s of page time in %d ms", op, inject, c.RealMs)
+				if c.Pending.State != "pending" || !c.Pending.ButtonDisabled {
+					t.Errorf("while pending: state %q, button disabled %v", c.Pending.State, c.Pending.ButtonDisabled)
+				}
+				o := c.After.Op
+				if o == nil || o.Method != op || !o.Signal || !o.Aborted || o.AbortedAfterDate == nil || *o.AbortedAfterDate != 240000 || o.AbortedAfter == nil || math.Abs(*o.AbortedAfter-240000) > 0.15 {
+					t.Fatalf("the operation: %s; want %s with a signal aborted at exactly 240000 ms", opString(o), op)
+				}
+				if inject == "late" && (o.ResolvedAfterDate == nil || *o.ResolvedAfterDate != 245000) {
+					t.Fatalf("the late result was not delivered at 245000 ms: %s", opString(o))
+				}
+				a := c.After
+				if a.State != "expired" || !a.StatusError || a.Status != "Timed out: no passkey answer within 4 minutes. Nothing can be saved from this attempt; try again." {
+					t.Errorf("after the timer: state %q, error %v, status %q", a.State, a.StatusError, a.Status)
+				}
+				if files := worldFiles(t, dl); !a.SaveHidden || a.SaveHref != "" || a.Result != "" || a.Code != "" || c.SaveClicked || len(c.Downloads) != 0 || len(files) != 0 {
+					t.Errorf("something can be saved after expiry: %+v, downloads %+v, files %q", a, c.Downloads, files)
+				}
+				if a.ButtonDisabled || len(a.Violations) != 0 || len(c.Errors) != 0 {
+					t.Errorf("button disabled %v after the operation settled; CSP violations %q; page errors %q", a.ButtonDisabled, a.Violations, c.Errors)
+				}
+			})
+		}
+	}
+}
--- a/cmd/ailang-worldd/passkey_page_chrome_test.go
+++ b/cmd/ailang-worldd/passkey_page_chrome_test.go
@@ -0,0 +1,171 @@
+package main
+
+// Row 201 A2 sprint 4 M7, XB1 (design §2.14, §6): the page in a live Chrome, end to end through the
+// CLI. The links it opens are the ones the CLI printed; the files it saves go back through the CLI.
+
+import (
+	"bytes"
+	"encoding/base64"
+	"encoding/hex"
+	"os"
+	"path/filepath"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+// XB1: the CLI's create link enrols, through approver enroll --format webauthn, the file the page
+// saved, with the code the page showed; the CLI's get link approves, through approve
+// --webauthn-assertion, the file the page saved, after a get on an authenticator without user
+// verification saved nothing; the finished attempt's timer never fires. A create on an authenticator
+// without user verification saves nothing, and a link for another page is refused before any button.
+// (Whether the button is disabled while pending is XB2's: a virtual authenticator answers too fast to
+// observe it.)
+func TestPasskeyPageInChrome(t *testing.T) {
+	b := startPasskeyBrowser(t)
+	origin := passkeyPageOrigin(t)
+	t.Logf("%s, page at %s", b.chrome, origin)
+
+	t.Run("enrol-then-approve", func(t *testing.T) {
+		b.do(t, nil, "tab", map[string]any{"authenticator": "uv"})
+		a := newApproveRig(t, "ep1", approveT0)
+		dl := t.TempDir()
+		reg := filepath.Join(dl, "world-registration-mark-mac.json")
+		var view pageView
+		var page pageCeremony
+		var out, errw bytes.Buffer
+		term := &ceremonyTerm{}
+		term.answer = func(n int) (string, error) {
+			if n > 1 {
+				return "", os.ErrClosed
+			}
+			if n == 1 {
+				return "y\n", nil
+			}
+			link, err := printedLink(out.String())
+			if err == nil {
+				err = b.call(&view, "open", map[string]any{"url": link})
+			}
+			if err == nil {
+				err = b.call(&page, "ceremony", map[string]any{"download": dl})
+			}
+			if err != nil {
+				t.Error(err)
+				return "", err
+			}
+			return page.After.Code + "\n", nil
+		}
+		args := []string{"enroll", "--name", "mark-mac", "--format", "webauthn", "--origin", origin, "--registration", reg, "--db", a.r.f.db}
+		if code := runApprover(args, &out, &errw, enrolEnv(term, nil)); code != exitOK {
+			t.Fatalf("enrol: rc %d\nstdout %s\nstderr %s", code, out.String(), errw.String())
+		}
+		if view.State != "ready" || view.What != "Create a passkey for approver mark-mac at localhost" || view.Button != "Create passkey" || view.ButtonDisabled {
+			t.Fatalf("the create page shows %+v", view)
+		}
+		file := saved(t, page, origin, dl, "world-registration-mark-mac.json")
+		if page.After.Code != codeOf(t, file) {
+			t.Fatalf("the page's code %q is not Go's %q for the saved file", page.After.Code, codeOf(t, file))
+		}
+		row, pk := enrolled(t, a.credPath)
+		if row.Name != "mark-mac" || pk.Origin != origin || pk.RPID != "localhost" || pk.Alg != approval.AlgES256 || !strings.Contains(out.String(), "enrolled approver mark-mac: passkey ") {
+			t.Fatalf("enrolled %+v %+v\nstdout %s", row, pk, out.String())
+		}
+
+		var lout, lerr bytes.Buffer
+		var opened bool
+		if code := runApprove(a.args("--webauthn-link"), &lout, &lerr, approveEnv(nil, &opened, atTime(approveT0+1))); code != exitOK || opened {
+			t.Fatalf("--webauthn-link: rc %d, terminal %v\nstderr %s", code, opened, lerr.String())
+		}
+		link, err := printedLink(lout.String())
+		if err != nil {
+			t.Fatal(err)
+		}
+		b.do(t, &view, "open", map[string]any{"url": link})
+		h := requestHashOf(t, a)
+		if view.State != "ready" || !strings.HasPrefix(view.What, "Approve request "+a.id+", expires ") || !strings.HasSuffix(view.What, ", message hash "+hex.EncodeToString(h[:8])) || view.Button != "Approve with passkey" {
+			t.Fatalf("the get page shows %+v", view)
+		}
+		var moved struct{ Moved int }
+		b.do(t, &moved, "uv", map[string]any{"verified": false}) // an authenticator with no user verification: nothing to save
+		b.do(t, &page, "ceremony", map[string]any{"download": dl})
+		if files := worldFiles(t, dl); page.After.State != "failed" || !page.After.StatusError || !page.After.SaveHidden || page.SaveClicked || len(page.Downloads) != 0 || len(files) != 1 || page.After.ButtonDisabled {
+			t.Fatalf("a get without user verification: %+v, downloads %+v, files %q", page.After, page.Downloads, files)
+		}
+		if moved.Moved != 1 {
+			t.Fatalf("%d credentials moved, want 1", moved.Moved)
+		}
+		b.do(t, nil, "uv", map[string]any{"verified": true}) // the retry, on the same page
+		b.do(t, &page, "ceremony", map[string]any{"download": dl})
+		name := "world-assertion-" + a.id + ".json"
+		saved(t, page, origin, dl, name)
+		b.do(t, &view, "advance", map[string]any{"budget": 250000}) // the finished attempt's timer was cleared
+		if view.State != "done" || view.StatusError || view.SaveHidden {
+			t.Fatalf("240 s after the result the page shows %+v", view)
+		}
+		if code := runApprove(a.args("--webauthn-assertion", filepath.Join(dl, name)), &lout, &lerr, approveEnv(nil, &opened, atTime(approveT0+2))); code != exitOK || opened {
+			t.Fatalf("--webauthn-assertion: rc %d, terminal %v\nstderr %s", code, opened, lerr.String())
+		}
+		if r := a.row(t); !r.Consumed || r.ApprovedBy != "webauthn:mark-mac" {
+			t.Fatalf("request consumed %v by %q, want webauthn:mark-mac", r.Consumed, r.ApprovedBy)
+		}
+		if after := passkeyOf(t, a.credPath, "mark-mac"); after.SignCount <= pk.SignCount {
+			t.Fatalf("sign count %d after the approval, %d at enrolment", after.SignCount, pk.SignCount)
+		}
+	})
+
+	t.Run("user-verification-required", func(t *testing.T) {
+		b.do(t, nil, "tab", map[string]any{"authenticator": "no-uv"})
+		db, creds := approverStore(t)
+		dl := t.TempDir()
+		var page pageCeremony
+		var out, errw bytes.Buffer
+		term := &ceremonyTerm{}
+		term.answer = func(n int) (string, error) {
+			if n > 0 {
+				return "", os.ErrClosed
+			}
+			link, err := printedLink(out.String())
+			if err == nil {
+				err = b.call(nil, "open", map[string]any{"url": link})
+			}
+			if err == nil {
+				err = b.call(&page, "ceremony", map[string]any{"download": dl})
+			}
+			if err != nil {
+				t.Error(err)
+			}
+			return page.After.Code + "\n", err
+		}
+		args := []string{"enroll", "--name", "mark-mac", "--format", "webauthn", "--origin", origin, "--registration", filepath.Join(dl, "world-registration-mark-mac.json"), "--db", db}
+		code := runApprover(args, &out, &errw, enrolEnv(term, nil))
+		if page.After.State != "failed" || !page.After.StatusError || !page.After.SaveHidden || page.SaveClicked || len(page.Downloads) != 0 || page.After.ButtonDisabled {
+			t.Fatalf("a UV-false authenticator: %+v, downloads %+v", page.After, page.Downloads)
+		}
+		if files := worldFiles(t, dl); code != exitUsage || len(files) != 0 || len(approverRows(t, creds)) != 0 {
+			t.Fatalf("rc %d, files saved %q, %d approvers; stderr %s", code, files, len(approverRows(t, creds)), errw.String())
+		}
+	})
+
+	t.Run("links-for-another-page", func(t *testing.T) {
+		b.do(t, nil, "tab", nil)
+		ok := ceremonyLink{V: 1, Op: "create", RP: "localhost", Origin: origin, Name: "mark-mac", Challenge: base64.RawURLEncoding.EncodeToString(make([]byte, 32))}
+		frag := func(l ceremonyLink) string { _, f, _ := strings.Cut(l.link(), "#"); return f }
+		edit := func(e func(*ceremonyLink)) string { l := ok; e(&l); return frag(l) }
+		for _, tc := range []struct{ name, fragment, state, what string }{
+			{"control", frag(ok), "ready", "Create a passkey for approver mark-mac at localhost"},
+			{"version", edit(func(l *ceremonyLink) { l.V = 2 }), "refused", "(not a version 1 create or get)"},
+			{"op", edit(func(l *ceremonyLink) { l.Op = "sign" }), "refused", "(not a version 1 create or get)"},
+			{"origin", edit(func(l *ceremonyLink) { l.Origin = "http://localhost:7644" }), "refused", "This link is for http://localhost:7644 (passkey site localhost): open it there, not at " + origin + "."},
+			{"rp", edit(func(l *ceremonyLink) { l.RP = "x.ts.net" }), "refused", "This link is for " + origin + " (passkey site x.ts.net)"},
+			{"not-base64url", "e30=", "refused", "(not base64url)"},
+			{"not-json", base64.RawURLEncoding.EncodeToString([]byte("{")), "refused", "This is not a World passkey link ("},
+		} {
+			var view pageView
+			b.do(t, &view, "open", map[string]any{"url": origin + "/surface/passkey#" + tc.fragment})
+			if view.State != tc.state || !strings.Contains(view.What, tc.what) || view.ButtonHidden != (tc.state == "refused") {
+				t.Errorf("%s: state %q, button hidden %v, what %q; want %q containing %q", tc.name, view.State, view.ButtonHidden, view.What, tc.state, tc.what)
+			}
+		}
+	})
+}
--- a/host/daemon/daemon.go
+++ b/host/daemon/daemon.go
@@ -1083,7 +1083,8 @@
 // pages, row 188 M2/M3) and GET /surface/surface.js (their script, M4a) are also unversioned
 // read routes outside the frozen table, as are GET /surface/follow-status.json, GET /surface/model.json and
 // GET /surface/stack.json (the app's view model, N1) and GET /surface with /surface/app.js,
-// /surface/app.css and /surface/ailang-logo.svg (the app, N2b).
+// /surface/app.css and /surface/ailang-logo.svg (the app, N2b), and GET /surface/passkey with
+// /surface/passkey.js (the passkey ceremony page, row 201 A2: static, it reads nothing).
 // The registry pattern deliberately uses a multi-segment wildcard: registry semantic IDs
 // such as "world/epoch-registry/v1" contain slashes.
 func (d *Daemon) Handler() http.Handler {
@@ -1111,6 +1112,8 @@
 	mux.HandleFunc("GET /surface/app.js", d.handleSurfaceAppScript)
 	mux.HandleFunc("GET /surface/app.css", d.handleSurfaceAppStyle)
 	mux.HandleFunc("GET /surface/ailang-logo.svg", d.handleSurfaceLogo)
+	mux.HandleFunc("GET /surface/passkey", d.handleSurfacePasskey)
+	mux.HandleFunc("GET /surface/passkey.js", d.handleSurfacePasskeyScript)
 	mux.HandleFunc("GET /surface/fonts/{name...}", d.handleSurfaceFont)
 	mux.HandleFunc("GET /surface/model.json", d.handleSurfaceModel)
 	mux.HandleFunc("GET /surface/stack.json", d.handleSurfaceStack)
--- a/host/daemon/passkey_page_test.go
+++ b/host/daemon/passkey_page_test.go
@@ -0,0 +1,68 @@
+package daemon
+
+import (
+	"bytes"
+	"net/http/httptest"
+	"os"
+	"path/filepath"
+	"regexp"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// PG1 (row 201 A2 M7, design §2.14, §6): GET /surface/passkey and /surface/passkey.js serve the
+// embedded page and script, byte-equal to host/surfacepage/passkey/, under surfaceCSP and the
+// surface headers, with the world store closed (they read nothing); POST is refused.
+func TestPasskeyPageRoutes(t *testing.T) {
+	d, err := New(boundedTestContext(t), Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: DefaultBindHost})
+	if err != nil {
+		t.Fatal(err)
+	}
+	h := d.Handler()
+	if err := d.Close(); err != nil { // no store read: the routes answer with the store closed
+		t.Fatal(err)
+	}
+	for _, tc := range []struct{ target, file, contentType string }{
+		{"/surface/passkey", "index.html", "text/html; charset=utf-8"},
+		{"/surface/passkey.js", "passkey.js", "text/javascript; charset=utf-8"},
+	} {
+		want, err := os.ReadFile(filepath.Join("..", "surfacepage", "passkey", tc.file))
+		if err != nil {
+			t.Fatal(err)
+		}
+		w := httptest.NewRecorder()
+		h.ServeHTTP(w, httptest.NewRequest("GET", tc.target, nil))
+		if w.Code != 200 || len(want) == 0 || !bytes.Equal(w.Body.Bytes(), want) {
+			t.Errorf("GET %s: %d, or not the bytes of host/surfacepage/passkey/%s", tc.target, w.Code, tc.file)
+		}
+		for k, v := range map[string]string{"Content-Security-Policy": surfaceCSP, "Content-Type": tc.contentType,
+			"Cache-Control": "no-store", "X-Content-Type-Options": "nosniff", "Referrer-Policy": "no-referrer"} {
+			if got := w.Header().Get(k); got != v {
+				t.Errorf("GET %s: %s = %q, want %q", tc.target, k, got, v)
+			}
+		}
+		w = httptest.NewRecorder()
+		h.ServeHTTP(w, httptest.NewRequest("POST", tc.target, nil))
+		if w.Code != 405 {
+			t.Errorf("POST %s = %d, want 405", tc.target, w.Code)
+		}
+	}
+
+	// The page loads exactly its own script, with no inline script or handler and no other URL.
+	page := string(surfacepage.PasskeyPage)
+	if scripts := regexp.MustCompile(`(?i)<script\b[^>]*>`).FindAllString(page, -1); len(scripts) != 1 || scripts[0] != `<script src="/surface/passkey.js">` || !strings.Contains(page, `<script src="/surface/passkey.js"></script>`) {
+		t.Errorf("scripts %q", scripts)
+	}
+	if m := regexp.MustCompile(`(?i)\son[a-z]+\s*=|<form|<img|<link|<iframe|https?://`).FindString(page); m != "" {
+		t.Errorf("the page has %q", m)
+	}
+	// Labelled instrument-health controls (not load-bearing; XB1 and XB2 are): the script makes no
+	// request and writes no markup.
+	for _, banned := range []string{"fetch(", "XMLHttpRequest", "innerHTML", "outerHTML", "insertAdjacentHTML", "document.write", "eval(", "new Function", "sendBeacon", "WebSocket"} {
+		if bytes.Contains(surfacepage.PasskeyScript, []byte(banned)) {
+			t.Errorf("passkey.js uses %s", banned)
+		}
+	}
+}
--- a/host/daemon/surface_app_test.go
+++ b/host/daemon/surface_app_test.go
@@ -15,7 +15,8 @@
 
 // TestSurfaceAppRoutes: row 188 N2b. GET /surface is the app shell (the
 // same bytes whatever the index holds), its assets are the exact embedded
-// bytes, and every other name under /surface/ is a 404.
+// bytes, and every other name under /surface/ is a 404 (row 201 A2 adds
+// /surface/passkey and /surface/passkey.js, PG1's, and no name under them).
 func TestSurfaceAppRoutes(t *testing.T) {
 	path := filepath.Join(t.TempDir(), "operator-index.json")
 	d, err := New(boundedTestContext(t), Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: DefaultBindHost, SurfaceIndex: path})
@@ -148,7 +149,8 @@
 
 	t.Run("404", func(t *testing.T) {
 		for _, target := range []string{"/surface/", "/surface/nope", "/surface/app.jsx", "/surface/index.html", "/surface/data.js",
-			"/surface/spike-data.json", "/surface/app/app.js", "/surface/static/x", "/surface/fonts.css"} {
+			"/surface/spike-data.json", "/surface/app/app.js", "/surface/static/x", "/surface/fonts.css",
+			"/surface/passkey/", "/surface/passkey/index.html", "/surface/passkey/passkey.js", "/surface/passkey.html"} {
 			if w := request("GET", target); w.Code != 404 {
 				t.Errorf("%s = %d, want 404", target, w.Code)
 			}
--- a/host/daemon/surface_page.go
+++ b/host/daemon/surface_page.go
@@ -170,6 +170,20 @@
 	writeSurfaceAsset(w, "image/svg+xml", surfacepage.AppLogo)
 }
 
+// handleSurfacePasskey is the passkey ceremony page (row 201 A2, design §2.14): embedded bytes
+// under surfaceCSP. It reads no store and no request; the page reads its own link's fragment.
+func (d *Daemon) handleSurfacePasskey(w http.ResponseWriter, _ *http.Request) {
+	setSurfaceHeaders(w)
+	_, _ = w.Write(surfacepage.PasskeyPage)
+}
+
+// handleSurfacePasskeyScript is the page's script, under the same headers.
+func (d *Daemon) handleSurfacePasskeyScript(w http.ResponseWriter, _ *http.Request) {
+	setSurfaceHeaders(w)
+	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
+	_, _ = w.Write(surfacepage.PasskeyScript)
+}
+
 func writeSurfaceAsset(w http.ResponseWriter, contentType string, body []byte) {
 	w.Header().Set("Content-Type", contentType)
 	w.Header().Set("X-Content-Type-Options", "nosniff")
````

### Patch M7f — generate.sh made real

````diff M7f
--- a/cmd/ailang-worldd/passkey_fixture_record_test.go
+++ b/cmd/ailang-worldd/passkey_fixture_record_test.go
@@ -0,0 +1,98 @@
+package main
+
+// Row 201 A2 sprint 4 M7 (design §2.15): host/approval/testdata/webauthn/generate.sh, its recorder
+// and its refusals.
+
+import (
+	"bytes"
+	"encoding/base64"
+	"encoding/hex"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"os"
+	"os/exec"
+	"path/filepath"
+	"testing"
+)
+
+// TestPasskeyPageRecordsTheFixtures is host/approval/testdata/webauthn/generate.sh's recorder:
+// with WORLD_PASSKEY_RECORD_DIR set it records, through the page and a virtual authenticator, an
+// ES256 and an EdDSA create and two gets over the fixture message hash (design §2.15), and writes
+// chrome-es256.json and chrome-eddsa.json there in the committed fixtures' layout. Unset, it skips.
+func TestPasskeyPageRecordsTheFixtures(t *testing.T) {
+	dir := os.Getenv("WORLD_PASSKEY_RECORD_DIR")
+	if dir == "" {
+		t.Skip("not recording: host/approval/testdata/webauthn/generate.sh sets WORLD_PASSKEY_RECORD_DIR")
+	}
+	b := startPasskeyBrowser(t)
+	origin := passkeyPageOrigin(t)
+	const hashHex = "29181d224d982f659d1bd4703238d3f8489e26968ffe2a22bd786fa3bf15fef5" // host/approval request_test.go
+	h, _ := hex.DecodeString(hashHex)
+	challenge := base64.RawURLEncoding.EncodeToString(h)
+	record := func(link, dl, name string) json.RawMessage {
+		var c pageCeremony
+		b.do(t, nil, "open", map[string]any{"url": link})
+		b.do(t, &c, "ceremony", map[string]any{"download": dl})
+		return json.RawMessage(bytes.TrimSuffix(saved(t, c, origin, dl, name), []byte("\n")))
+	}
+	for _, alg := range []struct{ name, inject string }{{"es256", ""}, {"eddsa", "eddsa"}} {
+		b.do(t, nil, "tab", map[string]any{"authenticator": "uv", "inject": alg.inject})
+		reg := record(ceremonyLink{V: 1, Op: "create", RP: "localhost", Origin: origin, Name: "mark-mac", Challenge: challenge}.link(), t.TempDir(), "world-registration-mark-mac.json")
+		var r struct{ RawID string }
+		if err := json.Unmarshal(reg, &r); err != nil {
+			t.Fatal(err)
+		}
+		get := ceremonyLink{V: 1, Op: "get", RP: "localhost", Origin: origin, Request: "mr-fixture", Expires: 1767225600, Hash16: hashHex[:16], Challenge: challenge, Allow: []string{r.RawID}}.link()
+		asr := record(get, t.TempDir(), "world-assertion-mr-fixture.json")
+		var asr2 struct {
+			Response struct{ AuthenticatorData string }
+		}
+		if err := json.Unmarshal(record(get, t.TempDir(), "world-assertion-mr-fixture.json"), &asr2); err != nil {
+			t.Fatal(err)
+		}
+		a2, _ := json.Marshal(map[string]string{"authenticatorData": asr2.Response.AuthenticatorData})
+		out := fmt.Sprintf("{\n  \"source\": %q,\n  \"chrome\": %q,\n  \"origin\": %q,\n  \"rp_id\": \"localhost\",\n  \"challenge_hex\": %q,\n  \"registration\": %s,\n  \"assertion\": %s,\n  \"assertion2\": %s\n}\n",
+			"row 201 A2, generate.sh: Chrome headless, the passkey page, CDP virtual authenticator (a test key, D-WORLD-109); see SOURCE.txt", b.chrome, origin, hashHex, reg, asr, a2)
+		if err := os.WriteFile(filepath.Join(dir, "chrome-"+alg.name+".json"), []byte(out), 0o644); err != nil {
+			t.Fatal(err)
+		}
+	}
+}
+
+// generate.sh refuses, writing nothing and running nothing, a call that does not name exactly one
+// existing empty directory, and one without WORLD_CHROME. It needs no browser.
+func TestPasskeyFixtureGeneratorRefuses(t *testing.T) {
+	script, err := filepath.Abs(filepath.Join("..", "..", "host", "approval", "testdata", "webauthn", "generate.sh"))
+	if err != nil {
+		t.Fatal(err)
+	}
+	full := t.TempDir()
+	if err := os.WriteFile(filepath.Join(full, "keep"), nil, 0o600); err != nil {
+		t.Fatal(err)
+	}
+	empty := t.TempDir()
+	for _, tc := range []struct {
+		name, chrome, want string
+		args               []string
+	}{
+		{"no-directory", "chrome", "usage: sh generate.sh <an empty directory>\n", nil},
+		{"two-directories", "chrome", "usage: sh generate.sh <an empty directory>\n", []string{empty, empty}},
+		{"missing", "chrome", "generate.sh: " + filepath.Join(empty, "nope") + " is not an empty directory\n", []string{filepath.Join(empty, "nope")}},
+		{"not-empty", "chrome", "generate.sh: " + full + " is not an empty directory\n", []string{full}},
+		{"no-chrome", "", "generate.sh: set WORLD_CHROME to Chrome's binary (only its virtual authenticator is used)\n", []string{empty}},
+	} {
+		cmd := exec.CommandContext(boundedTestContext(t), "sh", append([]string{script}, tc.args...)...)
+		cmd.Env = append(os.Environ(), "WORLD_CHROME="+tc.chrome, "PATH=/usr/bin:/bin") // a fall-through finds no go and no node there
+		var stdout, stderr bytes.Buffer
+		cmd.Stdout, cmd.Stderr = &stdout, &stderr
+		err := cmd.Run()
+		var exit *exec.ExitError
+		if !errors.As(err, &exit) || exit.ExitCode() != 2 || stderr.String() != tc.want || stdout.Len() != 0 {
+			t.Errorf("%s: %v, stderr %q, stdout %q; want exit 2 and %q", tc.name, err, stderr.String(), stdout.String(), tc.want)
+		}
+		if entries, _ := os.ReadDir(empty); len(entries) != 0 {
+			t.Fatalf("%s wrote into %s", tc.name, empty)
+		}
+	}
+}
--- a/host/approval/testdata/webauthn/SOURCE.txt
+++ b/host/approval/testdata/webauthn/SOURCE.txt
@@ -17,4 +17,5 @@
   chrome-es256.json  from the probe's es256.json, sha256 c0370982dbca416a7677b4fa1ff2cd64e9e533cefc608d3dcab42a1000d6caf6
   chrome-eddsa.json  from the probe's eddsa.json, sha256 1544721fc45dd65a553a5baab8a43fcd4df0375b0b8c2c897a199d2526f6f719
 
-Regeneration lands in sprint 3 (M7, the committed DevTools driver); see generate.sh.
+Regenerate with generate.sh (row 201 A2 sprint 4, M7: the passkey page and its DevTools driver) into
+an empty directory; copying its output over these files re-pins XW0 in the same commit.
--- a/host/approval/testdata/webauthn/generate.sh
+++ b/host/approval/testdata/webauthn/generate.sh
@@ -1,7 +1,24 @@
 #!/bin/sh
-# Run as: sh generate.sh. Re-records chrome-es256.json and chrome-eddsa.json with Chrome's
-# virtual authenticator (a test key, D-WORLD-109). Row 201 A2 sprint 1 commits the iteration-379
-# recording (SOURCE.txt); the driver this script runs, host/daemon/testdata/passkey_cdp.mjs,
-# lands in sprint 3 (M7).
-echo "generate.sh: regeneration lands in row 201 A2 sprint 3 (M7); see SOURCE.txt" >&2
-exit 1
+# Run as: sh generate.sh <an empty directory>, with WORLD_CHROME set to Chrome's binary and node on
+# PATH. Re-records chrome-es256.json and chrome-eddsa.json through the passkey page (row 201 A2 M7)
+# with Chrome's virtual authenticator (a test key, D-WORLD-109): cmd/ailang-worldd's
+# TestPasskeyPageRecordsTheFixtures, driving cmd/ailang-worldd/testdata/passkey_cdp.mjs. It writes
+# only into that directory. Copying the files over this directory's is a re-pin: XW0's SHA-256s and
+# every golden value taken from them change in the same commit, and SOURCE.txt names the recording.
+set -eu
+if [ "$#" -ne 1 ]; then
+	echo "usage: sh generate.sh <an empty directory>" >&2
+	exit 2
+fi
+if [ ! -d "$1" ] || [ -n "$(ls -A "$1")" ]; then
+	echo "generate.sh: $1 is not an empty directory" >&2
+	exit 2
+fi
+if [ -z "${WORLD_CHROME:-}" ]; then
+	echo "generate.sh: set WORLD_CHROME to Chrome's binary (only its virtual authenticator is used)" >&2
+	exit 2
+fi
+out=$(cd "$1" && pwd)
+cd "$(dirname "$0")/../../../.."
+WORLD_PASSKEY_RECORD_DIR="$out" WORLD_REQUIRE_PASSKEY_BROWSER=1 go test ./cmd/ailang-worldd/ -run '^TestPasskeyPageRecordsTheFixtures$' -count=1 -v
+ls "$out/chrome-es256.json" "$out/chrome-eddsa.json"
````

### Patch M7g — QUICKSTART §16

````diff M7g
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -1807,3 +1807,33 @@
 ```
 
 Compare the fingerprint `ssh-keygen -l` prints with the one the y/N shows before you type `y`. `enroll` refuses, before any terminal opens, a private key file (without showing it), a key type other than ssh-ed25519 or ecdsa-sha2-nistp256, a name or a key already enrolled (revoked ones included), and a credentials store below schema 5. A revoked name and key can never be enrolled again, and a revoked key approves nothing: `approve --signature` refuses it, and a revoke that lands while an approve is minting refuses that mint. `approve --signature` also refuses a file over 16 KiB or not a regular file, a signature in another namespace or over another request's message, and a key that is not enrolled.
+
+## 16. Approve with a passkey *(attended — pending first verbatim run)*
+
+A passkey approver is a name and one passkey, kept in the credentials store next to §15's SSH approvers. World never sees the passkey: the daemon serves a static page, `/surface/passkey`, that runs the browser's own passkey ceremony (Touch ID, Face ID, a PIN or a security key) and offers the result as a file to save; the CLI reads that file and judges it. The page makes no network request and reads no store. This guide uses the §0 scratch store `/tmp/se-world/world.db` with its daemon serving at `127.0.0.1:7644` (`serve` upgrades the credentials store to schema 6), opened in the browser as `http://localhost:7644`, the only http origin a passkey accepts; never a real store. Use a browser profile and passkey provider you are willing to enrol a test passkey from.
+
+```bash
+ailang-worldd approver enroll --name demo-mac --format webauthn --origin http://localhost:7644 \
+  --registration ~/Downloads/world-registration-demo-mac.json --db /tmp/se-world/world.db
+```
+
+`enroll` refuses first, before any link prints, a registration path that already exists, an origin other than `http://localhost:<port>` or `https://<name>.ts.net`, a name that is not lowercase letters, digits, `.`, `_` and `-`, or a name already enrolled. It then prints a link to `http://localhost:7644/surface/passkey#…` and waits at the terminal. Open the link, check the page says *Create a passkey for approver demo-mac at localhost*, press **Create passkey**, pass your device's check, and press **Save** (the browser saves `world-registration-demo-mac.json` to `~/Downloads`). The page shows a 26-character enrolment code in five groups, `XXXXX-XXXXX-XXXXX-XXXXX-XXXXXX`. Type it at the terminal, or paste it; case, dashes and spaces do not matter. The code is a fingerprint of the file the browser saved, not a secret: it proves that the file the CLI read is the one your browser made, so a file swapped in by anything else is refused. The CLI then shows the passkey (algorithm, rp, `synced passkey` or `device-bound`) and asks `[y/N]` within 2 minutes; the code itself must arrive within 5 minutes. The page gives up a ceremony after 4 minutes and then saves nothing.
+
+```bash
+ailang-worldd approver list --db /tmp/se-world/world.db
+ailang-worldd episode start pk-demo --repo "$PWD" --db /tmp/se-world/world.db --request
+ailang-worldd approve mr-<the 16 hex it printed> --webauthn-link --db /tmp/se-world/world.db
+```
+
+`--webauthn-link` opens no terminal and writes nothing. It prints the request's id, its expiry and the first 16 hex of its message hash, then one link per origin of the enrolled passkeys. Open the link, check the page shows the same request, expiry and hash, press **Approve with passkey**, pass your device's check, and **Save** (`world-assertion-mr-<id>.json`). Then run the command `--webauthn-link` printed, which names the request's own stores:
+
+```bash
+ailang-worldd approve mr-<the same id> --webauthn-assertion ~/Downloads/world-assertion-mr-<the same id>.json \
+  --db /tmp/se-world/world.db --credentials-db /tmp/se-world/credentials.db
+ailang-worldd approve mr-<the same id> --webauthn-assertion ~/Downloads/world-assertion-mr-<the same id>.json \
+  --db /tmp/se-world/world.db --credentials-db /tmp/se-world/credentials.db
+ailang-worldd approve --list --decided --db /tmp/se-world/world.db
+ailang-worldd approver revoke demo-mac --db /tmp/se-world/world.db
+```
+
+The first `--webauthn-assertion` mints with no terminal and records `approved_by=webauthn:demo-mac`; it checks the assertion before any network call, against the enrolled, unrevoked passkey, the request's own message hash, the origin and a signature counter that must move forward. The second, the same file again, is refused as already used. `approve --list --decided` shows `webauthn:demo-mac`; `approver revoke` (a typed `y`) lists the sessions the passkey approved and revokes none of them, and from then on that passkey approves nothing. `approve --webauthn-assertion` also refuses a file over 16 KiB or not a regular file, an assertion for another request, a passkey that is not enrolled, and one made without user verification.
````
