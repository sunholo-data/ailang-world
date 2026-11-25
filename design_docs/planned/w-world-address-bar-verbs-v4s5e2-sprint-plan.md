# Sprint plan — w-world-address-bar-verbs-v4, sprint S5e2 (iteration 436, row 243)

Authority is D-WORLD-119 = A and D-WORLD-120 = A (attended 82f05004). The reviewed r1 design has quorum r2 PROCEED 3/3 with no absent seats. S1 #606, S2 #609, S2b #611, S3 #643, S4 #651, S4b #654, S5 #657, S5b #659, S5c #661, S5d #666 and S5e #669 landed. S5e merged at 2df9dc80 with judge Sonnet PASS 89 and zero blocking findings. Base is 5cac9316b755bf4b5b33263448366ec877926c8b. The planner and executor lane is codex:gpt-6.2-sol.

## 0. Binding constraints

- Land M1 as the reservation-path test milestone. Keep private file publication, independent readback and attended `service runner reconcile-files ACTION-ID` together. Appendix C records the over-cap M2 candidate.
- Design §§2.2 and 5.5 keep production effect activation gated. Preserve the existing worktree and mint resolvers, #667 episode GC and the row-271 se-tools baseline. The approved design stays at r1.
- The executor applies Appendix P with `/usr/bin/patch -p1 -E` and snapshots. The controller commits. Use migrated credentials and the repository signer, real Git and real writes only in `t.TempDir()` fixtures, scratch for logs and AILANG state. Leave TMPDIR and GOTMPDIR as inherited.

## 1. Pre-flight

### 1.1 Baseline at 5cac9316

Run §6 in archive mode on the pristine `git archive 5cac9316` export at `$S/base`. The first run before staging reached the final marker with a missing scope allowlist. After supplying the three exact M1 paths, the pristine replay reached `GATES_DONE 169s FAILS: none` in `$S/gate-base-replay/legs.txt`. Baseline build, vet, compile fence, verify_ail, plain/race store, routes and named inventories passed. The initial verifier used its default transient scratch filenames. The final gate sends shell verifier temporaries to the gate log directory via `$S/gate-bin/mktemp`, with TMPDIR and GOTMPDIR inherited. The final script also records Linux vet and a matched base CLI race selector. The final fenced base replay ended at `GATES_DONE 176s FAILS: none` in `$S/gate-base-final-script`. The Linux vet leg returned rc=0.

Whole CLI, CLI race and QUICKSTART §9 hit actual loopback denial and are **UNINFORMATIVE UNDER SANDBOX**. Real service path observations are controller-only. Archive `verify_go` is **SKIPPED** because it needs the tracked-blob Git gate. `$S/base-fails.expected` is the prior controller reference with 11 CLI and three daemon row-271 tests. The controller captures iteration-436 sets in an existing pristine product Git checkout via §7. Plain CLI, race CLI and plain `verify_go` each get an independent file. The verified prior rig result was rc=1 for the known tests. Verdicts use fail-set equality, parser controls and completion evidence.

### 1.2 Current inputs and coupling inventory

`/bin/bash "$S/anchor-head.sh" > "$S/preflight.txt"` records the full HEAD results and the exact `git grep -n -E 'func .*SYMBOLS\(' -- cmd host` queries. `$S/stage-anchors.txt` records stage lines. Each row below gives the identifying command and observed result.

| ID | Command | Observed at 5cac9316 |
|---|---|---|
| C1 | `git rev-parse HEAD`, `git status --short`, `git diff --name-only 2df9dc80 5cac9316 -- cmd host docs app \| wc -l`, `git ls-tree --name-only 5cac9316 -- cmd host docs app` | Full base SHA above, initially empty status, 0 product paths since the S5e merge, roots cmd/docs/host. Ancestry from 2df9dc80 lists record 89f45ae4 and merge #670. |
| C2 | `git show --format='%h %s' --name-only e6606ffe`, `git diff -U2 ec4565ec 5cac9316 -- cmd/ailang-worldd/episode.go host/store/context_roots_test.go` | #667 edits episode.go, episode_gc.go, episode_gc_test.go, auto_gc.go, daemon.go and context_roots_test.go. episode.go:182 has the new gc arm. Context pins at context_roots_test.go:53,87 belong to #667. |
| C3 | `grep -nE "now >= row.ExpiresAt\|row.Owner != at.Runner.Nonce\|identity ==\|count\(\\*\)\|kind = 'start'" host/store/operator_worktree.go` | Reserve expiry :28, owner and identity :31, singleton read :91. Dispose SQL action kind :137. `nl -ba host/store/operator_worktree.go` gives entry lines :13,73,108 and read state :80. |
| C4 | `grep -n -e 'if addErr != nil' -e 'exact != 1' -e 'exact++' -e '!filepath.IsAbs(raw)' -e 'registration +=' cmd/ailang-worldd/operator_worktree.go` | Add/absence disposition :101, path check :198, registration accounting :206,208, downstream count :235. `grep -n -e 'signed.Kind' -e 'IsAbs(target' -e 'current, _, currentEpisode' cmd/ailang-worldd/operator_worktree_reconcile.go` returns :39,42,67. |
| C5 | `git grep -n -E 'func .*(operatorCASRow\|ReserveOperatorStep\|observeOperatorStep\|operatorMintSession\|readOperatorMint\|RecheckClaimedOperatorAction)\(' -- cmd host` | Store anchors are operator_cas.go:32, operator_steps.go:17,55 and operator_mint.go:189,249. Authority verifier is host/authority/judge_operator.go:129. `grep -n -F "kind NOT IN ('operation','worktree','mint')" host/store/operator_mint.go` returns :266. |
| C6 | `grep -nE 'PRIMARY KEY \(kind\|identity.*NOT NULL' host/store/operator_schema.sql` and `grep -nE 'finish.*Headless\|func DecodeOperatorAction\|a.Validate\(' host/approval/operator_action.go` | :51,58 permit distinct worktree identities per action. Kind and identity together form the key. Canonical decoder :202,231 invokes validation. Finish with headless is refused at :151. Decode row/kind equality is host/store/operator_mint.go:183. |
| C7 | `git grep -n -E 'func .*(operatorStartFacts\|reconstructOperatorStart\|operatorPhysicalTarget\|canonicalLaunchPath\|writeTokenFile\|writeMCPConfig\|resolvePresetDaemon)\(' -- cmd` and `grep -nE 'sessionFile :=\|taskCopy :=\|mcp =\|afterToken' cmd/ailang-worldd/episode_preset.go` | CLI anchors are operator_start_input.go:16,68, operator_worktree.go:131, episode_launch.go:73, session_new.go:674, episode.go:556 and preset.go:298. Preset output paths :194,195,198,201,340. Claude sets MCP when headless, Codex's MCP path is empty. Task record assignment is operator_start_input.go:39–42. |
| C8 | `git grep -n -E 'reconcile-(mint\|worktree)\|TestOperatorFilesReconcile' -- cmd/ailang-worldd/{main,service,operator_service}.go docs/QUICKSTART.md design_docs/planned/w-world-address-bar-verbs-v4.md` | Top help main.go:86, service.go:37, runner help operator_service.go:24 and dispatch :275,286,289. QUICKSTART :90,103. The reviewed file row at design :234 names `TestOperatorFilesReconcile` and a repair path. |
| C9 | `git grep -n -E 'func .*(operatorClaimFixture\|operatorExec\|getOperator\|assertReservedOperatorStep\|newWorktreeStartRig\|observeProvision\|presetGit)\(' -- cmd host` | Migrated store fixture operator_cas_test.go:12,279 and operator_actions_test.go:37,46. CLI worktree fixture is operator_worktree_test.go:31, state reader is operator_provision_test.go:109 and real Git helper is preset_test.go:28. |
| C10 | `grep -nE 'WorkspaceRoot\|IsActive\|age <\|autoGCOlderThan\|checkWorktreeUnsafe\|removeWorktree' host/daemon/auto_gc.go` | :46=72 h. Scan configuration is :77–81. Directory, session, age, Git cleanliness and removal checks are :117,129,136,143,151. `isEpisodeActive` at :89–98 reads the activity snapshot. §6's live finish verifier supplies the positive caller control. |
| C11 | `go env GOOS GOARCH GOVERSION` and `grep -nF 'ubuntu-24.04' .github/workflows/ci.yml` | darwin/arm64 go1.26.6. Ubuntu 24.04 jobs at :121,446. |
| C12 | `python3 "$S/stage-diff.py" M1 M2`, `diff -u "$S/base/cmd/ailang-worldd/episode.go" "$S/st-M2/cmd/ailang-worldd/episode.go"` | Rebased unmodified probe 967 incremental, 1,172 with M1. Hardened probe 1,259 incremental, 1,464 with M1. The MCP wrapper hunk only starts near :554. The gc arm and #667 context-root pins match base. See §1.5 and Appendix C. |

