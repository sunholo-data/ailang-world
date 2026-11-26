# Evaluator report, iteration 437, row 243 sprint S5e3 (judge: sonnet, independent)

Verdict: **PASS, 79/100**. No BLOCKING finding.
Head `de608652`, base `eaa202b1`. My worktree `.wt-world-iter437-row243-v4s5e3-eval` is clean (`git status --short` empty). Ledger `$HOME/.ailang/state/world-iter437/judge/mutants.tsv` has 70 lines (harness `mut.py`, `plan.py`, `mine.py` beside it).

## Scores
| Category | Score | Note |
|---|---|---|
| Correctness (30) | 26 | Effect trace holds; two un-pinned directory-swap windows, both fail safe |
| Tests (30) | 20 | Real fixtures and killing tests. My own sample killed only 9 of 32 valid runs |
| Security / fail-closed (20) | 18 | Token never stored raw; modes exact under umask 0777; refusals before prompt |
| Coding standards (10) | 8 | Dormant, one production entry, seams minimal. A few defensive predicates are untested |
| Scope and size (10) | 7 | Exactly 1,650; the exception was earned. Whole-package race is 456 s against an 8-minute bound, and the planner's 122 AUDIT lines were not measured |

## Effect correctness (code read plus runs)
- **Reservation first.** `publishOperatorStartFiles` (`operator_files.go:262-407`) holds the episode lock throughout. It re-reads state under the lock (:274) and requires an exact observed-created worktree and a matching observed mint with the current owner (store join, `host/store/operator_files.go:39-48`).
  - It checks the raw token hash against the mint's credential id (:279), and checks the signed address, the approved address and the preset-derived address together (:231).
  - It derives MCP bytes from the pinned preset and approved address, and writes nothing for Codex (no MCP file).
  - It requires every target to be absent (:318-325), calls `recheckFileApproval` (live signature, action expiry) and then `RecordOperatorFiles` (reserve). The first `installOperatorFile` is at :372, after reservation.
- **Install and readback.** Install uses `openat` with `O_CREAT|O_EXCL|O_NOFOLLOW|O_CLOEXEC` (0600), then `fchmod` 0600, write, fsync, fstat (:205-222). Before each file the code re-reads state, checks expiry (:361-366) and re-opens the directory path (:367).
  - Readback is independent and directory-relative: `fstatat` NOFOLLOW, 0600 mode, owner, nlink=1, open NOFOLLOW|NONBLOCK, fstat identity equal, bounded read, final stat, digest (:127-165).
  - The readback inode must equal the installed inode (:388). The code re-inspects the worktree and re-reads the set (:395-398), rechecks approval, then records `published`.
  - Every failure after reservation becomes `uncertain` plus the attended remedy. Nothing is removed or rewritten anywhere.
- **No unreserved write, overwrite or symlink follow, and no observation without readback.** I found no path for any of these. The only `O_CREAT` is :205, reached only after :339 `record(wire)` succeeded. `published` is recorded only at :403, after readback and the re-read.
- **Production callers.** `publishOperatorStartFiles`: 0 non-test callers. `reserveOperatorStartWorktree`: 0. `authority.MintOperatorStart`: 0 (definition at `host/authority/mint.go:204` only).
  - `reconstructOperatorStart` is called only from `operator_worktree.go:45`, inside the unreferenced `reserveOperatorStartWorktree`.
  - `reconcileRunnerFiles` is called only from `operator_service.go:293`, the attended verb. Positive control: `grep` finds that reconcile-files call, and the five test callers of `publishOperatorStartFiles` in `operator_files_test.go`.
  - No daemon route changed. Of the `service`/`main` edits, only `service.go` (a help string) and `main.go` (a usage line) touch served surfaces.
