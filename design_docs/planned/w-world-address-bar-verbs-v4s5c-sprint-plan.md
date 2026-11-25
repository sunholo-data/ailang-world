# Sprint plan — w-world-address-bar-verbs-v4, sprint S5c (iteration 433, row 243)

> **Scope.** M1 ships killing tests for five of the seven judge-432 survivors and for the two S5 parity survivors. It changes tests only. M2 ships the full live-input reconstruction of a signed headless start (design §5.3 `TestOperatorStartChanged`) over the S5 shared read seam. That reconstruction locks, writes, mints and launches nothing, and it has no production caller. Worktree reservation and its attended disposition, file publication and the paired-process consumer/drain are cut to S5d and later, with measured sizes in Appendix C. No new `uncertain` state ships in this sprint.
>
> **Authority.** Reviewed design r1 §§2.2, 3.1–3.3, 4.2, 5.3, 5.5, 6–8 (quorum r2 PROCEED 3/3, no absent seats; not reopened). D-WORLD-119=A and D-WORLD-120=A, attended `82f05004`. Landed: S1 #606, S2 #609, S2b #611, S3 #643, S4 #651, S4b #654, S5 #657, S5b #659 (merge `5d3b0581`, sonnet PASS 84, 0 blocking). Base `606f6e2c57c79e98975153f055e213df143638d0` = `5d3b0581` + record merge #660 (design_docs only, C1).

## 0. Binding constraints

- The controller commits the two plan records, then applies M1 and M2 in order with `/usr/bin/patch -p1 -E`, and owns every Git write. The planner built product changes only in archive stages under `$HOME/.ailang/state/world-iter433/plan-scratch/`. The planner made no Git writes and did not touch the shared main checkout.
- **Activation stays closed.** This sprint adds no production caller of `authority.MintOperatorStart`, `CommitOperatorMint`, `ClaimOperatorAction` or any reserve/observe API. The new `reconstructOperatorStart` and `operatorStartFacts` are reachable only from tests. §6 `callers` proves this with a positive control (`authority.RecheckVerb` at `episode_finish_surface.go:239`) in the same leg. The seven-survivor prerequisite in directive item 1 applies only to a sprint that gives the mint a production caller. M1 still discharges five of the seven now, and Appendix C carries the other two to the consumer sprint.
- **Attended resolver gate.** M2 introduces no new uncertain outcome. A stale reconstruction is a refusal with no effect, so it needs no resolver. Worktree and file uncertainty travel with their attended commands in S5d (Appendix C).
- Tests use disposable real migrations, the repository test signer, real `git` in `t.TempDir()` repositories, a stubbed daemon health and injected clocks. No test writes outside `t.TempDir()`. No new test is added to `host/daemon` or `host/verifygate`.

## 1. Pre-flight

### 1.1 Baseline at 606f6e2c

`git archive 606f6e2c | tar -x -C "$S/base"` exported the pinned tree. §6's gate ran on it before any stage with `GATES_GIT_REF="$W" GATES_SKIP_VERIFY_GO=1 bash "$S/gates.sh" "$S/base" "$S/gate-base" base 606f6e2c57c79e98975153f055e213df143638d0`. It ran outside any sandbox, so loopback, the whole CLI package and the read-only real observations are informative. `verify_go` is SKIPPED in archive runs because the archive has no Git metadata (S5b §1.1 F5). The controller runs it in the committed checkout. The base leg results are in §6 next to the final run.

### 1.2 Current inputs and coupling inventory

`W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter433-row243-v4s5c`, `S=$HOME/.ailang/state/world-iter433/plan-scratch`, `B=606f6e2c57c79e98975153f055e213df143638d0`. Commands run in `W` with `GIT_OPTIONAL_LOCKS=0` unless a stage is named.

| ID | Command | Observed output and consequence |
|---|---|---|
| C1 | `git rev-parse HEAD; git diff --stat 5d3b0581 606f6e2c -- cmd host docs app \| wc -l` | `606f6e2c57c79e98975153f055e213df143638d0`; `0`. The product is S5b's landed tree. |
| C2 | `git grep -n 'MintOperatorStart(' -- '*.go' \| grep -v _test.go; git grep -n '\.CommitOperatorMint(\|\.AbortOperatorMint(\|\.ReadOperatorMint(' -- '*.go' \| grep -v _test.go; git grep -n 'authority.RecheckVerb(' -- '*.go' \| grep -v _test.go` | Only the definition `host/authority/mint.go:204`. `CommitOperatorMint` at `host/authority/mint.go:219` only. `AbortOperatorMint`/`ReadOperatorMint` at `cmd/ailang-worldd/operator_mint.go:64,81` only (attended resolver). Positive control `cmd/ailang-worldd/episode_finish_surface.go:239`. The mint has no production caller (directive item 1 premise confirmed). |
| C3 | `git grep -n 'resolvePresetStart(' -- '*.go' \| grep -v _test.go; git grep -n 'provisionEpisode(' -- 'cmd/*.go' \| grep -v _test.go` | Read seam definition `episode_preset.go:79`, one caller `episode_preset.go:328`. Provision definition `session_new.go:405`; callers `approve.go:398`, `episode.go:377`, `episode_preset.go:365`, `episode_resume.go:215`, `session_new.go:660`. M2 adds a second read-seam caller in `operator_start_input.go` and touches no provision caller. |
| C4 | `git grep -n 'worktree\.\(Add\|AddExisting\|Remove\)(' -- 'cmd/*.go' \| grep -v _test.go; git grep -n 'writeTokenFile(\|writeMCPConfig(' -- 'cmd/*.go' \| grep -v _test.go` | Worktree creation `session_new.go:510,512`, removal `:577`. Token writer `session_new.go:549` (def `:674`). MCP writer def `episode.go:552`, calls `episode.go:373`, `episode_launch.go:500`, `episode_preset.go:358`, `mint_request.go:96`. Preset task copy `episode_preset.go:342`. S5c touches none of these. They are S5d's integration points (Appendix C). |
| C5 | `git grep -n 'syscall\.Flock.*LOCK_NB' -- cmd host \| grep -v _test.go \| wc -l; git grep -n '\.\(ClaimOperatorAction\|ReserveOperatorStep\|ObserveOperatorStep\|RecheckClaimedOperatorAction\)(' -- cmd host \| grep -v _test.go` | `7` (unchanged, includes the writer-lock comment). Zero production claim/reserve/observe/recheck sites. `ls host/store/operator_steps.go`: absent. The S3b bank is still unlanded. |
| C6 | `git grep -n 'row.Kind != "start"\|a.Store != c.lock.dbPath' host/store/operator_mint.go; git grep -n 'r.Store != c.Credentials\|current.Episode != r.Episode\|IsRegular' cmd/ailang-worldd/operator_mint.go; git grep -n 'agent: req.start.Agent\|overrides = append(overrides, override.name)' cmd/ailang-worldd/episode_preset.go` | Survivor anchors `operator_mint.go:115,132` (store); resolver `:38,94,114`, token basename `:95`, reread `:41`; parity `episode_preset.go:130,375`. |
| C7 | `git grep -n 'func \(TestDaemonDependencyAllowlist\|TestWorldBoundaryDependencyAllowlist\|TestProductionGoSurface\|TestProductionContextRoots\|TestCLICredentialWriterSites\|TestEveryVerbHasHelp\|TestQuickstartServiceSectionMatchesTheCLIHelp\)(' -- '*_test.go'` | `host/daemon/daemon_test.go:899`, `host/boundary/allowlist_world_test.go:758`, `host/store/context_roots_test.go:234,129`, `cmd/ailang-worldd/cli_credential_writers_test.go:18`, `help_test.go:21`, `quickstart_service_test.go:15`. |
| C8 | `"$HOME/.pinned-ailang/ailang" version; go version` | Pinned AILANG v0.41.0 (S5b C9). go1.26.6 darwin/arm64. Linux CI ubuntu-24.04. |

| Inventory | Disposition (measured on `st-M2`) |
|---|---|
| CLI commands, help, flags, QUICKSTART goldens | Unchanged. No new verb, flag, runner operation or help text. No docs change. |
| Writer sites and locks | `TestCLICredentialWriterSites` PASS unchanged (accept). No new open, writer or `Flock` (C5 count stays 7, checked in §6 `callers`). |
| Context roots and Go surface | `go test -v -run '^(TestProductionContextRoots\|TestProductionGoSurface)$' ./host/store`: `root census: 54 named roots, 327 production files`, `Go surface: host/cmd=327 allow-listed reproducers=6`. The base has 326 production files, and the one new production file adds one. Both tests count dynamically and pass. No new context root. The reconstruction takes the caller's context. |
| Dependencies and imports | The new file imports only `context`, `encoding/json`, `host/approval` and `host/store`, all already in `cmd/ailang-worldd`. `TestWorldBoundaryDependencyAllowlist`, `TestBoundaryASTWriteGuard` and `TestEveryCommandDirectoryIsAProtectedGroup` PASS on `st-M2`. `TestDaemonDependencyAllowlist` passes in §6 `inv-host-daemon`. |
| Daemon routes and served behaviour | Unchanged. §6 `routes` compares `daemon.go`, `surface_page.go` and `app.js` byte-for-byte with base. |

### 1.3 Decisions