Inventory names are the regex used in `$S/anchor-head.sh`. §6 exercises them by name with `-v` and checks the PASS lines. M1 edits three test files and adds no production roots, callers, packages or help. M2's measurement is recorded separately.

| Inventory | M1 | M2 probe and reason |
|---|---|---|
| `TestDaemonDependencyAllowlist` | Pass | Same daemon packages. Probe CLI uses the existing unix module. |
| `TestWorldBoundaryDependencyAllowlist` | Pass | Existing CLI/store/authority groups. |
| `TestProductionGoSurface` | Pass | Same Go package set. |
| `TestProductionContextRoots` | Pass | Reuses `cli.go:90 budgetContext` and derives child contexts. Retains the two #667 roots. If later adding a root, add its exact pin and run `go test ./host/store/ -run '^TestProductionContextRoots$'`. |
| `TestCLICredentialWriterSites` | Pass | File disposition uses the existing `operator_runner.go:162 runnerWrite` opener. |
| `TestEveryVerbHasHelp` | Pass | Probe adds nested runner help/usage example and help enumeration, with existing top-level verb. |
| `TestVerbTableMatchesTheCLISwitch` | Pass | Runner grammar extends under the existing `service` verb. |
| `TestQuickstartServiceSectionMatchesTheCLIHelp` | Pass | New command appears in main, service and runner help and QUICKSTART. Its exact example is executed by the probe. |
| `TestCredentialsMethodsRequireADeadline` | Pass | Probe adds guards and test inventory entries for `ReadOperatorFiles` and `RecordOperatorFiles`. |
| `TestStoreSurfaceTakesContext` | Pass | Both methods take a caller context with a deadline. |
| `TestBoundaryASTWriteGuard` | Pass | Probe writes at the protected CLI boundary and uses the existing store group. |
| `TestOperatorNeverDispatch` | Pass | Dormant file effect and existing typed attended CLI branch only. |
| `TestPluginCommandsMatchTheCLIHelp` in both cmd/ailang-worldd and cmd/world-publish | Pass | Existing plugin command IDs and routes. The probe verifies both. |

All 14 named probe inventory results passed with `-v`, including both plugin packages (`$S/probe-M2-inventories.log`). `/bin/bash "$S/stage-anchors.sh"` prints one definition and zero production callers of `publishOperatorStartFiles`, plus the live finish control. The only caller of `reconcileRunnerFiles` is `operator_service.go`, and its read set is shared only with dormant publication. The MCP byte formatter also supplies the existing terminal writer. The #667 files and route/app bytes compared equal; the gc dispatch arm was preserved.

### 1.3 Decisions

1. Retain the existing worktree production predicates. Use the repository signer in reservation and CLI fixtures and the existing migrated store claim fixture for read and disposition CAS. Inject disagreement through real rows, real Git and the existing reader seam. Read back the refusal class and preserved evidence, including rollback after a successful step update.
2. Test duplicate worktree rows with independently reserved identities through the step bank. Check the exact-expiry class even though a later SQL expiry predicate also prevents the effect.
3. Classify two redundant comparisons by their decoder and cardinality invariants (§4). Exercise malformed relative Git provenance separately by modifying a real captured porcelain registration.
4. Take M1 only. The hardened M2 publication and attended disposition unit is 64 cumulative lines over cap before all remaining flag and timing mutants. Carry it together.

### 1.4 Premises found false or incomplete at HEAD

| ID | Command and observed result | Resolution |
|---|---|---|
| F1 | C2 and the MCP-only hunk in C12 show same-file #667 overlap in episode.go at a different function. A byte comparison of the gc branch before/after produced rc=0. | Transplanted the MCP function region from the old stage into the archive and applied the other hunks with `patch -p1 -E -F 0`. No offset fuzz or rejected product hunks (`$S/proto-rebase.log`). |
| F2 | A pristine replay of §6 names loopback denial for CLI, race CLI and QUICKSTART §9. grep positive is `grep -F 'bind: operation not permitted' "$S/gate-base-replay/pkg-cmd-ailang-worldd.raw"`. | Controller baseline uses fresh same-rig sets. The 11 CLI and 3 daemon references remain row 271. |
| F3 | C5 and C6 show a per-kind/identity key and count-by-action query. `TestReadOperatorWorktreeGuards/duplicate` reserves both rows and reads both back. | Q-ST-READCOUNT is reachable. The SQL key permits multiple distinct identities for one action. |
| F4 | C4 shows `canonicalLaunchPath` (C7) may turn relative input into an equal physical path. `TestOperatorWorktreeReadbackGuards/relative_registration` transforms only the path field from real Git output. | Q-RB-RAWABS is a behavioural boundary against invalid Git evidence. The independent count still masks Q-RB-EXACT. |
| F5 | C8 finds the repair test in the design. The command name comes from this directive and the adjacent attended CLI patterns. | Keep the command and its executable example with M2. |
| F6 | `$S/early-numstat-fixed.txt` reports 1,172 combined. `$S/numstat.txt` reports 1,464 after authority, reopen and disposition work. | The old prototype fits on paper. The measured hardened unit exceeds the cap. |
| F7 | `grep -E '^--- FAIL: TestOperator(FileReadbackMetadata\|StartFacts)' "$S/verify-M2-harden.log"` printed both fixture failures. `grep -E '^--- PASS: TestOperator(FileReadbackMetadata\|StartFacts)' "$S/verify-M2-harden-fixed.log"` prints both passes. | Physical path comparisons normalize both sides. No product fallback for /var aliases. |
| F8 | In `$S/probe-M2-GC` with the §7 environment, `go test -v -count=1 -timeout 60s ./cmd/ailang-worldd -run '^TestDeferredFilePublicationGCProbe$'` passed both arms. `$S/verify-M2-harden-fixed.log` also passes `TestEpisodeGCFailSafeAndReclaims`, and `TestOperatorFilesReconcile/all_files` checks an identical workspace/Git snapshot with outputs under the sessions directory. | The scratch-only GC probe reported both reserved and published examples as candidates. Keep the GC gate on S5f/S6. |

### 1.5 Size and CI race cost

Controller command is exactly `git diff -M --numstat 5cac9316 HEAD -- cmd host docs app`. Sum both counts. The unfiltered product tally scans the full archive. With both plan records, the expected unfiltered M1 branch total is **1,323**. No renames or binary edits are in either candidate. Read-only equivalent is `python3 "$S/stage-diff.py" M1 M2` with output in `$S/numstat.txt`. Save any deferred patch separately from the executor fences.

| Milestone | Added | Deleted | Cumulative | Increment | Unfiltered product | STATUS/DELETED/MANIFEST |
|---|---:|---:|---:|---:|---:|---:|
| M1 | 164 | 41 | **205** | 205 | 205 | 3/0/3 |
| M2 hardened probe | 1,403 | 61 | **1,464** | 1,259 | 1,464 | Deferred |

