# Sprint plan — w-world-address-bar-verbs-v4, sprint S4b (iteration 430, row 243)

> **Scope.** M1 adds terminal episode exclusion, including preset and resume. M2 gates companion startup on a newer compatible lease and extends the landed home-upgrade refusal test. M3 pins the carried service and store guards. Cumulative product changes are **638 → 878 → 978** against `c42fe698`.
>
> **Authority.** Reviewed design r1 §§2.2, 3.1–3.2 and 5.5. Quorum r2 PROCEED 3/3, independent audit PASS 81. D-WORLD-119=A and D-WORLD-120=A were attended in `82f05004`. S1 #606, S2 #609, S2b #611, S3 #643 and S4 #651 are landed. S4 merge `ee8ef1c0` received FAIL 72 → PASS 81 → PASS 81.

## 0. Binding constraints

- Apply M1 through M3 in order. Product changes come exclusively from Appendix P through `/usr/bin/patch -p1 -E`. The controller owns Git writes.
- D-WORLD-109 governs fixture keys and stores. Use real disposable SQLite migration, repository test signatures, fake HTTP/controller/terminal, and files under `t.TempDir()` or mission scratch. Observe operator files and supervisor jobs read-only.
- Preserve the F1 terminal fence, signed action ownership, attended service confirmation, schema rollback guards and runner inertness. Dispatch depends on the paired admission, drain and recovery gates in Appendix C.

## 1. Pre-flight

### 1.1 Baseline at c42fe698

`git archive c42fe698 | tar -x -C "$S/base"` under `$HOME/.ailang/state/world-iter430/plan-scratch`. The first pass ran the S4 §6 script with its base ref updated before any stage edits. Every Go leg used the pinned `AILANG_BIN`. The initial script displayed a skipped `verify_go` leg with rc 0. The leg was intentionally skipped and provided no pass evidence. Output from the original baseline follows.

````text
LEG observation-before rc=77 0s UNINFORMATIVE UNDER SANDBOX
LEG paths rc=0 0s
TOUCHED_PACKAGES: 
LEG touched rc=0 0s
LEG status rc=0 0s
LEG gofmt rc=0 0s
LEG build rc=0 2s
LEG vet rc=0 3s
LEG fence rc=0 25s
LEG verify_ail rc=0 42s
LEG pkg-cmd-ailang-worldd rc=1 2s UNINFORMATIVE UNDER SANDBOX
LEG race-cmd-ailang-worldd rc=1 20s UNINFORMATIVE UNDER SANDBOX
LEG pkg-host-store rc=0 13s
LEG race-host-store rc=0 37s
LEG inv-host-approval rc=0 0s
LEG inv-host-approval-approvaltest rc=0 1s
LEG inv-host-authority rc=0 0s
LEG inv-host-boundary rc=0 1s
LEG inv-host-daemon rc=0 1s
LEG inv-host-verbs rc=0 1s
LEG inv-cmd-ailang-worldd rc=0 0s
LEG inv-host-store rc=0 1s
LEG routes rc=0 0s
LEG scope rc=0 0s
LEG callers rc=0 0s
LEG verify_go rc=0 0s
LEG observation-after rc=77 0s UNINFORMATIVE UNDER SANDBOX
GATES_DONE 149s FAILS: none UNINFORMATIVE: observation-before pkg-cmd-ailang-worldd race-cmd-ailang-worldd observation-after
````

The §6 script with explicit skip output was also run on the base archive. Generated compile caches were removed from the exports before measuring patches. Exit 77 means the listed sandbox legs require a controller run.

````text
LEG observation-before rc=77 0s UNINFORMATIVE UNDER SANDBOX
LEG paths rc=0 0s
TOUCHED_PACKAGES: 
LEG touched rc=0 0s
LEG status rc=0 0s
LEG gofmt rc=0 0s
LEG build rc=0 1s
LEG vet rc=0 0s
LEG fence rc=0 1s
LEG verify_ail rc=0 49s
LEG pkg-cmd-ailang-worldd rc=1 3s UNINFORMATIVE UNDER SANDBOX
LEG race-cmd-ailang-worldd rc=1 14s UNINFORMATIVE UNDER SANDBOX
LEG pkg-host-store rc=0 15s
LEG race-host-store rc=0 36s
LEG inv-host-approval rc=0 0s
LEG inv-host-approval-approvaltest rc=0 1s
LEG inv-host-authority rc=0 1s
LEG inv-host-boundary rc=0 1s
LEG inv-host-daemon rc=0 1s
LEG inv-host-verbs rc=0 1s
LEG inv-cmd-ailang-worldd rc=0 0s
LEG inv-host-store rc=0 1s
LEG routes rc=0 0s
LEG scope rc=0 0s
LEG callers rc=0 1s
LEG verify_go SKIPPED 0s (GATES_SKIP_VERIFY_GO=1)
LEG observation-after rc=77 0s UNINFORMATIVE UNDER SANDBOX
GATES_DONE 126s FAILS: none UNINFORMATIVE: observation-before pkg-cmd-ailang-worldd race-cmd-ailang-worldd observation-after
````

### 1.2 Codebase claims checked

Run the commands here at `c42fe69888ce0c7f37a2ef967738085b1a27c1e6`, with `S=$HOME/.ailang/state/world-iter430/plan-scratch`.

| ID | Command | Observed at HEAD |
|---|---|---|
| C1 | `git rev-parse HEAD` and `git show -s --format='%h %s' ee8ef1c0 82f05004` | `c42fe69888ce0c7f37a2ef967738085b1a27c1e6`. `ee8ef1c0` merges #651. `82f05004` records the two attended A rulings. The three reports in `design_docs/verification/world-iter429/` match the three controller copies. |
| C2 | `git log --oneline 22473391..c42fe698` and `git diff --stat 22473391 c42fe698 -- cmd host` | 12 files, 1,374 additions and 30 deletions. #646–#649 alter preset and profile setup, coordinator, and allowance files. `git diff --name-only 22473391 c42fe698 -- 'cmd/ailang-worldd/episode*.go' 'cmd/ailang-worldd/operator_*.go' 'cmd/ailang-worldd/service*.go' 'host/store/operator*.go'` prints only `operator_service.go`, `operator_service_test.go` and `host/store/operator_runner_test.go`. M5 accounts for the touched production hunk in `runnerPath`. |
| C3 | `rg -n 'terminalEpisodeLock\|operatorEpisodeLockPath\|OperatorOwnedEpisode' cmd host` and `rg -n 'checkDaemon\|checkPresetDaemon\|launchEpisode\(\|fetchReport\(' cmd/ailang-worldd/episode*.go` | First query returns zero. M1 work is at the real calls in `episode.go:361`, `episode_preset.go:279,311`, `episode_resume.go:172`, `episode_launch.go:253,258,277,424` and `episode_finish.go:125`. Start, preset, resume each fence before prompt/provisioning, launch probes an existing descriptor before report and creates a missing one before the PID/child, finish after F1 and store health binding and before report. Two callers of `launchEpisode` are the public launch and post-mint preset start. |
| C4 | `rg -n 'func (runnerKernelLock\|runOperatorRunner\|runnerDB\|budgetContext\|credentialsFor)\|func \(e serviceEnv\) (runnerTick\|runnerStatus\|checkRunnerPair\|changeRunner)' cmd/ailang-worldd` | `operator_runner.go:21,38,76,98,150` contains flock/start/tick/status/DSN, `cli.go:90` context budget, `session.go:121` canonical credential selection, `operator_service.go:177,388` pair and lifecycle. `service.go:73,87,208,690` has env, real env, health and loopback. The startup hook point is `operator_service.go:430–442` after pair validation. |
| C5 | `rg -n 'func (operatorService\|newServiceFixture\|newFinishFixture\|newPresetFixture\|newEpisodeRun\|resumeFixture\|newLaunchFixture\|credSQL)\(' cmd/ailang-worldd` and `rg -n 'func .*OperatorHeartbeat\|func requireDeadline\|func readOnlyDSN\|func CanonicalPath' host/store` | `operator_service_test.go:32` builds the fake, `service_test.go:75` its service, `episode_finish_test.go:58` F1, `preset_test.go:44` preset, `episode_test.go:85` plain start, `episode_resume_test.go:16` resume, `episode_launch_test.go:26` launch, `episode_finish_surface_test.go:75` fixture SQL. Store `operator_runner.go:21,28,43`, `store.go:250`, `writer_lock.go:269`, `credentials.go:185` supply heartbeat, deadline and canonical read DSN. `operator_schema.sql:1–32` ties all durable ownership to `owner IS NOT NULL` and at most one active action per episode. |
| C6 | `git show --stat 976f4028` and `rg -n 'TestOperatorServiceUpgradePairFence\|TestOperatorServiceCancelConfirmation\|TestRunnerPathRefusalClasses\|TestOperatorCancelConfirmed' cmd host` | M4 +110 across the two test files. Existing `operator_service_test.go:280,332`, `host/store/operator_runner_test.go:146` pin four upgrade arms and both cancel checks. M5 `operator_service_test.go:174` already pins type-before-mode. M2 adds the query-error arm to the existing upgrade test. `service_upgrade.go:195–214` holds the lifecycle lock across installation/supervisor validation and the swap path. |
| C7 | `rg -n 'func Test(DaemonDependencyAllowlist\|WorldBoundaryDependencyAllowlist\|EveryCommandDirectoryIsAProtectedGroup\|ProductionGoSurface\|ProductionContextRoots\|StoreSurfaceTakesContext\|CLICredentialWriterSites\|DaemonCredentialWriterSites)' cmd host` | `host/daemon/daemon_test.go:899`, `host/boundary/allowlist_world_test.go:758,1121`, `host/store/context_roots_test.go:129,234`, `host/store/deadline_guard_test.go:43`, `cmd/ailang-worldd/cli_credential_writers_test.go:18` and `host/daemon/surface_approve_post_guard_test.go` cover the inventories. No allow-list or context-root addition applies. `TestCredentialsMethodsRequireADeadline` gains the three bounded reader calls. No new production dependency edges out of CLI/store, new subpackage, credential write site or `*Store` method. |
| C8 | `rg -n 'func Test(EveryVerbHasHelp\|VerbTableMatchesTheCLISwitch\|OperatorFlagParity\|QuickstartServiceSectionMatchesTheCLIHelp\|HelpIsNeverBlockedByTheRegistryKeyGuard\|PluginCommandsMatchTheCLIHelp)' cmd/ailang-worldd` and `rg -n 'syscall.Flock.*LOCK_NB' cmd host --glob '!**/*_test.go'` | Existing help/parity, QUICKSTART and plugin-help goldens keep the same expectations. Six production flock sites at base and final, two in `service.go`, two in `writer_lock_unix.go`, one each in `operator_runner.go` and `operator_service.go`. Terminal callers reuse `runnerKernelLock`. §6 extends caller inventory and runs behavioral exclusive/non-waiting tests. |
| C9 | `rg -n '\.(ClaimOperatorAction\|ReserveOperatorStep\|ObserveOperatorStep)\(' cmd host --glob '!**/*_test.go'` and `rg -n 'authority\.RecheckVerb\(' cmd host --glob '!**/*_test.go'` | Zero runner claim/reserve/observe production callers. V3 terminal `RecheckVerb` is a live positive control. Runner loop has a singleton lock and heartbeat publication only. §6 finds 50 registrations and 3 surface POSTs at base and M3. Route table, surface page and app compare byte-identical. The context scan finds 54 named roots and 324 production files at M3. |
| C10 | `rg -n 'addrGiven\|operator-runner' cmd/ailang-worldd/main.go` and `rg -n '0022\|O_NOFOLLOW\|se-tools\|world == creds\|h.Status\|net.ParseIP' cmd/ailang-worldd/operator_{runner,service}.go` | `main.go:382–387` already rejects global `--addr`; leading flag and zero-argument startup reach distinct observable paths. `operator_service.go:77,91,129,140,164,188` retains the mode, config open, loopback, reserved preset, distinct stores and health predicates. `operator_runner.go:22,27` retains the lock nofollow and mode/owner guards. `host/store/operator_runner.go:29` guards heartbeat read with `operatorAccess`. See §4 for the HEAD rechecks using named runtime tests. |
| C11 | `rg -n 'ApprovedBy !=\|Digest' host/store/operator_cas.go` and apply the S4c replay bank in the archive | Around `:54–60`, digest is checked but a confirmed row is returned on differing bytes. Archive replay patch still measures 13 lines. Appendix C assigns it before S7. |
| C12 | `cat "$HOME/.ailang/state/world-iter430/commit_ms.sh" "$HOME/.ailang/state/world-iter430/ctl_roundtrip.sh"` and `rg -n 'ubuntu-24.04\|timeout 8m' .github/workflows/ci.yml scripts/verify_go.sh` | Commit stages copy cumulative manifests, verify gofmt/build/vet/test and measure cmd/host/docs/app. Controller extracts exact four-backtick `diff M` fences. Linux is ubuntu-24.04. Race timeout 8m. Controller measured `host/daemon` race 342 s typical and three 480 s timeouts in ten runs. This patch adds no tests to it or `host/verifygate`. |

| C13 | `python3 "$S/mutants/recheck_head_focused.py"` and `python3 "$S/mutants/recheck_head.py"` run on separate pristine HEAD copies | The nine carried CLI variants each return rc 0 with `go test -count=1 -timeout 8m ./cmd/ailang-worldd -run '^Test(Operator\|RunnerPath)'`. The historical wider Quickstart subset is UNINFORMATIVE UNDER SANDBOX because of a loopback denial at `TestQuickstartSection9FlagsMatchTheCLI`. J-ST-HB-ACCESS returns rc 0 for the actual whole `go test -count=1 -timeout 8m ./host/store` on the HEAD variant. §4 rechecks all ten with new named tests. |
| C14 | `"$HOME/.pinned-ailang/ailang" version`, `go version` and `rg -n 'func .*ValidOperatorActionID\|func .*Available\|func .*WithContext\|func .*serviceFileHash' cmd/ailang-worldd host/store host/verbs` | AILANG v0.41.0, `24ee108`, Go 1.26.6 on darwin/arm64. `host/verbs/operator.go:54` validates exact action IDs, `host/store/operator_runner.go:21` provides generation/build/protocol/clock/readiness law, `service.go:103` reads hashes, `service_test.go:33,48,53` injects bootstrap/bootout/print. |

### 1.3 Planner decisions

- Take the episode lock for the first terminal operation even before the first companion startup. This closes the no-marker race for current binaries and preserves the canonical lock across stop/uninstall. Uncontended legacy or absent stores still proceed. Standalone launch checks for an existing descriptor and durable owners before report I/O. With no sidecar it acquires at the boundary before the PID or child. The inert runner has no action consumer. Appendix C gates any future consumer on holding the same lock before its CAS and on mixed-build drain.
- Hold the same open file description across nested preset start → launch. Validate the served credentials store at standalone launch before acquiring. The path is anchored by `store.CanonicalPath` and the validated episode ID. Separate credentials stores or episodes can proceed independently.
- Use a four-second read context for ownership resolution after the direct `LOCK_EX|LOCK_NB` attempt. Both the version and owner query use it. SQLite busy retry is disabled for the bounded readers because cancellation does not interrupt busy-retry sleep in the pinned driver. On unknown/corrupt evidence, decline with service-runner status as the remedy. An expired lease cannot clear a durable owner.
- `service runner start` captures the previous lease before bootstrap and requires a higher generation, the expected build/protocol, a current heartbeat and start readiness before the startup deadline. Initial publication with no predecessor qualifies. An already-loaded job checks a live heartbeat and retains its process on failure. A job bootstrapped by this start is stopped on timeout where the fake or real supervisor permits. Stop failure reports that distinction and preserves the lease evidence.