- **Cut by measurement.** The compiled S5c candidate (M1 91 + M2 310) measures 401 lines. Re-applying the iteration-432 S5c probes to `st-M2` measured worktree 59, files 53 and fd 72, plus the S3b bank 132, all compiled and green (Appendix C). The full worktree reservation, attended worktree disposition, publication and consumer are still unmeasured beyond these probes. Each one needs an injected failure at every step boundary, plus its own attended remedy and mutant matrix. Folding any of them into this sprint would put a new uncertain outcome in the same patch as unmeasured recovery. This sprint ships the two complete, effect-free milestones. Each is independently green and safe at its commit.
- **Reconstruction shape.** `operatorStartFacts` runs the S5 read seam with the signed preset, the signed episode, `launching`/`headless` forced true and the **signed task bytes** as the only task source. It projects the read into `store.OperatorMintInput` and twelve name-ordered facts: the six `mint_*` facts the mint transaction already rechecks, plus `launch_binary`, `launch_wall`, `mcp_path`, `served_build`, `served_interpreter` and `worktree_path`. `reconstructOperatorStart` requires exact list equality with the signed resolution. A count difference is `BindingError{resolution}`. A value difference is `BindingError{<name>}`. A served-store difference is `BindingError{store}`. A non-start or non-headless action is `ErrOperatorConflict`, and a profile-less preset is `ErrOperatorInput`. A change that would lower authority is still a stale review (design §4.2 step 2).
- **Survivor dispositions.** `E-TX-STORE-ID` is pinned by minting into a `VACUUM INTO` copy of the store. `E-TX-KIND` is pinned by refusal order: a finish action is a conflict before runner evaluation. With a current runner the mutant is equivalent, because approval validation forbids a headless finish (`host/approval/operator_action.go:151`) and the headless guard refuses it. `E-RES-STORE-ID` and `E-RES-TOKENPATH-BASE` are pinned by rewriting the stored action and its digest together. `E-RES-REGULAR` is pinned by the portable type class (FIFO and directory are refused as a type, before mode). `E-RES-EPISODE-CONFLICT` and `E-RES-REREAD` require an interleaving between the first read and the lock. They survive (§4) and go to the paired-process sprint (Appendix C). This plan claims no test for them.

### 1.4 Premises found false or incomplete at HEAD

| ID | Measurement and consequence |
|---|---|
| F1 | Directive item 1 lists `E-TX-KIND` as an unconditional survivor. Measured: the mutant is killed only by refusal ordering (no current runner). With a current runner, approval validation O8b (`operator_action.go:151`, "a finish has no headless") plus the transaction's `!a.Options.Headless` guard make it equivalent. `TestOperatorMintScope/finish_kind` pins the ordering, and §4 records both facts. |
| F2 | `E-RES-EPISODE-CONFLICT` and `E-RES-REREAD` cannot be killed by a single-process test without a production interleaving hook. The resolver reads, opens the lock and re-reads with no injectable step. For the only effect path, `AbortOperatorMint` already compares `current != reviewed` inside the transaction (`host/store/operator_mint.go:306`). The residual difference is the retained-token message under a race. Measured SURVIVED with the M1 tests in place (§4). |
| F3 | A profile-less preset with Workspace.Exec grants is refused by the seam itself (`preset.go:441–442`, "Workspace.Exec requires a preset profile"). M2's `req.prof == nil` branch is reachable only with a non-exec grant set and no sandboxed checks. `TestOperatorStartChanged/no_profile` uses `Workspace.Read=worktree:1` with no checks. |
| F4 | Daemon/runner build drift is not visible to the S5 seam alone. `checkDaemon` (`episode.go:465`) accepts any daemon whose build equals this binary's build, so a paired upgrade between review and execution passes it. `served_build` makes that a stale review (`TestOperatorStartChanged/build`). |
| F5 | The S5b plan's Appendix C names "full live-input reconstruction of every signed field". At HEAD, no producer of signed resolution facts exists (the catalogue is still inert, S4). M2 defines the reconstruction projection, so the future catalogue must sign exactly `operatorStartFacts`'s list. `TestOperatorStartFacts` pins the names and independently derived values, and `Validate()` proves the projection is a signable action. |

### 1.5 Size and CI race cost

Count with `git diff -M --numstat 606f6e2c HEAD -- cmd host docs app` after each commit, and report the sum of both columns. The planner measured `git diff --no-index -M --numstat "$S/base/<d>" "$S/st-<M>/<d>"` for `d` in `cmd host` (`docs` and `app` unchanged: `diff -rq` empty / absent). Every product line is inside the filter, so the unfiltered product sum equals the filtered one. The two plan records add their own line counts to an unfiltered diff.

| Milestone | Added | Deleted | Cumulative | Increment | Snap STATUS / DELETED / manifest |
|---|---:|---:|---:|---:|---|
| M1 | 90 | 1 | 91 | 91 | 3/0/3 |
| M2 | 400 | 1 | 401 | 310 | 5/0/5 |

Per file at M2: `episode_preset_read_test.go` 8, `operator_mint_test.go` 43 (42/1), `operator_start_input.go` 81, `operator_start_input_test.go` 229, `host/authority/operator_mint_test.go` 40. The touched packages are `cmd/ailang-worldd` and `host/authority`. The new tests add about 6 s to the CLI package (each `TestOperatorStartChanged` arm builds one preset fixture in about 0.3 s) and about 0.02 s to authority. They contain no sleeps.

## 2. How the executor applies a sub-milestone

````bash
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter433-row243-v4s5c
S=$HOME/.ailang/state/world-iter433/plan-scratch
T=$HOME/.ailang/state/world-iter433/executor
B=606f6e2c57c79e98975153f055e213df143638d0
PLAN=$W/design_docs/planned/w-world-address-bar-verbs-v4s5c-sprint-plan.md
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

Repeat from patch extraction with `M=M2` after M1 is committed, or after M1 acceptance and snapshot if the controller stages cumulatively. Each patch is incremental: M1 applies to the base and M2 to the M1 tree. Keep the two plan records out of the executor snapshot (the snapshot pathspec is `cmd host docs app`).

## 3. Milestones

### M1 — survivor and parity pins, tests only (91 cumulative)

`TestOperatorMintScope` (`host/authority`) covers `store_copy` and `finish_kind`. In `store_copy`, a sealed, reserved mint is committed into a `VACUUM INTO` copy of the store. It must refuse as a conflict, mint no session or start in the copy and leave the original exactly unchanged. In `finish_kind`, a confirmed, running, owned finish action with a reserved mint step and no runner row must refuse as a conflict. Three new arms in `cmd/ailang-worldd` cover the resolver. `TestOperatorMintRecovery/store_identity` and `/token_name` rewrite the stored action and its digest together, so only `inspectRunnerMintToken`'s store and basename checks can refuse. Each arm requires the exact refusal text, unchanged store bytes, files, worktree, zero writes and zero controller calls. `TestOperatorMintTokenReadback/fifo` and the strengthened `/directory` arm require the refusal text `mint token is not a regular file`. That is a portable type class: a FIFO's and a directory's `Mode()` carry type bits on both Darwin and Linux, and `O_RDONLY|O_NONBLOCK` opens a FIFO on both. `TestPresetStartReadParity` asserts the persisted override list `[agent branch ttl]` against an independent slice. `TestPresetStartReadHeadless` runs the attended Codex start and requires its printed ` codex exec --ignore-user-config ` launch line. M1 changes no production file.

### M2 — live-input reconstruction of a signed headless start (401 cumulative)

`cmd/ailang-worldd/operator_start_input.go` adds `operatorStartFacts` and `reconstructOperatorStart` (§1.3). `TestOperatorStartFacts` builds the review from a first live read, as the catalogue would sign it. It requires `Validate()` to pass on the resulting action, plus the exact fact-name list and each non-mint value derived independently of the projection: the resolved `codex` path, `1h0m0s`, an empty MCP path for headless Codex, `json(testBuild)`, the fixture daemon's `ToolRef` and the worktree path. The test also requires the effective input's start, task path and digest, token path, profile, grants (against `sessionNewGrants("se-tools", nil, 50)`) and attempt. Its mint facts must equal `in.Resolution()`.

`TestOperatorStartChanged` changes one live input after review per arm: preset TTL, preset grants, base ref advance, profile file and its daemon reload, signed task bytes, catalogue generation, paired daemon/runner build, interpreter fact, a different agent binary on PATH, MCP and worktree path facts, an extra or missing fact, served store, finish kind, interactive mode and a profile-less preset. Each arm requires the exact field or class. Every arm, including `unchanged` (which must return the reviewed input exactly), reads back `observeProvision` (credentials DB, branch refs, Git registration, workspace and session files) and asserts that the episode lock, the worktree and the fake agent's `ran` marker are all absent.

````bash accept.sh
#!/bin/bash
set -euo pipefail
W=$1; M=$2; PART=$3
export AILANG_BIN="$HOME/.pinned-ailang/ailang"
export WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
S="$HOME/.ailang/state/world-iter433/plan-scratch"
export AILANG_STATE_DIR="$S/accept-state" AILANG_CACHE_DIR="$S/accept-cache"
mkdir -p "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR"
cd "$W"
test "$PART" = fast
run() {
  local pkg=$1 names=$2 name log
  log=$(mktemp "$S/accept.XXXXXX")
  go test -v -count=1 -timeout 2m "./$pkg" -run "^($names)\$" | tee "$log"
  while IFS= read -r name; do
    grep -Fq -- "--- PASS: $name (" "$log"
  done < <(printf '%s\n' "$names" | tr '|' '\n')
}
case "$M" in
  M1|M2) ;;
  *) exit 2 ;;
esac
run cmd/ailang-worldd 'TestOperatorMintRecovery|TestOperatorMintTokenReadback|TestOperatorMintRecoveryArgs|TestQuickstartRunnerMintRecovery|TestPresetStartReadParity|TestPresetStartReadHeadless|TestPresetStartReadInteractive|TestCLICredentialWriterSites'
run host/authority 'TestOperatorMintScope|TestOperatorMintBinding|TestOperatorMintCollision|TestOperatorMintRollback|TestOperatorMintReceiptWindow|TestOperatorRunnerRechecks'
if [ "$M" = M2 ]; then
  run cmd/ailang-worldd 'TestOperatorStartFacts|TestOperatorStartChanged'
