# Row 243 V4 S5f — start pre-caller pins and dormant consumer selection

Iteration 439, planner lane `codex:gpt-6.2-sol`. Base `9610f5de0921fe1c8eeef939928ffe639bd2259b`. Apply M1, M2, then M3. Appendix C records the compiled S5f2 caller prototype and the outstanding activation gates.

## 0. Binding constraints

- Reviewed V4 r1, quorum r2 PROCEED 3/3 with all seats present. D-WORLD-119=A and D-WORLD-120=A are attended rulings at `82f05004`. S5e4 landed at `e554db6b` via #676, 929 changed lines, Sonnet PASS 88. The latest merge is the #677 record.
- There is no live consumer in the staged milestones. M3 stages a read-only start selector, a lock-scoped visit, and the signed executable and catalogue-generation prerequisites. The live claim, worktree, mint, file publication and drain are a single activation unit. S6 owns the confined launch. The routes, app, SQL schema, and AILANG inputs have no diff.
- Real fixtures use a migrated credentials database, the repository's signing helper, real filesystem evidence and a re-executed Go test binary with fixture-only arguments. Use temp roots and compare physical paths. Use pipe EOF and a child context deadline for the process barrier. The controller performs any eventual service replacement under the attended service contract.

## 1. Pre-flight

### 1.1 Baseline at 9610f5de

`git rev-parse HEAD` returned the full base above. Both base-reference numstat commands returned zero rows. `git diff --name-only e554db6b HEAD` returned seven paths under `design_docs/` and `git show -s --oneline 82f05004` named the two attended A rulings. `$S/anchors-base.txt` records the exact commands, output and status. A pristine `git archive` export is at `$S/base-tree`. Its gate started before fixture edits. The TEXT raw logs, parsed failure sets, parser controls, load and timing are under `$S/gate-base`.

### 1.2 Current inputs and coupling inventory

All paths in the measurements are relative to the base unless a stage is named. `cmd/` below means `cmd/ailang-worldd/`. Rerun `/bin/bash "$S/anchors.sh" "$W"` for the full readback.

| Input | Command and output |
|---|---|
| Installed runner | `nl -ba cmd/ailang-worldd/operator_runner.go` printed `runOperatorRunner` :42, kernel runner lock :54, SIGTERM/Interrupt registration :64, heartbeat-only `runnerTick` :80, pair checks :81–86 and `No effect dispatch or GitHub probe` :89. |
| Mint and process census | `git grep -nF 'MintOperatorStart(' -- cmd host ':!*_test.go'` printed only `host/authority/mint.go:204:func MintOperatorStart`. The positive control in `cmd/operator_mint_test.go:82` calls it. `git grep -nE 'TestOperator(FirstManagedEpisodeProcess|OldBuildTerminalDrain|CurrentOwnerProcessRecovery)\|active_child' -- cmd host` exited 1 with no matches. `git grep -nF 'func TestOperatorServiceUpgrade' -- cmd host` printed `operator_service_test.go:302` and `operator_upgrade_test.go:13`. |
| Store lookup fixture | `nl -ba host/store/operator_episode_gc{,_test}.go` printed SQL `s.kind <> 'operation'` :23 in the lookup and the effect row created by `f.reserveWorktree(t)` :34 in the table. The secondary episode inserts a `cleanup` row :55. No existing same-episode `mint` or `files` row protects the lookup independently. |
| Mint recovery seam | `nl -ba cmd/ailang-worldd/operator_mint.go` printed first `readRunnerMint` :25, `openRunnerLock(..., 0)` :29, deferred close :33, second read :34, episode comparison :38, fresh assignment :41, token check :42 and expiry branch :85. `rg -n 'reconcileRunnerMint\('` found `operator_service.go:287`. `nl -ba host/store/operator_mint.go` printed the independent transaction reread/identity comparison :302–307. |
| Existing ownership and I/O seams | `git grep -nE 'func .* (OperatorAction\|ClaimOperatorAction\|ReserveOperatorWorktree\|ObserveOperatorStep\|ReserveOperatorStep\|PutOperatorHeartbeat\|ReadOperatorMint\|ReadOperatorFiles\|ResolveSession)\(' -- host/store` found the action reader :60, CAS :119, worktree reservation :13, generic reservation :17 and observe :45, heartbeat :46, mint reader :233, files reader :67, and session reader `credentials.go:561`. The same command with the full selector in the anchor receipt gives every referenced store method. `nl -ba` found start reconstruction `cmd/operator_start_input.go:84`, worktree `cmd/operator_worktree.go:31`, files `cmd/operator_files.go:262` and held/open/name separation at `cmd/operator_episode_lock.go:14,18,22,36`. |
| Case and liveness clocks | `nl -ba cmd/ailang-worldd/operator_runner.go` printed `min(remaining, time.Second)` :115, fake clock injection :109, :111 and :116 and poll injection :124. `rg -n 'func budgetContext' cmd` printed `cli.go:90`. `nl -ba host/store/operator_runner.go` printed catalogue generation increment :54 and :62. `nl -ba host/store/operator_{worktree,mint}.go` printed exact generation comparison in worktree :53 and mint :120. |
| Terminal callers | `git grep -nF 'terminalEpisodeLock(' -- cmd/ailang-worldd/session.go cmd/ailang-worldd/session_new.go` found mint :269 and shared `provisionEpisode` :427. `nl -ba cmd/ailang-worldd/session.go` printed `runSessionRevoke` :361, the episode branch :413, resolution by credential hash :422 and direct revoke :433. That function contains no episode-lock call, with mint as the same-file positive control. |
| Activation and GC | `git grep -nE 'HasOpenOperatorStep\|removeWorktree\(ctx' -- host/daemon/auto_gc.go cmd/ailang-worldd/episode_gc.go` printed daemon lookup :85–86, core eligibility :158–167, removal :175 and explicit CLI lookup :113. The `operatorEpisodeLock` search over both GC paths exited 1 with no matches, controlled by `cmd/operator_episode_lock.go:14`. |
| Upgrade | `git grep -nE 'runner is still installed\|cannot verify that the operator runner is unloaded\|rollbackUpgrade\(\|Bootout\(e.domain, runnerLabel\)\|ExitTimeOut' -- cmd/ailang-worldd/{service_upgrade,service,operator_service}.go` returned refusal at `service_upgrade.go:208,212`, rollback entry :257 and guard :264–268, stop at `operator_service.go:418` and the supervisor's 30 s exit bound at `service.go:197`. |
| Reviewed executable and confinement | `git grep -nE 'launch_binary\|Executables\[\|launch requires\|read-only", "-C"' -- cmd/ailang-worldd/{operator_start_input,operator_service,episode_launch,episode_launch_text}.go` printed the action fact `launch_binary` :70, configured executable hashes :145–154, path resolution :151–165, and Codex flags `-s read-only` :164. The exact `launch_binary_sha256` search was empty with the path fact as its positive control. |
| Worktree/file recovery windows | `nl -ba host/store/operator_worktree.go` printed the read gate `phase = 'reserved'` :89 and later-kind exclusion :90. `nl -ba host/store/operator_mint.go` printed observed mint only :245 with the start/credential join :254. `nl -ba cmd/ailang-worldd/operator_files_reconcile.go` printed the reserved-file gate :32. `git grep -nE 'func .* (NextOperatorStart\|.*RecoverOperatorClaim)\|reconcile-start\|launch_binary_sha256' -- cmd host` returned empty. Same-file controls were `OperatorAction`, `ClaimOperatorAction`, `reconcile-mint`, and `launch_binary`. |

| Inventory at base | Staged change and reason |
|---|---|
| `TestDaemonDependencyAllowlist`, `host/daemon/daemon_test.go:899` | No snapshot change. The daemon package retains its imports. Gate `inv-daemon`. |
| `TestWorldBoundaryDependencyAllowlist`, `host/boundary/allowlist_world_test.go:758` | No expectation change. Store uses its existing SQLite boundary and the CLI uses existing store and verbs edges. Gate `inv-boundary`. |
| `TestProductionGoSurface`, `host/store/context_roots_test.go:239` | No expectation change. Added production source is under `cmd/`. Gate `inv-store`. |
| `TestProductionContextRoots`, `host/store/context_roots_test.go:134` | No root addition. The selector accepts its caller context and demands a deadline via `operatorAccess`. The wait still uses `budgetContext`. Name this inventory in M3 acceptance. |
| `TestCredentialsMethodsRequireADeadline`, `host/store/credentials_test.go:277` | M3 adds `NextOperatorStart` to the credential deadline matrix and expects `ErrNoDeadline`. Run it and the positive bounded reader. |
| `TestStoreSurfaceTakesContext`, `host/store/deadline_guard_test.go:43` | No expectation change. The new exported credentials reader has a first-argument context. Gate `inv-store`. |
| `TestBoundaryASTWriteGuard`, `host/boundary/allowlist_world_test.go:1323` | No expectation change. The new command-side seam opens with the existing kernel-lock helper and opens SQLite read-only. Gate `inv-boundary`. |
| `TestOperatorNeverDispatch`, `host/verbs/operator_policy_test.go:69` | No expectation change even when the internal consumer is wired. :70–78 forbid app dispatch of approver, session mint/revoke, service, serve, setup and world-publish verbs with `OperatorNever`. :79–84 require `OperatorOther` for non-operator commands and :85–88 refuse shell syntax via `ErrShell`. A timer consuming a durable signed action uses the internal operator-runner mode. Keep all those exact negative arms and run the census for the internal caller separately. |
| `TestCLICredentialWriterSites`, `cmd/cli_credential_writers_test.go:18` | No expectation change. The seam opens a reader. The existing control logs read-only opens and pins the writer site map including `runnerWrite`. |

