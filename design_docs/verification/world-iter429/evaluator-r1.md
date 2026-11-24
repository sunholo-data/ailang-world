# Evaluator r1 - iter 429, row 243 V4 S4 (head 22473391, base f074dd49)

## Verdict: FAIL (one BLOCKING finding). Score 72/100 (would be a PASS on points; the blocking finding forces FAIL)

| Category | Score |
|---|---|
| Correctness vs design / plan section 9 | 24/30 |
| Tests (non-vacuity, mutants) | 18/30 |
| Security / fail-closed | 15/20 |
| Coding standards | 8/10 |
| Scope and size | 7/10 |
| Total | 72/100 |

## BLOCKING

B1. Install-time authority and fail-closed gates are undefended by the whole package, and the `service runner install|cancel` commands become reachable at merge (design 5.5 / section 8: "never cut an authority or recovery gate").
Four gates survived the WHOLE `./cmd/ailang-worldd/` package (rc=0, ledger lines J-CANCEL-CONFIRM, J-UPGRADE-RUNNER-PLIST, J-UPGRADE-SUPERVISOR):
- `cmd/ailang-worldd/operator_service.go` `runServiceRunner`: removing `|| verb == "cancel"` from the typed-y set lets `service runner cancel ACTION-ID` cancel an operator action with no confirmation. This is the attended-cancel authority gate and has no CLI test (grep of `_test` shows no cancel CLI test).
- `cmd/ailang-worldd/service_upgrade.go:~206`: removing the Plist half of the "runner still installed" check survives (Record alone suffices, so partly equivalent). Removing the supervisor check (`!errors.Is(err, errServiceNotLoaded)` -> `err == nil && false`) also survives the whole package. No test anywhere asserts the strings "operator runner is still installed", "lifecycle is busy" or "cannot verify that the operator runner is unloaded". So the S4-D3 / plan 9.3 claim "Companion install fences unpaired home upgrade" has NO behavioural test. The upgrade-time refusal and its lifecycle lock are shipped untested.
- The plan's S4b patch (`TestOperatorServiceUpgradePairFence`) that would pin this is parked, so S4 ships the fence without its test.
Reproduce (mutants.tsv has the exact edits): in the eval worktree, edit `service_upgrade.go` to replace `!errors.Is(err, errServiceNotLoaded) {\n\t\treturn refuse("cannot verify` with `err == nil && false {...`; `go test -count=1 ./cmd/ailang-worldd/` -> rc=0 (PASS = SURVIVED). Out files: `$HOME/.ailang/state/world-iter429/judge/out-J-*-whole.txt`.
Fix: add (a) a CLI cancel test asserting refusal on `n` and no row change, (b) an upgrade test with a fake installed runner record/plist and a loaded supervisor asserting refusal with nothing swapped. Both are small.

## Non-blocking findings (ranked)
1. Survived the cmd subset (Operator|HomeService|HomeVersion|DoctorService|Quickstart|Help|Keyguard|CredentialWriters), NOT run whole-package for time (the whole package survivors above prove the subset is a good proxy for these, but they are SURVIVED-in-subset only): J-OWN-UID (the ownership UID comparison; untestable as non-root, only the group/other-write and wrong-type checks are pinned), J-WRITABLE-OTHERS (the `0002` bit: only group-write 0020 is tested), J-LOCKFILE-MODE (kernel lock file mode/owner), J-LOCKFILE-NOFOLLOW, J-CONFIG-NOFOLLOW, J-PRESET-SE (reserved name `se-tools`), J-DISTINCT-STORES, J-HEALTH-STATUS, J-LOOPBACK-IP (second loopback check), J-ADDR-GIVEN (`operator-runner` with `--addr` is accepted). The plan's own table lists ~76 more named survivors; most are the same class.
2. Store survivors (whole `host/store` run, genuinely undefended): CANCEL of a `confirmed` row (`state IN ('requested','confirmed')` -> `('requested')` survives; the only success path is a `requested` row); `owner IS NULL` in cancel (plan calls it equivalent - probably unreachable state, still unpinned); `OperatorHeartbeat` read without the deadline check (`operatorAccess` removed survives; the write and cancel deadlines are pinned, the read is not). Plan's HB-SINGLETON-WHERE equivalent.
3. M1 compatible-rollback case is covered only by the pre-existing `TestHomeServiceUpgradeRollback` (kept green with schema fixture added); no new equal-schema named test. Behaviourally fine, but the pre-existing test now depends on the probe fixture returning the current schema for both builds.
4. Rollback refusal leaves the candidate serving/loaded and does not stop the service (exit 3, evidence retained). That matches design intent ("retain both"), but the home record keeps the old digest, so `service status` will report drift; remedy text is generic.
5. `version --json` `credentials_schema` is `omitempty` on the shared struct: a build lacking the field decodes as 0 and is refused (fail-closed, correct). It is self-reported by the binary, so a lying build is outside the threat model; a v9 build cannot be confused (field absent -> 0).
6. Size: plan appendix is huge (3,274 lines) but outside the measured unit. Measured unit (cmd host docs app): see below, within cap.

