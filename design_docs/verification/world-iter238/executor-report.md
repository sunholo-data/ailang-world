# Row 141 executor report (iteration 238)

Branch `sprint/row141-workspace-project-layouts`, base `c17e002`. Gate binary v0.41.0
(`AILANG_BIN`), tool binary v0.52.1 (`WORLD_TOOL_AILANG_BIN`, RIG only). GOCACHE and TMPDIR stayed
outside the tree; raw outputs are under `~/.ailang/state/world-iter238-row141/exec/`.

## Milestones

| M | Commit | Gates (all rc 0) |
|---|---|---|
| M1 | `588838d` | vet 0; compile fence (`go test ./... -run '^$'`) 0; `-race -v` daemon 438 RUN / 135 PASS / 22 SKIP / 0 FAIL (ok 148 s); broker 627 / 203 / 19 / 0 (ok 145 s); cmd 114 / 70 / 3 / 0. Named M1 tests 5 top-level, 34 RUN with subtests, all PASS. |
| M2 | `475eed2` | vet 0; fence 0; `-race -v` daemon 471 / 141 / 23 / 0 (ok 147 s); broker 627 / 203 / 19 / 0 (ok 146.8 s); cmd 114 / 70 / 3 / 0. RIG after M2: see below. |
| M3 | this commit (the report ships with it) | vet 0; fence 0; `verify_ail.sh` rc 0 (`16 required identities verified, 40 named tests pass, 444 se-tools named tests pass`, equal to base); cmd `-race -v` 115 / 71 / 3 / 0; full `go test -race -count=1 ./...` rc 0, 27 `ok` lines, no FAIL (daemon and broker code unchanged since M2). RIG at the end: see below. |

`git diff 3c76e99 -- host/daemon/workspace_test.go` is empty (AC1.6: `TestWorkspacePolicyRenderedPerEpisode`
byte-unchanged). The only edit to an existing test is the P-C3 one-liner in
`cmd/ailang-worldd/quickstart_exec_test.go`.

## RIG (`WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.52.1/ailang`)

After M2 and again at the end, identical:

    === RUN   TestSeToolsModuleRootResolvesBareImports
    --- PASS: TestSeToolsModuleRootResolvesBareImports (2.15s)
    === RUN   TestSeToolsPackageCacheResolvesRegistryImports
    --- PASS: TestSeToolsPackageCacheResolvesRegistryImports (2.21s)

2 `=== RUN`, 2 `--- PASS`, 0 `--- SKIP`. Arm (i) of each asserts the measured failure first
(`LDR001 module not found: b`; `cache not found`), so the pass is not vacuous.

## Mutation table

Every mutant was applied in place on the intended code, the named test run with `-count=1 -v`, then the
file restored from a byte copy (`cmp` identical, printed as `REVERTED`). "RUN" counts are `=== RUN` lines.

