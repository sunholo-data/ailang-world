# Sprint plan — w-world-address-bar-verbs-v4, sprint S5e3 (iteration 437, row 243)

D-WORLD-119 = A and D-WORLD-120 = A were attended at `82f05004`. The reviewed r1 design has quorum r2 PROCEED 3/3 with all seats present. S1 #606, S2 #609, S2b #611, S3 #643, S4 #651, S4b #654, S5 #657, S5b #659, S5c #661, S5d #666, S5e #669 and S5e2 #672 landed. The S5e2 merge `d9c24a3b` contains tests and received Sonnet PASS 84 with zero blocking findings. The base is `eaa202b1d23f601ce597f063034a061b5bf5dea3` and the lane is `codex:gpt-6.2-sol`.

## 0. Binding constraints

- M1 ships private publication of the signed task, the in-memory mint token and the Claude MCP configuration, plus independent readback and the attended `service runner reconcile-files ACTION-ID` in the same patch. This is the whole sprint.
- The reviewed design remains at r1. Activation follows §§2.2 and 5.5. Use the canonical episode lock and existing step bank, signed proofs, and writer seam. Keep the worktree, branch, start, credential and all files as found during disposition.
- Tests use the real signer, migrated stores, real Git and real file installation under `t.TempDir()`. Scratch holds archive trees and evidence. Leave `TMPDIR` and `GOTMPDIR` inherited. All commands in this plan use fixture stores and services or read-only observations. The controller owns Git writes. The executor applies Appendix P and produces a cumulative snapshot.

## 1. Pre-flight

### 1.1 Baseline at eaa202b1

The first §6 run preceded staging. Its input was the pristine `git archive eaa202b1` export at `$S/base`. `$S/gate-base/ALL_DONE` recorded `GATES_DONE 266s FAILS: none UNINFORMATIVE: pkg-cmd-ailang-worldd race-cmd-ailang-worldd`. Build, gofmt, vet, Linux vet, compile, verify_ail, plain and race store, and named inventories passed. The two whole CLI runs hit `bind: operation not permitted`. The initial archive run skipped `verify_go` because the tracked-blob gate requires Git metadata. The extended run uses a scratch Git bridge for the verifier’s two tracked-blob reads against the pinned base while all test source and scripts still come from the archive.

The extended archive run reached `GATES_DONE 727s FAILS: none UNINFORMATIVE: pkg-cmd-ailang-worldd race-cmd-ailang-worldd verify_go`. The plain verifier read 7,387 tracked blobs, reached the evidence manifest and `go test ./... -count=1`, and returned raw rc 1 with actual bind-denial stacks. The fail-set parser printed `FAILSET_CONTROL host/store TestParserPositive`. Its scratch raw failures are uninformative for whole-suite attribution. The CLI and verifier base sets required by §7 must be recorded outside the sandbox.

The historical row 271 census has 11 CLI and three daemon tests in `$S/base-fails.expected`. Iteration 436’s recorded plain verifier set also contains `TestWorldContextAndStatusline` and `TestMissionConfigDefaultFlowIntegration` (`$S/verify_go.reference-fails`). The daemon `TestSeToolsRunBenchmarkTasksEndToEnd` flaked once in four controller base runs. Use current same-rig sets. Archive failures associated with denied binds or absent metadata cannot establish the operator’s fail set.

### 1.2 Current inputs and coupling inventory

`/bin/bash "$S/anchor-head.sh" "$W"` produces `$S/head-anchors.txt`. The table retains the primary identifying commands and their outputs. `$S/stage-anchors.txt` gives changed symbol locations. `$S/branch-anchors-M1.txt` and §4 anchor branch probes.

| ID | Command | Observed result |
|---|---|---|
| C1 | `git rev-parse HEAD`, `git status --short`, `git diff --name-only d9c24a3b eaa202b1 -- cmd host docs app`, `git ls-tree --name-only eaa202b1 -- cmd host docs app` | Full base above, initially empty status, zero product paths between #672 and #673, roots `cmd`, `docs`, `host`. `git diff --stat d9c24a3b eaa202b1` gives seven documentation/ledger files, 248 added and 15 deleted. |
| C2 | `git show --stat --oneline 90e97b1f` and the sorted overlap in §1.4 | #671 adds allowance schema and CLI, touching 39 paths (1,450 added / 50 deleted). `git grep -nE 'credentialsSchemaVersion =\|allowanceSchemaSQL\|\*StartAllowance' -- host/store` gives `credentials.go:42,52,414,437` and `episode_starts.go:36`. Schema is v11, with allowance DDL on fresh creation and upgrade. The new optional start JSON field is excluded when nil. |
| C3 | `git grep -nE 'func .*MintFromAllowance\|func .*MintSessionFromAllowance\|start.Allowance\|INSERT INTO (session_credentials\|episode_starts\|allowance_draws)' -- host/authority/allowance.go host/store/allowances.go` | Authority entry `:31`. Store entry `:256`, provenance `:312`, three inserts `:329,338,344`. Allowance transactions create no operator mint step. The positive control below shows the signed operator mint writer. |
| C4 | `git grep -nE 'func .*CommitOperatorMint\|func .*readOperatorMint\|func .*operatorMintSession\|UPDATE operator_action_steps' -- host/store/operator_mint.go` | `:94,189,249` are the mint, session binding and receipt read. `:166` observes the signed mint. `sed -n '249,287p' host/store/operator_mint.go` shows the observed mint, action owner, episode/credential join and start JSON digest. M1 keeps that read and adds an explicit file-stage allowance for a file reservation. |
| C5 | `git grep -nE 'func .*operatorCASRow\|func .*reserveOperatorStep\|func .*observeOperatorStep\|func .*RecheckClaimedOperatorAction' -- cmd host` | `operator_cas.go:32`, `operator_steps.go:31,55`, `judge_operator.go:129`. Reused authority and transaction boundaries. `git grep -nF 'proof.Bind(' -- host/store/operator_mint.go` gives the mint proof binding. M1 verifies again after reserve/observe step changes. |
| C6 | `git grep -nE 'func .* (ReadOperatorWorktree\|DisposeOperatorWorktree)\(|now-10\|NOT EXISTS\|phase = .reserved.' -- host/store/operator_worktree.go` | `:73,108` read/dispose; optional pin locations `:56–57,59,89–90,138`. M1 adds a separate file transition that requires the exact observed-created worktree row. |
| C7 | `git grep -nE 'func (operatorStartFacts\|reconstructOperatorStart\|operatorPhysicalTarget\|canonicalLaunchPath\|resolvePresetDaemon\|writeMCPConfig)\(' -- cmd` | `operator_start_input.go:16,68`, `operator_worktree.go:131`, `episode_launch.go:73`, `preset.go:298`, `episode.go:564`. M1 records the canonical served address with the resolved mint/launch facts and shares the former MCP writer’s bytes with independent publication. |
| C8 | `git grep -nE 'func (openRunnerLock\|runnerWrite\|operatorEpisodeLockPath)\(' -- cmd` and `git grep -nF 'syscall.LOCK_EX\|syscall.LOCK_NB' -- cmd/ailang-worldd/operator_runner.go` | `operator_runner.go:25,162`, `operator_episode_lock.go:14`, Flock at `:35`. The episode lock is an existing, nonblocking exclusive lock. The resolver reuses the bounded read, installed config and typed terminal seams. |
| C9 | `git grep -nE 'func (newWorktreeStartRig\|observeProvision\|boundedTestContext\|presetTestContext)\|func .* confirm\(' -- cmd` | `operator_worktree_test.go:31`, `operator_provision_test.go:109`, `m7b_deadline_test.go:10`, `preset_test.go:22`, `service.go:171`. The main file fixture starts from the signed worktree rig and mints through authority. It snapshots Git, file bytes and metadata and the credential DB. `go run "$S/observe-journal.go" "$S"` in M1 printed `fixture credentials schema=11 journal_mode=delete`; committed changes are visible in the main-file digest and key states also have direct SQL assertions. |
| C10 | `go env GOOS GOARCH GOVERSION` and `git grep -nF ubuntu-24.04 -- .github/workflows/ci.yml` | Darwin arm64, Go 1.26.6. Native Go CI uses Ubuntu 24.04 at workflow lines 121 and 446. |
| C11 | `/bin/bash "$S/callers.sh" "$S/st-M1" "$W" scratch` | Positive control `authority.RecheckVerb` in `episode_finish_surface.go:239`. Zero production entry callers of publishing, claiming, minting, or reserving worktree. Four set readbacks and five file-state reads are confined to the dormant publication path and the resolver. One resolver dispatch in `operator_service.go`. 54 daemon routes and three surface POSTs. Daemon, served assets and GC compare byte-identical to base. |

| Inventory | Change and evidence |
|---|---|
| `TestDaemonDependencyAllowlist` | Existing daemon dependencies. Publication and disposition live at the CLI/store boundary. `inv-daemon`. |
| `TestWorldBoundaryDependencyAllowlist` | Existing Go groups and x/sys/unix dependency. `inv-boundary`. |
| `TestProductionGoSurface` | Existing package set. `inv-store`. |
| `TestProductionContextRoots` | Resolver reuses `cli.go:90 budgetContext` (existing pin at `context_roots_test.go:46`); publication derives a child from its caller context. The #671 pins, existing #667 pins and production root inventory remain byte-equal. Explicit gate `go test ./host/store/ -run '^TestProductionContextRoots$'`. |
| `TestCLICredentialWriterSites` | Reuses `runnerWrite` and its current credentials-store writer opening site. `inv-cli`. |
| `TestEveryVerbHasHelp` | Pins the nested new command and typed disposition copy in existing help. `inv-cli`. |
| `TestVerbTableMatchesTheCLISwitch` | Grammar extends within `service runner`; the top-level verb table stays the same. `inv-cli`. |
| `TestQuickstartServiceSectionMatchesTheCLIHelp` | M1 changes top-level, service and runner help and the QUICKSTART service section together. `TestOperatorFilesReconcile` loads and executes the exact example. `inv-cli`. |
| `TestCredentialsMethodsRequireADeadline` | Add the read and record cases to the deadline inventory in `credentials_test.go`. `inv-store`. |
| `TestStoreSurfaceTakesContext` | Both entry points take a caller context. The existing inventory runs against the extended surface. `inv-store`. |
| `TestBoundaryASTWriteGuard` | File writes sit in the existing protected CLI group and SQL transitions stay within store. `inv-boundary`. |
| `TestOperatorNeverDispatch` | Dormant publication. Resolver is confined to the typed attended CLI branch. `inv-verbs` and C11. |
| `TestPluginCommandsMatchTheCLIHelp` (CLI and publisher) | Plugin IDs unchanged. Both named tests run on the final tree via `inv-cli` and `inv-publish`. |

### 1.3 Decisions

1. Require the same-owner observed mint receipt, exact session row and start JSON digest. A standing-allowance draw is unreachable as publication authority through its real producer. `TestOperatorFileAllowanceOrigin` draws an actual allowance for the same episode in an already-created worktree, proves its allowance provenance and credential exist, and sees the missing observed mint receipt with unchanged store, Git and file evidence. Both path admission and transaction admission also refuse allowance provenance.
2. Sign the resolved address as `served_address`, which can be reconstructed from the existing preset read. Compare that address with the pinned preset and installed runner config before generating MCP bytes. The exact prior signed start JSON continues to pin the preset, task, base and token path. S5e2 had no production producer of these signed facts (C11).
3. Take each write under a held per-episode lock. Reserve once, use private exclusive creation and reopen independently for readback. Compare real file identity and timestamps before accepting bytes, re-read again before observing, and revalidate proof after step changes. A stop can settle reserved evidence even after expiry, with a fresh bounded write context after human review.
4. Add the resolver as a scoped stop path. It requires an existing reserved row before a prompt, reads evidence on both sides of typed approval, and stops only wholly absent or exact evidence. It preserves the evidence and consumes the durable receipt transactionally.
5. M1 uses the explicit safe-unit overage. Appendix C receives the optional store pins.

### 1.4 Premises found false or incomplete at HEAD

| ID | Command and observed result | Resolution |
|---|---|---|
| F1 | `comm -12 "$S/changed-base" "$S/proto-files"` gives exactly `cmd/ailang-worldd/episode.go`, `cmd/ailang-worldd/main.go`, `docs/QUICKSTART.md`, `host/store/credentials_test.go`. Inputs are sorted `git diff --name-only 5cac9316 eaa202b1 -- cmd host docs app` and old patch headers. Patch with `-F 0` accepted all 14 paths but offset one episode hunk and emitted `episode.go.orig`. | `diff -u "$S/base/cmd/ailang-worldd/episode.go" "$S/st-M1/cmd/ailang-worldd/episode.go"` showed only the MCP encoding extraction. Removed the scratch backup. A fresh diff against this base has exact offsets. |
| F2 | `python3 "$S/stage-diff.py" M1` immediately after the rebase and cleanup printed 1,239 additions, 20 deletions, 1,259 total, 14 files. | Inherited incremental size reproduced on this HEAD. Hardened M1 grew by 391 changed lines. |
| F3 | C2–C4 identify the v11 credential origin and retained S5b receipt join. `TestOperatorFileAllowanceOrigin` reaches the real allowance mint and refuses before reservation. | Explicit origin decision in §1.3. Preserve the schema change and allow nullable legacy start provenance. |
| F4 | `sed -n '24,86p' "$S/st-M1/cmd/ailang-worldd/operator_files_reconcile.go"` shows the pre-review budget and a fresh post-review write budget. The inherited prototype derived the post-review write deadline from the five-second leader. | Separate human think time from I/O budget, retaining the lock. Expired signed actions can be stopped, including at the credential expiry boundary. |
| F5 | `grep -m 1 -F 'bind: operation not permitted' "$S/gate-base/pkg-cmd-ailang-worldd.raw"` fires using the known positive. | Use fail-set evidence outside the sandbox for whole CLI, CLI race and plain verifier. |
| F6 | `git grep -nF 'AND NOT EXISTS' -- host/store/operator_worktree.go` identifies both absence predicates separately at :56 and :57. | The inherited list groups start/session absence as one survivor. Recheck both operations when pinning this family. |
| F7 | `grep -F 'clean & merged' "$S/probe-M1-GC.log"` is the dry-run evidence at 73 h, with `IsActive` stubbed false. | Appendix C records its activation consequence. |

### 1.5 Size and CI race cost

The controller runs exactly `git diff -M --numstat eaa202b1 HEAD -- cmd host docs app` and sums additions plus deletions. The unfiltered product sum scans the whole archive. The planner used `python3 "$S/stage-diff.py" M1`, a read-only, base-archive comparison with `git diff --no-index -M --numstat` on each changed file. There are no renames or binary edits.

| Milestone | Added | Deleted | Cumulative | Increment | Unfiltered product | STATUS/DELETED/MANIFEST |
|---|---:|---:|---:|---:|---:|---:|
| M1 | 1625 | 25 | **1650** | 1650 | 1650 | 16/0/16 |

M1 is 250 lines over 1,400 and uses the 1,650-line safe-unit ceiling. The added lines bind the rebase’s new allowance origin and signed address, make authorization and action/step evidence independent after the step update, and test the expiry, metadata, reopen, CAS and OS-flag boundaries. Splitting the attended resolver from publication would leave a stranded uncertain effect. Available optional-pin capacity at 1,400 is zero. With the plan and JSON committed, the unfiltered iteration diff is 4308 changed lines.

The file-only CLI race selector added 46.692s of package time in the exact latest stage (`$S/gate-latest/files-race.txt`). The same inherited S5e2 selector took 37.113s at the pristine base and 29.080s at M1. The timing varied with load; the file-only selector is the direct added-work measurement. Whole CLI race remains uninformative under the sandbox.

### 1.6 Portability

Use the same `openat`, `fstatat`, `fchmod`, `fsync` and `flock` code on Darwin and Linux. Exclusive creation avoids Darwin’s lack of `RENAME_NOREPLACE`. Directory mode is exactly 0700 and installed file mode is 0600 even under the fixture’s restrictive umask. Refusal tests assert platform-neutral classes or unchanged output. Each flag and each stat call is isolated in §4, including actual descriptor flags with `fcntl`. Native Ubuntu 24.04 execution remains required alongside Linux vet. With the CI environment's working pinned binaries, run `WORLD_REQUIRE_SSH_KEYGEN=1 go test -v -count=1 -timeout 2m ./cmd/ailang-worldd -run '^TestOperatorFile'`. The CLI fixture supplies a synthetic Darwin service environment while the file and lock calls run on the native kernel.

## 2. How the executor applies a sub-milestone

````bash
export GIT_OPTIONAL_LOCKS=0
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter437-row243-v4s5e3
STATE=$HOME/.ailang/state/world-iter437
S=$STATE/plan-scratch
T=$STATE/executor
B=eaa202b1d23f601ce597f063034a061b5bf5dea3
PLAN=$W/design_docs/planned/w-world-address-bar-verbs-v4s5e3-sprint-plan.md
mkdir -p "$T"
unset SNAP_PARENT SNAP_STATUS_LIST
git -C "$W" merge-base --is-ancestor "$B" HEAD || exit 2
git -C "$W" diff --quiet "$B" -- cmd host docs app || exit 2
test -z "$(git -C "$W" ls-files --others --exclude-standard -- cmd host docs app)" || exit 2
cp "$S"/{bounded.py,failset.py,callers.sh} "$T/" || exit 2
cp -R "$S/gate-bin" "$T/" || exit 2
for name in accept.sh gates.sh snapshot.sh; do
  awk -v n="$name" '$0 == "````bash " n {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$name"
  test -s "$T/$name" && /bin/bash -n "$T/$name" || exit 2
done
M=M1
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch" || exit 2
(cd "$W" && /usr/bin/patch -p1 -E < "$T/$M.patch") || exit 2
/bin/bash "$T/accept.sh" "$W" "$M" fast || exit 2
/bin/bash "$T/snapshot.sh" "$W" "$M" || exit 2
````

Apply M1 once. The resulting snapshot excludes both plan records.

## 3. Milestones

### M1 — Signed start files, independent readback and attended disposition

