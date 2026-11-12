# Sprint plan — w-world-approvals-by-key-a3a, sprint 4 (iteration 387, row 201)

> **THE SPRINT — A3a sprint 4, the last A3a sprint: the scope the sprint-3 plan's fallback cut moved here
> (its Appendix C, verbatim: the opt-in browser notification of design §2.9 with XA4, the Decided list of
> §2.8 with XA3's Decided half, the `N` key) plus the iteration-386 judge's test gap (A02/A03/A05: three
> `esc()` calls in `needsYouHTML` no test pinned), and the optional live-Chrome hostile-value case (judge
> C03/C04), which fits and is built.**
> - **Measured: 389 changed lines under `cmd/ host/`** — `git diff --numstat <base> HEAD -- cmd host`,
>   insertions + deletions, after a commit of each of the six scratch stages (58 production: `app.js` 58;
>   331 test: `app_test.go` 37, the driver `passkey_cdp.mjs` 43, three Chrome test files 251).
>   `docs/QUICKSTART.md` §17 is **4** lines more (outside the filter). Margin to the 2,000 cap: **1,611**.
>   No fallback cut.
> - **No Go production change, no route, no CSP change.** The data was already served: `GET
>   /surface/requests.json` carries `decided` (the last 20) since sprint 2 (§1.2 C3). Everything new is in
>   `host/surfacepage/app/app.js`: two pure, node-tested functions (`decidedHTML`, `notifications`) and one pure
>   renderer change (`needsYouHTML` gains the opt-in), plus their effects at the page boundary (the click
>   handlers, the `Notification` constructor, `requestPermission`, the `N` key). `surfaceCSP` is byte-unchanged
>   (a notification is a JS API, not a fetched resource; XA4 sees zero CSP violations).
> - **The judge's survivors are closed and proved** (§4, rows `ESC-*`): at the base, removing `esc()` on
>   `q.agent`, the link's `id=`, its `href`, `q.recorded`, the refused row's id and the `unavailable` text each
>   leaves `TestSurfaceAppPure` **green** (6 survivors, measured); at M9a (6 test lines), each fails it.
> - **Every line of code, test, JS and doc is a unified diff in Appendix P**, generated LAST from compiled,
>   vetted, gofmt-clean, green scratch stages; the patches were then extracted **from this file** with §2's
>   `awk`, applied to a fresh `git archive bed9d135` export, each sub-milestone's acceptance run on it, and the
>   result compared byte-for-byte with the mutant tree (§1.2 C27). The executor applies each patch and runs
>   its acceptance; it writes no code by hand.

- **Design**: [w-world-approvals-by-key-a3.md](w-world-approvals-by-key-a3.md), **r2** (quorum r1+r2 +
  controller carve-out + independent audit VALID 84; **not reopened**). D-WORLD-114 (no Tauri half in A3a; the
  CLI minter re-judges and mints) is OPEN with default A = keep both, applied since iteration 384: built as
  designed. The notification here is the **web** notification of design §2.9 (A3a), not the native/Tauri one
  D-WORLD-114 (a) defers. Where this plan and the design differ, this plan was measured against `bed9d135` and
  wins; each difference is a decision (§1.3) or a premise found false (§1.4).
- **Starting ground**: A3a sprints 1–3 merged (`b44b5f7d` PR #554; `d1027477` PR #556; `70823879` PR #558,
  recorded `a141f0ba`): schema 7, the shared judge, `approve <id> --await`, the read routes, the origin set, the
  POST (both formats), the card page, Needs you's pending list. The sprint-3 plan's §1.2–§1.4 (S3-D1…S3-D16,
  P1–P8) are facts of HEAD this plan builds on.
- **Planner base**: `bed9d135fe38a1c8f7e2784414aaac348b76879f` (= origin/dev). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter387-a3a-s4`, branch `sprint/world-iter387-row201-a3a-s4`.
  Planner scratch (outside the repo, never under `/tmp`): `$HOME/.ailang/state/world-iter387/plan-scratch/`
  (`base/` the pristine export, `work/` the stage repo (tags `M9a` … `M10c`) and the tree the mutants were cut
  from, `m9a/` the M9a tree (the kill proofs), `patches/`, `tools/`, `logs/`, `roundtrip/`, `mut/`).
- **Executor**: as routed by the controller. **Two runs**: Run 1 (M9a, M9b, M9c — `app.js` and its node
  test), Run 2 (M10a, M10b, M10c — the driver and the live-Chrome tests). Each sub-milestone is one `patch -p1`
  (§2) plus its acceptance (§3). **Apply, never retype.** From M9b on the acceptance runs **Chrome headless**
  (temporary profiles, a CDP virtual authenticator, test keys only).
- **Touched** (§5): `host/surfacepage/app/app.js` (M9b, M9c, M10a), `host/surfacepage/app_test.go` (M9a, M9b,
  M9c), `cmd/ailang-worldd/needs_you_chrome_test.go` (M9c one line, M10a), `cmd/ailang-worldd/testdata/
  passkey_cdp.mjs` (M10a, M10b), new `cmd/ailang-worldd/needs_you_notify_chrome_test.go` (M10b), new
  `cmd/ailang-worldd/approve_page_hostile_test.go` and `docs/QUICKSTART.md` §17 (M10c). **Outside the design's
  §8 file list**: `app_test.go` (the node harness, as in sprint 3), `needs_you_chrome_test.go` (sprint 3's
  split of the design's "`approve_page_*_test.go`"), `needs_you_notify_chrome_test.go` and
  `approve_page_hostile_test.go` (the same split; the latter is the directive's optional item 5).
- **Not touched**: every Go production file (`git diff --stat bed9d135 M10c -- '*.go' ':!*_test.go'` → empty),
  `host/approval/*`, `host/store`, `host/authority`, `host/daemon/*` (so `surfaceCSP`, `isProtected`, every route,
  `/v1`), `host/surfacepage/surface.js`, `approve.js`, `passkey.js`, both `index.html`s, `app.css`, `model.go`
  (`needsYou` stays `[]`), `world/*.ail`, `packages/**`, `scripts/**`, `tools/launchd/*`, `cmd/world-publish/*`,
  any skill, `~/.ssh`, `~/.ailang/world`, `~/.ailang/state/mission-v1*`.

---

## 0. Binding constraints (each is restated where it governs)

1. **Test keys, temp stores and virtual authenticators only (D-WORLD-109 = A).** Every store is under
   `t.TempDir()`; the one passkey is `approvaltest.NewPasskey`'s fixed ES256 seed handed to Chrome's CDP
   virtual authenticator (`virtualPasskey`, sprint 3 S3-D11). The daemon is `httptest.NewServer` on the
   loopback. Chrome runs `--headless=new` with `--user-data-dir` under the test's temp dir, killed and removed
   at the end. **No notification is ever raised**: XA4's `notify` recorder replaces `window.Notification` in the
   test tab and never calls the real constructor (S4-D8); the CDP grant is to the test's loopback origin in the
   temporary profile. Nothing opens `~/.ailang/world`, `~/.ssh`, a real store, a real browser profile, a real
   `/dev/tty`, or runs `ssh-keygen`; nothing mints into a real store or publishes.
2. **Pure core, effects at the boundary (coding standards S1–S3).** What the page shows and what it would
   notify are pure functions of `requests.json` (`needsYouHTML`, `decidedHTML`, `notifications`), exported for
   node and pinned by `TestSurfaceAppPure` (CI runs it with `WORLD_EXEC_NODE`). The effects — the click
   handlers, `Notification.requestPermission`, `new Notification`, its click, the `N` key — are three short
   lines at the page boundary, pinned by the live-Chrome tests. No `.ail` (design §3), no Go change, no new
   module, package, exec site, route or context root.
3. **No `innerHTML` with unescaped data.** Every new value that lands in HTML goes through the file's `esc()`
   (`decidedHTML`'s episode and approver; the times are `hm()` digits; the counts are numbers); the notification
   body is plain text given to the `Notification` constructor, never HTML. No `eval(`, no `document.cookie`.
4. **One mutant per new production branch, anchored to the diff** (§4), each run against its single named test
   **and** its package suite, inventories excluded from every suite; named equivalents with reasons; every run
   appends a ledger line.
5. **No `t.Parallel()`** in any new test (`TestNoParallelWriteForkPackagesOutsideVerifygate`).
6. **Bisectable and safe at every commit.** Each sub-milestone ends gofmt-clean, `go build ./...` and `go vet`
   green, the `_test.go` compile fence (`go test -run '^$' ./...`) green, its own tests and every earlier one
   green, **including the live-Chrome XA3** from M9b on (M9c carries XA3's one-line expectation for the opt-in).
7. **No git writes by the executor** (§7). The controller commits from snapshots.

---

## 1. Pre-flight

### 1.1 Baseline at `bed9d135` (measured by the planner, 2026-11-12)

Every command ran with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` (AILANG v0.41.0), `go1.26.6
darwin/arm64`, node `v26.13.0`, `Google Chrome 157.0.8103.51`, in a `git archive bed9d135` export at
`plan-scratch/base` (a throwaway git repo; `diff -r cmd host` against the worktree: identical,
`BASE_EQUALS_WORKTREE`; the worktree itself was only read). The pre-commit legs are `tools/gates.sh <tree>
<logs> base` = §6 without the sprint's new Chrome test. Logs: `plan-scratch/logs/baseline/`.

| Command | rc | Time | Result |
|---|---|---|---|
| `git -C <worktree> status --porcelain \| wc -l` | 0 | — | `0` (clean at `bed9d135`) |
| `./scripts/verify_ail.sh` | **0** | — | `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` |
| `gofmt -l cmd host`; `go build ./...`; `go vet ./...`; `go test -run '^$' ./...` | **0** | — | empty / clean / empty / no FAIL line |
| `go test -count=1 -timeout 20m ./host/surfacepage/ ./host/store/ ./host/authority/ ./host/approval/... ./cmd/ailang-worldd/ ./host/daemon/` | **0** | — | all `ok` (`host/surfacepage` 0.54 s, `host/store` 12.0 s, `host/authority` 0.68 s, `cmd/ailang-worldd` 124.2 s, `host/daemon` 284.8 s) |
| `-race` legs (§6's three) | **0** | — | `ok host/surfacepage` 3.0 s, `ok host/daemon` 5.0 s, `ok cmd/ailang-worldd` 24.8 s; no race report |
| §6's named inventory, schema, reflective, help, quickstart and route tests, `-v` | **0** | — | **57** `--- PASS`, 0 SKIP, 0 FAIL; the AW leg 9 `--- PASS: TestApproveAwait` |
| the live-Chrome leg (`WORLD_CHROME=… WORLD_REQUIRE_PASSKEY_BROWSER=1`, A2's and sprint 3's browser tests) | **0** | — | `--- PASS: TestApproveCardNetworkDeadlines (79.31s)`, `TestApproveCardPasskeyDeadline (1.29s)`, `TestApprovePageInChrome (1.44s)`, `TestApprovePageRefusesASwappedRequest (1.62s)`, `TestApprovePageSSHSIGInChrome (1.08s)`, `TestNeedsYouListsPendingRequests (1.06s)`, `TestPasskeyPageAbortsAPendingCeremony (3.22s)`, `TestPasskeyPageInChrome (7.75s)` |
| the whole pre-commit script (`tools/gates.sh base`) | **0** | **521 s** | `PRE_COMMIT_GATES_GREEN 521s` |
| `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` (the base export, committed) | **1** (each red green alone) | — | measured while the planner's round trip and mutant runs loaded the machine (649 s): `✓ 0 other binary blobs among 6822 tracked files`, `✓ toolchain floor gate: go1.26.6 >= root module floor go1.26.6`, `✓ all 64 required top-level evidence tests passed exactly once`, then three untouched reds, all on the known list: `--- FAIL: TestWorldContextAndStatusline` (cmd/ailang-worldd), `--- FAIL: TestWorktreeSnapshotRefusesAConfiguredFilter` (host/broker, "control: plain git status did not run the clean filter"), `--- FAIL: TestMissionConfigDefaultFlowIntegration/valid_ledger` (host/verifygate, "default-flow fixture timed out … context deadline exceeded"); every other package `ok` (`logs/base_verify_go.log`). Re-run alone (`logs/rerun_alone.log`): `--- PASS: TestWorldContextAndStatusline (2.82s)`, `--- PASS: TestWorktreeSnapshotRefusesAConfiguredFilter (0.51s)`, `--- PASS: TestMissionConfigDefaultFlowIntegration (9.96s)` with `/valid_ledger` and `/invalid_ledger` PASS. A red here measures the repo, not the sprint (the rule below) |

**Rule (carried):** a red at the base measures the repo, not the sprint. The controller re-runs a red
untouched package alone once and does not chase it. **Never start a gate under `nohup`.**
`VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` reads committed blobs and runs `go test ./...`; it is the
after-commit leg (§6).

### 1.2 Codebase claims checked (at `bed9d135` unless "work" — the planner's final scratch tree)

| # | Claim | Command → observed output |
|---|---|---|
| C1 | The base | `git rev-parse HEAD origin/dev` (worktree) → `bed9d135fe38a1c8f7e2784414aaac348b76879f` twice; `git status --porcelain \| wc -l` → `0` |
| C2 | **Nothing of sprint 4 has landed** | `grep -c Notification host/surfacepage/app/app.js` → `0`; `grep -rln grantNotifications cmd host` → no match; control in the same call: `grep -rln 'getRecorder\|fetchRewrite' cmd host` → `testdata/passkey_cdp.mjs`, `approve_page_driver_test.go`, `approve_page_swap_test.go`, `approve_page_chrome_test.go`; `grep -n decided host/surfacepage/app/app.js` → no match (control: `grep -c pending app.js` → `2`) |
| C3 | **`requests.json` already carries the decision list** (no Go change needed) | `host/surfacepage/approve_model.go:24 const MaxDecidedListed = 20`, `:68 Decided []DecidedRequest \`json:"decided"\`` with `DecidedRequest{ID, Episode, ApprovedBy \`json:"approvedBy"\`, At \`json:"at"\` (consumed_at, unix s), Credential16}` (`:54-60`); `host/daemon/surface_approve.go:128 decided, err := ro.DecidedMintRequests(ctx, "", surfacepage.MaxDecidedListed)`, `:133 out.Decided = append(...)`; `host/store/mint_requests.go:229 … ORDER BY consumed_at DESC, request_id DESC LIMIT ?` (newest first) |
| C4 | **`needsYouHTML` and its export** (design drift: §2.8 cites `app.js:605-606`) | `app.js:426 function needsYouHTML(r, now) {`; `:435 if (typeof module !== 'undefined') module.exports = { needsYouHTML, …`; the now card's call `:616 ${needsYouHTML(REQS, REQS ? REQS.now : 0)}`; `:893 let REQS = null;`, `:894 async function loadRequests()`, called in `loadModel` (`:902`) and at the start of `refresh` (`:993`) |
| C5 | **The judge's survivors are real at the base** (A02/A03/A05) | the node harness's only `needsYouHTML` rows use `agent: 'claude'`, id `mr-0123456789abcdef`, `recorded: 'webauthn:mark'`, refused id `mr-00000000000000ff`, unavailable `'the mint requests could not be read'` (`app_test.go:180-183`): no character `esc()` changes. Measured (`SRC=base`, `logs/mutant_ledger.tsv` run `base-esc`): removing `esc()` on `q.agent`, the `id=`, the `href`, `q.recorded`, the refused row's id and `unavailable` → `TestSurfaceAppPure` **SURVIVED** (named and suite, 6 rows); on `q.episode` and `q.refused` → KILLED (sprint 3's `w<1>` / `digest <bad>` rows) |
| C6 | **The `N` key is not the app's** (design §2.8 "the key already exists, `surface.js:257`") | `host/surfacepage/surface.js:156 case "n": case "N": return {ui: "needs"};` and `:257`/`:259` the help rows — the **static pages'** keys; `:238 if (document.documentElement && document.documentElement.getAttribute("data-app") != null) return;` (the app shell `index.html` has `<html lang="en" data-app>`, so `surface.js` stops before its key handler); `grep -n "'n'\|'N'" app.js` → no match; control: `app.js:1380 if (k === '?') toggle('keys');` |
| C7 | **The now card's Enter is the timeline's** (found building `N`, §1.4 P2) | `app.js:1379 if (k === 'Enter' && e.target && e.target.id === 'home') return;` then `:1389 else if (st.d === 0 && k === 'Enter') { if (st.focus) open(…st.focus) }` and `:1407 if (done) e.preventDefault();` — Enter on a focused Needs-you link opens the focused timeline state and cancels the link's own activation (work: with S4-D7's line removed, mutant `ENTER-NOW`, XA3's `key Enter` never reaches the card, §4) |
| C8 | **The app refreshes on the live log, and polls only in follow mode** | `app.js:945 if (!followMode \|\| !liveVisible()) return;` (`schedulePoll`), `:897 POLL_MS = 15000`; `refresh('live')` on a committed log entry (`queueRefresh`); `:1413 poll: () => refresh('poll')` (the driver seam); `:988 if (refreshing) { if (trigger === 'live') queued = true; return; }` — a `poll` during a refresh is dropped (§1.4 P4) |
| C9 | **The app's clock is UTC `hh:mm`** | `app.js:13 const hm = iso => { const d = new Date(iso); return pad(d.getUTCHours()) + ':' + pad(d.getUTCMinutes()); };` used by the now card's `Now · hh:mm` (`:608`) |
| C10 | **`surfaceCSP` and the shell are unchanged** | `host/daemon/surface_page.go:19 const surfaceCSP = "default-src 'none'; style-src 'unsafe-inline'; script-src 'self'; connect-src 'self'; font-src 'self'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'"`; work: `git diff --stat base M10c -- host/daemon host/surfacepage/app/index.html host/surfacepage/app/app.css host/surfacepage/surface.js host/surfacepage/approve` → empty. CSP has no directive for the Notifications API (it governs fetched resources, scripts, styles, frames); work: XA4 asserts `window.__cspv` empty at every step (`--- PASS`) |
| C11 | **`app.css` is pinned to the working copy above its "After v3" marker** | `app_test.go:39 cut := strings.Index(style, "\n/* ---------- After v3: the brand is the way home")`; so the sprint adds **no CSS**: the two buttons reuse `class="nl run"` (`openRunningHTML`'s running-list buttons) and inline styles stay as they are |
| C12 | **Headless Chrome 157: the Notifications API, the CDP grant, a gesture-bound request** (directive retro 4; design V25/V26 re-measured live) | work, XA4 (`--- PASS`): before any grant `Notification.permission` = `default` and the page never calls `requestPermission` (`asks: 0`); `Browser.grantPermissions({permissions: ['notifications'], origin: http://localhost:<port>})` → the page reads `granted`; the opt-in click calls `requestPermission` once with `navigator.userActivation.isActive === true`, it resolves, and the control repaints to "Notifications on"; a constructed notification's `onclick` (the recorder's) runs the page's handler: `window.focus` counted 1, the tab at `/surface/approve/<id>` |
| C13 | **A Promise through the driver's `waitFor` is not awaited** (harness trap, §1.4 P4) | `passkey_cdp.mjs` `evaluate` sends `JSON.stringify(${expression})` with `awaitPromise: true`: a Promise becomes `"{}"` before it settles. Work, first XA4: `(window.World.poll(), 1)` returned before the refresh; the next `poll` was dropped (C8) and the wait for "Needs you · 3" timed out once in three runs; `pollSurface` (a flag set in `.then`) fixed it: 3 of 3 green |
| C14 | **`World.open` opens a brane asynchronously** | work: `window.World.open('e2')` returned `true` while `stage.dataset.view` was still `timeline`; a key pressed at once was handled on the timeline and the brane then opened over it (`N in a brane: {… Active:root View:brane}`); the test waits for `view === 'brane'` first |
| C15 | **The card renders signed values as text** (judge C03/C04) | `host/surfacepage/approve/approve.js:133 dd.textContent = shown(f, l.value);` (and `:6` "Text is set with textContent only"); at the base no Chrome test renders a hostile value (sprint 3's fields are a path, a model `claude-opus-5-5`, …): mutant `CARD-TEXTCONTENT` (`textContent` → `innerHTML`) at the base **SURVIVES** `TestApprovePageInChrome` and the Chrome suite (§4) |
| C16 | **Direct SQL marks decided rows in tests** (precedent for XA3's three extra decided rows) | `cmd/ailang-worldd/approve_test.go:489 func rawCredentials`; `approve_await_refusals_test.go:92 raw.Exec(\`UPDATE mint_requests SET consumed_at = ?, approved_by = 'tty', credential_id = 'elsewhere' WHERE request_id = ?\`, …)`; the schema's `CHECK ((consumed_at IS NULL) = (approved_by IS NULL) AND (consumed_at IS NULL) = (credential_id IS NULL))` (`credentials_schema.sql:55`) |
| C17 | **No test pins QUICKSTART §17** | `grep -n '^## 17' docs/QUICKSTART.md` → `1841: ## 17. Approve from the surface *(attended — pending first verbatim run)*`; `grep -rln 'Approve from the surface\|## 17\.' --include='*_test.go' cmd host scripts` → no file (rc 1); control: `grep -rln '## 1[01]\.' --include='*_test.go' cmd host` → `quickstart_exec_test.go`, `quickstart_attended_test.go`; the edit adds no code fence and no link (the quickstart and marketing tests that read fences and links stay green, §6) |
| C18 | **CI runs node but not Chrome** (design R3) | `grep -cE 'WORLD_CHROME\|google-chrome\|chrome' .github/workflows/ci.yml` → `0`; `ci.yml:250 echo "WORLD_EXEC_NODE=$(command -v node)" >> "$GITHUB_ENV"`, `:289 : "${WORLD_EXEC_NODE:?WORLD_EXEC_NODE is required}"` — `TestSurfaceAppPure` runs in CI; the XA tests are counted SKIPs there |
| C19 | **Tests that enumerate what this sprint touches** (directive retro 5) | `grep -rn "aria-keyshortcuts\|KEYROWS\|these keys" --include='*_test.go' --include='*.mjs' host cmd scripts` → no match (the keys table and the shell's key list are unpinned); the export list is pinned by the node harness's `assert.deepEqual(Object.keys(app).sort(), […])` (`app_test.go:174`, edited M9b/M9c); `TestSurfaceAppShellAssembly` checks seam names and the `module.exports` line (unchanged shape); embedded files: `TestSurfaceAppRoutes`/`TestPasskeyPageRoutes`/`TestApprovePageRoutes` (no file added or renamed) |
| C20 | **Every repo-wide inventory, on the final tree** | the §6 script (`tools/gates.sh <tree> <logs> final` = §6's pre-commit block, `WORLD_EXEC_NODE` exported) on the round-trip tree (byte-identical to the mutant tree, C27) → `PRE_COMMIT_GATES_GREEN 548s`; pkgs: `host/surfacepage` ok 1.922s, `host/store` ok 14.128s, `host/authority` ok 1.747s, `host/approval` ok 2.100s, `host/approval/approvaltest` ok 1.328s, `cmd/ailang-worldd` ok 158.328s, `host/daemon` ok 302.395s (`logs/final_gates/`): design §7's inventories by name — `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` (no exec site), `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist`, `TestBareNetHTTPExemptionIsPerGroup`, `TestEveryCommandDirectoryIsAProtectedGroup`, `TestProductionContextRoots` (no root added), `TestProductionGoSurface`, `TestNoNewDeadlineFreeStoreReads`, `TestReadGettersHonorContext`, `TestCredentialsMethodsRequireADeadline`, `TestFrozenV1RouteTableMatchesMux` (13), `TestSurfaceAppRoutes`, `TestMCPMountMethodSource`, `TestReadCtxCancelledAfterHandler`, `TestDaemonCredentialWriterSites` (unchanged), `TestApprovalImportsOnlyTheStandardLibrary`, `TestApprovalTestIsTestOnly`, `TestNoParallelWriteForkPackagesOutsideVerifygate`, `TestSubprocessSinkGate`, `TestEveryVerbHasHelp`, `TestServeProofFlags`, the quickstart tests (`TestQuickstartAttendedStepsBoxMatchesSection9`, `…Section9/10/11FlagsMatchTheCLI`, `TestQuickstartServiceSectionMatchesTheCLIHelp`, `TestQuickstartSection10ProfilesLoad`, `TestQuickstartProtectedContract`, `TestQuickstartDecisionPacketSection`), the schema/judge pins, sprints 2–3's route and page tests, `TestSurfaceAppPure` — **57 named `--- PASS`**, 0 SKIP/FAIL; the AW leg 9 `--- PASS: TestApproveAwait` |
| C21 | **The live-Chrome tests, on the final tree** | §6's live leg (`WORLD_CHROME=… WORLD_REQUIRE_PASSKEY_BROWSER=1 go test -v ./cmd/ailang-worldd/ -run '^(TestPasskeyPageInChrome\|TestPasskeyPageAbortsAPendingCeremony\|TestApprovePage\|TestApproveCard\|TestNeedsYou)'`, Chrome 157.0.8103.51, `logs/final_gates/chrome.log`): **new** `--- PASS: TestNeedsYouNotifiesOnlyWhenAskedOnce (1.66s)` (XA4), `--- PASS: TestNeedsYouListsPendingRequests (1.93s)` (XA3 with its Decided half, Show all and `N`), `--- PASS: TestApprovePageShowsAHostileValueAsText (1.25s)`; carried `--- PASS: TestApproveCardNetworkDeadlines (79.61s)` with `/get (22.27s)` and `/post (57.34s)`, `TestApproveCardPasskeyDeadline (1.42s)`, `TestApprovePageInChrome (1.65s)`, `TestApprovePageRefusesASwappedRequest (2.14s)`, `TestApprovePageSSHSIGInChrome (1.72s)`, A2's `TestPasskeyPageAbortsAPendingCeremony (3.58s)` (4 subtests PASS) and `TestPasskeyPageInChrome (3.43s)` (3 subtests PASS); `ok cmd/ailang-worldd 98.703s`; 0 SKIP, 0 FAIL. CI skips the three new ones (counted, C18) |
| C22 | **Every mutant, on the final tree** | 60 rows, the last ledger line of each on `work` (= M10c): **59 KILLED (named and suite), 1 named equivalent SURVIVED (named and suite)** (`NP-NOAPI`), 0 not as expected, 0 `NOCOMPILE`; 77 ledger lines in all (runs `run1-node` 36, `run1-chrome` 24, `base-esc` 8, `m9a-esc` 8, `base-card` 1). Every kill's failing tests are the named behavioural test(s) (`TestSurfaceAppPure`, `TestNeedsYouListsPendingRequests`, `TestNeedsYouNotifiesOnlyWhenAskedOnce`, `TestApprovePageShowsAHostileValueAsText`); no suite includes an inventory, and no row is killed by `TestSurfaceAppShellAssembly` or `TestSurfaceAppMatchesWorkingCopy` |
| C23 | `go vet ./...`, gofmt, `-race`, `verify_ail.sh` on the final tree | same §6 run: `go vet ./...` empty; `gofmt -l cmd host` empty; the compile fence clean; `-race` → `ok host/surfacepage 5.045s`, `ok host/daemon 7.048s`, `ok cmd/ailang-worldd 28.199s`, no race report; `./scripts/verify_ail.sh` → `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` (no `.ail` change) |
| C24 | **Nothing outside §5 moved** | work: `git diff --stat bed9d135 M10c -- '*.go' ':!*_test.go' \| wc -l` → `0`; `git diff --stat bed9d135 M10c -- host/daemon host/surfacepage/app/index.html host/surfacepage/app/app.css host/surfacepage/surface.js host/surfacepage/approve host/surfacepage/passkey host/approval host/store host/authority \| wc -l` → `0`; added lines matching `innerHTML\|eval(\|document.cookie` in `app.js` → `0`; added `t.Parallel` → `0`; `git diff --stat bed9d135 M10c` → 7 files (§5), `357 insertions(+), 36 deletions(-)` |
| C25 | **The whole repository on the final tree** (committed) | `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` on the committed round-trip tree at M10c, nothing else running → **rc 1** (514 s), one untouched package red, on the directive's known list: `host/verifygate` (`--- FAIL: TestMissionConfigDefaultFlowIntegration/valid_ledger`: `default-flow fixture timed out … instrument failure: context deadline exceeded`, the same failure the base showed, §1.1); every other leg `✓` (`✓ 0 other binary blobs among 6824 tracked files`; `✓ toolchain floor gate: go1.26.6 >= root module floor go1.26.6`; `✓ all 64 required top-level evidence tests passed exactly once`) and every other package `ok` (`logs/final_verify_go.log`). **Re-run alone** (§6's rule; `logs/rerun_alone.log`): `--- PASS: TestMissionConfigDefaultFlowIntegration (7.96s)` with `/valid_ledger` and `/invalid_ledger` PASS, `ok host/verifygate 8.253s`. `host/verifygate` is not touched by the sprint (§5). **The whole gate re-run once, alone** → **rc 1** (470 s) on a different untouched test: `host/broker` `--- FAIL: TestPublishApprovalRefusalSetWithALandedPositiveControl` (`registry_publish_test.go:2399: AC9 positive control counters = POST=2 dispatches=1 credentialLoads=1, want POST=1 …`; `host/verifygate` was `ok` this time); every leg `✓` as above, every other package `ok` (`logs/final_verify_go_2.log`). Re-run alone on the final tree and at the base: `--- PASS: TestPublishApprovalRefusalSetWithALandedPositiveControl (0.36s)` / `(0.38s)`, all 21 subtests PASS, `ok host/broker` (`logs/rerun_alone.log`). It is **not** on the directive's known list: a load-sensitive counter in a package the sprint does not touch (§1.4 is unaffected; flagged for the controller). Across the two runs every package was `ok` at least once and every red was green alone. |
| C26 | **Every sub-milestone is bisectable** | the fresh export of C27, after **each** patch: `tools/accept.sh <tree> <M>` = §3 verbatim (incl. the `[chrome]` legs, live) → `ACCEPT M9a rc=0` 71s, `ACCEPT M9b rc=0` 33s, `ACCEPT M9c rc=0` 31s, `ACCEPT M10a rc=0` 40s, `ACCEPT M10b rc=0` 34s, `ACCEPT M10c rc=0` 255s (`logs/roundtrip.log`); the cumulative numstat after a commit of each = §7: M9a **6**, M9b **36**, M9c **93**, M10a **239**, M10b **357**, M10c **389** (docs 4) (`NUMSTAT` lines of `logs/roundtrip.log`) |
| C27 | **The patches reproduce the mutant tree** | a fresh `git archive bed9d135` export + `patch -p1` of Appendix P's six patches **extracted from this file with §2's `awk`** (each `cmp`-equal to the generated patch; `--dry-run` first; no `.orig`/`.rej`) → `diff -r cmd host docs` against `plan-scratch/work` (the tree §4's mutants were cut from): **identical** (`ROUNDTRIP_IDENTICAL_TO_MUTANT_TREE`) (`tools/roundtrip.sh`, `logs/roundtrip.log`); nothing uncommitted after M10c. A last extraction-only re-check after the final assembly (`SKIP_ACCEPT=1`, a second fresh export) again printed `ROUNDTRIP_IDENTICAL_TO_MUTANT_TREE` with the same six numstat values. |
| C28 | Size per stage (work) | `tools/numstat.sh` → M9a own 6 (0 prod + 6 test), cumulative **6**; M9b own 30 (19 + 11), **36**; M9c own 61 (39 + 22), **93**; M10a own 148 (4 + 144), **239**; M10b own 120 (0 + 120), **357**; M10c own 32 (0 + 32), **389**; docs 4 (M10c) |

### 1.3 Planner decisions where the design leaves a choice

- **S4-D1 — No cut.** The sprint-3 plan estimated the moved part at ≈ 300; measured, with the judge's test gap
  and the optional hostile-value test, it is 389 (§1.5). Every item of the directive's scope is built,
  including item 5.
- **S4-D2 — What a read of `requests.json` notifies is one pure function**, `notifications(r, seen,
  permission)` → `{notes, seen}`: a note `{title: "World needs you", tag: <id>, body: "Approve mint request for
  <episode> (<n> grant(s)), expires <hh:mm>", href: "/surface/approve/<id>"}` for each pending request whose id
  is not in `seen`, only when `permission === 'granted'`; `seen` gains every pending id read, whatever the
  permission (so a request seen before the opt-in never notifies later: design XA4's "a new request").
  - `seen` is `null` before the page's **first** read, which notifies nothing: design §2.9 says "on a
    **refresh**", and the page the operator just opened shows those requests itself. A failed read (no
    `pending` array) leaves `seen` as it was.
  - Refused rows (no summary to name, a card that refuses) and rows already `recorded` (approved elsewhere,
    waiting for the minter) are seen and never notify.
  - `loadRequests` (the one place `requests.json` is read: the page's load and every refresh, no new timer)
    calls it with the browser's current `Notification.permission` and raises each note: `new
    Notification(title, {body, tag})`, whose `onclick` is `window.focus(); location.assign(href)` (the same
    tab; design §2.9 "focuses the tab and opens"). One per id per tab: `seen` holds it.
- **S4-D3 — The opt-in is a renderer argument.** `needsYouHTML(r, now, permission)` ends its section with
  `notifyHTML(permission)`: `default` → a **Notify me** button (`id="notify-me"`); `granted` → "Notifications
  on"; `denied` → "Notifications are blocked for this page"; anything else (no API) → nothing. The two-argument
  calls of sprint 3's node rows render exactly as before. The button's click (and only it) calls
  `askNotify()` = `Notification.requestPermission()` then a repaint of the now card with the answer.
- **S4-D4 — Times are the app's clock, UTC `hh:mm`** (`hm`, C9): the notification's expiry and each Decided
  line's time, as the now card's own "Now · hh:mm". The design's literal `(<n> grants)` agrees in number
  ("1 grant"), as the Needs-you line already does.
- **S4-D5 — Decided is `decidedHTML(r, all)`**, a second section under Needs you: "Decided", then one line per
  decided request, newest first (the daemon's order, C3): `You · Minted · <episode> · <approvedBy> · <hh:mm>`
  (design §2.8; every decided row is a mint: `consumed_at IS NOT NULL`); the last **3**, and, when there are
  more, a **Show all <n>** / **Show fewer** button (`id="decided-all"`, `aria-expanded`) that toggles the page's
  `decidedAll` and repaints the now card. Nothing at all when there is no decided request.
- **S4-D6 — The page's state is three variables beside `REQS`**: `decidedAll`, `seenIds`, and
  `notifyPermission()` read live; the two click handlers join the stage's existing click delegation (as the
  quiet-lane toggle does). No timer, no storage, no route.
- **S4-D7 — `N` is new in `app.js`** (§1.4 P1). In the app's keydown handler, `n`/`N` focuses the first
  `[data-needs] a` (the first pending request's link), leaving a brane first (`out()`) if one is open; with no
  link it does nothing. The keys popover gains the row `N · what needs you`. **And Enter inside the now card is
  the browser's** (`if (k === 'Enter' && e.target.closest('[data-now]')) return;`, beside the existing `#home`
  rule): without it the timeline's Enter swallows the link `N` just focused (C7). This also lets Enter press
  the now card's other buttons (Notify me, Show all, the running list), which previously opened the focused
  timeline state instead.
- **S4-D8 — XA4's recorder raises nothing.** The driver's `notify` injection replaces `window.Notification`
  in the test tab with a recorder: the constructor records `{title, body, tag}` in `window.__notes` and
  returns a plain object whose `onclick` the page sets; `permission` and `requestPermission` pass through to
  the browser (`requestPermission` also records `navigator.userActivation.isActive`); `window.focus` calls are
  counted in `sessionStorage` (which outlives the navigation the click starts). The real constructor is never
  called (D-WORLD-109, the directive's "never trigger a real desktop notification"); design V26 measured it
  constructs under the CDP grant, and **whether one is displayed stays attended** (design §2.9, §10 step 4).
- **S4-D9 — XA3 is extended in place** (`TestNeedsYouListsPendingRequests`, design §6 "XA3 … after the mint,
  Decided shows `webauthn:<name>`"): the request is approved on its card through the test passkey
  (`virtualPasskey`) and minted by `approve --await` (a real mint, `requireMintedBy`), then `/surface` shows
  `DecidedYou · Minted · ep1 · webauthn:mark-mac · <its consumed_at as UTC hh:mm>`; three more requests are
  marked decided by SQL (C16; they are fixtures for the toggle, not mints) and Decided shows 3 lines and **Show
  all 4**, toggling to 4 and **Show fewer** and back. The same test drives `N` on the timeline and `N` (Shift)
  in a brane, then Enter on the focused link opens the card (it replaces sprint 3's mouse click, whose
  `href` assertion stays).
- **S4-D10 — The optional live-Chrome hostile value is built** (directive item 5, judge C03/C04):
  `TestApprovePageShowsAHostileValueAsText` makes a request whose `model` is
  `<img src=x onerror="window.__xss=1"><script>window.__xss=2</script>"&amp;` (through `episode start
  --model`, which accepts it), opens its card and asserts the model field's text is exactly that string, the
  card is `ready` with no page error and no CSP violation, `window.__xss` is unset, the document has no `<img>`
  and only its own `<script>`. 32 lines.
- **S4-D11 — Driver additions** (the A2 driver, extended, never replaced): commands `key` (CDP
  `Input.dispatchKeyEvent` for `n`, `N` with Shift, `Enter`; optionally waits for the next document; returns
  the location, the focused id and the app's view), `tap` (a click by element id on any page; `press` is the
  card's), `grantNotifications` (`Browser.grantPermissions` for an origin), `notification` (runs a recorded
  notification's `onclick` and waits for the next document); injection `notify` (S4-D8). Nothing existing
  changes.
- **S4-D12 — Test helpers live with XA3** (`rigSurface`, `sectionText`, `pollSurface`, `keyView` in
  `needs_you_chrome_test.go`); `pollSurface` sets a page flag in `World.poll()`'s `.then` and waits for it, so
  a refresh has finished before the next one starts (C8, C13).
- **S4-D13 — Sub-milestones** (§1.5): M9a the A02 kill test (6), M9b Decided (30), M9c the opt-in and the
  notification (61), M10a `N`, Enter, the driver's `key`/`tap` and XA3 extended (148), M10b XA4 and its driver
  (120), M10c the hostile value and QUICKSTART §17 (32 + 4 docs). Each is under 350.

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

- **P1 — Design §2.8 "`N` jumps to the first one (the key already exists, `surface.js:257`)".** The key and its
  help row exist only in the static pages' `surface.js`, which returns before its key handler on the app shell
  (`data-app`, C6); `app.js` has no `N`. Built new in `app.js` (S4-D7). Carried from the sprint-3 plan's P7.
- **P2 — An `N` that focuses a link presumes Enter follows it.** In the app the timeline's Enter handles every
  Enter not on `#home` and cancels the default (C7). Routed: Enter inside the now card is the browser's
  (S4-D7), pinned by XA3 (mutant `ENTER-NOW`).
- **P3 — Design §2.8 "fetched on the app's existing refresh tick".** The app has a 15 s tick only in follow mode;
  in v1 mode it refreshes on a committed log entry (C8). So a new request is listed, and notified, on the next
  refresh: within 15 s under `serve` with a follow config, else at the next world-log commit or reload. Kept as
  designed (no new timer, design §2.8); the tests drive the app's own refresh (`window.World.poll()`), and
  Appendix C records it for A4/attended review.
- **P4 — Design §6 XA4 "the next refresh → none" presumes a refresh can be awaited.** Through the driver a
  Promise is not awaited (C13), and `refresh` drops a `poll` while one runs (C8). Routed in the test only:
  `pollSurface` (S4-D12).
- **P5 — Design §2.8's line numbers** (`app.js:605-606`, measured at `0140bd7f`) are now `app.js:616` (sprint 3
  moved them, C4); `surface.js:257` is unchanged (C6).
- **P6 — Design §6 "the driver gains … `grantNotifications` … and the `notify` … injection" and design V26
  "the constructor succeeds".** Built as named; the recorder does not construct a real notification (S4-D8),
  so this sprint does not re-measure the constructor (V26 did, at design time); it re-measures the grant, the
  gesture-bound request and the click path live (C12).
- **P7 — `World.open` is synchronous in its answer only** (C14): a test waits for the brane view.

### 1.5 The size

Lines are insertions + deletions under `cmd/ host/`: **`git diff --numstat <base> HEAD -- cmd host`**
(`docs/` and `design_docs/` are outside the filter; `.js`/`.mjs` under `host/`/`cmd/` count). **Every row is
measured** from the scratch stages, each committed in a throwaway repo (`plan-scratch/work`, tags `M9a` …
`M10c`; `tools/numstat.sh`).

| Sub-milestone | Production | Test | Cumulative vs `bed9d135` (after its commit) | Run |
|---|---|---|---|---|
| M9a — the A02/A03/A05 kill test (`TestSurfaceAppPure`'s hostile rows) | 0 | 6 | **6** | 1 |
| M9b — Decided (`decidedHTML`, the toggle); node rows | 19 | 11 | **36** | 1 |
| M9c — the opt-in and the notification (`notifyHTML`, `notifications`, `loadRequests`, `askNotify`); node rows; XA3's one-line expectation | 39 | 22 | **93** | 1 |
| M10a — `N`, Enter in the now card, the keys row; the driver's `key`/`tap`; XA3 extended (Decided half, Show all, `N`) | 4 | 144 | **239** | 2 |
| M10b — XA4 and the driver's `notify`/`grantNotifications`/`notification` | 0 | 120 | **357** | 2 |
| M10c — the card's hostile value in Chrome; QUICKSTART §17 (+4 docs lines, not counted) | 0 | 32 | **389** | 2 |
| **Sprint 4** (cumulative, `git diff --numstat bed9d135 M10c -- cmd host`) | **58** | **331** | **389 ≤ 2,000** (margin 1,611) | |

The per-stage columns sum to 62 + 335 = 397: M9c and M10a re-edit lines M9b and M9c added (`needsYouHTML`'s
signature, the export line, the driver's `commands` line), which the cumulative diff counts once. Per file,
cumulative: `app.js` +49 −9, `app_test.go` +35 −2, `passkey_cdp.mjs` +42 −1, `needs_you_chrome_test.go` +103 −23,
`needs_you_notify_chrome_test.go` +93, `approve_page_hostile_test.go` +32; `docs/QUICKSTART.md` +3 −1. The
sprint-3 plan's estimate for the moved part was ≈ 300 (≈ 135 `app.js`, XA4 ≈ 115, the driver ≈ 30, XA3's
Decided half ≈ 20); measured, the moved part is ≈ 351 (all but the A02 test's 6 and the hostile test's 32).

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of
four backticks, in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; a new file's patch has an
`@@ -0,0 …` hunk. Apply it **only** this way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF'
set -euo pipefail
M=M9a   # <- the sub-milestone, exactly as in the Appendix P heading: M9a M9b M9c M10a M10b M10c
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter387-a3a-s4
PLAN="$W/design_docs/planned/w-world-approvals-by-key-a3a-sprint4-plan.md"
T="$HOME/.ailang/state/world-iter387/exec-patches"; mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
patch -p1 --dry-run < "$T/$M.patch"
patch -p1 < "$T/$M.patch"
if find cmd host docs \( -name '*.orig' -o -name '*.rej' \) | grep .; then echo "STOP: patch left .orig/.rej"; exit 1; fi
git status --short
EOF
```

- If `--dry-run` fails, **stop and report** (no hand edits, no fuzz options). It can only fail if an earlier
  sub-milestone was skipped or changed.
- Then run that sub-milestone's **Acceptance** (§3) in one bash shell with `export
  AILANG_BIN=$HOME/.pinned-ailang/ailang` and **`unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER
  WORLD_PASSKEY_RECORD_DIR`** (the Chrome legs set the two variables on their own command line), each
  command's output `tee`'d to a log under `$HOME/.ailang/state/world-iter387/exec-logs/` (never only through
  `tail`). **Never `nohup`, never background a test leg.**
- Then snapshot (§7). The controller commits and compares the numstat with §7's expected value.

---

## 3. Milestones

Every sub-milestone's acceptance starts with the same three lines (abbreviated **[base]** below):

```
test -z "$(gofmt -l cmd host)"
go build ./... && go vet ./host/surfacepage/ ./host/daemon/ ./cmd/ailang-worldd/
go test -run '^$' ./...
```

Each `-v` run below must show `--- PASS` for every test it names and **no `--- FAIL` and no `--- SKIP`**;
each package run must end `ok` with no `FAIL` line. **[node]** means the command is prefixed with
`WORLD_EXEC_NODE=$(command -v node)` (so `TestSurfaceAppPure` cannot skip). **[chrome]** means the command is
prefixed with `WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
WORLD_REQUIRE_PASSKEY_BROWSER=1` and gets `-timeout 600s` (headless, temporary profiles, a virtual
authenticator, test keys). The planner's `tools/accept.sh <tree> <M> <logs>` is exactly these lines (§1.2
C26).

### M9a — the judge's survivors are pinned · 6 test · Run 1

**Patch** `M9a`: `host/surfacepage/app_test.go` (`TestSurfaceAppPure`'s harness): one `assert.equal` of
`needsYouHTML` over a pending row whose id is `mr-"><img src=x onerror=1>`, agent `<img src=x onerror=2>` and
recorded approver `"><script>3</script>`, and a refused row whose id is `mr-<b>"4</b>` and reason
`<i>5</i>` — the whole section, every value escaped where it lands (the link's `id=` and `href`, the text); and
the `unavailable` text `<i>"down"</i>` escaped. **No production change.**

**Acceptance** (seconds)
```
[base]
[node] go test ./host/surfacepage/ -run '^TestSurfaceAppPure$' -count=1 -v
[node] go test ./host/surfacepage/ -count=1
```

**Kill proof** (the directive's item 4; §4 rows `ESC-*`, ledger runs `base-esc` and `m9a-esc`): each of the
eight `esc()` calls in `needsYouHTML`, removed alone: at the base 6 SURVIVE (`agent`, `id=`, `href`,
`recorded`, the refused id, `unavailable`) and 2 are KILLED (`episode`, the refused reason); on the M9a tree
**all 8 are KILLED** by `TestSurfaceAppPure` (named and the whole `host/surfacepage` suite).

**Expected cumulative numstat after the controller's commit of M9a: 6**.

### M9b — the decision list · 19 prod + 11 test · Run 1

**Patch** `M9b`: `host/surfacepage/app/app.js`: `decidedHTML(r, all)` (pure, exported, S4-D5), the now card's
`${decidedHTML(REQS, decidedAll)}` under Needs you, `decidedAll`, and the stage click's `#decided-all` toggle
(repaints the now card). `app_test.go`: the export list +`decidedHTML`; rows: `null` and `[]` → nothing; four
decided (hostile episode and approver) → three escaped lines with their UTC times and **Show all 4**
(`aria-expanded="false"`); all → four lines and **Show fewer** (`"true"`); three → no button.

**Acceptance** (~1 min)
```
[base]
[node] go test ./host/surfacepage/ -run '^(TestSurfaceAppPure|TestSurfaceAppMatchesWorkingCopy|TestSurfaceAppShellAssembly)$' -count=1 -v
[node] go test ./host/surfacepage/ -count=1
go test ./host/daemon/ -run '^(TestSurfaceAppRoutes|TestApprovePageRoutes)$' -count=1 -v
[chrome] go test ./cmd/ailang-worldd/ -run '^TestNeedsYouListsPendingRequests$' -count=1 -v
```

**Expected cumulative numstat after the controller's commit of M9b: 36**.

### M9c — the opt-in and the notification · 39 prod + 22 test · Run 1

**Patch** `M9c`: `host/surfacepage/app/app.js`: `needsYouHTML(r, now, permission)` + `notifyHTML` (S4-D3);
`notifications(r, seen, permission)` (pure, exported, S4-D2); `seenIds`, `notifyPermission()`, `loadRequests`
raising each note with its `onclick`; `askNotify()` and the stage click's `#notify-me`; the now card passes
`notifyPermission()`. `app_test.go`: the export list +`notifications`; rows: the opt-in for `default` /
`granted` / `denied` (and none for the two-argument calls); `notifications` — a failed read changes nothing,
the first read raises nothing and seeds `seen`, a new id raises exactly one note (title, tag, body with
"1 grant" and the UTC expiry, href), a refused and a recorded new id raise nothing, the same read again raises
nothing, under `default` nothing is raised but `seen` still grows, "2 grants". `cmd/ailang-worldd/
needs_you_chrome_test.go`: XA3's expected Needs-you text gains the opt-in (`… expires in 10 minNotify me`), so
the live test stays green at this commit.

**Acceptance** (~1 min)
```
[base]
[node] go test ./host/surfacepage/ -run '^(TestSurfaceAppPure|TestSurfaceAppMatchesWorkingCopy|TestSurfaceAppShellAssembly)$' -count=1 -v
[node] go test ./host/surfacepage/ -count=1
go test ./host/daemon/ -run '^(TestSurfaceAppRoutes|TestApprovePageRoutes)$' -count=1 -v
[chrome] go test ./cmd/ailang-worldd/ -run '^TestNeedsYouListsPendingRequests$' -count=1 -v
```

**Expected cumulative numstat after the controller's commit of M9c: 93**.

### M10a — `N`, Enter, and XA3's Decided half · 4 prod + 144 test · Run 2

**Patch** `M10a`: `host/surfacepage/app/app.js`: the `n`/`N` key and the now card's Enter (S4-D7), the keys
row. `cmd/ailang-worldd/testdata/passkey_cdp.mjs`: commands `key` and `tap` (S4-D11).
`cmd/ailang-worldd/needs_you_chrome_test.go`: the helpers (S4-D12) and **XA3** extended (S4-D9): "Needs you
· 1ep1 · claude · 10 grants · expires in 10 minNotify me"; `n` focuses `needs-<id>` on the timeline; in the
brane of episode `e2`, `N` (Shift) returns to the timeline and focuses it; a second request appears on the
app's refresh ("Needs you · 2"); Enter on the focused link opens `/surface/approve/<id>`; Approve there
(the test passkey) records `webauthn:mark-mac`; `approve --await` mints it; `/surface` then shows "DecidedYou
· Minted · ep1 · webauthn:mark-mac · <hh:mm>" and "Needs you · 1ep2 · claude …"; with three more decided,
3 lines and **Show all 4** → 4 and **Show fewer** → 3 and **Show all 4**.

**Acceptance** (~1 min)
```
[base]
[chrome] go test ./cmd/ailang-worldd/ -run '^(TestNeedsYouListsPendingRequests|TestApprovePageInChrome|TestPasskeyPageInChrome|TestPasskeyPageAbortsAPendingCeremony)$' -count=1 -v
```

**Expected cumulative numstat after the controller's commit of M10a: 239**.

### M10b — XA4, the opt-in notification · 120 test · Run 2

**Patch** `M10b`: `cmd/ailang-worldd/testdata/passkey_cdp.mjs`: the `notify` injection (S4-D8) and the
commands `grantNotifications`, `notification`. New `cmd/ailang-worldd/needs_you_notify_chrome_test.go`:
**XA4** `TestNeedsYouNotifiesOnlyWhenAskedOnce` — never asked: 0 notifications, 0 permission requests,
`default`, "Notify me", also after a new request (`ep2`); `grantNotifications` → `granted`; tap **Notify me**
→ exactly one request, made under a user gesture, and "Notifications on"; a refresh → none (both ids seen); a
new request `ep3` → exactly one `{World needs you, "Approve mint request for ep3 (10 grants), expires <its UTC
hh:mm>", tag <ep3's id>}`; the next refresh → still one; reloaded with the grant → 0 notifications and 0
requests for the three it shows; a new `ep4` → one, tag `ep4`'s id; its click → the tab focused once and at
`/surface/approve/<ep4's id>`, the card ready; zero CSP violations at every step.

**Acceptance** (~1 min)
```
[base]
[chrome] go test ./cmd/ailang-worldd/ -run '^(TestNeedsYouNotifiesOnlyWhenAskedOnce|TestNeedsYouListsPendingRequests|TestPasskeyPageInChrome)$' -count=1 -v
```

**Expected cumulative numstat after the controller's commit of M10b: 357**.

### M10c — the card shows a hostile value as text; QUICKSTART §17 · 32 test + 4 docs · Run 2

**Patch** `M10c`: new `cmd/ailang-worldd/approve_page_hostile_test.go`: `TestApprovePageShowsAHostileValueAsText`
(S4-D10). `docs/QUICKSTART.md` §17 (still *attended — pending first verbatim run*): after the mint, Decided
and its line; **Show all**; `N` and Enter; a paragraph on **Notify me** (asked only from the click; one
notification per new request the tab reads on its own refresh; the click opens the card; no authority; display
is the browser's and the system's, checked once per browser in the exit run).

**Acceptance** (~3–4 min: the whole `cmd/ailang-worldd`, the XA tests there as counted skips)
```
[base]
[chrome] go test ./cmd/ailang-worldd/ -run '^(TestApprovePageShowsAHostileValueAsText|TestApprovePageInChrome)$' -count=1 -v
go test ./cmd/ailang-worldd/ -run '^(TestQuickstartAttendedStepsBoxMatchesSection9|TestQuickstartSection9FlagsMatchTheCLI|TestQuickstartSection10FlagsMatchTheCLI|TestQuickstartSection11FlagsMatchTheCLI|TestQuickstartServiceSectionMatchesTheCLIHelp)$' -count=1 -v
go test -count=1 -timeout 20m ./cmd/ailang-worldd/
```

**Expected cumulative numstat after the controller's commit of M10c: 389** (and `-- docs` **4**).

---

## 4. Mutation table

**60 rows on the planner's final tree: 59 KILLED (named and suite), 1 named equivalent(s) SURVIVED as predicted, 0 not as expected.** 77 ledger lines in all (`plan-scratch/logs/mutant_ledger.tsv`; runs: base-card, base-esc, m9a-esc, run1-chrome, run1-node).

**The kill proofs** (directive item 4, judge A02/A03/A05 and C03/C04): the same edit at the base, on the M9a tree (the six test lines alone) and on the final tree. At the base the card row's named test is `TestApprovePageInChrome` (the hostile-value test does not exist there).

| Id | At the base (`bed9d135`) named / suite | At M9a named / suite | Final tree named / suite |
|---|---|---|---|
| ESC-AGENT | SURVIVED / SURVIVED | KILLED / KILLED | KILLED / KILLED |
| ESC-ID-ATTR | SURVIVED / SURVIVED | KILLED / KILLED | KILLED / KILLED |
| ESC-ID-HREF | SURVIVED / SURVIVED | KILLED / KILLED | KILLED / KILLED |
| ESC-RECORDED | SURVIVED / SURVIVED | KILLED / KILLED | KILLED / KILLED |
| ESC-EPISODE | KILLED / KILLED | KILLED / KILLED | KILLED / KILLED |
| ESC-REFUSED-ID | SURVIVED / SURVIVED | KILLED / KILLED | KILLED / KILLED |
| ESC-REFUSED | KILLED / KILLED | KILLED / KILLED | KILLED / KILLED |
| ESC-UNAVAILABLE | SURVIVED / SURVIVED | KILLED / KILLED | KILLED / KILLED |
| CARD-TEXTCONTENT | SURVIVED / SURVIVED (`TestApprovePageInChrome`) | — | KILLED / KILLED |

**Derivation.** Anchored to Appendix P's production hunks, one row per new branch, enumerated from the diff
(`git diff bed9d135 M10c -- host/surfacepage/app/app.js`, the only production file): `decidedHTML` (its input
guard, the empty case, the few/all slice, both `esc()`s, the time, the button's condition, label and
`aria-expanded`) (M9b, node); the now card's Decided call, the toggle's wiring, its negation and its repaint
(M9b, Chrome); `needsYouHTML`'s opt-in on both of its branches and `notifyHTML`'s three arms; `notifications`'
input guard, the first-read rule, the seen check, the refused and recorded exclusions, the seen update, the
granted check, the title, tag, body (episode, number, time) and href (M9c, node); `loadRequests`' permission
argument, the seen assignment, the raise, the tag passed to the constructor, the click's focus and its target,
`notifyPermission`'s read and its no-API guard, the opt-in's click wiring, asking only from the click, the
repaint after the answer, and the now card's permission argument (M9c, Chrome); the `N` branch, each of its two
keys, its target, its leaving a brane, its focus, and the now card's Enter (M10a, Chrome). Plus the eight
existing `esc()` calls of `needsYouHTML` that M9a pins (the directive's item 4), run at the base, at M9a and on
the final tree, and the card's `textContent` that M10c pins, run at the base and on the final tree.
**Procedure**: `plan-scratch/tools/mutate.py <run> --workers N` (the table is `tools/mutants.py`; every anchor
asserted to occur exactly once in the tree it edits): for each row, in its own worker copy (`SRC` = the tree:
`work` (= M10c) unless stated), apply **one** edit, run **the single named test** and then **its package suite**
with the inventories excluded — `host/surfacepage` whole with `WORLD_EXEC_NODE` set; `cmd/ailang-worldd`
`^(TestApprovePage|TestApproveCard|TestNeedsYou)` in live headless Chrome with `WORLD_CHROME=…
WORLD_REQUIRE_PASSKEY_BROWSER=1` — with `-vet=off` (a mutant that does not compile or parse is `NOCOMPILE`,
never counted), restore, assert the restore byte-identical, and append one ledger line (UTC time, run, tree,
id, stage, file, from → to, named test, named, suite, the suite's failing tests, seconds). A row is KILLED
only when both fail. **No kill in this table rests on an inventory or static test**: every suite excludes the
§7 inventories, and the failing-test column shows each kill is the named behavioural test (no row is killed
only by `TestSurfaceAppShellAssembly`/`TestSurfaceAppMatchesWorkingCopy`, the package's static tests).
**Never commit or snapshot a mutant.** In the table a newline prints as ⏎, a backquote as ' and a pipe as \|.

**Named equivalents** (survive by construction):
- `NP-NOAPI` — the `typeof Notification === 'undefined'` arm is for a browser without the Notifications API; Chrome (the only browser the tests run) and the notify recorder both define it, and node never calls `notifyPermission` (it is page code), so removing the guard changes nothing a test can observe; it stays so a browser without the API renders no opt-in instead of throwing

| Id | Stage | Mutant (single edit) | Named test | Named | Suite | Suite failures (inventories excluded) |
|---|---|---|---|---|---|---|
| ESC-AGENT | M9a | app.js: · ${esc(q.agent)} · → · ${q.agent} · | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| ESC-ID-ATTR | M9a | app.js: id="needs-${esc(q.id)}" → id="needs-${q.id}" | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| ESC-ID-HREF | M9a | app.js: href="/surface/approve/${esc(q.id)}" → href="/surface/approve/${q.id}" | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| ESC-RECORDED | M9a | app.js: approved by ${esc(q.recorded)} → approved by ${q.recorded} | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| ESC-EPISODE | M9a | app.js: <b>${esc(q.episode)}</b> → <b>${q.episode}</b> | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| ESC-REFUSED-ID | M9a | app.js: <div class="nl sub">${esc(q.id)} · refused: → <div class="nl sub">${q.id} · refused: | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| ESC-REFUSED | M9a | app.js: refused: ${esc(q.refused)} → refused: ${q.refused} | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| ESC-UNAVAILABLE | M9a | app.js: ${esc(r && r.unavailable \|\| 'Nothing needs you')} → ${r && r.unavailable \|\| 'Nothing needs you'} | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| CARD-TEXTCONTENT | M10c | approve.js: dd.textContent = shown(f, l.value); → dd.innerHTML = shown(f, l.value); | `TestApprovePageShowsAHostileValueAsText` | KILLED | KILLED | TestApprovePageShowsAHostileValueAsText |
| DEC-GUARD | M9b | app.js: const d = r && Array.isArray(r.decided) ? r.decided : []; → const d = r.decided; | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| DEC-EMPTY | M9b | app.js:     if (!d.length) return '';⏎    return '<div class="nsec" data-decided> →     return '<div class="nsec" data-decided> | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| DEC-FEW | M9b | app.js: (all ? d : d.slice(0, DECIDED_FEW)) → (all ? d : d) | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| DEC-ALL | M9b | app.js: (all ? d : d.slice(0, DECIDED_FEW)) → (false ? d : d.slice(0, DECIDED_FEW)) | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| DEC-ESC-EPISODE | M9b | app.js: You · Minted · ${esc(x.episode)} · → You · Minted · ${x.episode} · | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| DEC-ESC-APPROVER | M9b | app.js: · ${esc(x.approvedBy)} · → · ${x.approvedBy} · | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| DEC-TIME | M9b | app.js: ${hm(new Date(x.at * 1000))} → ${hm(new Date(x.at))} | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| DEC-BUTTON | M9b | app.js: (d.length > DECIDED_FEW ? → (d.length >= DECIDED_FEW ? | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| DEC-LABEL | M9b | app.js: ${all ? 'Show fewer' : 'Show all ' + d.length} → ${'Show all ' + d.length} | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| DEC-ARIA | M9b | app.js: aria-expanded="${!!all}" → aria-expanded="false" | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| DEC-NOWCARD | M9b | app.js: ${decidedHTML(REQS, decidedAll)} → ${decidedHTML(null, decidedAll)} | `TestNeedsYouListsPendingRequests` | KILLED | KILLED | TestNeedsYouListsPendingRequests |
| DEC-WIRE | M9b | app.js: if (e.target.closest('#decided-all')) { → if (false) { | `TestNeedsYouListsPendingRequests` | KILLED | KILLED | TestNeedsYouListsPendingRequests |
| DEC-TOGGLE | M9b | app.js: decidedAll = !decidedAll; → decidedAll = true; | `TestNeedsYouListsPendingRequests` | KILLED | KILLED | TestNeedsYouListsPendingRequests |
| DEC-REPAINT | M9b | app.js: decidedAll = !decidedAll; return repaintNow(); } → decidedAll = !decidedAll; return; } | `TestNeedsYouListsPendingRequests` | KILLED | KILLED | TestNeedsYouListsPendingRequests |
| NY-NOTIFY-EMPTY | M9c | app.js: Nothing needs you')}</div>${notifyHTML(permission)}</div>'; → Nothing needs you')}</div></div>'; | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| NY-NOTIFY-LIST | M9c | app.js: .join('') + notifyHTML(permission) + '</div>'; → .join('') + '</div>'; | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| NH-DEFAULT | M9c | app.js: if (permission === 'default') return → if (false) return | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| NH-GRANTED | M9c | app.js: if (permission === 'granted') return → if (false) return | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| NH-DENIED | M9c | app.js: return permission === 'denied' ? → return false ? | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| NT-GUARD | M9c | app.js: if (!r \|\| !Array.isArray(r.pending)) return { notes: [], seen }; → if (!r) return { notes: [], seen }; | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| NT-FIRST-READ | M9c | app.js: const fresh = seen ? r.pending.filter(q => !seen.has(q.id) && !q.refused && !q.recorded) : []; → const fresh = r.pending.filter(q => !(seen && seen.has(q.id)) && !q.refused && !q.recorded); | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| NT-SEEN | M9c | app.js: r.pending.filter(q => !seen.has(q.id) && !q.refused → r.pending.filter(q => !q.refused | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| NT-REFUSED | M9c | app.js: !seen.has(q.id) && !q.refused && !q.recorded) → !seen.has(q.id) && !q.recorded) | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| NT-RECORDED | M9c | app.js: !seen.has(q.id) && !q.refused && !q.recorded) → !seen.has(q.id) && !q.refused) | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| NT-SEEN-UPDATE | M9c | app.js: return { seen: new Set([...(seen \|\| []), ...r.pending.map(q => q.id)]), → return { seen: seen \|\| new Set(), | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| NT-GRANTED | M9c | app.js: notes: permission !== 'granted' ? [] : → notes: false ? [] : | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| NT-TITLE | M9c | app.js: title: 'World needs you' → title: 'World' | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| NT-TAG | M9c | app.js: tag: q.id, → tag: q.episode, | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| NT-BODY-EPISODE | M9c | app.js: Approve mint request for ${q.episode} ( → Approve mint request for ${q.id} ( | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| NT-BODY-PLURAL | M9c | app.js: (${q.grants} grant${q.grants === 1 ? '' : 's'}), expires → (${q.grants} grants), expires | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| NT-BODY-TIME | M9c | app.js: expires ${hm(new Date(q.expires * 1000))} → expires ${hm(new Date(q.expires))} | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| NT-HREF | M9c | app.js: href: '/surface/approve/' + q.id → href: '/surface/approve/' + q.episode | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| LR-PERMISSION | M9c | app.js: notifications(REQS, seenIds, notifyPermission()) → notifications(REQS, seenIds, 'granted') | `TestNeedsYouNotifiesOnlyWhenAskedOnce` | KILLED | KILLED | TestNeedsYouNotifiesOnlyWhenAskedOnce |
| LR-SEEN | M9c | app.js:     seenIds = n.seen;⏎ →  | `TestNeedsYouNotifiesOnlyWhenAskedOnce` | KILLED | KILLED | TestNeedsYouNotifiesOnlyWhenAskedOnce |
| LR-RAISE | M9c | app.js: n.notes.forEach(x => { → [].forEach(x => { | `TestNeedsYouNotifiesOnlyWhenAskedOnce` | KILLED | KILLED | TestNeedsYouNotifiesOnlyWhenAskedOnce |
| LR-TAG | M9c | app.js: new Notification(x.title, { body: x.body, tag: x.tag }) → new Notification(x.title, { body: x.body }) | `TestNeedsYouNotifiesOnlyWhenAskedOnce` | KILLED | KILLED | TestNeedsYouNotifiesOnlyWhenAskedOnce |
| LR-CLICK-FOCUS | M9c | app.js: note.onclick = () => { window.focus(); location.assign(x.href); }; → note.onclick = () => { location.assign(x.href); }; | `TestNeedsYouNotifiesOnlyWhenAskedOnce` | KILLED | KILLED | TestNeedsYouNotifiesOnlyWhenAskedOnce |
| LR-CLICK-OPEN | M9c | app.js: note.onclick = () => { window.focus(); location.assign(x.href); }; → note.onclick = () => { window.focus(); location.assign('/surface'); }; | `TestNeedsYouNotifiesOnlyWhenAskedOnce` | KILLED | KILLED | TestNeedsYouNotifiesOnlyWhenAskedOnce |
| NP-PERMISSION | M9c | app.js: const notifyPermission = () => typeof Notification === 'undefined' ? '' : Notification.permission; → const notifyPermission = () => 'default'; | `TestNeedsYouNotifiesOnlyWhenAskedOnce` | KILLED | KILLED | TestNeedsYouNotifiesOnlyWhenAskedOnce |
| NP-NOAPI | M9c | app.js: const notifyPermission = () => typeof Notification === 'undefined' ? '' : Notification.permission; → const notifyPermission = () => Notification.permission; | `TestNeedsYouNotifiesOnlyWhenAskedOnce` | SURVIVED | SURVIVED |  |
| ASK-WIRE | M9c | app.js: if (e.target.closest('#notify-me')) return askNotify(); → if (false) return askNotify(); | `TestNeedsYouNotifiesOnlyWhenAskedOnce` | KILLED | KILLED | TestNeedsYouNotifiesOnlyWhenAskedOnce |
| ASK-UNASKED | M9c | app.js:     document.documentElement.dataset.ready = '1';⏎ →     askNotify(); document.documentElement.dataset.ready = '1';⏎ | `TestNeedsYouNotifiesOnlyWhenAskedOnce` | KILLED | KILLED | TestNeedsYouListsPendingRequests,TestNeedsYouNotifiesOnlyWhenAskedOnce |
| ASK-REPAINT | M9c | app.js: .then(repaintNow, repaintNow); → .then(() => {}, () => {}); | `TestNeedsYouNotifiesOnlyWhenAskedOnce` | KILLED | KILLED | TestNeedsYouNotifiesOnlyWhenAskedOnce |
| NOWCARD-PERMISSION | M9c | app.js: ${needsYouHTML(REQS, REQS ? REQS.now : 0, notifyPermission())} → ${needsYouHTML(REQS, REQS ? REQS.now : 0)} | `TestNeedsYouListsPendingRequests` | KILLED | KILLED | TestNeedsYouListsPendingRequests,TestNeedsYouNotifiesOnlyWhenAskedOnce |
| NKEY | M10a | app.js: else if (k === 'n' \|\| k === 'N') { → else if (false) { | `TestNeedsYouListsPendingRequests` | KILLED | KILLED | TestNeedsYouListsPendingRequests |
| NKEY-LOWER | M10a | app.js: else if (k === 'n' \|\| k === 'N') { → else if (k === 'N') { | `TestNeedsYouListsPendingRequests` | KILLED | KILLED | TestNeedsYouListsPendingRequests |
| NKEY-UPPER | M10a | app.js: else if (k === 'n' \|\| k === 'N') { → else if (k === 'n') { | `TestNeedsYouListsPendingRequests` | KILLED | KILLED | TestNeedsYouListsPendingRequests |
| NKEY-TARGET | M10a | app.js: const a = document.querySelector('[data-needs] a'); → const a = document.querySelector('a'); | `TestNeedsYouListsPendingRequests` | KILLED | KILLED | TestNeedsYouListsPendingRequests |
| NKEY-OUT | M10a | app.js: if (a) { if (st.d === 1) out(); a.focus(); } → if (a) { a.focus(); } | `TestNeedsYouListsPendingRequests` | KILLED | KILLED | TestNeedsYouListsPendingRequests |
| NKEY-FOCUS | M10a | app.js: if (a) { if (st.d === 1) out(); a.focus(); } → if (a) { if (st.d === 1) out(); } | `TestNeedsYouListsPendingRequests` | KILLED | KILLED | TestNeedsYouListsPendingRequests |
| ENTER-NOW | M10a | app.js: if (k === 'Enter' && e.target && e.target.closest && e.target.closest('[data-now]')) return; → if (false) return; | `TestNeedsYouListsPendingRequests` | KILLED | KILLED | TestNeedsYouListsPendingRequests |

---

## 5. Files (final `git status --short` after M10c)

| File | Sub-milestone | Change |
|---|---|---|
| `host/surfacepage/app_test.go` | M9a, M9b, M9c | the hostile rows; Decided's rows; the opt-in's and `notifications`' rows |
| `host/surfacepage/app/app.js` | M9b, M9c, M10a | Decided; the opt-in and the notification; `N` and the now card's Enter |
| `cmd/ailang-worldd/needs_you_chrome_test.go` | M9c, M10a | XA3's expectation; XA3 extended, the helpers |
| `cmd/ailang-worldd/testdata/passkey_cdp.mjs` | M10a, M10b | `key`, `tap`; `notify`, `grantNotifications`, `notification` |
| `cmd/ailang-worldd/needs_you_notify_chrome_test.go` | M10b | new: XA4 |
| `cmd/ailang-worldd/approve_page_hostile_test.go`, `docs/QUICKSTART.md` | M10c | new: the hostile value; §17 |
| this plan and `sprint_w-world-approvals-by-key-a3a-s4.json` | planner | |

---

## 6. Sprint gates (controller, after the last sub-milestone; each measured at the base, §1.1)

In **one** `bash` shell, from the worktree, **not under `nohup`**, nothing else running. Every leg `tee`s its
full output to a log under `$HOME/.ailang/state/world-iter387/gates/`; a verdict is read from its log, never
from `tail` alone. The live-Chrome leg runs on this machine only (CI has no Chrome step, C18).

```
bash <<'EOF'
set -euo pipefail
export AILANG_BIN=$HOME/.pinned-ailang/ailang
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
export WORLD_EXEC_NODE=$(command -v node)
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter387-a3a-s4; cd "$W"
L=$HOME/.ailang/state/world-iter387/gates; mkdir -p "$L"; MODE=final
s0=$(date +%s)
./scripts/verify_ail.sh 2>&1 | tee "$L/verify_ail.log"
out=$(gofmt -l cmd host); echo "$out" | tee "$L/gofmt.log"; test -z "$out"
go build ./... 2>&1 | tee "$L/build.log"
go vet ./... 2>&1 | tee "$L/vet.log"
go test -run '^$' ./... 2>&1 | tee "$L/fence.log"
if grep -E '^(FAIL|--- FAIL)' "$L/fence.log"; then exit 1; fi
go test -count=1 -timeout 20m ./host/surfacepage/ ./host/store/ ./host/authority/ ./host/approval/... ./cmd/ailang-worldd/ ./host/daemon/ 2>&1 | tee "$L/pkgs.log"
if grep -E '^(FAIL|--- FAIL)' "$L/pkgs.log"; then exit 1; fi
go test -race -count=1 ./host/surfacepage/ 2>&1 | tee "$L/race.log"
go test -race -count=1 ./host/daemon/ -run '^(TestApprov|TestRequest|TestDaemonCredentialWriterSites|TestSurfaceAppRoutes)' 2>&1 | tee -a "$L/race.log"
go test -race -count=1 ./cmd/ailang-worldd/ -run '^(TestApprove|TestServeSurfaceOrigin|TestSurfaceApprove|TestSurfaceSSHSIG|TestEveryVerbHasHelp)' 2>&1 | tee -a "$L/race.log"
if grep -E '^(FAIL|--- FAIL)|WARNING: DATA RACE' "$L/race.log"; then exit 1; fi
NEW_SP="TestApprovePageRendersEveryCanonicalField TestApprovePageFieldTagControls"
NEW_D="TestRequestCardPasskeyReadError TestApprovePOSTRecordsAVerifiedSSHSIG TestApprovePOSTSSHSIGViaTailnetOriginWithNoPasskey TestApprovePageRoutes"
NEW_C="TestSurfaceSSHSIGThenAwaitMints"
alt() { local x="$*"; x="${x// /|}"; echo "${x:+|$x}"; }
go test -v -count=1 ./host/broker/ -run '^(TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential|TestQuickstartSection10ProfilesLoad|TestQuickstartProtectedContract)$' 2>&1 | tee "$L/inv.log"
go test -v -count=1 ./host/daemon/ -run "^(TestDaemonDependencyAllowlist|TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestDaemonConstructorCallerInventory|TestQuickstartDecisionPacketSection|TestPasskeyPageRoutes|TestReadCtxCancelledAfterHandler|TestMCPMountMethodSource|TestDaemonCredentialWriterSites|TestApprovalRoutesRegistered|TestApprovalRoutesRefuseForeignHost|TestRequestsJSON|TestRequestCardJSON|TestApprovePOSTRefusalArms|TestApprovePOSTRecordsAVerifiedPasskey|TestApprovePOSTCrossOrigin|TestApprovalGuards$(alt $NEW_D))\$" 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/boundary/ -run '^(TestWorldBoundaryDependencyAllowlist|TestBareNetHTTPExemptionIsPerGroup|TestEveryCommandDirectoryIsAProtectedGroup)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/store/ -run '^(TestProductionContextRoots|TestProductionGoSurface|TestNoNewDeadlineFreeStoreReads|TestReadGettersHonorContext|TestCredentialsMethodsRequireADeadline|TestCredentialsSchemaMatchesCanonicalDDL|TestCredentialsSchema7SurfaceApprovals|TestRecordSurfaceApprovalGuards|TestSurfaceApprovalPrunedWithItsRequest)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/verifygate/ -run '^(TestNoParallelWriteForkPackagesOutsideVerifygate|TestSubprocessSinkGate)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/approval/... -run '^(TestApprovalImportsOnlyTheStandardLibrary|TestApprovalTestIsTestOnly)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/authority/ -run '^(TestJudgeAssertionIsTheCLIJudge|TestJudgeSSHSignatureIsTheCLIJudge)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./cmd/ailang-worldd/ -run "^(TestEveryVerbHasHelp|TestServeProofFlags|TestQuickstartAttendedStepsBoxMatchesSection9|TestQuickstartSection9FlagsMatchTheCLI|TestQuickstartSection10FlagsMatchTheCLI|TestQuickstartSection11FlagsMatchTheCLI|TestQuickstartServiceSectionMatchesTheCLIHelp|TestSurfaceApproveThenAwaitMints|TestServeSurfaceOriginFlag$(alt $NEW_C))\$" 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/surfacepage/ -run "^(TestSurfaceAppPure|TestApproveModelSummaries|TestApproveModelCard$(alt $NEW_SP))\$" 2>&1 | tee -a "$L/inv.log"
for t in TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential TestQuickstartSection10ProfilesLoad TestQuickstartProtectedContract \
  TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestDaemonConstructorCallerInventory TestQuickstartDecisionPacketSection TestPasskeyPageRoutes TestReadCtxCancelledAfterHandler TestMCPMountMethodSource TestDaemonCredentialWriterSites \
  TestApprovalRoutesRegistered TestApprovalRoutesRefuseForeignHost TestRequestsJSON TestRequestCardJSON TestApprovePOSTRefusalArms TestApprovePOSTRecordsAVerifiedPasskey TestApprovePOSTCrossOrigin TestApprovalGuards \
  TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup \
  TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext TestCredentialsMethodsRequireADeadline TestCredentialsSchemaMatchesCanonicalDDL \
  TestCredentialsSchema7SurfaceApprovals TestRecordSurfaceApprovalGuards TestSurfaceApprovalPrunedWithItsRequest \
  TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestApprovalImportsOnlyTheStandardLibrary TestApprovalTestIsTestOnly \
  TestJudgeAssertionIsTheCLIJudge TestJudgeSSHSignatureIsTheCLIJudge \
  TestEveryVerbHasHelp TestServeProofFlags TestQuickstartAttendedStepsBoxMatchesSection9 TestQuickstartSection9FlagsMatchTheCLI TestQuickstartSection10FlagsMatchTheCLI \
  TestQuickstartSection11FlagsMatchTheCLI TestQuickstartServiceSectionMatchesTheCLIHelp TestSurfaceApproveThenAwaitMints TestServeSurfaceOriginFlag TestSurfaceAppPure TestApproveModelSummaries TestApproveModelCard $NEW_SP $NEW_D $NEW_C; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || { echo "$t did not PASS"; exit 1; }; done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then exit 1; fi
go test -v -count=1 ./cmd/ailang-worldd/ -run '^TestApproveAwait' 2>&1 | tee "$L/aw.log"
test "$(grep -cE '^--- PASS: TestApproveAwait' "$L/aw.log")" = 9
if grep -E '^\s*--- (SKIP|FAIL)' "$L/aw.log"; then exit 1; fi
# the live-Chrome leg (this machine only: CI has no Chrome step)
CHROME_RUN='^(TestPasskeyPageInChrome|TestPasskeyPageAbortsAPendingCeremony|TestApprovePage|TestApproveCard|TestNeedsYou)'
WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" WORLD_REQUIRE_PASSKEY_BROWSER=1 go test -v -count=1 -timeout 600s ./cmd/ailang-worldd/ -run "$CHROME_RUN" 2>&1 | tee "$L/chrome.log"
if grep -E '^\s*--- (SKIP|FAIL)|^FAIL' "$L/chrome.log"; then exit 1; fi
WANT_CHROME="TestPasskeyPageInChrome TestPasskeyPageAbortsAPendingCeremony TestApprovePageInChrome TestApprovePageRefusesASwappedRequest TestApprovePageSSHSIGInChrome TestApproveCardNetworkDeadlines TestApproveCardPasskeyDeadline TestNeedsYouListsPendingRequests"
if [ "$MODE" != base ]; then WANT_CHROME="$WANT_CHROME TestNeedsYouNotifiesOnlyWhenAskedOnce TestApprovePageShowsAHostileValueAsText"; fi
for t in $WANT_CHROME; do grep -qE -- "^--- PASS: $t " "$L/chrome.log" || { echo "$t did not PASS in Chrome"; exit 1; }; done
echo "PRE_COMMIT_GATES_GREEN $(( $(date +%s) - s0 ))s" | tee "$L/pre_commit.ok"
EOF
# AFTER the controller commits (the tracked-binary gate reads committed blobs; this leg runs go test ./...):
bash -c 'set -o pipefail; export AILANG_BIN=$HOME/.pinned-ailang/ailang; unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER; cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter387-a3a-s4 && VERIFY_GO_LEGS=plain ./scripts/verify_go.sh 2>&1 | tee $HOME/.ailang/state/world-iter387/gates/verify_go.log; echo "verify_go rc=$?" | tee -a $HOME/.ailang/state/world-iter387/gates/verify_go.log'
git -C /Users/voightkampff/dev/sunholo-data/.wt-world-iter387-a3a-s4 status --short   # only the two plan files, until the controller commits them
```

The pre-commit block is the planner's `tools/gates.sh <tree> <logs> final`; `base` mode omits the sprint's
new Chrome test (`TestNeedsYouNotifiesOnlyWhenAskedOnce`) and expects the same 9 `TestApproveAwait` passes. The
hostile-value test is matched by the Chrome leg's `TestApprovePage` prefix and must PASS there too (its
`--- PASS` line is read from `chrome.log`, §1.2 C21).

**At the base**: every pre-commit leg is green at `bed9d135` in 521 s (§1.1). **On the final tree**: `PRE_COMMIT_GATES_GREEN 548s` on the round-trip tree: pkgs all `ok` (`host/daemon` 302.4 s, `cmd/ailang-worldd` 158.3 s), 57 named tests `--- PASS`, the AW leg 9, `-race` ok ×3, the live-Chrome leg every test and subtest `--- PASS` incl. XA3, XA4 and the hostile-value test (§1.2 C20, C21, C23); the after-commit `verify_go.sh` (§1.2 C25): rc 1 on its first run (one known-list test, green alone), and rc 1 on its second, alone (a different untouched `host/broker` test, green alone at the base and on the final tree; not on the known list, flagged for the controller).

**What CI pins and what it does not (design R3).** CI has no Chrome step: XA3, XA4 and the hostile-value test
are counted `SKIP`s there (`SKIP COUNTED: the passkey page browser test needs WORLD_CHROME …`), as sprint 3's
browser tests. CI **does** run `TestSurfaceAppPure` under node (`WORLD_EXEC_NODE`): the A02 hostile rows,
Decided's rendering and toggle states, the opt-in's three states and every `notifications` rule (first read,
once per id, granted only, refused/recorded, the title, tag, body and href) — 36 of the sprint's 60 mutants die
there. The effects at the page boundary (the click wiring, the constructor, the click handler, `N`, Enter)
are pinned by the live-Chrome leg above, which is the controller's.

**Known load- or environment-sensitive (re-run once, alone; do not chase)**:
`TestSeedRestartIgnoresASurvivingOldCopy` (host/broker), `TestHostsLaunchersExecHost`, `host/broker`'s
exec-probe timeouts (incl. `TestExecProtectedCheckTimeout/reserve`) and its local `-race`, the local `-race`
of `host/daemon` and `host/verifygate` (incl. `TestSeToolsGateArmsAreLoadBearing/plan-bytes`; this sprint's
`-race` leg is the approval subset only), `TestMCPPostBudgetProductionConstants/batch` (host/daemon),
`TestBackoff` (host/surfacefollow), `TestDemoSeedPreflightRefusals`, `TestWorktreeSnapshotRefusesAConfiguredFilter`,
`TestMissionConfigDefaultFlowIntegration/invalid_ledger`, `scripts/marketingcheck`'s supervisor tests and
`TestMarketingCaptureFakeChromeSixShots/inventory`, and anything run under `nohup`; **measured by this planner
under load** (each green alone at the base and on the final tree, §1.1, §1.2 C25): `TestWorldContextAndStatusline`
(cmd/ailang-worldd), `TestMissionConfigDefaultFlowIntegration/valid_ledger` (the same fixture timeout as
`/invalid_ledger`) and `host/broker`'s `TestPublishApprovalRefusalSetWithALandedPositiveControl` (AC9's counter
`POST=2`) — the last is new to the list. XA5's deadline windows
(9.5–25 s) assume nothing else heavy runs; a red there under load is re-run once alone. Any other red, any race
report, any red in `host/surfacepage`, in `host/daemon`'s `TestApprov*`/`TestRequest*`/`TestSurfaceApp*`/
`TestDaemonCredential*`, or in `cmd/ailang-worldd`'s `TestApprove*`/`TestNeedsYou*`/`TestPasskeyPage*`/
`TestSurfaceApprove*`/`TestSurfaceSSHSIG*`, is a regression. CI's `go-race-verify` is the race verdict.

**Must stay green, with only the named edits**: every a0a1, A1b, A2 and A3a-sprint-1/2/3 test; the only edits
to existing tests are `TestSurfaceAppPure`'s harness (`app_test.go`) and XA3 (`needs_you_chrome_test.go`,
extended in place). In particular `TestFrozenV1RouteTableMatchesMux` (13), `TestSurfaceAppRoutes`,
`TestApprovePageRoutes`, `TestPasskeyPageRoutes` (no embedded byte of theirs changes), A2's XB1/XB2 (the driver
only gains commands), XA1, XA2, XA5, the two sprint-3 planner browser tests, `TestDaemonCredentialWriterSites`
and every §7 inventory.

---

## 7. Executor directive notes

- **No git writes.** Never `git add`, `commit`, `stash`, `checkout`, `reset`, `restore`, `rebase`, `push`,
  `worktree` or `apply`. `patch -p1` (§2) is the only writer; read-only `git status`/`git diff` are fine.
- **Never `nohup`, never background a test leg.** Every acceptance command runs in the foreground of one bash
  shell, `tee`'d to `$HOME/.ailang/state/world-iter387/exec-logs/`. Poll nothing in the background: there is
  nothing to background.
- **Chrome only through the tests.** The Chrome legs (`[chrome]`) are `go test` with `WORLD_CHROME` and
  `WORLD_REQUIRE_PASSKEY_BROWSER=1` on that command line only: headless, a temporary profile under the test's
  temp dir, a CDP virtual authenticator, test keys, and the `notify` recorder (no notification is raised).
  Never open Chrome by hand, never a real profile, never grant a real profile notifications; no `ssh-keygen`, no
  real store.
- **Runs.** Run 1 = M9a, M9b, M9c; Run 2 = M10a, M10b, M10c. Each run reads §0, §2 and its own §3 sections,
  then applies its patches in order.
- **Expected numstat after each sub-milestone, measured AFTER the controller's commit of it** (an uncommitted
  tree hides the new untracked files from numstat):
  ```
  bash -c 'cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter387-a3a-s4 && git diff --numstat bed9d135fe38a1c8f7e2784414aaac348b76879f HEAD -- cmd host | awk "{s+=\$1+\$2} END {print s}"'
  ```
  After the controller's commit of M9a **6**, M9b **36**, M9c **93**, M10a **239**, M10b **357**, M10c **389**
  (and `-- docs` **4** after M10c). **Any other value means a patch was not applied as written: stop and
  report.**
- **Snapshot after every sub-milestone, cumulatively**:
  ```
  bash -c 'set -euo pipefail; M=M9a; W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter387-a3a-s4
    D=$HOME/.ailang/state/world-iter387/snap/$M; rm -rf "$D"; mkdir -p "$D"; cd "$W"
    git status --porcelain --untracked-files=all | while read -r st f; do
      case "$f" in design_docs/planned/w-world-approvals-by-key-a3a-sprint4-plan.md|design_docs/planned/sprint_w-world-approvals-by-key-a3a-s4.json) continue;; esac
      mkdir -p "$D/$(dirname "$f")"; cp -p "$f" "$D/$f"; done
    (cd "$D" && find . -type f ! -name MANIFEST.sha256 | sort | xargs shasum -a 256) > "$D/MANIFEST.sha256"
    cat "$D/MANIFEST.sha256"'
  ```
  (`M` = `M9a` … `M10c`.) No file is deleted, so the porcelain status has only `M` and `??`; stop on any other
  status. The controller rebuilds one commit per sub-milestone (`$HOME/.ailang/state/world-iter387/commit_ms.sh`
  with `WT` and the snapshot root pointed here; **commit the plan BEFORE running it — it deletes untracked
  files**).
- **Stop conditions.** A `patch --dry-run` failure; a numstat other than the expected value; a red
  earlier-milestone test; a `--- FAIL`, or a `--- SKIP` in a §3 `-v` run; any need to touch a file outside §5;
  running Chrome other than through the `[chrome]` test commands; any reference to `~/.ssh`, `~/.ailang/world`
  or a real `/dev/tty`. Stop and report; do not work around.
- **Timing.** M9a seconds; M9b–M10b ~1 min each (one or a few Chrome tests, ~2–5 s each, plus the compile
  fence ~25–60 s); M10c ~3–4 min (the whole `cmd/ailang-worldd`, ~100–125 s at the base).
- **Shell.** Wrap loops, globs and `$VAR` expansions in `bash -c` or a `bash <<'EOF'` heredoc (zsh does not
  word-split). Export `AILANG_BIN` in the same shell as the Go tests. `tee` every gate's output to a file;
  never pipe a gate's only copy through `tail`.

---

## 8. Risks

- **R1 — Live-browser tests are local only** (design R3): CI skips XA3, XA4 and the hostile-value test
  (counted); the controller's live-Chrome leg (§6) is their verdict. The node harness pins the renderers and
  every notification rule in CI.
- **R2 — A notification's display is not tested anywhere** (S4-D8, design §2.9): XA4 judges the page's logic
  through the recorder; whether Chrome, Safari or Firefox on Mark's machines show it, and whether its click
  focuses the tab there, is the attended exit run's (design §10 step 4).
- **R3 — In v1 mode a new request waits for the next refresh** (§1.4 P3): no new timer, as designed; follow mode
  polls every 15 s. Recorded in Appendix C.
- **R4 — `app.js` is under active style work** (design §8): this sprint touches `needsYouHTML` and the now
  card's two calls, `loadRequests`, the stage click's delegation (two lines), the keydown handler (two lines)
  and `KEYROWS`; a conflicting edit there is a controller rebase.
- **R5 — Safari, Firefox, a `*.ts.net` HTTPS origin, a real authenticator and real notifications are
  UNVERIFIED** (design R5, R9); the attended exit run checks them.

---

## 9. Acceptance criteria (the judge scores against these)

1. **The judge's survivors are dead** (M9a): each of `needsYouHTML`'s `esc()` calls on `q.agent`, the link's
   `id=` and `href`, `q.recorded`, the refused id and the `unavailable` text, removed alone, fails
   `TestSurfaceAppPure` and the `host/surfacepage` suite with `WORLD_EXEC_NODE` set; at the base the same six
   edits survive (both runs in the ledger).
2. **Decided** (M9b, XA3's half): under Needs you, "You · Minted · <episode> · <approvedBy> · <hh:mm>" for the
   last three decided requests of `requests.json`, newest first, escaped; **Show all <n>** / **Show fewer**
   toggles the daemon's 20; nothing when none; after a real passkey approval and `approve --await`, Decided
   shows `webauthn:mark-mac` in live Chrome.
3. **The opt-in** (M9c, XA4): "Notify me" while the permission is `default`; `requestPermission` is called only
   from its click (zero calls on load and refresh; the one call under a user gesture); "Notifications on" /
   "blocked" after the answer.
4. **The notification** (M9c, XA4): only while `granted`; one per pending id the tab has not seen, never on the
   page's first read, never twice, never for a refused or recorded request; title "World needs you", tag = the
   id, body "Approve mint request for <episode> (<n> grant(s)), expires <hh:mm>"; its click focuses the tab and
   opens `/surface/approve/<id>`; zero CSP violations; nothing is ever displayed by the tests.
5. **`N`** (M10a, XA3): `n` and `N` focus the first Needs-you link on the timeline and from a brane; Enter on it
   opens the card; the keys popover lists `N`.
6. **The card shows a hostile signed value as text** (M10c): the model field's text equals the hostile string;
   no element, script or CSP violation results.
7. **No Go production change**: `surfaceCSP`, every route, `host/approval`, `host/store`, `host/authority`,
   `host/daemon`, `model.go`, `surface.js`, `approve.js`, `passkey.js`, both `index.html`s and `app.css`
   unchanged; no new module, package, exec site, context root or `/v1` route; no CORS; no `t.Parallel`; no
   `innerHTML` with unescaped data, no `eval(`, no `document.cookie` in the new code.
8. Each §4 mutant, applied alone, makes its named test **and** its package suite fail (inventories excluded);
   the named equivalents survive for the stated reasons; the ledger has a line per run.
9. Every §6 gate passes, including every inventory by name, the `-race` legs, the live-Chrome leg (every
   test and subtest `--- PASS`, incl. XA3, XA4 and the hostile-value test) and `VERIFY_GO_LEGS=plain
   ./scripts/verify_go.sh` after the commit (modulo §6's re-run rule); `git status --short` lists only §5.
10. The numstat after each sub-milestone's commit equals §7's value; the sprint total is ≤ 2,000 (389).
11. **A3a is then code-complete except Mark's attended exit run** (design §10, QUICKSTART §17, still *attended —
    pending first verbatim run*); §17 describes Decided, `N`, and the notification step (M10c), so the exit
    run's step 4 ("turn on Notify me and see one notification for a new request") has its text.
12. No real key, store, `/dev/tty`, `~/.ssh` access, real browser profile, real notification, mint into a real
    store, publish or `ssh-keygen` anywhere in the sprint (D-WORLD-109); Chrome only headless with temporary
    profiles and a virtual authenticator.

---

## Appendix C — after this sprint (A4 or later, and the attended exit run)

- **Nothing of A3a's design is left to build.** With this sprint, every §6 test row of the design exists (XA1–
  XA5, PG2, PG3, PG3b, RT1–RT13, SA*, JG*, AW*, AM*, RG1, CF1, E2E1) and every §2 item is built, except what is
  attended by design.
- **The attended exit run** (design §10, Mark; NOT run by the loop, D-WORLD-109), with QUICKSTART §17 executed
  verbatim: (1) `serve --surface-origin https://<host>.ts.net` on a scratch store behind `tailscale serve`;
  (2) record the `Host` the daemon saw from the phone; (3) `episode start … --request`, `approve <id> --await`,
  Needs you → the card → Approve → Face ID; the minter mints; **Decided shows `webauthn:<name>`**; (4) on the
  Mac, Safari or Chrome with Touch ID; **Notify me**, allow, then a new request: record per browser whether a
  notification is displayed and whether its click opened the card (the only display check, R2); (5) SSHSIG
  through the tailnet origin with no passkey; (6) refusals: approve again (already used), let one expire; (7)
  A3b only if `desktop/` exists; (8) record public values only; (9) confirm or veto D-WORLD-114's two
  amendments. Then remove §17's *attended — pending* tag.
- **For A4 or a follow-up (not designed here)**: (a) in v1 mode a new request is listed and notified only at
  the next refresh (a world-log commit or a reload; follow mode polls every 15 s) — §1.4 P3; if the exit run
  finds that too slow, a `requests.json` poll is a design change (the design said "no new timer"); (b)
  `model.json`'s `needsYou` stays `[]` (design §2.8: A4 may use it for decision packets); (c) the notification's
  `hh:mm` is UTC like the rest of the app — a local-time display is an app-wide choice, not this row's; (d)
  design Q2 (`episode start --request --await`) stays a possible follow-up.
- **D-WORLD-114** stays open (default A applied); if Mark vetoes (b), the minter's mint location returns to
  design.

---

## Appendix P — the patches (apply with §2; generated from the compiled scratch stages; do not edit)

### Patch M9a — the judge's survivors are pinned

````diff M9a
--- a/host/surfacepage/app_test.go
+++ b/host/surfacepage/app_test.go
@@ -181,6 +181,12 @@ assert.equal(app.needsYouHTML(null, 0), '<div class="nsec" data-needs><div class
   assert(h.includes('Needs you · 2') && h.includes('<a id="needs-mr-0123456789abcdef" href="/surface/approve/mr-0123456789abcdef"><b>w&lt;1&gt;</b> · claude · 2 grants · expires in 11 min</a> · approved by webauthn:mark, waiting for the minter</div>'), h);
   assert(h.includes('mr-00000000000000ff · refused: digest &lt;bad&gt;') && !h.includes('/surface/approve/mr-00000000000000ff'), h); }
 assert(app.needsYouHTML({pending: [], unavailable: 'the mint requests could not be read'}, 0).includes('>the mint requests could not be read<'));
+// Every value is escaped where it lands: hostile text in each field and in the link's id and href (sprint-4 judge A02/A03/A05).
+assert.equal(app.needsYouHTML({pending: [{id: 'mr-"><img src=x onerror=1>', episode: 'e', agent: '<img src=x onerror=2>', grants: 1, expires: 60, recorded: '"><script>3</script>'},
+  {id: 'mr-<b>"4</b>', episode: 'f', refused: '<i>5</i>'}]}, 0), '<div class="nsec" data-needs><div class="nh" style="color:var(--youShell)">Needs you · 2</div>'
+  + '<div class="nl"><a id="needs-mr-&quot;&gt;&lt;img src=x onerror=1&gt;" href="/surface/approve/mr-&quot;&gt;&lt;img src=x onerror=1&gt;"><b>e</b> · &lt;img src=x onerror=2&gt; · 1 grant · expires in 1 min</a>'
+  + ' · approved by &quot;&gt;&lt;script&gt;3&lt;/script&gt;, waiting for the minter</div><div class="nl sub">mr-&lt;b&gt;&quot;4&lt;/b&gt; · refused: &lt;i&gt;5&lt;/i&gt;</div></div>');
+assert(app.needsYouHTML({pending: [], unavailable: '<i>"down"</i>'}, 0).includes('<div class="nl sub">&lt;i&gt;&quot;down&quot;&lt;/i&gt;</div>'));
 assert.equal(app.serviceLineHTML({line:'<build & file>',warn:false}), '<div class="nl sub" data-service>&lt;build &amp; file&gt;</div>');
 assert.equal(app.serviceLineHTML({line:'restart',warn:true}), '<div class="nl sub" data-service style="color:var(--you)">restart</div>');
 
````

### Patch M9b — the decision list

````diff M9b
--- a/host/surfacepage/app/app.js
+++ b/host/surfacepage/app/app.js
@@ -431,8 +431,18 @@
       : `<div class="nl"><a id="needs-${esc(q.id)}" href="/surface/approve/${esc(q.id)}"><b>${esc(q.episode)}</b> · ${esc(q.agent)} · ${q.grants} grant${q.grants === 1 ? '' : 's'} · expires in ${Math.max(0, Math.ceil((q.expires - now) / 60))} min</a>`
         + (q.recorded ? ` · approved by ${esc(q.recorded)}, waiting for the minter` : '') + '</div>').join('') + '</div>';
   }
+  // decidedHTML is the decision list under Needs you (design §2.8): the last few of GET
+  // /surface/requests.json's decided requests (newest first; the daemon sends 20), and all of them on request.
+  const DECIDED_FEW = 3;
+  function decidedHTML(r, all) {
+    const d = r && Array.isArray(r.decided) ? r.decided : [];
+    if (!d.length) return '';
+    return '<div class="nsec" data-decided><div class="nh">Decided</div>' + (all ? d : d.slice(0, DECIDED_FEW)).map(x =>
+      `<div class="nl sub">You · Minted · ${esc(x.episode)} · ${esc(x.approvedBy)} · ${hm(new Date(x.at * 1000))}</div>`).join('')
+      + (d.length > DECIDED_FEW ? `<button type="button" class="nl run" id="decided-all" aria-expanded="${!!all}">${all ? 'Show fewer' : 'Show all ' + d.length}</button>` : '') + '</div>';
+  }
 
-  if (typeof module !== 'undefined') module.exports = { needsYouHTML, isSharedKey, sourcesHTML, openRunningHTML, fetchCoherent, serviceLineHTML, splitFrames, applyDelta, nextBackoff, liveRequest, liveStatus, esc, cleanTitle, shortTitle, prettyModel, agentLine, summaryOf, rowLabel, verdict, LEVELS, levelXs, blendXS,
+  if (typeof module !== 'undefined') module.exports = { needsYouHTML, decidedHTML, isSharedKey, sourcesHTML, openRunningHTML, fetchCoherent, serviceLineHTML, splitFrames, applyDelta, nextBackoff, liveRequest, liveStatus, esc, cleanTitle, shortTitle, prettyModel, agentLine, summaryOf, rowLabel, verdict, LEVELS, levelXs, blendXS,
     adaptModel, layoutGraph, stackBars, barHTML, frontHTML, railHTML, appPath, stackPath, parseAppQuery, addressTarget, lookOf, crumbsHTML, geometry: { PL, PR, PW } };
   if (typeof document === 'undefined') return;
 
@@ -614,6 +624,7 @@
       ${today}
       ${running}
       ${needsYouHTML(REQS, REQS ? REQS.now : 0)}
+      ${decidedHTML(REQS, decidedAll)}
     </div>`;
   }
 
@@ -889,8 +900,9 @@
     try { body = await r.json(); } catch (_) { body = null; }
     return { ok: r.ok, status: r.status, body };
   }
-  // REQS is GET /surface/requests.json, read with the model and on every refresh (no timer of its own).
-  let REQS = null;
+  // REQS is GET /surface/requests.json, read with the model and on every refresh (no timer of its own);
+  // decidedAll is the decision list's "show all".
+  let REQS = null, decidedAll = false;
   async function loadRequests() {
     try { REQS = (await getJSON('/surface/requests.json')).body; } catch (_) { REQS = null; }
   }
@@ -1284,6 +1296,7 @@
     if (n) return open(n.dataset.id);
     const qb = e.target.closest('[data-quiet]');
     if (qb) { const l = [...document.querySelectorAll('[data-quiet-list]')].find(x => x.dataset.quietList === qb.dataset.quiet); if (l) { l.hidden = !l.hidden; qb.setAttribute('aria-expanded', String(!l.hidden)); } return; }
+    if (e.target.closest('#decided-all')) { decidedAll = !decidedAll; return repaintNow(); }
     const lab = e.target.closest('.lhead');
     if (lab) { const key = lab.dataset.lane, s = lastState(key); return s ? openState(key, s.id) : say(`${nameOf(key)} has no states yet`); }
     if (e.target.closest('#now-b')) return goNow();
--- a/host/surfacepage/app_test.go
+++ b/host/surfacepage/app_test.go
@@ -171,7 +171,7 @@ const seams = require('./surface.js');
 const model = JSON.parse(fs.readFileSync('model.json', 'utf8'));
 const stacks = JSON.parse(fs.readFileSync('stacks.json', 'utf8'));
 const index = JSON.parse(fs.readFileSync('index.json', 'utf8'));
-assert.deepEqual(Object.keys(app).sort(), ['LEVELS', 'adaptModel', 'addressTarget', 'agentLine', 'appPath', 'applyDelta', 'barHTML', 'blendXS', 'cleanTitle', 'crumbsHTML', 'esc', 'fetchCoherent', 'frontHTML',
+assert.deepEqual(Object.keys(app).sort(), ['LEVELS', 'adaptModel', 'addressTarget', 'agentLine', 'appPath', 'applyDelta', 'barHTML', 'blendXS', 'cleanTitle', 'crumbsHTML', 'decidedHTML', 'esc', 'fetchCoherent', 'frontHTML',
   'geometry', 'isSharedKey', 'layoutGraph', 'levelXs', 'liveRequest', 'liveStatus', 'lookOf', 'needsYouHTML', 'nextBackoff', 'openRunningHTML', 'parseAppQuery', 'prettyModel', 'railHTML', 'rowLabel', 'serviceLineHTML', 'shortTitle', 'sourcesHTML', 'splitFrames', 'stackBars', 'stackPath', 'summaryOf', 'verdict']);
 
 assert.equal(app.serviceLineHTML(undefined), '');
@@ -187,6 +187,15 @@ assert.equal(app.needsYouHTML({pending: [{id: 'mr-"><img src=x onerror=1>', epis
   + '<div class="nl"><a id="needs-mr-&quot;&gt;&lt;img src=x onerror=1&gt;" href="/surface/approve/mr-&quot;&gt;&lt;img src=x onerror=1&gt;"><b>e</b> · &lt;img src=x onerror=2&gt; · 1 grant · expires in 1 min</a>'
   + ' · approved by &quot;&gt;&lt;script&gt;3&lt;/script&gt;, waiting for the minter</div><div class="nl sub">mr-&lt;b&gt;&quot;4&lt;/b&gt; · refused: &lt;i&gt;5&lt;/i&gt;</div></div>');
 assert(app.needsYouHTML({pending: [], unavailable: '<i>"down"</i>'}, 0).includes('<div class="nl sub">&lt;i&gt;&quot;down&quot;&lt;/i&gt;</div>'));
+// The decision list: the last three decided requests, newest first, and a toggle for all of them; every value escaped.
+assert.equal(app.decidedHTML(null, false), '');
+assert.equal(app.decidedHTML({decided: []}, true), '');
+{ const d = [{episode: 'w<4>', approvedBy: 'webauthn:"m"', at: 50520}, {episode: 'w3', approvedBy: 'ssh:<k>', at: 50460}, {episode: 'w2', approvedBy: 'tty', at: 3600}, {episode: 'w1', approvedBy: 'tty', at: 0}];
+  const few = '<div class="nsec" data-decided><div class="nh">Decided</div><div class="nl sub">You · Minted · w&lt;4&gt; · webauthn:&quot;m&quot; · 14:02</div>'
+    + '<div class="nl sub">You · Minted · w3 · ssh:&lt;k&gt; · 14:01</div><div class="nl sub">You · Minted · w2 · tty · 01:00</div>';
+  assert.equal(app.decidedHTML({decided: d}, false), few + '<button type="button" class="nl run" id="decided-all" aria-expanded="false">Show all 4</button></div>');
+  assert.equal(app.decidedHTML({decided: d}, true), few + '<div class="nl sub">You · Minted · w1 · tty · 00:00</div><button type="button" class="nl run" id="decided-all" aria-expanded="true">Show fewer</button></div>');
+  assert.equal(app.decidedHTML({decided: d.slice(0, 3)}, false), few + '</div>'); }
 assert.equal(app.serviceLineHTML({line:'<build & file>',warn:false}), '<div class="nl sub" data-service>&lt;build &amp; file&gt;</div>');
 assert.equal(app.serviceLineHTML({line:'restart',warn:true}), '<div class="nl sub" data-service style="color:var(--you)">restart</div>');
 
````

### Patch M9c — the opt-in and the notification

````diff M9c
--- a/cmd/ailang-worldd/needs_you_chrome_test.go
+++ b/cmd/ailang-worldd/needs_you_chrome_test.go
@@ -31,7 +31,7 @@ func TestNeedsYouListsPendingRequests(t *testing.T) {
 	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface", "raw": true})
 	var needs string
 	b.do(t, &needs, "waitFor", map[string]any{"expression": `(() => { const e = document.querySelector('[data-needs]'); return e && e.textContent; })()`})
-	if want := "Needs you · 1ep1 · claude · 10 grants · expires in 10 min"; needs != want {
+	if want := "Needs you · 1ep1 · claude · 10 grants · expires in 10 minNotify me"; needs != want {
 		t.Fatalf("Needs you shows %q, want %q", needs, want)
 	}
 	newApproveRigIn(t, a.r, "ep2", now)
--- a/host/surfacepage/app/app.js
+++ b/host/surfacepage/app/app.js
@@ -422,14 +422,32 @@
   }
 
   // needsYouHTML is the now card's Needs you (row 201 A3, design §2.8): the pending mint requests of
-  // GET /surface/requests.json, each a link to its card; a request the daemon could not read shows why.
-  function needsYouHTML(r, now) {
+  // GET /surface/requests.json, each a link to its card; a request the daemon could not read shows why;
+  // then the notification opt-in for this browser's permission (design §2.9).
+  function needsYouHTML(r, now, permission) {
     const p = r && Array.isArray(r.pending) ? r.pending : [];
     const head = `<div class="nh" style="color:${p.length ? 'var(--youShell)' : 'var(--sub)'}">Needs you · ${p.length}</div>`;
-    if (!p.length) return `<div class="nsec" data-needs>${head}<div class="nl sub">${esc(r && r.unavailable || 'Nothing needs you')}</div></div>`;
+    if (!p.length) return `<div class="nsec" data-needs>${head}<div class="nl sub">${esc(r && r.unavailable || 'Nothing needs you')}</div>${notifyHTML(permission)}</div>`;
     return `<div class="nsec" data-needs>${head}` + p.map(q => q.refused ? `<div class="nl sub">${esc(q.id)} · refused: ${esc(q.refused)}</div>`
       : `<div class="nl"><a id="needs-${esc(q.id)}" href="/surface/approve/${esc(q.id)}"><b>${esc(q.episode)}</b> · ${esc(q.agent)} · ${q.grants} grant${q.grants === 1 ? '' : 's'} · expires in ${Math.max(0, Math.ceil((q.expires - now) / 60))} min</a>`
-        + (q.recorded ? ` · approved by ${esc(q.recorded)}, waiting for the minter` : '') + '</div>').join('') + '</div>';
+        + (q.recorded ? ` · approved by ${esc(q.recorded)}, waiting for the minter` : '') + '</div>').join('') + notifyHTML(permission) + '</div>';
+  }
+  // notifyHTML is the opt-in: "Notify me" while this browser has not been asked (permission "default"),
+  // its answer once it has, and nothing where the browser has no notifications.
+  function notifyHTML(permission) {
+    if (permission === 'default') return '<button type="button" class="nl run" id="notify-me">Notify me</button>';
+    if (permission === 'granted') return '<div class="nl sub" data-notify>Notifications on</div>';
+    return permission === 'denied' ? '<div class="nl sub" data-notify>Notifications are blocked for this page</div>' : '';
+  }
+  // notifications is what one read of GET /surface/requests.json raises (design §2.9): a "World needs you"
+  // for each pending request this tab has not seen, while notifications are granted; each opens its card.
+  // seen is null before the page's first read, which raises nothing (the page itself shows those); a read
+  // that failed changes nothing. Refused and already-recorded requests are seen but raise nothing.
+  function notifications(r, seen, permission) {
+    if (!r || !Array.isArray(r.pending)) return { notes: [], seen };
+    const fresh = seen ? r.pending.filter(q => !seen.has(q.id) && !q.refused && !q.recorded) : [];
+    return { seen: new Set([...(seen || []), ...r.pending.map(q => q.id)]), notes: permission !== 'granted' ? [] : fresh.map(q => ({ title: 'World needs you', tag: q.id,
+      body: `Approve mint request for ${q.episode} (${q.grants} grant${q.grants === 1 ? '' : 's'}), expires ${hm(new Date(q.expires * 1000))}`, href: '/surface/approve/' + q.id })) };
   }
   // decidedHTML is the decision list under Needs you (design §2.8): the last few of GET
   // /surface/requests.json's decided requests (newest first; the daemon sends 20), and all of them on request.
@@ -442,7 +460,7 @@
       + (d.length > DECIDED_FEW ? `<button type="button" class="nl run" id="decided-all" aria-expanded="${!!all}">${all ? 'Show fewer' : 'Show all ' + d.length}</button>` : '') + '</div>';
   }
 
-  if (typeof module !== 'undefined') module.exports = { needsYouHTML, decidedHTML, isSharedKey, sourcesHTML, openRunningHTML, fetchCoherent, serviceLineHTML, splitFrames, applyDelta, nextBackoff, liveRequest, liveStatus, esc, cleanTitle, shortTitle, prettyModel, agentLine, summaryOf, rowLabel, verdict, LEVELS, levelXs, blendXS,
+  if (typeof module !== 'undefined') module.exports = { needsYouHTML, decidedHTML, notifications, isSharedKey, sourcesHTML, openRunningHTML, fetchCoherent, serviceLineHTML, splitFrames, applyDelta, nextBackoff, liveRequest, liveStatus, esc, cleanTitle, shortTitle, prettyModel, agentLine, summaryOf, rowLabel, verdict, LEVELS, levelXs, blendXS,
     adaptModel, layoutGraph, stackBars, barHTML, frontHTML, railHTML, appPath, stackPath, parseAppQuery, addressTarget, lookOf, crumbsHTML, geometry: { PL, PR, PW } };
   if (typeof document === 'undefined') return;
 
@@ -623,7 +641,7 @@
       ${sourcesHTML(D.repos)}
       ${today}
       ${running}
-      ${needsYouHTML(REQS, REQS ? REQS.now : 0)}
+      ${needsYouHTML(REQS, REQS ? REQS.now : 0, notifyPermission())}
       ${decidedHTML(REQS, decidedAll)}
     </div>`;
   }
@@ -902,10 +920,16 @@
   }
   // REQS is GET /surface/requests.json, read with the model and on every refresh (no timer of its own);
   // decidedAll is the decision list's "show all".
-  let REQS = null, decidedAll = false;
+  let REQS = null, decidedAll = false, seenIds = null;
+  const notifyPermission = () => typeof Notification === 'undefined' ? '' : Notification.permission;
   async function loadRequests() {
     try { REQS = (await getJSON('/surface/requests.json')).body; } catch (_) { REQS = null; }
+    const n = notifications(REQS, seenIds, notifyPermission());
+    seenIds = n.seen;
+    n.notes.forEach(x => { const note = new Notification(x.title, { body: x.body, tag: x.tag }); note.onclick = () => { window.focus(); location.assign(x.href); }; });
   }
+  // askNotify asks this browser for notifications, only from the "Notify me" click (a user gesture).
+  function askNotify() { Promise.resolve(Notification.requestPermission()).then(repaintNow, repaintNow); }
   const FOLLOW_MODEL = 'world/surface-model/v2', POLL_MS = 15000;
   let followMode = false, shownGen = null, shownAt = '', pollTimer = null, liveStarted = false;
   const poll = { note: '' };
@@ -1297,6 +1321,7 @@
     const qb = e.target.closest('[data-quiet]');
     if (qb) { const l = [...document.querySelectorAll('[data-quiet-list]')].find(x => x.dataset.quietList === qb.dataset.quiet); if (l) { l.hidden = !l.hidden; qb.setAttribute('aria-expanded', String(!l.hidden)); } return; }
     if (e.target.closest('#decided-all')) { decidedAll = !decidedAll; return repaintNow(); }
+    if (e.target.closest('#notify-me')) return askNotify();
     const lab = e.target.closest('.lhead');
     if (lab) { const key = lab.dataset.lane, s = lastState(key); return s ? openState(key, s.id) : say(`${nameOf(key)} has no states yet`); }
     if (e.target.closest('#now-b')) return goNow();
--- a/host/surfacepage/app_test.go
+++ b/host/surfacepage/app_test.go
@@ -172,7 +172,7 @@ const model = JSON.parse(fs.readFileSync('model.json', 'utf8'));
 const stacks = JSON.parse(fs.readFileSync('stacks.json', 'utf8'));
 const index = JSON.parse(fs.readFileSync('index.json', 'utf8'));
 assert.deepEqual(Object.keys(app).sort(), ['LEVELS', 'adaptModel', 'addressTarget', 'agentLine', 'appPath', 'applyDelta', 'barHTML', 'blendXS', 'cleanTitle', 'crumbsHTML', 'decidedHTML', 'esc', 'fetchCoherent', 'frontHTML',
-  'geometry', 'isSharedKey', 'layoutGraph', 'levelXs', 'liveRequest', 'liveStatus', 'lookOf', 'needsYouHTML', 'nextBackoff', 'openRunningHTML', 'parseAppQuery', 'prettyModel', 'railHTML', 'rowLabel', 'serviceLineHTML', 'shortTitle', 'sourcesHTML', 'splitFrames', 'stackBars', 'stackPath', 'summaryOf', 'verdict']);
+  'geometry', 'isSharedKey', 'layoutGraph', 'levelXs', 'liveRequest', 'liveStatus', 'lookOf', 'needsYouHTML', 'nextBackoff', 'notifications', 'openRunningHTML', 'parseAppQuery', 'prettyModel', 'railHTML', 'rowLabel', 'serviceLineHTML', 'shortTitle', 'sourcesHTML', 'splitFrames', 'stackBars', 'stackPath', 'summaryOf', 'verdict']);
 
 assert.equal(app.serviceLineHTML(undefined), '');
 // Row 201 A3: Needs you is the pending requests, each a link to its card; a refused row shows why, with no link.
@@ -187,6 +187,24 @@ assert.equal(app.needsYouHTML({pending: [{id: 'mr-"><img src=x onerror=1>', epis
   + '<div class="nl"><a id="needs-mr-&quot;&gt;&lt;img src=x onerror=1&gt;" href="/surface/approve/mr-&quot;&gt;&lt;img src=x onerror=1&gt;"><b>e</b> · &lt;img src=x onerror=2&gt; · 1 grant · expires in 1 min</a>'
   + ' · approved by &quot;&gt;&lt;script&gt;3&lt;/script&gt;, waiting for the minter</div><div class="nl sub">mr-&lt;b&gt;&quot;4&lt;/b&gt; · refused: &lt;i&gt;5&lt;/i&gt;</div></div>');
 assert(app.needsYouHTML({pending: [], unavailable: '<i>"down"</i>'}, 0).includes('<div class="nl sub">&lt;i&gt;&quot;down&quot;&lt;/i&gt;</div>'));
+// The notification opt-in follows the browser's permission; nothing without the API (the calls above).
+assert(app.needsYouHTML(null, 0, 'default').endsWith('<div class="nl sub">Nothing needs you</div><button type="button" class="nl run" id="notify-me">Notify me</button></div>'));
+assert(app.needsYouHTML({pending: [{id: 'mr-1', episode: 'e', agent: 'a', grants: 1, expires: 60}]}, 0, 'granted').endsWith(' min</a></div><div class="nl sub" data-notify>Notifications on</div></div>'));
+assert(app.needsYouHTML(null, 0, 'denied').endsWith('<div class="nl sub" data-notify>Notifications are blocked for this page</div></div>'));
+// A notification per request this tab has not seen, while granted: none on the page's first read, none twice, none for a refused or recorded one.
+{ const r1 = {pending: [{id: 'mr-a', episode: 'e1', grants: 2, expires: 50520}]}, r2 = {pending: r1.pending.concat([{id: 'mr-b', episode: 'e<2>', grants: 1, expires: 50580},
+    {id: 'mr-c', episode: 'e3', refused: 'bad'}, {id: 'mr-d', episode: 'e4', grants: 3, expires: 60, recorded: 'webauthn:m'}])};
+  assert.deepEqual(app.notifications(null, null, 'granted'), {notes: [], seen: null});
+  const first = app.notifications(r1, null, 'granted');
+  assert.deepEqual(first.notes, []); assert.deepEqual([...first.seen], ['mr-a']);
+  assert.equal(app.notifications({}, first.seen, 'granted').seen, first.seen);
+  const second = app.notifications(r2, first.seen, 'granted');
+  assert.deepEqual(second.notes, [{title: 'World needs you', tag: 'mr-b', body: 'Approve mint request for e<2> (1 grant), expires 14:03', href: '/surface/approve/mr-b'}]);
+  assert.deepEqual([...second.seen], ['mr-a', 'mr-b', 'mr-c', 'mr-d']);
+  assert.deepEqual(app.notifications(r2, second.seen, 'granted').notes, []);
+  const asked = app.notifications(r2, first.seen, 'default');
+  assert.deepEqual(asked.notes, []); assert.deepEqual([...asked.seen], ['mr-a', 'mr-b', 'mr-c', 'mr-d']);
+  assert.equal(app.notifications({pending: [{id: 'mr-e', episode: 'e5', grants: 2, expires: 50520}]}, new Set(), 'granted').notes[0].body, 'Approve mint request for e5 (2 grants), expires 14:02'); }
 // The decision list: the last three decided requests, newest first, and a toggle for all of them; every value escaped.
 assert.equal(app.decidedHTML(null, false), '');
 assert.equal(app.decidedHTML({decided: []}, true), '');
````

### Patch M10a — `N`, Enter, and XA3's Decided half

````diff M10a
--- a/cmd/ailang-worldd/needs_you_chrome_test.go
+++ b/cmd/ailang-worldd/needs_you_chrome_test.go
@@ -1,9 +1,11 @@
 package main
 
-// Row 201 A3 sprint 3 M8d (design §2.8, §6 XA3's list half): Needs you in the surface app, in a
-// live Chrome.
+// Row 201 A3 sprint 3 M8d and sprint 4 (design §2.8, §6 XA3): Needs you and the decision list in the
+// surface app, in a live Chrome.
 
 import (
+	"bytes"
+	"fmt"
 	"os"
 	"path/filepath"
 	"strings"
@@ -11,14 +13,9 @@ import (
 	"time"
 )
 
-// XA3, the list half (Decided is sprint 4's, plan Appendix C): /surface's now card shows the pending
-// request under "Needs you · 1", from GET /surface/requests.json; a second request appears on the
-// app's next refresh (no timer of its own); the first one's link opens its card, ready.
-func TestNeedsYouListsPendingRequests(t *testing.T) {
-	b := startPasskeyBrowser(t)
-	b.do(t, nil, "tab", nil)
-	now := time.Now().Unix()
-	a := newApproveRig(t, "ep1", now)
+// rigSurface serves the daemon over credentials store creds with a copy of the rig's surface index.
+func rigSurface(t *testing.T, creds string) string {
+	t.Helper()
 	rig, err := os.ReadFile(filepath.Join("..", "..", "host", "surface", "testdata", "index-rig.json"))
 	if err != nil {
 		t.Fatal(err)
@@ -27,29 +24,112 @@ func TestNeedsYouListsPendingRequests(t *testing.T) {
 	if err := os.WriteFile(index, rig, 0o600); err != nil {
 		t.Fatal(err)
 	}
-	origin := cardOrigin(t, a.credPath, index, nil)
+	return cardOrigin(t, creds, index, nil)
+}
+
+// sectionText waits until the now card's section sel (data-needs, data-decided) has text starting
+// with prefix, and returns it.
+func sectionText(t *testing.T, b *passkeyBrowser, sel, prefix string) string {
+	t.Helper()
+	var s string
+	b.do(t, &s, "waitFor", map[string]any{"expression": fmt.Sprintf(`(() => { const e = document.querySelector('[%s]'); return e && e.textContent.startsWith(%q) && e.textContent; })()`, sel, prefix)})
+	return s
+}
+
+// pollSurface runs the app's own refresh (window.World.poll) and waits until it has finished.
+func pollSurface(t *testing.T, b *passkeyBrowser) {
+	t.Helper()
+	b.do(t, nil, "waitFor", map[string]any{"expression": `(window.__polled = 0, window.World.poll().then(() => { window.__polled = 1; }), 1)`})
+	b.do(t, nil, "waitFor", map[string]any{"expression": `window.__polled === 1`})
+}
+
+// keyView is the driver's key answer.
+type keyView struct{ Location, Active, View string }
+
+// XA3: /surface's now card shows the pending request under "Needs you · 1", from GET
+// /surface/requests.json; a second request appears on the app's next refresh (no timer of its own).
+// N focuses the first one's link, from the timeline and from a brane, and Enter opens its card.
+// Approved there and minted by approve --await, the request is under Decided as "You · Minted · ep1 ·
+// webauthn:mark-mac · <hh:mm>"; with four decided, Decided shows three and "Show all 4" toggles.
+func TestNeedsYouListsPendingRequests(t *testing.T) {
+	b := startPasskeyBrowser(t)
+	b.do(t, nil, "tab", map[string]any{"authenticator": "uv"})
+	now := time.Now().Unix()
+	a := newApproveRig(t, "ep1", now)
+	origin := rigSurface(t, a.credPath)
+	virtualPasskey(t, b, a.credPath, origin, "mark-mac", 7, now-100)
 	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface", "raw": true})
-	var needs string
-	b.do(t, &needs, "waitFor", map[string]any{"expression": `(() => { const e = document.querySelector('[data-needs]'); return e && e.textContent; })()`})
-	if want := "Needs you · 1ep1 · claude · 10 grants · expires in 10 minNotify me"; needs != want {
+	if needs, want := sectionText(t, b, "data-needs", "Needs you · 1"), "Needs you · 1ep1 · claude · 10 grants · expires in 10 minNotify me"; needs != want {
 		t.Fatalf("Needs you shows %q, want %q", needs, want)
 	}
-	newApproveRigIn(t, a.r, "ep2", now)
-	b.do(t, nil, "waitFor", map[string]any{"expression": `(window.World.poll(), 1)`})
-	b.do(t, &needs, "waitFor", map[string]any{"expression": `(() => { const e = document.querySelector('[data-needs]'); return e && e.textContent.startsWith('Needs you · 2') && e.textContent; })()`})
-	if !strings.Contains(needs, "ep1 · claude") || !strings.Contains(needs, "ep2 · claude") {
+	var k keyView
+	if b.do(t, &k, "key", map[string]any{"key": "n"}); k.Active != "needs-"+a.id || k.View != "timeline" {
+		t.Fatalf("N on the timeline: %+v", k)
+	}
+	var found bool
+	if b.do(t, &found, "waitFor", map[string]any{"expression": `window.World.open('e2')`}); !found {
+		t.Fatal("the rig has no episode e2")
+	}
+	b.do(t, nil, "waitFor", map[string]any{"expression": `document.getElementById('stage').dataset.view === 'brane'`})
+	if b.do(t, &k, "key", map[string]any{"key": "N"}); k.Active != "needs-"+a.id || k.View != "timeline" {
+		t.Fatalf("N in a brane: %+v", k)
+	}
+	second := newApproveRigIn(t, a.r, "ep2", now)
+	pollSurface(t, b)
+	if needs := sectionText(t, b, "data-needs", "Needs you · 2"); !strings.Contains(needs, "ep1 · claude") || !strings.Contains(needs, "ep2 · claude") {
 		t.Fatalf("after a refresh Needs you shows %q", needs)
 	}
 	var href string
 	b.do(t, &href, "waitFor", map[string]any{"expression": `document.getElementById('needs-` + a.id + `').getAttribute('href')`})
-	var at struct{ Location string }
-	b.do(t, &at, "press", map[string]any{"id": "needs-" + a.id, "page": true})
-	if href != "/surface/approve/"+a.id || at.Location != origin+"/surface/approve/"+a.id {
-		t.Fatalf("the link %q opened %q", href, at.Location)
+	b.do(t, nil, "waitFor", map[string]any{"expression": `(document.getElementById('needs-` + a.id + `').focus(), 1)`})
+	if b.do(t, &k, "key", map[string]any{"key": "Enter", "page": true}); href != "/surface/approve/"+a.id || k.Location != origin+"/surface/approve/"+a.id {
+		t.Fatalf("the link %q opened %q", href, k.Location)
 	}
 	b.do(t, nil, "waitFor", map[string]any{"expression": `document.body.dataset.state === 'ready'`})
 	var view cardView
-	if b.do(t, &view, "card", nil); !strings.HasPrefix(view.What, "Approve request "+a.id+" for episode ep1") || len(view.Errors) != 0 || len(view.Violations) != 0 {
+	if b.do(t, &view, "press", map[string]any{"id": "go", "states": []string{"recorded", "rejected", "failed", "unknown"}}); view.State != "recorded" ||
+		!strings.HasPrefix(view.What, "Approve request "+a.id+" for episode ep1") || len(view.Errors) != 0 || len(view.Violations) != 0 {
 		t.Fatalf("the card shows %+v", view)
 	}
+	var opened bool
+	aenv, _ := awaitEnv(time.Now().Unix(), &opened, nil)
+	var stdout, stderr bytes.Buffer
+	if code := runApprove(a.args("--await", "--await-for", "5s"), &stdout, &stderr, aenv); code != exitOK || opened {
+		t.Fatalf("approve --await: rc %d, terminal opened %v\nstderr %s", code, opened, stderr.String())
+	}
+	requireMintedBy(t, a, "webauthn:mark-mac", stdout.String())
+	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface", "raw": true})
+	at := time.Unix(a.row(t).ConsumedAt, 0).UTC().Format("15:04")
+	if d, want := sectionText(t, b, "data-decided", "Decided"), "DecidedYou · Minted · ep1 · webauthn:mark-mac · "+at; d != want {
+		t.Fatalf("Decided shows %q, want %q", d, want)
+	}
+	if needs := sectionText(t, b, "data-needs", "Needs you · 1"); !strings.HasPrefix(needs, "Needs you · 1ep2 · claude") {
+		t.Fatalf("after the mint Needs you shows %q", needs)
+	}
+	raw := rawCredentials(t, a.credPath)
+	defer func() { _ = raw.Close() }()
+	for i, r := range []approveRig{second, newApproveRigIn(t, a.r, "ep3", now), newApproveRigIn(t, a.r, "ep4", now)} {
+		if _, err := raw.Exec(`UPDATE mint_requests SET consumed_at = ?, approved_by = 'tty', credential_id = 'elsewhere' WHERE request_id = ?`, now-int64(600*(i+1)), r.id); err != nil {
+			t.Fatal(err)
+		}
+	}
+	pollSurface(t, b)
+	lines := func(want int, button string) {
+		t.Helper()
+		var got struct {
+			Lines  int
+			Button string
+		}
+		b.do(t, &got, "waitFor", map[string]any{"expression": fmt.Sprintf(`(() => { const d = document.querySelector('[data-decided]'), x = document.getElementById('decided-all');
+			return d && x && d.querySelectorAll('.nl.sub').length === %d && { lines: %d, button: x.textContent }; })()`, want, want)})
+		if got.Lines != want || got.Button != button {
+			t.Fatalf("Decided shows %d lines and %q, want %d and %q", got.Lines, got.Button, want, button)
+		}
+	}
+	lines(3, "Show all 4")
+	b.do(t, nil, "tap", map[string]any{"id": "decided-all"})
+	lines(4, "Show fewer")
+	b.do(t, nil, "tap", map[string]any{"id": "decided-all"})
+	lines(3, "Show all 4")
+	b.do(t, nil, "openCard", map[string]any{"url": "about:blank", "raw": true}) // ends the page's live stream before the server closes
 }
--- a/cmd/ailang-worldd/testdata/passkey_cdp.mjs
+++ b/cmd/ailang-worldd/testdata/passkey_cdp.mjs
@@ -269,9 +269,25 @@ async function waitFor({ expression, ms = 15000 }) {
   return v;
 }
 
+// key presses key ('n', 'N' (with Shift), 'Enter') at the focused element; with page, it waits for the new document. It
+// returns the focused element's id and the app's view (row 201 A3 sprint 4, the N key).
+const KEYS = { n: { code: 'KeyN', windowsVirtualKeyCode: 78, text: 'n' }, N: { code: 'KeyN', windowsVirtualKeyCode: 78, text: 'N', modifiers: 8 }, Enter: { code: 'Enter', windowsVirtualKeyCode: 13, text: '\r' } };
+async function key({ key, page = false }) {
+  const from = events.length;
+  await send('Input.dispatchKeyEvent', { type: 'keyDown', key, ...KEYS[key] });
+  await send('Input.dispatchKeyEvent', { type: 'keyUp', key, code: KEYS[key].code, windowsVirtualKeyCode: KEYS[key].windowsVirtualKeyCode });
+  if (page) await until(() => events.slice(from).some((e) => e.method === 'Page.loadEventFired'), 15000, 'the next page');
+  return evaluate(`({ location: location.href, active: document.activeElement ? document.activeElement.id : '', view: document.getElementById('stage') ? document.getElementById('stage').dataset.view : '' })`);
+}
+// tap clicks element id on any page.
+async function tap({ id }) {
+  if (!(await click(id))) throw new Error(id + ' is not visible');
+  return {};
+}
+
 const product = (await send('Browser.getVersion', {}, null)).product;
 console.log(JSON.stringify({ ready: true, chrome: product, pid: chrome.pid }));
-const commands = { tab, open, ceremony, uv, advance, fetchRewrite, openCard, press, waitFor, credential, fill, save, card: cardState };
+const commands = { tab, open, ceremony, uv, advance, fetchRewrite, openCard, press, waitFor, credential, fill, save, card: cardState, key, tap };
 for await (const line of readline.createInterface({ input: process.stdin })) {
   const { cmd, ...args } = JSON.parse(line);
   try {
--- a/host/surfacepage/app/app.js
+++ b/host/surfacepage/app/app.js
@@ -1190,7 +1190,7 @@
   }
 
   // ---------- Popovers ----------
-  const KEYROWS = [['← →', 'back · forward a page', '—'], ['⇧← ⇧→', 'previous · next state', '—'], ['↑ ↓', 'lane above · below', 'back in time · toward now'], ['+ −', 'zoom in · out', '—'], ['End', 'back to now', 'most recent'], ['Home', 'the oldest state', 'earliest'], ['Enter', 'open brane', '—'], ['Esc', '—', 'back to timeline'], ['/', 'address bar', 'address bar'], [',', 'appearance', 'appearance'], ['?', 'these keys', 'these keys']];
+  const KEYROWS = [['← →', 'back · forward a page', '—'], ['⇧← ⇧→', 'previous · next state', '—'], ['↑ ↓', 'lane above · below', 'back in time · toward now'], ['+ −', 'zoom in · out', '—'], ['End', 'back to now', 'most recent'], ['Home', 'the oldest state', 'earliest'], ['Enter', 'open brane', '—'], ['Esc', '—', 'back to timeline'], ['N', 'what needs you', 'what needs you'], ['/', 'address bar', 'address bar'], [',', 'appearance', 'appearance'], ['?', 'these keys', 'these keys']];
   function renderKeys() {
     $('keys').innerHTML = `<div class="kgrid"><span></span><span style="font-size:12px;font-weight:700;color:${st.d === 0 ? 'var(--text)' : 'var(--faint)'}">Timeline</span><span style="font-size:12px;font-weight:700;color:${st.d === 1 ? 'var(--text)' : 'var(--faint)'}">Brane</span>` +
       KEYROWS.map(([k, a, b]) => `<span><span class="kbd">${k}</span></span><span>${a}</span><span>${b}</span>`).join('') +
@@ -1415,10 +1415,12 @@
     if (e.target && e.target.tagName === 'INPUT') return;
     const k = e.key; let done = true;
     if (k === 'Enter' && e.target && e.target.id === 'home') return;   // the link's own activation: its click goes home
+    if (k === 'Enter' && e.target && e.target.closest && e.target.closest('[data-now]')) return;   // the now card's links and buttons activate themselves
     if (k === '?') toggle('keys');
     else if (k === ',') toggle('look');
     else if (k === '/') $('addr').focus();
     else if (k === 'Escape') { if (!$('keys').hidden || !$('look').hidden) { $('keys').hidden = true; $('look').hidden = true; } else out(); }
+    else if (k === 'n' || k === 'N') { const a = document.querySelector('[data-needs] a'); if (a) { if (st.d === 1) out(); a.focus(); } }   // design §2.8: N jumps to what needs you
     else if (st.d === 1 && k === 'ArrowUp') openRail(railIdx() - 1);
     else if (st.d === 1 && k === 'ArrowDown') openRail(railIdx() + 1);
     else if (st.d === 1 && k === 'Home') openRail(0);
````

### Patch M10b — XA4, the opt-in notification

````diff M10b
--- a/cmd/ailang-worldd/needs_you_notify_chrome_test.go
+++ b/cmd/ailang-worldd/needs_you_notify_chrome_test.go
@@ -0,0 +1,93 @@
+package main
+
+// Row 201 A3 sprint 4 (design §2.9, §6 XA4): the opt-in notification in a live Chrome, through the
+// driver's notify recorder (nothing is shown) and a CDP permission grant.
+
+import (
+	"reflect"
+	"testing"
+	"time"
+)
+
+// notifyView is what the notify recorder saw, the permission, and the opt-in's text.
+type notifyView struct {
+	Notes      []struct{ Title, Body, Tag string }
+	Asks       []struct{ Active bool }
+	Permission string
+	Control    string
+	Violations []string
+}
+
+// XA4: never asked, the page raises nothing and never requests permission, even for a new request.
+// The browser grants notifications and "Notify me" asks once, from the click (a user gesture). A
+// request the tab has already seen raises nothing; a new one raises exactly one "World needs you"
+// (tag = its id, the body names the episode, its grants and its expiry), and the next refresh none.
+// Reloaded with notifications granted, the page raises nothing for what it shows; a new request's
+// notification, clicked, focuses the tab and opens its card.
+func TestNeedsYouNotifiesOnlyWhenAskedOnce(t *testing.T) {
+	b := startPasskeyBrowser(t)
+	b.do(t, nil, "tab", map[string]any{"inject": "notify"})
+	now := time.Now().Unix()
+	a := newApproveRig(t, "ep1", now)
+	origin := rigSurface(t, a.credPath)
+	seen := func(notes, asks int, permission, control string) notifyView {
+		t.Helper()
+		var v notifyView
+		b.do(t, &v, "waitFor", map[string]any{"expression": `(() => { const c = document.querySelector('[data-needs] #notify-me, [data-needs] [data-notify]');
+			return { notes: window.__notes, asks: window.__asks, permission: Notification.permission, control: c ? c.textContent : '', violations: window.__cspv }; })()`})
+		if len(v.Notes) != notes || len(v.Asks) != asks || v.Permission != permission || v.Control != control || len(v.Violations) != 0 {
+			t.Fatalf("the page made %d notifications %+v and %d permission requests, permission %q, control %q, CSP violations %q; want %d, %d, %q, %q, none",
+				len(v.Notes), v.Notes, len(v.Asks), v.Permission, v.Control, v.Violations, notes, asks, permission, control)
+		}
+		return v
+	}
+	poll := func(prefix string) {
+		t.Helper()
+		pollSurface(t, b)
+		sectionText(t, b, "data-needs", prefix)
+	}
+	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface", "raw": true})
+	sectionText(t, b, "data-needs", "Needs you · 1")
+	seen(0, 0, "default", "Notify me")
+	newApproveRigIn(t, a.r, "ep2", now)
+	poll("Needs you · 2")
+	seen(0, 0, "default", "Notify me")
+
+	var permission string
+	if b.do(t, &permission, "grantNotifications", map[string]any{"origin": origin}); permission != "granted" {
+		t.Fatalf("after the grant the permission is %q", permission)
+	}
+	b.do(t, nil, "tap", map[string]any{"id": "notify-me"})
+	b.do(t, nil, "waitFor", map[string]any{"expression": `document.querySelector('[data-needs] [data-notify]') !== null`})
+	if v := seen(0, 1, "granted", "Notifications on"); !v.Asks[0].Active {
+		t.Fatal("permission was requested without a user gesture")
+	}
+	poll("Needs you · 2")
+	seen(0, 1, "granted", "Notifications on")
+	third := newApproveRigIn(t, a.r, "ep3", now)
+	poll("Needs you · 3")
+	v := seen(1, 1, "granted", "Notifications on")
+	want := struct{ Title, Body, Tag string }{"World needs you", "Approve mint request for ep3 (10 grants), expires " + time.Unix(third.row(t).ExpiresAt, 0).UTC().Format("15:04"), third.id}
+	if !reflect.DeepEqual(v.Notes[0], want) {
+		t.Fatalf("the notification is %+v, want %+v", v.Notes[0], want)
+	}
+	poll("Needs you · 3")
+	seen(1, 1, "granted", "Notifications on")
+
+	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface", "raw": true})
+	sectionText(t, b, "data-needs", "Needs you · 3")
+	seen(0, 0, "granted", "Notifications on")
+	fourth := newApproveRigIn(t, a.r, "ep4", now)
+	poll("Needs you · 4")
+	if v := seen(1, 0, "granted", "Notifications on"); v.Notes[0].Tag != fourth.id {
+		t.Fatalf("the notification is %+v, want tag %s", v.Notes[0], fourth.id)
+	}
+	var clicked struct {
+		Location string
+		Focused  int
+	}
+	if b.do(t, &clicked, "notification", map[string]any{"index": 0}); clicked.Location != origin+"/surface/approve/"+fourth.id || clicked.Focused != 1 {
+		t.Fatalf("the notification's click: %+v", clicked)
+	}
+	b.do(t, nil, "waitFor", map[string]any{"expression": `document.body.dataset.state === 'ready'`})
+}
--- a/cmd/ailang-worldd/testdata/passkey_cdp.mjs
+++ b/cmd/ailang-worldd/testdata/passkey_cdp.mjs
@@ -70,6 +70,17 @@ const injections = {
     const b = (x) => btoa(String.fromCharCode(...new Uint8Array(x.buffer || x))).replace(/\\+/g, '-').replace(/\\//g, '_').replace(/=+$/, '');
     navigator.credentials.get = (o) => { const k = o.publicKey; window.__gets.push({ challenge: b(k.challenge), rpId: k.rpId, allow: (k.allowCredentials || []).map((c) => b(c.id)),
       userVerification: k.userVerification, timeout: k.timeout, signal: !!o.signal }); return get(o); }; }`,
+  // notify (XA4): a recorder in place of window.Notification. It shows nothing (no notification is
+  // raised, even in this headless profile): it records each notification the page makes in
+  // window.__notes and each permission request, with whether a user gesture was active, in
+  // window.__asks, and passes the permission and its request through to the browser. The page's
+  // window.focus calls are counted in sessionStorage, which outlives a navigation.
+  notify: `window.__notes = []; window.__asks = []; if (window.Notification) { const N = window.Notification;
+    const R = function (title, o) { const n = { title, body: o && o.body, tag: o && o.tag, onclick: null }; window.__notes.push(n); return n; };
+    Object.defineProperty(R, 'permission', { get: () => N.permission });
+    R.requestPermission = () => { window.__asks.push({ active: navigator.userActivation.isActive }); return N.requestPermission(); };
+    window.Notification = R; }
+    { const focus = window.focus.bind(window); window.focus = () => { sessionStorage.setItem('__focus', String(+(sessionStorage.getItem('__focus') || 0) + 1)); focus(); }; }`,
   eddsa: `{ const create = navigator.credentials.create.bind(navigator.credentials);
     navigator.credentials.create = (o) => { o.publicKey.pubKeyCredParams = o.publicKey.pubKeyCredParams.filter((p) => p.alg === -8); return create(o); }; }`,
 };
@@ -279,6 +290,20 @@ async function key({ key, page = false }) {
   if (page) await until(() => events.slice(from).some((e) => e.method === 'Page.loadEventFired'), 15000, 'the next page');
   return evaluate(`({ location: location.href, active: document.activeElement ? document.activeElement.id : '', view: document.getElementById('stage') ? document.getElementById('stage').dataset.view : '' })`);
 }
+// grantNotifications grants origin notifications, as the browser does after its own prompt (XA4), and
+// returns the page's permission.
+async function grantNotifications({ origin }) {
+  await send('Browser.grantPermissions', { permissions: ['notifications'], origin }, null);
+  return evaluate('Notification.permission');
+}
+// notification clicks the index'th recorded notification and waits for the new document; it returns the
+// location and how many times the page focused its window.
+async function notification({ index }) {
+  const from = events.length;
+  await evaluate(`(window.__notes[${index}].onclick(), 1)`);
+  await until(() => events.slice(from).some((e) => e.method === 'Page.loadEventFired'), 15000, 'the next page');
+  return evaluate(`({ location: location.href, focused: +(sessionStorage.getItem('__focus') || 0) })`);
+}
 // tap clicks element id on any page.
 async function tap({ id }) {
   if (!(await click(id))) throw new Error(id + ' is not visible');
@@ -287,7 +312,7 @@ async function tap({ id }) {
 
 const product = (await send('Browser.getVersion', {}, null)).product;
 console.log(JSON.stringify({ ready: true, chrome: product, pid: chrome.pid }));
-const commands = { tab, open, ceremony, uv, advance, fetchRewrite, openCard, press, waitFor, credential, fill, save, card: cardState, key, tap };
+const commands = { tab, open, ceremony, uv, advance, fetchRewrite, openCard, press, waitFor, credential, fill, save, card: cardState, key, tap, grantNotifications, notification };
 for await (const line of readline.createInterface({ input: process.stdin })) {
   const { cmd, ...args } = JSON.parse(line);
   try {
````

### Patch M10c — the card shows a hostile value as text; QUICKSTART §17

````diff M10c
--- a/cmd/ailang-worldd/approve_page_hostile_test.go
+++ b/cmd/ailang-worldd/approve_page_hostile_test.go
@@ -0,0 +1,32 @@
+package main
+
+// Row 201 A3 sprint 4 (iteration 386 judge survivors C03/C04): a hostile signed value on the request
+// card, in a live Chrome, is shown as text.
+
+import (
+	"testing"
+	"time"
+)
+
+// A request whose model is markup and an attribute breaker: the card shows it verbatim as the
+// model field's text, ready to approve; nothing it carries becomes an element or runs.
+func TestApprovePageShowsAHostileValueAsText(t *testing.T) {
+	const hostile = `<img src=x onerror="window.__xss=1"><script>window.__xss=2</script>"&amp;`
+	b := startPasskeyBrowser(t)
+	b.do(t, nil, "tab", nil)
+	a := newApproveRig(t, "ep1", time.Now().Unix(), "--model", hostile)
+	origin := cardOrigin(t, a.credPath, "", nil)
+	var view cardView
+	if b.do(t, &view, "openCard", map[string]any{"url": origin + "/surface/approve/" + a.id}); view.State != "ready" || view.Fields["model"] != hostile ||
+		len(view.Errors) != 0 || len(view.Violations) != 0 {
+		t.Fatalf("the card shows model %q, state %q, errors %q, CSP violations %q", view.Fields["model"], view.State, view.Errors, view.Violations)
+	}
+	var page struct {
+		XSS             bool
+		Images, Scripts int
+	}
+	b.do(t, &page, "waitFor", map[string]any{"expression": `({ xss: window.__xss !== undefined, images: document.querySelectorAll('img').length, scripts: document.scripts.length })`})
+	if page.XSS || page.Images != 0 || page.Scripts != 1 {
+		t.Fatalf("the hostile value became markup: %+v", page)
+	}
+}
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -1847,7 +1847,9 @@ ailang-worldd episode start surface-demo --repo "$PWD" --db /tmp/se-world/world.
 ailang-worldd approve mr-<the 16 hex it printed> --await --db /tmp/se-world/world.db &
 ```
 
-`--await` prints the request's id, the first 16 hex of its message hash and the card's address, `open http://localhost:7644/surface/approve/mr-<id>`. Open `http://localhost:7644/surface`: **Needs you · 1** lists `surface-demo` with its agent, grant count and minutes to expiry. Follow it to the card, check every field and that the message hash matches the one `--await` printed, press **Approve with passkey** and pass your device's check. The card shows `Approved by webauthn:demo-mac`; the waiting `approve` mints and prints `approved by webauthn:demo-mac`, and `approve --list --decided` shows it.
+`--await` prints the request's id, the first 16 hex of its message hash and the card's address, `open http://localhost:7644/surface/approve/mr-<id>`. Open `http://localhost:7644/surface`: **Needs you · 1** lists `surface-demo` with its agent, grant count and minutes to expiry. Follow it to the card, check every field and that the message hash matches the one `--await` printed, press **Approve with passkey** and pass your device's check. The card shows `Approved by webauthn:demo-mac`; the waiting `approve` mints and prints `approved by webauthn:demo-mac`, and `approve --list --decided` shows it. Back on `/surface`, **Decided** under Needs you shows `You · Minted · surface-demo · webauthn:demo-mac · <hh:mm>` (UTC, the app's clock): the last three, and **Show all** for the rest of the daemon's 20. **N** jumps to the first Needs-you link; Enter opens its card.
+
+To be told about the next request, press **Notify me** under Needs you and allow notifications when the browser asks; the page never asks on its own. From then on the tab raises one *World needs you* notification ("Approve mint request for <episode> (<n> grants), expires <hh:mm>") for each new request it reads on the app's own refresh, never for one it already showed; clicking it focuses the tab and opens that request's card. A notification carries no authority: the card and the passkey still decide. Whether the browser and the system display it is theirs to decide; check it once in each browser you use.
 
 ```bash
 wait
````