| Mutant | Killing test | Failure line (trimmed) | Reverted |
|---|---|---|---|
| MR-DOTDOT | TestWorkspaceModuleRootGrammar | `New = <nil>, want a config StartupError containing "--workspace-module-root"` (subtests `..` and `../x`) | yes |
| MR-DUP | same | `New = <nil>, want ... containing "a second time"` (mapped twice) | yes |
| MR-POLICY-ROOT | TestWorkspaceModuleRootIsTheSandbox | `rendered policy = ... want (fs_sandbox = .../ep1/tools)` | yes |
| MR-CWD | same | `read ran in fake:.../ep1, want the module root .../ep1/tools` | yes |
| MR-NOSYMLINK | TestWorkspaceModuleRootSymlinkOrMissingIsR8 | `registry = [...], want empty` (symlink arms) | yes |
| MR-MKDIR | same | `registry = [...], want empty` (missing arm; see note 1) | yes |
| MR-OVERRIDE-IGNORED | TestWorkspaceEpisodeModuleRootOverrides | `registry(ep2) has no Workspace.Read handler` | yes |
| MR-SILENT | TestWorkspaceModuleRootNeedsWorkspaceRoot | `New = <nil>, want ... containing "need --workspace-root"` | yes |
| MR-DEFAULT | existing TestWorkspacePolicyRenderedPerEpisode | `registry(ep1) is empty` | yes |
| MR-CWD-RIG | TestSeToolsModuleRootResolvesBareImports (RIG) | `registry(ep1) has no Ailang.Check handler` | yes |
| PC-DIRS-ONLY | TestWorkspacePackageCacheStartupChecks | writable file row: `New = <nil>, want ... "is writable"` | yes |
| PC-SYMLINK | same | symlink row: error is `... neither a directory nor a regular file`, want `is a symlink` | yes |
| PC-EMPTY | same | zero-packages row: `New = <nil>, want ... "holds no package"` | yes |
| PC-ROOT-OK | same | uid 0 row: `New = <nil>, want ... "uid 0"` | yes |
| PC-NOLINK | TestWorkspacePackageCacheLinkedIntoEpisodeHome | `Lstat(.../.ailang/cache/registry) = ... no such file or directory, want a symlink` | yes |
| PC-SILENT-DELETE | TestWorkspacePackageCacheRefusesUnprovisionedRegistry | `registry = [...], want empty` on all 3 plants (see note 2) | yes |
| PC-KEEP-EMPTY | same | `registry(ep1) is empty with an empty registry dir in place` | yes |
| PC-UNSET-REMOVES-DIR | TestWorkspacePackageCacheUnsetLeavesHomeAlone | `keep.txt = "", open ...: no such file`; `empty registry dir after startup: <nil>, lstat ...` | yes |
| PC-UNSET-FOLLOWS-LINK | same | `the symlink's target was touched: "", open .../target/keep.txt: no such file` | yes |
| PC-STAMP-ALWAYS | TestWorkspacePackageCacheDigestIsStamped | `Workspace.Read result carries package_cache with the flag unset` | yes |
| PC-DIR-HASH | same | `startup digest = ..., the test's own recomputation = ...` | yes |
| PC-SORT | same | same digest-mismatch line (the `util` / `util-x` pair) | yes |
| PC-RUN-UNSTAMPED | same | `Ailang.Run result = map[admitted:true ...], want admitted with package_cache <digest>` | yes |
| PC-RESERVED | same | `a forged package_cache key = <nil>, want an error naming the reserved key` | yes |
| PC-NOLOCK-COVERAGE | TestWorkspacePackageCacheCoversLock | `registry = [...], summaries 1, want empty and 0` | yes |
| PC-LOCK-TRAVERSAL | same | traversal subtest: `registry = [...], summaries 1, want empty and 0` | yes |
| PC-NOLINK-RIG | TestSeToolsPackageCacheResolvesRegistryImports (RIG) | `with the cache, ai_check a.ail: ... registry package acme/util cache not found at ...` | yes |
| HELP-DROP | TestQuickstartSection11FlagsMatchTheCLI | `serve --help does not document --workspace-package-cache` | yes |
| WIRING | same | `serve built Config {... WorkspaceModuleRoot:tools ... WorkspacePackageCache: ...}` | yes |
| (extra) P-C3-CUT | existing TestQuickstartSection10FlagsMatchTheCLI | `QUICKSTART §10 holds 2 serve lines, want 1` (proves P-C3: the cut is needed) | yes |
| (extra) EMPTY-FLAG | TestQuickstartSection11FlagsMatchTheCLI | `--workspace-module-root "": exit 0` (P-D6) | yes |

No survivors. All 31 mutants killed.

Notes on how the plan's mutants fired:
1. MR-MKDIR: with the mutant the missing directory is created, so the registry forms. The
   `registry = ..., want empty` assertion fires before the not-created assertion is reached. Same test, same arm.
2. PC-SILENT-DELETE: the empty-registry assertion fires first on each plant, so the planted-file-survives
   assertion is never reached. Same test.