`git grep -nE 'func (TestDaemonDependencyAllowlist\|TestWorldBoundaryDependencyAllowlist\|TestProductionGoSurface\|TestProductionContextRoots\|TestCredentialsMethodsRequireADeadline\|TestStoreSurfaceTakesContext\|TestBoundaryASTWriteGuard\|TestOperatorNeverDispatch\|TestCLICredentialWriterSites)' -- cmd host` printed all nine anchors above. The base gate ran all nine by name and observed nine `--- PASS:` controls.

### 1.3 Decisions

- The selector sorts eligible confirmed starts by creation time and ID. It reads again under the physical episode lock and requires an equal row. SQL comparisons use an injected clock. The lock-scoped callback receives the exact open description. The read-only credentials handle allows a competing writer to open during the visit. The candidate itself writes neither ownership nor effect evidence.
- GC keeps its existing fail-safe behavior. If the credentials lookup is missing, old, locked or otherwise unavailable, the report skips reclamation with `operator step lookup unavailable`. Evaluator 438 finding 2 identifies the cost for installations without compatible evidence. No behavioral change to that policy occurs in this sprint.
- Activation requires mutual exclusion with GC for a physical episode from final eligibility inspection through removal, including reuse of a pathname while an earlier GC pass is in progress. Preserve creation-only worktree policy. Use the same canonical credentials DB plus episode key in GC, terminal and runner, and establish an unambiguous store-to-workspace binding. Both GC entry points must skip with an explicit reason on unavailable lock/ledger evidence, and recheck Git registration and eligibility while holding the fd. A post-lookup lease or an in-memory mutex cannot arbitrate another process. Pin both lock orders with a barrier and check that no branch, worktree or files changed on refusal. Appendix C carries this before activation.

### 1.4 Premises found false or incomplete

| Premise | Measured consequence |
|---|---|
| The existing upgrade path can drain and replace a paired runner | `service_upgrade.go:201–213` refuses an installed or still-loaded companion before touching a binary. `operator_service.go:412–429` only delegates to `Bootout` and then removes the installation on uninstall. The matched pair restart, wait for kernel quiescence and compatible pair rollback need a single attended transition. |
| Existing receipts cover all inter-effect drain and crash gaps | `ReadOperatorWorktree` :89 sees only a reservation. `ReadOperatorMint` :245 sees only an observation and :254 requires its session. `ReadOperatorFiles` and its CLI accept a narrower window. In particular a prepared worktree can be observed with no mint receipt, a mint reservation can survive a rollback, and the full files receipt leaves a running owner. There is no attended transition from an untouched claimed row and no constrained preparatory handoff to S6 from the later gaps. Never park a live owner without a runnable remedy. |
| The dormant helpers already use the runner's open file description | `cmd/operator_worktree.go:40–46` and `cmd/operator_files.go:269–273` independently open and close the canonical sidecar. A second open description conflicts with the first. The next stage must borrow and verify the runner's fd, including through the real Git child's completion. |
| Build pinning is signed for the agent bytes | The launch fact binds the path. `cmd/operator_service.go:148–155` checks configured hashes, with a distinct authority record. M3 adds the exact content digest to the signed facts and pins same-path replacement. The activating consumer still has to require the attended executable entry and supported confinement independently. |
| Catalogue can be treated as a continuously refreshed heartbeat timestamp | Every heartbeat advances the catalogue generation, and the current mint/worktree CAS requires the signed generation. The scratch core reproduced a `mint_catalogue` binding refusal on one identical heartbeat, with unchanged evidence. M3 stabilizes unchanged generations and pins update and restart invalidation. A complete bounded eligible-preset catalogue remains an activation prerequisite. |
| Policy inventory would have to permit a runner dispatch | The exact browser grammar/class tests above passed. Keep its refusal matrix and exercise the new internal mode separately. |
| Terminal mint request path remains unguarded as once noted in evaluator 430 | `git grep -nF 'terminalEpisodeLock(' -- cmd/ailang-worldd/session_new.go` returned :427 in shared `provisionEpisode`. The direct revoke paths still need the independent exclusion measured above. |

### 1.5 Size and CI race cost

| Milestone | Increment | Cumulative +/− | Filtered | Unfiltered product | STATUS/DELETED/hashes |
|---|---:|---:|---:|---:|---:|
| M1 | 20 | +17/−3 | 20 | 20 | 1/0/1 |
| M2 | 78 | +92/−6 | 98 | 98 | 4/0/4 |
| M3 | 480 | +557/−21 | 578 | 578 | 17/0/17 |

Size applies `git diff -M --numstat 9610f5de HEAD -- cmd host docs app` at each controller commit. The archive equivalent uses per-file `git diff --no-index -M --numstat`, base-to-stage, and the same file scope. No paths are removed or renamed. Unfiltered product counts use all archive paths. The two plan files are controller planning records excluded from snapshots and counted separately in the JSON unfiltered branch totals. The normal ceiling is 1,400.

Focused and whole-package race receipts are pending.

### 1.6 Portability

Compare physical locations with `filepath.EvalSymlinks` on each side. Native and Linux-target vet run on the touched packages. The re-exec process uses the test binary and an 8 s context, a pipe readiness sentinel, stdin EOF and bounded Wait. Its normal race runtime exits without the race runtime's fixed atexit delay. No real store, daemon, agent, launchctl, worktree outside a fixture, or loopback is used by any new named test.

## 2. How the executor applies a sub-milestone

Run inside the sprint worktree. Apply each patch once, in order. The base is an ancestor even if the controller has already committed an earlier snapshot.

````bash
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter439-row243-v4s5f
STATE=$HOME/.ailang/state/world-iter439
S=$STATE/plan-scratch
T=$STATE/executor
B=9610f5de0921fe1c8eeef939928ffe639bd2259b
PLAN=$W/design_docs/planned/w-world-address-bar-verbs-v4s5f-sprint-plan.md
M=M1  # Set M2 or M3 for the next unapplied milestone.
mkdir -p "$T"
unset SNAP_PARENT SNAP_STATUS_LIST
git -C "$W" merge-base --is-ancestor "$B" HEAD
if [ "$M" = M1 ]; then
  git -C "$W" diff --quiet "$B" -- cmd host docs app
  test -z "$(git -C "$W" ls-files --others --exclude-standard -- cmd host docs app)"
fi
cp "$S"/{bounded.py,failset.py,verify-scratch.py,callers.sh,callers-base.txt,probe.sh,operator_gc_after_test.go} "$T/"
for name in accept.sh gates.sh snapshot.sh; do
  awk -v n="$name" '$0 == "````bash " n {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$name"
  test -s "$T/$name"
  /bin/bash -n "$T/$name"
done
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
(cd "$W" && /usr/bin/patch -p1 -E < "$T/$M.patch")
/bin/bash "$T/accept.sh" "$W" "$M" fast
/bin/bash "$T/snapshot.sh" "$W" "$M"
````

## 3. Milestones

### M1 — Judge 438 GC kind coverage

Use the real credential migration and the real store reservation / observation API. Keep one effect kind per matching episode. A reserved effect is pinned even with a terminal parent. An observed effect is pinned while its parent is running/uncertain. `TestEpisodeHasUnresolvedOperatorSteps/{files_reserved,files_observed_running}` are the S9 rows, and `{mint_reserved,mint_observed_uncertain}` the S10 rows. Each pair independently detects its edit. Each lookup validates the returned boolean and unchanged store bytes. Keep the worktree observations and the independent cleanup episode as controls.

### M2 — Attended mint resolver read/lock/read

Pass the existing acquiring function into `reconcileRunnerMint`, as the dispatch seam passes readers to existing recovery commands. Production dispatch supplies `openRunnerLock` with flags 0 and retains the existing defer and readback. `TestOperatorMintRecoveryLockedRead/different_episode` interleaves a complete redacted receipt refit to a different physical episode and proves its private token matches. The resolver refuses the episode conflict under the first episode's lock. `/expired_on_locked_read` crosses the expiry check with the same valid token and a newly expired credential. Both expect the exact refusal, check a second real receipt read, compare credentials/fixture bytes against the state after interleaving, and re-acquire the released fd. The mutations target the exact old :38 and :41 edits.

### M3 — Dormant start selection, lock-scoped admission, and stable signed inputs

`TestNextOperatorStart` proves SQL kind/state, strict expiry, creation edge and stable oldest/tie order independently of the CLI. `TestOperatorStartCandidateProcess` contends with a true terminal lock holder against a symlinked credentials spelling. A later visit sees the same fd pointer and inode, gets the actual locked row, and proves the writer is free. `TestOperatorStartCandidateRefusals` rotates the row generation, advances a fake clock exactly to signature expiry, or supplies unsafe/missing evidence. Each refusal expects the class, verifies complete action readback and compares credential/worktree/file bytes. After any acquisition require successful independent reacquisition.

`TestOperatorTerminalCreateLockMissingParent` models a directory disappearing after canonicalization and makes the existing `flags == 0` distinction reachable. Its read-only legacy arm accepts absence with an unmodified namespace. `TestOperatorWaitReadBudget` receives budgets of 250, 150 and 50 ms using an injected clock and returns false after the deadline. The production wait keeps its existing bounded reader. M3 preserves the liveness update while retaining the catalogue generation on an identical payload and readiness. `TestOperatorCatalogueAdmission` verifies stability, changed bytes, changed readiness and restart invalidation against a real store. M3 does not invoke the visit from the runner.

`TestOperatorStartFacts` checks the exact signed `launch_binary_sha256`. Its `read_disappeared` arm refuses an executable that vanishes after the live resolution returns, with evidence unchanged. `TestOperatorStartChanged/launch_binary_bytes` replaces executable bytes at the signed path and checks the binding refusal and state. The effect-free reconstruction signs and rechecks the digest. An activating runner still has to match that digest against the attended configuration and prove the confined execution path.

