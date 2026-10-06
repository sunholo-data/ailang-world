# Sprint plan: w-workspace-project-layouts (iteration 238, row 141)

**Base:** branch `sprint/row141-workspace-project-layouts` at `3c76e99` (design r3).
**Authority:** [approved design](w-workspace-project-layouts.md). It went through quorum r1 (BLOCKED 3/3) and r2 (1 pass, 2 rejects), and r3 is the carve-out that applies the r2 fixes verbatim. This plan implements that design. Where the design is wrong or silent, the plan says so in **§ Planner corrections and decisions**. Each such item is labelled `P-C*` (a correction backed by evidence) or `P-D*` (a planner decision the controller may strike). Do not settle anything else silently: if you hit a new ambiguity, stop and write it in the iteration notes.
**Executor:** one commit per milestone, in order M1 → M2 → M3. Each commit must be green at its own boundary. Stage named files only, never `git add -A`. Do not push.

## Planner corrections and decisions (read first)

| Id | Kind | What | Evidence |
|---|---|---|---|
| P-C1 | correction | **AC1.1's killing row for MUT-MR-DOTDOT is `../x` (and `..`), not `tools/../../x`.** With the `..` check dropped, `tools/../../x` is still refused by the `filepath.Clean(REL) == REL` check, because `Clean` gives `../x`. `../x` is already clean, so it passes `Clean`, and with the `..` check dropped only it reaches acceptance. Keep the `tools/../../x` row, but assert the kill on `../x`. | by construction of the two checks |
| P-C2 | correction | **AC1.7 cannot use `a.ail`.** `a` imports `pkg/acme/util/greet`, so with the module root alone (no package cache) it fails `registry package acme/util cache not found`, not `passed`. AC1.7 uses a third module, `c.ail` (bare `import b` only). AC2.6 keeps `a.ail`. | Measured on a scratch copy of V16 (`~/.ailang/state/world-iter238-row141/plan/syn2`): root sandbox + `tools/c.ail` → `LDR001 module not found: b`; `tools/` sandbox + `c.ail` → `ok:true`, `passed:true`; `tools/` sandbox + `a.ail` (empty HOME) → `registry package acme… cache not found` |
| P-C3 | correction | **The §10 binding test is NOT unaffected (the Conflict Surface says it is).** `section10ServeLines` (`cmd/ailang-worldd/quickstart_exec_test.go:77`) cuts at the §10 heading and reads **to EOF**. It then requires exactly 1 `/tmp/ailang-worldd serve ` line. A §11 serve line makes it 2, and `TestQuickstartSection10FlagsMatchTheCLI` goes red. M3 adds one line after the Cut, `section, _, _ = strings.Cut(section, "### 11. ")`, the same pattern §9's test uses (`quickstart_se_test.go:41`). This is the only edit to an existing test. | `sed -n 70,90p cmd/ailang-worldd/quickstart_exec_test.go` |
| P-C4 | correction | **Ailang.Run results do not pass through `stamp`.** `executeRun` builds `runResult` and returns `json.Marshal(result)` (`handlers_ailang.go:563–598`). Only the policy-tool ops (read, write, ai_check, discover, cli) are stamped. As designed, an `ailang-run` result would not name the package cache it resolved against, and V15 shows that `run` does resolve registry imports. See P-D4. | `grep -n 'stamp(' host/broker/*.go` → only `handlers_ailang.go:440,465` |
| P-C5 | correction | **`execPathOk` is unexported** (`host/broker/execargs.go:153`), so `host/daemon` cannot call it. M1 adds an exported one-line wrapper, `func ExecPathOk(p string) bool { return execPathOk(p) }`, so both features share one grammar (design (a1)). No behaviour change. | grep |
| P-C6 | correction | **The design's M2 operator lines do not fit the existing `registry()` error path.** That path prints `workspace tools unavailable for episode %q: %v` (`workspace.go:304–307`). The M2 lines have their own prefixes and must appear **exactly once**. M2 adds a typed `*episodeRefusal{line string}`, which `registry()` prints verbatim (`errors.As`). The M1 module-root line does use the generic path, as the design says. | `sed -n 296,316p host/daemon/workspace.go` |
| P-C7 | measured (closes a design UNMEASURED) | **The `test` op under the snapshot works.** On the scratch copy, `{"op":"test","path":"t.ail"}` passes with the snapshot `HOME` and fails with an empty `HOME`. Write the test as a `tests [((), "hi")]` clause on a function. An inline `test "x" { … }` block fails for an unrelated reason: `test` copies the file to `/tmp/ailang-namedtest-*` and the block cannot see the module's own functions. `run --policy -- r.ail` prints `hi` (rc 0) with the snapshot `HOME`. | same scratch dir |
| P-D1 | decision (provenance) | **Add `package_cache` to the reserved-key list** in `Execute` (`handlers_ailang.go:460`), unconditionally. Without this, with the flag unset, a policy-tool response that carried a `package_cache` key would reach the record unstamped and look like World provenance. Real v0.52.1 responses never carry the key, so real results stay byte-identical (AC2.5 unset arm). | — |
| P-D2 | decision (behaviour, fail-closed) | **The startup walk also refuses any entry that is neither a directory nor a regular file** (FIFO, socket, device). Reading a FIFO for the digest would hang startup. Error: `<path> is neither a directory nor a regular file`. | — |
| P-D3 | decision (behaviour, fail-closed) | **The lock coverage check fails closed on a bad lock.** `ailang.lock` is agent-writable (design R-B), and World reads it with operator privileges. The check therefore: (a) `Lstat`s it and requires a regular file (a symlink or FIFO refuses); (b) reads at most 1 MiB; (c) refuses invalid JSON; (d) treats a `registry` entry whose `name/version` is not exactly 3 non-empty segments with no `.` or `..` segment as **not covered**, so a traversal like `../../x` never reaches `Stat`. Every case takes the coverage-refusal path with its own reason (§M2). | — |
| P-D4 | decision (provenance; the controller may strike) | **Stamp `ailang-run` results too.** Add `PackageCache string \`json:"package_cache,omitempty"\`` to `runResult` and set it in `executeRun`. Flag unset → omitted → byte-identical. Tested by AC2.5 arm (c). If struck, delete M2 step 7 and arm (c). | P-C4 |
| P-D5 | decision (spec precision) | **The digest's byte form.** Records are built for every entry **except DIR itself**. Each record uses the slash-separated relpath: a directory is `rel + "\x00d"`, a regular file is `rel + "\x00f\x00" + lowercase-hex(sha256(content))`. Records are sorted with `sort.Strings` (byte order of the **full relpath**, not walk order: `a/b` sorts after `a-c`), and each is terminated by `"\n"`. The digest is `"sha256:" + hex(sha256(concat))`. The AC2.5 test recomputes the digest with its own recursive `os.ReadDir` walker and its own sort. | the design says "in byte order", which walk order is not |
| P-D6 | decision (CLI) | **`--workspace-module-root ""` is refused at flag parse** (`fs.Func`, exit 1, `is empty (use . for the worktree root)`), because an empty `Config.WorkspaceModuleRoot` means "unset". AC1.1's `""` row is the Config-level `--workspace-episode-module-root ep1=` form. | — |
| P-D7 | decision (placement) | **The CLI flags, usage and `serve --help` land in M3**, together with their binding test (AC3.1, MUT-HELP-DROP). M1 and M2 land the `Config` fields and the daemon behaviour, driven by tests through `Config`. | design Conflict Surface lists `main.go` once |
| P-D8 | decision (test rig) | **The RIG tests drive the episode's registry handler directly**: `realToolBin(t)` (needs only `WORLD_TOOL_AILANG_BIN`) and `d.workspace.registry("ep1")[effect].Execute(…)`, as `TestWorkspaceSandboxIsTheEpisodeNotTheRoot` does. They do not use the full MCP `seRig`. They still live in `host/daemon/setools_e2e_test.go`, as the design names. | — |

