# Sprint plan — w-world-address-bar-verbs-v4, sprint S5e (iteration 435, row 243)

Authority is D-WORLD-119 = A and D-WORLD-120 = A (attended 82f05004). The r1 design has quorum r2 PROCEED 3/3 with no absent seats. S1 #606, S2 #609, S2b #611, S3 #643, S4 #651, S4b #654, S5 #657, S5b #659, S5c #661 and S5d #666 landed. S5d merged at bb214516 with judge Sonnet PASS 84 and zero blocking findings. Base is ec4565ec5eedfed693f9ea5f2c78b89b63bf432e.

## 0. Binding constraints

- Land and verify M1 before a consumer builds on the reservation path. File reservation, installation, independent readback and attended disposition travel in a single milestone in Appendix C.
- Design §2.2 and §5.5 govern activation. Keep the production reservation/start/mint chain dormant and retain the existing attended worktree command for an existing reserved step. Row 271 owns the 14 real-binary se-tools failures.
- The executor applies with `/usr/bin/patch -p1 -E` and snapshots. The controller commits. Use migrated fixture stores and the repo signer. Real Git checkout mutations and file writes in tests occur inside `t.TempDir()` fixtures. Gate logs and AILANG state go to scratch. Leave TMPDIR and GOTMPDIR alone.

## 1. Pre-flight

### 1.1 Baseline at ec4565ec

Before staging M1, the §6 baseline ran on `git archive ec4565ec` under `$S/base`. `$S/gate-base/legs.txt` and `ALL_DONE` recorded 160 s and zero informative failures. Build, vet, compile-only fence, AILANG verification, plain/race store, named inventories, routes, callers and scope passed. The gofmt leg accounted for the single inherited `host/daemon/daemon.go` drift by byte comparison with the base blob.

CLI whole package, CLI race and QUICKSTART §9 returned **UNINFORMATIVE UNDER SANDBOX** on a loopback denial. The extra archive probe `go test -json -count=1 -timeout 8m ./cmd/ailang-worldd ./host/daemon` recorded the same denial and its sandbox-only fail set in `$S/base-fails.archive` (CLI `TestApproveDecideWriteMessageAndSignature`, daemon `TestAGUISeesDirectStoreCommit` and `TestAGUIShutdownEndsRuns`). CLI execution stopped at the first bind panic. Archive `verify_go` was **SKIPPED** due to the Git metadata requirement. Real-path observation legs were **UNINFORMATIVE UNDER SANDBOX**.

The 14 names from iter434's controller are in `$S/base-fails.expected`. Run `/bin/bash "$S/baseline-controller.sh" "$BASE_TREE"` on a pristine ec4565ec product tree in a Git checkout outside the sandbox. A controller plan-only commit is accepted after an all-path comparison. It records the actual plain/race CLI and full plain `verify_go` fail sets separately. It compares them to the row-271 record and deposits `pkg-cmd-ailang-worldd.base-fails`, `race-cmd-ailang-worldd.base-fails` and `verify_go.base-fails` for the final comparisons.

### 1.2 Current inputs and coupling inventory

Full command output is in `$S/preflight.txt`. For function and inventory rows, use `git grep -n -E 'func .*SYMBOL_REGEX\(' -- cmd host` with the listed symbols. Literal searches use `git grep -n -E 'LITERAL_REGEX' -- PATHS`.

| ID | Command or regex | Observed at ec4565ec |
|---|---|---|
| C1 | `git rev-parse HEAD`, `git status --short`, `git diff --name-only bb214516 ec4565ec -- cmd host docs app \| wc -l`, `git ls-tree --name-only ec4565ec -- cmd host docs app` | ec4565ec5eedfed693f9ea5f2c78b89b63bf432e, empty status, zero changed paths, `cmd`, `docs`, `host`. |
| C2 | `reserveOperatorStartWorktree\(\|MintOperatorStart\(\|reconstructOperatorStart\(\|authority.RecheckVerb\(`, `-- '*.go'`, exclude tests | Definitions at `operator_worktree.go:31`, `host/authority/mint.go:204`, `operator_start_input.go:68`. Internal call at `operator_worktree.go:45`. Positive live control `episode_finish_surface.go:239`. |
| C3 | `ReserveOperatorWorktree\|ReadOperatorWorktree\|DisposeOperatorWorktree\|readRunnerWorktree\|reconcileRunnerWorktree\|operatorPhysicalTarget\|operatorGitCommon` | Store `host/store/operator_worktree.go:13,73,108`. CLI `operator_worktree_reconcile.go:19,48` and `operator_worktree.go:131,159`. Numbered read `nl -ba cmd/ailang-worldd/operator_worktree_reconcile.go` locates identities at 38–41 and the locked re-read at 66–71. |
| C4 | `ReserveOperatorStep\|ObserveOperatorStep\|OperatorStep\|operatorCASRow\|readOperatorMint` and `kind NOT IN\|NOT EXISTS` in mint store | `operator_steps.go:17,45,72`, `operator_cas.go:32` and `operator_mint.go:249,266`. Mint receipt reader excludes effects beyond operation/worktree/mint. Schema regex `length\(challenge\)\|PRIMARY KEY \(kind\|worktree.*mint.*files` returns `operator_schema.sql:9,49,58`. |
| C5 | `provisionEpisode\|writeTokenFile\|writeMCPConfig` plus `taskCopy :=\|req.afterToken =\|os.OpenFile\(taskCopy\|writeMCPConfig\(mcp,\|writeTokenFile\(r.outPath` | `session_new.go:405,549,674`, `episode_preset.go:195,340,342,358`, `episode.go:554`. `sed -n '185,212p' cmd/ailang-worldd/episode_preset.go` shows the MCP path for Claude and empty MCP for headless Codex. Signed facts in `operator_start_input.go:54–60` include MCP, worktree and mint resolution. |
| C6 | `reconcile-(mint\|worktree)\|installedRunner\|runnerWrite\|runnerDB\|openRunnerLock\|operatorEpisodeLockPath` | Mint resolver `operator_mint.go:17,49`, worktree remedy `operator_worktree_reconcile.go:15`, grammar/dispatch `operator_service.go:24,222,275,286,289`. Usage `main.go:86` and `service.go:37`, runbook `docs/QUICKSTART.md:90,103`. Openers `operator_runner.go:25,162,180` and `operator_episode_lock.go:14`. |
| C7 | `git grep -n -E 'syscall\.Flock.*LOCK_NB' -- cmd host \| grep -vF '_test.go:'` | Seven lines including one comment. CLI `operator_runner.go:35`, `operator_service.go:251`, `service.go:158,662`. Writer-lock `writer_lock_unix.go:6,41,86`. |
| C8 | `TestDaemonDependencyAllowlist\|TestWorldBoundaryDependencyAllowlist\|TestProductionGoSurface\|TestProductionContextRoots\|TestCLICredentialWriterSites\|TestEveryVerbHasHelp\|TestVerbTableMatchesTheCLISwitch\|TestQuickstartServiceSectionMatchesTheCLIHelp\|TestCredentialsMethodsRequireADeadline\|TestStoreSurfaceTakesContext\|TestBoundaryASTWriteGuard\|TestOperatorNeverDispatch\|TestPluginCommandsMatchTheCLIHelp` | Daemon `daemon_test.go:899`, boundary `allowlist_world_test.go:758,1323`, store `context_roots_test.go:129,234`, `deadline_guard_test.go:43`, `credentials_test.go:277`, CLI `cli_credential_writers_test.go:18`, `help_test.go:21`, `verbs_parity_test.go:66`, `quickstart_service_test.go:15`, `plugin_help_test.go:52`, publish plugin `cmd/world-publish/plugin_help_test.go:58`, verbs `operator_policy_test.go:69`. |
| C9 | `git grep -niE 'reconcile[- ]?(file\|worktree\|mint)\|TestOperatorFilesReconcile' -- design_docs/planned/w-world-address-bar-verbs-v4.md` | Sole match is line 234 naming `TestOperatorFilesReconcile` and a repair path. §3.3 lines 108–109 specify durable evidence and no overwrite. The reviewed design supplies no literal CLI spelling. The prototype uses `service runner reconcile-files ACTION-ID` following C6. |


The companion `$S/anchors.txt` records fixture and review anchors from `git grep -n -E 'func (canonicalLaunchPath|presetGit|newEpisodeRun|observeProvision|operatorClaimFixture|assertReservedOperatorStep|claimFixtureEpisodeAction|CanonicalOperatorAction|OperatorSignatureDigest|OperatorActionMessage|CheckOperatorAction|credentialsFor)\(' -- cmd host`. Representative results are `episode_launch.go:73`, `preset_test.go:28`, `episode_test.go:85`, `operator_provision_test.go:109`, `operator_cas_test.go:12,279`, `operator_episode_lock_test.go:30`, `operator_action.go:189,302,315`, `operator_signed.go:83,203` and `session.go:121`. These anchors identify the Git fixture, migrated store and signer/receipt seams used by the patch.

| Inventory | Impact and verification |
|---|---|
| Help, command and QUICKSTART | M1 threads a reader through the existing handler. §6 runs help and verb inventories, service runbook, both plugin parity tests and QUICKSTART §9. The file milestone needs nested runner help, top-level/service usage, help test enumeration and QUICKSTART updated together with its executable example. |
| Credential writer sites and flock flags | The same writer opener and kernel sidecar are used. §6 checks CLI/daemon writer inventories and the seven-line census. |
| Package, import and context roots | M1 adds two test files in existing CLI/store packages. §6 runs the daemon/boundary allowlists, protected command groups, AST write guard, store production surface, context roots and method deadlines. The prototype's new credentials methods have deadline-inventory entries and the prototype passed both boundary tests and named store/CLI inventories in `$S/probe-M2-inventories.log`. Repeat them for the final file patch. |
| Routes and served behaviour | §6 compares mux, surface source and app JS bytes and runs route/dispatch inventories. The live finish verifier is the positive caller control. |

### 1.3 Decisions

1. Run store guards in the migrated credentials store with the repository signer. Check the exact fast refusal and absence of a state change.
2. Move the store start-kind predicate to the row guard before decoding. In the landed post-decode form, every valid finish also fails headless. The earlier predicate gives damaged non-start rows a typed conflict before decode or proof binding. The existing attended worktree handler takes its current reader as a narrow parameter for deterministic between-read drift.
3. Test real Git matching branches and SHAs across separate common directories, real duplicate registration and deterministic filesystem timing around real Git reads. Persist stable, internally agreeing tampered rows in the CLI tests so their CAS cannot supply the early refusal.
4. M1 lands independently. The 1,462-line joint measured candidate exceeds the cap. Keep all file and resolver gates together in Appendix C. The signed Codex/Claude distinction supplies no approval to split Claude's three-file bundle.

### 1.4 Premises found false or incomplete at HEAD

