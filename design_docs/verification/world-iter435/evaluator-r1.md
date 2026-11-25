# Evaluator report r1 — iteration 435, row 243 V4 S5e M1 (head fd17279f on base ec4565ec)

## Verdict: PASS — 89 / 100 (no BLOCKING findings)

| Category | Score |
|---|---|
| Correctness vs design / judge-434 survivors / plan §9 | 27 / 30 |
| Tests (non-vacuous, kill, read-back, portable) | 25 / 30 |
| Security / fail-closed | 19 / 20 |
| Coding standards | 9 / 10 |
| Scope and size discipline | 9 / 10 |

## BLOCKING findings
None.

## Non-blocking findings (ranked)
1. Owner guard in `ReserveOperatorWorktree` is unpinned. `row.Owner != at.Runner.Nonce` → `false` (Q-ST-OWNER, host/store/operator_worktree.go:31) SURVIVES whole `host/store` and the whole `cmd/ailang-worldd` package (fail set == the 11-test base set, `comm -23` empty). It is not a judge-434 survivor, but it is a pre-cap admission guard on the path M1 claims to pin. The `runner_state` SQL checks `Nonce`, the UPDATE uses `row.Owner`, so no other check refuses a row claimed by a different runner. Add an `owner` arm to `TestReserveOperatorWorktreeGuards`.
2. `Q-W-ABSENTDISP` (cmd/ailang-worldd/operator_worktree.go:101; drop the `addErr != nil` clause so a successful add with an absent readback is recorded as "stopped: absent" instead of uncertain) SURVIVES the whole cmd package (fail set == base, no new failure). Unpinned and arguably wrong-direction: the added-but-absent case should stay uncertain. Not a judge-434 survivor.
3. Other unpinned guards that survived their narrowed selector. The store ones also survived whole `host/store`, but the whole cmd package was NOT run for them (time), so these are unverified survivors, not confirmed:
   - Q-ST-IDENT: `|| identity == ""`
   - Q-ST-EXPIRYEDGE: `>=` → `>`
   - Q-ST-READSTATE: allow only "running" in `ReadOperatorWorktree`
   - Q-ST-READCOUNT: `= 1` → `>= 1`
   - Q-ST-DISPKINDSQL: SQL `kind = 'start'` in dispose
   - Q-R-SIGNEDKIND, Q-R-REPOABS, Q-R-COMMONABS in `readRunnerWorktree`
   - Q-R-REREADERR: error ignored on the locked re-read
   
   Most are probably masked by SQL or key CAS; I did not prove each.
4. Q-RB-EXACT (`exact != 1` → `exact < 1`) and Q-RB-RAWABS (`!ok || !IsAbs(raw)` → `!ok`) survived named runs. Both are plausibly equivalent: the registration-count check covers exact≥2, and Git always prints absolute paths. I judge both equivalent but did not run them whole-package.
5. The plan's own table marks CLI whole-package rows "UNINFORMATIVE UNDER SANDBOX". The controller sweep and my reruns on the real rig fill that gap (see below).

## Judge-434 survivors and the four claimed equivalents
All 30 survivor ids are in the plan's §4 table; none was silently dropped. I re-measured 15 of the 31 named kills myself.

- Store kills: J-ST-KIND, J-ST-STORE, J-ST-APPROVEDBY, J-ST-CHALLENGE, J-ST-UNC-OWNER, J-STEP-RECEIPTMAX.
- CLI kills: J-W-PHYS-ABS, J-R-BRANCHCHK, J-ST-READPHASE, J-RB-WTCOMMON, J-RB-LINKREG, J-RB-ISLINKED, J-RB-REGCOUNT, J-R-CURRENTKEY.
- Hunk kills: M1-ATTENDED-DISPATCH, M1-LOCKED-FRESH-READ.

All of these were killed by behavioural tests that run real store, Git and CLI operations. None was killed by an AST, inventory or string test, so no re-run with a test excluded was needed.

The four claimed equivalents all SURVIVED their named test, and I accept each:

- **J-R-CURRENTEP:** the key's digest binds the canonical episode, and `decodeOperatorMintAction` enforces row/episode agreement. `current != key` already covers any episode drift.
- **J-W-CREATEDREAD:** every error return in `inspect` precedes the only assignment of `treeExact` (line 258). `readErr != nil` therefore implies a non-exact state, and the plan's masking invariant is true.
- **N-STORE-RAW-CRED:** `credentialsFor` canonicalises first (session.go:129), and `CanonicalPath` is idempotent.
- **N-RES-ABS:** a relative path is still refused downstream by `runnerPath` (operator_service.go:62), which is reached through `inspectRunnerMintToken`.

I ran these four only against their named tests. The whole-package fail-set run for them rests on the plan's and judge-434's earlier runs.