**Controller actions, not executor actions.** The executor does **not**:
- file the U1/U2 upstream issues;
- run `ailang messages send`;
- edit `design_docs/world-mission.md` (the R-A/R-B rows);
- touch the iteration log.

AC3.2 is evidence that the controller produces after the land.

## Gates and measured base (planner, at `3c76e99`, 2026-10-06, macOS)

Run every command from the worktree root, with this environment:

```sh
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN=$HOME/.pinned-ailang/ailang          # GATE binary, AILANG v0.41.0
export TOOL=$HOME/.pinned-ailang-tools/v0.52.1/ailang  # TOOL binary, AILANG v0.52.1 — only as WORLD_TOOL_AILANG_BIN
```

Do not mix them up. `AILANG_BIN` is always v0.41.0. `WORLD_TOOL_AILANG_BIN` is always the v0.52.1 tool binary.

| # | Gate | Base result |
|---|---|---|
| G1 | `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh` | **rc 0**, 31 s. `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 444 se-tools named tests pass` |
| G2 | `go vet ./...` | **rc 0**, 1 s (warm) |
| G3 | `AILANG_BIN=$HOME/.pinned-ailang/ailang go test -race -count=1 -v ./host/daemon/ ./host/broker/ ./cmd/ailang-worldd/` | **rc 1 at base**, 157 s wall (3 packages concurrently): daemon ok 148 s, cmd ok 7.8 s, **broker FAIL** on `TestRunBoundedHeadTailTimeoutKeepsPartialOutput` (`handlers_capture_test.go:122: partial stdout = ""`, 0.30 s). `=== RUN` lines 1152 in that run. **Sequential per-package re-run** (`-race -v`): daemon rc 0, 151 top-level / 411 total `=== RUN`, 130 PASS, 21 SKIP (cached; uncached it took 148 s); broker **rc 0**, 146 s, 222 / 627, 203 PASS, 19 SKIP; cmd rc 0, 6 s, 73 / 114, 70 PASS, 3 SKIP. So the base is green when the packages run one at a time. `TestExecSrtMCPEndToEnd` **SKIPs** under this env (no srt rig), so the known row-153 flake cannot show here. See the base-flake note below. |
| G4 | `WORLD_TOOL_AILANG_BIN=$TOOL AILANG_BIN=$HOME/.pinned-ailang/ailang go test -count=1 ./host/daemon/ -run '^(TestSeToolsModuleRootResolvesBareImports\|TestSeToolsPackageCacheResolvesRegistryImports)$' -v` | **rc 0 and VACUOUS at base**: `testing: warning: no tests to run`, 0 `=== RUN`. The tests do not exist yet. This is exactly why you must count `=== RUN` lines (want 2) and `--- PASS` lines (want 2, 0 SKIP) after M2. |

In the table cell, `\|` is Markdown's escaped pipe. Type it as a bare `|`. **Never write a Go `-run` alternation as `A\|B`: it runs zero tests and still exits 0.** Always pass `-count=1` when you want timing or a fresh result: the planner's second `-v` daemon run came back cached (`0 s`).

**Base flake (record it, do not chase it).** `TestRunBoundedHeadTailTimeoutKeepsPartialOutput` failed once under the 3-package concurrent `-race` load. Re-run it alone: `go test -race -count=1 -run '^TestRunBoundedHeadTailTimeoutKeepsPartialOutput$' ./host/broker/` (planner re-run, alone, 5/5 rc 0 at ~1.75 s each). It is a 300 ms timing test in code this row does not touch. If it reds at a milestone boundary, re-run it alone. If it passes alone, record "base flake, row-153 class" and continue. `TestExecSrtMCPEndToEnd` (row 153, `host callback timed out`) gets the same treatment wherever it runs.

**At every milestone boundary**, run:
1. `go vet ./...`
2. `go test ./... -run '^$'` (compile fence for `_test.go`; `go build` is not one)
3. G3
4. the milestone's named tests with `-count=1 -v`, counting `=== RUN`
5. G1 once, at the M3 boundary (no `.ail` changes, so it must equal the base)
6. G4 at the M2 and M3 boundaries, requiring `--- PASS: TestSeToolsModuleRootResolvesBareImports` and `--- PASS: TestSeToolsPackageCacheResolvesRegistryImports`, 2 `=== RUN`, 0 `--- SKIP`
7. at M3 closure, `AILANG_BIN=$HOME/.pinned-ailang/ailang go test -race -count=1 ./...`

## Executor hazards

