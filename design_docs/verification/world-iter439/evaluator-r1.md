# Evaluator r1 — iteration 439, row 243 V4 S5f (head 0ed2ab53, base 9610f5de)

## Verdict: PASS — 86 / 100. No BLOCKING findings.

| Category | Score |
|---|---|
| Correctness (30) | 27 |
| Tests (30) | 26 |
| Security / fail-closed (20) | 18 |
| Coding standards (10) | 8 |
| Scope and size (10) | 7 |
| **Total** | **86** |

## Measurements (mine)
- `git diff -M --numstat 9610f5de 0ed2ab53 -- cmd host docs app` = **578** (matches plan).
- `go test -count=1 ./host/store/` ok 12.8 s; `-race -run 'Operator|EpisodeHasUnresolved|NextOperatorStart' ./host/store/` ok 7.5 s; same selector `-race ./cmd/ailang-worldd/` ok 127 s (139 s wall; this is the whole `Operator` family, not new cost).
- New/changed cmd tests alone under `-race` (StartCandidate, MintRecoveryLockedRead, WaitReadBudget, StartFacts, StartChanged, TerminalCreateLockMissingParent): **9.3 s test / 11 s wall**, within the 12 s budget. New store tests under -race 2.6 s.
- `go vet ./cmd/... ./host/...` ok; `GOOS=linux go vet ./cmd/ailang-worldd ./host/store` ok; `go test -run '^$' ./...` clean.
- Controller gates (`gates-head/legs.txt`): store-plain/store-race 0 fails; cmd-plain fails=11 equal to base set (compare-cmd-plain rc=0); daemon-race head has 3 fails vs base 4 (the load-sensitive `TestSeToolsRunBenchmarkTasksEndToEnd` passed at head), so `compare-daemon-race rc=1` is a base-superset diff, not a regression. `gofmt` leg flags only `host/daemon/daemon.go`, which this sprint does not touch (`git diff 9610f5de HEAD -- host/daemon/daemon.go` is empty): pre-existing. `callers` leg rc=1 is a harness bug (`callers.sh: line 8: $4: unbound variable`); I redid the census myself (below).

## Questions
**M1/M2 kills (my re-run, original edits from mutants.json):**
- S9 (`NOT IN ('operation','files')`) KILLED by `TestEpisodeHasUnresolvedOperatorSteps/files_reserved` and `/files_observed_running`, host/store.
- S10 (`'mint'`) KILLED by `.../mint_reserved`, `.../mint_observed_uncertain`, host/store. Both are behavioural (real migration + reserve API). I ran the named test; whole host/store passes unmutated (12.8 s) and these are failing tests, so a whole run adds nothing to a kill.
- E-RES-EPISODE-CONFLICT KILLED by `TestOperatorMintRecoveryLockedRead/different_episode`; E-RES-REREAD KILLED by `.../expired_on_locked_read`, cmd/ailang-worldd. The test interleaves a real ledger rewrite inside the injected lock function and re-reads the receipt, then asserts store/file bytes unchanged and the lock released.

**M2 production semantics.** Production dispatch passes `openRunnerLock`, the function previously called inline with the same `(path, 0)` arguments; defer/close/readback and refusal strings unchanged. Mutants J13 (flags 0 -> O_CREATE), J14 (wrong episode key), J15 (dispatch wiring wrapped with O_CREATE) are all KILLED, so the wiring is pinned in the real `runServiceRunner` path (`TestOperatorMintRecovery/missing_lock`). No production change.

