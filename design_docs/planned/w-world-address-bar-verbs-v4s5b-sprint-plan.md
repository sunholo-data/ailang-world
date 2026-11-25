# Sprint plan — w-world-address-bar-verbs-v4, sprint S5b (iteration 432, row 243)

> **Scope.** M1 closes the direct terminal `session mint` exclusion gap. M2 compiles the atomic start, session and mint receipt under a distinct sealed proof, and ships the attended remedy for a committed orphan in the same patch. A matching token can be inspected and preserved. The transaction has no production caller. Worktree reservation, file publication and the paired-process consumer stay in Appendix C.
>
> **Authority.** Reviewed design r1 §§2.2, 3.1–3.3, 4.2, 5.3, 5.5, 6–8. Quorum r2 PROCEED 3/3, no absent seats. D-WORLD-119=A and D-WORLD-120=A, attended `82f05004`. S1 #606, S2 #609, S2b #611, S3 #643, S4 #651, S4b #654 and S5 #657 (`ab520385`, sonnet PASS 86, zero blocking) are landed. Base `f8d8e3ad27b8442e1cd62674052b0392f03089a6`.

## 0. Binding constraints

- The controller commits the two plan records before applying M1 and M2, in order, with `/usr/bin/patch -p1 -E`. The controller owns Git writes. The planner used only archive stages under the iteration scratch root for product changes.
- Use disposable real migrations, repository signers, an injected clock/terminal and fake supervisor. Product tests write within test fixture roots. Real credentials, LaunchAgents, World data and launchctl are observation-only. The resolver's live approval stays attended under AGENTS.md.
- Stop before any reachable claim, reserve, observe or operator mint. The recovery command has a separate exact CAS and can address an already observed mint only. All pair, drain and provenance prerequisites travel with activation.
- Each new failure assertion must read back the affected store and artifacts. Hold and acquire via independent descriptors. Keep context roots in existing bounded helpers. No new test in either `host/daemon` or `host/verifygate`.

## 1. Pre-flight

### 1.1 Baseline at f8d8e3ad

Before creating a stage, `git archive f8d8e3ad | tar -x -C "$S/base"` exported the pinned tree. The prior §6 gate with only the base SHA replaced ran first using `GATES_GIT_REF="$W" GATES_SANDBOX=1 bash "$S/gates-first.sh" "$S/base" "$S/gate-base" base "$B"`. It ended in 152 s. `verify_go` was red at base because the archive lacks Git metadata (`tracked_binary_gate.sh: FATAL: git ls-tree HEAD failed`, after 24 fixture passes). Both CLI whole-package legs and QUICKSTART §9 had inherited loopback denial, and launchctl observation was sandbox-uninformative. Base and final-format sweeps record the archive's `verify_go` leg as SKIPPED and leave it to the controller's Git checkout. The final-format pristine baseline ran with `GATES_GIT_REF="$W" GATES_SANDBOX=1 GATES_SKIP_VERIFY_GO=1 bash "$S/gates.sh" "$S/base" "$S/gate-base-final" base "$B"`.

````text
LEG observation-before rc=77 0s UNINFORMATIVE UNDER SANDBOX
LEG paths rc=0 0s
TOUCHED_PACKAGES: 
LEG touched rc=0 1s
LEG status rc=0 0s
LEG gofmt rc=0 0s
LEG build rc=0 1s
LEG vet rc=0 1s
LEG fence rc=0 1s
LEG verify_ail rc=0 79s
LEG pkg-cmd-ailang-worldd rc=1 3s UNINFORMATIVE UNDER SANDBOX
LEG race-cmd-ailang-worldd rc=1 21s UNINFORMATIVE UNDER SANDBOX
LEG pkg-host-approval rc=0 3s
LEG race-host-approval rc=0 7s
LEG pkg-host-authority rc=0 4s
LEG race-host-authority rc=0 19s
LEG pkg-host-store rc=0 14s
LEG race-host-store rc=0 43s
LEG inv-host-approval rc=0 1s
LEG inv-host-approval-approvaltest rc=0 1s
LEG inv-host-authority rc=0 1s
LEG inv-host-boundary rc=0 3s
LEG inv-host-daemon rc=0 2s
LEG inv-host-verbs rc=0 1s
LEG inv-host-broker rc=0 2s
LEG inv-cmd-ailang-worldd rc=0 2s
LEG inv-quickstart-section9 rc=1 2s UNINFORMATIVE UNDER SANDBOX
LEG inv-host-store rc=0 1s
LEG routes rc=0 0s
LEG scope rc=0 1s
LEG callers rc=0 1s
LEG verify_go SKIPPED 0s (GATES_SKIP_VERIFY_GO=1)
LEG observation-after rc=77 0s UNINFORMATIVE UNDER SANDBOX
GATES_DONE 217s FAILS: none UNINFORMATIVE: observation-before pkg-cmd-ailang-worldd race-cmd-ailang-worldd inv-quickstart-section9 observation-after
````

### 1.2 Current inputs and coupling inventory

The measurements use `F=cmd/ailang-worldd`, `W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter432-row243-v4s5b`, and `S=$HOME/.ailang/state/world-iter432/plan-scratch`. `B` is the full base SHA. Anchor commands run in `W` unless a stage is given.

| ID | Command | Observed output and consequence |
|---|---|---|
| C1 | `git rev-parse HEAD; git diff --name-status ab520385 HEAD; git log --oneline ab520385..HEAD` | `f8d8e3ad27b8442e1cd62674052b0392f03089a6`. Seven record paths, all under `design_docs/`, from #658. Product diff empty. |
| C2 | `rg -n 'resolvePresetStart\(|runPresetEpisodeStart\(|readPresetStartTask\(' "$F"/episode{,_preset}.go; rg -n 'prepareMint\(|provisionEpisode\(|terminalEpisodeLock\(' "$F" -g '*.go' -g '!**/*_test.go'` | M2 read seam at `episode_preset.go:79,288,304,328`; terminal caller `episode.go:272`. `prepareMint` and shared provision at `session_new.go:237,405`. Provision callers `approve.go:398`, `episode.go:377`, `episode_preset.go:365`, `episode_resume.go:215`, `session_new.go:660`. Lock helper `operator_episode_lock.go:18`; calls `episode.go:364`, `episode_finish.go:128`, `episode_launch.go:456`, `episode_preset.go:334`, `episode_resume.go:175`, `session_new.go:424`. |
| C3 | `rg -n 'runSessionMint\(|authority\.Mint\(|terminalEpisodeLock\(|openCredentialsForSession\(|os.WriteFile\(' "$F"/session.go; rg -n 'func .*Mint\b\|func .*MintStarted\b\|func .*MintForRequest\b' host/authority/mint.go` | `runSessionMint:199`, writer `:286`, `authority.Mint:299`, optional file write `:306`, zero episode lock hits. `authority/mint.go:30,50,112` keeps public Mint, start and A3 semantics. |
| C4 | `rg -n 'MintSessionWithStart\|mintSession\|insertSessionRow\|ResolveSession\|RevokeSession' host/store/credentials.go; rg -n 'operatorCASRow\|ConfirmOperatorAction\|ClaimOperatorAction\|operatorChangedOne' host/store/operator_cas.go; rg -n 'RecheckClaimedOperatorAction\|OperatorConfirmation' host/authority/judge_operator.go` | `credentials.go:469,473,526,545,580`; `episode_starts.go:70` reader. `operator_cas.go:32,42,119,157`. Existing action generation and digest CAS, and sealed raw law transport `judge_operator.go:120,129`. Existing store has no step reserve/observe API. `operator_schema.sql` has a reserved/observed step with NULL receipt while reserved and private owner. |
| C5 | `rg -n 'worktree\.(Add\|AddExisting\|Remove)\|afterToken\|writeTokenFile\|writeMCPConfig' "$F"/{session_new,episode,episode_preset}.go; rg -n '^func (Add\|Output)' host/worktree/worktree.go` | `session_new.go:510,512,549,552,577,674`, `episode.go:552`. Preset paths and task freeze `episode_preset.go:195–219`; writer closures `:339–362`. `worktree.go:63,83`. The S5 terminal writers are downstream of the read seam. |
| C6 | `git grep -n -E '\.(ClaimOperatorAction\|ReserveOperatorStep\|ObserveOperatorStep\|RecheckClaimedOperatorAction\|RecheckVerb)\(' -- cmd host \| grep -vF '_test.go:'` | Existing positive control `cmd/ailang-worldd/episode_finish_surface.go:239` invokes `authority.RecheckVerb`. Zero production claim, reserve or observe sites. Final §6 adds positive controls for inner mint and attended abort and requires zero outer mint caller. |
| C7 | `rg -n 'func Test.*(DependencyAllowlist\|ProductionGoSurface\|ProductionContextRoots\|WriterSites\|Quickstart\|VerbTable\|EveryVerb\|OperatorFlagParity\|PluginCommands)' cmd host -g '*test.go'; git grep -n LOCK_NB -- cmd/ailang-worldd host/store` | Inventories below. Six active `Flock(...LOCK_NB)` lines plus one writer-lock comment. No new primitive. |
| C8 | `rg -n 'runnerHelp\|case "cancel"\|func .*installedRunner\|func runnerWrite\|func openRunnerLock\|func .*confirm' "$F"/{operator_service,operator_runner,service}.go` | `operator_service.go:24,222,275`; `operator_runner.go:25,162`; `service.go:171`. Reuse exact private installation validation, bounded runner store open, terminal `y` and the existing lock primitive. Reconcile is a separate operation from unstarted cancel. |
| C9 | `"$HOME/.pinned-ailang/ailang" version; go version; rg -n 'runs-on: ubuntu' .github/workflows -g '*.yml'` | Pinned AILANG v0.41.0, `24ee1088776e21cd06a3781ed18e77f40be06db3`. go1.26.6 darwin/arm64. Linux CI selects ubuntu-24.04 in `.github/workflows/ci.yml:121,446`. |

| Inventory | Anchor | Disposition |
|---|---|---|
| CLI commands, help, flags | `help_test.go:21`, `verbs_parity_test.go:66`, `operator_parity_test.go:65`, `episode_launch_test.go:449`, `plugin_help_test.go:52`, `cmd/world-publish/plugin_help_test.go:58` | Nested runner help grows by `reconcile-mint ACTION-ID`. Top-level verb, flag and plugin enumerations and their expectations stay the same. |
| QUICKSTART | `quickstart_service_test.go:15`, `quickstart_se_test.go:56`, `quickstart_verbs_test.go:15`, `quickstart_finish_surface_test.go:17` | New resolver example and exact sidecar sentence assertion. Existing sections and goldens stay unchanged. The runbook remedy passes through real CLI parsing and an injected tty. |
| Writer sites and locks | `cli_credential_writers_test.go:18`, `host/daemon/surface_approve_post_guard_test.go:132`, `operator_runner.go:35`, `operator_service.go:251`, `service.go:158,662`, `host/store/writer_lock_unix.go:41,86` | Same CLI/daemon open allowlist. Resolver uses the existing `runnerWrite`, and the inventory counts 26 read-only opens (25 at base). One bounded episode lock on each added path; same `LOCK_EX|LOCK_NB` and descriptor ownership. |
| Context roots and packages | `host/store/context_roots_test.go:129,234`, `host/store/deadline_guard_test.go:43`, `host/store/credentials_test.go:277` | Two production Go files within the already-scanned packages. No new context roots. Add three credential deadline controls. |
| Dependencies and imports | `host/daemon/daemon_test.go:899`, `host/boundary/allowlist_world_test.go:758`, `host/approval/imports_test.go:27`, `host/broker/allowlist_test.go:89` | Existing package, per-group imports and route fences stand. `host/store` already imports `host/approval`, and `host/approval` stays stdlib-only. The mint imports the existing `host/hashref` dependency of `host/store`. |

| Stage anchor measurement | Selected output |
|---|---|
| `rg -n 'func (.* )?(MintProof\|Bind\|CommitOperatorMint\|ReadOperatorMint\|AbortOperatorMint\|MintOperatorStart\|insertEpisodeStart\|reconcileRunnerMint\|readMintToken)\b' "$S/st-M2"/{cmd,host} -g '*.go' -g '!**/*_test.go'` | `host/approval/operator_signed.go:194,203`, `host/store/operator_mint.go:94,242,290`, `host/authority/mint.go:204`, `host/store/credentials.go:516`, `cmd/ailang-worldd/operator_mint.go:17,104`. Confirm actual anchors in `stage-anchors.txt`. |
| `grep -F -e 'root census' -e 'Go surface' "$S/gate-scratch/inv-host-store.txt"; cat "$S/gate-scratch/callers.txt"` | 54 context roots, 326 production files, six explicit reproducers. One legacy read-seam caller, five shared provision callers, seven flock lines including the comment, zero outer mint callers, one inner mint caller in the inert authority entry, one attended abort caller. |

### 1.3 Decisions

- M1 locks admitted episode IDs after opening the tty and before displaying the mint prompt. Legacy opaque IDs remain supported because operator action admission requires the shared, narrower grammar. The episode lock fences token output and writer close. An attended decline can leave a stable private inode.
- M2 takes an independently reconstructed action and effective mint inputs through the existing runner signature judge. The store binds the exact artifact, signed envelope, action/runner/catalogue/step identity and current enrollment with a fresh clock. The sealed mint proof only accepts a runner recheck. The same transaction writes the start, credential and observed step. A fixture reserves the step because the production consumer is parked.
- The transaction projects the archived start fields (with transaction-owned timestamp), ordered semantic grants, profile binding, token path and catalogue generation into signed resolution facts. M2 earns this mint boundary. The CLI-to-catalogue reconstruction and worktree/file admission stay in S5c.
- The recovery command re-reads evidence under the same canonical sidecar, examines the exact private token with bounded fd readback and preserves the owner for a matching live token. An absent token beneath an established private directory is eligible for attended revoke only before any later reservation. Typed `y` gates a second readback and independent transactional owner/step compare. Ambiguity preserves the owner and identifies the blocker.
- The new reachable behavior is terminal mint exclusion and the attended resolver. The inert mint seam has zero production effect dispatch. The caller fence covers S5b and is an activation dependency for S5c.

### 1.4 Premises found false or incomplete at HEAD

| ID | Measurement and consequence |
|---|---|
| F1 | C2 and C3 reproduce the lock-free mint. A legacy raw episode is accepted by `session mint` but cannot become an operator action under `host/approval/operator_action.go:133–135` and `host/approval/request.go:93`. M1 tests both an opaque ID and exact durable owner exclusion after kernel lock release. |
| F2 | C8 shows cancel is only `CancelUnstartedOperatorAction` (`host/store/operator_runner.go:74–83`). Re-ran the prior `TestS5ResolverProbe` in an isolated base export with `AILANG_BIN=$HOME/.pinned-ailang/ailang go test -count=1 -v ./cmd/ailang-worldd -run '^TestS5ResolverProbe$'`: status rc=0, attended cancel rc=1 conflict, uncertain owner unchanged. M2 resolves the committed orphan only, with effect-specific evidence. Production mint remains unreachable. |
| F3 | `rg -n 'catalogue_generation' host/store/operator_runner.go` shows initial generation 1 and increment on later heartbeat. The store rechecks the independently signed catalogue generation. Advance after review makes that mint stale. |
| F4 | `rg --files host/store | grep -F operator_steps.go` found no file. Reapplying the S3b scratch bank measured 132 lines (§Appendix C). No production reserve/observe API or paired-process drain landed with S5. |
| F5 | §1.1 observed pristine archive failure taking the whole `verify_go` script before Git metadata exists. The sandboxes deny loopback and launchctl independently. §6 preserves the controller commands and classifies leg results explicitly. |
| F6 | Reproductions of the three S5 judge variants and their current disposition are in §4 and Appendix C. The judge's whole-CLI results are external measurements from iteration 431. |
| F7 | `rg -n 'Digest: "sha256:"\|func serviceHash\|SHA256: serviceHash' host/broker/exec_profile.go "$F"/{service,episode_preset}.go` yields `host/broker/exec_profile.go:158`, `service.go:104`, `episode_preset.go:220`. S5 supplies tagged profile and task SHA-256 values. The mint uses `host/hashref/hashref.go:127,171` to parse and hash these. The credential id and start receipt hash retain their separate raw-hex formats. |

### 1.5 Size and CI race cost

After each commit use `git diff -M --numstat f8d8e3ad HEAD -- cmd host docs app`. Report the sum of both numeric columns. The planner used `git diff --no-index -M --numstat "$S/base" "$S/st-$M"` on archive exports. Every product line is inside the filter, so the unfiltered product sum is the same. The final JSON records the unfiltered total including the two plan records. The controller may reproduce that count with the same command without pathspecs.

| Milestone | Added | Deleted | Cumulative | Increment | Snap STATUS / DELETED / manifest |
|---|---:|---:|---:|---:|---|
| M1 | 95 | 0 | 95 | 95 | 2/0/2 |
| M2 | 1376 | 24 | 1400 | 1305 | 15/0/15 |

The first compiled sizing pass measured 1,289 product lines with the core transaction and attended recovery. The design budgets 510 raw / 689 calibrated before resolver and consumer. The measured cut uses the full 1,400-line cap. The S3b bank alone is 132 lines before any process/drain gate. Appendix C measures the carried surfaces and identifies which sizes are estimates. The touched packages are `cmd/ailang-worldd`, `host/approval`, `host/authority` and `host/store`. Linux CI uses ubuntu-24.04 and the existing eight-minute race limit. Row 267 has five slow `host/daemon` race failures near 480.1 s. The sprint adds no tests to the two exposed packages.

## 2. How the executor applies a sub-milestone

````bash
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter432-row243-v4s5b
S=$HOME/.ailang/state/world-iter432/plan-scratch
T=$HOME/.ailang/state/world-iter432/executor
B=f8d8e3ad27b8442e1cd62674052b0392f03089a6
PLAN=$W/design_docs/planned/w-world-address-bar-verbs-v4s5b-sprint-plan.md
mkdir -p "$T"
for name in accept.sh gates.sh snapshot.sh; do
  awk -v n="$name" '$0 == "````bash " n {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$name"
  test -s "$T/$name"
done
M=M1
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
(cd "$W" && /usr/bin/patch -p1 -E < "$T/$M.patch")
AILANG_BIN=$HOME/.pinned-ailang/ailang bash "$T/accept.sh" "$W" "$M" fast
bash "$T/snapshot.sh" "$W" "$M"
````

Repeat from patch extraction for M2 after M1 is committed, or after M1 acceptance and snapshot if the controller elects cumulative staging. Keep the two plan records out of the executor snapshot.

## 3. Milestones

### M1 — terminal session mint exclusion (95 cumulative)

