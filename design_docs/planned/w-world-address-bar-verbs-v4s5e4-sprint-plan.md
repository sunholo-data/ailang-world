# Sprint plan — w-world-address-bar-verbs-v4, sprint S5e4 (iteration 438, row 243)

**Pre-caller pins and episode GC.** Base `684069f2118264d0797a0f7575b3c4d837c7fef3`. Product patches in Appendix P are incremental. The executor runs on `codex:gpt-6.2-sol`, applies them in the order M1, M2, M3, and hands the controller cumulative snapshots.

## 0. Binding constraints

- The reviewed [V4 design r1](w-world-address-bar-verbs-v4.md) and quorum r2 PROCEED 3/3 stand. D-WORLD-119=A and D-WORLD-120=A were attended at `82f05004`. Fixtures use temp stores and test keys under D-WORLD-109.
- S5e3 #674 merged at `22bb7804`. `git diff -M --numstat eaa202b1 22bb7804 -- cmd host docs app` sums to +1,625/−25, or 1,650 changed lines. [The iteration 437 judge](../verification/world-iter437/evaluator-r1.md) gave 79/100 and zero blocking findings. `git diff --name-status 22bb7804 HEAD` reports seven paths under `design_docs/` (two verification records and five charter/dashboard/index/log/status documents, record #675).
- M1 and M2 are test-only. M3 only reduces worktree removals, in the daemon and CLI. The start, reservation, mint and publication paths gain zero production callers. S5f owns activation. Routes, schema, and `.ail` have zero diff.
- Size cap 1,400 changed lines, with additions plus deletions from `git diff -M --numstat 684069f2 HEAD -- cmd host docs app`. Always report the unfiltered count too. If revision pushes the product above the cap, retain M2 and M3, and move the tail of M1 to Appendix C after exhausting optional M4.
- Apply only the exact Appendix P patches. Never edit the real checkout, launchd driver, operator homes, credential stores or services during verification. The controller owns Git commits. Treat bind-denied runs as **UNINFORMATIVE UNDER SANDBOX** and compare whole-package fail sets against the same-head external recordings.

## 1. Pre-flight

### 1.1 Baseline at 684069f2

`git rev-parse HEAD` printed `684069f2118264d0797a0f7575b3c4d837c7fef3`. Initial `git status --short` printed nothing. `git archive 684069f2 | tar -x -C "$S/base"` made the pristine baseline before staging any milestone. The first gate check completed with `ALL_DONE` at 152 s. Build, both vets, all inventory legs, the GC probe, and store plain/race passed. Whole daemon and CLI plain/race had actual loopback denials. The first run skipped the repository verifiers because they allocate implicit temp files and `verify_go.sh` needs a committed-blob Git reference. The supplementary gate run used scratch-contained paths for their shell allocations with the real Git binary and a read-only committed-blob reference. See §6 for the invocation and §1.5 for raw receipts.

### 1.2 Current inputs and coupling inventory

Commands refer to the exact base. The base output is in `$S/base-anchors.txt` and the stage `grep -n` output is in `$S/stage-anchors.txt`. All stage paths in this document are relative to an archive root.

| Input | Command and observed base output |
|---|---|
| File store N22/N23 | `git grep -nE 'Headless\|count\(\*\)' -- host/store/operator_files.go` gives :36 headless, :44 exact worktree count. The functions start at :23 and :67. |
| File store N24/N25 | `git grep -nE 'lock.dbPath\|value == ""' -- host/store/operator_files.go` gives lock comparison :91, empty identity/proof :109. The record function starts at :87. |
| File store N27/N28/N29/N31 | `git grep -nE 'now < row.CreatedAt\|ApprovedBy !=\|Start.Allowance\|now >= row.ExpiresAt' -- host/store/operator_files.go` gives creation :194, signer identity :204, allowance :219, expiry :191. `git grep -n 'func admitOperatorFiles' -- host/store/operator_files.go` starts :189. |
| Worktree reserve | `git grep -nE 'episode_starts WHERE\|session_credentials WHERE\|now-10\|func .*ReserveOperatorWorktree' -- host/store/operator_worktree.go` gives :13 entry, :56–57 existing-start and credential guards, :59 heartbeat window argument. |
| Worktree read and dispose | `git grep -nE 'func .*OperatorWorktree\|AND owner = \?\|NOT EXISTS' -- host/store/operator_worktree.go` gives :73 read, :89 phase, :90 read other-kind, :108 dispose, :138 dispose other-kind. |
| Mint files stage | `git grep -nE 'fileStage\|readOperatorMint\(ctx, (c.db\|tx)' -- host/store/operator_mint.go` gives :246 and :302 `false`, function :249, SQL placeholder :266 and bind :267. `grep -nF 'later effect reserved' host/store/operator_mint.go` gives :269 and wraps `sql.ErrNoRows`. |
| GC callers and age | `git grep -nE 'startAutoGC\|credentialsPath: credPath\|IsActive\|RunEpisodeGC\|clean & merged\|autoGCOlderThan' -- host/daemon/{auto_gc,daemon}.go cmd/ailang-worldd/episode_gc.go` gives daemon resolved path `daemon.go:844`, startup :957, `auto_gc.go:34` options, :46 72 h, :73 pass, :81 active predicate, :101 GC, :165 candidate, CLI :74 invocation. |
| Durable operator phases | `nl -ba host/store/operator_schema.sql` shows :14–15 state strings, :26 terminal invariant, :46–47 alias comment, :52–61 step ownership and phase CHECKs. `git grep -nE 'func .*OperatorStep\|phase' -- host/store/operator_steps.go` gives :17/:34 reserve, :45/:57–58 observe, and :72 read. `git grep -n 'Uncertainty retains' -- host/store/operator_worktree.go` gives :103–104. |
| Store read harness | `git grep -nE '^func (openOperatorReadOnly\|requireDeadline\|openTestCredentials\|operatorClaimFixture\|operatorExec\|getOperator\|assertReservedOperatorStep\|insertEpisodeStart\|testCtx)' -- host/store` gives helper declarations in `operator_runner.go:125`, `store.go:250`, `credentials_test.go:121`, `operator_cas_test.go:12,279`, `operator_actions_test.go:37,46`, `credentials.go:525`, `mint_requests_test.go:175`. `sed -n '125,154p' host/store/operator_runner.go` shows read-only DSN, busy timeout zero and context-bounded schema query. |
| CLI configuration | `git grep -nE '^func (defaultWorldPaths\|homeServicePaths\|readServiceRecord\|runEpisodeGC\|TestEpisodeGCFailSafeAndReclaims)' -- cmd/ailang-worldd` gives `setup.go:74`, `service.go:49,271`, `episode_gc.go:46`, and `episode_gc_test.go:27`. `sed -n '271,289p' cmd/ailang-worldd/service.go` prints the decoder and loopback, schema, label and home path validation. M3 reads that existing record with the same validation. |
| Production caller fence | `/bin/bash "$S/callers.sh" "$S/st-M3" "$S/callers-final.txt" "$S/callers-base.txt"` compares the production census. The two positive controls print the attended `return e.reconcileRunnerFiles(args[1], out, errOut)` and a test call to `publishOperatorStartFiles`. |

The stage measurement `grep -nE 'HasOpenOperatorStep|terminal_at|rootErr|func .*GC' "$S/st-M3"/{host/store/operator_episode_gc.go,host/daemon/auto_gc.go,cmd/ailang-worldd/episode_gc.go}` prints the read entry at `operator_episode_gc.go:8`, joined SQL at :22–24, option at `auto_gc.go:43`, daemon wiring :76/:85–86, guard :157–175, CLI resolution at `episode_gc.go:51–74` and CLI wiring :106/:113–114. The caller census is base-identical. The positive controls are the actual attended `reconcileRunnerFiles` call and the direct publication test call.

Inventory command pattern is `go test -v -count=1 -timeout 2m ./<package> -run '^<test>$'` with §6 environment. The source command is `git grep -nE '^func Test(DaemonDependencyAllowlist|WorldBoundaryDependencyAllowlist|ProductionGoSurface|ProductionContextRoots|CredentialsMethodsRequireADeadline|StoreSurfaceTakesContext|BoundaryASTWriteGuard|OperatorNeverDispatch|CLICredentialWriterSites)' -- cmd host`.

| Inventory and measured declaration | Change and reason |
|---|---|
| `TestDaemonDependencyAllowlist`, `host/daemon/daemon_test.go:899` | No expectation changes. Both production packages already depend on `host/store`. |
| `TestWorldBoundaryDependencyAllowlist`, `host/boundary/allowlist_world_test.go:758` | No expectation changes. The new read uses stdlib and the existing store package. |
| `TestProductionGoSurface`, `host/store/context_roots_test.go:239` | No expectation changes. The new Go file is in the allowed `host/store` tree. |
| `TestProductionContextRoots`, `host/store/context_roots_test.go:134` | Root pins stay identical. CLI and daemon pass their existing bounded GC contexts, and the lookup receives one. Named gate under §7, also reproducible as `go test ./host/store/ -run '^TestProductionContextRoots$'`. `grep -E 'root references\|root census' "$S"/gate-{base-final-script,M3}/inv-store.txt` reads 59 roots on both, with production-file counts 342 → 343. |
| `TestCredentialsMethodsRequireADeadline`, `host/store/credentials_test.go:277` | Adds one explicit exercised entry for the package-level credential lookup, like the existing path readers. A runtime store refusal also pins it. |
| `TestStoreSurfaceTakesContext`, `host/store/deadline_guard_test.go:43` | No expectation changes. The reflection scan is of `*Store`. M3 adds a package function with `context.Context` first and no new `Credentials` method. |
| `TestBoundaryASTWriteGuard`, `host/boundary/allowlist_world_test.go:1323` | No expectation changes. Store opens read-only and runs a SELECT. |
| `TestOperatorNeverDispatch`, `host/verbs/operator_policy_test.go:69` | No expectation changes. No verb table edit or dispatch. |
| `TestCLICredentialWriterSites`, `cmd/ailang-worldd/cli_credential_writers_test.go:18` | No expectation changes. CLI calls a bounded path reader. SQLite fixture writes are test-only. |

### 1.3 Decisions

- Treat observation and disposal separately. The durable action terminal set is `succeeded, failed, expired, cancelled` (DDL §1.2). The step table has two phases, `reserved` and `observed`. Disposal is denoted by the parent action. An effect blocks if it is reserved, even beside a terminal parent, or if it is observed and its parent still lacks `terminal_at`. `operation` rows are observed operation aliases without owner or receipt, and do not block by themselves. An expired timestamp alone does not dispose a running or uncertain action.
- Resolve attended CLI lookup from explicit `--credentials-db PATH` when supplied, else from a recorded home service whose workspace physically matches, else from the default store only for the default workspace when the record is absent. A malformed service record or unmatched custom workspace leaves evidence unavailable. Both paths use the same read-only store function. Its SQL is bounded by the existing two-minute pass context and opens an existing SQLite file with zero busy wait. Old, missing, malformed or timed-out evidence fails closed. The lookup runs immediately before deletion, after the existing dirty and merge tests.
- M2 selects the store operations directly. The read control succeeds before reservation. The abort control revokes and stops without a reservation. The abort test obtains the reviewed receipt before reserving files, so the public read cannot mask the abort guard.
- This is one host-side housekeeping predicate, operating on the installed credentials schema and Git worktrees. It belongs at the Go host boundary under Standards S2/S3.

### 1.4 Premises found false or incomplete at HEAD

| Premise | Remeasurement and consequence |
|---|---|
| N33/N34 edit direction | The base `git grep` in §1.2 prints `false` at :246 and :302. **The mutant flips each to true.** The CLI masking claim is consistent with those two independent barriers. |
| Allowance layer reachability | `git grep -nF 'sql.ErrNoRows' -- cmd/ailang-worldd/operator_files_test.go` prints the real allowance-path refusal assertion at :598. The exact original N01 edit still exercises that mask. M1 constructs a signed mixed-origin start with matching credential, observed mint receipt and checksum in the real database, so N29 becomes independently observable and its test expects `ErrOperatorConflict`. |
| Existing file-store coverage | `git grep -nE 'ReadOperatorFiles\|RecordOperatorFiles' -- 'host/store/*_test.go'` prints only `credentials_test.go:317–318`. Those calls exercise the no-deadline gate. Positive control is the real caller in `cmd/ailang-worldd/operator_files_test.go`. M1 is the first store behavioural coverage. |
| Creation boundary also checked in reservation SQL | `git grep -nF 'created_at <= ?' -- host/store/operator_steps.go` prints :37, a second clock check. M1 invokes `admitOperatorFiles` directly on a real readback with a single creation-field change and checks the durable rows again. |
| GC evidence boundary at base | `grep -F 'dry-run GC candidate' "$S/gate-base/gc-probe.txt"` prints both `reserved` and `published`, reason `clean & merged`, evidence unchanged. The M3 probe uses live GC in the fixture for the protected arms and a dry run for the terminal arm. |

### 1.5 Size and CI race cost

| Milestone | Increment | Cumulative +/− | Filtered | Unfiltered product | Snapshot status/deletions/hashes |
|---|---:|---:|---:|---:|---:|
| M1 | 320 | +320/−0 | 320 | 320 | 1/0/1 |
| M2 | 52 | +372/−0 | 372 | 372 | 2/0/2 |
| M3 | 557 | +917/−12 | 929 | 929 | 9/0/9 |

Product margin is 471 lines. If the controller commits both plan records first, add their **2025** lines to each unfiltered HEAD total.

