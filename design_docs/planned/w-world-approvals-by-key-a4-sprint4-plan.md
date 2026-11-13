# Sprint plan — w-world-approvals-by-key-a4, sprint 4 "the surface", headed by the sprint-3 judge's notes (iteration 392, row 201)

> **THE SPRINT — A4 sprint 4, the last A4 sprint (design §0 "Sprint 3", §1.12, §2.8, §2.9, §2.10, §6 "Sprint 3: the
> surface", §7, §8, §9 "Sprint 3", §10, §12 Q2), headed by the sprint-3 plan's Appendix C:** `requests.json` gains
> **`decisions` + `decisionsTruncated`** (DMD1–DMD3); **`approve --list`** prints them (Q2); ***Needs you*** lists
> each pending decision with a link to its card, the reserved text without a link, and the notification for a new
> decision id (XD3); **`DecisionSubject.Passkeys`** (sprint-1 P7); the **decision card** `GET /surface/decide/{ref}` +
> `GET /surface/decide.js` (`DECISION_FIELDS`, both hash checks and `requestRef`, Approve/Deny with a passkey over
> `DecisionMessageHash` with a **`BigInt`** expiry, **Save the message** + the signing line + a paste box; **PD1**,
> **PD2**) proven in **live Chrome** (XD1 approve and deny with Chrome's virtual authenticator, XD2 fails closed, XD4
> the saved 91 bytes and a pasted SSHSIG, with a real `ssh-keygen` sibling); **`TestSurfaceAppRoutes` +2**;
> **QUICKSTART §18's card paragraph**; and the sprint-3 judge's non-blocking notes 2–4 as tests (note 1 is a
> narrow-run survivor, §1.4 P4). **The card shows its own verb**; the CLI's `--webauthn-link` path keeps the frozen
> A2 page's "Approve request …" wording as a named residual (S4-D2).
> - **Measured: 1,947 changed lines** — `git diff --numstat 799660ac HEAD -- cmd host docs`, insertions + deletions,
>   after a commit of each of the nine scratch stages (**655 production** — 286 Go + 369 page/script, **1,284 test**,
>   **8 docs**; `-- cmd host` alone: **1,939**). Margin to the 2,000 cap: **53**. **Nothing is cut** (S4-D1).
> - **After this sprint A4 is code-complete** except Mark's attended exit run (design §10), whose step 4 (the phone,
>   *Needs you*, the card, Face ID) is now runnable.
> - **Nine premises of the design or of HEAD are false or incomplete and are routed minimally** (§1.4): `approve
>   --list` never contacts a daemon (P1, P2: S4-D3 reads the world store offline through the daemon's own listing
>   function); the listing's refusal needs `ValidatePacket`'s own scan (P3); the judge's `X-W15-OPERAND` survivor is
>   killed by an A2 test its narrow run skipped (P4); `passkey.js` is byte-pinned by an A3 test (P5); the
>   notification lives outside `needsYouHTML` (P6); two sprint-1 tests pin the subject's key set and the decision
>   routes and are edited by one line each (P7, P8); the design's PD1 "json tags" are of unexported wires (P9).
> - **Every line of code, test and doc is a unified diff in Appendix P**, generated from compiled, vetted,
>   gofmt-clean, green scratch stages; extracted **from this file** with §2's `awk`, applied to a fresh
>   `git archive 799660ac` export and compared byte-for-byte with the stage tree and the mutant tree (§1.2 C16).
>   The executor applies each patch and runs its acceptance; it writes no code by hand.

- **Design**: [w-world-approvals-by-key-a4.md](w-world-approvals-by-key-a4.md), **r2** (reviewed in iteration 388;
  **not reopened**). **D-WORLD-115** is OPEN; its default **A = keep it** is applied (S4-D9). Where this plan and
  the design differ, this plan was measured against `799660ac` and wins; each difference is a decision (§1.3) or a
  premise found false (§1.4).
