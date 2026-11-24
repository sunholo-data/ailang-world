# Evaluator report r1 — iteration 430, row 243 V4 S4b (head 40c21340, base c42fe698)

## Verdict: PASS — 85/100. No BLOCKING findings.

| Category | Pts | Score |
|---|---|---|
| Correctness vs design §2.2/§5.5 and plan §9 | 30 | 25 |
| Tests (non-vacuous, kills, portable) | 30 | 26 |
| Security / fail-closed | 20 | 16 |
| Coding standards | 10 | 9 |
| Scope and size | 10 | 9 |

## Non-blocking findings (ranked)
1. **Two terminal entry points still reach `provisionEpisode` (worktree add + session mint) without the episode lock**: `runSessionNew` (cmd/ailang-worldd/session_new.go:649) and `runApprove` (approve.go:398, the mint-request approval path). Only start, preset start, resume, launch and finish take it (episode.go:364, episode_preset.go:282, episode_resume.go:175, episode_launch.go:289/456, episode_finish.go:128). Design §2.2 names "start/launch/finish", so this is within the letter of S4b, and it is latent because S4b dispatches no effect and nothing can own an episode in production yet. But plan Appendix C does not list it, and S5/S6 must close it before any consumer exists. Repro: `grep -n terminalEpisodeLock cmd/ailang-worldd/session_new.go cmd/ailang-worldd/approve.go` prints nothing. Recommend adding it to Appendix C as a gate.
2. **Operator-trap remedy is thin**: the refusal reads `...active action <id>; run ailang-worldd service runner status and reconcile at the terminal` (operator_episode_lock.go:56). No attended cancel or resolver verb exists in cmd (grep "cancel|reconcil" finds none outside that string). It is latent now, since no claim path is reachable (zero production callers of Claim/Reserve/Observe, measured). The plan defers "M4 attended cancel" in Appendix C. A durable owner correctly never expires by staleness (T-OWNER-BLOCK and the resume test), so once dispatch lands an uncertain action blocks the terminal until that resolver exists. Make the resolver a hard precondition of S5.
3. **Permanent sidecar files beside the REAL credentials store**: every start/preset/resume/finish creates `<creds>.operator-episode-<ep>.lock` (mode 0600, never removed, one per episode). The plan documents this (§1.3, F5), but docs/QUICKSTART.md and the user docs do not. Tests only use temp dirs (the gate observation-before/after files are identical). Acceptable, but should be documented and ideally garbage-collected.
4. **Unpinned survivors (whole-package, my measurements)**: `Z-ABSENT-ANY` (operator_episode_lock.go:36, dropping `flags == 0` from `absent`) means that if the credentials directory does not exist, a creating acquire returns (nil, nil) instead of the open error, so no lock is held. `Z-WAIT-READBUDGET` (operator_runner.go:112, `min(remaining, time.Second)` → `time.Second`) is only observable in real time. `Z-RDWR` is an equivalent mutant, since flock works on O_RDONLY. Low severity.
5. `TestOperatorServiceWriteScope` edit (operator_service_test.go): `start` was moved out of the loop and `want[4]` re-baselined to the credentials hash after start, because the fake boot now writes a heartbeat. The assertion is adapted, not removed (the hash is still pinned across stop/uninstall), but "start leaves credentials unchanged" is no longer asserted there. Covered by the startup gate's `serviceFileHash(creds) != knownHash` arm.
6. Linux portability is by code reading only (I ran macOS only). I found no macOS-specific assertion. Symlink tests use `os.Stat`/`os.Readlink`, never `Lstat().Mode()` (the iteration-429 trap). Modes are checked with `Perm()` on regular files. The only symlink-to-directory use is the canonical alias (nonwaiting). Flock tests use independent open descriptions, which behave the same on Linux. Residual risk: `TestOperatorEpisodeLockRefusals` `unsafe lock` assumes the process is not root-squashing the chmod, which is fine on ubuntu runners.

