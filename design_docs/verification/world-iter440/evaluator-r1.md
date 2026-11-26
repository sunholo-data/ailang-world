# Evaluator report, iteration 440, row 243 V4 sprint S5f2 (head c4fb23f6, base b5f41458)

## Verdict: FAIL (score 67/100, and one BLOCKING finding)

| Category | Max | Score |
|---|---:|---:|
| Correctness (borrow verify, GC exclusion, live behaviour, inertness) | 30 | 23 |
| Tests | 30 | 14 |
| Security / fail-closed | 20 | 16 |
| Coding standards | 10 | 8 |
| Scope and size discipline | 10 | 6 |
| **Total** | 100 | **67** |

## BLOCKING

**B1. The sprint breaks three existing tests that pass at base. M2 changed the fixtures, so these arms no longer test what they claim.** The executor/planner gates could not see it because the CLI whole-package legs were "UNINFORMATIVE UNDER SANDBOX" and the controller's named-mutant replay only ran the named tests.

- `TestOperatorStartCollisions/held`, `TestOperatorStartCollisions/missing_lock` (`cmd/ailang-worldd/operator_worktree_test.go:330`, `:333`).
- `TestOperatorFilePublicationAdmission/lock_held` (`cmd/ailang-worldd/operator_files_test.go:545`).
- Cause: the `r.start` / `publishOperatorFilesFixture` helpers now call `r.borrow(t)` (`store.OpenOperatorEpisodeLock`) before the function under test. In the "held" arms the test already holds the lock, so the borrow itself fails with EAGAIN and calls `t.Fatal`.
- In `missing_lock` the borrow re-creates the removed sidecar, so the function under test never refuses and the expected "episode lock unavailable" refusal never happens. The only refusal arms for "another holder" and "key missing" on the dormant start helpers are therefore dead or red.
- Reproduction (about 4 s, quiet machine):
  `cd <eval worktree> && go test -count=1 -v -run 'TestOperatorStartCollisions$/^(held|missing_lock)$|TestOperatorFilePublicationAdmission$/^lock_held$' ./cmd/ailang-worldd/`
  Head output: `operator_files_test.go:545: resource temporarily unavailable`, `operator_worktree_test.go:330: resource temporarily unavailable`, `operator_worktree_test.go:333: collision missing_lock changed state or wrong refusal: <nil>`, then FAIL.
  The same command on a pristine base worktree printed PASS for all three (`lock_held` 0.43 s, `held` 0.22 s, `missing_lock` 0.22 s).
- The controller's own gate also shows them: `gates-head/cmd-plain.raw` lines 54-61 list `TestOperatorFilePublicationAdmission` and `TestOperatorStartCollisions` as new failures against `base/cmd-plain.fails`.
- CI runs `go test ./...`, so this lands red. Fix direction: the held arms should call the function with a lock opened before the test's competing holder, or assert the borrow's own refusal (`OpenOperatorEpisodeLock` EAGAIN, and re-creation on missing). Restore a real "wrong or absent lock" refusal assertion for the `held` arguments of `reserveOperatorStartWorktree` and `publishOperatorStartFiles`.

## Answers to the directive's questions

**M1 kills.** J01 is killed. Original edit `catalogue IS ?` to `catalogue = ?`: whole-`host/store` run fails only `TestOperatorCatalogueNullRefresh` and its `/null` arm (store). J05 is killed. Dropping the candidate visit's `defer db.Close()` fails only `TestOperatorStartCandidateClosesCredentials` (cmd/ailang-worldd). Both are behavioural (SQLite read-back through the original handle).

**M2 correctness.** Sound in the parts I mutated.
- `Verify` checks the pathname against the held description. It checks `l.file.Name()`, `Lstat` plus `os.SameFile`, regular file, 0600, uid, and `Nlink == 1`, and returns the same `*os.File`.
- `RemoveRegistered` and `Add` pass it as `cmd.ExtraFiles`, so the Git child holds the same open description until it exits. Borrowers never close. Owners close with `defer` (GC `defer lock.Close()`).
- Mutants on `Nlink`, uid, perm, `SameFile`, name, `O_NOFOLLOW`, the episode-id slash check, `ExtraFiles`, the held argument, and the pre-Add and post-Add `Verify` calls are all killed by named behavioural tests (table below).
- No leak or double close was found on the GC path.
- Terminal paths keep the raw `openRunnerLock`. `operatorEpisodeLockPath` is only a pass-through to `store.OperatorEpisodeLockPath`.

