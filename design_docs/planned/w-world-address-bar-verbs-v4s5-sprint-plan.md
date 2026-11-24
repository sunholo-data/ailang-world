# Sprint plan — w-world-address-bar-verbs-v4, sprint S5 (iteration 431, row 243)

> **Scope.** Take the design §8 fallback split at the read-seam boundary. M1 closes terminal mint exclusion and documents the sidecars. M2 makes terminal preset start and a future operator adapter share a typed, bounded read seam. M3 closes exact-artifact confirm replay. S5b carries the atomic mint and recovery together. Cumulative product changes are **163 → 730 → 743**.
>
> **Authority.** Reviewed design r1 §§2.2, 3.1–3.3, 4.2, 5.3, 5.5, 6–8. Quorum r2 PROCEED 3/3 with all seats present. Attended D-WORLD-119=A and D-WORLD-120=A (`82f05004`). S1 #606, S2 #609, S2b #611, S3 #643, S4 #651 (`ee8ef1c0`) and S4b #654 (`fa3ad190`, judge PASS 85 with zero blocking findings) are landed. Base `cd06d9500ec2d34e975d7566f8bbdb14b1529527`.

## 0. Binding constraints

- The controller commits the two plan records first. Apply M1–M3 in order using the Appendix P unified diffs and `/usr/bin/patch -p1 -E`. The controller owns Git writes.
- Use disposable stores from the real migration, signatures from the repository signer, fake HTTP/supervisor/tty and fixture Git roots. All writes stay within `t.TempDir()` or mission scratch. Existing operator paths and launchctl are observation-only.
- Preserve the public terminal fence, CLI read/prompt/effect order, approval forms and served route table. The runner remains inert for operator effects. Appendix C governs activation.

## 1. Pre-flight

### 1.1 Baseline at cd06d950

Before any stage edits, `git archive cd06d950 | tar -x -C "$S/base"` exported the pinned tree. The first run used the prior §6 script with the SHA changed. `GATES_GIT_REF="$W" GATES_SANDBOX=1 GATES_SKIP_VERIFY_GO=1 /bin/bash "$S/gates-first-base.sh" "$S/base" "$S/gates-base" base cd06d950` produced the following. The archive lacks the Git metadata required for `verify_go.sh`'s tracked-binary inventory. The controller runs that full leg in its Git checkout.

````text
LEG observation-before rc=77 0s UNINFORMATIVE UNDER SANDBOX
LEG paths rc=0 0s
TOUCHED_PACKAGES: 
LEG touched rc=0 0s
LEG status rc=0 0s
LEG gofmt rc=0 0s
LEG build rc=0 4s
LEG vet rc=0 3s
LEG fence rc=0 29s
LEG verify_ail rc=0 34s
LEG pkg-cmd-ailang-worldd rc=1 2s UNINFORMATIVE UNDER SANDBOX
LEG race-cmd-ailang-worldd rc=1 21s UNINFORMATIVE UNDER SANDBOX
LEG pkg-host-store rc=0 15s
LEG race-host-store rc=0 36s
LEG inv-host-approval rc=0 1s
LEG inv-host-approval-approvaltest rc=0 0s
LEG inv-host-authority rc=0 1s
LEG inv-host-boundary rc=0 1s
LEG inv-host-daemon rc=0 2s
LEG inv-host-verbs rc=0 0s
LEG inv-cmd-ailang-worldd rc=0 1s
LEG inv-host-store rc=0 1s
LEG routes rc=0 0s
LEG scope rc=0 0s
LEG callers rc=0 0s
LEG verify_go SKIPPED 0s (GATES_SKIP_VERIFY_GO=1)
LEG observation-after rc=77 0s UNINFORMATIVE UNDER SANDBOX
GATES_DONE 151s FAILS: none UNINFORMATIVE: observation-before pkg-cmd-ailang-worldd race-cmd-ailang-worldd observation-after
````

The final-format script also ran on that pristine base with all three affected packages. Whole CLI plain and race hit the inherited `TestApproveDecideWriteMessageAndSignature` loopback bind panic. These are sandbox-uninformative, with a controller rerun required. `AILANG_BIN` reports released v0.41.0 (`24ee108`). Go is `go1.26.6 darwin/arm64`, matching the CI version. Linux CI runs on ubuntu-24.04.

````text
LEG observation-before rc=77 0s UNINFORMATIVE UNDER SANDBOX
LEG paths rc=0 0s
TOUCHED_PACKAGES: 
LEG touched rc=0 0s
LEG status rc=0 1s
LEG gofmt rc=0 0s
LEG build rc=0 4s
LEG vet rc=0 1s
LEG fence rc=0 2s
LEG verify_ail rc=0 177s
LEG pkg-cmd-ailang-worldd rc=1 10s UNINFORMATIVE UNDER SANDBOX
LEG race-cmd-ailang-worldd rc=1 57s UNINFORMATIVE UNDER SANDBOX
LEG pkg-host-authority rc=0 4s
LEG race-host-authority rc=0 15s
LEG pkg-host-store rc=0 12s
LEG race-host-store rc=0 52s
LEG inv-host-approval rc=0 1s
LEG inv-host-approval-approvaltest rc=0 1s
LEG inv-host-authority rc=0 1s
LEG inv-host-boundary rc=0 2s
LEG inv-host-daemon rc=0 1s
LEG inv-host-verbs rc=0 1s
LEG inv-host-broker rc=0 1s
LEG inv-cmd-ailang-worldd rc=0 3s
LEG inv-quickstart-section9 rc=1 3s UNINFORMATIVE UNDER SANDBOX
LEG inv-host-store rc=0 1s
LEG routes rc=0 0s
LEG scope rc=0 1s
LEG callers rc=0 0s
LEG verify_go SKIPPED 0s (GATES_SKIP_VERIFY_GO=1)
LEG observation-after rc=77 0s UNINFORMATIVE UNDER SANDBOX
GATES_DONE 351s FAILS: none UNINFORMATIVE: observation-before pkg-cmd-ailang-worldd race-cmd-ailang-worldd inv-quickstart-section9 observation-after
````

### 1.2 Current inputs and coupling inventory

The commands below ran in the pinned worktree. `F=cmd/ailang-worldd`, `S=$HOME/.ailang/state/world-iter431/plan-scratch`, and `W` is the assigned worktree. Outputs are narrowed to the lines used by this plan.

| ID | Command | Observed output and consequence |
|---|---|---|
| C1 | `git rev-parse HEAD; git diff --stat fa3ad190 cd06d950 -- cmd host app docs; git diff --name-only fa3ad190 cd06d950` | `cd06d9500ec2d34e975d7566f8bbdb14b1529527`. Product stat empty. Eight paths under `design_docs`, including the S4b judge and readings. |
| C2 | `rg -n 'provisionEpisode\(' "$F"; rg -n 'terminalEpisodeLock\(' "$F" -g '*.go' -g '!**/*_test.go'` | Provision definition `session_new.go:405`. Five callers at `session_new.go:649`, `approve.go:398`, `episode.go:377`, `episode_preset.go:312`, `episode_resume.go:215`. The old effect lock call sites are `episode.go:364`, `episode_preset.go:282`, `episode_resume.go:175`, `episode_launch.go:456`, `episode_finish.go:128`. Acquires belong to the canonical credentials/episode pair via `operator_episode_lock.go:14–45`. |
| C3 | `nl -ba "$F/episode_preset.go" \| sed -n '21,325p'; rg -n 'loadLaunchPreset\|resolvePresetDaemon\|type resolvedPreset' "$F/preset.go"; rg -n 'runPresetEpisodeStart\(' "$F"` | Preset root/name/authority resolution `preset.go:298,329–502`; one terminal caller `episode.go:272`. Defaults/explicit overrides, tool and launch preflight `episode_preset.go:66–144`; Git branch probe `:148–175`; path and task probes `:177–221`; `prepareMint`, profile and v3 binding floor `:223–244`; fixed start, prompt and served pair `:245–281`; effect separation at `:282`. |
| C4 | `rg -n 'prepareMint\|openCredentialsForSession\|worktree\.(Add\|AddExisting\|Remove)\|MintForRequest\|MintStarted\|writeTokenFile\|writeMCPConfig\|afterToken' "$F" -g '*.go' -g '!**/*_test.go'; rg -n 'func (Add\|Output)\|func .*MintSessionWithStart\|func .*insertSessionRow\|func .*EpisodeStart' host/{worktree,store}` | `session_new.go:237` is read validation, `:487` opens the writer, `:499/501` selects Git add, `:519/523` select A3/start mint, `:538/541` install token/extra files, `:566` legacy rollback, `:663` token no-replace writer. `episode_preset.go:287–310` task copy and MCP hook. `episode.go:552` MCP no-replace writer. `approve.go:365–396` re-derives independent A3 bytes and constructs `RequestMint`. `worktree.go:63,83` Add/Output. `host/authority/mint.go:50,112` MintStarted/MintForRequest. `host/store/credentials.go:469–520` tentative session/start transaction, `episode_starts.go:70` reader. The session insert helper is at credentials.go:526. S5b must couple the receipt here and can share bounded Git and file mechanisms. The legacy rollback contract is separate from uncertain operator worktree provenance. |
| C5 | `rg -n 'runnerStatus\|runnerCancel\|case "cancel"\|CancelUnstartedOperatorAction\|uncertain\|reconcile at the terminal' "$F"/operator*.go host/store/operator_{runner.go,schema.sql}; AILANG_BIN=$HOME/.pinned-ailang/ailang go test -v -count=1 ./cmd/ailang-worldd -run '^TestS5ResolverProbe$'` in `$S/st-footprint` | Existing status dispatch `operator_service.go:285–289`; typed cancel `:275,291–329,423–424` → `operator_runner.go:171` → `host/store/operator_runner.go:74–83` with requested/confirmed, owner NULL, valid expiry and only operation aliases. `operator_runner.go:130–158` renders heartbeat and supervisor readiness. Exact status command rc=0, attended cancel rc=1 conflict, state still uncertain and owner unchanged. Fixture and log are `$S/st-footprint/cmd/ailang-worldd/s5_resolver_probe_test.go` and `$S/resolver-probe.log`. The real signer, migration and claim create the row, and fixture SQL sets the otherwise-unreachable phase. Hash, full row and supervisor counts match. |
| C6 | `git grep -n -E '\.(ClaimOperatorAction\|ReserveOperatorStep\|ObserveOperatorStep\|RecheckClaimedOperatorAction\|RecheckVerb)\(' -- cmd host \| grep -v '_test\.go:'` | The only match is positive control `cmd/ailang-worldd/episode_finish_surface.go:239` (`authority.RecheckVerb`). `operator_runner.go:77–100` publishes preset names/heartbeat. `host/store/operator_cas.go:116` implements claim. Zero live claim/reserve/observe/recheck-claimed call sites on base and final stages. |
| C7 | `rg -n 'func Test.*(DependencyAllowlist\|ProductionGoSurface\|ProductionContextRoots\|WriterSites\|Quickstart\|VerbTable\|EveryVerb\|OperatorFlagParity\|PluginCommands)' cmd host -g '*test.go'; git grep -n 'LOCK_NB' -- cmd/ailang-worldd host/store` | Source positions and why expectations stay stable are below. Six active flock calls and one writer-lock comment under production Go search. S5 adds one cooperating caller to the existing episode lock helper and leaves low-level open/flock flags as landed. |
| C8 | `(cd "$S/st-footprint" && /usr/bin/patch -p1 -E < "$HOME/.ailang/state/world-iter430/plan-scratch/S3b.patch"); git diff --no-index -M --numstat "$S/base/host" "$S/st-footprint/host"` | Both files applied. `operator_cas_test.go` +60 and `operator_steps.go` +72, total 132. `AILANG_BIN=$HOME/.pinned-ailang/ailang go test -count=1 -v ./host/store -run '^TestOperatorStepRestart$'` passes on the bank, covering only the fixture window. The S4c bank applied as +3 store and +10 authority tests, total 13. |
| C9 | `git diff --no-index -M --numstat "$S/base" "$S/st-M3"; rg -n '^.{0,4}S5\|Shared pinned-start' design_docs/planned/w-world-address-bar-verbs-v4.md` | Compiled product cut 647 added, 96 deleted. Design §8 allocates 510 raw to atomic mint and recovery. 510 × 1.35 rounds to 689. The measured 743 plus that calibrated 689 projects to 1,432, before the resolver, fd/process tests, current-owner recovery and drain. The measured 132-line sketch remains separately gated. |
| C10 | `rg -n 'provisionEpisode\(\|terminalEpisodeLock\(\|resolvePresetStart\(\|readPresetStartTask\(\|bytes.Equal\(row.Artifact' "$S/st-M3"/cmd/ailang-worldd/*.go "$S/st-M3/host/store/operator_cas.go"` | Final `session_new.go:405,424,660`; `approve.go:398`; `episode.go:364,377`; `episode_resume.go:175,215`; `episode_preset.go:79,288,326,328,334,365`. The Git reads take the caller context. The existing health and binding helpers retain their individual finite budgets. Terminal starts still acquire before writes and transfer the held description. Artifact replay equality is at `operator_cas.go:59`. |
| C11 | `rg -n 'fetchDaemonHealth\|func getDaemonHealth\|func checkDaemon\|pinnedToolRef' "$F/episode.go"; rg -n 'func launchMode\|func launchBinary\|func validateLaunchDir\|func launchLockCheck\|func launchLauncher' "$F"/episode_launch*.go` | `episode.go:429–449` bounded GET transport, `:455,465–514` tool/build/profile/credentials F6 predicates. `episode_launch.go:44,98,151,168` mode/cap, path checks, resolved executable and PID lock read. `episode_launch_interactive.go:43` confined wrapper validation with bounded `--version` probe. S5's typed read records those inspected values. S5b adds signed equality, current generation checks and private receipt provenance. |

