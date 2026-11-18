# Sprint plan — w-world-address-bar-verbs-v3, sprint S6 "the surface end to end" (iteration 405, row 243)

> **THE SPRINT — row 243 phase V3, sprint S6, the LAST V3 sprint (design §0, §1, §2.2, §2.5 (the app.js *Needs you*
> paragraph), §2.6, §2.7, §2.9, §3, §4, §5, §6 "Sprint 6" (and Sprint 2's XC1 for the live-Chrome pattern), §7, §8, §11,
> §12, §14, §16), D-WORLD-118 = A:**
> **M15** — app.js *Needs you* lists each `requests.json` `confirmations` entry after the decisions: "Confirm · episode
> finish <ep> · <repository> <head>@<sha12> · draft|ready · expires hh:mm", a link to `/surface/confirm/<id>`; a refused row
> (or one whose id is not `vr-` + 16 lower-case hex) shows its reason and **no link**; `confirmationsUnavailable` is said,
> never hidden; every value through `esc()`. **CP7** `TestNeedsYouListsConfirmations` (node) renders **the real handler's
> answers** — a golden written by `TestRequestsJSONConfirmationsGolden` (host/daemon, `WORLD_UPDATE_CONFIRMATIONS_GOLDEN=1`).
> **M16** — **XC2** `TestConfirmPageConfirmsWithAPasskeyInChrome` (the virtual authenticator confirms a request stored by the
> CLI's own put on the card a real `daemon.New` serves; the world decision is `webauthn:mark-key`'s and the row is consumed
> naming it; the CLI's own re-check and read-back accept both) and **XC3** `TestFinishConfirmOnSurfaceInChrome` (the real
> `episode finish … --confirm-on-surface` on S5's rig waits; Chrome confirms; the CLI prints the fake gh's PR URL).
> **M17** — the carried items: judge-404's four survivors killed by named tests (`E-LK-VR-LEN`, `E-JS-INCLUDES`/`E-JS-NODASH`,
> `E-RC-INPUT-NOW`, `E-AWAIT-NOW-STALE`), S5 §1.4 F8 `TestCLICredentialWriterSites`, judge-403 notes 2 and 3 (note 4 decided,
> S6-D11). **No schema change, no daemon route, no `.ail`, `surfaceCSP` unchanged, no CLI behaviour change: the only
> production byte S6 changes is `host/surfacepage/app/app.js`.**
> - **Measured: 425 changed lines** — `git diff -M --numstat 539cecd0 HEAD -- cmd host docs app` (insertions + deletions,
>   renames on; no rename occurs; **the same 425 with no path filter**; `app/` does not exist at the repo root — app.js is
>   `host/surfacepage/app/app.js`, under `host`), after a commit of each of the three scratch stages: **17 production**
>   (`app.js` +14 −3), **1 golden** (generated), **407 test**. **Cap 1,400: margin 975.**
> - **No cut (S6-D1).** Design ≈ 355 raw → **425 compiled (1.20×)** with every carried item (1 → 3 → 2) taken; nothing parked.
> - **Three sub-milestones, each green at its commit** (§1.2 C14): **M15** *Needs you* + CP7 + its golden (157) → **M16**
>   XC2 + XC3 (290) → **M17** the carried items (425).
> - **Every line ships as a compiled patch** (Appendix P), extracted from this file and round-tripped onto a fresh
>   `git archive 539cecd0` export, byte-identical to the stage tree (§1.2 C16). The executor applies, never retypes.
> - **Pre-existing test files changed: 4; deleted: 0** (§5, §9 criterion 12): `host/verbs/lookup_test.go` (E-LK-VR-LEN, design §7),
>   `host/surfacepage/app_verb_test.go` (E-JS-*, as S5 F6), `host/surfacepage/confirm_script_test.go` (judge-403 note 2) and
>   `host/daemon/surface_confirm_test.go` (judge-403 note 3) — the last two outside design §7 (§1.4 F7). **No
>   `episode_finish*_test.go`.**
> - **V3 closes with this sprint** (Appendix C states the V3 exit).

---

## 0. Binding constraints (each is restated where it governs)

- **D-WORLD-118 = A** (controller, iteration 405: `mission_decisions.sh --open` → no OPEN row; `--check` valid, 105 rows; A
  since iteration 400, attended). S6 is routable; no B-only branch exists or is built.
- **No daemon route, no schema change, no `.ail`, `surfaceCSP` unchanged, no CLI behaviour change**: every production file
  but `host/surfacepage/app/app.js` is untouched — the gates' `UNTOUCHED_VIOLATIONS` lists every changed file that is not a
  `_test.go`, not under a `testdata/` dir and not `app.js` (plus anything under `host/store`, `world`, `packages`, `scripts`,
  `tools`, any `.ail`), and must say `none`; `DAEMON_ROUTES` base = final, `SURFACE_POSTS: 3`;
  `TestDaemonCredentialWriterSites`, `TestFrozenV1RouteTableMatchesMux`, `TestConfirmRoutesRegistered`, `TestDecidePageAssets`
  (the CSP literal) green unedited.
- **D-WORLD-109**: every store under `t.TempDir()`; keys are `approvaltest`'s (`sigSigner` seed 11 as `mark-mac` in S5's
  fixtures; the passkey seeds 93 (XC2) and 94 (XC3) enrolled as `mark-key` through V2b's `virtualPasskey`, which hands the
  same private key to Chrome's CDP virtual authenticator); Chrome is headless on the driver's **temporary profile**
  (`testdata/passkey_cdp.mjs`, unchanged); the daemon is `daemon.New` over test stores on an `httptest` loopback listener;
  `gh` is the finish rig's `fakeFinishDispatch`; the terminal is S5's `countingTerminal`. No real key, `~/.ssh`,
  `~/.ailang/world`, `/dev/tty`, browser profile, mint or publish.
- Coding standards: the page's *Needs you* is pure (`needsYouHTML`, `confirmLineHTML`: data in, HTML out, no DOM, no fetch);
  no Go production change; "why not a package" is design §3 S3.