**M3, the live change.**
- (a) Ordinary installs. A missing credentials path still gives "operator step lookup unavailable", the same as base. Per-episode behaviour is otherwise preserved, and the existing and new tests reclaim a clean merged worktree (`candidate`, `daemon_wiring`, and the GC probe `disposed published` dry-run candidate). Two real narrowings (new skips, non-blocking, both fail-safe):
  1. `episode gc --credentials-db X` for a workspace that is not the default or home-service workspace used to be honoured. It now gives "operator step lookup unavailable" for every worktree (`episode_gc.go:50-77`; the explicit override must `SameDatabase`-match the paired store). It is documented and pinned in `TestEpisodeGCFailSafeAndReclaims`, but it is a behaviour change to a user-facing flag.
  2. The old `os.RemoveAll` fallback and repository-wide `git worktree prune` are gone. A worktree whose registration is broken or duplicated is now skipped ("worktree registration unavailable or changed") rather than recursively deleted.
- Side effect: GC now creates a persistent `<credentials>.operator-episode-<ep>.lock` sidecar for every scanned episode worktree, never unlinked, even on installs with no operator runner. This is unadvertised clutter.
- (b) Deadlock. None found. The acquire is non-blocking (`LOCK_EX|LOCK_NB`, `O_NONBLOCK`), so contention returns an immediate skip. The lock is held across the ledger read and the Git checks, which are bounded by the pass context (`autoGCPassBudget`) and the per-Git `worktree.Timeout`. `TestEpisodeGCOperatorSteps/gc_first` and `/activation_first` pin both lock orders. The lock is released on every path (`defer lock.Close()`; `J-GC-DEFERCLOSE` is killed).
- (c) Fail-safe and visible. Lookup error, unresolved step, canceled context, bad store path, bad episode grammar, a lock replaced under review, a worktree replaced during scan, and an unlinked, unregistered or duplicate registration all skip with a reason string in `rep.Skipped`. Every one of my removal-path mutants dies (table below). Dry-run is respected (`J-GC-DRYRUN` killed). One gap: `worktree remove` loses the `--force` guard untested (`J-RR-FORCE` survives, see non-blocking N1).
- (d) `D-LOCK-NO-NONBLOCK` (drop `O_NONBLOCK` on open). Reproduced: the whole `host/store` run passes with the mutant (`ledger-judge.txt`: `fails=[] ok`). It is equivalent for the accepted input. The open is `O_RDWR` on a regular single-link file, where `O_NONBLOCK` has no effect, and the FIFO case is refused by the regular-file check (`TestOperatorEpisodeLockFIFORefused` passes either way). I agree it is defence in depth only. `D-LOCK-FLOCK-NONBLOCK` (`LOCK_NB` dropped) is genuinely killed by `TestOperatorEpisodeLockContendsWithoutWaiting`.

**M4p.** N04, N12 and N15 assert the exact uncertain receipt (`"file publication uncertain; "+filesExample`), `Phase=="reserved"`, `Step.Version==2`, mint credential hash, unchanged worktree, branches, workspace and sessions snapshots, and moved-subject contents. The per-call opener seam `readOperatorFileSetWithOpener` is used only by `readOperatorFileSet` (the production default) and the test. It is production-inert. N04, N12 and N15 plus the opener mutants are killed (ledger).

**Inertness of the start path.** `git grep -nE '(MintOperatorStart|withOperatorStartCandidate|reserveOperatorStartWorktree|publishOperatorStartFiles)\(' -- cmd host app ':!*_test.go'` printed exactly four declarations (`operator_files.go:267`, `operator_start_candidate.go:15`, `operator_worktree.go:30`, `host/authority/mint.go:204`), which is the positive control, and zero call sites. `OpenOperatorEpisodeLock` has one production caller (`host/daemon/auto_gc.go:232`) and `RemoveRegistered` one (`auto_gc.go` via the `worktree.RemoveRegistered` default). `runnerTick` is untouched in the diff. No route, schema, `.ail` or app change (the diff touches only the 18 paths in `cmd host`). Only M3's GC change is live.

