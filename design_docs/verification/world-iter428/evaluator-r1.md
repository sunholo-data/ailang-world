# Evaluator report r1 — iteration 428, row 243 V4 sprint S3 (head 786f6743, base 92e1c44e)

## Verdict: PASS — 86 / 100, no BLOCKING findings

| Category | Max | Score |
|---|---|---|
| Correctness vs design §3.1-§3.3/§5.2, plan §9 | 30 | 26 |
| Tests: non-vacuous, kill mutants, class AND unchanged state | 30 | 23 |
| Security / fail-closed | 20 | 19 |
| Coding standards | 10 | 9 |
| Scope and size | 10 | 9 |
| **Total** | 100 | **86** |

## BLOCKING findings
None.

## Non-blocking findings (ranked)
1. **Surviving mutants = test gaps, whole store+authority packages run** (ledger `$HOME/.ailang/state/world-iter428/judge/mutants.tsv`):
   - `J-REQ-PENDING-CONF` (operator_actions.go:~140): pending-cap count `state IN ('requested','confirmed')` -> `('requested')` SURVIVES. The tests fill the cap with requested rows only, so whether a `confirmed` row counts toward the 16 is unpinned. Design §5 says "16 total", so confirmed arguably should count.
   - `J-DDL-STEPS-PK`: `PRIMARY KEY (kind, identity)` -> `(kind, identity, action_id)` SURVIVES. Cross-action uniqueness of an operation/effect identity is only enforced by the lookup in RequestOperatorAction, not pinned at schema level (the property S4 effect reservations rely on).
   - `J-REQ-ACTIVE-ID` (`OR id = ?` in the active count): SURVIVES; a colliding action id would surface as a raw PRIMARY KEY error instead of `ErrOperatorConflict`. Still a refusal; class untested.
   - `J-CONF-ARTSIZE` (`len(artifact) > MaxSurfaceApprovalArtifactBytes`): SURVIVES; the artifact size cap in confirm is untested.
   - `J-REQ-DEDUP-ORDER` (ORDER BY of the review-digest lookup): SURVIVES; probably unreachable (only one non-terminal plus terminal rows can match), likely equivalent.
   - `J-FLOOR-PRESET-REVERT` (cmd/ailang-worldd/episode_preset.go:241, v9 floor -> v10), whole `./cmd/ailang-worldd/` run: SURVIVES. The v9 floor at the two cmd sites has no test (the store-level floor in credentials.go mintSession IS pinned by the plan's V9-BOUND-FLOOR; my `J-FLOOR-MINT` (floor 9->8) survives only because a v8 store never reaches it).
2. **Replay shortcut returns success without consuming a proof**: operator_cas.go:58 `if row.ApprovedBy != "" { return row, nil }` runs after only the digest check, before any proof/artifact check. A replay with a nil or foreign proof returns the already-confirmed row. It changes no state and grants nothing new (the row was confirmed by a genuine proof), so not a forgery, but a replay with a *different* artifact also "succeeds". Tests pin only the same-proof replay. Consider requiring the same artifact bytes in the replay branch.
3. `go.mod` still lists `golang.org/x/sys v0.48.0 // indirect` though the new fd test imports it directly; `go mod tidy -diff` is non-empty. No CI tidy gate found (grep of .github/scripts), so cosmetic.
4. `judge_operator_silent_test.go` is now `//go:build unix` as a whole file, so the existing judge-silence test is also skipped on non-unix. CI runs ubuntu-24.04 only (workflows grep: runs-on ubuntu), and the dev rig is darwin, so no CI leg silently skips it. Acceptable.
5. Intermediate commits M1 (7281c090) and M2 (1e088583): verified green for `go vet ./host/... ./cmd/...` and `go test ./host/store ./host/authority ./host/approval/...` in a throwaway worktree. I did not run the whole `cmd/ailang-worldd` package on them.
6. Plan §4 table not run row-by-row: I ran my own 40 mutants instead (below); several overlap semantically (cap comparison, expiry comparisons, prune age/reservation, confirm generation/digest, replay digest, approver swallow). The planner's own ledger (186 lines) was read, not reproduced. All my ids are non-plan ids (`comm -23` = 45 of 45 mine; `comm -12` = 0), so this under-delivers the "sample >= 8 plan rows" instruction.

## Gates (controller's logs, `$HOME/.ailang/state/world-iter428/gates/`)
All legs rc=0 except `verify_go` rc=1: only `host/broker TestExecProtectedCheckTimeout/reserve-after-before` ("returned at/after the original deadline"), an exec-probe deadline in an area untouched by S3 and on the directive's known load-sensitive list. The summary line `FAILS: routes_changed caller_fence verify_go` is not backed by evidence in legs.txt/fence.txt (`DAEMON_ROUTES: (base ) empty`, `UNTOUCHED_VIOLATIONS: none`, `NEW_PRODUCTION_CALLERS: none`); I verified independently (below).

## My test runs
- `go test -count=1 ./host/store/ ./host/authority/ ./host/approval/...`: all ok (store 11.8s, authority 1.3s, approval 1.0s, approvaltest 1.2s). Output `judge/base-tests.txt`.
- `go test -race -count=2 ./host/store/ ./host/authority/`: ok / ok (73.8s, 17.4s). `judge/race.txt`.
- `go vet ./cmd/... ./host/...`: rc=0. `go test -run '^$' ./...`: rc=0 (no non-ok lines).
- Mutant runs: whole `./host/store/ ./host/authority/` (plus `./host/approval/` for approval edits; whole `./cmd/ailang-worldd/` for the floor mutant).

## Size (unit: `git diff -M --numstat 92e1c44e <c> -- cmd host docs app`, ins+del)
M1 241, M2 738, M3 1,396 (cap 1,400; 4 lines of margin).

## Purity / scope
- `git diff --name-status -M`: 7 added product files (host/store operator_actions.go, operator_cas.go, operator_schema.sql, and three test files, host/authority operator_store_test.go) + 2 plan files; modified: credentials.go, operator_signed.go, judge_operator.go, episode_preset.go, session_new.go, silent/judge tests, and the pre-existing test fixtures. No `.ail`, no `host/surfacepage` or `app.js`/CSP/`surface_page.go`/`daemon.go` change (`git diff --stat -- '*.ail' host/surfacepage` empty; no `mux.Handle` lines added by the diff).
- **No production caller**: `grep -rln` over non-test .go for RequestOperatorAction / ConfirmOperatorAction (store) / ClaimOperatorAction / PruneOperatorActions / RecheckClaimedOperatorAction / OperatorConfirmation / ConfirmationFor finds only the defining files (`approval.ConfirmOperatorAction` and `authority.JudgeOperatorAction` are pre-existing; the controller's positive control `authority.RecheckVerb` was 0 too).
- No `context.Background/TODO/WithoutCancel` in the new store files or new authority tests (grep empty). Store methods go through `requireDeadline`.
- One scope item the plan should have flagged more loudly: `openCredentialsSQLite` now enables `foreign_keys` on every credentials connection (credentials.go, new `_pragma=foreign_keys(1)` and an in-memory `PRAGMA`). It is needed for `operator_action_steps -> operator_actions` FK and pinned by plan mutant FK-CONNECTION-OFF, but it is a behaviour change for all credentials tables (no existing table has an FK; harmless). Not a finding beyond a note.
- Pre-existing test edits audited (`git diff` of every pre-existing _test.go): (a) cmd/ailang-worldd + host/daemon fixtures that step a store back to v5/v6/v8 now prefix `DROP TABLE IF EXISTS operator_*` (needed, otherwise the DROP of verb_requests / step-back leaves v10 tables in a "v8" store); the wanted refusal strings and assertions are unchanged, no assertion removed or loosened; (b) credentials_test/verb_requests_test/mint_requests_test/approvers_test: table-name lists gain the three tables, `user_version` 9 -> 10 (now `CredentialsSchemaVersion`), "newer store" fixture moves 10 -> 11 with text "10"; same strength; (c) judge_test.go `jgExec` rewrites DROP-prefix for the same reason, judge_operator_test.go adds the new subtest. All justified.
- `//go:build unix` + `golang.org/x/sys/unix`: see non-blocking #3/#4.

## Proof transport (S3-D3): forgery attempts (scratch test, deleted; log `judge/forge.txt`)
Scratch type implementing `OperatorConfirmation() *approval.SignedOperatorAction` (the unexported store interface is satisfiable from outside; only the returned sealed object carries authority):
- zero `&SignedOperatorAction{}`: refused `approval: not a verified signature approval`; row still `requested`, passkey counter 1, approver last_used 0.
- nil pointer proof / typed nil / untyped `nil` interface: same refusal, state unchanged.
- recheck result (`RecheckOperatorAction`, SSH) passed as confirmation: refused ErrNotVerified (`confirmation` snapshot is set only when `confirm`); approver use unchanged. My mutant `J-APPR-RECHECK-CONFIRMABLE` (set confirmation on rechecks too) is KILLED by `TestOperatorRunnerRechecks/claimed_store_row`.
- an outsider calling the public `approval.ConfirmOperatorAction` itself with a fabricated Approver (rogue key, enrolled name `mark-ssh`) and a genuine rogue signature gets a valid law object; the store refuses it: `store: the approver's key is not the key that signed: "mark-ssh"` (useApprover compares the enrolled public key / revocation / enrolment time in the writer tx; the passkey path additionally compares rp/origin/backup and `usePasskey` enforces counter > stored and credential id). So the store, not the proof alone, is the arbiter of enrolment; nothing was advanced.
- artifact substitution: `ConfirmationFor` hashes the artifact; mutant `J-APPR-ARTHASH` KILLED by `TestOperatorConfirmStoreBinding/artifact_substitution`.
- Positive control: genuine `JudgeOperatorAction` proof confirms (state `confirmed`, `approved_by=webauthn:mark-mac`).
- Limit: Go cannot stop `unsafe`/reflect writes to the unexported `ok` field of the sealed object; same as every sealed type in the repo.
- What `host/approval/operator_signed.go` newly exposes: exported type `OperatorApprovalBinding` and method `ConfirmationFor(artifact)`; private fields `confirmation`, `artifactHash`. Neither can build a proof (only `operatorProof` sets `ok`); it widens who can READ the verified snapshot, not who can construct one. `authority.VerifiedOperatorAction.OperatorConfirmation()` is the new exported transport and is nil-safe.

## Schema 10 / migration
- DDL + stamp: `upgradeCredentials` appends `operatorSchemaSQL` to the same statement list inside the one transaction as the version stamp; fresh init is `credentialsSchemaSQL + operatorSchemaSQL` in one tx. One shared embed => migration and fresh DDL cannot drift. Mutants killed by TestOperatorSchema10 / TestOperatorSchema10Constraints: `J-ACCESS-VER`, `J-DDL-EXPIRY-ORDER`, `J-DDL-V3-UNIQUE` (via TestOperatorClaimRace); plan ledger lists V10-PARTIAL-MIGRATION/V10-STAMP/V9-SKIP-MIGRATION killed.
- Newer-than-10 refused: test moved to `TestCredentialsNewerThanTenIsRefused` (user_version 11, same text). A pre-v10 binary refuses a v10 store by the existing `case version > credentialsSchemaVersion` branch (credentials.go ~344).
- v9 floor (F6): CORRECT. `CredentialsBindingWriteVersion = 9` keeps the pre-existing meaning of the `< current` check ("store has the binding columns"); v10 is a strict superset of v9, so any v9-assuming bound write on a v10 store writes only v9 tables; the new tables are touched only by the new unregistered methods guarded by `operatorAccess` (`version < 10`). The floor at the two cmd sites is untested (non-blocking #1).
- Parity test: compares `credentialsTables(v)` exact table set plus column/DDL of the upgraded v8->current store against a fresh store (`requireExactTableNames`, `tableDDL`), so it checks table sets and DDL text, not behaviour; it is not vacuous (V10-FRESH-OMIT / V10-TABLE-SET kills in the plan ledger), though behaviour parity is what Constraints test adds.

## Finish shape (S3-D2)
Schema 10 carries `kind IN ('start','finish')`, state `awaiting-pr-confirmation`, unique finish-only `v3_id`, and finish step kinds (`commit`,`v3-card`,`push`,`label`,`pull-request`,`record`,`cleanup`). That is within D-WORLD-120 = A and design §3.1-§3.3 (a v10 that omitted the finish shape would force another irreversible bump). Nothing makes a finish path reachable: `RequestOperatorAction` accepts `Kind` only via `approval.CanonicalOperatorAction`, no route/runner exists, the methods have no caller. The DDL guards (`v3_id` only with kind finish; awaiting-pr-confirmation needs v3_id) are killed by TestOperatorSchema10Constraints.

## Admission / dedup / expiry / prune (measured)
- Cap 16: `J-CAP-GT` (`>`), `J-CAP-15` both KILLED by TestOperatorAdmissionBoundaries (both sides of the boundary). Confirmed-row counting: gap (SURVIVED).
- Expiry: `now >= expires` is expired everywhere. `J-ADMIT-EXP-EQ` KILLED (TestOperatorAdmissionBoundaries at_expiry), `J-READ-EXP-LT` KILLED (TestOperatorPrune), `J-EXPIRE-SQL-LT` KILLED (TestOperatorAdmissionBoundaries), `J-CONF-CLOCK-EQ` KILLED (TestOperatorConfirmStoreBinding/at_expiry), claim expiry covered by TestOperatorClaimBeforeExpiry exact-expiry case.
- Op-id / review-digest dedup: `J-REQ-REVIEW-LINK` KILLED (TestOperatorAdmissionBoundaries/operation_changed, TestOperatorDuplicateTap); `J-REQ-DEDUP-EXCL` KILLED (TestOperatorPrune).
- Prune keeps tombstones/uncertain: `J-PRUNE-AGE` and `J-PRUNE-RESV` KILLED (TestOperatorPrune).
- Refusal assertions: the claim table test asserts the error class and re-reads the row (state, owner, generation) and row counts after every case; confirm tests assert class via `errors.Is` and call `assertUnconfirmed` (state, approved_by, artifact, generation, digest, expiry, passkey counter, approver use) after every refusal. Admission refusals use `operatorCounts`. Good.

## Confirm CAS / Claim CAS
- Counter exactly once / replay: TestOperatorConfirmReplay (ssh + passkey, store reopened) asserts DeepEqual row, `last_used_at` stays 1010, counter stays 2. A store error never maps to success: `J-CONF-SWALLOW-APPROVER` (approver error -> `return row, nil`) KILLED (TestOperatorConfirmCrash/revoked, TestOperatorConfirmReadOnly); crash test injects a trigger fault and asserts SQLite code 1811 and unchanged state.
- `ApprovedBy` form persisted: `J-CONF-APPROVEDBY-FORM` (persist `binding.Name`) KILLED (TestOperatorConfirmReplay, TestOperatorRunnerRechecks/claimed_store_row); the claimed-row recheck with the stored value succeeds, a replacement identity gives `ErrOperatorApprover`.
- Claim predicates each independently killed: state (`J-CLAIM-PRE-STATE` by TestOperatorClaimBeforeExpiry), digest (`J-CLAIM-CAS-PRE-DIGEST`), heartbeat freshness (`J-CLAIM-HB-FRESH`), future heartbeat (`J-CLAIM-HB-FUTURE`), readiness (`J-CLAIM-READY`), runner owner (`J-CLAIM-OWNER-PRED`), generation increment (`J-CLAIM-GENINC` by TestOperatorClaimRace), owned-sibling NOT EXISTS (`J-CLAIM-SIBLING-PRE` by TestOperatorClaimRace, with the unique index dropped), exactly-one-row (`J-CLAIM-ROWCOUNT`). Two concurrent claimers: TestOperatorClaimRace -> one nil, one `ErrOperatorConflict`. Note: the Go prechecks duplicate several SQL predicates inside one serialized tx, so single-predicate SQL mutants are largely equivalent by design (plan ledger marks them S as well).
- Claimed-row recheck uses the runner's own frozen action: `testOperatorClaimedRecheck` blanks the row's digest/challenge, bumps its expiry to 1601, and shows (a) the genuine path passes, (b) a different frozen intent fails, (c) recorded expiry fails at 1600, (d) a replacement approver string fails, (e) missing artifact fails. Good.

## Judge-409 note (a): reproduction
- Base 92e1c44e, `println("x")` inserted at the start of `JudgeOperatorAction`, `go test -count=1 ./host/authority/`: `ok` (SURVIVES) — throwaway worktree, ledger `J409A-BASE-println`.
- Head 786f6743, same edit: KILLED by `TestOperatorJudgesWriteNothing` (fd-level dup2 capture) — ledger `J409A-HEAD-println`. Also `TestOperatorStoreWritesNothing` covers the store methods.

## Mutant table (ledger: `$HOME/.ailang/state/world-iter428/judge/mutants.tsv`, 55 lines incl. 2 invalid/compile retries and 1 base-worktree line; 45 unique ids, all NOT in the plan's ids: `comm -23 my-ids plan-ids` = 45, `comm -12` = 0)
Controls: `J-CTRL-LIVE` (comment-only edit) SURVIVED as required; `J-CTRL-DIE` (cap 16 -> 1000) KILLED by TestOperatorAdmissionBoundaries.
KILLED (killer named; every killer is a behavioural Test*, none is an AST/static/inventory test): J409A-HEAD-println (TestOperatorJudgesWriteNothing), J-CAP-GT, J-CAP-15, J-ADMIT-EXP-EQ (TestOperatorAdmissionBoundaries); J-READ-EXP-LT, J-PRUNE-AGE, J-PRUNE-RESV, J-REQ-DEDUP-EXCL (TestOperatorPrune); J-EXPIRE-SQL-LT (TestOperatorAdmissionBoundaries); J-REQ-LIFETIME-GE (ten tests incl. TestOperatorAdmissionBoundaries, TestOperatorConfirmCrash...); J-CONF-CLOCK-EQ, J-CONF-GENPRED, J-CONF-CLOCKBEHIND, J-CONF-REPLAY-DIGEST (TestOperatorConfirmStoreBinding); J-CONF-APPROVEDBY-FORM (TestOperatorConfirmReplay, TestOperatorRunnerRechecks); J-CONF-SWALLOW-APPROVER (TestOperatorConfirmCrash/revoked); J-CONF-RP-ONLY (TestOperatorConfirmCrash/changed-origin,changed-backup); J-CONF-NILPROOF (TestOperatorStoreWritesNothing); J-CLAIM-HB-FRESH, J-CLAIM-PRE-STATE, J-CLAIM-CAS-PRE-DIGEST (TestOperatorClaimBeforeExpiry); J-CLAIM-HB-FUTURE, J-CLAIM-OWNER-PRED, J-CLAIM-READY (TestOperatorClaimRace, TestOperatorRunnerRechecks, TestOperatorStoreWritesNothing, ClaimBeforeExpiry); J-CLAIM-SIBLING-PRE, J-CLAIM-GENINC, J-CLAIM-ROWCOUNT (TestOperatorClaimRace, ConfirmCrash drift cases); J-APPR-ARTHASH (TestOperatorConfirmStoreBinding/artifact_substitution); J-APPR-RECHECK-CONFIRMABLE (TestOperatorRunnerRechecks/claimed_store_row); J-REQ-REVIEW-LINK (TestOperatorAdmissionBoundaries/operation_changed, TestOperatorDuplicateTap); J-ACCESS-VER (TestOperatorSchema10); J-DDL-V3-UNIQUE (TestOperatorClaimRace); J-DDL-EXPIRY-ORDER (TestOperatorSchema10Constraints).
SURVIVED (whole packages): J-CTRL-LIVE (control), J-REQ-ACTIVE-ID, J-REQ-PENDING-CONF, J-CONF-ARTSIZE, J-REQ-DEDUP-ORDER, J-DDL-STEPS-PK, J-FLOOR-MINT, J-FK-INMEM, J-FLOOR-PRESET-REVERT (whole cmd/ailang-worldd).
COMPILEERR/INVALID (not counted): `J409A-HEAD-println` first try (2 matches, no change; re-run with a unique anchor), `J-REQ-REVIEW-LINK` first try (unused import), both re-run and recorded.
Per-milestone non-vacuity: M1 guards (version access, DDL order/unique checks) are killed by TestOperatorSchema10(+Constraints); M2 guards (cap, expiry, dedup, prune) by TestOperatorAdmissionBoundaries/Prune/DuplicateTap; M3 guards (confirm, claim, proof) by TestOperatorConfirm*/ClaimBeforeExpiry/ClaimRace/RunnerRechecks — each named test failed with only its target guard changed in the ledger.

## Plan §1.4 false premises, S3-D1..D5 and §8 deployment risk
- §1.4 premises (e.g. "store may not import authority", "V3 row persists so the finish link can be an FK", "existing rollback restores the old build safely") are consistent with what I measured: the transport via `OperatorConfirmation()` is sound; v3_id is a logical reference (no FK), consistent with V3 pruning; the rollback premise is F4.
- D1 (shared DDL), D2 (finish shape), D3 (sealed-proof transport), D4 (operation identity in the three-table footprint) are all supported by the diff; D5 not separately exercised.
- §8 / F4: measured RIGHT. `cmd/ailang-worldd/service_upgrade.go`: swap -> `rotate`/`Kickstart` -> `waitForBuild`; on failure (or on a failed `saveUpgrade` after health) `rollbackUpgrade` renames `.prev` back and kickstarts it; a v9 binary refuses a v10 store (existing newer-store refusal), so the service stops with data intact and the candidate retained as `.failed`. Classification: NON-BLOCKING for merging S3 to dev (upgrade is a manual `--binary <file>` operator step; nothing deploys on merge; S3 has no consumer), but a FOLLOW-UP to queue before the next attended `service upgrade` of Mark's home service: backup of credentials.db before migrating (or rollback that restores the candidate and warns), and a runbook line. This is the same shape as every prior schema bump, it is documented in the plan §8, and the data is not lost.

## Files
- Report: `$HOME/.ailang/state/world-iter428/evaluator-r1.md`; ledger `judge/mutants.tsv`; forgery log `judge/forge.txt`; per-mutant outputs `judge/mut-*.txt`. Eval worktree `git status --short` clean; throwaway worktree removed; scratch forgery test deleted.