| ID | Command and observed result | Resolution |
|---|---|---|
| F1 | C9 returned only design line 234 and no command name. C6 returned the two precedents. | Use proposed `reconcile-files` in the file prototype. |
| F2 | `grep -F 'bind: operation not permitted' "$S/base-e2e.jsonl"` shows the CLI httptest panic and daemon listen failures. `$S/base-fails.archive` records three sandbox-specific test names. | Use the controller Git checkout for real-binary fail-set tracking. |
| F3 | `test -e "$HOME/.ailang/state/world-iter433/plan-scratch/st-s5d-footprint/cmd/ailang-worldd/s5c_files_probe.go"` returned 1. S5d C12 records 53 lines. | Carry that as prior evidence. The current complete-shape prototype determines the cap decision. |
| F4 | `git grep -nF "kind NOT IN ('operation','worktree','mint')" -- host/store/operator_mint.go` returns 266. `git grep -nF '[]byte("created")' -- cmd/ailang-worldd/operator_worktree.go` returns 106. | Require the observed-created receipt and matching mint, preserving the old mint resolver's later-effect refusal. |
| F5 | `grep -n 'finish.*Headless' host/approval/operator_action.go` returns 151. | The old kind-only mutation was masked for valid finish. M1's relocated guard refuses a damaged canonical finish record before decode. |
| F6 | `$S/numstat.txt` records 967 incremental and 1,462 cumulative. `$S/verify-M2-probe.log` shows the CLI and store focused tests passing. | The full file slice crosses the cap before the remaining mutation and Linux gates. |

### 1.5 Size and CI race cost

Controller cap command is exactly `git diff -M --numstat ec4565ec HEAD -- cmd host docs app`. Sum additions and deletions. Repeat without the pathspec for the literal unfiltered sum. `python3 "$S/stage-diff.py" M1 M2` measures read-only no-index diffs of archive trees and separately scans all product paths. No renames or binary changes. The separate plan commit adds both plan files to the literal unfiltered history. The JSON records the corresponding full expected line count.

| Milestone | Added | Deleted | Cumulative | Increment | Unfiltered product | STATUS/DELETED/MANIFEST |
|---|---:|---:|---:|---:|---:|---:|
| M1 | 489 | 6 | **495** | 495 | 495 | 5/0/5 |
| M2 probe, Appendix C | 1,442 | 20 | **1,462** | 967 | 1,462 | Deferred |

On the archive rebuilt from Appendix P, `$S/gate-final-archive/ALL_DONE` records 190 s, no informative failures and the same sandbox limitations as base. Named acceptance passed. The new CLI families under `-race` took **7.668 s** of package time (9 s gate wall). The same selector at base had no test names and took 1.648 s, a **6.020 s** measured increment. See `added-cli-race.txt` and `$S/base-cli-race-empty.log`. The whole CLI race delta requires the external run. The route census was 54 and three surface POSTs at base and M1. The caller leg found the live finish verifier and zero production file entries.

### 1.6 Portability

Darwin arm64, Go 1.26.6. Normalize physical paths on both sides with `filepath.EvalSymlinks`. Use typed refusal or unchanged snapshots for platform-dependent errors. The file prototype uses directory-relative Unix open/stat and exclusive create with no-follow semantics. Test restrictive umask, exact chmod, linked/unattributed files and independent descriptor evidence on both Darwin and Ubuntu 24.04 before that patch lands. M1 adds no tests to `host/daemon` or `host/verifygate`.

## 2. How the executor applies a sub-milestone

````bash
export GIT_OPTIONAL_LOCKS=0
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter435-row243-v4s5e
T=$HOME/.ailang/state/world-iter435/executor
S=$HOME/.ailang/state/world-iter435/plan-scratch
B=ec4565ec5eedfed693f9ea5f2c78b89b63bf432e
PLAN=$W/design_docs/planned/w-world-address-bar-verbs-v4s5e-sprint-plan.md
mkdir -p "$T"
unset SNAP_PARENT SNAP_STATUS_LIST
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

Apply M1 once to the base product tree. Snapshot scope excludes both plan records.

## 3. Milestones

### M1 — Reservation-path pre-caller pins

| Boundary / input | Named test and observed state |
|---|---|
| Store start kind, headless, canonical store, stored approver, length and claimed state | `TestReserveOperatorWorktreeGuards/{kind,headless,store,approver,challenge,claimed}`. Each refusal yields `ErrOperatorConflict`, exact action snapshot, existing operation row only and no credential/start. Valid control reads reserved version 1 and the incremented running generation. |
| Requested kind, replaced step owner and stale version | `TestDisposeOperatorUncertainStepGuards/{kind,owner,version}`. Wrong kind is input error, stale caller state is a conflict. Independent action and reserved step stay intact. |
| Receipt size | `TestOperatorStepReceiptAtLimit` observes exact byte-distinct maximum receipt and version 2. Inherited `TestOperatorStepGuards` covers empty and above maximum. |
| Wrong caller store, episode grammar, relative path | `TestOperatorWorktreeSurvivors` pins exact early refusal and unchanged store hash, Git registration/ref and both fixture file trees. |
| Initial and checkout common directory, link type/prefix, registration count, wrong directory, redirected ancestor | `TestOperatorWorktreeReadbackGuards`. Real SHA and symbolic HEAD are the same in the cross-repo decoy. The linked-prefix row puts deterministic repair behind the later Git read. |
| Stable tamper and already-observed row | `TestOperatorWorktreeRecoveryGuards` records conflict/missing, no prompt, and unchanged store/Git/files. The CAS sees the same tampered bytes as the reader. |
| Generation changes between initial and locked re-read | `TestOperatorWorktreeLockedReread` refuses before prompt and preserves the replacement rows and Git/file snapshot. |
| Last HEAD read fails | `TestOperatorReadbackErrorIsNeverCreated` reads a real exact positive control and a partial result on command failure. |

### 3.1 Design §5.3 row assignment

| Reviewed row | Assignment |
|---|---|
| `TestOperatorStartOnce` | M1 strengthens the existing worktree half. Files in Appendix C, receipt-before-spawn and confinement in S6, browser link in S8. |
| `TestOperatorStartChanged` | S5c comparison and S5d integration exercised by §6. |
| `TestOperatorStartCollisions` | M1 pins additional worktree/Git and early identity cases. Publication-time no-replace stays in the atomic file milestone. |
| `TestOperatorWorktreeWindow` | S5d fault windows and exact remedy in acceptance, plus new M1 readback/disposition cases. |
| `TestOperatorMintBinding`, `TestOperatorMintRollback`, `TestOperatorMintReceiptWindow` | S5b. §6 runs mint recovery regressions. Live mint consumption awaits S5f/S6. |
| `TestOperatorFilesReconcile` | Appendix C, reserve/install/readback and attended disposition together. |
| `TestOperatorLaunchSession`, `TestOperatorSpawnRefused`, `TestOperatorLaunchCrashWindow`, `TestOperatorLaunchExit` | S5f/S6 pairing first, then S6 launch. |
| `TestOperatorStartBrowser` | S8 fake Sol launcher and browser reconnect. |
| `TestOperatorHerdrUnavailable`, `TestOperatorHerdrHandoff`, `TestOperatorHerdrWindow` | S12 with E4 and typed ticket/TTY gates. |

### 3.2 Design §5.5 and §2.2 activation

| Reviewed gate | Assignment |
|---|---|
| `TestOperatorServiceFence` | §6 exercises typed approval and edited plist. M1 pins the existing attended resolver's pre-prompt refusals. |
| `TestOperatorServiceOwnership` | §6 privilege, owner, file modes/links and executable identity. |
| `TestOperatorServiceConfig` | §6 typed config and rejection of raw commands. |
| `TestOperatorServicePair` | §6 compatibility. Paired restart and readiness gate real activation in S5f/S6. |
| `TestOperatorServiceStartup` | §6 daemon unavailable, boot failure and bounded startup. Future pairing must preserve no queue across incompatible peers. |
| `TestOperatorServiceUpgrade` and `TestOperatorServiceUpgradePairFence` | §6 rollback and pair fence. Active-child drain, durable recovery and old-build exclusion remain S5f/S6 activation pins. |
| `TestOperatorServiceWriteScope` | §6 fixture lifecycle and real-path comparison in controller mode. |
| `TestOperatorServiceEnvironment` | §6 approved configuration and environment/credential boundaries. |
| `TestOperatorFinishReadiness` | §6 store readiness. Keep the E2-equivalent preflight before real finish claims. |
| `TestOperatorSchema10` | §6 and acceptance exercise fresh/v9/legacy-reader/migration-failure parity. |

Activation requires approved configuration, UID and matching executable/store/build, recent compatible heartbeat and catalogue, per-episode exclusion, exact prior worktree/mint/current credential, paired-process receipt and recovery, controlled drain and compatible rollback. §6 counts production callers and identifies the live `authority.RecheckVerb` control. The only carried worktree disposition entry is the attended CLI route.

## 4. Mutation table

Anchors are at cumulative M1. Exact old/new text and package paths are in `$S/mutants.json`. Named tests below drive runtime store, Git or CLI operations. A kill requires the exact named FAIL and absence of bind-denial text. Inventory-only failures cannot establish a kill.

31 named behavioural kills. 4 measured equivalent controls. 10 whole-package store kills.