## Findings by question
- Rollback guard (M1): read at base and head. Refuses unknown candidate capability, lower previous capability, missing / non-regular / symlinked / different-bytes `.prev`; retains candidate and `.prev`; compatible case rolls back (pre-existing test). Mutants J-ROLLBACK-ZERO-CAND, J-ROLLBACK-RET3, J-M1-GUARD-REVERT killed by `TestOperatorServiceUpgrade` subtests (named below). I found no path that installs a binary that refuses the migrated store, except that the guard compares self-reported schema numbers and checks the `.prev` hash only against `from`.
- Install / ownership / environment: refusals are tested with a shared `refused()` helper asserting class substring, exit 1, zero observed writes and zero controller boots, so both class and unchanged state are asserted for the tested refusals. PATH (tool dirs + system floor), HOME, GIT/GH prompts, config hash and no ambient token are pinned (J-ENV-SYSPATH, J-ENV-GITPROMPT killed). UID ownership and several mode/nofollow checks are not (non-blocking 1).
- Writes outside the worktree: tests inject HOME (`runnerLogin`, `t.TempDir` + EvalSymlinks), a fake controller and fake terminal; `TestOperatorServiceWriteScope` asserts every observed write is under the fake HOME. Gate sweep `observation-before.txt` vs `observation-after.txt`: identical (diff empty). Pre-existing edits to help/keyguard/quickstart tests add `t.Setenv("HOME", t.TempDir())`, which shrinks real-HOME exposure.
- Heartbeat / liveness: heartbeat is liveness only in the diff: nothing releases ownership or reservations; `Available` is a pure reader. Freshness operators are pinned (seen_at<=now and > now-10: HB-NOW-EQUAL, HB-STALE-EQUAL, my J-ST-HB-GTE killed by `TestOperatorRunnerUnavailable/exactly_stale`).
- Cancel: unstarted `requested` row advances once; reserved / owned / expired are refused with the row unchanged (plan CANCEL-* mutants killed in whole `host/store`; my J-ST-CANCEL-EXPIRES-ARG and J-ST-CANCEL-STEPS-ALL killed by `TestOperatorCancelInflight`). The plan's reported survivors: CANCEL-OWNER survives (equivalent in reachable states) and I additionally found the `confirmed` state undefended. The HB-*-DEADLINE survivors reproduce for the read path (my J-ST-HB-ACCESS). The cancel predicate is mostly tested; deadline/owner/confirmed edges are not.
- Inertness / scope: gate sweep `scope.txt`: routes 50 (base 50), surface POSTs 3 (3); daemon.go, surface_page.go, app.js unchanged; no `.ail` changed; `UNTOUCHED_VIOLATIONS: none`. The internal `operator-runner` mode rejects any args (exitUsage) and `--addr`; it only heartbeats (no dispatch). Catalogue is written to the store, not served. `budgetContext` is used for every bounded op; no new context.Background/TODO/WithoutCancel in production added.
- Pre-existing test edits: `service_test.go` (adds `runner` bit to the fake controller and a `CredentialsSchema` field in the probe; no assertion weakened), `cli_credential_writers_test.go` (adds the new `runnerWrite` site to the exact-map; strengthens the inventory), `help_test.go` (adds runner verbs), `keyguard_test.go` and `quickstart_service_test.go` (HOME isolation, one added mutant). None weakened.
- Carried items: S3b, judge-428 gaps and different-artifact replay are moved to Appendix C with measured sizes; honest and explicit. The upgrade fence's test (S4b) is among the carried items while the fence itself ships: that is the BLOCKING gap above.
- Plan false premises (F1-F10): spot-checks consistent with the diff (e.g. F7 fake controller, F8 EvalSymlinks fix visible in `newServiceFixture`). Decisions S4-D1..D4 honoured; risk 8 "Uncertain owner" acknowledged.
- Intermediate-commit green: not individually executed (time); M1 files `service_upgrade.go`, `version.go`, `service_test.go` and `operator_upgrade_test.go` changed together in the head and the M1 revert mutant is killed by M1's test. State this as UNMEASURED per-commit.

## Size (git diff -M --numstat f074dd49 HEAD -- cmd host docs app)
Sum of insertions+deletions: cmd 1,0.. see table: cli_credential_writers_test 2, doctor_service 17, help_test 4, keyguard 5, keyguard_test 1, main 7, operator_runner 153, operator_service 451, operator_service_test 353, operator_upgrade_test 83, quickstart_service_test 3, service 20, service_test 14, service_upgrade 27, version 12, store/operator_runner 80, store test 144, QUICKSTART 24 = 1,400 (at the cap; the plan states 1,400).