### 1.4 Premises found false or incomplete at HEAD

| ID | Evidence and consequence |
|---|---|
| F1 | C2/C6. The parked S4b patch included the upgrade pair-fence test. M4 already shipped all four arms and M5 added portable type ordering. Do not transplant the stale test hunk. |
| F2 | C10. The r1 J-ADDR-GIVEN result measured the removal of an existing guard, not base behavior. HEAD already rejects it. M3 adds a negative invocation that would reach the real internal entry on removal, with HOME isolated. |
| F3 | Compare the parked helper in `$S/parked-S4b.patch` with `operator_runner.go:51` and C9. A missing first-run marker bypassed exclusion in the parked draft. A separately opened descriptor collides even if a singleton marker appears later, as exercised in M1. |
| F4 | `sed -n '165,186p' host/store/writer_lock.go` documents that a context does not interrupt the driver's busy retry. The parked `time.Sleep` poll and version/schema open could extend beyond the remaining duration. M2 pins both deadline edges and M1 pins SQLite busy timeout zero. |
| F5 | C3. Old standalone launch bound the World store and could select credentials other than the served credentials DB. The M1 fixture serves both paths. A mismatch refuses before report lookup. Existing launch refusal snapshots exposed a new sidecar from the early acquisition draft. A no-create probe and a late acquire preserve ordinary refusal snapshots and close the preflight claim window. |
| F6 | C9. No paired-process episode-lock consumer or child-drain implementation exists at HEAD. No production action claim, reserve or observe became reachable in S4. Appendix C keeps dispatch gated. |
| F7 | C11. Different confirm replay bytes still return the prior confirmed row. Applying the 13-line bank passes `TestOperatorConfirmReplay`. Removing only its production hunk reproduces `different replay artifact returned <nil>` for both real signer variants. |
| F8 | `git ls-files host/store/operator_steps.go` prints nothing and the separate scratch reapply measures 132 lines. Receipt fencing and current-owner recovery are still absent. |

### 1.5 Size and CI race cost

Measure after each controller commit with `git diff -M --numstat c42fe698 HEAD -- cmd host docs app`. Sum additions plus deletions. Archive compare used `git diff --no-index -M --numstat "$S/base" "$S/st-$M"` and checked every changed path against `cmd host docs app`. Product totals are identical without the path filter. The unfiltered sum with both plan records included is **3137**. The records are separate controller inputs.

| Milestone | Added | Deleted | Cumulative | Increment | Snap STATUS / DELETED / manifest |
|---|---:|---:|---:|---:|---|
| M1 | 629 | 9 | 638 | 638 | 11/0/11 |
| M2 | 856 | 22 | 878 | 240 | 16/0/16 |
| M3 | 955 | 23 | 978 | 100 | 18/0/18 |

Scopes 1–3 fit near 900 lines at M3. Keep S4c-replay and S3b as separate milestones. Base `host/store` race was 37 s on the first pass and 36 s on the final-format baseline. Final was 45 s. Whole CLI race was sandbox-uninformative, and the focused CLI race passed in 11.016 s (12.67 s wall). The focused selection and output are in `$S/terminal-race-deliverable.log`. See §6 for the final result. Whole CLI race is sandbox-uninformative on bind. Keep host/daemon and verifygate race measurements separate from this sprint's fixture cost.

## 2. How the executor applies a sub-milestone

From the target worktree set `W=$PWD`, `T=$HOME/.ailang/state/world-iter430/exec-patches`, and `PLAN=$W/design_docs/planned/w-world-address-bar-verbs-v4s4b-sprint-plan.md`. The controller records the two plan inputs before snapshot commits. For each of M1, M2 and M3, set `M` to the ID and use this extraction.

````bash
mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
(cd "$W" && /usr/bin/patch -p1 -E < "$T/$M.patch")
for label in accept gates snapshot; do
  awk -v label="$label" '$0 == "````bash " label ".sh" {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$label.sh"
  test -s "$T/$label.sh"
  bash -n "$T/$label.sh"
done
````

Stop at any reject, unexpected edit, test/format failure or size drift. §7 gives the acceptance and cumulative snapshot commands. The executor never reformats or edits the patched files by hand.

## 3. Milestones

### M1 — terminal episode exclusion (638 cumulative)

Use the existing validated start, launch and F1 finish entry points. `TestOperatorTerminalExclusion` exercises plain, preset, resume, launch and finish with a genuinely migrated and signed claimed row. Resume releases the holder before the CLI invocation and uses a clock after expiry. `TestStartLaunchHeadlessMintsAndLaunchesOnce` proves the real nested transfer through mint, launch and receipt. The independent descriptor in `TestOperatorTerminalExclusion/nonwaiting` tests the first-managed transition. `TestOperatorEpisodeLockExclusive` distinguishes an exclusive lock from a shared lock.

### M2 — bounded companion startup and pair verification (878 cumulative)

Validate the previous owner generation and wait within the service startup budget. Use the fake bootstrap clock and exact expiry tests in `TestOperatorServiceStartup/gate`, with the existing down/schema/bootstrap failure coverage under `TestOperatorServiceStartup`. `TestOperatorServiceUpgradePairFence` retains the four landed arms and adds a supervisor query failure. Keep the attended service controller seams and pre-existing rollback example.

### M3 — carried service and store guards (978 cumulative)

Add portable lock/config symlink and unsafe mode fixtures, other-write-only setup, reserved preset, same stores, non-ok health, named loopback, explicit global address rejection and missing deadline. The fixture expectations include an absence of downstream writes and worker activity.

````bash accept.sh
#!/bin/bash
set -euo pipefail
W=$1
M=$2
test "$3" = fast
S=$HOME/.ailang/state/world-iter430/plan-scratch
cd "$W"
export AILANG_BIN="${AILANG_BIN:-$HOME/.pinned-ailang/ailang}"
export AILANG_STATE_DIR="$S/accept-ailang-state" AILANG_CACHE_DIR="$S/accept-ailang-cache"
export WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR WORLD_UPDATE_OPERATOR_ACTION_GOLDEN WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN
mkdir -p "$AILANG_STATE_DIR"
test -x "$AILANG_BIN"
cmd='TestOperatorTerminalExclusion|TestOperatorEpisodeLock(Refusals|Exclusive)|TestOperatorLaunch(ServedCredentials|ClaimDuringPreflight)|TestStartLaunchHeadlessMintsAndLaunchesOnce|TestLaunchHeadlessClaudeAC9AndSecondLaunch|TestLaunchRefusalsStartNothing|TestLaunchInteractiveRefusals|TestEpisodeFinishFenceAC1|TestCLICredentialWriterSites'
store='TestOperatorEpisodeOwnerReadSchema|TestProductionContextRoots|TestProductionGoSurface|TestOperatorSchema10'
case "$M" in
  M1) ;;
  M2|M3)
    cmd="$cmd|TestOperatorService(Startup|UpgradePairFence|WriteScope|CancelConfirmation)|TestHomeServiceUpgradeRollback"
    store="$store|TestOperatorRunnerUnavailable|TestOperatorCancelConfirmed"
    ;;
  *) exit 2 ;;
esac
if [ "$M" = M3 ]; then
  cmd="$cmd|TestOperatorRunner(LockFileRefusals|SetupRefusals)|TestOperatorService(Ownership|Pair)|TestRunnerPathRefusalClasses"
  store="$store|TestCredentialsMethodsRequireADeadline"
fi
for pkg in cmd/ailang-worldd host/store; do
  if [ "$pkg" = cmd/ailang-worldd ]; then pattern=$cmd; else pattern=$store; fi
  log=$(mktemp "$S/accept-$M.XXXXXX")
  go test -v -count=1 -timeout 8m "./$pkg" -run "^($pattern)$" | tee "$log"
  grep -q '^--- PASS:' "$log"
  if grep -q -- '--- SKIP:' "$log"; then echo "Unexpected skip in $log"; exit 1; fi
done
echo "ACCEPT_OK $M"
````

### 3.1 Branch traces

| Input or failure | Observable behavior and readbacks |
|---|---|
| Plain and preset start held by a valid claimed action | Existing usage/conflict status names the exact action plus `service runner status`. No tty, no worktree registration change, no session files, same credentials hash. |
| Resume held durably after descriptor close and expiry | Same exact ID. Old token, Git worktree list and store unchanged, no archive or tty. Different episode of the same store acquires. |
| Standalone launch held, or daemon credentials mismatch | Refuses before report client or child. Exact owner for contention, F6 pair diagnostic for mismatch. Store bytes, token and child evidence remain unchanged. The launch PID file remains absent. A claim during report preflight still blocks the late acquisition in `TestOperatorLaunchClaimDuringPreflight`. |
| F1 finish held | Terminal opened before health, one health request, zero report/dispatch/check/outcome traffic, same worktree status, HEAD and store. Existing no-terminal refusal is unchanged. |
| Pre-management terminal holder, second episode/store, transfer | Same-episode descriptor refuses without waiting during the singleton marker transition. Different keys acquire. Aliases share the canonical descriptor, wrong episode transfer refuses, and descriptor close permits reacquisition. A lock with no row reports contention. |
| Invalid mode/path, owner read failure, malformed ID and unclaimed row | Portable refusal for unusable lock. Unknown or corrupt ownership produces bounded refusal with no raw ID. A confirmed unowned row permits terminal use. Hash store and lock before/after and re-open the safe descriptor to verify cleanup. Legacy, absent, future and corrupt store cases are exercised by `TestOperatorEpisodeOwnerReadSchema`. |
| New startup, live prior job, bootstrap/print/lease faults | Fake supervisor counts, loaded state, heartbeat hash and exact generation readbacks distinguish initial/fresh/loaded success, missing baseline, stale/equal/late/wrong readiness, unreadable lease, stop failure and already-stopped cleanup. No repair path rewrites or clears the last heartbeat. |
| Installed record/plist, held lifecycle, loaded or unreadable supervisor | Specific upgrade diagnostic. Binary hash, fake kicks and write recorder show no staged/swap action. |

## 4. Mutation table

The named tests are runtime checks run against one mutation at a time on M3. Named and whole-package verdicts are independent. A successful mutated named test blocks publication of that claim. Full CLI variants hit the baseline loopback denial unless a complete run can be measured. Exact patch rewrites, commands, stdout and return codes are in `$S/mutants`. No inventory-only kill is counted as behavior.