| Run | Receipt |
|---|---|
| base first | `GATES_DONE 152s FAILS: none UNINFORMATIVE: pkg-host-daemon race-host-daemon pkg-cmd-ailang-worldd race-cmd-ailang-worldd` |
| base complete verifier | `GATES_DONE 675s FAILS: none UNINFORMATIVE: verify_go pkg-host-daemon race-host-daemon pkg-cmd-ailang-worldd race-cmd-ailang-worldd` |
| M1 | `GATES_DONE 100s FAILS: none UNINFORMATIVE: none` |
| M2 | `GATES_DONE 64s FAILS: none UNINFORMATIVE: none` |
| M3 initial | `GATES_DONE 928s FAILS: none UNINFORMATIVE: verify_go pkg-host-daemon race-host-daemon pkg-cmd-ailang-worldd race-cmd-ailang-worldd` |
| M3 final | `GATES_DONE 172s FAILS: none UNINFORMATIVE: pkg-host-daemon race-host-daemon pkg-cmd-ailang-worldd race-cmd-ailang-worldd` |

The Go toolchain was go1.26.6 darwin/arm64. The real Git executable was `/opt/homebrew/bin/git`. Native build/vet, Linux vet, the requested inventories and named acceptance all returned rc 0 for the compiled stages. The AILANG v0.41.0 gate on base and M3 verified 16 required identities, 40 named tests, 594 se-tools tests, and the 9/9 world-package checks. Those runs used the scratch-temp relocation described in §6.

The plain Go verifier saw its race control fire and ran `go test ./... -count=1`. The base and the initial M3 archive reported actual `bind: operation not permitted` output. The final M3 path-error simplification was replayed through the other gate legs. The broad verifier on the exact final archive remains for the controller. Whole daemon and CLI plain/race are **UNINFORMATIVE UNDER SANDBOX**. The recorded local store fail sets are empty. Unrelated load-sensitive rows are compared on the controller's same-head base in §7. `uptime` was recorded in each log, and the sandbox denied the process-list command. Base receipts stay in `$S/gate-base` and `$S/gate-base-final-script`. M2's first gate process was invalidated by updating the gate script while Bash was reading it. The frozen-script M2 replay is recorded independently in `$S/gate-M2-final`.

The historical 11-test CLI process baseline is `$HOME/.ailang/state/world-iter437/plan-scratch/controller-base/pkg-cmd-ailang-worldd.base-fails`. See §8 for its names. The controller compares with fresh `684069f2` captures in `$S/controller-base`.

Go package seconds, with the same command on each archive. Whole runs use `go test -count=1 -race -timeout 8m ./<pkg>` and the base archive's last run with the matching leg. Focused runs use `go test -race -v -count=1 -timeout 2m ./<pkg> -run <selector>`. Receipts are `$S/gate-base-final-script`, `$S/gate-M3-final` and `$S/race-focused-final.log`.

| Package | Whole base → M3 (delta) | Focused base → M3 (delta) |
|---|---|---|
| `host/store` | 44.317 → 41.085 s (-3.232 s) | 1.507 → 3.812 s (+2.305 s) |
| `host/daemon` | UNINFORMATIVE UNDER SANDBOX | 1.592 → 2.001 s (+0.409 s) |
| `cmd/ailang-worldd` | UNINFORMATIVE UNDER SANDBOX | 2.345 → 2.357 s (+0.012 s) |

The focused daemon selector is `^TestEpisodeGCOperatorSteps$` and the CLI selector is `^TestEpisodeGCFailSafeAndReclaims$`. The daemon selector discovers zero tests at base, so its package-time delta is the measured cost of the new test. The CLI filters the same existing test at both trees. The store selector and every exact invocation are in `$S/measure-race.sh`.

### 1.6 Portability

Store fixtures are migrated SQLite. The signed fixture uses the repo SSH signer. Git runs as the real executable with worktrees inside `t.TempDir()`. Path comparisons resolve both sides with `filepath.EvalSymlinks` before any removal. HEAD and porcelain output are compared as bytes for live refusals. The GC fixture clock is based on the observed directory mtime plus 73 hours. Other tests use bounded contexts and injected operation times. There are no new sleeps or asynchronous waits.

## 2. How the executor applies a sub-milestone

Run from the sprint worktree. Set `M` to the next unapplied milestone. The runner can use one working tree through M3, regardless of whether the controller commits after each snapshot. The base must be an ancestor, and these patches apply in the order M1, M2, M3.

````bash
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter438-row243-v4s5e4
STATE=$HOME/.ailang/state/world-iter438
S=$STATE/plan-scratch
T=$STATE/executor
B=684069f2118264d0797a0f7575b3c4d837c7fef3
PLAN=$W/design_docs/planned/w-world-address-bar-verbs-v4s5e4-sprint-plan.md
M=M1  # next unapplied milestone, in order
mkdir -p "$T"
unset SNAP_PARENT SNAP_STATUS_LIST
git -C "$W" merge-base --is-ancestor "$B" HEAD
if [ "$M" = M1 ]; then
  git -C "$W" diff --quiet "$B" -- cmd host docs app
  test -z "$(git -C "$W" ls-files --others --exclude-standard -- cmd host docs app)"
fi
cp "$S"/{bounded.py,failset.py,verify-scratch.py,callers.sh,callers-base.txt,probe.sh,operator_gc_probe_test.go,operator_gc_after_test.go} "$T/"
for name in accept.sh gates.sh snapshot.sh; do
  awk -v n="$name" '$0 == "````bash " n {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$name"
  test -s "$T/$name" && /bin/bash -n "$T/$name"
done
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
(cd "$W" && /usr/bin/patch -p1 -E < "$T/$M.patch")
/bin/bash "$T/accept.sh" "$W" "$M" fast
/bin/bash "$T/snapshot.sh" "$W" "$M"
````

Apply each milestone once. A `STATUS`, `DELETED` and `MANIFEST.sha256` receipt are required for each cumulative snapshot.

## 3. Milestones

### M1 — Real-store pins for file and worktree guards

| Exact judge target | Sole row or measured mask |
|---|---|
| N22 `Headless` | `TestOperatorFileStoreReadPins/headless` refuses with `ErrOperatorConflict`. The fixture re-signs a headless-disabled action over the same minted evidence. |
| N23 worktree count | `TestOperatorFileStoreReadPins/duplicate_worktree` refuses with `sql.ErrNoRows`. One matching observed worktree plus a second reservation separates the count from the join. |
| N24 path binding | `TestOperatorFileStoreReservePins/foreign_lock` refuses with `ErrOperatorConflict`. The reviewed mint and store rows remain intact. |
| N25 empty identity | `TestOperatorFileStoreReservePins/empty_identity` requires `ErrOperatorInput`. Proof, phase and all other fields are valid. |
| N27 creation boundary | `TestOperatorFileStoreReservePins/before_created` invokes admission with a future action creation and a live signed proof and session, bypassing the second reservation SQL guard. |
| N28 approver identity | `TestOperatorFileStoreReservePins/different_approver` leaves the signer enrolled and signed artifact intact, mismatches only the row identity, and requires conflict. |
| N29 allowance | `TestOperatorFileStoreReservePins/allowance` uses a canonical, signed mixed-origin fixture. `CommitOperatorMint` produces matching start and mint bytes under the fixture. File reserve independently rejects that allowance. |
| N31 row expiry | `TestOperatorFileStoreReservePins/action_expired` invokes the clock at exactly the action expiry, with the credential still live, and requires `ErrOperatorExpired`. |
| N01 CLI allowance | Exact edit replay with existing `TestOperatorFileAllowanceOrigin`. The missing observed mint row refuses with `sql.ErrNoRows` before `operatorFilePaths`. Store pin N29 exercises the later store layer independently. |
| Judge 436 heartbeat lower bound | `TestOperatorWorktreeStoreReservePins/heartbeat_edge` pins equality at `seen_at=now-10` after a successful fresh claim. |
| Judge 436 start/session absence | `TestOperatorWorktreeStoreReservePins/{existing_start,existing_credential}` independently injects only one existing row and expects conflict. The original composite mutant and each single predicate deletion are run. |
| Judge 436 read phase and read other-kind | `TestOperatorWorktreeStoreRecoveryPins/{observed,read_other_kind}` requires `ErrOperatorMissing`. |
| Judge 436 dispose other-kind | `TestOperatorWorktreeStoreRecoveryPins/dispose_other_kind` requires conflict after a real attempted observation and checks transactional rollback of every step. |

The tests reuse the migrated store, fixture request/claim, SQLite helper and real SSH confirmation primitives. Controls establish successful reserve and read. Every refusal snapshots action, step phase/receipt/version/owner, session rows and start before and after.

### M2 — Independently pin the mint read and abort after files reserve

| Exact judge target | Sole killing row |
|---|---|
| N33 ReadOperatorMint false → true | `TestOperatorMintStoreFilesFence/read` calls the read directly. It must return a zero receipt and `sql.ErrNoRows`. |
| N34 AbortOperatorMint false → true | `TestOperatorMintStoreFilesFence/abort` calls abort directly with an earlier genuine receipt. It must return `sql.ErrNoRows` and keep the exact action, token credential, episode start and reservations. |

Store a real observed mint, then reserve files through `RecordOperatorFiles`, check the reservation readback, and assert both refusals against the snapshot. `abort_without_files` tests the real revoke-and-stop control. An extra CLI test would spend a second CLI test slot and adds no independent evidence.

### M3 — #667 operator-step-aware episode GC

- Add `store.EpisodeHasUnresolvedOperatorSteps(ctx, path, episode)` and wire it into both existing GC paths via `EpisodeGCOptions.HasOpenOperatorStep`. The daemon passes its resolved credentials path. CLI path policy is in §1.3. On a failed, missing or canceled lookup the named reason is `operator step lookup unavailable`. The protected-worktree reason is `non-terminal operator step`. Clean merged old work with a successful negative lookup keeps the inherited candidate reason.
- `TestEpisodeHasUnresolvedOperatorSteps` covers no action, alias only, reserved, uncertain, observed active, each of the four terminal action states, terminal-with-reservation, and unrelated episodes. Its readback compares exact SQLite bytes. `TestEpisodeOperatorStepLookupRefusals` covers deadline missing/expired, absent database, old schema and failed SQL.
- One in-process daemon test `TestEpisodeGCOperatorSteps` runs real Git in a temp repository, and exercises true, false, nil, error, timeout, canceled and configured daemon wiring. Live refusals preserve HEAD and status as bytes. The daemon pass gets a supplied clock and returns its report, with its existing logging intact.
- Extend only the existing CLI test `TestEpisodeGCFailSafeAndReclaims`. It checks a real reserved ledger, matching and mismatched workspace, explicit override, unavailable HOME/record/path, and a subsequent terminal observation that permits the original Git/caches reclamation. The CLI help explains the credential selector. No socket or service is used.
- Use the before/after probe in §6. With the CLI each explicit custom path is scoped to the exact scratch fixtures. The only disposition the probe writes is fixture-local. It uses no operator prompts.

### 3.1 Design §5.3 row assignment

| Reviewed row | Assignment |
|---|---|
| Worktree window and collisions | M1 pins the fresh claim window, start/session absence and reserved-phase recovery through the real store. |
| Mint binding, rollback and receipt window | M1 pins admission to file publication. M2 protects the later files reservation from both mint recovery entries. |
| Task/token/MCP file mismatch and recovery | The existing S5e3 CLI pins remain. M1 adds independent store coverage. The three optional path-window pins are measured in Appendix C. |
| Start/launch, browser lifetime and herdr hand-off | S5f/S6, S8, S12 follow the reviewed design. M3 provides the retention gate needed before an operator can start creating live managed worktrees. |

### 3.2 Design §5.5 and §2.2 activation

The caller census in §1.2 covers the existing attended recovery and the dormant preparation, mint and publication helpers. M3 only reads effect evidence. The GC lookup is an eligibility check immediately before removal, and S5f must arrange exclusion between activation and GC if it permits concurrent reservations against the same eligible worktree. The paired process and owner/old-binary recovery are Appendix C scope.

## 4. Mutation table

The ledgers distinguish exact original judge substitutions from production-diff substitutions. The mutation counts cover executable predicates and the named judge edits. Every row names the substitution (literal old and new source), exact stage line, package, selector and outcome in `$S/mutants.json` and `$S/ledger-{named,whole}.jsonl`. Each raw log is indexed in its ledger. A behavioural kill requires a named `--- FAIL:` with no bind-denial and, for whole runs, a failure added to the package base fail set. Compile errors or unrelated failures do not qualify. `DIE-POSITIVE` and `LIVE-POSITIVE` are harness controls, with individual 436 start/session cuts recorded separately.

- **Diff** (37 substitutions). Named 37 KILLED. Whole 13 KILLED, 24 UNINFORMATIVE UNDER SANDBOX.
- **Judge** (19 substitutions). Named 15 KILLED, 4 SURVIVED. Whole 15 KILLED, 4 UNINFORMATIVE UNDER SANDBOX.
- **Controls** (4 substitutions). Named 3 KILLED, 1 SURVIVED. Whole 3 KILLED, 1 SURVIVED.

