VERDICT: PASS 91

# Row 141 evaluator r1 (iteration 238), judging `git diff 3c76e99..191d4e5` (M1 588838d, M2 475eed2, M3 191d4e5)

Evaluator worktree `.wt-world-iter238-eval` at 191d4e5: `git diff --quiet` is clean and nothing is committed.
Mutants were run on a `git archive 191d4e5` copy at `~/.ailang/state/world-iter238-row141/eval/mut`. Every file was restored, and `diff -r` against a fresh archive is empty.
Raw gate outputs are in `~/.ailang/state/world-iter238-row141/eval/`. The mutant scripts are `eval/bin/mut.py` and `eval/bin/multi.py`.

## Score (100 points, 70 to pass)
| Area | Pts | Notes |
|---|---|---|
| AC coverage vs design + plan (P-C1..7, P-D1..8 honoured) | 28/30 | All of AC1.1–AC3.1 are present and green. The 1 MiB lock bound (P-D3 b) has no test. |
| Test quality, non-vacuity, mutation | 22/25 | 31 of my mutants killed and 5 survived (all low severity, listed below). Behavioural reverts of each milestone fail every test named for that milestone. |
| Code quality / standards | 14/15 | Small and readable, with fail-closed defaults. The default path is unchanged (AC1.6 test file byte-unchanged). |
| Security | 13/15 | No blocking hole. There is a Lstat→Open race on the lock (finding N1). |
| Docs (QUICKSTART §11, help, website) | 9/10 | Accurate against the flags. The recipe's `set -euo pipefail` is an interactive-shell hazard (N4). |
| Gates | 5/5 | All green, with counts below. **Total 91/100.** |

## Blocking findings
None.

## Non-blocking findings (each with a reproduction)
- **N1. Lock read TOCTOU can hang episode construction while holding `w.mu`.**
  - `checkLockCoverage` calls `Lstat` and then `os.Open`. A FIFO swapped in between makes `Open` block, which stalls every episode's handler construction.
  - This needs an agent able to mkfifo, for example through an exec profile that runs agent code, plus a tight race. Construction happens only on the first build.
  - Fix: open with `O_RDONLY|O_NONBLOCK|O_NOFOLLOW`, then `f.Stat()`.
  - Repro: read `host/daemon/workspace_layout.go` lines 'fi, err := os.Lstat(lock)' through 'os.Open(lock)'.
- **N2. The 1 MiB lock bound is untested.** Mutant E8 (`case len(data) > lockMaxBytes:` → `case false:`) SURVIVED.
  - Repro: `python3 eval/bin/mut.py E8 host/daemon/workspace_layout.go 'case len(data) > lockMaxBytes:' 'case false:' '^(TestWorkspace.*)$'`
- **N3. Writability is only tested with owner-writable modes.** E10 (mask `0o222` → `0o200`) SURVIVED, so a group- or other-writable entry (for example `0o464`) has no test row. The code is correct today.
- **N4. QUICKSTART §11 recipe, `set -euo pipefail`.** Pasted into an interactive shell, it stays set. When the guard's `exit 1` fails the pipeline, pipefail plus `-e` closes the operator's terminal.
  - Suggest wrapping the recipe in `bash -c '…'` or `( … )`.
  - This is a docs issue only. The attended S7 verbatim run is still pending.
- **N5. Flag-unset "byte-identical" has no "absent stays absent" arm.** E25 (unset arm `MkdirAll(link)`) SURVIVED. AC2.4 covers dir, empty dir and symlink, but not a missing registry.
- **N6. A symlink to an older snapshot path is refused with a misleading message.** If the operator moves `SNAP` to a new path, an existing episode HOME link to the old path is refused as "not empty and was not provisioned by --workspace-package-cache", although World did make it.
  - This is acceptable as fail-closed: the documented refresh keeps the same path.
  - Consider noting it in §11's line table.
- **N7. Minor survivors.**
  - E9: `ailang.toml` in the snapshot could be a directory, because the `IsRegular` check is untested.
  - E16: the package count would accept a deeper `.../ailang.toml`.
- **N8. Small differences from the executor's report.** It reports the daemon at 23 SKIP; I got 24 SKIP and 447 PASS of 471 RUN. These are environment-dependent skips, not a failure.