**Scope and the third deferral.**
- I verified §1.4 from source. I did not re-measure `S5f3-prototype.patch` (UNMEASURED), so I rely on the planner's 1,984 cumulative figure, which implies about 584 changed lines of S5f3 on top of this sprint's 1,400.
- Appendix C and §1.4 are concrete and plannable as ONE sprint: about 584 lines is under the 1,400 cap with the paired drain and four process pins included. No further split is justified unless the paired upgrader work (`service_upgrade.go:208,212` refusing a loaded companion) turns out to be larger than the prototype shows. If it does, split out the upgrade/drain as the first half.
- Spending the full cap on prerequisites was defensible: M1 (73) and M2 (+574) are prerequisites for the caller, and M3 is the exclusion that makes it safe. But the sprint exhausted the cap while shipping a regression to its own tests (B1). That is the clearest evidence that the cap was too tight to also run a CLI whole-package gate.
- CI race cost: measured on a quiet rig, the new and changed cmd tests total `-race` 11.96 s (`EpisodeGCFailSafeAndReclaims` 6.26, `IdentityWindows` 1.11, `StartCandidate*` about 0.6, `BorrowedWorktreeGitProcess` 2.15). Daemon `TestEpisodeGCOperatorSteps` 3.12 s, within the +5 s allowance. Within the +12 s cmd allowance. Under controller load the same tests showed Git 30 s timeouts, which I attribute to load (they pass alone).

## Non-blocking findings (ranked)

1. **N1. `J-RR-FORCE` survives.** Adding `--force` to `RemoveRegistered`'s `git worktree remove` is not caught by `host/worktree` whole, `host/daemon` `EpisodeGC|GC|Operator`, or `cmd` `TestEpisodeGC`. The non-force refusal of a locked or dirty worktree after the checks is untested. Equivalent only while the earlier dirty check holds.
2. **N2.** The CLI explicit `--credentials-db` for non-default workspaces now skips all (documented; see M3 (a)).
3. **N3.** Persistent per-episode `.operator-episode-<ep>.lock` sidecars created by GC on every install. Undocumented clutter.
4. **N4.** The must-live control `J-CTL-LIVE-MSG` (changing the error text `wrong episode lock`) was killed. Tests pin error wording. Harmless but brittle.
5. **N5.** `host/daemon/daemon.go` fails `gofmt -l`; the sprint did not touch it (pre-existing).
6. Daemon race fail set differs from base2 (5 added): four pass alone; `TestSeToolsMCPEndToEnd` fails alone and belongs to the known SeTools family. The cmd-plain set also contained `TestWorldContextAndStatusline` (known load-sensitive).

## Tests review

- Real migration, signer, git, re-exec and `ExtraFiles` fixtures. Read-back assertions (SQLite, `git worktree list`, `show-ref`, file bytes). Injected clocks. No `time.Sleep` in the new tests.
- Portable: `/dev/fd` enumeration and the child `flock(3,...)` use portable semantics; Linux vet passes.
- Existing assertions are not weakened except through B1. The `episode_gc_test.go` hunks add checks (refs, registration, README bytes) and re-point one `assertSkip` reason, which is the documented behaviour change. The `auto_gc_s5e4_test.go` removals are the old `HasOpenOperatorStep` plumbing and `daemon_wiring`, which are re-added.
- Under controller load, Git timeouts fail several heavy cmd tests (`TestOperatorFilesReconcile/*`, `TestOperatorWorktreeRecoveryGuards/*`). These pass alone (`IdentityWindows`, `RecoveryGuards` 0.55 s and 2.96 s). Not attributed.

## Mutation (my own, ids not in mutants.json)

Ledger: `$HOME/.ailang/state/world-iter440/judge/ledger-judge.txt` (my 33 mutants plus 8 reruns; raw logs `judge/mut-*.txt`). Each non-compiling mutant was re-edited until it compiled (a "COMPILEFAIL" line precedes the fixed rerun). `status --porcelain` empty after each.

