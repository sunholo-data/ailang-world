# Sprint plan — w-world-address-bar-verbs-v4, sprint S4 (iteration 429, row 243)

> **THE SPRINT — row 243, phase V4, sprint S4.** Reviewed design r1 §§2.2, 3.1–3.2, 5.2 and 5.5. **M1** guards automatic home service rollback. **M2** adds the typed companion and its inert heartbeat and catalogue, with a fake walkthrough. **M3** pins kernel, lease, readiness and cancel boundaries. Appendix C queues terminal exclusion and the remaining gates before dispatch.
>
> **Size.** Expected cumulative changes **114 → 1232 → 1400**. The full terminal-exclusion draft measures **1748**. The 1,400-line cap fires at the S4b boundary. Design budget 880 raw / 1,188 calibrated.
>
> **Authority.** D-WORLD-119 A and D-WORLD-120 A, attended by Mark in 82f05004. Controller E1 and E3 passed under GUI UID 501. S1 #606, S2 #609, S2b #611 and S3 #643 have landed. S3 merge e83d1717, judge PASS 86. Design r1 remains ratified at quorum r2 PROCEED 3/3 and independent audit PASS 81.

## 0. Binding constraints

- D-WORLD-109 limits the tests to fixture stores, disposable directories and test keys. New service tests use a fake controller, fake terminal and in-process HTTP transport. The controller's two observation legs inspect real service files and query supervisor job count read-only.
- The internal process accepts its fixed zero-argument mode. It holds a singleton kernel lock and publishes liveness. No action dispatch, mint, agent spawn, finish, publish or listener. The attended cancel command can cancel an unstarted row.
- Keep daemon routes, serving code, app.js, surfaceCSP, .ail, launchd files, skills, scripts, packages, integrations, plugins and world-publish byte-identical. Respect #646 and #617 collision paths. Apply exactly the M1–M3 fences.

## 1. Pre-flight

### 1.1 Baseline at f074dd49

````text
LEG observation-before rc=77 0s UNINFORMATIVE UNDER SANDBOX
LEG paths rc=0 0s
TOUCHED_PACKAGES: 
LEG touched rc=0 0s
LEG status rc=0 1s
LEG gofmt rc=0 0s
LEG build rc=0 4s
LEG vet rc=0 1s
LEG fence rc=0 1s
LEG verify_ail rc=0 43s
LEG pkg-cmd-ailang-worldd rc=1 4s UNINFORMATIVE UNDER SANDBOX
LEG race-cmd-ailang-worldd rc=1 18s UNINFORMATIVE UNDER SANDBOX
LEG pkg-host-store rc=0 15s
LEG race-host-store rc=0 41s
LEG inv-host-approval rc=0 1s
LEG inv-host-approval-approvaltest rc=0 0s
LEG inv-host-authority rc=0 1s
LEG inv-host-boundary rc=0 1s
LEG inv-host-daemon rc=0 1s
LEG inv-host-verbs rc=0 1s
LEG inv-cmd-ailang-worldd rc=0 1s
LEG inv-host-store rc=0 1s
LEG routes rc=0 0s
LEG scope rc=0 0s
LEG callers rc=0 0s
LEG verify_go rc=0 0s
LEG observation-after rc=77 0s UNINFORMATIVE UNDER SANDBOX
GATES_DONE 135s FAILS: observation-before_sandbox_uninformative pkg-cmd-ailang-worldd_sandbox_uninformative race-cmd-ailang-worldd_sandbox_uninformative observation-after_sandbox_uninformative
````

The verify_go leg shown above was intentionally skipped using GATES_SKIP_VERIFY_GO=1. The ledger and experiments are from Voights-Mac-Studio.local. The gate uses a pristine git archive under plan-scratch and an isolated AILANG_STATE_DIR. The whole CLI plain and race tests encounter a denied httptest loopback bind. The controller runs the exact commands outside the sandbox.

### 1.2 Codebase claims checked

| # | Claim | Command | Observed |
|---|---|---|---|
| C1 | HEAD, predecessor and authority | `git rev-parse HEAD`, `git show -s --format='%h %s' HEAD e83d1717 82f05004`, `git show 82f05004 --stat` | f074dd49e8fe1f281792de4fb475e20b53cd1d92 is merge #645, after S3 #643. 82f05004 records Mark's attended D-WORLD-119=A and D-WORLD-120=A. Controller supplied empty open queue. |
| C2 | E1/E3 identity, E2 owner | Controller design §11 transcript and `cat ~/.ailang/state/world-iter429/e3.out` | E1 rc 0 on gui/501/com.sunholo.ailang-world.worldd, UID 501 and expected installed binary. Mission-world PID is the positive control. E3 rc 0, two receipt lines, both synthetic World paths absent before/after, cleanup count 0, real plist mtime Nov 21. E3 demonstrates environment and restart reentry only. W33 E2 read-level PASS belongs to S9/S11. |
| C3 | W04 service analogues | `rg -n 'const worldServiceLabel\|^func (realServiceEnv\|homeServicePaths\|renderServicePlist\|readServiceRecord\|runServiceWith\|validateServicePlist)\|func \(e serviceEnv\)' cmd/ailang-worldd/service.go`, `rg -n 'serviceController\|^func' cmd/ailang-worldd/{service_launchctl,service_upgrade,doctor_service}.go` | service.go:29,48,85,113,144,167,186,199,220,234,260,277,572,639,655,678,686 pins label, paths, env, atomic writes, lock, typed confirmation, plist, bounded health, rotation, boot, record, flags, status, held probe, path safety, loopback and plist parsing. service_launchctl.go:17,24,37,49,59–78 pins interface, backend, budgets and calls. doctor_service.go:12 supplies matching-service diagnostics. |
| C4 | W10/W24 schema and rollback | `rg -n 'CredentialsSchemaVersion\|OpenCredentials\|enforceSchema\|CanonicalPath' host/store/credentials.go`, `nl -ba cmd/ailang-worldd/service_upgrade.go`, `rg -n 'ensureCredentialsStore\|CredentialsDB:' host/daemon/{daemon,credentials}.go` | credentials.go:45,52,56 current 10 and floor 9, :176,219,230,271,334–356 canonical path, writer, migration, reader and newer-schema refusal. service_upgrade.go:21,34,40,95,125,210,243 version wire shape, sanitized probe, health wait, upgrade, swap and rollback. daemon.go:817,1197 and daemon/credentials.go:59–77 migrate before health. version.go:15–39 carries build/platform/tool/interpreter, no schema capability at base. |
| C5 | §§3.1–3.2 schema and CAS | `nl -ba host/store/operator_schema.sql`, `rg -n 'OperatorOwner\|operatorAccess\|operatorCASRow\|ClaimOperatorAction\|operator_runner_state\|operatorChangedOne' host/store/operator_{actions,cas}.go` | operator_schema.sql:34–44 has singleton PK (must be 1), integer generation, text owner, integer protocol, text build, integer seen_at, text readiness, integer catalogue_generation, nullable BLOB catalogue. :48–64 has PK(kind,identity) and step FK. operator_cas.go:21,32,42,116,136–149,154 owner, id/digest/generation CAS, confirm, claim, freshness/readiness and exactly one row. operator_actions.go:40,50,60,74,136,155,163 v10+deadline gate, scan, projected read, request, cap, expiry, prune. |
| C6 | §3.1 W15/W36 keys and claimed receipt | `rg -n 'sshApprover\|passkeyApprover\|JudgeVerbSSH\|JudgeVerbAssertion\|RecheckVerb' host/authority/judge_verb.go`, `rg -n 'JudgeOperatorAction\|OperatorConfirmation\|RecheckClaimedOperatorAction' host/authority/judge_operator.go`, `rg -n 'consume\|RecheckVerb' host/daemon/{surface_confirm,daemon}.go` | V3 judge_verb.go:30,44,63,82,107. V4 judge_operator.go:61,121–142 gives judge and proof accessor/recheck, correlates frozen bytes, persisted row and current enrollment. No daemon status input. V3 confirm/consume at surface_confirm.go:37–39,97–207; authority.RecheckVerb is the legacy caller control. |
| C7 | W10/W13/W14/W32 terminal and kernel locks | `git grep -n 'LOCK_NB' -- cmd host`, `rg -n 'func (runEpisodeStart\|runPresetEpisodeStart\|runEpisodeResume\|runEpisodeFinishWith\|launchEpisode\|launchLockCheck)\|O_EXCL\|Getpid' cmd/ailang-worldd/episode*.go` | writer_lock_unix.go:35,41 exclusive non-waiting, :77,86 shared probe. service.go:154,648 non-waiting. episode.go:194, episode_preset.go:53, episode_resume.go:24, episode_finish.go:63 and episode_launch.go:230,258 are terminal sites. launch check :168–184 and :382,442–447 create/write a distinct PID file. No common episode flock exists at base. |
| C8 | Inventories | `rg -n '^func Test.*(Allowlist\|ProductionContextRoots\|ProductionGoSurface\|SurfaceTakesContext\|CredentialWriterSites\|SeamOnly\|WriterLock\|VerbTable\|EveryVerb\|QuickstartService\|HelpIsNeverBlocked)' cmd host -g '*test.go'` and §6 inv legs | daemon_test.go:899 dependency allowlist; host/boundary/allowlist_world_test.go:758,1121 dependency groups/protected packages; context_roots_test.go:129,234 roots/surface; deadline_guard_test.go:43 reflects *Store. credentials_test.go verifies store deadlines/DDL; cli_credential_writers_test.go:18 four base opens, surface_approve_post_guard_test.go:132 daemon opens. service_test.go:359 sole launchctl exec, :156–192 plist, :377 version. S4 adds no daemon imports, module, context root, package, new *Store method or direct launchctl call. |
| C9 | §8 public help and QUICKSTART pins | `nl -ba cmd/ailang-worldd/{help_test,verbs_parity_test,quickstart_service_test,plugin_help_test,keyguard_test}.go`, `rg -n '^## Run World as a service\|^## Run your agents through World' docs/QUICKSTART.md` | QUICKSTART :23–57 pinned by literal heading, at least ten command lines, unknown verb/flag controls and help parse. plugindoc.UsagePaths extracts lowercase words in top-level usage. Existing public table scope is 29 switch words. Internal startup mode precedes the switch. Help includes all new service runner subcommands, and keyguard tests check both help forms. |
| C10 | Reused fixtures and symbols for patches | `rg -n 'func (newServiceFixture\|upgradeFixture\|openTestCredentials\|testCtx\|operatorClaimFixture\|newFinishFixture\|newEpisodeRun\|credSQL\|ojSSH\|ojAssert\|budgetContext\|within\|isHelpToken)' cmd host`, `rg -n 'type (serviceEnv\|servicePaths\|serviceRecord\|sessionEnv\|HealthResponse\|BuildIdentity) struct' cmd host/daemon` | service_test.go:74,123 fake home/controller/transport, service_upgrade_test.go:26 upgrade prompt. service.go:44–83 types, session.go login/clock, daemon.go:623–658 daemon status/store/build, build.go:27,93 identity/current. cli.go:90 budgets, setup.go:172 containment, keyguard.go help tokens. journal_test.go:121, mint_requests_test.go:175, operator_cas_test.go:9, operator_actions_test.go:15,28 store fixtures. episode_finish_test.go:58,241, episode_test.go:84–99 and episode_finish_surface_test.go:75 terminal cases. authority/operator_store_test.go:24–53 uses judge and real signer/assertion fixtures. |
| C11 | Judge-428 cap, identity, order and size survivors | `nl -ba host/store/operator_actions.go`, `nl -ba host/store/operator_schema.sql`, `rg -n 'MaxSurfaceApprovalArtifactBytes\|ApprovedBy !=' host/store/operator_cas.go` and prior judge/mutants.tsv rows 17,18,28,48,50 | :121 prefers active then most recent. :139–140 includes explicit ID collisions and confirmed pending. PK(kind,identity) remains cross-action, confirm artifact size at :82. Predicates still match the judged behavior. See Appendix C for individual dispositions. |
| C12 | Judge-428 v9 floor, FK and module survivors | `rg -n 'CredentialsBindingWriteVersion\|foreign_keys\|golang.org/x/sys' cmd/ailang-worldd/{episode_preset,session_new}.go host/store/credentials.go go.mod host/authority/judge_operator_silent_test.go` | v9 binding floor used at episode_preset.go:241, session_new.go:381 and credentials.go:483. Judge J-FLOOR-MINT mutated the store site, so carry it in addition to both command sites. FK set at :312 and in-memory :320. The module remains indirect at go.mod:17 despite the direct test import at judge_operator_silent_test.go:17. |
| C13 | Replay and current production call fence | `nl -ba host/store/operator_cas.go`, `git grep -nF '.ConfirmOperatorAction(' -- cmd host ':!*_test.go' \| grep -vF 'approval.ConfirmOperatorAction('`, `git grep -n 'authority.RecheckVerb(' -- cmd host ':!*_test.go'`; extended real signer replay at scratch | :55–59 compares digest and returns an approved row. Zero production store-confirm callers, one pure approval call (excluded) and one V3 recheck control. Both altered fixture artifacts returned nil before the follow-up guard, with unchanged durable row. The same-byte-only draft adds 13 lines. |
| C14 | Routes, serving, inert callers | `grep -Ec 'mux\.(HandleFunc\|Handle)\(' host/daemon/daemon.go`, `grep -Fc 'mux.HandleFunc("POST /surface/' host/daemon/daemon.go`, `grep -n surfaceCSP host/daemon/surface_page.go` and §6 fences | 50 registrations, 3 surface POSTs, CSP constant at :19. Daemon and app bytes identical. Heartbeat publication/query calls are inside operator_runner.go, with doctor through status. |
| C15 | Snapshot, cap, CI | `cat ~/.ailang/state/world-iter429/{commit_ms,ctl_roundtrip}.sh`, `rg -n 'VERIFY_GO_LEGS: race\|timeout 8m' .github/workflows/ci.yml scripts/verify_go.sh`, scratch `git diff --no-index -M --numstat <base> <stage>` | Controller commits only after gofmt/build/vet/test and prints base-relative numstat. Round trip loops M1–M3 and compares cmd/host. M2 includes docs, so the planner compares full trees. CI race per package 8m. |
| C16 | New symbols | Append P and `grep -n 'func\|type\|const' cmd/ailang-worldd/operator_{service,runner}.go host/store/operator_runner.go` on final scratch | New typed config, label, private path/identity/config/pair functions, immutable rendering, lifecycle and process locks, runner loop/tick/status/write/cancel, bounded DB URI, protocol and heartbeat helpers. Existing renderer gets optional ordered environment; serviceVersion gets credentials capability; serviceEnv gets login injection. Only parked S4b adds terminalEpisodeLock, operatorEpisodeLockPath and OperatorOwnedEpisode. |

| C17 | Judge-428 source continuity and controls | `git diff --name-only 786f6743 f074dd49`, and plan-scratch/premises.log | Changes after judged S3 include episode launch, hosts awareness, mcp stdio, daemon usage/trajectory and mission documents. Positive changed paths are listed in the log. None are in the carried predicate files from C11/C12. Actual source predicates were re-read at HEAD. |
| C18 | Store-open contention and clock | `nl -ba host/store/writer_lock.go \| sed -n '169,205p;256,275p'` and source in M2/M3 | Landed busy default is 2000 ms. The runner sets 250 ms and requires an existing file via mode=rw. The read-only opener replaces mode with ro. Each method has a bounded operation context. Test clocks control Unix tick, exact-current, future, 9-second fresh, and 10-second stale values without sleeping. S4b separately uses a bounded synthetic start timeout. |

| C19 | Carried whole-package replay at actual HEAD | `python3 $HOME/.ailang/state/world-iter429/plan-scratch/carried/recheck.py` | Whole host/store plus host/authority survive J-REQ-PENDING-CONF, J-REQ-ACTIVE-ID, J-REQ-DEDUP-ORDER, J-DDL-STEPS-PK, J-CONF-ARTSIZE, J-FLOOR-MINT and J-FK-INMEM at f074dd49. Both command-floor whole attempts are UNINFORMATIVE UNDER SANDBOX at the existing bind denial. The exact command is go test -count=1 -timeout 8m ./cmd/ailang-worldd against each cmd-floor variant. See carried/ledger.tsv. |

### 1.3 Planner decisions

The LaunchAgent, file permissions, kernel lock and SQLite write are host effects. Keep those operations in the existing CLI and store packages. Keep the readiness comparison pure and the AILANG kernel unchanged.

- **S4-D1.** Mirror the home service's typed prompt, private atomic writes, explicit plist/program, status and fake controller. The companion uses a separate label, a config digest, validated account HOME, appends the system PATH floor, and has no effect child. Install stages files, start bootstraps, uninstall keeps the spool and records.
- **S4-D2.** One fixed owner lease per active credentials-store kernel lock. Boot and 2 s polls publish complete generations. Freshness requires seen_at ≤ now and seen_at > now−10. Start readiness is a fixture from the validated pair. Finish is blocked pending E2. The catalogue is configured preset names; eligibility and Git reads wait on S5's pinned inputs.
- **S4-D3.** Version JSON advertises the migration ceiling. Candidate launch may migrate before health. Unknown candidate capability, lower previous capability, a missing/nonregular previous binary or wrong digest all refuse automatic rollback. Retain the current candidate and previous path with an attended remedy. Equal-schema rollback keeps existing history. While a companion is installed, home upgrade refuses under the lifecycle lock; attended stop/uninstall is required before the swap. Paired replacement and child drain remain gates before effects.
- **S4-D4.** Put terminal exclusion in a measured follow-up. Losing an episode kernel lock cannot clear a durable owner. With no production consumer, S4 cannot claim an action or make a new reservation.

### 1.4 Premises found false or incomplete at HEAD

| # | Premise | Measurement and consequence |
|---|---|---|
| F1 | Heartbeat already has a producer and can release ownership | C5/C14. Claim is a reader of S3's row; no producer or runner existed. Publish liveness only. |
| F2 | Existing PID/credentials lock excludes every terminal effect | C7. Launch-only O_EXCL and short credentials writer are distinct, neither shared across CLI commands. S4b carries the common lock. |
| F3 | There is an explicit rollback verb, or auto rollback is schema-safe | C4. No standalone verb. Candidate migrates before health, rollback unconditionally restored .prev and version JSON lacked capability. M1 guards it. |
| F4 | Claimed receipt recheck verifies current daemon health | C6. Inputs are frozen action, row, challenge and time. Pair check is separate. S5 must use both at effect admission. |
| F5 | S3b landed because the record says “with S3b” | `git ls-files host/store/operator_steps.go host/store/operator_cas.go` prints only operator_cas.go. The original 132-line patch re-applies and compiles; full fd/predicate/current-owner work remains. |
| F6 | Different confirmation bytes are refused as replay | C13. Real signer fixture red with different replay artifact returned <nil>, twice, with unchanged row. |
| F7 | One fake supervisor loaded bit describes both jobs | C3, service_test.go:52. A fixture with a loaded base service must distinguish the companion's supervisor label. |
| F8 | Darwin t.TempDir path spellings are always canonical | EvalSymlinks gave /var/folders versus /private/var/folders. Canonicalize injected fixture HOME and check production account ancestors separately. |
| F9 | Existing controller round trip includes every S4 directory | C15. Keep M1–M3 and include docs for M2. |

| F10 | Every .orig in a full archive is a failed patch | `git ls-files '*.orig' '*.rej'` returned two tracked iter362 evaluator backup fixtures. Full-tree equality includes those bytes. Rejection/backup checks target the edited cmd, host and docs trees. |

### 1.5 Size and CI race cost

Controller measure `git diff -M --numstat f074dd49 HEAD -- cmd host docs app`. Count additions plus deletions. Scratch measured clean archive trees with `git diff --no-index -M --numstat` and inspected the unfiltered roots. Unfiltered product total is **1400** (**1381 additions, 19 deletions**). The two plan files are separate controller inputs. Increment counts use the previous stage. They can exceed the increase in cumulative numstat when a fixture added in one stage is rewritten in the next.

| Stage | Increment | Expected cumulative numstat | STATUS / DELETED / MANIFEST |
|---|---:|---:|---|
| M1 | 114 | 114 | 4 / 0 / 4 |
| M2 | 1118 | 1232 | 17 / 0 / 17 |
| M3 | 238 | 1400 | 18 / 0 / 18 |
| S4b, separate proposal | 366 | 1748 | Outside iteration 429 execution |

| Package | Base race package time | Final race package time | Added | Measurement |
|---|---:|---:|---:|---|
| host/store | 39.955 s | 33.256 s | -6.699 s | Whole package passed. 446.744 s to CI's 480 s timeout on this run. |
| cmd/ailang-worldd | 3.698 s | 4.982 s | +1.284 s | Socket-free service and inventory selection, both passed. Whole race uninformative. |

Wall-clock measurements were 40.61 s / 34.83 s for the whole store and 19.15 s / 10.30 s for the CLI selections. The whole CLI race attempts stopped at the denied httptest bind. The controller measures that whole package outside the sandbox. Exact focused commands and output are in the plan-scratch race logs.

Build, vet, gofmt, compile fence, whole store plain/race, and all named inventories and route/caller/path fences passed on the exact final stage. The inventory logs report 56 Background roots, 0 TODO and 5 WithoutCancel, and 25 read-only CLI store sites. Both AILANG gates passed with 16 required identities, 40 named core tests, 503 se-tools tests and PUB011. verify_go was intentionally skipped under GATES_SKIP_VERIFY_GO=1. The controller runs the plain repo-wide leg. The sandbox final sentinel reads:

````text
GATES_DONE 151s FAILS: none UNINFORMATIVE: observation-before pkg-cmd-ailang-worldd race-cmd-ailang-worldd observation-after
````

The actual per-leg logs are under plan-scratch/base-gates and plan-scratch/gates-exact-final. The baseline was taken with the earlier sentinel formatting, which listed incomplete sandbox legs in FAILS. Both baseline and final distinguish the bind/observation results from genuine failures in each leg's status.

## 2. How the executor applies a sub-milestone

Extract exactly M1, M2 or M3. A failed dry run, patch rejection or unexpected path stops the executor.

````bash apply.sh
#!/bin/bash
set -euo pipefail
M=M1 # choose M1, M2, M3 in order
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter429-row243-v4s4
PLAN="$W/design_docs/planned/w-world-address-bar-verbs-v4s4-sprint-plan.md"
T="$HOME/.ailang/state/world-iter429/exec-patches"
mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
/usr/bin/patch -p1 -E --dry-run < "$T/$M.patch"
/usr/bin/patch -p1 -E < "$T/$M.patch"
if find cmd host docs \( -name '*.orig' -o -name '*.rej' \) -print | grep .; then exit 1; fi
git status --short --untracked-files=all
````

## 3. Milestones

````bash accept.sh
#!/bin/bash
set -euo pipefail
cd "$1"
M=$2
test "${3:-fast}" = fast
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export AILANG_STATE_DIR="$HOME/.ailang/state/world-iter429/accept-ailang-state"
mkdir -p "$AILANG_STATE_DIR"
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR WORLD_UPDATE_OPERATOR_ACTION_GOLDEN WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN
test -x "$AILANG_BIN"
cmd='TestOperatorServiceUpgrade|TestHomeServiceUpgrade.*|TestHomeVersionJSON'
store='TestProductionContextRoots|TestProductionGoSurface|TestOperatorSchema10'
case "$M" in
  M1) ;;
  M2|M3|S4b)
    cmd="$cmd|TestOperatorService.*|TestOperatorStatusDisconnected|TestHomeService.*|TestEveryVerbHasHelp|TestHelpIsNeverBlockedByTheRegistryKeyGuard|TestVerbTableMatchesTheCLISwitch|TestCLICredentialWriterSites|TestQuickstartServiceSectionMatchesTheCLIHelp|TestDoctorServiceCheck"
    ;;
  *) exit 2 ;;
esac
if [ "$M" = M3 ] || [ "$M" = S4b ]; then
  cmd="$cmd|TestOperatorRunnerKernelLock"
  store="$store|TestOperator(RunnerUnavailable|CatalogueAdmission|FinishReadiness|CancelInflight|TerminalExclusion|ClaimRace)"
fi
if [ "$M" = S4b ]; then cmd="$cmd|TestOperatorTerminalExclusion(Start)?"; fi
go test -v -count=1 -timeout 8m ./cmd/ailang-worldd -run "^($cmd)$"
go test -v -count=1 -timeout 8m ./host/store -run "^($store)$"
echo "ACCEPT_OK $M"
````