1. **Read-only fixtures in `t.TempDir()`.** `TempDir` cleanup fails on a `chmod a-w` tree, and the test then fails with `unlinkat … permission denied`. Every helper that makes a tree read-only must register, **right after it chmods**, a `t.Cleanup` that walks the tree and runs `os.Chmod(p, 0o755)` on directories and `0o644` on files. Cleanups run LIFO, so this one runs before `TempDir` removal, because `TempDir` registered first. Use `filepath.WalkDir` and ignore errors. Do **not** follow symlinks.
2. **Mode bits on Linux CI and macOS.** The CI runner is non-root. The checks are pure `Lstat().Mode().Perm() & 0o222` reads, so they behave the same on both. Set modes explicitly with `os.Chmod` (`0o555` dirs, `0o444` files) after `os.WriteFile`, and never rely on umask. Tests that expect **acceptance** set the `geteuid` seam to `func() int { return 1000 }` (`t.Cleanup` restores it), so a root container does not flip them. The uid-0 row sets it to `0`. The seam is a package var: those tests must not call `t.Parallel()`.
3. **`TMPDIR` and `GOCACHE` stay outside the repo.** At base, `go env GOCACHE` is `~/Library/Caches/go-build` and `$TMPDIR` is `/var/folders/…`. Never export a `TMPDIR` or `GOCACHE` inside the worktree: a temp dir in the tree reds the boundary scanners. Put scratch work under `~/.ailang/state/world-iter238-row141/exec/`. `/tmp` gets wiped.
4. **No `.snap/` Go copies, `*.orig` files or mutation backups in the tree.** Apply each mutation in place, run the test, revert with `git checkout -- <file>`, then prove it with `git diff --quiet -- <file> && echo IDENTICAL`. Before each commit, `git status --porcelain` must list only the milestone's named files.
5. **Flags unset must be byte-identical.**
   - AC1.6: `TestWorkspacePolicyRenderedPerEpisode` stays **byte-unchanged**. Check with `git diff 3c76e99 -- host/daemon/workspace_test.go` and confirm that function's hunk is absent.
   - AC2.5 unset arm: no `package_cache` key.
   - `runResult` uses `omitempty`.
   - With the flag unset, `linkPackageCache` only `Lstat`s a missing path.
   - Do not reorder or reword existing operator lines.
6. **Serial tests.** The fixtures and seams mutate package vars (`geteuid`). None of the new tests calls `t.Parallel()`.
7. **ErrorLog capture.** Pass `ErrorLog: &buf` (`bytes.Buffer`; `registry()` is called synchronously in these tests) and assert the **exact** line count with `strings.Count(buf.String(), "\n")` and the exact text.

## M1 — module root: a subdirectory sandbox (one commit)

**Files:**
- `host/broker/execargs.go`: add `ExecPathOk`, P-C5.
- `host/daemon/daemon.go`: `Config` fields and startup.
- `host/daemon/workspace.go`: struct fields, `registry`, `episodeHandler`.
- new `host/daemon/workspace_layout.go`: the M1 helpers; M2 extends the file.
- new `host/daemon/workspace_layout_test.go`: the CI tests.
- `host/daemon/setools_e2e_test.go`: AC1.7.

**Steps:**
1. `Config` gains `WorkspaceModuleRoot string` and `WorkspaceEpisodeModuleRoots []string`, with doc comments in the style of `ExecEpisodeProjects`.
2. `New()`: directly after the `DBPath` check, **before** `store.Open`:
   ```go
   if cfg.WorkspaceRoot == "" && (cfg.WorkspaceModuleRoot != "" || len(cfg.WorkspaceEpisodeModuleRoots) > 0) {
       return nil, &StartupError{Stage: StageConfig, Detail: "--workspace-module-root and --workspace-episode-module-root need --workspace-root"}
   }
   ```
   Inside the `if cfg.WorkspaceRoot != ""` block, after the examples dir:
   ```go
   def, byEp, err := resolveModuleRoots(cfg.WorkspaceModuleRoot, cfg.WorkspaceEpisodeModuleRoots, root)
   if err != nil { return nil, &StartupError{Stage: StageConfig, Detail: "the module root is refused", Err: err} }
   ```
   Then store them in `workspace.moduleRoot` and `workspace.episodeModuleRoot`, two new fields on `workspaceTools`.
3. `workspace_layout.go`:
   - `checkModuleRoot(flag, value, rel, root string) error`. It refuses `!broker.ExecPathOk(rel)` ("must be relative, with no leading / or -, and no .. segment") and `filepath.Clean(rel) != rel` ("is not clean"). It then calls `broker.RenderEpisodePolicy(filepath.Join(root, rel))`; an error there is a refusal (this catches `a"b`).
   - `resolveModuleRoots(def string, pairs []string, root string) (string, map[string]string, error)`. It checks `def` when it is non-empty. For each pair, `strings.Cut(pair, "=")` must succeed (`is not EP=REL`) and `episodeIDPattern` must match EP (`episode %q is outside the episode grammar`). A repeated EP is refused (`maps episode %q a second time`, the shape of `exec_startup.go:162–163`). Then `checkModuleRoot` runs. Every error names the flag and the raw value.
   - `(w *workspaceTools) sandboxRoot(episodeID, epRoot string) (string, error)`:
     - `rel := w.moduleRoot`; the episode's override wins; `if rel == "" { rel = "." }`;
     - `want := filepath.Join(epRoot, rel)`;
     - `got, err := filepath.EvalSymlinks(want)`: an error gives `module root %q: %v`; `got != want` gives `module root %q: resolves through a symlink to %q`;
     - `os.Stat(want)` must be a directory, else `module root %q: not a directory`.
     - **Never `MkdirAll`.**
4. `registry()`: after `episodeRoot` succeeds, call `sandbox, err := w.sandboxRoot(episodeID, epRoot)`. On error, print the existing generic line (`ailang-worldd: workspace tools unavailable for episode %q: %v\n`) and return an empty registry. Then call `w.episodeHandler(episodeID, sandbox)`. `w.execHandler(episodeID, epRoot)` is **unchanged**.
5. `episodeHandler(episodeID, sandbox string)`: rename the parameter, then:
   - `broker.RenderEpisodePolicy(sandbox)`;
   - `Root: sandbox`;
   - the cache key is `cached.root == sandbox`; store `episodeTool{root: sandbox, …}`;
   - `policyPath` and `cacheDir` are unchanged.