**72/72 named kills.** Whole-package results were 14 killed, 3 survived and 55 **UNMEASURED IN SANDBOX**. The store-only survivors `E-CURRENT`, `E-ROW-EP`, `E-ROW-OWNER` are killed by the indicated real CLI fixtures. The store package has no such CLI entry point. No static-only kill was used.
For every CLI row marked UNMEASURED, apply only that variant, then run `AILANG_BIN=$HOME/.pinned-ailang/ailang go test -count=1 -timeout 8m ./cmd/ailang-worldd` outside the sandbox. Use the fixture env in §6 and restore the original bytes between variants. The exact `old`, `new`, `named_command` and `whole_command` are in `$S/mutants/ledger.json`. The first T-IS-BUSY attempt was a non-compiling edit, and the first T-LAUNCH-LATE-CLOSE test run let GC close a leaked descriptor. The ledger records the compiling edit and the rerun with deterministic descriptor lifetime.
| ID | Anchor on M3 | Named test and single assertion | Named | Whole package |
|---|---|---|---|---|
| T-PLAIN | `cmd/episode.go:364` | `TestOperatorTerminalExclusion/plain` — exact active ID and no entry-point effects | KILLED | UNMEASURED IN SANDBOX |
| T-PRESET | `cmd/episode_preset.go:282` | `TestOperatorTerminalExclusion/preset` — exact active ID and no entry-point effects | KILLED | UNMEASURED IN SANDBOX |
| T-RESUME | `cmd/episode_resume.go:175` | `TestOperatorTerminalExclusion/resume` — exact active ID and no entry-point effects | KILLED | UNMEASURED IN SANDBOX |
| T-LAUNCH | `cmd/episode_launch.go:289` | `TestOperatorTerminalExclusion/launch` — exact active ID and no entry-point effects | KILLED | UNMEASURED IN SANDBOX |
| T-FINISH | `cmd/episode_finish.go:128` | `TestOperatorTerminalExclusion/finish` — exact active ID and no entry-point effects | KILLED | UNMEASURED IN SANDBOX |
| T-LAUNCH-PAIR | `cmd/episode_launch.go:284` | `TestOperatorLaunchServedCredentials` — served-pair refusal before report | KILLED | UNMEASURED IN SANDBOX |
| T-LAUNCH-PREFLIGHT | `cmd/episode_launch.go:289` | `TestLaunchRefusalsStartNothing/nonempty` — legacy preflight leaves the whole fixture byte-identical | KILLED | UNMEASURED IN SANDBOX |
| T-LAUNCH-LATE | `cmd/episode_launch.go:454` | `TestOperatorLaunchClaimDuringPreflight` — claim between report and PID leaves no child or PID | KILLED | UNMEASURED IN SANDBOX |
| T-LAUNCH-ALREADY-HELD | `cmd/episode_launch.go:454` | `TestLaunchHeadlessClaudeAC9AndSecondLaunch` — existing standalone lock is reused at the effect boundary | KILLED | UNMEASURED IN SANDBOX |
| T-LAUNCH-NESTED | `cmd/episode_launch.go:454` | `TestStartLaunchHeadlessMintsAndLaunchesOnce` — nested launch preserves original owner without reacquisition | KILLED | UNMEASURED IN SANDBOX |
| T-LAUNCH-LATE-ERR | `cmd/episode_launch.go:456` | `TestOperatorLaunchClaimDuringPreflight` — late contention exits before PID and child | KILLED | UNMEASURED IN SANDBOX |
| T-LAUNCH-LATE-CLOSE | `cmd/episode_launch.go:460` | `TestLaunchHeadlessClaudeAC9AndSecondLaunch` — first launch releases late-opened descriptor | KILLED | UNMEASURED IN SANDBOX |
| T-ABSENT-MODE | `cmd/operator_episode_lock.go:38` | `TestLaunchHeadlessClaudeAC9AndSecondLaunch` — unmanaged preflight can proceed without a sidecar | KILLED | UNMEASURED IN SANDBOX |
| T-ABSENT-OK | `cmd/operator_episode_lock.go:45` | `TestLaunchHeadlessClaudeAC9AndSecondLaunch` — unmanaged successful launch acquires at boundary | KILLED | UNMEASURED IN SANDBOX |
| T-LAUNCH-TRANSFER | `cmd/episode_preset.go:316` | `TestStartLaunchHeadlessMintsAndLaunchesOnce` — mint and child/receipt once with the original descriptor | KILLED | UNMEASURED IN SANDBOX |
| T-CANONICAL | `cmd/operator_episode_lock.go:23` | `TestOperatorTerminalExclusion/nonwaiting` — symlink spelling transfers the same descriptor | KILLED | UNMEASURED IN SANDBOX |
| T-STORE-KEY | `cmd/operator_episode_lock.go:15` | `TestOperatorTerminalExclusion/nonwaiting` — different store in the same directory still acquires | KILLED | UNMEASURED IN SANDBOX |
| T-EPISODE-KEY | `cmd/operator_episode_lock.go:15` | `TestOperatorTerminalExclusion/nonwaiting` — different episode still acquires | KILLED | UNMEASURED IN SANDBOX |
| T-HELD-KEY | `cmd/operator_episode_lock.go:29` | `TestOperatorTerminalExclusion/nonwaiting` — wrong episode transfer is refused | KILLED | UNMEASURED IN SANDBOX |
| T-HELD-REUSE | `cmd/operator_episode_lock.go:28` | `TestStartLaunchHeadlessMintsAndLaunchesOnce` — nested launch uses the held descriptor | KILLED | UNMEASURED IN SANDBOX |
| T-LAUNCH-RELEASE | `cmd/episode_launch.go:293` | `TestLaunchHeadlessClaudeAC9AndSecondLaunch` — successful standalone launch releases lock for the next launch | KILLED | UNMEASURED IN SANDBOX |
| T-NB | `cmd/operator_runner.go:35` | `TestOperatorTerminalExclusion/nonwaiting` — same-episode contender returns EWOULDBLOCK before fuse | KILLED | UNMEASURED IN SANDBOX |
| T-EX | `cmd/operator_runner.go:35` | `TestOperatorEpisodeLockExclusive` — shared holder prevents exclusive acquisition | KILLED | UNMEASURED IN SANDBOX |
| T-IS-BUSY | `cmd/operator_episode_lock.go:39` | `TestOperatorTerminalExclusion/plain` — busy diagnostic still resolves exact owner ID | KILLED | UNMEASURED IN SANDBOX |
| T-OPEN-FAULT | `cmd/operator_episode_lock.go:39` | `TestOperatorEpisodeLockRefusals/unsafe_lock` — unsafe-lock diagnostic and unchanged bytes | KILLED | UNMEASURED IN SANDBOX |
| T-LOCK-OK | `cmd/operator_episode_lock.go:45` | `TestOperatorTerminalExclusion/nonwaiting` — unknown busy holder refuses when no row exists | KILLED | UNMEASURED IN SANDBOX |
| T-READ-OK | `cmd/operator_episode_lock.go:45` | `TestOperatorEpisodeLockRefusals/ownership_read` — cannot read ownership and no lock handed off | KILLED | UNMEASURED IN SANDBOX |
| T-OWNER-BLOCK | `cmd/operator_episode_lock.go:45` | `TestOperatorTerminalExclusion/resume` — durable owner still refuses after releasing the kernel holder | KILLED | UNMEASURED IN SANDBOX |
| T-CLOSE | `cmd/operator_episode_lock.go:48` | `TestOperatorEpisodeLockRefusals/invalid_action_id` — refusal permits independent descriptor reopen | KILLED | UNMEASURED IN SANDBOX |
| T-READ-DIAG | `cmd/operator_episode_lock.go:53` | `TestOperatorEpisodeLockRefusals/ownership_read` — read failure uses ownership diagnostic | KILLED | UNMEASURED IN SANDBOX |
| T-ID-VALID | `cmd/operator_episode_lock.go:53` | `TestOperatorEpisodeLockRefusals/invalid_action_id` — corrupt id remains redacted | KILLED | UNMEASURED IN SANDBOX |
| T-ID-EMPTY | `cmd/operator_episode_lock.go:53` | `TestOperatorTerminalExclusion/nonwaiting` — empty action on locked pre-claim episode is ordinary contention | KILLED | UNMEASURED IN SANDBOX |
| T-ID-DIAG | `cmd/operator_episode_lock.go:55` | `TestOperatorTerminalExclusion/plain` — exact valid ID instead of generic contention | KILLED | UNMEASURED IN SANDBOX |
| E-LEGACY | `store/operator_runner.go:100` | `TestOperatorEpisodeOwnerReadSchema` — v9 does not query missing operator tables | KILLED | KILLED |
| E-CURRENT | `store/operator_runner.go:100` | `TestOperatorTerminalExclusion/resume` — valid v10 durable owner cannot be bypassed | KILLED | SURVIVED |
| E-FUTURE | `store/operator_runner.go:146` | `TestOperatorEpisodeOwnerReadSchema` — future version refuses without changing the file | KILLED | KILLED |
| E-FILE-TYPE | `store/operator_runner.go:134` | `TestOperatorEpisodeOwnerReadSchema` — ordinary directory cannot be read as SQLite | KILLED | KILLED |
| E-NO-ROW | `store/operator_runner.go:105` | `TestOperatorEpisodeOwnerReadSchema` — v10 no-owner case is normal | KILLED | KILLED |
| E-MISSING | `store/operator_runner.go:93` | `TestOperatorEpisodeOwnerReadSchema` — missing database does not block first legacy terminal | KILLED | KILLED |
| E-ROW-EP | `store/operator_runner.go:104` | `TestOperatorTerminalExclusion/resume` — other episode can acquire despite real durable action | KILLED | SURVIVED |
| E-ROW-OWNER | `store/operator_runner.go:104` | `TestOperatorEpisodeLockRefusals/unclaimed_action` — confirmed but unclaimed row permits terminal | KILLED | SURVIVED |
| E-BUSY | `store/operator_runner.go:137` | `TestOperatorEpisodeOwnerReadSchema` — query reports no SQLite busy retry | KILLED | KILLED |
| E-MODE | `store/operator_runner.go:137` | `TestOperatorEpisodeOwnerReadSchema` — SQLite write attempt denied | KILLED | KILLED |
| E-DEADLINE | `store/operator_runner.go:88` | `TestCredentialsMethodsRequireADeadline` — ownership read requires caller deadline | KILLED | KILLED |
| E-BAD-VERSION | `store/operator_runner.go:142` | `TestOperatorEpisodeOwnerReadSchema` — corrupt schema refuses | KILLED | KILLED |
| H-DEADLINE | `store/operator_runner.go:112` | `TestCredentialsMethodsRequireADeadline` — distinguish guard before I/O using invalid path | KILLED | KILLED |
| S-PRE-TIME | `cmd/operator_runner.go:112` | `TestOperatorServiceStartup/gate/same_generation` — no extra read beyond exact deadline | KILLED | UNMEASURED IN SANDBOX |
| S-POST-TIME | `cmd/operator_runner.go:117` | `TestOperatorServiceStartup/gate/late` — deadline equality is a timeout even with a fresh result | KILLED | UNMEASURED IN SANDBOX |
| S-GENERATION | `cmd/operator_runner.go:120` | `TestOperatorServiceStartup/gate/same_generation` — equal prior lease cannot report bootstrap ready | KILLED | UNMEASURED IN SANDBOX |
| S-READ | `cmd/operator_runner.go:120` | `TestOperatorServiceStartup/gate/bad_read` — partially scanned heartbeat cannot satisfy gate | KILLED | UNMEASURED IN SANDBOX |
| S-AVAILABLE | `cmd/operator_runner.go:121` | `TestOperatorServiceStartup/gate/finish_only` — fresh generation without start readiness fails | KILLED | UNMEASURED IN SANDBOX |
| S-SLEEP-CAP | `cmd/operator_runner.go:124` | `TestOperatorServiceStartup/gate/same_generation` — last delay clips at the deadline | KILLED | UNMEASURED IN SANDBOX |
| S-BASE-ERR | `cmd/operator_service.go:436` | `TestOperatorServiceStartup/gate/bad_previous` — bad prior read prevents boot | KILLED | UNMEASURED IN SANDBOX |
| S-NO-BASE | `cmd/operator_service.go:436` | `TestOperatorServiceStartup` — no previous heartbeat is valid on initial startup | KILLED | UNMEASURED IN SANDBOX |
| S-TIMEOUT | `cmd/operator_service.go:444` | `TestOperatorServiceStartup/gate/same_generation` — bootstrap success requires a new heartbeat | KILLED | UNMEASURED IN SANDBOX |
| S-CLEANUP-ERR | `cmd/operator_service.go:445` | `TestOperatorServiceStartup/gate/failed_stop` — stop failure has a distinct diagnostic | KILLED | UNMEASURED IN SANDBOX |
| S-ALREADY-STOPPED | `cmd/operator_service.go:445` | `TestOperatorServiceStartup/gate/already_stopped` — already gone is a successful cleanup | KILLED | UNMEASURED IN SANDBOX |
| S-LOADED | `cmd/operator_service.go:452` | `TestOperatorServiceStartup/gate/already_loaded_stale` — stale loaded job refuses and is left up | KILLED | UNMEASURED IN SANDBOX |
| S-NOW-FUTURE | `store/operator_runner.go:27` | `TestOperatorServiceStartup/gate/future` — future fresh generation refuses | KILLED | KILLED |
| S-NOW-EQUAL | `store/operator_runner.go:27` | `TestOperatorServiceStartup/gate/fresh_restart` — same-second new owner is ready | KILLED | KILLED |
| S-NOW-STALE | `store/operator_runner.go:27` | `TestOperatorServiceStartup/gate/stale` — exactly ten seconds is stale | KILLED | KILLED |
| P-PRINT-ERROR | `cmd/service_upgrade.go:211` | `TestOperatorServiceUpgradePairFence/inspect_error` — generic supervisor query error fences swap | KILLED | UNMEASURED IN SANDBOX |
| J-WRITABLE-OTHERS | `cmd/operator_service.go:78` | `TestOperatorServiceOwnership/other_writable` — chmod 0702 has other-write only | KILLED | UNMEASURED IN SANDBOX |
| J-LOCKFILE-MODE | `cmd/operator_runner.go:31` | `TestOperatorRunnerLockFileRefusals/mode` — 0644 file refused without mutation | KILLED | UNMEASURED IN SANDBOX |
| J-LOCKFILE-NOFOLLOW | `cmd/operator_runner.go:26` | `TestOperatorRunnerLockFileRefusals/symlink` — valid link target remains unopened | KILLED | UNMEASURED IN SANDBOX |
| J-CONFIG-NOFOLLOW | `cmd/operator_service.go:92` | `TestOperatorRunnerSetupRefusals/linked_config` — install from symlink refuses without writing | KILLED | UNMEASURED IN SANDBOX |
| J-PRESET-SE | `cmd/operator_service.go:141` | `TestOperatorRunnerSetupRefusals/reserved_preset` — reserved se-tools name refused before install | KILLED | UNMEASURED IN SANDBOX |
| J-DISTINCT-STORES | `cmd/operator_service.go:167` | `TestOperatorServicePair/distinct` — same canonical file refuses with specific class | KILLED | UNMEASURED IN SANDBOX |
| J-HEALTH-STATUS | `cmd/operator_service.go:189` | `TestOperatorServicePair/health_status` — HTTP 200 JSON status starting refuses | KILLED | UNMEASURED IN SANDBOX |
| J-LOOPBACK-IP | `cmd/operator_service.go:130` | `TestOperatorServicePair/named_loopback` — localhost passes broad service predicate but not the runner | KILLED | UNMEASURED IN SANDBOX |
| J-ADDR-GIVEN | `cmd/main.go:383` | `TestOperatorRunnerSetupRefusals/global_address` — explicit global address refuses before runner startup | KILLED | UNMEASURED IN SANDBOX |
| J-ST-HB-ACCESS | `store/operator_runner.go:31` | `TestCredentialsMethodsRequireADeadline` — raw heartbeat read refuses missing deadline | KILLED | KILLED |

`cmd/` in the table is `cmd/ailang-worldd/`, and `store/` is `host/store/`. The `sql.Open` error is a defensive propagation arm for a statically registered driver. Stable-path absent-file behavior is separately measured. An OS I/O failure between successful canonicalization and stat cannot be made with an ordinary stable fixture and gets no claimed kill. Invoking `Close` on a nil `*os.File` would have no visible effect, so removing only the nil skip is equivalent. J-OWN-UID has no non-root fixture and is carried in Appendix C.

## 5. Files

- `cmd/ailang-worldd/episode.go`
- `cmd/ailang-worldd/episode_finish.go`
- `cmd/ailang-worldd/episode_launch.go`
- `cmd/ailang-worldd/episode_launch_test.go`
- `cmd/ailang-worldd/episode_preset.go`
- `cmd/ailang-worldd/episode_resume.go`
- `cmd/ailang-worldd/operator_episode_lock.go`
- `cmd/ailang-worldd/operator_episode_lock_test.go`
- `cmd/ailang-worldd/operator_refusals_test.go`
- `cmd/ailang-worldd/operator_runner.go`
- `cmd/ailang-worldd/operator_service.go`
- `cmd/ailang-worldd/operator_service_test.go`
- `cmd/ailang-worldd/operator_startup_test.go`
- `cmd/ailang-worldd/service.go`
- `cmd/ailang-worldd/service_test.go`
- `host/store/credentials_test.go`
- `host/store/operator_runner.go`
- `host/store/operator_schema_test.go`

Only `cmd/ailang-worldd` and `host/store` Go packages change. The CLI's writer-site, context, dependency, commands/help and QUICKSTART inventories keep their expected entries. No tests are deleted. The paired-upgrade test is extended in place.

## 6. Sprint gates

Extract this script using §2. Base uses a pristine archive. For scratch export `GATES_GIT_REF` for the checkout and `GATES_CHANGED_FILE` for the cumulative path list. Use `GATES_SANDBOX=1` only in a workspace-write sandbox. The controller runs `bash "$T/gates.sh" "$W" "$HOME/.ailang/state/world-iter430/plan-scratch/controller-gates" final` outside the sandbox with `GATES_SKIP_VERIFY_GO` unset. The corresponding affected-package plain/race legs run with `-timeout 8m`. Linux CI owns the full race suite with its exact target toolchain and launchctl observations are Darwin read-only controller evidence.

