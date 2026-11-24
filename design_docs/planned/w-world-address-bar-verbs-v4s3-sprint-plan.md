# Sprint plan — w-world-address-bar-verbs-v4, sprint S3 (iteration 428, row 243)

> **THE SPRINT — row 243, phase V4, sprint S3.** Reviewed design r1, §§3.1–3.2 and 5.2. **M1** installs credentials schema 10 and pins transactional upgrade parity. **M2** adds inert requests, operation and review dedup, capacity, expiry, prune and store exclusion. **M3** adds confirmation and claim CAS, claimed receipt recheck and fd-level silence pins. The executor applies three cumulative patches. S3b's separately compiled receipt-write slice is in Appendix C/P.
>
> **Size.** Expected cumulative changes are **241 → 738 → 1,396** against `92e1c44e` with renames on and the `cmd host docs app` filter. The hard limit is 1,400. The full staged draft is **1,528**, including the 132-line S3b slice. Design budget is 890 raw / 1,202 calibrated. The tracked product diff stays inside `cmd` and `host`.
>
> **Authority.** D-WORLD-119 A and D-WORLD-120 A, both attended by Mark in `82f05004`. Schema 10 includes both kinds, the awaiting-PR phase, V3 link and finish step identities. The schema has no start or finish executor. No route is registered. Install, launch, mint, finish, publication and browser work remain in their assigned sprints.

## 0. Binding constraints

- D-WORLD-109. Test keys, generated artifacts and temporary stores only. The fixture probes never access studio or home service paths, a real key, terminal, GitHub or browser. Use the pinned `AILANG_BIN`. Keep network listener legs with the controller and classify bind denials in the sandbox.
- Use the reviewed r1 design and the S2b handoff from PR #611, merge `07b30a5e`, judge PASS 92. The existing new-route count, three surface POSTs, `surfaceCSP`, `app.js` and registered serving path are pinned in §6. The only startup effect is the authorized credentials migration.
- Keep `.ail`, scripts, launchd, skills, packages, plugins, integrations, world-publish and attended collision paths byte-identical. The two existing CLI lines pin the old accepted schema floor for bound writes on still-v9 stores. No CLI flags or help change.

## 1. Pre-flight

### 1.1 Baseline at `92e1c44e`

````text
TOUCHED_PACKAGES: 
LEG status rc=0 0s
LEG gofmt rc=0 1s
LEG build rc=0 4s
LEG vet rc=0 10s
LEG fence rc=0 36s
LEG verify_ail rc=0 54s
LEG pkg-host-approval rc=0 1s
LEG race-host-approval rc=0 16s
LEG pkg-host-approval-approvaltest rc=0 1s
LEG race-host-approval-approvaltest rc=0 3s
LEG pkg-host-authority rc=0 2s
LEG race-host-authority rc=0 42s
LEG pkg-host-store rc=0 25s
LEG race-host-store rc=0 38s
LEG pkg-cmd-ailang-worldd rc=1 2s UNINFORMATIVE UNDER SANDBOX (loopback bind; inspect for independent failures)
LEG race-cmd-ailang-worldd rc=1 22s UNINFORMATIVE UNDER SANDBOX (loopback bind; inspect for independent failures)
LEG inv-host-approval rc=0 1s
LEG inv-host-approval-approvaltest rc=0 1s
LEG inv-host-authority rc=0 1s
LEG inv-host-boundary rc=0 4s
LEG inv-host-daemon rc=0 2s
LEG inv-cmd-ailang-worldd rc=0 1s
LEG inv-host-store rc=0 2s
DAEMON_ROUTES: 50 (base 50) SURFACE_POSTS: 3
UNTOUCHED_VIOLATIONS: none
UNCHANGED: host/surfacepage/app/app.js
UNCHANGED: host/daemon/surface_page.go
UNCHANGED: host/daemon/daemon.go
NEW_PRODUCTION_CALLERS: none (authority.RecheckVerb control: 1)
LEG verify_go SKIPPED (GATES_SKIP_VERIFY_GO=1)
GATES_DONE 270s FAILS: pkg-cmd-ailang-worldd_sandbox_uninformative race-cmd-ailang-worldd_sandbox_uninformative
````


### 1.2 Codebase claims checked

| # | Claim | Command | Observed |
|---|---|---|---|
| C1 | Pinned base, accepted decisions and predecessor | `git rev-parse HEAD`, `git show --format=short --stat 82f05004`, `git log -1 --format=%s 07b30a5e` | Base `92e1c44ebd8dbbdb2b962ef54aed0dd9a0bf2456`; Mark is the cited decision commit's author and both ledger edits say RESOLVED. Controller's `--open` result and attended provenance are supplied inputs. Merge subject says #611 from iteration 409. |
| C2 | Design §3.1 cited seams and Appendix C handoff | `rg -n 'JudgeOperatorAction\|RecheckOperatorAction\|VerifiedOperatorAction\|OperatorJudgeInput\|Binds\|ApprovedBy\|PasskeyUse\|Rechecked' host/authority/judge_operator.go`; `rg -n 'sshApprover\|passkeyApprover\|JudgeVerbSSH\|JudgeVerbAssertion\|RecheckVerb' host/authority/judge_verb.go` | `judge_operator.go:19,33,61,77,95,102,109,116` gives proof, input, two judges, four accessors. The V3 functions named by r1 are still at `judge_verb.go:30,44,63,82,107`. `surface_confirm.go:97` handler and `:37–39` judge/consume seams. `ApproverByKey` at `approvers.go:140`, `PasskeyByCredentialID` at `passkeys.go:191`. |
| C3 | Canonical envelope and boundary times | `rg -n 'DecodeOperatorAction\|OperatorReviewDigest\|OperatorSignatureDigest\|OperatorActionLifetime\|CanonicalOperatorAction\|OperatorActionMessage' host/approval/operator_action.go` | `:32` 10 minutes, `:189` canonical, `:202` strict decode, `:288` review and `:302` signature digest, `:315,329` message and hash. The review digest erases only the assigned action id. S3 recomputes from the typed intent and uses a separate 32-byte challenge. |
| C4 | V3 is the closest analogue | `rg -n 'credentialsSchemaVersion\|credentialsTables\|upgradeCredentials\|initCredentialsSchema' host/store/credentials.go`; `rg -n 'verbRequestsSQL\|MaxPendingVerbRequests\|PutVerbRequest\|ConsumeVerbRequest\|useApprover\|usePasskey' host/store/{verb_requests,approvers,passkeys}.go` | Version 9 at `credentials.go:45`, exact sorted table list at `:74`, transactional migrations at `:353–398`, fresh at `:400–418`. V3 table at `credentials_schema.sql:118` and v8→v9 copy at `verb_requests.go:24`; cap 16 at `:42,146`; request at `:121`; consume CAS at `:239`; transaction helpers `approvers.go:245` and `passkeys.go:225`. V4 keeps the transactional format and typed refusal, declines supersession of a competing active action, retains distinct tap aliases and uses a shared DDL source. |
| C5 | Parity and old version gates | `rg -n 'TestCredentialsSchema\|TestCredentialsNewerThanNineIsRefused\|tableInfo\|tableDDL' host/store/*test.go` | V7 `surface_approvals_test.go:57`, V8 `key_decisions_test.go:59`, V9 `verb_requests_test.go:105,155` plus `tableInfo:83`. Canonical table set, DDL and user_version at `credentials_test.go:76,115`. M1 compares all normalized SQLite schema objects including indexes and autoindexes, injects mid-DDL failure with an index collision, checks a retained bound credential and repeated upgrade. Probe below exercises actual base binaries. |
| C6 | Real fixture generators | `rg -n '^func .*Sign\(|^func .*Assert\(|^func (jgRig\|ojSSH\|ojAssert\|boundedTestContext)' host/approval/approvaltest host/authority` | `approvaltest/signer.go:46` Sign, `passkey.go:96` Assert. `judge_test.go:24–25,30` seeds `jgPasskey`, `jgSigner` and enrolls their public facts in a temporary store. `judge_operator_test.go:50,59` `ojSSH` and `ojAssert` sign the real operator message and challenge hash. `approvaltest/sign_crosscheck_test.go` byte-checks against ssh-keygen. No confirmation artifact is synthesized by hand. SQL-only schema/corruption fixtures use specimen payloads without treating them as authority. |
| C7 | Dependency, symbol, route and store inventories | `rg -n '^func Test.*(Allowlist\|ProductionContextRoots\|ProductionGoSurface\|SurfaceTakesContext\|CredentialWriterSites\|WriterLock\|ImportsOnlyTheStandardLibrary\|TestOnly)' cmd host -g '*test.go'` and §6 `inv-*` legs | `daemon_test.go:899` TestDaemonDependencyAllowlist; `allowlist_world_test.go:758` TestWorldBoundaryDependencyAllowlist; `context_roots_test.go:129,234` TestProductionContextRoots and TestProductionGoSurface; `deadline_guard_test.go:43` TestStoreSurfaceTakesContext reflects `*Store`, separate from credentials; `credentials_test.go:277` TestCredentialsMethodsRequireADeadline; `cli_credential_writers_test.go:18`, `surface_approve_post_guard_test.go:132` writer sites; `writer_lock_test.go:205` and `writer_lock_held_test.go:12` lock behavior. Approval import and signer test-only inventories remain. No production context roots, modules, packages or writer opens are added. DDL table inventories, hardcoded versions and historical temporary-store fixtures change in M1. Existing AST inventories stay byte-identical. |
| C8 | §3.1 / §5.2 terminal lock and no fallback | `rg -n 'O_EXCL\|LOCK_NB\|lockEpisode\|func (runPresetEpisodeStart\|runEpisodeFinishWith\|launchEpisode)' cmd/ailang-worldd/{episode_preset,episode_finish,episode_launch,service}.go host/{store/writer_lock_unix,coordinator/effectful}.go` | `episode_preset.go:53`, `episode_finish.go:63`, `episode_launch.go:258,442`: first two have no common episode flock; launch uses O_EXCL PID file. Positive controls are the service lock at `service.go:154,648`, the context/channel lock at `effectful.go:98` and credentials non-waiting flock next to `writer_lock_unix.go:41`. S3 owns existing store-level contention and cross-kind admission. |
| C9 | Service install, upgrade and rollback | `rg -n 'upgrade\|rollback\|Credentials\|Kickstart\|Bootout\|Rename' cmd/ailang-worldd/service{,_upgrade}.go`; `rg -n 'ensureCredentialsStore\|UpgradeCredentials\|CredentialsSchemaVersion' host/daemon/{daemon,credentials}.go` | Install passes original credentials path to serve (`service.go:435,448`). `:322,348` admits upgrade; no standalone rollback verb. Candidate check is build / platform / two tool pins (`service_upgrade.go:37,57–61,165–174`); swap and restart `:210–240`. Serve migrates before exposing health (`daemon.go:817`, `credentials.go:59–77`). Automatic rollback renames .prev over current, restarts it and checks old build health (`service_upgrade.go:243–268`), then stops and returns 3 if old build does not serve (`:270–291`). |
| C10 | Scope fences and cut | §6 `paths`, `rg -c 'mux\.(Handle\|HandleFunc)\(' host/daemon/daemon.go`, `git diff --no-index -M --numstat <archive>/{cmd,host} <stage>/{cmd,host}` run per root | 50 registrations, 3 surface POSTs. Caller search includes positive control `authority.RecheckVerb`. Product diff at every stage is under `cmd` and `host`. See §1.5. |
| C11 | Binding symbols used by the patch and no new root | `rg -n '^func (requireDeadline\|testCtx\|rawDB\|oneBoundRow\|setPositiveFixtureVersion)\b\|^type (Credentials\|OlderCredentialsStoreError\|ApproverUse\|PasskeyUse)\b' host/store`; `rg -n 'type (OperatorAction\|OperatorOptions\|OperatorResolved\|SignedOperatorAction)\b\|func .*Canonical\(' host/approval` | Existing `requireDeadline` `store.go:250`, `testCtx` `mint_requests_test.go:175`, `oneBoundRow:55`, `rawDB` `schema_version_test.go:169`, version fixture setter `:297`. Existing store structs and SQL lifecycle are in credentials.go, mint_requests.go, approvers.go and passkeys.go. Approval action fields are in `operator_action.go:51–76`, signed proof/accessors at `operator_signed.go:17,132,153,162`, COSE canonical at `cose.go:102`. NEW schema SQL, action/CAS/owner/request types and files per §5, proof transport and binding plus test helpers in the patch. The positive root control is the existing `budgetContext`; the §6 inventory scans all production sources. |
| C12 | Strict time and deterministic contention | `rg -n 'now\|expires_at\|seen_at\|WaitGroup\|Sleep\(' host/store/operator*.go` on scratch | All action clocks are int64 inputs, including expiry, claim readiness and retention. Two simultaneous claim callers start on one barrier with a finite store context. No sleeps, no open-ended subprocess, no listener required for behavioral mutants. |
| C13 | Round trip and snapshot protocol | §2 extraction plus §7 snapshot, controller `ctl_roundtrip.sh` and `commit_ms.sh` read at HEAD | Extracted from the plan with §2 awk and applied to a fresh archive. M1, M2, M3 and S3b full-tree diff -r all empty. No patch rejects or backups under cmd or host. Snapshot manifests verified (M1 19, M2 21, M3 28 files, zero deletions). |
| C14 | Actual older binary reads v10 | Build `open-probe.go` from both the unmodified archive and M1, create a scratch v10 DB, call old probe with `read`, `open`, `upgrade` and compare SHA-256 | M1 probe reported version=10. The actual v9 build from the base returned CredentialsSchemaError with “user_version 10 is newer than this binary's 9; use a compatible binary” for each of read/open/upgrade, rc 1 each. Before/after SHA-256 ec243a1f689e1ee7f5881c4179899cf85a34928062dd4a1e0b68c9666f8908cd. Scratch source and binaries: open-probe.go and probe-v9/probe-v10 in plan-scratch. |
| C15 | Negative mechanism and service safety | `rg -n 'case version > credentialsSchemaVersion\|requireDeadline\|db.SetMaxOpenConns\(1\)' host/store/{credentials,store}.go` | `credentials.go:344–345` refuses newer versions, `:304` serialized per-handle pool; production second-writer prevention comes from canonical exclusive flock, not the test goroutines. §8 records where current rollback fails. |
| C16 | Exact fixture and readback references | `rg -n 'type (Credentials\|ApproverUse\|PasskeyUse\|OlderCredentialsStoreError) struct\|func (normalizeDDL\|openTestCredentials\|setPositiveFixtureVersion\|userVersion\|jgExec)\b' host/store host/authority` | `Credentials` and its display/db/version/lock fields at credentials.go:144–153. `CanonicalPath:176`, `Close:274`, `MintSession:441`, `ResolveSession:522`; `Approver:135` in approvers.go; `Passkey:200` and `PasskeyUse:216` in passkeys.go. `OlderCredentialsStoreError:67` and `ApproverUse:300` in mint_requests.go. `MaxSurfaceApprovalArtifactBytes:38` in surface_approvals.go is 16 KiB. `openTestCredentials:121`, `normalizeDDL:854` in journal_test.go, `userVersion:49` in credentials_binding_test.go and `jgExec:52` in authority/judge_test.go. The package-quality stderr advisory and isolated retry are in F10 and §1.5. |
| C17 | Remaining runner and herdr seams | `rg -n 'func (launchEpisode\|runEpisodeFinishWith\|runPresetEpisodeStart)\|operator-runner\|operator_actions' cmd host -g '*.go' -g '!**/*_test.go'`; `git show HEAD:integrations/herdr/ailang-world/bin/world-config-init \| nl -ba \| rg 'codex='` | At base the three terminal entry points are the positive controls and no operator runner/action code exists. The herdr model string is still `WORLD_MODELS="claude=claude-opus-5-5 codex=gpt-6.1-sol pi="` at line 96. The design's unattended uncertain and unowned-child cases remain runner recovery work. |

### 1.3 Planner decisions

- **S3-D1. One shared migration DDL.** Embed the schema-10 DDL once and execute it under the same transaction as the version stamp, both for fresh and upgraded stores. Keep V3 table DDL and version-9 consume semantics. Keep the previous version floor for bound writes and the CLI checks.
- **S3-D2. Carry the finish shape.** A start-only migration would make the forthcoming V3 link/awaiting phase require another irreversible bump. The v10 schema accepts finish actions and their checkpoint identities. A linked V3 id is unique and finish-only. It is a logical reference, preserving V3's existing removal of expired request rows and letting S10 classify a missing card. There is no prep or dispatch in this plan.
- **S3-D3. Proof transport preserves the authority boundary.** `host/authority` already imports `host/store`, so the store accepts the authority proof through its `OperatorConfirmation` method. The method returns the sealed law object. The store requires exact artifact hash, canonical bytes, digest, challenge and expiry, and compares the enrolled key and passkey metadata under its writer transaction. Counter and exact proof identity are read from the same object that backed `VerifiedOperatorAction.PasskeyUse` and `ApprovedBy`. Bare interface-supplied identity, public key and counter fields cannot authorize a write. Rechecks cannot be consumed as a confirmation.
- **S3-D4. Persist retry identities in the three-table footprint.** A reserved `operation` step kind permanently binds every tap's id to an action, even when the action id or challenge assigned to the second tap differs. The row retained after pruning is its tombstone and still returns a stable id and state. Reserve/observe step writes are S3b. Per-action effect identities are distinct from operation aliases.
- **S3-D5. Three patches and a vertical S3b cut.** Preserve independently green schema, admission and CAS checkpoints. The receipt API and its restart test are one measurable follow-up. No feature dispatch can activate before its gate. Full terminal lock and status cannot run against a runner that has yet to be built.

### 1.4 Premises found false or incomplete at HEAD

| # | Premise | Measurement and consequence |
|---|---|---|
| F1 | Audit start-only residual, and S2b's OPEN ledger | C1 `git show 82f05004` resolves both A. Design §3.1 calls for a linked V3 id, §3.2 explicitly contains Awaiting PR confirmation, §3.3 carries checks/commit/card/push/label/PR/record/cleanup, and §4.3 separates prep from the V3 proof. Land these fields/phases/kinds now. The direct DDL comparison in §1.5 measures **15 changed lines (9 add, 6 delete, 3 net)** for the finish shape. |
| F2 | Terminal already shares a kernel episode lock | C8 finds no such lock. The store lock half is real and exercised with two distinct open descriptions plus an episode collision after reopen. The full terminal refusal, naming the action and remedy, and zero-effect counter go to S4. |
| F3 | Authority proof can be a concrete store import | C2 shows the reverse import. A direct import cycles. S3-D3 passes the unforgeable signed object and binds fresh enrollment and bytes at consume. |
| F4 | Existing rollback can safely restore the previous build | C9 and C14 show old start fails after the committed migration. Rollback stops the service, leaving a valid v10 store and the compatible candidate. Plan §8 records recovery review for the controller. |
| F5 | S2b route-count inventory still says 47 | C10 measured 50 at HEAD and 3 surface POSTs. Final must equal that base. |
| F6 | The schema-version bump only affects schema code | `rg -n 'CredentialsSchemaVersion\|c.version < credentialsSchemaVersion' cmd host -g '*.go'` finds `episode_preset.go:241`, `session_new.go:381` and credentials `:460` using the current number as the binding write floor. Pin the existing floor separately so pre-upgrade v9 behavior stays intact. |
| F7 | S2b's fd capture catches builtin println | `nl -ba host/authority/judge_operator_silent_test.go` shows `:19–32` swapping Go file variables only. M3 uses dup/dup2 on fd 1 and 2 and proves builtin output capture. See §4. |
| F8 | Second daemon/runner writer inventory must move in S3 | C7/C8 and §6 caller count show only methods on the existing credentials handle. No new writer acquisition site. S4/S7 must move any additional writer sites and non-waiting coverage together. |
| F9 | Historical fixture downgrades remain valid after adding tables | Whole-package fixtures that start fresh and explicitly downgrade must remove v10 tables. The provisional whole runs failed `TestCredentialsUpgradeToV4`, `TestCredentialsUpgradeToV6`, `TestJudgeVerbSSH`, `TestJudgeVerbAssertion`, `TestRecheckVerb` and other old judge fixtures with the table gate. `rg -n 'DROP TABLE verb_requests|requireExactTableNames' cmd host -g '*test.go'` identified the fixtures. M1's exact test changes in §5 keep those branches runnable. |
| F10 | The bare AILANG quality gate consumes clean merged JSON | Direct pinned `pkg quality --json packages/world-core` succeeded on the base and staged package with JSON on stdout, but the default Observatory state emitted a size advisory on stderr. The existing script merges those streams; Python then refused “Extra data”. §6 and §7 now scope `AILANG_STATE_DIR` under iteration 428. The product scripts stayed byte-identical. |
| F11 | The unified-diff line count predicts Git numstat exactly | Final `git diff --no-index -M --numstat <clean base export> <compiled stage>` gives 241/738/1,396. The equivalent patch-text count was four lines higher because the two diff algorithms aligned fixture edits differently. The controller uses Git's count. |

### 1.5 Size and CI race cost

`git diff -M --numstat 92e1c44e HEAD -- cmd host docs app` is the controller's exact base-relative filter. Count additions plus deletions. Scratch used fresh archive and pathwise `git diff --no-index -M --numstat` because the planner performs no Git writes. The unfiltered product diff is the same. The planning artifacts are controller inputs kept outside the application diff. If counted in the branch's unfiltered diff, the two plan files add their own line counts separately.

| Package | Base -race package s | Final -race package s | Added s | Scope |
|---|---:|---:|---:|---|
| `host/store` | 36.866 | 33.547 | -3.319 | whole |
| `host/authority` | 5.575 | 7.859 | +2.284 | whole |
| `host/approval` | 2.027 | 2.012 | -0.015 | whole |
| `cmd/ailang-worldd` | 5.511 | 6.000 | +0.489 | focused, whole bind denied |
| `host/daemon` | 4.789 | 5.906 | +1.117 | focused, whole bind denied |

Sub-milestone increments (additions + deletions) are M1 241, M2 497, M3 658, cumulative 241/738/1,396. S3b is 132 additional lines. Each stage was vetted and formatted, its §3 command passed, and its snapshot manifest and full-tree round trip were verified. The root compare against a fresh archive confirms the unfiltered product total is also 1,396 (1,350 additions, 46 deletions). The sprint plan and its JSON are reported separately as planning artifacts.

The base host/store and host/authority first paired race run also passed at 31.424 s and 5.713 s. A loaded final race hit the pre-existing `TestCommitLandedReadsTheLogRowAtItsIndex` context deadline, then the isolated run passed. CI's 480 s per-package deadline has at least 440 s headroom in the three fully measured packages. Whole `host/daemon` and `cmd/ailang-worldd` were **UNINFORMATIVE UNDER SANDBOX** at the base and final. The exact whole-package commands are `go test -count=1 -timeout 8m ./host/daemon ./cmd/ailang-worldd` and the same with `-race`. Their v4/v5/v6/v8 fixture and inventory selections passed both plain and race.

Final build, vet, compile fence and named inventories passed. The first source walk overlapped cleanup of a generated `.ailang` cache and then passed on the clean stage. The original final `verify_ail` ended at the package-quality JSON parser because the default Observatory state wrote a size advisory to stderr. stdout alone parsed on both base and final. The rerun with `AILANG_STATE_DIR` under plan-scratch passed in 55.00 s, including 16 required identities, all 40 named core tests, all 503 named se-tools tests and PUB011. The §6 gate now keeps this state under the log root. Full `verify_go` stayed skipped in the planner's run, as allowed by `GATES_SKIP_VERIFY_GO=1`. The controller runs the full script outside the sandbox.