### 3.1 Design §5.3 branch assignment

| Reviewed row | Staged proof and activation gate |
|---|---|
| Changed preset, profile, base, task, daemon and executable | Existing `TestOperatorStartChanged` verifies reconstructed facts. M3 adds the locked durable row read and signed agent bytes. Before activation bind those bytes to the attended allowlist, prove the supported confined backend and trusted Git command/auth configuration, and recheck the current pair, typed proof and fresh expiry at each new effect. |
| Worktree collision/window | Existing reservation helper proves initial absence and independent Git readback. At activation, borrow the M3 fd across the worktree and the Git child. Never create while a terminal or GC operation can use that path. |
| Atomic mint and rollback/window | M2 pins locked reread and existing retained token behavior. Before first production call, pin the observed-created worktree dependency, claimed proof, signed grants/TTL, current session and bounded handling of an uncommitted mint reservation and a committed mint with lost token. |
| File recovery | Existing publication helper installs with no replacement and independently reads back the files. Before activation provide the same held fd and a runnable remedy in every reserved/observed gap. Recheck current session equivalence, expiry and revocation separately. |
| Launch and later browser paths | S6 proves one confined run with exact binary, run and birth identities plus positive session validity and no fallback. Browser-independent state precedes S7/S8. |

### 3.2 Design §5.5 and §2.2 process assignment

The first real operation takes the credentials runner kernel lock and canonical episode lock before CAS. Stop admitting new effects when a drain starts. Keep the fd until the bounded active Git or file operation is quiescent and its exact readback is recorded, or the saved state has a runnable attended remedy. Treat a dead process, elapsed heartbeat or returned PID separately from kernel quiescence. The upgrade operator reviews the two matched build identities before stopping, followed by paired startup, positive health and compatible pair rollback. See Appendix C for the four process checks and the measured caller.

## 4. Mutation table

Full mutation runs are pending. The baseline set recorded S9/S10 surviving whole store and both E-RES edits surviving their named tests.

### 4.1 Whole-package runs for the controller

Re-record the outside-sandbox base with the exact pinned tools before classifying any whole mutation. Save raw TEXT and all parser controls before comparing. Do not use row 271 failures from a different commit as a substitute base. Each `--- FAIL: <named test>` must be present for a kill. A daemon/CLI bind denial invalidates the run, and a Go timeout truncates it. The scripts print the independent parser control and reject an unattributed nonzero result.

````bash
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter439-row243-v4s5f
S=$HOME/.ailang/state/world-iter439/plan-scratch
/bin/bash "$S/gates.sh" "$S/base-tree" "$S/controller-base" base "" "$W"
/bin/bash "$S/gates.sh" "$S/st-M3" "$S/controller-head" M3 "$S/controller-base" "$W"
/bin/bash "$S/whole-mutants.sh" "$S/st-M3" "$S/controller-base"
````

These commands write only their fresh scratch directories. The final outside-sandbox ledger compares against `pkg-cmd-ailang-worldd.base-fails` and `pkg-host-store.base-fails` under `controller-base`. A leg exit alone does not determine its package result. Section 8 lists the known comparators.

## 5. Files

Final cumulative product numstat by path.

```text
4	0	cmd/ailang-worldd/operator_episode_lock.go
3	2	cmd/ailang-worldd/operator_mint.go
71	0	cmd/ailang-worldd/operator_mint_test.go
2	2	cmd/ailang-worldd/operator_runner.go
3	3	cmd/ailang-worldd/operator_service.go
47	0	cmd/ailang-worldd/operator_start_candidate.go
229	0	cmd/ailang-worldd/operator_start_candidate_test.go
5	0	cmd/ailang-worldd/operator_start_input.go
29	6	cmd/ailang-worldd/operator_start_input_test.go
62	0	cmd/ailang-worldd/operator_wait_pins_test.go
1	0	host/store/credentials_test.go
12	0	host/store/operator_actions.go
17	3	host/store/operator_episode_gc_test.go
4	3	host/store/operator_runner.go
2	0	host/store/operator_runner_test.go
8	2	host/store/operator_s5e4_test.go
58	0	host/store/operator_start_candidate_test.go
```

## 6. Sprint gates

The pristine gate and compiled stage checks are recorded in `$S/gate-base.log` and `$S/check-stages-final.log`. Final whole-stage gate results are pending.

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
  if grep -F 'panic: test timed out' "$L/$tag.raw"; then
    echo 'INSTRUMENT FAILURE: truncated package'
    return 2
  fi
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
  for path in host/store/operator*.go cmd/ailang-worldd/operator*.go; do
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
else leg accept env S5F_PIN_HOME="$ORIGINAL" /bin/bash "$S/accept.sh" "$T" "$M" fast; fi
if [ "$M" = base ]; then
  leg gc-probe env HOME="$ORIGINAL" /bin/bash "$S/probe.sh" "$T" operator_gc_after_test.go
elif [ "$M" = M3 ]; then
  leg gc-probe env HOME="$ORIGINAL" /bin/bash "$S/probe.sh" "$T" operator_gc_after_test.go
else leg gc-probe env HOME="$ORIGINAL" /bin/bash "$S/probe.sh" "$T" operator_gc_after_test.go; fi
if [ "$M" = base ] || [ "$M" = "${S5F_FINAL:-M3}" ]; then
  leg verify_ail bounded 510 python3 "$S/verify-scratch.py" "$T" "$L" verify_ail ""
  if [ -n "$REF" ]; then
    leg verify_go record verify_go bounded 1500 python3 "$S/verify-scratch.py" "$T" "$L" verify_go "$REF"
  else skip verify_go 'read-only committed-blob Git reference missing'; fi
else
  skip verify_ail 'AILANG tree byte-identical; base and M3 carry the verifier'
  skip verify_go 'pristine and final broad gate'
fi
leg callers /bin/bash "$S/callers.sh" "$T" "$L/callers.census" "$S/callers-base.txt" "$M"
for pkg in host/store host/daemon cmd/ailang-worldd; do
  name=${pkg//\//-}
  if { [ "$M" = M1 ] && [ "$pkg" != host/store ]; } || { [ "$M" != base ] && [ "$pkg" = host/daemon ]; }; then
    skip "pkg-$name" 'package files unchanged'
    skip "race-$name" 'package files unchanged'
    continue
  fi
  leg "pkg-$name" record "pkg-$name" bounded 1230 go test -count=1 -timeout 20m "./$pkg"
  leg "race-$name" record "race-$name" bounded 1230 go test -race -count=1 -timeout 20m "./$pkg"
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

The GC probe replays real fixture publication on fresh scratch copies. Both reserved and published-under-running arms report `non-terminal operator step` in the live GC core and CLI. The disposed published arm reappears as a clean and merged candidate at 73 h in dry run. The probe compares Git and store evidence. The M3 census uses the same positive controls and retains zero non-test calls to `MintOperatorStart`. Appendix C records the separate prototype's one call.

## 7. Executor directive notes

````bash accept.sh
#!/bin/bash
set -euo pipefail
W=$1; M=$2; PART=$3
case "$M" in M1|M2|M3) ;; *) exit 2;; esac
test "$PART" = fast
S=$(cd "$(dirname "$0")" && pwd -P)
PIN=${S5F_PIN_HOME:-$HOME}
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
  if grep -F 'panic: test timed out' "$log"; then
    echo 'INSTRUMENT FAILURE: truncated acceptance'
    exit 2
  fi
  if grep -Ei '(bind|listen tcp).*(: operation not permitted|: permission denied)' "$log"; then
    echo 'UNINFORMATIVE UNDER SANDBOX'
    exit 77
  fi
  test "$rc" -eq 0
  while IFS= read -r name; do
    grep -Fq -- "--- PASS: $name (" "$log"
  done < <(printf '%s\n' "$names" | tr '|' '\n')
}
pins='TestEpisodeHasUnresolvedOperatorSteps|TestEpisodeOperatorStepLookupRefusals'
if [ "$M" = M3 ]; then
  pins="$pins|TestNextOperatorStart|TestOperatorCatalogueAdmission|TestProductionContextRoots|TestCredentialsMethodsRequireADeadline|TestStoreSurfaceTakesContext"
fi
run host/store "$pins"
if [ "$M" != M1 ]; then
  run cmd/ailang-worldd 'TestOperatorMintRecovery|TestOperatorMintRecoveryLockedRead'
fi
if [ "$M" = M3 ]; then
  run cmd/ailang-worldd 'TestOperatorStartFacts|TestOperatorStartCandidateProcess|TestOperatorStartCandidateRefusals|TestOperatorTerminalCreateLockMissingParent|TestOperatorWaitReadBudget|TestOperatorServiceStartup'
  digest_log=$(mktemp "$run_dir/executable-digest.XXXXXX")
  python3 "$S/bounded.py" 180 go test -v -count=1 -timeout 2m ./cmd/ailang-worldd \
    -run '^TestOperatorStartChanged$/^launch_binary_bytes$' 2>&1 | tee "$digest_log"
  grep -Fq -- '--- PASS: TestOperatorStartChanged/launch_binary_bytes (' "$digest_log"
fi
printf 'ACCEPT_OK %s %s\n' "$M" "$PART"
````

````bash snapshot.sh
#!/bin/bash
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=$1
M=$2
case "$M" in M1|M2|M3) ;; *) exit 2;; esac
B=9610f5de0921fe1c8eeef939928ffe639bd2259b
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter439/snap}/$M
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
    $2 != "design_docs/planned/w-world-address-bar-verbs-v4s5f-sprint-plan.md" &&
    $2 != "design_docs/planned/sprint_w-world-address-bar-verbs-v4s5f.json"
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