````bash gates.sh
#!/bin/bash
# Usage: gates.sh <tree> <logdir> <base|final|scratch> [sha]
set -u
set -o pipefail
T=$1; L=$2; MODE=$3
B=c42fe69888ce0c7f37a2ef967738085b1a27c1e6
REF=${GATES_GIT_REF:-$T}
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
  if [ "$rc" -ne 0 ] && [ "${GATES_SANDBOX:-0}" = 1 ] &&
     grep -Eq 'listen .*bind: operation not permitted|OBSERVATION REQUIRES CONTROLLER' "$L/$name.txt"; then
    line "LEG $name rc=$rc $(( $(date +%s)-start ))s UNINFORMATIVE UNDER SANDBOX"
    UNKNOWN="$UNKNOWN $name"
  else
    line "LEG $name rc=$rc $(( $(date +%s)-start ))s"
    [ "$rc" -eq 0 ] || FAILS="$FAILS $name"
  fi
}
observe() {
  if [ "${GATES_SANDBOX:-0}" = 1 ]; then
    echo 'OBSERVATION REQUIRES CONTROLLER (operator paths and read-only launchctl list)'
    return 77
  fi
  local d file jobs n
  for d in "$HOME/Library/LaunchAgents" "$HOME/.ailang-world/bin" "$HOME/.ailang-world"; do
    if [ ! -d "$d" ]; then printf 'absent %s\n' "$d"; continue; fi
    echo "DIRECTORY $d"
    if [ "$d" = "$HOME/.ailang-world" ]; then
      find "$d" -maxdepth 1 \( -type f -o -type l \) -name '*.json' -print0
    else
      find "$d" \( -type f -o -type l \) -print0
    fi | LC_ALL=C sort -z | while IFS= read -r -d '' file; do
      /bin/ls -la -d "$file"
      if [ -L "$file" ]; then readlink "$file" | shasum -a 256; else shasum -a 256 "$file"; fi
    done
    if [ "$d" = "$HOME/.ailang-world" ]; then
      for file in "$d/service.json" "$d/operator-runner.json"; do
        if [ ! -e "$file" ] && [ ! -L "$file" ]; then printf 'absent %s\n' "$file"; fi
      done
    fi
  done
  jobs=$(/bin/launchctl list)
  printf '%s\n' "$jobs" | grep -Eq '^PID[[:space:]]+Status[[:space:]]+Label$'
  n=$(printf '%s\n' "$jobs" | grep -c 'ailang-world' || test "$?" -eq 1)
  printf 'ailang-world launchctl jobs: %s\n' "$n"
}
leg observation-before 'observe'
paths() {
  if [ "$MODE" = scratch ]; then cat "$GATES_CHANGED_FILE"
  else git -C "$REF" diff --name-only "$B" HEAD
  fi
}
leg paths 'paths'
TOUCHED=$(awk -F/ '/\.go$/ && ($1=="cmd" || $1=="host") {NF--; OFS="/"; $1=$1; print}' "$L/paths.txt" | sort -u | tr '\n' ' ')
EXPECTED='cmd/ailang-worldd host/store '
line "TOUCHED_PACKAGES: $TOUCHED"
leg touched 'if [ "$MODE" = base ]; then test -z "$TOUCHED"; else test "$TOUCHED" = "$EXPECTED"; fi'
leg status 'git -C "$REF" status --porcelain'
leg gofmt 'out=$(gofmt -l cmd host); printf "%s\n" "$out"; test -z "$out"'
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./..."
leg verify_ail './scripts/verify_ail.sh'
PACKAGES=$TOUCHED
if [ "$MODE" = base ]; then PACKAGES=$EXPECTED; fi
for p in $PACKAGES; do
  tag=${p//\//-}
  leg "pkg-$tag" "go test -count=1 -timeout 8m ./$p"
  leg "race-$tag" "/usr/bin/time -p go test -count=1 -race -timeout 8m ./$p"
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
inv cmd/ailang-worldd TestEveryVerbHasHelp TestVerbTableMatchesTheCLISwitch TestCLICredentialWriterSites TestOperatorFlagParity TestQuickstartServiceSectionMatchesTheCLIHelp TestPluginCommandsMatchTheCLIHelp TestHelpIsNeverBlockedByTheRegistryKeyGuard TestHomeServiceSeamOnly TestHomeServiceInstallAndPlist TestHomeServiceUpgradeRollback
inv host/store TestProductionContextRoots TestProductionGoSurface TestStoreSurfaceTakesContext TestCredentialsMethodsRequireADeadline TestCredentialsSchemaMatchesCanonicalDDL TestOperatorSchema10 TestOperatorClaimBeforeExpiry TestOperatorClaimRace TestOperatorTerminalExclusion TestWriterLockCrossProcessContention TestWriterLockHeldIsNonMutating
route_fence() {
  local routes posts old_routes old_posts
  routes=$(grep -Ec 'mux\.(HandleFunc|Handle)\(' host/daemon/daemon.go)
  posts=$(grep -Fc 'mux.HandleFunc("POST /surface/' host/daemon/daemon.go)
  old_routes=$(git -C "$REF" show "$B:host/daemon/daemon.go" | grep -Ec 'mux\.(HandleFunc|Handle)\(')
  old_posts=$(git -C "$REF" show "$B:host/daemon/daemon.go" | grep -Fc 'mux.HandleFunc("POST /surface/')
  printf 'DAEMON_ROUTES: %s (base %s); SURFACE_POSTS: %s (base %s)\n' "$routes" "$old_routes" "$posts" "$old_posts"
  test "$routes" -gt 0
  test "$routes" = "$old_routes"
  test "$posts" = 3
  test "$posts" = "$old_posts"
  local f
  for f in host/daemon/daemon.go host/daemon/surface_page.go host/surfacepage/app/app.js; do
    git -C "$REF" show "$B:$f" | cmp - "$f"
    printf 'UNCHANGED: %s\n' "$f"
  done
  grep -n 'surfaceCSP' host/daemon/surface_page.go
}
leg routes 'route_fence'
scope_fence() {
  local allowed extra count
  allowed='^(cmd/ailang-worldd/(episode(_(finish|launch|launch_test|preset|resume))?|operator_(episode_lock(_test)?|refusals_test|runner|service(_test)?|startup_test)|service(_test)?)\.go|host/store/(credentials_test|operator_(runner|schema_test))\.go|design_docs/planned/(w-world-address-bar-verbs-v4s4b-sprint-plan\.md|sprint_w-world-address-bar-verbs-v4s4b\.json))$'
  printf '%s\n' cmd/ailang-worldd/episode.go host/store/operator_runner.go | grep -E "$allowed"
  printf '%s\n' world/transitions.ail | grep -Ev "$allowed"
  extra=$(grep -Ev "$allowed" "$L/paths.txt" || test "$?" -eq 1)
  printf 'UNTOUCHED_VIOLATIONS: %s\n' "${extra:-none}"
  test -z "$extra"
}
leg scope 'scope_fence'
caller_fence() {
  local old callers external
  old=$(grep -R -nE --include='*.go' 'authority\.RecheckVerb\(' cmd host | grep -v '_test\.go:')
  printf 'Existing live caller control: %s\n' "$old"
  test -n "$old"
  if [ "$MODE" != base ]; then
    callers=$(grep -R -nE --include='*.go' '\.(PutOperatorHeartbeat|OperatorHeartbeat|runnerTick)\(' cmd host | grep -v '_test\.go:')
    printf 'Heartbeat and catalogue callers:\n%s\n' "$callers"
    printf '%s\n' "$callers" | grep -F 'cmd/ailang-worldd/operator_runner.go:'
    external=$(printf '%s\n' "$callers" | grep -Ev '^(cmd/ailang-worldd|host/store)/operator_runner.go:' || test "$?" -eq 1)
    test -z "$external"
  fi
  if grep -R -nE --include='*.go' '\.(ClaimOperatorAction|ReserveOperatorStep|ObserveOperatorStep)\(' cmd host | grep -v '_test\.go:'; then
    return 1
  fi
  grep -n 'syscall.Flock.*LOCK_NB' cmd/ailang-worldd/service.go host/store/writer_lock_unix.go
  if [ "$MODE" != base ]; then
    grep -n 'syscall.Flock.*LOCK_NB' cmd/ailang-worldd/operator_runner.go cmd/ailang-worldd/operator_service.go
    grep -nE '(openTerminalEpisodeLock|terminalEpisodeLock)\(' cmd/ailang-worldd/episode{,_launch,_finish,_preset,_resume}.go
    go test -v -count=1 -timeout 8m ./cmd/ailang-worldd -run '^(TestOperatorEpisodeLockExclusive|TestOperatorTerminalExclusion)$'
  fi
}
leg callers 'caller_fence'
if [ "${GATES_SKIP_VERIFY_GO:-0}" = 1 ]; then
  echo 'GATES_SKIP_VERIFY_GO=1' > "$L/verify_go.txt"
  line "LEG verify_go SKIPPED 0s (GATES_SKIP_VERIFY_GO=1)"
else
  leg verify_go 'VERIFY_GO_LEGS=plain AILANG_BIN="$HOME/.pinned-ailang/ailang" ./scripts/verify_go.sh'
fi
leg observation-after 'observe'
if ! cmp -s "$L/observation-before.txt" "$L/observation-after.txt"; then
  diff -u "$L/observation-before.txt" "$L/observation-after.txt" | tee "$L/observation-diff.txt"
  FAILS="$FAILS observation_changed"
fi
line "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none} UNINFORMATIVE:${UNKNOWN:- none}"
tail -n 1 "$L/legs.txt" > "$L/ALL_DONE.$SHA"
[ -z "$FAILS" ] || exit 1
[ -z "$UNKNOWN" ] || exit 77
````

Scratch results follow. The first candidate sweep, while concurrent mutation builds were running, returned an unexpected deadline error in the existing `TestDurableStallAtCutoffIsUncertainAndReconciles/Commit` race case. That test passed when run alone, and the full store race package passed on the next sweep. See `$S/gates-final/race-host-store.txt`, `$S/race-host-store-isolated.log`, and the final gate transcript below.

The broader terminal regression command with the known loopback and browser cases excluded passed. Its preceding run of the loopback arm was UNINFORMATIVE UNDER SANDBOX. The exact run was `go test -count=1 -timeout 8m -v ./cmd/ailang-worldd -run '^(Test(Episode(Start|Resume|Finish)|Preset(Start|Local)|StartLaunch|Launch|Operator|FinishConfirmOnSurface))' -skip '^(TestEpisodeStartWithoutAControllingTerminalRefusesEndToEnd|TestFinishConfirmOnSurfaceInChrome)$'` with the §6 fixture environment. Results are in `$S/terminal-regression-v2.log`.

````text
LEG observation-before rc=77 0s UNINFORMATIVE UNDER SANDBOX
LEG paths rc=0 0s
TOUCHED_PACKAGES: cmd/ailang-worldd host/store 
LEG touched rc=0 0s
LEG status rc=0 0s
LEG gofmt rc=0 0s
LEG build rc=0 1s
LEG vet rc=0 0s
LEG fence rc=0 1s
LEG verify_ail rc=0 53s
LEG pkg-cmd-ailang-worldd rc=1 3s UNINFORMATIVE UNDER SANDBOX
LEG race-cmd-ailang-worldd rc=1 18s UNINFORMATIVE UNDER SANDBOX
LEG pkg-host-store rc=0 15s
LEG race-host-store rc=0 45s
LEG inv-host-approval rc=0 1s
LEG inv-host-approval-approvaltest rc=0 1s
LEG inv-host-authority rc=0 1s
LEG inv-host-boundary rc=0 1s
LEG inv-host-daemon rc=0 2s
LEG inv-host-verbs rc=0 0s
LEG inv-cmd-ailang-worldd rc=0 1s
LEG inv-host-store rc=0 1s
LEG routes rc=0 0s
LEG scope rc=0 0s
LEG callers rc=0 2s
LEG verify_go SKIPPED 0s (GATES_SKIP_VERIFY_GO=1)
LEG observation-after rc=77 0s UNINFORMATIVE UNDER SANDBOX
GATES_DONE 146s FAILS: none UNINFORMATIVE: observation-before pkg-cmd-ailang-worldd race-cmd-ailang-worldd observation-after
````

## 7. Executor directive notes

- For each stage, run `AILANG_BIN=$HOME/.pinned-ailang/ailang bash "$T/accept.sh" "$W" M1 fast`, then the same exact command with `M2`, then `M3` after their respective patch applications. The script exports the other required fixture env vars and sends the pinned interpreter cache to mission scratch. Record test output with pipefail under `$HOME/.ailang/state/world-iter430/exec-logs/`.
- Run `bash "$T/snapshot.sh" "$W" "$M"`. For the controller to commit, pass `commit_ms.sh "$M" "$msgfile" "AILANG_BIN=$HOME/.pinned-ailang/ailang bash '$T/accept.sh' '$W' '$M' fast"`. Check the cumulative numstat and manifest against §1.5. The worktree snapshot enumerates the product pathspec, so absent product directories require no special setup.
- Run the controller round trip with `bash "$HOME/.ailang/state/world-iter430/ctl_roundtrip.sh" M1 M2 M3`. Planner additionally compares complete exports, including empty/absent path cases. `SNAP_STATUS_LIST` and `SNAP_PARENT` are only for archive verification. Executor snapshots live at the default `snap/<M>` path. Each compiled stage passed gofmt, build, vet and its marked `ACCEPT_OK` suite. The latest isolated runs are in `$S/st-M1.cache-accept.log`, `$S/st-M2.cache-accept.log` and `$S/st-M3.cache-accept.log`.

````bash snapshot.sh
#!/bin/bash
set -euo pipefail
W=$1
M=$2
B=c42fe69888ce0c7f37a2ef967738085b1a27c1e6
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter430/snap}/$M
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

````text
cmd/ailang-worldd/episode.go: OK
cmd/ailang-worldd/episode_finish.go: OK
cmd/ailang-worldd/episode_launch.go: OK
cmd/ailang-worldd/episode_launch_test.go: OK
cmd/ailang-worldd/episode_preset.go: OK
cmd/ailang-worldd/episode_resume.go: OK
cmd/ailang-worldd/operator_episode_lock.go: OK
cmd/ailang-worldd/operator_episode_lock_test.go: OK
cmd/ailang-worldd/operator_runner.go: OK
host/store/operator_runner.go: OK
host/store/operator_schema_test.go: OK
      11 /Users/voightkampff/.ailang/state/world-iter430/plan-scratch/snapshot-deliverable/M1/STATUS
       0 /Users/voightkampff/.ailang/state/world-iter430/plan-scratch/snapshot-deliverable/M1/DELETED
      11 /Users/voightkampff/.ailang/state/world-iter430/plan-scratch/snapshot-deliverable/M1/MANIFEST.sha256
      22 total
cmd/ailang-worldd/episode.go: OK
cmd/ailang-worldd/episode_finish.go: OK
cmd/ailang-worldd/episode_launch.go: OK
cmd/ailang-worldd/episode_launch_test.go: OK
cmd/ailang-worldd/episode_preset.go: OK
cmd/ailang-worldd/episode_resume.go: OK
cmd/ailang-worldd/operator_episode_lock.go: OK
cmd/ailang-worldd/operator_episode_lock_test.go: OK
cmd/ailang-worldd/operator_runner.go: OK
cmd/ailang-worldd/operator_service.go: OK
cmd/ailang-worldd/operator_service_test.go: OK
cmd/ailang-worldd/operator_startup_test.go: OK
cmd/ailang-worldd/service.go: OK
cmd/ailang-worldd/service_test.go: OK
host/store/operator_runner.go: OK
host/store/operator_schema_test.go: OK
      16 /Users/voightkampff/.ailang/state/world-iter430/plan-scratch/snapshot-deliverable/M2/STATUS
       0 /Users/voightkampff/.ailang/state/world-iter430/plan-scratch/snapshot-deliverable/M2/DELETED
      16 /Users/voightkampff/.ailang/state/world-iter430/plan-scratch/snapshot-deliverable/M2/MANIFEST.sha256
      32 total
cmd/ailang-worldd/episode.go: OK
cmd/ailang-worldd/episode_finish.go: OK
cmd/ailang-worldd/episode_launch.go: OK
cmd/ailang-worldd/episode_launch_test.go: OK
cmd/ailang-worldd/episode_preset.go: OK
cmd/ailang-worldd/episode_resume.go: OK
cmd/ailang-worldd/operator_episode_lock.go: OK
cmd/ailang-worldd/operator_episode_lock_test.go: OK
cmd/ailang-worldd/operator_refusals_test.go: OK
cmd/ailang-worldd/operator_runner.go: OK
cmd/ailang-worldd/operator_service.go: OK
cmd/ailang-worldd/operator_service_test.go: OK
cmd/ailang-worldd/operator_startup_test.go: OK
cmd/ailang-worldd/service.go: OK
cmd/ailang-worldd/service_test.go: OK
host/store/credentials_test.go: OK
host/store/operator_runner.go: OK
host/store/operator_schema_test.go: OK
      18 /Users/voightkampff/.ailang/state/world-iter430/plan-scratch/snapshot-deliverable/M3/STATUS
       0 /Users/voightkampff/.ailang/state/world-iter430/plan-scratch/snapshot-deliverable/M3/DELETED
      18 /Users/voightkampff/.ailang/state/world-iter430/plan-scratch/snapshot-deliverable/M3/MANIFEST.sha256
      36 total
host/store/operator_runner.go: OK
host/store/operator_schema_test.go: OK
       2 /Users/voightkampff/.ailang/state/world-iter430/plan-scratch/snapshot-deliverable/host-only/STATUS
       0 /Users/voightkampff/.ailang/state/world-iter430/plan-scratch/snapshot-deliverable/host-only/DELETED
       2 /Users/voightkampff/.ailang/state/world-iter430/plan-scratch/snapshot-deliverable/host-only/MANIFEST.sha256
       4 total
M1 applied (     804 lines)
  diff -r cmd vs st-M1: rc=0
  diff -r host vs st-M1: rc=0
  diff -r docs vs st-M1: rc=0