- **Not touched**: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/`, `cmd/world-publish/*`, every
  production `.go` file, `confirm.js`, `approve.js`, `decide.js`, `surface_page.go`, `daemon.go`, the `/v1` table,
  `testdata/passkey_cdp.mjs`, every pre-existing `episode_finish*_test.go`.

---

## 1. Pre-flight

### 1.1 Baseline at `539cecd0` (measured by the planner, 2026-11-18)

`tools/gates.sh <base-tree> base-gates base 539cecd0` (the §6 script in `base` mode, on a clone of the repo at `539cecd0`;
`GATES_SKIP_VERIFY_GO=1`), run beside the planner's compile and mutant work (load ≈ 9–35):

| leg | rc | time | note |
|---|---|---|---|
| `TOUCHED_PACKAGES` | — | — | *(empty at the base)* |
| `status` | 0 | 1 s |  |
| `verify_ail` | 0 | 43 s | `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` — the final must match |
| `gofmt`, `build`, `vet` | 0 | 0 / 7 / 4 s |  |
| `fence` | 0 | 24 s | `go test -run '^$' ./...` is the test-file compile fence |
| `pkgs` | 0 | 1,015 s | `host/verbs` 0.1 s, `host/surfacepage` 1.4 s, `cmd/ailang-worldd` 352 s, `host/daemon` 659 s (load ≈ 15–35) |
| `race` | 0 | 161 s | `host/verbs`, `host/surfacepage` whole; the `cmd` subset 80 s; the `host/daemon` subset 41 s |
| `node` | 0 | 6 s | the twelve named harnesses `--- PASS`, `WORLD_EXEC_NODE` set |
| `inv` | 0 | 26 s | **48 named `--- PASS`**, 0 SKIP, 0 FAIL |
| `DAEMON_ROUTES` / `SURFACE_POSTS` | — | — | **47** (base 47) / **3** |
| `UNTOUCHED_VIOLATIONS` | — | — | none |
| `chrome` | 0 | 80 s | **17 `--- PASS`** (V2b's fourteen incl. the needs-you pair, XC1, the two `host/daemon` drills) |
| **`GATES_DONE 1368s FAILS: none`** | | | the final tree is compared leg by leg with these (C17) |

### 1.2 Codebase claims checked (at `539cecd0` unless "work" — the planner's stage tree, tag `M17`)

| # | claim | command | observed |
|---|---|---|---|
| C1 | **Input 1 — the cites, re-measured at `539cecd0`: none moved since S5's plan** | `grep -n` each | `app.js:429` `function needsYouHTML(r, now, permission)`, `:496` the three-prefix navigate allowlist (S5's `/surface/confirm/vr-` third); `surface_approve.go:76` `approvalRoute`, `:92` `handleSurfaceRequests`, `:113-121` the `readConfirmations` call, `isOlderStore` → "the credentials store predates verb requests: restart \`ailang-worldd serve\` with this build once", other error → "the verb requests could not be read", nil → `[]`; `surface_confirm_read.go:88` `readConfirmations` (`PendingVerbRequests(ctx, now, store.MaxPendingVerbRequests)`), `:106` `pendingConfirmation` (digest mismatch / not a valid verb request / scope not a finish approval → `Refused`, else `Episode, Repository, Head, SHA, Draft`); `approve_model.go:75-76` `Confirmations`/`ConfirmationsUnavailable`, `:81` `PendingConfirmation{ID, Episode, Repository, Head, SHA, Draft *bool, Expires, Refused}` (omitempty but `id`/`expires`), `:149` `CardPasskey`; `lookup.go:129` `ApproveTarget`, `:136` the `vr-` arm `len(args) == 1 && verbRequestID.MatchString(args[0])`; `episode_finish_surface_rig_test.go:61` `newSurfaceRig`; `passkey_page_driver_test.go:63` `startPasskeyBrowser`; `confirm_page_chrome_test.go:32` XC1; `approve_page_driver_test.go:56` `cardOrigin`, `:77` `virtualPasskey`; `app_decide_test.go:18` `TestNeedsYouListsDecisions`; `confirm.go:15,20` and `approve.go:13,18` the embeds; `surface_page.go:19` `surfaceCSP` (last changed by `17c21dae`, row 201 A2). All as S5's Appendix C describes them |
| C2 | the store's id grammar the page mirrors | `grep -n 'VerbRequestIDPattern *=' host/store/verb_requests.go` | `:58` `regexp.MustCompile(\`^vr-[0-9a-f]{16}$\`)` — M15's `/^vr-[0-9a-f]{16}$/` is the same expression |
| C3 | **Input 2 — the closest landed analogues, and what was mirrored** | `sed -n` of each | **Mirrored**: `TestNeedsYouListsDecisions` (the node harness: `AppScript` + a harness file in `t.TempDir()`, `WORLD_EXEC_NODE` non-skippable, a 30 s context, a PASS manifest) → CP7; `TestSurfaceVerbOutcomes` (`host/daemon` writes a golden under `host/surfacepage/testdata`, regenerated by an env var, else byte-compared) → CP7's fixture; XC1 (`startPasskeyBrowser`, `openCard raw`, `waitFor` over the card's own ids) and A4's XD1 `TestDecidePageApproveWithPasskey` (`tab {authenticator: uv, inject: getRecorder}`, `virtualPasskey`, the get's challenge = the message hash) → XC2; S5's `newSurfaceRig` + `TestFinishConfirmOnSurfaceOpensThePR` (the finish fixture with a real `daemon.New`, `r.onWait(1)` where the card is played, `r.run`, the fake gh's URL) → XC3; `TestDaemonCredentialWriterSites`' AST walk → `TestCLICredentialWriterSites`. **Differs**: XC3 serves the rig's daemon on an `httptest` loopback listener for Chrome (S5's rig served it only through `Handler().ServeHTTP`), and plays the card in Chrome instead of POSTing an SSHSIG; the needs-you harness reads the daemon's own answers instead of hand-typed objects |
| C4 | nothing named like S6's code exists | `grep -rn 'confirmLineHTML\|requests-confirmations\|TestNeedsYouListsConfirmations\|TestCLICredentialWriterSites\|TestConfirmPageConfirmsWithAPasskeyInChrome\|TestFinishConfirmOnSurfaceInChrome\|WORLD_UPDATE_CONFIRMATIONS_GOLDEN' cmd host`; control `grep -rn 'function decisionLineHTML' host/surfacepage/app/app.js` | *(empty)*; control `:444` |
| C5 | every type, field and function the patches name | `grep -n` at HEAD | **existing**: `esc`, `hm`, `needsYouHTML`, `decisionLineHTML`, `notifyHTML`, `notifications` (app.js); `RequestsJSON.Confirmations/ConfirmationsUnavailable`, `PendingConfirmation` (`approve_model.go`); `AppScript`; `newApprovalRig`, `approvalRig.putVerb/serve`, `setVerbColumn`, `setVersion`, `apT0`, `apLoopback`, `requireRefusal` (host/daemon tests); `startPasskeyBrowser`, `passkeyBrowser.call/do`, `cardGet`, `cardOrigin`, `virtualPasskey`, `surfaceUnit`, `storeOwn`, `verbSig`, `credSQL`, `approveByHand`, `sfT0`, `sigSigner`, `newSurfaceRig`, `finishSurfaceRig.{run,row,object,onWait,creds,d,start,now,tty,f}`, `drawnID`, `boundedTestContext`, `exitOK` (cmd tests); `verbLandGrace`, `awaitVerbConfirmation`, `pollVerbRequest`, `readBackFinishDecision`, `worldReader`, `newClient`, `approverTime`, `sessionEnv{now, after}`, `surfaceFinish{scope, scopeText}`, `verbOwn{id, digest, challenge, expires}` (cmd production, S5); `approval.VerbMessageHash`, `approval.VerbNamespace`, `approval.ErrApproverClockBehind`, `broker.FormatFinishApprovalScope`, `store.OpenCredentialsReadOnly`, `(*Credentials).VerbRequest`, `store.VerbOutcomeApproved`; the driver's `tab`, `openCard {raw}`, `waitFor {expression, ms}`, `tap {id}` (`passkey_cdp.mjs:95,224,277,308`). **NEW**: `confirmLineHTML` (app.js, unexported); tests `TestRequestsJSONConfirmationsGolden`, `confirmGolden`, `TestNeedsYouListsConfirmations`, `needsYouConfirmationsHarness`, `confirmView`, `confirmInChrome`, `TestConfirmPageConfirmsWithAPasskeyInChrome`, `TestFinishConfirmOnSurfaceInChrome`, `TestAwaitVerbConfirmationEndsOnItsClock`, `TestPollVerbRequestRechecksAtItsOwnClock`, `TestCLICredentialWriterSites`; the golden `host/surfacepage/testdata/requests-confirmations.golden.json` |
| C6 | **Input 3 — fixtures from the real tool** | the test files | CP7: `requests.json` bodies written by the real `handleSurfaceRequests` (`a.serve("GET", "/surface/requests.json", …)`) over rows stored by S3's `putVerb` (`approval.CanonicalVerb` of `verbFixture`: `FinishPromptText`, `FormatFinishApprovalScope`), one edited by `setVerbColumn` (refused: digest), then the same read after S3's `DROP TABLE verb_requests` + `PRAGMA user_version = 8`; generating command **`WORLD_UPDATE_CONFIRMATIONS_GOLDEN=1 go test ./host/daemon/ -run '^TestRequestsJSONConfirmationsGolden$'`**; the harness's hostile rows are a real row with one field changed (CF7's form). XC2: the request from the CLI's `finishVerbRequest` stored by its production `putFinishVerbRequest` (`storeOwn`); XC3: the request the CLI under test stores. The passkey: `virtualPasskey` (`EnrolPasskey` of `approvaltest.NewPasskey(ES256, seed)`, the same key to `WebAuthn.addCredential`). The PR URL: `r.f.finishRecords[0].PR.URL`, the fake gh's own answer. No hand-typed bytes or JSON |
| C7 | **Input 4 — tests that enumerate app.js exports, harnesses, Chrome legs, roots, routes, writer sites** | `grep -rln 'module.exports\|Object.keys(app)' host/surfacepage/*_test.go`; `grep -rn 'func Test.*WriterSites\|func TestProductionContextRoots\|func TestFrozenV1RouteTableMatchesMux' cmd host` | **unchanged and green**: `TestSurfaceAppPure` (pins `Object.keys(app)`: `confirmLineHTML` is **not exported**, so the list is unchanged), `TestAddressHintsArePinned`, `TestSurfaceAppMatchesWorkingCopy`, `TestSurfaceAppShellAssembly`, `TestNeedsYouListsDecisions` (its "without confirmations, exactly as before" rows hold: no `confirmations` key → `cs = []`, no unavailable line), `TestProductionContextRoots` (no new store call in production), `TestDaemonCredentialWriterSites` (no daemon change), `TestFrozenV1RouteTableMatchesMux`, `TestConfirmRoutesRegistered`, `TestEveryVerbHasHelp`, `TestProductionGoSurface`. **New inventory**: `TestCLICredentialWriterSites` (S5 F8: none existed — `grep -rn 'func Test.*WriterSites' cmd` → empty at base; control: host/daemon's found). **Changed**: `TestApproveTarget` (+3 rows in one line), `TestSurfaceAppVerbs` (+4 refused hrefs), `TestConfirmPageShowsRefusals` (+1 answer row), `TestConfirmPOSTRefusesBeforeAnyWrite` (the busy reason) |
| C8 | **Input 5 — the instruments, bounded** | the test files; work runs | the CLI's wait in XC3 runs on the rig's fake clock (S5), bounded at `expires_at + 10` polls; every Chrome step is a driver command with its own wall bound (`openCard` 15 s per load, `waitFor` 15 s default / 30 s for the ceremony, `tap` immediate) inside `passkeyBrowserBudget` (120 s, the driver's `exec.CommandContext`); XC2's first step is the **injected-delay control**: `waitFor {expression: false, ms: 300}` must fail, and within 10 s (work: failed after ≈ 0.3 s); `TestAwaitVerbConfirmationEndsOnItsClock` is bounded by its injected clock (parks the wait past its deadline) and by a 20 s wall deadline; **cleanup on failure** C21 |
| C9 | the golden's bytes, and that it is the handler's | work: `WORLD_UPDATE_CONFIRMATIONS_GOLDEN=1 go test ./host/daemon/ -run '^TestRequestsJSONConfirmationsGolden$'`, then the same without the variable | one line: `{"v8":{… "confirmations":[],"confirmationsUnavailable":"the credentials store predates verb requests: restart \`ailang-worldd serve\` with this build once"},"v9":{… "confirmations":[{"id":"vr-00000000000000c6","expires":1800000580,"refused":"the stored request does not match its digest"},{"id":"vr-00000000000000c4","episode":"w1",…,"draft":true,"expires":1800000560},{"id":"vr-00000000000000c2","episode":"w2",…,"draft":false,"expires":1800000550}],"confirmationsUnavailable":""}}`; the second run `ok` (byte-equal) |
| C10 | the passkey path through the CLI's re-check works (a premise S5 tested only with SSHSIG) | work: XC2's last assertion (`pollVerbRequest` over the Chrome-made row) and XC3 | both `--- PASS` (C13): `authority.RecheckVerb` admits the stored WebAuthn assertion after C15 moved the counter; no production change needed |
| C11 | CI never sets `WORLD_CHROME`, so XC2/XC3 skip there, counted | `grep -n 'WORLD_CHROME\|WORLD_REQUIRE_PASSKEY' .github/workflows/*.yml`; control `grep -n 'WORLD_EXEC_NODE' .github/workflows/ci.yml` | *(empty)*; control `:250`, `:289`. The executor's `skipped` step (§3) proves both names `--- SKIP` with `SKIP COUNTED` without a browser and `--- FAIL` with `WORLD_REQUIRE_PASSKEY_BROWSER=1`; the gates' `chrome` leg and M16's commit command run them with Chrome and require `--- PASS` per name |
| C12 | **the four judge-404 survivors reproduced on the base first** | `tools/mut.py <clone at 539cecd0> <id> … --tag=base-*` (ledger lines `base-whole`/`base-named`) | `E-LK-VR-LEN` **SURVIVED** `go test ./host/verbs/ -run '.'`; `E-JS-INCLUDES` and `E-JS-NODASH` **SURVIVED** `go test ./host/surfacepage/ -run '.'`; `E-RC-INPUT-NOW` **SURVIVED** `go test ./cmd/ailang-worldd/ -run '.'` (380 s); `E-AWAIT-NOW-STALE`: **killed by S5's `TestAwaitVerbConfirmationEnds` run alone** (`--- FAIL … the wait did not end: unbounded`, 20 s) — the judge's "timeout only" is the whole package, where the rig tests (`episode_finish_surface_rig_test.go` sorts before `episode_finish_surface_test.go`) wait to expiry on the fake clock with no wall bound and hang first (§1.4 F5) |
| C13 | the new tests, work | `go test -count=1 -v` by name (`accept.sh`; Chrome: `WORLD_CHROME=… WORLD_REQUIRE_PASSKEY_BROWSER=1`) | all `--- PASS`, none SKIP; XC2 3.6 s, XC3 7.3 s (load ≈ 9) |
| C14 | **every sub-milestone is green when committed** | `tools/accept_all.sh` (worktrees of `work` at each stage tag running §3's `accept.sh`) | on the final stage trees (`logs/accept_all.txt`, `logs/accept/*.txt`; worktrees `st-M15`, `st-M16`, `st-M17`/`final` at the stage tags; load ≈ 5–35 beside the gates and the mutant shards): `M15 fast rc=0 ACCEPT_OK` (9–43 s), `M16 fast rc=0` (11–36 s: the `skipped` step saw both names `--- SKIP` with `SKIP COUNTED` and both `--- FAIL` under `WORLD_REQUIRE_PASSKEY_BROWSER=1`), `M17 fast rc=0` (41–65 s), `M17 cmd1 rc=0` (322 s), `M17 cmd2 rc=0` (68 s), `M17 daemon1 rc=0` (158 s), `M17 daemon2 rc=0` (251 s); the controller's `chrome` parts: `M15 chrome rc=0` (11 s), `M16 chrome rc=0` (18 s), `M17 chrome rc=0` (14 s), every named test `--- PASS`. **Measured and changed before shipping**: the whole-package parts first measured `M17 daemon` **649 s** and `M17 cmd` 493 s at load ≈ 11–13 (the base's `pkgs` leg: `host/daemon` 659 s) — over or near one 600 s call — so each is split by test name into two halves that together are the package (`cmd1`/`cmd2`, `daemon1`/`daemon2`); M15's whole-`host/daemon` part was dropped (M15 adds one daemon test, run by name in `fast`; the package runs whole at M17) |
| C15 | the snapshot block (§7) on every sub-milestone | `tools/snap_rehearse.sh` (a clone at `539cecd0`, the patches in order, §7's block verbatim) | `M15 4/0/4`, `M16 5/0/5`, `M17 11/0/11` (`STATUS` lines / `DELETED` lines / MANIFEST lines), each `MANIFEST_OK`, rc 0; dirs `host` for M15 and `cmd host` for M16–M17 (the `find` guard holds both ways) — `logs/snap.txt`, re-run on this final plan file |
| C16 | **the round trip** | `tools/roundtrip.sh <this plan>` | `CMP_EQUAL` for all three (each block extracted from this file with §2's awk equals the generated patch byte for byte); each applied with `/usr/bin/patch -p1 -E` after a `--dry-run` to a fresh `git archive 539cecd0` export, no `.orig`/`.rej`, committed and measured: `NUMSTAT M15 157`, `M16 290`, `M17 425` (each `unfiltered=` the same); then `diff -r -x .git` against the stage tree (`git archive M17`) → **`ROUNDTRIP_IDENTICAL_TO_STAGE_TREE`** (`logs/roundtrip.txt`, re-run on this final plan file); `accept.sh` and `gates.sh` extracted from this file are byte-equal to the scripts that produced C14 and C17 |
| C17 | **the gates on the final tree** | `tools/gates.sh <final-tree> final-gates final <sha>` | `tools/gates.sh <final-tree> final-gates final 7b68b3aa4` (`final-tree` = a worktree of the planner's clone at the stage tag `M17`, the three stage commits on `539cecd0`; `diff -r` equal to the round-trip tree, C16; `GATES_SKIP_VERIFY_GO=1`; load ≈ 5–15): **`TOUCHED_PACKAGES: cmd/ailang-worldd host/daemon host/surfacepage host/verbs`** (= the plan's list); `status` 0; **`verify_ail` rc 0** (`✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` — the base's line, byte for byte); `gofmt`, `build`, `vet`, `fence` rc 0; **`pkgs` rc 0** (766 s: `host/verbs` 0.8 s, `host/surfacepage` 2.7 s, `cmd/ailang-worldd` 342 s, `host/daemon` 416 s); **`race` rc 0** (70 s); **`race3` rc 0** (the two clock tests ×3, 1.7 s); **`node` rc 0** (the thirteen named harnesses incl. CP7); **`inv` rc 0: 52 named `--- PASS`** (the base's 48 + `TestRequestsJSONConfirmationsGolden`, the two clock tests, `TestCLICredentialWriterSites`), 0 SKIP, 0 FAIL, no `NOT PASSED`; **`DAEMON_ROUTES: 47 (base 47)`, `SURFACE_POSTS: 3`, `UNTOUCHED_VIOLATIONS: none`**; **`chrome` rc 0** (47 s: **19 `--- PASS`** = the base's 17 + `TestConfirmPageConfirmsWithAPasskeyInChrome` 1.25 s + `TestFinishConfirmOnSurfaceInChrome` 1.28 s, none SKIP; the needs-you pair green on the new app.js); **`GATES_DONE 988s FAILS: none`**. `verify_go` is the controller's |
| C18 | node, the pinned binaries, Chrome | `/opt/homebrew/bin/node --version`; `$HOME/.pinned-ailang/ailang --version`; the driver's ready line | `v26.13.0`; `AILANG v0.41.0`; Chrome (`/Applications/Google Chrome.app`), headless, the driver's temporary profile |
| C19 | the confirm card's ids vs the driver's `cardState` | `grep -o 'id="[^"]*"' host/surfacepage/confirm/index.html`; `sed -n 214,220p cmd/ailang-worldd/testdata/passkey_cdp.mjs` | the card has `what shown argv prompt status answer controls go save sshline sig post-sig signed`; `cardState` reads `again` and `origin-note` (A3's card) — so `openCard` (not raw) and `press` would throw on this card (§1.4 F3) |
| C20 | the needs-you line's time | `node -e 'console.log(new Date(1800000560*1000).toISOString())'`; `app.js:13` `hm` | `2027-01-15T08:09:20.000Z`; `hm` reads UTC hours and minutes → `08:09` (the harness's literal) |
| C21 | **no Chrome or daemon leak on failure** | work: XC2 with `X-RECHECK-FORMAT` applied (it fails), then `pgrep -fl 'passkey_cdp.mjs'`, `pgrep -fl 'remote-debugging-pipe'` and the driver's profile dirs | before: `passkey_cdp.mjs` drivers 0, headless Chrome processes 2 (not the test's: present before it ran); the run: both tests `--- FAIL` (XC2: "the CLI's re-check of the passkey's confirmation: … approver format must be ssh"; XC3: `rc 1` with the CLI's refusal), 5.5 s; after (+2 s): drivers **0**, headless Chrome **2** (the same), no `TestConfirmPage*`/`TestFinishConfirmOnSurfaceInChrome*` temp dir left — `startPasskeyBrowser`'s cleanup closes the driver's stdin, waits, and SIGKILLs Chrome's pid if the driver did not reap it; the daemons close in `t.Cleanup` (`cardOrigin`, `newSurfaceRig`) and the `httptest` servers in `t.Cleanup(srv.Close)` |

### 1.3 Planner decisions where the design or the directive leaves a choice

- **S6-D1 — no cut; everything carried is taken, by measurement.** M15 157, M16 +133, M17 +135: **425 ≤ 1,400**.
  Carried items in the directive's order 1 → 3 → 2: all of item 1 (four survivors), item 3 (`TestCLICredentialWriterSites`),
  item 2's notes 2 and 3; note 4 is decided (S6-D11). Nothing parked; no patch under `plan-scratch/patches/` is unshipped.
- **S6-D2 — three sub-milestones**, each independently green: **M15** (app.js + the golden's generator + the golden + CP7) —
  the page change ships with its only behavioural proof; **M16** (XC2 + XC3, one new test file, no production byte); **M17**
  (the carried items: two new test files, four one-to-two-line edits of existing tests).
- **S6-D3 — the *Needs you* line, verbatim from design §2.5**: `Confirm · episode finish <ep> · <repository> <head>@<sha12> ·
  draft|ready · expires hh:mm` (sha12 = the first 12 characters of the scope's sha; `hh:mm` = `hm`, UTC, as every Needs-you
  time), linked as `<a id="confirm-<id>" href="/surface/confirm/<id>">`; listed **after the decisions** (design: "under the
  decisions") and **counted** in "Needs you · n"; the order is the daemon's (newest first). `confirmationsUnavailable` is the
  line `Confirm · unavailable: <reason>` after the lines (before the truncation note), in both the empty and the non-empty
  branch. **A confirmation raises no notification** (design §2.5 lists it; it does not notify — the waiting CLI's terminal
  already shows the link); pinned by the harness, carried to V4 (Appendix C).
- **S6-D4 — the id shape check in the page.** The page's navigate allowlist already admits only `/surface/confirm/vr-`
  (S5 M13) and the store's grammar is `^vr-[0-9a-f]{16}$` (C2), so a summarised row is linked only when its id matches that
  expression; otherwise it is listed as `Confirm · <id> · refused: not a confirmation id`, no link. The daemon never emits such
  a row (the store's CHECK and Go grammar, S3), so the branch is defence in depth, and its fixture is the real row with the id
  changed. Every value is still passed through `esc()` (design: "all text through `esc()`"), so the two `esc()` calls on the
  checked id are **named equivalents** (`CP-ESC-ID-A`, `CP-ESC-ID-HREF`: a matching id has no character `esc` changes).
- **S6-D5 — CP7's fixture is the handler's answer**, as a golden: no landed harness consumes daemon-made `requests.json`
  (§1.4 F2); the `verb-route.golden.json` pattern (C3) makes the page's test replay the route's own bytes.
- **S6-D6 — XC2/XC3 live in `cmd/ailang-worldd`** (`confirm_surface_chrome_test.go`): they need `startPasskeyBrowser`,
  `virtualPasskey`, `cardOrigin`, `surfaceUnit`/`storeOwn`, `newSurfaceRig` — all `package main` test helpers; host/daemon's
  Chrome drills (`TestSurfaceLiveChrome`) use another driver and cannot import them. **CI**: no `WORLD_CHROME` (C11) → both
  skip, counted; the gates' `chrome` leg and M16's commit command run them with `WORLD_REQUIRE_PASSKEY_BROWSER=1`.
- **S6-D7 — XC3 plays the card from the CLI's first wait** (`r.onWait(1)`, S5's CF3 hook): the CLI is parked in `env.after`
  while Chrome opens the printed card, clicks *Approve with a passkey* and waits for the card to read `closed`; then the CLI's
  next poll sees the row. The Chrome steps return errors (`b.call`), never `t.Fatal`, from inside the CLI's call stack. The rig
  restarts its fake clock at the wall clock (S5 F3), so the daemon's real clock admits the request.
- **S6-D8 — `E-RC-INPUT-NOW` is killed, not left near-equivalent**, in the direction that matters: a row that *claims* a later
  consume (`consumed_at = +50`) must not admit a key enrolled at +20 to a poll at +19 — the re-check judges at **this
  process's** clock (design §2.7 6(a)); the boundary (+20) passes (P5 is `Now < EnrolledAt`). The mutant (`Now:
  row.ConsumedAt`) admits it.
- **S6-D9 — `E-AWAIT-NOW-STALE` gets a named test on the injected clock**, not a tick count: the wait starts one second before
  its deadline and must end at the deadline after exactly one wait; once the injected clock passes the deadline the test fails
  and the wait is **parked** (a nil channel), never spun, so the kill is named and takes < 1 s. The file sorts before the
  finish rig's, so in a whole-package run the named failure is printed before the rig's (unchanged, design §7) unbounded wait
  meets the package timeout (§4 records both).
- **S6-D10 — `TestCLICredentialWriterSites` (S5 §1.4 F8) taken**: the daemon's walk, over `cmd/ailang-worldd`'s production
  files: exactly `openCredentialsForSession`, `launchEpisode`, `finishRevoke`, `putFinishVerbRequest` open the writer
  (`OpenCredentials`); ≥ 20 read-only opens (23 measured) as the control.
- **S6-D11 — judge-403 note 4 is decided, not built**: a second POST route is a static fact of `daemon.go`'s mux; the
  instruments that see it are static by nature — `TestFrozenV1RouteTableMatchesMux`, `TestConfirmRoutesRegistered` and the
  gates' `DAEMON_ROUTES`/`SURFACE_POSTS: 3` count. A behavioural probe could only POST to paths it already knows, which is
  weaker than the table. Notes 2 and 3 are taken (one row each).
- **S6-D12 — mutants run against the named test AND the WHOLE package of the mutated file** (every test, static included);
  when the whole-package kill names only an inventory/AST/static test, a third leg re-runs the whole package without it and the
  table records both (directive retro 3).
- **S6-D13 — no docs change**: QUICKSTART §20 (S5) already walks the CLI's card link; *Needs you* is a second way to the same
  card and adds no command, flag or output line (S7 binds the sprint that adds a CLI flag).

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

| # | premise | measured | routing |
|---|---|---|---|
| F1 | design §2.5 "a refused row shows its reason and no link" covers every unlinkable row | a summarised row (no `refused`) whose id is not `vr-` + 16 hex is not addressed; the daemon cannot emit one (C2), but the page's allowlist (S5) rejects such an href | S6-D4: such a row is listed refused ("not a confirmation id"), no link |
| F2 | directive: "the existing app.js harness pattern" with fixtures "generated by the real daemon handler" | the existing needs-you harnesses (`app_decide_test.go`, `app_hints_test.go`) use hand-typed objects; none reads a daemon answer | S6-D5: a host/daemon generator writes the golden; the node harness reads it (the verb-route pattern) |
| F3 | S5 Appendix C: XC3 "replacing the SSHSIG POST with the virtual authenticator on `/surface/confirm/<id>`" with the A2/A3 driver | the driver's `cardState` (used by `openCard` non-raw and `press`) reads ids the confirm card does not have (C19) | XC2/XC3 use `openCard {raw}`, `waitFor` and `tap` (XC1's form); the driver is unchanged |
| F4 | S5 Appendix C: "XC2 needs only a stored request (`storeOwn`'s form)" | `surfaceUnit`'s scope expires at `sfT0 + 900` (2030); C10's F2 refuses a scope expiring later than `now + 3600` (`host/broker/finish_confirm.go:28`) | XC2 sets `in.scope.ExpiresAt = now + 900` and re-formats `scopeText` (two lines) before `storeOwn` |
| F5 | judge-404: `E-AWAIT-NOW-STALE` "is killed only by the package timeout" | C12: S5's named test kills it alone in 20 s; the hang is the rig tests', which run first | S6-D9: a named test that kills it in < 1 s and sorts first; the rig file stays unedited |
| F6 | judge-403 note 2: "`confirm.js` checks `confirmed:true` without status 200" | the code checks both (`confirm.js:198`: `res.status === 200 && b.confirmed === true`); what is missing is a test of the status half | one `answers` row: a 500 whose body says `confirmed: true` → "Not confirmed (HTTP 500 Internal): r" (`N2-STATUS` killed) |
| F7 | design §7: the pre-existing tests S6 edits | §7 names `lookup_test.go` (and "the page tests"); judge-403 notes 2–3 land in `confirm_script_test.go` and `surface_confirm_test.go` | both edited, one row/condition each; listed in §5 and §9 criterion 12 |
| F8 | design §6 "XC2 … the world decision exists" via the daemon only | the decision is checked through the CLI's own `readBackFinishDecision` over a real HTTP client to the daemon (`/v1/objects`), plus the row read from the credentials store | stronger than the design's row; no production change |

**Files outside design §7**: `host/daemon/surface_needs_you_golden_test.go` + `host/surfacepage/testdata/requests-confirmations.golden.json`
(S6-D5), `cmd/ailang-worldd/confirm_surface_clock_test.go` and `cli_credential_writers_test.go` (the carried items),
`host/surfacepage/confirm_script_test.go` and `host/daemon/surface_confirm_test.go` (F7). §7's "XC3" lives in
`cmd/ailang-worldd/confirm_surface_chrome_test.go` with XC2 (S6-D6).

### 1.5 The size and the race-mode cost

| stage | production | test / golden | changed lines | cumulative (numstat after commit) |
|---|---|---|---|---|
| M15 | 17 (`app.js` +14 −3) | 140 (`app_confirm_test.go` 91, `surface_needs_you_golden_test.go` 48, golden 1) | 157 | **157** |
| M16 | 0 | 133 (`confirm_surface_chrome_test.go`) | 133 | **290** |
| M17 | 0 | 135 (`confirm_surface_clock_test.go` 67, `cli_credential_writers_test.go` 57, `lookup_test.go` 3, `app_verb_test.go` 3, `confirm_script_test.go` 1, `surface_confirm_test.go` 4) | 135 | **425** |

Design ≈ 355 raw → compiled **425 (1.20×)**: M15 (≈ 85 raw) and M16 (≈ 270 raw) are under their raw estimates; the
carried items (≈ 135) are extra to the design's S6 row.

**Race-mode wall time of the new tests** (directive 8; measured on the work tree at load ≈ 35, `go test -race -count=1 -v`,
`logs/race-new.txt`): `cmd/ailang-worldd` — `TestCLICredentialWriterSites` 0.52 s, `TestAwaitVerbConfirmationEndsOnItsClock`
0.14 s, `TestPollVerbRequestRechecksAtItsOwnClock` 0.55 s, XC2/XC3 **0.00 s (SKIP: CI and the race job set no
`WORLD_CHROME`)** — **≈ 1.2 s added**; `host/daemon` — `TestRequestsJSONConfirmationsGolden` 0.41 s (and the busy-reason
condition 0 s); `host/surfacepage` — `TestNeedsYouListsConfirmations` 0.32 s; `host/verbs` — two table rows. CI's bound is 8
minutes per package; host/daemon's 446–480 s on slow runners gains < 0.1 %. No sleep anywhere in S6's tests.

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of four backticks,
in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; a new file's patch has `--- /dev/null`; **no
sub-milestone deletes a file**. Apply it **only** this way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF2'
set -euo pipefail
M=M15   # <- the sub-milestone, exactly as in the Appendix P heading, in order: M15 M16 M17
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter405-row243-v3s6
PLAN="$W/design_docs/planned/w-world-address-bar-verbs-v3s6-sprint-plan.md"
T="$HOME/.ailang/state/world-iter405/exec-patches"; mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
/usr/bin/patch -p1 -E --dry-run < "$T/$M.patch"
/usr/bin/patch -p1 -E < "$T/$M.patch"
if find cmd host \( -name '*.orig' -o -name '*.rej' \) 2>/dev/null | grep .; then echo "STOP: patch left .orig/.rej"; exit 1; fi
git status --short --untracked-files=all
EOF2
```

- If `--dry-run` fails, **stop and report** (no hand edits, no fuzz options). It can only fail if an earlier
  sub-milestone was skipped, reordered or changed.
- Then run that sub-milestone's **Acceptance** (§3): extract `accept.sh` from this plan (`awk '$0 == "````bash
  accept.sh" {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/accept.sh"`) and run **each part §3 lists for M, one
  per Bash call**: `bash "$T/accept.sh" "$W" <M> <part> 2>&1 | tee
  "$HOME/.ailang/state/world-iter405/exec-logs/accept_<M>_<part>.txt"` (never only through `tail`). **Never `nohup`,
  never background a test leg; pass a 600000 ms timeout to every long Bash call.** The executor never runs the `chrome`
  part (it launches a browser).
- Then snapshot (§7). The controller commits and compares the numstat with §3's expected value.

---

## 3. Milestones

**Acceptance.** Every part runs the same `[base]` checks first (gofmt, build, vet of the four touched packages, the
test-file compile fence `go test -run '^$' ./...`). `v` requires every named test to `--- PASS` and none to SKIP or FAIL;
`p` is a whole-package run; `skipped` proves XC2/XC3 skip **counted** without a browser and **fail** with
`WORLD_REQUIRE_PASSKEY_BROWSER=1` (the CI behaviour and its guard); `vc` runs named live-Chrome tests with `WORLD_CHROME` and
`WORLD_REQUIRE_PASSKEY_BROWSER=1` and requires `--- PASS` (the controller's `chrome` part only). Parts: `fast` (≈ 0.5–1 min),
`cmd1` + `cmd2` (`cmd/ailang-worldd` whole, in two halves by test name: `^Test[A-E]` and `^Test[^A-E]`), `daemon1` +
`daemon2` (`host/daemon` whole, `^Test[A-R]` and `^Test[^A-R]`), `chrome` (controller only). Each part is one 600 s call;
measured in C14 — the halves exist because a whole package measured over 600 s at load (`host/daemon` 649–659 s,
`cmd/ailang-worldd` 493 s).

````bash accept.sh
#!/bin/bash
# accept.sh <tree> <M> [fast|cmd1|cmd2|daemon1|daemon2|chrome]: exactly plan §3's acceptance for sub-milestone M (one bash
# shell; rc 0 = ACCEPT_OK). Each part fits one 600 s call; run every part §3 lists for M, in order. The
# executor never runs `chrome` (it launches a browser); the controller's commit command for M16 does.
T=$1; M=$2; P=${3:-fast}
cd "$T" || exit 2
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang
unset WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER
set -o pipefail
chk() { tee /dev/stderr | { out=$(cat); echo "$out" | grep -qE -- '--- (FAIL|SKIP)|^FAIL' && exit 1; echo "$out" | grep -q -- '--- PASS' || exit 1; }; }
v() { go test -count=1 -v "$@" 2>&1 | chk; }   # -v: every named test PASS, none SKIP/FAIL
p() { go test -count=1 "$@"; }
vc() { WORLD_CHROME='/Applications/Google Chrome.app/Contents/MacOS/Google Chrome' WORLD_REQUIRE_PASSKEY_BROWSER=1 go test -count=1 -v "$@" 2>&1 | chk; }   # live Chrome: PASS required
step() { echo "== $*"; "$@" || { echo "ACCEPT_RED $M $P: $*"; exit 1; }; }
# skipped <names>: without WORLD_CHROME each named Chrome test SKIPs, counted; with WORLD_REQUIRE_PASSKEY_BROWSER=1 each FAILs
skipped() {
  local re="^($(echo "$*" | tr ' ' '|'))\$" out n
  out=$(go test -count=1 -v ./cmd/ailang-worldd/ -run "$re" 2>&1); echo "$out" | grep -E -- '^--- |SKIP COUNTED' | head -20
  for n in "$@"; do echo "$out" | grep -qE -- "^--- SKIP: $n " || { echo "not a counted skip: $n"; return 1; }; done
  echo "$out" | grep -q 'SKIP COUNTED' || return 1
  out=$(WORLD_REQUIRE_PASSKEY_BROWSER=1 go test -count=1 -v ./cmd/ailang-worldd/ -run "$re" 2>&1)
  for n in "$@"; do echo "$out" | grep -qE -- "^--- FAIL: $n " || { echo "WORLD_REQUIRE_PASSKEY_BROWSER=1 did not fail $n"; return 1; }; done
  echo "required: both FAIL without a browser"
}
# [base]
step test -z "$(gofmt -l cmd host)"
step go build ./...
step go vet ./cmd/ailang-worldd/ ./host/daemon/ ./host/surfacepage/ ./host/verbs/
if go test -run '^$' ./... 2>&1 | grep -E '^(FAIL|--- FAIL)'; then echo "ACCEPT_RED $M $P: FENCE_RED"; exit 1; fi
NODE='TestSurfaceAppPure|TestSurfaceScriptPure|TestSurfaceAppVerbs|TestNeedsYouListsDecisions|TestAddressHintsArePinned|TestSurfaceAppMatchesWorkingCopy|TestSurfaceAppShellAssembly|TestConfirmPage.*'
INV='TestEveryVerbHasHelp|TestVerbTableMatchesTheCLISwitch|TestEpisodeFinishHelp|TestQuickstartSection19MatchesTheBar|TestQuickstartSection20MatchesTheFinishFlag'
XC='TestConfirmPageConfirmsWithAPasskeyInChrome TestFinishConfirmOnSurfaceInChrome'
XCHROME='TestConfirmPageShowsTheFinishPromptInChrome|TestNeedsYouListsPendingRequests|TestNeedsYouNotifiesOnlyWhenAskedOnce|TestDecidePageApproveWithPasskey|TestVerbBarEpisodeReportInChrome'
case "$M:$P" in
M15:fast)   step v ./host/surfacepage/ -run "^(TestNeedsYouListsConfirmations|$NODE)\$"
            step p ./host/surfacepage/
            step v ./host/daemon/ -run '^(TestRequestsJSONConfirmationsGolden|TestRequestsJSONConfirmations|TestSurfaceAppRoutes|TestDecidePageAssets)$' ;;
M16:fast)   step skipped $XC
            step v ./cmd/ailang-worldd/ -run "^(TestSurfaceCardLinks|TestPollVerbRequestRechecks|TestFinishConfirmOnSurfaceOpensThePR|$INV)\$" ;;
M17:fast)   step v ./cmd/ailang-worldd/ -run "^(TestAwaitVerbConfirmationEndsOnItsClock|TestPollVerbRequestRechecksAtItsOwnClock|TestCLICredentialWriterSites|TestAwaitVerbConfirmationEnds|TestPollVerbRequestRechecks|$INV)\$"
            step skipped $XC
            step p ./host/verbs/
            step v ./host/verbs/ -run '^(TestLookup|TestApproveTarget)$'
            step v ./host/surfacepage/ -run "^(TestNeedsYouListsConfirmations|$NODE)\$"
            step v ./host/daemon/ -run '^(TestConfirmPOSTRefusesBeforeAnyWrite|TestConfirmPOSTGuards|TestDaemonCredentialWriterSites|TestRequestsJSONConfirmationsGolden)$' ;;
M16:chrome|M17:chrome) step vc ./cmd/ailang-worldd/ -run "^($(echo $XC | tr ' ' '|')|$XCHROME)\$" ;;
M15:chrome) step vc ./cmd/ailang-worldd/ -run "^($XCHROME)\$" ;;
M17:cmd1)    step p -timeout 30m ./cmd/ailang-worldd/ -run '^Test[A-E]' ;;    # the two halves are the whole package
M17:cmd2)    step p -timeout 30m ./cmd/ailang-worldd/ -run '^Test[^A-E]' ;;
M17:daemon1) step p -timeout 30m ./host/daemon/ -run '^Test[A-R]' ;;
M17:daemon2) step p -timeout 30m ./host/daemon/ -run '^Test[^A-R]' ;;
*) echo "unknown $M $P"; exit 2 ;;
esac
echo "ACCEPT_OK $M $P"
````

### M15 — *Needs you* lists the finish confirmations; CP7 on the handler's answers · 17 prod + 140 test · part: `fast`

- `host/surfacepage/app/app.js`: in `needsYouHTML`, `cs` (the `confirmations` array or `[]`), `n` counts it, the lines are
  `… + ds.map(decisionLineHTML).join('') + cs.map(confirmLineHTML).join('')`, and `more` starts with the
  `Confirm · unavailable: <esc(reason)>` line when `confirmationsUnavailable` is set (S6-D3); `confirmLineHTML(c)` (NEW,
  unexported): refused or not `/^vr-[0-9a-f]{16}$/` → `<div class="nl sub">Confirm · <esc(id)> · refused: <esc(reason or "not a
  confirmation id")></div>`; else the link line (S6-D3, S6-D4). Nothing else in app.js changes; `module.exports` unchanged.
- `host/daemon/surface_needs_you_golden_test.go` (NEW): `TestRequestsJSONConfirmationsGolden` — S3's rig (`newApprovalRig`,
  `apT0`), a ready (`vr-…c2`, w2), a draft (`vr-…c4`, w1) and a digest-refused (`vr-…c6`) request; the handler's answer, then
  the schema-8 answer; `{"v8":…,"v9":…}` byte-compared with the golden (`WORLD_UPDATE_CONFIRMATIONS_GOLDEN=1` writes it).
- `host/surfacepage/testdata/requests-confirmations.golden.json` (NEW, generated, C9).
- `host/surfacepage/app_confirm_test.go` (NEW): **CP7** `TestNeedsYouListsConfirmations` — exact HTML for the handler's
  answer (three lines, the refused one unlinked, two links), the schema-8 answer (said, alone), unavailable + truncation +
  notify (order, escaping), an empty list, after mint requests and decisions (count 5), every field escaped, eight ids that
  are not confirmation ids (15/17 hex, upper case, `xvr-`, no dash, `vr_`, non-hex, a quote-injection suffix) unlinked, and
  no notification.
- Expected numstat after commit **157**.

### M16 — XC2 and XC3 in live Chrome · 0 prod + 133 test · part: `fast` (the executor); `chrome` (the controller's commit)

- `cmd/ailang-worldd/confirm_surface_chrome_test.go` (NEW): `confirmView`, `confirmInChrome` (open raw → wait for `ready`
  with the passkey button → `tap go` → wait ≤ 30 s for `closed`, returning status, answer, the recorded gets and CSP
  violations); **XC2** `TestConfirmPageConfirmsWithAPasskeyInChrome` (the bounded-step control first; `surfaceUnit` with the
  scope at `now + 900` (F4); `cardOrigin` (a real `daemon.New` on an `httptest` loopback); `virtualPasskey` `mark-key` seed 93;
  `storeOwn`; then: the row approved by `webauthn:mark-key` with a `sha256:` decision and `consumed_at` within the run; the
  card's answer and closed status name that decision; one get whose challenge is `VerbMessageHash` and UV `required`; no CSP
  violation; `readBackFinishDecision` over `newClient(origin)` accepts the decision for this scope and key; `pollVerbRequest`
  re-checks the row); **XC3** `TestFinishConfirmOnSurfaceInChrome` (`newSurfaceRig`; its daemon on an `httptest` loopback;
  `virtualPasskey` seed 94 at that origin; `onWait(1)` → `confirmInChrome` on `origin/surface/confirm/<drawnID(1)>`; then: rc
  0; stdout carries `  open <that card>`, `confirmed by webauthn:mark-key (decision <ref>)` and ends with the fake gh's PR URL;
  one finish record naming the row's decision, outcome `done`; the row's `created_at` = the rig's start (the wall clock); the
  terminal never read; no typed approve phase; the decision object `approve` by `webauthn:mark-key`).
- Expected numstat after commit **290**.

### M17 — the carried items · 0 prod + 135 test · parts: `fast`, `cmd1`, `cmd2`, `daemon1`, `daemon2`

- `cmd/ailang-worldd/confirm_surface_clock_test.go` (NEW): `TestAwaitVerbConfirmationEndsOnItsClock` (S6-D9) and
  `TestPollVerbRequestRechecksAtItsOwnClock` (S6-D8).
- `cmd/ailang-worldd/cli_credential_writers_test.go` (NEW): `TestCLICredentialWriterSites` (S6-D10).
- `host/verbs/lookup_test.go`: `"vr-0123456789abcdef x": ""`, `"vr-0123456789abcdef --deny": ""` (kill `E-LK-VR-LEN`) and
  `"--deny vr-0123456789abcdef": ""` (kills its sibling `E-LK-VR-ARGS1`, the id taken from the last argument), on one line.
- `host/surfacepage/app_verb_test.go`: four refused hrefs — `//evil/surface/confirm/vr-…`, `/x/surface/confirm/vr-…`
  (substring), `/surface/confirm/vrx…` (the delimiter), `/surface/confirm/VR-…` (case) (kills `E-JS-INCLUDES`, `E-JS-NODASH`,
  `E-JS-INDEXOF`, `E-JS-CASE`).
- `host/surfacepage/confirm_script_test.go`: the `answers` row "500 that says confirmed" (judge-403 note 2, F6).
- `host/daemon/surface_confirm_test.go`: C13's busy refusal's reason, exact (judge-403 note 3).
- Expected numstat after commit **425**.

### 3.1 Every input and failure class S6's code reaches, traced (design §2.9 and CP7's states; directive retro 4)

| input / state | branch | test |
|---|---|---|
| CP7 pending, draft | `confirmLineHTML` link, `draft` | `TestNeedsYouListsConfirmations` (the handler's `vr-…c4`) |
| CP7 pending, ready | link, `ready` | same (`vr-…c2`) |
| CP7 refused (the daemon's digest refusal) | `c.refused` → no link | same (`vr-…c6`: its reason, no `<a`, no `confirm-` id, no href) |
| CP7 a summarised row with a non-confirmation id | the shape check → refused, no link | same (eight ids: length 15/17, case, prefix, delimiter ×2, hex class, suffix) |
| CP7 unavailable (schema 8), nothing else | the empty branch + the unavailable line | same (the handler's v8 answer) |
| CP7 unavailable with lines, truncation, notify | `more`'s order | same |
| CP7 empty list / no key | `cs = []`, no line | same + `TestNeedsYouListsDecisions` (unchanged "exactly as before" rows) |
| CP7 after mint requests and decisions | order, count | same (count 5, the decision line then the confirmation lines) |
| CP7 hostile episode, repository, head, sha, id, reason, unavailable | `esc()` each | same |
| CP7 notification | none | same (`notifications(G.v9…)` → `[]`) |
| P1 "card confirms" (design §2.9 P1, `confirmed by webauthn:<name>`; then the PR URL) | S5's path, passkey arm | **XC3** (and XC2 for the card and the daemon) |
| P4 the card with a passkey (CP9's success) | `confirm.js` `confirmWithPasskey` → POST → re-read | **XC2**, XC3 |
| P5 the POST, WebAuthn arm, ok | S4's C11 `JudgeVerbAssertion` … C15 | **XC2** (row consumed, decision recorded) |
| P6 `approve vr-… x`, `approve vr-… --deny`, `approve --deny vr-…` | `ApproveTarget` arity | `TestApproveTarget` (+3 rows) |
| P6 a navigate answer with a substring/no-delimiter/upper-case `vr-` target | the allowlist | `TestSurfaceAppVerbs` (+4 refused) |
| P1 the wait's clock never advances | `awaitVerbConfirmation` | `TestAwaitVerbConfirmationEndsOnItsClock` |
| P1 the re-check's clock (a row claiming a later consume) | `pollVerbRequest` `Now` | `TestPollVerbRequestRechecksAtItsOwnClock` |
| P4 a POST answer that is not 200 but says `confirmed` | `confirm.js:198` | `TestConfirmPageShowsRefusals` (+1 row) |
| P5 C13 busy | the reason text | `TestConfirmPOSTRefusesBeforeAnyWrite` (exact reason) |
| CLI writer sites | — (inventory) | `TestCLICredentialWriterSites` |

---

## 4. Mutation table

Enumerated **from the diff**: every new or changed branch of `app.js` (the list, its count, its order, the unavailable line, the refused/no-link branch, the id check per boundary class — anchors, length 15/17, case, delimiter, hex class — every field's `esc()`, the href, the sha cut, the draft word, the time), every carried survivor and its boundary siblings (arity, substring, delimiter, case), each judge-403 note, the CLI writer sites, and two production mutants only a live Chrome test can see (XC2/XC3 are non-vacuous). Each is one exact replacement (`tools/mut.py`; text not found exactly once, or a mutant that does not compile, is **NM** and is re-run in a compiling `-2` form). Each runs against its **named** test, then against the **WHOLE package** of the mutated file (`-run '.'`: every test, static included); a whole-package kill by a static test only is re-run without it (`whole w/o static`). The Chrome mutants run named only, with `WORLD_CHROME` and `WORLD_REQUIRE_PASSKEY_BROWSER=1` (without Chrome the two tests skip, so a whole-package leg would be vacuous). Ledger: `$HOME/.ailang/state/world-iter405/plan-scratch/mutants/ledger.txt` (one line per run; the ids below are the ledger's; the table shows each id's last run per leg — `E-LK-VR-ARGS1` SURVIVED the first draft and was killed after the `--deny vr-…` row was added, both runs in the ledger).

**Base reproductions (C12; tags `base-*`, on a clone at `539cecd0`)**:

| id | file | mutation | command | verdict |
|---|---|---|---|---|
| `E-LK-VR-LEN` | `host/verbs/lookup.go` | `if len(args) == 1 && verbRequestID.MatchString(args[0]) -> if len(args) >= 1 && verbRequestID.MatchString(args[0])` | `go test -count=1 -timeout 900s ./host/verbs/ -run '.'` | **SURVIVED**  |
| `E-JS-INCLUDES` | `host/surfacepage/app/app.js` | `b.href.startsWith('/surface/confirm/vr-') -> b.href.includes('/surface/confirm/vr-')` | `go test -count=1 -timeout 900s ./host/surfacepage/ -run '.'` | **SURVIVED**  |
| `E-JS-NODASH` | `host/surfacepage/app/app.js` | `b.href.startsWith('/surface/confirm/vr-') -> b.href.startsWith('/surface/confirm/vr')` | `go test -count=1 -timeout 900s ./host/surfacepage/ -run '.'` | **SURVIVED**  |
| `E-AWAIT-NOW-STALE` | `cmd/ailang-worldd/episode_finish_surface.go` | `  now := env.now()\n  row, refusal, err := pollVerbRequest(credPath, own, now) ->   now := own.expires - 600\n  row, ref` | `go test -count=1 -timeout 120s ./cmd/ailang-worldd/ -run '^TestAwaitVerbConfirmationEnds$'` | **KILLED** TestAwaitVerbConfirmationEnds |
| `E-RC-INPUT-NOW` | `cmd/ailang-worldd/episode_finish_surface.go` | `ExpiresAt: own.expires, Now: now} -> ExpiresAt: own.expires, Now: row.ConsumedAt}` | `go test -count=1 -timeout 900s ./cmd/ailang-worldd/ -run '.'` | **SURVIVED**  |

| id | file | mutation | named test | named | whole package (killing tests) | whole w/o static | class |
|---|---|---|---|---|---|---|---|
| `E-LK-VR-LEN` | `lookup.go` | `if len(args) == 1 && verbRequestID.MatchString(args[0]) -> if len(args) >= 1 && verbRequestID.MatchString(args[0])` | `^TestApproveTarget$` | killed (TestApproveTarget) | killed (TestApproveTarget) | — | K |
| `E-JS-INCLUDES` | `app.js` | `b.href.startsWith('/surface/confirm/vr-') -> b.href.includes('/surface/confirm/vr-')` | `^TestSurfaceAppVerbs$` | killed (TestSurfaceAppVerbs) | killed (TestSurfaceAppVerbs) | — | K |
| `E-JS-NODASH` | `app.js` | `b.href.startsWith('/surface/confirm/vr-') -> b.href.startsWith('/surface/confirm/vr')` | `^TestSurfaceAppVerbs$` | killed (TestSurfaceAppVerbs) | killed (TestSurfaceAppVerbs) | — | K |
| `E-AWAIT-NOW-STALE` | `episode_finish_surface.go` | `  now := env.now()\n  row, refusal, err := pollVerbRequest(credPath, own, now) ->   now := own.expires - 600\n  row, refusal, err ` | `^TestAwaitVerbConfirmationEndsOnItsClock$` | killed (TestAwaitVerbConfirmationEndsOnItsClock) | killed(timeout) (TestAwaitVerbConfirmationEndsOnItsClock) | — | K |
| `E-RC-INPUT-NOW` | `episode_finish_surface.go` | `ExpiresAt: own.expires, Now: now} -> ExpiresAt: own.expires, Now: row.ConsumedAt}` | `^TestPollVerbRequestRechecksAtItsOwnClock$` | killed (TestPollVerbRequestRechecksAtItsOwnClock) | killed (TestPollVerbRequestRechecksAtItsOwnClock) | — | K |
| `CP-CS-NONE` | `app.js` | `const cs = r && Array.isArray(r.confirmations) ? r.confirmations : []; -> const cs = [];` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-COUNT` | `app.js` | `const n = p.length + ds.length + cs.length; -> const n = p.length + ds.length;` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-JOIN-DROP` | `app.js` | `+ ds.map(decisionLineHTML).join('') + cs.map(confirmLineHTML).join('') -> + ds.map(decisionLineHTML).join('')` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-ORDER` | `app.js` | `+ ds.map(decisionLineHTML).join('') + cs.map(confirmLineHTML).join('') -> + cs.map(confirmLineHTML).join('') + ds.map(decisionLine` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-UNAV-HIDE` | `app.js` | `(r && r.confirmationsUnavailable ? ' -> (false ? '` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-UNAV-ESC` | `app.js` | `unavailable: ${esc(r.confirmationsUnavailable)} -> unavailable: ${r.confirmationsUnavailable}` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-UNAV-AFTER` | `app.js` | `const more = (r && r.confirmationsUnavailable ? '<div class="nl sub">Confirm · unavailable: ${esc(r.confirmationsUnavailable)}</di` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-REFUSED-LINK` | `app.js` | `if (c.refused \|\| !/^vr-[0-9a-f]{16}$/.test(c.id)) -> if (!/^vr-[0-9a-f]{16}$/.test(c.id))` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-SHAPE-NONE` | `app.js` | `if (c.refused \|\| !/^vr-[0-9a-f]{16}$/.test(c.id)) -> if (c.refused)` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-RE-END` | `app.js` | `/^vr-[0-9a-f]{16}$/.test(c.id) -> /^vr-[0-9a-f]{16}/.test(c.id)` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-RE-START` | `app.js` | `/^vr-[0-9a-f]{16}$/.test(c.id) -> /vr-[0-9a-f]{16}$/.test(c.id)` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-RE-LEN15` | `app.js` | `/^vr-[0-9a-f]{16}$/.test(c.id) -> /^vr-[0-9a-f]{15,16}$/.test(c.id)` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-RE-LEN17` | `app.js` | `/^vr-[0-9a-f]{16}$/.test(c.id) -> /^vr-[0-9a-f]{16,17}$/.test(c.id)` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-RE-CASE` | `app.js` | `/^vr-[0-9a-f]{16}$/.test(c.id) -> /^vr-[0-9a-f]{16}$/i.test(c.id)` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-RE-DASH-OPT` | `app.js` | `/^vr-[0-9a-f]{16}$/.test(c.id) -> /^vr-?[0-9a-f]{16}$/.test(c.id)` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-RE-DELIM-ANY` | `app.js` | `/^vr-[0-9a-f]{16}$/.test(c.id) -> /^vr.[0-9a-f]{16}$/.test(c.id)` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-RE-HEX` | `app.js` | `/^vr-[0-9a-f]{16}$/.test(c.id) -> /^vr-[0-9a-z]{16}$/.test(c.id)` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-ESC-ID-REFUSED` | `app.js` | `Confirm · ${esc(c.id)} · refused: -> Confirm · ${c.id} · refused:` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-ESC-REFUSED` | `app.js` | `refused: ${esc(c.refused \|\| 'not a confirmation id')} -> refused: ${c.refused \|\| 'not a confirmation id'}` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-REFUSED-TEXT` | `app.js` | `'not a confirmation id' -> 'not an id'` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-ESC-ID-A` | `app.js` | `id="confirm-${esc(c.id)}" -> id="confirm-${c.id}"` | `^TestNeedsYouListsConfirmations$` | survived (-) | survived (-) | — | E — the id is linked only after `/^vr-[0-9a-f]{16}$/` holds, so it has no character `esc` changes (S6-D4) |
| `CP-ESC-ID-HREF` | `app.js` | `href="/surface/confirm/${esc(c.id)}" -> href="/surface/confirm/${c.id}"` | `^TestNeedsYouListsConfirmations$` | survived (-) | survived (-) | — | E — as `CP-ESC-ID-A` |
| `CP-HREF` | `app.js` | `href="/surface/confirm/${esc(c.id)}" -> href="/surface/approve/${esc(c.id)}"` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-A-ID` | `app.js` | `id="confirm-${esc(c.id)}" -> id="needs-${esc(c.id)}"` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-ESC-EP` | `app.js` | `episode finish ${esc(c.episode)} -> episode finish ${c.episode}` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-ESC-REPO` | `app.js` | `${esc(c.repository)} -> ${c.repository}` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-ESC-HEAD` | `app.js` | `${esc(c.head)}@ -> ${c.head}@` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-ESC-SHA` | `app.js` | `@${esc(String(c.sha \|\| '').slice(0, 12))} -> @${String(c.sha \|\| '').slice(0, 12)}` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-SHA-11` | `app.js` | `String(c.sha \|\| '').slice(0, 12) -> String(c.sha \|\| '').slice(0, 11)` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-SHA-13` | `app.js` | `String(c.sha \|\| '').slice(0, 12) -> String(c.sha \|\| '').slice(0, 13)` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-DRAFT-SWAP` | `app.js` | `${c.draft ? 'draft' : 'ready'} -> ${c.draft ? 'ready' : 'draft'}` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-DRAFT-TRUE` | `app.js` | `${c.draft ? 'draft' : 'ready'} -> ${true ? 'draft' : 'ready'}` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-EXPIRES` | `app.js` | `hm(new Date(c.expires * 1000)) -> hm(new Date(c.expires))` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CP-TEXT` | `app.js` | `Confirm · episode finish ${ -> Confirm · ${` | `^TestNeedsYouListsConfirmations$` | killed (TestNeedsYouListsConfirmations) | killed (TestNeedsYouListsConfirmations) | — | K |
| `CTL-COMMENT` | `app.js` | `// confirmLineHTML is one finish confirmation's -> // confirmLineHTML is one finish confirmations` | `^TestNeedsYouListsConfirmations$` | survived (-) | survived (-) | — | CTL — must stay green: green |
| `E-JS-INDEXOF` | `app.js` | `b.href.startsWith('/surface/confirm/vr-') -> b.href.indexOf('/surface/confirm/vr-') >= 0` | `^TestSurfaceAppVerbs$` | killed (TestSurfaceAppVerbs) | killed (TestSurfaceAppVerbs) | — | K |
| `E-JS-CASE` | `app.js` | `b.href.startsWith('/surface/confirm/vr-') -> b.href.toLowerCase().startsWith('/surface/confirm/vr-')` | `^TestSurfaceAppVerbs$` | killed (TestSurfaceAppVerbs) | killed (TestSurfaceAppVerbs) | — | K |
| `E-JS-ANY` | `app.js` | `b.href.startsWith('/surface/confirm/vr-') -> b.href.startsWith('/surface/confirm/')` | `^TestSurfaceAppVerbs$` | killed (TestSurfaceAppVerbs) | killed (TestSurfaceAppVerbs) | — | K |
| `N2-STATUS` | `confirm.js` | `res.status === 200 && b.confirmed === true -> b.confirmed === true` | `^TestConfirmPageShowsRefusals$` | killed (TestConfirmPageShowsRefusals) | killed (TestConfirmPageShowsRefusals) | — | K |
| `E-LK-VR-ARGS1` | `lookup.go` | `if len(args) == 1 && verbRequestID.MatchString(args[0]) -> if len(args) > 0 && verbRequestID.MatchString(args[len(args)-1])` | `^TestApproveTarget$` | killed (TestApproveTarget) | killed (TestApproveTarget) | — | K |
| `X-RECHECK-FORMAT` | `judge_verb.go` | `check.Format, check.Approver = approval.FormatWebAuthn, a -> check.Format, check.Approver = approval.FormatSSH, a` | `^(TestConfirmPageConfirmsWithAPasskeyInChrome\|TestFinishConfirmOnSurfaceInChrome)$` | killed (TestConfirmPageConfirmsWithAPasskeyInChrome,TestFinishConfirmOnSurfaceInChrome) | — (Chrome: named only) | — | K |
| `X-LINK-ORIGINS` | `episode_finish_surface.go` | `links = append(links, e.Passkey.Origin+path) -> links = append(links)` | — | — | NOT-A-MUTANT(compile) | — | NM (re-run as its `-2` form) |
| `X-CTL-COMMENT` | `episode_finish_surface.go` | `// putFinishVerbRequest is design -> // putFinishVerbRequest is the design` | `^(TestConfirmPageConfirmsWithAPasskeyInChrome\|TestFinishConfirmOnSurfaceInChrome)$` | survived (-) | — (Chrome: named only) | — | CTL — must stay green: green |
| `X-LINK-ORIGINS-2` | `episode_finish_surface.go` | `links = append(links, e.Passkey.Origin+path) -> links = append(links[:0:0], links...)` | `^(TestConfirmPageConfirmsWithAPasskeyInChrome\|TestFinishConfirmOnSurfaceInChrome)$` | killed (TestFinishConfirmOnSurfaceInChrome) | — (Chrome: named only) | — | K |
| `N3-BUSY-TEXT` | `surface_confirm.go` | `the credentials store is busy (a mint or a session verb is writing it); try again -> the credentials store is busy; try again` | `^TestConfirmPOSTRefusesBeforeAnyWrite$` | killed (TestConfirmPOSTRefusesBeforeAnyWrite) | killed (TestConfirmPOSTRefusesBeforeAnyWrite) | — | K |
| `W-SITE-POLL` | `episode_finish_surface.go` | ` ro, err := store.OpenCredentialsReadOnly(credPath)\n if err != nil {\n  return store.VerbRequestRow{}, "", err ->  ro, err := sto` | `^TestCLICredentialWriterSites$` | killed (TestCLICredentialWriterSites) | killed (TestCLICredentialWriterSites) | survived (-) | K-static (an inventory by nature: an extra writer is invisible to behaviour) |
| `W-SITE-DROP` | `episode_finish_surface.go` | ` st, err := store.OpenCredentials(credPath) ->  st, err := store.OpenCredentialsReadOnly(credPath)` | `^TestCLICredentialWriterSites$` | killed (TestCLICredentialWriterSites) | killed (TestAwaitVerbConfirmationEnds,TestAwaitVerbConfirmationEndsOnItsClock,TestCLICredentialWriterSites,TestFinishConfirmOnSurfaceDispatchesTheNamedDecision,TestFinishConfirmOnSurfaceKeepsTheCommit,TestFinishConfirmOnSurfaceOpensThePR,TestFinishConfirmOnSurfaceRefusesAnUnrecordedDecision,TestFinishConfirmOnSurfaceRefusesBeforeTheStore,TestFinishConfirmOnSurfaceRerunsAndTypedPrompt,TestFinishConfirmOnSurfaceResumeReusesTheApproval,TestPollVerbRequestRechecks,TestPollVerbRequestRechecksAtItsOwnClock,TestPutFinishVerbRequestRefusals) | — | K |

**Totals: 49 mutants run (controls and NM excluded), 47 killed (1 of them only by an inventory, by nature), 2 named equivalents, 0 survivors; 1 NM line(s) re-run as `-2` forms; 2 control(s) green.**

Notes on three rows. **`E-AWAIT-NOW-STALE`** (whole package): `--- FAIL: TestAwaitVerbConfirmationEndsOnItsClock (0.01s)` is
printed first; the package then meets its 7 m timeout inside S5's `TestFinishConfirmOnSurfaceRerunsAndTypedPrompt`, whose rig
waits on the fake clock with no wall bound (`mutants/out-E-AWAIT-NOW-STALE-whole.txt`; S6-D9; the rig file is design-§7-frozen,
Appendix C). **`W-SITE-DROP`** is killed behaviourally as well (the put cannot write: twelve behavioural tests fail beside the
inventory). **`X-*`** are the Chrome-only proof that XC2/XC3 see production behaviour: a passkey re-check that judges the wrong
key format (both fail), and a card link for the passkey's origin dropped from the CLI's stdout (XC3 fails); `X-CTL-COMMENT` is
their green control.

---

## 5. Files (final `git diff --name-status 539cecd0 HEAD`)

```
A	cmd/ailang-worldd/cli_credential_writers_test.go
A	cmd/ailang-worldd/confirm_surface_chrome_test.go
A	cmd/ailang-worldd/confirm_surface_clock_test.go
M	host/daemon/surface_confirm_test.go
A	host/daemon/surface_needs_you_golden_test.go
M	host/surfacepage/app/app.js
A	host/surfacepage/app_confirm_test.go
M	host/surfacepage/app_verb_test.go
M	host/surfacepage/confirm_script_test.go
A	host/surfacepage/testdata/requests-confirmations.golden.json
M	host/verbs/lookup_test.go
```

Plus this plan and its JSON under `design_docs/planned/` (the controller's commit). **Production files**:
`host/surfacepage/app/app.js` only — no Go production byte, no `host/store`, `host/daemon`, `host/broker`, `host/approval`,
`host/authority` or `cmd/` production file, no `.ail`, no docs.

---

## 6. Sprint gates (controller, after the last sub-milestone's commit; each measured at the base, §1.1)

One bash shell, never under `nohup`: `bash gates.sh <worktree> $HOME/.ailang/state/world-iter405/gates final <HEAD sha>`
(extract it as `accept.sh` is extracted, with the fence `` ````bash gates.sh ``). Every leg is `tee`'d to
`$HOME/.ailang/state/world-iter405/gates/<leg>.txt`; `ALL_DONE.<sha>` is written last with `FAILS: …`. **Each leg runs in its
own `bash -c` with its own rc; inside a leg, every package run ORs its rc (`r=$(( r | $? ))`) — no `&&` chain can skip a
leg or a package.** The executor sets `GATES_SKIP_VERIFY_GO=1`; the controller runs `verify_go` (`VERIFY_GO_LEGS=plain
AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh`) after the commit. Legs: the **touched-package derivation**
(printed: `git diff --name-only 539cecd0 HEAD -- cmd host` → the `.go` files' directories; the plan's list is
`cmd/ailang-worldd host/daemon host/surfacepage host/verbs`; a difference FAILS), `status`, `verify_ail`, `gofmt`, `build`,
`vet` (`./...`), `fence` (`go test -run '^$' ./...`), `pkgs` (the four touched packages whole, each its own rc), **`race`**
(`host/verbs`, `host/surfacepage` whole; `cmd/ailang-worldd` on the finish/approve/surface/confirm-page tests and the new ones;
`host/daemon` on the approval-surface and verb-route tests incl. `TestRequestsJSON*` — the two big packages whole under
`-race` exceed the leg's budget at load; CI's race job runs them whole), **`race3`** (`-race -count=3` over the two new
clock tests), **`node`** (every page harness incl. CP7, non-skippable, `--- PASS` per name), **`inv`** (the inventories by name
with `-v`, incl. `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots`,
`TestProductionGoSurface`, `TestEveryVerbHasHelp`, `TestFrozenV1RouteTableMatchesMux`, `TestConfirmRoutesRegistered`,
`TestDaemonCredentialWriterSites`, `TestCLICredentialWriterSites`, `TestSurfaceAppPure` (in `node`), the golden test),
**`DAEMON_ROUTES`** (base = final) and **`SURFACE_POSTS: 3`**, **`UNTOUCHED_VIOLATIONS: none`** (§0), **`chrome`** (V2b's
fourteen + XC1 + the needs-you pair + **XC2 + XC3**, `WORLD_CHROME` + `WORLD_REQUIRE_PASSKEY_BROWSER=1`, `--- PASS` per name;
the two host/daemon drills), then `verify_go` unless skipped.

````bash gates.sh
#!/bin/bash
# Row 243 V3 S6 gates (plan §6). usage: gates.sh <tree> <logdir> <base|final> [<sha>]
# The controller: bash gates.sh <worktree> $HOME/.ailang/state/world-iter405/gates final <HEAD sha>
# macOS /bin/bash 3.2: no declare -A. Every leg records its own rc; every leg runs (no set -e, no leg
# chained to another with &&); every leg is tee'd to <logdir>/<leg>.txt. GATES_SKIP_VERIFY_GO=1 skips
# verify_go (the executor sets it; the controller runs it). Never start this under nohup. Writes
# <logdir>/ALL_DONE.<sha> with FAILS: when every leg has finished.
T=$1; L=$2; MODE=$3
mkdir -p "$L"; cd "$T" || exit 2
SHA=${4:-$(git rev-parse HEAD 2>/dev/null || echo unknown)}
B=539cecd0d37930fb7ab0021bc4e529d794c653b2
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_REQUIRE_SSH_KEYGEN WORLD_PASSKEY_RECORD_DIR WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN
: > "$L/legs.txt"
FAILS=""
s0=$(date +%s)
leg() { # leg <name> <command string>: tee'd to $L/<name>.txt, rc recorded in legs.txt
  local name=$1 cmd=$2 s rc
  s=$(date +%s)
  bash -c "$cmd" 2>&1 | tee "$L/$name.txt"
  rc=${PIPESTATUS[0]}
  echo "LEG $name rc=$rc $(( $(date +%s) - s ))s" | tee -a "$L/legs.txt"
  [ "$rc" -eq 0 ] || FAILS="$FAILS $name"
}
need() { # need <leg> <test id>...: each must show --- PASS in <leg>.txt
  local lg=$1 t; shift
  for t in "$@"; do grep -qE -- "^ *--- PASS: $t " "$L/$lg.txt" || { echo "NOT PASSED ($lg): $t" | tee -a "$L/legs.txt"; FAILS="$FAILS ${lg}_missing"; }; done
}
alt() { local x="$*"; x="${x// /|}"; echo "${x:+|$x}"; }
bar() { echo "$*" | tr ' ' '|'; }
# The touched Go packages, derived from the diff (printed so the controller can diff it). The plan's list is
# cmd/ailang-worldd host/daemon host/surfacepage host/verbs; a difference FAILS.
TOUCHED=$(git diff --name-only "$B" HEAD -- cmd host 2>/dev/null | grep '\.go$' | xargs -n1 dirname 2>/dev/null | sort -u | tr '\n' ' ')
echo "TOUCHED_PACKAGES: $TOUCHED" | tee -a "$L/legs.txt"
if [ "$MODE" = final ] && [ "$TOUCHED" != "cmd/ailang-worldd host/daemon host/surfacepage host/verbs " ]; then FAILS="$FAILS touched_packages_differ"; fi
leg status 'git status --porcelain | wc -l'
leg verify_ail './scripts/verify_ail.sh'
leg gofmt 'out=$(gofmt -l cmd host); echo "$out"; test -z "$out"'
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./... > '$L/fence.raw' 2>&1; rc=\$?; cat '$L/fence.raw'; if grep -E '^(FAIL|--- FAIL)' '$L/fence.raw' >/dev/null; then exit 1; fi; exit \$rc"
# the four touched packages whole, each its own rc (one red never skips another)
leg pkgs "r=0; for p in ./host/verbs/ ./host/surfacepage/ ./cmd/ailang-worldd/ ./host/daemon/; do go test -count=1 -timeout 30m \$p; r=\$(( r | \$? )); done; exit \$r"
# -race on every touched package: host/verbs, host/surfacepage whole; cmd/ailang-worldd on the finish, approve,
# surface and Chrome-family tests (the Chrome ones skip here: no WORLD_CHROME); host/daemon on the approval-surface
# and verb-route tests (both whole under -race exceed the leg's budget at load; CI's race job runs them whole)
CRACE='^(TestFinishConfirmOnSurface|TestCanonicalFinishArgv|TestFinishVerbRequestFields|TestSurfaceCardLinks|TestPutFinishVerbRequestRefusals|TestAwaitVerbConfirmationEnds|TestPollVerbRequestRechecks|TestReadBackFinishDecision|TestEpisodeFinish|TestApproveVerbRequest|TestSurfaceSSHSIGThenAwaitMints|TestQuickstartSection|TestConfirmPage|TestCLICredentialWriterSites)'
DRACE='^(TestConfirm|TestConsumeRefusal|TestRequestsJSON|TestApprovePOST|TestApprovalRoutes|TestDecidePOST|TestDecideConcurrent|TestDecideRefuses|TestDecisionRoutes|TestDaemonCredentialWriterSites|TestSurfaceVerb)'
leg race "r=0; go test -race -count=1 ./host/verbs/ ./host/surfacepage/; r=\$(( r | \$? )); go test -race -count=1 -v ./cmd/ailang-worldd/ -run '$CRACE' > '$L/race_cmd.raw' 2>&1; r=\$(( r | \$? )); grep -E '^(--- |ok|FAIL)' '$L/race_cmd.raw'; go test -race -count=1 -v ./host/daemon/ -run '$DRACE' > '$L/race_daemon.raw' 2>&1; r=\$(( r | \$? )); grep -E '^(--- |ok|FAIL)' '$L/race_daemon.raw'; exit \$r"
# the two new wait/clock tests under -race -count=3 (directive 8, 12): the injected clock, SQLite reads
NEW_WAIT="TestAwaitVerbConfirmationEndsOnItsClock TestPollVerbRequestRechecksAtItsOwnClock"
if [ "$MODE" = final ]; then leg race3 "go test -race -count=3 ./cmd/ailang-worldd/ -run '^($(bar $NEW_WAIT))\$'"; fi
# node: every page harness, non-skippable (WORLD_EXEC_NODE is set above; a harness that cannot run FAILS)
NODE="TestSurfaceAppPure TestSurfaceScriptPure TestSurfaceAppVerbs TestNeedsYouListsDecisions TestAddressHintsArePinned TestConfirmPageFailsClosed TestConfirmPageRendersText TestConfirmPageMessage TestConfirmPageStates TestConfirmPageSubmitsAPasskey TestConfirmPageShowsRefusals TestConfirmPageSubmitsAnSSHSIG"
[ "$MODE" = final ] && NODE="$NODE TestNeedsYouListsConfirmations"
leg node "set -o pipefail; go test -v -count=1 ./host/surfacepage/ -run '^($(bar $NODE))\$'"
need node $NODE
if grep -E '^\s*--- (SKIP|FAIL)' "$L/node.txt"; then FAILS="$FAILS node_skip_or_fail"; fi
NEW_CMD="TestPollVerbRequestRechecksAtItsOwnClock TestAwaitVerbConfirmationEndsOnItsClock TestCLICredentialWriterSites"
NEW_D="TestRequestsJSONConfirmationsGolden"
if [ "$MODE" = base ]; then NEW_CMD=""; NEW_D=""; fi
INV_DAEMON="TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestConfirmRoutesRegistered TestDaemonCredentialWriterSites TestSurfaceVerbRouteRegistered TestSurfaceVerbIsReadOnly TestSurfaceVerbOutcomes TestDecidePageAssets TestConfirmPageAssetsServed TestConfirmationCardServed TestRequestsJSONConfirmations TestDecideRefusesReservedScopesByKey TestConfirmPOSTConfirms TestConfirmPOSTRetriesAFailedConsume TestConfirmPOSTRecordsAtItsClock"
INV_BOUNDARY="TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard TestWorkbenchPackageRemainsTransportFree TestPinfetchStaysOutOfStoreReplayAndPublish"
INV_STORE="TestProductionContextRoots TestProductionGoSurface TestStoreSurfaceTakesContext TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext TestCredentialsSchemaMatchesCanonicalDDL TestCredentialsMethodsRequireADeadline TestCredentialsSchema9VerbRequests TestCredentialsNewerThanNineIsRefused"
INV_SP="TestSurfaceAppMatchesWorkingCopy TestSurfaceAppShellAssembly TestConfirmPageFieldTable TestConfirmCardGolden TestConfirmPageAssets"
INV_VERBS="TestLookup TestApproveTarget"
INV_CMD="TestEveryVerbHasHelp TestVerbTableMatchesTheCLISwitch TestFinishPromptIsTheGolden TestEpisodeFinishHelp TestQuickstartSection19MatchesTheBar TestQuickstartSection20MatchesTheFinishFlag TestPollVerbRequestRechecks TestAwaitVerbConfirmationEnds"
leg inv "set -o pipefail; r=0
go test -v -count=1 ./host/daemon/ -run '^($(bar $INV_DAEMON)$(alt $NEW_D))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/boundary/ -run '^($(bar $INV_BOUNDARY))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/store/ -run '^($(bar $INV_STORE))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/surfacepage/ -run '^($(bar $INV_SP))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/verbs/ -run '^($(bar $INV_VERBS))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(bar $INV_CMD)$(alt $NEW_CMD))\$'; r=\$(( r | \$? ))
exit \$r"
need inv $INV_DAEMON $NEW_D $INV_BOUNDARY $INV_STORE $INV_SP $INV_VERBS $INV_CMD $NEW_CMD
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.txt"; then FAILS="$FAILS inv_skip_or_fail"; fi
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.txt")" | tee -a "$L/legs.txt"
# No new daemon route: every mux registration in daemon.go, counted (base = final), and the /surface POSTs (3)
ROUTES=$(grep -cE 'mux\.(HandleFunc|Handle)\(' host/daemon/daemon.go); POSTS=$(grep -cE 'mux\.HandleFunc\("POST /surface/' host/daemon/daemon.go)
BROUTES=$(git show "$B:host/daemon/daemon.go" | grep -cE 'mux\.(HandleFunc|Handle)\(')
echo "DAEMON_ROUTES: $ROUTES (base $BROUTES) SURFACE_POSTS: $POSTS" | tee -a "$L/legs.txt"
if [ "$ROUTES" != "$BROUTES" ] || [ "$POSTS" != 3 ]; then FAILS="$FAILS routes_changed"; fi
# Production bytes: the only production file S6 changes is host/surfacepage/app/app.js. Every other changed
# file is a _test.go or under a testdata/ dir; nothing under host/store, world, packages, scripts, tools, no .ail.
NOPROD=$(git diff --name-only "$B" HEAD | grep -vE '(_test\.go$|/testdata/|^host/surfacepage/app/app\.js$|^design_docs/planned/(w-world-address-bar-verbs-v3s6-sprint-plan\.md|sprint_w-world-address-bar-verbs-v3s6\.json)$)' | tr '\n' ' ')
NOPROD="$NOPROD$(git diff --name-only "$B" HEAD -- host/store world packages scripts tools '*.ail' | tr '\n' ' ')"
echo "UNTOUCHED_VIOLATIONS: ${NOPROD:-none}" | tee -a "$L/legs.txt"
if [ "$MODE" = final ] && [ -n "$NOPROD" ]; then FAILS="$FAILS untouched_violated"; fi
# live Chrome: V2b's fourteen, XC1, the needs-you pair and the two drills stay green on the new app.js; XC2 and XC3 run
XB="TestApprovePageInChrome TestApprovePageRefusesASwappedRequest TestApprovePageSSHSIGInChrome TestApprovePageShowsAHostileValueAsText TestPasskeyPageAbortsAPendingCeremony TestPasskeyPageInChrome TestApproveDecideWebAuthnLinkInChrome TestDecidePageApproveWithPasskey TestDecidePageFailsClosed TestDecidePageSSHSIGPaste TestNeedsYouListsPendingRequests TestNeedsYouNotifiesOnlyWhenAskedOnce TestVerbBarEpisodeReportInChrome TestVerbBarWhyInChrome TestConfirmPageShowsTheFinishPromptInChrome"
[ "$MODE" = final ] && XB="$XB TestConfirmPageConfirmsWithAPasskeyInChrome TestFinishConfirmOnSurfaceInChrome"
XD="TestSurfaceLiveChrome TestSurfaceFollowChrome"
CH="'/Applications/Google Chrome.app/Contents/MacOS/Google Chrome'"
leg chrome "set -o pipefail; r=0; WORLD_CHROME=$CH WORLD_REQUIRE_PASSKEY_BROWSER=1 WORLD_REQUIRE_SSH_KEYGEN=1 go test ./cmd/ailang-worldd/ -run '^($(bar $XB))\$' -count=1 -v; r=\$(( r | \$? )); WORLD_CHROME=$CH go test ./host/daemon/ -run '^($(bar $XD))\$' -count=1 -v; r=\$(( r | \$? )); exit \$r"
need chrome $XB $XD
if grep -E '^\s*--- (SKIP|FAIL)' "$L/chrome.txt"; then FAILS="$FAILS chrome_skip_or_fail"; fi
if [ "${GATES_SKIP_VERIFY_GO:-0}" = 1 ]; then echo "LEG verify_go SKIPPED (GATES_SKIP_VERIFY_GO=1)" | tee -a "$L/legs.txt"
else leg verify_go 'VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh'; fi
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" | tee -a "$L/legs.txt"
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" > "$L/ALL_DONE.$SHA"
````

---

## 7. Executor directive notes

- **No git writes of any kind** (no add, rm, mv, commit, stash, checkout, reset, branch). **`patch -p1 -E` (§2) is the
  only writer** of the tree; no hand edit, no `gofmt -w`, no golden regeneration (the golden ships in M15's patch), no
  fuzz. A failed `--dry-run` or a red acceptance → **stop and report** with the log; do not repair.
- **Never `nohup`, never background a test leg; pass a 600000 ms timeout to every long Bash call.** Each acceptance
  part is one call (C14 measured each). **The executor never launches Chrome**: it never runs the `chrome` part and never sets
  `WORLD_CHROME`.
- **Runs**: **one run, M15 → M16 → M17** (each patch assumes the previous tree; the controller may commit between
  sub-milestones). Per sub-milestone the parts, in order: **M15** `fast`; **M16** `fast`; **M17** `fast`, `cmd1`, `cmd2`,
  `daemon1`, `daemon2` (the four halves are the two big packages whole; each half is one 600 s call). `accept.sh` exports `AILANG_BIN=$HOME/.pinned-ailang/ailang`, `WORLD_EXEC_NODE=/opt/homebrew/bin/node` and
  `WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang`, and unsets `WORLD_CHROME`,
  `WORLD_REQUIRE_PASSKEY_BROWSER` and every golden-update variable.
- **Commit-time test commands** (`commit_ms.sh <M> <msgfile> <test-cmd…>`; each carries its env through `accept.sh`; with
  `T=$HOME/.ailang/state/world-iter405/exec-patches W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter405-row243-v3s6`):
  **M15** `bash $T/accept.sh $W M15 fast`; **M16** `bash $T/accept.sh $W M16 fast && bash $T/accept.sh $W M16 chrome` (the
  `chrome` part sets `WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
  WORLD_REQUIRE_PASSKEY_BROWSER=1` itself and requires `--- PASS` for XC2, XC3, XC1, the needs-you pair, XD1 and the bar's
  report test — a skip is a failure); **M17** `bash $T/accept.sh $W M17 fast`. The `cmd1`/`cmd2`/`daemon1`/`daemon2` parts ran at
  acceptance and are not repeated at commit. `commit_ms.sh` vets the directories of the changed `.go` files only; `app.js` and the golden
  (M15) are not `.go` — M15's `fast` part runs the node harnesses over the new `app.js` and the generator over the new golden.
- After each sub-milestone's acceptance, snapshot **cumulatively against the base** to
  `$HOME/.ailang/state/world-iter405/snap/<M>/`: every file added or changed since `539cecd0` (repo-relative paths), a
  **`DELETED`** file (one repo path per line; **empty for every S6 stage**), a **`STATUS`** file, and `MANIFEST.sha256`
  (the two plan files, `STATUS` and `DELETED` are not in it). **Never an `rsync` of whole directories.** The `find` takes
  only the top-level dirs that exist (M15 has only `host/`; M16 and M17 have `cmd/` and `host/`; no stage has `docs/`),
  rehearsed on all three (C15):
  ```bash
  bash <<'EOF2'
  set -euo pipefail
  M=M15   # the sub-milestone just accepted
  B=539cecd0d37930fb7ab0021bc4e529d794c653b2
  W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter405-row243-v3s6
  D="$HOME/.ailang/state/world-iter405/snap/$M"; rm -rf "${D:?}"; mkdir -p "$D"
  cd "$W"
  { git diff --name-status --no-renames "$B" -- cmd host docs; git ls-files --others --exclude-standard -- cmd host docs | sed 's/^/A\t/'; } | sort -u -t "$(printf '\t')" -k2 > "$D/STATUS"
  awk -F'\t' '$1 != "D" {print $2}' "$D/STATUS" | while read -r p; do mkdir -p "$D/$(dirname "$p")"; cp -p "$p" "$D/$p"; done
  awk -F'\t' '$1 == "D" {print $2}' "$D/STATUS" > "$D/DELETED"
  (cd "$D" && find $(for d in cmd host docs; do [ -d "$d" ] && echo "$d"; done) -type f -print0 | sort -z | xargs -0 shasum -a 256) > "$D/MANIFEST.sha256"
  cat "$D/STATUS"; echo "DELETED: $(wc -l < "$D/DELETED")"; echo "MANIFEST: $(wc -l < "$D/MANIFEST.sha256")"
  EOF2
  ```
  Expected (`STATUS` lines / `DELETED` / MANIFEST): **M15** 4 / 0 / 4; **M16** 5 / 0 / 5; **M17** 11 / 0 / 11. The controller
  copies each snapshot over the worktree, commits it as one commit, and checks the numstat (`git diff -M --numstat 539cecd0
  HEAD -- cmd host docs app | awk '{s+=$1+$2} END {print s}'`): **M15 157, M16 290, M17 425**. A different
  number → stop.
- **Stop conditions**: a dry-run failure; any `--- FAIL`, `--- SKIP` or `FAIL` line in an acceptance (`ACCEPT_RED`; the
  `skipped` step's expected SKIPs/FAILs are inside it and do not count); gofmt output; a numstat mismatch; any file outside §5
  changed (`git status --short` after each patch must equal the cumulative list of §3's paragraphs); a numstat above 1,400.
- **Known load-sensitive tests** (the directive's list) are not this sprint's: a red there is re-run once alone
  (`TestWorkspaceExecConfiguredBindsTheEpisodeHandler`, `TestAGUISlowReaderCutIsResumable`,
  `TestWorkspaceReadThroughProductionWiring`, `TestSeToolsRunBenchmarkTasksEndToEnd` in the `daemon` halves;
  `TestCLIRealSubprocessEpisode`, `TestCapsuleTemplateIsCopiedNotShared` in the `cmd` halves).
- D-WORLD-109: nothing touches `~/.ailang/world`, `~/.ssh`, a real `/dev/tty` or a real key; every store is a
  `t.TempDir()` file.

---

## 8. Risks

- **R1 — Chrome under load.** XC2/XC3 add ≈ 11 s to the `chrome` leg; each driver step is bounded (C8) and the driver is
  reaped by its context (`passkeyBrowserBudget`), so a stuck Chrome fails the test, never hangs the leg (C21).
- **R2 — XC3 holds the CLI's wait while Chrome works.** The CLI is parked in `env.after` (S6-D7); the daemon's wall clock moves
  ≈ 5–10 s during the ceremony, far inside the request's 600 s. If Chrome fails, the CLI's fake-clock wait ends at its deadline
  (≈ 610 polls, < 2 s) and the test reports Chrome's error first.
- **R3 — a daemon without a credentials store now shows a *Needs you* line** "Confirm · unavailable: the verb requests could not
  be read" (design: "shown, not hidden"); every Chrome test that pins Needs-you text ran green on the new app.js (C17's
  `chrome` leg: V2b's fourteen incl. `TestNeedsYouListsPendingRequests`, XC1, the host/daemon drills).
- **R4 — the parked goroutine** in `TestAwaitVerbConfirmationEndsOnItsClock` exists only when the code under test is broken (a
  mutant): the test has already failed.

---

## 9. Acceptance criteria (the judge scores against these)

1. **No new daemon route, no schema change, no `.ail`, `surfaceCSP` unchanged, no CLI behaviour change**: `DAEMON_ROUTES` base
   = final, `SURFACE_POSTS: 3`, `UNTOUCHED_VIOLATIONS: none` (the only production file changed is `app.js`), `verify_ail`
   identical to the base's, `TestDecidePageAssets`, `TestFrozenV1RouteTableMatchesMux`, `TestConfirmRoutesRegistered`,
   `TestDaemonCredentialWriterSites` green unedited.
2. **CP7**: *Needs you* lists each confirmation of the handler's answer with the design's text and a link; **a refused row has no
   link** (`CP-REFUSED-LINK` killed); a non-confirmation id has no link (every `CP-RE-*` boundary mutant killed);
   `confirmationsUnavailable` is shown (`CP-UNAV-HIDE` killed); every field escaped (`CP-ESC-*` killed but the two named
   equivalents on the checked id).
3. **CP7's fixture is the handler's**: `TestRequestsJSONConfirmationsGolden` green; the golden regenerates byte-identical with
   `WORLD_UPDATE_CONFIRMATIONS_GOLDEN=1`.
4. **XC2 `--- PASS` in the chrome leg (not SKIP)**: the decision exists, decided by `webauthn:mark-key`, and the row is consumed
   naming it; the CLI's own re-check and read-back accept it.
5. **XC3 `--- PASS` in the chrome leg (not SKIP)**: the CLI prints the card link, `confirmed by webauthn:mark-key (decision …)`
   and the fake gh's PR URL last; the record names the decision; the terminal is never read.
6. **CI skips XC2/XC3, counted**, and `WORLD_REQUIRE_PASSKEY_BROWSER=1` turns the skip into a failure (`skipped` in M16/M17
   `fast`).
7. **Every judge-404 survivor is killed by a named test, whole package**: `E-LK-VR-LEN` (`TestApproveTarget`), `E-JS-INCLUDES`,
   `E-JS-NODASH` (`TestSurfaceAppVerbs`), `E-RC-INPUT-NOW` (`TestPollVerbRequestRechecksAtItsOwnClock`), `E-AWAIT-NOW-STALE`
   (`TestAwaitVerbConfirmationEndsOnItsClock`, < 1 s); each reproduced on the base first (C12).
8. **Judge-403 notes 2 and 3** killed (`N2-STATUS`, `N3-BUSY-TEXT`); note 4 decided (S6-D11).
9. **`TestCLICredentialWriterSites`** pins the CLI's four writer sites (S5 F8); `W-SITE-POLL`/`W-SITE-DROP` killed.
10. **Every §3.1 row is driven by a named test.**
11. **Size**: `git diff -M --numstat 539cecd0 HEAD -- cmd host docs app` sums to **425** (≤ 1,400; the same unfiltered);
    each sub-milestone's numstat as §7.
12. **Exactly these pre-existing test files change**: `host/verbs/lookup_test.go` (+3 rows), `host/surfacepage/app_verb_test.go`
    (+4 refused hrefs), `host/surfacepage/confirm_script_test.go` (+1 answers row), `host/daemon/surface_confirm_test.go` (the
    busy reason). **None deleted; no `episode_finish*_test.go`; `testdata/passkey_cdp.mjs` unchanged.**
13. **Mutants**: every `K` row killed by its named test and by the whole package; every `CTL` green; the `E` rows named with
    reasons; static-only kills re-run without the static test and recorded; the plan's ids are the ledger's.
14. **Gates**: §6's script `FAILS: none` (or only a directive-listed load-sensitive test, green alone), `verify_go` rc 0.

---

## Appendix C — after this sprint (S6's close; the V3 exit)

- **Moved out: nothing** (S6-D1: 425 ≤ 1,400). No patch is parked.
- **Closed here**: judge-404's survivors `E-LK-VR-LEN`, `E-JS-INCLUDES`, `E-JS-NODASH`, `E-RC-INPUT-NOW`, `E-AWAIT-NOW-STALE`;
  judge-403 notes 2 and 3 (note 4 decided, S6-D11); S5 §1.4 F8 (`TestCLICredentialWriterSites`); design §6 Sprint 6 (CP7, XC2,
  XC3).
- **Residuals carried**: S4's judge-402 note 2 (a day-expired request of the same episode is superseded, not pruned); S3 §1.4 F1
  (the DDL accepts a superseded row with a `decision_ref`); S6-D3's "a confirmation raises no notification" (V4 may add one); the
  finish rig's unbounded fake-clock wait in `episode_finish_surface_rig_test.go` (S6-D9: a broken wait now fails by name first,
  but a whole-package run still meets the package timeout after it — that file is design-§7-frozen); design RV1–RV9 unchanged
  (RV1 Ctrl-C leaves a pending request; RV3 unsandboxed checks refuse surface mode — still UNVERIFIED against Mark's real
  presets; RV4 F1 kept; RV5 no deny on the card).
- **The row 243 V3 exit.** **Built (S1–S6)**: the verb-approval law (`world/verb-request/v1`, the 86-byte `world-verb/v1`
  message, `CheckConfirmation`, `JudgeVerbSSH`/`JudgeVerbAssertion`/`RecheckVerb`), `FinishPromptText` as the terminal's own
  lines, credentials schema 9 `verb_requests`, the read routes and the card `/surface/confirm/<id>` (passkey and SSHSIG),
  `POST /surface/confirm/{id}` (key-decided finish approval, one credentials transaction), `episode finish …
  --confirm-on-surface` (store, print, wait ≤ 610 s, re-check over its own bytes, read back, dispatch), the bar's finish reason
  and `approve vr-`, QUICKSTART §20, *Needs you*'s confirmations, and live-Chrome proof of the whole path with a virtual
  passkey (XC3). **Remains attended**: Mark's own run — parent §10 step 4, the no-terminal exit episode with his real enrolled
  key on his real surface — is **never run by the loop** (D-WORLD-109); QUICKSTART §14–§20 stay *(attended — pending first
  verbatim run)*. **What V4 inherits**: F1 is kept — finish still needs a terminal (RV4, Q1): V4 decides who runs the CLI with
  no terminal (relaxing F1 for surface mode only), launching it from the bar (which still answers `terminal` for `episode
  finish`), `episode start` with a launch preset (P5), handing interactive launches to herdr, per-check confirmation of
  unsandboxed checks (RV3), and whether a pending confirmation should notify (S6-D3).
- The controller records; this plan moves no design doc and edits no charter.

---

## Appendix P — the patches (apply with §2, in order; generated from the compiled scratch stages; do not edit)

### Patch M15 — Needs you lists the finish confirmations; CP7 and its golden

````diff M15
diff --git a/host/daemon/surface_needs_you_golden_test.go b/host/daemon/surface_needs_you_golden_test.go
new file mode 100644
index 00000000..4a841403
--- /dev/null
+++ b/host/daemon/surface_needs_you_golden_test.go
@@ -0,0 +1,48 @@
+package daemon
+
+// Row 243 V3 sprint 6 M15 (design §2.5, §6 CP7): the fixtures of app.js's Needs-you confirmations are
+// GET /surface/requests.json as this handler answers it — a draft and a ready confirmation and a refused
+// one on a schema-9 store, then the same read on a schema-8 store — never hand-typed JSON. The page's
+// node harness (surfacepage TestNeedsYouListsConfirmations) renders them.
+
+import (
+	"bytes"
+	"encoding/json"
+	"os"
+	"path/filepath"
+	"testing"
+)
+
+// confirmGolden is the two answers: regenerate with WORLD_UPDATE_CONFIRMATIONS_GOLDEN=1; otherwise they
+// must equal it byte for byte.
+var confirmGolden = filepath.Join("..", "surfacepage", "testdata", "requests-confirmations.golden.json")
+
+func TestRequestsJSONConfirmationsGolden(t *testing.T) {
+	a := newApprovalRig(t)
+	a.putVerb("vr-00000000000000c2", "w2", apT0-50, false)
+	a.putVerb("vr-00000000000000c4", "w1", apT0-40, true)
+	edited := a.putVerb("vr-00000000000000c6", "w6", apT0-20, true)
+	setVerbColumn(t, a.d.credentialsPath, edited.RequestID, "request", append(edited.Request, ' '))
+	get := func() json.RawMessage {
+		rec := a.serve("GET", "/surface/requests.json", apLoopback, nil)
+		if rec.Code != 200 || !json.Valid(rec.Body.Bytes()) {
+			t.Fatalf("requests.json = %d %s", rec.Code, rec.Body.String())
+		}
+		return bytes.TrimSpace(rec.Body.Bytes())
+	}
+	v9 := get()
+	setVersion(t, a.d.credentialsPath, `DROP TABLE verb_requests`, `PRAGMA user_version = 8`)
+	want, err := json.Marshal(map[string]json.RawMessage{"v9": v9, "v8": get()})
+	if err != nil {
+		t.Fatal(err)
+	}
+	want = append(want, '\n')
+	if os.Getenv("WORLD_UPDATE_CONFIRMATIONS_GOLDEN") == "1" {
+		if err := os.WriteFile(confirmGolden, want, 0o644); err != nil {
+			t.Fatal(err)
+		}
+	}
+	if got, err := os.ReadFile(confirmGolden); err != nil || !bytes.Equal(got, want) {
+		t.Errorf("%s is not the handler's answers (WORLD_UPDATE_CONFIRMATIONS_GOLDEN=1 regenerates it): %v", confirmGolden, err)
+	}
+}
diff --git a/host/surfacepage/app/app.js b/host/surfacepage/app/app.js
index 36a1a366..a3c2ff22 100644
--- a/host/surfacepage/app/app.js
+++ b/host/surfacepage/app/app.js
@@ -426,16 +426,20 @@
   // then the notification opt-in for this browser's permission (design §2.9).
   // Row 201 A4 (design §2.8): then the pending Human.Approve decisions of the same read, each a link to
   // its decision card; a reserved one names its typed owner and a refused one its reason, with no link.
+  // Row 243 V3 (design §2.5): then the pending finish confirmations, each a link to its card; a refused one
+  // shows its reason with no link; confirmationsUnavailable is said, never hidden.
   function needsYouHTML(r, now, permission) {
     const p = r && Array.isArray(r.pending) ? r.pending : [];
     const ds = r && Array.isArray(r.decisions) ? r.decisions : [];
-    const n = p.length + ds.length;
+    const cs = r && Array.isArray(r.confirmations) ? r.confirmations : [];
+    const n = p.length + ds.length + cs.length;
     const head = `<div class="nh" style="color:${n ? 'var(--youShell)' : 'var(--sub)'}">Needs you · ${n}</div>`;
-    const more = r && r.decisionsTruncated ? '<div class="nl sub">Older approval requests are not listed; decide them at a terminal with ailang-worldd approve &lt;request-ref&gt;</div>' : '';
+    const more = (r && r.confirmationsUnavailable ? `<div class="nl sub">Confirm · unavailable: ${esc(r.confirmationsUnavailable)}</div>` : '')
+      + (r && r.decisionsTruncated ? '<div class="nl sub">Older approval requests are not listed; decide them at a terminal with ailang-worldd approve &lt;request-ref&gt;</div>' : '');
     if (!n) return `<div class="nsec" data-needs>${head}<div class="nl sub">${esc(r && r.unavailable || 'Nothing needs you')}</div>${more}${notifyHTML(permission)}</div>`;
     return `<div class="nsec" data-needs>${head}` + p.map(q => q.refused ? `<div class="nl sub">${esc(q.id)} · refused: ${esc(q.refused)}</div>`
       : `<div class="nl"><a id="needs-${esc(q.id)}" href="/surface/approve/${esc(q.id)}"><b>${esc(q.episode)}</b> · ${esc(q.agent)} · ${q.grants} grant${q.grants === 1 ? '' : 's'} · expires in ${Math.max(0, Math.ceil((q.expires - now) / 60))} min</a>`
-        + (q.recorded ? ` · approved by ${esc(q.recorded)}, waiting for the minter` : '') + '</div>').join('') + ds.map(decisionLineHTML).join('')
+        + (q.recorded ? ` · approved by ${esc(q.recorded)}, waiting for the minter` : '') + '</div>').join('') + ds.map(decisionLineHTML).join('') + cs.map(confirmLineHTML).join('')
       + more + notifyHTML(permission) + '</div>';
   }
   // RESERVED_OWNERS is the typed command that owns a reserved decision (design §2.6 L1).
@@ -447,6 +451,13 @@
     const ref = String(x.request || '').replace(/^sha256:/, '');
     return `<div class="nl"><a id="decide-${esc(ref)}" href="/surface/decide/sha256:${esc(ref)}">Decide · ${esc(x.scope)} · requested by ${esc(x.requester)} · deadline ${hm(new Date(x.deadline * 1000))}</a></div>`;
   }
+  // confirmLineHTML is one finish confirmation's Needs-you line (row 243 V3 §2.5): a link to its card only for
+  // a summarised row whose id is a confirmation id (store.VerbRequestIDPattern); otherwise its reason, no link.
+  function confirmLineHTML(c) {
+    if (c.refused || !/^vr-[0-9a-f]{16}$/.test(c.id)) return `<div class="nl sub">Confirm · ${esc(c.id)} · refused: ${esc(c.refused || 'not a confirmation id')}</div>`;
+    return `<div class="nl"><a id="confirm-${esc(c.id)}" href="/surface/confirm/${esc(c.id)}">Confirm · episode finish ${esc(c.episode)} · ${esc(c.repository)} ${esc(c.head)}@${esc(String(c.sha || '').slice(0, 12))}`
+      + ` · ${c.draft ? 'draft' : 'ready'} · expires ${hm(new Date(c.expires * 1000))}</a></div>`;
+  }
   // notifyHTML is the opt-in: "Notify me" while this browser has not been asked (permission "default"),
   // its answer once it has, and nothing where the browser has no notifications.
   function notifyHTML(permission) {
diff --git a/host/surfacepage/app_confirm_test.go b/host/surfacepage/app_confirm_test.go
new file mode 100644
index 00000000..a4df44cc
--- /dev/null
+++ b/host/surfacepage/app_confirm_test.go
@@ -0,0 +1,91 @@
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
+// TestNeedsYouListsConfirmations (row 243 V3 S6, design §2.5, §6 CP7) runs app.js's pure part under node
+// on the daemon's own GET /surface/requests.json answers (testdata/requests-confirmations.golden.json,
+// written by host/daemon TestRequestsJSONConfirmationsGolden): Needs you lists each pending finish
+// confirmation after the decisions as a link to its card; a refused row shows its reason and no link, as
+// does a row whose id is not a confirmation id; confirmationsUnavailable is said; every value is escaped.
+// It skips only when no node exists; with WORLD_EXEC_NODE set, a harness that cannot run fails.
+func TestNeedsYouListsConfirmations(t *testing.T) {
+	node := os.Getenv("WORLD_EXEC_NODE")
+	if node == "" {
+		var err error
+		node, err = exec.LookPath("node")
+		if err != nil {
+			t.Skip("node unavailable locally")
+		}
+	}
+	golden, err := os.ReadFile(filepath.Join("testdata", "requests-confirmations.golden.json"))
+	if err != nil {
+		t.Fatal(err)
+	}
+	dir := t.TempDir()
+	for name, data := range map[string][]byte{"app.js": AppScript, "harness.js": []byte(needsYouConfirmationsHarness), "golden.json": golden} {
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
+	if !strings.Contains(string(out), "NEEDS YOU CONFIRMATIONS PASS") {
+		t.Fatalf("harness manifest missing: %s", out)
+	}
+}
+
+const needsYouConfirmationsHarness = `
+const assert = require('node:assert/strict');
+const app = require('./app.js');
+const G = JSON.parse(require('node:fs').readFileSync('golden.json', 'utf8'));
+const HEAD = (n) => '<div class="nsec" data-needs><div class="nh" style="color:var(--' + (n ? 'youShell' : 'sub') + ')">Needs you · ' + n + '</div>';
+const refused = '<div class="nl sub">Confirm · vr-00000000000000c6 · refused: the stored request does not match its digest</div>';
+const link = (id, ep, word) => '<div class="nl"><a id="confirm-' + id + '" href="/surface/confirm/' + id + '">Confirm · episode finish ' + ep
+  + ' · sunholo-data/ailang-world world/' + ep + '@0a0a0a0a0a0a · ' + word + ' · expires 08:09</a></div>';
+const lines = refused + link('vr-00000000000000c4', 'w1', 'draft') + link('vr-00000000000000c2', 'w2', 'ready');
+// The handler's answer: one line each, newest first, after the decisions; the refused one has no link.
+const h = app.needsYouHTML(G.v9, G.v9.now);
+assert.equal(h, HEAD(3) + lines + '</div>');
+assert(!h.includes('/surface/confirm/vr-00000000000000c6') && !h.includes('confirm-vr-00000000000000c6') && h.split('<a ').length === 3, h);
+// Below schema 9: the reason is said, not hidden — alone, and after the lines (escaped).
+assert.equal(app.needsYouHTML(G.v8, G.v8.now), HEAD(0) + '<div class="nl sub">Nothing needs you</div><div class="nl sub">Confirm · unavailable: '
+  + 'the credentials store predates verb requests: restart ` + "`ailang-worldd serve`" + ` with this build once</div></div>');
+assert.equal(app.needsYouHTML({...G.v9, confirmationsUnavailable: '<i>u</i>', decisionsTruncated: true}, 0, 'granted'), HEAD(3) + lines
+  + '<div class="nl sub">Confirm · unavailable: &lt;i&gt;u&lt;/i&gt;</div><div class="nl sub">Older approval requests are not listed; decide them at a terminal with ailang-worldd approve &lt;request-ref&gt;</div>'
+  + '<div class="nl sub" data-notify>Notifications on</div></div>');
+// None listed: exactly the empty Needs you.
+assert.equal(app.needsYouHTML({...G.v9, confirmations: []}, 0), HEAD(0) + '<div class="nl sub">Nothing needs you</div></div>');
+// After the mint requests and the decisions, counted with them.
+{ const Q = 'ab'.repeat(32), dec = {request: 'sha256:' + Q, scope: 'worktree:1', requester: 'agent-1', deadline: 50520, reserved: '', refusal: ''};
+  const m = app.needsYouHTML({...G.v9, pending: [{id: 'mr-1', episode: 'e', agent: 'a', grants: 1, expires: 60}], decisions: [dec]}, 0);
+  assert(m.startsWith(HEAD(5) + '<div class="nl"><a id="needs-mr-1"'), m);
+  assert(m.endsWith('<a id="decide-' + Q + '" href="/surface/decide/sha256:' + Q + '">Decide · worktree:1 · requested by agent-1 · deadline 14:02</a></div>' + lines + '</div>'), m); }
+// Every value escaped where it lands (the sha cut to its first 12 characters, then escaped).
+{ const real = G.v9.confirmations[1];
+  assert.equal(app.needsYouHTML({confirmations: [{...real, episode: '<e>', repository: 'a&b', head: '"h"', sha: '<sha>0123456789'}]}, 0), HEAD(1)
+    + '<div class="nl"><a id="confirm-vr-00000000000000c4" href="/surface/confirm/vr-00000000000000c4">Confirm · episode finish &lt;e&gt; · a&amp;b &quot;h&quot;@&lt;sha&gt;0123456 · draft · expires 08:09</a></div></div>');
+  assert.equal(app.needsYouHTML({confirmations: [{...G.v9.confirmations[0], id: '<x>', refused: '<i>r</i>'}]}, 0), HEAD(1)
+    + '<div class="nl sub">Confirm · &lt;x&gt; · refused: &lt;i&gt;r&lt;/i&gt;</div></div>');
+  // A summarised row whose id is not exactly vr- and 16 lower-case hex: its id, escaped, and no link.
+  for (const id of ['vr-' + 'a'.repeat(15), 'vr-' + 'a'.repeat(17), 'vr-' + 'A'.repeat(16), 'xvr-' + 'a'.repeat(16), 'vr' + 'a'.repeat(16), 'vr_' + 'a'.repeat(16),
+    'vr-' + 'g'.repeat(16), 'vr-' + 'a'.repeat(16) + '"><b>']) {
+    const x = app.needsYouHTML({confirmations: [{...real, id}]}, 0);
+    assert.equal(x, HEAD(1) + '<div class="nl sub">Confirm · ' + app.esc(id) + ' · refused: not a confirmation id</div></div>', id); } }
+// A confirmation raises no notification (design §2.5 lists it; it does not notify).
+assert.deepEqual(app.notifications(G.v9, new Set(), 'granted').notes, []);
+console.log('NEEDS YOU CONFIRMATIONS PASS');
+`
diff --git a/host/surfacepage/testdata/requests-confirmations.golden.json b/host/surfacepage/testdata/requests-confirmations.golden.json
new file mode 100644
index 00000000..a357d947
--- /dev/null
+++ b/host/surfacepage/testdata/requests-confirmations.golden.json
@@ -0,0 +1 @@
+{"v8":{"v":1,"now":1800000000,"pending":[],"decided":[],"unavailable":"","decisions":[],"decisionsTruncated":false,"confirmations":[],"confirmationsUnavailable":"the credentials store predates verb requests: restart `ailang-worldd serve` with this build once"},"v9":{"v":1,"now":1800000000,"pending":[],"decided":[],"unavailable":"","decisions":[],"decisionsTruncated":false,"confirmations":[{"id":"vr-00000000000000c6","expires":1800000580,"refused":"the stored request does not match its digest"},{"id":"vr-00000000000000c4","episode":"w1","repository":"sunholo-data/ailang-world","head":"world/w1","sha":"0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a","draft":true,"expires":1800000560},{"id":"vr-00000000000000c2","episode":"w2","repository":"sunholo-data/ailang-world","head":"world/w2","sha":"0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a","draft":false,"expires":1800000550}],"confirmationsUnavailable":""}}
````

### Patch M16 — XC2 and XC3 in live Chrome

````diff M16
diff --git a/cmd/ailang-worldd/confirm_surface_chrome_test.go b/cmd/ailang-worldd/confirm_surface_chrome_test.go
new file mode 100644
index 00000000..69df4cac
--- /dev/null
+++ b/cmd/ailang-worldd/confirm_surface_chrome_test.go
@@ -0,0 +1,133 @@
+package main
+
+// Row 243 V3 sprint 6 M16 (design §6 XC2, XC3): a finish confirmation given by a passkey in a live Chrome,
+// on the A2/A3 driver (testdata/passkey_cdp.mjs, unchanged; its virtual authenticator holds approvaltest's
+// key, enrolled by virtualPasskey), gated as XB1 (WORLD_CHROME, WORLD_REQUIRE_PASSKEY_BROWSER=1; skipped
+// and counted without them). The card is S3/S4's /surface/confirm/<id> served by a real daemon.New on the
+// loopback; the request is stored by the CLI's production put (XC2) or by the CLI itself (XC3).
+
+import (
+	"encoding/base64"
+	"encoding/hex"
+	"fmt"
+	"net/http/httptest"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// confirmView is the confirmation card after its passkey ceremony settles.
+type confirmView struct {
+	State, Status, Answer string
+	Gets                  []cardGet
+	Violations            []string
+}
+
+// confirmInChrome opens card, waits until it offers the passkey, clicks it and waits (each step bounded by
+// the driver) until the card reads its request again, closed. It never calls t.Fatal (XC3 calls it from
+// the CLI's wait).
+func confirmInChrome(b *passkeyBrowser, card string) (confirmView, error) {
+	var v confirmView
+	if err := b.call(nil, "openCard", map[string]any{"url": card, "raw": true}); err != nil {
+		return v, err
+	}
+	if err := b.call(nil, "waitFor", map[string]any{"expression": `document.body.dataset.state === 'ready' && !document.getElementById('go').hidden`}); err != nil {
+		return v, err
+	}
+	if err := b.call(nil, "tap", map[string]any{"id": "go"}); err != nil {
+		return v, err
+	}
+	err := b.call(&v, "waitFor", map[string]any{"ms": 30000, "expression": `document.body.dataset.state === 'closed' && { state: document.body.dataset.state,
+		status: document.getElementById('status').textContent, answer: document.getElementById('answer').textContent, gets: window.__gets || [], violations: window.__cspv || [] }`})
+	return v, err
+}
+
+// XC2: the virtual authenticator confirms a request stored by the CLI's own put on the card a real daemon
+// serves; the daemon records a world decision decided by webauthn:mark-key and consumes the row naming it;
+// the CLI's own re-check (a) and read-back (b) accept both. Control: a driver step is bounded.
+func TestConfirmPageConfirmsWithAPasskeyInChrome(t *testing.T) {
+	b := startPasskeyBrowser(t)
+	b.do(t, nil, "tab", map[string]any{"authenticator": "uv", "inject": "getRecorder"})
+	if s := time.Now(); b.call(nil, "waitFor", map[string]any{"expression": "false", "ms": 300}) == nil || time.Since(s) > 10*time.Second {
+		t.Fatalf("a waitFor that never holds did not fail within its bound (%s)", time.Since(s))
+	}
+	now := time.Now().Unix()
+	creds, in := surfaceUnit(t)
+	in.scope.ExpiresAt = now + 900 // a finish scope expires 900 s after the CLI builds it
+	in.scopeText = broker.FormatFinishApprovalScope(in.scope)
+	origin := cardOrigin(t, creds, "", nil)
+	virtualPasskey(t, b, creds, origin, "mark-key", 93, now-100)
+	own := storeOwn(t, in, "vr-00000000000000c7", 0xc7, now)
+	v, err := confirmInChrome(b, origin+"/surface/confirm/"+own.id)
+	if err != nil {
+		t.Fatal(err)
+	}
+	ro, err := store.OpenCredentialsReadOnly(creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	row, ok, err := ro.VerbRequest(boundedTestContext(t), own.id)
+	_ = ro.Close()
+	if err != nil || !ok || row.Outcome != store.VerbOutcomeApproved || row.ApprovedBy != "webauthn:mark-key" || !strings.HasPrefix(row.DecisionRef, "sha256:") ||
+		row.ConsumedAt < now || row.ConsumedAt > time.Now().Unix() {
+		t.Fatalf("row %+v (%v %v)", row, ok, err)
+	}
+	h, _ := approval.VerbMessageHash(own.digest, own.challenge, own.expires)
+	if v.Answer != "Confirmed by webauthn:mark-key: decision "+row.DecisionRef+"; the waiting episode finish opens the pull request." ||
+		v.Status != "Already confirmed by webauthn:mark-key: decision "+row.DecisionRef+"; nothing to confirm here." || len(v.Violations) != 0 ||
+		len(v.Gets) != 1 || v.Gets[0].Challenge != base64.RawURLEncoding.EncodeToString(h[:]) || v.Gets[0].UserVerification != "required" {
+		t.Fatalf("the card shows %+v; want the confirmation of %s with challenge %s", v, row.DecisionRef, hex.EncodeToString(h[:]))
+	}
+	if refusal := readBackFinishDecision(worldReader{ctx: boundedTestContext(t), c: newClient(origin)}, row.DecisionRef, "webauthn:mark-key", in.scopeText); refusal != "" {
+		t.Fatal(refusal)
+	}
+	if _, refusal, err := pollVerbRequest(creds, own, time.Now().Unix()); refusal != "" || err != nil {
+		t.Fatalf("the CLI's re-check of the passkey's confirmation: %q %v", refusal, err)
+	}
+}
+
+// XC3: episode finish --confirm-on-surface (the real CLI path on S5's rig, fake gh) stores its request and
+// prints the card's link for the passkey's origin; the page confirms with the virtual authenticator while
+// the CLI waits; the CLI re-checks, reads the decision back and prints the fake gh's PR URL. The terminal
+// is never read; the rig's clock starts at the wall clock the daemon reads.
+func TestFinishConfirmOnSurfaceInChrome(t *testing.T) {
+	b := startPasskeyBrowser(t)
+	b.do(t, nil, "tab", map[string]any{"authenticator": "uv"})
+	r := newSurfaceRig(t)
+	srv := httptest.NewServer(r.d.Handler())
+	t.Cleanup(srv.Close)
+	origin := strings.Replace(srv.URL, "127.0.0.1", "localhost", 1)
+	virtualPasskey(t, b, r.creds, origin, "mark-key", 94, r.now-100)
+	id, card := drawnID(1), ""
+	var v confirmView
+	var chromeErr error
+	r.onWait = func(n int64) {
+		if n == 1 {
+			card = origin + "/surface/confirm/" + id
+			v, chromeErr = confirmInChrome(b, card)
+		}
+	}
+	rc, out, errs := r.run()
+	if chromeErr != nil || rc != exitOK {
+		t.Fatalf("chrome: %v; rc %d\nstdout %s\nstderr %s", chromeErr, rc, out, errs)
+	}
+	row := r.row(id)
+	pr := r.f.finishRecords[0].PR.URL // the fake gh's own answer
+	if !strings.Contains(out, "\n  open "+card+"\n") || !strings.Contains(out, fmt.Sprintf("\nconfirmed by webauthn:mark-key (decision %s)\n", row.DecisionRef)) ||
+		!strings.HasSuffix(out, "\n"+pr+"\n") || pr == "" || !strings.HasPrefix(v.Answer, "Confirmed by webauthn:mark-key: decision "+row.DecisionRef) {
+		t.Fatalf("stdout\n%s\ncard %+v; want the link %s, the confirmation and the PR %q last", out, v, card, pr)
+	}
+	if rec := r.f.finishRecords; len(rec) != 1 || rec[0].Approval != row.DecisionRef || rec[0].Outcome != "done" || row.ApprovedBy != "webauthn:mark-key" ||
+		row.CreatedAt != r.start || r.tty.reads != 0 || len(r.f.scopes) != 0 {
+		t.Fatalf("records %+v, row %+v, terminal reads %d, typed approve phases %d", rec, row, r.tty.reads, len(r.f.scopes))
+	}
+	var decision struct{ Decision, DecidedBy string }
+	r.object(row.DecisionRef, &decision)
+	if decision.Decision != "approve" || decision.DecidedBy != "webauthn:mark-key" {
+		t.Fatalf("decision %+v", decision)
+	}
+}
````

### Patch M17 — the carried items

````diff M17
diff --git a/cmd/ailang-worldd/cli_credential_writers_test.go b/cmd/ailang-worldd/cli_credential_writers_test.go
new file mode 100644
index 00000000..42bbe4c2
--- /dev/null
+++ b/cmd/ailang-worldd/cli_credential_writers_test.go
@@ -0,0 +1,57 @@
+package main
+
+// Row 243 V3 sprint 6 (S5 §1.4 F8): the CLI's production code opens the credentials writer at exactly
+// these sites, as host/daemon's TestDaemonCredentialWriterSites pins the daemon's: a session's open
+// (mints, approver verbs), launch's record, finish's revoke and S5's verb-request put. Control: the
+// read-only opens are found by the same walk.
+
+import (
+	"go/ast"
+	"go/parser"
+	"go/token"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+)
+
+func TestCLICredentialWriterSites(t *testing.T) {
+	files, err := filepath.Glob("*.go")
+	if err != nil {
+		t.Fatal(err)
+	}
+	sites, readOnly := map[string]string{}, 0
+	for _, f := range files {
+		if strings.HasSuffix(f, "_test.go") {
+			continue
+		}
+		file, err := parser.ParseFile(token.NewFileSet(), f, nil, 0)
+		if err != nil {
+			t.Fatal(err)
+		}
+		for _, decl := range file.Decls {
+			fn, ok := decl.(*ast.FuncDecl)
+			if !ok {
+				continue
+			}
+			ast.Inspect(fn, func(n ast.Node) bool {
+				if sel, ok := n.(*ast.SelectorExpr); ok {
+					if pkg, isIdent := sel.X.(*ast.Ident); isIdent && pkg.Name == "store" {
+						switch sel.Sel.Name {
+						case "OpenCredentials", "UpgradeCredentials":
+							sites[fn.Name.Name] += sel.Sel.Name
+						case "OpenCredentialsReadOnly":
+							readOnly++
+						}
+					}
+				}
+				return true
+			})
+		}
+	}
+	want := map[string]string{"openCredentialsForSession": "OpenCredentials", "launchEpisode": "OpenCredentials", "finishRevoke": "OpenCredentials", "putFinishVerbRequest": "OpenCredentials"}
+	if !reflect.DeepEqual(sites, want) || readOnly < 20 {
+		t.Fatalf("writer sites %v (read-only opens %d), want exactly %v and at least 20 read-only opens", sites, readOnly, want)
+	}
+	t.Logf("read-only opens: %d", readOnly)
+}
diff --git a/cmd/ailang-worldd/confirm_surface_clock_test.go b/cmd/ailang-worldd/confirm_surface_clock_test.go
new file mode 100644
index 00000000..d0949569
--- /dev/null
+++ b/cmd/ailang-worldd/confirm_surface_clock_test.go
@@ -0,0 +1,67 @@
+package main
+
+// Row 243 V3 sprint 6 (judge-404 survivors E-AWAIT-NOW-STALE and E-RC-INPUT-NOW, design §2.7 steps 5 and
+// 6(a)): the confirm-on-surface wait reads its own injected clock on every poll and ends at its deadline on
+// that clock, and the re-check judges the approver at this process's clock, never the row's. This file
+// sorts before the finish rig's, so a wait that never ends fails here, by name, in under a second.
+
+import (
+	"fmt"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+// The wait polls at expires_at+9 and at +10 (its deadline) and refuses there, on the injected clock: once
+// that clock passes the deadline the test fails and the wait is parked (never spun), whatever the wall clock.
+func TestAwaitVerbConfirmationEndsOnItsClock(t *testing.T) {
+	creds, in := surfaceUnit(t)
+	own := storeOwn(t, in, "vr-00000000000000e1", 0xe1, sfT0)
+	deadline := own.expires + int64(verbLandGrace/time.Second)
+	now, waits, late := deadline-1, 0, make(chan int64, 1)
+	env := sessionEnv{now: func() int64 { return now }, after: func(d time.Duration) <-chan time.Time {
+		waits++
+		if now += int64(d / time.Second); now > deadline+1 {
+			late <- now
+			return nil // parks the wait: a nil channel never fires
+		}
+		ch := make(chan time.Time, 1)
+		ch <- time.Unix(now, 0)
+		return ch
+	}}
+	done := make(chan string, 1)
+	go func() { _, refusal := awaitVerbConfirmation(env, creds, own); done <- refusal }()
+	select {
+	case refusal := <-done:
+		if want := "no confirmation of request vr-00000000000000e1 by " + approverTime(deadline); refusal != want || waits != 1 || now != deadline {
+			t.Fatalf("%q after %d waits at %d, want %q after 1 wait at %d", refusal, waits, now, want, deadline)
+		}
+	case at := <-late:
+		t.Fatalf("the wait was still polling at %d on its own clock, past its deadline %d", at, deadline)
+	case <-time.After(20 * time.Second):
+		t.Fatal("the wait did not end within 20 s of wall time")
+	}
+}
+
+// The approver is judged at the CLI's own clock: a row that says it was consumed at +50 does not admit a
+// key enrolled at +20 to a poll at +19; the same row re-checks at +20 (P5's boundary, the CLI's clock).
+func TestPollVerbRequestRechecksAtItsOwnClock(t *testing.T) {
+	for _, c := range []struct {
+		at   int64
+		want string
+	}{{sfT0 + 19, approval.ErrApproverClockBehind.Error()}, {sfT0 + 20, ""}} {
+		t.Run(fmt.Sprint(c.at-sfT0), func(t *testing.T) {
+			creds, in := surfaceUnit(t)
+			own := storeOwn(t, in, "vr-00000000000000e2", 0xe2, sfT0)
+			approveByHand(t, creds, own.id, verbSig(t, sigSigner, approval.VerbNamespace, own.digest, own.challenge, own.expires), "ssh:mark-mac", sfT0+50)
+			credSQL(t, creds, fmt.Sprintf("UPDATE approvers SET enrolled_at = %d WHERE name = 'mark-mac'", sfT0+20))
+			row, refusal, err := pollVerbRequest(creds, own, c.at)
+			if err != nil || (c.want == "") != (refusal == "") || !strings.Contains(refusal, c.want) || (c.want == "" && row.ConsumedAt != sfT0+50) ||
+				(c.want != "" && !strings.HasPrefix(refusal, "the confirmation of request vr-00000000000000e2 does not verify: ")) {
+				t.Fatalf("poll at %d: %q %v, want %q", c.at, refusal, err, c.want)
+			}
+		})
+	}
+}
diff --git a/host/daemon/surface_confirm_test.go b/host/daemon/surface_confirm_test.go
index 8f7a2826..988909a2 100644
--- a/host/daemon/surface_confirm_test.go
+++ b/host/daemon/surface_confirm_test.go
@@ -275,8 +275,8 @@ func TestConfirmPOSTRefusesBeforeAnyWrite(t *testing.T) {
 	}
 	rec := a.confirm(other.RequestID, otherBody)
 	_ = held.Close()
-	if requireRefusal(t, "C13 busy", rec, 503, "Busy"); rec.Header().Get("Retry-After") != "1" {
-		t.Errorf("C13 Retry-After %q", rec.Header().Get("Retry-After"))
+	if got := requireRefusal(t, "C13 busy", rec, 503, "Busy"); rec.Header().Get("Retry-After") != "1" || got.Reason != "the credentials store is busy (a mint or a session verb is writing it); try again" {
+		t.Errorf("C13 Retry-After %q, reason %q", rec.Header().Get("Retry-After"), got.Reason)
 	}
 	if after := a.confirmState(); after != before {
 		t.Fatalf("C12/C13 wrote:\n%s\n%s", before, after)
diff --git a/host/surfacepage/app_verb_test.go b/host/surfacepage/app_verb_test.go
index 16f18955..b9f5334c 100644
--- a/host/surfacepage/app_verb_test.go
+++ b/host/surfacepage/app_verb_test.go
@@ -123,7 +123,8 @@ async function run(raw, fetchImpl) {
     assert.equal(B(q).href, href); assert.deepEqual(act, { assign: href }); assert.deepEqual(log, [['clear'], ['assign', href]]);
   }
   // VP5: any other target is refused.
-  for (const href of ['https://evil/x', '/surface/model.json', '//evil/surface/approve/mr-', '/surface/confirm/x', '/surface/confirm/mr-0123456789abcdef']) {
+  for (const href of ['https://evil/x', '/surface/model.json', '//evil/surface/approve/mr-', '/surface/confirm/x', '/surface/confirm/mr-0123456789abcdef',
+    '//evil/surface/confirm/vr-0123456789abcdef', '/x/surface/confirm/vr-0123456789abcdef', '/surface/confirm/vrx0123456789abcdef', '/surface/confirm/VR-0123456789abcdef']) {
     const { log } = await run(MR, reply(200, JSON.stringify({ ...B(MR), href })));
     assert.deepEqual(log, [['say', 'refused: unexpected target']], href);
   }
diff --git a/host/surfacepage/confirm_script_test.go b/host/surfacepage/confirm_script_test.go
index 9979251d..ffb8c3af 100644
--- a/host/surfacepage/confirm_script_test.go
+++ b/host/surfacepage/confirm_script_test.go
@@ -300,6 +300,7 @@ const sections = {
       ['409', answer(409, { v: 1, error: 'Conflict', reason: 'the request changed since the card was shown; reload it' }), 'Not confirmed (HTTP 409 Conflict): the request changed since the card was shown; reload it'],
       ['500', answer(500, { v: 1, error: 'Internal', reason: 'the confirmation did not complete; approve again' }), 'Not confirmed (HTTP 500 Internal): the confirmation did not complete; approve again'],
       ['200 not confirmed', answer(200, { v: 1, confirmed: false }), 'Not confirmed (HTTP 200 ): no reason given'],
+      ['500 that says confirmed', answer(500, { v: 1, confirmed: true, approvedBy: 'ssh:x', decision: 'd', next: 'n', error: 'Internal', reason: 'r' }), 'Not confirmed (HTTP 500 Internal): r'],
       ['lost', () => Promise.reject(new TypeError('Failed to fetch')), 'Confirmation outcome unknown: Failed to fetch. Reading the request again.']]) {
       const p = await open(card(c), daemon([c], post));
       p.$('sig').value = 'sig';
diff --git a/host/verbs/lookup_test.go b/host/verbs/lookup_test.go
index 23cedc9d..bdd98d1d 100644
--- a/host/verbs/lookup_test.go
+++ b/host/verbs/lookup_test.go
@@ -86,7 +86,8 @@ func TestApproveTarget(t *testing.T) {
 	ref := "sha256:" + strings.Repeat("ab", 32)
 	for args, want := range map[string]string{"mr-0123456789abcdef": "/surface/approve/mr-0123456789abcdef", ref: "/surface/decide/" + ref,
 		"mr-0123456789abcdef --deny": "", "--list": "", "": "", ref + "0": "", "sha256:" + strings.Repeat("AB", 32): "", "mr-0123456789abcdef x": "", "mr-0123456789ABCDEF": "",
-		"vr-0123456789abcdef": "/surface/confirm/vr-0123456789abcdef", "vr-0123": "", "vr-0123456789abcdef0": "", "vr-0123456789ABCDEF": "", "xvr-0123456789abcdef": ""} {
+		"vr-0123456789abcdef": "/surface/confirm/vr-0123456789abcdef", "vr-0123": "", "vr-0123456789abcdef0": "", "vr-0123456789ABCDEF": "", "xvr-0123456789abcdef": "",
+		"vr-0123456789abcdef x": "", "vr-0123456789abcdef --deny": "", "--deny vr-0123456789abcdef": ""} {
 		href, reason := ApproveTarget(strings.Fields(args))
 		if href != want || (href == "") == (reason == "") {
 			t.Errorf("ApproveTarget(%q) = %q, %q", args, href, reason)
````