The final archive replay ended at `GATES_DONE 230s FAILS: none` in `$S/gate-scratch-final-script`. Same four-family CLI race selector took 16.682 s at base and 20.171 s at M1, a measured +3.489 s package-time difference. Whole CLI race timing remains uninformative under sandbox. Route census was 54 and three surface POSTs at base and M1, with equal served bytes and no production file callers.

### 1.6 Portability

Darwin uses physical paths for comparisons, including /private/var, and the real Git wrapper in M1 uses POSIX shell and a captured NUL-delimited listing. The file probe uses directory-relative Unix open/stat with portable O_EXCL and no-follow flags. `TestOperatorFilesReconcile/all_files` exercises restrictive umask. Use refusal plus byte and metadata snapshots for OS-dependent errno. `GOOS=linux go vet ./cmd/ailang-worldd ./host/store` passed on the probe. Native Ubuntu tests and the per-flag behavioural matrix remain gates for the deferred patch.

## 2. How the executor applies a sub-milestone

````bash
export GIT_OPTIONAL_LOCKS=0
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter436-row243-v4s5e2
T=$HOME/.ailang/state/world-iter436/executor
S=$HOME/.ailang/state/world-iter436/plan-scratch
B=5cac9316b755bf4b5b33263448366ec877926c8b
PLAN=$W/design_docs/planned/w-world-address-bar-verbs-v4s5e2-sprint-plan.md
mkdir -p "$T"
unset SNAP_PARENT SNAP_STATUS_LIST
git -C "$W" merge-base --is-ancestor "$B" HEAD || exit 2
git -C "$W" diff --quiet "$B" -- cmd host docs app || exit 2
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

Apply M1 once to the base product tree. The snapshot excludes both plan records. Appendix C is scratch evidence for a later iteration.

## 3. Milestones

### M1 — Reservation-path pre-caller pins

| Boundary and input | Named test and state read back |
|---|---|
| Wrong current runner, empty identity and exact expiry | `TestReserveOperatorWorktreeGuards/{owner,identity,expiry}` uses a real proof and fresh heartbeat. It reads action, operation row count, empty sessions/start and zero reservation, plus the valid control. Exact expiry returns `ErrOperatorExpired`. |
| Resume unresolved uncertainty and reject multiple worktree rows | `TestReadOperatorWorktreeGuards/{uncertain,duplicate}` reads the exact uncertain key/version and the two distinct reservations, with an unchanged action and `ErrOperatorMissing` in the duplicate case. |
| SQL action-kind CAS after step update | `TestDisposeOperatorUncertainStepGuards/action_kind` reads exact action and reserved step after refusal. |
| Real successful add followed by absent readback | `TestOperatorWorktreeWindow/success_then_absent` proves a real created checkout, removes the fixture worktree and ref before readback, and checks `uncertain`, reserved v2 and the exact remedy, then executes the remedy and reads `stopped: absent` v3 after typed fixture approval. |
| Stored signed kind, absolute repository and common dir | `TestOperatorWorktreeRecoveryGuards/{signed_kind,repo_abs,common_abs}` reads stable persisted fixtures through the CLI, requires refusal before prompt and snapshots store, branch and files. |
| Locked read returns an error despite an unchanged key | `TestOperatorWorktreeLockedReread/error` supplies the real key plus a distinct injected error on call two. No prompt, no state change. The generation arm remains independently exercised. |
| Relative or duplicate Git registration | `TestOperatorWorktreeReadbackGuards/{relative_registration,registration_count}` uses a real Git fixture, controlled registration text and exact branch/file readback. |

### 3.1 Design §5.3 row assignment

| Reviewed row | Assignment |
|---|---|
| `TestOperatorStartOnce` | M1 reinforces worktree evidence. File prerequisites in Appendix C. One live child, confinement and receipt-before-spawn in S6. |
| `TestOperatorStartChanged` | Landed S5c/S5d, exercised by §6. File binding and expiry arms in the M2 probe. |
| `TestOperatorStartCollisions` | Existing worktree/file preflight tests plus M1's readback. The file-install no-replace guarantee stays with Appendix C. |
| `TestOperatorWorktreeWindow` | M1's success-then-absent arm pins uncertainty and exact attended remedy. Existing S5d injection arms run alongside it. |
| `TestOperatorMintBinding`, `TestOperatorMintRollback`, `TestOperatorMintReceiptWindow` | Landed mint proof, rollback and recovery. §6 runs focused regression gates. Live mint consumer awaits process pairing. |
| `TestOperatorFilesReconcile` | M2 in Appendix C, one reserve/install/readback/disposition unit. |
| `TestOperatorLaunchSession`, `TestOperatorSpawnRefused`, `TestOperatorLaunchCrashWindow`, `TestOperatorLaunchExit` | S5f/S6 paired process and S6 launch. |
| `TestOperatorStartBrowser` | S8, fake Sol launcher with reconnect/restart. |
| `TestOperatorHerdrUnavailable`, `TestOperatorHerdrHandoff`, `TestOperatorHerdrWindow` | S12 E4, typed ticket/TTY and socket/readiness gates. |

### 3.2 Design §5.5 and §2.2 activation

| Reviewed gate | Assignment |
|---|---|
| `TestOperatorServiceFence` | §6 exercises typed approval and edited plist. M1 checks pre-prompt refusal on existing worktree resolution. |
| `TestOperatorServiceOwnership` | §6 UID, modes, links and executable identity. |
| `TestOperatorServiceConfig` | §6 typed config and raw-argv refusal. |
| `TestOperatorServicePair` | §6 existing compatibility. Paired restart and admission must precede start activation in S5f/S6. |
| `TestOperatorServiceStartup` | §6 daemon, boot failure and bounded startup. Maintain no queued work on unavailable or skewed peers. |
| `TestOperatorServiceUpgrade`, `TestOperatorServiceUpgradePairFence` | §6 rollback and build/store compatibility. Active-child drain and recovery in S5f/S6. |
| `TestOperatorServiceWriteScope` | §6 fixture lifecycle, real-profile observations outside sandbox. |
| `TestOperatorServiceEnvironment` | §6 tests approved environment and rejects browser-supplied authority. |
| `TestOperatorFinishReadiness` | §6 store readiness and E2-equivalent checks before any real finish claim. |
| `TestOperatorSchema10` | §6 and commit-time fixture cover fresh/v9, legacy reader and rollback arms. |

Keep installation attended and immutable, config/spool private and user-owned, executables and stores pinned, PATH/HOME and tool auth approved, and protocol/build/schema compatible. Activation requires recent heartbeat/catalogue, shared runner/episode exclusion, exact previous worktree/mint/file receipts and current scoped session, verified ownership and quiescence on restart, bounded old-build/terminal drain and compatible rollback. Before the first production caller, protect non-terminal operator evidence from #667 GC. Service/doctor and E2, E3, E4/E5 keep the design's published sequencing.

## 4. Mutation table

Anchors are at HEAD and remain identical at M1. M1 edits tests only, so its diff has zero new production branches. `$S/mutants.json` starts from the judge's exact replacements. Named selectors distinguish `--- FAIL: <test>` from compile errors, binds and unrelated package failure. Every claimed kill drives the real store, Git or CLI and uses no static inventory as its proof. `$S/mutation-ledger.jsonl` records the actual command, log and verdict per phase.