## Mutants (ledger `$HOME/.ailang/state/world-iter429/judge/mutants.tsv`, 40 lines; per-run outputs `out-*.txt`)
cmd runs use the Operator|HomeService|... subset unless "whole". store runs are whole `host/store`.
| id | verdict |
|---|---|
| J-OWN-UID (first edit) | COMPILEERR (not counted); rerun SURVIVED subset |
| J-WRITABLE-GROUP (must-die) | KILLED by TestOperatorServiceOwnership/writable_tool_dir |
| J-WRITABLE-OTHERS | SURVIVED subset |
| J-LOCKFILE-NB | KILLED by TestOperatorRunnerKernelLock |
| J-LOCKFILE-MODE | SURVIVED subset |
| J-LOCKFILE-NOFOLLOW | SURVIVED subset |
| J-ENV-SYSPATH | KILLED by TestOperatorServiceEnvironment |
| J-ENV-GITPROMPT (1st anchor INVALID, rerun) | KILLED by TestOperatorServiceEnvironment |
| J-BACKEND | KILLED by TestOperatorServiceConfig |
| J-PRESET-SE | SURVIVED subset |
| J-LOOPBACK-IP (1st COMPILEERR, rerun) | SURVIVED subset |
| J-DISTINCT-STORES | SURVIVED subset |
| J-HEALTH-STATUS | SURVIVED subset |
| J-SCHEMA-EQ | KILLED by TestOperatorServiceStartup |
| J-UNINSTALL-LOG | KILLED by TestOperatorServiceWriteScope |
| J-CANCEL-CONFIRM | SURVIVED subset AND whole package |
| J-UNKNOWN-FIELDS | KILLED by TestOperatorServiceConfig |
| J-CONFIG-NOFOLLOW | SURVIVED subset |
| J-TYPECHECK | KILLED by TestOperatorServiceOwnership/symlink_tool |
| J-READINESS-BOTH | KILLED by TestOperatorStatusDisconnected |
| J-ADDR-GIVEN | SURVIVED subset |
| J-ROLLBACK-ZERO-CAND | KILLED by TestOperatorServiceUpgrade/missing_candidate_schema |
| J-ROLLBACK-RET3 | KILLED by TestOperatorServiceUpgrade (4+ subtests) |
| J-UPGRADE-RUNNER-PLIST | SURVIVED subset AND whole |
| J-UPGRADE-SUPERVISOR (= plan PAIR-UPGRADE-SUPERVISOR) | SURVIVED whole |
| J-STATUS-FINISH | KILLED by TestOperatorStatusDisconnected |
| J-STATUS-STALE-WINDOW | KILLED by TestOperatorStatusDisconnected |
| J-CTL-LIVE-CONTROL-LIVE (must-live: status read budget 1s->2s) | SURVIVED (as expected) |
| J-ST-HB-GTE | KILLED by TestOperatorRunnerUnavailable/exactly_stale |
| J-ST-CANCEL-EXPIRES-ARG | KILLED by TestOperatorCancelInflight |
| J-ST-CANCEL-OWNER-ISNULL | SURVIVED whole host/store |
| J-ST-CANCEL-STEPS-ALL | KILLED by TestOperatorCancelInflight |
| J-ST-CANCEL-CONFIRMED | SURVIVED whole host/store |
| J-ST-HB-ACCESS | SURVIVED whole host/store |
| J-M1-GUARD-REVERT (non-vacuity of M1 guard, whole guard disabled) | KILLED by TestOperatorServiceUpgrade (linked / missing candidate schema / missing previous binary ...) |
All kills are behavioural (runtime assertions); no kill relied on an AST/string inventory test. 33 distinct ids, none in the plan's id list (all prefixed J-; the one overlap in substance is J-UPGRADE-SUPERVISOR vs PAIR-UPGRADE-SUPERVISOR, which I re-measured on purpose). `comm -23` was not run as a file; the id sets are disjoint by prefix (plan ids never start with `J-`). Whole-package runs for survivors were limited to 3 cmd survivors plus all store mutants (time budget; subset-only survivors are flagged above).

## Plan section 4 re-measurement
Of ~76 plan named-survivors I re-measured by whole package the upgrade-fence group (INSTALLED/SUPERVISOR), and by the subset the ownership, config, loopback, lock-file, cancel-confirm groups: the groups survive outside the sandbox too, i.e. NM for those is genuinely undefended. The plan's host/store survivors reproduce (CANCEL-OWNER-style, HB deadline read).

## Test runs (after ALL_DONE.22473391d9b8...)
- Gate sweep (controller): all legs rc=0 including pkg cmd/ailang-worldd 237 s, `-race` cmd 346 s, host/store, `-race` host/store, verify_ail, fence, scope. `verify_go` rc=1 on untouched packages only (host/daemon TestSeToolsRunReplayIsByteEqual "close of closed channel" panic, host/verifygate TestMissionConfigDefaultFlowIntegration); the latter is on the known load-sensitive list; I did not re-run them alone (time). UNMEASURED whether the daemon failure is load-only.
- Mine: `go vet ./cmd/... ./host/...` rc=0; `go test -count=1 ./host/store/` ok; `go test -count=1 -run 'Operator|HomeService|HomeVersion|DoctorService|Quickstart|Help|Keyguard|CredentialWriters' ./cmd/ailang-worldd/` ok. `go test -race ./host/store/` and `go test -run '^$' ./...` not re-run (the controller's sweep covers race on host/store). Worktree clean (`git status --short` empty); throwaway worktrees removed.