| ID | Production anchor | Mutation | Test package and named test | Base named | M1 named | Whole |
|---|---|---|---|---|---|---|
| J-ST-KIND | host/store/operator_worktree.go:31 | `row.Kind != "start" \|\|` → `false \|\|` | host/store / TestReserveOperatorWorktreeGuards/kind | SURVIVED | KILLED | KILLED |
| J-ST-HEADLESS | host/store/operator_worktree.go:38 | `!action.Options.Headless {` → `false {` | host/store / TestReserveOperatorWorktreeGuards/headless | SURVIVED | KILLED | KILLED |
| J-ST-STORE | host/store/operator_worktree.go:38 | `action.Store != c.lock.dbPath \|\|` → `false \|\|` | host/store / TestReserveOperatorWorktreeGuards/store | SURVIVED | KILLED | KILLED |
| J-ST-APPROVEDBY | host/store/operator_worktree.go:45 | `row.ApprovedBy != b.Format+":"+b.Name` → `false` | host/store / TestReserveOperatorWorktreeGuards/approver | SURVIVED | KILLED | KILLED |
| J-ST-CHALLENGE | host/store/operator_worktree.go:31 | `len(row.Challenge) != 32 \|\|` → `false \|\|` | host/store / TestReserveOperatorWorktreeGuards/challenge | SURVIVED | KILLED | KILLED |
| J-ST-CLAIMEDSTATE | host/store/operator_worktree.go:31 | `row.State != "claimed" \|\|` → `false \|\|` | host/store / TestReserveOperatorWorktreeGuards/claimed | SURVIVED | KILLED | KILLED |
| J-ST-DISPKIND | host/store/operator_worktree.go:112 | `key.Kind != "worktree" {` → `false {` | host/store / TestDisposeOperatorUncertainStepGuards/kind | SURVIVED | KILLED | KILLED |
| J-ST-UNC-OWNER | host/store/operator_worktree.go:123 | `identity = ? AND owner = ? AND version = ? AND phase = 'reserved'`,` → `identity = ? AND ? IS NOT NULL AND version = ? AND phase = 'reserved'`,` | host/store / TestDisposeOperatorUncertainStepGuards/owner | SURVIVED | KILLED | KILLED |
| J-ST-UNC-VERSIONCAS | host/store/operator_worktree.go:123 | `identity = ? AND owner = ? AND version = ? AND phase = 'reserved'`,` → `identity = ? AND owner = ? AND ? IS NOT NULL AND phase = 'reserved'`,` | host/store / TestDisposeOperatorUncertainStepGuards/version | SURVIVED | KILLED | KILLED |
| J-STEP-RECEIPTMAX | host/store/operator_steps.go:49 | `len(receipt) > MaxSurfaceApprovalArtifactBytes` → `len(receipt) >= MaxSurfaceApprovalArtifactBytes` | host/store / TestOperatorStepReceiptAtLimit | SURVIVED | KILLED | KILLED |
| J-W-STORECHK | cmd/ailang-worldd/operator_worktree.go:37 | `signed.Store != canonical \|\|` → `false \|\|` | cmd/ailang-worldd / TestOperatorWorktreeSurvivors/store | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| J-W-EPPATTERN | cmd/ailang-worldd/operator_worktree.go:37 | `!verbs.EpisodeIDPattern.MatchString(signed.Episode) {` → `verbs.EpisodeIDPattern == nil {` | cmd/ailang-worldd / TestOperatorWorktreeSurvivors/episode | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| J-W-PHYS-ABS | cmd/ailang-worldd/operator_worktree.go:132 | `!filepath.IsAbs(path) \|\| filepath.Clean(path) != path` → `filepath.Clean(path) != path` | cmd/ailang-worldd / TestOperatorWorktreeSurvivors/relative | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| J-R-STORECHK | cmd/ailang-worldd/operator_worktree_reconcile.go:39 | `signed.Store != credentials \|\|` → `false \|\|` | cmd/ailang-worldd / TestOperatorWorktreeRecoveryGuards/store | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| J-R-BRANCHCHK | cmd/ailang-worldd/operator_worktree_reconcile.go:41 | `target.Branch != "world/"+signed.Episode \|\|` → `false \|\|` | cmd/ailang-worldd / TestOperatorWorktreeRecoveryGuards/branch | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| J-R-PATHCHK | cmd/ailang-worldd/operator_worktree_reconcile.go:41 | `target.Path != path \|\|` → `false \|\|` | cmd/ailang-worldd / TestOperatorWorktreeRecoveryGuards/path | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| J-R-ACTIONCHK | cmd/ailang-worldd/operator_worktree_reconcile.go:41 | `target.Action != id \|\|` → `false \|\|` | cmd/ailang-worldd / TestOperatorWorktreeRecoveryGuards/action | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| J-R-BASECHK | cmd/ailang-worldd/operator_worktree_reconcile.go:40 | `filepath.Base(path) != signed.Episode \|\|` → `false \|\|` | cmd/ailang-worldd / TestOperatorWorktreeRecoveryGuards/basename | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| J-R-HEADLESS | cmd/ailang-worldd/operator_worktree_reconcile.go:39 | `!signed.Options.Headless \|\|` → `false \|\|` | cmd/ailang-worldd / TestOperatorWorktreeRecoveryGuards/headless | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| J-ST-READPHASE | host/store/operator_worktree.go:89 | `owner = ? AND phase = 'reserved'` → `owner = ? AND phase IS NOT NULL` | cmd/ailang-worldd / TestOperatorWorktreeRecoveryGuards/observed | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| J-RB-COMMONID | cmd/ailang-worldd/operator_worktree.go:185 | `if err != nil \|\| common != target.Common {` → `if err != nil \|\| common == "" {` | cmd/ailang-worldd / TestOperatorWorktreeReadbackGuards/repo_common | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| J-RB-WTCOMMON | cmd/ailang-worldd/operator_worktree.go:252 | `e.common, err = operatorGitCommon(ctx, git, target.Path)` → `e.common, err = operatorGitCommon(ctx, git, target.Path)` | cmd/ailang-worldd / TestOperatorWorktreeReadbackGuards/checkout_common | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| J-RB-LINKREG | cmd/ailang-worldd/operator_worktree.go:249 | `!link.Mode().IsRegular() \|\| !worktree.IsLinked(target.Path)` → `!link.Mode().IsRegular() && false \|\| !worktree.IsLinked(target.Path)` | cmd/ailang-worldd / TestOperatorWorktreeReadbackGuards/link_regular | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| J-RB-ISLINKED | cmd/ailang-worldd/operator_worktree.go:249 | `!link.Mode().IsRegular() \|\| !worktree.IsLinked(target.Path)` → `!link.Mode().IsRegular()` | cmd/ailang-worldd / TestOperatorWorktreeReadbackGuards/link_prefix | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| J-RB-REGCOUNT | cmd/ailang-worldd/operator_worktree.go:235 | ` \|\| strings.Count(e.registration, "\x00\x00") != 1` → `(delete)` | cmd/ailang-worldd / TestOperatorWorktreeReadbackGuards/registration_count | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| J-RB-ISDIR | cmd/ailang-worldd/operator_worktree.go:235 | `!info.IsDir() \|\| exact != 1` → `exact != 1` | cmd/ailang-worldd / TestOperatorWorktreeReadbackGuards/directory | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| J-RB-SYMLINKEVAL | cmd/ailang-worldd/operator_worktree.go:239 | `physical != target.Path {` → `false && physical != target.Path {` | cmd/ailang-worldd / TestOperatorWorktreeReadbackGuards/physical | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| J-R-CURRENTKEY | cmd/ailang-worldd/operator_worktree_reconcile.go:71 | `current != key \|\| currentEpisode != episode` → `current.Action.ID == "" \|\| currentEpisode != episode` | cmd/ailang-worldd / TestOperatorWorktreeLockedReread | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| J-R-CURRENTEP | cmd/ailang-worldd/operator_worktree_reconcile.go:71 | `current != key \|\| currentEpisode != episode` → `current != key \|\| currentEpisode == ""` | cmd/ailang-worldd / TestOperatorWorktreeLockedReread | SURVIVED | SURVIVED (equivalent) | UNINFORMATIVE UNDER SANDBOX |
| J-W-CREATEDREAD | cmd/ailang-worldd/operator_worktree.go:99 | `created := addErr == nil && readErr == nil && got.state == treeExact` → `created := addErr == nil && got.state == treeExact` | cmd/ailang-worldd / TestOperatorWorktreeWindow/readback_error | SURVIVED | SURVIVED (equivalent) | UNINFORMATIVE UNDER SANDBOX |
| N-STORE-RAW-CRED | cmd/ailang-worldd/operator_start_input.go:32 | `served, err := store.CanonicalPath(req.credPath)` → `served, err := req.credPath, error(nil)` | cmd/ailang-worldd / TestOperatorStartFacts | S5d equivalent | SURVIVED (equivalent) | UNINFORMATIVE UNDER SANDBOX |
| N-RES-ABS | cmd/ailang-worldd/operator_mint.go:95 | `!filepath.IsAbs(path) \|\| filepath.Clean(path) != path` → `filepath.Clean(path) != path` | cmd/ailang-worldd / TestOperatorMintRecovery | S5d equivalent | SURVIVED (equivalent) | UNINFORMATIVE UNDER SANDBOX |
| M1-ATTENDED-DISPATCH | cmd/ailang-worldd/operator_service.go:290 | `e.reconcileRunnerWorktree(args[1], out, errOut, readRunnerWorktree)` → `e.reconcileRunnerWorktree(args[0], out, errOut, readRunnerWorktree)` | cmd/ailang-worldd / TestOperatorWorktreeWindow/before_add | M1 hunk | KILLED | UNINFORMATIVE UNDER SANDBOX |
| M1-FIRST-READ-KEY | cmd/ailang-worldd/operator_worktree_reconcile.go:58 | `key, target, episode, err := readReservation(c.Credentials, id)` → `key, target, episode, err := readReservation(c.Credentials, id+"-other")` | cmd/ailang-worldd / TestOperatorWorktreeWindow/before_add | M1 hunk | KILLED | UNINFORMATIVE UNDER SANDBOX |
| M1-LOCKED-FRESH-READ | cmd/ailang-worldd/operator_worktree_reconcile.go:67 | `current, _, currentEpisode, err := readReservation(c.Credentials, id)` → `current, currentEpisode, err := key, episode, error(nil)` | cmd/ailang-worldd / TestOperatorWorktreeLockedReread | M1 hunk | KILLED | UNINFORMATIVE UNDER SANDBOX |

Equivalent predicates after replay

- J-R-CURRENTEP. An unchanged key digest binds the exact canonical episode and decoder enforces row/episode agreement. Key drift still triggers the independent key check.
- J-W-CREATEDREAD. All inspect error exits precede treeExact. Successful exact and last Git read failure are opposing controls.
- N-STORE-RAW-CRED. Existing S5d equivalent under intact configured paths. credentialsFor canonicalizes before request construction at session.go:129 and session_new.go:247.
- N-RES-ABS. Existing S5d equivalent. runnerPath checks the absolute clean parent at operator_service.go:62 after inspectRunnerMintToken pins the basename at operator_mint.go:95.

J-ST-KIND replays the judge's original predicate at base. The M1 anchor deletes the relocated row-kind guard. J-RB-COMMONID and J-W-EPPATTERN retain a compile-safe use of the operand/import; the remaining target comparison is removed.

### 4.1 Whole-package runs for the controller

Run `/bin/bash "$S/whole-mutants.sh" "$W" "$S/pkg-cmd-ailang-worldd.base-fails"` after M1 commits. W defaults to the sprint checkout and may name a controller-supplied clean checkout at M1. It runs the IDs in `$S/whole-mutant-ids.txt` in that Git checkout, one at a time, with exact content restore. Each row runs `go test -json -count=1 -timeout 8m ./cmd/ailang-worldd`. A kill requires a new failure of the named behavioural test against the recorded same-rig base set. An equivalent requires the same base set and a named PASS. Build failure, timeout, missing named result or bind denial invalidate attribution.

CLI rows for the external run (25)