M2 applied (     368 lines)
  diff -r cmd vs st-M2: rc=0
  diff -r host vs st-M2: rc=0
  diff -r docs vs st-M2: rc=0
M3 applied (     138 lines)
  diff -r cmd vs st-M3: rc=0
  diff -r host vs st-M3: rc=0
  diff -r docs vs st-M3: rc=0
DONE
M1 diff -r full-tree rc=0
M2 diff -r full-tree rc=0
M3 diff -r full-tree rc=0
````

## 8. Risks

- Flock acts on open descriptions for the tested independent opens and transfer. Keep the path inode stable. Deleting a lock file while any descriptor owns it would break exclusion. Darwin may expose different symlink mode bits. Tests classify wrong type before mode or check refusal and unchanged evidence. Ubuntu 24.04 reruns the behavioral pins on Linux.
- Exact durable ID depends on a readable valid row. The immediate lock attempt can coincide with a holder before its claim or with a corrupt/unavailable read. The operation refuses with the available evidence and explicit runner-status remedy. A stale heartbeat cannot distinguish dead children from live ones.
- A new runner needs an attended drain of terminals started by binaries that predate this lock and of any resumed child with uncertain outcome. Appendix C retains this production admission gate. The configured wait bounds readiness polling after bootstrap. Supervisor command budgets remain independent.
- The production open-description owner check authenticates the file UID. J-OWN-UID requires a legitimately alternate owner and cannot be measured in this non-root fixture. A mocked stat UID would not prove the kernel/path condition.

## 9. Acceptance criteria

1. M1–M3 patch from `c42fe698`, match the compiled archive stages byte-for-byte, and verify cumulative manifest and numstat under the controller filter.
2. Real entry points name the exact action, acquire immediately or refuse, preserve durable ownership through expiry and descriptor close, and preserve isolated progress and F1.
3. Startup waits for compatible fresh generation, fails closed at the deadline and preserves the observed lease through cleanup. All five pair-fence arms retain the binary and supervisor evidence.
4. Scope 3 predicates have measured named runtime outcomes. Controller runs the unmeasured whole-package legs. Keep the paired runner consumer, drain and receipts before dispatch.
5. Surface/route inventories, writer inventory, context roots and dependency allow-lists run by name. New code stays in host Go at the CLI/SQLite boundaries. No AILANG package can own these file descriptors, SQLite connections or local service lifecycle commands.

## Appendix C — after this sprint

- **Runner consumer and drain before S5/S6.** HEAD and final stage each have zero production claim/reserve/observe call sites (C9). M1 covers first-managed lock ordering using two independent descriptors and migrates actual stores for ownership fixtures. The paired process itself must acquire that same exact lock before claiming, perform the post-claim proof and owner checks, then hold it across effects or durably fenced waits. Require a process-to-process named test for first-managed transition, old-build terminal drain, current-owner recovery and upgrade with an active child. Keep stop/uninstall, matched restart and rollback to verified schema-compatible pair ahead of any dispatch. S4b consumes no reservations or actions.
- **S4c-replay.** Re-derived against `c42fe698`; independent scope 4 patch measures 13 changed lines (`$S/S4c-replay.patch`). `TestOperatorConfirmReplay` covers identical replay and signed-artifact mismatch, preserving the row and passkey count. Keep the focus of the roughly 900-line cut on scopes 1–3. Before S7, land the fix and make an exact replay return a link with zero second claim.
- **S3b.** Re-applied and measured 132 lines (`$S/S3b.patch`). `TestOperatorStepRestart` covers the fixture reopen, reservation and receipt case only. The separate sub-milestone must include fd-transfer tests, each receipt-predicate mutation and the current-owner recovery gates. Its current sketch has no production consumer, no fd test and no new handoff gate. The 132-line measurement bounds the sketch only.
- **J-OWN-UID.** At HEAD the guard checks the actual `Stat_t.Uid` and ancestor trust in `runnerPath`. Keep an independently provisioned alternate-owner fixture for a controller that can create one; no privileged test or mock is included. Retain M4 attended cancel and existing kernel/start pairing authority.
- **Earlier S4 Appendix C items.** Keep untouched/reserved recovery, in-flight stop/reconcile, last-action disconnect status, auth readiness E2, and the listed judge-428 cap/DDL/identity/write-floor/FK coverage open on their original gates. Uncertain results require attended reconciliation. S7 owns route registration and replay links, S8 the panel and its browser exits, S9–S11 the dual-proof finish path, S12 the attended handoff. Each must consume guard evidence from the corresponding prior milestones.

## Appendix P — compiled patches

### Patch M1

````diff M1
diff --git a/cmd/ailang-worldd/episode.go b/cmd/ailang-worldd/episode.go
--- a/cmd/ailang-worldd/episode.go
+++ b/cmd/ailang-worldd/episode.go
@@ -361,6 +361,11 @@
 	if err := checkDaemon(health, req, hostPort, self); err != nil {
 		return fail("refusing: %v", err)
 	}
+	effectLock, err := terminalEpisodeLock(req.credPath, episode, nil)
+	if err != nil {
+		return fail("refusing: %v", err)
+	}
+	defer effectLock.Close()
 	if *request {
 		return createMintRequest(req, health.Build, hostPort, *agent, *model, taskPath, mcpFile, stdout, stderr, env)
 	}
diff --git a/cmd/ailang-worldd/episode_finish.go b/cmd/ailang-worldd/episode_finish.go
--- a/cmd/ailang-worldd/episode_finish.go
+++ b/cmd/ailang-worldd/episode_finish.go
@@ -125,6 +125,11 @@
 	if same, _ := store.SameDatabase(h.CredentialsDB, credentials); !same {
 		return fail(fmt.Sprintf("the daemon serves credentials store %s; preset %s names %s (F6)", h.CredentialsDB, *preset, credentials), "the operator proof must be beside the served credentials store", "fix the preset daemon")
 	}
+	effectLock, err := terminalEpisodeLock(credentials, ep, nil)
+	if err != nil {
+		return fail(err.Error(), "the terminal could not acquire exclusive access to the episode", "ailang-worldd service runner status")
+	}
+	defer effectLock.Close()
 	_, rep, err := fetchReport(ctx, c, ep)
 	if err != nil && !errors.Is(err, errNoCalls) {
 		return fail(err.Error(), "finish needs the episode report", "ailang-worldd episode report "+ep)
diff --git a/cmd/ailang-worldd/episode_launch.go b/cmd/ailang-worldd/episode_launch.go
--- a/cmd/ailang-worldd/episode_launch.go
+++ b/cmd/ailang-worldd/episode_launch.go
@@ -250,12 +250,12 @@
 	if len(pos) != 1 || !workspaceEpisodePattern.MatchString(pos[0]) || *name == "" || *name == "se-tools" {
 		return fail(errors.New("episode and --preset <name> are required"))
 	}
-	return launchEpisode(pos[0], *name, *dir, *headless, stdout, stderr, *hostsDir)
+	return launchEpisode(pos[0], *name, *dir, *headless, stdout, stderr, nil, *hostsDir)
 }
 
 // launchEpisode is shared by episode launch and the post-mint start path.
 // Its preflight has no writes; only a fully validated launch takes the lock.
