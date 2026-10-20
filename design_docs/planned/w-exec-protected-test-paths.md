# w-exec-protected-test-paths — protected test paths for `workspace-exec`: refused to World's write tools, denied in the sandbox, and proven intact in every exec result

**Status**: PLANNED; **M1c equivalence premise REFUTED in iteration 281 (2026-10-19, V16) → [r4] detection redesign (iteration 282, V17–V22): Git is a candidate generator with a provably wider pathspec, World's mirror is the authority. Quorum r1 BLOCKED → [r4.1] (V23, ignored-baseline snapshot, shared phase timeout); quorum r2 BLOCKED on one surface (the phase budget) → **[r4.2]** controller narrow-refinement carve-out applying the reviewers' fixes verbatim (2 500 ms phase bound, a named construction bound, M6 gated on V24). M1c may proceed; M6 waits for V24.** **[V25, iteration 289] The v0.54.15 tool pin closes P6 for `ailang-write`/`ailang-edit` upstream (resolved-name match + multi-link refusal, ailang#1873); §4.3.2's World pre-check is held on D-WORLD-86.** Unicode semantics are unchanged.  Designer draft 2026-10-17 (mission iteration 277, designer role; Agent `opus`). Queue row 156. **r2** (2026-10-17): quorum round 1 BLOCKED on one objection (index flags and a moved HEAD); revised with V13 (§12).
- **Direction:** D-WORLD-81 = A (attended 2026-10-16): order 156 → 150 → 151 → v1.0.0 tag. Row filed attended 2026-10-06 as row 140's verification gap 1.
- **Extends:** row 140 `design_docs/planned/w-workspace-exec-toolchain-effect.md` (§4.3 profile, §4.4 rendered settings, §4.5 handler and record, §4.7 replay) and row 184 `design_docs/planned/w-toolchain-worktree-delta.md` (the `worktree_delta` result field, git-status basis).
- **Base:** ailang-world `dev` `b7bae72`.
- **Rig:** macOS 26.6.2 arm64 (Apple M4 Max), node v26.0.0, srt **0.0.78** (`cli.js` sha256 `3c3092bd…d96e`, the row-140 pin), tool binary AILANG **v0.52.1** (`~/.pinned-ailang-tools/v0.52.1/ailang`, the `RoleTool` pin). **Linux is not measured** (no bubblewrap on the rig). Every Linux statement below is sourced from the pinned srt's README and labelled so. M0 measures Linux on CI before any production code lands, as row 140 M0 did.

**Clauses**: 5 (an honest, checkable record of what an agent's run proves) plus 3 (the protection is part of the brokered effect's record).

**Estimate**: ≈ 1.5 executor days (M0–M7), each milestone ≤ ~30 min of executor work and ≤ ~250 LOC.

**Probe conditions**:
- Every V-row was run first-party on 2026-10-17. Scripts and raw outputs are in `design_docs/verification/world-row156-design-2026-10-17/` (`q*.sh` / `q*.py` / `*.go.txt` and the matching `*.out.txt`).
- The three real checkouts were only **read** (`git ls-files`, `git grep`, `git ls-tree`, a read-only file walk). Every test run and every write ran in a throwaway sibling copy under `/Users/voightkampff/dev/sunholo-data/.scratch-iter277-row156/` (APFS `cp -c` or `git clone --shared`), deleted afterwards.
- This doc claims no `.ail` text (§4.8 says why none is needed).

## 1. Why (measured)

- **P1 — An agent can rewrite the tests that judge it.**
  - `workspace-exec` renders `allowWrite = [worktree, episode exec cache]`. `denyWrite` holds only World's deny names and srt's default write paths (V1).
  - World's `ailang-write`/`ailang-edit` refuse only `.git` and the fixed deny names (`.github`, `.pi`, `.claude`, `.ailang`, `.gitmodules`, `.gitattributes`) (V1).
  - With no protection, every write arm succeeded: modify, create, delete and rename a test file, both through srt and through policy-tool (V5 rendering NONE, V7 control block).
  - So "tests passed through World" proves only that *some* tests passed.
- **P2 — The three M4 projects keep their tests in two different shapes (V2).**
  - **TwilightGame (TS):** all 246 test files are under the root `tests/`. 244 are `*.test.ts`.
  - **sunholo-platform (Python + TS):** `backend/tests/` (317 files), `cli/tests/` (39), `services/browser/tests/` (2). Plus 184 `*.test.ts(x)` in `__tests__/` folders under `frontend/`. Four `conftest.py` files.
  - **ailang (Go):** 1 924 `*_test.go` files sit **next to the code they test**. There are also 208 `*.golden` files, 337 files under `testdata/` folders, and 103 files under the root `tests/`.
  - So directory patterns alone cannot protect Go. A profile needs **basename patterns** (`*_test.go`, `*.golden`) that match at any depth.
- **P3 — Targeted test runs write nothing under these paths (V3, V4).**
  - Each project ran one targeted command in a throwaway copy with World's child environment. Results:
    - `vitest run <file>`: 3 tests passed. Nothing was written outside `node_modules` (1 file there).
    - `pytest -p no:cacheprovider <file>` (the M4 profile's form): 5 passed, 6 skipped, nothing written.
    - `go test ./internal/lexer/`: `ok`, nothing written.
  - Without `-p no:cacheprovider`, pytest writes `.pytest_cache/` at its rootdir (`cli/`), not under `cli/tests/`.
  - **Snapshots:** 0 test files in either TS project use a snapshot assertion. The control `expect(` matches 245 and 182 files.
  - **Go goldens:** the ailang repo rewrites goldens only behind 4 opt-in `-update-*` flags. A profile's flag allowlist does not list them, so the agent cannot pass them (row 140 §4.3).
  - World already sends `__pycache__` to the exec cache (`PYTHONPYCACHEPREFIX`, V1).
  - So protecting test paths breaks none of the measured runs.
- **P4 — The two matchers World must drive do not read globs the same way. The row's example `**/*.test.ts` is a trap.**
  - **World's write tools use AILANG's `fileguard.MatchDenyWrite` (V6).** A pattern ending `/**` is a root-relative directory prefix. Any other pattern is `path.Match`ed against the whole path and against the base name.
    - So `**/*.test.ts` protects exactly one directory level. `a.test.ts` and `src/deep/x.test.ts` were both **written** under it (V7).
    - `tests/**` does not protect `cli/tests/x.py` (V7).
    - `*.test.ts` (a base-name glob) protects every depth (V7).
  - **srt on macOS (Seatbelt) reads `**` as any number of directories.** `<wt>/**/*_test.go` refused a new root-level `c_test.go` as well as `pkg/b_test.go` (V5).
  - **srt on Linux (bubblewrap) takes literal `denyWrite` paths only.** A trailing `/**` is dropped (the directory is denied). **Any other glob is skipped** (V8, srt README line 438, not measured here).
  - So World must define its own small pattern grammar, which means the same thing in every matcher. It must refuse the forms that do not.
- **P5 — Expanding globs at render time does not stop creation (V5, rendering LITERAL).**
  - The row proposed expanding globs at render time, because Linux has no globs. So the existing matches were rendered as literal paths. It was measured on macOS, where both forms can be compared.
  - Existing files were protected (modify, delete, rename-away and hard link were all refused).
  - But **six creation arms were written**:
    - a new `pkg/b_test.go`;
    - a root `c_test.go`;
    - `cli/sub/conftest.py`;
    - `src/bar.test.ts`;
    - `mv pkg/a.go pkg/z_test.go`;
    - `pkg/B_TEST.GO`.
  - A new `zz_test.go` with a `TestMain` that never calls `m.Run()` is enough to make `go test` print `ok`. On Linux this cannot be prevented by the sandbox. It can only be **detected**. That is why this design records a per-call integrity check (§4.5) and does not rely on the sandbox alone.
- **P6 — World's write tools can be steered around a protected pattern by a symlink (V7).**
  - policy-tool checks the **lexical** path (`h.root.Rel(p)`) and then opens it through `os.Root`, which follows a symlink that stays inside the root (V6).
  - With `tests/**` and `*_test.go` protected, `lnk/new5.py` (`lnk → tests`) and `alias.txt` (`→ pkg/a_test.go`) were **written**.
  - An agent can make such a link with any profiled command that runs project code, since the sandbox does not deny creating `lnk`. A repository can also check one in.
  - AILANG's FS effect (`ailang-run` with `FS`) uses the same lexical check (V6).
- **P7 — `git status`, limited to the protected pathspecs, proves "protected files equal the base commit" at ~176 ms on the largest project (V9).**
  - The command is `git status --ignored=matching` with the protected pathspecs. On a 28 699-file clone of the ailang repo, it caught every arm:
    - a modified test (` M`);
    - a deleted test (` D`);
    - a new test (`??`);
    - a new test hidden by an edit to `.gitignore` (`!!`);
    - a new test in an `info/exclude`d directory (`!! zzign/`).
  - **[r4]** That command misses non-ASCII spellings (V16) and flags excluded directories with no match (V20). §4.5 now uses widened candidate pathspecs, `status --ignored=no`, `ls-files -o -i` and `ls-files -v`, filtered by the mirror (V17–V22).
  - It listed nothing for a clean tree, or for a change to `go.mod`.
  - Row 184's existing full status (`--ignored=no`) takes 360 ms. It **missed** both ignored creations.
  - By contrast, walking and hashing the worktree costs 0.9–4.8 s on the same repo (V10). That is far too slow per call.

## 2. Goals and non-goals

**Goals.**
- **A profile may declare `protected` patterns** (§4.1): an optional list in the exec profile, in a grammar that means the same thing to all four matchers (AILANG's `fileguard`, Seatbelt, bubblewrap and git pathspecs) or is refused at load.
- **Refused to World's write tools for that episode.** The episode policy's `fs_deny_write` gains the patterns, re-based to the AILANG sandbox root. The same applies to `ailang-write`, `ailang-edit`, `ailang-cli` writes and `ailang-run`'s FS effect. World also refuses a write whose path resolves through a symlink, or to a hard link, into a protected path (§4.3, on D-156-4).
- **Denied to the sandbox.**
  - **macOS:** prevented, creation included: directory patterns as literal `denyWrite` dirs, basename patterns as `<worktree>/**/<glob>`.
  - **Linux:** directory patterns are prevented. Basename patterns are **detected only** (D-156-3), and the result says so.
- **Every exec result of a protected episode carries a `protected` block** (§4.5):
  - the pattern-set digest;
  - the base commit and the digest of its protected entries;
  - whether the protected paths equalled the base commit **before** and **after** the command;
  - the first changed paths when they did not.
- **A reader can tell a protected run apart without trusting the agent:**
  - `ailang-worldd why` prints a `protected:` line;
  - a result with no `protected` key ran with no protected set;
  - a result with `intact_after: true` ran against test bytes identical to the base commit.

**Non-goals (with owners).**
- **Test semantics, not only test bytes.** Project code still runs inside the tests. It can `os.Exit(0)` in an `init()`, or read an environment variable and skip. Protecting config files (`vitest.config.ts`, `pyproject.toml`, `package.json`) is the operator's choice in the pattern list. This row protects bytes and records their integrity; it does not prove the run executed every test (R-156-5).
- **A default protected set:** D-156-2 (default: none; the operator writes the list).
- **Failing a call or an episode because protection was breached:** D-156-1 (default: report, never judge).
- **Dependency digests** are row 157. **Resource limits** are row 158.
- **Fixing AILANG's lexical `fileguard` check** (P6): filed upstream (§10). World backstops it with the per-call integrity check.

## 3. Verification Log (2026-10-17; outputs in `design_docs/verification/world-row156-design-2026-10-17/`)

| V | Command | Observed (short) |
|---|---|---|
| V1 | `q0_anchors.out.txt`: `grep -n` over `host/broker/exec_settings.go`, `exec_profile.go`, `handlers_exec.go`, `handlers_ailang*.go` at `b7bae72` | `renderExecSettings` (`:91`) appends `execDenyNames()` under the worktree (`:99`) and srt's defaults (`:101`); `AllowWrite: [Worktree, Cache]` (`:113`). `ExecProfileVersion = "world/exec-profile/v1"` (`:28`); `execWorldEnv` includes `PYTHONPYCACHEPREFIX` and `npm_config_cache` (`:53`); the profile key allowlist is `strictObject("profile", …)` (`:128`), with no `protected` key. `episodeDenyDirs = [.github .pi .claude .ailang]` (`handlers_ailang.go:57`); `renderPolicy` writes `fs_deny_write` (`handlers_ailang_run.go:196`); `RenderEpisodePolicy(sandbox)` per episode (`workspace.go:522`); the exec handler is built per episode (`workspace.go:382`) |
| V2 | `q1_layout.sh` (`git ls-files` per glob on the three checkouts; control = all tracked files) | ailang (`37d6641b3`, 28 699 files): `*_test.go` 1 924, `testdata/` 337, `*.golden` 208, root `tests/` 103, `examples/tests/` 65, `*.test.ts` 30. TwilightGame (`57655894`, 3 141): root `tests/` 246, `*.test.ts(x)` 244, no `__tests__`/snapshots. sunholo-platform (`fe80714f`, 1 682): `test_*.py` 334, `backend/tests/` 317, `cli/tests/` 39, `__tests__/` 191, `*.test.ts(x)` 184, `conftest.py` 4. Profiles' roots (`world-row140-m4/*.json`): ailang `.`, twilightgame `.`, sunholo-cli `cli` |
| V3 | `q2b_dynamic.sh`: one targeted command per project in a throwaway APFS clone (`cp -c` / `git clone --shared`), World's env (`CI=1`, `HOME`/`TMPDIR` in a scratch cache, `PYTHONPYCACHEPREFIX`, `npm_config_cache`), `find -newer marker` | vitest `tests/<first>.test.ts`: `3 passed`, **0** paths written outside `node_modules` (1 inside). pytest `-p no:cacheprovider tests/test_cli_bq.py`: `5 passed, 6 skipped`, **0** paths. Default cacheprovider: `.pytest_cache/v/cache/nodeids` (at `cli/`, not under `tests/`). `go test -count=1 ./internal/lexer/`: `ok`, **0** paths |
| V4 | `q2a_static.out.txt`: `git grep` for golden-update flags and snapshot asserts (control `expect(`) | ailang: 4 `flag.Bool("update-…")` sites (opt-in; e.g. `-update-tailcall-goldens`). Snapshot asserts: TwilightGame 0 files (control 245), sunholo frontend 0 (control 182). sunholo `cli/pyproject.toml`: `testpaths = ["tests"]`, `addopts = "-v --tb=short"` |
| V5 | `q3_srt_denywrite.py`: srt 0.0.78 on macOS, 16 arms × 3 renderings (NONE / GLOB `<wt>/tests`, `<wt>/**/*_test.go`, `<wt>/**/conftest.py`, `<wt>/**/*.test.ts` / LITERAL = existing matches), each arm in a fresh tree, host bytes read unsandboxed | **NONE** (control): all 16 arms written (15 protected + the control). **GLOB**: control `pkg/a.go` written; **every** protected arm refused with `Operation not permitted`, host unchanged. That includes new `pkg/b_test.go`, root `c_test.go` (`**/` matches zero dirs), `cli/sub/conftest.py`, `src/bar.test.ts`, `rm`, `mv` onto a test name, `mv` away, `mv tests tests2`, a symlink `lnk→tests` write, a hard link, `PKG/A_TEST.GO` and a new `pkg/B_TEST.GO` (case-insensitive). One gap: `mv pkg pkg2` **succeeded** (the test file moved with its dir, bytes unchanged; the write after it was refused). **LITERAL**: existing files and the `tests` dir refused (modify, rm, mv away, dir renames, symlink, hard link, case); **all six creation arms written** (`b_test.go`, `c_test.go`, `cli/sub/conftest.py`, `src/bar.test.ts`, `mv a.go z_test.go`, `B_TEST.GO`) |
| V6 | `q0_anchors.out.txt`: `git -C ailang show v0.52.1:` `internal/fileguard/protect.go`, `internal/policytool/fs_ops.go`, `internal/effects/fs_root.go`, `fs_dir.go` | `MatchDenyWrite` (`protect.go:79`): `/**` suffix → folded prefix (`:87`); else `path.Match` on the whole key (`:93`) and the base name (`:96`). policy-tool `protected()` takes `h.root.Rel(p)`, the lexical path (`fs_ops.go:63-70`), then `h.root.OpenFile` (`:93`). The FS effect's `fsCheckMutation` uses `fsRelToRoot`, also lexical (`fs_root.go:145,164`). Rename checks both operands (`fs_dir.go:135,138`) |
| V7 | `q4_policytool_denywrite.py`: the pinned `v0.52.1 policy-tool --policy P` with `fs_deny_write` = `[.github/**]` (control) and then `+ [tests/**, *_test.go, *.test.ts, conftest.py, **/*.test.ts]`; then `[**/*.test.ts, tests/**]` alone | Control: 13 of 14 writes `ok:true`; the 14th (`cli/sub/conftest.py`, parent dir absent) failed `openat`. **With the patterns:** write/edit existing test, new `pkg/b_test.go`, new `tests/new.py`, `cli/sub/conftest.py`, `src/deep/x.test.ts`, `pkg/../tests/new2.py`, an absolute path, `TESTS/new4.py` and `pkg/A_TEST.GO` were all **refused** (`matches fs_deny_write "…"`). **Written through a symlink alias:** `lnk/new5.py` (`lnk→tests`) → `tests/new5.py` NEW; `lnk/test_x.py` → CHANGED; `alias.txt` (`→ pkg/a_test.go`) → CHANGED. **`**/*.test.ts` alone:** `a.test.ts` (depth 0) written, `src/foo.test.ts` (depth 1) refused, `src/deep/x.test.ts` (depth 2) written; `tests/**` does not cover `cli/tests/x.py` (written) |
| V8 | `q0_anchors.out.txt`: installed srt 0.0.78 `README.md`; `command -v bwrap` | **Source, not a measurement:** Linux `denyWrite` takes literal paths; a trailing `/**` is dropped; any other glob is skipped (`:438`). An absent `denyWrite` path gets a temporary mount stub (`:745`). bubblewrap parses at most 9 000 args (≈ 3 000 mounts) (`:753`). Ancestors of a protected path are pinned: `mv`/`rmdir` → `EBUSY` (`:784`). macOS globs "cover existing and new files alike" (`:782`). `bwrap: absent`, Darwin arm64: **Linux unmeasured** |
| V9 | `q6_git_protected.sh`: `git clone --shared` of ailang (28 699 tracked) in scratch; `GIT_OPTIONAL_LOCKS=0 git -c core.fsmonitor=false status --porcelain=v1 -z --untracked-files=all` with `--ignored=no` (row 184) vs `--ignored=matching -- ':(glob,icase)**/*_test.go' ':(glob,icase)**/*.golden' ':(glob,icase)tests/**' ':(glob,icase)examples/tests/**' ':(glob,icase)**/.gitignore'`; six mutation arms + two controls | Cost (3 runs): full status **359–362 ms**; protected pathspecs **174–176 ms**. Clean tree: both empty. Modify / delete / create untracked test: both list it. **`.gitignore += zz*` then new `zz_ign_test.go`**: full lists only ` M .gitignore`; protected lists `!! cmd/ailang/zz_ign_test.go`. **Test in an `info/exclude`d dir**: full lists **nothing**; protected lists `!! zzign/`. Control `go.mod` change: full ` M go.mod`, protected **empty** |
| V10 | `q5_walk_main.go.txt` built in scratch: read-only `WalkDir` of the real checkouts (skipping `.git`, optionally `node_modules`), fold-match the patterns, sha256 every match; controls: no pattern → 0 matches, `*` → every file | ailang (`*_test.go`, `*.golden`, `tests/**`, `examples/tests/**`): 221 201 entries visited, 12 651 matched; **4 817 ms cold**, 896–945 ms warm. TwilightGame: 18.5 ms warm (737 ms with `node_modules`). sunholo-platform: 206–214 ms warm, but `test_*.py` matched **2 924** files: 2 002 `test_*.py` sit inside `.venv` site-packages and 334 outside, equal to the 334 tracked (`q5b_venv_matches.out.txt`). Controls: `[]` → 0 matches; `*` → 207 279 matches |
| V11 | `q7_lstree_cost.out.txt`: `git ls-tree -r -z --full-tree HEAD` on ailang (read-only); `git rev-parse HEAD`; `q2a`: `git ls-files -s` with the protected pathspecs | ls-tree: 28 699 entries in **132–134 ms**; rev-parse 24 ms; `ls-files -s` with pathspecs: 2 283 entries in 37 ms |
| V12 | `q0_anchors.out.txt`: `worktree_delta.go`, `why.go`, `exec_probe.go`, `pins/manifest_test.go` | Row 184 snapshots with `--ignored=no` (`worktree_delta.go:357`); the after-snapshot runs under `WithoutCancel` with a 60 s bound (`:135,142`). `why` decodes a result field lazily, and an old record yields nil (`why.go:585`). Probe arms 1–7 (`exec_probe.go:44-50`). Tool pin `v0.52.1` (`manifest_test.go:31`) |
| V13 | **[r2]** `q8_git_index_flags.sh`: scratch sibling repo + `git worktree add --detach` episode; World's srt rendering (`allowWrite=[worktree, cache]`, `denyWrite` incl. `<wt>/.git`, `<wt>/.git/hooks`, `<wt>/.git/config`; `denyRead` left EMPTY, i.e. more permissive than World's fence); `git update-index --skip-worktree`, `--assume-unchanged`, `git commit -am x`, unsandboxed (control) and under srt 0.0.78. Then, with each flag set unsandboxed and the test file modified: `status --ignored=matching`, `diff-index --cached HEAD`, `diff-files`, `ls-files -v`, all on the protected pathspec. Cost on a scratch clone of ailang | The episode's `.git` is a **pointer file** (`gitdir: …/src/.git/worktrees/ep1`). **Control:** `S` set, `h` set, HEAD **MOVED**. **Under srt:** all three fail with `fatal: Unable to create '…/src/.git/worktrees/ep1/index.lock': Operation not permitted`; the host flag stays `H` and HEAD is unchanged. The index and HEAD live outside `allowWrite`. **Detection with a flag set:** `status`, `diff-index --cached` and `diff-files` all print **nothing** for the modified file; only `ls-files -v` shows it (`h pkg/a_test.go`, `S pkg/a_test.go`). Cost (3 warm runs, V9's four protected pathspecs): `status --porcelain=v2 --branch --ignored=matching` **96–104 ms**; `ls-files -v` **36–42 ms**; `diff-index --cached` + `diff-files` 28–34 + 48–62 ms |
| V14 | **[M0b, iteration 278]** `m0_matrix.py --protected-only` on GitHub Actions `ubuntu-latest` (green CI run 39000212431, job 117301358374, head `4c12f2c`; the first run 38999371033 / job 117298837511 on `4c3b878` gave identical rows but left its fixture trees in the workspace; Linux 6.17 azure x86_64, bubblewrap 0.9.0, node v22.23.3, srt 0.0.78 at the pinned cli.js sha256 `3c3092bd…96e`). V5's 16 arms × NONE / DIRS_ONLY (`denyWrite=[<wt>/tests]`) / LITERAL (`<wt>/tests` + every existing fixture file whose base name case-folds to `*_test.go`, `conftest.py` or `*.test.ts`) = 48 rows, a fresh tree per row, host bytes hashed unsandboxed before and after, outcome taken from the host. Output: `m0b_linux_ci.json` | **NONE:** all 16 arms written (the matrix is not vacuous). **DIRS_ONLY:** refused only inside the directory: create `tests/new.py` (`EROFS`), `mv tests tests2` (`EBUSY`), and a write through a new symlink `lnk→tests` (the link is created, the write gets `EROFS`, `tests/x.py` stays absent). Every base-name arm is **written**: modify, create (nested, root, conftest, `.test.ts`), `rm`, `mv` onto and away, `mv pkg pkg2`, hard link, `PKG/A_TEST.GO` and `pkg/B_TEST.GO`. **LITERAL:** existing literals and the directory refused (modify `EROFS`, `rm` / `mv` away / `mv pkg pkg2` / `mv tests tests2` `EBUSY`, hard link `EXDEV`, the separate upper-case file `PKG/A_TEST.GO` `EROFS`; the symlink arm is PARTIAL as in DIRS_ONLY: link created, write `EROFS`), and the **same six creation arms written** as V5 (`pkg/b_test.go`, `c_test.go`, `cli/sub/conftest.py`, `src/bar.test.ts`, `mv a.go z_test.go`, `pkg/B_TEST.GO`). Gating rows (NONE all written, dir-modify refused in both, LITERAL modify refused), the 48-row validator and `fixtures-removed` all pass |
| V15 | **[M0b]** the same harness on the macOS rig (srt 0.0.78, same cli.js sha256, node v26.0.0), scratch root inside the sprint worktree. Output: `m0b_darwin.json`. Control for the gate: with the `<wt>/tests` deny removed from DIRS_ONLY, the run exits 1 with exactly `DIRS_ONLY/dir-modify` failing | Outcome-for-outcome **identical to V14** in DIRS_ONLY and LITERAL; only the errno differs (`EPERM` everywhere on macOS). So on macOS the dir-only and literal renderings behave as on Linux, and only V5's GLOB rendering adds base-name protection. (On case-insensitive APFS `PKG/` and `pkg/` are one directory, so the case-existing arm hits the same file there; on Linux it is a separate file.) |
| V16 | **[M1c, iteration 281]** `python3 design_docs/verification/world-iter281/repro.py`, Git 2.54.0 on darwin/arm64 at `a0c73053f10092fd511ad22fcf6b95d6a3f896aa`; four separate disposable repositories, NUL-safe output and same-repo unfiltered controls; existing `TestProtectedMatchFixture/long-s` and `/kelvin` rows PASS | **Equivalence FALSE:** ASCII controls are selected; `tests/**` misses `teſts/x.py`, and `*k_test.go` misses `pkg/K_test.go`, in both selected `status` and `ls-files`. Both patterns are already ASCII, so restricting pattern bytes as R-156-7 suggested cannot solve these misses. M1c remains incomplete; preserve Unicode semantics and route a bounded detection design revision, quorum and replan. Full commands/output and scope: `../verification/world-iter281/README.md`, `result.json`; no production change. |
| V17 | **[r4]** `world-row156-r4/v17_fold_scan.go.txt` (go1.26.6, `golang.org/x/text v0.42.0` from `go.mod`; `foldProtectedKey`'s fold copied verbatim): every Unicode scalar value from U+0080 to U+10FFFF, surrogates skipped; print each one whose key contains an ASCII rune. Output `v17_fold_scan.out.txt` | **1 111 936 scanned, exactly 4 hits:** U+017F `ſ` → `S`, U+037E (Greek question mark, NFC singleton) → `;`, U+1FEF (Greek varia, NFC singleton) → `` ` ``, U+212A (Kelvin, NFC singleton) → `K`. U+212B (Ångström) gives `Å`, not ASCII. The fold maps each ASCII letter to its upper case and leaves other ASCII alone. NFC composition only makes non-ASCII runes, so these 4 are the only non-ASCII spellings an ASCII pattern literal can match. **[r4.1]** The fold is NFC (not NFKC), then per-rune `unicode.SimpleFold` orbit minimum: it never expands one rune into several. Controller measurement (iteration 282, `go run` with the repo's x/text and the exact `foldProtectedRune`/NFC code): key(`ß`) = `ß`, key(`ﬀ`) = `ﬀ`, key(`ﬁ`) = `ﬁ`, key(`buﬀ`) = `BUﬀ`, key(`tésts`) = `TÉSTS`, key(`ḋata`) = `ḊATA`. So multi-rune folds (`ß`→`SS`, `ﬀ`→`FF`), compatibility ligatures and composed letters do not fold onto ASCII |
| V18 | **[r4]** `world-row156-r4/v18_wildmatch.sh`, Git 2.54.0, one disposable repo per spelling in the worktree scratch dir, fake HOME, no system/global config. Repo `pkg/é_test.go` (NFC, 2 bytes), modified, `ls-files` and `status` per pathspec. Output `v18_wildmatch.out.txt` | `**/*_test.go` **hits**; `**/?_test.go` **misses** and `**/??_test.go` hits, so `?` matches one **byte**; `**/[^a]_test.go` and `**/[é]_test.go` **miss**; literal `**/é_test.go` hits. So `*` spans multi-byte runes, while `?` and bracket classes do not match like the mirror's per-rune `path.Match` |
| V19 | **[r4]** same script: widened candidate pathspecs on V16's counterexamples plus U+037E, and r3's replacement/`top` obligations | `teſts/x.py`: `tests`, `tests/**`, bare `te*t*` miss; **`te*t*/**` hits**. `pkg/K_test.go` (U+212A): `**/*k_test.go` misses; **`**/*_te*t.go` hits**. `pkg/a;_test.go` (U+037E): `**/a;_test.go` misses; **`**/a*_te*t.go` hits**. `tests/` replaced by a symlink, run from `sub/`: `:(top,glob,icase)te*t*` gives the status entry (1); without `top` it gives 0 |
| V20 | **[r4]** same script, repo with `zzign/a_test.go`, `zzign/b.txt` and `nm/c.txt`, both dirs in `info/exclude`; pathspec `**/*_te*t.go` | `status --ignored=matching -- <pathspec>` prints **`! nm/` and `! zzign/`**: an excluded directory is shown whole, even `nm/`, which holds no matching file, and the entry names no file. `ls-files -z -o -i --exclude-standard -- <pathspec>` prints only **`zzign/a_test.go`**. With no pathspec it prints all 3 files. So r3's ignored check would flag any excluded directory (e.g. `node_modules/`) as a protected change, and a filter cannot judge a bare directory entry |
| V21 | **[r4]** `world-row156-r4/v20_cost.py`: `git clone --shared` of ailang (28 713 tracked) in scratch; 4 runs each, first discarded; then again with a 30 000-file `node_modules/` (2 `*_test.go` inside) in `info/exclude`. Output `v20_cost.out.txt`, clone removed | Clean / with 30k ignored, warm ms: r3 `status --ignored=matching` 109–112 / 112–118 (with 30k ignored it adds the `! node_modules/` entry, V20). **r4 candidate:** `status --ignored=no` 100–104 / 101–113; `ls-files -o -i` 44–47 / 63–64 (0 / **2** entries); `ls-files -v` 18–19 (2 292 entries). **No pathspec:** `status` 193–195 / 194–208; `ls-files -o -i` 43–45 / 65–66 with **30 002 entries (736 054 bytes)**; `ls-files -v` 12–13 with 28 713 entries (2.3 MB) |
| V22 | **[r4]** `world-row156-r4/v22_nfd.sh`: repo with `core.precomposeunicode=false`, NFD `pkg/e`+U+0301`_test.go` in the index via `update-index --cacheinfo` (as a Linux-made commit stores it); `ls-files` under both precompose settings | Index bytes are NFD (`e 314 201`). The NFC literal `**/é_test.go` **misses under both** settings; widened `**/*_te*t.go` **hits under both**. (V18's own "NFD" repo is void as an NFD test: with `precomposeunicode=true`, `update-index` precomposed its argument and stored NFC) |
| V23 | **[r4.1]** `world-row156-r4/v23_ambient_ignored.py` on the **real** checkouts, read-only (`ls-files -o -i` only, `GIT_OPTIONAL_LOCKS=0`, fake HOME): sunholo-platform `fe80714f`, TwilightGame `57655894`, ailang `38525b06b`. Candidate pathspecs per §4.2; the filter is an ASCII stand-in for `matchProtected` (all patterns ASCII). 3 runs each (ailang re-timed with 5); snapshot, then re-run and diff; `lstat` and sha256 of the snapshot timed. Output `v23_ambient_ignored.out.txt` | **Before the fix (standing entries a pristine checkout would report):** sunholo-platform **3 259** (`.venv` 2 032, `node_modules` 877, other 350); TwilightGame **165** (all `node_modules`); ailang **1 240** (ignored outputs under `tests/` and `examples/tests/`). **After (snapshot set-difference):** **0 / 0 / 0**. Cost of the ignored leg, warm minimum: **346 ms** / **177 ms** / **406–504 ms** (cold up to 1 680 / 823 / 2 114 ms). Content identity of the snapshot: `lstat` 155 / 8 / 13 ms; sha256 534 / 22 / 155 ms (49.8 MB / 1.0 MB / 8.7 MB). Not measured: the r4 `status` and `ls-files -v` legs on these checkouts, and the agent-created arms on them (no writes into real repos) |
| V25 | **[iteration 289, controller]** `design_docs/verification/world-iter289/v7r_policytool_aliases.py` and `v7r_doublestar.py` on the **v0.54.15 tool pin** (digest = `host/pins/ailang-releases.json` `tool`) and on v0.52.1 as the control arm; V7's fixture and patterns plus hard-link arms V7 never ran; each arm a fresh tree, verdict = host bytes read back | **v0.52.1 reproduces V7:** `lnk/new5.py`, `lnk/test_x.py`, `alias.txt` (write and edit) and both hard-link aliases are **written** under the protected set. **v0.54.15 refuses all of them, host bytes unchanged:** symlink arms with `(resolves to tests/…)` / `(resolves to pkg/a_test.go)`: matches fs_deny_write; hard-link arms with "the file has more than one hard link". **Controls:** without the protected set every v0.54.15 symlink arm is written (so the refusal is the resolved-name match), and the hard-link refusal fires even for an unprotected file and with an EMPTY `fs_deny_write` (v0.52.1 writes it): it is a property of the pin, not of the protected set. **`**` is doublestar on v0.54.15:** `**/*.test.ts` refuses depth 0, 1 and 2 (V7 on v0.52.1: depth 1 only); `tests/**` still does not match `cli/tests/x.py` |

## 4. Design

### 4.1 The profile key `protected` and its grammar

```json
"protected": ["tests/**", "*.test.ts", "*.test.tsx", "vitest.config.ts"]
```

- Optional. Absent or `[]` means no protected set. Behaviour and result bytes are then identical to today (AC6.4).
- At most 64 entries. Each is 1–256 bytes, NFC, with no NUL, control character or `\`.
- **Every pattern is one of exactly two forms.** Both mean the same in all four matchers (§4.2):
  1. **Directory form `D/**`.** `D` is a clean relative path, worktree-relative: no leading `/`, no `.`/`..` segment, no glob metacharacter. It protects `D` and everything beneath it.
  2. **Base-name form `G`.** It has no `/`. `G` is a `path.Match` pattern using only `*`, `?` and `[…]`. It must compile (`path.Match(G, "")` is not `ErrBadPattern`) and must not contain `**`. It protects every file whose base name matches, at any depth.
- **Refused at load, naming the pattern and the fix:**
  - `**/x` or any other `**` not at the end: "use the base-name form `x`; World's write tools read `**/` as exactly one directory" (V7);
    - **[V25]** That reason is true of v0.52.1 only: the v0.54.15 tool pin reads `**` as doublestar. The refusal stays (the other renderings in §4.2 were designed around the two forms), but its message must stop citing the one-level reading — carried to M7b's guidance;
  - a pattern with `/` that is not `D/**`, e.g. `src/*.ts`: no form means the same in `fileguard` and Seatbelt;
  - a pattern that covers everything: `*`, `**`, or a `D/**` where `D` is `.`;
  - `{`, `}`, `\`, a leading `!`;
  - a `D` that names, or lies inside, `.git` or one of World's deny names (already read-only).
- **Folding.** Patterns and paths are compared case-folded and NFC-normalized, exactly like `fileguard.foldKey` (V6). APFS folds (V5 case arms). On a case-sensitive Linux volume, folding over-denies an exotic spelling, which is `fileguard`'s own fail-closed rule.
- **Canonical form and digest.**
  - Each pattern is NFC-normalized. Duplicates are dropped by fold key, keeping the byte-smallest NFC spelling (so the digest does not depend on input order: `TESTS/**` and `tests/**` give one entry, `TESTS/**`). The list is sorted by fold key bytes. Entries with leading or trailing whitespace are refused (fileguard trims patterns, so a padded entry would mean different things to the two matchers; M1a judge N1). A `D` with a `.git` segment at any depth is refused, as fileguard's `UnderGitDir` already denies it (M1a judge N5).
  - `protected.digest = "sha256:" + hex(sha256("world/exec-protected/v1\n" + join(canonical, "\n") + "\n"))`.
  - The profile's own `digest` (row 140) changes too, because it digests the file bytes. This digest is separate, so two profiles with the same protected set can be compared.
- **`ExecProfileVersion` stays `world/exec-profile/v1`.** The key is optional and additive. An older World refuses a profile that carries it (unknown key, fail closed). That is the safe direction.

### 4.2 One grammar, four renderings (pure functions, `host/broker/exec_protected.go`)

| Consumer | Directory form `D/**` | Base-name form `G` | Evidence |
|---|---|---|---|
| AILANG `fs_deny_write` (World's write tools, `ailang-run` FS) | `D'/**`, where `D'` is `D` re-based to the AILANG sandbox root (below) | `G` verbatim | V6, V7 |
| srt `denyWrite`, darwin | `<worktree>/D` (literal dir) | `<worktree>/**/G` | V5 GLOB |
| srt `denyWrite`, linux | `<worktree>/D` (literal dir; srt binds it read-only) | **not rendered** (D-156-3 = A): srt would skip it silently (V8) | V8 (README) |
| git **candidate** pathspec (the integrity check, §4.5) **[r4]** | `:(top,glob,icase)W(D)` and `:(top,glob,icase)W(D)/**` (r3's bare-D + `/**` pair, kept) | `:(top,glob,icase)**/W(G)` | V9, V17–V22 |

- **Re-basing for AILANG.**
  - The episode's AILANG sandbox is `<worktree>/<module root>` (row 141, `workspace_layout.go:85`), and `fs_deny_write` is relative to it.
  - With module root `.`, nothing changes.
  - With module root `M`:
    - `D/**` with `D` under `M` becomes `D-without-M/**`;
    - `D/**` with `D` equal to `M`, or an ancestor of `M`, protects the whole sandbox, so it renders as `*` (fileguard matches `*` against every base name);
    - `D/**` outside `M` is dropped, because AILANG cannot write outside its sandbox anyway;
    - `G` is unchanged.
- **[r4] Git selects candidates; the mirror decides.** Git's `icase` folds ASCII only, and its `?` and `[…]` match bytes (V16, V18). So the git rendering is no longer an equivalent of the mirror. It is a pathspec that matches **at least** every path the mirror protects, and §4.5 keeps only entries for which `matchProtected` is true. `W` widens a pattern rune by rune: each of `s S k K ; `` ` `` (the ASCII keys of V17's 4 non-ASCII runes), every non-ASCII rune, every `?` and every whole `[…]` class becomes `*`; runs of `*` collapse to one; `/` and all other ASCII stay. Examples: `tests/**` → `te*t*`, `te*t*/**`; `*k_test.go` → `**/*_te*t.go`; `é_test.go` → `**/*_te*t.go`.
- **Why `W` is a superset.** `foldProtectedKey` is NFC (not NFKC), then the per-rune `unicode.SimpleFold` orbit minimum, so no rune folds into several runes and no ligature folds onto ASCII (`ß`, `ﬀ`, `ﬁ` keep non-ASCII keys, V17 [r4.1]). A literal ASCII rune kept by `W` can only be matched in the mirror by itself or its ASCII case pair, because V17 shows no other rune has that key, and `icase` covers the case pair. Every other position is `*`, which in glob mode matches any byte run without `/`, multi-byte and NFD spellings included (V18, V22). `D` has no glob, so `W(D)` keeps its `/` and the r3 bare-`D` + `D/**` pair still covers `D` itself and its contents (V19). The proof rests on V17's list. The M1c test recomputes that list, so a Go or x/text upgrade that adds a rune fails the build instead of shrinking the superset.
- **Equivalence is tested, not assumed:**
  - `fileguard` ↔ World's mirror: a fixture table banked from V7's measured rows, plus a pinned-binary arm that skips without `WORLD_TOOL_AILANG_BIN`.
  - World's mirror ↔ git **[r4]** candidates + filter: the three §4.5 invocations on real git in temp repos, then the filter. The result must equal the mirror exactly (M1c). git is on CI.
  - World's mirror ↔ Seatbelt: real srt, skips without it.
  - Each test drives the same path × pattern matrix.

### 4.3 Refused to World's write tools

1. **Policy.** `RenderEpisodePolicy(root)` becomes `RenderEpisodePolicy(root, protected)`. `renderPolicy` appends the re-based patterns after `episodeDenyWrite` in **every** variant: base, Declassify, Net and Env.
   - So `ailang-write`, `ailang-edit`, the `cli` ops' known writes (policy-tool's `protected()`, V6) and `ailang-run`'s `FS.writeFile`/`renameFile` (`fsCheckMutation`, V6) all refuse a protected path. The refusal is AILANG's own, already worded "matches fs_deny_write …".
   - `workspace.go` passes the episode profile's protected set (`profileFor(episodeID)`). An episode with no profile, or no `protected`, renders exactly today's policy bytes (AC4.3).
2. **World's own pre-check (on D-156-4 = A).**
   - Before forwarding a `Workspace.Write` (`write`/`edit`) request to policy-tool, `AilangToolHandler` resolves the request `path` against the sandbox root, using the existing `resolveExistingPrefix` (`handlers_ailang.go:203`).
   - It then matches the worktree-relative **lexical** path and the **resolved** path against the protected set. A hit on either is a typed refusal: `path %q resolves to protected %q (pattern %q)`.
   - A regular file whose `Lstat` shows `Nlink > 1` is also refused while the episode has a protected set. A hard link can alias a protected inode, and no path check can see that.
   - This closes P6's three measured aliases for the write tools.
   - It cannot reach inside an `ailang-run` program, which picks its paths at run time. That route stays lexical upstream: it is filed (§10) and detected by §4.5.
   - **[V25] Upstream now does this for the write tools.** On the v0.54.15 tool pin, policy-tool matches the spelled AND the resolved name and refuses any multiply-linked file (ailang#1873, `db9719bb7`). Every alias this pre-check targets is already refused, with host bytes unchanged. Whether World still builds its own pre-check (defence against a future tool regression) or relies on the pin plus a real-tool regression arm is D-WORLD-86; until it is answered the pre-check is neither built nor dropped.
   - **Why this is not a "local workaround" of a language gap (CLAUDE.md):** World does not patch, fork or wrap AILANG's matcher. It applies **its own** operator policy to **its own** tool before calling AILANG, as `cliWriteFlags` already does for `fmt --write` (V67 of row 134). The upstream ticket is filed either way. D-156-4 asks Mark to confirm.

### 4.4 Denied to the sandbox

- `execSettingsSpec` gains `Protected []string`, the canonical patterns. `renderExecSettings` appends §4.2's srt column for `spec.GOOS` **after** World's deny names and srt's defaults. The order of the existing entries is unchanged (V50 of row 140 pins `.git/hooks` and `.git/config` first).
- **What is enforced, per OS:**

| Attempt | darwin (measured, V5) | linux (measured, V14; dir form only) | Integrity check (§4.5) |
|---|---|---|---|
| modify / truncate an existing protected file | refused | dir form: refused (`EROFS`); base-name form: **allowed** | `M` → `intact_after: false` |
| create a new file matching a pattern | refused, any depth, case variants too | dir form: refused; base-name form: **allowed** | `??` / `!!` → false |
| delete / rename away | refused | dir form: refused (`EBUSY` for `mv` of the dir); base-name: **allowed** | `D` → false |
| rename a parent dir of a base-name match (`mv pkg pkg2`) | **allowed** (V5) | allowed (V14 DIRS_ONLY; only the non-production LITERAL rendering refuses it, `EBUSY`) | `D` (old path) → false |
| write through a symlink or hard link | refused (V5) | dir form: refused (V14: the link itself is created, the write through it gets `EROFS`); a hard link to a base-name match: **allowed** | content change → false |

- **Linux errnos (V14).** bubblewrap refuses with `EROFS` (writes), `EBUSY` (rename or delete of a bind-mounted path) and `EXDEV` (hard link across a mount), never `EPERM`. Probe arm 8 and every refusal check from M3b on must accept those codes, not macOS's `EPERM` alone.
- **Linux mounts.** Only directory-form patterns are rendered, at most 64, so the profile stays far below bubblewrap's ~3 000-mount limit (V8). Rendering all of ailang's 2 283 protected entries as literals would approach that limit.
- **Settings digest.** The rendered bytes change, so `settings_digest` changes. That is correct: the record already names the exact settings that ran.
- **Startup probe arm 8 `arm8-protected-write`.** It runs only when the profile declares `protected`.
  - World plants, unsandboxed, one file in the probe's scratch episode for the first directory-form pattern, and one for the first base-name pattern.
  - The sandboxed client must print `WORLD-PROBE-REFUSED <errno>` when modifying the planted file, and the host bytes must be unchanged.
  - On darwin it also creates a new file matching the base-name pattern: refused, absent on the host.
  - On linux the base-name leg is **expected to succeed**. The arm records that, and passes only if the integrity check then reports `intact: false` for it. So the detect-only path is proven at startup, not assumed.
  - Control: the same writes unsandboxed succeed. A pass-through shim fails the arm (the row-140 AC3.1 pattern).

### 4.5 Proven intact in the record

**The integrity check [r4].** A no-filter preflight config query runs before the THREE candidate data invocations, inside the same provisional 2500 ms phase bound [r4.x]. It refuses configured filters and query errors; no host filter execution is allowed. The three bounded data invocations are followed by one pure filter. All three run with `cmd.Dir = h.worktree` and the root-anchored candidate pathspecs of §4.2 [r3, r4], each through `runBounded` with row 184's minimal git environment (`gitEnv()`, `GIT_OPTIONAL_LOCKS=0`, `core.fsmonitor=false`, no filters, V12):

```
git status --porcelain=v2 -z --branch --untracked-files=all --ignored=no --no-renames \
  --ignore-submodules=all -- <candidate pathspecs>
git ls-files -z -o -i --exclude-standard -- <candidate pathspecs>
git ls-files -v -z -- <candidate pathspecs>
```

- `status` gives `# branch.oid` (HEAD) and the candidate paths that differ from HEAD: modified, deleted, type-changed (a directory replaced by a file or symlink shows as bare `D`, V19) in `1`/`2` lines, and untracked files in `?` lines.
- **[r4.1] The ignored leg is judged against an episode-start snapshot.** At handler construction, World runs the same pathspec'd `git ls-files -z -o -i --exclude-standard -- <candidates>` once, keeps the mirror-true entries, and stores their fold-keyed path set as `base_ignored_digest` (`sha256` over the sorted fold keys, one per line) and `base_ignored_files` in the `protected` block alongside `base_digest`. The ignored leg reports only entries whose fold key is absent from that snapshot. Ambient dependency content (venv site-packages, `node_modules`, build caches) is then exempt: without it a pristine sunholo-platform checkout reports 3 259 standing entries (2 032 in `.venv`, 877 in `node_modules`), TwilightGame 165 and ailang 1 240; with it all three report 0 (V23). V9's agent-created ignored files (a `.gitignore` edit, an `info/exclude`d dir) still trip, because they postdate the snapshot.
- **[r4] Ignored files come from `ls-files -o -i`, one line per file.** r3 used `status --ignored=matching`. That shows an excluded directory as one entry (`! nm/`) whether or not it holds a match, so it would flag `node_modules/` as a protected change and give the filter no file name to judge (V20, V21). `ls-files -o -i` with the pathspec lists exactly the matching ignored files (V20: `zzign/a_test.go` only).
- **[r2] Index flags:** in the `ls-files -v` output, an entry tagged with a lower-case letter (assume-unchanged) or `S` (skip-worktree) is reported with `status: "index-flag"`. With either flag set, `git status`, `diff-index --cached` and `diff-files` all report **nothing** for a modified file (V13), so `ls-files -v` is chosen by measurement.
- **[r4] The filter.** Every path entry from the three outputs is kept only if `matchProtected(canonical, path)` is true. Git's raw bytes are passed as they are (with `-z` git does not quote paths); the mirror does its own NFC fold, so an NFD path in the index is judged like its NFC spelling (V22). `branch.oid` is not filtered. Git can over-select (e.g. `te*t*` also selects `tests2/`); it cannot under-select (§4.2, V17–V19). A bare directory entry cannot occur, because `status` runs with `--ignored=no` and `--untracked-files=all`.
- **[r4] Fail closed on any doubt.** Each invocation's stdout is read through a cap, `protectedCheckMaxOutput` = 8 MiB, 3.5× V21's largest no-pathspec output. A non-zero exit, a timeout, output over the cap, or a record that does not parse (a missing NUL, an unknown `status` line type, invalid `ls-files -v` tag) sets `error`, with intact `null`. It never yields `intact`, never drops the excess silently, and never fails the call.
- **[r2] Why the flags and HEAD should never move (prevention, measured).** Episodes are `git worktree add` checkouts, so the worktree's `.git` is a pointer file. The index and HEAD live in the source repo's `.git/worktrees/<name>/`, outside srt's `allowWrite`. Under World's rendering, `update-index --skip-worktree`, `--assume-unchanged` and `commit` all fail on `index.lock` with `Operation not permitted`, and HEAD stays put (V13). World's write tools cannot reach that directory either: it is outside the AILANG `fs_sandbox`, and any `.git` component is read-only to policy-tool (V6). The flag check is defence in depth against a route not yet found.
- **[r2] intact** requires **both** `branch.oid == base` **and** zero entries left after the filter. When `branch.oid != base`, `changes` gets `{"path": "HEAD", "status": "head-moved", "phase": …}` and intact is unconditionally `false`. The path entries are still listed, against the current HEAD, but they cannot say which protected paths differ from `base`. The `head-moved` entry says exactly that.
- It runs **before** the command (after stack re-verification) and **after** it (before row 184's `finish`). **[r4.1]** Each phase (before and after) runs all three git invocations under a single shared context (`protectedCheckTimeout`), reserving that bound from the child budget for the after-phase. **[r4.2, quorum r2 fix verbatim (gemini-4-argon, oc-kimi-k3)]** "raise `protectedCheckTimeout` from `500 ms` to `2500 ms`": V23 measures the ignored leg alone at 406–504 ms warm and 2 114 ms cold on ailang, so 500 ms would time out on every call on a dependency-installed checkout.
  - The child's budget becomes `min(timeout_ms, until(deadline) − protectedCheckTimeout)`. The after-check therefore always fits the handler's cap. **[r4.2]** Recomputed for the 2 500 ms bound (oc-kimi-k3 fix (1)): with today's 9 000 ms cap and 10 s handler deadline, a command keeps at most `10 s − (time already used, including the before-phase) − 2.5 s`, i.e. ≈ 7.3 s after a warm ~0.2–0.65 s before-phase and ≈ 5 s after a cold 2.5 s before-phase. This only binds a profile that declares `protected`; an unprotected profile keeps its 9 s unchanged. V24 (M6 gate) re-sizes the bound from measurement.
  - A check that times out or fails records `error`, with `intact_before`/`intact_after` `null`. It never fails the call and is never silent (row 184's rule).
- **[r4] Measured cost** on the largest project (ailang, 28 713 files, warm): `status` 100–113 ms + `ls-files -o -i` 44–64 ms + `ls-files -v` 18–19 ms ≈ **165–195 ms** per phase, with or without a 30 000-file ignored tree (V21). The filter is pure Go over at most a few thousand entries (2 292 on ailang, V21). All three share the phase's one context [r4.1], bounded at 2 500 ms [r4.2]. Two phases add ≈ 0.35–0.4 s per call there, and less on the smaller projects. **[r4.2, judge r1 fix]** That figure is V21's clean shared clone only. On the real dependency-installed checkouts V23 measures the ignored leg alone at 406–504 ms warm and 2 114 ms cold on ailang (346 ms / 1 680 ms on sunholo-platform), which is why the phase bound is 2 500 ms and V24 gates M6. The handler-construction step below warms the cache.
- **[r4] Why not drop the pathspec.** Running the three commands with no pathspec and filtering everything is also complete, and also fits (V21: ≈ 250–290 ms). It lost because its output grows with every file in the tree: `ls-files -v` returns 28 713 entries (2.3 MB) on ailang, and `ls-files -o -i` returns every ignored file (30 002 entries for one 30k `node_modules/`), so a larger ignored tree would hit the output cap and turn a clean run into `error`. The candidate pathspec returned 2 292 and 2 entries for the same trees.

**Baseline at handler construction** (`NewExecHandler`, once per episode × project). **[r4.2, quorum r2 fix (oc-glm-5-3 (3))]** The construction step (`rev-parse`, `ls-tree`, the ignored snapshot) runs under its own named bound, `protectedBaselineTimeout` = 3 200 ms (V23's 2 114 ms cold ignored leg × 1.5, GLM's sizing rule; re-sized from V24's max × 1.5 before M6). A timeout there sets `error` for the episode, never intact:
- `base` = `git rev-parse HEAD` of the worktree (24 ms, V11).
- `base_digest` is built from `git ls-tree -r -z --full-tree <base>` (≈ 133 ms on 28 699 entries, V11), filtered by **World's own matcher**: `sha256` over each matching entry's `<mode> <oid>\t<path>\n`, in ls-tree order. `base_files` is the number of entries.
- **[r4.1]** `base_ignored_digest` / `base_ignored_files`: the ignored snapshot above, taken in the same construction step. A snapshot that fails or times out sets `error` for the episode, as a failed `base` does. **[r4]** Unchanged: this step already filters with the mirror, not with a git glob, so V16 does not touch it.
- So the reader gets a content identity of the protected tree **as of the base commit**. Anyone holding the repo can recompute it, and git computes no glob in it.
- Not a git worktree, or git missing: `base` is empty and every result carries `protected.error = "worktree is not a git checkout"`, with intact `null`. Every World episode today is a git worktree (row 184 D2).

**The result field** (added to `execResult`, `omitempty`; absent when the profile declares no protected set):

```json
"protected": {
  "digest": "sha256:…", "patterns": ["*.test.ts", "tests/**"],
  "enforced": "glob",
  "base": "37d6641b308a3e8ffb4d5a2023ba4faa89a5254a", "base_digest": "sha256:…", "base_files": 247,
  "base_ignored_digest": "sha256:…", "base_ignored_files": 0,
  "intact_before": true, "intact_after": false,
  "changes": [{"path": "tests/zz_new.test.ts", "status": "untracked", "phase": "after"}],
  "truncated": false, "total": 1,
  "error": ""
}
```

- `enforced`:
  - `"glob"` (darwin): both forms denied in the sandbox;
  - `"dirs-only"` (linux): base-name forms detected only.
- `changes` holds at most 32 entries, then `truncated`/`total`. `status` is one of `modified`, `deleted`, `untracked`, `ignored`, `index-flag` [r2] or `head-moved` (path `"HEAD"`).
- Paths are cut the way row 184 cuts them (`worktreeDeltaMaxPathJSON`), so the worst case stays far under the 1 MiB output cap. A test measures it (row 184's `TestWorktreeDeltaWorstCaseFitsTheOutputCap` pattern).
- `intact_before: false` means the protected tree was already changed when the call started. This makes the state **sticky** with no extra bookkeeping: once a protected file differs from the base commit, every later call says so until the file is restored.

**Backward compatibility and replay.**
- Results are stored verbatim (`world/effect-result/v1`; row 184 V2), and replay returns them without re-executing (row 140 §4.7). An old record has no `protected` key and decodes and replays as before.
- No record schema, `InterfaceHash` or kernel type changes.
- `execDecode` (a test helper using `DisallowUnknownFields`) learns the field, as it learned `worktree_delta`.

**`why`.**
- `resultProtected(payload)` decodes the field lazily, nil for old records (the `resultWorktreeDelta` pattern, V12). The text form prints `protected: intact (247 files @ 37d6641, digest sha256:…, enforced glob)`, or `protected: CHANGED after the call — tests/zz_new.test.ts (untracked) …`, or `protected: not checked — <error>`. `--json` carries the whole block.
- An exec result with no `protected` key prints `protected: none declared`, so an unprotected run is visible as such.

### 4.6 Where each guard sits (summary)

| Route to a protected path | Prevention | Detection |
|---|---|---|
| `ailang-write` / `ailang-edit`, lexical path | AILANG `fs_deny_write` (§4.3.1) | §4.5 |
| the same through a symlink or hard link | World pre-check (§4.3.2, on D-156-4) | §4.5 |
| `ailang-cli` op that writes | AILANG `protected()` (§4.3.1) | §4.5 |
| `ailang-run` FS effect | AILANG `fsCheckMutation`; lexical only (upstream ticket) | §4.5 |
| `workspace-exec`, darwin | Seatbelt `denyWrite` (§4.4) | §4.5 |
| `workspace-exec`, linux | dir form only (§4.4) | §4.5 |

### 4.7 Why not a package (S3), and why not the kernel

- No `world/` change and no new effect, tool or grant.
- The pattern grammar is operator configuration, parsed at the host boundary. The enforcement lives in three external engines (AILANG `fileguard`, Seatbelt, bubblewrap) and git, and only the host can drive them (S2).
- The `se-tools` package is untouched: the `workspace-exec` plan is unchanged.

### 4.8 Why no `.ail` in this row (S1/S5)

The law this row must honour is **AILANG's own `fileguard` matcher**, which is Go upstream. A World `.ail` re-statement would be a third copy with no Z3-checkable invariant beyond the grammar predicate.

The canonical text form (the digest input) is pinned by a named Go golden test. Each matcher equivalence is bound by a test that drives the real engine where it exists (§4.2), so the truth source is measured, not restated. The evaluator may challenge this. The alternative is a `design_docs/sketches/protected.ail` holding a contracted `protectedPatternOk` and named tests for the digest's text form, about +40 LOC in M1.

## 5. Premises

| P | Premise | Evidence | vs the charter row |
|---|---|---|---|
| P1 | The exec sandbox and World's write tools leave test files writable today | V1, V5 NONE, V7 control | as stated |
| P2 | Go tests sit beside the code (1 924 files), so basename patterns are required; TS/Python tests sit in `tests/` dirs | V2 | (new; the row asked to measure it) |
| P3 | Targeted test runs write nothing under the protected paths; snapshots are unused; Go goldens update only behind opt-in flags | V3, V4 | answers "measure first" |
| P4 | `**/*.test.ts` means one directory level to World's write tools; the grammar must be restricted | V7 | **FALSE as written:** the row's example glob does not protect what it names |
| P5 | Linux srt skips non-`/**` globs in `denyWrite`; render-time expansion protects existing files but cannot deny creation | V5 LITERAL, V8 | **INCOMPLETE:** "expand at render time" leaves creation open (6/6 arms written), and per-file expansion risks bubblewrap's mount limit; detection is required. **Linux confirms (V14):** LITERAL leaves the same 6/6 creation arms written, and DIRS_ONLY protects only the directory form, as §4.4 states |
| P6 | AILANG's deny check is lexical; a symlink alias bypasses it | V6, V7, V25 | **FALSE as implied:** adding patterns to the write tools alone does not refuse every write. **[V25] True on v0.52.1 only:** on the v0.54.15 tool pin, policy-tool write/edit refuse every measured symlink and hard-link alias |
| P7 | A pathspec-limited `git status --ignored=matching` detects every measured content tamper, ignored creations included, in ~176 ms on 28 699 files. **[r2]** It is blind to a file under assume-unchanged/skip-worktree, as are `diff-index`/`diff-files`; `ls-files -v` sees the flag. **[r4] FALSE as written for non-ASCII spellings and excluded directories:** the r3 pathspecs miss `teſts/x.py` and `pkg/K_test.go` (V16) and a byte-wise `?`/`[…]` misses multi-byte names (V18), while `--ignored=matching` flags whole excluded dirs that hold no match (V20). **Replaced by:** widened candidate pathspecs (a measured superset, V17–V19, V22) + `status --ignored=no` + `ls-files -o -i` + `ls-files -v`, filtered by `matchProtected`, ≈ 165–195 ms per phase on 28 713 files (V21) | V9, V13, V16–V22 | (new) |
| P12 | **[r2]** Sandboxed code cannot set index flags or move HEAD: the linked worktree's index and HEAD are outside `allowWrite` | V13 | (new) |
| P8 | A worktree walk-and-hash is too slow per call (0.9–4.8 s) | V10 | (new) |
| P9 | "Linux srt has no globs (V33)" | V8 | **PARTLY STALE:** on 0.0.78 Linux expands `denyRead`/`allowRead` globs; only `denyWrite`/`allowWrite` are literal-only |
| P10 | `PYTHONPYCACHEPREFIX` and `npm_config_cache` are already in World's child env | V1 | confirmed |
| P11 | A new optional result field is replay- and record-compatible | V12, row 184 V2/V4 | confirmed |

## 6. Milestones

Every milestone's tests **pass with `WORLD_TOOL_AILANG_BIN` unset**. A test that needs real srt, the pinned tool binary or bubblewrap **skips cleanly** when it is absent, and has a pin-free companion that drives the same logic. Gates are `go vet ./...`, `go test -race ./...` and `./scripts/verify_ail.sh`, with `AILANG_BIN` at the gate pin.

**M0 — Linux measurement first (CI only; no production code).**
- AC0.1: extend the row-140 M0 matrix step (`world-row140-m0/m0_matrix.py`, run in `ci.yml`'s go job on `ubuntu-latest`) with V5's 16 arms under the dir-only and literal renderings.
- AC0.2: record the observed tokens as V-rows in this doc. If Linux differs from §4.4's table, amend §4.4 and P5 before M3.
- Negative control: the NONE rendering must show every arm written. Otherwise the matrix is vacuous and the step fails.

**M1 — the pure core (`host/broker/exec_protected.go`, ≤ 200 LOC + tests).**
- `ParseProtectedPatterns` (the grammar), `canonicalProtected`, `protectedDigest`, `matchProtected` (fold like fileguard), `rebaseForSandbox`, `renderSrtProtected(goos)` and `gitPathspecs`.
- Tests:
  - `TestProtectedGrammar`: every §4.1 refusal by name;
  - `TestProtectedDigestGolden`: fixed bytes;
  - `TestProtectedMatchFixture`: V7's measured rows, with outcomes copied from `q4_policytool_denywrite.out.txt`;
  - `TestProtectedRebase`;
  - ~~`TestProtectedGitPathspecEquivalence`~~ **[r4]** replaced by M1c below (V16 refuted plain equivalence).
- **Negative controls:**
  - `TestProtectedGrammarRefusesDoubleStarPrefix` dies if `**/` is accepted;
  - the fold row `PKG/A_TEST.GO` dies if folding is dropped.

**M1c [r4] — git candidate renderer + mirror filter (pure code, real-git equivalence test).** Waits for the r4 quorum review.
- **File allowlist:** `host/broker/exec_protected.go`, `host/broker/exec_protected_test.go`, new `host/broker/exec_protected_git_test.go`. Nothing else. **Cap:** ≤ 140 LOC production, ≤ 350 LOC tests.
- **[r4.1] Signature:** `parseProtectedCandidates` also takes the episode-start ignored snapshot (a fold-key set) and drops ignored entries found in it. **Code:** `gitCandidatePathspecs(canonical)` replaces `gitPathspecs` (§4.2's `W`; M1b's `TestProtectedRenderers/git` golden is re-derived to the widened form, since the old rendering is not a superset, V16). `protectedWidenRunes` = `s S k K ; `` ` ``. `parseProtectedCandidates(statusOut, ignoredOut, flagsOut []byte) (head string, changes []protectedChange, err error)` parses the three NUL-separated outputs, applies `matchProtected`, enforces `protectedCheckMaxOutput`, and returns an error on any malformed record. No I/O in production code: the invocations are wired in M6.
- AC1c.1 `TestProtectedFoldToASCIIClosed`: scans every scalar value U+0080–U+10FFFF (surrogates skipped) with `foldProtectedKey`, which is NFC (not NFKC) then per-rune `unicode.SimpleFold` orbit minimum [r4.1]; it also asserts key(`ß`) = `ß`, key(`ﬀ`) = `ﬀ` and key(`ﬁ`) = `ﬁ` (no multi-rune fold). The non-ASCII runes whose key holds an ASCII rune are exactly U+017F, U+037E, U+1FEF, U+212A (V17), and each such ASCII key, in both cases, is in `protectedWidenRunes`.
- AC1c.2 `TestProtectedCandidatePathspecsGolden`: `tests/**` → `:(top,glob,icase)te*t*`, `:(top,glob,icase)te*t*/**`; `*k_test.go` → `:(top,glob,icase)**/*_te*t.go`; `é_test.go` and `?_test.go` → `:(top,glob,icase)**/*_te*t.go`; `[ab]_x.go` → `:(top,glob,icase)**/*_x.go`; runs of `*` collapse.
- AC1c.3 `TestProtectedGitCandidateEquivalence` (real git; fails if git is absent on CI, skips only locally): one disposable repo per spelling under `t.TempDir()`, fake `HOME`, `GIT_CONFIG_NOSYSTEM=1`, `GIT_CONFIG_GLOBAL=/dev/null`. The bank is **every** `protectedFixtureRows` row unchanged (including `long-s`, `kelvin`, `pkg-A_TEST.GO` and the three non-hits), plus: `pkg/a;_test.go` with U+037E (pattern `a;_test.go`); NFD `pkg/e`+U+0301`_test.go` stored with `core.precomposeunicode=false` (pattern `é_test.go`); `pkg/é_test.go` (pattern `?_test.go`); an excluded `zzign/a_test.go` next to an excluded `nm/c.txt`; `tests` replaced by a symlink and, separately, by a file; `tests/x.py` deleted; an assume-unchanged and a skip-worktree protected file; each case run from the root and from a subdirectory CWD. Each row makes its change, runs the three §4.5 commands with `gitCandidatePathspecs`, and feeds `parseProtectedCandidates`. The reported path set must **equal** {changed paths for which `matchProtected` is true}: every mirror-true row reported with the right `status`, every mirror-false row absent (`tests2/x.py`, `nm/c.txt`, and the [r4.1] no-multi-rune-fold controls `pkg/a_teßt.go` against `*_test.go`, which `W` selects and the filter drops, and `buﬀ/x.py` against `buff/**`).
- AC1c.5 **[r4.1]** `TestProtectedIgnoredBaseline`: `parseProtectedCandidates` takes the snapshot set. A standing ignored protected file present in the snapshot is not reported; a new one (absent from it) is reported as `ignored`; a case-variant spelling of a snapshot entry is treated as present (fold key). It dies under MUT-NO-IGNORED-BASELINE.
- AC1c.4 `TestProtectedCandidateParseFailsClosed`: output over the cap, a record missing its NUL, an unknown line type and a bad `ls-files -v` tag each return an error, never an empty change list.
- **Negative controls (each must kill a named row):** MUT-ASCII-ONLY-SUPERSET kills `long-s`, `kelvin` and U+037E; MUT-QMARK-KEPT kills `?_test.go`; MUT-NO-MIRROR-FILTER kills `non-hit-sibling` (`te*t*` selects `tests2/`); MUT-IGNORED-MATCHING kills `nm/c.txt`; MUT-NO-TOP kills the subdirectory row; MUT-NO-BARE-D kills the symlink/file replacement rows; MUT-CAP-SILENT kills AC1c.4; [r4.1] MUT-NO-IGNORED-BASELINE (snapshot ignored) kills AC1c.5.

**M2 — the profile key (≤ 80 LOC + tests).**
- `strictObject` gains `protected`. `ExecProfile.Protected` holds the canonical list. Load refusals name `protected[i]`.
- `TestExecStartupRefusalTable` gains the grammar rows.
- Negative control: `TestExecProfileProtectedUnknownFormRefused` (`src/*.ts`) dies if the load check is removed.

**M3 — sandbox rendering (≤ 120 LOC + tests).**
- `execSettingsSpec.Protected`, rendered per GOOS.
- Tests:
  - goldens for darwin and linux renderings;
  - `TestExecSettingsNoProtectedIsByteIdentical`;
  - `TestExecSrtProtectedWrites` (real srt, skips without it): V5's GLOB arms through the handler, host bytes unchanged.
- Negative control: `TestExecSrtProtectedWrites/modify` dies if the append is deleted. The pin-free companion `TestExecSettingsProtectedRendered` dies too.

**M4 — probe arm 8 (≤ 150 LOC + tests).**
- `arm8-protected-write` as in §4.4, including the linux detect leg.
- Tests:
  - `TestExecProbeArm8` (real srt, skips);
  - pin-free: the pass-through shim fails arm 8 (`TestExecProbeArm8FailsUnsandboxed`), which is the negative control;
  - a profile without `protected` runs arms 1–7 only (`TestExecProbeArm8Absent`).

**M5 — World's write tools (≤ 180 LOC + tests).**
- `RenderEpisodePolicy(root, protected)` in all four variants, wired in `workspace.go`. The §4.3.2 pre-check is behind D-156-4 (if D-156-4 = B, M5 is the policy half only).
- Tests:
  - policy goldens with and without a protected set (`TestEpisodePolicyNoProtectedIsByteIdentical`);
  - `TestEpisodePolicyProtectedRebased` (module root `sub`);
  - pin-free `TestWorkspaceWriteProtectedPrecheck`: a symlink `lnk→tests`, `alias.txt→pkg/a_test.go`, a hard link and a lexical hit, each refused before any subprocess, with a spawn counter of 0;
  - pinned `TestWorkspaceWriteProtectedPolicyTool`, which skips without `WORLD_TOOL_AILANG_BIN` and replays V7's arms.
- **Negative controls:**
  - deleting the pre-check kills the symlink row;
  - deleting the policy append kills the pinned test and the golden.

**M6 — the record (≤ 220 LOC + tests).**
- The baseline at `NewExecHandler`, the before/after checks (**[r4]** M1c's three invocations and `parseProtectedCandidates`, unchanged), the `protected` field, the budget reserve and `protectedCheckTimeout`.
- Tests:
  - `TestExecProtectedIntact`: a profile command that touches nothing gives `intact_before` and `intact_after` both true, with `base`, `base_digest` and `base_files` asserted against an independent `git ls-tree`;
  - `TestExecProtectedDetectsCreation`: a command that creates `zz_test.go` gives `intact_after: false`, `untracked`, and the next call `intact_before: false`;
  - `TestExecProtectedDetectsIgnoredCreation` (`.gitignore` arm);
  - `TestExecProtectedNotGit`: `error` set, intact null;
  - **[r2]** `TestExecProtectedIndexFlag`: the test sets `--assume-unchanged` and, separately, `--skip-worktree` on a protected file host-side, then modifies it. The next call gives `intact_before: false` with `status: "index-flag"`. This is the **negative control for the flag guard**: it dies under MUT-NO-INDEX-FLAG-CHECK, because `status` alone reports nothing (V13);
  - **[r2]** `TestExecProtectedHeadMoved`: a host-side commit moves HEAD. The result is intact `false` with `{"path":"HEAD","status":"head-moved"}`, even though the entry list is empty;
  - **[r2]** `TestExecSrtIndexWriteRefused` (real srt; skips without it): `git update-index --skip-worktree` and `git commit` under the rendered settings fail, and the host index flag and HEAD are unchanged (V13 as a regression arm);
  - **[r4.2] M6 gate — V24 first (quorum r2 fix verbatim, oc-glm-5-3 (1)/(4) and oc-kimi-k3 (2)):** "git worktree add of sunholo-platform and ailang with .venv/node_modules installed, measuring all three §4.5 legs plus the construction-time snapshot together, warm and cold, with the real matchProtected filter (not the ASCII stand-in)", "requiring p95 phase time below the budget with intact ≠ null". If V24 does not fit, the bounds are re-sized from V24's max × 1.5 and this doc is amended before M6 starts. V24 also records whether every M4 episode installs its dependencies before `NewExecHandler` (oc-kimi-k3 fix (3), first option); if not, a mid-episode install shows as ignored changes (`intact_after: false`, a false CHANGED, never a false intact) and is a disclosed residual (R-156-6);
  - `TestExecProtectedCheckTimeout`: a fake git sleeping past `protectedCheckTimeout` (2 500 ms [r4.2]) gives `error`, the call still succeeds, and stays inside the handler deadline;
  - `TestExecNoProtectedIsByteIdentical`;
  - **[r4.1]** `TestExecProtectedAmbientIgnoredIntact`: an ignored `.venv/lib/site-packages/pkg/test_a.py` present before the handler is built gives `intact_before`/`intact_after` true; a command that creates a new `test_b.py` in that site-packages dir gives `intact_after: false` with `status: "ignored"`;
  - the worst-case-size test.
- These tests run the command through a shim srt that runs it unsandboxed, so they need neither srt nor the tool binary.
- **Negative control:** `MUT-ALWAYS-INTACT` (the check returns true) kills `TestExecProtectedDetectsCreation`.

**M7 — `why` and docs (≤ 150 LOC + docs).**
- `resultProtected` and the `protected:` line (three states plus `none declared`), with `TestWhyProtectedLine` and an old-record row.
- QUICKSTART §10 gains the `protected` key, the grammar table, the per-OS enforcement table, and worked patterns for the three projects (V2):
  - ailang: `*_test.go`, `*.golden`, `tests/**`;
  - TwilightGame: `tests/**`, `*.test.ts`, `vitest.config.ts`;
  - sunholo-cli: `cli/tests/**`, `conftest.py`, `pyproject.toml`. The last is a base-name pattern, so it protects every `pyproject.toml` in the worktree; the grammar has no single-file form.
- `TestQuickstartSection10ProfilesLoad` loads them.
- Negative control: `TestWhyProtectedLine/none-declared` dies if the absent case prints nothing.

## 7. Load-bearing mutations

| Mutant | Killer |
|---|---|
| MUT-DOUBLESTAR-ACCEPTED: the grammar admits `**/x` | `TestProtectedGrammarRefusesDoubleStarPrefix` |
| MUT-NO-FOLD: `matchProtected` compares bytes | `TestProtectedMatchFixture` row `pkg/A_TEST.GO` |
| MUT-PATHSPEC-NO-ICASE: git pathspecs without `icase` | `TestProtectedGitCandidateEquivalence` row `pkg-A_TEST.GO` [r4] |
| **[r4]** MUT-ASCII-ONLY-SUPERSET: candidates rendered as r3's unwidened pathspecs | `TestProtectedGitCandidateEquivalence` rows `long-s`, `kelvin`, U+037E (V16, V19) |
| **[r4]** MUT-QMARK-KEPT: `?` / `[…]` passed to git unwidened | `TestProtectedGitCandidateEquivalence` row `?_test.go` on `pkg/é_test.go` (V18) |
| **[r4]** MUT-NO-MIRROR-FILTER: every candidate reported | `TestProtectedGitCandidateEquivalence` row `non-hit-sibling` |
| **[r4]** MUT-IGNORED-MATCHING: ignored files from `status --ignored=matching` | `TestProtectedGitCandidateEquivalence` row `nm/c.txt` (V20) |
| **[r4]** MUT-NO-TOP / MUT-NO-BARE-D: `top` or the bare `W(D)` pathspec dropped | subdirectory-CWD row / symlink and file replacement rows (V19) |
| **[r4]** MUT-WIDEN-SET-STALE: `protectedWidenRunes` misses a rune V17's scan finds | `TestProtectedFoldToASCIIClosed` |
| **[r4]** MUT-CAP-SILENT: output over the cap truncated to what fits | `TestProtectedCandidateParseFailsClosed` |
| **[r4.1]** MUT-NO-IGNORED-BASELINE: ignored entries not checked against the episode-start snapshot | `TestProtectedIgnoredBaseline`; `TestExecProtectedAmbientIgnoredIntact` |
| MUT-NO-SRT-PROTECTED: `renderExecSettings` ignores `Protected` | `TestExecSettingsProtectedRendered` + `TestExecSrtProtectedWrites/modify` |
| MUT-LINUX-GLOB: linux renders `<wt>/**/G` (silently skipped by srt) | linux golden; M0 arm shows it does nothing |
| MUT-PROTECTED-BEFORE-GITHOOKS: protected entries inserted before `.git/hooks` | `TestExecSettingsGoldenDarwin/Linux` order (row-140 V50) |
| MUT-NO-POLICY-APPEND: `fs_deny_write` without the patterns | `TestEpisodePolicyProtectedGolden`; pinned `TestWorkspaceWriteProtectedPolicyTool` |
| MUT-VARIANT-MISSED: the Env/Net variant policy omits the patterns | `TestEpisodePolicyProtectedAllVariants` |
| MUT-NO-REBASE: patterns not re-based to module root `sub` | `TestEpisodePolicyProtectedRebased` |
| MUT-PRECHECK-LEXICAL-ONLY: the pre-check skips resolution | `TestWorkspaceWriteProtectedPrecheck/symlink-dir` |
| MUT-NO-NLINK: the hard-link refusal removed | `TestWorkspaceWriteProtectedPrecheck/hardlink` |
| MUT-ALWAYS-INTACT: the check reports intact without running git | `TestExecProtectedDetectsCreation` |
| MUT-IGNORED-NO: **[r4]** the `ls-files -o -i` ignored-file invocation removed (`status` itself runs `--ignored=no` by design) | `TestExecProtectedDetectsIgnoredCreation` |
| MUT-HEAD-UNCHECKED: `branch.oid` ignored (intact from entries alone) | `TestExecProtectedHeadMoved` |
| MUT-NO-INDEX-FLAG-CHECK: the `ls-files -v` check removed | `TestExecProtectedIndexFlag` |
| MUT-NO-RESERVE: the child gets the full remaining budget | `TestExecProtectedCheckTimeout` (the after-check must fit the deadline) |
| MUT-CHECK-FAILS-CALL: a check error fails the call | `TestExecProtectedNotGit` (the call must succeed with `error`) |
| MUT-FIELD-ALWAYS: `protected` emitted with no patterns | `TestExecNoProtectedIsByteIdentical` |
| MUT-ARM8-RC-ONLY: arm 8 accepts any non-zero rc | `TestExecProbeArm8FailsUnsandboxed` (missing-client variant) |

## 8. Decisions for Mark (each with a recommended default and a default-while-unanswered)

- **D-156-1 — A protected path changed: what happens to the call?**
  - **A (recommended): report, never judge.** The call returns the command's result with `intact_after: false`, and `why` shows it. This matches row 140's "rc 0 is reported, never judged" and row 184's "a snapshot never fails the call".
  - B: the call fails (handler error, recorded `failed`). The test output is lost, and an honest agent whose test happened to write a snapshot is punished.
  - C: A, and every later `workspace-exec` call in the episode is refused.
  - *Default while unanswered: A.*
- **D-156-2 — A profile declares no `protected`.**
  - **A (recommended): no protected set.** The result carries no `protected` key and `why` prints `none declared`. No authority changes.
  - B: a built-in default per detected toolchain (`*_test.go`, `tests/**`, …). It would guess, and a wrong guess silently refuses legitimate writes.
  - *Default while unanswered: A.*
- **D-156-3 — Base-name patterns on Linux.**
  - **A (recommended): detected only.** Not rendered for srt; `enforced: "dirs-only"`; arm 8 proves detection at startup.
  - B: render the existing matches as literal paths up to a cap (e.g. 512). This protects modify/delete of existing files, but not creation (V5 LITERAL). It adds mounts (V8), and ailang alone has 2 283 entries.
  - C: refuse base-name patterns on Linux (Go projects then cannot be protected there).
  - *Default while unanswered: A.*
- **D-156-4 — World's own pre-check on its write tools (symlink / hard-link aliases, §4.3.2).**
  - **A (recommended): yes.** World applies its own operator policy to its own tool before calling AILANG (precedent `cliWriteFlags`), and the upstream ticket is filed too.
  - B: no pre-check. Rely on the upstream fix plus §4.5 detection. This respects the strictest reading of "no local workarounds" in CLAUDE.md.
  - *Default while unanswered: B.* M5 lands the policy half, and the pre-check waits for A.
  - **Answered A (attended 2026-10-17, D-WORLD-82).** **[V25, iteration 289]** The tool pin moved to v0.54.15 the same day and closes every measured alias upstream; D-WORLD-86 asks whether M5c/M5d are still wanted.

## 9. Risks and residuals

- **R-156-1 (Linux creation, D-156-3 = A).** On Linux an agent's project code can create a base-name-matching file, e.g. a `zz_test.go` with a `TestMain` that skips `m.Run()`. It is **detected** (`intact_after: false`), not prevented. `enforced: "dirs-only"` states this in every result.
- **R-156-2 (directory rename on darwin, V5).** `mv pkg pkg2` moves a base-name-protected file with its directory. The bytes are unchanged, but the package no longer sits where its tests expect. The check reports the old path `deleted`.
- **R-156-3 (`ailang-run` aliases).** An AILANG program with `FS` can write through a symlink into a protected path (V6). It is detected by §4.5 at the next exec call. That call may never happen: an episode can end with no exec after the write. Upstream ticket (§10). An episode-end check is a possible add-on.
- **R-156-4 (TOCTOU).** A concurrent exec could create a symlink between World's pre-check and policy-tool's open. The window is milliseconds, and §4.5 detects the result.
- **R-156-5 (semantics, not bytes).** Protected test bytes do not stop project code from changing what the tests do at run time (row 140 R-140-3). Config files (`vitest.config.ts`, `pyproject.toml`, `package.json`, `go.mod`) steer which tests run. QUICKSTART's worked patterns include them where the grammar allows it. A directory's config file is protected only by its base name (e.g. `pyproject.toml` at any depth).
- **R-156-6 (gitignored check scope).** `--ignored=matching` reports a protected-matching file in an ignored location (V9); **[r4]** `ls-files -o -i` now does this per file (V20). A protected file inside a nested repository or submodule is out of scope (`--ignore-submodules=all`). **[r4.1]** Ignored legs are judged relative to the episode-start ignored snapshot, not the base commit alone. **Residual:** a change to a pre-existing ignored protected file (one in the snapshot) is not detected by the path-set difference. A content identity was measured and not adopted: on sunholo-platform's 3 259 standing files, `lstat` took 155 ms and sha256 534 ms over 49.8 MB (V23), which does not fit the 500 ms phase. **Budget risk (V23):** on the real, dependency-installed checkouts the ignored leg alone is 177–504 ms warm (sunholo-platform 346 ms, ailang 406–504 ms), against 44–64 ms on V21's clean clone. Added to `status` and `ls-files -v`, a phase on such a tree can exceed 500 ms and then records `error` (fail closed, never intact). The planner must re-measure this on an episode worktree. **[r4.2]** Resolved by quorum r2's fix: `protectedCheckTimeout` = 2 500 ms and `protectedBaselineTimeout` = 3 200 ms, both re-sized from V24 before M6. A content identity (534 ms sha256) still does not fit and is not adopted. **[r4.2] Residual (oc-kimi-k3 (3)):** dependencies installed after `NewExecHandler` are not in the snapshot, so they report as ignored changes until V24 shows M4 episodes install first; this fails toward CHANGED, never toward intact.
- **R-156-7 (Unicode and case across matchers) [r4, rewritten to the measured state].** Git's `icase` folds ASCII only and its `?`/`[…]` match bytes, so git pathspecs cannot be made equal to the mirror's NFC + full-orbit fold (V16, V18), and restricting patterns to ASCII does not help (both V16 misses are ASCII patterns). [r4] stops asking git to agree: git only proposes candidates through a widened pathspec, and `matchProtected` (the same fold as fileguard) decides (§4.2, §4.5). Completeness then rests on two measured facts, each pinned by a test: only U+017F, U+037E, U+1FEF and U+212A fold onto ASCII (V17, AC1c.1), and git's `*` spans multi-byte and NFD spellings (V18, V22). **Residuals:** (a) a future Go/x/text that folds a new rune onto ASCII is caught by AC1c.1 at build time, not at run time; (b) paths that are not valid UTF-8 are UNMEASURED (the mirror's NFC of invalid bytes was not probed); (c) Linux git and `core.precomposeunicode` on a Linux checkout are UNMEASURED on the rig, and AC1c.3 runs on CI's Linux runner. Prevention (fileguard, Seatbelt) is unchanged. **[r4.1]** The prevention matcher is the pinned AILANG binary's fileguard, with its own x/text. Agreement is asserted only on the banked fixture rows (`TestProtectedMatchFixturePinned`), so a Unicode-version skew between that binary and World's toolchain is a residual.
- **R-156-9 (very short patterns widen to the whole tree) [r4.2, judge r1].** A one-character pattern made only of widened characters (`s`, `?`, `é`, `[ab]`) renders as `:(top,glob,icase)**/*`, so its candidate set is every path. It is still correct, because the mirror filters it, and it is still bounded: output past `protectedCheckMaxOutput` (8 MiB) is `error`, never intact. Profiles should not declare such patterns; the grammar does not refuse them.
- **R-156-10 (hard-linked files are unwritable through the write tools) [V25].** Since the v0.54.15 tool move (row 139 D7), `ailang-write`/`ailang-edit` refuse ANY file with more than one hard link, protected or not; V25 shows the refusal does not depend on `fs_deny_write` at all (it fires with an empty list). A tree whose package manager hard-links installed files cannot be edited there through these tools (pnpm stores and uv caches are plausible examples, not measured). This is a property of the pin, recorded here because row 156 measured it; it is not a row-156 change.
- **R-156-8 (cost).** About 0.35 s per call on the largest project (V9 × 2). It is bounded at 500 ms per check. A cold first call can exceed it (row 140 F4) and then records `error`, not a false `intact`. **[r4.2]** Now bounded at 2 500 ms per phase and 3 200 ms at construction (V23; re-sized from V24); a protected profile's command budget shrinks accordingly (§4.5).

## 10. Conflict surface

- **Files:**
  - `host/broker/exec_protected.go` (new), `exec_profile.go`, `exec_settings.go`, `handlers_exec.go`, `exec_probe.go`;
  - `handlers_ailang.go` and `handlers_ailang_run.go` (the `RenderEpisodePolicy` signature);
  - `host/daemon/workspace.go`, `cmd/ailang-worldd/why.go`, `docs/QUICKSTART.md` §10;
  - `.github/workflows/ci.yml` and `design_docs/verification/world-row140-m0/m0_matrix.py` (M0).
- **Rows:**
  - row 157 (dependency digest) also adds an exec result field. It is a disjoint key, so both can land in either order.
  - row 150 M5b / 151 (workbench) and row 184 (landed) share `handlers_exec.go`/`why.go` hunks.
  - re-fetch before merge.
- **Upstream (to file at landing):** `sunholo-data/ailang` issue: "`fileguard.Protection.Check` matches the lexical path; `policytool` write and the FS effect follow in-root symlinks, so `fs_deny_write` is bypassable through a symlink or hard link (measured v0.52.1)". Evidence: `q4_policytool_denywrite.out.txt`. Then an `ailang messages send mission-control` note. **[V25]** Filed as ailang#1873 (2026-10-17), closed as fixed by `db9719bb7` in v0.54.15; World re-measured it on the pinned binary (V25).

## 11. Non-goals (restated)

- Docker or containers.
- A default protected set (D-156-2).
- Failing the episode on a breach (D-156-1).
- Dependency digests (row 157).
- Memory/CPU limits (row 158).
- Proving that every test ran (R-156-5).

## 12. Quorum log

- **Round 1 (2026-10-17): BLOCKED.**
  - **Present:** `gemini-4-argon`. **Absent (unreachable):** `gpt6-1-sol`, `oc-glm-5-3`, `oc-kimi-k3`.
  - **Objection:** §4.5's `git status` runs against the current HEAD and honours `assume-unchanged`/`skip-worktree`, so it does not prove the worktree equals `base`. If HEAD moves, `changes` cannot name the paths that differ from `base`.
  - **Proposed fix:** `diff-index --cached <base>` + `diff-files`, or intact = `branch.oid == base` AND zero entries, with `.git` index writes shown to be denied, and a `head-moved` entry.
  - **Measured (V13):** the premise is real. With either flag set, `status` misses a modified test file, **and so do `diff-index` and `diff-files`**. Under World's srt rendering, sandboxed `update-index` and `commit` are refused (`index.lock: Operation not permitted`) and HEAD does not move. The control (unsandboxed) sets both flags and moves HEAD.
  - **Changed (r2):**
    - §4.5 adopts the reviewer's intact rule and `head-moved` entry verbatim in substance;
    - it adds a `git ls-files -v` flag check (`index-flag`, 36–42 ms), chosen over `diff-index`/`diff-files` because those are blind to the flags too;
    - V13, P7 and P12 are amended;
    - M6 gains `TestExecProtectedIndexFlag` (negative control), `TestExecProtectedHeadMoved` and `TestExecSrtIndexWriteRefused`;
    - §7 gains MUT-NO-INDEX-FLAG-CHECK.

- **Round 2 (2026-10-17): BLOCKED** (same seats: `gemini-4-argon` present; `gpt6-1-sol`, `oc-glm-5-3`, `oc-kimi-k3` absent, unreachable). The round-1 objection was not repeated.
  - **Objection:** pathspecs are CWD-relative unless anchored, and `D/**` does not match `D` itself when the directory is replaced by a file or symlink.
  - **Proposed fix (verbatim):** "In §4.2, change the git pathspec rendering for directory form `D/**` to emit two pathspecs: `:(top,glob,icase)D` and `:(top,glob,icase)D/**` (and `:(top,glob,icase)**/G` for base-name form `G`), and specify in §4.5 that both git invocations (`status` and `ls-files -v`) run with `cmd.Dir = h.worktree`."
  - **Disposition (r3, controller, narrow-refinement carve-out):** the fix is concrete, reviewer-authored and does not dispute the direction, so it is applied verbatim to §4.2's table and §4.5's first sentence, with no third quorum round. The executor's `TestProtectedGitPathspecEquivalence` (M1) must include a "directory replaced by a symlink/file" row and a run from a subdirectory CWD, each killing a mutant that drops `top` or the bare `D` pathspec.

- **r4 (iteration 282, 2026-10-19, designer `opus`): detection redesign after V16; not yet quorum-reviewed.**
  - **Why:** M1c measured that r3's `:(top,glob,icase)` pathspecs miss `teſts/x.py` and `pkg/K_test.go` (V16). Restricting patterns to ASCII cannot fix it.
  - **Measured (V17–V22):** exactly 4 non-ASCII scalar values fold onto ASCII (V17); git `?` and bracket classes match bytes and `*` spans multi-byte runes (V18); widened pathspecs select every V16 counterexample and keep r3's bare-`D`/`top` behaviour (V19); `status --ignored=matching` flags whole excluded directories that hold no match (V20, an r3 false-positive found on the way); NFC literals miss NFD index entries and the widened form does not (V22); the r4 commands cost 165–195 ms per phase on ailang, and the no-pathspec alternative is ≈ 250–290 ms with output that grows with the tree (V21).
  - **Changed:** §4.2 git row and two bullets (candidate generator + mirror authority, the widening `W` and why it is a superset); §4.5 (three invocations: `status --ignored=no`, `ls-files -o -i`, `ls-files -v`; the filter; the output cap and fail-closed rules; cost; why not drop the pathspec); `base_digest` noted unchanged; P7; M1 gains M1c [r4] with a file allowlist and LOC cap; M6 uses M1c's code; §7 gains seven mutations; R-156-7 rewritten. Unchanged: the grammar, the mirror and its fixtures, prevention (§4.3/§4.4), the record schema.
  - **Evidence:** `design_docs/verification/world-row156-r4/` (scripts and outputs; fixture repos deleted).

- **r4 quorum round 1 (iteration 282): BLOCKED** (`gemini-4-argon` reject, `oc-glm-5-3` reject; `gpt6-1-sol`, `oc-kimi-k3` unreachable). Revised as **[r4.1]**.
  - **Gemini objection 1:** multi-rune folds and ligatures (`ß`→`SS`, `ﬀ`) could fold onto ASCII, so `W` would not be a superset. **Premise measured FALSE** for World's mirror by the controller: key(`ß`) = `ß`, key(`ﬀ`) = `ﬀ`, key(`ﬁ`) = `ﬁ`, key(`buﬀ`) = `BUﬀ`, key(`tésts`) = `TÉSTS`, key(`ḋata`) = `ḊATA` (SimpleFold per rune, NFC not NFKC). **Disposition:** the fix's intent is applied as clarification: V17, §4.2 and AC1c.1 state the fold exactly; AC1c.3 gains `a_teßt.go` and `buﬀ/` as mirror-false controls.
  - **Gemini objection 2:** "each run has its own 500 ms bound" contradicted the per-phase budget. **Disposition:** the fix is applied verbatim in §4.5; the contradictory wording is removed.
  - **GLM objection:** ambient ignored dependency trees make a pristine checkout report changes forever. **Premise REAL** (V10's 2 002 `test_*.py` in `.venv`; V23: 3 259 / 165 / 1 240 standing entries). **Disposition:** the fix is applied verbatim: an episode-start ignored snapshot (`base_ignored_digest`, `base_ignored_files`), the set-difference rule, AC1c.5 with MUT-NO-IGNORED-BASELINE, M6's `TestExecProtectedAmbientIgnoredIntact`, and V23 (0 / 0 / 0 after). R-156-6 now states the snapshot basis, the residual (a change to a snapshot file is not detected; content identity measured at 155 ms `lstat` / 534 ms sha256 on sunholo-platform and not adopted) and a **new measured budget risk:** the ignored leg alone is 177–504 ms warm on real dependency-installed checkouts (V23).
  - **GLM catch 3:** the prevention matcher's x/text is the pinned binary's own. **Disposition:** one sentence added to R-156-7.
  - **Unchanged:** everything else in r4.
- **r4 quorum round 2 (iteration 282, artifact `w-exec-protected-test-paths-2026-10-19T07-43-02Z.json`)**: BLOCKED, 3 present (gemini-4-argon, oc-glm-5-3, oc-kimi-k3 as the reserve), `gpt6-1-sol` absent (unreachable); the author's vendor (Claude) sat out. Doc-level round count: 5 (r1, r2, r3 earlier; r4 r1 and r2 now).
  - **Surface per objection (round tracking):** all three objections land on ONE surface, §4.5's phase budget, which V23 contradicted at 500 ms. No objection disputes the r4 direction (candidate generator + mirror filter) or the ignored snapshot.
  - **gemini-4-argon fix (verbatim):** "raise `protectedCheckTimeout` from `500 ms` to `2500 ms`". Applied.
  - **oc-kimi-k3 fix (1):** "set the shared phase budget to ≈2,500 ms … and recompute the child reserve's min() interaction with the 9 s cap and 10 s handler deadline". Applied (§4.5). **Fix (2):** "Gate M6 on a new V-row: all three git invocations per phase on a dependency-installed ailang episode worktree … requiring p95 phase time below the budget with intact ≠ null". Applied as V24 in M6. **Fix (3):** first option taken ("add a verification-log row proving all dependency installation precedes NewExecHandler in M4 episodes"), folded into V24, with the fail-toward-CHANGED residual in R-156-6.
  - **oc-glm-5-3 fix:** (1) V24 as above, applied; (2) a named bound sized "measured max × 1.5", applied to the construction step (`protectedBaselineTimeout` = 3 200 ms from V23's 2 114 ms cold leg); the phase bound takes the other two reviewers' 2 500 ms, which already exceeds V23's cold max; both are re-sized from V24's max × 1.5 before M6; (3) the construction bound is stated in §4.5's baseline paragraph, applied; (4) R-156-6's "left open" replaced by the chosen numbers and M6's timeout test gated on V24, applied.
  - **Disposition (r4.2, controller, narrow-refinement carve-out):** every remaining objection carries a concrete reviewer-authored fix and none disputes the direction, so the fixes are applied verbatim with no third round. M1c (pure candidate renderer + filter + real-git equivalence) does not depend on the bound; M6 is gated on V24.