| Effect boundary | Expected readback and named proof in `cmd/ailang-worldd` |
|---|---|
| Admission and lock | `TestOperatorFilePublicationAdmission` checks missing/foreign mint, bad worktree owner/phase/receipt, moved worktree, existing exact file, linked and nonprivate destination, wrong token/grants/TTL/address, revoked credential/key, exact/future timing, held lock and reserve owner drift. Read store, file and Git snapshots on refusal. `TestOperatorFileAllowanceOrigin` exercises the #671 origin. |
| Before step reserve | `TestOperatorFilesReconcile/before_reservation` keeps action running, with zero file rows and zero new files. The remedy refuses an absent reservation before reaching any terminal. |
| After reserve and before first file | `before_first` reads one exact reserved row v1, no files and no new mint. It proves the episode lock is held and the row contains hashes and no raw token. Actions/credentials crossing exact expiry before the first write stay reserved/uncertain with nothing installed. |
| Between files | `between_task_token` reads task only. `between_token_mcp` reads task and token only. Both preserve a reserved step and refuse typed stop on partial evidence. `competing_file` injects an existing token pathname after preflight, checks bytes and metadata, and requires the uncertainty remedy. |
| Before readback and before observation | `before_readback` and `before_observe` read all three exact outputs, live action owner and reserved step. Neither treats return status as proof. `replaced_before_readback` and `replaced_after_readback` replace exact bytes at a new inode. `changed_file` changes the task at the effect boundary. |
| Successful observation | `all_files` and `no_mcp` check the exact Claude three-file and Codex two-file sets. Check every file’s body, mode, link count, owner and persisted step receipt, plus unchanged Git and service evidence. Reject replay and observed-row disposition before a prompt. |
| Post-step CAS and late authority | `publish_owner_drift`, `publish_mint_drift`, `observe_error`, `late_key_revoked` and `late_signature_drift` prove refusal after the last readback. `owner_changed`, `post_step_version`, `version_changed`, `later_step` and changed file/parent prove refusal and rollback in attended disposition. |
| Recovery | `expired_reopen`, `before_first`, `before_readback`, `before_observe` run the QUICKSTART CLI example with no TTY, explicit refusal and typed approval. Known absence produces `stopped: absent`; a complete exact set produces `stopped: retained`. Partial, foreign and inaccessible evidence remains reserved. |
| Metadata and portability | `TestOperatorFileReadbackMetadata`, `TestOperatorFileSyscalls`, `TestOperatorFileLockFlags` and `TestOperatorFileReopen`. The tests cover regular type, no-follow and no-replace behavior, UID/mode/link/digest, exact and oversized limits, before-open/final-stat replacement, unreadable open and independent legal lock flags. |

Transaction order: read the exact active action plus mint and observed-created worktree, acquire the canonical episode lock and reread, verify unsigned execution inputs against signed facts, bind one directory inode and expected output hashes, reserve via the S3b step bank, recheck live proof/session before each write, install privately without replacement, independently reopen and compare output paths/metadata/content, then compare fresh Git/file evidence and observe. The store rechecks proof and mint after the reservation and successful observation step update. The action CAS and resulting step identity/owner/phase/version/receipt are both predicates in the action update. Failure rolls back both writes. A simulated power loss leaves the last reached durable reservation. A detected uncertain outcome records its version and generation while retaining the reservation. Disposition uses effect-specific readbacks around the prompt and a fresh bounded transaction based on the reviewed row. Each stop preserves the worktree, branch, credentials, start and files.

### 3.1 Design §5.3 row assignment

| Reviewed row | Sprint assignment |
|---|---|
| `TestOperatorFilesReconcile` | M1. §4 tests truncation, digest, regular-file metadata, exclusive creation, file identity and path/parent reopen. |
| `TestOperatorStartCollisions`, `TestOperatorStartChanged` | M1 adds coverage at the signed file boundary and reruns the previous worktree/input tests. |
| `TestOperatorStartOnce`, `TestOperatorWorktreeWindow` | Existing worktree prerequisite and receipt-order tests stay green. Launch portion stays at S6. |
| `TestOperatorMintBinding`, `TestOperatorMintRollback`, `TestOperatorMintReceiptWindow` | Existing mint authority path supplies the fixture and publication proof. M1 exercises the newly rechecked binding at file reservation and before/after writes. |
| `TestOperatorLaunchSession`, `TestOperatorSpawnRefused`, `TestOperatorLaunchCrashWindow`, `TestOperatorLaunchExit` | S6 and paired consumer. No process dispatch from M1. |
| `TestOperatorStartBrowser`, `TestOperatorHerdrUnavailable`, `TestOperatorHerdrHandoff`, `TestOperatorHerdrWindow` | Approved S8/S12 slices. Appendix C. |

### 3.2 Design §5.5 and §2.2 activation

| Reviewed gate | M1 evidence or later boundary |
|---|---|
| `TestOperatorServiceFence`, `TestOperatorServiceOwnership`, `TestOperatorServiceConfig` | Keep typed prompt and installed immutable runner validation. `inv-cli` and the resolver fixture verify the reachable recovery CLI. |
| `TestOperatorServicePair`, `TestOperatorServiceStartup`, `TestOperatorServiceUpgrade` | Existing tests run in §6. S5f/S6 must prove the compatible pair, controlled startup and active drain before any production publication caller. |
| `TestOperatorServiceWriteScope`, `TestOperatorServiceEnvironment` | Fixture checks zero service boot/stop/write effects. §6 adds optional read-only rig observations and inventory checks. Browser still has no file-authority input path. |
| `TestOperatorFinishReadiness`, `TestOperatorSchema10` | Named in §6. Keep E2 readiness and v10 compatibility under the v11 schema. No new DDL in M1. |
| §2.2 canonical lock, expiry and catalogue | Existing per-episode lock, typed CLI, bounded caller and freshness checks. Keep heartbeat/readiness, lease claims, catalogue activation and GC action-aware skip gated until S5f/S6. |

## 4. Mutation table

Named row failures require an actual Go `--- FAIL: <test>` line and no bind denial. Test-only failures at a preexisting positive control or compilation failures are instrumentation results. The mutation ledger includes the changed production path and line, exact replacement, Go selector, raw output, rc, and control plus mutant verdict. Its package is `cmd/ailang-worldd` for the file publisher, resolver, lock helper and the cross-package store changes. The store deadline inventory and help inventories validate coupling. Behavioural proof comes from the runtime tests above.

| ID | Production anchor | Named test in CLI | Named result |
|---|---|---|---|
| `F-TRUNC` | `cmd/ailang-worldd/operator_files.go:205` | `TestOperatorFilesReconcile/competing_file` | KILLED |
| `F-EXCL` | `cmd/ailang-worldd/operator_files.go:205` | `TestOperatorFilesReconcile/competing_file` | KILLED |
| `F-CREATE-NOFOLLOW` | `cmd/ailang-worldd/operator_files.go:205` | `TestOperatorFileSyscalls/create_alias` | MASKED — native exclusive-create alias control: O_CREAT|O_EXCL rejects the existing symlink |
| `F-DIR-NOFOLLOW` | `cmd/ailang-worldd/operator_files.go:105` | `TestOperatorFileSyscalls/directory_open` | KILLED |
| `F-DIR-DIRECTORY` | `cmd/ailang-worldd/operator_files.go:105` | `TestOperatorFileSyscalls/directory_open` | KILLED |
| `F-DIR-CLOEXEC` | `cmd/ailang-worldd/operator_files.go:105` | `TestOperatorFileSyscalls/directory_open` | KILLED |
| `F-READ-NOFOLLOW` | `cmd/ailang-worldd/operator_files.go:145` | `TestOperatorFileSyscalls/read_open` | KILLED |
| `F-READ-NONBLOCK` | `cmd/ailang-worldd/operator_files.go:145` | `TestOperatorFileSyscalls/read_open` | KILLED |
| `F-READ-CLOEXEC` | `cmd/ailang-worldd/operator_files.go:145` | `TestOperatorFileSyscalls/read_open` | KILLED |
| `F-CREATE-CLOEXEC` | `cmd/ailang-worldd/operator_files.go:205` | `TestOperatorFileSyscalls` | KILLED |
| `F-STATAT-FIRST` | `cmd/ailang-worldd/operator_files.go:134` | `TestOperatorFileSyscalls/read_open` | KILLED |
| `F-STATAT-LAST` | `cmd/ailang-worldd/operator_files.go:156` | `TestOperatorFileSyscalls/read_open` | KILLED |
| `F-LOCK-EX` | `cmd/ailang-worldd/operator_runner.go:35` | `TestOperatorFileLockFlags` | KILLED |
| `F-LOCK-NB` | `cmd/ailang-worldd/operator_runner.go:35` | `TestOperatorFileLockFlags` | KILLED |
| `F-DIR-MODE` | `cmd/ailang-worldd/operator_files.go:101` | `TestOperatorFilePublicationAdmission/wrong_destination_mode` | KILLED |
| `F-FILE-MODE` | `cmd/ailang-worldd/operator_files.go:142` | `TestOperatorFilesReconcile/metadata/mode` | KILLED |
| `F-FILE-UID` | `cmd/ailang-worldd/operator_files.go:142` | `TestOperatorFilesReconcile/metadata/owner` | KILLED |
| `F-FILE-NLINK` | `cmd/ailang-worldd/operator_files.go:142` | `TestOperatorFilesReconcile/metadata/hardlink` | KILLED |
| `F-DIGEST` | `cmd/ailang-worldd/operator_files.go:160` | `TestOperatorFilesReconcile/metadata/digest` | KILLED |
| `F-LIMIT-EDGE` | `cmd/ailang-worldd/operator_files.go:156` | `TestOperatorFilesReconcile/metadata/at_limit` | KILLED |
| `F-LIMIT-TRUNCATION` | `cmd/ailang-worldd/operator_files.go:155` | `TestOperatorFilesReconcile/metadata/oversize` | KILLED |
| `F-OPEN-INODE` | `cmd/ailang-worldd/operator_files.go:151` | `TestOperatorFileReopen/open` | KILLED |
| `F-FINAL-INODE` | `cmd/ailang-worldd/operator_files.go:156` | `TestOperatorFileReopen/last_stat` | KILLED |
| `F-FCHMOD` | `cmd/ailang-worldd/operator_files.go:211` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `F-ABSENCE` | `cmd/ailang-worldd/operator_files.go:321` | `TestOperatorFilePublicationAdmission/existing_exact_file` | KILLED |
| `F-NAMES-CODEX` | `cmd/ailang-worldd/operator_files.go:79` | `TestOperatorFilesReconcile/no_mcp` | KILLED |
| `F-MCP-DERIVE` | `cmd/ailang-worldd/operator_files.go:288` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `F-WRITE-IDENTITY` | `cmd/ailang-worldd/operator_files.go:388` | `TestOperatorFilesReconcile/replaced_before_readback` | KILLED |
| `F-POST-READ` | `cmd/ailang-worldd/operator_files.go:396` | `TestOperatorFilesReconcile/replaced_after_readback` | KILLED |
| `F-RAW-TOKEN` | `cmd/ailang-worldd/operator_files.go:279` | `TestOperatorFilePublicationAdmission/wrong_token` | KILLED |
| `F-CONFIG-ADDRESS` | `cmd/ailang-worldd/operator_files.go:231` | `TestOperatorFilePublicationAdmission/wrong_approved_address` | KILLED |
| `F-ACTION-EXPIRE` | `cmd/ailang-worldd/operator_files.go:247` | `TestOperatorFilePublicationAdmission/expired_action` | KILLED |
| `F-CREDENTIAL-EXPIRE` | `cmd/ailang-worldd/operator_files.go:361` | `TestOperatorFilePublicationAdmission/credential_expiry_after_reservation` | KILLED |
| `F-LOCK-READ` | `cmd/ailang-worldd/operator_files.go:275` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `F-GIT-EXACT` | `cmd/ailang-worldd/operator_files.go:309` | `TestOperatorFilePublicationAdmission/worktree_moved` | KILLED |
| `R-RESERVED` | `cmd/ailang-worldd/operator_files_reconcile.go:32` | `TestOperatorFilesReconcile/before_reservation` | KILLED |
| `R-POST-READ` | `cmd/ailang-worldd/operator_files_reconcile.go:79` | `TestOperatorFilesReconcile/changed_during_review` | KILLED |
| `R-POST-PARENT` | `cmd/ailang-worldd/operator_files_reconcile.go:79` | `TestOperatorFilesReconcile/parent_changed_during_review` | MASKED — a changed parent leaves no child proofs, so child comparison independently refuses |
| `R-PARTIAL` | `cmd/ailang-worldd/operator_files_reconcile.go:67` | `TestOperatorFilesReconcile/between_task_token` | KILLED |
| `R-STALE-STEP` | `host/store/operator_files.go:103` | `TestOperatorFilesReconcile/version_changed` | KILLED |
| `R-ACTION-OWNER` | `host/store/operator_files.go:174` | `TestOperatorFilesReconcile/owner_changed` | KILLED |
| `R-STEP-VERSION` | `host/store/operator_files.go:176` | `TestOperatorFilesReconcile/post_step_version` | KILLED |
| `R-OTHER-EFFECT` | `host/store/operator_files.go:177` | `TestOperatorFilesReconcile/later_step` | KILLED |
| `R-CHANGED-ERROR` | `cmd/ailang-worldd/operator_files_reconcile.go:87` | `TestOperatorFilesReconcile/owner_changed` | KILLED |
| `R-TYPED` | `cmd/ailang-worldd/operator_files_reconcile.go:75` | `TestOperatorFilesReconcile/before_first` | KILLED |
| `S-WORKTREE-OWNER` | `host/store/operator_files.go:43` | `TestOperatorFilePublicationAdmission/wrong_worktree_owner` | KILLED |
| `S-WORKTREE-PHASE` | `host/store/operator_files.go:43` | `TestOperatorFilePublicationAdmission/wrong_worktree_phase` | MASKED — reserved schema forces receipt NULL, which independently refuses |
| `S-WORKTREE-RECEIPT` | `host/store/operator_files.go:43` | `TestOperatorFilePublicationAdmission/wrong_worktree_receipt` | KILLED |
| `S-CREDENTIAL-EXPIRE` | `host/store/operator_files.go:191` | `TestOperatorFilePublicationAdmission/expired_credential` | KILLED |
| `S-CREDENTIAL-CREATED` | `host/store/operator_files.go:194` | `TestOperatorFilePublicationAdmission/future_mint` | KILLED |
| `S-SESSION-GRANTS` | `host/store/operator_files.go:228` | `TestOperatorFilePublicationAdmission/wrong_grants` | KILLED |
| `S-SESSION-TTL` | `host/store/operator_files.go:228` | `TestOperatorFilePublicationAdmission/wrong_ttl` | KILLED |
| `S-POST-RESERVE` | `host/store/operator_files.go:151` | `TestOperatorFilePublicationAdmission/reserve_mint_drift` | KILLED |
| `S-POST-OBSERVE` | `host/store/operator_files.go:151` | `TestOperatorFilesReconcile/publish_mint_drift` | KILLED |
| `S-LATE-APPROVAL` | `cmd/ailang-worldd/operator_files.go:399` | `TestOperatorFilesReconcile/late_key_revoked` | KILLED |
| `S-LATE-SIGNATURE` | `cmd/ailang-worldd/operator_files.go:399` | `TestOperatorFilesReconcile/late_signature_drift` | KILLED |
| `S-MINT-OWNER` | `host/store/operator_mint.go:265` | `TestOperatorFilePublicationAdmission/wrong_mint_owner` | KILLED |
| `F-FILE-TYPE` | `cmd/ailang-worldd/operator_files.go:142` | `TestOperatorFilesReconcile/metadata/fifo` | KILLED |
| `F-SIGNED-TASK-DATA` | `cmd/ailang-worldd/operator_files.go:286` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `F-TOKEN-DATA` | `cmd/ailang-worldd/operator_files.go:286` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `F-SKIP-CLAUDE-MCP` | `cmd/ailang-worldd/operator_files.go:287` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `F-BEFORE-RESERVE` | `cmd/ailang-worldd/operator_files.go:326` | `TestOperatorFilesReconcile/before_reservation` | KILLED |
| `F-BEFORE-WRITE` | `cmd/ailang-worldd/operator_files.go:356` | `TestOperatorFilesReconcile/between_token_mcp` | KILLED |
| `F-BEFORE-READBACK` | `cmd/ailang-worldd/operator_files.go:380` | `TestOperatorFilesReconcile/before_readback` | KILLED |
| `F-BEFORE-OBSERVE` | `cmd/ailang-worldd/operator_files.go:392` | `TestOperatorFilesReconcile/before_observe` | KILLED |
| `F-AFTER-RESERVE-SIGNATURE` | `cmd/ailang-worldd/operator_files.go:364` | `TestOperatorFilePublicationAdmission/expiry_after_reservation` | KILLED |
| `F-RECORD-UNCERTAIN` | `cmd/ailang-worldd/operator_files.go:346` | `TestOperatorFilesReconcile/competing_file` | KILLED |
| `F-SIGNED-ADDRESS-FACT` | `cmd/ailang-worldd/operator_start_input.go:73` | `TestOperatorStartFacts` | KILLED |
| `R-NEEDS-RESERVED-AFTER-OBSERVE` | `cmd/ailang-worldd/operator_files_reconcile.go:32` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `S-MINT-FILE-STAGE` | `host/store/operator_mint.go:266` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `R-WRONG-ATTENDED-MODE` | `cmd/ailang-worldd/operator_service.go:293` | `TestOperatorFilesReconcile/before_first` | KILLED |
| `R-PARENT-COMPLETE` | `cmd/ailang-worldd/operator_files_reconcile.go:79` | `TestOperatorFilesReconcile/parent_changed_during_review` | KILLED |
| `F-TOKEN-LEAF` | `cmd/ailang-worldd/operator_files.go:71` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `F-TOKEN-PARENT` | `cmd/ailang-worldd/operator_files.go:71` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `F-TASK-PATH` | `cmd/ailang-worldd/operator_files.go:72` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `F-TASK-HASH` | `cmd/ailang-worldd/operator_files.go:73` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `F-MCP-PRESENCE` | `cmd/ailang-worldd/operator_files.go:74` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `F-MCP-PATH` | `cmd/ailang-worldd/operator_files.go:75` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `F-TARGET-ACTION` | `cmd/ailang-worldd/operator_files.go:297` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `F-TARGET-PATH` | `cmd/ailang-worldd/operator_files.go:297` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `F-TARGET-BRANCH` | `cmd/ailang-worldd/operator_files.go:297` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `F-TARGET-SHA` | `cmd/ailang-worldd/operator_files.go:297` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `F-TARGET-REPO` | `cmd/ailang-worldd/operator_files.go:305` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `F-PARENT-UID` | `cmd/ailang-worldd/operator_files.go:101` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `F-PARENT-DEV` | `cmd/ailang-worldd/operator_files.go:102` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `F-PARENT-INODE` | `cmd/ailang-worldd/operator_files.go:102` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `R-IDENTITY-DIR` | `cmd/ailang-worldd/operator_files_reconcile.go:46` | `TestOperatorFilesReconcile/before_first` | KILLED |
| `R-IDENTITY-UID` | `cmd/ailang-worldd/operator_files_reconcile.go:47` | `TestOperatorFilesReconcile/before_first` | KILLED |
| `R-IDENTITY-NAMES` | `cmd/ailang-worldd/operator_files_reconcile.go:47` | `TestOperatorFilesReconcile/before_first` | KILLED |
| `R-IDENTITY-TASK-HASH` | `cmd/ailang-worldd/operator_files_reconcile.go:48` | `TestOperatorFilesReconcile/before_first` | KILLED |
| `R-IDENTITY-TOKEN-HASH` | `cmd/ailang-worldd/operator_files_reconcile.go:48` | `TestOperatorFilesReconcile/before_first` | KILLED |
| `S-ACTION-CAS-DIGEST` | `host/store/operator_files.go:174` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `S-ACTION-CAS-GENERATION` | `host/store/operator_files.go:174` | `TestOperatorFilesReconcile/all_files` | KILLED |
| `S-STEP-CAS-IDENTITY` | `host/store/operator_files.go:176` | `TestOperatorFilesReconcile/before_first` | KILLED |
| `S-STEP-CAS-OWNER` | `host/store/operator_files.go:176` | `TestOperatorFilesReconcile/before_first` | KILLED |
| `F-MCP-RENDER-ADDRESS` | `cmd/ailang-worldd/episode.go:566` | `TestEpisodeStartProvisionsAndPrintsTheLaunchWithTheTokenInTwoFilesOnly` | KILLED |
| `F-MCP-RENDER-BEARER` | `cmd/ailang-worldd/episode.go:567` | `TestEpisodeStartProvisionsAndPrintsTheLaunchWithTheTokenInTwoFilesOnly` | KILLED |
| `F-MCP-RENDER-TYPE` | `cmd/ailang-worldd/episode.go:566` | `TestEpisodeStartProvisionsAndPrintsTheLaunchWithTheTokenInTwoFilesOnly` | KILLED |