-func launchEpisode(ep, name, presetsDir string, headless bool, stdout, stderr io.Writer, hostsDirs ...string) int {
+func launchEpisode(ep, name, presetsDir string, headless bool, stdout, stderr io.Writer, held *os.File, hostsDirs ...string) int {
 	fail := func(err error) int {
 		msg := err.Error()
 		if !strings.Contains(msg, "\n  why:") {
@@ -276,6 +276,22 @@
 	}
 	if same, _ := store.SameDatabase(h.DBPath, p.DB); !same {
 		return fail(presetRefusal(fmt.Sprintf("the daemon serves store %s; preset %s names %s (F6)", h.DBPath, name, p.DB), "a preset must bind the served store", "fix the preset daemon, or run service status"))
+	}
+	credsPath, err := credentialsFor(p.DB, p.CredentialsDB)
+	if err != nil {
+		return fail(err)
+	}
+	if same, _ := store.SameDatabase(h.CredentialsDB, credsPath); !same {
+		return fail(presetRefusal("daemon credentials store differs from the preset (F6)", "a preset must bind the served credentials store", "fix the preset daemon, or run service status"))
+	}
+	// Check a present owner before report I/O. An unmanaged preflight may
+	// refuse without leaving a new sidecar behind.
+	effectLock, err := openTerminalEpisodeLock(credsPath, ep, held, 0)
+	if err != nil {
+		return fail(err)
+	}
+	if effectLock != nil {
+		defer effectLock.Close()
 	}
 	c := launchClient(p.Addr)
 	_, rep, err := fetchReport(ctx, c, ep)
@@ -421,10 +437,6 @@
 			return fail(errors.New("launch argv would carry the session token"))
 		}
 	}
-	credsPath, err := credentialsFor(p.DB, p.CredentialsDB)
-	if err != nil {
-		return fail(err)
-	}
 	// Check local history before starting; persistence uses the writer lock
 	// only after exit, so the daemon can continue reading while the agent runs.
 	ro, err := store.OpenCredentialsReadOnly(credsPath)
@@ -438,6 +450,14 @@
 	}
 	if localStart == nil {
 		return fail(errors.New("local credentials store has no start record"))
+	}
+	if effectLock == nil && held == nil {
+		// The runner may have acquired ownership while preflight was reading.
+		effectLock, err = terminalEpisodeLock(credsPath, ep, nil)
+		if err != nil {
+			return fail(err)
+		}
+		defer effectLock.Close()
 	}
 	f, err := os.OpenFile(lock, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
 	if err != nil {
diff --git a/cmd/ailang-worldd/episode_launch_test.go b/cmd/ailang-worldd/episode_launch_test.go
--- a/cmd/ailang-worldd/episode_launch_test.go
+++ b/cmd/ailang-worldd/episode_launch_test.go
@@ -8,6 +8,7 @@
 	"os"
 	"path/filepath"
 	"reflect"
+	"runtime/debug"
 	"strings"
 	"testing"
 	"time"
@@ -67,7 +68,9 @@
 	oldClient, oldHealth := launchClient, fetchDaemonHealth
 	t.Cleanup(func() { launchClient, fetchDaemonHealth = oldClient, oldHealth })
 	launchClient = func(string) *client { return memoryDaemon(f.serve) }
-	fetchDaemonHealth = func(string) (daemon.HealthResponse, error) { return daemon.HealthResponse{DBPath: f.db}, nil }
+	fetchDaemonHealth = func(string) (daemon.HealthResponse, error) {
+		return daemon.HealthResponse{DBPath: f.db, CredentialsDB: f.creds}, nil
+	}
 	return f
 }
 
@@ -163,6 +166,9 @@
 }
 
 func TestLaunchHeadlessClaudeAC9AndSecondLaunch(t *testing.T) {
+	// Keep a leaked lock descriptor observable instead of letting GC close it.
+	gc := debug.SetGCPercent(-1)
+	t.Cleanup(func() { debug.SetGCPercent(gc) })
 	f := newLaunchFixture(t, "claude")
 	// The stale mint config is a regular private file and must be replaced.
 	mcp := filepath.Join(f.sessions, "ep.mcp.json")
@@ -234,6 +240,11 @@
 	if code, _, errOut := f.run(); code != 0 {
 		t.Fatal(code, errOut)
 	}
+	effectLock, lockErr := terminalEpisodeLock(f.creds, "ep", nil)
+	if lockErr != nil || effectLock == nil {
+		t.Fatalf("standalone launch retained the episode lock: %v", lockErr)
+	}
+	effectLock.Close()
 	second := f.read()
 	if len(second.Launches) != 2 || !reflect.DeepEqual(first.Launches[0], second.Launches[0]) {
 		t.Fatal("earlier launch changed", second)
diff --git a/cmd/ailang-worldd/episode_preset.go b/cmd/ailang-worldd/episode_preset.go
--- a/cmd/ailang-worldd/episode_preset.go
+++ b/cmd/ailang-worldd/episode_preset.go
@@ -279,6 +279,11 @@
 	if err := checkPresetDaemon(h, req, p, p.Daemon.Addr); err != nil {
 		return fail(err)
 	}
+	effectLock, err := terminalEpisodeLock(req.credPath, ep, nil)
+	if err != nil {
+		return fail(err)
+	}
+	defer effectLock.Close()
 	req.afterToken = func(token string) error {
 		if taskRecord != nil {
 			f, err := os.OpenFile(taskCopy, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
@@ -308,7 +313,7 @@
 		return code
 	}
 	if launching {
-		return launchEpisode(ep, name, dir, headless, out, errOut, fs.Lookup("hosts-dir").Value.String())
+		return launchEpisode(ep, name, dir, headless, out, errOut, effectLock, fs.Lookup("hosts-dir").Value.String())
 	}
 	taskPath := ""
 	if taskRecord != nil {
diff --git a/cmd/ailang-worldd/episode_resume.go b/cmd/ailang-worldd/episode_resume.go
--- a/cmd/ailang-worldd/episode_resume.go
+++ b/cmd/ailang-worldd/episode_resume.go
@@ -172,6 +172,11 @@
 	if err := checkDaemon(h, req, p.Addr, episodeSelfBuild()); err != nil {
 		return fail(err)
 	}
+	effectLock, err := terminalEpisodeLock(req.credPath, ep, nil)
+	if err != nil {
+		return fail(err)
+	}
+	defer effectLock.Close()
 	presetPath, err := expandPresetPath(filepath.Join(presetsDir, name+".json"))
 	if err != nil {
 		return fail(err)
diff --git a/cmd/ailang-worldd/operator_episode_lock.go b/cmd/ailang-worldd/operator_episode_lock.go
new file mode 100644
--- /dev/null
+++ b/cmd/ailang-worldd/operator_episode_lock.go
@@ -0,0 +1,60 @@
+package main
+
+import (
+	"errors"
+	"fmt"
+	"os"
+	"syscall"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+func operatorEpisodeLockPath(canonical, episode string) string {
+	return canonical + ".operator-episode-" + episode + ".lock"
+}
+
+func terminalEpisodeLock(creds, episode string, held *os.File) (*os.File, error) {
+	return openTerminalEpisodeLock(creds, episode, held, os.O_CREATE)
+}
+
+func openTerminalEpisodeLock(creds, episode string, held *os.File, flags int) (*os.File, error) {
+	canonical, err := store.CanonicalPath(creds)
+	if err != nil {
+		return nil, err
+	}
+	path := operatorEpisodeLockPath(canonical, episode)
+	if held != nil {
+		if held.Name() != path {
+			return nil, errors.New("operator store or episode changed during episode start")
+		}
+		return nil, nil // start retains and closes the same open file description
+	}
+
+	// Take the same lock even before the first companion startup. A first
+	// runner must not race a terminal operation that saw an unmanaged store.
+	lock, lockErr := openRunnerLock(path, flags)
+	absent := flags == 0 && errors.Is(lockErr, os.ErrNotExist)
+	if lockErr != nil && !absent && !errors.Is(lockErr, syscall.EWOULDBLOCK) {
+		return nil, lockErr
+	}
+	ctx, cancel := budgetContext(4 * time.Second)
+	defer cancel()
+	id, readErr := store.ReadOperatorEpisodeOwner(ctx, canonical, episode)
+	if (lockErr == nil || absent) && readErr == nil && id == "" {
+		return lock, nil
+	}
+	if lock != nil {
+		lock.Close()
+	}
+
+	switch {
+	case readErr != nil || id != "" && !verbs.ValidOperatorActionID(id):
+		return nil, fmt.Errorf("%w: could not read action ownership; run ailang-worldd service runner status", store.ErrOperatorConflict)
+	case id != "":
+		return nil, fmt.Errorf("%w: active action %s; run ailang-worldd service runner status and reconcile at the terminal", store.ErrOperatorConflict, id)
+	default:
+		return nil, fmt.Errorf("%w: episode lock held; run ailang-worldd service runner status", store.ErrOperatorConflict)
+	}
+}
diff --git a/cmd/ailang-worldd/operator_episode_lock_test.go b/cmd/ailang-worldd/operator_episode_lock_test.go
new file mode 100644
--- /dev/null
+++ b/cmd/ailang-worldd/operator_episode_lock_test.go
@@ -0,0 +1,382 @@
+package main
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"errors"
+	"io"
+	"net/http"
+	"os"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"syscall"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+const fixtureActionID = "oa-0123456789abcdef0123456789abcdef"
+
+func claimFixtureEpisode(t *testing.T, creds, episode string) *os.File {
+	t.Helper()
+	canonical, err := store.CanonicalPath(creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	held, err := runnerKernelLock(operatorEpisodeLockPath(canonical, episode))
+	if err != nil {
+		t.Fatal(err)
+	}
+	// Release a holder if a regression makes the CLI wait on flock.
+	fuse := time.AfterFunc(2*time.Second, func() { held.Close() })
+	t.Cleanup(func() { fuse.Stop(); held.Close() })
+	if _, err := store.UpgradeCredentials(creds); err != nil {
+		t.Fatal(err)
+	}
+	db, err := store.OpenCredentials(creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer db.Close()
+	ctx := boundedTestContext(t)
+	key := approvaltest.NewSigner(51)
+	if err := db.EnrolApprover(ctx, store.ApproverRow{Name: "fixture", Format: "ssh", PublicKey: key.PublicKey().Blob, EnrolledAt: 900}); err != nil {
+		t.Fatal(err)
+	}
+	a := approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: fixtureActionID, Kind: "start", Episode: episode,
+		Store: canonical, Options: approval.OperatorOptions{Preset: "p", Headless: true}, Task: "fixture", Review: "fixture",
+		Resolution: []approval.OperatorResolved{{Name: "agent", Value: "codex"}}}
+	a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
+	req := store.OperatorRequest{Action: a, OperationID: strings.Repeat("1", 32), Challenge: sha256.Sum256([]byte(t.Name())), CreatedAt: 1000, ExpiresAt: 1600}
+	row, err := db.RequestOperatorAction(ctx, req, 1000)
+	if err != nil {
+		t.Fatal(err)
+	}
+	cas := store.OperatorCAS{ID: row.ID, Digest: [32]byte(row.Digest), Generation: row.Generation}
+	in := authority.OperatorJudgeInput{Action: a, Digest: cas.Digest, BytesDigest: sha256.Sum256(row.Action), Challenge: req.Challenge,
+		CreatedAt: req.CreatedAt, ExpiresAt: req.ExpiresAt, Now: 1010}
+	message, err := approval.OperatorActionMessage(cas.Digest, req.Challenge, req.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	artifact := key.Sign(approval.OperatorActionNamespace, message)
+	proof, err := authority.JudgeOperatorAction(ctx, db, in, artifact)
+	if err != nil {
+		t.Fatal(err)
+	}
+	row, err = db.ConfirmOperatorAction(ctx, cas, artifact, proof, 1010)
+	if err != nil {
+		t.Fatal(err)
+	}
+	h := store.OperatorHeartbeat{Owner: store.OperatorOwner{Nonce: "test-runner"}, Protocol: store.OperatorRunnerProtocol, Build: "fixture", Readiness: "start", SeenAt: 1010}
+	if h.Owner.Generation, err = db.PutOperatorHeartbeat(ctx, h); err != nil {
+		t.Fatal(err)
+	}
+	cas.Generation = row.Generation
+	if err := db.ClaimOperatorAction(ctx, cas, h.Owner, 1011); err != nil {
+		t.Fatal(err)
+	}
+	return held
+}
+
+func expectEpisodeRefusal(t *testing.T, code int, output string, before, creds string) {
+	t.Helper()
+	if code != exitUsage || !strings.Contains(output, "active action "+fixtureActionID+";") ||
+		!strings.Contains(output, "service runner status") || serviceFileHash(creds) != before {
+		t.Fatalf("expected exact owner refusal with unchanged credentials: %d %s", code, output)
+	}
+}
+
+func TestOperatorTerminalExclusion(t *testing.T) {
+	t.Run("nonwaiting", testOperatorEpisodeLockTransition)
+	t.Run("plain", func(t *testing.T) {
+		r := newEpisodeRun(t)
+		health := matchingHealth(t, r.f, r.profiles)
+		stubDaemon(t, health, nil)
+		claimFixtureEpisode(t, health.CredentialsDB, "ep1")
+		before, trees := serviceFileHash(health.CredentialsDB), r.f.worktrees(t)
+		var out bytes.Buffer
+		opened := false
+		code := runEpisodeStart(r.args("ep1"), &out, &out, spyEnv(&opened))
+		expectEpisodeRefusal(t, code, out.String(), before, health.CredentialsDB)
+		r.assertNoSessionFiles(t)
+		if opened || r.f.worktrees(t) != trees {
+			t.Fatal("terminal opened or worktree changed")
+		}
+	})
+	t.Run("preset", func(t *testing.T) {
+		f := newPresetFixture(t)
+		creds, _ := credentialsFor(f.r.f.db, "")
+		claimFixtureEpisode(t, creds, "p-2610281352")
+		before, trees := serviceFileHash(creds), f.r.f.worktrees(t)
+		var out bytes.Buffer
+		opened := false
+		env := spyEnv(&opened)
+		env.now = presetClock
+		code := runEpisodeStart(f.args("--task", f.task), &out, &out, env)
+		expectEpisodeRefusal(t, code, out.String(), before, creds)
+		f.r.assertNoSessionFiles(t)
+		if opened || f.r.f.worktrees(t) != trees {
+			t.Fatal("terminal opened or worktree changed")
+		}
+	})
+	t.Run("resume", func(t *testing.T) {
+		f, ep, _ := resumeFixture(t)
+		creds, _ := credentialsFor(f.r.f.db, "")
+		held := claimFixtureEpisode(t, creds, ep)
+		other, err := terminalEpisodeLock(creds, "other", nil)
+		if err != nil || other == nil {
+			t.Fatalf("durable owner of %s blocked another episode: %v", ep, err)
+		}
+		other.Close()
+		if err := held.Close(); err != nil {
+			t.Fatal(err)
+		}
+		before, trees := serviceFileHash(creds), f.r.f.worktrees(t)
+		token := filepath.Join(f.r.sessions, ep+".session")
+		tokenHash := serviceFileHash(token)
+		now := presetClock() + 43201
+		var out bytes.Buffer
+		opened := false
+		env := spyEnv(&opened)
+		env.now = func() int64 { return now }
+		code := runEpisodeStart(f.args("--resume", ep), &out, &out, env)
+		expectEpisodeRefusal(t, code, out.String(), before, creds)
+		archived, err := filepath.Glob(token + ".expired-*")
+		if err != nil || len(archived) != 0 || opened || serviceFileHash(token) != tokenHash || f.r.f.worktrees(t) != trees {
+			t.Fatal("expired approval or missing kernel holder permitted a resume effect")
+		}
+	})
+	t.Run("launch", func(t *testing.T) {
+		f := newLaunchFixture(t, "claude")
+		claimFixtureEpisode(t, f.creds, "ep")
+		before := serviceFileHash(f.creds)
+		token := serviceFileHash(filepath.Join(f.sessions, "ep.session"))
+		launchClient = func(string) *client { t.Fatal("launch reached report dispatch"); return nil }
+		code, _, output := f.run()
+		expectEpisodeRefusal(t, code, output, before, f.creds)
+		if files, err := os.ReadDir(f.facts); err != nil || len(files) != 0 || token != serviceFileHash(filepath.Join(f.sessions, "ep.session")) {
+			t.Fatal("launch modified session files or spawned")
+		}
+		if _, err := os.Lstat(filepath.Join(f.sessions, "ep.launch.lock")); !errors.Is(err, os.ErrNotExist) {
+			t.Fatalf("launch refusal left a PID file: %v", err)
+		}
+	})
+	t.Run("finish", func(t *testing.T) {
+		t.Setenv("HOME", t.TempDir())
+		f := newFinishFixture(t)
+		creds := f.health.CredentialsDB
+		claimFixtureEpisode(t, creds, "ep1")
+		before, status := serviceFileHash(creds), f.gitOutput("-C", f.wt, "status", "--porcelain")
+		opened := false
+		env := sessionEnv{openTerminal: func() (io.ReadWriteCloser, error) {
+			if len(f.requests) != 0 {
+				t.Fatal("health preceded F1")
+			}
+			opened = true
+			return &fakeTerm{r: strings.NewReader("n\n")}, nil
+		}}
+		code, _, output := f.run(env)
+		expectEpisodeRefusal(t, code, output, before, creds)
+		if !opened || !reflect.DeepEqual(f.requests, []string{"health"}) || f.checks != 0 || f.outcomes != 0 ||
+			strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != f.before || f.gitOutput("-C", f.wt, "status", "--porcelain") != status {
+			t.Fatalf("F1 or finish effect order changed: %v", f.requests)
+		}
+	})
+}
+
+func testOperatorEpisodeLockTransition(t *testing.T) {
+	creds := filepath.Join(t.TempDir(), "credentials.db")
+	canonical, err := store.CanonicalPath(creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	initial, err := terminalEpisodeLock(creds, "ep", nil)
+	if err != nil || initial == nil {
+		t.Fatalf("pre-management acquisition: %v %v", initial, err)
+	}
+	defer initial.Close()
+	process, err := runnerKernelLock(canonical + ".operator-runner.lock")
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer process.Close()
+	// An independently opened fd represents the future paired process's
+	// acquisition. Release on timeout so a blocking mutation cannot hang.
+	results := make(chan error, 1)
+	go func() {
+		contender, err := runnerKernelLock(operatorEpisodeLockPath(canonical, "ep"))
+		if contender != nil {
+			contender.Close()
+		}
+		results <- err
+	}()
+	select {
+	case err := <-results:
+		if !errors.Is(err, syscall.EWOULDBLOCK) {
+			t.Fatalf("same episode did not refuse: %v", err)
+		}
+	case <-time.After(time.Second):
+		initial.Close()
+		t.Fatal("episode acquisition waited on the first terminal holder")
+	}
+	if lock, err := terminalEpisodeLock(creds, "ep", nil); lock != nil || !errors.Is(err, store.ErrOperatorConflict) ||
+		!strings.Contains(err.Error(), "episode lock held") {
+		t.Fatalf("unidentified holder: %v %v", lock, err)
+	}
+	if _, err := os.Stat(creds); !errors.Is(err, os.ErrNotExist) {
+		t.Fatalf("refusal created a database: %v", err)
+	}
+	for _, tc := range []struct{ creds, episode string }{{creds, "other"}, {filepath.Join(filepath.Dir(creds), "other.db"), "ep"}} {
+		fd, err := terminalEpisodeLock(tc.creds, tc.episode, nil)
+		if err != nil || fd == nil {
+			t.Fatalf("independent store/episode denied: %v %v", fd, err)
+		}
+		fd.Close()
+	}
+	if _, err := terminalEpisodeLock(creds, "other", initial); err == nil {
+		t.Fatal("transferred lock to another episode")
+	}
+	alias := filepath.Join(t.TempDir(), "alias")
+	if err := os.Symlink(filepath.Dir(creds), alias); err != nil {
+		t.Fatal(err)
+	}
+	if fd, err := terminalEpisodeLock(filepath.Join(alias, filepath.Base(creds)), "ep", initial); fd != nil || err != nil {
+		t.Fatalf("canonical transfer: %v %v", fd, err)
+	}
+	if err := initial.Close(); err != nil {
+		t.Fatal(err)
+	}
+	fd, err := runnerKernelLock(operatorEpisodeLockPath(canonical, "ep"))
+	if err != nil {
+		t.Fatal(err)
+	}
+	fd.Close()
+}
+
+func TestOperatorLaunchServedCredentials(t *testing.T) {
+	f := newLaunchFixture(t, "claude")
+	fetchDaemonHealth = func(string) (daemon.HealthResponse, error) {
+		return daemon.HealthResponse{DBPath: f.db, CredentialsDB: filepath.Join(f.root, "other.db")}, nil
+	}
+	launchClient = func(string) *client { t.Fatal("mismatch reached report dispatch"); return nil }
+	before := serviceFileHash(f.creds)
+	code, _, output := f.run()
+	if code != exitUsage || !strings.Contains(output, "credentials store differs") || serviceFileHash(f.creds) != before {
+		t.Fatalf("credentials pair mismatch: %d %s", code, output)
+	}
+}
+
+func TestOperatorEpisodeLockRefusals(t *testing.T) {
+	for _, arm := range []string{"unsafe lock", "ownership read", "invalid action id", "unclaimed action"} {
+		t.Run(arm, func(t *testing.T) {
+			creds := filepath.Join(t.TempDir(), "credentials.db")
+			canonical, err := store.CanonicalPath(creds)
+			if err != nil {
+				t.Fatal(err)
+			}
+			lockPath := operatorEpisodeLockPath(canonical, "ep")
+			want, change := "could not read action ownership", ""
+			switch arm {
+			case "unsafe lock":
+				want = "unsafe runner lock file"
+				if err := os.WriteFile(lockPath, []byte("keep"), 0600); err != nil {
+					t.Fatal(err)
+				}
+				if err := os.Chmod(lockPath, 0644); err != nil {
+					t.Fatal(err)
+				}
+			case "ownership read":
+				change = "DROP TABLE operator_actions"
+			case "invalid action id":
+				change = "PRAGMA foreign_keys=OFF; UPDATE operator_actions SET id='invalid'"
+			case "unclaimed action":
+				change = "UPDATE operator_actions SET state='confirmed', owner=NULL"
+			}
+			if change != "" {
+				held := claimFixtureEpisode(t, creds, "ep")
+				held.Close()
+				credSQL(t, creds, change)
+			}
+			before, lockBefore := serviceFileHash(creds), serviceFileHash(lockPath)
+			lock, err := terminalEpisodeLock(creds, "ep", nil)
+			if lock != nil {
+				lock.Close()
+			}
+			if arm == "unclaimed action" {
+				if err != nil || lock == nil {
+					t.Fatalf("confirmed action without an owner excludes terminal: %v", err)
+				}
+			} else if err == nil || lock != nil || !strings.Contains(err.Error(), want) {
+				t.Fatalf("refusal: %v %v", lock, err)
+			}
+			if serviceFileHash(creds) != before || serviceFileHash(lockPath) != lockBefore {
+				t.Fatal("store or lock bytes changed")
+			}
+			if arm != "unsafe lock" {
+				fd, err := runnerKernelLock(lockPath)
+				if err != nil {
+					t.Fatalf("terminal failed to close its lock: %v", err)
+				}
+				fd.Close()
+			}
+		})
+	}
+}
+
+func TestOperatorEpisodeLockExclusive(t *testing.T) {
+	path := filepath.Join(t.TempDir(), "lock")
+	holder, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer holder.Close()
+	if err := syscall.Flock(int(holder.Fd()), syscall.LOCK_SH|syscall.LOCK_NB); err != nil {
+		t.Fatal(err)
+	}
+	result := make(chan error, 1)
+	go func() {
+		lock, err := runnerKernelLock(path)
+		if lock != nil {
+			lock.Close()
+		}
+		result <- err
+	}()
+	select {
+	case err := <-result:
+		if !errors.Is(err, syscall.EWOULDBLOCK) {
+			t.Fatalf("runner lock did not exclude a shared holder: %v", err)
+		}
+	case <-time.After(time.Second):
+		holder.Close()
+		t.Fatal("runner lock waited on a shared holder")
+	}
+}
+
+func TestOperatorLaunchClaimDuringPreflight(t *testing.T) {
+	f := newLaunchFixture(t, "claude")
+	var before string
+	launchClient = func(string) *client {
+		return memoryDaemon(func(w http.ResponseWriter, r *http.Request) {
+			if r.URL.Path == "/v1/episodes/ep/report" {
+				claimFixtureEpisode(t, f.creds, "ep")
+				before = serviceFileHash(f.creds)
+			}
+			f.serve(w, r)
+		})
+	}
+	code, _, output := f.run()
+	expectEpisodeRefusal(t, code, output, before, f.creds)
+	if facts, err := os.ReadDir(f.facts); err != nil || len(facts) != 0 {
+		t.Fatal("claim during preflight allowed a child")
+	}
+	if _, err := os.Lstat(filepath.Join(f.sessions, "ep.launch.lock")); !errors.Is(err, os.ErrNotExist) {
+		t.Fatalf("claim during preflight left a PID file: %v", err)
+	}
+}
diff --git a/cmd/ailang-worldd/operator_runner.go b/cmd/ailang-worldd/operator_runner.go
--- a/cmd/ailang-worldd/operator_runner.go
+++ b/cmd/ailang-worldd/operator_runner.go
@@ -19,7 +19,11 @@
 const runnerPoll = 2 * time.Second
 
 func runnerKernelLock(path string) (*os.File, error) {
-	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0600)
+	return openRunnerLock(path, os.O_CREATE)
+}
+
+func openRunnerLock(path string, flags int) (*os.File, error) {
+	f, err := os.OpenFile(path, flags|os.O_RDWR|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0600)
 	if err != nil {
 		return nil, err
 	}
diff --git a/host/store/operator_runner.go b/host/store/operator_runner.go
--- a/host/store/operator_runner.go
+++ b/host/store/operator_runner.go
@@ -4,9 +4,12 @@
 	"context"
 	"database/sql"
 	"errors"
+	"net/url"
+	"os"
 )
 
 const OperatorRunnerProtocol = 1
+const operatorSchemaVersion = 10
 
 type OperatorHeartbeat struct {
 	Owner                       OperatorOwner
@@ -78,3 +81,57 @@
 			WHERE action_id=? AND kind<>'operation')`, now, id, now, id)
 	return operatorChangedOne(result, err)
 }
+
+// ReadOperatorEpisodeOwner reads exclusion independently of lease liveness.
+// The caller holds the episode lock or has just received a non-waiting
+// contention error. A missing store or older schema has no operator rows.
+func ReadOperatorEpisodeOwner(ctx context.Context, path, episode string) (string, error) {
+	if err := requireDeadline(ctx); err != nil {
+		return "", err
+	}
+	db, version, err := openOperatorReadOnly(ctx, path)
+	if errors.Is(err, os.ErrNotExist) {
+		return "", nil
+	}
+	if err != nil {
+		return "", err
+	}
+	defer db.Close()
+	if version < operatorSchemaVersion {
+		return "", nil
+	}
+	var id string
+	err = db.QueryRowContext(ctx, `SELECT id FROM operator_actions WHERE episode_id=? AND owner IS NOT NULL`, episode).Scan(&id)
+	if errors.Is(err, sql.ErrNoRows) {
+		return "", nil
+	}
+	return id, err
+}
+
+func openOperatorReadOnly(ctx context.Context, path string) (*sql.DB, int, error) {
+	canonical, _, err := resolveDSN(path)
+	if err != nil {
+		return nil, 0, err
+	}
+	st, err := os.Stat(canonical)
+	if err != nil {
+		return nil, 0, err
+	}
+	if !st.Mode().IsRegular() {
+		return nil, 0, errors.New("expected a regular operator credentials file")
+	}
+	db, err := sql.Open("sqlite", readOnlyDSN(canonical, url.Values{"_pragma": {"busy_timeout(0)"}}))
+	if err != nil {
+		return nil, 0, err
+	}
+	var version int
+	if err := db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); err != nil {
+		db.Close()
+		return nil, 0, err
+	}
+	if version > credentialsSchemaVersion {
+		db.Close()
+		return nil, 0, &CredentialsSchemaError{Path: canonical, Detail: "unsupported operator credentials schema"}
+	}
+	return db, version, nil
+}
diff --git a/host/store/operator_schema_test.go b/host/store/operator_schema_test.go
--- a/host/store/operator_schema_test.go
+++ b/host/store/operator_schema_test.go
@@ -1,8 +1,10 @@
 package store
 
 import (
+	"bytes"
 	"database/sql"
 	"errors"
+	"os"
 	"path/filepath"
 	"reflect"
 	"testing"
@@ -85,3 +87,67 @@
 		t.Fatalf("lost v9 credential: %+v %v %v", got, ok, err)
 	}
 }
+
+func TestOperatorEpisodeOwnerReadSchema(t *testing.T) {
+	ctx := testCtx(t)
+	path := filepath.Join(t.TempDir(), CredentialsFileName)
+	if id, err := ReadOperatorEpisodeOwner(ctx, path, "ep"); id != "" || err != nil {
+		t.Fatalf("absent store: %q %v", id, err)
+	}
+	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
+		t.Fatalf("absent store changed: %v", err)
+	}
+	if _, err := ReadOperatorEpisodeOwner(ctx, "", "ep"); err == nil {
+		t.Fatal("accepted an empty database path")
+	}
+	if _, err := ReadOperatorEpisodeOwner(ctx, filepath.Dir(path), "ep"); err == nil || err.Error() != "expected a regular operator credentials file" {
+		t.Fatal("accepted a directory as a database")
+	}
+	raw := rawDB(t, path)
+	if _, err := raw.Exec(credentialsSchemaSQL); err != nil {
+		t.Fatal(err)
+	}
+	setPositiveFixtureVersion(t, raw, 9)
+	if id, err := ReadOperatorEpisodeOwner(ctx, path, "ep"); id != "" || err != nil {
+		t.Fatalf("v9 owner lookup: %q %v", id, err)
+	}
+	if _, err := UpgradeCredentials(path); err != nil {
+		t.Fatal(err)
+	}
+	if id, err := ReadOperatorEpisodeOwner(ctx, path, "ep"); id != "" || err != nil {
+		t.Fatalf("empty v10 owner lookup: %q %v", id, err)
+	}
+	ro, _, err := openOperatorReadOnly(ctx, path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	var busy int
+	busyErr := ro.QueryRowContext(ctx, "PRAGMA busy_timeout").Scan(&busy)
+	_, writeErr := ro.ExecContext(ctx, "DELETE FROM operator_actions")
+	ro.Close()
+	if busyErr != nil || busy != 0 || writeErr == nil {
+		t.Fatalf("deadline or read-only open drift: busy=%d %v write=%v", busy, busyErr, writeErr)
+	}
+	setPositiveFixtureVersion(t, raw, 11)
+	if _, err := ReadOperatorEpisodeOwner(ctx, path, "ep"); err == nil || userVersion(t, path) != 11 {
+		t.Fatalf("future schema accepted or changed: %v", err)
+	}
+	corrupt := filepath.Join(t.TempDir(), "corrupt.db")
+	if _, err := UpgradeCredentials(corrupt); err != nil {
+		t.Fatal(err)
+	}
+	data, err := os.ReadFile(corrupt)
+	if err != nil || len(data) == 0 {
+		t.Fatalf("read migrated fixture: %v", err)
+	}
+	data[0] ^= 0xff
+	if err := os.WriteFile(corrupt, data, 0600); err != nil {
+		t.Fatal(err)
+	}
+	if _, err := ReadOperatorEpisodeOwner(ctx, corrupt, "ep"); err == nil {
+		t.Fatal("corrupt database accepted")
+	}
+	if got, err := os.ReadFile(corrupt); err != nil || !bytes.Equal(got, data) {
+		t.Fatalf("corrupt database changed: %v", err)
+	}
+}
````

### Patch M2

````diff M2
diff --git a/cmd/ailang-worldd/operator_runner.go b/cmd/ailang-worldd/operator_runner.go
--- a/cmd/ailang-worldd/operator_runner.go
+++ b/cmd/ailang-worldd/operator_runner.go
@@ -99,6 +99,32 @@
 	})
 }
 
+func runnerHeartbeat(path string, budget time.Duration) (store.OperatorHeartbeat, error) {
+	ctx, cancel := budgetContext(budget)
+	defer cancel()
+	return store.ReadOperatorHeartbeat(ctx, path)
+}
+
+func (e serviceEnv) waitForRunner(c runnerConfig, previous int64) bool {
+	deadline := e.runnerNow().Add(e.wait)
+	for {
+		remaining := deadline.Sub(e.runnerNow())
+		if remaining <= 0 {
+			return false
+		}
+		h, err := runnerHeartbeat(c.Credentials, min(remaining, time.Second))
+		now := e.runnerNow()
+		if !now.Before(deadline) {
+			return false
+		}
+		if err == nil && h.Owner.Generation > previous &&
+			h.Available(c.Executables[e.paths.Binary], "start", now.Unix()) {
+			return true
+		}
+		e.runnerSleep(min(deadline.Sub(now), 100*time.Millisecond))
+	}
+}
+
 func (e serviceEnv) runnerStatus(w io.Writer) error {
 	c, _, err := e.installedRunner()
 	if err != nil {
diff --git a/cmd/ailang-worldd/operator_service.go b/cmd/ailang-worldd/operator_service.go
--- a/cmd/ailang-worldd/operator_service.go
+++ b/cmd/ailang-worldd/operator_service.go
@@ -13,6 +13,7 @@
 	"strconv"
 	"strings"
 	"syscall"
+	"time"
 
 	"github.com/sunholo-data/ailang-world/host/store"
 	"github.com/sunholo-data/ailang-world/host/verbs"
@@ -431,13 +432,25 @@
 	_, err = e.controller.Print(e.domain, runnerLabel)
 	switch {
 	case errors.Is(err, errServiceNotLoaded):
+		previous, err := runnerHeartbeat(c.Credentials, time.Second)
+		if err != nil && !errors.Is(err, store.ErrOperatorRunner) {
+			return errors.New("cannot determine the previous runner lease")
+		}
 		if err := e.controller.Bootstrap(e.domain, p.Plist); err != nil {
 			_ = e.controller.Bootout(e.domain, runnerLabel)
 
 			return errors.New("could not bootstrap the companion")
 		}
+		if !e.waitForRunner(c, previous.Owner.Generation) {
+			if err := e.controller.Bootout(e.domain, runnerLabel); err != nil && !errors.Is(err, errServiceNotLoaded) {
+				return errors.New("runner readiness timed out and companion could not be stopped")
+			}
+			return errors.New("runner readiness timed out; companion stopped")
+		}
 	case err != nil:
 		return errors.New("could not inspect the companion")
+	case !e.waitForRunner(c, 0):
+		return errors.New("loaded companion heartbeat unavailable")
 	}
 	return nil
 }
diff --git a/cmd/ailang-worldd/operator_service_test.go b/cmd/ailang-worldd/operator_service_test.go
--- a/cmd/ailang-worldd/operator_service_test.go
+++ b/cmd/ailang-worldd/operator_service_test.go
@@ -55,6 +55,14 @@
 		health:    daemon.HealthResponse{Status: "ok", DBPath: world, CredentialsDB: creds, Build: daemon.BuildIdentity{ExecutableSHA256: r.BinarySHA256}}}
 	f.ctl = &fakeServiceController{runner: true, binary: r.Binary}
 	f.env.controller = f.ctl
+	f.env.wait = 250 * time.Millisecond
+	clock := time.Unix(f.env.terminal.now(), 0)
+	f.env.runnerNow = func() time.Time { return clock }
+	f.env.runnerSleep = func(d time.Duration) { clock = clock.Add(d) }
+	f.ctl.onBoot = func() error {
+		o.heartbeat.Owner.Generation = 0
+		return f.env.runnerTick(o.config, &o.heartbeat)
+	}
 	f.env.client = &http.Client{Transport: serviceTransport(func(req *http.Request) (*http.Response, error) {
 		if req.URL.Host != r.Address || req.URL.Path != "/v1/health" {
 			t.Fatalf("unexpected health URL: %s", req.URL)
@@ -275,6 +283,7 @@
 	if c, s := f.command("start"); c != 0 || !f.ctl.loaded {
 		t.Fatalf("%d %s", c, s)
 	}
+	t.Run("gate", testOperatorServiceStartupGate)
 }
 
 func TestOperatorServiceUpgradePairFence(t *testing.T) {
@@ -283,6 +292,7 @@
 		{"plist-only", "operator runner is still installed"},
 		{"registered", "cannot verify that the operator runner is unloaded"},
 		{"locked", "operator runner lifecycle is busy or unavailable"},
+		{"inspect error", "cannot verify that the operator runner is unloaded"},
 	} {
 		t.Run(tc.arm, func(t *testing.T) {
 			f := operatorService(t)
@@ -304,6 +314,9 @@
 			}
 			f.env.controller = f.homeJob
 			f.homeJob.runner = tc.arm == "registered"
+			if tc.arm == "inspect error" {
+				f.homeJob.runnerPrintErr = errors.New("supervisor query failed")
+			}
 			f.homeJob.onKick = func() error {
 				f.health.Build.ExecutableSHA256 = f.homeJob.running
 				return nil
@@ -414,7 +427,11 @@
 	if info, err := os.Stat(filepath.Dir(p.Log)); err != nil || info.Mode().Perm() != 0700 {
 		t.Fatalf("spool: %v %v", info, err)
 	}
-	for _, verb := range []string{"start", "stop", "uninstall"} {
+	if c, s := f.command("start"); c != 0 {
+		t.Fatalf("start: %d %s", c, s)
+	}
+	want[4] = serviceFileHash(f.config.Credentials)
+	for _, verb := range []string{"stop", "uninstall"} {
 		if c, s := f.command(verb); c != 0 {
 			t.Fatalf("%s: %d %s", verb, c, s)
 		}
diff --git a/cmd/ailang-worldd/operator_startup_test.go b/cmd/ailang-worldd/operator_startup_test.go
new file mode 100644
--- /dev/null
+++ b/cmd/ailang-worldd/operator_startup_test.go
@@ -0,0 +1,131 @@
+package main
+
+import (
+	"errors"
+	"strings"
+	"testing"
+	"time"
+)
+
+func testOperatorServiceStartupGate(t *testing.T) {
+	for _, tc := range []struct {
+		arm, want        string
+		loaded, stopFail bool
+		boots, stops     int
+	}{
+		{arm: "same generation", want: "readiness timed out; companion stopped", boots: 1, stops: 1},
+		{arm: "late", want: "readiness timed out; companion stopped", boots: 1, stops: 1},
+		{arm: "stale", want: "readiness timed out; companion stopped", boots: 1, stops: 1},
+		{arm: "future", want: "readiness timed out; companion stopped", boots: 1, stops: 1},
+		{arm: "wrong build", want: "readiness timed out; companion stopped", boots: 1, stops: 1},
+		{arm: "wrong protocol", want: "readiness timed out; companion stopped", boots: 1, stops: 1},
+		{arm: "finish only", want: "readiness timed out; companion stopped", boots: 1, stops: 1},
+		{arm: "empty owner", want: "readiness timed out; companion stopped", boots: 1, stops: 1},
+		{arm: "bad read", want: "readiness timed out; companion stopped", boots: 1, stops: 1},
+		{arm: "failed stop", want: "companion could not be stopped", boots: 1, stops: 1, stopFail: true},
+		{arm: "already stopped", want: "readiness timed out; companion stopped", boots: 1, stops: 1},
+		{arm: "already loaded stale", want: "loaded companion heartbeat unavailable", loaded: true},
+		{arm: "already loaded live", loaded: true},
+		{arm: "fresh restart", boots: 1},
+		{arm: "bad previous", want: "cannot determine the previous runner lease"},
+		{arm: "unknown supervisor", want: "could not inspect the companion"},
+	} {
+		t.Run(tc.arm, func(t *testing.T) {
+			f := operatorService(t)
+			f.mustInstall(t)
+			current := f.env.terminal.now()
+			clock := time.Unix(current, 0)
+			reads := 0
+			f.env.runnerNow = func() time.Time {
+				reads++
+				if tc.arm == "late" && reads >= 3 {
+					return clock.Add(f.env.wait)
+				}
+				return clock
+			}
+			sleeps := 0
+			f.env.runnerSleep = func(d time.Duration) {
+				sleeps++
+				if sleeps > 3 || d <= 0 || d > 100*time.Millisecond {
+					t.Fatalf("unbounded poll: %d sleeps, delay %v", sleeps, d)
+				}
+				clock = clock.Add(d)
+			}
+			if err := f.env.runnerTick(f.config, &f.heartbeat); err != nil {
+				t.Fatal(err)
+			}
+			creds := f.config.Credentials
+			f.ctl.loaded = tc.loaded
+			if tc.stopFail {
+				f.ctl.stopErr = errors.New("fake stop failed")
+			}
+			if tc.arm == "already stopped" {
+				f.ctl.stopErr = errServiceNotLoaded
+			}
+			if tc.arm == "unknown supervisor" {
+				f.ctl.runnerPrintErr = errors.New("fake supervisor error")
+			}
+			if tc.arm == "already loaded stale" {
+				credSQL(t, creds, "UPDATE operator_runner_state SET seen_at=113")
+			}
+			if tc.arm == "bad previous" {
+				credSQL(t, creds, "UPDATE operator_runner_state SET generation='broken'")
+			}
+			knownHash := serviceFileHash(creds)
+			newBoot := f.ctl.onBoot
+			f.ctl.onBoot = func() error {
+				switch tc.arm {
+				case "same generation", "failed stop", "already stopped":
+					return nil
+				}
+				if err := newBoot(); err != nil {
+					return err
+				}
+				switch tc.arm {
+				case "stale":
+					credSQL(t, creds, "UPDATE operator_runner_state SET seen_at=113")
+				case "future":
+					credSQL(t, creds, "UPDATE operator_runner_state SET seen_at=124")
+				case "wrong build":
+					credSQL(t, creds, "UPDATE operator_runner_state SET build='wrong'")
+				case "wrong protocol":
+					credSQL(t, creds, "UPDATE operator_runner_state SET protocol=0")
+				case "finish only":
+					credSQL(t, creds, "UPDATE operator_runner_state SET readiness='finish'")
+				case "empty owner":
+					credSQL(t, creds, "UPDATE operator_runner_state SET owner=''")
+				case "bad read":
+					credSQL(t, creds, "UPDATE operator_runner_state SET catalogue_generation='broken'")
+				}
+				knownHash = serviceFileHash(creds)
+				return nil
+			}
+			f.writes = nil
+			code, output := f.command("start")
+			switch {
+			case tc.want == "":
+				h, err := runnerHeartbeat(creds, time.Second)
+				wantGen := int64(1 + tc.boots)
+				if code != 0 || !f.ctl.loaded || err != nil || h.Owner.Generation != wantGen || sleeps != 0 {
+					t.Fatalf("fresh success: %d %q loaded=%v heartbeat=%+v err=%v", code, output, f.ctl.loaded, h, err)
+				}
+			case code != 1 || !strings.Contains(output, tc.want):
+				t.Fatalf("startup refusal: %d %s", code, output)
+			}
+			if tc.arm == "same generation" {
+				if sleeps != 3 || reads != 8 || !clock.Equal(time.Unix(current, 0).Add(f.env.wait)) || f.ctl.loaded {
+					t.Fatalf("startup timeout drift: sleeps=%d clock reads=%d time=%v loaded=%v", sleeps, reads, clock, f.ctl.loaded)
+				}
+				h, err := runnerHeartbeat(creds, time.Second)
+				if err != nil || h.Owner.Generation != 1 || h.SeenAt != current {
+					t.Fatalf("old liveness was changed: %+v %v", h, err)
+				}
+			}
+			wantLoaded := tc.loaded || tc.stopFail || tc.arm == "fresh restart"
+			if f.ctl.loaded != wantLoaded || f.ctl.boots != tc.boots || f.ctl.stops != tc.stops || f.ctl.kicks != 0 ||
+				len(f.writes) != 0 || serviceFileHash(creds) != knownHash {
+				t.Fatalf("unexpected start effects: boots=%d stops=%d kicks=%d writes=%v", f.ctl.boots, f.ctl.stops, f.ctl.kicks, f.writes)
+			}
+		})
+	}
+}
diff --git a/cmd/ailang-worldd/service.go b/cmd/ailang-worldd/service.go
--- a/cmd/ailang-worldd/service.go
+++ b/cmd/ailang-worldd/service.go
@@ -82,13 +82,15 @@
 	observe     func(string)
 	probe       func(string, string) ([]byte, error)
 	runnerLogin func() (int, string, error)
+	runnerNow   func() time.Time
+	runnerSleep func(time.Duration)
 }
 
 func realServiceEnv() serviceEnv {
 	home, _ := os.UserHomeDir()
 	return serviceEnv{paths: homeServicePaths(home), controller: systemServiceController{}, terminal: realSessionEnv(),
 		goos: runtime.GOOS, domain: "gui/" + strconv.Itoa(os.Getuid()), label: worldServiceLabel,
-		wait: 30 * time.Second, client: &http.Client{Timeout: time.Second, CheckRedirect: func(*http.Request, []*http.Request) error {
+		wait: 30 * time.Second, runnerNow: time.Now, runnerSleep: time.Sleep, client: &http.Client{Timeout: time.Second, CheckRedirect: func(*http.Request, []*http.Request) error {
 			return errors.New("health redirects are refused")
 		}}}
 }
diff --git a/cmd/ailang-worldd/service_test.go b/cmd/ailang-worldd/service_test.go
--- a/cmd/ailang-worldd/service_test.go
+++ b/cmd/ailang-worldd/service_test.go
@@ -19,21 +19,27 @@
 )
 
 type fakeServiceController struct {
-	loaded       bool
-	runner       bool
-	boots, stops int
-	key          string
-	bootErr      error
-	kicks        int
-	binary       string
-	running      string
-	onKick       func() error
+	loaded         bool
+	runner         bool
+	boots, stops   int
+	key            string
+	bootErr        error
+	kicks          int
+	binary         string
+	running        string
+	onKick         func() error
+	onBoot         func() error
+	stopErr        error
+	runnerPrintErr error
 }
 
 func (f *fakeServiceController) Bootstrap(domain, plist string) error {
 	f.boots++
 	f.loaded = true
 	f.running = serviceFileHash(f.binary)
+	if f.bootErr == nil && f.onBoot != nil {
+		return f.onBoot()
+	}
 	return f.bootErr
 }
 func (f *fakeServiceController) Kickstart(domain, label string) error {
@@ -47,10 +53,15 @@
 }
 func (f *fakeServiceController) Bootout(domain, label string) error {
 	f.stops++
-	f.loaded = false
-	return nil
+	if f.stopErr == nil || errors.Is(f.stopErr, errServiceNotLoaded) {
+		f.loaded = false
+	}
+	return f.stopErr
 }
 func (f *fakeServiceController) Print(domain, label string) (string, error) {
+	if label == runnerLabel && f.runnerPrintErr != nil {
+		return "", f.runnerPrintErr
+	}
 	if !f.loaded || (label == runnerLabel && !f.runner) {
 		return "not loaded", errServiceNotLoaded
 	}
diff --git a/host/store/operator_runner.go b/host/store/operator_runner.go
--- a/host/store/operator_runner.go
+++ b/host/store/operator_runner.go
@@ -108,6 +108,20 @@
 	return id, err
 }
 
+// ReadOperatorHeartbeat bounds opening and reading by the same context.
+func ReadOperatorHeartbeat(ctx context.Context, path string) (OperatorHeartbeat, error) {
+	if err := requireDeadline(ctx); err != nil {
+		return OperatorHeartbeat{}, err
+	}
+	db, version, err := openOperatorReadOnly(ctx, path)
+	if err != nil {
+		return OperatorHeartbeat{}, err
+	}
+	defer db.Close()
+	c := Credentials{db: db, version: version, display: path}
+	return c.OperatorHeartbeat(ctx)
+}
+
 func openOperatorReadOnly(ctx context.Context, path string) (*sql.DB, int, error) {
 	canonical, _, err := resolveDSN(path)
 	if err != nil {
````

### Patch M3

````diff M3
diff --git a/cmd/ailang-worldd/operator_refusals_test.go b/cmd/ailang-worldd/operator_refusals_test.go
new file mode 100644
--- /dev/null
+++ b/cmd/ailang-worldd/operator_refusals_test.go
@@ -0,0 +1,82 @@
+package main
+
+import (
+	"bytes"
+	"os"
+	"path/filepath"
+	"reflect"
+	"testing"
+)
+
+func TestOperatorRunnerLockFileRefusals(t *testing.T) {
+	path := filepath.Join(t.TempDir(), "lock")
+	if err := os.WriteFile(path, []byte("keep"), 0600); err != nil {
+		t.Fatal(err)
+	}
+	before := serviceFileHash(path)
+	alias := path + "-alias"
+	if err := os.Symlink(path, alias); err != nil {
+		t.Fatal(err)
+	}
+	for _, tc := range []struct {
+		name, path string
+		mode       os.FileMode
+	}{
+		{"symlink", alias, 0600},
+		{"mode", path, 0644},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			if err := os.Chmod(path, tc.mode); err != nil {
+				t.Fatal(err)
+			}
+			lock, err := runnerKernelLock(tc.path)
+			if lock != nil {
+				lock.Close()
+			}
+			info, statErr := os.Stat(path)
+			target, linkErr := os.Readlink(alias)
+			if err == nil || lock != nil || statErr != nil || info.Mode().Perm() != tc.mode ||
+				serviceFileHash(path) != before || linkErr != nil || target != path {
+				t.Fatalf("unsafe lock did not refuse without mutation: %v %v", lock, err)
+			}
+		})
+	}
+}
+
+func TestOperatorRunnerSetupRefusals(t *testing.T) {
+	t.Run("linked config", func(t *testing.T) {
+		f := operatorService(t)
+		f.mustInstall(t)
+		p := f.env.runnerPaths()
+		alias := f.configFile + "-alias"
+		if err := os.Symlink(f.configFile, alias); err != nil {
+			t.Fatal(err)
+		}
+		before := []string{serviceFileHash(f.configFile), serviceFileHash(p.Record), serviceFileHash(p.Plist), serviceFileHash(f.config.Credentials)}
+		f.writes = nil
+		code, msg := f.command("install", "--config", alias)
+		f.refused(t, "invalid runner setup file", code, msg)
+		after := []string{serviceFileHash(f.configFile), serviceFileHash(p.Record), serviceFileHash(p.Plist), serviceFileHash(f.config.Credentials)}
+		if target, err := os.Readlink(alias); err != nil || target != f.configFile || !reflect.DeepEqual(after, before) {
+			t.Fatal("refusal modified configuration or credentials")
+		}
+	})
+	t.Run("reserved preset", func(t *testing.T) {
+		f := operatorService(t)
+		f.config.PresetNames = []string{"se-tools"}
+		code, msg := f.install(t)
+		f.refused(t, "invalid runner preset name", code, msg)
+	})
+	t.Run("global address", func(t *testing.T) {
+		home := t.TempDir()
+		t.Setenv("HOME", home)
+		t.Setenv("AILANG_REGISTRY_API_KEY", "")
+		var out, errOut bytes.Buffer
+		if code := run([]string{"--addr", "127.0.0.1:9", "operator-runner"}, &out, &errOut); code != exitUsage || out.Len()+errOut.Len() != 0 {
+			t.Fatalf("explicit address reached runner startup: %d %s", code, errOut.String())
+		}
+		if files, err := os.ReadDir(home); err != nil || len(files) != 0 {
+			t.Fatal("refused invocation changed its home")
+		}
+	})
+}
diff --git a/cmd/ailang-worldd/operator_service_test.go b/cmd/ailang-worldd/operator_service_test.go
--- a/cmd/ailang-worldd/operator_service_test.go
+++ b/cmd/ailang-worldd/operator_service_test.go
@@ -153,6 +153,16 @@
 			}
 			f.config.ToolDirs = []string{p}
 		}},
+		{"other writable", "writable by others", func(t *testing.T, f *operatorServiceFixture) {
+			p := filepath.Join(f.home, "other-writable")
+			if err := os.Mkdir(p, 0700); err != nil {
+				t.Fatal(err)
+			}
+			if err := os.Chmod(p, 0702); err != nil {
+				t.Fatal(err)
+			}
+			f.config.ToolDirs = []string{p}
+		}},
 		{"symlink tool", "wrong type", func(t *testing.T, f *operatorServiceFixture) {
 			p := filepath.Join(f.home, "linked")
 			if err := os.Symlink(f.env.paths.Binary, p); err != nil {
@@ -240,13 +250,16 @@
 		{"credentials", "credentials store mismatch", func(f *operatorServiceFixture) { f.health.CredentialsDB = f.config.World }},
 		{"build", "restart pending", func(f *operatorServiceFixture) { f.health.Build.ExecutableSHA256 = "sha256:other" }},
 		{"loopback", "numeric loopback", func(f *operatorServiceFixture) { f.config.Address = "192.0.2.1:7644" }},
+		{"named loopback", "numeric loopback", func(f *operatorServiceFixture) { f.config.Address = "localhost:7644" }},
+		{"distinct", "distinct canonical files", func(f *operatorServiceFixture) { f.config.Credentials = f.config.World }},
+		{"health status", "daemon unavailable", func(f *operatorServiceFixture) { f.health.Status = "starting" }},
 	} {
 		t.Run(tc.name, func(t *testing.T) {
 			f := operatorService(t)
 			f.mustInstall(t)
 			f.writes = nil
 			tc.change(f)
-			if tc.name == "loopback" {
+			if tc.name == "loopback" || tc.name == "named loopback" || tc.name == "distinct" {
 				c, msg := f.install(t)
 				f.refused(t, tc.want, c, msg)
 				return
diff --git a/host/store/credentials_test.go b/host/store/credentials_test.go
--- a/host/store/credentials_test.go
+++ b/host/store/credentials_test.go
@@ -308,6 +308,9 @@
 	_, errs["PendingVerbRequests"] = c.PendingVerbRequests(ctx, 1000, 0)
 	errs["ConsumeVerbRequest"] = c.ConsumeVerbRequest(ctx, "x", "y", ApproverUse{}, "", nil, "", 1000)
 	_, errs["OperatorAction"] = c.OperatorAction(ctx, "x", 1000)
+	_, errs["OperatorHeartbeat"] = c.OperatorHeartbeat(ctx)
+	_, errs["ReadOperatorHeartbeat"] = ReadOperatorHeartbeat(ctx, "")
+	_, errs["ReadOperatorEpisodeOwner"] = ReadOperatorEpisodeOwner(ctx, c.display, "ep")
 	for name, err := range errs {
 		if !errors.Is(err, ErrNoDeadline) {
 			t.Errorf("%s(Background) = %v, want ErrNoDeadline", name, err)
````