**M3 correctness.**
- `NextOperatorStart` SQL: kind/state/created `<=`/expires `>`/ORDER BY created_at,id are each pinned (J18, J20, D-SELECT-EXPIRES killed by `TestNextOperatorStart`). `LIMIT 1` is equivalent: `QueryRowContext.Scan` reads only the first row; J21 (LIMIT 2) SURVIVED in named and whole host/store (expected equivalent). `kind='start'` -> `kind<>'finish'` (J19) survives whole host/store: equivalent because the schema CHECK closes the domain to {start, finish} (control, not a gap).
- Visit: the same open fd is handed to the callback and exact kernel exclusion is asserted with a real child terminal holder and a symlinked credentials spelling; lock released via `defer` on every path (J05 note below); equality check on the whole row (`reflect.DeepEqual`), second read uses a fresh `now()` (J08 killed), lock-held returns `ErrOperatorConflict` wrapping EWOULDBLOCK (J09 killed), invalid episode refused before locking (J10 killed). Read-only DB handle leaves the SQLite writer free (asserted).
- `launch_binary_sha256`: adds a signed fact; changes the canonical start envelope for any newly signed start (resolution slice grows from 13 to 14 entries). Any start signed before this change would no longer reconstruct (binding refusal on the extra/missing fact). Acceptable here: nothing has a production caller of start (census below) and no persisted golden fixture is involved (the shifted `Resolution[...]` indices in `operator_start_input_test.go` are in-test). Digest source (J11 hash of path instead of bytes) KILLED by `TestOperatorStartFacts` + `TestOperatorStartChanged/launch_binary_bytes`; fact order (J22) KILLED, though only by wide-ranging file-publication tests, which is incidental but behavioural.
- Catalogue generation (LIVE heartbeat path). SQLite evaluates every expression on the RHS of `UPDATE ... SET` against the pre-update row, so `readiness IS ?`/`catalogue IS ?` compare the incoming values to the stored (old) ones; bind order verified by reading: `Readiness, SeenAt, Readiness, Catalogue, Catalogue, owner, generation, protocol, build, seenAt` matches the placeholders. Readers of `catalogue_generation` today: `ReadOperatorHeartbeat` (op runner status) and the CAS terms in `host/store/operator_mint.go:120` and `operator_worktree.go:53` (`catalogue_generation = ?`). Nothing depends on it advancing every tick; before this change every tick invalidated a signed mint/worktree CAS (premise in plan §1.4 is true, base SQL visibly did `+1` unconditionally). New behaviour only advances on a changed readiness or catalogue bytes, and the new-owner upsert path still increments. Residual: `catalogue` is a nullable BLOB; with a NULL on both sides `IS` treats it as equal (intended). Mutant J01 (`catalogue IS ?` -> `catalogue = ?`, which advances on NULL==NULL) SURVIVES whole host/store: no test covers an unchanged NULL catalogue. Non-blocking.

**Inertness.** Non-test production references (git grep, positive control = test callers exist): `authority.MintOperatorStart` 0 callers (definition only; test callers at `operator_mint_test.go:82`, `operator_files_test.go:50`); `NextOperatorStart` called only from `withOperatorStartCandidate`; `withOperatorStartCandidate` 0 callers; `openCanonicalTerminalEpisodeLock` called only from `openTerminalEpisodeLock`. Nothing new is reachable from `runOperatorRunner`/`runnerTick` besides the heartbeat SQL. No route, schema, `.ail` or served-behaviour change (diff limited to the listed files).

**Scope / deferral.** Priority order M1 > M2 > M3-seam was followed. The deferral rests on §1.4: I verified the catalogue premise directly; the dormant-helper premise I confirmed (`operator_worktree.go:40` calls `openRunnerLock` on the canonical episode path itself, so a second description would conflict); the upgrade-refusal premise is cited by line in the plan and not independently disproved by me. A caller without the four process pins is not a safe unit, so deferring is defensible and consistent with the directive's fallback. Weak points: (a) Appendix C says the prototype "builds and vets" but states no measured line count; the 1,082 cumulative cmd+host lines vs base is only in the controller's measurement. I re-measured `proto-S5f2-current` with per-file `git diff --no-index` against the base tree: **1082** (cmd host). So S5f2 at 1,082 on top of 578 is > 1,400 cumulative, consistent with a separate sprint, but the doc itself should carry the number. (b) Plan §1.5 race receipts, §4 mutation table, §6/§7 are "pending" (planner refusal); I supplied them here. (c) J-OWN-UID stays carried: Appendix C records the original ledger as unavailable.

## Mutation (mine, ledger `~/.ailang/state/world-iter439/judge/ledger-judge.jsonl`, 21 lines; J06, J17 first NOCOMPILE and were corrected and re-run; raw `J*.raw`)
Ids J01-J22 are not in `mutants.json` (ids D-*/S*/E-*/Z-*). Selector `Operator|EpisodeHasUnresolved|NextOperatorStart`.

