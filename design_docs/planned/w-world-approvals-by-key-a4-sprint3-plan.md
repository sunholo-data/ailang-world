# Sprint plan — w-world-approvals-by-key-a4, sprint 3 "the passkey arm", headed by the sprint-2 judge's gaps (iteration 391, row 201)

> **THE SPRINT — A4 sprint 3 (design §0 "Sprint 3", §2.3, §2.7, §2.8, §2.10, §6 "Sprint 2"/"Sprint 3" rows that
> moved here, §9 "Sprint 3"), headed by the sprint-2 plan's Appendix C:** planner obligation **P1** measured first
> (Chrome rebuilds the 91-byte decision message and its hash byte-equal to Go); the sprint-2 judge's test gaps
> **JM22, K0's approver operand, JM01–JM04, JM46**; the **webauthn arm** — `DecideByAssertion` (the A2 assertion
> law's W1–W14 through one shared function, **W15 deferred**), `SignedDecision.PasskeyUse`,
> `JudgeDecisionAssertion`, `RecordKeyDecision`'s **K4** (`usePasskey` in the transaction, off the exact-match
> retry only), the POST's **`webauthn`** grammar member and arm (KD2, KD7, KD8, RD10/RD11 passkey halves, RD13
> (b)/(c), RD14); the CLI's **CD6 `--write-message`, CD7 `--signature --expiry`, CD8 `--webauthn-link`, CD9
> `--webauthn-assertion --expiry`** on a request ref (CD6–CD7 re-derived from sprint 2's `CD67-CUT` on M5c as
> landed); a **live-Chrome** sibling (Chrome's virtual authenticator signs a decision through the unchanged A2
> passkey page, the CLI posts it, the daemon's law decides it); **QUICKSTART §18**.
> - **Measured: 1,493 changed lines** — `git diff --numstat 7004f591 HEAD -- cmd host docs`, insertions +
>   deletions, after a commit of each of the seven scratch stages (**376 production, 1,088 test, 29 docs**;
>   `-- cmd host` alone: **1,464**). Margin to the 2,000 cap: **507**.
> - **The design's priced fallback fired, and further (S3-D1):** the card page **and** the rest of the surface
>   move to a fourth sprint — `requests.json` `decisions` (DMD1–DMD3), *Needs you* (XD3), `GET /surface/decide/{ref}`
>   + `decide.js` (PD1, PD2, XD1, XD2, XD4), `approve --list` decisions (design Q2) and `DecisionSubject.Passkeys`
>   (its only consumer is the card). Not for the line cap (507 left) but for the time cap and coherence: the
>   surface is one ≈ 1,250-line unit whose *Needs you* line should link to the card it ships with, not to a CLI
>   line re-pointed a sprint later; every arm this sprint ships is **reachable end to end by an operator from
>   the CLI**, the passkey arm included (§3.1, S3-D2: `--webauthn-link` → the A2 passkey page → `--webauthn-assertion`).
>   Never cut and built: P1, the J-gap tests, KD2, RD13 (b)/(c), RD14.
> - **P1 PASSED (V-P1, §1.2 C4):** Chrome 157.0.8103.51 headless, a page under exactly `surfaceCSP` on
>   `http://localhost`, rebuilt DM1's two golden decision messages (91 bytes) and their SHA-256 with
>   `TextEncoder`, `DataView.setBigUint64` and `crypto.subtle.digest`: **`msg=EQUAL hash=EQUAL` for both**, one
>   subject built from the payload bytes' SHA-256 (as the card will), one from hex with an expiry above 2^53
>   (so `decide.js` must use `BigInt`, recorded for sprint 4).
> - **Eight premises of the design or of sprint 2 are false or incomplete at HEAD and are routed minimally**
>   (§1.4): the A2 passkey page signs whatever challenge its link carries and always reads "Approve request …"
>   (P1, the operator path S3-D2 rests on it, `passkey.js` unchanged); sprint 2's tests pin its own cut and are
>   edited (P2); `DecideBody` with a `json.RawMessage` breaks a sprint-2 test's struct comparison (P3, routed by a
>   pointer); the W15 refusal moves from the judge to the writer with the store's own error text (P4); and four
>   smaller ones (P5–P8).
> - **Every line of code, test and doc is a unified diff in Appendix P**, generated from compiled, vetted,
>   gofmt-clean, green scratch stages; extracted **from this file** with §2's `awk`, applied to a fresh
>   `git archive 7004f591` export and compared byte-for-byte with the stage tree and the mutant tree (§1.2 C20).
>   The executor applies each patch and runs its acceptance; it writes no code by hand.

- **Design**: [w-world-approvals-by-key-a4.md](w-world-approvals-by-key-a4.md), **r2** (reviewed in iteration 388;
  **not reopened**). **D-WORLD-115** (the decision message with a derived challenge) is OPEN; its default **A = keep
  it** is applied (S3-D8). Where this plan and the design differ, this plan was measured against `7004f591` and
  wins; each difference is a decision (§1.3) or a premise found false (§1.4).
- **Starting ground**: A4 sprint 1 landed (PR #563, merge `338ff821`), sprint 2 landed (PR #566, merge
  `9ea7b45f`; judge sonnet PASS 88). Sprint-2 plan: [w-world-approvals-by-key-a4-sprint2-plan.md](w-world-approvals-by-key-a4-sprint2-plan.md)
  (its Appendix C is this sprint's head; its judge's non-blocking findings 1–4 are M6a).