- **#671 allowance origin.** `TestOperatorFileAllowanceOrigin` (`operator_files_test.go:547`) uses real `store.GrantAllowance` and real `authority.MintFromAllowance`. It asserts the `allowance_draws` row exists, `minted.Allowance != nil`, and that `publishOperatorStartFiles` returns `sql.ErrNoRows` before any checkpoint, with provision state byte-equal.
  - Refusal is layered: the missing mint step stops it first, with `start.Allowance != nil` checks at `operator_files.go:62` and `store:219` behind it.
  - Mutants N01 (:62 removed) and N29 (store :219 removed) survive. They are redundant behind the ErrNoRows layer, so the refusal holds, but those two layers are un-pinned.
- **Resolver.** `reconcile-files` refuses before the prompt for a non-reserved phase (:32, which covers both "no reserved row" and "observed row"). Tests assert the terminal opener is never called for the observed row and for the before-reservation case.
  - Evidence is read before the prompt (:57) and after (:79), and any difference refuses.
  - Only a wholly absent or wholly exact set can be stopped (:58-69). Partial or foreign sets refuse and the test asserts no prompt.
  - The store CAS is independent: it re-reads and compares the whole state, then runs an action-owner/digest/generation CAS with a step CAS (`store:169-180`) and a second post-update re-read for `published` and reserve.
  - **Remedy verbatim.** `operator_files_test.go:62-70` extracts the example line from `docs/QUICKSTART.md` and requires it to equal `filesExample`. The `example` is then run through `r.run(t, example)` at :298, :314 and :321+, and the same string is compared to the error's remedy at :260-263 and to the "Retry at a terminal:" line at :326.

## MASKED claims (re-run, named then WHOLE `cmd/ailang-worldd`)
All three survive the named test and the whole package, with fail set equal to the controller's 11-test base set. All three are genuine equivalents:
- `F-CREATE-NOFOLLOW` (:205): `O_CREAT|O_EXCL` already rejects an existing symlink on both Darwin and Linux. Whole run: SURVIVED, fail set = base.
- `R-POST-PARENT` (reconcile :79): a changed parent leaves no child proofs, so the child comparison still differs. SURVIVED, fail set = base.
- `S-WORKTREE-PHASE` (store :43): the schema CHECK `phase <> 'reserved' OR receipt IS NULL` (`host/store/operator_schema.sql:61`) makes `receipt='created'` fail for a reserved row. SURVIVED, fail set = base.