fi
printf 'ACCEPT_OK %s %s\n' "$M" "$PART"
````

Planner rehearsal: `bash "$S/accept.sh" "$S/st-M1" M1 fast` and `bash "$S/accept.sh" "$S/st-M2" M2 fast` both ended `ACCEPT_OK` (`$S/accept-M1.log`, `$S/accept-M2.log`). Both also passed again on the round-trip exports (§7).

### 3.1 Branch traces

| Input or fault | Decision and readback |
|---|---|
| Signed action is a finish, or a start without `--headless` | `ErrOperatorConflict` before the seam runs. Nothing is read, locked or written. |
| Signed preset, episode or task drives the seam | The seam runs with exactly the signed preset name, the signed episode, launch+headless and the signed task bytes. A changed user task file is never read. |
| Seam refusal: existing branch, worktree, token/task/MCP path, daemon store/build/tool mismatch, launch dir/lock, absent agent | Returned unchanged from the S5 seam (S5 tests). No effect. |
| Preset without a profile | `ErrOperatorInput`, because the mint requires a profile binding. No effect. |
| Served store differs from the signed store | `BindingError{store}`. No effect. |
| Preset TTL/bytes, base ref, task bytes changed | `BindingError{mint_start}`. No effect. |
| Preset grants changed | `BindingError{mint_grants}` (sorted before `mint_start`). No effect. |
| Profile digest changed (daemon reloaded) | `BindingError{mint_profile}`. No effect. |
| Catalogue generation changed | `BindingError{mint_catalogue}`. No effect. |
| Paired daemon+runner build changed | `BindingError{served_build}`. No effect. |
| Agent binary path, interpreter, MCP path or worktree path differs | `BindingError{<that fact>}`. No effect. |
| Signed resolution has an extra or missing fact | `BindingError{resolution}`. No effect. |
| Every fact equal | Returns the effective `OperatorMintInput` for the inert mint. Still no effect. Nothing calls it in production. |
| Mint into a copied store, finish-kind mint | `ErrOperatorConflict`, with the original store unchanged and no copy session or start (M1). |
| Resolver: rewritten store identity or token basename | `mint token path or store identity differs`, with store, files and worktree unchanged and no writes or controller calls (M1). |
| Token is a FIFO or directory | `mint token is not a regular file` (portable type class), with the files unchanged (M1). |

### 3.2 Design §5.3 row assignment

| Design row | Disposition |
|---|---|
| `TestOperatorStartChanged` | **M2, complete for the reconstruction.** Preset/base/profile/tool/task/daemon/catalogue/paths changes each refuse with no effect. Activation must call `reconstructOperatorStart` before the worktree reservation (S5d). |
| `TestOperatorStartCollisions` | The seam already refuses an existing episode worktree, branch and token/task/MCP path (S5 tests), and M2 inherits those refusals unchanged. The reservation-time absence proofs, Git registration and no-adoption cases are S5d. |
| `TestOperatorStartOnce` | S5d (one worktree and file set) and S6 (receipt-before-spawn, confinement). |
| `TestOperatorWorktreeWindow` | S5d, with the attended worktree disposition command, effect-specific readback and owner/step CAS in the same patch. |
| `TestOperatorFilesReconcile` | S5e (no-overwrite publication and attended reconciliation). |
| `TestOperatorMintBinding`, `Rollback`, `ReceiptWindow` | Landed in S5b. M1 adds `TestOperatorMintScope` for the two transaction scope guards. |
| `TestOperatorLaunchSession`, `SpawnRefused`, `LaunchCrashWindow`, `LaunchExit` | S6. |
| `TestOperatorStartBrowser` | S6/S8. |
| `TestOperatorHerdrUnavailable`, `HerdrHandoff`, `HerdrWindow` | S12. |
| §5.5 `TestOperatorServiceUpgrade/active_child` and the three other process pins | S5f/S6 consumer (Appendix C). Activation gate. |

## 4. Mutation table

The ledger has **40 / 40 named kills** for the rows this plan names against a test. The two interleaving survivors appear separately below, with no claimed test, so there are 42 specs in `$S/mutants.json`. Each mutant was enumerated from the M2 diff (one per new or changed production branch and per projected value) or from the judge-432/S5 survivor list for M1. `$S/mutate.py` (the iteration-432 replayer) applied each edit in a disposable copy of its stage with the §3 pinned environment and `go test -count=1 -timeout 1m -v -run '^TEST$/^subtest$' ./PACKAGE`. The specs are in `$S/mutants.json` and the results in `$S/mutation-ledger.jsonl` (first pass `$S/mut-b_6obnfi`, retargeted pass `$S/mut-v7sr2vnq`). The **Package** column names the package holding the killing test. For `E-TX-*` the production file is in `host/store` and the test in `host/authority`, so they were run against `host/authority` (iteration-432 retro rule). The whole-package column is the `go test -count=1 -timeout 8m -v ./PACKAGE` verdict. §4.1 gives the command for every row still pending. In the first pass, four M2 rows did not compile (`X-STORE`, `X-MINT-FACTS`, `X-BUILD`, `X-RETURN` left a variable unused). They were rewritten to compile and re-run. Four rows first survived their named arm: `X-GRANT-EFFECT` and `X-GRANT-SCOPE` (the test asserted only grant count), and `X-INTERPRETER` and `X-WORKTREE` (the changed-fact arms pin the comparison, not the value source). The test was strengthened (grants compared with `sessionNewGrants`) or retargeted to `TestOperatorStartFacts`'s independent values, and the rows were re-run killed.

| ID | Production anchor | Package of killing test | Named behavioural test | Named | Whole package |
|---|---|---|---|---|---|
| E-TX-STORE-ID | `host/store/operator_mint.go:132` | `host/authority` | `TestOperatorMintScope/store_copy` | KILLED | KILLED |
| E-TX-KIND | `host/store/operator_mint.go:115` | `host/authority` | `TestOperatorMintScope/finish_kind` | KILLED | KILLED |
| E-RES-STORE-ID | `cmd/ailang-worldd/operator_mint.go:94` | `cmd/ailang-worldd` | `TestOperatorMintRecovery/store_identity` | KILLED | §4.1 |
| E-RES-TOKENPATH-BASE | `cmd/ailang-worldd/operator_mint.go:95` | `cmd/ailang-worldd` | `TestOperatorMintRecovery/token_name` | KILLED | §4.1 |
| E-RES-REGULAR | `cmd/ailang-worldd/operator_mint.go:114` | `cmd/ailang-worldd` | `TestOperatorMintTokenReadback/fifo` | KILLED | §4.1 |
| E-RES-REGULAR-DIR | `cmd/ailang-worldd/operator_mint.go:114` | `cmd/ailang-worldd` | `TestOperatorMintTokenReadback/directory` | KILLED | §4.1 |
| J-PRINT-AGENT | `cmd/ailang-worldd/episode_preset.go:375` | `cmd/ailang-worldd` | `TestPresetStartReadHeadless` | KILLED | §4.1 |
| J-DOUBLE-OVERRIDE-APPEND | `cmd/ailang-worldd/episode_preset.go:130` | `cmd/ailang-worldd` | `TestPresetStartReadParity` | KILLED | §4.1 |
| X-KIND | `CLI/operator_start_input.go:16` | `cmd/ailang-worldd` | `TestOperatorStartChanged/finish` | KILLED | §4.1 |
| X-HEADLESS | `CLI/operator_start_input.go:16` | `cmd/ailang-worldd` | `TestOperatorStartChanged/interactive` | KILLED | §4.1 |
| X-PRESET | `CLI/operator_start_input.go:19` | `cmd/ailang-worldd` | `TestOperatorStartFacts` | KILLED | §4.1 |
| X-EPISODE | `CLI/operator_start_input.go:20` | `cmd/ailang-worldd` | `TestOperatorStartFacts` | KILLED | §4.1 |
| X-HEADLESS-INPUT | `CLI/operator_start_input.go:20` | `cmd/ailang-worldd` | `TestOperatorStartFacts` | KILLED | §4.1 |
| X-TASK-BYTES | `CLI/operator_start_input.go:21` | `cmd/ailang-worldd` | `TestOperatorStartChanged/task` | KILLED | §4.1 |
| X-PROF-NIL | `CLI/operator_start_input.go:27` | `cmd/ailang-worldd` | `TestOperatorStartChanged/no_profile` | KILLED | §4.1 |
| X-PROF-CLASS | `CLI/operator_start_input.go:28` | `cmd/ailang-worldd` | `TestOperatorStartChanged/no_profile` | KILLED | §4.1 |
| X-STORE | `CLI/operator_start_input.go:34` | `cmd/ailang-worldd` | `TestOperatorStartChanged/store` | KILLED | §4.1 |
| X-STORE-FIELD | `CLI/operator_start_input.go:35` | `cmd/ailang-worldd` | `TestOperatorStartChanged/store` | KILLED | §4.1 |
| X-TASK-PATH | `CLI/operator_start_input.go:38` | `cmd/ailang-worldd` | `TestOperatorStartFacts` | KILLED | §4.1 |
| X-ATTEMPT | `CLI/operator_start_input.go:40` | `cmd/ailang-worldd` | `TestOperatorStartFacts` | KILLED | §4.1 |
| X-TOKEN-PATH | `CLI/operator_start_input.go:40` | `cmd/ailang-worldd` | `TestOperatorStartFacts` | KILLED | §4.1 |
| X-PROJECT | `CLI/operator_start_input.go:41` | `cmd/ailang-worldd` | `TestOperatorStartFacts` | KILLED | §4.1 |
| X-PROFILE-DIGEST | `CLI/operator_start_input.go:41` | `cmd/ailang-worldd` | `TestOperatorStartFacts` | KILLED | §4.1 |
| X-GRANTS | `CLI/operator_start_input.go:43` | `cmd/ailang-worldd` | `TestOperatorStartFacts` | KILLED | §4.1 |
| X-GRANT-EFFECT | `CLI/operator_start_input.go:43` | `cmd/ailang-worldd` | `TestOperatorStartFacts` | KILLED | §4.1 |
| X-GRANT-SCOPE | `CLI/operator_start_input.go:43` | `cmd/ailang-worldd` | `TestOperatorStartFacts` | KILLED | §4.1 |
| X-GRANT-BUDGET | `CLI/operator_start_input.go:43` | `cmd/ailang-worldd` | `TestOperatorStartChanged/preset_grants` | KILLED | §4.1 |
| X-LAUNCH-BINARY | `CLI/operator_start_input.go:53` | `cmd/ailang-worldd` | `TestOperatorStartFacts` | KILLED | §4.1 |
| X-LAUNCH-WALL | `CLI/operator_start_input.go:54` | `cmd/ailang-worldd` | `TestOperatorStartFacts` | KILLED | §4.1 |
| X-MCP | `CLI/operator_start_input.go:54` | `cmd/ailang-worldd` | `TestOperatorStartChanged/mcp_path` | KILLED | §4.1 |
| X-MINT-FACTS | `CLI/operator_start_input.go:55` | `cmd/ailang-worldd` | `TestOperatorStartFacts` | KILLED | §4.1 |
| X-BUILD | `CLI/operator_start_input.go:56` | `cmd/ailang-worldd` | `TestOperatorStartFacts` | KILLED | §4.1 |
| X-INTERPRETER | `CLI/operator_start_input.go:57` | `cmd/ailang-worldd` | `TestOperatorStartFacts` | KILLED | §4.1 |
| X-WORKTREE | `CLI/operator_start_input.go:58` | `cmd/ailang-worldd` | `TestOperatorStartFacts` | KILLED | §4.1 |
| NV-M2 | `CLI/operator_start_input.go:68` | `cmd/ailang-worldd` | `TestOperatorStartChanged/base` | KILLED | §4.1 |
| X-LEN | `CLI/operator_start_input.go:72` | `cmd/ailang-worldd` | `TestOperatorStartChanged/extra_fact` | KILLED | §4.1 |
| X-LEN-FIELD | `CLI/operator_start_input.go:73` | `cmd/ailang-worldd` | `TestOperatorStartChanged/missing_fact` | KILLED | §4.1 |
| X-COMPARE | `CLI/operator_start_input.go:76` | `cmd/ailang-worldd` | `TestOperatorStartChanged/preset_ttl` | KILLED | §4.1 |
| X-COMPARE-FIELD | `CLI/operator_start_input.go:77` | `cmd/ailang-worldd` | `TestOperatorStartChanged/base` | KILLED | §4.1 |
| X-RETURN | `CLI/operator_start_input.go:80` | `cmd/ailang-worldd` | `TestOperatorStartChanged/unchanged` | KILLED | §4.1 |