## Mutant table (my own, none listed by the executor)
| # | Edit | Result / killing test |
|---|---|---|
| E1 | flag-set: `case fi.IsDir()` + `RemoveAll` (any dir) | KILLED, RefusesUnprovisionedRegistry/non-empty_dir |
| E2 | keep any symlink (drop `readlinkIs`) | KILLED, RefusesUnprovisioned/foreign_symlink |
| E3 | relink a foreign symlink (remove + link) | KILLED, RefusesUnprovisioned/foreign_symlink |
| E4 | unset arm also removes an empty dir | KILLED, UnsetLeavesHomeAlone/empty_dir_untouched |
| E5 | unset arm no-op (keeps stale symlink) | KILLED, UnsetLeavesHomeAlone/stale_symlink |
| E6 | lock `Lstat` → `Stat` (follows symlink) | KILLED, CoversLock/lock_is_a_symlink |
| E7 | invalid-JSON lock → `return nil` (fail-open) | KILLED, CoversLock/invalid_JSON |
| E8 | 1 MiB bound → `false` | **SURVIVED** (N2) |
| E9 | covered `IsRegular` → `st != nil` | **SURVIVED** (N7) |
| E10 | writable mask `0o222` → `0o200` | **SURVIVED** (N3) |
| E11 | file record drops content hash | KILLED, DigestIsStamped (own recomputation) |
| E12 | digest includes a DIR-itself record | KILLED, DigestIsStamped |
| E13 | drop the `Clean(rel)==rel` check | KILLED, ModuleRootGrammar/tools/ |
| E14 | sandboxRoot drops `IsDir` | KILLED, SymlinkOrMissingIsR8/regular_file (line text) |
| E15 | state-dir containment → `false` | KILLED, StartupChecks/inside_the_state_dir |
| E16 | package count accepts deep tomls | **SURVIVED** (N7) |
| E17 | lock read from the worktree, not the module root | KILLED, CoversLock (3 arms) |
| E18 | uid-0 check only when n>0 (moved semantics) | KILLED, StartupChecks/uid_0 |
| E19 | symlink entry skipped, not refused | KILLED, StartupChecks/symlink_inside |
| E20 | allow a symlinked module root resolving inside epRoot | KILLED, SymlinkOrMissingIsR8/symlink_to_a_real_dir |
| E21 | episodeRefusal printed via the generic line | KILLED, RefusesUnprovisioned (exact line) |
| E22 | digest not passed to the handler | KILLED, DigestIsStamped |
| E24 | digest records not `\n`-terminated | KILLED, DigestIsStamped |
| E25 | unset arm `MkdirAll(link)` | **SURVIVED** (N5) |
| E26 | skip one registry lock entry | KILLED, CoversLock/one_registry_entry_absent |
| C1 | `--workspace-episode-module-root` last-wins | KILLED, Section11 (wiring DeepEqual) |
| C2 | accept empty `--workspace-module-root` | KILLED, Section11 (P-D6) |
| C4 | remove the P-C3 §10 cut | KILLED, Section10 ("2 serve lines, want 1") |
| R1 | snapshot greet returns "ho" (RIG test-arm strength) | KILLED, TestSeToolsPackageCache… (test op ok:false) |

## Spot-checks of the executor's mutants: all CONFIRMED
MR-MKDIR, PC-SORT, PC-RESERVED, PC-LOCK-TRAVERSAL, PC-RUN-UNSTAMPED and HELP-DROP were killed by the claimed tests (CI). PC-NOLINK-RIG ("cache not found") and MR-CWD-RIG ("no Ailang.Check handler") were killed under RIG.

## Non-vacuity per milestone (behavioural revert, applied at 191d4e5)
- **M1** (sandboxRoot forced to ".", resolveModuleRoots accepts all, pre-store check off): all 5 M1 tests FAIL.
- **M2** (resolvePackageCache, link and lock no-op; stamp, reserved key, run stamp and pre-check removed): all 6 M2 CI tests FAIL.
- **M3** (`main.go` at 3c76e99): TestQuickstartSection11FlagsMatchTheCLI FAILS on help and on the flags.

## Gates (at 191d4e5, packages run one at a time)
- `go vet ./...`: rc 0.
- daemon `-race -count=1 -v`: rc 0. 471 RUN, 447 PASS, 24 SKIP, 0 FAIL (187 s). Both RIG tests SKIP, as expected without `WORLD_TOOL_AILANG_BIN`.
- broker: rc 0. 627 RUN, 604 PASS, 23 SKIP, 0 FAIL (168 s). The base flake did not appear.
- cmd: rc 0. 115 RUN, 112 PASS, 3 SKIP, 0 FAIL.
- `verify_ail.sh`: rc 0. "16 required identities verified, 40 named tests pass, 444 se-tools named tests pass", equal to base.
- RIG pair: rc 0. 2 RUN, 2 PASS, 0 SKIP.

## CI reality
- Each RIG mutant has a CI companion that kills the same mutant:
  - MR-CWD and MR-POLICY-ROOT → AC1.2.
  - PC-NOLINK → AC2.2.
  - The RIG run stamp is also covered in CI by AC2.5(c).
- Linux portability: the checks are pure `Lstat` `Perm()` reads, `syscall.Mkfifo` exists on Linux, `geteuid` is a seam, and the state dir and root are canonicalized. I found no darwin-only assumption.

## Security read
- **Writes into the snapshot.** It is reached only through the HOME symlink. The mode bits make `pkg_docs` writes fail with EACCES, and AILANG FS writes are confined to the sandbox. A process with the same uid could still `chmod u+w`; that is design risk R2, accepted and documented.
- **Deleting operator data.** World never deletes content. The flag-set path removes only an empty dir. The unset path removes only a symlink, and never follows it (E5 and the executor's FOLLOWS-LINK mutant).
- **Redirecting the link.** A foreign link is refused (E2, E3).
- **Module root.** It only narrows the sandbox. Symlinked, missing or `..` roots are refused, and sandboxRoot runs again on every `registry()` bind, so a later symlink swap is caught on the next bind. Overrides are operator-only, so one episode cannot reach another's worktree.

## Docs
- QUICKSTART §11 matches the flags and is bound by AC3.1. The §10 cut adds one `strings.Cut` at `### 11. ` and leaves §10's checks intact (C4 proves the cut is needed).
- `website/docs/reference/cli.md`'s help block is byte-equal to `ailang-worldd help`.