The finish-ready DDL costs 9 insertions and 6 deletions compared with the equivalent start-only schema (15 changed lines, 3 net). Measured with git diff --no-index --numstat plan-scratch/{start,finish}-shape.sql. It carries one nullable linked V3 id and its constraints, the finish kind and awaiting phase, and 8 finish step kinds.

## 2. How the executor applies a sub-milestone

Each Appendix P block begins with four backticks, `diff` and the exact milestone id. Use one bash shell.

````bash apply.sh
#!/bin/bash
set -euo pipefail
M=M1 # choose M1, M2, M3 in order
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter428-row243-v4s3
PLAN="$W/design_docs/planned/w-world-address-bar-verbs-v4s3-sprint-plan.md"
T="$HOME/.ailang/state/world-iter428/exec-patches"
mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
/usr/bin/patch -p1 -E --dry-run < "$T/$M.patch"
/usr/bin/patch -p1 -E < "$T/$M.patch"
if find cmd host \( -name '*.orig' -o -name '*.rej' \) -print | grep .; then exit 1; fi
git status --short --untracked-files=all
````


A failed dry run, patch failure or unexpected path stops the executor. Fetch the matching acceptance script from §3 and follow §7.

## 3. Milestones

````bash accept.sh
#!/bin/bash
set -euo pipefail
cd "$1"
M=$2
test "${3:-fast}" = fast
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export AILANG_STATE_DIR="$HOME/.ailang/state/world-iter428/accept-ailang-state"
mkdir -p "$AILANG_STATE_DIR"
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR WORLD_UPDATE_OPERATOR_ACTION_GOLDEN WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN
test -x "$AILANG_BIN"
case "$M" in
  M1) store='TestOperatorSchema10|TestCredentialsSchemaMatchesCanonicalDDL|TestCredentialsSchema9VerbRequests|TestCredentialsNewerThanTenIsRefused' ;;
  M2) store='TestOperatorSchema10(Constraints)?|TestOperatorAdmissionBoundaries|TestOperatorDuplicateTap|TestOperatorReconnect|TestOperatorPrune|TestOperatorTerminalExclusion' ;;
  M3) store='TestOperatorSchema10(Constraints)?|TestOperatorAdmissionBoundaries|TestOperatorDuplicateTap|TestOperatorReconnect|TestOperatorPrune|TestOperatorTerminalExclusion|TestOperatorClaimBeforeExpiry|TestOperatorClaimRace' ;;
  *) exit 2 ;;
esac
go test -v -count=1 -timeout 8m ./host/store -run "^($store|TestProductionContextRoots|TestProductionGoSurface)$"
if [ "$M" = M3 ]; then
  go test -v -count=1 -timeout 8m ./host/authority -run '^TestOperator(ConfirmReplay|ConfirmCrash|ConfirmReadOnly|ConfirmStoreBinding|RunnerRechecks|JudgesWriteNothing|StoreWritesNothing)$'
  go test -count=1 -timeout 8m ./host/approval ./host/approval/approvaltest
fi
echo "ACCEPT_OK $M"
````


### M1 — credentials schema 10 and migration parity (241 cumulative)

Three user-version-10 tables, all schema objects equivalent on fresh and upgraded fixtures, one-way old-open refusal and atomic abort. Existing v9 app/session paths retain their v9 floor. Advances schema without enabling effects.

### M2 — inert admission, expiry, dedup and pruning (738 cumulative)

Typed request/read/prune, exact review-based aliases, per-episode exclusion across kinds, 16 unexpired pending rows, projected expiry without write-on-read, retained tombstones/reservations after prune. Existing credentials writer lock remains the only writer gate. Includes direct SQLite constraint refusals.

### M3 — authority proof consumption, claim and fd tests (1,396 cumulative)

Atomic passkey and approval use on confirmation with exact stored identity and replay. Claim checks digest, phase generation, live current runner nonce/generation/heartbeat/per-kind readiness, strict expiry and exclusive ownership. Claimed data is rechecked with the runner's own frozen action and times. Competitors are finite and single-winner. All these entry points stay unregistered.

### 3.1 Branch traces

| Input / failure | Store or authority result | Durable assertion |
|---|---|---|
| Incorrect operation id, store, action, issued-at/future time or >600 s | input or approval validation sentinel | Existing request remains inert, count/steps unchanged. |
| Existing operation id, same review, or independently assigned tap | original receipt/status | Same original id and challenge, alias written only once, persists across reopen. |
| Changed review on same operation, competing episode, seventeenth pending | conflict or capacity | Existing action and aliases unchanged. Exact expiry frees pending and episode on next successful write. |
| Pending row and proper enrolled SSHSIG or WebAuthn proof | exact artifact and format:name, one state transaction | Handed back state and proof; one use; source-store key and metadata pinned. |
| Missing id, wrong digest/generation, clock or expiry, nil/swapped artifact | typed refusal | Requested state, no artifact, approver last-use and counter unchanged. |
| Claim/counter races, read-only DB or failed last statement | typed SQLite or store refusal | State and approver-use updates roll back together. Replayed confirmation returns durable result without increment. |
| Clock before issue, expiry equality, wrong runner or phase, cancellation | claim refusal | Action generation, owner, state and alias/step counts pinned. Two callers race and exactly one changes row. |
| Recheck on a claimed row with wrong frozen action, approver, live deadline or artifact | proof refusal | No `VerifiedOperatorAction`; only the signed artifact can produce proof. |
| 30 day boundary, observed effect, uncertain reservation, old operation retry | prune old evidence only | Terminal ids/state and reservations survive. Exact boundary is retained. |
| Existing credentials writer while another attempts open | `WriterAlreadyActive` plus no competing rows | Release and reopen is the positive lock control; same active episode still refuses. |

## 4. Mutation table

The ledger records the named test and whole-package verdict for each materialized mutant. Existing static inventory kills are counted separately when relevant. Ship the compiled patch and verify the round trip before running the mutant loop. See the result legend immediately above the table.

**Final discrimination.** 138 scoped mutants. **129 killed by both their named behavioral test and its whole package, 9 named equivalents, 0 non-equivalent survivors, 0 unmeasured, 0 pending.** The separate v9-base control `BASE-LOG-IN-JUDGE-NOFAIL` survived both named and whole authority tests. The final `LOG-IN-JUDGE-NOFAIL` and the five store fd canaries fail their behavior tests. `CONFIRM-IDENTITY` fails both the replay test and claimed recheck.

The scratch ledger is `$HOME/.ailang/state/world-iter428/plan-scratch/mutants/ledger.txt`. `ledger.jsonl` retains every run and anchor, including retests. Rows below use the latest run. The whole test package includes its static tests. Confirmation's real store/authority tests live in `host/authority`, which supplies the sealed proof without an import cycle. For cross-package mutations the extra whole source package column shows the separately run package with the modified line. All 28 separate source-package runs survived because the discriminating integration tests are in host/authority. Those verdicts are listed separately from the whole test-package results. No scored kill depends on static checks. `RECHECK-STATUS` also hit the constructor inventory but independently failed `TestOperatorRunnerRechecks`. A named writer-lock mutant finishes via bounded refusal check. Its whole-package mutation hits the bounded Go test timeout in the existing held-writer migration test.

**Legend.** K killed, S survived, E independently justified equivalent on the supported stored shape and transaction, NM unmeasured. The four confirming SQL tests separately perturb the transaction: trigger-abort, generation drift, digest drift and exact expiry drift.