| Inventory | Base anchor | Disposition |
|---|---|---|
| Commands, flags and help | `F/help_test.go:21` `TestEveryVerbHasHelp`, `verbs_parity_test.go:66` `TestVerbTableMatchesTheCLISwitch`, `operator_parity_test.go:65` `TestOperatorFlagParity`, `episode_launch_test.go:449` `TestLaunchAndStartHelp`, `plugin_help_test.go:52` | Existing switch, flag set, help and plugin command records stay byte-identical. M2 adapts already parsed flags. |
| QUICKSTART | `quickstart_service_test.go:15`; `quickstart_se_test.go:56`; `quickstart_verbs_test.go:15`; `quickstart_finish_surface_test.go:17` | Four prose lines at `docs/QUICKSTART.md:77–80`. Existing command inventory expectations stand. §3 executes the existing fake-supervisor example command text verbatim. No golden regenerated. |
| CLI/daemon writers | `F/cli_credential_writers_test.go:18`; `host/daemon/surface_approve_post_guard_test.go:132` | Same writer functions and call-site counts. The preset read-only open moves from the terminal function into its resolver. M3 adds a comparison inside the existing transaction. |
| Flock | `operator_runner.go:35`, `operator_service.go:251`, `service.go:158,662`, `host/store/writer_lock_unix.go:41,86`; `TestOperatorEpisodeLockExclusive`, `TestOperatorTerminalExclusion/nonwaiting`, `TestWriterLockCrossProcessContention` | Same exclusive and nonwaiting primitives, checked with independent opens. M1 passes the already-held description into shared provisioning from three callers and uses the same fallback for the other two. |
| Context roots and Go surface | `host/store/context_roots_test.go:129,234` `TestProductionContextRoots`, `TestProductionGoSurface` | Existing `runPresetEpisodeStart` root stays in the same owner. `resolvePresetStart` accepts the caller context. Production files and package allowlist stay the same. |
| Dependencies and imports | `host/daemon/daemon_test.go:899`, `host/boundary/allowlist_world_test.go:758,1323`, `host/broker/allowlist_test.go:89`, `host/approval/imports_test.go:27` | The final inventories pass with 54 named context roots, 324 production Go files, six explicit reproducers, both cmd groups and 25 CLI read-only opens. Dependency and writer inventories retain their base expectations. Test-only stdlib and existing module imports remain under the existing test package boundaries. |
| Serving and effect dispatch | `host/daemon/surface_verb_test.go:133`, `TestFrozenV1RouteTableMatchesMux`, `TestSurfaceVerbRouteRegistered`, `TestOperatorNeverDispatch` | §6 compares daemon mux, CSP and app bytes, runs the inventories and positively controls the production caller search. |

### 1.3 Decisions

- Use a single missing-holder fallback in `provisionEpisode`. It runs after the existing terminal-open/verified-signature fence and before prompt, writer, worktree and mint. A valid signed approve can reach that fallback without a tty. Awaiting a key stays before the effect lock. An attended terminal that declines may leave the private sidecar. A session new with no tty refuses before creating it.
- The shared preset resolver accepts fixed typed input, the terminal's explicit overrides and a task reader. The terminal callback reads its own original file once. The future operator caller supplies persisted inline bytes and the rig-derived copy path. A resolved mint request carries the legacy start fields and prompt. Its presence provides no key or capability.
- The operator's read-stage fixtures call the actual resolver with inline bytes, then check the terminal separately using temporary files and injected clocks. The runner's catalogue, launch adapter and status paths receive no new callers or effects.
- S5b must build the attended resolver before any production worktree/mint uncertainty is possible. Keep proof construction, transaction identity, file recovery and their refusal tests together. Activation requires the same paired-process lock and drain even if the transaction can be tested in an inert store harness first.

### 1.4 Premises found false or incomplete at HEAD

| ID | Evidence and consequence |
|---|---|
| F1 | C2 confirms both missing terminal mint callers. A synthetic operator owner persists through releasing its descriptor and the injected clock passing expiry. M1 still refuses session new, and the signed approve arm refuses without consuming request or approver use. |
| F2 | C5 corrects the judge's statement about cancel. A typed `service runner cancel` exists for unstarted work. Status and cancel provide no attended resolution for an owned uncertain row. Existing status may even report startup readiness while that row retains its owner. |
| F3 | `grep -Fc operator-episode docs/QUICKSTART.md` returned 0 before M1. M1 documents the private, stable sidecar alongside the canonical credentials DB. The same fixed-string count is 1 on M1 through M3. |
| F4 | C6 and C8. A schema table and a claim CAS exist, while no paired-process effect consumer or drain is present. The 132-line bank adds reserve/observe and one restart fixture. It omits real fd coordination and current-owner recovery. |
| F5 | `nl -ba host/store/operator_cas.go \| sed -n '54,63p'` reads digest, then `if row.ApprovedBy != "" { return row, nil }`. M3 requires byte equality with the stored artifact. Identical SSHSIG and passkey replays still return the prior row. |
| F6 | The read seam exposes existing terminal preflight. Its result carries no durable reservation. It allows a terminal explicit free existing branch and allows a legacy typed task path. The operator's stricter absence checks and signed-input comparison remain S5b gates. |

### 1.5 Size and CI race cost

After each controller commit measure `git diff -M --numstat cd06d950 HEAD -- cmd host docs app`. Sum insertions and deletions, with renames enabled. Scratch used `git diff --no-index -M --numstat "$S/base" "$S/st-$M"`, then verified every product path is inside the same filter. The unfiltered product sum is also 743. Including the two plan records, the unfiltered sum is **2,599**.

| Milestone | Added | Deleted | Cumulative | Increment | Snap STATUS / DELETED / manifest |
|---|---:|---:|---:|---:|---|
| M1 | 157 | 6 | 163 | 163 | 7/0/7 |
| M2 | 634 | 96 | 730 | 567 | 8/0/8 |
| M3 | 647 | 96 | 743 | 13 | 10/0/10 |

Design §8 budgets 390 raw lines for the read seam and 510 for mint and recovery, totaling 900 raw and 1,215 calibrated. The 743-line compiled cut includes the read seam, carried lock/doc work and replay. It leaves 657 to the cap. The separately calibrated design allocation for mint/recovery is 689. That projection crosses the cap even before the carried resolver and paired-process gates. The compiled read seam and gates are a safe intermediate commit under §8. No tests land in `host/daemon` or `host/verifygate`. The row 267 charter entry at line 6362 records four slow-runner daemon race timeouts at roughly 480.1 seconds. Sprint cost comes from the CLI, authority and store packages.

## 2. How the executor applies a sub-milestone

Set `W=$PWD`, `T=$HOME/.ailang/state/world-iter431/exec-patches` and `PLAN=$W/design_docs/planned/w-world-address-bar-verbs-v4s5-sprint-plan.md` in the assigned worktree. The controller records the two plan inputs before the snapshot commits. Set `M` to each of M1, M2 and M3 in order.

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

Stop at any reject, unexpected file, format/test failure or size drift. §7 gives exact acceptance and snapshot commands.

## 3. Milestones

### M1 — terminal provision exclusion and sidecar documentation (163 cumulative)

Use the same lock in all five `provisionEpisode` call paths. `TestOperatorProvisionEpisodeExclusion` drives the real session-new and signed approve entry points using an actual signed row and canonical lock. Each refuses with the exact owner and unchanged effects. The no-tty test checks absence of the sidecar, and independent session and approve fixtures exercise successful mint and release. Continue to use `TestOperatorTerminalExclusion` for the previously guarded verbs and `TestStartLaunchHeadlessMintsAndLaunchesOnce` for the nested descriptor transfer. `TestSessionVerbsOnAPreSplitStore` covers the existing migration path.

### M2 — shared read seam and semantic parity (730 cumulative)

`resolvePresetStart` carries the resolved preset, launch mode/cap/executable, selected store/health, grants/checks/TTL, literal task data and legacy start record plus prompt. The terminal uses that result for provisioning. The write block retains exclusive token/task/MCP behavior. `TestPresetStartReadParity` independently checks prompt inputs, fixture Git base/worktree, parsed credential and byte-identical task after the source file changes at the terminal fence. `TestPresetStartReadHeadless` tests literal inline data and resolved headless metadata with a no-spawn marker. `TestPresetStartReadInteractive` tests the bounded fake host preflight, explicit model and MCP suppression, then drives the actual CLI decline. `TestPresetStartReadRefusals` covers local collisions and preflight faults with complete fixture readbacks.

### M3 — exact confirm replay (743 cumulative)

Land the 13-line S4c bank. `TestOperatorConfirmReplay` exercises real SSHSIG and WebAuthn artifacts against the durable row after reopening. Exact replay returns it unchanged. Different bytes return `ErrOperatorConflict` with the row and signer counters intact. S7 remains responsible for the status link and route.

````bash accept.sh
#!/bin/bash
set -euo pipefail
W=$1
M=$2
test "$3" = fast
S=$HOME/.ailang/state/world-iter431/plan-scratch
cd "$W"
export AILANG_BIN="${AILANG_BIN:-$HOME/.pinned-ailang/ailang}"
export AILANG_STATE_DIR="$S/accept-ailang-state" AILANG_CACHE_DIR="$S/accept-ailang-cache"
export WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR WORLD_UPDATE_OPERATOR_ACTION_GOLDEN WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN
mkdir -p "$AILANG_STATE_DIR"
test -x "$AILANG_BIN"
cli='TestOperatorProvisionEpisodeExclusion TestOperatorTerminalExclusion TestOperatorEpisodeLockExclusive TestSessionNewProvisionsTheWorktreeAndTheSession TestSessionNewWithoutATerminalProvisionsNothing TestSessionNewDeclinedProvisionsNothing TestSessionNewMintFailureRollsBackTheWorktree TestSessionVerbsOnAPreSplitStore TestApproveBySignatureMintsWithNoTerminal TestApproveMintsExactlyTheRequest TestStartLaunchHeadlessMintsAndLaunchesOnce TestQuickstartServiceSectionMatchesTheCLIHelp'
case "$M" in
  M1) ;;
  M2|M3)
    cli="$cli TestPresetStartReadParity TestPresetStartReadHeadless TestPresetStartReadInteractive TestPresetStartReadRefusals TestPresetStartLocalBaseAndPinnedHistory TestPresetDeclineAndNoTerminal TestPresetDerivedCollisionsAndExistingBranch TestPresetRefusalsBeforeDaemon TestStartLaunchValidatesBeforeFence TestStartLaunchPreFenceTaskBinaryAndDirectory TestOperatorFlagParity"
    ;;
  *) exit 2 ;;
esac
run_tests() {
  local pkg=$1 names=$2 log pattern test
  log=$(mktemp "$S/accept-$M.XXXXXX")
  pattern=$(printf '%s' "$names" | tr ' ' '|')
  go test -v -count=1 -timeout 8m "./$pkg" -run "^($pattern)$" | tee "$log"
  for test in $names; do
    grep -Fq -- "--- PASS: $test (" "$log"
  done
  if grep -Fq -- '--- SKIP:' "$log"; then echo "Unexpected skip in $log"; exit 1; fi
}
run_tests cmd/ailang-worldd "$cli"
# Existing QUICKSTART example, executed with its flags verbatim.
AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./cmd/ailang-worldd -run '^TestOperatorServiceWriteScope$' -count=1 -v
if [ "$M" = M3 ]; then
  run_tests host/authority 'TestOperatorConfirmReplay TestOperatorConfirmCrash TestOperatorRunnerRechecks'
  run_tests host/store 'TestOperatorSchema10 TestOperatorClaimBeforeExpiry TestOperatorClaimRace TestOperatorCancelConfirmed'
fi
echo "ACCEPT_OK $M"
````

### 3.1 Branch traces

| Input or failure | Observable behavior and readbacks |
|---|---|
| Session new lacks controlling tty | Existing refusal, no new sidecar, Git refs and worktrees, credentials and disposable session/workspace contents match. |
| Session new has a tty, owned row after fd close/lease expiry | Exact action ID with conflict and no approval prompt. The token path remains absent, complete Git/store/workspace and files match. An uncontended episode provisions once and releases its lock. |
| Valid signed approve during contention | Exact action ID with conflict, no tty. Request, approver use, SQLite bytes, worktree registration, refs and session paths match. The uncontended signed fixture consumes once and permits a subsequent independent descriptor. |
| Existing lock holders in plain start/preset/resume | Existing single acquire remains in control. Nested start/launch uses the same description. The prior F1/health tests and same/different key tests still pass. |
| Local preset, base, launch and task read | Git queries, fixture health and task read return the typed plan. No new sidecar on direct read, zero child marker, unchanged fixture store/refs/files. The attended path uses the read prompt and bytes, pins the base despite a different checkout HEAD and writes only fixture artifacts. |
| Branch exists or explicitly chosen branch is checked out | Class-specific refusal. Worktree registration and refs match. Independent explicit free existing branch remains valid for terminal start. |
| Existing worktree, token, task or MCP | Specific refusal for each path with Git refs, workspace contents, session bytes and symlink target unchanged. Worktree input is a real linked fixture. MCP uses a dangling symlink to pin `Lstat`. |
| Path loop, nonprivate sessions, task directory or missing task | Portable loop class using `syscall.ELOOP`, private-mode or task class, no tty/Git/store/file write. No symlink permission assumption. |
| Nonnegative TTL validation, grammar, ref format and agent | Invalid input refuses before querying health. Successful semantic fixture confirms explicit TTL, branch, agent and terminal-only model behavior. |
| Oversized wall cap, missing task or executable, invalid launch dir, stale PID file | Preflight-specific error, zero child/mint and no configuration output. Interactive validates the confined host manifest under fake binaries. |
| Edited profile across two reads, missing health or mismatched world, build, tool, credentials or served profile | Distinct refusal and health-call count, zero prompt or effects. Shared terminal preflight only. Signed stale action comparison and catalogue generation recheck are S5b. |
| Exact replay vs changed artifact after confirmation | Matching row, approver use and passkey count for exact bytes. Different bytes get conflict and the persisted row still equals the original. |

### 3.2 Design §5.3 row assignment