`TestOperatorSessionMintEpisodeExclusion` uses an actual signed and claimed synthetic episode, a held description, a released description with stale lease, a no-tty arm, attended decline and uncontended mints. The storage view, output/prompt, branch refs, Git registration and file snapshots surround the refusal. `TestOperatorProvisionEpisodeExclusion` and `TestSessionMint_*` are existing controls.

### M2 — atomic mint proof, journal and attended recovery (1,400 cumulative)

`TestOperatorMintBinding` separates exact proof and time, field equality, runner/action identity and SSH/passkey current-state pins. `TestOperatorMintRollback` injects SQLite faults before the credential insert, before the start insert and after receipt update before commit. Each checks exact DB bytes, complete action/session/start reads and untouched paths. The test proves a bounded retry after verified rollback.

`TestOperatorMintReceiptWindow/transaction_separation/receipt_update` turns a premature commit into an observed partial write. The test injects SQLite failure after the attempted receipt UPDATE. The separate restart arm models a lost handoff after a successful atomic commit and withholds token publication. A reopened handle finds the exact receipt and credential with initial `grants_json` and owner. A second mint fails with no additional credential; approver use and passkey counter remain fixed. The receipt has no plaintext. At the receipt/readback-before-publication boundary, `TestOperatorMintRecovery/revocation` withholds file publication, reads the command from the refusal, executes it with a real typed `y` fixture and reads back exact action generation, receipt version/bytes and targeted credential revocation. This is a terminal recovery fixture without an integrated publisher. S5c injects individual staging/publish windows in its real caller.

`TestOperatorMintRecovery` also pins matching-token retention, exact expiry, contention, missing/unsafe evidence, the owner/step race across the prompt, a newly appeared file and abort rollback. `TestOperatorMintTokenReadback` tests actual regular-file fd type, mode, ownership, link count, size and byte digest. `TestQuickstartRunnerMintRecovery` pins the exact documented clause and fixture command.

````bash accept.sh
#!/bin/bash
set -euo pipefail
W=$1; M=$2; PART=$3
export AILANG_BIN="$HOME/.pinned-ailang/ailang"
export WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
S="$HOME/.ailang/state/world-iter432/plan-scratch"
export AILANG_STATE_DIR="$S/accept-state" AILANG_CACHE_DIR="$S/accept-cache"
mkdir -p "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR"
cd "$W"
test "$PART" = fast
run() {
  local pkg=$1 names=$2 name log
  log=$(mktemp "$S/accept.XXXXXX")
  go test -v -count=1 -timeout 1m "./$pkg" -run "^($names)\$" | tee "$log"
  while IFS= read -r name; do
    if [[ "$name" == *'.*' ]]; then continue; fi
    grep -Fq -- "--- PASS: $name (" "$log"
  done < <(printf '%s\n' "$names" | tr '|' '\n')
}
case "$M" in
  M1|M2) run cmd/ailang-worldd 'TestOperatorSessionMintEpisodeExclusion|TestOperatorProvisionEpisodeExclusion|TestSessionMint_.*' ;;
  *) exit 2 ;;
esac
if [ "$M" = M2 ]; then
  run cmd/ailang-worldd 'TestOperatorMintRecovery|TestOperatorMintRecoveryArgs|TestOperatorMintTokenReadback|TestQuickstartRunnerMintRecovery|TestQuickstartServiceSectionMatchesTheCLIHelp|TestOperatorServiceCancelConfirmation|TestCLICredentialWriterSites'
  run host/authority 'TestOperatorMintBinding|TestOperatorMintCollision|TestOperatorMintRollback|TestOperatorMintReceiptWindow|TestOperatorConfirmReplay|TestOperatorRunnerRechecks'
  run host/approval 'TestOperatorRunnerRecheckLaw|TestOperatorPureEntryPointsWriteNothing'
  run host/store 'TestCredentialsMethodsRequireADeadline|TestProductionGoSurface|TestProductionContextRoots'
fi
printf 'ACCEPT_OK %s %s\n' "$M" "$PART"
````

### 3.1 Branch traces

For M1 the no-effect comparison includes credentials DB bytes, Git refs and registration, workspace files, the output file, stdout and prompt. M2's transaction fixture compares full action/session/start rows, DB bytes (including receipt, approver and passkey evidence) and task/MCP/token paths. The recovery fixture also compares the workspace marker, fake supervisor calls and unrelated session. A Git worktree registration, actual publisher and child dispatch are outside the M2 entry.

| Input or fault | Decision and readback |
|---|---|
| Direct mint cannot open tty | Usage refusal; no sidecar, prompt, store write or token output. |
| Direct mint admitted episode, existing holder or durable owner after holder exits | Existing nonwaiting conflict and exact action id before prompt. The stale lease cannot bypass the owner. |
| Independent valid episode, attended decline, and legacy opaque episode | Valid episode acquires/releases the canonical description. Decline releases with stable evidence. Opaque episode follows the prior typed terminal path and preserves its exact session episode. |
| Mint credentials handle missing deadline, schema or writer | Structured store refusal before writes. `TestCredentialsMethodsRequireADeadline` exercises the new I/O APIs. |
| Missing, zero or confirmation-only seal, short challenge, changed approval bytes or reconstructed review | Proof refusal with no credentials or receipt. The real signer fixture exercises existing SSHSIG and passkey binding. |
| Missing or stale action id/digest/generation, wrong row episode, phase or kind | Missing/conflict with prior owner and step intact. The seal binds id, kind, challenge, expiry and canonical action bytes and the mint requires headless start. |
| Wrong runner nonce/generation or current catalogue generation | Runner conflict. A current catalogue changed independently of the signed review is an approval binding conflict. |
| Clock before action creation or at approval expiry | Conflict or expiry with no mint, including expiry between the initial recheck and transaction clock. |
| Changed approver name, format, key, revocation or enrollment time | In-transaction conflict. Successful mint leaves last use unchanged. |
| Passkey credential, RP, origin, backup or recorded count changed | In-transaction conflict. Successful mint preserves the confirmed counter. |
| Effective grant effect/scope/budget, TTL, start digest, profile project/digest or signed token path differs | Exact binding field. Start uses a transaction-owned timestamp, publishes the established tagged profile/task digests and uses a raw-hex credential id. |
| Existing branch/history, stale task digest, noncanonical token or task path, invalid binding, zero TTL or overflow | Reconcile nothing and issue no credential. A new action must pass the stricter review-to-effect checks in Appendix C before activation. |
| Existing episode start/session, stale mint reservation identity/owner/version/phase | No adoption. Session absence or start uniqueness blocks the collision. Any later CAS refusal rolls the same session/start transaction back. |
| SQLite error at session insert, start insert or receipt update | Prior reserved state, session list and start preserved. A second test moves commit before receipt and the fault reveals partial persisted writes. |
| Fresh signed start with exact effective mint and a successful commit | One exact session/start/observed mint receipt. Only the bound call receives the raw token. Receipt, approver, passkey, action owner and time are read back. |
| Lost handoff after commit, retry before publication | A reopened store locates the atomic receipt; the attempted remint returns no token and preserves the one credential and start. |
| Resolver malformed args, absent approved configuration, wrong action state, damaged start/receipt, missing credential, or later reserved effect | Typed usage or evidence blocker. Retain owner and all files. The resolver lacks a dispatch path. |
| Missing/contended descriptor or nonprivate parent | Lock/path refusal before revocation. Scope remains the selected installed configuration and action. |
| Token symlink, nonregular, wrong UID/mode/links, size or bytes | Portable refusal with unchanged file evidence. The bound command reads the descriptor itself. |
| Exact live token, expiry or appearance after review | Exact live token inspection preserves the owner for file recovery. Expiry and newly appeared evidence refuse. |
| Missing token beneath an established private directory, no tty or declined approval | No revoke. Exact retry text names the identified action and is executed verbatim in the fixture. |
| Missing token, typed approval, same owner and mint step | Delete just the receipt credential under the episode descriptor. Fail the action, advance its generation, retain start, receipt, workspace and unrelated session. |
| Recovery owner or step changed, deletion fails, or final action CAS fails | Independent transaction read and CAS refuse. Any abort SQL failure rolls back the targeted deletion along with the action transition. |

### 3.2 Design §5.3 row assignment

| Design row | Disposition |
|---|---|
| `TestOperatorStartOnce` | M2 pins one start, session and stable mint receipt in `TestOperatorMintReceiptWindow`. S5c earns one worktree and file set. S6 earns receipt-before-spawn, confinement and stable launch link. |
| `TestOperatorStartChanged` | M2 pins independent effective mint proof/clock checks in `TestOperatorMintBinding`. S5c earns live preset/base/profile/tool/task/daemon/world/store, paths and catalogue reconstruction over the S5 seam. |
| `TestOperatorStartCollisions` | M2 denies an existing credential, start PK and wrong/observed reserved identity in `TestOperatorMintCollision` and `TestOperatorMintBinding`. Full branch/registration/path provenance/no-adoption tests stay S5c. |
| `TestOperatorWorktreeWindow` | S5c, with attended worktree-specific evidence and owner/step CAS in its first uncertain write. |
| `TestOperatorMintBinding`, `TestOperatorMintRollback` | M2, with real SSH/passkey signer and real SQLite step faults. |
| `TestOperatorMintReceiptWindow` | M2, receipt in tx, recovery after reopen, no remint, exact private token or attended orphan revoke. Multi-file publication windows stay S5c. |
| `TestOperatorFilesReconcile` | M2 takes only bounded private token inspection. S5c owns full task/token/MCP staging, no replacement, provenance and attended reconciliation. |
| `TestOperatorLaunchSession`, `TestOperatorSpawnRefused`, `TestOperatorLaunchCrashWindow`, `TestOperatorLaunchExit` | S6, with E3, process receipt and pre-dispatch gates. |
| `TestOperatorStartBrowser` | S6 synthetic restart/launch half, S8 browser close/reconnect and one durable start. S7 owns routes/status. |
| `TestOperatorHerdrUnavailable`, `TestOperatorHerdrHandoff`, `TestOperatorHerdrWindow` | S12 under D-WORLD-119=A and E4. |

## 4. Mutation table

The diff-driven ledger has **123 / 123 named kills**. The mutation replayer applies each edit in a disposable copy of its recorded stage with the same pinned environment as §3. Full commands and patches are in `$S/mutation-ledger.json`; `$S/mutants.json` and `$S/mutate.py` replay them. The named command takes the form `go test -count=1 -timeout 1m -v -run '^TEST$/^subtest$' ./PACKAGE`. Whole-package results name the package holding the behavioral test. The implementation-only column gives the additional run for shared host code. A whole-package survivor there means the behavioral pin is in a downstream package.

Run the unmeasured whole-package versions outside the sandbox from the validated snapshot stages. For each row in the CLI list the replayer invokes `go test -count=1 -timeout 8m -v ./cmd/ailang-worldd` with the specified edit. The legacy replayer runs the same command against the three old mutations at base.

````bash
export AILANG_BIN="$HOME/.pinned-ailang/ailang"
xargs python3 "$S/mutate.py" --whole-only --jobs 3 < "$S/controller-cli-mutants.txt"
python3 "$S/mutate.py" --legacy --whole-only --jobs 1
````

Legacy repro on the pristine export: the targeted S5 tests survived each exact judge edit. Base whole-package runs hit inherited loopback denial for all three and are **UNMEASURED IN SANDBOX**. J-PRINT-AGENT and J-DOUBLE-OVERRIDE-APPEND remain parked with their measured repairs in Appendix C. The M2 J-QUICKSTART-b named kill is a documentation pin. With that pin excluded, the exact remedy and existing quickstart behavior tests still pass. Executing the remedy is pinned independently by R-REMEDY. The two focused commands and logs are retained in scratch.

Non-vacuity. NV-M1 restores the entire changed production file from the base. NV-M2 restores the prior runner dispatch and makes the new command fail through the real CLI switch. The complete M2 production revert was also tested on a fresh M1 copy with the unchanged, independently compilable `TestOperatorMintRecoveryArgs` case extracted from its final test file. `go test -v -count=1 -timeout 1m ./cmd/ailang-worldd -run '^TestOperatorMintRecoveryArgs$'` fails with `unknown runner operation`. The command passes on M2. See `$S/M2-full-production-revert.log`. Removing the transaction commit or splitting it before receipt independently reds the new store behavior and leaves the buildable API shape intact. The exact `M-LATE-RECEIPT` mutation commits credential and start first, then attempts a separate receipt write. `TestOperatorMintReceiptWindow/transaction_separation/receipt_update` detects the surviving partial writes after the injected fault.

Defensive and shared cases. `encoding/json` can't fail for slices of the plain grant structs. `crypto/rand.Read` in the pinned Go version either fills the buffer or terminates, which leaves the Go error-return arm unreachable. The SQL receipt read and final CAS overlap before the only in-transaction write, the precise credential delete. Removing a repeated predicate from the final UPDATE while leaving the independent snapshot equality would be equivalent under the shipped schema with no triggers. The retained CAS keeps the exact owner and step in the effect write itself. The action kind and episode are immutable through the supported APIs. Start-only and terminal-only guards stay explicit at boundaries. Distinct operational faults are represented by the earlier proof/refusal class and injected SQL failures, with no release from generic database corruption.