Named replay totals: 95 KILLED, 3 MASKED. Whole-package planner results: 1 UNINFORMATIVE UNDER SANDBOX. The diff census has 122 lines marked AUDIT with no independent measured replacement. Some are context or supporting expressions. Reproduce with `python3 "$S/diff-anchors.py"` and inspect the AUDIT rows; the census does not imply those conditions are killed.

The v10 schema at `host/store/operator_schema.sql:9,20,53,61` validates a live action’s challenge width and separates reserved rows from observed receipts. The directory and file tests pin expected UID against actual descriptor metadata. Remaining raw syscall error returns and some compound admission predicates lack independent fault injection.

The raw branch census is `$S/branch-anchors-M1.txt`. It includes context lines from reused primitives. Use `$S/diff-branch-anchors-M1.txt` for changed hunks and `$S/mutants.json` for executable replacements. Each tested flag has an independent replacement. The equality inversions prove a predicate uses its input on the admitted path, and the refusal cases separately test the named negative inputs. In particular the two `AT_SYMLINK_NOFOLLOW` uses and the close-on-exec flags on directory walking, file read and file installation have separate IDs. The combined type/permission checks deliberately classify unusual platform-specific errno through the exposed refusal or unchanged-state assertion.

### 4.1 Whole-package runs for the controller

````bash
S=$HOME/.ailang/state/world-iter437/plan-scratch
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter437-row243-v4s5e3
/bin/bash "$S/whole-mutants.sh" "$W" "$S/controller-base/pkg-cmd-ailang-worldd.base-fails"
````

`whole-mutant-ids.txt` provides the list and the wrapper accepts explicit trailing IDs for a bounded subset. The selected checkout must contain committed M1 and a clean product tree. The wrapper uses the named checkout, restores each changed file from scratch, and compares the whole CLI fail set to the recorded pristine same-rig base. Run whole-package mutation legs after the broad verifier. The local named replay command is `python3 "$S/mutate.py" named --jobs 2`. §6 and the controller’s real-binary cases have socket-bound tests and must run outside the sandbox for an attributable fail set. The explicit mutations are separate from the additional unmeasured audit anchors, which require independently targeted fixtures or identified masks.

## 5. Files

````text
7	2	cmd/ailang-worldd/episode.go
2	2	cmd/ailang-worldd/help_test.go
1	1	cmd/ailang-worldd/main.go
203	0	cmd/ailang-worldd/operator_file_syscalls_test.go
427	0	cmd/ailang-worldd/operator_files.go
92	0	cmd/ailang-worldd/operator_files_reconcile.go
606	0	cmd/ailang-worldd/operator_files_test.go
6	3	cmd/ailang-worldd/operator_service.go
20	4	cmd/ailang-worldd/operator_start_input.go
10	6	cmd/ailang-worldd/operator_start_input_test.go
4	1	cmd/ailang-worldd/operator_worktree_test.go
1	1	cmd/ailang-worldd/service.go
7	0	docs/QUICKSTART.md
2	0	host/store/credentials_test.go
232	0	host/store/operator_files.go
5	5	host/store/operator_mint.go
````

The file and SQLite effects belong in the trusted host boundary because they operate on OS descriptors, filesystem metadata and local transactional rows. The signed action law remains in the reviewed approval package. The patch touches no daemon route, served assets, row-271 tests, GitHub auth, driver, launch agent, schema, or frozen AILANG core.

## 6. Sprint gates

`$S/gate-latest/ALL_DONE` recorded `GATES_DONE 269s FAILS: none UNINFORMATIVE: pkg-cmd-ailang-worldd race-cmd-ailang-worldd`. Native build/vet, Linux vet, compile, verify_ail, store plain/race, named inventories, both named acceptance invocations and file race passed. The earlier, combined CLI acceptance run timed out at 120.415 s under load during the S5e2 worktree cases. The final acceptance splits the named publication and regression selectors, each with its original two-minute bound.

````bash gates.sh
#!/bin/bash
# gates.sh TREE LOGDIR base|scratch|final [GIT-REF] [BASE-LOGDIR]
set -uo pipefail
export GIT_OPTIONAL_LOCKS=0
export WORLD_GATE_GIT="$(command -v git)"
T=$(cd "$1" && pwd -P); L=$2; MODE=$3; REF=${4:-$T}; BASE_LOG=${5:-}
S=$(cd "$(dirname "$0")" && pwd -P)
B=eaa202b1d23f601ce597f063034a061b5bf5dea3
case "$MODE" in base|scratch|final) ;; *) exit 2;; esac
mkdir -p "$L"/{tmp,state,cache} || exit 2
L=$(cd "$L" && pwd -P)
cd "$T" || exit 2
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
export AILANG_STATE_DIR="$L/state" AILANG_CACHE_DIR="$L/cache"
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
unset WORLD_UPDATE_OPERATOR_ACTION_GOLDEN WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN
: > "$L/legs.txt"
FAILS=(); UNKNOWN=(); START=$(date +%s)
line() { printf '%s\n' "$*" | tee -a "$L/legs.txt"; }
skip() { line "LEG $1 SKIPPED 0s ($2)"; }
bounded() { python3 "$S/bounded.py" "$@"; }
denied() { grep -Ei 'bind: (operation not permitted|permission denied)|listen tcp.*(operation not permitted|permission denied)' "$1"; }
leg() {
  local name=$1 rc=0 start; shift
  start=$(date +%s)
  "$@" > "$L/$name.txt" 2>&1 || rc=$?
  local result="rc=$rc $(( $(date +%s)-start ))s"
  if denied "$L/$name.txt" >/dev/null; then
    result="$result UNINFORMATIVE UNDER SANDBOX"
    UNKNOWN+=("$name")
  elif [ "$rc" -ne 0 ]; then
    FAILS+=("$name")
  fi
  line "LEG $name $result"
}
observations() {
  local path
  for path in "$HOME/.ailang-world"/{service.json,operator-runner.json,bin/ailang-worldd} "$HOME/.ailang/world/credentials.db"; do
    if [ -f "$path" ]; then shasum -a 256 "$path"; wc -l "$path"
    elif [ -e "$path" ] || [ -L "$path" ]; then ls -ld "$path"
    else printf 'absent %s\n' "$path"; fi
  done
  path="$HOME/Library/LaunchAgents"
  if [ -d "$path" ]; then
    printf 'directory %s\n' "$path"
    find "$path" -maxdepth 1 -type f -name '*ailang*.plist' -print | LC_ALL=C sort |
      while IFS= read -r path; do shasum -a 256 "$path"; wc -l "$path"; done
  else printf 'absent %s\n' "$path"; fi
  path="$HOME/.ailang/world-ws"
  if [ -d "$path" ]; then
    printf 'directory %s\n' "$path"
    find "$path" -maxdepth 1 -type d -print | LC_ALL=C sort | shasum -a 256
    find "$path" -maxdepth 1 -type d -print | wc -l
  else printf 'absent %s\n' "$path"; fi
}
if [ "${WORLD_GATE_OBSERVE:-0}" = 1 ]; then leg observation-before observations
else skip observation-before 'enable read-only rig observations with WORLD_GATE_OBSERVE=1'; fi
{ uptime; ps -A -o pid,pcpu,comm -r | head -12; } > "$L/load.txt" 2>&1
fmt() {
  local output; output=$(gofmt -l cmd host) || return
  printf '%s\n' "$output"
  [ -z "$output" ] || {
    [ "$output" = host/daemon/daemon.go ] &&
      git -C "$REF" show "$B:$output" | cmp - "$output"
  }
}
leg gofmt fmt
leg build bounded 240 go build ./...
leg vet bounded 240 go vet ./...
leg linux-vet env GOOS=linux python3 "$S/bounded.py" 240 go vet ./cmd/ailang-worldd ./host/store
leg compile bounded 240 go test -run '^$' ./...
leg verify_ail env PATH="$S/gate-bin:$PATH" WORLD_GATE_TMP="$L/tmp" python3 "$S/bounded.py" 380 ./scripts/verify_ail.sh
leg pkg-host-store bounded 510 go test -count=1 -timeout 8m ./host/store
leg race-host-store bounded 510 go test -count=1 -race -timeout 8m ./host/store
record() {
  local tag=$1 rc=0 parser=0; shift
  "$@" > "$L/$tag.raw" 2>&1 || rc=$?
  printf '%s\n' "$rc" > "$L/$tag.rc"
  python3 "$S/failset.py" "$L/$tag.raw" "$rc" > "$L/$tag.fails" || parser=$?
  if denied "$L/$tag.raw"; then return 77; fi
  if [ "$parser" -ne 0 ] || [ "$rc" -gt 1 ] ||
    grep -Eq '\[build failed\]|test timed out|BOUND EXCEEDED' "$L/$tag.raw"; then
    echo 'INSTRUMENT FAILURE'; return 2
  fi
  cat "$L/$tag.fails"
}
leg pkg-cmd-ailang-worldd record pkg-cmd-ailang-worldd bounded 510 go test -count=1 -timeout 8m ./cmd/ailang-worldd
leg race-cmd-ailang-worldd record race-cmd-ailang-worldd bounded 510 go test -count=1 -race -timeout 8m ./cmd/ailang-worldd
if [ "${GATES_VERIFY_ARCHIVE:-0}" = 1 ]; then
  export WORLD_GATE_ARCHIVE="$T" WORLD_GATE_GIT_REF="$REF" WORLD_GATE_BASE="$B"
fi
if [ "$MODE" = final ] || [ "${GATES_CHECKOUT_BASE:-0}" = 1 ] || [ "${GATES_VERIFY_ARCHIVE:-0}" = 1 ]; then
  leg verify_go record verify_go env VERIFY_GO_LEGS=plain PATH="$S/gate-bin:$PATH" WORLD_GATE_TMP="$L/tmp" python3 "$S/bounded.py" 1500 ./scripts/verify_go.sh
else skip verify_go 'archive has no Git metadata for tracked-binary gate'; fi
if [ -n "$BASE_LOG" ]; then
  # Recording above must finish before any comparison can fail.
  for tag in pkg-cmd-ailang-worldd race-cmd-ailang-worldd verify_go; do
    if [ ! -f "$L/$tag.raw" ]; then skip "compare-$tag" 'recording leg skipped'
    elif denied "$L/$tag.raw" >/dev/null; then skip "compare-$tag" 'loopback denial'
    else leg "compare-$tag" diff -u "$BASE_LOG/$tag.base-fails" "$L/$tag.fails"; fi
  done
else
  for tag in pkg-cmd-ailang-worldd race-cmd-ailang-worldd verify_go; do
    if [ -f "$L/$tag.fails" ] && ! denied "$L/$tag.raw" >/dev/null &&
      ! grep -Fq 'INSTRUMENT FAILURE' "$L/$tag.txt"; then
      cp "$L/$tag.fails" "$L/$tag.base-fails"
    fi
  done
fi
inventory() {
  local pkg=$1 names=$2 name output
  output=$(mktemp "$L/inventory.XXXXXX") || return
  bounded 150 go test -v -count=1 -timeout 2m "./$pkg" -run "^($names)\$" > "$output" 2>&1
  local rc=$?
  cat "$output"
  [ "$rc" -eq 0 ] || return "$rc"
  while IFS= read -r name; do
    grep -Fq -- "--- PASS: $name (" "$output" || return 2
  done < <(printf '%s\n' "$names" | tr '|' '\n')
}
leg inv-boundary inventory host/boundary 'TestWorldBoundaryDependencyAllowlist|TestBoundaryASTWriteGuard'
leg inv-daemon inventory host/daemon 'TestDaemonDependencyAllowlist|TestDaemonCredentialWriterSites|TestFrozenV1RouteTableMatchesMux'
leg inv-verbs inventory host/verbs 'TestOperatorNeverDispatch|TestOperatorOtherVerbParity'
leg inv-cli inventory cmd/ailang-worldd 'TestEveryVerbHasHelp|TestVerbTableMatchesTheCLISwitch|TestCLICredentialWriterSites|TestQuickstartServiceSectionMatchesTheCLIHelp|TestPluginCommandsMatchTheCLIHelp|TestOperatorServiceFence|TestOperatorServiceOwnership|TestOperatorServiceConfig|TestOperatorServicePair|TestOperatorServiceStartup|TestOperatorServiceUpgrade|TestOperatorServiceWriteScope|TestOperatorServiceEnvironment'
leg inv-store inventory host/store 'TestProductionContextRoots|TestProductionGoSurface|TestStoreSurfaceTakesContext|TestCredentialsMethodsRequireADeadline|TestOperatorSchema10|TestOperatorFinishReadiness'
leg inv-publish inventory cmd/world-publish TestPluginCommandsMatchTheCLIHelp
leg s5e2-race bounded 150 go test -race -v -count=1 -timeout 2m ./cmd/ailang-worldd -run '^(TestOperatorWorktreeWindow|TestOperatorWorktreeRecoveryGuards|TestOperatorWorktreeReadbackGuards|TestOperatorWorktreeLockedReread)$'
if [ "$MODE" = base ]; then
  skip accept 'publication tests absent at base'
  skip files-race 'publication tests absent at base'
else
  leg accept /bin/bash "$S/accept.sh" "$T" M1 fast
  leg files-race bounded 150 go test -race -v -count=1 -timeout 2m ./cmd/ailang-worldd -run '^TestOperatorFile'
fi
leg callers /bin/bash "$S/callers.sh" "$T" "$REF" "$MODE"
if [ "${WORLD_GATE_OBSERVE:-0}" = 1 ]; then
  leg observation-after observations
  leg observation-compare cmp "$L/observation-before.txt" "$L/observation-after.txt"
else skip observation-after 'enable read-only rig observations with WORLD_GATE_OBSERVE=1'; fi
line "GATES_DONE $(( $(date +%s)-START ))s FAILS: ${FAILS[*]:-none} UNINFORMATIVE: ${UNKNOWN[*]:-none}"
tail -1 "$L/legs.txt" > "$L/ALL_DONE"
[ "${#FAILS[@]}" -eq 0 ] || exit 1
[ "${#UNKNOWN[@]}" -eq 0 ] || exit 77
````

## 7. Executor directive notes

````bash accept.sh
#!/bin/bash
set -euo pipefail
W=$1; M=$2; PART=$3
test "$M" = M1
test "$PART" = fast
S=$(cd "$(dirname "$0")" && pwd -P)
export AILANG_BIN="$HOME/.pinned-ailang/ailang"
export WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
unset WORLD_UPDATE_OPERATOR_ACTION_GOLDEN WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN
run_dir=$(mktemp -d "$S/accept.XXXXXX")
export AILANG_STATE_DIR="$run_dir/state" AILANG_CACHE_DIR="$run_dir/cache"
mkdir -p "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR"
cd "$W"
run() {
  local pkg=$1 names=$2 name log rc=0
  log=$(mktemp "$run_dir/${pkg//\//-}.XXXXXX")
  python3 "$S/bounded.py" 180 go test -v -count=1 -timeout 2m "./$pkg" -run "^($names)\$" 2>&1 | tee "$log" || rc=$?
  if grep -Ei 'bind: (operation not permitted|permission denied)' "$log"; then
    echo 'UNINFORMATIVE UNDER SANDBOX'
    exit 77
  fi
  test "$rc" -eq 0
  while IFS= read -r name; do
    grep -Fq -- "--- PASS: $name (" "$log"
  done < <(printf '%s\n' "$names" | tr '|' '\n')
}
run host/store 'TestReserveOperatorWorktreeGuards|TestDisposeOperatorUncertainStepGuards|TestReadOperatorWorktreeGuards|TestOperatorStepReceiptAtLimit|TestOperatorStepRestart|TestOperatorStepGuards|TestCredentialsMethodsRequireADeadline|TestOperatorSchema10|TestProductionContextRoots'
run host/authority 'TestOperatorMintBinding|TestOperatorMintRollback|TestOperatorMintReceiptWindow'
run cmd/ailang-worldd 'TestOperatorFilesReconcile|TestOperatorFileReadbackMetadata|TestOperatorFilePublicationAdmission|TestOperatorFileAllowanceOrigin|TestOperatorFileSyscalls|TestOperatorFileLockFlags|TestOperatorFileReopen|TestEveryVerbHasHelp|TestQuickstartServiceSectionMatchesTheCLIHelp'
run cmd/ailang-worldd 'TestEpisodeStartProvisionsAndPrintsTheLaunchWithTheTokenInTwoFilesOnly|TestOperatorMintRecovery|TestQuickstartRunnerMintRecovery|TestOperatorStartFacts|TestOperatorStartChanged|TestOperatorWorktreeSurvivors|TestOperatorWorktreeRecoveryGuards|TestOperatorWorktreeReadbackGuards|TestOperatorReadbackErrorIsNeverCreated|TestOperatorWorktreeLockedReread|TestOperatorWorktreeWindow|TestOperatorStartOnce|TestOperatorStartCollisions|TestOperatorWorktreeDispositionCAS'
printf 'ACCEPT_OK %s %s\n' "$M" "$PART"
````

Capture the cumulative snapshot at `$HOME/.ailang/state/world-iter437/snap/M1`. Git reads the product pathspec relative to base and skips absent roots. The two plan files live outside the product pathspec.

````bash snapshot.sh
#!/bin/bash
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=$1
M=$2
B=eaa202b1d23f601ce597f063034a061b5bf5dea3
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter437/snap}/$M
mkdir -p "$(dirname "$D")"
mkdir "$D"
cd "$W"
if [ -n "${SNAP_STATUS_LIST:-}" ]; then
  cp "$SNAP_STATUS_LIST" "$D/STATUS"
else
  {
    git diff --name-status --no-renames "$B" -- cmd host docs app
    git ls-files --others --exclude-standard -- cmd host docs app | sed $'s/^/A\t/'
  } | LC_ALL=C sort -u -t "$(printf '\t')" -k2 > "$D/STATUS"
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

Before the first product application, the controller records the current outside-sandbox base sets on an existing named checkout at `eaa202b1` with the base product tree. The recorder permits a subsequent commit containing only the two sprint plan files. Keep each recording even when a comparison reports drift.

````bash
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter437-row243-v4s5e3
STATE=$HOME/.ailang/state/world-iter437
T=$STATE/executor
S=$STATE/plan-scratch
BASE_CHECKOUT=${BASE_CHECKOUT:-$W}
/bin/bash "$S/baseline-controller.sh" "$BASE_CHECKOUT" "$S/controller-base"
````

After the plan commit and M1 snapshot, commit M1, then run the final gates when the rig has capacity. If running with rig observations enabled, use read-only snapshots and permit missing paths.