````text
J-W-STORECHK
J-W-EPPATTERN
J-W-PHYS-ABS
J-R-STORECHK
J-R-BRANCHCHK
J-R-PATHCHK
J-R-ACTIONCHK
J-R-BASECHK
J-R-HEADLESS
J-ST-READPHASE
J-RB-COMMONID
J-RB-WTCOMMON
J-RB-LINKREG
J-RB-ISLINKED
J-RB-REGCOUNT
J-RB-ISDIR
J-RB-SYMLINKEVAL
J-R-CURRENTKEY
J-R-CURRENTEP
J-W-CREATEDREAD
N-STORE-RAW-CRED
N-RES-ABS
M1-ATTENDED-DISPATCH
M1-FIRST-READ-KEY
M1-LOCKED-FRESH-READ
````

## 5. Files

The measured M1 product paths are below. Reproduce with `python3 "$S/stage-diff.py" M1 M2`.

````text
1   1  cmd/ailang-worldd/operator_service.go
315 0  cmd/ailang-worldd/operator_worktree_guards_test.go
4   3  cmd/ailang-worldd/operator_worktree_reconcile.go
2   2  host/store/operator_worktree.go
167 0  host/store/operator_worktree_guards_test.go
````

`$S/status-M1` records their snapshot statuses. M1 moves one reservation-kind predicate, threads the real reservation reader through the existing attended entry, and adds store and CLI guard tests. §6 measures route/served bytes, package scope, lock flags and production callers.

## 6. Sprint gates

````bash gates.sh
#!/bin/bash
# Usage: gates.sh <tree> <logdir> <base|final|scratch> [sha]
set -u
set -o pipefail
export GIT_OPTIONAL_LOCKS=0
T=$1; L=$2; MODE=$3
B=ec4565ec5eedfed693f9ea5f2c78b89b63bf432e
REF=${GATES_GIT_REF:-$T}
S="$HOME/.ailang/state/world-iter435/plan-scratch"
SCRIPTS=$(cd "$(dirname "$0")" && pwd)
SHA=${4:-$(git -C "$REF" rev-parse HEAD)}
mkdir -p "$L"
cd "$T" || exit 2
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export AILANG_STATE_DIR="$L/ailang-state" AILANG_CACHE_DIR="$L/ailang-cache"
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR WORLD_UPDATE_OPERATOR_ACTION_GOLDEN WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN
mkdir -p "$AILANG_STATE_DIR"
: > "$L/legs.txt"
FAILS=
UNKNOWN=
s0=$(date +%s)
line() { echo "$*" | tee -a "$L/legs.txt"; }
leg() {
  local name=$1 cmd=$2 start rc
  start=$(date +%s)
  (set -eo pipefail; eval "$cmd") 2>&1 | tee "$L/$name.txt"
  rc=${PIPESTATUS[0]}
  if [ "${GATES_SANDBOX:-0}" = 1 ] &&
     grep -Eq 'bind: operation not permitted|OBSERVATION REQUIRES CONTROLLER|ARCHIVE REQUIRES CONTROLLER' "$L/$name.txt"; then
    line "LEG $name rc=$rc $(( $(date +%s)-start ))s UNINFORMATIVE UNDER SANDBOX"
    UNKNOWN="$UNKNOWN $name"
  else
    line "LEG $name rc=$rc $(( $(date +%s)-start ))s"
    [ "$rc" -eq 0 ] || FAILS="$FAILS $name"
  fi
}
observe() {
  if [ "${GATES_SANDBOX:-0}" = 1 ]; then
    echo 'OBSERVATION REQUIRES CONTROLLER (operator file observations)'
    return 77
  fi
  local d file creds
  for d in "$HOME/Library/LaunchAgents" "$HOME/.ailang-world/bin" "$HOME/.ailang-world"; do
    if [ ! -d "$d" ]; then printf 'absent %s\n' "$d"; continue; fi
    echo "DIRECTORY $d"
    if [ "$d" = "$HOME/.ailang-world" ]; then
      find "$d" -maxdepth 1 \( -type f -o -type l \) -name '*.json' -print0
    else
      find "$d" \( -type f -o -type l \) -print0
    fi | LC_ALL=C sort -z | while IFS= read -r -d '' file; do
      /bin/ls -la -d "$file"
      if [ -L "$file" ]; then readlink "$file" | shasum -a 256; else shasum -a 256 "$file"; wc -l "$file"; fi
    done
    if [ "$d" = "$HOME/.ailang-world" ]; then
      for file in "$d/service.json" "$d/operator-runner.json"; do
        if [ ! -e "$file" ] && [ ! -L "$file" ]; then printf 'absent %s\n' "$file"; fi
      done
    fi
  done
  creds=${GATES_REAL_CREDS:-$HOME/.ailang/world/credentials.db}
  for file in "$creds" "$creds".operator-episode-*.lock; do
    if [ -f "$file" ]; then shasum -a 256 "$file"; wc -l "$file"
    else printf 'absent %s\n' "$file"; fi
  done
  echo 'Service process write scope is checked with fixture supervisor counters'
}
leg observation-before 'observe'
paths() {
  case "$MODE" in
    base) git -C "$REF" diff --name-only "$B" "$B" ;;
    scratch) cat "$GATES_CHANGED_FILE" ;;
    final) git -C "$REF" diff --name-only "$B" HEAD ;;
    *) return 2 ;;
  esac
}
leg paths 'paths'
TOUCHED=$(awk -F/ '/\.go$/ && ($1=="cmd" || $1=="host") {NF--; OFS="/"; $1=$1; print}' "$L/paths.txt" | sort -u | tr '\n' ' ')
EXPECTED='cmd/ailang-worldd host/store '
line "TOUCHED_PACKAGES: $TOUCHED"
leg touched 'if [ "$MODE" = base ]; then test -z "$TOUCHED"; else test "$TOUCHED" = "$EXPECTED"; fi'
leg status 'git -C "$REF" status --porcelain
  if [ "$MODE" = final ]; then
    test -z "$(git -C "$REF" status --porcelain -- cmd host docs app)"
  fi' 
leg gofmt 'out=$(gofmt -l cmd host); printf "%s\n" "$out"
  if [ -n "$out" ]; then
    test "$out" = host/daemon/daemon.go
    git -C "$REF" show "$B:host/daemon/daemon.go" | cmp - "$out"
    echo BASELINE-FORMAT-DRIFT
  fi' 