| Design test | Disposition |
|---|---|
| `TestOperatorStartOnce` | S5b earns one minted start/worktree/credential and stable receipt. S6 earns original launch-before-receipt and confinement assertions. M2's terminal-only positive `TestStartLaunchHeadlessMintsAndLaunchesOnce` is a legacy control. |
| `TestOperatorStartChanged` | M2 pins live preset resolution and independent health/profile guards. S5b adds the reviewed-to-current digest/generation/identity comparison and its named test. |
| `TestOperatorStartCollisions` | M2 takes the shared local branch/worktree/token/task/MCP refusals in `TestPresetStartReadRefusals`. S5b adds reserved/current episode checks, positive evidence of Git ref absence and canonical no-adoption at the effect boundary. |
| `TestOperatorWorktreeWindow` | S5b. Require worktree reservation and provenance-backed recovery with an attended escape for ambiguous evidence. |
| `TestOperatorMintBinding`, `TestOperatorMintRollback` | S5b. Bind the typed proof to the credential/start/receipt in one transaction. Pin each fresh predicate and inject rollback before commit. |
| `TestOperatorMintReceiptWindow`, `TestOperatorFilesReconcile` | S5b. Reopen after atomic commit, adopt only exact private material or revoke the orphan, and preserve unmatched paths. |
| `TestOperatorLaunchSession`, `TestOperatorSpawnRefused`, `TestOperatorLaunchCrashWindow`, `TestOperatorLaunchExit` | S6 process/receipt and launch boundary gates. |
| `TestOperatorStartBrowser` | S6 takes the process restart/launch half against fakes. S8 owns the live browser close/reconnect end-to-end acceptance with the S7 route. |
| `TestOperatorHerdrUnavailable`, `TestOperatorHerdrHandoff`, `TestOperatorHerdrWindow` | S12 under D-WORLD-119=A and E4. |

## 4. Mutation table

**79/79 named behavioral kills.** The 77 CLI whole-package variants are **UNMEASURED IN SANDBOX**. The base and final whole CLI runs hit the same loopback bind panic. The two replay variants are killed by whole `host/authority`. The two replay variants survive the isolated `host/store` package, which lacks the cross-package real-signer replay test. No inventory-only kill is counted.

For each CLI row, apply only its recorded edit on a disposable final tree, run `AILANG_BIN=$HOME/.pinned-ailang/ailang go test -count=1 -timeout 8m ./cmd/ailang-worldd` outside the sandbox using the §6 exports, and restore the original bytes before the next. `$S/mutants/ledger.json` stores exact edits and argument vectors. `$S/mutants/runs/<id>/named.log` records the failing named test. The initial episode-grammar test let the later `prepareMint` duplicate check hide the skipped early check. The strengthened case enters launch preflight with a missing fixture executable and still requires the earlier grammar refusal. The final grammar variant fails that assertion.

`cmd` in the table means `cmd/ailang-worldd`. The two inherited flock flag variants and four health guards are included as regression pins for the shared read/lock callers.

| ID | M3 diff anchor or called guard | Named test and assertion | Named | Whole package |
|---|---|---|---|---|
| L-SNEW | `cmd/session_new.go:423` | `TestOperatorProvisionEpisodeExclusion/session_new` — durable owner refusal without a prompt | KILLED | UNMEASURED IN SANDBOX |
| L-APPROVE | `cmd/approve.go:398` | `TestOperatorProvisionEpisodeExclusion/approve` — signed approve cannot skip exclusion | KILLED | UNMEASURED IN SANDBOX |
| L-DUP | `cmd/session_new.go:423` | `TestStartLaunchHeadlessMintsAndLaunchesOnce` — held description avoids a second acquire | KILLED | UNMEASURED IN SANDBOX |
| L-KEY-CREDS | `cmd/session_new.go:424` | `TestOperatorProvisionEpisodeExclusion/session_new` — owner comes from served credentials db | KILLED | UNMEASURED IN SANDBOX |
| L-KEY-EP | `cmd/session_new.go:424` | `TestOperatorProvisionEpisodeExclusion/session_new` — owner lookup uses exact episode | KILLED | UNMEASURED IN SANDBOX |
| L-OPEN-ERR | `cmd/session_new.go:424` | `TestOperatorProvisionEpisodeExclusion/approve` — lock refusal precedes consumed A3 proof | KILLED | UNMEASURED IN SANDBOX |
| L-REFUSAL-CLASS | `cmd/session_new.go:426` | `TestOperatorProvisionEpisodeExclusion/session_new` — conflict uses the existing usage exit | KILLED | UNMEASURED IN SANDBOX |
| L-REFUSAL-ID | `cmd/session_new.go:426` | `TestOperatorProvisionEpisodeExclusion/session_new` — exact active action is printed | KILLED | UNMEASURED IN SANDBOX |
| L-RELEASE | `cmd/session_new.go:429` | `TestOperatorProvisionEpisodeExclusion/approve_releases` — normal success releases the new descriptor | KILLED | UNMEASURED IN SANDBOX |
| L-OPEN-CREATE | `cmd/session_new.go:424` | `TestOperatorProvisionEpisodeExclusion/session_new` — first unmanaged terminal creates the same stable lock | KILLED | UNMEASURED IN SANDBOX |
| L-DURABLE-OWNER | `cmd/session_new.go:424` | `TestOperatorProvisionEpisodeExclusion/session_new` — a free kernel lock cannot clear a durable owner | KILLED | UNMEASURED IN SANDBOX |
| L-TRANSFER-PLAIN | `cmd/episode.go:377` | `TestEpisodeStartProvisionsAndPrintsTheLaunchWithTheTokenInTwoFilesOnly` — plain start retains its held descriptor | KILLED | UNMEASURED IN SANDBOX |
| L-TRANSFER-PRESET | `cmd/episode_preset.go:365` | `TestPresetStartReadParity` — preset retains its held descriptor | KILLED | UNMEASURED IN SANDBOX |
| L-TRANSFER-RESUME | `cmd/episode_resume.go:215` | `TestEpisodeResumeChainUsesCurrentCredential` — resume retains its held descriptor | KILLED | UNMEASURED IN SANDBOX |
| L-SIDECAR-LOCK-EX | `cmd/operator_runner.go:35` | `TestOperatorEpisodeLockExclusive` — shared holder prevents the exclusive acquire | KILLED | UNMEASURED IN SANDBOX |
| L-SIDECAR-LOCK-NB | `cmd/operator_runner.go:35` | `TestOperatorTerminalExclusion/nonwaiting` — independent contender returns immediately | KILLED | UNMEASURED IN SANDBOX |
| L-F1-ORDER | `cmd/session_new.go:409` | `TestOperatorProvisionEpisodeExclusion/no_terminal` — no tty means neither sidecar nor durable writes | KILLED | UNMEASURED IN SANDBOX |
| R-READ-MODE | `cmd/episode_preset.go:86` | `TestPresetStartReadHeadless` — headless mode resolves from requested launch | KILLED | UNMEASURED IN SANDBOX |
| R-WALL-ERR | `cmd/episode_preset.go:87` | `TestPresetStartReadRefusals/invalid_wall` — huge wall cap fails closed | KILLED | UNMEASURED IN SANDBOX |
| R-MISSING-TASK | `cmd/episode_preset.go:91` | `TestStartLaunchPreFenceTaskBinaryAndDirectory/task` — headless needs a task before checking tools | KILLED | UNMEASURED IN SANDBOX |
| R-TASK-MODE | `cmd/episode_preset.go:91` | `TestPresetStartReadInteractive` — interactive availability without a task | KILLED | UNMEASURED IN SANDBOX |
| R-EP-COUNT | `cmd/episode_preset.go:95` | `TestPresetStartReadRefusals/two_episodes` — two episode inputs refuse | KILLED | UNMEASURED IN SANDBOX |
| R-EP-GRAMMAR | `cmd/episode_preset.go:107` | `TestPresetStartReadRefusals/bad_episode` — explicit episode grammar | KILLED | UNMEASURED IN SANDBOX |
| R-OVERRIDE-AGENT | `cmd/episode_preset.go:118` | `TestPresetStartReadParity` — resolved agent uses actual explicit flag | KILLED | UNMEASURED IN SANDBOX |
| R-OVERRIDE-MODEL | `cmd/episode_preset.go:121` | `TestPresetStartReadInteractive` — explicit model survives agent change | KILLED | UNMEASURED IN SANDBOX |
| R-OVERRIDE-TTL | `cmd/episode_preset.go:125` | `TestPresetStartReadParity` — resolved and persisted TTL match the flag | KILLED | UNMEASURED IN SANDBOX |
| R-OVERRIDE-BRANCH | `cmd/episode_preset.go:128` | `TestPresetStartReadParity` — resolved branch follows the CLI input | KILLED | UNMEASURED IN SANDBOX |
| R-CLEAR-MODEL | `cmd/episode_preset.go:134` | `TestPresetStartReadParity` — changing agent clears inherited model | KILLED | UNMEASURED IN SANDBOX |
| R-HONOR-MODEL | `cmd/episode_preset.go:132` | `TestPresetStartReadInteractive` — explicit cross-agent model is kept | KILLED | UNMEASURED IN SANDBOX |
| R-BAD-AGENT | `cmd/episode_preset.go:136` | `TestPresetStartReadRefusals/bad_agent` — unknown agent refuses before health | KILLED | UNMEASURED IN SANDBOX |
| R-LAUNCH-PROBE | `cmd/episode_preset.go:140` | `TestPresetStartReadHeadless` — headless resolves the executable | KILLED | UNMEASURED IN SANDBOX |
| R-BINARY-ERR | `cmd/episode_preset.go:141` | `TestStartLaunchPreFenceTaskBinaryAndDirectory/binary` — missing fake binary has a specific refusal | KILLED | UNMEASURED IN SANDBOX |
| R-INTERACTIVE-PROBE | `cmd/episode_preset.go:145` | `TestStartLaunchValidatesBeforeFence` — missing interactive wrapper fails before tty | KILLED | UNMEASURED IN SANDBOX |
| R-INTERACTIVE-ERR | `cmd/episode_preset.go:146` | `TestStartLaunchValidatesBeforeFence` — refuse bad confined wrapper | KILLED | UNMEASURED IN SANDBOX |
| R-LAUNCH-DIR-ERR | `cmd/episode_preset.go:151` | `TestStartLaunchPreFenceTaskBinaryAndDirectory/directory` — nonempty launch directory refusal | KILLED | UNMEASURED IN SANDBOX |
| R-LAUNCH-LOCK-ERR | `cmd/episode_preset.go:154` | `TestPresetStartReadRefusals/launch_lock` — existing process lock refuses | KILLED | UNMEASURED IN SANDBOX |
| R-TTL-BOUND | `cmd/episode_preset.go:158` | `TestPresetStartReadRefusals/bad_ttl` — negative TTL receives the preset-specific refusal | KILLED | UNMEASURED IN SANDBOX |
| R-BRANCH-GRAMMAR | `cmd/episode_preset.go:169` | `TestPresetStartReadRefusals/bad_branch` — invalid branch refuses before tty | KILLED | UNMEASURED IN SANDBOX |
| R-BRANCH-EXISTS | `cmd/episode_preset.go:182` | `TestPresetStartReadRefusals/branch` — derived existing branch must refuse | KILLED | UNMEASURED IN SANDBOX |
| R-BRANCH-EXPLICIT | `cmd/episode_preset.go:182` | `TestPresetDerivedCollisionsAndExistingBranch` — free explicit existing branch accepted at terminal | KILLED | UNMEASURED IN SANDBOX |
| R-BRANCH-CHECKED-OUT | `cmd/episode_preset.go:190` | `TestPresetStartReadRefusals/checked_out` — branch assigned to another worktree refuses | KILLED | UNMEASURED IN SANDBOX |
| R-WT-REUSE | `cmd/episode_preset.go:232` | `TestPresetStartReadRefusals/worktree` — preset cannot adopt a linked worktree | KILLED | UNMEASURED IN SANDBOX |
| R-CLAUDE-MCP | `cmd/episode_preset.go:197` | `TestPresetStartReadParity` — Claude creates the selected MCP path | KILLED | UNMEASURED IN SANDBOX |
| R-INTERACTIVE-MCP | `cmd/episode_preset.go:200` | `TestPresetStartReadInteractive` — interactive resolves without terminal MCP file | KILLED | UNMEASURED IN SANDBOX |
| R-COLLISION-TOKEN | `cmd/episode_preset.go:203` | `TestPresetStartReadRefusals/token` — early token collision class before task or health | KILLED | UNMEASURED IN SANDBOX |
| R-COLLISION-TASK | `cmd/episode_preset.go:203` | `TestPresetStartReadRefusals/task` — reserved task path independently checked | KILLED | UNMEASURED IN SANDBOX |
| R-COLLISION-MCP | `cmd/episode_preset.go:203` | `TestPresetStartReadRefusals/mcp` — dangling MCP link is occupied | KILLED | UNMEASURED IN SANDBOX |
| R-LSTAT-FAULT | `cmd/episode_preset.go:209` | `TestPresetStartReadRefusals/loop` — other stat errors refuse | KILLED | UNMEASURED IN SANDBOX |
| R-TASK-READ | `cmd/episode_preset.go:215` | `TestPresetStartReadParity` — task bytes and digest are supplied by the read callback | KILLED | UNMEASURED IN SANDBOX |
| R-TASK-STAT-ERR | `cmd/episode_preset.go:293` | `TestPresetStartReadRefusals/unreadable_task` — missing task refusal | KILLED | UNMEASURED IN SANDBOX |
| R-TASK-TYPE | `cmd/episode_preset.go:297` | `TestPresetStartReadRefusals/task_directory` — only a regular task input read | KILLED | UNMEASURED IN SANDBOX |
| R-TASK-ERR | `cmd/episode_preset.go:217` | `TestPresetStartReadRefusals/unreadable_task` — task-reader fault stops before health | KILLED | UNMEASURED IN SANDBOX |
| R-PROFILE-DRIFT | `cmd/episode_preset.go:229` | `TestPresetStartReadRefusals/project_race` — profile bytes changing inside resolution refuse | KILLED | UNMEASURED IN SANDBOX |
| R-MINT-INPUT-ERR | `cmd/episode_preset.go:225` | `TestPresetStartReadRefusals/nonprivate_sessions` — preserve private-path refusal class | KILLED | UNMEASURED IN SANDBOX |
| R-PATH-PRIVATE | `cmd/episode_preset.go:225` | `TestPresetStartReadRefusals/nonprivate_sessions` — unsafe sessions mode refuses before prompt | KILLED | UNMEASURED IN SANDBOX |
| R-HEALTH-ERR | `cmd/episode_preset.go:277` | `TestPresetStartReadRefusals/health_missing` — transport error class preserved | KILLED | UNMEASURED IN SANDBOX |
| R-DAEMON-ERR | `cmd/episode_preset.go:281` | `TestPresetStartReadRefusals/world` — daemon validation failure cannot return a read plan | KILLED | UNMEASURED IN SANDBOX |
| R-WORLD-ID | `cmd/episode_preset.go:47` | `TestPresetStartReadRefusals/world` — served world identity exact | KILLED | UNMEASURED IN SANDBOX |
| R-DAEMON-BUILD | `cmd/episode.go:466` | `TestPresetStartReadRefusals/build` — served build equality | KILLED | UNMEASURED IN SANDBOX |
| R-DAEMON-TOOL | `cmd/episode.go:483` | `TestPresetStartReadRefusals/tool` — served tool hash equality | KILLED | UNMEASURED IN SANDBOX |
| R-DAEMON-CREDS | `cmd/episode.go:492` | `TestPresetStartReadRefusals/credentials` — served credentials identity exact | KILLED | UNMEASURED IN SANDBOX |
| R-DAEMON-PROFILE | `cmd/episode.go:509` | `TestPresetStartReadRefusals/profile` — served profile digest equality | KILLED | UNMEASURED IN SANDBOX |
| R-START-OID | `cmd/episode_preset.go:247` | `TestPresetStartReadParity` — mint uses resolved base rather than unrelated HEAD | KILLED | UNMEASURED IN SANDBOX |
| R-RESOLVED-WALL | `cmd/episode_preset.go:285` | `TestPresetStartReadHeadless` — exact launch duration returned | KILLED | UNMEASURED IN SANDBOX |
| R-RESOLVED-BIN | `cmd/episode_preset.go:285` | `TestPresetStartReadHeadless` — resolved binary identity is absolute | KILLED | UNMEASURED IN SANDBOX |
| R-TASK-SNAPSHOT | `cmd/episode_preset.go:285` | `TestPresetStartReadParity` — copy retains exact frozen bytes | KILLED | UNMEASURED IN SANDBOX |
| R-FLAG-LAUNCH | `cmd/episode_preset.go:317` | `TestStartLaunchHeadlessMintsAndLaunchesOnce` — terminal request actually launches once | KILLED | UNMEASURED IN SANDBOX |
| R-FLAG-HEADLESS | `cmd/episode_preset.go:317` | `TestStartLaunchHeadlessMintsAndLaunchesOnce` — explicit headless reaches terminal preflight | KILLED | UNMEASURED IN SANDBOX |
| R-FLAG-HOSTS | `cmd/episode_preset.go:316` | `TestPresetStartReadInteractive` — CLI passes exact custom host dir | KILLED | UNMEASURED IN SANDBOX |
| R-FLAG-AGENT | `cmd/episode_preset.go:321` | `TestPresetStartReadParity` — explicit agent flag reaches resolved plan | KILLED | UNMEASURED IN SANDBOX |
| R-FLAG-MODEL | `cmd/episode_preset.go:321` | `TestPresetStartReadInteractive` — explicit model flag reaches resolved plan | KILLED | UNMEASURED IN SANDBOX |
| R-FLAG-TTL | `cmd/episode_preset.go:321` | `TestPresetStartReadParity` — explicit ttl flag reaches resolved plan | KILLED | UNMEASURED IN SANDBOX |
| R-FLAG-BRANCH | `cmd/episode_preset.go:321` | `TestPresetStartReadParity` — explicit branch flag reaches resolved plan | KILLED | UNMEASURED IN SANDBOX |
| R-CLI-TASK | `cmd/episode_preset.go:325` | `TestPresetStartReadParity` — terminal task reaches prompt and copy | KILLED | UNMEASURED IN SANDBOX |
| R-CLI-NO-TASK | `cmd/episode_preset.go:325` | `TestPresetStartReadInteractive` — optional absent task leaves no copy | KILLED | UNMEASURED IN SANDBOX |
| R-MINT-ONLY | `cmd/episode_preset.go:368` | `TestPresetStartReadParity` — unrequested launch never reaches fake client | KILLED | UNMEASURED IN SANDBOX |
| R-LAUNCH-ONCE | `cmd/episode_preset.go:368` | `TestStartLaunchHeadlessMintsAndLaunchesOnce` — requested launch remains on terminal path | KILLED | UNMEASURED IN SANDBOX |
| P-ARTIFACT-CHECK | `host/store/operator_cas.go:59` | `TestOperatorConfirmReplay` — different proof must conflict and keep prior row | KILLED | authority KILLED, store SURVIVED |
| P-ARTIFACT-POLARITY | `host/store/operator_cas.go:59` | `TestOperatorConfirmReplay` — identical proof must preserve idempotence | KILLED | authority KILLED, store SURVIVED |