`CLI/` is `cmd/ailang-worldd/`. The ledger file also keeps the 8 superseded first-pass rows. Every row above is a behavioural test. None is an inventory, AST or string-grep test.

**Survivors with no claimed test (carried, Appendix C).** `E-RES-EPISODE-CONFLICT` (`CLI/operator_mint.go:38`, `if current.Episode != r.Episode` → `if false`) and `E-RES-REREAD` (`CLI/operator_mint.go:41`, `r = current` → `_ = current`) each SURVIVED `go test -count=1 -timeout 1m -v -run '^TestOperatorMintRecovery$' ./cmd/ailang-worldd` with the M1 tests in place (ledger). Reason: no single-process test can interleave between the first read and the lock. Abort's in-transaction `current != reviewed` (`host/store/operator_mint.go:306`) refuses a stale review on the only effect path. Both go to the paired-process harness sprint.

**Non-vacuity, per milestone.** M1 has no production diff, so its non-vacuity is the eight M1 rows above: each M1 test fails against the exact production regression it names. M2: `NV-M2` makes `reconstructOperatorStart` return the projection without comparing it, which is the M2 production behaviour reverted to "trust live inputs". The named `TestOperatorStartChanged/base` fails. Deleting `operator_start_input.go` makes the M2 test file fail to compile, so that check was not used as the non-vacuity proof.

**Equivalent / defensive (no mutant claimed).** The `json.Marshal(read.health.Build)` error arm is unreachable: `BuildIdentity` has only string and bool fields. The `in.Resolution()` error arm is unreachable for the same reason (`store.EpisodeStart` and `approval.Grant` are plain data, as S5b §4 records). The `store.CanonicalPath(req.credPath)` error arm can fail only if the store vanishes after `prepareMint` canonicalized it. That is a race no single-process test can arrange, and the error is returned before any effect either way.

### 4.1 Whole-package runs for the controller

The planner ran host-package whole runs inline (`E-TX-*` above, both KILLED by `go test -count=1 -timeout 8m -v ./host/authority`). The 38 CLI whole-package rows are **UNMEASURED**. Reason: one whole `cmd/ailang-worldd` run takes about 240 s (S5b readings), so 38 rows at 3 jobs is about 50 minutes, beyond the planning budget. They were deliberately not started in the background, so they cannot compete with the controller's gates. The controller replays them from the planner's stages (`mutate.py` copies `st-M1`/`st-M2`) with:

````bash
export AILANG_BIN="$HOME/.pinned-ailang/ailang"
S=$HOME/.ailang/state/world-iter433/plan-scratch
python3 "$S/mutate.py" --whole-only --jobs 3 $(python3 -c 'import json;print(" ".join(m["id"] for m in json.load(open("'"$S"'/mutants.json")) if m["pkg"]=="cmd/ailang-worldd"))')
````

## 5. Files

- M1: `cmd/ailang-worldd/{episode_preset_read_test,operator_mint_test}.go`, `host/authority/operator_mint_test.go` (tests only).
- M2: NEW `cmd/ailang-worldd/operator_start_input.go`, NEW `cmd/ailang-worldd/operator_start_input_test.go`.
- Not touched: any `.ail`, `docs/`, `app/`, daemon routes and serving code, `host/store`, `host/approval`, `go.mod`, `tools/launchd/*`, the S3b bank, every production writer in C4.

## 6. Sprint gates

Use §2's extracted script. Scratch runs name the changed file list through `GATES_CHANGED_FILE` (`$S/changed.txt`, the five §5 paths). The controller runs the final gate in the committed checkout with `GATES_SKIP_VERIFY_GO` unset: `bash "$T/gates.sh" "$W" "$HOME/.ailang/state/world-iter433/controller-gates" final`. The `ALL_DONE.<sha>` marker is the last file the gate writes. A skipped leg prints `SKIPPED` in its `LEG` line. Path predicates use `grep -F`. No leg uses `rg`.

````bash gates.sh
#!/bin/bash
# Usage: gates.sh <tree> <logdir> <base|final|scratch> [sha]
set -u
set -o pipefail
export GIT_OPTIONAL_LOCKS=0
T=$1; L=$2; MODE=$3
B=606f6e2c57c79e98975153f055e213df143638d0
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
EXPECTED='cmd/ailang-worldd host/authority '
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
  leg start-input "go test -v -count=1 -timeout 2m ./cmd/ailang-worldd -run '^(TestOperatorStartFacts|TestOperatorStartChanged)$'"
