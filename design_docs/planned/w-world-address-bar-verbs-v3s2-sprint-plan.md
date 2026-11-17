# Sprint plan — w-world-address-bar-verbs-v3, sprint S2 "the card page, against fixtures" (iteration 401, row 243)

> **THE SPRINT — row 243 phase V3, sprint S2 (design §0, §1, §2.2, §2.5, §2.6 (the seam only), §2.9 P4, §3–§9, §11,
> §12, §14, §16), with S1b folded in:** the confirmation card's page and its pure data, tested against fixtures, with
> **no route registered, no store read or written, no credentials schema change, no daemon change, no `.ail`, the CLI's
> output unchanged and `surfaceCSP` unchanged**. `host/surfacepage` gains **`ConfirmationInput`**, the wire type
> **`Confirmation`** and **`ConfirmationCard(in, passkeys []CardPasskey)`** (design §2.5's data JSON: the stored bytes
> base64url, the digest, the raw challenge, the times, the state; **never a message hash, never an artifact**), and the
> embedded page **`confirm/{index.html,confirm.js}`** (`ConfirmPage`, `ConfirmScript`, embedded as `approve.go:13-18`
> embeds A3a's): a 10 s aborted GET, the digest check, one `CONFIRM_FIELDS` table, `textContent` rendering of §2.2's
> layout, the 86-byte `world-verb/v1` message and its hash built **in the page**, *Save the message* + the `ssh-keygen
> -Y sign -n world-verb@sunholo` line, and no controls in any state but `pending`. M0 takes the two V2b judge notes
> (VB8, VP12, tests only). **S1b** (the judges `JudgeVerbSSH`, `JudgeVerbAssertion`, `RecheckVerb`, VJ1) rides as the
> last sub-milestone: iteration 400's patch, which still applies cleanly to `07d7c90e`, plus three VJ1 test lines
> that kill the two survivors the mutant run found (S2-D1).
> - **Measured: 1,380 changed lines** — `git diff -M --numstat 07d7c90e HEAD -- cmd host docs` (insertions + deletions,
>   renames on; no rename occurs; the same 1,380 with no path filter, so nothing lands outside `cmd host docs`), after a
>   commit of each of the six scratch stages: **429 production** (`confirm.js` 167, `judge_verb.go` 129,
>   `confirm_model.go` 72, `index.html` 40, `confirm.go` 21), **951 test** (incl. the 5-line generated golden), **0
>   docs**. **Cap 1,400: margin 20.**
> - **S2 alone: 1,051** (M0 68 + M4 225 + M5a 418 + M5b 231 + M6 109) against the design's ≈ 825 raw: **1.27×**.
>   **S1b folded (S2-D1):** iteration 400's patch (sha256 prefix `775b927ad56aa4c0`) re-applied to `07d7c90e` cleanly
>   and re-measured at **326**: 1,051 + 326 = 1,377 ≤ 1,400, the directive's rule; with the three VJ1 lines that kill
>   R2's two S1b survivors it is **329**, total **1,380**. VJ1's 22 mutants run.
> - **Six sub-milestones, each green at its commit** (§1.2 C14): **M0** VB8+VP12 (68) → **M4** `ConfirmationCard`, CP0
>   (225) → **M5a** the page, the script, CP2's table test, the generated golden, the no-route proof (418) → **M5b**
>   `confirm.js` under node, CP1–CP6 (231) → **M6** XC1 live Chrome (109) → **S1b** the judges, VJ1 (329).
> - **Every line ships as a compiled patch** (Appendix P), extracted from this file and round-tripped onto a fresh
>   `git archive 07d7c90e` export, byte-identical to the stage tree (§1.2 C16). The executor applies, never retypes.
> - **Pre-existing test files changed: exactly one**, `host/verbs/why_admit_test.go` (+1 row, VB8). Every other test file
>   is new (§5, §9 criterion 11).

---

## 0. Binding constraints (each is restated where it governs)

- **D-WORLD-118 = A** (verified by the controller at iteration 400, `mission_ledger_provenance.sh --id D-WORLD-118` →
  `verdict=attended`, flip `72007b37`). S2 does not depend on it; nothing B-only (F6, `DraftOnly`) is built.
- **No route** (`daemon.go` untouched; proven by `TestConfirmRoutesNotRegistered`, a route-table test that also serves
  the three paths through the real Handler → 404, with a 200 control), **no store read or write** (the page's data is a
  fixture, served by `httptest` in XC1 and by a stand-in `fetch` under node), **no credentials schema change**, **no
  daemon process change**, **no `.ail`**, **the CLI's output unchanged** (no `cmd/ailang-worldd` production file
  changes), **`surfaceCSP` unchanged** (no `host/daemon` production file changes).
- **D-WORLD-109**: test keys and temp dirs only. XC1 runs Chrome headless with the passkey driver's temporary profile,
  against a loopback `httptest` server; S1b's tests use `approvaltest` keys and `t.TempDir()` credentials stores.
- Coding standards: the card's data is a pure function (`ConfirmationCard`, primitives in, a struct out); the page's
  only effect is one GET; S1b's judges read a read-only handle.