| id | edit | result | killer (package) |
|---|---|---|---|
| J01 | `catalogue IS ?` -> `catalogue = ?` | SURVIVED (named + whole host/store) | none; NULL-catalogue case untested |
| J02 | `readiness IS ?` -> `IS NOT ?` (must-die) | KILLED | TestOperatorCatalogueAdmission (host/store) |
| J03 | `ELSE 1` -> `ELSE 2` | KILLED | TestOperatorCatalogueAdmission |
| J04 | `catalogue=?` -> `catalogue=catalogue` | KILLED | TestOperatorCatalogueAdmission |
| J05 | drop `defer db.Close()` | SURVIVED (named + `Operator` family; no whole cmd run) | fd/handle leak only; not behavioural |
| J06 | DeepEqual(row) -> DeepEqual(ID) | KILLED | TestOperatorStartCandidateRefusals/changed_generation (cmd) |
| J07 | `visit(lock, initial)` (must-live) | SURVIVED | equivalent after the DeepEqual gate |
| J08 | second read uses `initial.CreatedAt` as the clock | KILLED | .../expired_under_lock |
| J09 | lock-unavailable class -> ErrOperatorInput | KILLED | TestOperatorStartCandidateProcess |
| J10 | invalid episode returns nil | KILLED | .../invalid_episode |
| J11 | digest = hash of path string | KILLED | TestOperatorStartFacts, TestOperatorStartChanged/launch_binary_bytes |
| J13 | mint lock flags 0 -> O_CREATE | KILLED | TestOperatorMintRecovery/missing_lock, TestOperatorMintRecoveryLockedRead |
| J14 | mint lock path keyed by `id` | KILLED | TestOperatorMintRecovery (many arms) |
| J15 | dispatch wraps `openRunnerLock` with O_CREATE | KILLED | TestOperatorMintRecovery/missing_lock |
| J16 | `held.Name() != path` -> `== ""` | KILLED | TestOperatorTerminalExclusion/nonwaiting |
| J17 | ignore `held` in canonical lock helper | KILLED | TestOperatorTerminalExclusion/nonwaiting |
| J18 | ORDER BY id, created_at | KILLED | TestNextOperatorStart |
| J19 | `kind <> 'finish'` (must-live) | SURVIVED whole host/store | equivalent by schema CHECK |
| J20 | state IN ('confirmed','requested') | KILLED | TestNextOperatorStart |
| J21 | LIMIT 2 (must-live) | SURVIVED whole host/store | equivalent (QueryRow) |
| J22 | move digest fact after launch_wall | KILLED | TestOperatorFilePublicationAdmission and others (behavioural; incidental) |

Whole-package cmd runs for survivors were not done (J05, J07 are leak/equivalent by construction); the named + `Operator` selector was run for each.

**Controller KILLED sample re-run (11, all behavioural, all confirmed KILLED; ledger `judge/ledger-rerun.jsonl`):** S9, S10, E-RES-EPISODE-CONFLICT, E-RES-REREAD, D-CATALOGUE-UNCHANGED, D-CAND-LOCK-RELEASE, D-EXEC-DIGEST-READ, Z-ABSENT-ANY (TestOperatorTerminalCreateLockMissingParent), Z-WAIT-READBUDGET (TestOperatorWaitReadBudget), D-SELECT-EXPIRES, D-CAND-LOCK-IDENTITY. Controller ledger `plan-scratch/ledger-ctl-named.jsonl` = 45 lines; my reruns agree on all 11.

## Tests review
Non-vacuous: real migration, real signer, real `exec` re-exec child holding a terminal lock (8 s ctx, pipe barrier, `atexit_sleep_ms=0`), injected clocks, no real sleeps, state asserted by read-back (hash of credentials DB + file/worktree trees + lock reacquire). Portable (flock via existing helpers; vet clean on linux). Existing-test hunks read: GC table gained `kind` column and four rows (additive); `credentials_test.go` +1 line, `operator_s5e4_test.go`/`operator_runner_test.go` fixture signing now stamps `state='confirmed'` and a catalogue-liveness assertion; start-input tests only index-shifted for the new fact. No assertion weakened. Gaps: unchanged-NULL catalogue (J01); `Close()` leak (J05); S5f2's four process pins remain deferred by design.

## Findings
BLOCKING: none.
Non-blocking, ranked:
1. Appendix C carries no measured S5f2 size in the document (I measured 1,082 changed cmd+host lines vs base); plan §1.5/§4/§6/§7 still "pending" (planner refusal). Doc should be completed in the record step.
2. J01: heartbeat catalogue comparison for NULL catalogue untested (`IS` vs `=`), whole host/store survives.
3. Existing signed start actions (none in production) no longer reconstruct after the new `launch_binary_sha256` fact; fine while uncalled, but S5f2 must state the compatibility rule.
4. Dormant `withOperatorStartCandidate`/`NextOperatorStart` have no caller: judged a justified staged seam (pinned by real-process tests, 47+12 production lines), not dead code to reject.
5. Controller `callers` leg harness bug (`$4: unbound variable`) and pre-existing gofmt hit on `host/daemon/daemon.go` are not attributable to this sprint.

Worktree restored: `git status --porcelain` is empty.