| Ledger | ID and exact stage anchor | Named test | Named | Whole |
|---|---|---|---|---|
| diff | `GC-ST-DEADLINE` `host/store/operator_episode_gc.go:9` | `TestEpisodeOperatorStepLookupRefusals/no_deadline` | KILLED | KILLED |
| diff | `GC-ST-OPEN-ERROR` `host/store/operator_episode_gc.go:13` | `TestEpisodeOperatorStepLookupRefusals/missing_store` | KILLED | KILLED |
| diff | `GC-ST-OLD-SCHEMA` `host/store/operator_episode_gc.go:17` | `TestEpisodeOperatorStepLookupRefusals/older_schema` | KILLED | KILLED |
| diff | `GC-ST-QUERY-ERROR` `host/store/operator_episode_gc.go:25` | `TestEpisodeOperatorStepLookupRefusals/query_error` | KILLED | KILLED |
| diff | `GC-ST-SCALAR` `host/store/operator_episode_gc.go:25` | `TestEpisodeHasUnresolvedOperatorSteps/reserved` | KILLED | KILLED |
| diff | `GC-ST-ACTION-JOIN` `host/store/operator_episode_gc.go:22` | `TestEpisodeHasUnresolvedOperatorSteps/operation_only` | KILLED | KILLED |
| diff | `GC-ST-EPISODE` `host/store/operator_episode_gc.go:23` | `TestEpisodeHasUnresolvedOperatorSteps/none` | KILLED | KILLED |
| diff | `GC-ST-OP-ALIAS` `host/store/operator_episode_gc.go:23` | `TestEpisodeHasUnresolvedOperatorSteps/operation_only` | KILLED | KILLED |
| diff | `GC-ST-RESERVED` `host/store/operator_episode_gc.go:24` | `TestEpisodeHasUnresolvedOperatorSteps/terminal_with_reservation` | KILLED | KILLED |
| diff | `GC-ST-TERMINAL` `host/store/operator_episode_gc.go:24` | `TestEpisodeHasUnresolvedOperatorSteps/observed_running` | KILLED | KILLED |
| diff | `GC-ST-OR` `host/store/operator_episode_gc.go:24` | `TestEpisodeHasUnresolvedOperatorSteps/terminal_with_reservation` | KILLED | KILLED |
| diff | `GC-ST-OPEN-TERMINAL-FILTER` `host/store/operator_episode_gc.go:24` | `TestEpisodeHasUnresolvedOperatorSteps/observed_failed` | KILLED | KILLED |
| diff | `GC-ST-PARAM-EPISODE` `host/store/operator_episode_gc.go:24` | `TestEpisodeHasUnresolvedOperatorSteps/reserved` | KILLED | KILLED |
| diff | `GC-D-WIRING` `host/daemon/auto_gc.go:86` | `TestEpisodeGCOperatorSteps/daemon_wiring` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-D-WIRING-BOOL` `host/daemon/auto_gc.go:86` | `TestEpisodeGCOperatorSteps/daemon_wiring` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-D-CLOCK` `host/daemon/auto_gc.go:83` | `TestEpisodeGCOperatorSteps/daemon_wiring` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-D-NIL` `host/daemon/auto_gc.go:158` | `TestEpisodeGCOperatorSteps/candidate` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-D-ERR` `host/daemon/auto_gc.go:161` | `TestEpisodeGCOperatorSteps/error` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-D-CANCELED` `host/daemon/auto_gc.go:161` | `TestEpisodeGCOperatorSteps/canceled` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-D-ERR-OR` `host/daemon/auto_gc.go:161` | `TestEpisodeGCOperatorSteps/timeout` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-D-TRUE` `host/daemon/auto_gc.go:162` | `TestEpisodeGCOperatorSteps/protected` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-D-FALSE` `host/daemon/auto_gc.go:164` | `TestEpisodeGCOperatorSteps/candidate` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-D-SKIP` `host/daemon/auto_gc.go:168` | `TestEpisodeGCOperatorSteps/protected` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-D-CONTINUE` `host/daemon/auto_gc.go:172` | `TestEpisodeGCOperatorSteps/protected` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-D-EPISODE` `host/daemon/auto_gc.go:159` | `TestEpisodeGCOperatorSteps/protected` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-C-EXPLICIT` `cmd/ailang-worldd/episode_gc.go:52` | `TestEpisodeGCFailSafeAndReclaims` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-C-HOME` `cmd/ailang-worldd/episode_gc.go:56` | `TestEpisodeGCFailSafeAndReclaims` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-C-SERVICE` `cmd/ailang-worldd/episode_gc.go:63` | `TestEpisodeGCFailSafeAndReclaims` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-C-MALFORMED` `cmd/ailang-worldd/episode_gc.go:65` | `TestEpisodeGCFailSafeAndReclaims` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-C-ABS-INPUT` `cmd/ailang-worldd/episode_gc.go:68` | `TestEpisodeGCFailSafeAndReclaims` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-C-PATH-ERRORS` `cmd/ailang-worldd/episode_gc.go:71` | `TestEpisodeGCFailSafeAndReclaims` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-C-SAME-ROOT` `cmd/ailang-worldd/episode_gc.go:71` | `TestEpisodeGCFailSafeAndReclaims` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-C-PHYSICAL-SCAN` `cmd/ailang-worldd/episode_gc.go:69` | `TestEpisodeGCFailSafeAndReclaims` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-C-PHYSICAL-RECORDED` `cmd/ailang-worldd/episode_gc.go:70` | `TestEpisodeGCFailSafeAndReclaims` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-C-WIRING` `cmd/ailang-worldd/episode_gc.go:114` | `TestEpisodeGCFailSafeAndReclaims` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-C-WIRING-BOOL` `cmd/ailang-worldd/episode_gc.go:114` | `TestEpisodeGCFailSafeAndReclaims` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| diff | `GC-C-FLAG` `cmd/ailang-worldd/episode_gc.go:106` | `TestEpisodeGCFailSafeAndReclaims` | KILLED | UNINFORMATIVE UNDER SANDBOX |
| judge-437 | `N01-ALLOWANCE-PATHS` `cmd/ailang-worldd/operator_files.go:62` | `TestOperatorFileAllowanceOrigin` | SURVIVED | UNINFORMATIVE UNDER SANDBOX |
| judge-437 | `N04-REOPEN-CMP` `cmd/ailang-worldd/operator_files.go:194` | `TestOperatorFilesReconcile/replaced_before_readback` | SURVIVED | UNINFORMATIVE UNDER SANDBOX |
| judge-437 | `N12-FRESH-DIR` `cmd/ailang-worldd/operator_files.go:367` | `TestOperatorFilesReconcile/all_files` | SURVIVED | UNINFORMATIVE UNDER SANDBOX |
| judge-437 | `N15-POST-TREE` `cmd/ailang-worldd/operator_files.go:396` | `TestOperatorFilesReconcile/replaced_after_readback` | SURVIVED | UNINFORMATIVE UNDER SANDBOX |
| judge-437 | `N22-STORE-HEADLESS` `host/store/operator_files.go:36` | `TestOperatorFileStoreReadPins/headless` | KILLED | KILLED |
| judge-437 | `N23-STORE-WTCOUNT` `host/store/operator_files.go:44` | `TestOperatorFileStoreReadPins/duplicate_worktree` | KILLED | KILLED |
| judge-437 | `N24-STORE-LOCKCHK` `host/store/operator_files.go:91` | `TestOperatorFileStoreReservePins/foreign_lock` | KILLED | KILLED |
| judge-437 | `N25-STORE-RESERVE-EMPTY` `host/store/operator_files.go:109` | `TestOperatorFileStoreReservePins/empty_identity` | KILLED | KILLED |
| judge-437 | `N27-STORE-CREATED` `host/store/operator_files.go:194` | `TestOperatorFileStoreReservePins/before_created` | KILLED | KILLED |
| judge-437 | `N28-STORE-APPROVEDBY` `host/store/operator_files.go:204` | `TestOperatorFileStoreReservePins/different_approver` | KILLED | KILLED |
| judge-437 | `N29-STORE-ALLOW` `host/store/operator_files.go:219` | `TestOperatorFileStoreReservePins/allowance` | KILLED | KILLED |
| judge-437 | `N31-STORE-ROWEXP` `host/store/operator_files.go:191` | `TestOperatorFileStoreReservePins/action_expired` | KILLED | KILLED |
| judge-437 | `N33-MINT-READ-STAGE` `host/store/operator_mint.go:246` | `TestOperatorMintStoreFilesFence/read` | KILLED | KILLED |
| judge-437 | `N34-MINT-ABORT-STAGE` `host/store/operator_mint.go:302` | `TestOperatorMintStoreFilesFence/abort` | KILLED | KILLED |
| judge-436 | `J436-ST-HEARTBEATWIN` `host/store/operator_worktree.go:59` | `TestOperatorWorktreeStoreReservePins/heartbeat_edge` | KILLED | KILLED |
| judge-436 | `J436-ST-NOEPISODESTART` `host/store/operator_worktree.go:56` | `TestOperatorWorktreeStoreReservePins/(existing_start\|existing_credential)` | KILLED | KILLED |
| judge-436 | `J436-ST-READPHASE` `host/store/operator_worktree.go:89` | `TestOperatorWorktreeStoreRecoveryPins/observed` | KILLED | KILLED |
| judge-436 | `J436-ST-READOTHERKIND` `host/store/operator_worktree.go:90` | `TestOperatorWorktreeStoreRecoveryPins/read_other_kind` | KILLED | KILLED |
| judge-436 | `J436-ST-DISPOTHERKIND` `host/store/operator_worktree.go:138` | `TestOperatorWorktreeStoreRecoveryPins/dispose_other_kind` | KILLED | KILLED |
| controls | `W-START-ONLY` `host/store/operator_worktree.go:56` | `TestOperatorWorktreeStoreReservePins/existing_start` | KILLED | KILLED |
| controls | `W-CREDENTIAL-ONLY` `host/store/operator_worktree.go:57` | `TestOperatorWorktreeStoreReservePins/existing_credential` | KILLED | KILLED |
| controls | `DIE-POSITIVE` `host/store/operator_worktree.go:31` | `TestOperatorFileStoreReadPins/valid` | KILLED | KILLED |
| controls | `LIVE-POSITIVE` `host/store/operator_worktree.go:31` | `TestOperatorFileStoreReadPins/valid` | SURVIVED | SURVIVED |

An initial standalone absolute-path error mutant survived. Planning folded the path errors into one fail-closed comparison and added the relative workspace control. The final `GC-C-PATH-ERRORS` and `GC-C-ABS-INPUT` substitutions both fail the existing named CLI test. The initial experiment is preserved at `$S/ledger-named-initial.jsonl`.

**Survivors and scope.** `N01-ALLOWANCE-PATHS`: Real allowance has no observed mint receipt and `TestOperatorFileAllowanceOrigin` observes the earlier `sql.ErrNoRows` refusal with exact evidence retained. `N04-REOPEN-CMP`: Optional window pin is measured in Appendix C. `N12-FRESH-DIR`: Optional window pin is measured in Appendix C. `N15-POST-TREE`: Optional window pin is measured in Appendix C. `LIVE-POSITIVE`: Live-positive harness control.

**Diff audit.** `$S/diff-audit.tsv` enumerates all 98 production changed lines. It labels 59 added lines with no independent mutation and 6 removed/replaced base lines. The components covered by the surviving aliases are reported above. Each SQL term and added decision has a literal substitution in the ledger. Help text, comments, formatting, option fields, context forwarding, cleanup calls and report construction are listed as audit lines where they have no separate mutant. The raw logs contain the parser controls.

### 4.1 Whole-package runs for the controller

````bash
S=$HOME/.ailang/state/world-iter438/plan-scratch
/bin/bash "$S/whole-mutants.sh" "$S/st-M3" "$S/controller-base"
````

The script uses `whole-mutant-ids.txt`, makes a fresh scratch product copy and mutates only that copy. The new base fail sets are `$S/controller-base/pkg-{cmd-ailang-worldd,host-daemon,host-store}.base-fails`. A historical base at `eaa202b1` is evidence only. The parser prints a failing control and rejects rc 1 with no attributable failure. A known flake drift is rerun in isolation and is documented alongside the comparison.

## 5. Files

````text
39	3	cmd/ailang-worldd/episode_gc.go
178	6	cmd/ailang-worldd/episode_gc_test.go
27	3	host/daemon/auto_gc.go
125	0	host/daemon/auto_gc_s5e4_test.go
1	0	host/store/credentials_test.go
26	0	host/store/operator_episode_gc.go
149	0	host/store/operator_episode_gc_test.go
52	0	host/store/operator_mint_s5e4_test.go
320	0	host/store/operator_s5e4_test.go
````

The only planner writes in the worktree are this plan and `sprint_w-world-address-bar-verbs-v4s5e4.json`. Everything else in §5 is shipped as patches. No stage adds a route, schema change, or new operator start caller.

## 6. Sprint gates

Run the following with the support files copied in §2. For a pristine archive the committed-blob reference must resolve to `684069f2` or a descendant that changes only the two plan files. For final, use the actual committed product reference for that read-only verification. The package legs unset GOFLAGS and record all outputs before comparing fail sets. `verify-scratch.py` relocates shell temp paths in a copy and supplies the committed-blob reference only to the tracked-binary subprocesses. Go test sources are those of the requested tree, and Git and AILANG are the real installed binaries. No PATH shims and no TMPDIR/GOTMPDIR override.