fi
scope_fence() {
  local allowed="$L/allowed.txt" extra
  cat > "$allowed" <<'EOF'
cmd/ailang-worldd/episode_preset_read_test.go
cmd/ailang-worldd/operator_mint_test.go
cmd/ailang-worldd/operator_start_input.go
cmd/ailang-worldd/operator_start_input_test.go
host/authority/operator_mint_test.go
design_docs/planned/w-world-address-bar-verbs-v4s5c-sprint-plan.md
design_docs/planned/sprint_w-world-address-bar-verbs-v4s5c.json
EOF
  printf '%s\n' cmd/ailang-worldd/operator_start_input.go | grep -Fx -f "$allowed"
  printf '%s\n' cmd/ailang-worldd/episode_preset.go | grep -Fxv -f "$allowed"
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
  printf 'Read seam callers:\n%s\n' "$callers"
  printf '%s\n' "$callers" | grep -F 'cmd/ailang-worldd/episode_preset.go:'
  mint=$(grep -R -nF --include='*.go' 'MintOperatorStart(' cmd host | grep -vF '_test.go:' | grep -vF 'func MintOperatorStart(' || :)
  test -z "$mint"
  printf 'Production MintOperatorStart callers: 0\n'
  commit=$(grep -R -nF --include='*.go' '.CommitOperatorMint(' cmd host | grep -vF '_test.go:' || :)
  resolver=$(grep -R -nF --include='*.go' '.AbortOperatorMint(' cmd host | grep -vF '_test.go:' || :)
  test "$(printf '%s\n' "$commit" | wc -l | tr -d ' ')" = 1
  printf '%s\n' "$commit" | grep -F 'host/authority/mint.go:'
  test "$(printf '%s\n' "$resolver" | wc -l | tr -d ' ')" = 1
  printf '%s\n' "$resolver" | grep -F 'cmd/ailang-worldd/operator_mint.go:'
  rebuild=$(grep -R -nF --include='*.go' 'reconstructOperatorStart(' cmd host | grep -vF '_test.go:' | grep -vF 'func reconstructOperatorStart(' || :)
  facts=$(grep -R -nF --include='*.go' 'operatorStartFacts(' cmd host | grep -vF '_test.go:' | grep -vF 'func operatorStartFacts(' || :)
  test -z "$rebuild"
  printf 'Production reconstructOperatorStart callers: 0\n'
  if [ "$MODE" = base ]; then
    test "$(printf '%s\n' "$callers" | wc -l | tr -d ' ')" = 1
    test -z "$facts"
  else
    test "$(printf '%s\n' "$callers" | wc -l | tr -d ' ')" = 2
    printf '%s\n' "$callers" | grep -F 'cmd/ailang-worldd/operator_start_input.go:'
    test "$(printf '%s\n' "$facts" | wc -l | tr -d ' ')" = 1
    printf '%s\n' "$facts" | grep -F 'cmd/ailang-worldd/operator_start_input.go:'
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

**Base run (pristine `git archive 606f6e2c`, mode `base`, outside any sandbox).** All 27 legs rc=0, `verify_go SKIPPED`, `GATES_DONE 764s FAILS: none UNINFORMATIVE: none`. That includes `pkg-cmd-ailang-worldd` (286 s), the race legs and every inventory, so no leg is red at base. Log: `$S/gate-base/legs.txt`.

**Final scratch run (`st-M2`, mode `scratch`).**

````text
LEG observation-before rc=0 1s
LEG paths rc=0 0s
TOUCHED_PACKAGES: cmd/ailang-worldd host/authority 
LEG touched rc=0 0s
LEG status rc=0 0s
LEG gofmt rc=0 0s
LEG build rc=0 2s
LEG vet rc=0 1s
LEG fence rc=0 17s
LEG verify_ail rc=0 36s
LEG pkg-cmd-ailang-worldd rc=0 288s
LEG race-cmd-ailang-worldd rc=0 395s
LEG pkg-host-authority rc=0 2s
LEG race-host-authority rc=0 16s
LEG inv-host-approval rc=0 1s
LEG inv-host-approval-approvaltest rc=0 1s
LEG inv-host-authority rc=0 1s
LEG inv-host-boundary rc=0 1s
LEG inv-host-daemon rc=0 1s
LEG inv-host-verbs rc=0 1s
LEG inv-host-broker rc=0 1s
LEG inv-cmd-ailang-worldd rc=0 1s
LEG inv-quickstart-section9 rc=0 0s
LEG inv-host-store rc=0 1s
LEG routes rc=0 0s
LEG inv-cmd-world-publish rc=0 1s
LEG mint-quickstart rc=0 2s
LEG start-input rc=0 6s
LEG scope rc=0 0s
LEG callers rc=0 1s
LEG verify_go SKIPPED 0s (GATES_SKIP_VERIFY_GO=1)
LEG observation-after rc=0 1s
GATES_DONE 778s FAILS: none UNINFORMATIVE: none
````

`callers` printed the positive control `cmd/ailang-worldd/episode_finish_surface.go:239` (`authority.RecheckVerb`), 7 `Flock LOCK_NB` lines and 5 provision callers. It printed two read-seam callers (`episode_preset.go:328`, `operator_start_input.go:19`) and 0 production `MintOperatorStart` and `reconstructOperatorStart` callers. It printed one `CommitOperatorMint` caller (`host/authority/mint.go:219`), one `AbortOperatorMint` caller (`cmd/ailang-worldd/operator_mint.go:64`) and one `operatorStartFacts` caller, inside `operator_start_input.go:68`. In the `callers` output, `episode_preset.go:328` is printed twice: once for the caller list and once for the `grep -F` predicate echo. `routes`: `DAEMON_ROUTES: 50 (base 50); SURFACE_POSTS: 3 (base 3)`, with `daemon.go`, `surface_page.go` and `app.js` byte-identical to base. `scope`: `UNTOUCHED_VIOLATIONS: none`. `observation-before` and `observation-after` (real `~/Library/LaunchAgents`, `~/.ailang-world`, the real credentials store and its sidecars, and the read-only `launchctl list` count) are byte-identical. Darwin `race-cmd-ailang-worldd` took 395 s under concurrent load (the base gate ran in parallel), against the 8 m limit. Logs: `$S/gate-final/`.

The controller repeats the gate in the committed checkout with `verify_go`, which covers CLI plain and race, QUICKSTART §9 and full plain plus Linux race in CI: `go test -count=1 -timeout 8m ./cmd/ailang-worldd`, `go test -race -count=1 -timeout 8m ./cmd/ailang-worldd`, `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh`, `VERIFY_GO_LEGS=race ./scripts/verify_go.sh`.

## 7. Executor directive notes

- Commit-time commands, one per milestone, run from the controller-owned snapshot in the same pinned fixture environment as scratch (`accept.sh` exports `AILANG_BIN`, `WORLD_REQUIRE_SSH_KEYGEN=1`, `WORLD_EXEC_NODE`, `WORLD_TOOL_AILANG_BIN` and private `AILANG_STATE_DIR`/`AILANG_CACHE_DIR`, and unsets the browser variables). Do not set `TMPDIR` or `GOTMPDIR`.

````bash
T=$HOME/.ailang/state/world-iter433/executor
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter433-row243-v4s5c
bash "$HOME/.ailang/state/world-iter433/commit_ms.sh" M1 "$msgfile" "AILANG_BIN=$HOME/.pinned-ailang/ailang bash '$T/accept.sh' '$W' M1 fast"
bash "$HOME/.ailang/state/world-iter433/commit_ms.sh" M2 "$msgfile" "AILANG_BIN=$HOME/.pinned-ailang/ailang bash '$T/accept.sh' '$W' M2 fast"
````

- Expected `NUMSTAT` printed by `commit_ms.sh` (`git diff -M --numstat 606f6e2c HEAD -- cmd host docs app`, summed): **91** after M1 and **401** after M2. No renames are expected, so the `=>` listing is empty.
- The snapshot is cumulative from the base. It accepts missing `docs/`, `app/` or `cmd/` directories and excludes the plan records through its pathspec. Never reuse a snapshot directory. Expected STATUS/DELETED/MANIFEST counts are 3/0/3 (M1) and 5/0/5 (M2).

````bash snapshot.sh
#!/bin/bash
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=$1
M=$2
B=606f6e2c57c79e98975153f055e213df143638d0
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter433/snap}/$M
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

**Planner rehearsal (run).** `SNAP_STATUS_LIST="$S/status-M1" SNAP_PARENT="$S/snap-rehearsal" bash "$S/snapshot.sh" "$S/st-M1" M1` and the same for M2 gave 3/0/3 and 5/0/5, with every manifest line `OK`. The planner's own round trip (`$S/rt1`, not the controller's `ctl-roundtrip` directory, which `ctl_roundtrip.sh` refuses to reuse) extracted `M1` and `M2` from this plan file with the same awk as `ctl_roundtrip.sh`. It applied them in order with `/usr/bin/patch -p1 -E -s` to a fresh `git archive 606f6e2c` export and ran `diff -r` on `cmd`, `host` and `docs` against `st-M1` after M1 and `st-M2` after M2. All six comparisons were empty (rc=0), with no `.orig`/`.rej` files. The controller runs `bash "$HOME/.ailang/state/world-iter433/ctl_roundtrip.sh" M1 M2`.

## 8. Risks

- **Fact-list contract.** The future catalogue (S7) must sign exactly `operatorStartFacts`'s list. Any producer drift fails closed as `BindingError{resolution}` or `{<name>}`, a stale review with no effect. That is intended, and `TestOperatorStartFacts` pins the list.
- **Reconstruction is not provenance.** Equal facts prove the review matches the live inputs at this instant. They do not prove absence of a worktree, branch or files at reservation. S5d's reservation must still demand initial absence under the episode lock and re-run the reconstruction inside that lock.
- **Linux/Darwin.** M1's FIFO arm relies on `syscall.Mkfifo` and on `open(O_RDONLY|O_NONBLOCK)` succeeding on a FIFO with no writer, which holds on both kernels. The assertion is the type class only. M1 relies on no `flock` behaviour. `VACUUM INTO` is SQLite ≥3.27 and is available in the pinned driver (test passes).
- **CI time.** About 6 s is added to the CLI package plain run and somewhat more under `-race`, against a CLI race leg of 347 s with an 8 m limit (S5b readings). Nothing is added to `host/daemon` or `host/verifygate`.

## 9. Acceptance criteria

1. M1 and M2 apply in order from the base and reproduce the compiled stages under the controller's round trip. `commit_ms.sh` prints NUMSTAT 91 and 401. Snapshots match §7.
2. M1: `TestOperatorMintScope` (`host/authority`), `TestOperatorMintRecovery/{store_identity,token_name}`, `TestOperatorMintTokenReadback/{fifo,directory}` and the two parity assertions pass on the base production code, and each fails against its named survivor (§4).
3. M2: `TestOperatorStartFacts` and every `TestOperatorStartChanged` arm pass. Each changed input refuses with its exact field or class, and every arm shows unchanged provision state and no lock, worktree or agent run.
4. `callers` leg: zero production callers of `MintOperatorStart`, `reconstructOperatorStart`, claim/reserve/observe; exactly one `CommitOperatorMint` caller and one `AbortOperatorMint` caller. There are two read-seam callers (terminal, inert reconstruction), and the positive control is present. Routes, served bytes, writer sites and the lock count are unchanged.
5. Every Appendix C gate travels with S5d and later. No activation happens in this sprint.

## Appendix C — after this sprint (measured)

