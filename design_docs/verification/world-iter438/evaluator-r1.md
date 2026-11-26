# Evaluator r1 — iteration 438, row 243 V4 S5e4 (base 684069f2, head 1b78d779)

VERDICT: PASS — 88 / 100. No BLOCKING findings.

## Scores
- Correctness 28/30
- Tests 25/30
- Security / fail-closed 19/20
- Coding standards 9/10
- Scope / size 7/10   (-3: Appendix C measured but N04/N12/N15 left carried; whole daemon-race sits at the bound; test gap noted below counted in Tests)
Total 88.

## Measured
- numstat `git diff -M --numstat 684069f2 1b78d779 -- cmd host docs app` = 929 changed lines (planner 929; cap 1,400).
- Own runs: `go test -count=1 ./host/store/` ok 12.2 s. Focused race: store 3.5 s total (new tests ~2 s), daemon `TestEpisodeGCOperatorSteps` 2.3 s package / 0.65 s test, cmd `TestEpisodeGCFailSafeAndReclaims` 2.6 s package / 0.98 s test. `go vet ./cmd/... ./host/...` ok; `GOOS=linux go vet` of the three packages ok; `go test -run '^$' ./...` clean.
- Controller legs (read): store plain/race 0 fails (36-40 s); daemon-race fails = 4 (the 3 known + TestSeToolsRunBenchmarkTasksEndToEnd, known load-sensitive), 479 s at head vs 467 s at base (rig load; whole-package race is at the 480 s bound at base too, so not attributable to this sprint, but the margin is nil); cmd-plain 11 fails == known set. gofmt flags only `host/daemon/daemon.go`, unchanged by this sprint and flagged at base too.
- Did NOT run whole daemon/cmd packages (UNMEASURED by me; relied on controller legs + selectors `EpisodeGC|Episode|GC`).

## Judge-437 / judge-436 survivors (re-applied by me, whole `host/store` run for each, named test = killer; all package host/store)
N22 TestOperatorFileStoreReadPins/headless; N23 .../duplicate_worktree; N24 TestOperatorFileStoreReservePins/foreign_lock; N25 .../empty_identity; N27 .../before_created; N28 .../different_approver; N29 .../allowance; N31 .../action_expired; N33 TestOperatorMintStoreFilesFence/read; N34 .../abort; J436-heartbeat lower edge (seen_at > -> >=) TestOperatorWorktreeStoreReservePins/heartbeat_edge; J436 no-episode-start TestOperatorWorktreeStoreReservePins/existing_start; no-credential .../existing_credential; J436 READPHASE TestOperatorWorktreeStoreRecoveryPins/observed; READOTHERKIND .../read_other_kind; DISPOTHERKIND .../dispose_other_kind. All KILLED, all behavioural (no static/AST test involved). Note: the upper heartbeat edge (seen_at <= -> <) is killed only incidentally by the GC lookup fixtures, not by heartbeat_edge.
N01 (cmd, `start.Allowance != nil` -> false): re-run SURVIVES `OperatorFile|OperatorStart|Allowance` in cmd — confirmed MASK: operator_files_test.go:598 shows the refusal is `sql.ErrNoRows` (no observed mint receipt) before operatorFilePaths; the later layer is independently pinned by store N29 (killed). Accepted.
N04/N12/N15: not re-run by me (UNMEASURED); parked in Appendix C.