| ID | Diff or reused guard anchor | Named behavioral test | Named | Whole package | Implementation only |
|---|---|---|---|---|---|
| L-FENCE | `CLI/session.go:268` | `TestOperatorSessionMintEpisodeExclusion/durably_owned` | KILLED | UNMEASURED IN SANDBOX | — |
| L-OPAQUE | `CLI/session.go:268` | `TestOperatorSessionMintEpisodeExclusion/free_and_opaque` | KILLED | UNMEASURED IN SANDBOX | — |
| L-EPISODE | `CLI/session.go:269` | `TestOperatorSessionMintEpisodeExclusion/durably_owned` | KILLED | UNMEASURED IN SANDBOX | — |
| L-STORE | `CLI/session.go:269` | `TestOperatorSessionMintEpisodeExclusion/durably_owned` | KILLED | UNMEASURED IN SANDBOX | — |
| L-EXIT | `CLI/session.go:269` | `TestOperatorSessionMintEpisodeExclusion/contended` | KILLED | UNMEASURED IN SANDBOX | — |
| L-TEXT | `CLI/session.go:269` | `TestOperatorSessionMintEpisodeExclusion/durably_owned` | KILLED | UNMEASURED IN SANDBOX | — |
| L-CLOSE | `CLI/session.go:274` | `TestOperatorSessionMintEpisodeExclusion/declined` | KILLED | UNMEASURED IN SANDBOX | — |
| NV-M1 | `CLI/session.go:1` | `TestOperatorSessionMintEpisodeExclusion/durably_owned` | KILLED | UNMEASURED IN SANDBOX | — |
| P-MINT-NIL | `host/approval/operator_signed.go:195` | `TestOperatorMintBinding/proof` | KILLED | KILLED | `host/approval` SURVIVED |
| P-MINT-ZERO | `host/approval/operator_signed.go:195` | `TestOperatorMintBinding/proof` | KILLED | KILLED | `host/approval` SURVIVED |
| P-MINT-CONFIRM | `host/approval/operator_signed.go:195` | `TestOperatorMintBinding/proof` | KILLED | KILLED | `host/approval` SURVIVED |
| P-SEAL-NIL | `host/approval/operator_signed.go:204` | `TestOperatorMintBinding/proof` | KILLED | KILLED | `host/approval` SURVIVED |
| P-SEAL-ZERO | `host/approval/operator_signed.go:204` | `TestOperatorMintBinding/proof` | KILLED | KILLED | `host/approval` SURVIVED |
| P-SEAL-ARTIFACT | `host/approval/operator_signed.go:204` | `TestOperatorMintBinding/artifact` | KILLED | KILLED | `host/approval` SURVIVED |
| P-SEAL-BINDS | `host/approval/operator_signed.go:207` | `TestOperatorMintBinding/different_challenge` | KILLED | KILLED | `host/approval` SURVIVED |
| A-CHALLENGE | `host/authority/mint.go:206` | `TestOperatorMintBinding/thin_challenge` | KILLED | KILLED | — |
| A-RECHECK | `host/authority/mint.go:209` | `TestOperatorMintBinding/changed_review` | KILLED | KILLED | — |
| A-ACK | `host/authority/mint.go:219` | `TestOperatorMintBinding/fresh_clock` | KILLED | KILLED | — |
| D-COMMIT-CONTEXT | `host/store/operator_mint.go:95` | `TestCredentialsMethodsRequireADeadline` | KILLED | KILLED | — |
| D-READ-CONTEXT | `host/store/operator_mint.go:242` | `TestCredentialsMethodsRequireADeadline` | KILLED | KILLED | — |
| D-ABORT-CONTEXT | `host/store/operator_mint.go:290` | `TestCredentialsMethodsRequireADeadline` | KILLED | KILLED | — |
| M-WRITER | `host/store/operator_mint.go:98` | `TestOperatorMintBinding/read_only` | KILLED | KILLED | `host/store` SURVIVED |
| M-EXPIRY | `host/store/operator_mint.go:112` | `TestOperatorMintBinding/expiry_boundary` | KILLED | KILLED | `host/store` SURVIVED |
| M-CREATED | `host/store/operator_mint.go:115` | `TestOperatorMintBinding/clock_behind` | KILLED | KILLED | `host/store` SURVIVED |
| M-STATE | `host/store/operator_mint.go:115` | `TestOperatorMintBinding/action_phase` | KILLED | KILLED | `host/store` SURVIVED |
| M-OWNER | `host/store/operator_mint.go:115` | `TestOperatorMintBinding/action_owner` | KILLED | KILLED | `host/store` SURVIVED |
| M-RUN-GEN | `host/store/operator_mint.go:120` | `TestOperatorMintBinding/runner_generation` | KILLED | KILLED | `host/store` SURVIVED |
| M-RUN-OWNER | `host/store/operator_mint.go:120` | `TestOperatorMintBinding/runner_owner` | KILLED | KILLED | `host/store` SURVIVED |
| M-RUN-CAT | `host/store/operator_mint.go:120` | `TestOperatorMintBinding/catalogue_generation` | KILLED | KILLED | `host/store` SURVIVED |
| M-HEADLESS | `host/store/operator_mint.go:132` | `TestOperatorMintBinding/start_options` | KILLED | KILLED | `host/store` SURVIVED |
| M-APPROVED-BY | `host/store/operator_mint.go:139` | `TestOperatorMintBinding/approver_identity` | KILLED | KILLED | `host/store` SURVIVED |
| M-APP-KEY | `host/store/operator_mint.go:72` | `TestOperatorMintBinding/replaced_key` | KILLED | KILLED | `host/store` SURVIVED |
| M-APP-REVOKED | `host/store/operator_mint.go:72` | `TestOperatorMintBinding/revoked` | KILLED | KILLED | `host/store` SURVIVED |
| M-APP-ENROL | `host/store/operator_mint.go:72` | `TestOperatorMintBinding/future_enrolment` | KILLED | KILLED | `host/store` SURVIVED |
| M-PK-ID | `host/store/operator_mint.go:81` | `TestOperatorMintBinding/passkey_credential` | KILLED | KILLED | `host/store` SURVIVED |
| M-PK-RP | `host/store/operator_mint.go:81` | `TestOperatorMintBinding/passkey_rp` | KILLED | KILLED | `host/store` SURVIVED |
| M-PK-ORIGIN | `host/store/operator_mint.go:81` | `TestOperatorMintBinding/passkey_origin` | KILLED | KILLED | `host/store` SURVIVED |
| M-PK-BACKUP | `host/store/operator_mint.go:81` | `TestOperatorMintBinding/passkey_backup` | KILLED | KILLED | `host/store` SURVIVED |
| M-PK-COUNT | `host/store/operator_mint.go:81` | `TestOperatorMintBinding/passkey_counter` | KILLED | KILLED | `host/store` SURVIVED |
| M-PK-USE | `host/store/operator_mint.go:79` | `TestOperatorMintBinding/passkey_counter` | KILLED | KILLED | `host/store` SURVIVED |
| M-EPISODE-ROW | `host/store/operator_mint.go:183` | `TestOperatorMintBinding/episode_row` | KILLED | KILLED | `host/store` SURVIVED |
| M-BRANCH | `host/store/operator_mint.go:191` | `TestOperatorMintBinding/existing_branch` | KILLED | KILLED | `host/store` SURVIVED |
| M-RESUMES | `host/store/operator_mint.go:191` | `TestOperatorMintBinding/resume_history` | KILLED | KILLED | `host/store` SURVIVED |
| M-LAUNCHES | `host/store/operator_mint.go:191` | `TestOperatorMintBinding/launch_history` | KILLED | KILLED | `host/store` SURVIVED |
| M-PRESET | `host/store/operator_mint.go:192` | `TestOperatorMintBinding/preset_name` | KILLED | KILLED | `host/store` SURVIVED |
| M-TASK-NIL | `host/store/operator_mint.go:192` | `TestOperatorMintBinding/missing_task` | KILLED | KILLED | `host/store` SURVIVED |
| M-TASK-DIGEST | `host/store/operator_mint.go:193` | `TestOperatorMintBinding/task_digest` | KILLED | KILLED | `host/store` SURVIVED |
| M-PATH-ABS | `host/store/operator_mint.go:196` | `TestOperatorMintBinding/relative_token` | KILLED | KILLED | `host/store` SURVIVED |
| M-PATH-CLEAN | `host/store/operator_mint.go:196` | `TestOperatorMintBinding/unclean_token` | KILLED | KILLED | `host/store` SURVIVED |
| M-PATH-LEAF | `host/store/operator_mint.go:197` | `TestOperatorMintBinding/wrong_token` | KILLED | KILLED | `host/store` SURVIVED |
| M-PATH-TASK | `host/store/operator_mint.go:198` | `TestOperatorMintBinding/wrong_task` | KILLED | KILLED | `host/store` SURVIVED |
| M-TTL-FLOOR | `host/store/operator_mint.go:201` | `TestOperatorMintBinding/zero_ttl` | KILLED | KILLED | `host/store` SURVIVED |
| M-TTL-LIMIT | `host/store/operator_mint.go:201` | `TestOperatorMintBinding/overflow` | KILLED | KILLED | `host/store` SURVIVED |
| M-GRANT-COUNT | `host/store/operator_mint.go:201` | `TestOperatorMintBinding/no_grants` | KILLED | KILLED | `host/store` SURVIVED |
| M-PROFILE-PROJECT | `host/store/operator_mint.go:202` | `TestOperatorMintBinding/profile_blank` | KILLED | KILLED | `host/store` SURVIVED |
| M-PROFILE-SHA | `host/store/operator_mint.go:202` | `TestOperatorMintBinding/profile_hash` | KILLED | KILLED | `host/store` SURVIVED |
| M-TOKEN-SHA | `host/store/operator_mint.go:202` | `TestOperatorMintBinding/credential_hash` | KILLED | KILLED | `host/store` SURVIVED |
| M-HASH-LOWER | `host/store/operator_mint.go:226` | `TestOperatorMintBinding/credential_upper` | KILLED | KILLED | `host/store` SURVIVED |
| M-FACT-CAT | `host/store/operator_mint.go:210` | `TestOperatorMintBinding/catalogue_review` | KILLED | KILLED | `host/store` SURVIVED |
| M-FACT-TTL | `host/store/operator_mint.go:210` | `TestOperatorMintBinding/ttl` | KILLED | KILLED | `host/store` SURVIVED |
| M-FACT-BASE | `host/store/operator_mint.go:210` | `TestOperatorMintBinding/base` | KILLED | KILLED | `host/store` SURVIVED |
| M-FACT-GRANT-EFFECT | `host/store/operator_mint.go:210` | `TestOperatorMintBinding/grant_effect` | KILLED | KILLED | `host/store` SURVIVED |
| M-FACT-GRANT-SCOPE | `host/store/operator_mint.go:210` | `TestOperatorMintBinding/grant_scope` | KILLED | KILLED | `host/store` SURVIVED |
| M-FACT-GRANT-BUDGET | `host/store/operator_mint.go:210` | `TestOperatorMintBinding/grant_budget` | KILLED | KILLED | `host/store` SURVIVED |
| M-FACT-PROJECT | `host/store/operator_mint.go:210` | `TestOperatorMintBinding/project` | KILLED | KILLED | `host/store` SURVIVED |
| M-FACT-PROFILE | `host/store/operator_mint.go:210` | `TestOperatorMintBinding/profile` | KILLED | KILLED | `host/store` SURVIVED |
| M-STEP-VERSION | `host/store/operator_mint.go:167` | `TestOperatorMintBinding/step_version` | KILLED | KILLED | `host/store` SURVIVED |
| M-STEP-OWNER | `host/store/operator_mint.go:167` | `TestOperatorMintBinding/step_owner` | KILLED | KILLED | `host/store` SURVIVED |
| M-STEP-PHASE | `host/store/operator_mint.go:167` | `TestOperatorMintBinding/step_phase` | KILLED | KILLED | `host/store` SURVIVED |
| M-STEP-IDENT | `host/store/operator_mint.go:167` | `TestOperatorMintBinding/step_identity` | KILLED | KILLED | `host/store` SURVIVED |
| M-COLLISION | `host/store/operator_mint.go:153` | `TestOperatorMintCollision/credential` | KILLED | KILLED | `host/store` SURVIVED |
| M-SESSION-ERROR | `host/store/operator_mint.go:156` | `TestOperatorMintRollback/session_insert` | KILLED | KILLED | `host/store` SURVIVED |
| M-START-ERROR | `host/store/operator_mint.go:161` | `TestOperatorMintRollback/start_insert` | KILLED | KILLED | `host/store` SURVIVED |
| M-RECEIPT-ERROR | `host/store/operator_mint.go:169` | `TestOperatorMintRollback/receipt_update` | KILLED | KILLED | `host/store` SURVIVED |
| M-COMMIT | `host/store/operator_mint.go:172` | `TestOperatorMintReceiptWindow` | KILLED | KILLED | `host/store` SURVIVED |
| M-RECEIPT | `host/store/operator_mint.go:166` | `TestOperatorMintReceiptWindow` | KILLED | KILLED | `host/store` SURVIVED |
| M-RECEIPT-VERSION | `host/store/operator_mint.go:166` | `TestOperatorMintReceiptWindow` | KILLED | KILLED | `host/store` SURVIVED |
| M-GRANT-EXPIRY | `host/store/operator_mint.go:217` | `TestOperatorMintReceiptWindow` | KILLED | KILLED | `host/store` SURVIVED |
| M-START-CLOCK | `host/store/operator_mint.go:160` | `TestOperatorMintReceiptWindow` | KILLED | KILLED | `host/store` SURVIVED |
| R-CLOSE | `CLI/operator_mint.go:33` | `TestOperatorMintRecovery/revocation` | KILLED | UNMEASURED IN SANDBOX | — |
| R-LOCK-CREATE | `CLI/operator_mint.go:29` | `TestOperatorMintRecovery/missing_lock` | KILLED | UNMEASURED IN SANDBOX | — |
| R-LOCK-CONTENTION | `CLI/operator_mint.go:29` | `TestOperatorMintRecovery/busy` | KILLED | UNMEASURED IN SANDBOX | — |
| R-TTY | `CLI/operator_mint.go:50` | `TestOperatorMintRecovery/revocation` | KILLED | UNMEASURED IN SANDBOX | — |
| R-REMEDY | `CLI/operator_mint.go:49` | `TestOperatorMintRecovery/revocation` | KILLED | UNMEASURED IN SANDBOX | — |
| R-TOKEN-RECHECK | `CLI/operator_mint.go:61` | `TestOperatorMintRecovery/appeared` | KILLED | UNMEASURED IN SANDBOX | — |
| R-TOKEN-FIRST | `CLI/operator_mint.go:46` | `TestOperatorMintRecovery/adopt` | KILLED | UNMEASURED IN SANDBOX | — |
| R-EXPIRED | `CLI/operator_mint.go:85` | `TestOperatorMintRecovery/expired` | KILLED | UNMEASURED IN SANDBOX | — |
| R-DIR | `CLI/operator_mint.go:98` | `TestOperatorMintRecovery/unsafe_parent` | KILLED | UNMEASURED IN SANDBOX | — |
| R-OWNER | `host/store/operator_mint.go:265` | `TestOperatorMintRecovery/wrong_step_owner` | KILLED | UNMEASURED IN SANDBOX | `host/store` SURVIVED |
| R-PHASE | `host/store/operator_mint.go:265` | `TestOperatorMintRecovery/wrong_receipt_phase` | KILLED | UNMEASURED IN SANDBOX | `host/store` SURVIVED |
| R-LATER | `host/store/operator_mint.go:266` | `TestOperatorMintRecovery/later_effect` | KILLED | UNMEASURED IN SANDBOX | `host/store` SURVIVED |
| R-STATE | `host/store/operator_mint.go:254` | `TestOperatorMintRecovery/unowned` | KILLED | UNMEASURED IN SANDBOX | `host/store` SURVIVED |
| R-UNCERTAIN | `host/store/operator_mint.go:254` | `TestOperatorMintRecovery/revocation` | KILLED | UNMEASURED IN SANDBOX | `host/store` SURVIVED |
| R-ACTION-BYTES | `host/store/operator_mint.go:176` | `TestOperatorMintRecovery/action_changed` | KILLED | UNMEASURED IN SANDBOX | `host/store` SURVIVED |
| R-START-BYTES | `host/store/operator_mint.go:282` | `TestOperatorMintRecovery/start_changed` | KILLED | UNMEASURED IN SANDBOX | `host/store` SURVIVED |
| R-RECEIPT-SHA | `host/store/operator_mint.go:273` | `TestOperatorMintRecovery/receipt_changed` | KILLED | UNMEASURED IN SANDBOX | `host/store` SURVIVED |
| R-CAS-OWNER | `host/store/operator_mint.go:306` | `TestOperatorMintRecovery/stale_owner` | KILLED | UNMEASURED IN SANDBOX | `host/store` SURVIVED |
| R-CAS-STEP | `host/store/operator_mint.go:306` | `TestOperatorMintRecovery/stale_step` | KILLED | UNMEASURED IN SANDBOX | `host/store` SURVIVED |
| R-REVOKE-ID | `host/store/operator_mint.go:309` | `TestOperatorMintRecovery/revocation` | KILLED | UNMEASURED IN SANDBOX | `host/store` SURVIVED |
| R-ABORT-ERROR | `host/store/operator_mint.go:319` | `TestOperatorMintRecovery/abort_rollback` | KILLED | UNMEASURED IN SANDBOX | `host/store` SURVIVED |
| F-NOFOLLOW | `CLI/operator_mint.go:105` | `TestOperatorMintTokenReadback/symlink_valid` | KILLED | UNMEASURED IN SANDBOX | — |
| F-ABSENCE | `CLI/operator_mint.go:106` | `TestOperatorMintTokenReadback/symlink` | KILLED | UNMEASURED IN SANDBOX | — |
| F-UID | `CLI/operator_mint.go:118` | `TestOperatorMintTokenReadback/uid` | KILLED | UNMEASURED IN SANDBOX | — |
| F-MODE | `CLI/operator_mint.go:118` | `TestOperatorMintTokenReadback/mode` | KILLED | UNMEASURED IN SANDBOX | — |
| F-LINK | `CLI/operator_mint.go:118` | `TestOperatorMintTokenReadback/link` | KILLED | UNMEASURED IN SANDBOX | — |
| F-BYTES | `CLI/operator_mint.go:122` | `TestOperatorMintTokenReadback/bytes` | KILLED | UNMEASURED IN SANDBOX | — |
| F-READ-CAP | `CLI/operator_mint.go:121` | `TestOperatorMintTokenReadback/length` | KILLED | UNMEASURED IN SANDBOX | — |
| R-ARG-COUNT | `CLI/operator_service.go:276` | `TestOperatorMintRecoveryArgs` | KILLED | UNMEASURED IN SANDBOX | — |
| R-ARG-ID | `CLI/operator_service.go:276` | `TestOperatorMintRecoveryArgs` | KILLED | UNMEASURED IN SANDBOX | — |
| J-QUICKSTART-b | `docs/QUICKSTART.md:79` | `TestQuickstartRunnerMintRecovery` | KILLED | UNMEASURED IN SANDBOX | — |
| NV-M2 | `CLI/operator_service.go:1` | `TestOperatorMintRecovery/revocation` | KILLED | UNMEASURED IN SANDBOX | — |
| M-APP-NAME | `host/store/operator_mint.go:72` | `TestOperatorMintBinding/approver_name` | KILLED | KILLED | `host/store` SURVIVED |
| M-APP-FORMAT | `host/store/operator_mint.go:72` | `TestOperatorMintBinding/approver_format` | KILLED | KILLED | `host/store` SURVIVED |
| F-EXACT-LENGTH | `CLI/operator_mint.go:122` | `TestOperatorMintTokenReadback/short` | KILLED | UNMEASURED IN SANDBOX | — |
| L-TTY-ORDER | `CLI/session.go:255` | `TestOperatorSessionMintEpisodeExclusion/tty` | KILLED | UNMEASURED IN SANDBOX | — |
| R-LOCK-NB | `CLI/operator_runner.go:35` | `TestOperatorEpisodeLockExclusive` | KILLED | UNMEASURED IN SANDBOX | — |
| R-LOCK-EX | `CLI/operator_runner.go:35` | `TestOperatorEpisodeLockExclusive` | KILLED | UNMEASURED IN SANDBOX | — |
| M-CAS-DIGEST | `host/store/operator_cas.go:34` | `TestOperatorMintBinding/digest` | KILLED | KILLED | `host/store` KILLED |
| M-CAS-GENERATION | `host/store/operator_cas.go:34` | `TestOperatorMintBinding/generation` | KILLED | KILLED | `host/store` KILLED |
| M-LATE-RECEIPT | `host/store/operator_mint.go:165` | `TestOperatorMintReceiptWindow/transaction_separation/receipt_update` | KILLED | KILLED | `host/store` SURVIVED |
| M-REMINT | `host/authority/mint.go:219` | `TestOperatorMintReceiptWindow` | KILLED | KILLED | — |
| M-TASK-TAG | `host/store/operator_mint.go:193` | `TestOperatorMintReceiptWindow` | KILLED | KILLED | `host/store` SURVIVED |
| M-FACT-TOKEN | `host/store/operator_mint.go:210` | `TestOperatorMintBinding/token_fact` | KILLED | KILLED | `host/store` SURVIVED |

## 5. Files

- `cmd/ailang-worldd/{session,service,operator_service,operator_mint}.go`
- `cmd/ailang-worldd/{operator_provision_test,operator_episode_lock_test,operator_mint_test}.go`
- `host/approval/operator_signed.go`
- `host/authority/{mint,operator_store_test,operator_mint_test}.go`
- `host/store/{credentials,credentials_test,operator_mint}.go`
- `docs/QUICKSTART.md`

## 6. Sprint gates

Use §2's extracted script. Scratch names the changed file list through `GATES_CHANGED_FILE`. The controller runs the final gate outside the sandbox with `GATES_SKIP_VERIFY_GO` and `GATES_SANDBOX` unset. Run Linux CI's existing race leg with ubuntu-24.04, and read-only launchctl observations on Darwin. In the committed checkout run `env -u GATES_SKIP_VERIFY_GO -u GATES_SANDBOX bash "$T/gates.sh" "$W" "$HOME/.ailang/state/world-iter432/controller-gates" final`. The `ALL_DONE` marker is the final file written by the gate.

````bash gates.sh
#!/bin/bash
# Usage: gates.sh <tree> <logdir> <base|final|scratch> [sha]
set -u
set -o pipefail
export GIT_OPTIONAL_LOCKS=0
T=$1; L=$2; MODE=$3
B=f8d8e3ad27b8442e1cd62674052b0392f03089a6
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
EXPECTED='cmd/ailang-worldd host/approval host/authority host/store '
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
if [ "$MODE" != base ]; then
  inv cmd/world-publish TestPluginCommandsMatchTheCLIHelp
  leg mint-quickstart "go test -v -count=1 -timeout 1m ./cmd/ailang-worldd -run '^(TestQuickstartRunnerMintRecovery|TestOperatorMintRecovery)$'"