- **Planner base**: `7004f5913d85cc1f307068a7f4c4e4e31aa190ea` (= origin/dev). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter391-a4-s3`, branch `sprint/world-iter391-row201-a4-s3`.
  Planner scratch (outside the repo, never under `/tmp`): `$HOME/.ailang/state/world-iter391/plan-scratch/`
  (`base/` the pristine export, `work/` the tree the mutants were cut from (tag `FINAL`), `work2/` the stage repo
  (tags `M6a` … `M6g`), `patches/`, `tools/` (`gates.sh`, `accept.sh`, `accept_all.sh`, `stages.sh`,
  `mutants.py`, `mutate.py`, `roundtrip.sh`, `assemble.py`), `logs/`, `roundtrip/`, `mut/`, `p1/` (the P1 probe
  and its DOM dump)).
- **Executor**: as routed by the controller. **Three runs**: Run 1 (M6a, M6b, M6c — the gaps, the law, the
  judge and the writer's counter), Run 2 (M6d, M6e — the POST's passkey arm and the CLI's artifact flags), Run 3
  (M6f, M6g — the live-Chrome sibling and QUICKSTART §18). Each sub-milestone is one `patch -p1` (§2) plus its
  acceptance (§3). **Apply, never retype.**
- **Touched** (§5): `host/approval/{passkey_assert.go, decision_signed.go}` (+ new `decision_assert_test.go`),
  `host/authority/judge_decision.go` (+ new `judge_decision_assert_test.go`), `host/store/key_decisions.go` (+ new
  `key_decisions_jgaps_test.go`, `key_decisions_passkey_test.go`), `host/surfacepage/decide_model.go`,
  `host/daemon/{surface_decide.go, surface_decide_key.go}` (+ new `surface_decide_passkey_test.go`,
  `surface_decide_race_test.go`), `cmd/ailang-worldd/{approve.go, approve_decide.go}` (+ new
  `approve_decide_artifact_test.go`, `approve_decide_chrome_test.go`, `approve_decide_jgaps_test.go`,
  `quickstart_decide_test.go`), edited sprint-2 tests `cmd/ailang-worldd/{approve_decide_key_test.go,
  approve_decide_rig_test.go}` (P2), `docs/QUICKSTART.md` (§18). **Outside the design's §8 file list**:
  `host/approval/passkey_assert.go` (the refactor the design's §2.3 names, behaviour-preserving: the A2 matrix
  runs unedited), `host/approval/decision_signed.go`, `host/authority/judge_decision.go`,
  `host/store/key_decisions.go`, `host/daemon/surface_decide_key.go`, `host/surfacepage/decide_model.go`,
  `cmd/ailang-worldd/approve_decide.go` (sprint 2's own files) and every `_test.go` above.
- **Not touched**: `host/approval/{request,signed,sshsig,passkey}.go`, `host/authority/judge.go` (the mint judge),
  `host/store/{passkeys.go, approvers.go, credentials.go, credentials_schema.sql}` (no schema 9),
  `host/broker/**`, `host/coordinator/**`, `host/daemon/{daemon.go, episode_finish.go, surface_approve*.go,
  surface_page.go}`, `host/surfacepage/{approve_model.go, app/**, approve/**, passkey/**}` (`passkey.js`,
  `approve.js`, `app.js` unchanged), `host/sshsign/**`, `cmd/ailang-worldd/{main.go, session.go, help_test.go}`,
  `world/*.ail`, `packages/**`, `scripts/**`, `tools/launchd/*`, `cmd/world-publish/*`, any skill, `~/.ssh`,
  `~/.ailang/world`, `~/.ailang/state/mission-v1*`.

---

## 0. Binding constraints (each is restated where it governs)

1. **Test keys and temp stores only (D-WORLD-109 = A).** Every store is under `t.TempDir()`; every passkey is
   `approvaltest.NewPasskey(alg, seed)` or Chrome's **virtual** authenticator in a temp `--user-data-dir` (the A2
   driver, `testdata/passkey_cdp.mjs`); every SSH key is `approvaltest.NewSigner(seed)` or, for CD7's real-tool
   sibling, an ed25519 key made by `ssh-keygen -q -N ""` in a `t.TempDir()` under a temp `HOME` with
   `SSH_AUTH_SOCK=`. The CLI tests talk to an `httptest` loopback daemon; JM46's test points a **temp `HOME`'s**
   service record at it and **asserts `defaultServiceAddr()` resolves to it before running** (§1.4 P7: without
   that guard a broken record falls through to `127.0.0.1:7644`, the real service). No real store, `~/.ssh`,
   `/dev/tty`, browser profile, mint or publish.
2. **Pure core, effects at the boundary (coding standards S1–S3).** `host/approval` stays pure and
   standard-library only (`TestApprovalImportsOnlyTheStandardLibrary`); `now` is an input; `DecideByAssertion`
   reads no clock. The judge reads on a read-only handle; the counter (W15) is judged **in the writer's
   transaction** (`RecordKeyDecision` K4 → `usePasskey`, unchanged), which the daemon opens under `humanMu` and
   closes after `LandKeyDecision` (`TestDaemonCredentialWriterSites` unchanged: two sites). No `.ail`, no package.
3. **`decidedBy` is derived, never accepted.** The strict grammar still refuses `decidedBy` (RD11 on a
   `webauthn` body); the passkey arm passes `SignedDecision.DecidedBy()` (`webauthn:<name>`).
4. **No `/v1` route, no new route at all, no CSP change, no CORS header, no new Go dependency, no schema change.**
   `TestFrozenV1RouteTableMatchesMux` stays 13; `TestSurfaceAppRoutes` unchanged (the card's +2 move with it).
5. **One mutant per new production branch, anchored to the diff** (§4), each run with its named test and its
   package suite (inventory and static tests excluded); named equivalents and survivors with reasons; every run
   appends a ledger line; the table's ids are the ledger's ids.
6. **No `t.Parallel()`** in any new test (`TestNoParallelWriteForkPackagesOutsideVerifygate`).
7. **Bisectable and safe at every commit.** Each sub-milestone ends gofmt-clean, `go build ./...` and `go vet`
   green, the `_test.go` compile fence green, its own tests and every earlier one green (§1.2 C18).
8. **No git writes by the executor** (§7). The controller commits from snapshots.
9. **The A2/A3 mint law is not re-opened.** `ApproveByAssertion` keeps P1–P5, P3w, W1–W15 in order; the design's
   §2.3 refactor moves W1–W14 into `assertionSteps` and keeps W15 in `ApproveByAssertion`; every A2 test
   (`TestApproveByAssertion*`, `TestChromeAssertion*`, `TestSignedApprovalPasskeyUse`, the PA matrix) runs
   **unedited** and green; `judge.go`, `passkeys.go` and `usePasskey` are unchanged.

---

## 1. Pre-flight

### 1.1 Baseline at `7004f591` (measured by the planner, 2026-11-13)

Every command ran with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` (AILANG v0.41.0), `go1.26.6 darwin/arm64`,
in a `git archive 7004f591` export at `plan-scratch/base` (a throwaway git repo; `diff -r cmd` and `diff -r host`
against the worktree: identical, `BASE_EQUALS_WORKTREE`; the worktree itself was only read). The pre-commit legs
are `tools/gates.sh <tree> <logs> base` = §6's script with the sprint's new-test lists emptied. Load average
9–10 on 16 cores at the start (`uptime` 01:35Z), 20–45 while this planner's own acceptance and mutant runs ran
alongside. Logs: `plan-scratch/logs/base_gates/` (`legs.txt`, one log per leg).

| Leg (`tools/gates.sh <base export> … base`) | rc | Time | Result |
|---|---|---|---|
| `git -C <worktree> status --porcelain \| wc -l` | 0 | — | `0` (clean at `7004f591`) |
| `verify_ail` (`./scripts/verify_ail.sh`) | **0** | 82 s | `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` |
| `gofmt`, `build`, `vet`, `fence` | **0** | 1 / 11 / 10 / 35 s | empty / clean / clean / no FAIL line |
| `pkgs` (`host/approval/...`, `host/authority`, `host/store`, `host/sshsign`, `host/surfacepage`, `host/daemon`, `cmd/ailang-worldd`) | **1** | 546 s | every package `ok` but `host/daemon`, red **only** on `--- FAIL: TestSessionSelfGit` ("local git metadata timed out (2s budget)", the directive's known list, load ≈ 45); **re-run alone: `ok host/daemon 6.442s`** (`logs/base_sessionselfgit_alone.log`); `cmd/ailang-worldd 235.1s` |
| `broker_named` | **0** | 47 s | every named test PASS |
| `broker_whole_reported` (reported, not gating) | **1** | 597 s | only `--- FAIL: TestExecProtectedCheckTimeout (45.05s)` (known list; the same red as iteration 390's base) |
| `race` (three `-race` runs) | **0** | 252 s | `ok` ×6 packages, no race report |
| `inv` (every named inventory, `-v`, base mode: no sprint test) | **0** | 78 s | **49** named `--- PASS`, 0 SKIP, 0 FAIL, none missing |
| `chrome` (`WORLD_CHROME=… WORLD_REQUIRE_PASSKEY_BROWSER=1`, the six A2/A3 page tests) | **0** | 20 s | `TestApprovePageInChrome`, `…ShowsAHostileValueAsText`, `…RefusesASwappedRequest`, `…SSHSIGInChrome`, `TestPasskeyPageAbortsAPendingCeremony`, `TestPasskeyPageInChrome`: `--- PASS` |
| `inv2` (surfacepage/daemon asset and route inventories: `TestSurfaceAppPure`, `TestSurfaceScriptPure`, `TestSurfaceFontEmbed`, `TestSurfaceScriptAssetSeams`, `TestApprovePageRoutes`, `TestPasskeyPageRoutes`, `TestApprovalRoutesRefuseForeignHost`, `TestGETRoutesRejectOtherMethods`, `TestDaemonCallerInventoryAliasControl`; added to §6 after this run began, run alone with `tools/inv2.sh`) | **0** | — | **9** `--- PASS`, 0 SKIP/FAIL (`logs/inv2_base.log`) |
| whole script | — | 1,679 s | `GATES_DONE 1679s FAILS: pkgs broker_whole_reported` (each a known-list red, the first green alone) |
| `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` (after commit) | — | — | **not measured by this planner** (it reads committed blobs; the planner makes no commit in the worktree); the controller's after-commit leg (§6) |

**Rule (carried):** a red at the base measures the repo, not the sprint. The controller re-runs a red untouched
package alone once and does not chase it. **Never start a gate under `nohup`.** `VERIFY_GO_LEGS=plain
./scripts/verify_go.sh` reads committed blobs and runs `go test ./...`; it is the after-commit leg (§6).

### 1.2 Codebase claims checked (at `7004f591` unless "work" — the planner's final scratch tree)

| # | Claim | Command → observed output |
|---|---|---|
| C1 | The base | `git rev-parse HEAD` (worktree) → `7004f5913d85cc1f307068a7f4c4e4e31aa190ea`; `git status --porcelain \| wc -l` → `0` |
| C2 | Controller-verified anchors hold | `grep -n "const credentialsSchemaVersion = " host/store/credentials.go` → `44:const credentialsSchemaVersion = 8`; `grep -n "^func SignIn" host/sshsign/sshsign.go` → `50:`; `grep -n -i -c webauthn host/daemon/surface_decide.go` → `0` (the key known-positive in the same call: `grep -n "GET /surface/requests.json" host/daemon/daemon.go` → `1138:`); `ls host/surfacepage/app/` → `ailang-logo.svg app.css app.js index.html`; `grep -n "^## 17" docs/QUICKSTART.md` → `1841:`; `$HOME/.pinned-ailang/ailang --version` → `AILANG v0.41.0` |
| C3 | **Drift since the design's base is exactly sprints 1–2 and their records** | `git diff --stat b0fbcd3e 7004f591 -- cmd host \| tail -1` → `55 files changed, 3889 insertions(+), 35 deletions(-)` (sprint 1's 1,982 + sprint 2's 1,983 lines, plus sprint 2's fixture edits); the design's anchors this sprint reads re-measure: `passkey_assert.go:88` W15 `if (ad.SignCount != 0 \|\| p.SignCount != 0) && ad.SignCount <= p.SignCount` (design: `:88`), `passkeys.go:226` `usePasskey`'s `UPDATE webauthn_credentials SET sign_count = ? … (sign_count < ? OR (sign_count = 0 AND ? = 0))` (design: `:226`) |
| C4 | **P1 (design §1.12, V-P1): Chrome builds the 91-byte message and its hash byte-equal to Go** | `plan-scratch/p1/srv/main.go` (a throwaway Go server: `GET /p1`, `/p1.js`, `/goldens.json`, each with exactly `surfaceCSP` from `surface_page.go:19`) on `127.0.0.1:0`, opened as `http://localhost:<port>/p1` by `"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" --headless=new --user-data-dir=<mktemp -d under plan-scratch/p1, removed after> --virtual-time-budget=10000 --dump-dom`; the page computes `P`/`Q` with `crypto.subtle.digest('SHA-256', …)` of the base64url payloads (`packet-1`, `request-1`) for subject 1 and from hex for subject 2, assembles `TextEncoder('world-decision/v1') ‖ 0 ‖ verb ‖ P ‖ Q ‖ setBigUint64(83, BigInt(expiry))`, hashes it, and compares with DM1's goldens (`decision_test.go:28-33`) → `DM1-subject-1 … len=91 msg=EQUAL hash=EQUAL`, `DM1-subject-2 (… expiry 0x0102030405060708 > 2^53) len=91 msg=EQUAL hash=EQUAL`, `P1_DONE` (`p1/P1_RESULT.html`). Chrome `157.0.8103.51`. A first run with a non-deferred script printed `pending` (the script ran before `<pre id=out>` existed): a probe bug, fixed by `defer`, recorded |
| C5 | **The A2 passkey page signs the challenge its link carries, whatever it is, and always says "Approve"** (§1.4 P1) | `passkey.js:54` `const challenge = unb64u(f.challenge);`, `:62` `navigator.credentials.get({ publicKey: { challenge, rpId: f.rp,`; `:108` `save.download = … \`world-assertion-${f.request}.json\``; `:131` `` `Approve request ${f.request}, expires …, message hash ${f.hash16}` ``; `:118` refuses only a non-v1 or non-create/get link and `:124` a foreign origin/rp. work: `TestApproveDecideWebAuthnLinkInChrome` (live Chrome) — the CLI's decision link opens to `state ready`, `What` contains `decide-approve-<16 hex>` and `message hash <16 hex of the decision hash>`, and the saved `world-assertion-decide-approve-<16 hex>.json` decides through `DecideByAssertion` |
| C6 | **The JM22 gap** | `key_decisions.go:143` `UPDATE key_decisions SET decision_ref = ? WHERE request_ref = ? AND decision_ref IS NULL`; every sprint-2 store test holds one pending row (judge finding 1); work: `TestLandKeyDecisionLandsOnlyItsRequest` (two unlanded rows; ledger `JM22-FILTER`) |
| C7 | **K0's approver/format operands at the store API** (judge finding 2) | `key_decisions.go` K0 `old.Format == row.Format && old.Approver == row.Approver`; K2 `row.Approver != use.Format+":"+use.Name` → the format operand is implied by the approver operand once K2 passed (named equivalent `K0-FORMAT`); work: `TestRecordKeyDecisionExactMatchNeedsTheSameApprover` (ledger `K0-APPROVER`) |
| C8 | **The JM01–JM04 gap** | `surface_decide_key.go:22 func isApproverRace` four `errors.Is` operands; unreachable through the POST (the read-only judge refuses such an approver first; sprint-2 survivor `D-RACE`); work: `TestIsApproverRaceOperands` (unit; ledger `JM01`–`JM04`) |
| C9 | **The JM46 gap** | `approve_decide.go:117 if o.daemon == "" {` → `defaultServiceAddr()` (`service_pointers.go:28`: the home service record, else `daemon.DefaultAddr` = `http://127.0.0.1:7644`, `daemon.go:169`); `readServiceRecord` (`service.go:258-273`) requires `schema world/service/v1`, the label, the plist and binary paths and a loopback address; work: `TestApproveDecideDefaultDaemonIsTheService` (ledger `JM46`) |
| C10 | **JM27 is an equivalent mutant, not a gap** | `grep -n "LastUsedAt" host/approval/*.go \| grep -v _test` → `approver.go:33: LastUsedAt int64 // 0: never used` only: no law reads it (the judge passes it through; ledger `J-A-LASTUSED`, `E`, the passkey judge's twin) |
| C11 | **W15 and `usePasskey` are the same rule** (S3-D3) | W15 refuses `(ad != 0 \|\| p != 0) && ad <= p`; `usePasskey` updates only when `sign_count < new OR (sign_count = 0 AND new = 0)` — the complement of the same predicate, evaluated on the **stored** counter inside the transaction (C3's two anchors) |
| C12 | **The sprint-2 cut is pinned by two sprint-2 tests** (§1.4 P2) | `approve_decide_key_test.go:68-72` (`--signature`, `--write-message` with a request ref → "the other key flags do not decide it"), `approve_decide_rig_test.go:169` (`{f.ref, "--webauthn-link"}` → "there is nothing to await"); `approve.go:219` the refusal text both read |
| C13 | **`DecideBody` must stay comparable** (§1.4 P3) | `approve_decide_test.go:50 if b != want \|\|` compares two `surfacepage.DecideBody` values; work: with `WebAuthn json.RawMessage` the cmd test package does not compile (`invalid operation: b != want (struct containing json.RawMessage cannot be compared)`); with `*json.RawMessage` it does and passes unedited |
| C14 | The A3a body's member is a raw JSON assertion envelope | `approve_model.go:185 WebAuthn json.RawMessage \`json:"webauthn"\``; the decision body uses the same JSON member name and envelope (A2 §2.7), as a pointer (C13) |
| C15 | `host/approval` stays standard-library only | work: `decision_signed.go` imports `bytes` only; `passkey_assert.go` unchanged imports; `TestApprovalImportsOnlyTheStandardLibrary` PASS |
| C16 | No new context root | work: `decideByArtifact` and `decisionPasskeyLinks` run under `runApproveDecide`'s existing 10 s root (`approve_decide.go`, pinned `runApproveDecide\|Background`); `TestProductionContextRoots` PASS unedited |
| C17 | No new credentials method, writer site or route | work: `git diff 7004f591 FINAL -- host/store/credentials.go host/daemon/daemon.go host/store/passkeys.go \| wc -l` → `0`; `TestCredentialsMethodsRequireADeadline`, `TestDaemonCredentialWriterSites`, `TestSurfaceAppRoutes`, `TestFrozenV1RouteTableMatchesMux` PASS unedited |
| C18 | **Every sub-milestone is green when committed** | `tools/accept_all.sh` (a clone of the stage repo checked out at each tag in turn; `logs/accept_<M>.log`): `ACCEPT_OK` for **every one of the seven** (`M6a` … `M6g`; run 1 on the stage repo before the last two `TestApproveDecideWriteMessageAndSignature` row edits, and `M6e`, `M6f`, `M6g` again on the final stage repo, `logs/accept2_<M>.log`): each stage gofmt-clean, `go build ./...` and `go vet` of the touched packages clean, the `_test.go` compile fence without a FAIL line, its §3 `-v` tests `--- PASS` with no SKIP/FAIL (M6e's real-`ssh-keygen` sibling ran; M6f's live Chrome ran with its env: `--- PASS: TestApproveDecideWebAuthnLinkInChrome`, `--- PASS: TestPasskeyPageInChrome`), and its package runs `ok` |
| C19 | **Every sub-milestone is committed and measured** | the round-trip tree commits each patch in turn: `NUMSTAT M6a 124`, `M6b 381`, `M6c 628`, `M6d 963`, `M6e 1328`, `M6f 1411`, `M6g 1493 (cmd host: 1464)` (`logs/roundtrip_final.log`) |
| C20 | **The patches reproduce the stage tree and the mutant tree** | `tools/roundtrip.sh <this plan>`: a fresh `git archive 7004f591` export + `patch -p1` of Appendix P's seven patches **extracted from this file with §2's `awk`** (each `cmp`-equal to the generated patch; `--dry-run` first; no `.orig`/`.rej`) → `diff -r` of `cmd`, `host`, `docs` against `plan-scratch/work2` (`M6g`) and of `cmd`, `host` against `plan-scratch/work` (`FINAL`, the mutant tree): **identical** (`ROUNDTRIP_IDENTICAL_TO_STAGE_TREE_AND_MUTANT_TREE`); `work2` is built from `work` by `tools/stages.sh`, which ends `WORK2_EQUALS_FINAL` (`diff -r cmd host docs`) |
| C21 | **Nothing outside §5 moved** | stage repo: `git diff --stat <base> M6g -- world packages scripts tools design_docs \| wc -l` → `0`; production: `git diff --stat <base> M6g -- '*.go' ':!*_test.go'` → `9 files changed, 315 insertions(+), 61 deletions(-)`; added lines matching `t.Parallel\|Access-Control\|OpenCredentials(\|context.Background()\|/dev/tty\|.ssh` under `cmd host` → four, all in `_test.go` (two test writers and two test contexts in `approve_decide_artifact_test.go`/`approve_decide_chrome_test.go` building the fake daemon's credentials store); none in production, no `t.Parallel`, no CORS header |
| C22 | **Every mutant, on the final tree; the table's ids are the ledger's** | 76 rows, the last ledger line of each: 65 K, 5 U, 4 E, 2 L, 0 not as expected, 0 `NOCOMPILE` (§4, `logs/mutant_ledger.tsv`, generated by `tools/mtable.py`); `comm -3 <(cut -f4 ledger \| sort -u) <(the table's ids)` → empty (`logs/ledger_ids.txt`, `logs/plan_ids.txt`) |
| C23 | **The gates, on the final tree** (`plan-scratch/final-tree` = `work2` at `M6g`, `diff -r cmd host docs` equal) | `tools/gates.sh <final-tree> logs/final_gates final`: `verify_ail` → `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` (no `.ail` changed); `gofmt`, `build`, `vet`, `fence` rc 0; `pkgs` rc 0 (every package `ok`: `host/daemon 430.9s`, `cmd/ailang-worldd 245.7s`); `broker_named` rc 0; `broker_whole_reported` rc 0 (`ok host/broker 230.0s`: the base's known-list red did not recur); `race` rc 0 (`ok` ×6: host/approval, host/authority, host/store, host/sshsign, host/daemon (the decide, approve and finish tests), cmd/ailang-worldd (the approve tests); no `DATA RACE`); `inv` rc 0: **67** named `--- PASS` (49 as at the base + 18 new), 0 SKIP, 0 FAIL, none missing — incl. `TestFrozenV1RouteTableMatchesMux` (13), `TestSurfaceAppRoutes`, `TestDaemonCredentialWriterSites` (two sites), `TestProductionContextRoots`, `TestCredentialsMethodsRequireADeadline`, `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist`, `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`, `TestEveryVerbHasHelp`, `TestApprovalImportsOnlyTheStandardLibrary`, the quickstart tests — all unedited; `inv2` rc 0 (9 `--- PASS`; bash read the leg from the edited script mid-run, and `tools/inv2.sh` on the same tree agrees: `logs/inv2_final.log`); `chrome` rc 0: `TestApproveDecideWebAuthnLinkInChrome` and the six A2/A3 page tests `--- PASS`; **`GATES_DONE 1001s FAILS: none`** (`logs/final_gates/legs.txt`). One test file (`approve_decide_artifact_test.go`, two rows) was synced into the gate tree at the start of `pkgs` (the R3 fix); `inv`, `race` and the acceptance re-run (`logs/accept2_M6e.log`) ran it in its final form | |

### 1.3 Planner decisions where the design leaves a choice

- **S3-D1 — The cut: the whole surface moves to sprint 4; the passkey arm, CD6–CD9, the gaps, P1 and §18 ship.**
  Measured: this set is **1,493** (`-- cmd host docs`). The design's sprint-3 fallback keeps the listing, the
  *Needs you* line "pointing at `ailang-worldd approve <ref>`" and QUICKSTART, and moves the card. This plan keeps
  QUICKSTART (§18 documents the CLI path that ships) and moves the listing and *Needs you* with the card, because
  (a) the directive's time cap (≈ 100 min) does not hold a proven listing, `app.js` block and its live-Chrome
  XD3 on top of the passkey arm; (b) a *Needs you* line pointing at a CLI command would be re-pointed to the card
  one sprint later (churn in `app.js`, an A3a-era 1,470-line file); (c) the operator does not need the surface to
  reach any shipped arm (S3-D2). Priced in Appendix C: ≈ 1,250. Never cut, built: P1, JM22/K0/JM01–04/JM46, KD2,
  RD13 (b)/(c), RD14, KD7/KD8.
- **S3-D2 — How an operator produces a passkey decision without the card (the directive's interaction).**
  `approve <request-ref> [--deny] --webauthn-link` prints, per origin of the active passkeys, an A2 `get` link
  (`/surface/passkey#…`, unchanged page) whose `challenge` is `SHA-256(DecisionMessage(verb, P, Q, window.max))`,
  `request` = `decide-<verb>-<16 hex of Q>`, `hash16` = the hash's first 8 bytes; the page signs it with the
  device's check and saves `world-assertion-decide-<verb>-<16 hex>.json`; `approve <request-ref> [--deny]
  --webauthn-assertion <file> --expiry <t>` posts it as `webauthn`. Proven end to end with Chrome's own bytes
  (`TestApproveDecideWebAuthnLinkInChrome`, M6f). So no POST arm ships without an operator path (§3.1). The page's
  "Approve request decide-deny-…" wording for a deny is §1.4 P1.
- **S3-D3 — W15 lives in the writer as `usePasskey`, not as a second Go check.** The design's "W15 + `usePasskey`
  in `RecordKeyDecision`'s transaction on the non-exact path" is one rule (C11): `usePasskey`'s SQL is W15
  evaluated on the stored counter inside the transaction. K4 = `if !exact && use.Passkey != nil { usePasskey }`;
  the exact-match path skips only it (`useApprover` still runs, KD9's rule). A refused counter answers 422
  `NotApproved` with the store's text (`store: the passkey's sign counter did not advance …`; §1.4 P4).
- **S3-D4 — K2 gains the passkey rule:** a `webauthn` row needs `use.Passkey`, an `ssh` row must not carry one
  (A2's C13 for the consume), so the counter can never be skipped by a caller that forgot the use.
- **S3-D5 — The POST's `webauthn` member** is `*json.RawMessage` (C13, C14): `null` decodes to nil (counted as
  no format); the arm stores the envelope bytes as `artifact` (K0 compares them byte for byte). The CLI posts the
  saved file through `json.Marshal`, which compacts it deterministically, so an identical re-post is byte-equal.
- **S3-D6 — CD7/CD9 take `--expiry`** (the expiry is in the signed message; the CLI cannot recover it from an
  SSHSIG or an assertion); `--signature`/`--webauthn-assertion` without it, and `--expiry` without one of them or
  with a mint id, are `exitUsage` before any daemon call. CD6/CD8 print the exact next command with `--expiry`,
  `--daemon`, `--db` and `--credentials-db`.
- **S3-D7 — CD8 reads the passkeys from the credentials store read-only**, as the mint `--webauthn-link` does
  (`writePasskeyLinks`), not from the daemon: `DecisionSubject.Passkeys` (sprint-1 P7) has the card as its only
  consumer and moves with it (Appendix C). The CLI has already proved the daemon serves this credentials store
  (sprint 2 S2-D2's health check).
- **S3-D8 — D-WORLD-115.** Under B (a stored daemon challenge), the decision message's P slot carries it; this
  sprint's passkey law hashes the same `DecisionMessage`, so B changes `DecideByAssertion`'s input with
  `DecideBySignature`'s (sprint-2 S2-D9) and CD6/CD8's message, nothing else. Built under **A**.
- **S3-D9 — The real-tool siblings (directive retro 11).** DM7, CD8 and CD9 drive approvaltest's authenticator;
  their sibling is **live Chrome** (`TestApproveDecideWebAuthnLinkInChrome`, gated as XB1: `WORLD_CHROME`,
  `WORLD_REQUIRE_PASSKEY_BROWSER=1`, the A2 driver's temp profile), and the A2 recorded Chrome fixtures still drive
  the mint law unedited. CD6/CD7's sibling is the real `ssh-keygen` signing the file `--write-message` wrote
  (`TestApproveDecideWriteMessageSignedBySSHKeygen`). P1 measures the card's future arithmetic. No new
  `testdata/` capture is checked in: the live test re-records each run (S3-D10's risk, §8).
- **S3-D10 — Sub-milestones** (§1.5): M6e (365) is over the 250–350 band by 15: CD6–CD9 share one function and
  one test file.

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

- **P1 — Design §2.8/§9 "the webauthn arm is needed for XD1" and "a passkey on the card" assume the card is the
  only passkey surface.** The A2 passkey page (`/surface/passkey`, unchanged) signs **any** challenge a link
  carries and saves the assertion (C5). Routed: CD8 builds a decision link for it (S3-D2), so the passkey arm is
  reachable without the card. Its limit, recorded for the operator (QUICKSTART §18) and sprint 4: the page always
  reads **"Approve request <request>"**, so a deny link reads "Approve request decide-deny-…"; the verb is in the
  link's `request` text, the CLI's capitals and the signed hash the page shows. `passkey.js` is not edited
  (design §2.8: "approve.js, /surface/passkey and passkey.js are not changed").
- **P2 — Design §6 "Must stay green, unedited: every `TestApprove*`".** Two sprint-2 tests pin sprint 2's own
  cut, which this sprint reverses (C12): `TestApproveDecideSignWith`'s last loop (5 lines deleted) and one row of
  `TestApproveDecideFlagRules` (the `--webauthn-link` row becomes "`--webauthn-assertion` without `--expiry`").
  Both are A4 tests written in iteration 390; no A2/A3 test is edited.
- **P3 — Design §2.8 "Body, strict, exactly one of `sshsig`, `webauthn`, `typed`" with A3a's `json.RawMessage`
  member** breaks `approve_decide_test.go:50`'s struct comparison (C13). Routed: `*json.RawMessage` (S3-D5).
- **P4 — Design §2.7 "A counter refused here, in the writer, answers 422 with the same `ErrSignCountRegression`
  text the read-only judge uses".** The writer's refusal is `store.ErrSignCountRegression` (`passkeys.go:57`,
  "store: the passkey's sign counter did not advance …"), not `approval.ErrSignCountRegression` ("approval: …"):
  the same sentence after a different package prefix, and the store's carries the stored and posted counters and
  the revoke hint. Routed: 422 `NotApproved` "the passkey assertion does not decide request …: store: …"
  (`isCounterRefusal`), pinned by `TestDecidePOSTPasskey`.
- **P5 — Design §6 RD13 (b) "a fresh assertion at N+1 instead → 200 by replacement"** after a failed world write
  that already moved the counter to N: N+1 lands only if it is above the stored N — true for a real
  authenticator. The test posts N+2 after N+1's failed write (stored N+1) to keep both steps distinct
  (`TestDecidePOSTPasskeyRetryAfterAFailedWorldWrite`: 8 fails, 9 lands).
- **P6 — Two identical `Human.Approve` requests share a ref** (design §1.5, RK2): the daemon tests' packets
  need distinct requesters (`agent-2`, `agent-3`), or D4b refuses (observed in this planner's first run: 409
  "carried by 2 decision packets").
- **P7 — JM46 cannot be tested with an empty `HOME`.** `defaultServiceAddr()` falls back to
  `http://127.0.0.1:7644` when no valid record exists; this planner's **first draft** of the test wrote a record
  missing `schema`/`label`/paths, the CLI fell through and its `GET /v1/health` reached the machine's real
  service on 7644 (read-only; it refused "serves store /Users/…/.ailang/world/world.db, not <temp>"; nothing was
  written). Routed: the record now carries every field `readServiceRecord` requires and the test **fails before
  running the CLI unless `defaultServiceAddr()` names the fake daemon** (§0.1). Recorded as a D-WORLD-109 near
  miss, not repeated.
- **P8 — Design §7 "`TestSurfaceAppRoutes` +2 (sprint 3)"**: unchanged this sprint (the card's routes move,
  S3-D1).

### 1.5 The size

Lines are insertions + deletions: **`git diff --numstat <base> HEAD -- cmd host docs`** (the directive's filter;
`-- cmd host` alone in the last column). **Every row is measured** from the scratch stages, each committed in a
throwaway repo (`plan-scratch/work2`, tags `M6a` … `M6g`), and re-measured on the round-trip tree (§1.2 C19).

| Sub-milestone | Production | Test | Docs | Own | Cumulative vs `7004f591` (after its commit) | `-- cmd host` | Run |
|---|---|---|---|---|---|---|---|
| M6a — JM22, K0 approver, JM01–JM04, JM46 (test only) | 0 | 124 | 0 | 124 | **124** | 124 | 1 |
| M6b — `assertionSteps` (the refactor), `DecideByAssertion`, `PasskeyUse`; DM7 | 111 | 146 | 0 | 257 | **381** | 381 | 1 |
| M6c — `JudgeDecisionAssertion`, K4 + K2's passkey rule; JD (passkey), KD2, KD7 (passkey), K2 | 51 | 196 | 0 | 247 | **628** | 628 | 1 |
| M6d — the POST's `webauthn` member and arm, `isCounterRefusal`; RD10/RD11 (passkey), KD2/KD7 end to end, KD8, RD13 (b)(c), RD14 | 84 | 251 | 0 | 335 | **963** | 963 | 2 |
| M6e — CD6–CD9 (`decideByArtifact`, `decisionPasskeyLinks`, `--expiry`); P2's edits | 130 | 235 | 0 | 365 | **1,328** | 1,328 | 2 |
| M6f — the live-Chrome sibling | 0 | 83 | 0 | 83 | **1,411** | 1,411 | 3 |
| M6g — QUICKSTART §18 and its flag test | 0 | 53 | 29 | 82 | **1,493** | 1,464 | 3 |
| **Sprint 3** | **376** | **1,088** | **29** | | **1,493 ≤ 2,000** (margin 507) | 1,464 | |

Per file, cumulative (+/−): `approve.go` +19 −5, `approve_decide.go` +106, `approve_decide_artifact_test.go` +222,
`approve_decide_chrome_test.go` +83, `approve_decide_jgaps_test.go` +39, `approve_decide_key_test.go` +1 −6,
`approve_decide_rig_test.go` +3 −3, `quickstart_decide_test.go` +53, `QUICKSTART.md` +29,
`decision_assert_test.go` +146, `decision_signed.go` +45 −2, `passkey_assert.go` +39 −25, `judge_decision.go`
+34 −1, `judge_decision_assert_test.go` +83, `surface_decide.go` +3 −3, `surface_decide_key.go` +29 −8,
`surface_decide_passkey_test.go` +251, `surface_decide_race_test.go` +27, `key_decisions.go` +13 −3,
`key_decisions_jgaps_test.go` +58, `key_decisions_passkey_test.go` +113, `decide_model.go` +27 −14. Against the
design's ≈ 500 for the webauthn arm and the sprint-2 plan's 63 for CD6–CD7: the arm with CD6–CD9 came in at 1,204
(M6b–M6e: 376 production, 828 test).

---
## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of four
backticks, in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; a new file's patch has `--- /dev/null`
and an `@@ -0,0 …` hunk. Apply it **only** this way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF2'
set -euo pipefail
M=M6a   # <- the sub-milestone, exactly as in the Appendix P heading: M6a M6b M6c M6d M6e M6f M6g
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter391-a4-s3
PLAN="$W/design_docs/planned/w-world-approvals-by-key-a4-sprint3-plan.md"
T="$HOME/.ailang/state/world-iter391/exec-patches"; mkdir -p "$T"
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
  `$HOME/.ailang/state/world-iter391/exec-logs/` (never only through `tail`). **Never `nohup`, never background a
  test leg.** The planner's `tools/accept.sh <tree> <M>` is exactly §3's lines; it printed `ACCEPT_OK` for every
  stage on the stage repo (§1.2 C18).
- Then snapshot (§7). The controller commits and compares the numstat with §3's expected value.

---

## 3. Milestones

Every sub-milestone's acceptance starts with the same three lines (abbreviated **[base]** below):

```
test -z "$(gofmt -l cmd host)"
go build ./... && go vet ./host/approval/... ./host/authority/ ./host/store/ ./host/sshsign/ ./host/surfacepage/ ./host/daemon/ ./cmd/ailang-worldd/ ./host/broker/
go test -run '^$' ./... 2>&1 | grep -E '^(FAIL|--- FAIL)' && echo FENCE_RED   # must print nothing
```

Each `-v` run below must show `--- PASS` for every test it names and **no `--- FAIL` and no `--- SKIP`**; each
package run must end `ok` with no `FAIL` line.

### M6a — the sprint-2 judge's test gaps · 124 test · Run 1

**Patch** `M6a`: new `host/store/key_decisions_jgaps_test.go` — `TestLandKeyDecisionLandsOnlyItsRequest` (JM22:
two unlanded rows; landing one leaves the other unlanded, then it lands on its own) and
`TestRecordKeyDecisionExactMatchNeedsTheSameApprover` (finding 2: the same artifact and signed fields for another
approver replace the row, never kept as the first approver's). New `host/daemon/surface_decide_race_test.go` —
`TestIsApproverRaceOperands` (JM01–JM04: each of the four wrapped store errors is a race, three other errors are
not). New `cmd/ailang-worldd/approve_decide_jgaps_test.go` — `TestApproveDecideDefaultDaemonIsTheService` (JM46: a
temp `HOME`'s complete service record names the fake daemon; the test asserts `defaultServiceAddr()` resolves to
it **before** running the CLI (§1.4 P7); `approve <ref> --sign-with` with no `--daemon` reaches it: 1 GET, 1 POST).

**Acceptance**
```
[base]
go test ./host/store/ -run '^(TestLandKeyDecisionLandsOnlyItsRequest|TestRecordKeyDecisionExactMatchNeedsTheSameApprover)$' -count=1 -v
go test ./host/daemon/ -run '^TestIsApproverRaceOperands$' -count=1 -v
go test ./cmd/ailang-worldd/ -run '^TestApproveDecideDefaultDaemonIsTheService$' -count=1 -v
```
**Expected cumulative numstat after the controller's commit of M6a: 124** (`-- cmd host docs`; `-- cmd host`: 124).

### M6b — the decision assertion law · 111 prod + 146 test · Run 1

**Patch** `M6b`: `host/approval/passkey_assert.go` — W1–W14 move into `assertionSteps(response, a, h, hErr)
(AuthData, error)` (the design §2.3 refactor; `hErr` is refused at W6, in W6's place, so the order of every
refusal is unchanged); `ApproveByAssertion` computes the mint `MessageHash`, calls it, and keeps W15 and its
`SignedApproval` as before. `host/approval/decision_signed.go` — `SignedDecision` gains `credentialID`,
`signCount`; `DecisionAssertionInput`; `DecideByAssertion` (P2–P5, P3w via `checkApprover`, then
`assertionSteps` over `DecisionMessageHash(verb, P, Q, expiry)`; **no W15**); `PasskeyUse()` (ok only for a
verified webauthn decision). New `host/approval/decision_assert_test.go` — `TestDecideByAssertionSteps` (DM7: P2,
P3w ×2, P4, P5, W1, W2, W3, W5, W6 ×7 incl. the mint message, the other verb, another P/Q/expiry, no message and
the zero hash an absent message would leave, W7–W14 each → its error, nil decision),
`TestDecideByAssertionDefersTheCounter` (counters 2, 1 (= stored) and 0 all verify: W15 is not the law's; the
`PasskeyUse` and `Use` hand-off; a decision assertion as a mint → `ErrChallenge` (DM4 passkey half); zero and SSH
decisions have no passkey use).

**Acceptance**
```
[base]
go test ./host/approval/... -run '^(TestDecideByAssertionSteps|TestDecideByAssertionDefersTheCounter|TestApprovalImportsOnlyTheStandardLibrary|TestApproveByAssertionOrder|TestApproveByAssertionAcceptsChrome)$' -count=1 -v
go test ./host/approval/... -count=1
```
**Expected cumulative numstat after the controller's commit of M6b: 381**.

### M6c — the passkey judge and the writer's counter · 51 prod + 196 test · Run 1

**Patch** `M6c`: `host/authority/judge_decision.go` — `JudgeDecisionAssertion(ctx, ro, in, response)`:
`ParseAssertionResponse` (parse) → `PasskeyByCredentialID` (store; unknown-key with `RawID`) →
`ParseStoredCOSEKey` (stored-key) → `DecideByAssertion` (law), judge.go's step names. `host/store/key_decisions.go`
— K2 gains "a `webauthn` row exactly when `use.Passkey` is set" (S3-D4); **K4**: `if !exact && use.Passkey != nil
{ usePasskey }` between K3 and K5 (S3-D3); comments. New `host/authority/judge_decision_assert_test.go` —
`TestJudgeDecisionAssertion` (parse, store below v6, unknown key with the credential ID, stored key, law: mint
assertion, other verb, revoked, backup eligibility, origin; the control at counters 2 and 1 over stored 1). New
`host/store/key_decisions_passkey_test.go` — `TestRecordKeyDecisionCounter` (KD2: 7 and 6 over stored 7 →
`ErrSignCountRegression`, nothing written or moved; 8 → recorded, counter 8, `last_used_at`; stored 0 + 0 →
recorded), `TestRecordKeyDecisionIdempotentPasskeyRetry` (KD7 passkey half: the exact retry at 7 over stored 7
records without the counter, row kept, `useApprover` ran; another artifact at 7 → refused, row unchanged; 8 →
replaced, counter 8), `TestRecordKeyDecisionPasskeyUseMatchesTheFormat` (K2's rule both ways; nothing written,
counter unmoved).

**Acceptance**
```
[base]
go test ./host/authority/ -run '^(TestJudgeDecisionAssertion|TestJudgeDecisionSSH|TestJudgeAssertionIsTheCLIJudge)$' -count=1 -v
go test ./host/store/ -run '^(TestRecordKeyDecisionCounter|TestRecordKeyDecisionIdempotentPasskeyRetry|TestRecordKeyDecisionPasskeyUseMatchesTheFormat|TestRecordKeyDecision|TestCredentialsMethodsRequireADeadline)' -count=1 -v
go test ./host/store/ ./host/authority/ -count=1
```
**Expected cumulative numstat after the controller's commit of M6c: 628**.

### M6d — the POST's passkey arm · 84 prod + 251 test · Run 2

**Patch** `M6d`: `host/surfacepage/decide_model.go` — `DecideBody.WebAuthn *json.RawMessage \`json:"webauthn,omitempty"\``
(S3-D5), `decideFormats` (exactly one of typed, sshsig, webauthn), the grammar's text. `host/daemon/surface_decide.go`
— the arm slot takes the key arm for `sshsig` **or** `webauthn` (one line + comment). `host/daemon/surface_decide_key.go`
— `keyArm` selects `JudgeDecisionAssertion` and the envelope bytes for `webauthn` (else sprint 2's SSH judge),
names the artifact in its refusals, passes `ApproverUse.Passkey` from `PasskeyUse()` to `RecordKeyDecision`, and
maps `isCounterRefusal` (`store.ErrSignCountRegression`, `store.ErrPasskeyChanged`) to 422 `NotApproved` beside
sprint 2's `isApproverRace`. New `host/daemon/surface_decide_passkey_test.go` — `TestDecidePOSTPasskey` (RD10,
RD11 passkey halves: approve and deny land as `webauthn:mark-mac`, the row landed, the counter moved; `decidedBy`,
`webauthn`+`sshsig`, `webauthn`+`typed` → 400; a mint assertion → 422; KD2 end to end: counter 3 over stored 3 →
422, nothing decided, recorded or moved), `TestDecidePOSTCounterZeroAssertionLandsOnce` (KD8 + RD13 (c)),
`TestDecidePOSTPasskeyRetryAfterAFailedWorldWrite` (RD13 (b) + KD7 end to end: the failed write leaves one
unlanded row and counter 7; another assertion at 7 → 422; the identical body → 200, counter still 7, `decided_at`
kept; a fresh assertion lands by replacement), `TestDecidePOSTRetryAfterAFailedWorldWriteAcrossRestart` (RD14:
both reviewers' scenarios through a restarted daemon on the same stores), `TestIsCounterRefusalOperands`.

**Acceptance**
```
[base]
go test ./host/daemon/ -run '^(TestDecidePOSTPasskey|TestDecidePOSTCounterZeroAssertionLandsOnce|TestDecidePOSTPasskeyRetryAfterAFailedWorldWrite|TestDecidePOSTRetryAfterAFailedWorldWriteAcrossRestart|TestIsCounterRefusalOperands|TestDecidePOSTGuards|TestDecidePOSTSSHSIG|TestDaemonCredentialWriterSites)$' -count=1 -v
go test ./host/surfacepage/ -count=1
go test ./host/daemon/ -run '^(TestDecision|TestDecide|TestSurfaceApprove|TestApprov|TestIs)' -count=1
```
**Expected cumulative numstat after the controller's commit of M6d: 963**.

### M6e — the CLI's artifact flags on a request ref · 130 prod + 235 test · Run 2

**Patch** `M6e`: `cmd/ailang-worldd/approve.go` — `--expiry`; CD3 widened (`--expiry` needs a request ref);
the decision mode refuses only `--await`, and `--signature`/`--webauthn-assertion` without `--expiry` (or
`--expiry` alone) (S3-D6); the options passed; the help (two usage lines, one paragraph).
`cmd/ailang-worldd/approve_decide.go` — `decideOptions` gains the flags; `runApproveDecide` dispatches them to
`decideByArtifact` after the subject checks; `decideByArtifact` (CD6: the 91-byte message for `window.max` to a NEW
file (`O_EXCL`), the `ssh-keygen -Y sign -n world-decision@sunholo` line and the exact `--signature … --expiry`
command; CD8: `decisionPasskeyLinks`; CD7: the capped file as `sshsig` with `--expiry`; CD9: the capped file,
parsed as an assertion before any POST, as `webauthn` with `--expiry`); `decisionPasskeyLinks` (S3-D2, S3-D7). New
`cmd/ailang-worldd/approve_decide_artifact_test.go` — `TestApproveDecideWriteMessageAndSignature` (CD6, CD7 and
the flag rules; missing files), `TestApproveDecideWriteMessageSignedBySSHKeygen` (CD6+CD7 with the real
`ssh-keygen` on a temp key; verifies with `VerifyDecisionSSHSIG`, not as a mint),
`TestApproveDecideWebAuthnLinkAndAssertion` (CD8: no passkey refused; the link's op, origin, rp, challenge =
the decision hash, expiry, `decide-deny-<16>`, allow, hash16; the next command; CD9: an approvaltest assertion of
that challenge posted as `webauthn` and decided by `DecideByAssertion`; a non-assertion refused before the POST).
Edited (§1.4 P2): `approve_decide_key_test.go` (the cut's loop removed), `approve_decide_rig_test.go` (one row).

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApproveDecide|TestEveryVerbHasHelp)' -count=1 -v
go test ./host/store/ -run '^TestProductionContextRoots$' -count=1 -v
go test ./cmd/ailang-worldd/ -run '^(TestApprove|TestApprover)' -count=1
```
(`TestApproveDecideWebAuthnLinkInChrome` does not exist until M6f.) **Expected cumulative numstat after the
controller's commit of M6e: 1,328**.

### M6f — the live-Chrome sibling · 83 test · Run 3

**Patch** `M6f`: new `cmd/ailang-worldd/approve_decide_chrome_test.go` — `TestApproveDecideWebAuthnLinkInChrome`
(S3-D9): the A2 driver (`startPasskeyBrowser`, a temp profile, Chrome's virtual authenticator with UV) creates a
passkey through the A2 page at a `passkeyPageOrigin` daemon; it is enrolled in the fake daemon's credentials store;
`approve <ref> --webauthn-link` prints one link; the page opens `ready`, showing `decide-approve-<16 hex>` and the
decision hash's 16 hex; its saved `world-assertion-decide-approve-<16 hex>.json`, posted by `--webauthn-assertion
--expiry 1760`, is byte-for-byte the POST's `webauthn` and decides through `DecideByAssertion` as
`webauthn:mark-mac`. Without `WORLD_CHROME` it skips, counted; with `WORLD_REQUIRE_PASSKEY_BROWSER=1` a missing
browser fails.

**Acceptance** (gated as XB1, as A3; on the controller's machine; CI has no Chrome step)
```
[base]
WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" WORLD_REQUIRE_PASSKEY_BROWSER=1 go test ./cmd/ailang-worldd/ -run '^(TestApproveDecideWebAuthnLinkInChrome|TestPasskeyPageInChrome)$' -count=1 -v
```
**Expected cumulative numstat after the controller's commit of M6f: 1,411**.

### M6g — QUICKSTART §18 · 29 docs + 53 test · Run 3

**Patch** `M6g`: `docs/QUICKSTART.md` — new **§18 "Decide an agent's approval request *(attended — pending
first verbatim run)*"** (typed `--deny`, `--sign-with`, `--write-message` → `ssh-keygen -Y sign` →
`--signature --expiry`, `--webauthn-link` → the page → `--webauthn-assertion --expiry`, the replay and retry
behaviour; the page's "Approve request decide-…" wording, §1.4 P1). New `cmd/ailang-worldd/quickstart_decide_test.go`
— `TestQuickstartSection18FlagsMatchTheCLI` (every flag on §18's six `approve sha256:` lines is in `approveHelp`;
each decision flag appears; the namespace and the 91-byte length are the code's constants).

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestQuickstartSection18FlagsMatchTheCLI|TestEveryVerbHasHelp|TestQuickstartAttendedStepsBoxMatchesSection9|TestQuickstartSection9FlagsMatchTheCLI|TestQuickstartSection10FlagsMatchTheCLI|TestQuickstartSection11FlagsMatchTheCLI)$' -count=1 -v
go test ./cmd/ailang-worldd/ -count=1
```
**Expected cumulative numstat after the controller's commit of M6g: 1,493** (`-- cmd host`: 1,464).

### 3.1 Reachability: every new refusal and branch, from the entry point (directive retro 4)

Entry points: **POST** = `POST /surface/decide/{ref}` through `d.Handler()` (R0, R3–R6, grammar, `humanMu`,
`readDecision`, law D1–D10, arm, `DecideApproval`, `LandKeyDecision`); **CLI** = `runApprove` → decision mode.
"E2E" = driven by a test through that entry point; "unit" = only by a direct call.

| Branch (file) | Input class → outcome | Trace from the entry point | Driven |
|---|---|---|---|
| grammar `decideFormats == 0` (`decide_model.go`) | no format, `"webauthn":null` → 400 | POST → `ParseDecideBody` | E2E `TestDecidePOSTGuards` "R6 no format" |
| grammar `decideFormats > 1` | `webauthn`+`sshsig`, `webauthn`+`typed` → 400 | POST → `ParseDecideBody` | E2E `TestDecidePOSTPasskey` |
| arm slot `SSHSIG != nil \|\| WebAuthn != nil` (`surface_decide.go`) | a `webauthn` body → `keyArm` | POST → law → slot | E2E `TestDecidePOSTPasskey` |
| `keyArm` judge select (`surface_decide_key.go`) | `webauthn` → `JudgeDecisionAssertion`; `sshsig` → `JudgeDecisionSSH` | POST → slot → `keyArm` | E2E `TestDecidePOSTPasskey`, `TestDecidePOSTSSHSIG` |
| judge `parse` | not an envelope → 422 | POST → `keyArm` → judge | unit `TestJudgeDecisionAssertion`; E2E: the CLI refuses a non-assertion before any POST (CD9), a hand POST reaches it (same 422 as sprint 2's `sshsig` parse) |
| judge `store` | credentials below v6 → 503 `Restart` | POST → `keyArm` → judge | unit only: the daemon upgrades the store at `serve` (sprint-2 precedent for `sshsig`) |
| judge `unknown-key` | a credential not enrolled → 422 | POST → judge | unit `TestJudgeDecisionAssertion` (the HTTP mapping is sprint 2's `case err != nil`) |
| judge `stored-key` | an unparsable stored key → 422 | POST → judge | unit only (a hand-edited row) |
| judge `law` (DM7 steps) | mint assertion, other verb, revoked, BE, origin… → 422 | POST → judge → `DecideByAssertion` | E2E `TestDecidePOSTPasskey` (mint assertion); unit `TestDecideByAssertionSteps`, `TestJudgeDecisionAssertion` |
| `PasskeyUse()` → `ApproverUse.Passkey` | every verified passkey decision | POST → `keyArm` | E2E `TestDecidePOSTPasskey` |
| K2 passkey rule (`key_decisions.go`) | webauthn without use / ssh with use → `ErrBadKeyDecision` | POST → `keyArm` → `RecordKeyDecision` | **unit only**: `keyArm` always pairs them (the use comes from `PasskeyUse()`, ok exactly for webauthn) — defence in depth, `TestRecordKeyDecisionPasskeyUseMatchesTheFormat` |
| K4 off the exact match | counter ≤ stored → refused; above → moved | POST → `RecordKeyDecision` | E2E `TestDecidePOSTPasskey` (3 over 3), `…RetryAfterAFailedWorldWrite` (another at 7), RD14; unit KD2 |
| K4 skipped on the exact match | the identical unlanded artifact → no counter move | POST → `RecordKeyDecision` | E2E RD13 (b), RD14, KD8; unit KD7 |
| `isCounterRefusal` → 422 | `ErrSignCountRegression` | POST → `keyArm` | E2E `TestDecidePOSTPasskey`; `ErrPasskeyChanged` **unit only** (`TestIsCounterRefusalOperands`): the judge found the passkey by this credential ID in the same request, so a change between judge and writer is a race like `D-RACE` |
| `isApproverRace` → 422 (sprint 2) | revoked/re-keyed/unknown/clock between judge and writer | POST → `keyArm` | **unit only** (`TestIsApproverRaceOperands`, JM01–04): the read-only judge refuses such an approver first (`D-RACE`) |
| CLI `--expiry` rule (`approve.go`) | `--expiry` with a mint id / `--list` → exitUsage | CLI | E2E `TestApproveDecideWriteMessageAndSignature` |
| CLI `--await` with a ref | → exitUsage | CLI | E2E `TestApproveDecideFlagRules` |
| CLI artifact ↔ `--expiry` | either without the other → exitUsage, no GET | CLI | E2E `TestApproveDecideWriteMessageAndSignature`, `TestApproveDecideFlagRules` |
| `decideByArtifact` CD6 | new file → 91 bytes + lines; existing → refused | CLI → GET → subject checks → CD6 | E2E `…WriteMessageAndSignature`, real tool `…SignedBySSHKeygen` |
| CD7 | file → `sshsig` + `--expiry`; missing → refused | CLI → … → CD7 → POST | E2E same tests |
| CD8 `decisionPasskeyLinks` | no active passkey → refused; else links | CLI → … → CD8 | E2E `…WebAuthnLinkAndAssertion`; live Chrome `…InChrome` |
| CD8 `ActivePasskeys` error | below schema 6 → refused | CLI → CD8 | **not reachable** after the health check: the daemon serving this credentials store has upgraded it (survivor `C-LINK-ERR`, §4) |
| CD9 | not an assertion → refused before POST; missing → refused; else `webauthn` | CLI → … → CD9 → POST | E2E `…WebAuthnLinkAndAssertion`, `…WriteMessageAndSignature`; live Chrome |
| JM46 default address | no `--daemon` → the service record's address | CLI | E2E `TestApproveDecideDefaultDaemonIsTheService` |

**No POST arm ships without an operator path** (the directive's rule): `webauthn` is reached by CD8 → the A2 page →
CD9 (§1.3 S3-D2, proven with Chrome's bytes); `sshsig` by `--sign-with` (sprint 2) and CD6 → `ssh-keygen` → CD7.

---

## 4. Mutation table

**Procedure** (`tools/mutate.py <run> <tree> <workers> [ids…]`, rows in `tools/mutants.py`): for each row, an APFS
clone of the tree (`cp -cR`), one exact-string replacement whose anchor occurs exactly once in the file (else
`ANCHOR<n>`, not a run), then `go test -count=1` of the **named test** in its package and of the **package
suite** with the inventory and static tests excluded (`SUITES`: the daemon suite skips `TestDecisionRoutesRegistered`,
`TestApprovalRoutesRegistered` and `TestDaemonCredentialWriterSites`; the cmd suite skips `TestEveryVerbHasHelp`;
the approval suite skips `TestApprovalImportsOnlyTheStandardLibrary`; the store suite runs no deadline or
context-root inventory), each with `GOFLAGS=-vet=off` and `WORLD_CHROME` unset. A build failure is `NOCOMPILE`
and is **not a mutant**: the row is rewritten to compile and re-run. Every run appends one line per row to
`plan-scratch/logs/mutant_ledger.tsv` (time, run, tree, id, stage, file, from → to, named regex, named result,
suite result, failing tests, seconds, expected). **The ids below are the ledger's ids** (`cut -f4 ledger | sort -u`
equals this table's first column). Expected: `K` killed by both; `U` killed by its **unit** test only (the branch
is unreachable end to end, §3.1 — the unit test is behavioural, not an inventory); `E` a named equivalent; `L` a
named survivor (unreachable, with the reason).

**76 rows on the planner's final tree (`work`, tag `FINAL` = M6g), the last ledger line of each: 65 KILLED by the named test and the package suite, 5 `U` rows killed by their unit test (unreachable end to end, reason given), 4 `E` rows (named equivalents, each with its reason), 2 `L` rows (named survivors, unreachable, reason given); 0 not as expected; 0 `NOCOMPILE` among the last lines.** 80 ledger lines in all (`plan-scratch/logs/mutant_ledger.tsv`; runs: R1 68, R2 11, R3 1). 2 ledger lines are `NOCOMPILE` (R1: an edit left `ad` or `id` unused) and are **not counted**: each row was rewritten to compile (`+ 1`, `append(id, 0)`) and re-run in R2; the NOCOMPILE lines stay in the ledger. R2 also ran the M6a rows (JM22-FILTER, K0-APPROVER, K0-FORMAT, JM01–JM04, JM46), which mutate sprint-2 production code the new tests pin. `C-EXPIRY-RULE` SURVIVED in R1 and R2 and exposed two weak test rows in `TestApproveDecideWriteMessageAndSignature` (first: any `--expiry` in the printed help matched; then: the row's own `--daemon` tripped CD3 by itself); the rows now check the refusal text and drop `--daemon`, and R3 KILLED it on the final tree. R1's other cmd rows ran on the tree before those two test edits (they only add kills). Three test arms were added after the table was drafted and before R1 ran, so that a branch had a killing test (`W6-no-message-signed-zero` for `A-W6-HERR`'s decision half, `JD-law-backup-eligible` for `J-A-BE`, the missing-file rows for `C-SIG-READ-ERR`/`C-ASR-READ-ERR`); the patches in Appendix P carry them. **No kill rests on an inventory or static test**: every named test above is behavioural, and the suites exclude the inventories.

| Id | Stage | file: from → to | Named test | Named | Suite | Failing tests (suite, else named) | Expected |
|---|---|---|---|---|---|---|---|
| A-W6-HERR | M6b | passkey_assert.go: if hErr != nil \|\| cd.Challenge != b64Challenge(h[:]) { // W6 → if cd.Challenge != b64Challenge(h[:]) { // W6 | `TestDecideByAssertionSteps` | KILLED | KILLED | TestApproveByAssertionOrder,TestApproveByAssertionOrder/-7/W6-negative-expiry-has-no-hash,TestApproveByAssertionOrder/-8/W6-negative-expiry-has-no-hash,TestDeci | K |
| A-W15-MINT | M6b | passkey_assert.go: if (ad.SignCount != 0 \|\| p.SignCount != 0) && ad.SignCount <= p.SignCount { // W15 → if false { // W15 | `TestChromeAssertionAfterAssertion2IsRefused` | KILLED | KILLED | TestApproveByAssertionOrder,TestApproveByAssertionOrder/-7/W15-counter-not-advanced,TestApproveByAssertionOrder/-8/W15-counter-not-advanced,TestChromeAssertionA | K |
| A-MINT-STEPS-ERR | M6b | passkey_assert.go: ad, err := assertionSteps(in.Response, in.Approver, h, hErr) // W1-W14⏎ if err != nil { → ad, err := assertionSteps(in.Response, in.Approver, h, hErr) // W1-W14⏎ if false && err != nil { | `TestApproveByAssertionOrder` | KILLED | KILLED | TestApproveByAssertionOrder,TestApproveByAssertionOrder/-7/W1-envelope,TestApproveByAssertionOrder/-7/W10-rp-id,TestApproveByAssertionOrder/-7/W11-no-UP,TestApp | K |
| A-W3-NAME | M6b | passkey_assert.go: !bytes.Equal(resp.UserHandle, []byte(a.Name)) { // W3 → !bytes.Equal(resp.UserHandle, []byte(a.Name+"x")) { // W3 | `TestApproveByAssertion` | KILLED | KILLED | TestApproveByAssertionOrder,TestApproveByAssertionOrder/-7/user-handle-is-the-name,TestApproveByAssertionOrder/-8/user-handle-is-the-name | K |
| D-ASSERT-APPROVER | M6b | decision_signed.go: if err := checkApprover(AdmitInput{Now: in.Now}, in.Approver, FormatWebAuthn); err != nil { // P2-P5, P3w → if err := error(nil); err != nil { // P2-P5, P3w | `TestDecideByAssertionSteps` | KILLED | KILLED | TestDecideByAssertionSteps | K |
| D-ASSERT-FORMAT | M6b | decision_signed.go: checkApprover(AdmitInput{Now: in.Now}, in.Approver, FormatWebAuthn) → checkApprover(AdmitInput{Now: in.Now}, in.Approver, in.Approver.Format) | `TestDecideByAssertionSteps` | KILLED | KILLED | TestDecideByAssertionSteps | K |
| D-ASSERT-NOW | M6b | decision_signed.go: checkApprover(AdmitInput{Now: in.Now}, in.Approver, FormatWebAuthn) → checkApprover(AdmitInput{Now: in.Now + 1000}, in.Approver, FormatWebAuthn) | `TestDecideByAssertionSteps` | KILLED | KILLED | TestDecideByAssertionSteps | K |
| D-ASSERT-VERB | M6b | decision_signed.go: h, hErr := DecisionMessageHash(in.Verb, in.Packet, in.Request, in.Expiry) → h, hErr := DecisionMessageHash(VerbApprove, in.Packet, in.Request, in.Expiry) | `TestDecideByAssertionSteps` | KILLED | KILLED | TestDecideByAssertionSteps | K |
| D-ASSERT-PQ | M6b | decision_signed.go: h, hErr := DecisionMessageHash(in.Verb, in.Packet, in.Request, in.Expiry) → h, hErr := DecisionMessageHash(in.Verb, in.Request, in.Packet, in.Expiry) | `TestDecideByAssertionDefersTheCounter` | KILLED | KILLED | TestDecideByAssertionDefersTheCounter,TestDecideByAssertionSteps | K |
| D-ASSERT-EXPIRY | M6b | decision_signed.go: h, hErr := DecisionMessageHash(in.Verb, in.Packet, in.Request, in.Expiry) → h, hErr := DecisionMessageHash(in.Verb, in.Packet, in.Request, in.Expiry+1) | `TestDecideByAssertionDefersTheCounter` | KILLED | KILLED | TestDecideByAssertionDefersTheCounter,TestDecideByAssertionSteps | K |
| D-ASSERT-ERR | M6b | decision_signed.go: ad, err := assertionSteps(in.Response, in.Approver, h, hErr) // W1-W14⏎ if err != nil { → ad, err := assertionSteps(in.Response, in.Approver, h, hErr) // W1-W14⏎ if false && err != nil { | `TestDecideByAssertionSteps` | KILLED | KILLED | TestDecideByAssertionSteps | K |
| D-ASSERT-COUNT | M6b | decision_signed.go: credentialID: bytes.Clone(p.CredentialID), signCount: ad.SignCount} → credentialID: bytes.Clone(p.CredentialID), signCount: ad.SignCount + 1} | `TestDecideByAssertionDefersTheCounter` | KILLED | KILLED | TestDecideByAssertionDefersTheCounter | K |
| D-ASSERT-CRED | M6b | decision_signed.go: credentialID: bytes.Clone(p.CredentialID), signCount: ad.SignCount} → credentialID: nil, signCount: ad.SignCount} | `TestDecideByAssertionDefersTheCounter` | KILLED | KILLED | TestDecideByAssertionDefersTheCounter | K |
| D-ASSERT-KEY | M6b | decision_signed.go: format: FormatWebAuthn, key: p.Key.Canonical(), → format: FormatWebAuthn, key: nil, | `TestDecideByAssertionDefersTheCounter` | KILLED | KILLED | TestDecideByAssertionDefersTheCounter | K |
| D-ASSERT-VERBFIELD | M6b | decision_signed.go: return &SignedDecision{ok: true, verb: in.Verb, packet: in.Packet, request: in.Request, expiry: in.Expiry,⏎  name: in.Approver.Name, format: FormatWebAuthn → return &SignedDecision{ok: true, verb: VerbApprove, packet: in.Packet, request: in.Request, expiry: in.Expiry,⏎  name: in.Approver.Name, format: FormatWebAuthn | `TestDecidePOSTPasskey` | KILLED | KILLED | TestDecidePOSTCounterZeroAssertionLandsOnce,TestDecidePOSTPasskey,TestDecidePOSTPasskeyRetryAfterAFailedWorldWrite | K |
| D-PKUSE-FORMAT | M6b | decision_signed.go: if s == nil \|\| !s.ok \|\| s.format != FormatWebAuthn { → if s == nil \|\| !s.ok { | `TestDecideByAssertionDefersTheCounter` | KILLED | KILLED | TestDecideByAssertionDefersTheCounter | K |
| D-PKUSE-OK | M6b | decision_signed.go: if s == nil \|\| !s.ok \|\| s.format != FormatWebAuthn { → if s == nil \|\| s.format != FormatWebAuthn { | `TestDecideByAssertionDefersTheCounter` | SURVIVED | SURVIVED | — | E: equivalent: an unverified or zero `SignedDecision` has format "", so the format operand refuses it already; only `DecideBySignature`/`DecideByAssertion` set `ok`, and both set the format |
| D-PKUSE-COUNT | M6b | decision_signed.go: return bytes.Clone(s.credentialID), s.signCount, true → return bytes.Clone(s.credentialID), s.signCount + 1, true | `TestDecideByAssertionDefersTheCounter` | KILLED | KILLED | TestDecideByAssertionDefersTheCounter | K |
| J-A-PARSE | M6c | judge_decision.go: resp, err := approval.ParseAssertionResponse(response)⏎ if err != nil {⏎  return nil, &JudgeError{Step: JudgeStepParse, Err: err} → resp, err := approval.ParseAssertionResponse(response)⏎ if err != nil {⏎  return nil, &JudgeError{Step: JudgeStepLaw, Err: err} | `TestJudgeDecisionAssertion` | KILLED | KILLED | TestJudgeDecisionAssertion | K |
| J-A-STORE | M6c | judge_decision.go: a, p, ok, err := ro.PasskeyByCredentialID(ctx, resp.RawID)⏎ if err != nil {⏎  return nil, &JudgeError{Step: JudgeStepStore, Err: err} → a, p, ok, err := ro.PasskeyByCredentialID(ctx, resp.RawID)⏎ if err != nil {⏎  return nil, &JudgeError{Step: JudgeStepLaw, Err: err} | `TestJudgeDecisionAssertion` | KILLED | KILLED | TestJudgeDecisionAssertion | K |
| J-A-UNKNOWN | M6c | judge_decision.go:   return nil, &JudgeError{Step: JudgeStepUnknownKey, RawID: resp.RawID, Err: ErrNotEnrolled}⏎ }⏎ key, err := approval.ParseStoredCOSEKey(a.PublicKey)⏎ if err != nil {⏎  return nil, &JudgeError{Step: JudgeStepStoredKey, Approver: a.Name, Err: err}⏎ }⏎ signed, err := approval.DecideByAssertion →   return nil, &JudgeError{Step: JudgeStepLaw, RawID: resp.RawID, Err: ErrNotEnrolled}⏎ }⏎ key, err := approval.ParseStoredCOSEKey(a.PublicKey)⏎ if err != nil {⏎  return nil, &JudgeError{Step: JudgeStepStoredKey, Approver: a.Name, Err: err}⏎ }⏎ signed, err := approval.DecideByAssertion | `TestJudgeDecisionAssertion` | KILLED | KILLED | TestJudgeDecisionAssertion | K |
| J-A-STOREDKEY | M6c | judge_decision.go:   return nil, &JudgeError{Step: JudgeStepStoredKey, Approver: a.Name, Err: err}⏎ }⏎ signed, err := approval.DecideByAssertion →   return nil, &JudgeError{Step: JudgeStepLaw, Approver: a.Name, Err: err}⏎ }⏎ signed, err := approval.DecideByAssertion | `TestJudgeDecisionAssertion` | KILLED | KILLED | TestJudgeDecisionAssertion | K |
| J-A-LAW | M6c | judge_decision.go: Response: response,⏎  Approver: approval.Approver{Name: a.Name, Format: a.Format, EnrolledAt: a.EnrolledAt, LastUsedAt: a.LastUsedAt, RevokedAt: a.RevokedAt,⏎   Passkey: &approval.Passkey{RPID: p.RPID, Origin: p.Origin, CredentialID: p.CredentialID, Key: key,⏎    BackupEligible: p.BackupEligible, SignCount: uint32(p.SignCount)}}})⏎ if err != nil {⏎  return nil, &JudgeError{Step: JudgeStepLaw, Err: err} → Response: response,⏎  Approver: approval.Approver{Name: a.Name, Format: a.Format, EnrolledAt: a.EnrolledAt, LastUsedAt: a.LastUsedAt, RevokedAt: a.RevokedAt,⏎   Passkey: &approval.Passkey{RPID: p.RPID, Origin: p.Origin, CredentialID: p.CredentialID, Key: key,⏎    BackupEligible: p.BackupEligible, SignCount: uint32(p.SignCount)}}})⏎ if err != nil {⏎  return nil, &JudgeError{Step: JudgeStepParse, Err: err} | `TestJudgeDecisionAssertion` | KILLED | KILLED | TestJudgeDecisionAssertion | K |
| J-A-REVOKED | M6c | judge_decision.go: Expiry: in.Expiry, Now: in.Now, Response: response,⏎  Approver: approval.Approver{Name: a.Name, Format: a.Format, EnrolledAt: a.EnrolledAt, LastUsedAt: a.LastUsedAt, RevokedAt: a.RevokedAt, → Expiry: in.Expiry, Now: in.Now, Response: response,⏎  Approver: approval.Approver{Name: a.Name, Format: a.Format, EnrolledAt: a.EnrolledAt, LastUsedAt: a.LastUsedAt, | `TestJudgeDecisionAssertion` | KILLED | KILLED | TestJudgeDecisionAssertion | K |
| J-A-BE | M6c | judge_decision.go:     BackupEligible: p.BackupEligible, SignCount: uint32(p.SignCount)}}})⏎ if err != nil {⏎  return nil, &JudgeError{Step: JudgeStepLaw, Err: err}⏎ }⏎ return signed, nil⏎}⏎ →     SignCount: uint32(p.SignCount)}}})⏎ if err != nil {⏎  return nil, &JudgeError{Step: JudgeStepLaw, Err: err}⏎ }⏎ return signed, nil⏎}⏎ | `TestJudgeDecisionAssertion` | KILLED | KILLED | TestJudgeDecisionAssertion | K |
| J-A-ORIGIN | M6c | judge_decision.go: Passkey: &approval.Passkey{RPID: p.RPID, Origin: p.Origin, CredentialID: p.CredentialID, Key: key,⏎    BackupEligible: p.BackupEligible, SignCount: uint32(p.SignCount)}}})⏎ if err != nil {⏎  return nil, &JudgeError{Step: JudgeStepLaw, Err: err}⏎ }⏎ return signed, nil⏎}⏎ → Passkey: &approval.Passkey{RPID: p.RPID, Origin: "", CredentialID: p.CredentialID, Key: key,⏎    BackupEligible: p.BackupEligible, SignCount: uint32(p.SignCount)}}})⏎ if err != nil {⏎  return nil, &JudgeError{Step: JudgeStepLaw, Err: err}⏎ }⏎ return signed, nil⏎}⏎ | `TestJudgeDecisionAssertion` | KILLED | KILLED | TestJudgeDecisionAssertion | K |
| J-A-SIGNCOUNT | M6c | judge_decision.go:     BackupEligible: p.BackupEligible, SignCount: uint32(p.SignCount)}}})⏎ if err != nil {⏎  return nil, &JudgeError{Step: JudgeStepLaw, Err: err}⏎ }⏎ return signed, nil⏎}⏎ →     BackupEligible: p.BackupEligible, SignCount: 0}}})⏎ if err != nil {⏎  return nil, &JudgeError{Step: JudgeStepLaw, Err: err}⏎ }⏎ return signed, nil⏎}⏎ | `TestJudgeDecisionAssertion` | SURVIVED | SURVIVED | — | E: equivalent by design (the deferral's evidence): `DecideByAssertion` does not run W15, so the law never reads the stored counter; the writer reads it itself in K4 (`usePasskey`) |
| J-A-LASTUSED | M6c | judge_decision.go: Expiry: in.Expiry, Now: in.Now, Response: response,⏎  Approver: approval.Approver{Name: a.Name, Format: a.Format, EnrolledAt: a.EnrolledAt, LastUsedAt: a.LastUsedAt, → Expiry: in.Expiry, Now: in.Now, Response: response,⏎  Approver: approval.Approver{Name: a.Name, Format: a.Format, EnrolledAt: a.EnrolledAt, | `TestJudgeDecisionAssertion` | SURVIVED | SURVIVED | — | E: equivalent: no law reads `Approver.LastUsedAt` (§1.2 C10); the passkey judge's twin of sprint 2's JM27 |
| K-PASSKEY-GRAMMAR | M6c | key_decisions.go: case (use.Format == "webauthn") != (use.Passkey != nil): → case false: | `TestRecordKeyDecisionPasskeyUseMatchesTheFormat` | KILLED | KILLED | TestRecordKeyDecisionPasskeyUseMatchesTheFormat | K |
| K4-DELETE | M6c | key_decisions.go: if !exact && use.Passkey != nil { // K4 → if false { // K4 | `TestRecordKeyDecisionCounter` | KILLED | KILLED | TestRecordKeyDecisionCounter,TestRecordKeyDecisionIdempotentPasskeyRetry | K |
| K4-EXACT | M6c | key_decisions.go: if !exact && use.Passkey != nil { // K4 → if use.Passkey != nil { // K4 | `TestRecordKeyDecisionIdempotentPasskeyRetry` | KILLED | KILLED | TestRecordKeyDecisionIdempotentPasskeyRetry | K |
| K4-ERR | M6c | key_decisions.go: if err := usePasskey(ctx, tx, use.Name, *use.Passkey); err != nil {⏎   return err → if err := usePasskey(ctx, tx, use.Name, *use.Passkey); err != nil {⏎   _ = err | `TestRecordKeyDecisionCounter` | KILLED | KILLED | TestRecordKeyDecisionCounter,TestRecordKeyDecisionIdempotentPasskeyRetry | K |
| KA-COUNTER-422 | M6d | surface_decide_key.go:   case isCounterRefusal(err): →   case false: | `TestDecidePOSTPasskey` | KILLED | KILLED | TestDecidePOSTPasskey,TestDecidePOSTPasskeyRetryAfterAFailedWorldWrite,TestDecidePOSTRetryAfterAFailedWorldWriteAcrossRestart | K |
| KA-COUNTER-REGRESS | M6d | surface_decide_key.go: return errors.Is(err, store.ErrSignCountRegression) \|\| errors.Is(err, store.ErrPasskeyChanged) → return errors.Is(err, store.ErrPasskeyChanged) | `TestDecidePOSTPasskey` | KILLED | KILLED | TestDecidePOSTPasskey,TestDecidePOSTPasskeyRetryAfterAFailedWorldWrite,TestDecidePOSTRetryAfterAFailedWorldWriteAcrossRestart,TestIsCounterRefusalOperands | K |
| KA-COUNTER-CHANGED | M6d | surface_decide_key.go: return errors.Is(err, store.ErrSignCountRegression) \|\| errors.Is(err, store.ErrPasskeyChanged) → return errors.Is(err, store.ErrSignCountRegression) | `TestIsCounterRefusalOperands` | KILLED | KILLED | TestIsCounterRefusalOperands | U: unit only: `ErrPasskeyChanged` needs the credential to change between the judge's lookup by this credential ID and the writer, a race like `D-RACE` (§3.1) |
| KA-PASSKEYUSE | M6d | surface_decide_key.go: if id, count, ok := signed.PasskeyUse(); ok { → if id, count, ok := signed.PasskeyUse(); false && ok { | `TestDecidePOSTPasskey` | KILLED | KILLED | TestDecidePOSTCounterZeroAssertionLandsOnce,TestDecidePOSTPasskey,TestDecidePOSTPasskeyRetryAfterAFailedWorldWrite,TestDecidePOSTRetryAfterAFailedWorldWriteAcro | K |
| KA-PASSKEYUSE-COUNT | M6d | surface_decide_key.go: use.Passkey = &store.PasskeyUse{CredentialID: id, SignCount: int64(count)} → use.Passkey = &store.PasskeyUse{CredentialID: id, SignCount: int64(count) + 1} | `TestDecidePOSTPasskeyRetryAfterAFailedWorldWrite` | KILLED | KILLED | TestDecidePOSTCounterZeroAssertionLandsOnce,TestDecidePOSTPasskey,TestDecidePOSTPasskeyRetryAfterAFailedWorldWrite,TestDecidePOSTRetryAfterAFailedWorldWriteAcro | K |
| KA-PASSKEYUSE-ID | M6d | surface_decide_key.go: use.Passkey = &store.PasskeyUse{CredentialID: id, SignCount: int64(count)} → use.Passkey = &store.PasskeyUse{CredentialID: append(id, 0), SignCount: int64(count)} | `TestDecidePOSTPasskey` | KILLED | KILLED | TestDecidePOSTCounterZeroAssertionLandsOnce,TestDecidePOSTPasskey,TestDecidePOSTPasskeyRetryAfterAFailedWorldWrite,TestDecidePOSTRetryAfterAFailedWorldWriteAcro | K |
| KA-JUDGE-SELECT | M6d | surface_decide_key.go: if body.WebAuthn != nil { // the grammar admits exactly one key format → if false { // the grammar admits exactly one key format | `TestDecidePOSTPasskey` | KILLED | KILLED | TestDecidePOSTPasskey | K |
| KA-WHAT | M6d | surface_decide_key.go: []byte(*body.WebAuthn), "passkey assertion" → []byte(*body.WebAuthn), "SSH signature" | `TestDecidePOSTPasskey` | KILLED | KILLED | TestDecidePOSTPasskey | K |
| KA-RACE-TEXT | M6d | surface_decide_key.go: "the approver changed since the "+what+" was judged: " → "the approver changed since the signature was judged: " | `TestIsApproverRaceOperands` | SURVIVED | SURVIVED | — | L: survivor: the race branch's text is unreachable end to end (`D-RACE`) and the unit test pins the predicate, not the handler's text |
| SD-ARM | M6d | surface_decide.go: if body.SSHSIG != nil \|\| body.WebAuthn != nil { → if body.SSHSIG != nil { | `TestDecidePOSTPasskey` | KILLED | KILLED | TestDecidePOSTCounterZeroAssertionLandsOnce,TestDecidePOSTPasskey,TestDecidePOSTPasskeyRetryAfterAFailedWorldWrite,TestDecidePOSTRetryAfterAFailedWorldWriteAcro | K |
| DM-FORMATS-WEBAUTHN | M6d | decide_model.go: for _, on := range []bool{body.Typed, body.SSHSIG != nil, body.WebAuthn != nil} { → for _, on := range []bool{body.Typed, body.SSHSIG != nil} { | `TestDecidePOSTPasskey` | KILLED | KILLED | TestDecidePOSTCounterZeroAssertionLandsOnce,TestDecidePOSTPasskey,TestDecidePOSTPasskeyRetryAfterAFailedWorldWrite,TestDecidePOSTRetryAfterAFailedWorldWriteAcro | K |
| DM-FORMATS-TWO | M6d | decide_model.go: case decideFormats(body) > 1: → case decideFormats(body) > 2: | `TestDecidePOSTPasskey` | KILLED | KILLED | TestDecidePOSTGuards,TestDecidePOSTPasskey | K |
| DM-FORMATS-ZERO | M6d | decide_model.go: case decideFormats(body) == 0: → case decideFormats(body) < 0: | `TestDecidePOSTGuards` | KILLED | KILLED | TestDecidePOSTGuards | K |
| C-EXPIRY-RULE | M6e | approve.go: if (*deny \|\| *daemonAddr != "" \|\| *expiry != 0) && (*list → if (*deny \|\| *daemonAddr != "") && (*list | `TestApproveDecideWriteMessageAndSignature` | KILLED | KILLED | TestApproveDecideWriteMessageAndSignature | K |
| C-AWAIT | M6e | approve.go: if *await { // CD3 → if false { // CD3 | `TestApproveDecideFlagRules` | KILLED | KILLED | TestApproveDecideFlagRules | K |
| C-ARTIFACT-EXPIRY | M6e | approve.go: if (*sigPath != "" \|\| *assertion != "") != (*expiry != 0) { → if false { | `TestApproveDecideWriteMessageAndSignature` | KILLED | KILLED | TestApproveDecideWriteMessageAndSignature | K |
| C-ARTIFACT-EXPIRY-ASR | M6e | approve.go: if (*sigPath != "" \|\| *assertion != "") != (*expiry != 0) { → if (*sigPath != "") != (*expiry != 0) { | `TestApproveDecideFlagRules` | KILLED | KILLED | TestApproveDecideFlagRules,TestApproveDecideWebAuthnLinkAndAssertion,TestApproveDecideWriteMessageAndSignature | K |
| C-OPTS-LINK | M6e | approve.go: passkeyLink: *passkeyLink, expiry: *expiry} → passkeyLink: false, expiry: *expiry} | `TestApproveDecideWebAuthnLinkAndAssertion` | KILLED | KILLED | TestApproveDecideWebAuthnLinkAndAssertion | K |
| C-ART-DISPATCH | M6e | approve_decide.go: if o.sigPath != "" \|\| o.msgPath != "" \|\| o.assertion != "" \|\| o.passkeyLink { → if o.sigPath != "" \|\| o.msgPath != "" \|\| o.assertion != "" { | `TestApproveDecideWebAuthnLinkAndAssertion` | KILLED | KILLED | TestApproveDecideWebAuthnLinkAndAssertion | K |
| C-ART-EXPIRY | M6e | approve_decide.go:  if o.expiry != 0 {⏎  body.Expiry = o.expiry⏎ } →  if false {⏎  body.Expiry = o.expiry⏎ } | `TestApproveDecideWriteMessageAndSignature` | KILLED | KILLED | TestApproveDecideWriteMessageAndSignature | K |
| C-WM-EXCL | M6e | approve_decide.go: os.O_WRONLY\|os.O_CREATE\|os.O_EXCL, 0o644) → os.O_WRONLY\|os.O_CREATE\|os.O_TRUNC, 0o644) | `TestApproveDecideWriteMessageAndSignature` | KILLED | KILLED | TestApproveDecideWriteMessageAndSignature | K |
| C-WM-ERR | M6e | approve_decide.go:   if err != nil {⏎   return refuse("--write-message: %v; nothing was written", err) →   if false && err != nil {⏎   return refuse("--write-message: %v; nothing was written", err) | `TestApproveDecideWriteMessageAndSignature` | KILLED | KILLED | TestApproveDecideWriteMessageAndSignature | K |
| C-WM-MSG | M6e | approve_decide.go:    _, err = f.Write(msg) →    _, err = f.Write(msg[1:]) | `TestApproveDecideWriteMessageAndSignature` | KILLED | KILLED | TestApproveDecideWriteMessageAndSignature,TestApproveDecideWriteMessageSignedBySSHKeygen | K |
| C-WM-DENY | M6e | approve_decide.go: map[string]string{"approve": "", "deny": " --deny"}[o.verb] → map[string]string{"approve": "", "deny": ""}[o.verb] | `TestApproveDecideWriteMessageAndSignature` | KILLED | KILLED | TestApproveDecideWebAuthnLinkAndAssertion,TestApproveDecideWriteMessageAndSignature | K |
| C-WM-NAMESPACE | M6e | approve_decide.go: approverTime(body.Expiry), approval.DecisionNamespace, shellQuote(o.msgPath) → approverTime(body.Expiry), approval.Namespace, shellQuote(o.msgPath) | `TestApproveDecideWriteMessageSignedBySSHKeygen` | KILLED | KILLED | TestApproveDecideWriteMessageAndSignature,TestApproveDecideWriteMessageSignedBySSHKeygen | K |
| C-LINK-HASH | M6e | approve_decide.go:  h := sha256.Sum256(msg)⏎ →  h := sha256.Sum256(msg[1:])⏎ | `TestApproveDecideWebAuthnLinkAndAssertion` | KILLED | KILLED | TestApproveDecideWebAuthnLinkAndAssertion | K |
| C-LINK-REQUEST | M6e | approve_decide.go: "decide-"+o.verb+"-"+strings.TrimPrefix(o.ref, "sha256:")[:16] → "decide-approve-"+strings.TrimPrefix(o.ref, "sha256:")[:16] | `TestApproveDecideWebAuthnLinkAndAssertion` | KILLED | KILLED | TestApproveDecideWebAuthnLinkAndAssertion | K |
| C-LINK-NONE | M6e | approve_decide.go:  if len(active) == 0 {⏎  return refuse("no passkey approver is enrolled →  if false {⏎  return refuse("no passkey approver is enrolled | `TestApproveDecideWebAuthnLinkAndAssertion` | KILLED | KILLED | TestApproveDecideWebAuthnLinkAndAssertion | K |
| C-LINK-ERR | M6e | approve_decide.go:  active, err := ro.ActivePasskeys(ctx)⏎ _ = ro.Close()⏎ if err != nil {⏎  return refuse("%v; nothing was decided", err) →  active, err := ro.ActivePasskeys(ctx)⏎ _ = ro.Close()⏎ if false && err != nil {⏎  return refuse("%v; nothing was decided", err) | `TestApproveDecideWebAuthnLinkAndAssertion` | SURVIVED | SURVIVED | — | L: survivor, unreachable: `ActivePasskeys` fails only below schema 6, and the CLI has proved (health) that a daemon serves this credentials store, which `serve` upgraded |
| C-LINK-ALLOW | M6e | approve_decide.go:   l.Allow = append(l.Allow, base64.RawURLEncoding.EncodeToString(e.Passkey.CredentialID))⏎ }⏎ for _, origin := range order { →   _ = l⏎ }⏎ for _, origin := range order { | `TestApproveDecideWebAuthnLinkAndAssertion` | KILLED | KILLED | TestApproveDecideWebAuthnLinkAndAssertion | K |
| C-LINK-EXPIRES | M6e | approve_decide.go: Request: request, Expires: expiry, Hash16: hash16, → Request: request, Expires: 0, Hash16: hash16, | `TestApproveDecideWebAuthnLinkAndAssertion` | KILLED | KILLED | TestApproveDecideWebAuthnLinkAndAssertion | K |
| C-SIG-READ-ERR | M6e | approve_decide.go:   if err != nil {⏎   return refuse("--signature: %v; nothing was decided", err) →   if false && err != nil {⏎   return refuse("--signature: %v; nothing was decided", err) | `TestApproveDecideWriteMessageAndSignature` | KILLED | KILLED | TestApproveDecideWriteMessageAndSignature | K |
| C-SIG-BODY | M6e | approve_decide.go:   sig := string(armored)⏎  body.SSHSIG = &sig⏎ default: // CD9 →   _ = armored⏎ default: // CD9 | `TestApproveDecideWriteMessageAndSignature` | KILLED | KILLED | TestApproveDecideWriteMessageAndSignature,TestApproveDecideWriteMessageSignedBySSHKeygen | K |
| C-ASR-READ-ERR | M6e | approve_decide.go:   if err != nil {⏎   return refuse("--webauthn-assertion: %v; nothing was decided", err) →   if false && err != nil {⏎   return refuse("--webauthn-assertion: %v; nothing was decided", err) | `TestApproveDecideWriteMessageAndSignature` | KILLED | KILLED | TestApproveDecideWriteMessageAndSignature | K |
| C-ASR-PARSE | M6e | approve_decide.go: if _, err := approval.ParseAssertionResponse(file); err != nil { → if _, err := approval.ParseAssertionResponse(file); false && err != nil { | `TestApproveDecideWebAuthnLinkAndAssertion` | KILLED | KILLED | TestApproveDecideWebAuthnLinkAndAssertion | K |
| C-ASR-BODY | M6e | approve_decide.go:   asr := json.RawMessage(file)⏎  body.WebAuthn = &asr →   _ = file | `TestApproveDecideWebAuthnLinkAndAssertion` | KILLED | KILLED | TestApproveDecideWebAuthnLinkAndAssertion | K |
| JM22-FILTER | M6a | key_decisions.go: `UPDATE key_decisions SET decision_ref = ? WHERE request_ref = ? AND decision_ref IS NULL`, decisionRef, request) → `UPDATE key_decisions SET decision_ref = ? WHERE decision_ref IS NULL`, decisionRef) | `TestLandKeyDecisionLandsOnlyItsRequest` | KILLED | KILLED | TestLandKeyDecisionLandsOnlyItsRequest | K |
| K0-APPROVER | M6a | key_decisions.go: old.Format == row.Format && old.Approver == row.Approver && old.Expiry == row.Expiry → old.Format == row.Format && old.Expiry == row.Expiry | `TestRecordKeyDecisionExactMatchNeedsTheSameApprover` | KILLED | KILLED | TestRecordKeyDecisionExactMatchNeedsTheSameApprover | K |
| K0-FORMAT | M6a | key_decisions.go: old.Format == row.Format && old.Approver == row.Approver && old.Expiry == row.Expiry → old.Approver == row.Approver && old.Expiry == row.Expiry | `TestRecordKeyDecisionExactMatchNeedsTheSameApprover` | SURVIVED | SURVIVED | — | E: equivalent at the store API too: K2 holds `row.Approver == use.Format + ":" + use.Name` and `row.Format == use.Format`, so equal approvers imply equal formats (§1.2 C7) |
| JM01 | M6a | surface_decide_key.go: return errors.Is(err, store.ErrApproverRevoked) \|\| errors.Is(err, store.ErrApproverKeyChanged) \|\| → return errors.Is(err, store.ErrApproverKeyChanged) \|\| | `TestIsApproverRaceOperands` | KILLED | KILLED | TestIsApproverRaceOperands | U: unit only: `D-RACE` is unreachable through the POST (the read-only judge refuses first, §3.1) |
| JM02 | M6a | surface_decide_key.go: return errors.Is(err, store.ErrApproverRevoked) \|\| errors.Is(err, store.ErrApproverKeyChanged) \|\| → return errors.Is(err, store.ErrApproverRevoked) \|\| | `TestIsApproverRaceOperands` | KILLED | KILLED | TestIsApproverRaceOperands | U: unit only (as JM01) |
| JM03 | M6a | surface_decide_key.go:   errors.Is(err, store.ErrApproverUnknown) \|\| errors.Is(err, store.ErrApproverClockBehind) →   errors.Is(err, store.ErrApproverClockBehind) | `TestIsApproverRaceOperands` | KILLED | KILLED | TestIsApproverRaceOperands | U: unit only (as JM01) |
| JM04 | M6a | surface_decide_key.go:   errors.Is(err, store.ErrApproverUnknown) \|\| errors.Is(err, store.ErrApproverClockBehind) →   errors.Is(err, store.ErrApproverUnknown) | `TestIsApproverRaceOperands` | KILLED | KILLED | TestIsApproverRaceOperands | U: unit only (as JM01) |
| JM46 | M6a | approve_decide.go:  if o.daemon == "" {⏎  o.daemon = defaultServiceAddr() →  if false {⏎  o.daemon = defaultServiceAddr() | `TestApproveDecideDefaultDaemonIsTheService` | KILLED | KILLED | TestApproveDecideDefaultDaemonIsTheService | K |

**Sprint-2 named survivors (its §4) — closed where this sprint changed the guarded code:** `D-RACE` — its four
operands are now pinned by `TestIsApproverRaceOperands` (JM01–JM04, `U`); the branch itself stays unreachable
through the POST (§3.1). `D-LAND-ERR`, `C-PROOF-ERR`, `C-SIGN-NIL` — the code they guard
(`LandKeyDecision`'s error arm in the handler, `writeOperatorProof`'s error arm, the nil decision signer) is not
changed by this sprint; they stay named survivors with sprint 2's reasons.

---

## 5. Files (final `git status --short` after M6g)

```
 M cmd/ailang-worldd/approve.go
 M cmd/ailang-worldd/approve_decide.go
 M cmd/ailang-worldd/approve_decide_key_test.go
 M cmd/ailang-worldd/approve_decide_rig_test.go
 M docs/QUICKSTART.md
 M host/approval/decision_signed.go
 M host/approval/passkey_assert.go
 M host/authority/judge_decision.go
 M host/daemon/surface_decide.go
 M host/daemon/surface_decide_key.go
 M host/store/key_decisions.go
 M host/surfacepage/decide_model.go
?? cmd/ailang-worldd/approve_decide_artifact_test.go
?? cmd/ailang-worldd/approve_decide_chrome_test.go
?? cmd/ailang-worldd/approve_decide_jgaps_test.go
?? cmd/ailang-worldd/quickstart_decide_test.go
?? host/approval/decision_assert_test.go
?? host/authority/judge_decision_assert_test.go
?? host/daemon/surface_decide_passkey_test.go
?? host/daemon/surface_decide_race_test.go
?? host/store/key_decisions_jgaps_test.go
?? host/store/key_decisions_passkey_test.go
```
plus the two plan files (`design_docs/planned/w-world-approvals-by-key-a4-sprint3-plan.md`,
`design_docs/planned/sprint_w-world-approvals-by-key-a4-s3.json`), which the controller commits. Nothing under
`world/`, `packages/`, `scripts/`, `tools/`, `host/broker/`, `host/coordinator/`, `host/surfacepage/{app,approve,passkey}/`.

---

## 6. Sprint gates (controller, after the last sub-milestone; each measured at the base, §1.1)

Run in **one** `/bin/bash` shell (macOS bash 3.2: no `declare -A`), **not under `nohup`**, in the foreground or
with the Bash tool's own background + a bounded poll of `legs.txt` for `GATES_DONE`; `mkdir -p
$HOME/.ailang/state/world-iter391/gates` first. Every leg records its own rc and **every leg runs** (no `set -e`);
every leg is `tee`'d to `<logdir>/<leg>.log`. Usage on the worktree: `gates.sh
/Users/voightkampff/dev/sunholo-data/.wt-world-iter391-a4-s3 $HOME/.ailang/state/world-iter391/gates final`.
The script (`plan-scratch/tools/gates.sh`, verbatim):

```bash
#!/bin/bash
# A4 sprint 3 pre-commit gates (plan §6). usage: gates.sh <tree> <logdir> <base|final>
# macOS /bin/bash 3.2: no declare -A. Every leg records its own rc; every leg runs (no set -e).
T=$1; L=$2; MODE=$3
mkdir -p "$L"; cd "$T" || exit 2
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
leg verify_ail './scripts/verify_ail.sh'
leg gofmt 'out=$(gofmt -l cmd host); echo "$out"; test -z "$out"'
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./... > '$L/fence.raw' 2>&1; rc=\$?; cat '$L/fence.raw'; if grep -E '^(FAIL|--- FAIL)' '$L/fence.raw' >/dev/null; then exit 1; fi; exit \$rc"
leg pkgs 'go test -count=1 -timeout 30m ./host/approval/... ./host/authority/ ./host/store/ ./host/sshsign/ ./host/surfacepage/ ./host/daemon/ ./cmd/ailang-worldd/'
leg broker_named 'go test -count=1 ./host/broker/ -run "^(TestReservedApprovalScope|TestDecisionRefusal|TestDecideApproval|TestPublishApproval|TestApprovalContextIdentity|TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential)" -v'
leg broker_whole_reported 'go test -count=1 -timeout 30m ./host/broker/'
leg race 'go test -race -count=1 ./host/approval/ ./host/authority/ ./host/store/ ./host/sshsign/ && go test -race -count=1 ./host/daemon/ -run "^(TestDecision|TestDecide|TestApprov|TestSurfaceApprove|TestDaemonCredentialWriterSites|TestEpisodeFinish)" && go test -race -count=1 ./cmd/ailang-worldd/ -run "^(TestApproveDecide|TestApprove|TestApprover)"'
NEW_A="TestDecideByAssertionSteps TestDecideByAssertionDefersTheCounter"
NEW_AU="TestJudgeDecisionAssertion"
NEW_S="TestLandKeyDecisionLandsOnlyItsRequest TestRecordKeyDecisionExactMatchNeedsTheSameApprover TestRecordKeyDecisionCounter TestRecordKeyDecisionIdempotentPasskeyRetry TestRecordKeyDecisionPasskeyUseMatchesTheFormat"
NEW_D="TestIsApproverRaceOperands TestIsCounterRefusalOperands TestDecidePOSTPasskey TestDecidePOSTCounterZeroAssertionLandsOnce TestDecidePOSTPasskeyRetryAfterAFailedWorldWrite TestDecidePOSTRetryAfterAFailedWorldWriteAcrossRestart"
NEW_C="TestApproveDecideDefaultDaemonIsTheService TestApproveDecideWriteMessageAndSignature TestApproveDecideWriteMessageSignedBySSHKeygen TestApproveDecideWebAuthnLinkAndAssertion"
NEW_X="TestApproveDecideWebAuthnLinkInChrome"
if [ "$MODE" = base ]; then NEW_A=""; NEW_AU=""; NEW_S=""; NEW_D=""; NEW_C=""; NEW_X=""; fi
alt() { local x="$*"; x="${x// /|}"; echo "${x:+|$x}"; }
INV_APPROVAL="TestApprovalImportsOnlyTheStandardLibrary TestApprovalTestIsTestOnly TestDecisionMessageLayout TestDecisionMessageRefusesBadInput"
INV_AUTH="TestJudgeSSHSignatureIsTheCLIJudge TestJudgeAssertionIsTheCLIJudge"
INV_BROKER="TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential TestQuickstartSection10ProfilesLoad TestQuickstartProtectedContract TestApprovalContextIdentity TestReservedApprovalScope TestDecideApprovalCallerCensus"
INV_DAEMON="TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestReadCtxCancelledAfterHandler TestMCPMountMethodSource TestDaemonCredentialWriterSites TestQuickstartDecisionPacketSection TestDaemonConstructorCallerInventory TestDecidePOSTGuards TestDecideDiscoveryFindsAHiddenSecondPacket TestDecideDiscoveryRefusesAnOverflow"
INV_BOUNDARY="TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard TestWorkbenchPackageRemainsTransportFree"
INV_STORE="TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext TestCredentialsMethodsRequireADeadline TestCredentialsSchemaMatchesCanonicalDDL TestCredentialsSchema7SurfaceApprovals"
INV_VG="TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestVerifygateTestWritesAreForkLocked"
INV_CMD="TestEveryVerbHasHelp TestQuickstartAttendedStepsBoxMatchesSection9 TestQuickstartSection9FlagsMatchTheCLI TestQuickstartSection10FlagsMatchTheCLI TestQuickstartSection11FlagsMatchTheCLI TestQuickstartServiceSectionMatchesTheCLIHelp TestRealSessionEnvSignsWithSSHKeygen TestSignWithSSHKeygenEndToEnd"
INV_RUNBOOK="TestNoCIStepOrScriptReachesThePublishEntrypoint"
leg inv "set -o pipefail
go test -v -count=1 ./host/approval/... -run '^($(echo $INV_APPROVAL | tr ' ' '|')$(alt $NEW_A))\$' &&
go test -v -count=1 ./host/authority/ -run '^($(echo $INV_AUTH | tr ' ' '|')$(alt $NEW_AU))\$' &&
go test -v -count=1 ./host/broker/ -run '^($(echo $INV_BROKER | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/daemon/ -run '^($(echo $INV_DAEMON | tr ' ' '|')$(alt $NEW_D))\$' &&
go test -v -count=1 ./host/boundary/ -run '^($(echo $INV_BOUNDARY | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/store/ -run '^($(echo $INV_STORE | tr ' ' '|')$(alt $NEW_S))\$' &&
go test -v -count=1 ./host/verifygate/ -run '^($(echo $INV_VG | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/runbook/ -run '^($(echo $INV_RUNBOOK | tr ' ' '|'))\$' &&
go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(echo $INV_CMD | tr ' ' '|')$(alt $NEW_C))\$'"
INV2_SP="TestSurfaceAppPure TestSurfaceScriptPure TestSurfaceFontEmbed TestSurfaceScriptAssetSeams"
INV2_D="TestApprovePageRoutes TestPasskeyPageRoutes TestApprovalRoutesRefuseForeignHost TestGETRoutesRejectOtherMethods TestDaemonCallerInventoryAliasControl"
leg inv2 "set -o pipefail
go test -v -count=1 ./host/surfacepage/ -run '^($(echo $INV2_SP | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/daemon/ -run '^($(echo $INV2_D | tr ' ' '|'))\$'"
for t in $INV2_SP $INV2_D; do grep -qE -- "^--- PASS: $t " "$L/inv2.log" || { echo "INV2 NOT PASSED: $t" | tee -a "$L/legs.txt"; FAILS="$FAILS inv2_missing"; }; done
XB="TestApprovePageInChrome TestApprovePageRefusesASwappedRequest TestApprovePageSSHSIGInChrome TestApprovePageShowsAHostileValueAsText TestPasskeyPageAbortsAPendingCeremony TestPasskeyPageInChrome"
leg chrome "WORLD_CHROME='/Applications/Google Chrome.app/Contents/MacOS/Google Chrome' WORLD_REQUIRE_PASSKEY_BROWSER=1 go test ./cmd/ailang-worldd/ -run '^($(echo $XB | tr ' ' '|')$(alt $NEW_X))\$' -count=1 -v"
for t in $XB $NEW_X; do grep -qE -- "^--- PASS: $t " "$L/chrome.log" || { echo "CHROME NOT PASSED: $t" | tee -a "$L/legs.txt"; FAILS="$FAILS chrome_missing"; }; done
# every named test must PASS (none SKIP, none FAIL)
MISSING=""
for t in $INV_APPROVAL $INV_AUTH $INV_BROKER $INV_DAEMON $INV_BOUNDARY $INV_STORE $INV_VG $INV_CMD $INV_RUNBOOK $NEW_A $NEW_AU $NEW_S $NEW_D $NEW_C; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || MISSING="$MISSING $t"
done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then FAILS="$FAILS inv_skip_or_fail"; fi
[ -z "$MISSING" ] || { echo "NOT PASSED:$MISSING" | tee -a "$L/legs.txt"; FAILS="$FAILS inv_missing"; }
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.log")" | tee -a "$L/legs.txt"
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" | tee -a "$L/legs.txt"
```

Then, **after the controller's commit** (it reads committed blobs):

```
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter391-a4-s3 && VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh 2>&1 | tee $HOME/.ailang/state/world-iter391/gates/verify_go.log
```

Pass rule: every leg rc 0 except `broker_whole_reported` (reported, not gating: the known-flaky list); `inv` names
every test `--- PASS`, none SKIP or FAIL; the `chrome` leg names every XB test and
`TestApproveDecideWebAuthnLinkInChrome` `--- PASS`. A red in a package this sprint does not touch is re-run once
alone before it is believed (§1.1's rule; the directive's known-flaky list). No production context root is added
(§1.2 C16), so `TestProductionContextRoots` is an `inv` row, not a new gate.

---

## 7. Executor directive notes

- **No git writes of any kind** (no add, commit, stash, checkout, reset, branch). **`patch -p1` (§2) is the only
  writer** of the tree; no hand edit, no `gofmt -w`, no fuzz. A failed `--dry-run` or a red acceptance → **stop
  and report** with the log; do not repair.
- **Runs**: Run 1 = M6a, M6b, M6c; Run 2 = M6d, M6e; Run 3 = M6f, M6g. After each sub-milestone's acceptance,
  snapshot the tree to `$HOME/.ailang/state/world-iter391/snap/<M>/` (`rsync -a --exclude .git` of `cmd host docs`)
  with a `MANIFEST.sha256` (`find cmd host docs -type f -print0 | sort -z | xargs -0 shasum -a 256`), **excluding
  the two plan files**. The controller commits each snapshot as one commit and checks the numstat against §3's
  expected value (`git diff --numstat 7004f591 HEAD -- cmd host docs | awk '{s+=$1+$2} END {print s}'`): M6a 124,
  M6b 381, M6c 628, M6d 963, M6e 1,328, M6f 1,411, M6g 1,493. A different number → stop.
- **Every Go Bash call passes `timeout: 600000`, runs in the foreground, no Monitor.** **Never `nohup`, never
  background a test leg.** `export AILANG_BIN=$HOME/.pinned-ailang/ailang` in every shell that runs Go tests. The
  controller's shell is zsh: run each acceptance block under `bash <<'EOF' … EOF`.
- **Exact acceptance commands** are §3's, verbatim; M6f's carries `WORLD_CHROME="/Applications/Google
  Chrome.app/Contents/MacOS/Google Chrome" WORLD_REQUIRE_PASSKEY_BROWSER=1` on the command line itself (copy it
  verbatim; the driver makes its own temp `--user-data-dir`; never a real profile). `ssh-keygen` (`/usr/bin`) must
  be on PATH for M6e's real-tool sibling (it skips, counted, otherwise; `WORLD_REQUIRE_SSH_KEYGEN=1` makes that a
  failure).
- **Stop conditions**: a dry-run failure; any `--- FAIL`, `--- SKIP` (M6f with its env must not skip) or `FAIL`
  line in an acceptance; gofmt output; a numstat mismatch; any file outside §5 changed (`git status --short` after
  each patch must equal §5's list up to that stage).
- **Known load-sensitive tests** (the directive's list) are not this sprint's: a red there is re-run once alone.
- D-WORLD-109: nothing touches `~/.ailang/world`, `~/.ssh`, a real `/dev/tty`, a real browser profile, or a real
  daemon (§1.4 P7: JM46's test refuses to run unless the temp service record resolves to the fake daemon).

---

## 8. Risks

- **R1 — The A2 page's wording for a deny** (§1.4 P1): "Approve request decide-deny-…" can mislead a hurried
  operator. Mitigations here: the verb in the link's `request`, in capitals in the CLI, and the signed hash the
  page shows; the cryptography refuses a cross-verb use (DM5, W6 rows). Sprint 4's card replaces this path for the
  phone; `passkey.js` is not edited (design §2.8).
- **R2 — The refactor of the A2 law.** Mitigated by running the A2 matrix unedited (M6b acceptance), the
  recorded Chrome fixtures (`TestApproveByAssertionAcceptsChrome`, `TestChromeAssertionAfterAssertion2IsRefused`)
  and the mutants `A-W6-HERR`, `A-W15-MINT`, `A-MINT-STEPS-ERR`, `A-W3-NAME`, all killed by A2 tests.
- **R3 — Live Chrome only on the controller's machine** (CI has no Chrome step, A3 V13): M6f is gated as XB1; a
  CI run skips it, counted. Its approvaltest twin (CD8/CD9) runs everywhere.
- **R4 — `json.Marshal` compacts the posted envelope** (S3-D5): the `artifact` the daemon stores is the compacted
  file, deterministic for the same file, so K0's byte equality holds for an identical re-post; a re-formatted
  file with the same content is also byte-equal after compaction (harmless: same signature, same counter).
- **R5 — Load.** The rig ran at load 20–45 during this planner's runs; the daemon passkey tests take < 0.1 s each
  but share the package's other slow tests. Re-run once alone before believing a red (§1.1).
- **R6 — D-WORLD-115 = B** would change the message's P slot for both laws and CD6/CD8 (S3-D8).
- **R7 — The surface is a fourth sprint** (S3-D1): until it lands an operator finds request refs on the packet
  page (QUICKSTART §11) and decides from a terminal; the phone can sign via the CD8 link but the file must come
  back to the terminal.

---

## 9. Acceptance criteria (the judge scores against these)

1. **P1 is measured and recorded** as a V-row (§1.2 C4): Chrome rebuilds DM1's two 91-byte messages and hashes
   byte-equal to Go under `surfaceCSP`.
2. **The sprint-2 judge's findings 1–4 are closed or named**: JM22 (`TestLandKeyDecisionLandsOnlyItsRequest`),
   K0's approver operand (`TestRecordKeyDecisionExactMatchNeedsTheSameApprover`; the format operand a named
   equivalent with its reason), JM01–JM04 (`TestIsApproverRaceOperands`, unit, the branch's reachability stated),
   JM46 (`TestApproveDecideDefaultDaemonIsTheService`), JM27 (named equivalent, C10); each with its ledger row
   killed as expected.
3. **The A2 mint law is behaviourally unchanged**: every A2 test runs unedited and green; `judge.go`,
   `passkeys.go`, `usePasskey` untouched (`git diff` empty for them).
4. **`DecideByAssertion` runs P2–P5, P3w and W1–W14 over the decision message's hash and not W15** (DM7: every
   step refuses; counters 1 and 0 over stored 1 verify at the law).
5. **The counter is judged in the writer's transaction, off the exact match only** (KD2, KD7 passkey half, K2's
   passkey rule): a non-advancing counter writes nothing; the exact unlanded artifact records without moving the
   counter; any other artifact at the same counter is refused.
6. **The POST decides with a passkey** as `webauthn:<name>`, derived (RD10/RD11 passkey halves), with 400 for two
   formats and a forged `decidedBy`, 422 for a mint assertion and a regressed counter, and **nothing decided,
   recorded or moved on any refusal**.
7. **Single use and retry hold with passkeys**: KD8 (counter 0 lands once; after a failed write, once), RD13 (b)
   (counter N: the identical retry lands without moving the counter; another assertion at N is refused; a fresh
   one lands by replacement), RD14 (the same across a restart).
8. **The CLI decides with every artifact flag on a request ref**: CD6 (91 bytes to a new file, the namespace, the
   exact next command), CD7 (`--signature --expiry`), CD8 (the decision link: challenge = the decision hash),
   CD9 (`--webauthn-assertion --expiry`, parsed before any POST); the flag rules refuse before any daemon call;
   **the real-tool siblings pass**: `ssh-keygen` signs CD6's file and it verifies; **live Chrome** signs CD8's link
   and the daemon's law decides Chrome's bytes (M6f, with its env).
9. **QUICKSTART §18** documents the decision path with the CLI's own flags (`TestQuickstartSection18FlagsMatchTheCLI`).
10. **Inventories**: `TestFrozenV1RouteTableMatchesMux` 13, `TestSurfaceAppRoutes`, `TestDaemonCredentialWriterSites`
    (two sites), `TestCredentialsMethodsRequireADeadline`, `TestProductionContextRoots`, `TestDaemonDependencyAllowlist`,
    `TestWorldBoundaryDependencyAllowlist`, `TestApprovalImportsOnlyTheStandardLibrary`, `TestEveryVerbHasHelp`,
    the quickstart tests, `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` — all PASS **unedited**.
11. **Size** ≤ 2,000 (`-- cmd host docs`): 1,493 measured; each sub-milestone's numstat as §3.
12. **Every mutant in §4 behaves as expected** on the final tree, its id in the ledger; no kill rests only on an
    inventory test.
13. **The round trip is identical** (§1.2 C20) and the gates (§6) pass, `VERIFY_GO_LEGS=plain verify_go.sh` after
    commit included.
14. **No POST arm without an operator path** (§3.1).

**What an operator can do after this sprint:** decide a pending `Human.Approve` request (the ref from the packet
page) at a terminal with `y` (sprint 1–2), with an enrolled SSH key by `--sign-with` (sprint 2) or by
`--write-message` → `ssh-keygen -Y sign` anywhere → `--signature --expiry`, and **with an enrolled passkey** by
`--webauthn-link` → the passkey page (on the studio or, over the tailnet origin, the phone) → `--webauthn-assertion
--expiry`; a replay is refused; a failed world write is retried with the same artifact.

**What the attended exit run (design §10) still waits for:** step 4 as written ("open *Needs you* on the phone,
open the card, Approve with Face ID") waits for sprint 4's card and *Needs you*; steps 1–3 and 5 are runnable after
this sprint (step 5's saved body from `--write-message`/`--signature`), and step 4 can be rehearsed through the
CD8 link on the phone. §10's run itself is Mark's (D-WORLD-109).

---

## Appendix C — after this sprint (sprint 4's head and handoff)

- **Moved to sprint 4 "the surface" (S3-D1), priced:** `requests.json` `decisions` + `decisionsTruncated`
  (DMD1–DMD3; ≈ 120 Go + 200 tests), *Needs you* in `app.js` linking each decision to its card, the reserved text
  without a link, the notification for a new decision id (XD3; ≈ 80 JS + 120 tests), the decision card
  `GET /surface/decide/{ref}` + `GET /surface/decide.js` (`DECISION_FIELDS`, the hash and `requestRef` checks, the
  passkey ceremony with `challenge = DecisionMessageHash` built with **`BigInt`** for the expiry (P1), **Save the
  message** + the signing line + a paste box; PD1, PD2, XD1, XD2, XD4; ≈ 300 page + 60 Go + 300 tests),
  `DecisionSubject.Passkeys` (sprint-1 P7, ≈ 30), `approve --list` decisions (design Q2, ≈ 80),
  `TestSurfaceAppRoutes` +2, QUICKSTART §18's card paragraph (≈ 20). **≈ 1,250 in all**, under one cap.
- **Its head:** the card reuses this sprint's POST arm unchanged (`webauthn` body = the page's envelope);
  XD1 can follow `TestApproveDecideWebAuthnLinkInChrome`'s driver use; `decide.js` repeats A3's encoding helpers
  (design §2.8) and P1's construction (C4's probe page is in `plan-scratch/p1/srv/main.go`).
- **Named survivors carried**: `D-LAND-ERR`, `C-PROOF-ERR`, `C-SIGN-NIL` (sprint 2's reasons), `C-LINK-ERR`,
  `KA-RACE-TEXT` (§4).
- **For the attended exit run** (design §10, Mark; NOT run by the loop): steps 1–3 and 5 now; step 4 after sprint 4
  (or rehearsed via CD8 on the phone).
- **D-WORLD-115** stays open (default A applied); S3-D8 names what B changes.

---

## Appendix P — the patches (apply with §2; generated from the compiled scratch stages; do not edit)

### Patch M6a — the sprint-2 judge's test gaps (JM22, K0 approver, JM01–JM04, JM46)

````diff M6a
diff --git a/cmd/ailang-worldd/approve_decide_jgaps_test.go b/cmd/ailang-worldd/approve_decide_jgaps_test.go
new file mode 100644
index 0000000..c34413d
--- /dev/null
+++ b/cmd/ailang-worldd/approve_decide_jgaps_test.go
@@ -0,0 +1,39 @@
+package main
+
+// Row 201 A4 sprint 3 M6a (the sprint-2 judge's finding 4, JM46): with no --daemon, the decision
+// mode reads the installed service's address from $HOME/.ailang-world/service.json. HOME is a temp
+// dir whose record names the fake daemon on the loopback (D-WORLD-109: no real service, no real port).
+
+import (
+	"bytes"
+	"encoding/json"
+	"os"
+	"path/filepath"
+	"testing"
+)
+
+func TestApproveDecideDefaultDaemonIsTheService(t *testing.T) {
+	f := newFakeDecider(t)
+	home := t.TempDir()
+	t.Setenv("HOME", home)
+	paths := homeServicePaths(home)
+	rec, err := json.Marshal(serviceRecord{Schema: "world/service/v1", Label: worldServiceLabel, Plist: paths.Plist, Binary: paths.Binary, Address: f.addr()})
+	if err != nil {
+		t.Fatal(err)
+	}
+	if err := os.MkdirAll(filepath.Join(home, ".ailang-world"), 0o700); err != nil {
+		t.Fatal(err)
+	}
+	if err := os.WriteFile(filepath.Join(home, ".ailang-world", "service.json"), rec, 0o600); err != nil {
+		t.Fatal(err)
+	}
+	if got := defaultServiceAddr(); got != "http://"+f.addr() { // never fall through to a real daemon's port
+		t.Fatalf("the temp service record is not read: %q", got)
+	}
+	var stdout, stderr bytes.Buffer
+	env := approveEnv(nil, nil, atTime(1710))
+	env.signDecision = (&fakeSigner{ret: func([]byte) ([]byte, error) { return []byte("ARMORED"), nil }}).sign
+	if code := runApprove([]string{f.ref, "--db", f.db, "--sign-with", "/k/id"}, &stdout, &stderr, env); code != exitOK || f.gets != 1 || len(f.posts) != 1 {
+		t.Fatalf("no --daemon: rc %d, %d GETs, %d POSTs\nstderr %s", code, f.gets, len(f.posts), stderr.String())
+	}
+}
diff --git a/host/daemon/surface_decide_race_test.go b/host/daemon/surface_decide_race_test.go
new file mode 100644
index 0000000..dc3b31d
--- /dev/null
+++ b/host/daemon/surface_decide_race_test.go
@@ -0,0 +1,27 @@
+package daemon
+
+// Row 201 A4 sprint 3 M6a (the sprint-2 judge's finding 3, JM01-JM04): the writer's approver race
+// (K3 in RecordKeyDecision's transaction) answers 422 NotApproved; any other store error is not a
+// race. Through the POST the read-only judge refuses such an approver first (plan §3.1, D-RACE), so
+// the mapping is pinned here, operand by operand, as the errors RecordKeyDecision wraps them.
+
+import (
+	"errors"
+	"fmt"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func TestIsApproverRaceOperands(t *testing.T) {
+	for _, base := range []error{store.ErrApproverRevoked, store.ErrApproverKeyChanged, store.ErrApproverUnknown, store.ErrApproverClockBehind} {
+		if err := fmt.Errorf("%w: %q", base, "mark"); !isApproverRace(err) {
+			t.Errorf("%v: not a race", err)
+		}
+	}
+	for _, err := range []error{store.ErrBadKeyDecision, store.ErrKeyDecisionLanded, errors.New("store: record key decision: disk I/O error")} {
+		if isApproverRace(err) {
+			t.Errorf("%v: a race", err)
+		}
+	}
+}
diff --git a/host/store/key_decisions_jgaps_test.go b/host/store/key_decisions_jgaps_test.go
new file mode 100644
index 0000000..7eb1b18
--- /dev/null
+++ b/host/store/key_decisions_jgaps_test.go
@@ -0,0 +1,58 @@
+package store
+
+// Row 201 A4 sprint 3 M6a (the sprint-2 judge's findings 1 and 2): LandKeyDecision lands only the
+// request it names, and K0's exact match needs the same approver. The fixtures are
+// key_decisions_test.go's; every store is under t.TempDir() (D-WORLD-109).
+
+import (
+	"reflect"
+	"testing"
+)
+
+// JM22: with two unlanded rows, landing one leaves the other unlanded (the UPDATE's request_ref
+// filter), and the other then lands with its own decision.
+func TestLandKeyDecisionLandsOnlyItsRequest(t *testing.T) {
+	c := kdStore(t)
+	other := kdRow()
+	other.RequestRef = kdD
+	for _, row := range []KeyDecisionRow{kdRow(), other} {
+		if err := c.RecordKeyDecision(testCtx(t), row, kdUse, 1000); err != nil {
+			t.Fatal(err)
+		}
+	}
+	if err := c.LandKeyDecision(testCtx(t), kdQ, kdP); err != nil {
+		t.Fatal(err)
+	}
+	got, ok, err := c.KeyDecision(testCtx(t), kdD)
+	if err != nil || !ok || got.DecisionRef != "" {
+		t.Fatalf("the other request after landing %s: %+v %v %v, want it unlanded", kdQ, got, ok, err)
+	}
+	if first, _ := kdRead(t, c); first.DecisionRef != kdP {
+		t.Fatalf("the landed request: %+v", first)
+	}
+	if err := c.LandKeyDecision(testCtx(t), kdD, kdQ); err != nil {
+		t.Fatalf("the other request lands on its own: %v", err)
+	}
+}
+
+// Finding 2 (K0's approver operand at the store API): the same artifact and signed fields recorded
+// for another approver is not the exact match: the unlanded row is replaced (the new approver,
+// the new decided_at), never kept as the first approver's.
+func TestRecordKeyDecisionExactMatchNeedsTheSameApprover(t *testing.T) {
+	c := kdStore(t)
+	if err := c.EnrolApprover(testCtx(t), ApproverRow{Name: "ann-ssh", Format: "ssh", PublicKey: []byte("k2"), EnrolledAt: 900}); err != nil {
+		t.Fatal(err)
+	}
+	if err := c.RecordKeyDecision(testCtx(t), kdRow(), kdUse, 1000); err != nil {
+		t.Fatal(err)
+	}
+	row := kdRow()
+	row.Approver = "ssh:ann-ssh"
+	if err := c.RecordKeyDecision(testCtx(t), row, ApproverUse{Name: "ann-ssh", Format: "ssh", PublicKey: []byte("k2")}, 1010); err != nil {
+		t.Fatal(err)
+	}
+	row.DecidedAt = 1010
+	if got, _ := kdRead(t, c); !reflect.DeepEqual(got, row) {
+		t.Fatalf("another approver's record: %+v, want the row replaced by %+v", got, row)
+	}
+}
````

### Patch M6b — the decision assertion law (DM7)

````diff M6b
diff --git a/host/approval/decision_assert_test.go b/host/approval/decision_assert_test.go
new file mode 100644
index 0000000..c06364e
--- /dev/null
+++ b/host/approval/decision_assert_test.go
@@ -0,0 +1,146 @@
+package approval_test
+
+// Row 201 A4 sprint 3 M6b (design §2.3, §2.7 [r2], §6 DM7, the passkey halves of DM3-DM6): the
+// decision assertion law with the approvaltest authenticator (D-WORLD-109) and, for the real tool,
+// Chrome's recorded assertion (A2's fixture): P2-P5, P3w and W1-W14 over the decision message's
+// hash, and no W15 (the counter is the credentials writer's). The A2 matrix stays unedited.
+
+import (
+	"bytes"
+	"errors"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+)
+
+var dmPasskey = approvaltest.NewPasskey(approval.AlgES256, 61)
+
+// dmHash is the decision message hash of (verb, dmP, dmQ, expiry).
+func dmHash(t *testing.T, verb byte, expiry int64) [32]byte {
+	t.Helper()
+	h, err := approval.DecisionMessageHash(verb, dmP, dmQ, expiry)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return h
+}
+
+// dmaInput is a valid decision assertion: mark-mac's passkey (enrolled at 900, stored counter 1)
+// signs approve of (P, Q, expiry) with counter count at paOrigin.
+func dmaInput(t *testing.T, count uint32, e ...approvaltest.Edit) approval.DecisionAssertionInput {
+	t.Helper()
+	a := paApprover(t, dmPasskey)
+	return approval.DecisionAssertionInput{Verb: approval.VerbApprove, Packet: dmP, Request: dmQ, Expiry: dmExpiry, Now: dmNow,
+		Response: dmPasskey.Assert(paOrigin, paRPID, dmHash(t, approval.VerbApprove, dmExpiry), count, e...), Approver: a}
+}
+
+// DM7: each of P2-P5, P3w and W1-W14 refuses an otherwise valid decision assertion, with the
+// decision challenge; a nil decision on every refusal.
+func TestDecideByAssertionSteps(t *testing.T) {
+	mint, err := approval.MessageHash(dmQ, dmP, dmExpiry)
+	if err != nil {
+		t.Fatal(err)
+	}
+	for _, tc := range []struct {
+		name string
+		edit func(t *testing.T, in *approval.DecisionAssertionInput)
+		want error
+	}{
+		{"P2-name", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Approver.Name = "Mark" }, approval.ErrApproverName},
+		{"P3w-format", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Approver.Format = approval.FormatSSH }, approval.ErrApproverNotPasskey},
+		{"P3w-no-passkey", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Approver.Passkey = nil }, approval.ErrApproverNotPasskey},
+		{"P4-revoked", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Approver.RevokedAt = 950 }, approval.ErrApproverRevoked},
+		{"P5-clock-behind", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Approver.EnrolledAt = dmNow + 1 }, approval.ErrApproverClockBehind},
+		{"W1-envelope", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Response = []byte("{}") }, approval.ErrPasskeyResponse},
+		{"W2-raw-id", func(t *testing.T, in *approval.DecisionAssertionInput) {
+			in.Response = dmPasskey.Assert(paOrigin, paRPID, dmHash(t, approval.VerbApprove, dmExpiry), 2, approvaltest.Edit{RawID: []byte("another")})
+		}, approval.ErrCredentialIDMismatch},
+		{"W3-user-handle", func(t *testing.T, in *approval.DecisionAssertionInput) {
+			in.Response = dmPasskey.Assert(paOrigin, paRPID, dmHash(t, approval.VerbApprove, dmExpiry), 2, approvaltest.Edit{UserHandle: []byte("mallory")})
+		}, approval.ErrUserHandle},
+		{"W5-type", func(t *testing.T, in *approval.DecisionAssertionInput) {
+			in.Response = dmPasskey.Assert(paOrigin, paRPID, dmHash(t, approval.VerbApprove, dmExpiry), 2, approvaltest.Edit{Type: "webauthn.create"})
+		}, approval.ErrClientDataType},
+		{"W6-the-mint-message", func(_ *testing.T, in *approval.DecisionAssertionInput) {
+			in.Response = dmPasskey.Assert(paOrigin, paRPID, mint, 2)
+		}, approval.ErrChallenge},
+		{"W6-the-other-verb", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Verb = approval.VerbDeny }, approval.ErrChallenge},
+		{"W6-another-packet", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Packet[0] ^= 1 }, approval.ErrChallenge},
+		{"W6-another-request", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Request[0] ^= 1 }, approval.ErrChallenge},
+		{"W6-another-expiry", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Expiry++ }, approval.ErrChallenge},
+		{"W6-no-message", func(_ *testing.T, in *approval.DecisionAssertionInput) { in.Verb = 3 }, approval.ErrChallenge},
+		{"W6-no-message-signed-zero", func(_ *testing.T, in *approval.DecisionAssertionInput) { // the zero hash an absent message would leave
+			in.Verb, in.Response = 3, dmPasskey.Assert(paOrigin, paRPID, [32]byte{}, 2)
+		}, approval.ErrChallenge},
+		{"W7-origin", func(t *testing.T, in *approval.DecisionAssertionInput) {
+			in.Response = dmPasskey.Assert("https://x.ts.net", paRPID, dmHash(t, approval.VerbApprove, dmExpiry), 2)
+		}, approval.ErrOriginMismatch},
+		{"W8-cross-origin", func(t *testing.T, in *approval.DecisionAssertionInput) {
+			in.Response = dmPasskey.Assert(paOrigin, paRPID, dmHash(t, approval.VerbApprove, dmExpiry), 2, approvaltest.Edit{CrossOrigin: true})
+		}, approval.ErrCrossOrigin},
+		{"W9-flags", func(t *testing.T, in *approval.DecisionAssertionInput) {
+			in.Response = dmPasskey.Assert(paOrigin, paRPID, dmHash(t, approval.VerbApprove, dmExpiry), 2, approvaltest.Edit{Set: 0x40})
+		}, approval.ErrAuthDataFlags},
+		{"W10-rp-id", func(t *testing.T, in *approval.DecisionAssertionInput) {
+			in.Response = dmPasskey.Assert(paOrigin, "y.ts.net", dmHash(t, approval.VerbApprove, dmExpiry), 2)
+		}, approval.ErrRPIDHash},
+		{"W11-user-presence", func(t *testing.T, in *approval.DecisionAssertionInput) {
+			in.Response = dmPasskey.Assert(paOrigin, paRPID, dmHash(t, approval.VerbApprove, dmExpiry), 2, approvaltest.Edit{Clear: 0x01})
+		}, approval.ErrUserPresence},
+		{"W12-user-verification", func(t *testing.T, in *approval.DecisionAssertionInput) {
+			in.Response = dmPasskey.Assert(paOrigin, paRPID, dmHash(t, approval.VerbApprove, dmExpiry), 2, approvaltest.Edit{Clear: 0x04})
+		}, approval.ErrUserVerification},
+		{"W13-backup-eligibility", func(t *testing.T, in *approval.DecisionAssertionInput) {
+			in.Response = dmPasskey.Assert(paOrigin, paRPID, dmHash(t, approval.VerbApprove, dmExpiry), 2, approvaltest.Edit{Set: 0x08})
+		}, approval.ErrBackupEligibility},
+		{"W14-bad-signature", func(t *testing.T, in *approval.DecisionAssertionInput) {
+			in.Response = dmPasskey.Assert(paOrigin, paRPID, dmHash(t, approval.VerbApprove, dmExpiry), 2, approvaltest.Edit{BadSignature: true})
+		}, approval.ErrBadSignature},
+	} {
+		in := dmaInput(t, 2)
+		tc.edit(t, &in)
+		got, err := approval.DecideByAssertion(in)
+		if !errors.Is(err, tc.want) || got != nil {
+			t.Errorf("%s: %v (decision %v), want %v", tc.name, err, got, tc.want)
+		}
+	}
+}
+
+// DM7 (W15 deferred) + the control: a valid decision assertion decides as webauthn:mark-mac, binds
+// its subject, and hands the writer the credential ID and the assertion's counter, which the law
+// did not judge: counter 1 (= the stored 1) and 0 are verified here too; the writer refuses them
+// (KD2). The mint law over the same response is refused (DM4's passkey half).
+func TestDecideByAssertionDefersTheCounter(t *testing.T) {
+	for _, count := range []uint32{2, 1, 0} {
+		in := dmaInput(t, count)
+		got, err := approval.DecideByAssertion(in)
+		if err != nil {
+			t.Fatalf("counter %d: %v", count, err)
+		}
+		id, n, ok := got.PasskeyUse()
+		if got.DecidedBy() != "webauthn:mark-mac" || got.Binds(approval.VerbApprove, dmP, dmQ, dmExpiry) != nil || !ok || n != count || !bytes.Equal(id, dmPasskey.CredentialID()) {
+			t.Fatalf("counter %d: %q, binds %v, use %x %d %v", count, got.DecidedBy(), got.Binds(approval.VerbApprove, dmP, dmQ, dmExpiry), id, n, ok)
+		}
+		if _, format, key, ok := got.Use(); format != approval.FormatWebAuthn || !ok || !bytes.Equal(key, in.Approver.Passkey.Key.Canonical()) {
+			t.Fatalf("counter %d: use %s %x %v", count, format, key, ok)
+		}
+	}
+	in := dmaInput(t, 2)
+	mint := approval.AssertionInput{Admit: approval.AdmitInput{Now: dmNow, CreatedAt: dmNow - 10, ExpiresAt: dmExpiry, StoredDigest: dmQ, BytesDigest: dmQ},
+		Challenge: dmP, Response: in.Response, Approver: in.Approver}
+	if got, err := approval.ApproveByAssertion(mint); !errors.Is(err, approval.ErrChallenge) || got != nil {
+		t.Fatalf("a decision assertion as a mint: %v", err)
+	}
+	if _, _, ok := (&approval.SignedDecision{}).PasskeyUse(); ok {
+		t.Fatal("the zero decision has a passkey use")
+	}
+	s := approvaltest.NewSigner(41)
+	ssh, err := approval.DecideBySignature(dmInput(t, s, approval.VerbApprove))
+	if err != nil {
+		t.Fatal(err)
+	}
+	if _, _, ok := ssh.PasskeyUse(); ok {
+		t.Fatal("an ssh decision has a passkey use")
+	}
+}
diff --git a/host/approval/decision_signed.go b/host/approval/decision_signed.go
index 4f6e8fc..94a7eb7 100644
--- a/host/approval/decision_signed.go
+++ b/host/approval/decision_signed.go
@@ -2,12 +2,14 @@ package approval
 
 // Row 201 A4 sprint 2 (design §2.3, §2.5 D11): the pure law of a decision signed by an enrolled SSH
 // key: the mint law's P2-P5 and VerifySSHSIG over the 91-byte message under DecisionNamespace (the
-// mint law is unchanged). The passkey half is sprint 3's (plan S2-D1). No I/O; now is an input.
+// mint law is unchanged). Sprint 3 adds the passkey half, DecideByAssertion: the assertion law's
+// W1-W14 over the decision message's hash, with W15 deferred to the credentials writer (design
+// §2.7 [r2]). No I/O; now is an input.
 
 import "bytes"
 
 // SignedDecision is a verified key decision. Its fields are unexported, so only DecideBySignature
-// can return one with ok set; the zero value binds nothing.
+// and DecideByAssertion can return one with ok set; the zero value binds nothing.
 type SignedDecision struct {
 	ok              bool
 	verb            byte
@@ -15,6 +17,8 @@ type SignedDecision struct {
 	expiry          int64
 	name, format    string
 	key             []byte
+	credentialID    []byte // webauthn: the passkey that signed
+	signCount       uint32 // webauthn: the assertion's counter, not yet judged (W15 is the writer's)
 }
 
 // DecisionSignatureInput is what DecideBySignature judges: the decision's subject (verb, packet
@@ -57,6 +61,35 @@ func DecideBySignature(in DecisionSignatureInput) (*SignedDecision, error) {
 		name: in.Approver.Name, format: in.Approver.Format, key: bytes.Clone(in.Approver.Key.Blob)}, nil
 }
 
+// DecisionAssertionInput is what DecideByAssertion judges: the decision's subject, the clock, the
+// assertion file and the enrolled approver its credential ID names, with Passkey set.
+type DecisionAssertionInput struct {
+	Verb            byte
+	Packet, Request [32]byte
+	Expiry, Now     int64
+	Response        []byte
+	Approver        Approver
+}
+
+// DecideByAssertion is the decision assertion law (design §2.3, §2.7 [r2]): the approver's state
+// (P2-P5, P3w), then the assertion law's W1-W14 with the challenge DecisionMessageHash(verb, packet,
+// request, expiry). It does not run W15: the counter is judged in RecordKeyDecision's transaction
+// against the stored counter (usePasskey), and skipped only for the exact artifact of an unlanded
+// row, so a retry after a failed world write is not refused by its own first use.
+func DecideByAssertion(in DecisionAssertionInput) (*SignedDecision, error) {
+	if err := checkApprover(AdmitInput{Now: in.Now}, in.Approver, FormatWebAuthn); err != nil { // P2-P5, P3w
+		return nil, err
+	}
+	h, hErr := DecisionMessageHash(in.Verb, in.Packet, in.Request, in.Expiry)
+	ad, err := assertionSteps(in.Response, in.Approver, h, hErr) // W1-W14
+	if err != nil {
+		return nil, err
+	}
+	p := in.Approver.Passkey
+	return &SignedDecision{ok: true, verb: in.Verb, packet: in.Packet, request: in.Request, expiry: in.Expiry,
+		name: in.Approver.Name, format: FormatWebAuthn, key: p.Key.Canonical(), credentialID: bytes.Clone(p.CredentialID), signCount: ad.SignCount}, nil
+}
+
 // Binds reports whether s is a verified decision of exactly this verb, packet, request and expiry.
 func (s *SignedDecision) Binds(verb byte, packet, request [32]byte, expiry int64) error {
 	switch {
@@ -89,3 +122,13 @@ func (s *SignedDecision) Use() (name, format string, key []byte, ok bool) {
 	}
 	return s.name, s.format, bytes.Clone(s.key), true
 }
+
+// PasskeyUse is what the evidence writer binds for a passkey decision: the credential ID (a copy)
+// and the assertion's counter, which the writer judges (W15) against the stored one; ok only for a
+// verified webauthn decision.
+func (s *SignedDecision) PasskeyUse() (credentialID []byte, signCount uint32, ok bool) {
+	if s == nil || !s.ok || s.format != FormatWebAuthn {
+		return nil, 0, false
+	}
+	return bytes.Clone(s.credentialID), s.signCount, true
+}
diff --git a/host/approval/passkey_assert.go b/host/approval/passkey_assert.go
index 4aa6843..7b80cf2 100644
--- a/host/approval/passkey_assert.go
+++ b/host/approval/passkey_assert.go
@@ -37,59 +37,73 @@ func ApproveByAssertion(in AssertionInput) (*SignedApproval, error) {
 	if err := checkApprover(in.Admit, in.Approver, FormatWebAuthn); err != nil { // P2-P5, P3w
 		return nil, err
 	}
-	p := in.Approver.Passkey
-	resp, err := ParseAssertionResponse(in.Response) // W1
+	h, hErr := MessageHash(in.Admit.StoredDigest, in.Challenge, in.Admit.ExpiresAt)
+	ad, err := assertionSteps(in.Response, in.Approver, h, hErr) // W1-W14
 	if err != nil {
 		return nil, err
 	}
+	p := in.Approver.Passkey
+	if (ad.SignCount != 0 || p.SignCount != 0) && ad.SignCount <= p.SignCount { // W15
+		return nil, fmt.Errorf("%w: enrolled %d, this assertion %d; if you did not cause this, run: ailang-worldd approver revoke %s",
+			ErrSignCountRegression, p.SignCount, ad.SignCount, in.Approver.Name)
+	}
+	return &SignedApproval{ok: true, digest: in.Admit.StoredDigest, challenge: in.Challenge,
+		expiry: in.Admit.ExpiresAt, name: in.Approver.Name, format: FormatWebAuthn, key: p.Key.Canonical(),
+		credentialID: bytes.Clone(p.CredentialID), signCount: ad.SignCount}, nil
+}
+
+// assertionSteps is the assertion law's W1-W14 for approver a over the expected challenge hash h
+// (hErr: h could not be computed, refused at W6, in W6's place), shared by ApproveByAssertion (the
+// mint message) and DecideByAssertion (the decision message, row 201 A4 §2.3). It returns the
+// verified authenticator data; the counter (W15) is the caller's. The caller has run P2-P5 and
+// P3w, so a.Passkey is set.
+func assertionSteps(response []byte, a Approver, h [32]byte, hErr error) (AuthData, error) {
+	p := a.Passkey
+	resp, err := ParseAssertionResponse(response) // W1
+	if err != nil {
+		return AuthData{}, err
+	}
 	if !bytes.Equal(resp.RawID, p.CredentialID) { // W2
-		return nil, ErrCredentialIDMismatch
+		return AuthData{}, ErrCredentialIDMismatch
 	}
-	if resp.HasUserHandle && !bytes.Equal(resp.UserHandle, []byte(in.Approver.Name)) { // W3
-		return nil, ErrUserHandle
+	if resp.HasUserHandle && !bytes.Equal(resp.UserHandle, []byte(a.Name)) { // W3
+		return AuthData{}, ErrUserHandle
 	}
 	cd, err := ParseClientData(resp.ClientDataJSON) // W4
 	if err != nil {
-		return nil, err
+		return AuthData{}, err
 	}
 	if cd.Type != "webauthn.get" { // W5
-		return nil, ErrClientDataType
+		return AuthData{}, ErrClientDataType
 	}
-	h, err := MessageHash(in.Admit.StoredDigest, in.Challenge, in.Admit.ExpiresAt)
-	if err != nil || cd.Challenge != b64Challenge(h[:]) { // W6
-		return nil, ErrChallenge
+	if hErr != nil || cd.Challenge != b64Challenge(h[:]) { // W6
+		return AuthData{}, ErrChallenge
 	}
 	if cd.Origin != p.Origin { // W7
-		return nil, ErrOriginMismatch
+		return AuthData{}, ErrOriginMismatch
 	}
 	if cd.CrossOrigin || cd.HasTopOrigin { // W8
-		return nil, ErrCrossOrigin
+		return AuthData{}, ErrCrossOrigin
 	}
 	ad, err := ParseAuthData(resp.AuthenticatorData, false) // W9
 	if err != nil {
-		return nil, err
+		return AuthData{}, err
 	}
 	if ad.RPIDHash != sha256.Sum256([]byte(p.RPID)) { // W10
-		return nil, ErrRPIDHash
+		return AuthData{}, ErrRPIDHash
 	}
 	if !ad.UP { // W11
-		return nil, ErrUserPresence
+		return AuthData{}, ErrUserPresence
 	}
 	if !ad.UV { // W12
-		return nil, ErrUserVerification
+		return AuthData{}, ErrUserVerification
 	}
 	if ad.BE != p.BackupEligible { // W13
-		return nil, ErrBackupEligibility
+		return AuthData{}, ErrBackupEligibility
 	}
 	cdh := sha256.Sum256(resp.ClientDataJSON)
 	if err := p.Key.verify(append(bytes.Clone(resp.AuthenticatorData), cdh[:]...), resp.Signature); err != nil { // W14
-		return nil, err
+		return AuthData{}, err
 	}
-	if (ad.SignCount != 0 || p.SignCount != 0) && ad.SignCount <= p.SignCount { // W15
-		return nil, fmt.Errorf("%w: enrolled %d, this assertion %d; if you did not cause this, run: ailang-worldd approver revoke %s",
-			ErrSignCountRegression, p.SignCount, ad.SignCount, in.Approver.Name)
-	}
-	return &SignedApproval{ok: true, digest: in.Admit.StoredDigest, challenge: in.Challenge,
-		expiry: in.Admit.ExpiresAt, name: in.Approver.Name, format: FormatWebAuthn, key: p.Key.Canonical(),
-		credentialID: bytes.Clone(p.CredentialID), signCount: ad.SignCount}, nil
+	return ad, nil
 }
````

### Patch M6c — the passkey judge and the writer's counter (JD passkey, KD2, KD7, K2)

````diff M6c
diff --git a/host/authority/judge_decision.go b/host/authority/judge_decision.go
index f24c36f..59fd81a 100644
--- a/host/authority/judge_decision.go
+++ b/host/authority/judge_decision.go
@@ -2,7 +2,8 @@ package authority
 
 // Row 201 A4 sprint 2 (design §2.7): the decision judge beside the mint judge (judge.go, unchanged),
 // on a read-only handle, with judge.go's step names, over a decision's subject rather than a
-// MintRequestRow. The passkey judge (W15 deferred to the writer) is sprint 3's (plan S2-D1).
+// MintRequestRow. JudgeDecisionAssertion (sprint 3) defers the counter (W15) to the credentials
+// writer's transaction (design §2.7 [r2]).
 
 import (
 	"context"
@@ -43,3 +44,35 @@ func JudgeDecisionSSH(ctx context.Context, ro *store.Credentials, in DecisionJud
 	}
 	return signed, nil
 }
+
+// JudgeDecisionAssertion judges a passkey decision assertion: ParseAssertionResponse, the passkey its
+// credential ID names (PasskeyByCredentialID on ro), ParseStoredCOSEKey, then
+// approval.DecideByAssertion, which runs W1-W14 and not W15: the counter is judged by
+// RecordKeyDecision against the stored counter, off the exact-match retry. It returns the verified
+// decision or a *JudgeError, with judge.go's step names.
+func JudgeDecisionAssertion(ctx context.Context, ro *store.Credentials, in DecisionJudgeInput, response []byte) (*approval.SignedDecision, error) {
+	resp, err := approval.ParseAssertionResponse(response)
+	if err != nil {
+		return nil, &JudgeError{Step: JudgeStepParse, Err: err}
+	}
+	a, p, ok, err := ro.PasskeyByCredentialID(ctx, resp.RawID)
+	if err != nil {
+		return nil, &JudgeError{Step: JudgeStepStore, Err: err}
+	}
+	if !ok {
+		return nil, &JudgeError{Step: JudgeStepUnknownKey, RawID: resp.RawID, Err: ErrNotEnrolled}
+	}
+	key, err := approval.ParseStoredCOSEKey(a.PublicKey)
+	if err != nil {
+		return nil, &JudgeError{Step: JudgeStepStoredKey, Approver: a.Name, Err: err}
+	}
+	signed, err := approval.DecideByAssertion(approval.DecisionAssertionInput{Verb: in.Verb, Packet: in.Packet, Request: in.Request,
+		Expiry: in.Expiry, Now: in.Now, Response: response,
+		Approver: approval.Approver{Name: a.Name, Format: a.Format, EnrolledAt: a.EnrolledAt, LastUsedAt: a.LastUsedAt, RevokedAt: a.RevokedAt,
+			Passkey: &approval.Passkey{RPID: p.RPID, Origin: p.Origin, CredentialID: p.CredentialID, Key: key,
+				BackupEligible: p.BackupEligible, SignCount: uint32(p.SignCount)}}})
+	if err != nil {
+		return nil, &JudgeError{Step: JudgeStepLaw, Err: err}
+	}
+	return signed, nil
+}
diff --git a/host/authority/judge_decision_assert_test.go b/host/authority/judge_decision_assert_test.go
new file mode 100644
index 0000000..8da3791
--- /dev/null
+++ b/host/authority/judge_decision_assert_test.go
@@ -0,0 +1,83 @@
+package authority
+
+// Row 201 A4 sprint 3 M6c (design §2.7 [r2], §6 JD1-JD3, the passkey half): JudgeDecisionAssertion's
+// steps, causes and fields agree with the mint judge's (JudgeAssertion): parse, store, unknown key,
+// stored key, law; it does not judge the counter (W15 is the writer's). Every passkey is
+// approvaltest's (D-WORLD-109); the rig is judge_test.go's (mark-mac enrolled with counter 1).
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
+// jdAssert is pk's assertion of in's decision message hash with counter count.
+func jdAssert(t *testing.T, pk *approvaltest.Passkey, in DecisionJudgeInput, count uint32) []byte {
+	t.Helper()
+	h, err := approval.DecisionMessageHash(in.Verb, in.Packet, in.Request, in.Expiry)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return pk.Assert(jgOrigin, jgRPID, h, count)
+}
+
+func TestJudgeDecisionAssertion(t *testing.T) {
+	stranger := approvaltest.NewPasskey(approval.AlgEdDSA, 62)
+	mint, _ := approval.MessageHash(jdIn.Request, jdIn.Packet, jdIn.Expiry)
+	other := jdIn
+	other.Verb = approval.VerbApprove
+	for _, tc := range []struct {
+		name     string
+		setup    []string
+		artifact []byte
+		step     string
+		cause    error
+		rawID    []byte
+		approver string
+	}{
+		{name: "JD-parse", artifact: []byte("{}"), step: JudgeStepParse, cause: approval.ErrPasskeyResponse},
+		{name: "JD-store-v5", setup: []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`},
+			artifact: jdAssert(t, jgPasskey, jdIn, 2), step: JudgeStepStore},
+		{name: "JD-unknown-key", artifact: jdAssert(t, stranger, jdIn, 2), step: JudgeStepUnknownKey, cause: ErrNotEnrolled, rawID: stranger.CredentialID()},
+		{name: "JD-stored-key", setup: []string{`UPDATE approvers SET public_key = X'00' WHERE name = 'mark-mac'`}, artifact: jdAssert(t, jgPasskey, jdIn, 2),
+			step: JudgeStepStoredKey, approver: "mark-mac"},
+		{name: "JD-law-mint-assertion", artifact: jgPasskey.Assert(jgOrigin, jgRPID, mint, 2), step: JudgeStepLaw, cause: approval.ErrChallenge},
+		{name: "JD-law-other-verb", artifact: jdAssert(t, jgPasskey, other, 2), step: JudgeStepLaw, cause: approval.ErrChallenge},
+		{name: "JD-law-revoked", setup: []string{`UPDATE approvers SET revoked_at = 950 WHERE name = 'mark-mac'`}, artifact: jdAssert(t, jgPasskey, jdIn, 2),
+			step: JudgeStepLaw, cause: approval.ErrApproverRevoked},
+		{name: "JD-law-backup-eligible", setup: []string{`UPDATE webauthn_credentials SET backup_eligible = 1 WHERE name = 'mark-mac'`}, artifact: jdAssert(t, jgPasskey, jdIn, 2),
+			step: JudgeStepLaw, cause: approval.ErrBackupEligibility},
+		{name: "JD-law-origin", setup: []string{`UPDATE webauthn_credentials SET origin = 'https://x.ts.net' WHERE name = 'mark-mac'`}, artifact: jdAssert(t, jgPasskey, jdIn, 2),
+			step: JudgeStepLaw, cause: approval.ErrOriginMismatch},
+	} {
+		path, _ := jgRig(t)
+		jgExec(t, path, tc.setup...)
+		ro, err := store.OpenCredentialsReadOnly(path)
+		if err != nil {
+			t.Fatal(err)
+		}
+		d, err := JudgeDecisionAssertion(boundedTestContext(t), ro, jdIn, tc.artifact)
+		_ = ro.Close()
+		var je *JudgeError
+		if d != nil || !errors.As(err, &je) || je.Step != tc.step || (tc.cause != nil && !errors.Is(err, tc.cause)) || !bytes.Equal(je.RawID, tc.rawID) || je.Approver != tc.approver {
+			t.Errorf("%s: %v, %v; want step %s cause %v", tc.name, d, err, tc.step, tc.cause)
+		}
+	}
+	path, _ := jgRig(t)
+	ro, err := store.OpenCredentialsReadOnly(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = ro.Close() }()
+	for _, count := range []uint32{2, 1} { // the stored counter is 1: the judge leaves 1 to the writer (W15 deferred)
+		d, err := JudgeDecisionAssertion(boundedTestContext(t), ro, jdIn, jdAssert(t, jgPasskey, jdIn, count))
+		id, n, ok := d.PasskeyUse()
+		if err != nil || d.DecidedBy() != "webauthn:mark-mac" || d.Binds(jdIn.Verb, jdIn.Packet, jdIn.Request, jdIn.Expiry) != nil || !ok || n != count || !bytes.Equal(id, jgPasskey.CredentialID()) {
+			t.Fatalf("counter %d: %v, %v", count, d, err)
+		}
+	}
+}
diff --git a/host/store/key_decisions.go b/host/store/key_decisions.go
index 15adb54..e87ef71 100644
--- a/host/store/key_decisions.go
+++ b/host/store/key_decisions.go
@@ -2,7 +2,8 @@ package store
 
 // Row 201 A4 sprint 2 (design §2.7): credentials schema 8's key_decisions, the evidence of a key
 // decision. The authority is the world store's decision, never this row; decision_ref links the two
-// once it landed (NULL until then). Format-agnostic: sprint 3's passkey arm adds no column (S2-D3).
+// once it landed (NULL until then). Format-agnostic: sprint 3's passkey arm adds no column (S2-D3);
+// its counter is webauthn_credentials.sign_count, moved by usePasskey in RecordKeyDecision (K4).
 
 import (
 	"bytes"
@@ -71,6 +72,8 @@ func checkKeyDecisionRow(row KeyDecisionRow, use ApproverUse) error {
 		return fmt.Errorf("%w: verb %q", ErrBadKeyDecision, row.Verb)
 	case row.Format != use.Format || row.Approver != use.Format+":"+use.Name:
 		return fmt.Errorf("%w: %s by %q is not the approver %s:%s", ErrBadKeyDecision, row.Format, row.Approver, use.Format, use.Name)
+	case (use.Format == "webauthn") != (use.Passkey != nil): // the counter binding is exactly a passkey's (A2 C13)
+		return fmt.Errorf("%w: a %s decision with passkey use %v", ErrBadKeyDecision, use.Format, use.Passkey != nil)
 	case len(row.Artifact) < 1 || len(row.Artifact) > MaxSurfaceApprovalArtifactBytes:
 		return fmt.Errorf("%w: an artifact is 1 to %d bytes, not %d", ErrBadKeyDecision, MaxSurfaceApprovalArtifactBytes, len(row.Artifact))
 	}
@@ -85,8 +88,10 @@ func scanKeyDecision(scan func(...any) error) (KeyDecisionRow, error) {
 
 // RecordKeyDecision records row for the approver use at now, in one transaction (design §2.7): K1
 // deadline, schema 8; K2 grammar; K0 a landed row refuses, an unlanded one equal in artifact and
-// signed fields is kept (the exact-match retry); K3 useApprover on both paths; K5 insert or replace
-// the unlanded row. The passkey counter (K4) is sprint 3's.
+// signed fields is kept (the exact-match retry); K3 useApprover on both paths; K4, for a passkey and
+// only off the exact match, usePasskey: the assertion law's W15 against the stored counter, in this
+// transaction (design §2.7 [r2]; the decision judge defers it), so a retry of the exact artifact of
+// an unlanded row is not refused by its own first use; K5 insert or replace the unlanded row.
 func (c *Credentials) RecordKeyDecision(ctx context.Context, row KeyDecisionRow, use ApproverUse, now int64) error {
 	if err := requireDeadline(ctx); err != nil { // K1
 		return err
@@ -115,6 +120,11 @@ func (c *Credentials) RecordKeyDecision(ctx context.Context, row KeyDecisionRow,
 	if err := useApprover(ctx, tx, use, now); err != nil { // K3, on both paths
 		return err
 	}
+	if !exact && use.Passkey != nil { // K4: W15, forward only (both zero: unchanged), off the exact match
+		if err := usePasskey(ctx, tx, use.Name, *use.Passkey); err != nil {
+			return err
+		}
+	}
 	if !exact { // K5
 		if _, err := tx.ExecContext(ctx, `INSERT INTO key_decisions (request_ref, packet_ref, verb, format, approver, artifact, expiry, decided_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
 ON CONFLICT(request_ref) DO UPDATE SET packet_ref = excluded.packet_ref, verb = excluded.verb, format = excluded.format, approver = excluded.approver,
diff --git a/host/store/key_decisions_passkey_test.go b/host/store/key_decisions_passkey_test.go
new file mode 100644
index 0000000..3a4f164
--- /dev/null
+++ b/host/store/key_decisions_passkey_test.go
@@ -0,0 +1,113 @@
+package store
+
+// Row 201 A4 sprint 3 M6c (design §2.7 [r2], §6 KD2, KD7 passkey half, K2's passkey rule):
+// RecordKeyDecision's counter (K4: usePasskey, the assertion law's W15 against the stored counter)
+// off the exact match, and only the counter skipped on it. The fixtures are passkey_consume_test.go's
+// (approver x, credential credX); every store is under t.TempDir() (D-WORLD-109).
+
+import (
+	"errors"
+	"reflect"
+	"testing"
+)
+
+// pkdUse is x's passkey use with counter count; pkdRow its decision row with artifact a.
+func pkdUse(count int64) ApproverUse {
+	return ApproverUse{Name: "x", Format: "webauthn", PublicKey: []byte("cose-x"), Passkey: &PasskeyUse{CredentialID: credX, SignCount: count}}
+}
+
+func pkdRow(a string) KeyDecisionRow {
+	return KeyDecisionRow{RequestRef: kdQ, PacketRef: kdP, Verb: "approve", Format: "webauthn", Approver: "webauthn:x", Artifact: []byte(a), Expiry: 1760}
+}
+
+// KD2: a counter that does not advance past the stored one is ErrSignCountRegression with nothing
+// written (no row, the counter and last_used_at unchanged); an advancing one is recorded and moves
+// it; stored 0 and assertion 0 (an authenticator without a counter) is recorded, counter still 0.
+func TestRecordKeyDecisionCounter(t *testing.T) {
+	c, _ := passkeyAt(t, credX, 7)
+	for _, n := range []int64{7, 6} {
+		if err := c.RecordKeyDecision(testCtx(t), pkdRow("asr-1"), pkdUse(n), 1000); !errors.Is(err, ErrSignCountRegression) {
+			t.Fatalf("counter %d over stored 7: %v", n, err)
+		}
+	}
+	if count, used := passkeyState(t, c); count != 7 || used != 0 {
+		t.Fatalf("refused: counter %d, last used %d", count, used)
+	}
+	if _, ok := kdRead(t, c); ok {
+		t.Fatal("refused: a row was written")
+	}
+	if err := c.RecordKeyDecision(testCtx(t), pkdRow("asr-1"), pkdUse(8), 1000); err != nil {
+		t.Fatal(err)
+	}
+	if count, used := passkeyState(t, c); count != 8 || used != 1000 {
+		t.Fatalf("recorded: counter %d, last used %d", count, used)
+	}
+	z, _ := passkeyAt(t, credX, 0)
+	if err := z.RecordKeyDecision(testCtx(t), pkdRow("asr-0"), pkdUse(0), 1000); err != nil {
+		t.Fatalf("counter 0 over stored 0: %v", err)
+	}
+	if count, _ := passkeyState(t, z); count != 0 {
+		t.Fatalf("counter 0: stored %d", count)
+	}
+}
+
+// KD7 (passkey half): an unlanded row from a counter-7 assertion (stored now 7); the same artifact
+// retried → the exact match: recorded without the counter (still 7), the row kept as first
+// recorded; a second, independently made assertion that also carries counter 7 → ErrSignCountRegression,
+// the row unchanged; one carrying 8 replaces it (K5) and moves the counter.
+func TestRecordKeyDecisionIdempotentPasskeyRetry(t *testing.T) {
+	c, _ := passkeyAt(t, credX, 6)
+	if err := c.RecordKeyDecision(testCtx(t), pkdRow("asr-7"), pkdUse(7), 1000); err != nil {
+		t.Fatal(err)
+	}
+	first, _ := kdRead(t, c)
+	if err := c.RecordKeyDecision(testCtx(t), pkdRow("asr-7"), pkdUse(7), 1010); err != nil {
+		t.Fatalf("the exact retry: %v", err)
+	}
+	if got, _ := kdRead(t, c); !reflect.DeepEqual(got, first) {
+		t.Fatalf("the exact retry rewrote the row: %+v, was %+v", got, first)
+	}
+	if count, used := passkeyState(t, c); count != 7 || used != 1010 {
+		t.Fatalf("the exact retry: counter %d, last used %d (useApprover runs, usePasskey does not)", count, used)
+	}
+	if err := c.RecordKeyDecision(testCtx(t), pkdRow("asr-7b"), pkdUse(7), 1020); !errors.Is(err, ErrSignCountRegression) {
+		t.Fatalf("another assertion at 7: %v", err)
+	}
+	if got, _ := kdRead(t, c); !reflect.DeepEqual(got, first) {
+		t.Fatalf("a refused assertion changed the row: %+v", got)
+	}
+	if err := c.RecordKeyDecision(testCtx(t), pkdRow("asr-8"), pkdUse(8), 1030); err != nil {
+		t.Fatal(err)
+	}
+	if got, _ := kdRead(t, c); string(got.Artifact) != "asr-8" || got.DecidedAt != 1030 {
+		t.Fatalf("a fresh assertion: %+v", got)
+	}
+	if count, _ := passkeyState(t, c); count != 8 {
+		t.Fatalf("a fresh assertion: counter %d", count)
+	}
+}
+
+// K2 (passkey rule): a webauthn decision without a passkey use, and an ssh one with one, are
+// ErrBadKeyDecision; nothing is written and the counter does not move.
+func TestRecordKeyDecisionPasskeyUseMatchesTheFormat(t *testing.T) {
+	c, _ := passkeyAt(t, credX, 1)
+	noPasskey := pkdUse(2)
+	noPasskey.Passkey = nil
+	if err := c.RecordKeyDecision(testCtx(t), pkdRow("asr"), noPasskey, 1000); !errors.Is(err, ErrBadKeyDecision) {
+		t.Errorf("webauthn without a passkey use: %v", err)
+	}
+	if err := c.EnrolApprover(testCtx(t), ApproverRow{Name: "mark-ssh", Format: "ssh", PublicKey: []byte("k1"), EnrolledAt: 900}); err != nil {
+		t.Fatal(err)
+	}
+	withPasskey := kdUse
+	withPasskey.Passkey = &PasskeyUse{CredentialID: credX, SignCount: 2}
+	if err := c.RecordKeyDecision(testCtx(t), kdRow(), withPasskey, 1000); !errors.Is(err, ErrBadKeyDecision) {
+		t.Errorf("ssh with a passkey use: %v", err)
+	}
+	if _, ok := kdRead(t, c); ok {
+		t.Fatal("a refused row was written")
+	}
+	if count, _ := passkeyState(t, c); count != 1 {
+		t.Fatalf("the counter moved: %d", count)
+	}
+}
````

### Patch M6d — the POST's passkey arm (RD10/RD11, KD8, RD13 b/c, RD14)

````diff M6d
diff --git a/host/daemon/surface_decide.go b/host/daemon/surface_decide.go
index df32a1b..2e36d26 100644
--- a/host/daemon/surface_decide.go
+++ b/host/daemon/surface_decide.go
@@ -176,10 +176,10 @@ func (d *Daemon) handleSurfaceDecide(w http.ResponseWriter, r *http.Request, o [
 		return
 	}
 	q := hashref.MustParse(ref) // D1 held
-	// 2. the arm (D11): a key's judge and evidence (sprint 2, surface_decide_key.go), or the typed
-	// proof, whose phase carries the verb.
+	// 2. the arm (D11): a key's judge and evidence (an SSHSIG, sprint 2; a passkey, sprint 3;
+	// surface_decide_key.go), or the typed proof, whose phase carries the verb.
 	decidedBy, evidence := typedDecider, (*store.Credentials)(nil)
-	if body.SSHSIG != nil {
+	if body.SSHSIG != nil || body.WebAuthn != nil {
 		var ok bool
 		if evidence, decidedBy, ok = d.keyArm(ctx, w, r, f, body, now); !ok {
 			return
diff --git a/host/daemon/surface_decide_key.go b/host/daemon/surface_decide_key.go
index 29a7c0a..62a07f9 100644
--- a/host/daemon/surface_decide_key.go
+++ b/host/daemon/surface_decide_key.go
@@ -2,7 +2,8 @@ package daemon
 
 // Row 201 A4 sprint 2 (design §2.5 D11, §2.7): the POST's key arm, in sprint 1's arm slot (S1-D9):
 // after the law, under humanMu, the decision judge on a read-only handle, Binds, then a non-waiting
-// writer for RecordKeyDecision (K0-K5). SSHSIG only; the passkey arm is sprint 3's (plan S2-D1).
+// writer for RecordKeyDecision (K0-K5). An SSHSIG (sprint 2) or a passkey assertion (sprint 3),
+// whose counter the writer judges (K4; the judge defers W15, design §2.7 [r2]).
 
 import (
 	"context"
@@ -24,18 +25,31 @@ func isApproverRace(err error) bool {
 		errors.Is(err, store.ErrApproverUnknown) || errors.Is(err, store.ErrApproverClockBehind)
 }
 
+// isCounterRefusal is a passkey decision the writer's counter refused (K4, the deferred W15): a
+// counter that did not advance past the stored one, or a passkey replaced since the judge.
+func isCounterRefusal(err error) bool {
+	return errors.Is(err, store.ErrSignCountRegression) || errors.Is(err, store.ErrPasskeyChanged)
+}
+
 // keyArm is D11 for a key format. It returns the open credentials writer (the caller closes it
 // after LandKeyDecision) and the decidedBy the verified key derives, or ok false when it has
 // answered a refusal: nothing is decided, and nothing is recorded unless RecordKeyDecision committed.
 func (d *Daemon) keyArm(ctx context.Context, w http.ResponseWriter, r *http.Request, f broker.DecisionFacts, body surfacepage.DecideBody, now int64) (*store.Credentials, string, bool) {
 	verb, _ := approval.DecisionVerb(body.Verb) // D9 held
-	artifact := []byte(*body.SSHSIG)
+	in := authority.DecisionJudgeInput{Verb: verb, Packet: f.Packet, Request: f.Request, Expiry: body.Expiry, Now: now}
+	judge, what := authority.JudgeDecisionSSH, "SSH signature"
+	var artifact []byte
+	if body.WebAuthn != nil { // the grammar admits exactly one key format
+		judge, artifact, what = authority.JudgeDecisionAssertion, []byte(*body.WebAuthn), "passkey assertion"
+	} else {
+		artifact = []byte(*body.SSHSIG)
+	}
 	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
 	if err != nil {
 		d.approvalInternal(w, r, err)
 		return nil, "", false
 	}
-	signed, err := authority.JudgeDecisionSSH(ctx, ro, authority.DecisionJudgeInput{Verb: verb, Packet: f.Packet, Request: f.Request, Expiry: body.Expiry, Now: now}, artifact)
+	signed, err := judge(ctx, ro, in, artifact)
 	_ = ro.Close()
 	if err == nil {
 		err = signed.Binds(verb, f.Packet, f.Request, body.Expiry)
@@ -46,7 +60,7 @@ func (d *Daemon) keyArm(ctx context.Context, w http.ResponseWriter, r *http.Requ
 		d.approvalStoreError(w, r, err)
 		return nil, "", false
 	case err != nil:
-		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the SSH signature does not decide request "+f.Ref+": "+err.Error())
+		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the "+what+" does not decide request "+f.Ref+": "+err.Error())
 		return nil, "", false
 	}
 	wr, err := store.OpenCredentials(d.credentialsPath) // non-waiting: a busy writer is 503 Busy
@@ -55,13 +69,20 @@ func (d *Daemon) keyArm(ctx context.Context, w http.ResponseWriter, r *http.Requ
 		return nil, "", false
 	}
 	name, format, key, _ := signed.Use()
+	use := store.ApproverUse{Name: name, Format: format, PublicKey: key}
+	if id, count, ok := signed.PasskeyUse(); ok { // K4: the writer judges the counter, off the exact-match retry
+		use.Passkey = &store.PasskeyUse{CredentialID: id, SignCount: int64(count)}
+	}
 	err = wr.RecordKeyDecision(ctx, store.KeyDecisionRow{RequestRef: f.Ref, PacketRef: "sha256:" + hex.EncodeToString(f.Packet[:]), Verb: body.Verb,
-		Format: format, Approver: signed.DecidedBy(), Artifact: artifact, Expiry: body.Expiry}, store.ApproverUse{Name: name, Format: format, PublicKey: key}, now)
+		Format: format, Approver: signed.DecidedBy(), Artifact: artifact, Expiry: body.Expiry}, use, now)
 	if err != nil {
 		_ = wr.Close()
-		if isApproverRace(err) {
-			writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the approver changed since the signature was judged: "+err.Error())
-		} else {
+		switch {
+		case isApproverRace(err):
+			writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the approver changed since the "+what+" was judged: "+err.Error())
+		case isCounterRefusal(err):
+			writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the "+what+" does not decide request "+f.Ref+": "+err.Error())
+		default:
 			d.approvalStoreError(w, r, err)
 		}
 		return nil, "", false
diff --git a/host/daemon/surface_decide_passkey_test.go b/host/daemon/surface_decide_passkey_test.go
new file mode 100644
index 0000000..7f46abd
--- /dev/null
+++ b/host/daemon/surface_decide_passkey_test.go
@@ -0,0 +1,251 @@
+package daemon
+
+// Row 201 A4 sprint 3 M6d (design §2.7 [r2], §2.8, §6 RD10 and RD11 passkey halves, KD2 and KD7 end
+// to end, KD8, RD13 (b) (c), RD14): POST /surface/decide/{ref} with a passkey assertion by an
+// enrolled approvaltest authenticator (D-WORLD-109), through d.Handler(). The judge defers the
+// counter (W15); the writer judges it against the stored one, off the exact-match retry.
+
+import (
+	"context"
+	"database/sql"
+	"encoding/json"
+	"fmt"
+	"io"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+var kdPasskey = approvaltest.NewPasskey(approval.AlgES256, 72)
+
+const kdOrigin, kdRP = "http://localhost:7644", "localhost"
+
+// enrolPasskeyDecider enrols pk as name at packetNow (before the rig's clock) with counter stored.
+func (r decideRig) enrolPasskeyDecider(name string, pk *approvaltest.Passkey, stored int64) {
+	r.t.Helper()
+	k := pk.Key()
+	r.write(func(ctx context.Context, c *store.Credentials) error {
+		return c.EnrolPasskey(ctx, store.ApproverRow{Name: name, Format: approval.FormatWebAuthn, PublicKey: k.Canonical(), EnrolledAt: packetNow},
+			store.PasskeyRow{Name: name, RPID: kdRP, Origin: kdOrigin, CredentialID: pk.CredentialID(), AAGUID: make([]byte, 16), Alg: k.Alg, SignCount: stored})
+	})
+}
+
+// decidePasskeyBody is a POST body deciding verb on (p, q, expiry) with pk's assertion of the
+// decision message hash with counter count (and edits e).
+func decidePasskeyBody(t *testing.T, pk *approvaltest.Passkey, verb string, p store.Object, q hashref.HashRef, expiry int64, count uint32, e ...approvaltest.Edit) []byte {
+	t.Helper()
+	v, _ := approval.DecisionVerb(verb)
+	h, err := approval.DecisionMessageHash(v, [32]byte(mustDigest(t, p.Hash)), [32]byte(mustDigest(t, q)), expiry)
+	if err != nil {
+		t.Fatal(err)
+	}
+	b, _ := json.Marshal(map[string]any{"v": 1, "verb": verb, "packet": p.Hash.Digest(), "request": q.Digest(), "expiry": expiry,
+		"webauthn": json.RawMessage(pk.Assert(kdOrigin, kdRP, h, count, e...))})
+	return b
+}
+
+// signCount is name's stored passkey counter.
+func (r decideRig) signCount(name string) int64 {
+	r.t.Helper()
+	ro, err := store.OpenCredentialsReadOnly(r.d.credentialsPath)
+	if err != nil {
+		r.t.Fatal(err)
+	}
+	defer func() { _ = ro.Close() }()
+	p, ok, err := ro.Passkey(boundedTestContext(r.t), name)
+	if err != nil || !ok {
+		r.t.Fatal(name, ok, err)
+	}
+	return p.SignCount
+}
+
+// failWorldWrites makes every approval decision write fail until the returned func drops the trigger.
+func (r decideRig) failWorldWrites() func() {
+	r.t.Helper()
+	db, err := sql.Open("sqlite", r.d.cfg.DBPath)
+	if err != nil {
+		r.t.Fatal(err)
+	}
+	if _, err := db.Exec(`CREATE TRIGGER fail_decision BEFORE INSERT ON objects WHEN NEW.semantic_id = '` + broker.ApprovalDecisionV1 + `' BEGIN SELECT RAISE(ABORT, 'injected world-write failure'); END`); err != nil {
+		r.t.Fatal(err)
+	}
+	return func() {
+		r.t.Helper()
+		defer func() { _ = db.Close() }()
+		if _, err := db.Exec(`DROP TRIGGER fail_decision`); err != nil {
+			r.t.Fatal(err)
+		}
+	}
+}
+
+// RD10, RD11 (passkey halves), KD2 end to end: an enrolled passkey's assertion decides approve and
+// deny as webauthn:<name> (derived; a body naming decidedBy is 400), the evidence row is landed and
+// the counter moved; a webauthn member beside another format is 400; a mint assertion and a
+// counter that did not advance are 422 with nothing decided, recorded or moved.
+func TestDecidePOSTPasskey(t *testing.T) {
+	r := newDecideRig(t)
+	r.enrolPasskeyDecider("mark-mac", kdPasskey, 1)
+	q, p := r.packet("passkey", broker.WorkspaceScope, "agent")
+	body := decidePasskeyBody(t, kdPasskey, "approve", p, q, decNow+50, 2)
+	for name, b := range map[string][]byte{
+		"decidedBy":          []byte(strings.Replace(string(body), `"v":1`, `"v":1,"decidedBy":"webauthn:mallory"`, 1)),
+		"webauthn and ssh":   []byte(strings.Replace(string(body), `"v":1`, `"v":1,"sshsig":"-----BEGIN SSH SIGNATURE-----"`, 1)),
+		"webauthn and typed": []byte(strings.Replace(string(body), `"v":1`, `"v":1,"typed":true`, 1)),
+	} {
+		requireRefusal(t, name, r.postKey(q, b), 400, "BadRequest")
+	}
+	mint, _ := approval.MessageHash([32]byte(mustDigest(t, q)), [32]byte(mustDigest(t, p.Hash)), decNow+50)
+	mintBody, _ := json.Marshal(map[string]any{"v": 1, "verb": "approve", "packet": p.Hash.Digest(), "request": q.Digest(), "expiry": decNow + 50,
+		"webauthn": json.RawMessage(kdPasskey.Assert(kdOrigin, kdRP, mint, 2))})
+	if got := requireRefusal(t, "a mint assertion", r.postKey(q, mintBody), 422, "NotApproved"); !strings.Contains(got.Reason, "passkey assertion does not decide") {
+		t.Errorf("a mint assertion: %q", got.Reason)
+	}
+	if n, _ := r.decisions(); n != 0 || r.signCount("mark-mac") != 1 {
+		t.Fatalf("refused: %d decisions, counter %d", n, r.signCount("mark-mac"))
+	}
+	rec := r.postKey(q, body)
+	var made surfacepage.DecisionMade
+	if err := json.Unmarshal(rec.Body.Bytes(), &made); rec.Code != 200 || err != nil || made.DecidedBy != "webauthn:mark-mac" || made.Decided != "approve" {
+		t.Fatalf("approve: %d %s", rec.Code, rec.Body.String())
+	}
+	row, ok := r.keyRow(q)
+	if !ok || row.DecisionRef != made.Decision || row.Format != "webauthn" || row.Approver != "webauthn:mark-mac" || r.signCount("mark-mac") != 2 {
+		t.Fatalf("evidence %+v %v, counter %d", row, ok, r.signCount("mark-mac"))
+	}
+	q2, p2 := r.packet("passkey-deny", broker.WorkspaceScope, "agent-2")
+	if rec := r.postKey(q2, decidePasskeyBody(t, kdPasskey, "deny", p2, q2, decNow+50, 3)); rec.Code != 200 || !strings.Contains(rec.Body.String(), `"decided":"deny"`) {
+		t.Fatalf("deny: %d %s", rec.Code, rec.Body.String())
+	}
+	q3, p3 := r.packet("passkey-regress", broker.WorkspaceScope, "agent-3")
+	n, head := r.decisions()
+	if got := requireRefusal(t, "counter 3 over stored 3", r.postKey(q3, decidePasskeyBody(t, kdPasskey, "approve", p3, q3, decNow+50, 3)), 422, "NotApproved"); !strings.Contains(got.Reason, "sign counter did not advance") {
+		t.Errorf("counter: %q", got.Reason)
+	}
+	r.requireNothingDecidedSince(n, head)
+	if _, ok := r.keyRow(q3); ok || r.signCount("mark-mac") != 3 {
+		t.Fatalf("a refused counter recorded evidence or moved the counter (%d)", r.signCount("mark-mac"))
+	}
+}
+
+// KD8 and RD13 (c): an authenticator without a counter (stored 0, assertion 0) decides once: the
+// identical POST again is D6's 409 with one decision and one row; after a failed world write the
+// identical POST lands that one decision (K0's exact match) and the next is 409.
+func TestDecidePOSTCounterZeroAssertionLandsOnce(t *testing.T) {
+	r := newDecideRig(t)
+	r.enrolPasskeyDecider("mark-mac", kdPasskey, 0)
+	q, p := r.packet("zero", broker.WorkspaceScope, "agent")
+	body := decidePasskeyBody(t, kdPasskey, "approve", p, q, decNow+50, 0)
+	if rec := r.postKey(q, body); rec.Code != 200 {
+		t.Fatalf("first: %d %s", rec.Code, rec.Body.String())
+	}
+	requireRefusal(t, "the identical assertion again", r.postKey(q, body), 409, "Conflict")
+	if n, _ := r.decisions(); n != 1 {
+		t.Fatalf("%d decisions", n)
+	}
+	q2, p2 := r.packet("zero-crash", broker.WorkspaceScope, "agent-2")
+	body2 := decidePasskeyBody(t, kdPasskey, "deny", p2, q2, decNow+50, 0)
+	release := r.failWorldWrites()
+	requireRefusal(t, "failed world write", r.postKey(q2, body2), 500, "Internal")
+	if row, ok := r.keyRow(q2); !ok || row.DecisionRef != "" {
+		t.Fatalf("after the failure: row %+v %v", row, ok)
+	}
+	release()
+	if rec := r.postKey(q2, body2); rec.Code != 200 {
+		t.Fatalf("retry: %d %s", rec.Code, rec.Body.String())
+	}
+	requireRefusal(t, "after the retry", r.postKey(q2, body2), 409, "Conflict")
+	if n, _ := r.decisions(); n != 2 || r.signCount("mark-mac") != 0 {
+		t.Fatalf("%d decisions, counter %d", n, r.signCount("mark-mac"))
+	}
+}
+
+// RD13 (b), KD7 end to end: a counter-7 assertion (stored 6) whose world write fails → 500, one
+// unlanded row, the counter 7, 0 decisions; a second, independently made assertion that also
+// carries 7 → 422 (the writer's W15), nothing changed; the identical body → 200, the row landed,
+// the counter still 7. Then, for a second request, a fresh assertion at N+1 lands by replacement.
+func TestDecidePOSTPasskeyRetryAfterAFailedWorldWrite(t *testing.T) {
+	r := newDecideRig(t)
+	r.enrolPasskeyDecider("mark-mac", kdPasskey, 6)
+	q, p := r.packet("pk-retry", broker.WorkspaceScope, "agent")
+	body := decidePasskeyBody(t, kdPasskey, "approve", p, q, decNow+50, 7)
+	release := r.failWorldWrites()
+	requireRefusal(t, "failed world write", r.postKey(q, body), 500, "Internal")
+	first, ok := r.keyRow(q)
+	if n, _ := r.decisions(); n != 0 || !ok || first.DecisionRef != "" || r.signCount("mark-mac") != 7 {
+		t.Fatalf("after the failure: %d decisions, row %+v %v, counter %d", n, first, ok, r.signCount("mark-mac"))
+	}
+	release()
+	other := decidePasskeyBody(t, kdPasskey, "approve", p, q, decNow+50, 7, approvaltest.Edit{UserHandle: []byte("mark-mac")})
+	requireRefusal(t, "another assertion at 7", r.postKey(q, other), 422, "NotApproved")
+	if row, _ := r.keyRow(q); string(row.Artifact) != string(first.Artifact) {
+		t.Fatal("a refused assertion replaced the row")
+	}
+	rec := r.postKey(q, body)
+	row, _ := r.keyRow(q)
+	if n, _ := r.decisions(); rec.Code != 200 || n != 1 || row.DecisionRef == "" || row.DecidedAt != first.DecidedAt || r.signCount("mark-mac") != 7 {
+		t.Fatalf("retry: %d %s; %d decisions; row %+v; counter %d", rec.Code, rec.Body.String(), n, row, r.signCount("mark-mac"))
+	}
+	q2, p2 := r.packet("pk-fresh", broker.WorkspaceScope, "agent-2")
+	release = r.failWorldWrites()
+	requireRefusal(t, "failed world write 2", r.postKey(q2, decidePasskeyBody(t, kdPasskey, "deny", p2, q2, decNow+50, 8)), 500, "Internal")
+	release()
+	if rec := r.postKey(q2, decidePasskeyBody(t, kdPasskey, "deny", p2, q2, decNow+50, 9)); rec.Code != 200 || r.signCount("mark-mac") != 9 {
+		t.Fatalf("a fresh assertion: %d %s, counter %d", rec.Code, rec.Body.String(), r.signCount("mark-mac"))
+	}
+	if row, _ := r.keyRow(q2); row.DecisionRef == "" {
+		t.Fatalf("a fresh assertion: row %+v", row)
+	}
+}
+
+// RD14: the counter-7 retry across a daemon restart: stored 7 in SQLite after the failed world
+// write; the restarted daemon refuses another assertion carrying 7 and lands the identical one
+// through the full HTTP/judge path without moving the counter.
+func TestDecidePOSTRetryAfterAFailedWorldWriteAcrossRestart(t *testing.T) {
+	r := newDecideRig(t)
+	r.enrolPasskeyDecider("mark-mac", kdPasskey, 6)
+	q, p := r.packet("pk-restart", broker.WorkspaceScope, "agent")
+	body := decidePasskeyBody(t, kdPasskey, "approve", p, q, decNow+50, 7)
+	release := r.failWorldWrites()
+	requireRefusal(t, "failed world write", r.postKey(q, body), 500, "Internal")
+	release()
+	db := r.d.cfg.DBPath
+	if err := r.d.Close(); err != nil {
+		t.Fatal(err)
+	}
+	d2, err := New(boundedTestContext(t), Config{DBPath: db, BindHost: DefaultBindHost, ErrorLog: io.Discard})
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { _ = d2.Close() })
+	d2.approveNow = func() int64 { return decNow + 1 }
+	r2 := decideRig{approvalRig{t: t, d: d2}}
+	if r2.signCount("mark-mac") != 7 {
+		t.Fatalf("stored counter after the restart: %d", r2.signCount("mark-mac"))
+	}
+	requireRefusal(t, "another assertion at 7", r2.postKey(q, decidePasskeyBody(t, kdPasskey, "approve", p, q, decNow+50, 7, approvaltest.Edit{UserHandle: []byte("mark-mac")})), 422, "NotApproved")
+	if rec := r2.postKey(q, body); rec.Code != 200 || r2.signCount("mark-mac") != 7 {
+		t.Fatalf("the identical retry after the restart: %d %s, counter %d", rec.Code, rec.Body.String(), r2.signCount("mark-mac"))
+	}
+	if row, _ := r2.keyRow(q); row.DecisionRef == "" {
+		t.Fatalf("not landed: %+v", row)
+	}
+}
+
+// The writer's counter refusals answer 422 (both operands); an approver race is not one of them.
+// ErrPasskeyChanged is unit-only: the judge finds the passkey by its credential ID first (plan §3.1).
+func TestIsCounterRefusalOperands(t *testing.T) {
+	for _, base := range []error{store.ErrSignCountRegression, store.ErrPasskeyChanged} {
+		if err := fmt.Errorf("%w: %q", base, "mark-mac"); !isCounterRefusal(err) {
+			t.Errorf("%v: not a counter refusal", err)
+		}
+	}
+	if isCounterRefusal(store.ErrApproverRevoked) {
+		t.Error("a revoked approver is a counter refusal")
+	}
+}
diff --git a/host/surfacepage/decide_model.go b/host/surfacepage/decide_model.go
index ffada1d..1639c0f 100644
--- a/host/surfacepage/decide_model.go
+++ b/host/surfacepage/decide_model.go
@@ -1,8 +1,8 @@
 package surfacepage
 
 // Row 201 A4 sprint 1 (design_docs/planned/w-world-approvals-by-key-a4.md §2.8): the decision
-// routes' JSON shapes and the POST body's grammar, pure. Sprint 2 admits exactly one of typed and
-// sshsig; webauthn arrives with sprint 3's passkey arm, so until then it is an unknown field.
+// routes' JSON shapes and the POST body's grammar, pure. The body carries exactly one of typed,
+// sshsig (sprint 2) and webauthn (sprint 3's passkey arm).
 
 import (
 	"bytes"
@@ -53,20 +53,22 @@ type DecisionMade struct {
 }
 
 // DecideBody is the POST body (design §2.8): the verb, the packet digest P and request digest Q
-// the caller saw (64 lowercase hex), the expiry, and exactly one format: typed: true (sprint 1) or
-// sshsig, an armored SSH signature over the decision message (sprint 2).
+// the caller saw (64 lowercase hex), the expiry, and exactly one format: typed: true (sprint 1),
+// sshsig, an armored SSH signature over the decision message (sprint 2), or webauthn, a passkey
+// assertion envelope (A2 §2.7) over its hash (sprint 3).
 type DecideBody struct {
-	V       int     `json:"v"`
-	Verb    string  `json:"verb"`
-	Packet  string  `json:"packet"`
-	Request string  `json:"request"`
-	Expiry  int64   `json:"expiry"`
-	Typed   bool    `json:"typed,omitempty"`
-	SSHSIG  *string `json:"sshsig,omitempty"`
+	V        int              `json:"v"`
+	Verb     string           `json:"verb"`
+	Packet   string           `json:"packet"`
+	Request  string           `json:"request"`
+	Expiry   int64            `json:"expiry"`
+	Typed    bool             `json:"typed,omitempty"`
+	SSHSIG   *string          `json:"sshsig,omitempty"`
+	WebAuthn *json.RawMessage `json:"webauthn,omitempty"`
 }
 
 // ErrDecideBody refuses a body that is not DecideBody's grammar.
-var ErrDecideBody = errors.New(`the body is not {"v":1,"verb":…,"packet":<64 hex>,"request":<64 hex>,"expiry":…} with exactly one of "typed":true and "sshsig":<armored SSHSIG>`)
+var ErrDecideBody = errors.New(`the body is not {"v":1,"verb":…,"packet":<64 hex>,"request":<64 hex>,"expiry":…} with exactly one of "typed":true, "sshsig":<armored SSHSIG> and "webauthn":<assertion>`)
 
 // ParseDecideBody is the POST's grammar: strict JSON (no unknown field, so no decidedBy; exactly
 // one value), v 1, exactly one format, both digests 64 lowercase hex. The verb and the expiry are
@@ -83,9 +85,9 @@ func ParseDecideBody(b []byte) (body DecideBody, packet, request [32]byte, err e
 	switch {
 	case body.V != 1:
 		return DecideBody{}, packet, request, fmt.Errorf("%w: v is %d", ErrDecideBody, body.V)
-	case !body.Typed && body.SSHSIG == nil:
+	case decideFormats(body) == 0:
 		return DecideBody{}, packet, request, fmt.Errorf("%w: it has no format", ErrDecideBody)
-	case body.Typed && body.SSHSIG != nil:
+	case decideFormats(body) > 1:
 		return DecideBody{}, packet, request, fmt.Errorf("%w: it has two formats", ErrDecideBody)
 	case !lowerHex64(body.Packet) || !lowerHex64(body.Request):
 		return DecideBody{}, packet, request, fmt.Errorf("%w: packet and request are each 64 lowercase hex", ErrDecideBody)
@@ -94,3 +96,14 @@ func ParseDecideBody(b []byte) (body DecideBody, packet, request [32]byte, err e
 	request, _ = hex32(body.Request)
 	return body, packet, request, nil
 }
+
+// decideFormats is how many of typed, sshsig and webauthn body carries.
+func decideFormats(body DecideBody) int {
+	n := 0
+	for _, on := range []bool{body.Typed, body.SSHSIG != nil, body.WebAuthn != nil} {
+		if on {
+			n++
+		}
+	}
+	return n
+}
````

### Patch M6e — the CLI's artifact flags on a request ref (CD6–CD9)

````diff M6e
diff --git a/cmd/ailang-worldd/approve.go b/cmd/ailang-worldd/approve.go
index 36b9662..1754240 100644
--- a/cmd/ailang-worldd/approve.go
+++ b/cmd/ailang-worldd/approve.go
@@ -31,6 +31,8 @@ const approveHelp = `usage: ailang-worldd approve <request-id> [--db <path>] [--
        ailang-worldd approve --list [--decided] [--db <path>] [--credentials-db <path>]
        ailang-worldd approve <request-ref> [--deny] [--daemon <host:port>] [--db <path>] [--credentials-db <path>]
        ailang-worldd approve <request-ref> [--deny] --sign-with <key> [--daemon <host:port>] [--db <path>] [--credentials-db <path>]
+       ailang-worldd approve <request-ref> [--deny] --write-message <file> | --webauthn-link [--daemon <host:port>] [--db <path>] [--credentials-db <path>]
+       ailang-worldd approve <request-ref> [--deny] --signature <file> | --webauthn-assertion <file> --expiry <t> [--daemon <host:port>] [...]
 
 ATTENDED, or SIGNED. Mints the session a stored mint request describes (row
 201), after a y/N typed at /dev/tty, or, with --signature or --sign-with,
@@ -105,6 +107,13 @@ decides at once); publish and finish approvals stay with world-publish
 approve and episode finish. --sign-with <key> opens no terminal: ssh-keygen
 -Y sign -n world-decision@sunholo signs the 91-byte decision message, and the
 decision is by ssh:<name>, an approver enrolled with 'approver enroll'.
+--write-message <file> writes that message (expiry: the latest the daemon
+allows) to a NEW file and prints the signing line and the --signature <file>
+--expiry <t> step that posts it. --webauthn-link prints, per passkey origin, a
+link to the passkey page over the message's hash (the page says "Approve
+request decide-<verb>-…" for either verb: check the hash it shows), and the
+--webauthn-assertion <file> --expiry <t> step that posts the saved file; that
+decision is by webauthn:<name>. The daemon judges every artifact.
 `
 
 // approveAdmission is approve's admission of a stored row, before any terminal: the decoded
@@ -171,6 +180,7 @@ func runApprove(args []string, stdout, stderr io.Writer, env sessionEnv) int {
 	awaitFor := fs.Duration("await-for", 0, "with --await: stop waiting after this long (default: the request's expiry)")
 	deny := fs.Bool("deny", false, "with a request ref: deny it instead of approving it")
 	daemonAddr := fs.String("daemon", "", "with a request ref: the daemon serving --db (default: the installed service's address)")
+	expiry := fs.Int64("expiry", 0, "with a request ref and --signature or --webauthn-assertion: the expiry the signed decision message carries")
 	positional, err := parseInterleaved(fs, args)
 	if errors.Is(err, flag.ErrHelp) {
 		fmt.Fprint(stdout, approveHelp)
@@ -208,21 +218,25 @@ func runApprove(args []string, stdout, stderr io.Writer, env sessionEnv) int {
 		fmt.Fprint(stderr, approveHelp)
 		return exitUsage
 	}
-	if (*deny || *daemonAddr != "") && (*list || !decisionRefPattern.MatchString(positional[0])) { // CD3: --deny and --daemon decide a request ref
-		return refuse("--deny and --daemon take a Human.Approve request ref (sha256: and 64 hex); a mint request has no deny: it expires")
+	if (*deny || *daemonAddr != "" || *expiry != 0) && (*list || !decisionRefPattern.MatchString(positional[0])) { // CD3: they decide a request ref
+		return refuse("--deny, --daemon and --expiry take a Human.Approve request ref (sha256: and 64 hex); a mint request has no deny: it expires")
 	}
 	id := ""
 	if !*list {
 		id = positional[0]
 		if decisionRefPattern.MatchString(id) { // row 201 A4: a request ref selects the decision mode (design §2.10)
-			if modes != 0 && *signWith == "" { // CD3: --sign-with is the key flag that decides a request ref (plan S2-D1)
-				return refuse("a request ref is decided at once by the daemon, by the typed y or --sign-with <key>: there is nothing to await, and the other key flags do not decide it")
+			if *await { // CD3
+				return refuse("a request ref is decided at once by the daemon: there is nothing to await")
+			}
+			if (*sigPath != "" || *assertion != "") != (*expiry != 0) { // CD7, CD9: an artifact is posted with the expiry it signed
+				return refuse("--signature and --webauthn-assertion go with --expiry <t> on a request ref, and --expiry with one of them: the expiry is part of the signed message")
 			}
 			verb := "approve"
 			if *deny {
 				verb = "deny"
 			}
-			return runApproveDecide(decideOptions{ref: id, verb: verb, db: *db, credsDB: *credsDB, daemon: *daemonAddr, signWith: *signWith}, stdout, stderr, env, refuse)
+			return runApproveDecide(decideOptions{ref: id, verb: verb, db: *db, credsDB: *credsDB, daemon: *daemonAddr, signWith: *signWith,
+				sigPath: *sigPath, msgPath: *msgPath, assertion: *assertion, passkeyLink: *passkeyLink, expiry: *expiry}, stdout, stderr, env, refuse)
 		}
 		if !store.MintRequestIDPattern.MatchString(id) {
 			return refuse("%q is not a mint request id (mr- and 16 lowercase hex)", id)
diff --git a/cmd/ailang-worldd/approve_decide.go b/cmd/ailang-worldd/approve_decide.go
index 77fcc2a..6dac75b 100644
--- a/cmd/ailang-worldd/approve_decide.go
+++ b/cmd/ailang-worldd/approve_decide.go
@@ -12,12 +12,14 @@ package main
 import (
 	"bytes"
 	"context"
+	"crypto/sha256"
 	"encoding/base64"
 	"encoding/hex"
 	"encoding/json"
 	"fmt"
 	"io"
 	"net/http"
+	"os"
 	"regexp"
 	"strings"
 	"time"
@@ -39,6 +41,9 @@ const approveDecideBudget = 10 * time.Second
 // decideOptions are approve's flags in the decision mode.
 type decideOptions struct {
 	ref, verb, db, credsDB, daemon, signWith string
+	sigPath, msgPath, assertion              string // the other key flags (CD6, CD7, CD9)
+	passkeyLink                              bool   // CD8
+	expiry                                   int64  // with sigPath or assertion: the signed expiry
 }
 
 // decisionSubject is a subject the CLI checked: its request and packet, decoded from bytes that
@@ -149,6 +154,9 @@ func runApproveDecide(o decideOptions, stdout, stderr io.Writer, env sessionEnv,
 	if o.signWith != "" {
 		return decideBySignWith(o, s, c, base, env, stdout, stderr, refuse)
 	}
+	if o.sigPath != "" || o.msgPath != "" || o.assertion != "" || o.passkeyLink {
+		return decideByArtifact(ctx, o, credPath, hostPort, s, c, base, stdout, refuse)
+	}
 	const nothing = "nothing was decided; decide with an enrolled key instead (--sign-with <key> or --signature <file>)"
 	term, code := openApproverTerminal(env, stderr, "ailang-worldd approve: ", nothing)
 	if code != exitOK {
@@ -200,6 +208,104 @@ func signDecisionWithSSHKeygen(ctx context.Context, key string, msg []byte, stde
 	return sshsign.SignIn(ctx, "ssh-keygen", approval.DecisionNamespace, key, msg, stderr)
 }
 
+// decideByArtifact is the decision mode's artifact flags (design §2.10): --write-message (CD6)
+// writes the 91-byte decision message for the window's latest expiry to a NEW file and prints the
+// signing line and the --signature step; --webauthn-link (CD8) prints, per origin of the active
+// passkeys, a link to the A2 passkey page (unchanged) over the message's hash, and the
+// --webauthn-assertion step; --signature (CD7) and --webauthn-assertion (CD9) post a file made for
+// --expiry. The daemon judges the artifact (design §2.7); no terminal opens and no proof is written.
+func decideByArtifact(ctx context.Context, o decideOptions, credPath, hostPort string, s decisionSubject, c *client, base string, stdout io.Writer, refuse func(string, ...any) int) int {
+	verb, _ := approval.DecisionVerb(o.verb)
+	body := surfacepage.DecideBody{V: 1, Verb: o.verb, Packet: strings.TrimPrefix(s.Packet, "sha256:"), Request: strings.TrimPrefix(o.ref, "sha256:"), Expiry: s.Window.Max}
+	if o.expiry != 0 {
+		body.Expiry = o.expiry
+	}
+	msg, err := approval.DecisionMessage(verb, refDigest(s.Packet), refDigest(o.ref), body.Expiry)
+	if err != nil {
+		return refuse("%v", err)
+	}
+	next := fmt.Sprintf("ailang-worldd approve %s%s %%s %%s --expiry %d --daemon %s --db %s --credentials-db %s", o.ref,
+		map[string]string{"approve": "", "deny": " --deny"}[o.verb], body.Expiry, hostPort, shellQuote(o.db), shellQuote(credPath))
+	switch {
+	case o.msgPath != "": // CD6
+		f, err := os.OpenFile(o.msgPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644) // never overwrites
+		if err == nil {
+			_, err = f.Write(msg)
+			if cerr := f.Close(); err == nil {
+				err = cerr
+			}
+		}
+		if err != nil {
+			return refuse("--write-message: %v; nothing was written", err)
+		}
+		fmt.Fprintf(stdout, "wrote the %d-byte message deciding %s request %s to %s\nsign it with an enrolled key before %s:\n  ssh-keygen -Y sign -n %s -f <your key> < %s > %s.sig\nthen decide it:\n  %s\n",
+			len(msg), strings.ToUpper(o.verb), o.ref, o.msgPath, approverTime(body.Expiry), approval.DecisionNamespace, shellQuote(o.msgPath), shellQuote(o.msgPath),
+			fmt.Sprintf(next, "--signature", shellQuote(o.msgPath+".sig")))
+		return exitOK
+	case o.passkeyLink: // CD8
+		return decisionPasskeyLinks(ctx, o, credPath, body.Expiry, msg, fmt.Sprintf(next, "--webauthn-assertion", "<the saved file>"), stdout, refuse)
+	case o.sigPath != "": // CD7
+		armored, err := readRegularFileCapped(o.sigPath, approveMaxSignatureBytes)
+		if err != nil {
+			return refuse("--signature: %v; nothing was decided", err)
+		}
+		sig := string(armored)
+		body.SSHSIG = &sig
+	default: // CD9
+		file, err := readRegularFileCapped(o.assertion, approval.MaxPasskeyResponseBytes)
+		if err != nil {
+			return refuse("--webauthn-assertion: %v; nothing was decided", err)
+		}
+		if _, err := approval.ParseAssertionResponse(file); err != nil {
+			return refuse("--webauthn-assertion %s is not a passkey assertion: %v; nothing was decided", o.assertion, err)
+		}
+		asr := json.RawMessage(file)
+		body.WebAuthn = &asr
+	}
+	return postDecision(c, base, o.ref, body, "", stdout, refuse)
+}
+
+// decisionPasskeyLinks is --webauthn-link on a request ref (CD8): one get link per distinct origin
+// among the active passkeys (read-only), each allowing that origin's credentials, whose challenge is
+// SHA-256 of the decision message. The A2 page signs the challenge its link carries and reads
+// "Approve request <request>": the link's request names the verb (decide-<verb>-<16 hex of Q>), and
+// the page shows the message hash printed here. It writes nothing and posts nothing.
+func decisionPasskeyLinks(ctx context.Context, o decideOptions, credPath string, expiry int64, msg []byte, next string, stdout io.Writer, refuse func(string, ...any) int) int {
+	ro, err := store.OpenCredentialsReadOnly(credPath)
+	if err != nil {
+		return refuse("cannot read the credentials store %s: %v; nothing was decided", credPath, err)
+	}
+	active, err := ro.ActivePasskeys(ctx)
+	_ = ro.Close()
+	if err != nil {
+		return refuse("%v; nothing was decided", err)
+	}
+	if len(active) == 0 {
+		return refuse("no passkey approver is enrolled; see ailang-worldd approver list")
+	}
+	h := sha256.Sum256(msg)
+	hash16, request := hex.EncodeToString(h[:8]), "decide-"+o.verb+"-"+strings.TrimPrefix(o.ref, "sha256:")[:16]
+	fmt.Fprintf(stdout, "%s request %s with a passkey (expires %s, message hash %s): open the link for its origin; the page reads \"Approve request %s\" for either verb: check it shows message hash %s, and save the file it offers:\n",
+		strings.ToUpper(o.verb), o.ref, approverTime(expiry), hash16, request, hash16)
+	links := map[string]*ceremonyLink{}
+	var order []string
+	for _, e := range active { // ActivePasskeys is ordered by name
+		l := links[e.Passkey.Origin]
+		if l == nil {
+			l = &ceremonyLink{V: 1, Op: "get", RP: e.Passkey.RPID, Origin: e.Passkey.Origin, Request: request, Expires: expiry, Hash16: hash16,
+				Challenge: base64.RawURLEncoding.EncodeToString(h[:])}
+			links[e.Passkey.Origin] = l
+			order = append(order, e.Passkey.Origin)
+		}
+		l.Allow = append(l.Allow, base64.RawURLEncoding.EncodeToString(e.Passkey.CredentialID))
+	}
+	for _, origin := range order {
+		fmt.Fprintf(stdout, "  %s\n", links[origin].link())
+	}
+	fmt.Fprintf(stdout, "then decide it:\n  %s\n", next)
+	return exitOK
+}
+
 // postDecision POSTs body to /surface/decide/{ref} with the daemon's own loopback Origin (in O by
 // construction; design §2.8) and, for the typed arm, the operator proof; it prints the decision.
 func postDecision(c *client, origin, ref string, body surfacepage.DecideBody, proof string, stdout io.Writer, refuse func(string, ...any) int) int {
diff --git a/cmd/ailang-worldd/approve_decide_artifact_test.go b/cmd/ailang-worldd/approve_decide_artifact_test.go
new file mode 100644
index 0000000..f23e3f1
--- /dev/null
+++ b/cmd/ailang-worldd/approve_decide_artifact_test.go
@@ -0,0 +1,222 @@
+package main
+
+// Row 201 A4 sprint 3 M6e (design §2.10, §6 CD6-CD9): approve <request-ref> with an artifact made
+// elsewhere, against the fake daemon (approve_decide_rig_test.go): --write-message and --signature
+// --expiry (with the real ssh-keygen signing the written file, on a key made in a temp dir), and
+// --webauthn-link and --webauthn-assertion --expiry (an approvaltest passkey signing the link's
+// challenge, judged by the daemon's own law). D-WORLD-109: temp dirs, no real key, no agent.
+
+import (
+	"bytes"
+	"context"
+	"crypto/sha256"
+	"encoding/base64"
+	"encoding/json"
+	"os"
+	"os/exec"
+	"path/filepath"
+	"regexp"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// CD6, CD7: --write-message writes the message for the window's max to a NEW file (never over an
+// existing one) and prints the signing line in the decision namespace and the --signature step with
+// its expiry; --signature posts the file with the --expiry it was signed for, as sshsig; --signature
+// without --expiry, --expiry alone and --expiry with a mint id are refused before any daemon call.
+func TestApproveDecideWriteMessageAndSignature(t *testing.T) {
+	f := newFakeDecider(t)
+	env := approveEnv(nil, nil, atTime(1710))
+	msgPath := filepath.Join(t.TempDir(), "decide.msg")
+	code, stdout, stderr := f.keyRun(t, env, "--deny", "--write-message", msgPath)
+	got, _ := os.ReadFile(msgPath)
+	if code != exitOK || !bytes.Equal(got, f.decisionMessageOf(t, approval.VerbDeny, 1760)) || len(got) != 91 || len(f.posts) != 0 ||
+		!strings.Contains(stdout, "ssh-keygen -Y sign -n world-decision@sunholo -f <your key> < ") || !strings.Contains(stdout, f.ref+" --deny --signature ") ||
+		!strings.Contains(stdout, "--expiry 1760 --daemon "+f.addr()) || !strings.Contains(stdout, "deciding DENY") {
+		t.Fatalf("write-message: rc %d, %d POSTs\nstdout %s\nstderr %s", code, len(f.posts), stdout, stderr)
+	}
+	if code, _, stderr := f.keyRun(t, env, "--write-message", msgPath); code != exitUsage || !strings.Contains(stderr, "nothing was written") {
+		t.Fatalf("an existing file: rc %d %s", code, stderr)
+	}
+	sigPath := filepath.Join(t.TempDir(), "decide.sig")
+	if err := os.WriteFile(sigPath, []byte("SIGNED-ELSEWHERE"), 0o600); err != nil {
+		t.Fatal(err)
+	}
+	if code, _, stderr := f.keyRun(t, env, "--signature", sigPath, "--expiry", "1755"); code != exitOK || len(f.posts) != 1 ||
+		f.postBodies[0].SSHSIG == nil || *f.postBodies[0].SSHSIG != "SIGNED-ELSEWHERE" || f.postBodies[0].Expiry != 1755 || f.postBodies[0].Typed ||
+		f.postBodies[0].WebAuthn != nil || f.posts[0].Header.Get("X-World-Operator-Proof") != "" {
+		t.Fatalf("--signature: rc %d, %d POSTs\nstderr %s", code, len(f.posts), stderr)
+	}
+	missing := filepath.Join(t.TempDir(), "absent")
+	for _, flag := range []string{"--signature", "--webauthn-assertion"} {
+		if code, _, stderr := f.keyRun(t, env, flag, missing, "--expiry", "1755"); code != exitUsage || len(f.posts) != 1 || !strings.Contains(stderr, flag+": ") {
+			t.Fatalf("%s of a missing file: rc %d, %d POSTs\nstderr %s", flag, code, len(f.posts), stderr)
+		}
+	}
+	gets := f.gets
+	for _, tc := range []struct {
+		args []string
+		want string
+	}{
+		{[]string{f.ref, "--signature", sigPath}, "go with --expiry <t>"},
+		{[]string{f.ref, "--expiry", "1755"}, "go with --expiry <t>"},
+		{[]string{"mr-0123456789abcdef", "--expiry", "1755"}, "--expiry take a Human.Approve request ref"},
+	} {
+		var stdout, stderr bytes.Buffer
+		if code := runApprove(append(tc.args, "--db", f.db), &stdout, &stderr, env); code != exitUsage || f.gets != gets || !strings.Contains(stderr.String(), tc.want) { // no --daemon: it would trip CD3 itself
+			t.Errorf("%v: rc %d, GETs %d\nstderr %s", tc.args, code, f.gets, stderr.String())
+		}
+	}
+}
+
+// CD6 + CD7 with the real tool: the file --write-message wrote, signed by the real ssh-keygen with
+// the printed namespace on a key made in a temp dir (no agent, temp HOME), posted by --signature
+// with the printed expiry, verifies with the daemon's own check (VerifyDecisionSSHSIG) and is not
+// a mint signature. With WORLD_REQUIRE_SSH_KEYGEN=1 a missing ssh-keygen fails (CI).
+func TestApproveDecideWriteMessageSignedBySSHKeygen(t *testing.T) {
+	if _, err := exec.LookPath("ssh-keygen"); err != nil {
+		if os.Getenv("WORLD_REQUIRE_SSH_KEYGEN") == "1" {
+			t.Fatal("WORLD_REQUIRE_SSH_KEYGEN=1 but ssh-keygen was not found")
+		}
+		t.Skip("SKIP COUNTED: ssh-keygen not found; the decision --write-message end-to-end test did NOT run")
+	}
+	t.Setenv("HOME", t.TempDir())
+	t.Setenv("SSH_AUTH_SOCK", "")
+	dir := t.TempDir()
+	key := filepath.Join(dir, "k")
+	sshKeygen(t, "-q", "-t", "ed25519", "-N", "", "-C", "cd7", "-f", key)
+	pub, err := os.ReadFile(key + ".pub")
+	if err != nil {
+		t.Fatal(err)
+	}
+	k, err := approval.ParseAuthorizedKey(pub)
+	if err != nil {
+		t.Fatal(err)
+	}
+	f := newFakeDecider(t)
+	env := approveEnv(nil, nil, atTime(1710))
+	msgPath := filepath.Join(dir, "decide.msg")
+	code, stdout, stderr := f.keyRun(t, env, "--write-message", msgPath)
+	m := regexp.MustCompile(`ssh-keygen -Y sign -n (\S+) .* --expiry (\d+) `).FindStringSubmatch(strings.ReplaceAll(stdout, "\n", " "))
+	if code != exitOK || m == nil {
+		t.Fatalf("write-message: rc %d\nstdout %s\nstderr %s", code, stdout, stderr)
+	}
+	sshKeygen(t, "-Y", "sign", "-n", m[1], "-f", key, msgPath) // writes decide.msg.sig
+	if code, _, stderr := f.keyRun(t, env, "--signature", msgPath+".sig", "--expiry", m[2]); code != exitOK || len(f.posts) != 1 {
+		t.Fatalf("--signature: rc %d, %d POSTs\nstderr %s", code, len(f.posts), stderr)
+	}
+	armored := []byte(*f.postBodies[0].SSHSIG)
+	if err := approval.VerifyDecisionSSHSIG(armored, approval.VerbApprove, refDigest(f.subject.Packet), refDigest(f.ref), 1760, k); err != nil {
+		t.Fatalf("the real signature of the written file does not decide: %v", err)
+	}
+	if err := approval.VerifyApproval(armored, refDigest(f.ref), refDigest(f.subject.Packet), 1760, k); err == nil {
+		t.Fatal("the decision signature verifies as a mint")
+	}
+}
+
+var cdPasskey = approvaltest.NewPasskey(approval.AlgEdDSA, 73)
+
+const cdOrigin, cdRP = "https://studio.x.ts.net", "studio.x.ts.net"
+
+// enrolDecisionPasskey makes f's credentials store (replacing its empty file) and, with enrol,
+// enrols mark-mac's passkey in it.
+func (f *fakeDecider) enrolDecisionPasskey(t *testing.T, enrol bool) {
+	t.Helper()
+	if fi, err := os.Stat(f.creds); err == nil && fi.Size() == 0 {
+		if err := os.Remove(f.creds); err != nil {
+			t.Fatal(err)
+		}
+	}
+	c, err := store.OpenCredentials(f.creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = c.Close() }()
+	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
+	defer cancel()
+	k := cdPasskey.Key()
+	if !enrol {
+		return
+	}
+	if err := c.EnrolPasskey(ctx, store.ApproverRow{Name: "mark-mac", Format: approval.FormatWebAuthn, PublicKey: k.Canonical(), EnrolledAt: 1000},
+		store.PasskeyRow{Name: "mark-mac", RPID: cdRP, Origin: cdOrigin, CredentialID: cdPasskey.CredentialID(), AAGUID: make([]byte, 16), Alg: k.Alg, SignCount: 1}); err != nil {
+		t.Fatal(err)
+	}
+}
+
+// CD8, CD9: --webauthn-link prints the passkey page's get link for the active passkey, whose
+// challenge is the decision message's hash for the window's max (the verb in the link's request and
+// in capitals, the hash the page shows) and the --webauthn-assertion step, and posts nothing; the
+// assertion an authenticator makes over that challenge, posted by --webauthn-assertion with the
+// printed expiry, is the daemon's law's decision (DecideByAssertion) as webauthn:mark-mac. A file
+// that is not an assertion is refused before the POST; no passkey enrolled is refused.
+func TestApproveDecideWebAuthnLinkAndAssertion(t *testing.T) {
+	f := newFakeDecider(t)
+	env := approveEnv(nil, nil, atTime(1710))
+	f.enrolDecisionPasskey(t, false)
+	if code, _, stderr := f.keyRun(t, env, "--webauthn-link"); code != exitUsage || !strings.Contains(stderr, "no passkey approver is enrolled") {
+		t.Fatalf("no passkey: rc %d %s", code, stderr)
+	}
+	f.enrolDecisionPasskey(t, true)
+	code, stdout, stderr := f.keyRun(t, env, "--deny", "--webauthn-link")
+	m := regexp.MustCompile(`(?m)^  (https://\S+/surface/passkey#(\S+))$`).FindStringSubmatch(stdout)
+	if code != exitOK || m == nil || len(f.posts) != 0 || !strings.Contains(stdout, "DENY request "+f.ref) ||
+		!strings.Contains(stdout, f.ref+" --deny --webauthn-assertion <the saved file> --expiry 1760 --daemon "+f.addr()) {
+		t.Fatalf("link: rc %d, %d POSTs\nstdout %s\nstderr %s", code, len(f.posts), stdout, stderr)
+	}
+	raw, err := base64.RawURLEncoding.DecodeString(m[2])
+	if err != nil {
+		t.Fatal(err)
+	}
+	var l ceremonyLink
+	if err := json.Unmarshal(raw, &l); err != nil {
+		t.Fatal(err)
+	}
+	h := sha256.Sum256(f.decisionMessageOf(t, approval.VerbDeny, 1760))
+	ref16 := strings.TrimPrefix(f.ref, "sha256:")[:16]
+	if l.Op != "get" || l.Origin != cdOrigin || l.RP != cdRP || l.Challenge != base64.RawURLEncoding.EncodeToString(h[:]) || l.Expires != 1760 ||
+		l.Request != "decide-deny-"+ref16 || len(l.Allow) != 1 || l.Allow[0] != base64.RawURLEncoding.EncodeToString(cdPasskey.CredentialID()) ||
+		!strings.Contains(stdout, "message hash "+l.Hash16) || l.Hash16 != base64Hex(h[:8]) {
+		t.Fatalf("link %+v\nstdout %s", l, stdout)
+	}
+	asrPath := filepath.Join(t.TempDir(), "world-assertion-"+l.Request+".json")
+	if err := os.WriteFile(asrPath, append(cdPasskey.Assert(cdOrigin, cdRP, h, 2), '\n'), 0o600); err != nil {
+		t.Fatal(err)
+	}
+	if code, _, stderr := f.keyRun(t, env, "--deny", "--webauthn-assertion", asrPath, "--expiry", "1760"); code != exitOK || len(f.posts) != 1 {
+		t.Fatalf("assertion: rc %d, %d POSTs\nstderr %s", code, len(f.posts), stderr)
+	}
+	b := f.postBodies[0]
+	if b.WebAuthn == nil || b.SSHSIG != nil || b.Typed || b.Verb != "deny" || b.Expiry != 1760 || f.posts[0].Header.Get("X-World-Operator-Proof") != "" {
+		t.Fatalf("POST %+v", b)
+	}
+	k := cdPasskey.Key()
+	got, err := approval.DecideByAssertion(approval.DecisionAssertionInput{Verb: approval.VerbDeny, Packet: refDigest(f.subject.Packet), Request: refDigest(f.ref),
+		Expiry: b.Expiry, Now: 1710, Response: *b.WebAuthn,
+		Approver: approval.Approver{Name: "mark-mac", Format: approval.FormatWebAuthn, EnrolledAt: 1000,
+			Passkey: &approval.Passkey{RPID: cdRP, Origin: cdOrigin, CredentialID: cdPasskey.CredentialID(), Key: k, SignCount: 1}}})
+	if err != nil || got.DecidedBy() != "webauthn:mark-mac" {
+		t.Fatalf("the posted assertion does not decide: %v", err)
+	}
+	notAsr := filepath.Join(t.TempDir(), "x.json")
+	if err := os.WriteFile(notAsr, []byte(`{"id":"x"}`), 0o600); err != nil {
+		t.Fatal(err)
+	}
+	if code, _, stderr := f.keyRun(t, env, "--webauthn-assertion", notAsr, "--expiry", "1760"); code != exitUsage || len(f.posts) != 1 || !strings.Contains(stderr, "is not a passkey assertion") {
+		t.Fatalf("not an assertion: rc %d, %d POSTs\nstderr %s", code, len(f.posts), stderr)
+	}
+}
+
+func base64Hex(b []byte) string {
+	const digits = "0123456789abcdef"
+	out := make([]byte, 0, 2*len(b))
+	for _, c := range b {
+		out = append(out, digits[c>>4], digits[c&15])
+	}
+	return string(out)
+}
diff --git a/cmd/ailang-worldd/approve_decide_key_test.go b/cmd/ailang-worldd/approve_decide_key_test.go
index 02de42f..0a684a9 100644
--- a/cmd/ailang-worldd/approve_decide_key_test.go
+++ b/cmd/ailang-worldd/approve_decide_key_test.go
@@ -1,6 +1,6 @@
 package main
 
-// Row 201 A4 sprint 2 M5b (design §2.10, §6 CD5; CD6-CD9 are sprint 3's, plan S2-D1): approve
+// Row 201 A4 sprint 2 M5b (design §2.10, §6 CD5; CD6-CD9 are sprint 3's, approve_decide_artifact_test.go): approve
 // <request-ref> --sign-with against the fake daemon (approve_decide_rig_test.go): the 91-byte
 // decision message through the signer seam, and the real ssh-keygen on a key made in a temp dir
 // (D-WORLD-109: no real key, no agent).
@@ -65,11 +65,6 @@ func TestApproveDecideSignWith(t *testing.T) {
 	if got, want := reflect.ValueOf(realSessionEnv().signDecision).Pointer(), reflect.ValueOf(signDecisionWithSSHKeygen).Pointer(); got != want {
 		t.Fatal("realSessionEnv().signDecision is not signDecisionWithSSHKeygen")
 	}
-	for _, flag := range []string{"--signature", "--write-message"} { // not a request ref's flags in this sprint
-		if code, _, stderr := g.keyRun(t, env, flag, "/x"); code != exitUsage || g.gets != 1 || !strings.Contains(stderr, "the other key flags do not decide it") {
-			t.Errorf("%s: rc %d, GETs %d\nstderr %s", flag, code, g.gets, stderr)
-		}
-	}
 }
 
 // CD5 with the real tool: the production decision signer runs the real ssh-keygen on a key made
diff --git a/cmd/ailang-worldd/approve_decide_rig_test.go b/cmd/ailang-worldd/approve_decide_rig_test.go
index 4fbf220..0e29487 100644
--- a/cmd/ailang-worldd/approve_decide_rig_test.go
+++ b/cmd/ailang-worldd/approve_decide_rig_test.go
@@ -153,8 +153,8 @@ func TestApproveDecideNoTerminal(t *testing.T) {
 	}
 }
 
-// CD3: --deny and --daemon need a request ref (an mr- id, --list); --await, a passkey flag and a
-// key flag are refused with a request ref; two modes are F1. Each before any daemon call.
+// CD3: --deny and --daemon need a request ref (an mr- id, --list); --await, and an assertion
+// without --expiry (CD9), are refused with a request ref; two modes are F1. Each before any daemon call.
 func TestApproveDecideFlagRules(t *testing.T) {
 	f := newFakeDecider(t)
 	for _, tc := range []struct {
@@ -166,7 +166,7 @@ func TestApproveDecideFlagRules(t *testing.T) {
 		{[]string{"mr-0123456789abcdef", "--daemon", f.addr()}, "a mint request has no deny"},
 		{[]string{"--list", "--deny"}, "a mint request has no deny"},
 		{[]string{f.ref, "--await"}, "there is nothing to await"},
-		{[]string{f.ref, "--webauthn-link"}, "there is nothing to await"},
+		{[]string{f.ref, "--webauthn-assertion", "x"}, "go with --expiry <t>"},
 		{[]string{f.ref, "--signature", "x", "--sign-with", "y"}, "usage: ailang-worldd approve"},
 	} {
 		var stdout, stderr bytes.Buffer
````

### Patch M6f — the live-Chrome sibling

````diff M6f
diff --git a/cmd/ailang-worldd/approve_decide_chrome_test.go b/cmd/ailang-worldd/approve_decide_chrome_test.go
new file mode 100644
index 0000000..6e3392b
--- /dev/null
+++ b/cmd/ailang-worldd/approve_decide_chrome_test.go
@@ -0,0 +1,83 @@
+package main
+
+// Row 201 A4 sprint 3 M6f (design §2.10, §6 CD8-CD9; the real-tool sibling of DM7 and CD8-CD9, plan
+// §0 retro 11): approve <request-ref> --webauthn-link in a live Chrome, gated as XB1 (WORLD_CHROME,
+// WORLD_REQUIRE_PASSKEY_BROWSER=1; skipped and counted without them). Chrome's virtual authenticator
+// (a test key, D-WORLD-109, temp profile) creates a passkey through the A2 page, the decision link
+// the CLI printed signs the decision message's hash through the same page, unchanged, and the file
+// it saved, posted by --webauthn-assertion, is the daemon's law's decision: Chrome's bytes, not
+// approvaltest's.
+
+import (
+	"bytes"
+	"context"
+	"crypto/sha256"
+	"encoding/base64"
+	"path/filepath"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func TestApproveDecideWebAuthnLinkInChrome(t *testing.T) {
+	b := startPasskeyBrowser(t)
+	origin := passkeyPageOrigin(t)
+	b.do(t, nil, "tab", map[string]any{"authenticator": "uv"})
+	enrol := sha256.Sum256([]byte("decide-enrol"))
+	dl := t.TempDir()
+	var page pageCeremony
+	b.do(t, nil, "open", map[string]any{"url": ceremonyLink{V: 1, Op: "create", RP: "localhost", Origin: origin, Name: "mark-mac",
+		Challenge: base64.RawURLEncoding.EncodeToString(enrol[:])}.link()})
+	b.do(t, &page, "ceremony", map[string]any{"download": dl})
+	reg := saved(t, page, origin, dl, "world-registration-mark-mac.json")
+	a, err := approval.NewPasskeyApprover(approval.RegistrationInput{Name: "mark-mac", Origin: origin, Challenge: enrol, Response: reg, Now: 1000})
+	if err != nil {
+		t.Fatal(err)
+	}
+	f := newFakeDecider(t)
+	f.enrolDecisionPasskey(t, false)
+	c, err := store.OpenCredentials(f.creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
+	defer cancel()
+	p := a.Passkey
+	if err := c.EnrolPasskey(ctx, store.ApproverRow{Name: a.Name, Format: a.Format, PublicKey: p.Key.Canonical(), EnrolledAt: a.EnrolledAt},
+		store.PasskeyRow{Name: a.Name, RPID: p.RPID, Origin: p.Origin, CredentialID: p.CredentialID, AAGUID: p.AAGUID[:], Alg: p.Key.Alg,
+			BackupEligible: p.BackupEligible, SignCount: int64(p.SignCount)}); err != nil {
+		t.Fatal(err)
+	}
+	_ = c.Close()
+	env := approveEnv(nil, nil, atTime(1710))
+	code, stdout, stderr := f.keyRun(t, env, "--webauthn-link")
+	link, err := printedLink(stdout)
+	if code != exitOK || err != nil {
+		t.Fatalf("--webauthn-link: rc %d, %v\nstdout %s\nstderr %s", code, err, stdout, stderr)
+	}
+	var view pageView
+	b.do(t, &view, "open", map[string]any{"url": link})
+	ref16 := strings.TrimPrefix(f.ref, "sha256:")[:16]
+	h := sha256.Sum256(f.decisionMessageOf(t, approval.VerbApprove, 1760))
+	if view.State != "ready" || !strings.Contains(view.What, "decide-approve-"+ref16) || !strings.Contains(view.What, "message hash "+base64Hex(h[:8])) {
+		t.Fatalf("the page shows %+v", view)
+	}
+	dl2 := t.TempDir()
+	b.do(t, &page, "ceremony", map[string]any{"download": dl2})
+	asr := saved(t, page, origin, dl2, "world-assertion-decide-approve-"+ref16+".json")
+	path := filepath.Join(dl2, "world-assertion-decide-approve-"+ref16+".json")
+	if code, _, stderr := f.keyRun(t, env, "--webauthn-assertion", path, "--expiry", "1760"); code != exitOK || len(f.posts) != 1 || f.postBodies[0].WebAuthn == nil {
+		t.Fatalf("--webauthn-assertion: rc %d, %d POSTs\nstderr %s", code, len(f.posts), stderr)
+	}
+	got, err := approval.DecideByAssertion(approval.DecisionAssertionInput{Verb: approval.VerbApprove, Packet: refDigest(f.subject.Packet), Request: refDigest(f.ref),
+		Expiry: 1760, Now: 1710, Response: *f.postBodies[0].WebAuthn, Approver: a})
+	if !bytes.Equal(bytes.TrimSpace(asr), []byte(*f.postBodies[0].WebAuthn)) {
+		t.Fatalf("the POST did not carry the saved file:\n%s\n%s", asr, *f.postBodies[0].WebAuthn)
+	}
+	if err != nil || got.DecidedBy() != "webauthn:mark-mac" {
+		t.Fatalf("Chrome's decision assertion does not decide: %v", err)
+	}
+}
````

### Patch M6g — QUICKSTART §18

````diff M6g
diff --git a/cmd/ailang-worldd/quickstart_decide_test.go b/cmd/ailang-worldd/quickstart_decide_test.go
new file mode 100644
index 0000000..409d23f
--- /dev/null
+++ b/cmd/ailang-worldd/quickstart_decide_test.go
@@ -0,0 +1,53 @@
+package main
+
+// Row 201 A4 sprint 3 M6g (design §2.10, §10; S7): QUICKSTART §18 documents the decision mode
+// with the CLI's own flags, the decision namespace and the message length the code uses.
+
+import (
+	"fmt"
+	"os"
+	"regexp"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+func TestQuickstartSection18FlagsMatchTheCLI(t *testing.T) {
+	guide, err := os.ReadFile("../../docs/QUICKSTART.md")
+	if err != nil {
+		t.Fatal(err)
+	}
+	_, section, ok := strings.Cut(string(guide), "## 18. Decide an agent's approval request")
+	if !ok {
+		t.Fatal("QUICKSTART has no §18")
+	}
+	section = strings.ReplaceAll(section, "\\\n", " ")
+	seen := map[string]bool{}
+	lines := 0
+	for _, line := range strings.Split(section, "\n") {
+		if !strings.HasPrefix(line, "ailang-worldd approve sha256:") {
+			continue
+		}
+		lines++
+		for _, flag := range regexp.MustCompile(`(?:^|\s)(--[a-z-]+)`).FindAllStringSubmatch(line, -1) {
+			seen[flag[1]] = true
+			if !strings.Contains(approveHelp, flag[1]+" ") && !strings.Contains(approveHelp, flag[1]+"]") {
+				t.Errorf("§18 uses %s, which approve's help does not document: %s", flag[1], line)
+			}
+		}
+	}
+	for _, flag := range []string{"--deny", "--sign-with", "--write-message", "--signature", "--expiry", "--webauthn-link", "--webauthn-assertion", "--daemon", "--db"} {
+		if !seen[flag] {
+			t.Errorf("§18 shows no approve line with %s", flag)
+		}
+	}
+	if lines != 6 {
+		t.Errorf("§18 holds %d approve sha256: lines, want 6", lines)
+	}
+	for _, want := range []string{"-n " + approval.DecisionNamespace, fmt.Sprintf("%d-byte decision message", approval.DecisionMessageLen)} {
+		if !strings.Contains(section, want) {
+			t.Errorf("§18 does not say %q", want)
+		}
+	}
+}
diff --git a/docs/QUICKSTART.md b/docs/QUICKSTART.md
index 73429f1..99b5418 100644
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -1857,3 +1857,32 @@ ailang-worldd approve --list --decided --db /tmp/se-world/world.db
 ```
 
 Open the same card again and press **Approve with passkey**: the daemon refuses it as already used. With an SSH approver (§15) instead, **Save the message** on the card, sign it with the `ssh-keygen -Y sign -n world-approval@sunholo` line the card shows, paste the signature and press **Approve with this signature**. Every request and answer has a 10-second limit: a card that cannot be read says so and offers to read it again, and an approval whose answer never arrives is shown as *outcome unknown* with a status check, never as refused or minted. A busy credentials store is tried again at most three times, a second apart. Reached through `tailscale serve`, start `serve` with `--surface-origin https://<host>.ts.net`, or the card and its Approve are refused (421, 403) naming that flag.
+
+## 18. Decide an agent's approval request *(attended — pending first verbatim run)*
+
+When an agent's plan requests `Human.Approve` under a live grant (§11's **Inspect a decision packet**), the daemon commits a decision packet and the request stays `approvals chain: pending` until you decide it. `ailang-worldd approve <request-ref>` decides it: the ref is the packet page's approval request, `sha256:` and 64 hex. It reads the request and its one packet from the daemon serving `--db` (`--daemon <host:port>`, default the installed service's address), checks that both hash to their refs, shows every field and the verb in capitals, and has the daemon record the decision. A request is decided once, approve or deny; publish- and finish-shaped requests stay with `world-publish approve` and `episode finish`. This guide uses the §0 scratch store `/tmp/se-world/world.db` with its daemon serving at `127.0.0.1:7644`, a §15 SSH approver and a §16 passkey `demo-mac`; never a real store.
+
+```bash
+ailang-worldd approve sha256:<64 hex> --deny --daemon 127.0.0.1:7644 --db /tmp/se-world/world.db
+ailang-worldd approve sha256:<another ref> --sign-with ~/.ssh/<your enrolled key> --daemon 127.0.0.1:7644 --db /tmp/se-world/world.db
+```
+
+The first asks `y` at the terminal and records `deny by operator:terminal`. The second opens no terminal: `ssh-keygen -Y sign -n world-decision@sunholo` signs the 91-byte decision message (`world-decision/v1`, the verb, the packet and request digests, an expiry at most 10 minutes ahead and never past the packet's deadline), and the daemon records `approve by ssh:<name>` after checking the signature against the enrolled, unrevoked key. A decision message is not a mint message: a §15 mint signature never decides, and a decision signature never mints.
+
+To sign elsewhere, write the message to a new file and post the signature with the expiry the CLI printed:
+
+```bash
+ailang-worldd approve sha256:<ref> --write-message /tmp/se-world/decide.msg --daemon 127.0.0.1:7644 --db /tmp/se-world/world.db
+ssh-keygen -Y sign -n world-decision@sunholo -f ~/.ssh/<your enrolled key> /tmp/se-world/decide.msg
+ailang-worldd approve sha256:<ref> --signature /tmp/se-world/decide.msg.sig --expiry <the printed expiry> --daemon 127.0.0.1:7644 --db /tmp/se-world/world.db
+```
+
+With a passkey, `--webauthn-link` prints one link per passkey origin to the §16 page. That page reads *Approve request decide-approve-<16 hex>* (or *decide-deny-…*) whichever verb you chose: check that it shows the message hash the CLI printed, press **Approve with passkey**, pass your device's check and **Save**. Then post the saved file with the printed expiry:
+
+```bash
+ailang-worldd approve sha256:<ref> --webauthn-link --daemon 127.0.0.1:7644 --db /tmp/se-world/world.db
+ailang-worldd approve sha256:<ref> --webauthn-assertion ~/Downloads/world-assertion-decide-approve-<16 hex>.json --expiry <the printed expiry> \
+  --daemon 127.0.0.1:7644 --db /tmp/se-world/world.db
+```
+
+The daemon records `approve by webauthn:demo-mac` after checking the assertion against the enrolled passkey, the decision message's hash, the origin and a signature counter that must move forward. Posting the same file again is refused as already decided. If the world write fails after the daemon checked a key decision, the same file (or signature) posted again lands that one decision without moving the counter twice; a different assertion carrying the same counter is refused. The packet page then reads `approvals chain: approve by …`.
````