````bash
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter437-row243-v4s5e3
STATE=$HOME/.ailang/state/world-iter437
T=$STATE/executor
S=$STATE/plan-scratch
printf '%s\n' 'feat(operator): publish signed start files and reconcile the reserved file effect' > "$T/M1.msg"
/bin/bash "$STATE/commit_ms.sh" M1 "$T/M1.msg" 'AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang /bin/bash "'"$T"'/accept.sh" "'"$W"'" M1 fast' || exit
WORLD_GATE_OBSERVE=1 /bin/bash "$T/gates.sh" "$W" "$S/gate-controller-final" final "$W" "$S/controller-base" || exit
/bin/bash "$S/whole-mutants.sh" "$W" "$S/controller-base/pkg-cmd-ailang-worldd.base-fails"
````

The verifier runs with GOFLAGS unset and `VERIFY_GO_LEGS=plain` on the rig. Its fail-set parser reads plain Go output and prints a firing known-positive on every leg. Any nonzero rc with an empty attributable set is an instrument failure. All outputs are recorded before comparisons. Compare build and package verdicts separately, preserve raw artifacts and re-run a load-sensitive discrepancy against the untouched package and test in isolation. Use native Ubuntu 24.04 CI for OS execution. The commit-time command names socket-free tests only.

The exact Appendix P round trip in `$S/rt-OzVjqJ` applied the 1,895-line patch to a fresh base archive with `/usr/bin/patch -p1 -E`. Full archive and each of `cmd`, `host` and `docs` returned `diff -r` rc 0. App was absent on both sides. The inherited backup census stayed identical. The exact script fences matched the tested support copies and passed Bash syntax checks. Extracted snapshot script reproduced the tracked status and verified 16 file hashes (16/0/16) through read-only base Git and scratch output. See `$S/roundtrip-final.log`.

## 8. Risks

- The rig rejects loopback while CI skips the row-271 real-binary tests. Keep the known-red controls and use same-rig fail-set equality around M1. Native Ubuntu is the publication syscall arbiter.
- Opening exclusive files permits a crash to leave a private partial set. The existing reservation and retained evidence keep it distinguishable from failure before effect. The resolver refuses that partial state and allows stop only when the set is wholly absent or exact.
- File work runs against local private paths. Descriptor comparison includes modification and change timestamps and re-walks the directory path. The GC dry run shows why activation still requires the #667 gate.

## 9. Acceptance criteria

1. Apply Appendix P to base with the exact executor patch invocation. Snapshot and validate its cumulative manifest. Confirm the expected filtered and unfiltered counts.
2. Publish exact Codex and Claude bundles only after signed admission and reservation. Confirm the exact independent readback and observed step, including every tested failure boundary.
3. Execute the exact offered remedy via the QUICKSTART example. Read back an absent, partial, foreign, changed or exact set and enforce the independent transactional action and step comparisons. Preserve all previous evidence.
4. Run the named socket-free acceptance and inventories, then the same-rig broad gates, brand-new Linux vet and native Ubuntu checks. Preserve the caller fence for §2.2/§5.5 and retain the #667 activation gate.

## Appendix C — after this sprint (measured)

- **Optional M2 store pins.** `host/store/operator_worktree.go:59` heartbeat lower bound, `:56–57` start/session absence, `:89` reserved phase, `:90` other-kind read, and `:138` other-kind disposal. Dropping the read phase predicate survives at the current base as well. On the pristine export with just that substitution, `go test -count=1 -timeout 2m ./host/store/` gave rc 0 and 27.954 s (`$S/pin-phase-probe.log`). No pin patch added. M1 already uses 1650 lines and exceeds the normal cap by 250, leaving zero measured capacity. The related production worktree reservation path has zero entry callers.
- **Carried pins.** `E-RES-EPISODE-CONFLICT`, `E-RES-REREAD`, `J-OWN-UID`, `Z-ABSENT-ANY`, `Z-WAIT-READBUDGET` and the four process pins `TestOperatorFirstManagedEpisodeProcess`, `TestOperatorOldBuildTerminalDrain`, `TestOperatorCurrentOwnerProcessRecovery`, `TestOperatorServiceUpgrade/active_child`.
- **S5f/S6 paired consumer and drain.** Prior FD harness was 72 changed lines (S5e2 Appendix C). Full implementation is UNMEASURED. S6 890, S8 860 and S12 700 are inherited raw design counts. Preserve fresh launch session, bounded processes, old-binary drain, browser lifecycle, interactive authorization and current owner recovery.
- **#667 GC activation gate.** `/bin/bash "$S/gc-probe.sh"` ran the inherited probe against a copy of rebased M1. `$S/probe-M1-GC.log` reports both reserved and fully published worktrees as `clean & merged` dry-run candidates with unchanged evidence at an injected 73 h age. The fixture stubs `IsActive` false, so it establishes candidate eligibility under that precondition. The picture is unchanged. Add an operator-step-aware skip before a production caller publishes files or creates managed worktrees. `git grep -nE "func RunEpisodeGC|OlderThan|IsActive|clean & merged" -- host/daemon/auto_gc.go` gives the 72 h default at :46, entry :101, activity and age gates at :129,136 and candidate at :165. GC change size is UNMEASURED. `$S/callers-M1.log` records base-identical GC bytes.
- **Row 271 and rig flake.** Whole-suite baselines only. Use the historical sets for comparison and the newly recorded sets for mutation and product attribution.

## Appendix P — compiled patches

### Patch M1