| id | pkg run | killer (behavioural unless stated) |
|---|---|---|
| J-GC-LOOKUPERR (ignore lookup error) | daemon | `TestEpisodeGCOperatorSteps/error` |
| J-GC-UNRESOLVED | daemon | `.../protected` |
| J-GC-SAMEFILE | daemon | `.../replaced_during_scan`, `/not_linked` |
| J-GC-ACTIVE | daemon | `.../active` |
| J-GC-AGE (re-check after lock) | daemon | `.../recheck_after_ledger`, `/registration_rechecked`, `/unregistered`, `/duplicate_registration` |
| J-GC-DRYRUN | daemon | `.../candidate` |
| J-GC-DEFERCLOSE (lock leak) | daemon | `.../protected` |
| J-GC-BACKPATH | daemon | `.../registration_rechecked` |
| J-GC-DUPMATCH | daemon | `.../duplicate_registration` |
| J-GC-DIRTY | daemon: SURVIVED the `EpisodeGC|GC|Operator` selector; cmd `TestEpisodeGC` (whole daemon not run) | killed by `TestEpisodeGCFailSafeAndReclaims` (cmd) |
| J-GC-NOCREDS | daemon | `.../unavailable` |
| J-GC-FINALVERIFY | daemon | `.../lock_changed_under_review`, `/replaced_during_scan`, `/not_linked` |
| J-GC-INVALIDEP | daemon | `.../invalid_episode` and others |
| J-DAEMON-NOW | daemon | `.../daemon_wiring` |
| J-LOCK-NLINK | store (whole) | `TestOperatorEpisodeLockIdentity/hardlink` |
| J-LOCK-UID | store (whole) | `.../uid` |
| J-LOCK-PERM | store (whole) | `.../permissions`, `TestOperatorEpisodeLockAcquireRefusals` |
| J-LOCK-SAMEFILE | store (whole) | `.../replaced` |
| J-LOCK-NAME | store (whole) | `.../wrong_key` |
| J-LOCK-NOFOLLOW | store (whole) | `TestOperatorEpisodeLockAcquireRefusals` |
| J-LOCK-EPSLASH | store (whole) | `TestOperatorEpisodeLockAcquireRefusals` |
| J-RR-HELD (drop held fd) | worktree (whole) | `TestRemoveRegisteredLendsDescriptor` |
| J-ADD-INHERIT (`ExtraFiles` off) | worktree (whole) | `TestRemoveRegisteredLendsDescriptor` |
| **J-RR-FORCE** | worktree (whole), daemon selector, cmd `TestEpisodeGC` | **SURVIVED** (N1) |
| J-FILES-VERIFY | cmd | `TestOperatorFilesReconcile/no_mcp` |
| J-WT-PREVERIFY | cmd | `TestOperatorBorrowedWorktreeGitProcess/*` (confounded: `StartCollisions` already red, B1) |
| J-WT-POSTVERIFY | cmd | `TestOperatorBorrowedWorktreeGitProcess/during_git` |
| J-CLI-SAMEDB | cmd | `TestEpisodeGCFailSafeAndReclaims` |
| J-CLI-ROOTBIND | cmd | `TestEpisodeGCFailSafeAndReclaims` |
| J-CTL-LIVE-COMMENT (must-live control) | daemon | survived as intended |
| J-CTL-LIVE-MSG (must-live control, error text) | store | KILLED (error wording is pinned, N4) |
| must-die control | `J-GC-UNRESOLVED` / `J-GC-LOOKUPERR` | killed as required |

Controller ledger `plan-scratch/ledger-ctl-named.jsonl`: 102 lines, 101 KILLED, 1 SURVIVED (`D-LOCK-NO-NONBLOCK`), `ctl-named-mutants.rc` = `named rc=0`. Planner's ids vs mine: no overlap (`mutants.json` has 102 ids, mine are all `J-*`).

KILLED sample I re-ran from the planner/controller ledger on head: J01 (store, `TestOperatorCatalogueNullRefresh/null`), J05 (cmd, `TestOperatorStartCandidateClosesCredentials`), D-LOCK-FLOCK-NONBLOCK (store, `TestOperatorEpisodeLockContendsWithoutWaiting`), D-GC-RELEASE (daemon, `/gc_first`), D-GC-LOCK-FAIL (daemon, `/activation_first`), N12 (cmd, `IdentityWindows/N12`), D-GC-REMOVE-GIT (worktree, `TestRemoveRegisteredLendsDescriptor`). All seven confirmed KILLED. D-LOCK-NO-NONBLOCK SURVIVED, confirmed.

## Measurements

- `git diff -M --numstat b5f41458 c4fb23f6 -- cmd host docs app`: +1240 / -160 = **1,400**, 18 files (matches the plan).
- `go vet ./cmd/... ./host/...` clean (38 s wall); `GOOS=linux go vet` of the four packages rc=0; `go test -run '^$' ./...` compiles clean.
- My `-race -run 'Operator|EpisodeGC|GC|Lock|Remove'`: store 19.97 s ok, worktree 2.73 s ok, daemon 11.34 s ok, cmd FAIL (B1 plus load-timeouts), `judge/racecmd.txt`.
- Controller gate-head: store plain/race, worktree plain/race empty fail sets, 28.7/59.1 s and 25.6/3.4 s; cmd-plain 600 s and daemon-race 797 s with load-attributed differences as above. Compare-store and compare-worktree legs rc=0.

## Recommendation

Fix B1 (restore a real refusal assertion for the held and missing-lock arms), run the whole `cmd/ailang-worldd` plain leg so the red is visible before landing, and add a locked-worktree arm for `RemoveRegistered` (N1). The production diff itself (M1 to M4p) is, in my own measurement, sound and mutation-tight on the live GC path; the sprint fails on its test edits, not on its production change.