````bash gates.sh
#!/bin/bash
# gates.sh TREE LOGDIR base|M1|M2|M3 [BASE-LOGDIR] [READ-ONLY-GIT-REF-TREE]
set -uo pipefail
S=$(cd "$(dirname "$0")" && pwd -P)
T=$(cd "$1" && pwd -P)
mkdir -p "$2" || exit 2
L=$(cd "$2" && pwd -P)
M=$3
BASE_LOG=${4:-}
REF=${5:-}
case "$M" in base|M1|M2|M3) ;; *) exit 2;; esac
ORIGINAL=$HOME
export PATH=/opt/homebrew/bin:$PATH GIT_OPTIONAL_LOCKS=0
export AILANG_BIN="$ORIGINAL/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$ORIGINAL/.pinned-ailang-tools/v0.54.15/ailang"
export AILANG_STATE_DIR="$L/state" AILANG_CACHE_DIR="$L/cache"
export GOCACHE="$(go env GOCACHE)" GOMODCACHE="$(go env GOMODCACHE)" HOME="$L/home"
export XDG_CONFIG_HOME="$HOME/.config" XDG_CACHE_HOME="$HOME/.cache" XDG_DATA_HOME="$HOME/.local/share"
mkdir -p "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR" "$HOME"
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
unset WORLD_UPDATE_OPERATOR_ACTION_GOLDEN WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN
cd "$T" || exit 2
: > "$L/legs.txt"
FAILS=() UNKNOWN=() START=$(date +%s)
line() { printf '%s\n' "$*" | tee -a "$L/legs.txt"; }
skip() { line "LEG $1 SKIPPED 0s ($2)"; }
bounded() { python3 "$S/bounded.py" "$@"; }
denied() { grep -Eiq '(bind|listen tcp).*(: operation not permitted|: permission denied)' "$1"; }
leg() {
  local name=$1 rc=0 start; shift
  start=$(date +%s)
  "$@" > "$L/$name.txt" 2>&1 || rc=$?
  local result="rc=$rc $(( $(date +%s)-start ))s"
  if denied "$L/$name.txt" || [ "$rc" -eq 77 ]; then
    result="$result UNINFORMATIVE UNDER SANDBOX"
    UNKNOWN+=("$name")
  elif [ "$rc" -ne 0 ]; then
    FAILS+=("$name")
  fi
  line "LEG $name $result"
}
record() {
  local tag=$1 rc=0 parser=0; shift
  "$@" > "$L/$tag.raw" 2>&1 || rc=$?
  printf '%s\n' "$rc" > "$L/$tag.rc"
  python3 "$S/failset.py" "$L/$tag.raw" "$rc" > "$L/$tag.fails" || parser=$?
  if denied "$L/$tag.raw"; then grep -Ei '(bind|listen tcp).*(: operation not permitted|: permission denied)' "$L/$tag.raw"; return 77; fi
  if [ "$parser" -ne 0 ] || [ "$rc" -gt 1 ] ||
    grep -Eq '\[build failed\]|test timed out|BOUND EXCEEDED' "$L/$tag.raw"; then
    echo 'INSTRUMENT FAILURE'; return 2
  fi
  cp "$L/$tag.fails" "$L/$tag.attributable"
  cat "$L/$tag.fails"
}
inventory() {
  local pkg=$1 names=$2 name output rc=0
  output=$(mktemp "$L/inventory.XXXXXX") || return
  bounded 180 go test -v -count=1 -timeout 2m "./$pkg" -run "^($names)\$" > "$output" 2>&1 || rc=$?
  cat "$output"
  [ "$rc" -eq 0 ] || return "$rc"
  while IFS= read -r name; do
    grep -Fq -- "--- PASS: $name (" "$output" || return 2
  done < <(printf '%s\n' "$names" | tr '|' '\n')
}
fmt() {
  local path
  for path in host/store/credentials_test.go host/store/*s5e4*.go host/store/operator_episode_gc*.go host/daemon/auto_gc*.go cmd/ailang-worldd/episode_gc*.go; do
    [ ! -f "$path" ] || gofmt -l "$path" || return
  done
}
check_fmt() { local out; out=$(fmt) || return; printf '%s\n' "$out"; [ -z "$out" ]; }
{ uptime; ps -A -o pid,pcpu,comm -r | head -12; go version; command -v git; } > "$L/load.txt" 2>&1
leg gofmt check_fmt
leg build bounded 240 go build ./...
leg vet bounded 240 go vet ./host/store ./host/daemon ./cmd/ailang-worldd
leg linux-vet env GOOS=linux python3 "$S/bounded.py" 240 go vet ./host/store ./host/daemon ./cmd/ailang-worldd
leg inv-store inventory host/store 'TestProductionContextRoots|TestProductionGoSurface|TestStoreSurfaceTakesContext|TestCredentialsMethodsRequireADeadline'
leg inv-boundary inventory host/boundary 'TestWorldBoundaryDependencyAllowlist|TestBoundaryASTWriteGuard'
leg inv-cli inventory cmd/ailang-worldd 'TestCLICredentialWriterSites|TestEveryVerbHasHelp'
leg inv-verbs inventory host/verbs TestOperatorNeverDispatch
leg inv-daemon inventory host/daemon TestDaemonDependencyAllowlist
if [ "$M" = base ]; then skip accept 'new pins absent at base'
else leg accept env S5E4_PIN_HOME="$ORIGINAL" /bin/bash "$S/accept.sh" "$T" "$M" fast; fi
if [ "$M" = base ]; then
  leg gc-probe env HOME="$ORIGINAL" /bin/bash "$S/probe.sh" "$T"
elif [ "$M" = M3 ]; then
  leg gc-probe env HOME="$ORIGINAL" /bin/bash "$S/probe.sh" "$T" operator_gc_after_test.go
else skip gc-probe 'GC is unchanged; see base and M3'; fi
if [ "$M" = base ] || [ "$M" = M3 ]; then
  leg verify_ail bounded 510 python3 "$S/verify-scratch.py" "$T" "$L" verify_ail ""
  if [ -n "$REF" ]; then
    leg verify_go record verify_go bounded 1500 python3 "$S/verify-scratch.py" "$T" "$L" verify_go "$REF"
  else skip verify_go 'read-only committed-blob Git reference missing'; fi
else
  skip verify_ail 'AILANG tree byte-identical; base and M3 carry the verifier'
  skip verify_go 'pristine and final broad gate'
fi
leg callers /bin/bash "$S/callers.sh" "$T" "$L/callers.census" "$S/callers-base.txt"
for pkg in host/store host/daemon cmd/ailang-worldd; do
  name=${pkg//\//-}
  if { [ "$M" = M1 ] || [ "$M" = M2 ]; } && [ "$pkg" != host/store ]; then
    skip "pkg-$name" 'package files unchanged'
    skip "race-$name" 'package files unchanged'
    continue
  fi
  leg "pkg-$name" record "pkg-$name" bounded 510 go test -count=1 -timeout 8m "./$pkg"
  leg "race-$name" record "race-$name" bounded 510 go test -race -count=1 -timeout 8m "./$pkg"
done
for tag in $(printf '%s\n' pkg-host-store race-host-store pkg-host-daemon race-host-daemon pkg-cmd-ailang-worldd race-cmd-ailang-worldd verify_go); do
  if [ -f "$L/$tag.attributable" ]; then
    if [ -z "$BASE_LOG" ]; then cp "$L/$tag.attributable" "$L/$tag.base-fails"
    elif [ -f "$BASE_LOG/$tag.base-fails" ]; then
      leg "compare-$tag" diff -u "$BASE_LOG/$tag.base-fails" "$L/$tag.attributable"
    else skip "compare-$tag" 'attributable base recording unavailable'; fi
  else skip "compare-$tag" 'attributable stage recording unavailable'; fi
done
line "GATES_DONE $(( $(date +%s)-START ))s FAILS: ${FAILS[*]:-none} UNINFORMATIVE: ${UNKNOWN[*]:-none}"
tail -1 "$L/legs.txt" > "$L/ALL_DONE"
[ "${#FAILS[@]}" -eq 0 ] || exit 1
[ "${#UNKNOWN[@]}" -eq 0 ] || exit 77
````

The base probe uses the unchanged iteration 437 fixture in a scratch archive copy. The M3 probe uses the same signed publication entry fixture and adds real store/CLI wiring checks. The named verification saves its input/output receipts in the selected log directory.

| Fixture arm | Pristine 684069f2 | Compiled M3 |
|---|---|---|
| Real reserved before first publication | Injected 73 h, dry-run candidate `clean & merged` | At 73 h, live core skips as `non-terminal operator step`. Real CLI also skips. |
| Real published, action left running | Injected 73 h, dry-run candidate `clean & merged` | At 73 h, live core skips as `non-terminal operator step`. Real CLI also skips. |
| Same published evidence after terminal disposition | Not part of the unchanged base probe | Injected 73 h, dry-run candidate `clean & merged`. |

Both sides verify a clean Git status and preserved evidence. `$S/gate-base-final-script/gc-probe.txt` and `$S/gate-M3-final/gc-probe.txt` contain the exact isolated before/after. The CLI arm supplies an explicit fixture credential path and a zero age threshold to isolate store wiring from its wall clock. The core arms use the same 72-hour cutoff on both trees.

## 7. Executor directive notes

````bash accept.sh
#!/bin/bash
set -euo pipefail
W=$1; M=$2; PART=$3
case "$M" in M1|M2|M3) ;; *) exit 2;; esac
test "$PART" = fast
S=$(cd "$(dirname "$0")" && pwd -P)
PIN=${S5E4_PIN_HOME:-$HOME}
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN="$PIN/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$PIN/.pinned-ailang-tools/v0.54.15/ailang"
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
unset WORLD_UPDATE_OPERATOR_ACTION_GOLDEN WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN
run_dir=$(mktemp -d "$S/accept.XXXXXX")
export AILANG_STATE_DIR="$run_dir/state" AILANG_CACHE_DIR="$run_dir/cache"
export GOCACHE="$(go env GOCACHE)" GOMODCACHE="$(go env GOMODCACHE)" HOME="$run_dir/home"
mkdir -p "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR" "$HOME"
cd "$W"
run() {
  local pkg=$1 names=$2 name log rc=0
  log=$(mktemp "$run_dir/${pkg//\//-}.XXXXXX")
  python3 "$S/bounded.py" 180 go test -v -count=1 -timeout 2m "./$pkg" -run "^($names)\$" 2>&1 | tee "$log" || rc=$?
  if grep -Ei '(bind|listen tcp).*(: operation not permitted|: permission denied)' "$log"; then
    echo 'UNINFORMATIVE UNDER SANDBOX'
    exit 77
  fi
  test "$rc" -eq 0
  while IFS= read -r name; do
    grep -Fq -- "--- PASS: $name (" "$log"
  done < <(printf '%s\n' "$names" | tr '|' '\n')
}
pins='TestOperatorFileStoreReadPins|TestOperatorFileStoreReservePins|TestOperatorWorktreeStoreReservePins|TestOperatorWorktreeStoreRecoveryPins'
if [ "$M" != M1 ]; then pins="$pins|TestOperatorMintStoreFilesFence"; fi
if [ "$M" = M3 ]; then
  pins="$pins|TestEpisodeHasUnresolvedOperatorSteps|TestEpisodeOperatorStepLookupRefusals|TestCredentialsMethodsRequireADeadline|TestProductionContextRoots"
fi
run host/store "$pins"
if [ "$M" = M3 ]; then
  run host/daemon TestEpisodeGCOperatorSteps
  run cmd/ailang-worldd TestEpisodeGCFailSafeAndReclaims
fi
printf 'ACCEPT_OK %s %s\n' "$M" "$PART"
````

Snapshot control accepts an archive status fixture when explicitly requested for tests. For the executor, unset both override variables and snapshot from the actual worktree. Paths `docs`, `app` and `cmd` need not exist. Both plan records are excluded even after the controller commits them.

````bash snapshot.sh
#!/bin/bash
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=$1
M=$2
case "$M" in M1|M2|M3) ;; *) exit 2;; esac
B=684069f2118264d0797a0f7575b3c4d837c7fef3
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter438/snap}/$M
mkdir -p "$(dirname "$D")"
mkdir "$D"
D=$(cd "$D" && pwd -P)
cd "$W"
if [ -n "${SNAP_STATUS_LIST:-}" ]; then
  cp "$SNAP_STATUS_LIST" "$D/STATUS"
else
  {
    git diff --name-status --no-renames "$B" --
    git ls-files --others --exclude-standard | sed $'s/^/A\t/'
  } | awk -F '\t' '
    $2 != "design_docs/planned/w-world-address-bar-verbs-v4s5e4-sprint-plan.md" &&
    $2 != "design_docs/planned/sprint_w-world-address-bar-verbs-v4s5e4.json"
  ' | LC_ALL=C sort -u -t "$(printf '\t')" -k2 > "$D/STATUS"
fi
test -s "$D/STATUS"
awk -F '\t' '$1 != "D" {print $2}' "$D/STATUS" | while IFS= read -r p; do
  mkdir -p "$D/$(dirname "$p")"
  cp -p "$p" "$D/$p"
done
awk -F '\t' '$1 == "D" {print $2}' "$D/STATUS" > "$D/DELETED"
(
  cd "$D"
  while IFS="$(printf '\t')" read -r kind p; do
    if [ "$kind" != D ]; then shasum -a 256 "$p"; fi
  done < STATUS > MANIFEST.sha256
  shasum -a 256 -c MANIFEST.sha256
)
wc -l "$D"/{STATUS,DELETED,MANIFEST.sha256}
````

Controller commands for each milestone (commit the two plan records before the product snapshot script removes unsnapshotted files). The controller runs these commands.

````bash
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter438-row243-v4s5e4
STATE=$HOME/.ailang/state/world-iter438
S=$STATE/plan-scratch
T=$STATE/executor
B=684069f2118264d0797a0f7575b3c4d837c7fef3

# Before product commits, record the base outside the sandbox. Its product tree
# must still match B; baseline tests run strictly in the pristine archive.
git -C "$W" diff --quiet "$B" HEAD -- cmd host docs app || exit 2
base_rc=0
/bin/bash "$T/gates.sh" "$S/base" "$S/controller-base" base "" "$W" || base_rc=$?
test -f "$S/controller-base/ALL_DONE" || exit 2
printf 'Base gate rc=%s\n' "$base_rc"
# Record every leg before attributing failures. The controller reads ALL_DONE.

printf '%s\n' 'test(store): pin operator file and worktree admission' > "$T/M1.msg"
printf '%s\n' 'test(store): independently pin both mint recovery fences' > "$T/M2.msg"
printf '%s\n' 'fix(gc): preserve worktrees with unresolved operator effects' > "$T/M3.msg"
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
for M in M1 M2 M3; do
  test -f "$STATE/snap/$M/STATUS" || exit 2
  test_cmd="AILANG_BIN=\"$AILANG_BIN\" WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=\"$WORLD_EXEC_NODE\" WORLD_TOOL_AILANG_BIN=\"$WORLD_TOOL_AILANG_BIN\" /bin/bash \"$T/accept.sh\" \"$W\" $M fast"
  /bin/bash "$STATE/commit_ms.sh" "$M" "$T/$M.msg" "$test_cmd" || exit
done
git -C "$W" diff -M --numstat "$B" HEAD -- cmd host docs app
git -C "$W" diff -M --numstat "$B" HEAD
final_rc=0
/bin/bash "$T/gates.sh" "$S/st-M3" "$S/gate-controller-final" M3 "$S/controller-base" "$W" || final_rc=$?
test -f "$S/gate-controller-final/ALL_DONE" || exit 2
printf 'Final gate rc=%s\n' "$final_rc"
/bin/bash "$S/whole-mutants.sh" "$S/st-M3" "$S/controller-base"
````

No commit-time test command runs a whole CLI or daemon package. The full package gates and `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` use the controller's socket-capable baseline and compare fail sets. Race uses `go test -race -count=1 -timeout 8m` per package on that same rig. Keep the baseline and final outputs before evaluating any rc. Native Ubuntu CI remains the Linux execution gate.

The exact Appendix P fences were extracted from this plan and replayed under `$S/rt-0zhgJj`. M1, M2 and M3 returned `diff -r` rc 0 for the full trees and each of `cmd`, `host` and `docs`. `app` was absent in both trees. The rehearsed snapshot used the main executor path with the sprint worktree's Git metadata read-only and the replay tree as `GIT_WORK_TREE`. Cumulative STATUS/deletion/hash counts were `1/0/1`, `2/0/2` and `9/0/9`. No new patch reject or backup files appeared. Extracted shell fences matched the tested copies and passed Bash syntax checks. Receipt `$S/roundtrip-refined.log`.

## 8. Risks

- The eleven row-271 CLI real-binary baseline rows are `TestCLIAgainstSeToolsDaemon`, `TestCLINoGenesisIsNoWorldHead`, `TestDemoSeedApprovalPacket`, `TestDemoSeedEndToEnd`, `TestDemoSeedManifestReferences`, `TestDemoSeedProofStage`, `TestDemoSeedReadTestCheck`, `TestDemoSeedRejectsToolReportRefusals`, `TestDemoSeedRosterKindMismatchRefuses`, `TestDemoSeedWhyNamedLinks`, and `TestMarketingStillRoutes`. The daemon baseline rows are `TestSeToolsDescriptorLiveExactEight`, `TestSeToolsMCPVisibilityFollowsGrants`, `TestSeToolsA2AEndToEnd`. Treat `TestWorldContextAndStatusline`, `host/broker TestExecProtectedCheckTimeout` and `TestSeedRestartIgnoresASurvivingOldCopy`, `host/daemon TestSeToolsRunBenchmarkTasksEndToEnd`, and `host/verifygate TestMissionConfigDefaultFlowIntegration` as load-sensitive baseline candidates. See exact historical test names and new fail-set paths in §1.5.
- Row 267 provides a 480 s bound for both slow-runner packages. The user-supplied CLI runs reached 477.6 and 469.0 s. Daemon reached 480.0x s. Named race costs and full package status are separate measurements. No socket or launched daemon is needed for the new named tests.
- A successful negative lookup applies only to the observed step ledger. Future activation must preserve a worktree exclusion boundary across effects. The retained GC probe is a pre-caller test.

## 9. Acceptance criteria

1. Apply the fenced diffs in order with `/usr/bin/patch -p1 -E`, and match each cumulative numstat and snapshot manifest.
2. Run the M1 and M2 named store rows, retain every refusal's exact state, and record every original judge edit. Accept N01 only with its measured missing-mint mask.
3. Prove GC keeps both reserved and observed undisposed worktrees and preserves Git evidence on lookup failures. Prove terminal or empty ledgers can still be reclaimed. Run the real GC probe before and after with the approved temporary fixture.
4. Run the named inventories, native build/vet, Linux vet, pinned AILANG gate and production caller census. Send bind-denied commands through the controller recipe with the fresh base fail sets.

## Appendix C — after this sprint (measured)

- **M4 parked.** The exact N04, N12 and N15 substitutions survive their inherited named selectors on final M3. The scratch prototype under `$S/opt-M4` shows that N04 (`cmd/ailang-worldd/operator_files.go:194`) can use a 9-changed-line production seam (+7/−2) that supplies the existing directory opener. N12 (:367) and N15 (:396) can use the existing `before ep.task.md` and `before observe` checkpoints with zero production seam lines. The combined prototype test adds 73 lines, for a total footprint of **82**.
- The prototype's `TestOptionalM4PinCost/{N04,N12,N15}` each rejects its exact judge edit in the scratch copy. The paired command `go test -race -v -count=1 -timeout 2m ./cmd/ailang-worldd -run '^TestOptionalM4PinCost$'` took 1.632 s of package time on base (zero discovered tests) and 6.031 s on the prototype, or **+4.399 s**. M3 uses the single permitted CLI test slot and remains under the row-267 time allowance for that test. M4 has zero staged lines. See `$S/m4-measured.log` and its raw named logs.

- **S5f/S6 paired caller and drain.** No caller in this sprint. S5e3 recorded a prior 72-line process FD harness and raw design counts of S6 890, S8 860 and S12 700. Full S5f/S6 implementation is **UNMEASURED** here. The four process pins from S5e3 are still carried: `TestOperatorFirstManagedEpisodeProcess`, `TestOperatorOldBuildTerminalDrain`, `TestOperatorCurrentOwnerProcessRecovery`, and `TestOperatorServiceUpgrade/active_child`. The older `E-RES-EPISODE-CONFLICT`, `E-RES-REREAD`, `J-OWN-UID`, `Z-ABSENT-ANY` and `Z-WAIT-READBUDGET` are outside this requested survivor set. New patch lines for these deferrals: zero.
- **Row 271 and load-sensitive tests.** Their outcomes are controller baseline evidence. Sprint patch lines aimed at those tests: zero.

## Appendix P — compiled patches

### Patch M1

````diff M1
diff --git a/host/store/operator_s5e4_test.go b/host/store/operator_s5e4_test.go
new file mode 100644
index 00000000..6261c8e2
--- /dev/null
+++ b/host/store/operator_s5e4_test.go
@@ -0,0 +1,320 @@
+package store
+
+import (
+	"context"
+	"database/sql"
+	"errors"
+	"path/filepath"
+	"reflect"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+)
+
+const pinNow int64 = 1590
+
+type operatorPinFixture struct {
+	c     *Credentials
+	ctx   context.Context
+	in    OperatorMintInput
+	proof *approval.OperatorMintProof
+}
+
+func newOperatorPinFixture(t *testing.T, allowance bool) operatorPinFixture {
+	t.Helper()
+	c, ctx := openTestCredentials(t), testCtx(t)
+	cas, runner := operatorClaimFixture(t, c)
+	row := getOperator(t, c, cas.ID, pinNow)
+	a, err := approval.DecodeOperatorAction(row.Action)
+	if err != nil {
+		t.Fatal(err)
+	}
+	dir := filepath.Dir(c.lock.dbPath)
+	in := OperatorMintInput{
+		Attempt:      OperatorMintAttempt{CAS: cas, Runner: runner, Catalogue: 1},
+		TokenPath:    filepath.Join(dir, a.Episode+".session"),
+		Profile:      approval.ProfileRef{Project: "world", Digest: hashref.SumSHA256([]byte("pin profile")).String()},
+		Grants:       []approval.Grant{{Effect: "exec", Scope: "repo", Budget: 4}},
+		CredentialID: hex64("operator pin"),
+		Start: EpisodeStart{Schema: "world/episode-start/v1", Episode: a.Episode, Agent: "codex",
+			Preset: StartPreset{Name: a.Options.Preset}, Repo: dir, TTLS: 600,
+			Task: &StartTask{Path: filepath.Join(dir, a.Episode+".task.md"), SHA256: hashref.SumSHA256([]byte(a.Task)).String()}},
+	}
+	if allowance {
+		// Give the fixture a signed, internally consistent mixed-origin mint.
+		// A normal allowance draw alone lacks a matching operator mint receipt.
+		in.Start.Allowance = &StartAllowance{AllowanceID: "sa-0123456789abcdef", DrawID: "ad-0123456789abcdef", ApprovedBy: "tty", DayKey: "1970-01-01", DrawNumber: 1}
+	}
+	a.Resolution, err = in.Resolution()
+	if err != nil {
+		t.Fatal(err)
+	}
+	proof := signOperatorPinAction(t, c, row, a)
+	in.Attempt.CAS.Digest, err = approval.OperatorSignatureDigest(a)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if err := c.ClaimOperatorAction(ctx, in.Attempt.CAS, runner, pinNow); err != nil {
+		t.Fatal(err)
+	}
+	in.Attempt.CAS.Generation++
+	return operatorPinFixture{c, ctx, in, proof}
+}
+
+func signOperatorPinAction(t *testing.T, c *Credentials, row OperatorActionRow, a approval.OperatorAction) *approval.OperatorMintProof {
+	t.Helper()
+	body, err := approval.CanonicalOperatorAction(a)
+	if err != nil {
+		t.Fatal(err)
+	}
+	digest, err := approval.OperatorSignatureDigest(a)
+	if err != nil {
+		t.Fatal(err)
+	}
+	challenge := [32]byte(row.Challenge)
+	message, err := approval.OperatorActionMessage(digest, challenge, row.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	signer := approvaltest.NewSigner(42)
+	ap := approval.Approver{Name: "pin", Format: approval.FormatSSH, Key: signer.PublicKey(), EnrolledAt: 900}
+	_, enrolled, err := c.Approver(testCtx(t), ap.Name)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if !enrolled {
+		if err := c.EnrolApprover(testCtx(t), ApproverRow{Name: ap.Name, Format: ap.Format, PublicKey: ap.Key.Blob, EnrolledAt: ap.EnrolledAt}); err != nil {
+			t.Fatal(err)
+		}
+	}
+	artifact := signer.Sign(approval.OperatorActionNamespace, message)
+	signed, err := approval.CheckOperatorAction(approval.OperatorCheckInput{Action: a, Challenge: challenge,
+		Expiry: row.ExpiresAt, Now: pinNow, Artifact: artifact, Format: ap.Format, Approver: ap})
+	if err != nil {
+		t.Fatal(err)
+	}
+	operatorExec(t, c, "UPDATE operator_actions SET action=?, digest=?, artifact=?, approved_by=? WHERE id=?",
+		body, digest[:], artifact, signed.ApprovedBy(), row.ID)
+	return signed.MintProof()
+}
+
+func (f operatorPinFixture) reserveWorktree(t *testing.T) OperatorStepKey {
+	t.Helper()
+	key, err := f.c.ReserveOperatorWorktree(f.ctx, f.in.Attempt, "fixture-worktree", f.proof, func() int64 { return pinNow })
+	if err != nil {
+		t.Fatal(err)
+	}
+	return key
+}
+
+func (f operatorPinFixture) mint(t *testing.T) {
+	t.Helper()
+	worktree := f.reserveWorktree(t)
+	if err := f.c.ObserveOperatorStep(f.ctx, worktree, []byte("created")); err != nil {
+		t.Fatal(err)
+	}
+	in := f.in
+	in.Attempt.CAS, in.Attempt.StepVersion = worktree.Action, 1
+	key := OperatorStepKey{in.Attempt.CAS, worktree.Owner, "mint", worktree.Action.ID, 1}
+	if err := f.c.ReserveOperatorStep(f.ctx, key, pinNow); err != nil {
+		t.Fatal(err)
+	}
+	if err := f.c.CommitOperatorMint(f.ctx, in, f.proof, func() int64 { return pinNow }); err != nil {
+		t.Fatal(err)
+	}
+}
+
+type operatorPinState struct {
+	action   OperatorActionRow
+	steps    string
+	sessions []SessionRow
+	start    *EpisodeStart
+}
+
+func (f operatorPinFixture) snapshot(t *testing.T) operatorPinState {
+	t.Helper()
+	id := f.in.Attempt.CAS.ID
+	s := operatorPinState{action: getOperator(t, f.c, id, pinNow)}
+	err := f.c.db.QueryRowContext(f.ctx, `SELECT json_group_array(json_array(kind, identity, phase, version, owner, created_at,
+		CASE WHEN receipt IS NULL THEN NULL ELSE hex(receipt) END))
+		FROM (SELECT * FROM operator_action_steps WHERE action_id=? ORDER BY kind, identity)`, id).Scan(&s.steps)
+	if err != nil {
+		t.Fatal(err)
+	}
+	s.sessions, err = f.c.ListSessions(f.ctx, f.in.Start.Episode, 20)
+	if err != nil {
+		t.Fatal(err)
+	}
+	s.start, err = f.c.EpisodeStart(f.ctx, f.in.Start.Episode)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return s
+}
+
+func (f operatorPinFixture) assertRefused(t *testing.T, before operatorPinState, got, want error) {
+	t.Helper()
+	if after := f.snapshot(t); !errors.Is(got, want) || !reflect.DeepEqual(before, after) {
+		t.Fatalf("refusal = %v, want %v; durable rows unchanged = %t", got, want, reflect.DeepEqual(before, after))
+	}
+}
+
+func TestOperatorFileStoreReadPins(t *testing.T) {
+	for _, name := range []string{"valid", "headless", "duplicate_worktree"} {
+		t.Run(name, func(t *testing.T) {
+			f := newOperatorPinFixture(t, false)
+			f.mint(t)
+			id := f.in.Attempt.CAS.ID
+			var want error
+			switch name {
+			case "headless":
+				row := getOperator(t, f.c, id, pinNow)
+				a, err := approval.DecodeOperatorAction(row.Action)
+				if err != nil {
+					t.Fatal(err)
+				}
+				a.Options.Headless = false
+				a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
+				signOperatorPinAction(t, f.c, row, a)
+				want = ErrOperatorConflict
+			case "duplicate_worktree":
+				row := getOperator(t, f.c, id, pinNow)
+				key := OperatorStepKey{OperatorCAS{id, [32]byte(row.Digest), row.Generation}, row.Owner, "worktree", "duplicate", 1}
+				if err := f.c.ReserveOperatorStep(f.ctx, key, pinNow); err != nil {
+					t.Fatal(err)
+				}
+				want = sql.ErrNoRows
+			}
+			before := f.snapshot(t)
+			got, _, err := f.c.ReadOperatorFiles(f.ctx, id)
+			f.assertRefused(t, before, err, want)
+			if want == nil && (got.Worktree != "fixture-worktree" || got.Phase != "" || got.Mint.CredentialID != f.in.CredentialID) {
+				t.Fatalf("positive read = %+v", got)
+			}
+		})
+	}
+}
+
+func TestOperatorFileStoreReservePins(t *testing.T) {
+	for _, name := range []string{"valid", "foreign_lock", "empty_identity", "before_created", "different_approver", "allowance", "action_expired"} {
+		t.Run(name, func(t *testing.T) {
+			f := newOperatorPinFixture(t, name == "allowance")
+			f.mint(t)
+			id := f.in.Attempt.CAS.ID
+			state, action, err := f.c.ReadOperatorFiles(f.ctx, id)
+			if err != nil {
+				t.Fatal(err)
+			}
+			value, now, want := "files-identity", pinNow, ErrOperatorConflict
+			switch name {
+			case "foreign_lock":
+				dbPath := f.c.lock.dbPath
+				f.c.lock.dbPath = filepath.Join(filepath.Dir(dbPath), "other.db")
+				defer func() { f.c.lock.dbPath = dbPath }()
+			case "empty_identity":
+				value, want = "", ErrOperatorInput
+			case "different_approver":
+				operatorExec(t, f.c, "UPDATE operator_actions SET approved_by='ssh:other'")
+			case "action_expired":
+				now, want = getOperator(t, f.c, id, pinNow).ExpiresAt, ErrOperatorExpired
+				if now >= state.Session.ExpiresAt {
+					t.Fatal("fixture requires a live credential at action expiry")
+				}
+			}
+			before := f.snapshot(t)
+			if name == "before_created" {
+				// The reservation SQL also rejects a future created_at. Pin
+				// the earlier admission independently of that second guard.
+				row := before.action
+				row.CreatedAt = now + 1
+				tx, err := f.c.db.BeginTx(f.ctx, &sql.TxOptions{ReadOnly: true})
+				if err != nil {
+					t.Fatal(err)
+				}
+				defer tx.Rollback()
+				err = admitOperatorFiles(f.ctx, tx, row, state, action, f.proof, now)
+				if closeErr := tx.Rollback(); closeErr != nil {
+					t.Fatal(closeErr)
+				}
+				f.assertRefused(t, before, err, want)
+				return
+			}
+			err = f.c.RecordOperatorFiles(f.ctx, state, value, f.proof, func() int64 { return now })
+			if name != "valid" {
+				f.assertRefused(t, before, err, want)
+				return
+			}
+			if err != nil {
+				t.Fatal(err)
+			}
+			key := state.Step
+			key.Identity, key.Version = value, 1
+			assertReservedOperatorStep(t, f.c, key)
+		})
+	}
+}
+
+func TestOperatorWorktreeStoreReservePins(t *testing.T) {
+	for _, name := range []string{"heartbeat_edge", "existing_start", "existing_credential"} {
+		t.Run(name, func(t *testing.T) {
+			f := newOperatorPinFixture(t, false)
+			switch name {
+			case "heartbeat_edge":
+				operatorExec(t, f.c, "UPDATE operator_runner_state SET seen_at=?", pinNow-10)
+			case "existing_start":
+				tx, err := f.c.db.BeginTx(f.ctx, nil)
+				if err != nil {
+					t.Fatal(err)
+				}
+				defer tx.Rollback()
+				if _, err := insertEpisodeStart(f.ctx, tx, f.in.Start.Episode, f.in.Start); err != nil {
+					t.Fatal(err)
+				}
+				if err := tx.Commit(); err != nil {
+					t.Fatal(err)
+				}
+			case "existing_credential":
+				row := SessionRow{CredentialID: f.in.CredentialID, EpisodeID: f.in.Start.Episode, GrantsJSON: "[]",
+					CreatedAt: pinNow, ExpiresAt: pinNow + 600}
+				if err := f.c.MintSession(f.ctx, row); err != nil {
+					t.Fatal(err)
+				}
+			}
+			before := f.snapshot(t)
+			key, err := f.c.ReserveOperatorWorktree(f.ctx, f.in.Attempt, "fixture-worktree", f.proof, func() int64 { return pinNow })
+			f.assertRefused(t, before, err, ErrOperatorConflict)
+			if key != (OperatorStepKey{}) {
+				t.Fatalf("refusal returned worktree key %+v", key)
+			}
+		})
+	}
+}
+
+func TestOperatorWorktreeStoreRecoveryPins(t *testing.T) {
+	for _, name := range []string{"observed", "read_other_kind", "dispose_other_kind"} {
+		t.Run(name, func(t *testing.T) {
+			f := newOperatorPinFixture(t, false)
+			key := f.reserveWorktree(t)
+			if name == "observed" {
+				if err := f.c.ObserveOperatorStep(f.ctx, key, []byte("created")); err != nil {
+					t.Fatal(err)
+				}
+			} else {
+				other := OperatorStepKey{key.Action, key.Owner, "mint", key.Action.ID, 1}
+				if err := f.c.ReserveOperatorStep(f.ctx, other, pinNow); err != nil {
+					t.Fatal(err)
+				}
+			}
+			before := f.snapshot(t)
+			if name == "dispose_other_kind" {
+				err := f.c.DisposeOperatorWorktree(f.ctx, key, OperatorWorktreeRetained, pinNow)
+				f.assertRefused(t, before, err, ErrOperatorConflict)
+				return
+			}
+			_, _, err := f.c.ReadOperatorWorktree(f.ctx, key.Action.ID)
+			f.assertRefused(t, before, err, ErrOperatorMissing)
+		})
+	}
+}
````

### Patch M2

````diff M2
diff --git a/host/store/operator_mint_s5e4_test.go b/host/store/operator_mint_s5e4_test.go
new file mode 100644
index 00000000..085b9259
--- /dev/null
+++ b/host/store/operator_mint_s5e4_test.go
@@ -0,0 +1,52 @@
+package store
+
+import (
+	"database/sql"
+	"testing"
+)
+
+func TestOperatorMintStoreFilesFence(t *testing.T) {
+	for _, name := range []string{"read", "abort", "abort_without_files"} {
+		t.Run(name, func(t *testing.T) {
+			f := newOperatorPinFixture(t, false)
+			f.mint(t)
+			id := f.in.Attempt.CAS.ID
+			reviewed, err := f.c.ReadOperatorMint(f.ctx, id)
+			if err != nil || reviewed.CredentialID != f.in.CredentialID {
+				t.Fatalf("read before later reservation: %v", err)
+			}
+			if name == "abort_without_files" {
+				if err := f.c.AbortOperatorMint(f.ctx, reviewed, pinNow); err != nil {
+					t.Fatal(err)
+				}
+				after := f.snapshot(t)
+				if after.action.State != "failed" || len(after.sessions) != 0 || after.start == nil {
+					t.Fatalf("control failed to revoke and stop while retaining the start: %s, %d sessions, start present %t",
+						after.action.State, len(after.sessions), after.start != nil)
+				}
+				return
+			}
+			files, _, err := f.c.ReadOperatorFiles(f.ctx, id)
+			if err != nil {
+				t.Fatal(err)
+			}
+			if err := f.c.RecordOperatorFiles(f.ctx, files, "files-identity", f.proof, func() int64 { return pinNow }); err != nil {
+				t.Fatal(err)
+			}
+			key := files.Step
+			key.Identity, key.Version = "files-identity", 1
+			assertReservedOperatorStep(t, f.c, key)
+			before := f.snapshot(t)
+			if name == "abort" {
+				err := f.c.AbortOperatorMint(f.ctx, reviewed, pinNow)
+				f.assertRefused(t, before, err, sql.ErrNoRows)
+				return
+			}
+			got, err := f.c.ReadOperatorMint(f.ctx, id)
+			f.assertRefused(t, before, err, sql.ErrNoRows)
+			if got != (OperatorMintReceipt{}) {
+				t.Fatalf("refusal returned mint evidence: %+v", got)
+			}
+		})
+	}
+}
````

### Patch M3

````diff M3
diff --git a/cmd/ailang-worldd/episode_gc.go b/cmd/ailang-worldd/episode_gc.go
index a1ac2476..8bd6b967 100644
--- a/cmd/ailang-worldd/episode_gc.go
+++ b/cmd/ailang-worldd/episode_gc.go
@@ -12,18 +12,21 @@ import (
 	"time"
 
 	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/store"
 )
 
 const episodeGCHelp = `usage: ailang-worldd episode gc [--workspace-root D] [--exec-cache-dir D]
-    [--older-than 72h] [--dry-run] [--json]
+    [--credentials-db PATH] [--older-than 72h] [--dry-run] [--json]
 
 Safely reclaim disk space from finished/stale episode worktrees and orphaned
 exec caches. Never removes dirty worktrees (uncommitted or untracked files)
-or worktrees with unmerged commits. Also runs automatically inside ailang-worldd
-serve (2m after boot, then every 6h).
+or worktrees with unmerged commits. Worktrees with unresolved operator steps
+or unavailable step evidence are skipped. Also runs automatically inside
+ailang-worldd serve (2m after boot, then every 6h).
 
   --workspace-root D   episode worktrees root (default ~/.ailang/world-ws)
   --exec-cache-dir D   per-episode exec caches (default ~/.ailang/world/exec-cache)
+  --credentials-db PATH   credentials for this workspace (matching home service or default workspace inferred)
   --older-than DUR     minimum age before a clean merged worktree is removed (default 72h)
   --dry-run            report what would be removed without deleting anything
   --json               emit JSON report
@@ -43,12 +46,41 @@ func defaultExecCacheDir() string {
 	return ""
 }
 
+// episodeGCCredentialsPath uses a store known to belong to the scanned
+// workspace. An explicit path is required for other workspace layouts.
+func episodeGCCredentialsPath(root, explicit string) string {
+	if explicit != "" {
+		return explicit
+	}
+	home, err := os.UserHomeDir()
+	if err != nil {
+		return ""
+	}
+	def := defaultWorldPaths()
+	workspace, credentials := def.workspaceRoot, filepath.Join(filepath.Dir(def.db), store.CredentialsFileName)
+	record, err := readServiceRecord(homeServicePaths(home))
+	switch {
+	case err == nil:
+		workspace, credentials = record.Workspace, record.Credentials
+	case !errors.Is(err, os.ErrNotExist):
+		return ""
+	}
+	root, absErr := filepath.Abs(root)
+	root, rootErr := filepath.EvalSymlinks(root)
+	workspace, workspaceErr := filepath.EvalSymlinks(workspace)
+	if errors.Join(absErr, rootErr, workspaceErr) != nil || root != workspace {
+		return ""
+	}
+	return credentials
+}
+
 func runEpisodeGC(args []string, stdout, stderr io.Writer) int {
 	fs := flag.NewFlagSet("ailang-worldd episode gc", flag.ContinueOnError)
 	fs.SetOutput(stderr)
 	fs.Usage = func() {}
 	wsRoot := fs.String("workspace-root", defaultWorkspaceRootDir(), "episode worktrees directory")
 	cacheDir := fs.String("exec-cache-dir", defaultExecCacheDir(), "per-episode exec cache directory")
+	creds := fs.String("credentials-db", "", "credentials store for the scanned workspace")
 	olderThan := fs.Duration("older-than", 72*time.Hour, "minimum age before removing clean merged worktrees")
 	dryRun := fs.Bool("dry-run", false, "preview without deleting")
 	asJSON := fs.Bool("json", false, "emit JSON")
@@ -71,12 +103,16 @@ func runEpisodeGC(args []string, stdout, stderr io.Writer) int {
 
 	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
 	defer cancel()
+	credentials := episodeGCCredentialsPath(*wsRoot, *creds)
 	rep := daemon.RunEpisodeGC(ctx, daemon.EpisodeGCOptions{
 		WorkspaceRoot: *wsRoot,
 		ExecCacheDir:  *cacheDir,
 		OlderThan:     *olderThan,
 		DryRun:        *dryRun,
 		Now:           time.Now(),
+		HasOpenOperatorStep: func(ctx context.Context, episode string) (bool, error) {
+			return store.EpisodeHasUnresolvedOperatorSteps(ctx, credentials, episode)
+		},
 	})
 	if *asJSON {
 		enc := json.NewEncoder(stdout)
diff --git a/cmd/ailang-worldd/episode_gc_test.go b/cmd/ailang-worldd/episode_gc_test.go
index c88dfdf9..0c5b18b7 100644
--- a/cmd/ailang-worldd/episode_gc_test.go
+++ b/cmd/ailang-worldd/episode_gc_test.go
@@ -2,14 +2,18 @@ package main
 
 import (
 	"bytes"
+	"database/sql"
 	"encoding/json"
+	"errors"
 	"os"
 	"os/exec"
 	"path/filepath"
 	"testing"
 	"time"
 
+	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/store"
 )
 
 func runGit(t *testing.T, dir string, args ...string) {
@@ -26,6 +30,8 @@ func runGit(t *testing.T, dir string, args ...string) {
 
 func TestEpisodeGCFailSafeAndReclaims(t *testing.T) {
 	base := t.TempDir()
+	home := filepath.Join(base, "home")
+	t.Setenv("HOME", home)
 	repo := filepath.Join(base, "repo")
 	wsRoot := filepath.Join(base, "world-ws")
 	cacheDir := filepath.Join(base, "exec-cache")
@@ -64,7 +70,7 @@ func TestEpisodeGCFailSafeAndReclaims(t *testing.T) {
 	runGit(t, unmergedWT, "commit", "-m", "unmerged work")
 
 	// Backdate all 3 worktrees by 100h so they pass the 72h threshold.
-	oldTime := time.Now().Add(-100 * time.Hour)
+	oldTime := time.Unix(1_000_000_000, 0)
 	for _, d := range []string{cleanWT, dirtyWT, unmergedWT} {
 		if err := os.Chtimes(d, oldTime, oldTime); err != nil {
 			t.Fatal(err)
@@ -88,7 +94,173 @@ func TestEpisodeGCFailSafeAndReclaims(t *testing.T) {
 		}
 	}
 
-	// Dry run leaves everything intact.
+	// An unrelated default store must never clear a custom workspace.
+	defaultCreds := filepath.Join(home, ".ailang", "world", store.CredentialsFileName)
+	if err := os.MkdirAll(filepath.Dir(defaultCreds), 0700); err != nil {
+		t.Fatal(err)
+	}
+	if _, err := store.UpgradeCredentials(defaultCreds); err != nil {
+		t.Fatal(err)
+	}
+	creds := filepath.Join(base, "custom.db")
+	c, err := store.OpenCredentials(creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer c.Close()
+	ctx := boundedTestContext(t)
+	a := approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: "oa-00000000000000000000000000000001",
+		Kind: "start", Episode: "ep-clean", Store: creds, Task: "fixture", Review: "fixture",
+		Options:    approval.OperatorOptions{Preset: "sol", Headless: true},
+		Resolution: []approval.OperatorResolved{{Name: "agent", Value: "codex"}}}
+	a.Store, err = store.CanonicalPath(creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
+	if _, err := c.RequestOperatorAction(ctx, store.OperatorRequest{Action: a, OperationID: "00000000000000000000000000000001",
+		CreatedAt: 1000, ExpiresAt: 1600}, 1000); err != nil {
+		t.Fatal(err)
+	}
+	// A reserved effect can survive its parent action. Use the real migrated
+	// schema to model that case without claiming or running an action.
+	raw, err := sql.Open("sqlite", creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer raw.Close()
+	if _, err := raw.ExecContext(ctx, `UPDATE operator_actions SET state='failed', terminal_at=1001 WHERE id=?;
+		INSERT INTO operator_action_steps (kind, identity, action_id, phase, version, owner, created_at)
+		VALUES ('worktree','fixture-wt',?,'reserved',1,'fixture',1001)`, a.ID, a.ID); err != nil {
+		t.Fatal(err)
+	}
+	gitState := func(args ...string) []byte {
+		t.Helper()
+		cmd := exec.CommandContext(ctx, "git", args...)
+		cmd.Dir = cleanWT
+		got, err := cmd.CombinedOutput()
+		if err != nil {
+			t.Fatalf("git %v: %v: %s", args, err, got)
+		}
+		return got
+	}
+	head, status := gitState("rev-parse", "HEAD"), gitState("status", "--porcelain=v1")
+	physical := func(path string) string {
+		t.Helper()
+		got, err := filepath.EvalSymlinks(path)
+		if err != nil {
+			t.Fatal(err)
+		}
+		return got
+	}
+	record := homeServicePaths(home)
+	if err := os.MkdirAll(filepath.Dir(record.Record), 0700); err != nil {
+		t.Fatal(err)
+	}
+	setRecord := func(root string) {
+		t.Helper()
+		data, err := json.Marshal(serviceRecord{Schema: "world/service/v1", Label: worldServiceLabel,
+			Address: "127.0.0.1:7644", Plist: record.Plist, Binary: record.Binary, Credentials: creds, Workspace: root})
+		if err != nil {
+			t.Fatal(err)
+		}
+		if err := os.WriteFile(record.Record, data, 0600); err != nil {
+			t.Fatal(err)
+		}
+	}
+	assertSkip := func(reason string, extra ...string) {
+		t.Helper()
+		var out, errw bytes.Buffer
+		args := []string{"episode", "gc", "--workspace-root", wsRoot, "--exec-cache-dir", filepath.Join(base, "absent-cache"), "--json"}
+		code := run(append(args, extra...), &out, &errw)
+		var rep daemon.EpisodeGCReport
+		if err := json.Unmarshal(out.Bytes(), &rep); code != exitOK || err != nil || len(rep.Removed) != 0 {
+			t.Fatalf("GC refusal: code %d, report %+v, stderr %q, decode %v", code, rep, errw.String(), err)
+		}
+		found := false
+		for _, item := range rep.Skipped {
+			if item.Name == "ep-clean" && item.Reason == reason && physical(item.Path) == physical(cleanWT) {
+				found = true
+			}
+		}
+		if _, err := os.Stat(cleanWT); err != nil || !found || !bytes.Equal(head, gitState("rev-parse", "HEAD")) ||
+			!bytes.Equal(status, gitState("status", "--porcelain=v1")) {
+			t.Fatalf("want unchanged ep-clean and skip %q: %+v, stat %v", reason, rep, err)
+		}
+	}
+	assertSkip("operator step lookup unavailable")
+	setRecord(repo)
+	if got := episodeGCCredentialsPath(wsRoot, ""); got != "" {
+		t.Fatalf("different configured workspace selected credentials %q", got)
+	}
+	assertSkip("non-terminal operator step", "--credentials-db", creds)
+	alias := filepath.Join(base, "workspace-alias")
+	if err := os.Symlink(wsRoot, alias); err != nil {
+		t.Fatal(err)
+	}
+	setRecord(wsRoot)
+	if got := episodeGCCredentialsPath(alias, ""); got != creds {
+		t.Fatalf("symlink scan root selected %q, want %q", got, creds)
+	}
+	setRecord(alias)
+	assertSkip("non-terminal operator step")
+	if _, err := raw.ExecContext(ctx, "UPDATE operator_action_steps SET phase='observed',receipt=? WHERE kind='worktree'", []byte("stopped: retained")); err != nil {
+		t.Fatal(err)
+	}
+	if err := os.Remove(record.Record); err != nil {
+		t.Fatal(err)
+	}
+	if codePath := episodeGCCredentialsPath(defaultWorkspaceRootDir(), ""); codePath != "" {
+		t.Fatalf("missing workspace selected a credentials path: %q", codePath)
+	}
+	if err := os.MkdirAll(defaultWorkspaceRootDir(), 0700); err != nil {
+		t.Fatal(err)
+	}
+	if got := episodeGCCredentialsPath(defaultWorkspaceRootDir(), ""); got != defaultCreds {
+		t.Fatalf("default credentials path = %q, want %q", got, defaultCreds)
+	}
+	if err := os.WriteFile(record.Record, []byte("{"), 0600); err != nil {
+		t.Fatal(err)
+	}
+	if got := episodeGCCredentialsPath(defaultWorkspaceRootDir(), ""); got != "" {
+		t.Fatalf("malformed record fell back to credentials %q", got)
+	}
+	t.Run("missing_home", func(t *testing.T) {
+		t.Setenv("HOME", "")
+		t.Chdir(t.TempDir())
+		p := homeServicePaths("")
+		if err := os.MkdirAll(filepath.Dir(p.Record), 0700); err != nil {
+			t.Fatal(err)
+		}
+		data, err := json.Marshal(serviceRecord{Schema: "world/service/v1", Label: worldServiceLabel,
+			Address: "127.0.0.1:7644", Plist: p.Plist, Binary: p.Binary, Credentials: creds, Workspace: wsRoot})
+		if err != nil {
+			t.Fatal(err)
+		}
+		if err := os.WriteFile(p.Record, data, 0600); err != nil {
+			t.Fatal(err)
+		}
+		if got := episodeGCCredentialsPath(wsRoot, ""); got != "" {
+			t.Fatalf("missing HOME selected credentials from cwd: %q", got)
+		}
+	})
+	setRecord(alias)
+	if got := episodeGCCredentialsPath(filepath.Join(base, "missing"), ""); got != "" {
+		t.Fatalf("unresolvable scan root selected credentials %q", got)
+	}
+	setRecord(filepath.Join(base, "missing"))
+	if got := episodeGCCredentialsPath(wsRoot, ""); got != "" {
+		t.Fatalf("unresolvable recorded root selected credentials %q", got)
+	}
+	setRecord(alias)
+	t.Run("relative_root", func(t *testing.T) {
+		t.Chdir(base)
+		if got := episodeGCCredentialsPath(filepath.Base(wsRoot), ""); got != creds {
+			t.Fatalf("relative scan selected %q, want %q", got, creds)
+		}
+	})
+
+	// Dry run leaves everything intact and a terminal-only ledger allows GC.
 	var out, errw bytes.Buffer
 	code := run([]string{"episode", "gc", "--workspace-root", wsRoot, "--exec-cache-dir", cacheDir, "--dry-run", "--json"}, &out, &errw)
 	if code != exitOK {
@@ -112,7 +284,7 @@ func TestEpisodeGCFailSafeAndReclaims(t *testing.T) {
 	if code != exitOK {
 		t.Fatalf("live exit %d, stderr %q", code, errw.String())
 	}
-	if _, err := os.Stat(cleanWT); !os.IsNotExist(err) {
+	if _, err := os.Stat(cleanWT); !errors.Is(err, os.ErrNotExist) {
 		t.Errorf("cleanWT still exists after gc")
 	}
 	if _, err := os.Stat(dirtyWT); err != nil {
@@ -121,13 +293,13 @@ func TestEpisodeGCFailSafeAndReclaims(t *testing.T) {
 	if _, err := os.Stat(unmergedWT); err != nil {
 		t.Errorf("unmergedWT was deleted: %v", err)
 	}
-	if _, err := os.Stat(probeCache); !os.IsNotExist(err) {
+	if _, err := os.Stat(probeCache); !errors.Is(err, os.ErrNotExist) {
 		t.Errorf("probeCache still exists")
 	}
-	if _, err := os.Stat(orphanCache); !os.IsNotExist(err) {
+	if _, err := os.Stat(orphanCache); !errors.Is(err, os.ErrNotExist) {
 		t.Errorf("orphanCache still exists")
 	}
-	if _, err := os.Stat(cleanCache); !os.IsNotExist(err) {
+	if _, err := os.Stat(cleanCache); !errors.Is(err, os.ErrNotExist) {
 		t.Errorf("cleanCache still exists")
 	}
 	if _, err := os.Stat(dirtyCache); err != nil {
diff --git a/host/daemon/auto_gc.go b/host/daemon/auto_gc.go
index 26a503e0..6c44bd04 100644
--- a/host/daemon/auto_gc.go
+++ b/host/daemon/auto_gc.go
@@ -9,6 +9,7 @@ import (
 	"strings"
 	"time"
 
+	"github.com/sunholo-data/ailang-world/host/store"
 	"github.com/sunholo-data/ailang-world/host/toolprobe"
 )
 
@@ -38,6 +39,8 @@ type EpisodeGCOptions struct {
 	DryRun        bool
 	Now           time.Time
 	IsActive      func(episode string) bool
+	// A missing lookup, a lookup error, or an unresolved effect bars worktree removal.
+	HasOpenOperatorStep func(ctx context.Context, episode string) (bool, error)
 }
 
 const (
@@ -63,27 +66,31 @@ func (d *Daemon) startAutoGC(dbPath string) {
 			case <-d.aguiStop:
 				return
 			case <-timer.C:
-				d.runAutoGCPass(wsRoot, cacheDir)
+				d.runAutoGCPass(wsRoot, cacheDir, time.Now())
 				timer.Reset(autoGCInterval)
 			}
 		}
 	}()
 }
 
-func (d *Daemon) runAutoGCPass(wsRoot, cacheDir string) {
+func (d *Daemon) runAutoGCPass(wsRoot, cacheDir string, now time.Time) EpisodeGCReport {
 	ctx, cancel := context.WithTimeout(context.Background(), autoGCPassBudget)
 	defer cancel()
 	rep := RunEpisodeGC(ctx, EpisodeGCOptions{
 		WorkspaceRoot: wsRoot,
 		ExecCacheDir:  cacheDir,
 		OlderThan:     autoGCOlderThan,
-		Now:           time.Now(),
+		Now:           now,
 		IsActive:      d.isEpisodeActive,
+		HasOpenOperatorStep: func(ctx context.Context, episode string) (bool, error) {
+			return store.EpisodeHasUnresolvedOperatorSteps(ctx, d.credentialsPath, episode)
+		},
 	})
 	if len(rep.Removed) > 0 && d.errLog != nil {
 		fmt.Fprintf(d.errLog, "ailang-worldd: auto-gc reclaimed %s (%d removed, %d skipped)\n",
 			FormatGCBytes(rep.ReclaimedBytes), len(rep.Removed), len(rep.Skipped))
 	}
+	return rep
 }
 
 func (d *Daemon) isEpisodeActive(episode string) bool {
@@ -147,6 +154,23 @@ func RunEpisodeGC(ctx context.Context, opts EpisodeGCOptions) EpisodeGCReport {
 					continue
 				}
 				commonDir := gitCommonDir(ctx, wtPath)
+				reason := "operator step lookup unavailable"
+				if opts.HasOpenOperatorStep != nil {
+					unresolved, err := opts.HasOpenOperatorStep(ctx, name)
+					switch {
+					case err != nil || ctx.Err() != nil:
+					case unresolved:
+						reason = "non-terminal operator step"
+					default:
+						reason = ""
+					}
+				}
+				if reason != "" {
+					rep.Skipped = append(rep.Skipped, EpisodeGCItem{
+						Kind: "worktree", Name: name, Path: wtPath, Bytes: size, Reason: reason,
+					})
+					continue
+				}
 				if !opts.DryRun {
 					if err := removeWorktree(ctx, wtPath, commonDir); err != nil {
 						rep.Skipped = append(rep.Skipped, EpisodeGCItem{
diff --git a/host/daemon/auto_gc_s5e4_test.go b/host/daemon/auto_gc_s5e4_test.go
new file mode 100644
index 00000000..0e9f49e1
--- /dev/null
+++ b/host/daemon/auto_gc_s5e4_test.go
@@ -0,0 +1,125 @@
+package daemon
+
+import (
+	"bytes"
+	"context"
+	"errors"
+	"os"
+	"os/exec"
+	"path/filepath"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func TestEpisodeGCOperatorSteps(t *testing.T) {
+	root := t.TempDir()
+	ws := filepath.Join(root, "workspaces")
+	if err := os.Mkdir(ws, 0700); err != nil {
+		t.Fatal(err)
+	}
+	wt := filepath.Join(ws, "ep")
+	git := func(dir string, args ...string) []byte {
+		t.Helper()
+		cmd := exec.CommandContext(boundedTestContext(t), "git", args...)
+		cmd.Dir = dir
+		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=Fixture", "GIT_AUTHOR_EMAIL=fixture@example.net",
+			"GIT_COMMITTER_NAME=Fixture", "GIT_COMMITTER_EMAIL=fixture@example.net")
+		out, err := cmd.CombinedOutput()
+		if err != nil {
+			t.Fatalf("git %v: %v: %s", args, err, out)
+		}
+		return out
+	}
+	git(root, "init", "-b", "main")
+	git(root, "commit", "--allow-empty", "-m", "base")
+	git(root, "worktree", "add", "--detach", wt, "main")
+	info, err := os.Stat(wt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	physical := func(path string) string {
+		t.Helper()
+		got, err := filepath.EvalSymlinks(path)
+		if err != nil {
+			t.Fatal(err)
+		}
+		return got
+	}
+	path := physical(wt)
+	head, status := git(wt, "rev-parse", "HEAD"), git(wt, "status", "--porcelain=v1")
+	if len(status) != 0 {
+		t.Fatalf("fixture is dirty: %q", status)
+	}
+	opts := EpisodeGCOptions{WorkspaceRoot: ws, Now: info.ModTime().Add(73 * time.Hour), OlderThan: 72 * time.Hour}
+	checkSkip := func(t *testing.T, rep EpisodeGCReport, reason string) {
+		t.Helper()
+		if len(rep.Removed) != 0 || len(rep.Skipped) != 1 || rep.Skipped[0].Name != "ep" ||
+			rep.Skipped[0].Reason != reason || physical(rep.Skipped[0].Path) != path {
+			t.Fatalf("want sole worktree skip %q: %+v", reason, rep)
+		}
+		if _, err := os.Stat(wt); err != nil || !bytes.Equal(head, git(wt, "rev-parse", "HEAD")) ||
+			!bytes.Equal(status, git(wt, "status", "--porcelain=v1")) {
+			t.Fatalf("refusal changed worktree or Git evidence: %v", err)
+		}
+	}
+	t.Run("candidate", func(t *testing.T) {
+		o := opts
+		o.DryRun = true
+		o.HasOpenOperatorStep = func(context.Context, string) (bool, error) { return false, nil }
+		rep := RunEpisodeGC(boundedTestContext(t), o)
+		if len(rep.Skipped) != 0 || len(rep.Removed) != 1 || rep.Removed[0].Reason != "clean & merged" ||
+			physical(rep.Removed[0].Path) != path {
+			t.Fatalf("successful negative lookup must allow the candidate: %+v", rep)
+		}
+	})
+	for _, name := range []string{"protected", "error", "timeout", "canceled", "unavailable"} {
+		t.Run(name, func(t *testing.T) {
+			ctx, cancel := context.WithCancel(boundedTestContext(t))
+			defer cancel()
+			o := opts
+			reason := "operator step lookup unavailable"
+			switch name {
+			case "protected":
+				reason = "non-terminal operator step"
+				o.HasOpenOperatorStep = func(_ context.Context, episode string) (bool, error) {
+					return episode == "ep", nil
+				}
+			case "error":
+				o.HasOpenOperatorStep = func(context.Context, string) (bool, error) {
+					return false, errors.New("fixture store error")
+				}
+			case "timeout":
+				o.HasOpenOperatorStep = func(context.Context, string) (bool, error) {
+					return false, context.DeadlineExceeded
+				}
+			case "canceled":
+				o.HasOpenOperatorStep = func(context.Context, string) (bool, error) {
+					cancel()
+					return false, nil
+				}
+			}
+			checkSkip(t, RunEpisodeGC(ctx, o), reason)
+		})
+		if t.Failed() {
+			return
+		}
+	}
+	t.Run("daemon_wiring", func(t *testing.T) {
+		creds := filepath.Join(root, "credentials.db")
+		if _, err := store.UpgradeCredentials(creds); err != nil {
+			t.Fatal(err)
+		}
+		// Invoke the synchronous daemon pass with a temporary configured store.
+		// The timers, server, launchctl and network are absent from this fixture.
+		d := &Daemon{credentialsPath: creds}
+		rep := d.runAutoGCPass(ws, "", opts.Now)
+		if len(rep.Skipped) != 0 || len(rep.Removed) != 1 || filepath.Join(physical(filepath.Dir(rep.Removed[0].Path)), filepath.Base(rep.Removed[0].Path)) != path {
+			t.Fatalf("daemon did not pass its readable credentials store to GC: %+v", rep)
+		}
+		if _, err := os.Stat(wt); !errors.Is(err, os.ErrNotExist) {
+			t.Fatalf("daemon did not remove unprotected worktree: %v", err)
+		}
+	})
+}
diff --git a/host/store/credentials_test.go b/host/store/credentials_test.go
index 6d1744de..a4a55695 100644
--- a/host/store/credentials_test.go
+++ b/host/store/credentials_test.go
@@ -322,6 +322,7 @@ func TestCredentialsMethodsRequireADeadline(t *testing.T) {
 	_, errs["OperatorHeartbeat"] = c.OperatorHeartbeat(ctx)
 	_, errs["ReadOperatorHeartbeat"] = ReadOperatorHeartbeat(ctx, "")
 	_, errs["ReadOperatorEpisodeOwner"] = ReadOperatorEpisodeOwner(ctx, c.display, "ep")
+	_, errs["EpisodeHasUnresolvedOperatorSteps"] = EpisodeHasUnresolvedOperatorSteps(ctx, c.display, "ep")
 	for name, err := range errs {
 		if !errors.Is(err, ErrNoDeadline) {
 			t.Errorf("%s(Background) = %v, want ErrNoDeadline", name, err)
diff --git a/host/store/operator_episode_gc.go b/host/store/operator_episode_gc.go
new file mode 100644
index 00000000..95dbb2e3
--- /dev/null
+++ b/host/store/operator_episode_gc.go
@@ -0,0 +1,26 @@
+package store
+
+import "context"
+
+// EpisodeHasUnresolvedOperatorSteps reads the durable effect ledger without a
+// writer lock. An observation alone does not dispose of an effect. A reserved
+// row still protects the worktree even if its action was marked terminal.
+func EpisodeHasUnresolvedOperatorSteps(ctx context.Context, path, episode string) (bool, error) {
+	if err := requireDeadline(ctx); err != nil {
+		return false, err
+	}
+	db, version, err := openOperatorReadOnly(ctx, path)
+	if err != nil {
+		return false, err
+	}
+	defer db.Close()
+	if version < operatorSchemaVersion {
+		return false, &OlderCredentialsStoreError{Path: path, Version: version, Need: operatorSchemaVersion, Feature: "operator actions"}
+	}
+	var unresolved bool
+	err = db.QueryRowContext(ctx, `SELECT EXISTS (
+		SELECT 1 FROM operator_actions a JOIN operator_action_steps s ON s.action_id = a.id
+		WHERE a.episode_id = ? AND s.kind <> 'operation'
+		AND (s.phase = 'reserved' OR a.terminal_at IS NULL))`, episode).Scan(&unresolved)
+	return unresolved, err
+}
diff --git a/host/store/operator_episode_gc_test.go b/host/store/operator_episode_gc_test.go
new file mode 100644
index 00000000..6767b8b2
--- /dev/null
+++ b/host/store/operator_episode_gc_test.go
@@ -0,0 +1,149 @@
+package store
+
+import (
+	"bytes"
+	"context"
+	"errors"
+	"os"
+	"path/filepath"
+	"strings"
+	"testing"
+	"time"
+)
+
+func TestEpisodeHasUnresolvedOperatorSteps(t *testing.T) {
+	for _, tc := range []struct {
+		name, phase, state string
+		want               bool
+	}{
+		{name: "none"},
+		{name: "operation_only", state: "running"},
+		{name: "reserved", phase: "reserved", state: "running", want: true},
+		{name: "uncertain", phase: "reserved", state: "uncertain", want: true},
+		{name: "observed_running", phase: "observed", state: "running", want: true},
+		{name: "observed_uncertain", phase: "observed", state: "uncertain", want: true},
+		{name: "observed_succeeded", phase: "observed", state: "succeeded"},
+		{name: "observed_failed", phase: "observed", state: "failed"},
+		{name: "observed_expired", phase: "observed", state: "expired"},
+		{name: "observed_cancelled", phase: "observed", state: "cancelled"},
+		{name: "terminal_with_reservation", phase: "reserved", state: "failed", want: true},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			f := newOperatorPinFixture(t, false)
+			if tc.phase != "" {
+				key := f.reserveWorktree(t)
+				if tc.phase == "observed" {
+					if err := f.c.ObserveOperatorStep(f.ctx, key, []byte("created")); err != nil {
+						t.Fatal(err)
+					}
+				}
+			}
+			// Materialize durable outcomes (including abnormal combinations)
+			// subject to the same state and terminal_at CHECK as production.
+			switch tc.state {
+			case "running", "uncertain":
+				operatorExec(t, f.c, "UPDATE operator_actions SET state=?", tc.state)
+			case "succeeded", "failed", "expired", "cancelled":
+				operatorExec(t, f.c, "UPDATE operator_actions SET state=?, owner=NULL, terminal_at=?", tc.state, pinNow+1)
+			}
+			episode := f.in.Start.Episode
+			if tc.name == "none" {
+				episode = "no-action"
+			}
+			// A second episode's reservation must never affect this result.
+			other := putOperator(t, f.c, operatorRequest(t, f.c, 2, "other-episode", "finish"), 1000)
+			operatorExec(t, f.c, `INSERT INTO operator_action_steps (kind,identity,action_id,phase,version,owner,created_at)
+				VALUES ('cleanup','other-effect',?,'reserved',1,'other-runner',?)`, other.ID, pinNow)
+			path := f.c.lock.dbPath
+			before, err := os.ReadFile(path)
+			if err != nil {
+				t.Fatal(err)
+			}
+			got, err := EpisodeHasUnresolvedOperatorSteps(f.ctx, path, episode)
+			if got != tc.want || err != nil {
+				t.Fatalf("GC ledger lookup = %t, %v; want %t, nil", got, err, tc.want)
+			}
+			assertGCLedgerBytes(t, path, before)
+		})
+	}
+}
+
+func TestEpisodeOperatorStepLookupRefusals(t *testing.T) {
+	t.Run("no_deadline", func(t *testing.T) {
+		c := openTestCredentials(t)
+		path := c.lock.dbPath
+		before, err := os.ReadFile(path)
+		if err != nil {
+			t.Fatal(err)
+		}
+		got, err := EpisodeHasUnresolvedOperatorSteps(context.Background(), path, "ep")
+		if got || !errors.Is(err, ErrNoDeadline) {
+			t.Fatalf("unbounded lookup = %t %v", got, err)
+		}
+		assertGCLedgerBytes(t, path, before)
+	})
+	t.Run("missing_store", func(t *testing.T) {
+		path := filepath.Join(t.TempDir(), CredentialsFileName)
+		got, err := EpisodeHasUnresolvedOperatorSteps(testCtx(t), path, "ep")
+		if got || !errors.Is(err, os.ErrNotExist) {
+			t.Fatalf("missing store = %t %v", got, err)
+		}
+		if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
+			t.Fatalf("lookup created or replaced credentials: %v", err)
+		}
+	})
+	t.Run("older_schema", func(t *testing.T) {
+		path := filepath.Join(t.TempDir(), CredentialsFileName)
+		raw := rawDB(t, path)
+		if _, err := raw.Exec(credentialsSchemaSQL); err != nil {
+			t.Fatal(err)
+		}
+		setPositiveFixtureVersion(t, raw, 9)
+		before, err := os.ReadFile(path)
+		if err != nil {
+			t.Fatal(err)
+		}
+		got, err := EpisodeHasUnresolvedOperatorSteps(testCtx(t), path, "ep")
+		old, ok := errors.AsType[*OlderCredentialsStoreError](err)
+		if got || !ok || old.Version != 9 || old.Need != operatorSchemaVersion {
+			t.Fatalf("older schema = %t %v", got, err)
+		}
+		assertGCLedgerBytes(t, path, before)
+	})
+	t.Run("query_error", func(t *testing.T) {
+		c := openTestCredentials(t)
+		path := c.lock.dbPath
+		operatorExec(t, c, "ALTER TABLE operator_action_steps RENAME TO unavailable_steps")
+		before, err := os.ReadFile(path)
+		if err != nil {
+			t.Fatal(err)
+		}
+		got, err := EpisodeHasUnresolvedOperatorSteps(testCtx(t), path, "ep")
+		if got || err == nil || !strings.Contains(err.Error(), "no such table") {
+			t.Fatalf("missing ledger table = %t %v", got, err)
+		}
+		assertGCLedgerBytes(t, path, before)
+	})
+	t.Run("deadline", func(t *testing.T) {
+		c := openTestCredentials(t)
+		path := c.lock.dbPath
+		before, err := os.ReadFile(path)
+		if err != nil {
+			t.Fatal(err)
+		}
+		ctx, cancel := context.WithDeadline(context.Background(), time.Unix(1, 0))
+		defer cancel()
+		got, err := EpisodeHasUnresolvedOperatorSteps(ctx, path, "ep")
+		if got || !errors.Is(err, context.DeadlineExceeded) {
+			t.Fatalf("expired lookup = %t %v", got, err)
+		}
+		assertGCLedgerBytes(t, path, before)
+	})
+}
+
+func assertGCLedgerBytes(t *testing.T, path string, before []byte) {
+	t.Helper()
+	if after, err := os.ReadFile(path); err != nil || !bytes.Equal(before, after) {
+		t.Fatalf("credentials changed during GC lookup: %v", err)
+	}
+}
````