````diff M1
diff --git a/cmd/ailang-worldd/episode.go b/cmd/ailang-worldd/episode.go
index 325aff7e..005d59df 100644
--- a/cmd/ailang-worldd/episode.go
+++ b/cmd/ailang-worldd/episode.go
@@ -561,10 +561,15 @@ func serveLine(r *mintRequest, addr string) string {
 
 // writeMCPConfig writes the claude MCP config once, mode 0600: arms.py's
 // claude_mcp_config, the one place besides the session file the token lives.
-func writeMCPConfig(path, addr, token string) error {
+func mcpConfigBytes(addr, token string) ([]byte, error) {
 	body, err := json.Marshal(map[string]any{"mcpServers": map[string]any{worldMCPServer: map[string]any{
 		"type": "http", "url": "http://" + addr + "/mcp/",
 		"headers": map[string]string{"Authorization": "Bearer " + token}}}})
+	return append(body, '\n'), err
+}
+
+func writeMCPConfig(path, addr, token string) error {
+	body, err := mcpConfigBytes(addr, token)
 	if err != nil {
 		return err
 	}
@@ -572,7 +577,7 @@ func writeMCPConfig(path, addr, token string) error {
 	if err != nil {
 		return fmt.Errorf("write the MCP config %s: %w", path, errors.Unwrap(err))
 	}
-	if _, err := f.Write(append(body, '\n')); err != nil {
+	if _, err := f.Write(body); err != nil {
 		_ = f.Close()
 		_ = os.Remove(path)
 		return fmt.Errorf("write the MCP config %s: %w", path, err)
diff --git a/cmd/ailang-worldd/help_test.go b/cmd/ailang-worldd/help_test.go
index d0ae5963..062d1a97 100644
--- a/cmd/ailang-worldd/help_test.go
+++ b/cmd/ailang-worldd/help_test.go
@@ -13,7 +13,7 @@ var everyVerb = [][]string{
 	{"tools", "list"}, {"call"}, {"why"}, {"provenance"}, {"setup"}, {"doctor"}, {"hosts", "init"},
 	{"session"}, {"session", "new"}, {"session", "list"}, {"session", "revoke"}, {"session", "mint"},
 	{"approve"}, {"service", "runner"}, {"service", "runner", "install"}, {"service", "runner", "start"},
-	{"service", "runner", "stop"}, {"service", "runner", "status"}, {"service", "runner", "cancel"}, {"service", "runner", "reconcile-worktree"}, {"service", "runner", "uninstall"}, {"approver"}, {"approver", "list"}, {"approver", "enroll"}, {"approver", "revoke"},
+	{"service", "runner", "stop"}, {"service", "runner", "status"}, {"service", "runner", "cancel"}, {"service", "runner", "reconcile-worktree"}, {"service", "runner", "reconcile-files"}, {"service", "runner", "uninstall"}, {"approver"}, {"approver", "list"}, {"approver", "enroll"}, {"approver", "revoke"},
 }
 
 // AC5.7: `<verb> --help` and `help <verb>` exit 0 with that verb's help on
@@ -42,7 +42,7 @@ func TestEveryVerbHasHelp(t *testing.T) {
 		"serve --db", "] health", "] head", "world get", "object get", "object find", "log get", "log range",
 		"log tail", "registry get", "commit --file <commit.json> [--session <file|token>]",
 		"tools list", "] call <tool>", "] why ", "provenance", "setup", "doctor",
-		"service runner install|start|stop|status|cancel|reconcile-mint|reconcile-worktree|uninstall", "session new <episode>", "session list", "session revoke", "session mint", "help [<verb>]",
+		"service runner install|start|stop|status|cancel|reconcile-mint|reconcile-worktree|reconcile-files|uninstall", "session new <episode>", "session list", "session revoke", "session mint", "help [<verb>]",
 		"approve <request-id>", "approve <request-id> --signature <file>", "approve <request-id> --sign-with <key>", "approve <request-id> --write-message <file>", "approve <request-id> --webauthn-link", "approve <request-id> --webauthn-assertion <file>", "approve <request-id> --await [--await-for <duration>]", "approve --list [--decided]", "approve <request-ref> [--deny] [--daemon <host:port>]", "approver enroll --name", "approver enroll --name <name> --format webauthn --origin <origin> --registration <file>", "approver list", "approver revoke <name>",
 	} {
 		if !strings.Contains(out.String(), want) {
diff --git a/cmd/ailang-worldd/main.go b/cmd/ailang-worldd/main.go
index b3eaf65e..5fd70aee 100644
--- a/cmd/ailang-worldd/main.go
+++ b/cmd/ailang-worldd/main.go
@@ -83,7 +83,7 @@ Usage:
   ailang-worldd hosts init [--out-dir <dir>] [--addr <url>] [--replace]
   ailang-worldd version [--json]
   ailang-worldd service install|status|upgrade|stop|start|uninstall [flags; upgrade: --binary FILE --pause]
-  ailang-worldd service runner install|start|stop|status|cancel|reconcile-mint|reconcile-worktree|uninstall
+  ailang-worldd service runner install|start|stop|status|cancel|reconcile-mint|reconcile-worktree|reconcile-files|uninstall
   ailang-worldd serve --db <path> [--credentials-db <path>] [--bind host:port] [--ailang-bin <path>]
                       [--workbench-still]
                       [--workspace-root <dir> --tool-ailang-bin <path>]
diff --git a/cmd/ailang-worldd/operator_file_syscalls_test.go b/cmd/ailang-worldd/operator_file_syscalls_test.go
new file mode 100644
index 00000000..8efb8455
--- /dev/null
+++ b/cmd/ailang-worldd/operator_file_syscalls_test.go
@@ -0,0 +1,203 @@
+package main
+
+import (
+	"errors"
+	"os"
+	"path/filepath"
+	"reflect"
+	"testing"
+	"time"
+
+	"golang.org/x/sys/unix"
+)
+
+func TestOperatorFileSyscalls(t *testing.T) {
+	root, err := filepath.EvalSymlinks(t.TempDir())
+	if err != nil {
+		t.Fatal(err)
+	}
+	base := filepath.Join(root, "sessions-fixture")
+	if err := os.Mkdir(base, 0700); err != nil {
+		t.Fatal(err)
+	}
+	set := operatorFileSet{Dir: base, UID: uint64(os.Geteuid())}
+	dir, _ := openOperatorFilesDir(set)
+	if dir == nil {
+		t.Fatal("cannot open private fixture")
+	}
+	defer dir.Close()
+	body := []byte(operatorStartTask)
+	cloexec := func(t *testing.T, fd int) {
+		t.Helper()
+		got, err := unix.FcntlInt(uintptr(fd), unix.F_GETFD, 0)
+		if err != nil || got&unix.FD_CLOEXEC == 0 {
+			t.Errorf("descriptor can leak into a child: %d %v", got, err)
+		}
+	}
+	if _, err := installOperatorFileAt(dir, "task", body, func(fd int, name string, flags int, mode uint32) (int, error) {
+		f, err := unix.Openat(fd, name, flags, mode)
+		if err == nil {
+			cloexec(t, f)
+		}
+		return f, err
+	}); err != nil {
+		t.Fatal(err)
+	}
+	for alias, target := range map[string]string{"link": "task", "dirlink": filepath.Base(base)} {
+		parent := base
+		if alias == "dirlink" {
+			parent = root
+		}
+		if err := os.Symlink(target, filepath.Join(parent, alias)); err != nil {
+			t.Fatal(err)
+		}
+	}
+	before := launchFileSnapshot(t, root)
+	refuses := func(t *testing.T, fd int, name string, flags int) {
+		t.Helper()
+		if got, err := unix.Openat(fd, name, flags, 0); err == nil {
+			unix.Close(got)
+			t.Errorf("kernel opened %s with publication flags", name)
+		}
+	}
+	t.Run("directory_open", func(t *testing.T) {
+		probed := false
+		f, p := openOperatorFilesDirAt(set, func(fd int, name string, flags int, mode uint32) (int, error) {
+			if name == filepath.Base(base) && !probed {
+				probed = true
+				refuses(t, fd, "dirlink", flags&^unix.O_DIRECTORY)
+				refuses(t, int(dir.Fd()), "task", flags)
+			}
+			child, err := unix.Openat(fd, name, flags, mode)
+			if err == nil {
+				cloexec(t, child)
+			}
+			return child, err
+		})
+		if f == nil || !probed || p.Status != "private directory" {
+			t.Fatalf("directory positive control failed: %+v", p)
+		}
+		f.Close()
+	})
+	t.Run("read_open", func(t *testing.T) {
+		opened, stats := 0, 0
+		p, data := readOperatorFileAt(dir, "task", serviceHash(body), set.UID,
+			func(fd int, name string, flags int, mode uint32) (int, error) {
+				opened++
+				refuses(t, fd, "link", flags)
+				raw, err := unix.Openat(fd, name, flags, mode)
+				if err == nil {
+					cloexec(t, raw)
+					got, err := unix.FcntlInt(uintptr(raw), unix.F_GETFL, 0)
+					if err != nil || got&unix.O_NONBLOCK == 0 {
+						t.Errorf("blocking read-open: %d %v", got, err)
+					}
+				}
+				return raw, err
+			}, func(fd int, name string, st *unix.Stat_t, flags int) error {
+				stats++
+				var link unix.Stat_t
+				if err := unix.Fstatat(fd, "link", &link, flags); err != nil || link.Mode&unix.S_IFMT != unix.S_IFLNK {
+					t.Errorf("stat %d followed the alias: %o %v", stats, link.Mode, err)
+				}
+				return unix.Fstatat(fd, name, st, flags)
+			})
+		if p.Status != "exact" || string(data) != string(body) || opened != 1 || stats != 2 {
+			t.Fatalf("independent readback did not run: %s %d/%d", p.Status, opened, stats)
+		}
+	})
+	t.Run("create_alias", func(t *testing.T) {
+		// O_CREAT|O_EXCL itself refuses a symlink to an existing file.
+		if _, err := installOperatorFile(dir, "link", []byte("foreign")); !errors.Is(err, os.ErrExist) {
+			t.Fatalf("exclusive create did not refuse an alias: %v", err)
+		}
+	})
+	if !reflect.DeepEqual(before, launchFileSnapshot(t, root)) {
+		t.Fatal("syscall checks changed fixture evidence")
+	}
+}
+
+func TestOperatorFileLockFlags(t *testing.T) {
+	path := filepath.Join(t.TempDir(), "lock")
+	hold, err := openRunnerLock(path, os.O_CREATE)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer hold.Close()
+	other, err := os.Open(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer other.Close()
+	if err := unix.Flock(int(other.Fd()), unix.LOCK_SH|unix.LOCK_NB); !errors.Is(err, unix.EWOULDBLOCK) {
+		t.Fatalf("the episode lock allows a shared entrant: %v", err)
+	}
+	out := make(chan error, 1)
+	go func() {
+		second, err := openRunnerLock(path, 0)
+		if second != nil {
+			second.Close()
+		}
+		out <- err
+	}()
+	select {
+	case err := <-out:
+		if !errors.Is(err, unix.EWOULDBLOCK) {
+			t.Fatalf("unexpected competing lock result: %v", err)
+		}
+	case <-time.After(time.Second):
+		hold.Close()
+		t.Fatal("episode lock blocked instead of refusing")
+	}
+}
+
+func TestOperatorFileReopen(t *testing.T) {
+	for _, phase := range []string{"open", "last_stat", "unreadable"} {
+		t.Run(phase, func(t *testing.T) {
+			base, err := filepath.EvalSymlinks(t.TempDir())
+			if err != nil {
+				t.Fatal(err)
+			}
+			if err := os.Chmod(base, 0700); err != nil {
+				t.Fatal(err)
+			}
+			d, _ := openOperatorFilesDir(operatorFileSet{Dir: base, UID: uint64(os.Geteuid())})
+			if d == nil {
+				t.Fatal("missing private directory")
+			}
+			defer d.Close()
+			body := []byte(operatorStartTask)
+			for _, name := range []string{"task", "replacement"} {
+				if _, err := installOperatorFile(d, name, body); err != nil {
+					t.Fatal(err)
+				}
+			}
+			replace := func(fd int, name string) {
+				t.Helper()
+				if err := unix.Renameat(fd, "replacement", fd, name); err != nil {
+					t.Fatal(err)
+				}
+			}
+			stats := 0
+			p, got := readOperatorFileAt(d, "task", serviceHash(body), uint64(os.Geteuid()),
+				func(fd int, name string, flags int, mode uint32) (int, error) {
+					switch phase {
+					case "open":
+						return unix.Openat(fd, "replacement", flags, mode)
+					case "unreadable":
+						return -1, unix.EACCES
+					}
+					return unix.Openat(fd, name, flags, mode)
+				}, func(fd int, name string, st *unix.Stat_t, flags int) error {
+					stats++
+					if phase == "last_stat" && stats == 2 {
+						replace(fd, name)
+					}
+					return unix.Fstatat(fd, name, st, flags)
+				})
+			if p.Status != "unattributed file" || got != nil || string(readStartFixtureFile(t, filepath.Join(base, "task"))) != string(body) {
+				t.Fatalf("unattributed file received a proof or changed: %s", p.Status)
+			}
+		})
+	}
+}
diff --git a/cmd/ailang-worldd/operator_files.go b/cmd/ailang-worldd/operator_files.go
new file mode 100644
index 00000000..1b4e0f38
--- /dev/null
+++ b/cmd/ailang-worldd/operator_files.go
@@ -0,0 +1,427 @@
+package main
+
+import (
+	"context"
+	"encoding/hex"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"io"
+	"os"
+	"path/filepath"
+	"strings"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"golang.org/x/sys/unix"
+)
+
+type operatorFileID struct {
+	Dev, Inode, Mode, UID, Links uint64
+	Size                         int64
+	Modified, Changed            unix.Timespec
+}
+
+func operatorFileStat(s unix.Stat_t) operatorFileID {
+	return operatorFileID{uint64(s.Dev), s.Ino, uint64(s.Mode), uint64(s.Uid), uint64(s.Nlink), s.Size, s.Mtim, s.Ctim}
+}
+
+type operatorFileSet struct {
+	Dir           string
+	Device, Inode uint64
+	UID           uint64
+	Names, Hashes [3]string
+}
+
+type operatorFileProof struct {
+	ID     operatorFileID
+	Status string
+	Digest string
+}
+
+type operatorFilesProof struct {
+	Dir   operatorFileProof
+	Files [3]operatorFileProof
+}
+
+func operatorFileFact(a approval.OperatorAction, name string) string {
+	for _, fact := range a.Resolution {
+		if fact.Name == name {
+			return fact.Value
+		}
+	}
+	return ""
+}
+
+func operatorFilePaths(s store.OperatorFileState, a approval.OperatorAction) (operatorFileSet, store.EpisodeStart, error) {
+	var set operatorFileSet
+	var start store.EpisodeStart
+	mcp := operatorFileFact(a, "mcp_path")
+	if err := json.Unmarshal([]byte(s.StartJSON), &start); err != nil || start.Task == nil || start.Allowance != nil {
+		return set, start, store.ErrOperatorConflict
+	}
+	token, err := operatorPhysicalTarget(s.Mint.TokenPath)
+	if err != nil {
+		return set, start, err
+	}
+	set.Dir = filepath.Dir(token)
+	set.Names = [3]string{a.Episode + ".task.md", a.Episode + ".session", a.Episode + ".mcp.json"}
+	if filepath.Base(token) != set.Names[1] || filepath.Dir(s.Mint.TokenPath) != set.Dir ||
+		start.Task.Path != filepath.Join(filepath.Dir(s.Mint.TokenPath), set.Names[0]) ||
+		start.Task.SHA256 != serviceHash([]byte(a.Task)) ||
+		(mcp != "") != (start.Agent == "claude") ||
+		mcp != "" && mcp != filepath.Join(filepath.Dir(s.Mint.TokenPath), set.Names[2]) {
+		return set, start, store.ErrOperatorConflict
+	}
+	if mcp == "" {
+		set.Names[2] = ""
+	}
+	set.Hashes[0], set.Hashes[1] = start.Task.SHA256, "sha256:"+s.Mint.CredentialID
+	set.UID = uint64(os.Geteuid())
+	return set, start, nil
+}
+
+type operatorOpenat func(int, string, int, uint32) (int, error)
+
+func openOperatorFilesDir(set operatorFileSet) (*os.File, operatorFileProof) {
+	return openOperatorFilesDirAt(set, unix.Openat)
+}
+
+func openOperatorFilesDirAt(set operatorFileSet, open operatorOpenat) (*os.File, operatorFileProof) {
+	var st unix.Stat_t
+	var p operatorFileProof
+	if err := unix.Lstat(set.Dir, &st); err != nil {
+		p.Status = err.Error()
+		return nil, p
+	}
+	p.ID = operatorFileStat(st)
+	p.Status = "unattributed directory"
+	if p.ID.Mode != unix.S_IFDIR|0700 || p.ID.UID != set.UID ||
+		set.Inode != 0 && (p.ID.Dev != set.Device || p.ID.Inode != set.Inode) {
+		return nil, p
+	}
+	const flags = unix.O_RDONLY | unix.O_DIRECTORY | unix.O_NOFOLLOW | unix.O_CLOEXEC
+	fd, err := unix.Open("/", flags, 0)
+	if err != nil {
+		return nil, p
+	}
+	for _, part := range strings.Split(strings.TrimPrefix(set.Dir, "/"), "/") {
+		next, err := open(fd, part, flags, 0)
+		unix.Close(fd)
+		if err != nil {
+			return nil, p
+		}
+		fd = next
+	}
+	f := os.NewFile(uintptr(fd), set.Dir)
+	if err := unix.Fstat(int(f.Fd()), &st); err != nil || operatorFileStat(st) != p.ID {
+		f.Close()
+		return nil, p
+	}
+	p.Status = "private directory"
+	return f, p
+}
+
+func readOperatorFile(dir *os.File, name, want string, uid uint64) (operatorFileProof, []byte) {
+	return readOperatorFileAt(dir, name, want, uid, unix.Openat, unix.Fstatat)
+}
+
+func readOperatorFileAt(dir *os.File, name, want string, uid uint64, open operatorOpenat, stat func(int, string, *unix.Stat_t, int) error) (operatorFileProof, []byte) {
+	var st unix.Stat_t
+	fd := int(dir.Fd())
+	err := stat(fd, name, &st, unix.AT_SYMLINK_NOFOLLOW)
+	if errors.Is(err, os.ErrNotExist) {
+		return operatorFileProof{Status: "absent"}, nil
+	}
+	if err != nil {
+		return operatorFileProof{Status: err.Error()}, nil
+	}
+	p := operatorFileProof{ID: operatorFileStat(st), Status: "unattributed file"}
+	if p.ID.Mode != unix.S_IFREG|0600 || p.ID.UID != uid || p.ID.Links != 1 {
+		return p, nil
+	}
+	raw, err := open(fd, name, unix.O_RDONLY|unix.O_NOFOLLOW|unix.O_NONBLOCK|unix.O_CLOEXEC, 0)
+	if err != nil {
+		return p, nil
+	}
+	f := os.NewFile(uintptr(raw), name)
+	defer f.Close()
+	if unix.Fstat(raw, &st) != nil || operatorFileStat(st) != p.ID {
+		return p, nil
+	}
+	const limit = approval.MaxOperatorTaskBytes + 4096
+	body, err := io.ReadAll(io.LimitReader(f, limit+1))
+	if err != nil || len(body) > limit || stat(fd, name, &st, unix.AT_SYMLINK_NOFOLLOW) != nil || operatorFileStat(st) != p.ID {
+		return p, nil
+	}
+	p.Digest = serviceHash(body)
+	if p.Digest == want {
+		p.Status = "exact"
+		return p, body
+	}
+	return p, nil
+}
+
+func readOperatorFileSet(set operatorFileSet, addr string) operatorFilesProof {
+	f, dir := openOperatorFilesDir(set)
+	p := operatorFilesProof{Dir: dir}
+	if f == nil {
+		return p
+	}
+	defer f.Close()
+	var token []byte
+	for i, name := range set.Names {
+		if name != "" {
+			var body []byte
+			p.Files[i], body = readOperatorFile(f, name, set.Hashes[i], set.UID)
+			if i == 1 {
+				token = body
+			}
+		}
+	}
+	if p.Files[2].Status == "exact" {
+		body, err := mcpConfigBytes(addr, string(token))
+		if token == nil || err != nil || serviceHash(body) != set.Hashes[2] {
+			p.Files[2].Status = "unattributed file"
+		}
+	}
+	current, reopened := openOperatorFilesDir(set)
+	if current != nil {
+		current.Close()
+	}
+	if reopened != dir {
+		p.Dir.Status = "directory changed during readback"
+	}
+	return p
+}
+
+func installOperatorFile(dir *os.File, name string, data []byte) (operatorFileID, error) {
+	return installOperatorFileAt(dir, name, data, unix.Openat)
+}
+
+func installOperatorFileAt(dir *os.File, name string, data []byte, open operatorOpenat) (operatorFileID, error) {
+	raw, err := open(int(dir.Fd()), name, unix.O_WRONLY|unix.O_CREAT|unix.O_EXCL|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0600)
+	if err != nil {
+		return operatorFileID{}, err
+	}
+	f := os.NewFile(uintptr(raw), name)
+	defer f.Close()
+	if err := f.Chmod(0600); err != nil {
+		return operatorFileID{}, err
+	}
+	if _, err := f.Write(data); err != nil {
+		return operatorFileID{}, err
+	}
+	if err := f.Sync(); err != nil {
+		return operatorFileID{}, err
+	}
+	var st unix.Stat_t
+	err = unix.Fstat(raw, &st)
+	return operatorFileStat(st), err
+}
+
+func operatorFileAddress(start store.EpisodeStart, signed, approved string) (string, error) {
+	var preset launchPreset
+	if err := json.Unmarshal([]byte(start.Preset.Bytes), &preset); err != nil {
+		return "", store.ErrOperatorConflict
+	}
+	d, err := resolvePresetDaemon(preset.Daemon)
+	if err != nil || d.Addr != signed || signed != approved {
+		return "", store.ErrOperatorConflict
+	}
+	return d.Addr, nil
+}
+
+func recheckFileApproval(ctx context.Context, credentials string, a approval.OperatorAction, now int64) (*approval.OperatorMintProof, error) {
+	db, err := store.OpenCredentialsReadOnly(runnerDB(credentials))
+	if err != nil {
+		return nil, err
+	}
+	defer db.Close()
+	row, err := db.OperatorAction(ctx, a.ID, now)
+	if err != nil {
+		return nil, err
+	}
+	if now >= row.ExpiresAt {
+		return nil, store.ErrOperatorExpired
+	}
+	if len(row.Challenge) != 32 {
+		return nil, approval.ErrNotVerified
+	}
+	proof, err := authority.RecheckClaimedOperatorAction(ctx, db, a, [32]byte(row.Challenge), row.ExpiresAt, now, row)
+	if err != nil {
+		return nil, err
+	}
+	return proof.OperatorConfirmation().MintProof(), nil
+}
+
+// publishOperatorStartFiles is dormant until the paired-consumer gates open.
+// rawToken is the existing mint output. Neither it nor the MCP bytes enter SQL.
+func publishOperatorStartFiles(ctx context.Context, c runnerConfig, id, rawToken string, clock func() int64, checkpoint func(string) error) error {
+	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
+	defer cancel()
+	s, a, err := readRunnerFiles(ctx, c.Credentials, id)
+	if err != nil {
+		return err
+	}
+	lock, err := openRunnerLock(operatorEpisodeLockPath(c.Credentials, a.Episode), 0)
+	if err != nil {
+		return err
+	}
+	defer lock.Close()
+	current, a, err := readRunnerFiles(ctx, c.Credentials, id)
+	if err != nil || current != s || s.Phase != "" {
+		return store.ErrOperatorConflict
+	}
+	set, start, err := operatorFilePaths(s, a)
+	if err != nil || serviceHash([]byte(rawToken)) != set.Hashes[1] {
+		return store.ErrOperatorConflict
+	}
+	addr, err := operatorFileAddress(start, operatorFileFact(a, "served_address"), c.Address)
+	if err != nil {
+		return err
+	}
+	data := [3][]byte{[]byte(a.Task), []byte(rawToken), nil}
+	if set.Names[2] != "" {
+		data[2], err = mcpConfigBytes(addr, rawToken)
+		if err != nil {
+			return err
+		}
+		set.Hashes[2] = serviceHash(data[2])
+	}
+	var target operatorWorktreeTarget
+	path, err := operatorPhysicalTarget(operatorWorktreeFact(a))
+	if err != nil || json.Unmarshal([]byte(s.Worktree), &target) != nil ||
+		target.Action != id || target.Path != path || target.Branch != start.Branch || target.SHA != start.Base.SHA {
+		return store.ErrOperatorConflict
+	}
+	git, err := sessionGit()
+	if err != nil {
+		return err
+	}
+	repo, err := filepath.EvalSymlinks(start.Repo)
+	if err != nil || repo != target.Repo {
+		return store.ErrOperatorConflict
+	}
+	worktreeBefore, err := target.inspect(ctx, git)
+	if err != nil || worktreeBefore.state != treeExact {
+		return store.ErrOperatorConflict
+	}
+	dir, proof := openOperatorFilesDir(set)
+	if dir == nil {
+		return store.ErrOperatorConflict
+	}
+	defer dir.Close()
+	set.Device, set.Inode = proof.ID.Dev, proof.ID.Inode
+	for i, name := range set.Names {
+		if name != "" {
+			file, _ := readOperatorFile(dir, name, set.Hashes[i], set.UID)
+			if file.Status != "absent" {
+				return store.ErrOperatorConflict
+			}
+		}
+	}
+	if err := checkpoint("before reservation"); err != nil {
+		return err
+	}
+	approvalProof, err := recheckFileApproval(ctx, c.Credentials, a, clock())
+	if err != nil {
+		return err
+	}
+	record := func(value string) error {
+		return runnerWrite(c.Credentials, func(db *store.Credentials) error {
+			return db.RecordOperatorFiles(ctx, s, value, approvalProof, clock)
+		})
+	}
+	wire, _ := json.Marshal(set) // scalar identity, no raw credential
+	if err := record(string(wire)); err != nil {
+		return err
+	}
+	s.Step.Identity, s.Step.Version, s.Phase = string(wire), 1, "reserved"
+	remedy := reconcileFilesCommand(id)
+	outcome := func(err error) error { return fmt.Errorf("%w; %s", err, remedy) }
+	uncertain := func() error {
+		if err := record(store.OperatorFilesUncertain); err != nil {
+			return outcome(err)
+		}
+		return outcome(errors.New("file publication uncertain"))
+	}
+	var installed [3]operatorFileID
+	for i, name := range set.Names {
+		if name == "" {
+			continue
+		}
+		if err := checkpoint("before " + name); err != nil {
+			return outcome(err)
+		}
+		current, action, err := readRunnerFiles(ctx, c.Credentials, id)
+		now := clock()
+		if err != nil || current != s || now >= current.Session.ExpiresAt {
+			return uncertain()
+		}
+		if _, err := recheckFileApproval(ctx, c.Credentials, action, now); err != nil {
+			return uncertain()
+		}
+		fresh, _ := openOperatorFilesDir(set)
+		if fresh == nil {
+			return uncertain()
+		}
+		fresh.Close()
+		installed[i], err = installOperatorFile(dir, name, data[i])
+		if err != nil {
+			return uncertain()
+		}
+	}
+	if err := dir.Sync(); err != nil {
+		return uncertain()
+	}
+	if err := checkpoint("before readback"); err != nil {
+		return outcome(err)
+	}
+	got := readOperatorFileSet(set, addr)
+	if got.Dir.Status != "private directory" {
+		return uncertain()
+	}
+	for i, name := range set.Names {
+		if name != "" && (got.Files[i].Status != "exact" || got.Files[i].ID != installed[i]) {
+			return uncertain()
+		}
+	}
+	if err := checkpoint("before observe"); err != nil {
+		return outcome(err)
+	}
+	worktreeAfter, err := target.inspect(ctx, git)
+	if err != nil || worktreeAfter != worktreeBefore || readOperatorFileSet(set, addr) != got {
+		return uncertain()
+	}
+	approvalProof, err = recheckFileApproval(ctx, c.Credentials, a, clock())
+	if err != nil {
+		return uncertain()
+	}
+	if err := record(store.OperatorFilesPublished); err != nil {
+		return outcome(err)
+	}
+	return nil
+}
+
+func readRunnerFiles(ctx context.Context, creds, id string) (store.OperatorFileState, approval.OperatorAction, error) {
+	db, err := store.OpenCredentialsReadOnly(runnerDB(creds))
+	if err != nil {
+		return store.OperatorFileState{}, approval.OperatorAction{}, err
+	}
+	defer db.Close()
+	read, cancel := context.WithTimeout(ctx, time.Second)
+	defer cancel()
+	s, a, err := db.ReadOperatorFiles(read, id)
+	if err == nil && (s.Mint.Store != creds || ctx.Err() != nil) {
+		err = store.ErrOperatorConflict
+	}
+	return s, a, err
+}
+
+func operatorFileDigest(s string) bool {
+	b, err := hex.DecodeString(strings.TrimPrefix(s, "sha256:"))
+	return strings.HasPrefix(s, "sha256:") && err == nil && len(b) == 32
+}
diff --git a/cmd/ailang-worldd/operator_files_reconcile.go b/cmd/ailang-worldd/operator_files_reconcile.go
new file mode 100644
index 00000000..5d4abed2
--- /dev/null
+++ b/cmd/ailang-worldd/operator_files_reconcile.go
@@ -0,0 +1,92 @@
+package main
+
+import (
+	"encoding/json"
+	"errors"
+	"fmt"
+	"io"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func reconcileFilesCommand(id string) string {
+	return "ailang-worldd service runner reconcile-files " + id
+}
+
+func (e serviceEnv) reconcileRunnerFiles(id string, out, stderr io.Writer) int {
+	remedy := reconcileFilesCommand(id)
+	refuse := func(err error) int {
+		return serviceMessage(stderr, "file recovery refused", err.Error(), remedy, exitUsage)
+	}
+	c, _, err := e.installedRunner()
+	if err != nil {
+		return refuse(err)
+	}
+	ctx, cancel := budgetContext(5 * time.Second)
+	defer cancel()
+	s, a, err := readRunnerFiles(ctx, c.Credentials, id)
+	if err != nil {
+		return refuse(err)
+	}
+	if s.Phase != "reserved" {
+		return refuse(store.ErrOperatorMissing)
+	}
+	lock, err := openRunnerLock(operatorEpisodeLockPath(c.Credentials, a.Episode), 0)
+	if err != nil {
+		return refuse(errors.New("episode recovery lock held or unavailable"))
+	}
+	defer lock.Close()
+	current, a, err := readRunnerFiles(ctx, c.Credentials, id)
+	if err != nil || current != s {
+		return refuse(store.ErrOperatorConflict)
+	}
+	want, start, err := operatorFilePaths(s, a)
+	var set operatorFileSet
+	if err != nil || json.Unmarshal([]byte(s.Step.Identity), &set) != nil || set.Dir != want.Dir ||
+		set.UID != want.UID || set.Inode == 0 || set.Names != want.Names ||
+		set.Hashes[0] != want.Hashes[0] || set.Hashes[1] != want.Hashes[1] ||
+		(set.Names[2] != "" && !operatorFileDigest(set.Hashes[2])) ||
+		(set.Names[2] == "" && set.Hashes[2] != "") {
+		return refuse(store.ErrOperatorConflict)
+	}
+	addr, err := operatorFileAddress(start, operatorFileFact(a, "served_address"), c.Address)
+	if err != nil {
+		return refuse(err)
+	}
+	before := readOperatorFileSet(set, addr)
+	absent, exact := before.Dir.Status == "private directory", before.Dir.Status == "private directory"
+	fmt.Fprintln(stderr, "File directory:", before.Dir.Status)
+	for i, name := range set.Names {
+		if name != "" {
+			fmt.Fprintln(stderr, name, before.Files[i].Status)
+			absent = absent && before.Files[i].Status == "absent"
+			exact = exact && before.Files[i].Status == "exact"
+		}
+	}
+	if !absent && !exact {
+		return refuse(errors.New("file evidence remains uncertain"))
+	}
+	cancel()
+	disposition := store.OperatorFilesStoppedRetained
+	if absent {
+		disposition = store.OperatorFilesStoppedAbsent
+	}
+	if !e.confirm(stderr, fmt.Sprintf("Stop action %s and keep its files, credential and worktree as found?", id)) {
+		fmt.Fprintln(stderr, "Retry at a terminal:", remedy)
+		return exitUsage
+	}
+	if after := readOperatorFileSet(set, addr); after != before {
+		return refuse(errors.New("file evidence changed after review"))
+	}
+	err = runnerWrite(c.Credentials, func(db *store.Credentials) error {
+		ctx, cancel := budgetContext(time.Second)
+		defer cancel()
+		return db.RecordOperatorFiles(ctx, s, disposition, nil, e.terminal.now)
+	})
+	if err != nil {
+		return refuse(err)
+	}
+	fmt.Fprintln(out, "File reservation disposed for", id)
+	return exitOK
+}
diff --git a/cmd/ailang-worldd/operator_files_test.go b/cmd/ailang-worldd/operator_files_test.go
new file mode 100644
index 00000000..c183f736
--- /dev/null
+++ b/cmd/ailang-worldd/operator_files_test.go
@@ -0,0 +1,606 @@
+package main
+
+import (
+	"database/sql"
+	"encoding/json"
+	"errors"
+	"io"
+	"os"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"syscall"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+const filesExample = "ailang-worldd service runner reconcile-files oa-0123456789abcdef0123456789abcdef"
+
+func mintOperatorFileFixture(t *testing.T, r *worktreeStartRig) string {
+	t.Helper()
+	if err := r.start(t, func(string) error { return nil }); err != nil {
+		t.Fatal(err)
+	}
+	r.state(t, "running", "observed", "created", 2, true)
+	ctx := boundedTestContext(t)
+	ro, err := store.OpenCredentialsReadOnly(runnerDB(r.creds))
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer ro.Close()
+	row, err := ro.OperatorAction(ctx, r.action.ID, 1012)
+	if err != nil {
+		t.Fatal(err)
+	}
+	in := r.in
+	in.Attempt = r.at
+	in.Attempt.CAS.Generation = row.Generation
+	in.Attempt.StepVersion = 1
+	var token string
+	err = runnerWrite(r.creds, func(db *store.Credentials) error {
+		key := store.OperatorStepKey{Action: in.Attempt.CAS, Owner: row.Owner, Kind: "mint", Identity: row.ID, Version: 1}
+		if err := db.ReserveOperatorStep(ctx, key, 1012); err != nil {
+			return err
+		}
+		token, err = authority.MintOperatorStart(ctx, db, ro, in, r.action, row, r.src.clock)
+		return err
+	})
+	if err != nil {
+		t.Fatal(err)
+	}
+	if err := os.Mkdir(r.f.r.sessions, 0700); err != nil {
+		t.Fatal(err)
+	}
+	return token
+}
+
+func TestOperatorFilesReconcile(t *testing.T) {
+	t.Run("metadata", TestOperatorFileReadbackMetadata)
+	doc, err := os.ReadFile("../../docs/QUICKSTART.md")
+	_, example, found := strings.Cut(string(doc), "\nailang-worldd service runner reconcile-files ")
+	line, _, _ := strings.Cut(example, "\n")
+	example = "ailang-worldd service runner reconcile-files " + line
+	if err != nil || !found || example != filesExample {
+		t.Fatalf("missing executable QUICKSTART example: %v", err)
+	}
+	for _, tc := range []struct {
+		name, at string
+		count    int
+	}{
+		{"no_mcp", "", 2}, {"all_files", "", 3},
+		{"before_reservation", "before reservation", 0}, {"before_first", "before ep.task.md", 0},
+		{"between_task_token", "before ep.session", 1}, {"between_token_mcp", "before ep.mcp.json", 2},
+		{"before_readback", "before readback", 3}, {"before_observe", "before observe", 3},
+		{"competing_file", "", 2}, {"changed_file", "", 3}, {"replaced_after_readback", "", 3}, {"replaced_before_readback", "", 3},
+		{"changed_during_review", "before observe", 3}, {"parent_changed_during_review", "before observe", 3},
+		{"expired_reopen", "before ep.task.md", 0},
+		{"owner_changed", "before observe", 3}, {"version_changed", "before observe", 3},
+		{"later_step", "before observe", 3}, {"post_step_version", "before observe", 3}, {"observe_error", "", 3},
+		{"publish_owner_drift", "", 3}, {"publish_mint_drift", "", 3}, {"late_key_revoked", "", 3}, {"late_signature_drift", "", 3},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			r := newWorktreeStartRig(t, tc.name != "no_mcp")
+			token := mintOperatorFileFixture(t, r)
+			mcpRef := filepath.Join(t.TempDir(), "reference")
+			if err := writeMCPConfig(mcpRef, r.config.Address, token); err != nil {
+				t.Fatal(err)
+			}
+			names := [3]string{"ep.task.md", "ep.session", "ep.mcp.json"}
+			data := [3]string{r.action.Task, token, string(readStartFixtureFile(t, mcpRef))}
+			ctx := presetTestContext(t)
+			initial := observeProvision(t, r.f.r, r.creds)
+			var baseline, want store.OperatorFileState
+			baseline, _, err := readRunnerFiles(ctx, r.creds, r.action.ID)
+			if err != nil {
+				t.Fatal(err)
+			}
+			step := func(phase, state, receipt string, version int64) {
+				t.Helper()
+				var gotPhase, gotState, gotReceipt, owner string
+				var gotVersion, generation, starts, sessions, terminal int64
+				q := `SELECT coalesce(f.phase,''), a.state, coalesce(f.receipt,''), coalesce(f.version,0),
+					coalesce(a.owner,''), a.generation, coalesce(a.terminal_at,0), (SELECT count(*) FROM episode_starts),
+					(SELECT count(*) FROM session_credentials) FROM operator_actions a
+					LEFT JOIN operator_action_steps f ON f.kind='files' AND f.action_id=a.id`
+				err := r.db.QueryRow(q).Scan(&gotPhase, &gotState, &gotReceipt, &gotVersion, &owner, &generation, &terminal, &starts, &sessions)
+				wantGeneration, wantOwner, wantTerminal := baseline.Mint.CAS.Generation, baseline.Mint.Owner, int64(0)
+				if state == "uncertain" {
+					wantGeneration++
+				}
+				if state == "failed" {
+					wantOwner, wantTerminal, wantGeneration = "", r.env.terminal.now(), wantGeneration+version-1
+				}
+				if err != nil || gotPhase != phase || gotState != state || gotReceipt != receipt || gotVersion != version ||
+					owner != wantOwner || terminal != wantTerminal || generation != wantGeneration || starts != 1 || sessions != 1 {
+					t.Fatalf("unexpected persisted state %s/%s/%s/%d ownership or mint counts: %v", gotPhase, gotState, gotReceipt, gotVersion, err)
+				}
+			}
+			checkFiles := func() {
+				t.Helper()
+				for i, name := range names {
+					p := filepath.Join(r.f.r.sessions, name)
+					st, err := os.Lstat(p)
+					if i >= tc.count {
+						if !errors.Is(err, os.ErrNotExist) {
+							t.Fatalf("unexpected file %s: %v", name, err)
+						}
+						continue
+					}
+					if err != nil || st.Mode() != 0600 || uint64(st.Sys().(*syscall.Stat_t).Nlink) != 1 ||
+						st.Sys().(*syscall.Stat_t).Uid != uint32(os.Geteuid()) || string(readStartFixtureFile(t, p)) != data[i] {
+						t.Fatalf("wrong bytes or metadata for %s: %v", name, err)
+					}
+				}
+				snapshot := observeProvision(t, r.f.r, r.creds)
+				if snapshot.worktrees != initial.worktrees || snapshot.branches != initial.branches ||
+					!reflect.DeepEqual(snapshot.workspace, initial.workspace) || len(r.writes) != 0 || r.ctl.boots != 0 || r.ctl.stops != 0 {
+					t.Fatal("publication or disposition changed Git or the service")
+				}
+			}
+			if tc.name == "all_files" {
+				mask := syscall.Umask(0777)
+				defer syscall.Umask(mask)
+			}
+			crash := errors.New("injected file interruption")
+			err = publishOperatorStartFiles(ctx, r.config, r.action.ID, token, r.src.clock, func(point string) error {
+				if point == "before ep.task.md" {
+					step("reserved", "running", "", 1)
+					want, _, err = readRunnerFiles(ctx, r.creds, r.action.ID)
+					if err != nil || strings.Contains(want.Step.Identity, token) {
+						t.Fatal("reservation missing or stores raw credentials")
+					}
+					var target operatorFileSet
+					if json.Unmarshal([]byte(want.Step.Identity), &target) != nil ||
+						target.Hashes[0] != serviceHash([]byte(data[0])) || target.Hashes[1] != serviceHash([]byte(data[1])) ||
+						(tc.name != "no_mcp" && target.Hashes[2] != serviceHash([]byte(data[2]))) {
+						t.Fatal("reservation does not bind actual file bytes")
+					}
+					lock, err := openRunnerLock(operatorEpisodeLockPath(r.creds, r.action.Episode), 0)
+					if err == nil {
+						lock.Close()
+						t.Fatal("episode lock was released during publication")
+					}
+				}
+				switch {
+				case tc.name == "competing_file" && point == "before ep.session":
+					data[1] = "foreign"
+					seedStartFixtureFile(t, filepath.Join(r.f.r.sessions, names[1]), []byte(data[1]))
+				case tc.name == "changed_file" && point == "before readback":
+					data[0] = strings.Repeat("x", len(r.action.Task))
+					seedStartFixtureFile(t, filepath.Join(r.f.r.sessions, names[0]), []byte(data[0]))
+				case tc.name == "replaced_after_readback" && point == "before observe",
+					tc.name == "replaced_before_readback" && point == "before readback":
+					path := filepath.Join(r.f.r.sessions, names[1])
+					old := filepath.Join(r.f.r.sessions, "retained-original")
+					if err := os.Rename(path, old); err != nil {
+						t.Fatal(err)
+					}
+					seedStartFixtureFile(t, path, readStartFixtureFile(t, old))
+				case point == "before observe":
+					var drift string
+					switch tc.name {
+					case "publish_owner_drift":
+						drift = "UPDATE operator_actions SET owner='other'"
+					case "publish_mint_drift":
+						drift = "UPDATE session_credentials SET expires_at=expires_at+1"
+					case "late_key_revoked":
+						r.exec(t, "UPDATE approvers SET revoked_at=1012 WHERE name='fixture'")
+					case "late_signature_drift":
+						r.exec(t, "UPDATE operator_actions SET artifact=X'01'")
+					case "observe_error":
+						drift = "SELECT RAISE(ABORT,'file observation refused')"
+					}
+					if drift != "" {
+						r.exec(t, "CREATE TRIGGER fail_file_observe AFTER UPDATE OF phase ON operator_action_steps WHEN old.kind='files' BEGIN "+drift+"; END")
+					}
+					if tc.at == point {
+						return crash
+					}
+				case tc.at == point:
+					return crash
+				}
+				return nil
+			})
+			checkFiles()
+			if tc.name == "no_mcp" || tc.name == "all_files" {
+				if err != nil {
+					t.Fatal(err)
+				}
+				step("observed", "running", "published", 2)
+				exact := observeProvision(t, r.f.r, r.creds)
+				if err := publishOperatorStartFiles(ctx, r.config, r.action.ID, token, r.src.clock, func(string) error { return errors.New("unexpected replay") }); !errors.Is(err, store.ErrOperatorConflict) ||
+					!reflect.DeepEqual(exact, observeProvision(t, r.f.r, r.creds)) {
+					t.Fatal("replay modified a published start")
+				}
+				r.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) { t.Error("prompt for observed row"); return nil, crash }
+				if code, msg := r.run(t, example); code != exitUsage || !strings.Contains(msg, store.ErrOperatorMissing.Error()) ||
+					!reflect.DeepEqual(exact, observeProvision(t, r.f.r, r.creds)) {
+					t.Fatalf("observed row reached disposition: %d %s", code, msg)
+				}
+				return
+			}
+			if tc.name == "before_reservation" {
+				step("", "running", "", 0)
+				before := observeProvision(t, r.f.r, r.creds)
+				r.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) { t.Error("prompt without reservation"); return nil, crash }
+				if c, msg := r.run(t, example); !errors.Is(err, crash) || c != exitUsage || !strings.Contains(msg, store.ErrOperatorMissing.Error()) ||
+					!reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
+					t.Fatalf("wrong refusal before reservation: %v %d %s", err, c, msg)
+				}
+				return
+			}
+			version := int64(1)
+			if tc.name == "competing_file" || tc.name == "changed_file" || tc.name == "replaced_after_readback" || tc.name == "replaced_before_readback" || strings.HasPrefix(tc.name, "late_") {
+				version++
+				step("reserved", "uncertain", "", version)
+				if err == nil || !strings.Contains(err.Error(), "file publication uncertain") {
+					t.Fatalf("wrong mismatch outcome: %v", err)
+				}
+			} else {
+				step("reserved", "running", "", version)
+				if tc.name == "observe_error" {
+					if err == nil || !strings.Contains(err.Error(), "file observation refused") {
+						t.Fatalf("wrong database failure: %v", err)
+					}
+					r.exec(t, "DROP TRIGGER fail_file_observe")
+				} else if tc.name == "publish_owner_drift" || tc.name == "publish_mint_drift" {
+					if !errors.Is(err, store.ErrOperatorConflict) {
+						t.Fatalf("wrong post-step CAS refusal: %v", err)
+					}
+					r.exec(t, "DROP TRIGGER fail_file_observe")
+				} else if !errors.Is(err, crash) {
+					t.Fatalf("lost interruption class: %v", err)
+				}
+			}
+			_, remedy, ok := strings.Cut(err.Error(), "; ")
+			if !ok || remedy != filesExample {
+				t.Fatalf("wrong file remedy: %v", err)
+			}
+			before := observeProvision(t, r.f.r, r.creds)
+			if tc.name == "between_task_token" || tc.name == "between_token_mcp" || tc.name == "competing_file" || tc.name == "changed_file" {
+				r.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
+					t.Error("unattributed or partial bundle reached attended stop")
+					return nil, crash
+				}
+				code, msg := r.run(t, example)
+				if code != exitUsage || !strings.Contains(msg, "file evidence remains uncertain") ||
+					!reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
+					t.Fatalf("partial or foreign evidence was disposed: %d %s", code, msg)
+				}
+				return
+			}
+			for _, answer := range []string{"no tty", "n\n", "y\n"} {
+				r.env.terminal = testMintEnv(answer, 1012)
+				if answer == "no tty" {
+					r.env.terminal = noTerminalEnv(1012)
+				}
+				if tc.name == "expired_reopen" {
+					r.env.terminal.now = func() int64 { return baseline.Session.ExpiresAt }
+				}
+				confirm := r.env.terminal.openTerminal
+				if answer == "y\n" {
+					r.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
+						switch tc.name {
+						case "changed_during_review":
+							data[0] += "changed"
+							seedStartFixtureFile(t, filepath.Join(r.f.r.sessions, names[0]), []byte(data[0]))
+						case "parent_changed_during_review":
+							old := r.f.r.sessions + "-retained"
+							if err := os.Rename(r.f.r.sessions, old); err != nil {
+								t.Fatal(err)
+							}
+							if err := os.Mkdir(r.f.r.sessions, 0700); err != nil {
+								t.Fatal(err)
+							}
+							// Same signed bytes and modes in a new parent do not make it the reserved parent.
+							for _, name := range names {
+								seedStartFixtureFile(t, filepath.Join(r.f.r.sessions, name), readStartFixtureFile(t, filepath.Join(old, name)))
+							}
+							files := launchFileSnapshot(t, old)
+							t.Cleanup(func() {
+								if !reflect.DeepEqual(files, launchFileSnapshot(t, old)) {
+									t.Error("disposition changed the original directory")
+								}
+							})
+						case "owner_changed":
+							r.exec(t, "CREATE TRIGGER owner_after_file AFTER UPDATE OF phase ON operator_action_steps WHEN old.kind='files' BEGIN UPDATE operator_actions SET owner='other'; END")
+						case "version_changed":
+							r.exec(t, "UPDATE operator_action_steps SET version=version+1 WHERE kind='files'")
+						case "later_step":
+							r.exec(t, "CREATE TRIGGER later_after_file AFTER UPDATE OF phase ON operator_action_steps WHEN old.kind='files' BEGIN INSERT INTO operator_action_steps(kind,identity,action_id,phase,version,owner,created_at) SELECT 'launch',id,id,'reserved',1,owner,1012 FROM operator_actions; END")
+						case "post_step_version":
+							r.exec(t, "CREATE TRIGGER version_after_file AFTER UPDATE OF phase ON operator_action_steps WHEN old.kind='files' BEGIN UPDATE operator_action_steps SET version=version+1 WHERE kind='files'; END")
+						}
+						before = observeProvision(t, r.f.r, r.creds)
+						return confirm()
+					}
+				}
+				code, msg := r.run(t, example)
+				switch {
+				case answer != "y\n":
+					if code != exitUsage || !strings.Contains(msg, "Retry at a terminal: "+remedy+"\n") ||
+						!reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
+						t.Fatalf("unattended disposition changed state: %d %s", code, msg)
+					}
+				case tc.name == "changed_during_review" || tc.name == "parent_changed_during_review" || tc.name == "owner_changed" || tc.name == "version_changed" || tc.name == "later_step" || tc.name == "post_step_version":
+					if code != exitUsage || !strings.Contains(msg, "file recovery refused") ||
+						!reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
+						t.Fatalf("stale disposition changed evidence: %d %s", code, msg)
+					}
+				default:
+					if code != exitOK || !strings.Contains(msg, "File reservation disposed for "+r.action.ID) {
+						t.Fatalf("attended disposition: %d %s", code, msg)
+					}
+					receipt := "stopped: retained"
+					if tc.count == 0 {
+						receipt = "stopped: absent"
+					}
+					step("observed", "failed", receipt, version+1)
+				}
+			}
+			checkFiles()
+		})
+	}
+}
+
+func TestOperatorFileReadbackMetadata(t *testing.T) {
+	for _, name := range []string{"exact", "digest", "owner", "mode", "symlink", "hardlink", "fifo", "at_limit", "oversize"} {
+		t.Run(name, func(t *testing.T) {
+			base, err := filepath.EvalSymlinks(t.TempDir())
+			if err != nil {
+				t.Fatal(err)
+			}
+			if err := os.Chmod(base, 0700); err != nil {
+				t.Fatal(err)
+			}
+			path, alternate := filepath.Join(base, "file"), filepath.Join(base, "alternate")
+			set := operatorFileSet{Dir: base, UID: uint64(os.Geteuid())}
+			f, _ := openOperatorFilesDir(set)
+			if f == nil {
+				t.Fatal("fixture directory unavailable")
+			}
+			defer f.Close()
+			body := []byte("from publisher")
+			if name == "at_limit" || name == "oversize" {
+				body = make([]byte, approval.MaxOperatorTaskBytes+4096)
+			}
+			installed := body
+			if name == "oversize" {
+				installed = append(body, 1)
+			}
+			if _, err := installOperatorFile(f, "alternate", installed); err != nil {
+				t.Fatal(err)
+			}
+			uid := set.UID
+			want := serviceHash(body)
+			switch name {
+			case "owner":
+				uid++
+				err = os.Rename(alternate, path)
+			case "mode":
+				err = os.Rename(alternate, path)
+				if err == nil {
+					err = os.Chmod(path, 0640)
+				}
+			case "symlink":
+				err = os.Symlink(alternate, path)
+			case "hardlink":
+				err = os.Link(alternate, path)
+			case "fifo":
+				want, err = serviceHash(nil), syscall.Mkfifo(path, 0600)
+			default:
+				err = os.Rename(alternate, path)
+				if name == "digest" {
+					want = serviceHash([]byte("other contents"))
+				}
+			}
+			if err != nil {
+				t.Fatal(err)
+			}
+			before := launchFileSnapshot(t, base)
+			got, _ := readOperatorFile(f, "file", want, uid)
+			if (got.Status == "exact") != (name == "exact" || name == "at_limit") {
+				t.Fatalf("wrong readback class: %s", got.Status)
+			}
+			if _, err := installOperatorFile(f, "file", []byte("replacement")); !errors.Is(err, os.ErrExist) ||
+				!reflect.DeepEqual(before, launchFileSnapshot(t, base)) {
+				t.Fatalf("install accepted existing pathname or changed evidence: %v", err)
+			}
+		})
+	}
+}
+
+func TestOperatorFilePublicationAdmission(t *testing.T) {
+	for _, name := range []string{"missing_mint", "wrong_worktree_receipt", "wrong_worktree_owner", "wrong_worktree_phase", "wrong_mint_owner", "expired_credential",
+		"expired_action", "expiry_after_reservation", "credential_expiry_after_reservation", "wrong_grants", "wrong_ttl", "revoked_approver", "revoked_credential",
+		"wrong_token", "wrong_approved_address", "reserve_owner_drift", "reserve_mint_drift", "future_mint", "existing_exact_file", "linked_destination", "wrong_destination_mode", "worktree_moved", "lock_held"} {
+		t.Run(name, func(t *testing.T) {
+			ttl := int64(43200)
+			if name == "credential_expiry_after_reservation" {
+				ttl = 60
+			}
+			r := newWorktreeStartRig(t, true, ttl)
+			token := mintOperatorFileFixture(t, r)
+			clock := r.src.clock
+			ctx := presetTestContext(t)
+			want := store.ErrOperatorConflict
+			switch name {
+			case "missing_mint":
+				r.exec(t, "DELETE FROM operator_action_steps WHERE kind='mint'")
+				want = nil
+			case "wrong_worktree_receipt":
+				r.exec(t, "UPDATE operator_action_steps SET receipt=? WHERE kind='worktree'", []byte("stopped: retained"))
+				want = nil
+			case "wrong_worktree_owner":
+				r.exec(t, "UPDATE operator_action_steps SET owner='other' WHERE kind='worktree'")
+				want = sql.ErrNoRows
+			case "wrong_worktree_phase":
+				r.exec(t, "UPDATE operator_action_steps SET phase='reserved',receipt=NULL WHERE kind='worktree'")
+				want = sql.ErrNoRows
+			case "wrong_mint_owner":
+				r.exec(t, "UPDATE operator_action_steps SET owner='other' WHERE kind='mint'")
+				want = nil
+			case "expired_credential":
+				r.exec(t, "UPDATE session_credentials SET expires_at=1012")
+				want = store.ErrOperatorExpired
+			case "expired_action":
+				clock = func() int64 { return 1600 }
+				want = store.ErrOperatorExpired
+			case "expiry_after_reservation", "credential_expiry_after_reservation":
+				now := int64(1012)
+				clock = func() int64 { return now }
+				err := publishOperatorStartFiles(ctx, r.config, r.action.ID, token, clock, func(point string) error {
+					if point == "before ep.task.md" {
+						now = 1600
+						if name == "credential_expiry_after_reservation" {
+							now = 1012 + ttl
+						}
+					}
+					return nil
+				})
+				if err == nil || err.Error() != "file publication uncertain; "+filesExample {
+					t.Fatalf("expired reserved publication: %v", err)
+				}
+				st, _, err := readRunnerFiles(ctx, r.creds, r.action.ID)
+				files, listErr := os.ReadDir(r.f.r.sessions)
+				if err != nil || st.Phase != "reserved" || st.Step.Version != 2 || listErr != nil || len(files) != 0 {
+					t.Fatalf("expired reservation changed its files or lost the receipt: %v %v", err, listErr)
+				}
+				return
+			case "wrong_grants":
+				r.exec(t, "UPDATE session_credentials SET grants_json='[]'")
+			case "wrong_ttl":
+				r.exec(t, "UPDATE session_credentials SET expires_at=expires_at+1")
+			case "revoked_approver":
+				r.exec(t, "UPDATE approvers SET revoked_at=1012 WHERE name='fixture'")
+				want = nil
+			case "revoked_credential":
+				r.exec(t, "DELETE FROM session_credentials")
+				want = nil
+			case "wrong_token":
+				token += "other"
+			case "wrong_approved_address":
+				r.config.Address = "127.0.0.1:1"
+			case "reserve_owner_drift":
+				r.exec(t, "CREATE TRIGGER owner_after_file_reserve AFTER INSERT ON operator_action_steps WHEN new.kind='files' BEGIN UPDATE operator_actions SET owner='other'; END")
+			case "reserve_mint_drift":
+				r.exec(t, "CREATE TRIGGER mint_after_file_reserve AFTER INSERT ON operator_action_steps WHEN new.kind='files' BEGIN UPDATE session_credentials SET expires_at=expires_at+1; END")
+			case "future_mint":
+				clock = func() int64 { return 1011 }
+			case "existing_exact_file":
+				seedStartFixtureFile(t, r.in.TokenPath, []byte(token))
+			case "linked_destination":
+				if err := os.Remove(r.f.r.sessions); err != nil {
+					t.Fatal(err)
+				}
+				decoy := t.TempDir()
+				if err := os.Chmod(decoy, 0700); err != nil {
+					t.Fatal(err)
+				}
+				if err := os.Symlink(decoy, r.f.r.sessions); err != nil {
+					t.Fatal(err)
+				}
+				t.Cleanup(func() {
+					if listed, err := os.ReadDir(decoy); err != nil || len(listed) != 0 {
+						t.Errorf("linked directory changed: %v", err)
+					}
+				})
+			case "wrong_destination_mode":
+				if err := os.Chmod(r.f.r.sessions, 0750); err != nil {
+					t.Fatal(err)
+				}
+			case "worktree_moved":
+				presetGit(t, r.f.r.f, "worktree", "move", r.path(), filepath.Join(t.TempDir(), "other"))
+			case "lock_held":
+				lock, err := openRunnerLock(operatorEpisodeLockPath(r.creds, r.action.Episode), 0)
+				if err != nil {
+					t.Fatal(err)
+				}
+				defer lock.Close()
+				want = nil
+			}
+			before := observeProvision(t, r.f.r, r.creds)
+			err := publishOperatorStartFiles(ctx, r.config, r.action.ID, token, clock, func(point string) error {
+				if point != "before reservation" ||
+					(name != "expired_credential" && name != "expired_action" && name != "wrong_grants" &&
+						name != "wrong_ttl" && name != "revoked_approver" && !strings.HasPrefix(name, "reserve_") && name != "future_mint") {
+					t.Errorf("unexpected effect boundary: %s", point)
+				}
+				return nil
+			})
+			if err == nil || want != nil && !errors.Is(err, want) || !reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
+				t.Fatalf("expected refusal with unchanged store, files and worktree: %v", err)
+			}
+			var files int
+			if err := r.db.QueryRow("SELECT count(*) FROM operator_action_steps WHERE kind='files'").Scan(&files); err != nil || files != 0 {
+				t.Fatalf("file step created by refused start: %d %v", files, err)
+			}
+		})
+	}
+}
+
+func TestOperatorFileAllowanceOrigin(t *testing.T) {
+	r := newWorktreeStartRig(t, true)
+	if err := r.start(t, func(string) error { return nil }); err != nil {
+		t.Fatal(err)
+	}
+	start := r.in.Start
+	if err := os.Mkdir(r.f.r.sessions, 0700); err != nil {
+		t.Fatal(err)
+	}
+	ctx := boundedTestContext(t)
+	now := r.src.clock()
+	start.StartedAt = time.Unix(now, 0).UTC().Format(time.RFC3339)
+	al := approval.StandingAllowance{Schema: approval.AllowanceSchema, PresetName: start.Preset.Name,
+		PresetSHA256: strings.TrimPrefix(start.Preset.SHA256, "sha256:"), Profile: r.in.Profile,
+		Repos: []string{start.Repo}, Grants: r.in.Grants, TTLSeconds: start.TTLS, MaxPerDay: 1, ValidDays: 1,
+		WorldDB: r.config.World, CredentialsDB: r.creds, WorkspaceRoot: r.f.r.f.ws, ExcludeMissionLoop: true}
+	payload, err := approval.CanonicalAllowance(al)
+	if err != nil {
+		t.Fatal(err)
+	}
+	var token string
+	err = runnerWrite(r.creds, func(db *store.Credentials) error {
+		if err := db.GrantAllowance(ctx, store.StandingAllowanceRow{AllowanceID: "sa-0123456789abcdef",
+			PresetName: al.PresetName, PresetSHA256: al.PresetSHA256, ProfileProject: al.Profile.Project,
+			ProfileDigest: al.Profile.Digest, Payload: payload, Digest: strings.TrimPrefix(serviceHash(payload), "sha256:"),
+			Challenge: strings.Repeat("a", 64), MaxPerDay: 1, TTLS: start.TTLS, ApprovedAt: now - 1, ExpiresAt: now + 3600, ApprovedBy: "tty"}); err != nil {
+			return err
+		}
+		var grants []broker.Capability
+		for _, g := range r.in.Grants {
+			grants = append(grants, broker.Capability{Effect: g.Effect, Scope: g.Scope, Budget: g.Budget})
+		}
+		token, _, _, _, _, err = authority.MintFromAllowance(ctx, db, start.Episode,
+			authority.ProfileBinding{Project: al.Profile.Project, Digest: al.Profile.Digest}, grants, start.TTLS, now, nil,
+			authority.AllowanceMintRequest{PresetName: al.PresetName, PresetSHA256: al.PresetSHA256, Repo: start.Repo, Start: start})
+		return err
+	})
+	if err != nil || token == "" {
+		t.Fatalf("real allowance draw: %v", err)
+	}
+	var minted store.EpisodeStart
+	var startJSON string
+	if err := r.db.QueryRow("SELECT start_json FROM episode_starts WHERE episode=?", start.Episode).Scan(&startJSON); err != nil ||
+		json.Unmarshal([]byte(startJSON), &minted) != nil || minted.Allowance == nil {
+		t.Fatalf("allowance provenance missing: %v", err)
+	}
+	before := observeProvision(t, r.f.r, r.creds)
+	err = publishOperatorStartFiles(ctx, r.config, r.action.ID, token, r.src.clock, func(point string) error {
+		t.Errorf("allowance entered publication at %s", point)
+		return nil
+	})
+	if !errors.Is(err, sql.ErrNoRows) || !reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
+		t.Fatalf("allowance draw was adopted: %v", err)
+	}
+	var draws, steps int
+	if err := r.db.QueryRow("SELECT (SELECT count(*) FROM allowance_draws), (SELECT count(*) FROM operator_action_steps WHERE kind IN ('mint','files'))").Scan(&draws, &steps); err != nil ||
+		draws != 1 || steps != 0 {
+		t.Fatalf("allowance or stage evidence changed: %d/%d %v", draws, steps, err)
+	}
+}
diff --git a/cmd/ailang-worldd/operator_service.go b/cmd/ailang-worldd/operator_service.go
index 5274a6f1..b47600c4 100644
--- a/cmd/ailang-worldd/operator_service.go
+++ b/cmd/ailang-worldd/operator_service.go
@@ -21,8 +21,8 @@ import (
 
 const runnerLabel = "com.sunholo.ailang-world.operator-runner"
 const runnerSystemPath = "/usr/bin:/bin:/usr/sbin:/sbin"
-const runnerHelp = `usage: ailang-worldd service runner install --config FILE | start | stop | status | cancel ACTION-ID | reconcile-mint ACTION-ID | reconcile-worktree ACTION-ID | uninstall
-Install, stop, cancel, orphan mint revocation, worktree disposition and uninstall require typed y. Start uses the approved configuration.
+const runnerHelp = `usage: ailang-worldd service runner install --config FILE | start | stop | status | cancel ACTION-ID | reconcile-mint ACTION-ID | reconcile-worktree ACTION-ID | reconcile-files ACTION-ID | uninstall
+Install, stop, cancel, orphan mint revocation, worktree and file disposition and uninstall require typed y. Start uses the approved configuration.
 Config JSON keys: home, uid, address, world, credentials, executables (absolute path to sha256 digest), preset_names, tool_dirs, backend ("headless").
 The World service must already be running. Use absolute, canonical paths and include its installed binary in executables.
 `
@@ -272,7 +272,7 @@ func runServiceRunner(args []string, out, errOut io.Writer, e serviceEnv) int {
 		if len(args) != 3 || args[1] != "--config" {
 			return refuse(errors.New("expected install --config FILE"))
 		}
-	case "cancel", "reconcile-mint", "reconcile-worktree":
+	case "cancel", "reconcile-mint", "reconcile-worktree", "reconcile-files":
 		if len(args) != 2 || !verbs.ValidOperatorActionID(args[1]) {
 			return refuse(fmt.Errorf("expected %s ACTION-ID", verb))
 		}
@@ -289,6 +289,9 @@ func runServiceRunner(args []string, out, errOut io.Writer, e serviceEnv) int {
 	if verb == "reconcile-worktree" {
 		return e.reconcileRunnerWorktree(args[1], out, errOut, readRunnerWorktree)
 	}
+	if verb == "reconcile-files" {
+		return e.reconcileRunnerFiles(args[1], out, errOut)
+	}
 	p := e.runnerPaths()
 	if verb == "status" {
 		if err := e.runnerStatus(out); err != nil {
diff --git a/cmd/ailang-worldd/operator_start_input.go b/cmd/ailang-worldd/operator_start_input.go
index 5576a2e1..95e77cd3 100644
--- a/cmd/ailang-worldd/operator_start_input.go
+++ b/cmd/ailang-worldd/operator_start_input.go
@@ -36,10 +36,25 @@ func operatorStartFacts(ctx context.Context, signed approval.OperatorAction, att
 	if served != signed.Store {
 		return store.OperatorMintInput{}, nil, &approval.BindingError{Field: "store"}
 	}
+	tokenPath, err := operatorPhysicalTarget(req.outPath)
+	if err != nil {
+		return store.OperatorMintInput{}, nil, err
+	}
+	taskPath, err := operatorPhysicalTarget(read.taskCopy)
+	if err != nil {
+		return store.OperatorMintInput{}, nil, err
+	}
+	mcp := read.mcp
+	if mcp != "" {
+		mcp, err = operatorPhysicalTarget(mcp)
+		if err != nil {
+			return store.OperatorMintInput{}, nil, err
+		}
+	}
 	start, task := *req.start, *req.start.Task
-	task.Path = read.taskCopy
+	task.Path = taskPath
 	start.Task = &task
-	in := store.OperatorMintInput{Attempt: attempt, Start: start, TokenPath: req.outPath,
+	in := store.OperatorMintInput{Attempt: attempt, Start: start, TokenPath: tokenPath,
 		Profile: approval.ProfileRef{Project: req.profile, Digest: req.prof.Digest}}
 	for _, g := range req.grants {
 		in.Grants = append(in.Grants, approval.Grant{Effect: g.Effect, Scope: g.Scope, Budget: g.Budget})
@@ -53,9 +68,10 @@ func operatorStartFacts(ctx context.Context, signed approval.OperatorAction, att
 		return store.OperatorMintInput{}, nil, err
 	}
 	facts := []approval.OperatorResolved{{Name: "launch_binary", Value: read.launchBinary},
-		{Name: "launch_wall", Value: read.launchWall.String()}, {Name: "mcp_path", Value: read.mcp}}
+		{Name: "launch_wall", Value: read.launchWall.String()}, {Name: "mcp_path", Value: mcp}}
 	facts = append(facts, mint...)
-	facts = append(facts, approval.OperatorResolved{Name: "served_build", Value: string(build)},
+	facts = append(facts, approval.OperatorResolved{Name: "served_address", Value: read.preset.Daemon.Addr},
+		approval.OperatorResolved{Name: "served_build", Value: string(build)},
 		approval.OperatorResolved{Name: "served_interpreter", Value: read.health.ToolRef},
 		approval.OperatorResolved{Name: "worktree_path", Value: req.dir})
 	return in, facts, nil
diff --git a/cmd/ailang-worldd/operator_start_input_test.go b/cmd/ailang-worldd/operator_start_input_test.go
index 3f056ca0..3726ca54 100644
--- a/cmd/ailang-worldd/operator_start_input_test.go
+++ b/cmd/ailang-worldd/operator_start_input_test.go
@@ -84,7 +84,7 @@ func TestOperatorStartFacts(t *testing.T) {
 		value[fact.Name] = fact.Value
 	}
 	want := []string{"launch_binary", "launch_wall", "mcp_path", "mint_catalogue", "mint_grants", "mint_profile", "mint_project",
-		"mint_start", "mint_token_path", "served_build", "served_interpreter", "worktree_path"}
+		"mint_start", "mint_token_path", "served_address", "served_build", "served_interpreter", "worktree_path"}
 	if !reflect.DeepEqual(names, want) {
 		t.Fatalf("reviewed fact names %v, want %v", names, want)
 	}
@@ -94,9 +94,13 @@ func TestOperatorStartFacts(t *testing.T) {
 	}
 	build, _ := json.Marshal(testBuild)
 	h := matchingHealth(t, f.r.f, f.r.profiles, f.profile)
+	tokenPath, err := operatorPhysicalTarget(filepath.Join(f.r.sessions, "ep.session"))
+	if err != nil {
+		t.Fatal(err)
+	}
 	for name, v := range map[string]string{"launch_binary": binary, "launch_wall": "1h0m0s", "mcp_path": "",
 		"mint_catalogue": "3", "mint_profile": f.profile.Digest, "mint_project": "proj",
-		"mint_token_path": filepath.Join(f.r.sessions, "ep.session"), "served_build": string(build),
+		"mint_token_path": tokenPath, "served_address": defaultEpisodeAddr, "served_build": string(build),
 		"served_interpreter": h.ToolRef, "worktree_path": filepath.Join(f.r.f.ws, "ep")} {
 		if value[name] != v {
 			t.Fatalf("fact %s = %q, want %q", name, value[name], v)
@@ -105,8 +109,8 @@ func TestOperatorStartFacts(t *testing.T) {
 	in := r.in
 	if in.Start.StartedAt != time.Unix(presetClock(), 0).Format(time.RFC3339) || in.Start.Episode != "ep" || in.Start.Base.SHA != f.base || in.Start.Preset.SHA256 != serviceHash(f.data) ||
 		in.Start.Branch != "world/ep" || in.Start.BranchExisted || in.Start.Agent != "codex" || in.Start.TTLS != 43200 ||
-		in.Start.Task == nil || in.Start.Task.Path != filepath.Join(f.r.sessions, "ep.task.md") ||
-		in.Start.Task.SHA256 != serviceHash([]byte(operatorStartTask)) || in.TokenPath != filepath.Join(f.r.sessions, "ep.session") ||
+		in.Start.Task == nil || in.Start.Task.Path != filepath.Join(filepath.Dir(tokenPath), "ep.task.md") ||
+		in.Start.Task.SHA256 != serviceHash([]byte(operatorStartTask)) || in.TokenPath != tokenPath ||
 		in.Profile != (approval.ProfileRef{Project: "proj", Digest: f.profile.Digest}) || len(in.Grants) == 0 || in.Attempt != r.at {
 		t.Fatalf("effective input differs from the fixture: %+v", in)
 	}
@@ -161,7 +165,7 @@ func TestOperatorStartChanged(t *testing.T) {
 			episodeSelfBuild = func() daemon.BuildIdentity { return other }
 		}},
 		{name: "interpreter", field: "served_interpreter", change: func(t *testing.T, r *operatorStartRig) {
-			r.action.Resolution[10].Value = "sha256:" + serviceHash([]byte("other"))
+			r.action.Resolution[11].Value = "sha256:" + serviceHash([]byte("other"))
 		}},
 		{name: "launch_binary", field: "launch_binary", change: func(t *testing.T, r *operatorStartRig) {
 			other := startReadFakeAgent(t, r.f)
@@ -173,7 +177,7 @@ func TestOperatorStartChanged(t *testing.T) {
 			r.action.Resolution[2].Value = filepath.Join(r.f.r.sessions, "ep.mcp.json")
 		}},
 		{name: "worktree_path", field: "worktree_path", change: func(t *testing.T, r *operatorStartRig) {
-			r.action.Resolution[11].Value += "-other"
+			r.action.Resolution[12].Value += "-other"
 		}},
 		{name: "extra_fact", field: "resolution", change: func(t *testing.T, r *operatorStartRig) {
 			r.action.Resolution = append(r.action.Resolution, approval.OperatorResolved{Name: "zz_extra", Value: "x"})
diff --git a/cmd/ailang-worldd/operator_worktree_test.go b/cmd/ailang-worldd/operator_worktree_test.go
index 7513c2ce..8ab9820e 100644
--- a/cmd/ailang-worldd/operator_worktree_test.go
+++ b/cmd/ailang-worldd/operator_worktree_test.go
@@ -28,7 +28,7 @@ type worktreeStartRig struct {
 	src operatorStartSource
 }
 
-func newWorktreeStartRig(t *testing.T, claude bool) *worktreeStartRig {
+func newWorktreeStartRig(t *testing.T, claude bool, ttl ...int64) *worktreeStartRig {
 	t.Helper()
 	s := operatorService(t)
 	s.mustInstall(t)
@@ -51,6 +51,9 @@ func newWorktreeStartRig(t *testing.T, claude bool) *worktreeStartRig {
 		if claude {
 			p.Agent = "claude"
 		}
+		if len(ttl) != 0 {
+			p.TTLS = ttl[0]
+		}
 	})
 	h := matchingHealth(t, r.f.r.f, r.f.r.profiles, r.f.profile)
 	h.DBPath = s.config.World
diff --git a/cmd/ailang-worldd/service.go b/cmd/ailang-worldd/service.go
index 8aeb906b..e26eaf05 100644
--- a/cmd/ailang-worldd/service.go
+++ b/cmd/ailang-worldd/service.go
@@ -34,7 +34,7 @@ const serviceHelp = `usage: ailang-worldd service install|status|upgrade|stop|st
     [--examples-dir DIR] [--surface-index PATH] [--surface-follow PATH]
     [--surface-origin https://<host>.ts.net[:port]]... [--binary PATH] [--replace]
   status [--json] [--verbose]
-  runner install|start|stop|status|cancel|reconcile-mint|reconcile-worktree|uninstall (see service runner --help)
+  runner install|start|stop|status|cancel|reconcile-mint|reconcile-worktree|reconcile-files|uninstall (see service runner --help)
   upgrade [--binary FILE] [--pause]
 Install, stop and uninstall require a typed y at the controlling terminal.
 Start uses the approved plist; edited plists require confirmation.
diff --git a/docs/QUICKSTART.md b/docs/QUICKSTART.md
index c4defbed..4ba3f894 100644
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -107,6 +107,13 @@ Review whether Git evidence is absent or retained, then type `y` to stop the
 action and release its owner. Any worktree, branch or files stay in place.
 Run the disposable example with `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./cmd/ailang-worldd -run '^TestOperatorWorktreeWindow$' -count=1 -v`.
 
+For an interrupted file reservation, review the output of this command. Type `y` to stop when the entire signed set is absent or exact. Keep the action reserved when any file remains uncertain. The command preserves the credential, worktree and all files.
+
+```sh
+ailang-worldd service runner reconcile-files oa-0123456789abcdef0123456789abcdef
+```
+
+
 ## Run your agents through World
 
 World launchers now tell Claude Code, Codex and pi their episode, worktree, session
diff --git a/host/store/credentials_test.go b/host/store/credentials_test.go
index db3785c7..6d1744de 100644
--- a/host/store/credentials_test.go
+++ b/host/store/credentials_test.go
@@ -314,6 +314,8 @@ func TestCredentialsMethodsRequireADeadline(t *testing.T) {
 	_, errs["ReserveOperatorWorktree"] = c.ReserveOperatorWorktree(ctx, OperatorMintAttempt{}, "", nil, func() int64 { return 1000 })
 	_, _, errs["ReadOperatorWorktree"] = c.ReadOperatorWorktree(ctx, "x")
 	errs["DisposeOperatorWorktree"] = c.DisposeOperatorWorktree(ctx, OperatorStepKey{}, OperatorWorktreeAbsent, 1000)
+	errs["ReadOperatorFiles"] = func() error { _, _, err := c.ReadOperatorFiles(ctx, "action"); return err }()
+	errs["RecordOperatorFiles"] = c.RecordOperatorFiles(ctx, OperatorFileState{}, "", nil, nil)
 	errs["CommitOperatorMint"] = c.CommitOperatorMint(ctx, OperatorMintInput{}, nil, nil)
 	_, errs["ReadOperatorMint"] = c.ReadOperatorMint(ctx, "x")
 	errs["AbortOperatorMint"] = c.AbortOperatorMint(ctx, OperatorMintReceipt{}, 1000)
diff --git a/host/store/operator_files.go b/host/store/operator_files.go
new file mode 100644
index 00000000..786bc0f8
--- /dev/null
+++ b/host/store/operator_files.go
@@ -0,0 +1,232 @@
+package store
+
+import (
+	"context"
+	"database/sql"
+	"encoding/json"
+	"errors"
+	"strconv"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+type OperatorFileState struct {
+	Mint      OperatorMintReceipt
+	Session   SessionRow
+	Worktree  string
+	StartJSON string
+	Step      OperatorStepKey
+	Phase     string
+}
+
+func readOperatorFiles(ctx context.Context, tx *sql.Tx, id string) (OperatorFileState, approval.OperatorAction, error) {
+	var state OperatorFileState
+	var action approval.OperatorAction
+	mint, err := readOperatorMint(ctx, tx, id, true)
+	if err != nil {
+		return state, action, err
+	}
+	state.Mint = mint
+	row, err := operatorCASRow(ctx, tx, mint.CAS)
+	if err != nil {
+		return state, action, err
+	}
+	action, err = decodeOperatorMintAction(row)
+	if err != nil || !action.Options.Headless {
+		return state, action, ErrOperatorConflict
+	}
+	err = tx.QueryRowContext(ctx, `SELECT w.identity, s.start_json, c.credential_id, c.episode_id, c.grants_json,
+		c.expires_at, c.created_at, coalesce(c.project,''), coalesce(c.profile_digest,'')
+		FROM operator_action_steps w JOIN episode_starts s ON s.episode = ?
+		JOIN session_credentials c ON c.episode_id=s.episode AND c.credential_id=?
+		WHERE w.action_id = ? AND w.kind = 'worktree' AND w.owner = ? AND w.phase = 'observed' AND w.receipt = ?
+		AND (SELECT count(*) FROM operator_action_steps WHERE action_id = ? AND kind = 'worktree') = 1
+		AND (SELECT count(*) FROM operator_action_steps WHERE action_id = ? AND kind = 'mint') = 1`,
+		mint.Episode, mint.CredentialID, id, mint.Owner, []byte("created"), id, id).Scan(
+		&state.Worktree, &state.StartJSON, &state.Session.CredentialID, &state.Session.EpisodeID, &state.Session.GrantsJSON,
+		&state.Session.ExpiresAt, &state.Session.CreatedAt, &state.Session.Project, &state.Session.ProfileDigest)
+	if err != nil {
+		return state, action, err
+	}
+	state.Step = OperatorStepKey{Action: mint.CAS, Owner: mint.Owner, Kind: "files"}
+	err = tx.QueryRowContext(ctx, `SELECT identity, version, phase FROM operator_action_steps
+		WHERE action_id = ? AND kind = 'files' AND owner = ?`, id, mint.Owner).
+		Scan(&state.Step.Identity, &state.Step.Version, &state.Phase)
+	if err != nil && !errors.Is(err, sql.ErrNoRows) {
+		return state, action, err
+	}
+	var count int
+	err = tx.QueryRowContext(ctx, "SELECT count(*) FROM operator_action_steps WHERE action_id = ? AND kind = 'files'", id).Scan(&count)
+	if err != nil || (count != 0 && (count != 1 || state.Step.Identity == "")) {
+		return state, action, ErrOperatorConflict
+	}
+	return state, action, nil
+}
+
+func (c *Credentials) ReadOperatorFiles(ctx context.Context, id string) (OperatorFileState, approval.OperatorAction, error) {
+	if err := c.operatorAccess(ctx); err != nil {
+		return OperatorFileState{}, approval.OperatorAction{}, err
+	}
+	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
+	if err != nil {
+		return OperatorFileState{}, approval.OperatorAction{}, err
+	}
+	defer tx.Rollback()
+	return readOperatorFiles(ctx, tx, id)
+}
+
+const (
+	OperatorFilesUncertain       = ""
+	OperatorFilesPublished       = "published"
+	OperatorFilesStoppedAbsent   = "stopped: absent"
+	OperatorFilesStoppedRetained = "stopped: retained"
+)
+
+// RecordOperatorFiles reserves an empty phase, or observes or disposes the exact reservation.
+func (c *Credentials) RecordOperatorFiles(ctx context.Context, reviewed OperatorFileState, value string, proof *approval.OperatorMintProof, clock func() int64) error {
+	if err := c.operatorAccess(ctx); err != nil {
+		return err
+	}
+	if c.lock == nil || c.lock.dbPath != reviewed.Mint.Store {
+		return ErrOperatorConflict
+	}
+	tx, err := c.db.BeginTx(ctx, nil)
+	if err != nil {
+		return err
+	}
+	defer tx.Rollback()
+	current, action, err := readOperatorFiles(ctx, tx, reviewed.Mint.CAS.ID)
+	if err != nil {
+		return err
+	}
+	if current != reviewed {
+		return ErrOperatorConflict
+	}
+	key, now := current.Step, clock()
+	switch current.Phase {
+	case "":
+		if value == "" || proof == nil {
+			return ErrOperatorInput
+		}
+		row, err := operatorCASRow(ctx, tx, key.Action)
+		if err != nil {
+			return err
+		}
+		if err := admitOperatorFiles(ctx, tx, row, current, action, proof, now); err != nil {
+			return err
+		}
+		key.Identity, key.Version = value, 1
+		if err := reserveOperatorStep(ctx, tx, key, now); err != nil {
+			return err
+		}
+	case "reserved":
+		switch value {
+		case OperatorFilesPublished, OperatorFilesStoppedAbsent, OperatorFilesStoppedRetained:
+			if err := observeOperatorStep(ctx, tx, key, []byte(value)); err != nil {
+				return err
+			}
+		case OperatorFilesUncertain:
+			res, err := tx.ExecContext(ctx, `UPDATE operator_action_steps SET version=version+1 WHERE action_id=?
+				AND kind='files' AND identity=? AND owner=? AND version=? AND phase='reserved'`,
+				key.Action.ID, key.Identity, key.Owner, key.Version)
+			if err := operatorChangedOne(res, err); err != nil {
+				return err
+			}
+		default:
+			return ErrOperatorInput
+		}
+	default:
+		return ErrOperatorConflict
+	}
+	// The step update can succeed before a trigger or other drift invalidates
+	// the action. Enforce the independent owner CAS before committing either.
+	phase, version, receipt := "reserved", key.Version, []byte(nil)
+	if current.Phase == "reserved" {
+		version++
+		if value != OperatorFilesUncertain {
+			phase, receipt = "observed", []byte(value)
+		}
+	}
+	if current.Phase == "" || value == OperatorFilesPublished {
+		after, action, err := readOperatorFiles(ctx, tx, key.Action.ID)
+		expected := current
+		expected.Step = key
+		expected.Phase, expected.Step.Version = phase, version
+		if err != nil || after != expected {
+			return ErrOperatorConflict
+		}
+		row, err := operatorCASRow(ctx, tx, key.Action)
+		if err != nil {
+			return err
+		}
+		if err := admitOperatorFiles(ctx, tx, row, after, action, proof, clock()); err != nil {
+			return err
+		}
+	}
+	stopped := current.Phase == "reserved" && (value == OperatorFilesStoppedAbsent || value == OperatorFilesStoppedRetained)
+	uncertain := current.Phase == "reserved" && value == OperatorFilesUncertain
+	res, err := tx.ExecContext(ctx, `UPDATE operator_actions
+		SET state=CASE WHEN ? THEN 'failed' WHEN ? THEN 'uncertain' ELSE state END,
+		    owner=CASE WHEN ? THEN NULL ELSE owner END,
+		    terminal_at=CASE WHEN ? THEN ? ELSE terminal_at END,
+		    generation=generation+CASE WHEN ? THEN 1 ELSE 0 END
+		WHERE id=? AND digest=? AND generation=? AND owner=? AND kind='start' AND state IN ('running','uncertain')
+		AND EXISTS (SELECT 1 FROM operator_action_steps WHERE action_id=? AND kind='files'
+			AND identity=? AND owner=? AND phase=? AND version=? AND receipt IS ?)
+		AND NOT EXISTS (SELECT 1 FROM operator_action_steps WHERE action_id=? AND kind NOT IN ('operation','worktree','mint','files'))`,
+		stopped, uncertain, stopped, stopped, now, stopped || uncertain,
+		key.Action.ID, key.Action.Digest[:], key.Action.Generation, key.Owner,
+		key.Action.ID, key.Identity, key.Owner, phase, version, receipt, key.Action.ID)
+	if err := operatorChangedOne(res, err); err != nil {
+		return err
+	}
+	return tx.Commit()
+}
+
+// admitOperatorFiles validates the persisted mint against the same law that
+// made it, including the present approver, session and original expiry.
+func admitOperatorFiles(ctx context.Context, tx *sql.Tx, row OperatorActionRow, current OperatorFileState,
+	action approval.OperatorAction, proof *approval.OperatorMintProof, now int64) error {
+	if now >= row.ExpiresAt || now >= current.Session.ExpiresAt {
+		return ErrOperatorExpired
+	}
+	if row.State != "running" || now < row.CreatedAt || now < current.Session.CreatedAt || len(row.Challenge) != 32 {
+		return ErrOperatorConflict
+	}
+	if proof == nil {
+		return ErrOperatorInput
+	}
+	b, err := proof.Bind(action, [32]byte(row.Challenge), row.ExpiresAt, row.Artifact)
+	if err != nil {
+		return err
+	}
+	if row.ApprovedBy != b.Format+":"+b.Name {
+		return ErrOperatorConflict
+	}
+	if err := operatorMintApprover(ctx, tx, proof, b, now); err != nil {
+		return err
+	}
+	catalogue, err := strconv.ParseInt(mintFact(action, "mint_catalogue"), 10, 64)
+	if err != nil {
+		return ErrOperatorConflict
+	}
+	in := OperatorMintInput{
+		Attempt:      OperatorMintAttempt{Catalogue: catalogue},
+		Profile:      approval.ProfileRef{Project: mintFact(action, "mint_project"), Digest: mintFact(action, "mint_profile")},
+		CredentialID: current.Mint.CredentialID, TokenPath: current.Mint.TokenPath,
+	}
+	if json.Unmarshal([]byte(current.StartJSON), &in.Start) != nil || in.Start.Allowance != nil ||
+		json.Unmarshal([]byte(mintFact(action, "mint_grants")), &in.Grants) != nil ||
+		in.Start.StartedAt != time.Unix(current.Session.CreatedAt, 0).UTC().Format(time.RFC3339) {
+		return ErrOperatorConflict
+	}
+	want, err := operatorMintSession(in, action, current.Session.CreatedAt)
+	if err != nil {
+		return err
+	}
+	if want != current.Session {
+		return ErrOperatorConflict
+	}
+	return nil
+}
diff --git a/host/store/operator_mint.go b/host/store/operator_mint.go
index 1f7eee47..0267140f 100644
--- a/host/store/operator_mint.go
+++ b/host/store/operator_mint.go
@@ -243,10 +243,10 @@ func (c *Credentials) ReadOperatorMint(ctx context.Context, id string) (Operator
 	if err := c.operatorAccess(ctx); err != nil {
 		return OperatorMintReceipt{}, err
 	}
-	return readOperatorMint(ctx, c.db, id)
+	return readOperatorMint(ctx, c.db, id, false)
 }
 
-func readOperatorMint(ctx context.Context, db operatorMintReader, id string) (OperatorMintReceipt, error) {
+func readOperatorMint(ctx context.Context, db operatorMintReader, id string, fileStage bool) (OperatorMintReceipt, error) {
 	row, err := scanOperatorAction(db.QueryRowContext(ctx, `SELECT `+operatorColumns+` FROM operator_actions WHERE id = ?`, id))
 	if err != nil {
 		return OperatorMintReceipt{}, err
@@ -263,8 +263,8 @@ func readOperatorMint(ctx context.Context, db operatorMintReader, id string) (Op
 	var receipt []byte
 	err = db.QueryRowContext(ctx, `SELECT version, receipt FROM operator_action_steps WHERE kind = 'mint'
 		AND identity = ? AND action_id = ? AND phase = 'observed' AND owner = ?
-		AND NOT EXISTS (SELECT 1 FROM operator_action_steps WHERE action_id = ? AND kind NOT IN ('operation','worktree','mint'))`,
-		id, id, row.Owner, id).Scan(&r.StepVersion, &receipt)
+		AND NOT EXISTS (SELECT 1 FROM operator_action_steps WHERE action_id = ? AND kind NOT IN ('operation','worktree','mint') AND NOT (? AND kind = 'files'))`,
+		id, id, row.Owner, id, fileStage).Scan(&r.StepVersion, &receipt)
 	if err != nil {
 		return OperatorMintReceipt{}, fmt.Errorf("mint evidence unavailable or later effect reserved: %w", err)
 	}
@@ -299,7 +299,7 @@ func (c *Credentials) AbortOperatorMint(ctx context.Context, reviewed OperatorMi
 		return err
 	}
 	defer tx.Rollback()
-	current, err := readOperatorMint(ctx, tx, reviewed.CAS.ID)
+	current, err := readOperatorMint(ctx, tx, reviewed.CAS.ID, false)
 	if err != nil {
 		return err
 	}
````