### M1 — compatible home rollback (114 cumulative)

Extend version JSON and guard the pre-rename rollback seam. Fake failures prove unknown/lower/missing/linked/substituted previous builds leave the candidate, service record and disposable v10 store intact. Existing tests cover equal capability, rollback success and previous restart failure.

### M2 — companion lifecycle, status and inert publication (1232 cumulative)

Reuse confirmation, renderer, three-part diagnostics and process controller. Pin separate label, exact configuration, private spool, numeric loopback, canonical store pair and both digests. Attended cancel writes only unstarted status. QUICKSTART fake walkthrough and doctor status land with the public CLI.

### M3 — kernel, heartbeat, readiness and cancellation (1400 cumulative)

Exercise the real kernel open-description lock, future/current/exact-stale clocks, immutable catalogue generation, stale writer refusal, finish readiness at CAS and preservation of in-flight/reserved effects.

### 3.1 Branch traces

| Boundary | Positive readback | Refusal and persistence |
|---|---|---|
| Upgrade | Existing compatible fake restart restores the previous build and history. | Missing or incompatible capability, missing/linked/changed previous preserves the installed candidate and record/DB hashes, with no second kick. |
| Install and config | Typed y approves and persists exact config, renderer output and file modes. | Negative approval, bad config/backend, UID/HOME, writable components and changed executables refuse before boot. Read back fake controller calls and allowed write paths. |
| Pair startup | Same canonical paths and build, migrated v10, fake loaded job. | World, credentials or build mismatch, nonloopback, down daemon and bootstrap failure yield diagnostics. Bootstrap failure attempts companion stop. S4b waits for a new lease after bootstrap and stops its job on timeout. |
| Heartbeat/status | Known last-seen and expected lease; fresh start-only readiness. | Future/exact-stale, identity mismatch and disconnect refuse. Renewals compare owner, generation, protocol, build and monotonic time. Failed CAS keeps the earlier row. |
| Catalogue | Complete immutable generation readback. | Stale publication cannot change previous bytes, time or generation. No route calls it. |
| Readiness and cancel | Start CAS succeeds with a fresh start-only row. Unstarted cancel increments once. | Finish CAS refuses without finish readiness. In-flight cancel leaves action and reserved step. |

| Uncovered source predicates | See §4 named survivors and Appendix C. | The next fake fixture slice must distinguish approval-time config changes, exact supervisor labels, staged-file I/O errors, and every remaining file/path/mode branch. These are approval and deployment gates. |

## 4. Mutation table

K = behavioral kill, S = survived, E = explained equivalent, NM = unmeasured in sandbox. A named K with whole NM remains unmeasured for release. Static inventories do not count as behavioral evidence.

**Final discrimination.** 156 valid scoped variants. 31 killed by the named behavioral test and the whole package, 2 named equivalents, 2 whole-test survivors and 121 whole-package sandbox-unmeasured. The named tests killed 80, and 76 named tests survived.

The ledger is at `$HOME/.ailang/state/world-iter429/plan-scratch/mutants/ledger.txt`. Full old/new anchors and named/whole return codes and killers are in ledger.jsonl. Repeated ids retain their history in the ledger. This table uses their latest verdict. No passing behavioral result is credited from an inventory or string census.

Every valid named run used `go test -count=1 -timeout 30s ./<package> -v -run '^<NamedTest>$'`. Each whole run used `go test -count=1 -timeout 8m ./<package>` on the same sole mutation. Set `AILANG_BIN=$HOME/.pinned-ailang/ailang` for each run. For NM entries the controller re-applies the exact listed old/new change on the indicated stage and reruns the whole package command outside the sandbox. The loopback failure was in the existing CLI suite and does not evaluate the mutation.

| id | Diff anchor | Mutation | Named test | Named | Whole | Verdict |
|---|---|---|---|---|---|---|
| ROLLBACK-UNKNOWN-CANDIDATE | cmd/ailang-worldd/service_upgrade.go:264 | `candidateSchema <= 0 &#124;&#124; ` → `` | TestOperatorServiceUpgrade | K | NM | NM |
| ROLLBACK-CAPACITY | cmd/ailang-worldd/service_upgrade.go:264 | `previousSchema < candidateSchema` → `previousSchema+1 < candidateSchema` | TestOperatorServiceUpgrade | K | NM | NM |
| ROLLBACK-EQUAL | cmd/ailang-worldd/service_upgrade.go:264 | `previousSchema < candidateSchema` → `previousSchema <= candidateSchema` | TestHomeServiceUpgradeRollback | K | NM | NM |
| ROLLBACK-LSTAT | cmd/ailang-worldd/service_upgrade.go:264 | `previousErr != nil &#124;&#124;` → `(previousErr != nil && false) &#124;&#124;` | TestOperatorServiceUpgrade | K | NM | NM |
| ROLLBACK-REGULAR | cmd/ailang-worldd/service_upgrade.go:265 | `!previousInfo.Mode().IsRegular() &#124;&#124;` → `(!previousInfo.Mode().IsRegular() && false) &#124;&#124;` | TestOperatorServiceUpgrade | K | NM | NM |
| ROLLBACK-HASH | cmd/ailang-worldd/service_upgrade.go:265 | ` &#124;&#124; serviceFileHash(previous) != from` → `` | TestOperatorServiceUpgrade | K | NM | NM |
| PAIR-UPGRADE-INSTALLED | cmd/ailang-worldd/service_upgrade.go:206 | `[]string{e.runnerPaths().Record, e.runnerPaths().Plist}` → `[]string{}` | TestOperatorServiceUpgrade | S | NM | NM |
| PAIR-UPGRADE-SUPERVISOR | cmd/ailang-worldd/service_upgrade.go:211 | `!errors.Is(err, errServiceNotLoaded) { ↵   return refuse("cannot verify` → `err == nil && false { ↵   return refuse("cannot verify` | TestOperatorServiceUpgrade | S | NM | NM |
| PAIR-UPGRADE-LOCK | cmd/ailang-worldd/service_upgrade.go:201 | `runnerLifecycleLock(filepath.Dir(e.paths.Record))` → `os.Open(filepath.Dir(e.paths.Record))` | TestOperatorServiceUpgrade | S | NM | NM |
| VERSION-SCHEMA-JSON | cmd/ailang-worldd/version.go:30 | `CredentialsSchema: store.CredentialsSchemaVersion` → `CredentialsSchema: store.CredentialsSchemaVersion * 0` | TestHomeVersionJSON | K | NM | NM |
| HB-READ-MISSING | host/store/operator_runner.go:36 | `err = ErrOperatorRunner` → `err = nil` | TestOperatorRunnerUnavailable | K | K | K |
| HB-PROTOCOL | host/store/operator_runner.go:22 | `h.Protocol == OperatorRunnerProtocol && ` → `` | TestOperatorRunnerUnavailable | K | K | K |
| HB-BUILD | host/store/operator_runner.go:22 | ` && h.Build == build` → `` | TestOperatorRunnerUnavailable | K | K | K |
| HB-GENERATION-POSITIVE | host/store/operator_runner.go:23 | `h.Owner.Generation > 0` → `h.Owner.Generation >= 0` | TestOperatorRunnerUnavailable | K | K | K |
| HB-NONCE | host/store/operator_runner.go:23 | ` && h.Owner.Nonce != ""` → `` | TestOperatorRunnerUnavailable | K | K | K |
| HB-FUTURE-REMOVE | host/store/operator_runner.go:24 | `h.SeenAt <= now && ` → `` | TestOperatorRunnerUnavailable | K | K | K |
| HB-NOW-EQUAL | host/store/operator_runner.go:24 | `h.SeenAt <= now` → `h.SeenAt < now` | TestOperatorRunnerUnavailable | K | K | K |
| HB-STALE-EQUAL | host/store/operator_runner.go:24 | `h.SeenAt > now-10` → `h.SeenAt >= now-10` | TestOperatorRunnerUnavailable | K | K | K |
| HB-STALE-REMOVE | host/store/operator_runner.go:24 | ` && h.SeenAt > now-10` → `` | TestOperatorRunnerUnavailable | K | K | K |
| HB-READINESS | host/store/operator_runner.go:25 | `(h.Readiness == "both" &#124;&#124; h.Readiness == kind)` → `kind != ""` | TestOperatorRunnerUnavailable | K | K | K |
| HB-READINESS-BOTH | host/store/operator_runner.go:25 | `h.Readiness == "both" &#124;&#124; ` → `` | TestOperatorRunnerUnavailable | K | K | K |
| HB-INSERT-VS-RENEW | host/store/operator_runner.go:48 | `if h.Owner.Generation == 0` → `if h.Owner.Generation < 0` | TestOperatorCatalogueAdmission | K | K | K |
| HB-BOOT-GENINC | host/store/operator_runner.go:52 | `generation = generation+1` → `generation = generation+0` | TestOperatorCatalogueAdmission | K | K | K |
| HB-BOOT-CATINC | host/store/operator_runner.go:54 | `catalogue_generation = catalogue_generation+1` → `catalogue_generation = catalogue_generation+0` | TestOperatorCatalogueAdmission | K | K | K |
| HB-TICK-CATINC | host/store/operator_runner.go:60 | `catalogue_generation=catalogue_generation+1` → `catalogue_generation=catalogue_generation+0` | TestOperatorCatalogueAdmission | K | K | K |
| HB-TICK-NONCE | host/store/operator_runner.go:61 | `AND owner=? AND` → `AND ? IS NOT NULL AND` | TestOperatorCatalogueAdmission | K | K | K |
| HB-TICK-GEN | host/store/operator_runner.go:61 | `AND generation=? AND` → `AND ? IS NOT NULL AND` | TestOperatorCatalogueAdmission | K | K | K |
| HB-TICK-PROTOCOL | host/store/operator_runner.go:61 | `AND protocol=? AND` → `AND ? IS NOT NULL AND` | TestOperatorCatalogueAdmission | K | K | K |
| HB-TICK-BUILD | host/store/operator_runner.go:61 | `AND build=? AND` → `AND ? IS NOT NULL AND` | TestOperatorCatalogueAdmission | K | K | K |
| HB-TICK-PAST | host/store/operator_runner.go:61 | `AND seen_at<=?` → `AND ? IS NOT NULL` | TestOperatorCatalogueAdmission | K | K | K |
| HB-TICK-EQUAL | host/store/operator_runner.go:61 | `seen_at<=?` → `seen_at<?` | TestOperatorCatalogueAdmission | K | K | K |
| HB-TICK-REFUSAL | host/store/operator_runner.go:65 | `err = ErrOperatorRunner` → `err = nil` | TestOperatorCatalogueAdmission | K | K | K |
| HB-CATALOGUE-BYTES | host/store/operator_runner.go:60 | `catalogue_generation=catalogue_generation+1, catalogue=?` → `catalogue_generation=catalogue_generation+1, catalogue=substr(?,1,0)` | TestOperatorCatalogueAdmission | K | K | K |
| HB-SINGLETON-WHERE | host/store/operator_runner.go:61 | `WHERE singleton=1 AND owner=?` → `WHERE owner=?` | TestOperatorCatalogueAdmission | S | S | E |
| FINISH-CLAIM-READINESS | host/store/operator_cas.go:137 | `(readiness = 'both' OR readiness = ?)` → `(? IS NOT NULL)` | TestOperatorFinishReadiness | K | K | K |
| CANCEL-STATE | host/store/operator_runner.go:76 | ` AND state IN ('requested','confirmed')` → `` | TestOperatorCancelInflight | K | K | K |
| CANCEL-OWNER | host/store/operator_runner.go:77 | `AND owner IS NULL AND` → `AND` | TestOperatorCancelInflight | S | S | E |
| CANCEL-RESERVATION | host/store/operator_runner.go:77 | `AND NOT EXISTS (SELECT 1 FROM operator_action_steps ↵    WHERE action_id=? AND kind<>'operation')` → `AND ? IS NOT NULL` | TestOperatorCancelInflight | K | K | K |
| CANCEL-ID | host/store/operator_runner.go:76 | `WHERE id=? AND state` → `WHERE ? IS NOT NULL AND state` | TestOperatorCancelInflight | K | K | K |
| CANCEL-EXPIRY-EQUAL | host/store/operator_runner.go:77 | `expires_at>?` → `expires_at>=?` | TestOperatorCancelInflight | K | K | K |
| CANCEL-GENERATION | host/store/operator_runner.go:76 | `generation=generation+1` → `generation=generation+0` | TestOperatorCancelInflight | K | K | K |
| CANCEL-TIME | host/store/operator_runner.go:75 | `SET state='cancelled', terminal_at=?` → `SET state='cancelled', terminal_at=(?+1)` | TestOperatorCancelInflight | K | K | K |
| CANCEL-DEADLINE | host/store/operator_runner.go:72 | `if err := c.operatorAccess(ctx); err != nil {` → `if err := c.operatorAccess(ctx); err != nil && false {` | TestOperatorCancelInflight | K | K | K |
| CLI-RUNNER-DISPATCH | cmd/ailang-worldd/service.go:296 | `if verb == "runner"` → `if verb == "runner-disabled"` | TestOperatorServiceWriteScope | K | NM | NM |
| CONFIG-UNKNOWN | cmd/ailang-worldd/operator_service.go:105 | `d.DisallowUnknownFields()` → `// allow unknown` | TestOperatorServiceConfig | K | NM | NM |
| CONFIG-DECODE | cmd/ailang-worldd/operator_service.go:106 | `d.Decode(&c) != nil &#124;&#124; ` → `(d.Decode(&c) != nil && false) &#124;&#124; ` | TestOperatorServiceConfig | K | NM | NM |
| CONFIG-TRAILING | cmd/ailang-worldd/operator_service.go:106 | `d.Decode(new(any)) != io.EOF` → `(d.Decode(new(any)) != io.EOF && false)` | TestOperatorServiceConfig | K | NM | NM |
| CONFIG-SIZE | cmd/ailang-worldd/operator_service.go:101 | `len(b) > 16384` → `len(b) > 16385` | TestOperatorServiceConfig | S | NM | NM |
| CONFIG-NOFOLLOW | cmd/ailang-worldd/operator_service.go:91 | `os.O_RDONLY&#124;syscall.O_NOFOLLOW&#124;syscall.O_NONBLOCK, 0)` → `os.O_RDONLY&#124;syscall.O_NONBLOCK, 0)` | TestOperatorServiceConfig | S | NM | NM |
| CONFIG-REGULAR | cmd/ailang-worldd/operator_service.go:97 | `if err != nil &#124;&#124; !st.Mode().IsRegular()` → `if err != nil &#124;&#124; (!st.Mode().IsRegular() && false)` | TestOperatorServiceConfig | S | NM | NM |
| ACCOUNT-ROOT | cmd/ailang-worldd/operator_service.go:119 | ` &#124;&#124; uid == 0 &#124;&#124; uid != c.UID` → ` &#124;&#124; uid != c.UID` | TestOperatorServiceOwnership | K | NM | NM |
| ACCOUNT-UID | cmd/ailang-worldd/operator_service.go:119 | ` &#124;&#124; uid != c.UID` → `` | TestOperatorServiceOwnership | K | NM | NM |
| ACCOUNT-LOOKUP | cmd/ailang-worldd/operator_service.go:119 | `case err != nil &#124;&#124; uid == 0` → `case (err != nil && false) &#124;&#124; uid == 0` | TestOperatorServiceOwnership | S | NM | NM |
| ACCOUNT-APPROVED-HOME | cmd/ailang-worldd/operator_service.go:121 | `case home != c.Home &#124;&#124; home != e.paths.Home` → `case home != e.paths.Home` | TestOperatorServiceOwnership | K | NM | NM |
| ACCOUNT-SERVICE-HOME | cmd/ailang-worldd/operator_service.go:121 | ` &#124;&#124; home != e.paths.Home` → `` | TestOperatorServiceOwnership | S | NM | NM |
| CONFIG-BACKEND | cmd/ailang-worldd/operator_service.go:123 | `case c.Backend != "headless"` → `case c.Backend != c.Backend` | TestOperatorServiceConfig | K | NM | NM |
| CONFIG-PRESET-EMPTY | cmd/ailang-worldd/operator_service.go:125 | `len(c.PresetNames) == 0 &#124;&#124; ` → `` | TestOperatorServiceConfig | S | NM | NM |
| CONFIG-PRESET-BOUND | cmd/ailang-worldd/operator_service.go:125 | `len(c.PresetNames) > 16` → `len(c.PresetNames) > 17` | TestOperatorServiceConfig | S | NM | NM |
| CONFIG-EXEC-BOUND | cmd/ailang-worldd/operator_service.go:125 | `len(c.Executables) > 16` → `len(c.Executables) > 17` | TestOperatorServiceConfig | S | NM | NM |
| CONFIG-PATH-BOUND | cmd/ailang-worldd/operator_service.go:125 | `len(c.ToolDirs) > 16` → `len(c.ToolDirs) > 17` | TestOperatorServiceConfig | S | NM | NM |
| CONFIG-PRESET-NAME | cmd/ailang-worldd/operator_service.go:140 | ` &#124;&#124; !verbs.EpisodeIDPattern.MatchString(preset)` → `` | TestOperatorServiceConfig | S | NM | NM |
| CONFIG-PRESET-LENGTH | cmd/ailang-worldd/operator_service.go:140 | `len(preset) > verbs.MaxPresetName &#124;&#124; ` → `` | TestOperatorServiceConfig | S | NM | NM |
| CONFIG-PRESET-RESERVED | cmd/ailang-worldd/operator_service.go:140 | `preset == "se-tools" &#124;&#124; ` → `` | TestOperatorServiceConfig | S | NM | NM |
| CONFIG-LOOPBACK | cmd/ailang-worldd/operator_service.go:129 | `if !serviceLoopback(c.Address) &#124;&#124; !net.ParseIP(host).IsLoopback()` → `if !serviceLoopback(c.Address) && !net.ParseIP(host).IsLoopback()` | TestOperatorServicePair | S | NM | NM |
| CONFIG-NUMERIC-LOOPBACK | cmd/ailang-worldd/operator_service.go:129 | ` &#124;&#124; !net.ParseIP(host).IsLoopback()` → ` &#124;&#124; (net.ParseIP(host).IsLoopback() && false)` | TestOperatorServicePair | S | NM | NM |
| PATH-ABSOLUTE | cmd/ailang-worldd/operator_service.go:61 | `!filepath.IsAbs(path) &#124;&#124; ` → `` | TestOperatorServiceOwnership | S | NM | NM |
| PATH-CLEAN | cmd/ailang-worldd/operator_service.go:61 | ` &#124;&#124; filepath.Clean(path) != path` → `` | TestOperatorServiceOwnership | S | NM | NM |
| PATH-OWNER | cmd/ailang-worldd/operator_service.go:70 | `(int(owner.Uid) != uid && (within(home, p) &#124;&#124; p == home &#124;&#124; owner.Uid != 0))` → `(int(owner.Uid) != uid && false)` | TestOperatorServiceOwnership | S | NM | NM |
| PATH-PERMISSION | cmd/ailang-worldd/operator_service.go:73 | `st.Mode().Perm()&0022 != 0` → `st.Mode().Perm()&0002 != 0` | TestOperatorServiceOwnership | K | NM | NM |
| PATH-REGULAR | cmd/ailang-worldd/operator_service.go:77 | ` &#124;&#124; !wantDir && !st.Mode().IsRegular()` → `` | TestOperatorServiceOwnership | K | NM | NM |
| PATH-DIRECTORY | cmd/ailang-worldd/operator_service.go:77 | `wantDir && !st.IsDir() &#124;&#124; ` → `` | TestOperatorServiceOwnership | S | NM | NM |
| PATH-PRIVATE-MODE | cmd/ailang-worldd/operator_service.go:80 | `private != 0 && st.Mode().Perm() != private` → `private != 0 && (st.Mode().Perm() != private && false)` | TestOperatorServiceOwnership | S | NM | NM |
| PATH-LSTAT | cmd/ailang-worldd/operator_service.go:65 | `st, err := os.Lstat(p)` → `st, err := os.Stat(p)` | TestOperatorServiceOwnership | K | NM | NM |
| PATH-TOOLS-SEPARATOR | cmd/ailang-worldd/operator_service.go:157 | `strings.ContainsAny(dir, ":\n")` → `strings.ContainsAny(dir, "\n")` | TestOperatorServiceEnvironment | S | NM | NM |
| EXEC-REQUIRED | cmd/ailang-worldd/operator_service.go:144 | `if c.Executables[e.paths.Binary] == ""` → `if c.Executables[e.paths.Binary] == "" && false` | TestOperatorServiceOwnership | K | NM | NM |
| EXEC-MODE | cmd/ailang-worldd/operator_service.go:152 | ` &#124;&#124; st.Mode().Perm()&0100 == 0` → ` &#124;&#124; (st.Mode().Perm()&0100 == 0 && false)` | TestOperatorServiceOwnership | K | NM | NM |
| EXEC-HASH | cmd/ailang-worldd/operator_service.go:152 | ` &#124;&#124; serviceFileHash(path) != sum` → `` | TestOperatorServiceOwnership | K | NM | NM |
| EXEC-NONEMPTY-HASH | cmd/ailang-worldd/operator_service.go:152 | ` &#124;&#124; sum == ""` → `` | TestOperatorServiceOwnership | S | NM | NM |
| PAIR-RECORD-UNINSTALLED | cmd/ailang-worldd/operator_service.go:179 | ` &#124;&#124; r.Uninstalled` → `` | TestOperatorServicePair | S | NM | NM |
| PAIR-RECORD-ADDRESS | cmd/ailang-worldd/operator_service.go:179 | ` &#124;&#124; r.Address != c.Address` → `` | TestOperatorServicePair | S | NM | NM |
| PAIR-RECORD-DIGEST | cmd/ailang-worldd/operator_service.go:179 | ` &#124;&#124; r.BinarySHA256 != c.Executables[e.paths.Binary]` → `` | TestOperatorServicePair | S | NM | NM |
| PAIR-RECORD-WORLD | cmd/ailang-worldd/operator_service.go:184 | ` &#124;&#124; worldRecord != c.World` → ` &#124;&#124; (worldRecord != c.World && false)` | TestOperatorServicePair | S | NM | NM |
| PAIR-RECORD-CREDS | cmd/ailang-worldd/operator_service.go:184 | ` &#124;&#124; credsRecord != c.Credentials` → ` &#124;&#124; (credsRecord != c.Credentials && false)` | TestOperatorServicePair | S | NM | NM |
| PAIR-HEALTH-ERROR | cmd/ailang-worldd/operator_service.go:188 | `if err != nil &#124;&#124; h.Status != "ok"` → `if err != nil && h.Status != "ok"` | TestOperatorServiceStartup | S | NM | NM |
| PAIR-HEALTH-WORLD | cmd/ailang-worldd/operator_service.go:192 | ` &#124;&#124; world != c.World {` → ` &#124;&#124; (world != c.World && false) {` | TestOperatorServicePair | K | NM | NM |
| PAIR-HEALTH-CREDS | cmd/ailang-worldd/operator_service.go:196 | ` &#124;&#124; creds != c.Credentials {` → ` &#124;&#124; (creds != c.Credentials && false) {` | TestOperatorServicePair | K | NM | NM |
| PAIR-HEALTH-DIGEST | cmd/ailang-worldd/operator_service.go:199 | `if h.Build.ExecutableSHA256 != c.Executables[e.paths.Binary]` → `if h.Build.ExecutableSHA256 != c.Executables[e.paths.Binary] && false` | TestOperatorServicePair | K | NM | NM |
| PAIR-SCHEMA | cmd/ailang-worldd/operator_service.go:207 | `db.Version() != store.CredentialsSchemaVersion` → `db.Version() > store.CredentialsSchemaVersion` | TestOperatorServiceStartup | K | NM | NM |
| PLIST-APPROVED-HASH | cmd/ailang-worldd/operator_service.go:238 | `if serviceFileHash(p.Plist) != serviceHash([]byte(runnerPlist(p, c, b)))` → `if serviceFileHash(p.Plist) != serviceHash([]byte(runnerPlist(p, c, b))) && false` | TestOperatorServiceFence | K | NM | NM |
| PLIST-HOME | cmd/ailang-worldd/operator_service.go:216 | `[2]string{"HOME", c.Home}` → `[2]string{"HOME", os.Getenv("HOME")}` | TestOperatorServiceEnvironment | K | NM | NM |
| PLIST-PATH | cmd/ailang-worldd/operator_service.go:216 | `strings.Join(paths, ":")` → `strings.Join(paths[:len(paths)-1], ":")` | TestOperatorServiceEnvironment | K | NM | NM |
| PLIST-TOKEN | cmd/ailang-worldd/operator_service.go:218 | `[2]string{"GH_PROMPT_DISABLED", "1"})` → `[2]string{"GH_PROMPT_DISABLED", "1"}, [2]string{"GH_TOKEN", os.Getenv("GH_TOKEN")})` | TestOperatorServiceEnvironment | K | NM | NM |
| PLIST-ARGS | cmd/ailang-worldd/operator_service.go:215 | `[]string{p.Binary, "operator-runner"}` → `[]string{p.Binary, "serve"}` | TestOperatorServiceEnvironment | K | NM | NM |
| PLIST-CONFIG-HASH | cmd/ailang-worldd/operator_service.go:217 | `[2]string{"WORLD_OPERATOR_CONFIG_SHA256", serviceHash(configBytes)}` → `[2]string{"WORLD_OPERATOR_CONFIG_SHA256", serviceHash(nil)}` | TestOperatorServiceEnvironment | K | NM | NM |
| INSTALL-APPROVAL | cmd/ailang-worldd/operator_service.go:306 | `if !e.confirm(errOut, change) {` → `if !e.confirm(errOut, change) && false {` | TestOperatorServiceFence | K | NM | NM |
| INSTALL-RECHECK | cmd/ailang-worldd/operator_service.go:338 | `if err := e.checkRunnerConfig(c); err != nil {` → `if err := e.checkRunnerConfig(c); err != nil && false {` | TestOperatorServiceOwnership | S | NM | NM |
| INSTALL-PAIR | cmd/ailang-worldd/operator_service.go:341 | `if err := e.checkRunnerPair(c); err != nil {` → `if err := e.checkRunnerPair(c); err != nil && false {` | TestOperatorServicePair | S | NM | NM |
| INSTALL-LOADED | cmd/ailang-worldd/operator_service.go:344 | `if _, err := e.controller.Print(e.domain, runnerLabel); !errors.Is(err, errServiceNotLoaded) {` → `if _, err := e.controller.Print(e.domain, runnerLabel); !errors.Is(err, errServiceNotLoaded) && false {` | TestOperatorServiceFence | S | NM | NM |
| INSTALL-IMMUTABLE | cmd/ailang-worldd/operator_service.go:351 | `if _, err := os.Lstat(path); !errors.Is(err, os.ErrNotExist) {` → `if _, err := os.Lstat(path); !errors.Is(err, os.ErrNotExist) && false {` | TestOperatorServiceFence | S | NM | NM |
| INSTALL-SPOOL-MODE | cmd/ailang-worldd/operator_service.go:364 | `os.Mkdir(spool, 0700)` → `os.Mkdir(spool, 0750)` | TestOperatorServiceWriteScope | K | NM | NM |
| INSTALL-PLIST-MODE | cmd/ailang-worldd/operator_service.go:381 | `e.atomic(p.Plist, []byte(runnerPlist(p, c, b)), 0600)` → `e.atomic(p.Plist, []byte(runnerPlist(p, c, b)), 0644)` | TestOperatorServiceWriteScope | K | NM | NM |
| INSTALL-CONFIG-MODE | cmd/ailang-worldd/operator_service.go:378 | `e.atomic(p.Record, b, 0600)` → `e.atomic(p.Record, b, 0644)` | TestOperatorServiceWriteScope | K | NM | NM |
| INSTALL-LOG-MODE | cmd/ailang-worldd/operator_service.go:371 | `e.atomic(p.Log, nil, 0600)` → `e.atomic(p.Log, nil, 0644)` | TestOperatorServiceWriteScope | K | NM | NM |
| INSTALL-SAFE-PATH | cmd/ailang-worldd/operator_service.go:348 | `if err := serviceSafePath(p.Home, path); err != nil {` → `if err := serviceSafePath(p.Home, path); err != nil && false {` | TestOperatorServiceWriteScope | S | NM | NM |
| INSTALL-ROOT-MODE | cmd/ailang-worldd/operator_service.go:358 | `filepath.Dir(p.Record), c.UID, true, 0700` → `filepath.Dir(p.Record), c.UID, true, 0` | TestOperatorServiceOwnership | S | NM | NM |
| INSTALL-SPOOL-CHECK | cmd/ailang-worldd/operator_service.go:367 | `c.Home, spool, c.UID, true, 0700` → `c.Home, spool, c.UID, true, 0` | TestOperatorServiceOwnership | S | NM | NM |
| INSTALL-LOG-CHECK | cmd/ailang-worldd/operator_service.go:375 | `c.Home, p.Log, c.UID, false, 0600` → `c.Home, p.Log, c.UID, false, 0` | TestOperatorServiceOwnership | S | NM | NM |
| INSTALL-ROLLBACK-CONFIG | cmd/ailang-worldd/operator_service.go:383 | `return errors.Join(err, os.Remove(p.Record))` → `return err` | TestOperatorServiceWriteScope | S | NM | NM |
| LIFECYCLE-LOCK-NB | cmd/ailang-worldd/operator_service.go:250 | `syscall.LOCK_EX&#124;syscall.LOCK_NB` → `syscall.LOCK_EX` | TestOperatorServiceFence | S | NM | NM |
| LIFECYCLE-PLATFORM | cmd/ailang-worldd/operator_service.go:266 | `e.goos != "darwin" &#124;&#124; e.label != worldServiceLabel` → `e.goos != "darwin" && e.label != worldServiceLabel` | TestOperatorServiceConfig | S | NM | NM |
| LIFECYCLE-INSTALL-ARGV | cmd/ailang-worldd/operator_service.go:271 | `args[1] != "--config"` → `(args[1] != "--config" && false)` | TestOperatorServiceConfig | S | NM | NM |
| LIFECYCLE-CANCEL-ID | cmd/ailang-worldd/operator_service.go:275 | ` &#124;&#124; !verbs.ValidOperatorActionID(args[1])` → `` | TestOperatorServiceConfig | S | NM | NM |
| LIFECYCLE-ARGV-ARITY | cmd/ailang-worldd/operator_service.go:279 | `if len(args) != 1 {` → `if len(args) < 1 {` | TestOperatorServiceConfig | S | NM | NM |
| LIFECYCLE-STOP-APPROVAL | cmd/ailang-worldd/operator_service.go:320 | `if !e.confirm(errOut, change) {` → `if !e.confirm(errOut, change) && false {` | TestOperatorServiceFence | S | NM | NM |
| LIFECYCLE-STOP-LABEL | cmd/ailang-worldd/operator_service.go:408 | `if err := e.controller.Bootout(e.domain, runnerLabel);` → `if err := e.controller.Bootout(e.domain, e.label);` | TestOperatorServiceWriteScope | S | NM | NM |
| LIFECYCLE-KEEP-SPOOL | cmd/ailang-worldd/operator_service.go:412 | `for _, path := range []string{p.Plist, p.Record} { ↵     e.note(path)` → `for _, path := range []string{p.Plist, p.Record, p.Log} { ↵     e.note(path)` | TestOperatorServiceWriteScope | K | NM | NM |
| LIFECYCLE-START-PAIR | cmd/ailang-worldd/operator_service.go:428 | `if err := e.checkRunnerPair(c); err != nil {` → `if err := e.checkRunnerPair(c); err != nil && false {` | TestOperatorServicePair | K | NM | NM |
| LIFECYCLE-BOOT-CLEANUP | cmd/ailang-worldd/operator_service.go:435 | `_ = e.controller.Bootout(e.domain, runnerLabel)` → `// leave a partially bootstrapped job` | TestOperatorServiceStartup | K | NM | NM |
| LIFECYCLE-BOOT-ERROR | cmd/ailang-worldd/operator_service.go:437 | `return errors.New("could not bootstrap the companion")` → `return nil` | TestOperatorServiceStartup | K | NM | NM |
| RUNNER-INTERNAL-ARGV | cmd/ailang-worldd/operator_runner.go:39 | `if len(args) != 0 {` → `if len(args) < 0 {` | TestOperatorServiceConfig | K | NM | NM |
| RUNNER-CONFIG-ENV | cmd/ailang-worldd/operator_runner.go:43 | `os.Getenv("WORLD_OPERATOR_CONFIG_SHA256") != serviceHash(b)` → `os.Getenv("WORLD_OPERATOR_CONFIG_SHA256") == serviceHash(b)` | TestOperatorServiceConfig | S | NM | NM |
| RUNNER-SELF-DIGEST | cmd/ailang-worldd/operator_runner.go:44 | `daemon.CurrentBuild().ExecutableSHA256 != c.Executables[e.paths.Binary]` → `daemon.CurrentBuild().ExecutableSHA256 == c.Executables[e.paths.Binary]` | TestOperatorServiceConfig | S | NM | NM |
| RUNNER-LOCK-NB | cmd/ailang-worldd/operator_runner.go:31 | `syscall.LOCK_EX&#124;syscall.LOCK_NB` → `syscall.LOCK_EX` | TestOperatorRunnerKernelLock | K | NM | NM |
| RUNNER-LOCK-UNIQUE | cmd/ailang-worldd/operator_runner.go:31 | `int(f.Fd()), syscall.LOCK_EX&#124;syscall.LOCK_NB` → `int(f.Fd()), syscall.LOCK_SH&#124;syscall.LOCK_NB` | TestOperatorRunnerKernelLock | K | NM | NM |
| RUNNER-LOCK-NOFOLLOW | cmd/ailang-worldd/operator_runner.go:22 | `&#124;syscall.O_NOFOLLOW` → `` | TestOperatorRunnerKernelLock | S | NM | NM |
| RUNNER-LOCK-OWNER | cmd/ailang-worldd/operator_runner.go:27 | ` &#124;&#124; int(st.Sys().(*syscall.Stat_t).Uid) != os.Getuid()` → `` | TestOperatorRunnerKernelLock | S | NM | NM |
| RUNNER-LOCK-MODE | cmd/ailang-worldd/operator_runner.go:27 | ` &#124;&#124; st.Mode().Perm() != 0600` → `` | TestOperatorRunnerKernelLock | S | NM | NM |
| RUNNER-LOCK-REGULAR | cmd/ailang-worldd/operator_runner.go:27 | ` &#124;&#124; !st.Mode().IsRegular()` → `` | TestOperatorRunnerKernelLock | S | NM | NM |
| RUNNER-POLL-INTERVAL | cmd/ailang-worldd/operator_runner.go:19 | `const runnerPoll = 2 * time.Second` → `const runnerPoll = 20 * time.Second` | TestOperatorStatusDisconnected | S | NM | NM |
| HB-READ-DEADLINE | host/store/operator_runner.go:29 | `if err := c.operatorAccess(ctx); err != nil {` → `if err := c.operatorAccess(ctx); err != nil && false {` | TestOperatorRunnerUnavailable | S | S | S |
| RUNNER-HALT-PAIR | cmd/ailang-worldd/operator_runner.go:80 | `if err := e.checkRunnerPair(c); err != nil {` → `if err := e.checkRunnerPair(c); err != nil && false {` | TestOperatorServicePair | K | NM | NM |
| RUNNER-CATALOGUE-REDACTION | cmd/ailang-worldd/operator_runner.go:86 | `json.Marshal(c.PresetNames)` → `json.Marshal(c)` | TestOperatorStatusDisconnected | K | NM | NM |
| HB-WRITE-DEADLINE | host/store/operator_runner.go:44 | `if err := c.operatorAccess(ctx); err != nil {` → `if err := c.operatorAccess(ctx); err != nil && false {` | TestOperatorCatalogueAdmission | S | S | S |
| RUNNER-GENERATION-ADVANCE | cmd/ailang-worldd/operator_runner.go:92 | `h.Owner.Generation = generation` → `h.Owner.Generation = generation+1` | TestOperatorStatusDisconnected | S | NM | NM |
| HOME-PLIST-DEFAULT | cmd/ailang-worldd/service.go:198 | `if len(environment) == 0 {` → `if len(environment) < 0 {` | TestHomeServiceInstallAndPlist | K | NM | NM |
| RUNNER-STATUS-DISCONNECT | cmd/ailang-worldd/operator_runner.go:116 | `return fmt.Errorf("daemon disconnected: %w", err)` → `return nil` | TestOperatorStatusDisconnected | K | NM | NM |
| HELP-RUNNER-KEYGUARD | cmd/ailang-worldd/keyguard.go:74 | `rest[1] == "runner"` → `rest[1] == "runner-disabled"` | TestHelpIsNeverBlockedByTheRegistryKeyGuard | K | NM | NM |
| RUNNER-STATUS-SUPERVISOR | cmd/ailang-worldd/operator_runner.go:119 | `return errors.New("runner disconnected")` → `return nil` | TestOperatorStatusDisconnected | S | NM | NM |
| RUNNER-MODE-DISPATCH | cmd/ailang-worldd/main.go:382 | `if rest[0] == "operator-runner"` → `if rest[0] == "operator-runner-disabled"` | TestOperatorServiceConfig | S | NM | NM |
| RUNNER-STATUS-STALE | cmd/ailang-worldd/operator_runner.go:122 | `if !h.Available(build, "start", now)` → `if !h.Available(build, "start", now) && false` | TestOperatorStatusDisconnected | S | NM | NM |
| DOCTOR-COMPANION-DIAGNOSTIC | cmd/ailang-worldd/doctor_service.go:48 | `if serviceAddressMatches(addr, r) { ↵   var status bytes.Buffer` → `if serviceAddressMatches(addr, r) && false { ↵   var status bytes.Buffer` | TestDoctorServiceCheck | S | NM | NM |
| RUNNER-STATUS-E2 | cmd/ailang-worldd/operator_runner.go:126 | `if !h.Available(build, "finish", now)` → `if h.Available(build, "finish", now)` | TestOperatorStatusDisconnected | K | NM | NM |
| CONFIG-READ-SOURCE-IO | cmd/ailang-worldd/operator_service.go:101 | `if err != nil &#124;&#124; len(b) > 16384 {` → `if (err != nil && false) &#124;&#124; len(b) > 16384 {` | TestOperatorServiceConfig | S | NM | NM |
| CONFIG-STORE-IDENTITY-DISTINCT | cmd/ailang-worldd/operator_service.go:166 | ` &#124;&#124; world == creds` → `` | TestOperatorServicePair | S | NM | NM |
| CONFIG-STORE-WORLD-CANONICAL | cmd/ailang-worldd/operator_service.go:166 | ` &#124;&#124; world != c.World &#124;&#124; creds != c.Credentials &#124;&#124; world == creds` → ` &#124;&#124; creds != c.Credentials &#124;&#124; world == creds` | TestOperatorServicePair | S | NM | NM |
| CONFIG-STORE-CREDS-CANONICAL | cmd/ailang-worldd/operator_service.go:166 | ` &#124;&#124; creds != c.Credentials &#124;&#124; world == creds` → ` &#124;&#124; world == creds` | TestOperatorServicePair | S | NM | NM |
| INSTALL-PLIST-DIRECTORY | cmd/ailang-worldd/operator_service.go:355 | `runnerPath(p.Home, filepath.Dir(p.Plist), c.UID, true, 0)` → `runnerPath(p.Home, filepath.Dir(p.Plist), c.UID, true, 0755)` | TestOperatorServiceOwnership | S | NM | NM |
| INSTALL-WRITE-LOG-ERROR | cmd/ailang-worldd/operator_service.go:371 | `if err := e.atomic(p.Log, nil, 0600); err != nil {` → `if err := e.atomic(p.Log, nil, 0600); err != nil && false {` | TestOperatorServiceWriteScope | S | NM | NM |
| INSTALL-WRITE-RECORD-ERROR | cmd/ailang-worldd/operator_service.go:378 | `if err := e.atomic(p.Record, b, 0600); err != nil {` → `if err := e.atomic(p.Record, b, 0600); err != nil && false {` | TestOperatorServiceWriteScope | S | NM | NM |
| INSTALL-WRITE-PLIST-ERROR | cmd/ailang-worldd/operator_service.go:381 | `if err := e.atomic(p.Plist, []byte(runnerPlist(p, c, b)), 0600); err != nil {` → `if err := e.atomic(p.Plist, []byte(runnerPlist(p, c, b)), 0600); err != nil && false {` | TestOperatorServiceWriteScope | S | NM | NM |
| INSTALL-STOP-ERROR | cmd/ailang-worldd/operator_service.go:408 | `if err := e.controller.Bootout(e.domain, runnerLabel); err != nil && !errors.Is(err, errServiceNotLoaded) {` → `if err := e.controller.Bootout(e.domain, runnerLabel); err != nil && false {` | TestOperatorServiceWriteScope | S | NM | NM |
| UNINSTALL-REMOVE-ERROR | cmd/ailang-worldd/operator_service.go:414 | `if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {` → `if err := os.Remove(path); err != nil && false {` | TestOperatorServiceWriteScope | S | NM | NM |
| RUNNER-TICK-PRECHECK | cmd/ailang-worldd/operator_runner.go:77 | `if err := e.checkRunnerConfig(c); err != nil {` → `if err := e.checkRunnerConfig(c); err != nil && false {` | TestOperatorServicePair | S | NM | NM |
| RUNNER-STATUS-READ-ERROR | cmd/ailang-worldd/operator_runner.go:111 | `if err != nil { ↵   return errors.New("runner heartbeat unavailable")` → `if err != nil && false { ↵   return errors.New("runner heartbeat unavailable")` | TestOperatorStatusDisconnected | S | NM | NM |
| RUNNER-WRITER-EXISTING | cmd/ailang-worldd/operator_runner.go:151 | `RawQuery: "mode=rw&_pragma=busy_timeout%28250%29"` → `RawQuery: "_pragma=busy_timeout%28250%29"` | TestOperatorStatusDisconnected | S | NM | NM |
| RUNNER-WRITER-BUSY-BOUND | cmd/ailang-worldd/operator_runner.go:151 | `busy_timeout%28250%29` → `busy_timeout%282000%29` | TestOperatorStatusDisconnected | S | NM | NM |