3. Mutant text differs from the plan's literal one-liner where my code shape differs (the plan's text assumed
   a shape): PC-LOCK-TRAVERSAL is `if validPkgRel(rel) {` -> `if true {`; PC-SORT is `sort.Strings(records)` ->
   `_ = sort.Strings` (keeps the import compiling; the plain delete fails to build); PC-UNSET-REMOVES-DIR is
   `if err == nil && fi.Mode()&os.ModeSymlink != 0 { return os.Remove(link)` -> `if err == nil && fi != nil {
   return os.RemoveAll(link)`; PC-NOLINK adds `if w.packageCache != "" { return nil }` as the first flag-set
   statement; HELP-DROP renames the `--workspace-package-cache` entry in the `serve flags:` block (the only
   column-2 line); WIRING sets the field to `""`.

## Deviations and decisions beyond the plan

- **AC1.7 / AC2.6 `a.ail`/`c.ail` bodies.** The plan gives `b.ail`, `c.ail`, `t.ail`, `r.ail` verbatim but not
  `a.ail`'s body. `a.ail` (`import b`, `import pkg/acme/util/greet`) returns `"${b()}${greet()}"`: a first draft
  used `++`, which v0.52.1 rejects ("`++` is for lists only").
- **Lock refusal wording.** The bad-lock reasons print `lock <path> <reason>` (`is not a regular file`,
  `is larger than 1 MiB`, `is not valid JSON: ...`, `cannot be read: ...`), as P-D3 says. A malformed
  registry name/version (traversal) prints the ordinary `lock requires <name>@<version>, absent from <DIR>` line.
- **Wiring and the empty-flag arm** (plan step 5 bullets) are in `checkLayoutWiring`, called from
  `TestQuickstartSection11FlagsMatchTheCLI`, so the plan's named test covers them.
- **One extra pre-existing-doc fix:** `website/docs/reference/cli.md` was regenerated from
  `ailang-worldd help` (verified byte-equal to the block), date stamp 2026-10-06. It also picks up row 140's exec
  flags, as the plan said it would.
- **S7 recipe check.** The snapshot recipe from §11 was run verbatim (bash) on a synthetic lock and operator cache
  in the scratch dir: rc 0, snapshot has the one package, 0 writable entries. The attended verbatim run of the whole
  section remains pending, as §11 says.
- **G3 ran per package, sequentially** (daemon, broker, cmd), not 3-concurrent, so the known broker flake did not
  appear. No flake re-runs were needed; `TestRunBoundedHeadTailTimeoutKeepsPartialOutput` and
  `TestExecSrtMCPEndToEnd` were green/skipped at every boundary and in the final full `-race` run.
- **Flake seen once, after the M3 commit.** A plain (non-race) `go test -count=1 ./...` run on the committed
  tree failed `TestExecHandlerTimeoutKillsTheGroup` (`handlers_exec_m2_test.go:398`, `TimedOut:true` with an
  empty head, 0.82 s) in `host/broker`, under the all-packages parallel load. Re-run alone with `-race -count=1`
  three times: 3/3 `ok`. It is a sleep/timing test in exec code this row does not touch (my broker diff is
  `ExecPathOk` plus the `package_cache` stamp), the same class as the known
  `TestRunBoundedHeadTailTimeoutKeepsPartialOutput` flake. The earlier full `-race` run was green.
- Not done, by instruction: U1/U2 upstream issues, `ailang messages send`, the R-A/R-B rows, the iteration log.

## r2 (fix round for judge findings N1-N7, one commit on top of `f680e28`)

Method per finding: apply the judge's mutant to the current suite (SURVIVE), add the row, re-apply (KILL),
restore from a byte copy (`cmp` identical). Mutants run against `^TestWorkspacePackageCache` (32 RUN before, 42 after).