Run snapshot after each `ACCEPT_OK`, with `B=9610f5de0921fe1c8eeef939928ffe639bd2259b` in the fenced script. Snapshot parent defaults to `$HOME/.ailang/state/world-iter439/snap`. Each milestone contains a cumulative `STATUS`, `DELETED` and `MANIFEST.sha256`. The two plan files are excluded. The controller records the plan and JSON before milestone commits so the commit script sees no untracked planning files. Use the commit test below from the controller, where shell variables have the same values as §2. The commit script runs gofmt, build, vet and the named test before it records numstat.

````bash
M=M1 # M2 and M3 use the same named entry with their own milestone argument.
printf 'test(operator): V4 S5f %s — pin pre-dispatch preparation\n' "$M" > "$T/$M.msg"
/bin/bash "$STATE/commit_ms.sh" "$M" "$T/$M.msg" \
  "AILANG_BIN='$HOME/.pinned-ailang/ailang' WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN='$HOME/.pinned-ailang-tools/v0.54.15/ailang' /bin/bash '$T/accept.sh' '$W' '$M' fast"
````

The fresh archive replay and snapshot rehearsal are pending.

## 8. Risks and pre-existing test failures

Row 271's CLI comparator includes `TestCLIAgainstSeToolsDaemon`, `TestCLINoGenesisIsNoWorldHead`, eight `TestDemoSeed*` and `TestMarketingStillRoutes`. Daemon includes `TestSeToolsDescriptorLiveExactEight`, `TestSeToolsMCPVisibilityFollowsGrants`, and `TestSeToolsA2AEndToEnd`. Separately identify the known load-sensitive `TestWorldContextAndStatusline`, `host/broker TestExecProtectedCheckTimeout`, `host/daemon TestSeToolsRunBenchmarkTasksEndToEnd`, `host/verifygate TestMissionConfigDefaultFlowIntegration`, and `host/broker TestSeedRestartIgnoresASurvivingOldCopy`. Preserve new failure differences for examination. Row 267 has CLI race observations of 262.9 s and 477.6/469.0 s against a 480 s bound, and daemon 355.8 s with slow runs at 480 s. Use the focused CLI budget and whole-package receipts separately.

The aged GC read-to-remove window is still present, and controller bind-denied tests stay uninformative. Command and service lifecycle activation remain parked until the evidence in Appendix C is available.

## 9. Acceptance criteria

1. Three incremental patch fences apply using `/usr/bin/patch -p1 -E` and reproduce their cumulative stage trees and snapshot receipts.
2. M1 and M2 kill the exact carried S9, S10 and E-RES pair. M3 kills the Z budget and absence edits, proves that a real terminal holder excludes selection before handing out the same fd, and pins the executable digest and stable catalogue generation.
3. Run native build/vet, Linux vet, the named deadline and boundary inventories, pinned AILANG gate and GC probe. Compare all attributable package failures to the recorded base and follow §4.1 for the sandbox-limited legs.
4. Census keeps the staged caller count at zero. The first caller can progress only together with the four process checks and the missing activation prerequisites below.

## Appendix C — S5f2 activation and optional pins

The prototype builds and vets in `$S/proto-S5f2`. Activation is blocked on the prerequisite and process matrix below. The staged M3 has no production caller.

| Carried process result | Definition and remaining proof |
|---|---|
| `TestOperatorFirstManagedEpisodeProcess` | `design_docs/planned/w-world-address-bar-verbs-v4s5b-sprint-plan.md:728` names the first managed claim against a terminal operation. S5c :593, S5d :808, S5e :646 and S5e4 :582 re-carry the requirement. M3's `TestOperatorStartCandidateProcess` proves the pre-CAS FD seam with real processes. The first caller test must pause a real runner before CAS, prove that a terminal holder leaves the row confirmed with zero reservations, then prove one claim and only the intended receipt/file materialization while terminal and GC contenders refuse. |
| `TestOperatorOldBuildTerminalDrain` | S5b :728 names retiring the old terminal build and S4b plan :682 requires process-to-process quiescence. Use a real old compatible terminal or re-exec with the old operation semantics. Pin exclusion of a writer which bypasses the new path until retirement is proved. Demonstrate unchanged binary, record and effect rows if quiescence cannot be proved. |
| `TestOperatorCurrentOwnerProcessRecovery` | S5b :728 defines recovery while a competing process owns the episode. With a real holder and a fresh/stale heartbeat split, observe the first owner or refuse with no changes. After the child exits, verify canonical descriptor exclusion and a fresh proof/credential reread before any narrowly allowed recovery. Record exact failure classes for current-session mismatch, expiry and revocation. No stale lease based owner transfer. |
| `TestOperatorServiceUpgrade/active_child` | S5b :728 and S5d :74, :190 describe the paired drain test. Reviewed V4 :292 gives active child, swap and failed pair restart arms. Exercise the attended service controller seam with a real bounded re-exec child, finish/uncertain with its executable remedy, pair identity and rollback. Refuse on unprovable exit. Target fresh/old signed pair and one compatible rollback. |
| `J-OWN-UID` | S4b plan :670, :685 and S5b plan :730 define dropping the actual `Stat_t.Uid` comparison. Judge 429 `evaluator-r1.md:25,51` reports the first edit failed compilation and a corrected non-root subset survived. The old mutable test ledger and before/after source named by S5b are unavailable in this rig's iteration-429 state. Preserve real alternate-owner testing for both `cmd/operator_service.go:71` and `cmd/operator_runner.go:30` until the original amended substitutions and fixture are recovered. |
| `Z-ABSENT-ANY` and `Z-WAIT-READBUDGET` | S5 plan :702 and S5b :730 re-carry these. Exact edits in `design_docs/verification/world-iter430/evaluator-r1.md:17`: remove `flags == 0` from `absent` in `operator_episode_lock.go:36` and replace `min(remaining, time.Second)` with `time.Second` in `operator_runner.go:115` at this head. Both have M3 named tests. |
| Terminal `session revoke` | S5c :593 and S5d :808 carry explicit terminal exclusion. Cover by-ID and by-episode routes under a real same-episode holder after read-only lookup, re-read on the locked side, and preserve unrelated sessions. Terminal worktree, start, task, token and credential bytes must stay unchanged on refusal. |

A new activation unit must account for expired and revoked keys at each new effect; prove independently current session ID/metadata, session expiry and revocation at post-mint and recovery; sign the exact approved executable identity and confinement path; and preserve the §1.3 GC exclusion. For a failed or interrupted claim with no reservations, prove quiescence and no effect before releasing it. For observed worktree only, reserved mint without commit, observed mint without private files, partially published files, and fully verified prep awaiting S6, prove a permissible reconcile/disposition with no blind retries. Recover from disk receipts under the lock, with the exact current owner CAS. Only a later launch sprint reserves and executes the agent process.

Optional N04/N12/N15 retain historical S5e4 evidence of 82 changed lines and +4.399 s focused CLI race. Current-head remeasurement is pending.

## Appendix P — compiled patches

### Patch M1

