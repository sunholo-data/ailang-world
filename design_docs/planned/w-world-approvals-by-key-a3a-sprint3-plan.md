# Sprint plan — w-world-approvals-by-key-a3a, sprint 3 (iteration 386, row 201)

> **THE SPRINT — A3a sprint 3, "the page and the app": M6 the request card `GET /surface/approve/{id}` +
> `GET /surface/approve.js` (PG2, PG3, PG3b, the card's network deadlines), M7's *Needs you* list in `app.js`
> (XA3's list half), M8 the live-Chrome tests XA1, XA2 (arms a–c), XA3's list half, XA5 on the A2 driver
> (extended) and QUICKSTART §17 *(attended — pending first verbatim run)*, as the A3 child design r2 §2.7–§2.10,
> §6 "Sprint 3" and §9 list them — plus both carried obligations: sprint 2's Appendix C (R11's `sshsig` arm,
> RT6, RT13, AW6's route half) and iteration 385's judge survivor J36.**
> - **Measured: 1,843 changed lines under `cmd/ host/`** — `git diff --numstat <base> HEAD -- cmd host`,
>   insertions + deletions, after a commit of each of the nine scratch stages (463 production: `approve.js`
>   + `index.html` + embed + page handlers + routes 411, the SSHSIG arm 27, `app.js` 25;
>   1,380 test, incl. the driver `passkey_cdp.mjs`). Margin to the 2,000 cap: **157**. `docs/QUICKSTART.md`
>   §17 is 18 lines more (outside the filter).
> - **The design's fallback cut fires, and only its first step** (§1.3 S3-D1, §1.5): M7's notification and the
>   Decided list (XA4, XA3's Decided half, the "Notify me" control, the `N` key) move to a fourth sprint
>   (Appendix C). Measured, the sprint as built is 1,843; the moved part is ≈ 300 lines at this tree's
>   density (≈ 135 `app.js` lines left of the design's ≈ 160 M7 estimate after the 25-line Needs-you list, XA4
>   ≈ 115 like XA2, the driver's `grantNotifications`/`notify` ≈ 30, XA3's Decided half ≈ 20), which would put
>   the sprint at ≈ 2,143, over the cap. **Appendix C is NOT cut**: with it the sprint is 1,843 ≤ 2,000.
>   Never cut, all present: PG2, PG3, PG3b, XA1, XA2 incl. arm (c), XA5, the J36 test.
> - **Safety order is the design's**: the J36 test (M6a) and the POST's SSHSIG arm (M5d) land before any page;
>   the SSHSIG arm **judges with the A1b law before it records**, exactly as the passkey arm (one shared R11 →
>   R12 → R13 → R14 path; `MUT-SSH-RECORD-BEFORE-LAW` and `MUT-SSH-SKIP-LAW` killed); the card page and its script
>   register **through `approvalRoute`** (the Host check, RT1/RT7 extended); `surfaceCSP` is unchanged (PG2 pins
>   its exact text); the page **verifies the digest itself, renders every signed field from the bytes it hashed,
>   refuses any field it cannot show, and computes the 90-byte message and its hash itself** — it never takes a
>   hash from the server (XA1: its `get` challenge equals the A2 page's for the same request; the SSH arm's saved
>   message equals Go's `approval.SignedMessage`).
> - **Every line of code, test, JS, HTML and doc is a unified diff in Appendix P**, generated LAST from
>   compiled, vetted, gofmt-clean, green scratch stages; the patches were then extracted **from this file** with
>   §2's `awk`, applied to a fresh `git archive 82fb9afa` export, and the result compared byte-for-byte with the
>   mutant tree (§1.2 C30). The executor applies each patch and runs its acceptance; it writes no code by hand.

- **Design**: [w-world-approvals-by-key-a3.md](w-world-approvals-by-key-a3.md), **r2** (quorum r1+r2 +
  controller carve-out + independent audit VALID 84; **not reopened**). D-WORLD-114 (no Tauri half in A3a; the
  CLI minter re-judges and mints) is OPEN with default A = keep both, applied since iteration 384: built as
  designed. Where this plan and the design differ, this plan was measured against `82fb9afa` and wins; each
  difference is a decision (§1.3) or a premise found false (§1.4).
- **Starting ground**: A3a sprints 1 and 2 merged (`b44b5f7d` PR #554; `d1027477` PR #556, recorded `9ad741bb`):
  schema 7 `surface_approvals`, the shared judge, the minter `approve <id> --await`, the read routes
  `GET /surface/requests.json` / `GET /surface/requests/{id}`, the surface origin set `O`, `serve
  --surface-origin`, and `POST /surface/approve/{id}` for passkeys. The sprint-2 plan's §1.4 P1–P7 and its
  S2-D1…S2-D17 are facts of HEAD this plan builds on (S2-D3: every approval route goes through
  `d.approvalRoute`, whose handler signature makes an unwrapped one a compile error).
- **Planner base**: `82fb9afa9086eed2f781bb0d4fe54d4bae3b70d7` (= origin/dev). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter386-a3a-s3`, branch `sprint/world-iter386-row201-a3a-s3`.
  Planner scratch (outside the repo, never under `/tmp`): `$HOME/.ailang/state/world-iter386/plan-scratch/`
  (`base/` the pristine export, `work/` the stage repo (tags `M6a` … `M8d`) and the tree the mutants were cut
  from, `forkA/` the planner's fork that built M6a/M5d, `patches/`, `tools/`, `logs/`, `roundtrip/`, `mut/`).
- **Executor**: as routed by the controller. **Three runs**: Run 1 (M6a, M5d, M6b), Run 2 (M6c, M7a, M8a),
  Run 3 (M8b, M8c, M8d). Each sub-milestone is one `patch -p1` (§2) plus its acceptance (§3). **Apply, never
  retype.** This sprint's acceptance runs **Chrome headless** (temporary profiles, a CDP virtual authenticator,
  test keys only) from M8a on.
- **Touched** (§5): new `host/daemon/surface_approve_card_error_test.go` (M6a); `host/daemon/surface_approve.go`
  (the SSHSIG arm), `surface_approve_post_test.go` (RT4's cut arm → three SSHSIG law arms), new
  `host/daemon/surface_approve_sshsig_test.go`, new `cmd/ailang-worldd/approve_surface_sshsig_e2e_test.go` (M5d);
  `host/daemon/daemon.go` (two `HandleFunc` + the doc comment), new `host/daemon/{surface_approve_page.go,
  surface_approve_page_test.go}`, `host/daemon/{surface_approve_test.go, surface_app_test.go}` (RT1/RT7 and the
  404 list), new `host/surfacepage/{approve.go, approve/index.html, approve/approve.js}` (M6b); new
  `host/surfacepage/approve_page_test.go` (M6c); `host/surfacepage/app/app.js`, `host/surfacepage/app_test.go`
  (M7a); `cmd/ailang-worldd/testdata/passkey_cdp.mjs`, new `cmd/ailang-worldd/{approve_page_driver_test.go,
  approve_page_chrome_test.go}` (M8a); new `cmd/ailang-worldd/approve_page_swap_test.go` (M8b); new
  `cmd/ailang-worldd/approve_card_deadline_test.go` (M8c); new `cmd/ailang-worldd/needs_you_chrome_test.go`,
  `docs/QUICKSTART.md` §17 (M8d). **Outside the design's §8 file list**: `surface_approve_card_error_test.go`
  (J36), `surface_approve.go`/`surface_approve_post_test.go`/`surface_approve_sshsig_test.go`/
  `approve_surface_sshsig_e2e_test.go` (Appendix C, a sprint-2 surface), `surface_approve_page.go` (the page
  handlers, kept out of the POST's file), `host/surfacepage/app_test.go` (the node harness's export list and
  `needsYouHTML` assertions), and the design's "`approve_page_*_test.go` (new)" split by sub-milestone
  (`approve_page_driver_test.go`, `approve_page_chrome_test.go`, `approve_page_swap_test.go`,
  `approve_card_deadline_test.go`, `needs_you_chrome_test.go`); `passkey_page_driver_test.go` is **not** edited
  (the design listed it as "extended": the card's helpers live in their own file).
- **Not touched**: every `host/approval` file (the law stays pure and standard-library-only), `host/store`,
  `host/authority`, `surfaceCSP`, `isProtected`, `/surface/passkey`, `passkey.js`, `passkey/index.html` (PG2 pins
  their sha256), `host/surfacepage/model.go` (`needsYou` stays `[]`), every `/v1` route, `world/*.ail`,
  `packages/**`, `scripts/**`, `tools/launchd/*`, `cmd/world-publish/*`, any skill, `~/.ssh`, `~/.ailang/world`,
  `~/.ailang/state/mission-v1*`.

---

## 0. Binding constraints (each is restated where it governs)

1. **Test keys, temp stores and virtual authenticators only (D-WORLD-109 = A).** Every store is under
   `t.TempDir()`; every SSH key is `approvaltest.NewSigner`'s fixed seed; every passkey is either created by
   Chrome's CDP **virtual** authenticator in a temporary profile (XA1, through the A2 page and the CLI) or is
   `approvaltest.NewPasskey`'s fixed ES256 seed handed to that virtual authenticator with
   `WebAuthn.addCredential` (S3-D11). The daemon is served by `httptest.NewServer` on the loopback (as A2's
   `passkeyPageOrigin`); Chrome runs `--headless=new` with `--user-data-dir` under the test's temp dir, killed and
   removed at the end. Nothing opens `~/.ailang/world`, `~/.ssh`, a real credentials store, a real browser
   profile, a real `/dev/tty`, or runs `ssh-keygen`; nothing mints into a real store or publishes.
2. **Pure core, effects at the boundary (coding standards S1–S3).** The page is static embedded bytes; its
   field table is tied to the Go struct by a pure reflection test (PG3, PG3b); the Needs-you renderer is a pure
   exported function of `app.js` tested under node (`TestSurfaceAppPure`, which CI runs with
   `WORLD_EXEC_NODE`). Effects (HTTP, the stores, the browser) stay in the daemon, the minter and the tests.
   `host/approval` is not touched. No `.ail` (design §3), no new Go module, package, production exec site, `/v1`
   route or context root.
3. **The daemon never keeps the credentials writer; it never mints.** The SSHSIG arm reuses the passkey arm's
   R12–R14 path unchanged: one non-waiting `store.OpenCredentials` after the law, closed before the answer
   (`TestDaemonCredentialWriterSites` unchanged and green). No CORS header anywhere; `isProtected` is unchanged;
   no cookie or session is read; `approve.js` has no `innerHTML`, `eval(`, `document.cookie` (PG2's labelled
   controls).
4. **One mutant per new production branch, anchored to the diff — Go AND JS** (§4), each run against its single
   named test **and** its package suite, inventories excluded from every suite; named equivalents with reasons;
   every run appends a ledger line.
5. **No `t.Parallel()`** in any new test (`TestNoParallelWriteForkPackagesOutsideVerifygate`).
6. **Bisectable and safe at every commit.** Each sub-milestone ends gofmt-clean, `go build ./...` and `go vet`
   green, the `_test.go` compile fence (`go test -run '^$' ./...`) green, its own tests and every earlier one
   green. The SSHSIG arm lands before the page that posts to it; the page lands behind the Host check.
7. **No git writes by the executor** (§7). The controller commits from snapshots.

---

## 1. Pre-flight

### 1.1 Baseline at `82fb9afa` (measured by the planner, 2026-11-12)

Every command ran with `export AILANG_BIN=$HOME/.pinned-ailang/ailang`, `go1.26.6 darwin/arm64`, node
`v26.13.0`, `Google Chrome 157.0.8103.51`, in a `git archive 82fb9afa` export at `plan-scratch/base` (a throwaway
git repo; `diff -r cmd host` against the worktree: identical, `BASE_EQUALS_WORKTREE`; the worktree itself was
only read). The pre-commit legs are `tools/gates.sh <tree> <logs> base` = §6 without the sprint's new tests
(they do not exist at the base). Logs: `plan-scratch/logs/baseline/`.

| Command | rc | Time | Result |
|---|---|---|---|
| `git -C <worktree> status --porcelain \| wc -l` | 0 | — | `0` (clean at `82fb9afa`) |
| `./scripts/verify_ail.sh` | **0** | — | `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` |
| `gofmt -l cmd host`; `go build ./...`; `go vet ./...`; `go test -run '^$' ./...` | **0** | — | empty / clean / empty / no FAIL line |
| `go test -count=1 -timeout 20m ./host/surfacepage/ ./host/store/ ./host/authority/ ./host/approval/... ./cmd/ailang-worldd/ ./host/daemon/` | **0** | — | all `ok` (`host/surfacepage` 0.56 s, `host/store` 12.0 s, `host/authority` 0.85 s, `cmd/ailang-worldd` 127.9 s, `host/daemon` 277.9 s) |
| `-race` legs (§6's three, at the base's tests of the same patterns) | **0** | — | `ok host/surfacepage` 2.9 s, `ok host/daemon` 4.1 s, `ok cmd/ailang-worldd` 20.7 s; no race report |
| §6's named inventory, schema, reflective, help, quickstart and route tests, `-v` | **0** | — | 50 `--- PASS`, 0 SKIP, 0 FAIL; the AW leg 9 `--- PASS: TestApproveAwait` |
| the live-Chrome leg (`WORLD_CHROME=… WORLD_REQUIRE_PASSKEY_BROWSER=1`, the A2 browser tests) | **0** | — | `--- PASS: TestPasskeyPageAbortsAPendingCeremony (4.21s)`, `--- PASS: TestPasskeyPageInChrome (3.36s)` |
| the whole pre-commit script (`tools/gates.sh base`) | **0** | **420 s** | `PRE_COMMIT_GATES_GREEN 420s` |
| `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` (the base export, committed) | **0** | — | `✓ 0 other binary blobs among 6791 tracked files`; `✓ toolchain floor gate: go1.26.6 >= root module floor go1.26.6`; `✓ all 64 required top-level evidence tests passed exactly once`; `✓ go gate PASSED (plain)` (nothing else running; `logs/base_verify_go.log`) |

**Rule (carried):** a red at the base measures the repo, not the sprint. The controller re-runs a red
untouched package alone once and does not chase it. **Never start a gate under `nohup`.**
`VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` reads committed blobs and runs `go test ./...`; it is the
after-commit leg (§6).

### 1.2 Codebase claims checked (at `82fb9afa` unless "work" — the planner's final scratch tree)

| # | Claim | Command → observed output |
|---|---|---|
| C1 | The base | `git rev-parse HEAD origin/dev` (worktree) → `82fb9afa9086eed2f781bb0d4fe54d4bae3b70d7` twice; `git status --porcelain \| wc -l` → `0` |
| C2 | **Nothing of sprint 3 has landed** | `grep -rn 'approve\.js\|"GET /surface/approve' --include='*.go' cmd host \| grep -v _test \| wc -l` → `0`; control in the same call: `grep -rn '/surface/passkey' --include='*.go' cmd host \| grep -v _test \| wc -l` → `11`; `ls host/surfacepage/` → no `approve/`, no `approve.go` |
| C3 | **Design drift since `0140bd7f`** | `shasum -a 256` of design V1's twelve files: `surface_page.go f769ee06b047f113`, `surfacepage/model.go 4ace112b84f959b4`, **`app.js 8ba97c83d7e78a69`**, **`passkey.js 0064c6c90099a356`**, `session_new.go 587f39f99ace013d`, **`passkey_cdp.mjs 2a336e5e810f8ed6`** equal the design's; `daemon.go a8ff6175f1c9584d`, `approve.go 6a38843161c7d75f`, `approve_passkey.go e17e59d114ef5648`, `approve_sig.go c50db2326f6a52ef`, `mint_requests.go a1501959398c48eb`, `credentials.go 4273fad7537d9905` differ (sprints 1–2: `git diff --stat 0140bd7f 82fb9afa -- cmd host` → `41 files changed, 3458 insertions(+), 80 deletions(-)`). So every `app.js`, `passkey.js` and driver line the design cites is current; `daemon.go`'s route block moved (C4) |
| C4 | The route block (design §1.6) | `grep -n 'approvalRoute(d.handle\|^func (d \*Daemon) Handler\|return r.Method == http.MethodPost' host/daemon/daemon.go` → `1104:func (d *Daemon) Handler()`, `1135`–`1137` the three approval routes (`GET /surface/requests.json`, `GET /surface/requests/{id}`, `POST /surface/approve/{id}`, each `d.approvalRoute(...)`), `1158: return r.Method == http.MethodPost && r.URL.Path == "/v1/commit"` |
| C5 | **The sprint-2 SSHSIG refusal** (carried obligation 1) | `grep -n 'if body.SSHSIG != nil\|authority.JudgeAssertion(' host/daemon/surface_approve.go` → `255: if body.SSHSIG != nil { // R11, sprint 2's fallback cut: the SSHSIG arm is sprint 3's`, `260: signed, err := authority.JudgeAssertion(ctx, ro, row, body.WebAuthn, now) // R11: the law`; RT4's arm `R11-sshsig-cut` at `surface_approve_post_test.go:131` |
| C6 | **J36's line** (carried obligation 2) | `grep -n 'active, err := ro.ActivePasskeys\|if err != nil && !isOlderStore(err) {' host/daemon/surface_approve.go` → `163: active, err := ro.ActivePasskeys(ctx)`, `164: if err != nil && !isOlderStore(err) {` (in `handleSurfaceRequestCard`; `123:` is `readRequests`' other check) |
| C7 | **J36 was a survivor at the base, and dies to M6a's test** (fork A, `logs/j36_proof.txt`) | base + `164: if false {` → `go test ./host/daemon/ -run '^(TestApprov\|TestRequest\|TestSurface)' -v` → rc 0, 29 `--- PASS` incl. `--- PASS: TestRequestCardJSON` (the survivor); M6a's tree + the same edit → `--- FAIL: TestRequestCardPasskeyReadError` with `surface_approve_card_error_test.go:23: a passkey read error: 200 {"v":1,"id":"mr-00000000000000e1",…,"passkeys":[],…}`; `go test -vet=off ./host/daemon/ -run '^(TestApprov\|TestRequest)'` (inventories excluded) → `FAIL`, only that test; restored (`cmp` equal) → `--- PASS: TestRequestCardPasskeyReadError (0.02s)` |
| C8 | **The J36 injection must leave the open working** (§1.4 P5) | work, fork A run 1: `DROP TABLE webauthn_credentials` (version kept 7) → the card GET's 500 comes from `OpenCredentialsReadOnly` (its table-set check), so the mutant **survived** that test; shipped: `ALTER TABLE webauthn_credentials RENAME COLUMN origin TO origin_gone` → the open and `MintRequest` succeed and `ActivePasskeys` fails with a non-`OlderCredentialsStoreError` |
| C9 | `surfaceCSP` and the headers | `host/daemon/surface_page.go:19 const surfaceCSP = "default-src 'none'; style-src 'unsafe-inline'; script-src 'self'; connect-src 'self'; font-src 'self'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'"`; `setSurfaceHeaders` (`:21`) = `setWorkbenchHeaders` + the CSP; unchanged on the final tree (`git diff base M8d -- host/daemon/surface_page.go host/daemon/workbench.go \| wc -l` → `0`, work) |
| C10 | **The canonical request has 21 keys and no `omitempty`** (design V24, V28) | `grep -c omitempty host/approval/request.go` → `0`; `sed -n '/^type MintRequest struct/,/^}/p' host/approval/request.go \| grep -c 'json:"'` → `21`; `Grant` (`effect scope budget`), `ProfileRef` (`project digest`), `DaemonRef` (`addr vcs_revision vcs_modified executable_sha256 go_version`): 28 leaf paths (work: PG3 `--- PASS`) |
| C11 | **The tests that pin the sprint-2 route set** | `grep -n '"/surface/approve.js", "/surface/approve"' host/daemon/surface_app_test.go host/daemon/surface_approve_test.go` → `surface_app_test.go:156` (`TestSurfaceAppRoutes`' 404 list) and `surface_approve_test.go:53` (RT1's "not served this sprint" 404 list); `surface_approve_test.go:20 var approvalRoutes = []string{"GET /surface/requests.json", "GET /surface/requests/{id}", "POST /surface/approve/{id}"}` — the three edits S2-D15 named for this sprint |
| C12 | Route-enumerating tests (sprint 2 C14, re-checked) | `TestFrozenV1RouteTableMatchesMux` (`/v1/` only, 13), `TestSurfaceAppRoutes`, `TestApprovalRoutesRegistered` (a regex over `Handler`'s body: every `/surface/(requests\|approve)` registration must be `d.approvalRoute(`), `TestMCPMountMethodSource` (AST, `d.projection.*` only), `TestReadCtxCancelledAfterHandler` (world-store reads); work: all `--- PASS` with the two new routes |
| C13 | **No test pins QUICKSTART §14–§17** | `grep -n '^## 1[4-7]\.' docs/QUICKSTART.md` → `1773: ## 14.`, `1785: ## 15.`, `1811: ## 16. Approve with a passkey *(attended — pending first verbatim run)*` (the last section); `grep -rln 'Section1[4-7]\|## 1[4-7]\.\|Approve with a passkey\|Mint by request' --include='*_test.go' cmd host` → no file; the quickstart tests pin §9–§11 and the service section (work: all `--- PASS`, §6) |
| C14 | **`app.js`'s Needs you has no source** (design §1.8) | `wc -l < host/surfacepage/app/app.js` → `1413`; `app.js:583 … needs = ne.needsYou \|\| [] …`, `:605 Needs you · ${needs.length}`; `:885 async function loadModel()`, `:970 async function refresh(trigger)` (the existing refresh: live entries and follow mode's poll); `:1395 window.World = {…`, `:1396 poll: () => refresh('poll'), …` (the screenshot driver's seam, which XA3 uses) |
| C15 | **The `N` key is the static pages', not the app's** (§1.4 P7) | `sed -n 255,259p host/surfacepage/surface.js` → the `KEYS` help table `["N", "jump to what needs you"]` (served at `/surface/surface.js` for `/surface/static`, `daemon.go:1125`); `grep -n "'N'" host/surfacepage/app/app.js` → no match |
| C16 | **The A2 harness** (design §1.10) | `cmd/ailang-worldd/testdata/passkey_cdp.mjs` 172 lines, `:159 const commands = { tab, open, ceremony, uv, advance };`, injections `abort`, `late`, `eddsa`; `click(id)` used the element's bounding box (`:46`); `passkey_page_driver_test.go` `startPasskeyBrowser` (gated, `SKIP COUNTED` unless `WORLD_REQUIRE_PASSKEY_BROWSER=1`), `passkeyBrowserBudget = 120 s` per driver |
| C17 | **A wrapped inline link's bounding-box centre can miss it** (S3-D10) | work: the Needs-you `<a>` wraps to two line boxes; its `getBoundingClientRect()` centre hit `<div class="nl">` (`document.elementFromPoint` → `<div class="nl"><a id="needs-…`), so `click` (bounding box) did not navigate (`driver press: timed out waiting for the next page`); with `getClientRects()[0]` it does; A2's XB1/XB2 stay green with the change (`--- PASS: TestPasskeyPageInChrome`, `--- PASS: TestPasskeyPageAbortsAPendingCeremony`, work) |
| C18 | **Decoding into a reused struct keeps old map keys** (a test-harness trap, S3-D12) | work: XA2 decoded each arm's view into one `cardView`; a refused card (`"fields":{}`) still showed the previous arm's 28 fields; the tests decode each view into a fresh value |
| C19 | **CI runs node but not Chrome** (design R3) | `grep -cE 'WORLD_CHROME\|google-chrome\|chrome' .github/workflows/ci.yml` → `0`; `ci.yml:250 echo "WORLD_EXEC_NODE=$(command -v node)" >> "$GITHUB_ENV"`, `:289 : "${WORLD_EXEC_NODE:?WORLD_EXEC_NODE is required}"` — so `TestSurfaceAppPure` (node) runs in CI, the XA tests are counted SKIPs there |
| C20 | The store APIs the J36 injection and the SSHSIG tests use | `host/store`: `ActivePasskeys`, `EnrolPasskey`, `EnrolApprover`/`seedApprover` (`cmd/ailang-worldd/approver_rig_test.go:93`), `SurfaceApproval`; `host/approval/sshsig.go:19 const Namespace = "world-approval@sunholo"`; `request.go:211 func SignedMessage(digest, challenge [32]byte, expiry int64) ([]byte, error)` |
| C21 | `approvaltest`'s passkey key is reproducible from its seed (S3-D11) | `host/approval/approvaltest/passkey.go:30-46`: `raw := bytes.Repeat([]byte{seed}, 32)`; ES256 → `ecdsa.ParseRawPrivateKey(elliptic.P256(), raw)`; credential ID `sha256("approvaltest-credential" ‖ alg ‖ seed)`; work: `virtualPasskey`'s control (the PKCS#8 key's public X equals `pk.Key().X`) holds, and Chrome's assertions verify under the daemon's law (XA2 control, XA5, `--- PASS`) |
| C22 | **Every repo-wide inventory, on the final tree** | the §6 script (`tools/gates.sh <tree> <logs> final` = §6's pre-commit block) on the round-trip tree (byte-identical to the mutant tree, C30) → `PRE_COMMIT_GATES_GREEN 473s`; pkgs: `host/surfacepage` ok 0.862s, `host/store` ok 12.390s, `host/authority` ok 1.053s, `host/approval` ok 0.840s, `host/approval/approvaltest` ok 1.272s, `cmd/ailang-worldd` ok 118.284s, `host/daemon` ok 278.434s (`logs/final_gates/`): design §7's inventories by name — `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` (no exec site added), `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist`, `TestBareNetHTTPExemptionIsPerGroup`, `TestEveryCommandDirectoryIsAProtectedGroup`, `TestProductionContextRoots` (no root added), `TestProductionGoSurface`, `TestNoNewDeadlineFreeStoreReads`, `TestReadGettersHonorContext`, `TestCredentialsMethodsRequireADeadline`, `TestFrozenV1RouteTableMatchesMux` (13), `TestSurfaceAppRoutes` (edited), `TestMCPMountMethodSource`, `TestReadCtxCancelledAfterHandler`, `TestDaemonCredentialWriterSites` (unchanged), `TestApprovalImportsOnlyTheStandardLibrary`, `TestApprovalTestIsTestOnly`, `TestNoParallelWriteForkPackagesOutsideVerifygate`, `TestSubprocessSinkGate`, `TestEveryVerbHasHelp`, `TestServeProofFlags`, the quickstart tests, the schema/judge pins, sprint 2's route tests and the sprint's 7 new Go tests — **57 named `--- PASS`**, 0 SKIP/FAIL; the AW leg 9 `--- PASS` |
| C23 | **The live-Chrome tests, on the final tree** | §6's live leg (`WORLD_CHROME=… WORLD_REQUIRE_PASSKEY_BROWSER=1 go test -v ./cmd/ailang-worldd/ -run '^(TestPasskeyPageInChrome|TestPasskeyPageAbortsAPendingCeremony|TestApprovePage|TestApproveCard|TestNeedsYou)'`, Chrome 157.0.8103.51, `logs/final_gates/chrome.log`): `--- PASS: TestApproveCardNetworkDeadlines (79.17s)` with `--- PASS: …/get (21.74s)` and `--- PASS: …/post (57.43s)`; `--- PASS: TestApproveCardPasskeyDeadline (1.78s)`; `--- PASS: TestApprovePageInChrome (1.76s)`; `--- PASS: TestApprovePageRefusesASwappedRequest (1.73s)`; `--- PASS: TestApprovePageSSHSIGInChrome (1.20s)`; `--- PASS: TestNeedsYouListsPendingRequests (1.04s)`; A2's `--- PASS: TestPasskeyPageAbortsAPendingCeremony (3.33s)` (4 subtests PASS) and `--- PASS: TestPasskeyPageInChrome (3.17s)` (3 subtests PASS); `ok cmd/ailang-worldd 93.455s`; 0 SKIP, 0 FAIL |
| C24 | **Every mutant, on the final tree** | 88 rows, the last ledger line of each: **87 KILLED (named and suite), 1 named equivalent SURVIVED (named and suite)** (`JS-PAGEHIDE`), 0 not as expected, 0 `NOCOMPILE`; 109 ledger lines in all (fork A 18: `logs/mutant_ledger_forkA.tsv`; the planner's run1 67 and run2-final 24: `logs/mutant_ledger.tsv`). run1 found one gap, fixed in the test: `JS-WALK-TYPE` (the walk's type check removed) survived because XA2's `type` arm turned `ttl_s` into a string, which the integer check refuses too; the arm now turns `reuse` into a string, and run2-final kills it. No kill rests on an inventory: every suite excludes `TestDaemonCredentialWriterSites` and the other §7 inventories |
| C25 | `go vet ./...`, gofmt, `-race`, `verify_ail.sh` on the final tree | same §6 run: `go vet ./...` empty; `gofmt -l cmd host` empty; the compile fence clean; `-race` → `ok host/surfacepage 2.859s`, `ok host/daemon 4.703s`, `ok cmd/ailang-worldd 22.134s`, no race report; `./scripts/verify_ail.sh` → `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` (no `.ail` change) |
| C26 | **The whole repository on the final tree** (committed) | `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` on the committed round-trip tree at M8d, nothing else running → **rc 1**, two untouched packages red, both on the directive's known list: `host/broker` (`--- FAIL: TestWorktreeSnapshotRefusesAConfiguredFilter`: `worktree_delta_test.go:432: control: plain git status did not run the clean filter`) and `host/verifygate` (`--- FAIL: TestMissionConfigDefaultFlowIntegration/invalid_ledger`: `default-flow fixture timed out … context deadline exceeded`); every other leg `✓` (`✓ 0 other binary blobs among 6805 tracked files`, `✓ toolchain floor gate: go1.26.6 >= root module floor go1.26.6`, `✓ all 64 required top-level evidence tests passed exactly once`) and every other package `ok` (`logs/final_verify_go.log`). **Re-run alone** (§6's rule; `logs/final_rerun_alone.log`): `--- PASS: TestWorktreeSnapshotRefusesAConfiguredFilter (0.24s)`, `ok host/broker`; `--- PASS: TestMissionConfigDefaultFlowIntegration (7.97s)` with `/valid_ledger` and `/invalid_ledger` PASS, `ok host/verifygate`. Neither package is touched by the sprint (§5). **Re-run once, alone, the whole gate** → **rc 0**: `✓ 0 other binary blobs among 6805 tracked files`; `✓ toolchain floor gate: go1.26.6 >= root module floor go1.26.6`; `✓ all 64 required top-level evidence tests passed exactly once`; `✓ go gate PASSED (plain): build clean; legs run pass with pinned AILANG_BIN (AILANG v0.41.0` (`logs/final_verify_go_2.log`) |
| C27 | Size per stage (work) | `tools/numstat.sh` → M6a own 29 (0 prod + 29 test), cumulative **29**, M5d own 251 (27 prod + 224 test), cumulative **280**, M6b own 508 (411 prod + 97 test), cumulative **788**, M6c own 155 (0 prod + 155 test), cumulative **943**, M7a own 33 (25 prod + 8 test), cumulative **976**, M8a own 407 (0 prod + 407 test), cumulative **1383**, M8b own 119 (0 prod + 119 test), cumulative **1502**, M8c own 286 (0 prod + 286 test), cumulative **1788**, M8d own 55 (0 prod + 55 test), cumulative **1843**; docs 18 (M8d) |
| C28 | The SSHSIG arm judges before it records (fork A) | `git diff M6a M5d -- host/daemon/surface_approve.go`: `format, artifact, what := approval.FormatWebAuthn, []byte(body.WebAuthn), "the passkey's assertion"`; `if body.SSHSIG != nil { format, artifact, what = approval.FormatSSH, []byte(*body.SSHSIG), "the SSH signature" }`; `if format == approval.FormatSSH { // R11: the law, before anything is recorded` → `JudgeSSHSignature(ctx, ro, row, artifact, now)` else `JudgeAssertion(…)`; then the unchanged R12 `Binds`, R13 writer, R14 record with `Format: format, … Artifact: artifact` |
| C29 | **The SSH judge's store step is a 500, not a 503** (§1.4 P6) | fork A: the `approvers` table exists at every schema the daemon opens, so `JudgeSSHSignature`'s `store` step fails only on a damaged store → `approvalStoreError` → sanitized `500 Internal`; RT4 arm `R11-sshsig-store` (an `approvers` column renamed) pins it |
| C30 | **The patches reproduce the mutant tree** | a fresh `git archive 82fb9afa` export + `patch -p1` of Appendix P's nine patches **extracted from this file with §2's `awk`** (each `cmp`-equal to the generated patch; `--dry-run` first; no `.orig`/`.rej`) → `diff -r cmd host docs` against `plan-scratch/work` (the tree §4's mutants were cut from): **identical** (`ROUNDTRIP_IDENTICAL_TO_MUTANT_TREE`) (`tools/roundtrip.sh`, `logs/roundtrip.log`); nothing uncommitted after M8d |
| C31 | **Every sub-milestone is bisectable** | the same fresh export, after **each** patch: `tools/accept.sh <tree> <M>` = §3 verbatim (incl. the `[chrome]` legs, live) → `ACCEPT M6a rc=0` 43s, `ACCEPT M5d rc=0` 324s, `ACCEPT M6b rc=0` 295s, `ACCEPT M6c rc=0` 4s, `ACCEPT M7a rc=0` 10s, `ACCEPT M8a rc=0` 20s, `ACCEPT M8b rc=0` 29s, `ACCEPT M8c rc=0` 95s, `ACCEPT M8d rc=0` 145s (`logs/roundtrip.log`; a last extraction-only re-check after the final assembly, `SKIP_ACCEPT=1`, again printed `ROUNDTRIP_IDENTICAL_TO_MUTANT_TREE`); the cumulative numstat after a commit of each = §7: M6a **29**, M5d **280**, M6b **788**, M6c **943**, M7a **976**, M8a **1383**, M8b **1502**, M8c **1788**, M8d **1843** (docs 18) (`NUMSTAT` lines of `logs/roundtrip.log`) |

### 1.3 Planner decisions where the design leaves a choice

- **S3-D1 — The fallback cut fires, first step only (design §9).** Moved to a fourth sprint: the opt-in
  notification (§2.9: the "Notify me" control, the per-id `Notification`, XA4, the driver's
  `grantNotifications` and `notify` recorder), the Decided list (§2.8's "You · Minted · …" lines and "show all",
  XA3's Decided half) and the `N` jump (§1.4 P7). Kept: the card, Approve and the Needs-you link (XA1, XA2, XA3's
  list half). Measured (§1.5): 1,843 with everything else; the moved part ≈ 300 → ≈ 2,143 > 2,000. The
  second cut step (Appendix C to sprint 4) is **not** needed: Appendix C is 251 lines and the total stays
  ≤ 2,000. `model.json`'s `needsYou` stays `[]` (design §2.8).
- **S3-D2 — Order of the carried obligations.** M6a (J36) is the earliest sub-milestone: one test, no
  production change. M5d (the SSHSIG arm) comes next, before any page, so the card's SSH path posts to an arm
  that already judges, and the POST's file is finished before the page lands. Its id keeps sprint 2's M5 number:
  it completes M5.
- **S3-D3 — The SSHSIG arm shares the passkey arm's path** (C28): `format`/`artifact`/`what` are chosen from the
  body, the matching `authority.Judge*` runs at the poll's `now`, every judge error but `store` is
  `422 NotApproved` naming the format ("the SSH signature does not approve request …"), and R12–R14 are the
  passkey's, unchanged. The recorded row is `ssh` / `ssh:<name>` / the armored bytes as POSTed.
- **S3-D4 — The card page and its script are approval routes.** `GET /surface/approve/{id}` and
  `GET /surface/approve.js` register through `d.approvalRoute` (R2's Host check, 421 before anything), in a new
  `host/daemon/surface_approve_page.go`; the page answers a malformed id with R1's `400 BadRequest` JSON refusal
  and otherwise the same embedded bytes for every id, under `setSurfaceHeaders` (`surfaceCSP` unchanged). RT1's
  `approvalRoutes` gains both; RT7's targets gain both (a foreign Host → 421; loopback → 200);
  `TestSurfaceAppRoutes` drops `/surface/approve.js` from its 404 list, keeps `/surface/approve`, adds
  `/surface/approve/`, and expects 421 for both new routes under its non-surface Host.
- **S3-D5 — `FIELDS` is a JSON literal between `/* FIELDS */` markers**, one entry `{path, label, group, type}`
  per leaf key path of `approval.MintRequest` (`[]` for a slice element), `type` ∈ `string`/`bool`/`int`, and a
  boolean's two words (`no`, `yes`) in its entry ("creates a new worktree"/"uses an existing worktree",
  "git worktree"/"NON-GIT, DEGRADED", "clean build"/"modified build"). PG3 compares the `{path → type}` map with a
  reflection walk of the struct's `json` tags (structs recurse, slices add `[]`), and requires a label, a shown
  group and both words exactly for booleans. The 28 paths render in the design's five groups; each value is one
  `<dd id="f-<concrete path>">` (`grants.0.effect` …) whose text is the value (a boolean's word), read back by XA1.
  `ttl_s`'s label carries the hours ("Session lifetime (s) (1.0 h)").
- **S3-D6 — The fail-closed walk.** Every leaf of the parsed bytes must be a `FIELDS` path of the listed JSON
  type (`int` also `Number.isSafeInteger`), and every `FIELDS` path must occur; otherwise the page refuses,
  naming the first bad path, before any `credentials.get`, message link or button. A request that is not an
  object has a leaf at `''` or under `[]`, which no entry names, so it is refused by the same rule (no separate
  branch). The bytes are decoded with `TextDecoder('utf-8', {fatal: true})`; the card's `id` must be the page's.
- **S3-D7 — The network deadline is one helper.** `net(url, init)` runs `fetch` **and** `r.text()` under one
  `AbortController` with a 10 s timer (each attempt its own), and every outstanding controller is aborted on
  `pagehide`. A GET that fails or times out is state `failed` with "The request could not be read (no answer
  within 10 s). Nothing to approve yet." and a **Read it again** button; a non-200 card answer names its HTTP
  status and reason with the same button. A POST with no answer (timeout or network error) is state `unknown`:
  "Approval recording outcome unknown: … It may or may not have been recorded; nothing was minted by this page.
  Check its status before approving again." and a **Check status** button, which is one deadline-bound GET of
  `/surface/requests.json` per click (the only route that shows `recorded`, §1.4 P2) and says "Recorded: approved
  by …, waiting for the minter (approve <id> --await)." or "Not recorded: the request is still pending with no
  approval. You may approve it again." (re-enabling the buttons). Only `503` with `error: "Busy"` is retried,
  1 s apart, **at most three retries (four attempts)**, each with its own deadline; any other answer is shown
  as is (`rejected`, the class and reason) and never retried; nothing switches formats.
- **S3-D8 — After a 200 the page shows the POST's own answer** (`Approved by <approver>: <next>.`) and does not
  re-read the card (§1.4 P2); a card opened after the mint shows `used` ("Already used: approved by …") with no
  button (XA1's last step).
- **S3-D9 — Needs you is a pure function.** `needsYouHTML(requests, now)` (exported for node): `Needs you · N`;
  per pending request a link `/surface/approve/<id>` with id `needs-<id>`, "<episode> · <agent> · N grant(s) ·
  expires in M min", and ", approved by <approver>, waiting for the minter" when `recorded` is set; a `refused`
  row shows its reason and no link; an empty list shows `unavailable` or "Nothing needs you". All text through
  the existing `esc()`. `requests.json` is read before the model in `loadModel` and at the start of every
  `refresh` (the existing tick; no timer), and `nowCard` renders `needsYouHTML(REQS, REQS.now)`.
- **S3-D10 — The A2 driver is extended, never replaced.** New commands `fetchRewrite`, `openCard`, `press`
  (= the design's `click`, plus a virtual-time budget), `waitFor`, `credential`, `fill`, `save`, `card`; new
  injection `getRecorder`; `pageState`, `open`, `ceremony`, `uv`, `advance` unchanged. One shared helper changed:
  `click` aims at the element's first line box (`getClientRects()[0]`, falling back to the bounding box) — the
  same point for every one-line button A2 clicks (C17).
- **S3-D11 — How the browser holds a test passkey.** XA1 enrols through the A2 page and the CLI exactly as XB1
  does (the design's text). XA2, XA5 and the SSH test give the tab's virtual authenticator
  `approvaltest.NewPasskey(AlgES256, seed)`'s key through `WebAuthn.addCredential` (PKCS#8 of the same 32-byte
  seed; a control asserts its public key equals `pk.Key()`), and enrol the matching row in the temp store at
  `now-100`: Chrome signs, the daemon's law verifies, with no ceremony to drive (C21).
- **S3-D12 — Each browser view decodes into a fresh value** (C18).
- **S3-D13 — Two planner tests beyond the design's list** (each kills JS branches no XA test reaches):
  `TestApprovePageSSHSIGInChrome` (the page half of AW6: **Save the message** saves exactly
  `approval.SignedMessage` through a `blob:` of the card's origin; a pasted signature over another message →
  `422`, nothing recorded; the right one → `ssh:mark-mac` recorded, then minted by `approve --await`; with the only
  passkey at another origin the card names it and hides the passkey button) and
  `TestApproveCardPasskeyDeadline` (A2's 240 s bound on the card: an `abort`-injected get is aborted at 240 s of
  page time, a `late` result after it is ignored, an authenticator without user verification makes no response;
  each sends **no POST**).
- **S3-D14 — XA5's network is a test-only wrapper** around the production handler (`cardNet`): `get-header` /
  `post-header` / `list-header` (the status check's `requests.json`) never answer, `get-body` / `post-body` answer a header and `{"v":1,` then nothing,
  `post-recorded` runs the daemon's POST (it records) and withholds the answer, `get-500` / `restart` answer
  500 / `503 Restart`. The 503 Busy arms use the daemon's own refusal: the test holds the credentials writer
  (`store.OpenCredentials`), released after 1.5 s (3 POSTs, recorded) or never (4 POSTs, refused). Every stall
  ends at the browser's abort or the test's cleanup.
- **S3-D15 — The Chrome gate's names** (§1.4 P4): `-run '^(TestPasskeyPageInChrome|TestPasskeyPageAbortsAPendingCeremony|TestApprovePage|TestApproveCard|TestNeedsYou)'`.
- **S3-D16 — Sub-milestones** (§1.5): M6a J36 (29), M5d the SSHSIG arm (251), M6b the page and its
  routes with PG2 (508), M6c PG3/PG3b (155), M7a Needs you (33), M8a the driver, its Go
  helpers and XA1 (407), M8b XA2 and the SSH page test (119), M8c XA5 and the passkey bound
  (286), M8d XA3's list half and QUICKSTART §17 (55 + 18 docs). M6b and M8a are above
  ~350: M6b is one production unit (the page, its script, its two routes; `approve.js` alone is 316) with the
  route tests it changes; M8a is the driver extension (121) with the Go helpers (177) and the one test that
  proves both (109) — splitting either lands code no test runs.

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

- **P1 — Design §6 XA2 arm (b): "bytes and digest both forged → the page signs, the POST is 422".** R9 (the
  body's digest equals the row's) runs before R11 (sprint 2 C-order, `surface_approve.go:246`), so the forged
  card's POST is **`409 Conflict` "the request changed since the card was shown; reload it"**, 0 rows. T4's own
  text ("The body's `digest` (R9) names the mismatch") agrees. XA2 asserts the 409 and no recorded row.
- **P2 — Design §2.7 step 5: "…then re-reads `/surface/requests/{id}`".** The card JSON carries no recorded
  approval (RT3's shape: `consumed`/`approvedBy` move only at the mint), so a re-read cannot show the recording.
  Routed by S3-D7/S3-D8: the POST's answer is shown; the outcome-unknown status check reads `requests.json`'s
  `recorded`; a later open shows `used`.
- **P3 — Design §6 "the driver gains commands `fetchRewrite`, `grantNotifications`, `click`, and the
  `notify`/`getRecorder` injections".** `grantNotifications` and `notify` move with XA4 (S3-D1); `click` already
  exists as the driver's internal helper, so the command is `press`; the card also needs `openCard`, `waitFor`,
  `credential`, `fill`, `save`, `card` (S3-D10).
- **P4 — Design §9: the Chrome run `-run '^(TestApprovePage|TestApproveCard|TestNeedsYou|TestPasskeyPage)'` "must
  print `--- PASS` for every subtest".** `TestPasskeyPageRecordsTheFixtures` matches and is a `--- SKIP` at the
  base by design (it records only under `WORLD_PASSKEY_RECORD_DIR`, `passkey_fixture_record_test.go:26`). The
  gate names the other two A2 browser tests (S3-D15).
- **P5 — The J36 directive's injection.** `DROP TABLE webauthn_credentials` makes `OpenCredentialsReadOnly`
  refuse the store (its table-set check), so the card's 500 came from the open and the J36 mutant **survived**
  that test (C8); the test renames a column instead, so the error is `ActivePasskeys`' own.
- **P6 — Design §2.4 R11/R13 and Appendix C "the same law … 503 below the schema" for SSHSIG.** A store step
  error of `JudgeSSHSignature` cannot be "below schema" (the `approvers` table exists at every version the daemon
  opens), so it is the sanitized `500 Internal`, pinned by RT4 `R11-sshsig-store` (C29).
- **P7 — Design §2.8 "`N` jumps to the first one (the key already exists, `surface.js:257`)".** `surface.js` is
  the static timeline's script and its `N` is a help-table entry; `app.js` has no `N` key (C15). Moved with M7's
  cut (sprint 4), where it is new app code.
- **P8 — Design §6 RT1/RT7 "five routes" are now five** (sprint 2 P2 routed them to three; this sprint adds the
  card and its script) — no longer false; RT1/RT7 are extended (S3-D4).

### 1.5 The size

Lines are insertions + deletions under `cmd/ host/`: **`git diff --numstat <base> HEAD -- cmd host`**
(`docs/` and `design_docs/` are outside the filter; `.js`/`.html`/`.mjs` under `host/`/`cmd/` count). **Every row
is measured** from the scratch stages, each committed in a throwaway repo (`plan-scratch/work`, tags `M6a` …
`M8d`; `tools/numstat.sh`).

| Sub-milestone | Production | Test | Cumulative vs `82fb9afa` (after its commit) | Run |
|---|---|---|---|---|
| M6a — J36 `TestRequestCardPasskeyReadError` (carried) | 0 | 29 | **29** | 1 |
| M5d — the POST's SSHSIG arm; RT4's SSHSIG arms, RT6, RT13, AW6's route half (carried, sprint-2 Appendix C) | 27 | 224 | **280** | 1 |
| M6b — `host/surfacepage/approve/{index.html,approve.js}`, `approve.go`, the page handlers and two routes; PG2; RT1/RT7/`TestSurfaceAppRoutes` edits | 411 | 97 | **788** | 1 |
| M6c — PG3, PG3b | 0 | 155 | **943** | 2 |
| M7a — Needs you in `app.js` (`needsYouHTML`, `requests.json` on load and refresh); node assertions | 25 | 8 | **976** | 2 |
| M8a — the driver (`passkey_cdp.mjs`), the card's Go helpers; XA1 | 0 | 407 | **1,383** | 2 |
| M8b — XA2 (arms a–c); the card's SSH path | 0 | 119 | **1,502** | 3 |
| M8c — XA5; the card's passkey bound | 0 | 286 | **1,788** | 3 |
| M8d — XA3's list half; QUICKSTART §17 (+18 docs lines, not counted) | 0 | 55 | **1,843** | 3 |
| **Sprint 3** | **463** | **1,380** | **1,843 ≤ 2,000** (margin 157) | |

The design estimated sprint 3 at ≈ 2,200 (≈ 1,950 + ≈ 250 for XA5) **before** the carried Appendix C. As built,
without the cut part and with Appendix C (251) and two planner tests, it measured 1,843.

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of
four backticks, in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; a new file's patch has an
`@@ -0,0 …` hunk. Apply it **only** this way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF'
set -euo pipefail
M=M6a   # <- the sub-milestone, exactly as in the Appendix P heading: M6a M5d M6b M6c M7a M8a M8b M8c M8d
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter386-a3a-s3
PLAN="$W/design_docs/planned/w-world-approvals-by-key-a3a-sprint3-plan.md"
T="$HOME/.ailang/state/world-iter386/exec-patches"; mkdir -p "$T"
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
  command's output `tee`'d to a log under `$HOME/.ailang/state/world-iter386/exec-logs/` (never only through
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
each package run must end `ok` with no `FAIL` line. **[chrome]** means the command is prefixed with
`WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" WORLD_REQUIRE_PASSKEY_BROWSER=1`
and gets `-timeout 600s` (headless, temporary profiles, a virtual authenticator, test keys). The planner's
`tools/accept.sh <tree> <M> <logs>` is exactly these lines (§1.2 C31).

### M6a — J36: the card route's passkey read error · 0 prod + 29 test · Run 1

**Patch** `M6a`: new `host/daemon/surface_approve_card_error_test.go`: `TestRequestCardPasskeyReadError` — a
schema-7 store whose `webauthn_credentials.origin` column was renamed (the open and the request read still work;
`ActivePasskeys` fails, not as an older store): the card GET is `500 Internal` (`requireRefusal`) with no
`passkeys`, `request`, `challenge` or digest in the body; control: `200` with `"passkeys":[]` before the edit.
It kills iteration 385's J36 (`164: if err != nil && !isOlderStore(err) {` → `if false {`), §1.2 C7.

**Acceptance**
```
[base]
go test ./host/daemon/ -run '^(TestRequestCardPasskeyReadError|TestRequestCardJSON|TestApprovalRoutesRegistered)$' -count=1 -v
go test ./host/daemon/ -run '^(TestApprov|TestRequest|TestSurface|TestPasskeyPage)' -count=1
```

**Expected cumulative numstat after the controller's commit of M6a: 29**.

### M5d — the POST's SSHSIG arm (sprint-2 Appendix C) · 27 prod + 224 test · Run 1

**Patch** `M5d`:
- `host/daemon/surface_approve.go`: the sprint-2 refusal at `:255` is replaced by S3-D3's arm (`format`,
  `artifact`, `what`; `JudgeSSHSignature` or `JudgeAssertion` at `now`; the shared refusal; R12–R14 with
  `Format: format`, `Artifact: artifact`); the file and handler comments name both formats.
- `host/daemon/surface_approve_post_test.go` (RT4): the arm `R11-sshsig-cut` becomes `R11-sshsig-garbage` (422
  "the SSH signature does not approve request"), `R11-sshsig-unenrolled` (an `approvaltest.NewSigner(74)` key
  not enrolled → 422) and `R11-sshsig-store` (an `approvers` column renamed → 500 `Internal`), each with 0 rows
  and 0 sessions.
- new `host/daemon/surface_approve_sshsig_test.go`: rig helpers `apSigner` (`NewSigner(73)`), `enrolSSH` (at
  `apT0-10`, after the requests, so a judge at the request's creation refuses), `signatureFor`, `sshBody`,
  `recorded`; **RT6** `TestApprovePOSTRecordsAVerifiedSSHSIG` (200, one row `ssh`/`ssh:mark`, the armored bytes,
  0 sessions, `requests.json` shows `recorded: ssh:mark`), **RT13**
  `TestApprovePOSTSSHSIGViaTailnetOriginWithNoPasskey` (one SSH approver, no passkey; `SurfaceOrigins:
  ["https://studio.tail0.ts.net"]`; the POST with that Origin under Host `studio.tail0.ts.net` and Host
  `127.0.0.1:7644` → 200, one row `ssh:<name>`; the card GET under both → 200; unconfigured → 403 `CrossOrigin`
  naming `--surface-origin` (loopback Host) and 421 (ts.net Host), 0 rows; `Origin: https://other.tail0.ts.net` →
  403).
- new `cmd/ailang-worldd/approve_surface_sshsig_e2e_test.go`: **AW6's route half**
  `TestSurfaceSSHSIGThenAwaitMints` (E2E1's shape: the card GET, the 90-byte message from the card's digest,
  nonce and expiry, `sigSigner`'s SSHSIG POSTed → 200; `approve --await` mints `ssh:mark-mac` with no terminal;
  `requests.json` lists it decided by `ssh:mark-mac`).

**Acceptance** (~5–7 min: the whole `host/daemon`, because the POST changed)
```
[base]
go test ./host/daemon/ -run '^(TestApprovePOSTRefusalArms|TestApprovePOSTRecordsAVerifiedSSHSIG|TestApprovePOSTSSHSIGViaTailnetOriginWithNoPasskey|TestApprovePOSTRecordsAVerifiedPasskey|TestApprovePOSTCrossOrigin|TestDaemonCredentialWriterSites)$' -count=1 -v
go test ./cmd/ailang-worldd/ -run '^(TestSurfaceSSHSIGThenAwaitMints|TestSurfaceApproveThenAwaitMints|TestApproveAwaitSSHSIG)$' -count=1 -v
go test -count=1 -timeout 20m ./host/daemon/
```

**Expected cumulative numstat after the controller's commit of M5d: 280**.

### M6b — the card page and its routes · 411 prod + 97 test · Run 1

**Patch** `M6b`:
- new `host/surfacepage/approve/index.html` (the card: `#what`, `#card`, `#origin-note`, `#go`, `#status`,
  `#again`, the SSH box `#save`/`#sshline`/`#sig`/`#post-sig`; one `<script src="/surface/approve.js">`, no inline
  script or handler), new `host/surfacepage/approve/approve.js` (S3-D5–S3-D8: `FIELDS`, the walk, the digest and
  UTF-8 checks, the 90-byte message and its SHA-256, A2's get with its 240 s bound and late-result rule, the
  network helper with its 10 s deadlines, the Busy retries, outcome-unknown and its status check, the SSH path),
  new `host/surfacepage/approve.go` (`ApprovePage`, `ApproveScript`, `go:embed`).
- new `host/daemon/surface_approve_page.go`: `handleSurfaceApprovePage` (R1's 400, else the page under
  `setSurfaceHeaders`), `handleSurfaceApproveScript`; `host/daemon/daemon.go`: `mux.HandleFunc("GET
  /surface/approve/{id}", d.approvalRoute(d.handleSurfaceApprovePage))`, `("GET /surface/approve.js", …)` and
  the `Handler` doc comment.
- new `host/daemon/surface_approve_page_test.go`: **PG2** `TestApprovePageRoutes` (the embedded bytes at both
  routes, byte-equal to the files; CSP = the literal `surfaceCSP` text; `text/html`/`text/javascript`, `no-store`,
  `nosniff`, `no-referrer`; no `Set-Cookie`, no CORS; three malformed ids → 400; the A2 page's `index.html` and
  `passkey.js` sha256 unchanged; exactly one script, no inline handler/form/img/link/iframe/URL; labelled
  controls: `approve.js` has no `innerHTML`, `outerHTML`, `insertAdjacentHTML`, `document.write`, `eval(`,
  `new Function`, `document.cookie`, `XMLHttpRequest`, `sendBeacon`, `WebSocket`).
- `host/daemon/surface_approve_test.go`: RT1's `approvalRoutes` +2 and its 404 list (`/surface/approve`,
  `/surface/approve/`, `/surface/approve.jsx`, `/surface/approve/<id>/x`, …); RT7's targets +2.
  `host/daemon/surface_app_test.go`: S3-D4's 404/421 edit.

**Acceptance** (~5–7 min: the whole `host/daemon`, because `Handler()` changed)
```
[base]
go test ./host/daemon/ -run '^(TestApprovePageRoutes|TestApprovalRoutesRegistered|TestApprovalRoutesRefuseForeignHost|TestSurfaceAppRoutes|TestPasskeyPageRoutes|TestFrozenV1RouteTableMatchesMux|TestMCPMountMethodSource|TestDaemonDependencyAllowlist)$' -count=1 -v
go test ./host/surfacepage/ -count=1
go test -count=1 -timeout 20m ./host/daemon/
```

**Expected cumulative numstat after the controller's commit of M6b: 788**.

### M6c — every signed field has a renderer · 155 test · Run 2

**Patch** `M6c`: new `host/surfacepage/approve_page_test.go`: **PG3** `TestApprovePageRendersEveryCanonicalField`
(the `FIELDS` literal parsed as JSON equals the reflection walk of `approval.MintRequest` as `{path → type}`; every
entry a label, a shown group, and the two words exactly for a boolean; labelled control: 21 top-level keys,
design V24), **PG3b** `TestApprovePageFieldTagControls` (no `json` tag option anywhere in `MintRequest`, `Grant`,
`ProfileRef`, `DaemonRef`; a degraded request with `reuse: false` and every string empty still serializes 21
keys and every `FIELDS` path and no other leaf; `branch`, which XA2 (c) deletes, is a non-optional string).

**Acceptance**
```
[base]
go test ./host/surfacepage/ -run '^(TestApprovePageRendersEveryCanonicalField|TestApprovePageFieldTagControls|TestApproveModelSummaries)$' -count=1 -v
go test ./host/surfacepage/ -count=1
```

**Expected cumulative numstat after the controller's commit of M6c: 943**.

### M7a — Needs you lists the pending requests · 25 prod + 8 test · Run 2

**Patch** `M7a`: `host/surfacepage/app/app.js`: `needsYouHTML` (pure, exported, S3-D9); `REQS` and
`loadRequests()` (one `getJSON('/surface/requests.json')`), called in `loadModel` and at the start of `refresh`;
`nowCard` renders `needsYouHTML(REQS, REQS ? REQS.now : 0)` and no longer reads `ne.needsYou`.
`host/surfacepage/app_test.go` (`TestSurfaceAppPure`'s node harness): the export list +`needsYouHTML`, and its
assertions (empty → "Nothing needs you"; two rows: the link, the escaped episode, "2 grants · expires in 11
min", the recorded approver; a refused row with its escaped reason and no link; `unavailable`).

**Acceptance** (`TestSurfaceAppPure` needs node; CI sets `WORLD_EXEC_NODE`, C19)
```
[base]
WORLD_EXEC_NODE=$(command -v node) go test ./host/surfacepage/ -run '^(TestSurfaceAppPure|TestApprovePageRendersEveryCanonicalField)$' -count=1 -v
go test ./host/surfacepage/ -count=1
go test ./host/daemon/ -run '^(TestSurfaceAppRoutes|TestApprovePageRoutes)$' -count=1 -v
```

**Expected cumulative numstat after the controller's commit of M7a: 976**.

### M8a — the driver, the card's helpers, XA1 · 407 test · Run 2

**Patch** `M8a`: `cmd/ailang-worldd/testdata/passkey_cdp.mjs` (S3-D10: `fetchRewrite` with `rewritePaused`,
`cardState`/`card`, `openCard` (fresh document; `raw`), `press` (states, page, virtual time), `waitFor`,
`credential`, `fill`, `save`; the `getRecorder` injection; `click` on the first line box); new
`cmd/ailang-worldd/approve_page_driver_test.go` (`cardView`, `cardGet`, `cardOrigin` (the production handler,
optionally wrapped, on the loopback over the request's store, with an optional surface index),
`virtualPasskey` (S3-D11), `cardFields` (every leaf of the stored bytes by concrete path, a boolean as its
`FIELDS` word), `surfaceRecorded`); new `cmd/ailang-worldd/approve_page_chrome_test.go`: **XA1**
`TestApprovePageInChrome` (mark-mac enrolled through the A2 page and the CLI; the A2 page's get for the same
request via `--webauthn-link`; the card: ready, "Approve request <id> for episode ep1, expires …, message hash
<the CLI's hash16>", **every field's text equals the stored request's** and the TTL label; Approve → `recorded`
"Approved by webauthn:mark-mac: the minter (approve <id> --await) mints it…", one row `webauthn:mark-mac`; zero
CSP violations and zero page errors; the card's get = the A2 page's get (challenge, rpId `localhost`, the one
credential, `required`, 240000, a signal); `approve --await` mints it with no terminal; reopened → `used`).

**Acceptance** (~1 min)
```
[base]
[chrome] go test ./cmd/ailang-worldd/ -run '^(TestApprovePageInChrome|TestPasskeyPageInChrome|TestPasskeyPageAbortsAPendingCeremony)$' -count=1 -v
```

**Expected cumulative numstat after the controller's commit of M8a: 1,383**.

### M8b — the card refuses what it cannot bind; its SSH path · 119 test · Run 3

**Patch** `M8b`: new `cmd/ailang-worldd/approve_page_swap_test.go`: **XA2** `TestApprovePageRefusesASwappedRequest`
(`fetchRewrite` on the card's GET: control ready; (a) `bytes` and `id` → "The daemon's request does not match its
digest; nothing to approve.", no get; (b) `forge` → ready, Approve → `rejected` "(HTTP 409 Conflict): the request
changed since the card was shown" with one get, nothing recorded (§1.4 P1); (c) `extra`, `drop:branch`, `type`,
`float` → "This request has a field the page cannot show (<path>); nothing to approve; approve it at a
terminal.", `utf8` → "The daemon's request is not UTF-8 JSON (…", `nonce`, `expiry` → "The daemon's challenge or
expiry is malformed; nothing to approve.", `b64` → "The daemon's request is not base64url; nothing to approve."
— each refused card: no get, no button, no message link, no field; 11 responses rewritten) and **`TestApprovePageSSHSIGInChrome`** (S3-D13).

**Acceptance** (~1 min)
```
[base]
[chrome] go test ./cmd/ailang-worldd/ -run '^(TestApprovePageRefusesASwappedRequest|TestApprovePageSSHSIGInChrome|TestApprovePageInChrome)$' -count=1 -v
```

**Expected cumulative numstat after the controller's commit of M8b: 1,502**.

### M8c — the card's network deadlines and its passkey bound · 286 test · Run 3

**Patch** `M8c`: new `cmd/ailang-worldd/approve_card_deadline_test.go`: `cardNet` (S3-D14) and **XA5**
`TestApproveCardNetworkDeadlines` — subtest `get`: a 500 card → `failed` with the daemon's reason and **Read it
again**; `get-header` and `get-body` each end after the 10 s deadline (9.5–25 s measured from the open) as
`failed` "The request could not be read (no answer within 10 s). Nothing to approve yet." with no field and no
button, and **Read it again** reads the card; subtest `post`: `post-header`, `post-body`, `post-recorded` each end
after the deadline as `unknown` "Approval recording outcome unknown: no answer within 10 s." (never `rejected` or
`recorded`), nothing / nothing / one row recorded, and **Check status** answers "Not recorded: …" / "Not
recorded: …" / "Recorded: approved by webauthn:mark-mac, waiting for the minter (approve <id> --await)."; a
stalled status check ends at its own deadline with "The status could not be read (no answer within 10 s)." and
offers **Check again**; recorded with no answer and then minted by `approve --await` before the check → "The
request is no longer pending (used or expired): read the card again."; then, on a second request (`ep2`), a
writer held 1.5 s → 3 POSTs then `recorded`; held throughout → 4 POSTs then `rejected` "(HTTP 503 Busy): the
credentials store is busy"; a `503 Restart` → 1 POST, `rejected`; no POST ever carries `sshsig`. And
**`TestApproveCardPasskeyDeadline`** (S3-D13).

**Acceptance** (~2 min: XA5 waits out seven 10 s deadlines)
```
[base]
[chrome] go test ./cmd/ailang-worldd/ -run '^(TestApproveCardNetworkDeadlines|TestApproveCardPasskeyDeadline)$' -count=1 -v
```

**Expected cumulative numstat after the controller's commit of M8c: 1,788**.

### M8d — Needs you in a live Chrome; QUICKSTART §17 · 55 test + 18 docs · Run 3

**Patch** `M8d`: new `cmd/ailang-worldd/needs_you_chrome_test.go`: **XA3's list half**
`TestNeedsYouListsPendingRequests` (a daemon with `host/surface/testdata/index-rig.json` (copied to a temp file)
and the request's store; `/surface`'s now card shows exactly "Needs you · 1ep1 · claude · 10 grants · expires
in 10 min"; a second request appears as "Needs you · 2" after the app's own `window.World.poll()` refresh; the
first link is `/surface/approve/<id>` and opens the card, ready, with no page error or CSP violation).
`docs/QUICKSTART.md`: **§17 Approve from the surface *(attended — pending first verbatim run)*** (design §2.10:
the §0 scratch store and the §16 passkey; `episode start … --request`; `approve <id> --await &`; Needs you → the
card → Approve; the minter mints; `--list --decided`; a second Approve refused as already used; the SSH path;
the 10 s limits and the Busy retries; `--surface-origin` behind `tailscale serve`).

**Acceptance** (~4–5 min: the whole `cmd/ailang-worldd`, the XA tests there as counted skips)
```
[base]
[chrome] go test ./cmd/ailang-worldd/ -run '^(TestNeedsYouListsPendingRequests|TestApprovePageInChrome)$' -count=1 -v
go test -count=1 -timeout 20m ./cmd/ailang-worldd/
```

**Expected cumulative numstat after the controller's commit of M8d: 1,843**.

---

## 4. Mutation table

**88 rows run on the planner's final tree: 87 KILLED (named and suite), 1 named equivalent(s) SURVIVED as predicted, 0 not as expected.** 109 ledger lines in all (fork A run1 18 (M6a, M5d); planner run1 67 (M6b, M7a: 67 rows), run2-final 24 (every row whose named test or driver mode changed after run1 started: XA2's and XA5's rows, the run1 survivor JS-WALK-TYPE, and the three rows added for the new arms)).

**Derivation.** Anchored to Appendix P's production hunks, one row per new branch: `surface_approve_card_error_test.go` kills J36 (M6a); `surface_approve.go`'s SSHSIG arm (the format dispatch, the artifact source, the judge dispatch, skipping the SSH judge, judging at the request's creation, the store-step arm, a record inserted before the law, the refusal text, the recorded format/artifact/approver for both formats, and sprint 2's passkey-arm mutants re-anchored to the restructured lines) (M5d, fork A's `mutA/`); `daemon.go` (each of the two registrations) and `surface_approve_page.go` (the id check, both handlers' headers, the script's content type, the page's bytes) (M6b, Go); `approve.js` (the digest check, the card's id, the fatal UTF-8 decode, the walk's unlisted / missing / type / integer arms, the challenge and expiry checks, each of the message's four parts and its hash, the get's rpId / allowCredentials / user verification / signal, the 240 s timer, the late-result rule, the failure arm, the origin group and its note, the consumed arm, the network deadline and its body coverage, the pagehide abort, the card GET's status arm and its retry offer, the Busy predicate / count / pause, the unknown outcome, both status-check arms, the 200 arm, re-enabling after a refusal, the message link's bytes / name / visibility, the namespace, the SSH POST's field, the passkey POST's digest, the field ids, the boolean words, the TTL label, the hash16 shown, the passkey button's visibility) (M6b, JS); `app.js` (the load and refresh fetches, the now card's call and its clock, and `needsYouHTML`'s refused / recorded / escape / rounding / link / unavailable / count / plural arms) (M7a). **Procedure**: `plan-scratch/tools/mutate.py <run> --workers 4` (the table is `tools/mutants.py`; every anchor asserted to occur exactly once in the final tree): for each row, in its own worker copy, apply **one** edit, build with `-vet=off` (a mutant that does not compile is `NOCOMPILE`, never counted), run **the single named test** and then **its package suite** with the inventories excluded — `host/daemon` `^(TestApprov|TestRequest|TestSurfaceAppRoutes|TestPasskeyPageRoutes)`, `host/surfacepage` whole (node via `WORLD_EXEC_NODE`), `cmd/ailang-worldd` `^(TestApprovePage|TestApproveCard|TestNeedsYou)` in live headless Chrome (fork A's M5d suites: `host/daemon` `^(TestApprov|TestRequest|TestSurfaceAppRoutes)`, `cmd/ailang-worldd` `^(TestApprove|TestSurfaceApprove|TestSurfaceSSHSIG)`) — restore, assert the restore byte-identical, and append one ledger line (UTC time, run, id, stage, file, from → to, named test, named, suite, the suite's failing tests). A row is KILLED only when both fail. No kill in this table rests on an inventory (`TestDaemonCredentialWriterSites` and the other §7 inventories are outside every suite). **Never commit or snapshot a mutant.** In the table a newline prints as ⏎, a backquote as ' and a pipe as \|.

**Named equivalents** (survive by construction):
- `JS-PAGEHIDE` — the browser itself aborts a document's in-flight fetches when the document is unloaded (and Chrome keeps no page with an in-flight fetch in the back/forward cache), so removing the explicit abort on pagehide changes nothing a test can observe; it stays as the design's explicit statement ("Navigation cancels outstanding operations", §2.7 [r2])

| Id | Stage | Mutant (single edit) | Named test | Named | Suite | Suite failures (inventories excluded) |
|---|---|---|---|---|---|---|
| MUT-J36-CARD-PASSKEY-ERR | M6a | surface_approve.go: active, err := ro.ActivePasskeys(ctx)⏎if err != nil && !isOlderStore(err) { → active, err := ro.ActivePasskeys(ctx)⏎if false { | `TestRequestCardPasskeyReadError` | KILLED | KILLED | TestRequestCardPasskeyReadError |
| MUT-SSH-DISPATCH | M5d | surface_approve.go: if body.SSHSIG != nil {⏎format, artifact, what → if false {⏎format, artifact, what | `TestApprovePOSTRecordsAVerifiedSSHSIG` | KILLED | KILLED | TestApprovePOSTRefusalArms,TestApprovePOSTRecordsAVerifiedSSHSIG,TestApprovePOSTSSHSIGViaTailnetOriginWithNoPasskey |
| MUT-SSH-ARTIFACT-SRC | M5d | surface_approve.go: approval.FormatSSH, []byte(*body.SSHSIG), "the SSH signature" → approval.FormatSSH, []byte(body.WebAuthn), "the SSH signature" | `TestApprovePOSTRecordsAVerifiedSSHSIG` | KILLED | KILLED | TestApprovePOSTRefusalArms,TestApprovePOSTRecordsAVerifiedSSHSIG,TestApprovePOSTSSHSIGViaTailnetOriginWithNoPasskey |
| MUT-SSH-JUDGE-DISPATCH | M5d | surface_approve.go: if format == approval.FormatSSH { // R11 → if false { // R11 | `TestApprovePOSTRecordsAVerifiedSSHSIG` | KILLED | KILLED | TestApprovePOSTRefusalArms,TestApprovePOSTRecordsAVerifiedSSHSIG,TestApprovePOSTSSHSIGViaTailnetOriginWithNoPasskey |
| MUT-SSH-SKIP-LAW | M5d | surface_approve.go: signed, err = authority.JudgeSSHSignature( → signed, _ = authority.JudgeSSHSignature( | `TestApprovePOSTRefusalArms` | KILLED | KILLED | TestApprovePOSTRefusalArms |
| MUT-SSH-NOW | M5d | surface_approve.go: authority.JudgeSSHSignature(ctx, ro, row, artifact, now) → authority.JudgeSSHSignature(ctx, ro, row, artifact, row.CreatedAt) | `TestApprovePOSTRecordsAVerifiedSSHSIG` | KILLED | KILLED | TestApprovePOSTRecordsAVerifiedSSHSIG,TestApprovePOSTSSHSIGViaTailnetOriginWithNoPasskey |
| MUT-SSH-RECORD-BEFORE-LAW | M5d | surface_approve.go: signed, err = authority.JudgeSSHSignature( → if w2, e2 := store.OpenCredentials(d.credentialsPath); e2 == nil {⏎_ = w2.RecordSurfaceApproval(ctx, store.SurfaceApprovalRow{RequestID: id, Format: approval.FormatSSH, Approver: "ssh:unjudged", Artifact: artifact, Digest: row.Digest}, now)⏎_ = w2.Close()⏎}⏎signed, err = authority.JudgeSSHSignature( | `TestApprovePOSTRefusalArms` | KILLED | KILLED | TestApprovePOSTRefusalArms |
| MUT-R11-STORE-STEP | M5d | surface_approve.go: je.Step == authority.JudgeStepStore { → je.Step == "never" { | `TestApprovePOSTRefusalArms` | KILLED | KILLED | TestApprovePOSTRefusalArms |
| MUT-SSH-REFUSAL-TEXT | M5d | surface_approve.go: []byte(*body.SSHSIG), "the SSH signature" → []byte(*body.SSHSIG), "the passkey's assertion" | `TestApprovePOSTRefusalArms` | KILLED | KILLED | TestApprovePOSTRefusalArms |
| MUT-REC-FORMAT-SSH | M5d | surface_approve.go: Format: format, Approver → Format: approval.FormatWebAuthn, Approver | `TestApprovePOSTRecordsAVerifiedSSHSIG` | KILLED | KILLED | TestApprovePOSTRecordsAVerifiedSSHSIG,TestApprovePOSTSSHSIGViaTailnetOriginWithNoPasskey |
| MUT-REC-FORMAT | M5d | surface_approve.go: Format: format, Approver → Format: approval.FormatSSH, Approver | `TestApprovePOSTRecordsAVerifiedPasskey` | KILLED | KILLED | TestApprovePOSTCrossOrigin,TestApprovePOSTReplay,TestApprovePOSTWriterBusy,TestApprovePOSTRefusalArms,TestApprovePOSTRecordsAVerifiedPasskey |
| MUT-REC-ARTIFACT | M5d | surface_approve.go: Artifact: artifact, Digest: row.Digest}, d.approveNow()) → Artifact: []byte("{}"), Digest: row.Digest}, d.approveNow()) | `TestApprovePOSTRecordsAVerifiedSSHSIG` | KILLED | KILLED | TestApprovePOSTRecordsAVerifiedPasskey,TestApprovePOSTRecordsAVerifiedSSHSIG |
| MUT-REC-APPROVER | M5d | surface_approve.go: Approver: signed.ApprovedBy(),⏎ → Approver: "ssh:someone",⏎ | `TestApprovePOSTRecordsAVerifiedSSHSIG` | KILLED | KILLED | TestApprovePOSTCrossOrigin,TestApprovePOSTReplay,TestApprovePOSTWriterBusy,TestApprovePOSTRefusalArms,TestApprovePOSTRecordsAVerifiedPasskey,TestApprovePOSTRecordsAVerifiedSSHSIG,TestApprovePOSTSSHSIGViaTailnetOriginWithNoPasskey |
| MUT-R11-NOW | M5d | surface_approve.go: authority.JudgeAssertion(ctx, ro, row, artifact, now) → authority.JudgeAssertion(ctx, ro, row, artifact, row.CreatedAt) | `TestApprovePOSTRecordsAVerifiedPasskey` | KILLED | KILLED | TestApprovePOSTCrossOrigin,TestApprovePOSTReplay,TestApprovePOSTWriterBusy,TestApprovePOSTRefusalArms,TestApprovePOSTRecordsAVerifiedPasskey |
| MUT-R11-SKIP-LAW | M5d | surface_approve.go: signed, err = authority.JudgeAssertion( → signed, _ = authority.JudgeAssertion( | `TestApprovePOSTRefusalArms` | KILLED | KILLED | TestApprovePOSTRefusalArms |
| MUT-RECORD-BEFORE-LAW | M5d | surface_approve.go: signed, err = authority.JudgeAssertion( → if w2, e2 := store.OpenCredentials(d.credentialsPath); e2 == nil {⏎_ = w2.RecordSurfaceApproval(ctx, store.SurfaceApprovalRow{RequestID: id, Format: approval.FormatWebAuthn, Approver: "webauthn:unjudged", Artifact: artifact, Digest: row.Digest}, now)⏎_ = w2.Close()⏎}⏎signed, err = authority.JudgeAssertion( | `TestApprovePOSTRefusalArms` | KILLED | KILLED | TestApprovePOSTNoSessionNoCookie,TestApprovePOSTRefusalArms,TestApprovePOSTRecordsAVerifiedPasskey |
| MUT-PK-ARTIFACT-SRC | M5d | surface_approve.go: approval.FormatWebAuthn, []byte(body.WebAuthn), "the passkey's assertion" → approval.FormatWebAuthn, []byte("{}"), "the passkey's assertion" | `TestApprovePOSTRecordsAVerifiedPasskey` | KILLED | KILLED | TestApprovePOSTCrossOrigin,TestApprovePOSTReplay,TestApprovePOSTWriterBusy,TestApprovePOSTRefusalArms,TestApprovePOSTRecordsAVerifiedPasskey |
| MUT-SSH-E2E-FORMAT | M5d | surface_approve.go: Format: format, Approver → Format: approval.FormatWebAuthn, Approver | `TestSurfaceSSHSIGThenAwaitMints` | KILLED | KILLED | TestSurfaceSSHSIGThenAwaitMints |
| MUT-ROUTE-PAGE | M6b | daemon.go:  mux.HandleFunc("GET /surface/approve/{id}", d.approvalRoute(d.handleSurfaceApprovePage))⏎ →  | `TestApprovalRoutesRegistered` | KILLED | KILLED | TestSurfaceAppRoutes,TestApprovePageRoutes,TestApprovalRoutesRegistered,TestApprovalRoutesRefuseForeignHost |
| MUT-ROUTE-SCRIPT | M6b | daemon.go:  mux.HandleFunc("GET /surface/approve.js", d.approvalRoute(d.handleSurfaceApproveScript))⏎ →  | `TestApprovalRoutesRegistered` | KILLED | KILLED | TestSurfaceAppRoutes,TestApprovePageRoutes,TestApprovalRoutesRegistered,TestApprovalRoutesRefuseForeignHost |
| MUT-PAGE-HEADERS | M6b | surface_approve_page.go:  setSurfaceHeaders(w)⏎ _, _ = w.Write(surfacepage.ApprovePage) →  _, _ = w.Write(surfacepage.ApprovePage) | `TestApprovePageRoutes` | KILLED | KILLED | TestApprovePageRoutes |
| MUT-PAGE-ID | M6b | surface_approve_page.go: if !store.MintRequestIDPattern.MatchString(r.PathValue("id")) { → if false && !store.MintRequestIDPattern.MatchString(r.PathValue("id")) { | `TestApprovePageRoutes` | KILLED | KILLED | TestApprovePageRoutes |
| MUT-SCRIPT-HEADERS | M6b | surface_approve_page.go:  setSurfaceHeaders(w)⏎ w.Header().Set("Content-Type", "text/javascript; charset=utf-8") →  w.Header().Set("Content-Type", "text/javascript; charset=utf-8") | `TestApprovePageRoutes` | KILLED | KILLED | TestApprovePageRoutes |
| MUT-SCRIPT-CT | M6b | surface_approve_page.go:  w.Header().Set("Content-Type", "text/javascript; charset=utf-8")⏎ _, _ = w.Write(surfacepage.ApproveScript) →  _, _ = w.Write(surfacepage.ApproveScript) | `TestApprovePageRoutes` | KILLED | KILLED | TestApprovePageRoutes |
| MUT-PAGE-BYTES | M6b | surface_approve_page.go: _, _ = w.Write(surfacepage.ApprovePage) → _, _ = w.Write(surfacepage.PasskeyPage) | `TestApprovePageRoutes` | KILLED | KILLED | TestApprovePageRoutes |
| JS-DIGEST | M6b | approve.js: hex(await sha256(bytes)) !== card.digest → false | `TestApprovePageRefusesASwappedRequest` | KILLED | KILLED | TestApprovePageRefusesASwappedRequest |
| JS-CARD-ID | M6b | approve.js: if (card.id !== id \|\| hex( → if (hex( | `TestApprovePageRefusesASwappedRequest` | KILLED | KILLED | TestApprovePageRefusesASwappedRequest |
| JS-UTF8-FATAL | M6b | approve.js: new TextDecoder('utf-8', { fatal: true }) → new TextDecoder('utf-8', {}) | `TestApprovePageRefusesASwappedRequest` | KILLED | KILLED | TestApprovePageRefusesASwappedRequest |
| JS-WALK-UNLISTED | M6b | approve.js: const f = byPath.get(l.pattern); → const f = byPath.get(l.pattern) \|\| { type: typeof l.value === 'boolean' ? 'bool' : typeof l.value === 'number' ? 'int' : 'string' }; | `TestApprovePageRefusesASwappedRequest` | KILLED | KILLED | TestApprovePageRefusesASwappedRequest |
| JS-WALK-MISSING | M6b | approve.js:     for (const f of FIELDS) if (!seen.has(f.path)) return f.path;⏎ →  | `TestApprovePageRefusesASwappedRequest` | KILLED | KILLED | TestApprovePageRefusesASwappedRequest |
| JS-WALK-TYPE | M6b | approve.js: typeof l.value !== TYPES[f.type] \|\|  →  | `TestApprovePageRefusesASwappedRequest` | KILLED | KILLED | TestApprovePageRefusesASwappedRequest |
| JS-WALK-INT | M6b | approve.js: (f.type === 'int' && !Number.isSafeInteger(l.value)) → false | `TestApprovePageRefusesASwappedRequest` | KILLED | KILLED | TestApprovePageRefusesASwappedRequest |
| JS-CHALLENGE-FMT | M6b | approve.js: if (!/^[0-9a-f]{64}$/.test(card.challenge) \|\|  → if ( | `TestApprovePageRefusesASwappedRequest` | KILLED | KILLED | TestApprovePageRefusesASwappedRequest |
| JS-EXPIRY-FMT | M6b | approve.js:  \|\| !Number.isSafeInteger(card.expires) \|\| card.expires < 0) { → ) { | `TestApprovePageRefusesASwappedRequest` | KILLED | KILLED | TestApprovePageRefusesASwappedRequest |
| JS-MSG-LABEL | M6b | approve.js: message.set(new TextEncoder().encode('world-approval/v1'), 0); → message.set(new TextEncoder().encode('world-approval/v2'), 0); | `TestApprovePageInChrome` | KILLED | KILLED | TestApproveCardNetworkDeadlines,TestApprovePageInChrome,TestApprovePageSSHSIGInChrome |
| JS-MSG-DIGEST | M6b | approve.js: message.set(unhex(card.digest), 18); → message.set(unhex(card.digest), 17); | `TestApprovePageInChrome` | KILLED | KILLED | TestApproveCardNetworkDeadlines,TestApprovePageInChrome,TestApprovePageSSHSIGInChrome |
| JS-MSG-NONCE | M6b | approve.js: message.set(unhex(card.challenge), 50); → message.set(unhex(card.challenge), 49); | `TestApprovePageInChrome` | KILLED | KILLED | TestApproveCardNetworkDeadlines,TestApprovePageInChrome,TestApprovePageSSHSIGInChrome |
| JS-MSG-EXPIRY | M6b | approve.js: setBigUint64(82, BigInt(card.expires)); → setBigUint64(82, BigInt(card.expires + 1)); | `TestApprovePageInChrome` | KILLED | KILLED | TestApproveCardNetworkDeadlines,TestApprovePageInChrome,TestApprovePageSSHSIGInChrome |
| JS-RPID | M6b | approve.js: rpId: location.hostname, → rpId: 'example.com', | `TestApprovePageInChrome` | KILLED | KILLED | TestApproveCardNetworkDeadlines,TestApproveCardPasskeyDeadline,TestApprovePageInChrome,TestApprovePageRefusesASwappedRequest |
| JS-ALLOW | M6b | approve.js: allowCredentials: group.allow.map((a) => ({ type: 'public-key', id: unb64u(a) })), → allowCredentials: [], | `TestApprovePageInChrome` | KILLED | KILLED | TestApproveCardNetworkDeadlines,TestApprovePageInChrome,TestApprovePageRefusesASwappedRequest |
| JS-CHALLENGE-HASH | M6b | approve.js: challenge = await sha256(message); → challenge = message.slice(0, 32); | `TestApprovePageInChrome` | KILLED | KILLED | TestApproveCardNetworkDeadlines,TestApprovePageInChrome |
| JS-UV | M6b | approve.js: userVerification: 'required', timeout: TIMEOUT_MS }, signal: controller.signal }); → userVerification: 'preferred', timeout: TIMEOUT_MS }, signal: controller.signal }); | `TestApprovePageInChrome` | KILLED | KILLED | TestApprovePageInChrome |
| JS-SIGNAL | M6b | approve.js: timeout: TIMEOUT_MS }, signal: controller.signal }); → timeout: TIMEOUT_MS } }); | `TestApproveCardPasskeyDeadline` | KILLED | KILLED | TestApproveCardPasskeyDeadline,TestApprovePageInChrome |
| JS-PASSKEY-TIMER | M6b | approve.js:       buttons(true);⏎    }, TIMEOUT_MS); →       buttons(true);⏎    }, 10 * TIMEOUT_MS); | `TestApproveCardPasskeyDeadline` | KILLED | KILLED | TestApproveCardPasskeyDeadline |
| JS-LATE | M6b | approve.js: if (expired) return; // a result after the timer is ignored → if (false) return; // a result after the timer is ignored | `TestApproveCardPasskeyDeadline` | KILLED | KILLED | TestApproveCardPasskeyDeadline |
| JS-GROUP | M6b | approve.js: .find((p) => p.origin === location.origin) \|\| null; → .find(() => true) \|\| null; | `TestApprovePageSSHSIGInChrome` | KILLED | KILLED | TestApprovePageSSHSIGInChrome |
| JS-NOGROUP-NOTE | M6b | approve.js: others.length ? 'No passkey is enrolled for → false ? 'No passkey is enrolled for | `TestApprovePageSSHSIGInChrome` | KILLED | KILLED | TestApprovePageSSHSIGInChrome |
| JS-FAILURE | M6b | approve.js:     if (failure !== null) { →     if (false) { | `TestApproveCardPasskeyDeadline` | KILLED | KILLED | TestApproveCardPasskeyDeadline |
| JS-CONSUMED | M6b | approve.js:     if (card.consumed) { →     if (false) { | `TestApprovePageInChrome` | KILLED | KILLED | TestApprovePageInChrome |
| JS-NET-DEADLINE | M6b | approve.js: const timer = setTimeout(() => c.abort(), NET_MS); → const timer = 0; | `TestApproveCardNetworkDeadlines` | KILLED | KILLED | TestApproveCardNetworkDeadlines |
| JS-GET-STATUS | M6b | approve.js:     if (res.status !== 200 \|\| !res.body) { →     if (false) { | `TestApproveCardNetworkDeadlines` | KILLED | KILLED | TestApproveCardNetworkDeadlines |
| JS-PAGEHIDE | M6b | approve.js: addEventListener('pagehide', () => { for (const c of outstanding) c.abort(); }); → addEventListener('pagehide', () => {}); | `TestApproveCardNetworkDeadlines` | SURVIVED | SURVIVED |  |
| JS-NET-BODY | M6b | approve.js:       const text = await r.text(); →       clearTimeout(timer);⏎      const text = await r.text(); | `TestApproveCardNetworkDeadlines` | KILLED | KILLED | TestApproveCardNetworkDeadlines |
| JS-GET-RETRY | M6b | approve.js:       say('The request could not be read (${why(e)}). Nothing to approve yet.', true);⏎      return offer('Read it again', load); →       say('The request could not be read (${why(e)}). Nothing to approve yet.', true);⏎      return; | `TestApproveCardNetworkDeadlines` | KILLED | KILLED | TestApproveCardNetworkDeadlines |
| JS-RETRY-PREDICATE | M6b | approve.js: if (res.status === 503 && b.error === 'Busy' && attempt < RETRIES) { → if (res.status !== 200 && attempt < RETRIES) { | `TestApproveCardNetworkDeadlines` | KILLED | KILLED | TestApproveCardNetworkDeadlines |
| JS-RETRY-COUNT | M6b | approve.js: b.error === 'Busy' && attempt < RETRIES) { → b.error === 'Busy' && attempt <= RETRIES) { | `TestApproveCardNetworkDeadlines` | KILLED | KILLED | TestApproveCardNetworkDeadlines |
| JS-RETRY-SLEEP | M6b | approve.js:         await sleep(RETRY_MS);⏎ →  | `TestApproveCardNetworkDeadlines` | KILLED | KILLED | TestApproveCardNetworkDeadlines |
| JS-UNKNOWN | M6b | approve.js:         setState('unknown'); →         setState('rejected'); | `TestApproveCardNetworkDeadlines` | KILLED | KILLED | TestApproveCardNetworkDeadlines |
| JS-POST-OK | M6b | approve.js:       if (res.status === 200) { →       if (true) { | `TestApprovePageRefusesASwappedRequest` | KILLED | KILLED | TestApproveCardNetworkDeadlines,TestApprovePageRefusesASwappedRequest,TestApprovePageSSHSIGInChrome |
| JS-STATUS-RECORDED | M6b | approve.js:     if (p && p.recorded) { →     if (false) { | `TestApproveCardNetworkDeadlines` | KILLED | KILLED | TestApproveCardNetworkDeadlines |
| JS-REJECT-BUTTONS | M6b | approve.js:       say('Not recorded (HTTP ${res.status} ${b.error \|\| ''}): ${b.reason \|\| 'no reason given'}', true);⏎      buttons(true); →       say('Not recorded (HTTP ${res.status} ${b.error \|\| ''}): ${b.reason \|\| 'no reason given'}', true); | `TestApprovePageSSHSIGInChrome` | KILLED | KILLED | TestApprovePageSSHSIGInChrome |
| JS-STATUS-PENDING | M6b | approve.js:     } else if (p) { →     } else if (false) { | `TestApproveCardNetworkDeadlines` | KILLED | KILLED | TestApproveCardNetworkDeadlines |
| JS-SAVE-BLOB | M6b | approve.js: new Blob([message], { type: 'application/octet-stream' }) → new Blob([unhex(card.digest)], { type: 'application/octet-stream' }) | `TestApprovePageSSHSIGInChrome` | KILLED | KILLED | TestApprovePageSSHSIGInChrome |
| JS-SAVE-NAME | M6b | approve.js: save.download = 'world-approval-${id}.msg'; → save.download = 'world-approval-${id}.bin'; | `TestApprovePageSSHSIGInChrome` | KILLED | KILLED | TestApprovePageInChrome,TestApprovePageSSHSIGInChrome |
| JS-NAMESPACE | M6b | approve.js: const NAMESPACE = 'world-approval@sunholo'; → const NAMESPACE = 'world-approval'; | `TestApprovePageSSHSIGInChrome` | KILLED | KILLED | TestApprovePageSSHSIGInChrome |
| JS-POST-SIG | M6b | approve.js: post({ v: 1, digest, sshsig: $('sig').value }); → post({ v: 1, digest, webauthn: $('sig').value }); | `TestApprovePageSSHSIGInChrome` | KILLED | KILLED | TestApprovePageSSHSIGInChrome |
| JS-POST-DIGEST | M6b | approve.js: await post({ v: 1, digest, webauthn: assertionJSON(cred) }); → await post({ v: 1, digest: card.challenge, webauthn: assertionJSON(cred) }); | `TestApprovePageInChrome` | KILLED | KILLED | TestApproveCardNetworkDeadlines,TestApprovePageInChrome |
| JS-RENDER-ID | M6b | approve.js: dd.id = 'f-' + l.concrete; → dd.id = 'f-' + l.pattern; | `TestApprovePageInChrome` | KILLED | KILLED | TestApprovePageInChrome |
| JS-SHOWN-BOOL | M6b | approve.js: (v ? f.yes : f.no) → (v ? f.no : f.yes) | `TestApprovePageInChrome` | KILLED | KILLED | TestApprovePageInChrome |
| JS-TTL-LABEL | M6b | approve.js: (f.path === 'ttl_s' ?  → (false ?  | `TestApprovePageInChrome` | KILLED | KILLED | TestApprovePageInChrome |
| JS-WHAT-HASH16 | M6b | approve.js: message hash ${hex(challenge.slice(0, 8))} → message hash ${hex(challenge.slice(0, 7))} | `TestApprovePageInChrome` | KILLED | KILLED | TestApprovePageInChrome |
| JS-BUTTON-SHOWN | M6b | approve.js:       $('go').hidden = false;⏎ →  | `TestApprovePageInChrome` | KILLED | KILLED | TestApproveCardNetworkDeadlines,TestApproveCardPasskeyDeadline,TestApprovePageInChrome,TestApprovePageRefusesASwappedRequest |
| JS-SAVE-SHOWN | M6b | approve.js:     save.hidden = false;⏎ →  | `TestApprovePageSSHSIGInChrome` | KILLED | KILLED | TestApprovePageInChrome,TestApprovePageSSHSIGInChrome |
| JS-B64 | M6b | approve.js:       return refuse('The daemon\'s request is not base64url; nothing to approve.'); →       bytes = new Uint8Array(0); | `TestApprovePageRefusesASwappedRequest` | KILLED | KILLED | TestApprovePageRefusesASwappedRequest |
| JS-STATUS-CATCH | M6b | approve.js:       say('The status could not be read (${why(e)}).', true);⏎      return offer('Check again', status); →       return; | `TestApproveCardNetworkDeadlines` | KILLED | KILLED | TestApproveCardNetworkDeadlines |
| JS-STATUS-GONE | M6b | approve.js: say(res.status === 200 ? 'The request is no longer pending → say(false ? 'The request is no longer pending | `TestApproveCardNetworkDeadlines` | KILLED | KILLED | TestApproveCardNetworkDeadlines |
| APP-LOAD | M7a | app.js:     await loadRequests();⏎    try { res = await getJSON('/surface/model.json'); →     try { res = await getJSON('/surface/model.json'); | `TestNeedsYouListsPendingRequests` | KILLED | KILLED | TestNeedsYouListsPendingRequests |
| APP-PURE-REFUSED | M7a | app.js: p.map(q => q.refused ?  → p.map(q => false ?  | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| APP-PURE-RECORDED | M7a | app.js: (q.recorded ? ' · approved by → (false ? ' · approved by | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| APP-PURE-ESC | M7a | app.js: <b>${esc(q.episode)}</b> → <b>${q.episode}</b> | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| APP-PURE-CEIL | M7a | app.js: Math.max(0, Math.ceil((q.expires - now) / 60)) → Math.max(0, Math.floor((q.expires - now) / 60)) | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| APP-PURE-LINK | M7a | app.js: href="/surface/approve/${esc(q.id)}" → href="/surface/approve/" | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| APP-PURE-UNAVAILABLE | M7a | app.js: esc(r && r.unavailable \|\| 'Nothing needs you') → esc('Nothing needs you') | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| APP-PURE-COUNT | M7a | app.js: Needs you · ${p.length}</div>'; → Needs you · ${p.length ? 1 : 0}</div>'; | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| APP-PURE-PLURAL | M7a | app.js: grant${q.grants === 1 ? '' : 's'} → grant | `TestSurfaceAppPure` | KILLED | KILLED | TestSurfaceAppPure |
| APP-REFRESH | M7a | app.js:       await loadRequests();⏎      const r = await fetchCoherent( →       const r = await fetchCoherent( | `TestNeedsYouListsPendingRequests` | KILLED | KILLED | TestNeedsYouListsPendingRequests |
| APP-NOWCARD | M7a | app.js: ${needsYouHTML(REQS, REQS ? REQS.now : 0)} → ${needsYouHTML(null, 0)} | `TestNeedsYouListsPendingRequests` | KILLED | KILLED | TestNeedsYouListsPendingRequests |
| APP-NOW-ARG | M7a | app.js: ${needsYouHTML(REQS, REQS ? REQS.now : 0)} → ${needsYouHTML(REQS, 0)} | `TestNeedsYouListsPendingRequests` | KILLED | KILLED | TestNeedsYouListsPendingRequests |

---

## 5. Files (final `git status --short` after M8d)

| File | Sub-milestone | Change |
|---|---|---|
| `host/daemon/surface_approve_card_error_test.go` | M6a | new: J36 |
| `host/daemon/surface_approve.go`, `surface_approve_post_test.go` | M5d | the SSHSIG arm; RT4's SSHSIG arms |
| `host/daemon/surface_approve_sshsig_test.go`, `cmd/ailang-worldd/approve_surface_sshsig_e2e_test.go` | M5d | new: RT6, RT13; AW6's route half |
| `host/surfacepage/approve.go`, `approve/index.html`, `approve/approve.js` | M6b | new: the card |
| `host/daemon/surface_approve_page.go`, `surface_approve_page_test.go` | M6b | new: the page handlers; PG2 |
| `host/daemon/daemon.go`, `surface_approve_test.go`, `surface_app_test.go` | M6b | two routes + comment; RT1/RT7; the 404/421 list |
| `host/surfacepage/approve_page_test.go` | M6c | new: PG3, PG3b |
| `host/surfacepage/app/app.js`, `app_test.go` | M7a | Needs you; node assertions |
| `cmd/ailang-worldd/testdata/passkey_cdp.mjs`, `approve_page_driver_test.go`, `approve_page_chrome_test.go` | M8a | the driver; new: helpers, XA1 |
| `cmd/ailang-worldd/approve_page_swap_test.go` | M8b | new: XA2, the SSH page test |
| `cmd/ailang-worldd/approve_card_deadline_test.go` | M8c | new: XA5, the passkey bound |
| `cmd/ailang-worldd/needs_you_chrome_test.go`, `docs/QUICKSTART.md` | M8d | new: XA3's list half; §17 |
| this plan and `sprint_w-world-approvals-by-key-a3a-s3.json` | planner | |

---

## 6. Sprint gates (controller, after the last sub-milestone; each measured at the base, §1.1)

In **one** `bash` shell, from the worktree, **not under `nohup`**, nothing else running. Every leg `tee`s its
full output to a log under `$HOME/.ailang/state/world-iter386/gates/`; a verdict is read from its log, never
from `tail` alone. The live-Chrome leg runs on this machine only (CI has no Chrome step, C19).

```
bash <<'EOF'
set -euo pipefail
export AILANG_BIN=$HOME/.pinned-ailang/ailang
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter386-a3a-s3; cd "$W"
L=$HOME/.ailang/state/world-iter386/gates; mkdir -p "$L"; MODE=final
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
if [ "$MODE" = base ]; then NEW_SP=""; NEW_D=""; NEW_C=""; fi
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
WANT_CHROME="TestPasskeyPageInChrome TestPasskeyPageAbortsAPendingCeremony"
if [ "$MODE" != base ]; then WANT_CHROME="$WANT_CHROME TestApprovePageInChrome TestApprovePageRefusesASwappedRequest TestApprovePageSSHSIGInChrome TestApproveCardNetworkDeadlines TestNeedsYouListsPendingRequests"; fi
for t in $WANT_CHROME; do grep -qE -- "^--- PASS: $t " "$L/chrome.log" || { echo "$t did not PASS in Chrome"; exit 1; }; done
echo "PRE_COMMIT_GATES_GREEN $(( $(date +%s) - s0 ))s" | tee "$L/pre_commit.ok"
EOF
# AFTER the controller commits (the tracked-binary gate reads committed blobs; this leg runs go test ./...):
bash -c 'set -o pipefail; export AILANG_BIN=$HOME/.pinned-ailang/ailang; unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER; cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter386-a3a-s3 && VERIFY_GO_LEGS=plain ./scripts/verify_go.sh 2>&1 | tee $HOME/.ailang/state/world-iter386/gates/verify_go.log; echo "verify_go rc=$?" | tee -a $HOME/.ailang/state/world-iter386/gates/verify_go.log'
git -C /Users/voightkampff/dev/sunholo-data/.wt-world-iter386-a3a-s3 status --short   # only the two plan files, until the controller commits them
```

The pre-commit block is the planner's `tools/gates.sh <tree> <logs> final`; `base` mode omits the sprint's new
tests and Chrome names and expects the same 9 `TestApproveAwait` passes.

**At the base**: every pre-commit leg is green at `82fb9afa` in 420 s (§1.1). **On the final tree**: `PRE_COMMIT_GATES_GREEN 473s` on the round-trip tree: pkgs all `ok` (`host/daemon` 278.4 s, `cmd/ailang-worldd` 118.3 s), 57 named tests `--- PASS`, the AW leg 9, `-race` ok ×3, the live-Chrome leg every test and subtest `--- PASS` (§1.2 C22, C23, C25); the after-commit `verify_go.sh` rc 1 on its first run (two known-list tests, each green alone) and **rc 0** on its second, alone (§1.2 C26).

**What CI pins and what it does not (design R3).** CI has no Chrome step: XA1, XA2, XA3, XA5 and the two
planner browser tests are counted `SKIP`s there (`SKIP COUNTED: the passkey page browser test needs
WORLD_CHROME …`), exactly as A2's XB1/XB2. CI **does** run PG2 (the page's bytes, routes, headers, the unchanged
`surfaceCSP`, the A2 page's hashes, the labelled controls), PG3 and PG3b (every signed field has a renderer; no
`omitempty`), RT1/RT7 (both new routes behind the Host check), J36, RT4/RT6/RT13 and AW6's route half (the
SSHSIG arm), and `TestSurfaceAppPure` (the Needs-you renderer under node, `WORLD_EXEC_NODE`). The live-Chrome leg
above is the controller's.

**Known load- or environment-sensitive (re-run once, alone; do not chase)**:
`TestSeedRestartIgnoresASurvivingOldCopy` (host/broker), `TestHostsLaunchersExecHost`, `host/broker`'s
exec-probe timeouts (incl. `TestExecProtectedCheckTimeout/reserve`) and its local `-race`, the local `-race`
of `host/daemon` and `host/verifygate` (this sprint's `-race` leg is the approval subset only),
`TestDemoSeedPreflightRefusals`, `TestWorktreeSnapshotRefusesAConfiguredFilter`,
`TestMissionConfigDefaultFlowIntegration/invalid_ledger`, `scripts/marketingcheck`'s supervisor tests and
`TestMarketingCaptureFakeChromeSixShots/inventory`, `host/daemon`'s `TestAGUISlowReaderCutIsResumable`,
`TestSessionSelfGit`, `TestServeStartsWhenTheProbeSpendsTheStartupClock`, `TestRegistryPhaseHasItsOwnClock`,
`cmd/ailang-worldd`'s `TestWorldContextAndStatusline` and `TestCLIRealSubprocessEpisode`, and anything run
under `nohup`. The XA5 deadline windows (9.5–25 s) assume nothing else heavy runs; a red there under load is
re-run once alone. Any other red, any race report, any red in `host/surfacepage`, `host/store`,
`host/authority`, `host/approval`, in `host/daemon`'s `TestApprov*`/`TestRequest*`/`TestSurfaceApp*`/
`TestDaemonCredential*` or in `cmd/ailang-worldd`'s `TestApprove*`/`TestServe*`/`TestSurfaceApprove*`/
`TestSurfaceSSHSIG*`/`TestNeedsYou*`/`TestPasskeyPage*`, is a regression. CI's `go-race-verify` is the race
verdict.

**Must stay green, with only the named edits**: every a0a1, A1b, A2 and A3a-sprint-1/2 test; the only edits to
existing tests are RT1/RT7 (`surface_approve_test.go`), `TestSurfaceAppRoutes`, RT4's SSHSIG arm
(`surface_approve_post_test.go`) and `TestSurfaceAppPure`'s harness (`app_test.go`). In particular
`TestFrozenV1RouteTableMatchesMux` (13, untouched), `TestPasskeyPageRoutes` and XB1/XB2 (the A2 page unchanged;
the driver's `click` change keeps them green), `TestDaemonCredentialWriterSites` (unchanged: no writer site
added) and every §7 inventory.

---

## 7. Executor directive notes

- **No git writes.** Never `git add`, `commit`, `stash`, `checkout`, `reset`, `restore`, `rebase`, `push`,
  `worktree` or `apply`. `patch -p1` (§2) is the only writer; read-only `git status`/`git diff` are fine.
- **Never `nohup`, never background a test leg.** Every acceptance command runs in the foreground of one bash
  shell, `tee`'d to `$HOME/.ailang/state/world-iter386/exec-logs/`. (Iteration 385's executor broke this once.)
- **Chrome only through the tests.** The Chrome legs (`[chrome]`) are `go test` with `WORLD_CHROME` and
  `WORLD_REQUIRE_PASSKEY_BROWSER=1` on that command line only: headless, a temporary profile under the test's
  temp dir, a CDP virtual authenticator, test keys. Never open Chrome by hand, never a real profile; no
  `ssh-keygen`, no real store.
- **Runs.** Run 1 = M6a, M5d, M6b; Run 2 = M6c, M7a, M8a; Run 3 = M8b, M8c, M8d. Each run reads §0, §2 and its
  own §3 sections, then applies its patches in order.
- **Expected numstat after each sub-milestone, measured AFTER the controller's commit of it** (an uncommitted
  tree hides the new untracked files from numstat):
  ```
  bash -c 'cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter386-a3a-s3 && git diff --numstat 82fb9afa9086eed2f781bb0d4fe54d4bae3b70d7 HEAD -- cmd host | awk "{s+=\$1+\$2} END {print s}"'
  ```
  After the controller's commit of M6a **29**, M5d **280**, M6b **788**, M6c **943**, M7a
  **976**, M8a **1,383**, M8b **1,502**, M8c **1,788**, M8d **1,843** (and `-- docs` **18** after
  M8d). **Any other value means a patch was not applied as written: stop and report.**
- **Snapshot after every sub-milestone, cumulatively**:
  ```
  bash -c 'set -euo pipefail; M=M6a; W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter386-a3a-s3
    D=$HOME/.ailang/state/world-iter386/snap/$M; rm -rf "$D"; mkdir -p "$D"; cd "$W"
    git status --porcelain --untracked-files=all | while read -r st f; do
      case "$f" in design_docs/planned/w-world-approvals-by-key-a3a-sprint3-plan.md|design_docs/planned/sprint_w-world-approvals-by-key-a3a-s3.json) continue;; esac
      mkdir -p "$D/$(dirname "$f")"; cp -p "$f" "$D/$f"; done
    (cd "$D" && find . -type f ! -name MANIFEST.sha256 | sort | xargs shasum -a 256) > "$D/MANIFEST.sha256"
    cat "$D/MANIFEST.sha256"'
  ```
  (`M` = `M6a` … `M8d`.) No file is deleted, so the porcelain status has only `M` and `??`; stop on any other
  status. The controller rebuilds one commit per sub-milestone (`$HOME/.ailang/state/world-iter386/commit_ms.sh`
  with `WT` and the snapshot root pointed here; **commit the plan BEFORE running it — it deletes untracked
  files**).
- **Stop conditions.** A `patch --dry-run` failure; a numstat other than the expected value; a red
  earlier-milestone test; a `--- FAIL`, or a `--- SKIP` in a §3 `-v` run; any need to touch a file outside §5;
  running Chrome other than through the `[chrome]` test commands; any reference to `~/.ssh`, `~/.ailang/world`
  or a real `/dev/tty`. Stop and report; do not work around.
- **Timing.** M5d's and M6b's acceptance ~5–7 min each (the whole `host/daemon`, 278 s at the base); M8d's whole
  `cmd/ailang-worldd` ~3–5 min; M8c ~2 min (seven 10 s deadlines); every other acceptance under 1–2 min plus
  the compile fence (~25–60 s).
- **Shell.** Wrap loops, globs and `$VAR` expansions in `bash -c` or a `bash <<'EOF'` heredoc (zsh does not
  word-split). Export `AILANG_BIN` in the same shell as the Go tests. `tee` every gate's output to a file;
  never pipe a gate's only copy through `tail`. Poll nothing in the background: there is nothing to background.

---

## 8. Risks

- **R1 — Live-browser tests are local only** (design R3): CI skips XA1–XA5 (counted); the controller's
  live-Chrome leg (§6) is the verdict for them. PG2/PG3/PG3b and the node harness pin the page and the list in CI.
- **R2 — XA5 measures real time** (seven 10 s deadlines, 9.5–25 s windows): heavy machine load could stretch an
  arm; re-run once alone (§6).
- **R3 — The notification and the Decided list are not built** (S3-D1): an operator sees pending requests and
  approves them; who approved what stays visible through `approve --list --decided` and `requests.json`'s
  `decided` until sprint 4. A3a's exit run waits for sprint 4 only for its notification step (design §9).
- **R4 — `app.js` is under active style work** (design §8): M7a touches only `nowCard`'s Needs-you line,
  `loadModel` and `refresh` (+ the pure renderer); a conflicting edit to those lines is a controller rebase.
- **R5 — Safari, Firefox, a `*.ts.net` HTTPS origin and a real authenticator are UNVERIFIED** (design R5, R9):
  the attended exit run (design §10) checks them.
- **R6 — Behind `tailscale serve`, which `Host` arrives is UNVERIFIED** (design §2.3); RT13 pins both Hosts for
  the configured origin.

---

## 9. Acceptance criteria (the judge scores against these)

1. **J36 is killed** (M6a): a passkey read error on the card route is `500 Internal` with no card; the
   mutation `if err != nil && !isOlderStore(err)` → `if false` fails `TestRequestCardPasskeyReadError` and the
   non-inventory `host/daemon` subset (§1.2 C7).
2. **The SSHSIG arm judges before it records** (M5d, RT4, RT6, RT13, AW6's route half): an `sshsig` body is
   judged by `authority.JudgeSSHSignature` at the daemon's `now`; every refusal leaves 0 rows and 0 sessions;
   a verified one records `ssh`/`ssh:<name>`/the armored bytes and mints nothing; through a configured tailnet
   origin with no passkey, under both Hosts, it records; unconfigured or another origin, it is refused;
   `approve --await` mints it.
3. **The card page is an approval route** (PG2, RT1, RT7): `GET /surface/approve/{id}` and
   `GET /surface/approve.js` serve the embedded bytes under the unchanged `surfaceCSP` and the surface headers,
   behind the Host check (421 for a foreign Host), 400 for a malformed id; `/v1` stays 13; the A2 page's bytes
   are unchanged; no cookie, no CORS.
4. **Every signed field has a renderer** (PG3, PG3b): `FIELDS` equals `approval.MintRequest`'s key paths and
   types; no tag has `omitempty`; a degraded request still carries every path.
5. **What the card shows is what it signs** (XA1): every field's text equals the stored request's; its get's
   challenge, rpId, credential and user verification equal the A2 page's for the same request; the recorded
   assertion is minted by `approve --await`; no CSP violation or page error; a used request shows `used`.
6. **The card refuses what it cannot bind** (XA2): changed bytes or another id → refused with no get; forged
   bytes and digest → the POST is refused (409, R9) and nothing recorded; an extra, missing, mistyped or
   fractional field, non-UTF-8 or non-base64url bytes, a non-hex nonce or a negative expiry → refused naming the
   problem, no get, no message link.
7. **The card's network deadlines** (XA5): every GET and POST, body included, ends at 10 s; a GET shows a
   retryable error; a POST shows "outcome unknown" (never refused or minted) and a status check that tells
   recorded from not recorded and from minted meanwhile (itself deadline-bound); only `503 Busy` is retried, 1 s apart, at most three times; nothing switches
   format.
8. **The card's passkey bound and SSH path** (planner tests): a get is aborted at 240 s, a late result and a
   no-UV failure send nothing; the saved message is `approval.SignedMessage`; a pasted SSHSIG is recorded and
   minted, a wrong one refused.
9. **Needs you** (M7a, XA3's list half): the now card lists each pending request from `requests.json` with a
   link to its card (refused rows with their reason), refreshed on the app's existing tick, and the link opens
   the card.
10. **The fallback cut is explicit and minimal**: the notification, the Decided list and the `N` key are
    Appendix C's sprint-4 scope with the numbers that forced them; Appendix C of sprint 2 is built.
11. Each §4 mutant, applied alone, makes its named test **and** its package suite fail (inventories excluded);
    the named equivalents survive for the stated reasons; the ledger has a line per run.
12. Every §6 gate passes, including every inventory by name, the `-race` legs, the live-Chrome leg and
    `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` after the commit (modulo §6's re-run rule); `git status
    --short` lists only §5.
13. The numstat after each sub-milestone's commit equals §7's value; the sprint total is ≤ 2,000.
14. `host/approval`, `host/store`, `host/authority`, `surfaceCSP`, `isProtected`, `/surface/passkey`,
    `passkey.js`, `model.go` unchanged; no new production exec site, module, package, context root or `/v1`
    route; no CORS header; no `t.Parallel`; `approve.js` has no `innerHTML`, `eval(`, `document.cookie`.
15. No real key, store, `/dev/tty`, `~/.ssh` access, real browser profile, mint into a real store, publish or
    `ssh-keygen` anywhere in the sprint (D-WORLD-109); Chrome only headless with temporary profiles and a
    virtual authenticator.

---

## Appendix C — after this sprint (sprint-4 scope moved by the fallback cut; no re-design)

- **Moved here by design §9's fallback cut (first step), unchanged from design §2.8/§2.9/§6**: the opt-in
  browser notification (the "Notify me" control calling `Notification.requestPermission()` only from its click;
  on a refresh with a pending id not seen before in the tab and permission `granted`, one `new Notification("World
  needs you", {body: "Approve mint request for <episode> (<n> grants), expires <hh:mm>", tag: <id>})` whose click
  focuses the tab and opens `/surface/approve/<id>`); **XA4** `TestNeedsYouNotifiesOnlyWhenAskedOnce` with the
  driver's `grantNotifications` command and `notify` recorder injection; the **Decided** list under Needs you
  ("You · Minted · <episode> · <approvedBy> · <hh:mm>", the last few, a "show all" toggling the 20 of
  `requests.json`'s `decided`) and **XA3's Decided half** ("after the mint, Decided shows `webauthn:<name>`");
  the `N` key jumping to the first Needs-you link (new in `app.js`, §1.4 P7).
- **Why (measured)**: this sprint as built is 1,843 lines (`git diff --numstat 82fb9afa HEAD -- cmd host`);
  the moved part is ≈ 300 at this tree's density (the design's ≈ 160 `app.js` for M7 minus the 25 built here
  ≈ 135; XA4 ≈ 115, the size of XA2; the driver ≈ 30; XA3's Decided half ≈ 20) → ≈ 2,143 > 2,000.
- **Carry**: `needsYouHTML(requests, now)` is the place the Decided list and the notification's trigger read
  (`requests.json` is already fetched on load and on every refresh); the card's address is
  `/surface/approve/<id>`; `model.json`'s `needsYou` stays `[]` (A4 may use it).
- **Not moved**: sprint 2's Appendix C (built here as M5d); every never-cut item (PG2, PG3, PG3b, XA1, XA2 incl.
  arm (c), XA5, J36).
- **D-WORLD-114** stays open (default A applied); if Mark vetoes (b), the minter's mint location returns to
  design.

---

## Appendix P — the patches (apply with §2; generated from the compiled scratch stages; do not edit)

### Patch M6a — J36: the card route's passkey read error

````diff M6a
--- a/host/daemon/surface_approve_card_error_test.go
+++ b/host/daemon/surface_approve_card_error_test.go
@@ -0,0 +1,29 @@
+package daemon
+
+// Row 201 A3 sprint 3 M6a (design §2.3; iteration 385's judge survivor J36): the card route's
+// passkey read failing is the sanitized 500, never a card that lists no passkeys.
+
+import (
+	"strings"
+	"testing"
+)
+
+// TestRequestCardPasskeyReadError: a schema-7 store whose webauthn_credentials table lost a column
+// (the table set still opens, the request row still reads, and ActivePasskeys fails, not as an
+// older store): the card GET is 500 Internal with no card in the body. Control: 200 before.
+func TestRequestCardPasskeyReadError(t *testing.T) {
+	a := newApprovalRig(t)
+	row := a.putRequest("mr-00000000000000e1", "w201e", apT0-30)
+	target := "/surface/requests/" + row.RequestID
+	if rec := a.serve("GET", target, apLoopback, nil); rec.Code != 200 || !strings.Contains(rec.Body.String(), `"passkeys":[]`) {
+		t.Fatalf("control: the card = %d %s", rec.Code, rec.Body.String())
+	}
+	setVersion(t, a.d.credentialsPath, `ALTER TABLE webauthn_credentials RENAME COLUMN origin TO origin_gone`) // still version 7
+	rec := a.serve("GET", target, apLoopback, nil)
+	requireRefusal(t, "a passkey read error", rec, 500, "Internal")
+	for _, key := range []string{`"passkeys"`, `"request"`, `"challenge"`, row.Digest} {
+		if strings.Contains(rec.Body.String(), key) {
+			t.Fatalf("the 500 carries %s: %s", key, rec.Body.String())
+		}
+	}
+}
````

### Patch M5d — the POST's SSHSIG arm (sprint-2 Appendix C)

````diff M5d
--- a/cmd/ailang-worldd/approve_surface_sshsig_e2e_test.go
+++ b/cmd/ailang-worldd/approve_surface_sshsig_e2e_test.go
@@ -0,0 +1,84 @@
+package main
+
+// Row 201 A3 sprint 3 M5d (design §2.4, §6 AW6's route half; the sprint-2 plan's Appendix C): an
+// SSHSIG POSTed to the surface is recorded, and approve --await mints it, as E2E1 does a passkey's.
+
+import (
+	"bytes"
+	"context"
+	"encoding/hex"
+	"encoding/json"
+	"io"
+	"net"
+	"net/http"
+	"net/http/httptest"
+	"path/filepath"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/daemon"
+)
+
+// TestSurfaceSSHSIGThenAwaitMints: episode start --request; the card GET; the 90-byte message
+// from the card's digest, nonce and expiry; sigSigner's SSHSIG POSTed (200, recorded); approve
+// --await mints it with no terminal; requests.json lists it under decided, by ssh:mark-mac.
+func TestSurfaceSSHSIGThenAwaitMints(t *testing.T) {
+	now := time.Now().Unix() // the daemon's routes read the wall clock
+	a := newApproveRig(t, "ep1", now)
+	seedApprover(t, a.credPath, "mark-mac", sigSigner.PublicKey(), now-100)
+	d, err := daemon.New(boundedTestContext(t), daemon.Config{DBPath: filepath.Join(t.TempDir(), "world.db"), CredentialsPath: a.credPath,
+		BindHost: daemon.DefaultBindHost, ErrorLog: io.Discard})
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = d.Close() }()
+	serve := func(method, target string, body []byte, out any) int {
+		req := httptest.NewRequest(method, target, bytes.NewReader(body))
+		req = req.WithContext(context.WithValue(req.Context(), http.LocalAddrContextKey, net.Addr(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 7644})))
+		req.Host = "localhost:7644"
+		for k, v := range map[string]string{"Origin": "http://localhost:7644", "Sec-Fetch-Site": "same-origin", "Content-Type": "application/json"} {
+			req.Header.Set(k, v)
+		}
+		rec := httptest.NewRecorder()
+		d.Handler().ServeHTTP(rec, req)
+		if err := json.Unmarshal(rec.Body.Bytes(), out); err != nil {
+			t.Fatalf("%s %s: %d %s", method, target, rec.Code, rec.Body.String())
+		}
+		return rec.Code
+	}
+	var card struct {
+		Digest, Challenge string
+		Expires           int64
+	}
+	if code := serve("GET", "/surface/requests/"+a.id, nil, &card); code != 200 || card.Digest != a.row(t).Digest {
+		t.Fatalf("card = %d %+v", code, card)
+	}
+	dg, _ := hex.DecodeString(card.Digest)
+	nonce, _ := hex.DecodeString(card.Challenge)
+	msg, err := approval.SignedMessage([32]byte(dg), [32]byte(nonce), card.Expires)
+	if err != nil {
+		t.Fatal(err)
+	}
+	body, _ := json.Marshal(map[string]any{"v": 1, "digest": card.Digest, "sshsig": string(sigSigner.Sign(approval.Namespace, msg))})
+	var recorded struct {
+		Recorded   bool
+		ApprovedBy string
+	}
+	if code := serve("POST", "/surface/approve/"+a.id, body, &recorded); code != 200 || !recorded.Recorded || recorded.ApprovedBy != "ssh:mark-mac" {
+		t.Fatalf("POST = %d %+v", code, recorded)
+	}
+	var opened bool
+	env, _ := awaitEnv(time.Now().Unix(), &opened, nil)
+	var stdout, stderr bytes.Buffer
+	if code := runApprove(a.args("--await", "--await-for", "5s"), &stdout, &stderr, env); code != exitOK || opened {
+		t.Fatalf("approve --await: rc %d, terminal opened %v\nstderr %s", code, opened, stderr.String())
+	}
+	requireMintedBy(t, a, "ssh:mark-mac", stdout.String())
+	var list struct {
+		Decided []struct{ ID, ApprovedBy string }
+	}
+	if code := serve("GET", "/surface/requests.json", nil, &list); code != 200 || len(list.Decided) != 1 || list.Decided[0].ID != a.id || list.Decided[0].ApprovedBy != "ssh:mark-mac" {
+		t.Fatalf("requests.json = %d %+v", code, list)
+	}
+}
--- a/host/daemon/surface_approve.go
+++ b/host/daemon/surface_approve.go
@@ -3,8 +3,8 @@ package daemon
 // Row 201 A3 sprint 2 (design_docs/planned/w-world-approvals-by-key-a3.md §2.3, §2.4): the
 // approval routes under /surface, none under /v1. GET /surface/requests.json (Needs you and the
 // decision list) and GET /surface/requests/{id} (the card's data) read the credentials store on a
-// fresh read-only handle; POST /surface/approve/{id} records a passkey approval only after the
-// A2 law verified it. Each route answers only a Host in the surface origin set O (R2); every
+// fresh read-only handle; POST /surface/approve/{id} records a passkey or SSH approval only after
+// the A1b/A2 law verified it. Each route answers only a Host in the surface origin set O (R2); every
 // existing route stays exactly as row 199 leaves it.
 
 import (
@@ -189,7 +189,8 @@ func (d *Daemon) approvalStoreError(w http.ResponseWriter, r *http.Request, err
 }
 
 // handleSurfaceApprove is POST /surface/approve/{id} (design §2.4): after R2 (approvalRoute), the
-// checks R1, R3-R12 in order on a read-only handle, each refusal writing nothing; then R13 a
+// checks R1, R3-R12 in order on a read-only handle (R11 judges the passkey assertion or the SSH
+// signature with the CLI's own judge), each refusal writing nothing; then R13 a
 // non-waiting writer, opened for R14's RecordSurfaceApproval and closed at once. It reads no
 // Cookie and no Authorization, answers no CORS, and never mints: the minter (approve --await)
 // judges the recorded artifact again.
@@ -252,19 +253,23 @@ func (d *Daemon) handleSurfaceApprove(w http.ResponseWriter, r *http.Request, o
 		writeApprovalRefusal(w, http.StatusConflict, "Conflict", "request "+id+"'s stored challenge is not 64 hex; make a new request")
 		return
 	}
-	if body.SSHSIG != nil { // R11, sprint 2's fallback cut: the SSHSIG arm is sprint 3's
-		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved",
-			"this daemon records only passkey approvals from the surface; approve an SSH signature at a terminal: ailang-worldd approve "+id+" --signature <file>")
-		return
+	format, artifact, what := approval.FormatWebAuthn, []byte(body.WebAuthn), "the passkey's assertion"
+	if body.SSHSIG != nil {
+		format, artifact, what = approval.FormatSSH, []byte(*body.SSHSIG), "the SSH signature"
+	}
+	var signed *approval.SignedApproval
+	if format == approval.FormatSSH { // R11: the law, before anything is recorded
+		signed, err = authority.JudgeSSHSignature(ctx, ro, row, artifact, now)
+	} else {
+		signed, err = authority.JudgeAssertion(ctx, ro, row, artifact, now)
 	}
-	signed, err := authority.JudgeAssertion(ctx, ro, row, body.WebAuthn, now) // R11: the law
 	if err != nil {
 		var je *authority.JudgeError
 		if errors.As(err, &je) && je.Step == authority.JudgeStepStore {
 			d.approvalStoreError(w, r, err)
 			return
 		}
-		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the passkey's assertion does not approve request "+id+": "+err.Error())
+		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", what+" does not approve request "+id+": "+err.Error())
 		return
 	}
 	digest, _ := hex.DecodeString(row.Digest)                                                  // R8 has shown it is 64 hex
@@ -277,8 +282,8 @@ func (d *Daemon) handleSurfaceApprove(w http.ResponseWriter, r *http.Request, o
 		d.approvalStoreError(w, r, err)
 		return
 	}
-	err = wr.RecordSurfaceApproval(ctx, store.SurfaceApprovalRow{RequestID: id, Format: approval.FormatWebAuthn, Approver: signed.ApprovedBy(),
-		Artifact: body.WebAuthn, Digest: row.Digest}, d.approveNow()) // R14, at the time it is received
+	err = wr.RecordSurfaceApproval(ctx, store.SurfaceApprovalRow{RequestID: id, Format: format, Approver: signed.ApprovedBy(),
+		Artifact: artifact, Digest: row.Digest}, d.approveNow()) // R14, at the time it is received
 	_ = wr.Close()
 	switch {
 	case isRequestRace(err):
--- a/host/daemon/surface_approve_post_test.go
+++ b/host/daemon/surface_approve_post_test.go
@@ -128,8 +128,15 @@ func TestApprovePOSTRefusalArms(t *testing.T) {
 		{"R11-garbage", "", "NotApproved", "does not approve request", 422, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
 			return approveBody(row.Digest, []byte(`{"id":"x"}`)), nil
 		}},
-		{"R11-sshsig-cut", "", "NotApproved", "--signature", 422, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
-			return []byte(`{"v":1,"digest":"` + row.Digest + `","sshsig":"-----BEGIN SSH SIGNATURE-----\n"}`), nil
+		{"R11-sshsig-garbage", "", "NotApproved", "the SSH signature does not approve request", 422, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			return sshBody(row.Digest, []byte("-----BEGIN SSH SIGNATURE-----\n")), nil
+		}},
+		{"R11-sshsig-unenrolled", "", "NotApproved", "no enrolled approver", 422, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			return sshBody(row.Digest, signatureFor(t, approvaltest.NewSigner(74), row)), nil
+		}},
+		{"R11-sshsig-store", "", "Internal", "internal store failure", 500, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			setVersion(t, a.d.credentialsPath, `ALTER TABLE approvers RENAME COLUMN public_key TO public_key_gone`)
+			return sshBody(row.Digest, signatureFor(t, apSigner, row)), nil
 		}},
 		{"R11-store-v5", "", "Restart", "restart `ailang-worldd serve`", 503, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
 			body := approveBody(row.Digest, assertionFor(t, apPasskey, row, 2))
--- a/host/daemon/surface_approve_sshsig_test.go
+++ b/host/daemon/surface_approve_sshsig_test.go
@@ -0,0 +1,129 @@
+package daemon
+
+// Row 201 A3 sprint 3 M5d (design §2.4, §6 RT6, RT13; the sprint-2 plan's Appendix C): the POST's
+// SSHSIG arm judges an approvaltest signature with the A1b law before it records.
+
+import (
+	"bytes"
+	"context"
+	"encoding/hex"
+	"encoding/json"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// apSigner is the SSH approver (mark-ssh) of the SSHSIG tests.
+var apSigner = approvaltest.NewSigner(73)
+
+// enrolSSH enrols s as name at apT0-10: after the tests' requests (apT0-30), as enrolPasskey.
+func (a approvalRig) enrolSSH(name string, s *approvaltest.Signer) {
+	a.t.Helper()
+	a.write(func(ctx context.Context, c *store.Credentials) error {
+		return c.EnrolApprover(ctx, store.ApproverRow{Name: name, Format: approval.FormatSSH, PublicKey: s.PublicKey().Blob, EnrolledAt: apT0 - 10})
+	})
+}
+
+// signatureFor is s's armored SSHSIG over row's 90-byte message.
+func signatureFor(t *testing.T, s *approvaltest.Signer, row store.MintRequestRow) []byte {
+	t.Helper()
+	d, _ := hex.DecodeString(row.Digest)
+	c, _ := hex.DecodeString(row.Challenge)
+	msg, err := approval.SignedMessage([32]byte(d), [32]byte(c), row.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return s.Sign(approval.Namespace, msg)
+}
+
+// sshBody is the page's POST body for digest with an armored SSHSIG.
+func sshBody(digest string, armored []byte) []byte {
+	b, _ := json.Marshal(map[string]any{"v": 1, "digest": digest, "sshsig": string(armored)})
+	return b
+}
+
+// recorded is request id's surface_approvals row.
+func (a approvalRig) recorded(id string) (store.SurfaceApprovalRow, bool) {
+	a.t.Helper()
+	ro, err := store.OpenCredentialsReadOnly(a.d.credentialsPath)
+	if err != nil {
+		a.t.Fatal(err)
+	}
+	defer func() { _ = ro.Close() }()
+	ctx, cancel := context.WithTimeout(context.Background(), 5e9)
+	defer cancel()
+	sa, ok, err := ro.SurfaceApproval(ctx, id)
+	if err != nil {
+		a.t.Fatal(err)
+	}
+	return sa, ok
+}
+
+// RT6: an enrolled key's SSHSIG over the request's message is recorded (ssh, ssh:mark-ssh, the
+// armored bytes, the request's digest, received at the daemon's now); nothing mints.
+func TestApprovePOSTRecordsAVerifiedSSHSIG(t *testing.T) {
+	a, row := newPostRig(t)
+	a.enrolSSH("mark-ssh", apSigner)
+	sig := signatureFor(t, apSigner, row)
+	a.d.approveNow = func() int64 { return apT0 + 4 }
+	rec := a.post(row.RequestID, sshBody(row.Digest, sig))
+	if rec.Code != 200 || !strings.Contains(rec.Body.String(), `"approvedBy":"ssh:mark-ssh"`) {
+		t.Fatalf("POST = %d %s", rec.Code, rec.Body.String())
+	}
+	sa, ok := a.recorded(row.RequestID)
+	if !ok || sa.Format != "ssh" || sa.Approver != "ssh:mark-ssh" || !bytes.Equal(sa.Artifact, sig) || sa.Digest != row.Digest || sa.ReceivedAt != apT0+4 {
+		t.Fatalf("recorded = %+v %v", sa, ok)
+	}
+	if n, s := a.recordedAndSessions(); n != 1 || s != 0 {
+		t.Fatalf("%d recorded, %d sessions; want 1 and 0", n, s)
+	}
+	if body := a.serve("GET", "/surface/requests.json", apLoopback, nil).Body.String(); !strings.Contains(body, `"recorded":"ssh:mark-ssh"`) {
+		t.Fatalf("requests.json does not show the approver: %s", body)
+	}
+}
+
+// RT13: one SSH approver and no passkey; with --surface-origin https://studio.tail0.ts.net the
+// card and the SSHSIG POST work through that origin under both its own Host and the loopback one
+// (a proxy); without it they are refused naming the flag (403 via the loopback Host, 421 via the
+// ts.net Host), and another ts.net origin is refused even when one is configured.
+func TestApprovePOSTSSHSIGViaTailnetOriginWithNoPasskey(t *testing.T) {
+	const ts, tsHost = "https://studio.tail0.ts.net", "studio.tail0.ts.net"
+	rig := func(configured bool) (approvalRig, store.MintRequestRow, []byte) {
+		a := newApprovalRig(t)
+		a.enrolSSH("mark-ssh", apSigner)
+		if configured {
+			a.d.cfg.SurfaceOrigins = []string{ts}
+		}
+		row := a.putRequest("mr-00000000000000d1", "w201d", apT0-30)
+		return a, row, sshBody(row.Digest, signatureFor(t, apSigner, row))
+	}
+	a, row, body := rig(true)
+	for _, host := range []string{tsHost, apLoopback} {
+		if rec := a.serve("GET", "/surface/requests/"+row.RequestID, host, nil); rec.Code != 200 {
+			t.Fatalf("card via Host %s = %d %s", host, rec.Code, rec.Body.String())
+		}
+		rec := a.serve("POST", "/surface/approve/"+row.RequestID, host, body, "Origin", ts, "Sec-Fetch-Site", "same-origin", "Content-Type", "application/json")
+		if sa, ok := a.recorded(row.RequestID); rec.Code != 200 || !ok || sa.Approver != "ssh:mark-ssh" {
+			t.Fatalf("POST via Host %s = %d %s; recorded %+v", host, rec.Code, rec.Body.String(), sa)
+		}
+	}
+	if n, _ := a.recordedAndSessions(); n != 1 {
+		t.Fatalf("%d rows after two POSTs of one request, want 1", n)
+	}
+	other := a.serve("POST", "/surface/approve/"+row.RequestID, apLoopback, body, "Origin", "https://other.tail0.ts.net", "Sec-Fetch-Site", "same-origin", "Content-Type", "application/json")
+	requireRefusal(t, "another ts.net origin", other, 403, "CrossOrigin")
+
+	a, row, body = rig(false)
+	r := requireRefusal(t, "unconfigured", a.serve("POST", "/surface/approve/"+row.RequestID, apLoopback, body, "Origin", ts, "Sec-Fetch-Site", "same-origin", "Content-Type", "application/json"), 403, "CrossOrigin")
+	if !strings.Contains(r.Reason, "--surface-origin "+ts) {
+		t.Fatalf("the refusal %q does not name the flag", r.Reason)
+	}
+	requireRefusal(t, "unconfigured ts.net Host", a.serve("POST", "/surface/approve/"+row.RequestID, tsHost, body, "Origin", ts, "Sec-Fetch-Site", "same-origin", "Content-Type", "application/json"), 421, "MisdirectedRequest")
+	requireRefusal(t, "unconfigured card", a.serve("GET", "/surface/requests/"+row.RequestID, tsHost, nil), 421, "MisdirectedRequest")
+	if n, s := a.recordedAndSessions(); n != 0 || s != 0 {
+		t.Fatalf("%d recorded, %d sessions unconfigured; want 0 and 0", n, s)
+	}
+}
````

### Patch M6b — the card page and its routes

````diff M6b
--- a/host/daemon/daemon.go
+++ b/host/daemon/daemon.go
@@ -1098,7 +1098,9 @@ func releaseFromVersion(version string) string {
 // read routes outside the frozen table, as are GET /surface/follow-status.json, GET /surface/model.json and
 // GET /surface/stack.json (the app's view model, N1) and GET /surface with /surface/app.js,
 // /surface/app.css and /surface/ailang-logo.svg (the app, N2b), and GET /surface/passkey with
-// /surface/passkey.js (the passkey ceremony page, row 201 A2: static, it reads nothing).
+// /surface/passkey.js (the passkey ceremony page, row 201 A2: static, it reads nothing), and the
+// approval routes of row 201 A3 (requests.json, requests/{id}, the approve POST, the card page
+// GET /surface/approve/{id} and its approve.js), each behind approvalRoute's Host check.
 // The registry pattern deliberately uses a multi-segment wildcard: registry semantic IDs
 // such as "world/epoch-registry/v1" contain slashes.
 func (d *Daemon) Handler() http.Handler {
@@ -1135,6 +1137,8 @@ func (d *Daemon) Handler() http.Handler {
 	mux.HandleFunc("GET /surface/requests.json", d.approvalRoute(d.handleSurfaceRequests))
 	mux.HandleFunc("GET /surface/requests/{id}", d.approvalRoute(d.handleSurfaceRequestCard))
 	mux.HandleFunc("POST /surface/approve/{id}", d.approvalRoute(d.handleSurfaceApprove))
+	mux.HandleFunc("GET /surface/approve/{id}", d.approvalRoute(d.handleSurfaceApprovePage))
+	mux.HandleFunc("GET /surface/approve.js", d.approvalRoute(d.handleSurfaceApproveScript))
 	// The two A2A projection routes (w-a2a-session-projection P6.B-A2A-CARD)
 	// are ADDITIVE: the frozen /v1/ table above is untouched, and the routes
 	// are NOT in isProtected — the projection handler resolves the session
--- a/host/daemon/surface_app_test.go
+++ b/host/daemon/surface_app_test.go
@@ -17,8 +17,9 @@ import (
 // same bytes whatever the index holds), its assets are the exact embedded
 // bytes, and every other name under /surface/ is a 404 (row 201 A2 adds
 // /surface/passkey and /surface/passkey.js, PG1's, and no name under them; row 201 A3
-// sprint 2 adds the approval routes /surface/requests.json and /surface/requests/{id},
-// which answer only an allowed Host, and /surface/approve.js stays a 404 until sprint 3).
+// sprint 2 adds the approval routes /surface/requests.json and /surface/requests/{id}, and
+// sprint 3 the card page /surface/approve/{id} and /surface/approve.js, which answer only an
+// allowed Host).
 func TestSurfaceAppRoutes(t *testing.T) {
 	path := filepath.Join(t.TempDir(), "operator-index.json")
 	d, err := New(boundedTestContext(t), Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: DefaultBindHost, SurfaceIndex: path})
@@ -153,12 +154,12 @@ func TestSurfaceAppRoutes(t *testing.T) {
 		for _, target := range []string{"/surface/", "/surface/nope", "/surface/app.jsx", "/surface/index.html", "/surface/data.js",
 			"/surface/spike-data.json", "/surface/app/app.js", "/surface/static/x", "/surface/fonts.css",
 			"/surface/passkey/", "/surface/passkey/index.html", "/surface/passkey/passkey.js", "/surface/passkey.html",
-			"/surface/approve.js", "/surface/approve", "/surface/requests", "/surface/requests.json/x"} {
+			"/surface/approve", "/surface/approve/", "/surface/requests", "/surface/requests.json/x"} {
 			if w := request("GET", target); w.Code != 404 {
 				t.Errorf("%s = %d, want 404", target, w.Code)
 			}
 		}
-		for _, target := range []string{"/surface/requests.json", "/surface/requests/mr-0123456789abcdef"} {
+		for _, target := range []string{"/surface/requests.json", "/surface/requests/mr-0123456789abcdef", "/surface/approve/mr-0123456789abcdef", "/surface/approve.js"} {
 			if w := request("GET", target); w.Code != 421 { // registered; example.com is no Host of the surface's
 				t.Errorf("%s = %d, want 421", target, w.Code)
 			}
--- a/host/daemon/surface_approve_page.go
+++ b/host/daemon/surface_approve_page.go
@@ -0,0 +1,29 @@
+package daemon
+
+// Row 201 A3 sprint 3 (design §2.7): the request card page and its script. Static embedded bytes
+// under surfaceCSP, behind approvalRoute's Host check like every approval route; the page reads the
+// request itself from GET /surface/requests/{id}.
+
+import (
+	"net/http"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// handleSurfaceApprovePage is GET /surface/approve/{id}: the same bytes for every well-formed id.
+func (d *Daemon) handleSurfaceApprovePage(w http.ResponseWriter, r *http.Request, _ []string) {
+	if !store.MintRequestIDPattern.MatchString(r.PathValue("id")) {
+		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "a request id is mr- and 16 lowercase hex")
+		return
+	}
+	setSurfaceHeaders(w)
+	_, _ = w.Write(surfacepage.ApprovePage)
+}
+
+// handleSurfaceApproveScript is the card's script, under the same headers.
+func (d *Daemon) handleSurfaceApproveScript(w http.ResponseWriter, _ *http.Request, _ []string) {
+	setSurfaceHeaders(w)
+	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
+	_, _ = w.Write(surfacepage.ApproveScript)
+}
--- a/host/daemon/surface_approve_page_test.go
+++ b/host/daemon/surface_approve_page_test.go
@@ -0,0 +1,74 @@
+package daemon
+
+// Row 201 A3 sprint 3 M6b (design §2.7, §6 PG2): the request card page and its script.
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"encoding/hex"
+	"os"
+	"path/filepath"
+	"regexp"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// PG2: GET /surface/approve/{id} and /surface/approve.js serve the embedded bytes, byte-equal to
+// host/surfacepage/approve/, under the unchanged surfaceCSP and the surface headers; a malformed id
+// is 400; the A2 page's bytes are unchanged. The page loads exactly its own script.
+func TestApprovePageRoutes(t *testing.T) {
+	const csp = "default-src 'none'; style-src 'unsafe-inline'; script-src 'self'; connect-src 'self'; font-src 'self'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'"
+	if surfaceCSP != csp {
+		t.Fatalf("surfaceCSP changed: %q", surfaceCSP)
+	}
+	a := newApprovalRig(t)
+	for _, tc := range []struct{ target, file, contentType string }{
+		{"/surface/approve/mr-0123456789abcdef", "index.html", "text/html; charset=utf-8"},
+		{"/surface/approve/mr-00000000000000b1", "index.html", "text/html; charset=utf-8"},
+		{"/surface/approve.js", "approve.js", "text/javascript; charset=utf-8"},
+	} {
+		want, err := os.ReadFile(filepath.Join("..", "surfacepage", "approve", tc.file))
+		if err != nil {
+			t.Fatal(err)
+		}
+		rec := a.serve("GET", tc.target, apLoopback, nil)
+		if rec.Code != 200 || len(want) == 0 || !bytes.Equal(rec.Body.Bytes(), want) {
+			t.Errorf("GET %s: %d, or not the bytes of host/surfacepage/approve/%s", tc.target, rec.Code, tc.file)
+		}
+		for k, v := range map[string]string{"Content-Security-Policy": csp, "Content-Type": tc.contentType,
+			"Cache-Control": "no-store", "X-Content-Type-Options": "nosniff", "Referrer-Policy": "no-referrer"} {
+			if got := rec.Header().Get(k); got != v {
+				t.Errorf("GET %s: %s = %q, want %q", tc.target, k, got, v)
+			}
+		}
+		if rec.Header().Get("Set-Cookie") != "" || rec.Header().Get("Access-Control-Allow-Origin") != "" {
+			t.Errorf("GET %s sets a cookie or a CORS header", tc.target)
+		}
+	}
+	for _, id := range []string{"mr-0123", "mr-0123456789ABCDEF", "x"} {
+		requireRefusal(t, id, a.serve("GET", "/surface/approve/"+id, apLoopback, nil), 400, "BadRequest")
+	}
+	for file, sum := range map[string]string{"index.html": "1695ebc2ae20f303dbfa3b617b3bec4c76035447ec290deba91e695fec9cc06d",
+		"passkey.js": "0064c6c90099a3565b178cf170618b80a67bbf1f3444a1018cc4ee29f51ed53f"} {
+		b, err := os.ReadFile(filepath.Join("..", "surfacepage", "passkey", file))
+		if h := sha256.Sum256(b); err != nil || hex.EncodeToString(h[:]) != sum {
+			t.Errorf("the A2 page's %s changed (%v)", file, err)
+		}
+	}
+	page := string(surfacepage.ApprovePage)
+	if scripts := regexp.MustCompile(`(?i)<script\b[^>]*>`).FindAllString(page, -1); len(scripts) != 1 || !strings.Contains(page, `<script src="/surface/approve.js"></script>`) {
+		t.Errorf("scripts %q", scripts)
+	}
+	if m := regexp.MustCompile(`(?i)\son[a-z]+\s*=|<form|<img|<link|<iframe|https?://`).FindString(page); m != "" {
+		t.Errorf("the page has %q", m)
+	}
+	// Labelled instrument-health controls (not load-bearing; XA1 and XA2 are): the script writes no
+	// markup, evaluates no string and reads no cookie.
+	for _, banned := range []string{"innerHTML", "outerHTML", "insertAdjacentHTML", "document.write", "eval(", "new Function", "document.cookie", "XMLHttpRequest", "sendBeacon", "WebSocket"} {
+		if bytes.Contains(surfacepage.ApproveScript, []byte(banned)) {
+			t.Errorf("approve.js uses %s", banned)
+		}
+	}
+}
--- a/host/daemon/surface_approve_test.go
+++ b/host/daemon/surface_approve_test.go
@@ -15,12 +15,12 @@ import (
 	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
 )
 
-// approvalRoutes is every approval route Handler registers this sprint (the card page and its
-// script are sprint 3's).
-var approvalRoutes = []string{"GET /surface/requests.json", "GET /surface/requests/{id}", "POST /surface/approve/{id}"}
+// approvalRoutes is every approval route Handler registers (sprint 3 adds the card page and its
+// script).
+var approvalRoutes = []string{"GET /surface/approve.js", "GET /surface/approve/{id}", "GET /surface/requests.json", "GET /surface/requests/{id}", "POST /surface/approve/{id}"}
 
 // RT1: exactly the approval routes are registered, under /surface, each behind approvalRoute; no
-// other method is answered on them; the card page and its script are not served yet.
+// other method is answered on them; no other name under them is served.
 func TestApprovalRoutesRegistered(t *testing.T) {
 	src, err := os.ReadFile("daemon.go")
 	if err != nil {
@@ -50,9 +50,9 @@ func TestApprovalRoutesRegistered(t *testing.T) {
 			}
 		}
 	}
-	for _, target := range []string{"/surface/approve.js", "/surface/approve", "/surface/requests", "/surface/requests/mr-0123456789abcdef/x"} {
+	for _, target := range []string{"/surface/approve", "/surface/approve/", "/surface/approve.jsx", "/surface/approve/mr-0123456789abcdef/x", "/surface/requests", "/surface/requests/mr-0123456789abcdef/x"} {
 		if rec := a.serve("GET", target, apLoopback, nil); rec.Code != 404 {
-			t.Errorf("GET %s = %d, want 404 (not served this sprint)", target, rec.Code)
+			t.Errorf("GET %s = %d, want 404", target, rec.Code)
 		}
 	}
 }
@@ -63,7 +63,7 @@ func TestApprovalRoutesRegistered(t *testing.T) {
 func TestApprovalRoutesRefuseForeignHost(t *testing.T) {
 	a := newApprovalRig(t)
 	row := a.putRequest("mr-00000000000000b1", "w201b", apT0-30)
-	targets := []string{"/surface/requests.json", "/surface/requests/" + row.RequestID}
+	targets := []string{"/surface/requests.json", "/surface/requests/" + row.RequestID, "/surface/approve/" + row.RequestID, "/surface/approve.js"}
 	for _, target := range targets {
 		for _, host := range []string{"rebind.example:7644", "localhost:7645", "localhost", "127.0.0.1:7644.evil", "studio.tail0.ts.net", "cfg.tail0.ts.net"} {
 			requireRefusal(t, host+target, a.serve("GET", target, host, nil), 421, "MisdirectedRequest")
--- a/host/surfacepage/approve.go
+++ b/host/surfacepage/approve.go
@@ -0,0 +1,19 @@
+package surfacepage
+
+import _ "embed"
+
+// Row 201 A3 (design_docs/planned/w-world-approvals-by-key-a3.md §2.7): the request card, served
+// at GET /surface/approve/{id} and GET /surface/approve.js. The page reads the request from GET
+// /surface/requests/{id}, shows every signed field of its stored bytes after checking their digest,
+// computes the message it signs itself and POSTs a passkey assertion or an SSH signature to POST
+// /surface/approve/{id}; the daemon judges it. The bytes are the same for every id.
+
+// ApprovePage is GET /surface/approve/{id}.
+//
+//go:embed approve/index.html
+var ApprovePage []byte
+
+// ApproveScript is GET /surface/approve.js.
+//
+//go:embed approve/approve.js
+var ApproveScript []byte
--- a/host/surfacepage/approve/approve.js
+++ b/host/surfacepage/approve/approve.js
@@ -0,0 +1,316 @@
+'use strict';
+// AILANG World's request card (row 201 A3, design §2.7). It reads one mint request from the daemon,
+// checks that the stored bytes hash to the request's digest, shows EVERY signed field of those bytes
+// (and refuses a request with a field it cannot show), computes the 90-byte message and its hash
+// itself (never a hash from the server), and approves with a passkey or a pasted SSH signature by a
+// POST the daemon judges with the A2 law. Text is set with textContent only.
+(() => {
+  const TIMEOUT_MS = 240000; // A2's ceremony bound (A2 §2.14 [r2])
+  const NET_MS = 10000; // [r2] every GET and POST, its body included
+  const RETRIES = 3, RETRY_MS = 1000; // [r2] only after an explicit 503 Busy
+  const NAMESPACE = 'world-approval@sunholo';
+  // Every JSON key path of approval.MintRequest ([] = each element), with its label, group and type.
+  // TestApprovePageRendersEveryCanonicalField compares this list with the Go struct.
+  const FIELDS = /* FIELDS */[
+    {"path": "episode", "label": "Episode", "group": "what", "type": "string"},
+    {"path": "agent", "label": "Agent", "group": "what", "type": "string"},
+    {"path": "model", "label": "Model", "group": "what", "type": "string"},
+    {"path": "task", "label": "Task file", "group": "what", "type": "string"},
+    {"path": "grants[].effect", "label": "Grant effect", "group": "what", "type": "string"},
+    {"path": "grants[].scope", "label": "Grant scope", "group": "what", "type": "string"},
+    {"path": "grants[].budget", "label": "Grant budget", "group": "what", "type": "int"},
+    {"path": "ttl_s", "label": "Session lifetime (s)", "group": "what", "type": "int"},
+    {"path": "worktree", "label": "Worktree", "group": "where", "type": "string"},
+    {"path": "reuse", "label": "Worktree use", "group": "where", "type": "bool", "no": "creates a new worktree", "yes": "uses an existing worktree"},
+    {"path": "degraded", "label": "Mode", "group": "where", "type": "bool", "no": "git worktree", "yes": "NON-GIT, DEGRADED"},
+    {"path": "repo", "label": "Repository", "group": "where", "type": "string"},
+    {"path": "branch", "label": "Branch", "group": "where", "type": "string"},
+    {"path": "git", "label": "Git binary", "group": "where", "type": "string"},
+    {"path": "workspace_root", "label": "Workspace root", "group": "where", "type": "string"},
+    {"path": "profile.project", "label": "Exec profile", "group": "bound", "type": "string"},
+    {"path": "profile.digest", "label": "Profile digest", "group": "bound", "type": "string"},
+    {"path": "profiles_dir", "label": "Profiles directory", "group": "bound", "type": "string"},
+    {"path": "world_db", "label": "World store", "group": "stores", "type": "string"},
+    {"path": "credentials_db", "label": "Credentials store", "group": "stores", "type": "string"},
+    {"path": "session_file", "label": "Token written to", "group": "stores", "type": "string"},
+    {"path": "mcp_config", "label": "MCP config", "group": "stores", "type": "string"},
+    {"path": "daemon.addr", "label": "Daemon address", "group": "daemon", "type": "string"},
+    {"path": "daemon.vcs_revision", "label": "Daemon revision", "group": "daemon", "type": "string"},
+    {"path": "daemon.vcs_modified", "label": "Daemon build", "group": "daemon", "type": "bool", "no": "clean build", "yes": "modified build"},
+    {"path": "daemon.executable_sha256", "label": "Daemon binary SHA-256", "group": "daemon", "type": "string"},
+    {"path": "daemon.go_version", "label": "Daemon Go version", "group": "daemon", "type": "string"},
+    {"path": "schema", "label": "Request schema", "group": "daemon", "type": "string"}
+  ]/* /FIELDS */;
+  const GROUPS = [['what', 'What runs'], ['where', 'Where'], ['bound', 'Bound to'], ['stores', 'Stores and files'], ['daemon', 'Requesting daemon']];
+  const TYPES = { string: 'string', bool: 'boolean', int: 'number' };
+
+  const $ = (id) => document.getElementById(id);
+  const setState = (s) => { document.body.dataset.state = s; };
+  const say = (text, error) => {
+    $('status').textContent = text;
+    $('status').toggleAttribute('data-error', !!error);
+  };
+  const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
+  const b64u = (buf) => {
+    let s = '';
+    for (const b of new Uint8Array(buf)) s += String.fromCharCode(b);
+    return btoa(s).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
+  };
+  const unb64u = (s) => {
+    if (typeof s !== 'string' || !/^[A-Za-z0-9_-]*$/.test(s)) throw new Error('not base64url');
+    return Uint8Array.from(atob(s.replace(/-/g, '+').replace(/_/g, '/')), (c) => c.charCodeAt(0));
+  };
+  const hex = (b) => Array.from(b, (x) => x.toString(16).padStart(2, '0')).join('');
+  const unhex = (s) => Uint8Array.from(s.match(/../g), (x) => parseInt(x, 16));
+  const sha256 = async (bytes) => new Uint8Array(await crypto.subtle.digest('SHA-256', bytes));
+  // A2 §2.7's assertion file, built here rather than by toJSON() (as passkey.js does).
+  const assertionJSON = (c) => ({ id: c.id, rawId: b64u(c.rawId), type: c.type,
+    response: { clientDataJSON: b64u(c.response.clientDataJSON), authenticatorData: b64u(c.response.authenticatorData),
+      signature: b64u(c.response.signature), userHandle: c.response.userHandle ? b64u(c.response.userHandle) : null } });
+
+  // [r2] Each GET and POST, response body included, has its own 10 s deadline; leaving the page
+  // cancels every outstanding one.
+  const outstanding = new Set();
+  addEventListener('pagehide', () => { for (const c of outstanding) c.abort(); });
+  async function net(url, init) {
+    const c = new AbortController();
+    outstanding.add(c);
+    const timer = setTimeout(() => c.abort(), NET_MS);
+    try {
+      const r = await fetch(url, Object.assign({ cache: 'no-store', signal: c.signal }, init));
+      const text = await r.text();
+      let body = null;
+      try { body = JSON.parse(text); } catch (_) { body = null; }
+      return { status: r.status, body };
+    } finally {
+      clearTimeout(timer);
+      outstanding.delete(c);
+    }
+  }
+  const why = (e) => (e && e.name === 'AbortError' ? `no answer within ${NET_MS / 1000} s` : String(e && e.message || e));
+
+  // leaves is every leaf of v: its key path pattern ([] for an element), its concrete path (the
+  // element's index) and its value. An empty object or array has none, so its fields are missing.
+  function leaves(v, pattern, concrete, out, index) {
+    if (Array.isArray(v)) {
+      v.forEach((e, i) => leaves(e, pattern + '[]', concrete + '.' + i, out, i));
+    } else if (v !== null && typeof v === 'object') {
+      for (const k of Object.keys(v)) leaves(v[k], pattern ? pattern + '.' + k : k, concrete ? concrete + '.' + k : k, out, index);
+    } else {
+      out.push({ pattern, concrete, value: v, index });
+    }
+    return out;
+  }
+  // unshown is the first key path the card cannot show (unlisted, of another type, or missing), or ''.
+  // A request that is not an object has a leaf at '' or under '[]', which no FIELDS entry names.
+  function unshown(ls) {
+    const byPath = new Map(FIELDS.map((f) => [f.path, f]));
+    const seen = new Set();
+    for (const l of ls) {
+      const f = byPath.get(l.pattern);
+      if (!f || typeof l.value !== TYPES[f.type] || (f.type === 'int' && !Number.isSafeInteger(l.value))) return l.pattern;
+      seen.add(l.pattern);
+    }
+    for (const f of FIELDS) if (!seen.has(f.path)) return f.path;
+    return '';
+  }
+  const shown = (f, v) => (f.type === 'bool' ? (v ? f.yes : f.no) : String(v));
+
+  function render(ls) {
+    const card = $('card');
+    for (const [group, title] of GROUPS) {
+      const section = document.createElement('section');
+      section.dataset.group = group;
+      const h = document.createElement('h2');
+      h.textContent = title;
+      const dl = document.createElement('dl');
+      for (const f of FIELDS.filter((x) => x.group === group)) {
+        for (const l of ls.filter((x) => x.pattern === f.path)) {
+          const dt = document.createElement('dt');
+          dt.textContent = f.label + (l.index === undefined ? '' : ' ' + (l.index + 1)) + (f.path === 'ttl_s' ? ` (${(l.value / 3600).toFixed(1)} h)` : '');
+          const dd = document.createElement('dd');
+          dd.id = 'f-' + l.concrete;
+          dd.textContent = shown(f, l.value);
+          dl.append(dt, dd);
+        }
+      }
+      section.append(h, dl);
+      card.append(section);
+    }
+    card.hidden = false;
+  }
+
+  const id = (location.pathname.match(/^\/surface\/approve\/(mr-[0-9a-f]{16})$/) || [])[1];
+  let card = null, digest = '', message = null, challenge = null, group = null, again = null;
+  const buttons = (on) => { $('go').disabled = !on; $('post-sig').disabled = !on; };
+  const offer = (label, action) => { again = action; $('again').textContent = label; $('again').hidden = false; $('again').disabled = false; };
+  const refuse = (text) => { setState('refused'); $('what').textContent = text; };
+
+  async function load() {
+    $('again').hidden = true;
+    setState('loading');
+    let res;
+    try {
+      res = await net('/surface/requests/' + id, { headers: { Accept: 'application/json' } });
+    } catch (e) {
+      setState('failed');
+      say(`The request could not be read (${why(e)}). Nothing to approve yet.`, true);
+      return offer('Read it again', load);
+    }
+    if (res.status !== 200 || !res.body) {
+      setState('failed');
+      say(`The daemon did not give the request (HTTP ${res.status}${res.body && res.body.reason ? ': ' + res.body.reason : ''}).`, true);
+      return offer('Read it again', load);
+    }
+    card = res.body;
+    let bytes, obj;
+    try {
+      bytes = unb64u(card.request);
+    } catch (e) {
+      return refuse('The daemon\'s request is not base64url; nothing to approve.');
+    }
+    if (card.id !== id || hex(await sha256(bytes)) !== card.digest) {
+      return refuse('The daemon\'s request does not match its digest; nothing to approve.');
+    }
+    try {
+      obj = JSON.parse(new TextDecoder('utf-8', { fatal: true }).decode(bytes));
+    } catch (e) {
+      return refuse(`The daemon's request is not UTF-8 JSON (${e.message}); nothing to approve.`);
+    }
+    const ls = leaves(obj, '', '', []);
+    const bad = unshown(ls);
+    if (bad) {
+      return refuse(`This request has a field the page cannot show (${bad}); nothing to approve; approve it at a terminal.`);
+    }
+    if (!/^[0-9a-f]{64}$/.test(card.challenge) || !Number.isSafeInteger(card.expires) || card.expires < 0) {
+      return refuse('The daemon\'s challenge or expiry is malformed; nothing to approve.');
+    }
+    // The 90-byte message (design §2.7 step 3): "world-approval/v1" ‖ 0x00 ‖ digest ‖ nonce ‖ u64be(expires).
+    digest = card.digest;
+    message = new Uint8Array(90);
+    message.set(new TextEncoder().encode('world-approval/v1'), 0);
+    message.set(unhex(card.digest), 18);
+    message.set(unhex(card.challenge), 50);
+    new DataView(message.buffer).setBigUint64(82, BigInt(card.expires));
+    challenge = await sha256(message);
+    render(ls);
+    $('what').textContent = `Approve request ${id} for episode ${obj.episode}, expires ${new Date(card.expires * 1000).toLocaleString()}, message hash ${hex(challenge.slice(0, 8))}`;
+    if (card.consumed) {
+      setState('used');
+      say(`Already used: approved by ${card.approvedBy}.`, false);
+      return;
+    }
+    group = (card.passkeys || []).find((p) => p.origin === location.origin) || null;
+    if (group) {
+      $('go').hidden = false;
+    } else {
+      const others = (card.passkeys || []).map((p) => p.origin);
+      $('origin-note').textContent = others.length ? `No passkey is enrolled for ${location.origin}: open this card at ${others.join(' or ')}.` : 'No passkey is enrolled; approve with an SSH key below.';
+      $('origin-note').hidden = false;
+    }
+    const save = $('save');
+    save.href = URL.createObjectURL(new Blob([message], { type: 'application/octet-stream' }));
+    save.download = `world-approval-${id}.msg`;
+    save.hidden = false;
+    $('sshline').textContent = `ssh-keygen -Y sign -n ${NAMESPACE} -f <key> world-approval-${id}.msg`;
+    $('ssh').hidden = false;
+    buttons(true);
+    setState('ready');
+  }
+
+  // status reads the request's state from the decision list: one GET, under its deadline.
+  async function status() {
+    $('again').disabled = true;
+    let res;
+    try {
+      res = await net('/surface/requests.json', { headers: { Accept: 'application/json' } });
+    } catch (e) {
+      say(`The status could not be read (${why(e)}).`, true);
+      return offer('Check again', status);
+    }
+    const b = res.body || {};
+    const p = (b.pending || []).find((x) => x.id === id);
+    $('again').hidden = true;
+    if (p && p.recorded) {
+      setState('recorded');
+      say(`Recorded: approved by ${p.recorded}, waiting for the minter (approve ${id} --await).`, false);
+    } else if (p) {
+      setState('ready');
+      say('Not recorded: the request is still pending with no approval. You may approve it again.', true);
+      buttons(true);
+    } else {
+      say(res.status === 200 ? 'The request is no longer pending (used or expired): read the card again.' : `The status could not be read (HTTP ${res.status}).`, true);
+      offer('Check again', status);
+    }
+  }
+
+  async function post(body) {
+    buttons(false);
+    setState('posting');
+    say('Recording the approval…', false);
+    for (let attempt = 0; ; attempt++) {
+      let res;
+      try {
+        res = await net('/surface/approve/' + id, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) });
+      } catch (e) {
+        // [r2] No answer: the daemon may or may not have recorded it. Never retried, never another format.
+        setState('unknown');
+        say(`Approval recording outcome unknown: ${why(e)}. It may or may not have been recorded; nothing was minted by this page. Check its status before approving again.`, true);
+        return offer('Check status', status);
+      }
+      const b = res.body || {};
+      if (res.status === 503 && b.error === 'Busy' && attempt < RETRIES) {
+        say(`The credentials store is busy; trying again (${attempt + 1} of ${RETRIES})…`, false);
+        await sleep(RETRY_MS);
+        continue;
+      }
+      if (res.status === 200) {
+        setState('recorded');
+        say(`Approved by ${b.approvedBy}: ${b.next}.`, false);
+        return;
+      }
+      setState('rejected');
+      say(`Not recorded (HTTP ${res.status} ${b.error || ''}): ${b.reason || 'no reason given'}`, true);
+      buttons(true);
+      return;
+    }
+  }
+
+  async function approvePasskey() {
+    buttons(false);
+    setState('pending');
+    say('Waiting for your passkey…', false);
+    const controller = new AbortController();
+    let expired = false;
+    const timer = setTimeout(() => {
+      expired = true;
+      controller.abort();
+      setState('expired');
+      say(`Timed out: no passkey answer within ${TIMEOUT_MS / 60000} minutes. Nothing was sent; try again.`, true);
+      buttons(true);
+    }, TIMEOUT_MS);
+    let cred = null, failure = null;
+    try {
+      cred = await navigator.credentials.get({ publicKey: { challenge, rpId: location.hostname,
+        allowCredentials: group.allow.map((a) => ({ type: 'public-key', id: unb64u(a) })),
+        userVerification: 'required', timeout: TIMEOUT_MS }, signal: controller.signal });
+    } catch (e) {
+      failure = e;
+    } finally {
+      clearTimeout(timer);
+    }
+    if (expired) return; // a result after the timer is ignored: nothing is sent
+    if (failure !== null) {
+      setState('failed');
+      say(`The browser did not make a passkey response (${failure.name}: ${failure.message}). Nothing was sent; try again.`, true);
+      buttons(true);
+      return;
+    }
+    await post({ v: 1, digest, webauthn: assertionJSON(cred) });
+  }
+
+  $('go').addEventListener('click', () => { approvePasskey(); });
+  $('post-sig').addEventListener('click', () => { post({ v: 1, digest, sshsig: $('sig').value }); });
+  $('again').addEventListener('click', () => { if (again) again(); });
+  load(); // the daemon serves this page only at /surface/approve/mr-<16 hex>
+})();
--- a/host/surfacepage/approve/index.html
+++ b/host/surfacepage/approve/index.html
@@ -0,0 +1,41 @@
+<!doctype html>
+<html lang="en">
+<head>
+<meta charset="utf-8">
+<meta name="viewport" content="width=device-width, initial-scale=1">
+<title>AILANG World: approve a mint request</title>
+<style>
+body { font: 16px/1.5 system-ui, sans-serif; max-width: 46rem; margin: 2rem auto; padding: 0 1rem; color: #1d1d1f; }
+h1 { font-size: 1.4rem; }
+h2 { font-size: 1.05rem; margin: 1.2rem 0 0.3rem; }
+#what, #status { white-space: pre-wrap; }
+#status[data-error] { color: #b00020; font-weight: 600; }
+dl { display: grid; grid-template-columns: max-content 1fr; gap: 0.15rem 1rem; margin: 0; }
+dt { color: #555; }
+dd { margin: 0; font: 14px ui-monospace, monospace; overflow-wrap: anywhere; }
+dd:empty::after { content: "(none)"; color: #888; }
+section[data-group="what"] dd { font-size: 16px; font-weight: 600; }
+#sshline { font: 13px ui-monospace, monospace; overflow-wrap: anywhere; }
+#sig { font: 12px ui-monospace, monospace; width: 100%; height: 7rem; }
+button, #save { font-size: 1rem; padding: 0.4rem 1rem; }
+[hidden] { display: none !important; }
+</style>
+</head>
+<body data-state="loading">
+<h1>Approve a mint request</h1>
+<p id="what">Reading the request…</p>
+<div id="card" hidden></div>
+<p id="origin-note" hidden></p>
+<p><button id="go" type="button" hidden disabled>Approve with passkey</button></p>
+<p id="status"></p>
+<p><button id="again" type="button" hidden>Check again</button></p>
+<div id="ssh" hidden>
+<h2>Or approve with an SSH key</h2>
+<p><a id="save" hidden>Save the message</a></p>
+<p id="sshline"></p>
+<p><label for="sig">Paste the signature</label><br><textarea id="sig"></textarea></p>
+<p><button id="post-sig" type="button" disabled>Approve with this signature</button></p>
+</div>
+<script src="/surface/approve.js"></script>
+</body>
+</html>
````

### Patch M6c — every signed field has a renderer

````diff M6c
--- a/host/surfacepage/approve_page_test.go
+++ b/host/surfacepage/approve_page_test.go
@@ -0,0 +1,155 @@
+package surfacepage
+
+// Row 201 A3 sprint 3 M6c (design §2.7, §5.2, §6 PG3, PG3b): the card's field table is the
+// canonical request's, field for field.
+
+import (
+	"bytes"
+	"encoding/json"
+	"reflect"
+	"sort"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+// pageField is one entry of approve.js's FIELDS.
+type pageField struct {
+	Path, Label, Group, Type, No, Yes string
+}
+
+// pageFields is approve.js's FIELDS literal, between its /* FIELDS */ markers, parsed as JSON.
+func pageFields(t *testing.T) []pageField {
+	t.Helper()
+	_, rest, ok := bytes.Cut(ApproveScript, []byte("/* FIELDS */"))
+	lit, _, ok2 := bytes.Cut(rest, []byte("/* /FIELDS */"))
+	var out []pageField
+	if !ok || !ok2 || json.Unmarshal(lit, &out) != nil {
+		t.Fatalf("approve.js has no FIELDS JSON between its markers")
+	}
+	return out
+}
+
+// structFields is every JSON key path of t's json tags ([] for a slice's element) with its kind,
+// and every tag with an option (",omitempty" or any other).
+func structFields(typ reflect.Type, prefix string, paths map[string]string, options *[]string) {
+	for i := 0; i < typ.NumField(); i++ {
+		f := typ.Field(i)
+		name, opts, _ := strings.Cut(f.Tag.Get("json"), ",")
+		if opts != "" {
+			*options = append(*options, typ.Name()+"."+f.Name+": "+opts)
+		}
+		p, ft := prefix+name, f.Type
+		if ft.Kind() == reflect.Slice {
+			p, ft = p+"[]", ft.Elem()
+		}
+		switch ft.Kind() {
+		case reflect.Struct:
+			structFields(ft, p+".", paths, options)
+		case reflect.String:
+			paths[p] = "string"
+		case reflect.Bool:
+			paths[p] = "bool"
+		case reflect.Int, reflect.Int64:
+			paths[p] = "int"
+		default:
+			paths[p] = "unsupported " + ft.Kind().String()
+		}
+	}
+}
+
+// PG3: FIELDS names exactly approval.MintRequest's JSON key paths, each with its type, a label and
+// a group the page shows (a boolean with both of its words). Control: the struct has 21 top-level
+// keys, the count the browser measured in the stored bytes (design V24).
+func TestApprovePageRendersEveryCanonicalField(t *testing.T) {
+	want := map[string]string{}
+	var options []string
+	structFields(reflect.TypeOf(approval.MintRequest{}), "", want, &options)
+	if n := reflect.TypeOf(approval.MintRequest{}).NumField(); n != 21 {
+		t.Fatalf("control: MintRequest has %d top-level keys, want 21 (design V24)", n)
+	}
+	got := map[string]string{}
+	groups := map[string]bool{"what": true, "where": true, "bound": true, "stores": true, "daemon": true}
+	for _, f := range pageFields(t) {
+		if _, dup := got[f.Path]; dup {
+			t.Errorf("FIELDS lists %s twice", f.Path)
+		}
+		got[f.Path] = f.Type
+		if f.Label == "" || !groups[f.Group] || (f.Type == "bool") != (f.No != "" && f.Yes != "") {
+			t.Errorf("FIELDS entry %+v: a label, a shown group, and both words for a boolean only", f)
+		}
+	}
+	if !reflect.DeepEqual(got, want) {
+		var missing, extra []string
+		for p, k := range want {
+			if got[p] != k {
+				missing = append(missing, p+" "+k)
+			}
+		}
+		for p, k := range got {
+			if want[p] != k {
+				extra = append(extra, p+" "+k)
+			}
+		}
+		sort.Strings(missing)
+		sort.Strings(extra)
+		t.Fatalf("FIELDS does not match approval.MintRequest: no renderer for %q; not in the struct %q", missing, extra)
+	}
+}
+
+// leafPaths is every leaf key path of a decoded JSON value ([] for an element).
+func leafPaths(v any, prefix string, out map[string]bool) {
+	switch x := v.(type) {
+	case map[string]any:
+		for k, e := range x {
+			p := k
+			if prefix != "" {
+				p = prefix + "." + k
+			}
+			leafPaths(e, p, out)
+		}
+	case []any:
+		for _, e := range x {
+			leafPaths(e, prefix+"[]", out)
+		}
+	default:
+		out[prefix] = true
+	}
+}
+
+// PG3b: no tag of the canonical request carries an option (so no omitempty: design V28), and a
+// degraded, non-git request with reuse false and every string empty still serializes all 21 keys
+// and every FIELDS path, so the page's missing-path refusal never fires on a canonical request;
+// branch, which XA2 arm (c) deletes, is one of them.
+func TestApprovePageFieldTagControls(t *testing.T) {
+	paths := map[string]string{}
+	var options []string
+	structFields(reflect.TypeOf(approval.MintRequest{}), "", paths, &options)
+	if len(options) != 0 {
+		t.Fatalf("tags with options on the canonical request: %q", options)
+	}
+	b, err := approval.Canonical(approval.MintRequest{Schema: approval.RequestSchema, Episode: "w201d", Degraded: true,
+		Grants: []approval.Grant{{Effect: "Workspace.Read"}}})
+	if err != nil {
+		t.Fatal(err)
+	}
+	var top map[string]any
+	if err := json.Unmarshal(b, &top); err != nil || len(top) != 21 {
+		t.Fatalf("the degraded request has %d top-level keys (%v), want 21", len(top), err)
+	}
+	if top["reuse"] != false || top["degraded"] != true || top["branch"] != "" {
+		t.Fatalf("the degraded request: %s", b)
+	}
+	have := map[string]bool{}
+	leafPaths(top, "", have)
+	for _, f := range pageFields(t) {
+		if !have[f.Path] {
+			t.Errorf("the degraded request lacks %s", f.Path)
+		}
+		delete(have, f.Path)
+	}
+	if len(have) != 0 || paths["branch"] != "string" {
+		t.Errorf("leaf paths with no FIELDS entry %v; branch %q", have, paths["branch"])
+	}
+}
````

### Patch M7a — Needs you lists the pending requests

````diff M7a
--- a/host/surfacepage/app/app.js
+++ b/host/surfacepage/app/app.js
@@ -421,7 +421,18 @@
     return '<div class="nl sub" data-service' + (service.warn ? ' style="color:var(--you)"' : '') + '>' + esc(service.line) + '</div>';
   }
 
-  if (typeof module !== 'undefined') module.exports = { isSharedKey, sourcesHTML, openRunningHTML, fetchCoherent, serviceLineHTML, splitFrames, applyDelta, nextBackoff, liveRequest, liveStatus, esc, cleanTitle, shortTitle, prettyModel, agentLine, summaryOf, rowLabel, verdict, LEVELS, levelXs, blendXS,
+  // needsYouHTML is the now card's Needs you (row 201 A3, design §2.8): the pending mint requests of
+  // GET /surface/requests.json, each a link to its card; a request the daemon could not read shows why.
+  function needsYouHTML(r, now) {
+    const p = r && Array.isArray(r.pending) ? r.pending : [];
+    const head = `<div class="nh" style="color:${p.length ? 'var(--youShell)' : 'var(--sub)'}">Needs you · ${p.length}</div>`;
+    if (!p.length) return `<div class="nsec" data-needs>${head}<div class="nl sub">${esc(r && r.unavailable || 'Nothing needs you')}</div></div>`;
+    return `<div class="nsec" data-needs>${head}` + p.map(q => q.refused ? `<div class="nl sub">${esc(q.id)} · refused: ${esc(q.refused)}</div>`
+      : `<div class="nl"><a id="needs-${esc(q.id)}" href="/surface/approve/${esc(q.id)}"><b>${esc(q.episode)}</b> · ${esc(q.agent)} · ${q.grants} grant${q.grants === 1 ? '' : 's'} · expires in ${Math.max(0, Math.ceil((q.expires - now) / 60))} min</a>`
+        + (q.recorded ? ` · approved by ${esc(q.recorded)}, waiting for the minter` : '') + '</div>').join('') + '</div>';
+  }
+
+  if (typeof module !== 'undefined') module.exports = { needsYouHTML, isSharedKey, sourcesHTML, openRunningHTML, fetchCoherent, serviceLineHTML, splitFrames, applyDelta, nextBackoff, liveRequest, liveStatus, esc, cleanTitle, shortTitle, prettyModel, agentLine, summaryOf, rowLabel, verdict, LEVELS, levelXs, blendXS,
     adaptModel, layoutGraph, stackBars, barHTML, frontHTML, railHTML, appPath, stackPath, parseAppQuery, addressTarget, lookOf, crumbsHTML, geometry: { PL, PR, PW } };
   if (typeof document === 'undefined') return;
 
@@ -580,7 +591,7 @@
   }
 
   function nowCard() {
-    const ne = D.nowEdge || {}, now = NOW_D, sess = ne.sessions || [], needs = ne.needsYou || [], td = ne.today;
+    const ne = D.nowEdge || {}, now = NOW_D, sess = ne.sessions || [], td = ne.today;
     const healthy = /healthy/i.test(ne.health);
     let today = '';
     if (td) {
@@ -602,8 +613,7 @@
       ${sourcesHTML(D.repos)}
       ${today}
       ${running}
-      <div class="nsec"><div class="nh" style="color:${needs.length ? 'var(--youShell)' : 'var(--sub)'}">Needs you · ${needs.length}</div>
-        ${needs.length ? '' : `<div class="nl sub">Nothing needs you</div>`}</div>
+      ${needsYouHTML(REQS, REQS ? REQS.now : 0)}
     </div>`;
   }
 
@@ -879,11 +889,17 @@
     try { body = await r.json(); } catch (_) { body = null; }
     return { ok: r.ok, status: r.status, body };
   }
+  // REQS is GET /surface/requests.json, read with the model and on every refresh (no timer of its own).
+  let REQS = null;
+  async function loadRequests() {
+    try { REQS = (await getJSON('/surface/requests.json')).body; } catch (_) { REQS = null; }
+  }
   const FOLLOW_MODEL = 'world/surface-model/v2', POLL_MS = 15000;
   let followMode = false, shownGen = null, shownAt = '', pollTimer = null, liveStarted = false;
   const poll = { note: '' };
   async function loadModel() {
     let res;
+    await loadRequests();
     try { res = await getJSON('/surface/model.json'); } catch (e) { return renderVoid('error', 'the surface model could not be fetched: ' + e.message); }
     const b = res.body || {};
     followMode = b.schema === FOLLOW_MODEL;
@@ -974,6 +990,7 @@
     const count = live.count;
     try {
       const selection = T && st.d === 1 && F ? {brane: F.brane, at: F.id, seq: openSeq} : null;
+      await loadRequests();
       const r = await fetchCoherent(getJSON, selection, T ? shownGen : null);
       if (r.kind === 'fail') {
         if (!followMode) throw new Error('refresh failed (HTTP ' + r.status + ')');
--- a/host/surfacepage/app_test.go
+++ b/host/surfacepage/app_test.go
@@ -172,9 +172,15 @@ const model = JSON.parse(fs.readFileSync('model.json', 'utf8'));
 const stacks = JSON.parse(fs.readFileSync('stacks.json', 'utf8'));
 const index = JSON.parse(fs.readFileSync('index.json', 'utf8'));
 assert.deepEqual(Object.keys(app).sort(), ['LEVELS', 'adaptModel', 'addressTarget', 'agentLine', 'appPath', 'applyDelta', 'barHTML', 'blendXS', 'cleanTitle', 'crumbsHTML', 'esc', 'fetchCoherent', 'frontHTML',
-  'geometry', 'isSharedKey', 'layoutGraph', 'levelXs', 'liveRequest', 'liveStatus', 'lookOf', 'nextBackoff', 'openRunningHTML', 'parseAppQuery', 'prettyModel', 'railHTML', 'rowLabel', 'serviceLineHTML', 'shortTitle', 'sourcesHTML', 'splitFrames', 'stackBars', 'stackPath', 'summaryOf', 'verdict']);
+  'geometry', 'isSharedKey', 'layoutGraph', 'levelXs', 'liveRequest', 'liveStatus', 'lookOf', 'needsYouHTML', 'nextBackoff', 'openRunningHTML', 'parseAppQuery', 'prettyModel', 'railHTML', 'rowLabel', 'serviceLineHTML', 'shortTitle', 'sourcesHTML', 'splitFrames', 'stackBars', 'stackPath', 'summaryOf', 'verdict']);
 
 assert.equal(app.serviceLineHTML(undefined), '');
+// Row 201 A3: Needs you is the pending requests, each a link to its card; a refused row shows why, with no link.
+assert.equal(app.needsYouHTML(null, 0), '<div class="nsec" data-needs><div class="nh" style="color:var(--sub)">Needs you · 0</div><div class="nl sub">Nothing needs you</div></div>');
+{ const h = app.needsYouHTML({now: 1000, pending: [{id: 'mr-0123456789abcdef', episode: 'w<1>', agent: 'claude', grants: 2, expires: 1601, recorded: 'webauthn:mark'}, {id: 'mr-00000000000000ff', episode: 'x', refused: 'digest <bad>', recorded: null}]}, 1000);
+  assert(h.includes('Needs you · 2') && h.includes('<a id="needs-mr-0123456789abcdef" href="/surface/approve/mr-0123456789abcdef"><b>w&lt;1&gt;</b> · claude · 2 grants · expires in 11 min</a> · approved by webauthn:mark, waiting for the minter</div>'), h);
+  assert(h.includes('mr-00000000000000ff · refused: digest &lt;bad&gt;') && !h.includes('/surface/approve/mr-00000000000000ff'), h); }
+assert(app.needsYouHTML({pending: [], unavailable: 'the mint requests could not be read'}, 0).includes('>the mint requests could not be read<'));
 assert.equal(app.serviceLineHTML({line:'<build & file>',warn:false}), '<div class="nl sub" data-service>&lt;build &amp; file&gt;</div>');
 assert.equal(app.serviceLineHTML({line:'restart',warn:true}), '<div class="nl sub" data-service style="color:var(--you)">restart</div>');
 
````

### Patch M8a — the driver, the card's helpers, XA1

````diff M8a
--- a/cmd/ailang-worldd/approve_page_chrome_test.go
+++ b/cmd/ailang-worldd/approve_page_chrome_test.go
@@ -0,0 +1,109 @@
+package main
+
+// Row 201 A3 sprint 3 M8a-M8b (design §2.7, §6 XA1, XA2): the request card in a live Chrome.
+
+import (
+	"bytes"
+	"encoding/hex"
+	"os"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+	"time"
+)
+
+// XA1: mark-mac enrolled through the A2 page and the CLI (as XB1); the card at
+// /surface/approve/<id> shows every field of the stored request as its text; Approve records the
+// passkey's assertion (200) and approve --await mints it; no CSP violation and no page error; the
+// card's get has the same challenge, rpId, credential and user verification as the A2 page's get for
+// the same request (the CLI's --webauthn-link). Opened again after the mint, the card says it is used.
+func TestApprovePageInChrome(t *testing.T) {
+	b := startPasskeyBrowser(t)
+	b.do(t, nil, "tab", map[string]any{"authenticator": "uv", "inject": "getRecorder"})
+	now := time.Now().Unix() // the daemon's routes read the wall clock
+	a := newApproveRig(t, "ep1", now)
+	origin := cardOrigin(t, a.credPath, "", nil)
+	t.Logf("%s, card at %s", b.chrome, origin)
+
+	dl := t.TempDir()
+	var page pageCeremony
+	var out, errw bytes.Buffer
+	term := &ceremonyTerm{}
+	term.answer = func(n int) (string, error) {
+		if n > 1 {
+			return "", os.ErrClosed
+		}
+		if n == 1 {
+			return "y\n", nil
+		}
+		link, err := printedLink(out.String())
+		if err == nil {
+			err = b.call(nil, "open", map[string]any{"url": link})
+		}
+		if err == nil {
+			err = b.call(&page, "ceremony", map[string]any{"download": dl})
+		}
+		if err != nil {
+			t.Error(err)
+			return "", err
+		}
+		return page.After.Code + "\n", nil
+	}
+	env := enrolEnv(term, nil)
+	env.now = atTime(now - 100)
+	args := []string{"enroll", "--name", "mark-mac", "--format", "webauthn", "--origin", origin, "--registration", filepath.Join(dl, "world-registration-mark-mac.json"), "--db", a.r.f.db}
+	if code := runApprover(args, &out, &errw, env); code != exitOK {
+		t.Fatalf("enrol: rc %d\nstdout %s\nstderr %s", code, out.String(), errw.String())
+	}
+
+	var lout, lerr bytes.Buffer
+	if code := runApprove(a.args("--webauthn-link"), &lout, &lerr, approveEnv(nil, nil, atTime(now))); code != exitOK {
+		t.Fatalf("--webauthn-link: rc %d\nstderr %s", code, lerr.String())
+	}
+	link, err := printedLink(lout.String())
+	if err != nil {
+		t.Fatal(err)
+	}
+	b.do(t, nil, "open", map[string]any{"url": link})
+	b.do(t, &page, "ceremony", nil)
+	var a2 []cardGet
+	b.do(t, &a2, "waitFor", map[string]any{"expression": "window.__gets"})
+
+	var view cardView
+	b.do(t, &view, "openCard", map[string]any{"url": origin + "/surface/approve/" + a.id})
+	h := requestHashOf(t, a)
+	if view.State != "ready" || view.ButtonHidden || view.ButtonDisabled || !strings.HasPrefix(view.What, "Approve request "+a.id+" for episode ep1, expires ") ||
+		!strings.HasSuffix(view.What, ", message hash "+hex.EncodeToString(h[:8])) || view.SaveHidden || view.SaveName != "world-approval-"+a.id+".msg" {
+		t.Fatalf("the card shows %+v", view)
+	}
+	if want := cardFields(t, a.row(t).Request); len(want) < 28 || !reflect.DeepEqual(view.Fields, want) || view.TTLLabel != "Session lifetime (s) (1.0 h)" {
+		t.Fatalf("the card's fields\n%v\nwant every leaf of the stored request\n%v", view.Fields, want)
+	}
+	var after cardView // a fresh value: decoding into view would keep its fields
+	b.do(t, &after, "press", map[string]any{"id": "go", "states": []string{"recorded", "rejected", "failed", "unknown"}})
+	if after.State != "recorded" || after.StatusError || !strings.HasPrefix(after.Status, "Approved by webauthn:mark-mac: the minter (approve "+a.id+" --await) mints it") {
+		t.Fatalf("after Approve the card shows %+v", after)
+	}
+	if rec, ok := surfaceRecorded(t, a.credPath, a.id); !ok || rec.Approver != "webauthn:mark-mac" || rec.Format != "webauthn" {
+		t.Fatalf("recorded %+v %v", rec, ok)
+	}
+	if len(after.Violations) != 0 || len(after.Errors) != 0 {
+		t.Fatalf("CSP violations %q, page errors %q", after.Violations, after.Errors)
+	}
+	if len(a2) != 1 || len(after.Gets) != 1 || !reflect.DeepEqual(after.Gets[0], a2[0]) || after.Gets[0].RPID != "localhost" ||
+		after.Gets[0].UserVerification != "required" || after.Gets[0].Timeout != 240000 || !after.Gets[0].Signal || len(after.Gets[0].Allow) != 1 {
+		t.Fatalf("the card's get %+v is not the A2 page's %+v", after.Gets, a2)
+	}
+	var opened bool
+	aenv, _ := awaitEnv(time.Now().Unix(), &opened, nil)
+	var stdout, stderr bytes.Buffer
+	if code := runApprove(a.args("--await", "--await-for", "5s"), &stdout, &stderr, aenv); code != exitOK || opened {
+		t.Fatalf("approve --await: rc %d, terminal opened %v\nstderr %s", code, opened, stderr.String())
+	}
+	requireMintedBy(t, a, "webauthn:mark-mac", stdout.String())
+	var used cardView
+	if b.do(t, &used, "openCard", map[string]any{"url": origin + "/surface/approve/" + a.id}); used.State != "used" || used.Status != "Already used: approved by webauthn:mark-mac." || !used.ButtonHidden {
+		t.Fatalf("the used card shows %+v", used)
+	}
+}
--- a/cmd/ailang-worldd/approve_page_driver_test.go
+++ b/cmd/ailang-worldd/approve_page_driver_test.go
@@ -0,0 +1,177 @@
+package main
+
+// Row 201 A3 sprint 3 M8a (design §2.7, §6 XA1-XA5): the request card's browser harness, on A2's
+// driver (testdata/passkey_cdp.mjs, extended): Chrome headless, a temporary profile and a CDP virtual
+// authenticator holding a test key (D-WORLD-109). The card is served by the production daemon
+// handler on the loopback, over the request's temporary credentials store.
+
+import (
+	"bytes"
+	"context"
+	"crypto/ecdsa"
+	"crypto/elliptic"
+	"crypto/x509"
+	"encoding/base64"
+	"encoding/json"
+	"fmt"
+	"io"
+	"net"
+	"net/http"
+	"net/http/httptest"
+	"path/filepath"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// cardView is what the card page shows (the driver's card command).
+type cardView struct {
+	State, What, Status, Again, SaveHref, SaveName, SSHLine, Note      string
+	StatusError, ButtonHidden, ButtonDisabled, AgainHidden, SaveHidden bool
+	Fields                                                             map[string]string
+	TTLLabel                                                           string
+	Gets                                                               []cardGet
+	Violations, Errors                                                 []string
+	Op                                                                 *struct {
+		Signal, Aborted             bool
+		AbortedAfter, ResolvedAfter *float64
+	}
+}
+
+// cardGet is one navigator.credentials.get call the getRecorder injection saw.
+type cardGet struct {
+	Challenge, RPID, UserVerification string
+	Allow                             []string
+	Timeout                           int
+	Signal                            bool
+}
+
+// cardOrigin serves the production daemon handler (wrapped by wrap, if any) over credentials store
+// creds, with surface index index (or none), on the loopback, and returns its origin.
+func cardOrigin(t *testing.T, creds, index string, wrap func(http.Handler) http.Handler) string {
+	t.Helper()
+	d, err := daemon.New(boundedTestContext(t), daemon.Config{DBPath: filepath.Join(t.TempDir(), "world.db"), CredentialsPath: creds,
+		SurfaceIndex: index, BindHost: daemon.DefaultBindHost, ErrorLog: io.Discard})
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { _ = d.Close() })
+	h := d.Handler()
+	if wrap != nil {
+		h = wrap(h)
+	}
+	srv := httptest.NewServer(h)
+	t.Cleanup(srv.Close)
+	_, port, _ := net.SplitHostPort(srv.Listener.Addr().String())
+	return "http://localhost:" + port
+}
+
+// virtualPasskey enrols approvaltest's ES256 passkey seed as name at origin (its rpID: the origin's
+// host) in creds (at `at`) and gives its private key to the tab's virtual authenticator, so the
+// browser signs as that test key.
+func virtualPasskey(t *testing.T, b *passkeyBrowser, creds, origin, name string, seed byte, at int64) {
+	t.Helper()
+	pk := approvaltest.NewPasskey(approval.AlgES256, seed)
+	k, err := ecdsa.ParseRawPrivateKey(elliptic.P256(), bytes.Repeat([]byte{seed}, 32))
+	if err != nil {
+		t.Fatal(err)
+	}
+	if pub, _ := k.PublicKey.Bytes(); !bytes.Equal(pub[1:33], pk.Key().X) { // control: the same key
+		t.Fatal("the virtual credential is not approvaltest's key")
+	}
+	der, err := x509.MarshalPKCS8PrivateKey(k)
+	if err != nil {
+		t.Fatal(err)
+	}
+	rp, err := approval.ParseCeremonyOrigin(origin)
+	if err != nil {
+		t.Fatal(err)
+	}
+	st, err := store.OpenCredentials(creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = st.Close() }()
+	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
+	defer cancel()
+	if err := st.EnrolPasskey(ctx, store.ApproverRow{Name: name, Format: approval.FormatWebAuthn, PublicKey: pk.Key().Canonical(), EnrolledAt: at},
+		store.PasskeyRow{Name: name, RPID: rp, Origin: origin, CredentialID: pk.CredentialID(), AAGUID: make([]byte, 16), Alg: approval.AlgES256, SignCount: 1}); err != nil {
+		t.Fatal(err)
+	}
+	b.do(t, nil, "credential", map[string]any{"id": base64.StdEncoding.EncodeToString(pk.CredentialID()), "rpId": rp,
+		"privateKey": base64.StdEncoding.EncodeToString(der)})
+}
+
+// cardFields is what the card must show for canonical request bytes req: every leaf, by its
+// concrete key path, as text (a boolean as FIELDS' word for it).
+func cardFields(t *testing.T, req []byte) map[string]string {
+	t.Helper()
+	_, rest, _ := bytes.Cut(surfacepage.ApproveScript, []byte("/* FIELDS */"))
+	lit, _, _ := bytes.Cut(rest, []byte("/* /FIELDS */"))
+	var fields []struct{ Path, Type, No, Yes string }
+	if err := json.Unmarshal(lit, &fields); err != nil || len(fields) == 0 {
+		t.Fatalf("approve.js's FIELDS: %v", err)
+	}
+	words := map[string][2]string{}
+	for _, f := range fields {
+		words[f.Path] = [2]string{f.No, f.Yes}
+	}
+	dec := json.NewDecoder(bytes.NewReader(req))
+	dec.UseNumber()
+	var v any
+	if err := dec.Decode(&v); err != nil {
+		t.Fatal(err)
+	}
+	out := map[string]string{}
+	var walk func(v any, pattern, concrete string)
+	walk = func(v any, pattern, concrete string) {
+		dot := func(p string) string {
+			if p == "" {
+				return p
+			}
+			return p + "."
+		}
+		switch x := v.(type) {
+		case map[string]any:
+			for k, e := range x {
+				walk(e, dot(pattern)+k, dot(concrete)+k)
+			}
+		case []any:
+			for i, e := range x {
+				walk(e, pattern+"[]", fmt.Sprintf("%s.%d", concrete, i))
+			}
+		case bool:
+			w := words[pattern][0]
+			if x {
+				w = words[pattern][1]
+			}
+			out[concrete] = w
+		default:
+			out[concrete] = fmt.Sprint(x)
+		}
+	}
+	walk(v, "", "")
+	return out
+}
+
+// surfaceRecorded is the surface approval recorded for request id in creds, if any.
+func surfaceRecorded(t *testing.T, creds, id string) (store.SurfaceApprovalRow, bool) {
+	t.Helper()
+	ro, err := store.OpenCredentialsReadOnly(creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = ro.Close() }()
+	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
+	defer cancel()
+	row, ok, err := ro.SurfaceApproval(ctx, id)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return row, ok
+}
--- a/cmd/ailang-worldd/testdata/passkey_cdp.mjs
+++ b/cmd/ailang-worldd/testdata/passkey_cdp.mjs
@@ -1,9 +1,10 @@
 // The passkey page's browser driver (row 201 A2 M7; XB1, XB2 and host/approval/testdata/webauthn/
-// generate.sh). Test-only: Chrome headless with a temporary profile under argv[2], driven over the
+// generate.sh), extended for the request card (row 201 A3 sprint 3; XA1-XA3, XA5). Test-only: Chrome headless with a temporary profile under argv[2], driven over the
 // DevTools protocol, with a CDP VIRTUAL authenticator, a test key (D-WORLD-109): never a real
 // profile, keychain, passkey or platform authenticator. Commands arrive as JSON lines on stdin;
 // each answer is one JSON line on stdout. EOF on stdin kills Chrome and removes the profile.
 import { spawn } from 'node:child_process';
+import { createHash } from 'node:crypto';
 import fs from 'node:fs';
 import path from 'node:path';
 import readline from 'node:readline';
@@ -31,6 +32,7 @@ const replies = new Map(), events = [];
 ws.onmessage = (m) => {
   const msg = JSON.parse(m.data);
   if (msg.id && replies.has(msg.id)) { replies.get(msg.id)(msg); replies.delete(msg.id); } else if (msg.method) events.push(msg);
+  if (msg.method === 'Fetch.requestPaused') rewritePaused(msg.params).catch((e) => { rewriteErrors.push(String(e && e.message || e)); });
 };
 const send = async (method, params = {}, sessionId = session) => {
   const msg = await new Promise((r) => { const id = ++seq; replies.set(id, r); ws.send(JSON.stringify({ id, method, params, sessionId: sessionId || undefined })); });
@@ -43,7 +45,7 @@ const evaluate = async (expression) => {
   return JSON.parse(r.result.value);
 };
 const click = async (id) => {
-  const box = await evaluate(`(() => { const e = document.getElementById('${id}'); e.scrollIntoView({ block: 'center' }); const r = e.getBoundingClientRect(); return { x: r.x + r.width / 2, y: r.y + r.height / 2, w: r.width }; })()`);
+  const box = await evaluate(`(() => { const e = document.getElementById('${id}'); e.scrollIntoView({ block: 'center' }); const r = e.getClientRects()[0] || e.getBoundingClientRect(); return { x: r.x + r.width / 2, y: r.y + r.height / 2, w: r.width }; })()`);
   if (box.w === 0) return false;
   for (const type of ['mousePressed', 'mouseReleased']) await send('Input.dispatchMouseEvent', { type, x: box.x, y: box.y, button: 'left', clickCount: 1 });
   return true;
@@ -63,6 +65,11 @@ const injections = {
     setTimeout(() => { op.resolvedAfter = performance.now() - op.t0; op.resolvedAfterDate = Date.now() - op.d0; resolve({ id: 'AA', rawId: b(1), type: 'public-key',
       response: { clientDataJSON: b(2), attestationObject: b(3), authenticatorData: b(37), signature: b(8), userHandle: null } }); }, 245000);
   });`,
+  // getRecorder (XA1, XA2): every navigator.credentials.get call's options, then the real get.
+  getRecorder: `window.__gets = []; if (navigator.credentials) { const get = navigator.credentials.get.bind(navigator.credentials);
+    const b = (x) => btoa(String.fromCharCode(...new Uint8Array(x.buffer || x))).replace(/\\+/g, '-').replace(/\\//g, '_').replace(/=+$/, '');
+    navigator.credentials.get = (o) => { const k = o.publicKey; window.__gets.push({ challenge: b(k.challenge), rpId: k.rpId, allow: (k.allowCredentials || []).map((c) => b(c.id)),
+      userVerification: k.userVerification, timeout: k.timeout, signal: !!o.signal }); return get(o); }; }`,
   eddsa: `{ const create = navigator.credentials.create.bind(navigator.credentials);
     navigator.credentials.create = (o) => { o.publicKey.pubKeyCredParams = o.publicKey.pubKeyCredParams.filter((p) => p.alg === -8); return create(o); }; }`,
 };
@@ -154,9 +161,117 @@ async function advance({ budget }) {
   return pageState();
 }
 
+// fetchRewrite (XA2) rewrites the card's GET /surface/requests/<id> answer in the browser: 'bytes'
+// changes the stored bytes and keeps the digest; 'forge' changes both; 'extra' adds a key, 'drop:<key>'
+// deletes one, 'type' makes reuse a string, 'float' makes ttl_s 3600.5 and 'utf8' puts a 0xff byte in the
+// model, each re-hashed so the digest matches; 'id' names another request and 'nonce' sends a
+// non-hex challenge, 'expiry' a negative expiry and 'b64' bytes that are not base64url. Every paused response is counted.
+let rewriteMode = null, rewritten = 0;
+const rewriteErrors = [];
+async function rewritePaused(p) {
+  const { body, base64Encoded } = await send('Fetch.getResponseBody', { requestId: p.requestId });
+  const card = JSON.parse(base64Encoded ? Buffer.from(body, 'base64').toString('utf8') : body);
+  const req = JSON.parse(Buffer.from(card.request, 'base64url').toString('utf8'));
+  if (rewriteMode === 'extra') req.extra = 'x';
+  else if (rewriteMode.startsWith('drop:')) delete req[rewriteMode.slice(5)];
+  else if (rewriteMode === 'type') req.reuse = String(req.reuse);
+  else if (rewriteMode === 'float') req.ttl_s += 0.5;
+  else if (rewriteMode === 'bytes' || rewriteMode === 'forge') req.ttl_s += 3600;
+  let bytes = Buffer.from(JSON.stringify(req), 'utf8');
+  if (rewriteMode === 'utf8') { const k = bytes.indexOf('"model":"') + 9; bytes = Buffer.concat([bytes.subarray(0, k), Buffer.from([0xff]), bytes.subarray(k)]); }
+  if (rewriteMode === 'id') card.id = 'mr-0000000000000000';
+  else if (rewriteMode === 'nonce') card.challenge = 'zz' + card.challenge.slice(2);
+  else if (rewriteMode === 'expiry') card.expires = -1;
+  else if (rewriteMode === 'b64') card.request = '*' + card.request;
+  else card.request = bytes.toString('base64url');
+  if (!['bytes', 'id', 'nonce', 'expiry', 'b64'].includes(rewriteMode)) card.digest = createHash('sha256').update(bytes).digest('hex');
+  rewritten++;
+  await send('Fetch.fulfillRequest', { requestId: p.requestId, responseCode: p.responseStatusCode, responseHeaders: p.responseHeaders,
+    body: Buffer.from(JSON.stringify(card)).toString('base64') });
+}
+async function fetchRewrite({ mode }) {
+  rewriteMode = mode;
+  if (mode) await send('Fetch.enable', { patterns: [{ urlPattern: '*/surface/requests/mr-*', requestStage: 'Response' }] });
+  else await send('Fetch.disable');
+  return { rewritten, errors: rewriteErrors };
+}
+
+// The card page's view: its state and texts, every rendered field (f-<path> → text), the buttons,
+// the message link, the recorded get calls, CSP violations and the page's errors since the tab.
+const pageErrors = () => events.filter((e) => e.method === 'Runtime.exceptionThrown' || (e.method === 'Log.entryAdded' && e.params.entry.level === 'error'))
+  .map((e) => e.params.entry ? e.params.entry.text : e.params.exceptionDetails.text);
+const cardState = async () => Object.assign(await evaluate(`(() => { const $ = (id) => document.getElementById(id); return {
+  state: document.body.dataset.state, what: $('what').textContent, status: $('status').textContent, statusError: $('status').hasAttribute('data-error'),
+  buttonHidden: $('go').hidden, buttonDisabled: $('go').disabled, againHidden: $('again').hidden, again: $('again').textContent,
+  saveHidden: $('save').hidden, saveHref: $('save').getAttribute('href') || '', saveName: $('save').download, sshLine: $('sshline').textContent,
+  fields: Object.fromEntries(Array.from(document.querySelectorAll('[id^="f-"]'), (e) => [e.id.slice(2), e.textContent])),
+  note: $('origin-note').hidden ? '' : $('origin-note').textContent, ttlLabel: $('f-ttl_s') ? $('f-ttl_s').previousElementSibling.textContent : '', op: window.__op || null,
+  gets: window.__gets || null, violations: window.__cspv || [] }; })()`), { errors: pageErrors() });
+
+// openCard opens a card page and waits until it leaves 'loading' (at most settleMs); raw opens any
+// page and returns its location.
+async function openCard({ url, settleMs = 15000, raw = false }) {
+  for (const u of ['about:blank', url]) { // a new document each time
+    const from = events.length;
+    await send('Page.navigate', { url: u });
+    await until(() => events.slice(from).some((e) => e.method === 'Page.loadEventFired'), 15000, 'the page load');
+  }
+  if (raw) return { location: await evaluate('location.href') };
+  await until(async () => (await evaluate(`document.body && document.body.dataset.state`)) !== 'loading', settleMs, 'the card script');
+  return cardState();
+}
+
+// press clicks element id; with page, it waits for the new document; else, until the card's state
+// is one of states (at most settleMs), or with virtualBudget, after that much page time from a clock
+// paused before the click. It returns the card's view, or with page only the location.
+async function press({ id, page = false, states = null, settleMs = 30000, virtualBudget = 0 }) {
+  const from = events.length;
+  if (virtualBudget > 0) await send('Emulation.setVirtualTimePolicy', { policy: 'pause' });
+  if (!(await click(id))) throw new Error(id + ' is not visible');
+  if (virtualBudget > 0) {
+    await send('Emulation.setVirtualTimePolicy', { policy: 'advance', budget: virtualBudget });
+    await until(() => events.slice(from).some((e) => e.method === 'Emulation.virtualTimeBudgetExpired'), settleMs, 'the virtual time budget');
+  }
+  if (page) {
+    await until(() => events.slice(from).some((e) => e.method === 'Page.loadEventFired'), 15000, 'the next page');
+    return { location: await evaluate('location.href') };
+  }
+  if (states) await until(async () => states.includes(await evaluate(`document.body.dataset.state`)), settleMs, 'a state in ' + states);
+  return cardState();
+}
+
+// credential gives the virtual authenticator a test key's credential (base64 id and PKCS#8 key).
+async function credential({ id, rpId, privateKey }) {
+  await send('WebAuthn.addCredential', { authenticatorId, credential: { credentialId: id, isResidentCredential: false, rpId, privateKey, signCount: 1 } });
+  return {};
+}
+
+// fill types value into element id (as a paste: one input event).
+async function fill({ id, value }) {
+  await evaluate(`(() => { const e = document.getElementById(${JSON.stringify(id)}); e.value = ${JSON.stringify(value)}; e.dispatchEvent(new Event('input')); return 1; })()`);
+  return cardState();
+}
+
+// save clicks element id and waits for its download into dir; it returns the download's events.
+async function save({ id, dir }) {
+  const from = events.length;
+  await send('Browser.setDownloadBehavior', { behavior: 'allow', downloadPath: dir, eventsEnabled: true }, null);
+  if (!(await click(id))) throw new Error(id + ' is not visible');
+  await until(() => events.slice(from).some((e) => e.method === 'Browser.downloadProgress' && e.params.state !== 'inProgress'), 15000, 'the download');
+  return events.slice(from).filter((e) => e.method.startsWith('Browser.download'))
+    .map((e) => ({ event: e.method, state: e.params.state || '', name: e.params.suggestedFilename || '', url: e.params.url || '' }));
+}
+
+// waitFor evaluates expression until it is truthy (at most ms) and returns its value.
+async function waitFor({ expression, ms = 15000 }) {
+  let v;
+  await until(async () => (v = await evaluate(expression)), ms, expression);
+  return v;
+}
+
 const product = (await send('Browser.getVersion', {}, null)).product;
 console.log(JSON.stringify({ ready: true, chrome: product, pid: chrome.pid }));
-const commands = { tab, open, ceremony, uv, advance };
+const commands = { tab, open, ceremony, uv, advance, fetchRewrite, openCard, press, waitFor, credential, fill, save, card: cardState };
 for await (const line of readline.createInterface({ input: process.stdin })) {
   const { cmd, ...args } = JSON.parse(line);
   try {
````

### Patch M8b — the card refuses what it cannot bind; its SSH path

````diff M8b
--- a/cmd/ailang-worldd/approve_page_swap_test.go
+++ b/cmd/ailang-worldd/approve_page_swap_test.go
@@ -0,0 +1,119 @@
+package main
+
+// Row 201 A3 sprint 3 M8b (design §2.7, §4 T4, §6 XA2; AW6's page half): the card refuses what it
+// cannot bind, and its SSH path, in a live Chrome.
+
+import (
+	"bytes"
+	"os"
+	"path/filepath"
+	"strings"
+	"testing"
+	"time"
+)
+
+// XA2: the card's response rewritten in the browser. (a) bytes changed, digest kept (or another
+// request's id): refused, no get; (b) bytes and digest forged together: the page signs, the daemon
+// refuses the POST (409, R9: the body's digest is not the request's), nothing is recorded; (c) a key
+// added, branch deleted, reuse a string, ttl_s a fraction, or a byte that is not UTF-8, each re-hashed:
+// refused naming the path, no get and no message link; a non-hex nonce, a negative expiry or bytes that
+// are not base64url: refused.
+// Control: no rewrite, ready.
+func TestApprovePageRefusesASwappedRequest(t *testing.T) {
+	b := startPasskeyBrowser(t)
+	b.do(t, nil, "tab", map[string]any{"authenticator": "uv", "inject": "getRecorder"})
+	now := time.Now().Unix()
+	a := newApproveRig(t, "ep1", now)
+	origin := cardOrigin(t, a.credPath, "", nil)
+	virtualPasskey(t, b, a.credPath, origin, "mark-mac", 81, now-100)
+	card := origin + "/surface/approve/" + a.id
+	for _, tc := range []struct{ mode, state, what string }{
+		{"", "ready", "Approve request " + a.id},
+		{"bytes", "refused", "The daemon's request does not match its digest; nothing to approve."},
+		{"forge", "ready", "Approve request " + a.id},
+		{"extra", "refused", "This request has a field the page cannot show (extra); nothing to approve; approve it at a terminal."},
+		{"drop:branch", "refused", "This request has a field the page cannot show (branch); nothing to approve; approve it at a terminal."},
+		{"type", "refused", "This request has a field the page cannot show (reuse); nothing to approve; approve it at a terminal."},
+		{"float", "refused", "This request has a field the page cannot show (ttl_s); nothing to approve; approve it at a terminal."},
+		{"utf8", "refused", "The daemon's request is not UTF-8 JSON ("},
+		{"id", "refused", "The daemon's request does not match its digest; nothing to approve."},
+		{"nonce", "refused", "The daemon's challenge or expiry is malformed; nothing to approve."},
+		{"expiry", "refused", "The daemon's challenge or expiry is malformed; nothing to approve."},
+		{"b64", "refused", "The daemon's request is not base64url; nothing to approve."},
+	} {
+		var view cardView // a fresh value: decoding into an old one would keep its fields
+		b.do(t, nil, "fetchRewrite", map[string]any{"mode": tc.mode})
+		b.do(t, &view, "openCard", map[string]any{"url": card})
+		if view.State != tc.state || !strings.HasPrefix(view.What, tc.what) {
+			t.Fatalf("%q: state %q, what %q; want %q, %q", tc.mode, view.State, view.What, tc.state, tc.what)
+		}
+		if tc.state == "refused" && (len(view.Gets) != 0 || !view.ButtonHidden || !view.SaveHidden || view.SaveHref != "" || len(view.Fields) != 0) {
+			t.Fatalf("%q: a refused card offers something: %+v", tc.mode, view)
+		}
+		if tc.mode != "forge" {
+			continue
+		}
+		b.do(t, &view, "press", map[string]any{"id": "go", "states": []string{"recorded", "rejected", "failed", "unknown"}})
+		if view.State != "rejected" || len(view.Gets) != 1 || !strings.Contains(view.Status, "HTTP 409 Conflict): the request changed since the card was shown") {
+			t.Fatalf("forged: the page shows %+v", view)
+		}
+		if _, ok := surfaceRecorded(t, a.credPath, a.id); ok {
+			t.Fatal("forged: an approval was recorded")
+		}
+	}
+	var n struct{ Rewritten int }
+	b.do(t, &n, "fetchRewrite", map[string]any{"mode": ""})
+	if n.Rewritten != 11 {
+		t.Fatalf("%d card responses rewritten, want 11", n.Rewritten)
+	}
+}
+
+// The card's SSH path (planner, the page half of AW6): Save the message saves exactly the request's
+// 90-byte message (approval.SignedMessage) through a blob: download of the card's origin; its
+// approvaltest SSHSIG, pasted, is recorded as ssh:mark-mac (200), and approve --await mints it. A
+// signature over another message is refused (422) and records nothing. The only passkey is enrolled
+// at another origin: the card names it and hides the passkey button.
+func TestApprovePageSSHSIGInChrome(t *testing.T) {
+	b := startPasskeyBrowser(t)
+	b.do(t, nil, "tab", map[string]any{"authenticator": "uv"})
+	now := time.Now().Unix()
+	a := newApproveRig(t, "ep1", now)
+	seedApprover(t, a.credPath, "mark-mac", sigSigner.PublicKey(), now-100)
+	origin := cardOrigin(t, a.credPath, "", nil)
+	virtualPasskey(t, b, a.credPath, "https://studio.tail0.ts.net", "mark-phone", 83, now-100)
+	var view cardView
+	b.do(t, &view, "openCard", map[string]any{"url": origin + "/surface/approve/" + a.id})
+	if view.State != "ready" || !view.ButtonHidden || view.Note != "No passkey is enrolled for "+origin+": open this card at https://studio.tail0.ts.net." || view.SaveHidden || view.SSHLine != "ssh-keygen -Y sign -n world-approval@sunholo -f <key> world-approval-"+a.id+".msg" {
+		t.Fatalf("the card shows %+v", view)
+	}
+	dl := t.TempDir()
+	var downloads []struct{ Event, State, Name, URL string }
+	b.do(t, &downloads, "save", map[string]any{"id": "save", "dir": dl})
+	name := "world-approval-" + a.id + ".msg"
+	msg, err := os.ReadFile(filepath.Join(dl, name))
+	if want := requestMessage(t, a); err != nil || !bytes.Equal(msg, want) || len(downloads) == 0 || !strings.HasPrefix(downloads[0].URL, "blob:"+origin+"/") {
+		t.Fatalf("saved %x (%v) via %+v, want the request's message %x", msg, err, downloads, want)
+	}
+	var after cardView
+	b.do(t, nil, "fill", map[string]any{"id": "sig", "value": string(sigSigner.Sign("world-approval@sunholo", append([]byte("x"), msg...)))})
+	b.do(t, &after, "press", map[string]any{"id": "post-sig", "states": []string{"recorded", "rejected", "unknown"}})
+	if after.State != "rejected" || !strings.Contains(after.Status, "(HTTP 422 NotApproved): the SSH signature does not approve request "+a.id) {
+		t.Fatalf("a signature over another message: %+v", after)
+	}
+	if _, ok := surfaceRecorded(t, a.credPath, a.id); ok {
+		t.Fatal("a refused signature was recorded")
+	}
+	b.do(t, nil, "fill", map[string]any{"id": "sig", "value": string(sigSigner.Sign("world-approval@sunholo", msg))})
+	var done cardView
+	b.do(t, &done, "press", map[string]any{"id": "post-sig", "states": []string{"recorded", "rejected", "unknown"}})
+	if done.State != "recorded" || !strings.HasPrefix(done.Status, "Approved by ssh:mark-mac: ") || len(done.Violations) != 0 {
+		t.Fatalf("the signature: %+v", done)
+	}
+	var opened bool
+	aenv, _ := awaitEnv(time.Now().Unix(), &opened, nil)
+	var stdout, stderr bytes.Buffer
+	if code := runApprove(a.args("--await", "--await-for", "5s"), &stdout, &stderr, aenv); code != exitOK || opened {
+		t.Fatalf("approve --await: rc %d, terminal opened %v\nstderr %s", code, opened, stderr.String())
+	}
+	requireMintedBy(t, a, "ssh:mark-mac", stdout.String())
+}
````

### Patch M8c — the card's network deadlines and its passkey bound

````diff M8c
--- a/cmd/ailang-worldd/approve_card_deadline_test.go
+++ b/cmd/ailang-worldd/approve_card_deadline_test.go
@@ -0,0 +1,286 @@
+package main
+
+// Row 201 A3 sprint 3 M8c (design §2.7 [r2], §6 XA5): the card's network deadlines in a live Chrome.
+
+import (
+	"bytes"
+	"io"
+	"net/http"
+	"strings"
+	"sync"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// cardNet stands between the browser and the daemon: in a stall mode the matching request gets no
+// header, or a header and part of its body, and then nothing; post-recorded runs the daemon's POST
+// (it records) and withholds its answer; restart answers a POST 503 Restart, get-500 the card's GET a
+// 500, and list-header stalls GET /surface/requests.json (the status check). It counts POSTs.
+type cardNet struct {
+	mu          sync.Mutex
+	mode        string
+	posts, ssh  int
+	release     chan struct{}
+	recordedRun int
+}
+
+func (n *cardNet) set(mode string) { n.mu.Lock(); n.mode = mode; n.mu.Unlock() }
+func (n *cardNet) count() int      { n.mu.Lock(); defer n.mu.Unlock(); return n.posts }
+
+func (n *cardNet) wrap(h http.Handler) http.Handler {
+	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
+		post := r.Method == "POST" && strings.HasPrefix(r.URL.Path, "/surface/approve/")
+		get := r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/surface/requests/")
+		list := r.Method == "GET" && r.URL.Path == "/surface/requests.json"
+		if post {
+			body, _ := io.ReadAll(r.Body)
+			r.Body = io.NopCloser(bytes.NewReader(body))
+			n.mu.Lock()
+			n.posts++
+			if bytes.Contains(body, []byte(`"sshsig"`)) {
+				n.ssh++
+			}
+			n.mu.Unlock()
+		}
+		n.mu.Lock()
+		mode := n.mode
+		n.mu.Unlock()
+		stall := func() {
+			select {
+			case <-r.Context().Done():
+			case <-n.release:
+			}
+		}
+		switch {
+		case (get && mode == "get-header") || (post && mode == "post-header") || (list && mode == "list-header"):
+			stall()
+		case (get && mode == "get-body") || (post && mode == "post-body"):
+			w.Header().Set("Content-Type", "application/json")
+			w.WriteHeader(http.StatusOK)
+			_, _ = io.WriteString(w, `{"v":1,`)
+			w.(http.Flusher).Flush()
+			stall()
+		case post && mode == "post-recorded":
+			h.ServeHTTP(newDiscardWriter(), r)
+			n.mu.Lock()
+			n.recordedRun++
+			n.mu.Unlock()
+			stall()
+		case get && mode == "get-500":
+			w.Header().Set("Content-Type", "application/json")
+			w.WriteHeader(http.StatusInternalServerError)
+			_, _ = io.WriteString(w, `{"v":1,"error":"Internal","reason":"internal error (test)"}`)
+		case post && mode == "restart":
+			w.Header().Set("Content-Type", "application/json")
+			w.WriteHeader(http.StatusServiceUnavailable)
+			_, _ = io.WriteString(w, `{"v":1,"error":"Restart","reason":"restart serve (test)"}`)
+		default:
+			h.ServeHTTP(w, r)
+		}
+	})
+}
+
+// discardWriter is a ResponseWriter whose answer nobody reads.
+type discardWriter struct{ h http.Header }
+
+func newDiscardWriter() *discardWriter               { return &discardWriter{h: http.Header{}} }
+func (d *discardWriter) Header() http.Header         { return d.h }
+func (d *discardWriter) Write(b []byte) (int, error) { return len(b), nil }
+func (d *discardWriter) WriteHeader(int)             {}
+
+// XA5: every card GET and POST, body included, ends at its 10 s deadline. A stalled GET (no header;
+// a header and part of the body) shows a retryable error, and the retry reads the card. A stalled
+// POST (no header; part of the body; recorded by the daemon but unanswered) shows "outcome
+// unknown", never refused or minted, and its status check tells the truth (not recorded; recorded;
+// minted meanwhile: no longer pending), itself under the deadline (a stalled check offers another).
+// A card GET answered 500 shows the daemon's reason and the same retry.
+// Only an explicit 503 Busy is retried: 1 s apart, at most 3 times (a writer held 1.5 s: 3 POSTs,
+// recorded; held throughout: 4 POSTs, then refused); a 503 Restart is not retried. No timeout or
+// refusal ever sends the other format.
+func TestApproveCardNetworkDeadlines(t *testing.T) {
+	setup := func(t *testing.T) (*passkeyBrowser, approveRig, *cardNet, string) {
+		b := startPasskeyBrowser(t)
+		b.do(t, nil, "tab", map[string]any{"authenticator": "uv"})
+		now := time.Now().Unix()
+		a := newApproveRig(t, "ep1", now)
+		n := &cardNet{release: make(chan struct{})}
+		origin := cardOrigin(t, a.credPath, "", n.wrap)
+		t.Cleanup(func() { close(n.release) }) // before the server's Close, which waits for its handlers
+		virtualPasskey(t, b, a.credPath, origin, "mark-mac", 82, now-100)
+		return b, a, n, origin + "/surface/approve/" + a.id
+	}
+	timed := func(t *testing.T, b *passkeyBrowser, cmd string, args map[string]any) (cardView, time.Duration) {
+		t.Helper()
+		var v cardView
+		t0 := time.Now()
+		b.do(t, &v, cmd, args)
+		return v, time.Since(t0)
+	}
+	deadline := func(t *testing.T, what string, took time.Duration) {
+		t.Helper()
+		if took < 9500*time.Millisecond || took > 25*time.Second {
+			t.Fatalf("%s ended after %v, want the 10 s deadline", what, took)
+		}
+	}
+
+	t.Run("get", func(t *testing.T) {
+		b, _, n, card := setup(t)
+		n.set("get-500")
+		if v, _ := timed(t, b, "openCard", map[string]any{"url": card}); v.State != "failed" || v.Status != "The daemon did not give the request (HTTP 500: internal error (test))." || v.Again != "Read it again" || v.AgainHidden {
+			t.Fatalf("get-500: the card shows %+v", v)
+		}
+		for _, mode := range []string{"get-header", "get-body"} {
+			n.set(mode)
+			v, took := timed(t, b, "openCard", map[string]any{"url": card, "settleMs": 25000})
+			deadline(t, mode, took)
+			if v.State != "failed" || !v.StatusError || v.Status != "The request could not be read (no answer within 10 s). Nothing to approve yet." ||
+				v.AgainHidden || v.Again != "Read it again" || !v.ButtonHidden || len(v.Fields) != 0 {
+				t.Fatalf("%s: the card shows %+v", mode, v)
+			}
+			n.set("")
+			if v, _ = timed(t, b, "press", map[string]any{"id": "again", "states": []string{"ready", "failed", "refused"}}); v.State != "ready" {
+				t.Fatalf("%s: the retry shows %+v", mode, v)
+			}
+		}
+	})
+
+	t.Run("post", func(t *testing.T) {
+		b, a, n, card := setup(t)
+		if v, _ := timed(t, b, "openCard", map[string]any{"url": card}); v.State != "ready" {
+			t.Fatalf("the card shows %+v", v)
+		}
+		for _, tc := range []struct{ mode, after, status string }{
+			{"post-header", "ready", "Not recorded: the request is still pending with no approval. You may approve it again."},
+			{"post-body", "ready", "Not recorded: the request is still pending with no approval. You may approve it again."},
+			{"post-recorded", "recorded", "Recorded: approved by webauthn:mark-mac, waiting for the minter (approve " + a.id + " --await)."},
+		} {
+			n.set(tc.mode)
+			v, took := timed(t, b, "press", map[string]any{"id": "go", "states": []string{"unknown", "recorded", "rejected", "failed"}, "settleMs": 25000})
+			deadline(t, tc.mode, took)
+			if v.State != "unknown" || !strings.HasPrefix(v.Status, "Approval recording outcome unknown: no answer within 10 s.") || v.AgainHidden || v.Again != "Check status" {
+				t.Fatalf("%s: the card shows %+v", tc.mode, v)
+			}
+			if _, ok := surfaceRecorded(t, a.credPath, a.id); ok != (tc.mode == "post-recorded") {
+				t.Fatalf("%s: recorded %v", tc.mode, ok)
+			}
+			if tc.mode == "post-header" { // the status check has its own deadline, and offers another
+				n.set("list-header")
+				b.do(t, nil, "press", map[string]any{"id": "again"})
+				var said string
+				b.do(t, &said, "waitFor", map[string]any{"expression": `(() => { const s = document.getElementById('status').textContent; return s.startsWith('The status could not be read') && s; })()`, "ms": 25000})
+				if said != "The status could not be read (no answer within 10 s)." {
+					t.Fatalf("a stalled status check says %q", said)
+				}
+			}
+			n.set("")
+			if v, _ = timed(t, b, "press", map[string]any{"id": "again", "states": []string{"ready", "recorded", "used"}}); v.State != tc.after || v.Status != tc.status {
+				t.Fatalf("%s: the status check shows %+v", tc.mode, v)
+			}
+		}
+
+		// Recorded with no answer, then minted before the check: the check says it is no longer pending.
+		if v, _ := timed(t, b, "openCard", map[string]any{"url": card}); v.State != "ready" {
+			t.Fatalf("the card shows %+v", v)
+		}
+		n.set("post-recorded")
+		if v, _ := timed(t, b, "press", map[string]any{"id": "go", "states": []string{"unknown", "recorded", "rejected"}, "settleMs": 25000}); v.State != "unknown" {
+			t.Fatalf("post-recorded again: the card shows %+v", v)
+		}
+		n.set("")
+		var opened bool
+		aenv, _ := awaitEnv(time.Now().Unix(), &opened, nil)
+		var stdout, stderr bytes.Buffer
+		if code := runApprove(a.args("--await", "--await-for", "5s"), &stdout, &stderr, aenv); code != exitOK || opened {
+			t.Fatalf("approve --await: rc %d\nstderr %s", code, stderr.String())
+		}
+		b.do(t, nil, "press", map[string]any{"id": "again"})
+		var said string
+		b.do(t, &said, "waitFor", map[string]any{"expression": `(() => { const s = document.getElementById('status').textContent; return s.startsWith('The request is no longer') && s; })()`})
+		if said != "The request is no longer pending (used or expired): read the card again." {
+			t.Fatalf("after the mint the status check says %q", said)
+		}
+
+		a2 := newApproveRigIn(t, a.r, "ep2", time.Now().Unix()) // a pending request for the Busy arms
+		card = strings.TrimSuffix(card, a.id) + a2.id
+		busy := func(hold time.Duration) (cardView, int) {
+			t.Helper()
+			if v, _ := timed(t, b, "openCard", map[string]any{"url": card}); v.State != "ready" {
+				t.Fatalf("the card shows %+v", v)
+			}
+			w, err := store.OpenCredentials(a.credPath) // the daemon's non-waiting writer is refused: 503 Busy
+			if err != nil {
+				t.Fatal(err)
+			}
+			var once sync.Once
+			free := func() { once.Do(func() { _ = w.Close() }) }
+			defer free()
+			if hold > 0 {
+				time.AfterFunc(hold, free)
+			}
+			before := n.count()
+			v, _ := timed(t, b, "press", map[string]any{"id": "go", "states": []string{"recorded", "rejected", "unknown", "failed"}, "settleMs": 20000})
+			return v, n.count() - before
+		}
+		if v, posts := busy(1500 * time.Millisecond); v.State != "recorded" || posts != 3 {
+			t.Fatalf("a writer held 1.5 s: %d POSTs, the card shows %+v", posts, v)
+		}
+		if v, posts := busy(0); v.State != "rejected" || posts != 4 || !strings.Contains(v.Status, "(HTTP 503 Busy): the credentials store is busy") {
+			t.Fatalf("a writer held throughout: %d POSTs, the card shows %+v", posts, v)
+		}
+		if v, _ := timed(t, b, "openCard", map[string]any{"url": card}); v.State != "ready" {
+			t.Fatalf("the card shows %+v", v)
+		}
+		n.set("restart")
+		before := n.count()
+		if v, _ := timed(t, b, "press", map[string]any{"id": "go", "states": []string{"recorded", "rejected", "unknown"}}); v.State != "rejected" || n.count()-before != 1 || !strings.Contains(v.Status, "HTTP 503 Restart") {
+			t.Fatalf("a 503 Restart: %d POSTs, the card shows %+v", n.count()-before, v)
+		}
+		n.mu.Lock()
+		defer n.mu.Unlock()
+		if n.ssh != 0 || n.recordedRun != 2 {
+			t.Fatalf("%d POSTs sent an SSH signature; the withheld POST ran %d times", n.ssh, n.recordedRun)
+		}
+	})
+}
+
+// XA5, the passkey's own bound (A2 §2.14 [r2], design §2.7 step 4): a get that never answers is
+// aborted at 240 s of page time and nothing is sent; a result after the 240 s is ignored and
+// nothing is sent; an authenticator without user verification makes no response and nothing is
+// sent. Each leaves the button usable.
+func TestApproveCardPasskeyDeadline(t *testing.T) {
+	b := startPasskeyBrowser(t)
+	now := time.Now().Unix()
+	a := newApproveRig(t, "ep1", now)
+	n := &cardNet{release: make(chan struct{})}
+	origin := cardOrigin(t, a.credPath, "", n.wrap)
+	t.Cleanup(func() { close(n.release) })
+	card := origin + "/surface/approve/" + a.id
+	for i, tc := range []struct{ inject, authenticator, state, status string }{
+		{"abort", "uv", "expired", "Timed out: no passkey answer within 4 minutes. Nothing was sent; try again."},
+		{"late", "uv", "expired", "Timed out: no passkey answer within 4 minutes. Nothing was sent; try again."},
+		{"", "no-uv", "failed", "The browser did not make a passkey response (NotAllowedError: "},
+	} {
+		b.do(t, nil, "tab", map[string]any{"authenticator": tc.authenticator, "inject": tc.inject})
+		if i == 0 {
+			virtualPasskey(t, b, a.credPath, origin, "mark-mac", 84, now-100)
+		}
+		var v cardView
+		if b.do(t, &v, "openCard", map[string]any{"url": card}); v.State != "ready" {
+			t.Fatalf("%s: the card shows %+v", tc.inject, v)
+		}
+		args := map[string]any{"id": "go", "states": []string{"expired", "failed", "rejected", "recorded"}}
+		if tc.inject != "" {
+			args = map[string]any{"id": "go", "virtualBudget": 250000}
+		}
+		var after cardView
+		b.do(t, &after, "press", args)
+		if after.State != tc.state || !strings.HasPrefix(after.Status, tc.status) || after.ButtonDisabled || n.count() != 0 {
+			t.Fatalf("%s/%s: %d POSTs, the card shows %+v", tc.inject, tc.authenticator, n.count(), after)
+		}
+		if tc.inject != "" && (after.Op == nil || !after.Op.Signal || after.Op.AbortedAfter == nil || *after.Op.AbortedAfter < 239000 || *after.Op.AbortedAfter > 241000) {
+			t.Fatalf("%s: the get was not aborted at 240 s: %s", tc.inject, opString(after.Op))
+		}
+	}
+}
````

### Patch M8d — Needs you in a live Chrome; QUICKSTART §17

````diff M8d
--- a/cmd/ailang-worldd/needs_you_chrome_test.go
+++ b/cmd/ailang-worldd/needs_you_chrome_test.go
@@ -0,0 +1,55 @@
+package main
+
+// Row 201 A3 sprint 3 M8d (design §2.8, §6 XA3's list half): Needs you in the surface app, in a
+// live Chrome.
+
+import (
+	"os"
+	"path/filepath"
+	"strings"
+	"testing"
+	"time"
+)
+
+// XA3, the list half (Decided is sprint 4's, plan Appendix C): /surface's now card shows the pending
+// request under "Needs you · 1", from GET /surface/requests.json; a second request appears on the
+// app's next refresh (no timer of its own); the first one's link opens its card, ready.
+func TestNeedsYouListsPendingRequests(t *testing.T) {
+	b := startPasskeyBrowser(t)
+	b.do(t, nil, "tab", nil)
+	now := time.Now().Unix()
+	a := newApproveRig(t, "ep1", now)
+	rig, err := os.ReadFile(filepath.Join("..", "..", "host", "surface", "testdata", "index-rig.json"))
+	if err != nil {
+		t.Fatal(err)
+	}
+	index := filepath.Join(t.TempDir(), "index.json")
+	if err := os.WriteFile(index, rig, 0o600); err != nil {
+		t.Fatal(err)
+	}
+	origin := cardOrigin(t, a.credPath, index, nil)
+	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface", "raw": true})
+	var needs string
+	b.do(t, &needs, "waitFor", map[string]any{"expression": `(() => { const e = document.querySelector('[data-needs]'); return e && e.textContent; })()`})
+	if want := "Needs you · 1ep1 · claude · 10 grants · expires in 10 min"; needs != want {
+		t.Fatalf("Needs you shows %q, want %q", needs, want)
+	}
+	newApproveRigIn(t, a.r, "ep2", now)
+	b.do(t, nil, "waitFor", map[string]any{"expression": `(window.World.poll(), 1)`})
+	b.do(t, &needs, "waitFor", map[string]any{"expression": `(() => { const e = document.querySelector('[data-needs]'); return e && e.textContent.startsWith('Needs you · 2') && e.textContent; })()`})
+	if !strings.Contains(needs, "ep1 · claude") || !strings.Contains(needs, "ep2 · claude") {
+		t.Fatalf("after a refresh Needs you shows %q", needs)
+	}
+	var href string
+	b.do(t, &href, "waitFor", map[string]any{"expression": `document.getElementById('needs-` + a.id + `').getAttribute('href')`})
+	var at struct{ Location string }
+	b.do(t, &at, "press", map[string]any{"id": "needs-" + a.id, "page": true})
+	if href != "/surface/approve/"+a.id || at.Location != origin+"/surface/approve/"+a.id {
+		t.Fatalf("the link %q opened %q", href, at.Location)
+	}
+	b.do(t, nil, "waitFor", map[string]any{"expression": `document.body.dataset.state === 'ready'`})
+	var view cardView
+	if b.do(t, &view, "card", nil); !strings.HasPrefix(view.What, "Approve request "+a.id+" for episode ep1") || len(view.Errors) != 0 || len(view.Violations) != 0 {
+		t.Fatalf("the card shows %+v", view)
+	}
+}
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -1837,3 +1837,21 @@ ailang-worldd approver revoke demo-mac --db /tmp/se-world/world.db
 ```
 
 The first `--webauthn-assertion` mints with no terminal and records `approved_by=webauthn:demo-mac`; it checks the assertion before any network call, against the enrolled, unrevoked passkey, the request's own message hash, the origin and a signature counter that must move forward. The second, the same file again, is refused as already used. `approve --list --decided` shows `webauthn:demo-mac`; `approver revoke` (a typed `y`) lists the sessions the passkey approved and revokes none of them, and from then on that passkey approves nothing. `approve --webauthn-assertion` also refuses a file over 16 KiB or not a regular file, an assertion for another request, a passkey that is not enrolled, and one made without user verification.
+
+## 17. Approve from the surface *(attended — pending first verbatim run)*
+
+The surface app lists the mint requests waiting for you under **Needs you**, and each one opens a request card at `/surface/approve/<request-id>`. The card shows every field of the request the daemon stored (the episode, agent, model, task, every grant, the worktree, the stores, the requesting daemon), checks that those bytes match the request's digest, and computes the message to sign itself; it refuses a request it cannot show in full. **Approve with passkey** signs it with a §16 passkey, and the daemon records the approval only after checking it as `approve --webauthn-assertion` does. The daemon never mints: `approve <request-id> --await`, running in another shell, waits for the recorded approval, checks it again, and mints with no terminal. This guide uses the §0 scratch store `/tmp/se-world/world.db` with its daemon serving at `127.0.0.1:7644` and the §16 passkey `demo-mac`, enrolled at `http://localhost:7644`; never a real store.
+
+```bash
+ailang-worldd episode start surface-demo --repo "$PWD" --db /tmp/se-world/world.db --request
+ailang-worldd approve mr-<the 16 hex it printed> --await --db /tmp/se-world/world.db &
+```
+
+`--await` prints the request's id, the first 16 hex of its message hash and the card's address, `open http://localhost:7644/surface/approve/mr-<id>`. Open `http://localhost:7644/surface`: **Needs you · 1** lists `surface-demo` with its agent, grant count and minutes to expiry. Follow it to the card, check every field and that the message hash matches the one `--await` printed, press **Approve with passkey** and pass your device's check. The card shows `Approved by webauthn:demo-mac`; the waiting `approve` mints and prints `approved by webauthn:demo-mac`, and `approve --list --decided` shows it.
+
+```bash
+wait
+ailang-worldd approve --list --decided --db /tmp/se-world/world.db
+```
+
+Open the same card again and press **Approve with passkey**: the daemon refuses it as already used. With an SSH approver (§15) instead, **Save the message** on the card, sign it with the `ssh-keygen -Y sign -n world-approval@sunholo` line the card shows, paste the signature and press **Approve with this signature**. Every request and answer has a 10-second limit: a card that cannot be read says so and offers to read it again, and an approval whose answer never arrives is shown as *outcome unknown* with a status check, never as refused or minted. A busy credentials store is tried again at most three times, a second apart. Reached through `tailscale serve`, start `serve` with `--surface-origin https://<host>.ts.net`, or the card and its Approve are refused (421, 403) naming that flag.
````