## Answers to the questions
- **Completeness**: five enumerated entry points are guarded, each BEFORE its first effect and non-waiting (`LOCK_EX|LOCK_NB`, then a bounded 4 s read-only owner query with busy_timeout 0). Plain start is guarded after daemon checks and before `--request`/`provisionEpisode`. Preset start is guarded before `afterToken` and `provisionEpisode`, and passes the open file to nested `launchEpisode` (name-checked transfer, `held.Name()!=path` refuses). Resume is guarded before the preset-path/token work. Standalone launch probes with flags 0 (no sidecar created on legacy refusal), re-acquires at the PID boundary, and validates the served credentials pair first. Finish is guarded after F1 and health, before report/dispatch/check. Exact action id appears in the message (`expectEpisodeRefusal` asserts `active action oa-… ;`) with class `ErrOperatorConflict` and exit 2. A durable owner still refuses after the holder fd is closed and the clock is past expiry (resume subtest; the query reads `owner IS NOT NULL`, never the lease). Other episodes and other stores acquire; symlinked spellings canonicalize. Missed paths: finding 1.
- **Safety**: every acquire has a `defer Close()` or a documented transfer. I mutated each (Z-NO-DEFER-*): launch, preset, finish, start and resume leaks are all killed by whole-package tests (finish by `TestEpisodeFinishRerunAC10`, `…FinishedRerun`; start/resume by `TestEpisodeResumeChainUsesCurrentCredential` and `TestEpisodeStartRefusesAMismatchedDaemonBeforeTheTerminal`). The combined runs were not attributed per edit for start vs resume. Absent or pre-v10 store: `ReadOperatorEpisodeOwner` returns "" (store test arms: absent/v9/v10-empty/future/corrupt), so "no exclusion" is justified because no operator rows can exist.
- **Startup gate**: bounded by injected `runnerNow`/`runnerSleep` (no real sleep in tests; sleeps asserted ≤3, each ≤100 ms, clock lands exactly on `wait`). A timeout runs Bootout and leaves the lease untouched (hash and generation/seen_at asserted). Generation, build, protocol, readiness and freshness are each pinned by separate arms (same generation, stale, future, wrong build, wrong protocol, finish only, empty owner, bad read).
- **CI cost**: all new tests together run in 3.0 s test time (9 s wall including compile). No real sleeps (`time.AfterFunc` fuses are 2 s safety timers that are stopped). Focused cmd selection 47.9 s, race selection 22.7 s.