Equivalent or unchanged cases. Ignoring the optional empty MCP pathname during the path scan saves an `Lstat("")` that would return ENOENT and take the existing absence arm. Repeated pure `launchMode` calls on identical parsed values are now one calculation. The old `FlagSet.Lookup` nil checks were defensive because the only caller creates these flags unconditionally. The real tty fence, worktree add/remove logic, task/MCP exclusive writers, credentials migration and transaction code retain their previous predicates. No kill is claimed for unmodified branches without a ledger row.

## 5. Files

- `cmd/ailang-worldd/{approve,episode,episode_preset,episode_resume,session_new}.go`
- `cmd/ailang-worldd/{operator_provision_test,episode_preset_read_test}.go` (added)
- `docs/QUICKSTART.md`
- `host/authority/operator_store_test.go` (replay extended)
- `host/store/operator_cas.go`

Package set is `cmd/ailang-worldd`, `host/authority` and `host/store`. C7 inventories verify the unchanged command/help, context roots, imports, writer list and Go package surface.

## 6. Sprint gates

Extract the script with §2. Base and scratch use the repository read-only for Git refs. Pass the scratch changed-file list explicitly. The controller runs `bash "$T/gates.sh" "$W" "$HOME/.ailang/state/world-iter431/controller-gates" final` outside the sandbox with `GATES_SKIP_VERIFY_GO` and `GATES_SANDBOX` unset. `GATES_REAL_CREDS` can name the actual credentials path if it differs from the default. Observations compare hashes and line counts or report absence. Use Linux CI for the exact full race suite and Darwin for the read-only launchctl observation.

````bash gates.sh
#!/bin/bash
# Usage: gates.sh <tree> <logdir> <base|final|scratch> [sha]
set -u
set -o pipefail
export GIT_OPTIONAL_LOCKS=0
T=$1; L=$2; MODE=$3
B=cd06d9500ec2d34e975d7566f8bbdb14b1529527
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
     { grep -Fq 'OBSERVATION REQUIRES CONTROLLER' "$L/$name.txt" ||
       { grep -Eq 'listen .*bind: operation not permitted' "$L/$name.txt" &&
         ! grep -E '^ *--- FAIL:' "$L/$name.txt" | grep -Fv -e '--- FAIL: TestApproveDecideWriteMessageAndSignature ' -e '--- FAIL: TestQuickstartSection9FlagsMatchTheCLI '; }; }; then
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
  local d file jobs n creds
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
  if [ ! -x /bin/launchctl ]; then echo 'launchctl absent'; return 0; fi
  jobs=$(/bin/launchctl list)
  printf '%s\n' "$jobs" | grep -Eq '^PID[[:space:]]+Status[[:space:]]+Label$'
  n=$(printf '%s\n' "$jobs" | grep -c 'ailang-world' || test "$?" -eq 1)
  printf 'ailang-world launchctl jobs: %s\n' "$n"
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
EXPECTED='cmd/ailang-worldd host/authority host/store '
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
inv host/broker TestBrokerDependencyAllowlist
inv cmd/ailang-worldd TestEveryVerbHasHelp TestVerbTableMatchesTheCLISwitch TestCLICredentialWriterSites TestOperatorFlagParity TestLaunchAndStartHelp TestQuickstartSection19MatchesTheBar TestQuickstartSection20MatchesTheFinishFlag TestQuickstartServiceSectionMatchesTheCLIHelp TestPluginCommandsMatchTheCLIHelp TestHelpIsNeverBlockedByTheRegistryKeyGuard TestHomeServiceSeamOnly TestHomeServiceInstallAndPlist TestHomeServiceUpgradeRollback
leg inv-quickstart-section9 "go test -v -count=1 -timeout 8m ./cmd/ailang-worldd -run '^TestQuickstartSection9FlagsMatchTheCLI\$'"
if [ "${GATES_SANDBOX:-0}" != 1 ] && ! grep -Fq -- '--- PASS: TestQuickstartSection9FlagsMatchTheCLI (' "$L/inv-quickstart-section9.txt"; then
  FAILS="$FAILS inv-quickstart-section9_missing"
fi
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
  local allowed="$L/allowed.txt" extra
  cat > "$allowed" <<'EOF'
cmd/ailang-worldd/approve.go
cmd/ailang-worldd/episode.go
cmd/ailang-worldd/episode_preset.go
cmd/ailang-worldd/episode_preset_read_test.go
cmd/ailang-worldd/episode_resume.go
cmd/ailang-worldd/operator_provision_test.go
cmd/ailang-worldd/session_new.go
docs/QUICKSTART.md
host/authority/operator_store_test.go
host/store/operator_cas.go
design_docs/planned/w-world-address-bar-verbs-v4s5-sprint-plan.md
design_docs/planned/sprint_w-world-address-bar-verbs-v4s5.json
EOF
  printf '%s\n' cmd/ailang-worldd/episode_preset.go | grep -Fx -f "$allowed"
  printf '%s\n' world/transitions.ail | grep -Fxv -f "$allowed"
  extra=$(grep -Fxv -f "$allowed" "$L/paths.txt" || test "$?" -eq 1)
  printf 'UNTOUCHED_VIOLATIONS: %s\n' "${extra:-none}"
  test -z "$extra"
}
leg scope 'scope_fence'
caller_fence() {
  local old callers lock_count cases
  old=$(grep -R -nF --include='*.go' 'authority.RecheckVerb(' cmd host | grep -vF '_test.go:')
  printf 'Existing live caller control: %s\n' "$old"
  test -n "$old"
  if grep -R -nE --include='*.go' '\.(ClaimOperatorAction|ReserveOperatorStep|ObserveOperatorStep|RecheckClaimedOperatorAction)\(' cmd host | grep -vF '_test.go:'; then
    return 1
  fi
  lock_count=$(grep -R -nE --include='*.go' 'syscall\.Flock.*LOCK_NB' cmd host | grep -vF '_test.go:' | wc -l | tr -d ' ')
  printf 'Production Flock LOCK_NB lines (including the writer-lock comment): %s\n' "$lock_count"
  test "$lock_count" = 7
  cases=$(grep -R -nF --include='*.go' 'provisionEpisode(' cmd | grep -vF '_test.go:' | grep -vF 'func provisionEpisode(')
  printf 'Terminal provision callers:\n%s\n' "$cases"
  test "$(printf '%s\n' "$cases" | wc -l | tr -d ' ')" = 5
  if [ "$MODE" != base ]; then
    callers=$(grep -R -nF --include='*.go' 'resolvePresetStart(' cmd host | grep -vF '_test.go:' | grep -vF 'func resolvePresetStart(')
    printf 'Read seam caller:\n%s\n' "$callers"
    test "$(printf '%s\n' "$callers" | wc -l | tr -d ' ')" = 1
    printf '%s\n' "$callers" | grep -F 'cmd/ailang-worldd/episode_preset.go:'
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

````text
LEG observation-before rc=77 0s UNINFORMATIVE UNDER SANDBOX
LEG paths rc=0 0s
TOUCHED_PACKAGES: cmd/ailang-worldd host/authority host/store 
LEG touched rc=0 0s
LEG status rc=0 0s
LEG gofmt rc=0 0s
LEG build rc=0 1s
LEG vet rc=0 2s
LEG fence rc=0 1s
LEG verify_ail rc=0 51s
LEG pkg-cmd-ailang-worldd rc=1 4s UNINFORMATIVE UNDER SANDBOX
LEG race-cmd-ailang-worldd rc=1 22s UNINFORMATIVE UNDER SANDBOX
LEG pkg-host-authority rc=0 2s
LEG race-host-authority rc=0 12s
LEG pkg-host-store rc=0 15s
LEG race-host-store rc=0 37s
LEG inv-host-approval rc=0 1s
LEG inv-host-approval-approvaltest rc=0 0s
LEG inv-host-authority rc=0 1s
LEG inv-host-boundary rc=0 1s
LEG inv-host-daemon rc=0 2s
LEG inv-host-verbs rc=0 0s
LEG inv-host-broker rc=0 1s
LEG inv-cmd-ailang-worldd rc=0 1s
LEG inv-quickstart-section9 rc=1 1s UNINFORMATIVE UNDER SANDBOX
LEG inv-host-store rc=0 1s
LEG routes rc=0 0s
LEG scope rc=0 0s
LEG callers rc=0 0s
LEG verify_go SKIPPED 0s (GATES_SKIP_VERIFY_GO=1)
LEG observation-after rc=77 0s UNINFORMATIVE UNDER SANDBOX
GATES_DONE 156s FAILS: none UNINFORMATIVE: observation-before pkg-cmd-ailang-worldd race-cmd-ailang-worldd inv-quickstart-section9 observation-after
````

The first expanded sweep grouped the loopback-dependent QUICKSTART §9 test with the local inventories. It aborted subsequent tests on both base and scratch. The delivered script isolates that test and requires a real pass outside the sandbox. The local inventories now pass with 50 daemon routes, three surface POSTs, byte-identical daemon/CSP/app files and one legacy preset caller of the shared seam. The slot has no new daemon or verifygate fixtures.

The targeted CLI race selection is recorded in `$S/focused-race.log` (40.059 seconds with concurrent planner work). Run the same selection with this command.

````bash
AILANG_BIN=$HOME/.pinned-ailang/ailang go test -race -count=1 -timeout 8m -v ./cmd/ailang-worldd -run '^(TestOperatorProvisionEpisodeExclusion|TestPresetStartRead.*|TestStartLaunchHeadlessMintsAndLaunchesOnce|TestOperatorTerminalExclusion|TestOperatorEpisodeLockExclusive)$'
````

The controller runs the outside-sandbox script with its exported fixture settings. The following Go invocations require a controller run.

````bash
export AILANG_BIN="$HOME/.pinned-ailang/ailang"
go test -count=1 -timeout 8m ./cmd/ailang-worldd
go test -count=1 -race -timeout 8m ./cmd/ailang-worldd
go test -v -count=1 -timeout 8m ./cmd/ailang-worldd -run '^TestQuickstartSection9FlagsMatchTheCLI$'
VERIFY_GO_LEGS=plain ./scripts/verify_go.sh
VERIFY_GO_LEGS=race ./scripts/verify_go.sh
````

The race command for the repository uses the existing Linux CI settings. The controller also runs the §6 before/after file and launchctl observation on Darwin.

## 7. Executor directive notes

- For each stage run `AILANG_BIN=$HOME/.pinned-ailang/ailang bash "$T/accept.sh" "$W" M1 fast`, then M2 and M3 after applying their patches. The script supplies all fixture environment and stores interpreter caches under scratch. All three stage trees passed gofmt, `go build ./...`, targeted vet and their named acceptance. The extracted acceptance logs are `$S/compiled-final-M*.log`.
- Run `bash "$T/snapshot.sh" "$W" "$M"`. Commit-time command is `bash "$HOME/.ailang/state/world-iter431/commit_ms.sh" "$M" "$msgfile" "AILANG_BIN=$HOME/.pinned-ailang/ailang bash '$T/accept.sh' '$W' '$M' fast"`. Each commit's reported numstat must match §1.5. The snapshot enumerates product pathspecs without assuming corresponding directories exist.
- Run `bash "$HOME/.ailang/state/world-iter431/ctl_roundtrip.sh" M1 M2 M3` on a fresh export. Check each product directory diff status, plus the full export comparison. `SNAP_STATUS_LIST` and `SNAP_PARENT` are scratch rehearsal overrides. The executor uses the default destination and Git read-only path discovery.

````bash snapshot.sh
#!/bin/bash
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=$1
M=$2
B=cd06d9500ec2d34e975d7566f8bbdb14b1529527
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter431/snap}/$M
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
M1 applied (     249 lines)
  diff -r cmd vs st-M1: rc=0
  diff -r host vs st-M1: rc=0
  diff -r docs vs st-M1: rc=0