| Finding | Mutant (judge's) | Before: current suite | New row | After: killed by | Reverted |
|---|---|---|---|---|---|
| N2 / E8 | `case len(data) > lockMaxBytes:` -> `case false:` | SURVIVED (ok) | CoversLock `lock over 1 MiB` (valid JSON that would be covered, 1 MiB + padding) | `operator log ... is not valid JSON: unexpected end ..., want "is larger than 1 MiB; rebuild ..."` | yes |
| N3 / E10 | `Perm()&0o222` -> `Perm()&0o200` | SURVIVED | StartupChecks: group-writable file `0o464`, other-writable file `0o446`, group-writable dir `0o575`, other-writable dir `0o557` | 4x `New = <nil>, want ... "is writable"` | yes |
| N5 / E9 | `covered = err == nil && st.Mode().IsRegular()` -> `err == nil && st != nil` | SURVIVED | CoversLock `ailang.toml of a locked package is a directory` | `registry = [...], summaries 1, want empty and 0` | yes |
| N7 / E16 | count `len(segs) == 4 && segs[3] == ...` -> `len(segs) >= 4 && last == ...` | SURVIVED | StartupChecks `ailang.toml one level too deep`, `ailang.toml is a directory` | `New = <nil>, want ... "holds no package"` | yes |
| N5b / E25 | flag-unset arm `MkdirAll(link)` | SURVIVED | UnsetLeavesHomeAlone `missing registry stays missing` | `flag unset: the absent registry path was created (lstat: <nil>)` | yes |
| N1 | (code change; see below) | n/a | `TestWorkspacePackageCacheLockFIFOIsRefusedPromptly` (FIFO `ailang.lock`, refusal within 2 s) | mutant drop `O_NONBLOCK`: `a FIFO as ailang.lock hung registry construction (no refusal within 2s)`; mutant drop the Fstat regular check: `registry size 0, log ...` assertion fails | yes |
| N6 | (code change) mutant `case fi.Mode()&os.ModeSymlink != 0:` -> `case false:` | n/a | foreign-symlink arm now expects the new line | exact-line mismatch (old "is not empty ..." line printed) | yes |

Code changes:
- **N5, N7:** no code change was needed. The code already required a regular file for coverage (`IsRegular`)
  and counted only exactly 4-segment regular `ailang.toml` files; the survivors were missing tests, now added.
  CoversLock now builds the snapshot, runs the row's setup, then locks it down, so a row can shape the snapshot.
- **N1:** `checkLockCoverage` no longer does `Lstat` then `Open`. One `OpenFile(O_RDONLY|O_NONBLOCK|O_NOFOLLOW)`
  (ELOOP -> "is not a regular file"; not-exist -> covered/no lock), then `f.Stat()` must be a regular file. The
  1 MiB bound is enforced once, on the bytes read (`LimitReader` + `len(data) > lockMaxBytes`), not also by Fstat
  size, so E8 stays a real mutant. A pre-planted FIFO is refused by the new code and was refused by the old
  `Lstat` too, so that test cannot fail on the old code (the race itself is not deterministic); it fails on the
  two mutants above, which reintroduce the blocking open or drop the type check. Runs on Linux and macOS.
- **N6:** an existing symlink to a different path is refused with `... <link> is a symlink to <old>, not to
  --workspace-package-cache <DIR>; remove it if the snapshot moved` (still refused, never deleted). §11's
  operator-line table has the new row.
- **N4:** the §11 recipe is wrapped in `( set -euo pipefail; ... )`. Re-run verbatim (extracted from the doc)
  on the synthetic lock: rc 0, 1 package, 0 writable entries, `shell survived rc=0`; missing-package arm
  prints `not in operator cache: acme/nope/1`, `shell survived rc=1`.

Gates (one package at a time, `-race -count=1 -v`, all rc 0, 0 FAIL): `go vet ./...` 0; daemon 481 RUN / 142 PASS /
23 SKIP (ok 151 s); broker 627 / 203 / 19 (ok 145 s); cmd 115 / 71 / 3. RIG pair: 2 `=== RUN`, 2 `--- PASS`, 0 SKIP.
No flake appeared in this round.