## Gates and runs (mine)
- `go test -count=1 ./host/store/` ok 12.3 s.
- Focused cmd selection (directive pattern) ok 47.9 s.
- `go test -race -count=1 -run 'Operator|Launch' ./cmd/ailang-worldd/` ok 22.7 s.
- `go vet ./cmd/... ./host/...` rc 0.
- `go test -run '^$' ./...` clean.
- Whole `./cmd/ailang-worldd` was run only under mutation (the unmutated whole-package pass is the controller's: `pkg-cmd-ailang-worldd rc=0 235s`, race rc=0).
- Controller sweep: legs.txt all rc=0 except `verify_go` rc=1, whose only failures are the known load-sensitive `TestSeToolsRunBenchmarkTasksEndToEnd` (host/daemon) and `TestMissionConfigDefaultFlowIntegration` (host/verifygate), in untouched packages. `observation-before.txt` and `observation-after.txt` are identical (real LaunchAgents, `~/.ailang-world`, launchctl count untouched).

## Size (my measurement, unit `git diff -M --numstat … -- cmd host docs app`)
M1 638, M2 878, M3 978 cumulative (matches the plan). Total 978 ≤ 1,400. Increments are 638/240/100. Files changed match plan §5: cmd 15, host/store 4, plus the plan and JSON. No `.ail`, `app.js`, `surfaceCSP`, `surface_page.go` or `daemon.go` change. No new route. Zero production callers of `ClaimOperatorAction|ReserveOperatorStep|ObserveOperatorStep` outside the store. Zero added `context.Background/TODO/WithoutCancel` lines (the deadline helper `budgetContext` pre-exists).

## Mutation ledger
`$HOME/.ailang/state/world-iter430/judge/mutants.tsv`, 70 lines (every run, including re-runs and compile errors; combined whole-package runs are labelled COMBINED). All edits were restored; `git status --short` is clean in my worktree.

### (a) Plan §4 rows re-measured (27 of the UNMEASURED rows, edits taken from the plan ledger, applied by me to HEAD): ALL KILLED by behavioural tests (each runs real CLI entry points, the real migration or signer fixtures, or a fake supervisor)
T-PLAIN, T-PRESET, T-RESUME, T-LAUNCH, T-FINISH → `TestOperatorTerminalExclusion/{plain,preset,resume,launch,finish}`. T-LAUNCH-LATE, T-LAUNCH-LATE-ERR → `TestOperatorLaunchClaimDuringPreflight`. T-LAUNCH-NESTED → `TestStartLaunchHeadlessMintsAndLaunchesOnce`. T-CANONICAL, T-HELD-KEY, T-NB, T-OWNER-BLOCK → `TestOperatorTerminalExclusion/{nonwaiting,nonwaiting,nonwaiting,resume}`. S-POST-TIME → `…StartupGate/late`. S-GENERATION, S-SLEEP-CAP → `…/same_generation`. S-AVAILABLE → `…/finish_only`. S-LOADED → `…/already_loaded_stale`. P-PRINT-ERROR → `TestOperatorServiceUpgradePairFence/inspect_error`. J-WRITABLE-OTHERS → `TestOperatorServiceOwnership/other_writable`. J-LOCKFILE-MODE → `TestOperatorRunnerLockFileRefusals/mode`. J-LOCKFILE-NOFOLLOW → `…/symlink`. J-CONFIG-NOFOLLOW → `TestOperatorRunnerSetupRefusals/linked_config`. J-PRESET-SE → `…/reserved_preset`. J-DISTINCT-STORES → `TestOperatorServicePair/distinct`. J-HEALTH-STATUS → `…/health_status`. J-LOOPBACK-IP → `…/named_loopback`. J-ADDR-GIVEN → `TestOperatorRunnerSetupRefusals/global_address`.
The `OperatorHeartbeat` read deadline is pinned: Z-HB-DEADLINE is killed by `TestCredentialsMethodsRequireADeadline`.
None of the 27 was killed only by an AST/static/inventory test, so no re-run with a test excluded was needed.

### (b) My own mutants (ids absent from the plan: `comm -23 mine plan` → 38 ids; the files are `judge/mine-ids.txt` and `judge/plan-ids.txt`)
| id | edit | verdict / killer |
|---|---|---|
| Z-CTRL-DIE (must-die) | `case id != ""` → `case false` | KILLED `TestOperatorTerminalExclusion/{plain,preset,…}` |
| Z-CTRL-LIVE (must-live) | comment text only | SURVIVED (named); not rerun on the whole package, since it is a comment |
| Z-NOCREATE | creating acquire passes flags 0 | KILLED `TestLaunchHeadlessClaudeAC9AndSecondLaunch`, `…TerminalExclusion/nonwaiting` |
| Z-TIMEOUT-NS(-excl-only) | 4 s → 4 ns read budget | KILLED `TestOperatorTerminalExclusion/{plain,preset}` (also `TestLaunchInteractiveWallCapIsSessionLeft`) |
| Z-RAW-LOCKERR (compile error), -fixed (wrong mutant), -fixed2 | contention returns the raw errno | KILLED (fixed2) `TestOperatorTerminalExclusion/{nonwaiting,plain}` |
| Z-NO-CONFLICT-WRAP | `%w` → `%v` | KILLED `…/nonwaiting` |
| Z-HELD-RETURN | transfer returns the held fd | KILLED `…/nonwaiting` |
| Z-LAUNCH-PREFLIGHT-CREATE | preflight creates a sidecar | KILLED `TestLaunchInteractiveRefusals/*` |
| Z-NO-DEFER-LAUNCH | drop launch close | KILLED `TestLaunchHeadlessClaudeAC9AndSecondLaunch` |
| Z-NO-DEFER-PRESET | drop preset close | KILLED `TestOperatorTerminalExclusion/resume` |
| Z-NO-DEFER-FINISH / -START / -RESUME | drop the close | SURVIVED named; KILLED whole-package (combined) by `TestEpisodeFinishRerunAC10` etc. and `TestEpisodeResumeChainUsesCurrentCredential` / `TestEpisodeStartRefusesAMismatched…` |
| Z-WAIT-DOUBLE | deadline = 2×wait | KILLED `…Startup/gate/same_generation` |
| Z-WAIT-SLEEP-1S | sleep cap 100 ms → 1 s | KILLED `…/same_generation` |
| Z-PREV-ZERO (compile error), -fixed | previous generation ignored | KILLED `…/same_generation` |
| Z-LOADED-NOWAIT | skip the loaded-job check | KILLED `…/already_loaded_stale` |
| Z-NO-BOOTOUT (wrong site, INVALID for the intent), -timeout | no stop on timeout | KILLED `…/same_generation` |
| Z-HB-DEADLINE | heartbeat read without a deadline | KILLED `TestCredentialsMethodsRequireADeadline` |
| Z-READ-BUSY | busy_timeout 0 → 5000 | KILLED `TestOperatorEpisodeOwnerReadSchema` |
| Z-VERSION-LE | v10 treated as legacy | SURVIVED the store package alone (expected; the CLI resume test kills it, as the plan states for E-CURRENT) |
| Z-ABSENT-ANY, Z-RDWR, Z-WAIT-READBUDGET | see finding 4 | SURVIVED whole package (3-edit combined run, rc=0, `whole-combined3.txt`) |

Counts: 27 plan rows re-measured (27 killed); my own mutants 24 distinct edits, of which 2 were compile errors that I corrected and re-ran, and 1 wrong site that I re-ran; 3 survivors on the whole package (all low severity, finding 4); 1 control must-live survived as intended; 1 must-die killed. Every KILLED entry above is a behavioural test (real CLI entry points, real migrated store, fake supervisor, fake clock).

## Per-milestone non-vacuity
M1's production hunks cannot be reverted wholesale against the later tests (compile break), so I measured at the call sites instead. Each M1 acceptance test fails with its own call site removed (T-PLAIN, T-PRESET, T-RESUME, T-LAUNCH, T-FINISH, T-LAUNCH-*; `TestOperatorEpisodeLockExclusive` via T-EX is the plan's row, not re-run by me). M2: `…Startup/gate/*` fail under Z-WAIT-*, Z-LOADED-NOWAIT, Z-NO-BOOTOUT-timeout, S-*. M3: J-* pins fail under their original mutations (all nine plus the heartbeat deadline).

## Plan false premises, decisions and §8 risks
F1–F8 hold up as far as I checked (F3 and F5 are confirmed in the code: late acquisition plus a no-create preflight; F6 confirmed, no production Claim/Reserve/Observe caller). §1.3 decisions are all realized in the code. The §8 risks are real and honestly stated: sidecar inode stability, the stale heartbeat cannot prove a dead child, and the runner drain gate (Appendix C). I add findings 1–3 as additions.

## Carried items
S4c-replay (13 lines), S3b (132 lines) and J-OWN-UID are dispositioned honestly in Appendix C (patch sizes measured, with gates attached; J-OWN-UID correctly noted as needing a non-root fixture). I did not re-measure the S4c-replay or S3b patches.

## Pre-existing test edits audit
`episode_launch_test.go` (fixture health now reports the served credentials DB; GC disabled for the descriptor-leak observation; new reacquire assert), `operator_service_test.go` (injected clock, new arms, finding 5), `service_test.go` (fake controller gains onBoot, stopErr and runnerPrintErr; the semantics of `Bootout` are preserved on the default path), `credentials_test.go` (+3 deadline arms), `operator_schema_test.go` (+66, new test). No assertion is weakened apart from the re-baseline in finding 5.