- Not touched: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/`, `cmd/world-publish/*`,
  `approve.js`, `decide.js`, `app.js`, `index.html` (the app's), `daemon.go`, `surface_page.go`.

---

## 1. Pre-flight

### 1.1 Baseline at `07d7c90e` (measured by the planner, 2026-11-17)

`tools/gates.sh <base-tree> base-gates base` (the §6 script in `base` mode; a `git archive 07d7c90e` export with a
`git init`; `GATES_SKIP_VERIFY_GO=1` — S1's record ran `verify_go` on this tree's `cmd/host`, `b0817474` = `07d7c90e`
for `cmd host`):

| leg | rc | time | note |
|---|---|---|---|
| `status` | 0 | 0 s | clean |
| `verify_ail` | 0 | 71 s | no `.ail` changes in S2; the final must match |
| `gofmt` / `build` / `vet` / `fence` | 0 / 0 / 0 / 0 | 0 / 5 / 3 / 25 s | `go test -run '^$' ./...` is the test-file compile fence |
| `pkgs` | 0 | 309 s | `host/verbs` 0.1 s, `host/surfacepage` 1.2 s, `host/daemon` 307.6 s, `cmd/ailang-worldd` 175.0 s |
| `race` | 0 | 20 s | |
| `node` | 0 | 1 s | 4 named `--- PASS` (the base's node harnesses), `WORLD_EXEC_NODE` set |
| `inv` | 0 | 19 s | **37 named `--- PASS`**, 0 SKIP |
| `chrome` | 0 | 45 s | the 14 passkey/approve/decide/Needs-you/verb-bar Chrome tests + 2 `host/daemon` drills `--- PASS` |
| **`GATES_DONE 499s FAILS: none`** | | | |

### 1.2 Codebase claims checked (at `07d7c90e` unless "work" — the planner's final stage tree, tag `S1b`)

| # | claim | command | observed |
|---|---|---|---|
| C1 | **Nothing the design cites for S2 moved** between the design's base and HEAD | `git diff --stat 994ba123 07d7c90e -- host/daemon/surface_approve.go host/surfacepage/ host/daemon/surface_page.go host/daemon/surface_decide_page_test.go host/verbs/flags.go host/daemon/daemon.go cmd/ailang-worldd/passkey_page_driver_test.go cmd/ailang-worldd/testdata/passkey_cdp.mjs`; control `git diff --stat 994ba123 07d7c90e -- cmd host \| tail -1` | empty; control `12 files changed, 1280 insertions(+), 20 deletions(-)` (S1) |
| C2 | the cites, re-measured | `grep -n` each | `surface_approve.go:76 func (d *Daemon) approvalRoute`, `:92 handleSurfaceRequests`, `:146 handleSurfaceRequestCard` (`:178` calls `surfacepage.RequestCard`); `approve.go:13 //go:embed approve/index.html`, `:18 approve/approve.js`; `approve_model.go:132 type CardPasskey struct`, `:138 CardPasskeys`, `:145 CardSSH`, `:166 func RequestCard`; `surface_page.go:19 const surfaceCSP = "default-src 'none'; … connect-src 'self'; …"`; `decide.js:10 const NET_MS = 10000`, `:60 setTimeout(() => c.abort(), NET_MS)` (`approve.js:9`, `:78`); `surface_decide_page_test.go:38 TestDecidePageFieldTable` (A4's PD1), `:84 TestDecidePageAssets`; `app.js:1015 … + ' · episode report <episode> · why <target>';`; `app.js:1303 const KEYROWS = [… ['Enter', 'in the bar: a place, episode report E, why T, approve ID', 'the same'] …]`; `index.html:67 … placeholder="row-189@6 · shared@0ffbb9c · back · episode report &lt;episode&gt; · why &lt;target&gt;">`; `flags.go:137 cli := Quote(append([]string{Binary, "why"}, args...))` — **all unmoved** |
| C3 | the V2b notes are unpinned at HEAD (W27, W28) | `grep -rlE "why &lt;target&gt;\|why <target>'\|why T, approve ID" host/surfacepage/*_test.go host/surfacepage/testdata cmd/ailang-worldd/*_test.go host/daemon/*_test.go \| wc -l`; control `grep -rl "did not answer in time" host/surfacepage/*_test.go`; `TestAdmitWhy`'s rows `why_admit_test.go:83-88` | `0`; control `host/surfacepage/app_verb_test.go`; no row's argument needs quoting |
| C4 | nothing named confirm exists yet | `grep -rln "surface/confirm\|ConfirmationCard\|CONFIRM_FIELDS" cmd host \| wc -l`; control `grep -rln "surface/decide" cmd host \| wc -l` | `0`; control `25` |
| C5 | S1's API, which S2 consumes, is at HEAD | `grep -n` | `verb.go:22 VerbRequestSchema`, `:32 VerbFinish`, `:35 VerbPrefix = "world-verb/v1"`, `:38 VerbNamespace = "world-verb@sunholo"`, `:41 VerbMessageLen = len(VerbPrefix) + 1 + 32 + 32 + 8`, `:51 type VerbRequest struct` (ten json tags `schema verb argv episode prompt scope title body committer protected`), `:69 Validate`, `:120 CanonicalVerb`, `:152 VerbMessage`, `:166 VerbMessageHash`; `finish_prompt.go:14 DisplayCommitter`, `:24 FinishPromptText`; `approve.go:689 type FinishApprovalScope`, `:698 FormatFinishApprovalScope`; `split.go:100 func Quote`; `episode_finish_operator.go:153 func finishPrompt` |
| C6 | `Admit`'s expiry rule (the card's `expired` mirrors it) | `grep -n "if in.Now >= in.ExpiresAt" host/approval/admit.go` | `51` |
| C7 | the closest landed analogues and how they are tested | `grep -rln "DecideScript\|DecidePage\|DECISION_FIELDS\|decide.js" --include="*.go" cmd host`; `ls design_docs/planned \| grep a4-sprint4` | A4's decide page is tested by **live Chrome only** (`decide_page_{chrome,closed,sshsig}_test.go`, rig `decide_page_rig_test.go`, all `startPasskeyBrowser`-gated) plus PD1/PD2 in `host/daemon`; **no node test reads `decide.js`**; plan `w-world-approvals-by-key-a4-sprint4-plan.md`. The node pattern is `app.js`'s (`app_test.go:124`, `app_decide_test.go:17`: `WORLD_EXEC_NODE`, else `exec.LookPath("node")`, skip only without node), made non-skippable in CI by `ci.yml:250` (`WORLD_EXEC_NODE=$(command -v node)`) |
| C8 | how live Chrome runs and skips | `grep -n "WORLD_REQUIRE_PASSKEY_BROWSER\|func startPasskeyBrowser" cmd/ailang-worldd/passkey_page_driver_test.go`; V2b's `gates.sh` `chrome` leg | `:63 startPasskeyBrowser`: needs `WORLD_CHROME` and node; `:69` `WORLD_REQUIRE_PASSKEY_BROWSER=1` turns the skip into a failure; CI sets neither (Chrome tests skip there); V2b's leg runs them with both set and requires `--- PASS` per name. XC1 is gated identically |
| C9 | the daemon's script headers (XC1 serves the card under them, not under a copied literal) | `grep -n "func (d \*Daemon) handleSurfacePasskeyScript" -A2 host/daemon/surface_page.go`; `grep -n '"GET /surface/passkey.js"' host/daemon/daemon.go` | `:181-182 … setSurfaceHeaders(w)` (= `surfaceCSP` + the workbench headers, `:21-24`); route `daemon.go:1134`, not behind `approvalRoute` |
| C10 | `host/verbs` is pure, so `host/surfacepage`'s tests may import it and `host/broker` without a cycle | `go list -f '{{.Imports}}' ./host/verbs`; `go list -deps ./host/broker ./host/verbs \| grep -c surfacepage` | `[errors flag fmt io regexp strings time unicode/utf8]`; `0` |
| C11 | S1b's patch (iteration 400, `775b927ad56aa4c0`) applies cleanly to HEAD and to S2's tip | `patch -p1 -E --dry-run < S1b-M2.patch` on the base export and on `M6` | `patching file 'host/authority/judge_verb.go'` / `'host/authority/judge_verb_test.go'` both times (`logs/s1b_dryrun_base.txt`); applied and committed at `M6` it measured **326** and regenerated **byte-identical** (`cmp`); R2's mutants then found two VJ1 survivors (§4: `JV-ADMIT-BYTES`, `JV-RCK-NOW`), killed by three new lines in `judge_verb_test.go` — Appendix P's `S1b` is that patch (329 lines changed) |
| C12 | node | `/opt/homebrew/bin/node --version` | `v26.13.0` |
| C13 | the new tests, work | `go test -count=1 -v ./host/surfacepage/ -run '^(TestAddressHintsArePinned\|TestConfirm.*)$'` (and the others by name) | all `--- PASS`, none SKIP (`WORLD_EXEC_NODE` set) |
| C14 | **every sub-milestone is green when committed** | `tools/accept_all.sh M0 M4 M5a M5b M6 S1b` (an APFS clone of `work` at each tag running §3's `accept.sh`) | first run (the first shipped stages) and **re-run on the final stages** (14:36Z, after R2's test strengthening): both times `M0 rc=0 ACCEPT_OK`, `M4 rc=0 ACCEPT_OK`, `M5a rc=0 ACCEPT_OK`, `M5b rc=0 ACCEPT_OK`, `M6 rc=0 ACCEPT_OK` (XC1 `--- PASS`), `S1b rc=0 ACCEPT_OK` (`logs/accept_all.txt`) |
| C15 | the snapshot block (§7) on every sub-milestone, incl. M0–M5b (no `cmd/`) | `tools/snap_rehearse.sh` (a local clone at `07d7c90e`, the patches in order, §7's block verbatim) | `M0 2/0/2`, `M4 4/0/4`, `M5a 10/0/10`, `M5b 11/0/11`, `M6 12/0/12 (cmd host)`, `S1b 14/0/14`, each `MANIFEST_OK` (STATUS/DELETED/MANIFEST lines) |
| C16 | **the round trip** | `tools/roundtrip.sh <this plan>` | `CMP_EQUAL` for all six (each block extracted from this file with §2's awk equals the generated patch byte for byte); each applied with `/usr/bin/patch -p1 -E` after a `--dry-run`, no `.orig`/`.rej`, committed and measured: `NUMSTAT M0 68`, `M4 293`, `M5a 711`, `M5b 942`, `M6 1051`, `S1b 1380`; then `diff -r -x .git` against the stage tree → **`ROUNDTRIP_IDENTICAL_TO_STAGE_TREE`** (`logs/roundtrip_last.txt`; re-run on the final plan file, below in §1.2's note) |
| C17 | **the gates on the final tree** | `tools/gates.sh <final-tree> final-gates final` | `tools/gates.sh <final-tree> final-gates final 80dd837e` (`final-tree` = the round-trip tree, every stage committed, byte-identical to the stage tree; run beside nothing else after the mutant runs; load 6–29 from other sessions): `status` 0; **`verify_ail` rc 0** (`✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` — the base's line, byte for byte); `gofmt`, `build`, `vet`, `fence` rc 0; **`pkgs` rc 0** (`host/verbs`, `host/surfacepage` 1.3 s, `host/daemon` 300.1 s, `cmd/ailang-worldd` 152.3 s); **`race` rc 0**; **`node` rc 0: 9 named `--- PASS`** (the base's 4 + `TestAddressHintsArePinned` + the four CP harnesses), none SKIP; **`inv` rc 0: 43 named `--- PASS`** (the base's 37 + `TestConfirmRoutesNotRegistered` + the five new `host/surfacepage` tests), 0 SKIP, 0 FAIL, none missing; **`chrome` rc 0: 17 `--- PASS`** (V2b's 14 + the two `host/daemon` drills + **XC1 `TestConfirmPageShowsTheFinishPromptInChrome (1.66s)`**); **`verify_go` rc 1** — 63 packages `ok`, the one red **`host/verifygate` `panic: test timed out after 10m0s`** (running `TestExecGateArmsAreLoadBearing/sketch_tests`, `TestSeToolsGateArmsAreLoadBearing/plan-bytes`; a package S2 does not touch — no `.ail`, no doc, no gate input changes — and red in iteration 400's controller gate run too (`world-iter400/gates/verify_go.txt`: `--- FAIL: TestMissionConfigDefaultFlowIntegration`, `FAIL host/verifygate 458.268s`, `GATES_DONE 965s FAILS: verify_go`), on the directive's load-sensitive list); re-run alone on the final tree (`logs/verifygate_alone.txt`, load ≈ 20): `FAIL host/verifygate 514.9s` with only **`TestMissionConfigDefaultFlowIntegration`** red (its `valid_ledger`/`invalid_ledger` 10 s arms; the directive's list), and that test alone is **red at the base too** (`--- FAIL … (19.17s)` on `base-tree`) and **green on the final tree** the next minute (`ok host/verifygate 19.3s`, `logs/mission_config_flow_base_vs_final.txt`) — load-sensitive, not S2's; **`GATES_DONE 1167s FAILS: verify_go`**; `ALL_DONE.80dd837e…` written (`plan-scratch/final-gates/`) |

### 1.3 Planner decisions where the design or the directive leaves a choice

- **S2-D1 — S1b is folded, as the last sub-milestone.** S2 compiled to **1,051**; 1,051 + 326 = 1,377 ≤ 1,400 (the
  directive's rule; S1's Appendix C's stricter "≤ 1,074" also holds). Re-applied to `07d7c90e` it applies cleanly and
  measures 326; VJ1's mutants are enumerated per S1 Appendix C (step labels, lookup branches, `RecheckVerb`'s kind
  detection and its use of `in.ExpiresAt`, plus its `in.Digest`/`in.Now`, the approvers' fields and `admit()`) and run
  (§4). Two survived R2 — `JV-ADMIT-BYTES` (`admit()` passing `in.Digest` as the bytes' digest: no judge test had a row
  whose bytes do not hash to its digest) and `JV-RCK-NOW` (`RecheckVerb` checking P5 at `in.CreatedAt`: every approver
  was enrolled before both) — and three test lines kill them (`VJ-law-bytes-digest`; a recheck control with the
  approver enrolled at 1300, between the row's creation 1000 and the caller's now 1605): S1b = **329**, total **1,380**; `TestJudgeVerbSSH|TestJudgeVerbAssertion|
  TestRecheckVerb` join `inv`. Last, not first, so that a red S1b at commit cannot block the card: S1b touches only
  `host/authority`, which no S2 file imports. **No cut was needed**; XC1 stays in S2.
- **S2-D2 — the wire type is `Confirmation`.** Design §2.5 names the function `ConfirmationCard`; Go cannot give a type
  and a function one name, so the JSON's type is `Confirmation` (A3a's analogue is `Card` + `RequestCard`). Design §7's
  `Confirmed` is S4's.
- **S2-D3 — the state is derived from primitives.** `ConfirmationInput` carries `Outcome` (`""`, `approved`,
  `superseded`: S3's `verb_requests.outcome` column, design §2.4) and `Now`; `confirmationState` is the outcome if any,
  else `expired` when `Now >= Expires` (`Admit`'s rule, `admit.go:51`), else `pending`. `approvedBy`/`decision` are
  emitted only for `approved`. An unknown non-empty outcome passes through and the page shows it with no controls.
- **S2-D4 — `CP1–CP6` run under node, XC1 in Chrome.** A4's card is Chrome-only (C7), so CI never runs its behaviour.
  S2 runs `confirm.js` itself — the embedded bytes, unmodified — under node in a `vm` context with a stand-in DOM built
  from `index.html`'s own ids and `hidden` attributes, a stand-in `fetch` and controllable timers, against the golden
  (four tests, one per CP group). They are CI-run and non-skippable there (`WORLD_EXEC_NODE`), and they see every sink
  (`innerHTML` and its relatives are recorded by the stand-in and asserted empty). XC1 proves the real browser: CSP,
  layout, `textContent` in a `<pre>`, the prompt equal to `FinishPromptText`.
- **S2-D5 — the golden is generated, never typed.** `TestConfirmCardGolden` builds the cards from
  `approval.CanonicalVerb` over requests whose `Prompt` is `broker.FinishPromptText` and `Scope`
  `broker.FormatFinishApprovalScope`, with Go's `approval.VerbMessage`/`VerbMessageHash` and `verbs.Quote(argv)`; it
  compares byte for byte, and **`WORLD_UPDATE_CONFIRM_GOLDEN=1 go test ./host/surfacepage/ -run
  '^TestConfirmCardGolden$'`** rewrites it. Five JSON lines: `draft`, `hostile` (pending), `draft-approved`,
  `draft-superseded`, `draft-expired` (each from `ConfirmationCard`).
- **S2-D6 — XC1 serves the card under the daemon's own headers.** No route exists, so a loopback `httptest` handler
  serves `ConfirmPage`, `ConfirmScript` and the fixture; it copies the headers the real daemon Handler gives `GET
  /surface/passkey.js` (`setSurfaceHeaders` → `surfaceCSP`), asserting they carry `script-src 'self'` and `connect-src
  'self'`. XC1's expected prompt is `broker.FinishPromptText(scope, title, body, protected, committer)`, with a control
  that the CLI's own `finishPrompt` prints exactly those bytes before its question.
- **S2-D7 — "no route" is proven by a test, not prose.** `TestConfirmRoutesNotRegistered` (host/daemon) cuts the
  Handler's body from `daemon.go` (as `TestApprovalRoutesRegistered` does) and finds no `/surface/confirm`, and serves
  `GET /surface/confirm/vr-…`, `/surface/confirm.js`, `/surface/confirmations/vr-…` through the real Handler → 404, with
  the decide page → 200 as its control. **S3 deletes it** when it registers the routes (Appendix C).
- **S2-D8 — VP12 is a string pin, and says so.** The two `app.js` strings live in the DOM half (after `:572 if (typeof
  document === 'undefined') return;`), which node cannot execute; the node harness cuts the placeholder expression and
  the `KEYROWS` literal from `app.js`'s own source and evaluates them in a `vm` (each placeholder branch), and Go reads
  `index.html`'s placeholder from the embedded bytes. Their mutants are class **S** in §4 (killed only by that pin, as
  the design's "Kills: each string edited" asks); no behavioural test of the app shows these hints (C3).
- **S2-D9 — the S4 seam.** `confirm.js` leaves `#controls` as the one place the submit half goes (S4: the passkey
  button, the SSHSIG paste box, A2's 240 s abort, the re-read), and the GET as one `net()` with the 10 s abort; the
  states are one `CLOSED` table. Nothing of S4 is stubbed.
- **S2-D10 — the CP3 NUL row (S1 Appendix C, P6).** The `hostile` fixture's title is `<img src=x onerror=alert(1)>
  a\x00b & 'q'` and its body carries `NUL \x00 here &amp; after — ünï` (multibyte, so the size line's byte count is
  tested); both pass `Validate`; the page shows the text after each NUL in the title, the prompt and the body.

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

| # | premise | measured | routing |
|---|---|---|---|
| F1 | design §6 CP2: "a Go test reads both, as A4's PD1" | PD1 (`surface_decide_page_test.go:38`) reads the **served** wire through a route and a store (C2, C7); S2 has neither | CP2's Go half reads `approval.VerbRequest`'s json tags **by reflection** in `host/surfacepage` (as S1 Appendix C foresaw); the JS half is the node `closed` section |
| F2 | design §6 "Chrome: the A3/A4 driver" for the card's behaviour | A4's page has no CI-run behavioural test (C7, C8) | S2-D4: node for CP1–CP6, Chrome for XC1 |
| F3 | design §2.5 function `ConfirmationCard` + design §7's type list | one Go identifier cannot be both | S2-D2: type `Confirmation` |
| F4 | design §9(b)/§6 VP12 "node for the two `app.js` strings" | the strings are unreachable by node's execution of `app.js` (C3, S2-D8) | evaluated from the source in a `vm`; class S, said in §4 |
| F5 | design §2.5 "unchanged `surfaceCSP`" for the page | no S2 route applies any header | XC1 serves under the daemon's own headers, copied at run time (S2-D6) |
| F6 | design §7's file list | S2 also adds `host/surfacepage/app_hints_test.go` (VP12), `host/surfacepage/confirm_script_test.go` (node), `host/surfacepage/testdata/confirm-card.golden.jsonl`, `host/daemon/surface_confirm_unrouted_test.go` (S2-D7), `cmd/ailang-worldd/confirm_page_chrome_test.go` (XC1), and S1b's `host/authority/judge_verb{,_test}.go` (§7's sprint-1 list) | **files outside design §7: the five test files above** — all new tests; no production file outside §7 |
| F7 | S1 Appendix C: "S2's M0 (VB8, VP12) is unaffected by S1" | holds (C1, C2) | — |

### 1.5 The size

| stage | production | test | changed lines | cumulative (numstat after commit) |
|---|---|---|---|---|
| M0 | 0 | 68 | 68 | **68** |
| M4 | 72 | 153 | 225 | **293** |
| M5a | 228 (`confirm.js` 167, `index.html` 40, `confirm.go` 21) | 190 (`confirm_page_test.go` 152, golden 5, `surface_confirm_unrouted_test.go` 33) | 418 | **711** |
| M5b | 0 | 231 | 231 | **942** |
| M6 | 0 | 109 | 109 | **1,051** |
| S1b | 129 | 200 | 329 | **1,380** |

Design ≈ 825 raw for S2 → compiled 1,051 (**1.27×**: M0 68 vs ≈ 45; M4 225 vs ≈ 140; M5 649 vs ≈ 530; M6 109 vs ≈ 110).
Below the 1.7× law sprints measured: a page sprint's tests are harnesses, not tables. Cap 1,400, margin 20 — any growth
on the executor's side is a stop (§7).

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of four backticks,
in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; a new file's patch has `--- /dev/null` and an
`@@ -0,0 …` hunk. S2 deletes no file and edits one (`why_admit_test.go`); `-E` is kept. Apply it **only** this way, in
one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF2'
set -euo pipefail
M=M0   # <- the sub-milestone, exactly as in the Appendix P heading, in order: M0 M4 M5a M5b M6 S1b
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter401-row243-v3s2
PLAN="$W/design_docs/planned/w-world-address-bar-verbs-v3s2-sprint-plan.md"
T="$HOME/.ailang/state/world-iter401/exec-patches"; mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
/usr/bin/patch -p1 -E --dry-run < "$T/$M.patch"
/usr/bin/patch -p1 -E < "$T/$M.patch"
if find cmd host docs \( -name '*.orig' -o -name '*.rej' \) | grep .; then echo "STOP: patch left .orig/.rej"; exit 1; fi
git status --short --untracked-files=all
EOF2
```

- If `--dry-run` fails, **stop and report** (no hand edits, no fuzz options). It can only fail if an earlier
  sub-milestone was skipped, reordered or changed.
- Then run that sub-milestone's **Acceptance** (§3): extract `accept.sh` from this plan (`awk '$0 == "````bash
  accept.sh" {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/accept.sh"`) and run `bash "$T/accept.sh" "$W" <M>
  2>&1 | tee "$HOME/.ailang/state/world-iter401/exec-logs/accept_<M>.txt"` in one bash shell (never only through
  `tail`). **Never `nohup`, never background a test leg; pass a 600000 ms timeout to every long Bash call.**
- Then snapshot (§7). The controller commits and compares the numstat with §3's expected value.

---

## 3. Milestones

**Acceptance.** Every sub-milestone's acceptance is the arm of this script named for it, preceded by the same
`[base]` checks (gofmt, build, vet of the five touched packages, the test-file compile fence). `v` requires every named
test to `--- PASS` and none to SKIP or FAIL; `x` is live Chrome (temp profile, the passkey driver), required to PASS.
Measured (C14): each arm ≈ 1–4 minutes, one 600 s call each.

````bash accept.sh
#!/bin/bash
# accept.sh <tree> <M>: exactly plan §3's acceptance for sub-milestone M (one bash shell; rc 0 = ACCEPT_OK).
T=$1; M=$2
cd "$T" || exit 2
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
unset WORLD_UPDATE_CONFIRM_GOLDEN WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER
set -o pipefail
chk() { tee /dev/stderr | { out=$(cat); echo "$out" | grep -qE -- '--- (FAIL|SKIP)|^FAIL' && exit 1; echo "$out" | grep -q -- '--- PASS' || exit 1; }; }
v() { go test -count=1 -v "$@" 2>&1 | chk; }   # -v: every named test PASS, none SKIP/FAIL
p() { go test -count=1 "$@"; }
x() { WORLD_CHROME='/Applications/Google Chrome.app/Contents/MacOS/Google Chrome' WORLD_REQUIRE_PASSKEY_BROWSER=1 go test -count=1 -v ./cmd/ailang-worldd/ -run "$1" 2>&1 | chk; }   # live Chrome, temp profile
step() { echo "== $*"; "$@" || { echo "ACCEPT_RED $M: $*"; exit 1; }; }
# [base]
step test -z "$(gofmt -l cmd host)"
step go build ./...
step go vet ./host/verbs/ ./host/surfacepage/ ./host/daemon/ ./host/authority/ ./cmd/ailang-worldd/
if go test -run '^$' ./... 2>&1 | grep -E '^(FAIL|--- FAIL)'; then echo "ACCEPT_RED $M: FENCE_RED"; exit 1; fi
M0T='TestAddressHintsArePinned'
CP0='TestConfirmationCard|TestConfirmationCardWire'
PAGE='TestConfirmPageFieldTable|TestConfirmCardGolden|TestConfirmPageAssets'
NODE='TestConfirmPageFailsClosed|TestConfirmPageRendersText|TestConfirmPageMessage|TestConfirmPageStates'
ROUTES='TestConfirmRoutesNotRegistered|TestApprovalRoutesRegistered|TestDecisionRoutesRegistered|TestDecidePageAssets|TestFrozenV1RouteTableMatchesMux|TestDaemonDependencyAllowlist|TestSurfaceAppRoutes'
case $M in
M0)  step v ./host/verbs/ -run '^(TestAdmitWhy|TestSplit|TestQuoteRoundTrip|TestLookup)$'
     step v ./host/surfacepage/ -run "^($M0T|TestSurfaceAppPure|TestSurfaceAppVerbs)\$"
     step p ./host/verbs/ ./host/surfacepage/ ;;
M4)  step v ./host/surfacepage/ -run "^($M0T|$CP0)\$"
     step p ./host/surfacepage/ ;;
M5a) step v ./host/surfacepage/ -run "^($M0T|$CP0|$PAGE)\$"
     step v ./host/daemon/ -run "^($ROUTES)\$"
     step p ./host/surfacepage/ ;;
M5b) step v ./host/surfacepage/ -run "^($M0T|$CP0|$PAGE|$NODE)\$"
     step p ./host/surfacepage/ ;;
M6)  step x "^(TestConfirmPageShowsTheFinishPromptInChrome|TestDecidePageFailsClosed)\$"
     step v ./host/surfacepage/ -run "^($M0T|$CP0|$PAGE|$NODE)\$"
     step v ./cmd/ailang-worldd/ -run '^(TestEveryVerbHasHelp|TestFinishPromptIsTheGolden|TestVerbTableMatchesTheCLISwitch)$' ;;
S1b) step v ./host/authority/ -run '^(TestJudgeVerbSSH|TestJudgeVerbAssertion|TestRecheckVerb|TestJudgeDecisionSSH|TestJudgeDecisionAssertion)$'
     step v ./host/surfacepage/ -run "^($M0T|$CP0|$PAGE|$NODE)\$"
     step p ./host/authority/ ./host/approval/... ;;
*) echo "unknown $M"; exit 2 ;;
esac
echo "ACCEPT_OK $M"
````

### M0 — the V2b pins (design §9, §6 VB8, VP12) · 0 prod + 68 test · Run 1

- `host/verbs/why_admit_test.go` (**edited, +1 row**): `{"--result 'a b'", "", 0, false, "--result is not run from the
  app: …" + at + "--result 'a b'"}` — the refusal's CLI line is `Quote`'s form; a plain join (`flags.go:137` →
  `strings.Join`) gives `--result a b` (VB8-JOIN, §4: killed by `TestAdmitWhy`).
- `host/surfacepage/app_hints_test.go` (NEW) **`TestAddressHintsArePinned`** (VP12): Go asserts `index.html`'s
  placeholder once in `appIndex` and in `AppShell`; node evaluates `app.js`'s placeholder expression in its three
  branches (`st.d` ≠ 0 → `''`; no groups; two groups) and its `KEYROWS` literal (13 rows; row 11 is the bar's).
- **Expected numstat after commit: 68.**

### M4 — `ConfirmationCard` (design §2.5, §6 CP0) · 72 prod + 153 test · Run 1

- `host/surfacepage/confirm_model.go` (NEW): `ConfirmationPending|Approved|Superseded|Expired`, `ConfirmationInput`
  (`ID, Request, Digest, Challenge, Created, Expires, Now, Outcome, ApprovedBy, Decision` — primitives; no artifact, no
  hash), `Confirmation` (the 13 keys of design §2.5, in its order), `ConfirmationCard` (passkeys grouped by reusing
  `RequestCard`'s grouping, `approve_model.go:166`; `ssh.namespace = approval.VerbNamespace`), `confirmationState`.
- `host/surfacepage/confirm_model_test.go` (NEW): `confirmFixtures` (two requests built by `approval.CanonicalVerb`,
  their `Prompt` by `broker.FinishPromptText`, their `Scope` by `broker.FormatFinishApprovalScope`, both `Validate`d;
  `hostile` carries markup, quotes, NUL and multibyte text) — **`TestConfirmationCard`** (bytes, digest, challenge,
  times, the base64url form (`-__-`, not base64's `+//+`), passkey grouping, `[]` for none, the namespace literal, and the six state rows incl. `Now == Expires` →
  `expired`, approver/decision only when approved) and **`TestConfirmationCardWire`** (`ConfirmationInput`'s exact field
  list; the JSON's exact 13 keys in order; neither the 86-byte message nor its hash appears in hex, base64url or base64).
- **Expected numstat after commit: 293.**

### M5a — the page, the script, the table, the golden, no route (design §2.5) · 228 prod + 190 test · Run 1

- `host/surfacepage/confirm.go` (NEW): `ConfirmPage`, `ConfirmScript` (`//go:embed confirm/index.html`,
  `confirm/confirm.js`, as `approve.go:13-18`).
- `host/surfacepage/confirm/index.html` (NEW): `#what`; `#shown` (hidden) holding the command `<pre id="argv">`, the
  terminal's prompt `<pre id="prompt">`, "Open this pull request?" with `#status` and `#controls` (hidden: `#save`,
  `#sshline` — S4 adds its controls here), and *What your key signs* `<dl id="signed">`; one `<script
  src="/surface/confirm.js">`.
- `host/surfacepage/confirm/confirm.js` (NEW): `NET_MS = 10000` with an `AbortController` (A3 §2.7's rule); the id from
  the path (`^/surface/confirm/(vr-[0-9a-f]{16})$`, else refused before any fetch); GET `/surface/confirmations/<id>`
  (`no-store`, `omit`); refusals, in order: fetch error/timeout → `failed`; non-200 or non-JSON → `failed`; request not
  base64url → `refused`; another id, digest or challenge not 64 lowercase hex, expiry not a safe integer ≥ 0, state not a
  string → `refused`; `SHA-256(bytes) ≠ digest` → `refused`; not UTF-8 JSON (`fatal`) → `refused`; a field
  `CONFIRM_FIELDS` cannot show (unlisted, mistyped, missing, not an object) → `refused` "…confirm at the terminal". Then
  the 86-byte message (`world-verb/v1` ‖ 0x00 ‖ digest@14 ‖ challenge@46 ‖ u64be(expiry)@78) and its SHA-256; render
  (argv as `verbs.Quote`'s line; the prompt verbatim; digest, challenge, expiry, message hash; every `CONFIRM_FIELDS`
  row, lists one per line, the body with its byte count in a `<details>`) — **`textContent` only**; then `pending` →
  Save the message (`blob:` of the 86 bytes, `world-verb-<id>.msg`) + `ssh-keygen -Y sign -n world-verb@sunholo -f
  <key> world-verb-<id>.msg`, state `ready`; any other state → `closed`, its text, no controls.
- `host/surfacepage/confirm_page_test.go` (NEW): **`TestConfirmPageFieldTable`** (CP2: `CONFIRM_FIELDS` == the ten json
  tags of `approval.VerbRequest` in order, with `string`/`strings`, each labelled), **`TestConfirmCardGolden`**
  (S2-D5), **`TestConfirmPageAssets`** (embedded bytes = the files; one script; no inline handler, form, image, link,
  frame or remote URL; every `$('id')` the script reads exists; banned sinks — labelled instrument health).
- `host/surfacepage/testdata/confirm-card.golden.jsonl` (NEW, **generated** by `TestConfirmCardGolden`).
- `host/daemon/surface_confirm_unrouted_test.go` (NEW): **`TestConfirmRoutesNotRegistered`** (S2-D7).
- **Expected numstat after commit: 711.**

### M5b — `confirm.js` under node (design §2.9 P4, §6 CP1–CP6) · 0 prod + 231 test · Run 1

- `host/surfacepage/confirm_script_test.go` (NEW): `runConfirmHarness` (as `app_decide_test.go`: `WORLD_EXEC_NODE`, else
  node on `PATH`, skip only without node) runs one section of `confirmHarness` over the embedded `confirm.js` and
  `index.html` and the golden — **`TestConfirmPageFailsClosed`** (CP1 + CP2's JS half: the control is `ready` with one
  GET of the right URL and options and one cleared 10 000 ms timer; then network error, the 10 s timeout fired, 404 with
  its reason, 503, non-JSON, another id, digest case, short challenge, negative and fractional expiry, a non-string
  state, `*`-prefixed and **padded** base64, a wrong digest, another card's bytes, a lone `0xff` inside the title,
  unlisted / missing / argv-a-string / protected-a-number / title-a-number / body-null / an array, a bad path with no
  GET — each `refused`/`failed` with its reason, nothing shown, no controls, no message), **`TestConfirmPageRendersText`**
  (CP3: every field, the argv line = Go's `verbs.Quote`, the prompt, the body and its byte count, digest, challenge,
  expiry as text; no sink written; the NUL rows), **`TestConfirmPageMessage`** (CP4, CP6: the shown hash = Go's
  `VerbMessageHash`; the blob = Go's 86-byte `VerbMessage`, `application/octet-stream`, `world-verb-<id>.msg`; the
  signing line), **`TestConfirmPageStates`** (CP5: approved / superseded / expired / an unknown state → `closed`, the
  state's text, the fields shown, no controls, no blob; control `ready`).
- **Expected numstat after commit: 942.**

### M6 — XC1 live Chrome (design §6 XC1) · 0 prod + 109 test · Run 1

- `cmd/ailang-worldd/confirm_page_chrome_test.go` (NEW): **`TestConfirmPageShowsTheFinishPromptInChrome`** (S2-D6),
  gated by `startPasskeyBrowser` exactly as XB1/XD1 (`WORLD_CHROME`, `WORLD_REQUIRE_PASSKEY_BROWSER=1`). It asserts state
  `ready`, `#prompt` = `broker.FinishPromptText(…)` of the fixture, the title, the argv line, `#m-hash` = Go's hash, the
  save name, the signing line and **zero CSP violations** (`window.__cspv`).
- **Expected numstat after commit: 1,051.**

### S1b — the judges (design §2.3, §6 VJ1; S1 Appendix C) · 129 prod + 200 test · Run 1

- `host/authority/judge_verb.go`, `judge_verb_test.go` (NEW): iteration 400's `S1b-M2.patch` (C11) plus three test
  lines (S2-D1: `badBytes` + the `VJ-law-bytes-digest` row in `TestJudgeVerbSSH`; the enrolled-at-1300 control in
  `TestRecheckVerb`).
- **Expected numstat after commit: 1,380.**

### 3.1 Every input and failure class S2's code reaches, traced (design §2.9 P4; directive retro 4)

| design §2.9 P4 row / input | branch (`confirm.js` / Go) | test (section) |
|---|---|---|
| data fetch: network error | `net` throws → `refuse('failed', …could not be read (<message>))` | `TestConfirmPageFailsClosed` (`network`) |
| data fetch: timeout (10 s) | `setTimeout(() => c.abort(), NET_MS)` → `AbortError` → "no answer within 10 s" | `…FailsClosed` (`timeout`: the 10 000 ms timer fired) |
| data fetch: non-2xx | `res.status !== 200` → `failed`, "HTTP <n>: <reason>" | `…FailsClosed` (`404`, `503`) |
| data fetch: a body that is not JSON | `!res.body` → `failed` | `…FailsClosed` (`not JSON`) |
| `SHA-256(bytes) ≠ digest` | `hex(digest) !== s.digest` → `refused` | `…FailsClosed` (`digest`, `swapped bytes`) |
| the answer malformed (id, digest/challenge grammar, expiry, state) | the shape check | `…FailsClosed` (six rows) |
| bytes not base64url / not UTF-8 JSON | `unb64u` / `TextDecoder({fatal})` | `…FailsClosed` (`base64url`, `padded base64`, `a lone 0xff`) |
| unlisted / missing / mistyped field | `unshown` over `CONFIRM_FIELDS` | `…FailsClosed` (seven rows); `TestConfirmPageFieldTable` (the table = the Go type) |
| hostile title, body, prompt, argv (incl. NUL, P6) | `textContent` sinks; `quote` | `TestConfirmPageRendersText`; XC1 (`Fix the <bar>` in Chrome) |
| state approved / superseded / expired (and unknown) | `s.state !== 'pending'` → `closed` | `TestConfirmPageStates`; the state itself: `TestConfirmationCard` |
| the message and its hash | `message()`; `sha256` | `TestConfirmPageMessage`; XC1 |
| not a card's path | the id regex → `refused`, no GET | `…FailsClosed` (`path`) |
| passkey ceremony aborted / 240 s; POST refusal | **S4** (no submit half in S2) | — (CP9, S4) |
| the card data: bytes, raw challenge, no hash, no artifact | `ConfirmationCard` | `TestConfirmationCard`, `TestConfirmationCardWire` |
| no route serves any of it | `daemon.go` unchanged | `TestConfirmRoutesNotRegistered` |
| VJ1 (S1b): each `JudgeError.Step` of the three judges | `judge_verb.go` | `TestJudgeVerbSSH`, `TestJudgeVerbAssertion`, `TestRecheckVerb` |

---

## 4. Mutation table

**Procedure** (`tools/mutate.py <run> <workers> [ids…]`, rows in `tools/mutants.py`, enumerated from the diff
`BASE..S1b`, not from the test plan): every new or changed production branch, constant, string, field and sink —
`ConfirmationCard`'s every field (bytes encoding, digest, challenge, created, expires, now, v, id, the namespace, the
passkey grouping and its `[]`), its state rule (outcome, expiry, the `>=` boundary), the approved-only guard and each of
its two fields, and **"no hash" / "no artifact" inserted as mutants** (`CC-HASH-FIELD`, `CC-ARTIFACT-FIELD`);
`confirm.js`'s abort, its timeout constant, its `clearTimeout`, each fetch option, the URL, the catch's state, the
abort's wording, the non-2xx and no-body branches, the path regex and its check, the base64url check, each shape check
(id, digest, challenge, integer and negative expiry, state type), the digest check, `fatal` UTF-8, the table's three
checks (unlisted, mistyped, missing) and its object check, both type predicates, **each `CONFIRM_FIELDS` entry removed
and its type flipped** (20 rows), every `textContent` sink turned to `innerHTML` (10 rows), the argv quoting (whole,
escape, safe set), the list join, the body's byte count and id, the four signed-detail rows, the message's prefix,
length, each field position, the field swap, the expiry and its byte order, the digest's source, the hash shown twice,
the blob's content and type, the save name, the namespace and the signing line's file, the controls, the ready state,
the non-pending check, the closed state, its early return, each state's text, the fallback and the approver clause;
`index.html`'s two `hidden` attributes and its script URL; the **M0 strings** (VB8's `Quote`, VP12's three hints); and
**S1b's VJ1** per S1 Appendix C (each `JudgeError.Step` label of the three judges and both lookups, `RecheckVerb`'s kind
detection — both formats — its use of `in.ExpiresAt`, `in.Digest`, `in.Now`, the approvers' `RevokedAt`, the passkey
counter, and the `admit()` mapping). For each row: an APFS clone of `work` at `S1b` (`cp -cR`; byte-identical to the
round-trip tree, C16), one exact-string replacement whose anchor occurs exactly once (checked before the run: 0 bad
anchors), a `node --check` of a mutated `.js` (a script that does not parse is `NOPARSE`, not a mutant), then **one** `go test
-count=1 -v` of the row's package, from whose `--- FAIL`/`--- PASS` lines the **named test**'s result, the package
**`suite`** (every failing test) and, for `host/surfacepage` rows, **`suite-nostatic`** (the failures other than the
static and string tests `TestConfirmPageAssets`, `TestConfirmPageFieldTable`, `TestAddressHintsArePinned`) are read —
the behavioural proof the directive asks for; the `X` row also runs XC1 in live Chrome. (Run R1 ran the three as
separate invocations and was stopped after 15 rows at load 29; its lines stay in the ledger; R2 is the full set.)
Rows ran with `GOFLAGS=-vet=off`, `WORLD_EXEC_NODE` set, six workers, at load ≈ 20–30. Every run appends one line per row
to the ledger (time, run, id, file, expected, named result, named failures, suites, named regex, seconds); per-row output
in `logs/mut/<run>-<id>.log`. **The ids below are the ledger's ids**; the table shows each row's last line. Classes:
`K` killed by the named test **and** by the behavioural suite; `S` killed only by a string pin (VP12, by design); `X`
killed by a static test and by live Chrome only; `E` equivalent (reason given); `CTL` the unmutated tree, green.

**Summary:** 145 rows (3 `CTL`, 138 `K`, 3 `S`, 1 `X`); 305 ledger lines (runs R1, R2, R3); 145 of 145 rows run; the last line of each row behaves as its class says for **145 of 145**. Table ids = ledger ids (`tools/mktable.py` fills §4 from the ledger's last line per id); ledger `$HOME/.ailang/state/world-iter401/plan-scratch/logs/mutant_ledger.tsv`

**R2 → R3 (the survivors, and what changed).** R2 (all 145 rows, on the first shipped stages) left five rows not as
expected: **`J-QUOTE-SAFE`** survived everything (no fixture argv word had a space and nothing else unsafe) → the
`draft` fixture's argv gains `--title "Fix the bar"` (M4's fixture; the golden regenerated by its Go test);
**`CC-REQUEST-ENC`** was killed only by `TestConfirmCardGolden`, not by its named `TestConfirmationCard` → that test
gains the `-__-`/`+//+` row; **`J-SHOWN`** was killed by `…Message` and `…States` but named `…RendersText` (which does
not look at `#shown`) → re-pointed to `…Message`, no test change; **`JV-ADMIT-BYTES`** and **`JV-RCK-NOW`** (S1b)
survived → three lines in `judge_verb_test.go` (S2-D1). R3 re-ran all 145 rows on the final stages; the table shows R3.
No row is equivalent.

| id | file | from → to | named test → result | package suites (`suite-nostatic`: the static and string tests skipped) | class | as expected |
|---|---|---|---|---|---|---|
| `CTL-V` | `flags.go` | control (no change) | `TestAdmitWhy` → **PASS** | suite=PASS | CTL | yes |
| `CTL-S` | `confirm.js` | control (no change) | `(TestConfirm|TestAddressHints)` → **PASS** | suite=PASS;suite-nostatic=PASS | CTL | yes |
| `CTL-A` | `judge_verb.go` | control (no change) | `(TestJudgeVerbSSH|TestJudgeVerbAssertion|Tes` → **PASS** | suite=PASS | CTL | yes |
| `VB8-JOIN` | `flags.go` | `cli := Quote(append([]string{Binary, "why"}, args...))` → `cli := strings.Join(append([]string{Binary, "why"}, args.…` | `TestAdmitWhy` → **KILLED** | suite=KILLED[TestAdmitWhy] | K | yes |
| `VP12-SUFFIX` | `app.js` | `' · episode report <episode> · why <target>';` → `' · episode report <episode> · why <ref>';` | `TestAddressHintsArePinned` → **KILLED** | suite=KILLED[TestAddressHintsArePinned];suite-nostatic=PASS | S — killed only by the string pin (`TestAddressHintsArePinned`); behavioural suites pass (as designed, S2-D8) | yes |
| `VP12-KEYROW` | `app.js` | `'in the bar: a place, episode report E, why T, approve ID'` → `'in the bar: a place, episode report E, approve ID'` | `TestAddressHintsArePinned` → **KILLED** | suite=KILLED[TestAddressHintsArePinned];suite-nostatic=PASS | S — killed only by the string pin (`TestAddressHintsArePinned`); behavioural suites pass (as designed, S2-D8) | yes |
| `VP12-INDEX` | `index.html` | ` · why &lt;target&gt;">` → `">` | `TestAddressHintsArePinned` → **KILLED** | suite=KILLED[TestAddressHintsArePinned];suite-nostatic=PASS | S — killed only by the string pin (`TestAddressHintsArePinned`); behavioural suites pass (as designed, S2-D8) | yes |
| `CC-OUTCOME` | `confirm_model.go` | `if in.Outcome != "" {` → `if false {` | `TestConfirmationCard` → **KILLED** | suite=KILLED[TestConfirmCardGolden,TestConfirmationCard,TestConfirmationCardWire];suite-nostatic=KILLED[TestConfirmCardGolden,TestConfirmationCard,Te… | K | yes |
| `CC-EXPIRY` | `confirm_model.go` | `if in.Now >= in.Expires {` → `if false {` | `TestConfirmationCard` → **KILLED** | suite=KILLED[TestConfirmCardGolden,TestConfirmationCard];suite-nostatic=KILLED[TestConfirmCardGolden,TestConfirmationCard] | K | yes |
| `CC-EXPIRY-BOUND` | `confirm_model.go` | `if in.Now >= in.Expires {` → `if in.Now > in.Expires {` | `TestConfirmationCard` → **KILLED** | suite=KILLED[TestConfirmCardGolden,TestConfirmationCard];suite-nostatic=KILLED[TestConfirmCardGolden,TestConfirmationCard] | K | yes |
| `CC-APPROVED-ONLY` | `confirm_model.go` | `if c.State == ConfirmationApproved {` → `if true {` | `TestConfirmationCard` → **KILLED** | suite=KILLED[TestConfirmationCard];suite-nostatic=KILLED[TestConfirmationCard] | K | yes |
| `CC-DECISION` | `confirm_model.go` | `c.ApprovedBy, c.Decision = in.ApprovedBy, in.Decision` → `c.ApprovedBy, c.Decision = in.ApprovedBy, ""` | `TestConfirmationCard` → **KILLED** | suite=KILLED[TestConfirmCardGolden,TestConfirmationCard,TestConfirmationCardWire];suite-nostatic=KILLED[TestConfirmCardGolden,TestConfirmationCard,Te… | K | yes |
| `CC-APPROVEDBY` | `confirm_model.go` | `c.ApprovedBy, c.Decision = in.ApprovedBy, in.Decision` → `c.ApprovedBy, c.Decision = "", in.Decision` | `TestConfirmationCard` → **KILLED** | suite=KILLED[TestConfirmCardGolden,TestConfirmationCard,TestConfirmationCardWire];suite-nostatic=KILLED[TestConfirmCardGolden,TestConfirmationCard,Te… | K | yes |
| `CC-NAMESPACE` | `confirm_model.go` | `SSH: CardSSH{Namespace: approval.VerbNamespace}` → `SSH: CardSSH{Namespace: approval.Namespace}` | `TestConfirmationCard` → **KILLED** | suite=KILLED[TestConfirmCardGolden,TestConfirmationCard,TestConfirmationCardWire];suite-nostatic=KILLED[TestConfirmCardGolden,TestConfirmationCard,Te… | K | yes |
| `CC-PASSKEYS` | `confirm_model.go` | `Passkeys: RequestCard(StoredRequest{}, passkeys).Passkeys` → `Passkeys: []CardPasskeys{}` | `TestConfirmationCard` → **KILLED** | suite=KILLED[TestConfirmCardGolden,TestConfirmationCard];suite-nostatic=KILLED[TestConfirmCardGolden,TestConfirmationCard] | K | yes |
| `CC-PASSKEYS-NIL` | `confirm_model.go` | `Passkeys: RequestCard(StoredRequest{}, passkeys).Passkeys` → `Passkeys: nil` | `TestConfirmationCard` → **KILLED** | suite=KILLED[TestConfirmCardGolden,TestConfirmationCard];suite-nostatic=KILLED[TestConfirmCardGolden,TestConfirmationCard] | K | yes |
| `CC-REQUEST-ENC` | `confirm_model.go` | `Request: base64.RawURLEncoding.EncodeToString(in.Request)` → `Request: base64.StdEncoding.EncodeToString(in.Request)` | `TestConfirmationCard` → **KILLED** | suite=KILLED[TestConfirmCardGolden,TestConfirmationCard];suite-nostatic=KILLED[TestConfirmCardGolden,TestConfirmationCard] | K | yes |
| `CC-DIGEST` | `confirm_model.go` | `Digest: in.Digest, Challenge: in.Challenge,` → `Digest: in.Challenge, Challenge: in.Challenge,` | `TestConfirmationCard` → **KILLED** | suite=KILLED[TestConfirmCardGolden,TestConfirmationCard];suite-nostatic=KILLED[TestConfirmCardGolden,TestConfirmationCard] | K | yes |
| `CC-CHALLENGE` | `confirm_model.go` | `Digest: in.Digest, Challenge: in.Challenge,` → `Digest: in.Digest, Challenge: in.Digest,` | `TestConfirmationCard` → **KILLED** | suite=KILLED[TestConfirmCardGolden,TestConfirmationCard];suite-nostatic=KILLED[TestConfirmCardGolden,TestConfirmationCard] | K | yes |
| `CC-CREATED` | `confirm_model.go` | `Created: in.Created, Expires: in.Expires, Now: in.Now,` → `Created: in.Expires, Expires: in.Expires, Now: in.Now,` | `TestConfirmationCard` → **KILLED** | suite=KILLED[TestConfirmCardGolden,TestConfirmationCard];suite-nostatic=KILLED[TestConfirmCardGolden,TestConfirmationCard] | K | yes |
| `CC-EXPIRES` | `confirm_model.go` | `Created: in.Created, Expires: in.Expires, Now: in.Now,` → `Created: in.Created, Expires: in.Created, Now: in.Now,` | `TestConfirmationCard` → **KILLED** | suite=KILLED[TestConfirmCardGolden,TestConfirmationCard];suite-nostatic=KILLED[TestConfirmCardGolden,TestConfirmationCard] | K | yes |
| `CC-NOW` | `confirm_model.go` | `Created: in.Created, Expires: in.Expires, Now: in.Now,` → `Created: in.Created, Expires: in.Expires, Now: in.Created,` | `TestConfirmationCard` → **KILLED** | suite=KILLED[TestConfirmCardGolden,TestConfirmationCard];suite-nostatic=KILLED[TestConfirmCardGolden,TestConfirmationCard] | K | yes |
| `CC-V` | `confirm_model.go` | `Confirmation{V: 1,` → `Confirmation{V: 2,` | `TestConfirmationCard` → **KILLED** | suite=KILLED[TestConfirmCardGolden,TestConfirmationCard];suite-nostatic=KILLED[TestConfirmCardGolden,TestConfirmationCard] | K | yes |
| `CC-ID` | `confirm_model.go` | `ID: in.ID, Request:` → `ID: "", Request:` | `TestConfirmationCard` → **KILLED** | suite=KILLED[TestConfirmCardGolden,TestConfirmationCard];suite-nostatic=KILLED[TestConfirmCardGolden,TestConfirmationCard] | K | yes |
| `CC-HASH-FIELD` | `confirm_model.go` | `Decision string `json:"decision"`⏎` → `Decision string `json:"decision"`⏎Hash string `json:"hash…` | `TestConfirmationCardWire` → **KILLED** | suite=KILLED[TestConfirmCardGolden,TestConfirmationCardWire];suite-nostatic=KILLED[TestConfirmCardGolden,TestConfirmationCardWire] | K | yes |
| `CC-ARTIFACT-FIELD` | `confirm_model.go` | `ApprovedBy, Decision string // the approver` → `Artifact []byte⏎ApprovedBy, Decision string // the approv…` | `TestConfirmationCardWire` → **KILLED** | suite=KILLED[TestConfirmationCardWire];suite-nostatic=KILLED[TestConfirmationCardWire] | K | yes |
| `J-ABORT` | `confirm.js` | `const timer = setTimeout(() => c.abort(), NET_MS);` → `const timer = setTimeout(() => {}, NET_MS);` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-NET-MS` | `confirm.js` | `const NET_MS = 10000;` → `const NET_MS = 20000;` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-CLEAR` | `confirm.js` | ` } finally {⏎ clearTimeout(timer);⏎ }` → ` } finally {⏎ }` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-NOSTORE` | `confirm.js` | `cache: 'no-store', ` → `` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-OMIT` | `confirm.js` | `credentials: 'omit', ` → `` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-SIGNAL` | `confirm.js` | `signal: c.signal, ` → `` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-URL` | `confirm.js` | `net('/surface/confirmations/' + id)` → `net('/surface/requests/' + id)` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-CATCH-STATE` | `confirm.js` | ` return refuse('failed', `The request could not be read` → ` return refuse('refused', `The request could not be read` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-WHY-ABORT` | `confirm.js` | `e.name === 'AbortError'` → `e.name === 'TimeoutError'` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-NON2XX` | `confirm.js` | `if (res.status !== 200 \|\| !res.body) {` → `if (!res.body) {` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-NOBODY` | `confirm.js` | `if (res.status !== 200 \|\| !res.body) {` → `if (res.status !== 200) {` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-PATH-RE` | `confirm.js` | `(vr-[0-9a-f]{16})$/` → `(vr-[0-9a-f]+)$/` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-PATH-CHECK` | `confirm.js` | `if (!id) return refuse(` → `if (false) return refuse(` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-B64-CHECK` | `confirm.js` | `if (typeof s !== 'string' \|\| !/^[A-Za-z0-9_-]*$/.test(s…` → `` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-ID` | `confirm.js` | `if (s.id !== id \|\| !unhex(s.digest)` → `if (!unhex(s.digest)` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-DIGEST-FMT` | `confirm.js` | ` \|\| !unhex(s.digest) \|\| !challenge` → ` \|\| !challenge` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-CHALLENGE-FMT` | `confirm.js` | ` \|\| !challenge \|\| !Number` → ` \|\| !Number` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-EXP-INT` | `confirm.js` | ` \|\| !Number.isSafeInteger(s.expires) \|\| s.expires < 0` → ` \|\| s.expires < 0` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-EXP-NEG` | `confirm.js` | ` \|\| !Number.isSafeInteger(s.expires) \|\| s.expires < 0` → ` \|\| !Number.isSafeInteger(s.expires)` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-STATE-TYPE` | `confirm.js` | ` \|\| typeof s.state !== 'string') {` → `) {` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-DIGEST-CHECK` | `confirm.js` | `if (hex(digest) !== s.digest) return refuse(` → `if (false) return refuse(` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-UTF8-FATAL` | `confirm.js` | `{ fatal: true }` → `{ fatal: false }` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-UNLISTED` | `confirm.js` | `if (!f \|\| !typed(f, r[k])) return k;` → `if (f && !typed(f, r[k])) return k;` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-MISTYPED` | `confirm.js` | `if (!f \|\| !typed(f, r[k])) return k;` → `if (!f) return k;` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-MISSING` | `confirm.js` | ` for (const f of CONFIRM_FIELDS) if (!Object.prototype.ha…` → `` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-OBJECT` | `confirm.js` | `if (r === null \|\| typeof r !== 'object' \|\| Array.isAr…` → `if (r === null) return 'the request';` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-TYPED-STRINGS` | `confirm.js` | `Array.isArray(v) && v.every((x) => typeof x === 'string')` → `Array.isArray(v)` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-TYPED-STRING` | `confirm.js` | `f.type === 'string' ? typeof v === 'string'` → `f.type === 'string' ? v !== undefined` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-FIELD-schema` | `confirm.js` | ` {"key": "schema", "label": "Schema", "type": "string"},⏎` → `` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageFieldTable,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nost… | K | yes |
| `J-FIELD-verb` | `confirm.js` | ` {"key": "verb", "label": "Verb", "type": "string"},⏎` → `` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageFieldTable,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nost… | K | yes |
| `J-FIELD-argv` | `confirm.js` | ` {"key": "argv", "label": "Command words", "type": "strin…` → `` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageFieldTable,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nost… | K | yes |
| `J-FIELD-episode` | `confirm.js` | ` {"key": "episode", "label": "Episode", "type": "string"}…` → `` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageFieldTable,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nost… | K | yes |
| `J-FIELD-prompt` | `confirm.js` | ` {"key": "prompt", "label": "Prompt", "type": "string"},⏎` → `` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageFieldTable,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nost… | K | yes |
| `J-FIELD-scope` | `confirm.js` | ` {"key": "scope", "label": "Approval scope", "type": "str…` → `` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageFieldTable,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nost… | K | yes |
| `J-FIELD-title` | `confirm.js` | ` {"key": "title", "label": "Title", "type": "string"},⏎` → `` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageFieldTable,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nost… | K | yes |
| `J-FIELD-body` | `confirm.js` | ` {"key": "body", "label": "Body", "type": "string"},⏎` → `` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageFieldTable,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nost… | K | yes |
| `J-FIELD-committer` | `confirm.js` | ` {"key": "committer", "label": "Committer", "type": "stri…` → `` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageFieldTable,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nost… | K | yes |
| `J-FIELD-protected` | `confirm.js` | ` {"key": "protected", "label": "Protected paths changed",…` → `⏎` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageFieldTable,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nost… | K | yes |
| `J-TYPE-schema` | `confirm.js` | `{"key": "schema", "label": "Schema", "type": "string"}` → `{"key": "schema", "label": "Schema", "type": "strings"}` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageFieldTable,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nost… | K | yes |
| `J-TYPE-verb` | `confirm.js` | `{"key": "verb", "label": "Verb", "type": "string"}` → `{"key": "verb", "label": "Verb", "type": "strings"}` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageFieldTable,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nost… | K | yes |
| `J-TYPE-argv` | `confirm.js` | `… "type": "strings"}` → `… "type": "string"}` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageFieldTable,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nost… | K | yes |
| `J-TYPE-episode` | `confirm.js` | `{"key": "episode", "label": "Episode", "type": "string"}` → `{"key": "episode", "label": "Episode", "type": "strings"}` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageFieldTable,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nost… | K | yes |
| `J-TYPE-prompt` | `confirm.js` | `{"key": "prompt", "label": "Prompt", "type": "string"}` → `{"key": "prompt", "label": "Prompt", "type": "strings"}` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageFieldTable,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nost… | K | yes |
| `J-TYPE-scope` | `confirm.js` | `… "type": "string"}` → `… "type": "strings"}` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageFieldTable,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nost… | K | yes |
| `J-TYPE-title` | `confirm.js` | `{"key": "title", "label": "Title", "type": "string"}` → `{"key": "title", "label": "Title", "type": "strings"}` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageFieldTable,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nost… | K | yes |
| `J-TYPE-body` | `confirm.js` | `{"key": "body", "label": "Body", "type": "string"}` → `{"key": "body", "label": "Body", "type": "strings"}` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageFieldTable,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nost… | K | yes |
| `J-TYPE-committer` | `confirm.js` | `… "type": "string"}` → `… "type": "strings"}` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageFieldTable,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nost… | K | yes |
| `J-TYPE-protected` | `confirm.js` | `… "type": "strings"}` → `… "type": "string"}` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageFieldTable,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nost… | K | yes |
| `J-SINK-ARGV` | `confirm.js` | `$('argv').textContent = quote(r.argv);` → `$('argv').innerHTML = quote(r.argv);` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageAssets,TestConfirmPageRendersText];suite-nostatic=KILLED[TestConfirmPageRendersText] | K | yes |
| `J-SINK-PROMPT` | `confirm.js` | `$('prompt').textContent = r.prompt;` → `$('prompt').innerHTML = r.prompt;` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageAssets,TestConfirmPageRendersText];suite-nostatic=KILLED[TestConfirmPageRendersText] | K | yes |
| `J-SINK-DD` | `confirm.js` | ` dd.textContent = text;` → ` dd.innerHTML = text;` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageAssets,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nostatic=KILLED[TestConfirmPageMess… | K | yes |
| `J-SINK-DT` | `confirm.js` | ` dt.textContent = label;` → ` dt.innerHTML = label;` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageAssets,TestConfirmPageRendersText];suite-nostatic=KILLED[TestConfirmPageRendersText] | K | yes |
| `J-SINK-BODY` | `confirm.js` | ` pre.textContent = v;` → ` pre.innerHTML = v;` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageAssets,TestConfirmPageRendersText];suite-nostatic=KILLED[TestConfirmPageRendersText] | K | yes |
| `J-SINK-SUMMARY` | `confirm.js` | `summary.textContent = 'the full body';` → `summary.innerHTML = 'the full body';` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageAssets,TestConfirmPageRendersText];suite-nostatic=KILLED[TestConfirmPageRendersText] | K | yes |
| `J-SINK-WHAT` | `confirm.js` | `$('what').textContent = `Confirm request` → `$('what').innerHTML = `Confirm request` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageAssets,TestConfirmPageMessage,TestConfirmPageRendersText];suite-nostatic=KILLED[TestConfirmPageMessage,TestConfirmPageRen… | K | yes |
| `J-SINK-REFUSE` | `confirm.js` | `setState(state); $('what').textContent = text;` → `setState(state); $('what').innerHTML = text;` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageAssets,TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `J-SINK-STATUS` | `confirm.js` | ` $('status').textContent = `${CLOSED[` → ` $('status').innerHTML = `${CLOSED[` | `TestConfirmPageStates` → **KILLED** | suite=KILLED[TestConfirmPageAssets,TestConfirmPageStates];suite-nostatic=KILLED[TestConfirmPageStates] | K | yes |
| `J-SINK-SSHLINE` | `confirm.js` | `$('sshline').textContent = `ssh-keygen` → `$('sshline').innerHTML = `ssh-keygen` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageAssets,TestConfirmPageMessage,TestConfirmPageRendersText];suite-nostatic=KILLED[TestConfirmPageMessage,TestConfirmPageRen… | K | yes |
| `J-QUOTE` | `confirm.js` | `argv.map((a) => (SAFE.test(a) ? a : "'" + a.split("'").jo…` → `argv.join(' ')` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageRendersText];suite-nostatic=KILLED[TestConfirmPageRendersText] | K | yes |
| `J-QUOTE-ESC` | `confirm.js` | `a.split("'").join(`'"'"'`)` → `a` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageRendersText];suite-nostatic=KILLED[TestConfirmPageRendersText] | K | yes |
| `J-QUOTE-SAFE` | `confirm.js` | `const SAFE = /^[A-Za-z0-9@%+=:,./_-]+$/;` → `const SAFE = /^[A-Za-z0-9@%+=:,./_ -]+$/;` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageRendersText];suite-nostatic=KILLED[TestConfirmPageRendersText] | K | yes |
| `J-LIST-JOIN` | `confirm.js` | `f.type === 'strings' ? v.join('\n') : v` → `f.type === 'strings' ? v.join(' ') : v` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageRendersText];suite-nostatic=KILLED[TestConfirmPageRendersText] | K | yes |
| `J-BODY-SIZE` | `confirm.js` | `${new TextEncoder().encode(v).length} bytes` → `${v.length} bytes` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageRendersText];suite-nostatic=KILLED[TestConfirmPageRendersText] | K | yes |
| `J-BODY-ID` | `confirm.js` | `pre.id = 'f-body';` → `pre.id = 'f-body-text';` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageRendersText];suite-nostatic=KILLED[TestConfirmPageRendersText] | K | yes |
| `J-M-DIGEST` | `confirm.js` | `row(dl, 'm-digest', 'Request digest (SHA-256)', s.digest);` → `row(dl, 'm-digest', 'Request digest (SHA-256)', s.challen…` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageRendersText];suite-nostatic=KILLED[TestConfirmPageRendersText] | K | yes |
| `J-M-CHALLENGE` | `confirm.js` | `row(dl, 'm-challenge', 'Challenge', s.challenge);` → `row(dl, 'm-challenge', 'Challenge', s.digest);` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageRendersText];suite-nostatic=KILLED[TestConfirmPageRendersText] | K | yes |
| `J-M-EXPIRES` | `confirm.js` | `row(dl, 'm-expires', 'Expires', time(s.expires));` → `row(dl, 'm-expires', 'Expires', time(s.created));` | `TestConfirmPageRendersText` → **KILLED** | suite=KILLED[TestConfirmPageRendersText];suite-nostatic=KILLED[TestConfirmPageRendersText] | K | yes |
| `J-SHOWN` | `confirm.js` | ` $('shown').hidden = false;⏎` → `` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageMessage,TestConfirmPageStates];suite-nostatic=KILLED[TestConfirmPageMessage,TestConfirmPageStates] | K | yes |
| `J-MSG-PREFIX` | `confirm.js` | `PREFIX = 'world-verb/v1'` → `PREFIX = 'world-vexb/v1'` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageMessage];suite-nostatic=KILLED[TestConfirmPageMessage] | K | yes |
| `J-MSG-LEN` | `confirm.js` | `LEN = 86;` → `LEN = 87;` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageMessage];suite-nostatic=KILLED[TestConfirmPageMessage] | K | yes |
| `J-MSG-DIGEST-AT` | `confirm.js` | `m.set(digest, 14);` → `m.set(digest, 13);` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageMessage];suite-nostatic=KILLED[TestConfirmPageMessage] | K | yes |
| `J-MSG-CHALLENGE-AT` | `confirm.js` | `m.set(challenge, 46);` → `m.set(challenge, 45);` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageMessage];suite-nostatic=KILLED[TestConfirmPageMessage] | K | yes |
| `J-MSG-SWAP` | `confirm.js` | ` m.set(digest, 14);⏎ m.set(challenge, 46);` → ` m.set(challenge, 14);⏎ m.set(digest, 46);` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageMessage];suite-nostatic=KILLED[TestConfirmPageMessage] | K | yes |
| `J-MSG-EXPIRY-AT` | `confirm.js` | `setBigUint64(78, BigInt(expiry))` → `setBigUint64(77, BigInt(expiry))` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageMessage];suite-nostatic=KILLED[TestConfirmPageMessage] | K | yes |
| `J-MSG-EXPIRY` | `confirm.js` | `setBigUint64(78, BigInt(expiry))` → `setBigUint64(78, BigInt(expiry + 1))` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageMessage];suite-nostatic=KILLED[TestConfirmPageMessage] | K | yes |
| `J-MSG-LE` | `confirm.js` | `setBigUint64(78, BigInt(expiry))` → `setBigUint64(78, BigInt(expiry), true)` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageMessage];suite-nostatic=KILLED[TestConfirmPageMessage] | K | yes |
| `J-MSG-DIGEST-SRC` | `confirm.js` | `const m = message(digest, challenge, s.expires)` → `const m = message(unhex(s.challenge), challenge, s.expire…` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageMessage];suite-nostatic=KILLED[TestConfirmPageMessage] | K | yes |
| `J-M-HASH` | `confirm.js` | `…hallenge)', hex(h));` → `…hallenge)', hex(digest));` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nostatic=KILLED[TestConfirmPag… | K | yes |
| `J-WHAT-HASH` | `confirm.js` | `; message hash ${hex(h)}`;` → `; message hash ${hex(digest)}`;` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageMessage];suite-nostatic=KILLED[TestConfirmPageMessage] | K | yes |
| `J-BLOB` | `confirm.js` | `new Blob([m], {` → `new Blob([h], {` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageMessage];suite-nostatic=KILLED[TestConfirmPageMessage] | K | yes |
| `J-BLOB-TYPE` | `confirm.js` | `{ type: 'application/octet-stream' }` → `{ type: 'text/plain' }` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageMessage];suite-nostatic=KILLED[TestConfirmPageMessage] | K | yes |
| `J-SAVE-NAME` | `confirm.js` | `$('save').download = `world-verb-${id}.msg`;` → `$('save').download = `world-verb.msg`;` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageMessage];suite-nostatic=KILLED[TestConfirmPageMessage] | K | yes |
| `J-NAMESPACE` | `confirm.js` | `NAMESPACE = 'world-verb@sunholo'` → `NAMESPACE = 'world-approval@sunholo'` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageMessage];suite-nostatic=KILLED[TestConfirmPageMessage] | K | yes |
| `J-SSHLINE-FILE` | `confirm.js` | `-f <key> world-verb-${id}.msg`;` → `-f <key> world-verb.msg`;` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageMessage];suite-nostatic=KILLED[TestConfirmPageMessage] | K | yes |
| `J-CONTROLS` | `confirm.js` | ` $('controls').hidden = false;⏎` → `` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageMessage,TestConfirmPageStates];suite-nostatic=KILLED[TestConfirmPageMessage,TestConfirmPageStates] | K | yes |
| `J-READY` | `confirm.js` | ` setState('ready');⏎ }` → ` setState('pending');⏎ }` | `TestConfirmPageMessage` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageMessage,TestConfirmPageRendersText,TestConfirmPageStates];suite-nostatic=KILLED[TestConfirmPag… | K | yes |
| `J-STATE-CHECK` | `confirm.js` | `if (s.state !== 'pending') {` → `if (false) {` | `TestConfirmPageStates` → **KILLED** | suite=KILLED[TestConfirmPageStates];suite-nostatic=KILLED[TestConfirmPageStates] | K | yes |
| `J-STATE-CLOSED` | `confirm.js` | ` setState('closed');` → ` setState('ready');` | `TestConfirmPageStates` → **KILLED** | suite=KILLED[TestConfirmPageStates];suite-nostatic=KILLED[TestConfirmPageStates] | K | yes |
| `J-CLOSED-RETURN` | `confirm.js` | `nothing to confirm here.`;⏎ return;⏎` → `nothing to confirm here.`;⏎` | `TestConfirmPageStates` → **KILLED** | suite=KILLED[TestConfirmPageStates];suite-nostatic=KILLED[TestConfirmPageStates] | K | yes |
| `J-CLOSED-APPROVED` | `confirm.js` | `approved: 'Already confirmed'` → `approved: 'Confirmed'` | `TestConfirmPageStates` → **KILLED** | suite=KILLED[TestConfirmPageStates];suite-nostatic=KILLED[TestConfirmPageStates] | K | yes |
| `J-CLOSED-SUPERSEDED` | `confirm.js` | `superseded: 'Superseded by a newer finish of this episode'` → `superseded: 'Superseded'` | `TestConfirmPageStates` → **KILLED** | suite=KILLED[TestConfirmPageStates];suite-nostatic=KILLED[TestConfirmPageStates] | K | yes |
| `J-CLOSED-EXPIRED` | `confirm.js` | `expired: 'Expired' }` → `expired: 'Gone' }` | `TestConfirmPageStates` → **KILLED** | suite=KILLED[TestConfirmPageStates];suite-nostatic=KILLED[TestConfirmPageStates] | K | yes |
| `J-CLOSED-FALLBACK` | `confirm.js` | `${CLOSED[s.state] \|\| 'This request is ' + s.state}` → `${CLOSED[s.state]}` | `TestConfirmPageStates` → **KILLED** | suite=KILLED[TestConfirmPageStates];suite-nostatic=KILLED[TestConfirmPageStates] | K | yes |
| `J-APPROVED-BY` | `confirm.js` | `${s.state === 'approved' ? ` by ${s.approvedBy}: decision…` → `` | `TestConfirmPageStates` → **KILLED** | suite=KILLED[TestConfirmPageStates];suite-nostatic=KILLED[TestConfirmPageStates] | K | yes |
| `H-SHOWN-HIDDEN` | `index.html` | `<div id="shown" hidden>` → `<div id="shown">` | `TestConfirmPageFailsClosed` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed];suite-nostatic=KILLED[TestConfirmPageFailsClosed] | K | yes |
| `H-CONTROLS-HIDDEN` | `index.html` | `<div id="controls" hidden>` → `<div id="controls">` | `TestConfirmPageStates` → **KILLED** | suite=KILLED[TestConfirmPageFailsClosed,TestConfirmPageStates];suite-nostatic=KILLED[TestConfirmPageFailsClosed,TestConfirmPageStates] | K | yes |
| `H-SCRIPT-SRC` | `index.html` | `<script src="/surface/confirm.js"></script>` → `<script src="confirm.js"></script>` | `TestConfirmPageAssets` → **KILLED** | suite=KILLED[TestConfirmPageAssets];suite-nostatic=PASS;chrome=KILLED[TestConfirmPageShowsTheFinishPromptInChrome] | X — killed by the static assets test and by live Chrome (XC1); the node suites cannot see a script URL | yes |
| `JV-SSH-PARSE` | `judge_verb.go` | `…{Step: JudgeStepParse` → `…{Step: JudgeStepLaw` | `TestJudgeVerbSSH` → **KILLED** | suite=KILLED[TestJudgeVerbSSH] | K | yes |
| `JV-ASR-PARSE` | `judge_verb.go` | `…{Step: JudgeStepParse` → `…{Step: JudgeStepLaw` | `TestJudgeVerbAssertion` → **KILLED** | suite=KILLED[TestJudgeVerbAssertion] | K | yes |
| `JV-SSH-STORE` | `judge_verb.go` | `…{Step: JudgeStepStore` → `…{Step: JudgeStepLaw` | `TestJudgeVerbSSH` → **KILLED** | suite=KILLED[TestJudgeVerbSSH,TestRecheckVerb] | K | yes |
| `JV-PK-STORE` | `judge_verb.go` | `…{Step: JudgeStepStore` → `…{Step: JudgeStepLaw` | `TestJudgeVerbAssertion` → **KILLED** | suite=KILLED[TestJudgeVerbAssertion,TestRecheckVerb] | K | yes |
| `JV-SSH-UNKNOWN` | `judge_verb.go` | `Step: JudgeStepUnknownKey, Fingerprint:` → `Step: JudgeStepLaw, Fingerprint:` | `TestJudgeVerbSSH` → **KILLED** | suite=KILLED[TestJudgeVerbSSH,TestRecheckVerb] | K | yes |
| `JV-PK-UNKNOWN` | `judge_verb.go` | `Step: JudgeStepUnknownKey, RawID:` → `Step: JudgeStepLaw, RawID:` | `TestJudgeVerbAssertion` → **KILLED** | suite=KILLED[TestJudgeVerbAssertion,TestRecheckVerb] | K | yes |
| `JV-PK-STOREDKEY` | `judge_verb.go` | `Step: JudgeStepStoredKey, Approver: a.Name` → `Step: JudgeStepLaw, Approver: a.Name` | `TestJudgeVerbAssertion` → **KILLED** | suite=KILLED[TestJudgeVerbAssertion,TestRecheckVerb] | K | yes |
| `JV-SSH-LAW` | `judge_verb.go` | `…{Step: JudgeStepLaw` → `…{Step: JudgeStepParse` | `TestJudgeVerbSSH` → **KILLED** | suite=KILLED[TestJudgeVerbSSH] | K | yes |
| `JV-ASR-LAW` | `judge_verb.go` | `…{Step: JudgeStepLaw` → `…{Step: JudgeStepParse` | `TestJudgeVerbAssertion` → **KILLED** | suite=KILLED[TestJudgeVerbAssertion] | K | yes |
| `JV-RCK-LAW` | `judge_verb.go` | `…{Step: JudgeStepLaw` → `…{Step: JudgeStepParse` | `TestRecheckVerb` → **KILLED** | suite=KILLED[TestRecheckVerb] | K | yes |
| `JV-RCK-NEITHER` | `judge_verb.go` | `…{Step: JudgeStepParse, Err: errors.Join(` → `…{Step: JudgeStepLaw, Err: errors.Join(` | `TestRecheckVerb` → **KILLED** | suite=KILLED[TestRecheckVerb] | K | yes |
| `JV-RCK-EXPIRY` | `judge_verb.go` | `Expiry: in.ExpiresAt, Now: in.Now,` → `Expiry: in.ExpiresAt + 1, Now: in.Now,` | `TestRecheckVerb` → **KILLED** | suite=KILLED[TestRecheckVerb] | K | yes |
| `JV-RCK-DIGEST` | `judge_verb.go` | `check := approval.VerbCheckInput{Digest: in.Digest,` → `check := approval.VerbCheckInput{Digest: in.BytesDigest,` | `TestRecheckVerb` → **KILLED** | suite=KILLED[TestRecheckVerb] | K | yes |
| `JV-RCK-NOW` | `judge_verb.go` | `Expiry: in.ExpiresAt, Now: in.Now,` → `Expiry: in.ExpiresAt, Now: in.CreatedAt,` | `TestRecheckVerb` → **KILLED** | suite=KILLED[TestRecheckVerb] | K | yes |
| `JV-RCK-SSH-FORMAT` | `judge_verb.go` | `check.Format, check.Approver = approval.FormatSSH, a` → `check.Format, check.Approver = approval.FormatWebAuthn, a` | `TestRecheckVerb` → **KILLED** | suite=KILLED[TestRecheckVerb] | K | yes |
| `JV-RCK-PK-FORMAT` | `judge_verb.go` | `check.Format, check.Approver = approval.FormatWebAuthn, a` → `check.Format, check.Approver = approval.FormatSSH, a` | `TestRecheckVerb` → **KILLED** | suite=KILLED[TestRecheckVerb] | K | yes |
| `JV-SSH-REVOKED` | `judge_verb.go` | `…olled.LastUsedAt, RevokedAt: enrolled.RevokedAt}, nil` → `…olled.LastUsedAt}, nil` | `TestJudgeVerbSSH` → **KILLED** | suite=KILLED[TestJudgeVerbSSH,TestRecheckVerb] | K | yes |
| `JV-PK-REVOKED` | `judge_verb.go` | `LastUsedAt: a.LastUsedAt, RevokedAt: a.RevokedAt,` → `LastUsedAt: a.LastUsedAt,` | `TestRecheckVerb` → **KILLED** | suite=KILLED[TestRecheckVerb] | K | yes |
| `JV-PK-COUNT` | `judge_verb.go` | `SignCount: uint32(p.SignCount)}}, nil` → `SignCount: 0}}, nil` | `TestJudgeVerbAssertion` → **KILLED** | suite=KILLED[TestJudgeVerbAssertion] | K | yes |
| `JV-ADMIT-CONSUMED` | `judge_verb.go` | `Consumed: in.Consumed,` → `Consumed: false,` | `TestJudgeVerbSSH` → **KILLED** | suite=KILLED[TestJudgeVerbSSH] | K | yes |
| `JV-ADMIT-BYTES` | `judge_verb.go` | `StoredDigest: in.Digest, BytesDigest: in.BytesDigest}` → `StoredDigest: in.Digest, BytesDigest: in.Digest}` | `TestJudgeVerbSSH` → **KILLED** | suite=KILLED[TestJudgeVerbSSH] | K | yes |
| `JV-ADMIT-EXPIRES` | `judge_verb.go` | `CreatedAt: in.CreatedAt, ExpiresAt: in.ExpiresAt, Consumed` → `CreatedAt: in.CreatedAt, ExpiresAt: in.ExpiresAt + 100000…` | `TestJudgeVerbAssertion` → **KILLED** | suite=KILLED[TestJudgeVerbAssertion,TestJudgeVerbSSH] | K | yes |

---

## 5. Files (final `git status --short --untracked-files=all` after S1b, against the base)

```
A  cmd/ailang-worldd/confirm_page_chrome_test.go
A  host/authority/judge_verb.go
A  host/authority/judge_verb_test.go
A  host/daemon/surface_confirm_unrouted_test.go
A  host/surfacepage/app_hints_test.go
A  host/surfacepage/confirm.go
A  host/surfacepage/confirm/confirm.js
A  host/surfacepage/confirm/index.html
A  host/surfacepage/confirm_model.go
A  host/surfacepage/confirm_model_test.go
A  host/surfacepage/confirm_page_test.go
A  host/surfacepage/confirm_script_test.go
A  host/surfacepage/testdata/confirm-card.golden.jsonl
M  host/verbs/why_admit_test.go
```

Plus this plan and its JSON under `design_docs/planned/` (the controller's commit). Nothing deleted. **No production
file outside `host/surfacepage` (new files only) and `host/authority` (new files only).**

---

## 6. Sprint gates (controller, after the last sub-milestone's commit; each measured at the base, §1.1)

One bash shell, never under `nohup`: `bash gates.sh <worktree> $HOME/.ailang/state/world-iter401/gates final <HEAD sha>`
(extract it as `accept.sh` is extracted, with the fence `` ````bash gates.sh ``). Every leg is `tee`'d to
`$HOME/.ailang/state/world-iter401/gates/<leg>.txt`; `ALL_DONE.<sha>` is written last with `FAILS: …`. The executor sets
`GATES_SKIP_VERIFY_GO=1`; the controller runs `verify_go`. Legs: `status`, `verify_ail`, `gofmt`, `build`, `vet`,
`fence` (`go test -run '^$' ./...`), `pkgs` (the four touched packages), `race`, **`node`** (every `host/surfacepage`
node harness, `WORLD_EXEC_NODE` set, `--- PASS` required per name), **`inv`** (the inventories by name with `-v`: route
tables incl. `TestConfirmRoutesNotRegistered`, `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist`,
`TestProductionContextRoots`, `TestProductionGoSurface`, `TestEveryVerbHasHelp`, `TestFrozenV1RouteTableMatchesMux`,
the A2/A3/A4 page and route tests, the embedded-font test, S2's static tests, S1b's judges), **`chrome`** (V2b's
fourteen + the two `host/daemon` drills + XC1; `WORLD_CHROME`, `WORLD_REQUIRE_PASSKEY_BROWSER=1`; temp profiles), then
`verify_go` (`VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh`) unless skipped.

````bash gates.sh
#!/bin/bash
# Row 243 V3 S2 gates (plan §6). usage: gates.sh <tree> <logdir> <base|final> [<sha>]
# The controller: bash gates.sh <worktree> $HOME/.ailang/state/world-iter401/gates final <HEAD sha>
# macOS /bin/bash 3.2: no declare -A. Every leg records its own rc; every leg runs (no set -e); every leg is tee'd
# to <logdir>/<leg>.txt. GATES_SKIP_VERIFY_GO=1 skips verify_go (the executor sets it; the controller runs it).
# Never start this under nohup. Writes <logdir>/ALL_DONE.<sha> with FAILS: when every leg has finished.
T=$1; L=$2; MODE=$3
mkdir -p "$L"; cd "$T" || exit 2
SHA=${4:-$(git rev-parse HEAD 2>/dev/null || echo unknown)}
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_REQUIRE_SSH_KEYGEN WORLD_PASSKEY_RECORD_DIR WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_TOOL_AILANG_BIN
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
leg status 'git status --porcelain | wc -l'
leg verify_ail './scripts/verify_ail.sh'
leg gofmt 'out=$(gofmt -l cmd host); echo "$out"; test -z "$out"'
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./... > '$L/fence.raw' 2>&1; rc=\$?; cat '$L/fence.raw'; if grep -E '^(FAIL|--- FAIL)' '$L/fence.raw' >/dev/null; then exit 1; fi; exit \$rc"
leg pkgs "go test -count=1 -timeout 30m ./host/verbs/ ./host/surfacepage/ ./host/daemon/ ./cmd/ailang-worldd/"
leg race "go test -race -count=1 ./host/verbs/ ./host/surfacepage/ && go test -race -count=1 ./host/daemon/ -run '^(TestConfirm|TestApprovalRoutesRegistered|TestDecisionRoutesRegistered|TestDecidePage|TestApprovePage)'"
# node: every page harness, non-skippable (WORLD_EXEC_NODE is set above; a harness that cannot run FAILS)
NODE_BASE="TestSurfaceAppPure TestSurfaceScriptPure TestSurfaceAppVerbs TestNeedsYouListsDecisions"
NODE_NEW="TestAddressHintsArePinned TestConfirmPageFailsClosed TestConfirmPageRendersText TestConfirmPageMessage TestConfirmPageStates"
[ "$MODE" = base ] && NODE_NEW=""
leg node "set -o pipefail; go test -v -count=1 ./host/surfacepage/ -run '^($(bar $NODE_BASE)$(alt $NODE_NEW))\$'"
need node $NODE_BASE $NODE_NEW
if grep -E '^\s*--- (SKIP|FAIL)' "$L/node.txt"; then FAILS="$FAILS node_skip_or_fail"; fi
NEW_V="TestAdmitWhy"
NEW_SP="TestConfirmationCard TestConfirmationCardWire TestConfirmPageFieldTable TestConfirmCardGolden TestConfirmPageAssets"
NEW_D="TestConfirmRoutesNotRegistered"
if [ "$MODE" = base ]; then NEW_SP=""; NEW_D=""; fi
INV_DAEMON="TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestDaemonCredentialWriterSites TestSurfaceVerbRouteRegistered TestSurfaceVerbIsReadOnly TestDecidePageAssets TestDecidePageFieldTable TestApprovePageRoutes TestPasskeyPageRoutes"
INV_BOUNDARY="TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard TestWorkbenchPackageRemainsTransportFree TestPinfetchStaysOutOfStoreReplayAndPublish"
INV_STORE="TestProductionContextRoots TestProductionGoSurface TestStoreSurfaceTakesContext TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext"
INV_VG="TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestVerifygateTestWritesAreForkLocked"
INV_SP="TestSurfaceAppMatchesWorkingCopy TestSurfaceAppShellAssembly TestSurfaceFontEmbed"
INV_VERBS="TestSplit TestQuoteRoundTrip TestLookup TestApproveTarget TestAdmitWhy"
INV_CMD="TestEveryVerbHasHelp TestVerbTableMatchesTheCLISwitch TestFinishPromptIsTheGolden"
leg inv "set -o pipefail
go test -v -count=1 ./host/daemon/ -run '^($(bar $INV_DAEMON)$(alt $NEW_D))\$' &&
go test -v -count=1 ./host/boundary/ -run '^($(bar $INV_BOUNDARY))\$' &&
go test -v -count=1 ./host/store/ -run '^($(bar $INV_STORE))\$' &&
go test -v -count=1 ./host/verifygate/ -run '^($(bar $INV_VG))\$' &&
go test -v -count=1 ./host/surfacepage/ -run '^($(bar $INV_SP)$(alt $NEW_SP))\$' &&
go test -v -count=1 ./host/verbs/ -run '^($(bar $INV_VERBS))\$' &&
go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(bar $INV_CMD))\$'"
need inv $INV_DAEMON $NEW_D $INV_BOUNDARY $INV_STORE $INV_VG $INV_SP $NEW_SP $INV_VERBS $INV_CMD
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.txt"; then FAILS="$FAILS inv_skip_or_fail"; fi
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.txt")" | tee -a "$L/legs.txt"
# live Chrome as V2b ran it (temp profiles; the passkey driver; skips are failures here)
XB="TestApprovePageInChrome TestApprovePageRefusesASwappedRequest TestApprovePageSSHSIGInChrome TestApprovePageShowsAHostileValueAsText TestPasskeyPageAbortsAPendingCeremony TestPasskeyPageInChrome TestApproveDecideWebAuthnLinkInChrome TestDecidePageApproveWithPasskey TestDecidePageFailsClosed TestDecidePageSSHSIGPaste TestNeedsYouListsPendingRequests TestNeedsYouNotifiesOnlyWhenAskedOnce TestVerbBarEpisodeReportInChrome TestVerbBarWhyInChrome"
NEW_X="TestConfirmPageShowsTheFinishPromptInChrome"; [ "$MODE" = base ] && NEW_X=""
XD="TestSurfaceLiveChrome TestSurfaceFollowChrome"
CH="'/Applications/Google Chrome.app/Contents/MacOS/Google Chrome'"
leg chrome "set -o pipefail; WORLD_CHROME=$CH WORLD_REQUIRE_PASSKEY_BROWSER=1 WORLD_REQUIRE_SSH_KEYGEN=1 go test ./cmd/ailang-worldd/ -run '^($(bar $XB)$(alt $NEW_X))\$' -count=1 -v && WORLD_CHROME=$CH go test ./host/daemon/ -run '^($(bar $XD))\$' -count=1 -v"
need chrome $XB $NEW_X $XD
if grep -E '^\s*--- (SKIP|FAIL)' "$L/chrome.txt"; then FAILS="$FAILS chrome_skip_or_fail"; fi
if [ "${GATES_SKIP_VERIFY_GO:-0}" = 1 ]; then echo "LEG verify_go SKIPPED (GATES_SKIP_VERIFY_GO=1)" | tee -a "$L/legs.txt"
else leg verify_go 'VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh'; fi
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" | tee -a "$L/legs.txt"
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" > "$L/ALL_DONE.$SHA"
````

---

## 7. Executor directive notes

- **No git writes of any kind** (no add, rm, mv, commit, stash, checkout, reset, branch). **`patch -p1 -E` (§2) is the
  only writer** of the tree; no hand edit, no `gofmt -w`, no fuzz, **never `WORLD_UPDATE_CONFIRM_GOLDEN=1`** (the golden
  ships in M5a's patch; regenerating it would launder a changed card). A failed `--dry-run` or a red acceptance → **stop
  and report** with the log; do not repair.
- **Never `nohup`, never background a test leg; pass a 600000 ms timeout to every long Bash call.** Each acceptance arm
  fits one call (C14: ≈ 1–4 minutes each).
- **One run, in this fixed order: M0 → M4 → M5a → M5b → M6 → S1b** (each patch assumes the previous tree). After each
  sub-milestone's acceptance, snapshot **cumulatively against the base** to
  `$HOME/.ailang/state/world-iter401/snap/<M>/`: every file added or changed since `07d7c90e` (repo-relative paths), a
  **`DELETED`** file (one repo path per line; S2 deletes nothing, so it is empty), a **`STATUS`** file, and
  `MANIFEST.sha256` (the two plan files, `STATUS` and `DELETED` are not in it). **Never an `rsync` of whole
  directories.** The `find` takes only the top-level dirs that exist (M0–M5b have only `host/`; M6 and S1b `cmd/` and
  `host/`; none has `docs/`), rehearsed on all six (C15):
  ```bash
  bash <<'EOF2'
  set -euo pipefail
  M=M0   # the sub-milestone just accepted
  B=07d7c90e5979893154d2244037bedeffd42386fe
  W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter401-row243-v3s2
  D="$HOME/.ailang/state/world-iter401/snap/$M"; rm -rf "$D"; mkdir -p "$D"
  cd "$W"
  { git diff --name-status --no-renames "$B" -- cmd host docs; git ls-files --others --exclude-standard -- cmd host docs | sed 's/^/A\t/'; } | sort -u -t "$(printf '\t')" -k2 > "$D/STATUS"
  awk -F'\t' '$1 != "D" {print $2}' "$D/STATUS" | while read -r p; do mkdir -p "$D/$(dirname "$p")"; cp -p "$p" "$D/$p"; done
  awk -F'\t' '$1 == "D" {print $2}' "$D/STATUS" > "$D/DELETED"
  (cd "$D" && find $(for d in cmd host docs; do [ -d "$d" ] && echo "$d"; done) -type f -print0 | sort -z | xargs -0 shasum -a 256) > "$D/MANIFEST.sha256"
  cat "$D/STATUS"; echo "DELETED: $(wc -l < "$D/DELETED")"; echo "MANIFEST: $(wc -l < "$D/MANIFEST.sha256")"
  EOF2
  ```
  Expected (`STATUS` lines / `DELETED` / MANIFEST): **M0** 2 / 0 / 2; **M4** 4 / 0 / 4; **M5a** 10 / 0 / 10; **M5b**
  11 / 0 / 11; **M6** 12 / 0 / 12; **S1b** 14 / 0 / 14. The controller copies each snapshot over the worktree, `git rm`s
  each `DELETED` path (none), commits it as one commit, and checks the numstat (`git diff -M --numstat 07d7c90e HEAD --
  cmd host docs | awk '{s+=$1+$2} END {print s}'`): **M0 68, M4 293, M5a 711, M5b 942, M6 1,051, S1b 1,380**. A
  different number → stop.
- `export AILANG_BIN=$HOME/.pinned-ailang/ailang` and `WORLD_EXEC_NODE=/opt/homebrew/bin/node` in every shell that runs Go
  tests (`accept.sh` and `gates.sh` set both). The controller's shell is zsh: run each block under `bash <<'EOF2' …
  EOF2`.
- **Stop conditions**: a dry-run failure; any `--- FAIL`, `--- SKIP` or `FAIL` line in an acceptance (`ACCEPT_RED`); gofmt
  output; a numstat mismatch; any file outside §5 changed (`git status --short` after each patch must equal the
  cumulative list of §3's paragraphs); a numstat above 1,400.
- **Known load-sensitive tests** (the directive's list) are not this sprint's: a red there is re-run once alone. XC1 and
  the other Chrome tests are load-sensitive under heavy parallel load (V2b C11); re-run once alone before stopping.
- D-WORLD-109: nothing touches `~/.ailang/world`, `~/.ssh`, a real `/dev/tty` or a real key; Chrome runs headless on a
  temporary profile against loopback `httptest` servers; S1b's stores are `t.TempDir()` files.

---

## 8. Risks

- **R1 — the 20-line margin.** Any executor-side growth breaks the cap; the plan forbids edits (§7). If a gate forces a
  change, S1b leaves first (back to its own sprint; 1,051 then), never a CP test.
- **R2 — the stand-in DOM is not a browser.** It records sinks and models `hidden`, `textContent`, `append` and ids
  from `index.html`; layout and CSP are XC1's (one Chrome test). A behaviour only a browser shows would be class X
  (one row, `H-SCRIPT-SRC`).
- **R3 — CI has no Chrome** (as for every A3/A4/V2 card test): XC1 runs in the controller's `chrome` leg only. CP1–CP6
  run in CI under node.
- **R4 — S3 must delete `TestConfirmRoutesNotRegistered`** when it registers the routes (Appendix C); otherwise S3 is
  red at its first registration, loudly.
- **R6 — `host/verifygate` under load.** The planner's final `verify_go` hit the package's 10-minute test timeout at
  load ≈ 20–29 (C17); S2 changes nothing it reads (no `.ail`, no doc, no gate input), and iteration 400's controller
  run was red there too. The controller re-runs it alone before judging a red (the directive's list).
- **R5 — the golden is a contract with S3.** S3's `GET /surface/confirmations/{id}` must emit `ConfirmationCard`'s
  output unchanged; its CR1 should compare the served bytes with `ConfirmationCard` of the row.

---

## 9. Acceptance criteria (the judge scores against these)

1. **No route registered**: `daemon.go` unchanged; `TestConfirmRoutesNotRegistered` `--- PASS` (Handler body has no
   `/surface/confirm`; the three paths → 404 through the real Handler; control 200), and `TestApprovalRoutesRegistered`,
   `TestDecisionRoutesRegistered`, `TestSurfaceAppRoutes`, `TestFrozenV1RouteTableMatchesMux` unchanged and green.
2. **No store read or write; no credentials schema change**: no `host/store` file changes; `ConfirmationCard` takes
   primitives (`TestConfirmationCardWire` pins `ConfirmationInput`'s fields); the page's data in tests is a fixture.
3. **No daemon change, `surfaceCSP` unchanged**: no production file under `host/daemon` changes; `TestDecidePageAssets`
   (which pins `surfaceCSP`'s literal) green; XC1 runs the page under the daemon's own headers with zero CSP violations.
4. **The CLI's output unchanged**: no production file under `cmd/` changes; `TestFinishPromptIsTheGolden`,
   `TestEveryVerbHasHelp`, `TestVerbTableMatchesTheCLISwitch` green.
5. **No `.ail`**: `verify_ail` identical to the base's.
6. **CP0**: `ConfirmationCard` emits design §2.5's 13 keys in order, the bytes and the raw challenge, never the message or
   its hash, never an artifact; the state from the row (`TestConfirmationCard`, `TestConfirmationCardWire`).
7. **CP1–CP6** under node, non-skippable with `WORLD_EXEC_NODE`: `TestConfirmPageFailsClosed`, `…RendersText` (incl.
   the NUL rows), `…Message` (CP4 against the **Go-generated** golden; CP6's blob and signing line), `…States`;
   `TestConfirmPageFieldTable` (CP2's Go half).
8. **XC1**: `TestConfirmPageShowsTheFinishPromptInChrome` `--- PASS` in the `chrome` leg; the prompt shown equals
   `broker.FinishPromptText` of the fixture (control: equals `finishPrompt`'s bytes before its question).
9. **VB8, VP12**: `TestAdmitWhy` (the quoting row) and `TestAddressHintsArePinned` green; VB8-JOIN killed; the VP12
   mutants killed (class S, §4).
10. **S1b**: the three judges' tests green; VJ1's mutants as §4 records.
11. **Exactly these pre-existing test files change**: `host/verbs/why_admit_test.go` (+1 row: VB8, design §6/§9(a)).
    Every other test file in §5 is new.
12. **Size**: `git diff -M --numstat 07d7c90e HEAD -- cmd host docs` sums to **1,380** (≤ 1,400); each sub-milestone's
    numstat as §7.
13. **Mutants**: every `K` row killed by its named test; every `CTL` green; §4's equivalents named with reasons; the
    plan's ids are the ledger's.
14. **Gates**: §6's script `FAILS: none` (or only a directive-listed load-sensitive test, green alone), `verify_go` rc 0.

---

## Appendix C — after this sprint (S2's close; what S3, S4, S5 and S6 inherit)

- **Moved out: nothing.** S1b is folded (S2-D1); XC1 stays (no cut fired). **S1b's state: landed with S2** (iteration
  400's patch + three VJ1 test lines, VJ1's mutants run, `inv` names its three tests).
- **Residuals carried, unchanged**: S1 Appendix C's audit notes — note 4 and note 5 / gemini r2 (the ≈ 1 s window after
  C15) → **S5**; notes 1–3 (the "four GETs", CR numbering, the CX15→CX17 gap) → **S3**; the raced supersession row →
  **S3/S4**; design RV1–RV9; VT2 `FinishConfirmRefusal` F1–F5 → **S3's M8**. F6/`DraftOnly` are never built (A).
- **What S3 (the store and the reads) inherits**: `surfacepage.ConfirmationInput` (fill it from a `verb_requests` row:
  `outcome` → `Outcome`, `approved_by`, `decision_ref` → `Decision`, the clock → `Now`) and `ConfirmationCard` for `GET
  /surface/confirmations/{id}` (its CR1 can compare the served bytes with `ConfirmationCard` of the row; the golden's
  shape is the contract); `ConfirmPage`/`ConfirmScript` for `GET /surface/confirm/{id}` and `GET /surface/confirm.js`
  (served with `setSurfaceHeaders`, as `surface_decide_page.go`); **S3 deletes
  `host/daemon/surface_confirm_unrouted_test.go`** in the commit that registers the routes, and adds the confirm page's
  bytes to an assets test (`TestDecidePageAssets`-style, PD2). The `requests.json` *Needs you* list is S3/S6's.
- **What S4 (the submit half) inherits**: `#controls` as the one place for the passkey button and the SSHSIG paste box;
  `net()` (add the POST and A2's 240 s ceremony abort, as `decide.js`); the message `m` and its hash `h` already built
  in `load()` (the passkey's `challenge` is `h`); `CLOSED` for the post-confirm re-read; `passkeys`/`ssh` already in
  the data (`Confirmation.Passkeys`, `.SSH`). The node harness takes new sections (CP8–CP10) the same way.
- **What S5/S6 inherit**: from S1b, `JudgeVerbSSH`/`JudgeVerbAssertion` (S4's C11) and `RecheckVerb` (S5's CF6); the
  card's URL shape `/surface/confirm/vr-<16 hex>` (the bar's `approve vr-` navigation, CB2) and XC1's rig as XC2/XC3's
  starting point.

---

## Appendix P — the patches (apply with §2, in order; generated from the compiled scratch stages; do not edit)

### Patch M0 — the V2b pins (design §9, §6 VB8, VP12)

````diff M0
diff --git a/host/surfacepage/app_hints_test.go b/host/surfacepage/app_hints_test.go
new file mode 100644
index 0000000..0790858
--- /dev/null
+++ b/host/surfacepage/app_hints_test.go
@@ -0,0 +1,67 @@
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
+// TestAddressHintsArePinned (row 243 V3 sprint 2 M0, design §9 (b), §6 VP12): the address bar's three
+// hints. node evaluates app.js's own placeholder expression (each branch) and its KEYROWS literal, cut
+// from the script's source; Go reads index.html's placeholder. It skips only when no node exists; with
+// WORLD_EXEC_NODE set, a harness that cannot run fails.
+func TestAddressHintsArePinned(t *testing.T) {
+	const placeholder = `placeholder="row-189@6 · shared@0ffbb9c · back · episode report &lt;episode&gt; · why &lt;target&gt;">`
+	if n := strings.Count(string(appIndex), placeholder); n != 1 || !strings.Contains(string(AppShell), placeholder) {
+		t.Fatalf("index.html's address placeholder: %d of %q", n, placeholder)
+	}
+	node := os.Getenv("WORLD_EXEC_NODE")
+	if node == "" {
+		var err error
+		node, err = exec.LookPath("node")
+		if err != nil {
+			t.Skip("node unavailable locally")
+		}
+	}
+	dir := t.TempDir()
+	for name, data := range map[string][]byte{"app.js": AppScript, "harness.js": []byte(hintsHarness)} {
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
+	if !strings.Contains(string(out), "ADDRESS HINTS PASS") {
+		t.Fatalf("harness manifest missing: %s", out)
+	}
+}
+
+const hintsHarness = `
+const assert = require('node:assert/strict');
+const fs = require('node:fs');
+const vm = require('node:vm');
+const src = fs.readFileSync('app.js', 'utf8');
+const one = (re) => { const m = src.match(re); assert(m && src.split(m[0]).length === 2, String(re)); return m[1]; };
+// The placeholder: app.js's own expression, evaluated in each of its three branches.
+const expr = one(/\n *\$\('addr'\)\.placeholder = (.*);\n/);
+const at = (st, T) => vm.runInNewContext(expr, { st, T });
+const tail = ' · episode report <episode> · why <target>';
+assert.equal(at({ d: 1 }, null), '');
+assert.equal(at({ d: 0 }, null), 'go to a place — row-189@6 · shared@0ffbb9c' + tail);
+assert.equal(at({ d: 0 }, { GROUPS: [{ id: 'row' }, { id: 'shared' }] }), 'go to a place — row · shared · row:row-N@1' + tail);
+// The keys popover: app.js's own KEYROWS literal; the bar's row names every form the bar takes.
+const rows = JSON.parse(JSON.stringify(vm.runInNewContext(one(/\n *const KEYROWS = (\[.*\]);\n/)))); // out of the vm realm
+assert.equal(rows.length, 13);
+assert.deepEqual(rows[10], ['Enter', 'in the bar: a place, episode report E, why T, approve ID', 'the same']);
+console.log('ADDRESS HINTS PASS');
+`
diff --git a/host/verbs/why_admit_test.go b/host/verbs/why_admit_test.go
index 3e0ac29..ed0f98b 100644
--- a/host/verbs/why_admit_test.go
+++ b/host/verbs/why_admit_test.go
@@ -83,6 +83,7 @@ func TestAdmitWhy(t *testing.T) {
 		{"-", "", 0, false, "why - is not run from the app: the app has no stdin" + at + "-"},
 		{"- --scan 3", "", 0, false, "why - is not run from the app: the app has no stdin" + at + "- --scan 3"},
 		{"--result f", "", 0, false, "--result is not run from the app: the app reads no file on the daemon's host" + at + "--result f"},
+		{"--result 'a b'", "", 0, false, "--result is not run from the app: the app reads no file on the daemon's host" + at + "--result 'a b'"}, // VB8: Quote's form, not a plain join
 		{"--expect-store p head", "", 0, false, "--expect-store is not run from the app: the app reads its own daemon's store" + at + "--expect-store p head"},
 		{"- --result f --expect-store p", "", 0, false, "--result is not run from the app: the app reads no file on the daemon's host" + at + "- --result f --expect-store p"},
 		{"- --expect-store p", "", 0, false, "--expect-store is not run from the app: the app reads its own daemon's store" + at + "- --expect-store p"},
````

### Patch M4 — `ConfirmationCard` (design §2.5, §6 CP0)

````diff M4
diff --git a/host/surfacepage/confirm_model.go b/host/surfacepage/confirm_model.go
new file mode 100644
index 0000000..ed52fc5
--- /dev/null
+++ b/host/surfacepage/confirm_model.go
@@ -0,0 +1,72 @@
+package surfacepage
+
+// Row 243 V3 sprint 2 M4 (design_docs/planned/w-world-address-bar-verbs-v3.md §2.5, §6 CP0): the
+// confirmation card's data, the answer of GET /surface/confirmations/{id} (sprint 3 registers that
+// route and reads the row), built here from primitives without a daemon or a store: the stored
+// canonical bytes and the raw challenge, from which the page computes the digest and the 86-byte
+// message itself. There is no message-hash field and no artifact field.
+
+import (
+	"encoding/base64"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+// The states a confirmation card shows (design §2.5); only ConfirmationPending offers controls.
+const (
+	ConfirmationPending    = "pending"
+	ConfirmationApproved   = "approved"
+	ConfirmationSuperseded = "superseded"
+	ConfirmationExpired    = "expired"
+)
+
+// ConfirmationInput is one stored verb request as a read route reads it: primitives only.
+type ConfirmationInput struct {
+	ID                   string // vr-<16 hex>
+	Request              []byte // the stored canonical world/verb-request/v1 bytes
+	Digest, Challenge    string // 64 hex each, as stored
+	Created, Expires     int64
+	Now                  int64
+	Outcome              string // "" while unconsumed, else ConfirmationApproved or ConfirmationSuperseded
+	ApprovedBy, Decision string // the approver and the world decision's ref, for an approved row
+}
+
+// Confirmation is GET /surface/confirmations/{id} (design §2.5).
+type Confirmation struct {
+	V          int            `json:"v"`
+	ID         string         `json:"id"`
+	Request    string         `json:"request"` // base64url of the stored canonical bytes
+	Digest     string         `json:"digest"`
+	Challenge  string         `json:"challenge"`
+	Created    int64          `json:"created"`
+	Expires    int64          `json:"expires"`
+	Now        int64          `json:"now"`
+	State      string         `json:"state"`
+	ApprovedBy string         `json:"approvedBy"`
+	Decision   string         `json:"decision"`
+	Passkeys   []CardPasskeys `json:"passkeys"`
+	SSH        CardSSH        `json:"ssh"`
+}
+
+// ConfirmationCard is in's card data, its active passkeys grouped by origin as RequestCard groups them.
+func ConfirmationCard(in ConfirmationInput, passkeys []CardPasskey) Confirmation {
+	c := Confirmation{V: 1, ID: in.ID, Request: base64.RawURLEncoding.EncodeToString(in.Request), Digest: in.Digest, Challenge: in.Challenge,
+		Created: in.Created, Expires: in.Expires, Now: in.Now, State: confirmationState(in),
+		Passkeys: RequestCard(StoredRequest{}, passkeys).Passkeys, SSH: CardSSH{Namespace: approval.VerbNamespace}}
+	if c.State == ConfirmationApproved {
+		c.ApprovedBy, c.Decision = in.ApprovedBy, in.Decision
+	}
+	return c
+}
+
+// confirmationState is the row's outcome if it has one, else expired at or after its expiry (as
+// approval.Admit), else pending.
+func confirmationState(in ConfirmationInput) string {
+	if in.Outcome != "" {
+		return in.Outcome
+	}
+	if in.Now >= in.Expires {
+		return ConfirmationExpired
+	}
+	return ConfirmationPending
+}
diff --git a/host/surfacepage/confirm_model_test.go b/host/surfacepage/confirm_model_test.go
new file mode 100644
index 0000000..5f3253c
--- /dev/null
+++ b/host/surfacepage/confirm_model_test.go
@@ -0,0 +1,153 @@
+package surfacepage
+
+// Row 243 V3 sprint 2 M4 (design §2.5, §6 CP0): the confirmation card's data, over fixtures made by
+// the real tools (approval.CanonicalVerb, broker.FinishPromptText, broker.FormatFinishApprovalScope).
+
+import (
+	"crypto/sha256"
+	"encoding/base64"
+	"encoding/hex"
+	"encoding/json"
+	"reflect"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/broker"
+)
+
+// confirmFixture is one stored verb request: the request, its canonical bytes and the row.
+type confirmFixture struct {
+	name string
+	req  approval.VerbRequest
+	in   ConfirmationInput
+}
+
+// confirmFixtures are a draft finish with protected paths and the repository's committer, and a ready
+// one creating its label whose title, body and argv are hostile (markup, quotes, a NUL).
+func confirmFixtures(t *testing.T) []confirmFixture {
+	t.Helper()
+	const created = 1800000000
+	var out []confirmFixture
+	for i, c := range []struct {
+		name, id, title, body, committer string
+		argv, protected                  []string
+		draft, createLabel               bool
+	}{
+		{"draft", "vr-0123456789abcdef", "Fix the bar", "## Summary\nline 2\nline 3\nline 4\nline 5\nline 6\nline 7 is not in the prompt\n", "",
+			[]string{"--preset", "p", "--title", "Fix the bar"}, []string{"host/x.go", "docs/a b.md"}, true, false},
+		{"hostile", "vr-00000000000000ff", "<img src=x onerror=alert(1)> a\x00b & 'q'", "</pre><script>alert(2)</script>\nNUL \x00 here &amp; after — ünï\n", "Mark <m@x>",
+			[]string{"--preset", "p", "--title", "it's <b>", "--ready"}, []string{}, false, true},
+	} {
+		scope := broker.FinishApprovalScope{Repository: "sunholo-data/ailang-world", Remote: "origin", Head: "world/w1", SHA: strings.Repeat("ab", 20),
+			Base: "dev", TitleSHA256: strings.Repeat("0", 64), BodySHA256: strings.Repeat("1", 64), Draft: c.draft, CreateLabel: c.createLabel, ExpiresAt: created + 900}
+		r := approval.VerbRequest{Schema: approval.VerbRequestSchema, Verb: approval.VerbFinish, Argv: append([]string{"ailang-worldd", "episode", "finish", "w1"}, c.argv...),
+			Episode: "w1", Prompt: broker.FinishPromptText(scope, c.title, c.body, c.protected, c.committer), Scope: broker.FormatFinishApprovalScope(scope),
+			Title: c.title, Body: c.body, Committer: c.committer, Protected: c.protected}
+		b, err := approval.CanonicalVerb(r)
+		if err == nil {
+			err = r.Validate() // the fixtures are admissible requests
+		}
+		if err != nil {
+			t.Fatalf("%s: %v", c.name, err)
+		}
+		d, ch := sha256.Sum256(b), sha256.Sum256([]byte("challenge "+c.name))
+		out = append(out, confirmFixture{c.name, r, ConfirmationInput{ID: c.id, Request: b, Digest: hex.EncodeToString(d[:]), Challenge: hex.EncodeToString(ch[:]),
+			Created: created + int64(i), Expires: created + int64(i) + 600, Now: created + 5}})
+	}
+	return out
+}
+
+// CP0: the card carries the stored bytes (base64url), the digest and the raw challenge as stored, the
+// times, the state from the row (its outcome, else expired from its expiry as Admit, else pending), the
+// approver and decision only for an approved row, the passkeys grouped by origin and the verb namespace.
+func TestConfirmationCard(t *testing.T) {
+	f := confirmFixtures(t)[0]
+	pk := []CardPasskey{{Origin: "http://localhost:1", RPID: "localhost", CredentialID: []byte{1}}, {Origin: "http://127.0.0.1:1", RPID: "127.0.0.1", CredentialID: []byte{2}},
+		{Origin: "http://localhost:1", RPID: "localhost", CredentialID: []byte{3}}}
+	c := ConfirmationCard(f.in, pk)
+	if b, err := base64.RawURLEncoding.DecodeString(c.Request); err != nil || string(b) != string(f.in.Request) {
+		t.Fatalf("request %q does not carry the stored bytes (%v)", c.Request, err)
+	}
+	if u := ConfirmationCard(ConfirmationInput{Request: []byte{0xfb, 0xff, 0xfe}}, nil).Request; u != "-__-" { // base64url, unpadded ("+//+" in base64)
+		t.Fatalf("request %q is not base64url", u)
+	}
+	if c.V != 1 || c.ID != f.in.ID || c.Digest != f.in.Digest || c.Challenge != f.in.Challenge || c.Created != f.in.Created || c.Expires != f.in.Expires || c.Now != f.in.Now {
+		t.Fatalf("card %+v, input %+v", c, f.in)
+	}
+	if want := []CardPasskeys{{Origin: "http://localhost:1", RP: "localhost", Allow: []string{"AQ", "Aw"}}, {Origin: "http://127.0.0.1:1", RP: "127.0.0.1", Allow: []string{"Ag"}}}; !reflect.DeepEqual(c.Passkeys, want) {
+		t.Fatalf("passkeys %+v", c.Passkeys)
+	}
+	if c.SSH.Namespace != "world-verb@sunholo" { // approval.VerbNamespace, as a literal
+		t.Fatalf("namespace %q", c.SSH.Namespace)
+	}
+	if none := ConfirmationCard(f.in, nil); none.Passkeys == nil || len(none.Passkeys) != 0 {
+		t.Fatalf("no passkeys: %#v", none.Passkeys)
+	}
+	for _, tc := range []struct {
+		name, outcome     string
+		now               int64
+		state, by, decide string
+	}{
+		{"pending", "", f.in.Expires - 1, ConfirmationPending, "", ""},
+		{"expired at its expiry", "", f.in.Expires, ConfirmationExpired, "", ""},
+		{"expired after", "", f.in.Expires + 100, ConfirmationExpired, "", ""},
+		{"approved", ConfirmationApproved, f.in.Created + 1, ConfirmationApproved, "webauthn:mark", "sha256:" + strings.Repeat("cd", 32)},
+		{"approved, read after its expiry", ConfirmationApproved, f.in.Expires + 1, ConfirmationApproved, "webauthn:mark", "sha256:" + strings.Repeat("cd", 32)},
+		{"superseded", ConfirmationSuperseded, f.in.Created + 1, ConfirmationSuperseded, "", ""},
+	} {
+		in := f.in
+		in.Outcome, in.Now, in.ApprovedBy, in.Decision = tc.outcome, tc.now, "webauthn:mark", "sha256:"+strings.Repeat("cd", 32) // set on every row: only approved shows them
+		c := ConfirmationCard(in, nil)
+		if c.State != tc.state || c.ApprovedBy != tc.by || c.Decision != tc.decide || c.Now != tc.now {
+			t.Errorf("%s: state %q by %q decision %q now %d; want %q %q %q %d", tc.name, c.State, c.ApprovedBy, c.Decision, c.Now, tc.state, tc.by, tc.decide, tc.now)
+		}
+	}
+}
+
+// CP0: the wire is exactly design §2.5's keys in order; the input has no artifact; nothing on the wire
+// is the message or its hash in any encoding.
+func TestConfirmationCardWire(t *testing.T) {
+	var names []string
+	for i, ty := 0, reflect.TypeOf(ConfirmationInput{}); i < ty.NumField(); i++ {
+		names = append(names, ty.Field(i).Name)
+	}
+	if want := []string{"ID", "Request", "Digest", "Challenge", "Created", "Expires", "Now", "Outcome", "ApprovedBy", "Decision"}; !reflect.DeepEqual(names, want) {
+		t.Fatalf("ConfirmationInput's fields %v, want %v (no artifact, no hash)", names, want)
+	}
+	for _, f := range confirmFixtures(t) {
+		in := f.in
+		in.Outcome, in.ApprovedBy, in.Decision = ConfirmationApproved, "ssh:mark", "sha256:"+strings.Repeat("ef", 32)
+		b, err := json.Marshal(ConfirmationCard(in, []CardPasskey{{Origin: "http://localhost:1", RPID: "localhost", CredentialID: []byte{9}}}))
+		if err != nil {
+			t.Fatal(err)
+		}
+		var keys []string
+		dec := json.NewDecoder(strings.NewReader(string(b)))
+		_, _ = dec.Token()
+		for dec.More() {
+			k, _ := dec.Token()
+			keys = append(keys, k.(string))
+			var skip json.RawMessage
+			_ = dec.Decode(&skip)
+		}
+		if want := []string{"v", "id", "request", "digest", "challenge", "created", "expires", "now", "state", "approvedBy", "decision", "passkeys", "ssh"}; !reflect.DeepEqual(keys, want) {
+			t.Fatalf("%s: keys %v, want %v", f.name, keys, want)
+		}
+		var d, ch [32]byte
+		hex.Decode(d[:], []byte(in.Digest))
+		hex.Decode(ch[:], []byte(in.Challenge))
+		m, _ := approval.VerbMessage(d, ch, in.Expires)
+		h, _ := approval.VerbMessageHash(d, ch, in.Expires)
+		for _, x := range [][]byte{m, h[:]} {
+			for _, enc := range []string{hex.EncodeToString(x), base64.RawURLEncoding.EncodeToString(x), base64.StdEncoding.EncodeToString(x)} {
+				if strings.Contains(string(b), enc) {
+					t.Fatalf("%s: the wire carries the message or its hash (%s)", f.name, enc)
+				}
+			}
+		}
+		if !strings.Contains(string(b), `"state":"approved","approvedBy":"ssh:mark","decision":"sha256:`) || !strings.Contains(string(b), `"ssh":{"namespace":"world-verb@sunholo"}`) {
+			t.Fatalf("%s: %s", f.name, b)
+		}
+	}
+}
````

### Patch M5a — the page, the script, the table, the golden, no route (design §2.5; CP2 Go half)

````diff M5a
diff --git a/host/daemon/surface_confirm_unrouted_test.go b/host/daemon/surface_confirm_unrouted_test.go
new file mode 100644
index 0000000..5d2d195
--- /dev/null
+++ b/host/daemon/surface_confirm_unrouted_test.go
@@ -0,0 +1,33 @@
+package daemon
+
+// Row 243 V3 sprint 2 (design §2.5 [r1], §8): the confirmation card's page and its data exist only in
+// host/surfacepage and tests; no daemon route serves them before sprint 3, which replaces this test
+// with the routes' own (CR1).
+
+import (
+	"os"
+	"strings"
+	"testing"
+)
+
+// The Handler registers nothing under /surface/confirm, and the daemon answers 404 for the card, its
+// script and its data. Control: the decision card's page, registered beside them, answers 200.
+func TestConfirmRoutesNotRegistered(t *testing.T) {
+	src, err := os.ReadFile("daemon.go")
+	if err != nil {
+		t.Fatal(err)
+	}
+	body := strings.SplitN(strings.SplitN(string(src), "func (d *Daemon) Handler() http.Handler {", 2)[1], "\n}\n", 2)[0]
+	if strings.Contains(body, "/surface/confirm") || !strings.Contains(body, `"GET /surface/decide/{ref}"`) {
+		t.Fatal("Handler registers a /surface/confirm route (or lost the decide control)")
+	}
+	a := newApprovalRig(t)
+	for _, target := range []string{"/surface/confirm/vr-0123456789abcdef", "/surface/confirm.js", "/surface/confirmations/vr-0123456789abcdef"} {
+		if rec := a.serve("GET", target, apLoopback, nil); rec.Code != 404 {
+			t.Errorf("GET %s = %d, want 404 (no route before sprint 3)", target, rec.Code)
+		}
+	}
+	if rec := a.serve("GET", "/surface/decide/sha256:"+strings.Repeat("ab", 32), apLoopback, nil); rec.Code != 200 {
+		t.Fatalf("control: GET /surface/decide/… = %d", rec.Code)
+	}
+}
diff --git a/host/surfacepage/confirm.go b/host/surfacepage/confirm.go
new file mode 100644
index 0000000..112324b
--- /dev/null
+++ b/host/surfacepage/confirm.go
@@ -0,0 +1,21 @@
+package surfacepage
+
+import _ "embed"
+
+// Row 243 V3 sprint 2 M5 (design_docs/planned/w-world-address-bar-verbs-v3.md §2.5): the confirmation
+// card, to be served at GET /surface/confirm/{id} and GET /surface/confirm.js (sprint 3 registers the
+// routes; until then only tests serve these bytes). The page reads the card's data from GET
+// /surface/confirmations/{id}, checks the request's digest, shows every field of the signed bytes
+// through one table (refusing a field it cannot show), and builds the 86-byte message itself: Save
+// the message and the ssh-keygen line. Sprint 4 adds the submit half. The bytes are the same for every
+// id.
+
+// ConfirmPage is GET /surface/confirm/{id}.
+//
+//go:embed confirm/index.html
+var ConfirmPage []byte
+
+// ConfirmScript is GET /surface/confirm.js.
+//
+//go:embed confirm/confirm.js
+var ConfirmScript []byte
diff --git a/host/surfacepage/confirm/confirm.js b/host/surfacepage/confirm/confirm.js
new file mode 100644
index 0000000..12296d7
--- /dev/null
+++ b/host/surfacepage/confirm/confirm.js
@@ -0,0 +1,167 @@
+'use strict';
+// AILANG World's confirmation card (row 243 V3, design §2.5). It reads one stored verb request from
+// the daemon, checks that its bytes hash to its digest, shows EVERY field of the signed bytes through
+// one table (and refuses a request with a field it cannot show), and builds the 86-byte world-verb/v1
+// message itself (never a hash from the server): Save the message and the ssh-keygen line. A state
+// other than pending shows that state and no controls. Text is set with textContent only. Sprint 4
+// adds the submit half (a passkey and a pasted SSH signature) inside #controls.
+(() => {
+  const NET_MS = 10000; // A3 §2.7 [r2]: the data GET, its body included
+  const PREFIX = 'world-verb/v1', NAMESPACE = 'world-verb@sunholo', LEN = 86; // approval.VerbPrefix, VerbNamespace, VerbMessageLen
+  // Every key of a world/verb-request/v1 (approval.VerbRequest's json tags, in order), with its label
+  // and JSON type. TestConfirmPageFieldTable compares this list with the Go type.
+  const CONFIRM_FIELDS = /* CONFIRM_FIELDS */[
+    {"key": "schema", "label": "Schema", "type": "string"},
+    {"key": "verb", "label": "Verb", "type": "string"},
+    {"key": "argv", "label": "Command words", "type": "strings"},
+    {"key": "episode", "label": "Episode", "type": "string"},
+    {"key": "prompt", "label": "Prompt", "type": "string"},
+    {"key": "scope", "label": "Approval scope", "type": "string"},
+    {"key": "title", "label": "Title", "type": "string"},
+    {"key": "body", "label": "Body", "type": "string"},
+    {"key": "committer", "label": "Committer", "type": "string"},
+    {"key": "protected", "label": "Protected paths changed", "type": "strings"}
+  ]/* /CONFIRM_FIELDS */;
+  const CLOSED = { approved: 'Already confirmed', superseded: 'Superseded by a newer finish of this episode', expired: 'Expired' };
+
+  const $ = (id) => document.getElementById(id);
+  const setState = (s) => { document.body.dataset.state = s; };
+  const refuse = (state, text) => { setState(state); $('what').textContent = text; $('what').toggleAttribute('data-error', true); };
+  const unb64u = (s) => {
+    if (typeof s !== 'string' || !/^[A-Za-z0-9_-]*$/.test(s)) throw new Error('not base64url');
+    return Uint8Array.from(atob(s.replace(/-/g, '+').replace(/_/g, '/')), (c) => c.charCodeAt(0));
+  };
+  const hex = (b) => Array.from(b, (x) => x.toString(16).padStart(2, '0')).join('');
+  const unhex = (s) => (typeof s === 'string' && /^[0-9a-f]{64}$/.test(s) ? Uint8Array.from(s.match(/../g), (x) => parseInt(x, 16)) : null);
+  const sha256 = async (bytes) => new Uint8Array(await crypto.subtle.digest('SHA-256', bytes));
+  const time = (v) => `${new Date(v * 1000).toISOString()} (${v})`;
+  // quote is verbs.Quote: a word outside safeWord's set is single-quoted, a ' inside as '"'"'.
+  const SAFE = /^[A-Za-z0-9@%+=:,./_-]+$/;
+  const quote = (argv) => argv.map((a) => (SAFE.test(a) ? a : "'" + a.split("'").join(`'"'"'`) + "'")).join(' ');
+
+  async function net(url) {
+    const c = new AbortController();
+    const timer = setTimeout(() => c.abort(), NET_MS);
+    try {
+      const r = await fetch(url, { cache: 'no-store', credentials: 'omit', signal: c.signal, headers: { Accept: 'application/json' } });
+      const text = await r.text();
+      let body = null;
+      try { body = JSON.parse(text); } catch (_) { body = null; }
+      return { status: r.status, body };
+    } finally {
+      clearTimeout(timer);
+    }
+  }
+  const why = (e) => (e && e.name === 'AbortError' ? `no answer within ${NET_MS / 1000} s` : String(e && e.message || e));
+
+  // unshown is the first key of request r the card cannot show (unlisted, of another type, or
+  // missing), or ''. A value that is not a JSON object is itself unshowable.
+  const typed = (f, v) => (f.type === 'string' ? typeof v === 'string' : Array.isArray(v) && v.every((x) => typeof x === 'string'));
+  function unshown(r) {
+    if (r === null || typeof r !== 'object' || Array.isArray(r)) return 'the request';
+    for (const k of Object.keys(r)) {
+      const f = CONFIRM_FIELDS.find((x) => x.key === k);
+      if (!f || !typed(f, r[k])) return k;
+    }
+    for (const f of CONFIRM_FIELDS) if (!Object.prototype.hasOwnProperty.call(r, f.key)) return f.key;
+    return '';
+  }
+
+  // message is the 86-byte verb message (design §2.2): PREFIX ‖ 0x00 ‖ digest ‖ challenge ‖ u64be(expiry).
+  function message(digest, challenge, expiry) {
+    const m = new Uint8Array(LEN);
+    m.set(new TextEncoder().encode(PREFIX), 0);
+    m.set(digest, 14);
+    m.set(challenge, 46);
+    new DataView(m.buffer).setBigUint64(78, BigInt(expiry));
+    return m;
+  }
+
+  const row = (dl, id, label, text) => {
+    const dt = document.createElement('dt');
+    dt.textContent = label;
+    const dd = document.createElement('dd');
+    dd.id = id;
+    dd.textContent = text;
+    dl.append(dt, dd);
+    return dd;
+  };
+  function render(r, s, h) {
+    $('argv').textContent = quote(r.argv);
+    $('prompt').textContent = r.prompt;
+    const dl = $('signed');
+    dl.textContent = '';
+    row(dl, 'm-digest', 'Request digest (SHA-256)', s.digest);
+    row(dl, 'm-challenge', 'Challenge', s.challenge);
+    row(dl, 'm-expires', 'Expires', time(s.expires));
+    row(dl, 'm-hash', 'Message hash (the passkey challenge)', hex(h));
+    for (const f of CONFIRM_FIELDS) {
+      const v = r[f.key];
+      if (f.key === 'body') {
+        const dd = row(dl, 'f-body-size', f.label, `${new TextEncoder().encode(v).length} bytes`);
+        const details = document.createElement('details');
+        const summary = document.createElement('summary');
+        summary.textContent = 'the full body';
+        const pre = document.createElement('pre');
+        pre.id = 'f-body';
+        pre.textContent = v;
+        details.append(summary, pre);
+        dd.append(details);
+        continue;
+      }
+      row(dl, 'f-' + f.key, f.label, f.type === 'strings' ? v.join('\n') : v);
+    }
+    $('shown').hidden = false;
+  }
+
+  const id = (location.pathname.match(/^\/surface\/confirm\/(vr-[0-9a-f]{16})$/) || [])[1];
+  async function load() {
+    if (!id) return refuse('refused', 'This is not a confirmation card\'s address; nothing to confirm.');
+    let res;
+    try {
+      res = await net('/surface/confirmations/' + id);
+    } catch (e) {
+      return refuse('failed', `The request could not be read (${why(e)}). Nothing to confirm here; reload, or answer at the terminal.`);
+    }
+    if (res.status !== 200 || !res.body) {
+      return refuse('failed', `The daemon did not give the request (HTTP ${res.status}${res.body && res.body.reason ? ': ' + res.body.reason : ''}). Nothing to confirm here.`);
+    }
+    const s = res.body;
+    let bytes;
+    try {
+      bytes = unb64u(s.request);
+    } catch (_) {
+      return refuse('refused', 'The daemon\'s request is not base64url; nothing to confirm.');
+    }
+    const challenge = unhex(s.challenge);
+    if (s.id !== id || !unhex(s.digest) || !challenge || !Number.isSafeInteger(s.expires) || s.expires < 0 || typeof s.state !== 'string') {
+      return refuse('refused', 'The daemon\'s answer is malformed; nothing to confirm.');
+    }
+    const digest = await sha256(bytes);
+    if (hex(digest) !== s.digest) return refuse('refused', 'The request\'s bytes do not match its digest; nothing to confirm.');
+    let r;
+    try {
+      r = JSON.parse(new TextDecoder('utf-8', { fatal: true }).decode(bytes));
+    } catch (e) {
+      return refuse('refused', `The request is not UTF-8 JSON (${e.message}); nothing to confirm.`);
+    }
+    const bad = unshown(r);
+    if (bad) return refuse('refused', `This request has a field the page cannot show (${bad}); nothing to confirm here; confirm at the terminal.`);
+    const m = message(digest, challenge, s.expires), h = await sha256(m);
+    render(r, s, h);
+    $('what').textContent = `Confirm request ${id}: ${r.verb} ${r.episode}; expires ${time(s.expires)}; message hash ${hex(h)}`;
+    if (s.state !== 'pending') {
+      setState('closed');
+      $('status').textContent = `${CLOSED[s.state] || 'This request is ' + s.state}${s.state === 'approved' ? ` by ${s.approvedBy}: decision ${s.decision}` : ''}; nothing to confirm here.`;
+      return;
+    }
+    $('save').href = URL.createObjectURL(new Blob([m], { type: 'application/octet-stream' }));
+    $('save').download = `world-verb-${id}.msg`;
+    $('sshline').textContent = `ssh-keygen -Y sign -n ${NAMESPACE} -f <key> world-verb-${id}.msg`;
+    $('status').textContent = 'Waiting for your key: save the message and sign it with your enrolled SSH key.';
+    $('controls').hidden = false;
+    setState('ready');
+  }
+
+  load(); // the daemon will serve this page only at /surface/confirm/vr-<16 hex>
+})();
diff --git a/host/surfacepage/confirm/index.html b/host/surfacepage/confirm/index.html
new file mode 100644
index 0000000..dc80518
--- /dev/null
+++ b/host/surfacepage/confirm/index.html
@@ -0,0 +1,40 @@
+<!doctype html>
+<html lang="en">
+<head>
+<meta charset="utf-8">
+<meta name="viewport" content="width=device-width, initial-scale=1">
+<title>AILANG World: confirm an episode finish</title>
+<style>
+body { font: 16px/1.5 system-ui, sans-serif; max-width: 52rem; margin: 2rem auto; padding: 0 1rem; color: #1d1d1f; }
+h1 { font-size: 1.4rem; }
+h2 { font-size: 1.05rem; margin: 1.2rem 0 0.3rem; }
+#what, #status { white-space: pre-wrap; }
+#what[data-error] { color: #b00020; font-weight: 600; }
+pre { font: 13px ui-monospace, monospace; white-space: pre-wrap; overflow-wrap: anywhere; background: #f5f5f7; padding: 0.6rem; margin: 0; }
+dl { display: grid; grid-template-columns: max-content 1fr; gap: 0.15rem 1rem; margin: 0; }
+dt { color: #555; }
+dd { margin: 0; font: 13px ui-monospace, monospace; white-space: pre-wrap; overflow-wrap: anywhere; }
+#sshline { font: 13px ui-monospace, monospace; overflow-wrap: anywhere; }
+[hidden] { display: none !important; }
+</style>
+</head>
+<body data-state="loading">
+<h1>Confirm an episode finish</h1>
+<p id="what">Reading the request…</p>
+<div id="shown" hidden>
+<h2>The command</h2>
+<pre id="argv"></pre>
+<h2>The terminal's prompt</h2>
+<pre id="prompt"></pre>
+<h2>Open this pull request?</h2>
+<p id="status"></p>
+<div id="controls" hidden>
+<p><a id="save">Save the message</a></p>
+<p id="sshline"></p>
+</div>
+<h2>What your key signs</h2>
+<dl id="signed"></dl>
+</div>
+<script src="/surface/confirm.js"></script>
+</body>
+</html>
diff --git a/host/surfacepage/confirm_page_test.go b/host/surfacepage/confirm_page_test.go
new file mode 100644
index 0000000..56d8af3
--- /dev/null
+++ b/host/surfacepage/confirm_page_test.go
@@ -0,0 +1,152 @@
+package surfacepage
+
+// Row 243 V3 sprint 2 M5 (design §2.5, §6 CP2, CP4): the confirmation card's field table, the golden its
+// node tests read (written by this Go test from the real tools, never by hand), and the page's assets.
+
+import (
+	"bytes"
+	"encoding/hex"
+	"encoding/json"
+	"os"
+	"path/filepath"
+	"reflect"
+	"regexp"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+// confirmGolden is the node tests' fixture file; WORLD_UPDATE_CONFIRM_GOLDEN=1 rewrites it.
+var confirmGolden = filepath.Join("testdata", "confirm-card.golden.jsonl")
+
+// confirmFields is confirm.js's CONFIRM_FIELDS table.
+func confirmFields(t *testing.T) []struct{ Key, Label, Type string } {
+	t.Helper()
+	_, rest, _ := bytes.Cut(ConfirmScript, []byte("/* CONFIRM_FIELDS */"))
+	lit, _, _ := bytes.Cut(rest, []byte("/* /CONFIRM_FIELDS */"))
+	var fields []struct{ Key, Label, Type string }
+	if err := json.Unmarshal(lit, &fields); err != nil || len(fields) == 0 {
+		t.Fatalf("confirm.js's CONFIRM_FIELDS: %v", err)
+	}
+	return fields
+}
+
+// CP2: CONFIRM_FIELDS names exactly approval.VerbRequest's json tags, in order, each with its JSON type
+// (string or a list of strings) and a label: a field added to the request without the page fails here.
+func TestConfirmPageFieldTable(t *testing.T) {
+	var want, got []string
+	ty := reflect.TypeOf(approval.VerbRequest{})
+	for i := 0; i < ty.NumField(); i++ {
+		kind := map[reflect.Kind]string{reflect.String: "string", reflect.Slice: "strings"}[ty.Field(i).Type.Kind()]
+		if ty.Field(i).Type.Kind() == reflect.Slice && ty.Field(i).Type.Elem().Kind() != reflect.String {
+			kind = ""
+		}
+		want = append(want, strings.Split(ty.Field(i).Tag.Get("json"), ",")[0]+":"+kind)
+	}
+	for _, f := range confirmFields(t) {
+		if f.Label == "" {
+			t.Errorf("CONFIRM_FIELDS %s has no label", f.Key)
+		}
+		got = append(got, f.Key+":"+f.Type)
+	}
+	if !reflect.DeepEqual(got, want) || len(want) != 10 { // control: design §2.2's ten fields
+		t.Fatalf("CONFIRM_FIELDS %v, approval.VerbRequest %v", got, want)
+	}
+}
+
+// confirmGoldenLine is one card the node tests open, and what the page must make of it.
+type confirmGoldenLine struct {
+	Name     string       `json:"name"`
+	Card     Confirmation `json:"card"`
+	ArgvLine string       `json:"argvLine"` // verbs.Quote(argv)
+	Message  string       `json:"message"`  // hex of approval.VerbMessage(digest, challenge, expires)
+	Hash     string       `json:"hash"`     // hex of approval.VerbMessageHash
+}
+
+// CP4, CP6: the golden is the fixtures' cards (ConfirmationCard over approval.CanonicalVerb's bytes) in
+// each state, with Go's own message, its hash and the quoted argv; the node tests check the page's
+// against it. WORLD_UPDATE_CONFIRM_GOLDEN=1 go test ./host/surfacepage/ -run '^TestConfirmCardGolden$'
+// writes it.
+func TestConfirmCardGolden(t *testing.T) {
+	var out bytes.Buffer
+	fx := confirmFixtures(t)
+	for _, c := range []struct {
+		name    string
+		f       confirmFixture
+		outcome string
+		now     int64
+	}{
+		{"draft", fx[0], "", 0}, {"hostile", fx[1], "", 0}, {"draft-approved", fx[0], ConfirmationApproved, 0},
+		{"draft-superseded", fx[0], ConfirmationSuperseded, 0}, {"draft-expired", fx[0], "", fx[0].in.Expires},
+	} {
+		in := c.f.in
+		if c.outcome != "" {
+			in.Outcome = c.outcome
+		}
+		if c.now != 0 {
+			in.Now = c.now
+		}
+		if c.outcome == ConfirmationApproved {
+			in.ApprovedBy, in.Decision = "webauthn:mark-mac", "sha256:"+strings.Repeat("5e", 32)
+		}
+		var d, ch [32]byte
+		if _, err := hex.Decode(d[:], []byte(in.Digest)); err != nil {
+			t.Fatal(err)
+		}
+		if _, err := hex.Decode(ch[:], []byte(in.Challenge)); err != nil {
+			t.Fatal(err)
+		}
+		m, err := approval.VerbMessage(d, ch, in.Expires)
+		if err != nil {
+			t.Fatal(err)
+		}
+		h, _ := approval.VerbMessageHash(d, ch, in.Expires)
+		line, err := json.Marshal(confirmGoldenLine{Name: c.name, Card: ConfirmationCard(in, []CardPasskey{{Origin: "http://localhost:1", RPID: "localhost", CredentialID: []byte{7, 7}}}),
+			ArgvLine: verbs.Quote(c.f.req.Argv), Message: hex.EncodeToString(m), Hash: hex.EncodeToString(h[:])})
+		if err != nil {
+			t.Fatal(err)
+		}
+		out.Write(append(line, '\n'))
+	}
+	if os.Getenv("WORLD_UPDATE_CONFIRM_GOLDEN") == "1" {
+		if err := os.WriteFile(confirmGolden, out.Bytes(), 0o644); err != nil {
+			t.Fatal(err)
+		}
+	}
+	got, err := os.ReadFile(confirmGolden)
+	if err != nil || !bytes.Equal(got, out.Bytes()) {
+		t.Fatalf("%s is not the cards the fixtures make (%v); regenerate it with WORLD_UPDATE_CONFIRM_GOLDEN=1", confirmGolden, err)
+	}
+}
+
+// The page loads exactly its own script and has no inline handler, form, image, link, frame or remote
+// URL; the ids confirm.js reads exist; the embedded bytes are the files'. Labelled instrument-health
+// controls (not load-bearing; the node tests and XC1 are): the script uses no HTML sink, eval, cookie
+// or other transport.
+func TestConfirmPageAssets(t *testing.T) {
+	for file, b := range map[string][]byte{"index.html": ConfirmPage, "confirm.js": ConfirmScript} {
+		want, err := os.ReadFile(filepath.Join("confirm", file))
+		if err != nil || len(want) == 0 || !bytes.Equal(b, want) {
+			t.Fatalf("confirm/%s is not embedded (%v)", file, err)
+		}
+	}
+	page := string(ConfirmPage)
+	if scripts := regexp.MustCompile(`(?i)<script\b[^>]*>`).FindAllString(page, -1); len(scripts) != 1 || !strings.Contains(page, `<script src="/surface/confirm.js"></script>`) {
+		t.Errorf("scripts %q", scripts)
+	}
+	if m := regexp.MustCompile(`(?i)\son[a-z]+\s*=|<form|<img|<link|<iframe|https?://`).FindString(page); m != "" {
+		t.Errorf("the page has %q", m)
+	}
+	for _, m := range regexp.MustCompile(`\$\('([a-z-]+)'\)`).FindAllStringSubmatch(string(ConfirmScript), -1) {
+		if !strings.Contains(page, ` id="`+m[1]+`"`) {
+			t.Errorf("confirm.js reads #%s, which the page lacks", m[1])
+		}
+	}
+	for _, banned := range []string{"innerHTML", "outerHTML", "insertAdjacentHTML", "document.write", "eval(", "new Function", "document.cookie", "XMLHttpRequest", "sendBeacon", "WebSocket"} {
+		if bytes.Contains(ConfirmScript, []byte(banned)) {
+			t.Errorf("confirm.js uses %s", banned)
+		}
+	}
+}
diff --git a/host/surfacepage/testdata/confirm-card.golden.jsonl b/host/surfacepage/testdata/confirm-card.golden.jsonl
new file mode 100644
index 0000000..fd5b1f5
--- /dev/null
+++ b/host/surfacepage/testdata/confirm-card.golden.jsonl
@@ -0,0 +1,5 @@
+{"name":"draft","card":{"v":1,"id":"vr-0123456789abcdef","request":"eyJzY2hlbWEiOiJ3b3JsZC92ZXJiLXJlcXVlc3QvdjEiLCJ2ZXJiIjoiZXBpc29kZSBmaW5pc2giLCJhcmd2IjpbImFpbGFuZy13b3JsZGQiLCJlcGlzb2RlIiwiZmluaXNoIiwidzEiLCItLXByZXNldCIsInAiLCItLXRpdGxlIiwiRml4IHRoZSBiYXIiXSwiZXBpc29kZSI6IncxIiwicHJvbXB0IjoiY29tbWl0dGVyOiAocmVwb3NpdG9yeSBnaXQgY29uZmlnKVxucmVwb3NpdG9yeTogc3VuaG9sby1kYXRhL2FpbGFuZy13b3JsZFxuYnJhbmNoOiB3b3JsZC93MVxuc2hhOiBhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiXG5iYXNlOiBkZXZcbnRpdGxlOiBGaXggdGhlIGJhclxuYm9keSAoNzQgYnl0ZXMpLCBmaXJzdCBsaW5lczpcbiMjIFN1bW1hcnlcbmxpbmUgMlxubGluZSAzXG5saW5lIDRcbmxpbmUgNVxubGluZSA2XG5sYWJlbDogYWlsYW5nLXdvcmxkXG5kcmFmdFxucHJvdGVjdGVkIHBhdGhzIGNoYW5nZWQ6IGhvc3QveC5nbywgZG9jcy9hIGIubWRcbiIsInNjb3BlIjoiZ2l0aHViLXByOnN1bmhvbG8tZGF0YS9haWxhbmctd29ybGQjcmVtb3RlPW9yaWdpbiZoZWFkPXdvcmxkL3cxJnNoYT1hYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiJmJhc2U9ZGV2JnRpdGxlLXNoYTI1Nj0wMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwJmJvZHktc2hhMjU2PTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTEmbGFiZWw9YWlsYW5nLXdvcmxkJmRyYWZ0PXRydWUmY3JlYXRlLWxhYmVsPWZhbHNlJmV4cGlyZXM9MTgwMDAwMDkwMCIsInRpdGxlIjoiRml4IHRoZSBiYXIiLCJib2R5IjoiIyMgU3VtbWFyeVxubGluZSAyXG5saW5lIDNcbmxpbmUgNFxubGluZSA1XG5saW5lIDZcbmxpbmUgNyBpcyBub3QgaW4gdGhlIHByb21wdFxuIiwiY29tbWl0dGVyIjoiIiwicHJvdGVjdGVkIjpbImhvc3QveC5nbyIsImRvY3MvYSBiLm1kIl19","digest":"68e53444f32639d1676d16f5505f632cfae431352b3a68405a67ac3ec04b3f17","challenge":"5eddea8e87657285501aed0a6e790a046cbb45d190713679cc31e38ac2db2d05","created":1800000000,"expires":1800000600,"now":1800000005,"state":"pending","approvedBy":"","decision":"","passkeys":[{"origin":"http://localhost:1","rp":"localhost","allow":["Bwc"]}],"ssh":{"namespace":"world-verb@sunholo"}},"argvLine":"ailang-worldd episode finish w1 --preset p --title 'Fix the bar'","message":"776f726c642d766572622f76310068e53444f32639d1676d16f5505f632cfae431352b3a68405a67ac3ec04b3f175eddea8e87657285501aed0a6e790a046cbb45d190713679cc31e38ac2db2d05000000006b49d458","hash":"197a0f1dad97b06c3830f1ed3ab5165de3fd215df6dfe9536c00211aa6b792ab"}
+{"name":"hostile","card":{"v":1,"id":"vr-00000000000000ff","request":"eyJzY2hlbWEiOiJ3b3JsZC92ZXJiLXJlcXVlc3QvdjEiLCJ2ZXJiIjoiZXBpc29kZSBmaW5pc2giLCJhcmd2IjpbImFpbGFuZy13b3JsZGQiLCJlcGlzb2RlIiwiZmluaXNoIiwidzEiLCItLXByZXNldCIsInAiLCItLXRpdGxlIiwiaXQncyA8Yj4iLCItLXJlYWR5Il0sImVwaXNvZGUiOiJ3MSIsInByb21wdCI6ImNvbW1pdHRlcjogTWFyayA8bUB4PlxucmVwb3NpdG9yeTogc3VuaG9sby1kYXRhL2FpbGFuZy13b3JsZFxuYnJhbmNoOiB3b3JsZC93MVxuc2hhOiBhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiXG5iYXNlOiBkZXZcbnRpdGxlOiA8aW1nIHNyYz14IG9uZXJyb3I9YWxlcnQoMSk-IGFcdTAwMDBiICYgJ3EnXG5ib2R5ICg2NSBieXRlcyksIGZpcnN0IGxpbmVzOlxuPC9wcmU-PHNjcmlwdD5hbGVydCgyKTwvc2NyaXB0PlxuTlVMIFx1MDAwMCBoZXJlICZhbXA7IGFmdGVyIOKAlCDDvG7Dr1xuXG5sYWJlbDogYWlsYW5nLXdvcmxkICh3aWxsIGJlIGNyZWF0ZWQpXG5yZWFkeVxucHJvdGVjdGVkIHBhdGhzIGNoYW5nZWQ6IFxuIiwic2NvcGUiOiJnaXRodWItcHI6c3VuaG9sby1kYXRhL2FpbGFuZy13b3JsZCNyZW1vdGU9b3JpZ2luJmhlYWQ9d29ybGQvdzEmc2hhPWFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWImYmFzZT1kZXYmdGl0bGUtc2hhMjU2PTAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAmYm9keS1zaGEyNTY9MTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMSZsYWJlbD1haWxhbmctd29ybGQmZHJhZnQ9ZmFsc2UmY3JlYXRlLWxhYmVsPXRydWUmZXhwaXJlcz0xODAwMDAwOTAwIiwidGl0bGUiOiI8aW1nIHNyYz14IG9uZXJyb3I9YWxlcnQoMSk-IGFcdTAwMDBiICYgJ3EnIiwiYm9keSI6IjwvcHJlPjxzY3JpcHQ-YWxlcnQoMik8L3NjcmlwdD5cbk5VTCBcdTAwMDAgaGVyZSAmYW1wOyBhZnRlciDigJQgw7xuw69cbiIsImNvbW1pdHRlciI6Ik1hcmsgPG1AeD4iLCJwcm90ZWN0ZWQiOltdfQ","digest":"9f8cb56913ffb59e40efb482d528e9d00d5bdc3b9feb5887afdd6f6a6bffab4a","challenge":"a7f5c7cb4acece64fb98b2dfce7b4db39ae30e73d55322ed4ddfbd5e61048c69","created":1800000001,"expires":1800000601,"now":1800000005,"state":"pending","approvedBy":"","decision":"","passkeys":[{"origin":"http://localhost:1","rp":"localhost","allow":["Bwc"]}],"ssh":{"namespace":"world-verb@sunholo"}},"argvLine":"ailang-worldd episode finish w1 --preset p --title 'it'\"'\"'s \u003cb\u003e' --ready","message":"776f726c642d766572622f7631009f8cb56913ffb59e40efb482d528e9d00d5bdc3b9feb5887afdd6f6a6bffab4aa7f5c7cb4acece64fb98b2dfce7b4db39ae30e73d55322ed4ddfbd5e61048c69000000006b49d459","hash":"24746ae3ed0fbbf353267e4213afc87e60641cb7f771910782a44a53d9d33dc1"}
+{"name":"draft-approved","card":{"v":1,"id":"vr-0123456789abcdef","request":"eyJzY2hlbWEiOiJ3b3JsZC92ZXJiLXJlcXVlc3QvdjEiLCJ2ZXJiIjoiZXBpc29kZSBmaW5pc2giLCJhcmd2IjpbImFpbGFuZy13b3JsZGQiLCJlcGlzb2RlIiwiZmluaXNoIiwidzEiLCItLXByZXNldCIsInAiLCItLXRpdGxlIiwiRml4IHRoZSBiYXIiXSwiZXBpc29kZSI6IncxIiwicHJvbXB0IjoiY29tbWl0dGVyOiAocmVwb3NpdG9yeSBnaXQgY29uZmlnKVxucmVwb3NpdG9yeTogc3VuaG9sby1kYXRhL2FpbGFuZy13b3JsZFxuYnJhbmNoOiB3b3JsZC93MVxuc2hhOiBhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiXG5iYXNlOiBkZXZcbnRpdGxlOiBGaXggdGhlIGJhclxuYm9keSAoNzQgYnl0ZXMpLCBmaXJzdCBsaW5lczpcbiMjIFN1bW1hcnlcbmxpbmUgMlxubGluZSAzXG5saW5lIDRcbmxpbmUgNVxubGluZSA2XG5sYWJlbDogYWlsYW5nLXdvcmxkXG5kcmFmdFxucHJvdGVjdGVkIHBhdGhzIGNoYW5nZWQ6IGhvc3QveC5nbywgZG9jcy9hIGIubWRcbiIsInNjb3BlIjoiZ2l0aHViLXByOnN1bmhvbG8tZGF0YS9haWxhbmctd29ybGQjcmVtb3RlPW9yaWdpbiZoZWFkPXdvcmxkL3cxJnNoYT1hYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiJmJhc2U9ZGV2JnRpdGxlLXNoYTI1Nj0wMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwJmJvZHktc2hhMjU2PTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTEmbGFiZWw9YWlsYW5nLXdvcmxkJmRyYWZ0PXRydWUmY3JlYXRlLWxhYmVsPWZhbHNlJmV4cGlyZXM9MTgwMDAwMDkwMCIsInRpdGxlIjoiRml4IHRoZSBiYXIiLCJib2R5IjoiIyMgU3VtbWFyeVxubGluZSAyXG5saW5lIDNcbmxpbmUgNFxubGluZSA1XG5saW5lIDZcbmxpbmUgNyBpcyBub3QgaW4gdGhlIHByb21wdFxuIiwiY29tbWl0dGVyIjoiIiwicHJvdGVjdGVkIjpbImhvc3QveC5nbyIsImRvY3MvYSBiLm1kIl19","digest":"68e53444f32639d1676d16f5505f632cfae431352b3a68405a67ac3ec04b3f17","challenge":"5eddea8e87657285501aed0a6e790a046cbb45d190713679cc31e38ac2db2d05","created":1800000000,"expires":1800000600,"now":1800000005,"state":"approved","approvedBy":"webauthn:mark-mac","decision":"sha256:5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e","passkeys":[{"origin":"http://localhost:1","rp":"localhost","allow":["Bwc"]}],"ssh":{"namespace":"world-verb@sunholo"}},"argvLine":"ailang-worldd episode finish w1 --preset p --title 'Fix the bar'","message":"776f726c642d766572622f76310068e53444f32639d1676d16f5505f632cfae431352b3a68405a67ac3ec04b3f175eddea8e87657285501aed0a6e790a046cbb45d190713679cc31e38ac2db2d05000000006b49d458","hash":"197a0f1dad97b06c3830f1ed3ab5165de3fd215df6dfe9536c00211aa6b792ab"}
+{"name":"draft-superseded","card":{"v":1,"id":"vr-0123456789abcdef","request":"eyJzY2hlbWEiOiJ3b3JsZC92ZXJiLXJlcXVlc3QvdjEiLCJ2ZXJiIjoiZXBpc29kZSBmaW5pc2giLCJhcmd2IjpbImFpbGFuZy13b3JsZGQiLCJlcGlzb2RlIiwiZmluaXNoIiwidzEiLCItLXByZXNldCIsInAiLCItLXRpdGxlIiwiRml4IHRoZSBiYXIiXSwiZXBpc29kZSI6IncxIiwicHJvbXB0IjoiY29tbWl0dGVyOiAocmVwb3NpdG9yeSBnaXQgY29uZmlnKVxucmVwb3NpdG9yeTogc3VuaG9sby1kYXRhL2FpbGFuZy13b3JsZFxuYnJhbmNoOiB3b3JsZC93MVxuc2hhOiBhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiXG5iYXNlOiBkZXZcbnRpdGxlOiBGaXggdGhlIGJhclxuYm9keSAoNzQgYnl0ZXMpLCBmaXJzdCBsaW5lczpcbiMjIFN1bW1hcnlcbmxpbmUgMlxubGluZSAzXG5saW5lIDRcbmxpbmUgNVxubGluZSA2XG5sYWJlbDogYWlsYW5nLXdvcmxkXG5kcmFmdFxucHJvdGVjdGVkIHBhdGhzIGNoYW5nZWQ6IGhvc3QveC5nbywgZG9jcy9hIGIubWRcbiIsInNjb3BlIjoiZ2l0aHViLXByOnN1bmhvbG8tZGF0YS9haWxhbmctd29ybGQjcmVtb3RlPW9yaWdpbiZoZWFkPXdvcmxkL3cxJnNoYT1hYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiJmJhc2U9ZGV2JnRpdGxlLXNoYTI1Nj0wMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwJmJvZHktc2hhMjU2PTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTEmbGFiZWw9YWlsYW5nLXdvcmxkJmRyYWZ0PXRydWUmY3JlYXRlLWxhYmVsPWZhbHNlJmV4cGlyZXM9MTgwMDAwMDkwMCIsInRpdGxlIjoiRml4IHRoZSBiYXIiLCJib2R5IjoiIyMgU3VtbWFyeVxubGluZSAyXG5saW5lIDNcbmxpbmUgNFxubGluZSA1XG5saW5lIDZcbmxpbmUgNyBpcyBub3QgaW4gdGhlIHByb21wdFxuIiwiY29tbWl0dGVyIjoiIiwicHJvdGVjdGVkIjpbImhvc3QveC5nbyIsImRvY3MvYSBiLm1kIl19","digest":"68e53444f32639d1676d16f5505f632cfae431352b3a68405a67ac3ec04b3f17","challenge":"5eddea8e87657285501aed0a6e790a046cbb45d190713679cc31e38ac2db2d05","created":1800000000,"expires":1800000600,"now":1800000005,"state":"superseded","approvedBy":"","decision":"","passkeys":[{"origin":"http://localhost:1","rp":"localhost","allow":["Bwc"]}],"ssh":{"namespace":"world-verb@sunholo"}},"argvLine":"ailang-worldd episode finish w1 --preset p --title 'Fix the bar'","message":"776f726c642d766572622f76310068e53444f32639d1676d16f5505f632cfae431352b3a68405a67ac3ec04b3f175eddea8e87657285501aed0a6e790a046cbb45d190713679cc31e38ac2db2d05000000006b49d458","hash":"197a0f1dad97b06c3830f1ed3ab5165de3fd215df6dfe9536c00211aa6b792ab"}
+{"name":"draft-expired","card":{"v":1,"id":"vr-0123456789abcdef","request":"eyJzY2hlbWEiOiJ3b3JsZC92ZXJiLXJlcXVlc3QvdjEiLCJ2ZXJiIjoiZXBpc29kZSBmaW5pc2giLCJhcmd2IjpbImFpbGFuZy13b3JsZGQiLCJlcGlzb2RlIiwiZmluaXNoIiwidzEiLCItLXByZXNldCIsInAiLCItLXRpdGxlIiwiRml4IHRoZSBiYXIiXSwiZXBpc29kZSI6IncxIiwicHJvbXB0IjoiY29tbWl0dGVyOiAocmVwb3NpdG9yeSBnaXQgY29uZmlnKVxucmVwb3NpdG9yeTogc3VuaG9sby1kYXRhL2FpbGFuZy13b3JsZFxuYnJhbmNoOiB3b3JsZC93MVxuc2hhOiBhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiXG5iYXNlOiBkZXZcbnRpdGxlOiBGaXggdGhlIGJhclxuYm9keSAoNzQgYnl0ZXMpLCBmaXJzdCBsaW5lczpcbiMjIFN1bW1hcnlcbmxpbmUgMlxubGluZSAzXG5saW5lIDRcbmxpbmUgNVxubGluZSA2XG5sYWJlbDogYWlsYW5nLXdvcmxkXG5kcmFmdFxucHJvdGVjdGVkIHBhdGhzIGNoYW5nZWQ6IGhvc3QveC5nbywgZG9jcy9hIGIubWRcbiIsInNjb3BlIjoiZ2l0aHViLXByOnN1bmhvbG8tZGF0YS9haWxhbmctd29ybGQjcmVtb3RlPW9yaWdpbiZoZWFkPXdvcmxkL3cxJnNoYT1hYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiJmJhc2U9ZGV2JnRpdGxlLXNoYTI1Nj0wMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwJmJvZHktc2hhMjU2PTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTExMTEmbGFiZWw9YWlsYW5nLXdvcmxkJmRyYWZ0PXRydWUmY3JlYXRlLWxhYmVsPWZhbHNlJmV4cGlyZXM9MTgwMDAwMDkwMCIsInRpdGxlIjoiRml4IHRoZSBiYXIiLCJib2R5IjoiIyMgU3VtbWFyeVxubGluZSAyXG5saW5lIDNcbmxpbmUgNFxubGluZSA1XG5saW5lIDZcbmxpbmUgNyBpcyBub3QgaW4gdGhlIHByb21wdFxuIiwiY29tbWl0dGVyIjoiIiwicHJvdGVjdGVkIjpbImhvc3QveC5nbyIsImRvY3MvYSBiLm1kIl19","digest":"68e53444f32639d1676d16f5505f632cfae431352b3a68405a67ac3ec04b3f17","challenge":"5eddea8e87657285501aed0a6e790a046cbb45d190713679cc31e38ac2db2d05","created":1800000000,"expires":1800000600,"now":1800000600,"state":"expired","approvedBy":"","decision":"","passkeys":[{"origin":"http://localhost:1","rp":"localhost","allow":["Bwc"]}],"ssh":{"namespace":"world-verb@sunholo"}},"argvLine":"ailang-worldd episode finish w1 --preset p --title 'Fix the bar'","message":"776f726c642d766572622f76310068e53444f32639d1676d16f5505f632cfae431352b3a68405a67ac3ec04b3f175eddea8e87657285501aed0a6e790a046cbb45d190713679cc31e38ac2db2d05000000006b49d458","hash":"197a0f1dad97b06c3830f1ed3ab5165de3fd215df6dfe9536c00211aa6b792ab"}
````

### Patch M5b — `confirm.js` under node (design §6 CP1–CP6)

````diff M5b
diff --git a/host/surfacepage/confirm_script_test.go b/host/surfacepage/confirm_script_test.go
new file mode 100644
index 0000000..64c21bd
--- /dev/null
+++ b/host/surfacepage/confirm_script_test.go
@@ -0,0 +1,231 @@
+package surfacepage
+
+// Row 243 V3 sprint 2 M5 (design §2.5, §2.9 P4, §6 CP1-CP6): confirm.js under node, on a stand-in DOM
+// built from index.html's own ids, against the cards of testdata/confirm-card.golden.jsonl (written by
+// TestConfirmCardGolden from the real tools). Each test skips only when no node exists; with
+// WORLD_EXEC_NODE set, a harness that cannot run fails.
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
+// runConfirmHarness runs section of confirmHarness under node and requires its pass line.
+func runConfirmHarness(t *testing.T, section string) {
+	t.Helper()
+	node := os.Getenv("WORLD_EXEC_NODE")
+	if node == "" {
+		var err error
+		node, err = exec.LookPath("node")
+		if err != nil {
+			t.Skip("node unavailable locally")
+		}
+	}
+	golden, err := os.ReadFile(confirmGolden)
+	if err != nil {
+		t.Fatal(err)
+	}
+	dir := t.TempDir()
+	for name, data := range map[string][]byte{"confirm.js": ConfirmScript, "index.html": ConfirmPage, "golden.jsonl": golden, "harness.js": []byte(confirmHarness)} {
+		if err := os.WriteFile(filepath.Join(dir, name), data, 0o600); err != nil {
+			t.Fatal(err)
+		}
+	}
+	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
+	defer cancel()
+	cmd := exec.CommandContext(ctx, node, filepath.Join(dir, "harness.js"), section)
+	cmd.Dir = dir
+	out, err := cmd.CombinedOutput()
+	if err != nil {
+		t.Fatalf("node harness %s: %v\n%s", section, err, out)
+	}
+	if !strings.Contains(string(out), "CONFIRM "+section+" PASS") {
+		t.Fatalf("harness manifest missing: %s", out)
+	}
+}
+
+// CP1, CP2: a fetch that fails, times out (the 10 s abort) or answers non-2xx or non-JSON; an answer for
+// another id, or with a malformed digest, challenge or expiry; bytes that are not base64url, do not hash
+// to the digest or are not UTF-8 JSON; an unlisted, missing or mistyped field; a path that is not a card's
+// → a refusal with its reason, nothing shown, no controls, no message made. Control: the card is ready.
+func TestConfirmPageFailsClosed(t *testing.T) { runConfirmHarness(t, "closed") }
+
+// CP3: hostile values (markup, quotes, a NUL in the title and the body) are shown as text, whole, in every
+// field, the argv line, the prompt and the body; no HTML sink is written.
+func TestConfirmPageRendersText(t *testing.T) { runConfirmHarness(t, "text") }
+
+// CP4, CP6: the page's message is Go's VerbMessage (86 bytes) and its hash Go's VerbMessageHash; Save
+// the message is exactly those bytes as world-verb-<id>.msg; the signing line names world-verb@sunholo.
+func TestConfirmPageMessage(t *testing.T) { runConfirmHarness(t, "message") }
+
+// CP5: an approved, superseded, expired or unknown state shows that state and the fields, and no controls.
+func TestConfirmPageStates(t *testing.T) { runConfirmHarness(t, "states") }
+
+const confirmHarness = `
+const assert = require('node:assert/strict');
+const fs = require('node:fs');
+const vm = require('node:vm');
+const nodeCrypto = require('node:crypto');
+const SRC = fs.readFileSync('confirm.js', 'utf8'), HTML = fs.readFileSync('index.html', 'utf8');
+const G = Object.fromEntries(fs.readFileSync('golden.jsonl', 'utf8').trim().split('\n').map((l) => JSON.parse(l)).map((x) => [x.name, x]));
+const sha = (b) => nodeCrypto.createHash('sha256').update(b).digest('hex');
+const bytesOf = (card) => Buffer.from(card.request, 'base64url');
+const reqOf = (card) => JSON.parse(bytesOf(card).toString('utf8'));
+// withBytes is card carrying bytes b, its digest made to match.
+const withBytes = (card, b) => Object.assign({}, card, { request: Buffer.from(b).toString('base64url'), digest: sha(b) });
+
+// A stand-in DOM: index.html's elements by id (with their hidden attribute), and created ones.
+class El {
+  constructor(doc, tag, id, hidden) { Object.assign(this, { doc, tagName: tag, id: id || '', hidden: !!hidden, dataset: {}, attrs: {}, kids: [], text: '' }); }
+  set textContent(v) { this.text = String(v); this.kids = []; }
+  get textContent() { return this.text + this.kids.map((k) => k.textContent).join(''); }
+  set innerHTML(v) { this.doc.sinks.push(String(v)); }
+  set outerHTML(v) { this.doc.sinks.push(String(v)); }
+  insertAdjacentHTML(_, v) { this.doc.sinks.push(String(v)); }
+  append(...k) { this.kids.push(...k); }
+  toggleAttribute(n, on) { if (on) this.attrs[n] = ''; else delete this.attrs[n]; }
+  hasAttribute(n) { return n in this.attrs; }
+}
+function makeDoc() {
+  const doc = { sinks: [], byId: {} };
+  doc.body = new El(doc, 'body');
+  doc.body.dataset.state = HTML.match(/<body data-state="([a-z]+)">/)[1];
+  for (const m of HTML.matchAll(/<([a-z0-9]+) id="([a-z-]+)"( hidden)?/g)) doc.byId[m[2]] = new El(doc, m[1], m[2], !!m[3]);
+  const find = (els, id) => { for (const e of els) { if (e.id === id) return e; const k = find(e.kids, id); if (k) return k; } return null; };
+  doc.getElementById = (id) => doc.byId[id] || find(Object.values(doc.byId), id);
+  doc.createElement = (tag) => new El(doc, tag);
+  return doc;
+}
+const answer = (status, body) => () => Promise.resolve({ status, text: async () => (typeof body === 'string' ? body : JSON.stringify(body)) });
+const hang = (init) => new Promise((_, reject) => init.signal.addEventListener('abort', () => reject(new DOMException('The operation was aborted.', 'AbortError'))));
+
+// open runs confirm.js at path against fetch; fire runs the page's pending timers once it has asked.
+async function open(path, fetch, fire) {
+  const doc = makeDoc(), timers = [], blobs = [], fetches = [];
+  vm.runInNewContext(SRC, { document: doc, location: { pathname: path, origin: 'http://localhost:1' }, crypto: globalThis.crypto, TextEncoder, TextDecoder, atob, btoa, Blob, AbortController, console,
+    URL: { createObjectURL: (b) => { blobs.push(b); return 'blob:http://localhost:1/' + blobs.length; } },
+    setTimeout: (f, ms) => timers.push({ f, ms, cleared: false }), clearTimeout: (n) => { if (timers[n - 1]) timers[n - 1].cleared = true; },
+    fetch: (url, init) => { fetches.push({ url, init }); return fetch(init); } });
+  for (let i = 0; i < 400 && doc.body.dataset.state === 'loading'; i++) {
+    if (fire && i === 20) for (const t of timers) if (!t.cleared) t.f();
+    await new Promise((r) => setImmediate(r));
+  }
+  const $ = (id) => doc.getElementById(id);
+  return { doc, $, timers, blobs, fetches, state: doc.body.dataset.state, what: $('what').textContent };
+}
+const card = (c) => '/surface/confirm/' + c.id;
+// refused asserts a fail-closed page: state, the reason, nothing shown, no controls, no message made.
+function refused(p, state, reason, name) {
+  assert.equal(p.state, state, name + ': ' + p.what);
+  assert(p.what.includes(reason), name + ': ' + p.what);
+  assert(p.$('what').hasAttribute('data-error'), name);
+  assert(p.$('shown').hidden && p.$('controls').hidden && p.blobs.length === 0 && p.$('sshline').textContent === '' && p.$('prompt').textContent === '', name);
+}
+
+const sections = {
+  async closed() {
+    const c = G.draft.card, other = G.hostile.card;
+    const ok = await open(card(c), answer(200, c));
+    assert.equal(ok.state, 'ready', ok.what);
+    assert.equal(ok.fetches.length, 1);
+    assert.equal(ok.fetches[0].url, '/surface/confirmations/' + c.id);
+    assert(ok.fetches[0].init.cache === 'no-store' && ok.fetches[0].init.credentials === 'omit' && ok.fetches[0].init.signal);
+    assert(ok.timers.length === 1 && ok.timers[0].ms === 10000 && ok.timers[0].cleared, JSON.stringify(ok.timers));
+    refused(await open(card(c), () => Promise.reject(new TypeError('Failed to fetch'))), 'failed', 'The request could not be read (Failed to fetch).', 'network');
+    const slow = await open(card(c), hang, true);
+    refused(slow, 'failed', 'The request could not be read (no answer within 10 s).', 'timeout');
+    assert(slow.timers.length === 1 && slow.timers[0].ms === 10000);
+    refused(await open(card(c), answer(404, { reason: 'no such verb request' })), 'failed', 'The daemon did not give the request (HTTP 404: no such verb request).', '404');
+    refused(await open(card(c), answer(503, 'down')), 'failed', '(HTTP 503)', '503');
+    refused(await open(card(c), answer(200, '<html>')), 'failed', '(HTTP 200)', 'not JSON');
+    const malformed = 'The daemon\'s answer is malformed; nothing to confirm.';
+    for (const [name, edit] of [['another id', { id: other.id }], ['digest case', { digest: c.digest.toUpperCase() }], ['short challenge', { challenge: c.challenge.slice(1) }],
+      ['negative expiry', { expires: -1 }], ['fractional expiry', { expires: c.expires + 0.5 }], ['no state', { state: 1 }]]) {
+      refused(await open(card(c), answer(200, Object.assign({}, c, edit))), 'refused', malformed, name);
+    }
+    refused(await open(card(c), answer(200, Object.assign({}, c, { request: '*' + c.request }))), 'refused', 'is not base64url', 'base64url');
+    const digest = 'The request\'s bytes do not match its digest; nothing to confirm.';
+    refused(await open(card(c), answer(200, Object.assign({}, c, { digest: sha('other') }))), 'refused', digest, 'digest');
+    refused(await open(card(c), answer(200, Object.assign({}, c, { request: other.request }))), 'refused', digest, 'swapped bytes');
+    const padded = other.request + '='.repeat((4 - other.request.length % 4) % 4); // standard base64's padding: atob takes it
+    assert(padded !== other.request && Buffer.from(padded, 'base64').equals(bytesOf(other)));
+    refused(await open(card(other), answer(200, Object.assign({}, other, { request: padded }))), 'refused', 'is not base64url', 'padded base64');
+    const latin = Buffer.from(JSON.stringify(Object.assign({}, reqOf(c), { title: 'TITLE' })).replace('TITLE', '\\u00ff'), 'utf8');
+    refused(await open(card(c), answer(200, withBytes(c, Buffer.concat([latin.subarray(0, latin.indexOf('\\u00ff')), Buffer.from([0xff]), latin.subarray(latin.indexOf('\\u00ff') + 6)])))),
+      'refused', 'is not UTF-8 JSON', 'a lone 0xff in the title');
+    const r = reqOf(c), unshown = 'This request has a field the page cannot show (';
+    for (const [name, obj, key] of [['unlisted', Object.assign({}, r, { extra: 'x' }), 'extra'], ['missing', Object.assign({}, r, { committer: undefined }), 'committer'],
+      ['argv a string', Object.assign({}, r, { argv: 'ailang-worldd episode finish w1' }), 'argv'], ['protected a number', Object.assign({}, r, { protected: [1] }), 'protected'],
+      ['title a number', Object.assign({}, r, { title: 7 }), 'title'], ['body null', Object.assign({}, r, { body: null }), 'body'], ['an array', [r], 'the request']]) {
+      refused(await open(card(c), answer(200, withBytes(c, JSON.stringify(obj)))), 'refused', unshown + key + '); nothing to confirm here; confirm at the terminal.', name);
+    }
+    const lost = await open('/surface/confirm/vr-0123', answer(200, c));
+    refused(lost, 'refused', 'This is not a confirmation card\'s address; nothing to confirm.', 'path');
+    assert.equal(lost.fetches.length, 0);
+  },
+  async text() {
+    for (const name of ['hostile', 'draft']) {
+      const c = G[name].card, r = reqOf(c), p = await open(card(c), answer(200, c));
+      assert.equal(p.state, 'ready', p.what);
+      assert.equal(p.$('argv').textContent, G[name].argvLine);
+      assert.equal(p.$('prompt').textContent, r.prompt);
+      for (const k of ['schema', 'verb', 'episode', 'prompt', 'scope', 'title', 'committer']) assert.equal(p.$('f-' + k).textContent, r[k], k);
+      for (const k of ['argv', 'protected']) assert.equal(p.$('f-' + k).textContent, r[k].join('\n'), k);
+      assert.equal(p.$('f-body').textContent, r.body);
+      assert.equal(p.$('f-body-size').textContent, Buffer.byteLength(r.body) + ' bytes' + 'the full body' + r.body);
+      assert.equal(p.$('m-digest').textContent, c.digest);
+      assert.equal(p.$('m-challenge').textContent, c.challenge);
+      assert.equal(p.$('m-expires').textContent, new Date(c.expires * 1000).toISOString() + ' (' + c.expires + ')');
+      assert.deepEqual(p.doc.sinks, []);
+    }
+    const r = reqOf(G.hostile.card), p = await open(card(G.hostile.card), answer(200, G.hostile.card));
+    // The NUL rows (S1 Appendix C, P6): the text after a NUL is shown, in the title, the prompt and the body.
+    assert(r.title.includes('a\u0000b & ') && r.body.includes('NUL \u0000 here &amp; after'));
+    assert(p.$('f-title').textContent.endsWith('a\u0000b & \'q\''));
+    assert(p.$('prompt').textContent.includes('title: <img src=x onerror=alert(1)> a\u0000b & \'q\'\n'));
+    assert(p.$('f-body').textContent.includes('</pre><script>alert(2)</script>\nNUL \u0000 here &amp; after — ünï\n'));
+    assert(Buffer.byteLength(r.body) > r.body.length); // the size line counts bytes, not characters
+    assert(p.$('argv').textContent.endsWith(' --title \'it\'"\'"\'s <b>\' --ready'), p.$('argv').textContent);
+  },
+  async message() {
+    for (const name of ['draft', 'hostile']) {
+      const g = G[name], c = g.card, p = await open(card(c), answer(200, c));
+      assert.equal(p.state, 'ready', p.what);
+      assert.equal(p.$('m-hash').textContent, g.hash);
+      assert(p.what.endsWith('; message hash ' + g.hash), p.what);
+      assert.equal(p.blobs.length, 1);
+      assert.equal(p.blobs[0].type, 'application/octet-stream');
+      const m = Buffer.from(await p.blobs[0].arrayBuffer());
+      assert.equal(m.length, 86);
+      assert.equal(m.toString('hex'), g.message);
+      assert.equal(sha(m), g.hash);
+      assert.equal(p.$('save').href, 'blob:http://localhost:1/1');
+      assert.equal(p.$('save').download, 'world-verb-' + c.id + '.msg');
+      assert.equal(p.$('sshline').textContent, 'ssh-keygen -Y sign -n world-verb@sunholo -f <key> world-verb-' + c.id + '.msg');
+      assert(!p.$('controls').hidden && !p.$('shown').hidden);
+    }
+  },
+  async states() {
+    const ready = await open(card(G.draft.card), answer(200, G.draft.card));
+    assert(ready.state === 'ready' && !ready.$('controls').hidden && ready.blobs.length === 1);
+    for (const [c, status] of [[G['draft-approved'].card, 'Already confirmed by webauthn:mark-mac: decision sha256:' + '5e'.repeat(32) + '; nothing to confirm here.'],
+      [G['draft-superseded'].card, 'Superseded by a newer finish of this episode; nothing to confirm here.'],
+      [G['draft-expired'].card, 'Expired; nothing to confirm here.'],
+      [Object.assign({}, G.draft.card, { state: 'odd' }), 'This request is odd; nothing to confirm here.']]) {
+      const p = await open(card(c), answer(200, c));
+      assert.equal(p.state, 'closed', c.state + ': ' + p.what);
+      assert.equal(p.$('status').textContent, status);
+      assert(p.$('controls').hidden && p.blobs.length === 0 && p.$('sshline').textContent === '' && p.$('save').download === undefined, c.state);
+      assert(!p.$('shown').hidden && p.$('f-title').textContent === reqOf(c).title, c.state);
+    }
+    assert.deepEqual([G['draft-approved'].card.state, G['draft-superseded'].card.state, G['draft-expired'].card.state], ['approved', 'superseded', 'expired']);
+  },
+};
+const which = process.argv[2];
+sections[which]().then(() => console.log('CONFIRM ' + which + ' PASS'), (e) => { console.error(e); process.exit(1); });
+`
````

### Patch M6 — XC1 live Chrome (design §6 XC1)

````diff M6
diff --git a/cmd/ailang-worldd/confirm_page_chrome_test.go b/cmd/ailang-worldd/confirm_page_chrome_test.go
new file mode 100644
index 0000000..833ac46
--- /dev/null
+++ b/cmd/ailang-worldd/confirm_page_chrome_test.go
@@ -0,0 +1,109 @@
+package main
+
+// Row 243 V3 sprint 2 M6 (design §2.5, §6 XC1): the confirmation card in a live Chrome, on the A2/A3
+// driver (testdata/passkey_cdp.mjs, unchanged), gated as XB1 (WORLD_CHROME,
+// WORLD_REQUIRE_PASSKEY_BROWSER=1; skipped and counted without them). No daemon route serves the card
+// before sprint 3, so a loopback test handler serves surfacepage.ConfirmPage and ConfirmScript under
+// the headers the daemon gives its own page scripts (surfaceCSP among them), and the card's data, made
+// by surfacepage.ConfirmationCard from approval.CanonicalVerb's bytes. No store is read.
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"encoding/hex"
+	"encoding/json"
+	"io"
+	"net/http"
+	"net/http/httptest"
+	"path/filepath"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// XC1: the card shows, in its <pre>, exactly broker.FinishPromptText of the request's scope, title,
+// body, protected paths and committer (the terminal's lines before its question), and Go's
+// VerbMessageHash, under the daemon's CSP with no violation.
+func TestConfirmPageShowsTheFinishPromptInChrome(t *testing.T) {
+	b := startPasskeyBrowser(t)
+	b.do(t, nil, "tab", map[string]any{})
+	now := time.Now().Unix()
+	scope := broker.FinishApprovalScope{Repository: "sunholo-data/ailang-world", Remote: "origin", Head: "world/w1", SHA: strings.Repeat("ab", 20), Base: "dev",
+		TitleSHA256: strings.Repeat("0", 64), BodySHA256: strings.Repeat("1", 64), Draft: true, ExpiresAt: now + 900}
+	title, body, protected, committer := "Fix the <bar>", "## Summary\nline 2\nline 3\nline 4\nline 5\nline 6\nline 7 is not in the prompt\n", []string{"host/x.go"}, ""
+	prompt := broker.FinishPromptText(scope, title, body, protected, committer)
+	var tty bytes.Buffer // control: the terminal prints exactly these lines, then its question
+	finishPrompt(&tty, scope, title, body, protected, committer)
+	if tty.String() != prompt+"Open this pull request? [y/N] " {
+		t.Fatalf("FinishPromptText is not the terminal's prompt: %q", tty.String())
+	}
+	r := approval.VerbRequest{Schema: approval.VerbRequestSchema, Verb: approval.VerbFinish, Argv: []string{"ailang-worldd", "episode", "finish", "w1", "--preset", "p"},
+		Episode: "w1", Prompt: prompt, Scope: broker.FormatFinishApprovalScope(scope), Title: title, Body: body, Committer: committer, Protected: protected}
+	req, err := approval.CanonicalVerb(r)
+	if err == nil {
+		err = r.Validate()
+	}
+	if err != nil {
+		t.Fatal(err)
+	}
+	digest, challenge := sha256.Sum256(req), sha256.Sum256([]byte("xc1 challenge"))
+	const id = "vr-00000000000000c1"
+	card, err := json.Marshal(surfacepage.ConfirmationCard(surfacepage.ConfirmationInput{ID: id, Request: req, Digest: hex.EncodeToString(digest[:]),
+		Challenge: hex.EncodeToString(challenge[:]), Created: now, Expires: now + 600, Now: now}, nil))
+	if err != nil {
+		t.Fatal(err)
+	}
+	hash, _ := approval.VerbMessageHash(digest, challenge, now+600)
+
+	d, err := daemon.New(boundedTestContext(t), daemon.Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: daemon.DefaultBindHost, ErrorLog: io.Discard})
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { _ = d.Close() })
+	own := httptest.NewRecorder() // the headers the daemon serves its passkey page's script with
+	d.Handler().ServeHTTP(own, httptest.NewRequest("GET", "/surface/passkey.js", nil))
+	if csp := own.Header().Get("Content-Security-Policy"); own.Code != 200 || !strings.Contains(csp, "script-src 'self'") || !strings.Contains(csp, "connect-src 'self'") {
+		t.Fatalf("the daemon's own script headers: %d %q", own.Code, csp)
+	}
+	serve := func(w http.ResponseWriter, contentType string, b []byte) {
+		for k, v := range own.Header() {
+			w.Header()[k] = v
+		}
+		w.Header().Set("Content-Type", contentType)
+		_, _ = w.Write(b)
+	}
+	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, q *http.Request) {
+		switch q.URL.Path {
+		case "/surface/confirm/" + id:
+			serve(w, "text/html; charset=utf-8", surfacepage.ConfirmPage)
+		case "/surface/confirm.js":
+			serve(w, "text/javascript; charset=utf-8", surfacepage.ConfirmScript)
+		case "/surface/confirmations/" + id:
+			serve(w, "application/json", card)
+		default:
+			http.NotFound(w, q)
+		}
+	}))
+	t.Cleanup(srv.Close)
+	origin := strings.Replace(srv.URL, "127.0.0.1", "localhost", 1)
+
+	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface/confirm/" + id, "raw": true})
+	var view struct {
+		State, Prompt, Argv, Hash, Save, SSHLine, Title string
+		Violations                                      []string
+	}
+	b.do(t, &view, "waitFor", map[string]any{"expression": `document.body.dataset.state !== 'loading' && (() => { const $ = (id) => document.getElementById(id); return {
+		state: document.body.dataset.state, prompt: $('prompt').textContent, argv: $('argv').textContent, hash: $('m-hash') ? $('m-hash').textContent : '',
+		save: $('save').download, sshLine: $('sshline').textContent, title: $('f-title') ? $('f-title').textContent : '', violations: window.__cspv || [] }; })()`})
+	if view.State != "ready" || view.Prompt != prompt || view.Title != title || view.Argv != "ailang-worldd episode finish w1 --preset p" {
+		t.Fatalf("the card shows %+v; want the prompt %q", view, prompt)
+	}
+	if view.Hash != hex.EncodeToString(hash[:]) || view.Save != "world-verb-"+id+".msg" || view.SSHLine != "ssh-keygen -Y sign -n world-verb@sunholo -f <key> world-verb-"+id+".msg" || len(view.Violations) != 0 {
+		t.Fatalf("the card's message, link or CSP: %+v", view)
+	}
+}
````