**Named survivors with whole NM.** PAIR-UPGRADE-INSTALLED, PAIR-UPGRADE-SUPERVISOR, PAIR-UPGRADE-LOCK, CONFIG-SIZE, CONFIG-NOFOLLOW, CONFIG-REGULAR, ACCOUNT-LOOKUP, ACCOUNT-SERVICE-HOME, CONFIG-PRESET-EMPTY, CONFIG-PRESET-BOUND, CONFIG-EXEC-BOUND, CONFIG-PATH-BOUND, CONFIG-PRESET-NAME, CONFIG-PRESET-LENGTH, CONFIG-PRESET-RESERVED, CONFIG-LOOPBACK, CONFIG-NUMERIC-LOOPBACK, PATH-ABSOLUTE, PATH-CLEAN, PATH-OWNER, PATH-DIRECTORY, PATH-PRIVATE-MODE, PATH-TOOLS-SEPARATOR, EXEC-NONEMPTY-HASH, PAIR-RECORD-UNINSTALLED, PAIR-RECORD-ADDRESS, PAIR-RECORD-DIGEST, PAIR-RECORD-WORLD, PAIR-RECORD-CREDS, PAIR-HEALTH-ERROR, INSTALL-RECHECK, INSTALL-PAIR, INSTALL-LOADED, INSTALL-IMMUTABLE, INSTALL-SAFE-PATH, INSTALL-ROOT-MODE, INSTALL-SPOOL-CHECK, INSTALL-LOG-CHECK, INSTALL-ROLLBACK-CONFIG, LIFECYCLE-LOCK-NB, LIFECYCLE-PLATFORM, LIFECYCLE-INSTALL-ARGV, LIFECYCLE-CANCEL-ID, LIFECYCLE-ARGV-ARITY, LIFECYCLE-STOP-APPROVAL, LIFECYCLE-STOP-LABEL, RUNNER-CONFIG-ENV, RUNNER-SELF-DIGEST, RUNNER-LOCK-NOFOLLOW, RUNNER-LOCK-OWNER, RUNNER-LOCK-MODE, RUNNER-LOCK-REGULAR, RUNNER-POLL-INTERVAL, RUNNER-GENERATION-ADVANCE, RUNNER-STATUS-SUPERVISOR, RUNNER-MODE-DISPATCH, RUNNER-STATUS-STALE, DOCTOR-COMPANION-DIAGNOSTIC, CONFIG-READ-SOURCE-IO, CONFIG-STORE-IDENTITY-DISTINCT, CONFIG-STORE-WORLD-CANONICAL, CONFIG-STORE-CREDS-CANONICAL, INSTALL-PLIST-DIRECTORY, INSTALL-WRITE-LOG-ERROR, INSTALL-WRITE-RECORD-ERROR, INSTALL-WRITE-PLIST-ERROR, INSTALL-STOP-ERROR, UNINSTALL-REMOVE-ERROR, RUNNER-TICK-PRECHECK, RUNNER-STATUS-READ-ERROR, RUNNER-WRITER-EXISTING, RUNNER-WRITER-BUSY-BOUND. These are open behavioral coverage gaps. The outside-sandbox suite is still required. Do not treat a sandbox bind failure as proof or arm dispatch on this evidence.

**Whole-package survivors.** HB-READ-DEADLINE, HB-WRITE-DEADLINE. Preserve the production guards and add discriminating fixtures before using their results to admit effects.

**HB-SINGLETON-WHERE (E).** The v10 PK CHECK makes singleton=1 mandatory for every valid row and permits at most one.

**CANCEL-OWNER (E).** The state predicate admits requested/confirmed, whose owners must be NULL under the v10 CHECK.

**Positive harness control.** LIVE-COMMENT S/S. Comment-only edits are expected to survive. Excluded from the scoped tally.

**Parked follow-up variants.** Separate counts for S4b and the replay draft. The replay also runs the source package host/store.

| id | Diff anchor | Mutation | Named test | Named | Whole | Verdict |
|---|---|---|---|---|---|---|
| S4B-PAIR-UPGRADE-INSTALLED | cmd/ailang-worldd/service_upgrade.go:206 | `[]string{e.runnerPaths().Record, e.runnerPaths().Plist}` → `[]string{}` | TestOperatorServiceUpgradePairFence | K | NM | NM |
| S4B-PAIR-UPGRADE-SUPERVISOR | cmd/ailang-worldd/service_upgrade.go:211 | `!errors.Is(err, errServiceNotLoaded) { ↵   return refuse("cannot verify` → `err == nil && false { ↵   return refuse("cannot verify` | TestOperatorServiceUpgradePairFence | K | NM | NM |
| S4B-PAIR-UPGRADE-LOCK | cmd/ailang-worldd/service_upgrade.go:201 | `runnerLifecycleLock(filepath.Dir(e.paths.Record))` → `os.Open(filepath.Dir(e.paths.Record))` | TestOperatorServiceUpgradePairFence | K | NM | NM |
| S4B-STARTUP-NEW-GENERATION | cmd/ailang-worldd/operator_runner.go:113 | `h.Owner.Generation > previous` → `h.Owner.Generation >= previous` | TestOperatorServiceStartupTimeout | K | NM | NM |
| S4B-STARTUP-WAIT | cmd/ailang-worldd/operator_service.go:443 | `if !e.waitForRunner(c, previous.Owner.Generation)` → `if !e.waitForRunner(c, previous.Owner.Generation) && false` | TestOperatorServiceStartupTimeout | K | NM | NM |
| S4B-STARTUP-STOP | cmd/ailang-worldd/operator_service.go:444 | `if err := e.controller.Bootout(e.domain, runnerLabel); err != nil && !errors.Is(err, errServiceNotLoaded) {` → `if err := errors.New("bootout omitted"); err != nil && !errors.Is(err, errServiceNotLoaded) {` | TestOperatorServiceStartupTimeout | K | NM | NM |
| TERMINAL-NB | cmd/ailang-worldd/operator_runner.go:31 | `syscall.LOCK_EX&#124;syscall.LOCK_NB` → `syscall.LOCK_EX` | TestOperatorTerminalExclusion | K | NM | NM |
| TERMINAL-DURABLE-OWNER | cmd/ailang-worldd/operator_episode_lock.go:54 | `id, readErr := read()` → `id, readErr := read(); if lockErr == nil { id = "" }` | TestOperatorTerminalExclusion | K | NM | NM |
| TERMINAL-NOLOCK | cmd/ailang-worldd/operator_episode_lock.go:40 | `lock, lockErr := runnerKernelLock(path)` → `var lock *os.File; var lockErr error` | TestOperatorTerminalExclusion | K | NM | NM |
| TERMINAL-ACTION-ID | cmd/ailang-worldd/operator_episode_lock.go:64 | `if verbs.ValidOperatorActionID(id)` → `if verbs.ValidOperatorActionID(id) && false` | TestOperatorTerminalExclusion | K | NM | NM |
| TERMINAL-START-GATE | cmd/ailang-worldd/episode.go:364 | `terminalEpisodeLock(req.credPath, episode, nil)` → `terminalEpisodeLock(req.credPath, episode+"-other", nil)` | TestOperatorTerminalExclusionStart | K | NM | NM |
| REPLAY-IDENTICAL | host/store/operator_cas.go:59 | `if !bytes.Equal(row.Artifact, artifact) {` → `if !bytes.Equal(row.Artifact, artifact) && false {` | TestOperatorConfirmReplay | K | K | K |