| id | file | mutation | named test | named | whole test package | separate whole source package | verdict |
|---|---|---|---|---|---|---|---|
| `V10-STAMP` | `host/store/credentials.go` | `const credentialsSchemaVersion = 10` → `const credentialsSchemaVersion = 9` | `TestOperatorSchema10` | K (TestOperatorSchema10) | K (TestAppendEpisodeLaunchPreservesMintJSON, TestConcurrentApprovalsMintOnce, TestCredentialsMintAndRevokeBesideALiveWorldWriter, +23 in ledger) | same package | **K** |
| `V10-TABLE-SET` | `host/store/credentials.go` | `case version == 10:` → `case version == 11:` | `TestOperatorSchema10` | K (TestOperatorSchema10) | K (TestAppendEpisodeLaunchPreservesMintJSON, TestConcurrentApprovalsMintOnce, TestCredentialsMintAndRevokeBesideALiveWorldWriter, +15 in ledger) | same package | **K** |
| `V9-SKIP-MIGRATION` | `host/store/credentials.go` | `if version < 10 {` → `if version < 9 {` | `TestOperatorSchema10` | K (TestOperatorSchema10) | K (TestOperatorSchema10) | same package | **K** |
| `V10-FRESH-OMIT` | `host/store/credentials.go` | `tx.Exec(credentialsSchemaSQL + operatorSchemaSQL)` → `tx.Exec(credentialsSchemaSQL)` | `TestOperatorSchema10` | K (TestOperatorSchema10) | K (TestAppendEpisodeLaunchPreservesMintJSON, TestConcurrentApprovalsMintOnce, TestCredentialsMintAndRevokeBesideALiveWorldWriter, +16 in ledger) | same package | **K** |
| `V10-PARTIAL-MIGRATION` | `host/store/credentials.go` | `if _, err := tx.Exec(stmt); err != nil { return` → `if _, err := tx.Exec(stmt); err != nil { _ = tx.Commit() return` | `TestOperatorSchema10` | K (TestOperatorSchema10) | K (TestOperatorSchema10) | same package | **K** |
| `V9-BOUND-FLOOR` | `host/store/credentials.go` | `c.version < CredentialsBindingWriteVersion {` → `c.version < credentialsSchemaVersion {` | `TestOperatorSchema10` | K (TestOperatorSchema10) | K (TestOperatorSchema10) | same package | **K** |
| `FK-CONNECTION-OFF` | `host/store/credentials.go` | `dsn += "&_pragma=foreign_keys(1)"` → `dsn += "&_pragma=foreign_keys(0)"` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-ID-PRIMARY` | `host/store/operator_schema.sql` | `TEXT NOT NULL PRIMARY KEY,` → `TEXT NOT NULL,` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorAdmissionBoundaries, TestOperatorClaimBeforeExpiry, TestOperatorClaimRace, +5 in ledger) | same package | **K** |
| `DDL-ACTION-ID-NOTNULL` | `host/store/operator_schema.sql` | `TEXT NOT NULL PRIMARY KEY,` → `TEXT PRIMARY KEY,` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-EPISODE` | `host/store/operator_schema.sql` | `episode_id TEXT NOT NULL` → `episode_id TEXT` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-KIND-NOTNULL` | `host/store/operator_schema.sql` | `kind TEXT NOT NULL` → `kind TEXT` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-KIND` | `host/store/operator_schema.sql` | `CHECK (kind IN ('start', 'finish'))` → `` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestConvertedStoreMethodsM5a, TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-REVIEW-NOTNULL` | `host/store/operator_schema.sql` | `review BLOB NOT NULL` → `review BLOB` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-REVIEW-LEN` | `host/store/operator_schema.sql` | `CHECK (length(review) = 32)` → `` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-DIGEST-NOTNULL` | `host/store/operator_schema.sql` | `digest BLOB NOT NULL` → `digest BLOB` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-DIGEST-LEN` | `host/store/operator_schema.sql` | `CHECK (length(digest) = 32)` → `` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-CHALLENGE-LEN` | `host/store/operator_schema.sql` | `CHECK (length(challenge) = 32)` → `` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-CREATED` | `host/store/operator_schema.sql` | `created_at INTEGER NOT NULL` → `created_at INTEGER` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-EXPIRES-NOTNULL` | `host/store/operator_schema.sql` | `expires_at INTEGER NOT NULL` → `expires_at INTEGER` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-EXPIRY-ORDER` | `host/store/operator_schema.sql` | `CHECK (expires_at > created_at)` → `CHECK (expires_at >= created_at)` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-STATE-NOTNULL` | `host/store/operator_schema.sql` | `state TEXT NOT NULL` → `state TEXT` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-STATE-ENUM` | `host/store/operator_schema.sql` | `CHECK (state IN ('requested', 'confirmed', 'claimed', 'running', 'awaiting-pr-confirmation',` → `` (full anchor in ledger) | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-GENERATION` | `host/store/operator_schema.sql` | `generation INTEGER NOT NULL` → `generation INTEGER` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-PAYLOAD` | `host/store/operator_schema.sql` | `CHECK (terminal_at IS NOT NULL OR (action IS NOT NULL AND challenge IS NOT NULL)),` → `` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-ARTIFACT` | `host/store/operator_schema.sql` | `CHECK ((approved_by IS NULL) = (artifact IS NULL)),` → `` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-REQUESTED` | `host/store/operator_schema.sql` | `CHECK (state <> 'requested' OR (approved_by IS NULL AND owner IS NULL)),` → `` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-APPROVED` | `host/store/operator_schema.sql` | `CHECK (state NOT IN ('confirmed', 'claimed', 'running', 'awaiting-pr-confirmation', 'uncerta` → `` (full anchor in ledger) | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-OWNER` | `host/store/operator_schema.sql` | `CHECK ((owner IS NOT NULL) = (state IN ('claimed', 'running', 'awaiting-pr-confirmation', 'u` → `` (full anchor in ledger) | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-ACTION-TERMINAL` | `host/store/operator_schema.sql` | `CHECK ((terminal_at IS NOT NULL) = (state IN ('succeeded', 'failed', 'expired', 'cancelled')` → `` (full anchor in ledger) | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-FINISH-LINK-KIND` | `host/store/operator_schema.sql` | `CHECK (v3_id IS NULL OR kind = 'finish'),` → `` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-FINISH-LINK-REQUIRED` | `host/store/operator_schema.sql` | `, CHECK (state <> 'awaiting-pr-confirmation' OR v3_id IS NOT NULL)` → `` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-FINISH-LINK-UNIQUE` | `host/store/operator_schema.sql` | `v3_id TEXT UNIQUE` → `v3_id TEXT` | `TestOperatorClaimRace` | K (TestOperatorClaimRace) | K (TestOperatorClaimRace) | same package | **K** |
| `DDL-EPISODE-UNIQUE` | `host/store/operator_schema.sql` | `CREATE UNIQUE INDEX operator_active_episode` → `CREATE INDEX operator_active_episode` | `TestOperatorAdmissionBoundaries` | K (TestOperatorAdmissionBoundaries) | K (TestOperatorAdmissionBoundaries) | same package | **K** |
| `DDL-RUNNER-SINGLETON` | `host/store/operator_schema.sql` | `CHECK (singleton = 1)` → `` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-RUNNER-PK` | `host/store/operator_schema.sql` | `singleton INTEGER PRIMARY KEY` → `singleton INTEGER` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-RUNNER-GENERATION-NOTNULL` | `host/store/operator_schema.sql` | `generation INTEGER NOT NULL` → `generation INTEGER` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-RUNNER-OWNER-NOTNULL` | `host/store/operator_schema.sql` | `owner TEXT NOT NULL` → `owner TEXT` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-RUNNER-PROTOCOL-NOTNULL` | `host/store/operator_schema.sql` | `protocol INTEGER NOT NULL` → `protocol INTEGER` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-RUNNER-BUILD-NOTNULL` | `host/store/operator_schema.sql` | `build TEXT NOT NULL` → `build TEXT` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-RUNNER-SEEN-NOTNULL` | `host/store/operator_schema.sql` | `seen_at INTEGER NOT NULL` → `seen_at INTEGER` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-RUNNER-READINESS-NOTNULL` | `host/store/operator_schema.sql` | `readiness TEXT NOT NULL` → `readiness TEXT` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-RUNNER-CATALOGUE-GENERATION-NOTNULL` | `host/store/operator_schema.sql` | `catalogue_generation INTEGER NOT NULL` → `catalogue_generation INTEGER` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-STEPS-ACTION-FK` | `host/store/operator_schema.sql` | `REFERENCES operator_actions(id)` → `` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-STEPS-ACTION-NOTNULL` | `host/store/operator_schema.sql` | `action_id TEXT NOT NULL` → `action_id TEXT` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-STEPS-KIND-NOTNULL` | `host/store/operator_schema.sql` | `kind TEXT NOT NULL` → `kind TEXT` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-STEPS-IDENTITY-NOTNULL` | `host/store/operator_schema.sql` | `identity TEXT NOT NULL` → `identity TEXT` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-STEPS-PHASE-NOTNULL` | `host/store/operator_schema.sql` | `phase TEXT NOT NULL` → `phase TEXT` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-STEPS-VERSION-NOTNULL` | `host/store/operator_schema.sql` | `version INTEGER NOT NULL` → `version INTEGER` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-STEPS-CREATED-NOTNULL` | `host/store/operator_schema.sql` | `created_at INTEGER NOT NULL` → `created_at INTEGER` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-STEPS-KIND-ENUM` | `host/store/operator_schema.sql` | `CHECK (kind IN ('operation', 'worktree', 'mint', 'files', 'launch', 'herdr-handoff', 'checks` → `` (full anchor in ledger) | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-STEPS-PHASE` | `host/store/operator_schema.sql` | `CHECK (phase IN ('reserved', 'observed'))` → `` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-STEPS-UNIQUE` | `host/store/operator_schema.sql` | `PRIMARY KEY (kind, identity),` → `` | `TestOperatorDuplicateTap` | K (TestOperatorDuplicateTap) | K (TestOperatorDuplicateTap) | same package | **K** |
| `DDL-STEPS-OP-SHAPE` | `host/store/operator_schema.sql` | `CHECK (kind <> 'operation' OR (phase = 'observed' AND owner IS NULL AND receipt IS NULL)),` → `` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-STEPS-OWNER` | `host/store/operator_schema.sql` | `CHECK (kind = 'operation' OR owner IS NOT NULL),` → `` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `DDL-STEPS-RESERVATION` | `host/store/operator_schema.sql` | `, CHECK (phase <> 'reserved' OR receipt IS NULL)` → `` | `TestOperatorSchema10Constraints` | K (TestOperatorSchema10Constraints) | K (TestOperatorSchema10Constraints) | same package | **K** |
| `REQ-OLD` | `host/store/operator_actions.go` | `if c.version < 10 {` → `if c.version < 9 {` | `TestOperatorSchema10` | K (TestOperatorSchema10) | K (TestOperatorSchema10) | same package | **K** |
| `REQ-DEADLINE` | `host/store/operator_actions.go` | `if err := requireDeadline(ctx); err != nil {` → `if err := requireDeadline(ctx); err != nil && false {` | `TestCredentialsMethodsRequireADeadline` | K (TestCredentialsMethodsRequireADeadline) | K (TestCredentialsMethodsRequireADeadline) | same package | **K** |
| `REQ-OP-ID` | `host/store/operator_actions.go` | `!operatorOperationID.MatchString(request.OperationID)` → `(!operatorOperationID.MatchString(request.OperationID) && false)` | `TestOperatorAdmissionBoundaries` | K (TestOperatorAdmissionBoundaries) | K (TestOperatorAdmissionBoundaries) | same package | **K** |
| `REQ-AT-ZERO` | `host/store/operator_actions.go` | `request.CreatedAt <= 0` → `request.CreatedAt < 0` | `TestOperatorAdmissionBoundaries` | K (TestOperatorAdmissionBoundaries) | K (TestOperatorAdmissionBoundaries) | same package | **K** |
| `REQ-ISSUE-CMP` | `host/store/operator_actions.go` | `now < request.CreatedAt` → `now <= request.CreatedAt` | `TestOperatorAdmissionBoundaries` | K (TestOperatorAdmissionBoundaries) | K (TestOperatorAdmissionBoundaries, TestOperatorClaimBeforeExpiry, TestOperatorClaimRace, +5 in ledger) | same package | **K** |
| `REQ-MAX-CMP` | `host/store/operator_actions.go` | `request.ExpiresAt-request.CreatedAt >` → `request.ExpiresAt-request.CreatedAt >=` | `TestOperatorAdmissionBoundaries` | K (TestOperatorAdmissionBoundaries) | K (TestOperatorAdmissionBoundaries, TestOperatorClaimBeforeExpiry, TestOperatorClaimRace, +5 in ledger) | same package | **K** |
| `REQ-STORE` | `host/store/operator_actions.go` | `a.Store != c.lock.dbPath` → `(a.Store != c.lock.dbPath && false)` | `TestOperatorAdmissionBoundaries` | K (TestOperatorAdmissionBoundaries) | K (TestOperatorAdmissionBoundaries) | same package | **K** |
| `REQ-EXPIRY-CMP` | `host/store/operator_actions.go` | `if now >= request.ExpiresAt` → `if now > request.ExpiresAt` | `TestOperatorAdmissionBoundaries` | K (TestOperatorAdmissionBoundaries) | K (TestOperatorAdmissionBoundaries) | same package | **K** |
| `REQ-REVIEW-HASH` | `host/store/operator_actions.go` | `review, err := approval.OperatorReviewDigest(a)` → `review, err := approval.OperatorSignatureDigest(a)` | `TestOperatorDuplicateTap` | K (TestOperatorDuplicateTap) | K (TestOperatorDuplicateTap, TestOperatorPrune) | same package | **K** |
| `REQ-SIGNATURE-DIGEST-REFUSAL` | `host/store/operator_actions.go` | `digest, err := approval.OperatorSignatureDigest(a) if err != nil {` → `digest, err := approval.OperatorSignatureDigest(a) if err != nil && false {` | `TestOperatorAdmissionBoundaries` | K (TestOperatorAdmissionBoundaries) | K (TestOperatorAdmissionBoundaries) | same package | **K** |
| `WRITE-EXPIRY-STATE` | `host/store/operator_actions.go` | `WHERE state IN ('requested', 'confirmed') AND expires_at <= ?` → `WHERE expires_at <= ?` | `TestOperatorPrune` | K (TestOperatorPrune) | K (TestOperatorPrune) | same package | **K** |
| `REQ-OP-MISMATCH` | `host/store/operator_actions.go` | `if !bytes.Equal(row.Review, review[:]) {` → `if !bytes.Equal(row.Review, review[:]) && false {` | `TestOperatorAdmissionBoundaries` | K (TestOperatorAdmissionBoundaries) | K (TestOperatorAdmissionBoundaries, TestOperatorDuplicateTap) | same package | **K** |
| `REQ-OP-LOOKUP` | `host/store/operator_actions.go` | `kind = 'operation' AND identity = ?)` → `kind = 'operation' AND identity = ? AND 0)` | `TestOperatorReconnect` | K (TestOperatorReconnect) | K (TestOperatorAdmissionBoundaries, TestOperatorDuplicateTap, TestOperatorPrune, +1 in ledger) | same package | **K** |
| `REQ-REVIEW-LOOKUP` | `host/store/operator_actions.go` | `AND state NOT IN ('expired', 'cancelled')` → `AND state NOT IN ('expired', 'cancelled') AND 0` | `TestOperatorDuplicateTap` | K (TestOperatorDuplicateTap) | K (TestOperatorDuplicateTap) | same package | **K** |
| `REQ-REVIEW-EXPIRED` | `host/store/operator_actions.go` | `AND state NOT IN ('expired', 'cancelled')` → `` | `TestOperatorPrune` | K (TestOperatorPrune) | K (TestOperatorPrune) | same package | **K** |
| `REQ-RESERVE-OP` | `host/store/operator_actions.go` | `request.OperationID, row.ID, now);` → `request.OperationID + "x", row.ID, now);` | `TestOperatorDuplicateTap` | K (TestOperatorDuplicateTap) | K (TestOperatorAdmissionBoundaries, TestOperatorDuplicateTap, TestOperatorPrune, +1 in ledger) | same package | **K** |
| `REQ-ACTIVE` | `host/store/operator_actions.go` | `case active != 0:` → `case active != 0 && false:` | `TestOperatorAdmissionBoundaries` | K (TestOperatorAdmissionBoundaries) | K (TestOperatorAdmissionBoundaries, TestOperatorTerminalExclusion) | same package | **K** |
| `REQ-CAP-CMP` | `host/store/operator_actions.go` | `case pending >= 16:` → `case pending > 16:` | `TestOperatorAdmissionBoundaries` | K (TestOperatorAdmissionBoundaries) | K (TestOperatorAdmissionBoundaries) | same package | **K** |
| `READ-EXPIRY-CMP` | `host/store/operator_actions.go` | `row.ExpiresAt <= now` → `row.ExpiresAt < now` | `TestOperatorPrune` | K (TestOperatorPrune) | K (TestOperatorPrune) | same package | **K** |
| `WRITE-EXPIRY-CMP` | `host/store/operator_actions.go` | `AND expires_at <= ?` → `AND expires_at < ?` | `TestOperatorAdmissionBoundaries` | K (TestOperatorAdmissionBoundaries) | K (TestOperatorAdmissionBoundaries) | same package | **K** |
| `PRUNE-AGE-CMP` | `host/store/operator_actions.go` | `WHERE terminal_at < ?` → `WHERE terminal_at <= ?` | `TestOperatorPrune` | K (TestOperatorPrune) | K (TestOperatorPrune) | same package | **K** |
| `PRUNE-RESERVATION` | `host/store/operator_actions.go` | `AND phase = 'reserved')` → `AND phase = 'reserved' AND 0)` | `TestOperatorPrune` | K (TestOperatorPrune) | K (TestOperatorPrune) | same package | **K** |
| `PRUNE-UNCERTAIN` | `host/store/operator_actions.go` | `WHERE terminal_at < ? AND NOT EXISTS` → `WHERE state = 'uncertain' OR terminal_at < ? AND NOT EXISTS` | `TestOperatorPrune` | K (TestOperatorPrune) | K (TestOperatorPrune) | same package | **K** |
| `PRUNE-TOMBSTONE` | `host/store/operator_actions.go` | `SET action = NULL, challenge = NULL, artifact = NULL, approved_by = NULL` → `SET action = NULL, challenge = NULL, artifact = NULL, approved_by = NULL, review = zeroblob(` (full anchor in ledger) | `TestOperatorPrune` | K (TestOperatorPrune) | K (TestOperatorPrune) | same package | **K** |
| `LOCK-EXCLUSIVE-OFF` | `host/store/writer_lock_unix.go` | `syscall.LOCK_EX\|syscall.LOCK_NB` → `syscall.LOCK_SH\|syscall.LOCK_NB` | `TestOperatorTerminalExclusion` | K (TestOperatorTerminalExclusion) | K (TestCloseKeepsWriterLockUntilDurableWorkerSettles, TestCredentialsMintAndRevokeBesideALiveWorldWriter, TestMigrateCredentialsRefusesAHeldWorldLock, +6 in ledger) | same package | **K** |
| `LOCK-NONWAITING-OFF` | `host/store/writer_lock_unix.go` | `syscall.LOCK_EX\|syscall.LOCK_NB` → `syscall.LOCK_EX` | `TestOperatorTerminalExclusion` | K (TestOperatorTerminalExclusion) | K (TestMigrateCredentialsRefusesAHeldWorldLock) | same package | **K** |
| `CONFIRM-READ-ID` | `host/store/operator_cas.go` | `FROM operator_actions WHERE id = ?, cas.ID)) if err != nil` → `FROM operator_actions WHERE ? IS NOT NULL, cas.ID)) if err != nil` | `TestOperatorConfirmStoreBinding` | K (TestOperatorConfirmStoreBinding) | K (TestOperatorConfirmStoreBinding) | S | **K** |
| `CONFIRM-DIGEST` | `host/store/operator_cas.go` | `if !bytes.Equal(row.Digest, cas.Digest[:]) {` → `if !bytes.Equal(row.Digest, cas.Digest[:]) && false {` | `TestOperatorConfirmStoreBinding` | K (TestOperatorConfirmStoreBinding) | K (TestOperatorConfirmStoreBinding) | S | **K** |
| `CONFIRM-GENERATION` | `host/store/operator_cas.go` | `row.State != "requested" \|\| row.Generation != cas.Generation` → `row.State != "requested"` | `TestOperatorConfirmStoreBinding` | K (TestOperatorConfirmStoreBinding) | K (TestOperatorConfirmStoreBinding) | S | **K** |
| `CONFIRM-REPLAY` | `host/store/operator_cas.go` | `if row.ApprovedBy != "" {` → `if row.ApprovedBy != "" && false {` | `TestOperatorConfirmReplay` | K (TestOperatorConfirmReplay) | K (TestOperatorConfirmReplay) | S | **K** |
| `CONFIRM-EXPIRY-CMP` | `host/store/operator_cas.go` | `if now >= row.ExpiresAt { return OperatorActionRow{}, ErrOperatorExpired` → `if now > row.ExpiresAt { return OperatorActionRow{}, ErrOperatorExpired` | `TestOperatorConfirmStoreBinding` | K (TestOperatorConfirmStoreBinding) | K (TestOperatorConfirmStoreBinding) | S | **K** |
| `CONFIRM-CREATED-CMP` | `host/store/operator_cas.go` | `if now < row.CreatedAt \|\|` → `if now <= row.CreatedAt \|\|` | `TestOperatorConfirmStoreBinding` | K (TestOperatorConfirmStoreBinding) | K (TestOperatorConfirmStoreBinding) | S | **K** |
| `CONFIRM-SQL-ID` | `host/store/operator_cas.go` | `WHERE id = ? AND digest = ? AND generation = ? AND state = 'requested'` → `WHERE ? IS NOT NULL AND digest = ? AND generation = ? AND state = 'requested'` | `TestOperatorConfirmCrash` | S | S | S | **E** |
| `CONFIRM-SQL-DIGEST` | `host/store/operator_cas.go` | `AND digest = ? AND generation = ? AND state = 'requested'` → `AND ? IS NOT NULL AND generation = ? AND state = 'requested'` | `TestOperatorConfirmCrash` | K (TestOperatorConfirmCrash) | K (TestOperatorConfirmCrash) | S | **K** |
| `CONFIRM-SQL-GEN` | `host/store/operator_cas.go` | `AND generation = ? AND state = 'requested'` → `AND ? IS NOT NULL AND state = 'requested'` | `TestOperatorConfirmCrash` | K (TestOperatorConfirmCrash) | K (TestOperatorConfirmCrash) | S | **K** |
| `CONFIRM-SQL-STATE` | `host/store/operator_cas.go` | `AND state = 'requested' AND expires_at > ?` → `AND expires_at > ?` | `TestOperatorConfirmCrash` | S | S | S | **E** |
| `CONFIRM-SQL-EXPIRY-CMP` | `host/store/operator_cas.go` | `AND state = 'requested' AND expires_at > ?` → `AND state = 'requested' AND expires_at >= ?` | `TestOperatorConfirmCrash` | K (TestOperatorConfirmCrash) | K (TestOperatorConfirmCrash) | S | **K** |
| `CONFIRM-KEY-USE` | `host/store/operator_cas.go` | `if err := useApprover(ctx, tx, ApproverUse{Name: binding.Name, Format: binding.Format, Publi` → `if err := useApprover(ctx, tx, ApproverUse{Name: binding.Name, Format: binding.Format, Publi` (full anchor in ledger) | `TestOperatorConfirmCrash` | K (TestOperatorConfirmCrash) | K (TestOperatorConfirmCrash) | S | **K** |
| `CONFIRM-SPLIT` | `host/store/operator_cas.go` | `res, err := tx.ExecContext(ctx, UPDATE operator_actions SET state = 'confirmed'` → `if err := tx.Commit(); err != nil { return OperatorActionRow{}, err } res, err := c.db.ExecC` (full anchor in ledger) | `TestOperatorConfirmCrash` | K (TestOperatorConfirmCrash) | K (TestOperatorConfirmCrash, TestOperatorConfirmReplay, TestOperatorConfirmStoreBinding, +2 in ledger) | S | **K** |
| `CONFIRM-SWALLOW` | `host/store/operator_cas.go` | `if err := operatorChangedOne(res, err); err != nil { return OperatorActionRow{}, err` → `if err := operatorChangedOne(res, err); err != nil { return row, nil` | `TestOperatorConfirmCrash` | K (TestOperatorConfirmCrash) | K (TestOperatorConfirmCrash) | S | **K** |
| `CONFIRM-IDENTITY` | `host/store/operator_cas.go` | `signed.ApprovedBy(), artifact, cas.ID` → `binding.Name, artifact, cas.ID` | `TestOperatorConfirmReplay` | K (TestOperatorConfirmReplay) | K (TestOperatorConfirmReplay, TestOperatorRunnerRechecks) | S | **K** |
| `CONFIRM-RP` | `host/store/operator_cas.go` | `rp != binding.RPID \|\|` → `(rp != binding.RPID && false) \|\|` | `TestOperatorConfirmCrash` | K (TestOperatorConfirmCrash) | K (TestOperatorConfirmCrash) | S | **K** |
| `CONFIRM-ORIGIN` | `host/store/operator_cas.go` | `origin != binding.Origin \|\|` → `(origin != binding.Origin && false) \|\|` | `TestOperatorConfirmCrash` | K (TestOperatorConfirmCrash) | K (TestOperatorConfirmCrash) | S | **K** |
| `CONFIRM-BACKUP` | `host/store/operator_cas.go` | `backup != binding.BackupEligible {` → `(backup != binding.BackupEligible && false) {` | `TestOperatorConfirmCrash` | K (TestOperatorConfirmCrash) | K (TestOperatorConfirmCrash) | S | **K** |
| `CONFIRM-PASSKEY-SKIP` | `host/store/operator_cas.go` | `signed.PasskeyUse(); ok {` → `signed.PasskeyUse(); ok && false {` | `TestOperatorConfirmReplay` | K (TestOperatorConfirmReplay) | K (TestOperatorConfirmCrash, TestOperatorConfirmReplay, TestOperatorRunnerRechecks) | S | **K** |
| `CONFIRM-COUNTER-ERR` | `host/store/operator_cas.go` | `SignCount: int64(count)}); err != nil {` → `SignCount: int64(count)}); err != nil && false {` | `TestOperatorConfirmCrash` | K (TestOperatorConfirmCrash) | K (TestOperatorConfirmCrash) | S | **K** |
| `CONFIRM-BINDS` | `host/store/operator_cas.go` | `signed.Binds(a, [32]byte(row.Challenge), row.ExpiresAt); err != nil {` → `signed.Binds(a, [32]byte(row.Challenge), row.ExpiresAt); err != nil && false {` | `TestOperatorConfirmStoreBinding` | K (TestOperatorConfirmStoreBinding) | K (TestOperatorConfirmStoreBinding) | S | **K** |
| `CLAIM-READ-ID` | `host/store/operator_cas.go` | `FROM operator_actions WHERE id = ?, cas.ID)) if err == nil` → `FROM operator_actions WHERE ? IS NOT NULL, cas.ID)) if err == nil` | `TestOperatorClaimBeforeExpiry` | K (TestOperatorClaimBeforeExpiry) | K (TestOperatorClaimBeforeExpiry) | same package | **K** |
| `CLAIM-READ-DIGEST` | `host/store/operator_cas.go` | `!bytes.Equal(row.Digest, cas.Digest[:]) \|\| row.Generation` → `(!bytes.Equal(row.Digest, cas.Digest[:]) && false) \|\| row.Generation` | `TestOperatorClaimBeforeExpiry` | K (TestOperatorClaimBeforeExpiry) | K (TestOperatorClaimBeforeExpiry) | same package | **K** |
| `CLAIM-READ-GEN` | `host/store/operator_cas.go` | `\|\| row.Generation != cas.Generation) {` → `\|\| (row.Generation != cas.Generation && false)) {` | `TestOperatorClaimBeforeExpiry` | K (TestOperatorClaimBeforeExpiry) | K (TestOperatorClaimBeforeExpiry) | same package | **K** |
| `CLAIM-PRE-STATE` | `host/store/operator_cas.go` | `row.State != "confirmed" \|\|` → `(row.State != "confirmed" && false) \|\|` | `TestOperatorClaimBeforeExpiry` | K (TestOperatorClaimBeforeExpiry) | K (TestOperatorClaimBeforeExpiry) | same package | **K** |
| `CLAIM-PRE-OWNER` | `host/store/operator_cas.go` | `\|\| row.Owner != "" \|\|` → `\|\| (row.Owner != "" && false) \|\|` | `TestOperatorClaimBeforeExpiry` | S | S | same package | **E** |
| `CLAIM-ISSUE-CMP` | `host/store/operator_cas.go` | `row.Owner != "" \|\| now < row.CreatedAt` → `row.Owner != "" \|\| now <= row.CreatedAt` | `TestOperatorClaimRace` | K (TestOperatorClaimRace) | K (TestOperatorClaimRace) | same package | **K** |
| `CLAIM-EXPIRY-CMP` | `host/store/operator_cas.go` | `if now >= row.ExpiresAt { return ErrOperatorExpired` → `if now > row.ExpiresAt { return ErrOperatorExpired` | `TestOperatorClaimBeforeExpiry` | K (TestOperatorClaimBeforeExpiry) | K (TestOperatorClaimBeforeExpiry) | same package | **K** |
| `CLAIM-RUNNER-GEN` | `host/store/operator_cas.go` | `generation = ? AND owner = ? AND seen_at` → `? IS NOT NULL AND owner = ? AND seen_at` | `TestOperatorClaimBeforeExpiry` | K (TestOperatorClaimBeforeExpiry) | K (TestOperatorClaimBeforeExpiry) | same package | **K** |
| `CLAIM-RUNNER-OWNER` | `host/store/operator_cas.go` | `owner = ? AND seen_at` → `? IS NOT NULL AND seen_at` | `TestOperatorClaimBeforeExpiry` | K (TestOperatorClaimBeforeExpiry) | K (TestOperatorClaimBeforeExpiry) | same package | **K** |
| `CLAIM-SEEN-STALE-CMP` | `host/store/operator_cas.go` | `seen_at > ?` → `seen_at >= ?` | `TestOperatorClaimBeforeExpiry` | K (TestOperatorClaimBeforeExpiry) | K (TestOperatorClaimBeforeExpiry) | same package | **K** |
| `CLAIM-SEEN-FUTURE-CMP` | `host/store/operator_cas.go` | `seen_at <= ?` → `seen_at < ?` | `TestOperatorClaimRace` | K (TestOperatorClaimRace) | K (TestOperatorClaimRace) | same package | **K** |
| `CLAIM-READINESS` | `host/store/operator_cas.go` | `(readiness = 'both' OR readiness = ?)` → `(? IS NOT NULL)` | `TestOperatorClaimBeforeExpiry` | K (TestOperatorClaimBeforeExpiry) | K (TestOperatorClaimBeforeExpiry) | same package | **K** |
| `CLAIM-RUNNER-IGNORED` | `host/store/operator_cas.go` | `if !current {` → `if !current && false {` | `TestOperatorClaimBeforeExpiry` | K (TestOperatorClaimBeforeExpiry) | K (TestOperatorClaimBeforeExpiry) | same package | **K** |
| `CLAIM-SIBLING` | `host/store/operator_cas.go` | `AND b.owner IS NOT NULL)` → `AND b.owner IS NOT NULL AND 0)` | `TestOperatorClaimRace` | K (TestOperatorClaimRace) | K (TestOperatorClaimRace) | same package | **K** |
| `CLAIM-SQL-ID` | `host/store/operator_cas.go` | `WHERE id = ? AND digest = ? AND generation = ? AND state = 'confirmed'` → `WHERE ? IS NOT NULL AND digest = ? AND generation = ? AND state = 'confirmed'` | `TestOperatorClaimRace` | S | S | same package | **E** |
| `CLAIM-SQL-DIGEST` | `host/store/operator_cas.go` | `AND digest = ? AND generation = ? AND state = 'confirmed'` → `AND ? IS NOT NULL AND generation = ? AND state = 'confirmed'` | `TestOperatorClaimBeforeExpiry` | S | S | same package | **E** |
| `CLAIM-SQL-GEN` | `host/store/operator_cas.go` | `AND generation = ? AND state = 'confirmed'` → `AND ? IS NOT NULL AND state = 'confirmed'` | `TestOperatorClaimBeforeExpiry` | S | S | same package | **E** |
| `CLAIM-SQL-STATE` | `host/store/operator_cas.go` | `AND state = 'confirmed' AND expires_at > ?` → `AND expires_at > ?` | `TestOperatorClaimBeforeExpiry` | S | S | same package | **E** |
| `CLAIM-SQL-EXPIRY-CMP` | `host/store/operator_cas.go` | `AND state = 'confirmed' AND expires_at > ?` → `AND state = 'confirmed' AND expires_at >= ?` | `TestOperatorClaimBeforeExpiry` | S | S | same package | **E** |
| `CLAIM-SQL-OWNER` | `host/store/operator_cas.go` | `AND owner IS NULL AND NOT EXISTS` → `AND NOT EXISTS` | `TestOperatorClaimBeforeExpiry` | S | S | same package | **E** |
| `CLAIM-AFFECTED` | `host/store/operator_cas.go` | `if n != 1 {` → `if n != 1 && false {` | `TestOperatorClaimRace` | K (TestOperatorClaimRace) | K (TestOperatorClaimRace) | same package | **K** |
| `CLAIM-ADVANCE` | `host/store/operator_cas.go` | `generation = generation + 1` → `generation = generation` | `TestOperatorClaimRace` | K (TestOperatorClaimRace) | K (TestOperatorClaimRace) | same package | **K** |
| `CONFIRM-ARTIFACT-HASH` | `host/approval/operator_signed.go` | `\|\| sha256.Sum256(artifact) != s.artifactHash` → `` | `TestOperatorConfirmStoreBinding` | K (TestOperatorConfirmStoreBinding) | K (TestOperatorConfirmStoreBinding) | S | **K** |
| `CONFIRM-ROLE` | `host/approval/operator_signed.go` | `if confirm {` → `if true {` | `TestOperatorRunnerRechecks` | K (TestOperatorRunnerRechecks) | K (TestOperatorRunnerRechecks) | S | **K** |
| `AUTH-TRANSPORT-NIL` | `host/authority/judge_operator.go` | `return v.signed }` → `return nil }` | `TestOperatorConfirmReplay` | K (TestOperatorConfirmReplay) | K (TestOperatorConfirmCrash, TestOperatorConfirmReadOnly, TestOperatorConfirmReplay, +3 in ledger) | same package | **K** |
| `RECHECK-ROW-ACTION` | `host/authority/judge_operator.go` | `return RecheckOperatorAction(ctx, ro, frozen, challenge, expiry, now, row.Artifact, row.Appr` → `frozen, _ = approval.DecodeOperatorAction(row.Action) return RecheckOperatorAction(ctx, ro, ` (full anchor in ledger) | `TestOperatorRunnerRechecks` | K (TestOperatorRunnerRechecks) | K (TestOperatorRunnerRechecks) | same package | **K** |
| `RECHECK-ROW-EXPIRY` | `host/authority/judge_operator.go` | `ro, frozen, challenge, expiry, now, row.Artifact` → `ro, frozen, challenge, row.ExpiresAt, now, row.Artifact` | `TestOperatorRunnerRechecks` | K (TestOperatorRunnerRechecks) | K (TestOperatorRunnerRechecks) | same package | **K** |
| `RECHECK-ROW-CHALLENGE` | `host/authority/judge_operator.go` | `ro, frozen, challenge, expiry, now, row.Artifact` → `ro, frozen, [32]byte{}, expiry, now, row.Artifact` | `TestOperatorRunnerRechecks` | K (TestOperatorRunnerRechecks) | K (TestOperatorRunnerRechecks) | same package | **K** |
| `RECHECK-ROW-APPROVER` | `host/authority/judge_operator.go` | `if s.ApprovedBy() != approvedBy {` → `if s.ApprovedBy() != approvedBy && false {` | `TestOperatorRunnerRechecks` | K (TestOperatorRunnerRechecks) | K (TestOperatorRunnerRechecks) | same package | **K** |
| `RECHECK-STATUS` | `host/authority/judge_operator.go` | `return RecheckOperatorAction(ctx, ro, frozen, challenge, expiry, now, row.Artifact, row.Appr` → `if row.State == "claimed" { return verifiedOperatorAction(nil, true), nil }; return RecheckO` (full anchor in ledger) | `TestOperatorRunnerRechecks` | K (TestOperatorRunnerRechecks) | K (TestOperatorRunnerRechecks, TestVerifiedOperatorActionOnlyFromTheJudges) | same package | **K** |
| `LOG-IN-JUDGE-NOFAIL` | `host/authority/judge_operator.go` | `(*VerifiedOperatorAction, error) { format, a, err := operatorApprover(ctx, ro, artifact)` → `(*VerifiedOperatorAction, error) { println("x") format, a, err := operatorApprover(ctx, ro, ` (full anchor in ledger) | `TestOperatorJudgesWriteNothing` | K (TestOperatorJudgesWriteNothing) | K (TestOperatorJudgesWriteNothing) | same package | **K** |
| `LOG-IN-READ` | `host/store/operator_actions.go` | `OperatorAction(ctx context.Context, id string, now int64) (OperatorActionRow, error) { if er` → `OperatorAction(ctx context.Context, id string, now int64) (OperatorActionRow, error) { print` (full anchor in ledger) | `TestOperatorStoreWritesNothing` | K (TestOperatorStoreWritesNothing) | K (TestOperatorStoreWritesNothing) | S | **K** |
| `LOG-IN-REQUEST` | `host/store/operator_actions.go` | `RequestOperatorAction(ctx context.Context, request OperatorRequest, now int64) (OperatorActi` → `RequestOperatorAction(ctx context.Context, request OperatorRequest, now int64) (OperatorActi` (full anchor in ledger) | `TestOperatorStoreWritesNothing` | K (TestOperatorStoreWritesNothing) | K (TestOperatorStoreWritesNothing) | S | **K** |
| `LOG-IN-CONFIRM` | `host/store/operator_cas.go` | `proof operatorConfirmation, now int64) (OperatorActionRow, error) { if err := c.operatorAcce` → `proof operatorConfirmation, now int64) (OperatorActionRow, error) { println("x") if err := c` (full anchor in ledger) | `TestOperatorStoreWritesNothing` | K (TestOperatorStoreWritesNothing) | K (TestOperatorStoreWritesNothing) | S | **K** |
| `LOG-IN-CLAIM` | `host/store/operator_cas.go` | `now int64) error { if err := c.operatorAccess(ctx); err != nil {` → `now int64) error { println("x") if err := c.operatorAccess(ctx); err != nil {` | `TestOperatorStoreWritesNothing` | K (TestOperatorStoreWritesNothing) | K (TestOperatorStoreWritesNothing) | S | **K** |
| `LOG-IN-PRUNE` | `host/store/operator_actions.go` | `now int64) error { if err := c.operatorAccess(ctx); err != nil {` → `now int64) error { println("x") if err := c.operatorAccess(ctx); err != nil {` | `TestOperatorStoreWritesNothing` | K (TestOperatorStoreWritesNothing) | K (TestOperatorStoreWritesNothing) | S | **K** |

**Named equivalents**

- `CONFIRM-SQL-ID`. A valid signature digest contains the action id, so distinct valid rows cannot match this digest. The mandatory read-by-id and affected-row check remain.
- `CONFIRM-SQL-STATE`. The same transaction already refused non-requested state. The production statements between read and write only touch enrollment and counter rows. The CAS check retains a second fence.
- `CLAIM-PRE-OWNER`. The v10 CHECK makes a non-NULL owner invalid on a confirmed row. The SQL owner predicate still refuses a malformed owned row.
- `CLAIM-SQL-ID`. A valid signature digest uniquely identifies the signed action id. The read-by-id, digest check and affected-row fence remain.
- `CLAIM-SQL-STATE`. The same transaction previously read and required confirmed. No statement writes the action between read and update.
- `CLAIM-SQL-GEN`. The exact phase generation was read and compared in the same transaction before the update.
- `CLAIM-SQL-OWNER`. The same transaction required an empty owner before the update. The schema also requires NULL on confirmed rows.
- `CLAIM-SQL-EXPIRY-CMP`. The precheck rejects now at or after expiry with the same injected time in the same transaction. The SQL comparison stays strict.
- `CLAIM-SQL-DIGEST`. The same transaction already compared the exact digest. No write can change that row between its own read and update.

## 5. Files

````text
M	cmd/ailang-worldd/approve_await_refusals_test.go
M	cmd/ailang-worldd/approve_passkey_refusals_test.go
M	cmd/ailang-worldd/approve_passkey_test.go
M	cmd/ailang-worldd/approve_sig_test.go
M	cmd/ailang-worldd/approver_label_test.go
M	cmd/ailang-worldd/approver_passkey_test.go
M	cmd/ailang-worldd/approver_rig_test.go
M	cmd/ailang-worldd/episode_finish_surface_test.go
M	cmd/ailang-worldd/episode_preset.go
M	cmd/ailang-worldd/session_new.go
M	host/approval/operator_signed.go
M	host/authority/judge_operator.go
M	host/authority/judge_operator_silent_test.go
M	host/authority/judge_operator_test.go
M	host/authority/judge_test.go
A	host/authority/operator_store_test.go
M	host/daemon/surface_approve_rig_test.go
M	host/store/approvers_test.go
M	host/store/credentials.go
M	host/store/credentials_test.go
M	host/store/mint_requests_test.go
A	host/store/operator_actions.go
A	host/store/operator_actions_test.go
A	host/store/operator_cas.go
A	host/store/operator_cas_test.go
A	host/store/operator_schema.sql
A	host/store/operator_schema_test.go
M	host/store/verb_requests_test.go
````


Pre-existing test files changed (zero deletions): `host/store/credentials_test.go` pins the v10 table/version set and the request context guard, `host/store/verb_requests_test.go` retains the v8-through-v9 path and advances the newer-version refusal, `host/authority/judge_operator_test.go` adds the real claimed-recheck branch and `host/authority/judge_operator_silent_test.go` captures fd 1 and 2 plus the store entry points. The remaining changes preserve legacy fixture tests:
- `host/store/approvers_test.go` and `host/store/mint_requests_test.go` add the three new table names to independent expected upgrade table sets.
- `host/authority/judge_test.go` and `host/daemon/surface_approve_rig_test.go` strip v10 tables when the fixture steps back through v9. They retain all actual row and protocol assertions.
- `cmd/ailang-worldd/approve_await_refusals_test.go`, `approve_passkey_refusals_test.go`, `approve_passkey_test.go`, `approve_sig_test.go`, `approver_label_test.go`, `approver_passkey_test.go`, `approver_rig_test.go`, `episode_finish_surface_test.go` strip the same three v10 tables from v4/v5/v6/v8 fixtures. The three v10 table drops run before the older tables and preserve their historical version fences.
Existing production files are only credentials.go, operator_signed.go, judge_operator.go and the two compatibility floor call sites. New DDL lives beside the credentials schema in the same package. No collision files from #641 or #617 are touched.

## 6. Sprint gates

Run one bash shell with `bash "$HOME/.ailang/state/world-iter428/gates.sh" "$W" "$HOME/.ailang/state/world-iter428/gates" final "$(git -C "$W" rev-parse HEAD)"`. Each leg tees to its own log, records its rc and continues for complete evidence. `ALL_DONE.<sha>` contains the final failures. The controller reruns every sandbox-uninformative leg on the real checkout. Inspect only the exact bind failures when classifying sandbox evidence.

````bash gates.sh
#!/bin/bash
# Usage: gates.sh <tree> <logdir> <base|final|scratch> [sha]
# scratch uses a measured path list and a separate read-only git reference.
T=$1; L=$2; MODE=$3
mkdir -p "$L"
cd "$T" || exit 2
B=92e1c44ebd8dbbdb2b962ef54aed0dd9a0bf2456
REF=${GATES_GIT_REF:-$T}
SHA=${4:-$(git -C "$REF" rev-parse HEAD)}
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export AILANG_STATE_DIR="$L/ailang-state"
mkdir -p "$AILANG_STATE_DIR"
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR WORLD_UPDATE_OPERATOR_ACTION_GOLDEN WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN
: > "$L/legs.txt"
FAILS=
s0=$(date +%s)
line() { echo "$*" | tee -a "$L/legs.txt"; }
leg() {
  local name=$1 cmd=$2 s rc
  s=$(date +%s)
  bash -c "$cmd" 2>&1 | tee "$L/$name.txt"
  rc=${PIPESTATUS[0]}
  if [ "$rc" -ne 0 ] && [ "${GATES_SANDBOX:-0}" = 1 ] &&
      grep -Ei 'listen .*bind: operation not permitted' "$L/$name.txt" >/dev/null; then
    line "LEG $name rc=$rc $(( $(date +%s)-s ))s UNINFORMATIVE UNDER SANDBOX (loopback bind; inspect for independent failures)"
    FAILS="$FAILS ${name}_sandbox_uninformative"
  else
    line "LEG $name rc=$rc $(( $(date +%s)-s ))s"
    [ "$rc" -eq 0 ] || FAILS="$FAILS $name"
  fi
}
paths() {
  if [ "$MODE" = scratch ]; then cat "$GATES_CHANGED_FILE"
  else git -C "$REF" diff --name-only "$B" HEAD
  fi
}
TOUCHED=$(paths | awk -F/ '/\.go$/ && ($1=="cmd" || $1=="host") {NF--; OFS="/"; $1=$1; print}' | sort -u | tr '\n' ' ')
EXPECTED='cmd/ailang-worldd host/approval host/authority host/daemon host/store '
line "TOUCHED_PACKAGES: $TOUCHED"
if [ "$MODE" != base ] && [ "$TOUCHED" != "$EXPECTED" ]; then FAILS="$FAILS touched_packages_differ"; fi
leg status "git -C '$REF' status --porcelain | wc -l"
leg gofmt 'out=$(gofmt -l cmd host); echo "$out"; test -z "$out"'
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./..."
leg verify_ail './scripts/verify_ail.sh'
PKGS="host/approval host/approval/approvaltest host/authority host/store cmd/ailang-worldd host/daemon"
for p in $PKGS; do
  leg "pkg-${p//\//-}" "go test -count=1 -timeout 8m ./$p"
  leg "race-${p//\//-}" "/usr/bin/time -p go test -count=1 -race -timeout 8m ./$p"
done
inv() {
  local p=$1 n i
  shift
  n=$(echo "$*" | tr ' ' '|')
  i="inv-${p//\//-}"
  leg "$i" "go test -v -count=1 ./$p -run '^($n)\$'"
  local t
  for t in "$@"; do
    grep -qE "^--- PASS: $t " "$L/$i.txt" ||
      { line "NOT PASSED ($i): $t"; FAILS="$FAILS ${i}_missing"; }
  done
}
inv host/approval TestApprovalImportsOnlyTheStandardLibrary TestSSHSIGCrossCheckWithSSHKeygen TestVerbDomainSeparation TestOperatorConfirmBothKeys TestOperatorRunnerRecheckLaw TestOperatorPureEntryPointsWriteNothing
inv host/approval/approvaltest TestApprovalTestIsTestOnly TestTestSignerAgreesWithSSHKeygen
AUTH_TESTS=""
if [ "$MODE" != base ]; then AUTH_TESTS="TestOperatorConfirmReplay TestOperatorConfirmCrash TestOperatorConfirmStoreBinding TestOperatorConfirmReadOnly TestOperatorStoreWritesNothing"; fi
inv host/authority TestJudgeVerbSSH TestJudgeVerbAssertion TestRecheckVerb TestReadOnlyAtFailuresAreErrors TestVerifiedOperatorActionOnlyFromTheJudges TestOperatorRunnerRechecks TestOperatorJudgesWriteNothing $AUTH_TESTS
inv host/boundary TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard
inv host/daemon TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestConfirmRoutesRegistered TestDaemonCredentialWriterSites TestSurfaceVerbRouteRegistered TestSurfaceVerbIsReadOnly TestSurfaceVerbOutcomes TestApprovePOSTRefusalArms TestConfirmPOSTRefusesBeforeAnyWrite TestRequestsJSONConfirmations TestRequestsJSONConfirmationsGolden TestDecidePOSTSSHSIGRefusals
inv cmd/ailang-worldd TestEveryVerbHasHelp TestVerbTableMatchesTheCLISwitch TestCLICredentialWriterSites TestOperatorFlagParity TestApproveAwaitNeedsSchema7AndItsRequest TestApprovePasskeyRefusalsBeforeAnyEffect TestApproveBySignatureRefusalsBeforeAnyEffect TestPutFinishVerbRequestRefusals
STORE_TESTS='TestProductionContextRoots TestProductionGoSurface TestStoreSurfaceTakesContext TestCredentialsMethodsRequireADeadline TestCredentialsSchemaMatchesCanonicalDDL TestCredentialsSchema9VerbRequests TestCredentialsUpgradeToV4 TestCredentialsUpgradeToV6 TestWriterLockCrossProcessContention TestWriterLockHeldIsNonMutating'
if [ "$MODE" = base ]; then inv host/store $STORE_TESTS TestCredentialsNewerThanNineIsRefused
else inv host/store $STORE_TESTS TestCredentialsNewerThanTenIsRefused TestOperatorSchema10 TestOperatorSchema10Constraints TestOperatorAdmissionBoundaries TestOperatorDuplicateTap TestOperatorReconnect TestOperatorClaimBeforeExpiry TestOperatorClaimRace TestOperatorTerminalExclusion TestOperatorPrune
fi
ROUTES=$(rg -c 'mux\.(HandleFunc|Handle)\(' host/daemon/daemon.go)
POSTS=$(rg -c 'mux\.HandleFunc\("POST /surface/' host/daemon/daemon.go)
BROUTES=$(git -C "$REF" show "$B:host/daemon/daemon.go" | rg -c 'mux\.(HandleFunc|Handle)\(')
line "DAEMON_ROUTES: $ROUTES (base $BROUTES) SURFACE_POSTS: $POSTS"
[ "$ROUTES" = "$BROUTES" ] && [ "$POSTS" = 3 ] || FAILS="$FAILS routes_changed"
ALLOWED='^(cmd/ailang-worldd/(episode_preset|session_new|approve_await_refusals_test|approve_passkey_refusals_test|approve_passkey_test|approve_sig_test|approver_label_test|approver_passkey_test|approver_rig_test|episode_finish_surface_test)\.go|host/approval/operator_signed\.go|host/authority/(judge_operator\.go|judge_operator(_silent)?_test\.go|judge_test\.go|operator_store_test\.go)|host/daemon/surface_approve_rig_test\.go|host/store/((approvers|mint_requests)_test\.go|credentials\.go|credentials_test\.go|operator_(actions(_test)?\.go|cas(_test)?\.go|schema\.sql|schema_test\.go)|verb_requests_test\.go)|design_docs/planned/(w-world-address-bar-verbs-v4s3-sprint-plan\.md|sprint_w-world-address-bar-verbs-v4s3\.json))$'
NOPROD=$(paths | rg -v "$ALLOWED")
line "UNTOUCHED_VIOLATIONS: ${NOPROD:-none}"
[ -z "$NOPROD" ] || FAILS="$FAILS untouched_violated"
for f in host/surfacepage/app/app.js host/daemon/surface_page.go host/daemon/daemon.go; do
  if git -C "$REF" show "$B:$f" | cmp -s - "$f"; then line "UNCHANGED: $f"
  else FAILS="$FAILS served_bytes_changed"; fi
done
CALLERS=$(rg -n -g '*.go' -g '!**/*_test.go' -g '!host/store/**' '\.(RequestOperatorAction|OperatorAction|ConfirmOperatorAction|ClaimOperatorAction|PruneOperatorActions|RecheckClaimedOperatorAction)\(' cmd host | rg -v 'approval\.ConfirmOperatorAction')
CONTROL=$(rg -c -g '*.go' -g '!**/*_test.go' 'authority\.RecheckVerb\(' cmd host | wc -l | tr -d ' ')
line "NEW_PRODUCTION_CALLERS: ${CALLERS:-none} (authority.RecheckVerb control: $CONTROL)"
[ -z "$CALLERS" ] && [ "$CONTROL" -gt 0 ] || FAILS="$FAILS caller_fence"
if [ "${GATES_SKIP_VERIFY_GO:-0}" = 1 ]; then line "LEG verify_go SKIPPED (GATES_SKIP_VERIFY_GO=1)"
else leg verify_go 'VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh'; fi
line "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}"
tail -n 1 "$L/legs.txt" > "$L/ALL_DONE.$SHA"
[ -z "$FAILS" ]
````


## 7. Executor directive notes

- One run M1 → M2 → M3. `/usr/bin/patch -p1 -E` is the only worktree writer. No Git writes, hand edits, gofmt writes, Chrome, golden regeneration, `nohup`, background test legs, or real store operations.
- With `W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter428-row243-v4s3` and `T=$HOME/.ailang/state/world-iter428/exec-patches` and `PLAN=$W/design_docs/planned/w-world-address-bar-verbs-v4s3-sprint-plan.md`, extract the script with `awk '$0 == "````bash accept.sh" {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/accept.sh"`. Run `bash "$T/accept.sh" "$W" "$M" fast 2>&1 | tee "$HOME/.ailang/state/world-iter428/exec-logs/accept_$M.txt"` under pipefail.
- Commit-time commands for M1, M2, M3 respectively are `bash "$T/accept.sh" "$W" M1 fast`, `bash "$T/accept.sh" "$W" M2 fast`, `bash "$T/accept.sh" "$W" M3 fast`. Each exports the exact pinned AILANG_BIN, isolated AILANG_STATE_DIR and test fixture env. The controller expands the paths before forwarding to `commit_ms.sh <M> <message-file> <test-cmd…>`. Record the two plan files first, because its cleanup refuses untracked design docs outside the snapshot.
- Snapshot only declared changes cumulatively against base, excluding the plans. The optional `SNAP_STATUS_LIST` was used only to rehearse with an archive tree. Extract the snapshot fence with §2's pattern adjusted to the fence label. Stop if a snapshot already exists until the controller checks it.
- Stop on an acceptance failure, skipped expected test, scope/status discrepancy, numstat mismatch or a filtered cumulative product diff over the hard cap. Loopback bind denials are **UNINFORMATIVE UNDER SANDBOX**. Keep their output and hand the exact gate over to the controller. A test failing independently of that denial is a real failure.
- Re-run a documented load-sensitive failure once alone, including the broker exec probes, `TestSeedRestartIgnoresASurvivingOldCopy`, `TestMissionConfigDefaultFlowIntegration`, `TestSeToolsRunBenchmarkTasksEndToEnd`, `TestHostsLaunchersExecHost`, and `TestBackoff`. Leave transport or scheduling failures visible in the gate log.

````bash snapshot.sh
#!/bin/bash
set -euo pipefail
W=$1
M=$2
B=92e1c44ebd8dbbdb2b962ef54aed0dd9a0bf2456
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter428/snap}/$M
mkdir "$D"
cd "$W"
if [ -n "${SNAP_STATUS_LIST:-}" ]; then
  cp "$SNAP_STATUS_LIST" "$D/STATUS" # scratch rehearsal only
else
  {
    git diff --name-status --no-renames "$B" -- cmd host docs app
    git ls-files --others --exclude-standard -- cmd host docs app | sed 's/^/A\t/'
  } | sort -u -t "$(printf '\t')" -k2 > "$D/STATUS"
fi
awk -F '\t' '$1 != "D" {print $2}' "$D/STATUS" | while IFS= read -r p; do
  mkdir -p "$D/$(dirname "$p")"
  cp -p "$p" "$D/$p"
done
awk -F '\t' '$1 == "D" {print $2}' "$D/STATUS" > "$D/DELETED"
(
  cd "$D"
  for d in cmd host docs app; do
    if [ -d "$d" ]; then find "$d" -type f -print0; fi
  done | sort -z | xargs -0 shasum -a 256 > MANIFEST.sha256
  shasum -a 256 -c MANIFEST.sha256
)
wc -l "$D"/{STATUS,DELETED,MANIFEST.sha256}
````


## 8. Risks

- **Attended deployment risk.** With a free credentials writer, a successful candidate start commits v10 before health. The existing busy-writer retry-on-restart path can leave v9 in place until another start. If upgrade then fails, including a failure to save history after health, automatic rollback can install .prev and stop the service because v9 refuses that credentials store. The v10 data remains intact and the candidate binary is retained as .failed when the rollback swap succeeds. A compatibility-safe v10 build can open the store. No current service command automatically restores the v10 candidate or old data. Controller must hold live deployment until an attended compatible recovery and backup procedure is ready, or route a compatibility guard separately. The current fallback next-command message points at the older binary. No service code or real store was touched.
- S4 must hold the kernel runner and episode locks before recovery, check live process birth provenance and own status freshness. Heartbeat and the phase generation alone are insufficient evidence to release a reservation. S3 retains uncertain rows and their identities across prune and restart.
- S10 must pin the V3 association, arbitration and read-back. A V3 row can age out; preserving the logical link permits an explicit missing-card refusal. No V3 supersession/prune behavior changes in v10.
- S3 methods introduce no writer-open site. S4 and S7 must keep held credentials handles short and preserve existing `LOCK_NB` and writer-site inventories. Do not share an unbounded writer across waits.

## 9. Acceptance criteria

1. Only the exact §5 paths change. No daemon route, `.ail`, CSP, app.js, collision path or reachable operator behavior. Caller searches have a legacy positive control.
2. Fresh v10 and upgraded v9 have identical SQLite schema objects, and failed migration rolls back DDL/version. A base reader and writer refuse a v10 fixture; v9 bound writes retain their existing floor.
3. Admission cap, two dedup paths, stable reopen receipt, episode exclusion, clock boundaries and pruning tombstones are behaviorally pinned. No proof or effect exists at requested admission.
4. Real SSH and WebAuthn artifacts round trip with exact `format:name`. Counter/use and state are atomic, replay is idempotent and ambiguous errors never return success.
5. Claim requires one correct owned current runner, a same-digest live confirmed generation and one affected row. Refusal leaves the durable row and steps unchanged. The runner rechecks its frozen intent independently of the receipt.
6. Judge-409 (a) fd `println` control is caught and all newly staged server-facing store methods are silent. Judge-409 (b) exact stored identity rechecks successfully and a changed form refuses.
7. Score the §4 named and whole-package verdicts, including explicitly reported surviving or unmeasured mutants. The controller checks every required inventory and any sandbox-uninformative legs before merge.
8. Each cumulative Appendix P patch round-trips byte-identically to the corresponding scratch tree with the controller fence protocol, and snapshots verify manifests. S3b's separate patch is not in the three-milestone execution loop.

## Appendix C — after this sprint

- **S3b, measured 132 changed lines.** Compile-tested draft step reserve, observe and read, including reopen with retained reservation and exact step version. Apply after M3 only in an independently budgeted follow-up. Add fd coverage for these new entry points, full receipt predicate mutations and current-owner restart/recovery gates before enabling an effect. `TestOperatorStepRestart` separates reservation persistence and observation from child identity/recovery work. Patch label `S3b` and stage `st-S3b`. The compiled draft passes `AILANG_BIN=$HOME/.pinned-ailang/ailang go test -count=1 -v ./host/store -run '^TestOperatorStepRestart$'` and `go vet ./host/store`. S3 already keeps operation and effect identities and uncertain reservations when pruning. No step dispatch can use the dormant table alone.
- **S4** inherits the v10 runner row, owner and generation CAS, and claimed receipt recheck. Owns install/process/status/upgrade seam, heartbeat/catalogue (including per-verb E2 gate), `TestOperatorRunnerUnavailable`, `TestOperatorCatalogueAdmission`, `TestOperatorRecoverUntouched`, `TestOperatorRecoverReserved`, `TestOperatorStatusDisconnected`, `TestOperatorCancelInflight`, kernel and real terminal half of `TestOperatorTerminalExclusion`, daemon-status half of RunnerRechecks, and §5.5 service tests beyond schema. The terminal runner half needs the actual shared lock owner and the start/launch/finish refusal path, including exact action id, ≤5 s non-waiting behavior and zero effect calls. The existing PID file and short store writer cannot establish that property independently.
- **S5–S6** own §5.3 headless start/mint/file/launch and reconciliation rows from StartOnce through StartBrowser, using the durable receipt and common locks first. **S7** owns guarded routes and §5.1 Misdirected, CrossOrigin, FetchSite, MediaType, BodyBound, UnsupportedMethod, RequestLimits and LiteralTask at route, BarStartReview, BarFinishReview and OtherVerbParity. Invoke `JudgeOperatorAction` from trusted read-only store state and use the confirmation CAS afterward. **S8** owns PanelNetwork, PanelLatest, LinkAllowlist, WireTypes and browser fixtures. **S9–S11** own §5.4 finish, linked independent V3 confirmation, exact PR/remote recovery and scoped record; schema-10 has the finish shape. **S12** owns HerdrUnavailable, HerdrHandoff and HerdrWindow under 119 A plus E4.
- **Carried from S2b:** judge-408 note 4 field-drop redundancy through envelope validation; judge-407 note 3 numeric string flags, usage rendering at S7; judge-407 note 2 / S1 F8 null arrays at S8. Original two-store V3 window and judge-402 v9 prune/DDL waits remain separate. Uncertain resolution lacks an attended resolver; an outliving headless child is unowned without proof; new credentials writer sites move together with `LOCK_NB` coverage; herdr plugin `codex=gpt-6.1-sol` hard-code stays with S12.
- **Controller:** leave the `ctl_roundtrip.sh` loop at M1 M2 M3. Keep S3b ahead of runner effect activation. Resolve the attended compatible rollout path in §8 before upgrading Mark's service. The reviewed design and ledger stay closed.

## Appendix P — compiled patches

### Patch M1

````diff M1
--- a/cmd/ailang-worldd/approve_await_refusals_test.go
+++ b/cmd/ailang-worldd/approve_await_refusals_test.go
@@ -145,8 +145,8 @@
 		gone  bool     // delete the request after the first wait
 		want  string
 	}{
-		{name: "v6-store", setup: []string{`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `PRAGMA user_version = 6`}, want: "surface approvals need schema 7: restart `ailang-worldd serve`"},
-		{name: "v5-store", setup: []string{`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`}, want: "passkeys need schema 6"},
+		{name: "v6-store", setup: []string{`DROP TABLE IF EXISTS operator_action_steps; DROP TABLE IF EXISTS operator_actions; DROP TABLE IF EXISTS operator_runner_state; DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `PRAGMA user_version = 6`}, want: "surface approvals need schema 7: restart `ailang-worldd serve`"},
+		{name: "v5-store", setup: []string{`DROP TABLE IF EXISTS operator_action_steps; DROP TABLE IF EXISTS operator_actions; DROP TABLE IF EXISTS operator_runner_state; DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`}, want: "passkeys need schema 6"},
 		{name: "deleted", gone: true, want: "is no longer in the store"},
 	} {
 		t.Run(arm.name, func(t *testing.T) {
--- a/cmd/ailang-worldd/approve_passkey_refusals_test.go
+++ b/cmd/ailang-worldd/approve_passkey_refusals_test.go
@@ -56,7 +56,7 @@
 		{name: "WF2-missing", path: func(t *testing.T) string { return filepath.Join(t.TempDir(), "absent.json") }, want: "--webauthn-assertion: open "},
 		{name: "WF3-not-an-assertion", file: func(*testing.T, approveRig) []byte { return []byte(`{"id":"x"}`) },
 			want: "is not a passkey assertion: " + approval.ErrPasskeyResponse.Error()},
-		{name: "WF4-v5-store", setup: sql(`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`), file: good,
+		{name: "WF4-v5-store", setup: sql(`DROP TABLE IF EXISTS operator_action_steps; DROP TABLE IF EXISTS operator_actions; DROP TABLE IF EXISTS operator_runner_state; DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`), file: good,
 			want: "passkeys need schema 6"},
 		{name: "WF5-unknown-credential", file: func(t *testing.T, a approveRig) []byte {
 			return approvaltest.NewPasskey(approval.AlgEdDSA, 55).Assert(pkOrigin, pkRPID, requestHashOf(t, a), 2)
--- a/cmd/ailang-worldd/approve_passkey_test.go
+++ b/cmd/ailang-worldd/approve_passkey_test.go
@@ -96,7 +96,7 @@
 			t.Fatalf("rc %d, terminal %v, daemon calls %d, stderr %q", code, opened, *calls, stderr.String())
 		}
 		raw := rawCredentials(t, a.credPath)
-		for _, stmt := range []string{`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
+		for _, stmt := range []string{`DROP TABLE IF EXISTS operator_action_steps; DROP TABLE IF EXISTS operator_actions; DROP TABLE IF EXISTS operator_runner_state; DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
 			if _, err := raw.Exec(stmt); err != nil {
 				t.Fatal(err)
 			}
--- a/cmd/ailang-worldd/approve_sig_test.go
+++ b/cmd/ailang-worldd/approve_sig_test.go
@@ -216,7 +216,7 @@
 		{name: "v4-store", setup: func(t *testing.T, a approveRig) {
 			raw := rawCredentials(t, a.credPath)
 			defer func() { _ = raw.Close() }()
-			for _, stmt := range []string{`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`} {
+			for _, stmt := range []string{`DROP TABLE IF EXISTS operator_action_steps; DROP TABLE IF EXISTS operator_actions; DROP TABLE IF EXISTS operator_runner_state; DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`} {
 				if _, err := raw.Exec(stmt); err != nil {
 					t.Fatal(err)
 				}
--- a/cmd/ailang-worldd/approver_label_test.go
+++ b/cmd/ailang-worldd/approver_label_test.go
@@ -82,7 +82,7 @@
 	db5, creds5 := approverStore(t)
 	seedApprover(t, creds5, "s", fixtureKeyOf(t, "ed25519.pub"), approveT0)
 	raw = rawCredentials(t, creds5)
-	for _, stmt := range []string{`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
+	for _, stmt := range []string{`DROP TABLE IF EXISTS operator_action_steps; DROP TABLE IF EXISTS operator_actions; DROP TABLE IF EXISTS operator_runner_state; DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
 		if _, err := raw.Exec(stmt); err != nil {
 			t.Fatal(err)
 		}
--- a/cmd/ailang-worldd/approver_passkey_test.go
+++ b/cmd/ailang-worldd/approver_passkey_test.go
@@ -154,7 +154,7 @@
 		{name: "N6-v5", setup: func(t *testing.T, r enrolRig) {
 			raw := rawCredentials(t, r.creds)
 			defer func() { _ = raw.Close() }()
-			for _, stmt := range []string{`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
+			for _, stmt := range []string{`DROP TABLE IF EXISTS operator_action_steps; DROP TABLE IF EXISTS operator_actions; DROP TABLE IF EXISTS operator_runner_state; DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
 				if _, err := raw.Exec(stmt); err != nil {
 					t.Fatal(err)
 				}
--- a/cmd/ailang-worldd/approver_rig_test.go
+++ b/cmd/ailang-worldd/approver_rig_test.go
@@ -64,7 +64,7 @@
 	db, creds = approverStore(t)
 	raw := rawCredentials(t, creds)
 	defer func() { _ = raw.Close() }()
-	for _, stmt := range []string{`DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`} {
+	for _, stmt := range []string{`DROP TABLE IF EXISTS operator_action_steps; DROP TABLE IF EXISTS operator_actions; DROP TABLE IF EXISTS operator_runner_state; DROP TABLE verb_requests`, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`} {
 		if _, err := raw.Exec(stmt); err != nil {
 			t.Fatal(err)
 		}
--- a/cmd/ailang-worldd/episode_finish_surface_test.go
+++ b/cmd/ailang-worldd/episode_finish_surface_test.go
@@ -250,7 +250,7 @@
 	if got := putFinishVerbRequest(creds, row, sfT0); got != "cannot store request vr-00000000000000f1: "+store.ErrTooManyPendingVerbRequests.Error() {
 		t.Fatalf("cap: %q", got)
 	}
-	credSQL(t, creds, `DROP TABLE verb_requests`, `PRAGMA user_version = 8`)
+	credSQL(t, creds, `DROP TABLE IF EXISTS operator_action_steps; DROP TABLE IF EXISTS operator_actions; DROP TABLE IF EXISTS operator_runner_state; DROP TABLE verb_requests`, `PRAGMA user_version = 8`)
 	if got := putFinishVerbRequest(creds, row, sfT0); !strings.Contains(got, "verb requests need schema 9") || !strings.Contains(got, "restart `ailang-worldd serve`") {
 		t.Fatalf("v8: %q", got)
 	}
--- a/cmd/ailang-worldd/episode_preset.go
+++ b/cmd/ailang-worldd/episode_preset.go
@@ -238,7 +238,7 @@
 		}
 		version := ro.Version()
 		_ = ro.Close()
-		if version < store.CredentialsSchemaVersion {
+		if version < store.CredentialsBindingWriteVersion {
 			return fail(&store.UnboundCredentialsStoreError{Path: req.credPath})
 		}
 	}
--- a/cmd/ailang-worldd/session_new.go
+++ b/cmd/ailang-worldd/session_new.go
@@ -378,7 +378,7 @@
 		return nil // a store this can't read yet (pre-split): the mint moves it first
 	}
 	defer func() { _ = ro.Close() }()
-	if ro.Version() < store.CredentialsSchemaVersion {
+	if ro.Version() < store.CredentialsBindingWriteVersion {
 		return &store.UnboundCredentialsStoreError{Path: credPath}
 	}
 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
--- a/host/authority/judge_test.go
+++ b/host/authority/judge_test.go
@@ -57,6 +57,9 @@
 	}
 	defer func() { _ = db.Close() }()
 	for _, s := range stmts {
+		if strings.HasPrefix(s, "DROP TABLE verb_requests") {
+			s = `DROP TABLE IF EXISTS operator_action_steps; DROP TABLE IF EXISTS operator_actions; DROP TABLE IF EXISTS operator_runner_state; ` + s
+		}
 		if _, err := db.Exec(s); err != nil {
 			t.Fatal(err)
 		}
--- a/host/daemon/surface_approve_rig_test.go
+++ b/host/daemon/surface_approve_rig_test.go
@@ -171,6 +171,9 @@
 	}
 	defer func() { _ = db.Close() }()
 	for _, s := range statements {
+		if strings.HasPrefix(s, "DROP TABLE verb_requests") {
+			s = `DROP TABLE IF EXISTS operator_action_steps; DROP TABLE IF EXISTS operator_actions; DROP TABLE IF EXISTS operator_runner_state; ` + s
+		}
 		if _, err := db.Exec(s); err != nil {
 			t.Fatal(err)
 		}
--- a/host/store/approvers_test.go
+++ b/host/store/approvers_test.go
@@ -131,7 +131,7 @@
 			}
 			db := rawDB(t, path)
 			ddl := tableDDL(t, db)
-			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "key_decisions", "mint_requests", "session_credentials", "surface_approvals", "verb_requests", "webauthn_credentials"})
+			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "key_decisions", "mint_requests", "operator_action_steps", "operator_actions", "operator_runner_state", "session_credentials", "surface_approvals", "verb_requests", "webauthn_credentials"})
 			if got, want := normalizeDDL(ddl["approvers"]), normalizeDDL(canonicalApproversDDL); got != want {
 				t.Fatalf("upgraded approvers DDL:\n got: %s\nwant: %s", got, want)
 			}
--- a/host/store/credentials.go
+++ b/host/store/credentials.go
@@ -30,10 +30,13 @@
 	"slices"
 )
 
-// credentialsSchemaSQL is the credentials store's whole schema.
+// credentialsSchemaSQL is the credentials schema through version 9.
 //
 //go:embed credentials_schema.sql
 var credentialsSchemaSQL string
+
+//go:embed operator_schema.sql
+var operatorSchemaSQL string
 
 // credentialsSchemaVersion is the credentials store's own user_version,
 // independent of the world store's currentSchemaVersion. Version 2 added
@@ -41,8 +44,12 @@
 // version 4 adds mint requests (row 201); version 5 adds approvers (row 201 A1b); version 6 adds
 // passkeys (row 201 A2); version 7 adds the surface's recorded approvals (row 201 A3); version 8
 // adds the key decisions of Human.Approve requests (row 201 A4); version 9 adds the verb requests a
-// key confirms on the surface (row 243 V3).
-const credentialsSchemaVersion = 9
+// key confirms on the surface (row 243 V3); version 10 adds operator actions.
+const credentialsSchemaVersion = 10
+
+// CredentialsBindingWriteVersion preserves the existing floor for bound
+// session, start and resume writes on the previous schema.
+const CredentialsBindingWriteVersion = 9
 
 // CredentialsSchemaVersion exports credentialsSchemaVersion: serve compares a
 // store's version with it to decide whether to upgrade.
@@ -73,7 +80,9 @@
 // credentialsTables is the exact, sorted application table set of a credentials store at version.
 func credentialsTables(version int) []string {
 	switch {
-	case version >= 9:
+	case version == 10:
+		return []string{approversTable, "episode_starts", keyDecisionsTable, mintRequestsTable, "operator_action_steps", "operator_actions", "operator_runner_state", credentialsTable, surfaceApprovalsTable, verbRequestsTable, passkeysTable}
+	case version == 9:
 		return []string{approversTable, "episode_starts", keyDecisionsTable, mintRequestsTable, credentialsTable, surfaceApprovalsTable, verbRequestsTable, passkeysTable}
 	case version == 8:
 		return []string{approversTable, "episode_starts", keyDecisionsTable, mintRequestsTable, credentialsTable, surfaceApprovalsTable, passkeysTable}
@@ -297,11 +306,22 @@
 // version 1 is upgraded only with credOpenUpgrade, and anything else is
 // refused unmodified.
 func openCredentialsSQLite(display, dsn string, mode credOpenMode) (*Credentials, error) {
+	// Apply foreign-key enforcement to every credentials connection, including
+	// connections replaced by database/sql.
+	if !isInMemoryDSN(dsn) {
+		dsn += "&_pragma=foreign_keys(1)"
+	}
 	db, err := sql.Open("sqlite", dsn)
 	if err != nil {
 		return nil, fmt.Errorf("store: open credentials %q: %w", display, err)
 	}
 	db.SetMaxOpenConns(1)
+	if isInMemoryDSN(dsn) {
+		if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
+			_ = db.Close()
+			return nil, err
+		}
+	}
 	c := &Credentials{db: db, display: display}
 	if err := c.enforceSchema(mode); err != nil {
 		_ = db.Close()
@@ -381,6 +401,9 @@
 	if version < 9 {
 		statements = append(statements, verbRequestsSQL)
 	}
+	if version < 10 {
+		statements = append(statements, operatorSchemaSQL)
+	}
 	for _, stmt := range statements {
 		if _, err := tx.Exec(stmt); err != nil {
 			return fmt.Errorf("store: upgrade credentials %q: %w", display, err)
@@ -402,7 +425,7 @@
 	if err != nil {
 		return fmt.Errorf("store: initialize credentials %q: begin: %w", display, err)
 	}
-	if _, err := tx.Exec(credentialsSchemaSQL); err != nil {
+	if _, err := tx.Exec(credentialsSchemaSQL + operatorSchemaSQL); err != nil {
 		_ = tx.Rollback()
 		return fmt.Errorf("store: initialize credentials %q: %w", display, err)
 	}
@@ -457,7 +480,7 @@
 	if (row.Project == "") != (row.ProfileDigest == "") {
 		return fmt.Errorf("store: mint session: a profile binding needs both the project and its digest")
 	}
-	if (row.Bound() || start != nil || resume != nil) && c.version < credentialsSchemaVersion {
+	if (row.Bound() || start != nil || resume != nil) && c.version < CredentialsBindingWriteVersion {
 		return &UnboundCredentialsStoreError{Path: c.display}
 	}
 	tx, err := c.db.BeginTx(ctx, nil)
--- a/host/store/credentials_test.go
+++ b/host/store/credentials_test.go
@@ -76,7 +76,7 @@
 func TestCredentialsSchemaMatchesCanonicalDDL(t *testing.T) {
 	c := openTestCredentials(t)
 	ddl := tableDDL(t, c.db)
-	requireExactTableNames(t, "fresh credentials store", ddl, []string{"approvers", "episode_starts", "key_decisions", "mint_requests", "session_credentials", "surface_approvals", "verb_requests", "webauthn_credentials"})
+	requireExactTableNames(t, "fresh credentials store", ddl, []string{"approvers", "episode_starts", "key_decisions", "mint_requests", "operator_action_steps", "operator_actions", "operator_runner_state", "session_credentials", "surface_approvals", "verb_requests", "webauthn_credentials"})
 	for name, got := range map[string]string{"schema file": ddl["webauthn_credentials"], "passkeysSQL": passkeysSQL} {
 		if normalizeDDL(got) != normalizeDDL(canonicalPasskeysDDL) {
 			t.Fatalf("webauthn_credentials DDL (%s):\n got: %s\nwant: %s", name, got, canonicalPasskeysDDL)
@@ -112,8 +112,8 @@
 		t.Fatalf("v4 world table differs from the version-1 credentials table: %s", got)
 	}
 	var version int
-	if err := c.db.QueryRow("PRAGMA user_version").Scan(&version); err != nil || version != 9 {
-		t.Fatalf("credentials user_version = %d %v, want 9 (row 243 V3)", version, err)
+	if err := c.db.QueryRow("PRAGMA user_version").Scan(&version); err != nil || version != 10 {
+		t.Fatalf("credentials user_version = %d %v, want 10 (row 243 V4)", version, err)
 	}
 }
 
--- a/host/store/mint_requests_test.go
+++ b/host/store/mint_requests_test.go
@@ -129,7 +129,7 @@
 			}
 			db := rawDB(t, path)
 			ddl := tableDDL(t, db)
-			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "key_decisions", "mint_requests", "session_credentials", "surface_approvals", "verb_requests", "webauthn_credentials"})
+			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "key_decisions", "mint_requests", "operator_action_steps", "operator_actions", "operator_runner_state", "session_credentials", "surface_approvals", "verb_requests", "webauthn_credentials"})
 			if got, want := normalizeDDL(ddl["mint_requests"]), normalizeDDL(canonicalMintRequestsDDL); got != want {
 				t.Fatalf("upgraded mint_requests DDL:\n got: %s\nwant: %s", got, want)
 			}
--- /dev/null
+++ b/host/store/operator_schema.sql
@@ -0,0 +1,62 @@
+-- Credentials schema 10. Used for both a fresh store and the v9 upgrade.
+CREATE TABLE operator_actions (
+    id          TEXT NOT NULL PRIMARY KEY,
+    episode_id  TEXT NOT NULL,
+    kind        TEXT NOT NULL CHECK (kind IN ('start', 'finish')),
+    review      BLOB NOT NULL CHECK (length(review) = 32),
+    digest      BLOB NOT NULL CHECK (length(digest) = 32),
+    action      BLOB,
+    challenge   BLOB CHECK (length(challenge) = 32),
+    created_at  INTEGER NOT NULL,
+    expires_at  INTEGER NOT NULL CHECK (expires_at > created_at),
+    state       TEXT NOT NULL CHECK (state IN ('requested', 'confirmed', 'claimed', 'running',
+                 'awaiting-pr-confirmation', 'succeeded', 'failed', 'uncertain', 'expired', 'cancelled')),
+    generation  INTEGER NOT NULL,
+    owner       TEXT,
+    approved_by TEXT,
+    artifact    BLOB,
+    terminal_at INTEGER,
+    v3_id       TEXT UNIQUE,
+    CHECK (terminal_at IS NOT NULL OR (action IS NOT NULL AND challenge IS NOT NULL)),
+    CHECK ((approved_by IS NULL) = (artifact IS NULL)),
+    CHECK (state <> 'requested' OR (approved_by IS NULL AND owner IS NULL)),
+    CHECK (state NOT IN ('confirmed', 'claimed', 'running', 'awaiting-pr-confirmation', 'uncertain')
+           OR approved_by IS NOT NULL),
+    CHECK ((owner IS NOT NULL) = (state IN ('claimed', 'running', 'awaiting-pr-confirmation', 'uncertain'))),
+    CHECK ((terminal_at IS NOT NULL) = (state IN ('succeeded', 'failed', 'expired', 'cancelled'))),
+    CHECK (v3_id IS NULL OR kind = 'finish'),
+    CHECK (state <> 'awaiting-pr-confirmation' OR v3_id IS NOT NULL)
+);
+CREATE UNIQUE INDEX operator_active_episode ON operator_actions(episode_id)
+WHERE terminal_at IS NULL;
+CREATE INDEX operator_review ON operator_actions(episode_id, kind, review);
+
+CREATE TABLE operator_runner_state (
+    singleton  INTEGER PRIMARY KEY CHECK (singleton = 1),
+    generation INTEGER NOT NULL,
+    owner      TEXT NOT NULL,
+    protocol   INTEGER NOT NULL,
+    build      TEXT NOT NULL,
+    seen_at    INTEGER NOT NULL,
+    readiness  TEXT NOT NULL,
+    catalogue_generation INTEGER NOT NULL,
+    catalogue BLOB
+);
+
+-- Operations use the same identity ledger as effect reservations. An operation
+-- alias carries its caller's operation id even if another tap created the action.
+CREATE TABLE operator_action_steps (
+    kind       TEXT NOT NULL CHECK (kind IN ('operation', 'worktree', 'mint', 'files', 'launch',
+                 'herdr-handoff', 'checks', 'commit', 'v3-card', 'push', 'label', 'pull-request', 'record', 'cleanup')),
+    identity   TEXT NOT NULL,
+    action_id  TEXT NOT NULL REFERENCES operator_actions(id),
+    phase      TEXT NOT NULL CHECK (phase IN ('reserved', 'observed')),
+    version    INTEGER NOT NULL,
+    owner      TEXT,
+    created_at INTEGER NOT NULL,
+    receipt    BLOB,
+    PRIMARY KEY (kind, identity),
+    CHECK (kind <> 'operation' OR (phase = 'observed' AND owner IS NULL AND receipt IS NULL)),
+    CHECK (kind = 'operation' OR owner IS NOT NULL),
+    CHECK (phase <> 'reserved' OR receipt IS NULL)
+);
--- /dev/null
+++ b/host/store/operator_schema_test.go
@@ -0,0 +1,82 @@
+package store
+
+import (
+	"database/sql"
+	"path/filepath"
+	"reflect"
+	"testing"
+)
+
+func operatorSchema(t *testing.T, db *sql.DB) []string {
+	t.Helper()
+	rows, err := db.Query(`SELECT type || ':' || name || ':' || coalesce(sql, '') FROM sqlite_schema ORDER BY type, name`)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer rows.Close()
+	var schema []string
+	for rows.Next() {
+		var entry string
+		if err := rows.Scan(&entry); err != nil {
+			t.Fatal(err)
+		}
+		schema = append(schema, normalizeDDL(entry))
+	}
+	if err := rows.Err(); err != nil {
+		t.Fatal(err)
+	}
+	return schema
+}
+
+func TestOperatorSchema10(t *testing.T) {
+	path := filepath.Join(t.TempDir(), CredentialsFileName)
+	raw := rawDB(t, path)
+	if _, err := raw.Exec(credentialsSchemaSQL); err != nil {
+		t.Fatal(err)
+	}
+	setPositiveFixtureVersion(t, raw, 9)
+	old, err := OpenCredentials(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	row := oneBoundRow()
+	if old.Version() != 9 || old.MintSession(testCtx(t), row) != nil {
+		t.Fatal("v9 terminal open must retain its version and bound writes")
+	}
+	if err := old.Close(); err != nil {
+		t.Fatal(err)
+	}
+	// Fail after CREATE TABLE succeeds. Partial DDL and version must roll back.
+	if _, err := raw.Exec(`CREATE INDEX operator_active_episode ON session_credentials(episode_id)`); err != nil {
+		t.Fatal(err)
+	}
+	before := operatorSchema(t, raw)
+	if ok, err := UpgradeCredentials(path); err == nil || ok {
+		t.Fatalf("colliding migration: %v %v", ok, err)
+	}
+	if userVersion(t, path) != 9 || !reflect.DeepEqual(before, operatorSchema(t, raw)) {
+		t.Fatal("failed migration changed schema")
+	}
+	if _, err := raw.Exec(`DROP INDEX operator_active_episode`); err != nil {
+		t.Fatal(err)
+	}
+	if ok, err := UpgradeCredentials(path); err != nil || !ok || userVersion(t, path) != 10 {
+		t.Fatalf("v9 upgrade: %v %v", ok, err)
+	}
+	if ok, err := UpgradeCredentials(path); err != nil || ok {
+		t.Fatalf("repeat upgrade: %v %v", ok, err)
+	}
+	fresh := openTestCredentials(t)
+	if fresh.Version() != 10 || !reflect.DeepEqual(operatorSchema(t, raw), operatorSchema(t, fresh.db)) {
+		t.Fatal("fresh and upgraded SQLite schema differ")
+	}
+	ro, err := OpenCredentialsReadOnly(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer ro.Close()
+	got, ok, err := ro.ResolveSession(testCtx(t), row.CredentialID)
+	if err != nil || !ok || got != row {
+		t.Fatalf("lost v9 credential: %+v %v %v", got, ok, err)
+	}
+}
--- a/host/store/verb_requests_test.go
+++ b/host/store/verb_requests_test.go
@@ -98,16 +98,16 @@
 	return out
 }
 
-// CS1, CS2, CS7(b): a fresh store is schema 9; a real schema-8 store (makeV6Credentials, then the
+// CS1, CS2, CS7(b): a real schema-8 store (makeV6Credentials, then the
 // version 7 and 8 tables) opens at 8 and refuses every verb request method naming the serve
-// restart; UpgradeCredentials takes it to 9, and then every table's columns, and verb_requests'
+// restart; UpgradeCredentials takes it through 9 to the current version, and every table's columns and verb_requests'
 // whole DDL, equal the fresh store's.
 func TestCredentialsSchema9VerbRequests(t *testing.T) {
 	fresh := openTestCredentials(t)
 	if got, want := credentialsTables(9), []string{"approvers", "episode_starts", "key_decisions", "mint_requests", "session_credentials", "surface_approvals", "verb_requests", "webauthn_credentials"}; !reflect.DeepEqual(got, want) {
 		t.Fatalf("credentialsTables(9) = %v", got)
 	}
-	if fresh.Version() != 9 || CredentialsSchemaVersion != 9 || MaxVerbRequestBytes != approval.MaxVerbRequestBytes {
+	if fresh.Version() != CredentialsSchemaVersion || MaxVerbRequestBytes != approval.MaxVerbRequestBytes {
 		t.Fatalf("fresh version %d, constant %d", fresh.Version(), CredentialsSchemaVersion)
 	}
 	path := filepath.Join(t.TempDir(), CredentialsFileName)
@@ -134,14 +134,14 @@
 		}
 	}
 	_ = v8.Close()
-	if upgraded, err := UpgradeCredentials(path); err != nil || !upgraded || userVersion(t, path) != 9 {
+	if upgraded, err := UpgradeCredentials(path); err != nil || !upgraded || userVersion(t, path) != CredentialsSchemaVersion {
 		t.Fatalf("UpgradeCredentials: %v %v, version %d", upgraded, err, userVersion(t, path))
 	}
 	if again, err := UpgradeCredentials(path); err != nil || again {
 		t.Fatalf("second UpgradeCredentials: %v %v", again, err)
 	}
 	up := rawDB(t, path)
-	requireExactTableNames(t, "upgraded v8 store", tableDDL(t, up), credentialsTables(9))
+	requireExactTableNames(t, "upgraded v8 store", tableDDL(t, up), credentialsTables(CredentialsSchemaVersion))
 	if got, want := tableInfo(t, up), tableInfo(t, fresh.db); !reflect.DeepEqual(got, want) || strings.Count(strings.Join(want, "\n"), "verb_requests.") != 12 {
 		t.Fatalf("upgraded columns differ from a fresh store's:\n got: %v\nwant: %v", got, want)
 	}
@@ -150,9 +150,9 @@
 	}
 }
 
-// CS7(a): a store newer than this binary (user_version 10) is refused by both opens with the
-// existing "newer than this binary's" text, now at 9, and the refusal changes nothing.
-func TestCredentialsNewerThanNineIsRefused(t *testing.T) {
+// CS7(a): a store newer than this binary (user_version 11) is refused by both opens with the
+// existing "newer than this binary's" text, now at 10, and the refusal changes nothing.
+func TestCredentialsNewerThanTenIsRefused(t *testing.T) {
 	path := filepath.Join(t.TempDir(), CredentialsFileName)
 	c, err := OpenCredentials(path)
 	if err != nil {
@@ -160,7 +160,7 @@
 	}
 	_ = c.Close()
 	db := rawDB(t, path)
-	setPositiveFixtureVersion(t, db, 10)
+	setPositiveFixtureVersion(t, db, 11)
 	before, _ := schemaState(t, db)
 	for name, open := range map[string]func(string) (*Credentials, error){"writable": OpenCredentials, "read-only": OpenCredentialsReadOnly, "upgrade": func(p string) (*Credentials, error) { _, err := UpgradeCredentials(p); return nil, err }} {
 		c, err := open(path)
@@ -168,11 +168,11 @@
 			_ = c.Close()
 		}
 		var se *CredentialsSchemaError
-		if !errors.As(err, &se) || !strings.Contains(err.Error(), "user_version 10 is newer than this binary's 9; use a compatible binary") {
-			t.Errorf("%s open of a v10 store: %v", name, err)
-		}
-	}
-	if after, v := schemaState(t, db); !reflect.DeepEqual(after, before) || v != 10 {
+		if !errors.As(err, &se) || !strings.Contains(err.Error(), "user_version 11 is newer than this binary's 10; use a compatible binary") {
+			t.Errorf("%s open of a v11 store: %v", name, err)
+		}
+	}
+	if after, v := schemaState(t, db); !reflect.DeepEqual(after, before) || v != 11 {
 		t.Fatalf("refused opens changed the store: %v -> %v (version %d)", before, after, v)
 	}
 }
````

### Patch M2

````diff M2
--- a/host/store/credentials_test.go
+++ b/host/store/credentials_test.go
@@ -307,6 +307,7 @@
 	_, _, errs["VerbRequest"] = c.VerbRequest(ctx, "x")
 	_, errs["PendingVerbRequests"] = c.PendingVerbRequests(ctx, 1000, 0)
 	errs["ConsumeVerbRequest"] = c.ConsumeVerbRequest(ctx, "x", "y", ApproverUse{}, "", nil, "", 1000)
+	_, errs["OperatorAction"] = c.OperatorAction(ctx, "x", 1000)
 	for name, err := range errs {
 		if !errors.Is(err, ErrNoDeadline) {
 			t.Errorf("%s(Background) = %v, want ErrNoDeadline", name, err)
--- /dev/null
+++ b/host/store/operator_actions.go
@@ -0,0 +1,181 @@
+package store
+
+import (
+	"bytes"
+	"context"
+	"database/sql"
+	"errors"
+	"regexp"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+var (
+	ErrOperatorInput    = errors.New("store: invalid operator request")
+	ErrOperatorConflict = errors.New("store: operator action conflicts with existing work")
+	ErrOperatorCapacity = errors.New("store: sixteen operator actions already pending")
+	ErrOperatorMissing  = errors.New("store: operator action missing")
+	ErrOperatorExpired  = errors.New("store: operator action expired")
+)
+
+const operatorColumns = `id, episode_id, kind, review, digest, action, challenge, created_at, expires_at, state,
+	generation, coalesce(owner, ''), coalesce(approved_by, ''), artifact, coalesce(terminal_at, 0), coalesce(v3_id, '')`
+
+type OperatorActionRow struct {
+	ID, Episode, Kind, State, Owner, ApprovedBy, V3ID string
+	Review, Digest, Action, Challenge, Artifact       []byte
+	CreatedAt, ExpiresAt, Generation, TerminalAt      int64
+}
+
+type OperatorRequest struct {
+	Action               approval.OperatorAction
+	OperationID          string
+	Challenge            [32]byte
+	CreatedAt, ExpiresAt int64
+}
+
+var operatorOperationID = regexp.MustCompile(`^[0-9a-f]{32}$`)
+
+func (c *Credentials) operatorAccess(ctx context.Context) error {
+	if err := requireDeadline(ctx); err != nil {
+		return err
+	}
+	if c.version < 10 {
+		return &OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: 10, Feature: "operator actions"}
+	}
+	return nil
+}
+
+func scanOperatorAction(row *sql.Row) (a OperatorActionRow, err error) {
+	err = row.Scan(&a.ID, &a.Episode, &a.Kind, &a.Review, &a.Digest, &a.Action, &a.Challenge,
+		&a.CreatedAt, &a.ExpiresAt, &a.State, &a.Generation, &a.Owner, &a.ApprovedBy, &a.Artifact, &a.TerminalAt, &a.V3ID)
+	if errors.Is(err, sql.ErrNoRows) {
+		err = ErrOperatorMissing
+	}
+	return a, err
+}
+
+// OperatorAction reads the durable row and projects expiry for unclaimed work.
+func (c *Credentials) OperatorAction(ctx context.Context, id string, now int64) (OperatorActionRow, error) {
+	if err := c.operatorAccess(ctx); err != nil {
+		return OperatorActionRow{}, err
+	}
+	row, err := scanOperatorAction(c.db.QueryRowContext(ctx, `SELECT `+operatorColumns+` FROM operator_actions WHERE id = ?`, id))
+	if (row.State == "requested" || row.State == "confirmed") && row.ExpiresAt <= now {
+		row.State, row.TerminalAt = "expired", row.ExpiresAt
+	}
+	return row, err
+}
+
+// RequestOperatorAction records inert work or permanently binds another tap's
+// operation ID to the existing action. The active-episode index arbitrates starts
+// and finishes together.
+func (c *Credentials) RequestOperatorAction(ctx context.Context, request OperatorRequest, now int64) (OperatorActionRow, error) {
+	if err := c.operatorAccess(ctx); err != nil {
+		return OperatorActionRow{}, err
+	}
+	a := request.Action
+	if !operatorOperationID.MatchString(request.OperationID) ||
+		request.CreatedAt <= 0 || now < request.CreatedAt ||
+		request.ExpiresAt-request.CreatedAt > int64(approval.OperatorActionLifetime/time.Second) ||
+		c.lock == nil || a.Store != c.lock.dbPath {
+		return OperatorActionRow{}, ErrOperatorInput
+	}
+	if now >= request.ExpiresAt {
+		return OperatorActionRow{}, ErrOperatorExpired
+	}
+	action, err := approval.CanonicalOperatorAction(a)
+	if err != nil {
+		return OperatorActionRow{}, err
+	}
+	review, err := approval.OperatorReviewDigest(a)
+	if err != nil {
+		return OperatorActionRow{}, err
+	}
+	digest, err := approval.OperatorSignatureDigest(a)
+	if err != nil {
+		return OperatorActionRow{}, err
+	}
+	tx, err := c.db.BeginTx(ctx, nil)
+	if err != nil {
+		return OperatorActionRow{}, err
+	}
+	defer tx.Rollback()
+	if err := expireOperatorActions(ctx, tx, now); err != nil {
+		return OperatorActionRow{}, err
+	}
+	row, err := scanOperatorAction(tx.QueryRowContext(ctx, `SELECT `+operatorColumns+` FROM operator_actions WHERE id =
+		(SELECT action_id FROM operator_action_steps WHERE kind = 'operation' AND identity = ?)`, request.OperationID))
+	switch {
+	case err == nil:
+		if !bytes.Equal(row.Review, review[:]) {
+			return OperatorActionRow{}, ErrOperatorConflict
+		}
+		return row, tx.Commit()
+	case !errors.Is(err, ErrOperatorMissing):
+		return OperatorActionRow{}, err
+	}
+	row, err = scanOperatorAction(tx.QueryRowContext(ctx, `SELECT `+operatorColumns+` FROM operator_actions
+		WHERE episode_id = ? AND kind = ? AND review = ? AND state NOT IN ('expired', 'cancelled')
+		ORDER BY terminal_at IS NULL DESC, created_at DESC LIMIT 1`, a.Episode, a.Kind, review[:]))
+	if errors.Is(err, ErrOperatorMissing) {
+		row, err = insertOperatorAction(ctx, tx, request, action, digest, review)
+	}
+	if err != nil {
+		return OperatorActionRow{}, err
+	}
+	if _, err := tx.ExecContext(ctx, `INSERT INTO operator_action_steps
+		(kind, identity, action_id, phase, version, created_at) VALUES ('operation', ?, ?, 'observed', 0, ?)`,
+		request.OperationID, row.ID, now); err != nil {
+		return OperatorActionRow{}, err
+	}
+	return row, tx.Commit()
+}
+
+func insertOperatorAction(ctx context.Context, tx *sql.Tx, r OperatorRequest, action []byte, digest, review [32]byte) (OperatorActionRow, error) {
+	var active, pending int
+	a := r.Action
+	if err := tx.QueryRowContext(ctx, `SELECT count(CASE WHEN episode_id = ? AND terminal_at IS NULL OR id = ? THEN 1 END),
+		count(CASE WHEN state IN ('requested', 'confirmed') THEN 1 END) FROM operator_actions`, a.Episode, a.ID).Scan(&active, &pending); err != nil {
+		return OperatorActionRow{}, err
+	}
+	switch {
+	case active != 0:
+		return OperatorActionRow{}, ErrOperatorConflict
+	case pending >= 16:
+		return OperatorActionRow{}, ErrOperatorCapacity
+	}
+	return scanOperatorAction(tx.QueryRowContext(ctx, `INSERT INTO operator_actions
+		(id, episode_id, kind, review, digest, action, challenge, created_at, expires_at, state, generation)
+		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 'requested', 1) RETURNING `+operatorColumns,
+		a.ID, a.Episode, a.Kind, review[:], digest[:], action, r.Challenge[:], r.CreatedAt, r.ExpiresAt))
+}
+
+func expireOperatorActions(ctx context.Context, tx *sql.Tx, now int64) error {
+	_, err := tx.ExecContext(ctx, `UPDATE operator_actions SET state = 'expired', terminal_at = expires_at
+		WHERE state IN ('requested', 'confirmed') AND expires_at <= ?`, now)
+	return err
+}
+
+// PruneOperatorActions drops old presentation and proof bytes while retaining
+// the action link, all operation/effect identities and any uncertain reservation.
+func (c *Credentials) PruneOperatorActions(ctx context.Context, now int64) error {
+	if err := c.operatorAccess(ctx); err != nil {
+		return err
+	}
+	tx, err := c.db.BeginTx(ctx, nil)
+	if err != nil {
+		return err
+	}
+	defer tx.Rollback()
+	if err := expireOperatorActions(ctx, tx, now); err != nil {
+		return err
+	}
+	if _, err := tx.ExecContext(ctx, `UPDATE operator_actions SET action = NULL, challenge = NULL, artifact = NULL, approved_by = NULL
+		WHERE terminal_at < ? AND NOT EXISTS
+			(SELECT 1 FROM operator_action_steps WHERE action_id = operator_actions.id AND phase = 'reserved')`, now-30*86400); err != nil {
+		return err
+	}
+	return tx.Commit()
+}
--- /dev/null
+++ b/host/store/operator_actions_test.go
@@ -0,0 +1,310 @@
+package store
+
+import (
+	"bytes"
+	"context"
+	"errors"
+	"fmt"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+func operatorRequest(t *testing.T, c *Credentials, n int, episode, kind string) OperatorRequest {
+	t.Helper()
+	a := approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: fmt.Sprintf("oa-%032x", n),
+		Kind: kind, Episode: episode, Store: c.lock.dbPath, Review: "review",
+		Options:    approval.OperatorOptions{Preset: "sol"},
+		Resolution: []approval.OperatorResolved{{Name: "agent", Value: "codex"}}}
+	if kind == "start" {
+		a.Options.Headless, a.Task = true, "test task"
+	}
+	a.Argv = approval.OperatorArgv(kind, episode, a.Options)
+	return OperatorRequest{Action: a, OperationID: fmt.Sprintf("%032x", n), CreatedAt: 1000, ExpiresAt: 1600}
+}
+
+func putOperator(t *testing.T, c *Credentials, r OperatorRequest, now int64) OperatorActionRow {
+	t.Helper()
+	row, err := c.RequestOperatorAction(testCtx(t), r, now)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return row
+}
+
+func getOperator(t *testing.T, c *Credentials, id string, now int64) OperatorActionRow {
+	t.Helper()
+	row, err := c.OperatorAction(testCtx(t), id, now)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return row
+}
+
+func operatorExec(t *testing.T, c *Credentials, query string, args ...any) {
+	t.Helper()
+	if _, err := c.db.Exec(query, args...); err != nil {
+		t.Fatal(err)
+	}
+}
+
+func operatorCounts(t *testing.T, c *Credentials, actions, steps int) {
+	t.Helper()
+	var a, s int
+	err := c.db.QueryRow(`SELECT (SELECT count(*) FROM operator_actions), (SELECT count(*) FROM operator_action_steps)`).Scan(&a, &s)
+	if err != nil || a != actions || s != steps {
+		t.Fatalf("actions/steps: %d/%d, want %d/%d (err %v)", a, s, actions, steps, err)
+	}
+}
+
+func TestOperatorAdmissionBoundaries(t *testing.T) {
+	c := openTestCredentials(t)
+	ctx := testCtx(t)
+	a := operatorRequest(t, c, 1, "ep-1", "start")
+	accepted := putOperator(t, c, a, 1000)
+	if accepted.State != "requested" || accepted.ApprovedBy != "" || len(accepted.Artifact) != 0 || accepted.Owner != "" {
+		t.Fatalf("request is not inert: %+v", accepted)
+	}
+
+	cases := []struct {
+		name string
+		edit func(*OperatorRequest)
+		now  int64
+		want error
+	}{
+		{"operation ID", func(r *OperatorRequest) { r.OperationID = "not-128-bit-hex" }, 1000, ErrOperatorInput},
+		{"other store", func(r *OperatorRequest) { r.Action.Store += ".other" }, 1000, ErrOperatorInput},
+		{"creation zero", func(r *OperatorRequest) { r.CreatedAt, r.ExpiresAt = 0, 1 }, 0, ErrOperatorInput},
+		{"future creation", func(r *OperatorRequest) { r.CreatedAt++ }, 1000, ErrOperatorInput},
+		{"lifetime 601", func(r *OperatorRequest) { r.ExpiresAt++ }, 1000, ErrOperatorInput},
+		{"at expiry", func(r *OperatorRequest) {}, 1600, ErrOperatorExpired},
+		{"past expiry", func(r *OperatorRequest) {}, 1601, ErrOperatorExpired},
+		{"action ID", func(r *OperatorRequest) { r.Action.ID = "invalid" }, 1000, approval.ErrOperatorAction},
+		{"operation changed", func(r *OperatorRequest) { r.OperationID = a.OperationID }, 1000, ErrOperatorConflict},
+	}
+	for _, tc := range cases {
+		t.Run(tc.name, func(t *testing.T) {
+			r := operatorRequest(t, c, 2, "ep-2", "start")
+			tc.edit(&r)
+			if _, err := c.RequestOperatorAction(ctx, r, tc.now); !errors.Is(err, tc.want) {
+				t.Fatalf("refusal %v, want %v", err, tc.want)
+			}
+			operatorCounts(t, c, 1, 1)
+			if got := getOperator(t, c, a.Action.ID, 1000); got.State != "requested" || !bytes.Equal(got.Action, accepted.Action) {
+				t.Fatal("refusal changed the existing request")
+			}
+		})
+	}
+	finish := operatorRequest(t, c, 2, "ep-1", "finish")
+	if _, err := c.RequestOperatorAction(ctx, finish, 1000); !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("cross-kind collision: %v", err)
+	}
+	// The DB must also exclude competing writes which reach it after preflight.
+	if _, err := c.db.Exec(`INSERT INTO operator_actions
+		(id, episode_id, kind, review, digest, action, challenge, created_at, expires_at, state, generation)
+		SELECT ?, episode_id, 'finish', review, digest, action, challenge, created_at, expires_at, state, generation
+		FROM operator_actions WHERE id = ?`, finish.Action.ID, a.Action.ID); err == nil {
+		t.Fatal("database admitted a second active action")
+	}
+	operatorCounts(t, c, 1, 1)
+	for i := 2; i <= 16; i++ {
+		putOperator(t, c, operatorRequest(t, c, i, fmt.Sprintf("ep-%d", i), "start"), 1599)
+	}
+	if _, err := c.RequestOperatorAction(ctx, operatorRequest(t, c, 17, "ep-17", "start"), 1599); !errors.Is(err, ErrOperatorCapacity) {
+		t.Fatalf("17th pending: %v", err)
+	}
+	operatorCounts(t, c, 16, 16)
+	finish.CreatedAt, finish.ExpiresAt = 1600, 2200
+	finish.Action.ID, finish.OperationID = "oa-00000000000000000000000000000011", "00000000000000000000000000000011"
+	if got := putOperator(t, c, finish, 1600); got.Kind != "finish" || got.State != "requested" {
+		t.Fatalf("finish shape or exact-expiry reclamation lost: %+v", got)
+	}
+	operatorCounts(t, c, 17, 17)
+}
+
+func TestOperatorDuplicateTap(t *testing.T) {
+	c := openTestCredentials(t)
+	first := operatorRequest(t, c, 1, "ep", "start")
+	second := operatorRequest(t, c, 2, "ep", "start")
+	second.Challenge[0] = 2
+	row := putOperator(t, c, first, 1000)
+	want, err := approval.OperatorReviewDigest(first.Action)
+	if err != nil || !bytes.Equal(row.Review, want[:]) || bytes.Equal(row.Review, row.Digest) {
+		t.Fatal("store must key on the review digest")
+	}
+	for _, r := range []OperatorRequest{first, second} {
+		if got := putOperator(t, c, r, 1000); got.ID != row.ID || !bytes.Equal(got.Challenge, row.Challenge) {
+			t.Fatal("duplicate did not return the original link and challenge")
+		}
+	}
+	if _, err := c.db.Exec(`INSERT INTO operator_action_steps SELECT * FROM operator_action_steps WHERE identity = ?`, first.OperationID); err == nil {
+		t.Fatal("database admitted a duplicate operation ID")
+	}
+	second.Action.Review += " changed"
+	if _, err := c.RequestOperatorAction(testCtx(t), second, 1000); !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("second tap did not acquire durable id binding: %v", err)
+	}
+	operatorCounts(t, c, 1, 2)
+}
+
+func TestOperatorReconnect(t *testing.T) {
+	first := openTestCredentials(t)
+	path := first.display
+	req := operatorRequest(t, first, 1, "restart", "start")
+	want := putOperator(t, first, req, 1000)
+	if err := first.Close(); err != nil {
+		t.Fatal(err)
+	}
+	second, err := OpenCredentials(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer second.Close()
+	if got := putOperator(t, second, req, 1010); got.ID != want.ID || !bytes.Equal(got.Action, want.Action) || !bytes.Equal(got.Challenge, want.Challenge) {
+		t.Fatal("lost request on reopen")
+	}
+	operatorCounts(t, second, 1, 1)
+}
+
+func TestOperatorPrune(t *testing.T) {
+	c := openTestCredentials(t)
+	const day30 = 30 * 86400
+	r := operatorRequest(t, c, 1, "expired", "start")
+	putOperator(t, c, r, 1000)
+	if getOperator(t, c, r.Action.ID, 1599).State != "requested" ||
+		getOperator(t, c, r.Action.ID, 1600).State != "expired" {
+		t.Fatal("logical expiry boundary")
+	}
+	var stored string
+	if err := c.db.QueryRow(`SELECT state FROM operator_actions WHERE id = ?`, r.Action.ID).Scan(&stored); err != nil || stored != "requested" {
+		t.Fatalf("logical expiry read changed durable state: %q %v", stored, err)
+	}
+	for i, state := range []string{"uncertain", "failed", "succeeded"} {
+		q := operatorRequest(t, c, i+2, state, "start")
+		putOperator(t, c, q, 1000)
+		if state == "uncertain" {
+			operatorExec(t, c, `UPDATE operator_actions SET state = ?, owner = 'runner', approved_by = 'ssh:mark', artifact = x'01' WHERE id = ?`, state, q.Action.ID)
+		} else {
+			operatorExec(t, c, `UPDATE operator_actions SET state = ?, terminal_at = 1600 WHERE id = ?`, state, q.Action.ID)
+		}
+		phase := "reserved"
+		if state == "succeeded" {
+			phase = "observed"
+		}
+		operatorExec(t, c, `INSERT INTO operator_action_steps VALUES ('launch', ?, ?, ?, 1, 'runner', 1000, NULL)`, state, q.Action.ID, phase)
+	}
+	for _, now := range []int64{1600 + day30 - 1, 1600 + day30, 1600 + day30 + 1} {
+		if err := c.PruneOperatorActions(testCtx(t), now); err != nil {
+			t.Fatal(err)
+		}
+		for i, state := range []string{"expired", "uncertain", "failed", "succeeded"} {
+			got := getOperator(t, c, operatorRequest(t, c, i+1, state, "start").Action.ID, now)
+			pruned := now > 1600+day30 && (state == "expired" || state == "succeeded")
+			if got.State != state || (len(got.Action) == 0) != pruned {
+				t.Fatalf("prune %s at %d: %+v", state, now, got)
+			}
+		}
+		operatorCounts(t, c, 4, 7)
+	}
+	retry := operatorRequest(t, c, 10, "expired", "start")
+	retry.OperationID, retry.CreatedAt, retry.ExpiresAt = r.OperationID, 1600+day30, 2200+day30
+	if got := putOperator(t, c, retry, retry.CreatedAt); got.ID != r.Action.ID || got.State != "expired" {
+		t.Fatal("prune lost the retired operation identity")
+	}
+	retry.OperationID = "00000000000000000000000000000010"
+	if got := putOperator(t, c, retry, retry.CreatedAt); got.ID != retry.Action.ID || got.State != "requested" {
+		t.Fatal("fresh review could not supersede expired intent")
+	}
+}
+
+func TestOperatorTerminalExclusion(t *testing.T) {
+	c := openTestCredentials(t)
+	path := c.display
+	r := operatorRequest(t, c, 1, "lock", "start")
+	putOperator(t, c, r, 1000)
+	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
+	defer cancel()
+	contender := make(chan error, 1)
+	go func() {
+		other, err := OpenCredentials(path)
+		if other != nil {
+			other.Close()
+		}
+		contender <- err
+	}()
+	select {
+	case err := <-contender:
+		if !IsWriterAlreadyActive(err) {
+			t.Fatalf("second writer: %v", err)
+		}
+	case <-ctx.Done():
+		_ = c.Close()
+		t.Fatal("credentials writer open waited on the incumbent")
+	}
+	operatorCounts(t, c, 1, 1)
+	if err := c.Close(); err != nil {
+		t.Fatal(err)
+	}
+	second, err := OpenCredentials(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer second.Close()
+	if _, err := second.RequestOperatorAction(testCtx(t), operatorRequest(t, second, 2, "lock", "finish"), 1001); !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("exclusive episode after the writer reopened: %v", err)
+	}
+	operatorCounts(t, second, 1, 1)
+}
+
+func TestOperatorSchema10Constraints(t *testing.T) {
+	c := openTestCredentials(t)
+	r := operatorRequest(t, c, 1, "schema", "start")
+	putOperator(t, c, r, 1000)
+	operatorExec(t, c, `INSERT INTO operator_runner_state VALUES (1, 7, 'runner', 1, 'test', 1000, 'start', 1, NULL)`)
+	cases := []string{
+		`INSERT INTO operator_runner_state SELECT * FROM operator_runner_state`,
+		`INSERT INTO operator_actions SELECT NULL, 'other', kind, review, digest, action, challenge,
+			created_at, expires_at, state, generation, owner, approved_by, artifact, terminal_at, v3_id FROM operator_actions`,
+		`UPDATE operator_actions SET kind = 'other'`,
+		`UPDATE operator_actions SET review = zeroblob(31)`,
+		`UPDATE operator_actions SET digest = zeroblob(31)`,
+		`UPDATE operator_actions SET action = NULL`,
+		`UPDATE operator_actions SET challenge = zeroblob(31)`,
+		`UPDATE operator_actions SET expires_at = created_at`,
+		`UPDATE operator_actions SET state = 'other'`,
+		`UPDATE operator_actions SET artifact = action`,
+		`UPDATE operator_actions SET approved_by = 'ssh:fixture', artifact = action`,
+		`UPDATE operator_actions SET state = 'confirmed'`,
+		`UPDATE operator_actions SET state = 'claimed', approved_by = 'ssh:fixture', artifact = action`,
+		`UPDATE operator_actions SET terminal_at = 1600`,
+		`UPDATE operator_actions SET v3_id = 'vr-0000000000000001'`,
+		`UPDATE operator_actions SET kind = 'finish', state = 'awaiting-pr-confirmation',
+			approved_by = 'ssh:fixture', artifact = action, owner = 'runner'`,
+		`UPDATE operator_action_steps SET action_id = 'missing'`,
+		`INSERT INTO operator_actions SELECT id, 'other', kind, review, digest, action, challenge,
+			created_at, expires_at, state, generation, owner, approved_by, artifact, terminal_at, v3_id FROM operator_actions`,
+		`UPDATE operator_action_steps SET kind = 'unknown', owner = 'runner'`,
+		`UPDATE operator_action_steps SET kind = 'launch', owner = 'runner', phase = 'unknown'`,
+		`UPDATE operator_action_steps SET owner = 'x'`,
+		`UPDATE operator_action_steps SET kind = 'launch'`,
+		`UPDATE operator_action_steps SET kind = 'launch', phase = 'reserved', owner = 'x', receipt = zeroblob(1)`,
+		`INSERT INTO operator_runner_state VALUES (2, 7, 'runner', 1, 'test', 1000, 'ready', 1, NULL)`,
+	}
+	for _, column := range []string{"episode_id", "kind", "review", "digest", "created_at", "expires_at", "state", "generation"} {
+		cases = append(cases, "UPDATE operator_actions SET "+column+" = NULL")
+	}
+	for _, column := range []string{"action_id", "identity", "kind", "phase", "version", "created_at"} {
+		cases = append(cases, "UPDATE operator_action_steps SET "+column+" = NULL")
+	}
+	for _, column := range []string{"generation", "owner", "protocol", "build", "seen_at", "readiness", "catalogue_generation"} {
+		cases = append(cases, "UPDATE operator_runner_state SET "+column+" = NULL")
+	}
+	for _, query := range cases {
+		if _, err := c.db.Exec(query); err == nil || !strings.Contains(err.Error(), "constraint failed") {
+			t.Errorf("%s: %v, want SQLite constraint failure", query, err)
+		}
+	}
+	operatorCounts(t, c, 1, 1)
+}
--- a/host/store/operator_schema_test.go
+++ b/host/store/operator_schema_test.go
@@ -2,6 +2,7 @@
 
 import (
 	"database/sql"
+	"errors"
 	"path/filepath"
 	"reflect"
 	"testing"
@@ -38,6 +39,10 @@
 	old, err := OpenCredentials(path)
 	if err != nil {
 		t.Fatal(err)
+	}
+	_, err = old.OperatorAction(testCtx(t), "missing", 1000)
+	if oldSchema, ok := errors.AsType[*OlderCredentialsStoreError](err); !ok || oldSchema.Version != 9 || oldSchema.Need != 10 {
+		t.Fatalf("schema 9 operator read: %v", err)
 	}
 	row := oneBoundRow()
 	if old.Version() != 9 || old.MintSession(testCtx(t), row) != nil {
````

### Patch M3

````diff M3
--- a/host/approval/operator_signed.go
+++ b/host/approval/operator_signed.go
@@ -7,6 +7,7 @@
 
 import (
 	"bytes"
+	"crypto/sha256"
 	"errors"
 	"fmt"
 	"time"
@@ -21,6 +22,8 @@
 	name, format      string
 	credentialID      []byte
 	signCount         uint32
+	confirmation      *OperatorApprovalBinding
+	artifactHash      [32]byte
 }
 
 // OperatorConfirmInput is what ConfirmOperatorAction judges: the stored action (strictly decoded), the
@@ -96,6 +99,10 @@
 		return nil, err
 	}
 	s := &SignedOperatorAction{ok: true, digest: d, challenge: challenge, expiry: expiry, name: ap.Name, format: format}
+	if confirm {
+		s.confirmation = &OperatorApprovalBinding{Name: ap.Name, Format: format, PublicKey: string(ap.Key.Blob)}
+		s.artifactHash = sha256.Sum256(artifact)
+	}
 	if format == FormatSSH {
 		sig, err := ParseSSHSIG(artifact)
 		if err != nil {
@@ -122,6 +129,10 @@
 	}
 	if !confirm && ad.SignCount > p.SignCount { // R2
 		return nil, fmt.Errorf("%w: enrolled %d, this assertion %d", ErrOperatorCounter, p.SignCount, ad.SignCount)
+	}
+	if confirm {
+		s.confirmation.PublicKey = string(p.Key.Canonical())
+		s.confirmation.RPID, s.confirmation.Origin, s.confirmation.BackupEligible = p.RPID, p.Origin, p.BackupEligible
 	}
 	s.credentialID, s.signCount = bytes.Clone(p.CredentialID), ad.SignCount
 	return s, nil
@@ -162,3 +173,19 @@
 	}
 	return bytes.Clone(s.credentialID), s.signCount, true
 }
+
+// OperatorApprovalBinding describes the verified enrollment snapshot for the
+// store's transactional comparison. It is obtained from a signed proof.
+type OperatorApprovalBinding struct {
+	Name, Format, PublicKey, RPID, Origin string
+	BackupEligible                        bool
+}
+
+// ConfirmationFor binds the proof to the exact artifact to be stored and
+// distinguishes initial confirmation from a runner recheck.
+func (s *SignedOperatorAction) ConfirmationFor(artifact []byte) (OperatorApprovalBinding, error) {
+	if s == nil || !s.ok || s.confirmation == nil || sha256.Sum256(artifact) != s.artifactHash {
+		return OperatorApprovalBinding{}, ErrNotVerified
+	}
+	return *s.confirmation, nil
+}
--- a/host/authority/judge_operator.go
+++ b/host/authority/judge_operator.go
@@ -114,3 +114,19 @@
 }
 
 func (v *VerifiedOperatorAction) Rechecked() bool { return v != nil && v.signed != nil && v.rechecked }
+
+// OperatorConfirmation transports the sealed approval law proof to the store
+// without making store import authority. The store rechecks its enrollment.
+func (v *VerifiedOperatorAction) OperatorConfirmation() *approval.SignedOperatorAction {
+	if v == nil {
+		return nil
+	}
+	return v.signed
+}
+
+// RecheckClaimedOperatorAction uses the runner's frozen inputs. The receipt
+// supplies only the artifact and the identity to compare after verification.
+func RecheckClaimedOperatorAction(ctx context.Context, ro *store.Credentials, frozen approval.OperatorAction, challenge [32]byte,
+	expiry, now int64, row store.OperatorActionRow) (*VerifiedOperatorAction, error) {
+	return RecheckOperatorAction(ctx, ro, frozen, challenge, expiry, now, row.Artifact, row.ApprovedBy)
+}
--- a/host/authority/judge_operator_silent_test.go
+++ b/host/authority/judge_operator_silent_test.go
@@ -1,39 +1,67 @@
+//go:build unix
+
 package authority
 
 // Row 243 V4 sprint 2b (judge-407's retro rule, as sprint 2 M4): the two operator judges S7 and the
 // runner will call write nothing to stdout, stderr or the log, accepting or refusing.
 
 import (
-	"io"
 	"log"
 	"os"
+	"path/filepath"
 	"strings"
 	"testing"
 
 	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/store"
+	"golang.org/x/sys/unix"
 )
 
-// ojSilent runs f with os.Stdout, os.Stderr and the log on one pipe and returns what was written.
+// ojSilent captures writes to descriptors 1 and 2 as well as the log writer.
+// Run only in sequential tests, because the descriptors are process-wide.
 func ojSilent(t *testing.T, f func()) string {
 	t.Helper()
-	r, w, err := os.Pipe()
+	p := filepath.Join(t.TempDir(), "fd-output")
+	capture, err := os.Create(p)
 	if err != nil {
 		t.Fatal(err)
 	}
-	out, errOut := os.Stdout, os.Stderr
-	os.Stdout, os.Stderr = w, w
-	log.SetOutput(w)
-	f()
-	os.Stdout, os.Stderr = out, errOut
-	log.SetOutput(os.Stderr)
-	_ = w.Close()
-	b, _ := io.ReadAll(r)
+	defer capture.Close()
+	func() {
+		for _, fd := range []int{1, 2} {
+			saved, err := unix.Dup(fd)
+			if err != nil {
+				t.Fatal(err)
+			}
+			defer unix.Close(saved)
+			defer func() {
+				if err := unix.Dup2(saved, fd); err != nil {
+					t.Error(err)
+				}
+			}()
+			if err := unix.Dup2(int(capture.Fd()), fd); err != nil {
+				t.Fatal(err)
+			}
+		}
+		logWriter := log.Writer()
+		defer log.SetOutput(logWriter)
+		log.SetOutput(os.Stderr)
+		f()
+	}()
+	b, err := os.ReadFile(p)
+	if err != nil {
+		t.Fatal(err)
+	}
 	return string(b)
 }
 
 func TestOperatorJudgesWriteNothing(t *testing.T) {
-	if got := ojSilent(t, func() { _, _ = os.Stderr.WriteString("x"); _, _ = os.Stdout.WriteString("y"); log.Print("z") }); !strings.HasPrefix(got, "xy") || !strings.HasSuffix(got, " z\n") {
+	if got := ojSilent(t, func() {
+		_, _ = os.Stderr.WriteString("x")
+		_, _ = os.Stdout.WriteString("y")
+		log.Print("z")
+		println("builtin")
+	}); !strings.HasPrefix(got, "xy") || !strings.HasSuffix(got, " z\nbuiltin\n") {
 		t.Fatalf("control: the harness saw %q, want x, y and the log line", got)
 	}
 	in, ssh, pk := ojIn(t)
@@ -55,3 +83,27 @@
 		t.Fatalf("an operator judge wrote %q", got)
 	}
 }
+
+func TestOperatorStoreWritesNothing(t *testing.T) {
+	r := newOperatorStoreRig(t, true)
+	ctx := boundedTestContext(t)
+	const now = 1010
+	got := ojSilent(t, func() {
+		_, _ = r.RequestOperatorAction(ctx, store.OperatorRequest{Action: r.action, OperationID: "1234567890abcdef1234567890abcdef", CreatedAt: 1000, ExpiresAt: 1600}, now)
+		_, _ = r.OperatorAction(ctx, r.cas.ID, now)
+		_, _ = r.OperatorAction(ctx, "missing", now)
+		_, _ = r.ConfirmOperatorAction(ctx, r.cas, r.artifact, nil, now)
+		r.confirm(t)
+		_ = r.ClaimOperatorAction(ctx, r.cas, store.OperatorOwner{}, now)
+		jgExec(t, r.path, `INSERT INTO operator_runner_state VALUES(1, 7, 'runner', 1, 'test', 1010, 'start', 1, NULL)`)
+		if err := r.ClaimOperatorAction(ctx, r.cas, store.OperatorOwner{Nonce: "runner", Generation: 7}, now); err != nil {
+			t.Fatal(err)
+		}
+		if err := r.PruneOperatorActions(ctx, now); err != nil {
+			t.Fatal(err)
+		}
+	})
+	if got != "" {
+		t.Fatalf("an operator store method wrote %q", got)
+	}
+}
--- a/host/authority/judge_operator_test.go
+++ b/host/authority/judge_operator_test.go
@@ -173,6 +173,7 @@
 
 // TestOperatorRunnerRechecks: the runner trusts no row digest, stored approver or daemon status.
 func TestOperatorRunnerRechecks(t *testing.T) {
+	t.Run("claimed store row", testOperatorClaimedRecheck)
 	in, ssh, pk := ojIn(t)
 	moved := `UPDATE webauthn_credentials SET sign_count = 2 WHERE name = 'mark-mac'`
 	for by, c := range map[string]struct {
--- /dev/null
+++ b/host/authority/operator_store_test.go
@@ -0,0 +1,249 @@
+package authority
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"errors"
+	"reflect"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"modernc.org/sqlite"
+)
+
+type operatorStoreRig struct {
+	path string
+	*store.Credentials
+	action   approval.OperatorAction
+	cas      store.OperatorCAS
+	proof    *VerifiedOperatorAction
+	artifact []byte
+}
+
+func newOperatorStoreRig(t *testing.T, passkey bool) operatorStoreRig {
+	t.Helper()
+	path, _ := jgRig(t)
+	ctx := boundedTestContext(t)
+	c := operatorStoreOpen(t, path, store.OpenCredentials)
+	canonical, err := store.CanonicalPath(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	a := ojAction()
+	a.Store = canonical
+	req := store.OperatorRequest{Action: a, OperationID: "1234567890abcdef1234567890abcdef", CreatedAt: 1000, ExpiresAt: 1600, Challenge: sha256.Sum256([]byte(t.Name()))}
+	row, err := c.RequestOperatorAction(ctx, req, 1000)
+	if err != nil {
+		t.Fatal(err)
+	}
+	in := OperatorJudgeInput{Action: a, Digest: [32]byte(row.Digest), BytesDigest: sha256.Sum256(row.Action),
+		Challenge: [32]byte(row.Challenge), CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt, Now: 1010}
+	var artifact []byte
+	if passkey {
+		artifact = ojAssert(t, jgPasskey, in, 2)
+	} else {
+		artifact = ojSSH(t, jgSigner, approval.OperatorActionNamespace, in)
+	}
+	ro := operatorStoreOpen(t, path, store.OpenCredentialsReadOnly)
+	proof, err := JudgeOperatorAction(ctx, ro, in, artifact)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return operatorStoreRig{path, c, a, store.OperatorCAS{ID: row.ID, Digest: in.Digest, Generation: row.Generation}, proof, artifact}
+}
+
+func operatorStoreOpen(t *testing.T, path string, open func(string) (*store.Credentials, error)) *store.Credentials {
+	t.Helper()
+	c, err := open(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { c.Close() })
+	return c
+}
+
+func (r operatorStoreRig) confirm(t *testing.T) store.OperatorActionRow {
+	t.Helper()
+	row, err := r.ConfirmOperatorAction(boundedTestContext(t), r.cas, r.artifact, r.proof, 1010)
+	if err != nil || row.State != "confirmed" || row.ApprovedBy != r.proof.ApprovedBy() {
+		t.Fatalf("confirm: %s %q %v", row.State, row.ApprovedBy, err)
+	}
+	return row
+}
+
+func TestOperatorConfirmReplay(t *testing.T) {
+	for _, passkey := range []bool{false, true} {
+		r := newOperatorStoreRig(t, passkey)
+		want := r.confirm(t)
+		if err := r.Close(); err != nil {
+			t.Fatal(err)
+		}
+		c := operatorStoreOpen(t, r.path, store.OpenCredentials)
+		ctx := boundedTestContext(t)
+		got, err := c.ConfirmOperatorAction(ctx, r.cas, r.artifact, r.proof, 1020)
+		if err != nil || !reflect.DeepEqual(got, want) {
+			t.Fatalf("confirm replay changed the row or refused: %v", err)
+		}
+		name := "mark-ssh"
+		if passkey {
+			name = "mark-mac"
+		}
+		a, ok, err := c.Approver(ctx, name)
+		if err != nil || !ok || a.LastUsedAt != 1010 {
+			t.Fatalf("replay changed approver use: %+v %v", a, err)
+		}
+		p, ok, err := c.Passkey(ctx, "mark-mac")
+		if err != nil || !ok || (passkey && p.SignCount != 2) || (!passkey && p.SignCount != 1) {
+			t.Fatalf("unexpected counter: %d %v", p.SignCount, err)
+		}
+	}
+}
+
+func TestOperatorConfirmCrash(t *testing.T) {
+	const drift = `CREATE TRIGGER drift BEFORE UPDATE OF last_used_at ON approvers BEGIN UPDATE operator_actions SET `
+	for _, tc := range []struct {
+		name, sql string
+		want      error
+		counter   int64
+	}{
+		{"update-fault", `CREATE TRIGGER crash_operator BEFORE UPDATE OF approved_by ON operator_actions
+		   BEGIN SELECT raise(abort, 'test: interrupted before state commit'); END`, nil, 1},
+		{"generation drift", drift + `generation = 9; END`, store.ErrOperatorConflict, 1},
+		{"digest drift", drift + `digest = zeroblob(32); END`, store.ErrOperatorConflict, 1},
+		{"expiry drift", drift + `expires_at = 1010; END`, store.ErrOperatorConflict, 1},
+		{"counter-race", `UPDATE webauthn_credentials SET sign_count = 2`, store.ErrSignCountRegression, 2},
+		{"revoked", `UPDATE approvers SET revoked_at = 1010 WHERE name = 'mark-mac'`, store.ErrApproverRevoked, 1},
+		{"changed-rp", `UPDATE webauthn_credentials SET rp_id = 'replacement.example'`, store.ErrApproverKeyChanged, 1},
+		{"changed-origin", `UPDATE webauthn_credentials SET origin = 'https://replacement.example'`, store.ErrApproverKeyChanged, 1},
+		{"changed-backup", `UPDATE webauthn_credentials SET backup_eligible = 1`, store.ErrApproverKeyChanged, 1},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			r := newOperatorStoreRig(t, true)
+			jgExec(t, r.path, tc.sql)
+			ctx := boundedTestContext(t)
+			_, err := r.ConfirmOperatorAction(ctx, r.cas, r.artifact, r.proof, 1010)
+			if tc.want == nil {
+				operatorSQLiteCode(t, err, 1811) // SQLITE_CONSTRAINT_TRIGGER
+			} else if !errors.Is(err, tc.want) {
+				t.Fatalf("confirmation fault: %v, want %v", err, tc.want)
+			}
+			r.assertUnconfirmed(t, tc.counter)
+		})
+	}
+}
+
+func operatorSQLiteCode(t *testing.T, err error, code int) {
+	t.Helper()
+	var se *sqlite.Error
+	if !errors.As(err, &se) || se.Code() != code {
+		t.Fatalf("confirmation error %v, want SQLite code %d", err, code)
+	}
+}
+
+func TestOperatorConfirmReadOnly(t *testing.T) {
+	r := newOperatorStoreRig(t, true)
+	ro := operatorStoreOpen(t, r.path, store.OpenCredentialsReadOnly)
+	_, err := ro.ConfirmOperatorAction(boundedTestContext(t), r.cas, r.artifact, r.proof, 1010)
+	operatorSQLiteCode(t, err, 8) // SQLITE_READONLY
+	r.assertUnconfirmed(t, 1)
+}
+
+func (r operatorStoreRig) assertUnconfirmed(t *testing.T, count int64) {
+	t.Helper()
+	ctx := boundedTestContext(t)
+	row, err := r.OperatorAction(ctx, r.cas.ID, 1010)
+	if err != nil || row.State != "requested" || row.ApprovedBy != "" || row.Artifact != nil || row.Generation != r.cas.Generation || !bytes.Equal(row.Digest, r.cas.Digest[:]) || row.ExpiresAt != 1600 {
+		t.Fatalf("failure changed request: %s %q %v", row.State, row.ApprovedBy, err)
+	}
+	p, ok, err := r.Passkey(ctx, "mark-mac")
+	if err != nil || !ok || p.SignCount != count {
+		t.Fatalf("failure changed counter: %d %v", p.SignCount, err)
+	}
+	a, ok, err := r.Approver(ctx, "mark-mac")
+	if err != nil || !ok || a.LastUsedAt != 0 {
+		t.Fatalf("failure changed approver use: %d %v", a.LastUsedAt, err)
+	}
+}
+
+func TestOperatorConfirmStoreBinding(t *testing.T) {
+	r := newOperatorStoreRig(t, true)
+	for _, tc := range []struct {
+		name     string
+		cas      store.OperatorCAS
+		artifact []byte
+		proof    *VerifiedOperatorAction
+		now      int64
+		want     error
+	}{
+		{"id", store.OperatorCAS{ID: "missing", Digest: r.cas.Digest, Generation: r.cas.Generation}, r.artifact, r.proof, 1010, store.ErrOperatorMissing},
+		{"digest", store.OperatorCAS{ID: r.cas.ID, Generation: r.cas.Generation}, r.artifact, r.proof, 1600, store.ErrOperatorConflict},
+		{"generation", store.OperatorCAS{ID: r.cas.ID, Digest: r.cas.Digest, Generation: r.cas.Generation + 1}, r.artifact, nil, 1010, store.ErrOperatorConflict},
+		{"nil", r.cas, r.artifact, nil, 1010, approval.ErrNotVerified},
+		{"artifact substitution", r.cas, bytes.Clone(r.artifact[:len(r.artifact)-1]), r.proof, 1010, approval.ErrNotVerified},
+		{"at expiry", r.cas, r.artifact, r.proof, 1600, store.ErrOperatorExpired},
+		{"clock behind", r.cas, r.artifact, r.proof, 999, store.ErrOperatorConflict},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			if _, err := r.ConfirmOperatorAction(boundedTestContext(t), tc.cas, tc.artifact, tc.proof, tc.now); !errors.Is(err, tc.want) {
+				t.Fatalf("confirmation: %v, want %v", err, tc.want)
+			}
+			r.assertUnconfirmed(t, 1)
+		})
+	}
+	foreign := newOperatorStoreRig(t, true)
+	_, err := r.ConfirmOperatorAction(boundedTestContext(t), r.cas, foreign.artifact, foreign.proof, 1010)
+	if binding, ok := errors.AsType[*approval.BindingError](err); !ok || binding.Field != "digest" {
+		t.Fatalf("proof for different store: %v", err)
+	}
+	r.assertUnconfirmed(t, 1)
+
+	if row, err := r.ConfirmOperatorAction(boundedTestContext(t), r.cas, r.artifact, r.proof, 1000); err != nil || row.State != "confirmed" {
+		t.Fatalf("confirm at creation: %s %v", row.State, err)
+	}
+}
+
+func testOperatorClaimedRecheck(t *testing.T) {
+	r := newOperatorStoreRig(t, true)
+	confirmed := r.confirm(t)
+	ctx := boundedTestContext(t)
+	jgExec(t, r.path, `INSERT INTO operator_runner_state VALUES(1, 7, 'runner', 1, 'test', 1010, 'start', 1, NULL)`)
+	if err := r.ClaimOperatorAction(ctx, r.cas, store.OperatorOwner{Nonce: "runner", Generation: 7}, 1010); err != nil {
+		t.Fatal(err)
+	}
+	row, err := r.OperatorAction(ctx, r.cas.ID, 1010)
+	if err != nil || row.State != "claimed" {
+		t.Fatalf("claim read: %s %v", row.State, err)
+	}
+	ro := operatorStoreOpen(t, r.path, store.OpenCredentialsReadOnly)
+	challenge := [32]byte(confirmed.Challenge)
+	check := func(a approval.OperatorAction, now int64, row store.OperatorActionRow) (*VerifiedOperatorAction, error) {
+		return RecheckClaimedOperatorAction(ctx, ro, a, challenge, confirmed.ExpiresAt, now, row)
+	}
+	// The row may disagree with the runner's frozen evidence.
+	row.Digest, row.Challenge, row.ExpiresAt = nil, nil, 1601
+	proof, err := check(r.action, 1599, row)
+	if err != nil || !proof.Rechecked() || proof.ApprovedBy() != r.proof.ApprovedBy() {
+		t.Fatalf("claimed receipt proof or exact approver lost: %v", err)
+	}
+	if _, err := proof.OperatorConfirmation().ConfirmationFor(row.Artifact); !errors.Is(err, approval.ErrNotVerified) {
+		t.Fatalf("runner proof was accepted as an initial confirmation: %v", err)
+	}
+	wrong := r.action
+	wrong.Task += "changed"
+	if proof, err := check(wrong, 1599, row); proof != nil || err == nil {
+		t.Fatal("row supplied authority for a different frozen intent")
+	}
+	if proof, err := check(r.action, 1600, row); proof != nil || err == nil {
+		t.Fatal("row supplied authority after recorded expiry")
+	}
+	row.ApprovedBy = "webauthn:replacement"
+	if proof, err := check(r.action, 1599, row); proof != nil || !errors.Is(err, ErrOperatorApprover) {
+		t.Fatalf("trusted row identity: %v", err)
+	}
+	row.ApprovedBy = confirmed.ApprovedBy
+	row.Artifact = nil
+	if proof, err := check(r.action, 1599, row); proof != nil || !errors.Is(err, ErrUnknownArtifact) {
+		t.Fatalf("claimed state supplied authority without an artifact: %v", err)
+	}
+}
--- /dev/null
+++ b/host/store/operator_cas.go
@@ -0,0 +1,166 @@
+package store
+
+import (
+	"bytes"
+	"context"
+	"crypto/sha256"
+	"database/sql"
+	"errors"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+var ErrOperatorRunner = errors.New("store: operator runner is absent, stale or replaced")
+
+type OperatorCAS struct {
+	ID         string
+	Digest     [32]byte
+	Generation int64
+}
+
+type OperatorOwner struct {
+	Nonce      string
+	Generation int64
+}
+
+// The interface transports a sealed law proof. A method implemented by a caller
+// cannot supply the verified binding or approver-use fields as plain data.
+type operatorConfirmation interface {
+	OperatorConfirmation() *approval.SignedOperatorAction
+}
+
+func operatorCASRow(ctx context.Context, tx *sql.Tx, cas OperatorCAS) (OperatorActionRow, error) {
+	row, err := scanOperatorAction(tx.QueryRowContext(ctx, `SELECT `+operatorColumns+` FROM operator_actions WHERE id = ?`, cas.ID))
+	if err == nil && (!bytes.Equal(row.Digest, cas.Digest[:]) || row.Generation != cas.Generation) {
+		return OperatorActionRow{}, ErrOperatorConflict
+	}
+	return row, err
+}
+
+// ConfirmOperatorAction consumes an authority proof and persists its exact
+// artifact and approver identity atomically with approver use and counter.
+func (c *Credentials) ConfirmOperatorAction(ctx context.Context, cas OperatorCAS, artifact []byte, proof operatorConfirmation, now int64) (OperatorActionRow, error) {
+	if err := c.operatorAccess(ctx); err != nil {
+		return OperatorActionRow{}, err
+	}
+	tx, err := c.db.BeginTx(ctx, nil)
+	if err != nil {
+		return OperatorActionRow{}, err
+	}
+	defer tx.Rollback()
+	row, err := scanOperatorAction(tx.QueryRowContext(ctx, `SELECT `+operatorColumns+` FROM operator_actions WHERE id = ?`, cas.ID))
+	if err != nil {
+		return OperatorActionRow{}, err
+	}
+	if !bytes.Equal(row.Digest, cas.Digest[:]) {
+		return OperatorActionRow{}, ErrOperatorConflict
+	}
+	if row.ApprovedBy != "" {
+		return row, nil
+	}
+	if row.State != "requested" || row.Generation != cas.Generation {
+		return OperatorActionRow{}, ErrOperatorConflict
+	}
+	if now >= row.ExpiresAt {
+		return OperatorActionRow{}, ErrOperatorExpired
+	}
+	if now < row.CreatedAt || sha256.Sum256(row.Action) != cas.Digest || len(row.Challenge) != 32 {
+		return OperatorActionRow{}, ErrOperatorConflict
+	}
+	a, err := approval.DecodeOperatorAction(row.Action)
+	if err != nil {
+		return OperatorActionRow{}, err
+	}
+	if proof == nil {
+		return OperatorActionRow{}, approval.ErrNotVerified
+	}
+	signed := proof.OperatorConfirmation()
+	if err := signed.Binds(a, [32]byte(row.Challenge), row.ExpiresAt); err != nil {
+		return OperatorActionRow{}, err
+	}
+	binding, err := signed.ConfirmationFor(artifact)
+	if err != nil || len(artifact) > MaxSurfaceApprovalArtifactBytes {
+		return OperatorActionRow{}, approval.ErrNotVerified
+	}
+	if err := useApprover(ctx, tx, ApproverUse{Name: binding.Name, Format: binding.Format, PublicKey: []byte(binding.PublicKey)}, now); err != nil {
+		return OperatorActionRow{}, err
+	}
+	if id, count, ok := signed.PasskeyUse(); ok {
+		var rp, origin string
+		var backup bool
+		if err := tx.QueryRowContext(ctx, `SELECT rp_id, origin, backup_eligible FROM webauthn_credentials WHERE name = ?`, binding.Name).Scan(&rp, &origin, &backup); err != nil {
+			return OperatorActionRow{}, err
+		}
+		if rp != binding.RPID || origin != binding.Origin || backup != binding.BackupEligible {
+			return OperatorActionRow{}, ErrApproverKeyChanged
+		}
+		if err := usePasskey(ctx, tx, binding.Name, PasskeyUse{CredentialID: id, SignCount: int64(count)}); err != nil {
+			return OperatorActionRow{}, err
+		}
+	}
+	res, err := tx.ExecContext(ctx, `UPDATE operator_actions SET state = 'confirmed', approved_by = ?, artifact = ?
+		WHERE id = ? AND digest = ? AND generation = ? AND state = 'requested' AND expires_at > ?`,
+		signed.ApprovedBy(), artifact, cas.ID, cas.Digest[:], cas.Generation, now)
+	if err := operatorChangedOne(res, err); err != nil {
+		return OperatorActionRow{}, err
+	}
+	if err := tx.Commit(); err != nil {
+		return OperatorActionRow{}, err
+	}
+	row.State, row.ApprovedBy, row.Artifact = "confirmed", signed.ApprovedBy(), bytes.Clone(artifact)
+	return row, nil
+}
+
+// ClaimOperatorAction arbitrates one fresh claim. Runner effect readiness and
+// the shared kernel episode lock are prerequisites of the future runner.
+func (c *Credentials) ClaimOperatorAction(ctx context.Context, cas OperatorCAS, runner OperatorOwner, now int64) error {
+	if err := c.operatorAccess(ctx); err != nil {
+		return err
+	}
+	tx, err := c.db.BeginTx(ctx, nil)
+	if err != nil {
+		return err
+	}
+	defer tx.Rollback()
+	row, err := operatorCASRow(ctx, tx, cas)
+	if err != nil {
+		return err
+	}
+	if row.State != "confirmed" || row.Owner != "" || now < row.CreatedAt {
+		return ErrOperatorConflict
+	}
+	if now >= row.ExpiresAt {
+		return ErrOperatorExpired
+	}
+	var current bool
+	if err := tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM operator_runner_state WHERE
+		generation = ? AND owner = ? AND seen_at > ? AND seen_at <= ? AND (readiness = 'both' OR readiness = ?))`,
+		runner.Generation, runner.Nonce, now-10, now, row.Kind).Scan(&current); err != nil {
+		return err
+	}
+	if !current {
+		return ErrOperatorRunner
+	}
+	res, err := tx.ExecContext(ctx, `UPDATE operator_actions AS a SET state = 'claimed', owner = ?, generation = generation + 1
+		WHERE id = ? AND digest = ? AND generation = ? AND state = 'confirmed' AND expires_at > ? AND owner IS NULL
+		AND NOT EXISTS (SELECT 1 FROM operator_actions b WHERE b.episode_id = a.episode_id AND b.id <> a.id AND b.owner IS NOT NULL)`,
+		runner.Nonce, cas.ID, cas.Digest[:], cas.Generation, now)
+	if err := operatorChangedOne(res, err); err != nil {
+		return err
+	}
+	return tx.Commit()
+}
+
+func operatorChangedOne(result sql.Result, err error) error {
+	if err != nil {
+		return err
+	}
+	n, err := result.RowsAffected()
+	if err != nil {
+		return err
+	}
+	if n != 1 {
+		return ErrOperatorConflict
+	}
+	return nil
+}
--- /dev/null
+++ b/host/store/operator_cas_test.go
@@ -0,0 +1,123 @@
+package store
+
+import (
+	"errors"
+	"sync"
+	"testing"
+)
+
+func operatorClaimFixture(t *testing.T, c *Credentials) (OperatorCAS, OperatorOwner) {
+	t.Helper()
+	row := putOperator(t, c, operatorRequest(t, c, 1, "claim", "start"), 1000)
+	// This store-only fixture exercises the SQL claim. The actual proof and
+	// signed artifact are exercised end-to-end in host/authority.
+	operatorExec(t, c, `UPDATE operator_actions SET state = 'confirmed', approved_by = 'ssh:fixture', artifact = action`)
+	operatorExec(t, c, `INSERT INTO operator_runner_state VALUES(1, 7, 'runner', 1, 'test', 1590, 'start', 1, NULL)`)
+	return OperatorCAS{ID: row.ID, Digest: [32]byte(row.Digest), Generation: row.Generation}, OperatorOwner{Nonce: "runner", Generation: 7}
+}
+
+func TestOperatorClaimBeforeExpiry(t *testing.T) {
+	c := openTestCredentials(t)
+	cas, runner := operatorClaimFixture(t, c)
+	ctx := testCtx(t)
+	for _, tc := range []struct {
+		name   string
+		cas    OperatorCAS
+		runner OperatorOwner
+		now    int64
+		want   error
+	}{
+		{"exact expiry", cas, runner, 1600, ErrOperatorExpired},
+		{"past expiry", cas, runner, 1601, ErrOperatorExpired},
+		{"clock behind", cas, runner, 999, ErrOperatorConflict},
+		{"future heartbeat", cas, runner, 1589, ErrOperatorRunner},
+		{"readiness", cas, runner, 1599, ErrOperatorRunner},
+		{"past heartbeat", cas, runner, 1599, ErrOperatorRunner},
+		{"owner", cas, OperatorOwner{Nonce: "other", Generation: 7}, 1590, ErrOperatorRunner},
+		{"runner generation", cas, OperatorOwner{Nonce: "runner", Generation: 6}, 1590, ErrOperatorRunner},
+		{"action generation", OperatorCAS{ID: cas.ID, Digest: cas.Digest, Generation: 2}, OperatorOwner{}, 1590, ErrOperatorConflict},
+		{"digest", OperatorCAS{ID: cas.ID, Generation: 1}, OperatorOwner{}, 1590, ErrOperatorConflict},
+		{"id", OperatorCAS{ID: "missing", Digest: cas.Digest, Generation: 1}, runner, 1590, ErrOperatorMissing},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			if tc.name == "readiness" {
+				operatorExec(t, c, `UPDATE operator_runner_state SET readiness = 'none'`)
+				defer operatorExec(t, c, `UPDATE operator_runner_state SET readiness = 'start'`)
+			}
+			if tc.name == "past heartbeat" {
+				operatorExec(t, c, `UPDATE operator_runner_state SET seen_at = 1589`)
+				defer operatorExec(t, c, `UPDATE operator_runner_state SET seen_at = 1590`)
+			}
+			if err := c.ClaimOperatorAction(ctx, tc.cas, tc.runner, tc.now); !errors.Is(err, tc.want) {
+				t.Fatalf("claim: %v want %v", err, tc.want)
+			}
+			got := getOperator(t, c, cas.ID, 1599)
+			if got.State != "confirmed" || got.Owner != "" || got.Generation != cas.Generation {
+				t.Fatalf("claim refusal changed ownership: %+v", got)
+			}
+			operatorCounts(t, c, 1, 1)
+		})
+	}
+	operatorExec(t, c, `UPDATE operator_actions SET state = 'cancelled', terminal_at = 1590`)
+	if err := c.ClaimOperatorAction(ctx, cas, OperatorOwner{}, 1599); !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("cancelled: %v", err)
+	}
+	operatorCounts(t, c, 1, 1)
+	if got := getOperator(t, c, cas.ID, 1599); got.State != "cancelled" || got.Owner != "" {
+		t.Fatal("cancel won without retaining state")
+	}
+}
+
+func TestOperatorClaimRace(t *testing.T) {
+	c := openTestCredentials(t)
+	cas, runner := operatorClaimFixture(t, c)
+	operatorExec(t, c, `UPDATE operator_actions SET created_at = 1599; UPDATE operator_runner_state SET seen_at = 1599`)
+	const contenders = 2
+	var ready sync.WaitGroup
+	ready.Add(contenders)
+	start, done := make(chan struct{}), make(chan error, contenders)
+	ctx := testCtx(t)
+	for range contenders {
+		go func() {
+			ready.Done()
+			<-start
+			done <- c.ClaimOperatorAction(ctx, cas, runner, 1599)
+		}()
+	}
+	ready.Wait()
+	close(start)
+	a, b := <-done, <-done
+	if a != nil {
+		a, b = b, a
+	}
+	if a != nil || !errors.Is(b, ErrOperatorConflict) {
+		t.Fatalf("race: %v / %v", a, b)
+	}
+	got := getOperator(t, c, cas.ID, 1599)
+	if got.State != "claimed" || got.Owner != runner.Nonce || got.Generation != cas.Generation+1 {
+		t.Fatalf("claim result: %+v", got)
+	}
+	operatorCounts(t, c, 1, 1)
+	// Simulate corruption of the active-episode index. The CAS's owned-
+	// sibling predicate must still refuse a second claim.
+	operatorExec(t, c, `DROP INDEX operator_active_episode`)
+	sibling := operatorRequest(t, c, 2, "claim", "finish")
+	operatorExec(t, c, `INSERT INTO operator_actions
+		(id, episode_id, kind, review, digest, action, challenge, created_at, expires_at, state, generation, approved_by, artifact)
+		SELECT ?, episode_id, 'finish', review, digest, action, challenge, created_at, expires_at,
+			'confirmed', 1, approved_by, artifact FROM operator_actions WHERE id = ?`, sibling.Action.ID, cas.ID)
+	operatorExec(t, c, `UPDATE operator_runner_state SET readiness = 'both'`)
+	cas.ID = sibling.Action.ID
+	if err := c.ClaimOperatorAction(ctx, cas, runner, 1599); !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("claim alongside owned sibling: %v", err)
+	}
+	if second := getOperator(t, c, cas.ID, 1599); second.State != "confirmed" || second.Owner != "" {
+		t.Fatalf("second claim changed state: %+v", second)
+	}
+
+	operatorCounts(t, c, 2, 1)
+	operatorExec(t, c, `UPDATE operator_actions SET kind = 'finish', v3_id = 'vr-0000000000000001' WHERE state = 'claimed'`)
+	if _, err := c.db.Exec(`UPDATE operator_actions SET v3_id = 'vr-0000000000000001' WHERE state = 'confirmed'`); err == nil {
+		t.Fatal("the same V3 card was linked to two actions")
+	}
+}
````

### Patch S3b

````diff S3b
--- a/host/store/operator_cas_test.go
+++ b/host/store/operator_cas_test.go
@@ -1,7 +1,9 @@
 package store
 
 import (
+	"bytes"
 	"errors"
+	"path/filepath"
 	"sync"
 	"testing"
 )
@@ -121,3 +123,61 @@
 		t.Fatal("the same V3 card was linked to two actions")
 	}
 }
+
+func TestOperatorStepRestart(t *testing.T) {
+	path := filepath.Join(t.TempDir(), CredentialsFileName)
+	open := func() *Credentials {
+		c, err := OpenCredentials(path)
+		if err != nil {
+			t.Fatal(err)
+		}
+		t.Cleanup(func() { c.Close() })
+		return c
+	}
+	c := open()
+	cas, runner := operatorClaimFixture(t, c)
+	ctx := testCtx(t)
+	if err := c.ClaimOperatorAction(ctx, cas, runner, 1599); err != nil {
+		t.Fatal(err)
+	}
+	cas.Generation++
+	key := OperatorStepKey{Action: cas, Owner: runner.Nonce, Kind: "launch", Identity: "run-one", Version: 1}
+	if err := c.ReserveOperatorStep(ctx, key, 1600); !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("reserve at expiry: %v", err)
+	}
+	operatorCounts(t, c, 1, 1)
+	if err := c.ReserveOperatorStep(ctx, key, 1599); err != nil {
+		t.Fatal(err)
+	}
+	if err := c.Close(); err != nil {
+		t.Fatal(err)
+	}
+	c = open()
+	if err := c.ReserveOperatorStep(ctx, key, 1599); !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("reservation replay: %v", err)
+	}
+	reservation, err := c.OperatorStep(ctx, key)
+	if err != nil || reservation.Phase != "reserved" || reservation.Version != 1 || reservation.Owner != runner.Nonce || reservation.Receipt != nil {
+		t.Fatalf("lost reservation on restart: %+v %v", reservation, err)
+	}
+	receipt := []byte("observed-test-launch")
+	wrong := key
+	wrong.Version++
+	if err := c.ObserveOperatorStep(ctx, wrong, receipt); !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("wrong step CAS: %v", err)
+	}
+	if got, err := c.OperatorStep(ctx, key); err != nil || got.Phase != "reserved" || got.Receipt != nil {
+		t.Fatal("failed CAS changed reservation")
+	}
+	if err := c.ObserveOperatorStep(ctx, key, receipt); err != nil {
+		t.Fatal(err)
+	}
+	if err := c.ObserveOperatorStep(ctx, key, receipt); !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("observation replay: %v", err)
+	}
+	got, err := c.OperatorStep(ctx, key)
+	if err != nil || got.Phase != "observed" || got.Version != 2 || !bytes.Equal(got.Receipt, receipt) {
+		t.Fatalf("observed receipt: %+v %v", got, err)
+	}
+	operatorCounts(t, c, 1, 2)
+}
--- /dev/null
+++ b/host/store/operator_steps.go
@@ -0,0 +1,72 @@
+package store
+
+import (
+	"context"
+	"database/sql"
+	"errors"
+)
+
+type OperatorStepKey struct {
+	Action                OperatorCAS
+	Owner, Kind, Identity string
+	Version               int64
+}
+
+// ReserveOperatorStep records identity before a future caller crosses an
+// effect boundary. The caller must first recheck its own frozen proof.
+func (c *Credentials) ReserveOperatorStep(ctx context.Context, key OperatorStepKey, now int64) error {
+	if err := c.operatorAccess(ctx); err != nil {
+		return err
+	}
+	if key.Kind == "operation" || key.Identity == "" || key.Owner == "" || key.Version != 1 {
+		return ErrOperatorInput
+	}
+	result, err := c.db.ExecContext(ctx, `INSERT INTO operator_action_steps
+		(kind, identity, action_id, phase, version, owner, created_at)
+		SELECT ?, ?, id, 'reserved', ?, owner, ? FROM operator_actions
+		WHERE id = ? AND digest = ? AND generation = ? AND owner = ?
+		AND state IN ('claimed', 'running') AND expires_at > ?
+		AND (SELECT count(*) FROM operator_action_steps WHERE action_id = ? AND kind <> 'operation') < 32
+		ON CONFLICT(kind, identity) DO NOTHING`, key.Kind, key.Identity, key.Version, now,
+		key.Action.ID, key.Action.Digest[:], key.Action.Generation, key.Owner, now, key.Action.ID)
+	return operatorChangedOne(result, err)
+}
+
+// ObserveOperatorStep records the outcome of a previously reserved effect.
+// Expiry prevents new reservations but does not prevent exact observation.
+func (c *Credentials) ObserveOperatorStep(ctx context.Context, key OperatorStepKey, receipt []byte) error {
+	if err := c.operatorAccess(ctx); err != nil {
+		return err
+	}
+	if len(receipt) == 0 || len(receipt) > MaxSurfaceApprovalArtifactBytes {
+		return ErrOperatorInput
+	}
+	result, err := c.db.ExecContext(ctx, `UPDATE operator_action_steps
+		SET phase = 'observed', version = version + 1, receipt = ?
+		WHERE action_id = ? AND kind = ? AND identity = ? AND version = ? AND owner = ? AND phase = 'reserved'
+		AND EXISTS (SELECT 1 FROM operator_actions WHERE id = ? AND digest = ? AND generation = ? AND owner = ?)`,
+		receipt, key.Action.ID, key.Kind, key.Identity, key.Version, key.Owner,
+		key.Action.ID, key.Action.Digest[:], key.Action.Generation, key.Owner)
+	return operatorChangedOne(result, err)
+}
+
+type OperatorStepReceipt struct {
+	Phase, Owner string
+	Version      int64
+	Receipt      []byte
+}
+
+// OperatorStep reads a reservation or recorded observation across a restart.
+func (c *Credentials) OperatorStep(ctx context.Context, key OperatorStepKey) (OperatorStepReceipt, error) {
+	if err := c.operatorAccess(ctx); err != nil {
+		return OperatorStepReceipt{}, err
+	}
+	var step OperatorStepReceipt
+	err := c.db.QueryRowContext(ctx, `SELECT phase, coalesce(owner,''), version, receipt FROM operator_action_steps
+		WHERE action_id = ? AND kind = ? AND identity = ?`, key.Action.ID, key.Kind, key.Identity).
+		Scan(&step.Phase, &step.Owner, &step.Version, &step.Receipt)
+	if errors.Is(err, sql.ErrNoRows) {
+		err = ErrOperatorMissing
+	}
+	return step, err
+}
````