### Patch S1b — the judges (design §2.3, §6 VJ1; iteration 400's patch + three VJ1 test lines)

````diff S1b
diff --git a/host/authority/judge_verb.go b/host/authority/judge_verb.go
new file mode 100644
index 0000000..b2a64ba
--- /dev/null
+++ b/host/authority/judge_verb.go
@@ -0,0 +1,129 @@
+package authority
+
+// Row 243 V3 sprint 1 M2 (design §2.3): the verb judges beside the mint and decision judges, on a
+// read-only handle, with judge.go's step names. JudgeVerbSSH and JudgeVerbAssertion are the daemon's
+// C11 over a stored verb request's primitives (no store type: sprint 1 adds none); RecheckVerb is
+// the CLI's re-check of a recorded confirmation (approval.CheckConfirmation, design §2.7 step 6(a)).
+
+import (
+	"context"
+	"errors"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// VerbJudgeInput is a stored verb request at now: its digest column and challenge, the SHA-256 of
+// its stored bytes, its times and whether it is consumed.
+type VerbJudgeInput struct {
+	Digest, Challenge, BytesDigest [32]byte
+	CreatedAt, ExpiresAt, Now      int64
+	Consumed                       bool
+}
+
+func (in VerbJudgeInput) admit() approval.AdmitInput {
+	return approval.AdmitInput{Now: in.Now, CreatedAt: in.CreatedAt, ExpiresAt: in.ExpiresAt, Consumed: in.Consumed,
+		StoredDigest: in.Digest, BytesDigest: in.BytesDigest}
+}
+
+// sshApprover is the enrolled approver an SSHSIG's key names (ApproverByKey on ro), or a *JudgeError.
+func sshApprover(ctx context.Context, ro *store.Credentials, sig *approval.SSHSig) (approval.Approver, error) {
+	enrolled, ok, err := ro.ApproverByKey(ctx, sig.PublicKey.Blob)
+	if err != nil {
+		return approval.Approver{}, &JudgeError{Step: JudgeStepStore, Err: err}
+	}
+	if !ok {
+		return approval.Approver{}, &JudgeError{Step: JudgeStepUnknownKey, Fingerprint: sig.PublicKey.Fingerprint(), Err: ErrNotEnrolled}
+	}
+	return approval.Approver{Name: enrolled.Name, Format: enrolled.Format, Key: approval.SSHPublicKey{Blob: enrolled.PublicKey},
+		EnrolledAt: enrolled.EnrolledAt, LastUsedAt: enrolled.LastUsedAt, RevokedAt: enrolled.RevokedAt}, nil
+}
+
+// passkeyApprover is the enrolled approver, with its passkey, an assertion's credential ID names
+// (PasskeyByCredentialID on ro, ParseStoredCOSEKey), or a *JudgeError.
+func passkeyApprover(ctx context.Context, ro *store.Credentials, resp approval.AssertionResponse) (approval.Approver, error) {
+	a, p, ok, err := ro.PasskeyByCredentialID(ctx, resp.RawID)
+	if err != nil {
+		return approval.Approver{}, &JudgeError{Step: JudgeStepStore, Err: err}
+	}
+	if !ok {
+		return approval.Approver{}, &JudgeError{Step: JudgeStepUnknownKey, RawID: resp.RawID, Err: ErrNotEnrolled}
+	}
+	key, err := approval.ParseStoredCOSEKey(a.PublicKey)
+	if err != nil {
+		return approval.Approver{}, &JudgeError{Step: JudgeStepStoredKey, Approver: a.Name, Err: err}
+	}
+	return approval.Approver{Name: a.Name, Format: a.Format, EnrolledAt: a.EnrolledAt, LastUsedAt: a.LastUsedAt, RevokedAt: a.RevokedAt,
+		Passkey: &approval.Passkey{RPID: p.RPID, Origin: p.Origin, CredentialID: p.CredentialID, Key: key,
+			BackupEligible: p.BackupEligible, SignCount: uint32(p.SignCount)}}, nil
+}
+
+// JudgeVerbSSH judges an armored SSHSIG confirmation of in: ParseSSHSIG, the approver its key names,
+// then approval.ConfirmBySignature. It returns the verified confirmation or a *JudgeError.
+func JudgeVerbSSH(ctx context.Context, ro *store.Credentials, in VerbJudgeInput, armored []byte) (*approval.SignedVerb, error) {
+	sig, err := approval.ParseSSHSIG(armored)
+	if err != nil {
+		return nil, &JudgeError{Step: JudgeStepParse, Err: err}
+	}
+	a, err := sshApprover(ctx, ro, sig)
+	if err != nil {
+		return nil, err
+	}
+	signed, err := approval.ConfirmBySignature(approval.VerbSignatureInput{Admit: in.admit(), Challenge: in.Challenge, Signature: armored, Approver: a})
+	if err != nil {
+		return nil, &JudgeError{Step: JudgeStepLaw, Err: err}
+	}
+	return signed, nil
+}
+
+// JudgeVerbAssertion judges a passkey confirmation of in: ParseAssertionResponse, the passkey its
+// credential ID names, its stored key, then approval.ConfirmByAssertion (W15 in full). It returns the
+// verified confirmation or a *JudgeError.
+func JudgeVerbAssertion(ctx context.Context, ro *store.Credentials, in VerbJudgeInput, response []byte) (*approval.SignedVerb, error) {
+	resp, err := approval.ParseAssertionResponse(response)
+	if err != nil {
+		return nil, &JudgeError{Step: JudgeStepParse, Err: err}
+	}
+	a, err := passkeyApprover(ctx, ro, resp)
+	if err != nil {
+		return nil, err
+	}
+	signed, err := approval.ConfirmByAssertion(approval.VerbAssertionInput{Admit: in.admit(), Challenge: in.Challenge, Response: response, Approver: a})
+	if err != nil {
+		return nil, &JudgeError{Step: JudgeStepLaw, Err: err}
+	}
+	return signed, nil
+}
+
+// ErrUnknownArtifact is a recorded artifact that is neither an armored SSHSIG nor an assertion.
+var ErrUnknownArtifact = errors.New("authority: the artifact is neither an SSH signature nor a passkey assertion")
+
+// RecheckVerb is the CLI's re-check of a recorded confirmation (design §2.3 [r1], §2.7 step 6(a)):
+// the artifact's kind (an armored SSHSIG, else a passkey assertion), the enrolled approver it names
+// on ro, then approval.CheckConfirmation over in's digest, challenge and expiry at in.Now — the
+// caller's own values, never the row's. in.CreatedAt, BytesDigest and Consumed are not read: the
+// row is consumed by then (no P1) and the counter has moved (no W15). It returns the verified
+// confirmation or a *JudgeError.
+func RecheckVerb(ctx context.Context, ro *store.Credentials, in VerbJudgeInput, artifact []byte) (*approval.SignedVerb, error) {
+	check := approval.VerbCheckInput{Digest: in.Digest, Challenge: in.Challenge, Expiry: in.ExpiresAt, Now: in.Now, Artifact: artifact}
+	if sig, sigErr := approval.ParseSSHSIG(artifact); sigErr == nil {
+		a, err := sshApprover(ctx, ro, sig)
+		if err != nil {
+			return nil, err
+		}
+		check.Format, check.Approver = approval.FormatSSH, a
+	} else if resp, respErr := approval.ParseAssertionResponse(artifact); respErr == nil {
+		a, err := passkeyApprover(ctx, ro, resp)
+		if err != nil {
+			return nil, err
+		}
+		check.Format, check.Approver = approval.FormatWebAuthn, a
+	} else {
+		return nil, &JudgeError{Step: JudgeStepParse, Err: errors.Join(ErrUnknownArtifact, sigErr, respErr)}
+	}
+	signed, err := approval.CheckConfirmation(check)
+	if err != nil {
+		return nil, &JudgeError{Step: JudgeStepLaw, Err: err}
+	}
+	return signed, nil
+}
diff --git a/host/authority/judge_verb_test.go b/host/authority/judge_verb_test.go
new file mode 100644
index 0000000..9136e5d
--- /dev/null
+++ b/host/authority/judge_verb_test.go
@@ -0,0 +1,200 @@
+package authority
+
+// Row 243 V3 sprint 1 M2 (design §2.3, §6 VJ1): JudgeVerbSSH, JudgeVerbAssertion and RecheckVerb
+// name each step as the mint and decision judges do — parse, store, unknown key, stored key, law —
+// in that order, on a read-only handle. Every key is approvaltest's (D-WORLD-109); the rig is
+// judge_test.go's (mark-mac's passkey enrolled with counter 1, mark-ssh, both at 900).
+
+import (
+	"bytes"
+	"errors"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+var vjIn = VerbJudgeInput{Digest: [32]byte{0x76, 1}, Challenge: [32]byte{0x76, 2}, BytesDigest: [32]byte{0x76, 1}, CreatedAt: 1000, ExpiresAt: 1600, Now: 1010}
+
+const vjDropV4 = `DROP TABLE key_decisions; DROP TABLE surface_approvals; DROP TABLE webauthn_credentials; DROP TABLE approvers; PRAGMA user_version = 4`
+const vjDropV5 = `DROP TABLE key_decisions; DROP TABLE surface_approvals; DROP TABLE webauthn_credentials; PRAGMA user_version = 5`
+
+// vjSign and vjAssert are a confirmation of in's request by s or pk.
+func vjSign(t *testing.T, s *approvaltest.Signer, namespace string, in VerbJudgeInput) []byte {
+	t.Helper()
+	m, err := approval.VerbMessage(in.Digest, in.Challenge, in.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return s.Sign(namespace, m)
+}
+
+func vjAssert(t *testing.T, pk *approvaltest.Passkey, in VerbJudgeInput, count uint32) []byte {
+	t.Helper()
+	h, err := approval.VerbMessageHash(in.Digest, in.Challenge, in.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return pk.Assert(jgOrigin, jgRPID, h, count)
+}
+
+// vjCase is one refusal: setup runs on the rig, then judge(artifact) must fail at step with cause
+// (when set) and the unknown-key / stored-key fields.
+type vjCase struct {
+	name     string
+	setup    string
+	in       VerbJudgeInput
+	artifact []byte
+	step     string
+	cause    error
+	field    string // the unknown-key SSHSIG fingerprint
+	rawID    []byte // the unknown-key credential ID
+	approver string // the stored-key approver
+}
+
+type vjJudge func(t *testing.T, ro *store.Credentials, in VerbJudgeInput, artifact []byte) (*approval.SignedVerb, error)
+
+func vjRun(t *testing.T, judge vjJudge, cases []vjCase) {
+	t.Helper()
+	for _, tc := range cases {
+		path, _ := jgRig(t)
+		if tc.setup != "" {
+			jgExec(t, path, tc.setup)
+		}
+		ro, err := store.OpenCredentialsReadOnly(path)
+		if err != nil {
+			t.Fatal(err)
+		}
+		v, err := judge(t, ro, tc.in, tc.artifact)
+		_ = ro.Close()
+		var je *JudgeError
+		if v != nil || !errors.As(err, &je) || je.Step != tc.step || (tc.cause != nil && !errors.Is(err, tc.cause)) ||
+			je.Fingerprint != tc.field || !bytes.Equal(je.RawID, tc.rawID) || je.Approver != tc.approver {
+			t.Errorf("%s: %v, %v; want step %s cause %v", tc.name, v, err, tc.step, tc.cause)
+		}
+	}
+}
+
+// vjOK opens the rig after setup and returns judge's verified confirmation.
+func vjOK(t *testing.T, judge vjJudge, setup string, in VerbJudgeInput, artifact []byte) *approval.SignedVerb {
+	t.Helper()
+	path, _ := jgRig(t)
+	if setup != "" {
+		jgExec(t, path, setup)
+	}
+	ro, err := store.OpenCredentialsReadOnly(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = ro.Close() }()
+	v, err := judge(t, ro, in, artifact)
+	if err != nil || v.Binds(in.Digest, in.Challenge, in.ExpiresAt) != nil {
+		t.Fatalf("control: %v, %v", v, err)
+	}
+	return v
+}
+
+func TestJudgeVerbSSH(t *testing.T) {
+	judge := func(t *testing.T, ro *store.Credentials, in VerbJudgeInput, a []byte) (*approval.SignedVerb, error) {
+		return JudgeVerbSSH(boundedTestContext(t), ro, in, a)
+	}
+	good := vjSign(t, jgSigner, approval.VerbNamespace, vjIn)
+	stranger := approvaltest.NewSigner(66)
+	mintMsg, _ := approval.SignedMessage(vjIn.Digest, vjIn.Challenge, vjIn.ExpiresAt)
+	consumed, otherChallenge, badBytes := vjIn, vjIn, vjIn
+	consumed.Consumed = true
+	otherChallenge.Challenge[0] ^= 1
+	badBytes.BytesDigest[0] ^= 1 // the stored bytes do not hash to the digest column (P1)
+	vjRun(t, judge, []vjCase{
+		{name: "VJ-parse", in: vjIn, artifact: []byte("y\n"), step: JudgeStepParse, cause: approval.ErrArmor},
+		{name: "VJ-parse before store", setup: vjDropV4, in: vjIn, artifact: []byte("y\n"), step: JudgeStepParse},
+		{name: "VJ-store-v4", setup: vjDropV4, in: vjIn, artifact: good, step: JudgeStepStore},
+		{name: "VJ-unknown-key", in: vjIn, artifact: vjSign(t, stranger, approval.VerbNamespace, vjIn), step: JudgeStepUnknownKey, cause: ErrNotEnrolled,
+			field: stranger.PublicKey().Fingerprint()},
+		{name: "VJ-unknown-key before law", in: consumed, artifact: vjSign(t, stranger, approval.VerbNamespace, vjIn), step: JudgeStepUnknownKey,
+			field: stranger.PublicKey().Fingerprint()},
+		{name: "VJ-law-consumed", in: consumed, artifact: good, step: JudgeStepLaw},
+		{name: "VJ-law-bytes-digest", in: badBytes, artifact: good, step: JudgeStepLaw},
+		{name: "VJ-law-mint-signature", in: vjIn, artifact: jgSigner.Sign(approval.Namespace, mintMsg), step: JudgeStepLaw, cause: approval.ErrNamespace},
+		{name: "VJ-law-other-challenge", in: otherChallenge, artifact: good, step: JudgeStepLaw, cause: approval.ErrBadSignature},
+		{name: "VJ-law-revoked", setup: `UPDATE approvers SET revoked_at = 950 WHERE name = 'mark-ssh'`, in: vjIn, artifact: good, step: JudgeStepLaw, cause: approval.ErrApproverRevoked},
+	})
+	var na *approval.NotAdmittedError
+	path, _ := jgRig(t)
+	ro, _ := store.OpenCredentialsReadOnly(path)
+	if _, err := JudgeVerbSSH(boundedTestContext(t), ro, consumed, good); !errors.As(err, &na) || na.Verdict != approval.RefuseConsumed {
+		t.Errorf("consumed: %v, want the P1 verdict through the law step", err)
+	}
+	_ = ro.Close()
+	if v := vjOK(t, judge, "", vjIn, good); v.ApprovedBy() != "ssh:mark-ssh" {
+		t.Fatalf("ApprovedBy = %q", v.ApprovedBy())
+	}
+}
+
+func TestJudgeVerbAssertion(t *testing.T) {
+	judge := func(t *testing.T, ro *store.Credentials, in VerbJudgeInput, a []byte) (*approval.SignedVerb, error) {
+		return JudgeVerbAssertion(boundedTestContext(t), ro, in, a)
+	}
+	stranger := approvaltest.NewPasskey(approval.AlgEdDSA, 67)
+	mint, _ := approval.MessageHash(vjIn.Digest, vjIn.Challenge, vjIn.ExpiresAt)
+	expired := vjIn
+	expired.Now = vjIn.ExpiresAt
+	vjRun(t, judge, []vjCase{
+		{name: "VJ-parse", in: vjIn, artifact: []byte("{}"), step: JudgeStepParse, cause: approval.ErrPasskeyResponse},
+		{name: "VJ-store-v5", setup: vjDropV5, in: vjIn, artifact: vjAssert(t, jgPasskey, vjIn, 2), step: JudgeStepStore},
+		{name: "VJ-unknown-key", in: vjIn, artifact: vjAssert(t, stranger, vjIn, 2), step: JudgeStepUnknownKey, cause: ErrNotEnrolled, rawID: stranger.CredentialID()},
+		{name: "VJ-stored-key", setup: `UPDATE approvers SET public_key = X'00' WHERE name = 'mark-mac'`, in: vjIn, artifact: vjAssert(t, jgPasskey, vjIn, 2),
+			step: JudgeStepStoredKey, approver: "mark-mac"},
+		{name: "VJ-stored-key before law", setup: `UPDATE approvers SET public_key = X'00' WHERE name = 'mark-mac'`, in: expired, artifact: vjAssert(t, jgPasskey, vjIn, 2),
+			step: JudgeStepStoredKey, approver: "mark-mac"},
+		{name: "VJ-law-expired", in: expired, artifact: vjAssert(t, jgPasskey, vjIn, 2), step: JudgeStepLaw},
+		{name: "VJ-law-mint-assertion", in: vjIn, artifact: jgPasskey.Assert(jgOrigin, jgRPID, mint, 2), step: JudgeStepLaw, cause: approval.ErrChallenge},
+		{name: "VJ-law-W15", in: vjIn, artifact: vjAssert(t, jgPasskey, vjIn, 1), step: JudgeStepLaw, cause: approval.ErrSignCountRegression},
+		{name: "VJ-law-origin", setup: `UPDATE webauthn_credentials SET origin = 'https://x.ts.net' WHERE name = 'mark-mac'`, in: vjIn, artifact: vjAssert(t, jgPasskey, vjIn, 2),
+			step: JudgeStepLaw, cause: approval.ErrOriginMismatch},
+	})
+	v := vjOK(t, judge, "", vjIn, vjAssert(t, jgPasskey, vjIn, 2))
+	if id, n, ok := v.PasskeyUse(); v.ApprovedBy() != "webauthn:mark-mac" || !ok || n != 2 || !bytes.Equal(id, jgPasskey.CredentialID()) {
+		t.Fatalf("control: %q %x %d %v", v.ApprovedBy(), id, n, ok)
+	}
+}
+
+// RecheckVerb runs after the consume: the row is consumed and may be past its expiry (no P1), and
+// the stored counter is already the artifact's (no W15); a revocation since then is refused.
+func TestRecheckVerb(t *testing.T) {
+	judge := func(t *testing.T, ro *store.Credentials, in VerbJudgeInput, a []byte) (*approval.SignedVerb, error) {
+		return RecheckVerb(boundedTestContext(t), ro, in, a)
+	}
+	after := vjIn
+	after.Consumed, after.Now, after.BytesDigest = true, vjIn.ExpiresAt+5, [32]byte{}
+	ssh := vjSign(t, jgSigner, approval.VerbNamespace, vjIn)
+	asr := vjAssert(t, jgPasskey, vjIn, 2)
+	moved := `UPDATE webauthn_credentials SET sign_count = 2 WHERE name = 'mark-mac'`
+	otherDigest := after
+	otherDigest.Digest[0] ^= 1
+	vjRun(t, judge, []vjCase{
+		{name: "VJ-parse neither", in: after, artifact: []byte("typed y\n"), step: JudgeStepParse, cause: ErrUnknownArtifact},
+		{name: "VJ-parse an empty artifact", in: after, artifact: nil, step: JudgeStepParse, cause: ErrUnknownArtifact},
+		{name: "VJ-store ssh", setup: vjDropV4, in: after, artifact: ssh, step: JudgeStepStore},
+		{name: "VJ-store passkey", setup: vjDropV5, in: after, artifact: asr, step: JudgeStepStore},
+		{name: "VJ-unknown-key ssh", in: after, artifact: vjSign(t, approvaltest.NewSigner(68), approval.VerbNamespace, vjIn), step: JudgeStepUnknownKey,
+			cause: ErrNotEnrolled, field: approvaltest.NewSigner(68).PublicKey().Fingerprint()},
+		{name: "VJ-unknown-key passkey", in: after, artifact: vjAssert(t, approvaltest.NewPasskey(approval.AlgES256, 69), vjIn, 2), step: JudgeStepUnknownKey,
+			cause: ErrNotEnrolled, rawID: approvaltest.NewPasskey(approval.AlgES256, 69).CredentialID()},
+		{name: "VJ-stored-key", setup: `UPDATE approvers SET public_key = X'00' WHERE name = 'mark-mac'`, in: after, artifact: asr, step: JudgeStepStoredKey, approver: "mark-mac"},
+		{name: "VJ-law ssh revoked", setup: `UPDATE approvers SET revoked_at = 1601 WHERE name = 'mark-ssh'`, in: after, artifact: ssh, step: JudgeStepLaw, cause: approval.ErrApproverRevoked},
+		{name: "VJ-law passkey revoked", setup: moved + `; UPDATE approvers SET revoked_at = 1601 WHERE name = 'mark-mac'`, in: after, artifact: asr, step: JudgeStepLaw,
+			cause: approval.ErrApproverRevoked},
+		{name: "VJ-law ssh another digest", in: otherDigest, artifact: ssh, step: JudgeStepLaw, cause: approval.ErrBadSignature},
+		{name: "VJ-law passkey another digest", setup: moved, in: otherDigest, artifact: asr, step: JudgeStepLaw, cause: approval.ErrChallenge},
+		{name: "VJ-law mint namespace", in: after, artifact: vjSign(t, jgSigner, approval.Namespace, vjIn), step: JudgeStepLaw, cause: approval.ErrNamespace},
+	})
+	if v := vjOK(t, judge, "", after, ssh); v.ApprovedBy() != "ssh:mark-ssh" {
+		t.Fatalf("ssh: %q", v.ApprovedBy())
+	}
+	if v := vjOK(t, judge, moved, after, asr); v.ApprovedBy() != "webauthn:mark-mac" {
+		t.Fatalf("passkey, counter moved: %q", v.ApprovedBy())
+	}
+	vjOK(t, judge, `UPDATE approvers SET enrolled_at = 1300 WHERE name = 'mark-ssh'`, after, ssh) // P5 at the caller's now (1605), not the row's creation (1000)
+}
````