## My M3 mutants (ids M-*, disjoint from the planner's GC-* ids; scripts/logs in $HOME/.ailang/state/world-iter438/judge/, `s3-store.out`, `s3-dc.out`, mut-*.log; 31 mutants incl. controls)
Store (killing package host/store, test TestEpisodeHasUnresolvedOperatorSteps unless noted):
- S1 drop `kind <> 'operation'` KILLED /operation_only
- S2 phase-only KILLED /observed_running,/observed_uncertain (MUST-DIE control: killed)
- S3 terminal-only KILLED /terminal_with_reservation
- S4 OR->AND KILLED (3 subtests)
- S5 episode filter dropped KILLED (/none and others)
- S7 version gate `<0` KILLED TestEpisodeOperatorStepLookupRefusals/older_schema
- S8 requireDeadline removed KILLED .../no_deadline + TestCredentialsMethodsRequireADeadline
- S11 exclude worktree kind KILLED
- S12 open error swallowed KILLED .../missing_store,/deadline
- S6 JOIN->LEFT JOIN: SURVIVED named + whole host/store (equivalent given the kind filter; MUST-LIVE control: lived)
- **S9 exclude kind 'files' from protection: SURVIVED (named + whole host/store)**
- **S10 exclude kind 'mint': SURVIVED (named + whole host/store)**
Daemon (host/daemon, TestEpisodeGCOperatorSteps; named selector `EpisodeGC|Episode|GC`, no whole-package run): D1 drop ctx.Err arm KILLED /canceled; D2 drop err arm KILLED /error; D3 invert unresolved KILLED /candidate,/protected; D4 default reason KILLED /candidate; D5 nil-lookup allows KILLED /error; D6 skip disabled KILLED /protected; D7 daemon wiring path "" KILLED /daemon_wiring; D8 reason text KILLED /protected; D9 real clock KILLED /daemon_wiring.
CLI (cmd/ailang-worldd, TestEpisodeGCFailSafeAndReclaims, selector `EpisodeGC|Episode`; no whole-package run): C1 explicit flag KILLED; C2 malformed record fallthrough KILLED; C3 root compare removed KILLED; C4 no EvalSymlinks(root) KILLED; C5 no EvalSymlinks(workspace) KILLED; C6 no Abs KILLED (/relative_root); C7 path-error join removed KILLED; C8 missing-home KILLED (/missing_home); C9 default workspace KILLED; C10 closure uses flag value KILLED; C11 record credentials ignored KILLED.
Planner-ledger sample re-run (own equivalents of the same edits): GC-D-EPISODE, GC-C-WIRING, GC-C-FLAG killed by the same named tests; GC-ST-TERMINAL/OR/RESERVED/ACTION-JOIN/EPISODE/DEADLINE/OLD-SCHEMA/OPEN-ERROR (via S1-S8, S12), GC-D-CANCELED/ERR/TRUE/SKIP/WIRING/CLOCK/NIL (D1-D9), GC-C-EXPLICIT/HOME/SAME-ROOT/PHYSICAL-*/MALFORMED (C1-C11) — all KILLED, matching the ledger (38 rows). `comm`: my ids ∩ ledger ids = empty by construction (M-* vs GC-*/P-*).

## Findings (non-blocking, ranked)
1. Test gap (Tests -4): the lookup table covers only `worktree`-kind steps in the scanned episode (the `cleanup` step is in another episode). Mutants excluding `mint` or `files` from protection survive whole host/store (S9, S10). Production SQL is kind-agnostic except `operation`, so behaviour is correct today (S5d/S5e3 effects: worktree/mint/files, reserved/observed/uncertain, all protected unless action terminal AND step observed), but a future narrowing of the kind list would not be caught. Suggest adding mint and files rows to the table.
2. Usefulness cost: any lookup error (credentials absent, schema < v10, locked, path unresolved) skips every worktree. Fail-safe and visible in the report reason, but a home with no operator-capable store never reclaims worktrees (daemon: d.credentialsPath is the daemon's own store; CLI: default home or record). Acceptable for safety; worth a doc line.
3. Read-only lookup: openOperatorReadOnly + busy_timeout(0) + ledger-bytes-unchanged assertions in all refusal tests; a locked store errors -> skip.
4. whole-package daemon race at 467 s (base) / 479 s (head) vs 480 s bound: pre-existing; sprint adds ~0.4 s.
5. Heartbeat upper edge only incidentally killed (see above).

## GC policy view
Rule: protect when any non-`operation` step exists for the episode's action and (step.phase='reserved' OR action.terminal_at IS NULL). Observation alone does not release a running/uncertain action; a reserved step protects even a terminal action; observed + terminal (succeeded/failed/expired/cancelled) releases; `operation` rows are alias identity rows (schema CHECK: observed, no owner/receipt), never an effect, so excluding them is safe; an episode with no actions is a candidate (correct). CLI path resolution: explicit `--credentials-db` > home service record (only if its workspace EvalSymlinks-equals the scanned root) > default `~/.ailang/world/credentials` for the default workspace; mismatch, malformed record, missing HOME, unresolvable roots -> "" -> lookup error -> skip. It cannot be steered to a different store by the default path (tests prove an unrelated default store does not clear a custom workspace); only an operator-supplied explicit flag overrides. Daemon uses d.credentialsPath. Both bounded by existing contexts (daemon autoGCPassBudget, CLI 2 min); no new context root (TestProductionContextRoots in whole host/store passes).

## Inertness
Production callers (non-test) of reserveOperatorStartWorktree, MintOperatorStart, publishOperatorStartFiles: each still has only its definition plus the pre-existing S5e3/service-reconcile references (controller callers.txt unchanged); no route/schema change; only change in served behaviour is GC skipping more.

## Test realism
Real migration, real signer/approval fixtures (approvaltest), refusals asserted by byte/row read-back, no sleeps/time.Now, injected clocks, t.Setenv HOME/TempDir so no real-path writes; the cmd test pins the mtime to a fixed epoch (strengthened, not weakened); os.IsNotExist->errors.Is hunks equivalent; credentials_test.go +1 line adds the new function to the deadline table. Fixture uses direct SQL UPDATE to model "reserved step outliving a terminal action" on the real schema — acceptable.

## Appendix C
Present, measured (N04/N12/N15 parked: 82 lines, +4.4 s CLI race; S5f consumer; four carried process pins; zero new patch lines). I did not re-measure the 82/4.4 s figures.