**Tests** (CI unless marked), in `workspace_layout_test.go` (all use `newWSFixture`, `fakeToolBin`, `mustWSDaemon`/`newWSDaemon`):

| AC | Test | Asserts |
|---|---|---|
| AC1.1 | `TestWorkspaceModuleRootGrammar` (table through `newWSDaemon`, with `ep1/tools` and `ep1/tools/sub` created) | Accepted: `tools`, `tools/sub`, `.` (as the default **and** as `ep1=.`). Refused with a `*StartupError`, `Stage == StageConfig`, the message naming the flag: `ep1=` (empty REL), `/abs`, `..`, `../x`, `tools/../../x`, `-x`, `tools/`, `a"b`, `EP1=tools` (EP outside the grammar), `ep1` (no `=`), and `ep1=tools` given twice. Refused rows leave `f.db` absent (`os.Stat` → `IsNotExist`). |
| AC1.2 | `TestWorkspaceModuleRootIsTheSandbox` | With `WorkspaceModuleRoot: "tools"` and `ep1/tools` created: (a) `policies/ep1.toml` == `RenderEpisodePolicy(<root>/ep1/tools)` byte for byte; (b) `registryNames(reg)` == `wantWorkspaceEffects` (the 9 names); (c) `reg[Workspace.Read].Execute(… {"op":"read","path":"x"})` returns content `fake:<root>/ep1/tools` (the fake prints `pwd -P`); (d) `f.summaries == 1`. |
| AC1.3 | `TestWorkspaceModuleRootSymlinkOrMissingIsR8` (subtests, fresh fixture each) | `ep1/tools` as (i) a symlink to `ep1/real` (a real dir), (ii) a symlink to `../ep2/tools` (create `ep2/tools`), (iii) a regular file, (iv) missing. Each arm gives: an empty registry, `summaries == 0`, `dispatches == 0`, exactly one ErrorLog line that contains `workspace tools unavailable for episode "ep1": module root "tools":`. Arm (iv): `ep1/tools` still does not exist afterwards. |
| AC1.4 | `TestWorkspaceEpisodeModuleRootOverrides` | `WorkspaceModuleRoot: "tools"`, `WorkspaceEpisodeModuleRoots: {"ep2=."}`, `ep1/tools` created. A read through ep1 returns `fake:<root>/ep1/tools`; a read through ep2 returns `fake:<root>/ep2`. |
| AC1.5 | `TestWorkspaceModuleRootNeedsWorkspaceRoot` | Each of `{WorkspaceModuleRoot:"tools"}` and `{WorkspaceEpisodeModuleRoots:{"ep1=tools"}}` without `WorkspaceRoot` gives a `*StartupError` with `StageConfig` and a message containing `need --workspace-root`, and `f.db` does not exist. Control: the same config **with** `WorkspaceRoot` and the tool bin starts, and `f.db` exists, which proves the instrument sees an opened store. |
| AC1.6 | existing `TestWorkspacePolicyRenderedPerEpisode` | **Unchanged**: it must pass byte-unchanged. |
| AC1.7 | `TestSeToolsModuleRootResolvesBareImports` in `setools_e2e_test.go` (**RIG**: `realToolBin`, SKIP in CI) | Fixture in `ep1/tools`: `b.ail` = `module b\n\nexport func b() -> string {\n  "b"\n}\n`; `c.ail` = `module c\n\nimport b (b)\n\nexport func c() -> string {\n  b()\n}\n` (P-C2). (i) No module root: `Ailang.Check` `{"op":"ai_check","path":"tools/c.ail"}` → the decoded inner `stdout` has `check.passed == false` and an error with `code == "LDR001"` and `message == "module not found: b"`. (ii) A fresh daemon (fresh fixture) with `WorkspaceModuleRoot: "tools"`: `{"op":"ai_check","path":"c.ail"}` → `ok == true` and `check.passed == true`. |

**Commit:** `row 141 M1: --workspace-module-root subdirectory sandbox (AC1.1–AC1.7)`, with the `Co-Authored-By` trailer the controller gives.

## M2 — read-only registry package cache (one commit)

**Files:**
- `host/daemon/daemon.go`;
- `host/daemon/workspace.go`;
- `host/daemon/workspace_layout.go`;
- `host/broker/handlers_ailang.go` (`AilangToolConfig`, the handler field, `stamp`, the reserved keys, `runResult`/`executeRun` per P-D4);
- `host/daemon/workspace_layout_test.go`;
- `host/daemon/setools_e2e_test.go`.

**Steps:**
1. Add `Config.WorkspacePackageCache string`. Extend the M1 pre-store check's condition with `|| cfg.WorkspacePackageCache != ""`, and change its Detail to `--workspace-module-root, --workspace-episode-module-root and --workspace-package-cache need --workspace-root`. AC1.5 asserts only the substring `need --workspace-root`, so it stays green.
2. In the root block, after the module roots:
   - call `resolvePackageCache(cfg.WorkspacePackageCache, root, stateDir)`;
   - on error, return `&StartupError{Stage: StageConfig, Detail: "the --workspace-package-cache directory is refused", Err: err}`;
   - on success, set `workspace.packageCache` and `workspace.packageCacheDigest`;
   - write one line to `resolveErrorLog(cfg.ErrorLog)`: `ailang-worldd: workspace package cache <canonical>: <n> packages, read-only, <digest>\n`.
3. `workspace_layout.go`: add `var geteuid = os.Geteuid` and `resolvePackageCache(dir, canonicalRoot, stateDir string) (canonical, digest string, n int, err error)`. It refuses, in this order:
   - (a) `filepath.Abs` + `EvalSymlinks` fail, or `Stat` is not a directory: `--workspace-package-cache %q is not a directory`;
   - (b) `broker.CheckPolicyOutsideRoot(canonical, canonicalRoot)` fails: `&WorkspaceStateError{What: "the package cache is inside --workspace-root", …}`;
   - (c) the cache is lexically within `stateDir`: `&WorkspaceStateError{What: "the package cache is inside the state directory", …}`. Write a local `withinDir(path, base)` with `filepath.Rel`, as `broker.within` does. Do not call `CheckPolicyOutsideRoot` here, because it needs the state dir to exist;
   - (d) `geteuid() == 0`: `the daemon runs as root (uid 0), where mode bits do not stop writes`;
   - (e) a `filepath.WalkDir(canonical, …)` over every entry, **including DIR itself**, using `d.Info()` (Lstat semantics). A symlink gives `%s is a symlink; the package cache must hold none`. A non-dir, non-regular entry gives P-D2's line. `Perm()&0o222 != 0` gives `%s is writable (mode %v); the package cache must be read-only (chmod -R a-w)`. Collect the P-D5 records and count packages: relpaths of exactly 4 slash segments whose last segment is `ailang.toml`, on a regular file;
   - (f) `n == 0`: `%s holds no package (<ns>/<name>/<ver>/ailang.toml)`.

   Then compute the digest per P-D5, as its own function `packageCacheDigest(records []string) string`.