leg build 'python3 "$S/bounded.py" 180 go build ./...'
leg vet 'python3 "$S/bounded.py" 180 go vet ./...'
leg fence "go test -run '^\$' ./..."
leg verify_ail 'python3 "$S/bounded.py" 300 ./scripts/verify_ail.sh'
PACKAGES=$TOUCHED
if [ "$MODE" = base ]; then PACKAGES=$EXPECTED; fi
test_set() {
  local tag=$1 test_rc=0
  shift
  python3 "$S/bounded.py" 510 go test -json -count=1 -timeout 8m "$@" > "$L/$tag.jsonl" 2>&1 || test_rc=$?
  if grep -Fq 'bind: operation not permitted' "$L/$tag.jsonl"; then
    grep -m 1 -F 'bind: operation not permitted' "$L/$tag.jsonl"
    return 77
  fi
  python3 "$S/failset.py" "$L/$tag.jsonl" > "$L/$tag.fails" || return 1
  cat "$L/$tag.fails"
  if [ "$test_rc" -gt 1 ] || grep -Eq '\[build failed\]|test timed out|BOUND EXCEEDED' "$L/$tag.jsonl"; then return 1; fi
  if [ "${GATES_ARCHIVE_ONLY:-0}" = 1 ]; then
    echo 'ARCHIVE REQUIRES CONTROLLER for whole CLI package verification'
    return 77
  fi
  case "$MODE" in
    base)
      grep -F 'cmd/ailang-worldd ' "$S/base-fails.expected" > "$L/$tag.expected"
      diff -u "$L/$tag.expected" "$L/$tag.fails"
      ;;
    final)
      cmp "$S/$tag.base-fails" "$L/$tag.fails"
      ;;
    *) return 2 ;;
  esac
}
for p in $PACKAGES; do
  tag=${p//\//-}
  if [ "$p" = cmd/ailang-worldd ]; then
    leg "pkg-$tag" "test_set pkg-$tag ./$p"
    leg "race-$tag" "test_set race-$tag -race ./$p"
  else
    leg "pkg-$tag" "go test -count=1 -timeout 8m ./$p"
    leg "race-$tag" "/usr/bin/time -p go test -count=1 -race -timeout 8m ./$p"
  fi
done
inv() {
  local p=$1 name n t
  shift
  n=$(printf '%s\n' "$*" | tr ' ' '|')
  name="inv-${p//\//-}"
  leg "$name" "go test -v -count=1 -timeout 8m ./$p -run '^($n)\$'"
  for t in "$@"; do
    if ! grep -qE "^--- PASS: $t " "$L/$name.txt"; then
      line "NOT PASSED ($name): $t"
      FAILS="$FAILS ${name}_missing"
    fi
  done
}
inv host/approval TestApprovalImportsOnlyTheStandardLibrary TestOperatorRunnerRecheckLaw TestOperatorPureEntryPointsWriteNothing
inv host/approval/approvaltest TestApprovalTestIsTestOnly TestTestSignerAgreesWithSSHKeygen
inv host/authority TestVerifiedOperatorActionOnlyFromTheJudges TestOperatorRunnerRechecks TestOperatorJudgesWriteNothing TestOperatorStoreWritesNothing
inv host/boundary TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard
inv host/daemon TestDaemonDependencyAllowlist TestDaemonCredentialWriterSites TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestConfirmRoutesRegistered TestSurfaceVerbRouteRegistered TestSurfaceVerbIsReadOnly TestSurfaceVerbOutcomes
inv host/verbs TestOperatorNeverDispatch TestOperatorOtherVerbParity
inv host/broker TestBrokerDependencyAllowlist
inv cmd/ailang-worldd TestEveryVerbHasHelp TestVerbTableMatchesTheCLISwitch TestCLICredentialWriterSites TestOperatorFlagParity TestLaunchAndStartHelp TestQuickstartSection19MatchesTheBar TestQuickstartSection20MatchesTheFinishFlag TestQuickstartServiceSectionMatchesTheCLIHelp TestPluginCommandsMatchTheCLIHelp TestHelpIsNeverBlockedByTheRegistryKeyGuard TestHomeServiceSeamOnly TestHomeServiceInstallAndPlist TestHomeServiceUpgradeRollback TestOperatorServiceFence TestOperatorServiceOwnership TestOperatorServiceConfig TestOperatorServicePair TestOperatorServiceStartup TestOperatorServiceUpgrade TestOperatorServiceUpgradePairFence TestOperatorServiceWriteScope TestOperatorServiceEnvironment
leg inv-quickstart-section9 "go test -v -count=1 -timeout 8m ./cmd/ailang-worldd -run '^TestQuickstartSection9FlagsMatchTheCLI\$'"
if [ "${GATES_SANDBOX:-0}" != 1 ] && ! grep -Fq -- '--- PASS: TestQuickstartSection9FlagsMatchTheCLI (' "$L/inv-quickstart-section9.txt"; then
  FAILS="$FAILS inv-quickstart-section9_missing"
fi
inv host/store TestProductionContextRoots TestProductionGoSurface TestStoreSurfaceTakesContext TestCredentialsMethodsRequireADeadline TestCredentialsSchemaMatchesCanonicalDDL TestOperatorSchema10 TestOperatorClaimBeforeExpiry TestOperatorClaimRace TestOperatorTerminalExclusion TestOperatorFinishReadiness TestWriterLockCrossProcessContention TestWriterLockHeldIsNonMutating
route_fence() {
  local routes posts old_routes old_posts
  routes=$(grep -Ec 'mux\.(HandleFunc|Handle)\(' host/daemon/daemon.go)
  posts=$(grep -Fc 'mux.HandleFunc("POST /surface/' host/daemon/daemon.go)
  old_routes=$(git -C "$REF" show "$B:host/daemon/daemon.go" | grep -Ec 'mux\.(HandleFunc|Handle)\(')
  old_posts=$(git -C "$REF" show "$B:host/daemon/daemon.go" | grep -Fc 'mux.HandleFunc("POST /surface/')
  printf 'DAEMON_ROUTES: %s (base %s); SURFACE_POSTS: %s (base %s)\n' "$routes" "$old_routes" "$posts" "$old_posts"
  test "$routes" -gt 0
  test "$routes" = "$old_routes"
  test "$posts" -gt 0
  test "$posts" = "$old_posts"
  local f
  for f in host/daemon/daemon.go host/daemon/surface_page.go host/surfacepage/app/app.js; do
    git -C "$REF" show "$B:$f" | cmp - "$f"
    printf 'UNCHANGED: %s\n' "$f"
  done
  grep -n 'surfaceCSP' host/daemon/surface_page.go
}
leg routes 'route_fence'
inv cmd/world-publish TestPluginCommandsMatchTheCLIHelp
leg mint-quickstart "go test -v -count=1 -timeout 1m ./cmd/ailang-worldd -run '^(TestQuickstartRunnerMintRecovery|TestOperatorMintRecovery)$'"
leg start-input "go test -v -count=1 -timeout 2m ./cmd/ailang-worldd -run '^(TestOperatorStartFacts|TestOperatorStartChanged)$'"
if [ "$MODE" = base ]; then
  line 'LEG s5e-accept SKIPPED 0s (base has no S5e tests)'
  line 'LEG added-cli-race SKIPPED 0s (base has no S5e tests)'
else
  leg s5e-accept '/bin/bash "$SCRIPTS/accept.sh" "$T" M1 fast'
  leg added-cli-race 'go test -race -v -count=1 -timeout 2m ./cmd/ailang-worldd -run "^(TestOperatorWorktreeSurvivors|TestOperatorWorktreeReadbackGuards|TestOperatorWorktreeRecoveryGuards|TestOperatorReadbackErrorIsNeverCreated|TestOperatorWorktreeLockedReread)$"'
fi
scope_fence() {
  local allowed="$L/allowed.txt" extra
  cat "$S/allowed.txt" > "$allowed"
  printf '%s\n' cmd/ailang-worldd/operator_worktree_guards_test.go | grep -Fx -f "$allowed"
  printf '%s\n' cmd/ailang-worldd/episode_preset.go | grep -Fxv -f "$allowed"
  extra=$(grep -Fxv -f "$allowed" "$L/paths.txt" || test "$?" -eq 1)
  printf 'UNTOUCHED_VIOLATIONS: %s\n' "${extra:-none}"
  test -z "$extra"
}
leg scope 'scope_fence'
production() {
  grep -R -nF --include='*.go' "$1" cmd host | grep -Fv '_test.go:'
}
callers() {
  production "$1(" | grep -Ev "func (\([^)]*\) )?$1\(" || :
}
one_site() {
  test "$(printf '%s\n' "$1" | wc -l | tr -d ' ')" = 1
  printf '%s\n' "$1" | grep -F "$2:"
}
caller_fence() {
  local count hits name
  production 'authority.RecheckVerb(' | grep -F 'cmd/ailang-worldd/episode_finish_surface.go:'
  production 'func MintOperatorStart('
  production ' ClaimOperatorAction('
  production 'func reconstructOperatorStart('
  for name in MintOperatorStart ClaimOperatorAction reserveOperatorStartWorktree; do
    test -z "$(callers "$name")"
    printf 'Production entry callers of %s: 0\n' "$name"
  done
  count=$(grep -R -nE --include='*.go' 'syscall\.Flock.*LOCK_NB' cmd host | grep -Fv '_test.go:' | wc -l | tr -d ' ')
  printf 'Production Flock LOCK_NB lines (with comment): %s\n' "$count"
  test "$count" = 7
  hits=$(callers provisionEpisode)
  printf '%s\n' "$hits"
  test "$(printf '%s\n' "$hits" | wc -l | tr -d ' ')" = 5
  one_site "$(production '.CommitOperatorMint(')" host/authority/mint.go
  one_site "$(production '.AbortOperatorMint(')" cmd/ailang-worldd/operator_mint.go
  one_site "$(callers operatorStartFacts)" cmd/ailang-worldd/operator_start_input.go
  production 'resolvePresetStart' | grep -F 'cmd/ailang-worldd/episode_preset.go:'
  production 'resolvePresetStart' | grep -F 'cmd/ailang-worldd/operator_start_input.go:'
  one_site "$(callers reconstructOperatorStart)" cmd/ailang-worldd/operator_worktree.go
  one_site "$(production '.ReserveOperatorWorktree(')" cmd/ailang-worldd/operator_worktree.go
  one_site "$(production '.ReadOperatorWorktree(')" cmd/ailang-worldd/operator_worktree_reconcile.go
  one_site "$(production '.reconcileRunnerWorktree(')" cmd/ailang-worldd/operator_service.go
  test -z "$(callers publishOperatorStartFiles)"
  test -z "$(callers reconcileRunnerFiles)"
  test -z "$(callers readOperatorFileSet)"
  printf '%s\n' 'Production file publication, readback and disposition entry callers: 0' 
}
leg callers 'caller_fence'
verify_go_set() {
  local rc=0
  GOFLAGS="${GOFLAGS:-} -json" VERIFY_GO_LEGS=plain python3 "$S/bounded.py" 1100 /bin/bash ./scripts/verify_go.sh > "$L/verify_go.raw" 2>&1 || rc=$?
  if grep -Fq 'bind: operation not permitted' "$L/verify_go.raw"; then
    grep -m 1 -F 'bind: operation not permitted' "$L/verify_go.raw"
    return 77
  fi
  python3 "$S/failset.py" "$L/verify_go.raw" > "$L/verify_go.fails"
  cat "$L/verify_go.fails"
  grep -Fq '── go test ./... -count=1' "$L/verify_go.raw" || return 1
  if [ "$rc" -gt 1 ] || grep -Eq '\[build failed\]|test timed out|BOUND EXCEEDED' "$L/verify_go.raw"; then return 1; fi
  cmp "$S/verify_go.base-fails" "$L/verify_go.fails"
}
if [ "$MODE" != final ]; then
  line 'LEG verify_go SKIPPED 0s (tracked-blob gate needs a git checkout)'
else
  leg verify_go verify_go_set
fi
leg observation-after 'observe'
if ! cmp -s "$L/observation-before.txt" "$L/observation-after.txt"; then
  diff -u "$L/observation-before.txt" "$L/observation-after.txt" | tee "$L/observation-diff.txt"
  FAILS="$FAILS observation_changed"
fi
line "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none} UNINFORMATIVE:${UNKNOWN:- none}"
tail -n 1 "$L/legs.txt" > "$L/ALL_DONE.$SHA"
cp "$L/ALL_DONE.$SHA" "$L/ALL_DONE"
[ -z "$FAILS" ] || exit 1
[ -z "$UNKNOWN" ] || exit 77
````

## 7. Executor directive notes

Extract the three exact fences as in §2. The named acceptance tests use no sockets and are all informative in the planner sandbox. Full CLI and daemon packages stay out of that command. A bind denial is an uninformative acceptance run and blocks the snapshot.

````bash accept.sh
#!/bin/bash
set -euo pipefail
W=$1; M=$2; PART=$3
export AILANG_BIN="$HOME/.pinned-ailang/ailang"
export WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
S="$HOME/.ailang/state/world-iter435/plan-scratch"
export AILANG_STATE_DIR="$S/accept-state" AILANG_CACHE_DIR="$S/accept-cache"
mkdir -p "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR"
cd "$W"
test "$M" = M1
test "$PART" = fast
run() {
  local pkg=$1 names=$2 name log rc=0
  log=$(mktemp "$S/accept.XXXXXX")
  python3 "$S/bounded.py" 150 go test -v -count=1 -timeout 2m "./$pkg" -run "^($names)\$" 2>&1 | tee "$log" || rc=$?
  if grep -Fq 'bind: operation not permitted' "$log"; then
    echo 'UNINFORMATIVE UNDER SANDBOX'
    exit 77
  fi
  test "$rc" -eq 0
  while IFS= read -r name; do
    grep -Fq -- "--- PASS: $name (" "$log"
  done < <(printf '%s\n' "$names" | tr '|' '\n')
}
run host/store 'TestReserveOperatorWorktreeGuards|TestDisposeOperatorUncertainStepGuards|TestOperatorStepReceiptAtLimit|TestOperatorStepRestart|TestOperatorStepGuards|TestCredentialsMethodsRequireADeadline|TestOperatorSchema10'
run cmd/ailang-worldd 'TestOperatorWorktreeSurvivors|TestOperatorWorktreeRecoveryGuards|TestOperatorWorktreeReadbackGuards|TestOperatorReadbackErrorIsNeverCreated|TestOperatorWorktreeLockedReread|TestOperatorWorktreeWindow|TestOperatorStartOnce|TestOperatorWorktreeDispositionCAS|TestEveryVerbHasHelp|TestQuickstartServiceSectionMatchesTheCLIHelp'
printf 'ACCEPT_OK %s %s\n' "$M" "$PART"
````

Snapshot path is `$HOME/.ailang/state/world-iter435/snap/M1`. Read-only Git lists modified/tracked and added/untracked product paths relative to ec4565ec. The script writes `STATUS`, `DELETED` and `MANIFEST.sha256`. Its product pathspec supports absent roots and excludes plan records. Use the scratch overrides only for archive rehearsal.

````bash snapshot.sh
#!/bin/bash
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=$1
M=$2
B=ec4565ec5eedfed693f9ea5f2c78b89b63bf432e
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter435/snap}/$M
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