| ID | Production anchor | Mutation | Killing test package / test | Base named | M1 named | Whole |
|---|---|---|---|---|---|---|
| Q-ST-OWNER | host/store/operator_worktree.go:31 | drop current runner owner check | host/store / TestReserveOperatorWorktreeGuards/owner | SURVIVED | KILLED | KILLED |
| Q-ST-IDENT | host/store/operator_worktree.go:31 | accept empty identity | host/store / TestReserveOperatorWorktreeGuards/identity | SURVIVED | KILLED | KILLED |
| Q-ST-EXPIRYEDGE | host/store/operator_worktree.go:28 | now >= expires → now > expires | host/store / TestReserveOperatorWorktreeGuards/expiry | SURVIVED | KILLED | KILLED |
| Q-ST-READSTATE | host/store/operator_worktree.go:80 | allow running only | host/store / TestReadOperatorWorktreeGuards/uncertain | SURVIVED | KILLED | KILLED |
| Q-ST-READCOUNT | host/store/operator_worktree.go:91 | count = 1 → count >= 1 | host/store / TestReadOperatorWorktreeGuards/duplicate | SURVIVED | KILLED | KILLED |
| Q-ST-DISPKINDSQL | host/store/operator_worktree.go:137 | drop start kind from disposition UPDATE | host/store / TestDisposeOperatorUncertainStepGuards/action_kind | SURVIVED | KILLED | KILLED |
| Q-RB-EXACT | cmd/ailang-worldd/operator_worktree.go:235 | exact != 1 → exact < 1 | cmd/ailang-worldd / TestOperatorWorktreeReadbackGuards/registration_count | SURVIVED | SURVIVED (equivalent) | UNINFORMATIVE UNDER SANDBOX |
| Q-RB-RAWABS | cmd/ailang-worldd/operator_worktree.go:198 | allow relative porcelain paths | cmd/ailang-worldd / TestOperatorWorktreeReadbackGuards/relative_registration | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| Q-W-ABSENTDISP | cmd/ailang-worldd/operator_worktree.go:101 | drop prior add error requirement | cmd/ailang-worldd / TestOperatorWorktreeWindow/success_then_absent | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| Q-R-SIGNEDKIND | cmd/ailang-worldd/operator_worktree_reconcile.go:39 | drop signed start kind check | cmd/ailang-worldd / TestOperatorWorktreeRecoveryGuards/signed_kind | SURVIVED | SURVIVED (equivalent) | UNINFORMATIVE UNDER SANDBOX |
| Q-R-REPOABS | cmd/ailang-worldd/operator_worktree_reconcile.go:42 | drop repo absolute check | cmd/ailang-worldd / TestOperatorWorktreeRecoveryGuards/repo_abs | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| Q-R-COMMONABS | cmd/ailang-worldd/operator_worktree_reconcile.go:42 | drop common absolute check | cmd/ailang-worldd / TestOperatorWorktreeRecoveryGuards/common_abs | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| Q-R-REREADERR | cmd/ailang-worldd/operator_worktree_reconcile.go:67 | discard locked reread error | cmd/ailang-worldd / TestOperatorWorktreeLockedReread/error | SURVIVED | KILLED | UNINFORMATIVE UNDER SANDBOX |
| CTL-DIE | cmd/ailang-worldd/operator_worktree.go:61 | bypass initial-state refusal | cmd/ailang-worldd / TestOperatorStartCollisions/registration | — | KILLED | — |
| CTL-LIVE | cmd/ailang-worldd/operator_worktree.go:199 | change diagnostic outside control selector | cmd/ailang-worldd / TestOperatorStartOnce | — | SURVIVED | — |
| CTL-LIVE-STORE | host/store/operator_worktree.go:31 | insert redundant false | host/store / TestReserveOperatorWorktreeGuards/valid | — | SURVIVED | SURVIVED |
| CTL-DIE-STORE | host/store/operator_worktree.go:31 | invert claimed state check | host/store / TestReserveOperatorWorktreeGuards/valid | — | KILLED | KILLED |

The mask for **Q-R-SIGNEDKIND** is `host/approval/operator_action.go:151,231` and `host/store/operator_mint.go:183`. The only other valid kind is finish, the decoder enforces row-kind equality and rejects finish plus headless, and the immediate headless guard still refuses the valid finish fixture. The mask for **Q-RB-EXACT** is the separate registration count at `cmd/ailang-worldd/operator_worktree.go:206,208,235`. Each exact increment follows a delimited append. When exact ≥2, the delimiter cardinality already refuses. The nonmatching and duplicate registrations remain refused. §4.1 records which whole-suite checks require the controller.

### 4.1 Whole-package runs for the controller

The archive's whole CLI encounters loopback denial. Recheck the seven rows below outside sandbox in the named clean Git checkout using `$S/whole-mutants.sh` and `$S/whole-mutant-ids.txt`. The default checkout is the sprint worktree after M1 commit. The script stages backup copies under scratch, restores file bytes on exit and compares against `$S/pkg-cmd-ailang-worldd.base-fails` recorded by §7. A named failing row must add its top-level test while preserving the collected base set. A measured equivalence must preserve the fail set and show the named arm passes. Iteration-435's raw rc=1 reports for Q-ST-OWNER and Q-W-ABSENTDISP matched the CLI base after set comparison, as recorded in the judge report.

````text
Q-RB-EXACT
Q-RB-RAWABS
Q-W-ABSENTDISP
Q-R-SIGNEDKIND
Q-R-REPOABS
Q-R-COMMONABS
Q-R-REREADERR
````

## 5. Files

````text
M  cmd/ailang-worldd/operator_worktree_guards_test.go
M  cmd/ailang-worldd/operator_worktree_test.go
M  host/store/operator_worktree_guards_test.go
````

All three are existing test files. The measured candidate paths and counts are in `$S/changed-M2` and `$S/numstat.txt`. Scope and caller legs in §6 check the exact shipped footprint.

## 6. Sprint gates