## Production changes
- **Store predicate.** `row.Kind != "start"` moves ahead of decode, and `action.Kind != "start"` is dropped. `decodeOperatorMintAction` already enforces `a.Kind == row.Kind` (operator_mint.go:183), so this is behaviour-preserving for every valid input. The only change is that a damaged non-start row now returns `ErrOperatorConflict` instead of a decode error. It adds no reachable path.
- **Reader seam.** `reconcileRunnerWorktree` gains a `readReservation` parameter. Its single production call site passes `readRunnerWorktree` (operator_service.go:290), so production only ever reaches the real reader.
- **Existing tests.** Numstat shows two new files and no edit to any existing `_test.go`, so no assertion was weakened.

## Non-vacuity per hunk
- **Store predicate, literal revert** (`git diff ... | git apply -R` on `host/store/operator_worktree.go`): `TestReserveOperatorWorktreeGuards/kind` FAILS (host/store).
- **Store predicate, deleted:** J-ST-KIND is KILLED by `TestReserveOperatorWorktreeGuards/kind` (host/store).
- **Reader seam ignored** (parameter replaced by `readRunnerWorktree`): H2-IGNORE-READER is KILLED by `TestOperatorWorktreeLockedReread` (cmd/ailang-worldd).
- **Dispatch argument:** M1-ATTENDED-DISPATCH is KILLED by `TestOperatorWorktreeWindow/before_add`, `add_error_after_checkout` and `before_observe`.

## Test quality
- **Fixtures:** they use real `git` in `t.TempDir()`, the migrated credentials store and the repo signer (`approvaltest.NewSigner`).
- **Decoy fixture:** the `checkout_common` case first asserts that HEAD and the symbolic ref of the decoy match the target. Only the common-dir check can refuse it, and it is killed by J-RB-WTCOMMON.
- **Refusals** compare `observeProvision` / row snapshots before and after, and count operator rows and sessions.
- **Portability:** no `Sleep` or timers (grep found none), no `GOOS` or darwin branches. Paths go through `EvalSymlinks`, and the shell wrappers use POSIX `sh`, `cp`, `mv` and `ln`. I did not run on Linux.
- **Gaps:** see findings 1–3. The test is thorough for the listed survivors but did not extend to neighbouring guards.

## Inertness and activation
- Production callers (excluding `_test.go`), with positive control:
  - `reserveOperatorStartWorktree`: 0 callers (definition only).
  - `ReserveOperatorWorktree`: called only inside `reserveOperatorStartWorktree`.
  - `MintOperatorStart`: 0 callers. `CommitOperatorMint`: only inside `MintOperatorStart`.
  - `ReserveOperatorStep`: 0 callers (definition only).
  - `ObserveOperatorStep`: only inside `reserveOperatorStartWorktree`.
  - `reconstructOperatorStart`: only inside `reserveOperatorStartWorktree`.
  - `ReadOperatorWorktree`: only inside `readRunnerWorktree`.
  - Positive control: `readRunnerWorktree` and `reconcileRunnerWorktree` (the attended verb) each have a live caller in `operator_service.go`.
- The controller's `routes.txt` shows DAEMON_ROUTES 54 (base 54) and SURFACE_POSTS 3 (base 3). `daemon.go`, `surface_page.go` and `app.js` are UNCHANGED, so no new daemon route and served bytes are unchanged.

## The split
Shipping M1 alone is defensible. Cap command `git diff -M --numstat ec4565ec HEAD -- cmd host docs app` gives, per file (adds/dels):
- operator_service.go: 1/1
- operator_worktree_guards_test.go (cmd): 315/0
- operator_worktree_reconcile.go: 4/3
- operator_worktree.go (store): 2/2
- operator_worktree_guards_test.go (store): 167/0

That is 489 added + 6 deleted = **495**, matching the plan. Appendix C is complete and measured: publication plus its attended disposition kept together (967 incremental, 1,462 cumulative against the 1,400 cap), the S5f/S6 paired-process pins, and the carried items E-RES-EPISODE-CONFLICT, E-RES-REREAD, J-OWN-UID, Z-ABSENT-ANY and Z-WAIT-READBUDGET. I did not re-measure the 1,462; the plan cites its archive.

## Race cost
`-race` of the five new cmd tests: 0.49 + 3.56 + 1.99 + 0.30 + 0.51 = **6.85 s** of test time, 8.5 s package wall for that selector. The planner measured +6.0 s, so the number is confirmed (+0.85 s). The three new store tests under `-race` took 1.5 s of test time. Cost is acceptable.

## Mutant ledger
Ledger: `$HOME/.ailang/state/world-iter435/judge/mutants.tsv`, 58 lines, 56 distinct ids. Runner `judge/mut.py`, specs `judge/my-mutants.json` and `plan-sample.json`, logs `judge/log-*.txt`. Worktree `git status --short` is clean.

`comm -23` of my ids vs the plan's ids: 25 ids not in the plan (CTL-DIE, CTL-DIE-STORE, CTL-LIVE, CTL-LIVE-STORE, Q-R-BEFOREAFTER, Q-R-COMMONABS, Q-R-DISPOSITION, Q-R-REPOABS, Q-R-REREADERR, Q-R-SIGNEDKIND, Q-R-UNMARSHAL, Q-RB-EXACT, Q-RB-HEADREF, Q-RB-RAWABS, Q-RB-REFS, Q-ST-DISPKINDSQL, Q-ST-EXPIRYEDGE, Q-ST-IDENT, Q-ST-OWNER, Q-ST-READCOUNT, Q-ST-READSTATE, Q-STEP-EMPTY, Q-W-ABSENTDISP, Q-W-CREATEDEXACT, Q-W-INITIAL).