````diff M1
diff --git a/host/store/operator_episode_gc_test.go b/host/store/operator_episode_gc_test.go
index 6767b8b2..5b25f18a 100644
--- a/host/store/operator_episode_gc_test.go
+++ b/host/store/operator_episode_gc_test.go
@@ -13,8 +13,8 @@ import (
 
 func TestEpisodeHasUnresolvedOperatorSteps(t *testing.T) {
 	for _, tc := range []struct {
-		name, phase, state string
-		want               bool
+		name, kind, phase, state string
+		want                     bool
 	}{
 		{name: "none"},
 		{name: "operation_only", state: "running"},
@@ -27,11 +27,25 @@ func TestEpisodeHasUnresolvedOperatorSteps(t *testing.T) {
 		{name: "observed_expired", phase: "observed", state: "expired"},
 		{name: "observed_cancelled", phase: "observed", state: "cancelled"},
 		{name: "terminal_with_reservation", phase: "reserved", state: "failed", want: true},
+		{name: "files_reserved", kind: "files", phase: "reserved", state: "failed", want: true},
+		{name: "files_observed_running", kind: "files", phase: "observed", state: "running", want: true},
+		{name: "mint_reserved", kind: "mint", phase: "reserved", state: "failed", want: true},
+		{name: "mint_observed_uncertain", kind: "mint", phase: "observed", state: "uncertain", want: true},
 	} {
 		t.Run(tc.name, func(t *testing.T) {
 			f := newOperatorPinFixture(t, false)
 			if tc.phase != "" {
-				key := f.reserveWorktree(t)
+				var key OperatorStepKey
+				if tc.kind == "" {
+					key = f.reserveWorktree(t)
+				} else {
+					// Isolate the lookup by giving this episode exactly one effect kind.
+					attempt := f.in.Attempt
+					key = OperatorStepKey{attempt.CAS, attempt.Runner.Nonce, tc.kind, attempt.CAS.ID, 1}
+					if err := f.c.ReserveOperatorStep(f.ctx, key, pinNow); err != nil {
+						t.Fatal(err)
+					}
+				}
 				if tc.phase == "observed" {
 					if err := f.c.ObserveOperatorStep(f.ctx, key, []byte("created")); err != nil {
 						t.Fatal(err)
````

### Patch M2

````diff M2
diff --git a/cmd/ailang-worldd/operator_mint.go b/cmd/ailang-worldd/operator_mint.go
index 26bd4aa2..91df289a 100644
--- a/cmd/ailang-worldd/operator_mint.go
+++ b/cmd/ailang-worldd/operator_mint.go
@@ -14,7 +14,8 @@ import (
 	"github.com/sunholo-data/ailang-world/host/verbs"
 )
 
-func (e serviceEnv) reconcileRunnerMint(id string, out, stderr io.Writer) int {
+func (e serviceEnv) reconcileRunnerMint(id string, out, stderr io.Writer,
+	acquire func(string, int) (*os.File, error)) int {
 	refuse := func(err error) int {
 		return serviceMessage(stderr, "mint recovery refused", err.Error(), "ailang-worldd service runner status", 1)
 	}
@@ -26,7 +27,7 @@ func (e serviceEnv) reconcileRunnerMint(id string, out, stderr io.Writer) int {
 	if err != nil {
 		return refuse(err)
 	}
-	lock, err := openRunnerLock(operatorEpisodeLockPath(c.Credentials, r.Episode), 0)
+	lock, err := acquire(operatorEpisodeLockPath(c.Credentials, r.Episode), 0)
 	if err != nil {
 		return refuse(errors.New("episode recovery lock held or unavailable"))
 	}
diff --git a/cmd/ailang-worldd/operator_mint_test.go b/cmd/ailang-worldd/operator_mint_test.go
index 8d976043..78afab64 100644
--- a/cmd/ailang-worldd/operator_mint_test.go
+++ b/cmd/ailang-worldd/operator_mint_test.go
@@ -303,6 +303,77 @@ func TestOperatorMintRecovery(t *testing.T) {
 	}
 }
 
+// Change real ledger evidence between the first read and the real lock.
+func TestOperatorMintRecoveryLockedRead(t *testing.T) {
+	for _, tc := range []struct{ name, want string }{
+		{"different_episode", store.ErrOperatorConflict.Error()},
+		{"expired_on_locked_read", "matching token expired"},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			r := newRecoveryFixture(t)
+			id := r.receipt.CAS.ID
+			token := r.path
+			if tc.name == "different_episode" {
+				token = filepath.Join(filepath.Dir(r.path), "retargeted.session")
+			}
+			if err := writeTokenFile(token, r.token); err != nil {
+				t.Fatal(err)
+			}
+			var out, stderr bytes.Buffer
+			var before string
+			var reviewed store.OperatorMintReceipt
+			files := provisionFiles(t, filepath.Dir(r.path))
+			worktree := provisionFiles(t, r.worktree)
+			acquire := func(path string, flags int) (*os.File, error) {
+				t.Helper()
+				if path != operatorEpisodeLockPath(r.config.Credentials, r.receipt.Episode) || flags != 0 {
+					t.Fatalf("recovery lock = %q flags %d", path, flags)
+				}
+				first, err := readRunnerMint(r.config.Credentials, id)
+				if err != nil || first != r.receipt {
+					t.Fatalf("initial receipt differs: %v", err)
+				}
+				if tc.name == "different_episode" {
+					r.rewriteAction(t, func(a *approval.OperatorAction) {
+						a.Episode = "retargeted"
+						a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
+						for i := range a.Resolution {
+							if a.Resolution[i].Name == "mint_token_path" {
+								a.Resolution[i].Value = token
+							}
+						}
+					})
+					r.exec(t, "UPDATE operator_actions SET episode_id='retargeted' WHERE id=?", id)
+					r.exec(t, "UPDATE episode_starts SET episode='retargeted' WHERE episode=?", first.Episode)
+					r.exec(t, "UPDATE session_credentials SET episode_id='retargeted' WHERE credential_id=?", first.CredentialID)
+				} else {
+					r.exec(t, "UPDATE session_credentials SET expires_at=? WHERE credential_id=?", r.env.terminal.now(), first.CredentialID)
+				}
+				before = serviceFileHash(r.config.Credentials)
+				reviewed, err = readRunnerMint(r.config.Credentials, id)
+				if err != nil || reviewed == first {
+					t.Fatalf("interleaving did not leave a fresh readable receipt: %v", err)
+				}
+				if present, err := inspectRunnerMintToken(r.config, reviewed); err != nil || !present {
+					t.Fatalf("interleaved token evidence not reachable: %t %v", present, err)
+				}
+				return openRunnerLock(path, flags)
+			}
+			code := r.env.reconcileRunnerMint(id, &out, &stderr, acquire)
+			if code != exitUsage || !strings.Contains(stderr.String(), "why: "+tc.want+"\n") || out.Len() != 0 {
+				t.Fatalf("locked read %s: %d stdout=%q stderr=%q", tc.name, code, &out, &stderr)
+			}
+			after, err := readRunnerMint(r.config.Credentials, id)
+			if err != nil || after != reviewed || serviceFileHash(r.config.Credentials) != before ||
+				!reflect.DeepEqual(files, provisionFiles(t, filepath.Dir(r.path))) ||
+				!reflect.DeepEqual(worktree, provisionFiles(t, r.worktree)) {
+				t.Fatalf("refusal changed store, worktree or files: %v", err)
+			}
+			requireProvisionLockReleased(t, r.config.Credentials, r.receipt.Episode)
+		})
+	}
+}
+
 func TestOperatorMintTokenReadback(t *testing.T) {
 	for _, name := range []string{"uid", "mode", "link", "symlink", "symlink_valid", "directory", "fifo", "bytes", "length", "short"} {
 		t.Run(name, func(t *testing.T) {
diff --git a/cmd/ailang-worldd/operator_service.go b/cmd/ailang-worldd/operator_service.go
index b47600c4..867f55ee 100644
--- a/cmd/ailang-worldd/operator_service.go
+++ b/cmd/ailang-worldd/operator_service.go
@@ -284,7 +284,7 @@ func runServiceRunner(args []string, out, errOut io.Writer, e serviceEnv) int {
 		return refuse(errors.New("unknown runner operation"))
 	}
 	if verb == "reconcile-mint" {
-		return e.reconcileRunnerMint(args[1], out, errOut)
+		return e.reconcileRunnerMint(args[1], out, errOut, openRunnerLock)
 	}
 	if verb == "reconcile-worktree" {
 		return e.reconcileRunnerWorktree(args[1], out, errOut, readRunnerWorktree)
````

### Patch M3

````diff M3
diff --git a/cmd/ailang-worldd/operator_episode_lock.go b/cmd/ailang-worldd/operator_episode_lock.go
index 758f7cbe..2ffbba1c 100644
--- a/cmd/ailang-worldd/operator_episode_lock.go
+++ b/cmd/ailang-worldd/operator_episode_lock.go
@@ -24,6 +24,10 @@ func openTerminalEpisodeLock(creds, episode string, held *os.File, flags int) (*
 	if err != nil {
 		return nil, err
 	}
+	return openCanonicalTerminalEpisodeLock(canonical, episode, held, flags)
+}
+
+func openCanonicalTerminalEpisodeLock(canonical, episode string, held *os.File, flags int) (*os.File, error) {
 	path := operatorEpisodeLockPath(canonical, episode)
 	if held != nil {
 		if held.Name() != path {
diff --git a/cmd/ailang-worldd/operator_runner.go b/cmd/ailang-worldd/operator_runner.go
index 07cc7691..ac031e5f 100644
--- a/cmd/ailang-worldd/operator_runner.go
+++ b/cmd/ailang-worldd/operator_runner.go
@@ -105,14 +105,14 @@ func runnerHeartbeat(path string, budget time.Duration) (store.OperatorHeartbeat
 	return store.ReadOperatorHeartbeat(ctx, path)
 }
 
-func (e serviceEnv) waitForRunner(c runnerConfig, previous int64) bool {
+func (e serviceEnv) waitForRunner(c runnerConfig, previous int64, read func(string, time.Duration) (store.OperatorHeartbeat, error)) bool {
 	deadline := e.runnerNow().Add(e.wait)
 	for {
 		remaining := deadline.Sub(e.runnerNow())
 		if remaining <= 0 {
 			return false
 		}
-		h, err := runnerHeartbeat(c.Credentials, min(remaining, time.Second))
+		h, err := read(c.Credentials, min(remaining, time.Second))
 		now := e.runnerNow()
 		if !now.Before(deadline) {
 			return false
diff --git a/cmd/ailang-worldd/operator_service.go b/cmd/ailang-worldd/operator_service.go
index 867f55ee..36d0041f 100644
--- a/cmd/ailang-worldd/operator_service.go
+++ b/cmd/ailang-worldd/operator_service.go
@@ -450,7 +450,7 @@ func (e serviceEnv) changeRunner(verb string, args []string) error {
 
 			return errors.New("could not bootstrap the companion")
 		}
-		if !e.waitForRunner(c, previous.Owner.Generation) {
+		if !e.waitForRunner(c, previous.Owner.Generation, runnerHeartbeat) {
 			if err := e.controller.Bootout(e.domain, runnerLabel); err != nil && !errors.Is(err, errServiceNotLoaded) {
 				return errors.New("runner readiness timed out and companion could not be stopped")
 			}
@@ -458,7 +458,7 @@ func (e serviceEnv) changeRunner(verb string, args []string) error {
 		}
 	case err != nil:
 		return errors.New("could not inspect the companion")
-	case !e.waitForRunner(c, 0):
+	case !e.waitForRunner(c, 0, runnerHeartbeat):
 		return errors.New("loaded companion heartbeat unavailable")
 	}
 	return nil
diff --git a/cmd/ailang-worldd/operator_start_candidate.go b/cmd/ailang-worldd/operator_start_candidate.go
new file mode 100644
index 00000000..944267bb
--- /dev/null
+++ b/cmd/ailang-worldd/operator_start_candidate.go
@@ -0,0 +1,47 @@
+package main
+
+import (
+	"context"
+	"fmt"
+	"os"
+	"reflect"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+// withOperatorStartCandidate owns the episode descriptor for the entire
+// visit. The runner may wire a writer here only once the process and recovery
+// gates are met. The visit must use a fresh clock and the locked row for CAS.
+func withOperatorStartCandidate(ctx context.Context, credentials string, now func() int64,
+	acquire func(string) (*os.File, error), visit func(*os.File, store.OperatorActionRow) error) error {
+	canonical, err := store.CanonicalPath(credentials)
+	if err != nil {
+		return err
+	}
+	db, err := store.OpenCredentialsReadOnly(runnerDB(canonical))
+	if err != nil {
+		return err
+	}
+	defer db.Close()
+	initial, err := db.NextOperatorStart(ctx, now())
+	if err != nil {
+		return err
+	}
+	if !verbs.EpisodeIDPattern.MatchString(initial.Episode) {
+		return store.ErrOperatorConflict
+	}
+	lock, err := acquire(operatorEpisodeLockPath(canonical, initial.Episode))
+	if err != nil {
+		return fmt.Errorf("%w: episode lock unavailable: %w", store.ErrOperatorConflict, err)
+	}
+	defer lock.Close()
+	current, err := db.NextOperatorStart(ctx, now())
+	if err != nil {
+		return err
+	}
+	if !reflect.DeepEqual(current, initial) {
+		return store.ErrOperatorConflict
+	}
+	return visit(lock, current)
+}
diff --git a/cmd/ailang-worldd/operator_start_candidate_test.go b/cmd/ailang-worldd/operator_start_candidate_test.go
new file mode 100644
index 00000000..4aac89c0
--- /dev/null
+++ b/cmd/ailang-worldd/operator_start_candidate_test.go
@@ -0,0 +1,229 @@
+package main
+
+import (
+	"bufio"
+	"context"
+	"errors"
+	"fmt"
+	"io"
+	"os"
+	"os/exec"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"syscall"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func TestOperatorStartCandidateProcess(t *testing.T) {
+	const holder = "WORLD_TEST_START_CANDIDATE_TERMINAL"
+	if path := os.Getenv(holder); path != "" {
+		fd, err := terminalEpisodeLock(path, "ep", nil)
+		if err != nil {
+			t.Fatal(err)
+		}
+		defer fd.Close()
+		fmt.Println("holding canonical terminal lock")
+		if _, err := io.Copy(io.Discard, os.Stdin); err != nil {
+			t.Fatal(err)
+		}
+		return
+	}
+	creds := filepath.Join(t.TempDir(), "credentials.db")
+	canonical, err := store.CanonicalPath(creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	lock := claimFixtureEpisode(t, creds, "ep")
+	lock.Close()
+	credSQL(t, creds, "UPDATE operator_actions SET state='confirmed', owner=NULL")
+	aliasPath := filepath.Join(t.TempDir(), "alias.db")
+	if err := os.Symlink(creds, aliasPath); err != nil {
+		t.Fatal(err)
+	}
+	bin, err := os.Executable()
+	if err != nil {
+		t.Fatal(err)
+	}
+	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
+	defer cancel()
+	cmd := exec.CommandContext(ctx, bin, "-test.run=^TestOperatorStartCandidateProcess$")
+	cmd.Env = append(os.Environ(), holder+"="+creds, "GORACE="+strings.TrimSpace(os.Getenv("GORACE")+" atexit_sleep_ms=0"))
+	cmd.WaitDelay = time.Second
+	stdin, err := cmd.StdinPipe()
+	if err != nil {
+		t.Fatal(err)
+	}
+	stdout, err := cmd.StdoutPipe()
+	if err != nil {
+		t.Fatal(err)
+	}
+	if err := cmd.Start(); err != nil {
+		t.Fatal(err)
+	}
+	waited := false
+	t.Cleanup(func() {
+		if !waited {
+			cancel()
+			stdin.Close()
+			cmd.Wait()
+		}
+	})
+	reader := bufio.NewReader(stdout)
+	if line, err := reader.ReadString('\n'); err != nil || line != "holding canonical terminal lock\n" {
+		t.Fatalf("child barrier: %q %v", line, err)
+	}
+	before := serviceFileHash(creds)
+	worktree := t.TempDir()
+	files := t.TempDir()
+	treeBefore, filesBefore := provisionFiles(t, worktree), provisionFiles(t, files)
+	expectUnchanged := func() {
+		t.Helper()
+		row := candidateRow(t, creds, 1012)
+		if serviceFileHash(creds) != before || row.State != "confirmed" || row.Owner != "" ||
+			!reflect.DeepEqual(provisionFiles(t, worktree), treeBefore) || !reflect.DeepEqual(provisionFiles(t, files), filesBefore) {
+			t.Fatal("candidate probe changed ownership, files or worktree")
+		}
+	}
+	now := func() int64 { return 1012 }
+	err = withOperatorStartCandidate(ctx, aliasPath, now, runnerKernelLock, func(*os.File, store.OperatorActionRow) error {
+		t.Fatal("visit ran while the real terminal process held the lock")
+		return nil
+	})
+	if !errors.Is(err, store.ErrOperatorConflict) || !errors.Is(err, syscall.EWOULDBLOCK) {
+		t.Fatalf("terminal holder returned wrong refusal: %v", err)
+	}
+	expectUnchanged()
+	stdin.Close()
+	if _, err := io.Copy(io.Discard, reader); err != nil {
+		t.Fatal(err)
+	}
+	err = cmd.Wait()
+	waited = true
+	if err != nil {
+		t.Fatalf("terminal release: %v", err)
+	}
+	handled := errors.New("test visit stopped")
+	var acquired *os.File
+	err = withOperatorStartCandidate(ctx, aliasPath, now, func(path string) (*os.File, error) {
+		var err error
+		acquired, err = runnerKernelLock(path)
+		return acquired, err
+	}, func(fd *os.File, row store.OperatorActionRow) error {
+		physical, err := filepath.EvalSymlinks(fd.Name())
+		want, wantErr := filepath.EvalSymlinks(operatorEpisodeLockPath(canonical, "ep"))
+		if err != nil || wantErr != nil || physical != want || fd != acquired || row.ID != fixtureActionID {
+			t.Fatalf("visit did not receive exact locked candidate: %v %v", err, wantErr)
+		}
+		if err := runnerWrite(canonical, func(*store.Credentials) error { return nil }); err != nil {
+			t.Fatalf("candidate held the SQLite writer across the visit: %v", err)
+		}
+		if contender, err := runnerKernelLock(fd.Name()); !errors.Is(err, syscall.EWOULDBLOCK) {
+			if contender != nil {
+				contender.Close()
+			}
+			t.Fatalf("visit lost kernel exclusion: %v", err)
+		}
+		return handled
+	})
+	if !errors.Is(err, handled) {
+		t.Fatalf("visit outcome: %v", err)
+	}
+	expectUnchanged()
+	requireProvisionLockReleased(t, canonical, "ep")
+}
+
+func candidateRow(t *testing.T, path string, now int64) store.OperatorActionRow {
+	t.Helper()
+	db, err := store.OpenCredentialsReadOnly(runnerDB(path))
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer db.Close()
+	row, err := db.OperatorAction(boundedTestContext(t), fixtureActionID, now)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return row
+}
+
+func TestOperatorStartCandidateRefusals(t *testing.T) {
+	for _, tc := range []struct {
+		name, before, interleave string
+		advance                  bool
+		want                     error
+	}{
+		{"changed_generation", "", "UPDATE operator_actions SET generation=generation+1", false, store.ErrOperatorConflict},
+		{"expired_under_lock", "", "", true, store.ErrOperatorMissing},
+		{"empty", "UPDATE operator_actions SET state='requested', approved_by=NULL, artifact=NULL", "", false, store.ErrOperatorMissing},
+		{"invalid_episode", "UPDATE operator_actions SET episode_id='ep/invalid'", "", false, store.ErrOperatorConflict},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			creds := filepath.Join(t.TempDir(), "credentials.db")
+			held := claimFixtureEpisode(t, creds, "ep")
+			held.Close()
+			canonical, err := store.CanonicalPath(creds)
+			if err != nil {
+				t.Fatal(err)
+			}
+			credSQL(t, creds, "UPDATE operator_actions SET state='confirmed', owner=NULL")
+			if tc.before != "" {
+				credSQL(t, creds, tc.before)
+			}
+			before, hash := candidateRow(t, creds, 1012), serviceFileHash(creds)
+			worktree, files := t.TempDir(), t.TempDir()
+			treeBefore, filesBefore := provisionFiles(t, worktree), provisionFiles(t, files)
+			now := int64(1012)
+			var acquired bool
+			err = withOperatorStartCandidate(boundedTestContext(t), creds, func() int64 { return now },
+				func(path string) (*os.File, error) {
+					acquired = true
+					if tc.advance {
+						now = 1600
+					}
+					if tc.interleave != "" {
+						credSQL(t, creds, tc.interleave)
+						before, hash = candidateRow(t, creds, 1012), serviceFileHash(creds)
+					}
+					return runnerKernelLock(path)
+				}, func(*os.File, store.OperatorActionRow) error {
+					t.Fatal("refused candidate reached its visit")
+					return nil
+				})
+			if !errors.Is(err, tc.want) || acquired != (tc.interleave != "" || tc.advance) ||
+				!reflect.DeepEqual(candidateRow(t, creds, 1012), before) || serviceFileHash(creds) != hash ||
+				!reflect.DeepEqual(provisionFiles(t, worktree), treeBefore) || !reflect.DeepEqual(provisionFiles(t, files), filesBefore) {
+				t.Fatalf("refusal changed evidence or returned wrong class: %v", err)
+			}
+			requireProvisionLockReleased(t, canonical, "ep")
+		})
+	}
+	t.Run("missing_parent", func(t *testing.T) {
+		root := t.TempDir()
+		before := provisionFiles(t, root)
+		missing := filepath.Join(root, "absent", "credentials.db")
+		err := withOperatorStartCandidate(boundedTestContext(t), missing, func() int64 { return 1012 },
+			func(string) (*os.File, error) {
+				t.Fatal("canonical path failure reached the lock")
+				return nil, nil
+			}, nil)
+		if !errors.Is(err, os.ErrNotExist) || !reflect.DeepEqual(provisionFiles(t, root), before) {
+			t.Fatalf("missing parent: %v", err)
+		}
+	})
+	t.Run("missing_store", func(t *testing.T) {
+		root := t.TempDir()
+		before := provisionFiles(t, root)
+		err := withOperatorStartCandidate(boundedTestContext(t), filepath.Join(root, "credentials.db"), func() int64 { return 1012 },
+			func(string) (*os.File, error) {
+				t.Fatal("missing store reached the lock")
+				return nil, nil
+			}, nil)
+		if err == nil || !strings.Contains(err.Error(), "unable to open database file") || !reflect.DeepEqual(provisionFiles(t, root), before) {
+			t.Fatalf("missing store: %v", err)
+		}
+	})
+}
diff --git a/cmd/ailang-worldd/operator_start_input.go b/cmd/ailang-worldd/operator_start_input.go
index 95e77cd3..f83eb66a 100644
--- a/cmd/ailang-worldd/operator_start_input.go
+++ b/cmd/ailang-worldd/operator_start_input.go
@@ -67,7 +67,12 @@ func operatorStartFacts(ctx context.Context, signed approval.OperatorAction, att
 	if err != nil {
 		return store.OperatorMintInput{}, nil, err
 	}
+	digest := serviceFileHash(read.launchBinary)
+	if digest == "" {
+		return store.OperatorMintInput{}, nil, &approval.BindingError{Field: "launch_binary_sha256"}
+	}
 	facts := []approval.OperatorResolved{{Name: "launch_binary", Value: read.launchBinary},
+		{Name: "launch_binary_sha256", Value: digest},
 		{Name: "launch_wall", Value: read.launchWall.String()}, {Name: "mcp_path", Value: mcp}}
 	facts = append(facts, mint...)
 	facts = append(facts, approval.OperatorResolved{Name: "served_address", Value: read.preset.Daemon.Addr},
diff --git a/cmd/ailang-worldd/operator_start_input_test.go b/cmd/ailang-worldd/operator_start_input_test.go
index 3726ca54..24d18f04 100644
--- a/cmd/ailang-worldd/operator_start_input_test.go
+++ b/cmd/ailang-worldd/operator_start_input_test.go
@@ -83,7 +83,7 @@ func TestOperatorStartFacts(t *testing.T) {
 		names = append(names, fact.Name)
 		value[fact.Name] = fact.Value
 	}
-	want := []string{"launch_binary", "launch_wall", "mcp_path", "mint_catalogue", "mint_grants", "mint_profile", "mint_project",
+	want := []string{"launch_binary", "launch_binary_sha256", "launch_wall", "mcp_path", "mint_catalogue", "mint_grants", "mint_profile", "mint_project",
 		"mint_start", "mint_token_path", "served_address", "served_build", "served_interpreter", "worktree_path"}
 	if !reflect.DeepEqual(names, want) {
 		t.Fatalf("reviewed fact names %v, want %v", names, want)
@@ -98,7 +98,7 @@ func TestOperatorStartFacts(t *testing.T) {
 	if err != nil {
 		t.Fatal(err)
 	}
-	for name, v := range map[string]string{"launch_binary": binary, "launch_wall": "1h0m0s", "mcp_path": "",
+	for name, v := range map[string]string{"launch_binary": binary, "launch_binary_sha256": serviceHash(readStartFixtureFile(t, binary)), "launch_wall": "1h0m0s", "mcp_path": "",
 		"mint_catalogue": "3", "mint_profile": f.profile.Digest, "mint_project": "proj",
 		"mint_token_path": tokenPath, "served_address": defaultEpisodeAddr, "served_build": string(build),
 		"served_interpreter": h.ToolRef, "worktree_path": filepath.Join(f.r.f.ws, "ep")} {
@@ -126,9 +126,27 @@ func TestOperatorStartFacts(t *testing.T) {
 		t.Fatalf("effective grants %+v, want the preset's %+v", in.Grants, grants)
 	}
 	mint, err := in.Resolution()
-	if err != nil || !reflect.DeepEqual(mint, r.action.Resolution[3:9]) {
+	if err != nil || !reflect.DeepEqual(mint, r.action.Resolution[4:10]) {
 		t.Fatalf("mint facts differ from the effective input: %v", err)
 	}
+	t.Run("read_disappeared", func(t *testing.T) {
+		before := observeProvision(t, f.r, r.creds)
+		read := func(ctx context.Context, name, presets, profiles string, in presetStartInput, clock func() int64) (*presetStartRead, error) {
+			live, err := resolvePresetStart(ctx, name, presets, profiles, in, clock)
+			if err == nil {
+				if err := os.Remove(live.launchBinary); err != nil {
+					t.Fatal(err)
+				}
+			}
+			return live, err
+		}
+		_, _, err := operatorStartFacts(presetTestContext(t), r.action, r.at, f.dir, f.r.profiles, f.r.f.ws, r.hosts, presetClock, read)
+		var binding *approval.BindingError
+		if !errors.As(err, &binding) || binding.Field != "launch_binary_sha256" ||
+			!reflect.DeepEqual(observeProvision(t, f.r, r.creds), before) {
+			t.Fatalf("vanished executable changed state or returned wrong class: %v", err)
+		}
+	})
 }
 
 // TestOperatorStartChanged changes one live input after review at a time.
@@ -165,7 +183,7 @@ func TestOperatorStartChanged(t *testing.T) {
 			episodeSelfBuild = func() daemon.BuildIdentity { return other }
 		}},
 		{name: "interpreter", field: "served_interpreter", change: func(t *testing.T, r *operatorStartRig) {
-			r.action.Resolution[11].Value = "sha256:" + serviceHash([]byte("other"))
+			r.action.Resolution[12].Value = "sha256:" + serviceHash([]byte("other"))
 		}},
 		{name: "launch_binary", field: "launch_binary", change: func(t *testing.T, r *operatorStartRig) {
 			other := startReadFakeAgent(t, r.f)
@@ -173,11 +191,16 @@ func TestOperatorStartChanged(t *testing.T) {
 				t.Fatal("fixture reused the agent directory")
 			}
 		}},
+		{name: "launch_binary_bytes", field: "launch_binary_sha256", change: func(t *testing.T, r *operatorStartRig) {
+			if err := os.WriteFile(filepath.Join(r.bin, "codex"), []byte("#!/bin/sh\nexit 71\n"), 0700); err != nil {
+				t.Fatal(err)
+			}
+		}},
 		{name: "mcp_path", field: "mcp_path", change: func(t *testing.T, r *operatorStartRig) {
-			r.action.Resolution[2].Value = filepath.Join(r.f.r.sessions, "ep.mcp.json")
+			r.action.Resolution[3].Value = filepath.Join(r.f.r.sessions, "ep.mcp.json")
 		}},
 		{name: "worktree_path", field: "worktree_path", change: func(t *testing.T, r *operatorStartRig) {
-			r.action.Resolution[12].Value += "-other"
+			r.action.Resolution[13].Value += "-other"
 		}},
 		{name: "extra_fact", field: "resolution", change: func(t *testing.T, r *operatorStartRig) {
 			r.action.Resolution = append(r.action.Resolution, approval.OperatorResolved{Name: "zz_extra", Value: "x"})
diff --git a/cmd/ailang-worldd/operator_wait_pins_test.go b/cmd/ailang-worldd/operator_wait_pins_test.go
new file mode 100644
index 00000000..a9cc01e1
--- /dev/null
+++ b/cmd/ailang-worldd/operator_wait_pins_test.go
@@ -0,0 +1,62 @@
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
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func TestOperatorTerminalCreateLockMissingParent(t *testing.T) {
+	root := t.TempDir()
+	parent := filepath.Join(root, "gone")
+	if err := os.Mkdir(parent, 0700); err != nil {
+		t.Fatal(err)
+	}
+	canonical, err := store.CanonicalPath(filepath.Join(parent, "credentials.db"))
+	if err != nil {
+		t.Fatal(err)
+	}
+	// This is the window after canonicalization but before lock open.
+	if err := os.Remove(parent); err != nil {
+		t.Fatal(err)
+	}
+	before := provisionFiles(t, root)
+	lock, err := openCanonicalTerminalEpisodeLock(canonical, "ep", nil, os.O_CREATE)
+	if lock != nil || !errors.Is(err, os.ErrNotExist) || !reflect.DeepEqual(provisionFiles(t, root), before) {
+		t.Fatalf("creating lock accepted the missing parent or changed it: %v %v", lock, err)
+	}
+	lock, err = openCanonicalTerminalEpisodeLock(canonical, "ep", nil, 0)
+	if lock != nil || err != nil || !reflect.DeepEqual(provisionFiles(t, root), before) {
+		t.Fatalf("read-only legacy absence control: %v %v", lock, err)
+	}
+}
+
+func TestOperatorWaitReadBudget(t *testing.T) {
+	dir := t.TempDir()
+	before := provisionFiles(t, dir)
+	c := runnerConfig{Credentials: filepath.Join(dir, "credentials.db")}
+	clock := time.Unix(1012, 0)
+	env := serviceEnv{wait: 250 * time.Millisecond,
+		runnerNow:   func() time.Time { return clock },
+		runnerSleep: func(step time.Duration) { clock = clock.Add(step) },
+	}
+	var budgets []time.Duration
+	read := func(path string, budget time.Duration) (store.OperatorHeartbeat, error) {
+		t.Helper()
+		budgets = append(budgets, budget)
+		if path != c.Credentials || len(budgets) > 3 {
+			t.Fatalf("unexpected heartbeat read: %s %v", path, budgets)
+		}
+		return store.OperatorHeartbeat{}, store.ErrOperatorRunner
+	}
+	want := []time.Duration{250 * time.Millisecond, 150 * time.Millisecond, 50 * time.Millisecond}
+	if env.waitForRunner(c, 0, read) || !reflect.DeepEqual(budgets, want) ||
+		!reflect.DeepEqual(provisionFiles(t, dir), before) {
+		t.Fatalf("wait did not refuse with the remaining read budget: %v", budgets)
+	}
+}
diff --git a/host/store/credentials_test.go b/host/store/credentials_test.go
index a4a55695..9656fdf7 100644
--- a/host/store/credentials_test.go
+++ b/host/store/credentials_test.go
@@ -308,6 +308,7 @@ func TestCredentialsMethodsRequireADeadline(t *testing.T) {
 	_, errs["PendingVerbRequests"] = c.PendingVerbRequests(ctx, 1000, 0)
 	errs["ConsumeVerbRequest"] = c.ConsumeVerbRequest(ctx, "x", "y", ApproverUse{}, "", nil, "", 1000)
 	_, errs["OperatorAction"] = c.OperatorAction(ctx, "x", 1000)
+	_, errs["NextOperatorStart"] = c.NextOperatorStart(ctx, 1000)
 	errs["ReserveOperatorStep"] = c.ReserveOperatorStep(ctx, OperatorStepKey{}, 1000)
 	errs["ObserveOperatorStep"] = c.ObserveOperatorStep(ctx, OperatorStepKey{}, nil)
 	_, errs["OperatorStep"] = c.OperatorStep(ctx, OperatorStepKey{})
diff --git a/host/store/operator_actions.go b/host/store/operator_actions.go
index 867a98ad..d7777c13 100644
--- a/host/store/operator_actions.go
+++ b/host/store/operator_actions.go
@@ -68,6 +68,18 @@ func (c *Credentials) OperatorAction(ctx context.Context, id string, now int64)
 	return row, err
 }
 
+// NextOperatorStart reads the oldest eligible start without assigning ownership.
+// The caller must take the canonical episode lock and compare the row again
+// before claiming it.
+func (c *Credentials) NextOperatorStart(ctx context.Context, now int64) (OperatorActionRow, error) {
+	if err := c.operatorAccess(ctx); err != nil {
+		return OperatorActionRow{}, err
+	}
+	return scanOperatorAction(c.db.QueryRowContext(ctx, `SELECT `+operatorColumns+` FROM operator_actions
+        WHERE kind = 'start' AND state = 'confirmed' AND created_at <= ? AND expires_at > ?
+        ORDER BY created_at, id LIMIT 1`, now, now))
+}
+
 // RequestOperatorAction records inert work or permanently binds another tap's
 // operation ID to the existing action. The active-episode index arbitrates starts
 // and finishes together.
diff --git a/host/store/operator_runner.go b/host/store/operator_runner.go
index a09d76bc..98dcbe77 100644
--- a/host/store/operator_runner.go
+++ b/host/store/operator_runner.go
@@ -59,10 +59,11 @@ func (c *Credentials) PutOperatorHeartbeat(ctx context.Context, h OperatorHeartb
 			h.Owner.Nonce, h.Protocol, h.Build, h.Readiness, h.SeenAt, h.Catalogue).Scan(&generation)
 		return generation, err
 	}
+	// Refresh liveness without superseding a review of the same catalogue.
 	err := c.db.QueryRowContext(ctx, `UPDATE operator_runner_state SET readiness=?, seen_at=?,
-		catalogue_generation=catalogue_generation+1, catalogue=?
-		WHERE singleton=1 AND owner=? AND generation=? AND protocol=? AND build=? AND seen_at<=?
-		RETURNING generation`, h.Readiness, h.SeenAt, h.Catalogue, h.Owner.Nonce,
+		catalogue_generation=catalogue_generation+CASE WHEN readiness IS ? AND catalogue IS ? THEN 0 ELSE 1 END,
+        catalogue=? WHERE singleton=1 AND owner=? AND generation=? AND protocol=? AND build=? AND seen_at<=?
+		RETURNING generation`, h.Readiness, h.SeenAt, h.Readiness, h.Catalogue, h.Catalogue, h.Owner.Nonce,
 		h.Owner.Generation, h.Protocol, h.Build, h.SeenAt).Scan(&generation)
 	if errors.Is(err, sql.ErrNoRows) {
 		err = ErrOperatorRunner
diff --git a/host/store/operator_runner_test.go b/host/store/operator_runner_test.go
index dcc20f0c..6a3aeb8d 100644
--- a/host/store/operator_runner_test.go
+++ b/host/store/operator_runner_test.go
@@ -61,6 +61,8 @@ func TestOperatorCatalogueAdmission(t *testing.T) {
 	a.Owner.Generation, a.SeenAt, a.Catalogue = 1, 11, []byte("two")
 	publish(a, 1, 2) // same-second tick is a valid refresh
 	a.SeenAt++
+	publish(a, 1, 2) // liveness only leaves the approved catalogue current
+	a.Readiness = "none"
 	second := publish(a, 1, 3)
 	if string(first.Catalogue) != "one" {
 		t.Fatal("published generation changed in-place")
diff --git a/host/store/operator_s5e4_test.go b/host/store/operator_s5e4_test.go
index 6261c8e2..ee7cd058 100644
--- a/host/store/operator_s5e4_test.go
+++ b/host/store/operator_s5e4_test.go
@@ -64,6 +64,11 @@ func newOperatorPinFixture(t *testing.T, allowance bool) operatorPinFixture {
 }
 
 func signOperatorPinAction(t *testing.T, c *Credentials, row OperatorActionRow, a approval.OperatorAction) *approval.OperatorMintProof {
+	t.Helper()
+	return signOperatorPinActionAt(t, c, row, a, pinNow)
+}
+
+func signOperatorPinActionAt(t *testing.T, c *Credentials, row OperatorActionRow, a approval.OperatorAction, now int64) *approval.OperatorMintProof {
 	t.Helper()
 	body, err := approval.CanonicalOperatorAction(a)
 	if err != nil {
@@ -91,11 +96,12 @@ func signOperatorPinAction(t *testing.T, c *Credentials, row OperatorActionRow,
 	}
 	artifact := signer.Sign(approval.OperatorActionNamespace, message)
 	signed, err := approval.CheckOperatorAction(approval.OperatorCheckInput{Action: a, Challenge: challenge,
-		Expiry: row.ExpiresAt, Now: pinNow, Artifact: artifact, Format: ap.Format, Approver: ap})
+		Expiry: row.ExpiresAt, Now: now, Artifact: artifact, Format: ap.Format, Approver: ap})
 	if err != nil {
 		t.Fatal(err)
 	}
-	operatorExec(t, c, "UPDATE operator_actions SET action=?, digest=?, artifact=?, approved_by=? WHERE id=?",
+	operatorExec(t, c, `UPDATE operator_actions SET action=?, digest=?, artifact=?, approved_by=?,
+        state=CASE WHEN state='requested' THEN 'confirmed' ELSE state END WHERE id=?`,
 		body, digest[:], artifact, signed.ApprovedBy(), row.ID)
 	return signed.MintProof()
 }
diff --git a/host/store/operator_start_candidate_test.go b/host/store/operator_start_candidate_test.go
new file mode 100644
index 00000000..f24699ed
--- /dev/null
+++ b/host/store/operator_start_candidate_test.go
@@ -0,0 +1,58 @@
+package store
+
+import (
+	"bytes"
+	"errors"
+	"fmt"
+	"os"
+	"reflect"
+	"testing"
+)
+
+func TestNextOperatorStart(t *testing.T) {
+	c, ctx := openTestCredentials(t), testCtx(t)
+	add := func(n int, kind string, created, expires int64, confirm bool) OperatorActionRow {
+		t.Helper()
+		req := operatorRequest(t, c, n, fmt.Sprintf("ep-%d", n), kind)
+		req.CreatedAt, req.ExpiresAt = created, expires
+		// Insert in reverse tie-break order to exercise the explicit sort.
+		row := putOperator(t, c, req, created)
+		if confirm {
+			signOperatorPinActionAt(t, c, row, req.Action, created)
+		}
+		return getOperator(t, c, row.ID, created)
+	}
+	path := c.lock.dbPath
+	check := func(now int64, want OperatorActionRow) {
+		t.Helper()
+		before, err := os.ReadFile(path)
+		if err != nil {
+			t.Fatal(err)
+		}
+		got, err := c.NextOperatorStart(ctx, now)
+		if want.ID == "" {
+			if !errors.Is(err, ErrOperatorMissing) || got.ID != "" {
+				t.Fatalf("no eligible start at %d: %q %v", now, got.ID, err)
+			}
+		} else if err != nil || !reflect.DeepEqual(got, want) {
+			t.Fatalf("next start at %d: %q %v, want %q", now, got.ID, err, want.ID)
+		}
+		assertGCLedgerBytes(t, path, before)
+		if want.ID != "" && !bytes.Equal(getOperator(t, c, want.ID, now).Digest, want.Digest) {
+			t.Fatal("read changed the action digest")
+		}
+	}
+	add(1, "finish", 1000, 1600, true)
+	add(2, "start", 1000, 1600, false)
+	add(3, "start", 1000, 1100, true)
+	check(1100, OperatorActionRow{})
+	add(4, "start", 1200, 1600, true)
+	check(1100, OperatorActionRow{})
+	first := add(6, "start", 1100, 1600, true)
+	winner := add(5, "start", 1100, 1600, true)
+	check(1100, winner)
+	check(1200, winner)
+	operatorExec(t, c, "UPDATE operator_actions SET state='claimed', owner='runner' WHERE id=?", winner.ID)
+	check(1100, first)
+	check(1600, OperatorActionRow{})
+}
````
