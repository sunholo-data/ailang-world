# Iteration 441 — controller readings (row 243 V4 S5f3)

Fire `world-1795710251-25803`, attempt 1. Controller `claude-opus-5-5` (session). Base `origin/dev` `2f573488` (= S5f2 merge `150d0d96` + record #686).

## Gates 0–2
- Gate 0: kill switch absent; `gh` = `sunholo-voight-kampff`; billing CLEAN; `mission_directives.sh` on #635: 0 allowlisted directives (21 comments; LEGACY watermark file); step 6a: 0 crash and 0 transcript notices on #635/#574 (control 107:4); `mission-world` inbox 0 unread.
- Gate 1: local `dev` was 9 commits behind `origin/dev`, so the charter and log were read from origin. Base `2f573488…@2026-11-26T16:26:44Z`. Skill drift: the resolved symlink and the pin copy were each byte-identical to ailang `origin/dev` on all 14 files. Dev CI at `2f573488` was in progress at Gate 1 (verify gate green); parent `150d0d96` 3/3 green.
- Gate 2: clause map 1.0 clauses 1–7 MET. Ledger 1 OPEN: D-WORLD-121 (row 267 before row 243?), default B. Pick = row 243 V4 S5f3, record 440's Next. No S5f3 plan, branch, PR or worktree existed (`git log origin/dev --grep S5f3`, `git branch -a`, `git worktree list`, open PRs). Controller census at HEAD: four start-path DECLARATIONS, zero callers; zero process-pin test names (positive control `operator_upgrade_test.go:13`).

## Routing
- Resolver, followed verbatim: planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin` (derive-planner-lane agrees), executor `recipe codex:gpt-6.2-sol declared:provider-pin`, evaluator `agent-tool sonnet declared:alias-pin`. Designer not spawned (reviewed design r1, quorum r2 PROCEED 3/3).
- The operator's standing request asked for every role through the Agent tool "exactly as the routing table specifies". The table returned `recipe` for the planner and executor, so the Agent tool is not a valid path for them, and they ran as `codex exec` recipes. The REQUIRED judge ran on the Agent tool, both rounds plus a round-3 addendum. No role failed to spawn.
- Codex probe rc 0.

## Baseline (controller, pristine `2f573488`, outside the sandbox, `-timeout 20m`)
- `base/` (co-loaded with the planner): `store-plain` 0, `store-race` 0, `worktree-plain` 0, `worktree-race` 0 fails; `cmd-plain` 12 (the 11 row-271 reds + `TestHostsLaunchersExecHost`, `codex: signal: killed` under load, passes alone at base); `daemon-race` TIMED OUT at 20 m (instrument failure).
- `base2/` (after the planner): `authority-plain` 0, `authority-race` 0; `daemon-race` TIMED OUT again at 20 m (1,235 s) — no valid local daemon base this iteration; CI arbitrated the daemon (no `host/daemon` diff).

## Planner
- `codex exec`, workspace-write, 17:34–18:52 (~78 min), **599,391 tok**, rc 0. Plan + JSON committed `ce74776b`.
- **Outcome B.** The compiled activation prototype `st-M4` measured **2,302** cumulative (> 1,650). S5f3 = M1 289, M2 686, M3 **1,385** cumulative. Appendix C: exactly one remaining sprint, **S5f4 = the caller, 917** measured on M3.
- §1.4 premises found false: the upgrader refused installed companions; two durable gaps had no working remedy; a digest alone did not prove executable permission/confinement; revoke had independent writer entries; an explicit GC DB path did not prove workspace binding.
- Controller round trip `ctl_roundtrip.sh M1 M2 M3`: `diff -r` rc 0 on `cmd`, `host`, `docs` for each.

## Executor (round 1)
- `codex exec` (~4 min, **73,830 tok**), rc 0. M1 `ACCEPT_OK` (15 PASS), M2 `ACCEPT_OK` (21 PASS), snapshots 4/0/4 and 8/0/8 cmp-equal to the stages.
- M3 applied cleanly; in-sandbox `accept.sh M3 fast` exited 1 on `TestOperatorServiceStartup/gate/{bad_read,already_loaded_stale}` (`context deadline exceeded`); executor stopped correctly. Controller outside the sandbox: `TestOperatorServiceStartup -count=3` ok at base and head; `accept.sh M3 fast` rc 0, `ACCEPT_OK`, 60 PASS / 0 FAIL. The controller took the M3 snapshot with the plan's `snapshot.sh` (18/0/18, cmp-equal to `st-M3`). False-green (3): the sandbox invented the red.
- `commit_ms.sh` attempt 1 refused M1 on EXISTING `TestSessionRevokeByEpisode` (setup `session new` → `mint session: context deadline exceeded`, 6.15 s). `-count=5`: base 0.10–3.53 s, M1 0.08–0.98 s, all PASS. Retry committed M1 `7b7acdf6`, M2 `b56285dc`, M3 `211524ae` (NUMSTAT 289 → 686 → 1,385, exactly §1.5). Head tree = M3 manifest 18/18.

## Controller gates on `211524ae`
- Clean: gofmt (touched), build, vet, Linux vet, test compile, five inventories, `verify_ail`, GC probe, caller census (four definitions; zero non-declaration production hits), `accept.sh M3 race`.
- `host/store` plain/race, `host/authority` plain/race: 0 fails.
- `cmd-plain` added `TestSessionRootBudgets` (deterministic; the judge's B1) plus four git-fixture tests that pass alone (load).
- `daemon-race` head 6 fails vs a TRUNCATED base: no attribution (no diff).
- NOT run locally: `host/verifygate` — the CI red below came from it (rule 3g).

## Judge round 1
- Agent `sonnet`, foreground, ~87 min, **149,830 tok**: **FAIL 71** (22/15/17/8/9). B1: M1 renamed the revoke root, so `TestSessionRootBudgets` (a source-text guard) went red; reproduced first-party (PASS at base, FAIL at head).
- Non-blocking: refused/failed upgrade left a loaded companion stopped silently; `service runner stop|uninstall` refuses without a service record; process census matches on name only; 8 cmd survivors + 2 whole-store survivors.
- 25 own mutants (15 killed, 8 survived, 2 controls); 8/8 planner kills re-confirmed.

## Round 2
- Codex (~13 min, **157,191 tok**): R1 restores `revokeCtx, cancelRevoke` in `runSessionRevoke` (guard unchanged); R2 `upgradeFailureNext` adds the stopped-companion remedy; R3 pins killing J-ROLLBACK-SCHEMA, J-RESTART-LOADED, J-PROMPT-HASHES, J-POST-CENSUS, J-DISPOSE-OWNER, J-RECOVERY-SINGLE-STEP. 169 lines (37 production). Commit `81789324`, NUMSTAT 1,528.
- Controller round-2 gates: all legs rc 0 except `compare-cmd-plain` (head 12 = base 12 with `TestHostsLaunchersExecHost` out and `TestOperatorWorktreeWindow` in; the latter passes alone at base 7.66 s and head 7.67 s; untouched files).
- **Judge round 2** (fresh Agent `sonnet`, ~9 min, **63,389 tok**): **PASS 83** (27/25/18/8/5), 0 blocking. B1 and N1 CLOSED (guard re-detects the round-1 form); all six ids killed by named tests; own mutants 4/5 killed (O-DROP-LOADED survives: an already-stopped companion's failure message is untested).

## Round 3 (CI-found)
- PR #687 CI on `81789324`: both Go jobs red on one test only — `host/verifygate TestSubprocessSinkGate`: M3's `operator_pair_upgrade_test.go:81` used a bare `bytes.Buffer` as a subprocess stderr sink. PASS at base, FAIL at head locally.
- Codex (test-only, 2 lines): `var stderr syncBuffer`. Commit `ed135376`, NUMSTAT 1,527. Round-2 judge resumed by name (**66,062 tok**): addendum confirms the diff, gate pass and race pass; PASS 83 stands.

## Gate 3b
- PR #687 CI on `ed135376` attempt 1: build+test and verify green; race job red ONLY on the per-package 480 s bound — `host/daemon` 480.134 s (no diff), `cmd/ailang-worldd` 478.231 s, `host/verifygate` 469.105 s. Previous head `81789324`'s race job: `cmd/ailang-worldd` 402.9 s, `host/daemon` 448.3 s. Fifth instance of the row-267 race-bound class (D-WORLD-121 open). Re-ran the failed job.
- Attempt 2 on `ed135376`: 3/3 (`cmd/ailang-worldd` 275.1 s, `host/daemon` 358.2 s, `host/verifygate` 328.2 s). `MERGEABLE CLEAN`; merged `--match-head-commit` → `a50b0d72`.
- Merge `a50b0d72` attempt 1: race job red only on `FAIL host/daemon 480.095s` (`panic: test timed out after 8m0s`; `cmd/ailang-worldd` 443.934 s, `host/verifygate` 394.219 s) — the sixth instance. Attempt 2: **3/3** (`cmd/ailang-worldd` 273.859 s, `host/daemon` 358.805 s). **LANDED.**
- Round-3 executor: 68,031 tok.