````bash
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter435-row243-v4s5e
T=$HOME/.ailang/state/world-iter435/executor
S=$HOME/.ailang/state/world-iter435/plan-scratch
STATE=$HOME/.ailang/state/world-iter435
# Before application, the controller supplies or uses the pristine base Git checkout.
BASE_TREE=${BASE_TREE:-$W}
/bin/bash "$S/baseline-controller.sh" "$BASE_TREE"

````

After the executor snapshots and the plan commit exists, run this block sequentially.

````bash
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter435-row243-v4s5e
T=$HOME/.ailang/state/world-iter435/executor
S=$HOME/.ailang/state/world-iter435/plan-scratch
STATE=$HOME/.ailang/state/world-iter435
mkdir -p "$T"
printf '%s\n' 'test(operator): pin reservation guards before file effects' > "$T/M1.msg"
bash "$STATE/commit_ms.sh" M1 "$T/M1.msg" 'AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang /bin/bash "'"$T"'/accept.sh" "'"$W"'" M1 fast'
/bin/bash "$T/gates.sh" "$W" "$S/gate-controller-final" final
/bin/bash "$S/whole-mutants.sh" "$W" "$S/pkg-cmd-ailang-worldd.base-fails"
````

The controller round trip `/bin/bash "$STATE/ctl_roundtrip.sh" M1` extracted the plan's exact diff fence and applied it to a pristine export. `diff -r` on cmd, host, docs and the full archive returned 0. App was absent. No new rejects or backups appeared. Evidence is in `$S/roundtrip.log` and the empty `$S/roundtrip-M1.diff`. The extracted three scripts matched their scratch copies byte-for-byte and passed `bash -n`.

For the snapshot rehearsal, the exact extracted script used `GIT_DIR=$(git rev-parse --absolute-git-dir)`, `GIT_WORK_TREE="$STATE/ctl-roundtrip/export"`, `GIT_OPTIONAL_LOCKS=0` and `SNAP_PARENT="$S/snap-rehearsal"`. No supplied status list was used. STATUS matched `$S/status-M1`, the counts were **5/0/5** and each hash verified. See `$S/snapshot-M1.log`.

## 8. Risks

- A detached decoy masks the common-directory defence. The fixture proves both correct SHA and symbolic HEAD before the target is inspected.
- The controller compares real-binary fail sets in Git checkouts. The base archive cannot establish a rig-specific e2e verdict, and Linux CI is a separate gate.
- The file slice requires effect-specific disposition and a current-owner/step CAS in its own patch. Deferral keeps new file uncertainty outside the shipped surface.

## 9. Acceptance criteria

1. M1 applies and byte-matches its compiled stage. Cap result is 495 changed product lines and the snapshot has the expected 5/0/5 records.
2. Every non-equivalent judge survivor has a named behavioural kill. The base replay reproduced all 30 judge survivor rows. Measured equivalents state their masking invariants.
3. Stable tampering and observed-created reservations refuse before the attended prompt. Readback and locked re-read retain Git/file/state evidence. The activation caller control and named service gates pass.
4. Controller runs same-rig whole-package fail-set comparisons and Linux CI. Appendix C gates remain in force before any production consumer.

## Appendix C — after this sprint (measured)

- **File publication and disposition as one milestone.** The compiled scratch increment is **967 changed lines**, **1,462 cumulative** with M1. Unfiltered product sum is 1,462. See `$S/numstat.txt`, `$S/st-M2` and `$S/verify-M2-probe.log`. The first full-shape draft was 1,375 cumulative and lacked the separate admission matrix and linked-directory refusal. Adding those gates put it 62 lines over cap. The next sprint can build on committed M1 with an independent 967-line footprint. The remaining branch matrix and Linux implementation cost are UNMEASURED. Keep the full signed three-file case together.
- **File prerequisite and effect trace.** Under the canonical episode lock, require an exact observed-created worktree, matching observed mint with current owner, live action/credential expiry, signed task bytes and canonical paths, and exact in-memory token/credential match. Derive MCP bytes from the approved address and pinned preset. Reserve one exact identity with the step bank before any write. Use private exclusive no-replace installation and independent directory-relative path/type/UID/mode/link-count/digest readback for each signed output. A known-absent set can stop by attended review. Partial, inaccessible or foreign evidence is uncertain. Keep the worktree, credential and start evidence.
- **File prototype tests exercised.** `TestOperatorFilesReconcile/{no_mcp,all_files,before_reservation,before_first,between_task_token,between_token_mcp,before_readback,before_observe,competing_file,changed_file,changed_during_review,owner_changed,version_changed,later_step,observe_error}` reads exact store, file, mode, hash and Git evidence. It reads the QUICKSTART example and executes the emitted `ailang-worldd service runner reconcile-files oa-0123456789abcdef0123456789abcdef`. `TestOperatorFileReadbackMetadata` covers owner, mode, digest, symlink, hardlink, FIFO and read bound. `TestOperatorFilePublicationAdmission` covers missing mint, wrong worktree receipt, mint owner, exact expiry, token mismatch, existing exact file, linked or public parent, moved checkout and held lock.
- **File gates before the next runnable patch.** Compare pre/post effect readbacks and independently CAS the action owner and step under the write transaction, including drift after step update. Require an existing reserved row before any prompt. Recheck the approved signed start and mint binding at publication. Exercise the remaining reopen and expiry arms, changed file/parent, bitwise open/lock flags, path/UID/mode/link/digest and identity comparisons, timing bounds and refusal branches. Prove new preparation/readback remain dormant except the resolver's CLI reads. Publish the tested remedy with every new uncertainty. Run `go build ./...`, `go vet ./cmd/ailang-worldd ./host/store`, the three named test families, §6 inventories, the diff-anchored mutation matrix and Linux CI. No prototype patch is in Appendix P.
- **S5f/S6 paired-process consumer, owner recovery and drain.** Carry `E-RES-EPISODE-CONFLICT`, `E-RES-REREAD`, `J-OWN-UID`, `Z-ABSENT-ANY`, `Z-WAIT-READBUDGET`, `TestOperatorFirstManagedEpisodeProcess`, `TestOperatorOldBuildTerminalDrain`, `TestOperatorCurrentOwnerProcessRecovery` and `TestOperatorServiceUpgrade/active_child`. The M1 deterministic re-read seam serves the worktree resolver. These other pins remain with their effect-specific prerequisites. S5d C12 recorded a 72-line prior FD harness. The referenced iter433 tree was absent and full size remains UNMEASURED. Require canonical lock and current session, pinned executable, confinement, terminal session-revoke exclusion and controlled drain.
- **Launch, browser, finish and Herdr.** Design §8 raw budgets are S6 890, S8 860 and S12 700. Full implementations remain UNMEASURED. Keep E3, E5/#604 app sequencing, E4 and D-WORLD-120/E2 readiness in their assigned slices. Row 271 stays with its se-tools baseline and design question.

## Appendix P — compiled patches

### Patch M1