M2 applied (     714 lines)
  diff -r cmd vs st-M2: rc=0
  diff -r host vs st-M2: rc=0
  diff -r docs vs st-M2: rc=0
M3 applied (      35 lines)
  diff -r cmd vs st-M3: rc=0
  diff -r host vs st-M3: rc=0
  diff -r docs vs st-M3: rc=0
DONE
M1 full-tree diff -r rc=0
M2 full-tree diff -r rc=0
M3 full-tree diff -r rc=0
````

Manifest checks passed for M1 7/0/7, M2 8/0/8 and M3 10/0/10. A host-only rehearsal produced 1/0/1 with the cmd, docs and app directories absent.

## 8. Risks

- The lock covers independent open descriptions on Darwin and Linux. Contention may be reported as `EWOULDBLOCK` or its `EAGAIN` alias. Existing shared-lock tests open independent descriptions while the holder remains live and verify reacquisition after close. Keep the inode stable. Symlink/mode refusal assertions use portable classes and compare captured state within the same OS.
- A successful read ceases to be current as soon as trusted operator inputs change. S5b compares frozen values under the episode lock and rechecks authority at the mint and reservation boundaries. Legacy CLI typed overrides remain available.
- The current preset read seam inspects local files, Git refs and health. The interactive terminal preflight retains its fake-testable bounded version probe. The future browser route depends on a short-lived runner catalogue and an independently revalidated exact snapshot.
- Base and final whole CLI package coverage is uninformative under sandbox loopback denial. The isolated CLI regression and race tests can run now. The controller still owes the full outside-sandbox leg and Linux race gate.

## 9. Acceptance criteria

1. M1–M3 apply from `cd06d950`, reproduce compiled stages byte-for-byte, and match the cumulative snapshot manifests and filtered totals.
2. Both formerly unguarded legacy provisioning routes refuse the exact durable owner, preserve input and persisted state, and release on successful mint. Prior terminal-only routes, the pre-split migration, command/help inventory and private sidecar documentation are covered.
3. A read-only typed preset resolution supplies the same effective start fields and frozen task used by the attended CLI. Fixtures exercise local and health refusal branches with no mutation. No production operator action, reserve, observe or launch dispatch becomes reachable.
4. Identical replay is stable and a different artifact conflicts without advancing approver or counter. Named behavioral mutations pass §4's standard. The controller reruns unmeasured legs before promotion.
5. S5b retains the §5.3 atomic mint, worktree, file and authority tests together with pre-dispatch reconciliation. Preserve all pending §2.2/§5.5 activation gates.

## Appendix C — after this sprint

- **S5b atomic proof and recovery.** The current product cut is 743, leaving 657. The reviewed design has 510 raw lines (689 calibrated) for this sub-milestone, before the resolver/drain and before the 132-line bank's remaining gates. Reuse the shared resolver to reconstruct every signed field over fresh local inputs. Compare preset, base/ref/OID, repository/remote, profile/tool identity, grants/budgets/TTL, canonical paths, task bytes/digest, daemon/store identities and catalogue/owner generation independently. Mint needs a distinct sealed proof from the runner recheck. The tentative transaction binds a fresh clock, current approver (including recorded passkey state), action/id/digest/generation/owner, credential binding and start row. Preserve the A3 consume and public typed paths. Reserve worktree using exact repo/branch/starting SHA and registration identity. An error with positively established absent state stops or retries under the original live proof, and partial/unattributable state enters attended reconciliation. Mint receipt commits with the credential and start, pins the credential id, and carries no plaintext. Recover an exact valid private token file or revoke only the committed orphan; never mint again for that action. Publish files exclusively at canonical paths, validate file type, owner/mode, exact bytes and path independently. The complete row mapping is in §3.2.
- **Attended resolver.** C5 measures a status success and cancel conflict while uncertainty stays owned. S5 introduces zero uncertain writes. The first milestone introducing such an outcome must supply an attended command with tested exact remedy text, effect-specific readback and independent current-owner/step CAS. Under the same lock, a matching durable receipt permits adoption, positive non-execution permits an explicit reviewed disposition, and unresolved evidence retains the owner and records the precise blocker. Preserve worktree and all unrelated credentials. Keep the finish resolver precondition under D-WORLD-120 before S9.
- **Runner consumer and drain before S5b/S6.** Current production caller count is zero. Before any reachable claim/reserve/observe, acquire the exact canonical episode lock before CAS and keep the ownership fenced through worktree, mint, file installs and any bounded waits. Require named process-to-process pins for first-managed transition, old-build terminal drain, current-owner recovery and upgrade with an active child. No stale-heartbeat shortcut. Gate stop, uninstall, matched restart and schema-compatible rollback using drain or retained uncertainty. Existing direct terminal session mint/revoke are outside `provisionEpisode`. The activation proof must account for concurrent credential changes and revalidate the exact current credential before launch or recovery.
- **S3b.** Scratch re-application measured 132 changed lines and `TestOperatorStepRestart` passed. The bank alone has no process fd-transfer pin, full receipt-predicate mutation matrix or current-owner recovery gate. Those must land with any production consumer. No S3b hunk is in Appendix P.
- **J-OWN-UID and S4b audit notes.** Keep the actual stat UID check and controller alternate-owner fixture. Carry the unchanged `Z-ABSENT-ANY` missing-parent race and `Z-WAIT-READBUDGET` budget granularity cases with S4b pre-dispatch hardening. `Z-RDWR` is equivalent for the tested platforms and the lock remains private. The schema-v10 backup remains attended before a home upgrade.
- **Launch and browser stages.** S6 owns supervised launch, expiry, confinement and positive process receipt/restart evidence. S7 owns routes/status and use of the S4c replay contract. S8 owns the end-to-end browser case and #604 sequencing, S9–S11 the two-proof finish path and E2/E5, S12 the herdr path and E4.

## Appendix P — compiled patches

### Patch M1