4. `(w *workspaceTools) linkPackageCache(episodeID, cacheDir string) error`. Let `link := filepath.Join(cacheDir, ".ailang", "cache", "registry")`.
   - **Flag unset** (`w.packageCache == ""`): if `os.Lstat(link)` is a symlink, `os.Remove(link)` and return. Anything else (absent, a dir, an error) returns nil, untouched.
   - **Flag set**:
     - `os.MkdirAll(filepath.Dir(link), 0o700)`, then `os.Lstat(link)`;
     - absent → go on to symlink;
     - a symlink whose `os.Readlink == w.packageCache` → return nil;
     - a directory with `len(os.ReadDir) == 0` → `os.Remove(link)`, then go on;
     - **anything else** → return `&episodeRefusal{line: fmt.Sprintf("ailang-worldd: workspace package cache refused for episode %q: %s is not empty and was not provisioned by --workspace-package-cache; clear it manually", episodeID, link)}`. This covers a non-empty dir, a file, and a foreign symlink;
     - finally `os.Symlink(w.packageCache, link)`.
     - Never `RemoveAll`.
5. `(w *workspaceTools) checkLockCoverage(episodeID, sandbox string) error`: returns nil when the flag is unset or `<sandbox>/ailang.lock` is absent. Otherwise apply P-D3 and decode `struct{ Packages []struct{ Name, Version, Source string } \`json:"packages"\` }`. For each entry in lock order with `Source == "registry"`, the entry is covered only if `rel := Name + "/" + Version` has a valid shape and `os.Stat(filepath.Join(w.packageCache, rel, "ailang.toml"))` is a regular file. The first miss returns `&episodeRefusal{line: fmt.Sprintf("ailang-worldd: workspace package cache does not cover episode %q: lock requires %s@%s, absent from %s; rebuild the snapshot per QUICKSTART §11", episodeID, Name, Version, w.packageCache)}`. The P-D3 bad-lock cases use the same prefix, `…does not cover episode %q: lock %s <reason>; rebuild the snapshot per QUICKSTART §11`.
6. `episodeHandler`: after the existing `MkdirAll` loop, call `w.linkPackageCache(episodeID, cacheDir)`, then `w.checkLockCoverage(episodeID, sandbox)`. Either error returns before the policy is rendered and before any summary. Then pass `PackageCacheDigest: w.packageCacheDigest` in `AilangToolConfig`. `registry()`: `var r *episodeRefusal; if errors.As(err, &r) { fmt.Fprintln(w.errLog, r.line) } else { <existing generic line> }`.
7. `handlers_ailang.go`:
   - `AilangToolConfig.PackageCacheDigest string` and handler field `packageCacheDigest`;
   - `stamp`: `if h.packageCacheDigest != "" { resp["package_cache"], _ = json.Marshal(h.packageCacheDigest) }`;
   - the reserved list at `:460` adds `"package_cache"` (P-D1);
   - `runResult` gains `PackageCache string \`json:"package_cache,omitempty"\``, and `executeRun` sets `result.PackageCache = h.packageCacheDigest` (P-D4).

**Test fixtures** (in `workspace_layout_test.go`):
- `readOnlyPackageCache(t, dir string, pkgs ...string)`: for each `ns/name/ver`, it writes `ailang.toml` (`[package]\nname = "ns/name"\nversion = "ver"\nedition = "1"\n`) and one `.ail` file. It then chmods files to `0o444` and dirs to `0o555`, bottom-up, DIR last, and registers the hazard-1 cleanup.
- `writeLock(t, path string, entries ...lockEntry)`: marshals `{"schema":"ailang.lock/v1","schema_version":"1.0.0","packages":[{name,version,source,content_hash:"sha256:"+64 zeros}]}` from a Go struct.

Place DIR at `f.base/pkgcache`. That is outside `f.root` and `f.stateDir`, since `newWSFixture` puts them at `base/ws` and `base/state`.