REPLAY-IDENTICAL source-only whole package: S. The behavioral test is in host/authority.

## 5. Files

| Stage | Paths |
|---|---|
| M1 | cmd/ailang-worldd/{service_upgrade.go,version.go,service_test.go,operator_upgrade_test.go}. |
| M2 | cmd/ailang-worldd/{main.go,keyguard.go,doctor_service.go,service.go,service_upgrade.go,operator_service.go,operator_runner.go,operator_service_test.go,service_test.go,cli_credential_writers_test.go,help_test.go,keyguard_test.go,quickstart_service_test.go}; host/store/operator_runner.go; docs/QUICKSTART.md. |
| M3 | cmd/ailang-worldd/operator_service_test.go; host/store/operator_runner_test.go. |

Pre-existing tests edited: **service_test.go** supplies schema capabilities, canonical fake HOME and independent fake labels; **cli_credential_writers_test.go** adds exactly runnerWrite; **help_test.go** enumerates the new family under synthetic HOME; **keyguard_test.go** validates the expanded inventory and key guard under synthetic HOME; **quickstart_service_test.go** adds a bad runner flag control and synthetic HOME. Zero existing tests deleted. No inventory edits for dependency allowlists, production context roots/Go surface, 29 public switch arms, daemon writers or signer/export guard. §6 monitors the extended LOCK_NB census.

## 6. Sprint gates

Same bytes as `$HOME/.ailang/state/world-iter429/gates.sh`. The first and last legs record read-only observations and the sentinel records every failure. Uninformative sandbox legs have a separate summary field and exit 77. Invoke with committed tree, `$HOME/.ailang/state/world-iter429/gates`, mode `final`, and final SHA. For scratch supply `GATES_CHANGED_FILE` and `GATES_GIT_REF`. The controller runs unmodified full package/race commands outside the sandbox and honors GATES_SKIP_VERIFY_GO=1.

````bash gates.sh
#!/bin/bash
# Usage: gates.sh <tree> <logdir> <base|final|scratch> [sha]
set -u
set -o pipefail
T=$1; L=$2; MODE=$3
B=f074dd49e8fe1f281792de4fb475e20b53cd1d92
REF=${GATES_GIT_REF:-$T}
SHA=${4:-$(git -C "$REF" rev-parse HEAD)}
mkdir -p "$L"
cd "$T" || exit 2
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export AILANG_STATE_DIR="$L/ailang-state"
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
inv cmd/ailang-worldd TestEveryVerbHasHelp TestVerbTableMatchesTheCLISwitch TestCLICredentialWriterSites TestOperatorFlagParity TestQuickstartServiceSectionMatchesTheCLIHelp TestHelpIsNeverBlockedByTheRegistryKeyGuard TestHomeServiceSeamOnly TestHomeServiceInstallAndPlist TestHomeServiceUpgradeRollback
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
  allowed='^(cmd/ailang-worldd/(main|keyguard|doctor_service|service|service_upgrade|version|cli_credential_writers_test|help_test|keyguard_test|quickstart_service_test|service_test|operator_runner|operator_service|operator_service_test|operator_upgrade_test)\.go|host/store/operator_runner(_test)?\.go|docs/QUICKSTART\.md|design_docs/planned/(w-world-address-bar-verbs-v4s4-sprint-plan\.md|sprint_w-world-address-bar-verbs-v4s4\.json))$'
  printf '%s\n' cmd/ailang-worldd/service.go host/store/operator_runner.go docs/QUICKSTART.md | grep -E "$allowed"
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
    external=$(printf '%s\n' "$callers" | grep -v '^cmd/ailang-worldd/operator_runner.go:' || test "$?" -eq 1)
    test -z "$external"
  fi
  grep -n 'syscall.Flock.*LOCK_NB' cmd/ailang-worldd/service.go host/store/writer_lock_unix.go
  if [ "$MODE" != base ]; then
    grep -n 'syscall.Flock.*LOCK_NB' cmd/ailang-worldd/operator_runner.go cmd/ailang-worldd/operator_service.go
  fi
}
leg callers 'caller_fence'
if [ "${GATES_SKIP_VERIFY_GO:-0}" = 1 ]; then
  leg verify_go 'echo "SKIPPED (GATES_SKIP_VERIFY_GO=1)"'
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

## 7. Executor directive notes

- Run M1 → M2 → M3. /usr/bin/patch -p1 -E is the only worktree writer. No Git writes, hand edits, gofmt writes, golden updates, Chrome, nohup or background test legs. Parked fences are outside the execution loop.
- With `W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter429-row243-v4s4`, `T=$HOME/.ailang/state/world-iter429/exec-patches` and `PLAN=$W/design_docs/planned/w-world-address-bar-verbs-v4s4-sprint-plan.md`, create T and extract the accept and snapshot fences with §2 awk, substituting the label `bash accept.sh` or `bash snapshot.sh`. Record the tests under pipefail to `$HOME/.ailang/state/world-iter429/exec-logs/accept_$M.txt`.
- Commit-time commands m1/m2/m3 are `AILANG_BIN=$HOME/.pinned-ailang/ailang bash "$T/accept.sh" "$W" M1 fast`, `AILANG_BIN=$HOME/.pinned-ailang/ailang bash "$T/accept.sh" "$W" M2 fast` and `AILANG_BIN=$HOME/.pinned-ailang/ailang bash "$T/accept.sh" "$W" M3 fast`. The script also exports the isolated AILANG_STATE_DIR and fixture flags. Controller expands paths before `commit_ms.sh <M> <message-file> <test-cmd…>`. Record both plan files first because the commit script refuses untracked design docs outside snapshots.
- Snapshot with `bash "$T/snapshot.sh" "$W" "$M"`. Stop on existing snapshot, manifest mismatch, unexpected path, skipped expected test, actual gofmt/build/vet/test failure or cumulative numstat excess. SNAP_STATUS_LIST is scratch-only. Check filtered and unfiltered numstat. 
- Socket and launchctl denials are **UNINFORMATIVE UNDER SANDBOX**. Preserve logs and run those exact gates outside. Re-run a documented load-sensitive test once alone. Includes broker exec tests/InFlight timeouts, TestSeedRestartIgnoresASurvivingOldCopy, TestMissionConfigDefaultFlowIntegration, TestSeToolsRunBenchmarkTasksEndToEnd, TestHostsLaunchersExecHost, TestBackoff and host/broker, host/daemon and host/verifygate race runs. Sequence full suite and mutations.

````bash snapshot.sh
#!/bin/bash
set -euo pipefail
W=$1
M=$2
B=f074dd49e8fe1f281792de4fb475e20b53cd1d92
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter429/snap}/$M
mkdir "$D"
cd "$W"
if [ -n "${SNAP_STATUS_LIST:-}" ]; then
  cp "$SNAP_STATUS_LIST" "$D/STATUS"
else
  {
    git diff --name-status --no-renames "$B" -- cmd host docs app
    git ls-files --others --exclude-standard -- cmd host docs app | sed 's/^/A	/'
  } | LC_ALL=C sort -u -t "$(printf '\t')" -k2 > "$D/STATUS"
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
  done | LC_ALL=C sort -z | xargs -0 shasum -a 256 > MANIFEST.sha256
  shasum -a 256 -c MANIFEST.sha256
)
wc -l "$D"/{STATUS,DELETED,MANIFEST.sha256}
````

M1 full-tree diff empty; snapshot 4/0/4
M2 full-tree diff empty; snapshot 17/0/17
M3 full-tree diff empty; snapshot 18/0/18
S4b full-tree diff empty against its independent follow-up stage
S3b full-tree diff empty against its independent follow-up stage
S4c-replay full-tree diff empty against its independent follow-up stage
PASS full-tree plan extraction round trip, including docs; artifacts /Users/voightkampff/.ailang/state/world-iter429/plan-scratch/roundtrip.dGdmDa

## 8. Risks

- **Outside-worktree writes.** installRunner stages config, companion plist and log via the existing atomic writer (temp under injected service bin), creates private spool and removes the uncommitted config if plist placement fails. changeRunner removes only config/plist after supervisor stop. TestOperatorServiceWriteScope checks fake home paths, private modes and exact base-file/DB hashes; TestOperatorServiceEnvironment reads real renderer output, config and log after placing a fake ambient token.
- **SQLite and locks.** runnerWrite can create/test the credentials writer lock and SQLite sidecars and write heartbeat/cancel rows. The internal mode opens the singleton lock file next to the canonical creds path. Store tests and TestOperatorStatusDisconnected use disposable migrated credentials; TestOperatorRunnerKernelLock uses a temp path. The lifecycle directory lock writes no marker. Bound open busy timeout and close writer before sleeps.
- **Upgrade.** Existing upgrade path writes .new/.prev/.failed/log/binary/record. Fake restart tests exercise the same writer seam and preserve credentials bytes. Unknown rollback compatibility may leave the service record on the old digest while the candidate remains installed. Inspect and recover with a verified compatible binary under attendance. Paired restart/drain stays gated before effects.
- **Uncertain owner.** Stopping the runner or losing heartbeat does not remove persisted ownership. S4 has no producer of new uncertain work. S4b terminal lockout preserves uncertain evidence even after kernel release. An uncertain outcome may block Mark's fallback; an attended resolver stays open. Never delete lock or receipt files as a workaround.
- **Observation.** §6 observes LaunchAgents, World bin, top-level JSON and launchctl count. Volatile service logs and stores are excluded.

## 9. Acceptance criteria

1. M1–M3 apply cumulatively, reproduce compiled scratch trees, and snapshots verify manifests. Base-relative cap uses the exact §1.5 filter.
2. Typed immutable setup, exact build/store pair, isolated label, private modes, explicit HOME/PATH and no token persistence are exercised with fake controller/terminal/transport and real disposable stores.
3. Incompatible or unverified previous binaries are never installed by rollback. Compatible rollback still passes. Companion install fences unpaired home upgrade.
4. Kernel singleton, time operators, owner/generation/protocol/build predicates, catalogue immutability and per-verb readiness have named pins. Durable last-seen status identifies disconnect.
5. Attended unstarted cancel advances once; in-flight cancellation preserves the action and reservation. Full child stop/reconciliation is a later gate.
6. Record §4 named and whole verdicts, including survivors or sandbox-unmeasured rows. Behavioral claims require a named assertion. Apply only §5 changes; exact prior test edits are listed there and none deleted.
7. No new daemon route, .ail, app.js/surfaceCSP edit, served change or reachable effect dispatch. Appendix C disposes each inherited S4 gate and 4(a)–(c).

## Appendix C — after this sprint

- **S4b, 366 additional lines.** Compiled terminal exclusion on M3: real plain/preset/resume start, launch and F1 finish hooks, nested start-to-launch lock transfer, canonical per-store/episode lock path, direct LOCK_NB, exact active id and durable owner check even after unlock. Tests drive real plain start, launch and F1 finish, zero downstream effects, and unrelated-episode acquisition. The same S4b patch adds a bounded fresh-generation startup gate, deterministic fake timeout and checks for installed, registered and locked upgrade pairing fences. Add preset/resume behavioral fixtures and a paired process's production episode-lock consumer before dispatch. Exercise the first-managed-store/drain transition before arming S5/S6. The 5 s bound applies to the refusal lookup. The lock attempt itself is non-waiting. Never use stale heartbeat to release ownership.
- **Open S4 rows.** TestOperatorRecoverUntouched needs independently frozen intent and positive no-reservation evidence; C5 finds no owner handoff logic. TestOperatorRecoverReserved needs S3b plus S6's process birth/run identity; C7 finds only the launch PID. TestOperatorRunnerRechecks daemon half is pinned here through TestOperatorServicePair; S5 couples that to the signed action at each new effect. TestOperatorCancelInflight here proves refusal and evidence retention; S6 adds bounded stop/reconcile. TestOperatorStatusDisconnected here shows last heartbeat and explicit disconnect; S7 supplies last action status on failed reads and enforces zero request inserts when heartbeat or catalogue admission fails. TestOperatorServiceStartup includes a v9 skew fixture. S4b adds a bounded bootstrap/next-generation gate and a fake timeout with old liveness preserved. TestOperatorServiceUpgrade still needs paired restart/active child drain and redacted rotation.
- **4(a), S3b.** 132 additional lines, independently applied and compiled. Reserve/observe/restart slice does not fit 0 lines of headroom, before fd tests, all receipt predicates, and current-owner recovery are added. TestOperatorStepRestart proves reservation survives reopen only. Route all of those gates before dispatch.
- **4(b), rollback.** M1 implements the guard and fake v10 fixture. No real store or deployment.
- **4(c), judge-428.** Carry J-REQ-PENDING-CONF, J-DDL-STEPS-PK, J-REQ-ACTIVE-ID, J-CONF-ARTSIZE, J-REQ-DEDUP-ORDER, J-FLOOR-PRESET-REVERT, J-FLOOR-MINT and J-FK-INMEM from C11/C12. Keep cmd and store v9 write-floor cases distinct. Same-artifact replay and wrong-byte refusal are compiled in parked S4c-replay (13 lines). Before S7, give its consumer an explicit replay result so a seen link grants no new claim; there is no production caller yet. go.mod's x/sys direct-test-use correction is one replacement (two changed lines under widened pathspec), also deferred.
- **S5** inherits canonical pair/status, bounded writes, action schema, and S4b/S3b gates. Owns pinned-start read seam, minted-start proof, and §5.3 StartOnce/Changed/Collisions, WorktreeWindow, MintBinding/Rollback/ReceiptWindow and FilesReconcile. **S6** owns LaunchSession, SpawnRefused, LaunchCrashWindow, LaunchExit, StartBrowser and child provenance. **S7** owns §5.1 Misdirected, CrossOrigin, FetchSite, MediaType, BodyBound, UnsupportedMethod, RequestLimits, LiteralTask at route, BarStartReview, BarFinishReview and OtherVerbParity, plus the stable action projection. Invoke the existing judge using trusted read-only state, consume confirm CAS, and keep replay distinct. **S8** owns PanelNetwork, PanelLatest, LinkAllowlist, WireTypes and fake Chrome. **S9–S11** own §5.4 finish with E2 readiness, linked V3 confirmation, exact remote reconciliation and scoped record. **S12** owns HerdrUnavailable/Handoff/Window under 119 A and E4, including the plugin hard-code.
- **Other carried items.** Judge-408 note 4 field redundancy and judge-407 note 3 numeric string flags/usage belong with S7; judge-407 note 2/S1 F8 null arrays with S8. Two-store V3 orphan window and judge-402 v9 prune/DDL remain separate. Uncertain resolution needs an attended resolver. A headless child outliving the runner requires independent supervisor proof. Update writer-site and non-waiting inventories together.
- **Controller.** Keep M1/M2/M3. Add docs to ctl_roundtrip.sh or compare full trees. Run §6 outside the sandbox and every NM whole-package mutant. Route S4b and receipt gates before S5/S6 effects, replay guard before S7. The reviewed design and ledger stay closed.

## Appendix P — compiled patches

### Patch M1