fi
scope_fence() {
  local allowed="$L/allowed.txt" extra
  cat > "$allowed" <<'EOF'
cmd/ailang-worldd/operator_episode_lock_test.go
cmd/ailang-worldd/operator_mint.go
cmd/ailang-worldd/operator_mint_test.go
cmd/ailang-worldd/operator_provision_test.go
cmd/ailang-worldd/operator_service.go
cmd/ailang-worldd/service.go
cmd/ailang-worldd/session.go
docs/QUICKSTART.md
host/approval/operator_signed.go
host/authority/mint.go
host/authority/operator_mint_test.go
host/authority/operator_store_test.go
host/store/credentials.go
host/store/credentials_test.go
host/store/operator_mint.go
design_docs/planned/w-world-address-bar-verbs-v4s5b-sprint-plan.md
design_docs/planned/sprint_w-world-address-bar-verbs-v4s5b.json
EOF
  printf '%s\n' cmd/ailang-worldd/session.go | grep -Fx -f "$allowed"
  printf '%s\n' world/transitions.ail | grep -Fxv -f "$allowed"
  extra=$(grep -Fxv -f "$allowed" "$L/paths.txt" || test "$?" -eq 1)
  printf 'UNTOUCHED_VIOLATIONS: %s\n' "${extra:-none}"
  test -z "$extra"
}
leg scope 'scope_fence'
caller_fence() {
  local old callers lock_count cases mint commit resolver
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
  callers=$(grep -R -nF --include='*.go' 'resolvePresetStart(' cmd host | grep -vF '_test.go:' | grep -vF 'func resolvePresetStart(')
  printf 'Read seam caller:\n%s\n' "$callers"
  test "$(printf '%s\n' "$callers" | wc -l | tr -d ' ')" = 1
  printf '%s\n' "$callers" | grep -F 'cmd/ailang-worldd/episode_preset.go:'
  mint=$(grep -R -nF --include='*.go' 'MintOperatorStart(' cmd host | grep -vF '_test.go:' | grep -vF 'func MintOperatorStart(' || :)
  test -z "$mint"
  printf 'Production MintOperatorStart callers: 0\n'
  commit=$(grep -R -nF --include='*.go' '.CommitOperatorMint(' cmd host | grep -vF '_test.go:' || :)
  resolver=$(grep -R -nF --include='*.go' '.AbortOperatorMint(' cmd host | grep -vF '_test.go:' || :)
  if [ "$MODE" = base ]; then
    test -z "$commit"
    test -z "$resolver"
  else
    test "$(printf '%s\n' "$commit" | wc -l | tr -d ' ')" = 1
    printf '%s\n' "$commit" | grep -F 'host/authority/mint.go:'
    test "$(printf '%s\n' "$resolver" | wc -l | tr -d ' ')" = 1
    printf '%s\n' "$resolver" | grep -F 'cmd/ailang-worldd/operator_mint.go:'
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
TOUCHED_PACKAGES: cmd/ailang-worldd host/approval host/authority host/store 
LEG touched rc=0 0s
LEG status rc=0 0s
LEG gofmt rc=0 0s
LEG build rc=0 0s
LEG vet rc=0 2s
LEG fence rc=0 18s
LEG verify_ail rc=0 228s
LEG pkg-cmd-ailang-worldd rc=1 15s UNINFORMATIVE UNDER SANDBOX
LEG race-cmd-ailang-worldd rc=1 177s UNINFORMATIVE UNDER SANDBOX
LEG pkg-host-approval rc=0 3s
LEG race-host-approval rc=0 13s
LEG pkg-host-authority rc=0 21s
LEG race-host-authority rc=0 256s
LEG pkg-host-store rc=0 29s
LEG race-host-store rc=0 50s
LEG inv-host-approval rc=0 0s
LEG inv-host-approval-approvaltest rc=0 1s
LEG inv-host-authority rc=0 1s
LEG inv-host-boundary rc=0 2s
LEG inv-host-daemon rc=0 5s
LEG inv-host-verbs rc=0 0s
LEG inv-host-broker rc=0 4s
LEG inv-cmd-ailang-worldd rc=0 2s
LEG inv-quickstart-section9 rc=1 1s UNINFORMATIVE UNDER SANDBOX
LEG inv-host-store rc=0 2s
LEG routes rc=0 0s
LEG inv-cmd-world-publish rc=0 2s
LEG mint-quickstart rc=0 2s
LEG scope rc=0 0s
LEG callers rc=0 2s
LEG verify_go SKIPPED 0s (GATES_SKIP_VERIFY_GO=1)
LEG observation-after rc=77 0s UNINFORMATIVE UNDER SANDBOX
GATES_DONE 836s FAILS: none UNINFORMATIVE: observation-before pkg-cmd-ailang-worldd race-cmd-ailang-worldd inv-quickstart-section9 observation-after
````

The route gate returned 50 registrations and three surface POSTs. `cmp` matched the daemon, surface page/CSP and app bytes. Caller controls recorded one legacy `authority.RecheckVerb` production call, one terminal call to the S5 read seam, five legacy provision calls, zero reachable operator mint/claim/reserve/observe calls and exactly one attended abort call.

Controller repeats CLI plain, CLI race, QUICKSTART §9, full plain and Linux full race, plus before/after observations by running §6 outside the sandbox. Exact invocations, with the pinned environment from the gate, are `go test -count=1 -timeout 8m ./cmd/ailang-worldd`, `go test -race -count=1 -timeout 8m ./cmd/ailang-worldd`, `go test -v -count=1 -timeout 8m ./cmd/ailang-worldd -run '^TestQuickstartSection9FlagsMatchTheCLI$'`, `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` and `VERIFY_GO_LEGS=race ./scripts/verify_go.sh`. See §4 for per-variant whole-package verdicts and commands.

## 7. Executor directive notes

- Commit-time commands, one per milestone, use the controller-owned snapshot and the same pinned fixture environment as scratch.

````bash
bash "$HOME/.ailang/state/world-iter432/commit_ms.sh" M1 "$msgfile" "AILANG_BIN=$HOME/.pinned-ailang/ailang bash '$T/accept.sh' '$W' M1 fast"
bash "$HOME/.ailang/state/world-iter432/commit_ms.sh" M2 "$msgfile" "AILANG_BIN=$HOME/.pinned-ailang/ailang bash '$T/accept.sh' '$W' M2 fast"
````

- Snapshot is cumulative from the base, accepts missing selected directories in an export rehearsal and excludes plan records via pathspec. Do not reuse a snapshot directory. A rename is included as deletion plus addition in the snapshot and detected with `-M` at commit-time.

````bash snapshot.sh
#!/bin/bash
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=$1
M=$2
B=f8d8e3ad27b8442e1cd62674052b0392f03089a6
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter432/snap}/$M
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

Rehearse using `mkdir -p "$S/executor-snap-rehearsal"; SNAP_STATUS_LIST="$S/status-M1" SNAP_PARENT="$S/executor-snap-rehearsal" bash "$S/snapshot.sh" "$S/st-M1" M1` and repeat for M2. The controller runs `bash "$HOME/.ailang/state/world-iter432/ctl_roundtrip.sh" M1 M2` against its own fresh archive. The planner compared the full fresh export to each stage after extracting from §Appendix P with the same awk used by `ctl_roundtrip.sh`. Both `diff -r` outputs were empty. Extracted acceptance passed again on both fresh exports. Snapshot rehearsal yielded 2/0/2 and 15/0/15 with all manifest checks OK. A separate host-only rehearsal yielded 1/0/1 with `cmd/`, `docs/` and `app/` absent. Logs are in scratch.

## 8. Risks

- Darwin and Linux use exclusive, non-waiting flock on independent descriptions. Use the existing EWOULDBLOCK/EAGAIN behavior and finite watchdog for any lock-flag mutation. Type, symlink and no-change assertions are portable across the two systems.
- Review-to-current equality is per field and temporal. The current signed catalogue and app/daemon/store pair must be reconstructed from the S5 read seam before any effect activation. A heartbeat updates liveness and a generation. It never proves a child died or releases a durable owner.
- A valid token check only preserves the exact token and owner for the later file stage. The proof of a worktree, task/MCP and launch is distinct. Missing directory or ambiguous path cannot establish token absence. Compatible processes must remain fenced until full recovery and drain exist.
- Readiness and GIT process policy live at the companion consumer. A stale build, active child, stopped/replaced runner or mismatched store remains gated by §§2.2 and 5.5.

## 9. Acceptance criteria

1. M1 and M2 apply from the base and reproduce the compiled stages using the controller's round trip. Filtered counts and cumulative manifests match §1.5.
2. M1 refuses exact owner before prompt or mint, releases the fd on decline and normal mint, and preserves previous typed paths and independent sessions. Existing public A3 and start users keep their separate approval and writer behavior.
3. M2 proves exact effective proof plus the fresh transactional read of owner, step and enrollment, atomic SQLite rollback and retained commit receipts. The attended remedy is executable verbatim and uses independent current-owner/step CAS and effect readback.
4. No production operator mint/claim/reserve/observe caller or served change. Inventory and race checks stand. §4 records named kill status and the controller replays unmeasured whole-package runs.
5. S5c carries the worktree/file uncertainty remedies and paired-process gates before any production activation.

## Appendix C — after this sprint

- **S5c worktree, changed inputs and collision boundary. 113 measured sizing lines.** The compiled `s5c_resolution_probe.go` projection adds 54 lines. The read-only NUL-delimited Git registration, path, common repository and exact starting SHA sketch in `s5c_worktree_probe.go` adds 59. These measurements cover input projection and observation. S5c still needs the full mutation matrix, reservation, partial effect and attended disposition code. Reuse S5's shared resolution and instrument a fresh reconstruction of repo/remote, preset, base/ref/OID, task bytes, tool/daemon/profile/build, catalogue, exact episode and all canonical paths. Keep terminal permissiveness for an explicitly named free branch, and require the operator action's new-branch and worktree absence proofs. Reserve the exact worktree before issuing Git add. Query registration identity, repo/branch/starting SHA and physical path. Known absence is bounded retry or reviewed stop. Partial or unattributable state requires the same-sprint attended remedy with effect-specific readback and independent owner/step CAS. Never remove an unattributed worktree. No production caller in S5b.
- **S5c files. 53 measured sizing lines.** The compiled `s5c_files_probe.go` is an fd readback sketch. The full no-overwrite publication and attended disposition remain to be measured. Use private exclusive/no-replace staging and unchanged-file readback for canonical task/token/MCP. Pin path, ownership, mode, digest, link count, file type, exact bytes and actual writer ordering independently. Bind partial file publication to the owner and step version, keep all repair text tested verbatim, and preserve the worktree and unrelated authority. (Re)validate the current credential and the verified token before any launch. The S5b handler preserves on ambiguous file evidence.
- **S5c/S6 paired-process consumer and drain, binding. 72 measured sizing lines plus the 132-line S3b bank.** `TestS5cEpisodeFDPrototype` passed with an inherited open description in a separate process, a barrier after the parent closed its descriptor and reacquisition after the child exited. It measures fd inheritance only. The four product pins below and the lifecycle changes still need implementation and a cap check. A live claim/reserve/observe/mint is prohibited until the process holds exactly the same canonical episode fd before CAS, across effect boundaries. Required named process pins are `TestOperatorFirstManagedEpisodeProcess` for claim vs a terminal operation, `TestOperatorOldBuildTerminalDrain` for retiring the old terminal build, `TestOperatorCurrentOwnerProcessRecovery` for a competing process, and `TestOperatorServiceUpgrade/active_child` for paired upgrade. Add fd lifecycle pins and mutation matrix for each reservation/receipt selector and current-owner transition. Test stop, uninstall, matched restart and schema-compatible rollback through drain or retained uncertainty. No missed heartbeat shortcut. Include direct terminal `session revoke`, concurrent credential changes and exact session revalidation before launch or recovery. The reviewed pre-v10 backup remains attended. The S5b read-only adoption and attended abort rely on the canonical episode lock. Future production adoption/recovery must integrate the same process gates.
- **S3b, 132 measured lines.** The bank adds the two step methods and one restart fixture. Bring the fd/process test, full receipt-predicate mutations and current-owner recovery together with any production reserve/observe consumer. S5b tests use SQL only to construct an inert reservation.
- **J-OWN-UID, existing runner lock/path check.** Dropping the actual UID comparison changes one source line, 2 in numstat (`git diff --no-index -M --numstat "$S/J-OWN-UID-before.go" "$S/J-OWN-UID-after.go"`). Actual UID was 501. The alternate-owner fixture is unmeasured and stays with the process gate. Keep the actual kernel UID check and provision a controller alternate-owner fixture with ordinary non-root execution. The new token fd readback tests expected UID against real stat. They do not replace the existing runner lock/path alternate-owner and fd proof. Carry `Z-ABSENT-ANY` and `Z-WAIT-READBUDGET` as in S5 Appendix C.
- **Legacy parity survivors, 8 measured fixture lines.** The scratch parity patch adds 5 lines to assert the printed Codex command and 3 to compare persisted overrides with an independent fixed slice. `printLaunch` uses the agent to select argv and prints no separate `agent:` line. The original Claude positive takes the default case too, so the parked fixture exercises Codex. Both parked tests killed their variant. `J-PRINT-AGENT` and `J-DOUBLE-OVERRIDE-APPEND` go with the S5c seam edits. S5b touches neither legacy preset production nor its parity test file. The QUICKSTART mutation is covered here because M2 touches the same runbook.
- **S5c measured sizing floor.** `git diff --no-index -M --numstat "$S/st-M2/cmd" "$S/st-s5c-footprint/cmd"` measured 246 changed lines (54 + 59 + 53 + 72 + 5 + 3). `git diff --no-index -M --numstat "$S/base/host" "$S/st-s3b-bank/host"` measured 132 (60 fixture + 72 implementation), and `TestOperatorStepRestart` passed against the reapplied bank. Together these yield 378 measured lines. Full S5c implementation size remains to be measured, including attended worktree and file dispositions and lifecycle consumer pins. Split S5c ahead of activation if the compiled total again exceeds the cap.
- **Launch, browser and finish.** The S6 890 raw / 1,202 calibrated, S8 860 / 1,161 and S12 700 / 945 design figures are planning budgets. Retain E3, E5 and #604 sequencing, E4, and the D-WORLD-120 finish resolver before S9. No new finish authorization in this cut.

## Appendix P — compiled patches

### Patch M1

````diff M1
diff --git a/cmd/ailang-worldd/operator_provision_test.go b/cmd/ailang-worldd/operator_provision_test.go
index 7152f78..8647ed6 100644
--- a/cmd/ailang-worldd/operator_provision_test.go
+++ b/cmd/ailang-worldd/operator_provision_test.go
@@ -134,3 +134,87 @@ func requireProvisionLockReleased(t *testing.T, cred, ep string) {
 	}
 	fd.Close()
 }
+
+func TestOperatorSessionMintEpisodeExclusion(t *testing.T) {
+	gc := debug.SetGCPercent(-1)
+	t.Cleanup(func() { debug.SetGCPercent(gc) })
+	r := newEpisodeRun(t)
+	creds, err := credentialsFor(r.f.db, "")
+	if err != nil {
+		t.Fatal(err)
+	}
+	owned := "ep1"
+	held := claimFixtureEpisode(t, creds, owned)
+	token := filepath.Join(r.sessions, "legacy.session")
+	if err := os.Mkdir(r.sessions, 0700); err != nil {
+		t.Fatal(err)
+	}
+	if err := os.WriteFile(token, []byte("existing handoff"), 0600); err != nil {
+		t.Fatal(err)
+	}
+	args := func(ep string, extra ...string) []string {
+		return append([]string{"--db", r.f.db, "--episode", ep, "--grant", "fs.read=repo:1"}, extra...)
+	}
+	for _, name := range []string{"tty", "contended", "durably_owned"} {
+		t.Run(name, func(t *testing.T) {
+			ep := owned
+			if name == "tty" {
+				ep = "tty-required"
+			}
+			if name == "durably_owned" {
+				if err := held.Close(); err != nil {
+					t.Fatal(err)
+				}
+			}
+			before := observeProvision(t, r, creds)
+			var out, stderr bytes.Buffer
+			term := &promptTerm{in: strings.NewReader("y\n")}
+			env := testMintEnv("y\n", 1800)
+			env.openTerminal = func() (io.ReadWriteCloser, error) { return term, nil }
+			if name == "tty" {
+				env = noTerminalEnv(1800)
+			}
+			code := runSessionMint(args(ep, "--out", token), &out, &stderr, env)
+			if name == "tty" {
+				if code != exitUsage || !strings.Contains(stderr.String(), "no controlling terminal") {
+					t.Fatalf("terminal refusal: %d %s", code, &stderr)
+				}
+				if _, err := os.Lstat(operatorEpisodeLockPath(creds, ep)); !os.IsNotExist(err) {
+					t.Fatalf("terminal fence sidecar: %v", err)
+				}
+			} else {
+				expectEpisodeRefusal(t, code, stderr.String(), before.credentials, creds)
+			}
+			if out.Len() != 0 || term.out.Len() != 0 || !reflect.DeepEqual(before, observeProvision(t, r, creds)) {
+				t.Fatal("refusal changed prompt, output, credentials, worktree, refs or files")
+			}
+		})
+	}
+	t.Run("declined", func(t *testing.T) {
+		before := observeProvision(t, r, creds)
+		var out, stderr bytes.Buffer
+		if code := runSessionMint(args("declined"), &out, &stderr, testMintEnv("n\n", 1800)); code != exitUsage ||
+			out.Len() != 0 || !strings.Contains(stderr.String(), "not confirmed") ||
+			!reflect.DeepEqual(before, observeProvision(t, r, creds)) {
+			t.Fatalf("decline: %d %s", code, &stderr)
+		}
+		requireProvisionLockReleased(t, creds, "declined")
+	})
+	t.Run("free_and_opaque", func(t *testing.T) {
+		for _, ep := range []string{"free", "legacy/opaque"} {
+			var out, stderr bytes.Buffer
+			if code := runSessionMint(args(ep), &out, &stderr, testMintEnv("y\n", 1800)); code != exitOK || len(strings.TrimSpace(out.String())) != 64 {
+				t.Fatalf("mint %q: %d %s", ep, code, &stderr)
+			}
+			rows := r.f.rows(t)
+			found := false
+			for _, row := range rows {
+				found = found || row.EpisodeID == ep
+			}
+			if !found {
+				t.Fatalf("missing persisted session for %q", ep)
+			}
+		}
+		requireProvisionLockReleased(t, creds, "free")
+	})
+}
diff --git a/cmd/ailang-worldd/session.go b/cmd/ailang-worldd/session.go
index a274755..672620c 100644
--- a/cmd/ailang-worldd/session.go
+++ b/cmd/ailang-worldd/session.go
@@ -263,6 +263,17 @@ func runSessionMint(args []string, stdout, stderr io.Writer, env sessionEnv) int
 	}
 	defer func() { _ = term.Close() }()
 
+	// Legacy session mint accepts opaque episode IDs. Operator admission
+	// rejects them, and only admitted IDs are safe as sidecar components.
+	if workspaceEpisodePattern.MatchString(*episode) {
+		lock, err := terminalEpisodeLock(credPath, *episode, nil)
+		if err != nil {
+			fmt.Fprintf(stderr, "ailang-worldd session mint: %v\n", err)
+			return exitUsage
+		}
+		defer lock.Close()
+	}
+
 	now := env.now()
 	// Each CLI-minted grant expires WITH the session (row 134 break 1): the
 	// broker's liveness rule is now < ExpiresAt in Unix seconds, so a grant left
````

### Patch M2