````diff M1
diff --git a/cmd/ailang-worldd/approve.go b/cmd/ailang-worldd/approve.go
index b360411..4492966 100644
--- a/cmd/ailang-worldd/approve.go
+++ b/cmd/ailang-worldd/approve.go
@@ -395,7 +395,7 @@ func runApprove(args []string, stdout, stderr io.Writer, env sessionEnv) int {
 	req.answerBy = row.ExpiresAt
 	req.consume = &authority.RequestMint{Request: mr, RequestCreatedAt: row.CreatedAt, Challenge: [32]byte(challenge),
 		ExpiresAt: row.ExpiresAt, Signed: signed, Consume: store.RequestConsume{ID: row.RequestID, Digest: row.Digest, Clock: env.now}}
-	if code := provisionEpisode(req, stdout, stderr, env); code != exitOK {
+	if code := provisionEpisode(req, nil, stdout, stderr, env); code != exitOK {
 		return code
 	}
 	printLaunch(stdout, launchInfo{episode: mr.Episode, agent: mr.Agent, model: mr.Model, addr: mr.Daemon.Addr,
diff --git a/cmd/ailang-worldd/episode.go b/cmd/ailang-worldd/episode.go
index a82ca94..b4f5e09 100644
--- a/cmd/ailang-worldd/episode.go
+++ b/cmd/ailang-worldd/episode.go
@@ -374,7 +374,7 @@ func runEpisodeStart(args []string, stdout, stderr io.Writer, env sessionEnv) in
 	}
 
 	// 2.-3. session new's path, unchanged.
-	if code := provisionEpisode(req, stdout, stderr, env); code != exitOK {
+	if code := provisionEpisode(req, effectLock, stdout, stderr, env); code != exitOK {
 		return code
 	}
 
diff --git a/cmd/ailang-worldd/episode_preset.go b/cmd/ailang-worldd/episode_preset.go
index e9f9f6d..1fed21a 100644
--- a/cmd/ailang-worldd/episode_preset.go
+++ b/cmd/ailang-worldd/episode_preset.go
@@ -309,7 +309,7 @@ func runPresetEpisodeStart(fs *flag.FlagSet, pos []string, name, dir, profilesDi
 		}
 		return nil
 	}
-	if code := provisionEpisode(req, out, errOut, env); code != exitOK {
+	if code := provisionEpisode(req, effectLock, out, errOut, env); code != exitOK {
 		return code
 	}
 	if launching {
diff --git a/cmd/ailang-worldd/episode_resume.go b/cmd/ailang-worldd/episode_resume.go
index 1d1f004..b58bc84 100644
--- a/cmd/ailang-worldd/episode_resume.go
+++ b/cmd/ailang-worldd/episode_resume.go
@@ -212,7 +212,7 @@ func runEpisodeResume(fs *flag.FlagSet, pos []string, ep, name, presetsDir, prof
 		}
 		return resumeAfterRename()
 	}
-	if code := provisionEpisode(req, out, errOut, env); code != exitOK {
+	if code := provisionEpisode(req, effectLock, out, errOut, env); code != exitOK {
 		return code
 	}
 	fmt.Fprintf(out, "next: ailang-worldd episode launch %s --preset %s\n", ep, name)
diff --git a/cmd/ailang-worldd/operator_provision_test.go b/cmd/ailang-worldd/operator_provision_test.go
new file mode 100644
index 0000000..7152f78
--- /dev/null
+++ b/cmd/ailang-worldd/operator_provision_test.go
@@ -0,0 +1,136 @@
+package main
+
+import (
+	"bytes"
+	"io"
+	"os"
+	"path/filepath"
+	"reflect"
+	"runtime/debug"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func TestOperatorProvisionEpisodeExclusion(t *testing.T) {
+	gc := debug.SetGCPercent(-1)
+	t.Cleanup(func() { debug.SetGCPercent(gc) })
+
+	t.Run("no_terminal", func(t *testing.T) {
+		r := newEpisodeRun(t)
+		creds, err := credentialsFor(r.f.db, "")
+		if err != nil {
+			t.Fatal(err)
+		}
+		before := observeProvision(t, r, creds)
+		var out bytes.Buffer
+		code := runSessionNew(r.f.args("ep", "--preset", "se-tools"), &out, &out, noTerminalEnv(1800))
+		_, sidecarErr := os.Lstat(operatorEpisodeLockPath(creds, "ep"))
+		if code != exitUsage || !strings.Contains(out.String(), "no controlling terminal") ||
+			!os.IsNotExist(sidecarErr) || !reflect.DeepEqual(observeProvision(t, r, creds), before) {
+			t.Fatalf("session new bypassed the terminal fence: %d %s", code, &out)
+		}
+	})
+
+	t.Run("session_new", func(t *testing.T) {
+		r := newEpisodeRun(t)
+		cred, err := credentialsFor(r.f.db, "")
+		if err != nil {
+			t.Fatal(err)
+		}
+		held := claimFixtureEpisode(t, cred, "ep1")
+		if err := held.Close(); err != nil {
+			t.Fatal(err)
+		}
+		// The lease is stale and the synthetic holder has exited. The
+		// durable action must still take precedence over a typed yes.
+		var out bytes.Buffer
+		term := &promptTerm{in: strings.NewReader("y\n")}
+		env := testMintEnv("y\n", 1800)
+		env.openTerminal = func() (io.ReadWriteCloser, error) { return term, nil }
+		token := filepath.Join(r.sessions, "custom.session")
+		if err := os.Mkdir(r.sessions, 0700); err != nil {
+			t.Fatal(err)
+		}
+		before := observeProvision(t, r, cred)
+		code := runSessionNew(r.f.args("ep1", "--preset", "se-tools", "--branch", "new", "--out", token), &out, &out, env)
+		expectEpisodeRefusal(t, code, out.String(), before.credentials, cred)
+		if term.out.Len() != 0 || !reflect.DeepEqual(observeProvision(t, r, cred), before) {
+			t.Fatal("contention opened a prompt or changed Git, store or session state")
+		}
+
+		// A different episode still mints, and the normal exit releases
+		// its own descriptor. Do not rely on GC to close a leaked file.
+		out.Reset()
+		if code := runSessionNew(r.f.args("other", "--preset", "se-tools", "--out", token), &out, &out, env); code != exitOK {
+			t.Fatalf("independent episode: %d %s", code, &out)
+		}
+		if rows := r.f.rows(t); len(rows) != 1 || rows[0].EpisodeID != "other" {
+			t.Fatalf("independent mint did not persist: %+v", rows)
+		}
+		requireProvisionLockReleased(t, cred, "other")
+	})
+	t.Run("approve", func(t *testing.T) {
+		a, sig := signedRig(t)
+		held := claimFixtureEpisode(t, a.credPath, a.ep)
+		before := observeProvision(t, a.r, a.credPath)
+		request, approvers := a.row(t), approverRows(t, a.credPath)
+		var out bytes.Buffer
+		var opened bool
+		env := approveEnv(nil, &opened, atTime(approveT0+1))
+		code := runApprove(a.args("--signature", sig), &out, &out, env)
+		expectEpisodeRefusal(t, code, out.String(), before.credentials, a.credPath)
+		if opened || !reflect.DeepEqual(observeProvision(t, a.r, a.credPath), before) ||
+			!reflect.DeepEqual(a.row(t), request) || !reflect.DeepEqual(approverRows(t, a.credPath), approvers) {
+			t.Fatal("contended signature changed request, approver, Git, store or session state")
+		}
+		if err := held.Close(); err != nil {
+			t.Fatal(err)
+		}
+	})
+	t.Run("approve_releases", func(t *testing.T) {
+		a, sig := signedRig(t)
+		var out bytes.Buffer
+		var opened bool
+		code := runApprove(a.args("--signature", sig), &out, &out, approveEnv(nil, &opened, atTime(approveT0+1)))
+		if code != exitOK || opened || !a.row(t).Consumed {
+			t.Fatalf("uncontended approval: %d %s", code, &out)
+		}
+		requireProvisionLockReleased(t, a.credPath, a.ep)
+	})
+}
+
+type provisionSnapshot struct {
+	credentials, worktrees, branches string
+	workspace, sessions              map[string]string
+}
+
+func observeProvision(t *testing.T, r episodeRun, cred string) provisionSnapshot {
+	t.Helper()
+	return provisionSnapshot{serviceFileHash(cred), presetGit(t, r.f, "worktree", "list", "--porcelain"),
+		presetGit(t, r.f, "for-each-ref", "--format=%(refname) %(objectname)", "refs/heads"),
+		provisionFiles(t, r.f.ws), provisionFiles(t, r.sessions)}
+}
+
+func provisionFiles(t *testing.T, dir string) map[string]string {
+	t.Helper()
+	if _, err := os.Lstat(dir); os.IsNotExist(err) {
+		return nil
+	}
+	return launchFileSnapshot(t, dir)
+}
+
+func requireProvisionLockReleased(t *testing.T, cred, ep string) {
+	t.Helper()
+	id, err := store.CanonicalPath(cred)
+	if err != nil {
+		t.Fatal(err)
+	}
+	path := operatorEpisodeLockPath(id, ep)
+	fd, err := openRunnerLock(path, 0)
+	if err != nil {
+		t.Fatalf("cannot reacquire %s: %v", path, err)
+	}
+	fd.Close()
+}
diff --git a/cmd/ailang-worldd/session_new.go b/cmd/ailang-worldd/session_new.go
index 5450bc8..4455c40 100644
--- a/cmd/ailang-worldd/session_new.go
+++ b/cmd/ailang-worldd/session_new.go
@@ -402,7 +402,7 @@ func checkBindingFits(credPath, episode string, p *broker.ExecProfile) error {
 // provisionEpisode is steps 2-7: the fence, the y/N, the credentials lock,
 // the worktree, the mint and the rollback, in that order, for both session
 // new and episode start. It returns the exit code.
-func provisionEpisode(r *mintRequest, stdout, stderr io.Writer, env sessionEnv) int {
+func provisionEpisode(r *mintRequest, held *os.File, stdout, stderr io.Writer, env sessionEnv) int {
 	prefix := "ailang-worldd " + r.verb + ": "
 	// 2. The mint fence, unchanged, before any side effect, unless a verified signature approved
 	// this request (row 201 A1b §2.8: Signed is set only by approve, from ApproveBySignature).
@@ -418,6 +418,17 @@ func provisionEpisode(r *mintRequest, stdout, stderr io.Writer, env sessionEnv)
 		defer func() { _ = term.Close() }()
 	}
 
+	// Start and resume already hold the lock. Session new and mint-request
+	// approval acquire it here, after the terminal fence and before the prompt.
+	if held == nil {
+		lock, err := terminalEpisodeLock(r.credPath, r.episode, nil)
+		if err != nil {
+			fmt.Fprintf(stderr, "%srefusing: %v\n", prefix, err)
+			return exitUsage
+		}
+		defer lock.Close()
+	}
+
 	// 3. y/N naming the worktree, every grant and the exec profile binding.
 	now := env.now()
 	grants := append([]broker.Capability(nil), r.grants...)
@@ -646,7 +657,7 @@ func runSessionNew(args []string, stdout, stderr io.Writer, env sessionEnv) int
 	if req.prof != nil && !hasGrant(grants, broker.EffectWorkspaceExec) {
 		fmt.Fprintf(stderr, "ailang-worldd session new: note: --profile %s binds workspace-exec, but no Workspace.Exec grant is minted\n", req.prof.Project)
 	}
-	return provisionEpisode(req, stdout, stderr, env)
+	return provisionEpisode(req, nil, stdout, stderr, env)
 }
 
 // hasGrant reports whether grants hold effect at any scope.
diff --git a/docs/QUICKSTART.md b/docs/QUICKSTART.md
index c67cfc6..49659ff 100644
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -74,6 +74,10 @@ ailang-worldd service runner stop
 ailang-worldd service runner uninstall
 ```
 Review and type `y` for install, stop and uninstall. Uninstall before upgrading or replacing the config. Spool and action records are retained.
+Terminal episode operations use private mode-0600 sidecars named
+`<creds>.operator-episode-<ep>.lock` beside the canonical credentials database.
+Keep these sidecars in place across stop and uninstall. Independent opens need
+the same file to coordinate access to the episode.
 Run the checked example with a fake supervisor and disposable stores:
 ```sh
 AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./cmd/ailang-worldd -run '^TestOperatorServiceWriteScope$' -count=1 -v
````

### Patch M2

````diff M2
diff --git a/cmd/ailang-worldd/episode_preset.go b/cmd/ailang-worldd/episode_preset.go
index 1fed21a..5a1459d 100644
--- a/cmd/ailang-worldd/episode_preset.go
+++ b/cmd/ailang-worldd/episode_preset.go
@@ -50,108 +50,124 @@ func checkPresetDaemon(h daemon.HealthResponse, req *mintRequest, p *resolvedPre
 	return checkDaemon(h, req, addr, episodeSelfBuild())
 }
 
-func runPresetEpisodeStart(fs *flag.FlagSet, pos []string, name, dir, profilesDir, task string, specs []string, out, errOut io.Writer, env sessionEnv) int {
-	fail := func(err error) int {
-		fmt.Fprintln(errOut, "ailang-worldd episode start: "+err.Error())
-		return exitUsage
-	}
-	if len(specs) != 0 {
-		return fail(presetRefusal("grants come from preset "+name+" (F7)", "grants and checks are closed", "edit the preset, or drop --preset"))
-	}
-	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
-	defer cancel()
-	p, err := loadLaunchPreset(ctx, name, dir, profilesDir, defaultWorldPaths().workspaceRoot)
+type presetStartOverride struct {
+	name, value string
+}
+
+type presetStartInput struct {
+	episodes                []string
+	workspaceRoot, hostsDir string
+	launching, headless     bool
+	overrides               []presetStartOverride
+	readTask                func() (source string, data []byte, err error)
+}
+
+type presetStartRead struct {
+	mint          *mintRequest
+	preset        *resolvedPreset
+	health        daemon.HealthResponse
+	taskCopy, mcp string
+	task          []byte
+	launchMode    string
+	launchWall    time.Duration
+	launchBinary  string
+}
+
+// resolvePresetStart reads the local preset, Git, credential and daemon inputs
+// for review. It does not lock, prompt, provision files, mint or launch. The
+// terminal supplies a task-file reader. A runner can supply frozen inline bytes.
+func resolvePresetStart(ctx context.Context, name, dir, profilesDir string, input presetStartInput, now func() int64) (*presetStartRead, error) {
+	p, err := loadLaunchPreset(ctx, name, dir, profilesDir, input.workspaceRoot)
 	if err != nil {
-		return fail(err)
+		return nil, err
 	}
-	launching := fs.Lookup("launch") != nil && fs.Lookup("launch").Value.String() == "true"
-	headless := fs.Lookup("headless") != nil && fs.Lookup("headless").Value.String() == "true"
-	if launching {
-		if _, _, err := launchMode(p.Preset, headless); err != nil {
-			return fail(err)
+	var mode string
+	var wall time.Duration
+	if input.launching {
+		mode, wall, err = launchMode(p.Preset, input.headless)
+		if err != nil {
+			return nil, err
 		}
-		mode, _, _ := launchMode(p.Preset, headless)
-		if mode == "headless" && task == "" {
-			return fail(presetRefusal("--launch needs --task", "headless launch runs the pinned task", "add --task <file>"))
+		if mode == "headless" && input.readTask == nil {
+			return nil, presetRefusal("--launch needs --task", "headless launch runs the pinned task", "add --task <file>")
 		}
 	}
-	if len(pos) > 1 {
-		return fail(errors.New("expected at most one episode id"))
+	if len(input.episodes) > 1 {
+		return nil, errors.New("expected at most one episode id")
 	}
 	ep := ""
-	if len(pos) == 1 {
-		ep = pos[0]
+	if len(input.episodes) == 1 {
+		ep = input.episodes[0]
 	} else {
-		ep, err = derivePresetEpisode(name, p.Daemon.WorkspaceRoot, p.Daemon.SessionsDir, env.now())
+		ep, err = derivePresetEpisode(name, p.Daemon.WorkspaceRoot, p.Daemon.SessionsDir, now())
 		if err != nil {
-			return fail(err)
+			return nil, err
 		}
 	}
 	if !workspaceEpisodePattern.MatchString(ep) {
-		return fail(fmt.Errorf("episode %q does not match %s", ep, workspaceEpisodePattern))
+		return nil, fmt.Errorf("episode %q does not match %s", ep, workspaceEpisodePattern)
 	}
 	agent, model, ttl, branch := p.Preset.Agent, p.Preset.Model, p.Preset.TTLS, "world/"+ep
 	overrides := []string{}
 	overrideLines := []string{}
 	modelOverridden := false
-	fs.Visit(func(f *flag.Flag) {
-		switch f.Name {
+	for _, override := range input.overrides {
+		switch override.name {
 		case "agent":
-			overrideLines = append(overrideLines, "agent "+f.Value.String()+" (preset: "+agent+")")
-			agent = f.Value.String()
+			overrideLines = append(overrideLines, "agent "+override.value+" (preset: "+agent+")")
+			agent = override.value
 		case "model":
-			overrideLines = append(overrideLines, "model "+f.Value.String()+" (preset: "+model+")")
-			model = f.Value.String()
+			overrideLines = append(overrideLines, "model "+override.value+" (preset: "+model+")")
+			model = override.value
 			modelOverridden = true
 		case "ttl":
-			overrideLines = append(overrideLines, "ttl "+f.Value.String()+" (preset: "+strconv.FormatInt(ttl, 10)+")")
-			ttl, _ = strconv.ParseInt(f.Value.String(), 10, 64)
+			overrideLines = append(overrideLines, "ttl "+override.value+" (preset: "+strconv.FormatInt(ttl, 10)+")")
+			ttl, _ = strconv.ParseInt(override.value, 10, 64)
 		case "branch":
-			overrideLines = append(overrideLines, "branch "+f.Value.String()+" (preset: "+branch+")")
-			branch = f.Value.String()
-		default:
-			return
+			overrideLines = append(overrideLines, "branch "+override.value+" (preset: "+branch+")")
+			branch = override.value
 		}
-		overrides = append(overrides, f.Name)
-	})
+		overrides = append(overrides, override.name)
+	}
 	if agent != p.Preset.Agent && !modelOverridden && model != "" {
 		overrideLines = append(overrideLines, "model <default for "+agent+"> (cleared preset "+p.Preset.Agent+" model: "+model+")")
 		model = ""
 	}
 	if agent != "claude" && agent != "codex" && agent != "agy" {
-		return fail(errors.New("agent must be claude, codex or agy"))
+		return nil, errors.New("agent must be claude, codex or agy")
 	}
-	if launching {
-		if _, err := launchBinary(agent); err != nil {
-			return fail(err)
+	var binary string
+	if input.launching {
+		binary, err = launchBinary(agent)
+		if err != nil {
+			return nil, err
 		}
-		mode, _, _ := launchMode(p.Preset, headless)
 		if mode == "interactive" {
-			if _, err := launchLauncher(agent, fs.Lookup("hosts-dir").Value.String(), p.Daemon.Addr); err != nil {
-				return fail(err)
+			if _, err := launchLauncher(agent, input.hostsDir, p.Daemon.Addr); err != nil {
+				return nil, err
 			}
 		}
 		launchDir := filepath.Join(filepath.Dir(p.Daemon.SessionsDir), "launch", ep)
 		if err := validateLaunchDir(launchDir, p.Daemon.WorkspaceRoot, filepath.Join(p.Daemon.WorkspaceRoot, ep), p.Daemon.SessionsDir); err != nil {
-			return fail(err)
+			return nil, err
 		}
 		if err := launchLockCheck(filepath.Join(p.Daemon.SessionsDir, ep+".launch.lock")); err != nil {
-			return fail(err)
+			return nil, err
 		}
 	}
 	if ttl < 0 {
-		return fail(errors.New("ttl must be non-negative"))
+		return nil, errors.New("ttl must be non-negative")
 	}
 	git, err := sessionGit()
 	if err != nil {
-		return fail(err)
+		return nil, err
 	}
 	query := func(args ...string) (string, error) {
 		s, err := worktree.Output(ctx, git, append([]string{"-C", p.Preset.Repo}, args...))
 		return strings.TrimSpace(s), err
 	}
 	if _, err := query("check-ref-format", "refs/heads/"+branch); err != nil {
-		return fail(err)
+		return nil, err
 	}
 	existing := false
 	if _, err := query("show-ref", "--verify", "refs/heads/"+branch); err == nil {
@@ -164,15 +180,15 @@ func runPresetEpisodeStart(fs *flag.FlagSet, pos []string, name, dir, profilesDi
 		}
 	}
 	if existing && !branchOverridden {
-		return fail(errors.New("derived branch already exists; pass --branch explicitly or choose another episode"))
+		return nil, errors.New("derived branch already exists; pass --branch explicitly or choose another episode")
 	}
 	if existing {
 		list, err := query("worktree", "list", "--porcelain")
 		if err != nil {
-			return fail(err)
+			return nil, err
 		}
 		if strings.Contains("\n"+list+"\n", "\nbranch refs/heads/"+branch+"\n") {
-			return fail(fmt.Errorf("existing branch %s is checked out in another worktree; choose a free branch", branch))
+			return nil, fmt.Errorf("existing branch %s is checked out in another worktree; choose a free branch", branch)
 		}
 	}
 	sessionFile := filepath.Join(p.Daemon.SessionsDir, ep+".session")
@@ -181,69 +197,55 @@ func runPresetEpisodeStart(fs *flag.FlagSet, pos []string, name, dir, profilesDi
 	if agent == "claude" {
 		mcp = filepath.Join(p.Daemon.SessionsDir, ep+".mcp.json")
 	}
-	if launching {
-		mode, _, _ := launchMode(p.Preset, headless)
-		if mode == "interactive" {
-			mcp = ""
-		}
+	if mode == "interactive" {
+		mcp = ""
 	}
 	for _, path := range []string{sessionFile, taskCopy, mcp} {
 		if path == "" {
 			continue
 		}
 		if _, err := os.Lstat(path); err == nil {
-			return fail(fmt.Errorf("%s already exists; choose another episode", path))
+			return nil, fmt.Errorf("%s already exists; choose another episode", path)
 		} else if !errors.Is(err, os.ErrNotExist) {
-			return fail(err)
+			return nil, err
 		}
 	}
 	var taskData []byte
 	var taskRecord *store.StartTask
-	if task != "" {
-		path, err := filepath.Abs(task)
-		if err != nil {
-			return fail(err)
-		}
-		info, err := os.Stat(path)
-		if err != nil {
-			return fail(err)
-		}
-		if !info.Mode().IsRegular() {
-			return fail(errors.New("task must be a regular file"))
-		}
-		taskData, err = os.ReadFile(path)
+	if input.readTask != nil {
+		var path string
+		path, taskData, err = input.readTask()
 		if err != nil {
-			return fail(err)
+			return nil, err
 		}
-		title := taskTitle(taskData)
-		taskRecord = &store.StartTask{Path: path, SHA256: serviceHash(taskData), Title: title}
+		taskRecord = &store.StartTask{Path: path, SHA256: serviceHash(taskData), Title: taskTitle(taskData)}
 	}
 	req, err := prepareMint(mintOptions{verb: "episode start", episode: ep, db: p.Daemon.DB, credsDB: p.Daemon.CredentialsDB,
 		wsRoot: p.Daemon.WorkspaceRoot, repo: p.Preset.Repo, branch: branch, grants: p.Grants, ttl: ttl,
 		outPath: sessionFile, outDirCreate: true, profile: p.Preset.Profile, profilesDir: profilesDir})
 	if err != nil {
-		return fail(err)
+		return nil, err
 	}
 	if req.prof != nil && req.prof.Digest != p.ProfileDigest {
-		return fail(errors.New("exec profile changed during validation; run episode start again"))
+		return nil, errors.New("exec profile changed during validation; run episode start again")
 	}
 	if req.reuse {
-		return fail(errors.New("launch preset episode already has a worktree; choose another episode"))
+		return nil, errors.New("launch preset episode already has a worktree; choose another episode")
 	}
 	// A CLI never upgrades an older credentials store.
 	if _, err := os.Stat(req.credPath); err == nil {
 		ro, err := store.OpenCredentialsReadOnly(req.credPath)
 		if err != nil {
-			return fail(err)
+			return nil, err
 		}
 		version := ro.Version()
 		_ = ro.Close()
 		if version < store.CredentialsBindingWriteVersion {
-			return fail(&store.UnboundCredentialsStoreError{Path: req.credPath})
+			return nil, &store.UnboundCredentialsStoreError{Path: req.credPath}
 		}
 	}
 	req.startPoint, req.existingBranch = p.BaseSHA, existing
-	req.start = &store.EpisodeStart{Schema: "world/episode-start/v1", Episode: ep, StartedAt: time.Unix(env.now(), 0).Format(time.RFC3339),
+	req.start = &store.EpisodeStart{Schema: "world/episode-start/v1", Episode: ep, StartedAt: time.Unix(now(), 0).Format(time.RFC3339),
 		Preset: store.StartPreset{Name: name, Path: p.Path, SHA256: p.SHA256, Bytes: p.File}, Repo: p.Preset.Repo,
 		Remote: store.StartRemote{Name: p.Preset.Remote, URL: p.RemoteURL}, Base: store.StartBase{Branch: p.Preset.Base, SHA: p.BaseSHA},
 		Branch: branch, BranchExisted: existing, Agent: agent, Model: model, Overrides: overrides, Task: taskRecord, Checks: p.Checks, TTLS: ttl}
@@ -274,16 +276,67 @@ func runPresetEpisodeStart(fs *flag.FlagSet, pos []string, name, dir, profilesDi
 	req.presetPrompt = prompt.String()
 	h, err := fetchDaemonHealth(p.Daemon.Addr)
 	if err != nil {
-		return fail(err)
+		return nil, err
 	}
 	if err := checkPresetDaemon(h, req, p, p.Daemon.Addr); err != nil {
+		return nil, err
+	}
+	return &presetStartRead{mint: req, preset: p, health: h, taskCopy: taskCopy, mcp: mcp,
+		task: taskData, launchMode: mode, launchWall: wall, launchBinary: binary}, nil
+}
+
+func readPresetStartTask(task string) (string, []byte, error) {
+	path, err := filepath.Abs(task)
+	if err != nil {
+		return "", nil, err
+	}
+	info, err := os.Stat(path)
+	if err != nil {
+		return "", nil, err
+	}
+	if !info.Mode().IsRegular() {
+		return "", nil, errors.New("task must be a regular file")
+	}
+	data, err := os.ReadFile(path)
+	return path, data, err
+}
+
+func runPresetEpisodeStart(fs *flag.FlagSet, pos []string, name, dir, profilesDir, task string, specs []string, out, errOut io.Writer, env sessionEnv) int {
+	fail := func(err error) int {
+		fmt.Fprintln(errOut, "ailang-worldd episode start: "+err.Error())
+		return exitUsage
+	}
+	if len(specs) != 0 {
+		return fail(presetRefusal("grants come from preset "+name+" (F7)", "grants and checks are closed", "edit the preset, or drop --preset"))
+	}
+	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
+	defer cancel()
+
+	input := presetStartInput{
+		episodes: pos, workspaceRoot: defaultWorldPaths().workspaceRoot, hostsDir: fs.Lookup("hosts-dir").Value.String(),
+		launching: fs.Lookup("launch").Value.String() == "true", headless: fs.Lookup("headless").Value.String() == "true",
+	}
+	fs.Visit(func(f *flag.Flag) {
+		switch f.Name {
+		case "agent", "model", "ttl", "branch":
+			input.overrides = append(input.overrides, presetStartOverride{f.Name, f.Value.String()})
+		}
+	})
+	if task != "" {
+		input.readTask = func() (string, []byte, error) { return readPresetStartTask(task) }
+	}
+	read, err := resolvePresetStart(ctx, name, dir, profilesDir, input, env.now)
+	if err != nil {
 		return fail(err)
 	}
+	req, p := read.mint, read.preset
+	ep := req.episode
 	effectLock, err := terminalEpisodeLock(req.credPath, ep, nil)
 	if err != nil {
 		return fail(err)
 	}
 	defer effectLock.Close()
+	taskRecord, taskCopy, taskData, mcp := req.start.Task, read.taskCopy, read.task, read.mcp
 	req.afterToken = func(token string) error {
 		if taskRecord != nil {
 			f, err := os.OpenFile(taskCopy, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
@@ -312,14 +365,14 @@ func runPresetEpisodeStart(fs *flag.FlagSet, pos []string, name, dir, profilesDi
 	if code := provisionEpisode(req, effectLock, out, errOut, env); code != exitOK {
 		return code
 	}
-	if launching {
-		return launchEpisode(ep, name, dir, headless, out, errOut, effectLock, fs.Lookup("hosts-dir").Value.String())
+	if input.launching {
+		return launchEpisode(ep, name, dir, input.headless, out, errOut, effectLock, input.hostsDir)
 	}
 	taskPath := ""
 	if taskRecord != nil {
 		taskPath = taskRecord.Path
 	}
-	printLaunch(out, launchInfo{episode: ep, agent: agent, model: model, addr: p.Daemon.Addr, worktree: req.dir,
-		branch: branch, sessionFile: sessionFile, mcpFile: mcp, task: taskPath, grants: p.Grants})
+	printLaunch(out, launchInfo{episode: ep, agent: req.start.Agent, model: req.start.Model, addr: p.Daemon.Addr, worktree: req.dir,
+		branch: req.branch, sessionFile: req.outPath, mcpFile: mcp, task: taskPath, grants: p.Grants})
 	return exitOK
 }
diff --git a/cmd/ailang-worldd/episode_preset_read_test.go b/cmd/ailang-worldd/episode_preset_read_test.go
new file mode 100644
index 0000000..5cee58b
--- /dev/null
+++ b/cmd/ailang-worldd/episode_preset_read_test.go
@@ -0,0 +1,334 @@
+package main
+
+import (
+	"bytes"
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
+	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/worktree"
+)
+
+func TestPresetStartReadParity(t *testing.T) {
+	f := newPresetFixture(t)
+	oldClient := launchClient
+	t.Cleanup(func() { launchClient = oldClient })
+	launchClient = func(string) *client {
+		t.Fatal("a start without --launch attempted a launch")
+		return nil
+	}
+	creds := matchingHealth(t, f.r.f, f.r.profiles, f.profile).CredentialsDB
+	task := readStartFixtureFile(t, f.task)
+	before := observeProvision(t, f.r, creds)
+	input := presetStartInput{
+		episodes: []string{"ep"}, workspaceRoot: f.r.f.ws,
+		overrides: []presetStartOverride{{"agent", "claude"}, {"branch", "pinned"}, {"ttl", "120"}},
+		readTask:  func() (string, []byte, error) { return f.task, task, nil },
+	}
+	read, err := resolvePresetStart(presetTestContext(t), "p", f.dir, f.r.profiles, input, presetClock)
+	if err != nil {
+		t.Fatal(err)
+	}
+	req := read.mint
+	if req.episode != "ep" || req.dir != filepath.Join(f.r.f.ws, "ep") || req.outPath != filepath.Join(f.r.sessions, "ep.session") ||
+		read.taskCopy != filepath.Join(f.r.sessions, "ep.task.md") || read.mcp != filepath.Join(f.r.sessions, "ep.mcp.json") ||
+		req.startPoint != f.base || req.branch != "pinned" || req.existingBranch || req.reuse || req.ttl != 120 ||
+		req.start.Agent != "claude" || req.start.Model != "" || req.start.Base.SHA != f.base ||
+		req.start.Preset.SHA256 != serviceHash(f.data) || req.prof.Digest != f.profile.Digest ||
+		req.start.Task == nil || req.start.Task.SHA256 != serviceHash(task) || req.start.Task.Title != "Fix fixture" ||
+		read.health.ToolRef == "" || read.health.CredentialsDB != creds || !bytes.Equal(read.task, task) ||
+		!reflect.DeepEqual(observeProvision(t, f.r, creds), before) {
+		t.Fatalf("read seam diverged or wrote state: %+v", req.start)
+	}
+	if _, err := os.Lstat(operatorEpisodeLockPath(creds, "ep")); !os.IsNotExist(err) {
+		t.Fatalf("read seam touched episode lock: %v", err)
+	}
+
+	term := &promptTerm{in: strings.NewReader("y\n")}
+	env := presetEnv(term)
+	env.openTerminal = func() (io.ReadWriteCloser, error) {
+		// Resolution has completed. Only the already-read bytes may reach
+		// the task copy, start row or prompt.
+		seedStartFixtureFile(t, f.task, []byte("changed after resolution"))
+		return term, nil
+	}
+	var out bytes.Buffer
+	if code := runEpisodeStart(f.args("ep", "--agent", "claude", "--branch", "pinned", "--ttl", "120", "--task", f.task), &out, &out, env); code != exitOK {
+		t.Fatal(code, out.String())
+	}
+	if !strings.HasPrefix(term.out.String(), req.presetPrompt) ||
+		!strings.Contains(req.presetPrompt, "model <default for claude>") ||
+		!strings.Contains(req.presetPrompt, "task      "+f.task+" \"Fix fixture\" ("+serviceHash(task)+")") ||
+		!bytes.Equal(readStartFixtureFile(t, read.taskCopy), task) {
+		t.Fatal("the attended path did not use the resolved review and bytes")
+	}
+	got := readPresetStart(t, f, "ep")
+	want := *req.start
+	dir := filepath.Join(f.r.f.ws, "ep")
+	head, err := worktree.Output(presetTestContext(t), f.r.f.git, []string{"-C", dir, "rev-parse", "HEAD"})
+	rows := f.r.f.rows(t)
+	token := readStartFixtureFile(t, req.outPath)
+	mcp := readStartFixtureFile(t, read.mcp)
+	if len(rows) != 1 || got == nil || err != nil || strings.TrimSpace(head) != f.base || !isHex64(string(token)) ||
+		!bytes.Contains(mcp, append([]byte("Bearer "), token...)) {
+		t.Fatalf("worktree, credentials or files differ from the resolved plan: rows=%d, HEAD=%s, err=%v", len(rows), head, err)
+	}
+	if !reflect.DeepEqual(got, &want) || rows[0].ExpiresAt != presetClock()+120 || rows[0].ProfileDigest != f.profile.Digest {
+		t.Fatalf("mint changed reviewed start or binding: %+v", got)
+	}
+}
+
+func TestPresetStartReadHeadless(t *testing.T) {
+	f := newPresetFixture(t)
+	bin := startReadFakeAgent(t, f)
+	creds := matchingHealth(t, f.r.f, f.r.profiles, f.profile).CredentialsDB
+	before := observeProvision(t, f.r, creds)
+	task := []byte("# Literal\n$(touch should-never-run)\n")
+	input := presetStartInput{episodes: []string{"ep"}, workspaceRoot: f.r.f.ws,
+		launching: true, headless: true, hostsDir: t.TempDir(),
+		readTask: func() (string, []byte, error) { return "inline", task, nil },
+	}
+	read, err := resolvePresetStart(presetTestContext(t), "p", f.dir, f.r.profiles, input, presetClock)
+	if err != nil {
+		t.Fatal(err)
+	}
+	binary, err := filepath.EvalSymlinks(filepath.Join(bin, "codex"))
+	if err != nil {
+		t.Fatal(err)
+	}
+	if read.launchMode != "headless" || read.launchWall != time.Hour || read.launchBinary != binary ||
+		read.mint.start.Task.Path != "inline" || read.mint.start.Task.SHA256 != serviceHash(task) ||
+		!bytes.Equal(read.task, task) || read.mcp != "" || !reflect.DeepEqual(before, observeProvision(t, f.r, creds)) {
+		t.Fatal("typed resolver lost headless inputs or changed durable state")
+	}
+	if _, err := os.Lstat(filepath.Join(bin, "ran")); !os.IsNotExist(err) {
+		t.Fatal("a read started the executable")
+	}
+}
+
+func TestPresetStartReadInteractive(t *testing.T) {
+	host, hosts, _ := interactiveFixture(t, "claude", "exit 99")
+	f := newPresetFixture(t)
+	t.Setenv("HOME", filepath.Join(host.root, "home"))
+	creds := matchingHealth(t, f.r.f, f.r.profiles, f.profile).CredentialsDB
+	before := observeProvision(t, f.r, creds)
+	input := presetStartInput{episodes: []string{"ep"}, workspaceRoot: f.r.f.ws, launching: true, hostsDir: hosts,
+		overrides: []presetStartOverride{{"agent", "claude"}, {"model", "explicit"}}}
+	read, err := resolvePresetStart(presetTestContext(t), "p", f.dir, f.r.profiles, input, presetClock)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if read.launchMode != "interactive" || read.launchWall != 0 || read.mcp != "" || read.mint.start.Task != nil ||
+		read.mint.start.Model != "explicit" || !reflect.DeepEqual(before, observeProvision(t, f.r, creds)) {
+		t.Fatal("interactive resolution changed mint semantics or durable state")
+	}
+	term := &promptTerm{in: strings.NewReader("n\n")}
+	var out bytes.Buffer
+	code := runEpisodeStart(f.args("ep", "--launch", "--hosts-dir", hosts, "--agent", "claude", "--model", "explicit"), &out, &out, presetEnv(term))
+	if code != exitUsage || !strings.Contains(out.String(), "aborted (not confirmed)") ||
+		!strings.Contains(term.out.String(), "agent     claude explicit") ||
+		!reflect.DeepEqual(before, observeProvision(t, f.r, creds)) {
+		t.Fatalf("terminal interactive validation drifted: %d %s", code, &out)
+	}
+	if _, err := os.Lstat(filepath.Join(host.facts, "started")); !os.IsNotExist(err) {
+		t.Fatal("interactive resolution started a child")
+	}
+}
+
+func TestPresetStartReadRefusals(t *testing.T) {
+	for _, tc := range []struct{ name, want string }{
+		{"branch", "derived branch already exists"},
+		{"checked_out", "checked out in another worktree"},
+		{"worktree", "already has a worktree"},
+		{"token", ".session already exists; choose another episode"},
+		{"task", ".task.md already exists"},
+		{"mcp", ".mcp.json already exists"},
+		{"loop", syscall.ELOOP.Error()},
+		{"project_race", "exec profile changed during validation"},
+		{"task_directory", "task must be a regular file"},
+		{"unreadable_task", "no such file"},
+		{"bad_ttl", "ttl must be non-negative"},
+		{"bad_agent", "agent must be"},
+		{"two_episodes", "expected at most one"},
+		{"bad_episode", "does not match"},
+		{"bad_branch", "check-ref-format"},
+		{"health_missing", "fixture daemon offline"},
+		{"world", "(F6)"},
+		{"build", "is build"},
+		{"tool", "serves the tool binary"},
+		{"credentials", "which it never reads"},
+		{"profile", "changed since the daemon"},
+		{"invalid_wall", "launch deadline is too large"},
+		{"launch_lock", "only one launch may run"},
+		{"nonprivate_sessions", "must be private"},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			f := newPresetFixture(t)
+			h := matchingHealth(t, f.r.f, f.r.profiles, f.profile)
+			h.DBPath = f.r.f.db
+			if _, err := store.UpgradeCredentials(h.CredentialsDB); err != nil {
+				t.Fatal(err)
+			}
+			if err := os.Mkdir(f.r.sessions, 0700); err != nil {
+				t.Fatal(err)
+			}
+			args := f.args("ep", "--agent", "claude", "--task", f.task)
+			var healthErr error
+			healthCalls, wantHealth := 0, 0
+			switch tc.name {
+			case "branch":
+				presetGit(t, f.r.f, "branch", "world/ep", "dev")
+			case "checked_out":
+				presetGit(t, f.r.f, "branch", "world/used", "dev")
+				presetGit(t, f.r.f, "worktree", "add", filepath.Join(t.TempDir(), "used"), "world/used")
+				args = append(args, "--branch", "world/used")
+			case "worktree":
+				presetGit(t, f.r.f, "worktree", "add", "--detach", filepath.Join(f.r.f.ws, "ep"), "dev")
+			case "token":
+				seedStartFixtureFile(t, filepath.Join(f.r.sessions, "ep.session"), []byte("keep"))
+			case "task":
+				seedStartFixtureFile(t, filepath.Join(f.r.sessions, "ep.task.md"), []byte("keep"))
+			case "mcp":
+				target := filepath.Join(t.TempDir(), "missing")
+				if err := os.Symlink(target, filepath.Join(f.r.sessions, "ep.mcp.json")); err != nil {
+					t.Fatal(err)
+				}
+			case "loop":
+				link := filepath.Join(filepath.Dir(f.r.sessions), "loop")
+				if err := os.Symlink(link, link); err != nil {
+					t.Fatal(err)
+				}
+				var p launchPreset
+				if err := json.Unmarshal(f.data, &p); err != nil {
+					t.Fatal(err)
+				}
+				p.Daemon, _ = json.Marshal(presetDaemon{Addr: defaultEpisodeAddr, DB: f.r.f.db, CredentialsDB: h.CredentialsDB,
+					WorkspaceRoot: f.r.f.ws, SessionsDir: filepath.Join(link, "sessions")})
+				changed, _ := json.Marshal(p)
+				seedStartFixtureFile(t, f.file, changed)
+			case "project_race":
+				profile := filepath.Join(f.r.profiles, "proj.json")
+				raw := readStartFixtureFile(t, profile)
+				input := presetStartInput{episodes: []string{"ep"}, workspaceRoot: f.r.f.ws, readTask: func() (string, []byte, error) {
+					seedStartFixtureFile(t, profile, append(raw, '\n'))
+					return f.task, nil, nil
+				}}
+				before := observeProvision(t, f.r, h.CredentialsDB)
+				if read, err := resolvePresetStart(presetTestContext(t), "p", f.dir, f.r.profiles, input, presetClock); read != nil ||
+					err == nil || !strings.Contains(err.Error(), tc.want) || !reflect.DeepEqual(observeProvision(t, f.r, h.CredentialsDB), before) {
+					t.Fatalf("profile drift allowed or wrote state: %v", err)
+				}
+				return
+			case "task_directory":
+				args = append(args, "--task", filepath.Dir(f.task))
+			case "unreadable_task":
+				args = append(args, "--task", f.task+"-missing")
+			case "bad_ttl":
+				args = append(args, "--ttl", "-1")
+			case "bad_agent":
+				args = append(args, "--agent", "other")
+			case "two_episodes":
+				args = append(args, "other")
+			case "bad_episode":
+				// Enter launch preflight before prepareMint can repeat the
+				// same grammar check. The resolver must reject the episode
+				// even though the requested agent would also be missing.
+				bin := t.TempDir()
+				if err := os.Symlink(f.r.f.git, filepath.Join(bin, "git")); err != nil {
+					t.Fatal(err)
+				}
+				t.Setenv("PATH", bin)
+				args = f.args("BAD", "--task", f.task, "--launch", "--headless")
+			case "bad_branch":
+				args = append(args, "--branch", "bad branch")
+			case "health_missing":
+				healthErr, wantHealth = errors.New("fixture daemon offline"), 1
+			case "world":
+				h.DBPath, wantHealth = filepath.Join(t.TempDir(), "other.db"), 1
+			case "build":
+				h.Build, wantHealth = daemon.BuildIdentity{}, 1
+			case "tool":
+				h.ToolRef, wantHealth = "sha256:other", 1
+			case "credentials":
+				h.CredentialsDB, wantHealth = filepath.Join(t.TempDir(), "other.db"), 1
+			case "profile":
+				h.ExecProfiles[0].Digest, wantHealth = "sha256:other", 1
+			case "invalid_wall":
+				var p launchPreset
+				if err := json.Unmarshal(f.data, &p); err != nil {
+					t.Fatal(err)
+				}
+				p.Launch = &presetLaunch{Mode: "headless", DeadlineS: 1e12}
+				data, _ := json.Marshal(p)
+				seedStartFixtureFile(t, f.file, data)
+				args = append(args, "--launch", "--headless")
+			case "launch_lock":
+				startReadFakeAgent(t, f)
+				args = append(args, "--agent", "codex", "--launch", "--headless")
+				seedStartFixtureFile(t, filepath.Join(f.r.sessions, "ep.launch.lock"), []byte("0"))
+			case "nonprivate_sessions":
+				if err := os.Chmod(f.r.sessions, 0750); err != nil {
+					t.Fatal(err)
+				}
+			}
+			stubDaemon(t, h, healthErr)
+			fetchDaemonHealth = func(addr string) (daemon.HealthResponse, error) {
+				healthCalls++
+				return h, healthErr
+			}
+			creds, err := credentialsFor(f.r.f.db, "")
+			if err != nil {
+				t.Fatal(err)
+			}
+			before := observeProvision(t, f.r, creds)
+			var opened bool
+			var out bytes.Buffer
+			code := runEpisodeStart(args, &out, &out, spyEnv(&opened))
+			if code != exitUsage || !strings.Contains(out.String(), tc.want) || opened || healthCalls != wantHealth ||
+				!reflect.DeepEqual(before, observeProvision(t, f.r, creds)) {
+				t.Fatalf("read refusal changed state, effect count or class: code=%d opened=%v health=%d %s", code, opened, healthCalls, &out)
+			}
+		})
+	}
+}
+
+func startReadFakeAgent(t *testing.T, f presetFixture) string {
+	t.Helper()
+	bin := t.TempDir()
+	gitLink := filepath.Join(bin, "git")
+	if err := os.Symlink(f.r.f.git, gitLink); err != nil {
+		t.Fatal(err)
+	}
+	binary := filepath.Join(bin, "codex")
+	if err := os.WriteFile(binary, []byte("#!/bin/sh\n: > \"$WORLD_FIXTURE_AGENT_RAN\"\nexit 99\n"), 0700); err != nil {
+		t.Fatal(err)
+	}
+	t.Setenv("PATH", bin)
+	t.Setenv("WORLD_FIXTURE_AGENT_RAN", filepath.Join(bin, "ran"))
+	return bin
+}
+
+func seedStartFixtureFile(t *testing.T, path string, data []byte) {
+	t.Helper()
+	if err := os.WriteFile(path, data, 0600); err != nil {
+		t.Fatal(err)
+	}
+}
+
+func readStartFixtureFile(t *testing.T, path string) []byte {
+	t.Helper()
+	data, err := os.ReadFile(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return data
+}
````

### Patch M3

````diff M3
diff --git a/host/authority/operator_store_test.go b/host/authority/operator_store_test.go
index 363cb31..b62c11b 100644
--- a/host/authority/operator_store_test.go
+++ b/host/authority/operator_store_test.go
@@ -85,6 +85,16 @@ func TestOperatorConfirmReplay(t *testing.T) {
 		if err != nil || !reflect.DeepEqual(got, want) {
 			t.Fatalf("confirm replay changed the row or refused: %v", err)
 		}
+		other := bytes.Clone(r.artifact)
+		other[len(other)-1] ^= 1
+		_, mismatch := c.ConfirmOperatorAction(ctx, r.cas, other, r.proof, 1021)
+		if !errors.Is(mismatch, store.ErrOperatorConflict) {
+			t.Errorf("different replay artifact returned %v", mismatch)
+		}
+		after, err := c.OperatorAction(ctx, r.cas.ID, 1021)
+		if err != nil || !reflect.DeepEqual(after, want) {
+			t.Fatal("invalid replay changed durable approval")
+		}
 		name := "mark-ssh"
 		if passkey {
 			name = "mark-mac"
diff --git a/host/store/operator_cas.go b/host/store/operator_cas.go
index d4e32a1..20db2e5 100644
--- a/host/store/operator_cas.go
+++ b/host/store/operator_cas.go
@@ -56,6 +56,9 @@ func (c *Credentials) ConfirmOperatorAction(ctx context.Context, cas OperatorCAS
 		return OperatorActionRow{}, ErrOperatorConflict
 	}
 	if row.ApprovedBy != "" {
+		if !bytes.Equal(row.Artifact, artifact) {
+			return OperatorActionRow{}, ErrOperatorConflict
+		}
 		return row, nil
 	}
 	if row.State != "requested" || row.Generation != cas.Generation {
````