Sizes come from re-applying the iteration-432 probes to `st-M2` (`$S/st-s5d-footprint`). Both `go vet ./cmd/ailang-worldd ./host/store` and the probe tests passed (`TestS5cEpisodeFDPrototype` ok, `TestOperatorStepRestart` PASS). `git diff --no-index -M --numstat "$S/st-M2/<d>" "$S/st-s5d-footprint/<d>"` gave `s5c_worktree_probe.go` 59, `s5c_files_probe.go` 53, `s5c_fd_probe_test.go` 72, `host/store/operator_steps.go` 72 and `operator_cas_test.go` +60. The measured floor is **316** lines, beyond the 401 landed here. These are sizing probes, not implementations. Each sprint below must re-measure its compiled total against the 1,400 cap and split before activation if it exceeds it.

- **S5d — worktree reservation and attended worktree disposition (59-line readback probe; full size unmeasured).** Run `reconstructOperatorStart` under the canonical episode lock, then reserve the `worktree` step (S3b bank: reserve/observe step methods, 132 measured lines with `TestOperatorStepRestart`) before `worktree.Add` (`session_new.go:510,512` precedent). Read back registration (NUL-delimited `worktree list --porcelain -z`), physical path, common dir, branch and starting SHA against the reservation. Known absence allows a bounded retry or a reviewed stop. A partial or unattributable result is `uncertain` and ships **in the same patch** with an attended `service runner reconcile-worktree ACTION-ID`. That command has tested exact remedy text that a test executes verbatim, effect-specific readback and an independent current-owner/step CAS. It never removes an unattributed worktree. Inject failures before reservation, after reservation before add, and after add before readback. Each injection asserts store rows, worktrees, branch refs and files by read-back. `readOperatorMint` already tolerates an observed `worktree` step (`host/store/operator_mint.go:266`).
- **S5e — task/token/MCP publication (53-line fd readback probe; full size unmeasured).** Private exclusive/no-replace staging to the canonical paths (`writeTokenFile` `session_new.go:674`, `writeMCPConfig` `episode.go:552`, task copy `episode_preset.go:342`). Pin path, owner, mode, digest, link count, type, exact bytes and writer order independently. An injected failure before and after each publish ships with an attended file disposition bound to owner and step version. Revalidate the current credential and token before any launch.
- **S5f/S6 — paired-process consumer and drain (72-line fd probe + 132-line S3b bank; full size unmeasured), binding.** No live claim/reserve/observe/mint until the process holds the same canonical episode fd before the CAS and across effect boundaries. The required pins are `TestOperatorFirstManagedEpisodeProcess`, `TestOperatorOldBuildTerminalDrain`, `TestOperatorCurrentOwnerProcessRecovery` and `TestOperatorServiceUpgrade/active_child`. Before giving `MintOperatorStart` a production caller, this sprint must also kill **`E-RES-EPISODE-CONFLICT`** and **`E-RES-REREAD`** (`CLI/operator_mint.go:38,41`; both SURVIVED with S5c's tests, §4). Its process harness can interleave a state change between the resolver's first read and its lock. Also carry direct terminal `session revoke`, concurrent credential changes and exact session revalidation before launch or recovery.
- **J-OWN-UID, `Z-ABSENT-ANY`, `Z-WAIT-READBUDGET`.** Carried unchanged from S5/S5b Appendix C. J-OWN-UID needs an alternate-owner fixture (ordinary non-root), which belongs with the process gate. S5c touches neither the runner lock/path check nor the wait/read budget code.
- **Closed here.** `J-PRINT-AGENT` and `J-DOUBLE-OVERRIDE-APPEND` (8 lines, M1). `E-TX-STORE-ID`, `E-TX-KIND`, `E-RES-STORE-ID`, `E-RES-TOKENPATH-BASE` and `E-RES-REGULAR` (M1). Live-input reconstruction (M2, replacing the 54-line projection probe).
- **Launch, browser and finish.** Design budgets S6 890 raw, S8 860 raw, S12 700 raw. Retain E3, E5 and #604 sequencing, E4, and the D-WORLD-120 finish resolver before S9.

## Appendix P — compiled patches

### Patch M1

````diff M1
diff -ruN base/cmd/ailang-worldd/episode_preset_read_test.go st-M1/cmd/ailang-worldd/episode_preset_read_test.go
--- base/cmd/ailang-worldd/episode_preset_read_test.go	2026-11-25 06:56:45
+++ st-M1/cmd/ailang-worldd/episode_preset_read_test.go	2026-11-25 08:30:20
@@ -72,6 +72,9 @@
 		t.Fatal("the attended path did not use the resolved review and bytes")
 	}
 	got := readPresetStart(t, f, "ep")
+	if got == nil || !reflect.DeepEqual(got.Overrides, []string{"agent", "branch", "ttl"}) {
+		t.Fatal("persisted overrides differ from explicit flags")
+	}
 	want := *req.start
 	dir := filepath.Join(f.r.f.ws, "ep")
 	head, err := worktree.Output(presetTestContext(t), f.r.f.git, []string{"-C", dir, "rev-parse", "HEAD"})
@@ -112,6 +115,11 @@
 	}
 	if _, err := os.Lstat(filepath.Join(bin, "ran")); !os.IsNotExist(err) {
 		t.Fatal("a read started the executable")
+	}
+	var out bytes.Buffer
+	if code := runEpisodeStart(f.args("ep", "--task", f.task), &out, &out, presetEnv(&promptTerm{in: strings.NewReader("y\n")})); code != exitOK ||
+		!strings.Contains(out.String(), " codex exec --ignore-user-config ") {
+		t.Fatalf("Codex start did not print its confined launch: %d %s", code, &out)
 	}
 }
 
diff -ruN base/cmd/ailang-worldd/operator_mint_test.go st-M1/cmd/ailang-worldd/operator_mint_test.go
--- base/cmd/ailang-worldd/operator_mint_test.go	2026-11-25 06:56:45
+++ st-M1/cmd/ailang-worldd/operator_mint_test.go	2026-11-25 08:29:52
@@ -10,6 +10,7 @@
 	"path/filepath"
 	"reflect"
 	"strings"
+	"syscall"
 	"testing"
 
 	"github.com/sunholo-data/ailang-world/host/approval"
@@ -108,6 +109,26 @@
 	}
 }
 
