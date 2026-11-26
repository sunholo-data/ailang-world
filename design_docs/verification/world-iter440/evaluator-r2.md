# Evaluator report, round 2, iteration 440, row 243 V4 sprint S5f2 (head c7c08e5b, base b5f41458)

## Verdict: PASS (77/100), no BLOCKING finding

| Category | Max | R1 | R2 | Basis |
|---|---:|---:|---:|---|
| Correctness | 30 | 23 | 23 | adopted: production code is byte-identical since c4fb23f6 |
| Tests | 30 | 14 | 25 | B1 closed, N1 closed; 5 pts held back, see New/residual |
| Security / fail-closed | 20 | 16 | 16 | adopted (production unchanged) |
| Coding standards | 10 | 8 | 8 | adopted; gofmt on touched files clean (gates-r2/legs.txt) |
| Scope and size | 10 | 6 | 5 | 1,476 changed lines vs cap 1,400 |
| **Total** | 100 | 67 | **77** | |

## Scope of c7c08e5b (verified myself)
`git diff --stat c4fb23f6 c7c08e5b`: 4 files, +66/-14: operator_episode_lock_test.go, operator_files_test.go, operator_worktree_test.go, host/worktree/worktree_test.go. `git diff --quiet c4fb23f6 c7c08e5b -- ':!*_test.go'` = rc 0 (PRODUNCHANGED). Test-only claim holds. I read every hunk.

Fixture timer (operator_episode_lock_test.go): the 2 s `time.AfterFunc` fuse that releases the held lock moved from before the fixture's migration/claim work to after `ClaimOperatorAction`. It still bounds a contending caller by 2 s after the fixture is ready and is still stopped in Cleanup. It fixes a fixture race (fuse could fire during slow setup under load, releasing the lock before the test ran); it hides nothing, since the fuse is only a liveness fallback and the assertions are unchanged.

## Adjudication

| id | verdict | evidence |
|---|---|---|
| B1 | CLOSED | Alone, `go test -count=1 -v -run 'TestOperatorStartCollisions$/^(held|missing_lock)$|TestOperatorFilePublicationAdmission$/^lock_held$' ./cmd/ailang-worldd/`: all three PASS (0.26/0.26/0.49 s). Full groups plus the two load-red tests also pass alone (45 s), and -race of Collisions|Admission|EpisodeLock* ok (60 s). The arms again call the PRODUCTION helpers (`reserveOperatorStartWorktree`, `publishOperatorStartFiles`) directly and assert `errors.Is(err, store.ErrOperatorConflict)` plus unchanged `observeProvision` state, no writes/boots/stops. Intent vs base: at base "held" meant another process holds the flock, "missing_lock" meant the sidecar is gone. Under M2's contract the helper no longer opens the lock, it receives a borrowed description, so the equivalent refusal is "the description lent is not the canonical file now" (held: original inode renamed away and a different description opened at the canonical path; missing_lock: sidecar unlinked after borrow; lock_held: another episode's description lent). That is a faithful translation, and the refusal is the production `Verify`, not the test helper. |
| N1 | CLOSED | `TestRemoveRegisteredRefusesLockedOrDirty` (host/worktree), real git. My mutant R2-FORCE (`--force` added in `RemoveRegistered`, worktree.go:118) is KILLED by `/dirty` (host/worktree). |
| N2 | OPEN, accepted | CLI `--credentials-db` override for non-default workspaces skips all; documented and pinned in R1. Unchanged. |
| N3 | OPEN, accepted | persistent `.operator-episode-<ep>.lock` sidecars from GC; undocumented clutter. Unchanged. |
| N4 | NOT-A-DEFECT (brittleness note) | tests pin error wording; round 2 pins `episode lock pathname changed` / `wrong episode lock` too, consistently. |
| N5 | NOT-A-DEFECT | `host/daemon/daemon.go` gofmt is pre-existing, untouched. |
| N6 | CLOSED/attributed to load | see whole-package comparison. |

## My own mutants for the B1 arms and N1 (ledger `judge/r2/ledger.txt`)
Each restored with `git checkout`; worktree status clean at the end.
- R2-SAMEFILE (`SameFile` check disabled in `store.verify`): `TestOperatorStartCollisions/held` KILLED (cmd/ailang-worldd; refusal degrades to "mint_start differs" mismatch). `lock_held` survives this one by design, because the name check refuses it first.
- R2-NAMECHECK (`l.file.Name() != path` disabled): `TestOperatorFilePublicationAdmission/lock_held` KILLED (cmd).
- R2-WT-ENTRY (entry `Verify` at operator_worktree.go:39 removed): `StartCollisions/held` and `/missing_lock` KILLED (cmd).
- R2-FILES-VERIFY (`Verify` at operator_files.go:274 removed): `FilePublicationAdmission/lock_held` KILLED (cmd).
- R2-FORCE: KILLED (host/worktree, above).
No arm is vacuous: every B1 arm dies to at least one production mutation.

## Whole-package comparison
Controller `gates-r2/legs.txt`: gofmt, build, vet, linux-vet rc 0; worktree-plain fails=0 (compare rc 0); cmd-plain 13 fails (1,122 s). Base set has 11 (the known e2e reds, all 11 present at head). The two extras are `TestPresetInitSeToolsModeAndModel` (`git init` "timed out after 30s") and `TestPresetStartReadRefusals/profile` ("preset base must name a local branch", a git fixture under load, load average was 18). Both pass when re-run alone (`judge/r2/targets.txt`, ok 45 s). Neither file is touched by the sprint. `race-changed-cmd` rc=1 was `TestOperatorFilePublicationAdmission` Git 30 s timeouts (load); the same tests pass under -race alone. The three B1 tests are not in the head fail set. Own runs: whole `host/worktree` ok 9.6 s, whole `host/store` ok 24.4 s.
Adopted from R1 (production unchanged): M1/M2/M3/M4p kill tables, 33-mutant ledger, inertness grep, daemon results, race costs. Re-measured: scope, B1 arms, N1, package comparisons above.

## Size
`git diff -M --numstat b5f41458 c7c08e5b -- cmd host docs app` = +1,304/-172 = **1,476** (cap 1,400, +76 over, all test-only fix-forward of a judge finding; 80 changed lines in this round). Acceptable as a repair of a regression the sprint caused, but over the cap, so Scope drops 1 point. Round-2 test cost is not on the production path.

## New / residual
No new blocking finding. Residual: cmd whole-package gates are only reliable on a quiet rig (30 s Git timeouts under load), so the controller's compare leg reports false reds; re-run-alone is required. The cmd `-race` group for Admission/Collisions costs about 60 s under load (not baselined at base), worth a look against the 8-minute CI bound.