- **Starting ground**: A4 sprints 1–3 landed (PRs #563, #566, #569; sprint 3 merge `e6761873`, judge sonnet PASS 91;
  recorded at `799660ac`, PR #570). Sprint-3 plan: [w-world-approvals-by-key-a4-sprint3-plan.md](w-world-approvals-by-key-a4-sprint3-plan.md)
  (its Appendix C is this sprint's head; its §4 survivors are carried in §4 here).
- **Planner base**: `799660ac4521919e2efe9c864fa8f5a0d7666ac8` (= origin/dev). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter392-a4-s4`, branch `sprint/world-iter392-row201-a4-s4`.
  Planner scratch (outside the repo, never under `/tmp`): `$HOME/.ailang/state/world-iter392/plan-scratch/`
  (`base/` the pristine export (a throwaway git repo, tag `BASE`), `wJ`, `wA`, `wB`, `wC` the four slice trees the
  first mutant runs used (each = base + one slice; §4), `work/` the final tree (all slices; the later mutant runs),
  `work2/` the stage repo (tags `M7a` … `M7i`), `final-tree/` (work2 at `M7i`, the final gate run), `patches/`,
  `tools/` (`gates.sh`, `gates_final.sh`, `accept.sh`, `accept_all.sh`, `stages.sh`, `STAGES`, `mutate.py`,
  `mutants_{J,A,B,B_js,C}.py`, `mtable.py`, `roundtrip.sh`, `assemble.py`, `mkjson.py`), `logs/`, `roundtrip/`).
- **Executor**: as routed by the controller. **Three runs**: Run 1 (M7a, M7b, M7c — the judge's notes, the listing,
  `--list`), Run 2 (M7d, M7e, M7f — *Needs you*, `Passkeys`, the card), Run 3 (M7g, M7h, M7i — live Chrome and
  QUICKSTART). Each sub-milestone is one `patch -p1` (§2) plus its acceptance (§3). **Apply, never retype.**
- **Touched** (§5): new `host/surfacepage/{decide_list.go, decide.go, decide/index.html, decide/decide.js}`,
  `host/daemon/{surface_decide_list.go, surface_decide_page.go}`, `cmd/ailang-worldd/approve_list_decisions.go`;
  edited `host/surfacepage/{approve_model.go (+3: two fields), decide_model.go (+7), app/app.js (+25 −5)}`,
  `host/daemon/{surface_approve.go (+6), surface_decide.go (+25 −1), daemon.go (+2 routes)}`,
  `cmd/ailang-worldd/approve.go (+1 −1)`, `docs/QUICKSTART.md` (+8, §18); new tests in `host/approval`,
  `host/surfacepage`, `host/daemon`, `cmd/ailang-worldd`; **edited existing tests**: `TestSurfaceAppRoutes` (the
  design's §7 +2), `TestDecisionSubjectJSON` (+1 key, P7), `TestDecisionRoutesRegistered` (P8). **Outside the
  design's §8 file list**: `host/surfacepage/approve_model.go` (A3a's `RequestsJSON` gains the two members the
  design names), `host/daemon/surface_approve.go` (A3a's handler calls the listing), `host/surfacepage/decide_model.go`,
  `host/daemon/surface_decide.go` (sprint 1–3's own files), the new files and every `_test.go` above.
- **Not touched**: `host/approval/*.go` production, `host/authority/**`, `host/store/**`, `host/broker/**`,
  `host/coordinator/**`, `host/surfacepage/{approve/**, passkey/**, approve.go, passkey.go, surface.js}`,
  `host/daemon/{surface_approve_page.go, surface_decide_key.go, episode_finish.go}`, every sprint 1–3 POST arm and law,
  `cmd/ailang-worldd/{approve_decide.go, main.go, session.go, help_test.go}`, `world/*.ail`, `packages/**`,
  `scripts/**`, `tools/launchd/*`, `cmd/world-publish/*`, any skill, `~/.ssh`, `~/.ailang/world`,
  `~/.ailang/state/mission-v1*`.

---

## 0. Binding constraints (each is restated where it governs)

1. **Test keys and temp stores only (D-WORLD-109 = A).** Every store is under `t.TempDir()`; every passkey is
   `approvaltest.NewPasskey(alg, seed)` or Chrome's **virtual** authenticator in a temp `--user-data-dir` (the A2
   driver, `testdata/passkey_cdp.mjs`); every SSH key is `approvaltest.NewSigner(seed)` or, for the real-tool
   sibling, an ed25519 key made by `ssh-keygen -q -N ""` in a `t.TempDir()` with `SSH_AUTH_SOCK=`. The pages talk to
   an `httptest` loopback daemon built on `d.Handler()`. **Any test that runs the CLI against a daemon resolves the
   daemon explicitly** (`--daemon <httptest URL>`) or first proves that `defaultServiceAddr()` names the fake daemon
   (sprint-3 §1.4 P7, `approve_decide_jgaps_test.go:31`); no test, and no planner run, reaches `127.0.0.1:7644`.
   No real store, `~/.ssh`, `/dev/tty`, browser profile, mint or publish.
2. **Pure core, effects at the boundary (coding standards S1–S3).** The `decisions` fold is pure (`host/surfacepage`,
   no store, no clock: `now` is an input); the daemon reads the world store on its read handle and the credentials
   store on a fresh read-only handle, as A3a's `requests.json` does. `host/approval` is not touched. The card page
   computes nothing the POST trusts: the POST re-runs the law (D1–D10) and the arm unchanged.
3. **`decidedBy` is derived, never accepted.** The card POSTs sprint 3's grammar unchanged (no `decidedBy` member
   exists; RD11 refuses one).
4. **No `/v1` route, no CSP change, no CORS header, no new Go dependency, no schema change.** Exactly the design's two
   new `/surface` routes (`GET /surface/decide/{ref}`, `GET /surface/decide.js`), both behind `approvalRoute`;
   `TestFrozenV1RouteTableMatchesMux` stays 13; `TestSurfaceAppRoutes` +2 (design §7); `surfaceCSP` unchanged.
5. **The A2/A3 pages and the sprint 1–3 laws are not re-opened.** `passkey.js`, `/surface/passkey`, `approve.js`,
   `/surface/approve/{id}` and every POST arm are unchanged (design §2.8, §13); their byte pins
   (`surface_approve_page_test.go:54`, PG2) stay unedited.
6. **One mutant per new production branch, anchored to the diff** (§4): the named test and the package suite;
   a row the fast suite does not kill is re-run on the **whole package** before it may be called SURVIVED; every
   run appends a ledger line; the table's ids are the ledger's ids.
7. **No `t.Parallel()`** in any new test (`TestNoParallelWriteForkPackagesOutsideVerifygate`).
8. **Bisectable and safe at every commit.** Each sub-milestone ends gofmt-clean, `go build ./...` and `go vet`
   green, the `_test.go` compile fence green, its own tests and every earlier one green.
9. **No git writes by the executor** (§7). The controller commits from snapshots.

---

## 1. Pre-flight

### 1.1 Baseline at `799660ac` (measured by the planner, 2026-11-13)

Every command ran with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` (AILANG v0.41.0), `go1.26.6 darwin/arm64`,
in a `git archive 799660ac` export at `plan-scratch/base` (a throwaway git repo; `diff -r cmd`, `host`, `docs`
against the worktree: identical, `BASE_EQUALS_WORKTREE`; the worktree itself was only read). The legs are
`tools/gates.sh <tree> <logs> base` = §6's script with the sprint's new-test lists emptied. Load average 9–10 at the
start (`uptime` 05:35Z), 18–28 while this planner's four slice trees compiled, ran Chrome and mutated alongside.
Logs: `plan-scratch/logs/base_gates/` (`legs.txt`, one log per leg).

| Leg (`tools/gates.sh <base export> logs/base_gates base`) | rc | Time | Result |
|---|---|---|---|
| `status` (`git status --porcelain \| wc -l` in the export; the worktree itself: `0` at `799660ac`) | 0 | 0 s | `0` |
| `verify_ail` (`./scripts/verify_ail.sh`) | **0** | 56 s | `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` |
| `gofmt`, `build`, `vet`, `fence` | **0** | 0 / 6 / 6 / 47 s | empty / clean / clean / no FAIL line |
| `pkgs` (`host/approval/...`, `host/authority`, `host/store`, `host/surfacepage`, `host/coordinator`, `host/daemon`, `cmd/ailang-worldd`) | **1** | 775 s | every package `ok` but `host/daemon`, red **only** on `--- FAIL: TestWorkbenchPacketRosterTestedAndProven/not-ok` ("store: put object outcome is uncertain … context deadline exceeded", load average 18–24 while three planner forks compiled and mutated); **re-run alone: `ok host/daemon 455.8s`, rc 0** (`logs/base_roster_alone.log`) — a load timeout in an untouched test |
| `broker_named` | **0** | 26 s | every named test PASS |
| `race` (three `-race` runs) | **1** | 218 s | `ok host/surfacepage`, `ok host/daemon`; `cmd/ailang-worldd` red on `TestApproveAwaitStopsAtExpiry/expiry` ("stopped at … after 507 waits") and `TestApproveBySignatureRefusalsBeforeAnyEffect/unenrolled-key` ("begin: context deadline exceeded") — both clock/deadline-under-load in untouched A1b/A3 tests (the directive's "local `-race` … host/daemon" class); no `DATA RACE` line |
| `inv` (every named inventory, `-v`, base mode: no sprint test) | **0** | 82 s | **51** named `--- PASS`, 0 SKIP, 0 FAIL, none missing |
| `chrome` (`WORLD_CHROME=… WORLD_REQUIRE_PASSKEY_BROWSER=1`, the six A2/A3 page tests + `TestApproveDecideWebAuthnLinkInChrome`) | **0** | 32 s | all seven `--- PASS` |
| `verify_go` (`VERIFY_GO_LEGS=plain AILANG_BIN=… ./scripts/verify_go.sh`, the export is a git repo with the base committed) | **1** | 719 s | 57 packages `ok`; red: `TestExecProtectedCheckTimeout` (host/broker), `TestRegistryPhaseHasItsOwnClock` (host/daemon), `TestMissionConfigDefaultFlowIntegration` (host/verifygate), `TestMarketingSubprocessDeadline` (scripts/marketingcheck) — all on the directive's known-flaky list — and `TestWorldContextAndStatusline` (cmd/ailang-worldd, `world-context: signal: killed`, a subprocess killed at its deadline under load) |
| whole script | — | 1,969 s | `GATES_DONE 1969s FAILS: pkgs race verify_go` (every red a load timeout in an untouched test; the `pkgs` one green alone) |

**Rule (carried):** a red at the base measures the repo, not the sprint. The controller re-runs a red untouched
package alone once and does not chase it. **Never start a gate under `nohup`.**

### 1.2 Codebase claims checked (at `799660ac` unless "work" — the planner's final scratch tree)

| # | Claim | Command → observed output |
|---|---|---|
| C1 | The base | `git rev-parse HEAD` (worktree) → `799660ac4521919e2efe9c864fa8f5a0d7666ac8`; `git status --porcelain \| wc -l` → `0` |
| C2 | Controller-verified anchors hold | `grep -n 'GET /surface/requests.json\|GET /surface/decisions/{ref}\|POST /surface/decide/{ref}' host/daemon/daemon.go` → `1138:`, `1143:`, `1144:`; `grep -n "^func TestSurfaceAppRoutes" host/daemon/surface_app_test.go` → `23:`; `grep -n "const credentialsSchemaVersion = " host/store/credentials.go` → `44:const credentialsSchemaVersion = 8`; `grep -c Passkeys host/surfacepage/decide_model.go` → `0` with the known positive in the same call `grep -n "type DecisionSubject struct" host/surfacepage/decide_model.go` → `18:`; `grep -n "^## 18" docs/QUICKSTART.md` → `1861:` (the last section; the file is 1,888 lines); `$HOME/.pinned-ailang/ailang --version` → `AILANG v0.41.0` |
| C3 | **Drift since the design's base is exactly sprints 1–3 and their records** | `git diff --stat b0fbcd3e 799660ac -- cmd host \| tail -1` → `66 files changed, 5238 insertions(+), 60 deletions(-)`; the anchors this sprint reads re-measure: `app.js` 1,470 lines (design §8: 1,470), `function needsYouHTML` at `app.js:427`, `function notifications` at `:446`; `packetScanLimit = 100` (`decisionpacket_validate.go:34`); `RecentApprovals` `maxHeads = min(maxHeads, 40)` (`approvals_view.go:61`) and `if len(rows) > 20` (`:103`) |
| C4 | **P1 (design §1.12) carries from sprint 3, and the card now proves it end to end** | sprint-3 plan §1.2 C4 (`world-iter391/plan-scratch/p1/`): Chrome 157 rebuilt DM1's two 91-byte messages and hashes under `surfaceCSP` with `TextEncoder`, `setBigUint64(83, BigInt(expiry))`, `crypto.subtle.digest` → `msg=EQUAL hash=EQUAL`. work: `decide.js` builds the message the same way (`m[18] = verb`, P at 19, Q at 51, `setBigUint64(83, BigInt(expiry))`, P and Q hashed from the served bytes); XD1 (`TestDecidePageApproveWithPasskey`) checks each `navigator.credentials.get` challenge equals Go's `approval.DecisionMessageHash(verb, P, Q, window.max)` for an approve and a deny, and XD4 checks the saved file equals `approval.DecisionMessage` byte for byte (91 bytes) |
| C5 | **`approve --list` never contacts a daemon** (§1.4 P1, P2) | `approve.go:264` `if *list && *decided {`, `:268` `rows, err := ro.PendingMintRequests(ctx, env.now(), 0)`; `grep -n "/surface/requests.json" cmd/ailang-worldd/*.go \| grep -v _test` → empty (rc 1), known positive in the same call `grep -n '/surface/decisions/' cmd/ailang-worldd/approve_decide.go` → `140:`; `approve.go:221` CD3 refuses `--daemon` with `--list` |
| C6 | **`passkey.js` is byte-pinned** (§1.4 P5) | `surface_approve_page_test.go:54` `"passkey.js": "0064c6c90099a3565b178cf170618b80a67bbf1f3444a1018cc4ee29f51ed53f"} {`, `:56-57` `t.Errorf("the A2 page's %s changed (%v)", …)`; `passkey.js:131` `` `Approve request ${f.request}, expires …, message hash ${f.hash16}` `` |
| C7 | **The wires PD1 ties to are unexported** (§1.4 P9) | `host/broker/approve.go:68` `type approvalRequestWire struct` (effect, scope, cost, requester, now); `host/coordinator/decisionpacket.go:98` `type packetWire struct` (packetHash, proposalHash, requestRef, createdAt, deadlineAt, escalationsRemaining, policy) — the design's twelve fields, exactly |
| C8 | **Two sprint-1 tests pin what this sprint adds** (§1.4 P7, P8) | `surface_decide_subject_test.go:39` pins `DecisionSubject`'s exact key list; `surface_decide_routes_test.go:31,45` pins exactly two decision routes and `GET /surface/decide/{ref}` → 405 |
| C9 | The *Needs you* change needs no export | `app_test.go:174-175` the exported-name list; work: `needsYouHTML` and `notifications` extended in place, two non-exported helpers (`decisionLineHTML`, `RESERVED_OWNERS`); `TestSurfaceAppPure` PASS unedited (`WORLD_EXEC_NODE=/opt/homebrew/bin/node`, node v26.13.0) |
| C10 | **The listing and the CLI share one function** | work: `host/daemon/surface_decide_list.go` `func PendingDecisions(ctx, s *store.Store, now int64)`; `readDecisionList` (the handler) and `approveListDecisions` (`cmd/ailang-worldd/approve_list_decisions.go`, on `store.OpenReadOnly(--db)`) both call it |
| C11 | No new context root, credentials method, writer site, `/v1` route or dependency | work: `--list` reuses `runApprove`'s `approveReadBudget` context; `git diff BASE M7i -- host/store host/broker host/coordinator 'host/approval/*.go' ':!*_test.go' \| wc -l` → `0`; `TestProductionContextRoots`, `TestCredentialsMethodsRequireADeadline`, `TestDaemonCredentialWriterSites`, `TestFrozenV1RouteTableMatchesMux` (13), `TestDaemonDependencyAllowlist` PASS unedited |
| C12 | **Nothing outside §5 moved** | stage repo: `git diff --stat BASE M7i -- world packages scripts tools design_docs \| wc -l` → `0`; production Go: `git diff --stat BASE M7i -- '*.go' ':!*_test.go' \| tail -1` → `11 files changed, 284 insertions(+), 2 deletions(-)`; added lines matching `t.Parallel\|Access-Control\|OpenCredentials(\|context.Background()\|/dev/tty\|.ssh` under `cmd host`: 5, **all in `_test.go`** (test contexts, the rig's credentials store, a CORS-absence assertion), 0 in production; `git diff --name-only BASE M7i -- host/surfacepage/passkey host/surfacepage/approve host/daemon/surface_approve_page.go` → empty |
| C13 | **With `WORLD_CHROME` set, an unrelated cmd test is red at the base tree** | `TestMarketingChromeSixShots` needs `AILANG_BIN and WORLD_TOOL_AILANG_BIN`; every Chrome line in §3/§6 names its tests with `-run` (never the whole package under `WORLD_CHROME`); the JS mutant suite skips `TestMarketingChrome*` (`mutants_B_js.py`) |
| C14 | **Every sub-milestone is green when committed** | `tools/accept_all.sh M7a … M7i` (an APFS clone of the stage repo checked out at each tag; `logs/accept_<M>.log`): **`ACCEPT_OK` for every one of the nine** — each stage gofmt-clean, `go build ./...` and `go vet` of the touched packages clean, the `_test.go` compile fence without a FAIL line, its §3 `-v` tests `--- PASS` with no SKIP/FAIL, its package runs `ok`; M7g's live Chrome: `--- PASS: TestDecidePageApproveWithPasskey`, `TestApproveDecideWebAuthnLinkInChrome`, `TestPasskeyPageInChrome`; M7h's: `--- PASS: TestDecidePageFailsClosed`, `TestDecidePageSSHSIGPaste` (its real-`ssh-keygen` arm ran with `WORLD_REQUIRE_SSH_KEYGEN=1`), `TestApprovePageInChrome`, `TestApprovePageRefusesASwappedRequest` |
| C15 | **Every sub-milestone is committed and measured** | `tools/roundtrip.sh <this plan>` commits each extracted patch in a throwaway repo and measures `git diff --numstat <base> HEAD -- cmd host docs`: `NUMSTAT M7a 106`, `M7b 437`, `M7c 587`, `M7d 719`, `M7e 810`, `M7f 1350`, `M7g 1635`, `M7h 1878`, `M7i 1947 (cmd host: 1939)` (`logs/roundtrip_final.log`); the stage repo's own `logs/stage_numstat.txt` agrees row for row |
| C16 | **The patches reproduce the stage tree and the mutant tree** | `tools/roundtrip.sh <this plan>`: a fresh `git archive 799660ac` export + `patch -p1` of Appendix P's nine patches **extracted from this file with §2's `awk`** (each `cmp`-equal to the generated patch; `--dry-run` first; no `.orig`/`.rej`) → `diff -r` of `cmd`, `host`, `docs` against `plan-scratch/work2` (`M7i`) and against `plan-scratch/work` (the final tree the last mutant runs used): **identical** (`ROUNDTRIP_IDENTICAL_TO_STAGE_TREE_AND_MUTANT_TREE`, `logs/roundtrip_final.log`); `work2` is built from `work` by `tools/stages.sh`, which ends `WORK2_EQUALS_FINAL` |
| C17 | **Every mutant; the table's ids are the ledger's** | 126 rows, the last ledger line of each: **113 `K`, 7 `E`, 6 `L`, 0 not as expected**; 148 ledger lines (runs J1–J2, A1, C1–C2, B1–B5js on the slice trees; F1, F2, F2js on the final tree for every `E`/`L` row), 1 `NOCOMPILE` line (`P-PASSKEYS-FIELD` in B1, rewritten and killed in B4); `comm -3 <(cut -f4 logs/mutant_ledger.tsv \| sort -u) logs/plan_ids.txt` → empty (§4) |
| C18 | **The gates, on the final tree** | `tools/gates_final.sh <final-tree> logs/final_gates final` (`final-tree` = `work2` at `M7i`, a git repo with every stage committed): `verify_ail` → `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` (no `.ail` changed); `gofmt`, `build`, `vet`, `fence` rc 0; `pkgs` rc 1: `ok` for host/approval, approvaltest, host/authority, host/store, host/surfacepage; red **only** on `TestSeToolsPlansAdmittedByPlanLaw` (host/coordinator), `TestExecStartupRefusalTable/episode-bad-id` (host/daemon), `TestDemoProofUncertainCommitStops/lost-response-after-commit` (cmd/ailang-worldd) — each "archived interpreter … `ailang --version`: timed out after 10s" on a freshly copied binary under load, none reaching this sprint's code; **all three re-run alone: `ok` (4.3 s, 22.6 s, 4.0 s; `logs/final_reds_alone.log`)**; `broker_named` rc 0; `race` rc 0 (`ok` host/surfacepage, host/daemon (decide/approve/requests/app tests), cmd/ailang-worldd (approve tests); no `DATA RACE`); `inv` rc 0: **69** named `--- PASS` (51 as at the base + 18 new), 0 SKIP, 0 FAIL, none missing — incl. `TestFrozenV1RouteTableMatchesMux` (13), `TestSurfaceAppRoutes`, `TestDaemonCredentialWriterSites`, `TestProductionContextRoots`, `TestProductionGoSurface`, `TestCredentialsMethodsRequireADeadline`, `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist`, `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`, `TestEveryVerbHasHelp`, `TestApprovalImportsOnlyTheStandardLibrary`, `TestSurfaceAppPure`, `TestSurfaceScriptAssetSeams`, `TestApprovePageRoutes`, `TestPasskeyPageRoutes`, `TestGETRoutesRejectOtherMethods`, the quickstart tests; `chrome` rc 0: **ten** `--- PASS` (the seven base page tests + `TestDecidePageApproveWithPasskey`, `TestDecidePageFailsClosed`, `TestDecidePageSSHSIGPaste`); `verify_go` rc 1: **61** packages `ok`, red only on `TestWorkspaceExecEpisodeSelectsItsProfile` (host/daemon, "ep1 ran under \"\", want alpha", an untouched exec-profile test of the directive's `TestWorkspaceProtectedPolicyWiring` family) — **re-run alone: `ok host/daemon 7.9s`** (`logs/final_wsexec_alone.log`); the base's five `verify_go` reds did not recur. **`GATES_DONE 1326s FAILS: pkgs verify_go`**, every red an untouched load/environment test green alone; `ALL_DONE.<sha>` written (`logs/final_gates/`) |

### 1.3 Planner decisions where the design leaves a choice

- **S4-D1 — The scope: all of sprint 3's Appendix C, plus three of the sprint-3 judge's notes; nothing cut.**
  Measured: **1,947** changed lines (`-- cmd host docs`; 1,939 `-- cmd host`), margin 53. The cap's least operator-visible part would have been `approve --list` (M7c) and the
  judge-gap tests (M7a); neither had to move. **PD1 ships with the page** (M7e), and **no POST arm or page ships
  without an operator path**: *Needs you* → the card → the unchanged `POST /surface/decide/{ref}` arms (§3.1).
- **S4-D2 — The deny verb on screen (scope item 7).** The **card** shows its own verb: its passkey buttons read
  **Approve with passkey** / **Deny with passkey**, its status line names the verb being signed, and its SSH half
  saves `world-decision-<verb>-<ref16>.msg` and posts with **Approve with this signature** / **Deny with this
  signature**; XD1 asserts the on-screen verb for an approve and for a deny (§3 M7f). The **CLI's
  `--webauthn-link` path keeps the A2 page's wording, as a named residual (RS-VERB)**, because: (a) the design
  freezes that page — §2.8 "`approve.js`, `/surface/passkey` and `passkey.js` are not changed", §0 "Out of A4: any
  change to … `/surface/passkey`", §13; (b) `passkey.js`'s bytes are pinned by SHA-256 in an A3 test the design
  lists as "Must stay green, unedited" (`TestApprovePageRoutes`, `surface_approve_page_test.go:54`, §1.4 P5); (c)
  the link grammar the page reads (`v`, `op`, `origin`, `rp`, `request`, `challenge`, `expires`, `hash16`) is the
  mint path's too. **Measured alternative** (not built): a `verb` link member read at `passkey.js:131` (≈ 4 lines),
  the CLI's link builder (≈ 2), the pin's hash and a Chrome arm (≈ 25) — ≈ 30 lines, all in A2/A3 files the design
  freezes. The path still carries the verb three ways (the CLI's capitals, the link's `request` text
  `decide-deny-<16 hex>`, and the signed hash the page shows); QUICKSTART §18 already says the page reads
  *Approve request decide-deny-…* for a deny; and the phone's path is now the card. A design amendment (one line in
  §2.8 and §13) would be needed to change it — recorded for the attended exit run (Appendix C).
- **S4-D3 — `approve --list` reads the world store offline, through the daemon's own listing (design Q2).** At HEAD
  `--list` never contacts a daemon (§1.4 P1) and refuses `--daemon` (CD3, P2). Making it GET `requests.json` would
  add a dial to every `--list`, including the existing tests that pass no `--daemon` and would fall back to
  `defaultServiceAddr()` — the sprint-3 P7 near miss (D-WORLD-109). Instead `--list` opens `--db` with
  `store.OpenReadOnly` and calls `daemon.PendingDecisions`, **the same function `requests.json` serves**, so the CLI
  and the surface list the same rows by construction; it never dials. A store that cannot be read is one stderr
  line; the mint list and the exit code are unchanged.
- **S4-D4 — The listing's reads (design §2.8).** One `RecentApprovals(40)` read (≤ 20 summaries), one packet scan of
  101 (`decisionListScan` = `packetScanLimit` + a sentinel: 101 → `decisionsTruncated`), `GetObject` + decode per
  scanned packet, then **`ValidatePacket` once per request carried by exactly one scanned packet (≤ 20)** so the row's
  refusal is the very sentence `GET /surface/decisions/{ref}` shows (`packetRefusalSentence`), and
  `broker.SubjectRefusal` for D4b/D5/D7 with the scan treated as complete (display only; the POST runs the
  complete discovery, design [r1]). No new JSON member: a world-store error is logged and `decisions` is `[]`,
  `decisionsTruncated` false, the mint lists untouched. The pure fold (`surfacepage.FoldDecisions`) keeps pending
  `Human.Approve` requests with ≥ 1 packet; `packet`/`deadline` only for exactly one.
- **S4-D5 — *Needs you*'s lines (design §2.8 "app.js").** Decidable: `Decide · <scope> · requested by <requester> ·
  deadline <hh:mm>` linking `/surface/decide/sha256:<64 hex>` (`id="decide-<64 hex>"`); reserved: "<owner> approval
  — use world-publish approve | ailang-worldd episode finish at a terminal", no link; a non-reserved refusal:
  "Decide · <scope> · refused: <reason>", no link (the card cannot load a D4b subject); truncated: one line naming
  `ailang-worldd approve <request-ref>`. The count includes decisions. The notification fires for a new
  **decidable** decision id only (reserved/refused are seen, never raised), as A3 does for refused mint requests.
  `needsYouHTML` and `notifications` are extended in place: no new export, so `TestSurfaceAppPure` is unedited.
- **S4-D6 — `DecisionSubject.Passkeys` reuses A3a's shape.** `passkeys: [{origin, rp, allow}]` is exactly design
  §2.8's JSON and A3a's `CardPasskeys`; `surfacepage.DecisionPasskeys` is one call into `RequestCard`'s grouping
  (not a copy), so `decide.js` copies `approve.js`'s ceremony with one change (`rpId: group.rp`). The daemon reads
  `ActivePasskeys` on a fresh read-only credentials handle in `handleSurfaceDecision`; a store below schema 6 lists
  none (A3a's rule); any other error is the sanitized 500. Only the card reads it (sprint-3 S3-D7: the CLI reads its
  own store).
- **S4-D7 — PD1 ties `DECISION_FIELDS` to the real wire, not to struct tags** (§1.4 P9). `TestDecidePageFieldTable`
  (`host/daemon`) reads the table between `/* DECISION_FIELDS */` markers in `decide.js` and compares it with the key
  and JSON-type set of a request and a packet **written by broker's and coordinator's own encoders** and served by
  `GET /surface/decisions/{ref}` (control: 12 fields). A field added to either wire without the page fails it.
  `decide.js` refuses an unlisted key, a missing listed key, a wrong JSON type, a non-safe integer, a fraction, a
  non-object and bad base64url — each a live-Chrome XD2 case.
- **S4-D11 — The live-Chrome rig posts to the real law and judge, not to the daemon's POST over a seeded world.**
  The cmd-package rig (`decide_page_rig_test.go`) serves the **real** card bytes from `d.Handler()` and, for the
  subject and the POST, applies A3a's browser guards (Origin, `Sec-Fetch-Site`, JSON, no cookie/authorization), the
  real `surfacepage.ParseDecideBody`, `broker.DecisionRefusal` and `authority.JudgeDecisionAssertion` /
  `JudgeDecisionSSH` over a temp credentials store; seeding a world store with committed packets from `cmd/` would
  port the daemon's test-only commit helpers (≈ +150 lines, over the cap). The daemon's own `POST
  /surface/decide/{ref}` arm is unchanged and already proven by `TestDecidePOST*` with approvaltest bytes and by
  sprint 3's live-Chrome `TestApproveDecideWebAuthnLinkInChrome`; the card's body passes the same
  `ParseDecideBody` in both. Residual RS-RIG (Appendix C; the attended run's step 4 covers it).
- **S4-D8 — The sprint-3 judge's non-blocking notes (evidence `world-iter391/evaluator-r1.md` §"Non-blocking
  findings").** 2 (`hErr` placement): `TestDecideByAssertionRefusesAMissingMessageAtW6` pins W6's place (a missing
  message with an earlier failing W1/W2/W5 is refused by that step; ledger `JG-HERR-TOP`, killed). 3 (the packet
  the CLI posts with `--signature`/`--webauthn-assertion`): `TestApproveDecideArtifactPostsTheSubjectsDigests`
  (`JG-BODY-PACKET`, `JG-BODY-REQUEST`). 4 (the two read caps): `TestApproveDecideArtifactReadBounds` (one byte over
  each bound refused before any POST, the bound itself read; `JG-SIG-CAP`, `JG-ASR-CAP`). 1 (the mint law's W15
  operand, `X-W15-OPERAND`): **not a gap** — the judge's run was `-run '^(TestDecideByAssertion|TestApproveByAssertion)'`;
  the whole `host/approval` package kills it with `TestSignCountLaw`'s `{stored 5, new 0}` arm (ledger
  `JG-W15-OPERAND`, §1.4 P4). No production line changes for any of them.
- **S4-D9 — D-WORLD-115 (OPEN, default A).** The card builds the same `DecisionMessage` the laws hash. Under B (a
  stored daemon challenge in the P slot) `decide.js`'s message builder changes with `DecideBySignature`,
  `DecideByAssertion` and the CLI's CD6/CD8 (sprint-3 S3-D8); the listing, *Needs you*, `--list` and the routes do
  not. Built under **A**; nothing in this sprint depends on the answer otherwise.
- **S4-D10 — Sub-milestones** (§1.5): nine, each ≤ 350 lines except **M7f (540)**: the routes, the page, `decide.js`
  (300 lines alone) and PD1/PD2, which the design forbids separating ("PD1 is never separated from the page"); the
  card's live-Chrome proof is split into M7g (XD1) and M7h (XD2, XD4).

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

- **P1 — Design §2.10 "`approve --list` adds a 'decisions' section from `requests.json`'s `decisions`".** At HEAD
  `--list` reads only the credentials store and never contacts a daemon (`approve.go:264-283`: `if *list && *decided
  { return approveListDecided(…) }`, `ro.PendingMintRequests(ctx, env.now(), 0)`; `grep -n "/surface/requests.json" cmd/ailang-worldd/*.go | grep -v _test` → empty (rc 1), with the known
  positive in the same call `grep -n '/surface/decisions/' cmd/ailang-worldd/approve_decide.go` →
  `140: status, body, err := c.do(ctx, http.MethodGet, "/surface/decisions/"+o.ref, nil)`). Routed: S4-D3.
- **P2 — `--daemon` with `--list` is refused at HEAD** (`approve.go:221` CD3: "--deny, --daemon and --expiry take a
  Human.Approve request ref …"). A requests.json route would have had to reverse a sprint-1 rule; S4-D3 does not.
- **P3 — Design §2.8 "one scan of ≤ 101 packets per read".** The scan is one, but the refusal the design asks for
  ("its validated packet … refusal") needs `ValidatePacket`, which runs its own ≤ 101-packet duplicate scan
  (`decisionpacket_validate.go:190`) per call. Routed: S4-D4 (≤ 20 calls, display only; risk R5).
- **P4 — The sprint-3 judge's finding 1 (`X-W15-OPERAND` SURVIVED) is a narrow-run survivor.**
  `passkey_law_test.go:198` `{5, 0, false}` refuses exactly the input the mutant accepts; the planner's run of the
  mutant against the whole `host/approval` package: `JG-W15-OPERAND KILLED KILLED TestSignCountLaw` (ledger run J1).
- **P5 — `passkey.js` is byte-pinned by an A3 test.** `surface_approve_page_test.go:54` `"passkey.js":
  "0064c6c90099a3565b178cf170618b80a67bbf1f3444a1018cc4ee29f51ed53f"` ("the A2 page's %s changed"). Any verb text in
  the A2 page edits that pin (S4-D2).
- **P6 — Design §8 "`app.js`: one new block in `needsYouHTML`".** The design's own §2.8 also asks for "the opt-in
  notification … fires for a new decision id too", which lives in `notifications` (`app.js:446`). Both are extended
  in place (+25 −5); no other `app.js` function changes.
- **P7 — Design §6 "Must stay green, unedited" lists no sprint-1 A4 test, but `TestDecisionSubjectJSON` pins the
  subject's exact key list** (`surface_decide_subject_test.go:39`), so `DecisionSubject.Passkeys` (sprint-1 P7's
  own item) needs `"passkeys"` added to it: +1 −1 (M7e).
- **P8 — `TestDecisionRoutesRegistered` pins exactly two decision routes and `GET /surface/decide/{ref}` → 405**
  (`surface_decide_routes_test.go:31,45`); the card makes that GET a page. Edited +2 −2 (M7f): the route list gains
  the card's two routes and `GET` leaves the 405 list for that path. `TestSurfaceAppRoutes` gains the design's +2
  (`/surface/decide/sha256:…` and `/surface/decide.js` in its 421 list; `/surface/decide`, `/surface/decide/`,
  `/surface/decide.js/x` in its 404 list; +4 −2). No A2/A3 test, no A2/A3 page and not the A2 Chrome driver
  (`passkey_cdp.mjs`) is edited; PD2 re-pins the A2/A3 page hashes.
- **P9 — Design §6 PD1 "`DECISION_FIELDS` equals the json tags of `approvalRequestWire` and `packetWire`".** Both
  are unexported in other packages (§1.2 C7), so no daemon or surfacepage test can reflect on them. Routed: S4-D7.

### 1.5 The size

Lines are insertions + deletions: **`git diff --numstat <base> HEAD -- cmd host docs`** (the directive's filter;
`-- cmd host` alone in its own column). **Every row is measured** from the scratch stages, each committed in a
throwaway repo (`plan-scratch/work2`, tags `M7a` … `M7i`), and re-measured on the round-trip tree (§1.2 C15).

| Sub-milestone | Production | Test | Docs | Own | Cumulative vs `799660ac` (after its commit) | `-- cmd host` | Run |
|---|---|---|---|---|---|---|---|
| M7a — the sprint-3 judge's notes 2–4 (test only) | 0 | 106 | 0 | 106 | **106** | 106 | 1 |
| M7b — `requests.json` `decisions` (`FoldDecisions`, `PendingDecisions`); DMD1–DMD3 | 145 | 186 | 0 | 331 | **437** | 437 | 1 |
| M7c — `approve --list` decisions (Q2) | 57 | 93 | 0 | 150 | **587** | 587 | 1 |
| M7d — *Needs you* lists decisions; the notification (XD3) | 30 | 102 | 0 | 132 | **719** | 719 | 2 |
| M7e — `DecisionSubject.Passkeys` | 33 | 58 | 0 | 91 | **810** | 810 | 2 |
| M7f — the card: two routes, page, `decide.js`; PD1, PD2 | 390 | 150 | 0 | 540 | **1,350** | 1,350 | 2 |
| M7g — live Chrome XD1 (approve and deny with a passkey) | 0 | 285 | 0 | 285 | **1,635** | 1,635 | 3 |
| M7h — live Chrome XD2 (fails closed), XD4 (SSH paste, real `ssh-keygen`) | 0 | 243 | 0 | 243 | **1,878** | 1,878 | 3 |
| M7i — QUICKSTART §18's card paragraph and its test | 0 | 61 | 8 | 69 | **1,947** | 1,939 | 3 |
| **Sprint 4** | **655** | **1,284** | **8** | | **1,947 ≤ 2,000** (margin 53) | 1,939 | |

Per file (+/−): `approve_decide_jgaps4_test.go` +69, `decision_herr_test.go` +37; `surface_approve.go` +6,
`surface_decide_list.go` +80, `surface_decide_list_test.go` +137, `approve_model.go` +3, `decide_list.go` +56,
`decide_list_test.go` +49; `approve.go` +1 −1, `approve_list_decisions.go` +55, `approve_list_decisions_test.go` +93;
`app.js` +25 −5, `app_decide_test.go` +102; `surface_decide.go` +25 −1, `surface_decide_passkeys_test.go` +56,
`surface_decide_subject_test.go` +1 −1, `decide_model.go` +7; `daemon.go` +2, `surface_app_test.go` +4 −2,
`surface_decide_page.go` +29, `surface_decide_page_test.go` +140, `surface_decide_routes_test.go` +2 −2, `decide.go`
+20, `decide/decide.js` +300, `decide/index.html` +39; `decide_page_chrome_test.go` +86, `decide_page_rig_test.go`
+199; `decide_page_closed_test.go` +137, `decide_page_sshsig_test.go` +106; `quickstart_decide_card_test.go` +61,
`QUICKSTART.md` +8. Against Appendix C's ≈ 1,250: the listing + `--list` came in at 481 (priced ≈ 520), *Needs you*
at 132 (≈ 200), `Passkeys` at 91 (≈ 30), the card with its live-Chrome proof at 1,068 (≈ 660: `decide.js` 300 vs
≈ 300, but the three Chrome tests and their rig 528 vs ≈ 300), QUICKSTART at 69 (≈ 20); plus the judge's notes, 106.
**If a later re-measure exceeded 2,000**, the least operator-visible part goes first: M7a (106, test-only notes),
then M7c (150, `--list`; *Needs you* is the operator's list); never PD1 apart from the page (S4-D1).

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of four
backticks, in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; a new file's patch has `--- /dev/null`
and an `@@ -0,0 …` hunk. Apply it **only** this way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF2'
set -euo pipefail
M=M7a   # <- the sub-milestone, exactly as in the Appendix P heading: M7a M7b M7c M7d M7e M7f M7g M7h M7i
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter392-a4-s4
PLAN="$W/design_docs/planned/w-world-approvals-by-key-a4-sprint4-plan.md"
T="$HOME/.ailang/state/world-iter392/exec-patches"; mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
patch -p1 --dry-run < "$T/$M.patch"
patch -p1 < "$T/$M.patch"
if find cmd host docs \( -name '*.orig' -o -name '*.rej' \) | grep .; then echo "STOP: patch left .orig/.rej"; exit 1; fi
git status --short
EOF2
```

- If `--dry-run` fails, **stop and report** (no hand edits, no fuzz options). It can only fail if an earlier
  sub-milestone was skipped or changed.
- Then run that sub-milestone's **Acceptance** (§3) in one bash shell with `export
  AILANG_BIN=$HOME/.pinned-ailang/ailang`, each command's output `tee`'d to a log under
  `$HOME/.ailang/state/world-iter392/exec-logs/` (never only through `tail`). **Never `nohup`, never background a
  test leg.** The planner's `tools/accept.sh <tree> <M>` is exactly §3's lines.
- Then snapshot (§7). The controller commits and compares the numstat with §3's expected value.

---

## 3. Milestones

Every sub-milestone's acceptance starts with the same three lines (abbreviated **[base]** below), in one bash shell
with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` and `export WORLD_EXEC_NODE=/opt/homebrew/bin/node`:

```
test -z "$(gofmt -l cmd host)"
go build ./... && go vet ./host/approval/... ./host/surfacepage/ ./host/daemon/ ./cmd/ailang-worldd/
go test -run '^$' ./... 2>&1 | grep -E '^(FAIL|--- FAIL)' && echo FENCE_RED   # must print nothing
```

Each `-v` run below must show `--- PASS` for every test it names and **no `--- FAIL` and no `--- SKIP`**; each
package run must end `ok` with no `FAIL` line. The planner's `tools/accept.sh <tree> <M>` is exactly these lines.

### M7a — the sprint-3 judge's notes 2–4 · 106 test · Run 1

**Patch** `M7a`: new `cmd/ailang-worldd/approve_decide_jgaps4_test.go` — `TestApproveDecideArtifactPostsTheSubjectsDigests`
(note 3: the `--signature` and the `--webauthn-assertion` bodies carry the subject's P and Q, verb and v; the fake
daemon accepts any packet, so only this pins them), `TestApproveDecideArtifactReadBounds` (note 4: a file one byte
over `approveMaxSignatureBytes` / `approval.MaxPasskeyResponseBytes` is refused naming the bound, before any POST;
a signature of exactly the bound is read and posted). New `host/approval/decision_herr_test.go` —
`TestDecideByAssertionRefusesAMissingMessageAtW6` (note 2: with no decision message (`Verb = 3`) an earlier failing
W1, W2 or W5 refuses with its own error, and W6 alone refuses `ErrChallenge`). Note 1 needs no test (§1.4 P4).

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApproveDecideArtifactPostsTheSubjectsDigests|TestApproveDecideArtifactReadBounds|TestApproveDecideWriteMessageAndSignature|TestApproveDecideWebAuthnLinkAndAssertion)$' -count=1 -v
go test ./host/approval/ -run '^(TestDecideByAssertionRefusesAMissingMessageAtW6|TestDecideByAssertionSteps|TestSignCountLaw|TestApprovalImportsOnlyTheStandardLibrary)$' -count=1 -v
go test ./host/approval/... -count=1
```
**Expected cumulative numstat after the controller's commit of M7a: 106**.

### M7b — `requests.json` `decisions` · 145 prod + 186 test · Run 1

**Patch** `M7b`: new `host/surfacepage/decide_list.go` — `PendingDecision` (the nine members of §2.8, exactly),
`DecisionCandidate`, the pure `FoldDecisions(cands, chainTruncated, scanTruncated)` (pending only, `Human.Approve`
only, ≥ 1 packet; packet and deadline only for exactly one; reserved and refusal passed through; truncated = either
bound). New `host/daemon/surface_decide_list.go` — `PendingDecisions(ctx, *store.Store, now)` (S4-D4: one
`RecentApprovals(40)`, one scan of `decisionListScan` = 101, `GetObject` + decode per scanned packet, `ValidatePacket`
per one-packet request, `ReservedApprovalScope`, `SubjectRefusal`) and `readDecisionList`. `approve_model.go`:
`RequestsJSON` gains `Decisions []PendingDecision \`json:"decisions"\`` and `DecisionsTruncated bool
\`json:"decisionsTruncated"\``. `surface_approve.go`: `handleSurfaceRequests` calls `readDecisionList` (an error is
logged; the list stays `[]`, never null). Tests: `host/surfacepage/decide_list_test.go` (DMD1–DMD3:
`TestDecideModelListsPendingDecisionsOnly`, `TestDecideModelListsATwoPacketRequestRefused`,
`TestDecideModelTruncation`), `host/daemon/surface_decide_list_test.go` (through `d.Handler()`:
`TestSurfaceRequestsListsDecisions` — an empty world `[]`; a decidable row exactly (packet, deadline, cost); a
reserved publish row with D5; a two-packet row with D4b and no packet; an uncommitted packet's refusal; decided and
packet-less requests absent; past the deadline D7; `TestSurfaceRequestsDecisionsTruncated` — 101 packets truncated,
100 not, 21 chain requests truncated; `TestSurfaceRequestsDecisionsUnreadable` — a closed world store still answers
200 with `decisions: []` and the mint lists).

**Acceptance**
```
[base]
go test ./host/surfacepage/ -run '^(TestDecideModelListsPendingDecisionsOnly|TestDecideModelListsATwoPacketRequestRefused|TestDecideModelTruncation)$' -count=1 -v
go test ./host/daemon/ -run '^(TestSurfaceRequestsListsDecisions|TestSurfaceRequestsDecisionsTruncated|TestSurfaceRequestsDecisionsUnreadable|TestSurfaceAppRoutes|TestFrozenV1RouteTableMatchesMux|TestDaemonDependencyAllowlist)$' -count=1 -v
go test ./host/surfacepage/ -count=1
go test ./host/daemon/ -run '^(TestDecision|TestDecide|TestSurface|TestApprov)' -count=1
```
**Expected cumulative numstat after the controller's commit of M7b: 437**.

### M7c — `approve --list` decisions · 57 prod + 93 test · Run 1

**Patch** `M7c`: new `cmd/ailang-worldd/approve_list_decisions.go` — `approveListDecisions` (S4-D3: `store.OpenReadOnly(--db)`,
`daemon.PendingDecisions`, never a dial; an unreadable store is one stderr line, stdout and the exit code unchanged)
and `printDecisions` (nothing for none; `pending decisions (Human.Approve requests):`; per row `<ref>  <scope>
requested by <r>` then `refused: <reason>` or `deadline <RFC3339>  decide: ailang-worldd approve <ref> [--deny]`; a
truncation line). `approve.go` (+1 −1): the `--list` branch no longer returns early on "no pending mint requests"
and calls it after the mint lines. New `approve_list_decisions_test.go` — `TestApproveListShowsDecisions` (a temp
`--db` with a refused request: the exact stdout, empty stderr; a missing `--db` world store: one stderr line, rc
0) and `TestApproveListDecisionLines` (empty, decidable, refused, truncated, exact text). The existing exact-output
`--list` tests (`TestApproveListUsesTheNamedStores` incl. its temp-`HOME` default-db arm, `TestApproveListShowsPending`)
run unedited.

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApproveListShowsDecisions|TestApproveListDecisionLines|TestApproveListUsesTheNamedStores|TestApproveListShowsPending|TestEveryVerbHasHelp)$' -count=1 -v
go test ./host/store/ -run '^(TestProductionContextRoots|TestNoNewDeadlineFreeStoreReads)$' -count=1 -v
go test ./cmd/ailang-worldd/ -run '^(TestApprove|TestApprover)' -count=1
```
**Expected cumulative numstat after the controller's commit of M7c: 587**.

### M7d — *Needs you* lists decisions (XD3) · 30 prod + 102 test · Run 2

**Patch** `M7d`: `host/surfacepage/app/app.js` (+25 −5) — `needsYouHTML` adds one line per decision after the mint
requests (S4-D5: decidable link, reserved owner text, refused text, truncation line; the count includes them; every
value through `esc()`), `notifications` raises "World needs you" for a new decidable decision id (href
`/surface/decide/<ref>`, tag the ref; reserved and refused are seen, never raised; nothing on the first read). New
`host/surfacepage/app_decide_test.go` — `TestNeedsYouListsDecisions` (node, `TestSurfaceAppPure`'s harness rule:
the exact decidable line and link, the count, order after a mint request, reserved publish/finish and refused lines
with no `href`, exactly two `<a `, the fallback owner, the truncation line with and without rows, hostile values in
ref/scope/requester/refusal escaped, a `requests.json` without `decisions` byte-identical to before; the
notification arms).

**Acceptance**
```
[base]
WORLD_EXEC_NODE=/opt/homebrew/bin/node go test ./host/surfacepage/ -run '^(TestNeedsYouListsDecisions|TestSurfaceAppPure|TestSurfaceScriptPure|TestSurfaceScriptAssetSeams)$' -count=1 -v
WORLD_EXEC_NODE=/opt/homebrew/bin/node go test ./host/surfacepage/ -count=1
```
**Expected cumulative numstat after the controller's commit of M7d: 719**.

### M7e — `DecisionSubject.Passkeys` · 33 prod + 58 test · Run 2

**Patch** `M7e`: `decide_model.go` — `Passkeys []CardPasskeys \`json:"passkeys"\`` and `DecisionPasskeys` (S4-D6);
`surface_decide.go` — `decisionPasskeys` on a fresh read-only credentials handle (below schema 6: none; other
errors: 500) filled into the subject. New `surface_decide_passkeys_test.go` — `TestDecisionSubjectPasskeys` (none
enrolled `[]`; two origins grouped with their rp; a revoked passkey excluded; an unopenable store 500). Edited
`surface_decide_subject_test.go` +1 −1 (P7).

**Acceptance**
```
[base]
go test ./host/daemon/ -run '^(TestDecisionSubjectPasskeys|TestDecisionSubjectJSON|TestDaemonCredentialWriterSites)$' -count=1 -v
go test ./host/daemon/ -run '^(TestDecision|TestDecide|TestSurfaceApprove|TestApprov)' -count=1
```
**Expected cumulative numstat after the controller's commit of M7e: 810**.

### M7f — the decision card (PD1, PD2) · 390 prod + 150 test · Run 2

**Patch** `M7f`: `daemon.go` +2 — `GET /surface/decide/{ref}` and `GET /surface/decide.js`, both
`d.approvalRoute(…)`; new `surface_decide_page.go` (the page: 400 unless D1's grammar, the same bytes for every
ref; the script; `setSurfaceHeaders`, as `surface_approve_page.go`); new `host/surfacepage/decide.go` (embeds) and
`decide/index.html`, `decide/decide.js` (design §2.8 steps 1–5: the subject fetch with a 10 s abort; both hashes
and the packet's `requestRef` checked from the served bytes, fail closed; `DECISION_FIELDS` renders the twelve
fields + entry + `/workbench?object=<proposal>`; a refused, reserved or decided subject shows why and no control;
**Approve with passkey** / **Deny with passkey** — challenge = SHA-256 of the 91-byte message with the `BigInt`
expiry `window.max`, re-read once when < 60 s remain, UV required, `allowCredentials` of this origin's group, A2's
240 s abort; **Save the approve/deny message** (`blob:` `world-decision-<verb>-<ref16>.msg`), the `ssh-keygen -Y
sign -n world-decision@sunholo` line and a paste box posting `sshsig`; POST with `credentials: 'omit'`; the
verb on every button and status; `decidedBy` or the refusal class and reason, then a re-read; every value by
`textContent`). New `surface_decide_page_test.go` — PD1 `TestDecidePageFieldTable` (S4-D7), PD2
`TestDecidePageAssets` (both routes' bytes, `surfaceCSP` unchanged, nosniff, no-store, no-referrer, the content
types, no cookie and no CORS header; a foreign Host 421; four malformed refs 400; the A2/A3 page hashes pinned; the
page loads only its own script). Edited `surface_app_test.go` (+4 −2, design §7 +2) and `surface_decide_routes_test.go`
(+2 −2, P8).

**Acceptance**
```
[base]
go test ./host/daemon/ -run '^(TestDecidePageFieldTable|TestDecidePageAssets|TestSurfaceAppRoutes|TestDecisionRoutesRegistered|TestFrozenV1RouteTableMatchesMux|TestApprovePageRoutes|TestPasskeyPageRoutes|TestGETRoutesRejectOtherMethods|TestApprovalRoutesRefuseForeignHost|TestDaemonDependencyAllowlist)$' -count=1 -v
go test ./host/surfacepage/ -count=1
go test ./host/daemon/ -run '^(TestDecision|TestDecide|TestSurface|TestApprov|TestPasskey|TestGETRoutes)' -count=1
```
**Expected cumulative numstat after the controller's commit of M7f: 1,350**.

### M7g — the card in live Chrome: XD1 · 285 test · Run 3

**Patch** `M7g`: new `cmd/ailang-worldd/decide_page_rig_test.go` (S4-D11: the real card bytes from `d.Handler()`;
the subject and POST with A3a's guards, the real grammar, law and judges over a temp credentials store; Chrome by
the A2 driver with a temp profile and a virtual authenticator with UV) and `decide_page_chrome_test.go` —
`TestDecidePageApproveWithPasskey` (XD1: Approve on one card, Deny on a second, Approve with 50 s left (must
re-read); each `get`'s challenge = Go's `DecisionMessageHash(verb, P, Q, window.max)`, rpId, allow list, UV, 240 s,
a signal; each POST's verb, P, Q, expiry, `webauthn` only; the decision `webauthn:mark-mac`; 13 fields, dates, the
proposal link, the signing lines, **the verb texts**; re-read to `decided` with no control; no CSP violation, no page
error). Without `WORLD_CHROME` it skips, counted; with `WORLD_REQUIRE_PASSKEY_BROWSER=1` a missing browser fails.

**Acceptance** (gated as XB1, on the controller's machine; CI has no Chrome step)
```
[base]
WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" WORLD_REQUIRE_PASSKEY_BROWSER=1 go test ./cmd/ailang-worldd/ -run '^(TestDecidePageApproveWithPasskey|TestApproveDecideWebAuthnLinkInChrome|TestPasskeyPageInChrome)$' -count=1 -v
```
(three `--- PASS`.) **Expected cumulative numstat after the controller's commit of M7g: 1,635**.

### M7h — the card in live Chrome: XD2, XD4 · 243 test · Run 3

**Patch** `M7h`: new `decide_page_closed_test.go` — `TestDecidePageFailsClosed` (XD2, 16 cases: swapped request
bytes, a wrong packet hash, a subject naming another request, a packet carrying another request, unlisted keys in
either object, a missing key, string-for-int, a fraction, int-for-string, a non-object, bad base64url, `window.max`
−1 → refused, no field, no control, **no `get`**; reserved, refused and decided subjects render and offer no
control; a control case is ready; no POST). New `decide_page_sshsig_test.go` — `TestDecidePageSSHSIGPaste` (XD4:
both saved files equal `approval.DecisionMessage` (91 bytes) as `blob:` downloads named
`world-decision-<verb>-<ref16>.msg`; an approvaltest SSHSIG posted as the other verb 422 and decides nothing, as its
own verb decides `ssh:mark-mac`; a passkey of another origin is named and the passkey buttons hidden; **real-tool
sibling** (directive retro 11): `ssh-keygen -Y sign` on a temp-dir ed25519 key over the saved deny file decides
`ssh:mark-key`; `WORLD_REQUIRE_SSH_KEYGEN=1` makes a missing `ssh-keygen` a failure).

**Acceptance**
```
[base]
WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" WORLD_REQUIRE_PASSKEY_BROWSER=1 WORLD_REQUIRE_SSH_KEYGEN=1 go test ./cmd/ailang-worldd/ -run '^(TestDecidePageFailsClosed|TestDecidePageSSHSIGPaste|TestApprovePageInChrome|TestApprovePageRefusesASwappedRequest)$' -count=1 -v
```
(four `--- PASS`.) **Expected cumulative numstat after the controller's commit of M7h: 1,878**.

### M7i — QUICKSTART §18's card paragraph · 8 docs + 61 test · Run 3

**Patch** `M7i`: `docs/QUICKSTART.md` — §18 gains **From the surface** (*Needs you*'s line, the `approve --list`
line, the card at `/surface/decide/<request-ref>`, its six controls, what it refuses; no new `approve sha256:` line,
so `TestQuickstartSection18FlagsMatchTheCLI`'s count of six holds unedited). New `quickstart_decide_card_test.go` —
`TestQuickstartSection18CardMatchesThePage` (every bold control the paragraph names is a control of
`decide/index.html`; the route is the daemon's; the one `approve --list` line's flags are in `approveHelp`).

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestQuickstartSection18CardMatchesThePage|TestQuickstartSection18FlagsMatchTheCLI|TestEveryVerbHasHelp|TestQuickstartAttendedStepsBoxMatchesSection9|TestQuickstartSection9FlagsMatchTheCLI|TestQuickstartSection10FlagsMatchTheCLI|TestQuickstartSection11FlagsMatchTheCLI)$' -count=1 -v
go test ./cmd/ailang-worldd/ -run '^(TestQuickstart|TestApproveDecide|TestApproveList)' -count=1
```
**Expected cumulative numstat after the controller's commit of M7i: 1,947** (`-- cmd host`: 1,939).

### 3.1 Reachability: every new refusal and branch, from the entry point (directive retro 4)

Entry points: **LIST** = `GET /surface/requests.json` through `d.Handler()` (`approvalRoute`, `handleSurfaceRequests`,
`readDecisionList`, `PendingDecisions`, `FoldDecisions`); **SUBJ** = `GET /surface/decisions/{ref}`; **PAGE** =
`GET /surface/decide/{ref}` and `/surface/decide.js`; **CARD** = the page's script in Chrome (subject → checks →
render → ceremony/paste → POST); **APP** = `app.js` (`needsYouHTML`, `notifications`) under node; **CLI** =
`runApprove --list`. "E2E" = driven through that entry point; "unit" = a direct call only.

| Branch (file) | Input class → outcome | Trace from the entry point | Driven |
|---|---|---|---|
| fold: decided / not `Human.Approve` / no packet → dropped (`decide_list.go`) | each → absent | LIST → `PendingDecisions` → `FoldDecisions` | E2E `TestSurfaceRequestsListsDecisions` (decided, packet-less); unit `TestDecideModelListsPendingDecisionsOnly` (all three; the effect arm unit only: the daemon lists the chain's `Human.Approve` summaries, and a non-`Human.Approve` request on the chain is made by other effects' tests) |
| fold: exactly one packet → packet + deadline; ≥ 2 → neither | one / two | LIST → fold | E2E `TestSurfaceRequestsListsDecisions` (one; two sharing a ref); unit `…TwoPacket…` |
| truncation: chain bound / scan bound | 21 chain requests / 101 packets → `true`; 100 → `false` | LIST → `RecentApprovals` / scan | E2E `TestSurfaceRequestsDecisionsTruncated`; unit `TestDecideModelTruncation` |
| `ValidatePacket` code ≠ 0 → refusal sentence | an uncommitted packet → "is refused: …" | LIST → `PendingDecisions` | E2E `TestSurfaceRequestsListsDecisions` |
| reserved → `reserved` + D5 | a publish-shaped scope | LIST | E2E same |
| deadline passed → D7 | `now ≥ deadline` | LIST | E2E same |
| listing error → logged, `[]` | a closed world store | LIST → handler | E2E `TestSurfaceRequestsDecisionsUnreadable` |
| `ValidatePacket` error arm | a store failure after the scan | LIST | **not reachable**: an earlier read fails first (survivor `L-VALIDATE-ERR`) |
| `--list` decisions section; unreadable `--db` → stderr line | a temp `--db` / a missing one | CLI | E2E `TestApproveListShowsDecisions`; text unit `TestApproveListDecisionLines` |
| `--list` `PendingDecisions` error arm | a malformed approvals chain in a store that opens | CLI | **not reachable** without a write outside the broker (survivor `CL-LIST-ERR`) |
| *Needs you*: decidable / reserved / refused / truncated lines; count; escaping | each | APP | E2E (node) `TestNeedsYouListsDecisions` |
| notification: new decidable only | first read / seen / reserved / refused / not granted | APP | E2E (node) `TestNeedsYouListsDecisions` |
| `Passkeys`: none / grouped / revoked / store error 500 | each | SUBJ → `decisionPasskeys` | E2E `TestDecisionSubjectPasskeys` |
| `Passkeys`: below schema 6 → none; non-older `ActivePasskeys` error | — | SUBJ | **not reachable**: `serve` upgrades the store (survivors `P-PASSKEYS-OLDER`, `P-PASSKEYS-ACTIVE-ERR`) |
| page route: ref not D1 → 400; foreign Host → 421 | 4 malformed refs; a foreign Host | PAGE | E2E `TestDecidePageAssets` |
| card: request hash / packet hash / `requestRef` / subject ref mismatch → refused, no `get` | swapped bytes etc. | CARD | E2E live Chrome `TestDecidePageFailsClosed` |
| card: unlisted key / missing key / wrong type / unsafe int / fraction / non-object / bad base64url / bad window → refused, no `get` | each | CARD | E2E live Chrome `TestDecidePageFailsClosed` |
| card: refused / reserved / decided subject → fields, reason, no control | each | CARD | E2E live Chrome `TestDecidePageFailsClosed`, `TestDecidePageApproveWithPasskey` (the decided re-read) |
| card: passkey ceremony (verb, P, Q, `BigInt` expiry, allow list, UV) → POST `webauthn` | approve / deny | CARD → rig POST (S4-D11) | E2E live Chrome `TestDecidePageApproveWithPasskey` |
| card: < 60 s left → one re-read before signing | 50 s left | CARD | E2E live Chrome same |
| card: no passkey for this origin → buttons hidden, named | a passkey of another origin | CARD | E2E live Chrome `TestDecidePageSSHSIGPaste` |
| card: Save the message (`blob:`, name) + paste → POST `sshsig` | approve / deny / mismatched verb (422) | CARD → rig POST | E2E live Chrome `TestDecidePageSSHSIGPaste` (approvaltest + real `ssh-keygen`) |
| card: 240 s abort | a ceremony left pending | CARD | **not driven** for `decide.js` (survivor `J-TIMEOUT`; A2's own test drives the same rule on `passkey.js`) |
| card: `credentials: 'omit'` | — | CARD | **not observable** with a cookie-less temp profile (survivor `J-CREDS`; R0 refuses a cookie at the daemon) |

**No page or POST arm ships without an operator path**: *Needs you* (M7d) links the card (M7f), whose controls post
to the unchanged sprint 1–3 arms; `approve --list` (M7c) prints the ref and the decide command for the terminal;
the CLI arms of sprints 1–3 are unchanged.

---

## 4. Mutation table

**Procedure** (`tools/mutate.py <run> <tree> <workers> <rows module> [ids…]`, rows in `tools/mutants_{J,A,C,B,B_js}.py`):
for each row, an APFS clone of the tree (`cp -cR`), one exact-string replacement whose anchor occurs exactly once
in the file (else `ANCHOR<n>`, not a run), then `go test -count=1` of the **named test** and of the **package suite**:
a fast regex first (`host/daemon`: `^(TestDecision|TestDecide|TestSurface|TestApprov|TestPasskey|TestIs)`;
`cmd/ailang-worldd`: `^(TestApprove|TestApprover|TestDecide|TestNeedsYou|TestPasskey)`; `host/surfacepage`,
`host/approval`: the whole package), and **a row the fast suite does not kill is re-run on the whole package**
(`^Test`, verdict suffixed `(whole)`) before it may be called SURVIVED (directive retro 3: 391's narrow `-run`
false survivor). Inventory and static tests are excluded from every suite (`TestDecisionRoutesRegistered`,
`TestApprovalRoutesRegistered`, `TestDaemonCredentialWriterSites`, `TestSurfaceAppRoutes`,
`TestFrozenV1RouteTableMatchesMux`, `TestDaemonDependencyAllowlist`, `TestGETRoutesRejectOtherMethods`,
`TestEveryVerbHasHelp`, `TestQuickstart*`, `TestApprovalImportsOnlyTheStandardLibrary`), so **no kill rests on an
inventory**; the PD2 page-asset checks that kill the route rows (`P-REF`, `P-PAGE-BYTES`, `P-SCRIPT-CT`,
`P-ROUTE-*`) are behavioural HTTP checks of the served bytes, status and headers, not a static scan. Go rows ran
with `GOFLAGS=-vet=off` and `WORLD_CHROME` unset; the `decide.js` rows (`mutants_B_js.py`) ran with
`WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" WORLD_REQUIRE_PASSKEY_BROWSER=1` and
skip `TestMarketingChrome*` (§1.2 C13): a JS mutant can only be killed by a live-Chrome test. A build failure is
`NOCOMPILE` and is **not a mutant**: the row is rewritten to compile and re-run. Every run appends one line per row
to `plan-scratch/logs/mutant_ledger.tsv` (time, run, tree, id, stage, file, from → to, named regex, named result,
suite result, failing tests, seconds, expected). **The ids below are the ledger's ids** (`comm -3 <(cut -f4 ledger
| sort -u) logs/plan_ids.txt` → empty). Expected: `K` killed by both; `E` a named equivalent (reason given); `L` a
named survivor (unreachable or unobservable, reason given).

**Trees** (directive retro 3): runs J1–J2 on `wJ`, A1 on `wA`, C1–C2 on `wC`, B1–B5js on `wB` — each = the base +
that one slice, whose files are **byte-identical** to the final tree's (they were copied from them, `diff -r` of
`work2` at `M7i` against `work`: `WORK2_EQUALS_FINAL`); the final tree only adds other slices' tests, which can only
add kills. **Every `E` and `L` row was re-run on the final tree `work`** (runs F1, F2, F2js), so no survivor verdict
rests on a slice tree. **How to read an `E`/`L` row whose suite column says `KILLED(whole)`:** the fast suite — every test that reads the
mutated code — SURVIVED, and the whole-package re-run was red **only on untouched tests that are red on the
unmutated final tree under the same load** (named in the row): `TestExecStartupRefusalTable`,
`TestDaemonInvokeConfig`, `TestNewBootstrapsEpochRegistryIdempotently` (an archived `ailang --version` timing out
after 10 s on a freshly copied binary — the final gate's own `pkgs` reds, §1.2 C18, all green alone),
`TestRegistryPhaseHasItsOwnClock` and `TestServeStartsWhenTheProbeSpendsTheStartupClock` (the directive's
known-flaky list; the latter red at base), and `TestWorldContextAndStatusline` (a subprocess killed at its deadline
under load, red at the base's `verify_go`). None of them reaches the listing, the card or `--list`; such a line is
not a kill, and the row's verdict is its reason. `L-RA-ERR`, `CL-LIST-ERR`, `J-SUBJECT-REF`, `J-TIMEOUT` and `J-CREDS` SURVIVED the whole package outright.

**126 rows, the last ledger line of each: 113 `K` (killed by the named test and by the package suite), 7 `E` (named equivalents, reason given), 6 `L` (named survivors, reason given); 0 not as expected.** 148 ledger lines in all (`plan-scratch/logs/mutant_ledger.tsv`; runs: A1 33, B1 13, B1js 20, B2js 18, B3js 3, B4 3, B5js 4, C1 33, C2 2, F1 7, F2 2, F2js 4, J1 5, J2 1); 1 ledger lines are `NOCOMPILE` and not counted (each row was rewritten to compile and re-run).

| Id | Stage | file: from → to | Named test | Run/tree | Named | Suite | Failing tests (suite, else named) | Expected |
|---|---|---|---|---|---|---|---|---|
| JG-BODY-PACKET | M7a | `approve_decide.go`: verb, _ := approval.DecisionVerb(o.verb)⏎ body := surfacepage.DecideBody{V: 1, Verb: o.verb, Packet: strings.TrimPrefix(s.Packet, "sha256:"), Request: strings.TrimPrefix(o.ref, "sha256:"), Expiry: s.Window.Max}⏎ if o.… | `TestApproveDecideArtifactPostsTheSubjectsDigests` | J1/wJ | KILLED | KILLED | TestApproveDecideArtifactPostsTheSubjectsDigests | K |
| JG-BODY-REQUEST | M7a | `approve_decide.go`: verb, _ := approval.DecisionVerb(o.verb)⏎ body := surfacepage.DecideBody{V: 1, Verb: o.verb, Packet: strings.TrimPrefix(s.Packet, "sha256:"), Request: strings.TrimPrefix(o.ref, "sha256:"), Expiry: s.Window.Max}⏎ if o.… | `TestApproveDecideArtifactPostsTheSubjectsDigests` | J1/wJ | KILLED | KILLED | TestApproveDecideArtifactPostsTheSubjectsDigests,TestApproveDecideArtifactReadBounds,TestApproveDecideWebAuthnLinkAndAss | K |
| JG-SIG-CAP | M7a | `approve_decide.go`: readRegularFileCapped(o.sigPath, approveMaxSignatureBytes) → readRegularFileCapped(o.sigPath, 1<<30) | `TestApproveDecideArtifactReadBounds` | J1/wJ | KILLED | KILLED | TestApproveDecideArtifactReadBounds | K |
| JG-ASR-CAP | M7a | `approve_decide.go`: readRegularFileCapped(o.assertion, approval.MaxPasskeyResponseBytes) → readRegularFileCapped(o.assertion, 1<<30) | `TestApproveDecideArtifactReadBounds` | J1/wJ | KILLED | KILLED | TestApproveDecideArtifactReadBounds | K |
| JG-HERR-TOP | M7a | `passkey_assert.go`:  p := a.Passkey⏎ resp, err := ParseAssertionResponse(response) // W1 →  p := a.Passkey⏎ if hErr != nil {⏎  return AuthData{}, ErrChallenge⏎ }⏎ resp, err := ParseAssertionResponse(response) // W1 | `TestDecideByAssertionRefusesAMissingMessageAtW6` | J2/wJ | KILLED | KILLED | TestDecideByAssertionRefusesAMissingMessageAtW6 | K |
| JG-W15-OPERAND | M7a | `passkey_assert.go`: if (ad.SignCount != 0 \|\| p.SignCount != 0) && ad.SignCount <= p.SignCount { // W15 → if ad.SignCount != 0 && ad.SignCount <= p.SignCount { // W15 | `TestSignCountLaw` | J1/wJ | KILLED | KILLED | TestSignCountLaw | K |
| L-DECIDED | M7b | `decide_list.go`: if c.Decided \|\| c.Effect != decisionEffect \|\| len(c.Packets) == 0 { → if c.Effect != decisionEffect \|\| len(c.Packets) == 0 { | `TestDecideModel` | A1/wA | KILLED | KILLED | TestDecideModelListsPendingDecisionsOnly | K |
| L-EFFECT | M7b | `decide_list.go`: if c.Decided \|\| c.Effect != decisionEffect \|\| len(c.Packets) == 0 { → if c.Decided \|\| len(c.Packets) == 0 { | `TestDecideModel` | A1/wA | KILLED | KILLED | TestDecideModelListsPendingDecisionsOnly | K |
| L-ZERO | M7b | `decide_list.go`: if c.Decided \|\| c.Effect != decisionEffect \|\| len(c.Packets) == 0 { → if c.Decided \|\| c.Effect != decisionEffect { | `TestDecideModel` | A1/wA | KILLED | KILLED | TestDecideModelListsPendingDecisionsOnly | K |
| L-ONE | M7b | `decide_list.go`:   if len(c.Packets) == 1 { →   if len(c.Packets) >= 1 { | `TestDecideModel` | A1/wA | KILLED | KILLED | TestDecideModelListsATwoPacketRequestRefused | K |
| L-RESERVED | M7b | `decide_list.go`: Reserved: c.Reserved, Refusal: c.Refusal} → Refusal: c.Refusal} | `TestDecideModel` | A1/wA | KILLED | KILLED | TestDecideModelListsPendingDecisionsOnly | K |
| L-REFUSAL | M7b | `decide_list.go`: Reserved: c.Reserved, Refusal: c.Refusal} → Reserved: c.Reserved} | `TestDecideModel` | A1/wA | KILLED | KILLED | TestDecideModelListsATwoPacketRequestRefused,TestDecideModelListsPendingDecisionsOnly | K |
| L-TRUNC-CHAIN | M7b | `decide_list.go`: return rows, chainTruncated \|\| scanTruncated → return rows, scanTruncated | `TestDecideModel` | A1/wA | KILLED | KILLED | TestDecideModelTruncation | K |
| L-TRUNC-SCAN | M7b | `decide_list.go`: return rows, chainTruncated \|\| scanTruncated → return rows, chainTruncated | `TestDecideModel` | A1/wA | KILLED | KILLED | TestDecideModelTruncation | K |
| L-SCAN-SENTINEL | M7b | `surface_decide_list.go`: scanTruncated := len(scanned) == decisionListScan → scanTruncated := false | `TestSurfaceRequests` | A1/wA | KILLED | KILLED | TestSurfaceRequestsDecisionsTruncated,TestSurfaceRequestsDecisionsTruncated/packet-scan | K |
| L-SCAN-BOUND | M7b | `surface_decide_list.go`: const decisionListScan = 101 → const decisionListScan = 100 | `TestSurfaceRequests` | A1/wA | KILLED | KILLED | TestSurfaceRequestsDecisionsTruncated,TestSurfaceRequestsDecisionsTruncated/packet-scan-at-the-bound | K |
| L-CHAIN-TRUNC | M7b | `surface_decide_list.go`: surfacepage.FoldDecisions(cands, page.Truncated, scanTruncated) → surfacepage.FoldDecisions(cands, false, scanTruncated) | `TestSurfaceRequests` | A1/wA | KILLED | KILLED | TestSurfaceRequestsDecisionsTruncated,TestSurfaceRequestsDecisionsTruncated/approvals-chain | K |
| L-SCAN-DECODE | M7b | `surface_decide_list.go`: if p, err := coordinator.DecodeDecisionPacket(full.Payload); ok && err == nil { → if p, _ := coordinator.DecodeDecisionPacket(full.Payload); ok { | `TestSurfaceRequests` | F1/work | SURVIVED | KILLED(whole) | TestDaemonInvokeConfig,TestExecStartupRefusalTable,TestExecStartupRefusalTable/node-in-root,TestNewBootstrapsEpochRegist | E: an undecodable scanned packet has a zero `RequestRef`, which matches no chain request; dropping the decode guard lists nothing new |
| L-STATUS | M7b | `surface_decide_list.go`: Cost: sm.Cost, Decided: sm.Status != "pending"} → Cost: sm.Cost, Decided: false} | `TestSurfaceRequests` | A1/wA | KILLED | KILLED | TestSurfaceRequestsListsDecisions | K |
| L-GUARD | M7b | `surface_decide_list.go`: if !c.Decided && sm.Effect == broker.EffectHumanApprove && len(packets) != 0 { → if len(packets) != 0 { | `TestSurfaceRequests` | F1/work | SURVIVED | KILLED(whole) | TestExecStartupRefusalTable,TestExecStartupRefusalTable/node-in-state,TestNewBootstrapsEpochRegistryIdempotently | E: an efficiency guard: the fold (`FoldDecisions`) drops decided and non-`Human.Approve` rows itself |
| L-ONE-VALIDATE | M7b | `surface_decide_list.go`:    if len(packets) == 1 { →    if len(packets) >= 1 { | `TestSurfaceRequests` | F1/work | SURVIVED | KILLED(whole) | TestDaemonInvokeConfig,TestExecStartupRefusalTable,TestExecStartupRefusalTable/node-in-root,TestNewBootstrapsEpochRegist | E: with two or more packets D4b precedes `PacketCode` in `SubjectRefusal`, and the fold drops the deadline |
| L-CODE | M7b | `surface_decide_list.go`: f.PacketCode, f.PacketSentence, f.Deadline = res.Code, packetRefusalSentence(res), res.Packet.DeadlineAt → f.PacketCode, f.PacketSentence, f.Deadline = 0, packetRefusalSentence(res), res.Packet.DeadlineAt | `TestSurfaceRequests` | A1/wA | KILLED | KILLED | TestSurfaceRequestsListsDecisions | K |
| L-FACT-DEADLINE | M7b | `surface_decide_list.go`: f.PacketCode, f.PacketSentence, f.Deadline = res.Code, packetRefusalSentence(res), res.Packet.DeadlineAt → f.PacketCode, f.PacketSentence, f.Deadline = res.Code, packetRefusalSentence(res), now+600 | `TestSurfaceRequests` | A1/wA | KILLED | KILLED | TestSurfaceRequestsListsDecisions | K |
| L-ROW-DEADLINE | M7b | `surface_decide_list.go`:     c.Deadline = res.Packet.DeadlineAt →     c.Deadline = 0 | `TestSurfaceRequests` | A1/wA | KILLED | KILLED | TestSurfaceRequestsListsDecisions | K |
| L-RESERVED-D | M7b | `surface_decide_list.go`: c.Reserved, _ = broker.ReservedApprovalScope(sm.Scope) → c.Reserved = "" | `TestSurfaceRequests` | A1/wA | KILLED | KILLED | TestSurfaceRequestsListsDecisions | K |
| L-REFUSAL-D | M7b | `surface_decide_list.go`: c.Refusal = broker.SubjectRefusal(f).Reason → c.Refusal = "" | `TestSurfaceRequests` | A1/wA | KILLED | KILLED | TestSurfaceRequestsListsDecisions | K |
| L-DISCOVERY | M7b | `surface_decide_list.go`: Discovery: broker.Discovery{Complete: true, Scanned: len(scanned), Carrying: len(packets)}} → Discovery: broker.Discovery{Complete: true, Scanned: len(scanned), Carrying: 1}} | `TestSurfaceRequests` | A1/wA | KILLED | KILLED | TestSurfaceRequestsListsDecisions | K |
| L-RA-ERR | M7b | `surface_decide_list.go`: page, err := broker.RecentApprovals(ctx, s, 40)⏎ if err != nil { → page, err := broker.RecentApprovals(ctx, s, 40)⏎ if false && err != nil { | `TestSurfaceRequests` | F1/work | SURVIVED | SURVIVED(whole) | — | E: a failed or malformed chain read and an empty page both list `[]` (the handler logs the error either way) |
| L-SCAN-ERR | M7b | `surface_decide_list.go`: scanned, err := s.ObjectsBySemanticID(ctx, coordinator.DecisionPacketV1, "", decisionListScan)⏎ if err != nil { → scanned, err := s.ObjectsBySemanticID(ctx, coordinator.DecisionPacketV1, "", decisionListScan)⏎ if fals… | `TestSurfaceRequests` | F1/work | SURVIVED | KILLED(whole) | TestRegistryPhaseHasItsOwnClock,TestRegistryPhaseHasItsOwnClock/probe-spends-the-startup-clock,TestServeStartsWhenThePro | E: no packets drops every row, so `[]` either way |
| L-VALIDATE-ERR | M7b | `surface_decide_list.go`: res, err := coordinator.ValidatePacket(ctx, s, packets[0])⏎    if err != nil { → res, err := coordinator.ValidatePacket(ctx, s, packets[0])⏎    if false && err != nil { | `TestSurfaceRequests` | F1/work | SURVIVED | KILLED(whole) | TestRegistryPhaseHasItsOwnClock,TestRegistryPhaseHasItsOwnClock/probe-spends-the-startup-clock,TestServeStartsWhenThePro | L: unreachable: any store failure trips the earlier `RecentApprovals`/scan/`GetObject` reads first |
| L-NIL-DEFAULT | M7b | `surface_approve.go`:  if out.Decisions == nil { →  if false { | `TestSurfaceRequests` | A1/wA | KILLED | KILLED | TestSurfaceRequestsDecisionsUnreadable,TestSurfaceRequestsListsDecisions | K |
| CL-CALL | M7c | `approve.go`:   approveListDecisions(ctx, *db, env.now(), stdout, stderr) // row 201 A4 (design Q2)⏎ →  | `TestApproveList(ShowsDecisions|DecisionLines)` | A1/wA | KILLED | KILLED | TestApproveListShowsDecisions | K |
| CL-OPEN-ERR | M7c | `approve_list_decisions.go`:  ro, err := store.OpenReadOnly(db)⏎ if err != nil { →  ro, err := store.OpenReadOnly(db)⏎ if false && err != nil { | `TestApproveList(ShowsDecisions|DecisionLines)` | A1/wA | KILLED | KILLED | TestApproveListShowsDecisions | K |
| CL-LIST-ERR | M7c | `approve_list_decisions.go`:  _ = ro.Close()⏎ if err != nil { →  _ = ro.Close()⏎ if false && err != nil { | `TestApproveList(ShowsDecisions|DecisionLines)` | F1/work | SURVIVED | SURVIVED(whole) | — | L: unreachable from the CLI: needs a world store that opens but holds a malformed approvals chain (only a write outside the broker makes one); the mutant only drops a stderr line |
| CL-EMPTY-ROWS | M7c | `approve_list_decisions.go`: if len(rows) == 0 && !truncated { → if !truncated { | `TestApproveList(ShowsDecisions|DecisionLines)` | A1/wA | KILLED | KILLED | TestApproveListShowsDecisions | K |
| CL-EMPTY-TRUNC | M7c | `approve_list_decisions.go`: if len(rows) == 0 && !truncated { → if len(rows) == 0 { | `TestApproveList(ShowsDecisions|DecisionLines)` | A1/wA | KILLED | KILLED | TestApproveListDecisionLines | K |
| CL-REFUSAL | M7c | `approve_list_decisions.go`:   if d.Refusal != "" { →   if false { | `TestApproveList(ShowsDecisions|DecisionLines)` | A1/wA | KILLED | KILLED | TestApproveListDecisionLines,TestApproveListShowsDecisions | K |
| CL-DEADLINE | M7c | `approve_list_decisions.go`: approverTime(d.Deadline), d.Request) → approverTime(0), d.Request) | `TestApproveList(ShowsDecisions|DecisionLines)` | A1/wA | KILLED | KILLED | TestApproveListDecisionLines | K |
| CL-TRUNC | M7c | `approve_list_decisions.go`:  if truncated { →  if false { | `TestApproveList(ShowsDecisions|DecisionLines)` | A1/wA | KILLED | KILLED | TestApproveListDecisionLines | K |
| N-DS-ARRAY | M7d | `app.js`: const ds = r && Array.isArray(r.decisions) ? r.decisions : []; → const ds = []; | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-COUNT | M7d | `app.js`: const n = p.length + ds.length; → const n = p.length; | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-HEAD-COLOR | M7d | `app.js`: style="color:${n ? 'var(--youShell)' → style="color:${p.length ? 'var(--youShell)' | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-EMPTY | M7d | `app.js`: if (!n) return `<div class="nsec" data-needs> → if (!p.length) return `<div class="nsec" data-needs> | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-MAP | M7d | `app.js`: .join('') + ds.map(decisionLineHTML).join('') → .join('') | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-MORE | M7d | `app.js`: const more = r && r.decisionsTruncated ? → const more = false ? | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-MORE-EMPTY | M7d | `app.js`: ${more}${notifyHTML(permission)}</div>`; → ${notifyHTML(permission)}</div>`; | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-MORE-LIST | M7d | `app.js`: + more + notifyHTML(permission) → + notifyHTML(permission) | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-RESERVED | M7d | `app.js`: if (x.reserved) return → if (false) return | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-RES-ESC | M7d | `app.js`: ${esc(x.reserved)} approval → ${x.reserved} approval | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-RES-OWNER | M7d | `app.js`: ${esc(RESERVED_OWNERS[x.reserved] \|\| 'its typed owner')} → ${esc('world-publish approve')} | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-RES-FALLBACK | M7d | `app.js`: RESERVED_OWNERS[x.reserved] \|\| 'its typed owner' → RESERVED_OWNERS[x.reserved] | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-OWNER-FINISH | M7d | `app.js`: finish: 'ailang-worldd episode finish' → finish: 'world-publish approve' | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-REFUSED | M7d | `app.js`: if (x.refusal) return → if (false) return | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-REF-SCOPE-ESC | M7d | `app.js`: Decide · ${esc(x.scope)} · refused → Decide · ${x.scope} · refused | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-REF-ESC | M7d | `app.js`: refused: ${esc(x.refusal)}</div> → refused: ${x.refusal}</div> | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-PREFIX | M7d | `app.js`: String(x.request \|\| '').replace(/^sha256:/, '') → String(x.request \|\| '') | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-ID-ESC | M7d | `app.js`: id="decide-${esc(ref)}" → id="decide-${ref}" | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-HREF-ESC | M7d | `app.js`: href="/surface/decide/sha256:${esc(ref)}" → href="/surface/decide/sha256:${ref}" | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-HREF | M7d | `app.js`: href="/surface/decide/sha256: → href="/surface/approve/sha256: | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-SCOPE-ESC | M7d | `app.js`: >Decide · ${esc(x.scope)} · requested → >Decide · ${x.scope} · requested | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-REQ-ESC | M7d | `app.js`: requested by ${esc(x.requester)} · deadline → requested by ${x.requester} · deadline | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-DEADLINE | M7d | `app.js`: deadline ${hm(new Date(x.deadline * 1000))}</a> → deadline ${hm(new Date(x.deadline))}</a> | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-NOTE-DS | M7d | `app.js`: const ds = Array.isArray(r.decisions) ? r.decisions : []; → const ds = []; | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-NOTE-FIRST | M7d | `app.js`: const freshD = seen ? → const freshD = true ? | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-NOTE-SEEN | M7d | `app.js`: !seen.has(x.request) &&  →  | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-NOTE-RES | M7d | `app.js`: && !x.reserved && !x.refusal) → && !x.refusal) | `TestNeedsYouListsDecisions` | C2/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-NOTE-REFUSAL | M7d | `app.js`: && !x.reserved && !x.refusal) → && !x.reserved) | `TestNeedsYouListsDecisions` | C2/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-NOTE-SEENSET | M7d | `app.js`: , ...ds.map(x => x.request)]) → ]) | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-NOTE-CONCAT | M7d | `app.js`: .concat(freshD.map(x => → .concat([].map(x => | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-NOTE-TAG | M7d | `app.js`: tag: x.request, → tag: x.scope, | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-NOTE-BODY | M7d | `app.js`: body: `Decide ${x.scope} requested by ${x.requester}, deadline ${hm(new Date(x.deadline * 1000))}` → body: `Decide ${x.scope} requested by ${x.requester}` | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| N-NOTE-HREF | M7d | `app.js`: href: '/surface/decide/' + x.request → href: '/surface/approve/' + x.request | `TestNeedsYouListsDecisions` | C1/wC | KILLED | KILLED | TestNeedsYouListsDecisions | K |
| P-REF | M7f | `surface_decide_page.go`: if !decisionRefPattern.MatchString(r.PathValue("ref")) { → if false { | `TestDecidePageAssets` | B1/wB | KILLED | KILLED | TestDecidePageAssets | K |
| P-PAGE-BYTES | M7f | `surface_decide_page.go`: _, _ = w.Write(surfacepage.DecidePage) → _, _ = w.Write(surfacepage.ApprovePage) | `TestDecidePageAssets` | B1/wB | KILLED | KILLED | TestDecidePageAssets | K |
| P-SCRIPT-CT | M7f | `surface_decide_page.go`: w.Header().Set("Content-Type", "text/javascript; charset=utf-8")⏎ _, _ = w.Write(surfacepage.DecideScript) → _, _ = w.Write(surfacepage.DecideScript) | `TestDecidePageAssets` | B1/wB | KILLED | KILLED | TestDecidePageAssets | K |
| P-ROUTE-PAGE | M7f | `daemon.go`:  mux.HandleFunc("GET /surface/decide/{ref}", d.approvalRoute(d.handleSurfaceDecidePage))⏎ →  | `TestDecidePageAssets` | B1/wB | KILLED | KILLED | TestDecidePOSTRetryAfterAFailedWorldWrite,TestDecidePageAssets | K |
| P-ROUTE-JS | M7f | `daemon.go`:  mux.HandleFunc("GET /surface/decide.js", d.approvalRoute(d.handleSurfaceDecideScript))⏎ →  | `TestDecidePageAssets` | B1/wB | KILLED | KILLED | TestDecidePageAssets | K |
| P-ROUTE-HOST | M7f | `daemon.go`: mux.HandleFunc("GET /surface/decide.js", d.approvalRoute(d.handleSurfaceDecideScript)) → mux.HandleFunc("GET /surface/decide.js", func(w http.ResponseWriter, r *http.Request) { d.handleSurfaceDecideScript(w, r, nil) }) | `TestDecidePageAssets` | B1/wB | KILLED | KILLED | TestDecidePOSTCounterZeroAssertionLandsOnce,TestDecidePageAssets | K |
| P-PASSKEYS-FIELD | M7e | `surface_decide.go`: Passkeys: passkeys, SSH: → Passkeys: passkeys[:0], SSH: | `TestDecisionSubjectPasskeys` | B4/wB | KILLED | KILLED | TestDecisionSubjectPasskeys | K |
| P-PASSKEYS-ERR | M7e | `surface_decide.go`: passkeys, err := d.decisionPasskeys(ctx)⏎ if err != nil { → passkeys, err := d.decisionPasskeys(ctx)⏎ if false && err != nil { | `TestDecisionSubjectPasskeys` | B1/wB | KILLED | KILLED | TestDecisionSubjectPasskeys | K |
| P-PASSKEYS-OPEN-ERR | M7e | `surface_decide.go`:  if err != nil {⏎  return nil, err⏎ }⏎ defer func() { _ = ro.Close() }()⏎ active, err := ro.ActivePasskeys(ctx) →  if false && err != nil {⏎  return nil, err⏎ }⏎ defer func() { _ = ro.Close() }()⏎ active, err := ro.Ac… | `TestDecisionSubjectPasskeys` | B1/wB | KILLED | KILLED | TestDecisionSubjectPasskeys | K |
| P-PASSKEYS-OLDER | M7e | `surface_decide.go`:  active, err := ro.ActivePasskeys(ctx)⏎ if err != nil && !isOlderStore(err) {⏎  return nil, err⏎ }⏎ var passkeys →  active, err := ro.ActivePasskeys(ctx)⏎ if err != nil {⏎  return nil, err⏎ }⏎ var passkeys | `TestDecisionSubjectPasskeys` | F2/work | SURVIVED | KILLED(whole) | TestRegistryPhaseHasItsOwnClock,TestRegistryPhaseHasItsOwnClock/probe-spends-the-startup-clock | L: unreachable: `serve` upgrades the credentials store, so it is never below schema 6 (sprint-2 precedent) |
| P-PASSKEYS-ACTIVE-ERR | M7e | `surface_decide.go`:  active, err := ro.ActivePasskeys(ctx)⏎ if err != nil && !isOlderStore(err) {⏎  return nil, err⏎ }⏎ var passkeys →  active, err := ro.ActivePasskeys(ctx)⏎ if false && err != nil && !isOlderStore(err) {⏎  return nil, e… | `TestDecisionSubjectPasskeys` | F2/work | SURVIVED | KILLED(whole) | TestRegistryPhaseHasItsOwnClock,TestRegistryPhaseHasItsOwnClock/probe-spends-the-startup-clock | L: unreachable: a non-'older' `ActivePasskeys` error on a store `serve` opened (as sprint 3's `C-LINK-ERR`) |
| P-PASSKEYS-RP | M7e | `surface_decide.go`: surfacepage.CardPasskey{Origin: e.Passkey.Origin, RPID: e.Passkey.RPID, CredentialID: e.Passkey.CredentialID})⏎ }⏎ return surfacepage.DecisionPasskeys → surfacepage.CardPasskey{Origin: e.Passkey.Origin, RPID: e.Passke… | `TestDecisionSubjectPasskeys` | B1/wB | KILLED | KILLED | TestDecisionSubjectPasskeys | K |
| P-PASSKEYS-GROUP | M7e | `decide_model.go`: return RequestCard(StoredRequest{}, passkeys).Passkeys → return RequestCard(StoredRequest{}, passkeys[:min(len(passkeys), 1)]).Passkeys | `TestDecisionSubjectPasskeys` | B1/wB | KILLED | KILLED | TestDecisionSubjectPasskeys | K |
| J-PREFIX | M7f | `decide.js`: m.set(new TextEncoder().encode(PREFIX), 0); → m.set(new TextEncoder().encode(PREFIX), 1); | `TestDecidePageApproveWithPasskey` | B1js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey,TestDecidePageSSHSIGPaste | K |
| J-VERB-BYTE | M7f | `decide.js`: m[18] = VERBS[verb]; → m[18] = VERBS.approve; | `TestDecidePageApproveWithPasskey` | B1js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey,TestDecidePageSSHSIGPaste | K |
| J-PQ | M7f | `decide.js`: m.set(P, 19);⏎    m.set(Q, 51); → m.set(Q, 19);⏎    m.set(P, 51); | `TestDecidePageApproveWithPasskey` | B1js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey,TestDecidePageSSHSIGPaste | K |
| J-EXPIRY | M7f | `decide.js`: setBigUint64(83, BigInt(expiry)) → setBigUint64(83, BigInt(expiry - 1)) | `TestDecidePageApproveWithPasskey` | B1js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey,TestDecidePageSSHSIGPaste | K |
| J-EXPIRY-LE | M7f | `decide.js`: setBigUint64(83, BigInt(expiry)) → setBigUint64(83, BigInt(expiry), true) | `TestDecidePageApproveWithPasskey` | B1js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey,TestDecidePageSSHSIGPaste | K |
| J-HASH-Q | M7f | `decide.js`: 'sha256:' + hex(q) !== ref \|\| → false \|\| | `TestDecidePageFailsClosed` | B1js/wB | KILLED | KILLED | TestDecidePageFailsClosed | K |
| J-HASH-P | M7f | `decide.js`: 'sha256:' + hex(p) !== s.packet) → false) | `TestDecidePageFailsClosed` | B1js/wB | KILLED | KILLED | TestDecidePageFailsClosed | K |
| J-SUBJECT-REF | M7f | `decide.js`: if (s.request !== ref \|\| 'sha256:' → if (false \|\| 'sha256:' | `TestDecidePageFailsClosed` | F2js/work | SURVIVED | SURVIVED(whole) | — | E: `s.request !== ref` is implied by the stronger `'sha256:'+hex(SHA-256(requestBytes)) !== ref` check |
| J-REQREF | M7f | `decide.js`: if (objects.packet.requestRef !== ref) { → if (false) { | `TestDecidePageFailsClosed` | B1js/wB | KILLED | KILLED | TestDecidePageFailsClosed | K |
| J-UNLISTED | M7f | `decide.js`: const f = fields.find((x) => x.path === k); → const f = fields.find((x) => x.path === k) \|\| fields[0]; | `TestDecidePageFailsClosed` | B1js/wB | KILLED | KILLED | TestDecidePageFailsClosed | K |
| J-TYPE | M7f | `decide.js`: typeof obj[k] !== (f.type === 'string' ? 'string' : 'number') \|\| → false \|\| | `TestDecidePageFailsClosed` | B1js/wB | KILLED | KILLED | TestDecidePageFailsClosed | K |
| J-SAFEINT | M7f | `decide.js`: (f.type !== 'string' && !Number.isSafeInteger(obj[k])) → false | `TestDecidePageFailsClosed` | B1js/wB | KILLED | KILLED | TestDecidePageFailsClosed | K |
| J-MISSING | M7f | `decide.js`: if (!Object.prototype.hasOwnProperty.call(obj, f.path)) return → if (false) return | `TestDecidePageFailsClosed` | B1js/wB | KILLED | KILLED | TestDecidePageFailsClosed | K |
| J-NOTOBJ | M7f | `decide.js`: if (obj === null \|\| typeof obj !== 'object' \|\| Array.isArray(obj)) return name; → if (obj === null) return name; | `TestDecidePageFailsClosed` | B1js/wB | KILLED | KILLED | TestDecidePageFailsClosed | K |
| J-B64 | M7f | `decide.js`: pb = unb64u(s.packetBytes); → pb = unb64u(s.packetBytes.replace(/^\*/, '')); | `TestDecidePageFailsClosed` | B1js/wB | KILLED | KILLED | TestDecidePageFailsClosed,TestPasskeyPageAbortsAPendingCeremony,TestPasskeyPageAbortsAPendingCeremony/get-abort | K |
| J-WINDOW | M7f | `decide.js`: !Number.isSafeInteger(s.window.max) \|\| s.window.max < 0) → !Number.isSafeInteger(s.window.max)) | `TestDecidePageFailsClosed` | B1js/wB | KILLED | KILLED | TestDecidePageFailsClosed | K |
| J-DECIDED | M7f | `decide.js`: if (s.decision) { → if (false) { | `TestDecidePageFailsClosed` | B1js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey,TestDecidePageFailsClosed,TestDecidePageSSHSIGPaste | K |
| J-RESERVED | M7f | `decide.js`: if (s.reserved) { → if (false) { | `TestDecidePageFailsClosed` | B1js/wB | KILLED | KILLED | TestDecidePageFailsClosed | K |
| J-REFUSAL | M7f | `decide.js`: if (s.refusal) {⏎      setState('refused'); → if (false) {⏎      setState('refused'); | `TestDecidePageFailsClosed` | B1js/wB | KILLED | KILLED | TestDecidePageFailsClosed | K |
| J-GROUP | M7f | `decide.js`: group = (s.passkeys \|\| []).find((g) => g.origin === location.origin) \|\| null; → group = (s.passkeys \|\| [])[0] \|\| null; | `TestDecidePageSSHSIGPaste` | B1js/wB | KILLED | KILLED | TestDecidePageSSHSIGPaste | K |
| J-RP | M7f | `decide.js`: rpId: group.rp, → rpId: location.hostname, | `TestDecidePageApproveWithPasskey` | F2js/work | SURVIVED | KILLED(whole) | TestWorldContextAndStatusline | E: `group.rp` equals `location.hostname`: the group is chosen by `origin === location.origin` and rp is derived from the origin at enrolment |
| J-ALLOW | M7f | `decide.js`: allowCredentials: group.allow.map((a) => ({ type: 'public-key', id: unb64u(a) })), → allowCredentials: [], | `TestDecidePageApproveWithPasskey` | B2js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey | K |
| J-UV | M7f | `decide.js`: userVerification: 'required', timeout: TIMEOUT_MS } → userVerification: 'preferred', timeout: TIMEOUT_MS } | `TestDecidePageApproveWithPasskey` | B2js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey | K |
| J-TIMEOUT | M7f | `decide.js`:       buttons(true);⏎    }, TIMEOUT_MS); →       buttons(true);⏎    }, TIMEOUT_MS * 2); | `TestDecidePageApproveWithPasskey` | F2js/work | SURVIVED | SURVIVED(whole) | — | L: survivor: decide.js's 240 s abort (A2's rule, copied) is not live-tested; A2's own `TestPasskeyPageAbortsAPendingCeremony` pins the rule on passkey.js; ≈ 40 test lines (Appendix C) |
| J-CHALLENGE | M7f | `decide.js`: const challenge = await sha256(messages[verb]); → const challenge = await sha256(messages.approve); | `TestDecidePageApproveWithPasskey` | B2js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey | K |
| J-POST-VERB | M7f | `decide.js`: Object.assign({ v: 1, verb, packet: → Object.assign({ v: 1, verb: 'approve', packet: | `TestDecidePageApproveWithPasskey` | B2js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey,TestDecidePageSSHSIGPaste | K |
| J-POST-PQ | M7f | `decide.js`: packet: hex(P), request: hex(Q), expiry } → packet: hex(Q), request: hex(P), expiry } | `TestDecidePageApproveWithPasskey` | B2js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey,TestDecidePageSSHSIGPaste | K |
| J-POST-EXPIRY | M7f | `decide.js`: request: hex(Q), expiry } → request: hex(Q), expiry: expiry - 1 } | `TestDecidePageApproveWithPasskey` | B2js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey,TestDecidePageSSHSIGPaste | K |
| J-POST-WEBAUTHN | M7f | `decide.js`: await post(verb, { webauthn: assertionJSON(cred) }); → await post(verb, { sshsig: JSON.stringify(assertionJSON(cred)) }); | `TestDecidePageApproveWithPasskey` | B2js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey | K |
| J-POST-DENY-SIG | M7f | `decide.js`: post('deny', { sshsig: $('sig').value }); → post('approve', { sshsig: $('sig').value }); | `TestDecidePageSSHSIGPaste` | B2js/wB | KILLED | KILLED | TestDecidePageSSHSIGPaste | K |
| J-CREDS | M7f | `decide.js`: credentials: 'omit', → credentials: 'include', | `TestDecidePageApproveWithPasskey` | F2js/work | SURVIVED | SURVIVED(whole) | — | L: survivor: `credentials: 'omit'` is unobservable with a cookie-less temp profile; R0 (`TestDecidePOSTGuards`) refuses any cookie-bearing POST at the daemon anyway |
| J-FRESH | M7f | `decide.js`: >= FRESH_S) return true; → >= 0) return true; | `TestDecidePageApproveWithPasskey` | B2js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey | K |
| J-BLOB | M7f | `decide.js`: new Blob([messages[verb]] → new Blob([messages.approve] | `TestDecidePageSSHSIGPaste` | B2js/wB | KILLED | KILLED | TestDecidePageSSHSIGPaste | K |
| J-SAVE-NAME | M7f | `decide.js`: $(id).download = `world-decision-${verb}-${ref16}.msg`; → $(id).download = `world-decision-${ref16}.msg`; | `TestDecidePageSSHSIGPaste` | B2js/wB | KILLED | KILLED | TestDecidePageSSHSIGPaste | K |
| J-SSHLINE-NS | M7f | `decide.js`: `ssh-keygen -Y sign -n ${NAMESPACE} -f → `ssh-keygen -Y sign -n world-approval@sunholo -f | `TestDecidePageApproveWithPasskey` | B2js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey | K |
| J-VERB-TEXT | M7f | `decide.js`: say(`${b.decided === 'deny' ? 'Denied' : 'Approved'} by → say(`${'Approved'} by | `TestDecidePageApproveWithPasskey` | B2js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey,TestDecidePageSSHSIGPaste | K |
| J-BTN-DENY | M7f | `decide.js`: $('deny').addEventListener('click', () => { decideWithPasskey('deny'); }); → $('deny').addEventListener('click', () => { decideWithPasskey('approve'); }); | `TestDecidePageApproveWithPasskey` | B2js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey | K |
| J-REREAD | M7f | `decide.js`: return load(true); → return setState('decided'); | `TestDecidePageApproveWithPasskey` | B2js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey | K |
| J-RENDER | M7f | `decide.js`: dd.textContent = shown(f, objects[name][f.path]); → dd.textContent = ''; | `TestDecidePageApproveWithPasskey` | B2js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey | K |
| J-PROPOSAL | M7f | `decide.js`: a.href = '/workbench?object=' + encodeURIComponent(objects.packet.proposalHash); → a.href = '/workbench'; | `TestDecidePageApproveWithPasskey` | B2js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey | K |
| J-HIDE | M7f | `decide.js`: for (const id of ['go', 'deny', 'ssh', 'save', 'save-deny', 'origin-note']) $(id).hidden = true; → for (const id of ['origin-note']) $(id).hidden = true; | `TestDecidePageApproveWithPasskey` | B2js/wB | KILLED | KILLED | TestDecidePageApproveWithPasskey | K |

**Sprint-3 named survivors (its §4) — carried, none of their guarded code changed by this sprint:** `D-LAND-ERR`
(`LandKeyDecision`'s error arm in the handler), `C-PROOF-ERR` (`writeOperatorProof`'s error arm), `C-SIGN-NIL` (the
nil decision signer), `C-LINK-ERR` (CD8's `ActivePasskeys` error, below schema 6), `KA-RACE-TEXT` (the race
branch's text) — `git diff BASE M7i -- host/daemon/surface_decide_key.go cmd/ailang-worldd/approve_decide.go
host/store` → empty, and `surface_decide.go`'s diff touches only `handleSurfaceDecision`'s subject (M7e), not the
POST; each keeps its sprint-2/3 reason. The sprint-3 judge's `X-HERR-TOP`, `X-CLI-BODY-PACKET`, `X-CLI-SIG-CAP`,
`X-CLI-ASR-CAP` are now killed (`JG-HERR-TOP`, `JG-BODY-PACKET`, `JG-SIG-CAP`, `JG-ASR-CAP`), and `X-W15-OPERAND`
was never a survivor of the package (`JG-W15-OPERAND`, §1.4 P4).

---

## 5. Files (final `git status --short` after M7i)

```
 M cmd/ailang-worldd/approve.go
 M docs/QUICKSTART.md
 M host/daemon/daemon.go
 M host/daemon/surface_app_test.go
 M host/daemon/surface_approve.go
 M host/daemon/surface_decide.go
 M host/daemon/surface_decide_routes_test.go
 M host/daemon/surface_decide_subject_test.go
 M host/surfacepage/app/app.js
 M host/surfacepage/approve_model.go
 M host/surfacepage/decide_model.go
?? cmd/ailang-worldd/approve_decide_jgaps4_test.go
?? cmd/ailang-worldd/approve_list_decisions.go
?? cmd/ailang-worldd/approve_list_decisions_test.go
?? cmd/ailang-worldd/decide_page_chrome_test.go
?? cmd/ailang-worldd/decide_page_closed_test.go
?? cmd/ailang-worldd/decide_page_rig_test.go
?? cmd/ailang-worldd/decide_page_sshsig_test.go
?? cmd/ailang-worldd/quickstart_decide_card_test.go
?? host/approval/decision_herr_test.go
?? host/daemon/surface_decide_list.go
?? host/daemon/surface_decide_list_test.go
?? host/daemon/surface_decide_page.go
?? host/daemon/surface_decide_page_test.go
?? host/daemon/surface_decide_passkeys_test.go
?? host/surfacepage/app_decide_test.go
?? host/surfacepage/decide.go
?? host/surfacepage/decide/
?? host/surfacepage/decide_list.go
?? host/surfacepage/decide_list_test.go
```
(`git status --short --untracked-files=all` lists `host/surfacepage/decide/decide.js` and `decide/index.html`.)
Plus the two plan files (`design_docs/planned/w-world-approvals-by-key-a4-sprint4-plan.md`,
`design_docs/planned/sprint_w-world-approvals-by-key-a4-s4.json`), which the controller commits. Nothing under
`world/`, `packages/`, `scripts/`, `tools/`, `host/broker/`, `host/coordinator/`, `host/store/`, `host/authority/`,
`host/surfacepage/{approve,passkey}/`.

---

## 6. Sprint gates (controller, after the last sub-milestone's commit; each measured at the base, §1.1)

Run it in **one bash shell, in the foreground, never under `nohup`** (macOS `/bin/bash` 3.2; no `declare -A`; every
leg records its own rc and every leg runs — no `set -e`; every leg is `tee`'d), after `mkdir -p
$HOME/.ailang/state/world-iter392/gates`:

```bash
/bin/bash "$HOME/.ailang/state/world-iter392/plan-scratch/tools/gates_final.sh" \
  /Users/voightkampff/dev/sunholo-data/.wt-world-iter392-a4-s4 "$HOME/.ailang/state/world-iter392/gates" final "$(git -C /Users/voightkampff/dev/sunholo-data/.wt-world-iter392-a4-s4 rev-parse HEAD)"
```

It writes `gates/legs.txt`, one log per leg, and **`gates/ALL_DONE.<sha>`** when every leg has finished (the judge
waits on it). The `verify_go` leg (`VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang
./scripts/verify_go.sh`) reads committed blobs, so the script runs **after** the last commit. The script, verbatim:

```bash
#!/bin/bash
# A4 sprint 4 gates (plan §6). usage: gates.sh <tree> <logdir> <base|final> [<sha>]
# The controller: bash gates.sh <worktree> $HOME/.ailang/state/world-iter392/gates final <HEAD sha>
# macOS /bin/bash 3.2: no declare -A. Every leg records its own rc; every leg runs (no set -e).
# Never start this under nohup. Writes <logdir>/ALL_DONE.<sha> when every leg has finished.
T=$1; L=$2; MODE=$3
mkdir -p "$L"; cd "$T" || exit 2
SHA=${4:-$(git rev-parse HEAD 2>/dev/null || echo unknown)}
export AILANG_BIN=$HOME/.pinned-ailang/ailang
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
: > "$L/legs.txt"
FAILS=""
s0=$(date +%s)
leg() { # leg <name> <command string>: tee'd to $L/<name>.log, rc recorded in legs.txt
  local name=$1 cmd=$2 s rc
  s=$(date +%s)
  bash -c "$cmd" 2>&1 | tee "$L/$name.log"
  rc=${PIPESTATUS[0]}
  echo "LEG $name rc=$rc $(( $(date +%s) - s ))s" | tee -a "$L/legs.txt"
  [ "$rc" -eq 0 ] || FAILS="$FAILS $name"
}
leg status 'git status --porcelain | wc -l'
leg verify_ail './scripts/verify_ail.sh'
leg gofmt 'out=$(gofmt -l cmd host); echo "$out"; test -z "$out"'
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./... > '$L/fence.raw' 2>&1; rc=\$?; cat '$L/fence.raw'; if grep -E '^(FAIL|--- FAIL)' '$L/fence.raw' >/dev/null; then exit 1; fi; exit \$rc"
leg pkgs 'go test -count=1 -timeout 30m ./host/approval/... ./host/authority/ ./host/store/ ./host/surfacepage/ ./host/coordinator/ ./host/daemon/ ./cmd/ailang-worldd/'
leg broker_named 'go test -count=1 ./host/broker/ -run "^(TestReservedApprovalScope|TestDecisionRefusal|TestDecideApproval|TestPublishApproval|TestApprovalContextIdentity|TestRecentApprovals|TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential)" -v'
leg race 'go test -race -count=1 ./host/surfacepage/ && go test -race -count=1 ./host/daemon/ -run "^(TestDecision|TestDecide|TestApprov|TestSurfaceApprove|TestSurfaceRequests|TestSurfaceApp|TestDaemonCredentialWriterSites)" && go test -race -count=1 ./cmd/ailang-worldd/ -run "^(TestApproveDecide|TestApproveList|TestApprove)"'
NEW_A="TestDecideByAssertionRefusesAMissingMessageAtW6 TestSignCountLaw"
NEW_SP="TestDecideModelListsPendingDecisionsOnly TestDecideModelListsATwoPacketRequestRefused TestDecideModelTruncation TestNeedsYouListsDecisions"
NEW_D="TestSurfaceRequestsListsDecisions TestSurfaceRequestsDecisionsTruncated TestSurfaceRequestsDecisionsUnreadable TestDecisionSubjectPasskeys TestDecisionSubjectJSON TestDecidePageFieldTable TestDecidePageAssets"
NEW_C="TestApproveDecideArtifactPostsTheSubjectsDigests TestApproveDecideArtifactReadBounds TestApproveListShowsDecisions TestApproveListDecisionLines TestQuickstartSection18CardMatchesThePage"
NEW_X="TestDecidePageApproveWithPasskey TestDecidePageFailsClosed TestDecidePageSSHSIGPaste"
if [ "$MODE" = base ]; then NEW_A=""; NEW_SP=""; NEW_D=""; NEW_C=""; NEW_X=""; fi
alt() { local x="$*"; x="${x// /|}"; echo "${x:+|$x}"; }
INV_APPROVAL="TestApprovalImportsOnlyTheStandardLibrary TestApprovalTestIsTestOnly TestDecisionMessageLayout TestDecisionMessageRefusesBadInput"
INV_BROKER="TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential TestQuickstartSection10ProfilesLoad TestQuickstartProtectedContract TestApprovalContextIdentity TestReservedApprovalScope TestDecideApprovalCallerCensus"
INV_DAEMON="TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestReadCtxCancelledAfterHandler TestMCPMountMethodSource TestDaemonCredentialWriterSites TestQuickstartDecisionPacketSection TestDaemonConstructorCallerInventory TestApprovePageRoutes TestPasskeyPageRoutes TestApprovalRoutesRefuseForeignHost TestGETRoutesRejectOtherMethods TestDaemonCallerInventoryAliasControl"
INV_BOUNDARY="TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard TestWorkbenchPackageRemainsTransportFree"
INV_STORE="TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext TestCredentialsMethodsRequireADeadline TestCredentialsSchemaMatchesCanonicalDDL"
INV_VG="TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestVerifygateTestWritesAreForkLocked"
INV_SP="TestSurfaceAppPure TestSurfaceScriptPure TestSurfaceFontEmbed TestSurfaceScriptAssetSeams"
INV_CMD="TestEveryVerbHasHelp TestQuickstartAttendedStepsBoxMatchesSection9 TestQuickstartSection9FlagsMatchTheCLI TestQuickstartSection10FlagsMatchTheCLI TestQuickstartSection11FlagsMatchTheCLI TestQuickstartSection18FlagsMatchTheCLI TestQuickstartServiceSectionMatchesTheCLIHelp"
INV_RUNBOOK="TestNoCIStepOrScriptReachesThePublishEntrypoint"
leg inv "set -o pipefail
go test -v -count=1 ./host/approval/... -run '^($(echo $INV_APPROVAL | tr ' ' '|')$(alt $NEW_A))\$' &&
go test -v -count=1 ./host/broker/ -run '^($(echo $INV_BROKER | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/daemon/ -run '^($(echo $INV_DAEMON | tr ' ' '|')$(alt $NEW_D))\$' &&
go test -v -count=1 ./host/boundary/ -run '^($(echo $INV_BOUNDARY | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/store/ -run '^($(echo $INV_STORE | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/verifygate/ -run '^($(echo $INV_VG | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/runbook/ -run '^($(echo $INV_RUNBOOK | tr ' ' '|'))\$' &&
WORLD_EXEC_NODE=/opt/homebrew/bin/node go test -v -count=1 ./host/surfacepage/ -run '^($(echo $INV_SP | tr ' ' '|')$(alt $NEW_SP))\$' &&
go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(echo $INV_CMD | tr ' ' '|')$(alt $NEW_C))\$'"
XB="TestApprovePageInChrome TestApprovePageRefusesASwappedRequest TestApprovePageSSHSIGInChrome TestApprovePageShowsAHostileValueAsText TestPasskeyPageAbortsAPendingCeremony TestPasskeyPageInChrome TestApproveDecideWebAuthnLinkInChrome"
leg chrome "WORLD_CHROME='/Applications/Google Chrome.app/Contents/MacOS/Google Chrome' WORLD_REQUIRE_PASSKEY_BROWSER=1 WORLD_REQUIRE_SSH_KEYGEN=1 go test ./cmd/ailang-worldd/ -run '^($(echo $XB | tr ' ' '|')$(alt $NEW_X))\$' -count=1 -v"
for t in $XB $NEW_X; do grep -qE -- "^--- PASS: $t " "$L/chrome.log" || { echo "CHROME NOT PASSED: $t" | tee -a "$L/legs.txt"; FAILS="$FAILS chrome_missing"; }; done
# every named test must PASS (none SKIP, none FAIL)
MISSING=""
for t in $INV_APPROVAL $NEW_A $INV_BROKER $INV_DAEMON $INV_BOUNDARY $INV_STORE $INV_VG $INV_SP $INV_CMD $INV_RUNBOOK $NEW_SP $NEW_D $NEW_C; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || MISSING="$MISSING $t"
done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then FAILS="$FAILS inv_skip_or_fail"; fi
[ -z "$MISSING" ] || { echo "NOT PASSED:$MISSING" | tee -a "$L/legs.txt"; FAILS="$FAILS inv_missing"; }
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.log")" | tee -a "$L/legs.txt"
leg verify_go 'VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh'
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" | tee -a "$L/legs.txt"
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" > "$L/ALL_DONE.$SHA"
```

**Expected:** every leg rc 0 except a known-load red, re-run alone once (§1.1's rule); `inv` prints **named PASS**
for every inventory and every new test (none missing, none SKIP); `chrome` prints `--- PASS` for the seven A2/A3/
sprint-3 page tests and the three card tests. The planner's run of this script on the final tree: §1.2 C18.

---

## 7. Executor directive notes

- **No git writes of any kind** (no add, commit, stash, checkout, reset, branch). **`patch -p1` (§2) is the only
  writer** of the tree; no hand edit, no `gofmt -w`, no fuzz. A failed `--dry-run` or a red acceptance → **stop
  and report** with the log; do not repair.
- **Runs**: Run 1 = M7a, M7b, M7c; Run 2 = M7d, M7e, M7f; Run 3 = M7g, M7h, M7i. After each sub-milestone's acceptance, snapshot to
  `$HOME/.ailang/state/world-iter392/snap/<M>/`, **copying only the files `git status --porcelain` reports**
  (cumulative against the base: modified and untracked files, plus a `DELETED` list of deleted paths — none are
  expected in this sprint), **never an `rsync` of whole directories** (the controller's `commit_ms.sh` commits
  every file present in a snapshot):
  ```bash
  bash <<'EOF2'
  set -euo pipefail
  M=M7a   # the sub-milestone just accepted
  W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter392-a4-s4
  D="$HOME/.ailang/state/world-iter392/snap/$M"; rm -rf "$D"; mkdir -p "$D"
  cd "$W"
  git status --porcelain --untracked-files=all | grep -v -e 'w-world-approvals-by-key-a4-sprint4-plan.md$' -e 'sprint_w-world-approvals-by-key-a4-s4.json$' > "$D/STATUS"
  awk '$1 != "D" {print $NF}' "$D/STATUS" | while read -r p; do mkdir -p "$D/$(dirname "$p")"; cp -p "$p" "$D/$p"; done
  awk '$1 == "D" {print $NF}' "$D/STATUS" > "$D/DELETED"
  (cd "$D" && find cmd host docs -type f -print0 2>/dev/null | sort -z | xargs -0 shasum -a 256) > "$D/MANIFEST.sha256"
  cat "$D/STATUS"; wc -l < "$D/MANIFEST.sha256"
  EOF2
  ```
  The two plan files are excluded. The controller commits each snapshot as one commit and checks the numstat
  against §3's expected value (`git diff --numstat 799660ac HEAD -- cmd host docs | awk '{s+=$1+$2} END {print s}'`):
  M7a 106, M7b 437, M7c 587, M7d 719, M7e 810, M7f 1,350, M7g 1,635, M7h 1,878, M7i 1,947. A different number → stop.
- **Every Go Bash call passes `timeout: 600000`, runs in the foreground, no Monitor.** **Never `nohup`, never
  background a test leg.** `export AILANG_BIN=$HOME/.pinned-ailang/ailang` in every shell that runs Go tests. The
  controller's shell is zsh: run each acceptance block under `bash <<'EOF' … EOF`.
- **Exact acceptance commands** are §3's, verbatim. The live-Chrome lines carry
  `WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" WORLD_REQUIRE_PASSKEY_BROWSER=1` on
  the command line itself (copy it verbatim; the A2 driver makes its own temp `--user-data-dir`; never a real
  profile). The node lines carry `WORLD_EXEC_NODE=/opt/homebrew/bin/node` (with it set, a node that cannot run the
  harness fails instead of skipping). `ssh-keygen` (`/usr/bin`) must be on PATH for the real-tool sibling.
- **Stop conditions**: a dry-run failure; any `--- FAIL`, `--- SKIP` (a Chrome or node line with its env must not
  skip) or `FAIL` line in an acceptance; gofmt output; a numstat mismatch; any file outside §5 changed (`git status
  --short` after each patch must equal §5's list up to that stage).
- **Known load-sensitive tests** (the directive's list, incl. `TestPasskeyPageAbortsAPendingCeremony/get-late`
  under load) are not this sprint's: a red there is re-run once alone.
- D-WORLD-109: nothing touches `~/.ailang/world`, `~/.ssh`, a real `/dev/tty`, a real browser profile, or a real
  daemon (§0.1).

---

## 8. Risks

- **R1 — RS-VERB: the A2 page's wording for a CLI deny link** (S4-D2): "Approve request decide-deny-…" can still
  mislead a hurried operator on the `--webauthn-link` path. Mitigated: the card (the phone's path now) shows its own
  verb (XD1 asserts it); the link path carries the verb in the CLI's capitals, the link's `request` text and the
  signed hash; the cryptography refuses a cross-verb use (DM5, W6 rows; XD4's mismatched-verb paste → 422). A fix
  needs a design amendment (the A2 page and its A3 pin are frozen).
- **R2 — RS-RIG: Chrome drives the card against the rig's POST, not the daemon's** (S4-D11). The rig applies the
  daemon's own grammar, law and judges, and A3a's guards; the daemon's arm is unchanged and proven by
  `TestDecidePOST*` and sprint 3's live-Chrome CLI path. The first end-to-end card → daemon POST with a real
  authenticator is the attended run's step 4.
- **R3 — Live Chrome only on the controller's machine** (CI has no Chrome step, A3 V13): M7g/M7h are gated as XB1;
  a CI run skips them, counted. `TestPasskeyPageAbortsAPendingCeremony/get-late` is a known load flake (391).
- **R4 — Size margin 53.** Any executor-side growth is a stop (§7's numstat check), never a hand edit.
- **R5 — The listing's cost** (§1.4 P3): ≤ 20 `ValidatePacket` calls, each with its own ≤ 101-packet scan, per
  `requests.json` read (the surface re-reads on every refresh). Bounded (≈ 2,000 object reads worst case) and read
  only; a busy world may slow *Needs you*. An indexed packet lookup (RK3b's fix) removes it.
- **R6 — `--list` reads the world store directly** (S4-D3). Read-only (`store.OpenReadOnly`, SQLite WAL readers do
  not block the daemon's writer); the same function the daemon serves, so the lists cannot disagree in kind. It
  needs `--db` to name the daemon's store (the default is the installed one, as for the mint list).
- **R7 — Load.** The rig ran at load 18–28 during this planner's runs; every base red was a load timeout in an
  untouched test (§1.1). Re-run once alone before believing a red.
- **R8 — D-WORLD-115 = B** changes `decide.js`'s message builder with the laws and CD6/CD8 (S4-D9).

---

## 9. Acceptance criteria (the judge scores against these)

1. **`requests.json` lists decisions** (DMD1–DMD3): pending `Human.Approve` requests with ≥ 1 packet; exactly one →
   packet + deadline; two → listed with D4b and no packet; reserved marked with D5; D7 past the deadline; decided
   and packet-less requests absent; `decisionsTruncated` for either bound (21 chain requests, 101 packets; not 100);
   `decisions` is `[]`, never null, and a world-store error leaves the mint lists intact.
2. **`approve --list` prints the same decisions** from the same function, offline (never a dial), with the
   decide command for each decidable one; the existing `--list` tests unedited.
3. ***Needs you*** (XD3): the decidable line links `/surface/decide/<ref>`; reserved and refused lines have no link;
   truncation is said; every value escaped; the notification fires once for a new decidable decision id and never for
   reserved/refused/seen ones or before the first read; inputs without `decisions` render byte-identically;
   `TestSurfaceAppPure` unedited.
4. **`DecisionSubject.Passkeys`** groups the active passkeys by origin with their rp and credential ids; revoked
   excluded; a store error is a 500.
5. **The card routes** (PD2): `GET /surface/decide/{ref}` and `GET /surface/decide.js` serve the embedded bytes under
   the **unchanged** `surfaceCSP` behind `approvalRoute`; malformed ref 400; foreign Host 421; no cookie, no CORS
   header; A2/A3 page hashes unchanged; `TestSurfaceAppRoutes` +2, `TestFrozenV1RouteTableMatchesMux` 13.
6. **PD1**: `DECISION_FIELDS` equals the key and type set of the real stored request and packet (12 fields).
7. **XD1 in live Chrome**: Approve and Deny with a passkey; each challenge equals Go's `DecisionMessageHash` with the
   `BigInt` expiry; the POST carries the verb, P, Q, expiry and `webauthn` only; the decision is
   `webauthn:<name>`; **the card's verb text is correct for approve and for deny**; the < 60 s re-read.
8. **XD2 in live Chrome**: every tampered subject (hashes, `requestRef`, unknown/missing key, types, base64url,
   window) is refused **before any `get`**; refused/reserved/decided subjects show no control.
9. **XD4 in live Chrome**: the saved files are Go's 91-byte `DecisionMessage`; a pasted approvaltest SSHSIG
   decides; a mismatched verb is refused; **the real `ssh-keygen` sibling decides**.
10. **The sprint-3 judge's notes 2–4 are closed by tests** (`JG-HERR-TOP`, `JG-BODY-PACKET`, `JG-BODY-REQUEST`,
    `JG-SIG-CAP`, `JG-ASR-CAP` killed) and note 1 is shown to be a narrow-run survivor (`JG-W15-OPERAND` killed by
    `TestSignCountLaw`).
11. **QUICKSTART §18** documents the card and `--list` with the card's own controls
    (`TestQuickstartSection18CardMatchesThePage`); `TestQuickstartSection18FlagsMatchTheCLI` unedited.
12. **Inventories**: `TestFrozenV1RouteTableMatchesMux` 13, `TestSurfaceAppRoutes` (+2), `TestDaemonDependencyAllowlist`,
    `TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots`, `TestProductionGoSurface`,
    `TestDaemonCredentialWriterSites`, `TestCredentialsMethodsRequireADeadline`, `TestSurfaceAppPure`,
    `TestSurfaceScriptAssetSeams`, `TestApprovePageRoutes`, `TestPasskeyPageRoutes`,
    `TestGETRoutesRejectOtherMethods`, `TestApprovalImportsOnlyTheStandardLibrary`, `TestEveryVerbHasHelp`, the
    quickstart tests, `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` — all PASS; edited only as §1.4
    P7/P8 and the design's +2.
13. **Size** ≤ 2,000 (`-- cmd host docs`): 1,947 measured; each sub-milestone's numstat as §3.
14. **Every mutant in §4 behaves as expected**, its id in the ledger; no kill rests only on an inventory test.
15. **The round trip is identical** (§1.2 C16) and the gates (§6) pass, `verify_go` after commit included.
16. **No page or POST arm without an operator path** (§3.1).

**What an operator can do after this sprint:** open the surface (studio, or the phone over the tailnet origin), see
each pending `Human.Approve` request under ***Needs you*** ("Decide · scope · requested by … · deadline hh:mm"), open
its **card**, read every field of the request and its packet, and **Approve or Deny with a passkey** (Face ID / Touch
ID), or **save the message**, sign it with an enrolled SSH key anywhere (`ssh-keygen -Y sign -n world-decision@sunholo`)
and paste the signature; see `decidedBy`, then the request leaves *Needs you* and the packet page reads `approvals
chain: approve by webauthn:<name>`. From a terminal: `ailang-worldd approve --list --db …` lists the same requests
with their decide command, and every sprint 1–3 CLI path is unchanged. A notification (opt-in) announces a new one.

**What the attended exit run (design §10) still waits for:** Mark's run itself (D-WORLD-109) — steps 1–5 are all
runnable now; step 4 is the first card → daemon POST with a real authenticator on the phone over the tailnet
(RS-RIG) and must record which `Host` arrived (A3 §10.2's open question); RS-VERB (the CLI link's page wording) is
for Mark to judge; D-WORLD-115's answer.

---

## Appendix C — after this sprint (A4's close and the attended run's handoff)

- **Moved out: nothing.** A4 is code-complete with this sprint (design §0 sprints 1–3 + this sprint's surface).
- **Named survivors carried**: `D-LAND-ERR`, `C-PROOF-ERR`, `C-SIGN-NIL` (sprint 2's reasons), `C-LINK-ERR`,
  `KA-RACE-TEXT` (sprint 3's); this sprint's `L-VALIDATE-ERR`, `CL-LIST-ERR`, `P-PASSKEYS-OLDER`,
  `P-PASSKEYS-ACTIVE-ERR` (unreachable), `J-TIMEOUT` (decide.js's 240 s abort is not live-tested; ≈ 40 test lines
  would close it), `J-CREDS` (unobservable; R0 holds at the daemon).
- **Residuals for the attended exit run (design §10, Mark; NOT run by the loop):**
  - **RS-RIG** — the first card → daemon `POST /surface/decide/{ref}` with a real authenticator is §10 step 4
    (record the `Host` that arrived); a cmd-package seeded-world Chrome test is ≈ +150 lines if wanted later.
  - **RS-VERB** — the CLI's `--webauthn-link` path still opens the A2 page reading "Approve request
    decide-deny-…" for a deny (the design freezes that page and an A3 test pins its bytes). A one-line design
    amendment (§2.8/§13) and ≈ 30 lines would give the A2 link a verb member; Mark decides whether it is worth it
    now that the card is the phone's path.
  - **R5** — the listing's ≤ 20 `ValidatePacket` calls per read; watch *Needs you*'s latency in a busy world.
- **D-WORLD-115** stays open (default A applied); S4-D9 names what B changes.

---

## Appendix P — the patches (apply with §2; generated from the compiled scratch stages; do not edit)

### Patch M7a — the sprint-3 judge's notes 2–4 (test only)

````diff M7a
diff --git a/cmd/ailang-worldd/approve_decide_jgaps4_test.go b/cmd/ailang-worldd/approve_decide_jgaps4_test.go
new file mode 100644
index 0000000..d1b2e8e
--- /dev/null
+++ b/cmd/ailang-worldd/approve_decide_jgaps4_test.go
@@ -0,0 +1,69 @@
+package main
+
+// Row 201 A4 sprint 4 M7a (the sprint-3 judge's non-blocking findings 3 and 4): the packet and
+// request digests the CLI posts with --signature and --webauthn-assertion are the subject's, and
+// both artifact reads refuse a file one byte over their bound before any POST. The fake daemon
+// (approve_decide_rig_test.go) answers any packet, so only these checks pin the body's digests.
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"os"
+	"path/filepath"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+// The --signature and --webauthn-assertion bodies carry P and Q as the subject served them.
+func TestApproveDecideArtifactPostsTheSubjectsDigests(t *testing.T) {
+	f := newFakeDecider(t)
+	env := approveEnv(nil, nil, atTime(1710))
+	f.enrolDecisionPasskey(t, true)
+	sigPath := filepath.Join(t.TempDir(), "decide.sig")
+	if err := os.WriteFile(sigPath, []byte("SIGNED-ELSEWHERE"), 0o600); err != nil {
+		t.Fatal(err)
+	}
+	h := sha256.Sum256(f.decisionMessageOf(t, approval.VerbApprove, 1760))
+	asrPath := filepath.Join(t.TempDir(), "world-assertion.json")
+	if err := os.WriteFile(asrPath, cdPasskey.Assert(cdOrigin, cdRP, h, 2), 0o600); err != nil {
+		t.Fatal(err)
+	}
+	p, q := strings.TrimPrefix(f.subject.Packet, "sha256:"), strings.TrimPrefix(f.ref, "sha256:")
+	for i, args := range [][]string{{"--signature", sigPath, "--expiry", "1760"}, {"--webauthn-assertion", asrPath, "--expiry", "1760"}} {
+		if code, _, stderr := f.keyRun(t, env, args...); code != exitOK || len(f.posts) != i+1 {
+			t.Fatalf("%v: rc %d, %d POSTs\nstderr %s", args, code, len(f.posts), stderr)
+		}
+		if b := f.postBodies[i]; b.Packet != p || b.Request != q || b.Verb != "approve" || b.V != 1 {
+			t.Errorf("%v posted packet %q request %q verb %q v %d, want %q %q approve 1", args, b.Packet, b.Request, b.Verb, b.V, p, q)
+		}
+	}
+}
+
+// Each artifact flag reads at most its bound: one byte over is refused with the bound named and
+// nothing is posted; exactly the bound is read (and, not being an artifact, refused or posted).
+func TestApproveDecideArtifactReadBounds(t *testing.T) {
+	f := newFakeDecider(t)
+	env := approveEnv(nil, nil, atTime(1710))
+	for _, c := range []struct {
+		flag string
+		max  int
+	}{{"--signature", approveMaxSignatureBytes}, {"--webauthn-assertion", approval.MaxPasskeyResponseBytes}} {
+		over := filepath.Join(t.TempDir(), "over")
+		if err := os.WriteFile(over, bytes.Repeat([]byte("x"), c.max+1), 0o600); err != nil {
+			t.Fatal(err)
+		}
+		code, _, stderr := f.keyRun(t, env, c.flag, over, "--expiry", "1760")
+		if code != exitUsage || len(f.posts) != 0 || !strings.Contains(stderr, "over the "+itoa(c.max)+"-byte bound") {
+			t.Errorf("%s one byte over %d: rc %d, %d POSTs\nstderr %s", c.flag, c.max, code, len(f.posts), stderr)
+		}
+	}
+	at := filepath.Join(t.TempDir(), "at")
+	if err := os.WriteFile(at, bytes.Repeat([]byte("x"), approveMaxSignatureBytes), 0o600); err != nil {
+		t.Fatal(err)
+	}
+	if code, _, stderr := f.keyRun(t, env, "--signature", at, "--expiry", "1760"); code != exitOK || len(f.posts) != 1 || len(*f.postBodies[0].SSHSIG) != approveMaxSignatureBytes {
+		t.Errorf("--signature at the bound: rc %d, %d POSTs\nstderr %s", code, len(f.posts), stderr)
+	}
+}
diff --git a/host/approval/decision_herr_test.go b/host/approval/decision_herr_test.go
new file mode 100644
index 0000000..6aaf425
--- /dev/null
+++ b/host/approval/decision_herr_test.go
@@ -0,0 +1,37 @@
+package approval_test
+
+// Row 201 A4 sprint 4 M7a (the sprint-3 judge's non-blocking finding 2): a message hash that cannot
+// be built (hErr) is refused at W6, in W6's place, never before W1-W5: a decision with no message
+// AND an earlier failing step is refused by that earlier step, as the A2 order has it.
+
+import (
+	"errors"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+)
+
+func TestDecideByAssertionRefusesAMissingMessageAtW6(t *testing.T) {
+	for _, tc := range []struct {
+		name string
+		edit func(t *testing.T, in *approval.DecisionAssertionInput)
+		want error
+	}{
+		{"W1-envelope", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Response = []byte("{}") }, approval.ErrPasskeyResponse},
+		{"W2-raw-id", func(t *testing.T, in *approval.DecisionAssertionInput) {
+			in.Response = dmPasskey.Assert(paOrigin, paRPID, [32]byte{}, 2, approvaltest.Edit{RawID: []byte("another")})
+		}, approval.ErrCredentialIDMismatch},
+		{"W5-type", func(t *testing.T, in *approval.DecisionAssertionInput) {
+			in.Response = dmPasskey.Assert(paOrigin, paRPID, [32]byte{}, 2, approvaltest.Edit{Type: "webauthn.create"})
+		}, approval.ErrClientDataType},
+		{"W6-itself", func(*testing.T, *approval.DecisionAssertionInput) {}, approval.ErrChallenge},
+	} {
+		in := dmaInput(t, 2)
+		in.Verb = 3 // no decision message: hErr
+		tc.edit(t, &in)
+		if d, err := approval.DecideByAssertion(in); d != nil || !errors.Is(err, tc.want) {
+			t.Errorf("%s: (%v, %v), want %v", tc.name, d, err, tc.want)
+		}
+	}
+}
````

### Patch M7b — requests.json decisions (DMD1–DMD3)

````diff M7b
diff --git a/host/daemon/surface_approve.go b/host/daemon/surface_approve.go
index 1a4b34c..b515dca 100644
--- a/host/daemon/surface_approve.go
+++ b/host/daemon/surface_approve.go
@@ -104,6 +104,12 @@ func (d *Daemon) handleSurfaceRequests(w http.ResponseWriter, r *http.Request, _
 	if out.Decided == nil {
 		out.Decided = []surfacepage.DecidedRequest{}
 	}
+	if err := d.readDecisionList(ctx, now, &out); err != nil { // row 201 A4 §2.8: the world store's decisions
+		d.writeInternalErrorLog(r, err)
+	}
+	if out.Decisions == nil {
+		out.Decisions = []surfacepage.PendingDecision{}
+	}
 	writeApprovalJSON(w, http.StatusOK, out)
 }
 
diff --git a/host/daemon/surface_decide_list.go b/host/daemon/surface_decide_list.go
new file mode 100644
index 0000000..ccb3e9d
--- /dev/null
+++ b/host/daemon/surface_decide_list.go
@@ -0,0 +1,80 @@
+package daemon
+
+// Row 201 A4 sprint 4 (design_docs/planned/w-world-approvals-by-key-a4.md §2.8): the decisions
+// GET /surface/requests.json lists and approve --list prints. One read of the approvals chain
+// (broker.RecentApprovals: at most 20 summaries of 40 heads) and one packet scan of at most
+// decisionListScan packets; each request one scanned packet carries gets ValidatePacket's verdict and
+// the law's display reason at now. Display only: the POST runs the complete discovery and the law.
+
+import (
+	"context"
+
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/coordinator"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// decisionListScan is the listing's packet scan: the packet validator's 100-packet bound plus one
+// sentinel, whose presence marks the list truncated.
+const decisionListScan = 101
+
+// PendingDecisions reads the decisions list from the world store s at now (the daemon's handle or
+// the CLI's read-only one). Store errors and a malformed approvals chain are returned, never a row.
+func PendingDecisions(ctx context.Context, s *store.Store, now int64) ([]surfacepage.PendingDecision, bool, error) {
+	page, err := broker.RecentApprovals(ctx, s, 40)
+	if err != nil {
+		return nil, false, err
+	}
+	scanned, err := s.ObjectsBySemanticID(ctx, coordinator.DecisionPacketV1, "", decisionListScan)
+	if err != nil {
+		return nil, false, err
+	}
+	scanTruncated := len(scanned) == decisionListScan
+	carrying := map[string][]store.Object{}
+	for _, o := range scanned[:min(len(scanned), decisionListScan-1)] {
+		full, ok, err := s.GetObject(ctx, o.Hash)
+		if err != nil {
+			return nil, false, err
+		}
+		if p, err := coordinator.DecodeDecisionPacket(full.Payload); ok && err == nil {
+			carrying[p.RequestRef.String()] = append(carrying[p.RequestRef.String()], full)
+		}
+	}
+	var cands []surfacepage.DecisionCandidate
+	for _, sm := range page.Summaries {
+		c := surfacepage.DecisionCandidate{Request: sm.RequestRef, Effect: sm.Effect, Scope: sm.Scope, Requester: sm.Requester,
+			Cost: sm.Cost, Decided: sm.Status != "pending"}
+		packets := carrying[sm.RequestRef]
+		for _, p := range packets {
+			c.Packets = append(c.Packets, p.Hash.String())
+		}
+		if !c.Decided && sm.Effect == broker.EffectHumanApprove && len(packets) != 0 {
+			f := broker.DecisionFacts{Ref: sm.RequestRef, RefOK: true, RequestOK: true, OnChain: true, Scope: sm.Scope, Now: now,
+				Discovery: broker.Discovery{Complete: true, Scanned: len(scanned), Carrying: len(packets)}}
+			if len(packets) == 1 {
+				res, err := coordinator.ValidatePacket(ctx, s, packets[0])
+				if err != nil {
+					return nil, false, err
+				}
+				f.PacketCode, f.PacketSentence, f.Deadline = res.Code, packetRefusalSentence(res), res.Packet.DeadlineAt
+				c.Deadline = res.Packet.DeadlineAt
+			}
+			c.Reserved, _ = broker.ReservedApprovalScope(sm.Scope)
+			c.Refusal = broker.SubjectRefusal(f).Reason
+		}
+		cands = append(cands, c)
+	}
+	rows, truncated := surfacepage.FoldDecisions(cands, page.Truncated, scanTruncated)
+	return rows, truncated, nil
+}
+
+// readDecisionList fills out's decisions; a read that fails is logged by the caller and lists none.
+func (d *Daemon) readDecisionList(ctx context.Context, now int64, out *surfacepage.RequestsJSON) error {
+	rows, truncated, err := PendingDecisions(ctx, d.store, now)
+	if err != nil {
+		return err
+	}
+	out.Decisions, out.DecisionsTruncated = rows, truncated
+	return nil
+}
diff --git a/host/daemon/surface_decide_list_test.go b/host/daemon/surface_decide_list_test.go
new file mode 100644
index 0000000..7600053
--- /dev/null
+++ b/host/daemon/surface_decide_list_test.go
@@ -0,0 +1,137 @@
+package daemon
+
+// Row 201 A4 sprint 4 (design §2.8, DMD1–DMD3 end to end): GET /surface/requests.json lists the
+// pending Human.Approve requests a decision card decides, through d.Handler(). Temp stores only.
+
+import (
+	"encoding/json"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/coordinator"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// requestsJSON is GET /surface/requests.json from the loopback, decoded, and its raw body.
+func (r decideRig) requestsJSON() (surfacepage.RequestsJSON, string) {
+	r.t.Helper()
+	rec := r.serve("GET", "/surface/requests.json", apLoopback, nil)
+	var out surfacepage.RequestsJSON
+	if rec.Code != 200 || json.Unmarshal(rec.Body.Bytes(), &out) != nil {
+		r.t.Fatalf("requests.json: %d %s", rec.Code, rec.Body.String())
+	}
+	return out, rec.Body.String()
+}
+
+// TestSurfaceRequestsListsDecisions: a decidable request with its packet and deadline, a reserved
+// publish-shaped one with the law's D5 reason, a shared one (two packets) with D4b and no packet; a
+// packet no commit carries listed with D4b's packet sentence; a decided request, a request no packet
+// carries, are absent; the mint lists are unchanged; past the
+// deadline the row carries D7.
+func TestSurfaceRequestsListsDecisions(t *testing.T) {
+	r := newDecideRig(t)
+	if out, raw := r.requestsJSON(); !strings.Contains(raw, `"decisions":[]`) || !strings.Contains(raw, `"decisionsTruncated":false`) || len(out.Decisions) != 0 {
+		t.Fatalf("an empty world: %s", raw)
+	}
+	q1, p1 := r.packet("list-ok", broker.WorkspaceScope, "agent-1")
+	qr, pr := r.packet("list-reserved", "registry:pkg#effect=Registry.Publish", "agent-r")
+	qs, _ := r.packet("list-shared-1", broker.WorkspaceScope, "agent-s")
+	r.packet("list-shared-2", broker.WorkspaceScope, "agent-s")
+	qd, _ := r.packet("list-decided", broker.WorkspaceScope, "agent-d")
+	qz := approvalRequestInChain(t, r.d, broker.WorkspaceScope, "agent-z", packetNow)
+	qu := approvalRequestInChain(t, r.d, broker.WorkspaceScope, "agent-u", packetNow) // its packet is stored, never committed (code 2)
+	pu := forgeObject(coordinator.DecisionPacketV1, coordinator.EncodeDecisionPacket(coordinator.DecisionPacket{ProposalHash: hashref.SumSHA256([]byte("list-uncarried")),
+		RequestRef: qu, CreatedAt: packetNow, DeadlineAt: packetNow + 60, Policy: "Cancel"}))
+	if err := r.d.store.PutObject(boundedTestContext(t), pu); err != nil {
+		t.Fatal(err)
+	}
+	if _, err := broker.DecideApproval(boundedTestContext(t), r.d.store, qd, "approve", typedDecider, decNow); err != nil {
+		t.Fatal(err)
+	}
+	out, raw := r.requestsJSON()
+	if out.DecisionsTruncated || out.Pending == nil || out.Decided == nil || !strings.Contains(raw, `"pending":[]`) {
+		t.Fatalf("truncated or mint lists changed: %s", raw)
+	}
+	by := map[string]surfacepage.PendingDecision{}
+	for _, d := range out.Decisions {
+		by[d.Request] = d
+	}
+	if len(out.Decisions) != 4 || len(by) != 4 {
+		t.Fatalf("decisions = %+v, want the decidable, the reserved, the shared and the uncarried one", out.Decisions)
+	}
+	want := surfacepage.PendingDecision{Request: q1.String(), Packet: p1.Hash.String(), Effect: broker.EffectHumanApprove, Scope: broker.WorkspaceScope,
+		Cost: 1, Requester: "agent-1", Deadline: packetNow + 60}
+	if by[q1.String()] != want {
+		t.Errorf("decidable: %+v, want %+v", by[q1.String()], want)
+	}
+	if got := by[qr.String()]; got.Packet != pr.Hash.String() || got.Reserved != "publish" || !strings.Contains(got.Refusal, "world-publish approve") {
+		t.Errorf("reserved: %+v", got)
+	}
+	if got := by[qs.String()]; got.Packet != "" || got.Deadline != 0 || !strings.Contains(got.Refusal, "is carried by 2 decision packets") {
+		t.Errorf("shared: %+v", got)
+	}
+	if got := by[qu.String()]; got.Packet != pu.Hash.String() || !strings.Contains(got.Refusal, "the decision packet carrying request "+qu.String()+" is refused: ") {
+		t.Errorf("uncarried: %+v", got)
+	}
+	for _, q := range []string{qd.String(), qz.String()} {
+		if _, ok := by[q]; ok {
+			t.Errorf("%s is listed", q)
+		}
+	}
+	r.d.approveNow = func() int64 { return packetNow + 60 }
+	out, _ = r.requestsJSON()
+	for _, d := range out.Decisions {
+		if d.Request == q1.String() && !strings.Contains(d.Refusal, "the packet's deadline passed") {
+			t.Errorf("past the deadline: %+v", d)
+		}
+	}
+	if len(out.Decisions) != 4 {
+		t.Errorf("past the deadline: %+v", out.Decisions)
+	}
+}
+
+// TestSurfaceRequestsDecisionsTruncated: more than 100 stored packets, or more than 20 requests on
+// the approvals chain, mark the list truncated.
+func TestSurfaceRequestsDecisionsTruncated(t *testing.T) {
+	t.Run("packet-scan", func(t *testing.T) {
+		r := newDecideRig(t)
+		r.packet("scan-ok", broker.WorkspaceScope, "agent")
+		r.fillPackets("scan", 100)
+		if out, _ := r.requestsJSON(); !out.DecisionsTruncated {
+			t.Errorf("101 packets: not truncated: %+v", out)
+		}
+	})
+	t.Run("packet-scan-at-the-bound", func(t *testing.T) {
+		r := newDecideRig(t)
+		r.packet("bound-ok", broker.WorkspaceScope, "agent")
+		r.fillPackets("bound", 99)
+		if out, _ := r.requestsJSON(); out.DecisionsTruncated || len(out.Decisions) != 1 {
+			t.Errorf("100 packets: %+v", out)
+		}
+	})
+	t.Run("approvals-chain", func(t *testing.T) {
+		r := newDecideRig(t)
+		for i := 0; i < 21; i++ {
+			approvalRequestInChain(t, r.d, broker.WorkspaceScope, "agent-"+string(rune('a'+i)), packetNow)
+		}
+		if out, _ := r.requestsJSON(); !out.DecisionsTruncated {
+			t.Errorf("21 requests: not truncated: %+v", out)
+		}
+	})
+}
+
+// TestSurfaceRequestsDecisionsUnreadable: a world store that cannot be read lists no decision (an
+// empty list, never null) and keeps the mint lists.
+func TestSurfaceRequestsDecisionsUnreadable(t *testing.T) {
+	r := newDecideRig(t)
+	r.packet("unreadable", broker.WorkspaceScope, "agent")
+	if err := r.d.store.Close(); err != nil {
+		t.Fatal(err)
+	}
+	_, raw := r.requestsJSON()
+	if !strings.Contains(raw, `"decisions":[]`) || !strings.Contains(raw, `"pending":[]`) || !strings.Contains(raw, `"decisionsTruncated":false`) {
+		t.Fatalf("unreadable world store: %s", raw)
+	}
+}
diff --git a/host/surfacepage/approve_model.go b/host/surfacepage/approve_model.go
index e4e490b..2c32b1b 100644
--- a/host/surfacepage/approve_model.go
+++ b/host/surfacepage/approve_model.go
@@ -67,6 +67,9 @@ type RequestsJSON struct {
 	Pending     []PendingRequest `json:"pending"`
 	Decided     []DecidedRequest `json:"decided"`
 	Unavailable string           `json:"unavailable"`
+	// Decisions are the pending Human.Approve requests a decision card decides (row 201 A4 §2.8).
+	Decisions          []PendingDecision `json:"decisions"`
+	DecisionsTruncated bool              `json:"decisionsTruncated"`
 }
 
 // hex32 decodes 64 hex characters.
diff --git a/host/surfacepage/decide_list.go b/host/surfacepage/decide_list.go
new file mode 100644
index 0000000..720e019
--- /dev/null
+++ b/host/surfacepage/decide_list.go
@@ -0,0 +1,56 @@
+package surfacepage
+
+// Row 201 A4 sprint 4 (design_docs/planned/w-world-approvals-by-key-a4.md §2.8): the decisions
+// GET /surface/requests.json lists under the mint requests, and approve --list prints. Pure: the
+// daemon reads the approvals chain and one bounded packet scan, computes each row's display verdict
+// with the decision law, and this fold only selects and shapes the rows. Display only: the POST
+// re-runs the law with the complete discovery.
+
+// decisionEffect is the only effect a decision card decides.
+const decisionEffect = "Human.Approve"
+
+// PendingDecision is one line of requests.json's decisions: a pending Human.Approve request and the
+// packet carrying it. Packet is "" and Deadline 0 unless exactly one scanned packet carries it;
+// Reserved names a typed owner ("publish", "finish"); Refusal is the law's reason at now, "" when
+// the request can be decided.
+type PendingDecision struct {
+	Request   string `json:"request"`
+	Packet    string `json:"packet"`
+	Effect    string `json:"effect"`
+	Scope     string `json:"scope"`
+	Cost      int64  `json:"cost"`
+	Requester string `json:"requester"`
+	Deadline  int64  `json:"deadline"`
+	Reserved  string `json:"reserved"`
+	Refusal   string `json:"refusal"`
+}
+
+// DecisionCandidate is one request of the approvals chain as the daemon read it: its summary, the
+// refs of the scanned packets carrying it (scan order), and, for those it computed them for, the one
+// packet's deadline, the reserved owner and the law's display reason.
+type DecisionCandidate struct {
+	Request, Effect, Scope, Requester string
+	Cost                              int64
+	Decided                           bool
+	Packets                           []string
+	Deadline                          int64
+	Reserved, Refusal                 string
+}
+
+// FoldDecisions is the decisions list (design §2.8): pending Human.Approve requests only; a request
+// no scanned packet carries is omitted; one carried by two or more is listed without a packet, with
+// its refusal. truncated is set when the approvals read or the packet scan hit its bound.
+func FoldDecisions(cands []DecisionCandidate, chainTruncated, scanTruncated bool) (rows []PendingDecision, truncated bool) {
+	for _, c := range cands {
+		if c.Decided || c.Effect != decisionEffect || len(c.Packets) == 0 {
+			continue
+		}
+		row := PendingDecision{Request: c.Request, Effect: c.Effect, Scope: c.Scope, Cost: c.Cost, Requester: c.Requester,
+			Reserved: c.Reserved, Refusal: c.Refusal}
+		if len(c.Packets) == 1 {
+			row.Packet, row.Deadline = c.Packets[0], c.Deadline
+		}
+		rows = append(rows, row)
+	}
+	return rows, chainTruncated || scanTruncated
+}
diff --git a/host/surfacepage/decide_list_test.go b/host/surfacepage/decide_list_test.go
new file mode 100644
index 0000000..3253e5f
--- /dev/null
+++ b/host/surfacepage/decide_list_test.go
@@ -0,0 +1,49 @@
+package surfacepage
+
+// Row 201 A4 sprint 4 (design §6 DMD1–DMD3): the decisions fold, without a daemon or a store.
+
+import (
+	"reflect"
+	"testing"
+)
+
+// DMD1: pending Human.Approve requests only; decided, other effects and unpacketed requests are
+// omitted; one packet gives the packet and its deadline, the reserved owner and refusal pass through.
+func TestDecideModelListsPendingDecisionsOnly(t *testing.T) {
+	cands := []DecisionCandidate{
+		{Request: "sha256:01", Effect: "Human.Approve", Scope: "workspace:w", Requester: "agent-1", Cost: 3, Packets: []string{"sha256:p1"}, Deadline: 1760},
+		{Request: "sha256:02", Effect: "Human.Approve", Scope: "workspace:w", Requester: "agent-2", Decided: true, Packets: []string{"sha256:p2"}, Deadline: 1760},
+		{Request: "sha256:03", Effect: "Registry.Publish", Scope: "registry:x", Requester: "agent-3", Packets: []string{"sha256:p3"}, Deadline: 1760},
+		{Request: "sha256:04", Effect: "Human.Approve", Scope: "workspace:w", Requester: "agent-4"},
+		{Request: "sha256:05", Effect: "Human.Approve", Scope: "registry:pkg#effect=Registry.Publish", Requester: "agent-5", Packets: []string{"sha256:p5"},
+			Deadline: 1770, Reserved: "publish", Refusal: "publish approvals are minted by `world-publish approve`"},
+	}
+	rows, truncated := FoldDecisions(cands, false, false)
+	want := []PendingDecision{
+		{Request: "sha256:01", Packet: "sha256:p1", Effect: "Human.Approve", Scope: "workspace:w", Cost: 3, Requester: "agent-1", Deadline: 1760},
+		{Request: "sha256:05", Packet: "sha256:p5", Effect: "Human.Approve", Scope: "registry:pkg#effect=Registry.Publish", Requester: "agent-5", Deadline: 1770,
+			Reserved: "publish", Refusal: "publish approvals are minted by `world-publish approve`"},
+	}
+	if !reflect.DeepEqual(rows, want) || truncated {
+		t.Fatalf("FoldDecisions = %+v, %v; want %+v, false", rows, truncated, want)
+	}
+}
+
+// DMD2: a request two packets carry is listed with its refusal and no packet or deadline.
+func TestDecideModelListsATwoPacketRequestRefused(t *testing.T) {
+	rows, _ := FoldDecisions([]DecisionCandidate{{Request: "sha256:06", Effect: "Human.Approve", Scope: "s", Requester: "a",
+		Packets: []string{"sha256:pa", "sha256:pb"}, Deadline: 1760, Refusal: "request sha256:06 is carried by 2 decision packets"}}, false, false)
+	want := []PendingDecision{{Request: "sha256:06", Effect: "Human.Approve", Scope: "s", Requester: "a", Refusal: "request sha256:06 is carried by 2 decision packets"}}
+	if !reflect.DeepEqual(rows, want) {
+		t.Fatalf("two packets: %+v, want %+v", rows, want)
+	}
+}
+
+// DMD3: truncated is either bound: the approvals read's or the packet scan's.
+func TestDecideModelTruncation(t *testing.T) {
+	for _, c := range []struct{ chain, scan, want bool }{{false, false, false}, {true, false, true}, {false, true, true}, {true, true, true}} {
+		if _, got := FoldDecisions(nil, c.chain, c.scan); got != c.want {
+			t.Errorf("chain %v scan %v: truncated %v, want %v", c.chain, c.scan, got, c.want)
+		}
+	}
+}
````

### Patch M7c — approve --list decisions (design Q2)

````diff M7c
diff --git a/cmd/ailang-worldd/approve.go b/cmd/ailang-worldd/approve.go
index 1754240..cc6e8c5 100644
--- a/cmd/ailang-worldd/approve.go
+++ b/cmd/ailang-worldd/approve.go
@@ -272,12 +272,12 @@ func runApprove(args []string, stdout, stderr io.Writer, env sessionEnv) int {
 		}
 		if len(rows) == 0 {
 			fmt.Fprintln(stdout, "no pending mint requests")
-			return exitOK
 		}
 		for _, r := range rows {
 			fmt.Fprintf(stdout, "%s  episode %s  digest %s  expires %s\n", r.RequestID, r.EpisodeID, r.Digest[:16],
 				time.Unix(r.ExpiresAt, 0).UTC().Format(time.RFC3339))
 		}
+		approveListDecisions(ctx, *db, env.now(), stdout, stderr) // row 201 A4 (design Q2)
 		return exitOK
 	}
 	row, ok, err := ro.MintRequest(ctx, id)
diff --git a/cmd/ailang-worldd/approve_list_decisions.go b/cmd/ailang-worldd/approve_list_decisions.go
new file mode 100644
index 0000000..866113d
--- /dev/null
+++ b/cmd/ailang-worldd/approve_list_decisions.go
@@ -0,0 +1,55 @@
+package main
+
+// Row 201 A4 sprint 4 (design_docs/planned/w-world-approvals-by-key-a4.md §2.10, Q2): approve --list
+// also prints the pending Human.Approve requests an operator decides with a request ref. It reads
+// the world store --db read-only, offline, through the daemon's own listing (daemon.PendingDecisions,
+// the rows GET /surface/requests.json serves), so it never dials a daemon.
+
+import (
+	"context"
+	"fmt"
+	"io"
+
+	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// approveListDecisions prints the decisions section of approve --list: nothing when there is none;
+// a world store that cannot be read is one line on stderr and changes neither the mint list nor the
+// exit code.
+func approveListDecisions(ctx context.Context, db string, now int64, stdout, stderr io.Writer) {
+	ro, err := store.OpenReadOnly(db)
+	if err != nil {
+		fmt.Fprintf(stderr, "ailang-worldd approve: decisions: cannot read the world store %s: %v\n", db, err)
+		return
+	}
+	rows, truncated, err := daemon.PendingDecisions(ctx, ro, now)
+	_ = ro.Close()
+	if err != nil {
+		fmt.Fprintf(stderr, "ailang-worldd approve: decisions: cannot list the decisions in %s: %v\n", db, err)
+		return
+	}
+	printDecisions(stdout, rows, truncated)
+}
+
+// printDecisions is the section's text: nothing for no row; one line per request, its refusal or
+// its deadline and the command that decides it; the bound when the list was cut.
+func printDecisions(stdout io.Writer, rows []surfacepage.PendingDecision, truncated bool) {
+	if len(rows) == 0 && !truncated {
+		return
+	}
+	fmt.Fprintln(stdout, "pending decisions (Human.Approve requests):")
+	for _, d := range rows {
+		line := fmt.Sprintf("%s  %s  requested by %s", d.Request, d.Scope, d.Requester)
+		if d.Refusal != "" {
+			line += "  refused: " + d.Refusal
+		} else {
+			line += fmt.Sprintf("  deadline %s  decide: ailang-worldd approve %s [--deny]", approverTime(d.Deadline), d.Request)
+		}
+		fmt.Fprintln(stdout, line)
+	}
+	if truncated {
+		fmt.Fprintln(stdout, "(more requests than listed: the approvals chain or the packet scan reached its bound)")
+	}
+}
diff --git a/cmd/ailang-worldd/approve_list_decisions_test.go b/cmd/ailang-worldd/approve_list_decisions_test.go
new file mode 100644
index 0000000..eea27bb
--- /dev/null
+++ b/cmd/ailang-worldd/approve_list_decisions_test.go
@@ -0,0 +1,93 @@
+package main
+
+// Row 201 A4 sprint 4 (design §2.10, Q2): approve --list prints the pending decisions of the world
+// store --db, read-only and offline (it never dials a daemon). Temp stores only (D-WORLD-109).
+
+import (
+	"bytes"
+	"context"
+	"path/filepath"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/coordinator"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// TestApproveListShowsDecisions: a Human.Approve request on the --db store's approvals chain whose
+// one stored packet is not carried by a commit is listed with the law's refusal; the mint list is
+// printed as before; a --db that cannot be read is one stderr line and exit 0.
+func TestApproveListShowsDecisions(t *testing.T) {
+	db, creds := elsewhereStores(t)
+	putPendingRequest(t, creds)
+	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
+	defer cancel()
+	st, err := store.Open(db)
+	if err != nil {
+		t.Fatal(err)
+	}
+	out, err := broker.NewHumanHandler(st).Execute(ctx, broker.EffectRequest{Effect: broker.EffectHumanApprove, Scope: broker.WorkspaceScope, Cost: 1, Now: 1700},
+		[]byte(`{"requester":"agent-list"}`))
+	if err != nil {
+		t.Fatal(err)
+	}
+	q, err := broker.ParsePendingApproval(out)
+	if err != nil {
+		t.Fatal(err)
+	}
+	payload := coordinator.EncodeDecisionPacket(coordinator.DecisionPacket{ProposalHash: hashref.SumSHA256([]byte("plan")), RequestRef: q,
+		CreatedAt: 1700, DeadlineAt: 1760, Policy: "Cancel"})
+	if err := st.PutObject(ctx, store.Object{Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte(coordinator.DecisionPacketV1)),
+		SemanticID: coordinator.DecisionPacketV1, Provenance: "list-test", Payload: payload}); err != nil {
+		t.Fatal(err)
+	}
+	if err := st.Close(); err != nil {
+		t.Fatal(err)
+	}
+	var stdout, stderr bytes.Buffer
+	if code := runApprove([]string{"--list", "--db", db, "--credentials-db", creds}, &stdout, &stderr, approveEnv(nil, nil, atTime(approveT0+1))); code != exitOK {
+		t.Fatalf("rc %d, stderr %q", code, stderr.String())
+	}
+	want := firstRequestID + "  episode ep1  digest abababababababab  expires 2027-01-15T08:10:00Z\n" +
+		"pending decisions (Human.Approve requests):\n" +
+		q.String() + "  " + broker.WorkspaceScope + "  requested by agent-list  refused: the decision packet carrying request " + q.String() + " is refused: "
+	if got := stdout.String(); !strings.HasPrefix(got, want) || strings.Count(got, "\n") != 3 || stderr.Len() != 0 {
+		t.Fatalf("stdout %q\nwant the prefix %q; stderr %q", got, want, stderr.String())
+	}
+	stdout.Reset()
+	missing := filepath.Join(t.TempDir(), "missing-world.db")
+	if code := runApprove([]string{"--list", "--db", missing, "--credentials-db", creds}, &stdout, &stderr, approveEnv(nil, nil, atTime(approveT0+1))); code != exitOK ||
+		!strings.HasPrefix(stdout.String(), firstRequestID) || strings.Contains(stdout.String(), "decisions") ||
+		!strings.Contains(stderr.String(), "ailang-worldd approve: decisions: cannot") {
+		t.Fatalf("missing world store: rc %d, stdout %q, stderr %q", code, stdout.String(), stderr.String())
+	}
+}
+
+// TestApproveListDecisionLines: the section's text for a decidable, a refused and a cut list.
+func TestApproveListDecisionLines(t *testing.T) {
+	var b bytes.Buffer
+	printDecisions(&b, nil, false)
+	if b.Len() != 0 {
+		t.Fatalf("no decision printed %q", b.String())
+	}
+	printDecisions(&b, []surfacepage.PendingDecision{
+		{Request: "sha256:aa", Packet: "sha256:pp", Scope: "workspace:w", Requester: "agent-1", Deadline: 1760},
+		{Request: "sha256:bb", Scope: "registry:x", Requester: "agent-2", Reserved: "publish", Refusal: "publish approvals are minted by `world-publish approve` at a terminal (R4)"},
+	}, true)
+	want := "pending decisions (Human.Approve requests):\n" +
+		"sha256:aa  workspace:w  requested by agent-1  deadline 1970-01-01T00:29:20Z  decide: ailang-worldd approve sha256:aa [--deny]\n" +
+		"sha256:bb  registry:x  requested by agent-2  refused: publish approvals are minted by `world-publish approve` at a terminal (R4)\n" +
+		"(more requests than listed: the approvals chain or the packet scan reached its bound)\n"
+	if b.String() != want {
+		t.Fatalf("got\n%s\nwant\n%s", b.String(), want)
+	}
+	b.Reset()
+	printDecisions(&b, nil, true)
+	if b.String() != "pending decisions (Human.Approve requests):\n(more requests than listed: the approvals chain or the packet scan reached its bound)\n" {
+		t.Fatalf("cut and empty: %q", b.String())
+	}
+}
````

### Patch M7d — Needs you lists decisions (XD3)

````diff M7d
diff --git a/host/surfacepage/app/app.js b/host/surfacepage/app/app.js
index 1e6fdf9..39b4502 100644
--- a/host/surfacepage/app/app.js
+++ b/host/surfacepage/app/app.js
@@ -424,13 +424,28 @@
   // needsYouHTML is the now card's Needs you (row 201 A3, design §2.8): the pending mint requests of
   // GET /surface/requests.json, each a link to its card; a request the daemon could not read shows why;
   // then the notification opt-in for this browser's permission (design §2.9).
+  // Row 201 A4 (design §2.8): then the pending Human.Approve decisions of the same read, each a link to
+  // its decision card; a reserved one names its typed owner and a refused one its reason, with no link.
   function needsYouHTML(r, now, permission) {
     const p = r && Array.isArray(r.pending) ? r.pending : [];
-    const head = `<div class="nh" style="color:${p.length ? 'var(--youShell)' : 'var(--sub)'}">Needs you · ${p.length}</div>`;
-    if (!p.length) return `<div class="nsec" data-needs>${head}<div class="nl sub">${esc(r && r.unavailable || 'Nothing needs you')}</div>${notifyHTML(permission)}</div>`;
+    const ds = r && Array.isArray(r.decisions) ? r.decisions : [];
+    const n = p.length + ds.length;
+    const head = `<div class="nh" style="color:${n ? 'var(--youShell)' : 'var(--sub)'}">Needs you · ${n}</div>`;
+    const more = r && r.decisionsTruncated ? '<div class="nl sub">Older approval requests are not listed; decide them at a terminal with ailang-worldd approve &lt;request-ref&gt;</div>' : '';
+    if (!n) return `<div class="nsec" data-needs>${head}<div class="nl sub">${esc(r && r.unavailable || 'Nothing needs you')}</div>${more}${notifyHTML(permission)}</div>`;
     return `<div class="nsec" data-needs>${head}` + p.map(q => q.refused ? `<div class="nl sub">${esc(q.id)} · refused: ${esc(q.refused)}</div>`
       : `<div class="nl"><a id="needs-${esc(q.id)}" href="/surface/approve/${esc(q.id)}"><b>${esc(q.episode)}</b> · ${esc(q.agent)} · ${q.grants} grant${q.grants === 1 ? '' : 's'} · expires in ${Math.max(0, Math.ceil((q.expires - now) / 60))} min</a>`
-        + (q.recorded ? ` · approved by ${esc(q.recorded)}, waiting for the minter` : '') + '</div>').join('') + notifyHTML(permission) + '</div>';
+        + (q.recorded ? ` · approved by ${esc(q.recorded)}, waiting for the minter` : '') + '</div>').join('') + ds.map(decisionLineHTML).join('')
+      + more + notifyHTML(permission) + '</div>';
+  }
+  // RESERVED_OWNERS is the typed command that owns a reserved decision (design §2.6 L1).
+  const RESERVED_OWNERS = { publish: 'world-publish approve', finish: 'ailang-worldd episode finish' };
+  // decisionLineHTML is one decision's Needs-you line (design §2.8).
+  function decisionLineHTML(x) {
+    if (x.reserved) return `<div class="nl sub">${esc(x.reserved)} approval — use ${esc(RESERVED_OWNERS[x.reserved] || 'its typed owner')} at a terminal</div>`;
+    if (x.refusal) return `<div class="nl sub">Decide · ${esc(x.scope)} · refused: ${esc(x.refusal)}</div>`;
+    const ref = String(x.request || '').replace(/^sha256:/, '');
+    return `<div class="nl"><a id="decide-${esc(ref)}" href="/surface/decide/sha256:${esc(ref)}">Decide · ${esc(x.scope)} · requested by ${esc(x.requester)} · deadline ${hm(new Date(x.deadline * 1000))}</a></div>`;
   }
   // notifyHTML is the opt-in: "Notify me" while this browser has not been asked (permission "default"),
   // its answer once it has, and nothing where the browser has no notifications.
@@ -443,11 +458,16 @@
   // for each pending request this tab has not seen, while notifications are granted; each opens its card.
   // seen is null before the page's first read, which raises nothing (the page itself shows those); a read
   // that failed changes nothing. Refused and already-recorded requests are seen but raise nothing.
+  // A decision (row 201 A4) is keyed by its request ref; a reserved or refused one is seen but raises nothing.
   function notifications(r, seen, permission) {
     if (!r || !Array.isArray(r.pending)) return { notes: [], seen };
+    const ds = Array.isArray(r.decisions) ? r.decisions : [];
     const fresh = seen ? r.pending.filter(q => !seen.has(q.id) && !q.refused && !q.recorded) : [];
-    return { seen: new Set([...(seen || []), ...r.pending.map(q => q.id)]), notes: permission !== 'granted' ? [] : fresh.map(q => ({ title: 'World needs you', tag: q.id,
-      body: `Approve mint request for ${q.episode} (${q.grants} grant${q.grants === 1 ? '' : 's'}), expires ${hm(new Date(q.expires * 1000))}`, href: '/surface/approve/' + q.id })) };
+    const freshD = seen ? ds.filter(x => !seen.has(x.request) && !x.reserved && !x.refusal) : [];
+    return { seen: new Set([...(seen || []), ...r.pending.map(q => q.id), ...ds.map(x => x.request)]), notes: permission !== 'granted' ? [] : fresh.map(q => ({ title: 'World needs you', tag: q.id,
+      body: `Approve mint request for ${q.episode} (${q.grants} grant${q.grants === 1 ? '' : 's'}), expires ${hm(new Date(q.expires * 1000))}`, href: '/surface/approve/' + q.id }))
+      .concat(freshD.map(x => ({ title: 'World needs you', tag: x.request,
+        body: `Decide ${x.scope} requested by ${x.requester}, deadline ${hm(new Date(x.deadline * 1000))}`, href: '/surface/decide/' + x.request }))) };
   }
   // decidedHTML is the decision list under Needs you (design §2.8): the last few of GET
   // /surface/requests.json's decided requests (newest first; the daemon sends 20), and all of them on request.
diff --git a/host/surfacepage/app_decide_test.go b/host/surfacepage/app_decide_test.go
new file mode 100644
index 0000000..e1d1b7f
--- /dev/null
+++ b/host/surfacepage/app_decide_test.go
@@ -0,0 +1,102 @@
+package surfacepage
+
+import (
+	"context"
+	"os"
+	"os/exec"
+	"path/filepath"
+	"strings"
+	"testing"
+	"time"
+)
+
+// TestNeedsYouListsDecisions (row 201 A4 XD3, design §2.8) runs app.js's pure part under node: Needs
+// you lists GET /surface/requests.json's decisions after the mint requests, a decidable one as a link
+// to its card, a reserved one naming its typed owner and a refused one its reason, neither with a
+// link; every value escaped; a notification for each new decidable decision. It skips only when no
+// node exists; with WORLD_EXEC_NODE set, a harness that cannot run fails.
+func TestNeedsYouListsDecisions(t *testing.T) {
+	node := os.Getenv("WORLD_EXEC_NODE")
+	if node == "" {
+		var err error
+		node, err = exec.LookPath("node")
+		if err != nil {
+			t.Skip("node unavailable locally")
+		}
+	}
+	dir := t.TempDir()
+	for name, data := range map[string][]byte{"app.js": AppScript, "harness.js": []byte(needsYouDecisionsHarness)} {
+		if err := os.WriteFile(filepath.Join(dir, name), data, 0o600); err != nil {
+			t.Fatal(err)
+		}
+	}
+	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
+	defer cancel()
+	cmd := exec.CommandContext(ctx, node, filepath.Join(dir, "harness.js"))
+	cmd.Dir = dir
+	out, err := cmd.CombinedOutput()
+	if err != nil {
+		t.Fatalf("node harness: %v\n%s", err, out)
+	}
+	if !strings.Contains(string(out), "NEEDS YOU DECISIONS PASS") {
+		t.Fatalf("harness manifest missing: %s", out)
+	}
+}
+
+const needsYouDecisionsHarness = `
+const assert = require('node:assert/strict');
+const app = require('./app.js');
+const Q = 'ab'.repeat(32), R = 'cd'.repeat(32), T = 'ef'.repeat(32), U = '01'.repeat(32);
+const ok = {request: 'sha256:' + Q, packet: 'sha256:' + R, effect: 'Human.Approve', scope: 'worktree:1', cost: 1, requester: 'agent-1', deadline: 50520, reserved: '', refusal: ''};
+const pub = {request: 'sha256:' + R, packet: 'sha256:' + T, effect: 'Human.Approve', scope: 'registry:pkg', cost: 1, requester: 'agent-2', deadline: 50520, reserved: 'publish', refusal: 'publish approvals are minted by world-publish approve'};
+const fin = {request: 'sha256:' + U, packet: '', effect: 'Human.Approve', scope: 'github-pr:1', cost: 1, requester: 'agent-4', deadline: 0, reserved: 'finish', refusal: 'x'};
+const two = {request: 'sha256:' + T, packet: '', effect: 'Human.Approve', scope: 'worktree:2', cost: 1, requester: 'agent-3', deadline: 0, reserved: '', refusal: 'carried by 2 decision packets'};
+const line = '<div class="nl"><a id="decide-' + Q + '" href="/surface/decide/sha256:' + Q + '">Decide · worktree:1 · requested by agent-1 · deadline 14:02</a></div>';
+// A decidable decision: one line, a link to its card; the count includes it; no mint request needed.
+assert.equal(app.needsYouHTML({now: 1000, pending: [], decisions: [ok], decisionsTruncated: false}, 1000),
+  '<div class="nsec" data-needs><div class="nh" style="color:var(--youShell)">Needs you · 1</div>' + line + '</div>');
+// After the mint requests, in order; reserved names its owner with no link; refused shows its reason with no link.
+{ const h = app.needsYouHTML({now: 0, pending: [{id: 'mr-1', episode: 'e', agent: 'a', grants: 1, expires: 60}], decisions: [ok, pub, fin, two]}, 0);
+  assert(h.startsWith('<div class="nsec" data-needs><div class="nh" style="color:var(--youShell)">Needs you · 5</div><div class="nl"><a id="needs-mr-1"'), h);
+  assert(h.includes(' min</a></div>' + line + '<div class="nl sub">publish approval — use world-publish approve at a terminal</div>'
+    + '<div class="nl sub">finish approval — use ailang-worldd episode finish at a terminal</div>'
+    + '<div class="nl sub">Decide · worktree:2 · refused: carried by 2 decision packets</div></div>'), h);
+  for (const ref of [R, U, T]) assert(!h.includes('/surface/decide/sha256:' + ref) && !h.includes('decide-' + ref), h);
+  assert.equal(h.split('<a ').length, 3, h); }
+// An unknown reserved owner still renders no link.
+assert.equal(app.needsYouHTML({pending: [], decisions: [{...ok, reserved: 'x<y'}]}, 0),
+  '<div class="nsec" data-needs><div class="nh" style="color:var(--youShell)">Needs you · 1</div><div class="nl sub">x&lt;y approval — use its typed owner at a terminal</div></div>');
+// Truncation is said once, after the lines.
+assert.equal(app.needsYouHTML({pending: [], decisions: [ok], decisionsTruncated: true}, 0, 'granted'),
+  '<div class="nsec" data-needs><div class="nh" style="color:var(--youShell)">Needs you · 1</div>' + line
+  + '<div class="nl sub">Older approval requests are not listed; decide them at a terminal with ailang-worldd approve &lt;request-ref&gt;</div><div class="nl sub" data-notify>Notifications on</div></div>');
+// A truncated read with nothing listed still says so.
+assert.equal(app.needsYouHTML({pending: [], decisions: [], decisionsTruncated: true}, 0),
+  '<div class="nsec" data-needs><div class="nh" style="color:var(--sub)">Needs you · 0</div><div class="nl sub">Nothing needs you</div><div class="nl sub">Older approval requests are not listed; decide them at a terminal with ailang-worldd approve &lt;request-ref&gt;</div></div>');
+// Every value escaped where it lands: the ref in the id and href, the scope, the requester, the refusal.
+{ const bad = '"><img src=x onerror=1>';
+  assert.equal(app.needsYouHTML({pending: [], decisions: [{...ok, request: 'sha256:' + bad, scope: '<img src=x onerror=2>', requester: '<b>3</b>'}]}, 0),
+    '<div class="nsec" data-needs><div class="nh" style="color:var(--youShell)">Needs you · 1</div><div class="nl"><a id="decide-&quot;&gt;&lt;img src=x onerror=1&gt;" href="/surface/decide/sha256:&quot;&gt;&lt;img src=x onerror=1&gt;">'
+    + 'Decide · &lt;img src=x onerror=2&gt; · requested by &lt;b&gt;3&lt;/b&gt; · deadline 14:02</a></div></div>');
+  assert.equal(app.needsYouHTML({pending: [], decisions: [{...two, scope: '<s>', refusal: '<i>4</i>'}]}, 0),
+    '<div class="nsec" data-needs><div class="nh" style="color:var(--youShell)">Needs you · 1</div><div class="nl sub">Decide · &lt;s&gt; · refused: &lt;i&gt;4&lt;/i&gt;</div></div>'); }
+// Without decisions (or with an empty list), Needs you is exactly as before.
+assert.equal(app.needsYouHTML({pending: [], decisions: []}, 0), '<div class="nsec" data-needs><div class="nh" style="color:var(--sub)">Needs you · 0</div><div class="nl sub">Nothing needs you</div></div>');
+{ const r = {now: 1000, pending: [{id: 'mr-0123456789abcdef', episode: 'w<1>', agent: 'claude', grants: 2, expires: 1601, recorded: 'webauthn:mark'}]};
+  assert.equal(app.needsYouHTML({...r, decisions: [], decisionsTruncated: false}, 1000, 'default'), app.needsYouHTML(r, 1000, 'default')); }
+// Notifications: one per new decidable decision while granted; none on the first read, none twice, none for a reserved or refused one.
+{ const r1 = {pending: [], decisions: [ok]}, r2 = {pending: [{id: 'mr-a', episode: 'e1', grants: 1, expires: 50520}], decisions: [ok, pub, two, {...ok, request: 'sha256:' + U, scope: 's<2>', requester: 'agent-9', deadline: 50580}]};
+  const first = app.notifications(r1, null, 'granted');
+  assert.deepEqual(first.notes, []); assert.deepEqual([...first.seen], ['sha256:' + Q]);
+  const second = app.notifications(r2, first.seen, 'granted');
+  assert.deepEqual(second.notes, [{title: 'World needs you', tag: 'mr-a', body: 'Approve mint request for e1 (1 grant), expires 14:02', href: '/surface/approve/mr-a'},
+    {title: 'World needs you', tag: 'sha256:' + U, body: 'Decide s<2> requested by agent-9, deadline 14:03', href: '/surface/decide/sha256:' + U}]);
+  assert.deepEqual([...second.seen], ['sha256:' + Q, 'mr-a', 'sha256:' + R, 'sha256:' + T, 'sha256:' + U]);
+  assert.deepEqual(app.notifications(r2, second.seen, 'granted').notes, []);
+  assert.deepEqual(app.notifications(r2, first.seen, 'default').notes, []);
+  assert.deepEqual(app.notifications({pending: [], decisions: [ok]}, new Set(), 'granted').notes, [{title: 'World needs you', tag: 'sha256:' + Q, body: 'Decide worktree:1 requested by agent-1, deadline 14:02', href: '/surface/decide/sha256:' + Q}]);
+  assert.deepEqual(app.notifications({pending: [], decisions: [pub, two]}, new Set(), 'granted').notes, []);
+  assert.deepEqual(app.notifications({pending: [], decisions: [{...ok, reserved: 'publish'}]}, new Set(), 'granted').notes, []); // reserved alone, refusal empty
+  assert.deepEqual(app.notifications({pending: []}, new Set(), 'granted'), {notes: [], seen: new Set()}); }
+console.log('NEEDS YOU DECISIONS PASS');
+`
````

### Patch M7e — DecisionSubject.Passkeys

````diff M7e
diff --git a/host/daemon/surface_decide.go b/host/daemon/surface_decide.go
index 2e36d26..a0ab89c 100644
--- a/host/daemon/surface_decide.go
+++ b/host/daemon/surface_decide.go
@@ -115,13 +115,37 @@ func (d *Daemon) handleSurfaceDecision(w http.ResponseWriter, r *http.Request, _
 	}
 	p, _ := coordinator.DecodeDecisionPacket(rd.packet.Payload) // validated: code 0
 	reserved, _ := broker.ReservedApprovalScope(rd.facts.Scope)
+	passkeys, err := d.decisionPasskeys(ctx)
+	if err != nil {
+		d.approvalInternal(w, r, err)
+		return
+	}
 	writeApprovalJSON(w, http.StatusOK, surfacepage.DecisionSubject{V: 1,
 		Request: ref, RequestBytes: base64.RawURLEncoding.EncodeToString(rd.request.Payload),
 		Packet: rd.packet.Hash.String(), PacketBytes: base64.RawURLEncoding.EncodeToString(rd.packet.Payload),
 		Entry: rd.entry, Proposal: p.ProposalHash.String(), Now: now,
 		Window:  surfacepage.DecisionWindow{Min: now + 1, Max: broker.DecisionWindowMax(now, p.DeadlineAt)},
 		Refusal: law.Reason, Reserved: reserved, DecidedBy: rd.facts.DecidedBy, Decision: rd.facts.DecidedVerb,
-		SSH: surfacepage.DecisionSSH{Namespace: approval.DecisionNamespace}})
+		Passkeys: passkeys, SSH: surfacepage.DecisionSSH{Namespace: approval.DecisionNamespace}})
+}
+
+// decisionPasskeys is the active passkeys by origin on a fresh read-only credentials handle, for
+// the card's ceremony (sprint 4); a store below schema 6 has none.
+func (d *Daemon) decisionPasskeys(ctx context.Context) ([]surfacepage.CardPasskeys, error) {
+	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
+	if err != nil {
+		return nil, err
+	}
+	defer func() { _ = ro.Close() }()
+	active, err := ro.ActivePasskeys(ctx)
+	if err != nil && !isOlderStore(err) {
+		return nil, err
+	}
+	var passkeys []surfacepage.CardPasskey
+	for _, e := range active {
+		passkeys = append(passkeys, surfacepage.CardPasskey{Origin: e.Passkey.Origin, RPID: e.Passkey.RPID, CredentialID: e.Passkey.CredentialID})
+	}
+	return surfacepage.DecisionPasskeys(passkeys), nil
 }
 
 // handleSurfaceDecide is POST /surface/decide/{ref} (design §2.8). After R2 (approvalRoute): R0 any
diff --git a/host/daemon/surface_decide_passkeys_test.go b/host/daemon/surface_decide_passkeys_test.go
new file mode 100644
index 0000000..ba38310
--- /dev/null
+++ b/host/daemon/surface_decide_passkeys_test.go
@@ -0,0 +1,56 @@
+package daemon
+
+// Row 201 A4 sprint 4 (design §2.8; sprint-1 plan P7): the subject's passkeys.
+
+import (
+	"context"
+	"encoding/base64"
+	"encoding/json"
+	"path/filepath"
+	"reflect"
+	"sort"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// The subject's passkeys (sprint-1 P7, sprint 4): the active passkeys grouped by origin, as the
+// request card's; none enrolled is an empty list; a credentials store that cannot be opened is the
+// sanitized 500, never a subject without its passkeys.
+func TestDecisionSubjectPasskeys(t *testing.T) {
+	r := newDecideRig(t)
+	q, _ := r.packet("passkeys", broker.WorkspaceScope, "agent")
+	subject := func() (int, surfacepage.DecisionSubject, string) {
+		rec := r.serve("GET", "/surface/decisions/"+q.String(), apLoopback, nil)
+		var s surfacepage.DecisionSubject
+		_ = json.Unmarshal(rec.Body.Bytes(), &s)
+		return rec.Code, s, rec.Body.String()
+	}
+	if code, s, body := subject(); code != 200 || s.Passkeys == nil || len(s.Passkeys) != 0 || !strings.Contains(body, `"passkeys":[]`) {
+		t.Fatalf("no passkey: %d %s", code, body)
+	}
+	a, b, c := approvaltest.NewPasskey(-7, 1), approvaltest.NewPasskey(-7, 2), approvaltest.NewPasskey(-7, 3)
+	r.enrolPasskey("mark-mac", a, "https://studio.tail.ts.net")
+	r.enrolPasskey("mark-phone", b, "http://localhost:7644")
+	r.enrolPasskey("mark-key", c, "https://studio.tail.ts.net")
+	r.write(func(ctx context.Context, cr *store.Credentials) error {
+		return cr.RevokeApprover(ctx, "mark-key", apT0)
+	})
+	code, s, body := subject()
+	got := map[string][]string{}
+	for _, g := range s.Passkeys {
+		sort.Strings(g.Allow)
+		got[g.Origin+" "+g.RP] = g.Allow
+	}
+	enc := func(p *approvaltest.Passkey) string { return base64.RawURLEncoding.EncodeToString(p.CredentialID()) }
+	want := map[string][]string{"https://studio.tail.ts.net studio.tail.ts.net": {enc(a)}, "http://localhost:7644 localhost": {enc(b)}}
+	if code != 200 || len(s.Passkeys) != 2 || !reflect.DeepEqual(got, want) {
+		t.Fatalf("passkeys %d %v, want %v\n%s", code, got, want, body)
+	}
+	r.d.credentialsPath = filepath.Join(t.TempDir(), "missing", "credentials.db")
+	requireRefusal(t, "no credentials store", r.serve("GET", "/surface/decisions/"+q.String(), apLoopback, nil), 500, "Internal")
+}
diff --git a/host/daemon/surface_decide_subject_test.go b/host/daemon/surface_decide_subject_test.go
index 43e731e..2f3be59 100644
--- a/host/daemon/surface_decide_subject_test.go
+++ b/host/daemon/surface_decide_subject_test.go
@@ -36,7 +36,7 @@ func TestDecisionSubjectJSON(t *testing.T) {
 		names = append(names, k)
 	}
 	sort.Strings(names)
-	if want := []string{"decidedBy", "decision", "entry", "now", "packet", "packetBytes", "proposal", "refusal", "request", "requestBytes", "reserved", "ssh", "v", "window"}; !reflect.DeepEqual(names, want) {
+	if want := []string{"decidedBy", "decision", "entry", "now", "packet", "packetBytes", "passkeys", "proposal", "refusal", "request", "requestBytes", "reserved", "ssh", "v", "window"}; !reflect.DeepEqual(names, want) {
 		t.Fatalf("subject keys %v, want %v (no message, no hash)", names, want)
 	}
 	var s surfacepage.DecisionSubject
diff --git a/host/surfacepage/decide_model.go b/host/surfacepage/decide_model.go
index 1639c0f..afd44af 100644
--- a/host/surfacepage/decide_model.go
+++ b/host/surfacepage/decide_model.go
@@ -29,9 +29,16 @@ type DecisionSubject struct {
 	Reserved     string         `json:"reserved"`
 	DecidedBy    string         `json:"decidedBy"`
 	Decision     string         `json:"decision"`
+	Passkeys     []CardPasskeys `json:"passkeys"` // the active passkeys by origin, as the request card's
 	SSH          DecisionSSH    `json:"ssh"`
 }
 
+// DecisionPasskeys is the active passkeys grouped by origin, in the order given: each origin's rpID
+// and credential IDs (base64url), as the request card offers them; never nil.
+func DecisionPasskeys(passkeys []CardPasskey) []CardPasskeys {
+	return RequestCard(StoredRequest{}, passkeys).Passkeys
+}
+
 // DecisionWindow is the expiry window of a decision made now: min = now + 1, max = min(now + 600, deadline).
 type DecisionWindow struct {
 	Min int64 `json:"min"`
````

### Patch M7f — the decision card: routes, page, decide.js (PD1, PD2)

````diff M7f
diff --git a/host/daemon/daemon.go b/host/daemon/daemon.go
index ec62186..d6cad9c 100644
--- a/host/daemon/daemon.go
+++ b/host/daemon/daemon.go
@@ -1142,6 +1142,8 @@ func (d *Daemon) Handler() http.Handler {
 	mux.HandleFunc("GET /surface/approve.js", d.approvalRoute(d.handleSurfaceApproveScript))
 	mux.HandleFunc("GET /surface/decisions/{ref}", d.approvalRoute(d.handleSurfaceDecision))
 	mux.HandleFunc("POST /surface/decide/{ref}", d.approvalRoute(d.handleSurfaceDecide))
+	mux.HandleFunc("GET /surface/decide/{ref}", d.approvalRoute(d.handleSurfaceDecidePage))
+	mux.HandleFunc("GET /surface/decide.js", d.approvalRoute(d.handleSurfaceDecideScript))
 	// The two A2A projection routes (w-a2a-session-projection P6.B-A2A-CARD)
 	// are ADDITIVE: the frozen /v1/ table above is untouched, and the routes
 	// are NOT in isProtected — the projection handler resolves the session
diff --git a/host/daemon/surface_app_test.go b/host/daemon/surface_app_test.go
index b3c6249..31cfe5c 100644
--- a/host/daemon/surface_app_test.go
+++ b/host/daemon/surface_app_test.go
@@ -154,12 +154,14 @@ func TestSurfaceAppRoutes(t *testing.T) {
 		for _, target := range []string{"/surface/", "/surface/nope", "/surface/app.jsx", "/surface/index.html", "/surface/data.js",
 			"/surface/spike-data.json", "/surface/app/app.js", "/surface/static/x", "/surface/fonts.css",
 			"/surface/passkey/", "/surface/passkey/index.html", "/surface/passkey/passkey.js", "/surface/passkey.html",
-			"/surface/approve", "/surface/approve/", "/surface/requests", "/surface/requests.json/x"} {
+			"/surface/approve", "/surface/approve/", "/surface/requests", "/surface/requests.json/x",
+			"/surface/decide", "/surface/decide/", "/surface/decide.js/x"} {
 			if w := request("GET", target); w.Code != 404 {
 				t.Errorf("%s = %d, want 404", target, w.Code)
 			}
 		}
-		for _, target := range []string{"/surface/requests.json", "/surface/requests/mr-0123456789abcdef", "/surface/approve/mr-0123456789abcdef", "/surface/approve.js"} {
+		for _, target := range []string{"/surface/requests.json", "/surface/requests/mr-0123456789abcdef", "/surface/approve/mr-0123456789abcdef", "/surface/approve.js",
+			"/surface/decide/sha256:" + strings.Repeat("ab", 32), "/surface/decide.js"} {
 			if w := request("GET", target); w.Code != 421 { // registered; example.com is no Host of the surface's
 				t.Errorf("%s = %d, want 421", target, w.Code)
 			}
diff --git a/host/daemon/surface_decide_page.go b/host/daemon/surface_decide_page.go
new file mode 100644
index 0000000..901dff9
--- /dev/null
+++ b/host/daemon/surface_decide_page.go
@@ -0,0 +1,29 @@
+package daemon
+
+// Row 201 A4 sprint 4 (design §2.8): the decision card page and its script. Static embedded bytes
+// under surfaceCSP, behind approvalRoute's Host check like every approval route; the page reads the
+// subject itself from GET /surface/decisions/{ref} and POSTs to POST /surface/decide/{ref}, which
+// runs the law again.
+
+import (
+	"net/http"
+
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// handleSurfaceDecidePage is GET /surface/decide/{ref}: the same bytes for every well-formed ref.
+func (d *Daemon) handleSurfaceDecidePage(w http.ResponseWriter, r *http.Request, _ []string) {
+	if !decisionRefPattern.MatchString(r.PathValue("ref")) {
+		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "a request ref is sha256: and 64 lowercase hex")
+		return
+	}
+	setSurfaceHeaders(w)
+	_, _ = w.Write(surfacepage.DecidePage)
+}
+
+// handleSurfaceDecideScript is the card's script, under the same headers.
+func (d *Daemon) handleSurfaceDecideScript(w http.ResponseWriter, _ *http.Request, _ []string) {
+	setSurfaceHeaders(w)
+	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
+	_, _ = w.Write(surfacepage.DecideScript)
+}
diff --git a/host/daemon/surface_decide_page_test.go b/host/daemon/surface_decide_page_test.go
new file mode 100644
index 0000000..df6248a
--- /dev/null
+++ b/host/daemon/surface_decide_page_test.go
@@ -0,0 +1,140 @@
+package daemon
+
+// Row 201 A4 sprint 4 (design §2.8, §6 PD1, PD2): the decision card page, its script and its field
+// table.
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"encoding/base64"
+	"encoding/hex"
+	"encoding/json"
+	"os"
+	"path/filepath"
+	"reflect"
+	"regexp"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// decisionFields is decide.js's DECISION_FIELDS table.
+func decisionFields(t *testing.T) []struct{ Object, Path, Label, Type string } {
+	t.Helper()
+	_, rest, _ := bytes.Cut(surfacepage.DecideScript, []byte("/* DECISION_FIELDS */"))
+	lit, _, _ := bytes.Cut(rest, []byte("/* /DECISION_FIELDS */"))
+	var fields []struct{ Object, Path, Label, Type string }
+	if err := json.Unmarshal(lit, &fields); err != nil || len(fields) == 0 {
+		t.Fatalf("decide.js's DECISION_FIELDS: %v", err)
+	}
+	return fields
+}
+
+// PD1: DECISION_FIELDS names exactly the keys, and their JSON types, of a stored
+// world/approval-request/v1 and world/decision-packet/v1 as the subject serves them (the broker's
+// and the coordinator's own encoders wrote them): a wire field added without the page fails here.
+func TestDecidePageFieldTable(t *testing.T) {
+	r := newDecideRig(t)
+	q, _ := r.packet("fields", broker.WorkspaceScope, "agent")
+	rec := r.serve("GET", "/surface/decisions/"+q.String(), apLoopback, nil)
+	var s surfacepage.DecisionSubject
+	if err := json.Unmarshal(rec.Body.Bytes(), &s); err != nil || rec.Code != 200 {
+		t.Fatalf("subject = %d %s", rec.Code, rec.Body.String())
+	}
+	wire := map[string]string{}
+	for object, b64 := range map[string]string{"request": s.RequestBytes, "packet": s.PacketBytes} {
+		b, err := base64.RawURLEncoding.DecodeString(b64)
+		if err != nil {
+			t.Fatal(err)
+		}
+		var m map[string]any
+		if err := json.Unmarshal(b, &m); err != nil {
+			t.Fatal(err)
+		}
+		for k, v := range m {
+			kind := "number"
+			if _, ok := v.(string); ok {
+				kind = "string"
+			}
+			wire[object+"."+k] = kind
+		}
+	}
+	table := map[string]string{}
+	for _, f := range decisionFields(t) {
+		kind := map[string]string{"string": "string", "int": "number", "time": "number"}[f.Type]
+		if kind == "" || f.Label == "" {
+			t.Errorf("DECISION_FIELDS %s.%s: type %q, label %q", f.Object, f.Path, f.Type, f.Label)
+		}
+		table[f.Object+"."+f.Path] = kind
+	}
+	if !reflect.DeepEqual(table, wire) {
+		t.Fatalf("DECISION_FIELDS %v, the stored wires %v", table, wire)
+	}
+	if len(wire) != 12 { // control: request 5 + packet 7 (design §2.8)
+		t.Fatalf("the wires have %d fields", len(wire))
+	}
+}
+
+// PD2: GET /surface/decide/{ref} and /surface/decide.js serve the embedded bytes, byte-equal to
+// host/surfacepage/decide/, under the unchanged surfaceCSP and the surface headers; a malformed ref
+// is 400 and a foreign Host 421; the A2 and A3 pages' bytes are unchanged; the page loads exactly
+// its own script.
+func TestDecidePageAssets(t *testing.T) {
+	const csp = "default-src 'none'; style-src 'unsafe-inline'; script-src 'self'; connect-src 'self'; font-src 'self'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'"
+	if surfaceCSP != csp {
+		t.Fatalf("surfaceCSP changed: %q", surfaceCSP)
+	}
+	a := newApprovalRig(t)
+	ref := "sha256:" + strings.Repeat("ab", 32)
+	for _, tc := range []struct{ target, file, contentType string }{
+		{"/surface/decide/" + ref, "index.html", "text/html; charset=utf-8"},
+		{"/surface/decide/sha256:" + strings.Repeat("0", 64), "index.html", "text/html; charset=utf-8"},
+		{"/surface/decide.js", "decide.js", "text/javascript; charset=utf-8"},
+	} {
+		want, err := os.ReadFile(filepath.Join("..", "surfacepage", "decide", tc.file))
+		if err != nil {
+			t.Fatal(err)
+		}
+		rec := a.serve("GET", tc.target, apLoopback, nil)
+		if rec.Code != 200 || len(want) == 0 || !bytes.Equal(rec.Body.Bytes(), want) {
+			t.Errorf("GET %s: %d, or not the bytes of host/surfacepage/decide/%s", tc.target, rec.Code, tc.file)
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
+		requireRefusal(t, "foreign host", a.serve("GET", tc.target, "evil.example", nil), 421, "MisdirectedRequest")
+	}
+	for _, bad := range []string{"sha256:" + strings.Repeat("AB", 32), "sha256:abc", strings.Repeat("ab", 32), "mr-0123456789abcdef"} {
+		requireRefusal(t, bad, a.serve("GET", "/surface/decide/"+bad, apLoopback, nil), 400, "BadRequest")
+	}
+	for file, sum := range map[string]string{"passkey/index.html": "1695ebc2ae20f303dbfa3b617b3bec4c76035447ec290deba91e695fec9cc06d",
+		"passkey/passkey.js": "0064c6c90099a3565b178cf170618b80a67bbf1f3444a1018cc4ee29f51ed53f",
+		"approve/index.html": "e353429ab873633b9aae73e26cd43cd8697f7f2af647ecf7be842cad27f4950e",
+		"approve/approve.js": "9c7fadd4e8a5985f53c2e95981e6746cf422c03f07a2b459d13ef67893ad4b25"} {
+		b, err := os.ReadFile(filepath.Join("..", "surfacepage", file))
+		if h := sha256.Sum256(b); err != nil || hex.EncodeToString(h[:]) != sum {
+			t.Errorf("the A2/A3 page's %s changed (%v)", file, err)
+		}
+	}
+	page := string(surfacepage.DecidePage)
+	if scripts := regexp.MustCompile(`(?i)<script\b[^>]*>`).FindAllString(page, -1); len(scripts) != 1 || !strings.Contains(page, `<script src="/surface/decide.js"></script>`) {
+		t.Errorf("scripts %q", scripts)
+	}
+	if m := regexp.MustCompile(`(?i)\son[a-z]+\s*=|<form|<img|<link|<iframe|https?://`).FindString(page); m != "" {
+		t.Errorf("the page has %q", m)
+	}
+	// Labelled instrument-health controls (not load-bearing; XD1, XD2 and XD4 are).
+	for _, banned := range []string{"innerHTML", "outerHTML", "insertAdjacentHTML", "document.write", "eval(", "new Function", "document.cookie", "XMLHttpRequest", "sendBeacon", "WebSocket"} {
+		if bytes.Contains(surfacepage.DecideScript, []byte(banned)) {
+			t.Errorf("decide.js uses %s", banned)
+		}
+	}
+}
diff --git a/host/daemon/surface_decide_routes_test.go b/host/daemon/surface_decide_routes_test.go
index 5c508c2..7fd2a8d 100644
--- a/host/daemon/surface_decide_routes_test.go
+++ b/host/daemon/surface_decide_routes_test.go
@@ -28,7 +28,7 @@ func TestDecisionRoutesRegistered(t *testing.T) {
 		got = append(got, m[1])
 	}
 	sort.Strings(got)
-	if want := []string{"GET /surface/decisions/{ref}", "POST /surface/decide/{ref}"}; !reflect.DeepEqual(got, want) {
+	if want := []string{"GET /surface/decide.js", "GET /surface/decide/{ref}", "GET /surface/decisions/{ref}", "POST /surface/decide/{ref}"}; !reflect.DeepEqual(got, want) { // sprint 4: the card
 		t.Fatalf("decision routes %v, want %v", got, want)
 	}
 	r := newApprovalRig(t)
@@ -42,7 +42,7 @@ func TestDecisionRoutesRegistered(t *testing.T) {
 		methods []string
 	}{
 		{"/surface/decisions/" + ref, []string{"POST", "PUT", "DELETE", "OPTIONS"}},
-		{"/surface/decide/" + ref, []string{"GET", "PUT", "DELETE", "OPTIONS"}},
+		{"/surface/decide/" + ref, []string{"PUT", "DELETE", "OPTIONS"}}, // sprint 4: GET is the card
 	} {
 		for _, method := range tc.methods {
 			rec := r.serve(method, tc.target, apLoopback, nil, "Origin", "http://localhost:7644", "Access-Control-Request-Method", "POST")
diff --git a/host/surfacepage/decide.go b/host/surfacepage/decide.go
new file mode 100644
index 0000000..2740d45
--- /dev/null
+++ b/host/surfacepage/decide.go
@@ -0,0 +1,20 @@
+package surfacepage
+
+import _ "embed"
+
+// Row 201 A4 sprint 4 (design_docs/planned/w-world-approvals-by-key-a4.md §2.8): the decision card,
+// served at GET /surface/decide/{ref} and GET /surface/decide.js. The page reads the subject from GET
+// /surface/decisions/{ref}, checks both objects' hashes and the packet's requestRef, shows every
+// field of both, builds the 91-byte decision message itself and POSTs a passkey assertion or a pasted
+// SSH signature to POST /surface/decide/{ref}; the daemon runs the law again. The bytes are the same
+// for every ref.
+
+// DecidePage is GET /surface/decide/{ref}.
+//
+//go:embed decide/index.html
+var DecidePage []byte
+
+// DecideScript is GET /surface/decide.js.
+//
+//go:embed decide/decide.js
+var DecideScript []byte
diff --git a/host/surfacepage/decide/decide.js b/host/surfacepage/decide/decide.js
new file mode 100644
index 0000000..3002f87
--- /dev/null
+++ b/host/surfacepage/decide/decide.js
@@ -0,0 +1,300 @@
+'use strict';
+// AILANG World's decision card (row 201 A4, design §2.8). It reads one decision subject from the
+// daemon, checks that the request's and the packet's bytes hash to their refs and that the packet
+// carries this request, shows EVERY field of both objects (and refuses an object with a field it
+// cannot show), builds the 91-byte decision message itself (never a hash from the server), and
+// approves or denies with a passkey or a pasted SSH signature by a POST the daemon judges with the
+// decision law. Text is set with textContent only.
+(() => {
+  const TIMEOUT_MS = 240000; // A2's ceremony bound (A2 §2.14 [r2])
+  const NET_MS = 10000; // A3 §2.7 [r2]: every GET and POST, its body included
+  const FRESH_S = 60; // re-read the subject before signing when less than this remains of its window
+  const PREFIX = 'world-decision/v1', NAMESPACE = 'world-decision@sunholo';
+  const VERBS = { approve: 1, deny: 2 }; // approval.VerbApprove, approval.VerbDeny
+  // Every key of the request (world/approval-request/v1) and of the packet (world/decision-packet/v1),
+  // with its label and JSON type. TestDecidePageFieldTable compares this list with the stored wires.
+  const DECISION_FIELDS = /* DECISION_FIELDS */[
+    {"object": "request", "path": "effect", "label": "Effect", "type": "string"},
+    {"object": "request", "path": "scope", "label": "Scope", "type": "string"},
+    {"object": "request", "path": "cost", "label": "Cost", "type": "int"},
+    {"object": "request", "path": "requester", "label": "Requested by", "type": "string"},
+    {"object": "request", "path": "now", "label": "Requested at", "type": "time"},
+    {"object": "packet", "path": "packetHash", "label": "Packet hash", "type": "string"},
+    {"object": "packet", "path": "proposalHash", "label": "Proposal", "type": "string"},
+    {"object": "packet", "path": "requestRef", "label": "Request", "type": "string"},
+    {"object": "packet", "path": "createdAt", "label": "Created at", "type": "time"},
+    {"object": "packet", "path": "deadlineAt", "label": "Deadline", "type": "time"},
+    {"object": "packet", "path": "escalationsRemaining", "label": "Escalations remaining", "type": "int"},
+    {"object": "packet", "path": "policy", "label": "Policy", "type": "string"}
+  ]/* /DECISION_FIELDS */;
+  const GROUPS = [['request', 'The request'], ['packet', 'The decision packet']];
+
+  const $ = (id) => document.getElementById(id);
+  const setState = (s) => { document.body.dataset.state = s; };
+  const say = (text, error) => {
+    $('status').textContent = text;
+    $('status').toggleAttribute('data-error', !!error);
+  };
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
+  const sha256 = async (bytes) => new Uint8Array(await crypto.subtle.digest('SHA-256', bytes));
+  const Verb = (verb) => (verb === 'deny' ? 'Deny' : 'Approve');
+  // A2 §2.7's assertion file, built here rather than by toJSON() (as approve.js does).
+  const assertionJSON = (c) => ({ id: c.id, rawId: b64u(c.rawId), type: c.type,
+    response: { clientDataJSON: b64u(c.response.clientDataJSON), authenticatorData: b64u(c.response.authenticatorData),
+      signature: b64u(c.response.signature), userHandle: c.response.userHandle ? b64u(c.response.userHandle) : null } });
+
+  const outstanding = new Set();
+  addEventListener('pagehide', () => { for (const c of outstanding) c.abort(); });
+  async function net(url, init) {
+    const c = new AbortController();
+    outstanding.add(c);
+    const timer = setTimeout(() => c.abort(), NET_MS);
+    try {
+      const r = await fetch(url, Object.assign({ cache: 'no-store', credentials: 'omit', signal: c.signal }, init));
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
+  // unshown is the first key of object name the card cannot show (unlisted, of another type, or
+  // missing), or ''. A value that is not a JSON object is itself unshowable.
+  function unshown(name, obj) {
+    if (obj === null || typeof obj !== 'object' || Array.isArray(obj)) return name;
+    const fields = DECISION_FIELDS.filter((f) => f.object === name);
+    for (const k of Object.keys(obj)) {
+      const f = fields.find((x) => x.path === k);
+      if (!f || typeof obj[k] !== (f.type === 'string' ? 'string' : 'number') || (f.type !== 'string' && !Number.isSafeInteger(obj[k]))) return name + '.' + k;
+    }
+    for (const f of fields) if (!Object.prototype.hasOwnProperty.call(obj, f.path)) return name + '.' + f.path;
+    return '';
+  }
+  const shown = (f, v) => (f.type === 'time' ? `${new Date(v * 1000).toISOString()} (${v})` : String(v));
+  const parse = (bytes) => JSON.parse(new TextDecoder('utf-8', { fatal: true }).decode(bytes));
+
+  // message is the 91-byte decision message (design §2.3): PREFIX ‖ 0x00 ‖ verb ‖ P ‖ Q ‖ u64be(expiry).
+  function message(verb, P, Q, expiry) {
+    const m = new Uint8Array(91);
+    m.set(new TextEncoder().encode(PREFIX), 0);
+    m[18] = VERBS[verb];
+    m.set(P, 19);
+    m.set(Q, 51);
+    new DataView(m.buffer).setBigUint64(83, BigInt(expiry));
+    return m;
+  }
+
+  function render(objects, entry) {
+    const card = $('card');
+    card.textContent = '';
+    for (const [name, title] of GROUPS) {
+      const section = document.createElement('section');
+      section.dataset.group = name;
+      const h = document.createElement('h2');
+      h.textContent = title;
+      const dl = document.createElement('dl');
+      for (const f of DECISION_FIELDS.filter((x) => x.object === name)) {
+        const dt = document.createElement('dt');
+        dt.textContent = f.label;
+        const dd = document.createElement('dd');
+        dd.id = 'f-' + name + '.' + f.path;
+        dd.textContent = shown(f, objects[name][f.path]);
+        dl.append(dt, dd);
+      }
+      if (name === 'packet') {
+        const dt = document.createElement('dt');
+        dt.textContent = 'Committed at entry';
+        const dd = document.createElement('dd');
+        dd.id = 'f-entry';
+        const a = document.createElement('a');
+        a.id = 'proposal';
+        a.href = '/workbench?object=' + encodeURIComponent(objects.packet.proposalHash);
+        a.textContent = `${entry} · the proposal`;
+        dd.append(a);
+        dl.append(dt, dd);
+      }
+      section.append(h, dl);
+      card.append(section);
+    }
+    card.hidden = false;
+  }
+
+  const ref = (location.pathname.match(/^\/surface\/decide\/(sha256:[0-9a-f]{64})$/) || [])[1];
+  const ref16 = ref ? ref.slice(7, 23) : '';
+  let subject = null, P = null, Q = null, expiry = 0, messages = {}, group = null, renderedAt = 0, again = null;
+  const buttons = (on) => { for (const b of ['go', 'deny', 'post-sig', 'post-sig-deny']) $(b).disabled = !on; };
+  const offer = (label, action) => { again = action; $('again').textContent = label; $('again').hidden = false; $('again').disabled = false; };
+  const refuse = (text) => { setState('refused'); $('what').textContent = text; };
+
+  // load reads the subject; keep leaves the status line as it is (after a decision).
+  async function load(keep) {
+    $('again').hidden = true;
+    for (const id of ['go', 'deny', 'ssh', 'save', 'save-deny', 'origin-note']) $(id).hidden = true;
+    buttons(false);
+    setState('loading');
+    let res;
+    try {
+      res = await net('/surface/decisions/' + ref, { headers: { Accept: 'application/json' } });
+    } catch (e) {
+      setState('failed');
+      say(`The request could not be read (${why(e)}). Nothing to decide yet.`, true);
+      return offer('Read it again', load);
+    }
+    if (res.status !== 200 || !res.body) {
+      setState('failed');
+      say(`The daemon did not give the request (HTTP ${res.status}${res.body && res.body.reason ? ': ' + res.body.reason : ''}).`, true);
+      return offer('Read it again', load);
+    }
+    const s = res.body;
+    let rb, pb, objects;
+    try {
+      rb = unb64u(s.requestBytes);
+      pb = unb64u(s.packetBytes);
+    } catch (e) {
+      return refuse('The daemon\'s objects are not base64url; nothing to decide.');
+    }
+    const q = await sha256(rb), p = await sha256(pb);
+    if (s.request !== ref || 'sha256:' + hex(q) !== ref || 'sha256:' + hex(p) !== s.packet) {
+      return refuse('The daemon\'s objects do not match their hashes; nothing to decide.');
+    }
+    try {
+      objects = { request: parse(rb), packet: parse(pb) };
+    } catch (e) {
+      return refuse(`The daemon's objects are not UTF-8 JSON (${e.message}); nothing to decide.`);
+    }
+    const bad = unshown('request', objects.request) || unshown('packet', objects.packet);
+    if (bad) {
+      return refuse(`This request has a field the page cannot show (${bad}); nothing to decide; decide it at a terminal: ailang-worldd approve ${ref}.`);
+    }
+    if (objects.packet.requestRef !== ref) {
+      return refuse('The packet does not carry this request; nothing to decide.');
+    }
+    if (!Number.isSafeInteger(s.now) || !s.window || !Number.isSafeInteger(s.window.max) || s.window.max < 0) {
+      return refuse('The daemon\'s decision window is malformed; nothing to decide.');
+    }
+    subject = s; P = p; Q = q; expiry = s.window.max; renderedAt = performance.now();
+    messages = { approve: message('approve', P, Q, expiry), deny: message('deny', P, Q, expiry) };
+    const h = { approve: await sha256(messages.approve), deny: await sha256(messages.deny) };
+    render(objects, s.entry);
+    const r = objects.request;
+    $('what').textContent = `Decide request ${ref}: ${r.effect} for ${r.scope}, requested by ${r.requester}; a decision made here expires ${new Date(expiry * 1000).toLocaleString()}; message hash approve ${hex(h.approve.slice(0, 8))}, deny ${hex(h.deny.slice(0, 8))}`;
+    if (s.decision) {
+      setState('decided');
+      if (!keep) say(`Already decided: ${s.decision} by ${s.decidedBy}.`, false);
+      return;
+    }
+    if (s.reserved) {
+      setState('reserved');
+      say(`A ${s.reserved} approval: use ${s.reserved === 'publish' ? 'world-publish approve' : 'ailang-worldd episode finish'} at a terminal. ${s.refusal}`, true);
+      return;
+    }
+    if (s.refusal) {
+      setState('refused');
+      say(`This request cannot be decided here: ${s.refusal}`, true);
+      return;
+    }
+    group = (s.passkeys || []).find((g) => g.origin === location.origin) || null;
+    if (group) {
+      $('go').hidden = false;
+      $('deny').hidden = false;
+    } else {
+      const others = (s.passkeys || []).map((g) => g.origin);
+      $('origin-note').textContent = others.length ? `No passkey is enrolled for ${location.origin}: open this card at ${others.join(' or ')}.` : 'No passkey is enrolled; decide with an SSH key below.';
+      $('origin-note').hidden = false;
+    }
+    for (const [id, verb] of [['save', 'approve'], ['save-deny', 'deny']]) {
+      $(id).href = URL.createObjectURL(new Blob([messages[verb]], { type: 'application/octet-stream' }));
+      $(id).download = `world-decision-${verb}-${ref16}.msg`;
+      $(id).hidden = false;
+    }
+    $('sshline').textContent = ['approve', 'deny'].map((v) => `ssh-keygen -Y sign -n ${NAMESPACE} -f <key> world-decision-${v}-${ref16}.msg`).join('\n');
+    $('ssh').hidden = false;
+    buttons(true);
+    setState('ready');
+  }
+
+  // fresh re-reads the subject once when less than FRESH_S seconds remain of its window.
+  async function fresh() {
+    if (subject.window.max - subject.now - (performance.now() - renderedAt) / 1000 >= FRESH_S) return true;
+    await load(false);
+    return document.body.dataset.state === 'ready';
+  }
+
+  async function post(verb, format) {
+    buttons(false);
+    setState('posting');
+    say(`${Verb(verb)} request ${ref}: sending the decision…`, false);
+    let res;
+    try {
+      res = await net('/surface/decide/' + ref, { method: 'POST', headers: { 'Content-Type': 'application/json' },
+        body: JSON.stringify(Object.assign({ v: 1, verb, packet: hex(P), request: hex(Q), expiry }, format)) });
+    } catch (e) {
+      setState('unknown');
+      say(`Decision outcome unknown: ${why(e)}. It may or may not have been decided; read the request again before deciding again.`, true);
+      return offer('Read it again', load);
+    }
+    const b = res.body || {};
+    if (res.status === 200) {
+      say(`${b.decided === 'deny' ? 'Denied' : 'Approved'} by ${b.decidedBy}: decision ${b.decision}.`, false);
+      return load(true);
+    }
+    setState('rejected');
+    say(`Not decided (HTTP ${res.status} ${b.error || ''}): ${b.reason || 'no reason given'}`, true);
+    buttons(true);
+  }
+
+  async function decideWithPasskey(verb) {
+    buttons(false);
+    if (!(await fresh())) return;
+    setState('pending');
+    say(`${Verb(verb)} request ${ref}: waiting for your passkey…`, false);
+    const challenge = await sha256(messages[verb]);
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
+      cred = await navigator.credentials.get({ publicKey: { challenge, rpId: group.rp,
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
+    await post(verb, { webauthn: assertionJSON(cred) });
+  }
+
+  $('go').addEventListener('click', () => { decideWithPasskey('approve'); });
+  $('deny').addEventListener('click', () => { decideWithPasskey('deny'); });
+  $('post-sig').addEventListener('click', () => { post('approve', { sshsig: $('sig').value }); });
+  $('post-sig-deny').addEventListener('click', () => { post('deny', { sshsig: $('sig').value }); });
+  $('again').addEventListener('click', () => { if (again) again(); });
+  load(false); // the daemon serves this page only at /surface/decide/sha256:<64 hex>
+})();
diff --git a/host/surfacepage/decide/index.html b/host/surfacepage/decide/index.html
new file mode 100644
index 0000000..51bf14a
--- /dev/null
+++ b/host/surfacepage/decide/index.html
@@ -0,0 +1,39 @@
+<!doctype html>
+<html lang="en">
+<head>
+<meta charset="utf-8">
+<meta name="viewport" content="width=device-width, initial-scale=1">
+<title>AILANG World: decide an approval request</title>
+<style>
+body { font: 16px/1.5 system-ui, sans-serif; max-width: 46rem; margin: 2rem auto; padding: 0 1rem; color: #1d1d1f; }
+h1 { font-size: 1.4rem; }
+h2 { font-size: 1.05rem; margin: 1.2rem 0 0.3rem; }
+#what, #status { white-space: pre-wrap; }
+#status[data-error] { color: #b00020; font-weight: 600; }
+dl { display: grid; grid-template-columns: max-content 1fr; gap: 0.15rem 1rem; margin: 0; }
+dt { color: #555; }
+dd { margin: 0; font: 14px ui-monospace, monospace; overflow-wrap: anywhere; }
+#sshline { font: 13px ui-monospace, monospace; overflow-wrap: anywhere; white-space: pre-wrap; }
+#sig { font: 12px ui-monospace, monospace; width: 100%; height: 7rem; }
+button, #save, #save-deny { font-size: 1rem; padding: 0.4rem 1rem; }
+[hidden] { display: none !important; }
+</style>
+</head>
+<body data-state="loading">
+<h1>Decide an approval request</h1>
+<p id="what">Reading the request…</p>
+<div id="card" hidden></div>
+<p id="origin-note" hidden></p>
+<p><button id="go" type="button" hidden disabled>Approve with passkey</button> <button id="deny" type="button" hidden disabled>Deny with passkey</button></p>
+<p id="status"></p>
+<p><button id="again" type="button" hidden>Read it again</button></p>
+<div id="ssh" hidden>
+<h2>Or decide with an SSH key</h2>
+<p><a id="save" hidden>Save the approve message</a> <a id="save-deny" hidden>Save the deny message</a></p>
+<p id="sshline"></p>
+<p><label for="sig">Paste the signature</label><br><textarea id="sig"></textarea></p>
+<p><button id="post-sig" type="button" disabled>Approve with this signature</button> <button id="post-sig-deny" type="button" disabled>Deny with this signature</button></p>
+</div>
+<script src="/surface/decide.js"></script>
+</body>
+</html>
````

### Patch M7g — the card in live Chrome: approve and deny with a passkey (XD1)

````diff M7g
diff --git a/cmd/ailang-worldd/decide_page_chrome_test.go b/cmd/ailang-worldd/decide_page_chrome_test.go
new file mode 100644
index 0000000..5d24c0b
--- /dev/null
+++ b/cmd/ailang-worldd/decide_page_chrome_test.go
@@ -0,0 +1,86 @@
+package main
+
+// Row 201 A4 sprint 4 (design §2.8, §6 XD1): the decision card decides with a passkey in a live Chrome, gated as
+// XB1 (WORLD_CHROME, WORLD_REQUIRE_PASSKEY_BROWSER=1; skipped and counted without them).
+
+import (
+	"encoding/base64"
+	"fmt"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// cardPasskeys is the subject's passkeys for approvaltest's ES256 seed at origin (rpID localhost).
+func cardPasskeys(origin string, seed byte) []surfacepage.CardPasskeys {
+	id := approvaltest.NewPasskey(approval.AlgES256, seed).CredentialID()
+	return []surfacepage.CardPasskeys{{Origin: origin, RP: "localhost", Allow: []string{base64.RawURLEncoding.EncodeToString(id)}}}
+}
+
+// XD1: mark-mac's passkey (Chrome's virtual authenticator, approvaltest's key, enrolled in the
+// credentials store) decides two packets from their cards: Approve on the first is one decision
+// approve by webauthn:mark-mac, Deny on the second one deny. Each get's challenge is the decision
+// message's hash for its own verb and the window's max, with the subject's rpID, its credential and
+// user verification; each POST carries the verb, P, Q, the expiry and the assertion, with A3a's
+// browser headers and no cookie; the page says the verb it did and then reads the request again,
+// decided, with no decide control. With less than a minute of the window left, the page reads the
+// subject again before it signs. No CSP violation, no page error.
+func TestDecidePageApproveWithPasskey(t *testing.T) {
+	b := startPasskeyBrowser(t)
+	b.do(t, nil, "tab", map[string]any{"authenticator": "uv", "inject": "getRecorder"})
+	r := newDecideCardRig(t)
+	virtualPasskey(t, b, r.creds, r.origin, "mark-mac", 91, 1000)
+	for i, tc := range []struct {
+		verb, button, said string
+		v                  byte
+		now, deadline      int64
+		gets               int
+	}{
+		{"approve", "go", "Approved by webauthn:mark-mac: decision sha256:", approval.VerbApprove, 1710, 2400, 1},
+		{"deny", "deny", "Denied by webauthn:mark-mac: decision sha256:", approval.VerbDeny, 1710, 2400, 1},
+		{"approve", "go", "Approved by webauthn:mark-mac: decision sha256:", approval.VerbApprove, 2350, 2400, 2}, // 50 s left: read again
+	} {
+		req, pkt := cardObjects(broker.WorkspaceScope, fmt.Sprintf("agent-%d", i), tc.deadline)
+		s := cardSubject(req, pkt, tc.now, tc.deadline)
+		s.Passkeys = cardPasskeys(r.origin, 91)
+		url := r.add(s, broker.WorkspaceScope)
+		var view cardView
+		b.do(t, &view, "openCard", map[string]any{"url": url})
+		ref16 := strings.TrimPrefix(s.Request, "sha256:")[:16]
+		if view.State != "ready" || view.ButtonHidden || view.SaveHidden || !strings.HasPrefix(view.What, "Decide request "+s.Request+": Human.Approve for "+broker.WorkspaceScope+", requested by agent-") ||
+			view.Fields["request.scope"] != broker.WorkspaceScope || view.Fields["packet.requestRef"] != s.Request || view.Fields["packet.policy"] != "Cancel" ||
+			view.Fields["packet.deadlineAt"] != fmt.Sprintf("1970-01-01T00:%02d:%02d.000Z (%d)", tc.deadline/60, tc.deadline%60, tc.deadline) || len(view.Fields) != 13 ||
+			view.SSHLine != "ssh-keygen -Y sign -n world-decision@sunholo -f <key> world-decision-approve-"+ref16+".msg\nssh-keygen -Y sign -n world-decision@sunholo -f <key> world-decision-deny-"+ref16+".msg" {
+			t.Fatalf("%d: the card shows %+v", i, view)
+		}
+		var label, href string
+		b.do(t, &label, "waitFor", map[string]any{"expression": "document.getElementById('" + tc.button + "').textContent"})
+		b.do(t, &href, "waitFor", map[string]any{"expression": "document.getElementById('proposal').getAttribute('href')"})
+		if label != map[string]string{"approve": "Approve with passkey", "deny": "Deny with passkey"}[tc.verb] || href != "/workbench?object=sha256%3A"+hashref.SumSHA256([]byte("plan")).Digest() {
+			t.Fatalf("%d: the button reads %q, the proposal link %q", i, label, href)
+		}
+		b.do(t, &view, "press", map[string]any{"id": tc.button, "states": []string{"decided", "rejected", "failed", "unknown", "refused"}})
+		if view.State != "decided" || !strings.HasPrefix(view.Status, tc.said) || !view.ButtonHidden || !view.SaveHidden || len(view.Violations) != 0 || len(view.Errors) != 0 {
+			t.Fatalf("%d: after %s the page shows %+v", i, tc.verb, view)
+		}
+		g := view.Gets[len(view.Gets)-1]
+		if len(view.Gets) != 1 || g.Challenge != decisionHash64(t, tc.v, s, s.Window.Max) || g.RPID != "localhost" || g.UserVerification != "required" ||
+			len(g.Allow) != 1 || g.Allow[0] != s.Passkeys[0].Allow[0] || g.Timeout != 240000 || !g.Signal {
+			t.Fatalf("%d: get %+v (of %d), want the %s hash", i, g, len(view.Gets), tc.verb)
+		}
+		p := r.posts[len(r.posts)-1]
+		if len(r.posts) != i+1 || p.status != 200 || p.body.Verb != tc.verb || p.body.Packet != strings.TrimPrefix(s.Packet, "sha256:") ||
+			p.body.Request != strings.TrimPrefix(s.Request, "sha256:") || p.body.Expiry != s.Window.Max || p.body.WebAuthn == nil || p.body.SSHSIG != nil || p.body.Typed ||
+			s.Decision != tc.verb || s.DecidedBy != "webauthn:mark-mac" {
+			t.Fatalf("%d: posts %d, last %+v; subject %+v", i, len(r.posts), p, s)
+		}
+		if gets := r.gets[s.Request]; gets != tc.gets+1 { // the reads before signing, and one after the decision
+			t.Fatalf("%d: the subject was read %d times, want %d", i, gets, tc.gets+1)
+		}
+	}
+}
diff --git a/cmd/ailang-worldd/decide_page_rig_test.go b/cmd/ailang-worldd/decide_page_rig_test.go
new file mode 100644
index 0000000..e63025e
--- /dev/null
+++ b/cmd/ailang-worldd/decide_page_rig_test.go
@@ -0,0 +1,199 @@
+package main
+
+// Row 201 A4 sprint 4 (design §2.8, §6 XD1, XD2, XD4): the decision card's browser rig, on the A2/A3
+// driver (testdata/passkey_cdp.mjs, unchanged): Chrome headless, a temporary profile and a CDP
+// virtual authenticator holding a test key (D-WORLD-109). The page and its script are served by the
+// production daemon handler on the loopback; GET /surface/decisions/{ref} and POST
+// /surface/decide/{ref} are this rig's, over the same temporary credentials store: the subject is
+// served as given (tests tamper with it), and a POST is held to A3a's browser guards, parsed by
+// the POST's own grammar (surfacepage.ParseDecideBody), put through the decision law
+// (broker.DecisionRefusal) and judged by the daemon's own judges (authority.JudgeDecisionAssertion,
+// JudgeDecisionSSH) at the subject's now. The daemon's POST arms over a world store are
+// TestDecidePOST*'s (host/daemon); this rig proves what the page sends them.
+
+import (
+	"context"
+	"encoding/base64"
+	"encoding/json"
+	"fmt"
+	"io"
+	"mime"
+	"net/http"
+	"path/filepath"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/coordinator"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// cardPost is one POST /surface/decide/{ref} the rig received, and its answer.
+type cardPost struct {
+	ref    string
+	body   surfacepage.DecideBody
+	status int
+	header http.Header
+}
+
+// decideCardRig serves the decision card for the subjects it holds.
+type decideCardRig struct {
+	t        *testing.T
+	creds    string
+	origin   string
+	subjects map[string]*surfacepage.DecisionSubject
+	scopes   map[string]string
+	gets     map[string]int
+	posts    []cardPost
+}
+
+// cardSubject is a subject for request bytes req and packet bytes pkt at now, its window's max
+// min(now + 600, deadline) as the daemon's.
+func cardSubject(req, pkt []byte, now, deadline int64) *surfacepage.DecisionSubject {
+	return &surfacepage.DecisionSubject{V: 1, Request: hashref.SumSHA256(req).String(), RequestBytes: base64.RawURLEncoding.EncodeToString(req),
+		Packet: hashref.SumSHA256(pkt).String(), PacketBytes: base64.RawURLEncoding.EncodeToString(pkt), Entry: 7,
+		Proposal: hashref.SumSHA256([]byte("plan")).String(), Now: now, Window: surfacepage.DecisionWindow{Min: now + 1, Max: broker.DecisionWindowMax(now, deadline)},
+		Passkeys: []surfacepage.CardPasskeys{}, SSH: surfacepage.DecisionSSH{Namespace: approval.DecisionNamespace}}
+}
+
+// cardObjects are a Human.Approve request of scope by requester and a packet carrying it with deadline.
+func cardObjects(scope, requester string, deadline int64) (req, pkt []byte) {
+	req = []byte(fmt.Sprintf(`{"effect":"Human.Approve","scope":%q,"cost":1,"requester":%q,"now":1700}`, scope, requester))
+	pkt = coordinator.EncodeDecisionPacket(coordinator.DecisionPacket{ProposalHash: hashref.SumSHA256([]byte("plan")), RequestRef: hashref.SumSHA256(req),
+		CreatedAt: 1700, DeadlineAt: deadline, Policy: "Cancel"})
+	return req, pkt
+}
+
+func newDecideCardRig(t *testing.T) *decideCardRig {
+	t.Helper()
+	r := &decideCardRig{t: t, creds: filepath.Join(t.TempDir(), "credentials.db"), subjects: map[string]*surfacepage.DecisionSubject{},
+		scopes: map[string]string{}, gets: map[string]int{}}
+	c, err := store.OpenCredentials(r.creds) // make the store before the daemon opens it
+	if err != nil {
+		t.Fatal(err)
+	}
+	_ = c.Close()
+	r.origin = cardOrigin(t, r.creds, "", func(h http.Handler) http.Handler {
+		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
+			switch {
+			case req.Method == "GET" && strings.HasPrefix(req.URL.Path, "/surface/decisions/"):
+				r.serveSubject(w, strings.TrimPrefix(req.URL.Path, "/surface/decisions/"))
+			case req.Method == "POST" && strings.HasPrefix(req.URL.Path, "/surface/decide/"):
+				r.decide(w, req, strings.TrimPrefix(req.URL.Path, "/surface/decide/"))
+			default:
+				h.ServeHTTP(w, req)
+			}
+		})
+	})
+	return r
+}
+
+// add holds subject s (for request scope) under its request ref and returns the card's URL.
+func (r *decideCardRig) add(s *surfacepage.DecisionSubject, scope string) string {
+	r.subjects[s.Request], r.scopes[s.Request] = s, scope
+	return r.origin + "/surface/decide/" + s.Request
+}
+
+func (r *decideCardRig) serveSubject(w http.ResponseWriter, ref string) {
+	r.gets[ref]++
+	s, ok := r.subjects[ref]
+	w.Header().Set("Content-Type", "application/json")
+	if !ok {
+		w.WriteHeader(404)
+		_, _ = w.Write([]byte(`{"v":1,"error":"NotFound","reason":"no such approval request"}`))
+		return
+	}
+	_ = json.NewEncoder(w).Encode(s)
+}
+
+func (r *decideCardRig) refuse(w http.ResponseWriter, p cardPost, status int, class, reason string) {
+	p.status = status
+	r.posts = append(r.posts, p)
+	w.Header().Set("Content-Type", "application/json")
+	w.WriteHeader(status)
+	_ = json.NewEncoder(w).Encode(surfacepage.ApproveRefusal{V: 1, Error: class, Reason: reason})
+}
+
+func (r *decideCardRig) decide(w http.ResponseWriter, req *http.Request, ref string) {
+	p := cardPost{ref: ref, header: req.Header.Clone()}
+	mt, _, _ := mime.ParseMediaType(req.Header.Get("Content-Type"))
+	if req.Header.Get("Origin") != r.origin || req.Header.Get("Sec-Fetch-Site") != "same-origin" || mt != "application/json" ||
+		req.Header.Get("Cookie") != "" || req.Header.Get("Authorization") != "" {
+		r.refuse(w, p, 403, "CrossOrigin", "the rig's guards: Origin, Sec-Fetch-Site, Content-Type, no cookie")
+		return
+	}
+	raw, _ := io.ReadAll(io.LimitReader(req.Body, surfacepage.MaxApproveBodyBytes))
+	body, packet, request, err := surfacepage.ParseDecideBody(raw)
+	p.body = body
+	if err != nil {
+		r.refuse(w, p, 400, "BadRequest", err.Error())
+		return
+	}
+	s, ok := r.subjects[ref]
+	if !ok {
+		r.refuse(w, p, 404, "NotFound", "no such approval request "+ref)
+		return
+	}
+	pkt, err := coordinator.DecodeDecisionPacket(mustB64(s.PacketBytes))
+	if err != nil {
+		r.refuse(w, p, 409, "Conflict", err.Error())
+		return
+	}
+	f := broker.DecisionFacts{Ref: ref, RefOK: true, RequestOK: true, OnChain: true, Discovery: broker.Discovery{Complete: true, Carrying: 1},
+		Scope: r.scopes[ref], Decided: s.Decision != "", DecidedVerb: s.Decision, DecidedBy: s.DecidedBy, Now: s.Now, Deadline: pkt.DeadlineAt,
+		Packet: refDigest(s.Packet), Request: refDigest(ref), PostedPacket: packet, PostedRequest: request, Verb: body.Verb, Expiry: body.Expiry}
+	if law := broker.DecisionRefusal(f); law.Step != "" {
+		r.refuse(w, p, law.Status, "Conflict", law.Reason)
+		return
+	}
+	verb, _ := approval.DecisionVerb(body.Verb) // D9 held
+	ro, err := store.OpenCredentialsReadOnly(r.creds)
+	if err != nil {
+		r.t.Error(err)
+		return
+	}
+	defer func() { _ = ro.Close() }()
+	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
+	defer cancel()
+	in := authority.DecisionJudgeInput{Verb: verb, Packet: packet, Request: request, Expiry: body.Expiry, Now: s.Now}
+	var signed *approval.SignedDecision
+	if body.WebAuthn != nil {
+		signed, err = authority.JudgeDecisionAssertion(ctx, ro, in, *body.WebAuthn)
+	} else if body.SSHSIG != nil {
+		signed, err = authority.JudgeDecisionSSH(ctx, ro, in, []byte(*body.SSHSIG))
+	} else {
+		r.refuse(w, p, 403, "Forbidden", "the card never posts typed")
+		return
+	}
+	if err != nil {
+		r.refuse(w, p, 422, "NotApproved", "the decision does not decide request "+ref+": "+err.Error())
+		return
+	}
+	s.Decision, s.DecidedBy = body.Verb, signed.DecidedBy()
+	s.Refusal = broker.SubjectRefusal(broker.DecisionFacts{Ref: ref, RefOK: true, RequestOK: true, OnChain: true, Discovery: broker.Discovery{Complete: true, Carrying: 1},
+		Decided: true, DecidedVerb: s.Decision, DecidedBy: s.DecidedBy, Now: s.Now, Deadline: pkt.DeadlineAt}).Reason
+	p.status = 200
+	r.posts = append(r.posts, p)
+	w.Header().Set("Content-Type", "application/json")
+	_ = json.NewEncoder(w).Encode(surfacepage.DecisionMade{V: 1, Decided: body.Verb, Request: ref, Decision: "sha256:" + strings.Repeat("d", 64), DecidedBy: s.DecidedBy})
+}
+
+func mustB64(s string) []byte {
+	b, _ := base64.RawURLEncoding.DecodeString(s)
+	return b
+}
+
+// decisionHash64 is the base64url SHA-256 of the decision message for (verb, s, expiry), as a get's challenge.
+func decisionHash64(t *testing.T, verb byte, s *surfacepage.DecisionSubject, expiry int64) string {
+	t.Helper()
+	h, err := approval.DecisionMessageHash(verb, refDigest(s.Packet), refDigest(s.Request), expiry)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return base64.RawURLEncoding.EncodeToString(h[:])
+}
````

### Patch M7h — the card in live Chrome: fails closed, SSH paste (XD2, XD4)

````diff M7h
diff --git a/cmd/ailang-worldd/decide_page_closed_test.go b/cmd/ailang-worldd/decide_page_closed_test.go
new file mode 100644
index 0000000..71c3d7a
--- /dev/null
+++ b/cmd/ailang-worldd/decide_page_closed_test.go
@@ -0,0 +1,137 @@
+package main
+
+// Row 201 A4 sprint 4 (design §2.8, §6 XD2): the decision card fails closed, in a live Chrome, gated
+// as XB1 (WORLD_CHROME, WORLD_REQUIRE_PASSKEY_BROWSER=1; skipped and counted without them).
+
+import (
+	"encoding/base64"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/coordinator"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// XD2: the card fails closed. A subject whose bytes do not hash to the refs, a packet carrying
+// another request, a key the table does not list in either object, a listed key missing, a value of
+// another type, bytes that are not base64url and a malformed window: refused with the reason, no
+// fields, no decide control, no get. A reserved, a refused and an already decided subject render
+// their fields but offer no decide control. Control: the untampered subject is ready.
+func TestDecidePageFailsClosed(t *testing.T) {
+	b := startPasskeyBrowser(t)
+	b.do(t, nil, "tab", map[string]any{"authenticator": "uv", "inject": "getRecorder"})
+	r := newDecideCardRig(t)
+	other, _ := cardObjects("worktree:other", "agent", 2400)
+	packetFor := func(req []byte, extra string) []byte {
+		p := coordinator.EncodeDecisionPacket(coordinator.DecisionPacket{ProposalHash: hashref.SumSHA256([]byte("plan")), RequestRef: hashref.SumSHA256(req),
+			CreatedAt: 1700, DeadlineAt: 2400, Policy: "Cancel"})
+		if extra != "" {
+			p = append(p[:len(p)-1], []byte(","+extra+"}")...)
+		}
+		return p
+	}
+	// subject is the card URL of request bytes req and the subject served for it, after edit.
+	subject := func(req []byte, edit func(s *surfacepage.DecisionSubject)) cardCase {
+		s := cardSubject(req, packetFor(req, ""), 1710, 2400)
+		s.Passkeys = cardPasskeys(r.origin, 92)
+		ref := s.Request
+		if edit != nil {
+			edit(s)
+		}
+		r.subjects[ref] = s
+		return cardCase{r.origin + "/surface/decide/" + ref, s}
+	}
+	const unshown = "This request has a field the page cannot show ("
+	for _, tc := range []struct {
+		name, state, what string
+		c                 cardCase
+	}{
+		{"control", "ready", "Decide request ", subject([]byte(`{"effect":"Human.Approve","scope":"worktree:1","cost":1,"requester":"a","now":1700}`), nil)},
+		{"request bytes swapped", "refused", "The daemon's objects do not match their hashes; nothing to decide.",
+			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:2","cost":1,"requester":"a","now":1700}`), func(s *surfacepage.DecisionSubject) {
+				s.RequestBytes = base64.RawURLEncoding.EncodeToString(other)
+			})},
+		{"packet hash", "refused", "The daemon's objects do not match their hashes; nothing to decide.",
+			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:3","cost":1,"requester":"a","now":1700}`), func(s *surfacepage.DecisionSubject) {
+				s.Packet = hashref.SumSHA256([]byte("another packet")).String()
+			})},
+		{"another request's subject", "refused", "The daemon's objects do not match their hashes; nothing to decide.",
+			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:4","cost":1,"requester":"a","now":1700}`), func(s *surfacepage.DecisionSubject) {
+				s.RequestBytes = base64.RawURLEncoding.EncodeToString(other) // its bytes and ref agree, but it is not this card's
+				s.Request = hashref.SumSHA256(other).String()
+			})},
+		{"packet for another request", "refused", "The packet does not carry this request; nothing to decide.",
+			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:5","cost":1,"requester":"a","now":1700}`), func(s *surfacepage.DecisionSubject) {
+				p := packetFor(other, "")
+				s.Packet, s.PacketBytes = hashref.SumSHA256(p).String(), base64.RawURLEncoding.EncodeToString(p)
+			})},
+		{"unlisted request key", "refused", unshown + "request.extra)",
+			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:6","cost":1,"requester":"a","now":1700,"extra":"x"}`), nil)},
+		{"unlisted packet key", "refused", unshown + "packet.extra)",
+			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:7","cost":1,"requester":"a","now":1700}`), func(s *surfacepage.DecisionSubject) {
+				p := packetFor(mustB64(s.RequestBytes), `"extra":1`)
+				s.Packet, s.PacketBytes = hashref.SumSHA256(p).String(), base64.RawURLEncoding.EncodeToString(p)
+			})},
+		{"missing key", "refused", unshown + "request.cost)",
+			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:8","requester":"a","now":1700}`), nil)},
+		{"string for int", "refused", unshown + "request.cost)",
+			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:9","cost":"1","requester":"a","now":1700}`), nil)},
+		{"fraction", "refused", unshown + "request.now)",
+			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:10","cost":1,"requester":"a","now":1700.5}`), nil)},
+		{"int for string", "refused", unshown + "request.scope)",
+			subject([]byte(`{"effect":"Human.Approve","scope":11,"cost":1,"requester":"a","now":1700}`), nil)},
+		{"not an object", "refused", unshown + "request)", subject([]byte(`["Human.Approve"]`), nil)},
+		{"not base64url", "refused", "The daemon's objects are not base64url; nothing to decide.",
+			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:12","cost":1,"requester":"a","now":1700}`), func(s *surfacepage.DecisionSubject) { s.PacketBytes = "*" + s.PacketBytes })},
+		{"window", "refused", "The daemon's decision window is malformed; nothing to decide.",
+			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:13","cost":1,"requester":"a","now":1700}`), func(s *surfacepage.DecisionSubject) { s.Window.Max = -1 })},
+		{"reserved", "reserved", "Decide request ",
+			subject([]byte(`{"effect":"Human.Approve","scope":"registry:x","cost":1,"requester":"a","now":1700}`), func(s *surfacepage.DecisionSubject) {
+				s.Reserved, s.Refusal = "publish", "publish approvals are minted by `world-publish approve` at a terminal (R4); this request cannot be decided here"
+			})},
+		{"refused", "refused", "Decide request ",
+			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:14","cost":1,"requester":"a","now":1700}`), func(s *surfacepage.DecisionSubject) {
+				s.Refusal = "the packet's deadline passed at 1970-01-01T00:40:00Z; its policy is Cancel"
+			})},
+		{"decided", "decided", "Decide request ",
+			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:15","cost":1,"requester":"a","now":1700}`), func(s *surfacepage.DecisionSubject) {
+				s.Decision, s.DecidedBy, s.Refusal = "deny", "ssh:mark", "request was already decided deny by ssh:mark"
+			})},
+	} {
+		var view cardView
+		b.do(t, &view, "openCard", map[string]any{"url": tc.c.url})
+		var controls string
+		b.do(t, &controls, "waitFor", map[string]any{"expression": "['go','deny','ssh','save','save-deny'].map((id) => document.getElementById(id).hidden ? '-' : id).join(' ') || 'x'"})
+		if view.State != tc.state || !strings.HasPrefix(view.What, tc.what) || len(view.Gets) != 0 || len(view.Violations) != 0 {
+			t.Fatalf("%s: state %q, what %q, %d gets; want %q, %q", tc.name, view.State, view.What, len(view.Gets), tc.state, tc.what)
+		}
+		switch tc.state {
+		case "ready":
+			if controls != "go deny ssh save save-deny" || len(view.Fields) != 13 {
+				t.Fatalf("%s: controls %q, %d fields", tc.name, controls, len(view.Fields))
+			}
+		case "refused", "reserved", "decided":
+			if controls != "- - - - -" {
+				t.Fatalf("%s: a %s card offers %q", tc.name, tc.state, controls)
+			}
+			wantFields := 0
+			if strings.HasPrefix(tc.what, "Decide request ") {
+				wantFields = 13
+			}
+			if len(view.Fields) != wantFields || (wantFields != 0 && !strings.Contains(view.Status, map[string]string{"reserved": "A publish approval: use world-publish approve at a terminal.",
+				"refused": "This request cannot be decided here: the packet's deadline passed", "decided": "Already decided: deny by ssh:mark."}[tc.state])) {
+				t.Fatalf("%s: %d fields, status %q", tc.name, len(view.Fields), view.Status)
+			}
+		}
+	}
+	if len(r.posts) != 0 {
+		t.Fatalf("%d POSTs from cards that offer nothing", len(r.posts))
+	}
+}
+
+// cardCase is one card's URL and the subject the rig serves for it.
+type cardCase struct {
+	url string
+	s   *surfacepage.DecisionSubject
+}
diff --git a/cmd/ailang-worldd/decide_page_sshsig_test.go b/cmd/ailang-worldd/decide_page_sshsig_test.go
new file mode 100644
index 0000000..517b1f2
--- /dev/null
+++ b/cmd/ailang-worldd/decide_page_sshsig_test.go
@@ -0,0 +1,106 @@
+package main
+
+// Row 201 A4 sprint 4 (design §2.8, §6 XD4): the decision card's SSH path, in a live Chrome, gated as
+// XB1 (WORLD_CHROME, WORLD_REQUIRE_PASSKEY_BROWSER=1; skipped and counted without them); its
+// real-tool sibling needs ssh-keygen (WORLD_REQUIRE_SSH_KEYGEN=1 makes a missing one a failure).
+
+import (
+	"bytes"
+	"os"
+	"os/exec"
+	"path/filepath"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+// XD4: the card's SSH path. Save the approve message and Save the deny message download exactly Go's
+// approval.DecisionMessage (91 bytes) for their verbs and the window's max, under the names the
+// signing lines give; approvaltest's SSHSIG over the approve message, pasted, decides approve as
+// ssh:mark-mac (200); a signature over the approve message posted as a deny is refused by the law
+// (422) and decides nothing. Real-tool sibling (directive retro 11): ssh-keygen -Y sign over the
+// saved deny file with a key made in a temp dir decides deny as ssh:mark-key.
+func TestDecidePageSSHSIGPaste(t *testing.T) {
+	b := startPasskeyBrowser(t)
+	b.do(t, nil, "tab", map[string]any{"authenticator": "uv"})
+	r := newDecideCardRig(t)
+	seedApprover(t, r.creds, "mark-mac", sigSigner.PublicKey(), 1000)
+	for i, scope := range []string{"worktree:a", "worktree:b"} {
+		req, pkt := cardObjects(scope, "agent", 2400)
+		s := cardSubject(req, pkt, 1710, 2400)
+		note := "No passkey is enrolled; decide with an SSH key below."
+		if i == 0 { // a passkey enrolled only at another origin: the card names it and offers no passkey
+			s.Passkeys = cardPasskeys("https://studio.tail0.ts.net", 93)
+			note = "No passkey is enrolled for " + r.origin + ": open this card at https://studio.tail0.ts.net."
+		}
+		url := r.add(s, scope)
+		var view cardView
+		b.do(t, &view, "openCard", map[string]any{"url": url})
+		ref16 := strings.TrimPrefix(s.Request, "sha256:")[:16]
+		if view.State != "ready" || !view.ButtonHidden || view.Note != note || view.SaveHidden {
+			t.Fatalf("%d: the card shows %+v", i, view)
+		}
+		dl := t.TempDir()
+		msgs := map[string][]byte{}
+		for id, verb := range map[string]byte{"save": approval.VerbApprove, "save-deny": approval.VerbDeny} {
+			var downloads []struct{ Event, State, Name, URL string }
+			b.do(t, &downloads, "save", map[string]any{"id": id, "dir": dl})
+			name := map[byte]string{approval.VerbApprove: "approve", approval.VerbDeny: "deny"}[verb]
+			got, err := os.ReadFile(filepath.Join(dl, "world-decision-"+name+"-"+ref16+".msg"))
+			want, _ := approval.DecisionMessage(verb, refDigest(s.Packet), refDigest(s.Request), s.Window.Max)
+			if err != nil || len(want) != 91 || !bytes.Equal(got, want) || !strings.HasPrefix(downloads[0].URL, "blob:"+r.origin+"/") {
+				t.Fatalf("%d: %s saved %x (%v), want %x; %+v", i, id, got, err, want, downloads)
+			}
+			msgs[name] = got
+		}
+		if i == 1 {
+			ssh, err := exec.LookPath("ssh-keygen")
+			if err != nil {
+				if os.Getenv("WORLD_REQUIRE_SSH_KEYGEN") == "1" {
+					t.Fatal("WORLD_REQUIRE_SSH_KEYGEN=1 but ssh-keygen was not found")
+				}
+				t.Log("SKIP COUNTED: ssh-keygen not found; the card's real-tool sibling did NOT run")
+				return
+			}
+			_ = ssh
+			t.Setenv("HOME", t.TempDir())
+			t.Setenv("SSH_AUTH_SOCK", "")
+			key := filepath.Join(t.TempDir(), "k")
+			sshKeygen(t, "-q", "-t", "ed25519", "-N", "", "-C", "xd4", "-f", key)
+			pub, err := os.ReadFile(key + ".pub")
+			if err != nil {
+				t.Fatal(err)
+			}
+			k, err := approval.ParseAuthorizedKey(pub)
+			if err != nil {
+				t.Fatal(err)
+			}
+			seedApprover(t, r.creds, "mark-key", k, 1000)
+			file := filepath.Join(dl, "world-decision-deny-"+ref16+".msg")
+			sshKeygen(t, "-Y", "sign", "-n", approval.DecisionNamespace, "-f", key, file)
+			sig, err := os.ReadFile(file + ".sig")
+			if err != nil {
+				t.Fatal(err)
+			}
+			b.do(t, nil, "fill", map[string]any{"id": "sig", "value": string(sig)})
+			b.do(t, &view, "press", map[string]any{"id": "post-sig-deny", "states": []string{"decided", "rejected", "failed", "unknown"}})
+			if view.State != "decided" || !strings.HasPrefix(view.Status, "Denied by ssh:mark-key: decision sha256:") || s.Decision != "deny" || s.DecidedBy != "ssh:mark-key" {
+				t.Fatalf("ssh-keygen: the page shows %+v; subject %+v", view, s)
+			}
+			continue
+		}
+		sig := sigSigner.Sign(approval.DecisionNamespace, msgs["approve"])
+		b.do(t, nil, "fill", map[string]any{"id": "sig", "value": string(sig)})
+		b.do(t, &view, "press", map[string]any{"id": "post-sig-deny", "states": []string{"decided", "rejected", "failed", "unknown"}})
+		if view.State != "rejected" || !strings.Contains(view.Status, "HTTP 422 NotApproved") || s.Decision != "" {
+			t.Fatalf("the approve signature posted as a deny: %+v", view)
+		}
+		b.do(t, &view, "press", map[string]any{"id": "post-sig", "states": []string{"decided", "rejected", "failed", "unknown"}})
+		p := r.posts[len(r.posts)-1]
+		if view.State != "decided" || !strings.HasPrefix(view.Status, "Approved by ssh:mark-mac: decision sha256:") || s.Decision != "approve" ||
+			p.body.SSHSIG == nil || *p.body.SSHSIG != string(sig) || p.body.Expiry != s.Window.Max || p.body.WebAuthn != nil {
+			t.Fatalf("the approve signature: %+v; post %+v", view, p)
+		}
+	}
+}
````

### Patch M7i — QUICKSTART §18's card paragraph

````diff M7i
diff --git a/cmd/ailang-worldd/quickstart_decide_card_test.go b/cmd/ailang-worldd/quickstart_decide_card_test.go
new file mode 100644
index 0000000..f806054
--- /dev/null
+++ b/cmd/ailang-worldd/quickstart_decide_card_test.go
@@ -0,0 +1,61 @@
+package main
+
+// Row 201 A4 sprint 4 (design §2.8, §10 step 4; S7): QUICKSTART §18's card paragraph names the
+// card's own controls, the card's route and an approve --list line the CLI documents.
+
+import (
+	"os"
+	"regexp"
+	"strings"
+	"testing"
+)
+
+func TestQuickstartSection18CardMatchesThePage(t *testing.T) {
+	guide, err := os.ReadFile("../../docs/QUICKSTART.md")
+	if err != nil {
+		t.Fatal(err)
+	}
+	page, err := os.ReadFile("../../host/surfacepage/decide/index.html")
+	if err != nil {
+		t.Fatal(err)
+	}
+	routes, err := os.ReadFile("../../host/daemon/daemon.go")
+	if err != nil {
+		t.Fatal(err)
+	}
+	_, section, ok := strings.Cut(string(guide), "## 18. Decide an agent's approval request")
+	if !ok {
+		t.Fatal("QUICKSTART has no §18")
+	}
+	_, card, ok := strings.Cut(section, "**From the surface.**")
+	if !ok {
+		t.Fatal("§18 has no card paragraph")
+	}
+	labels := regexp.MustCompile(`\*\*([^*]+)\*\*`).FindAllStringSubmatch(card, -1)
+	if len(labels) < 6 {
+		t.Fatalf("the card paragraph names %d controls, want the card's six and Needs you", len(labels))
+	}
+	for _, l := range labels {
+		if l[1] != "Needs you" && !strings.Contains(string(page), ">"+l[1]+"<") {
+			t.Errorf("§18 names **%s**, which the card page does not show", l[1])
+		}
+	}
+	if !strings.Contains(card, "/surface/decide/<request-ref>") || !strings.Contains(string(routes), `"GET /surface/decide/{ref}"`) {
+		t.Error("§18 does not name the card's route, or the daemon does not serve it")
+	}
+	lists := 0
+	for _, line := range strings.Split(card, "\n") {
+		if !strings.HasPrefix(line, "ailang-worldd approve --list") {
+			continue
+		}
+		lists++
+		for _, flag := range regexp.MustCompile(`\s(--[a-z-]+)`).FindAllStringSubmatch(line, -1) {
+			if !strings.Contains(approveHelp, flag[1]+" ") && !strings.Contains(approveHelp, flag[1]+"]") {
+				t.Errorf("§18 uses %s, which approve's help does not document: %s", flag[1], line)
+			}
+		}
+	}
+	if lists != 1 {
+		t.Errorf("the card paragraph holds %d approve --list lines, want 1", lists)
+	}
+}
diff --git a/docs/QUICKSTART.md b/docs/QUICKSTART.md
index 99b5418..25d6923 100644
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -1886,3 +1886,11 @@ ailang-worldd approve sha256:<ref> --webauthn-assertion ~/Downloads/world-assert
 ```
 
 The daemon records `approve by webauthn:demo-mac` after checking the assertion against the enrolled passkey, the decision message's hash, the origin and a signature counter that must move forward. Posting the same file again is refused as already decided. If the world write fails after the daemon checked a key decision, the same file (or signature) posted again lands that one decision without moving the counter twice; a different assertion carrying the same counter is refused. The packet page then reads `approvals chain: approve by …`.
+
+**From the surface.** The surface app's **Needs you** lists each pending request under the §17 mint requests as *Decide · <scope> · requested by <requester> · deadline <hh:mm>*, and `approve --list` prints the same requests from the store, offline:
+
+```bash
+ailang-worldd approve --list --db /tmp/se-world/world.db
+```
+
+The line opens the decision card at `/surface/decide/<request-ref>`. The card reads the request and its packet from the daemon, checks that both hash to their refs and that the packet carries the request, shows every field of both, and refuses a field it cannot show. **Approve with passkey** or **Deny with passkey** signs that verb's decision message with a §16 passkey, and the card says which verb it is signing. **Save the approve message** (or the deny message) downloads the 91 bytes for the `ssh-keygen -Y sign` line the card prints; paste the signature and press **Approve with this signature** (or **Deny with this signature**). The daemon checks every decision as it checks the CLI's. A publish- or finish-shaped request is listed with its typed owner and has no link. A request the card cannot decide shows the daemon's reason and no button.
````