Own mutants (verdicts are for the named selector unless marked whole):

| Id | Pkg | Verdict |
|---|---|---|
| Q-STEP-EMPTY | host/store | KILLED by `TestOperatorStepGuards` (existing) |
| Q-RB-HEADREF | cmd | KILLED by `TestOperatorWorktreeWindow/redirected_checkout` |
| Q-RB-REFS | cmd | KILLED by `TestOperatorWorktreeWindow/ref_read_skew` |
| Q-W-CREATEDEXACT | cmd | KILLED by `TestOperatorWorktreeWindow/head_changed` |
| Q-R-DISPOSITION | cmd | KILLED by `TestOperatorWorktreeWindow/before_add` |
| Q-W-INITIAL | cmd | KILLED by `TestOperatorStartCollisions` |
| Q-ST-OWNER, Q-W-ABSENTDISP | store / cmd | SURVIVED whole package (findings 1 and 2) |
| Q-ST-IDENT, Q-ST-EXPIRYEDGE, Q-ST-READSTATE, Q-ST-READCOUNT, Q-ST-DISPKINDSQL | host/store | SURVIVED whole host/store only |
| Q-R-SIGNEDKIND, Q-R-REPOABS, Q-R-COMMONABS, Q-R-REREADERR, Q-RB-EXACT, Q-RB-RAWABS | cmd | SURVIVED the named selector only |
| Q-R-BEFOREAFTER, Q-R-UNMARSHAL | cmd | COMPILEERR (not counted as mutants) |

Controls:
- Must-die CTL-DIE-STORE (`row.State == "claimed"`): KILLED by `TestReserveOperatorWorktreeGuards` (host/store).
- Must-die CTL-DIE (`initial.state == 99`): KILLED by `TestOperatorStartCollisions` (cmd).
- Must-live CTL-LIVE-STORE (`row.Kind != "start" || false`): SURVIVED, including whole host/store.
- Must-live CTL-LIVE (an unasserted error string): SURVIVED the named cmd selector.

Ledger caveat: the line `Q-ST-OWNER-WHOLECMD` was machine-labelled KILLED because the whole-cmd run showed the base's 11 e2e failures. I added a correcting line marking it SURVIVED after fail-set comparison (`comm -23` empty against `plan-scratch/pkg-cmd-ailang-worldd.base-fails`). `Q-W-ABSENTDISP-WHOLECMD`'s machine label is similarly noisy; its `comm -23` against the base set is also empty, so it SURVIVED whole.

## Test runs (after ALL_DONE.fd17279f)
- `go test -count=1 ./host/store/` ok, 16.8 s.
- `go test -count=1 -run 'Operator|Preset|Session|Quickstart|Help|CredentialWriters' ./cmd/ailang-worldd/` ok, 93.3 s.
- `go test -race -count=1 -run 'OperatorWorktree|OperatorStart|OperatorStep|OperatorMint|ReadbackError' ./cmd/ailang-worldd/ ./host/store/` ok, 98.8 s and 2.0 s.
- `go vet ./cmd/... ./host/...`: no output.
- `go test -run '^$' ./...`: no failures or compile errors.
- (I read the pipeline output rather than `PIPESTATUS`, because zsh did not expose it; every run printed `ok`/no output.)
- I did not run `verify_go.sh` or a whole `go test ./...`, per the directive.

## Gate copy and `verify_go` leg
`diff executor/gates.sh executor/gates_ctl.sh` shows three changed lines, all in the `verify_go` leg: no `GOFLAGS=-json`, the bound raised from 1100 to 1500 s, and the fail-set parsed by `text_failset.py` instead of `failset.py`.

The controller's `legs.txt` shows `verify_go rc=1 (659 s)`. Its fail set is the 14 known row-271 tests (se-tools 10→12) plus three known load-sensitive tests: `TestSeToolsRunBenchmarkTasksEndToEnd`, `TestSeToolsRunReplayIsByteEqual` and `TestMissionConfigDefaultFlowIntegration`. All are on the directive's known-red or load-sensitive lists, and none involves the touched files. I judge the leg's verdict sound for this purpose: rc=1 is the baseline, and there is no new failure. `GATES_DONE 1497s FAILS: verify_go UNINFORMATIVE: none`. `observation-before` and `observation-after` (real `~/Library/LaunchAgents`, `~/.ailang-world`, credentials store) are byte-identical.

Every other leg is rc=0: gofmt, build, vet, fence, `verify_ail`, `pkg-cmd-ailang-worldd`, `race-cmd-ailang-worldd`, `pkg-host-store`, `race-host-store`, the inventories, `routes`, `start-input`, `s5e-accept` and `added-cli-race`.