| AC | Test | Asserts |
|---|---|---|
| AC2.1 | `TestWorkspacePackageCacheStartupChecks` (table through `newWSDaemon`; geteuid seam 1000 except in its own row) | Precondition (kimi r2): before any startup check, a walk of the accepted tree counts **0** entries with `Lstat` mode `&os.ModeSymlink` and **>0** entries in total. **Accepted:** one package (`acme/util/0.1.0`) → daemon starts. ErrorLog holds `workspace package cache <DIR>: 1 packages, read-only, sha256:`. Two direct `resolvePackageCache` calls return the same digest. **Refused** (each a `*StartupError`, `StageConfig`, message containing `--workspace-package-cache` and the row's substring): not a dir (a file), `not a directory`; inside root (`f.root/pkgs`), `inside --workspace-root`; inside state (`f.stateDir/pkgs`), `inside the state directory`; one writable **file** (`0o644`), `is writable`; one writable **inner dir** (`0o755`), `is writable`; DIR itself `0o755`, `is writable`; a symlink inside (an `os.Symlink` to a read-only sibling, created before the chmod), `is a symlink`; zero packages (dirs only), `holds no package`; uid 0 (seam), `uid 0`; no `--workspace-root`, `need --workspace-root`, with `f.db` absent. P-D2 row: a FIFO (`syscall.Mkfifo`) inside → `neither a directory nor a regular file`. |
| AC2.2 | `TestWorkspacePackageCacheLinkedIntoEpisodeHome` | After `registry("ep1")`: `Lstat(<state>/cache/ep1/.ailang/cache/registry)` is a symlink, and `Readlink` == canonical DIR. The fake's logged `env` contains `HOME=<state>/cache/ep1\n`. `policies/ep1.toml` == `RenderEpisodePolicy(<root>/ep1)`. Calling `d.workspace.linkPackageCache("ep1", <state>/cache/ep1)` again returns nil, and the link is unchanged (idempotent keep arm). |
| AC2.3 | `TestWorkspacePackageCacheRefusesUnprovisionedRegistry` (subtests) | (i-a) a non-empty dir at `link` holding `sunholo/x/9.9.9/ailang.toml`; (i-b) a regular file at `link`; (i-c) a symlink to another dir. Each gives: an empty registry; `summaries == 0`; exactly one ErrorLog line equal to the design's line with `<link>` substituted (after the startup line, so assert on the lines after it); the planted file/dir/link **survives** byte-identical (the file's content read back; the foreign link's `Readlink` unchanged). (ii) An **empty** dir at `link` → replaced by the symlink, and the registry forms. |
| AC2.4 | `TestWorkspacePackageCacheUnsetLeavesHomeAlone` (flag unset; subtests) | (a) a non-empty `registry` dir and (b) an empty one: both untouched, and the registry forms. (c) A pre-existing symlink (to a writable dir `f.outside/target` holding `keep.txt`): removed (`Lstat` → not exist), its target and `keep.txt` intact, and the registry forms. |
| AC2.5 | `TestWorkspacePackageCacheDigestIsStamped` | (a) Flag set: the `Workspace.Read` result decodes to a map with `package_cache` == the startup digest == the test's own recomputation (recursive `os.ReadDir`, P-D5 records, `sort.Strings`, `\n`-terminated). Use a fixture with a `a/b`-vs-`a-c` naming pair so that walk order ≠ byte order: `acme/util/0.1.0` and `acme/util-x/0.1.0`. (b) Flag unset: the result has **no** `package_cache` key. (c) P-D4: a run through `Ailang.Run` with a run-capable fake (below), with the flag set, has `package_cache` == the digest; with the flag unset, the key is absent. (d) P-D1: a fake whose dispatch answers `{"ok":true,"package_cache":"forged"}` makes `Execute` return an error naming the reserved key, with the flag **unset**. |
| AC2.7 | `TestWorkspacePackageCacheCoversLock` (subtests; `WorkspaceModuleRoot:"tools"`; snapshot holds `acme/util/0.1.0`) | (i) lock `[acme/util@0.1.0 registry]` → the registry forms. (ii) lock `[acme/util registry, sunholo/oauth@0.1.0 registry]` → empty registry, `summaries == 0`, exactly one line `ailang-worldd: workspace package cache does not cover episode "ep1": lock requires sunholo/oauth@0.1.0, absent from <DIR>; rebuild the snapshot per QUICKSTART §11`. (iii) lock `[acme/util registry, local/thing@0.1.0 path]` → forms. (iv) no lock → forms. P-D3 arms: (v) `{` → refused with `does not cover episode "ep1": lock`; (vi) name `../esc` version `1.0.0`, with `f.base/esc/1.0.0/ailang.toml` present → refused; (vii) `ailang.lock` as a symlink to a valid lock outside → refused. |
| AC2.6 | `TestSeToolsPackageCacheResolvesRegistryImports` in `setools_e2e_test.go` (**RIG**) | V16 fixture in `ep1/tools`: `a.ail`, `b.ail`, `ailang.toml`, `ailang.lock` (from `writeLock`), plus `t.ail` = `module t\n\nimport pkg/acme/util/greet (greet)\n\nexport func hi() -> string tests [((), "hi")] {\n  greet()\n}\n` and `r.ail` = `module r\n\nimport std/io (println)\nimport pkg/acme/util/greet (greet)\n\nexport func main() -> () ! {IO} {\n  println(greet())\n}\n`. The snapshot holds `acme/util/0.1.0/{ailang.toml,greet.ail}`; `ailang.toml` has `[exports]\nmodules = ["acme/util/greet"]`, and `greet.ail` = `module acme/util/greet\n\nexport func greet() -> string {\n  "hi"\n}\n`. These are exactly the V16 files, `cat`ed in the planner's scratch dir. (i) `WorkspaceModuleRoot:"tools"`, no cache: `ai_check a.ail` → inner `check.passed == false` with a message containing `cache not found`. (ii) A fresh daemon with the cache flag: `ai_check a.ail` → `passed`; `Ailang.CLI` `{"op":"test","path":"t.ail"}` → `ok == true`; `Ailang.Run` `{"path":"r.ail"}` → `admitted == true`, `exit_code == 0`, stdout `hi\n`. The snapshot digest (the test's recomputation) is equal before and after all three. P-C7 measured each arm on the scratch copy. |

**Run-capable fake for AC2.5(c)/(d).** Define a new helper in `workspace_layout_test.go`, `fakeToolBinWith(t, logDir, dispatchJSON string)`. Do **not** edit `fakeToolBin`. Its script equals `fakeToolBin`'s, with two changes:
- the policy-tool dispatch prints `dispatchJSON`;
- a branch `if [ "$1" = "run" ]; then echo 'policy: {"ok":true,"decision":{"ok":true}}' >&2; exit 0; fi` sits before the final `dispatches` line.

That is outcome shape (a) of `composeRunOutcome`. A run payload needs only `{"path":"r.ail"}`. Check `parseRunPayload` for the required keys before writing the arm.

**Commit:** `row 141 M2: --workspace-package-cache read-only snapshot, lock coverage, package_cache stamp (AC2.1–AC2.7)`.

## M3 — CLI flags, docs and binding test (one commit)

**Files:**
- `cmd/ailang-worldd/main.go`;
- `docs/QUICKSTART.md`;
- `website/docs/guides/operating-the-daemon.md`;
- `website/docs/reference/cli.md`;
- `cmd/ailang-worldd/quickstart_exec_test.go` (the P-C3 one-liner only);
- new `cmd/ailang-worldd/quickstart_layout_test.go`.

**Steps:**
1. `runServe`:
   - `--workspace-module-root` via `fs.Func` (P-D6 refuses empty);
   - `--workspace-episode-module-root` via `fs.Func`, appending (repeatable);
   - `--workspace-package-cache` via `fs.String`;
   - pass all three into `daemon.Config`.