## Planner KILLED sample re-run (11 of 11 confirmed, killing test and package named)
All in package `cmd/ailang-worldd`:
| Mutant | Killed by |
|---|---|
| F-EXCL | `TestOperatorFilesReconcile/competing_file` |
| F-DIR-NOFOLLOW | `TestOperatorFileSyscalls/directory_open` (a behavioural test using the real kernel's refusal) |
| F-DIGEST | `TestOperatorFilesReconcile/metadata/digest` |
| F-ABSENCE | `TestOperatorFilePublicationAdmission/existing_exact_file` |
| R-RESERVED | `TestOperatorFilesReconcile/before_reservation` |
| S-LATE-APPROVAL | `TestOperatorFilesReconcile/late_key_revoked` |
| S-MINT-OWNER | `TestOperatorFilePublicationAdmission/wrong_mint_owner` |
| F-WRITE-IDENTITY | `TestOperatorFilesReconcile/replaced_before_readback` |
| R-STEP-VERSION | `TestOperatorFilesReconcile/post_step_version` |
| F-RAW-TOKEN | `TestOperatorFilePublicationAdmission/wrong_token` |
| F-MCP-RENDER-BEARER | `TestEpisodeStartProvisionsAndPrintsTheLaunchWithTheTokenInTwoFilesOnly` |

Notes:
- The planner's S-LATE-APPROVAL replacement (`approvalProof, err = approvalProof, nil`) is rejected by `go vet` (self-assignment) in my vet-first harness. I re-ran it as `approvalProof, err = func() (...) { return approvalProof, nil }()` and it is killed.
- No kill above rests only on a static or AST test. Every killer is a runtime test.

## My own mutants (ids not in the plan; `comm -23` of my 36 distinct ids vs `plan-ids.txt` leaves all 36)
Selector: `^(TestOperatorFile|TestOperatorStart|TestOperatorMint|TestOperatorWorktree|TestEpisodeStart)` in `cmd/ailang-worldd`. Whole package is run only where marked.

**Killed (9):**
- N02 `filepath.Dir(TokenPath) != set.Dir` removed: `TestOperatorFilePublicationAdmission/linked_destination`.
- N06 approved-address compare: `.../wrong_approved_address`.
- N10 `treeExact` check: `.../worktree_moved`.
- N11 directory-inode pin: `TestOperatorFilesReconcile/before_first`, `before_observe` and others.
- N18 disposition swap: `TestOperatorFilesReconcile/*`.
- N26 `state IN ('running','uncertain')` narrowed: `late_key_revoked` and others.
- N32 other-effect NOT EXISTS: `.../later_step`.
- N35 `operatorPhysicalTarget` on the token path (`operator_start_input.go:39`): `.../expired_action` and others.
- **Must-die control X-DIE-EMPTYWRITE** (`f.Write(data[:0])`): `TestOperatorFileReadbackMetadata/exact`, `TestOperatorFileReopen` and others.

**Survived the named selector (not run whole-package except N33 and N34; time-boxed):**
- N01, N03 (:119 fstat identity), N04 (:194 post-readback directory re-open compare; first run was a compile error because the variable was unused, re-run with `&& false` survived), N05 (:231 `d.Addr != signed`), N07, N08, N09, N12 (:367 per-file directory re-open), N13 (dir fsync error), N14 (:384 got.Dir.Status), N15 (:396 post-readback worktree compare), N16, N17, N19, N20.
- Store predicates N22 (headless), N23 (worktree count), N24 (`lock.dbPath`), N25, N27 (created-at), N28 (ApprovedBy), N29, N31 (row expiry). I re-ran each against the whole of `host/store` and all of them survive there too. **No `host/store` test exercises `operator_files.go`; the store change is proven only by the CLI tests**, apart from the deadline-signature test in `credentials_test.go`.
- **Must-live control X-LIVE-STDERR-LABEL** (the stderr label text): survives named, as expected. I did not run it whole-package.
- N30: compile error (invalid, not counted).

**Whole-package runs, fail set = base (11 tests):**
- N33 (`ReadOperatorMint` passes `fileStage=true`) survives.
- N34 (`AbortOperatorMint` passes `fileStage=true`) survives.
- Each is masked by the other: the CLI read refuses before the prompt, and the abort re-read refuses in the transaction. No test pins either one, and no test shows `reconcile-mint` refusing once a files reservation exists.

Reading of the survivors: most are redundant layers behind a second guard (N01/N29, N07, N08, N14, N16, N33/N34). The planner's 122 AUDIT lines were therefore not mutation-measured and are only partly pinned: roughly two thirds of my valid samples survive, against the planner's 95/98 on its own set.

## Non-blocking findings, ranked
1. **Untested directory-swap windows.** N04 (`readOperatorFileSet` directory re-open, :194) and N12 (:367, per-file directory path check) survive. A parent swap between files or during readback is detected only through later checks. It always ends `uncertain` and never unattributed removal, so it fails safe. No test has a seam inside `readOperatorFileSet`, so adding a pin is non-trivial.
2. **N33/N34 test gap.** Nothing asserts `reconcile-mint` refuses once a `files` step exists. A future edit to either guard would pass the full suite.
3. **Store admission predicates have no `host/store` test** (N22-N31). They are exercised only through CLI tests. This repeats the S5e2 pattern of store pins being deferred. The plan's Appendix C does not list these.
4. **N15** (post-readback worktree re-inspect, :396) un-pinned.
5. **Whole-package race time:** base 405 s, head 456 s (controller legs `race-cmd-ailang-worldd`), against the CI 8-minute (480 s) per-package bound. That is 95% of the bound on this loaded rig (load average about 9). It is a real CI risk, not measured on Ubuntu. Marginal cost is about +51 s.
6. **The plan's `gate-bin/git` PATH wrapper** (the one the controller found `exec`ing an env var) is gate-only scaffolding that tests with a scrubbed child env silently defeat. It should not exist as a gate input. The controller's recorded `.contaminated-gitwrapper` files show it already skewed a result once. I did not use it.
7. **Existing tests edited.** `operator_start_input_test.go`, `operator_worktree_test.go`, `help_test.go` and `credentials_test.go`: every hunk checked, no assertion weakened.
   - Hunks are index shifts (`Resolution[10]`→`[11]`, `[11]`→`[12]`), the new `served_address` fact name and value, `operatorPhysicalTarget` canonicalisation of the token path, an optional `ttl` parameter, and the `reconcile-files` help strings.
   - The deadline test adds two entries.

## Test quality and portability
- **Fixtures are real.** Real signer, real migration, real `git` in `t.TempDir()`, real publication and resolver code, real allowance producer. State is asserted by read-back: store rows, file bytes, modes and Nlink, and the Git/worktree snapshot via `observeProvision`.
- **No macOS-only assertions.** I found no `/private`, `runtime.GOOS` or `darwin` use in the new tests. They use `fcntl F_GETFD` and kernel `openat` refusals, which are portable. `GOOS=linux go vet ./cmd/ailang-worldd ./host/store` is rc 0.
- **No real sleeps.** The only timing is a bounded `time.After(time.Second)` at `operator_file_syscalls_test.go:148`, in a lock test.
- **Security checks present.** The test asserts the reservation identity never contains the raw token (`strings.Contains(want.Step.Identity, token)`) and runs the all-files arm under `Umask(0777)`. No test writes outside `t.TempDir()`, and I touched none of the D-WORLD-109 paths.

## Size and CI cost
- `git diff -M --numstat eaa202b1 de608652 -- cmd host docs app` sums to **1,650** (additions plus deletions, 16 files), matching the planner.
- The exception was earned: publication, readback and the attended disposition are one safe unit, and the S5e2 judge said a third deferral should be refused. The two test files are 809 lines of the total. Nothing looks like dead weight except possibly the redundant `start.Allowance` checks.
- `-race -run '^TestOperatorFile'`: 52.8 s package time (63 s wall) here, against the planner's +46.7 s.
- Whole-package race, from the controller's legs (not re-run by me): base 405 s, head 456 s, against the 480 s bound. See finding 5.

## Appendix C and the #667 GC probe
- Appendix C is complete: optional store pins (`operator_worktree.go:56-59,89,90,138`), carried pins, the S5f/S6 consumer (UNMEASURED, stated), and the #667 gate with the GC change size UNMEASURED.
- It is honest about the zero pin capacity at 1,650. It does not list the N22-N34 store predicates or N04/N12 as pins to carry. They should be added.
- I did not re-run the GC probe. The controller recorded `probe-M1-GC.log` showing both reserved and published worktrees as `clean & merged` candidates at 73 h with `IsActive` stubbed false. The activation gate (operator-step-aware GC skip before any production caller) stands and is the right deferral.

## Test runs (my own)
| Command | Result |
|---|---|
| `go test -count=1 ./host/store/` | ok, 12.5 s |
| `go vet ./cmd/... ./host/...` | rc 0 |
| `GOOS=linux go vet ./cmd/ailang-worldd ./host/store` | rc 0 |
| `go test -run '^$' ./...` | rc 0 |
| `go test -count=1 -run 'Operator\|Preset\|Session\|Quickstart\|Help\|CredentialWriters' ./cmd/ailang-worldd/` | ok, 87.6 s |
| `go test -race -count=1 -run '^TestOperatorFile' ./cmd/ailang-worldd/` | ok, 52.8 s |

Controller sweep (read from `gates-final.out`):
- Every leg is rc 0 except `compare-verify_go`. Its only difference is that the head fail set lacks `host/broker TestExecProtectedCheckTimeout` (a known load-sensitive test), so head is no worse than base.
- Whole-package `cmd/ailang-worldd` fails exactly the 11 row-271 tests at both base and head.
- Real `~/Library/LaunchAgents` and `~/.ailang-world` before/after compare: rc 0.