+// rewriteAction replaces the stored action and its digest together, so only
+// the resolver's own store and token-path identity checks can refuse it.
+func (r recoveryFixture) rewriteAction(t *testing.T, edit func(*approval.OperatorAction)) {
+	t.Helper()
+	var raw []byte
+	if err := r.db.QueryRowContext(boundedTestContext(t), "SELECT action FROM operator_actions").Scan(&raw); err != nil {
+		t.Fatal(err)
+	}
+	a, err := approval.DecodeOperatorAction(raw)
+	if err != nil {
+		t.Fatal(err)
+	}
+	edit(&a)
+	if raw, err = approval.CanonicalOperatorAction(a); err != nil {
+		t.Fatal(err)
+	}
+	digest := sha256.Sum256(raw)
+	r.exec(t, "UPDATE operator_actions SET action=?, digest=?", raw, digest[:])
+}
+
 func (r recoveryFixture) run(t *testing.T, command string) (int, string) {
 	t.Helper()
 	args := strings.Fields(command)
@@ -137,6 +158,8 @@
 		{"receipt_changed", "operator action conflicts"},
 		{"missing_credential", "mint start or credential unavailable"},
 		{"abort_rollback", "injected abort failure"},
+		{"store_identity", "mint token path or store identity differs"},
+		{"token_name", "mint token path or store identity differs"},
 	} {
 		t.Run(tc.name, func(t *testing.T) {
 			r := newRecoveryFixture(t)
@@ -188,6 +211,18 @@
 				r.exec(t, "DELETE FROM session_credentials WHERE credential_id=?", r.receipt.CredentialID)
 			case "abort_rollback":
 				r.exec(t, "CREATE TRIGGER failed_abort AFTER UPDATE ON operator_actions BEGIN SELECT RAISE(ABORT,'injected abort failure'); END")
+			case "store_identity", "token_name":
+				r.rewriteAction(t, func(a *approval.OperatorAction) {
+					if tc.name == "store_identity" {
+						a.Store = filepath.Join(filepath.Dir(a.Store), "other.db")
+						return
+					}
+					for i := range a.Resolution {
+						if a.Resolution[i].Name == "mint_token_path" {
+							a.Resolution[i].Value = filepath.Join(filepath.Dir(r.path), "other.session")
+						}
+					}
+				})
 			}
 			before := serviceFileHash(r.config.Credentials)
 			files := provisionFiles(t, filepath.Dir(r.path))
@@ -269,7 +304,7 @@
 }
 
 func TestOperatorMintTokenReadback(t *testing.T) {
-	for _, name := range []string{"uid", "mode", "link", "symlink", "symlink_valid", "directory", "bytes", "length", "short"} {
+	for _, name := range []string{"uid", "mode", "link", "symlink", "symlink_valid", "directory", "fifo", "bytes", "length", "short"} {
 		t.Run(name, func(t *testing.T) {
 			dir := t.TempDir()
 			p := filepath.Join(dir, "token")
@@ -287,6 +322,8 @@
 				err = os.Symlink(target, p)
 			case "directory":
 				err = os.Mkdir(p, 0700)
+			case "fifo":
+				err = syscall.Mkfifo(p, 0600)
 			default:
 				data := token
 				switch name {
@@ -314,6 +351,10 @@
 			present, err := readMintToken(p, uid, fmt.Sprintf("%x", sha256.Sum256([]byte(token))))
 			if present || err == nil || !reflect.DeepEqual(before, provisionFiles(t, dir)) {
 				t.Fatalf("unsafe file evidence changed or accepted: present=%v err=%v", present, err)
+			}
+			// The file type is refused as a type, before any mode or owner check.
+			if (name == "directory" || name == "fifo") && err.Error() != "mint token is not a regular file" {
+				t.Fatalf("%s refused for the wrong reason: %v", name, err)
 			}
 		})
 	}
diff -ruN base/host/authority/operator_mint_test.go st-M1/host/authority/operator_mint_test.go
--- base/host/authority/operator_mint_test.go	2026-11-25 06:56:45
+++ st-M1/host/authority/operator_mint_test.go	2026-11-25 08:29:12
@@ -381,3 +381,43 @@
 	}
 	r.unchanged(t, got, err, store.ErrOperatorConflict)
 }
+
+// TestOperatorMintScope pins the transaction's own scope guards. A signed
+// action mints only in the store it names, and only as a start.
+func TestOperatorMintScope(t *testing.T) {
+	t.Run("store_copy", func(t *testing.T) {
+		r := operatorMintRig(t, false, nil)
+		seal := r.seal(t)
+		other := filepath.Join(t.TempDir(), "copy.db")
+		jgExec(t, r.path, fmt.Sprintf("VACUUM INTO '%s'", other))
+		copied := operatorStoreOpen(t, other, store.OpenCredentials)
+		before := r.snapshot(t)
+		err := copied.CommitOperatorMint(r.ctx, r.in, seal, mintAt(1014))
+		sessions, listErr := copied.ListSessions(r.ctx, "", 20)
+		start, startErr := copied.EpisodeStart(r.ctx, r.action.Episode)
+		if listErr != nil || startErr != nil || len(sessions) != 0 || start != nil {
+			t.Fatalf("a copied store minted: %v %v %d %v", listErr, startErr, len(sessions), start)
+		}
+		r.unchanged(t, before, err, store.ErrOperatorConflict)
+	})
+	t.Run("finish_kind", func(t *testing.T) {
+		// No runner row exists, so only the kind refusal can answer conflict.
+		s := newOperatorStoreRigAction(t, false, func(a *approval.OperatorAction) {
+			a.Kind, a.Task, a.Options.Headless = "finish", "", false
+			a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
+		})
+		ctx := boundedTestContext(t)
+		s.confirm(t)
+		jgExec(t, s.path, `UPDATE operator_actions SET state = 'running', owner = 'fresh-runner'`,
+			fmt.Sprintf(`INSERT INTO operator_action_steps (kind, identity, action_id, phase, version, owner, created_at)
+			VALUES ('mint', '%s', '%s', 'reserved', 7, 'fresh-runner', 1012)`, s.cas.ID, s.cas.ID))
+		row, err := s.OperatorAction(ctx, s.cas.ID, 1012)
+		if err != nil {
+			t.Fatal(err)
+		}
+		r := mintRig{s, store.OperatorMintInput{Attempt: store.OperatorMintAttempt{CAS: s.cas, StepVersion: 7,
+			Runner: store.OperatorOwner{Nonce: "fresh-runner"}, Catalogue: 1}}, row, operatorStoreOpen(t, s.path, store.OpenCredentialsReadOnly), ctx}
+		before := r.snapshot(t)
+		r.unchanged(t, before, r.CommitOperatorMint(ctx, r.in, nil, mintAt(1014)), store.ErrOperatorConflict)
+	})
+}
````

### Patch M2

````diff M2
diff -ruN st-M1/cmd/ailang-worldd/operator_start_input.go st-M2/cmd/ailang-worldd/operator_start_input.go
--- st-M1/cmd/ailang-worldd/operator_start_input.go	1970-01-01 01:00:00
+++ st-M2/cmd/ailang-worldd/operator_start_input.go	2026-11-25 08:31:29
@@ -0,0 +1,81 @@
+package main
+
+import (
+	"context"
+	"encoding/json"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// operatorStartFacts projects one live read of a headless start into the
+// effective mint input and the complete, name-ordered resolution a reviewed
+// action signs. It derives the task only from the signed task bytes.
+func operatorStartFacts(ctx context.Context, signed approval.OperatorAction, attempt store.OperatorMintAttempt,
+	presets, profiles, workspace, hosts string, clock func() int64) (store.OperatorMintInput, []approval.OperatorResolved, error) {
+	if signed.Kind != "start" || !signed.Options.Headless {
+		return store.OperatorMintInput{}, nil, store.ErrOperatorConflict
+	}
+	read, err := resolvePresetStart(ctx, signed.Options.Preset, presets, profiles, presetStartInput{
+		episodes: []string{signed.Episode}, workspaceRoot: workspace, hostsDir: hosts, launching: true, headless: true,
+		readTask: func() (string, []byte, error) { return "inline", []byte(signed.Task), nil },
+	}, clock)
+	if err != nil {
+		return store.OperatorMintInput{}, nil, err
+	}
+	req := read.mint
+	if req.prof == nil {
+		return store.OperatorMintInput{}, nil, store.ErrOperatorInput
+	}
+	served, err := store.CanonicalPath(req.credPath)
+	if err != nil {
+		return store.OperatorMintInput{}, nil, err
+	}
+	if served != signed.Store {
+		return store.OperatorMintInput{}, nil, &approval.BindingError{Field: "store"}
+	}
+	start, task := *req.start, *req.start.Task
+	task.Path = read.taskCopy
+	start.Task = &task
+	in := store.OperatorMintInput{Attempt: attempt, Start: start, TokenPath: req.outPath,
+		Profile: approval.ProfileRef{Project: req.profile, Digest: req.prof.Digest}}
+	for _, g := range req.grants {
+		in.Grants = append(in.Grants, approval.Grant{Effect: g.Effect, Scope: g.Scope, Budget: g.Budget})
+	}
+	mint, err := in.Resolution()
+	if err != nil {
+		return store.OperatorMintInput{}, nil, err
+	}
+	build, err := json.Marshal(read.health.Build)
+	if err != nil {
+		return store.OperatorMintInput{}, nil, err
+	}
+	facts := []approval.OperatorResolved{{Name: "launch_binary", Value: read.launchBinary},
+		{Name: "launch_wall", Value: read.launchWall.String()}, {Name: "mcp_path", Value: read.mcp}}
+	facts = append(facts, mint...)
+	facts = append(facts, approval.OperatorResolved{Name: "served_build", Value: string(build)},
+		approval.OperatorResolved{Name: "served_interpreter", Value: read.health.ToolRef},
+		approval.OperatorResolved{Name: "worktree_path", Value: req.dir})
+	return in, facts, nil
+}
+
+// reconstructOperatorStart requires every signed resolution fact to equal its
+// live value before any effect. A change, even one lowering authority, needs a
+// fresh review. It locks, writes, mints and launches nothing, and it has no
+// production caller until the paired-process consumer lands.
+func reconstructOperatorStart(ctx context.Context, signed approval.OperatorAction, attempt store.OperatorMintAttempt,
+	presets, profiles, workspace, hosts string, clock func() int64) (store.OperatorMintInput, error) {
+	in, facts, err := operatorStartFacts(ctx, signed, attempt, presets, profiles, workspace, hosts, clock)
+	if err != nil {
+		return store.OperatorMintInput{}, err
+	}
+	if len(facts) != len(signed.Resolution) {
+		return store.OperatorMintInput{}, &approval.BindingError{Field: "resolution"}
+	}
+	for i, fact := range facts {
+		if signed.Resolution[i] != fact {
+			return store.OperatorMintInput{}, &approval.BindingError{Field: fact.Name}
+		}
+	}
+	return in, nil
+}
diff -ruN st-M1/cmd/ailang-worldd/operator_start_input_test.go st-M2/cmd/ailang-worldd/operator_start_input_test.go
--- st-M1/cmd/ailang-worldd/operator_start_input_test.go	1970-01-01 01:00:00
+++ st-M2/cmd/ailang-worldd/operator_start_input_test.go	2026-11-25 08:36:02
@@ -0,0 +1,229 @@
+package main
+
+import (
+	"encoding/json"
+	"errors"
+	"os"
+	"path/filepath"
+	"reflect"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+type operatorStartRig struct {
+	f      presetFixture
+	creds  string
+	hosts  string
+	bin    string
+	action approval.OperatorAction
+	at     store.OperatorMintAttempt
+	in     store.OperatorMintInput
+}
+
+const operatorStartTask = "# Fix fixture\nExact signed task bytes.\n"
+
+// newOperatorStartRig reviews one headless start: the signed resolution is the
+// projection of a first live read, as the catalogue would sign it.
+func newOperatorStartRig(t *testing.T) operatorStartRig {
+	t.Helper()
+	f := newPresetFixture(t)
+	r := operatorStartRig{f: f, bin: startReadFakeAgent(t, f), hosts: t.TempDir(), at: store.OperatorMintAttempt{Catalogue: 3}}
+	r.creds = matchingHealth(t, f.r.f, f.r.profiles, f.profile).CredentialsDB
+	canonical, err := store.CanonicalPath(r.creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	r.action = approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: fixtureActionID, Kind: "start", Episode: "ep",
+		Store: canonical, Options: approval.OperatorOptions{Preset: "p", Headless: true}, Task: operatorStartTask, Review: "fixture"}
+	r.action.Argv = approval.OperatorArgv(r.action.Kind, r.action.Episode, r.action.Options)
+	in, facts, err := operatorStartFacts(presetTestContext(t), r.action, r.at, f.dir, f.r.profiles, f.r.f.ws, r.hosts, presetClock)
+	if err != nil {
+		t.Fatal(err)
+	}
+	r.action.Resolution, r.in = facts, in
+	if err := r.action.Validate(); err != nil {
+		t.Fatalf("projection is not a signable action: %v", err)
+	}
+	return r
+}
+
+func (r operatorStartRig) reconstruct(t *testing.T) (store.OperatorMintInput, error) {
+	t.Helper()
+	return reconstructOperatorStart(presetTestContext(t), r.action, r.at, r.f.dir, r.f.r.profiles, r.f.r.f.ws, r.hosts, presetClock)
+}
+
+func (r operatorStartRig) editPreset(t *testing.T, edit func(*launchPreset)) {
+	t.Helper()
+	var p launchPreset
+	if err := json.Unmarshal(readStartFixtureFile(t, r.f.file), &p); err != nil {
+		t.Fatal(err)
+	}
+	edit(&p)
+	data, err := json.MarshalIndent(p, "", "  ")
+	if err != nil {
+		t.Fatal(err)
+	}
+	seedStartFixtureFile(t, r.f.file, data)
+}
+
+// TestOperatorStartFacts pins every reviewed fact against fixture values
+// derived independently of the projection.
+func TestOperatorStartFacts(t *testing.T) {
+	r := newOperatorStartRig(t)
+	f := r.f
+	names := []string{}
+	value := map[string]string{}
+	for _, fact := range r.action.Resolution {
+		names = append(names, fact.Name)
+		value[fact.Name] = fact.Value
+	}
+	want := []string{"launch_binary", "launch_wall", "mcp_path", "mint_catalogue", "mint_grants", "mint_profile", "mint_project",
+		"mint_start", "mint_token_path", "served_build", "served_interpreter", "worktree_path"}
+	if !reflect.DeepEqual(names, want) {
+		t.Fatalf("reviewed fact names %v, want %v", names, want)
+	}
+	binary, err := filepath.EvalSymlinks(filepath.Join(r.bin, "codex"))
+	if err != nil {
+		t.Fatal(err)
+	}
+	build, _ := json.Marshal(testBuild)
+	h := matchingHealth(t, f.r.f, f.r.profiles, f.profile)
+	for name, v := range map[string]string{"launch_binary": binary, "launch_wall": "1h0m0s", "mcp_path": "",
+		"mint_catalogue": "3", "mint_profile": f.profile.Digest, "mint_project": "proj",
+		"mint_token_path": filepath.Join(f.r.sessions, "ep.session"), "served_build": string(build),
+		"served_interpreter": h.ToolRef, "worktree_path": filepath.Join(f.r.f.ws, "ep")} {
+		if value[name] != v {
+			t.Fatalf("fact %s = %q, want %q", name, value[name], v)
+		}
+	}
+	in := r.in
+	if in.Start.Episode != "ep" || in.Start.Base.SHA != f.base || in.Start.Preset.SHA256 != serviceHash(f.data) ||
+		in.Start.Branch != "world/ep" || in.Start.BranchExisted || in.Start.Agent != "codex" || in.Start.TTLS != 43200 ||
+		in.Start.Task == nil || in.Start.Task.Path != filepath.Join(f.r.sessions, "ep.task.md") ||
+		in.Start.Task.SHA256 != serviceHash([]byte(operatorStartTask)) || in.TokenPath != filepath.Join(f.r.sessions, "ep.session") ||
+		in.Profile != (approval.ProfileRef{Project: "proj", Digest: f.profile.Digest}) || len(in.Grants) == 0 || in.Attempt != r.at {
+		t.Fatalf("effective input differs from the fixture: %+v", in)
+	}
+	caps, err := sessionNewGrants("se-tools", nil, 50)
+	if err != nil {
+		t.Fatal(err)
+	}
+	grants := []approval.Grant{}
+	for _, c := range caps {
+		grants = append(grants, approval.Grant{Effect: c.Effect, Scope: c.Scope, Budget: c.Budget})
+	}
+	if !reflect.DeepEqual(in.Grants, grants) {
+		t.Fatalf("effective grants %+v, want the preset's %+v", in.Grants, grants)
+	}
+	mint, err := in.Resolution()
+	if err != nil || !reflect.DeepEqual(mint, r.action.Resolution[3:9]) {
+		t.Fatalf("mint facts differ from the effective input: %v", err)
+	}
+}
+
+// TestOperatorStartChanged changes one live input after review at a time.
+// Each change is a stale review with no lock, file, store or launch effect.
+func TestOperatorStartChanged(t *testing.T) {
+	for _, tc := range []struct {
+		name, field string
+		class       error
+		change      func(*testing.T, *operatorStartRig)
+	}{
+		{name: "unchanged"},
+		{name: "preset_ttl", field: "mint_start", change: func(t *testing.T, r *operatorStartRig) {
+			r.editPreset(t, func(p *launchPreset) { p.TTLS-- })
+		}},
+		{name: "preset_grants", field: "mint_grants", change: func(t *testing.T, r *operatorStartRig) {
+			r.editPreset(t, func(p *launchPreset) { p.Grants[0].Budget-- })
+		}},
+		{name: "base", field: "mint_start", change: func(t *testing.T, r *operatorStartRig) {
+			presetGit(t, r.f.r.f, "branch", "-f", "dev", "feature")
+		}},
+		{name: "profile", field: "mint_profile", change: func(t *testing.T, r *operatorStartRig) {
+			p := writeTestProfile(t, r.f.r.profiles, "proj", `, "env": {"FIXTURE": "changed"}`)
+			h := matchingHealth(t, r.f.r.f, r.f.r.profiles, p)
+			h.DBPath = r.f.r.f.db
+			stubDaemon(t, h, nil)
+		}},
+		{name: "task", field: "mint_start", change: func(t *testing.T, r *operatorStartRig) { r.action.Task += "changed\n" }},
+		{name: "catalogue", field: "mint_catalogue", change: func(t *testing.T, r *operatorStartRig) { r.at.Catalogue++ }},
+		{name: "build", field: "served_build", change: func(t *testing.T, r *operatorStartRig) {
+			other := daemon.BuildIdentity{VCSRevision: "fedcba9876543210fedcba9876543210fedcba98"}
+			h := matchingHealth(t, r.f.r.f, r.f.r.profiles, r.f.profile)
+			h.DBPath, h.Build = r.f.r.f.db, other
+			stubDaemon(t, h, nil)
+			episodeSelfBuild = func() daemon.BuildIdentity { return other }
+		}},
+		{name: "interpreter", field: "served_interpreter", change: func(t *testing.T, r *operatorStartRig) {
+			r.action.Resolution[10].Value = "sha256:" + serviceHash([]byte("other"))
+		}},
+		{name: "launch_binary", field: "launch_binary", change: func(t *testing.T, r *operatorStartRig) {
+			other := startReadFakeAgent(t, r.f)
+			if other == r.bin {
+				t.Fatal("fixture reused the agent directory")
+			}
+		}},
+		{name: "mcp_path", field: "mcp_path", change: func(t *testing.T, r *operatorStartRig) {
+			r.action.Resolution[2].Value = filepath.Join(r.f.r.sessions, "ep.mcp.json")
+		}},
+		{name: "worktree_path", field: "worktree_path", change: func(t *testing.T, r *operatorStartRig) {
+			r.action.Resolution[11].Value += "-other"
+		}},
+		{name: "extra_fact", field: "resolution", change: func(t *testing.T, r *operatorStartRig) {
+			r.action.Resolution = append(r.action.Resolution, approval.OperatorResolved{Name: "zz_extra", Value: "x"})
+		}},
+		{name: "missing_fact", field: "resolution", change: func(t *testing.T, r *operatorStartRig) {
+			r.action.Resolution = r.action.Resolution[:len(r.action.Resolution)-1]
+		}},
+		{name: "store", field: "store", change: func(t *testing.T, r *operatorStartRig) {
+			r.action.Store = filepath.Join(filepath.Dir(r.action.Store), "other.db")
+		}},
+		{name: "finish", class: store.ErrOperatorConflict, change: func(t *testing.T, r *operatorStartRig) {
+			r.action.Kind = "finish"
+		}},
+		{name: "interactive", class: store.ErrOperatorConflict, change: func(t *testing.T, r *operatorStartRig) {
+			r.action.Options.Headless = false
+		}},
+		{name: "no_profile", class: store.ErrOperatorInput, change: func(t *testing.T, r *operatorStartRig) {
+			r.editPreset(t, func(p *launchPreset) {
+				p.Profile, p.Checks = "", nil
+				p.Grants = []presetGrant{{Effect: "Workspace.Read", Scope: "worktree", Budget: 1}}
+			})
+		}},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			r := newOperatorStartRig(t)
+			if tc.change != nil {
+				tc.change(t, &r)
+			}
+			before := observeProvision(t, r.f.r, r.creds)
+			got, err := r.reconstruct(t)
+			switch {
+			case tc.field != "":
+				var mismatch *approval.BindingError
+				if !errors.As(err, &mismatch) || mismatch.Field != tc.field {
+					t.Fatalf("expected stale review %s, got %v", tc.field, err)
+				}
+			case tc.class != nil:
+				if !errors.Is(err, tc.class) {
+					t.Fatalf("expected %v, got %v", tc.class, err)
+				}
+			default:
+				if err != nil || !reflect.DeepEqual(got, r.in) {
+					t.Fatalf("unchanged review did not reconstruct its input: %v", err)
+				}
+			}
+			if !reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
+				t.Fatal("reconstruction changed durable state")
+			}
+			for _, p := range []string{operatorEpisodeLockPath(r.creds, "ep"), filepath.Join(r.bin, "ran"), filepath.Join(r.f.r.f.ws, "ep")} {
+				if _, err := os.Lstat(p); !os.IsNotExist(err) {
+					t.Fatalf("reconstruction created %s: %v", p, err)
+				}
+			}
+		})
+	}
+}
````