````diff M2
diff --git a/cmd/ailang-worldd/operator_episode_lock_test.go b/cmd/ailang-worldd/operator_episode_lock_test.go
index dc0d8f9..840bc95 100644
--- a/cmd/ailang-worldd/operator_episode_lock_test.go
+++ b/cmd/ailang-worldd/operator_episode_lock_test.go
@@ -24,6 +24,10 @@ import (
 const fixtureActionID = "oa-0123456789abcdef0123456789abcdef"
 
 func claimFixtureEpisode(t *testing.T, creds, episode string) *os.File {
+	return claimFixtureEpisodeAction(t, creds, episode, func(*approval.OperatorAction) {})
+}
+
+func claimFixtureEpisodeAction(t *testing.T, creds, episode string, prepare func(*approval.OperatorAction)) *os.File {
 	t.Helper()
 	canonical, err := store.CanonicalPath(creds)
 	if err != nil {
@@ -52,6 +56,7 @@ func claimFixtureEpisode(t *testing.T, creds, episode string) *os.File {
 	a := approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: fixtureActionID, Kind: "start", Episode: episode,
 		Store: canonical, Options: approval.OperatorOptions{Preset: "p", Headless: true}, Task: "fixture", Review: "fixture",
 		Resolution: []approval.OperatorResolved{{Name: "agent", Value: "codex"}}}
+	prepare(&a)
 	a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
 	req := store.OperatorRequest{Action: a, OperationID: strings.Repeat("1", 32), Challenge: sha256.Sum256([]byte(t.Name())), CreatedAt: 1000, ExpiresAt: 1600}
 	row, err := db.RequestOperatorAction(ctx, req, 1000)
diff --git a/cmd/ailang-worldd/operator_mint.go b/cmd/ailang-worldd/operator_mint.go
new file mode 100644
index 0000000..26bd4aa
--- /dev/null
+++ b/cmd/ailang-worldd/operator_mint.go
@@ -0,0 +1,126 @@
+package main
+
+import (
+	"crypto/sha256"
+	"errors"
+	"fmt"
+	"io"
+	"os"
+	"path/filepath"
+	"syscall"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+func (e serviceEnv) reconcileRunnerMint(id string, out, stderr io.Writer) int {
+	refuse := func(err error) int {
+		return serviceMessage(stderr, "mint recovery refused", err.Error(), "ailang-worldd service runner status", 1)
+	}
+	c, _, err := e.installedRunner()
+	if err != nil {
+		return refuse(err)
+	}
+	r, err := readRunnerMint(c.Credentials, id)
+	if err != nil {
+		return refuse(err)
+	}
+	lock, err := openRunnerLock(operatorEpisodeLockPath(c.Credentials, r.Episode), 0)
+	if err != nil {
+		return refuse(errors.New("episode recovery lock held or unavailable"))
+	}
+	defer lock.Close()
+	current, err := readRunnerMint(c.Credentials, id)
+	if err != nil {
+		return refuse(err)
+	}
+	if current.Episode != r.Episode {
+		return refuse(store.ErrOperatorConflict)
+	}
+	r = current
+	present, err := inspectRunnerMintToken(c, r)
+	if err != nil {
+		return refuse(err)
+	}
+	if present {
+		return runnerMintRetained(out, stderr, r, e.terminal.now())
+	}
+	remedy := "ailang-worldd service runner reconcile-mint " + id
+	if !e.confirm(stderr, fmt.Sprintf("Revoke orphaned mint %s for action %s and release its owner? Start and worktree are retained.", r.CredentialID, id)) {
+		fmt.Fprintln(stderr, "Retry at a terminal:", remedy)
+		return exitUsage
+	}
+	err = runnerWrite(c.Credentials, func(db *store.Credentials) error {
+		ctx, cancel := budgetContext(time.Second)
+		defer cancel()
+		present, err := inspectRunnerMintToken(c, r)
+		if err != nil {
+			return err
+		}
+		if present {
+			return errors.New("token appeared after review; inspect it again")
+		}
+		return db.AbortOperatorMint(ctx, r, e.terminal.now())
+	})
+	if err != nil {
+		return refuse(err)
+	}
+	fmt.Fprintln(out, "Orphaned mint revoked for", id)
+	return exitOK
+}
+
+func readRunnerMint(path, id string) (store.OperatorMintReceipt, error) {
+	db, err := store.OpenCredentialsReadOnly(runnerDB(path))
+	if err != nil {
+		return store.OperatorMintReceipt{}, err
+	}
+	defer db.Close()
+	ctx, cancel := budgetContext(time.Second)
+	defer cancel()
+	return db.ReadOperatorMint(ctx, id)
+}
+
+func runnerMintRetained(out, stderr io.Writer, r store.OperatorMintReceipt, now int64) int {
+	if now >= r.ExpiresAt {
+		return serviceMessage(stderr, "mint recovery refused", "matching token expired", "ailang-worldd service runner status", 1)
+	}
+	fmt.Fprintln(out, "Matching private token retained for", r.CAS.ID, "(owner retained for file reconciliation)")
+	return exitOK
+}
+
+func inspectRunnerMintToken(c runnerConfig, r store.OperatorMintReceipt) (bool, error) {
+	path := r.TokenPath
+	if r.Store != c.Credentials || !verbs.EpisodeIDPattern.MatchString(r.Episode) ||
+		!filepath.IsAbs(path) || filepath.Clean(path) != path || filepath.Base(path) != r.Episode+".session" {
+		return false, errors.New("mint token path or store identity differs")
+	}
+	if err := runnerPath(c.Home, filepath.Dir(path), c.UID, true, 0700); err != nil {
+		return false, fmt.Errorf("cannot establish private token directory: %w", err)
+	}
+	return readMintToken(path, c.UID, r.CredentialID)
+}
+
+func readMintToken(path string, uid int, credentialID string) (bool, error) {
+	f, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
+	if errors.Is(err, os.ErrNotExist) {
+		return false, nil
+	}
+	if err != nil {
+		return false, errors.New("cannot open mint token safely")
+	}
+	defer f.Close()
+	info, err := f.Stat()
+	if err != nil || !info.Mode().IsRegular() {
+		return false, errors.New("mint token is not a regular file")
+	}
+	stat := info.Sys().(*syscall.Stat_t)
+	if int(stat.Uid) != uid || stat.Nlink != 1 || info.Mode() != 0600 {
+		return false, errors.New("mint token must be private, owned and unlinked")
+	}
+	body, err := io.ReadAll(io.LimitReader(f, 65))
+	if err != nil || len(body) != 64 || fmt.Sprintf("%x", sha256.Sum256(body)) != credentialID {
+		return false, errors.New("mint token bytes differ from receipt")
+	}
+	return true, nil
+}
diff --git a/cmd/ailang-worldd/operator_mint_test.go b/cmd/ailang-worldd/operator_mint_test.go
new file mode 100644
index 0000000..b444e9a
--- /dev/null
+++ b/cmd/ailang-worldd/operator_mint_test.go
@@ -0,0 +1,341 @@
+package main
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"database/sql"
+	"fmt"
+	"io"
+	"os"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+type recoveryFixture struct {
+	*operatorServiceFixture
+	db                    *sql.DB
+	path, token, worktree string
+	receipt               store.OperatorMintReceipt
+}
+
+const reconcileMintExample = "ailang-worldd service runner reconcile-mint oa-0123456789abcdef0123456789abcdef"
+
+func newRecoveryFixture(t *testing.T) recoveryFixture {
+	t.Helper()
+	f := operatorService(t)
+	f.mustInstall(t)
+	f.env.terminal.now = func() int64 { return 1020 }
+	f.writes = nil
+	ep := "minted"
+	sessions := filepath.Join(f.home, "mint-sessions")
+	if err := os.Mkdir(sessions, 0700); err != nil {
+		t.Fatal(err)
+	}
+	in := store.OperatorMintInput{
+		TokenPath: filepath.Join(sessions, ep+".session"),
+		Profile:   approval.ProfileRef{Project: "world", Digest: serviceHash([]byte("profile fixture bytes"))},
+		Grants:    []approval.Grant{{Effect: "fs.read", Scope: "repo", Budget: 2}},
+	}
+	in.Attempt.Catalogue = 1
+	var a approval.OperatorAction
+	fd := claimFixtureEpisodeAction(t, f.config.Credentials, ep, func(action *approval.OperatorAction) {
+		action.Options.Preset = "sol"
+		in.Start = store.EpisodeStart{Schema: "world/episode-start/v1", Episode: ep, Preset: store.StartPreset{Name: "sol"}, TTLS: 20,
+			Task: &store.StartTask{Path: filepath.Join(sessions, ep+".task.md"), SHA256: serviceHash([]byte(action.Task))}}
+		var err error
+		action.Resolution, err = in.Resolution()
+		if err != nil {
+			t.Fatal(err)
+		}
+		a = *action
+		a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
+	})
+	defer fd.Close()
+	db, err := sql.Open("sqlite", runnerDB(f.config.Credentials))
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { db.Close() })
+	r := recoveryFixture{operatorServiceFixture: f, db: db, path: in.TokenPath, worktree: filepath.Join(f.home, "workspaces", ep)}
+	r.exec(t, "UPDATE operator_actions SET state='running' WHERE id=?", a.ID)
+	r.exec(t, `INSERT INTO operator_action_steps (kind,identity,action_id,phase,version,owner,created_at)
+		VALUES ('mint',?,?,'reserved',7,'test-runner',1012)`, a.ID, a.ID)
+	ctx := boundedTestContext(t)
+	err = runnerWrite(f.config.Credentials, func(c *store.Credentials) error {
+		row, err := c.OperatorAction(ctx, a.ID, 1012)
+		if err != nil {
+			return err
+		}
+		h, err := c.OperatorHeartbeat(ctx)
+		if err != nil {
+			return err
+		}
+		in.Attempt.CAS = store.OperatorCAS{ID: row.ID, Digest: [32]byte(row.Digest), Generation: row.Generation}
+		in.Attempt.Runner, in.Attempt.StepVersion = h.Owner, 7
+		r.token, err = authority.MintOperatorStart(ctx, c, c, in, a, row, func() int64 { return 1012 })
+		return err
+	})
+	if err != nil {
+		t.Fatal(err)
+	}
+	r.receipt, err = readRunnerMint(f.config.Credentials, a.ID)
+	if err != nil {
+		t.Fatal(err)
+	}
+	r.exec(t, `INSERT INTO session_credentials (credential_id, episode_id, grants_json, expires_at, created_at)
+		VALUES ('unrelated', ?, '[]', 1100, 1013)`, ep)
+	if err := os.Mkdir(r.worktree, 0700); err != nil {
+		t.Fatal(err)
+	}
+	for _, p := range []string{filepath.Join(r.worktree, "keep"), in.Start.Task.Path, filepath.Join(sessions, ep+".mcp.json")} {
+		if err := os.WriteFile(p, []byte("preserve these bytes"), 0600); err != nil {
+			t.Fatal(err)
+		}
+	}
+	return r
+}
+
+func (r recoveryFixture) exec(t *testing.T, query string, args ...any) {
+	t.Helper()
+	if _, err := r.db.ExecContext(boundedTestContext(t), query, args...); err != nil {
+		t.Fatal(err)
+	}
+}
+
+func (r recoveryFixture) run(t *testing.T, command string) (int, string) {
+	t.Helper()
+	args := strings.Fields(command)
+	if len(args) < 2 || args[0] != "ailang-worldd" || args[1] != "service" {
+		t.Fatalf("unexpected remedy %q", command)
+	}
+	return r.operatorServiceFixture.run(args[2:]...)
+}
+
+func TestOperatorMintRecovery(t *testing.T) {
+	for _, tc := range []struct{ name, want string }{
+		{"revocation", "Orphaned mint revoked"},
+		{"unsafe_parent", "cannot establish private token directory"},
+		{"missing_lock", "episode recovery lock held or unavailable"},
+		{"wrong_receipt_phase", "mint evidence unavailable"},
+		{"unowned", "operator action conflicts"},
+		{"busy", "episode recovery lock held or unavailable"},
+		{"stale_owner", "operator action conflicts"},
+		{"stale_step", "operator action conflicts"},
+		{"appeared", "token appeared after review"},
+		{"adopt", "Matching private token retained"},
+		{"expired", "matching token expired"},
+		{"later_effect", "later effect reserved"},
+		{"wrong_step_owner", "mint evidence unavailable"},
+		{"start_changed", "mint start record differs"},
+		{"action_changed", "operator action conflicts"},
+		{"receipt_changed", "operator action conflicts"},
+		{"missing_credential", "mint start or credential unavailable"},
+		{"abort_rollback", "injected abort failure"},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			r := newRecoveryFixture(t)
+			id := r.receipt.CAS.ID
+			writeToken := func() {
+				t.Helper()
+				if err := writeTokenFile(r.path, r.token); err != nil {
+					t.Fatal(err)
+				}
+			}
+			switch tc.name {
+			case "revocation":
+				r.exec(t, "UPDATE operator_actions SET state='uncertain'")
+			case "action_changed":
+				r.exec(t, `UPDATE operator_actions SET action=CAST(json_set(CAST(action AS TEXT),'$.review','changed') AS BLOB)`)
+			case "unsafe_parent":
+				if err := os.Chmod(filepath.Dir(r.path), 0750); err != nil {
+					t.Fatal(err)
+				}
+			case "missing_lock":
+				if err := os.Remove(operatorEpisodeLockPath(r.config.Credentials, r.receipt.Episode)); err != nil {
+					t.Fatal(err)
+				}
+			case "unowned":
+				r.exec(t, "UPDATE operator_actions SET state='failed', owner=NULL, terminal_at=1019")
+			case "wrong_receipt_phase":
+				r.exec(t, "UPDATE operator_action_steps SET phase='reserved', receipt=NULL WHERE kind='mint'")
+			case "busy":
+				held, err := openRunnerLock(operatorEpisodeLockPath(r.config.Credentials, r.receipt.Episode), 0)
+				if err != nil {
+					t.Fatal(err)
+				}
+				defer held.Close()
+			case "adopt", "expired":
+				writeToken()
+				if tc.name == "expired" {
+					r.env.terminal.now = func() int64 { return r.receipt.ExpiresAt }
+				}
+			case "later_effect":
+				r.exec(t, `INSERT INTO operator_action_steps (kind,identity,action_id,phase,version,owner,created_at)
+					VALUES ('launch',?,?,'reserved',1,'test-runner',1013)`, id, id)
+			case "wrong_step_owner":
+				r.exec(t, "UPDATE operator_action_steps SET owner='other' WHERE kind='mint'")
+			case "start_changed":
+				r.exec(t, "UPDATE episode_starts SET start_json='{}'")
+			case "receipt_changed":
+				r.exec(t, "UPDATE operator_action_steps SET receipt=X'41' WHERE kind='mint'")
+			case "missing_credential":
+				r.exec(t, "DELETE FROM session_credentials WHERE credential_id=?", r.receipt.CredentialID)
+			case "abort_rollback":
+				r.exec(t, "CREATE TRIGGER failed_abort AFTER UPDATE ON operator_actions BEGIN SELECT RAISE(ABORT,'injected abort failure'); END")
+			}
+			before := serviceFileHash(r.config.Credentials)
+			files := provisionFiles(t, filepath.Dir(r.path))
+			worktree := provisionFiles(t, r.worktree)
+			if tc.name == "revocation" {
+				// The returned token is deliberately never published. Invoke the
+				// command printed by the refusal, byte for byte, after a no-tty run.
+				r.env.terminal.openTerminal = noTerminalEnv(1020).openTerminal
+				code, out := r.run(t, reconcileMintExample)
+				_, command, ok := strings.Cut(out, "Retry at a terminal: ")
+				if code != exitUsage || !ok || strings.TrimSpace(command) != reconcileMintExample || !strings.Contains(out, "no controlling terminal") ||
+					serviceFileHash(r.config.Credentials) != before {
+					t.Fatalf("no tty or changed remedy: %d %s", code, out)
+				}
+				r.env.terminal.openTerminal = testMintEnv("n\n", 1020).openTerminal
+				if code, out := r.run(t, strings.TrimSpace(command)); code != exitUsage || !strings.Contains(out, "aborted (not confirmed)") ||
+					serviceFileHash(r.config.Credentials) != before {
+					t.Fatalf("declined recovery: %d %s", code, out)
+				}
+				r.env.terminal.openTerminal = testMintEnv("y\n", 1020).openTerminal
+			}
+			if tc.name == "stale_owner" || tc.name == "stale_step" || tc.name == "appeared" {
+				r.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
+					switch tc.name {
+					case "stale_owner":
+						r.exec(t, "UPDATE operator_actions SET owner='new-owner', generation=generation+1 WHERE id=?", id)
+						r.exec(t, "UPDATE operator_action_steps SET owner='new-owner' WHERE kind='mint'")
+						before = serviceFileHash(r.config.Credentials)
+					case "stale_step":
+						r.exec(t, "UPDATE operator_action_steps SET version=version+1 WHERE kind='mint'")
+						before = serviceFileHash(r.config.Credentials)
+					case "appeared":
+						writeToken()
+						files = provisionFiles(t, filepath.Dir(r.path))
+					}
+					return testMintEnv("y\n", 1020).openTerminal()
+				}
+			}
+			code, out := r.run(t, reconcileMintExample)
+			wantCode := exitUsage
+			if tc.name == "revocation" || tc.name == "adopt" {
+				wantCode = exitOK
+			}
+			if code != wantCode || !strings.Contains(out, tc.want) || strings.Contains(out, r.token) ||
+				!reflect.DeepEqual(files, provisionFiles(t, filepath.Dir(r.path))) || !reflect.DeepEqual(worktree, provisionFiles(t, r.worktree)) ||
+				len(r.writes) != 0 || r.ctl.boots != 0 || r.ctl.stops != 0 {
+				t.Fatalf("recovery %s: %d %s", tc.name, code, out)
+			}
+			if tc.name != "revocation" {
+				if before != serviceFileHash(r.config.Credentials) {
+					t.Fatal("recovery refusal or matched-token read modified the store")
+				}
+				return
+			}
+			var state, owner string
+			var version, generation, minted, other, starts int64
+			var receipt []byte
+			for _, q := range []struct {
+				sql   string
+				slots []any
+			}{
+				{"SELECT state, coalesce(owner,''), generation FROM operator_actions", []any{&state, &owner, &generation}},
+				{"SELECT version, receipt FROM operator_action_steps WHERE kind='mint'", []any{&version, &receipt}},
+				{"SELECT count(*) FROM session_credentials WHERE credential_id != 'unrelated'", []any{&minted}},
+				{"SELECT count(*) FROM session_credentials WHERE credential_id = 'unrelated'", []any{&other}},
+				{"SELECT count(*) FROM episode_starts", []any{&starts}},
+			} {
+				if err := r.db.QueryRowContext(boundedTestContext(t), q.sql).Scan(q.slots...); err != nil {
+					t.Fatal(err)
+				}
+			}
+			if state != "failed" || owner != "" || generation != r.receipt.CAS.Generation+1 || version != r.receipt.StepVersion ||
+				string(receipt) != r.receipt.CredentialID+":"+r.receipt.StartDigest || minted != 0 || other != 1 || starts != 1 {
+				t.Fatalf("orphan disposition changed evidence or unrelated authority: %s %s %d %d %d %d", state, owner, version, minted, other, starts)
+			}
+			requireProvisionLockReleased(t, r.config.Credentials, r.receipt.Episode)
+		})
+	}
+}
+
+func TestOperatorMintTokenReadback(t *testing.T) {
+	for _, name := range []string{"uid", "mode", "link", "symlink", "symlink_valid", "directory", "bytes", "length", "short"} {
+		t.Run(name, func(t *testing.T) {
+			dir := t.TempDir()
+			p := filepath.Join(dir, "token")
+			token := strings.Repeat("a", 64)
+			uid, mode := os.Getuid(), os.FileMode(0600)
+			var err error
+			switch name {
+			case "symlink", "symlink_valid":
+				target := filepath.Join(dir, "target")
+				if name == "symlink_valid" {
+					if err := os.WriteFile(target, []byte(token), 0600); err != nil {
+						t.Fatal(err)
+					}
+				}
+				err = os.Symlink(target, p)
+			case "directory":
+				err = os.Mkdir(p, 0700)
+			default:
+				data := token
+				switch name {
+				case "uid":
+					uid++
+				case "mode":
+					mode = 0644
+				case "bytes":
+					data = strings.Repeat("b", 64)
+				case "length":
+					data += "a"
+				case "short":
+					token = token[:63]
+					data = token
+				}
+				err = os.WriteFile(p, []byte(data), mode)
+				if err == nil && name == "link" {
+					err = os.Link(p, filepath.Join(dir, "other"))
+				}
+			}
+			if err != nil {
+				t.Fatal(err)
+			}
+			before := provisionFiles(t, dir)
+			present, err := readMintToken(p, uid, fmt.Sprintf("%x", sha256.Sum256([]byte(token))))
+			if present || err == nil || !reflect.DeepEqual(before, provisionFiles(t, dir)) {
+				t.Fatalf("unsafe file evidence changed or accepted: present=%v err=%v", present, err)
+			}
+		})
+	}
+}
+
+func TestOperatorMintRecoveryArgs(t *testing.T) {
+	e := serviceEnv{goos: "darwin", label: worldServiceLabel}
+	for _, args := range [][]string{{"reconcile-mint"}, {"reconcile-mint", "wrong"}, {"reconcile-mint", fixtureActionID, "extra"}} {
+		var msg bytes.Buffer
+		if code := runServiceRunner(args, &msg, &msg, e); code != exitUsage || !strings.Contains(msg.String(), "expected reconcile-mint ACTION-ID") {
+			t.Fatalf("recovery syntax: %v: %d %s", args, code, &msg)
+		}
+	}
+}
+
+func TestQuickstartRunnerMintRecovery(t *testing.T) {
+	b, err := os.ReadFile("../../docs/QUICKSTART.md")
+	if err != nil {
+		t.Fatal(err)
+	}
+	const sidecar = "Terminal episode operations use private mode-0600 sidecars named\n`<creds>.operator-episode-<ep>.lock` beside the canonical credentials database.\nKeep these sidecars in place across stop and uninstall."
+	if !strings.Contains(string(b), sidecar) || !strings.Contains(string(b), "\n"+reconcileMintExample+"\n") {
+		t.Fatal("missing the canonical sidecar description or exact attended command")
+	}
+}
diff --git a/cmd/ailang-worldd/operator_service.go b/cmd/ailang-worldd/operator_service.go
index 77974cd..d837ef2 100644
--- a/cmd/ailang-worldd/operator_service.go
+++ b/cmd/ailang-worldd/operator_service.go
@@ -21,8 +21,8 @@ import (
 
 const runnerLabel = "com.sunholo.ailang-world.operator-runner"
 const runnerSystemPath = "/usr/bin:/bin:/usr/sbin:/sbin"
-const runnerHelp = `usage: ailang-worldd service runner install --config FILE | start | stop | status | cancel ACTION-ID | uninstall
-Install, stop, cancel and uninstall require typed y. Start uses the approved configuration.
+const runnerHelp = `usage: ailang-worldd service runner install --config FILE | start | stop | status | cancel ACTION-ID | reconcile-mint ACTION-ID | uninstall
+Install, stop, cancel, orphan mint revocation and uninstall require typed y. Start uses the approved configuration.
 Config JSON keys: home, uid, address, world, credentials, executables (absolute path to sha256 digest), preset_names, tool_dirs, backend ("headless").
 The World service must already be running. Use absolute, canonical paths and include its installed binary in executables.
 `
@@ -272,9 +272,9 @@ func runServiceRunner(args []string, out, errOut io.Writer, e serviceEnv) int {
 		if len(args) != 3 || args[1] != "--config" {
 			return refuse(errors.New("expected install --config FILE"))
 		}
-	case "cancel":
+	case "cancel", "reconcile-mint":
 		if len(args) != 2 || !verbs.ValidOperatorActionID(args[1]) {
-			return refuse(errors.New("expected cancel ACTION-ID"))
+			return refuse(fmt.Errorf("expected %s ACTION-ID", verb))
 		}
 	case "start", "stop", "status", "uninstall":
 		if len(args) != 1 {
@@ -283,6 +283,9 @@ func runServiceRunner(args []string, out, errOut io.Writer, e serviceEnv) int {
 	default:
 		return refuse(errors.New("unknown runner operation"))
 	}
+	if verb == "reconcile-mint" {
+		return e.reconcileRunnerMint(args[1], out, errOut)
+	}
 	p := e.runnerPaths()
 	if verb == "status" {
 		if err := e.runnerStatus(out); err != nil {
diff --git a/cmd/ailang-worldd/service.go b/cmd/ailang-worldd/service.go
index 7d3dff1..3556272 100644
--- a/cmd/ailang-worldd/service.go
+++ b/cmd/ailang-worldd/service.go
@@ -34,7 +34,7 @@ const serviceHelp = `usage: ailang-worldd service install|status|upgrade|stop|st
     [--examples-dir DIR] [--surface-index PATH] [--surface-follow PATH]
     [--surface-origin https://<host>.ts.net[:port]]... [--binary PATH] [--replace]
   status [--json] [--verbose]
-  runner install|start|stop|status|cancel|uninstall (see service runner --help)
+  runner install|start|stop|status|cancel|reconcile-mint|uninstall (see service runner --help)
   upgrade [--binary FILE] [--pause]
 Install, stop and uninstall require a typed y at the controlling terminal.
 Start uses the approved plist; edited plists require confirmation.
diff --git a/docs/QUICKSTART.md b/docs/QUICKSTART.md
index 49659ff..0a17d10 100644
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -83,6 +83,19 @@ Run the checked example with a fake supervisor and disposable stores:
 AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./cmd/ailang-worldd -run '^TestOperatorServiceWriteScope$' -count=1 -v
 ```
 
+For an owned start with a committed mint and a missing token, use the
+attended remedy with its action ID. The disposable example uses this ID.
+
+```sh
+ailang-worldd service runner reconcile-mint oa-0123456789abcdef0123456789abcdef
+```
+
+Review the credential hash and type `y` to revoke that one orphaned
+credential. The start and worktree remain. A matching, unexpired private
+token is retained for file reconciliation with the episode owner still
+held. A conflicting file or later effect reservation keeps the owner in
+place. Run the disposable example with `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./cmd/ailang-worldd -run '^TestOperatorMintRecovery/revocation$' -count=1 -v`.
+
 ## Run your agents through World
 
 World launchers now tell Claude Code, Codex and pi their episode, worktree, session
diff --git a/host/approval/operator_signed.go b/host/approval/operator_signed.go
index 2f96bcb..8808101 100644
--- a/host/approval/operator_signed.go
+++ b/host/approval/operator_signed.go
@@ -22,7 +22,8 @@ type SignedOperatorAction struct {
 	name, format      string
 	credentialID      []byte
 	signCount         uint32
-	confirmation      *OperatorApprovalBinding
+	binding           OperatorApprovalBinding
+	confirmed         bool
 	artifactHash      [32]byte
 }
 
@@ -99,10 +100,9 @@ func operatorProof(d, challenge [32]byte, expiry, now int64, artifact []byte, fo
 		return nil, err
 	}
 	s := &SignedOperatorAction{ok: true, digest: d, challenge: challenge, expiry: expiry, name: ap.Name, format: format}
-	if confirm {
-		s.confirmation = &OperatorApprovalBinding{Name: ap.Name, Format: format, PublicKey: string(ap.Key.Blob)}
-		s.artifactHash = sha256.Sum256(artifact)
-	}
+	s.confirmed = confirm
+	s.binding = OperatorApprovalBinding{Name: ap.Name, Format: format, PublicKey: string(ap.Key.Blob)}
+	s.artifactHash = sha256.Sum256(artifact)
 	if format == FormatSSH {
 		sig, err := ParseSSHSIG(artifact)
 		if err != nil {
@@ -130,10 +130,8 @@ func operatorProof(d, challenge [32]byte, expiry, now int64, artifact []byte, fo
 	if !confirm && ad.SignCount > p.SignCount { // R2
 		return nil, fmt.Errorf("%w: enrolled %d, this assertion %d", ErrOperatorCounter, p.SignCount, ad.SignCount)
 	}
-	if confirm {
-		s.confirmation.PublicKey = string(p.Key.Canonical())
-		s.confirmation.RPID, s.confirmation.Origin, s.confirmation.BackupEligible = p.RPID, p.Origin, p.BackupEligible
-	}
+	s.binding.PublicKey = string(p.Key.Canonical())
+	s.binding.RPID, s.binding.Origin, s.binding.BackupEligible = p.RPID, p.Origin, p.BackupEligible
 	s.credentialID, s.signCount = bytes.Clone(p.CredentialID), ad.SignCount
 	return s, nil
 }
@@ -184,8 +182,34 @@ type OperatorApprovalBinding struct {
 // ConfirmationFor binds the proof to the exact artifact to be stored and
 // distinguishes initial confirmation from a runner recheck.
 func (s *SignedOperatorAction) ConfirmationFor(artifact []byte) (OperatorApprovalBinding, error) {
-	if s == nil || !s.ok || s.confirmation == nil || sha256.Sum256(artifact) != s.artifactHash {
+	if s == nil || !s.ok || !s.confirmed || sha256.Sum256(artifact) != s.artifactHash {
 		return OperatorApprovalBinding{}, ErrNotVerified
 	}
-	return *s.confirmation, nil
+	return s.binding, nil
+}
+
+// OperatorMintProof seals a runner recheck for the start transaction.
+type OperatorMintProof struct{ signed *SignedOperatorAction }
+
+func (s *SignedOperatorAction) MintProof() *OperatorMintProof {
+	if s == nil || !s.ok || s.confirmed {
+		return nil
+	}
+	return &OperatorMintProof{signed: s}
+}
+
+// Bind verifies the exact stored action and artifact and returns the enrollment
+// that the transaction must still find current.
+func (p *OperatorMintProof) Bind(a OperatorAction, challenge [32]byte, expiry int64, artifact []byte) (OperatorApprovalBinding, error) {
+	if p == nil || p.signed == nil || sha256.Sum256(artifact) != p.signed.artifactHash {
+		return OperatorApprovalBinding{}, ErrNotVerified
+	}
+	if err := p.signed.Binds(a, challenge, expiry); err != nil {
+		return OperatorApprovalBinding{}, err
+	}
+	return p.signed.binding, nil
+}
+
+func (p *OperatorMintProof) PasskeyUse() ([]byte, uint32, bool) {
+	return p.signed.PasskeyUse()
 }
diff --git a/host/authority/mint.go b/host/authority/mint.go
index 32e4122..9fd8541 100644
--- a/host/authority/mint.go
+++ b/host/authority/mint.go
@@ -198,3 +198,26 @@ func clampGrantExpiry(grants []broker.Capability, sessionExpiry int64) []broker.
 	}
 	return out
 }
+
+// MintOperatorStart uses only an independently reconstructed action for the
+// runner recheck. The transaction revalidates proof, time and enrollment.
+func MintOperatorStart(ctx context.Context, st, ro *store.Credentials, in store.OperatorMintInput,
+	current approval.OperatorAction, row store.OperatorActionRow, clock func() int64) (string, error) {
+	if len(row.Challenge) != 32 {
+		return "", approval.ErrNotVerified
+	}
+	rechecked, err := RecheckClaimedOperatorAction(ctx, ro, current, [32]byte(row.Challenge), row.ExpiresAt, clock(), row)
+	if err != nil {
+		return "", err
+	}
+	var raw [32]byte
+	if _, err := rand.Read(raw[:]); err != nil {
+		return "", err
+	}
+	token := hex.EncodeToString(raw[:])
+	in.CredentialID = hashHexToken(token)
+	if err := st.CommitOperatorMint(ctx, in, rechecked.OperatorConfirmation().MintProof(), clock); err != nil {
+		return "", err
+	}
+	return token, nil
+}
diff --git a/host/authority/operator_mint_test.go b/host/authority/operator_mint_test.go
new file mode 100644
index 0000000..ac65f34
--- /dev/null
+++ b/host/authority/operator_mint_test.go
@@ -0,0 +1,383 @@
+package authority
+
+import (
+	"context"
+	"crypto/sha256"
+	"errors"
+	"fmt"
+	"math"
+	"os"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+type mintRig struct {
+	operatorStoreRig
+	in  store.OperatorMintInput
+	row store.OperatorActionRow
+	ro  *store.Credentials
+	ctx context.Context
+}
+
+func operatorMintRig(t *testing.T, passkey bool, edit func(*approval.OperatorAction, *store.OperatorMintInput)) mintRig {
+	t.Helper()
+	dir := t.TempDir()
+	in := store.OperatorMintInput{Attempt: store.OperatorMintAttempt{Catalogue: 1}, Grants: []approval.Grant{{Effect: "exec", Scope: "repo", Budget: 4}},
+		Profile: approval.ProfileRef{Project: "world", Digest: hashref.SumSHA256([]byte("profile fixture bytes")).String()}, CredentialID: strings.Repeat("2", 64)}
+	r := newOperatorStoreRigAction(t, passkey, func(a *approval.OperatorAction) {
+		in.TokenPath = filepath.Join(dir, a.Episode+".session")
+		in.Start = store.EpisodeStart{Schema: "world/episode-start/v1", Episode: a.Episode,
+			Preset: store.StartPreset{Name: a.Options.Preset, SHA256: "preset"}, Repo: dir,
+			Base: store.StartBase{Branch: "main", SHA: "base"}, Branch: "new", Agent: "codex", TTLS: 90,
+			Task: &store.StartTask{Path: filepath.Join(dir, a.Episode+".task.md"), SHA256: hashref.SumSHA256([]byte(a.Task)).String()}}
+		var err error
+		a.Resolution, err = in.Resolution()
+		if err != nil {
+			t.Fatal(err)
+		}
+		if edit != nil {
+			edit(a, &in)
+		}
+	})
+	ctx := boundedTestContext(t)
+	r.confirm(t)
+	h := store.OperatorHeartbeat{Owner: store.OperatorOwner{Nonce: "fresh-runner"}, Protocol: 1, Build: "fixture",
+		SeenAt: 1010, Readiness: "start", CatalogueGeneration: in.Attempt.Catalogue}
+	var err error
+	if h.Owner.Generation, err = r.PutOperatorHeartbeat(ctx, h); err != nil {
+		t.Fatal(err)
+	}
+	if err := r.ClaimOperatorAction(ctx, r.cas, h.Owner, 1011); err != nil {
+		t.Fatal(err)
+	}
+	jgExec(t, r.path, `UPDATE operator_actions SET state = 'running'`)
+	jgExec(t, r.path, fmt.Sprintf(`INSERT INTO operator_action_steps
+		(kind, identity, action_id, phase, version, owner, created_at)
+		VALUES ('mint', '%s', '%s', 'reserved', 7, 'fresh-runner', 1012)`, r.cas.ID, r.cas.ID))
+	row, err := r.OperatorAction(ctx, r.cas.ID, 1012)
+	if err != nil {
+		t.Fatal(err)
+	}
+	in.Attempt = store.OperatorMintAttempt{CAS: store.OperatorCAS{ID: row.ID, Digest: [32]byte(row.Digest), Generation: row.Generation},
+		Runner: h.Owner, Catalogue: h.CatalogueGeneration, StepVersion: 7}
+	return mintRig{r, in, row, operatorStoreOpen(t, r.path, store.OpenCredentialsReadOnly), ctx}
+}
+
+func (r mintRig) seal(t *testing.T) *approval.OperatorMintProof {
+	t.Helper()
+	p, err := RecheckClaimedOperatorAction(r.ctx, r.ro, r.action, [32]byte(r.row.Challenge), r.row.ExpiresAt, 1013, r.row)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return p.OperatorConfirmation().MintProof()
+}
+
+type mintSnapshot struct {
+	db               [32]byte
+	action           store.OperatorActionRow
+	sessions         []store.SessionRow
+	start            *store.EpisodeStart
+	token, task, mcp bool
+}
+
+func (r mintRig) snapshot(t *testing.T) mintSnapshot {
+	t.Helper()
+	b, err := os.ReadFile(r.path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	row, err := r.ro.OperatorAction(r.ctx, r.row.ID, 1015)
+	if err != nil {
+		t.Fatal(err)
+	}
+	sessions, err := r.ro.ListSessions(r.ctx, "", 20)
+	if err != nil {
+		t.Fatal(err)
+	}
+	start, err := r.ro.EpisodeStart(r.ctx, r.row.Episode)
+	if err != nil {
+		t.Fatal(err)
+	}
+	exists := func(p string) bool {
+		_, err := os.Lstat(p)
+		if err != nil && !os.IsNotExist(err) {
+			t.Fatal(err)
+		}
+		return err == nil
+	}
+	var token string
+	for _, fact := range r.action.Resolution {
+		if fact.Name == "mint_token_path" {
+			token = fact.Value
+		}
+	}
+	dir := filepath.Dir(token)
+	return mintSnapshot{sha256.Sum256(b), row, sessions, start, exists(token), exists(filepath.Join(dir, r.action.Episode+".task.md")),
+		exists(filepath.Join(dir, r.action.Episode+".mcp.json"))}
+}
+
+func mintAt(now int64) func() int64 { return func() int64 { return now } }
+
+func (r mintRig) unchanged(t *testing.T, before mintSnapshot, err, want error) {
+	t.Helper()
+	if !errors.Is(err, want) || !reflect.DeepEqual(before, r.snapshot(t)) {
+		t.Fatalf("refusal %v, want %v with exact state unchanged", err, want)
+	}
+}
+
+func TestOperatorMintBinding(t *testing.T) {
+	for _, tc := range []struct {
+		name, sql string
+		key       bool
+		change    func(*store.OperatorMintInput)
+		now       int64
+		class     error
+		field     string
+	}{
+		{name: "ttl", change: func(in *store.OperatorMintInput) { in.Start.TTLS++ }, field: "mint_start"},
+		{name: "base", change: func(in *store.OperatorMintInput) { in.Start.Base.SHA += "0" }, field: "mint_start"},
+		{name: "grant_effect", change: func(in *store.OperatorMintInput) { in.Grants[0].Effect += "0" }, field: "mint_grants"},
+		{name: "grant_scope", change: func(in *store.OperatorMintInput) { in.Grants[0].Scope += "0" }, field: "mint_grants"},
+		{name: "grant_budget", change: func(in *store.OperatorMintInput) { in.Grants[0].Budget++ }, field: "mint_grants"},
+		{name: "project", change: func(in *store.OperatorMintInput) { in.Profile.Project += "0" }, field: "mint_project"},
+		{name: "profile", change: func(in *store.OperatorMintInput) {
+			in.Profile.Digest = hashref.SumSHA256([]byte("modified profile bytes")).String()
+		}, field: "mint_profile"},
+		{name: "action_id", change: func(in *store.OperatorMintInput) { in.Attempt.CAS.ID = "unrelated" }, class: store.ErrOperatorMissing},
+		{name: "generation", change: func(in *store.OperatorMintInput) { in.Attempt.CAS.Generation++ }, class: store.ErrOperatorConflict},
+		{name: "digest", change: func(in *store.OperatorMintInput) { in.Attempt.CAS.Digest[0] ^= 1 }, class: store.ErrOperatorConflict},
+		{name: "runner_generation", change: func(in *store.OperatorMintInput) { in.Attempt.Runner.Generation++ }, class: store.ErrOperatorRunner},
+		{name: "runner_owner", sql: "UPDATE operator_runner_state SET owner = 'another'", class: store.ErrOperatorRunner},
+		{name: "catalogue_generation", sql: "UPDATE operator_runner_state SET catalogue_generation = 2", class: store.ErrOperatorRunner},
+		{name: "catalogue_review", sql: "UPDATE operator_runner_state SET catalogue_generation = 2", change: func(in *store.OperatorMintInput) { in.Attempt.Catalogue++ }, field: "mint_catalogue"},
+		{name: "action_owner", sql: "UPDATE operator_actions SET owner = 'different'; UPDATE operator_action_steps SET owner = 'different' WHERE kind='mint'", class: store.ErrOperatorConflict},
+		{name: "action_phase", sql: "UPDATE operator_actions SET state = 'claimed'", class: store.ErrOperatorConflict},
+		{name: "approver_identity", sql: "UPDATE operator_actions SET approved_by = 'ssh:different'", class: store.ErrOperatorConflict},
+		{name: "different_challenge", sql: "UPDATE operator_actions SET challenge=zeroblob(32)", field: "challenge"},
+		{name: "step_phase", sql: "UPDATE operator_action_steps SET phase='observed', receipt=X'41' WHERE kind='mint'", class: store.ErrOperatorConflict},
+		{name: "step_version", change: func(in *store.OperatorMintInput) { in.Attempt.StepVersion++ }, class: store.ErrOperatorConflict},
+		{name: "step_owner", sql: "UPDATE operator_action_steps SET owner='another' WHERE kind='mint'", class: store.ErrOperatorConflict},
+		{name: "step_identity", sql: "UPDATE operator_action_steps SET identity='another' WHERE kind='mint'", class: store.ErrOperatorConflict},
+		{name: "episode_row", sql: "UPDATE operator_actions SET episode_id='another'", class: store.ErrOperatorConflict},
+		{name: "record_episode", change: func(in *store.OperatorMintInput) { in.Start.Episode = "another" }, field: "mint_start"},
+		{name: "artifact", sql: "UPDATE operator_actions SET artifact = X'41'", class: approval.ErrNotVerified},
+		{name: "revoked", sql: "UPDATE approvers SET revoked_at = 1014", class: store.ErrOperatorConflict},
+		{name: "future_enrolment", sql: "UPDATE approvers SET enrolled_at = 1015, last_used_at = NULL WHERE format = 'ssh'", class: store.ErrOperatorConflict},
+		{name: "approver_name", sql: "UPDATE approvers SET name='renamed' WHERE format='ssh'", class: store.ErrOperatorConflict},
+		{name: "approver_format", sql: "UPDATE approvers SET format='webauthn' WHERE format='ssh'", class: store.ErrOperatorConflict},
+		{name: "replaced_key", sql: "UPDATE approvers SET public_key = X'41' WHERE format = 'ssh'", class: store.ErrOperatorConflict},
+		{name: "passkey_credential", key: true, sql: "UPDATE webauthn_credentials SET credential_id = X'41'", class: store.ErrOperatorConflict},
+		{name: "passkey_counter", key: true, sql: "UPDATE webauthn_credentials SET sign_count = 1", class: store.ErrOperatorConflict},
+		{name: "passkey_rp", key: true, sql: "UPDATE webauthn_credentials SET rp_id = 'different'", class: store.ErrOperatorConflict},
+		{name: "passkey_origin", key: true, sql: "UPDATE webauthn_credentials SET origin = 'https://different'", class: store.ErrOperatorConflict},
+		{name: "passkey_backup", key: true, sql: "UPDATE webauthn_credentials SET backup_eligible = 1 - backup_eligible", class: store.ErrOperatorConflict},
+		{name: "expiry_boundary", now: 1600, class: store.ErrOperatorExpired},
+		{name: "clock_behind", now: 999, class: store.ErrOperatorConflict},
+		{name: "overflow", change: func(in *store.OperatorMintInput) { in.Start.TTLS = math.MaxInt64 }, class: store.ErrOperatorInput},
+		{name: "zero_ttl", change: func(in *store.OperatorMintInput) { in.Start.TTLS = 0 }, class: store.ErrOperatorInput},
+		{name: "no_grants", change: func(in *store.OperatorMintInput) { in.Grants = nil }, class: store.ErrOperatorInput},
+		{name: "profile_blank", change: func(in *store.OperatorMintInput) { in.Profile.Project = "" }, class: store.ErrOperatorInput},
+		{name: "credential_upper", change: func(in *store.OperatorMintInput) { in.CredentialID = strings.Repeat("A", 64) }, class: store.ErrOperatorInput},
+		{name: "profile_hash", change: func(in *store.OperatorMintInput) { in.Profile.Digest = strings.Repeat("1", 64) }, class: store.ErrOperatorInput},
+		{name: "credential_hash", change: func(in *store.OperatorMintInput) { in.CredentialID = "sha" }, class: store.ErrOperatorInput},
+		{name: "existing_branch", change: func(in *store.OperatorMintInput) { in.Start.BranchExisted = true }, class: store.ErrOperatorConflict},
+		{name: "resume_history", change: func(in *store.OperatorMintInput) { in.Start.Resumes = []store.EpisodeResume{{}} }, class: store.ErrOperatorConflict},
+		{name: "launch_history", change: func(in *store.OperatorMintInput) { in.Start.Launches = []store.EpisodeLaunch{{}} }, class: store.ErrOperatorConflict},
+		{name: "preset_name", change: func(in *store.OperatorMintInput) { in.Start.Preset.Name += "0" }, class: store.ErrOperatorConflict},
+		{name: "missing_task", change: func(in *store.OperatorMintInput) { in.Start.Task = nil }, class: store.ErrOperatorConflict},
+		{name: "task_digest", change: func(in *store.OperatorMintInput) { in.Start.Task.SHA256 += "0" }, class: store.ErrOperatorConflict},
+		{name: "relative_token", change: func(in *store.OperatorMintInput) {
+			in.TokenPath = filepath.Base(in.TokenPath)
+			in.Start.Task.Path = filepath.Base(in.Start.Task.Path)
+		}, class: store.ErrOperatorConflict},
+		{name: "unclean_token", change: func(in *store.OperatorMintInput) {
+			in.TokenPath = filepath.Dir(in.TokenPath) + "/./" + filepath.Base(in.TokenPath)
+		}, class: store.ErrOperatorConflict},
+		{name: "token_fact", field: "mint_token_path"},
+		{name: "wrong_token", change: func(in *store.OperatorMintInput) { in.TokenPath += "0" }, class: store.ErrOperatorConflict},
+		{name: "wrong_task", change: func(in *store.OperatorMintInput) { in.Start.Task.Path += "0" }, class: store.ErrOperatorConflict},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			r := operatorMintRig(t, tc.key, func(a *approval.OperatorAction, _ *store.OperatorMintInput) {
+				if tc.name == "token_fact" {
+					a.Resolution[len(a.Resolution)-1].Value += ".different"
+				}
+			})
+			seal := r.seal(t)
+			if tc.sql != "" {
+				jgExec(t, r.path, tc.sql)
+			}
+			before := r.snapshot(t)
+			if tc.change != nil {
+				tc.change(&r.in)
+			}
+			now := tc.now
+			if now == 0 {
+				now = 1014
+			}
+			err := r.CommitOperatorMint(r.ctx, r.in, seal, mintAt(now))
+			if tc.field != "" {
+				var mismatch *approval.BindingError
+				if !errors.As(err, &mismatch) || mismatch.Field != tc.field || !reflect.DeepEqual(before, r.snapshot(t)) {
+					t.Fatalf("expected binding %s and unchanged state, got %v", tc.field, err)
+				}
+			} else {
+				r.unchanged(t, before, err, tc.class)
+			}
+		})
+	}
+	t.Run("proof", func(t *testing.T) {
+		r := operatorMintRig(t, false, nil)
+		before := r.snapshot(t)
+		for _, signed := range []*approval.SignedOperatorAction{nil, {}, r.proof.OperatorConfirmation()} {
+			if signed.MintProof() != nil {
+				t.Fatal("non-runner proof acquired a mint seal")
+			}
+		}
+		for _, proof := range []*approval.OperatorMintProof{nil, {}} {
+			r.unchanged(t, before, r.CommitOperatorMint(r.ctx, r.in, proof, mintAt(1014)), approval.ErrNotVerified)
+		}
+	})
+	t.Run("start_options", func(t *testing.T) {
+		r := operatorMintRig(t, false, func(a *approval.OperatorAction, in *store.OperatorMintInput) {
+			a.Options.Headless = false
+			a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
+		})
+		r.unchanged(t, r.snapshot(t), r.CommitOperatorMint(r.ctx, r.in, r.seal(t), mintAt(1014)), store.ErrOperatorConflict)
+	})
+	t.Run("read_only", func(t *testing.T) {
+		r := operatorMintRig(t, false, nil)
+		r.unchanged(t, r.snapshot(t), r.ro.CommitOperatorMint(r.ctx, r.in, r.seal(t), mintAt(1014)), store.ErrOperatorConflict)
+	})
+	t.Run("thin_challenge", func(t *testing.T) {
+		r := operatorMintRig(t, false, nil)
+		before := r.snapshot(t)
+		r.row.Challenge = nil
+		raw, err := MintOperatorStart(r.ctx, r.Credentials, r.ro, r.in, r.action, r.row, mintAt(1014))
+		if raw != "" {
+			t.Fatal("short challenge returned a token")
+		}
+		r.unchanged(t, before, err, approval.ErrNotVerified)
+	})
+	t.Run("changed_review", func(t *testing.T) {
+		r := operatorMintRig(t, false, nil)
+		before := r.snapshot(t)
+		other := r.action
+		other.Review += " changed"
+		raw, err := MintOperatorStart(r.ctx, r.Credentials, r.ro, r.in, other, r.row, mintAt(1014))
+		var judge *JudgeError
+		if raw != "" || !errors.As(err, &judge) || !reflect.DeepEqual(before, r.snapshot(t)) {
+			t.Fatal("changed effective action reached mint")
+		}
+	})
+	t.Run("fresh_clock", func(t *testing.T) {
+		r := operatorMintRig(t, false, nil)
+		before, calls := r.snapshot(t), 0
+		raw, err := MintOperatorStart(r.ctx, r.Credentials, r.ro, r.in, r.action, r.row, func() int64 {
+			calls++
+			if calls == 1 {
+				return 1013
+			}
+			return r.row.ExpiresAt
+		})
+		if raw != "" || calls != 2 {
+			t.Fatal("expected distinct precheck and transactional clock readings, with no token")
+		}
+		r.unchanged(t, before, err, store.ErrOperatorExpired)
+	})
+}
+
+func TestOperatorMintCollision(t *testing.T) {
+	for _, kind := range []string{"credential", "start"} {
+		t.Run(kind, func(t *testing.T) {
+			r := operatorMintRig(t, false, nil)
+			if kind == "credential" {
+				jgExec(t, r.path, fmt.Sprintf("INSERT INTO session_credentials VALUES ('old', '%s', '[]', 1002, 1001, NULL, NULL)", r.action.Episode))
+			} else {
+				jgExec(t, r.path, fmt.Sprintf("INSERT INTO episode_starts VALUES ('%s', '{}')", r.action.Episode))
+			}
+			before := r.snapshot(t)
+			token, err := MintOperatorStart(r.ctx, r.Credentials, r.ro, r.in, r.action, r.row, mintAt(1014))
+			if token != "" || err == nil || !reflect.DeepEqual(before, r.snapshot(t)) ||
+				kind == "start" && !strings.Contains(err.Error(), "UNIQUE constraint") ||
+				kind == "credential" && !errors.Is(err, store.ErrOperatorConflict) {
+				t.Fatalf("collision %s accepted or changed state: %v", kind, err)
+			}
+		})
+	}
+}
+
+func TestOperatorMintRollback(t *testing.T) {
+	for _, tc := range []struct{ name, sql, want string }{
+		{"session_insert", "CREATE TRIGGER fail BEFORE INSERT ON session_credentials BEGIN SELECT RAISE(ABORT, 'injected session failure'); END", "injected session failure"},
+		{"start_insert", "CREATE TRIGGER fail BEFORE INSERT ON episode_starts BEGIN SELECT RAISE(ABORT, 'injected start failure'); END", "injected start failure"},
+		{"receipt_update", "CREATE TRIGGER fail AFTER UPDATE ON operator_action_steps BEGIN SELECT RAISE(ABORT, 'injected receipt failure'); END", "injected receipt failure"},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			r := operatorMintRig(t, false, nil)
+			jgExec(t, r.path, tc.sql)
+			before := r.snapshot(t)
+			raw, err := MintOperatorStart(r.ctx, r.Credentials, r.ro, r.in, r.action, r.row, mintAt(1014))
+			if raw != "" || err == nil || !strings.Contains(err.Error(), tc.want) || !reflect.DeepEqual(before, r.snapshot(t)) {
+				t.Fatalf("fault %s did not roll back all writes: %v", tc.name, err)
+			}
+			jgExec(t, r.path, "DROP TRIGGER fail")
+			if _, err := MintOperatorStart(r.ctx, r.Credentials, r.ro, r.in, r.action, r.row, mintAt(1015)); err != nil {
+				t.Fatalf("verified rollback did not permit bounded retry: %v", err)
+			}
+		})
+	}
+}
+
+func TestOperatorMintReceiptWindow(t *testing.T) {
+	t.Run("transaction_separation", TestOperatorMintRollback)
+	r := operatorMintRig(t, true, nil)
+	before, _, err := r.ro.Approver(r.ctx, strings.TrimPrefix(r.row.ApprovedBy, "webauthn:"))
+	if err != nil {
+		t.Fatal(err)
+	}
+	pk, _, err := r.ro.Passkey(r.ctx, before.Name)
+	if err != nil {
+		t.Fatal(err)
+	}
+	token, err := MintOperatorStart(r.ctx, r.Credentials, r.ro, r.in, r.action, r.row, mintAt(1014))
+	if err != nil || len(token) != 64 {
+		t.Fatalf("mint: %v", err)
+	}
+	// Model a lost handoff by retaining only its hash, before any file writer can run.
+	credential := fmt.Sprintf("%x", sha256.Sum256([]byte(token)))
+	if err := r.Close(); err != nil {
+		t.Fatal(err)
+	}
+	db := operatorStoreOpen(t, r.path, store.OpenCredentials)
+	receipt, err := db.ReadOperatorMint(r.ctx, r.row.ID)
+	if err != nil || receipt.CredentialID != credential || receipt.StepVersion != 8 || receipt.Owner != r.row.Owner || receipt.ExpiresAt != 1104 {
+		t.Fatalf("commit did not preserve exact receipt: %+v %v", receipt, err)
+	}
+	got := r.snapshot(t)
+	if got.token || got.task || got.mcp || len(got.sessions) != 1 || got.sessions[0].CredentialID != credential ||
+		got.sessions[0].CreatedAt != 1014 || got.sessions[0].ExpiresAt != 1104 || got.sessions[0].Project != r.in.Profile.Project ||
+		got.sessions[0].ProfileDigest != r.in.Profile.Digest || got.sessions[0].GrantsJSON != `[{"Effect":"exec","Scope":"repo","ExpiresAt":1104,"Budget":4}]` ||
+		got.start == nil || got.start.StartedAt != "1970-01-01T00:16:54Z" || got.start.TTLS != 90 || !reflect.DeepEqual(got.action, r.row) {
+		t.Fatalf("commit did not preserve exact start, owner and credential: %+v", got)
+	}
+	after, _, err := r.ro.Approver(r.ctx, before.Name)
+	if err != nil || !reflect.DeepEqual(before, after) {
+		t.Fatal("mint advanced approver use")
+	}
+	pkAfter, _, err := r.ro.Passkey(r.ctx, before.Name)
+	if err != nil || !reflect.DeepEqual(pk, pkAfter) {
+		t.Fatal("mint advanced the passkey counter")
+	}
+	token, err = MintOperatorStart(r.ctx, db, r.ro, r.in, r.action, r.row, mintAt(1015))
+	if token != "" {
+		t.Fatal("committed orphan returned a second token")
+	}
+	r.unchanged(t, got, err, store.ErrOperatorConflict)
+}
diff --git a/host/authority/operator_store_test.go b/host/authority/operator_store_test.go
index b62c11b..ad435be 100644
--- a/host/authority/operator_store_test.go
+++ b/host/authority/operator_store_test.go
@@ -22,6 +22,10 @@ type operatorStoreRig struct {
 }
 
 func newOperatorStoreRig(t *testing.T, passkey bool) operatorStoreRig {
+	return newOperatorStoreRigAction(t, passkey, func(*approval.OperatorAction) {})
+}
+
+func newOperatorStoreRigAction(t *testing.T, passkey bool, prepare func(*approval.OperatorAction)) operatorStoreRig {
 	t.Helper()
 	path, _ := jgRig(t)
 	ctx := boundedTestContext(t)
@@ -32,6 +36,7 @@ func newOperatorStoreRig(t *testing.T, passkey bool) operatorStoreRig {
 	}
 	a := ojAction()
 	a.Store = canonical
+	prepare(&a)
 	req := store.OperatorRequest{Action: a, OperationID: "1234567890abcdef1234567890abcdef", CreatedAt: 1000, ExpiresAt: 1600, Challenge: sha256.Sum256([]byte(t.Name()))}
 	row, err := c.RequestOperatorAction(ctx, req, 1000)
 	if err != nil {
diff --git a/host/store/credentials.go b/host/store/credentials.go
index 657c7e5..9d8bb6a 100644
--- a/host/store/credentials.go
+++ b/host/store/credentials.go
@@ -498,16 +498,9 @@ func (c *Credentials) mintSession(ctx context.Context, row SessionRow, start *Ep
 		return err
 	}
 	if start != nil {
-		if start.Episode != row.EpisodeID || start.Schema != "world/episode-start/v1" {
-			return fmt.Errorf("store: mint session: invalid start record")
-		}
-		data, err := json.Marshal(start)
-		if err != nil {
+		if _, err := insertEpisodeStart(ctx, tx, row.EpisodeID, *start); err != nil {
 			return err
 		}
-		if _, err := tx.ExecContext(ctx, `INSERT INTO episode_starts (episode, start_json) VALUES (?, ?)`, row.EpisodeID, string(data)); err != nil {
-			return fmt.Errorf("store: mint start: %w", err)
-		}
 	}
 	if resume != nil {
 		if err := appendResume(ctx, tx, row.EpisodeID, *resume); err != nil {
@@ -520,6 +513,20 @@ func (c *Credentials) mintSession(ctx context.Context, row SessionRow, start *Ep
 	return nil
 }
 
+func insertEpisodeStart(ctx context.Context, tx *sql.Tx, episode string, start EpisodeStart) ([]byte, error) {
+	if start.Episode != episode || start.Schema != "world/episode-start/v1" {
+		return nil, fmt.Errorf("store: mint session: invalid start record")
+	}
+	data, err := json.Marshal(start)
+	if err != nil {
+		return nil, err
+	}
+	if _, err := tx.ExecContext(ctx, `INSERT INTO episode_starts (episode, start_json) VALUES (?, ?)`, episode, string(data)); err != nil {
+		return nil, fmt.Errorf("store: mint start: %w", err)
+	}
+	return data, nil
+}
+
 // insertSessionRow writes one session_credentials row inside tx. An unbound row names only the
 // five original columns, so the same statement works on a version-1 store; the binding columns
 // stay NULL. mintSession and MintSessionForRequest both insert through here.
diff --git a/host/store/credentials_test.go b/host/store/credentials_test.go
index f39d869..0fab29a 100644
--- a/host/store/credentials_test.go
+++ b/host/store/credentials_test.go
@@ -308,6 +308,9 @@ func TestCredentialsMethodsRequireADeadline(t *testing.T) {
 	_, errs["PendingVerbRequests"] = c.PendingVerbRequests(ctx, 1000, 0)
 	errs["ConsumeVerbRequest"] = c.ConsumeVerbRequest(ctx, "x", "y", ApproverUse{}, "", nil, "", 1000)
 	_, errs["OperatorAction"] = c.OperatorAction(ctx, "x", 1000)
+	errs["CommitOperatorMint"] = c.CommitOperatorMint(ctx, OperatorMintInput{}, nil, nil)
+	_, errs["ReadOperatorMint"] = c.ReadOperatorMint(ctx, "x")
+	errs["AbortOperatorMint"] = c.AbortOperatorMint(ctx, OperatorMintReceipt{}, 1000)
 	_, errs["OperatorHeartbeat"] = c.OperatorHeartbeat(ctx)
 	_, errs["ReadOperatorHeartbeat"] = ReadOperatorHeartbeat(ctx, "")
 	_, errs["ReadOperatorEpisodeOwner"] = ReadOperatorEpisodeOwner(ctx, c.display, "ep")
diff --git a/host/store/operator_mint.go b/host/store/operator_mint.go
new file mode 100644
index 0000000..1f7eee4
--- /dev/null
+++ b/host/store/operator_mint.go
@@ -0,0 +1,324 @@
+package store
+
+import (
+	"bytes"
+	"context"
+	"crypto/sha256"
+	"database/sql"
+	"encoding/hex"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"math"
+	"path/filepath"
+	"strconv"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+)
+
+// OperatorMintAttempt identifies an already reserved mint step.
+type OperatorMintAttempt struct {
+	CAS         OperatorCAS
+	Runner      OperatorOwner
+	Catalogue   int64
+	StepVersion int64
+}
+
+// OperatorMintInput is the effective mint reconstructed by the caller. Clock
+// is read after the transaction starts. No token bytes enter the store.
+type OperatorMintInput struct {
+	Attempt      OperatorMintAttempt
+	Start        EpisodeStart
+	Grants       []approval.Grant
+	Profile      approval.ProfileRef
+	CredentialID string
+	TokenPath    string
+}
+
+// Resolution supplies the mint fields of the reviewed action and is repeated
+// on the effective inputs at commit. StartedAt comes from the commit clock.
+func (in OperatorMintInput) Resolution() ([]approval.OperatorResolved, error) {
+	start := in.Start
+	start.StartedAt = ""
+	body, err := json.Marshal(start)
+	if err != nil {
+		return nil, err
+	}
+	grants, _ := json.Marshal(in.Grants) // approval.Grant contains only strings and an integer.
+	return []approval.OperatorResolved{
+		{Name: "mint_catalogue", Value: strconv.FormatInt(in.Attempt.Catalogue, 10)},
+		{Name: "mint_grants", Value: string(grants)},
+		{Name: "mint_profile", Value: in.Profile.Digest},
+		{Name: "mint_project", Value: in.Profile.Project},
+		{Name: "mint_start", Value: fmt.Sprintf("%x", sha256.Sum256(body))},
+		{Name: "mint_token_path", Value: in.TokenPath},
+	}, nil
+}
+
+func mintFact(a approval.OperatorAction, name string) string {
+	for _, fact := range a.Resolution {
+		if fact.Name == name {
+			return fact.Value
+		}
+	}
+	return ""
+}
+
+func operatorMintApprover(ctx context.Context, tx *sql.Tx, proof *approval.OperatorMintProof, b approval.OperatorApprovalBinding, now int64) error {
+	var current bool
+	if err := tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM approvers WHERE
+		name = ? AND format = ? AND public_key = ? AND revoked_at IS NULL AND enrolled_at <= ?)`,
+		b.Name, b.Format, []byte(b.PublicKey), now).Scan(&current); err != nil {
+		return err
+	}
+	if !current {
+		return ErrOperatorConflict
+	}
+	if id, count, ok := proof.PasskeyUse(); ok {
+		if err := tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM webauthn_credentials WHERE name = ?
+			AND credential_id = ? AND rp_id = ? AND origin = ? AND backup_eligible = ? AND sign_count >= ?)`,
+			b.Name, id, b.RPID, b.Origin, b.BackupEligible, count).Scan(&current); err != nil {
+			return err
+		}
+		if !current {
+			return ErrOperatorConflict
+		}
+	}
+	return nil
+}
+
+// CommitOperatorMint atomically records the session, start, and observed mint
+// receipt. It is an inert transaction seam pending paired-process admission.
+func (c *Credentials) CommitOperatorMint(ctx context.Context, in OperatorMintInput, proof *approval.OperatorMintProof, clock func() int64) error {
+	if err := c.operatorAccess(ctx); err != nil {
+		return err
+	}
+	if c.lock == nil {
+		return ErrOperatorConflict
+	}
+	tx, err := c.db.BeginTx(ctx, nil)
+	if err != nil {
+		return err
+	}
+	defer tx.Rollback()
+	now := clock()
+	try := in.Attempt
+	row, err := operatorCASRow(ctx, tx, try.CAS)
+	if err != nil {
+		return err
+	}
+	if now >= row.ExpiresAt {
+		return ErrOperatorExpired
+	}
+	if now < row.CreatedAt || row.State != "running" || row.Owner != try.Runner.Nonce || row.Kind != "start" {
+		return ErrOperatorConflict
+	}
+	var current bool
+	err = tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM operator_runner_state WHERE
+		generation = ? AND owner = ? AND catalogue_generation = ?)`,
+		try.Runner.Generation, try.Runner.Nonce, try.Catalogue).Scan(&current)
+	if err != nil {
+		return err
+	}
+	if !current {
+		return ErrOperatorRunner
+	}
+	a, err := decodeOperatorMintAction(row)
+	if err != nil {
+		return err
+	}
+	if a.Store != c.lock.dbPath || !a.Options.Headless {
+		return ErrOperatorConflict
+	}
+	b, err := proof.Bind(a, [32]byte(row.Challenge), row.ExpiresAt, row.Artifact)
+	if err != nil {
+		return err
+	}
+	if row.ApprovedBy != b.Format+":"+b.Name {
+		return ErrOperatorConflict
+	}
+	if err := operatorMintApprover(ctx, tx, proof, b, now); err != nil {
+		return err
+	}
+	session, err := operatorMintSession(in, a, now)
+	if err != nil {
+		return err
+	}
+	err = tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM session_credentials WHERE episode_id = ?)`, a.Episode).Scan(&current)
+	if err != nil {
+		return err
+	}
+	if current {
+		return ErrOperatorConflict
+	}
+	if err := insertSessionRow(ctx, tx, session); err != nil {
+		return err
+	}
+	start := in.Start
+	start.StartedAt = time.Unix(now, 0).UTC().Format(time.RFC3339)
+	startJSON, err := insertEpisodeStart(ctx, tx, a.Episode, start)
+	if err != nil {
+		return err
+	}
+	receipt := in.CredentialID + fmt.Sprintf(":%x", sha256.Sum256(startJSON))
+	res, err := tx.ExecContext(ctx, `UPDATE operator_action_steps SET phase = 'observed', version = version + 1, receipt = ?
+		WHERE kind = 'mint' AND identity = ? AND action_id = ? AND phase = 'reserved' AND version = ? AND owner = ?`,
+		[]byte(receipt), a.ID, a.ID, try.StepVersion, row.Owner)
+	if err := operatorChangedOne(res, err); err != nil {
+		return err
+	}
+	return tx.Commit()
+}
+
+func decodeOperatorMintAction(row OperatorActionRow) (approval.OperatorAction, error) {
+	if sha256.Sum256(row.Action) != [32]byte(row.Digest) {
+		return approval.OperatorAction{}, ErrOperatorConflict
+	}
+	a, err := approval.DecodeOperatorAction(row.Action)
+	if err != nil {
+		return approval.OperatorAction{}, err
+	}
+	if a.ID != row.ID || a.Episode != row.Episode || a.Kind != row.Kind {
+		return approval.OperatorAction{}, ErrOperatorConflict
+	}
+	return a, nil
+}
+
+func operatorMintSession(in OperatorMintInput, a approval.OperatorAction, now int64) (SessionRow, error) {
+	start := in.Start
+	if start.BranchExisted || len(start.Resumes) != 0 || len(start.Launches) != 0 ||
+		start.Preset.Name != a.Options.Preset || start.Task == nil ||
+		start.Task.SHA256 != hashref.SumSHA256([]byte(a.Task)).String() {
+		return SessionRow{}, ErrOperatorConflict
+	}
+	if !filepath.IsAbs(in.TokenPath) || filepath.Clean(in.TokenPath) != in.TokenPath ||
+		filepath.Base(in.TokenPath) != a.Episode+".session" ||
+		start.Task.Path != filepath.Join(filepath.Dir(in.TokenPath), a.Episode+".task.md") {
+		return SessionRow{}, ErrOperatorConflict
+	}
+	if _, profileErr := hashref.Parse(in.Profile.Digest); start.TTLS <= 0 || now > math.MaxInt64-start.TTLS || len(in.Grants) == 0 ||
+		in.Profile.Project == "" || profileErr != nil || !operatorDigestString(in.CredentialID) {
+		return SessionRow{}, ErrOperatorInput
+	}
+	facts, err := in.Resolution()
+	if err != nil {
+		return SessionRow{}, err
+	}
+	for _, fact := range facts {
+		if value := mintFact(a, fact.Name); value != fact.Value {
+			return SessionRow{}, &approval.BindingError{Field: fact.Name}
+		}
+	}
+	expiry := now + start.TTLS
+	grants := make([]approval.FactGrant, len(in.Grants))
+	for i, grant := range in.Grants {
+		grants[i] = approval.FactGrant{Effect: grant.Effect, Scope: grant.Scope, Budget: grant.Budget, ExpiresAt: expiry}
+	}
+	wire, _ := json.Marshal(grants) // approval.FactGrant contains only strings and integers.
+	return SessionRow{CredentialID: in.CredentialID, EpisodeID: a.Episode, CreatedAt: now, ExpiresAt: expiry,
+		GrantsJSON: string(wire), Project: in.Profile.Project, ProfileDigest: in.Profile.Digest}, nil
+}
+
+func operatorDigestString(s string) bool {
+	b, err := hex.DecodeString(s)
+	return err == nil && len(b) == sha256.Size && hex.EncodeToString(b) == s
+}
+
+// OperatorMintReceipt is redacted evidence for attended recovery. CAS and
+// step version pin the exact action and receipt read before the terminal prompt.
+type OperatorMintReceipt struct {
+	CAS                              OperatorCAS
+	Owner, Episode, Store, TokenPath string
+	StepVersion, ExpiresAt           int64
+	CredentialID, StartDigest        string
+}
+
+type operatorMintReader interface {
+	QueryRowContext(context.Context, string, ...any) *sql.Row
+}
+
+func (c *Credentials) ReadOperatorMint(ctx context.Context, id string) (OperatorMintReceipt, error) {
+	if err := c.operatorAccess(ctx); err != nil {
+		return OperatorMintReceipt{}, err
+	}
+	return readOperatorMint(ctx, c.db, id)
+}
+
+func readOperatorMint(ctx context.Context, db operatorMintReader, id string) (OperatorMintReceipt, error) {
+	row, err := scanOperatorAction(db.QueryRowContext(ctx, `SELECT `+operatorColumns+` FROM operator_actions WHERE id = ?`, id))
+	if err != nil {
+		return OperatorMintReceipt{}, err
+	}
+	if row.Kind != "start" || row.State != "running" && row.State != "uncertain" {
+		return OperatorMintReceipt{}, ErrOperatorConflict
+	}
+	a, err := decodeOperatorMintAction(row)
+	if err != nil {
+		return OperatorMintReceipt{}, err
+	}
+	r := OperatorMintReceipt{CAS: OperatorCAS{ID: id, Digest: [32]byte(row.Digest), Generation: row.Generation},
+		Owner: row.Owner, Episode: a.Episode, Store: a.Store, TokenPath: mintFact(a, "mint_token_path")}
+	var receipt []byte
+	err = db.QueryRowContext(ctx, `SELECT version, receipt FROM operator_action_steps WHERE kind = 'mint'
+		AND identity = ? AND action_id = ? AND phase = 'observed' AND owner = ?
+		AND NOT EXISTS (SELECT 1 FROM operator_action_steps WHERE action_id = ? AND kind NOT IN ('operation','worktree','mint'))`,
+		id, id, row.Owner, id).Scan(&r.StepVersion, &receipt)
+	if err != nil {
+		return OperatorMintReceipt{}, fmt.Errorf("mint evidence unavailable or later effect reserved: %w", err)
+	}
+	cred, startDigest, _ := bytes.Cut(receipt, []byte(":"))
+	r.CredentialID, r.StartDigest = string(cred), string(startDigest)
+	if !operatorDigestString(r.CredentialID) || !operatorDigestString(r.StartDigest) {
+		return OperatorMintReceipt{}, ErrOperatorConflict
+	}
+	var start []byte
+	err = db.QueryRowContext(ctx, `SELECT s.start_json, c.expires_at FROM episode_starts s JOIN session_credentials c
+		ON c.episode_id = s.episode WHERE s.episode = ? AND c.credential_id = ?`, a.Episode, r.CredentialID).Scan(&start, &r.ExpiresAt)
+	if err != nil {
+		return OperatorMintReceipt{}, fmt.Errorf("mint start or credential unavailable: %w", err)
+	}
+	if fmt.Sprintf("%x", sha256.Sum256(start)) != r.StartDigest {
+		return OperatorMintReceipt{}, errors.New("mint start record differs from receipt")
+	}
+	return r, nil
+}
+
+// AbortOperatorMint revokes only the receipt's credential and releases the
+// exact reviewed owner after the attended caller proves no token was installed.
+func (c *Credentials) AbortOperatorMint(ctx context.Context, reviewed OperatorMintReceipt, now int64) error {
+	if err := c.operatorAccess(ctx); err != nil {
+		return err
+	}
+	if c.lock == nil || reviewed.Store != c.lock.dbPath {
+		return ErrOperatorConflict
+	}
+	tx, err := c.db.BeginTx(ctx, nil)
+	if err != nil {
+		return err
+	}
+	defer tx.Rollback()
+	current, err := readOperatorMint(ctx, tx, reviewed.CAS.ID)
+	if err != nil {
+		return err
+	}
+	if current != reviewed {
+		return ErrOperatorConflict
+	}
+	res, err := tx.ExecContext(ctx, `DELETE FROM session_credentials WHERE credential_id = ? AND episode_id = ?`,
+		current.CredentialID, current.Episode)
+	if err := operatorChangedOne(res, err); err != nil {
+		return err
+	}
+	res, err = tx.ExecContext(ctx, `UPDATE operator_actions SET state = 'failed', owner = NULL, generation = generation + 1, terminal_at = ?
+		WHERE id = ? AND digest = ? AND generation = ? AND owner = ? AND state IN ('running','uncertain')
+		AND EXISTS (SELECT 1 FROM operator_action_steps WHERE kind = 'mint' AND identity = ?
+			AND action_id = ? AND phase = 'observed' AND version = ? AND owner = ?)`,
+		now, current.CAS.ID, current.CAS.Digest[:], current.CAS.Generation, current.Owner,
+		current.CAS.ID, current.CAS.ID, current.StepVersion, current.Owner)
+	if err := operatorChangedOne(res, err); err != nil {
+		return err
+	}
+	return tx.Commit()
+}
````