````bash gates.sh
#!/bin/bash
# Usage: gates.sh <tree> <logdir> <base|final|scratch> [sha]
set -u
set -o pipefail
export GIT_OPTIONAL_LOCKS=0
T=$1; L=$2; MODE=$3
B=5cac9316b755bf4b5b33263448366ec877926c8b
REF=${GATES_GIT_REF:-$T}
S="$HOME/.ailang/state/world-iter436/plan-scratch"
SCRIPTS=$(cd "$(dirname "$0")" && pwd)
SHA=${4:-$(git -C "$REF" rev-parse HEAD)}
mkdir -p "$L/tmp"
L=$(cd "$L" && pwd -P)
cd "$T" || exit 2
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export AILANG_STATE_DIR="$L/ailang-state" AILANG_CACHE_DIR="$L/ailang-cache"
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR WORLD_UPDATE_OPERATOR_ACTION_GOLDEN WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN
unset GOFLAGS
mkdir -p "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR"
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
    final) git -C "$REF" diff --name-only "$B" HEAD -- cmd host docs app ;;
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
leg linux-vet 'GOOS=linux python3 "$S/bounded.py" 180 go vet ./cmd/ailang-worldd ./host/store'
leg vet 'python3 "$S/bounded.py" 180 go vet ./...'
leg fence "go test -run '^\$' ./..."
leg verify_ail 'PATH="$S/gate-bin:$PATH" WORLD_GATE_TMP="$L/tmp" python3 "$S/bounded.py" 300 ./scripts/verify_ail.sh'
PACKAGES=$TOUCHED
if [ "$MODE" = base ]; then PACKAGES=$EXPECTED; fi
test_set() {
  local tag=$1 test_rc=0
  shift
  python3 "$S/bounded.py" 510 go test -count=1 -timeout 8m "$@" > "$L/$tag.raw" 2>&1 || test_rc=$?
  python3 "$S/failset.py" "$L/$tag.raw" "$test_rc" > "$L/$tag.fails" || {
    grep -m 1 -F 'bind: operation not permitted' "$L/$tag.raw" && return 77
    return 1
  }
  if grep -Fq 'bind: operation not permitted' "$L/$tag.raw"; then
    grep -m 1 -F 'bind: operation not permitted' "$L/$tag.raw"
    return 77
  fi
  cat "$L/$tag.fails"
  if [ "$test_rc" -gt 1 ] || grep -Eq '\[build failed\]|test timed out|BOUND EXCEEDED' "$L/$tag.raw"; then return 1; fi
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
  line 'LEG s5e2-accept SKIPPED 0s (base has no S5e2 tests)'
else
  leg s5e2-accept '/bin/bash "$SCRIPTS/accept.sh" "$T" M1 fast'

fi
leg added-cli-race 'python3 "$S/bounded.py" 150 go test -race -v -count=1 -timeout 2m ./cmd/ailang-worldd -run "^(TestOperatorWorktreeWindow|TestOperatorWorktreeRecoveryGuards|TestOperatorWorktreeReadbackGuards|TestOperatorWorktreeLockedReread)$"'
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
  PATH="$S/gate-bin:$PATH" WORLD_GATE_TMP="$L/tmp" VERIFY_GO_LEGS=plain python3 "$S/bounded.py" 1500 /bin/bash ./scripts/verify_go.sh > "$L/verify_go.raw" 2>&1 || rc=$?
  python3 "$S/failset.py" "$L/verify_go.raw" "$rc" > "$L/verify_go.fails" || {
    grep -m 1 -F 'bind: operation not permitted' "$L/verify_go.raw" && return 77
    return 1
  }
  if grep -Fq 'bind: operation not permitted' "$L/verify_go.raw"; then
    grep -m 1 -F 'bind: operation not permitted' "$L/verify_go.raw"
    return 77
  fi
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

Use §2's extracted scripts. The named acceptance uses the pinned AILANG binary, ssh-keygen and fixture node/tool paths appearing in the exact command below. A bind denial blocks the snapshot as uninformative. Full CLI and daemon tests stay in controller gates.

````bash accept.sh
#!/bin/bash
set -euo pipefail
W=$1; M=$2; PART=$3
test "$M" = M1
test "$PART" = fast
S="$HOME/.ailang/state/world-iter436/plan-scratch"
export AILANG_BIN="$HOME/.pinned-ailang/ailang"
export WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
run_dir=$(mktemp -d "$S/accept.XXXXXX")
export AILANG_STATE_DIR="$run_dir/state" AILANG_CACHE_DIR="$run_dir/cache"
mkdir -p "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR"
cd "$W"
run() {
  local pkg=$1 names=$2 name log rc=0
  log="$run_dir/${pkg//\//-}.txt"
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
run host/store 'TestReserveOperatorWorktreeGuards|TestDisposeOperatorUncertainStepGuards|TestReadOperatorWorktreeGuards|TestOperatorStepReceiptAtLimit|TestOperatorStepRestart|TestOperatorStepGuards|TestCredentialsMethodsRequireADeadline|TestOperatorSchema10'
run cmd/ailang-worldd 'TestOperatorWorktreeSurvivors|TestOperatorWorktreeRecoveryGuards|TestOperatorWorktreeReadbackGuards|TestOperatorReadbackErrorIsNeverCreated|TestOperatorWorktreeLockedReread|TestOperatorWorktreeWindow|TestOperatorStartOnce|TestOperatorWorktreeDispositionCAS|TestEveryVerbHasHelp|TestQuickstartServiceSectionMatchesTheCLIHelp'
printf 'ACCEPT_OK %s %s\n' "$M" "$PART"
````

The snapshot path is `$HOME/.ailang/state/world-iter436/snap/M1`. Read-only Git lists the cumulative paths against base, with absent product roots allowed. Script defaults use the real executor path. Overrides are for scratch rehearsals.

````bash snapshot.sh
#!/bin/bash
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=$1
M=$2
B=5cac9316b755bf4b5b33263448366ec877926c8b
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter436/snap}/$M
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

Record baseline sets outside sandbox before product application. Run the committed checkout's broad gates and mutation set after M1. The base verification script rejects product drift and accepts a committed plan-only addition.

````bash
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter436-row243-v4s5e2
STATE=$HOME/.ailang/state/world-iter436
T=$STATE/executor
S=$STATE/plan-scratch
BASE_TREE=${BASE_TREE:-$W}
/bin/bash "$S/baseline-controller.sh" "$BASE_TREE" || exit
````

After the executor snapshots and the plan commit exists:

````bash
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter436-row243-v4s5e2
STATE=$HOME/.ailang/state/world-iter436
T=$STATE/executor
S=$STATE/plan-scratch
printf '%s\n' 'test(operator): pin reservation survivors before file publication' > "$T/M1.msg"
/bin/bash "$STATE/commit_ms.sh" M1 "$T/M1.msg" 'AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang /bin/bash "'"$T"'/accept.sh" "'"$W"'" M1 fast' || exit
/bin/bash "$T/gates.sh" "$W" "$S/gate-controller-final" final || exit
/bin/bash "$S/whole-mutants.sh" "$W" "$S/pkg-cmd-ailang-worldd.base-fails"
````

`verify_go.sh` runs with GOFLAGS unset and parses text. The fail-set parser prints its known positive on each call, rejects rc ≠0 with no attributable failures, and requires package summaries. If a prior load-sensitive test reappears, preserve both logs and retest under the same rig conditions. A fail-set change remains a red leg until accounted for. The controller also runs Ubuntu CI.

M1 round trip in `$S/rt-2` extracted Appendix P with the controller's exact awk and applied it with `/usr/bin/patch -p1 -E` to a fresh export of 5cac9316. Full archive, cmd, host and docs `diff -r` each returned 0. App was absent on both sides. The two inherited .orig fixtures explain the initial all-file backup census. The relative before/after backup list matched with zero new files. The three exact scripts matched their scratch copies and passed `bash -n`. The exact extracted snapshot script read base Git via GIT_DIR and GIT_WORK_TREE, staged no Git changes, and produced 3/0/3 and three verified hashes. `$S/roundtrip.log` and `$S/rt-2/accept.log` record `ACCEPT_OK M1 fast`.

## 8. Risks

- The exact expiry class has two layers of protection. The fixture keeps the heartbeat fresh so the mutated comparison reaches the later SQL gate and exposes the wrong refusal.
- The base archive reproduces the sandbox limit. Only the external same-rig run can attribute row-271 CLI/daemon results or detect an independent new e2e failure.
- The M2 prototype executes real writes in test directories only. Its incomplete portability and per-flag evidence require review in the later file milestone before enabling a caller.

## 9. Acceptance criteria

1. M1 applies byte-for-byte from Appendix P against the fixed base. Expected filtered numstat 205, snapshot counts 3/0/3, each hash verifies.
2. All 11 nonequivalent judge-435 survivors fail their named behavioural test on M1. The two measured equivalent cases retain their masking invariants and unchanged state checks. Both must-die and both must-live controls fire.
3. The added CLI race seconds and whole store runs are recorded. Gate legs show no new informative failures or caller, route or served-byte differences against base.
4. Controller runs the package fail-set gates and Ubuntu CI. All file publication, readback and terminal disposition ship together after the Appendix C requirements.

## Appendix C — after this sprint (measured)

- **File publication and disposition as one milestone.** On this HEAD the pure rebase is **967 incremental / 1,172 cumulative**. The hardened scratch unit is **1,259 incremental / 1,464 cumulative**, 64 over cap, with the same unfiltered product count. It keeps S5e’s three-file Claude bundle and two-file headless Codex case. Use `$S/st-M2`, `$S/M2-deferred.patch`, `$S/numstat.txt`, `$S/verify-M2-probe-final.log`. Scratch replay after M1 was identical to `st-M2` by full `diff -r` (`$S/rt-deferred-M2/diff-M2.txt`). Go build, native CLI/store vet, Linux CLI/store vet and the targeted Darwin tests completed. Appendix P supplies the executor’s M1 patch.
- **File prerequisite and effect trace.** Under the held canonical episode lock, require the observed-created worktree reservation with exact Git evidence and a same-owner observed mint. Recheck current key, action signature and exact signed facts, live action and credential expiry, and transactional equality against the signed mint session, approved address and pinned preset. Bind signed task, token hash and derived MCP bytes to one canonical directory/inode and reservation before writes. Open privately with exclusive no-replace, independently reopen and validate path, regular type, owner, exact mode, single link and digest, compare the installed and re-read inodes and compare fresh evidence again before observing. Exact or wholly absent evidence can be stopped by attended review. Partial, inaccessible or foreign evidence stays reserved and uncertain. Keep original worktree, mint/start and file evidence. Resolution compares independent action owner and step CAS in one transaction, after the step update too.
- **File prototype tests exercised.** `TestOperatorFilesReconcile` covers all boundaries before reservation, after reservation before first write, both file gaps, after writes before readback, after readback before observation, same-byte replacement, concurrent file appearance, changed content/parent during review, expired reopen, after-step owner drift, stale step version, later step and observation refusal. It runs the exact QUICKSTART command `ailang-worldd service runner reconcile-files oa-0123456789abcdef0123456789abcdef` and verifies preserved evidence after no terminal, typed refusal and attended success. `TestOperatorFileReadbackMetadata` checks digest, UID, mode, symlink, hardlink, FIFO and bounded bytes. `TestOperatorFilePublicationAdmission` checks missing/wrong worktree and mint evidence, exact expiry, expiry between reservation and first write, revoked key/credential, wrong grants/TTL/token, existing file, linked or public directory, moved tree and held lock.
- **File gates before the next runnable patch.** The hardened draft makes files visible only in a 0700 directory, uses `openat` and re-opens it by path for readback, and reuses the existing read prompt/lock/store boundary. The remaining whole matrix is **UNMEASURED** on this scratch candidate. Branch candidates with concrete file and line anchors are in `$S/branch-anchors-M2.txt`. The five spot mutants in `$S/m2-spot-ledger.jsonl` were killed by `TestOperatorFilesReconcile` in cmd/ailang-worldd. `P2-EXCL` at `cmd/ailang-worldd/operator_files.go:187` by `competing_file`, `P2-DIGEST` at `cmd/ailang-worldd/operator_files.go:146` by `changed_file`, `P2-PARENT` at `cmd/ailang-worldd/operator_files_reconcile.go:79` by `parent_changed_during_review`, `P2-NEEDS-RESERVED` at `cmd/ailang-worldd/operator_files_reconcile.go:33` by `before_reservation`, `P2-ACTION-OWNER` at `host/store/operator_files.go:152` by `owner_changed`. Reproduce with `python3 "$S/m2-spot.py"`. Test removal of each bit of `O_NOFOLLOW`, `O_EXCL`, `O_NONBLOCK`, `O_DIRECTORY`, `O_CLOEXEC`, `AT_SYMLINK_NOFOLLOW` and `LOCK_EX\|LOCK_NB` with real or deterministic bounded interposition. State and demonstrate any mask such as exclusive create independently refusing symlinks. Pin each input comparison and SQL predicate separately, including path/UID/mode/type/link/digest mismatches, independent CAS at reserve and successful observe with drift after step update, revoked/signature/mint drift at the final authority boundary, absent/observed/wrong-phase rows before prompt, file and parent reopen races, and context/read limits. Give a named runtime killing test for every reviewed row and include the same-phase report in the later patch. Use `go test -v -count=1 -timeout 2m ./cmd/ailang-worldd -run '^(TestOperatorFilesReconcile|TestOperatorFileReadbackMetadata|TestOperatorFilePublicationAdmission)$'` with §7's pinned environment on Ubuntu 24.04 as well as Darwin. Retain independent publication readback and attended stop in a single safe patch.
- **S5f/S6 paired-process consumer, owner recovery and drain.** Carry `E-RES-EPISODE-CONFLICT`, `E-RES-REREAD`, `J-OWN-UID`, `Z-ABSENT-ANY`, `Z-WAIT-READBUDGET`, `TestOperatorFirstManagedEpisodeProcess`, `TestOperatorOldBuildTerminalDrain`, `TestOperatorCurrentOwnerProcessRecovery` and `TestOperatorServiceUpgrade/active_child`. Prior FD harness measured 72 lines by S5d, full implementation size is UNMEASURED. Require appropriate pair, UID, locks, current session, pinned exec, confinement and no terminal revoke race.
- **#667 GC gate before a first production caller.** C10 and F8 show that #667 consults the live-session snapshot and ordinary Git evidence and that the probe puts the file set under sessions. `$S/probe-M2-GC.log` records `TestDeferredFilePublicationGCProbe/{reserved,published}`. Both had empty real Git status and a `clean & merged` GC candidate at an injected 73 h age with no live activity. Readback confirmed preserved evidence after the dry run. Add an operator-step-aware skip before activation, with real Git time-injected coverage in an approved package. The GC change’s size is UNMEASURED. This sprint's scope gate keeps the landed #667 files byte-identical.
- **Launch, browser, finish and Herdr.** Keep S6 launch/session/PID and drain work, S8 browser lifecycle and S12 interactive Herdr work in their approved slices. Reviewed raw budgets are S6 890, S8 860 and S12 700. Full implementations are UNMEASURED. Preserve E3, E5/#604, E4 and D-WORLD-120/E2 readiness. Row 271 stays with its rig baseline.

## Appendix P — compiled patches

### Patch M1

````diff M1
diff --git a/cmd/ailang-worldd/operator_worktree_guards_test.go b/cmd/ailang-worldd/operator_worktree_guards_test.go
index 60673eca..4ef6a1cd 100644
--- a/cmd/ailang-worldd/operator_worktree_guards_test.go
+++ b/cmd/ailang-worldd/operator_worktree_guards_test.go
@@ -13,6 +13,7 @@ import (
 
 	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/worktree"
 )
 
 func TestOperatorWorktreeSurvivors(t *testing.T) {
@@ -46,7 +47,7 @@ func TestOperatorWorktreeSurvivors(t *testing.T) {
 }
 
 func TestOperatorWorktreeRecoveryGuards(t *testing.T) {
-	for _, name := range []string{"store", "branch", "path", "action", "basename", "headless", "observed"} {
+	for _, name := range []string{"store", "branch", "path", "action", "basename", "headless", "observed", "signed_kind", "repo_abs", "common_abs"} {
 		t.Run(name, func(t *testing.T) {
 			r := newWorktreeStartRig(t, false)
 			crash := errors.New("before add")
@@ -83,6 +84,16 @@ func TestOperatorWorktreeRecoveryGuards(t *testing.T) {
 					target.Path += "-other"
 				case "action":
 					target.Action += "-other"
+				case "signed_kind":
+					a.Kind, a.Task, a.Options.Headless = "finish", "", false
+					a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
+					if err := a.Validate(); err != nil {
+						t.Fatal(err)
+					}
+				case "repo_abs":
+					target.Repo = strings.TrimPrefix(target.Repo, string(filepath.Separator))
+				case "common_abs":
+					target.Common = strings.TrimPrefix(target.Common, string(filepath.Separator))
 				case "headless":
 					a.Options.Headless = false
 					a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
@@ -105,7 +116,7 @@ func TestOperatorWorktreeRecoveryGuards(t *testing.T) {
 					t.Fatal(err)
 				}
 				// Both reads see the same damaged rows. The CAS can agree with them.
-				r.exec(t, "UPDATE operator_actions SET action=?, digest=?", body, digest[:])
+				r.exec(t, "UPDATE operator_actions SET action=?, digest=?, kind=?", body, digest[:], a.Kind)
 				r.exec(t, "UPDATE operator_action_steps SET identity=? WHERE kind='worktree'", string(encoded))
 			}
 			before := observeProvision(t, r.f.r, r.creds)
@@ -126,7 +137,7 @@ func TestOperatorWorktreeRecoveryGuards(t *testing.T) {
 }
 
 func TestOperatorWorktreeReadbackGuards(t *testing.T) {
-	for _, name := range []string{"repo_common", "checkout_common", "link_regular", "link_prefix", "registration_count", "directory", "physical"} {
+	for _, name := range []string{"repo_common", "checkout_common", "link_regular", "link_prefix", "registration_count", "relative_registration", "directory", "physical"} {
 		t.Run(name, func(t *testing.T) {
 			f := newEpisodeRun(t)
 			ctx := presetTestContext(t)
@@ -201,6 +212,33 @@ func TestOperatorWorktreeReadbackGuards(t *testing.T) {
 					seedStartFixtureFile(t, link, []byte("interrupted gitfile\n"))
 					wrapper("if [ \"$2\" = " + shellQuote(target.Path) + " ]; then cp " + shellQuote(saved) + " " + shellQuote(link) + " || exit 1; fi")
 				}
+			case "relative_registration":
+				cwd, err := os.Getwd()
+				if err != nil {
+					t.Fatal(err)
+				}
+				cwd, err = filepath.EvalSymlinks(cwd)
+				if err != nil {
+					t.Fatal(err)
+				}
+				relative, err := filepath.Rel(cwd, target.Path)
+				if err != nil {
+					t.Fatal(err)
+				}
+				list, err := worktree.Output(ctx, git, []string{"-C", repo, "worktree", "list", "--porcelain", "-z"})
+				if err != nil || !strings.Contains(list, "worktree "+target.Path+"\x00") {
+					t.Fatalf("cannot read real registration: %v", err)
+				}
+				list = strings.Replace(list, "worktree "+target.Path+"\x00", "worktree "+relative+"\x00", 1)
+				capture := filepath.Join(root, "porcelain")
+				seedStartFixtureFile(t, capture, []byte(list))
+				wrapper(`if [ "$3" = worktree ]; then exec cat ` + shellQuote(capture) + "; fi")
+				before := launchFileSnapshot(t, target.Path)
+				e, err := target.inspect(ctx, git)
+				if err == nil || err.Error() != "invalid Git worktree registration" || e.state == treeExact ||
+					!reflect.DeepEqual(before, launchFileSnapshot(t, target.Path)) {
+					t.Fatalf("relative registration admitted or checkout changed: %+v %v", e, err)
+				}
 			case "registration_count":
 				second := filepath.Join(root, "second")
 				presetGit(t, f.f, "worktree", "add", "--detach", second, target.SHA)
@@ -230,9 +268,13 @@ func TestOperatorWorktreeReadbackGuards(t *testing.T) {
 				wrapper("if [ \"$3\" = for-each-ref ]; then mv " + shellQuote(victim) + " " + shellQuote(moved) + " && " + replacement + " || exit 1; fi")
 			}
 			e, err := target.inspect(ctx, git)
-			if name == "physical" {
-				if err == nil || err.Error() != "worktree physical path differs" {
-					t.Fatalf("redirected ancestor not refused: %v", err)
+			if name == "physical" || name == "relative_registration" {
+				want := "worktree physical path differs"
+				if name == "relative_registration" {
+					want = "invalid Git worktree registration"
+				}
+				if err == nil || err.Error() != want {
+					t.Fatalf("registration or path not refused: %v", err)
 				}
 			} else if err != nil || e.state != treePartial {
 				t.Fatalf("attribution accepted or wrong refusal: state=%d %v", e.state, err)
@@ -281,35 +323,46 @@ func TestOperatorReadbackErrorIsNeverCreated(t *testing.T) {
 }
 
 func TestOperatorWorktreeLockedReread(t *testing.T) {
-	r := newWorktreeStartRig(t, false)
-	crash := errors.New("before add")
-	if err := r.start(t, func(step string) error {
-		if step == "before add" {
-			return crash
-		}
-		return nil
-	}); !errors.Is(err, crash) {
-		t.Fatalf("expected reserved worktree: %v", err)
-	}
-	var before provisionSnapshot
-	calls := 0
-	read := func(path, id string) (store.OperatorStepKey, operatorWorktreeTarget, string, error) {
-		key, target, episode, err := readRunnerWorktree(path, id)
-		calls++
-		if calls == 1 && err == nil {
-			r.exec(t, "UPDATE operator_actions SET generation=generation+1")
-			before = observeProvision(t, r.f.r, r.creds)
-		}
-		return key, target, episode, err
-	}
-	r.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
-		t.Error("a changed reservation reached the prompt")
-		return nil, errors.New("unexpected prompt")
-	}
-	var stderr strings.Builder
-	code := r.env.reconcileRunnerWorktree(r.action.ID, io.Discard, &stderr, read)
-	if code != exitUsage || calls != 2 || !strings.Contains(stderr.String(), store.ErrOperatorConflict.Error()) ||
-		!reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
-		t.Fatalf("reread did not preserve the replacement row and its evidence: %d %s", code, &stderr)
+	for _, name := range []string{"generation", "error"} {
+		t.Run(name, func(t *testing.T) {
+			r := newWorktreeStartRig(t, false)
+			crash := errors.New("before add")
+			if err := r.start(t, func(step string) error {
+				if step == "before add" {
+					return crash
+				}
+				return nil
+			}); !errors.Is(err, crash) {
+				t.Fatalf("expected reserved worktree: %v", err)
+			}
+			before := observeProvision(t, r.f.r, r.creds)
+			want := store.ErrOperatorConflict
+			if name == "error" {
+				want = errors.New("fixture locked read error")
+			}
+			calls := 0
+			read := func(path, id string) (store.OperatorStepKey, operatorWorktreeTarget, string, error) {
+				key, target, episode, err := readRunnerWorktree(path, id)
+				calls++
+				if name == "error" && calls == 2 && err == nil {
+					return key, target, episode, want
+				}
+				if name == "generation" && calls == 1 && err == nil {
+					r.exec(t, "UPDATE operator_actions SET generation=generation+1")
+					before = observeProvision(t, r.f.r, r.creds)
+				}
+				return key, target, episode, err
+			}
+			r.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
+				t.Error("a changed reservation reached the prompt")
+				return nil, errors.New("unexpected prompt")
+			}
+			var stderr strings.Builder
+			code := r.env.reconcileRunnerWorktree(r.action.ID, io.Discard, &stderr, read)
+			if code != exitUsage || calls != 2 || !strings.Contains(stderr.String(), want.Error()) ||
+				!reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
+				t.Fatalf("reread did not preserve the replacement row and its evidence: %d %s", code, &stderr)
+			}
+		})
 	}
 }
diff --git a/cmd/ailang-worldd/operator_worktree_test.go b/cmd/ailang-worldd/operator_worktree_test.go
index 9557aef8..7513c2ce 100644
--- a/cmd/ailang-worldd/operator_worktree_test.go
+++ b/cmd/ailang-worldd/operator_worktree_test.go
@@ -327,7 +327,7 @@ func TestOperatorWorktreeWindow(t *testing.T) {
 	if err != nil || !strings.Contains(string(md), "\n"+worktreeExample+"\n") {
 		t.Fatalf("missing executable QUICKSTART example: %v", err)
 	}
-	for _, window := range []string{"before reservation", "before add", "before readback", "before observe", "git refused", "path only", "ref changed", "head changed", "readback error", "observe error", "add error after checkout", "redirected checkout", "registration locked", "ref read skew"} {
+	for _, window := range []string{"before reservation", "before add", "before readback", "before observe", "git refused", "success then absent", "path only", "ref changed", "head changed", "readback error", "observe error", "add error after checkout", "redirected checkout", "registration locked", "ref read skew"} {
 		t.Run(window, func(t *testing.T) {
 			r := newWorktreeStartRig(t, false)
 			initial := observeProvision(t, r.f.r, r.creds)
@@ -344,6 +344,11 @@ func TestOperatorWorktreeWindow(t *testing.T) {
 				case window == "add error after checkout" && step == "before readback":
 					r.state(t, "running", "reserved", "", 1, true)
 					atFailure = observeProvision(t, r.f.r, r.creds)
+				case window == "success then absent" && step == "before readback":
+					r.state(t, "running", "reserved", "", 1, true)
+					presetGit(t, r.f.r.f, "worktree", "remove", r.path())
+					presetGit(t, r.f.r.f, "branch", "-D", "world/ep")
+					atFailure = observeProvision(t, r.f.r, r.creds)
 				case window == "git refused" && step == "before add":
 					dir := filepath.Join(r.f.r.f.repo, ".git", "refs", "heads", "world")
 					if err := os.Mkdir(dir, 0700); err != nil {
@@ -433,6 +438,13 @@ func TestOperatorWorktreeWindow(t *testing.T) {
 					t.Fatal(err)
 				}
 				r.exec(t, "DROP TRIGGER fail_observe")
+			case "success then absent":
+				r.state(t, "uncertain", "reserved", "", 2, false)
+				got := observeProvision(t, r.f.r, r.creds)
+				got.credentials = atFailure.credentials
+				if !reflect.DeepEqual(atFailure, got) || err == nil || err.Error() != "worktree outcome uncertain; "+worktreeExample {
+					t.Fatalf("successful add with absent readback lost uncertainty or evidence: %v", err)
+				}
 			case "path only", "ref changed", "head changed", "readback error", "add error after checkout", "redirected checkout", "registration locked", "ref read skew":
 				r.rows(t, "uncertain", "reserved", "", 2)
 				got := observeProvision(t, r.f.r, r.creds)
@@ -472,6 +484,8 @@ func TestOperatorWorktreeWindow(t *testing.T) {
 			}
 			receipt, version := "stopped: retained", 2
 			switch window {
+			case "success then absent":
+				receipt, version = "stopped: absent", 3
 			case "before add":
 				receipt = "stopped: absent"
 			case "path only", "ref changed", "head changed", "readback error", "add error after checkout", "redirected checkout", "registration locked", "ref read skew":
diff --git a/host/store/operator_worktree_guards_test.go b/host/store/operator_worktree_guards_test.go
index 01d61856..61e24c04 100644
--- a/host/store/operator_worktree_guards_test.go
+++ b/host/store/operator_worktree_guards_test.go
@@ -11,7 +11,7 @@ import (
 )
 
 func TestReserveOperatorWorktreeGuards(t *testing.T) {
-	for _, name := range []string{"valid", "kind", "headless", "store", "approver", "challenge", "claimed"} {
+	for _, name := range []string{"valid", "kind", "headless", "store", "approver", "challenge", "claimed", "owner", "identity", "expiry"} {
 		t.Run(name, func(t *testing.T) {
 			c := openTestCredentials(t)
 			ctx := testCtx(t)
@@ -60,7 +60,15 @@ func TestReserveOperatorWorktreeGuards(t *testing.T) {
 			proof := signed.MintProof()
 			cas.Digest, cas.Generation = digest, row.Generation
 			operatorExec(t, c, "UPDATE operator_actions SET action=?, digest=?, kind=?, artifact=?, approved_by=?", body, digest[:], action.Kind, artifact, signed.ApprovedBy())
+			identity, now, want := "fixture", int64(1590), ErrOperatorConflict
 			switch name {
+			case "owner":
+				operatorExec(t, c, "UPDATE operator_actions SET owner='other'")
+			case "identity":
+				identity = ""
+			case "expiry":
+				now, want = row.ExpiresAt, ErrOperatorExpired
+				operatorExec(t, c, "UPDATE operator_runner_state SET seen_at=?", now)
 			case "kind", "claimed":
 				if name == "claimed" {
 					// The SQL CAS also rejects running. Make reaching Bind observable.
@@ -81,7 +89,7 @@ func TestReserveOperatorWorktreeGuards(t *testing.T) {
 				operatorExec(t, c, "PRAGMA ignore_check_constraints = OFF")
 			}
 			before := getOperator(t, c, cas.ID, 1590)
-			key, err := c.ReserveOperatorWorktree(ctx, OperatorMintAttempt{CAS: cas, Runner: runner, Catalogue: 1}, "fixture", proof, func() int64 { return 1590 })
+			key, err := c.ReserveOperatorWorktree(ctx, OperatorMintAttempt{CAS: cas, Runner: runner, Catalogue: 1}, identity, proof, func() int64 { return now })
 			if name == "valid" {
 				if err != nil {
 					t.Fatal(err)
@@ -93,7 +101,7 @@ func TestReserveOperatorWorktreeGuards(t *testing.T) {
 				}
 				return
 			}
-			if !errors.Is(err, ErrOperatorConflict) || key != (OperatorStepKey{}) ||
+			if !errors.Is(err, want) || key != (OperatorStepKey{}) ||
 				!reflect.DeepEqual(before, getOperator(t, c, cas.ID, 1590)) {
 				t.Fatalf("refusal changed action or returned the wrong class: %v", err)
 			}
@@ -108,7 +116,7 @@ func TestReserveOperatorWorktreeGuards(t *testing.T) {
 }
 
 func TestDisposeOperatorUncertainStepGuards(t *testing.T) {
-	for _, name := range []string{"kind", "owner", "version"} {
+	for _, name := range []string{"kind", "owner", "version", "action_kind"} {
 		t.Run(name, func(t *testing.T) {
 			c := openTestCredentials(t)
 			ctx := testCtx(t)
@@ -133,6 +141,9 @@ func TestDisposeOperatorUncertainStepGuards(t *testing.T) {
 				key.Owner = "different"
 			case "version":
 				bad.Version++
+			case "action_kind":
+				operatorExec(t, c, "UPDATE operator_actions SET kind='finish'")
+				before = getOperator(t, c, cas.ID, 1590)
 			}
 			if err := c.DisposeOperatorWorktree(ctx, bad, OperatorWorktreeUncertain, 1591); !errors.Is(err, want) ||
 				!reflect.DeepEqual(before, getOperator(t, c, cas.ID, 1590)) {
@@ -165,3 +176,48 @@ func TestOperatorStepReceiptAtLimit(t *testing.T) {
 		t.Fatalf("at-limit receipt was not stored exactly: %v, size %d", err, len(after.Receipt))
 	}
 }
+
+func TestReadOperatorWorktreeGuards(t *testing.T) {
+	for _, name := range []string{"uncertain", "duplicate"} {
+		t.Run(name, func(t *testing.T) {
+			c := openTestCredentials(t)
+			ctx := testCtx(t)
+			cas, runner := operatorClaimFixture(t, c)
+			if err := c.ClaimOperatorAction(ctx, cas, runner, 1590); err != nil {
+				t.Fatal(err)
+			}
+			cas.Generation++
+			key := OperatorStepKey{cas, runner.Nonce, "worktree", "first", 1}
+			if err := c.ReserveOperatorStep(ctx, key, 1590); err != nil {
+				t.Fatal(err)
+			}
+			operatorExec(t, c, "UPDATE operator_actions SET state='running'")
+			var second OperatorStepKey
+			var want error
+			if name == "uncertain" {
+				if err := c.DisposeOperatorWorktree(ctx, key, OperatorWorktreeUncertain, 1591); err != nil {
+					t.Fatal(err)
+				}
+				key.Action.Generation++
+				key.Version++
+			} else {
+				second = key
+				second.Identity = "second"
+				if err := c.ReserveOperatorStep(ctx, second, 1590); err != nil {
+					t.Fatal(err)
+				}
+				want = ErrOperatorMissing
+			}
+			before := getOperator(t, c, cas.ID, 1591)
+			got, signed, err := c.ReadOperatorWorktree(ctx, cas.ID)
+			if !errors.Is(err, want) || !reflect.DeepEqual(before, getOperator(t, c, cas.ID, 1591)) ||
+				(want == nil && (got != key || signed.ID != cas.ID || signed.Kind != "start")) {
+				t.Fatalf("wrong reservation read or action changed: %+v, %v", got, err)
+			}
+			assertReservedOperatorStep(t, c, key)
+			if name == "duplicate" {
+				assertReservedOperatorStep(t, c, second)
+			}
+		})
+	}
+}
````