2. `usage`:
   - Usage block: add a line `                      [--workspace-module-root REL] [--workspace-episode-module-root EP=REL ...]` and a line `                      [--workspace-package-cache <dir>]`, after `--examples-dir`.
   - `serve flags:` block: three entries, each starting `  --flag` at column 2, after `--examples-dir`. Say:
     - module root: the sandbox and module root are `<workspace-root>/<episode>/REL`; `.` is the default; REL is clean and relative and must already exist as a real directory, or that episode's tools are refused; it is independent of an exec profile's `root`.
     - episode override: repeatable, and it overrides the default.
     - package cache: a read-only snapshot of registry packages (QUICKSTART §11), linked as every episode's package cache; it must be outside `--workspace-root` and the state dir, hold no symlink and nothing writable; serve refuses as uid 0; restart serve after rebuilding it.
3. `docs/QUICKSTART.md`: append `### 11. Projects with a module root and registry packages`. Mark it `**Written in row 141 M3; the attended verbatim run is pending.**` and name the binding test `TestQuickstartSection11FlagsMatchTheCLI`. Content, in order:
   1. the problem in two sentences (bare imports resolve against the sandbox; registry packages resolve under `HOME`);
   2. `session new <episode> --repo <project>`, in the shape of the existing §9 line;
   3. the operator lock step: `(cd $WS/<episode>/tools && $TOOL lock)`. Say plainly that the agent cannot run `lock` and that the lock is an input to the recipe;
   4. the design's snapshot recipe **verbatim** (`set -euo pipefail` plus the `test -d … || { echo "not in operator cache: $p" >&2; exit 1; }` guard), the refresh rule (`chmod -R u+w`, rebuild, restart), and the multi-episode union rule;
   5. exactly **one** serve line, `/tmp/ailang-worldd serve …` (same shape and `\` continuations as §10's), using `--workspace-root`, `--tool-ailang-bin`, `--workspace-module-root tools`, `--workspace-episode-module-root <ep>=tools` and `--workspace-package-cache $SNAP`;
   6. what the agent sees (paths relative to the module root, `policy_digest` and `package_cache` in results);
   7. the three operator lines (module root unavailable, cache refused, lock not covered) with what to do;
   8. one sentence each on R-A (`pkg_docs` still fetches when the flag is unset) and on the exec profile's `root` being independent.
4. `quickstart_exec_test.go`: insert `section, _, _ = strings.Cut(section, "### 11. ")` right after the `if !ok` block of `section10ServeLines` (P-C3). Nothing else in that file changes.
5. `quickstart_layout_test.go`: add `section11ServeLines` (a copy of `section10ServeLines` with the §11 heading) and `TestQuickstartSection11FlagsMatchTheCLI`:
   - `serve --help` documents each of the three flags (regex `(?m)^\s+--flag\b` over the text after `serve flags:`);
   - §11 holds exactly 1 serve line;
   - every flag token on it is documented;
   - the line shows all three new flags;
   - with `--bind not-a-hostport` appended, `run` exits `exitUsage` with `is not host:port`;
   - wiring: with `captureServeConfig`, `serve --db /x/world.db --workspace-root /w --workspace-module-root tools --workspace-episode-module-root a=. --workspace-episode-module-root b=x --workspace-package-cache /p` yields exactly those three Config values;
   - `--workspace-module-root ""` exits `exitUsage` (P-D6).
6. Website:
   - `operating-the-daemon.md` `## Flags`: extend the synopsis code block and add three table rows.
   - `cli.md`: replace the `ailang-worldd` help block with the binary's actual `go run ./cmd/ailang-worldd help` output. It is stale at base, missing row 140's exec flags. Update "captured from a build of `dev` on <date>". No other website edits.
7. **Not the executor's** (controller): U1/U2 issues, the `ailang messages send` note, the R-A/R-B rows in `design_docs/world-mission.md`, and AC3.2.

| AC | Test | Asserts |
|---|---|---|
| AC3.1 | `TestQuickstartSection11FlagsMatchTheCLI` | as step 5 |
| (regression) | `TestQuickstartSection10FlagsMatchTheCLI`, `TestQuickstartSection9…`, attended-box test | still green after §11 is appended (P-C3) |

**Commit:** `row 141 M3: serve flags, QUICKSTART §11 and binding test (AC3.1)`.

## Test plan: AC → test → mutation

Apply each mutation, run the named test with `-count=1 -run '^Name$' -v` (count `=== RUN` ≥ 1), and see it **FAIL** at the named assertion. Then `git checkout -- <file>` and `git diff --quiet -- <file> && echo IDENTICAL`. Record each as `mutant → test → FAIL line → IDENTICAL` in the iteration notes. A mutant that survives is a stop: fix the test, not the mutant.

| AC | Test | Mutant | One-line source edit (file) |
|---|---|---|---|
| AC1.1 | `TestWorkspaceModuleRootGrammar` | MUT-MR-DOTDOT | In `checkModuleRoot`, replace `!broker.ExecPathOk(rel)` with `rel == "" \|\| strings.HasPrefix(rel, "/") \|\| strings.HasPrefix(rel, "-")` (`workspace_layout.go`). The `../x` row fires (P-C1). |
| AC1.1 | same | MUT-MR-DUP | Delete the duplicate-EP `case` line (`workspace_layout.go`). The twice-mapped row fires. |
| AC1.2 | `TestWorkspaceModuleRootIsTheSandbox` | MUT-MR-POLICY-ROOT | `broker.RenderEpisodePolicy(sandbox)` → `broker.RenderEpisodePolicy(filepath.Dir(sandbox))` in `episodeHandler`. Byte-equality fires. (`filepath.Dir` stands in for `epRoot`, which is no longer in scope.) |
| AC1.2 | same | MUT-MR-CWD | In `episodeHandler`'s `AilangToolConfig`, `Root: sandbox` → `Root: filepath.Dir(sandbox)`. The `fake:<…>/ep1/tools` assertion fires. |
| AC1.3 | `TestWorkspaceModuleRootSymlinkOrMissingIsR8` | MUT-MR-NOSYMLINK | In `sandboxRoot`, `if got != want {` → `if false && got != want {`. Symlink arm (i) fires. |
| AC1.3 | same | MUT-MR-MKDIR | Insert `_ = os.MkdirAll(want, 0o755)` as the first line after `want :=` in `sandboxRoot`. Arm (iv)'s not-created assertion fires. |
| AC1.4 | `TestWorkspaceEpisodeModuleRootOverrides` | MUT-MR-OVERRIDE-IGNORED | Comment out the override lookup line in `sandboxRoot`. The ep2 assertion fires. |
| AC1.5 | `TestWorkspaceModuleRootNeedsWorkspaceRoot` | MUT-MR-SILENT | Change the pre-store condition `cfg.WorkspaceRoot == "" &&` → `false &&` (`daemon.go`). The error assertion fires. |
| AC1.6 | `TestWorkspacePolicyRenderedPerEpisode` | MUT-MR-DEFAULT | `if rel == "" { rel = "." }` → `if rel == "" { rel = "src" }`. The existing test fires (empty registry). |
| AC1.7 | `TestSeToolsModuleRootResolvesBareImports` (RIG) | MUT-MR-CWD | as above; arm (ii) gives `LDR001`. Run with G4's env. |
| AC2.1 | `TestWorkspacePackageCacheStartupChecks` | MUT-PC-DIRS-ONLY | Wrap the writable check as `if d.IsDir() && info.Mode().Perm()&0o222 != 0 {`. The writable-file row fires. |
| AC2.1 | same | MUT-PC-SYMLINK | Symlink check `if info.Mode()&os.ModeSymlink != 0 {` → `if false {`. The symlink row fires: it asserts the substring `is a symlink`, and under the mutant the entry falls through to P-D2's `neither a directory nor a regular file` instead. Keep the check order symlink → P-D2 → writable. |
| AC2.1 | same | MUT-PC-EMPTY | `if n == 0 {` → `if false {`. The zero-packages row fires. |
| AC2.1 | same | (uid seam) MUT-PC-ROOT-OK | `if geteuid() == 0 {` → `if false {`. The uid-0 row fires. |
| AC2.2 | `TestWorkspacePackageCacheLinkedIntoEpisodeHome` | MUT-PC-NOLINK | First line of `linkPackageCache`'s flag-set arm: `return nil`. The `Lstat` assertion fires. |
| AC2.3 | `TestWorkspacePackageCacheRefusesUnprovisionedRegistry` | MUT-PC-SILENT-DELETE | Replace the `default`/refusal return with `if err := os.RemoveAll(link); err != nil { return err }` (falls through to the symlink). Planted-file-survives fires. |
| AC2.3 | same | MUT-PC-KEEP-EMPTY | `len(entries) == 0` → `len(entries) < 0`. Arm (ii) fires. |
| AC2.4 | `TestWorkspacePackageCacheUnsetLeavesHomeAlone` | MUT-PC-UNSET-REMOVES-DIR | Flag-unset arm: `if err == nil && fi.Mode()&os.ModeSymlink != 0` → `if err == nil`, with `os.RemoveAll`. The untouched-dir arms fire. |
| AC2.4 | same | MUT-PC-UNSET-FOLLOWS-LINK | Flag-unset arm `return os.Remove(link)` → `t, _ := os.Readlink(link); _ = os.RemoveAll(t); return os.Remove(link)` (one line, `;`-joined). `keep.txt` is gone, and target-untouched fires. |
| AC2.5 | `TestWorkspacePackageCacheDigestIsStamped` | MUT-PC-STAMP-ALWAYS | `if h.packageCacheDigest != "" {` → `if true {` in `stamp`. Arm (b) fires. |
| AC2.5 | same | MUT-PC-DIR-HASH | The directory record `rel + "\x00d"` → `rel + "\x00f\x00" + hex(sha256(nil))`. Arm (a)'s recomputation fires. |
| AC2.5 | same | MUT-PC-SORT (P-D5) | Delete the `sort.Strings(records)` line. Arm (a) fires (the `util` / `util-x` pair). |
| AC2.5 | same | MUT-PC-RUN-UNSTAMPED (P-D4) | Delete `result.PackageCache = h.packageCacheDigest`. Arm (c) fires. |
| AC2.5 | same | MUT-PC-RESERVED (P-D1) | Remove `"package_cache"` from the reserved list. Arm (d) fires. |
| AC2.7 | `TestWorkspacePackageCacheCoversLock` | MUT-PC-NOLOCK-COVERAGE | First line of `checkLockCoverage`: `return nil`. Arm (ii) fires. |
| AC2.7 | same | MUT-PC-LOCK-TRAVERSAL (P-D3) | Drop the shape check (`if !validPkgRel(rel) \|\| …` → `if …`). Arm (vi) fires. Arm (vi) must be built so the traversal **would** succeed: write `f.base/esc/1.0.0/ailang.toml` (DIR is `f.base/pkgcache`) and use lock name `../esc`, version `1.0.0`. Without the shape check, `Join(DIR, "../esc/1.0.0", "ailang.toml")` exists, the episode is covered, and the arm's refusal assertion fires. |
| AC2.6 | `TestSeToolsPackageCacheResolvesRegistryImports` (RIG) | MUT-PC-NOLINK | as above; arm (ii) gives `cache not found`. Run with G4's env. |
| AC3.1 | `TestQuickstartSection11FlagsMatchTheCLI` | MUT-HELP-DROP | Delete the `  --workspace-package-cache <dir>` line from `usage` in `main.go`. The documented-flag assertion fires. |
| AC3.1 | same | MUT-WIRING | Drop `WorkspacePackageCache: *packageCache` from the `daemon.Config` literal. The wiring arm fires. |

## Closure log (planner measurements)

- Base: `3c76e99`; G1 rc 0 (31 s); G2 rc 0; G3 rc 1 (broker flake above; daemon and cmd ok); G4 rc 0, vacuous (0 RUN).
- Scratch measurements for P-C2 and P-C7: `~/.ailang/state/world-iter238-row141/plan/syn2` (outside the repo), with the tool binary v0.52.1 and World's exact child env.
- Broker flake `TestRunBoundedHeadTailTimeoutKeepsPartialOutput`: red once under concurrent 3-package `-race` load. 5/5 green alone; broker package green alone (146 s). Treat as a load flake.