````diff M1
diff --git a/cmd/ailang-worldd/operator_upgrade_test.go b/cmd/ailang-worldd/operator_upgrade_test.go
new file mode 100644
index 0000000..359daaa
--- /dev/null
+++ b/cmd/ailang-worldd/operator_upgrade_test.go
@@ -0,0 +1,83 @@
+package main
+
+import (
+	"encoding/json"
+	"fmt"
+	"os"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func TestOperatorServiceUpgrade(t *testing.T) {
+	for _, tc := range []struct {
+		name         string
+		before, next int
+		tamper       string
+	}{
+		{"old reader", 9, 10, ""},
+		{"missing previous schema", 0, 10, ""},
+		{"missing candidate schema", 10, 0, ""},
+		{"substituted previous binary", 10, 10, "different bytes"},
+		{"linked previous binary", 10, 10, "link"},
+		{"missing previous binary", 10, 10, "missing"},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			f, candidate, _ := upgradeFixture(t)
+			original := serviceFileHash(f.env.paths.Binary)
+			r, err := readServiceRecord(f.env.paths)
+			if err != nil {
+				t.Fatal(err)
+			}
+			if _, err := store.UpgradeCredentials(r.Credentials); err != nil {
+				t.Fatal(err)
+			}
+			db, record := serviceFileHash(r.Credentials), serviceFileHash(f.env.paths.Record)
+			probe := f.env.probe
+			f.env.probe = func(path, home string) ([]byte, error) {
+				b, err := probe(path, home)
+				if err != nil {
+					return nil, err
+				}
+				var v serviceVersion
+				if err := json.Unmarshal(b, &v); err != nil {
+					return nil, err
+				}
+				v.CredentialsSchema = tc.before
+				if path == candidate {
+					v.CredentialsSchema = tc.next
+				}
+				return json.Marshal(v)
+			}
+			f.ctl.onKick = func() error {
+				if tc.tamper != "" {
+					prev := f.env.paths.Binary + ".prev"
+					if err := os.Remove(prev); err != nil {
+						t.Fatal(err)
+					}
+					switch tc.tamper {
+					case "link":
+						if err := os.Symlink(f.source, prev); err != nil {
+							t.Fatal(err)
+						}
+					case "different bytes":
+						if err := os.WriteFile(prev, []byte("substitution"), 0700); err != nil {
+							t.Fatal(err)
+						}
+					}
+				}
+				return fmt.Errorf("fixture restart failed")
+			}
+			code, msg := f.run("upgrade", "--binary", candidate)
+			if code != 3 || !strings.Contains(msg, "automatic rollback refused") ||
+				!strings.Contains(msg, "verified schema-compatible binary") || f.ctl.kicks != 1 || f.ctl.stops != 0 ||
+				serviceFileHash(f.env.paths.Binary) != serviceFileHash(candidate) ||
+				serviceFileHash(f.env.paths.Record) != record || serviceFileHash(r.Credentials) != db ||
+				serviceFileHash(f.env.paths.Binary+".failed") != "" ||
+				(tc.tamper == "" && serviceFileHash(f.env.paths.Binary+".prev") != original) {
+				t.Fatalf("upgrade %d, kicks %d, stops %d: %s", code, f.ctl.kicks, f.ctl.stops, msg)
+			}
+		})
+	}
+}
diff --git a/cmd/ailang-worldd/service_test.go b/cmd/ailang-worldd/service_test.go
index fc34568..e2d8429 100644
--- a/cmd/ailang-worldd/service_test.go
+++ b/cmd/ailang-worldd/service_test.go
@@ -80,7 +80,7 @@ func newServiceFixture(t *testing.T) *serviceFixture {
 	f.ctl = &fakeServiceController{binary: f.env.paths.Binary}
 	f.env.controller = f.ctl
 	f.env.probe = func(binary, home string) ([]byte, error) {
-		return json.Marshal(serviceVersion{Build: daemon.BuildIdentity{ExecutableSHA256: serviceFileHash(binary)}, Platform: runtime.GOOS + "/" + runtime.GOARCH, ToolRelease: toolRelease(), InterpreterRelease: interpreterRelease})
+		return json.Marshal(serviceVersion{Build: daemon.BuildIdentity{ExecutableSHA256: serviceFileHash(binary)}, Platform: runtime.GOOS + "/" + runtime.GOARCH, ToolRelease: toolRelease(), InterpreterRelease: interpreterRelease, CredentialsSchema: store.CredentialsSchemaVersion})
 	}
 	f.env.observe = func(p string) { f.writes = append(f.writes, p) }
 	f.env.client = &http.Client{Transport: serviceTransport(func(r *http.Request) (*http.Response, error) {
@@ -399,7 +399,8 @@ func TestHomeVersionJSON(t *testing.T) {
 			t.Fatalf("%s = %s, want %s", key, got, want)
 		}
 	}
-	if v.Build.ExecutableSHA256 == "" || v.Build.GoVersion == "" {
+	expectedSchema, _ := json.Marshal(store.CredentialsSchemaVersion)
+	if v.Build.ExecutableSHA256 == "" || v.Build.GoVersion == "" || !bytes.Equal(raw["credentials_schema"], expectedSchema) {
 		t.Fatal(b.String())
 	}
 }
diff --git a/cmd/ailang-worldd/service_upgrade.go b/cmd/ailang-worldd/service_upgrade.go
index d570803..4d44ce0 100644
--- a/cmd/ailang-worldd/service_upgrade.go
+++ b/cmd/ailang-worldd/service_upgrade.go
@@ -24,6 +24,7 @@ type serviceUpgrade struct {
 }
 
 type serviceVersion struct {
+	CredentialsSchema  int                  `json:"credentials_schema,omitempty"`
 	Build              daemon.BuildIdentity `json:"build"`
 	Platform           string               `json:"platform"`
 	ToolRelease        string               `json:"tool_release"`
@@ -237,12 +238,21 @@ func (e serviceEnv) upgrade(candidate string, pause bool, out, errOut io.Writer)
 		}
 		// A healthy but unrecorded swap is not a successful upgrade.
 	}
-	return e.rollbackUpgrade(r, from, to, candidate, errOut)
+	// A candidate can migrate before reporting health, or while health is
+	// unavailable. The highest schema it could open is the safe comparison.
+	// Older version outputs omitted this field and cannot prove compatibility.
+	return e.rollbackUpgrade(r, from, to, candidate, installed.CredentialsSchema, v.CredentialsSchema, errOut)
 }
 
-func (e serviceEnv) rollbackUpgrade(r serviceRecord, from, to, candidate string, w io.Writer) int {
+func (e serviceEnv) rollbackUpgrade(r serviceRecord, from, to, candidate string, previousSchema, candidateSchema int, w io.Writer) int {
 	tail := serviceTail(e.paths.Log, 20)
 	failed, previous := e.paths.Binary+".failed", e.paths.Binary+".prev"
+	previousInfo, previousErr := os.Lstat(previous)
+	if candidateSchema <= 0 || previousSchema < candidateSchema || previousErr != nil ||
+		!previousInfo.Mode().IsRegular() || serviceFileHash(previous) != from {
+		return serviceMessage(w, "automatic rollback refused", "the previous binary cannot be verified compatible with the possibly migrated credentials store; both builds were retained",
+			"inspect service status and recover at the terminal with a verified schema-compatible binary", 3)
+	}
 	restored := false
 	e.note(failed)
 	err := os.Remove(failed)
diff --git a/cmd/ailang-worldd/version.go b/cmd/ailang-worldd/version.go
index c0b0fcb..09b24b3 100644
--- a/cmd/ailang-worldd/version.go
+++ b/cmd/ailang-worldd/version.go
@@ -8,6 +8,7 @@ import (
 	"runtime"
 
 	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/store"
 )
 
 const versionHelp = "usage: ailang-worldd version [--json]\n"
@@ -23,12 +24,11 @@ func runVersion(args []string, stdout, stderr io.Writer) int {
 		fmt.Fprint(stderr, versionHelp)
 		return exitUsage
 	}
-	v := struct {
-		Build              daemon.BuildIdentity `json:"build"`
-		Platform           string               `json:"platform"`
-		ToolRelease        string               `json:"tool_release"`
-		InterpreterRelease string               `json:"interpreter_release"`
-	}{daemon.CurrentBuild(), runtime.GOOS + "/" + runtime.GOARCH, toolRelease(), interpreterRelease}
+	v := serviceVersion{
+		Build: daemon.CurrentBuild(), Platform: runtime.GOOS + "/" + runtime.GOARCH,
+		ToolRelease: toolRelease(), InterpreterRelease: interpreterRelease,
+		CredentialsSchema: store.CredentialsSchemaVersion,
+	}
 	if *asJSON {
 		if err := json.NewEncoder(stdout).Encode(v); err != nil {
 			return exitUsage
````

### Patch M2

````diff M2
diff --git a/cmd/ailang-worldd/cli_credential_writers_test.go b/cmd/ailang-worldd/cli_credential_writers_test.go
index 42bbe4c..95f52e2 100644
--- a/cmd/ailang-worldd/cli_credential_writers_test.go
+++ b/cmd/ailang-worldd/cli_credential_writers_test.go
@@ -49,7 +49,7 @@ func TestCLICredentialWriterSites(t *testing.T) {
 			})
 		}
 	}
-	want := map[string]string{"openCredentialsForSession": "OpenCredentials", "launchEpisode": "OpenCredentials", "finishRevoke": "OpenCredentials", "putFinishVerbRequest": "OpenCredentials"}
+	want := map[string]string{"runnerWrite": "OpenCredentials", "openCredentialsForSession": "OpenCredentials", "launchEpisode": "OpenCredentials", "finishRevoke": "OpenCredentials", "putFinishVerbRequest": "OpenCredentials"}
 	if !reflect.DeepEqual(sites, want) || readOnly < 20 {
 		t.Fatalf("writer sites %v (read-only opens %d), want exactly %v and at least 20 read-only opens", sites, readOnly, want)
 	}
diff --git a/cmd/ailang-worldd/doctor_service.go b/cmd/ailang-worldd/doctor_service.go
index 107a5de..b3dc51a 100644
--- a/cmd/ailang-worldd/doctor_service.go
+++ b/cmd/ailang-worldd/doctor_service.go
@@ -1,6 +1,8 @@
 package main
 
 import (
+	"bytes"
+	"errors"
 	"fmt"
 	"os"
 	"regexp"
@@ -43,6 +45,21 @@ func doctorService(report *doctorReport, addr string) {
 		report.add(markWarn, "service", fmt.Sprintf("this ailang-worldd (%s, sha %s) is not the installed one (%s)", path, pointerSHA7(cli), pointerSHA7(installed)),
 			fmt.Sprintf("run %s, or install this one: %s service upgrade", r.Binary, path))
 	}
+	if serviceAddressMatches(addr, r) {
+		var status bytes.Buffer
+		switch _, err := os.Lstat(e.runnerPaths().Record); {
+		case errors.Is(err, os.ErrNotExist):
+			report.add(markWarn, "operator runner", "companion absent", "ailang-worldd service runner install --help")
+		case err != nil:
+			report.add(markWarn, "operator runner", "companion configuration inaccessible", "ailang-worldd service runner status")
+		default:
+			if err := e.runnerStatus(&status); err != nil {
+				report.add(markWarn, "operator runner", status.String()+runnerDiagnostic(err), "ailang-worldd service runner status")
+			} else {
+				report.add(markOK, "operator runner", status.String(), "")
+			}
+		}
+	}
 	if info, err := os.Stat(e.paths.Log); err == nil && info.Size() > 100<<20 {
 		report.add(markWarn, "service", "the service log is over 100 MiB: "+e.paths.Log, "rotate the log while the service is stopped")
 	}
diff --git a/cmd/ailang-worldd/help_test.go b/cmd/ailang-worldd/help_test.go
index e64d356..4a9eb60 100644
--- a/cmd/ailang-worldd/help_test.go
+++ b/cmd/ailang-worldd/help_test.go
@@ -12,12 +12,14 @@ var everyVerb = [][]string{
 	{"log", "get"}, {"log", "range"}, {"log", "tail"}, {"registry", "get"}, {"commit"},
 	{"tools", "list"}, {"call"}, {"why"}, {"provenance"}, {"setup"}, {"doctor"}, {"hosts", "init"},
 	{"session"}, {"session", "new"}, {"session", "list"}, {"session", "revoke"}, {"session", "mint"},
-	{"approve"}, {"approver"}, {"approver", "list"}, {"approver", "enroll"}, {"approver", "revoke"},
+	{"approve"}, {"service", "runner"}, {"service", "runner", "install"}, {"service", "runner", "start"},
+	{"service", "runner", "stop"}, {"service", "runner", "status"}, {"service", "runner", "cancel"}, {"service", "runner", "uninstall"}, {"approver"}, {"approver", "list"}, {"approver", "enroll"}, {"approver", "revoke"},
 }
 
 // AC5.7: `<verb> --help` and `help <verb>` exit 0 with that verb's help on
 // stdout for every verb, and the top-level usage is complete.
 func TestEveryVerbHasHelp(t *testing.T) {
+	t.Setenv("HOME", t.TempDir())
 	for _, verb := range everyVerb {
 		name := strings.Join(verb, " ")
 		for _, argv := range [][]string{append(append([]string{}, verb...), "--help"), append([]string{"help"}, verb...)} {
diff --git a/cmd/ailang-worldd/keyguard.go b/cmd/ailang-worldd/keyguard.go
index 9a8f3d8..94788c0 100644
--- a/cmd/ailang-worldd/keyguard.go
+++ b/cmd/ailang-worldd/keyguard.go
@@ -70,6 +70,11 @@ func helpOnly(args []string) (string, bool) {
 // byte with what the verb itself prints, for every verb line.
 func ownVerbHelp(rest []string) (string, bool) {
 	switch rest[0] {
+	case "service":
+		if len(rest) > 1 && rest[1] == "runner" {
+			return runnerHelp, true
+		}
+		return serviceHelp, true
 	case "tools":
 		return toolsHelp, true
 	case "call":
diff --git a/cmd/ailang-worldd/keyguard_test.go b/cmd/ailang-worldd/keyguard_test.go
index 66e3393..76a85f0 100644
--- a/cmd/ailang-worldd/keyguard_test.go
+++ b/cmd/ailang-worldd/keyguard_test.go
@@ -30,6 +30,7 @@ func helpForms() [][]string {
 // without the key (the negative control is the same argv, key absent), and
 // neither stream carries the value.
 func TestHelpIsNeverBlockedByTheRegistryKeyGuard(t *testing.T) {
+	t.Setenv("HOME", t.TempDir())
 	for _, argv := range helpForms() {
 		name := strings.Join(argv, " ")
 		t.Setenv(broker.RegistryCredentialVariable, "")
diff --git a/cmd/ailang-worldd/main.go b/cmd/ailang-worldd/main.go
index da4b5c7..f215c5a 100644
--- a/cmd/ailang-worldd/main.go
+++ b/cmd/ailang-worldd/main.go
@@ -83,6 +83,7 @@ Usage:
   ailang-worldd hosts init [--out-dir <dir>] [--addr <url>] [--replace]
   ailang-worldd version [--json]
   ailang-worldd service install|status|upgrade|stop|start|uninstall [flags; upgrade: --binary FILE --pause]
+  ailang-worldd service runner install|start|stop|status|cancel|uninstall
   ailang-worldd serve --db <path> [--credentials-db <path>] [--bind host:port] [--ailang-bin <path>]
                       [--workbench-still]
                       [--workspace-root <dir> --tool-ailang-bin <path>]
@@ -378,6 +379,12 @@ func run(args []string, stdout, stderr io.Writer) int {
 		return exitOK
 	}
 
+	if rest[0] == "operator-runner" {
+		if addrGiven {
+			return exitUsage
+		}
+		return runOperatorRunner(rest[1:], stderr, realServiceEnv())
+	}
 	switch verb := rest[0]; verb {
 	case "version":
 		return runVersion(rest[1:], stdout, stderr)
diff --git a/cmd/ailang-worldd/operator_runner.go b/cmd/ailang-worldd/operator_runner.go
new file mode 100644
index 0000000..75815a0
--- /dev/null
+++ b/cmd/ailang-worldd/operator_runner.go
@@ -0,0 +1,153 @@
+package main
+
+import (
+	"crypto/rand"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"io"
+	"net/url"
+	"os"
+	"os/signal"
+	"syscall"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+const runnerPoll = 2 * time.Second
+
+func runnerKernelLock(path string) (*os.File, error) {
+	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0600)
+	if err != nil {
+		return nil, err
+	}
+	st, err := f.Stat()
+	if err != nil || !st.Mode().IsRegular() || st.Mode().Perm() != 0600 || int(st.Sys().(*syscall.Stat_t).Uid) != os.Getuid() {
+		f.Close()
+		return nil, errors.New("unsafe runner lock file")
+	}
+	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
+		f.Close()
+		return nil, err
+	}
+	return f, nil
+}
+
+func runOperatorRunner(args []string, errOut io.Writer, e serviceEnv) int {
+	if len(args) != 0 {
+		return exitUsage
+	}
+	c, b, err := e.installedRunner()
+	if err != nil || os.Getenv("WORLD_OPERATOR_CONFIG_SHA256") != serviceHash(b) ||
+		daemon.CurrentBuild().ExecutableSHA256 != c.Executables[e.paths.Binary] {
+		return serviceMessage(errOut, "runner startup refused", "installation or executable identity mismatch", "ailang-worldd service runner status", 2)
+	}
+	if err := e.checkRunnerPair(c); err != nil {
+		return serviceMessage(errOut, "runner startup refused", err.Error(), "ailang-worldd service runner status", 2)
+	}
+	lock, err := runnerKernelLock(c.Credentials + ".operator-runner.lock")
+	if err != nil {
+		return serviceMessage(errOut, "runner startup refused", "cannot acquire the kernel runner lock", "ailang-worldd service runner status", 2)
+	}
+	defer lock.Close()
+	h := store.OperatorHeartbeat{
+		Owner: store.OperatorOwner{Nonce: rand.Text()}, Protocol: store.OperatorRunnerProtocol,
+		Build: c.Executables[e.paths.Binary],
+	}
+	stops := make(chan os.Signal, 1)
+	signal.Notify(stops, syscall.SIGTERM, os.Interrupt)
+	defer signal.Stop(stops)
+	ticker := time.NewTicker(runnerPoll)
+	defer ticker.Stop()
+	for {
+		if err := e.runnerTick(c, &h); err != nil {
+			return serviceMessage(errOut, "runner halted", "pair or heartbeat unavailable", "ailang-worldd service runner status", 2)
+		}
+		select {
+		case <-ticker.C:
+		case <-stops:
+			return 0
+		}
+	}
+}
+
+func (e serviceEnv) runnerTick(c runnerConfig, h *store.OperatorHeartbeat) error {
+	if err := e.checkRunnerConfig(c); err != nil {
+		return err
+	}
+	if err := e.checkRunnerPair(c); err != nil {
+		return err
+	}
+	h.SeenAt = e.terminal.now()
+	h.Readiness = "start"
+	// S4 projects the configured names. No effect dispatch or GitHub probe.
+	h.Catalogue, _ = json.Marshal(c.PresetNames)
+	return runnerWrite(c.Credentials, func(db *store.Credentials) error {
+		ctx, cancel := budgetContext(time.Second)
+		defer cancel()
+		generation, err := db.PutOperatorHeartbeat(ctx, *h)
+		if err == nil {
+			h.Owner.Generation = generation
+		}
+		return err
+	})
+}
+
+func (e serviceEnv) runnerStatus(w io.Writer) error {
+	c, _, err := e.installedRunner()
+	if err != nil {
+		return err
+	}
+	db, err := store.OpenCredentialsReadOnly(runnerDB(c.Credentials))
+	if err != nil {
+		return errors.New("credentials disconnected")
+	}
+	defer db.Close()
+	ctx, cancel := budgetContext(time.Second)
+	defer cancel()
+	h, err := db.OperatorHeartbeat(ctx)
+	if err != nil {
+		return errors.New("runner heartbeat unavailable")
+	}
+	fmt.Fprintf(w, "runner last seen %d, generation %d\n", h.SeenAt, h.Owner.Generation)
+	if err := e.checkRunnerPair(c); err != nil {
+		return fmt.Errorf("daemon disconnected: %w", err)
+	}
+	if _, err := e.controller.Print(e.domain, runnerLabel); err != nil {
+		return errors.New("runner disconnected")
+	}
+	build, now := c.Executables[e.paths.Binary], e.terminal.now()
+	if !h.Available(build, "start", now) {
+		return errors.New("runner heartbeat stale or incompatible")
+	}
+	fmt.Fprintln(w, "start: ready")
+	if !h.Available(build, "finish", now) {
+		fmt.Fprintln(w, "finish blocked: github-auth; repair PATH, GitHub login or Git credential helper in runner setup")
+	}
+	return nil
+}
+
+func runnerWrite(path string, write func(*store.Credentials) error) error {
+	db, err := store.OpenCredentials(runnerDB(path))
+	if err != nil {
+		return err
+	}
+	defer db.Close()
+	return write(db)
+}
+
+func runnerCancel(c runnerConfig, id string, now int64) error {
+	return runnerWrite(c.Credentials, func(db *store.Credentials) error {
+		ctx, cancel := budgetContext(time.Second)
+		defer cancel()
+		return db.CancelUnstartedOperatorAction(ctx, id, now)
+	})
+}
+
+// Bound schema reads on open and require an existing database for writes.
+func runnerDB(path string) string {
+	u := url.URL{Scheme: "file", Path: path, RawQuery: "mode=rw&_pragma=busy_timeout%28250%29"}
+	return u.String()
+}
diff --git a/cmd/ailang-worldd/operator_service.go b/cmd/ailang-worldd/operator_service.go
new file mode 100644
index 0000000..608670e
--- /dev/null
+++ b/cmd/ailang-worldd/operator_service.go
@@ -0,0 +1,451 @@
+package main
+
+import (
+	"bytes"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"io"
+	"net"
+	"os"
+	"os/user"
+	"path/filepath"
+	"strconv"
+	"strings"
+	"syscall"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+const runnerLabel = "com.sunholo.ailang-world.operator-runner"
+const runnerSystemPath = "/usr/bin:/bin:/usr/sbin:/sbin"
+const runnerHelp = `usage: ailang-worldd service runner install --config FILE | start | stop | status | cancel ACTION-ID | uninstall
+Install, stop, cancel and uninstall require typed y. Start uses the approved configuration.
+Config JSON keys: home, uid, address, world, credentials, executables (absolute path to sha256 digest), preset_names, tool_dirs, backend ("headless").
+The World service must already be running. Use absolute, canonical paths and include its installed binary in executables.
+`
+
+type runnerConfig struct {
+	Home        string            `json:"home"`
+	UID         int               `json:"uid"`
+	Address     string            `json:"address"`
+	World       string            `json:"world"`
+	Credentials string            `json:"credentials"`
+	Executables map[string]string `json:"executables"`
+	PresetNames []string          `json:"preset_names"`
+	ToolDirs    []string          `json:"tool_dirs"`
+	Backend     string            `json:"backend"`
+}
+
+func (e serviceEnv) runnerPaths() servicePaths {
+	p := e.paths
+	root := filepath.Dir(p.Record)
+	p.Plist = filepath.Join(filepath.Dir(p.Plist), runnerLabel+".plist")
+	p.Record = filepath.Join(root, "operator-runner.json")
+	p.Log = filepath.Join(root, "operator-runner", "runner.log")
+	return p
+}
+
+func actualRunnerLogin() (int, string, error) {
+	u, err := user.LookupId(strconv.Itoa(os.Getuid()))
+	if err != nil {
+		return 0, "", err
+	}
+	return os.Getuid(), u.HomeDir, nil
+}
+
+// The injected home is the test trust anchor. Production additionally gets
+// its login home from the account database, independently of the environment.
+func runnerPath(home, path string, uid int, directory bool, private os.FileMode) error {
+	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
+		return errors.New("path is not absolute and clean")
+	}
+	for p := path; ; p = filepath.Dir(p) {
+		st, err := os.Lstat(p)
+		if err != nil {
+			return err
+		}
+		owner, ok := st.Sys().(*syscall.Stat_t)
+		if !ok || (int(owner.Uid) != uid && (within(home, p) || p == home || owner.Uid != 0)) {
+			return errors.New("path ownership mismatch")
+		}
+		if st.Mode().Perm()&0022 != 0 {
+			return errors.New("path writable by others")
+		}
+		wantDir := p != path || directory
+		if wantDir && !st.IsDir() || !wantDir && !st.Mode().IsRegular() {
+			return errors.New("path has wrong type or is a symlink")
+		}
+		if p == path && private != 0 && st.Mode().Perm() != private {
+			return errors.New("path permissions are not private")
+		}
+		if p == home || p == "/" {
+			return nil
+		}
+	}
+}
+
+func readRunnerConfig(path string) (runnerConfig, []byte, error) {
+	var c runnerConfig
+	f, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
+	if err != nil {
+		return c, nil, err
+	}
+	defer f.Close()
+	st, err := f.Stat()
+	if err != nil || !st.Mode().IsRegular() {
+		return c, nil, errors.New("expected a regular runner configuration")
+	}
+	b, err := io.ReadAll(io.LimitReader(f, 16385))
+	if err != nil || len(b) > 16384 {
+		return c, nil, errors.New("cannot read bounded runner configuration")
+	}
+	d := json.NewDecoder(bytes.NewReader(b))
+	d.DisallowUnknownFields()
+	if d.Decode(&c) != nil || d.Decode(new(any)) != io.EOF {
+		return c, nil, errors.New("invalid runner configuration JSON")
+	}
+	return c, b, nil
+}
+
+func (e serviceEnv) checkRunnerConfig(c runnerConfig) error {
+	login := e.runnerLogin
+	if login == nil {
+		login = actualRunnerLogin
+	}
+	uid, home, err := login()
+	switch {
+	case err != nil || uid == 0 || uid != c.UID:
+		return errors.New("runner requires the approved non-root login UID")
+	case home != c.Home || home != e.paths.Home:
+		return errors.New("runner HOME differs from the login home")
+	case c.Backend != "headless":
+		return errors.New("runner backend must be headless")
+	case len(c.PresetNames) == 0 || len(c.PresetNames) > 16 || len(c.Executables) > 16 || len(c.ToolDirs) > 16:
+		return errors.New("runner configuration list size is invalid")
+	}
+	host, _, _ := net.SplitHostPort(c.Address)
+	if !serviceLoopback(c.Address) || !net.ParseIP(host).IsLoopback() {
+		return errors.New("runner daemon must be numeric loopback")
+	}
+	anchor := home
+	if e.runnerLogin == nil {
+		anchor = "" // verify all login-home ancestors in production
+	}
+	if err := runnerPath(anchor, home, uid, true, 0); err != nil {
+		return err
+	}
+	for _, preset := range c.PresetNames {
+		if preset == "se-tools" || len(preset) > verbs.MaxPresetName || !verbs.EpisodeIDPattern.MatchString(preset) {
+			return errors.New("invalid runner preset name")
+		}
+	}
+	if c.Executables[e.paths.Binary] == "" {
+		return errors.New("missing installed World executable")
+	}
+	for path, sum := range c.Executables {
+		if err := runnerPath(home, path, uid, false, 0); err != nil {
+			return err
+		}
+		st, err := os.Stat(path)
+		if err != nil || st.Mode().Perm()&0100 == 0 || sum == "" || serviceFileHash(path) != sum {
+			return errors.New("runner executable changed or is not executable")
+		}
+	}
+	for _, dir := range c.ToolDirs {
+		if strings.ContainsAny(dir, ":\n") {
+			return errors.New("invalid tool directory")
+		}
+		if err := runnerPath(home, dir, uid, true, 0); err != nil {
+			return err
+		}
+	}
+	world, worldErr := store.CanonicalPath(c.World)
+	creds, credsErr := store.CanonicalPath(c.Credentials)
+	if worldErr != nil || credsErr != nil || world != c.World || creds != c.Credentials || world == creds {
+		return errors.New("runner stores must be distinct canonical files")
+	}
+	for _, path := range []string{c.World, c.Credentials} {
+		if err := runnerPath(home, path, uid, false, 0); err != nil {
+			return err
+		}
+	}
+	return nil
+}
+
+func (e serviceEnv) checkRunnerPair(c runnerConfig) error {
+	r, err := readServiceRecord(e.paths)
+	if err != nil || r.Uninstalled || r.Address != c.Address || r.BinarySHA256 != c.Executables[e.paths.Binary] {
+		return errors.New("runner differs from the installed World service")
+	}
+	worldRecord, worldErr := store.CanonicalPath(r.Store)
+	credsRecord, credsErr := store.CanonicalPath(r.Credentials)
+	if worldErr != nil || credsErr != nil || worldRecord != c.World || credsRecord != c.Credentials {
+		return errors.New("runner store pair differs from the installed World service")
+	}
+	h, err := e.health(c.Address)
+	if err != nil || h.Status != "ok" {
+		return errors.New("daemon unavailable")
+	}
+	world, err := store.CanonicalPath(h.DBPath)
+	if err != nil || world != c.World {
+		return errors.New("daemon World store mismatch")
+	}
+	creds, err := store.CanonicalPath(h.CredentialsDB)
+	if err != nil || creds != c.Credentials {
+		return errors.New("daemon credentials store mismatch")
+	}
+	if h.Build.ExecutableSHA256 != c.Executables[e.paths.Binary] {
+		return errors.New("daemon build mismatch; restart pending")
+	}
+	db, err := store.OpenCredentialsReadOnly(runnerDB(c.Credentials))
+	if err != nil {
+		return errors.New("cannot read daemon credentials schema")
+	}
+	defer db.Close()
+	if db.Version() != store.CredentialsSchemaVersion {
+		return errors.New("daemon credentials schema mismatch")
+	}
+	return nil
+}
+
+func runnerPlist(p servicePaths, c runnerConfig, configBytes []byte) string {
+	paths := append(append([]string(nil), c.ToolDirs...), runnerSystemPath)
+	return renderServicePlist(runnerLabel, p, []string{p.Binary, "operator-runner"}, c.Home,
+		[2]string{"HOME", c.Home}, [2]string{"PATH", strings.Join(paths, ":")},
+		[2]string{"WORLD_OPERATOR_CONFIG_SHA256", serviceHash(configBytes)},
+		[2]string{"GIT_TERMINAL_PROMPT", "0"}, [2]string{"GH_PROMPT_DISABLED", "1"})
+}
+
+func (e serviceEnv) installedRunner() (runnerConfig, []byte, error) {
+	p := e.runnerPaths()
+	c, b, err := readRunnerConfig(p.Record)
+	if err != nil {
+		return c, nil, err
+	}
+	if err := e.checkRunnerConfig(c); err != nil {
+		return c, nil, err
+	}
+	for _, path := range []string{p.Record, p.Plist, p.Log} {
+		if err := runnerPath(c.Home, path, c.UID, false, 0600); err != nil {
+			return c, nil, err
+		}
+	}
+	if err := runnerPath(c.Home, filepath.Dir(p.Log), c.UID, true, 0700); err != nil {
+		return c, nil, err
+	}
+	if serviceFileHash(p.Plist) != serviceHash([]byte(runnerPlist(p, c, b))) {
+		return c, nil, errors.New("runner plist differs from approved configuration")
+	}
+	return c, b, nil
+}
+
+// Lock the service root inode, without rewriting any base service file.
+func runnerLifecycleLock(root string) (*os.File, error) {
+	f, err := os.OpenFile(root, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_DIRECTORY, 0)
+	if err != nil {
+		return nil, err
+	}
+	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
+		f.Close()
+		return nil, err
+	}
+	return f, nil
+}
+
+func runServiceRunner(args []string, out, errOut io.Writer, e serviceEnv) int {
+	refuse := func(err error) int {
+		return serviceMessage(errOut, "runner refused", runnerDiagnostic(err), "ailang-worldd service runner --help; inspect the approved configuration", 1)
+	}
+	if len(args) == 0 || isHelpToken(args[len(args)-1]) {
+		fmt.Fprint(out, runnerHelp)
+		return 0
+	}
+	verb := args[0]
+	if e.goos != "darwin" || e.label != worldServiceLabel {
+		return refuse(errors.New("runner requires the macOS World service"))
+	}
+	switch verb {
+	case "install":
+		if len(args) != 3 || args[1] != "--config" {
+			return refuse(errors.New("expected install --config FILE"))
+		}
+	case "cancel":
+		if len(args) != 2 || !verbs.ValidOperatorActionID(args[1]) {
+			return refuse(errors.New("expected cancel ACTION-ID"))
+		}
+	case "start", "stop", "status", "uninstall":
+		if len(args) != 1 {
+			return refuse(errors.New("unexpected runner argument"))
+		}
+	default:
+		return refuse(errors.New("unknown runner operation"))
+	}
+	p := e.runnerPaths()
+	if verb == "status" {
+		if err := e.runnerStatus(out); err != nil {
+			return refuse(err)
+		}
+		return 0
+	}
+
+	if verb == "install" || verb == "stop" || verb == "cancel" || verb == "uninstall" {
+		change := "operator runner " + strings.Join(args, " ")
+		if verb == "install" {
+			c, b, err := readRunnerConfig(args[2])
+			if err != nil {
+				return refuse(errors.New("invalid runner setup file"))
+			}
+			if err := e.checkRunnerConfig(c); err != nil {
+				return refuse(err)
+			}
+			// Display only typed, validated fields from the approved bytes.
+			pretty, _ := json.MarshalIndent(c, "", "  ")
+			change = "install operator runner " + serviceHash(b) + "\n" + string(pretty)
+			if !e.confirm(errOut, change) {
+				return 1
+			}
+			lock, err := runnerLifecycleLock(filepath.Dir(p.Record))
+			if err != nil {
+				return refuse(errors.New("runner lifecycle is busy or unavailable"))
+			}
+			defer lock.Close()
+			if err := e.installRunner(c, b); err != nil {
+				return refuse(err)
+			}
+			fmt.Fprintln(out, "runner installed; next: ailang-worldd service runner start")
+			return 0
+		}
+		if !e.confirm(errOut, change) {
+			return 1
+		}
+	}
+	lock, err := runnerLifecycleLock(filepath.Dir(p.Record))
+	if err != nil {
+		return refuse(errors.New("runner lifecycle is busy or unavailable"))
+	}
+	defer lock.Close()
+	if err := e.changeRunner(verb, args[1:]); err != nil {
+		return refuse(err)
+	}
+	fmt.Fprintln(out, "runner operation complete")
+	return 0
+}
+
+func (e serviceEnv) installRunner(c runnerConfig, b []byte) error {
+	p := e.runnerPaths()
+	if err := e.checkRunnerConfig(c); err != nil {
+		return err
+	}
+	if err := e.checkRunnerPair(c); err != nil {
+		return err
+	}
+	if _, err := e.controller.Print(e.domain, runnerLabel); !errors.Is(err, errServiceNotLoaded) {
+		return errors.New("runner is loaded or its supervisor is unavailable")
+	}
+	for _, path := range []string{p.Plist, p.Record} {
+		if err := serviceSafePath(p.Home, path); err != nil {
+			return err
+		}
+		if _, err := os.Lstat(path); !errors.Is(err, os.ErrNotExist) {
+			return errors.New("runner installation already exists; stop and uninstall it before replacing")
+		}
+	}
+	if err := runnerPath(p.Home, filepath.Dir(p.Plist), c.UID, true, 0); err != nil {
+		return err
+	}
+	if err := runnerPath(p.Home, filepath.Dir(p.Record), c.UID, true, 0700); err != nil {
+		return err
+	}
+	// A pre-existing spool may contain durable evidence. Keep it, including its log.
+	spool := filepath.Dir(p.Log)
+	e.note(spool)
+	if err := os.Mkdir(spool, 0700); err != nil && !errors.Is(err, os.ErrExist) {
+		return err
+	}
+	if err := runnerPath(c.Home, spool, c.UID, true, 0700); err != nil {
+		return err
+	}
+	if _, err := os.Lstat(p.Log); errors.Is(err, os.ErrNotExist) {
+		if err := e.atomic(p.Log, nil, 0600); err != nil {
+			return err
+		}
+	}
+	if err := runnerPath(c.Home, p.Log, c.UID, false, 0600); err != nil {
+		return err
+	}
+	if err := e.atomic(p.Record, b, 0600); err != nil {
+		return err
+	}
+	if err := e.atomic(p.Plist, []byte(runnerPlist(p, c, b)), 0600); err != nil {
+		e.note(p.Record)
+		return errors.Join(err, os.Remove(p.Record))
+	}
+	return nil
+}
+
+func (e serviceEnv) changeRunner(verb string, args []string) error {
+	p := e.runnerPaths()
+	login := e.runnerLogin
+	if login == nil {
+		login = actualRunnerLogin
+	}
+	uid, home, err := login()
+	if err != nil || uid == 0 || home != p.Home {
+		return errors.New("runner requires the non-root login home")
+	}
+	if err := runnerPath(home, filepath.Dir(p.Record), uid, true, 0700); err != nil {
+		return err
+	}
+	switch verb {
+	case "stop", "uninstall":
+		for _, path := range []string{p.Plist, p.Record} {
+			if err := serviceSafePath(p.Home, path); err != nil {
+				return err
+			}
+		}
+		if err := e.controller.Bootout(e.domain, runnerLabel); err != nil && !errors.Is(err, errServiceNotLoaded) {
+			return errors.New("could not stop the companion")
+		}
+		if verb == "uninstall" {
+			for _, path := range []string{p.Plist, p.Record} {
+				e.note(path)
+				if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
+					return err
+				}
+			}
+		}
+		return nil
+	}
+	c, _, err := e.installedRunner()
+	if err != nil {
+		return err
+	}
+	if verb == "cancel" {
+		return runnerCancel(c, args[0], e.terminal.now())
+	}
+	if err := e.checkRunnerPair(c); err != nil {
+		return err
+	}
+	_, err = e.controller.Print(e.domain, runnerLabel)
+	switch {
+	case errors.Is(err, errServiceNotLoaded):
+		if err := e.controller.Bootstrap(e.domain, p.Plist); err != nil {
+			_ = e.controller.Bootout(e.domain, runnerLabel)
+
+			return errors.New("could not bootstrap the companion")
+		}
+	case err != nil:
+		return errors.New("could not inspect the companion")
+	}
+	return nil
+}
+
+func runnerDiagnostic(err error) string {
+	var pathErr *os.PathError
+	if errors.As(err, &pathErr) {
+		return "runner file is absent or inaccessible"
+	}
+	return err.Error()
+}
diff --git a/cmd/ailang-worldd/operator_service_test.go b/cmd/ailang-worldd/operator_service_test.go
new file mode 100644
index 0000000..9e4f35c
--- /dev/null
+++ b/cmd/ailang-worldd/operator_service_test.go
@@ -0,0 +1,329 @@
+package main
+
+import (
+	"bytes"
+	"encoding/json"
+	"errors"
+	"io"
+	"net/http"
+	"os"
+	"path/filepath"
+	"reflect"
+	"strconv"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+type operatorServiceFixture struct {
+	*serviceFixture
+	config     runnerConfig
+	configFile string
+	heartbeat  store.OperatorHeartbeat
+	homeJob    *fakeServiceController
+	health     daemon.HealthResponse
+}
+
+func operatorService(t *testing.T) *operatorServiceFixture {
+	t.Helper()
+	f := newServiceFixture(t)
+	closeWorld, err := storeForServiceTest(f.db)
+	if err != nil {
+		t.Fatal(err)
+	}
+	closeWorld()
+	if code, msg := f.run(f.args...); code != 0 {
+		t.Fatalf("World fixture: %d %s", code, msg)
+	}
+	r, err := readServiceRecord(f.env.paths)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if _, err := store.UpgradeCredentials(r.Credentials); err != nil {
+		t.Fatal(err)
+	}
+	world, err := store.CanonicalPath(r.Store)
+	if err != nil {
+		t.Fatal(err)
+	}
+	creds, err := store.CanonicalPath(r.Credentials)
+	if err != nil {
+		t.Fatal(err)
+	}
+	f.env.runnerLogin = func() (int, string, error) { return os.Getuid(), f.home, nil }
+	o := &operatorServiceFixture{serviceFixture: f, homeJob: f.ctl, configFile: filepath.Join(f.home, "proposed.json"),
+		config:    runnerConfig{f.home, os.Getuid(), r.Address, world, creds, map[string]string{r.Binary: r.BinarySHA256}, []string{"sol"}, []string{filepath.Dir(r.Binary)}, "headless"},
+		heartbeat: store.OperatorHeartbeat{Owner: store.OperatorOwner{Nonce: "fake-runner"}, Protocol: store.OperatorRunnerProtocol, Build: r.BinarySHA256},
+		health:    daemon.HealthResponse{Status: "ok", DBPath: world, CredentialsDB: creds, Build: daemon.BuildIdentity{ExecutableSHA256: r.BinarySHA256}}}
+	f.ctl = &fakeServiceController{runner: true, binary: r.Binary}
+	f.env.controller = f.ctl
+	f.env.client = &http.Client{Transport: serviceTransport(func(req *http.Request) (*http.Response, error) {
+		if req.URL.Host != r.Address || req.URL.Path != "/v1/health" {
+			t.Fatalf("unexpected health URL: %s", req.URL)
+		}
+		if !f.serving {
+			return nil, errors.New("fake daemon went away")
+		}
+		b, _ := json.Marshal(o.health)
+		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(b)), Header: make(http.Header)}, nil
+	})}
+	f.writes = nil
+	return o
+}
+
+func (f *operatorServiceFixture) command(args ...string) (int, string) {
+	return f.run(append([]string{"runner"}, args...)...)
+}
+func (f *operatorServiceFixture) install(t *testing.T) (int, string) {
+	t.Helper()
+	b, err := json.Marshal(f.config)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if err := os.WriteFile(f.configFile, b, 0600); err != nil {
+		t.Fatal(err)
+	}
+	return f.command("install", "--config", f.configFile)
+}
+func (f *operatorServiceFixture) mustInstall(t *testing.T) {
+	t.Helper()
+	if c, s := f.install(t); c != 0 {
+		t.Fatalf("install: %d %s", c, s)
+	}
+}
+
+func TestOperatorServiceFence(t *testing.T) {
+	f := operatorService(t)
+	f.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) { return &fakeTerm{r: strings.NewReader("n\n")}, nil }
+	if c, s := f.install(t); c != 1 || !strings.Contains(s, "aborted") || len(f.writes) != 0 || f.ctl.boots != 0 {
+		t.Fatalf("%d %s writes=%v", c, s, f.writes)
+	}
+	f.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) { return &fakeTerm{r: strings.NewReader("y\n")}, nil }
+	f.mustInstall(t)
+	p := f.env.runnerPaths()
+	b, err := os.ReadFile(p.Plist)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if err := os.WriteFile(p.Plist, append(b, ' '), 0600); err != nil {
+		t.Fatal(err)
+	}
+	before := serviceFileHash(p.Plist)
+	if c, s := f.command("start"); c != 1 || !strings.Contains(s, "differs from approved") || f.ctl.boots != 0 || serviceFileHash(p.Plist) != before {
+		t.Fatalf("%d %s", c, s)
+	}
+	for _, words := range [][]string{{"service", "runner", "start"}, {"operator-runner"}} {
+		got := verbs.AdmitOperator(words, verbs.OperatorPolicy{Start: "A", Finish: "A"})
+		if got.Intent != nil {
+			t.Fatalf("bar admitted %q", words)
+		}
+	}
+}
+
+func TestOperatorServiceOwnership(t *testing.T) {
+	for _, tc := range []struct {
+		name, want string
+		damage     func(*testing.T, *operatorServiceFixture)
+	}{
+		{"root", "non-root", func(t *testing.T, f *operatorServiceFixture) {
+			f.env.runnerLogin = func() (int, string, error) { return 0, f.home, nil }
+			f.config.UID = 0
+		}},
+		{"uid", "non-root", func(t *testing.T, f *operatorServiceFixture) { f.config.UID++ }},
+		{"home", "login home", func(t *testing.T, f *operatorServiceFixture) { f.config.Home = t.TempDir() }},
+		{"writable tool dir", "writable by others", func(t *testing.T, f *operatorServiceFixture) {
+			p := filepath.Join(f.home, "loose")
+			if err := os.Mkdir(p, 0700); err != nil {
+				t.Fatal(err)
+			}
+			if err := os.Chmod(p, 0770); err != nil {
+				t.Fatal(err)
+			}
+			f.config.ToolDirs = []string{p}
+		}},
+		{"symlink tool", "wrong type", func(t *testing.T, f *operatorServiceFixture) {
+			p := filepath.Join(f.home, "linked")
+			if err := os.Symlink(f.env.paths.Binary, p); err != nil {
+				t.Fatal(err)
+			}
+			f.config.Executables[p] = serviceFileHash(f.env.paths.Binary)
+		}},
+		{"changed hash", "executable changed", func(t *testing.T, f *operatorServiceFixture) {
+			f.config.Executables[f.env.paths.Binary] = "sha256:other"
+		}},
+		{"missing executable", "missing installed", func(t *testing.T, f *operatorServiceFixture) { delete(f.config.Executables, f.env.paths.Binary) }},
+		{"non executable", "executable changed", func(t *testing.T, f *operatorServiceFixture) {
+			if err := os.Chmod(f.env.paths.Binary, 0600); err != nil {
+				t.Fatal(err)
+			}
+		}},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			f := operatorService(t)
+			tc.damage(t, f)
+			if c, s := f.install(t); c != 1 || !strings.Contains(s, tc.want) || len(f.writes) != 0 || f.ctl.boots != 0 {
+				t.Fatalf("%d %s writes=%v", c, s, f.writes)
+			}
+		})
+	}
+}
+
+func TestOperatorServiceConfig(t *testing.T) {
+	f := operatorService(t)
+	var b bytes.Buffer
+	if c := runOperatorRunner([]string{"episode", "start", "evil"}, &b, f.env); c != exitUsage || b.Len() != 0 {
+		t.Fatalf("accepted verb argv: %d %q", c, b.String())
+	}
+	raw, err := json.Marshal(f.config)
+	if err != nil {
+		t.Fatal(err)
+	}
+	for _, payload := range []string{string(raw[:len(raw)-1]) + `,"token":"canary-secret"}`, string(raw) + " {}", "{"} {
+		if err := os.WriteFile(f.configFile, []byte(payload), 0600); err != nil {
+			t.Fatal(err)
+		}
+		if c, s := f.command("install", "--config", f.configFile); c != 1 || !strings.Contains(s, "invalid runner setup") || strings.Contains(s, "canary-secret") || len(f.writes) != 0 {
+			t.Fatalf("%d %s", c, s)
+		}
+	}
+	f.config.Backend = "shell"
+	if c, s := f.install(t); c != 1 || !strings.Contains(s, "backend") || len(f.writes) != 0 {
+		t.Fatalf("%d %s", c, s)
+	}
+}
+
+func TestOperatorServicePair(t *testing.T) {
+	for _, tc := range []struct {
+		name, want string
+		change     func(*operatorServiceFixture)
+	}{
+		{"world", "World store mismatch", func(f *operatorServiceFixture) { f.health.DBPath = f.config.Credentials }},
+		{"credentials", "credentials store mismatch", func(f *operatorServiceFixture) { f.health.CredentialsDB = f.config.World }},
+		{"build", "restart pending", func(f *operatorServiceFixture) { f.health.Build.ExecutableSHA256 = "sha256:other" }},
+		{"loopback", "numeric loopback", func(f *operatorServiceFixture) { f.config.Address = "192.0.2.1:7644" }},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			f := operatorService(t)
+			f.mustInstall(t)
+			f.writes = nil
+			tc.change(f)
+			if tc.name == "loopback" {
+				if c, s := f.install(t); c != 1 || !strings.Contains(s, tc.want) || len(f.writes) != 0 {
+					t.Fatalf("%d %s", c, s)
+				}
+				return
+			}
+			if c, s := f.command("start"); c != 1 || !strings.Contains(s, tc.want) || len(f.writes) != 0 || f.ctl.boots != 0 {
+				t.Fatalf("%d %s", c, s)
+			}
+			if err := f.env.runnerTick(f.config, &f.heartbeat); err == nil || f.heartbeat.Owner.Generation != 0 {
+				t.Fatal("incompatible daemon allowed heartbeat")
+			}
+		})
+	}
+}
+
+func TestOperatorServiceStartup(t *testing.T) {
+	f := operatorService(t)
+	f.mustInstall(t)
+	f.serving = false
+	if c, s := f.command("start"); c != 1 || !strings.Contains(s, "daemon unavailable") || f.ctl.boots != 0 {
+		t.Fatalf("%d %s", c, s)
+	}
+	f.serving = true
+	f.ctl.bootErr = errors.New("fake bootstrap failure")
+	if c, s := f.command("start"); c != 1 || !strings.Contains(s, "bootstrap") || f.ctl.boots != 1 || f.ctl.loaded {
+		t.Fatalf("%d %s", c, s)
+	}
+	f.ctl.bootErr = nil
+	if c, s := f.command("start"); c != 0 || !f.ctl.loaded {
+		t.Fatalf("%d %s", c, s)
+	}
+}
+
+func TestOperatorServiceWriteScope(t *testing.T) {
+	f := operatorService(t)
+	base := f.env.paths
+	p := f.env.runnerPaths()
+	files := []string{base.Plist, base.Record, base.Binary, f.config.World, f.config.Credentials}
+	read := func() []string {
+		var h []string
+		for _, p := range files {
+			h = append(h, serviceFileHash(p))
+		}
+		return h
+	}
+	want := read()
+	f.mustInstall(t)
+	for _, path := range []string{p.Plist, p.Record, p.Log} {
+		info, err := os.Stat(path)
+		if err != nil || info.Mode().Perm() != 0600 {
+			t.Fatalf("private mode %s: %v %v", path, info, err)
+		}
+	}
+	if info, err := os.Stat(filepath.Dir(p.Log)); err != nil || info.Mode().Perm() != 0700 {
+		t.Fatalf("spool: %v %v", info, err)
+	}
+	for _, verb := range []string{"start", "stop", "uninstall"} {
+		if c, s := f.command(verb); c != 0 {
+			t.Fatalf("%s: %d %s", verb, c, s)
+		}
+	}
+	if f.homeJob.stops != 0 || f.ctl.loaded || !reflect.DeepEqual(want, read()) || serviceFileHash(p.Plist) != "" || serviceFileHash(p.Record) != "" || serviceFileHash(p.Log) == "" {
+		t.Fatal("companion changed base files, jobs, or retained evidence")
+	}
+	for _, path := range f.writes {
+		if !within(f.home, path) {
+			t.Fatalf("write escaped fake HOME: %s", path)
+		}
+	}
+}
+
+func TestOperatorServiceEnvironment(t *testing.T) {
+	f := operatorService(t)
+	t.Setenv("GH_TOKEN", "canary-secret")
+	t.Setenv("PATH", "/canary-path")
+	f.mustInstall(t)
+	p := f.env.runnerPaths()
+	b, err := os.ReadFile(p.Plist)
+	if err != nil {
+		t.Fatal(err)
+	}
+	s := string(b)
+	for _, want := range []string{runnerLabel, "<string>operator-runner</string>", "<key>HOME</key><string>" + f.home + "</string>", "<key>PATH</key><string>" + f.config.ToolDirs[0] + ":" + runnerSystemPath + "</string>", "GH_PROMPT_DISABLED", "GIT_TERMINAL_PROMPT", serviceFileHash(p.Record)} {
+		if !strings.Contains(s, want) {
+			t.Fatalf("plist lacks %q: %s", want, s)
+		}
+	}
+	for _, path := range []string{p.Plist, p.Record, p.Log} {
+		raw, err := os.ReadFile(path)
+		if err != nil || bytes.Contains(raw, []byte("canary")) {
+			t.Fatalf("ambient environment leaked to %s: %v", path, err)
+		}
+	}
+}
+
+func TestOperatorStatusDisconnected(t *testing.T) {
+	f := operatorService(t)
+	f.mustInstall(t)
+	if c, s := f.command("status"); c != 1 || !strings.Contains(s, "heartbeat unavailable") {
+		t.Fatalf("%d %s", c, s)
+	}
+	if c, s := f.command("start"); c != 0 {
+		t.Fatalf("%d %s", c, s)
+	}
+	if err := f.env.runnerTick(f.config, &f.heartbeat); err != nil {
+		t.Fatal(err)
+	}
+	known := "runner last seen " + strconv.FormatInt(f.heartbeat.SeenAt, 10) + ", generation 1"
+	if c, s := f.command("status"); c != 0 || !strings.Contains(s, known) || !strings.Contains(s, "start: ready") || !strings.Contains(s, "finish blocked: github-auth") {
+		t.Fatalf("%d %s", c, s)
+	}
+	f.serving = false
+	if c, s := f.command("status"); c != 1 || !strings.Contains(s, known) || !strings.Contains(s, "daemon disconnected") || strings.Contains(s, "start: ready") {
+		t.Fatalf("%d %s", c, s)
+	}
+}
diff --git a/cmd/ailang-worldd/quickstart_service_test.go b/cmd/ailang-worldd/quickstart_service_test.go
index a9dbb20..62b4ed3 100644
--- a/cmd/ailang-worldd/quickstart_service_test.go
+++ b/cmd/ailang-worldd/quickstart_service_test.go
@@ -13,6 +13,7 @@ func quickstartServiceLines(md string) []plugindoc.Line {
 }
 
 func TestQuickstartServiceSectionMatchesTheCLIHelp(t *testing.T) {
+	t.Setenv("HOME", t.TempDir())
 	raw, err := os.ReadFile("../../docs/QUICKSTART.md")
 	if err != nil {
 		t.Fatal(err)
@@ -33,7 +34,7 @@ func TestQuickstartServiceSectionMatchesTheCLIHelp(t *testing.T) {
 	if len(base) != 0 {
 		t.Fatal(base)
 	}
-	for _, mutant := range []string{"ailang-worldd frobnicate", "ailang-worldd why 3 --expect-shop x", "ailang-worldd service status --bogus"} {
+	for _, mutant := range []string{"ailang-worldd frobnicate", "ailang-worldd why 3 --expect-shop x", "ailang-worldd service status --bogus", "ailang-worldd service runner status --bogus"} {
 		for _, form := range []string{"\n```sh\n" + mutant + "\n```\n", "\n`" + mutant + "`\n"} {
 			errs := useWorldErrors(quickstartServiceLines(md + form))
 			if len(errs) != len(base)+1 || !strings.Contains(strings.Join(errs, "\n"), mutant) {
diff --git a/cmd/ailang-worldd/service.go b/cmd/ailang-worldd/service.go
index fae44e7..5aa824d 100644
--- a/cmd/ailang-worldd/service.go
+++ b/cmd/ailang-worldd/service.go
@@ -34,6 +34,7 @@ const serviceHelp = `usage: ailang-worldd service install|status|upgrade|stop|st
     [--examples-dir DIR] [--surface-index PATH] [--surface-follow PATH]
     [--surface-origin https://<host>.ts.net[:port]]... [--binary PATH] [--replace]
   status [--json] [--verbose]
+  runner install|start|stop|status|cancel|uninstall (see service runner --help)
   upgrade [--binary FILE] [--pause]
 Install, stop and uninstall require a typed y at the controlling terminal.
 Start uses the approved plist; edited plists require confirmation.
@@ -78,8 +79,9 @@ type serviceEnv struct {
 	client              *http.Client
 	portFree            func(string) error
 	// observe records every intended filesystem mutation, including temp paths.
-	observe func(string)
-	probe   func(string, string) ([]byte, error)
+	observe     func(string)
+	probe       func(string, string) ([]byte, error)
+	runnerLogin func() (int, string, error)
 }
 
 func realServiceEnv() serviceEnv {
@@ -183,7 +185,7 @@ func serviceXML(s string) string {
 	_ = xml.EscapeText(&b, []byte(s))
 	return b.String()
 }
-func renderServicePlist(label string, p servicePaths, args []string, working string) string {
+func renderServicePlist(label string, p servicePaths, args []string, working string, environment ...[2]string) string {
 	var b strings.Builder
 	b.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n<plist version=\"1.0\"><dict>\n")
 	fmt.Fprintf(&b, "<key>Label</key><string>%s</string>\n<key>ProgramArguments</key><array>\n", serviceXML(label))
@@ -192,7 +194,14 @@ func renderServicePlist(label string, p servicePaths, args []string, working str
 	}
 	b.WriteString("</array>\n<key>RunAtLoad</key><true/>\n<key>KeepAlive</key><true/>\n<key>ThrottleInterval</key><integer>30</integer>\n<key>ExitTimeOut</key><integer>30</integer>\n<key>ProcessType</key><string>Standard</string>\n<key>Umask</key><integer>63</integer>\n")
 	fmt.Fprintf(&b, "<key>WorkingDirectory</key><string>%s</string>\n", serviceXML(working))
-	b.WriteString("<key>EnvironmentVariables</key><dict><key>PATH</key><string>/usr/bin:/bin:/usr/sbin:/sbin</string></dict>\n")
+	b.WriteString("<key>EnvironmentVariables</key><dict>")
+	if len(environment) == 0 {
+		environment = [][2]string{{"PATH", runnerSystemPath}}
+	}
+	for _, item := range environment {
+		fmt.Fprintf(&b, "<key>%s</key><string>%s</string>", serviceXML(item[0]), serviceXML(item[1]))
+	}
+	b.WriteString("</dict>\n")
 	fmt.Fprintf(&b, "<key>StandardOutPath</key><string>%s</string>\n<key>StandardErrorPath</key><string>%s</string>\n</dict></plist>\n", serviceXML(p.Log), serviceXML(p.Log))
 	return b.String()
 }
@@ -284,6 +293,9 @@ func runServiceWith(args []string, out, errOut io.Writer, e serviceEnv) int {
 		fmt.Fprint(out, serviceHelp)
 		return 0
 	}
+	if verb == "runner" {
+		return runServiceRunner(args[1:], out, errOut, e)
+	}
 	fs := flag.NewFlagSet("service "+verb, flag.ContinueOnError)
 	fs.SetOutput(errOut)
 	asJSON := fs.Bool("json", false, "JSON status")
diff --git a/cmd/ailang-worldd/service_test.go b/cmd/ailang-worldd/service_test.go
index e2d8429..08f2711 100644
--- a/cmd/ailang-worldd/service_test.go
+++ b/cmd/ailang-worldd/service_test.go
@@ -20,6 +20,7 @@ import (
 
 type fakeServiceController struct {
 	loaded       bool
+	runner       bool
 	boots, stops int
 	key          string
 	bootErr      error
@@ -50,7 +51,7 @@ func (f *fakeServiceController) Bootout(domain, label string) error {
 	return nil
 }
 func (f *fakeServiceController) Print(domain, label string) (string, error) {
-	if !f.loaded {
+	if !f.loaded || (label == runnerLabel && !f.runner) {
 		return "not loaded", errServiceNotLoaded
 	}
 	return "pid = 123\nruns = 1\nlast exit code = 0", nil
@@ -73,7 +74,11 @@ type serviceFixture struct {
 
 func newServiceFixture(t *testing.T) *serviceFixture {
 	t.Helper()
-	f := &serviceFixture{home: t.TempDir(), serving: true}
+	home, err := filepath.EvalSymlinks(t.TempDir())
+	if err != nil {
+		t.Fatal(err)
+	}
+	f := &serviceFixture{home: home, serving: true}
 	f.env = serviceEnv{paths: homeServicePaths(f.home), controller: nil, terminal: sessionEnv{
 		openTerminal: func() (io.ReadWriteCloser, error) { return &fakeTerm{r: strings.NewReader("y\n")}, nil }, now: func() int64 { return 123 }},
 		goos: "darwin", domain: "gui/123", label: worldServiceLabel, wait: 0, portFree: func(string) error { return nil }}
diff --git a/cmd/ailang-worldd/service_upgrade.go b/cmd/ailang-worldd/service_upgrade.go
index 4d44ce0..d921fe9 100644
--- a/cmd/ailang-worldd/service_upgrade.go
+++ b/cmd/ailang-worldd/service_upgrade.go
@@ -198,6 +198,19 @@ func (e serviceEnv) upgrade(candidate string, pause bool, out, errOut io.Writer)
 	if !e.confirm(errOut, change) {
 		return 1
 	}
+	pairLock, err := runnerLifecycleLock(filepath.Dir(e.paths.Record))
+	if err != nil {
+		return refuse("operator runner lifecycle is busy or unavailable", "inspect service runner status and retry")
+	}
+	defer pairLock.Close()
+	for _, path := range []string{e.runnerPaths().Record, e.runnerPaths().Plist} {
+		if _, err := os.Lstat(path); !errors.Is(err, os.ErrNotExist) {
+			return refuse("operator runner is still installed", "stop and uninstall the companion; upgrade World, then review and reinstall the runner")
+		}
+	}
+	if _, err := e.controller.Print(e.domain, runnerLabel); !errors.Is(err, errServiceNotLoaded) {
+		return refuse("cannot verify that the operator runner is unloaded", "inspect service runner status")
+	}
 	lock, err := e.lock()
 	if err != nil {
 		return serviceMessage(errOut, "a service change is in progress", err.Error(), "wait and retry", 1)
diff --git a/docs/QUICKSTART.md b/docs/QUICKSTART.md
index 4115d72..c67cfc6 100644
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -55,6 +55,30 @@ Uninstall keeps the binary, databases, sessions, worktrees and logs.
 An upgrade refuses active sessions unless you explicitly pass `--pause`.
 Linux service support is a later milestone.
 
+### Run the operator companion
+
+With World running, fill in the login UID, canonical paths, preset and `sha256:`-prefixed digest of its installed binary.
+```json
+{
+  "home": "/Users/operator", "uid": 501, "address": "127.0.0.1:7644",
+  "world": "/Users/operator/.ailang/world/world.db", "credentials": "/Users/operator/.ailang/world/credentials.db",
+  "executables": {"/Users/operator/.ailang-world/bin/ailang-worldd": "sha256:<digest>"},
+  "preset_names": ["sol"], "tool_dirs": [], "backend": "headless"
+}
+```
+```sh
+ailang-worldd service runner install --config /absolute/path/to/runner.json
+ailang-worldd service runner start
+ailang-worldd service runner status
+ailang-worldd service runner stop
+ailang-worldd service runner uninstall
+```
+Review and type `y` for install, stop and uninstall. Uninstall before upgrading or replacing the config. Spool and action records are retained.
+Run the checked example with a fake supervisor and disposable stores:
+```sh
+AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./cmd/ailang-worldd -run '^TestOperatorServiceWriteScope$' -count=1 -v
+```
+
 ## Run your agents through World
 
 World launchers now tell Claude Code, Codex and pi their episode, worktree, session
diff --git a/host/store/operator_runner.go b/host/store/operator_runner.go
new file mode 100644
index 0000000..4bebd72
--- /dev/null
+++ b/host/store/operator_runner.go
@@ -0,0 +1,80 @@
+package store
+
+import (
+	"context"
+	"database/sql"
+	"errors"
+)
+
+const OperatorRunnerProtocol = 1
+
+type OperatorHeartbeat struct {
+	Owner                       OperatorOwner
+	Protocol                    int
+	Build, Readiness            string
+	SeenAt, CatalogueGeneration int64
+	Catalogue                   []byte
+}
+
+// Available applies the S4 read gate. An unavailable heartbeat leaves
+// both the durable lease and all effect ownership exactly as recorded.
+func (h OperatorHeartbeat) Available(build, kind string, now int64) bool {
+	return h.Protocol == OperatorRunnerProtocol && h.Build == build &&
+		h.Owner.Generation > 0 && h.Owner.Nonce != "" &&
+		h.SeenAt <= now && h.SeenAt > now-10 &&
+		(h.Readiness == "both" || h.Readiness == kind)
+}
+
+func (c *Credentials) OperatorHeartbeat(ctx context.Context) (h OperatorHeartbeat, err error) {
+	if err := c.operatorAccess(ctx); err != nil {
+		return h, err
+	}
+	err = c.db.QueryRowContext(ctx, `SELECT owner, generation, protocol, build, readiness, seen_at,
+		catalogue_generation, catalogue FROM operator_runner_state WHERE singleton = 1`).
+		Scan(&h.Owner.Nonce, &h.Owner.Generation, &h.Protocol, &h.Build, &h.Readiness, &h.SeenAt, &h.CatalogueGeneration, &h.Catalogue)
+	if errors.Is(err, sql.ErrNoRows) {
+		err = ErrOperatorRunner
+	}
+	return h, err
+}
+
+// PutOperatorHeartbeat publishes a complete generation in one statement.
+// A new owner (generation zero) must hold the kernel runner lock for this store.
+func (c *Credentials) PutOperatorHeartbeat(ctx context.Context, h OperatorHeartbeat) (int64, error) {
+	if err := c.operatorAccess(ctx); err != nil {
+		return 0, err
+	}
+	var generation int64
+	if h.Owner.Generation == 0 {
+		err := c.db.QueryRowContext(ctx, `INSERT INTO operator_runner_state
+			(singleton,generation,owner,protocol,build,readiness,seen_at,catalogue_generation,catalogue)
+			VALUES(1,1,?,?,?,?,?,1,?)
+			ON CONFLICT(singleton) DO UPDATE SET generation = generation+1, owner = excluded.owner,
+			protocol = excluded.protocol, build = excluded.build, readiness = excluded.readiness,
+			seen_at = excluded.seen_at, catalogue_generation = catalogue_generation+1,
+			catalogue = excluded.catalogue RETURNING generation`,
+			h.Owner.Nonce, h.Protocol, h.Build, h.Readiness, h.SeenAt, h.Catalogue).Scan(&generation)
+		return generation, err
+	}
+	err := c.db.QueryRowContext(ctx, `UPDATE operator_runner_state SET readiness=?, seen_at=?,
+		catalogue_generation=catalogue_generation+1, catalogue=?
+		WHERE singleton=1 AND owner=? AND generation=? AND protocol=? AND build=? AND seen_at<=?
+		RETURNING generation`, h.Readiness, h.SeenAt, h.Catalogue, h.Owner.Nonce,
+		h.Owner.Generation, h.Protocol, h.Build, h.SeenAt).Scan(&generation)
+	if errors.Is(err, sql.ErrNoRows) {
+		err = ErrOperatorRunner
+	}
+	return generation, err
+}
+
+// CancelUnstartedOperatorAction never removes or releases a reserved effect.
+func (c *Credentials) CancelUnstartedOperatorAction(ctx context.Context, id string, now int64) error {
+	if err := c.operatorAccess(ctx); err != nil {
+		return err
+	}
+	result, err := c.db.ExecContext(ctx, `UPDATE operator_actions SET state='cancelled', terminal_at=?,
+		generation=generation+1 WHERE id=? AND state IN ('requested','confirmed')
+		AND owner IS NULL AND expires_at>? AND NOT EXISTS (SELECT 1 FROM operator_action_steps
+			WHERE action_id=? AND kind<>'operation')`, now, id, now, id)
+	return operatorChangedOne(result, err)
+}
````

### Patch M3

````diff M3
diff --git a/cmd/ailang-worldd/operator_service_test.go b/cmd/ailang-worldd/operator_service_test.go
index 9e4f35c..2da9ed9 100644
--- a/cmd/ailang-worldd/operator_service_test.go
+++ b/cmd/ailang-worldd/operator_service_test.go
@@ -8,10 +8,11 @@ import (
 	"net/http"
 	"os"
 	"path/filepath"
-	"reflect"
 	"strconv"
 	"strings"
+	"syscall"
 	"testing"
+	"time"
 
 	"github.com/sunholo-data/ailang-world/host/daemon"
 	"github.com/sunholo-data/ailang-world/host/store"
@@ -45,14 +46,7 @@ func operatorService(t *testing.T) *operatorServiceFixture {
 	if _, err := store.UpgradeCredentials(r.Credentials); err != nil {
 		t.Fatal(err)
 	}
-	world, err := store.CanonicalPath(r.Store)
-	if err != nil {
-		t.Fatal(err)
-	}
-	creds, err := store.CanonicalPath(r.Credentials)
-	if err != nil {
-		t.Fatal(err)
-	}
+	world, creds := r.Store, r.Credentials // fixture paths are physically named
 	f.env.runnerLogin = func() (int, string, error) { return os.Getuid(), f.home, nil }
 	o := &operatorServiceFixture{serviceFixture: f, homeJob: f.ctl, configFile: filepath.Join(f.home, "proposed.json"),
 		config:    runnerConfig{f.home, os.Getuid(), r.Address, world, creds, map[string]string{r.Binary: r.BinarySHA256}, []string{"sol"}, []string{filepath.Dir(r.Binary)}, "headless"},
@@ -95,12 +89,18 @@ func (f *operatorServiceFixture) mustInstall(t *testing.T) {
 	}
 }
 
+func (f *operatorServiceFixture) refused(t *testing.T, want string, code int, message string) {
+	t.Helper()
+	if code != 1 || !strings.Contains(message, want) || len(f.writes) != 0 || f.ctl.boots != 0 {
+		t.Fatalf("refusal %q: %d %s writes=%v", want, code, message, f.writes)
+	}
+}
+
 func TestOperatorServiceFence(t *testing.T) {
 	f := operatorService(t)
 	f.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) { return &fakeTerm{r: strings.NewReader("n\n")}, nil }
-	if c, s := f.install(t); c != 1 || !strings.Contains(s, "aborted") || len(f.writes) != 0 || f.ctl.boots != 0 {
-		t.Fatalf("%d %s writes=%v", c, s, f.writes)
-	}
+	c, msg := f.install(t)
+	f.refused(t, "aborted", c, msg)
 	f.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) { return &fakeTerm{r: strings.NewReader("y\n")}, nil }
 	f.mustInstall(t)
 	p := f.env.runnerPaths()
@@ -164,9 +164,8 @@ func TestOperatorServiceOwnership(t *testing.T) {
 		t.Run(tc.name, func(t *testing.T) {
 			f := operatorService(t)
 			tc.damage(t, f)
-			if c, s := f.install(t); c != 1 || !strings.Contains(s, tc.want) || len(f.writes) != 0 || f.ctl.boots != 0 {
-				t.Fatalf("%d %s writes=%v", c, s, f.writes)
-			}
+			c, msg := f.install(t)
+			f.refused(t, tc.want, c, msg)
 		})
 	}
 }
@@ -190,9 +189,8 @@ func TestOperatorServiceConfig(t *testing.T) {
 		}
 	}
 	f.config.Backend = "shell"
-	if c, s := f.install(t); c != 1 || !strings.Contains(s, "backend") || len(f.writes) != 0 {
-		t.Fatalf("%d %s", c, s)
-	}
+	c, msg := f.install(t)
+	f.refused(t, "backend", c, msg)
 }
 
 func TestOperatorServicePair(t *testing.T) {
@@ -211,14 +209,12 @@ func TestOperatorServicePair(t *testing.T) {
 			f.writes = nil
 			tc.change(f)
 			if tc.name == "loopback" {
-				if c, s := f.install(t); c != 1 || !strings.Contains(s, tc.want) || len(f.writes) != 0 {
-					t.Fatalf("%d %s", c, s)
-				}
+				c, msg := f.install(t)
+				f.refused(t, tc.want, c, msg)
 				return
 			}
-			if c, s := f.command("start"); c != 1 || !strings.Contains(s, tc.want) || len(f.writes) != 0 || f.ctl.boots != 0 {
-				t.Fatalf("%d %s", c, s)
-			}
+			c, msg := f.command("start")
+			f.refused(t, tc.want, c, msg)
 			if err := f.env.runnerTick(f.config, &f.heartbeat); err == nil || f.heartbeat.Owner.Generation != 0 {
 				t.Fatal("incompatible daemon allowed heartbeat")
 			}
@@ -234,6 +230,13 @@ func TestOperatorServiceStartup(t *testing.T) {
 		t.Fatalf("%d %s", c, s)
 	}
 	f.serving = true
+	credSQL(t, f.config.Credentials, "DROP TABLE operator_action_steps; DROP TABLE operator_actions; DROP TABLE operator_runner_state; PRAGMA user_version=9")
+	if c, msg := f.command("start"); c != 1 || !strings.Contains(msg, "schema mismatch") || f.ctl.boots != 0 {
+		t.Fatalf("schema skew: %d %s", c, msg)
+	}
+	if _, err := store.UpgradeCredentials(f.config.Credentials); err != nil {
+		t.Fatal(err)
+	}
 	f.ctl.bootErr = errors.New("fake bootstrap failure")
 	if c, s := f.command("start"); c != 1 || !strings.Contains(s, "bootstrap") || f.ctl.boots != 1 || f.ctl.loaded {
 		t.Fatalf("%d %s", c, s)
@@ -244,17 +247,38 @@ func TestOperatorServiceStartup(t *testing.T) {
 	}
 }
 
+func TestOperatorRunnerKernelLock(t *testing.T) {
+	p := filepath.Join(t.TempDir(), "runner.lock")
+	f, err := runnerKernelLock(p)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer f.Close()
+	done := make(chan error, 1)
+	go func() {
+		second, err := runnerKernelLock(p)
+		if second != nil {
+			second.Close()
+		}
+		done <- err
+	}()
+	select {
+	case err := <-done:
+		if !errors.Is(err, syscall.EWOULDBLOCK) {
+			t.Fatalf("contender: %v", err)
+		}
+	case <-time.After(time.Second):
+		t.Fatal("runner lock waited")
+	}
+}
+
 func TestOperatorServiceWriteScope(t *testing.T) {
 	f := operatorService(t)
 	base := f.env.paths
 	p := f.env.runnerPaths()
-	files := []string{base.Plist, base.Record, base.Binary, f.config.World, f.config.Credentials}
-	read := func() []string {
-		var h []string
-		for _, p := range files {
-			h = append(h, serviceFileHash(p))
-		}
-		return h
+	read := func() [5]string {
+		return [5]string{serviceFileHash(base.Plist), serviceFileHash(base.Record), serviceFileHash(base.Binary),
+			serviceFileHash(f.config.World), serviceFileHash(f.config.Credentials)}
 	}
 	want := read()
 	f.mustInstall(t)
@@ -272,7 +296,7 @@ func TestOperatorServiceWriteScope(t *testing.T) {
 			t.Fatalf("%s: %d %s", verb, c, s)
 		}
 	}
-	if f.homeJob.stops != 0 || f.ctl.loaded || !reflect.DeepEqual(want, read()) || serviceFileHash(p.Plist) != "" || serviceFileHash(p.Record) != "" || serviceFileHash(p.Log) == "" {
+	if f.homeJob.stops != 0 || f.ctl.loaded || want != read() || serviceFileHash(p.Plist) != "" || serviceFileHash(p.Record) != "" || serviceFileHash(p.Log) == "" {
 		t.Fatal("companion changed base files, jobs, or retained evidence")
 	}
 	for _, path := range f.writes {
@@ -309,15 +333,15 @@ func TestOperatorServiceEnvironment(t *testing.T) {
 func TestOperatorStatusDisconnected(t *testing.T) {
 	f := operatorService(t)
 	f.mustInstall(t)
-	if c, s := f.command("status"); c != 1 || !strings.Contains(s, "heartbeat unavailable") {
-		t.Fatalf("%d %s", c, s)
-	}
 	if c, s := f.command("start"); c != 0 {
 		t.Fatalf("%d %s", c, s)
 	}
 	if err := f.env.runnerTick(f.config, &f.heartbeat); err != nil {
 		t.Fatal(err)
 	}
+	if string(f.heartbeat.Catalogue) != `["sol"]` {
+		t.Fatalf("unexpected catalogue: %q", f.heartbeat.Catalogue)
+	}
 	known := "runner last seen " + strconv.FormatInt(f.heartbeat.SeenAt, 10) + ", generation 1"
 	if c, s := f.command("status"); c != 0 || !strings.Contains(s, known) || !strings.Contains(s, "start: ready") || !strings.Contains(s, "finish blocked: github-auth") {
 		t.Fatalf("%d %s", c, s)
diff --git a/host/store/operator_runner_test.go b/host/store/operator_runner_test.go
new file mode 100644
index 0000000..916800b
--- /dev/null
+++ b/host/store/operator_runner_test.go
@@ -0,0 +1,144 @@
+package store
+
+import (
+	"context"
+	"errors"
+	"reflect"
+	"testing"
+)
+
+func TestOperatorRunnerUnavailable(t *testing.T) {
+	c := openTestCredentials(t)
+	ctx := testCtx(t)
+	if _, err := c.OperatorHeartbeat(ctx); !errors.Is(err, ErrOperatorRunner) {
+		t.Fatalf("empty heartbeat: %v", err)
+	}
+	h := OperatorHeartbeat{Owner: OperatorOwner{Generation: 1, Nonce: "first"}, Protocol: OperatorRunnerProtocol, Build: "test", Readiness: "start", SeenAt: 10}
+	for _, tc := range []struct {
+		name, build, kind string
+		now               int64
+		change            func(*OperatorHeartbeat)
+		want              bool
+	}{
+		{"now", "test", "start", 10, nil, true},
+		{"future", "test", "start", 9, nil, false},
+		{"just fresh", "test", "start", 19, nil, true},
+		{"exactly stale", "test", "start", 20, nil, false},
+		{"wrong build", "other", "start", 10, nil, false},
+		{"wrong protocol", "test", "start", 10, func(h *OperatorHeartbeat) { h.Protocol++ }, false},
+		{"no generation", "test", "start", 10, func(h *OperatorHeartbeat) { h.Owner.Generation = 0 }, false},
+		{"no owner", "test", "start", 10, func(h *OperatorHeartbeat) { h.Owner.Nonce = "" }, false},
+		{"finish not ready", "test", "finish", 10, nil, false},
+		{"both ready", "test", "finish", 10, func(h *OperatorHeartbeat) { h.Readiness = "both" }, true},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			x := h
+			if tc.change != nil {
+				tc.change(&x)
+			}
+			if got := x.Available(tc.build, tc.kind, tc.now); got != tc.want {
+				t.Fatalf("available = %t", got)
+			}
+		})
+	}
+}
+
+func TestOperatorCatalogueAdmission(t *testing.T) {
+	c := openTestCredentials(t)
+	ctx := testCtx(t)
+	a := OperatorHeartbeat{Owner: OperatorOwner{Nonce: "first"}, Protocol: OperatorRunnerProtocol, Build: "test", SeenAt: 11, Readiness: "start", Catalogue: []byte("one")}
+	publish := func(h OperatorHeartbeat, gen, catalogueGen int64) OperatorHeartbeat {
+		t.Helper()
+		gotGen, err := c.PutOperatorHeartbeat(ctx, h)
+		got, readErr := c.OperatorHeartbeat(ctx)
+		h.Owner.Generation, h.CatalogueGeneration = gen, catalogueGen
+		if gotGen != gen || err != nil || readErr != nil || !reflect.DeepEqual(got, h) {
+			t.Fatalf("heartbeat %+v, generation %d: %v %v; want %+v", got, gotGen, err, readErr, h)
+		}
+		return h
+	}
+	first := publish(a, 1, 1)
+	a.Owner.Generation, a.SeenAt, a.Catalogue = 1, 11, []byte("two")
+	publish(a, 1, 2) // same-second tick is a valid refresh
+	a.SeenAt++
+	second := publish(a, 1, 3)
+	if string(first.Catalogue) != "one" {
+		t.Fatal("published generation changed in-place")
+	}
+	for _, damage := range []func(*OperatorHeartbeat){
+		func(h *OperatorHeartbeat) { h.Owner.Nonce = "stale" },
+		func(h *OperatorHeartbeat) { h.Owner.Generation++ },
+		func(h *OperatorHeartbeat) { h.Protocol++ },
+		func(h *OperatorHeartbeat) { h.Build = "other" },
+		func(h *OperatorHeartbeat) { h.SeenAt-- },
+	} {
+		x := second
+		damage(&x)
+		if _, err := c.PutOperatorHeartbeat(ctx, x); !errors.Is(err, ErrOperatorRunner) {
+			t.Fatalf("stale update admitted: %+v %v", x, err)
+		}
+		got, err := c.OperatorHeartbeat(ctx)
+		if err != nil || !reflect.DeepEqual(got, second) {
+			t.Fatalf("rejected write changed row: %+v %v", got, err)
+		}
+	}
+	a.Owner.Nonce, a.Owner.Generation = "restart", 0
+	publish(a, 2, 4)
+}
+
+func TestOperatorFinishReadiness(t *testing.T) {
+	c := openTestCredentials(t)
+	cas, runner := operatorClaimFixture(t, c)
+	ctx := testCtx(t)
+	operatorExec(t, c, `UPDATE operator_actions SET kind='finish'`)
+	before := getOperator(t, c, cas.ID, 1590)
+	if err := c.ClaimOperatorAction(ctx, cas, runner, 1590); !errors.Is(err, ErrOperatorRunner) {
+		t.Fatalf("finish claim: %v", err)
+	}
+	if got := getOperator(t, c, cas.ID, 1590); !reflect.DeepEqual(got, before) {
+		t.Fatalf("blocked claim changed row: %+v", got)
+	}
+	operatorExec(t, c, `UPDATE operator_actions SET kind='start'`)
+	if err := c.ClaimOperatorAction(ctx, cas, runner, 1590); err != nil {
+		t.Fatalf("start claim: %v", err)
+	}
+}
+
+func TestOperatorCancelInflight(t *testing.T) {
+	c := openTestCredentials(t)
+	cas, runner := operatorClaimFixture(t, c)
+	ctx := testCtx(t)
+	if err := c.ClaimOperatorAction(ctx, cas, runner, 1590); err != nil {
+		t.Fatal(err)
+	}
+	operatorExec(t, c, `INSERT INTO operator_action_steps
+		SELECT 'launch','attempt',id,'reserved',1,owner,1590,NULL FROM operator_actions`)
+	refuse := func(callCtx context.Context, id string, now int64, want error) {
+		t.Helper()
+		before := getOperator(t, c, id, now)
+		err := c.CancelUnstartedOperatorAction(callCtx, id, now)
+		if got := getOperator(t, c, id, now); !errors.Is(err, want) || !reflect.DeepEqual(got, before) {
+			t.Fatalf("cancel refusal changed %s: %+v, %v", id, got, err)
+		}
+	}
+	refuse(ctx, cas.ID, 1591, ErrOperatorConflict)
+	operatorExec(t, c, "UPDATE operator_actions SET state='confirmed', owner=NULL")
+	refuse(ctx, cas.ID, 1591, ErrOperatorConflict)
+	var owner string
+	if err := c.db.QueryRow(`SELECT owner FROM operator_action_steps WHERE kind='launch' AND phase='reserved'`).Scan(&owner); err != nil || owner != runner.Nonce {
+		t.Fatalf("reservation lost: %s %v", owner, err)
+	}
+	pending := putOperator(t, c, operatorRequest(t, c, 2, "pending", "start"), 1591)
+	refuse(ctx, pending.ID, pending.ExpiresAt, ErrOperatorConflict)
+	refuse(context.Background(), pending.ID, 1591, ErrNoDeadline)
+	if err := c.CancelUnstartedOperatorAction(ctx, "missing", 1591); !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("cancel of missing action: %v", err)
+	}
+	if err := c.CancelUnstartedOperatorAction(ctx, pending.ID, 1591); err != nil {
+		t.Fatal(err)
+	}
+	if got := getOperator(t, c, pending.ID, 1591); got.State != "cancelled" || got.Generation != pending.Generation+1 || got.TerminalAt != 1591 {
+		t.Fatalf("cancellation readback: %+v", got)
+	}
+	refuse(ctx, pending.ID, 1591, ErrOperatorConflict)
+}
````

### Patch S4b

````diff S4b
diff --git a/cmd/ailang-worldd/episode.go b/cmd/ailang-worldd/episode.go
index c802b69..19b9f0f 100644
--- a/cmd/ailang-worldd/episode.go
+++ b/cmd/ailang-worldd/episode.go
@@ -361,6 +361,13 @@ func runEpisodeStart(args []string, stdout, stderr io.Writer, env sessionEnv) in
 	if err := checkDaemon(health, req, hostPort, self); err != nil {
 		return fail("refusing: %v", err)
 	}
+	effectLock, err := terminalEpisodeLock(req.credPath, episode, nil)
+	if err != nil {
+		return fail("refusing: %v", err)
+	}
+	if effectLock != nil {
+		defer effectLock.Close()
+	}
 	if *request {
 		return createMintRequest(req, health.Build, hostPort, *agent, *model, taskPath, mcpFile, stdout, stderr, env)
 	}
diff --git a/cmd/ailang-worldd/episode_finish.go b/cmd/ailang-worldd/episode_finish.go
index da0d65f..4d42252 100644
--- a/cmd/ailang-worldd/episode_finish.go
+++ b/cmd/ailang-worldd/episode_finish.go
@@ -125,6 +125,13 @@ func runEpisodeFinishWith(args []string, stdout, stderr io.Writer, env sessionEn
 	if same, _ := store.SameDatabase(h.CredentialsDB, credentials); !same {
 		return fail(fmt.Sprintf("the daemon serves credentials store %s; preset %s names %s (F6)", h.CredentialsDB, *preset, credentials), "the operator proof must be beside the served credentials store", "fix the preset daemon")
 	}
+	effectLock, err := terminalEpisodeLock(credentials, ep, nil)
+	if err != nil {
+		return fail(err.Error(), "the companion owns this episode", "ailang-worldd service runner status")
+	}
+	if effectLock != nil {
+		defer effectLock.Close()
+	}
 	_, rep, err := fetchReport(ctx, c, ep)
 	if err != nil && !errors.Is(err, errNoCalls) {
 		return fail(err.Error(), "finish needs the episode report", "ailang-worldd episode report "+ep)
diff --git a/cmd/ailang-worldd/episode_launch.go b/cmd/ailang-worldd/episode_launch.go
index 95d0f57..a9453c9 100644
--- a/cmd/ailang-worldd/episode_launch.go
+++ b/cmd/ailang-worldd/episode_launch.go
@@ -250,12 +250,12 @@ func runEpisodeLaunch(args []string, stdout, stderr io.Writer) int {
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
@@ -277,6 +277,17 @@ func launchEpisode(ep, name, presetsDir string, headless bool, stdout, stderr io
 	if same, _ := store.SameDatabase(h.DBPath, p.DB); !same {
 		return fail(presetRefusal(fmt.Sprintf("the daemon serves store %s; preset %s names %s (F6)", h.DBPath, name, p.DB), "a preset must bind the served store", "fix the preset daemon, or run service status"))
 	}
+	credsPath, err := credentialsFor(p.DB, p.CredentialsDB)
+	if err != nil {
+		return fail(err)
+	}
+	effectLock, err := terminalEpisodeLock(credsPath, ep, held)
+	if err != nil {
+		return fail(err)
+	}
+	if effectLock != nil {
+		defer effectLock.Close()
+	}
 	c := launchClient(p.Addr)
 	_, rep, err := fetchReport(ctx, c, ep)
 	if err != nil && !errors.Is(err, errNoCalls) {
@@ -421,10 +432,6 @@ func launchEpisode(ep, name, presetsDir string, headless bool, stdout, stderr io
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
diff --git a/cmd/ailang-worldd/episode_preset.go b/cmd/ailang-worldd/episode_preset.go
index a8faa6b..772f04c 100644
--- a/cmd/ailang-worldd/episode_preset.go
+++ b/cmd/ailang-worldd/episode_preset.go
@@ -279,6 +279,13 @@ func runPresetEpisodeStart(fs *flag.FlagSet, pos []string, name, dir, profilesDi
 	if err := checkPresetDaemon(h, req, p, p.Daemon.Addr); err != nil {
 		return fail(err)
 	}
+	effectLock, err := terminalEpisodeLock(req.credPath, ep, nil)
+	if err != nil {
+		return fail(err)
+	}
+	if effectLock != nil {
+		defer effectLock.Close()
+	}
 	req.afterToken = func(token string) error {
 		if taskRecord != nil {
 			f, err := os.OpenFile(taskCopy, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
@@ -308,7 +315,7 @@ func runPresetEpisodeStart(fs *flag.FlagSet, pos []string, name, dir, profilesDi
 		return code
 	}
 	if launching {
-		return launchEpisode(ep, name, dir, headless, out, errOut, fs.Lookup("hosts-dir").Value.String())
+		return launchEpisode(ep, name, dir, headless, out, errOut, effectLock, fs.Lookup("hosts-dir").Value.String())
 	}
 	taskPath := ""
 	if taskRecord != nil {
diff --git a/cmd/ailang-worldd/episode_resume.go b/cmd/ailang-worldd/episode_resume.go
index 29118e8..5e2d1dd 100644
--- a/cmd/ailang-worldd/episode_resume.go
+++ b/cmd/ailang-worldd/episode_resume.go
@@ -172,6 +172,13 @@ func runEpisodeResume(fs *flag.FlagSet, pos []string, ep, name, presetsDir, prof
 	if err := checkDaemon(h, req, p.Addr, episodeSelfBuild()); err != nil {
 		return fail(err)
 	}
+	effectLock, err := terminalEpisodeLock(req.credPath, ep, nil)
+	if err != nil {
+		return fail(err)
+	}
+	if effectLock != nil {
+		defer effectLock.Close()
+	}
 	presetPath, err := expandPresetPath(filepath.Join(presetsDir, name+".json"))
 	if err != nil {
 		return fail(err)
diff --git a/cmd/ailang-worldd/operator_episode_lock.go b/cmd/ailang-worldd/operator_episode_lock.go
new file mode 100644
index 0000000..bc960bc
--- /dev/null
+++ b/cmd/ailang-worldd/operator_episode_lock.go
@@ -0,0 +1,68 @@
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
+	canonical, err := store.CanonicalPath(creds)
+	if err != nil {
+		return nil, err
+	}
+	// The singleton runner lock persists through stop and uninstall, and marks
+	// this canonical store as managed. Its existence never implies liveness.
+	if _, err := os.Lstat(canonical + ".operator-runner.lock"); errors.Is(err, os.ErrNotExist) {
+		if held == nil {
+			return nil, nil
+		}
+		return nil, errors.New("operator store changed during episode start")
+	} else if err != nil {
+		return nil, err
+	}
+	path := operatorEpisodeLockPath(canonical, episode)
+	if held != nil {
+		if held.Name() != path {
+			return nil, errors.New("operator store changed during episode start")
+		}
+		return nil, nil // caller retains the same open file description
+	}
+	lock, lockErr := runnerKernelLock(path)
+	if lockErr != nil && !errors.Is(lockErr, syscall.EWOULDBLOCK) {
+		return nil, lockErr
+	}
+	read := func() (string, error) {
+		db, err := store.OpenCredentialsReadOnly(runnerDB(canonical))
+		if err != nil {
+			return "", err
+		}
+		defer db.Close()
+		ctx, cancel := budgetContext(4 * time.Second)
+		defer cancel()
+		return db.OperatorOwnedEpisode(ctx, episode)
+	}
+	id, readErr := read()
+	if lockErr == nil && readErr == nil && id == "" {
+		return lock, nil
+	}
+	if lock != nil {
+		lock.Close()
+	}
+	if readErr != nil {
+		return nil, fmt.Errorf("%w: could not read action ownership; run ailang-worldd service runner status", store.ErrOperatorConflict)
+	}
+	if verbs.ValidOperatorActionID(id) {
+		return nil, fmt.Errorf("%w: active action %s; run ailang-worldd service runner status and reconcile at the terminal", store.ErrOperatorConflict, id)
+	}
+	return nil, fmt.Errorf("%w: episode lock held; run ailang-worldd service runner status", store.ErrOperatorConflict)
+}
diff --git a/cmd/ailang-worldd/operator_episode_lock_test.go b/cmd/ailang-worldd/operator_episode_lock_test.go
new file mode 100644
index 0000000..3e83210
--- /dev/null
+++ b/cmd/ailang-worldd/operator_episode_lock_test.go
@@ -0,0 +1,127 @@
+package main
+
+import (
+	"bytes"
+	"errors"
+	"io"
+	"os"
+	"path/filepath"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+const activeOperatorFixtureID = "oa-0123456789abcdef0123456789abcdef"
+
+func holdOperatorEpisode(t *testing.T, creds, ep string) *os.File {
+	t.Helper()
+	if _, err := store.UpgradeCredentials(creds); err != nil {
+		t.Fatal(err)
+	}
+	canonical, err := store.CanonicalPath(creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	// SQL-only ownership fixture. The sealed key and claim are tested in
+	// host/authority, and this row grants no effect in S4b.
+	credSQL(t, creds, `INSERT INTO operator_actions
+		(id,episode_id,kind,review,digest,action,challenge,created_at,expires_at,state,generation,owner,approved_by,artifact)
+		VALUES ('`+activeOperatorFixtureID+`','`+ep+`','start',zeroblob(32),zeroblob(32),x'01',zeroblob(32),1,99999,'claimed',2,'fixture-runner','ssh:fixture',x'01')`)
+	process, err := runnerKernelLock(canonical + ".operator-runner.lock")
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { process.Close() })
+	lock, err := runnerKernelLock(operatorEpisodeLockPath(canonical, ep))
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { lock.Close() })
+	return lock
+}
+
+func TestOperatorTerminalExclusion(t *testing.T) {
+	t.Setenv("HOME", t.TempDir())
+	f := newFinishFixture(t)
+	held := holdOperatorEpisode(t, f.health.CredentialsDB, "ep1")
+	before := serviceFileHash(f.health.CredentialsDB)
+	invoke := func(name string, run func() (int, string)) {
+		t.Helper()
+		type result struct {
+			code   int
+			output string
+		}
+		done := make(chan result, 1)
+		go func() { code, s := run(); done <- result{code, s} }()
+		select {
+		case got := <-done:
+			if got.code != exitUsage || !strings.Contains(got.output, "active action "+activeOperatorFixtureID+";") ||
+				!strings.Contains(got.output, "service runner status") {
+				t.Fatalf("%s: %d %s", name, got.code, got.output)
+			}
+		case <-time.After(5 * time.Second):
+			held.Close() // unblock a waiting mutation before leaving the fixture
+			select {
+			case <-done:
+			case <-time.After(time.Second):
+			}
+			t.Fatalf("%s waited for the episode lock", name)
+		}
+	}
+	finishEntered := false
+	finishEnv := sessionEnv{openTerminal: func() (io.ReadWriteCloser, error) {
+		finishEntered = true
+		return &fakeTerm{r: strings.NewReader("n\n")}, nil
+	}}
+	invoke("finish", func() (int, string) { c, _, s := f.run(finishEnv); return c, s })
+	if !finishEntered || len(f.requests) != 1 || f.requests[0] != "health" {
+		t.Fatalf("finish effects or F1 reorder: %v", f.requests)
+	}
+	invoke("launch", func() (int, string) {
+		var b bytes.Buffer
+		c := runEpisodeLaunch([]string{"ep1", "--preset", "fixture", "--presets-dir", f.dir, "--headless"}, &b, &b)
+		return c, b.String()
+	})
+	if len(f.requests) != 2 || f.requests[1] != "health" || serviceFileHash(f.health.CredentialsDB) != before {
+		t.Fatalf("launch effects or altered ownership: %v", f.requests)
+	}
+	if err := held.Close(); err != nil {
+		t.Fatal(err)
+	}
+	if lock, err := terminalEpisodeLock(f.health.CredentialsDB, "ep1", nil); lock != nil || !errors.Is(err, store.ErrOperatorConflict) {
+		t.Fatalf("lost kernel lock released durable claim: %v %v", lock, err)
+	}
+	if lock, err := terminalEpisodeLock(f.health.CredentialsDB, "other", nil); err != nil || lock == nil {
+		t.Fatalf("unrelated episode blocked: %v %v", lock, err)
+	} else {
+		lock.Close()
+	}
+	db, err := store.OpenCredentialsReadOnly(f.health.CredentialsDB)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer db.Close()
+	if got, err := db.OperatorOwnedEpisode(boundedTestContext(t), "ep1"); err != nil || got != activeOperatorFixtureID {
+		t.Fatalf("owner readback: %s %v", got, err)
+	}
+}
+
+func TestOperatorTerminalExclusionStart(t *testing.T) {
+	r := newEpisodeRun(t)
+	health := matchingHealth(t, r.f, r.profiles)
+	stubDaemon(t, health, nil)
+	holdOperatorEpisode(t, health.CredentialsDB, "ep1")
+	before := serviceFileHash(health.CredentialsDB)
+	tty := false
+	var b bytes.Buffer
+	c := runEpisodeStart(r.args("ep1"), &b, &b, spyEnv(&tty))
+	if c != exitUsage || tty || !strings.Contains(b.String(), "active action "+activeOperatorFixtureID+";") || serviceFileHash(health.CredentialsDB) != before {
+		t.Fatalf("start reached effect boundary: %d tty=%t %s", c, tty, b.String())
+	}
+	r.assertNoSessionFiles(t)
+	if _, err := os.Lstat(filepath.Join(r.f.ws, "ep1")); !errors.Is(err, os.ErrNotExist) {
+		t.Fatal("start created worktree")
+	}
+}
diff --git a/cmd/ailang-worldd/operator_runner.go b/cmd/ailang-worldd/operator_runner.go
index 75815a0..27771a2 100644
--- a/cmd/ailang-worldd/operator_runner.go
+++ b/cmd/ailang-worldd/operator_runner.go
@@ -95,19 +95,39 @@ func (e serviceEnv) runnerTick(c runnerConfig, h *store.OperatorHeartbeat) error
 	})
 }
 
-func (e serviceEnv) runnerStatus(w io.Writer) error {
-	c, _, err := e.installedRunner()
-	if err != nil {
-		return err
-	}
-	db, err := store.OpenCredentialsReadOnly(runnerDB(c.Credentials))
+func runnerHeartbeat(path string) (store.OperatorHeartbeat, error) {
+	db, err := store.OpenCredentialsReadOnly(runnerDB(path))
 	if err != nil {
-		return errors.New("credentials disconnected")
+		return store.OperatorHeartbeat{}, errors.New("credentials disconnected")
 	}
 	defer db.Close()
 	ctx, cancel := budgetContext(time.Second)
 	defer cancel()
-	h, err := db.OperatorHeartbeat(ctx)
+	return db.OperatorHeartbeat(ctx)
+}
+
+func (e serviceEnv) waitForRunner(c runnerConfig, previous int64) bool {
+	deadline := time.Now().Add(e.wait)
+	for {
+		h, err := runnerHeartbeat(c.Credentials)
+		if err == nil && h.Owner.Generation > previous &&
+			h.Available(c.Executables[e.paths.Binary], "start", e.terminal.now()) {
+			return true
+		}
+		delay := min(time.Until(deadline), 100*time.Millisecond)
+		if delay <= 0 {
+			return false
+		}
+		time.Sleep(delay)
+	}
+}
+
+func (e serviceEnv) runnerStatus(w io.Writer) error {
+	c, _, err := e.installedRunner()
+	if err != nil {
+		return err
+	}
+	h, err := runnerHeartbeat(c.Credentials)
 	if err != nil {
 		return errors.New("runner heartbeat unavailable")
 	}
diff --git a/cmd/ailang-worldd/operator_service.go b/cmd/ailang-worldd/operator_service.go
index 608670e..b1f3be7 100644
--- a/cmd/ailang-worldd/operator_service.go
+++ b/cmd/ailang-worldd/operator_service.go
@@ -431,13 +431,25 @@ func (e serviceEnv) changeRunner(verb string, args []string) error {
 	_, err = e.controller.Print(e.domain, runnerLabel)
 	switch {
 	case errors.Is(err, errServiceNotLoaded):
+		previous, err := runnerHeartbeat(c.Credentials)
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
index 2da9ed9..645b311 100644
--- a/cmd/ailang-worldd/operator_service_test.go
+++ b/cmd/ailang-worldd/operator_service_test.go
@@ -54,6 +54,10 @@ func operatorService(t *testing.T) *operatorServiceFixture {
 		health:    daemon.HealthResponse{Status: "ok", DBPath: world, CredentialsDB: creds, Build: daemon.BuildIdentity{ExecutableSHA256: r.BinarySHA256}}}
 	f.ctl = &fakeServiceController{runner: true, binary: r.Binary}
 	f.env.controller = f.ctl
+	f.ctl.onBoot = func() error {
+		o.heartbeat.Owner.Generation = 0
+		return f.env.runnerTick(o.config, &o.heartbeat)
+	}
 	f.env.client = &http.Client{Transport: serviceTransport(func(req *http.Request) (*http.Response, error) {
 		if req.URL.Host != r.Address || req.URL.Path != "/v1/health" {
 			t.Fatalf("unexpected health URL: %s", req.URL)
@@ -247,6 +251,53 @@ func TestOperatorServiceStartup(t *testing.T) {
 	}
 }
 
+func TestOperatorServiceStartupTimeout(t *testing.T) {
+	f := operatorService(t)
+	f.mustInstall(t)
+	if err := f.env.runnerTick(f.config, &f.heartbeat); err != nil {
+		t.Fatal(err)
+	}
+	before := serviceFileHash(f.config.Credentials)
+	f.ctl.onBoot = nil // no new lease, even though the old heartbeat remains fresh
+	if c, msg := f.command("start"); c != 1 || !strings.Contains(msg, "readiness timed out") ||
+		f.ctl.loaded || f.ctl.boots != 1 || f.ctl.stops != 1 || serviceFileHash(f.config.Credentials) != before {
+		t.Fatalf("timeout: %d %s", c, msg)
+	}
+}
+
+func TestOperatorServiceUpgradePairFence(t *testing.T) {
+	for _, tc := range []struct{ arm, want string }{
+		{"installed", "runner is still installed"},
+		{"registered", "cannot verify that the operator runner is unloaded"},
+		{"locked", "runner lifecycle is busy"},
+	} {
+		t.Run(tc.arm, func(t *testing.T) {
+			f := operatorService(t)
+			f.mustInstall(t)
+			if tc.arm != "installed" {
+				if c, msg := f.command("uninstall"); c != 0 {
+					t.Fatal(c, msg)
+				}
+			}
+			f.env.controller = f.homeJob
+			f.homeJob.runner = tc.arm == "registered"
+			if tc.arm == "locked" {
+				lock, err := runnerLifecycleLock(filepath.Dir(f.env.paths.Record))
+				if err != nil {
+					t.Fatal(err)
+				}
+				defer lock.Close()
+			}
+			before := serviceFileHash(f.env.paths.Binary)
+			f.writes = nil
+			if c, msg := f.run("upgrade", "--binary", f.configFile); c != 1 || !strings.Contains(msg, tc.want) ||
+				f.homeJob.kicks != 0 || serviceFileHash(f.env.paths.Binary) != before || len(f.writes) != 0 {
+				t.Fatalf("unpaired upgrade: %d %s writes=%v", c, msg, f.writes)
+			}
+		})
+	}
+}
+
 func TestOperatorRunnerKernelLock(t *testing.T) {
 	p := filepath.Join(t.TempDir(), "runner.lock")
 	f, err := runnerKernelLock(p)
@@ -291,7 +342,11 @@ func TestOperatorServiceWriteScope(t *testing.T) {
 	if info, err := os.Stat(filepath.Dir(p.Log)); err != nil || info.Mode().Perm() != 0700 {
 		t.Fatalf("spool: %v %v", info, err)
 	}
-	for _, verb := range []string{"start", "stop", "uninstall"} {
+	if c, s := f.command("start"); c != 0 {
+		t.Fatalf("start: %d %s", c, s)
+	}
+	want[4] = serviceFileHash(f.config.Credentials) // the allowed heartbeat transaction
+	for _, verb := range []string{"stop", "uninstall"} {
 		if c, s := f.command(verb); c != 0 {
 			t.Fatalf("%s: %d %s", verb, c, s)
 		}
diff --git a/cmd/ailang-worldd/service_test.go b/cmd/ailang-worldd/service_test.go
index 08f2711..63462d2 100644
--- a/cmd/ailang-worldd/service_test.go
+++ b/cmd/ailang-worldd/service_test.go
@@ -28,12 +28,16 @@ type fakeServiceController struct {
 	binary       string
 	running      string
 	onKick       func() error
+	onBoot       func() error
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
diff --git a/host/store/operator_runner.go b/host/store/operator_runner.go
index 4bebd72..60b881e 100644
--- a/host/store/operator_runner.go
+++ b/host/store/operator_runner.go
@@ -78,3 +78,16 @@ func (c *Credentials) CancelUnstartedOperatorAction(ctx context.Context, id stri
 			WHERE action_id=? AND kind<>'operation')`, now, id, now, id)
 	return operatorChangedOne(result, err)
 }
+
+// OperatorOwnedEpisode reads durable exclusion independently of heartbeat.
+func (c *Credentials) OperatorOwnedEpisode(ctx context.Context, episode string) (string, error) {
+	if err := c.operatorAccess(ctx); err != nil {
+		return "", err
+	}
+	var id string
+	err := c.db.QueryRowContext(ctx, `SELECT id FROM operator_actions WHERE episode_id=? AND owner IS NOT NULL`, episode).Scan(&id)
+	if errors.Is(err, sql.ErrNoRows) {
+		return "", nil
+	}
+	return id, err
+}
````

### Patch S3b

````diff S3b
diff --git a/host/store/operator_cas_test.go b/host/store/operator_cas_test.go
index 6d528ac..52048d2 100644
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
@@ -121,3 +123,61 @@ func TestOperatorClaimRace(t *testing.T) {
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
diff --git a/host/store/operator_steps.go b/host/store/operator_steps.go
new file mode 100644
index 0000000..3372945
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

### Patch S4c-replay

````diff S4c-replay
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