````diff M1
diff --git a/cmd/ailang-worldd/operator_service.go b/cmd/ailang-worldd/operator_service.go
index fd2978d0..5274a6f1 100644
--- a/cmd/ailang-worldd/operator_service.go
+++ b/cmd/ailang-worldd/operator_service.go
@@ -287,7 +287,7 @@ func runServiceRunner(args []string, out, errOut io.Writer, e serviceEnv) int {
 		return e.reconcileRunnerMint(args[1], out, errOut)
 	}
 	if verb == "reconcile-worktree" {
-		return e.reconcileRunnerWorktree(args[1], out, errOut)
+		return e.reconcileRunnerWorktree(args[1], out, errOut, readRunnerWorktree)
 	}
 	p := e.runnerPaths()
 	if verb == "status" {
diff --git a/cmd/ailang-worldd/operator_worktree_guards_test.go b/cmd/ailang-worldd/operator_worktree_guards_test.go
new file mode 100644
index 00000000..60673eca
--- /dev/null
+++ b/cmd/ailang-worldd/operator_worktree_guards_test.go
@@ -0,0 +1,315 @@
+package main
+
+import (
+	"crypto/sha256"
+	"encoding/json"
+	"errors"
+	"io"
+	"os"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func TestOperatorWorktreeSurvivors(t *testing.T) {
+	r := newWorktreeStartRig(t, false)
+	for _, name := range []string{"store", "episode", "relative"} {
+		t.Run(name, func(t *testing.T) {
+			before := observeProvision(t, r.f.r, r.creds)
+			if name == "relative" {
+				if _, err := operatorPhysicalTarget("ep"); !errors.Is(err, store.ErrOperatorConflict) {
+					t.Fatalf("relative path admitted: %v", err)
+				}
+			} else {
+				a := r.action
+				if name == "store" {
+					a.Store += ".other"
+				} else {
+					a.Episode = "BAD!"
+				}
+				if err := reserveOperatorStartWorktree(presetTestContext(t), runnerDB(r.creds), a, r.at, r.src, func(string) error {
+					t.Fatal("invalid start crossed a worktree boundary")
+					return nil
+				}); err != store.ErrOperatorConflict {
+					t.Fatalf("expected identity refusal before lock and reconstruction: %v", err)
+				}
+			}
+			if !reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
+				t.Fatal("identity refusal modified store, Git or files")
+			}
+		})
+	}
+}
+
+func TestOperatorWorktreeRecoveryGuards(t *testing.T) {
+	for _, name := range []string{"store", "branch", "path", "action", "basename", "headless", "observed"} {
+		t.Run(name, func(t *testing.T) {
+			r := newWorktreeStartRig(t, false)
+			crash := errors.New("before add")
+			err := r.start(t, func(step string) error {
+				if step == "before add" && name != "observed" {
+					return crash
+				}
+				return nil
+			})
+			if name == "observed" {
+				if err != nil {
+					t.Fatal(err)
+				}
+				r.state(t, "running", "observed", "created", 2, true)
+			} else {
+				if !errors.Is(err, crash) {
+					t.Fatalf("expected reserved worktree: %v", err)
+				}
+				a := r.action
+				var target operatorWorktreeTarget
+				var identity string
+				if err := r.db.QueryRow("SELECT identity FROM operator_action_steps WHERE kind='worktree'").Scan(&identity); err != nil {
+					t.Fatal(err)
+				}
+				if err := json.Unmarshal([]byte(identity), &target); err != nil {
+					t.Fatal(err)
+				}
+				switch name {
+				case "store":
+					a.Store += ".other"
+				case "branch":
+					target.Branch += "-other"
+				case "path":
+					target.Path += "-other"
+				case "action":
+					target.Action += "-other"
+				case "headless":
+					a.Options.Headless = false
+					a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
+				case "basename":
+					target.Path += "-other"
+					a.Resolution = append([]approval.OperatorResolved(nil), a.Resolution...)
+					for i := range a.Resolution {
+						if a.Resolution[i].Name == "worktree_path" {
+							a.Resolution[i].Value += "-other"
+						}
+					}
+				}
+				body, err := approval.CanonicalOperatorAction(a)
+				if err != nil {
+					t.Fatal(err)
+				}
+				digest := sha256.Sum256(body)
+				encoded, err := json.Marshal(target)
+				if err != nil {
+					t.Fatal(err)
+				}
+				// Both reads see the same damaged rows. The CAS can agree with them.
+				r.exec(t, "UPDATE operator_actions SET action=?, digest=?", body, digest[:])
+				r.exec(t, "UPDATE operator_action_steps SET identity=? WHERE kind='worktree'", string(encoded))
+			}
+			before := observeProvision(t, r.f.r, r.creds)
+			r.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
+				t.Error("invalid worktree reservation reached the prompt")
+				return nil, errors.New("unexpected prompt")
+			}
+			want := store.ErrOperatorConflict
+			if name == "observed" {
+				want = store.ErrOperatorMissing
+			}
+			if code, msg := r.run(t, worktreeExample); code != exitUsage || !strings.Contains(msg, want.Error()) ||
+				!reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
+				t.Fatalf("invalid worktree reservation changed state or returned wrong class: %d %s", code, msg)
+			}
+		})
+	}
+}
+
+func TestOperatorWorktreeReadbackGuards(t *testing.T) {
+	for _, name := range []string{"repo_common", "checkout_common", "link_regular", "link_prefix", "registration_count", "directory", "physical"} {
+		t.Run(name, func(t *testing.T) {
+			f := newEpisodeRun(t)
+			ctx := presetTestContext(t)
+			root, err := filepath.EvalSymlinks(t.TempDir())
+			if err != nil {
+				t.Fatal(err)
+			}
+			repo, err := filepath.EvalSymlinks(f.f.repo)
+			if err != nil {
+				t.Fatal(err)
+			}
+			target := operatorWorktreeTarget{Repo: repo, Common: filepath.Join(repo, ".git"), Path: filepath.Join(root, "ep"),
+				Branch: "world/ep", SHA: presetGit(t, f.f, "rev-parse", "HEAD")}
+			git := f.f.git
+			setLink := func(dst, src string) {
+				t.Helper()
+				seedStartFixtureFile(t, dst, readStartFixtureFile(t, src))
+			}
+			link := filepath.Join(target.Path, ".git")
+			wrapper := func(body string) {
+				t.Helper()
+				git = filepath.Join(root, "git")
+				if err := os.WriteFile(git, []byte("#!/bin/sh\n"+body+"\nexec "+shellQuote(f.f.git)+" \"$@\"\n"), 0700); err != nil {
+					t.Fatal(err)
+				}
+			}
+			if name == "repo_common" {
+				target.Common = root
+				if _, err := target.inspect(ctx, git); err == nil || err.Error() != "worktree repository identity differs" {
+					t.Fatalf("initial inspect used another repository: %v", err)
+				}
+				if _, err := os.Lstat(target.Path); !errors.Is(err, os.ErrNotExist) {
+					t.Fatalf("initial inspection wrote checkout: %v", err)
+				}
+				return
+			}
+			presetGit(t, f.f, "worktree", "add", "-b", target.Branch, target.Path, target.SHA)
+			if e, err := target.inspect(ctx, git); err != nil || e.state != treeExact {
+				t.Fatalf("fixture checkout is not exact: %v", err)
+			}
+			switch name {
+			case "checkout_common":
+				other := f.f
+				other.repo = filepath.Join(root, "other")
+				presetGit(t, f.f, "clone", "--quiet", repo, other.repo)
+				decoy := filepath.Join(root, "decoy")
+				presetGit(t, other, "worktree", "add", "-b", target.Branch, decoy, target.SHA)
+				setLink(link, filepath.Join(decoy, ".git"))
+				checkout := f.f
+				checkout.repo = target.Path
+				if head := presetGit(t, checkout, "rev-parse", "HEAD", "--symbolic-full-name", "HEAD"); head != target.SHA+"\nrefs/heads/"+target.Branch {
+					t.Fatalf("decoy would be caught by HEAD: %q", head)
+				}
+			case "link_regular", "link_prefix":
+				saved := filepath.Join(root, "saved-link")
+				setLink(saved, link)
+				if name == "link_regular" {
+					if err := os.Remove(link); err != nil {
+						t.Fatal(err)
+					}
+					if err := os.Symlink(saved, link); err != nil {
+						t.Fatal(err)
+					}
+					checkout := f.f
+					checkout.repo = target.Path
+					if head := presetGit(t, checkout, "rev-parse", "HEAD"); head != target.SHA {
+						t.Fatalf("symlink fixture has no real Git checkout: %q", head)
+					}
+				} else {
+					// Different reads can see different link contents. The later real Git
+					// invocation would succeed if the early gitdir prefix check were skipped.
+					seedStartFixtureFile(t, link, []byte("interrupted gitfile\n"))
+					wrapper("if [ \"$2\" = " + shellQuote(target.Path) + " ]; then cp " + shellQuote(saved) + " " + shellQuote(link) + " || exit 1; fi")
+				}
+			case "registration_count":
+				second := filepath.Join(root, "second")
+				presetGit(t, f.f, "worktree", "add", "--detach", second, target.SHA)
+				secondGitDir := strings.TrimPrefix(strings.TrimSpace(string(readStartFixtureFile(t, filepath.Join(second, ".git")))), "gitdir: ")
+				seedStartFixtureFile(t, filepath.Join(secondGitDir, "gitdir"), []byte(link+"\n"))
+			case "directory", "physical":
+				// List the original registration, then present the changed object after
+				// that read. Both branches of this fixture run the actual Git command.
+				moved := filepath.Join(root, "moved")
+				var replacement string
+				if name == "directory" {
+					replacement = "touch " + shellQuote(target.Path)
+				} else {
+					// Replace only the ancestor so the leaf's Lstat remains a directory.
+					parent := filepath.Join(root, "parent")
+					if err := os.Mkdir(parent, 0700); err != nil {
+						t.Fatal(err)
+					}
+					presetGit(t, f.f, "worktree", "move", target.Path, filepath.Join(parent, "ep"))
+					target.Path = filepath.Join(parent, "ep")
+					replacement = "ln -s " + shellQuote(moved) + " " + shellQuote(parent)
+				}
+				victim := target.Path
+				if name == "physical" {
+					victim = filepath.Dir(target.Path)
+				}
+				wrapper("if [ \"$3\" = for-each-ref ]; then mv " + shellQuote(victim) + " " + shellQuote(moved) + " && " + replacement + " || exit 1; fi")
+			}
+			e, err := target.inspect(ctx, git)
+			if name == "physical" {
+				if err == nil || err.Error() != "worktree physical path differs" {
+					t.Fatalf("redirected ancestor not refused: %v", err)
+				}
+			} else if err != nil || e.state != treePartial {
+				t.Fatalf("attribution accepted or wrong refusal: state=%d %v", e.state, err)
+			}
+			if sha := presetGit(t, f.f, "rev-parse", "refs/heads/"+target.Branch); sha != target.SHA {
+				t.Fatal("inspection changed pinned branch")
+			}
+			if name == "link_prefix" && string(readStartFixtureFile(t, link)) != "interrupted gitfile\n" {
+				t.Fatal("unchecked link reached later Git read")
+			}
+			if name == "registration_count" && strings.Count(e.registration, "\x00\x00") != 2 {
+				t.Fatalf("expected two real registrations: %q", e.registration)
+			}
+		})
+	}
+}
+
+func TestOperatorReadbackErrorIsNeverCreated(t *testing.T) {
+	// The inspector's error returns are all before the one assignment of treeExact.
+	// Exercise a real command failure at the last possible Git observation.
+	f := newEpisodeRun(t)
+	repo, err := filepath.EvalSymlinks(f.f.repo)
+	if err != nil {
+		t.Fatal(err)
+	}
+	root, err := filepath.EvalSymlinks(t.TempDir())
+	if err != nil {
+		t.Fatal(err)
+	}
+	target := operatorWorktreeTarget{Repo: repo, Common: filepath.Join(repo, ".git"), Path: filepath.Join(root, "ep"),
+		Branch: "world/ep", SHA: presetGit(t, f.f, "rev-parse", "HEAD")}
+	presetGit(t, f.f, "worktree", "add", "-b", target.Branch, target.Path, target.SHA)
+	git := filepath.Join(root, "git")
+	script := "#!/bin/sh\nif [ \"$4\" = HEAD ]; then exit 1; fi\nexec " + shellQuote(f.f.git) + " \"$@\"\n"
+	if err := os.WriteFile(git, []byte(script), 0700); err != nil {
+		t.Fatal(err)
+	}
+	e, err := target.inspect(presetTestContext(t), git)
+	if err != nil || e.state != treePartial || e.common != target.Common {
+		t.Fatalf("failed HEAD read must be partial: %+v %v", e, err)
+	}
+	e, err = target.inspect(presetTestContext(t), f.f.git)
+	if err != nil || e.state != treeExact || e.head != target.SHA+"\nrefs/heads/"+target.Branch+"\n" {
+		t.Fatalf("control HEAD read must be exact: %+v %v", e, err)
+	}
+}
+
+func TestOperatorWorktreeLockedReread(t *testing.T) {
+	r := newWorktreeStartRig(t, false)
+	crash := errors.New("before add")
+	if err := r.start(t, func(step string) error {
+		if step == "before add" {
+			return crash
+		}
+		return nil
+	}); !errors.Is(err, crash) {
+		t.Fatalf("expected reserved worktree: %v", err)
+	}
+	var before provisionSnapshot
+	calls := 0
+	read := func(path, id string) (store.OperatorStepKey, operatorWorktreeTarget, string, error) {
+		key, target, episode, err := readRunnerWorktree(path, id)
+		calls++
+		if calls == 1 && err == nil {
+			r.exec(t, "UPDATE operator_actions SET generation=generation+1")
+			before = observeProvision(t, r.f.r, r.creds)
+		}
+		return key, target, episode, err
+	}
+	r.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
+		t.Error("a changed reservation reached the prompt")
+		return nil, errors.New("unexpected prompt")
+	}
+	var stderr strings.Builder
+	code := r.env.reconcileRunnerWorktree(r.action.ID, io.Discard, &stderr, read)
+	if code != exitUsage || calls != 2 || !strings.Contains(stderr.String(), store.ErrOperatorConflict.Error()) ||
+		!reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
+		t.Fatalf("reread did not preserve the replacement row and its evidence: %d %s", code, &stderr)
+	}
+}
diff --git a/cmd/ailang-worldd/operator_worktree_reconcile.go b/cmd/ailang-worldd/operator_worktree_reconcile.go
index eadf23c3..874a8ce2 100644
--- a/cmd/ailang-worldd/operator_worktree_reconcile.go
+++ b/cmd/ailang-worldd/operator_worktree_reconcile.go
@@ -45,7 +45,8 @@ func readRunnerWorktree(credentials, id string) (store.OperatorStepKey, operator
 	return key, target, signed.Episode, nil
 }
 
