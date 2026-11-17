# Sprint plan — w-world-address-bar-verbs-v3, sprint S4 "the POST" (iteration 403, row 243)

> **THE SPRINT — row 243 phase V3, sprint S4 (design §0, §1, §2.2, §2.3, §2.4 (C15's store contract), §2.5, §2.6, §2.7
> (the seam S5 reads), §2.9 rows S4 reaches, §3, §4, §5, §6 "Sprint 4", §7, §8, §11, §12, §14, §16), D-WORLD-118 = A:**
> **M10** — `POST /surface/confirm/{id}` (`host/daemon/surface_confirm.go`, NEW): the guards C0–C6 (A4's R0, A3a's R1,
> R3–R6; R2 is `approvalRoute`), then under `humanMu` and within `commitBudget` C7–C12 on a read-only handle (row,
> `approval.Admit`, the digest, `DecodeVerb` + `Validate` + `broker.FinishConfirmRefusal`, S1b's judges, `Binds`) —
> every refusal writes nothing — then C13 a non-waiting writer, **C14** `broker.RecordOperatorApproval(…, r.Scope,
> "operator:episode-finish", signed.ApprovedBy(), now)` (unchanged), **C15** `ConsumeVerbRequest(…, decision)` in one
> credentials transaction, the writer held across C14 and C15 and closed after; 200 = **`surfacepage.Confirmed`**
> (NEW). **One route registered**, `POST /surface/confirm/{id}` behind `approvalRoute`, in the sub-milestone that makes
> it complete (M10d). **M11** — `confirm.js`'s submit half: **Approve with your passkey** (`challenge` = the message
> hash, `allowCredentials` from the origin's group, UV required, A2's 240 s abort) and the **SSHSIG paste**, POSTing
> `ParseApproveBody`'s grammar; the answer or the refusal is shown, then the card is read again. Plus the sprint-3
> judge's note 1 (two exact-boundary pins, test-only) and S3 audit note 3's daemon half (a supersession between C14
> and C15). **No CLI change, no schema change, no bar, no *Needs you* rendering, no `.ail`, `surfaceCSP` unchanged.**
> - **Measured: 1,097 changed lines** — `git diff -M --numstat 430668f8 HEAD -- cmd host docs` (insertions +
>   deletions, renames on; no rename occurs; **the same 1,097 with no path filter**: nothing lands outside `cmd host
>   docs`), after a commit of each of the five scratch stages: **311 production** (`surface_confirm.go` 215, `confirm.js`
>   70, `confirm_model.go` 17, `index.html` 6, `daemon.go` 3), **786 test**, **0 docs**. **Cap 1,400: margin 303.**
> - **No cut (S4-D1).** M10 was compiled first and measured **916** (design ≈ 765 raw: 1.20×); with M10 at
>   916, M11 (design ≈ 180 raw) fits, so it was compiled: **181**. M10 + M11 = **1,097 ≤ 1,400** → both ship;
>   the design's cut (M11 → S6) was measured and **not taken**.
> - **Five sub-milestones, each green at its commit** (§1.2 C14): **M10a** the handler (unregistered) + `Confirmed` +
>   the writer-site pin + CX0–CX6 (433) → **M10b** CX7–CX13 + judge-402 note 1's pins (605) → **M10c** CX14–CX21, audit
>   note 3, C15's mapping (917) → **M10d** the registration, the route-table edit, the tests through the mux (916) →
>   **M11** the submit half, CP8–CP10 (1,097).
> - **Every line ships as a compiled patch** (Appendix P), extracted from this file and round-tripped onto a fresh
>   `git archive 430668f8` export, byte-identical to the stage tree (§1.2 C16). The executor applies, never retypes.
> - **Pre-existing test files changed: 3; deleted: 0** (§5, §9 criterion 12): `surface_approve_post_guard_test.go`
>   (the fourth writer site, AM8), `surface_confirm_read_test.go` (the route table gains the POST),
>   `confirm_script_test.go` (the stand-in DOM gains events; three new harness sections).

---

## 0. Binding constraints (each is restated where it governs)

- **D-WORLD-118 = A** (verified by the controller at iteration 400, `verdict=attended`, flip `72007b37`;
  `mission_decisions.sh --open` at iteration 403 → no OPEN row). S4 is routable; **F6, `DraftOnly` and every "B +≈ 15"
  branch are not built** (CX17's passkey row confirms a **ready** pull request).
- **Exactly one route added**: `POST /surface/confirm/{id}`, behind `approvalRoute` (`TestConfirmRoutesRegistered`: the
  three GETs and this POST, nothing else under `/surface/confirm`; PUT/DELETE/OPTIONS on the card's path and POST on the
  script and data paths still 405 with no CORS header). **No CLI change** (no `cmd/` file changes at all). **No schema
  change** (no `host/store` production file changes). **`surfaceCSP` unchanged** (`TestDecidePageAssets` pins the
  literal; the card's page and script are served under it, CR3). **No `.ail`.**
- **D-WORLD-109**: every store under `t.TempDir()`; keys are `approvaltest`'s (`apSigner` seed 73, `apPasskey` EdDSA
  seed 71); no real store, key, `~/.ssh`, `/dev/tty`, mint or publish; no listener (`Handler().ServeHTTP` /
  `approvalRoute(…)` directly); the node harness stubs `navigator.credentials` and `fetch`.
- Coding standards: the law stays pure (`approval.Admit`, `broker.FinishConfirmRefusal`, the S1 judges are reused, not
  changed); the POST is the host boundary's effect (design §3 S2); "why not a package" is design §3 S3.
- Not touched: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/`, `cmd/world-publish/*`, every
  `host/store`, `host/broker`, `host/approval`, `host/authority` production byte, `surface_page.go`, `approvalRoute`,
  the `/v1` table, `approve.js`, `decide.js`, `/surface/passkey`, `app.js`.

---

## 1. Pre-flight

### 1.1 Baseline at `430668f8` (measured by the planner, 2026-11-17)

`tools/gates.sh <base-tree> base-gates base` (the §6 script in `base` mode — its first draft, whose `race` leg chained with `&&` — on a `git archive 430668f8` export with a `git init`; `GATES_SKIP_VERIFY_GO=1`), run beside the planner's mutant run:

| leg | rc | time | note |
|---|---|---|---|
| `status` | 0 | 0 s |  |
| `verify_ail` | 0 | 163 s | `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` — the final must match |
| `gofmt` | 0 | 1 s |  |
| `build` | 0 | 23 s |  |
| `vet` | 0 | 7 s |  |
| `fence` | 0 | 40 s | `go test -run '^$' ./...` is the test-file compile fence |
| `pkgs` | 1 | 555 s | `host/store` ok 39 s, `host/surfacepage` ok 10 s; **`host/daemon` red (552 s) on `TestAGUISlowReaderCutIsResumable`** ("unexpected EOF": the AG-UI slow-reader cut, timing under load ≈ 20–30 from the planner's mutant run) — at the base, before any S4 line; S4 touches no AG-UI code |
| `race` | 1 | 158 s | **`host/store` red on `TestConvertedStoreMethodsM5b/cutoff/AppendOutcome`** ("context deadline exceeded, want uncertain": S3 §1.1's base red, a world-store cut-off under load); `host/surfacepage` ok; the draft leg chained the daemon subset after the store with `&&`, so it did not run — **fixed in §6's script** (each part runs, rc OR-ed); that subset run alone at the base afterwards: rc 0, 65 `--- PASS`, 13 s |
| `node` | 0 | 7 s | the nine named harnesses `--- PASS`, `WORLD_EXEC_NODE` set |
| `inv` | 0 | 82 s | **60 named `--- PASS`**, 0 SKIP, 0 FAIL |
| `chrome` | 0 | 131 s | 24 `--- PASS` (V2b's fourteen, XC1, the two `host/daemon` drills) |
| `SURFACE_POSTS` | — | — | **2** (`POST /surface/approve/{id}`, `POST /surface/decide/{ref}`) |
| **`GATES_DONE 1169s FAILS: pkgs race`** | | | both reds are load timing at the base (the directive's and S3's lists); the final tree is compared leg by leg with these (C17) |

### 1.2 Codebase claims checked (at `430668f8` unless "work" — the planner's stage tree, tag `M11`)

| # | claim | command | observed |
|---|---|---|---|
| C1 | **Input 1 — the design's §2.6 cites, re-measured: none of the code moved; three files shifted lines** | `git diff --stat 994ba123 430668f8 -- host/daemon/surface_decide.go host/daemon/surface_decide_key.go host/daemon/surface_approve_guard.go host/broker/approve.go host/daemon/episode_finish.go host/daemon/surface_decide_race_test.go host/daemon/surface_decide_key_retry_test.go host/daemon/surface_approve_post_guard_test.go host/daemon/route_table_test.go host/store/writer_lock.go host/store/credentials.go host/surfacepage/approve_model.go host/daemon/daemon.go`; control `git diff --stat 994ba123 430668f8 -- cmd host \| tail -1` | only `daemon.go \| 7` (S3's three GETs + comment), `credentials.go \| 12` (S3's v9), `approve_model.go \| 17` (S3's `Confirmations`); control `53 files changed, 3955 insertions(+), 48 deletions(-)` (S1, S2, S1b, S3) |
| C2 | the cites at HEAD, line by line | `grep -n`/`sed -n` | **moved**: `ParseApproveBody` `approve_model.go:197 → :214` (`MaxApproveBodyBytes = 20 << 10` `:27`, unmoved); the `/surface` approval routes `daemon.go:1139-1148 → :1141-1153` (S3's three confirm GETs at `:1150-1152`). **Unmoved**: `surface_decide.go:159` R0, `:186-188` `humanMu.Lock`/`defer Unlock`/`context.WithTimeout(…, d.commitBudget)`, `:211` `defer func() { _ = evidence.Close() }()`, `:216` `broker.DecideApproval`, `:222` `LandKeyDecision`; `surface_decide_key.go:66` `store.OpenCredentials(…) // non-waiting`, `:90` `return wr, signed.DecidedBy(), true`; `surface_approve_guard.go:53` `allowedApprovalOrigin`, `:57` `sameOriginFetch`, `:61` `jsonMediaType`; `daemon.go:150` `commitBudget = 3 * time.Second`, `:440` `humanMu sync.Mutex`, `:464` `approveNow func() int64`; `approve.go:722` `RecordOperatorApproval(ctx, s, scope, requester, decidedBy, now) (request, decision hashref.HashRef, error)`, `:219-226` W59 (`decided && bytes.Equal(existing.Payload, payload)` → the existing ref; `decided` → `ErrApprovalAlreadyDecided`); `credentials.go:210` `OpenCredentials`, its lock `writer_lock_unix.go:41` `syscall.Flock(…, LOCK_EX\|LOCK_NB)` (non-waiting), busy = `IsWriterAlreadyActive` (`writer_lock.go:61`); `surface_page.go:19` `surfaceCSP`; `route_table_test.go:15` `TestFrozenV1RouteTableMatchesMux` (its regexp is `/v1/` only, `:20`); `surface_approve_post_guard_test.go:132` `TestDaemonCredentialWriterSites` (want `ensureCredentialsStore`, `handleSurfaceApprove`, `keyArm`); `daemon_test.go:899` `TestDaemonDependencyAllowlist` |
| C3 | **Input 2 — the closest analogue, A4's `POST /surface/decide/{ref}`, and what was mirrored** | `sed -n 151,228p host/daemon/surface_decide.go`; `sed -n 34,91p host/daemon/surface_decide_key.go` | **Mirrored**: the guard order and texts' form (R0 first, then A3a's R3–R6 verbatim helpers), `humanMu.Lock(); defer Unlock()` then `context.WithTimeout(r.Context(), d.commitBudget)` then `d.approveNow()`; the judge on a fresh read-only handle with `JudgeStepStore` → the store-error answer and every other judge error → 422 `NotApproved`; a non-waiting `store.OpenCredentials` opened before the first write and `defer`-closed after the last (A4: across `DecideApproval` and `LandKeyDecision`; S4: across C14 and C15); `isApproverRace`/`isCounterRefusal` reused for C15's 422; the tests' mechanisms (`decideRig{a}.failWorldWrites()` = RD13 (a)'s objects trigger, `decideRig{a}.decisions()`, the H-LOCK eight-concurrent-POSTs shape, the operands test of `surface_decide_race_test.go`). **Differs** (design §2.6 [r1]): world first, then **one** credentials transaction (A4 records its evidence first and lands it after); W15 runs in full in the law (no deferred-W15 judge); the writer is opened in the handler (not a helper), so the AST pin names `handleSurfaceConfirm` |
| C4 | what S4 consumes from S1–S3 is at HEAD | `grep -n` | `store.VerbRequestIDPattern` (`verb_requests.go:58`), `VerbRequest` (`:168`), `ConsumeVerbRequest(ctx, id, digest, use, approvedBy, artifact, decisionRef, now)` (`:239`) and its sentinels (`:60-69`), `VerbOutcomeSuperseded` (`:53`), `OlderCredentialsStoreError`; `approval.Admit` (`admit.go:44`), `DecodeVerb` (`verb.go:133`), `Validate` (`:69`); `broker.FinishConfirmRefusal`/`FinishConfirmFacts` (`finish_confirm.go`); `authority.JudgeVerbSSH`/`JudgeVerbAssertion`/`VerbJudgeInput` (`judge_verb.go:63`, `:82`, `:18`); `SignedVerb.Binds`/`ApprovedBy`/`Use`/`PasskeyUse` (`verb_signed.go:149-187`); `confirmIDReason` (`surface_confirm_read.go:19`); `writeApprovalRefusal`, `writeApprovalJSON`, `approvalInternal`, `isOlderStore`, `approvalRoute` (`surface_approve.go:32-89`); `isApproverRace`, `isCounterRefusal` (`surface_decide_key.go:23`, `:30`) |
| C5 | nothing named confirmation POST exists | `grep -rn 'func Confirmed\|handleSurfaceConfirm(\|POST /surface/confirm' host --include=*.go \| wc -l`; control `grep -rn 'func Recorded' host/surfacepage \| wc -l` | `0`; control `1` |
| C6 | **Input 3 — fixtures from the real tools**, each name existing at HEAD | `grep -n '^func\|^var' host/daemon/*_test.go` | requests: S3's `putVerb`/`putVerbBytes`/`verbFixture`/`canonVerb`/`readVerb`/`setVerbColumn` (`surface_confirm_read_test.go:33-66`, `:273`); keys: `apSigner` (`surface_approve_sshsig_test.go:20`), `enrolSSH` (`:23`), `sshBody` (`:43`), `apPasskey`/`approveBody` (`surface_approve_post_test.go:23`, `:48`), `enrolPasskey` (`surface_approve_rig_test.go:93`); messages: `approval.VerbMessage`/`VerbMessageHash` (S1); CX17's expected body: `json.Marshal(surfacepage.Confirmed(id, approvedBy, row.DecisionRef))`; C10's rows: `verbFixture` edited, then `approval.CanonicalVerb` (`broker.FormatFinishApprovalScope`/`FinishPromptText` for the F2 row). **NEW names** (all test-only but `Confirmed`/`ConfirmedJSON`, `handleSurfaceConfirm`, `confirmAdmission`, `confirmStoreError`, `consumeRefusal`, `confirmRequester`, `confirmIncomplete`, the three seams): `newConfirmRig`, `verbMessageOf`, `confirmSSH`, `confirmPasskey`, `confirmAt`, `confirm`, `confirmState`, `credState`, `worldDecision` |
| C7 | **Input 4 — tests that enumerate routes, writer sites, assets, CSP, symbols, packages** | `grep -ln 'func TestConfirmRoutesRegistered\|…'` (the list in §6's `inv`) | **change**: `TestConfirmRoutesRegistered` (its `want` gains `POST /surface/confirm/{id}`; its 405 loop skips POST on the card's path — PUT/DELETE/OPTIONS still 405, POST on `/surface/confirm.js` and `/surface/confirmations/{id}` still 405); `TestDaemonCredentialWriterSites` (+ `handleSurfaceConfirm`, AM8). **Unchanged and green** (§6 `inv`): `TestFrozenV1RouteTableMatchesMux` (`/v1` only), `TestDaemonDependencyAllowlist` (no new import: `surface_confirm.go` imports only packages the daemon already imports), `TestConfirmPageAssetsServed` and `TestConfirmPageAssets` (both compare the served/embedded bytes with the files, which change together), `TestDecidePageAssets` (pins `surfaceCSP`), `TestSurfaceAppRoutes`, `TestProductionContextRoots`/`TestProductionGoSurface` (no new store method, no new root: the POST's ctx is `commitBudget`'s, A4's) |
| C8 | **Input 5 — the failure instruments** | `grep -n 'func (r decideRig) failWorldWrites\|func (r decideRig) decisions' host/daemon/*_test.go`; `grep -rn 'TRIGGER.*verb_requests' --include=*_test.go host cmd \| wc -l`, control `grep -rn 'CREATE TRIGGER' --include=*_test.go host/daemon \| wc -l` | A4's RD13 (a) instrument exists and is reused as is: `failWorldWrites` (`surface_decide_passkey_test.go:69`, `BEFORE INSERT ON objects WHEN NEW.semantic_id = ApprovalDecisionV1`) and `decisions` (`surface_decide_rig_test.go:88`), on `decideRig{a}` (`decideRig` is `struct{ approvalRig }`); no `verb_requests` trigger exists (`0`; control `2`) → CX21 writes the minimal one inline (`BEFORE UPDATE ON verb_requests WHEN NEW.outcome = 'approved'`). A credentials-store trigger does not trip the store's table-set check (it lists `type = 'table'` only, `credentials.go:421`) — CX21 green |
| C9 | the typed approve phase takes the same `humanMu` (CX19's other side) | `sed -n 156,176p host/daemon/episode_finish.go` | `:158 d.humanMu.Lock()`, `:159 defer d.humanMu.Unlock()`, `:164 context.WithTimeout(r.Context(), d.commitBudget)`, `:173 broker.RecordOperatorApproval(ctx, d.store, input.Scope, "operator:episode-finish", "operator:terminal", now)` |
| C10 | the page's harness reaches the submit half | `grep -n 'addEventListener\|navigator' host/surfacepage/confirm_script_test.go` at HEAD; work | HEAD: none (S2's stand-in DOM has no events, its context no `navigator`); work: `El.addEventListener`/`click()` (a disabled or hidden button does not fire) and `open(…, navigator)`, three sections `passkey`, `answers`, `sshsig` |
| C11 | the new code adds no context root, no deadline-free read, no non-test Go file outside `host/`, no import | `inv` (§6) on the final tree | `TestProductionContextRoots`, `TestNoNewDeadlineFreeStoreReads`, `TestProductionGoSurface`, `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist` `--- PASS`, unedited |
| C12 | the judge-402 note 1 premise: consume is `now >= ExpiresAt`; list and count are `expires_at > now` — consistent, unpinned at the second | `grep -n 'expires_at > ?\|now >= row.ExpiresAt' host/store/verb_requests.go`; S3's `TestPutVerbRequestGrammarCapPrune` and CR2 rows | `:143` count `expires_at > ?`, `:197` list `expires_at > ?`, `:268` consume `now >= row.ExpiresAt` (and the prune `:139` `expires_at <= now - 86400`): **consistent** (a request at `expires_at == now` is expired in all three), no production change; S3's tests used rows expiring 600 s later only → the two operators survived. S4 pins both (`TestVerbRequestExpiryBoundary`) and C8's at the daemon (CX8 `expired at expires_at`) |
| C13 | the new tests, work | `go test -count=1 -v` by name (`accept.sh`) | all `--- PASS`, none SKIP |
| C14 | **every sub-milestone is green when committed** | `tools/accept_all.sh M10a:fast … M11:chrome M10a:daemon M10d:daemon` (an APFS clone of `work` at each tag running §3's `accept.sh`) | on the final stages (`logs/accept_all.txt`, second pass; M10a's stage is byte-identical to the first pass's, tree hash compared): `M10a fast rc=0 ACCEPT_OK` (72 s), `M10a daemon rc=0` (469 s), `M10b fast rc=0` (61 s), `M10c fast rc=0` (44 s), `M10d fast rc=0` (26 s), `M10d daemon rc=0` (**645 s**: run beside the planner's mutant run (5 workers) and the base gates, load ≈ 20–30; `host/daemon` whole is ≈ 260 s alone, §1.1 `pkgs`), `M11 fast rc=0` (46 s), `M11 chrome rc=0` (57 s, XC1 in live Chrome on the new page bytes). A first pass on earlier stages (before four test-only hardenings the mutant enumeration asked for: C8's writes-nothing on the superseded/consumed rows, C11's store class, C7 unreadable, the CP10 in-flight check) was green too |
| C15 | the snapshot block (§7) on every sub-milestone | `tools/snap_rehearse.sh` (a local clone at `430668f8`, the patches in order, §7's block verbatim) | `M10a 5/0/5`, `M10b 6/0/6`, `M10c 7/0/7`, `M10d 9/0/9`, `M11 12/0/12` (`STATUS` lines / `DELETED` lines / MANIFEST lines), each `MANIFEST_OK`, dirs `host` only (no `cmd`, no `docs`: the `find` guard holds) — `logs/snap.txt` |
| C16 | **the round trip** | `tools/roundtrip.sh <this plan>` | `CMP_EQUAL` for all five (each block extracted from this file with §2's awk equals the generated patch byte for byte); each applied with `/usr/bin/patch -p1 -E` after a `--dry-run` to a fresh `git archive 430668f8` export, no `.orig`/`.rej`, committed and measured: `NUMSTAT M10a 433`, `M10b 605`, `M10c 917`, `M10d 916`, `M11 1097` (each `unfiltered=` the same); then `diff -r -x .git` against the stage tree (`git archive M11`) → **`ROUNDTRIP_IDENTICAL_TO_STAGE_TREE`** (`logs/roundtrip_last.txt`, re-run on this final plan file) |
| C17 | **the gates on the final tree** | `tools/gates.sh <final-tree> final-gates final <sha>` | `tools/gates.sh <final-tree> final-gates final 574ab4da` (`final-tree` = a clone of the repo at `430668f8`, the five patches applied with `patch -p1 -E` and committed one per sub-milestone; `diff -r` equal to the round-trip tree; `GATES_SKIP_VERIFY_GO=1`; load ≈ 10–20): **`TOUCHED_PACKAGES: host/daemon host/store host/surfacepage`** (= the plan's list); `status` 0; **`verify_ail` rc 0** (`✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` — the base's line, byte for byte); `gofmt`, `build`, `vet`, `fence` rc 0; **`pkgs` rc 0** (303 s: `host/store`, `host/surfacepage`, `host/daemon` whole — the base's AG-UI load red did not recur); **`race` rc 0** (97 s: `host/store` 47 s, `host/surfacepage` 8 s, the `host/daemon` subset 31 s with **73** `--- PASS`; the same subset at the base, run alone after the fix to the draft's `&&`: rc 0, **65** `--- PASS`, 13 s — the eight new tests); **`race3` rc 0** (19 s: the eight new daemon tests ×3, 13.7 s; the boundary pin ×3); **`node` rc 0** (12 named `--- PASS`: S3's nine + CP8–CP10); **`inv` rc 0: 70 named `--- PASS`** (the base's 60 + the eight CX tests, `TestVerbRequestExpiryBoundary`, `TestConfirmedWire`), 0 SKIP, 0 FAIL; **`SURFACE_POSTS: 3`** (`approve`, `decide`, `confirm`; base 2); **`chrome` rc 0** (24 `--- PASS`: V2b's fourteen, XC1 on the new page bytes, the two drills); **`GATES_DONE 862s FAILS: none`**. `verify_go` is the controller's |
| C18 | node and the pinned binary | `/opt/homebrew/bin/node --version`; `$HOME/.pinned-ailang/ailang --version` | `v26.13.0`; `AILANG v0.41.0` |

### 1.3 Planner decisions where the design or the directive leaves a choice

- **S4-D1 — no cut; M10 and M11 ship, by measurement.** M10 compiled first: **916** changed lines (incl. the
  judge-402 pins, 32, and `TestConfirmedWire`, 18) against the design's ≈ 765 raw. With 484 left under the cap,
  M11 was compiled: **181** (design ≈ 180 raw). Total **1,097 ≤ 1,400**: the design's cut "S4 → M11 moves to S6"
  was measured and not taken. Ratio 1.16× (S1 1.74×, S3 1.42×): M10's tests are table rows over shared helpers.
- **S4-D2 — five sub-milestones; the route is registered last (M10d).** M10a–M10c ship the complete handler
  **unregistered** and test it through `d.approvalRoute(d.handleSurfaceConfirm)` with the request's path value set (the
  helper `confirmAt`; R2's 421 is `approvalRoute`'s, so CX's R2 row is real at every stage). Unreachable code is safe
  at commit; no commit registers a POST whose writes are unproven. **M10d** registers the POST, extends
  `TestConfirmRoutesRegistered`, and switches `confirmAt` to `a.serve(…)` — from then every CX row goes through the mux,
  the session middleware and `approvalRoute`. **M11** (the page) follows the route it POSTs to.
- **S4-D3 — CX ids: the design's, with no CX16.** Design §6 numbers CX0–CX15 (one per step), then CX17–CX21; there is
  no CX16 and the plan adds none (S1 Appendix C note 3). Mapping: CX0–CX6 `TestConfirmPOSTGuards`; CX7–CX13
  `TestConfirmPOSTRefusesBeforeAnyWrite`; CX14 (a failing world write) is CX20's first half; CX15 (C15's refusals: a
  concurrent consume lost, approver revoked, counter) `TestConfirmPOSTRefusedAtC15` + `TestConsumeRefusalOperands`;
  CX17 `TestConfirmPOSTConfirms`; CX18 the `TestDaemonCredentialWriterSites` edit + A4's `TestDecideRefusesReservedScopesByKey`
  unedited and green; CX19 `TestConfirmPOSTHoldsHumanMu`; CX20 `TestConfirmPOSTRetriesAFailedWorldWrite`; CX21
  `TestConfirmPOSTRetriesAFailedConsume`; audit note 3 is a row of `TestConfirmPOSTRefusedAtC15` (`superseded`).
- **S4-D4 — R2 runs first.** The design's table lists C0, C1, R2; `approvalRoute` wraps the handler, so R2's 421 is
  answered before C0, as for A3a (R1 listed before R2) and A4. CX's R2 row is a foreign Host with a valid body.
- **S4-D5 — C7 answers below schema 9 (503 `Restart`)**, before C13: the read-only `VerbRequest` refuses a v8 store
  with `OlderCredentialsStoreError{Need: 9}` first, so CX13's "v8" row is C7's (§1.4 F2). The Restart text names verb
  requests (`confirmStoreError`), not A3a's "surface approvals" (`approvalStoreError`), which is why the POST has its
  own store-error answer.
- **S4-D6 — the 500 of C14 and of C15's store error carries the design's fixed text** "the confirmation did not
  complete; approve again" (design §2.6 window table), class `Internal`; the store's error goes to `d.errLog` only.
  `internalErrorMessage`'s rule (`handlers.go:136`) is about server detail on the wire; this text carries none. A
  read-only open that fails before C7 keeps `approvalInternal`'s sanitized 500.
- **S4-D7 — three test seams in production (package variables)**: `judgeVerbSSH`/`judgeVerbAssertion` (C11; CX12 needs
  a *verified* confirmation of other bytes for `Binds` to refuse — the real judge, given the row's own input, can never
  produce one; and CX11's store-class row) and `consumeVerbRequest` (C15; a supersession, a won consume, a revocation
  and a counter move **between** C14 and C15 need the writer the handler holds — a second `OpenCredentials` would be
  refused by the flock). Each seam's default is the real function; tests restore it in `t.Cleanup`. The daemon's tests
  are not `t.Parallel`.
- **S4-D8 — CX19 by the lock's identity, deterministically.** The design asks that "a confirmation and a typed finish
  approve phase on one daemon serialize". Both take `d.humanMu` (C9). CX19 proves the confirmation holds it at three
  points — its clock read (after `Lock`), C11 and C15 — with `d.humanMu.TryLock()` probes in the `approveNow` hook and
  the two seams (each must fail while the handler holds it: no sleep, no timing), then runs A4's H-LOCK shape (eight
  concurrent POSTs of one body: one 200, seven 409, never a busy writer). A two-route race of the typed phase would
  need the operator-proof file and a timed interleaving; the identity proof is stronger and bounded (§1.5: 0.17 s under
  `-race`).
- **S4-D9 — C14 and C15's window table, with audit note 3's row** (the design has none for it):

  | crash or failure | state left | the user sees | recovery | test |
  |---|---|---|---|---|
  | before C14 (incl. C13 busy, every C0–C12 refusal) | nothing written (credentials and world compared) | its 4xx/503 | approve again; the same body → 200 | CX0–CX13 (`confirmState` before = after), CX20 |
  | C14 fails | the world's approval request object and head (no decision); row pending, counter unmoved | 500 "the confirmation did not complete; approve again" | the identical body → 200, 1 decision | CX20 |
  | after C14, C15 fails | one decision no row names (inert, W62); row pending; key unused | 500, same text | the same body in the same second → 200, still 1 decision (W59); a second later → 200, 2 decisions, the row names the second | CX21 |
  | **after C14, the request superseded by a newer finish of its episode before C15 (audit note 3)** | one decision no row names (inert, W62); **the row `superseded`**; key unused | **409 "request <id> was superseded by a newer finish of <episode>"** | none on this card; the newer request's card | `TestConfirmPOSTRefusedAtC15` `superseded` |
  | after C14, C15 refuses: approver revoked / counter moved / a concurrent consume won | one orphan decision; row pending (won consume: approved by the winner) | 422 / 422 / 409 | a fresh gesture (revoked: none) | `TestConfirmPOSTRefusedAtC15` |
  | after C15 | done; the 200 may be lost | the card re-reads `approved` | none | CX17 (replay → 409) |

- **S4-D10 — judge-402 note 1 by test only.** C12 measured no inconsistency; the pins are `TestVerbRequestExpiryBoundary`
  (store, new file: sixteen rows expiring at 1600 are listed and capped at 1599, neither at 1600, and a seventeenth is
  stored at 1600) and CX8's `expired at expires_at` row (the daemon's C8). Every time comparison is its own mutant
  (§4 class J), S3's included.
- **S4-D11 — judge-402 note 2 stays a documented residual** (a day-expired request of the same episode is superseded,
  not pruned, by the next put): unchanged, Appendix C.
- **S4-D12 — the card's submit half (M11) mirrors `decide.js`**: `b64u`/`assertionJSON` copied (as design §2.5 says);
  `net(url, init)` gains `init`; a new `#answer` line holds the POST's answer across the re-read (`#status` is the
  re-read's); `reread()` hides the controls, then `load()` shows them again only for a pending card. A passkey button
  is shown only when the card's data has a group for `location.origin`; the SSHSIG paste works everywhere.
- **S4-D13 — timing (directive 8, 12).** No new test sleeps or counts ticks: the daemon tests are synchronous through
  `ServeHTTP` (CX19's concurrency joins on a `WaitGroup`); the node sections wait with `until(ok)` on a 10 s wall-clock
  deadline (the S2 harness's rule), and the ceremony timeout is a stubbed timer fired by hand.

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

| # | premise | measured | routing |
|---|---|---|---|
| F1 | design §2.6/§6: CX19 "the `surface_decide_race_test.go` shape" | that file holds `TestIsApproverRaceOperands` (an operand table); A4's concurrency proof (H-LOCK) is `TestDecideConcurrentPOSTsDecideOnce` in `surface_decide_key_retry_test.go:106` | both shapes used: the operand table for C15's mapping (`TestConsumeRefusalOperands`), H-LOCK + the `TryLock` probes for CX19 (S4-D8) |
| F2 | design §2.6 C13: "below v9 → 503 `Restart`" at the writer's open | the read-only `VerbRequest` at C7 refuses a v8 store first (`verb_requests.go:172`) | C7 answers it (S4-D5), 503 `Restart` naming verb requests; nothing written |
| F3 | directive: `ParseApproveBody` `approve_model.go:197`; routes `daemon.go:1139-1148` | `:214`; `:1141-1153` (S3 added 17 and 3 lines) | cited at their HEAD lines (C2) |
| F4 | design §6 CX21: "the CLI rig waiting on the row dispatches with the second (CF3's rig)" | no CLI rig exists before S5 (`--confirm-on-surface` is S5's) | S4 proves the daemon half (the row names the second decision, the first is named by no row); the CLI half → **S5** (Appendix C) |
| F5 | design §2.6 C14's failure "no decision; row pending" | `RecordOperatorApproval` writes the request object and moves the approvals head before the decision write (`approve.go:723-732`) | CX20 asserts what the design means: 0 decisions, the credentials store unchanged, then 1 decision after the retry (the request object is content-addressed; a later retry writes a new one, as A4's RD13 (a)) |

**Files outside design §7**: `host/store/verb_requests_bound_test.go` (S4-D10) and
`host/surfacepage/confirm_confirmed_test.go` (`Confirmed`'s wire); `host/daemon/surface_confirm_write_test.go` is §7's
"(+ tests)". Pre-existing tests edited outside §7's list: `surface_confirm_read_test.go` (the route table, S3's) and
`confirm_script_test.go` (S2's harness).

### 1.5 The size and the race-mode cost

| stage | production | test | changed lines | cumulative (numstat after commit) |
|---|---|---|---|---|
| M10a | 232 (`surface_confirm.go` 215, `confirm_model.go` 17) | 201 (`surface_confirm_test.go` 175, `confirm_confirmed_test.go` 18, writer sites 8) | 433 | **433** |
| M10b | 0 | 172 (`surface_confirm_test.go` +140, `verb_requests_bound_test.go` 32) | 172 | **605** |
| M10c | 0 | 312 (`surface_confirm_write_test.go`) | 312 | **917** |
| M10d | 3 (`daemon.go`) | 26 (route table 9; `confirmAt` back to the mux, −15/+2) | 29 | **916** |
| M11 | 76 (`confirm.js` 70, `index.html` 6) | 105 (`confirm_script_test.go`) | 181 | **1,097** |

(M10d's cumulative number falls by one: the staged `confirmAt` is longer than the final one.) Design ≈ 945 raw →
compiled **1,097 (1.16×)**. Cap 1,400, margin 303.

**Race-mode wall time of the new host/daemon tests** (directive 8; measured on the work tree, load ≈ 6, `go test -race
-count=1 -v`): `TestConfirmPOSTGuards` 0.24 s, `…RefusesBeforeAnyWrite` 0.76 s, `…Confirms` 0.20 s, `…HoldsHumanMu`
0.17 s, `…RetriesAFailedWorldWrite` 0.15 s, `…RetriesAFailedConsume` 0.24 s, `…RefusedAtC15` 0.29 s,
`TestConsumeRefusalOperands` 0.06 s — **≈ 2.1 s added to host/daemon's race run** (CI's bound is 8 min per package;
host/daemon takes 446–480 s on slow runners: the new tests add < 0.5 %). `-race -count=3` of all eight: 8.0 s
(package wall). No sleep anywhere.

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of four backticks,
in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; a new file's patch has `--- /dev/null`; **no
sub-milestone deletes a file**. Apply it **only** this way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF2'
set -euo pipefail
M=M10a   # <- the sub-milestone, exactly as in the Appendix P heading, in order: M10a M10b M10c M10d M11
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter403-row243-v3s4
PLAN="$W/design_docs/planned/w-world-address-bar-verbs-v3s4-sprint-plan.md"
T="$HOME/.ailang/state/world-iter403/exec-patches"; mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
/usr/bin/patch -p1 -E --dry-run < "$T/$M.patch"
/usr/bin/patch -p1 -E < "$T/$M.patch"
if find cmd host docs \( -name '*.orig' -o -name '*.rej' \) 2>/dev/null | grep .; then echo "STOP: patch left .orig/.rej"; exit 1; fi
git status --short --untracked-files=all
EOF2
```

- If `--dry-run` fails, **stop and report** (no hand edits, no fuzz options). It can only fail if an earlier
  sub-milestone was skipped, reordered or changed.
- Then run that sub-milestone's **Acceptance** (§3): extract `accept.sh` from this plan (`awk '$0 == "````bash
  accept.sh" {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/accept.sh"`) and run **each part §3 lists for M, one
  per Bash call**: `bash "$T/accept.sh" "$W" <M> <part> 2>&1 | tee
  "$HOME/.ailang/state/world-iter403/exec-logs/accept_<M>_<part>.txt"` (never only through `tail`). **Never `nohup`,
  never background a test leg; pass a 600000 ms timeout to every long Bash call.**
- Then snapshot (§7). The controller commits and compares the numstat with §3's expected value.

---

## 3. Milestones

**Acceptance.** Every part runs the same `[base]` checks first (gofmt, build, vet of the three touched packages, the
test-file compile fence `go test -run '^$' ./...`). `v` requires every named test to `--- PASS` and none to SKIP or FAIL;
`p` is a whole-package run. Parts: `fast` (named tests + small packages; ≈ 1 min), `daemon` (`host/daemon` whole,
≈ 4–8 min at load), `chrome` (XC1 in live Chrome, ≈ 45 s). Each part is one 600 s call; measured in C14. Commit-time
test commands for `commit_ms.sh` are the same parts (§7).

````bash accept.sh
#!/bin/bash
# accept.sh <tree> <M> [fast|daemon|chrome]: exactly plan §3's acceptance for sub-milestone M (one bash shell;
# rc 0 = ACCEPT_OK). Each part fits one 600 s call; run every part §3 lists for M, in order.
T=$1; M=$2; P=${3:-fast}
cd "$T" || exit 2
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
unset WORLD_UPDATE_CONFIRM_GOLDEN WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_TOOL_AILANG_BIN
set -o pipefail
chk() { tee /dev/stderr | { out=$(cat); echo "$out" | grep -qE -- '--- (FAIL|SKIP)|^FAIL' && exit 1; echo "$out" | grep -q -- '--- PASS' || exit 1; }; }
v() { go test -count=1 -v "$@" 2>&1 | chk; }   # -v: every named test PASS, none SKIP/FAIL
p() { go test -count=1 "$@"; }
step() { echo "== $*"; "$@" || { echo "ACCEPT_RED $M $P: $*"; exit 1; }; }
# [base]
step test -z "$(gofmt -l cmd host)"
step go build ./...
step go vet ./host/daemon/ ./host/surfacepage/ ./host/store/
if go test -run '^$' ./... 2>&1 | grep -E '^(FAIL|--- FAIL)'; then echo "ACCEPT_RED $M $P: FENCE_RED"; exit 1; fi
ROUTES='TestApprovalRoutesRegistered|TestDecisionRoutesRegistered|TestDecidePageAssets|TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestDaemonCredentialWriterSites|TestDaemonDependencyAllowlist|TestConfirmRoutesRegistered|TestConfirmPageAssetsServed|TestConfirmationCardServed|TestRequestsJSONConfirmations'
A4='TestDecideRefusesReservedScopesByKey|TestDecidePOSTRetryAfterAFailedWorldWrite|TestDecideConcurrentPOSTsDecideOnce|TestDecidePOSTBusyWriter|TestApprovePOSTWriterBusy'
CXA='TestConfirmPOSTGuards'
CXB='TestConfirmPOSTRefusesBeforeAnyWrite'
CXC='TestConfirmPOSTConfirms|TestConfirmPOSTHoldsHumanMu|TestConfirmPOSTRetriesAFailedWorldWrite|TestConfirmPOSTRetriesAFailedConsume|TestConfirmPOSTRefusedAtC15|TestConsumeRefusalOperands'
CP='TestConfirmPageFailsClosed|TestConfirmPageRendersText|TestConfirmPageMessage|TestConfirmPageStates|TestConfirmPageAssets|TestConfirmCardGolden|TestConfirmPageFieldTable|TestConfirmationCard|TestConfirmationCardWire|TestConfirmedWire'
CPS='TestConfirmPageSubmitsAPasskey|TestConfirmPageShowsRefusals|TestConfirmPageSubmitsAnSSHSIG'
case "$M:$P" in
M10a:fast) step v ./host/daemon/ -run "^($ROUTES|$A4|$CXA)\$"
           step v ./host/surfacepage/ -run "^($CP)\$"
           step p ./host/surfacepage/ ;;
M10b:fast) step v ./host/daemon/ -run "^($ROUTES|$CXA|$CXB)\$"
           step v ./host/store/ -run '^(TestVerbRequestExpiryBoundary|TestPutVerbRequestGrammarCapPrune|TestConsumeVerbRequest|TestPutVerbRequestSupersedes)$'
           step p ./host/store/ ;;
M10c:fast|M10d:fast) step v ./host/daemon/ -run "^($ROUTES|$A4|$CXA|$CXB|$CXC)\$" ;;
M11:fast)  step v ./host/surfacepage/ -run "^($CP|$CPS)\$"
           step p ./host/surfacepage/
           step v ./host/daemon/ -run '^(TestConfirmPageAssetsServed|TestDecidePageAssets|TestConfirmRoutesRegistered)$' ;;
M11:chrome) export WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" WORLD_REQUIRE_PASSKEY_BROWSER=1
           step v ./cmd/ailang-worldd/ -run '^TestConfirmPageShowsTheFinishPromptInChrome$' ;;
M10a:daemon|M10d:daemon) step p -timeout 30m ./host/daemon/ ;;
*) echo "unknown $M $P"; exit 2 ;;
esac
echo "ACCEPT_OK $M $P"
````

### M10a — the handler, unregistered; `Confirmed`; the fourth writer site; CX0–CX6 · 232 prod + 201 test · Run 1 · parts: `fast`, `daemon`

- `host/daemon/surface_confirm.go` (NEW): `confirmRequester = "operator:episode-finish"`, `confirmIncomplete`, the
  seams (S4-D7), `confirmAdmission` (C8: `approval.Admit` over the row; consumed → "already confirmed by
  <approved_by>" or "superseded by a newer finish of <episode>"), `confirmStoreError` (older → 503 `Restart` naming verb
  requests; busy → 503 `Busy` + `Retry-After: 1`; else 500 `Internal` + the design's text), `consumeRefusal` (C15:
  superseded → 409 naming the episode; consumed/unknown/expired/clock/digest → 409; `isApproverRace`/`isCounterRefusal`
  → 422; else `confirmStoreError`), `handleSurfaceConfirm` (C0–C15 in the design's order; the writer `defer`-closed
  after C15).
- `host/surfacepage/confirm_model.go`: `ConfirmedJSON{V, Confirmed, Request, ApprovedBy, Decision, Next}` and
  `Confirmed(id, approvedBy, decision)` — design §2.6's "ok" body; `confirm_confirmed_test.go` (NEW): `TestConfirmedWire`
  (the design's literal JSON).
- `host/daemon/surface_approve_post_guard_test.go`: `TestDaemonCredentialWriterSites`' `want` gains
  `"handleSurfaceConfirm": "OpenCredentials"` (AM8), its comment names the four sites.
- `host/daemon/surface_confirm_test.go` (NEW, first half): `newConfirmRig`, `verbMessageOf`, `confirmSSH`,
  `confirmPasskey`, **`confirmAt` (staged: the handler through `approvalRoute`, the path value set)**, `confirm`,
  `confirmState`/`credState`; **`TestConfirmPOSTGuards`** (CX0–CX6: C0 ×2 and C0-before-C1, C1 ×3 and C1-before-C3, C3
  missing/`null`/foreign, C4 ×2, C5, C6 oversize/`typed`/`y`/both/neither/`v` 2/two values, R2; nothing written; the
  same body then 200).
- **Expected numstat after commit: 433.**

### M10b — CX7–CX13 and judge-402 note 1 · 0 prod + 172 test · Run 1 · part: `fast`

- `host/daemon/surface_confirm_test.go` (second half): **`TestConfirmPOSTRefusesBeforeAnyWrite`** — C7 unknown (404),
  C8 clock behind / expired **at exactly `expires_at`** / a hand-edited digest / superseded / consumed (409, texts),
  C9 (409) and C9-before-C11, C10 decode / `Validate` / F1–F5 (409 `NotConfirmable`; each row's signature valid), C11
  unknown key / namespace / parse / W6 (a mint-hash assertion) / W15 (counter 1) / revoked / the judge's store class
  (503), C12 through the judge seam (422 naming the field), C13 busy (503 `Busy`, `Retry-After: 1`), C7 below schema 9
  (503 `Restart`) and unreadable (500); every refusal compared with `confirmState` before.
- `host/store/verb_requests_bound_test.go` (NEW): **`TestVerbRequestExpiryBoundary`** (S4-D10).
- **Expected numstat after commit: 605.**

### M10c — the writes: CX14–CX21, audit note 3, C15's mapping · 0 prod + 312 test · Run 1 · part: `fast`

- `host/daemon/surface_confirm_write_test.go` (NEW): `worldDecision` (the decision's verb and decider, its request's
  effect, requester and scope, read back through `GetObject`, `broker.DecodeApprovalRequest` and
  `broker.ApprovalChainState`); **`TestConfirmPOSTConfirms`** (CX17: SSHSIG at `apT0` and a ready-PR passkey one second
  before expiry; body = `json.Marshal(surfacepage.Confirmed(id, by, row.DecisionRef))`; `decision = approve`,
  `decidedBy = ssh:mark-ssh` / `webauthn:mark`, `requester = operator:episode-finish`, `scope = verbFixture`'s; five
  columns; counter 1 → 2; last use; replay 409 writing nothing), **`TestConfirmPOSTHoldsHumanMu`** (CX19, S4-D8),
  **`TestConfirmPOSTRetriesAFailedWorldWrite`** (CX20, passkey), **`TestConfirmPOSTRetriesAFailedConsume`** (CX21: same
  second → 1 decision; a second later → 2, the row names the second), **`TestConfirmPOSTRefusedAtC15`** (audit note 3 and
  C15's refusals, S4-D9), **`TestConsumeRefusalOperands`** (each sentinel → its status).
- **Expected numstat after commit: 917.**

### M10d — the one POST registered · 3 prod + 26 test · Run 1 · parts: `fast`, `daemon`

- `host/daemon/daemon.go`: `mux.HandleFunc("POST /surface/confirm/{id}", d.approvalRoute(d.handleSurfaceConfirm))`
  after `GET /surface/confirm/{id}`; the Handler comment names it.
- `host/daemon/surface_confirm_read_test.go`: `TestConfirmRoutesRegistered`' `want` gains the POST; its 405 loop skips
  POST on the card's path only.
- `host/daemon/surface_confirm_test.go`: `confirmAt` → `a.serve("POST", "/surface/confirm/"+id, host, …)` (through the
  mux); the three staging imports go.
- **Expected numstat after commit: 916.**

### M11 — `confirm.js`'s submit half (design §2.5 sprint 4, §6 CP8–CP10) · 76 prod + 105 test · Run 2 · parts: `fast`, `chrome`

- `host/surfacepage/confirm/confirm.js`: `TIMEOUT_MS = 240000`, `b64u`, `assertionJSON`, `net(url, init)`;
  `digestHex`/`hash`/`group` kept from the pending card; `buttons(on)`; `submit(format)` (POST `{v: 1, digest} + format`,
  the answer or "Not confirmed (HTTP n class): reason" — a 200 counts only with `confirmed === true` — then `reread()`);
  `confirmWithPasskey()` (`navigator.credentials.get({publicKey: {challenge: hash, rpId: group.rp, allowCredentials,
  userVerification: 'required', timeout}, signal})`, aborted at 240 s; a failure sends nothing and re-reads).
- `host/surfacepage/confirm/index.html`: `#answer`; inside `#controls` the passkey button (`hidden` until a group
  exists), the signature textarea and its button.
- `host/surfacepage/confirm_script_test.go`: the stand-in DOM's `addEventListener`/`click()`, `open(…, navigator)`,
  `until`, `daemon`, `posts`, `CRED`; **`TestConfirmPageSubmitsAPasskey`** (CP8), **`TestConfirmPageShowsRefusals`**
  (CP9), **`TestConfirmPageSubmitsAnSSHSIG`** (CP10).
- **Expected numstat after commit: 1,097.**

### 3.1 Every input and failure class S4's code reaches, traced (design §2.9, §2.6's window table; directive retro 4)

| design row / input | S4 branch | test (row) |
|---|---|---|
| P5 C0 `Authorization` / `Cookie` | C0 → 403 `Forbidden` | CX0 `C0 authorization`, `C0 cookie`, `C0 before C1` |
| P5 C1 id grammar | C1 → 400 | CX1 `C1 short`/`upper`/`mint id`, `C1 before C3` |
| P5 R2 foreign Host (DNS rebinding, T7) | `approvalRoute` → 421 | CX `R2 foreign host` |
| P5 C3–C5 Origin missing/`null`/foreign, `Sec-Fetch-Site`, media type | C3/C4 → 403, C5 → 415 | CX3–CX5 rows |
| P5 C6 oversize; unknown key, `typed`, `y`, both, neither, `v`, two values (T6: no typed arm) | C6 → 413 / 400 | CX6 rows |
| P5 C7 unknown; below v9; store unreadable | 404; 503 `Restart`; 500 | CX7 `unknown`, `v8`, `unreadable` |
| P5 C8 consumed / superseded / expired (at `expires_at`) / clock / digest; T2 replay; T4 | `confirmAdmission` → 409 | CX8 rows; CX17 replay |
| P5 C9 digest ≠ the card's | 409 | CX9, `C9 before C11` |
| P5 C10 decode / `Validate` / F1–F5 (T8: the prompt is the terminal's rendering) | 409 `NotConfirmable` | CX10 rows (7) |
| P5 C11 unknown key (T1), revoked, namespace (T5), parse, W6 (T5), W15, store class | 422 / 503 | CX11 rows |
| P5 C12 binding (T3: another request's signature) | 422 | CX12 (judge seam) |
| P5 C13 busy | 503 `Busy`, `Retry-After: 1`, nothing written | CX13 |
| P5 C14 fails (window row 2) | 500, 0 decisions, row pending | CX20 |
| P5 C15 fails after C14 (window row 3, RV2) | 500, 1 orphan, row pending; retry → 200 | CX21 (same second, a second later) |
| **audit note 3: superseded between C14 and C15** | 409 naming the episode; orphan inert; row `superseded` | `TestConfirmPOSTRefusedAtC15` `superseded` |
| P5 C15 a concurrent consume lost / approver revoked (T11) / counter (T11) / store error | 409 / 422 / 422 / 500 | `TestConfirmPOSTRefusedAtC15`; each sentinel `TestConsumeRefusalOperands` |
| P5 ok (window row 4) | 200 `Confirmed`; `decidedBy` the key; requester `operator:episode-finish`; scope the request's | CX17 |
| T9 / §5.2 "the scope the daemon records is the request's" | C14 passes `req.Scope` | CX17 `scope == verbFixture's` |
| T10 restart mid-POST | every window state is rows; retries complete | CX20, CX21 |
| CX18 A4's decide still refuses a finish-shaped request (AM6, L1 unchanged) | unchanged | `TestDecideRefusesReservedScopesByKey` (unedited, in `fast`) |
| CX19 serialization with the typed approve phase | `humanMu` held clock → C11 → C15 | CX19 |
| P4 passkey ceremony aborted / 240 s / refused | "Not confirmed: the browser made no passkey response"; nothing sent; re-read, controls back | CP8 |
| P4 POST refusal (any C step), a lost POST, a 200 without `confirmed` | "Not confirmed (HTTP …)"; re-read, still pending | CP9 |
| P4 SSHSIG paste | `{v, digest, sshsig}` | CP10 |
| P4 no passkey group for this origin | no passkey button; the paste still works | CP8 `elsewhere` |
| judge-402 note 1 | `expires_at == now` expired in list, count and consume | `TestVerbRequestExpiryBoundary`; CX8 `expired at expires_at` |

---

## 4. Mutation table

**Procedure** (`tools/mutate.py <run> <workers> [ids…]`, rows in `tools/mutants.py`, **enumerated from the diff
`BASE..M11`, not from the test plan**): **each step C0–C15 deleted** (C14 by a world write that never happens, C15 by a
consume that never runs), **each status code** (C0–C11, the store classes, C14, C15's two classes, the 200), **the step
order where it is load-bearing** (C0 after C1, C9 after C11, C13 after C14, C15 before C14, `humanMu` taken after C11,
the writer closed before C14, never closed), **`humanMu`** (removed, unlocked at once), the `commitBudget` context,
**the requester `"operator:episode-finish"`**, C14's `decidedBy` and scope, **the `decision_ref` passed to consume** and
each other consume argument, **each sentinel → status mapping** (every `errors.Is` operand of `consumeRefusal`, the
older-store / busy / `Retry-After` / default branches of `confirmStoreError`, the judge's store class), **the strict body
decoder** (`typed`, `y`, both formats: C6's parse and size), each `confirmAdmission` input and text, each
`FinishConfirmFacts` field, each `VerbJudgeInput` field, `Binds`' expiry, **the route registration** (removed, without
`approvalRoute`, any method, an extra POST), **each `Confirmed` field and a tag**, and **each branch of `confirm.js`'s
submit half**; plus, per the directive's new rule, **every time comparison's operator as its own mutant** — S4's diff adds
none in Go (C8 reuses `approval.Admit`; its inputs are mutated by ±1 instead), so class J is S3's verb store (the list
and count `expires_at > ?` — the sprint-3 judge's two survivors —, the prune `<=`, consume's `<` and `>=`, the row
grammar's `<=`) and F2's `<` and `<=`. For each row: an APFS clone of `work` at `M11` (`cp -cR`; byte-identical to the
round-trip tree), its exact-string replacements (each anchor checked to occur exactly once: 0 bad anchors), then
**one** `go test -count=1 -v` of the row's package suite — `host/surfacepage` whole (node harnesses included, with
`WORLD_EXEC_NODE`), `host/store` whole, `host/broker` `^TestFinish`, `host/daemon` the confirm/approval/decision/route
subset (run R1: `^(TestConfirm|TestConsumeRefusal|TestDaemonCredentialWriterSites|TestDecide|TestApprove|…)`; run R2 on:
`^(TestConfirm|TestConsumeRefusal|TestDaemonCredentialWriterSites|TestDecideRefusesReservedScopesByKey|TestDecideConcurrentPOSTsDecideOnce|TestDecidePOSTBusyWriter|TestApprovePOSTWriterBusy|TestApprovalRoutes|TestDecisionRoutes|TestFrozen|TestRequestsJSON|TestIsApproverRace|TestIsCounterRefusal)`,
narrowed for time at load) — from whose `--- FAIL`/`--- PASS` lines the **named test**'s result, the **`suite`** (every
failing test) and **`suite-nostatic`** (the failures other than the route-table, inventory and static tests
`TestConfirmRoutesRegistered`, `TestApprovalRoutesRegistered`, `TestDecisionRoutesRegistered`,
`TestFrozenV1RouteTableMatchesMux`, `TestDaemonCredentialWriterSites`, `TestConfirmPageAssets`,
`TestConfirmPageFieldTable`, `TestDecidePageAssets`) are read — the behavioural proof the directive asks for. Rows ran
with `GOFLAGS=-vet=off`, 3–5 workers, beside the base gates. Every run appends one line per row to the ledger (time,
run, id, file, expected, named result, named failures, suites, named regex, seconds); per-row output in
`logs/mut/<run>-<id>.log`. **The ids below are the ledger's ids**; the table shows each row's last line. Classes: `K`
killed by the named test **and** by the behavioural suite; `R` the route-table proof (a registration only the
route-table test can see: any-method and an extra POST — the behavioural CX rows cannot observe a second route they
never call); `E` equivalent (reason given); `CTL` the unmutated tree, green.

**R1 → R2 → R3 (what changed).** R1 (21 rows, 3 workers, the wider daemon subset) was stopped and continued as R2
(the other 129 rows, 5 workers, the narrower subset) for time at load. R2 found four rows not as expected:
**`C11-CHALLENGE`** did not compile (deleting the copy left `c` unused; a non-compiling row is not a mutant) — re-made as
`copy(challenge[:], c[:0])`; **`C11-ARM`** panicked in an earlier test (a passkey body through the SSH arm dereferences a
nil `sshsig`), so its named test never ran — its named test is now `TestConfirmPOSTGuards`, whose control the swapped
arm refuses; **`C10-DECODE`** and **`C14-ERR`** survived everything and are **equivalent** (reasons in the table: the
zero request `DecodeVerb` returns is refused by `Validate`; the zero ref a failed C14 returns is refused by C15's
grammar before any write, with the same 500). R3 re-ran those four. Earlier, the enumeration itself hardened four test
rows before any run (C8's superseded and consumed rows compare the stores, C11's judge store class, C7 unreadable,
CP10's in-flight check, CP8's other-origin group and its `waiting` state). R1's and R2's lines stay in the ledger.

**Summary:** 150 rows (3 `CTL`, 4 `E`, 141 `K`, 2 `R`); 154 ledger lines (runs R1, R2, R3); 150 of 150 rows run; the last line of each row behaves as its class says for **150 of 150**. Table ids = ledger ids (`tools/mktable.py` fills §4 from the ledger's last line per id); ledger `$HOME/.ailang/state/world-iter403/plan-scratch/logs/mutant_ledger.tsv`

| id | file | from → to | named test → result | suites (`suite-nostatic`: route-table/static tests excluded) | class | as expected |
|---|---|---|---|---|---|---|
| `CTL-DAEMON` | `surface_confirm.go` | control (no change) | `(TestConfirmPOST|TestConsumeRefusal)` → **PASS** | suite=PASS; suite-nostatic=PASS | CTL | yes |
| `CTL-PAGE` | `confirm.js` | control (no change) | `TestConfirmPage` → **PASS** | suite=PASS; suite-nostatic=PASS | CTL | yes |
| `CTL-STORE` | `verb_requests.go` | control (no change) | `TestVerbRequestExpiryBoundary` → **PASS** | suite=PASS; suite-nostatic=PASS | CTL | yes |
| `C0-DEL` | `surface_confirm.go` | `if len(r.Header.Values("Authorization")) != 0 \|\| len(r.…` → `if false { // C0` | `TestConfirmPOSTGuards` → **KILLED** | suite=KILLED[TestConfirmPOSTGuards]; suite-nostatic=KILLED[TestConfirmPOSTGuards] | K | yes |
| `C0-AUTH` | `surface_confirm.go` | `len(r.Header.Values("Authorization")) != 0 \|\| len(r.Hea…` → `len(r.Header.Values("Cookie")) != 0 { // C0` | `TestConfirmPOSTGuards` → **KILLED** | suite=KILLED[TestConfirmPOSTGuards]; suite-nostatic=KILLED[TestConfirmPOSTGuards] | K | yes |
| `C0-COOKIE` | `surface_confirm.go` | `…ization")) != 0 \|\| len(r.Header.Values("Cookie")) != 0…` → `…ization")) != 0 { // C0` | `TestConfirmPOSTGuards` → **KILLED** | suite=KILLED[TestConfirmPOSTGuards]; suite-nostatic=KILLED[TestConfirmPOSTGuards] | K | yes |
| `C0-STATUS` | `surface_confirm.go` | `http.StatusForbidden, "Forbidden", "a confirmation is nev…` → `http.StatusBadRequest, "Forbidden", "a confirmation is ne…` | `TestConfirmPOSTGuards` → **KILLED** | suite=KILLED[TestConfirmPOSTGuards]; suite-nostatic=KILLED[TestConfirmPOSTGuards] | K | yes |
| `C0-AFTER-C1` | `surface_confirm.go` | `if len(r.Header.Values("Authorization")) != 0 \|\| len(r.…` → `` (+1 more) | `TestConfirmPOSTGuards` → **KILLED** | suite=KILLED[TestConfirmPOSTGuards]; suite-nostatic=KILLED[TestConfirmPOSTGuards] | K | yes |
| `C1-DEL` | `surface_confirm.go` | `if !store.VerbRequestIDPattern.MatchString(id) { // C1` → `if false { // C1` | `TestConfirmPOSTGuards` → **KILLED** | suite=KILLED[TestConfirmPOSTGuards]; suite-nostatic=KILLED[TestConfirmPOSTGuards] | K | yes |
| `C1-STATUS` | `surface_confirm.go` | `…l(w, http.StatusBadRequest, "BadRequest", confirmIDReaso…` → `…l(w, http.StatusNotFound, "BadRequest", confirmIDReason)` | `TestConfirmPOSTGuards` → **KILLED** | suite=KILLED[TestConfirmPOSTGuards]; suite-nostatic=KILLED[TestConfirmPOSTGuards] | K | yes |
| `C3-DEL` | `surface_confirm.go` | `….Get("Origin"); !allowedApprovalOrigin(origin, o) { // C3` → `….Get("Origin"); false && origin == "" { // C3` | `TestConfirmPOSTGuards` → **KILLED** | suite=KILLED[TestConfirmPOSTGuards]; suite-nostatic=KILLED[TestConfirmPOSTGuards] | K | yes |
| `C3-STATUS` | `surface_confirm.go` | `…l(w, http.StatusForbidden, "CrossOrigin", crossOriginRea…` → `…l(w, http.StatusBadRequest, "CrossOrigin", crossOriginRe…` | `TestConfirmPOSTGuards` → **KILLED** | suite=KILLED[TestConfirmPOSTGuards]; suite-nostatic=KILLED[TestConfirmPOSTGuards] | K | yes |
| `C4-DEL` | `surface_confirm.go` | `if !sameOriginFetch(r.Header.Get("Sec-Fetch-Site")) { // …` → `if false { // C4` | `TestConfirmPOSTGuards` → **KILLED** | suite=KILLED[TestConfirmPOSTGuards]; suite-nostatic=KILLED[TestConfirmPOSTGuards] | K | yes |
| `C4-STATUS` | `surface_confirm.go` | `http.StatusForbidden, "CrossOrigin", "Sec-Fetch-Site is n…` → `http.StatusBadRequest, "CrossOrigin", "Sec-Fetch-Site is …` | `TestConfirmPOSTGuards` → **KILLED** | suite=KILLED[TestConfirmPOSTGuards]; suite-nostatic=KILLED[TestConfirmPOSTGuards] | K | yes |
| `C5-DEL` | `surface_confirm.go` | `if !jsonMediaType(r.Header.Get("Content-Type")) { // C5` → `if false { // C5` | `TestConfirmPOSTGuards` → **KILLED** | suite=KILLED[TestConfirmPOSTGuards]; suite-nostatic=KILLED[TestConfirmPOSTGuards] | K | yes |
| `C5-STATUS` | `surface_confirm.go` | `http.StatusUnsupportedMediaType, "UnsupportedMediaType"` → `http.StatusBadRequest, "UnsupportedMediaType"` | `TestConfirmPOSTGuards` → **KILLED** | suite=KILLED[TestConfirmPOSTGuards]; suite-nostatic=KILLED[TestConfirmPOSTGuards] | K | yes |
| `C6-SIZE` | `surface_confirm.go` | `io.ReadAll(http.MaxBytesReader(w, r.Body, surfacepage.Max…` → `io.ReadAll(r.Body) // C6` | `TestConfirmPOSTGuards` → **KILLED** | suite=KILLED[TestConfirmPOSTGuards]; suite-nostatic=KILLED[TestConfirmPOSTGuards] | K | yes |
| `C6-SIZE-STATUS` | `surface_confirm.go` | `http.StatusRequestEntityTooLarge, "TooLarge"` → `http.StatusBadRequest, "TooLarge"` | `TestConfirmPOSTGuards` → **KILLED** | suite=KILLED[TestConfirmPOSTGuards]; suite-nostatic=KILLED[TestConfirmPOSTGuards] | K | yes |
| `C6-PARSE` | `surface_confirm.go` | `if err != nil {⏎writeApprovalRefusal(w, http.StatusBadReq…` → `if false {⏎writeApprovalRefusal(w, http.StatusBadRequest,…` | `TestConfirmPOSTGuards` → **KILLED** | suite=KILLED[TestConfirmPOSTGuards]; suite-nostatic=KILLED[TestConfirmPOSTGuards] | K | yes |
| `C6-PARSE-STATUS` | `surface_confirm.go` | `…l(w, http.StatusBadRequest, "BadRequest", err.Error())` → `…l(w, http.StatusUnprocessableEntity, "BadRequest", err.E…` | `TestConfirmPOSTGuards` → **KILLED** | suite=KILLED[TestConfirmPOSTGuards]; suite-nostatic=KILLED[TestConfirmPOSTGuards] | K | yes |
| `H-LOCK` | `surface_confirm.go` | `d.humanMu.Lock()⏎defer d.humanMu.Unlock()⏎` → `` | `TestConfirmPOSTHoldsHumanMu` → **KILLED** | suite=KILLED[TestConfirmPOSTHoldsHumanMu]; suite-nostatic=KILLED[TestConfirmPOSTHoldsHumanMu] | K | yes |
| `H-UNLOCK-EARLY` | `surface_confirm.go` | `defer d.humanMu.Unlock()⏎` → `d.humanMu.Unlock()⏎` | `TestConfirmPOSTHoldsHumanMu` → **KILLED** | suite=KILLED[TestConfirmPOSTHoldsHumanMu]; suite-nostatic=KILLED[TestConfirmPOSTHoldsHumanMu] | K | yes |
| `H-LOCK-AFTER-C11` | `surface_confirm.go` | `d.humanMu.Lock()⏎defer d.humanMu.Unlock()⏎` → `` (+1 more) | `TestConfirmPOSTHoldsHumanMu` → **KILLED** | suite=KILLED[TestConfirmPOSTHoldsHumanMu]; suite-nostatic=KILLED[TestConfirmPOSTHoldsHumanMu] | K | yes |
| `CTX-BUDGET` | `surface_confirm.go` | `… := context.WithTimeout(r.Context(), d.commitBudget)` → `… := context.WithCancel(r.Context())` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4] | K | yes |
| `C7-RO-STATUS` | `surface_confirm.go` | `if err != nil {⏎d.approvalInternal(w, r, err)⏎return⏎}⏎de…` → `if err != nil {⏎writeApprovalRefusal(w, http.StatusServic…` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `C7-OLDER` | `surface_confirm.go` | `case err != nil:⏎d.confirmStoreError(w, r, err)⏎return⏎ca…` → `case err != nil:⏎d.approvalInternal(w, r, err)⏎return⏎cas…` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `C7-DEL` | `surface_confirm.go` | `case !ok:⏎writeApprovalRefusal(w, http.StatusNotFound` → `case false && !ok:⏎writeApprovalRefusal(w, http.StatusNot…` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `C7-STATUS` | `surface_confirm.go` | `http.StatusNotFound, "NotFound", "no verb request "+id` → `http.StatusConflict, "NotFound", "no verb request "+id` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `C8-DEL` | `surface_confirm.go` | `…sion(row, now); reason != "" { // C8` → `…sion(row, now); false && reason != "" { // C8` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTRefusesBeforeAnyWrite,TestConfirmPOSTRetriesAFailedWorldWrite]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTRefusesBeforeAnyWrite,TestConfirmPOSTRetriesAFailedWorldWrite] | K | yes |
| `C8-STATUS` | `surface_confirm.go` | `…l(w, http.StatusConflict, "Conflict", reason)` → `…l(w, http.StatusUnprocessableEntity, "Conflict", reason)` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTHoldsHumanMu,TestConfirmPOSTRefusesBeforeAnyWrite,…+1]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTHoldsHumanMu,TestConfirmPOSTRefusesBeforeAnyWrite,…+1] | K | yes |
| `ADM-CONSUMED` | `surface_confirm.go` | `Consumed: row.Outcome != "",` → `Consumed: false,` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTRefusesBeforeAnyWrite,TestConfirmPOSTRetriesAFailedWorldWrite]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTRefusesBeforeAnyWrite,TestConfirmPOSTRetriesAFailedWorldWrite] | K | yes |
| `ADM-NOW` | `surface_confirm.go` | `approval.AdmitInput{Now: now,` → `approval.AdmitInput{Now: now - 1,` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `ADM-CREATED` | `surface_confirm.go` | `CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt, Consu…` → `CreatedAt: 0, ExpiresAt: row.ExpiresAt, Consumed:` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `ADM-EXPIRES` | `surface_confirm.go` | `…t: row.ExpiresAt, Consumed:` → `…t: row.ExpiresAt + 1, Consumed:` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `ADM-DIGEST` | `surface_confirm.go` | `StoredDigest: digest, BytesDigest: sha256.Sum256(row.Requ…` → `StoredDigest: sha256.Sum256(row.Request), BytesDigest: sh…` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `ADM-SUPERSEDED` | `surface_confirm.go` | `if row.Outcome == store.VerbOutcomeSuperseded {` → `if false {` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `ADM-BY` | `surface_confirm.go` | `" was already confirmed by " + row.ApprovedBy` → `" was already confirmed by " + row.EpisodeID` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `C9-DEL` | `surface_confirm.go` | `if body.Digest != row.Digest { // C9` → `if false { // C9` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `C9-STATUS` | `surface_confirm.go` | `http.StatusConflict, "Conflict", "the request changed sin…` → `http.StatusUnprocessableEntity, "Conflict", "the request …` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `ORD-C9-AFTER-C11` | `surface_confirm.go` | `if body.Digest != row.Digest { // C9⏎writeApprovalRefusal…` → `` (+1 more) | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `C10-DEL` | `surface_confirm.go` | `if err != nil {⏎writeApprovalRefusal(w, http.StatusConfli…` → `if false {⏎writeApprovalRefusal(w, http.StatusConflict, "…` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `C10-STATUS` | `surface_confirm.go` | `http.StatusConflict, "NotConfirmable"` → `http.StatusUnprocessableEntity, "NotConfirmable"` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `C10-DECODE` | `surface_confirm.go` | `req, err := approval.DecodeVerb(row.Request) // C10` → `req, _ := approval.DecodeVerb(row.Request) // C10` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **PASS** | suite=PASS; suite-nostatic=PASS | E — `DecodeVerb` returns the zero `VerbRequest` with its error, and `Validate` refuses the zero request at V1 (schema), so a decode failure is still C10's 409 with the request unchanged; the decode error is kept for its text (CX10's decode row names neither) | yes |
| `C10-VALIDATE` | `surface_confirm.go` | `err = req.Validate()` → `err = nil` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `C10-FINISH` | `surface_confirm.go` | `RequestExpires: row.ExpiresAt, Now: now}); reason != "" {` → `RequestExpires: row.ExpiresAt, Now: now}); false && reaso…` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `C10-REQEXP` | `surface_confirm.go` | `RequestExpires: row.ExpiresAt, Now: now})` → `RequestExpires: row.ExpiresAt - 1, Now: now})` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `C10-NOW` | `surface_confirm.go` | `RequestExpires: row.ExpiresAt, Now: now})` → `RequestExpires: row.ExpiresAt, Now: 0})` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4] | K | yes |
| `C10-PROTECTED` | `surface_confirm.go` | `Protected: req.Protected,` → `Protected: nil,` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4] | K | yes |
| `C10-COMMITTER` | `surface_confirm.go` | `Committer: req.Committer,` → `Committer: "x",` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4] | K | yes |
| `C10-PROMPT` | `surface_confirm.go` | `Prompt: req.Prompt,` → `Prompt: req.Title,` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4] | K | yes |
| `C10-TITLE` | `surface_confirm.go` | `Title: req.Title, Body:` → `Title: req.Body, Body:` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4] | K | yes |
| `C10-BODY` | `surface_confirm.go` | `Body: req.Body, Prompt:` → `Body: req.Title, Prompt:` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4] | K | yes |
| `C10-SCOPE` | `surface_confirm.go` | `broker.FinishConfirmFacts{Scope: req.Scope,` → `broker.FinishConfirmFacts{Scope: req.Title,` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4] | K | yes |
| `C11-ARM` | `surface_confirm.go` | `if body.SSHSIG != nil {⏎judge, artifact, what = judgeVerb…` → `if body.SSHSIG == nil {⏎judge, artifact, what = judgeVerb…` | `TestConfirmPOSTGuards` → **KILLED** | suite=KILLED[TestConfirmPOSTGuards,TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTGuards,TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `C11-CHALLENGE` | `surface_confirm.go` | `copy(challenge[:], c)⏎` → `copy(challenge[:], c[:0])⏎` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4] | K | yes |
| `C11-EXPIRES` | `surface_confirm.go` | `…t: row.ExpiresAt, Now: now}` → `…t: row.ExpiresAt + 1, Now: now}` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4] | K | yes |
| `C11-NOW` | `surface_confirm.go` | `…ExpiresAt, Now: now}` → `…ExpiresAt, Now: 0}` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4] | K | yes |
| `C11-BYTES` | `surface_confirm.go` | `BytesDigest: digest, CreatedAt` → `BytesDigest: [32]byte{}, CreatedAt` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4] | K | yes |
| `C11-CREATED` | `surface_confirm.go` | `BytesDigest: digest, CreatedAt: row.CreatedAt,` → `BytesDigest: digest, CreatedAt: 0,` | `TestConfirmPOSTConfirms` → **PASS** | suite=PASS; suite-nostatic=PASS | E — the judge's P1 clock check is `approval.Admit` over the same `created_at` that C8 has already refused when `now < created_at`; every row that reaches C11 has `now >= created_at`, so `CreatedAt: 0` admits exactly the rows the real value admits | yes |
| `C11-STORE-CLASS` | `surface_confirm.go` | `case errors.As(err, &je) && je.Step == authority.JudgeSte…` → `case false && errors.As(err, &je):` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `C11-STATUS` | `surface_confirm.go` | `http.StatusUnprocessableEntity, "NotApproved", "the "+wha…` → `http.StatusConflict, "NotApproved", "the "+what+" does no…` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `C11-DEL` | `surface_confirm.go` | `case err != nil:⏎writeApprovalRefusal(w, http.StatusUnpro…` → `case false && err != nil:⏎writeApprovalRefusal(w, http.St…` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `C12-DEL` | `surface_confirm.go` | `err = signed.Binds(digest, challenge, row.ExpiresAt) // C…` → `err = nil // C12` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `C12-EXPIRY` | `surface_confirm.go` | `…e, row.ExpiresAt) // C12` → `…e, row.ExpiresAt+1) // C12` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4] | K | yes |
| `C13-ERR` | `surface_confirm.go` | `if err != nil {⏎d.confirmStoreError(w, r, err)⏎return⏎}⏎d…` → `if false {⏎d.confirmStoreError(w, r, err)⏎return⏎}⏎defer …` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `ORD-C13-AFTER-C14` | `surface_confirm.go` | `wr, err := store.OpenCredentials(d.credentialsPath) // C1…` → `` (+1 more) | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `WR-CLOSE-EARLY` | `surface_confirm.go` | `defer func() { _ = wr.Close() }() // held` → `_ = wr.Close() // held` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4] | K | yes |
| `WR-NO-CLOSE` | `surface_confirm.go` | `defer func() { _ = wr.Close() }() // held` → `// held` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTHoldsHumanMu,TestConfirmPOSTRefusedAtC15,…+3]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTHoldsHumanMu,TestConfirmPOSTRefusedAtC15,…+3] | K | yes |
| `STORE-OLDER` | `surface_confirm.go` | `case isOlderStore(err):⏎writeApprovalRefusal(w, http.Stat…` → `case false:⏎writeApprovalRefusal(w, http.StatusServiceUna…` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite,TestConsumeRefusalOperands]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite,TestConsumeRefusalOperands] | K | yes |
| `STORE-BUSY` | `surface_confirm.go` | `case store.IsWriterAlreadyActive(err):⏎w.Header().Set("Re…` → `case false:⏎w.Header().Set("Retry-After", "1")⏎writeAppro…` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `STORE-RETRY-AFTER` | `surface_confirm.go` | `…adyActive(err):⏎w.Header().Set("Retry-After", "1")⏎` → `…adyActive(err):⏎` | `TestConfirmPOSTRefusesBeforeAnyWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite]; suite-nostatic=KILLED[TestConfirmPOSTRefusesBeforeAnyWrite] | K | yes |
| `STORE-DEFAULT-STATUS` | `surface_confirm.go` | `…l(w, http.StatusInternalServerError, "Internal", confirm…` → `…l(w, http.StatusServiceUnavailable, "Internal", confirmI…` | `TestConsumeRefusalOperands` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusedAtC15,TestConfirmPOSTRetriesAFailedConsume,TestConsumeRefusalOperands]; suite-nostatic=KILLED[TestConfirmPOSTRefusedAtC15,TestConfirmPOSTRetriesAFailedConsume,TestConsumeRefusalOperands] | K | yes |
| `C14-DEL` | `surface_confirm.go` | `_, decision, err := broker.RecordOperatorApproval(ctx, d.…` → `decision, err := hashref.SumSHA256([]byte(req.Scope)), er…` (+1 more) | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTHoldsHumanMu,TestConfirmPOSTRefusedAtC15,…+2]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTHoldsHumanMu,TestConfirmPOSTRefusedAtC15,…+2] | K | yes |
| `C14-REQUESTER` | `surface_confirm.go` | `const confirmRequester = "operator:episode-finish"` → `const confirmRequester = "operator:terminal"` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms]; suite-nostatic=KILLED[TestConfirmPOSTConfirms] | K | yes |
| `C14-DECIDEDBY` | `surface_confirm.go` | `…nfirmRequester, signed.ApprovedBy(), now) // C14` → `…nfirmRequester, "operator:terminal", now) // C14` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms]; suite-nostatic=KILLED[TestConfirmPOSTConfirms] | K | yes |
| `C14-SCOPE` | `surface_confirm.go` | `d.store, req.Scope, confirmRequester` → `d.store, req.Scope+" ", confirmRequester` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms]; suite-nostatic=KILLED[TestConfirmPOSTConfirms] | K | yes |
| `C14-STATUS` | `surface_confirm.go` | `…l(w, http.StatusInternalServerError, "Internal", confirm…` → `…l(w, http.StatusServiceUnavailable, "Internal", confirmI…` | `TestConfirmPOSTRetriesAFailedWorldWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTRetriesAFailedWorldWrite]; suite-nostatic=KILLED[TestConfirmPOSTRetriesAFailedWorldWrite] | K | yes |
| `C14-ERR` | `surface_confirm.go` | `if err != nil {⏎d.writeInternalErrorLog(r, err)⏎writeAppr…` → `if false {⏎d.writeInternalErrorLog(r, err)⏎writeApprovalR…` | `TestConfirmPOSTRetriesAFailedWorldWrite` → **PASS** | suite=PASS; suite-nostatic=PASS | E — a failed C14 returns the zero decision ref (every error path of `RecordOperatorApproval`/`decideApproval` does), and C15's grammar (`checkVerbConsume`, `isDecisionRef`) refuses it before any write, which `consumeRefusal` answers with the same 500 and text: nothing is consumed either way (defence in depth) | yes |
| `C15-DEL` | `surface_confirm.go` | `if err := consumeVerbRequest(wr, ctx, id, row.Digest, use…` → `if err := error(nil); err != nil { // C15` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTHoldsHumanMu,TestConfirmPOSTRefusedAtC15,…+3]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTHoldsHumanMu,TestConfirmPOSTRefusedAtC15,…+3] | K | yes |
| `ORD-C15-BEFORE-C14` | `surface_confirm.go` | `_, decision, err := broker.RecordOperatorApproval(ctx, d.…` → `name0, format0, key0, _ := signed.Use()⏎use0 := store.App…` (+1 more) | `TestConfirmPOSTRetriesAFailedWorldWrite` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTRefusedAtC15,TestConfirmPOSTRetriesAFailedConsume,…+1]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTRefusedAtC15,TestConfirmPOSTRetriesAFailedConsume,…+1] | K | yes |
| `C15-REF` | `surface_confirm.go` | `artifact, decision.String(), now); err != nil { // C15` → `artifact, "sha256:"+row.Digest, now); err != nil { // C15` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTRetriesAFailedConsume,TestConfirmPOSTRetriesAFailedWorldWrite]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTRetriesAFailedConsume,TestConfirmPOSTRetriesAFailedWorldWrite] | K | yes |
| `C15-DIGEST` | `surface_confirm.go` | `consumeVerbRequest(wr, ctx, id, row.Digest, use,` → `consumeVerbRequest(wr, ctx, id, body.Digest, use,` | `TestConfirmPOSTConfirms` → **PASS** | suite=PASS; suite-nostatic=PASS | E — C9 has just refused `body.Digest != row.Digest`, so at C15 the two strings are equal | yes |
| `C15-ARTIFACT` | `surface_confirm.go` | `signed.ApprovedBy(), artifact, decision.String(), now); e…` → `signed.ApprovedBy(), append(artifact, 0), decision.String…` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms]; suite-nostatic=KILLED[TestConfirmPOSTConfirms] | K | yes |
| `C15-NOW` | `surface_confirm.go` | `artifact, decision.String(), now); err != nil { // C15` → `artifact, decision.String(), now+1); err != nil { // C15` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms]; suite-nostatic=KILLED[TestConfirmPOSTConfirms] | K | yes |
| `C15-PASSKEY` | `surface_confirm.go` | `if cid, count, ok := signed.PasskeyUse(); ok {` → `if cid, count, ok := signed.PasskeyUse(); false && ok {` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTRefusedAtC15,TestConfirmPOSTRetriesAFailedConsume,…+1]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTRefusedAtC15,TestConfirmPOSTRetriesAFailedConsume,…+1] | K | yes |
| `C15-COUNT` | `surface_confirm.go` | `SignCount: int64(count)}` → `SignCount: int64(count) - 1}` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTRetriesAFailedConsume,TestConfirmPOSTRetriesAFailedWorldWrite]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTRetriesAFailedConsume,TestConfirmPOSTRetriesAFailedWorldWrite] | K | yes |
| `CR-SUPERSEDED` | `surface_confirm.go` | `case errors.Is(err, store.ErrVerbRequestSuperseded):⏎writ…` → `` | `TestConfirmPOSTRefusedAtC15` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusedAtC15,TestConsumeRefusalOperands]; suite-nostatic=KILLED[TestConfirmPOSTRefusedAtC15,TestConsumeRefusalOperands] | K | yes |
| `CR-CONSUMED` | `surface_confirm.go` | `…e.ErrVerbRequestConsumed), errors.Is(err, store.ErrVerbR…` → `…e.ErrVerbRequestUnknown),` | `TestConsumeRefusalOperands` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusedAtC15,TestConsumeRefusalOperands]; suite-nostatic=KILLED[TestConfirmPOSTRefusedAtC15,TestConsumeRefusalOperands] | K | yes |
| `CR-UNKNOWN` | `surface_confirm.go` | `…e.ErrVerbRequestUnknown), errors.Is(err, store.ErrVerbRe…` → `…e.ErrVerbRequestExpired),` | `TestConsumeRefusalOperands` → **KILLED** | suite=KILLED[TestConsumeRefusalOperands]; suite-nostatic=KILLED[TestConsumeRefusalOperands] | K | yes |
| `CR-EXPIRED` | `surface_confirm.go` | `…RequestUnknown), errors.Is(err, store.ErrVerbRequestExpi…` → `…RequestUnknown),` | `TestConsumeRefusalOperands` → **KILLED** | suite=KILLED[TestConsumeRefusalOperands]; suite-nostatic=KILLED[TestConsumeRefusalOperands] | K | yes |
| `CR-CLOCK` | `surface_confirm.go` | `…e.ErrVerbRequestClockBehind), errors.Is(err, store.ErrVe…` → `…e.ErrVerbRequestDigestChanged):` | `TestConsumeRefusalOperands` → **KILLED** | suite=KILLED[TestConsumeRefusalOperands]; suite-nostatic=KILLED[TestConsumeRefusalOperands] | K | yes |
| `CR-DIGEST` | `surface_confirm.go` | `…uestClockBehind), errors.Is(err, store.ErrVerbRequestDig…` → `…uestClockBehind):` | `TestConsumeRefusalOperands` → **KILLED** | suite=KILLED[TestConsumeRefusalOperands]; suite-nostatic=KILLED[TestConsumeRefusalOperands] | K | yes |
| `CR-RACE-STATUS` | `surface_confirm.go` | `"Conflict", "request "+row.RequestID+" changed while it w…` → `"Conflict2", "request "+row.RequestID+" changed while it …` | `TestConfirmPOSTRefusedAtC15` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusedAtC15]; suite-nostatic=KILLED[TestConfirmPOSTRefusedAtC15] | K | yes |
| `CR-APPROVER` | `surface_confirm.go` | `case isApproverRace(err), isCounterRefusal(err):` → `case isCounterRefusal(err):` | `TestConfirmPOSTRefusedAtC15` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusedAtC15,TestConsumeRefusalOperands]; suite-nostatic=KILLED[TestConfirmPOSTRefusedAtC15,TestConsumeRefusalOperands] | K | yes |
| `CR-COUNTER` | `surface_confirm.go` | `case isApproverRace(err), isCounterRefusal(err):` → `case isApproverRace(err):` | `TestConfirmPOSTRefusedAtC15` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusedAtC15,TestConsumeRefusalOperands]; suite-nostatic=KILLED[TestConfirmPOSTRefusedAtC15,TestConsumeRefusalOperands] | K | yes |
| `CR-KEY-STATUS` | `surface_confirm.go` | `http.StatusUnprocessableEntity, "NotApproved", "the key c…` → `http.StatusConflict, "NotApproved", "the key changed sinc…` | `TestConfirmPOSTRefusedAtC15` → **KILLED** | suite=KILLED[TestConfirmPOSTRefusedAtC15,TestConsumeRefusalOperands]; suite-nostatic=KILLED[TestConfirmPOSTRefusedAtC15,TestConsumeRefusalOperands] | K | yes |
| `OK-STATUS` | `surface_confirm.go` | `writeApprovalJSON(w, http.StatusOK, surfacepage.Confirmed(` → `writeApprovalJSON(w, http.StatusCreated, surfacepage.Conf…` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+3]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+3] | K | yes |
| `OK-DECISION` | `surface_confirm.go` | `…d.ApprovedBy(), decision.String())` → `…d.ApprovedBy(), row.DecisionRef)` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTRetriesAFailedConsume,TestConfirmPOSTRetriesAFailedWorldWrite]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTRetriesAFailedConsume,TestConfirmPOSTRetriesAFailedWorldWrite] | K | yes |
| `OK-BY` | `surface_confirm.go` | `…e.Confirmed(id, signed.ApprovedBy(), decision.String())` → `…e.Confirmed(id, name, decision.String())` | `TestConfirmPOSTConfirms` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms]; suite-nostatic=KILLED[TestConfirmPOSTConfirms] | K | yes |
| `CONF-V` | `confirm_model.go` | `return ConfirmedJSON{V: 1, Confirmed: true,` → `return ConfirmedJSON{V: 2, Confirmed: true,` | `TestConfirmedWire` → **KILLED** | suite=KILLED[TestConfirmedWire]; suite-nostatic=KILLED[TestConfirmedWire] | K | yes |
| `CONF-TRUE` | `confirm_model.go` | `return ConfirmedJSON{V: 1, Confirmed: true,` → `return ConfirmedJSON{V: 1, Confirmed: false,` | `TestConfirmedWire` → **KILLED** | suite=KILLED[TestConfirmedWire]; suite-nostatic=KILLED[TestConfirmedWire] | K | yes |
| `CONF-REQUEST` | `confirm_model.go` | `Request: id, ApprovedBy: approvedBy, Decision: decision,` → `Request: approvedBy, ApprovedBy: approvedBy, Decision: de…` | `TestConfirmedWire` → **KILLED** | suite=KILLED[TestConfirmedWire]; suite-nostatic=KILLED[TestConfirmedWire] | K | yes |
| `CONF-NEXT` | `confirm_model.go` | `"the waiting episode finish opens the pull request"` → `"the minter mints it"` | `TestConfirmedWire` → **KILLED** | suite=KILLED[TestConfirmedWire]; suite-nostatic=KILLED[TestConfirmedWire] | K | yes |
| `CONF-TAG` | `confirm_model.go` | ``json:"approvedBy"`⏎Decision string `json:"decision"`⏎Next` → ``json:"approved_by"`⏎Decision string `json:"decision"`⏎Ne…` | `TestConfirmedWire` → **KILLED** | suite=KILLED[TestConfirmedWire]; suite-nostatic=KILLED[TestConfirmedWire] | K | yes |
| `RTE-REMOVED` | `daemon.go` | `mux.HandleFunc("POST /surface/confirm/{id}", d.approvalRo…` → `` | `TestConfirmRoutesRegistered` → **KILLED** | suite=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+5]; suite-nostatic=KILLED[TestConfirmPOSTConfirms,TestConfirmPOSTGuards,TestConfirmPOSTHoldsHumanMu,…+4] | K | yes |
| `RTE-BARE` | `daemon.go` | `…/confirm/{id}", d.approvalRoute(d.handleSurfaceConfirm))` → `…/confirm/{id}", func(w http.ResponseWriter, r *http.Requ…` | `TestConfirmRoutesRegistered` → **KILLED** | suite=KILLED[TestConfirmPOSTGuards,TestConfirmRoutesRegistered]; suite-nostatic=KILLED[TestConfirmPOSTGuards] | K | yes |
| `RTE-ANY-METHOD` | `daemon.go` | `mux.HandleFunc("POST /surface/confirm/{id}", d.approvalRo…` → `mux.HandleFunc("/surface/confirm/{id}", d.approvalRoute(d…` | `TestConfirmRoutesRegistered` → **KILLED** | suite=KILLED[TestConfirmRoutesRegistered]; suite-nostatic=PASS | R — killed by the route-table test (`TestConfirmRoutesRegistered`), the proof the directive names | yes |
| `RTE-EXTRA` | `daemon.go` | `…SurfaceConfirm))` → `…SurfaceConfirm))⏎mux.HandleFunc("POST /surface/confirmat…` | `TestConfirmRoutesRegistered` → **KILLED** | suite=KILLED[TestConfirmRoutesRegistered]; suite-nostatic=PASS | R — killed by the route-table test (`TestConfirmRoutesRegistered`), the proof the directive names | yes |
| `JS-CHALLENGE` | `confirm.js` | `publicKey: { challenge: hash, rpId` → `publicKey: { challenge: new Uint8Array(32), rpId` | `TestConfirmPageSubmitsAPasskey` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey] | K | yes |
| `JS-RPID` | `confirm.js` | `rpId: group.rp,` → `rpId: location.hostname,` | `TestConfirmPageSubmitsAPasskey` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey] | K | yes |
| `JS-ALLOW` | `confirm.js` | `allowCredentials: group.allow.map((a) => ({ type: 'public…` → `allowCredentials: [],` | `TestConfirmPageSubmitsAPasskey` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey] | K | yes |
| `JS-UV` | `confirm.js` | `userVerification: 'required'` → `userVerification: 'preferred'` | `TestConfirmPageSubmitsAPasskey` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey] | K | yes |
| `JS-TIMEOUT` | `confirm.js` | `const TIMEOUT_MS = 240000;` → `const TIMEOUT_MS = 60000;` | `TestConfirmPageSubmitsAPasskey` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey] | K | yes |
| `JS-ABORT-TIMER` | `confirm.js` | `const timer = setTimeout(() => controller.abort(), TIMEOU…` → `const timer = 0;` | `TestConfirmPageSubmitsAPasskey` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey] | K | yes |
| `JS-CLEAR-TIMER` | `confirm.js` | ` } finally {⏎ clearTimeout(timer);⏎ }⏎ return submit` → ` } finally {⏎ }⏎ return submit` | `TestConfirmPageSubmitsAPasskey` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey] | K | yes |
| `JS-SIGNAL` | `confirm.js` | `timeout: TIMEOUT_MS }, signal: controller.signal });` → `timeout: TIMEOUT_MS } });` | `TestConfirmPageSubmitsAPasskey` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey] | K | yes |
| `JS-GROUP-ORIGIN` | `confirm.js` | `.find((g) => g.origin === location.origin)` → `.find((g) => true)` | `TestConfirmPageSubmitsAPasskey` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey] | K | yes |
| `JS-GO-HIDDEN` | `confirm.js` | `$('go').hidden = !group;` → `$('go').hidden = false;` | `TestConfirmPageSubmitsAPasskey` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey] | K | yes |
| `JS-GO-DISABLED` | `confirm.js` | `$('go').disabled = !on \|\| !group;` → `$('go').disabled = !on;` | `TestConfirmPageSubmitsAPasskey` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey] | K | yes |
| `JS-SIG-DISABLED` | `confirm.js` | `$('post-sig').disabled = !on;` → `$('post-sig').disabled = false;` | `TestConfirmPageSubmitsAnSSHSIG` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey,TestConfirmPageSubmitsAnSSHSIG]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey,TestConfirmPageSubmitsAnSSHSIG] | K | yes |
| `JS-SUBMIT-BUTTONS` | `confirm.js` | ` async function submit(format) {⏎ buttons(false);⏎` → ` async function submit(format) {⏎` | `TestConfirmPageSubmitsAnSSHSIG` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAnSSHSIG]; suite-nostatic=KILLED[TestConfirmPageSubmitsAnSSHSIG] | K | yes |
| `JS-SUBMIT-STATE` | `confirm.js` | ` setState('posting');⏎` → `` | `TestConfirmPageSubmitsAnSSHSIG` → **KILLED** | suite=KILLED[TestConfirmPageShowsRefusals,TestConfirmPageSubmitsAnSSHSIG]; suite-nostatic=KILLED[TestConfirmPageShowsRefusals,TestConfirmPageSubmitsAnSSHSIG] | K | yes |
| `JS-POST-URL` | `confirm.js` | `net('/surface/confirm/' + id, { method: 'POST'` → `net('/surface/approve/' + id, { method: 'POST'` | `TestConfirmPageSubmitsAnSSHSIG` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey,TestConfirmPageSubmitsAnSSHSIG]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey,TestConfirmPageSubmitsAnSSHSIG] | K | yes |
| `JS-POST-METHOD` | `confirm.js` | `{ method: 'POST', headers:` → `{ method: 'PUT', headers:` | `TestConfirmPageSubmitsAnSSHSIG` → **KILLED** | suite=KILLED[TestConfirmPageShowsRefusals,TestConfirmPageSubmitsAPasskey,TestConfirmPageSubmitsAnSSHSIG]; suite-nostatic=KILLED[TestConfirmPageShowsRefusals,TestConfirmPageSubmitsAPasskey,TestConfirmPageSubmitsAnSSHSIG] | K | yes |
| `JS-POST-CT` | `confirm.js` | `headers: { 'Content-Type': 'application/json' },⏎ body:` → `headers: {},⏎ body:` | `TestConfirmPageSubmitsAPasskey` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey] | K | yes |
| `JS-POST-DIGEST` | `confirm.js` | `Object.assign({ v: 1, digest: digestHex }, format)` → `Object.assign({ v: 1, digest: '' }, format)` | `TestConfirmPageSubmitsAnSSHSIG` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey,TestConfirmPageSubmitsAnSSHSIG]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey,TestConfirmPageSubmitsAnSSHSIG] | K | yes |
| `JS-POST-V` | `confirm.js` | `Object.assign({ v: 1, digest: digestHex }, format)` → `Object.assign({ v: 2, digest: digestHex }, format)` | `TestConfirmPageSubmitsAnSSHSIG` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey,TestConfirmPageSubmitsAnSSHSIG]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey,TestConfirmPageSubmitsAnSSHSIG] | K | yes |
| `JS-NET-INIT` | `confirm.js` | `headers: { Accept: 'application/json' } }, init));` → `headers: { Accept: 'application/json' } }));` | `TestConfirmPageSubmitsAnSSHSIG` → **KILLED** | suite=KILLED[TestConfirmPageShowsRefusals,TestConfirmPageSubmitsAPasskey,TestConfirmPageSubmitsAnSSHSIG]; suite-nostatic=KILLED[TestConfirmPageShowsRefusals,TestConfirmPageSubmitsAPasskey,TestConfirmPageSubmitsAnSSHSIG] | K | yes |
| `JS-SSHSIG-FIELD` | `confirm.js` | `submit({ sshsig: $('sig').value })` → `submit({ sshsig: '' })` | `TestConfirmPageSubmitsAnSSHSIG` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAnSSHSIG]; suite-nostatic=KILLED[TestConfirmPageSubmitsAnSSHSIG] | K | yes |
| `JS-WEBAUTHN-FIELD` | `confirm.js` | `return submit({ webauthn: assertionJSON(cred) });` → `return submit({ sshsig: JSON.stringify(assertionJSON(cred…` | `TestConfirmPageSubmitsAPasskey` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey] | K | yes |
| `JS-ASSERT-USERHANDLE` | `confirm.js` | `userHandle: c.response.userHandle ? b64u(c.response.userH…` → `userHandle: ''` | `TestConfirmPageSubmitsAPasskey` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey] | K | yes |
| `JS-ASSERT-CLIENTDATA` | `confirm.js` | `clientDataJSON: b64u(c.response.clientDataJSON)` → `clientDataJSON: b64u(c.response.authenticatorData)` | `TestConfirmPageSubmitsAPasskey` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey] | K | yes |
| `JS-SUCCESS-CONFIRMED` | `confirm.js` | `res.status === 200 && b.confirmed === true` → `res.status === 200` | `TestConfirmPageShowsRefusals` → **KILLED** | suite=KILLED[TestConfirmPageShowsRefusals]; suite-nostatic=KILLED[TestConfirmPageShowsRefusals] | K | yes |
| `JS-SUCCESS-ANY` | `confirm.js` | `res.status === 200 && b.confirmed === true` → `true` | `TestConfirmPageShowsRefusals` → **KILLED** | suite=KILLED[TestConfirmPageShowsRefusals]; suite-nostatic=KILLED[TestConfirmPageShowsRefusals] | K | yes |
| `JS-REFUSAL-TEXT` | `confirm.js` | ``Not confirmed (HTTP ${res.status} ${b.error \|\| ''}): $…` → ``Confirmed (HTTP ${res.status})`` | `TestConfirmPageShowsRefusals` → **KILLED** | suite=KILLED[TestConfirmPageShowsRefusals]; suite-nostatic=KILLED[TestConfirmPageShowsRefusals] | K | yes |
| `JS-REREAD` | `confirm.js` | `…n'}`;⏎ return reread();` → `…n'}`;⏎ return;` | `TestConfirmPageShowsRefusals` → **KILLED** | suite=KILLED[TestConfirmPageShowsRefusals,TestConfirmPageSubmitsAPasskey,TestConfirmPageSubmitsAnSSHSIG]; suite-nostatic=KILLED[TestConfirmPageShowsRefusals,TestConfirmPageSubmitsAPasskey,TestConfirmPageSubmitsAnSSHSIG] | K | yes |
| `JS-REREAD-LOST` | `confirm.js` | `Reading the request again.`;⏎ return reread();` → `Reading the request again.`;⏎ return;` | `TestConfirmPageShowsRefusals` → **KILLED** | suite=KILLED[TestConfirmPageShowsRefusals]; suite-nostatic=KILLED[TestConfirmPageShowsRefusals] | K | yes |
| `JS-REREAD-HIDE` | `confirm.js` | `const reread = () => { $('controls').hidden = true; retur…` → `const reread = () => load();` | `TestConfirmPageSubmitsAPasskey` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey] | K | yes |
| `JS-PASSKEY-FAIL-REREAD` | `confirm.js` | `Nothing was sent.`;⏎ return reread();` → `Nothing was sent.`;⏎ return;` | `TestConfirmPageSubmitsAPasskey` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey] | K | yes |
| `JS-PASSKEY-FAIL-SENDS` | `confirm.js` | `Nothing was sent.`;⏎ return reread();` → `Nothing was sent.`;⏎ return submit({ webauthn: {} });` | `TestConfirmPageSubmitsAPasskey` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey] | K | yes |
| `JS-PASSKEY-STATE` | `confirm.js` | ` setState('waiting');⏎` → `` | `TestConfirmPageSubmitsAPasskey` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey] | K | yes |
| `JS-GO-WIRED` | `confirm.js` | `$('go').addEventListener('click', () => { confirmWithPass…` → `` | `TestConfirmPageSubmitsAPasskey` → **KILLED** | suite=KILLED[TestConfirmPageSubmitsAPasskey]; suite-nostatic=KILLED[TestConfirmPageSubmitsAPasskey] | K | yes |
| `J-LIST-EXPIRY-BOUND` | `verb_requests.go` | `WHERE consumed_at IS NULL AND expires_at > ? ORDER BY` → `WHERE consumed_at IS NULL AND expires_at >= ? ORDER BY` | `TestVerbRequestExpiryBoundary` → **KILLED** | suite=KILLED[TestVerbRequestExpiryBoundary]; suite-nostatic=KILLED[TestVerbRequestExpiryBoundary] | K | yes |
| `J-PUT-COUNT-BOUND` | `verb_requests.go` | `…AND expires_at > ?` → `…AND expires_at >= ?` | `TestVerbRequestExpiryBoundary` → **KILLED** | suite=KILLED[TestVerbRequestExpiryBoundary]; suite-nostatic=KILLED[TestVerbRequestExpiryBoundary] | K | yes |
| `J-PRUNE-BOUND` | `verb_requests.go` | `…AND expires_at <= ?`, now-verbRequestPruneAfter` → `…AND expires_at < ?`, now-verbRequestPruneAfter` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |
| `J-CONSUME-CLOCK` | `verb_requests.go` | `case now < row.CreatedAt:⏎return fmt.Errorf("%w: %s", Err…` → `case now <= row.CreatedAt:⏎return fmt.Errorf("%w: %s", Er…` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `J-CONSUME-EXPIRY` | `verb_requests.go` | `case now >= row.ExpiresAt:⏎return fmt.Errorf("%w: %s", Er…` → `case now > row.ExpiresAt:⏎return fmt.Errorf("%w: %s", Err…` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `J-F2-LOWER` | `finish_confirm.go` | `f.RequestExpires < scope.ExpiresAt &&` → `f.RequestExpires <= scope.ExpiresAt &&` | `TestFinishConfirmRefusal` → **KILLED** | suite=KILLED[TestFinishConfirmRefusal]; suite-nostatic=KILLED[TestFinishConfirmRefusal] | K | yes |
| `J-F2-UPPER` | `finish_confirm.go` | `scope.ExpiresAt <= f.Now+3600` → `scope.ExpiresAt < f.Now+3600` | `TestFinishConfirmRefusal` → **KILLED** | suite=KILLED[TestFinishConfirmRefusal]; suite-nostatic=KILLED[TestFinishConfirmRefusal] | K | yes |
| `J-ROW-EXPIRY` | `verb_requests.go` | `case row.ExpiresAt <= row.CreatedAt:` → `case row.ExpiresAt < row.CreatedAt:` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |

---

## 5. Files (final `git diff --name-status 430668f8 HEAD`)

```
M  host/daemon/daemon.go
M  host/daemon/surface_approve_post_guard_test.go       (pre-existing: the fourth writer site, AM8)
A  host/daemon/surface_confirm.go
M  host/daemon/surface_confirm_read_test.go             (pre-existing: the route table gains the POST)
A  host/daemon/surface_confirm_test.go
A  host/daemon/surface_confirm_write_test.go
A  host/store/verb_requests_bound_test.go
M  host/surfacepage/confirm/confirm.js
M  host/surfacepage/confirm/index.html
A  host/surfacepage/confirm_confirmed_test.go
M  host/surfacepage/confirm_model.go
M  host/surfacepage/confirm_script_test.go              (pre-existing: the harness's events, three sections)
```

Plus this plan and its JSON under `design_docs/planned/` (the controller's commit). **Production files**:
`host/daemon/{surface_confirm.go,daemon.go}`, `host/surfacepage/{confirm_model.go,confirm/confirm.js,confirm/index.html}`
— no `cmd/` file, no `host/store`/`host/broker`/`host/approval`/`host/authority` production byte, no `.ail`.

---

## 6. Sprint gates (controller, after the last sub-milestone's commit; each measured at the base, §1.1)

One bash shell, never under `nohup`: `bash gates.sh <worktree> $HOME/.ailang/state/world-iter403/gates final <HEAD sha>`
(extract it as `accept.sh` is extracted, with the fence `` ````bash gates.sh ``). Every leg is `tee`'d to
`$HOME/.ailang/state/world-iter403/gates/<leg>.txt`; `ALL_DONE.<sha>` is written last with `FAILS: …`. The executor sets
`GATES_SKIP_VERIFY_GO=1`; the controller runs `verify_go` (`VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang
./scripts/verify_go.sh`) after the commit. Legs: the **touched-package derivation** (printed: `git diff --name-only
430668f8 HEAD -- cmd host` → the `.go` files' directories; the plan's list is `host/daemon host/store host/surfacepage`; a
difference FAILS), `status`, `verify_ail`, `gofmt`, `build`, `vet` (`./...`), `fence` (`go test -run '^$' ./...`),
`pkgs` (the three touched packages whole), **`race`** (`host/store` and `host/surfacepage` whole; `host/daemon` on the
approval-surface tests — the whole package under `-race` exceeds the leg's budget at load; CI's race job runs it whole),
**`race3`** (`-race -count=3` over the eight new daemon tests and the boundary pin: SQLite locking and `humanMu`),
**`node`** (S3's nine named harnesses + CP8–CP10, `WORLD_EXEC_NODE`, `--- PASS` per name), **`inv`** (the
inventories by name with `-v`: route tables incl. `TestConfirmRoutesRegistered`, `TestFrozenV1RouteTableMatchesMux`,
`TestDaemonCredentialWriterSites`, `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist`,
`TestProductionContextRoots`, `TestProductionGoSurface`, `TestNoNewDeadlineFreeStoreReads`, `TestEveryVerbHasHelp`, every
credentials-schema test, the A2/A3/A4 page and route tests, the confirm page/asset/CSP tests, S1b's judges, the broker's
finish tests, A4's `TestDecideRefusesReservedScopesByKey`, and the new CX tests), **`SURFACE_POSTS`** (the `/surface`
POST registrations counted from `daemon.go`: base 2 (A3a's approve, A4's decide), final **3** — the one new route), **`chrome`** (S4 changes the
card's page bytes: V2b's fourteen + XC1 + the two `host/daemon` drills), then `verify_go` unless skipped.

````bash gates.sh
#!/bin/bash
# Row 243 V3 S4 gates (plan §6). usage: gates.sh <tree> <logdir> <base|final> [<sha>]
# The controller: bash gates.sh <worktree> $HOME/.ailang/state/world-iter403/gates final <HEAD sha>
# macOS /bin/bash 3.2: no declare -A. Every leg records its own rc; every leg runs (no set -e); every leg is tee'd
# to <logdir>/<leg>.txt. GATES_SKIP_VERIFY_GO=1 skips verify_go (the executor sets it; the controller runs it).
# Never start this under nohup. Writes <logdir>/ALL_DONE.<sha> with FAILS: when every leg has finished.
T=$1; L=$2; MODE=$3
mkdir -p "$L"; cd "$T" || exit 2
SHA=${4:-$(git rev-parse HEAD 2>/dev/null || echo unknown)}
B=430668f86d4bdfb226f92ec39de99440abbbb397
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
# The touched Go packages, derived from the diff (printed so the controller can diff it). The plan's list is
# host/daemon host/store host/surfacepage; a difference FAILS.
TOUCHED=$(git diff --name-only "$B" HEAD -- cmd host 2>/dev/null | grep '\.go$' | xargs -n1 dirname 2>/dev/null | sort -u | tr '\n' ' ')
echo "TOUCHED_PACKAGES: $TOUCHED" | tee -a "$L/legs.txt"
if [ "$MODE" = final ] && [ "$TOUCHED" != "host/daemon host/store host/surfacepage " ]; then FAILS="$FAILS touched_packages_differ"; fi
leg status 'git status --porcelain | wc -l'
leg verify_ail './scripts/verify_ail.sh'
leg gofmt 'out=$(gofmt -l cmd host); echo "$out"; test -z "$out"'
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./... > '$L/fence.raw' 2>&1; rc=\$?; cat '$L/fence.raw'; if grep -E '^(FAIL|--- FAIL)' '$L/fence.raw' >/dev/null; then exit 1; fi; exit \$rc"
leg pkgs "go test -count=1 -timeout 30m ./host/store/ ./host/surfacepage/ ./host/daemon/"
# -race on every touched package: host/store and host/surfacepage whole; host/daemon on the approval-surface
# tests (its whole package under -race exceeds the leg's budget at load; CI's race job runs it whole)
DRACE='^(TestConfirm|TestConsumeRefusal|TestRequestsJSON|TestRequestCardJSON|TestApprovePOST|TestApprovalRoutes|TestDecidePOST|TestDecideConcurrent|TestDecideRefuses|TestDecisionRoutes|TestDecidePage|TestSurfaceRequests|TestDaemonCredentialWriterSites|TestIsApproverRace|TestIsCounterRefusal)'
leg race "go test -race -count=1 ./host/store/ ./host/surfacepage/; r1=\$?; go test -race -count=1 -v ./host/daemon/ -run '$DRACE' > '$L/race_daemon.raw' 2>&1; r2=\$?; grep -E '^(--- |ok|FAIL)' '$L/race_daemon.raw'; exit \$(( r1 | r2 ))"
# SQLite busy/lock and humanMu under -race -count=3 (directive 8, 12): the new daemon tests and the boundary pin
NEW_CX="TestConfirmPOSTGuards TestConfirmPOSTRefusesBeforeAnyWrite TestConfirmPOSTConfirms TestConfirmPOSTHoldsHumanMu TestConfirmPOSTRetriesAFailedWorldWrite TestConfirmPOSTRetriesAFailedConsume TestConfirmPOSTRefusedAtC15 TestConsumeRefusalOperands"
if [ "$MODE" = final ]; then leg race3 "go test -race -count=3 ./host/daemon/ -run '^($(bar $NEW_CX))\$'; r1=\$?; go test -race -count=3 ./host/store/ -run '^TestVerbRequestExpiryBoundary\$'; exit \$(( r1 | \$? ))"; fi
# node: every page harness, non-skippable (WORLD_EXEC_NODE is set above; a harness that cannot run FAILS)
NODE="TestSurfaceAppPure TestSurfaceScriptPure TestSurfaceAppVerbs TestNeedsYouListsDecisions TestAddressHintsArePinned TestConfirmPageFailsClosed TestConfirmPageRendersText TestConfirmPageMessage TestConfirmPageStates"
NEW_NODE="TestConfirmPageSubmitsAPasskey TestConfirmPageShowsRefusals TestConfirmPageSubmitsAnSSHSIG"
[ "$MODE" = base ] && NEW_NODE=""
leg node "set -o pipefail; go test -v -count=1 ./host/surfacepage/ -run '^($(bar $NODE)$(alt $NEW_NODE))\$'"
need node $NODE $NEW_NODE
if grep -E '^\s*--- (SKIP|FAIL)' "$L/node.txt"; then FAILS="$FAILS node_skip_or_fail"; fi
NEW_D="$NEW_CX"; NEW_SP="TestConfirmedWire"; NEW_ST="TestVerbRequestExpiryBoundary"
if [ "$MODE" = base ]; then NEW_D=""; NEW_SP=""; NEW_ST=""; fi
INV_DAEMON="TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestConfirmRoutesRegistered TestDaemonCredentialWriterSites TestSurfaceVerbRouteRegistered TestSurfaceVerbIsReadOnly TestDecidePageAssets TestDecidePageFieldTable TestApprovePageRoutes TestPasskeyPageRoutes TestConfirmPageAssetsServed TestConfirmationCardServed TestRequestsJSONConfirmations TestRequestsJSON TestRequestCardJSON TestDecideRefusesReservedScopesByKey"
INV_BOUNDARY="TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard TestWorkbenchPackageRemainsTransportFree TestPinfetchStaysOutOfStoreReplayAndPublish"
INV_STORE="TestProductionContextRoots TestProductionGoSurface TestStoreSurfaceTakesContext TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext TestCredentialsSchemaMatchesCanonicalDDL TestCredentialsMethodsRequireADeadline TestCredentialsSchema9VerbRequests TestCredentialsNewerThanNineIsRefused TestPutVerbRequestSupersedes TestPutVerbRequestGrammarCapPrune TestVerbRequestsChecks TestConsumeVerbRequest"
INV_VG="TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestVerifygateTestWritesAreForkLocked"
INV_SP="TestSurfaceAppMatchesWorkingCopy TestSurfaceAppShellAssembly TestSurfaceFontEmbed TestConfirmationCard TestConfirmationCardWire TestConfirmPageFieldTable TestConfirmCardGolden TestConfirmPageAssets"
INV_AUTH="TestJudgeVerbSSH TestJudgeVerbAssertion TestRecheckVerb TestJudgeDecisionSSH TestJudgeDecisionAssertion"
INV_BROKER="TestFinishPromptText TestFinishApprovalScopeStrict TestFinishConfirmRefusal"
INV_CMD="TestEveryVerbHasHelp TestVerbTableMatchesTheCLISwitch TestFinishPromptIsTheGolden"
leg inv "set -o pipefail
go test -v -count=1 ./host/daemon/ -run '^($(bar $INV_DAEMON)$(alt $NEW_D))\$' &&
go test -v -count=1 ./host/boundary/ -run '^($(bar $INV_BOUNDARY))\$' &&
go test -v -count=1 ./host/store/ -run '^($(bar $INV_STORE)$(alt $NEW_ST))\$' &&
go test -v -count=1 ./host/verifygate/ -run '^($(bar $INV_VG))\$' &&
go test -v -count=1 ./host/surfacepage/ -run '^($(bar $INV_SP)$(alt $NEW_SP))\$' &&
go test -v -count=1 ./host/authority/ -run '^($(bar $INV_AUTH))\$' &&
go test -v -count=1 ./host/broker/ -run '^($(bar $INV_BROKER))\$' &&
go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(bar $INV_CMD))\$'"
need inv $INV_DAEMON $NEW_D $INV_BOUNDARY $INV_STORE $NEW_ST $INV_VG $INV_SP $NEW_SP $INV_AUTH $INV_BROKER $INV_CMD
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.txt"; then FAILS="$FAILS inv_skip_or_fail"; fi
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.txt")" | tee -a "$L/legs.txt"
# The POST is the only new route: the /surface POST registrations, counted from daemon.go (base 2, final 3)
POSTS=$(grep -cE 'mux\.HandleFunc\("POST /surface/' host/daemon/daemon.go)
echo "SURFACE_POSTS: $POSTS $(grep -oE 'mux\.HandleFunc\("POST /surface/[^"]*"' host/daemon/daemon.go | tr '\n' ' ')" | tee -a "$L/legs.txt"
if [ "$MODE" = final ] && [ "$POSTS" != 3 ]; then FAILS="$FAILS surface_posts_not_3"; fi
# live Chrome: S4 changes the card's page bytes (the submit half); V2b's fourteen, XC1 and the two drills stay green
XB="TestApprovePageInChrome TestApprovePageRefusesASwappedRequest TestApprovePageSSHSIGInChrome TestApprovePageShowsAHostileValueAsText TestPasskeyPageAbortsAPendingCeremony TestPasskeyPageInChrome TestApproveDecideWebAuthnLinkInChrome TestDecidePageApproveWithPasskey TestDecidePageFailsClosed TestDecidePageSSHSIGPaste TestNeedsYouListsPendingRequests TestNeedsYouNotifiesOnlyWhenAskedOnce TestVerbBarEpisodeReportInChrome TestVerbBarWhyInChrome TestConfirmPageShowsTheFinishPromptInChrome"
XD="TestSurfaceLiveChrome TestSurfaceFollowChrome"
CH="'/Applications/Google Chrome.app/Contents/MacOS/Google Chrome'"
leg chrome "set -o pipefail; WORLD_CHROME=$CH WORLD_REQUIRE_PASSKEY_BROWSER=1 WORLD_REQUIRE_SSH_KEYGEN=1 go test ./cmd/ailang-worldd/ -run '^($(bar $XB))\$' -count=1 -v; r1=\$?; WORLD_CHROME=$CH go test ./host/daemon/ -run '^($(bar $XD))\$' -count=1 -v; r2=\$?; exit \$(( r1 | r2 ))"
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
  only writer** of the tree; no hand edit, no `gofmt -w`, no fuzz. A failed `--dry-run` or a red acceptance → **stop and
  report** with the log; do not repair.
- **Never `nohup`, never background a test leg; pass a 600000 ms timeout to every long Bash call.** Each acceptance
  part is one call (C14 measured each).
- **Runs**: **Run 1 = M10a → M10b → M10c → M10d**, **Run 2 = M11** (each patch assumes the previous tree; the controller
  may commit between runs). Per sub-milestone the parts, in order: **M10a** `fast`, `daemon`; **M10b** `fast`; **M10c**
  `fast`; **M10d** `fast`, `daemon`; **M11** `fast`, `chrome`. `accept.sh` exports `AILANG_BIN=$HOME/.pinned-ailang/ailang`
  and `WORLD_EXEC_NODE=/opt/homebrew/bin/node` and unsets `WORLD_CHROME`, `WORLD_REQUIRE_PASSKEY_BROWSER` and every
  golden-update variable; **M11's `chrome` part sets `WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google
  Chrome" WORLD_REQUIRE_PASSKEY_BROWSER=1` itself**.
- **Commit-time test commands** (`commit_ms.sh <M> <msgfile> <test-cmd…>`; each carries its env through `accept.sh`, and
  each fits one 600 s call — the `daemon` parts ran at acceptance and are not repeated at commit):
  `bash $T/accept.sh $W <M> fast` for **M10a, M10b, M10c, M10d**, and `bash $T/accept.sh $W M11 fast && bash
  $T/accept.sh $W M11 chrome` for **M11** (the `chrome` part sets `WORLD_CHROME` and `WORLD_REQUIRE_PASSKEY_BROWSER=1`), with
  `T=$HOME/.ailang/state/world-iter403/exec-patches W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter403-row243-v3s4`
  (`commit_ms.sh` runs the command in the worktree after copying the snapshot; `accept.sh` takes the tree as its first
  argument).
- After each sub-milestone's acceptance, snapshot **cumulatively against the base** to
  `$HOME/.ailang/state/world-iter403/snap/<M>/`: every file added or changed since `430668f8` (repo-relative paths), a
  **`DELETED`** file (one repo path per line; **empty for every S4 stage**), a **`STATUS`** file, and `MANIFEST.sha256`
  (the two plan files, `STATUS` and `DELETED` are not in it). **Never an `rsync` of whole directories.** The `find` takes
  only the top-level dirs that exist (every S4 stage has `host/` only — no `cmd/`, no `docs/`), rehearsed on all five
  (C15):
  ```bash
  bash <<'EOF2'
  set -euo pipefail
  M=M10a   # the sub-milestone just accepted
  B=430668f86d4bdfb226f92ec39de99440abbbb397
  W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter403-row243-v3s4
  D="$HOME/.ailang/state/world-iter403/snap/$M"; rm -rf "$D"; mkdir -p "$D"
  cd "$W"
  { git diff --name-status --no-renames "$B" -- cmd host docs; git ls-files --others --exclude-standard -- cmd host docs | sed 's/^/A\t/'; } | sort -u -t "$(printf '\t')" -k2 > "$D/STATUS"
  awk -F'\t' '$1 != "D" {print $2}' "$D/STATUS" | while read -r p; do mkdir -p "$D/$(dirname "$p")"; cp -p "$p" "$D/$p"; done
  awk -F'\t' '$1 == "D" {print $2}' "$D/STATUS" > "$D/DELETED"
  (cd "$D" && find $(for d in cmd host docs; do [ -d "$d" ] && echo "$d"; done) -type f -print0 | sort -z | xargs -0 shasum -a 256) > "$D/MANIFEST.sha256"
  cat "$D/STATUS"; echo "DELETED: $(wc -l < "$D/DELETED")"; echo "MANIFEST: $(wc -l < "$D/MANIFEST.sha256")"
  EOF2
  ```
  Expected (`STATUS` lines / `DELETED` / MANIFEST): **M10a** 5 / 0 / 5; **M10b** 6 / 0 / 6; **M10c** 7 / 0 / 7;
  **M10d** 9 / 0 / 9; **M11** 12 / 0 / 12. The controller copies each snapshot over the worktree, commits it as one
  commit, and checks the numstat (`git diff -M --numstat 430668f8 HEAD -- cmd host docs | awk '{s+=$1+$2} END {print
  s}'`): **M10a 433, M10b 605, M10c 917, M10d 916, M11 1,097**. A different number → stop.
- **Stop conditions**: a dry-run failure; any `--- FAIL`, `--- SKIP` or `FAIL` line in an acceptance (`ACCEPT_RED`);
  gofmt output; a numstat mismatch; any file outside §5 changed (`git status --short` after each patch must equal the
  cumulative list of §3's paragraphs); a numstat above 1,400.
- **Known load-sensitive tests** (the directive's list) are not this sprint's: a red there is re-run once alone
  (`TestWorkspaceExecConfiguredBindsTheEpisodeHandler`, `TestMarketingCaptureInventoryMutation` in the `daemon` part;
  `TestConvertedStoreMethodsM5a`/`M5b` in `host/store` — S3 measured both red under load at its base, green alone).
- D-WORLD-109: nothing touches `~/.ailang/world`, `~/.ssh`, a real `/dev/tty` or a real key; every store is a
  `t.TempDir()` file.

---

## 8. Risks

- **R1 — a POST that writes.** S4 adds the daemon's fourth credentials-writer site and the first page-reachable
  `Human.Approve` decided by a key for a finish scope (AM7, D-WORLD-118 = A). Its writes happen only after C0–C13 hold;
  CX0–CX13 compare both stores before and after every refusal. L2 still refuses a key-decided approval at publish
  (T12, unchanged).
- **R2 — the seams are production variables** (S4-D7). A test that forgets to restore one would leak into later tests
  in the package; every assignment is paired with a `t.Cleanup` restore, and no daemon test runs in parallel.
- **R3 — two stores, not atomic (RV2).** Each failed attempt after C14 leaves one inert approval no row names; CX21 and
  audit note 3's row pin that the row never names an orphan.
- **R4 — the page now POSTs.** `confirm.js` sends only `ParseApproveBody`'s grammar (no `typed`, no `y`; T6), with
  `credentials: 'omit'` (no cookie for C0 to refuse) and under the unchanged CSP's `connect-src 'self'`.
- **R5 — host/daemon under load.** The `daemon` part runs the whole package (≈ 260 s alone, 469 s at load ≈ 20 in C14);
  a red known-flaky test is re-run once alone (§7).

---

## 9. Acceptance criteria (the judge scores against these)

1. **Exactly one route added, a POST**: `TestConfirmRoutesRegistered` `--- PASS` (route table: `GET /surface/confirm.js`,
   `GET /surface/confirm/{id}`, `GET /surface/confirmations/{id}`, `POST /surface/confirm/{id}`, each behind
   `approvalRoute`; PUT/DELETE/OPTIONS on the card's path and POST on the script and data paths → 405 with no CORS
   header); `SURFACE_POSTS: 3` in the gates (base 2); `RTE-*` mutants killed (§4);
   `TestFrozenV1RouteTableMatchesMux`, `TestApprovalRoutesRegistered`, `TestDecisionRoutesRegistered` unedited, green.
2. **The route is registered only in M10d**, after the handler's every step is proven (M10a–M10c test it through
   `approvalRoute` directly).
3. **No CLI change**: no `cmd/` file in §5; `TestEveryVerbHasHelp`, `TestVerbTableMatchesTheCLISwitch`,
   `TestFinishPromptIsTheGolden` green.
4. **No schema change**: no `host/store` production file in §5; every credentials-schema test green unedited.
5. **`surfaceCSP` unchanged**: `TestDecidePageAssets` (pins the literal) green unedited; `TestConfirmPageAssetsServed`
   serves the new page bytes under it.
6. **No `.ail`**: `verify_ail` identical to the base's.
7. **Every refusal before C14 writes nothing**: CX0–CX13 compare `confirmState` (every `verb_requests` row, approvers'
   last use and revocation, passkey counters, world decisions and head) before and after each refusal.
8. **C14 then C15, the writer held across both, under `humanMu`**: CX17 (the world decision decided by the key, requester
   `operator:episode-finish`, the request's scope, the row's five columns, counter and last use), CX19 (the probes and
   H-LOCK), `ORD-C15-BEFORE-C14`, `ORD-C13-AFTER-C14`, `WR-CLOSE-EARLY`, `WR-NO-CLOSE`, `H-LOCK`, `H-UNLOCK-EARLY`,
   `H-LOCK-AFTER-C11` killed.
9. **The window table's every row is driven by a test** (S4-D9), incl. audit note 3's: CX20 (0 decisions, then 1),
   CX21 (1 orphan; same-second retry 1 decision; a second later 2, the row names the second), `superseded` (409 naming
   the episode, the row `superseded`, no row names the orphan).
10. **CX17's 200 body is `surfacepage.Confirmed` of the same inputs** (never hand-typed), and `Confirmed`'s wire is the
    design's literal (`TestConfirmedWire`).
11. **The submit half**: CP8 (challenge = the message hash, the group's rp and credentials, UV required, 240 s, the
    assertion file, the answer, the re-read; no button without a group; an aborted or refused ceremony sends nothing),
    CP9 (refusals never shown as a confirmation; the controls return), CP10 (`{v, digest, sshsig}`; no second submit in
    flight); XC1 green in live Chrome on the new bytes.
12. **Exactly these pre-existing test files change**: `host/daemon/surface_approve_post_guard_test.go` (AM8's fourth
    writer site), `host/daemon/surface_confirm_read_test.go` (the route table gains the POST), and
    `host/surfacepage/confirm_script_test.go` (the stand-in DOM gains `addEventListener`/`click()`, `open` a `navigator`;
    three new sections). **None deleted.** Every other test file in §5 is new.
13. **Judge-402 note 1**: `J-LIST-EXPIRY-BOUND` and `J-PUT-COUNT-BOUND` killed by `TestVerbRequestExpiryBoundary`; every
    other time-comparison operator of the verb store and F2 its own mutant, killed (§4 class J).
14. **Size**: `git diff -M --numstat 430668f8 HEAD -- cmd host docs` sums to **1,097** (≤ 1,400; the same unfiltered);
    each sub-milestone's numstat as §7.
15. **Mutants**: every `K` row killed by its named test and by `suite-nostatic`; `R` rows killed by the route-table
    test; every `CTL` green; the `E` rows named with reasons; the plan's ids are the ledger's.
16. **Gates**: §6's script `FAILS: none` (or only a directive-listed load-sensitive test, green alone), `verify_go` rc 0.

---

## Appendix C — after this sprint (S4's close; what S5 and S6 inherit)

- **Moved out: nothing** (S4-D1: 1,097 ≤ 1,400; the design's M11 → S6 cut measured and not taken). No patch is parked.
- **Residuals carried, unchanged**: S1 Appendix C's note 4 (CF12: the `resume-dispatch` reuse path after a
  surface-confirmed approval whose push failed skips CF6/CF7) and note 5 / gemini r2 (the ≈ 1 s revocation window after
  C15; fail-safe) → **S5**; judge-402 note 2 (a day-expired request of the same episode is superseded, not pruned, by
  the next put) → documented residual (S4-D11); S3 §1.4 F1 (the DDL accepts a superseded row with a `decision_ref`; no
  writer) → a later schema; design RV1–RV9.
- **Closed here**: the CX15→CX17 gap (S4-D3, the design's ids kept); audit note 3's daemon half (S4-D9); judge-402
  note 1 (S4-D10).
- **New for S5** (§1.4 F4): **CX21's CLI half** — the waiting `episode finish --confirm-on-surface` must dispatch with
  the decision the row names (the second, after a C15 failure), never the orphan: CF3's rig over this daemon.
- **What S5 (the CLI) inherits from S4**: a row is `approved` only with `decision_ref` set in the same UPDATE (C15); the
  decision is `RecordOperatorApproval`'s, requester `operator:episode-finish`, `decidedBy` = `row.approved_by`, scope
  = the request's `scope` field byte for byte (CX17) — CF7's read-back compares exactly these; a 200 means the row is
  consumed; a 500 after C14 leaves the row pending and an orphan no row names (CF's wait keeps polling until expiry +
  `verbLandGrace`); the card's links are `/surface/confirm/<id>` (served, POSTable). `POST /surface/confirm/{id}` is
  the only writer of `verb_requests.outcome = 'approved'`.
- **What S6 inherits**: the submit half in place (XC2 drives it with the virtual authenticator: `#go`, then the
  `#answer` line and the card's re-read state `closed`); `requests.json`'s `confirmations` for *Needs you* (CP7).

---

## Appendix P — the patches (apply with §2, in order; generated from the compiled scratch stages; do not edit)

### Patch M10a — the handler (unregistered), `Confirmed`, the fourth writer site, CX0–CX6

````diff M10a
diff --git a/host/daemon/surface_approve_post_guard_test.go b/host/daemon/surface_approve_post_guard_test.go
index 4003f681..87dde271 100644
--- a/host/daemon/surface_approve_post_guard_test.go
+++ b/host/daemon/surface_approve_post_guard_test.go
@@ -126,9 +126,9 @@ func TestApprovePOSTNoSessionNoCookie(t *testing.T) {
 	}
 }
 
-// RT12: the daemon's production code opens the credentials writer at exactly two sites:
-// ensureCredentialsStore's upgrade at serve start and the POST's record. Control: the read-only
-// opens are found by the same walk.
+// RT12: the daemon's production code opens the credentials writer at exactly these sites:
+// ensureCredentialsStore's upgrade at serve start, A3a's record, A4's key arm and V3's
+// confirmation (row 243 AM8). Control: the read-only opens are found by the same walk.
 func TestDaemonCredentialWriterSites(t *testing.T) {
 	files, err := filepath.Glob("*.go")
 	if err != nil {
@@ -165,7 +165,7 @@ func TestDaemonCredentialWriterSites(t *testing.T) {
 			})
 		}
 	}
-	if want := map[string]string{"ensureCredentialsStore": "UpgradeCredentials", "handleSurfaceApprove": "OpenCredentials", "keyArm": "OpenCredentials"}; !reflect.DeepEqual(sites, want) || readOnly < 5 {
+	if want := map[string]string{"ensureCredentialsStore": "UpgradeCredentials", "handleSurfaceApprove": "OpenCredentials", "keyArm": "OpenCredentials", "handleSurfaceConfirm": "OpenCredentials"}; !reflect.DeepEqual(sites, want) || readOnly < 5 {
 		t.Fatalf("writer sites %v (read-only opens %d), want exactly %v and at least 5 read-only opens", sites, readOnly, want)
 	}
 }
diff --git a/host/daemon/surface_confirm.go b/host/daemon/surface_confirm.go
new file mode 100644
index 00000000..24b39fac
--- /dev/null
+++ b/host/daemon/surface_confirm.go
@@ -0,0 +1,215 @@
+package daemon
+
+// Row 243 V3 sprint 4 M10 (design_docs/planned/w-world-address-bar-verbs-v3.md §2.6): POST
+// /surface/confirm/{id}, a verified enrolled key's confirmation of one stored verb request. After R2
+// (approvalRoute) the guards C0-C6 (A4's R0, A3a's R1, R3-R6); then, under humanMu and within
+// commitBudget, C7-C12 on a read-only handle (each refusal writes nothing), C13 a non-waiting writer,
+// C14 the world store's finish approval decided by the key, and C15 one credentials transaction that
+// consumes the row and names the decision. The writer is held across C14 and C15, as
+// handleSurfaceDecide holds keyArm's across DecideApproval and LandKeyDecision. A failure after C14
+// leaves an approval no row names (inert, design W62) and the row pending: the same body retries.
+
+import (
+	"context"
+	"crypto/sha256"
+	"encoding/hex"
+	"errors"
+	"fmt"
+	"io"
+	"net/http"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// confirmRequester is the finish approval's requester, the typed approve phase's (W5).
+const confirmRequester = "operator:episode-finish"
+
+// confirmIncomplete is the 500 of a failed C14 or C15: nothing the CLI uses was written.
+const confirmIncomplete = "the confirmation did not complete; approve again"
+
+// Test seams: C11's judges (C12's refusal needs a verified confirmation of other bytes) and C15's
+// consume (a request superseded, or its approver changed, between C14 and C15).
+var (
+	judgeVerbSSH       = authority.JudgeVerbSSH
+	judgeVerbAssertion = authority.JudgeVerbAssertion
+	consumeVerbRequest = (*store.Credentials).ConsumeVerbRequest
+)
+
+// confirmAdmission is C8: approval.Admit on row at now, "" when admitted.
+func confirmAdmission(row store.VerbRequestRow, now int64) string {
+	stored, _ := hex.DecodeString(row.Digest) // the CHECK holds it to 64 characters; not hex: no bytes hash to it
+	var digest [32]byte
+	copy(digest[:], stored)
+	switch approval.Admit(approval.AdmitInput{Now: now, CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt, Consumed: row.Outcome != "",
+		StoredDigest: digest, BytesDigest: sha256.Sum256(row.Request)}) {
+	case approval.Admitted:
+		return ""
+	case approval.RefuseConsumed:
+		if row.Outcome == store.VerbOutcomeSuperseded {
+			return "request " + row.RequestID + " was superseded by a newer finish of " + row.EpisodeID
+		}
+		return "request " + row.RequestID + " was already confirmed by " + row.ApprovedBy
+	case approval.RefuseClockBehind:
+		return "the clock is behind request " + row.RequestID + "'s creation time; refusing"
+	case approval.RefuseExpired:
+		return "request " + row.RequestID + " expired; run episode finish again for a new one"
+	}
+	return "request " + row.RequestID + "'s stored bytes do not match its digest; refusing"
+}
+
+// confirmStoreError answers a credentials-store failure: below schema 9 503 Restart, a busy writer
+// 503 Busy with Retry-After (never a wait), anything else the 500 that asks for a retry.
+func (d *Daemon) confirmStoreError(w http.ResponseWriter, r *http.Request, err error) {
+	switch {
+	case isOlderStore(err):
+		writeApprovalRefusal(w, http.StatusServiceUnavailable, "Restart",
+			"the credentials store predates verb requests: restart `ailang-worldd serve` with this build once (it upgrades the store), then approve again")
+	case store.IsWriterAlreadyActive(err):
+		w.Header().Set("Retry-After", "1")
+		writeApprovalRefusal(w, http.StatusServiceUnavailable, "Busy", "the credentials store is busy (a mint or a session verb is writing it); try again")
+	default:
+		d.writeInternalErrorLog(r, err)
+		writeApprovalRefusal(w, http.StatusInternalServerError, "Internal", confirmIncomplete)
+	}
+}
+
+// consumeRefusal is C15's answer: the row raced (409), the approver or the counter changed since
+// C11 (422), else confirmStoreError.
+func (d *Daemon) consumeRefusal(w http.ResponseWriter, r *http.Request, row store.VerbRequestRow, err error) {
+	switch {
+	case errors.Is(err, store.ErrVerbRequestSuperseded):
+		writeApprovalRefusal(w, http.StatusConflict, "Conflict", "request "+row.RequestID+" was superseded by a newer finish of "+row.EpisodeID)
+	case errors.Is(err, store.ErrVerbRequestConsumed), errors.Is(err, store.ErrVerbRequestUnknown), errors.Is(err, store.ErrVerbRequestExpired),
+		errors.Is(err, store.ErrVerbRequestClockBehind), errors.Is(err, store.ErrVerbRequestDigestChanged):
+		writeApprovalRefusal(w, http.StatusConflict, "Conflict", "request "+row.RequestID+" changed while it was confirmed: "+err.Error())
+	case isApproverRace(err), isCounterRefusal(err):
+		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the key changed since it was judged: "+err.Error())
+	default:
+		d.confirmStoreError(w, r, err)
+	}
+}
+
+// handleSurfaceConfirm is POST /surface/confirm/{id} (design §2.6, steps C0-C15 in order).
+func (d *Daemon) handleSurfaceConfirm(w http.ResponseWriter, r *http.Request, o []string) {
+	if len(r.Header.Values("Authorization")) != 0 || len(r.Header.Values("Cookie")) != 0 { // C0
+		writeApprovalRefusal(w, http.StatusForbidden, "Forbidden", "a confirmation is never made with a session or a cookie (the agent holds a session); confirm with your key on the card")
+		return
+	}
+	id := r.PathValue("id")
+	if !store.VerbRequestIDPattern.MatchString(id) { // C1
+		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", confirmIDReason)
+		return
+	}
+	if origin := r.Header.Get("Origin"); !allowedApprovalOrigin(origin, o) { // C3
+		writeApprovalRefusal(w, http.StatusForbidden, "CrossOrigin", crossOriginReason(origin))
+		return
+	}
+	if !sameOriginFetch(r.Header.Get("Sec-Fetch-Site")) { // C4
+		writeApprovalRefusal(w, http.StatusForbidden, "CrossOrigin", "Sec-Fetch-Site is not same-origin")
+		return
+	}
+	if !jsonMediaType(r.Header.Get("Content-Type")) { // C5
+		writeApprovalRefusal(w, http.StatusUnsupportedMediaType, "UnsupportedMediaType", "the body must be application/json")
+		return
+	}
+	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, surfacepage.MaxApproveBodyBytes)) // C6
+	if err != nil {
+		writeApprovalRefusal(w, http.StatusRequestEntityTooLarge, "TooLarge", fmt.Sprintf("the body is over %d bytes", surfacepage.MaxApproveBodyBytes))
+		return
+	}
+	body, err := surfacepage.ParseApproveBody(raw) // no typed arm: a typed or y key is unknown
+	if err != nil {
+		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", err.Error())
+		return
+	}
+	d.humanMu.Lock()
+	defer d.humanMu.Unlock()
+	ctx, cancel := context.WithTimeout(r.Context(), d.commitBudget)
+	defer cancel()
+	now := d.approveNow()
+	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
+	if err != nil {
+		d.approvalInternal(w, r, err)
+		return
+	}
+	defer func() { _ = ro.Close() }()
+	row, ok, err := ro.VerbRequest(ctx, id) // C7
+	switch {
+	case err != nil:
+		d.confirmStoreError(w, r, err)
+		return
+	case !ok:
+		writeApprovalRefusal(w, http.StatusNotFound, "NotFound", "no verb request "+id)
+		return
+	}
+	if reason := confirmAdmission(row, now); reason != "" { // C8
+		writeApprovalRefusal(w, http.StatusConflict, "Conflict", reason)
+		return
+	}
+	if body.Digest != row.Digest { // C9
+		writeApprovalRefusal(w, http.StatusConflict, "Conflict", "the request changed since the card was shown; reload it")
+		return
+	}
+	req, err := approval.DecodeVerb(row.Request) // C10
+	if err == nil {
+		err = req.Validate()
+	}
+	if err == nil {
+		if reason := broker.FinishConfirmRefusal(broker.FinishConfirmFacts{Scope: req.Scope, Title: req.Title, Body: req.Body, Prompt: req.Prompt,
+			Committer: req.Committer, Protected: req.Protected, RequestExpires: row.ExpiresAt, Now: now}); reason != "" {
+			err = errors.New(reason)
+		}
+	}
+	if err != nil {
+		writeApprovalRefusal(w, http.StatusConflict, "NotConfirmable", "request "+id+" cannot be confirmed: "+err.Error()+"; confirm at the terminal")
+		return
+	}
+	digest, challenge := sha256.Sum256(row.Request), [32]byte{} // C8 held: the bytes hash to the digest
+	c, _ := hex.DecodeString(row.Challenge)                     // PutVerbRequest's grammar: 64 lowercase hex
+	copy(challenge[:], c)
+	in := authority.VerbJudgeInput{Digest: digest, Challenge: challenge, BytesDigest: digest, CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt, Now: now}
+	judge, artifact, what := judgeVerbAssertion, []byte(body.WebAuthn), "passkey assertion"
+	if body.SSHSIG != nil {
+		judge, artifact, what = judgeVerbSSH, []byte(*body.SSHSIG), "SSH signature"
+	}
+	signed, err := judge(ctx, ro, in, artifact) // C11
+	if err == nil {
+		err = signed.Binds(digest, challenge, row.ExpiresAt) // C12
+	}
+	var je *authority.JudgeError
+	switch {
+	case errors.As(err, &je) && je.Step == authority.JudgeStepStore:
+		d.confirmStoreError(w, r, err)
+		return
+	case err != nil:
+		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the "+what+" does not confirm request "+id+": "+err.Error())
+		return
+	}
+	wr, err := store.OpenCredentials(d.credentialsPath) // C13: non-waiting, before any write
+	if err != nil {
+		d.confirmStoreError(w, r, err)
+		return
+	}
+	defer func() { _ = wr.Close() }() // held across C14 and C15, as handleSurfaceDecide holds keyArm's
+
+	_, decision, err := broker.RecordOperatorApproval(ctx, d.store, req.Scope, confirmRequester, signed.ApprovedBy(), now) // C14
+	if err != nil {
+		d.writeInternalErrorLog(r, err)
+		writeApprovalRefusal(w, http.StatusInternalServerError, "Internal", confirmIncomplete)
+		return
+	}
+	name, format, key, _ := signed.Use()
+	use := store.ApproverUse{Name: name, Format: format, PublicKey: key}
+	if cid, count, ok := signed.PasskeyUse(); ok {
+		use.Passkey = &store.PasskeyUse{CredentialID: cid, SignCount: int64(count)}
+	}
+	if err := consumeVerbRequest(wr, ctx, id, row.Digest, use, signed.ApprovedBy(), artifact, decision.String(), now); err != nil { // C15
+		d.consumeRefusal(w, r, row, err)
+		return
+	}
+	writeApprovalJSON(w, http.StatusOK, surfacepage.Confirmed(id, signed.ApprovedBy(), decision.String()))
+}
diff --git a/host/daemon/surface_confirm_test.go b/host/daemon/surface_confirm_test.go
new file mode 100644
index 00000000..e63045a3
--- /dev/null
+++ b/host/daemon/surface_confirm_test.go
@@ -0,0 +1,175 @@
+package daemon
+
+// Row 243 V3 sprint 4 M10 (design §2.6, §6 CX0-CX13): POST /surface/confirm/{id}'s guards and its
+// refusals before any write. Requests are S3's putVerb rows (approval.CanonicalVerb of verbFixture);
+// signatures are approvaltest's over approval.VerbMessage of the row; every refusal is compared with
+// the stores before it (nothing written). Stores and keys are tests' (D-WORLD-109).
+
+import (
+	"bytes"
+	"context"
+	"database/sql"
+	"encoding/hex"
+	"fmt"
+	"net"
+	"net/http"
+	"net/http/httptest"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+const cfEp, cfID = "w243", "vr-00000000000000c1"
+
+// newConfirmRig enrols apSigner as mark-ssh and apPasskey as mark (counter 1) at apT0-10 and stores
+// one pending draft request cfID for cfEp, created at apT0-30.
+func newConfirmRig(t *testing.T) (approvalRig, store.VerbRequestRow) {
+	t.Helper()
+	a := newApprovalRig(t)
+	a.enrolSSH("mark-ssh", apSigner)
+	a.enrolPasskey("mark", apPasskey, apOrigin)
+	return a, a.putVerb(cfID, cfEp, apT0-30, true)
+}
+
+// verbMessageOf is row's 86-byte world-verb/v1 message and its hash.
+func verbMessageOf(t *testing.T, row store.VerbRequestRow) ([]byte, [32]byte) {
+	t.Helper()
+	d, _ := hex.DecodeString(row.Digest)
+	c, _ := hex.DecodeString(row.Challenge)
+	m, err := approval.VerbMessage([32]byte(d), [32]byte(c), row.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	h, err := approval.VerbMessageHash([32]byte(d), [32]byte(c), row.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return m, h
+}
+
+// confirmSSH is the card's body for row with s's SSHSIG over its message under namespace.
+func confirmSSH(t *testing.T, s *approvaltest.Signer, namespace string, row store.VerbRequestRow) []byte {
+	t.Helper()
+	m, _ := verbMessageOf(t, row)
+	return sshBody(row.Digest, s.Sign(namespace, m))
+}
+
+// confirmPasskey is the card's body for row with pk's assertion (counter n) over its message hash.
+func confirmPasskey(t *testing.T, pk *approvaltest.Passkey, row store.VerbRequestRow, n uint32) []byte {
+	t.Helper()
+	_, h := verbMessageOf(t, row)
+	return approveBody(row.Digest, pk.Assert(apOrigin, "localhost", h, n))
+}
+
+// confirmAt POSTs body to /surface/confirm/<id> with Host host, as the card would (loopback Origin,
+// same-origin, JSON); hdr pairs override those headers. Until M10d registers the route, the request
+// reaches the handler through approvalRoute, as the mux will route it.
+func (a approvalRig) confirmAt(host, id string, body []byte, hdr ...string) *httptest.ResponseRecorder {
+	a.t.Helper()
+	req := httptest.NewRequest("POST", "/surface/confirm/"+id, bytes.NewReader(body))
+	req = req.WithContext(context.WithValue(req.Context(), http.LocalAddrContextKey, net.Addr(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 7644})))
+	req.Host = host
+	req.SetPathValue("id", id)
+	for h, i := append([]string{"Origin", apOrigin, "Sec-Fetch-Site", "same-origin", "Content-Type", "application/json"}, hdr...), 0; i+1 < len(h); i += 2 {
+		req.Header.Set(h[i], h[i+1])
+	}
+	rec := httptest.NewRecorder()
+	a.d.approvalRoute(a.d.handleSurfaceConfirm).ServeHTTP(rec, req)
+	return rec
+}
+
+// confirm is confirmAt on the loopback Host.
+func (a approvalRig) confirm(id string, body []byte, hdr ...string) *httptest.ResponseRecorder {
+	a.t.Helper()
+	return a.confirmAt(apLoopback, id, body, hdr...)
+}
+
+// confirmState is the "nothing written" witness: every verb_requests row, each approver's last use
+// and revocation, each passkey's counter, and the world store's approval decisions and head.
+func (a approvalRig) confirmState() string {
+	a.t.Helper()
+	n, head := decideRig{a}.decisions()
+	return fmt.Sprintf("%s decisions=%d head=%s", a.credState(), n, head)
+}
+
+// credState is confirmState's credentials half.
+func (a approvalRig) credState() string {
+	a.t.Helper()
+	db, err := sql.Open("sqlite", a.d.credentialsPath)
+	if err != nil {
+		a.t.Fatal(err)
+	}
+	defer func() { _ = db.Close() }()
+	var b strings.Builder
+	for _, q := range []string{`SELECT request_id, coalesce(consumed_at, 0), coalesce(outcome, ''), coalesce(approved_by, ''), hex(artifact), coalesce(decision_ref, '') FROM verb_requests ORDER BY request_id`,
+		`SELECT name, coalesce(last_used_at, 0), coalesce(revoked_at, 0), '', '', '' FROM approvers ORDER BY name`, `SELECT name, sign_count, '', '', '', '' FROM webauthn_credentials ORDER BY name`} {
+		rows, err := db.Query(q)
+		if err != nil { // a schema-8 store has no verb_requests
+			fmt.Fprintf(&b, "%v;", err)
+			continue
+		}
+		for rows.Next() {
+			var c [6]string
+			if err := rows.Scan(&c[0], &c[1], &c[2], &c[3], &c[4], &c[5]); err != nil {
+				a.t.Fatal(err)
+			}
+			fmt.Fprintf(&b, "%v;", c)
+		}
+		_ = rows.Close()
+	}
+	return b.String()
+}
+
+// CX0-CX6: C0 (a session or a cookie), C1 (the id), R2 (the Host, approvalRoute), C3-C5 (A3a's
+// guards), C6 (size, ParseApproveBody: no typed arm) each refuse with nothing written; the same
+// signature then confirms (control).
+func TestConfirmPOSTGuards(t *testing.T) {
+	a, row := newConfirmRig(t)
+	m, _ := verbMessageOf(t, row)
+	sig := string(apSigner.Sign(approval.VerbNamespace, m))
+	body := sshBody(row.Digest, []byte(sig))
+	before := a.confirmState()
+	for _, tc := range []struct {
+		name, id      string
+		body          []byte
+		hdr           []string
+		status        int
+		class, reason string
+	}{
+		{"C0 authorization", cfID, body, []string{"Authorization", "Bearer x"}, 403, "Forbidden", "never made with a session"},
+		{"C0 cookie", cfID, body, []string{"Cookie", "s=1"}, 403, "Forbidden", "never made with a session"},
+		{"C0 before C1", "vr-0123", body, []string{"Cookie", "s=1"}, 403, "Forbidden", ""},
+		{"C1 short", "vr-0123", body, nil, 400, "BadRequest", "vr- and 16 lowercase hex"},
+		{"C1 upper", "vr-00000000000000C1", body, nil, 400, "BadRequest", ""},
+		{"C1 mint id", "mr-00000000000000c1", body, nil, 400, "BadRequest", ""},
+		{"C1 before C3", "vr-0123", body, []string{"Origin", "http://evil.example"}, 400, "BadRequest", ""},
+		{"C3 missing", cfID, body, []string{"Origin", ""}, 403, "CrossOrigin", ""},
+		{"C3 null", cfID, body, []string{"Origin", "null"}, 403, "CrossOrigin", ""},
+		{"C3 foreign", cfID, body, []string{"Origin", "http://localhost:7645"}, 403, "CrossOrigin", ""},
+		{"C4 same-site", cfID, body, []string{"Sec-Fetch-Site", "same-site"}, 403, "CrossOrigin", "Sec-Fetch-Site"},
+		{"C4 cross-site", cfID, body, []string{"Sec-Fetch-Site", "cross-site"}, 403, "CrossOrigin", "Sec-Fetch-Site"},
+		{"C5 text", cfID, body, []string{"Content-Type", "text/plain"}, 415, "UnsupportedMediaType", ""},
+		{"C6 oversize", cfID, append(body[:len(body)-1:len(body)-1], []byte(`,"x":"`+strings.Repeat("a", 20<<10)+`"}`)...), nil, 413, "TooLarge", ""},
+		{"C6 typed", cfID, []byte(`{"v":1,"digest":"` + row.Digest + `","typed":"y"}`), nil, 400, "BadRequest", `unknown field "typed"`},
+		{"C6 y", cfID, []byte(`{"v":1,"digest":"` + row.Digest + `","sshsig":` + fmt.Sprintf("%q", sig) + `,"y":true}`), nil, 400, "BadRequest", "unknown field"},
+		{"C6 both formats", cfID, []byte(`{"v":1,"digest":"` + row.Digest + `","sshsig":` + fmt.Sprintf("%q", sig) + `,"webauthn":{}}`), nil, 400, "BadRequest", "both"},
+		{"C6 neither", cfID, []byte(`{"v":1,"digest":"` + row.Digest + `"}`), nil, 400, "BadRequest", "neither"},
+		{"C6 v2", cfID, []byte(`{"v":2,"digest":"` + row.Digest + `","sshsig":` + fmt.Sprintf("%q", sig) + `}`), nil, 400, "BadRequest", "v is 2"},
+		{"C6 two values", cfID, append(append([]byte{}, body...), body...), nil, 400, "BadRequest", "more than one"},
+	} {
+		got := requireRefusal(t, tc.name, a.confirm(tc.id, tc.body, tc.hdr...), tc.status, tc.class)
+		if !strings.Contains(got.Reason, tc.reason) {
+			t.Errorf("%s: reason %q, want %q", tc.name, got.Reason, tc.reason)
+		}
+	}
+	requireRefusal(t, "R2 foreign host", a.confirmAt("rebind.example:7644", cfID, body), 421, "MisdirectedRequest")
+	if after := a.confirmState(); after != before {
+		t.Fatalf("a guard wrote:\n%s\n%s", before, after)
+	}
+	if rec := a.confirm(cfID, body); rec.Code != 200 {
+		t.Fatalf("control: %d %s", rec.Code, rec.Body.String())
+	}
+}
diff --git a/host/surfacepage/confirm_confirmed_test.go b/host/surfacepage/confirm_confirmed_test.go
new file mode 100644
index 00000000..8a0e07e3
--- /dev/null
+++ b/host/surfacepage/confirm_confirmed_test.go
@@ -0,0 +1,18 @@
+package surfacepage
+
+// Row 243 V3 sprint 4 M10 (design §2.6 "ok"): the POST's success body, the design's wire.
+
+import (
+	"encoding/json"
+	"strings"
+	"testing"
+)
+
+func TestConfirmedWire(t *testing.T) {
+	ref := "sha256:" + strings.Repeat("5e", 32)
+	got, err := json.Marshal(Confirmed("vr-0123456789abcdef", "webauthn:mark", ref))
+	want := `{"v":1,"confirmed":true,"request":"vr-0123456789abcdef","approvedBy":"webauthn:mark","decision":"` + ref + `","next":"the waiting episode finish opens the pull request"}`
+	if err != nil || string(got) != want {
+		t.Fatalf("Confirmed = %s (%v), want %s", got, err, want)
+	}
+}
diff --git a/host/surfacepage/confirm_model.go b/host/surfacepage/confirm_model.go
index ed52fc59..e098477d 100644
--- a/host/surfacepage/confirm_model.go
+++ b/host/surfacepage/confirm_model.go
@@ -70,3 +70,20 @@ func confirmationState(in ConfirmationInput) string {
 	}
 	return ConfirmationPending
 }
+
+// ConfirmedJSON is POST /surface/confirm/{id}'s success (design §2.6): the world store's finish
+// approval, decided by the key, is named by the request.
+type ConfirmedJSON struct {
+	V          int    `json:"v"`
+	Confirmed  bool   `json:"confirmed"`
+	Request    string `json:"request"`
+	ApprovedBy string `json:"approvedBy"`
+	Decision   string `json:"decision"`
+	Next       string `json:"next"`
+}
+
+// Confirmed is the success body for request id confirmed by approvedBy with decision.
+func Confirmed(id, approvedBy, decision string) ConfirmedJSON {
+	return ConfirmedJSON{V: 1, Confirmed: true, Request: id, ApprovedBy: approvedBy, Decision: decision,
+		Next: "the waiting episode finish opens the pull request"}
+}
````

### Patch M10b — CX7–CX13 and the sprint-3 judge's note 1

````diff M10b
diff --git a/host/daemon/surface_confirm_test.go b/host/daemon/surface_confirm_test.go
index e63045a3..6df65ff9 100644
--- a/host/daemon/surface_confirm_test.go
+++ b/host/daemon/surface_confirm_test.go
@@ -14,11 +14,14 @@ import (
 	"net"
 	"net/http"
 	"net/http/httptest"
+	"os"
 	"strings"
 	"testing"
 
 	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/broker"
 	"github.com/sunholo-data/ailang-world/host/store"
 )
 
@@ -173,3 +176,140 @@ func TestConfirmPOSTGuards(t *testing.T) {
 		t.Fatalf("control: %d %s", rec.Code, rec.Body.String())
 	}
 }
+
+// CX7-CX13: each refusal after the lock and before any write — C7 unknown and below schema 9, C8
+// each Admit verdict (expiry at exactly expires_at), C9, C10 (decode, Validate, F1-F5), C11 the law,
+// C12 the binding (through the judge seam), C13 a busy writer — answers its status with the store
+// and the world untouched.
+func TestConfirmPOSTRefusesBeforeAnyWrite(t *testing.T) {
+	a, row := newConfirmRig(t)
+	sig := func(s *approvaltest.Signer, namespace string, of store.VerbRequestRow) []byte {
+		m, _ := verbMessageOf(t, of)
+		return s.Sign(namespace, m)
+	}
+	body := sshBody(row.Digest, sig(apSigner, approval.VerbNamespace, row))
+	d, _ := hex.DecodeString(row.Digest)
+	c, _ := hex.DecodeString(row.Challenge)
+	mintHash, _ := approval.MessageHash([32]byte(d), [32]byte(c), row.ExpiresAt)
+	type cx struct {
+		name, id      string
+		body          []byte
+		now           int64
+		status        int
+		class, reason string
+	}
+	cases := []cx{
+		{"C7 unknown", "vr-00000000000000ff", body, apT0, 404, "NotFound", "no verb request vr-00000000000000ff"},
+		{"C8 clock behind", cfID, body, row.CreatedAt - 1, 409, "Conflict", "clock is behind"},
+		{"C8 expired at expires_at", cfID, body, row.ExpiresAt, 409, "Conflict", "expired"},
+		{"C9 digest", cfID, sshBody(strings.Repeat("ab", 32), sig(apSigner, approval.VerbNamespace, row)), apT0, 409, "Conflict", "changed since the card was shown"},
+		{"C9 before C11", cfID, sshBody(strings.Repeat("ab", 32), []byte("not a signature")), apT0, 409, "Conflict", "changed since the card was shown"},
+		{"C11 unknown key", cfID, sshBody(row.Digest, sig(approvaltest.NewSigner(99), approval.VerbNamespace, row)), apT0, 422, "NotApproved", "SSH signature does not confirm"},
+		{"C11 namespace", cfID, sshBody(row.Digest, sig(apSigner, approval.Namespace, row)), apT0, 422, "NotApproved", "SSH signature does not confirm"},
+		{"C11 parse", cfID, sshBody(row.Digest, []byte("not a signature")), apT0, 422, "NotApproved", "SSH signature does not confirm"},
+		{"C11 W6 mint hash", cfID, approveBody(row.Digest, apPasskey.Assert(apOrigin, "localhost", mintHash, 2)), apT0, 422, "NotApproved", "passkey assertion does not confirm"},
+		{"C11 W15 counter", cfID, confirmPasskey(t, apPasskey, row, 1), apT0, 422, "NotApproved", "passkey assertion does not confirm"},
+	}
+	for _, f := range []struct {
+		id, tail, reason string
+		edit             func(r *approval.VerbRequest)
+	}{
+		{"vr-00000000000000a1", " ", "cannot be confirmed", func(r *approval.VerbRequest) {}},
+		{"vr-00000000000000a2", "", "verb is", func(r *approval.VerbRequest) { r.Verb = "episode start" }},
+		{"vr-00000000000000a3", "", "not a finish approval", func(r *approval.VerbRequest) { r.Scope = "x" }},
+		{"vr-00000000000000a4", "", "outlives its approval scope", func(r *approval.VerbRequest) {
+			s, _ := broker.ParseFinishApprovalScope(r.Scope)
+			s.ExpiresAt = apT0 - 30 + 600 // = the request's expiry
+			r.Scope, r.Prompt = broker.FormatFinishApprovalScope(s), broker.FinishPromptText(s, r.Title, r.Body, r.Protected, r.Committer)
+		}},
+		{"vr-00000000000000a5", "", "title is not", func(r *approval.VerbRequest) { r.Title += "!" }},
+		{"vr-00000000000000a6", "", "body is not", func(r *approval.VerbRequest) { r.Body += "!" }},
+		{"vr-00000000000000a7", "", "prompt is not", func(r *approval.VerbRequest) { r.Prompt += "!" }},
+	} {
+		r := verbFixture(apT0-30, "w-"+f.id[len(f.id)-2:], true)
+		f.edit(&r)
+		v := a.putVerbBytes(f.id, r.Episode, apT0-30, append(canonVerb(t, r), f.tail...))
+		cases = append(cases, cx{"C10 " + f.reason, v.RequestID, sshBody(v.Digest, sig(apSigner, approval.VerbNamespace, v)), apT0, 409, "NotConfirmable", f.reason})
+	}
+	for _, tc := range cases {
+		before := a.confirmState()
+		a.d.approveNow = func() int64 { return tc.now }
+		got := requireRefusal(t, tc.name, a.confirm(tc.id, tc.body), tc.status, tc.class)
+		if !strings.Contains(got.Reason, tc.reason) {
+			t.Errorf("%s: reason %q, want %q", tc.name, got.Reason, tc.reason)
+		}
+		if after := a.confirmState(); after != before {
+			t.Errorf("%s wrote:\n%s\n%s", tc.name, before, after)
+		}
+	}
+	a.d.approveNow = func() int64 { return apT0 }
+	// C8: a hand-edited digest (signed over the request's real one); a superseded request; a
+	// confirmed one. None writes.
+	edited := a.putVerb("vr-00000000000000c3", "w-c3", apT0-30, true)
+	setVerbColumn(t, a.d.credentialsPath, edited.RequestID, "digest", strings.Repeat("0", 64))
+	newer := a.putVerb("vr-00000000000000c2", cfEp, apT0-20, true)
+	newerBody := sshBody(newer.Digest, sig(apSigner, approval.VerbNamespace, newer))
+	before := a.confirmState()
+	requireRefusal(t, "C8 digest", a.confirm(edited.RequestID, sshBody(strings.Repeat("0", 64), sig(apSigner, approval.VerbNamespace, edited))), 409, "Conflict")
+	if got := requireRefusal(t, "C8 superseded", a.confirm(cfID, body), 409, "Conflict"); !strings.Contains(got.Reason, "request "+cfID+" was superseded by a newer finish of "+cfEp) || a.confirmState() != before {
+		t.Errorf("superseded: %q (or it wrote)", got.Reason)
+	}
+	if rec := a.confirm(newer.RequestID, newerBody); rec.Code != 200 {
+		t.Fatalf("newer: %d %s", rec.Code, rec.Body.String())
+	}
+	before = a.confirmState()
+	if got := requireRefusal(t, "C8 consumed", a.confirm(newer.RequestID, newerBody), 409, "Conflict"); !strings.Contains(got.Reason, "already confirmed by ssh:mark-ssh") || a.confirmState() != before {
+		t.Errorf("consumed: %q (or it wrote)", got.Reason)
+	}
+	// C12 through the judge seam (a verified confirmation of another request); C13 a held writer
+	// (503 Busy, Retry-After, never a wait); C11 an approver revoked since enrolment.
+	other := a.putVerb("vr-00000000000000c4", "w-c4", apT0-30, true)
+	unbound := a.putVerb("vr-00000000000000c5", "w-c5", apT0-30, true)
+	before = a.confirmState()
+	judgeVerbSSH = func(ctx context.Context, ro *store.Credentials, in authority.VerbJudgeInput, armored []byte) (*approval.SignedVerb, error) {
+		od, _ := hex.DecodeString(other.Digest)
+		oc, _ := hex.DecodeString(other.Challenge)
+		in.Digest, in.BytesDigest, in.Challenge = [32]byte(od), [32]byte(od), [32]byte(oc)
+		return authority.JudgeVerbSSH(ctx, ro, in, armored)
+	}
+	t.Cleanup(func() { judgeVerbSSH = authority.JudgeVerbSSH })
+	if got := requireRefusal(t, "C12", a.confirm(unbound.RequestID, sshBody(unbound.Digest, sig(apSigner, approval.VerbNamespace, other))), 422, "NotApproved"); !strings.Contains(got.Reason, "digest") {
+		t.Errorf("C12: %q", got.Reason)
+	}
+	judgeVerbSSH = func(context.Context, *store.Credentials, authority.VerbJudgeInput, []byte) (*approval.SignedVerb, error) {
+		return nil, &authority.JudgeError{Step: authority.JudgeStepStore, Err: &store.OlderCredentialsStoreError{Version: 4, Need: 5}}
+	}
+	requireRefusal(t, "C11 store", a.confirm(unbound.RequestID, sshBody(unbound.Digest, sig(apSigner, approval.VerbNamespace, unbound))), 503, "Restart")
+	judgeVerbSSH = authority.JudgeVerbSSH
+	otherBody := sshBody(other.Digest, sig(apSigner, approval.VerbNamespace, other))
+	held, err := store.OpenCredentials(a.d.credentialsPath)
+	if err != nil {
+		t.Fatal(err)
+	}
+	rec := a.confirm(other.RequestID, otherBody)
+	_ = held.Close()
+	if requireRefusal(t, "C13 busy", rec, 503, "Busy"); rec.Header().Get("Retry-After") != "1" {
+		t.Errorf("C13 Retry-After %q", rec.Header().Get("Retry-After"))
+	}
+	if after := a.confirmState(); after != before {
+		t.Fatalf("C12/C13 wrote:\n%s\n%s", before, after)
+	}
+	a.write(func(ctx context.Context, c *store.Credentials) error {
+		return c.RevokeApprover(ctx, "mark-ssh", apT0-5)
+	})
+	before = a.confirmState()
+	requireRefusal(t, "C11 revoked", a.confirm(other.RequestID, otherBody), 422, "NotApproved")
+	if after := a.confirmState(); after != before {
+		t.Fatalf("C11 revoked wrote:\n%s\n%s", before, after)
+	}
+	// Below schema 9: C7's read refuses first, 503 naming the serve restart.
+	setVersion(t, a.d.credentialsPath, `DROP TABLE verb_requests`, `PRAGMA user_version = 8`)
+	if got := requireRefusal(t, "C7 v8", a.confirm(cfID, body), 503, "Restart"); !strings.Contains(got.Reason, "restart `ailang-worldd serve`") {
+		t.Errorf("v8: %q", got.Reason)
+	}
+	// No credentials store to read: the sanitized 500.
+	if err := os.Remove(a.d.credentialsPath); err != nil {
+		t.Fatal(err)
+	}
+	requireRefusal(t, "C7 unreadable", a.confirm(cfID, body), 500, "Internal")
+}
diff --git a/host/store/verb_requests_bound_test.go b/host/store/verb_requests_bound_test.go
new file mode 100644
index 00000000..d77a987f
--- /dev/null
+++ b/host/store/verb_requests_bound_test.go
@@ -0,0 +1,32 @@
+package store
+
+// Row 243 V3 sprint 4 (the sprint-3 judge's note 1): a verb request whose expires_at equals now is
+// expired everywhere, as ConsumeVerbRequest's now >= ExpiresAt (CS6 "expired"): it is not listed as
+// pending, and it does not count toward the 16 cap. One second earlier it is live in both.
+
+import (
+	"errors"
+	"fmt"
+	"testing"
+)
+
+func TestVerbRequestExpiryBoundary(t *testing.T) {
+	c := vrStore(t)
+	for n := 1; n <= MaxPendingVerbRequests; n++ {
+		vrPut(t, c, vrRow(t, n, fmt.Sprintf("w%d", n), 1000), 1000) // each expires at 1600
+	}
+	before := vrState(t, c)
+	if rows, err := c.PendingVerbRequests(testCtx(t), 1599, 0); err != nil || len(rows) != MaxPendingVerbRequests {
+		t.Fatalf("at 1599: %d pending, %v; want 16", len(rows), err)
+	}
+	if err := c.PutVerbRequest(testCtx(t), vrRow(t, 17, "w17", 1599), 1599); !errors.Is(err, ErrTooManyPendingVerbRequests) || vrState(t, c) != before {
+		t.Fatalf("at 1599 the 17th: %v (or a row was written), want the cap", err)
+	}
+	if rows, err := c.PendingVerbRequests(testCtx(t), 1600, 0); err != nil || len(rows) != 0 {
+		t.Fatalf("at 1600 (= expires_at): %d pending, %v; want 0", len(rows), err)
+	}
+	vrPut(t, c, vrRow(t, 17, "w17", 1600), 1600) // the 16 expired at 1600 do not count
+	if rows, err := c.PendingVerbRequests(testCtx(t), 1600, 0); err != nil || len(rows) != 1 || rows[0].RequestID != fmt.Sprintf("vr-%016x", 17) {
+		t.Fatalf("after the put: %+v, %v", rows, err)
+	}
+}
````

### Patch M10c — the writes: CX14–CX21, audit note 3, C15's mapping

````diff M10c
diff --git a/host/daemon/surface_confirm_write_test.go b/host/daemon/surface_confirm_write_test.go
new file mode 100644
index 00000000..5195df0d
--- /dev/null
+++ b/host/daemon/surface_confirm_write_test.go
@@ -0,0 +1,312 @@
+package daemon
+
+// Row 243 V3 sprint 4 M10 (design §2.6, §6 CX14-CX21; S3 audit note 3's daemon half): the
+// confirmation's writes. C14 records the world store's finish approval decided by the key; C15
+// consumes the row and names it in one credentials transaction; the writer and humanMu are held
+// across both. Each crash window of the design's table is driven by its instrument: an objects
+// trigger refusing ApprovalDecisionV1 (A4's RD13 (a)), a verb_requests UPDATE trigger, and the C15
+// seam for a supersession or an approver change between C14 and C15.
+
+import (
+	"bytes"
+	"context"
+	"database/sql"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"net/http/httptest"
+	"strings"
+	"sync"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// worldDecision is the finish approval decision ref names: its decision and decidedBy, and its
+// request's effect, requester and scope, read back from the world store.
+func (a approvalRig) worldDecision(ref string) (verb, decidedBy, effect, requester, scope string) {
+	a.t.Helper()
+	ctx := boundedTestContext(a.t)
+	obj, ok, err := a.d.store.GetObject(ctx, hashref.MustParse(ref))
+	var wire struct {
+		RequestRef string `json:"requestRef"`
+		Decision   string `json:"decision"`
+		DecidedBy  string `json:"decidedBy"`
+	}
+	if err != nil || !ok || obj.SemanticID != broker.ApprovalDecisionV1 || json.Unmarshal(obj.Payload, &wire) != nil {
+		a.t.Fatalf("decision %s: %v %v %s", ref, ok, err, obj.SemanticID)
+	}
+	q := hashref.MustParse(wire.RequestRef)
+	reqObj, ok, err := a.d.store.GetObject(ctx, q)
+	req, derr := broker.DecodeApprovalRequest(reqObj.Payload)
+	if err != nil || !ok || derr != nil {
+		a.t.Fatalf("request %s: %v %v %v", q, ok, err, derr)
+	}
+	if onChain, decided, chainVerb, chainBy, err := broker.ApprovalChainState(ctx, a.d.store, q); err != nil || !onChain || !decided || chainVerb != wire.Decision || chainBy != wire.DecidedBy {
+		a.t.Fatalf("chain of %s: %v %v %q %q %v", q, onChain, decided, chainVerb, chainBy, err)
+	}
+	return wire.Decision, wire.DecidedBy, req.Effect, req.Requester, req.Scope
+}
+
+// CX17: an SSHSIG and a passkey assertion each confirm (the passkey one second before expiry): the
+// answer is surfacepage.Confirmed of the row's decision; the world decision approves the request's
+// own scope for operator:episode-finish, decided by the key; the row is consumed with five columns;
+// the approver is used and the passkey's counter moved; a replay is 409 and writes nothing.
+func TestConfirmPOSTConfirms(t *testing.T) {
+	a, row := newConfirmRig(t)
+	pk := a.putVerb("vr-00000000000000d1", "w-d1", apT0-30, false)
+	m, _ := verbMessageOf(t, row)
+	armored := apSigner.Sign(approval.VerbNamespace, m)
+	_, h := verbMessageOf(t, pk)
+	assertion := apPasskey.Assert(apOrigin, "localhost", h, 2)
+	for i, tc := range []struct {
+		row       store.VerbRequestRow
+		body      []byte
+		artifact  []byte
+		now       int64
+		by        string
+		ep, draft string
+	}{
+		{row, sshBody(row.Digest, armored), armored, apT0, "ssh:mark-ssh", cfEp, "draft"},
+		{pk, approveBody(pk.Digest, assertion), assertion, pk.ExpiresAt - 1, "webauthn:mark", "w-d1", "ready"},
+	} {
+		a.d.approveNow = func() int64 { return tc.now }
+		rec := a.confirm(tc.row.RequestID, tc.body)
+		got := a.readVerb(tc.row.RequestID)
+		want, _ := json.Marshal(surfacepage.Confirmed(tc.row.RequestID, tc.by, got.DecisionRef))
+		if rec.Code != 200 || !bytes.Equal(bytes.TrimSuffix(rec.Body.Bytes(), []byte("\n")), want) || approvalHeaders(rec) != nil {
+			t.Fatalf("%s: %d %s, want %s", tc.by, rec.Code, rec.Body.String(), want)
+		}
+		if got.Outcome != store.VerbOutcomeApproved || got.ApprovedBy != tc.by || !bytes.Equal(got.Artifact, tc.artifact) || got.ConsumedAt != tc.now || !strings.HasPrefix(got.DecisionRef, "sha256:") {
+			t.Fatalf("%s: row %+v", tc.by, got)
+		}
+		verb, decidedBy, effect, requester, scope := a.worldDecision(got.DecisionRef)
+		if wantScope := verbFixture(apT0-30, tc.ep, tc.draft == "draft").Scope; verb != "approve" || decidedBy != tc.by || effect != broker.EffectHumanApprove || requester != "operator:episode-finish" || scope != wantScope {
+			t.Fatalf("%s: world decision %q %q %q %q %q", tc.by, verb, decidedBy, effect, requester, scope)
+		}
+		if n, _ := (decideRig{a}).decisions(); n != i+1 {
+			t.Fatalf("%s: %d decisions", tc.by, n)
+		}
+		before := a.confirmState()
+		if r := requireRefusal(t, tc.by+" replay", a.confirm(tc.row.RequestID, tc.body), 409, "Conflict"); !strings.Contains(r.Reason, "already confirmed by "+tc.by) || a.confirmState() != before {
+			t.Fatalf("%s replay: %q", tc.by, r.Reason)
+		}
+	}
+	if s := a.credState(); !strings.Contains(s, "[mark-ssh 1800000000 0") || !strings.Contains(s, "[mark 2 ") {
+		t.Fatalf("approver use or counter not moved: %s", s)
+	}
+}
+
+// CX19: the confirmation holds humanMu, the mutex the typed finish approve phase takes
+// (host/daemon/episode_finish.go), from its clock read through C11 to C15 (each probe's TryLock fails
+// while the handler holds it); eight concurrent POSTs of one body confirm once, never a busy writer.
+func TestConfirmPOSTHoldsHumanMu(t *testing.T) {
+	a, row := newConfirmRig(t)
+	body := confirmSSH(t, apSigner, approval.VerbNamespace, row)
+	var mu sync.Mutex
+	var probes []string
+	probe := func(at string) {
+		held := !a.d.humanMu.TryLock()
+		if !held {
+			a.d.humanMu.Unlock()
+		}
+		mu.Lock()
+		probes = append(probes, fmt.Sprintf("%s:%v", at, held))
+		mu.Unlock()
+	}
+	a.d.approveNow = func() int64 { probe("clock"); return apT0 }
+	judgeVerbSSH = func(ctx context.Context, ro *store.Credentials, in authority.VerbJudgeInput, armored []byte) (*approval.SignedVerb, error) {
+		probe("C11")
+		return authority.JudgeVerbSSH(ctx, ro, in, armored)
+	}
+	consumeVerbRequest = func(c *store.Credentials, ctx context.Context, id, digest string, use store.ApproverUse, by string, artifact []byte, ref string, now int64) error {
+		probe("C15")
+		return c.ConsumeVerbRequest(ctx, id, digest, use, by, artifact, ref, now)
+	}
+	t.Cleanup(func() {
+		judgeVerbSSH, consumeVerbRequest = authority.JudgeVerbSSH, (*store.Credentials).ConsumeVerbRequest
+	})
+	if rec := a.confirm(cfID, body); rec.Code != 200 || fmt.Sprint(probes) != "[clock:true C11:true C15:true]" {
+		t.Fatalf("%d %s; probes %v", rec.Code, rec.Body.String(), probes)
+	}
+	next := a.putVerb("vr-00000000000000c6", "w-c6", apT0-30, true)
+	body = confirmSSH(t, apSigner, approval.VerbNamespace, next)
+	codes := make([]int, 8)
+	var wg sync.WaitGroup
+	for i := range codes {
+		wg.Add(1)
+		go func(i int) { defer wg.Done(); codes[i] = a.confirm(next.RequestID, body).Code }(i)
+	}
+	wg.Wait()
+	ok := 0
+	for _, c := range codes {
+		if c == 200 {
+			ok++
+		} else if c != 409 {
+			t.Errorf("a concurrent POST answered %d (codes %v)", c, codes)
+		}
+	}
+	if n, _ := (decideRig{a}).decisions(); ok != 1 || n != 2 {
+		t.Fatalf("%d confirmed, %d decisions (codes %v)", ok, n, codes)
+	}
+}
+
+// CX20 (window "C14 fails"): the world write fails (A4's objects trigger) → 500, no decision, the
+// row pending and the counter unmoved; the identical passkey body after the trigger is dropped → 200,
+// one decision, the row names it; a replay → 409.
+func TestConfirmPOSTRetriesAFailedWorldWrite(t *testing.T) {
+	a, row := newConfirmRig(t)
+	body := confirmPasskey(t, apPasskey, row, 2)
+	cred := a.credState()
+	drop := (decideRig{a}).failWorldWrites()
+	if r := requireRefusal(t, "C14 fails", a.confirm(cfID, body), 500, "Internal"); r.Reason != confirmIncomplete {
+		t.Fatalf("reason %q", r.Reason)
+	}
+	if n, _ := (decideRig{a}).decisions(); n != 0 || a.credState() != cred {
+		t.Fatalf("after the failure: %d decisions; credentials %s (was %s)", n, a.credState(), cred)
+	}
+	drop()
+	rec := a.confirm(cfID, body)
+	var c surfacepage.ConfirmedJSON
+	_ = json.Unmarshal(rec.Body.Bytes(), &c)
+	if n, _ := (decideRig{a}).decisions(); rec.Code != 200 || n != 1 || a.readVerb(cfID).DecisionRef != c.Decision {
+		t.Fatalf("retry: %d %s; %d decisions; row %+v", rec.Code, rec.Body.String(), n, a.readVerb(cfID))
+	}
+	requireRefusal(t, "replay", a.confirm(cfID, body), 409, "Conflict")
+}
+
+// CX21 (window "after C14, C15 fails"): a verb_requests UPDATE trigger fails the consume after the
+// world write → 500, one decision that no row names, the row pending, the counter unmoved. The
+// identical body after the trigger is dropped: in the same second the identical decision answers its
+// existing ref (W59: still one decision); a second later a second decision lands and the row names
+// it, the first left inert.
+func TestConfirmPOSTRetriesAFailedConsume(t *testing.T) {
+	a, row := newConfirmRig(t)
+	later := a.putVerb("vr-00000000000000c7", "w-c7", apT0-30, true)
+	db, err := sql.Open("sqlite", a.d.credentialsPath)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = db.Close() }()
+	trigger := func(sql string) {
+		if _, err := db.Exec(sql); err != nil {
+			t.Fatal(err)
+		}
+	}
+	for i, tc := range []struct {
+		row   store.VerbRequestRow
+		retry int64
+		n     int
+	}{{row, apT0, 1}, {later, apT0 + 1, 3}} {
+		body := confirmPasskey(t, apPasskey, tc.row, uint32(2+i))
+		a.d.approveNow = func() int64 { return apT0 }
+		trigger(`CREATE TRIGGER fail_consume BEFORE UPDATE ON verb_requests WHEN NEW.outcome = 'approved' BEGIN SELECT RAISE(ABORT, 'injected consume failure'); END`)
+		cred := a.credState()
+		n0, _ := (decideRig{a}).decisions()
+		if r := requireRefusal(t, "C15 fails", a.confirm(tc.row.RequestID, body), 500, "Internal"); r.Reason != confirmIncomplete {
+			t.Fatalf("reason %q", r.Reason)
+		}
+		if n, _ := (decideRig{a}).decisions(); n != n0+1 || a.credState() != cred {
+			t.Fatalf("after the failure: %d decisions (was %d); credentials %s (was %s)", n, n0, a.credState(), cred)
+		}
+		trigger(`DROP TRIGGER fail_consume`)
+		a.d.approveNow = func() int64 { return tc.retry }
+		rec := a.confirm(tc.row.RequestID, body)
+		var c surfacepage.ConfirmedJSON
+		_ = json.Unmarshal(rec.Body.Bytes(), &c)
+		if n, _ := (decideRig{a}).decisions(); rec.Code != 200 || n != tc.n || a.readVerb(tc.row.RequestID).DecisionRef != c.Decision {
+			t.Fatalf("retry at %d: %d %s; %d decisions, want %d", tc.retry, rec.Code, rec.Body.String(), n, tc.n)
+		}
+	}
+}
+
+// Audit note 3 and C15's refusals (window "after C14, C15 refuses"), through the C15 seam: the
+// request superseded by a newer finish of its episode between C14 and C15 → 409 naming it; a
+// concurrent consume that won → 409; the approver revoked, or the passkey's counter moved past the
+// assertion's → 422; any other store error → 500. Each leaves one decision no row names (inert, W62)
+// and, but for the won consume, the row unconsumed and the key unused.
+func TestConfirmPOSTRefusedAtC15(t *testing.T) {
+	a, _ := newConfirmRig(t)
+	t.Cleanup(func() { consumeVerbRequest = (*store.Credentials).ConsumeVerbRequest })
+	for i, tc := range []struct {
+		name          string
+		before        func(ctx context.Context, c *store.Credentials, row store.VerbRequestRow) error
+		passkey       bool
+		status        int
+		class, reason string
+		outcome       string
+	}{
+		{"superseded", func(ctx context.Context, c *store.Credentials, row store.VerbRequestRow) error {
+			n := row
+			n.RequestID = "vr-00000000000000e9" // the same episode
+			return c.PutVerbRequest(ctx, n, apT0)
+		}, false, 409, "Conflict", "superseded by a newer finish of w-e1", store.VerbOutcomeSuperseded},
+		{"consume lost", nil, false, 409, "Conflict", "changed while it was confirmed", store.VerbOutcomeApproved},
+		{"store error", func(context.Context, *store.Credentials, store.VerbRequestRow) error { return errors.New("injected") }, false, 500, "Internal", confirmIncomplete, ""},
+		{"counter", func(context.Context, *store.Credentials, store.VerbRequestRow) error {
+			db, err := sql.Open("sqlite", a.d.credentialsPath)
+			if err == nil {
+				_, err = db.Exec(`UPDATE webauthn_credentials SET sign_count = 9 WHERE name = 'mark'`)
+				_ = db.Close()
+			}
+			return err
+		}, true, 422, "NotApproved", "counter", ""},
+		{"approver revoked", func(ctx context.Context, c *store.Credentials, _ store.VerbRequestRow) error {
+			return c.RevokeApprover(ctx, "mark-ssh", apT0)
+		}, false, 422, "NotApproved", "revoked", ""},
+	} {
+		row := a.putVerb(fmt.Sprintf("vr-00000000000000e%d", i+1), fmt.Sprintf("w-e%d", i+1), apT0-30, true)
+		body := confirmSSH(t, apSigner, approval.VerbNamespace, row)
+		if tc.passkey {
+			body = confirmPasskey(t, apPasskey, row, 2)
+		}
+		consumeVerbRequest = func(c *store.Credentials, ctx context.Context, id, digest string, use store.ApproverUse, by string, artifact []byte, ref string, now int64) error {
+			if tc.before == nil { // another consume of the same row wins first
+				if err := c.ConsumeVerbRequest(ctx, id, digest, use, by, artifact, ref, now); err != nil {
+					return err
+				}
+			} else if err := tc.before(ctx, c, row); err != nil {
+				return err
+			}
+			return c.ConsumeVerbRequest(ctx, id, digest, use, by, artifact, ref, now)
+		}
+		n0, _ := (decideRig{a}).decisions()
+		if r := requireRefusal(t, tc.name, a.confirm(row.RequestID, body), tc.status, tc.class); !strings.Contains(r.Reason, tc.reason) {
+			t.Errorf("%s: reason %q", tc.name, r.Reason)
+		}
+		got := a.readVerb(row.RequestID)
+		if n, _ := (decideRig{a}).decisions(); n != n0+1 || got.Outcome != tc.outcome || (tc.outcome != store.VerbOutcomeApproved && got.DecisionRef != "") {
+			t.Errorf("%s: %d decisions (was %d), row %+v", tc.name, n, n0, got)
+		}
+	}
+	if s := a.credState(); !strings.Contains(s, "[mark 9 ") || strings.Count(s, " approved ") != 1 { // 9: the counter row's own edit
+		t.Errorf("a refused C15 used a key or consumed a row: %s", s)
+	}
+	consumeVerbRequest = (*store.Credentials).ConsumeVerbRequest
+}
+
+// C15's sentinel → status mapping, operand by operand, as ConsumeVerbRequest wraps them.
+func TestConsumeRefusalOperands(t *testing.T) {
+	a, row := newConfirmRig(t)
+	for _, tc := range []struct {
+		err    error
+		status int
+	}{
+		{store.ErrVerbRequestSuperseded, 409}, {store.ErrVerbRequestConsumed, 409}, {store.ErrVerbRequestUnknown, 409},
+		{store.ErrVerbRequestExpired, 409}, {store.ErrVerbRequestClockBehind, 409}, {store.ErrVerbRequestDigestChanged, 409},
+		{store.ErrApproverRevoked, 422}, {store.ErrApproverKeyChanged, 422}, {store.ErrApproverUnknown, 422}, {store.ErrApproverClockBehind, 422},
+		{store.ErrSignCountRegression, 422}, {store.ErrPasskeyChanged, 422},
+		{&store.OlderCredentialsStoreError{Version: 8, Need: 9}, 503}, {store.ErrBadVerbRequest, 500}, {errors.New("disk I/O error"), 500},
+	} {
+		rec := httptest.NewRecorder()
+		a.d.consumeRefusal(rec, httptest.NewRequest("POST", "/surface/confirm/"+cfID, nil), row, fmt.Errorf("%w: %s", tc.err, cfID))
+		if rec.Code != tc.status {
+			t.Errorf("%v → %d, want %d", tc.err, rec.Code, tc.status)
+		}
+	}
+}
````

### Patch M10d — the one POST registered; the tests through the mux

````diff M10d
diff --git a/host/daemon/daemon.go b/host/daemon/daemon.go
index 6d3e02c9..83e7a657 100644
--- a/host/daemon/daemon.go
+++ b/host/daemon/daemon.go
@@ -1103,7 +1103,7 @@ func releaseFromVersion(version string) string {
 // GET /surface/approve/{id} and its approve.js), each behind approvalRoute's Host check, and
 // row 201 A4's decision routes GET /surface/decisions/{ref} and POST /surface/decide/{ref}, and
 // row 243 V3's confirmation card GET /surface/confirm/{id}, its confirm.js and its data GET
-// /surface/confirmations/{id} (read-only; no confirmation POST before sprint 4).
+// /surface/confirmations/{id}, and its one write, POST /surface/confirm/{id} (sprint 4).
 // The registry pattern deliberately uses a multi-segment wildcard: registry semantic IDs
 // such as "world/epoch-registry/v1" contain slashes.
 func (d *Daemon) Handler() http.Handler {
@@ -1148,6 +1148,7 @@ func (d *Daemon) Handler() http.Handler {
 	mux.HandleFunc("GET /surface/decide/{ref}", d.approvalRoute(d.handleSurfaceDecidePage))
 	mux.HandleFunc("GET /surface/decide.js", d.approvalRoute(d.handleSurfaceDecideScript))
 	mux.HandleFunc("GET /surface/confirm/{id}", d.approvalRoute(d.handleSurfaceConfirmPage))
+	mux.HandleFunc("POST /surface/confirm/{id}", d.approvalRoute(d.handleSurfaceConfirm))
 	mux.HandleFunc("GET /surface/confirm.js", d.approvalRoute(d.handleSurfaceConfirmScript))
 	mux.HandleFunc("GET /surface/confirmations/{id}", d.approvalRoute(d.handleSurfaceConfirmation))
 	mux.HandleFunc("GET /surface/verb.json", d.approvalRoute(d.handleSurfaceVerb)) // row 243 V1: the bar's verbs, read-only
diff --git a/host/daemon/surface_confirm_read_test.go b/host/daemon/surface_confirm_read_test.go
index ebf41c0d..9b81bbce 100644
--- a/host/daemon/surface_confirm_read_test.go
+++ b/host/daemon/surface_confirm_read_test.go
@@ -69,8 +69,8 @@ func (a approvalRig) readVerb(id string) (row store.VerbRequestRow) {
 	return row
 }
 
-// The route table: exactly three GETs under /surface/confirm, each behind approvalRoute, and no
-// POST (sprint 4's); other methods 405 with no CORS header; a foreign Host 421; other names 404.
+// The route table: exactly three GETs and the one POST (sprint 4) under /surface/confirm, each behind
+// approvalRoute; other methods 405 with no CORS header; a foreign Host 421; other names 404.
 func TestConfirmRoutesRegistered(t *testing.T) {
 	src, err := os.ReadFile("daemon.go")
 	if err != nil {
@@ -85,7 +85,7 @@ func TestConfirmRoutesRegistered(t *testing.T) {
 		got = append(got, m[1])
 	}
 	sort.Strings(got)
-	if want := []string{"GET /surface/confirm.js", "GET /surface/confirm/{id}", "GET /surface/confirmations/{id}"}; !reflect.DeepEqual(got, want) {
+	if want := []string{"GET /surface/confirm.js", "GET /surface/confirm/{id}", "GET /surface/confirmations/{id}", "POST /surface/confirm/{id}"}; !reflect.DeepEqual(got, want) {
 		t.Fatalf("confirmation routes %v, want %v", got, want)
 	}
 	a := newApprovalRig(t)
@@ -93,6 +93,9 @@ func TestConfirmRoutesRegistered(t *testing.T) {
 	for _, target := range []string{"/surface/confirm/" + id, "/surface/confirm.js", "/surface/confirmations/" + id} {
 		requireRefusal(t, "foreign "+target, a.serve("GET", target, "rebind.example:7644", nil), 421, "MisdirectedRequest")
 		for _, method := range []string{"POST", "PUT", "DELETE", "OPTIONS"} {
+			if method == "POST" && target == "/surface/confirm/"+id {
+				continue // the confirmation (surface_confirm_test.go)
+			}
 			rec := a.serve(method, target, apLoopback, nil, "Origin", "http://localhost:7644", "Access-Control-Request-Method", "POST")
 			if rec.Code != 405 || rec.Header().Get("Access-Control-Allow-Origin") != "" {
 				t.Errorf("%s %s = %d, want 405 with no CORS header", method, target, rec.Code)
diff --git a/host/daemon/surface_confirm_test.go b/host/daemon/surface_confirm_test.go
index 6df65ff9..8f7a2826 100644
--- a/host/daemon/surface_confirm_test.go
+++ b/host/daemon/surface_confirm_test.go
@@ -6,13 +6,10 @@ package daemon
 // the stores before it (nothing written). Stores and keys are tests' (D-WORLD-109).
 
 import (
-	"bytes"
 	"context"
 	"database/sql"
 	"encoding/hex"
 	"fmt"
-	"net"
-	"net/http"
 	"net/http/httptest"
 	"os"
 	"strings"
@@ -68,20 +65,10 @@ func confirmPasskey(t *testing.T, pk *approvaltest.Passkey, row store.VerbReques
 }
 
 // confirmAt POSTs body to /surface/confirm/<id> with Host host, as the card would (loopback Origin,
-// same-origin, JSON); hdr pairs override those headers. Until M10d registers the route, the request
-// reaches the handler through approvalRoute, as the mux will route it.
+// same-origin, JSON); hdr pairs override those headers.
 func (a approvalRig) confirmAt(host, id string, body []byte, hdr ...string) *httptest.ResponseRecorder {
 	a.t.Helper()
-	req := httptest.NewRequest("POST", "/surface/confirm/"+id, bytes.NewReader(body))
-	req = req.WithContext(context.WithValue(req.Context(), http.LocalAddrContextKey, net.Addr(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 7644})))
-	req.Host = host
-	req.SetPathValue("id", id)
-	for h, i := append([]string{"Origin", apOrigin, "Sec-Fetch-Site", "same-origin", "Content-Type", "application/json"}, hdr...), 0; i+1 < len(h); i += 2 {
-		req.Header.Set(h[i], h[i+1])
-	}
-	rec := httptest.NewRecorder()
-	a.d.approvalRoute(a.d.handleSurfaceConfirm).ServeHTTP(rec, req)
-	return rec
+	return a.serve("POST", "/surface/confirm/"+id, host, body, append([]string{"Origin", apOrigin, "Sec-Fetch-Site", "same-origin", "Content-Type", "application/json"}, hdr...)...)
 }
 
 // confirm is confirmAt on the loopback Host.
````

### Patch M11 — `confirm.js`'s submit half, CP8–CP10

````diff M11
diff --git a/host/surfacepage/confirm/confirm.js b/host/surfacepage/confirm/confirm.js
index 12296d76..6fd20c77 100644
--- a/host/surfacepage/confirm/confirm.js
+++ b/host/surfacepage/confirm/confirm.js
@@ -3,10 +3,13 @@
 // the daemon, checks that its bytes hash to its digest, shows EVERY field of the signed bytes through
 // one table (and refuses a request with a field it cannot show), and builds the 86-byte world-verb/v1
 // message itself (never a hash from the server): Save the message and the ssh-keygen line. A state
-// other than pending shows that state and no controls. Text is set with textContent only. Sprint 4
-// adds the submit half (a passkey and a pasted SSH signature) inside #controls.
+// other than pending shows that state and no controls. Text is set with textContent only. The submit
+// half (sprint 4, design §2.6): a passkey whose challenge is the message hash, or a pasted SSH
+// signature, POSTed to /surface/confirm/<id>; the answer or the refusal is shown, then the card is
+// read again (after a refusal or a 500 it is still pending and the controls return).
 (() => {
-  const NET_MS = 10000; // A3 §2.7 [r2]: the data GET, its body included
+  const NET_MS = 10000; // A3 §2.7 [r2]: every GET and POST, its body included
+  const TIMEOUT_MS = 240000; // A2's ceremony bound (A2 §2.14 [r2])
   const PREFIX = 'world-verb/v1', NAMESPACE = 'world-verb@sunholo', LEN = 86; // approval.VerbPrefix, VerbNamespace, VerbMessageLen
   // Every key of a world/verb-request/v1 (approval.VerbRequest's json tags, in order), with its label
   // and JSON type. TestConfirmPageFieldTable compares this list with the Go type.
@@ -31,6 +34,15 @@
     if (typeof s !== 'string' || !/^[A-Za-z0-9_-]*$/.test(s)) throw new Error('not base64url');
     return Uint8Array.from(atob(s.replace(/-/g, '+').replace(/_/g, '/')), (c) => c.charCodeAt(0));
   };
+  const b64u = (buf) => {
+    let s = '';
+    for (const b of new Uint8Array(buf)) s += String.fromCharCode(b);
+    return btoa(s).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
+  };
+  // A2 §2.7's assertion file, built here rather than by toJSON() (as decide.js does).
+  const assertionJSON = (c) => ({ id: c.id, rawId: b64u(c.rawId), type: c.type,
+    response: { clientDataJSON: b64u(c.response.clientDataJSON), authenticatorData: b64u(c.response.authenticatorData),
+      signature: b64u(c.response.signature), userHandle: c.response.userHandle ? b64u(c.response.userHandle) : null } });
   const hex = (b) => Array.from(b, (x) => x.toString(16).padStart(2, '0')).join('');
   const unhex = (s) => (typeof s === 'string' && /^[0-9a-f]{64}$/.test(s) ? Uint8Array.from(s.match(/../g), (x) => parseInt(x, 16)) : null);
   const sha256 = async (bytes) => new Uint8Array(await crypto.subtle.digest('SHA-256', bytes));
@@ -39,11 +51,11 @@
   const SAFE = /^[A-Za-z0-9@%+=:,./_-]+$/;
   const quote = (argv) => argv.map((a) => (SAFE.test(a) ? a : "'" + a.split("'").join(`'"'"'`) + "'")).join(' ');
 
-  async function net(url) {
+  async function net(url, init) {
     const c = new AbortController();
     const timer = setTimeout(() => c.abort(), NET_MS);
     try {
-      const r = await fetch(url, { cache: 'no-store', credentials: 'omit', signal: c.signal, headers: { Accept: 'application/json' } });
+      const r = await fetch(url, Object.assign({ cache: 'no-store', credentials: 'omit', signal: c.signal, headers: { Accept: 'application/json' } }, init));
       const text = await r.text();
       let body = null;
       try { body = JSON.parse(text); } catch (_) { body = null; }
@@ -115,6 +127,8 @@
   }
 
   const id = (location.pathname.match(/^\/surface\/confirm\/(vr-[0-9a-f]{16})$/) || [])[1];
+  let digestHex = '', hash = null, group = null; // the pending card's, for the submit half
+  const buttons = (on) => { $('go').disabled = !on || !group; $('post-sig').disabled = !on; };
   async function load() {
     if (!id) return refuse('refused', 'This is not a confirmation card\'s address; nothing to confirm.');
     let res;
@@ -159,9 +173,55 @@
     $('save').download = `world-verb-${id}.msg`;
     $('sshline').textContent = `ssh-keygen -Y sign -n ${NAMESPACE} -f <key> world-verb-${id}.msg`;
     $('status').textContent = 'Waiting for your key: save the message and sign it with your enrolled SSH key.';
+    digestHex = s.digest;
+    hash = h;
+    group = (s.passkeys || []).find((g) => g.origin === location.origin) || null;
+    $('go').hidden = !group;
     $('controls').hidden = false;
+    buttons(true);
     setState('ready');
   }
 
+  // submit POSTs one key format; the answer or the refusal stays shown while the card is read again.
+  async function submit(format) {
+    buttons(false);
+    setState('posting');
+    let res;
+    try {
+      res = await net('/surface/confirm/' + id, { method: 'POST', headers: { 'Content-Type': 'application/json' },
+        body: JSON.stringify(Object.assign({ v: 1, digest: digestHex }, format)) });
+    } catch (e) {
+      $('answer').textContent = `Confirmation outcome unknown: ${why(e)}. Reading the request again.`;
+      return reread();
+    }
+    const b = res.body || {};
+    $('answer').textContent = res.status === 200 && b.confirmed === true
+      ? `Confirmed by ${b.approvedBy}: decision ${b.decision}; ${b.next}.`
+      : `Not confirmed (HTTP ${res.status} ${b.error || ''}): ${b.reason || 'no reason given'}`;
+    return reread();
+  }
+  const reread = () => { $('controls').hidden = true; return load(); };
+
+  async function confirmWithPasskey() {
+    buttons(false);
+    setState('waiting');
+    const controller = new AbortController();
+    const timer = setTimeout(() => controller.abort(), TIMEOUT_MS);
+    let cred;
+    try {
+      cred = await navigator.credentials.get({ publicKey: { challenge: hash, rpId: group.rp,
+        allowCredentials: group.allow.map((a) => ({ type: 'public-key', id: unb64u(a) })),
+        userVerification: 'required', timeout: TIMEOUT_MS }, signal: controller.signal });
+    } catch (e) {
+      $('answer').textContent = `Not confirmed: the browser made no passkey response (${e.name}: ${e.message}). Nothing was sent.`;
+      return reread();
+    } finally {
+      clearTimeout(timer);
+    }
+    return submit({ webauthn: assertionJSON(cred) });
+  }
+  $('go').addEventListener('click', () => { confirmWithPasskey(); });
+  $('post-sig').addEventListener('click', () => { submit({ sshsig: $('sig').value }); });
+
   load(); // the daemon will serve this page only at /surface/confirm/vr-<16 hex>
 })();
diff --git a/host/surfacepage/confirm/index.html b/host/surfacepage/confirm/index.html
index dc805185..e163d6c1 100644
--- a/host/surfacepage/confirm/index.html
+++ b/host/surfacepage/confirm/index.html
@@ -8,7 +8,7 @@
 body { font: 16px/1.5 system-ui, sans-serif; max-width: 52rem; margin: 2rem auto; padding: 0 1rem; color: #1d1d1f; }
 h1 { font-size: 1.4rem; }
 h2 { font-size: 1.05rem; margin: 1.2rem 0 0.3rem; }
-#what, #status { white-space: pre-wrap; }
+#what, #status, #answer { white-space: pre-wrap; }
 #what[data-error] { color: #b00020; font-weight: 600; }
 pre { font: 13px ui-monospace, monospace; white-space: pre-wrap; overflow-wrap: anywhere; background: #f5f5f7; padding: 0.6rem; margin: 0; }
 dl { display: grid; grid-template-columns: max-content 1fr; gap: 0.15rem 1rem; margin: 0; }
@@ -28,9 +28,13 @@ dd { margin: 0; font: 13px ui-monospace, monospace; white-space: pre-wrap; overf
 <pre id="prompt"></pre>
 <h2>Open this pull request?</h2>
 <p id="status"></p>
+<p id="answer"></p>
 <div id="controls" hidden>
+<p><button id="go" hidden type="button" disabled>Approve with your passkey</button></p>
 <p><a id="save">Save the message</a></p>
 <p id="sshline"></p>
+<p><label for="sig">Paste the signature</label><br><textarea id="sig"></textarea></p>
+<p><button id="post-sig" type="button" disabled>Confirm with this signature</button></p>
 </div>
 <h2>What your key signs</h2>
 <dl id="signed"></dl>
diff --git a/host/surfacepage/confirm_script_test.go b/host/surfacepage/confirm_script_test.go
index bac1dfe0..9979251d 100644
--- a/host/surfacepage/confirm_script_test.go
+++ b/host/surfacepage/confirm_script_test.go
@@ -66,6 +66,19 @@ func TestConfirmPageMessage(t *testing.T) { runConfirmHarness(t, "message") }
 // CP5: an approved, superseded, expired or unknown state shows that state and the fields, and no controls.
 func TestConfirmPageStates(t *testing.T) { runConfirmHarness(t, "states") }
 
+// CP8 (sprint 4): Approve with your passkey asks navigator.credentials.get for the message hash as the
+// challenge, the origin's group's rp and credentials, UV required, A2's 240 s abort; POSTs the
+// assertion file with the card's digest; shows the answer and reads the card again. No group for the
+// origin: no passkey button. An aborted or failed ceremony sends nothing and the controls return.
+func TestConfirmPageSubmitsAPasskey(t *testing.T) { runConfirmHarness(t, "passkey") }
+
+// CP9 (sprint 4): a refusal (422, 409, 500) or a lost POST is shown as not confirmed, never as a
+// confirmation, and the card is read again: still pending, the controls are back.
+func TestConfirmPageShowsRefusals(t *testing.T) { runConfirmHarness(t, "answers") }
+
+// CP10 (sprint 4): the pasted SSH signature is POSTed in ParseApproveBody's grammar, {v, digest, sshsig}.
+func TestConfirmPageSubmitsAnSSHSIG(t *testing.T) { runConfirmHarness(t, "sshsig") }
+
 const confirmHarness = `
 const assert = require('node:assert/strict');
 const fs = require('node:fs');
@@ -81,7 +94,9 @@ const withBytes = (card, b) => Object.assign({}, card, { request: Buffer.from(b)
 
 // A stand-in DOM: index.html's elements by id (with their hidden attribute), and created ones.
 class El {
-  constructor(doc, tag, id, hidden) { Object.assign(this, { doc, tagName: tag, id: id || '', hidden: !!hidden, dataset: {}, attrs: {}, kids: [], text: '' }); }
+  constructor(doc, tag, id, hidden) { Object.assign(this, { doc, tagName: tag, id: id || '', hidden: !!hidden, dataset: {}, attrs: {}, kids: [], text: '', on: {} }); }
+  addEventListener(t, f) { (this.on[t] = this.on[t] || []).push(f); }
+  click() { if (!this.disabled && !this.hidden) for (const f of this.on.click || []) f(); }
   set textContent(v) { this.text = String(v); this.kids = []; }
   get textContent() { return this.text + this.kids.map((k) => k.textContent).join(''); }
   set innerHTML(v) { this.doc.sinks.push(String(v)); }
@@ -107,9 +122,9 @@ const hang = (init) => new Promise((_, reject) => init.signal.addEventListener('
 // open runs confirm.js at path against fetch and waits, on the wall clock (at most 10 s), until the page leaves
 // 'loading': its SHA-256 digests resolve off the event loop, so no count of ticks bounds them (iteration 401 M7);
 // fire runs the page's stubbed timers once, as soon as it has asked (its fetch made and its abort timer set).
-async function open(path, fetch, fire) {
+async function open(path, fetch, fire, navigator) {
   const doc = makeDoc(), timers = [], blobs = [], fetches = [];
-  vm.runInNewContext(SRC, { document: doc, location: { pathname: path, origin: 'http://localhost:1' }, crypto: globalThis.crypto, TextEncoder, TextDecoder, atob, btoa, Blob, AbortController, console,
+  vm.runInNewContext(SRC, { document: doc, location: { pathname: path, origin: 'http://localhost:1' }, crypto: globalThis.crypto, TextEncoder, TextDecoder, atob, btoa, Blob, AbortController, console, navigator,
     URL: { createObjectURL: (b) => { blobs.push(b); return 'blob:http://localhost:1/' + blobs.length; } },
     setTimeout: (f, ms) => timers.push({ f, ms, cleared: false }), clearTimeout: (n) => { if (timers[n - 1]) timers[n - 1].cleared = true; },
     fetch: (url, init) => { fetches.push({ url, init }); return fetch(init); } });
@@ -121,6 +136,18 @@ async function open(path, fetch, fire) {
   return { doc, $, timers, blobs, fetches, state: doc.body.dataset.state, what: $('what').textContent };
 }
 const card = (c) => '/surface/confirm/' + c.id;
+// until waits, on the wall clock (at most 10 s), for ok() after a click: the page's work is promises.
+async function until(ok, name) {
+  for (const end = Date.now() + 10000; !ok();) {
+    assert(Date.now() < end, 'timed out: ' + name);
+    await new Promise((r) => setTimeout(r, 2));
+  }
+}
+// daemon answers the card's GETs with cards in turn (the last repeats) and its POSTs with post.
+const daemon = (cards, post) => (init) => (init.method === 'POST' ? post(init) : answer(200, cards.length > 1 ? cards.shift() : cards[0])());
+const posts = (p) => p.fetches.filter((f) => f.init.method === 'POST');
+const CRED = { id: 'Bwc', rawId: new Uint8Array([7, 7]).buffer, type: 'public-key',
+  response: { clientDataJSON: Buffer.from('cd'), authenticatorData: Buffer.from('ad'), signature: Buffer.from('sg'), userHandle: null } };
 // refused asserts a fail-closed page: state, the reason, nothing shown, no controls, no message made.
 function refused(p, state, reason, name) {
   assert.equal(p.state, state, name + ': ' + p.what);
@@ -227,6 +254,78 @@ const sections = {
     }
     assert.deepEqual([G['draft-approved'].card.state, G['draft-superseded'].card.state, G['draft-expired'].card.state], ['approved', 'superseded', 'expired']);
   },
+  async passkey() {
+    const g = G.draft, c = g.card, approved = G['draft-approved'].card, asks = [];
+    const nav = { credentials: { get: (o) => { asks.push(o); return Promise.resolve(CRED); } } };
+    const ok = { v: 1, confirmed: true, request: c.id, approvedBy: 'webauthn:mark', decision: 'sha256:' + '5e'.repeat(32), next: 'the waiting episode finish opens the pull request' };
+    const p = await open(card(c), daemon([c, approved], answer(200, ok)), false, nav);
+    assert(p.state === 'ready' && !p.$('go').hidden && !p.$('go').disabled, p.what);
+    p.$('go').click();
+    await until(() => p.state === 'closed' || p.doc.body.dataset.state === 'closed', 'passkey');
+    assert.equal(asks.length, 1);
+    const pk = asks[0].publicKey;
+    assert.equal(Buffer.from(pk.challenge).toString('hex'), g.hash); // the message hash, never the mint hash
+    assert.deepEqual([pk.rpId, pk.userVerification, pk.timeout, pk.allowCredentials.length, pk.allowCredentials[0].type, Buffer.from(pk.allowCredentials[0].id).toString('hex')],
+      ['localhost', 'required', 240000, 1, 'public-key', '0707']);
+    assert(asks[0].signal && p.timers.some((t) => t.ms === 240000 && t.cleared), 'the ceremony\'s abort timer is cleared');
+    const sent = posts(p);
+    assert.equal(sent.length, 1);
+    assert.equal(sent[0].url, '/surface/confirm/' + c.id);
+    assert.equal(sent[0].init.headers['Content-Type'], 'application/json');
+    assert.deepEqual(JSON.parse(sent[0].init.body), { v: 1, digest: c.digest, webauthn: { id: 'Bwc', rawId: 'Bwc', type: 'public-key',
+      response: { clientDataJSON: 'Y2Q', authenticatorData: 'YWQ', signature: 'c2c', userHandle: null } } });
+    assert.equal(p.$('answer').textContent, 'Confirmed by webauthn:mark: decision sha256:' + '5e'.repeat(32) + '; the waiting episode finish opens the pull request.');
+    assert(p.$('controls').hidden && p.fetches.length === 3 && p.$('status').textContent.startsWith('Already confirmed by webauthn:mark-mac'), p.$('status').textContent);
+    const elsewhere = [{ origin: 'https://other.tail0.ts.net', rp: 'other.tail0.ts.net', allow: ['Bwc'] }]; // a group for another origin only
+    const none = await open(card(c), answer(200, Object.assign({}, c, { passkeys: elsewhere })), false, nav);
+    assert(none.state === 'ready' && none.$('go').hidden && none.$('go').disabled && !none.$('post-sig').disabled);
+    for (const [name, get, fire] of [['aborted at 240 s', (o) => new Promise((_, reject) => o.signal.addEventListener('abort', () => reject(new DOMException('The operation was aborted.', 'AbortError')))), true],
+      ['refused by the browser', () => Promise.reject(new DOMException('The operation is not allowed.', 'NotAllowedError')), false]]) {
+      const q = await open(card(c), daemon([c], answer(500, {})), false, { credentials: { get } });
+      q.$('go').click();
+      if (fire) {
+        await until(() => q.timers.some((t) => t.ms === 240000), name + ' timer');
+        assert(q.doc.body.dataset.state === 'waiting' && q.$('go').disabled && q.$('post-sig').disabled, name + ': waiting for the passkey');
+        for (const t of q.timers) if (t.ms === 240000 && !t.cleared) t.f();
+      }
+      await until(() => q.$('answer').textContent !== '' && q.doc.body.dataset.state === 'ready', name);
+      assert(q.$('answer').textContent.startsWith('Not confirmed: the browser made no passkey response ('), q.$('answer').textContent);
+      assert(posts(q).length === 0 && !q.$('controls').hidden && !q.$('go').disabled, name);
+    }
+  },
+  async answers() {
+    const c = G.draft.card;
+    for (const [name, post, text] of [
+      ['422', answer(422, { v: 1, error: 'NotApproved', reason: 'the SSH signature does not confirm request' }), 'Not confirmed (HTTP 422 NotApproved): the SSH signature does not confirm request'],
+      ['409', answer(409, { v: 1, error: 'Conflict', reason: 'the request changed since the card was shown; reload it' }), 'Not confirmed (HTTP 409 Conflict): the request changed since the card was shown; reload it'],
+      ['500', answer(500, { v: 1, error: 'Internal', reason: 'the confirmation did not complete; approve again' }), 'Not confirmed (HTTP 500 Internal): the confirmation did not complete; approve again'],
+      ['200 not confirmed', answer(200, { v: 1, confirmed: false }), 'Not confirmed (HTTP 200 ): no reason given'],
+      ['lost', () => Promise.reject(new TypeError('Failed to fetch')), 'Confirmation outcome unknown: Failed to fetch. Reading the request again.']]) {
+      const p = await open(card(c), daemon([c], post));
+      p.$('sig').value = 'sig';
+      p.$('post-sig').click();
+      await until(() => posts(p).length === 1 && p.$('answer').textContent !== '' && p.doc.body.dataset.state === 'ready', name);
+      assert.equal(p.$('answer').textContent, text, name);
+      assert(!p.$('controls').hidden && !p.$('post-sig').disabled && p.fetches.length === 3, name);
+    }
+  },
+  async sshsig() {
+    const c = G.draft.card, approved = G['draft-approved'].card;
+    let release;
+    const held = new Promise((r) => { release = r; });
+    const p = await open(card(c), daemon([c, approved], () => held.then(answer(200, { v: 1, confirmed: true, request: c.id, approvedBy: 'ssh:mark', decision: 'sha256:' + '6f'.repeat(32), next: 'n' }))));
+    const armored = '-----BEGIN SSH SIGNATURE-----\nU1NIU0lH\n-----END SSH SIGNATURE-----\n';
+    p.$('sig').value = armored;
+    p.$('post-sig').click();
+    await until(() => posts(p).length === 1, 'sshsig sent');
+    assert(p.doc.body.dataset.state === 'posting' && p.$('post-sig').disabled && p.$('go').disabled, 'no second submit while one is in flight');
+    release();
+    await until(() => p.doc.body.dataset.state === 'closed', 'sshsig');
+    const sent = posts(p);
+    assert(sent.length === 1 && sent[0].url === '/surface/confirm/' + c.id && sent[0].init.credentials === 'omit');
+    assert.deepEqual(JSON.parse(sent[0].init.body), { v: 1, digest: c.digest, sshsig: armored });
+    assert.equal(p.$('answer').textContent, 'Confirmed by ssh:mark: decision sha256:' + '6f'.repeat(32) + '; n.');
+  },
 };
 const which = process.argv[2];
 sections[which]().then(() => console.log('CONFIRM ' + which + ' PASS'), (e) => { console.error(e); process.exit(1); });
````

