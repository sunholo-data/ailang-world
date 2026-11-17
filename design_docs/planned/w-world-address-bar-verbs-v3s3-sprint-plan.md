# Sprint plan — w-world-address-bar-verbs-v3, sprint S3 "the store and the read routes" (iteration 402, row 243)

> **THE SPRINT — row 243 phase V3, sprint S3 (design §0, §1, §2.2, §2.3 (M8/VT2), §2.4, §2.5, §2.6 (the seam only), §2.9
> rows S3 reaches, §3–§8, §11, §12, §14, §16), D-WORLD-118 = A:** credentials **schema 9** (`verb_requests`, the
> design's DDL verbatim; a fresh store equals one upgraded from 8 column for column), `store.PutVerbRequest` (supersede
> per episode, grammar, the 16 cap, the day-old prune), `VerbRequest`, `PendingVerbRequests`, `ConsumeVerbRequest` (one
> transaction: refusals in `Admit`'s order, `useApprover`, `usePasskey`, then one UPDATE of five columns incl.
> `decision_ref`; every refusal writes nothing), `OlderCredentialsStoreError{Feature: "verb requests", Need: 9}` below 9;
> **`broker.FinishConfirmRefusal`** F1–F5 (pure; no F6, no `DraftOnly`); and the **three GET registrations** behind
> `approvalRoute` — `GET /surface/confirm/{id}` (S2's `ConfirmPage`), `GET /surface/confirm.js` (S2's `ConfirmScript`),
> `GET /surface/confirmations/{id}` (S2's `ConfirmationCard` of the stored row) — plus `requests.json`'s
> `confirmations` list and `confirmationsUnavailable`. **No POST, no CLI change, no bar, no *Needs you* rendering, no
> `.ail`, `surfaceCSP` unchanged, no new credentials-store writer in the daemon.** S3 deletes
> `host/daemon/surface_confirm_unrouted_test.go` in the sub-milestone that registers the routes (M9a).
> - **Measured: 1,391 changed lines** — `git diff -M --numstat 9733813c HEAD -- cmd host docs` (insertions + deletions,
>   renames on; no rename occurs; **the same 1,391 with no path filter**, so nothing lands outside `cmd host docs`), after
>   a commit of each of the five scratch stages: **526 production** (`verb_requests.go` 293, `surface_confirm_read.go`
>   127, `finish_confirm.go` 41, `credentials_schema.sql` 20, `approve_model.go` 17, `credentials.go` 12,
>   `surface_approve.go` 9, `daemon.go` 7), **865 test** (incl. the 33-line deletion), **0 docs**. **Cap 1,400: margin 9.**
> - **No cut (S3-D1).** M7 was compiled first and measured **836** (design ≈ 570 raw: 1.47×); M8 108 (≈ 90); M9 464
>   (≈ 320); first full compile **1,408 > 1,400**. Trimmed by measurement (comments, one test duplicating
>   `TestCredentialsMethodsRequireADeadline`), then grown by the mutant-driven test rows (§4) to **1,391**. All three
>   milestones ship fully proven; nothing moves to an S3b.
> - **Five sub-milestones, each green at its commit** (§1.2 C14): **M7a** schema 9 + the store methods + every
>   pre-existing test edit the bump forces (377) → **M7b** CS1–CS7 (428) → **M8** `FinishConfirmRefusal`, VT2 (108) →
>   **M9a** the three GET routes, CR1/CR3/CR4, the unrouted test deleted (322) → **M9b** `requests.json`'s
>   `confirmations`, CR2 (156).
> - **Every line ships as a compiled patch** (Appendix P), extracted from this file and round-tripped onto a fresh
>   `git archive 9733813c` export, byte-identical to the stage tree (§1.2 C16). The executor applies, never retypes.
> - **Pre-existing test files changed: 18; deleted: 1** (§5, §9 criterion 12): four `host/store` files (five tests) whose
>   table lists, version or deadline list name the schema; fourteen downgrade-fixture files (19 sites) that must also `DROP TABLE verb_requests`
>   (a v9 store relabelled 4–7 still holding `verb_requests` is refused as "not a credentials store", C6); and
>   `surface_confirm_unrouted_test.go` deleted (S2 Appendix C).

---

## 0. Binding constraints (each is restated where it governs)

- **D-WORLD-118 = A** (verified by the controller at iteration 400, `mission_ledger_provenance.sh --id D-WORLD-118` →
  `verdict=attended`, flip `72007b37`; `mission_decisions.sh --open` at iteration 402 → no OPEN row). S3 is routable;
  **F6, `DraftOnly` and every "B +≈ 15" branch are not built** (VT2's `ready` row pins that a ready scope is admissible).
- **Three GET routes, no POST** (`TestConfirmRoutesRegistered`: exactly `GET /surface/confirm.js`, `GET
  /surface/confirm/{id}`, `GET /surface/confirmations/{id}`, each behind `approvalRoute`; POST/PUT/DELETE/OPTIONS → 405).
  **No CLI change** (no `cmd/` production file changes; seven `cmd/` test fixtures gain one DROP each). **No daemon
  process change beyond route registration and one read in `requests.json`** (no writer: `TestDaemonCredentialWriterSites`
  unchanged and green). **`surfaceCSP` unchanged** (`TestDecidePageAssets` pins its literal; the card is served under it).
  **No `.ail`.**
- **D-WORLD-109**: every store under `t.TempDir()`; keys are `approvaltest`'s or opaque test bytes; no real store is
  opened or migrated; no listener beyond `httptest` (the approval rig serves through `Handler()` with no socket).
- Coding standards: `FinishConfirmRefusal` is pure (S1); the store and the routes are host effects (S2); "why not a
  package" is design §3 S3.
- Not touched: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/`, `cmd/world-publish/*`, every
  `host/surfacepage` page and script byte (`confirm/`, `decide/`, `approve/`, `passkey/`, `app/`), `surface_page.go`,
  `approvalRoute`, the `/v1` table.

---

## 1. Pre-flight

### 1.1 Baseline at `9733813c` (measured by the planner, 2026-11-17)

`tools/gates.sh <base-tree> base-gates base` (the §6 script in `base` mode; a `git archive 9733813c` export with a
`git init`; `GATES_SKIP_VERIFY_GO=1` — iteration 401's controller ran `verify_go` on this tree's `cmd host`):

| leg | rc | time | note |
|---|---|---|---|
| `status` / `gofmt` / `build` / `vet` / `fence` | 0 / 0 / 0 / 0 / 0 | 0 / 2 / 5 / 2 / 2 s | `go test -run '^$' ./...` is the test-file compile fence |
| `verify_ail` | 0 | 206 s | no `.ail` changes in S3; the final must match |
| `pkgs` | **1** | 578 s | `host/daemon` 573.9 s, `cmd/ailang-worldd` 388.3 s, `host/authority`, `host/surfacepage` ok; **`host/store` red: `TestConvertedStoreMethodsM5b/cutoff/AppendOutcome` "post-cutoff result = store: append outcome: context deadline exceeded, want uncertain"** — a world-store context cut-off under load (≈ 25), at the base, before any S3 line (S3 touches no world-store code); `host/broker` by `^(TestFinish\|TestApprove\|TestParse\|TestRecordOperator)` ok. (An earlier base run, the first `gates.sh` draft with `host/broker` whole, was red there on the directive's `TestExecProtectedCheckTimeout/reserve` and `TestWorktreeSnapshotRefusesAConfiguredFilter` — hence §6's scoping.) |
| `race` | **1** | 310 s | `host/store` 56 s, `host/authority`, `host/surfacepage`, `host/broker` (`^TestFinish`), `host/daemon` (subset) ok; **`cmd/ailang-worldd` (`^(TestApprove\|TestApprover)`) red: `TestApproveBySignatureRefusalsBeforeAnyEffect/git-namespace` "read mint request: context deadline exceeded"** — load, at the base |
| `node` | 0 | 21 s | the nine named harnesses `--- PASS` (S2's), `WORLD_EXEC_NODE` set |
| `inv` | 0 | 64 s | **55 named `--- PASS`**, 0 SKIP, 0 FAIL (incl. `TestConfirmRoutesNotRegistered`, which S3 deletes) |
| `chrome` | **1** | 139 s | 14 of V2b's 14 + XC1 `--- PASS` except **`TestVerbBarWhyInChrome` "driver waitFor: timed out"** (load); the `host/daemon` drills then did not run (`&&`) |
| **`GATES_DONE 1331s FAILS: pkgs race chrome …`** | | | every red is a timing failure under load ≈ 20–25 from other sessions, at the base; the final tree is compared leg by leg with these (C17) |

### 1.2 Codebase claims checked (at `9733813c` unless "work" — the planner's stage tree, tag `M9b`)

| # | claim | command | observed |
|---|---|---|---|
| C1 | **Nothing the design cites for S3 moved** between the design's base and HEAD | `git diff --stat 994ba123 9733813c -- host/store/ host/daemon/ host/broker/approve.go host/surfacepage/approve_model.go`; control `git diff --stat 994ba123 9733813c -- cmd host \| tail -1` | only `host/daemon/surface_confirm_unrouted_test.go \| 33 +` (S2's, which S3 deletes); control `26 files changed, 2662 insertions(+), 20 deletions(-)` (S1, S2, S1b) |
| C2 | the store cites, re-measured | `grep -n` | `credentials.go:44 const credentialsSchemaVersion = 8`, `:320 known := version >= credentialsSchemaV1 && version <= credentialsSchemaVersion`, `:342 … "user_version %d is newer than this binary's %d; use a compatible binary"` (W53), `:375 if version < 8 {` (the ladder); `mint_requests.go:67 type OlderCredentialsStoreError`, `:49-54` the mint sentinels (texts name mint requests, W39), `:103` the per-feature refusal (W54), `:300 type ApproverUse`; `approvers.go:245 func useApprover`; `passkeys.go:225 func usePasskey`; `key_decisions.go` (schema 8, row 201 A4) `isDecisionRef`, `MaxSurfaceApprovalArtifactBytes`; `credentials_schema.sql:103 CREATE TABLE IF NOT EXISTS key_decisions` — **all unmoved** |
| C3 | the daemon cites, re-measured | `grep -n` | `daemon.go:813 ensureCredentialsStore(credPath, …)` → `credentials.go:59-66 store.UpgradeCredentials(path)` (W54: `serve` upgrades the real store at start); `surface_approve.go:76 approvalRoute` (421 before the handler), `:92 handleSurfaceRequests`, `:146 handleSurfaceRequestCard`; `approve_model.go:63-73 RequestsJSON` (`Pending`, `Decided`, `Unavailable`, `Decisions`, `DecisionsTruncated`); `surface_page.go:19 const surfaceCSP`, `:21 setSurfaceHeaders`; `surface_decide_page.go:15 handleSurfaceDecidePage` (id check → 400, `setSurfaceHeaders`, the embedded bytes), `:26` the script (`text/javascript; charset=utf-8`); routes `daemon.go:1139-1148`; `route_table_test.go:15 TestFrozenV1RouteTableMatchesMux` (`/v1` only, 13 pairs) — **all unmoved** |
| C4 | what S3 consumes from S1/S2 is at HEAD | `grep -n` | `approval/verb.go:69 Validate`, `:120 CanonicalVerb`, `:133 DecodeVerb`, `:166 VerbMessageHash(digest, challenge, expiry) ([32]byte, error)`, `:29 MaxVerbRequestBytes = 128 << 10`; `broker/finish_prompt.go:14 DisplayCommitter`, `:24 FinishPromptText`; `broker/approve.go:689 FinishApprovalScope`, `:698 FormatFinishApprovalScope`, `:704 ParseFinishApprovalScope`, `:722 RecordOperatorApproval` (returns the decision `hashref` that `decision_ref` names, S4's C14); `surfacepage/confirm_model.go:24 ConfirmationInput`, `:52 ConfirmationCard`; `confirm.go:16 ConfirmPage`, `:21 ConfirmScript` |
| C5 | nothing named verb request exists in the store or the daemon yet | `grep -rn 'verb_requests\|VerbRequestRow\|PutVerbRequest' cmd host --include='*.go' \| wc -l`; control `grep -rn key_decisions host/store/credentials_schema.sql \| head -1` | `0`; control `103:CREATE TABLE IF NOT EXISTS key_decisions (` |
| C6 | **the schema bump breaks every downgrade fixture that does not drop the new table** | a scratch clone at `9733813c` with only M7a's three production files, then `go test ./host/daemon/ -run '^TestRequestsJSON$'`, `./host/authority/ -run '^TestJudgeVerbSSH$'`, `./host/store/` | `--- FAIL: TestRequestsJSON … a schema-6 store: {… "unavailable":"the mint requests could not be read" …}`; `--- FAIL: TestJudgeVerbSSH … user_version 4 with tables [episode_starts mint_requests session_credentials verb_requests] is not a credentials store`; `--- FAIL: TestCredentialsUpgradeToV6`, `…SchemaMatchesCanonicalDDL`, `…Schema8KeyDecisions`, `…UpgradeToV4`. `grep -rn 'DROP TABLE key_decisions' --include='*_test.go' cmd host` → **19 sites in 14 files** (the A4 v8 bump, `3a4026e9`, edited the same fixtures for `key_decisions`) |
| C7 | the closest landed analogues | `git log --oneline -S keyDecisionsSQL -- host/store/`; `git show --stat 3a4026e9` | the v8 bump is **row 201 A4 sprint 2 M4e** (`3a4026e9`, plan `w-world-approvals-by-key-a4-sprint2-plan.md`): `credentials.go` (+version, `credentialsTables`, the `< 8` line), the schema file, `key_decisions.go` + `TestCredentialsSchema8KeyDecisions` (a v7 store built by `makeV6Credentials` + `surfaceApprovalsSQL`, the per-method refusal at 7, `UpgradeCredentials`, the canonical DDL) and the downgrade fixtures. **Mirrored**: its upgrade test's construction (`makeV6Credentials` + the v7 and v8 DDL → a real v8 store), its refusal assertion (`OlderCredentialsStoreError` Version/Need + the restart text), its fixture edits; **strengthened**: CS1 compares `pragma_table_info` of every table and the whole `verb_requests` DDL of a fresh v9 store against the upgraded one (not a canonical string). The decide page routes (`surface_decide_page.go` + `TestDecidePageAssets`, `TestDecisionRoutesRegistered`) are mirrored by M9a's handlers and `TestConfirmPageAssetsServed`/`TestConfirmRoutesRegistered`; A3's `TestRequestsJSON`/`TestRequestCardJSON` by CR1/CR2 |
| C8 | `host/surfacepage` cannot parse the finish scope (so the `requests.json` summary is decoded in `host/daemon`) | `go list -f '{{.Imports}}' ./host/surfacepage ./host/broker` | surfacepage: `[… host/approval host/surface …]` (no broker, no store); broker imports `host/store` |
| C9 | the "newer than this binary's" branch (W53) has **no test** at HEAD | `grep -rn 'newer than this binary' --include='*.go' cmd host` | only `host/store/credentials.go:342` (production); CS7(a) adds the first test |
| C10 | **the design's CS5 "a `decision_ref` on a `superseded` row" is not refused by the design's own DDL** | `TestVerbRequestsChecks` first run, row `superseded-with-decision` | `err <nil>, want ok false` — `(outcome IS 'approved') = (approved_by IS NOT NULL AND artifact IS NOT NULL AND decision_ref IS NOT NULL)` is `false = false` when only `decision_ref` is set (§1.4 F1) |
| C11 | the new code adds no context root, no deadline-free read, no non-test Go file outside `host/`/`cmd/`, no credentials writer in the daemon | `inv` (§6) on the final tree: `TestProductionContextRoots`, `TestNoNewDeadlineFreeStoreReads`, `TestProductionGoSurface`, `TestDaemonCredentialWriterSites` | all `--- PASS`, unedited (C17); the read routes use `d.readCtx(r)` (existing) and `OpenCredentialsReadOnly` only |
| C12 | `requests.json`'s existing pins survive two new trailing keys | `surface_approve_read_test.go:76` asserts the substring `"pending":[],"decided":[],"unavailable":"the mint requests could not be read"` | the new fields follow `decisionsTruncated`; `TestRequestsJSON` green unedited except its two DROP fixtures |
| C13 | the new tests, work | `go test -count=1 -v` by name (`accept.sh`) | all `--- PASS`, none SKIP |
| C14 | **every sub-milestone is green when committed** | `tools/accept_all.sh M7a:fast M7b:fast M8:fast M9a:fast M9b:fast M7a:daemon M9a:daemon M9b:daemon M7a:cmd` (an APFS clone of `work` at each tag running §3's `accept.sh`) | on the stages before R4's two CS3 lines (`logs/accept_all.txt`): `M7a fast rc=0 ACCEPT_OK` (97 s), `M7b fast rc=0` (76 s), `M8 fast rc=0` (36 s), `M9a fast rc=0` (37 s), **`M9b fast rc=1`** (102 s) on `host/store`'s **`TestConvertedStoreMethodsM5a`** ("post-cutoff result … context deadline exceeded, want uncertain": a world-store cut-off under load; its sibling `…M5b` is red the same way at the base, §1.1) — green alone (`-count=3`, 5.6 s) and on the part's re-run (`M9b fast rc=0`, 20 s); `M7a daemon rc=0` (694 s), `M9a daemon rc=0` (632 s), `M9b daemon rc=0` (582 s), `M7a cmd rc=0` (525 s) — each ≤ 600 s alone, at load ≈ 20 they ran 525–694 s **beside** the planner's mutant and gate runs (the executor runs one at a time). **Re-run on the final stages** after R4 (only M7b's patch changed: +2 lines in `verb_requests_test.go`; the `daemon`/`cmd` parts run no `host/store` test): `M7b fast rc=0 ACCEPT_OK` (32 s), `M9b fast rc=0 ACCEPT_OK` (19 s). An earlier draft of `accept.sh` ran `host/broker` whole in M8's part and went red on the directive's exec-probe flakies after 480 s (`logs/accept_all_run1_aborted.txt`): M8 now accepts `^TestFinish` (§6) |
| C15 | the snapshot block (§7) on every sub-milestone, incl. the deletion | `tools/snap_rehearse.sh` (a local clone at `9733813c`, the patches in order, §7's block verbatim) | `M7a 21/0/21`, `M7b 22/0/22`, `M8 24/0/24`, `M9a 28/1/27 [host/daemon/surface_confirm_unrouted_test.go]`, `M9b 30/1/29`, each `MANIFEST_OK`, dirs `cmd host` (no `docs`) |
| C16 | **the round trip** | `tools/roundtrip.sh <this plan>` | `CMP_EQUAL` for all five (each block extracted from this file with §2's awk equals the generated patch byte for byte); each applied with `/usr/bin/patch -p1 -E` after a `--dry-run` to a fresh `git archive 9733813c` export, no `.orig`/`.rej`, committed and measured: `NUMSTAT M7a 377`, `M7b 805`, `M8 913`, `M9a 1235`, `M9b 1391` (each `unfiltered=` the same); `DELETED_OK host/daemon/surface_confirm_unrouted_test.go`; then `diff -r -x .git` against the stage tree (`git archive M9b`) → **`ROUNDTRIP_IDENTICAL_TO_STAGE_TREE`** (`logs/roundtrip_last.txt`, re-run on this final plan file) |
| C17 | **the gates on the final tree** | `tools/gates.sh <final-tree> final-gates final <sha>` | `tools/gates.sh <final-tree> final-gates final b1d590e9` (`final-tree` = a clone of the repo at `9733813c`, the five patches applied with `patch -p1 -E` and committed; `diff -r` equal to the stage tree; `GATES_SKIP_VERIFY_GO=1`; load ≈ 20 from other sessions): **`TOUCHED_PACKAGES: cmd/ailang-worldd host/authority host/broker host/daemon host/store host/surfacepage`** (= the plan's list); `status` 0; **`verify_ail` rc 0** (`✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` — the base's line, byte for byte); `gofmt`, `build`, `vet`, `fence` rc 0; **`pkgs` rc 1**: `host/store` 38.8 s, `host/authority`, `host/surfacepage` ok, `host/daemon` red only on **`TestMarketingCaptureInventoryMutation`** (a 30 s marketing timeout panic) and `cmd/ailang-worldd` only on **`TestCLIRealSubprocessEpisode`** ("build subprocess binary: signal: killed", the directive's list) — **both `--- PASS` alone on the same tree** (3.27 s, 2.34 s; `logs/final_rerun_*.txt`), and the M7a/M9a/M9b `daemon` and M7a `cmd` acceptance parts ran both packages whole green (C14); **`race` rc 0** (`host/store` 77 s, `host/authority`, `host/surfacepage`, `host/broker` `^TestFinish`, `host/daemon` and `cmd/ailang-worldd` subsets — the base's `race` was red on `cmd`, §1.1); **`race3` rc 0**; **`node` rc 0** (nine named `--- PASS`); **`inv` rc 0: 65 named `--- PASS`** (the base's 55 − `TestConfirmRoutesNotRegistered` + `TestConfirmRoutesRegistered`, `TestConfirmPageAssetsServed`, `TestConfirmationCardServed`, `TestRequestsJSONConfirmations`, `TestFinishConfirmRefusal` and CS1–CS7's six), 0 SKIP, 0 FAIL, none missing, `TestConfirmRoutesNotRegistered` absent; **`chrome` rc 0** (24 `--- PASS` lines: V2b's fourteen, XC1, the two `host/daemon` drills — the base's `chrome` was red on `TestVerbBarWhyInChrome` under load); **`GATES_DONE 1462s FAILS: pkgs`** (the two load reds above). **Delta on the final tree `d6ac82f8`** (after R4: +2 lines in `verb_requests_test.go`, nothing else; `final-gates-delta/delta.txt`): `gofmt` clean, `build`, `vet ./host/store/`, the fence over `./host/... ./cmd/...` clean; `go test ./host/store/` ok 11.4 s; `-race ./host/store/` ok 35.3 s; `race3` ok; the 19 store inventories and schema tests by name `--- PASS` (incl. `TestProductionContextRoots`, `TestProductionGoSurface`, `TestNoNewDeadlineFreeStoreReads`). `verify_go` is the controller's |
| C18 | node | `/opt/homebrew/bin/node --version` | `v26.13.0` (the `node` leg's harnesses are S2's, unchanged) |

### 1.3 Planner decisions where the design or the directive leaves a choice

- **S3-D1 — no cut; M7, M8 and M9 ship, by measurement.** M7 compiled first: **836** (M7's design ≈ 570 raw: 1.47×,
  above the V-sprints' 1.35×: the 19 fixture sites and five store-test edits are not in the design's estimate). M8 108,
  M9 464 → **1,408 at first compile**. The directive's rule (ship all if ≤ 1,400) then needed −8: comment lines in
  `verb_requests.go` and `surface_confirm_read.go`, and one context test that duplicated
  `TestCredentialsMethodsRequireADeadline`, came out (1,408 → 1,387); the mutant enumeration then changed the tests
  (CS6's boundaries and its consume-then-put check, CR2's `Validate` and scope rows, the 421 lines, the hand-built CR1
  input, a tighter `vrState`/`tableInfo`, CS3's prune-keeps-consumed check) → **1,391**. The
  cut the directive names (S3a = M7 + M8 = 913; S3b = M9 = 478) was measured and **not taken**: every part is proven
  here (§4), and a store with no route is no safer than a store whose only reader is these three GETs.
- **S3-D2 — five sub-milestones.** M7 split by kind, not by feature, so each is green: **M7a** (production + every
  pre-existing test the bump turns red, 377) and **M7b** (the new store tests, 428; one new file + the deadline list).
  M9 split at the seam the directive fixes: **M9a** registers the three GETs **and deletes
  `surface_confirm_unrouted_test.go` in the same patch** (its `TestConfirmRoutesNotRegistered` would be red the moment a
  route exists), with CR1, CR3 and CR4; **M9b** adds `requests.json`'s list and CR2. M7b (428) is the one sub-milestone
  above ≈ 400: a single test file, accepted in ≈ 1–2 min.
- **S3-D3 — audit note 1: three GETs, not four.** §8 says "the four GET registrations"; §2.5 defines three
  (`/surface/confirm/{id}`, `/surface/confirm.js`, `/surface/confirmations/{id}`); the fourth registration in §7's "four
  registrations" is S4's `POST /surface/confirm/{id}`. Proven by `TestConfirmRoutesRegistered` (route table: exactly the
  three, each behind `approvalRoute`; every other method 405 with no CORS header) and its mutants `RTE-POST`,
  `RTE-PAGE-ANY` (§4, class R).
- **S3-D4 — audit note 2: four CR rows, by measurement.** §6 "Sprint 3" names CR1 (the data route) and CR2
  (`requests.json`); §2.9 P6 says "P2/P3 (read routes) are §2.5's rows, CR1–CR4". §2.5 has exactly four read-route rows:
  the data GET, `requests.json`, the page + script GETs, and the route table they share. S3 builds all four:
  **CR1** `TestConfirmationCardServed`, **CR2** `TestRequestsJSONConfirmations`, **CR3** `TestConfirmPageAssetsServed`
  (PD2-style: the served bytes = the files, `surfaceCSP` and the surface headers, 400, 421 — S2 Appendix C's "assets
  test"), **CR4** `TestConfirmRoutesRegistered` (the route table, no POST). `TestDecidePageAssets` stays unedited.
- **S3-D5 — audit note 3 (a raced supersession between C14 and C15): S3 owns the store half, and pins both orders.**
  (a) *put-then-consume*: a newer `episode finish` supersedes the pending row after the daemon's C14 world write; then
  C15's `ConsumeVerbRequest` refuses with `ErrVerbRequestSuperseded` **and writes nothing** (CS6 `superseded`, with the
  rows, the approver's `last_used_at` and the counter compared before and after). (b) *consume-then-put*: a
  `PutVerbRequest` for the same episode after a confirmation leaves the approved row byte-identical (`outcome`,
  `approved_by`, `artifact`, `decision_ref`), because the supersede is `WHERE episode_id = ? AND consumed_at IS NULL`
  (CS6 `confirm-*`, mutant `PUT-SUPERSEDE-PENDING`). SQLite serialises the two transactions, so there is no third order.
  **S4's**: mapping (a) to the POST's 409 ("superseded by a newer finish of <episode>"), the orphaned world decision it
  leaves (inert, W62 — the same state as the window table's "after C14, C15 fails" row), and the window-table row the
  design lacks (a CX test that supersedes between C14 and C15 through the C15 seam).
- **S3-D6 — CR1's expected body is `surfacepage.ConfirmationCard` of the row, input built by hand.** The test reads
  the stored row (`VerbRequest`), builds `ConfirmationInput` field by field from it (not through the production
  `confirmationInput`, so a mapping mutant cannot cancel out), and compares `json.Marshal(ConfirmationCard(…))` with the
  served body byte for byte, for a pending, an approved and a superseded row; then it asserts no message hash (hex,
  base64url, base64 of Go's `VerbMessageHash`), no artifact (raw and base64), no `"artifact"`/`"hash"` key.
- **S3-D7 — the statuses.** `GET /surface/confirmations/{id}`: bad id 400 `BadRequest`; below schema 9 **503 `Restart`**
  naming the serve restart (design §2.5); unknown 404; store unreadable 500 `Internal` (A3's card pattern, `RT3`).
  `requests.json`: below 9 → `"confirmations":[]` + `confirmationsUnavailable` naming the restart; unreadable → `[]` +
  "the verb requests could not be read"; the other lists untouched in both (CR2).
- **S3-D8 — the summary is decoded in `host/daemon`** (C8): `pendingConfirmation` refuses, in order, bytes that do not
  hash to the digest, bytes that `approval.DecodeVerb` or `Validate` refuse, and a scope `broker.ParseFinishApprovalScope`
  refuses; otherwise `{id, episode, repository, head, sha, draft, expires}` (design §2.5). `surfacepage` gains only the
  wire type `PendingConfirmation` and the two `RequestsJSON` fields.
- **S3-D9 — `ConsumeVerbRequest`'s order.** Grammar (`ErrBadVerbRequest`, new: id, digest hex, `approvedBy` =
  `format:name` of the use, passkey use iff webauthn, artifact 1–16 384, `decisionRef` sha256 grammar), then one
  transaction: read the row (unknown), superseded, consumed ("… by <approved_by>", design C8's text), clock behind,
  expired, digest (the caller's and the stored bytes'), `useApprover`, `usePasskey`, one UPDATE. Checks precede the
  UPDATE (the key-decision pattern) rather than following a guarded UPDATE (the mint pattern) because every check is in
  the same transaction; the UPDATE keeps `AND consumed_at IS NULL` and the RowsAffected check as defence in depth (both
  named equivalent, §4).
- **S3-D10 — `store.MaxVerbRequestBytes`** (= 131 072, the CHECK's bound) is a store constant because `host/store` does
  not import `host/approval`; CS1 asserts it equals `approval.MaxVerbRequestBytes`.
- **S3-D11 — store fixtures are `approval.CanonicalVerb` bytes** (`vrRow`, `Validate`d), daemon fixtures too, with
  `Prompt = broker.FinishPromptText(…)` and `Scope = broker.FormatFinishApprovalScope(…)` (`verbFixture`); VT2's
  fixtures are `FormatFinishApprovalScope` + `FinishPromptText` of one scope. Nothing is hand-typed JSON.
- **S3-D12 — timing.** No S3 test waits on asynchronous work (no goroutine, harness or poll): every store call is
  synchronous under a deadline context; the daemon tests call `Handler().ServeHTTP` directly. Directive 11's SQLite
  half is the `race3` gate leg (`-race -count=3` over the new store tests, §6).

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

| # | premise | measured | routing |
|---|---|---|---|
| F1 | design §6 CS5: "the table's CHECKs refuse … a `decision_ref` on a `superseded` row" | the design's DDL accepts it (C10) | **DDL kept verbatim** (no re-design). CS5 pins the actual behaviour (`superseded-with-decision-accepted`, ok) with a comment, adds `superseded-with-approver` (refused), and the Go side is the fence: only `ConsumeVerbRequest` writes `decision_ref`, with `outcome = 'approved'` in the same UPDATE (CS6). **For the controller**: a v10 CHECK `(decision_ref IS NULL) OR (outcome IS 'approved')` would close it; not S3's to add |
| F2 | design §6 CS7: "a store at `user_version` 9 opened with the constant set to 8 refuses" | one binary cannot hold two constants; the branch (`version > credentialsSchemaVersion`, `credentials.go:341`) is the same for any newer version | CS7(a) opens a **v10** store under the v9 build (writable, read-only and `UpgradeCredentials`): the W53 text "user_version 10 is newer than this binary's 9", nothing changed; CS7(b) (a v8 store under v9 refuses each feature with `Need: 9`) as designed |
| F3 | design §7 "Pre-existing tests edited" for S3: "the canonical-DDL test (v9)" and "`host/store/context_roots_test.go` (two roots)" | S3 adds no context root (the two roots are S5's CLI poll, design §2.7); the bump forces **19 fixture edits in 14 files** (C6) and edits to five store tests | `context_roots_test.go` unedited; the fixture and store-test edits listed in §5/§9.12, each one line (or a list/version) |
| F4 | design §8 "the four GET registrations" | §2.5 defines three (audit note 1) | three (S3-D3) |
| F5 | design §2.9 P6 "CR1–CR4" vs §6 "CR1–CR2" | §2.5 has four read-route rows (audit note 2) | four test rows (S3-D4) |
| F6 | design §2.4 sentinels `ErrVerbRequestUnknown\|Consumed\|Superseded\|Expired\|ClockBehind\|DigestChanged`, `ErrTooManyPendingVerbRequests` | the grammar refusals need one more, as the mint and key-decision tables have (`ErrBadMintRequest`, `ErrBadKeyDecision`) | `ErrBadVerbRequest` (new) |
| F7 | design §7 places `FinishConfirmRefusal` in `host/broker/approve.go` "or a new `finish_prompt.go`" | `finish_prompt.go` exists (S1) and holds the prompt only | a new `host/broker/finish_confirm.go` (+ test) beside it — the one production file outside design §7's names (a new file in a named package) |
| F8 | S2 Appendix C: "adds the confirm page's bytes to an assets test (`TestDecidePageAssets`-style, PD2)" | editing `TestDecidePageAssets` would add a pre-existing test edit | a new `TestConfirmPageAssetsServed` (CR3) in the new test file; `TestDecidePageAssets` unedited and still pins `surfaceCSP` |

**Files outside design §7**: `host/broker/finish_confirm.go` (F7) and its test; `host/daemon/surface_confirm_read_test.go`
is §7's "(+ tests)". Test files edited outside §7's list: the 14 fixture files and 4 of the 5 store tests (F3).

### 1.5 The size

| stage | production | test | changed lines | cumulative (numstat after commit) |
|---|---|---|---|---|
| M7a | 325 (`verb_requests.go` 293, schema 20, `credentials.go` 12) | 52 (19 fixture sites, 38 lines; 4 store-test files, 14) | 377 | **377** |
| M7b | 0 | 428 (`verb_requests_test.go` 424, `credentials_test.go` 4) | 428 | **805** |
| M8 | 41 | 67 | 108 | **913** |
| M9a | 86 (`surface_confirm_read.go` 79, `daemon.go` 7) | 236 (`surface_confirm_read_test.go` 203, −33 deleted) | 322 | **1,235** |
| M9b | 74 (`surface_confirm_read.go` +48, `approve_model.go` 17, `surface_approve.go` 9) | 82 | 156 | **1,391** |

Design ≈ 980 raw → compiled **1,391 (1.42×)**: M7 805 vs ≈ 570 (the fixtures), M8 108 vs ≈ 90, M9 478 vs ≈ 320. Cap
1,400, margin 9 — any growth on the executor's side is a stop (§7).

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of four backticks,
in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; a new file's patch has `--- /dev/null`; **M9a deletes
`host/daemon/surface_confirm_unrouted_test.go`** (`+++ /dev/null`, `@@ -1,33 +0,0 @@`), which `patch -E` removes (without
`-E` it would stay as a 0-byte file). Apply it **only** this way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF2'
set -euo pipefail
M=M7a   # <- the sub-milestone, exactly as in the Appendix P heading, in order: M7a M7b M8 M9a M9b
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter402-row243-v3s3
PLAN="$W/design_docs/planned/w-world-address-bar-verbs-v3s3-sprint-plan.md"
T="$HOME/.ailang/state/world-iter402/exec-patches"; mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
/usr/bin/patch -p1 -E --dry-run < "$T/$M.patch"
/usr/bin/patch -p1 -E < "$T/$M.patch"
if find cmd host docs \( -name '*.orig' -o -name '*.rej' \) | grep .; then echo "STOP: patch left .orig/.rej"; exit 1; fi
if [ "$M" = M9a ] || [ "$M" = M9b ]; then test ! -e host/daemon/surface_confirm_unrouted_test.go || { echo "STOP: the unrouted test survived"; exit 1; }; fi
git status --short --untracked-files=all
EOF2
```

- If `--dry-run` fails, **stop and report** (no hand edits, no fuzz options). It can only fail if an earlier
  sub-milestone was skipped, reordered or changed.
- Then run that sub-milestone's **Acceptance** (§3): extract `accept.sh` from this plan (`awk '$0 == "````bash
  accept.sh" {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/accept.sh"`) and run **each part §3 lists for M, one
  per Bash call**: `bash "$T/accept.sh" "$W" <M> <part> 2>&1 | tee
  "$HOME/.ailang/state/world-iter402/exec-logs/accept_<M>_<part>.txt"` (never only through `tail`). **Never `nohup`,
  never background a test leg; pass a 600000 ms timeout to every long Bash call.**
- Then snapshot (§7). The controller commits and compares the numstat with §3's expected value.

---

## 3. Milestones

**Acceptance.** Every part runs the same `[base]` checks first (gofmt, build, vet of the six touched packages, the
test-file compile fence `go test -run '^$' ./...`). `v` requires every named test to `--- PASS` and none to SKIP or FAIL;
`p` is a whole-package run. Parts: `fast` (named tests + the small packages; ≈ 1–2 min), `daemon` (`host/daemon` whole,
≈ 5–8 min at load), `cmd` (`cmd/ailang-worldd` whole, ≈ 3–6 min). Each part is one 600 s call; measured in C14.

````bash accept.sh
#!/bin/bash
# accept.sh <tree> <M> [fast|daemon|cmd]: exactly plan §3's acceptance for sub-milestone M (one bash shell;
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
step go vet ./host/store/ ./host/broker/ ./host/daemon/ ./host/surfacepage/ ./host/authority/ ./cmd/ailang-worldd/
if go test -run '^$' ./... 2>&1 | grep -E '^(FAIL|--- FAIL)'; then echo "ACCEPT_RED $M $P: FENCE_RED"; exit 1; fi
SCHEMA='TestCredentialsSchemaMatchesCanonicalDDL|TestCredentialsSchema8KeyDecisions|TestCredentialsUpgradeToV4|TestCredentialsUpgradeToV6|TestCredentialsSchema7SurfaceApprovals|TestCredentialsMethodsRequireADeadline'
CS='TestCredentialsSchema9VerbRequests|TestCredentialsNewerThanNineIsRefused|TestPutVerbRequestSupersedes|TestPutVerbRequestGrammarCapPrune|TestVerbRequestsChecks|TestConsumeVerbRequest'
ROUTES='TestApprovalRoutesRegistered|TestDecisionRoutesRegistered|TestDecidePageAssets|TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestDaemonCredentialWriterSites|TestDaemonDependencyAllowlist'
CR1='TestConfirmRoutesRegistered|TestConfirmPageAssetsServed|TestConfirmationCardServed'
CR2='TestRequestsJSONConfirmations|TestRequestsJSON|TestRequestCardJSON|TestSurfaceRequestsListsDecisions|TestSurfaceRequestsDecisionsTruncated|TestSurfaceRequestsDecisionsUnreadable'
case "$M:$P" in
M7a:fast) step v ./host/store/ -run "^($SCHEMA)\$"
          step p ./host/store/ ./host/authority/ ./host/approval/... ;;
M7a:daemon|M9a:daemon|M9b:daemon) step p -timeout 30m ./host/daemon/ ;;
M7a:cmd)  step p -timeout 30m ./cmd/ailang-worldd/ ;;
M7b:fast) step v ./host/store/ -run "^($SCHEMA|$CS)\$"
          step p ./host/store/ ;;
M8:fast)  step v ./host/broker/ -run '^(TestFinishConfirmRefusal|TestFinishPromptText|TestFinishApprovalScopeStrict)$'
          step v ./host/store/ -run "^($CS)\$"
          step v ./host/broker/ -run '^TestFinish' ;;   # host/broker whole is 400+ s with the directive's exec-probe flakies
M9a:fast) step v ./host/daemon/ -run "^($ROUTES|$CR1)\$"
          step v ./host/surfacepage/ -run '^(TestConfirmationCard|TestConfirmationCardWire|TestConfirmCardGolden|TestConfirmPageAssets|TestConfirmPageFieldTable)$'
          step v ./host/broker/ -run '^TestFinishConfirmRefusal$'
          step v ./host/store/ -run "^($CS)\$" ;;
M9b:fast) step v ./host/daemon/ -run "^($ROUTES|$CR1|$CR2)\$"
          step p ./host/surfacepage/ ./host/store/ ;;
*) echo "unknown $M $P"; exit 2 ;;
esac
echo "ACCEPT_OK $M $P"
````

### M7a — schema 9 and the store (design §2.4) · 325 prod + 52 test · Run 1 · parts: `fast`, `daemon`, `cmd`

- `host/store/credentials_schema.sql` (+20): design §2.4's `verb_requests` DDL verbatim (`CREATE TABLE IF NOT EXISTS`).
- `host/store/credentials.go`: `credentialsSchemaVersion = 9` (comment names row 243 V3); `credentialsTables`: a
  `version >= 9` set with `verb_requests`, the old set kept as `version == 8`; `upgradeCredentials`: `if version < 9 {
  … verbRequestsSQL }`.
- `host/store/verb_requests.go` (NEW): `verbRequestsSQL` (the same DDL as `CREATE TABLE`), `verbRequestsSchemaVersion =
  9`, `MaxPendingVerbRequests = 16`, `MaxVerbRequestBytes = 128 << 10`, `verbRequestPruneAfter = 86400`,
  `VerbOutcomeApproved|Superseded`, `VerbRequestIDPattern = ^vr-[0-9a-f]{16}$`, the sentinels (§1.4 F6),
  `VerbRequestRow`, `PutVerbRequest` (S1 deadline + v9 → S2 grammar (id, size, digest = SHA-256 of the bytes — hence 64
  lowercase hex —, challenge hex, episode, expiry, pending) → S3 supersede → S4 prune → S5 cap → S6 insert, one
  transaction), `VerbRequest`, `PendingVerbRequests` (newest first, clamp 16), `ConsumeVerbRequest` (S3-D9).
- Pre-existing tests (each: the bump turns it red, C6): `credentials_test.go` (fresh table list + `verb_requests`;
  version 8 → 9), `approvers_test.go`, `mint_requests_test.go` (upgraded table lists), `key_decisions_test.go` (the v7
  store now upgrades to `credentialsSchemaVersion`, tables `credentialsTables(credentialsSchemaVersion)`); and the 14
  downgrade fixtures: `` `DROP TABLE verb_requests`, `` before each `` `DROP TABLE key_decisions` `` (cmd: 7 files, 8
  sites; host/authority: 4 files, 6 sites incl. S1b's `judge_verb_test.go`; host/daemon: 3 files, 5 sites).
- **Expected numstat after commit: 377.**

### M7b — CS1–CS7 (design §6 Sprint 3) · 0 prod + 428 test · Run 1 · part: `fast`

- `host/store/verb_requests_test.go` (NEW): `vrRow` (CanonicalVerb bytes), `vrStore` (ssh `mark-ssh` + passkey `x`,
  counter 5), `vrState` (every `verb_requests` row, approvers' `last_used_at`/`revoked_at`, passkey counters — the
  "writes nothing" witness), `tableInfo`; **`TestCredentialsSchema9VerbRequests`** (CS1, CS2, CS7(b)),
  **`TestCredentialsNewerThanNineIsRefused`** (CS7(a)), **`TestPutVerbRequestSupersedes`** (CS3; and a put a day after the superseded row's expiry keeps it — the prune takes
  pending rows only),
  **`TestPutVerbRequestGrammarCapPrune`** (CS4: 15 grammar rows, the exact-128 KiB control (`approval.MaxVerbRequestBytes`), prune at exactly a day and a
  second short of it, the 16 cap refusing and writing nothing, a supersede freeing its own place),
  **`TestVerbRequestsChecks`** (CS5: 19 raw-INSERT rows), **`TestConsumeVerbRequest`** (CS6: 17 refusals each comparing
  `vrState` before/after; `confirm-ssh` at the creation second, `confirm-webauthn` one second before expiry; five
  columns; counter 5 → 6; `last_used_at`; then the consume-then-put check, S3-D5(b)).
- `host/store/credentials_test.go`: `TestCredentialsMethodsRequireADeadline` lists the four new methods.
- **Expected numstat after commit: 805.**

### M8 — `FinishConfirmRefusal` (design §2.3, §6 VT2) · 41 prod + 67 test · Run 1 · part: `fast`

- `host/broker/finish_confirm.go` (NEW): `FinishConfirmFacts{Scope, Title, Body, Prompt, Committer string; Protected
  []string; RequestExpires, Now int64}`; `FinishConfirmRefusal(f) string` — F1 `ParseFinishApprovalScope`, F2
  `RequestExpires < scope.ExpiresAt ≤ Now + 3600`, F3 title SHA-256, F4 body SHA-256, F5 `Prompt ==
  FinishPromptText(scope, Title, Body, Protected, Committer)`; the design's reasons verbatim.
- `host/broker/finish_confirm_test.go` (NEW; accepted by `^TestFinish`, not the whole package — §6): **`TestFinishConfirmRefusal`** — admissible, `ready` (A: admissible), F1 ×2,
  F2 at `scope == request`, `scope < request`, exactly an hour ahead (admissible), an hour + 1 s, F3, F4, F5 ×4 (prompt
  edited, protected path not shown, committer not shown, prompt of another scope).
- **Expected numstat after commit: 913.**

### M9a — the three GET routes (design §2.5; CR1, CR3, CR4) · 86 prod + 236 test · Run 1 · parts: `fast`, `daemon`

- `host/daemon/daemon.go`: three registrations after `GET /surface/decide.js`, each `d.approvalRoute(…)`; the Handler
  comment names them (and "no confirmation POST before sprint 4").
- `host/daemon/surface_confirm_read.go` (NEW): `handleSurfaceConfirmPage` (id ≠ `VerbRequestIDPattern` → 400;
  `setSurfaceHeaders`; `ConfirmPage`), `handleSurfaceConfirmScript` (+ `text/javascript; charset=utf-8`),
  `handleSurfaceConfirmation` (S3-D7; passkeys grouped as `handleSurfaceRequestCard` groups them), `confirmationInput`.
- `host/daemon/surface_confirm_read_test.go` (NEW): `putVerb`/`verbFixture`/`canonVerb`/`putVerbBytes`/`readVerb`;
  **`TestConfirmRoutesRegistered`** (CR4), **`TestConfirmPageAssetsServed`** (CR3), **`TestConfirmationCardServed`**
  (CR1, S3-D6).
- `host/daemon/surface_confirm_unrouted_test.go` **deleted** (S2's `TestConfirmRoutesNotRegistered`).
- **Expected numstat after commit: 1,235.**

### M9b — `requests.json`'s confirmations (design §2.5; CR2) · 74 prod + 82 test · Run 1 · parts: `fast`, `daemon`

- `host/surfacepage/approve_model.go`: `RequestsJSON.Confirmations []PendingConfirmation` (`"confirmations"`) and
  `ConfirmationsUnavailable string`; the type `PendingConfirmation{ID, Episode, Repository, Head, SHA, Draft *bool,
  Expires, Refused}` (summary fields `omitempty`; a refused row carries only id, expires, refused).
- `host/daemon/surface_approve.go` (`handleSurfaceRequests`): `readConfirmations` after the decisions; older store →
  `[]` + the restart text; other error → logged, `[]` + "the verb requests could not be read"; `nil` → `[]`.
- `host/daemon/surface_confirm_read.go`: `readConfirmations`, `pendingConfirmation` (S3-D8).
- `host/daemon/surface_confirm_read_test.go`: **`TestRequestsJSONConfirmations`** (CR2: six listed rows newest first —
  a scope-less row, an invalid `episode start` row, a hand-edited row, an undecodable row, a draft and a ready summary —
  and the superseded, approved and expired rows omitted; the mint list intact; no bytes, challenge, digest or artifact;
  v8 → `[]` + `confirmationsUnavailable`; missing store → `[]` + "could not be read"), `setVerbColumn`.
- **Expected numstat after commit: 1,391.**

### 3.1 Every input and failure class S3's code reaches, traced (design §2.9; directive retro 4)

| design §2.9 row / input | S3 branch | test (row) |
|---|---|---|
| P1 "store below v9" (the CLI's Put, CF8's store half) | `PutVerbRequest` S1 → `OlderCredentialsStoreError{Need: 9}` | `TestCredentialsSchema9VerbRequests` (`put`) |
| P1 "16 pending" | S5 → `ErrTooManyPendingVerbRequests`, nothing written | `TestPutVerbRequestGrammarCapPrune` |
| P1 "request over 128 KiB" (the store's bound behind §2.7 step 1) | S2 size → `ErrBadVerbRequest`; CHECK backstop | `…GrammarCapPrune` (`request-128KiB+1`, the 128 KiB control); `TestVerbRequestsChecks` |
| P1 "writer busy" | the existing non-waiting writer lock (`OpenCredentials`), unchanged | existing `writer_lock_test.go`; CF8 is S5's |
| P1 "superseded by a newer run", "re-run after any refusal" (store half) | S3 supersede `WHERE episode_id = ? AND consumed_at IS NULL` | `TestPutVerbRequestSupersedes`; CS6 `superseded`, consume-then-put |
| P5 C7 "the row exists" | `VerbRequest` ok=false | CR1 `unknown` 404 (the route); CS6 `unknown` (the consume) |
| P5 C8 consumed / superseded / expired / clock / digest | `ConsumeVerbRequest`'s refusals (the store's re-check inside C15) | CS6 `consumed` (text names the approver), `superseded`, `clock-behind`, `expired`, `digest-other`, `bytes-edited` |
| P5 C10 F1–F5 | `FinishConfirmRefusal` | `TestFinishConfirmRefusal` (F1 ×2, F2 ×4, F3, F4, F5 ×4) |
| P5 C13 "below v9 → 503 Restart" (store half) | `ConsumeVerbRequest` version → `OlderCredentialsStoreError` | `TestCredentialsSchema9VerbRequests` (`consume`) |
| P5 C15 "approver revoked / counter" (T11) | `useApprover` / `usePasskey` in the consume's transaction | CS6 `approver-revoked`, `approver-key-changed`, `counter-not-advancing` |
| P5 C15 "a concurrent consume lost" | the re-read inside the transaction (consumed) | CS6 `consumed` |
| P5 crash window "C15 fails → row pending, nothing consumed" (store half of CX21) | every refusal rolls back | CS6 `vrState` before = after for all 17 refusals |
| P5 "`decision_ref` in the same statement" (CS6's kill) | one UPDATE of five columns | CS6 `confirm-*`; `CON-DECISION-SPLIT` (§4) |
| P2/P3 `GET /surface/confirmations/{id}`: bad id / unknown / v8 / no store / 200 | S3-D7 | CR1 (400 ×2, 404, 503, 500, 200 ×3 states, 421) |
| P2/P3 the page and the script | id check, headers, bytes | CR3 (200 ×3, 400 ×4, 421 ×3) |
| P2/P3 the route table, no POST, no CORS | `daemon.go` | CR4 |
| P2/P3 `requests.json` `confirmations`: pending, refused (digest / decode / `Validate` / scope), omitted (approved, superseded, expired), v8, unreadable | `readConfirmations`, `pendingConfirmation` | CR2 |
| P4 the card's fetch of a non-2xx (CP1) | the route now answers 400/404/503 as `ApproveRefusal` JSON | CR1 (S2's `TestConfirmPageFailsClosed` already maps any non-200 to `failed`) |
| version skew (RV7, W53) | `enforceSchema`'s newer-than branch at 9 | CS7(a) |
| D-WORLD-118 = A (no F6) | no draft rule | VT2 `ready` row |

---

## 4. Mutation table

**Procedure** (`tools/mutate.py <run> <workers> [ids…]`, rows in `tools/mutants.py`, **enumerated from the diff
`BASE..M9b`, not from the test plan**): the version constant, both table sets, **the `< 9` upgrade line** and its
bound, the "newer than" branch now at 9; **each CHECK** of the DDL (id length and GLOB, request bounds, digest, challenge,
episode, expiry, the outcome set, approved_by, artifact, decision_ref's length and GLOB, consumed⇔outcome, and the
approved⇔(approved_by ∧ artifact ∧ decision_ref) CHECK and each of its three terms); **each DDL column** of the
upgrade copy (12, CS1's fresh-vs-upgraded comparison); `PutVerbRequest`'s deadline, **version check**, **each grammar
check** (id, both size bounds and the constant, digest hex, challenge hex, digest = SHA-256, episode, expiry, each of
the five "pending" fields), **the supersede**, **its `WHERE episode_id`**, its pending guard and its timestamp, **the
prune**, its bound, its pending guard and its constant, **the 16 cap** (`>=`, removed, the constant) and both count
filters; `VerbRequest`'s and `PendingVerbRequests`' deadline and **version check**, the no-rows mapping, both list
filters, the order, the limit clamp, three scan columns; `ConsumeVerbRequest`'s deadline, **version check**, each grammar
check, **each refusal** (unknown, superseded, consumed, clock behind and its boundary, expired and both boundaries,
digest by parameter, digest by bytes), **`useApprover`**, **`usePasskey`**, the outcome, the timestamp, the artifact
column and **`decision_ref` written in a second statement** (`CON-DECISION-SPLIT`); **`FinishConfirmRefusal` F1–F5**
(F2's two bounds, each removed, and the hour; F3 against the body hash; F5's protected paths and committer); **each route
registration** (removed; registered on the bare mux with no `approvalRoute`; **a POST added**; the method dropped);
each handler's id check, headers, content type and bytes; the open, **the 503 below v9** and its status, the 404 and its
status, the passkeys, the clock, and **each `ConfirmationInput` field**; `requests.json`'s older-store and error branches,
the empty list, the read itself, its limit and its clock, and `pendingConfirmation`'s digest, decode, `Validate`, scope
and each summary field; and both new JSON tags. For each row: an APFS clone of `work` at `M9b` (`cp -cR`;
byte-identical to the round-trip tree), one exact-string replacement whose anchor occurs exactly once (checked before the
run: 0 bad anchors), then **one** `go test -count=1 -v` of the row's package suite — `host/store` whole;
`host/broker` `^TestFinish`; `host/daemon` the confirm/approval/decision/route/requests subset
(`^(TestConfirm|TestRequestsJSON|TestRequestCardJSON|TestApprovalRoutes|TestDecisionRoutes|TestDecidePage|TestSurfaceApp|TestFrozen|TestSurfaceRequests|TestDaemonCredentialWriterSites)`,
the whole package takes ≈ 300 s) — from whose `--- FAIL`/`--- PASS` lines the **named test**'s result, the **`suite`**
(every failing test) and **`suite-nostatic`** (the failures other than the route-table and inventory tests
`TestConfirmRoutesRegistered`, `TestApprovalRoutesRegistered`, `TestDecisionRoutesRegistered`,
`TestFrozenV1RouteTableMatchesMux`, `TestDaemonCredentialWriterSites`) are read — the behavioural proof the directive asks
for. Rows ran with `GOFLAGS=-vet=off`, four workers, at load ≈ 20–30 (other sessions). Every run appends one line per
row to the ledger (time, run, id, file, expected, named result, named failures, suites, named regex, seconds); per-row
output in `logs/mut/<run>-<id>.log`. **The ids below are the ledger's ids**; the table shows each row's last line.
Classes: `K` killed by the named test **and** by the behavioural suite; `R` the route-table proof (the directive: "prove
each route's registration, and that no POST is registered, by the route-table test"); `E` equivalent (reason given);
`CTL` the unmutated tree, green. A run R0 (iteration 402, 18:04Z) was stopped after four rows when a test was added;
its lines are kept in `logs/mutant_ledger_R0_aborted.tsv` and are not in the table.

**R1 → R2 (what changed).** R1 (83 of 178 rows, the first shipped stages) found two rows not as expected and was
stopped to fix them: **`G-DIGEST-HEX`** survived everything — `PutVerbRequest` checked the digest's grammar *and* that
it equals the bytes' SHA-256, so the first check was subsumed (an equivalent branch): the production grammar now checks
`digest == sha256(request)` (which is 64 lowercase hex by construction) and then the challenge's grammar alone, and the
row is gone; **`G-MAX-CONST`** was killed only by CS1's constant equality, not by CS4 — CS4's exact-size control now
uses `approval.MaxVerbRequestBytes` (the bound the CLI's request obeys), so a store bound below it refuses the control.
R2 re-ran every row on the final stages; the table shows R2. R1's lines stay in the ledger.

**Summary:** 176 rows (3 `CTL`, 5 `E`, 166 `K`, 2 `R`); 276 ledger lines (runs R1, R2, R3, R4); 176 of 176 rows run; the last line of each row behaves as its class says for **176 of 176**. Table ids = ledger ids (`tools/mktable.py` fills §4 from the ledger's last line per id); ledger `$HOME/.ailang/state/world-iter402/plan-scratch/logs/mutant_ledger.tsv`

| id | file | from → to | named test → result | suites (`suite-nostatic`: route-table/static tests excluded) | class | as expected |
|---|---|---|---|---|---|---|
| `CTL-STORE` | `verb_requests.go` | control (no change) | `(TestCredentialsSchema9|TestPutVerb|TestCons` → **PASS** | suite=PASS; suite-nostatic=PASS | CTL | yes |
| `CTL-BROKER` | `finish_confirm.go` | control (no change) | `TestFinishConfirmRefusal` → **PASS** | suite=PASS; suite-nostatic=PASS | CTL | yes |
| `CTL-DAEMON` | `surface_confirm_read.go` | control (no change) | `(TestConfirm|TestRequestsJSONConfirmations)` → **PASS** | suite=PASS; suite-nostatic=PASS | CTL | yes |
| `CRD-VERSION` | `credentials.go` | `const credentialsSchemaVersion = 9` → `const credentialsSchemaVersion = 8` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestAppendEpisodeLaunchPreservesMintJSON,TestConcurrentApprovalsMintOnce,TestConsumeVerbRequest,…+17]; suite-nostatic=KILLED[TestAppendEpisodeLaunchPreservesMintJSON,TestConcurrentApprovalsMintOnce,TestConsumeVerbRequest,…+17] | K | yes |
| `CRD-UPGRADE-LINE` | `credentials.go` | `if version < 9 {⏎statements = append(statements, verbRequ…` → `` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema7SurfaceApprovals,TestCredentialsSchema8KeyDecisions,TestCredentialsSchema9VerbRequests,…+4]; suite-nostatic=KILLED[TestCredentialsSchema7SurfaceApprovals,TestCredentialsSchema8KeyDecisions,TestCredentialsSchema9VerbRequests,…+4] | K | yes |
| `CRD-UPGRADE-BOUND` | `credentials.go` | `if version < 9 {` → `if version < 8 {` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestConvertedStoreMethodsM5a,TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestConvertedStoreMethodsM5a,TestCredentialsSchema9VerbRequests] | K | yes |
| `CRD-TABLES9` | `credentials.go` | `surfaceApprovalsTable, verbRequestsTable, passkeysTable}` → `surfaceApprovalsTable, passkeysTable}` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestAppendEpisodeLaunchPreservesMintJSON,TestConcurrentApprovalsMintOnce,TestCredentialsMintAndRevokeBesideALiveWorldWriter,…+12]; suite-nostatic=KILLED[TestAppendEpisodeLaunchPreservesMintJSON,TestConcurrentApprovalsMintOnce,TestCredentialsMintAndRevokeBesideALiveWorldWriter,…+12] | K | yes |
| `CRD-TABLES9-CASE` | `credentials.go` | `case version >= 9:` → `case version >= 10:` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestAppendEpisodeLaunchPreservesMintJSON,TestConcurrentApprovalsMintOnce,TestCredentialsMintAndRevokeBesideALiveWorldWriter,…+12]; suite-nostatic=KILLED[TestAppendEpisodeLaunchPreservesMintJSON,TestConcurrentApprovalsMintOnce,TestCredentialsMintAndRevokeBesideALiveWorldWriter,…+12] | K | yes |
| `CRD-TABLES8` | `credentials.go` | `case version == 8:⏎return []string{approversTable, "episo…` → `` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema8KeyDecisions,TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema8KeyDecisions,TestCredentialsSchema9VerbRequests] | K | yes |
| `CRD-NEWER` | `credentials.go` | `case version > credentialsSchemaVersion:` → `case version > credentialsSchemaVersion+1:` | `TestCredentialsNewerThanNineIsRefused` → **KILLED** | suite=KILLED[TestCredentialsNewerThanNineIsRefused]; suite-nostatic=KILLED[TestCredentialsNewerThanNineIsRefused] | K | yes |
| `DDL-ID-LEN` | `credentials_schema.sql` | `CHECK (length(request_id) = 19 AND request_id GLOB 'vr-[0…` → `CHECK (request_id GLOB 'vr-[0-9a-f]*')` | `TestVerbRequestsChecks` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks] | K | yes |
| `DDL-ID-GLOB` | `credentials_schema.sql` | `…request_id) = 19 AND request_id GLOB 'vr-[0-9a-f]*')` → `…request_id) = 19)` | `TestVerbRequestsChecks` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks] | K | yes |
| `DDL-REQ-LEN` | `credentials_schema.sql` | `… BLOB NOT NULL CHECK (length(request) BETWEEN 1 AND 1310…` → `… BLOB NOT NULL,` | `TestVerbRequestsChecks` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks] | K | yes |
| `DDL-REQ-MAX` | `credentials_schema.sql` | `BETWEEN 1 AND 131072` → `BETWEEN 1 AND 131073` | `TestVerbRequestsChecks` → **KILLED** | suite=KILLED[TestCrashReceiptLawAtNamedStopPoints,TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks]; suite-nostatic=KILLED[TestCrashReceiptLawAtNamedStopPoints,TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks] | K | yes |
| `DDL-REQ-MIN` | `credentials_schema.sql` | `CHECK (length(request) BETWEEN 1 AND 131072)` → `CHECK (length(request) BETWEEN 0 AND 131072)` | `TestVerbRequestsChecks` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks] | K | yes |
| `DDL-DIGEST` | `credentials_schema.sql` | `… TEXT NOT NULL CHECK (length(digest) = 64),⏎ challenge` → `… TEXT NOT NULL,⏎ challenge` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests] | K | yes |
| `DDL-CHALLENGE` | `credentials_schema.sql` | `… TEXT NOT NULL CHECK (length(challenge) = 64),⏎ episode_…` → `… TEXT NOT NULL,⏎ episode_id TEXT NOT NULL CHECK (episode…` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests] | K | yes |
| `DDL-EPISODE` | `credentials_schema.sql` | `… TEXT NOT NULL CHECK (episode_id <> ''),⏎ created_at INT…` → `… TEXT NOT NULL,⏎ created_at INTEGER NOT NULL,⏎ expires_a…` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests] | K | yes |
| `DDL-EXPIRES` | `credentials_schema.sql` | `…INTEGER NOT NULL CHECK (expires_at > created_at),⏎ consu…` → `…INTEGER NOT NULL,⏎ consumed_at INTEGER,⏎ outcome` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests] | K | yes |
| `DDL-OUTCOME-IN` | `credentials_schema.sql` | `outcome IN ('approved', 'superseded')` → `outcome IN ('approved', 'superseded', 'denied')` | `TestVerbRequestsChecks` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks] | K | yes |
| `DDL-BY-EMPTY` | `credentials_schema.sql` | `approved_by TEXT CHECK (approved_by IS NULL OR approved_b…` → `approved_by TEXT,⏎ artifact BLOB CHECK (artifact IS NULL …` | `TestVerbRequestsChecks` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks] | K | yes |
| `DDL-ARTIFACT` | `credentials_schema.sql` | `artifact BLOB CHECK (artifact IS NULL OR length(artifact)…` → `artifact BLOB,⏎ decision_ref TEXT CHECK (decision_ref IS …` | `TestVerbRequestsChecks` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks] | K | yes |
| `DDL-DECISION-LEN` | `credentials_schema.sql` | `…ref IS NULL OR (length(decision_ref) = 71 AND decision_r…` → `…ref IS NULL OR (decision_ref GLOB 'sha256:*')),⏎ CHECK (…` | `TestVerbRequestsChecks` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks] | K | yes |
| `DDL-DECISION-GLOB` | `credentials_schema.sql` | `…cision_ref) = 71 AND decision_ref GLOB 'sha256:*')),⏎ CH…` → `…cision_ref) = 71)),⏎ CHECK ((consumed_at` | `TestVerbRequestsChecks` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests,TestQuarantineRejectsOperationsUntilWorkerSettles,TestVerbRequestsChecks]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests,TestQuarantineRejectsOperationsUntilWorkerSettles,TestVerbRequestsChecks] | K | yes |
| `DDL-CONSUMED-IFF` | `credentials_schema.sql` | ` CHECK ((consumed_at IS NULL) = (outcome IS NULL)),⏎ CHEC…` → ` CHECK ((outcome IS 'approved')` | `TestVerbRequestsChecks` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks] | K | yes |
| `DDL-APPROVED-IFF` | `credentials_schema.sql` | ` CHECK ((outcome IS 'approved') = (approved_by IS NOT NUL…` → ` CHECK (1)⏎);⏎` | `TestVerbRequestsChecks` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks] | K | yes |
| `DDL-APPROVED-BY` | `credentials_schema.sql` | `(approved_by IS NOT NULL AND artifact IS NOT NULL AND dec…` → `(artifact IS NOT NULL AND decision_ref IS NOT NULL))⏎);⏎` | `TestVerbRequestsChecks` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks] | K | yes |
| `DDL-APPROVED-ARTIFACT` | `credentials_schema.sql` | `…IS NOT NULL AND artifact IS NOT NULL AND decision_ref IS…` → `…IS NOT NULL AND decision_ref IS NOT NULL))⏎);⏎` | `TestVerbRequestsChecks` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks] | K | yes |
| `DDL-APPROVED-DECISION` | `credentials_schema.sql` | `…fact IS NOT NULL AND decision_ref IS NOT NULL))⏎);⏎` → `…fact IS NOT NULL))⏎);⏎` | `TestVerbRequestsChecks` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests,TestVerbRequestsChecks] | K | yes |
| `COL-REQUEST` | `verb_requests.go` | ` request BLOB NOT NULL CHECK` → ` request TEXT NOT NULL CHECK` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests] | K | yes |
| `COL-DIGEST` | `verb_requests.go` | ` digest TEXT NOT NULL CHECK` → ` digest TEXT CHECK` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests] | K | yes |
| `COL-CHALLENGE` | `verb_requests.go` | ` challenge TEXT NOT NULL CHECK` → ` challenge TEXT CHECK` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests] | K | yes |
| `COL-EPISODE` | `verb_requests.go` | ` episode_id TEXT NOT NULL CHECK (episode_id <> ''),` → ` episode_id TEXT CHECK (episode_id <> ''),` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests] | K | yes |
| `COL-CREATED` | `verb_requests.go` | ` created_at INTEGER NOT NULL,` → ` created_at INTEGER,` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests] | K | yes |
| `COL-EXPIRES` | `verb_requests.go` | ` expires_at INTEGER NOT NULL CHECK` → ` expires_at INTEGER CHECK` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests] | K | yes |
| `COL-CONSUMED` | `verb_requests.go` | ` consumed_at INTEGER,⏎` → ` consumed_at TEXT,⏎` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests] | K | yes |
| `COL-OUTCOME` | `verb_requests.go` | ` outcome TEXT CHECK` → ` outcome BLOB CHECK` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests] | K | yes |
| `COL-APPROVEDBY` | `verb_requests.go` | ` approved_by TEXT CHECK` → ` approved_by BLOB CHECK` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests] | K | yes |
| `COL-ARTIFACT` | `verb_requests.go` | ` artifact BLOB CHECK` → ` artifact TEXT CHECK` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests] | K | yes |
| `COL-DECISION` | `verb_requests.go` | ` decision_ref TEXT CHECK` → ` decision_ref BLOB CHECK` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests] | K | yes |
| `COL-ID` | `verb_requests.go` | ` request_id TEXT PRIMARY KEY CHECK` → ` request_id TEXT CHECK` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests] | K | yes |
| `PUT-DEADLINE` | `verb_requests.go` | `if err := requireDeadline(ctx); err != nil { // S1` → `if err := requireDeadline(ctx); false && err != nil { // …` | `TestCredentialsMethodsRequireADeadline` → **KILLED** | suite=KILLED[TestCredentialsMethodsRequireADeadline]; suite-nostatic=KILLED[TestCredentialsMethodsRequireADeadline] | K | yes |
| `PUT-VERSION` | `verb_requests.go` | `if c.version < verbRequestsSchemaVersion { // S1` → `if false { // S1` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests] | K | yes |
| `G-ID` | `verb_requests.go` | `case !VerbRequestIDPattern.MatchString(row.RequestID):` → `case false:` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |
| `G-REQ-MIN` | `verb_requests.go` | `case len(row.Request) < 1 \|\| len(row.Request) > MaxVerb…` → `case len(row.Request) > MaxVerbRequestBytes:` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |
| `G-REQ-MAX` | `verb_requests.go` | `…row.Request) < 1 \|\| len(row.Request) > MaxVerbRequestB…` → `…row.Request) < 1:` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |
| `G-MAX-CONST` | `verb_requests.go` | `const MaxVerbRequestBytes = 128 << 10` → `const MaxVerbRequestBytes = 128<<10 - 1` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests,TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests,TestPutVerbRequestGrammarCapPrune] | K | yes |
| `G-CHALLENGE-HEX` | `verb_requests.go` | `case !isLowerHex64(row.Challenge):` → `case false:` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |
| `G-DIGEST-SHA` | `verb_requests.go` | `case row.Digest != sha256Hex(row.Request):` → `case false:` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestCloseKeepsWriterLockUntilDurableWorkerSettles,TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestCloseKeepsWriterLockUntilDurableWorkerSettles,TestPutVerbRequestGrammarCapPrune] | K | yes |
| `G-EPISODE` | `verb_requests.go` | `case row.EpisodeID == "":` → `case false:` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestCloseKeepsWriterLockUntilDurableWorkerSettles,TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestCloseKeepsWriterLockUntilDurableWorkerSettles,TestPutVerbRequestGrammarCapPrune] | K | yes |
| `G-EXPIRES` | `verb_requests.go` | `case row.ExpiresAt <= row.CreatedAt:` → `case false:` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |
| `G-PENDING-OUTCOME` | `verb_requests.go` | `case row.ConsumedAt != 0 \|\| row.Outcome != "" \|\| ` → `case row.ConsumedAt != 0 \|\| ` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |
| `G-PENDING-CONSUMED` | `verb_requests.go` | `case row.ConsumedAt != 0 \|\| row.Outcome != "" \|\| ` → `case row.Outcome != "" \|\| ` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |
| `G-PENDING-BY` | `verb_requests.go` | `row.ApprovedBy != "" \|\| row.Artifact != nil \|\| row.De…` → `row.Artifact != nil \|\| row.DecisionRef != "":` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |
| `G-PENDING-ARTIFACT` | `verb_requests.go` | `…By != "" \|\| row.Artifact != nil \|\| row.DecisionRef !…` → `…By != "" \|\| row.DecisionRef != "":` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |
| `G-PENDING-DECISION` | `verb_requests.go` | `….Artifact != nil \|\| row.DecisionRef != "":` → `….Artifact != nil:` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |
| `PUT-GRAMMAR-CALL` | `verb_requests.go` | `if err := checkVerbRequestRow(row); err != nil { // S2` → `if err := checkVerbRequestRow(row); false && err != nil {…` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |
| `PUT-SUPERSEDE` | `verb_requests.go` | `…sode_id = ? AND consumed_at IS NULL`, now, row.EpisodeID)` → `…sode_id = ? AND 0`, now, row.EpisodeID)` | `TestPutVerbRequestSupersedes` → **KILLED** | suite=KILLED[TestConsumeVerbRequest,TestPutVerbRequestGrammarCapPrune,TestPutVerbRequestSupersedes]; suite-nostatic=KILLED[TestConsumeVerbRequest,TestPutVerbRequestGrammarCapPrune,TestPutVerbRequestSupersedes] | K | yes |
| `PUT-SUPERSEDE-WHERE` | `verb_requests.go` | `…perseded' WHERE episode_id = ? AND consumed_at IS NULL`,…` → `…perseded' WHERE (episode_id = ? OR 1) AND consumed_at IS…` | `TestPutVerbRequestSupersedes` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune,TestPutVerbRequestSupersedes]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune,TestPutVerbRequestSupersedes] | K | yes |
| `PUT-SUPERSEDE-PENDING` | `verb_requests.go` | `…E episode_id = ? AND consumed_at IS NULL`, now, row.Epis…` → `…E episode_id = ?`, now, row.EpisodeID)` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `PUT-SUPERSEDE-AT` | `verb_requests.go` | `SET consumed_at = ?, outcome = 'superseded' WHERE` → `SET consumed_at = ? + 1, outcome = 'superseded' WHERE` | `TestPutVerbRequestSupersedes` → **KILLED** | suite=KILLED[TestPutVerbRequestSupersedes]; suite-nostatic=KILLED[TestPutVerbRequestSupersedes] | K | yes |
| `PUT-PRUNE` | `verb_requests.go` | `…_requests WHERE consumed_at IS NULL AND expires_at <= ?`` → `…_requests WHERE 0 AND expires_at <= ?`` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |
| `PUT-PRUNE-BOUND` | `verb_requests.go` | `…AND expires_at <= ?`` → `…AND expires_at < ?`` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |
| `PUT-PRUNE-PENDING` | `verb_requests.go` | `…_requests WHERE consumed_at IS NULL AND expires_at <= ?`` → `…_requests WHERE expires_at <= ?`` | `TestPutVerbRequestSupersedes` → **KILLED** | suite=KILLED[TestPutVerbRequestSupersedes]; suite-nostatic=KILLED[TestPutVerbRequestSupersedes] | K | yes |
| `PUT-PRUNE-AFTER` | `verb_requests.go` | `const verbRequestPruneAfter = 86400` → `const verbRequestPruneAfter = 86401` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |
| `PUT-CAP` | `verb_requests.go` | `if live >= MaxPendingVerbRequests { // S5` → `if live > MaxPendingVerbRequests { // S5` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |
| `PUT-CAP-OFF` | `verb_requests.go` | `if live >= MaxPendingVerbRequests { // S5` → `if false { // S5` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |
| `PUT-CAP-CONST` | `verb_requests.go` | `const MaxPendingVerbRequests = 16` → `const MaxPendingVerbRequests = 17` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |
| `PUT-CAP-COUNT-PENDING` | `verb_requests.go` | `…_requests WHERE consumed_at IS NULL AND expires_at > ?`` → `…_requests WHERE expires_at > ?`` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |
| `PUT-CAP-COUNT-LIVE` | `verb_requests.go` | `…_at IS NULL AND expires_at > ?`` → `…_at IS NULL AND ? > 0`` | `TestPutVerbRequestGrammarCapPrune` → **KILLED** | suite=KILLED[TestPutVerbRequestGrammarCapPrune]; suite-nostatic=KILLED[TestPutVerbRequestGrammarCapPrune] | K | yes |
| `READ-DEADLINE` | `verb_requests.go` | `if err := requireDeadline(ctx); err != nil {⏎return VerbR…` → `` | `TestCredentialsMethodsRequireADeadline` → **KILLED** | suite=KILLED[TestCredentialsMethodsRequireADeadline]; suite-nostatic=KILLED[TestCredentialsMethodsRequireADeadline] | K | yes |
| `READ-VERSION` | `verb_requests.go` | `if c.version < verbRequestsSchemaVersion {⏎return VerbReq…` → `` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests] | K | yes |
| `READ-NOROWS` | `verb_requests.go` | `if errors.Is(err, sql.ErrNoRows) {⏎return VerbRequestRow{…` → `` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `LIST-DEADLINE` | `verb_requests.go` | `if err := requireDeadline(ctx); err != nil {⏎return nil, …` → `if c.version < verbRequestsSchemaVersion {` | `TestCredentialsMethodsRequireADeadline` → **KILLED** | suite=KILLED[TestCredentialsMethodsRequireADeadline]; suite-nostatic=KILLED[TestCredentialsMethodsRequireADeadline] | K | yes |
| `LIST-VERSION` | `verb_requests.go` | `if c.version < verbRequestsSchemaVersion {⏎return nil, c.…` → `` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests] | K | yes |
| `LIST-PENDING` | `verb_requests.go` | `…_requests WHERE consumed_at IS NULL AND expires_at > ? O…` → `…_requests WHERE expires_at > ? ORDER BY` | `TestPutVerbRequestSupersedes` → **KILLED** | suite=KILLED[TestConsumeVerbRequest,TestPutVerbRequestSupersedes]; suite-nostatic=KILLED[TestConsumeVerbRequest,TestPutVerbRequestSupersedes] | K | yes |
| `LIST-LIVE` | `verb_requests.go` | `…_at IS NULL AND expires_at > ? ORDER BY` → `…_at IS NULL AND ? > 0 ORDER BY` | `TestRequestsJSONConfirmations` → **KILLED** | suite=KILLED[TestRequestsJSONConfirmations]; suite-nostatic=KILLED[TestRequestsJSONConfirmations] | K | yes |
| `LIST-ORDER` | `verb_requests.go` | `ORDER BY created_at DESC, request_id DESC LIMIT ?` → `ORDER BY created_at, request_id LIMIT ?` | `TestPutVerbRequestSupersedes` → **KILLED** | suite=KILLED[TestPutVerbRequestSupersedes]; suite-nostatic=KILLED[TestPutVerbRequestSupersedes] | K | yes |
| `LIST-LIMIT` | `verb_requests.go` | `if limit <= 0 \|\| limit > MaxPendingVerbRequests {` → `if limit <= 0 {` | `TestPutVerbRequestGrammarCapPrune` → **PASS** | suite=PASS; suite-nostatic=PASS | E — equivalent through the API: PutVerbRequest refuses a 17th live row, so no store reachable by the store's methods at a monotonic clock holds more than 16 live rows; the clamp is the mint list's defence in depth | yes |
| `LIST-LIMIT-ZERO` | `verb_requests.go` | `if limit <= 0 \|\| limit > MaxPendingVerbRequests {` → `if limit > MaxPendingVerbRequests {` | `TestPutVerbRequestSupersedes` → **KILLED** | suite=KILLED[TestPutVerbRequestSupersedes]; suite-nostatic=KILLED[TestPutVerbRequestSupersedes] | K | yes |
| `SCAN-OUTCOME` | `verb_requests.go` | `coalesce(outcome, ''), coalesce(approved_by, '')` → `'', coalesce(approved_by, '')` | `TestPutVerbRequestSupersedes` → **KILLED** | suite=KILLED[TestConsumeVerbRequest,TestPutVerbRequestSupersedes]; suite-nostatic=KILLED[TestConsumeVerbRequest,TestPutVerbRequestSupersedes] | K | yes |
| `SCAN-CONSUMED` | `verb_requests.go` | `coalesce(consumed_at, 0), coalesce(outcome` → `0, coalesce(outcome` | `TestPutVerbRequestSupersedes` → **KILLED** | suite=KILLED[TestConsumeVerbRequest,TestPutVerbRequestSupersedes]; suite-nostatic=KILLED[TestConsumeVerbRequest,TestPutVerbRequestSupersedes] | K | yes |
| `SCAN-DECISION` | `verb_requests.go` | `artifact, coalesce(decision_ref, '')`` → `artifact, ''`` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-DEADLINE` | `verb_requests.go` | `if err := requireDeadline(ctx); err != nil {⏎return err⏎}…` → `if c.version < verbRequestsSchemaVersion {⏎return c.older…` | `TestCredentialsMethodsRequireADeadline` → **KILLED** | suite=KILLED[TestCredentialsMethodsRequireADeadline]; suite-nostatic=KILLED[TestCredentialsMethodsRequireADeadline] | K | yes |
| `CON-VERSION` | `verb_requests.go` | `if c.version < verbRequestsSchemaVersion {⏎return c.older…` → `if err := checkVerbConsume` | `TestCredentialsSchema9VerbRequests` → **KILLED** | suite=KILLED[TestCredentialsSchema9VerbRequests]; suite-nostatic=KILLED[TestCredentialsSchema9VerbRequests] | K | yes |
| `CON-G-ID` | `verb_requests.go` | `case !VerbRequestIDPattern.MatchString(id) \|\| !isLowerH…` → `case !isLowerHex64(digest):` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-G-DIGEST` | `verb_requests.go` | `….MatchString(id) \|\| !isLowerHex64(digest):` → `….MatchString(id):` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-G-APPROVER` | `verb_requests.go` | `case approvedBy != use.Format+":"+use.Name:` → `case false:` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-G-PASSKEY` | `verb_requests.go` | `case (use.Format == "webauthn") != (use.Passkey != nil):` → `case false:` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-G-ARTIFACT-MIN` | `verb_requests.go` | `case len(artifact) < 1 \|\| len(artifact) > MaxSurfaceApp…` → `case len(artifact) > MaxSurfaceApprovalArtifactBytes:` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-G-ARTIFACT-MAX` | `verb_requests.go` | `case len(artifact) < 1 \|\| len(artifact) > MaxSurfaceApp…` → `case len(artifact) < 1:` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-G-DECISION` | `verb_requests.go` | `case !isDecisionRef(decisionRef):` → `case false:` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-G-CALL` | `verb_requests.go` | `…, decisionRef); err != nil {` → `…, decisionRef); false && err != nil {` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-UNKNOWN` | `verb_requests.go` | `if errors.Is(err, sql.ErrNoRows) {⏎return fmt.Errorf("%w:…` → `if false {⏎return fmt.Errorf("%w: %s", ErrVerbRequestUnkn…` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-SUPERSEDED` | `verb_requests.go` | `case row.Outcome == VerbOutcomeSuperseded:` → `case false:` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest,TestConvertedStoreMethodsM5a]; suite-nostatic=KILLED[TestConsumeVerbRequest,TestConvertedStoreMethodsM5a] | K | yes |
| `CON-CONSUMED` | `verb_requests.go` | `case row.Outcome != "":` → `case false:` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-CLOCK` | `verb_requests.go` | `case now < row.CreatedAt:` → `case false:` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-CLOCK-BOUND` | `verb_requests.go` | `case now < row.CreatedAt:` → `case now <= row.CreatedAt:` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-EXPIRED` | `verb_requests.go` | `case now >= row.ExpiresAt:` → `case false:` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-EXPIRED-BOUND` | `verb_requests.go` | `case now >= row.ExpiresAt:` → `case now > row.ExpiresAt:` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-EXPIRED-EARLY` | `verb_requests.go` | `case now >= row.ExpiresAt:` → `case now >= row.ExpiresAt-1:` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-DIGEST-PARAM` | `verb_requests.go` | `case digest != row.Digest \|\| sha256Hex(row.Request) != …` → `case sha256Hex(row.Request) != row.Digest:` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-DIGEST-BYTES` | `verb_requests.go` | `…st != row.Digest \|\| sha256Hex(row.Request) != row.Dige…` → `…st != row.Digest:` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-USEAPPROVER` | `verb_requests.go` | `if err := useApprover(ctx, tx, use, now); err != nil {⏎re…` → `if use.Passkey` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-USEPASSKEY` | `verb_requests.go` | `if use.Passkey != nil {⏎if err := usePasskey` → `if false {⏎if err := usePasskey` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-OUTCOME` | `verb_requests.go` | `SET consumed_at = ?, outcome = 'approved', approved_by` → `SET consumed_at = ?, outcome = 'superseded', approved_by` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-AT` | `verb_requests.go` | `SET consumed_at = ?, outcome = 'approved', approved_by` → `SET consumed_at = ? + 1, outcome = 'approved', approved_by` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-DECISION-SPLIT` | `verb_requests.go` | `… ?, artifact = ?, decision_ref = ? WHERE request_id = ? …` → `… ?, artifact = ? WHERE request_id = ? AND consumed_at IS…` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-ARTIFACT-COL` | `verb_requests.go` | `approved_by = ?, artifact = ?, decision_ref = ? WHERE` → `approved_by = ?, artifact = ? \|\| 'x', decision_ref = ? …` | `TestConsumeVerbRequest` → **KILLED** | suite=KILLED[TestConsumeVerbRequest]; suite-nostatic=KILLED[TestConsumeVerbRequest] | K | yes |
| `CON-UPDATE-GUARD` | `verb_requests.go` | `WHERE request_id = ? AND consumed_at IS NULL`,⏎now, appro…` → `WHERE request_id = ?`,⏎now, approvedBy, artifact, decisio…` | `TestConsumeVerbRequest` → **PASS** | suite=PASS; suite-nostatic=PASS | E — equivalent: the row was re-read in this transaction and refused when consumed; SQLite serialises writers, so `AND consumed_at IS NULL` is defence in depth (as the mint consume's guarded UPDATE) | yes |
| `CON-ROWS` | `verb_requests.go` | `if n, err := res.RowsAffected(); err != nil \|\| n != 1 {` → `if n, err := res.RowsAffected(); false && (err != nil \|\…` | `TestConsumeVerbRequest` → **PASS** | suite=PASS; suite-nostatic=PASS | E — equivalent: the UPDATE runs on a row this transaction just read pending, so it changes exactly one row; the RowsAffected check is defence in depth | yes |
| `F1` | `finish_confirm.go` | `if err != nil { // F1` → `if err != nil && false { // F1` | `TestFinishConfirmRefusal` → **KILLED** | suite=KILLED[TestFinishConfirmRefusal]; suite-nostatic=KILLED[TestFinishConfirmRefusal] | K | yes |
| `F2-OFF` | `finish_confirm.go` | `if !(f.RequestExpires < scope.ExpiresAt && scope.ExpiresA…` → `if false { // F2` | `TestFinishConfirmRefusal` → **KILLED** | suite=KILLED[TestFinishConfirmRefusal]; suite-nostatic=KILLED[TestFinishConfirmRefusal] | K | yes |
| `F2-LOWER` | `finish_confirm.go` | `f.RequestExpires < scope.ExpiresAt &&` → `f.RequestExpires <= scope.ExpiresAt &&` | `TestFinishConfirmRefusal` → **KILLED** | suite=KILLED[TestFinishConfirmRefusal]; suite-nostatic=KILLED[TestFinishConfirmRefusal] | K | yes |
| `F2-LOWER-OFF` | `finish_confirm.go` | `if !(f.RequestExpires < scope.ExpiresAt && scope.ExpiresA…` → `if !(scope.ExpiresAt <= f.Now+3600)` | `TestFinishConfirmRefusal` → **KILLED** | suite=KILLED[TestFinishConfirmRefusal]; suite-nostatic=KILLED[TestFinishConfirmRefusal] | K | yes |
| `F2-UPPER` | `finish_confirm.go` | `scope.ExpiresAt <= f.Now+3600` → `scope.ExpiresAt < f.Now+3600` | `TestFinishConfirmRefusal` → **KILLED** | suite=KILLED[TestFinishConfirmRefusal]; suite-nostatic=KILLED[TestFinishConfirmRefusal] | K | yes |
| `F2-UPPER-OFF` | `finish_confirm.go` | `… scope.ExpiresAt && scope.ExpiresAt <= f.Now+3600)` → `… scope.ExpiresAt)` | `TestFinishConfirmRefusal` → **KILLED** | suite=KILLED[TestFinishConfirmRefusal]; suite-nostatic=KILLED[TestFinishConfirmRefusal] | K | yes |
| `F2-HOUR` | `finish_confirm.go` | `f.Now+3600` → `f.Now+3601` | `TestFinishConfirmRefusal` → **KILLED** | suite=KILLED[TestFinishConfirmRefusal]; suite-nostatic=KILLED[TestFinishConfirmRefusal] | K | yes |
| `F3` | `finish_confirm.go` | `hex.EncodeToString(title[:]) != scope.TitleSHA256 { // F3` → `hex.EncodeToString(title[:]) != scope.TitleSHA256 && fals…` | `TestFinishConfirmRefusal` → **KILLED** | suite=KILLED[TestFinishConfirmRefusal]; suite-nostatic=KILLED[TestFinishConfirmRefusal] | K | yes |
| `F3-BODYHASH` | `finish_confirm.go` | `hex.EncodeToString(title[:]) != scope.TitleSHA256 { // F3` → `hex.EncodeToString(title[:]) != scope.BodySHA256 { // F3` | `TestFinishConfirmRefusal` → **KILLED** | suite=KILLED[TestFinishConfirmRefusal]; suite-nostatic=KILLED[TestFinishConfirmRefusal] | K | yes |
| `F4` | `finish_confirm.go` | `hex.EncodeToString(body[:]) != scope.BodySHA256 { // F4` → `hex.EncodeToString(body[:]) != scope.BodySHA256 && false …` | `TestFinishConfirmRefusal` → **KILLED** | suite=KILLED[TestFinishConfirmRefusal]; suite-nostatic=KILLED[TestFinishConfirmRefusal] | K | yes |
| `F5` | `finish_confirm.go` | `if f.Prompt != FinishPromptText(scope, f.Title, f.Body, f…` → `if false { // F5` | `TestFinishConfirmRefusal` → **KILLED** | suite=KILLED[TestFinishConfirmRefusal]; suite-nostatic=KILLED[TestFinishConfirmRefusal] | K | yes |
| `F5-PROTECTED` | `finish_confirm.go` | `….Title, f.Body, f.Protected, f.Committer)` → `….Title, f.Body, nil, f.Committer)` | `TestFinishConfirmRefusal` → **KILLED** | suite=KILLED[TestFinishConfirmRefusal]; suite-nostatic=KILLED[TestFinishConfirmRefusal] | K | yes |
| `F5-COMMITTER` | `finish_confirm.go` | `…y, f.Protected, f.Committer)` → `…y, f.Protected, "")` | `TestFinishConfirmRefusal` → **KILLED** | suite=KILLED[TestFinishConfirmRefusal]; suite-nostatic=KILLED[TestFinishConfirmRefusal] | K | yes |
| `RTE-PAGE` | `daemon.go` | `mux.HandleFunc("GET /surface/confirm/{id}", d.approvalRou…` → `` | `TestConfirmPageAssetsServed` → **KILLED** | suite=KILLED[TestConfirmPageAssetsServed,TestConfirmRoutesRegistered]; suite-nostatic=KILLED[TestConfirmPageAssetsServed] | K | yes |
| `RTE-SCRIPT` | `daemon.go` | `mux.HandleFunc("GET /surface/confirm.js", d.approvalRoute…` → `` | `TestConfirmPageAssetsServed` → **KILLED** | suite=KILLED[TestConfirmPageAssetsServed,TestConfirmRoutesRegistered]; suite-nostatic=KILLED[TestConfirmPageAssetsServed] | K | yes |
| `RTE-DATA` | `daemon.go` | `mux.HandleFunc("GET /surface/confirmations/{id}", d.appro…` → `` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmRoutesRegistered,TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `RTE-PAGE-BARE` | `daemon.go` | `…/confirm/{id}", d.approvalRoute(d.handleSurfaceConfirmPa…` → `…/confirm/{id}", func(w http.ResponseWriter, r *http.Requ…` | `TestConfirmPageAssetsServed` → **KILLED** | suite=KILLED[TestConfirmPageAssetsServed,TestConfirmRoutesRegistered]; suite-nostatic=KILLED[TestConfirmPageAssetsServed] | K | yes |
| `RTE-SCRIPT-BARE` | `daemon.go` | `…ce/confirm.js", d.approvalRoute(d.handleSurfaceConfirmSc…` → `…ce/confirm.js", func(w http.ResponseWriter, r *http.Requ…` | `TestConfirmPageAssetsServed` → **KILLED** | suite=KILLED[TestConfirmPageAssetsServed,TestConfirmRoutesRegistered]; suite-nostatic=KILLED[TestConfirmPageAssetsServed] | K | yes |
| `RTE-DATA-BARE` | `daemon.go` | `…rmations/{id}", d.approvalRoute(d.handleSurfaceConfirmat…` → `…rmations/{id}", func(w http.ResponseWriter, r *http.Requ…` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmRoutesRegistered,TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `RTE-POST` | `daemon.go` | `…eConfirmation))⏎` → `…eConfirmation))⏎mux.HandleFunc("POST /surface/confirm/{i…` | `TestConfirmRoutesRegistered` → **KILLED** | suite=KILLED[TestConfirmRoutesRegistered]; suite-nostatic=PASS | R — killed by the route-table test (`TestConfirmRoutesRegistered`, its static and its 405 halves), the proof the directive names | yes |
| `RTE-PAGE-ANY` | `daemon.go` | `"GET /surface/confirm/{id}"` → `"/surface/confirm/{id}"` | `TestConfirmRoutesRegistered` → **KILLED** | suite=KILLED[TestConfirmRoutesRegistered]; suite-nostatic=PASS | R — killed by the route-table test (`TestConfirmRoutesRegistered`, its static and its 405 halves), the proof the directive names | yes |
| `PAGE-ID` | `surface_confirm_read.go` | `if !store.VerbRequestIDPattern.MatchString(r.PathValue("i…` → `if false {` | `TestConfirmPageAssetsServed` → **KILLED** | suite=KILLED[TestConfirmPageAssetsServed]; suite-nostatic=KILLED[TestConfirmPageAssetsServed] | K | yes |
| `PAGE-HEADERS` | `surface_confirm_read.go` | `setSurfaceHeaders(w)⏎_, _ = w.Write(surfacepage.ConfirmPa…` → `_, _ = w.Write(surfacepage.ConfirmPage)` | `TestConfirmPageAssetsServed` → **KILLED** | suite=KILLED[TestConfirmPageAssetsServed]; suite-nostatic=KILLED[TestConfirmPageAssetsServed] | K | yes |
| `PAGE-BYTES` | `surface_confirm_read.go` | `_, _ = w.Write(surfacepage.ConfirmPage)` → `_, _ = w.Write(surfacepage.DecidePage)` | `TestConfirmPageAssetsServed` → **KILLED** | suite=KILLED[TestConfirmPageAssetsServed]; suite-nostatic=KILLED[TestConfirmPageAssetsServed] | K | yes |
| `SCRIPT-HEADERS` | `surface_confirm_read.go` | `setSurfaceHeaders(w)⏎w.Header().Set("Content-Type", "text…` → `w.Header().Set("Content-Type", "text/javascript; charset=…` | `TestConfirmPageAssetsServed` → **KILLED** | suite=KILLED[TestConfirmPageAssetsServed]; suite-nostatic=KILLED[TestConfirmPageAssetsServed] | K | yes |
| `SCRIPT-TYPE` | `surface_confirm_read.go` | `w.Header().Set("Content-Type", "text/javascript; charset=…` → `_, _ = w.Write(surfacepage.ConfirmScript)` | `TestConfirmPageAssetsServed` → **KILLED** | suite=KILLED[TestConfirmPageAssetsServed]; suite-nostatic=KILLED[TestConfirmPageAssetsServed] | K | yes |
| `SCRIPT-BYTES` | `surface_confirm_read.go` | `_, _ = w.Write(surfacepage.ConfirmScript)` → `_, _ = w.Write(surfacepage.DecideScript)` | `TestConfirmPageAssetsServed` → **KILLED** | suite=KILLED[TestConfirmPageAssetsServed]; suite-nostatic=KILLED[TestConfirmPageAssetsServed] | K | yes |
| `DATA-ID` | `surface_confirm_read.go` | `if !store.VerbRequestIDPattern.MatchString(id) {` → `if false {` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `DATA-OPEN` | `surface_confirm_read.go` | `ro, err := store.OpenCredentialsReadOnly(d.credentialsPat…` → `ro, _ := store.OpenCredentialsReadOnly(d.credentialsPath)…` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `DATA-OLDER` | `surface_confirm_read.go` | `case isOlderStore(err):⏎` → `case false:⏎` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `DATA-OLDER-STATUS` | `surface_confirm_read.go` | `…l(w, http.StatusServiceUnavailable, "Restart",` → `…l(w, http.StatusInternalServerError, "Restart",` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `DATA-ERR` | `surface_confirm_read.go` | `case err != nil:⏎d.approvalInternal(w, r, err)⏎return⏎cas…` → `case !ok:` | `TestConfirmationCardServed` → **PASS** | suite=PASS; suite-nostatic=PASS | E — equivalent through the API: after a successful read-only open, VerbRequest fails only below schema 9 (the case before) or on an I/O fault no test can inject without a fault seam; the sanitized 500 is A3's card pattern | yes |
| `DATA-NOTFOUND` | `surface_confirm_read.go` | `case !ok:⏎` → `case !ok && false:⏎` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `DATA-NOTFOUND-STATUS` | `surface_confirm_read.go` | `…l(w, http.StatusNotFound, "NotFound", "no verb request "…` → `…l(w, http.StatusGone, "NotFound", "no verb request "+id)` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `DATA-PASSKEYS-ERR` | `surface_confirm_read.go` | `active, err := ro.ActivePasskeys(ctx)⏎if err != nil {` → `active, err := ro.ActivePasskeys(ctx)⏎if false && err != …` | `TestConfirmationCardServed` → **PASS** | suite=PASS; suite-nostatic=PASS | E — equivalent through the API: ActivePasskeys on a schema-9 handle fails only on an I/O fault (schema 6 passkeys are always present at 9) | yes |
| `DATA-PASSKEYS` | `surface_confirm_read.go` | `passkeys = append(passkeys, surfacepage.CardPasskey{` → `_ = append(passkeys, surfacepage.CardPasskey{` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `DATA-PASSKEY-ORIGIN` | `surface_confirm_read.go` | `surfacepage.CardPasskey{Origin: e.Passkey.Origin, RPID:` → `surfacepage.CardPasskey{Origin: e.Passkey.RPID, RPID:` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `DATA-NOW` | `surface_confirm_read.go` | `confirmationInput(row, d.approveNow())` → `confirmationInput(row, 0)` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `IN-ID` | `surface_confirm_read.go` | `ConfirmationInput{ID: row.RequestID,` → `ConfirmationInput{ID: row.EpisodeID,` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `IN-REQUEST` | `surface_confirm_read.go` | `Request: row.Request, Digest: row.Digest,` → `Request: row.Request[1:], Digest: row.Digest,` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `IN-DIGEST` | `surface_confirm_read.go` | `Digest: row.Digest, Challenge: row.Challenge,` → `Digest: row.Challenge, Challenge: row.Challenge,` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `IN-CHALLENGE` | `surface_confirm_read.go` | `Digest: row.Digest, Challenge: row.Challenge,` → `Digest: row.Digest, Challenge: row.Digest,` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `IN-CREATED` | `surface_confirm_read.go` | `Created: row.CreatedAt,` → `Created: row.ExpiresAt,` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `IN-EXPIRES` | `surface_confirm_read.go` | `Expires: row.ExpiresAt, Now: now,` → `Expires: row.CreatedAt, Now: now,` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `IN-OUTCOME` | `surface_confirm_read.go` | `Outcome: row.Outcome, ApprovedBy:` → `Outcome: "", ApprovedBy:` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `IN-APPROVEDBY` | `surface_confirm_read.go` | `ApprovedBy: row.ApprovedBy, Decision:` → `ApprovedBy: "", Decision:` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `IN-DECISION` | `surface_confirm_read.go` | `Decision: row.DecisionRef}` → `Decision: ""}` | `TestConfirmationCardServed` → **KILLED** | suite=KILLED[TestConfirmationCardServed]; suite-nostatic=KILLED[TestConfirmationCardServed] | K | yes |
| `RQ-OLDER` | `surface_approve.go` | `…tx, now, &out); isOlderStore(err) {` → `…tx, now, &out); false && isOlderStore(err) {` | `TestRequestsJSONConfirmations` → **KILLED** | suite=KILLED[TestRequestsJSONConfirmations]; suite-nostatic=KILLED[TestRequestsJSONConfirmations] | K | yes |
| `RQ-ERR` | `surface_approve.go` | `} else if err != nil {⏎d.writeInternalErrorLog(r, err)⏎ou…` → `} else if false {⏎d.writeInternalErrorLog(r, err)⏎out.Con…` | `TestRequestsJSONConfirmations` → **KILLED** | suite=KILLED[TestRequestsJSONConfirmations]; suite-nostatic=KILLED[TestRequestsJSONConfirmations] | K | yes |
| `RQ-EMPTY` | `surface_approve.go` | `if out.Confirmations == nil {⏎` → `if false {⏎` | `TestRequestsJSONConfirmations` → **KILLED** | suite=KILLED[TestRequestsJSONConfirmations]; suite-nostatic=KILLED[TestRequestsJSONConfirmations] | K | yes |
| `RQ-CALL` | `surface_approve.go` | `if err := d.readConfirmations(ctx, now, &out); isOlderSto…` → `if err := error(nil); isOlderStore(err) {` | `TestRequestsJSONConfirmations` → **KILLED** | suite=KILLED[TestRequestsJSONConfirmations]; suite-nostatic=KILLED[TestRequestsJSONConfirmations] | K | yes |
| `RQ-LIMIT` | `surface_confirm_read.go` | `…uests(ctx, now, store.MaxPendingVerbRequests)` → `…uests(ctx, now, 1)` | `TestRequestsJSONConfirmations` → **KILLED** | suite=KILLED[TestRequestsJSONConfirmations]; suite-nostatic=KILLED[TestRequestsJSONConfirmations] | K | yes |
| `RQ-NOW` | `surface_confirm_read.go` | `…rbRequests(ctx, now, store.MaxPendingVerbRequests)` → `…rbRequests(ctx, 0, store.MaxPendingVerbRequests)` | `TestRequestsJSONConfirmations` → **KILLED** | suite=KILLED[TestRequestsJSONConfirmations]; suite-nostatic=KILLED[TestRequestsJSONConfirmations] | K | yes |
| `PC-DIGEST` | `surface_confirm_read.go` | `…6(row.Request); hex.EncodeToString(sum[:]) != row.Digest…` → `…6(row.Request); false && hex.EncodeToString(sum[:]) != r…` | `TestRequestsJSONConfirmations` → **KILLED** | suite=KILLED[TestRequestsJSONConfirmations]; suite-nostatic=KILLED[TestRequestsJSONConfirmations] | K | yes |
| `PC-VALIDATE` | `surface_confirm_read.go` | `if err == nil {⏎err = r.Validate()⏎}⏎` → `` | `TestRequestsJSONConfirmations` → **KILLED** | suite=KILLED[TestRequestsJSONConfirmations]; suite-nostatic=KILLED[TestRequestsJSONConfirmations] | K | yes |
| `PC-DECODE` | `surface_confirm_read.go` | `if err != nil {⏎p.Refused = "the stored request is not a …` → `if false {⏎p.Refused = "the stored request is not a valid…` | `TestRequestsJSONConfirmations` → **KILLED** | suite=KILLED[TestRequestsJSONConfirmations]; suite-nostatic=KILLED[TestRequestsJSONConfirmations] | K | yes |
| `PC-SCOPE` | `surface_confirm_read.go` | `if err != nil {⏎p.Refused = "the stored request's scope i…` → `if false {⏎p.Refused = "the stored request's scope is not…` | `TestRequestsJSONConfirmations` → **KILLED** | suite=KILLED[TestRequestsJSONConfirmations]; suite-nostatic=KILLED[TestRequestsJSONConfirmations] | K | yes |
| `PC-EPISODE` | `surface_confirm_read.go` | `= r.Episode, s.Repository, s.Head, s.SHA, &s.Draft` → `= "", s.Repository, s.Head, s.SHA, &s.Draft` | `TestRequestsJSONConfirmations` → **KILLED** | suite=KILLED[TestRequestsJSONConfirmations]; suite-nostatic=KILLED[TestRequestsJSONConfirmations] | K | yes |
| `PC-REPOSITORY` | `surface_confirm_read.go` | `= r.Episode, s.Repository, s.Head, s.SHA, &s.Draft` → `= r.Episode, s.Remote, s.Head, s.SHA, &s.Draft` | `TestRequestsJSONConfirmations` → **KILLED** | suite=KILLED[TestRequestsJSONConfirmations]; suite-nostatic=KILLED[TestRequestsJSONConfirmations] | K | yes |
| `PC-HEAD` | `surface_confirm_read.go` | `= r.Episode, s.Repository, s.Head, s.SHA, &s.Draft` → `= r.Episode, s.Repository, s.Base, s.SHA, &s.Draft` | `TestRequestsJSONConfirmations` → **KILLED** | suite=KILLED[TestRequestsJSONConfirmations]; suite-nostatic=KILLED[TestRequestsJSONConfirmations] | K | yes |
| `PC-SHA` | `surface_confirm_read.go` | `= r.Episode, s.Repository, s.Head, s.SHA, &s.Draft` → `= r.Episode, s.Repository, s.Head, "", &s.Draft` | `TestRequestsJSONConfirmations` → **KILLED** | suite=KILLED[TestRequestsJSONConfirmations]; suite-nostatic=KILLED[TestRequestsJSONConfirmations] | K | yes |
| `PC-DRAFT` | `surface_confirm_read.go` | `= r.Episode, s.Repository, s.Head, s.SHA, &s.Draft` → `= r.Episode, s.Repository, s.Head, s.SHA, &s.CreateLabel` | `TestRequestsJSONConfirmations` → **KILLED** | suite=KILLED[TestRequestsJSONConfirmations]; suite-nostatic=KILLED[TestRequestsJSONConfirmations] | K | yes |
| `PC-EXPIRES` | `surface_confirm_read.go` | `…D, Expires: row.ExpiresAt}` → `…D, Expires: row.CreatedAt}` | `TestRequestsJSONConfirmations` → **KILLED** | suite=KILLED[TestRequestsJSONConfirmations]; suite-nostatic=KILLED[TestRequestsJSONConfirmations] | K | yes |
| `AM-TAG` | `approve_model.go` | ``json:"confirmations"`` → ``json:"confirmation"`` | `TestRequestsJSONConfirmations` → **KILLED** | suite=KILLED[TestRequestsJSONConfirmations]; suite-nostatic=KILLED[TestRequestsJSONConfirmations] | K | yes |
| `AM-TAG-UNAVAILABLE` | `approve_model.go` | ``json:"confirmationsUnavailable"`` → ``json:"unavailableConfirmations"`` | `TestRequestsJSONConfirmations` → **KILLED** | suite=KILLED[TestRequestsJSONConfirmations]; suite-nostatic=KILLED[TestRequestsJSONConfirmations] | K | yes |

---

## 5. Files (final `git diff --name-status 9733813c HEAD`)

```
M  cmd/ailang-worldd/approve_await_refusals_test.go      (2 fixture sites)
M  cmd/ailang-worldd/approve_passkey_refusals_test.go
M  cmd/ailang-worldd/approve_passkey_test.go
M  cmd/ailang-worldd/approve_sig_test.go
M  cmd/ailang-worldd/approver_label_test.go
M  cmd/ailang-worldd/approver_passkey_test.go
M  cmd/ailang-worldd/approver_rig_test.go
M  host/authority/judge_decision_assert_test.go
M  host/authority/judge_decision_test.go
M  host/authority/judge_test.go                         (2)
M  host/authority/judge_verb_test.go                    (2)
A  host/broker/finish_confirm.go
A  host/broker/finish_confirm_test.go
M  host/daemon/daemon.go
M  host/daemon/surface_approve.go
M  host/daemon/surface_approve_post_test.go             (2)
M  host/daemon/surface_approve_read_test.go             (2)
A  host/daemon/surface_confirm_read.go
A  host/daemon/surface_confirm_read_test.go
D  host/daemon/surface_confirm_unrouted_test.go
M  host/daemon/surface_decide_key_test.go
M  host/store/approvers_test.go
M  host/store/credentials.go
M  host/store/credentials_schema.sql
M  host/store/credentials_test.go
M  host/store/key_decisions_test.go
M  host/store/mint_requests_test.go
A  host/store/verb_requests.go
A  host/store/verb_requests_test.go
M  host/surfacepage/approve_model.go
```

Plus this plan and its JSON under `design_docs/planned/` (the controller's commit). **Production files**:
`host/store/{credentials.go,credentials_schema.sql,verb_requests.go}`, `host/broker/finish_confirm.go`,
`host/daemon/{daemon.go,surface_approve.go,surface_confirm_read.go}`, `host/surfacepage/approve_model.go` — no `cmd/`
production file, no page or script byte, no `.ail`.

---

## 6. Sprint gates (controller, after the last sub-milestone's commit; each measured at the base, §1.1)

One bash shell, never under `nohup`: `bash gates.sh <worktree> $HOME/.ailang/state/world-iter402/gates final <HEAD sha>`
(extract it as `accept.sh` is extracted, with the fence `` ````bash gates.sh ``). Every leg is `tee`'d to
`$HOME/.ailang/state/world-iter402/gates/<leg>.txt`; `ALL_DONE.<sha>` is written last with `FAILS: …`. The executor sets
`GATES_SKIP_VERIFY_GO=1`; the controller runs `verify_go`. Legs: the **touched-package derivation** (printed:
`git diff --name-only 9733813c HEAD -- cmd host` → the `.go` files' directories; the plan's list is `cmd/ailang-worldd
host/authority host/broker host/daemon host/store host/surfacepage`; a difference FAILS), `status`, `verify_ail`,
`gofmt`, `build`, `vet` (`./...`), `fence` (`go test -run '^$' ./...`), `pkgs` (five touched packages whole;
`host/broker` by its finish/approve/scope tests — S3 adds one pure file there, and the whole package is 400+ s with the
directive's exec-probe flakies, red at the base (§1.1); `verify_go` runs it whole), **`race`** (every touched package:
`host/store`, `host/authority`, `host/surfacepage` whole; `host/broker`, `host/daemon` and `cmd/ailang-worldd` on the
tests of the files S3 touches — the whole packages under `-race` exceed the leg's budget at load, and all three are on
the directive's load-sensitive list), **`race3`** (`-race -count=3` over the new store tests:
SQLite busy/lock), **`node`** (S2's nine harnesses, `WORLD_EXEC_NODE`, `--- PASS` per name), **`inv`** (the inventories
by name with `-v`: route tables incl. the new `TestConfirmRoutesRegistered` — and `TestConfirmRoutesNotRegistered` must
be **gone**; `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots`,
`TestProductionGoSurface`, `TestNoNewDeadlineFreeStoreReads`, `TestStoreSurfaceTakesContext`, every credentials-schema
test, `TestDaemonCredentialWriterSites`, `TestEveryVerbHasHelp`, `TestFrozenV1RouteTableMatchesMux`, the A2/A3/A4 page
and route tests, S2's surfacepage static tests, S1b's judges, the broker's finish tests), **`chrome`** (V2b's fourteen +
XC1 + the two `host/daemon` drills; S3 changes no page byte, but the card is now daemon-served), then `verify_go`
(`VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh`) unless skipped.

````bash gates.sh
#!/bin/bash
# Row 243 V3 S3 gates (plan §6). usage: gates.sh <tree> <logdir> <base|final> [<sha>]
# The controller: bash gates.sh <worktree> $HOME/.ailang/state/world-iter402/gates final <HEAD sha>
# macOS /bin/bash 3.2: no declare -A. Every leg records its own rc; every leg runs (no set -e); every leg is tee'd
# to <logdir>/<leg>.txt. GATES_SKIP_VERIFY_GO=1 skips verify_go (the executor sets it; the controller runs it).
# Never start this under nohup. Writes <logdir>/ALL_DONE.<sha> with FAILS: when every leg has finished.
T=$1; L=$2; MODE=$3
mkdir -p "$L"; cd "$T" || exit 2
SHA=${4:-$(git rev-parse HEAD 2>/dev/null || echo unknown)}
B=9733813cd8ab41313c158f83d3a0d2bc7c7d7d7f
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
# host/store host/broker host/daemon host/surfacepage host/authority cmd/ailang-worldd; a difference FAILS.
TOUCHED=$(git diff --name-only "$B" HEAD -- cmd host 2>/dev/null | grep '\.go$' | xargs -n1 dirname 2>/dev/null | sort -u | tr '\n' ' ')
echo "TOUCHED_PACKAGES: $TOUCHED" | tee -a "$L/legs.txt"
if [ "$MODE" = final ] && [ "$TOUCHED" != "cmd/ailang-worldd host/authority host/broker host/daemon host/store host/surfacepage " ]; then FAILS="$FAILS touched_packages_differ"; fi
leg status 'git status --porcelain | wc -l'
leg verify_ail './scripts/verify_ail.sh'
leg gofmt 'out=$(gofmt -l cmd host); echo "$out"; test -z "$out"'
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./... > '$L/fence.raw' 2>&1; rc=\$?; cat '$L/fence.raw'; if grep -E '^(FAIL|--- FAIL)' '$L/fence.raw' >/dev/null; then exit 1; fi; exit \$rc"
# host/broker: S3 adds one pure file; the whole package (400+ s, the directive's exec-probe flakies) is verify_go's
leg pkgs "go test -count=1 -timeout 30m ./host/store/ ./host/authority/ ./host/surfacepage/ ./host/daemon/ ./cmd/ailang-worldd/ && go test -count=1 ./host/broker/ -run '^(TestFinish|TestApprove|TestParse|TestRecordOperator)'"
# -race on every touched package: three whole; host/broker, host/daemon and cmd/ailang-worldd on the tests of the files S3 touches
leg race "go test -race -count=1 ./host/store/ ./host/authority/ ./host/surfacepage/ &&
go test -race -count=1 ./host/broker/ -run '^TestFinish' &&
go test -race -count=1 ./host/daemon/ -run '^(TestConfirm|TestRequestsJSON|TestRequestCardJSON|TestApprovePOST|TestApprovalRoutes|TestDecidePOST|TestDecisionRoutes|TestDecidePage|TestSurfaceRequests)' &&
go test -race -count=1 ./cmd/ailang-worldd/ -run '^(TestApprove|TestApprover)'"
# SQLite busy/lock behaviour of the new store tests under -race and -count=3 (directive 11)
NEW_STORE="TestCredentialsSchema9VerbRequests TestCredentialsNewerThanNineIsRefused TestPutVerbRequestSupersedes TestPutVerbRequestGrammarCapPrune TestVerbRequestsChecks TestConsumeVerbRequest"
[ "$MODE" = base ] && NEW_STORE=""
if [ -n "$NEW_STORE" ]; then leg race3 "go test -race -count=3 ./host/store/ -run '^($(bar $NEW_STORE)|TestCredentialsMethodsRequireADeadline)\$'"; fi
# node: every page harness, non-skippable (WORLD_EXEC_NODE is set above; a harness that cannot run FAILS)
NODE="TestSurfaceAppPure TestSurfaceScriptPure TestSurfaceAppVerbs TestNeedsYouListsDecisions TestAddressHintsArePinned TestConfirmPageFailsClosed TestConfirmPageRendersText TestConfirmPageMessage TestConfirmPageStates"
leg node "set -o pipefail; go test -v -count=1 ./host/surfacepage/ -run '^($(bar $NODE))\$'"
need node $NODE
if grep -E '^\s*--- (SKIP|FAIL)' "$L/node.txt"; then FAILS="$FAILS node_skip_or_fail"; fi
NEW_D="TestConfirmRoutesRegistered TestConfirmPageAssetsServed TestConfirmationCardServed TestRequestsJSONConfirmations"
OLD_D="TestConfirmRoutesNotRegistered"
NEW_B="TestFinishConfirmRefusal"
if [ "$MODE" = base ]; then NEW_D=""; NEW_B=""; else OLD_D=""; fi
INV_DAEMON="TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestDaemonCredentialWriterSites TestSurfaceVerbRouteRegistered TestSurfaceVerbIsReadOnly TestDecidePageAssets TestDecidePageFieldTable TestApprovePageRoutes TestPasskeyPageRoutes TestRequestsJSON TestRequestCardJSON"
INV_BOUNDARY="TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard TestWorkbenchPackageRemainsTransportFree TestPinfetchStaysOutOfStoreReplayAndPublish"
INV_STORE="TestProductionContextRoots TestProductionGoSurface TestStoreSurfaceTakesContext TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext TestCredentialsSchemaMatchesCanonicalDDL TestCredentialsMethodsRequireADeadline TestCredentialsSchema8KeyDecisions TestCredentialsSchema7SurfaceApprovals TestCredentialsUpgradeToV4 TestCredentialsUpgradeToV6 TestCredentialsV4WithoutMintRequestsIsRefused TestCredentialsV5WithoutApproversIsRefused"
INV_VG="TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestVerifygateTestWritesAreForkLocked"
INV_SP="TestSurfaceAppMatchesWorkingCopy TestSurfaceAppShellAssembly TestSurfaceFontEmbed TestConfirmationCard TestConfirmationCardWire TestConfirmPageFieldTable TestConfirmCardGolden TestConfirmPageAssets"
INV_AUTH="TestJudgeVerbSSH TestJudgeVerbAssertion TestRecheckVerb TestJudgeDecisionSSH TestJudgeDecisionAssertion"
INV_BROKER="TestFinishPromptText TestFinishApprovalScopeStrict"
INV_CMD="TestEveryVerbHasHelp TestVerbTableMatchesTheCLISwitch TestFinishPromptIsTheGolden"
leg inv "set -o pipefail
go test -v -count=1 ./host/daemon/ -run '^($(bar $INV_DAEMON)$(alt $NEW_D)$(alt $OLD_D))\$' &&
go test -v -count=1 ./host/boundary/ -run '^($(bar $INV_BOUNDARY))\$' &&
go test -v -count=1 ./host/store/ -run '^($(bar $INV_STORE)$(alt $NEW_STORE))\$' &&
go test -v -count=1 ./host/verifygate/ -run '^($(bar $INV_VG))\$' &&
go test -v -count=1 ./host/surfacepage/ -run '^($(bar $INV_SP))\$' &&
go test -v -count=1 ./host/authority/ -run '^($(bar $INV_AUTH))\$' &&
go test -v -count=1 ./host/broker/ -run '^($(bar $INV_BROKER)$(alt $NEW_B))\$' &&
go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(bar $INV_CMD))\$'"
need inv $INV_DAEMON $NEW_D $OLD_D $INV_BOUNDARY $INV_STORE $NEW_STORE $INV_VG $INV_SP $INV_AUTH $INV_BROKER $NEW_B $INV_CMD
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.txt"; then FAILS="$FAILS inv_skip_or_fail"; fi
if [ "$MODE" = final ] && grep -q 'TestConfirmRoutesNotRegistered' "$L/inv.txt"; then FAILS="$FAILS unrouted_test_still_present"; fi
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.txt")" | tee -a "$L/legs.txt"
# live Chrome: S3 changes no page bytes, but now serves the card from the daemon; V2b's fourteen, the two drills
# and XC1 (S2's card in Chrome) stay green
XB="TestApprovePageInChrome TestApprovePageRefusesASwappedRequest TestApprovePageSSHSIGInChrome TestApprovePageShowsAHostileValueAsText TestPasskeyPageAbortsAPendingCeremony TestPasskeyPageInChrome TestApproveDecideWebAuthnLinkInChrome TestDecidePageApproveWithPasskey TestDecidePageFailsClosed TestDecidePageSSHSIGPaste TestNeedsYouListsPendingRequests TestNeedsYouNotifiesOnlyWhenAskedOnce TestVerbBarEpisodeReportInChrome TestVerbBarWhyInChrome TestConfirmPageShowsTheFinishPromptInChrome"
XD="TestSurfaceLiveChrome TestSurfaceFollowChrome"
CH="'/Applications/Google Chrome.app/Contents/MacOS/Google Chrome'"
leg chrome "set -o pipefail; WORLD_CHROME=$CH WORLD_REQUIRE_PASSKEY_BROWSER=1 WORLD_REQUIRE_SSH_KEYGEN=1 go test ./cmd/ailang-worldd/ -run '^($(bar $XB))\$' -count=1 -v && WORLD_CHROME=$CH go test ./host/daemon/ -run '^($(bar $XD))\$' -count=1 -v"
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
  only writer** of the tree — including M9a's deletion; no hand edit, no `gofmt -w`, no fuzz. A failed `--dry-run` or a
  red acceptance → **stop and report** with the log; do not repair.
- **Never `nohup`, never background a test leg; pass a 600000 ms timeout to every long Bash call.** Each acceptance
  part is one call (C14 measured each).
- **Runs**: **Run 1 = M7a → M7b → M8**, **Run 2 = M9a → M9b** (each patch assumes the previous tree; the controller may
  commit between runs). Per sub-milestone the parts, in order: **M7a** `fast`, `daemon`, `cmd`; **M7b** `fast`; **M8**
  `fast`; **M9a** `fast`, `daemon`; **M9b** `fast`, `daemon`. Each command line carries its env: `accept.sh` exports
  `AILANG_BIN=$HOME/.pinned-ailang/ailang` and `WORLD_EXEC_NODE=/opt/homebrew/bin/node` and unsets `WORLD_CHROME`,
  `WORLD_REQUIRE_PASSKEY_BROWSER` and every golden-update variable — **no S3 sub-milestone runs a Chrome test** (no page
  byte changes; the `chrome` gate leg sets `WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
  WORLD_REQUIRE_PASSKEY_BROWSER=1` itself).
- After each sub-milestone's acceptance, snapshot **cumulatively against the base** to
  `$HOME/.ailang/state/world-iter402/snap/<M>/`: every file added or changed since `9733813c` (repo-relative paths), a
  **`DELETED`** file (one repo path per line; from M9a on it holds `host/daemon/surface_confirm_unrouted_test.go`), a
  **`STATUS`** file, and `MANIFEST.sha256` (the two plan files, `STATUS` and `DELETED` are not in it). **Never an
  `rsync` of whole directories.** The `find` takes only the top-level dirs that exist (every S3 stage has `cmd/` and
  `host/`, none `docs/`), rehearsed on all five (C15):
  ```bash
  bash <<'EOF2'
  set -euo pipefail
  M=M7a   # the sub-milestone just accepted
  B=9733813cd8ab41313c158f83d3a0d2bc7c7d7d7f
  W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter402-row243-v3s3
  D="$HOME/.ailang/state/world-iter402/snap/$M"; rm -rf "$D"; mkdir -p "$D"
  cd "$W"
  { git diff --name-status --no-renames "$B" -- cmd host docs; git ls-files --others --exclude-standard -- cmd host docs | sed 's/^/A\t/'; } | sort -u -t "$(printf '\t')" -k2 > "$D/STATUS"
  awk -F'\t' '$1 != "D" {print $2}' "$D/STATUS" | while read -r p; do mkdir -p "$D/$(dirname "$p")"; cp -p "$p" "$D/$p"; done
  awk -F'\t' '$1 == "D" {print $2}' "$D/STATUS" > "$D/DELETED"
  (cd "$D" && find $(for d in cmd host docs; do [ -d "$d" ] && echo "$d"; done) -type f -print0 | sort -z | xargs -0 shasum -a 256) > "$D/MANIFEST.sha256"
  cat "$D/STATUS"; echo "DELETED: $(wc -l < "$D/DELETED")"; echo "MANIFEST: $(wc -l < "$D/MANIFEST.sha256")"
  EOF2
  ```
  Expected (`STATUS` lines / `DELETED` / MANIFEST): **M7a** 21 / 0 / 21; **M7b** 22 / 0 / 22; **M8** 24 / 0 / 24;
  **M9a** 28 / 1 / 27; **M9b** 30 / 1 / 29. The controller copies each snapshot over the worktree, `git rm`s each
  `DELETED` path, commits it as one commit, and checks the numstat (`git diff -M --numstat 9733813c HEAD -- cmd host
  docs | awk '{s+=$1+$2} END {print s}'`): **M7a 377, M7b 805, M8 913, M9a 1,235, M9b 1,391**. A different number →
  stop.
- `export AILANG_BIN=$HOME/.pinned-ailang/ailang` and `WORLD_EXEC_NODE=/opt/homebrew/bin/node` in every shell that runs Go
  tests (`accept.sh` and `gates.sh` set both). The controller's shell is zsh: run each block under `bash <<'EOF2' … EOF2`.
- **Stop conditions**: a dry-run failure; any `--- FAIL`, `--- SKIP` or `FAIL` line in an acceptance (`ACCEPT_RED`);
  gofmt output; a numstat mismatch; any file outside §5 changed (`git status --short` after each patch must equal the
  cumulative list of §3's paragraphs); `surface_confirm_unrouted_test.go` present after M9a; a numstat above 1,400.
- **Known load-sensitive tests** (the directive's list) are not this sprint's: a red there is re-run once alone
  (`TestWorkspaceExecConfiguredBindsTheEpisodeHandler` in the `daemon` part, `TestCLIRealSubprocessEpisode` in the `cmd`
  part are the likely ones). **Add `TestConvertedStoreMethodsM5a` / `TestConvertedStoreMethodsM5b`** (`host/store`, world
  store context cut-offs, "context deadline exceeded, want uncertain"): measured red under load at the base (§1.1, `M5b`)
  and at M9b's `fast` part (`M5a`), green alone (`-count=3`, 5.6 s) — S3 touches no world-store code.
- D-WORLD-109: nothing touches `~/.ailang/world`, `~/.ssh`, a real `/dev/tty` or a real key; every store is a
  `t.TempDir()` file.

---

## 8. Risks

- **R1 — the 9-line margin.** Any executor-side growth breaks the cap; the plan forbids edits (§7). If a gate forces a
  change, M9b leaves first (`requests.json` is S6's to render; 1,235 then), never a CS or CR1 test.
- **R2 — schema 9 strands nothing, but it is one-way.** A v9 `serve` upgrades the real store at start (W54); older
  binaries then refuse it loudly (CS7, RV7). D-WORLD-118 = A authorises it; the loop never runs `serve` on the real store
  (D-WORLD-109) — the first real upgrade is Mark's next `serve` of a build containing S3.
- **R3 — the design's DDL does not refuse a superseded row with a `decision_ref`** (§1.4 F1): no code writes one; a
  later schema could add the CHECK.
- **R4 — a pending verb request is visible before anything confirms it.** The routes serve it read-only behind the Host
  allowlist (RV6); nothing in S3 can approve it (no POST; `TestDaemonCredentialWriterSites` unchanged).
- **R5 — host/daemon under load.** The `daemon` part runs the whole package (≈ 300 s alone); at load ≥ 25 it can near the
  600 s call budget. The executor re-runs a red known-flaky test alone (§7); if the part itself times out, it reruns once.

---

## 9. Acceptance criteria (the judge scores against these)

1. **Exactly three routes registered, all GET, no POST**: `TestConfirmRoutesRegistered` `--- PASS` (route table:
   `GET /surface/confirm.js`, `GET /surface/confirm/{id}`, `GET /surface/confirmations/{id}`, each behind
   `approvalRoute`; POST/PUT/DELETE/OPTIONS → 405, no CORS header; 421 on a foreign Host; other names 404); `RTE-POST`
   and `RTE-PAGE-ANY` killed by it (§4); `TestApprovalRoutesRegistered`, `TestDecisionRoutesRegistered`,
   `TestFrozenV1RouteTableMatchesMux`, `TestSurfaceAppRoutes` unedited and green.
2. **`surface_confirm_unrouted_test.go` deleted in M9a** (the same patch that registers the routes); `inv` shows no
   `TestConfirmRoutesNotRegistered`.
3. **No CLI change**: no `cmd/` production file in §5; `TestEveryVerbHasHelp`, `TestVerbTableMatchesTheCLISwitch`,
   `TestFinishPromptIsTheGolden` green.
4. **No daemon process change beyond the three registrations and `requests.json`'s read**: no new credentials writer
   (`TestDaemonCredentialWriterSites` unedited, green); every new daemon read is `OpenCredentialsReadOnly` under
   `d.readCtx`; `TestProductionContextRoots`, `TestNoNewDeadlineFreeStoreReads` green unedited.
5. **`surfaceCSP` unchanged**: `TestDecidePageAssets` (pins the literal) green unedited; CR3 asserts the card's page and
   script are served under it with the surface headers.
6. **No `.ail`**: `verify_ail` identical to the base's.
7. **Schema 9, fresh == upgraded**: `TestCredentialsSchema9VerbRequests` compares `pragma_table_info` of every table and
   the whole `verb_requests` DDL of a fresh v9 store with a real v8 store upgraded to 9; `user_version = 9`;
   `CRD-UPGRADE-LINE` and the 12 `COL-*` mutants killed.
8. **Version skew**: v8 → each of the four methods refuses `OlderCredentialsStoreError{Version: 8, Need: 9}` naming the
   restart; v10 → "user_version 10 is newer than this binary's 9" from all three opens, nothing changed.
9. **Every refusal writes nothing**: CS4 (grammar, cap) and CS6 (17 refusals) compare the store before and after; each
   refusal's mutant killed (§4).
10. **`ConsumeVerbRequest` is one transaction with `decision_ref`**: CS6 `confirm-*` (five columns, counter, last use);
    `CON-DECISION-SPLIT`, `CON-USEAPPROVER`, `CON-USEPASSKEY` killed.
11. **VT2**: `TestFinishConfirmRefusal` green; F1–F5 mutants killed; no F6/`DraftOnly` symbol exists (`grep -rn
    'DraftOnly\|F6' host/broker/finish_confirm.go` → nothing).
12. **Exactly these pre-existing test files change** (each: the schema bump turns it red without the edit, C6):
    `host/store/credentials_test.go` (fresh table list, version 9, four deadline rows), `host/store/approvers_test.go`,
    `host/store/mint_requests_test.go` (upgraded table lists), `host/store/key_decisions_test.go` (upgrade target);
    the downgrade fixtures (one `DROP TABLE verb_requests` per site) in `cmd/ailang-worldd/{approve_await_refusals,
    approve_passkey_refusals, approve_passkey, approve_sig, approver_label, approver_passkey, approver_rig}_test.go`,
    `host/authority/{judge_decision_assert, judge_decision, judge, judge_verb}_test.go`,
    `host/daemon/{surface_approve_post, surface_approve_read, surface_decide_key}_test.go`; and **deleted**:
    `host/daemon/surface_confirm_unrouted_test.go`. Every other test file in §5 is new.
13. **CR1–CR4**: the served card equals `surfacepage.ConfirmationCard` of the row (hand-built input) byte for byte in three
    states, no hash, no artifact; the statuses of S3-D7; `requests.json` per CR2.
14. **Size**: `git diff -M --numstat 9733813c HEAD -- cmd host docs` sums to **1,391** (≤ 1,400; the same unfiltered);
    each sub-milestone's numstat as §7.
15. **Mutants**: every `K` row killed by its named test and by `suite-nostatic`; `R` rows killed by the route-table test;
    every `CTL` green; the `E` rows named with reasons; the plan's ids are the ledger's.
16. **Gates**: §6's script `FAILS: none` (or only a directive-listed load-sensitive test, green alone), `verify_go` rc 0.

---

## Appendix C — after this sprint (S3's close; what S4, S5 and S6 inherit)

- **Moved out: nothing** (S3-D1: 1,391 ≤ 1,400). No patch is parked.
- **Residuals carried, unchanged**: S1 Appendix C's note 4 and note 5 / gemini r2 → **S5**; the CX15→CX17 numbering gap
  → **S4** (S3 used no CX ids); design RV1–RV9. **New for the controller**: §1.4 F1 (the DDL accepts a superseded row with a
  `decision_ref`; inert, no writer) — a candidate for a later schema, not a re-design of V3.
- **Audit notes 1–3 closed here** (S3-D3, S3-D4, S3-D5); note 3's daemon half → **S4** (below).
- **What S4 (the POST + the submit half) inherits**: `store.VerbRequestIDPattern` (C1), `VerbRequest` (C7, on the
  read-only handle), `ConsumeVerbRequest(ctx, id, digest, use, approvedBy, artifact, decisionRef, now)` (C15; its
  sentinels map to C15's 409/422 — `ErrVerbRequestConsumed`/`Superseded` → 409, `ErrApproverRevoked`/`KeyChanged`/
  `ErrSignCountRegression` → 422, anything else 500), `OlderCredentialsStoreError{Need: 9}` (C13's 503 Restart),
  `broker.FinishConfirmRefusal` (C10; feed it `DecodeVerb` + `Validate` output, `RequestExpires = row.ExpiresAt`, `Now`),
  the read routes' `handleSurfaceConfirmation` shape and `confirmationInput`; the route table test to extend with the one
  POST (`TestConfirmRoutesRegistered`'s `want` list and its 405 loop, which today asserts POST → 405 on the card's path);
  `putVerb`/`verbFixture` as CX fixtures; **audit note 3's daemon half**: a CX row that supersedes between C14 and C15 →
  409 "superseded by a newer finish of <episode>", one orphaned decision (inert, W62), the row `superseded`.
  `TestDaemonCredentialWriterSites` gains `handleSurfaceConfirm` (AM8). `confirm.js`'s submit half per S2 Appendix C.
- **What S5 (the CLI) inherits**: `PutVerbRequest` (its refusals: `OlderCredentialsStoreError` → CF8's restart text,
  `ErrTooManyPendingVerbRequests`, `ErrBadVerbRequest` for an over-size request — CF9 checks size first in the CLI),
  `VerbRequest` for the poll (on `OpenCredentialsReadOnly`), the row's `outcome`/`decision_ref`/`artifact` for CF6/CF7;
  the context roots `TestProductionContextRoots` will pin (design §7).
- **What S6 inherits**: `requests.json`'s `confirmations` (`id, episode, repository, head, sha, draft, expires` or
  `refused`) and `confirmationsUnavailable` for *Needs you* (CP7).

---

## Appendix P — the patches (apply with §2, in order; generated from the compiled scratch stages; do not edit)

### Patch M7a — schema 9 and the store (design §2.4), the pre-existing test edits the bump forces

````diff M7a
diff --git a/cmd/ailang-worldd/approve_await_refusals_test.go b/cmd/ailang-worldd/approve_await_refusals_test.go
index 496f68e8..a7bb0b4c 100644
--- a/cmd/ailang-worldd/approve_await_refusals_test.go
+++ b/cmd/ailang-worldd/approve_await_refusals_test.go
@@ -145,8 +145,8 @@ func TestApproveAwaitNeedsSchema7AndItsRequest(t *testing.T) {
 		gone  bool     // delete the request after the first wait
 		want  string
 	}{
-		{name: "v6-store", setup: []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `PRAGMA user_version = 6`}, want: "surface approvals need schema 7: restart `ailang-worldd serve`"},
-		{name: "v5-store", setup: []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`}, want: "passkeys need schema 6"},
+		{name: "v6-store", setup: []string{`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `PRAGMA user_version = 6`}, want: "surface approvals need schema 7: restart `ailang-worldd serve`"},
+		{name: "v5-store", setup: []string{`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`}, want: "passkeys need schema 6"},
 		{name: "deleted", gone: true, want: "is no longer in the store"},
 	} {
 		t.Run(arm.name, func(t *testing.T) {
diff --git a/cmd/ailang-worldd/approve_passkey_refusals_test.go b/cmd/ailang-worldd/approve_passkey_refusals_test.go
index de84d002..635fbe7e 100644
--- a/cmd/ailang-worldd/approve_passkey_refusals_test.go
+++ b/cmd/ailang-worldd/approve_passkey_refusals_test.go
@@ -56,7 +56,7 @@ func TestApprovePasskeyRefusalsBeforeAnyEffect(t *testing.T) {
 		{name: "WF2-missing", path: func(t *testing.T) string { return filepath.Join(t.TempDir(), "absent.json") }, want: "--webauthn-assertion: open "},
 		{name: "WF3-not-an-assertion", file: func(*testing.T, approveRig) []byte { return []byte(`{"id":"x"}`) },
 			want: "is not a passkey assertion: " + approval.ErrPasskeyResponse.Error()},
-		{name: "WF4-v5-store", setup: sql(`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`), file: good,
+		{name: "WF4-v5-store", setup: sql(`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`), file: good,
 			want: "passkeys need schema 6"},
 		{name: "WF5-unknown-credential", file: func(t *testing.T, a approveRig) []byte {
 			return approvaltest.NewPasskey(approval.AlgEdDSA, 55).Assert(pkOrigin, pkRPID, requestHashOf(t, a), 2)
diff --git a/cmd/ailang-worldd/approve_passkey_test.go b/cmd/ailang-worldd/approve_passkey_test.go
index cf2dfc38..2406a325 100644
--- a/cmd/ailang-worldd/approve_passkey_test.go
+++ b/cmd/ailang-worldd/approve_passkey_test.go
@@ -96,7 +96,7 @@ func TestApproveByPasskey(t *testing.T) {
 			t.Fatalf("rc %d, terminal %v, daemon calls %d, stderr %q", code, opened, *calls, stderr.String())
 		}
 		raw := rawCredentials(t, a.credPath)
-		for _, stmt := range []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
+		for _, stmt := range []string{`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
 			if _, err := raw.Exec(stmt); err != nil {
 				t.Fatal(err)
 			}
diff --git a/cmd/ailang-worldd/approve_sig_test.go b/cmd/ailang-worldd/approve_sig_test.go
index 3bdfcbf9..ffde9a73 100644
--- a/cmd/ailang-worldd/approve_sig_test.go
+++ b/cmd/ailang-worldd/approve_sig_test.go
@@ -216,7 +216,7 @@ func TestApproveBySignatureRefusalsBeforeAnyEffect(t *testing.T) {
 		{name: "v4-store", setup: func(t *testing.T, a approveRig) {
 			raw := rawCredentials(t, a.credPath)
 			defer func() { _ = raw.Close() }()
-			for _, stmt := range []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`} {
+			for _, stmt := range []string{`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`} {
 				if _, err := raw.Exec(stmt); err != nil {
 					t.Fatal(err)
 				}
diff --git a/cmd/ailang-worldd/approver_label_test.go b/cmd/ailang-worldd/approver_label_test.go
index 40ac7ea8..3afb23ea 100644
--- a/cmd/ailang-worldd/approver_label_test.go
+++ b/cmd/ailang-worldd/approver_label_test.go
@@ -82,7 +82,7 @@ func TestApproverListAndRevokeShowPasskeys(t *testing.T) {
 	db5, creds5 := approverStore(t)
 	seedApprover(t, creds5, "s", fixtureKeyOf(t, "ed25519.pub"), approveT0)
 	raw = rawCredentials(t, creds5)
-	for _, stmt := range []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
+	for _, stmt := range []string{`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
 		if _, err := raw.Exec(stmt); err != nil {
 			t.Fatal(err)
 		}
diff --git a/cmd/ailang-worldd/approver_passkey_test.go b/cmd/ailang-worldd/approver_passkey_test.go
index 4f95f654..daeb0e52 100644
--- a/cmd/ailang-worldd/approver_passkey_test.go
+++ b/cmd/ailang-worldd/approver_passkey_test.go
@@ -154,7 +154,7 @@ func TestPasskeyEnrollFlagRules(t *testing.T) {
 		{name: "N6-v5", setup: func(t *testing.T, r enrolRig) {
 			raw := rawCredentials(t, r.creds)
 			defer func() { _ = raw.Close() }()
-			for _, stmt := range []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
+			for _, stmt := range []string{`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
 				if _, err := raw.Exec(stmt); err != nil {
 					t.Fatal(err)
 				}
diff --git a/cmd/ailang-worldd/approver_rig_test.go b/cmd/ailang-worldd/approver_rig_test.go
index b3a176a0..f116ceeb 100644
--- a/cmd/ailang-worldd/approver_rig_test.go
+++ b/cmd/ailang-worldd/approver_rig_test.go
@@ -64,7 +64,7 @@ func approverStoreV4(t *testing.T) (db, creds string) {
 	db, creds = approverStore(t)
 	raw := rawCredentials(t, creds)
 	defer func() { _ = raw.Close() }()
-	for _, stmt := range []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`} {
+	for _, stmt := range []string{`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`} {
 		if _, err := raw.Exec(stmt); err != nil {
 			t.Fatal(err)
 		}
diff --git a/host/authority/judge_decision_assert_test.go b/host/authority/judge_decision_assert_test.go
index 8da37914..da754fb6 100644
--- a/host/authority/judge_decision_assert_test.go
+++ b/host/authority/judge_decision_assert_test.go
@@ -40,7 +40,7 @@ func TestJudgeDecisionAssertion(t *testing.T) {
 		approver string
 	}{
 		{name: "JD-parse", artifact: []byte("{}"), step: JudgeStepParse, cause: approval.ErrPasskeyResponse},
-		{name: "JD-store-v5", setup: []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`},
+		{name: "JD-store-v5", setup: []string{`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`},
 			artifact: jdAssert(t, jgPasskey, jdIn, 2), step: JudgeStepStore},
 		{name: "JD-unknown-key", artifact: jdAssert(t, stranger, jdIn, 2), step: JudgeStepUnknownKey, cause: ErrNotEnrolled, rawID: stranger.CredentialID()},
 		{name: "JD-stored-key", setup: []string{`UPDATE approvers SET public_key = X'00' WHERE name = 'mark-mac'`}, artifact: jdAssert(t, jgPasskey, jdIn, 2),
diff --git a/host/authority/judge_decision_test.go b/host/authority/judge_decision_test.go
index 18d127b8..19c7fad4 100644
--- a/host/authority/judge_decision_test.go
+++ b/host/authority/judge_decision_test.go
@@ -39,7 +39,7 @@ func TestJudgeDecisionSSH(t *testing.T) {
 		field    string
 	}{
 		{name: "JD-parse", artifact: []byte("not a signature\n"), in: jdIn, step: JudgeStepParse, cause: approval.ErrArmor},
-		{name: "JD-store-v4", setup: []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`},
+		{name: "JD-store-v4", setup: []string{`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`},
 			artifact: good, in: jdIn, step: JudgeStepStore},
 		{name: "JD-unknown-key", artifact: jdSign(t, stranger, approval.DecisionNamespace, jdIn), in: jdIn, step: JudgeStepUnknownKey, cause: ErrNotEnrolled,
 			field: stranger.PublicKey().Fingerprint()},
diff --git a/host/authority/judge_test.go b/host/authority/judge_test.go
index 2af94a6b..efef79eb 100644
--- a/host/authority/judge_test.go
+++ b/host/authority/judge_test.go
@@ -174,7 +174,7 @@ func TestJudgeAssertionIsTheCLIJudge(t *testing.T) {
 	arms := append(requestArms(good),
 		judgeArm{name: "WF3-not-an-assertion", artifact: func(*testing.T, store.MintRequestRow) []byte { return []byte(`{"id":"x"}`) },
 			step: JudgeStepParse, cause: approval.ErrPasskeyResponse},
-		judgeArm{name: "WF4-v5-store", setup: []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`}, artifact: good, step: JudgeStepStore},
+		judgeArm{name: "WF4-v5-store", setup: []string{`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`}, artifact: good, step: JudgeStepStore},
 		judgeArm{name: "WF5-unknown-credential", artifact: func(t *testing.T, row store.MintRequestRow) []byte {
 			return stranger.Assert(jgOrigin, jgRPID, jgHash(t, row), 2)
 		}, step: JudgeStepUnknownKey, cause: ErrNotEnrolled, field: string(stranger.CredentialID())},
@@ -218,7 +218,7 @@ func TestJudgeSSHSignatureIsTheCLIJudge(t *testing.T) {
 	arms := append(requestArms(good),
 		judgeArm{name: "F11-garbage", artifact: func(*testing.T, store.MintRequestRow) []byte { return []byte("not a signature\n") },
 			step: JudgeStepParse, cause: approval.ErrArmor},
-		judgeArm{name: "F12-v4-store", setup: []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`},
+		judgeArm{name: "F12-v4-store", setup: []string{`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`},
 			artifact: good, step: JudgeStepStore},
 		judgeArm{name: "F13-unenrolled-key", artifact: func(t *testing.T, row store.MintRequestRow) []byte {
 			return stranger.Sign(approval.Namespace, jgMessage(t, row))
diff --git a/host/authority/judge_verb_test.go b/host/authority/judge_verb_test.go
index 9136e5d7..bf7cdb72 100644
--- a/host/authority/judge_verb_test.go
+++ b/host/authority/judge_verb_test.go
@@ -17,8 +17,8 @@ import (
 
 var vjIn = VerbJudgeInput{Digest: [32]byte{0x76, 1}, Challenge: [32]byte{0x76, 2}, BytesDigest: [32]byte{0x76, 1}, CreatedAt: 1000, ExpiresAt: 1600, Now: 1010}
 
-const vjDropV4 = `DROP TABLE key_decisions; DROP TABLE surface_approvals; DROP TABLE webauthn_credentials; DROP TABLE approvers; PRAGMA user_version = 4`
-const vjDropV5 = `DROP TABLE key_decisions; DROP TABLE surface_approvals; DROP TABLE webauthn_credentials; PRAGMA user_version = 5`
+const vjDropV4 = `DROP TABLE verb_requests; DROP TABLE key_decisions; DROP TABLE surface_approvals; DROP TABLE webauthn_credentials; DROP TABLE approvers; PRAGMA user_version = 4`
+const vjDropV5 = `DROP TABLE verb_requests; DROP TABLE key_decisions; DROP TABLE surface_approvals; DROP TABLE webauthn_credentials; PRAGMA user_version = 5`
 
 // vjSign and vjAssert are a confirmation of in's request by s or pk.
 func vjSign(t *testing.T, s *approvaltest.Signer, namespace string, in VerbJudgeInput) []byte {
diff --git a/host/daemon/surface_approve_post_test.go b/host/daemon/surface_approve_post_test.go
index 0a4f3a0c..987bf1d9 100644
--- a/host/daemon/surface_approve_post_test.go
+++ b/host/daemon/surface_approve_post_test.go
@@ -140,7 +140,7 @@ func TestApprovePOSTRefusalArms(t *testing.T) {
 		}},
 		{"R11-store-v5", "", "Restart", "restart `ailang-worldd serve`", 503, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
 			body := approveBody(row.Digest, assertionFor(t, apPasskey, row, 2))
-			setVersion(t, a.d.credentialsPath, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`)
+			setVersion(t, a.d.credentialsPath, `DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`)
 			return body, nil
 		}},
 		{"R13-busy", "", "Busy", "busy", 503, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
@@ -152,7 +152,7 @@ func TestApprovePOSTRefusalArms(t *testing.T) {
 			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), nil
 		}},
 		{"R14-v6", "", "Restart", "restart `ailang-worldd serve`", 503, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
-			setVersion(t, a.d.credentialsPath, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `PRAGMA user_version = 6`)
+			setVersion(t, a.d.credentialsPath, `DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `PRAGMA user_version = 6`)
 			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), nil
 		}},
 		{"R14-race", "", "Conflict", "changed while it was approved", 409, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
diff --git a/host/daemon/surface_approve_read_test.go b/host/daemon/surface_approve_read_test.go
index becea07a..8a4d5151 100644
--- a/host/daemon/surface_approve_read_test.go
+++ b/host/daemon/surface_approve_read_test.go
@@ -63,7 +63,7 @@ func TestRequestsJSON(t *testing.T) {
 			t.Errorf("requests.json carries %q", leak)
 		}
 	}
-	setVersion(t, a.d.credentialsPath, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `PRAGMA user_version = 6`) // below 7: listed, nothing recorded
+	setVersion(t, a.d.credentialsPath, `DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `PRAGMA user_version = 6`) // below 7: listed, nothing recorded
 	if body := a.serve("GET", "/surface/requests.json", apLoopback, nil).Body.String(); strings.Count(body, `"recorded":null`) != 3 || !strings.Contains(body, `"unavailable":""`) {
 		t.Fatalf("a schema-6 store: %s", body)
 	}
@@ -128,7 +128,7 @@ func TestRequestCardJSON(t *testing.T) {
 	if err := json.Unmarshal(a.serve("GET", "/surface/requests/"+row.RequestID, apLoopback, nil).Body.Bytes(), &c); err != nil || !c.Consumed || c.ApprovedBy != "tty" {
 		t.Fatalf("consumed card = %+v, %v", c, err)
 	}
-	setVersion(t, a.d.credentialsPath, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`)
+	setVersion(t, a.d.credentialsPath, `DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`)
 	if rec := a.serve("GET", "/surface/requests/"+row.RequestID, apLoopback, nil); rec.Code != 200 || !strings.Contains(rec.Body.String(), `"passkeys":[]`) {
 		t.Fatalf("a schema-5 store's card = %d %s, want 200 with no passkeys", rec.Code, rec.Body.String())
 	}
diff --git a/host/daemon/surface_decide_key_test.go b/host/daemon/surface_decide_key_test.go
index 0218dc82..29282fd5 100644
--- a/host/daemon/surface_decide_key_test.go
+++ b/host/daemon/surface_decide_key_test.go
@@ -140,7 +140,7 @@ func TestDecidePOSTSSHSIGRefusals(t *testing.T) {
 			t.Fatalf("%s: evidence was recorded", tc.name)
 		}
 	}
-	for _, drop := range [][]string{{`DROP TABLE key_decisions`, `PRAGMA user_version = 7`}, // v7: the writer's K1, then v4: the judge's store step
+	for _, drop := range [][]string{{`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `PRAGMA user_version = 7`}, // v7: the writer's K1, then v4: the judge's store step
 		{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`}} {
 		setVersion(t, r.d.credentialsPath, drop...)
 		requireRefusal(t, drop[len(drop)-1], r.postKey(q, good), 503, "Restart")
diff --git a/host/store/approvers_test.go b/host/store/approvers_test.go
index 3d1baea8..e7c3f70a 100644
--- a/host/store/approvers_test.go
+++ b/host/store/approvers_test.go
@@ -131,7 +131,7 @@ func TestCredentialsUpgradeToV6(t *testing.T) {
 			}
 			db := rawDB(t, path)
 			ddl := tableDDL(t, db)
-			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "key_decisions", "mint_requests", "session_credentials", "surface_approvals", "webauthn_credentials"})
+			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "key_decisions", "mint_requests", "session_credentials", "surface_approvals", "verb_requests", "webauthn_credentials"})
 			if got, want := normalizeDDL(ddl["approvers"]), normalizeDDL(canonicalApproversDDL); got != want {
 				t.Fatalf("upgraded approvers DDL:\n got: %s\nwant: %s", got, want)
 			}
diff --git a/host/store/credentials.go b/host/store/credentials.go
index 092bb0f1..546076d9 100644
--- a/host/store/credentials.go
+++ b/host/store/credentials.go
@@ -40,8 +40,9 @@ var credentialsSchemaSQL string
 // profile bindings; version 3 adds immutable attended episode start history;
 // version 4 adds mint requests (row 201); version 5 adds approvers (row 201 A1b); version 6 adds
 // passkeys (row 201 A2); version 7 adds the surface's recorded approvals (row 201 A3); version 8
-// adds the key decisions of Human.Approve requests (row 201 A4).
-const credentialsSchemaVersion = 8
+// adds the key decisions of Human.Approve requests (row 201 A4); version 9 adds the verb requests a
+// key confirms on the surface (row 243 V3).
+const credentialsSchemaVersion = 9
 
 // CredentialsSchemaVersion exports credentialsSchemaVersion: serve compares a
 // store's version with it to decide whether to upgrade.
@@ -72,7 +73,9 @@ const credentialsTable = "session_credentials"
 // credentialsTables is the exact, sorted application table set of a credentials store at version.
 func credentialsTables(version int) []string {
 	switch {
-	case version >= 8:
+	case version >= 9:
+		return []string{approversTable, "episode_starts", keyDecisionsTable, mintRequestsTable, credentialsTable, surfaceApprovalsTable, verbRequestsTable, passkeysTable}
+	case version == 8:
 		return []string{approversTable, "episode_starts", keyDecisionsTable, mintRequestsTable, credentialsTable, surfaceApprovalsTable, passkeysTable}
 	case version == 7:
 		return []string{approversTable, "episode_starts", mintRequestsTable, credentialsTable, surfaceApprovalsTable, passkeysTable}
@@ -375,6 +378,9 @@ func upgradeCredentials(display string, db *sql.DB, version int) error {
 	if version < 8 {
 		statements = append(statements, keyDecisionsSQL)
 	}
+	if version < 9 {
+		statements = append(statements, verbRequestsSQL)
+	}
 	for _, stmt := range statements {
 		if _, err := tx.Exec(stmt); err != nil {
 			return fmt.Errorf("store: upgrade credentials %q: %w", display, err)
diff --git a/host/store/credentials_schema.sql b/host/store/credentials_schema.sql
index f3c4bd49..2e6f4cab 100644
--- a/host/store/credentials_schema.sql
+++ b/host/store/credentials_schema.sql
@@ -111,3 +111,23 @@ CREATE TABLE IF NOT EXISTS key_decisions (
     decided_at    INTEGER NOT NULL,
     decision_ref  TEXT CHECK (decision_ref IS NULL OR (length(decision_ref) = 71 AND decision_ref GLOB 'sha256:*'))
 );
+
+-- Verb requests (version 9, row 243 V3). One stored request per `episode finish
+-- --confirm-on-surface`; inert until a verified key confirms it on the surface. The authority is
+-- the world store's finish approval (decision_ref), never this row.
+CREATE TABLE IF NOT EXISTS verb_requests (
+    request_id   TEXT PRIMARY KEY CHECK (length(request_id) = 19 AND request_id GLOB 'vr-[0-9a-f]*'),
+    request      BLOB NOT NULL CHECK (length(request) BETWEEN 1 AND 131072),
+    digest       TEXT NOT NULL CHECK (length(digest) = 64),
+    challenge    TEXT NOT NULL CHECK (length(challenge) = 64),
+    episode_id   TEXT NOT NULL CHECK (episode_id <> ''),
+    created_at   INTEGER NOT NULL,
+    expires_at   INTEGER NOT NULL CHECK (expires_at > created_at),
+    consumed_at  INTEGER,
+    outcome      TEXT CHECK (outcome IS NULL OR outcome IN ('approved', 'superseded')),
+    approved_by  TEXT CHECK (approved_by IS NULL OR approved_by <> ''),
+    artifact     BLOB CHECK (artifact IS NULL OR length(artifact) BETWEEN 1 AND 16384),
+    decision_ref TEXT CHECK (decision_ref IS NULL OR (length(decision_ref) = 71 AND decision_ref GLOB 'sha256:*')),
+    CHECK ((consumed_at IS NULL) = (outcome IS NULL)),
+    CHECK ((outcome IS 'approved') = (approved_by IS NOT NULL AND artifact IS NOT NULL AND decision_ref IS NOT NULL))
+);
diff --git a/host/store/credentials_test.go b/host/store/credentials_test.go
index 39f7afd7..627b9679 100644
--- a/host/store/credentials_test.go
+++ b/host/store/credentials_test.go
@@ -76,7 +76,7 @@ var canonicalPasskeysDDL = `CREATE TABLE webauthn_credentials (
 func TestCredentialsSchemaMatchesCanonicalDDL(t *testing.T) {
 	c := openTestCredentials(t)
 	ddl := tableDDL(t, c.db)
-	requireExactTableNames(t, "fresh credentials store", ddl, []string{"approvers", "episode_starts", "key_decisions", "mint_requests", "session_credentials", "surface_approvals", "webauthn_credentials"})
+	requireExactTableNames(t, "fresh credentials store", ddl, []string{"approvers", "episode_starts", "key_decisions", "mint_requests", "session_credentials", "surface_approvals", "verb_requests", "webauthn_credentials"})
 	for name, got := range map[string]string{"schema file": ddl["webauthn_credentials"], "passkeysSQL": passkeysSQL} {
 		if normalizeDDL(got) != normalizeDDL(canonicalPasskeysDDL) {
 			t.Fatalf("webauthn_credentials DDL (%s):\n got: %s\nwant: %s", name, got, canonicalPasskeysDDL)
@@ -112,8 +112,8 @@ func TestCredentialsSchemaMatchesCanonicalDDL(t *testing.T) {
 		t.Fatalf("v4 world table differs from the version-1 credentials table: %s", got)
 	}
 	var version int
-	if err := c.db.QueryRow("PRAGMA user_version").Scan(&version); err != nil || version != 8 {
-		t.Fatalf("credentials user_version = %d %v, want 8 (row 201 A4)", version, err)
+	if err := c.db.QueryRow("PRAGMA user_version").Scan(&version); err != nil || version != 9 {
+		t.Fatalf("credentials user_version = %d %v, want 9 (row 243 V3)", version, err)
 	}
 }
 
diff --git a/host/store/key_decisions_test.go b/host/store/key_decisions_test.go
index d05d7368..e188c21f 100644
--- a/host/store/key_decisions_test.go
+++ b/host/store/key_decisions_test.go
@@ -85,11 +85,11 @@ func TestCredentialsSchema8KeyDecisions(t *testing.T) {
 		}
 	}
 	_ = v7.Close()
-	if upgraded, err := UpgradeCredentials(path); err != nil || !upgraded || userVersion(t, path) != 8 {
+	if upgraded, err := UpgradeCredentials(path); err != nil || !upgraded || userVersion(t, path) != credentialsSchemaVersion {
 		t.Fatalf("UpgradeCredentials: %v %v, version %d", upgraded, err, userVersion(t, path))
 	}
 	up := tableDDL(t, rawDB(t, path))
-	requireExactTableNames(t, "upgraded v7 store", up, credentialsTables(8))
+	requireExactTableNames(t, "upgraded v7 store", up, credentialsTables(credentialsSchemaVersion))
 	if normalizeDDL(up["key_decisions"]) != normalizeDDL(canonicalKeyDecisionsDDL) {
 		t.Fatalf("upgraded key_decisions DDL: %s", up["key_decisions"])
 	}
diff --git a/host/store/mint_requests_test.go b/host/store/mint_requests_test.go
index 3c3bb8e9..83460f50 100644
--- a/host/store/mint_requests_test.go
+++ b/host/store/mint_requests_test.go
@@ -129,7 +129,7 @@ func TestCredentialsUpgradeToV4(t *testing.T) {
 			}
 			db := rawDB(t, path)
 			ddl := tableDDL(t, db)
-			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "key_decisions", "mint_requests", "session_credentials", "surface_approvals", "webauthn_credentials"})
+			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "key_decisions", "mint_requests", "session_credentials", "surface_approvals", "verb_requests", "webauthn_credentials"})
 			if got, want := normalizeDDL(ddl["mint_requests"]), normalizeDDL(canonicalMintRequestsDDL); got != want {
 				t.Fatalf("upgraded mint_requests DDL:\n got: %s\nwant: %s", got, want)
 			}
diff --git a/host/store/verb_requests.go b/host/store/verb_requests.go
new file mode 100644
index 00000000..18fbbd94
--- /dev/null
+++ b/host/store/verb_requests.go
@@ -0,0 +1,293 @@
+package store
+
+// Row 243 V3 sprint 3 M7 (design_docs/planned/w-world-address-bar-verbs-v3.md §2.4): credentials
+// schema 9's verb_requests, one per `episode finish --confirm-on-surface`, inert until a verified key
+// confirms it; the authority is the world store's finish approval (decision_ref), never this row.
+
+import (
+	"context"
+	"crypto/sha256"
+	"database/sql"
+	"encoding/hex"
+	"errors"
+	"fmt"
+	"regexp"
+)
+
+const verbRequestsTable = "verb_requests"
+
+// verbRequestsSchemaVersion is the first credentials schema with verb_requests.
+const verbRequestsSchemaVersion = 9
+
+// verbRequestsSQL is the version 8 -> 9 upgrade: credentials_schema.sql's verb_requests table,
+// column for column (TestCredentialsSchema9VerbRequests compares a fresh store with an upgraded one).
+const verbRequestsSQL = `CREATE TABLE verb_requests (
+    request_id   TEXT PRIMARY KEY CHECK (length(request_id) = 19 AND request_id GLOB 'vr-[0-9a-f]*'),
+    request      BLOB NOT NULL CHECK (length(request) BETWEEN 1 AND 131072),
+    digest       TEXT NOT NULL CHECK (length(digest) = 64),
+    challenge    TEXT NOT NULL CHECK (length(challenge) = 64),
+    episode_id   TEXT NOT NULL CHECK (episode_id <> ''),
+    created_at   INTEGER NOT NULL,
+    expires_at   INTEGER NOT NULL CHECK (expires_at > created_at),
+    consumed_at  INTEGER,
+    outcome      TEXT CHECK (outcome IS NULL OR outcome IN ('approved', 'superseded')),
+    approved_by  TEXT CHECK (approved_by IS NULL OR approved_by <> ''),
+    artifact     BLOB CHECK (artifact IS NULL OR length(artifact) BETWEEN 1 AND 16384),
+    decision_ref TEXT CHECK (decision_ref IS NULL OR (length(decision_ref) = 71 AND decision_ref GLOB 'sha256:*')),
+    CHECK ((consumed_at IS NULL) = (outcome IS NULL)),
+    CHECK ((outcome IS 'approved') = (approved_by IS NOT NULL AND artifact IS NOT NULL AND decision_ref IS NOT NULL))
+)`
+
+// MaxPendingVerbRequests caps the unconsumed, unexpired verb requests one store holds (design T14).
+const MaxPendingVerbRequests = 16
+
+// MaxVerbRequestBytes is the table's request size bound (approval.MaxVerbRequestBytes, 128 KiB).
+const MaxVerbRequestBytes = 128 << 10
+
+// verbRequestPruneAfter: an unconsumed verb request is deleted once it has been expired this long.
+const verbRequestPruneAfter = 86400
+
+// The outcomes a consumed verb request carries.
+const (
+	VerbOutcomeApproved   = "approved"
+	VerbOutcomeSuperseded = "superseded"
+)
+
+// VerbRequestIDPattern is a verb request id: "vr-" and 16 lowercase hex. The table's CHECK is
+// weaker, so PutVerbRequest enforces this.
+var VerbRequestIDPattern = regexp.MustCompile(`^vr-[0-9a-f]{16}$`)
+
+var (
+	ErrBadVerbRequest             = errors.New("store: invalid verb request row")
+	ErrTooManyPendingVerbRequests = errors.New("store: 16 verb requests are already pending; confirm them or let them expire")
+	ErrVerbRequestUnknown         = errors.New("store: no such verb request")
+	ErrVerbRequestConsumed        = errors.New("store: the verb request was already confirmed")
+	ErrVerbRequestSuperseded      = errors.New("store: the verb request was superseded by a newer finish of its episode")
+	ErrVerbRequestExpired         = errors.New("store: the verb request expired before it was confirmed")
+	ErrVerbRequestClockBehind     = errors.New("store: the clock is behind the verb request's creation time")
+	ErrVerbRequestDigestChanged   = errors.New("store: the verb request's digest is not the one confirmed")
+)
+
+// VerbRequestRow is one verb_requests row.
+type VerbRequestRow struct {
+	RequestID   string
+	Request     []byte // the canonical world/verb-request/v1 bytes (approval.CanonicalVerb)
+	Digest      string // 64 lowercase hex: SHA-256 of Request
+	Challenge   string // 64 lowercase hex, random, single use
+	EpisodeID   string
+	CreatedAt   int64
+	ExpiresAt   int64
+	ConsumedAt  int64  // 0 while pending
+	Outcome     string // "" while pending, else VerbOutcomeApproved or VerbOutcomeSuperseded
+	ApprovedBy  string // "<format>:<name>" of an approved row
+	Artifact    []byte // the verified artifact of an approved row
+	DecisionRef string // the world store's finish approval decision of an approved row
+}
+
+func (c *Credentials) olderForVerbRequests() error {
+	return &OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: verbRequestsSchemaVersion, Feature: "verb requests"}
+}
+
+// sha256Hex is the lowercase hex SHA-256 of b.
+func sha256Hex(b []byte) string {
+	sum := sha256.Sum256(b)
+	return hex.EncodeToString(sum[:])
+}
+
+// checkVerbRequestRow is PutVerbRequest's grammar (S2): each refusal is ErrBadVerbRequest.
+func checkVerbRequestRow(row VerbRequestRow) error {
+	switch {
+	case !VerbRequestIDPattern.MatchString(row.RequestID):
+		return fmt.Errorf("%w: request id %q", ErrBadVerbRequest, row.RequestID)
+	case len(row.Request) < 1 || len(row.Request) > MaxVerbRequestBytes:
+		return fmt.Errorf("%w: a request is 1 to %d bytes, not %d", ErrBadVerbRequest, MaxVerbRequestBytes, len(row.Request))
+	case row.Digest != sha256Hex(row.Request): // so the digest is 64 lowercase hex
+		return fmt.Errorf("%w: the digest is not the request's SHA-256", ErrBadVerbRequest)
+	case !isLowerHex64(row.Challenge):
+		return fmt.Errorf("%w: the challenge is not 64 lowercase hex", ErrBadVerbRequest)
+	case row.EpisodeID == "":
+		return fmt.Errorf("%w: empty episode id", ErrBadVerbRequest)
+	case row.ExpiresAt <= row.CreatedAt:
+		return fmt.Errorf("%w: expires_at %d is not after created_at %d", ErrBadVerbRequest, row.ExpiresAt, row.CreatedAt)
+	case row.ConsumedAt != 0 || row.Outcome != "" || row.ApprovedBy != "" || row.Artifact != nil || row.DecisionRef != "":
+		return fmt.Errorf("%w: a new request is pending", ErrBadVerbRequest)
+	}
+	return nil
+}
+
+// PutVerbRequest stores one new, pending request in one transaction (design §2.4): S1 deadline and
+// schema 9; S2 grammar; S3 the episode's pending request is superseded; S4 requests expired a day
+// ago are pruned; S5 at MaxPendingVerbRequests live ones it refuses; S6 insert.
+func (c *Credentials) PutVerbRequest(ctx context.Context, row VerbRequestRow, now int64) error {
+	if err := requireDeadline(ctx); err != nil { // S1
+		return err
+	}
+	if c.version < verbRequestsSchemaVersion { // S1
+		return c.olderForVerbRequests()
+	}
+	if err := checkVerbRequestRow(row); err != nil { // S2
+		return err
+	}
+	tx, err := c.db.BeginTx(ctx, nil)
+	if err != nil {
+		return fmt.Errorf("store: put verb request: begin: %w", err)
+	}
+	defer func() { _ = tx.Rollback() }()
+	if _, err := tx.ExecContext(ctx, `UPDATE verb_requests SET consumed_at = ?, outcome = 'superseded' WHERE episode_id = ? AND consumed_at IS NULL`, now, row.EpisodeID); err != nil { // S3
+		return fmt.Errorf("store: put verb request: supersede: %w", err)
+	}
+	if _, err := tx.ExecContext(ctx, `DELETE FROM verb_requests WHERE consumed_at IS NULL AND expires_at <= ?`, now-verbRequestPruneAfter); err != nil { // S4
+		return fmt.Errorf("store: put verb request: prune: %w", err)
+	}
+	var live int
+	if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM verb_requests WHERE consumed_at IS NULL AND expires_at > ?`, now).Scan(&live); err != nil {
+		return fmt.Errorf("store: put verb request: count: %w", err)
+	}
+	if live >= MaxPendingVerbRequests { // S5
+		return ErrTooManyPendingVerbRequests
+	}
+	if _, err := tx.ExecContext(ctx, `INSERT INTO verb_requests (request_id, request, digest, challenge, episode_id, created_at, expires_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
+		row.RequestID, row.Request, row.Digest, row.Challenge, row.EpisodeID, row.CreatedAt, row.ExpiresAt); err != nil { // S6
+		return fmt.Errorf("store: put verb request: %w", err)
+	}
+	if err := tx.Commit(); err != nil {
+		return fmt.Errorf("store: put verb request: commit: %w", err)
+	}
+	return nil
+}
+
+const verbRequestColumns = `request_id, request, digest, challenge, episode_id, created_at, expires_at, coalesce(consumed_at, 0), coalesce(outcome, ''), coalesce(approved_by, ''), artifact, coalesce(decision_ref, '')`
+
+func scanVerbRequest(scan func(...any) error) (VerbRequestRow, error) {
+	var r VerbRequestRow
+	err := scan(&r.RequestID, &r.Request, &r.Digest, &r.Challenge, &r.EpisodeID, &r.CreatedAt, &r.ExpiresAt, &r.ConsumedAt, &r.Outcome, &r.ApprovedBy, &r.Artifact, &r.DecisionRef)
+	return r, err
+}
+
+// VerbRequest reads one request by id; ok is false when there is none.
+func (c *Credentials) VerbRequest(ctx context.Context, id string) (VerbRequestRow, bool, error) {
+	if err := requireDeadline(ctx); err != nil {
+		return VerbRequestRow{}, false, err
+	}
+	if c.version < verbRequestsSchemaVersion {
+		return VerbRequestRow{}, false, c.olderForVerbRequests()
+	}
+	row, err := scanVerbRequest(c.db.QueryRowContext(ctx, `SELECT `+verbRequestColumns+` FROM verb_requests WHERE request_id = ?`, id).Scan)
+	if errors.Is(err, sql.ErrNoRows) {
+		return VerbRequestRow{}, false, nil
+	}
+	if err != nil {
+		return VerbRequestRow{}, false, fmt.Errorf("store: read verb request: %w", err)
+	}
+	return row, true, nil
+}
+
+// PendingVerbRequests lists the live, pending requests, newest first. limit <= 0 or above
+// MaxPendingVerbRequests means MaxPendingVerbRequests.
+func (c *Credentials) PendingVerbRequests(ctx context.Context, now int64, limit int) ([]VerbRequestRow, error) {
+	if err := requireDeadline(ctx); err != nil {
+		return nil, err
+	}
+	if c.version < verbRequestsSchemaVersion {
+		return nil, c.olderForVerbRequests()
+	}
+	if limit <= 0 || limit > MaxPendingVerbRequests {
+		limit = MaxPendingVerbRequests
+	}
+	rows, err := c.db.QueryContext(ctx, `SELECT `+verbRequestColumns+` FROM verb_requests WHERE consumed_at IS NULL AND expires_at > ? ORDER BY created_at DESC, request_id DESC LIMIT ?`, now, limit)
+	if err != nil {
+		return nil, fmt.Errorf("store: list verb requests: %w", err)
+	}
+	defer func() { _ = rows.Close() }()
+	var out []VerbRequestRow
+	for rows.Next() {
+		row, err := scanVerbRequest(rows.Scan)
+		if err != nil {
+			return nil, fmt.Errorf("store: list verb requests: %w", err)
+		}
+		out = append(out, row)
+	}
+	if err := rows.Err(); err != nil {
+		return nil, fmt.Errorf("store: list verb requests: %w", err)
+	}
+	return out, nil
+}
+
+// checkVerbConsume is ConsumeVerbRequest's grammar: each refusal is ErrBadVerbRequest.
+func checkVerbConsume(id, digest string, use ApproverUse, approvedBy string, artifact []byte, decisionRef string) error {
+	switch {
+	case !VerbRequestIDPattern.MatchString(id) || !isLowerHex64(digest):
+		return fmt.Errorf("%w: request %q, digest %q", ErrBadVerbRequest, id, digest)
+	case approvedBy != use.Format+":"+use.Name:
+		return fmt.Errorf("%w: %q is not the approver %s:%s", ErrBadVerbRequest, approvedBy, use.Format, use.Name)
+	case (use.Format == "webauthn") != (use.Passkey != nil): // the counter binding is exactly a passkey's (A2 C13)
+		return fmt.Errorf("%w: a %s confirmation with passkey use %v", ErrBadVerbRequest, use.Format, use.Passkey != nil)
+	case len(artifact) < 1 || len(artifact) > MaxSurfaceApprovalArtifactBytes:
+		return fmt.Errorf("%w: an artifact is 1 to %d bytes, not %d", ErrBadVerbRequest, MaxSurfaceApprovalArtifactBytes, len(artifact))
+	case !isDecisionRef(decisionRef):
+		return fmt.Errorf("%w: decision %q", ErrBadVerbRequest, decisionRef)
+	}
+	return nil
+}
+
+// ConsumeVerbRequest confirms request id at now in one transaction, after the world store's finish
+// approval decisionRef landed (design §2.4, §2.6 C15): the row is re-read and refused in Admit's
+// order (consumed or superseded, clock behind, expired, digest); the approver is re-checked and
+// marked used (useApprover) and, for a passkey, its counter moved forward (usePasskey); then one
+// UPDATE sets consumed_at, outcome, approved_by, artifact and decision_ref. Any refusal writes
+// nothing.
+func (c *Credentials) ConsumeVerbRequest(ctx context.Context, id, digest string, use ApproverUse, approvedBy string, artifact []byte, decisionRef string, now int64) error {
+	if err := requireDeadline(ctx); err != nil {
+		return err
+	}
+	if c.version < verbRequestsSchemaVersion {
+		return c.olderForVerbRequests()
+	}
+	if err := checkVerbConsume(id, digest, use, approvedBy, artifact, decisionRef); err != nil {
+		return err
+	}
+	tx, err := c.db.BeginTx(ctx, nil)
+	if err != nil {
+		return fmt.Errorf("store: consume verb request: begin: %w", err)
+	}
+	defer func() { _ = tx.Rollback() }()
+	row, err := scanVerbRequest(tx.QueryRowContext(ctx, `SELECT `+verbRequestColumns+` FROM verb_requests WHERE request_id = ?`, id).Scan)
+	if errors.Is(err, sql.ErrNoRows) {
+		return fmt.Errorf("%w: %s", ErrVerbRequestUnknown, id)
+	}
+	if err != nil {
+		return fmt.Errorf("store: consume verb request: %w", err)
+	}
+	switch {
+	case row.Outcome == VerbOutcomeSuperseded:
+		return fmt.Errorf("%w: %s", ErrVerbRequestSuperseded, id)
+	case row.Outcome != "":
+		return fmt.Errorf("%w: %s by %s", ErrVerbRequestConsumed, id, row.ApprovedBy)
+	case now < row.CreatedAt:
+		return fmt.Errorf("%w: %s", ErrVerbRequestClockBehind, id)
+	case now >= row.ExpiresAt:
+		return fmt.Errorf("%w: %s", ErrVerbRequestExpired, id)
+	case digest != row.Digest || sha256Hex(row.Request) != row.Digest:
+		return fmt.Errorf("%w: %s", ErrVerbRequestDigestChanged, id)
+	}
+	if err := useApprover(ctx, tx, use, now); err != nil {
+		return err
+	}
+	if use.Passkey != nil {
+		if err := usePasskey(ctx, tx, use.Name, *use.Passkey); err != nil {
+			return err
+		}
+	}
+	res, err := tx.ExecContext(ctx, `UPDATE verb_requests SET consumed_at = ?, outcome = 'approved', approved_by = ?, artifact = ?, decision_ref = ? WHERE request_id = ? AND consumed_at IS NULL`,
+		now, approvedBy, artifact, decisionRef, id)
+	if err != nil {
+		return fmt.Errorf("store: consume verb request: %w", err)
+	}
+	if n, err := res.RowsAffected(); err != nil || n != 1 {
+		return fmt.Errorf("%w: %s (%d rows, %v)", ErrVerbRequestConsumed, id, n, err)
+	}
+	if err := tx.Commit(); err != nil {
+		return fmt.Errorf("store: consume verb request: commit: %w", err)
+	}
+	return nil
+}
````

### Patch M7b — CS1–CS7 (design §6 Sprint 3)

````diff M7b
diff --git a/host/store/credentials_test.go b/host/store/credentials_test.go
index 627b9679..c7b85c6f 100644
--- a/host/store/credentials_test.go
+++ b/host/store/credentials_test.go
@@ -303,6 +303,10 @@ func TestCredentialsMethodsRequireADeadline(t *testing.T) {
 	errs["RecordKeyDecision"] = c.RecordKeyDecision(ctx, KeyDecisionRow{}, ApproverUse{}, 1000)
 	errs["LandKeyDecision"] = c.LandKeyDecision(ctx, "x", "y")
 	_, _, errs["KeyDecision"] = c.KeyDecision(ctx, "x")
+	errs["PutVerbRequest"] = c.PutVerbRequest(ctx, VerbRequestRow{}, 1000)
+	_, _, errs["VerbRequest"] = c.VerbRequest(ctx, "x")
+	_, errs["PendingVerbRequests"] = c.PendingVerbRequests(ctx, 1000, 0)
+	errs["ConsumeVerbRequest"] = c.ConsumeVerbRequest(ctx, "x", "y", ApproverUse{}, "", nil, "", 1000)
 	for name, err := range errs {
 		if !errors.Is(err, ErrNoDeadline) {
 			t.Errorf("%s(Background) = %v, want ErrNoDeadline", name, err)
diff --git a/host/store/verb_requests_test.go b/host/store/verb_requests_test.go
new file mode 100644
index 00000000..ac29d374
--- /dev/null
+++ b/host/store/verb_requests_test.go
@@ -0,0 +1,424 @@
+package store
+
+// Row 243 V3 sprint 3 M7 (design_docs/planned/w-world-address-bar-verbs-v3.md §2.4, §6 CS1-CS7):
+// credentials schema 9's verb_requests. Every store is under t.TempDir() (D-WORLD-109); every
+// request's bytes are approval.CanonicalVerb's.
+
+import (
+	"bytes"
+	"database/sql"
+	"errors"
+	"fmt"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+var vrDecision = "sha256:" + strings.Repeat("d4", 32)
+
+// vrRow is a pending verb request vr-<16 hex of n> for episode ep, created at created and expiring
+// 600 s later: its bytes are approval.CanonicalVerb of a valid request, its digest their SHA-256.
+func vrRow(t *testing.T, n int, ep string, created int64) VerbRequestRow {
+	t.Helper()
+	r := approval.VerbRequest{Schema: approval.VerbRequestSchema, Verb: approval.VerbFinish, Argv: []string{"ailang-worldd", "episode", "finish", ep, "--preset", "p"},
+		Episode: ep, Prompt: fmt.Sprintf("committer: (repository git config)\nrun %d\n", n), Scope: "github-pr:o/r#scope", Title: "Fix the bar", Body: "body", Protected: []string{}}
+	if err := r.Validate(); err != nil {
+		t.Fatal(err)
+	}
+	b, err := approval.CanonicalVerb(r)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return VerbRequestRow{RequestID: fmt.Sprintf("vr-%016x", n), Request: b, Digest: sha256Hex(b), Challenge: sha256Hex([]byte(fmt.Sprint(n))),
+		EpisodeID: ep, CreatedAt: created, ExpiresAt: created + 600}
+}
+
+// vrStore is a fresh schema-9 store with ssh approver mark-ssh (key k1) and passkey x (counter 5)
+// enrolled at 900.
+func vrStore(t *testing.T) *Credentials {
+	t.Helper()
+	c := openTestCredentials(t)
+	if err := c.EnrolApprover(testCtx(t), ApproverRow{Name: "mark-ssh", Format: "ssh", PublicKey: []byte("k1"), EnrolledAt: 900}); err != nil {
+		t.Fatal(err)
+	}
+	a, p := testPasskey("x", 1)
+	a.EnrolledAt, p.CredentialID, p.SignCount = 900, credX, 5
+	if err := c.EnrolPasskey(testCtx(t), a, p); err != nil {
+		t.Fatal(err)
+	}
+	return c
+}
+
+func vrPut(t *testing.T, c *Credentials, row VerbRequestRow, now int64) {
+	t.Helper()
+	if err := c.PutVerbRequest(testCtx(t), row, now); err != nil {
+		t.Fatalf("put %s: %v", row.RequestID, err)
+	}
+}
+
+func vrGet(t *testing.T, c *Credentials, id string) VerbRequestRow {
+	t.Helper()
+	r, ok, err := c.VerbRequest(testCtx(t), id)
+	if err != nil || !ok {
+		t.Fatalf("read %s: ok %v err %v", id, ok, err)
+	}
+	return r
+}
+
+// vrState is every row of verb_requests, the approvers' last use and revocation, and the passkeys' counters.
+func vrState(t *testing.T, c *Credentials) string {
+	t.Helper()
+	var out string
+	if err := c.db.QueryRow(`SELECT coalesce((SELECT group_concat(quote(request_id)||quote(request)||quote(consumed_at)||quote(outcome)||quote(approved_by)||quote(artifact)||quote(decision_ref), ';') FROM verb_requests), '')
+		|| '|' || (SELECT group_concat(name||quote(last_used_at)||quote(revoked_at), ';') FROM approvers) || '|' || (SELECT group_concat(name||sign_count, ';') FROM webauthn_credentials)`).Scan(&out); err != nil {
+		t.Fatal(err)
+	}
+	return out
+}
+
+// tableInfo is PRAGMA table_info of every table of db, in order.
+func tableInfo(t *testing.T, db *sql.DB) (out []string) {
+	t.Helper()
+	rows, err := db.Query(`SELECT m.name, p.name, p.type, p."notnull", quote(p.dflt_value), p.pk FROM sqlite_master m, pragma_table_info(m.name) p WHERE m.type = 'table' ORDER BY m.name, p.cid`)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = rows.Close() }()
+	for rows.Next() {
+		var table, col, typ, def string
+		var notNull, pk int
+		if err := rows.Scan(&table, &col, &typ, &notNull, &def, &pk); err != nil {
+			t.Fatal(err)
+		}
+		out = append(out, fmt.Sprintf("%s.%s %s notnull=%d default=%s pk=%d", table, col, typ, notNull, def, pk))
+	}
+	return out
+}
+
+// CS1, CS2, CS7(b): a fresh store is schema 9; a real schema-8 store (makeV6Credentials, then the
+// version 7 and 8 tables) opens at 8 and refuses every verb request method naming the serve
+// restart; UpgradeCredentials takes it to 9, and then every table's columns, and verb_requests'
+// whole DDL, equal the fresh store's.
+func TestCredentialsSchema9VerbRequests(t *testing.T) {
+	fresh := openTestCredentials(t)
+	if got, want := credentialsTables(9), []string{"approvers", "episode_starts", "key_decisions", "mint_requests", "session_credentials", "surface_approvals", "verb_requests", "webauthn_credentials"}; !reflect.DeepEqual(got, want) {
+		t.Fatalf("credentialsTables(9) = %v", got)
+	}
+	if fresh.Version() != 9 || CredentialsSchemaVersion != 9 || MaxVerbRequestBytes != approval.MaxVerbRequestBytes {
+		t.Fatalf("fresh version %d, constant %d", fresh.Version(), CredentialsSchemaVersion)
+	}
+	path := filepath.Join(t.TempDir(), CredentialsFileName)
+	makeV6Credentials(t, path, []SessionRow{oneBoundRow()})
+	raw := rawDB(t, path)
+	for _, stmt := range []string{surfaceApprovalsSQL, keyDecisionsSQL} {
+		if _, err := raw.Exec(stmt); err != nil {
+			t.Fatal(err)
+		}
+	}
+	setPositiveFixtureVersion(t, raw, 8)
+	v8, err := OpenCredentials(path)
+	if err != nil || v8.Version() != 8 {
+		t.Fatalf("open v8: %v", err)
+	}
+	ctx := testCtx(t)
+	_, _, readErr := v8.VerbRequest(ctx, "vr-0000000000000001")
+	_, listErr := v8.PendingVerbRequests(ctx, 1000, 0)
+	for name, err := range map[string]error{"put": v8.PutVerbRequest(ctx, vrRow(t, 1, "w1", 1000), 1000), "read": readErr, "list": listErr,
+		"consume": v8.ConsumeVerbRequest(ctx, "vr-0000000000000001", strings.Repeat("a", 64), kdUse, "ssh:mark-ssh", []byte("sig"), vrDecision, 1000)} {
+		var oe *OlderCredentialsStoreError
+		if !errors.As(err, &oe) || oe.Version != 8 || oe.Need != 9 || !strings.Contains(err.Error(), "verb requests need schema 9: restart `ailang-worldd serve`") {
+			t.Errorf("v8 %s: %v, want the restart refusal naming schema 9", name, err)
+		}
+	}
+	_ = v8.Close()
+	if upgraded, err := UpgradeCredentials(path); err != nil || !upgraded || userVersion(t, path) != 9 {
+		t.Fatalf("UpgradeCredentials: %v %v, version %d", upgraded, err, userVersion(t, path))
+	}
+	if again, err := UpgradeCredentials(path); err != nil || again {
+		t.Fatalf("second UpgradeCredentials: %v %v", again, err)
+	}
+	up := rawDB(t, path)
+	requireExactTableNames(t, "upgraded v8 store", tableDDL(t, up), credentialsTables(9))
+	if got, want := tableInfo(t, up), tableInfo(t, fresh.db); !reflect.DeepEqual(got, want) || strings.Count(strings.Join(want, "\n"), "verb_requests.") != 12 {
+		t.Fatalf("upgraded columns differ from a fresh store's:\n got: %v\nwant: %v", got, want)
+	}
+	if got, want := normalizeDDL(tableDDL(t, up)["verb_requests"]), normalizeDDL(tableDDL(t, fresh.db)["verb_requests"]); got != want {
+		t.Fatalf("upgraded verb_requests DDL:\n got: %s\nwant: %s", got, want)
+	}
+}
+
+// CS7(a): a store newer than this binary (user_version 10) is refused by both opens with the
+// existing "newer than this binary's" text, now at 9, and the refusal changes nothing.
+func TestCredentialsNewerThanNineIsRefused(t *testing.T) {
+	path := filepath.Join(t.TempDir(), CredentialsFileName)
+	c, err := OpenCredentials(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	_ = c.Close()
+	db := rawDB(t, path)
+	setPositiveFixtureVersion(t, db, 10)
+	before, _ := schemaState(t, db)
+	for name, open := range map[string]func(string) (*Credentials, error){"writable": OpenCredentials, "read-only": OpenCredentialsReadOnly, "upgrade": func(p string) (*Credentials, error) { _, err := UpgradeCredentials(p); return nil, err }} {
+		c, err := open(path)
+		if c != nil {
+			_ = c.Close()
+		}
+		var se *CredentialsSchemaError
+		if !errors.As(err, &se) || !strings.Contains(err.Error(), "user_version 10 is newer than this binary's 9; use a compatible binary") {
+			t.Errorf("%s open of a v10 store: %v", name, err)
+		}
+	}
+	if after, v := schemaState(t, db); !reflect.DeepEqual(after, before) || v != 10 {
+		t.Fatalf("refused opens changed the store: %v -> %v (version %d)", before, after, v)
+	}
+}
+
+// CS3: a second request for an episode supersedes its pending one (consumed_at = now, outcome
+// superseded); another episode's request stays pending.
+func TestPutVerbRequestSupersedes(t *testing.T) {
+	c := vrStore(t)
+	a, b, n := vrRow(t, 1, "w1", 1000), vrRow(t, 2, "w2", 1000), vrRow(t, 3, "w1", 1010)
+	vrPut(t, c, a, 1000)
+	vrPut(t, c, b, 1000)
+	vrPut(t, c, n, 1010)
+	if got := vrGet(t, c, a.RequestID); got.Outcome != VerbOutcomeSuperseded || got.ConsumedAt != 1010 || got.ApprovedBy != "" || got.DecisionRef != "" || got.Artifact != nil {
+		t.Fatalf("superseded row = %+v", got)
+	}
+	for _, id := range []string{b.RequestID, n.RequestID} {
+		if got := vrGet(t, c, id); got.Outcome != "" || got.ConsumedAt != 0 {
+			t.Fatalf("%s = %+v, want pending", id, got)
+		}
+	}
+	got, err := c.PendingVerbRequests(testCtx(t), 1010, 0)
+	if err != nil || len(got) != 2 || got[0].RequestID != n.RequestID || got[1].RequestID != b.RequestID || !bytes.Equal(got[0].Request, n.Request) {
+		t.Fatalf("pending = %+v, %v", got, err)
+	}
+	vrPut(t, c, vrRow(t, 4, "w4", 1600+86400), 1600+86400) // a day after a's expiry: the prune takes pending rows only
+	vrGet(t, c, a.RequestID)                               // the superseded row is kept as the record
+}
+
+// CS4: the grammar (each refusal ErrBadVerbRequest, nothing written), the 16 cap (a supersede frees
+// its episode's place first), and the prune of requests expired a day ago.
+func TestPutVerbRequestGrammarCapPrune(t *testing.T) {
+	c := vrStore(t)
+	big := bytes.Repeat([]byte("x"), MaxVerbRequestBytes+1)
+	for name, edit := range map[string]func(*VerbRequestRow){
+		"id-upper": func(r *VerbRequestRow) { r.RequestID = "vr-000000000000000A" }, "id-mint": func(r *VerbRequestRow) { r.RequestID = "mr-0000000000000001" },
+		"id-short":         func(r *VerbRequestRow) { r.RequestID = "vr-0001" },
+		"request-empty":    func(r *VerbRequestRow) { r.Request, r.Digest = nil, sha256Hex(nil) },
+		"request-128KiB+1": func(r *VerbRequestRow) { r.Request, r.Digest = big, sha256Hex(big) },
+		"digest-upper":     func(r *VerbRequestRow) { r.Digest = strings.ToUpper(r.Digest) },
+		"challenge-short":  func(r *VerbRequestRow) { r.Challenge = r.Challenge[:62] },
+		"digest-not-sha":   func(r *VerbRequestRow) { r.Digest = r.Challenge },
+		"episode-empty":    func(r *VerbRequestRow) { r.EpisodeID = "" },
+		"expires-created":  func(r *VerbRequestRow) { r.ExpiresAt = r.CreatedAt },
+		"outcome-set":      func(r *VerbRequestRow) { r.Outcome = VerbOutcomeApproved },
+		"consumed-set":     func(r *VerbRequestRow) { r.ConsumedAt = 1000 },
+		"approver-set":     func(r *VerbRequestRow) { r.ApprovedBy = "ssh:mark-ssh" },
+		"artifact-set":     func(r *VerbRequestRow) { r.Artifact = []byte("sig") },
+		"decision-set":     func(r *VerbRequestRow) { r.DecisionRef = vrDecision },
+	} {
+		r := vrRow(t, 1, "w1", 1000)
+		edit(&r)
+		if err := c.PutVerbRequest(testCtx(t), r, 1000); !errors.Is(err, ErrBadVerbRequest) {
+			t.Errorf("%s: %v, want ErrBadVerbRequest", name, err)
+		}
+	}
+	if s := vrState(t, c); strings.Contains(s, "'vr-") {
+		t.Fatalf("a refused put wrote a row: %s", s)
+	}
+	max := vrRow(t, 99, "wmax", 1000)
+	max.Request, max.Digest = bytes.Repeat([]byte("y"), approval.MaxVerbRequestBytes), sha256Hex(bytes.Repeat([]byte("y"), approval.MaxVerbRequestBytes))
+	vrPut(t, c, max, 1000) // control: exactly 128 KiB, the approval law's bound
+	old, kept := vrRow(t, 40, "wold", 1000-87000), vrRow(t, 41, "wkept", 1000-86999)
+	vrPut(t, c, old, 1000-87000)  // expires 1000-86400: expired a day ago at 1000
+	vrPut(t, c, kept, 1000-86999) // expires 1000-86399
+	for i := 1; i <= 13; i++ {
+		vrPut(t, c, vrRow(t, i, fmt.Sprintf("w%d", i), 1000), 1000)
+	}
+	if err := c.PutVerbRequest(testCtx(t), vrRow(t, 14, "w14", 1000), 1000); err != nil { // the 15th live (with wmax), old pruned
+		t.Fatal(err)
+	}
+	if _, ok, _ := c.VerbRequest(testCtx(t), old.RequestID); ok {
+		t.Fatal("a request expired a day ago was not pruned")
+	}
+	if _, ok, _ := c.VerbRequest(testCtx(t), kept.RequestID); !ok {
+		t.Fatal("a request expired less than a day ago was pruned")
+	}
+	vrPut(t, c, vrRow(t, 15, "w15", 1000), 1000) // 16 live
+	before := vrState(t, c)
+	if err := c.PutVerbRequest(testCtx(t), vrRow(t, 16, "w16", 1000), 1000); !errors.Is(err, ErrTooManyPendingVerbRequests) {
+		t.Fatalf("the 17th live request: %v", err)
+	}
+	if vrState(t, c) != before {
+		t.Fatal("the refused 17th request wrote")
+	}
+	vrPut(t, c, vrRow(t, 17, "w1", 1000), 1000) // w1's request is superseded first: still 16 live
+	if live, err := c.PendingVerbRequests(testCtx(t), 1000, 100); err != nil || len(live) != MaxPendingVerbRequests {
+		t.Fatalf("live after a supersede = %d, %v", len(live), err)
+	}
+}
+
+// CS5: the table's CHECKs, by raw INSERTs (backstops behind the Go grammar).
+func TestVerbRequestsChecks(t *testing.T) {
+	c := openTestCredentials(t)
+	id := 0
+	ins := func(reqID string, request []byte, consumed, outcome, approvedBy, artifact, decision any) error {
+		id++
+		if reqID == "" {
+			reqID = fmt.Sprintf("vr-%016x", id)
+		}
+		_, err := c.db.Exec(`INSERT INTO verb_requests (request_id, request, digest, challenge, episode_id, created_at, expires_at, consumed_at, outcome, approved_by, artifact, decision_ref) VALUES (?, ?, ?, ?, 'w1', 1000, 1600, ?, ?, ?, ?, ?)`,
+			reqID, request, strings.Repeat("a", 64), strings.Repeat("b", 64), consumed, outcome, approvedBy, artifact, decision)
+		return err
+	}
+	r := []byte("r")
+	for _, tc := range []struct {
+		name                                     string
+		reqID                                    string
+		request                                  []byte
+		consumed, outcome, by, artifact, decided any
+		ok                                       bool
+	}{
+		{"pending", "", r, nil, nil, nil, nil, nil, true},
+		{"superseded", "", r, 1100, "superseded", nil, nil, nil, true},
+		{"approved", "", r, 1100, "approved", "ssh:m", []byte("sig"), vrDecision, true},
+		{"id-mint", "mr-0000000000000001", r, nil, nil, nil, nil, nil, false},
+		{"id-length", "vr-01", r, nil, nil, nil, nil, nil, false},
+		{"request-empty", "", []byte{}, nil, nil, nil, nil, nil, false},
+		{"request-128KiB+1", "", make([]byte, 131073), nil, nil, nil, nil, nil, false},
+		{"outcome-unknown", "", r, 1100, "denied", nil, nil, nil, false},
+		{"outcome-unconsumed", "", r, nil, "superseded", nil, nil, nil, false},
+		{"consumed-no-outcome", "", r, 1100, nil, nil, nil, nil, false},
+		{"approved-no-by", "", r, 1100, "approved", nil, []byte("sig"), vrDecision, false},
+		{"approved-no-artifact", "", r, 1100, "approved", "ssh:m", nil, vrDecision, false},
+		{"approved-no-decision", "", r, 1100, "approved", "ssh:m", []byte("sig"), nil, false},
+		// Design §6 CS5 says a decision_ref on a superseded row is refused; the design's DDL (kept
+		// verbatim) does not refuse it: the approved-iff CHECK holds when all three are not set. Only
+		// ConsumeVerbRequest writes decision_ref, with outcome approved (plan §1.4 F3).
+		{"superseded-with-decision-accepted", "", r, 1100, "superseded", nil, nil, vrDecision, true},
+		{"superseded-with-approver", "", r, 1100, "superseded", "ssh:m", []byte("sig"), vrDecision, false},
+		{"by-empty", "", r, 1100, "approved", "", []byte("sig"), vrDecision, false},
+		{"artifact-16385", "", r, 1100, "approved", "ssh:m", make([]byte, 16385), vrDecision, false},
+		{"decision-sha512", "", r, 1100, "approved", "ssh:m", []byte("sig"), "sha512:" + strings.Repeat("d4", 32), false},
+		{"decision-70", "", r, 1100, "approved", "ssh:m", []byte("sig"), vrDecision[:70], false},
+	} {
+		if err := ins(tc.reqID, tc.request, tc.consumed, tc.outcome, tc.by, tc.artifact, tc.decided); (err == nil) != tc.ok {
+			t.Errorf("%s: err %v, want ok %v", tc.name, err, tc.ok)
+		}
+	}
+}
+
+// CS6: ConsumeVerbRequest's refusals each write nothing (rows, approvers' last use, counters); a
+// confirmation sets the five columns, decision_ref included, marks the approver used and moves a
+// passkey's counter forward, in one transaction.
+func TestConsumeVerbRequest(t *testing.T) {
+	sshUse := kdUse
+	keyUse := ApproverUse{Name: "x", Format: "webauthn", PublicKey: []byte("cose-x"), Passkey: &PasskeyUse{CredentialID: credX, SignCount: 6}}
+	type call struct {
+		id, digest, by string
+		use            ApproverUse
+		artifact       []byte
+		decision       string
+		now            int64
+	}
+	setup := func(t *testing.T) (*Credentials, VerbRequestRow, call) {
+		c := vrStore(t)
+		r := vrRow(t, 1, "w1", 1000)
+		vrPut(t, c, r, 1000)
+		return c, r, call{r.RequestID, r.Digest, "ssh:mark-ssh", sshUse, []byte("sig"), vrDecision, 1100}
+	}
+	consume := func(c *Credentials, k call) error {
+		return c.ConsumeVerbRequest(testCtx(t), k.id, k.digest, k.use, k.by, k.artifact, k.decision, k.now)
+	}
+	for _, tc := range []struct {
+		name string
+		prep func(*testing.T, *Credentials, VerbRequestRow)
+		edit func(*call)
+		want error
+	}{
+		{"unknown", nil, func(k *call) { k.id = "vr-00000000000000ff" }, ErrVerbRequestUnknown},
+		{"consumed", func(t *testing.T, c *Credentials, r VerbRequestRow) {
+			if err := c.ConsumeVerbRequest(testCtx(t), r.RequestID, r.Digest, sshUse, "ssh:mark-ssh", []byte("sig"), vrDecision, 1050); err != nil {
+				t.Fatal(err)
+			}
+		}, nil, ErrVerbRequestConsumed},
+		{"superseded", func(t *testing.T, c *Credentials, _ VerbRequestRow) { vrPut(t, c, vrRow(t, 2, "w1", 1050), 1050) }, nil, ErrVerbRequestSuperseded},
+		{"clock-behind", nil, func(k *call) { k.now = 999 }, ErrVerbRequestClockBehind},
+		{"expired", nil, func(k *call) { k.now = 1600 }, ErrVerbRequestExpired},
+		{"digest-other", nil, func(k *call) { k.digest = strings.Repeat("e", 64) }, ErrVerbRequestDigestChanged},
+		{"bytes-edited", func(t *testing.T, c *Credentials, r VerbRequestRow) {
+			if _, err := c.db.Exec(`UPDATE verb_requests SET request = ? WHERE request_id = ?`, append(r.Request, ' '), r.RequestID); err != nil {
+				t.Fatal(err)
+			}
+		}, nil, ErrVerbRequestDigestChanged},
+		{"approver-revoked", func(t *testing.T, c *Credentials, _ VerbRequestRow) {
+			if err := c.RevokeApprover(testCtx(t), "mark-ssh", 1050); err != nil {
+				t.Fatal(err)
+			}
+		}, nil, ErrApproverRevoked},
+		{"approver-key-changed", nil, func(k *call) { k.use.PublicKey = []byte("k2") }, ErrApproverKeyChanged},
+		{"counter-not-advancing", nil, func(k *call) {
+			k.use, k.by = keyUse, "webauthn:x"
+			k.use.Passkey = &PasskeyUse{CredentialID: credX, SignCount: 5}
+		}, ErrSignCountRegression},
+		{"by-not-the-approver", nil, func(k *call) { k.by = "ssh:other" }, ErrBadVerbRequest},
+		{"webauthn-without-passkey", nil, func(k *call) { k.use, k.by = keyUse, "webauthn:x"; k.use.Passkey = nil }, ErrBadVerbRequest},
+		{"artifact-empty", nil, func(k *call) { k.artifact = nil }, ErrBadVerbRequest},
+		{"artifact-16385", nil, func(k *call) { k.artifact = make([]byte, 16385) }, ErrBadVerbRequest},
+		{"decision-upper", nil, func(k *call) { k.decision = strings.ToUpper(vrDecision) }, ErrBadVerbRequest},
+		{"id-bad", nil, func(k *call) { k.id = "mr-0000000000000001" }, ErrBadVerbRequest},
+		{"digest-bad", nil, func(k *call) { k.digest = "x" }, ErrBadVerbRequest},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			c, r, k := setup(t)
+			if tc.prep != nil {
+				tc.prep(t, c, r)
+			}
+			if tc.edit != nil {
+				tc.edit(&k)
+			}
+			before := vrState(t, c)
+			if err := consume(c, k); !errors.Is(err, tc.want) || tc.name == "consumed" && !strings.Contains(err.Error(), "by ssh:mark-ssh") {
+				t.Fatalf("err %v, want %v", err, tc.want)
+			}
+			if after := vrState(t, c); after != before {
+				t.Fatalf("a refusal wrote:\n%s->\n%s", before, after)
+			}
+		})
+	}
+	for _, tc := range []struct {
+		name, by   string
+		use        ApproverUse
+		count, now int64
+	}{{"ssh", "ssh:mark-ssh", sshUse, 5, 1000}, {"webauthn", "webauthn:x", keyUse, 6, 1599}} { // at creation; a second before expiry
+		t.Run("confirm-"+tc.name, func(t *testing.T) {
+			c, r, k := setup(t)
+			k.use, k.by, k.now = tc.use, tc.by, tc.now
+			if err := consume(c, k); err != nil {
+				t.Fatal(err)
+			}
+			got := vrGet(t, c, r.RequestID)
+			if got.ConsumedAt != tc.now || got.Outcome != VerbOutcomeApproved || got.ApprovedBy != tc.by || !bytes.Equal(got.Artifact, []byte("sig")) || got.DecisionRef != vrDecision {
+				t.Fatalf("confirmed row = %+v", got)
+			}
+			if p, _, _ := c.Passkey(testCtx(t), "x"); p.SignCount != tc.count {
+				t.Fatalf("counter %d, want %d", p.SignCount, tc.count)
+			}
+			if a, _, _ := c.Approver(testCtx(t), tc.use.Name); a.LastUsedAt != tc.now {
+				t.Fatalf("%s last used %d", tc.use.Name, a.LastUsedAt)
+			}
+			if live, _ := c.PendingVerbRequests(testCtx(t), tc.now, 0); len(live) != 0 {
+				t.Fatalf("still pending: %+v", live)
+			}
+			vrPut(t, c, vrRow(t, 2, "w1", tc.now), tc.now) // a later finish of w1 supersedes only a pending request
+			if again := vrGet(t, c, r.RequestID); !reflect.DeepEqual(again, got) {
+				t.Fatalf("the confirmed row changed: %+v", again)
+			}
+		})
+	}
+}
````

### Patch M8 — `FinishConfirmRefusal` F1–F5 (design §2.3, §6 VT2)

````diff M8
diff --git a/host/broker/finish_confirm.go b/host/broker/finish_confirm.go
new file mode 100644
index 00000000..6ba7b397
--- /dev/null
+++ b/host/broker/finish_confirm.go
@@ -0,0 +1,41 @@
+package broker
+
+// Row 243 V3 sprint 3 M8 (design_docs/planned/w-world-address-bar-verbs-v3.md §2.3, §6 VT2): the
+// daemon's step C10, pure. A stored verb request is confirmable only if its scope is a finish
+// approval that outlives the request by at most an hour from now, its title and body are the ones
+// the scope approves, and its prompt is the terminal's rendering of them (D-WORLD-118 = A: no
+// draft-only rule).
+
+import (
+	"crypto/sha256"
+	"encoding/hex"
+)
+
+// FinishConfirmFacts are the fields of a decoded, validated world/verb-request/v1 that C10 judges,
+// with the request's expiry and the daemon's clock.
+type FinishConfirmFacts struct {
+	Scope, Title, Body, Prompt, Committer string
+	Protected                             []string
+	RequestExpires, Now                   int64
+}
+
+// FinishConfirmRefusal is why f is not confirmable (F1-F5), or "" when it is.
+func FinishConfirmRefusal(f FinishConfirmFacts) string {
+	scope, err := ParseFinishApprovalScope(f.Scope)
+	if err != nil { // F1
+		return "the request's scope is not a finish approval"
+	}
+	if !(f.RequestExpires < scope.ExpiresAt && scope.ExpiresAt <= f.Now+3600) { // F2
+		return "the request outlives its approval scope, or the scope is over an hour ahead"
+	}
+	if title := sha256.Sum256([]byte(f.Title)); hex.EncodeToString(title[:]) != scope.TitleSHA256 { // F3
+		return "the title is not the one the scope approves"
+	}
+	if body := sha256.Sum256([]byte(f.Body)); hex.EncodeToString(body[:]) != scope.BodySHA256 { // F4
+		return "the body is not the one the scope approves"
+	}
+	if f.Prompt != FinishPromptText(scope, f.Title, f.Body, f.Protected, f.Committer) { // F5
+		return "the prompt is not the terminal's rendering of what is approved"
+	}
+	return ""
+}
diff --git a/host/broker/finish_confirm_test.go b/host/broker/finish_confirm_test.go
new file mode 100644
index 00000000..3f0688f5
--- /dev/null
+++ b/host/broker/finish_confirm_test.go
@@ -0,0 +1,67 @@
+package broker
+
+// Row 243 V3 sprint 3 M8 (design §2.3, §6 VT2): FinishConfirmRefusal's F1-F5, each from a fixture
+// whose scope is FormatFinishApprovalScope's and whose prompt is FinishPromptText's (never typed).
+
+import (
+	"crypto/sha256"
+	"encoding/hex"
+	"testing"
+)
+
+func vtHex(s string) string {
+	h := sha256.Sum256([]byte(s))
+	return hex.EncodeToString(h[:])
+}
+
+// vtFacts is an admissible request at now 1000 expiring 1600, its scope expiring 1900.
+func vtFacts() (FinishConfirmFacts, FinishApprovalScope) {
+	title, body, protected, committer := "Fix the bar", "## Summary\nline 2\n", []string{".github/workflows/ci.yml", "go.mod"}, ""
+	s := FinishApprovalScope{Repository: "sunholo-data/ailang-world", Remote: "origin", Head: "world/w1", SHA: "0123456789abcdef0123456789abcdef01234567", Base: "dev",
+		TitleSHA256: vtHex(title), BodySHA256: vtHex(body), Draft: true, ExpiresAt: 1900}
+	return FinishConfirmFacts{Scope: FormatFinishApprovalScope(s), Title: title, Body: body, Prompt: FinishPromptText(s, title, body, protected, committer),
+		Committer: committer, Protected: protected, RequestExpires: 1600, Now: 1000}, s
+}
+
+func TestFinishConfirmRefusal(t *testing.T) {
+	const f1, f2, f3, f4, f5 = "the request's scope is not a finish approval", "the request outlives its approval scope, or the scope is over an hour ahead",
+		"the title is not the one the scope approves", "the body is not the one the scope approves", "the prompt is not the terminal's rendering of what is approved"
+	rescope := func(edit func(*FinishApprovalScope)) func(*FinishConfirmFacts) {
+		return func(f *FinishConfirmFacts) {
+			_, s := vtFacts()
+			edit(&s)
+			f.Scope = FormatFinishApprovalScope(s)
+			f.Prompt = FinishPromptText(s, f.Title, f.Body, f.Protected, f.Committer)
+		}
+	}
+	for _, tc := range []struct {
+		name string
+		edit func(*FinishConfirmFacts)
+		want string
+	}{
+		{"admissible", func(*FinishConfirmFacts) {}, ""},
+		{"ready", rescope(func(s *FinishApprovalScope) { s.Draft = false }), ""}, // D-WORLD-118 = A: no draft-only rule
+		{"F1-not-a-scope", func(f *FinishConfirmFacts) { f.Scope = "github-pr:o/r#x" }, f1},
+		{"F1-trailing-byte", func(f *FinishConfirmFacts) { f.Scope += " " }, f1},
+		{"F2-scope-equals-request", rescope(func(s *FinishApprovalScope) { s.ExpiresAt = 1600 }), f2},
+		{"F2-scope-before-request", rescope(func(s *FinishApprovalScope) { s.ExpiresAt = 1599 }), f2},
+		{"F2-scope-an-hour-ahead", rescope(func(s *FinishApprovalScope) { s.ExpiresAt = 4600 }), ""},
+		{"F2-scope-over-an-hour-ahead", rescope(func(s *FinishApprovalScope) { s.ExpiresAt = 4601 }), f2},
+		{"F3-title", func(f *FinishConfirmFacts) { f.Title = "Fix the bar." }, f3},
+		{"F4-body", func(f *FinishConfirmFacts) { f.Body += "x" }, f4},
+		{"F5-prompt-edited", func(f *FinishConfirmFacts) { f.Prompt += "\n" }, f5},
+		{"F5-protected-not-shown", func(f *FinishConfirmFacts) { f.Protected = f.Protected[:1] }, f5},
+		{"F5-committer-not-shown", func(f *FinishConfirmFacts) { f.Committer = "Mark <m@x>" }, f5},
+		{"F5-prompt-of-another-scope", func(f *FinishConfirmFacts) {
+			_, s := vtFacts()
+			s.Base = "main"
+			f.Prompt = FinishPromptText(s, f.Title, f.Body, f.Protected, f.Committer)
+		}, f5},
+	} {
+		f, _ := vtFacts()
+		tc.edit(&f)
+		if got := FinishConfirmRefusal(f); got != tc.want {
+			t.Errorf("%s: %q, want %q", tc.name, got, tc.want)
+		}
+	}
+}
````

### Patch M9a — the three GET routes, CR1/CR3/CR4, the unrouted test deleted (design §2.5)

````diff M9a
diff --git a/host/daemon/daemon.go b/host/daemon/daemon.go
index 8a2c0d44..6d3e02c9 100644
--- a/host/daemon/daemon.go
+++ b/host/daemon/daemon.go
@@ -1101,7 +1101,9 @@ func releaseFromVersion(version string) string {
 // /surface/passkey.js (the passkey ceremony page, row 201 A2: static, it reads nothing), and the
 // approval routes of row 201 A3 (requests.json, requests/{id}, the approve POST, the card page
 // GET /surface/approve/{id} and its approve.js), each behind approvalRoute's Host check, and
-// row 201 A4's decision routes GET /surface/decisions/{ref} and POST /surface/decide/{ref}.
+// row 201 A4's decision routes GET /surface/decisions/{ref} and POST /surface/decide/{ref}, and
+// row 243 V3's confirmation card GET /surface/confirm/{id}, its confirm.js and its data GET
+// /surface/confirmations/{id} (read-only; no confirmation POST before sprint 4).
 // The registry pattern deliberately uses a multi-segment wildcard: registry semantic IDs
 // such as "world/epoch-registry/v1" contain slashes.
 func (d *Daemon) Handler() http.Handler {
@@ -1145,6 +1147,9 @@ func (d *Daemon) Handler() http.Handler {
 	mux.HandleFunc("POST /surface/decide/{ref}", d.approvalRoute(d.handleSurfaceDecide))
 	mux.HandleFunc("GET /surface/decide/{ref}", d.approvalRoute(d.handleSurfaceDecidePage))
 	mux.HandleFunc("GET /surface/decide.js", d.approvalRoute(d.handleSurfaceDecideScript))
+	mux.HandleFunc("GET /surface/confirm/{id}", d.approvalRoute(d.handleSurfaceConfirmPage))
+	mux.HandleFunc("GET /surface/confirm.js", d.approvalRoute(d.handleSurfaceConfirmScript))
+	mux.HandleFunc("GET /surface/confirmations/{id}", d.approvalRoute(d.handleSurfaceConfirmation))
 	mux.HandleFunc("GET /surface/verb.json", d.approvalRoute(d.handleSurfaceVerb)) // row 243 V1: the bar's verbs, read-only
 	// The two A2A projection routes (w-a2a-session-projection P6.B-A2A-CARD)
 	// are ADDITIVE: the frozen /v1/ table above is untouched, and the routes
diff --git a/host/daemon/surface_confirm_read.go b/host/daemon/surface_confirm_read.go
new file mode 100644
index 00000000..d22ccfe4
--- /dev/null
+++ b/host/daemon/surface_confirm_read.go
@@ -0,0 +1,79 @@
+package daemon
+
+// Row 243 V3 sprint 3 M9 (design_docs/planned/w-world-address-bar-verbs-v3.md §2.5): the
+// confirmation card's read routes, each behind approvalRoute's Host check, served as the decision
+// card's are; the data is surfacepage.ConfirmationCard of the row. Nothing here writes (POST: sprint 4).
+
+import (
+	"net/http"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+const confirmIDReason = "a verb request id is vr- and 16 lowercase hex"
+
+// handleSurfaceConfirmPage is GET /surface/confirm/{id}: the same bytes for every well-formed id.
+func (d *Daemon) handleSurfaceConfirmPage(w http.ResponseWriter, r *http.Request, _ []string) {
+	if !store.VerbRequestIDPattern.MatchString(r.PathValue("id")) {
+		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", confirmIDReason)
+		return
+	}
+	setSurfaceHeaders(w)
+	_, _ = w.Write(surfacepage.ConfirmPage)
+}
+
+// handleSurfaceConfirmScript is the card's script, under the same headers.
+func (d *Daemon) handleSurfaceConfirmScript(w http.ResponseWriter, _ *http.Request, _ []string) {
+	setSurfaceHeaders(w)
+	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
+	_, _ = w.Write(surfacepage.ConfirmScript)
+}
+
+// handleSurfaceConfirmation is GET /surface/confirmations/{id}: 400 for a bad id, 503 naming the
+// serve restart below schema 9, 404 for an unknown one, else the card's data with the active
+// passkeys by origin.
+func (d *Daemon) handleSurfaceConfirmation(w http.ResponseWriter, r *http.Request, _ []string) {
+	id := r.PathValue("id")
+	if !store.VerbRequestIDPattern.MatchString(id) {
+		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", confirmIDReason)
+		return
+	}
+	ctx, cancel := d.readCtx(r)
+	defer cancel()
+	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
+	if err != nil {
+		d.approvalInternal(w, r, err)
+		return
+	}
+	defer func() { _ = ro.Close() }()
+	row, ok, err := ro.VerbRequest(ctx, id)
+	switch {
+	case isOlderStore(err):
+		writeApprovalRefusal(w, http.StatusServiceUnavailable, "Restart",
+			"the credentials store predates verb requests: restart `ailang-worldd serve` with this build once (it upgrades the store)")
+		return
+	case err != nil:
+		d.approvalInternal(w, r, err)
+		return
+	case !ok:
+		writeApprovalRefusal(w, http.StatusNotFound, "NotFound", "no verb request "+id)
+		return
+	}
+	active, err := ro.ActivePasskeys(ctx)
+	if err != nil {
+		d.approvalInternal(w, r, err)
+		return
+	}
+	var passkeys []surfacepage.CardPasskey
+	for _, e := range active {
+		passkeys = append(passkeys, surfacepage.CardPasskey{Origin: e.Passkey.Origin, RPID: e.Passkey.RPID, CredentialID: e.Passkey.CredentialID})
+	}
+	writeApprovalJSON(w, http.StatusOK, surfacepage.ConfirmationCard(confirmationInput(row, d.approveNow()), passkeys))
+}
+
+// confirmationInput is row as the card reads it at now.
+func confirmationInput(row store.VerbRequestRow, now int64) surfacepage.ConfirmationInput {
+	return surfacepage.ConfirmationInput{ID: row.RequestID, Request: row.Request, Digest: row.Digest, Challenge: row.Challenge, Created: row.CreatedAt,
+		Expires: row.ExpiresAt, Now: now, Outcome: row.Outcome, ApprovedBy: row.ApprovedBy, Decision: row.DecisionRef}
+}
diff --git a/host/daemon/surface_confirm_read_test.go b/host/daemon/surface_confirm_read_test.go
new file mode 100644
index 00000000..70588cf6
--- /dev/null
+++ b/host/daemon/surface_confirm_read_test.go
@@ -0,0 +1,203 @@
+package daemon
+
+// Row 243 V3 sprint 3 M9 (design §2.5, §6 CR1, CR2): the confirmation card's read routes. Every
+// request is approval.CanonicalVerb of a request whose prompt is broker.FinishPromptText and whose
+// scope is broker.FormatFinishApprovalScope; every expected card is surfacepage.ConfirmationCard of
+// the same row (sprint 2's golden is the shape contract). Stores and keys are tests' (D-WORLD-109).
+
+import (
+	"bytes"
+	"context"
+	"crypto/sha256"
+	"encoding/base64"
+	"encoding/hex"
+	"encoding/json"
+	"os"
+	"path/filepath"
+	"reflect"
+	"regexp"
+	"sort"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// putVerb stores a pending verb request id for ep, created at created (expiring 600 s later), for a
+// draft (or ready) pull request of head world/<ep>.
+func (a approvalRig) putVerb(id, ep string, created int64, draft bool) store.VerbRequestRow {
+	a.t.Helper()
+	return a.putVerbBytes(id, ep, created, canonVerb(a.t, verbFixture(created, ep, draft)))
+}
+
+// verbFixture is a valid request: its prompt FinishPromptText's, its scope FormatFinishApprovalScope's.
+func verbFixture(created int64, ep string, draft bool) approval.VerbRequest {
+	title, body := "Fix the <bar> & "+ep, "## Summary\n"+ep+"\n"
+	sum := func(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }
+	s := broker.FinishApprovalScope{Repository: "sunholo-data/ailang-world", Remote: "origin", Head: "world/" + ep, SHA: strings.Repeat("0a", 20), Base: "dev",
+		TitleSHA256: sum(title), BodySHA256: sum(body), Draft: draft, ExpiresAt: created + 900}
+	r := approval.VerbRequest{Schema: approval.VerbRequestSchema, Verb: approval.VerbFinish, Argv: []string{"ailang-worldd", "episode", "finish", ep, "--preset", "p", "--confirm-on-surface"},
+		Episode: ep, Prompt: broker.FinishPromptText(s, title, body, []string{"go.mod"}, ""), Scope: broker.FormatFinishApprovalScope(s), Title: title, Body: body, Protected: []string{"go.mod"}}
+	return r
+}
+
+func canonVerb(t *testing.T, r approval.VerbRequest) []byte {
+	b, err := approval.CanonicalVerb(r)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return b
+}
+
+func (a approvalRig) putVerbBytes(id, ep string, created int64, b []byte) store.VerbRequestRow {
+	a.t.Helper()
+	d, c := sha256.Sum256(b), sha256.Sum256([]byte(id))
+	row := store.VerbRequestRow{RequestID: id, Request: b, Digest: hex.EncodeToString(d[:]), Challenge: hex.EncodeToString(c[:]), EpisodeID: ep, CreatedAt: created, ExpiresAt: created + 600}
+	a.write(func(ctx context.Context, c *store.Credentials) error { return c.PutVerbRequest(ctx, row, created) })
+	return row
+}
+
+// readVerb is the stored row id.
+func (a approvalRig) readVerb(id string) (row store.VerbRequestRow) {
+	a.t.Helper()
+	a.write(func(ctx context.Context, c *store.Credentials) error { row, _, _ = c.VerbRequest(ctx, id); return nil })
+	return row
+}
+
+// The route table: exactly three GETs under /surface/confirm, each behind approvalRoute, and no
+// POST (sprint 4's); other methods 405 with no CORS header; a foreign Host 421; other names 404.
+func TestConfirmRoutesRegistered(t *testing.T) {
+	src, err := os.ReadFile("daemon.go")
+	if err != nil {
+		t.Fatal(err)
+	}
+	body := strings.SplitN(strings.SplitN(string(src), "func (d *Daemon) Handler() http.Handler {", 2)[1], "\n}\n", 2)[0]
+	var got []string
+	for _, m := range regexp.MustCompile(`mux\.HandleFunc\("([A-Z]+ /surface/confirm[^"]*)", (d\.approvalRoute\()?`).FindAllStringSubmatch(body, -1) {
+		if m[2] == "" {
+			t.Errorf("%s is not behind approvalRoute", m[1])
+		}
+		got = append(got, m[1])
+	}
+	sort.Strings(got)
+	if want := []string{"GET /surface/confirm.js", "GET /surface/confirm/{id}", "GET /surface/confirmations/{id}"}; !reflect.DeepEqual(got, want) {
+		t.Fatalf("confirmation routes %v, want %v", got, want)
+	}
+	a := newApprovalRig(t)
+	id := "vr-0123456789abcdef"
+	for _, target := range []string{"/surface/confirm/" + id, "/surface/confirm.js", "/surface/confirmations/" + id} {
+		requireRefusal(t, "foreign "+target, a.serve("GET", target, "rebind.example:7644", nil), 421, "MisdirectedRequest")
+		for _, method := range []string{"POST", "PUT", "DELETE", "OPTIONS"} {
+			rec := a.serve(method, target, apLoopback, nil, "Origin", "http://localhost:7644", "Access-Control-Request-Method", "POST")
+			if rec.Code != 405 || rec.Header().Get("Access-Control-Allow-Origin") != "" {
+				t.Errorf("%s %s = %d, want 405 with no CORS header", method, target, rec.Code)
+			}
+		}
+	}
+	for _, target := range []string{"/surface/confirm", "/surface/confirm/", "/surface/confirm/" + id + "/x", "/surface/confirmations", "/surface/confirmations/" + id + "/x", "/surface/confirm.jsx"} {
+		if rec := a.serve("GET", target, apLoopback, nil); rec.Code != 404 {
+			t.Errorf("GET %s = %d, want 404", target, rec.Code)
+		}
+	}
+}
+
+// PD2 for the confirmation card: the page and the script are host/surfacepage/confirm's bytes,
+// under surfaceCSP and the surface headers; a bad id is 400.
+func TestConfirmPageAssetsServed(t *testing.T) {
+	a := newApprovalRig(t)
+	for _, tc := range []struct{ target, file, contentType string }{
+		{"/surface/confirm/vr-0123456789abcdef", "index.html", "text/html; charset=utf-8"},
+		{"/surface/confirm/vr-0000000000000000", "index.html", "text/html; charset=utf-8"},
+		{"/surface/confirm.js", "confirm.js", "text/javascript; charset=utf-8"},
+	} {
+		want, err := os.ReadFile(filepath.Join("..", "surfacepage", "confirm", tc.file))
+		if err != nil {
+			t.Fatal(err)
+		}
+		rec := a.serve("GET", tc.target, apLoopback, nil)
+		if rec.Code != 200 || len(want) == 0 || !bytes.Equal(rec.Body.Bytes(), want) {
+			t.Errorf("GET %s: %d, or not the bytes of host/surfacepage/confirm/%s", tc.target, rec.Code, tc.file)
+		}
+		for k, v := range map[string]string{"Content-Security-Policy": surfaceCSP, "Content-Type": tc.contentType,
+			"Cache-Control": "no-store", "X-Content-Type-Options": "nosniff", "Referrer-Policy": "no-referrer"} {
+			if got := rec.Header().Get(k); got != v {
+				t.Errorf("GET %s: %s = %q, want %q", tc.target, k, got, v)
+			}
+		}
+		if rec.Header().Get("Set-Cookie") != "" {
+			t.Errorf("GET %s sets a cookie", tc.target)
+		}
+		requireRefusal(t, "foreign host", a.serve("GET", tc.target, "evil.example", nil), 421, "MisdirectedRequest")
+	}
+	for _, bad := range []string{"vr-0123456789ABCDEF", "vr-0123", "mr-0123456789abcdef", "sha256:" + strings.Repeat("ab", 32)} {
+		requireRefusal(t, bad, a.serve("GET", "/surface/confirm/"+bad, apLoopback, nil), 400, "BadRequest")
+	}
+}
+
+// CR1: GET /surface/confirmations/{id} is surfacepage.ConfirmationCard of the stored row, byte for
+// byte, for a pending, an approved and a superseded row: the bytes and the raw challenge, never the
+// message hash, never the artifact. 400 bad id, 404 unknown, 503 below schema 9, 500 no store.
+func TestConfirmationCardServed(t *testing.T) {
+	a := newApprovalRig(t)
+	pk := approvaltest.NewPasskey(approval.AlgEdDSA, 61)
+	a.enrolPasskey("a-mac", pk, "http://localhost:7644")
+	row := a.putVerb("vr-00000000000000a1", "w243", apT0-30, true)
+	other := a.putVerb("vr-00000000000000b1", "w244", apT0-20, false)
+	artifact := []byte("ARTIFACT-NEVER-SERVED")
+	decision := "sha256:" + strings.Repeat("d4", 32)
+	passkeys := []surfacepage.CardPasskey{{Origin: "http://localhost:7644", RPID: "localhost", CredentialID: pk.CredentialID()}}
+	check := func(what string, r store.VerbRequestRow, state string) {
+		t.Helper()
+		rec := a.serve("GET", "/surface/confirmations/"+r.RequestID, apLoopback, nil)
+		want, err := json.Marshal(surfacepage.ConfirmationCard(surfacepage.ConfirmationInput{ID: r.RequestID, Request: r.Request, Digest: r.Digest, Challenge: r.Challenge,
+			Created: r.CreatedAt, Expires: r.ExpiresAt, Now: apT0, Outcome: r.Outcome, ApprovedBy: r.ApprovedBy, Decision: r.DecisionRef}, passkeys))
+		if err := approvalHeaders(rec); err != nil || rec.Code != 200 {
+			t.Fatalf("%s: %d %v", what, rec.Code, err)
+		}
+		if err != nil || !bytes.Equal(bytes.TrimSuffix(rec.Body.Bytes(), []byte("\n")), want) {
+			t.Fatalf("%s:\n got %s\nwant %s", what, rec.Body.Bytes(), want)
+		}
+		var c surfacepage.Confirmation
+		if err := json.Unmarshal(rec.Body.Bytes(), &c); err != nil || c.State != state || c.Request != base64.RawURLEncoding.EncodeToString(r.Request) || c.Challenge != r.Challenge {
+			t.Fatalf("%s: %+v %v, want state %s", what, c, err, state)
+		}
+		d, _ := hex.DecodeString(r.Digest)
+		n, _ := hex.DecodeString(r.Challenge)
+		h, err := approval.VerbMessageHash([32]byte(d), [32]byte(n), r.ExpiresAt)
+		if err != nil {
+			t.Fatal(err)
+		}
+		for _, leak := range []string{hex.EncodeToString(h[:]), base64.RawURLEncoding.EncodeToString(h[:]), base64.StdEncoding.EncodeToString(h[:]), string(artifact), base64.StdEncoding.EncodeToString(artifact), `"artifact"`, `"hash"`} {
+			if strings.Contains(rec.Body.String(), leak) {
+				t.Fatalf("%s: the card carries %q", what, leak)
+			}
+		}
+	}
+	check("pending", row, surfacepage.ConfirmationPending)
+	a.write(func(ctx context.Context, c *store.Credentials) error {
+		return c.ConsumeVerbRequest(ctx, other.RequestID, other.Digest, store.ApproverUse{Name: "a-mac", Format: approval.FormatWebAuthn, PublicKey: pk.Key().Canonical(),
+			Passkey: &store.PasskeyUse{CredentialID: pk.CredentialID(), SignCount: 2}}, "webauthn:a-mac", artifact, decision, apT0-10)
+	})
+	if got := a.readVerb(other.RequestID); got.Outcome != store.VerbOutcomeApproved || got.DecisionRef != decision {
+		t.Fatalf("approved row = %+v", got)
+	}
+	check("approved", a.readVerb(other.RequestID), surfacepage.ConfirmationApproved)
+	a.putVerb("vr-00000000000000a2", "w243", apT0-5, true) // supersedes row
+	check("superseded", a.readVerb(row.RequestID), surfacepage.ConfirmationSuperseded)
+	requireRefusal(t, "bad id", a.serve("GET", "/surface/confirmations/VR-00000000000000A1", apLoopback, nil), 400, "BadRequest")
+	requireRefusal(t, "mint id", a.serve("GET", "/surface/confirmations/mr-00000000000000a1", apLoopback, nil), 400, "BadRequest")
+	requireRefusal(t, "unknown", a.serve("GET", "/surface/confirmations/vr-00000000000000f1", apLoopback, nil), 404, "NotFound")
+	requireRefusal(t, "foreign host", a.serve("GET", "/surface/confirmations/"+row.RequestID, "evil.example", nil), 421, "MisdirectedRequest")
+	setVersion(t, a.d.credentialsPath, `DROP TABLE verb_requests`, `PRAGMA user_version = 8`)
+	if r := requireRefusal(t, "v8", a.serve("GET", "/surface/confirmations/"+row.RequestID, apLoopback, nil), 503, "Restart"); !strings.Contains(r.Reason, "restart `ailang-worldd serve`") {
+		t.Fatalf("v8 reason %q", r.Reason)
+	}
+	if err := os.Remove(a.d.credentialsPath); err != nil {
+		t.Fatal(err)
+	}
+	requireRefusal(t, "no store", a.serve("GET", "/surface/confirmations/"+row.RequestID, apLoopback, nil), 500, "Internal")
+}
diff --git a/host/daemon/surface_confirm_unrouted_test.go b/host/daemon/surface_confirm_unrouted_test.go
deleted file mode 100644
index 5d2d1952..00000000
--- a/host/daemon/surface_confirm_unrouted_test.go
+++ /dev/null
@@ -1,33 +0,0 @@
-package daemon
-
-// Row 243 V3 sprint 2 (design §2.5 [r1], §8): the confirmation card's page and its data exist only in
-// host/surfacepage and tests; no daemon route serves them before sprint 3, which replaces this test
-// with the routes' own (CR1).
-
-import (
-	"os"
-	"strings"
-	"testing"
-)
-
-// The Handler registers nothing under /surface/confirm, and the daemon answers 404 for the card, its
-// script and its data. Control: the decision card's page, registered beside them, answers 200.
-func TestConfirmRoutesNotRegistered(t *testing.T) {
-	src, err := os.ReadFile("daemon.go")
-	if err != nil {
-		t.Fatal(err)
-	}
-	body := strings.SplitN(strings.SplitN(string(src), "func (d *Daemon) Handler() http.Handler {", 2)[1], "\n}\n", 2)[0]
-	if strings.Contains(body, "/surface/confirm") || !strings.Contains(body, `"GET /surface/decide/{ref}"`) {
-		t.Fatal("Handler registers a /surface/confirm route (or lost the decide control)")
-	}
-	a := newApprovalRig(t)
-	for _, target := range []string{"/surface/confirm/vr-0123456789abcdef", "/surface/confirm.js", "/surface/confirmations/vr-0123456789abcdef"} {
-		if rec := a.serve("GET", target, apLoopback, nil); rec.Code != 404 {
-			t.Errorf("GET %s = %d, want 404 (no route before sprint 3)", target, rec.Code)
-		}
-	}
-	if rec := a.serve("GET", "/surface/decide/sha256:"+strings.Repeat("ab", 32), apLoopback, nil); rec.Code != 200 {
-		t.Fatalf("control: GET /surface/decide/… = %d", rec.Code)
-	}
-}
````

### Patch M9b — `requests.json`'s confirmations, CR2 (design §2.5)

````diff M9b
diff --git a/host/daemon/surface_approve.go b/host/daemon/surface_approve.go
index b515dcac..54e344bc 100644
--- a/host/daemon/surface_approve.go
+++ b/host/daemon/surface_approve.go
@@ -110,6 +110,15 @@ func (d *Daemon) handleSurfaceRequests(w http.ResponseWriter, r *http.Request, _
 	if out.Decisions == nil {
 		out.Decisions = []surfacepage.PendingDecision{}
 	}
+	if err := d.readConfirmations(ctx, now, &out); isOlderStore(err) { // row 243 V3 §2.5: the verb requests
+		out.Confirmations, out.ConfirmationsUnavailable = nil, "the credentials store predates verb requests: restart `ailang-worldd serve` with this build once"
+	} else if err != nil {
+		d.writeInternalErrorLog(r, err)
+		out.Confirmations, out.ConfirmationsUnavailable = nil, "the verb requests could not be read"
+	}
+	if out.Confirmations == nil {
+		out.Confirmations = []surfacepage.PendingConfirmation{}
+	}
 	writeApprovalJSON(w, http.StatusOK, out)
 }
 
diff --git a/host/daemon/surface_confirm_read.go b/host/daemon/surface_confirm_read.go
index d22ccfe4..1b9f7780 100644
--- a/host/daemon/surface_confirm_read.go
+++ b/host/daemon/surface_confirm_read.go
@@ -5,8 +5,13 @@ package daemon
 // card's are; the data is surfacepage.ConfirmationCard of the row. Nothing here writes (POST: sprint 4).
 
 import (
+	"context"
+	"crypto/sha256"
+	"encoding/hex"
 	"net/http"
 
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/broker"
 	"github.com/sunholo-data/ailang-world/host/store"
 	"github.com/sunholo-data/ailang-world/host/surfacepage"
 )
@@ -77,3 +82,46 @@ func confirmationInput(row store.VerbRequestRow, now int64) surfacepage.Confirma
 	return surfacepage.ConfirmationInput{ID: row.RequestID, Request: row.Request, Digest: row.Digest, Challenge: row.Challenge, Created: row.CreatedAt,
 		Expires: row.ExpiresAt, Now: now, Outcome: row.Outcome, ApprovedBy: row.ApprovedBy, Decision: row.DecisionRef}
 }
+
+// readConfirmations fills out.Confirmations with the live, pending verb requests on one read-only
+// handle: a summary decoded from the stored bytes, or Refused and no summary.
+func (d *Daemon) readConfirmations(ctx context.Context, now int64, out *surfacepage.RequestsJSON) error {
+	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
+	if err != nil {
+		return err
+	}
+	defer func() { _ = ro.Close() }()
+	rows, err := ro.PendingVerbRequests(ctx, now, store.MaxPendingVerbRequests)
+	if err != nil {
+		return err
+	}
+	for _, row := range rows {
+		out.Confirmations = append(out.Confirmations, pendingConfirmation(row))
+	}
+	return nil
+}
+
+// pendingConfirmation summarises row, or refuses it when its bytes do not hash to its digest, do
+// not decode as a valid world/verb-request/v1, or carry no finish approval scope.
+func pendingConfirmation(row store.VerbRequestRow) surfacepage.PendingConfirmation {
+	p := surfacepage.PendingConfirmation{ID: row.RequestID, Expires: row.ExpiresAt}
+	if sum := sha256.Sum256(row.Request); hex.EncodeToString(sum[:]) != row.Digest {
+		p.Refused = "the stored request does not match its digest"
+		return p
+	}
+	r, err := approval.DecodeVerb(row.Request)
+	if err == nil {
+		err = r.Validate()
+	}
+	if err != nil {
+		p.Refused = "the stored request is not a valid verb request: " + err.Error()
+		return p
+	}
+	s, err := broker.ParseFinishApprovalScope(r.Scope)
+	if err != nil {
+		p.Refused = "the stored request's scope is not a finish approval"
+		return p
+	}
+	p.Episode, p.Repository, p.Head, p.SHA, p.Draft = r.Episode, s.Repository, s.Head, s.SHA, &s.Draft
+	return p
+}
diff --git a/host/daemon/surface_confirm_read_test.go b/host/daemon/surface_confirm_read_test.go
index 70588cf6..ebf41c0d 100644
--- a/host/daemon/surface_confirm_read_test.go
+++ b/host/daemon/surface_confirm_read_test.go
@@ -9,6 +9,7 @@ import (
 	"bytes"
 	"context"
 	"crypto/sha256"
+	"database/sql"
 	"encoding/base64"
 	"encoding/hex"
 	"encoding/json"
@@ -201,3 +202,84 @@ func TestConfirmationCardServed(t *testing.T) {
 	}
 	requireRefusal(t, "no store", a.serve("GET", "/surface/confirmations/"+row.RequestID, apLoopback, nil), 500, "Internal")
 }
+
+// CR2: requests.json's confirmations are the live pending verb requests, newest first, each
+// summarised from its stored bytes; a row whose bytes do not hash to its digest, or do not decode,
+// is listed refused with no summary; approved, superseded and expired rows are not listed; below
+// schema 9 the list is empty with confirmationsUnavailable and the other lists intact.
+func TestRequestsJSONConfirmations(t *testing.T) {
+	a := newApprovalRig(t)
+	mint := a.putRequest("mr-00000000000000b1", "w201b", apT0-30)
+	a.putVerb("vr-00000000000000a0", "wold", apT0-700, true) // expired at apT0-100
+	sup := a.putVerb("vr-00000000000000a1", "w1", apT0-60, true)
+	ready := a.putVerb("vr-00000000000000a2", "w2", apT0-50, false)
+	done := a.putVerb("vr-00000000000000a3", "w3", apT0-45, true)
+	draft := a.putVerb("vr-00000000000000a4", "w1", apT0-40, true) // supersedes sup
+	undecodable := a.putVerbBytes("vr-00000000000000a5", "w5", apT0-30, []byte(`{"schema":"world/verb-request/v1"}`))
+	edited := a.putVerb("vr-00000000000000a6", "w6", apT0-20, true)
+	setVerbColumn(t, a.d.credentialsPath, edited.RequestID, "request", append(edited.Request, ' '))
+	start, noScope := verbFixture(apT0-15, "w7", true), verbFixture(apT0-12, "w8", true)
+	start.Verb, noScope.Scope = "episode start", "github-pr:o/r#x" // each decodes; Validate, then the scope parse, refuses
+	invalid, unscoped := a.putVerbBytes("vr-00000000000000a7", "w7", apT0-15, canonVerb(t, start)), a.putVerbBytes("vr-00000000000000a8", "w8", apT0-12, canonVerb(t, noScope))
+	a.write(func(ctx context.Context, c *store.Credentials) error {
+		if err := c.EnrolApprover(ctx, store.ApproverRow{Name: "m", Format: "ssh", PublicKey: []byte("k"), EnrolledAt: apT0 - 100}); err != nil {
+			return err
+		}
+		return c.ConsumeVerbRequest(ctx, done.RequestID, done.Digest, store.ApproverUse{Name: "m", Format: "ssh", PublicKey: []byte("k")}, "ssh:m", []byte("sig"), "sha256:"+strings.Repeat("d4", 32), apT0-10)
+	})
+	get := func() (surfacepage.RequestsJSON, string) {
+		rec := a.serve("GET", "/surface/requests.json", apLoopback, nil)
+		var out surfacepage.RequestsJSON
+		if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil || rec.Code != 200 {
+			t.Fatalf("requests.json = %d %v", rec.Code, err)
+		}
+		return out, rec.Body.String()
+	}
+	out, body := get()
+	yes, no := true, false
+	want := []surfacepage.PendingConfirmation{
+		{ID: unscoped.RequestID, Expires: apT0 + 588, Refused: "the stored request's scope is not a finish approval"},
+		{ID: invalid.RequestID, Expires: apT0 + 585, Refused: `the stored request is not a valid verb request: approval: bad verb request: verb is "episode start", want "episode finish"`},
+		{ID: edited.RequestID, Expires: apT0 + 580, Refused: "the stored request does not match its digest"},
+		{ID: undecodable.RequestID, Expires: apT0 + 570, Refused: out.Confirmations[3].Refused},
+		{ID: draft.RequestID, Episode: "w1", Repository: "sunholo-data/ailang-world", Head: "world/w1", SHA: strings.Repeat("0a", 20), Draft: &yes, Expires: apT0 + 560},
+		{ID: ready.RequestID, Episode: "w2", Repository: "sunholo-data/ailang-world", Head: "world/w2", SHA: strings.Repeat("0a", 20), Draft: &no, Expires: apT0 + 550},
+	}
+	if !reflect.DeepEqual(out.Confirmations, want) || out.ConfirmationsUnavailable != "" || !strings.HasPrefix(want[3].Refused, "the stored request is not a valid verb request: ") {
+		t.Fatalf("confirmations = %s", body)
+	}
+	if len(out.Pending) != 1 || out.Pending[0].ID != mint.RequestID || out.Unavailable != "" {
+		t.Fatalf("pending = %+v", out.Pending)
+	}
+	for _, leak := range []string{sup.RequestID, done.RequestID, draft.Challenge, draft.Digest, `"request"`, `"challenge"`, "sig"} {
+		if strings.Contains(body, leak) {
+			t.Errorf("requests.json carries %q", leak)
+		}
+	}
+	setVersion(t, a.d.credentialsPath, `DROP TABLE verb_requests`, `PRAGMA user_version = 8`)
+	out, body = get()
+	if len(out.Pending) != 1 || out.Unavailable != "" || !strings.Contains(body, `"confirmations":[],"confirmationsUnavailable":"the credentials store predates verb requests: restart`) {
+		t.Fatalf("a schema-8 store: %s", body)
+	}
+	if err := os.Remove(a.d.credentialsPath); err != nil {
+		t.Fatal(err)
+	}
+	if _, body = get(); !strings.Contains(body, `"confirmations":[],"confirmationsUnavailable":"the verb requests could not be read"`) {
+		t.Fatalf("a missing store: %s", body)
+	}
+}
+
+// setVerbColumn sets one column of verb request id straight in the credentials file (a hand edit).
+func setVerbColumn(t *testing.T, path, id, column string, v any) {
+	t.Helper()
+	db, err := sql.Open("sqlite", path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = db.Close() }()
+	if res, err := db.Exec(`UPDATE verb_requests SET `+column+` = ? WHERE request_id = ?`, v, id); err != nil {
+		t.Fatal(err)
+	} else if n, _ := res.RowsAffected(); n != 1 {
+		t.Fatalf("%d rows updated", n)
+	}
+}
diff --git a/host/surfacepage/approve_model.go b/host/surfacepage/approve_model.go
index 2c32b1bf..7ddc6e23 100644
--- a/host/surfacepage/approve_model.go
+++ b/host/surfacepage/approve_model.go
@@ -70,6 +70,23 @@ type RequestsJSON struct {
 	// Decisions are the pending Human.Approve requests a decision card decides (row 201 A4 §2.8).
 	Decisions          []PendingDecision `json:"decisions"`
 	DecisionsTruncated bool              `json:"decisionsTruncated"`
+	// Confirmations are the pending verb requests a confirmation card confirms (row 243 V3 §2.5);
+	// below credentials schema 9, ConfirmationsUnavailable says why there are none.
+	Confirmations            []PendingConfirmation `json:"confirmations"`
+	ConfirmationsUnavailable string                `json:"confirmationsUnavailable"`
+}
+
+// PendingConfirmation is one pending verb request: a summary decoded from its stored bytes, or,
+// when they do not hash to the digest or do not decode, Refused and no summary (never hidden).
+type PendingConfirmation struct {
+	ID         string `json:"id"`
+	Episode    string `json:"episode,omitempty"`
+	Repository string `json:"repository,omitempty"`
+	Head       string `json:"head,omitempty"`
+	SHA        string `json:"sha,omitempty"`
+	Draft      *bool  `json:"draft,omitempty"`
+	Expires    int64  `json:"expires"`
+	Refused    string `json:"refused,omitempty"`
 }
 
 // hex32 decodes 64 hex characters.
````