-func (e serviceEnv) reconcileRunnerWorktree(id string, out, stderr io.Writer) int {
+func (e serviceEnv) reconcileRunnerWorktree(id string, out, stderr io.Writer,
+	readReservation func(string, string) (store.OperatorStepKey, operatorWorktreeTarget, string, error)) int {
 	remedy := reconcileWorktreeCommand(id)
 	refuse := func(err error) int {
 		return serviceMessage(stderr, "worktree recovery refused", err.Error(), remedy, exitUsage)
@@ -54,7 +55,7 @@ func (e serviceEnv) reconcileRunnerWorktree(id string, out, stderr io.Writer) in
 	if err != nil {
 		return refuse(err)
 	}
-	key, target, episode, err := readRunnerWorktree(c.Credentials, id)
+	key, target, episode, err := readReservation(c.Credentials, id)
 	if err != nil {
 		return refuse(err)
 	}
@@ -63,7 +64,7 @@ func (e serviceEnv) reconcileRunnerWorktree(id string, out, stderr io.Writer) in
 		return refuse(errors.New("episode recovery lock held or unavailable"))
 	}
 	defer lock.Close()
-	current, _, currentEpisode, err := readRunnerWorktree(c.Credentials, id)
+	current, _, currentEpisode, err := readReservation(c.Credentials, id)
 	if err != nil {
 		return refuse(err)
 	}
diff --git a/host/store/operator_worktree.go b/host/store/operator_worktree.go
index 39f76344..c884e5cb 100644
--- a/host/store/operator_worktree.go
+++ b/host/store/operator_worktree.go
@@ -28,14 +28,14 @@ func (c *Credentials) ReserveOperatorWorktree(ctx context.Context, at OperatorMi
 	if now >= row.ExpiresAt {
 		return key, ErrOperatorExpired
 	}
-	if row.State != "claimed" || row.Owner != at.Runner.Nonce || len(row.Challenge) != 32 || identity == "" {
+	if row.Kind != "start" || row.State != "claimed" || row.Owner != at.Runner.Nonce || len(row.Challenge) != 32 || identity == "" {
 		return key, ErrOperatorConflict
 	}
 	action, err := decodeOperatorMintAction(row)
 	if err != nil {
 		return key, err
 	}
-	if c.lock == nil || action.Store != c.lock.dbPath || action.Kind != "start" || !action.Options.Headless {
+	if c.lock == nil || action.Store != c.lock.dbPath || !action.Options.Headless {
 		return key, ErrOperatorConflict
 	}
 	b, err := proof.Bind(action, [32]byte(row.Challenge), row.ExpiresAt, row.Artifact)
diff --git a/host/store/operator_worktree_guards_test.go b/host/store/operator_worktree_guards_test.go
new file mode 100644
index 00000000..01d61856
--- /dev/null
+++ b/host/store/operator_worktree_guards_test.go
@@ -0,0 +1,167 @@
+package store
+
+import (
+	"crypto/sha256"
+	"errors"
+	"reflect"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+)
+
+func TestReserveOperatorWorktreeGuards(t *testing.T) {
+	for _, name := range []string{"valid", "kind", "headless", "store", "approver", "challenge", "claimed"} {
+		t.Run(name, func(t *testing.T) {
+			c := openTestCredentials(t)
+			ctx := testCtx(t)
+			cas, runner := operatorClaimFixture(t, c)
+			if err := c.ClaimOperatorAction(ctx, cas, runner, 1590); err != nil {
+				t.Fatal(err)
+			}
+			row := getOperator(t, c, cas.ID, 1590)
+			action, err := approval.DecodeOperatorAction(row.Action)
+			if err != nil {
+				t.Fatal(err)
+			}
+			switch name {
+			case "kind":
+				action.Kind, action.Task, action.Options.Headless = "finish", "", false
+			case "headless":
+				action.Options.Headless = false
+			case "store":
+				action.Store += ".other"
+			}
+			action.Argv = approval.OperatorArgv(action.Kind, action.Episode, action.Options)
+			body, err := approval.CanonicalOperatorAction(action)
+			if err != nil {
+				t.Fatal(err)
+			}
+			digest, err := approval.OperatorSignatureDigest(action)
+			if err != nil {
+				t.Fatal(err)
+			}
+			challenge := [32]byte(row.Challenge)
+			message, err := approval.OperatorActionMessage(digest, challenge, row.ExpiresAt)
+			if err != nil {
+				t.Fatal(err)
+			}
+			signer := approvaltest.NewSigner(42)
+			ap := approval.Approver{Name: "fixture", Format: approval.FormatSSH, Key: signer.PublicKey(), EnrolledAt: 1000}
+			if err := c.EnrolApprover(ctx, ApproverRow{Name: ap.Name, Format: ap.Format, PublicKey: ap.Key.Blob, EnrolledAt: ap.EnrolledAt}); err != nil {
+				t.Fatal(err)
+			}
+			artifact := signer.Sign(approval.OperatorActionNamespace, message)
+			signed, err := approval.CheckOperatorAction(approval.OperatorCheckInput{Action: action, Challenge: challenge,
+				Expiry: row.ExpiresAt, Now: 1590, Artifact: artifact, Format: ap.Format, Approver: ap})
+			if err != nil {
+				t.Fatal(err)
+			}
+			proof := signed.MintProof()
+			cas.Digest, cas.Generation = digest, row.Generation
+			operatorExec(t, c, "UPDATE operator_actions SET action=?, digest=?, kind=?, artifact=?, approved_by=?", body, digest[:], action.Kind, artifact, signed.ApprovedBy())
+			switch name {
+			case "kind", "claimed":
+				if name == "claimed" {
+					// The SQL CAS also rejects running. Make reaching Bind observable.
+					proof = nil
+					operatorExec(t, c, "UPDATE operator_actions SET state='running'")
+				} else {
+					// A finish must refuse before decoding even when its stored bytes are damaged.
+					broken := body[:1]
+					cas.Digest = sha256.Sum256(broken)
+					operatorExec(t, c, "UPDATE operator_actions SET action=?, digest=?", broken, cas.Digest[:])
+				}
+			case "approver":
+				operatorExec(t, c, "UPDATE operator_actions SET approved_by='ssh:other'")
+			case "challenge":
+				// Model a damaged row; the original 32 bytes still bind to the proof.
+				operatorExec(t, c, "PRAGMA ignore_check_constraints = ON")
+				operatorExec(t, c, "UPDATE operator_actions SET challenge=?", append(row.Challenge, 1))
+				operatorExec(t, c, "PRAGMA ignore_check_constraints = OFF")
+			}
+			before := getOperator(t, c, cas.ID, 1590)
+			key, err := c.ReserveOperatorWorktree(ctx, OperatorMintAttempt{CAS: cas, Runner: runner, Catalogue: 1}, "fixture", proof, func() int64 { return 1590 })
+			if name == "valid" {
+				if err != nil {
+					t.Fatal(err)
+				}
+				assertReservedOperatorStep(t, c, key)
+				operatorCounts(t, c, 1, 2)
+				if after := getOperator(t, c, cas.ID, 1590); after.State != "running" || after.Generation != cas.Generation+1 || after.Owner != runner.Nonce {
+					t.Fatalf("incorrect admitted state: %s/%d/%s", after.State, after.Generation, after.Owner)
+				}
+				return
+			}
+			if !errors.Is(err, ErrOperatorConflict) || key != (OperatorStepKey{}) ||
+				!reflect.DeepEqual(before, getOperator(t, c, cas.ID, 1590)) {
+				t.Fatalf("refusal changed action or returned the wrong class: %v", err)
+			}
+			operatorCounts(t, c, 1, 1)
+			sessions, err := c.ListSessions(ctx, "", 20)
+			start, startErr := c.EpisodeStart(ctx, action.Episode)
+			if err != nil || len(sessions) != 0 || startErr != nil || start != nil {
+				t.Fatalf("reserve changed mint evidence: sessions=%d start present=%t: %v %v", len(sessions), start != nil, err, startErr)
+			}
+		})
+	}
+}
+
+func TestDisposeOperatorUncertainStepGuards(t *testing.T) {
+	for _, name := range []string{"kind", "owner", "version"} {
+		t.Run(name, func(t *testing.T) {
+			c := openTestCredentials(t)
+			ctx := testCtx(t)
+			cas, runner := operatorClaimFixture(t, c)
+			if err := c.ClaimOperatorAction(ctx, cas, runner, 1590); err != nil {
+				t.Fatal(err)
+			}
+			cas.Generation++
+			key := OperatorStepKey{cas, runner.Nonce, "worktree", "fixture", 1}
+			operatorExec(t, c, "UPDATE operator_actions SET state='running'")
+			if err := c.ReserveOperatorStep(ctx, key, 1590); err != nil {
+				t.Fatal(err)
+			}
+			before := getOperator(t, c, cas.ID, 1590)
+			bad := key
+			want := ErrOperatorConflict
+			switch name {
+			case "kind":
+				bad.Kind, want = "files", ErrOperatorInput
+			case "owner":
+				operatorExec(t, c, "UPDATE operator_action_steps SET owner='different' WHERE kind='worktree'")
+				key.Owner = "different"
+			case "version":
+				bad.Version++
+			}
+			if err := c.DisposeOperatorWorktree(ctx, bad, OperatorWorktreeUncertain, 1591); !errors.Is(err, want) ||
+				!reflect.DeepEqual(before, getOperator(t, c, cas.ID, 1590)) {
+				t.Fatalf("disposition changed action or returned the wrong class: %v", err)
+			}
+			assertReservedOperatorStep(t, c, key)
+		})
+	}
+}
+
+func TestOperatorStepReceiptAtLimit(t *testing.T) {
+	c := openTestCredentials(t)
+	ctx := testCtx(t)
+	cas, runner := operatorClaimFixture(t, c)
+	if err := c.ClaimOperatorAction(ctx, cas, runner, 1590); err != nil {
+		t.Fatal(err)
+	}
+	cas.Generation++
+	key := OperatorStepKey{cas, runner.Nonce, "checks", "fixture", 1}
+	if err := c.ReserveOperatorStep(ctx, key, 1590); err != nil {
+		t.Fatal(err)
+	}
+	body := make([]byte, MaxSurfaceApprovalArtifactBytes)
+	body[0], body[len(body)-1] = 1, 2
+	if err := c.ObserveOperatorStep(ctx, key, body); err != nil {
+		t.Fatal(err)
+	}
+	after, err := c.OperatorStep(ctx, key)
+	if err != nil || !reflect.DeepEqual(after, OperatorStepReceipt{"observed", key.Owner, 2, body}) {
+		t.Fatalf("at-limit receipt was not stored exactly: %v, size %d", err, len(after.Receipt))
+	}
+}
````
