# Sprint plan — w-verifygate-forklock-guard-residuals (iteration 229, row 142)

**Design**: `design_docs/planned/w-verifygate-forklock-guard-residuals.md` (r2, `9ad84f2`). **Branch**:
`sprint/w-verifygate-forklock-guard-residuals`, base `74247e8` (on `origin/dev` `2d3a255`). **Scope**:
`host/verifygate/forklocked_write_test.go` + the one `ci.yml` step. Test-only: no production code, no
`.ail`, no `tools/launchd/*`, no `scripts/*` change. **Size**: 3 milestones, ~1 executor session.
Commit each milestone separately; every drill reverts to that milestone's commit and proves byte identity.

## 0. Re-measured at `74247e8` (planner, 2026-10-03, go1.26.6 darwin/arm64)

| Fact | Command | Output |
|---|---|---|
| Pin | `~/.pinned-ailang/ailang --version` | `AILANG v0.41.0` |
| Wrapper signatures (V5d) | `cat -n host/verifygate/forklocked_write_test.go` | l.36 `forkLockedDo`, l.55 `forkLockedWrite`, l.60 `writeFileForkLocked`, l.69 `copyFileForkLocked`, l.91 `createTempForkLocked`, l.113 `createForkLocked` — **match the design exactly** |
| Open closures (V10) | same | l.56 single `return os.OpenFile(path, flag, mode)`; l.93–98 multi-statement `os.CreateTemp` closure — as the design says |
| Scan shape (V3/V5) | same | l.183 `bannedWriteFuncs`, l.186–192 `wrapperFuncs` (by name), l.197–201 `rawWriteAllowlist`, l.248 `Info{Uses}` only, l.264–296 CallExpr-only visit; 5 positives l.331–335, negative l.357 (defines its own `writeFileForkLocked`) |
| CI step | `sed -n 196,204p .github/workflows/ci.yml` | l.196–199 comment, l.200 `- name: Fork-locked write tests, verbose…`, l.201 `timeout-minutes: 3`, l.202 `run: >-`, l.203–204 `go test … -run '^(…3 tests…)$'`. No `defaults.run.shell` in ci.yml → GitHub's `bash -eo pipefail` |
| Package baseline | `go test -json -count=1 ./host/verifygate/` → count `Action` per `Test` | **150 pass (59 top-level), 1 skip, 0 fail**, 141.5 s. Artifact `~/.ailang/state/w229_vg_base_plain.json` |
| D3 sweep prototype (AST, selector on `X` ident name, run from `host/verifygate` over `../../host ../../cmd`; scratch, deleted) | throwaway `go run` | **26 dirs, 13 write+fork dirs** (V9's 12 + `host/boundary`, whose `var rawWrite = os.WriteFile` value use the V9 regex `\(` missed); `t.Parallel()` only in `host/verifygate` (10); 0 parse errors; no `testdata/*.go` under host/cmd |
| `umask` | `umask` | `022` (rig); the probe asserts modes umask-independently (§1.3) |

## 1. Plan overrides (design text under-specified or inconsistent; see §6 for the controller list)

1. **Wrapper-identity check is live-only.** D2 says the scan "fails with a scan error if any wrapper is
   missing or not declared in `forklocked_write_test.go`". Every fixture package lacks the wrappers, so
   applied unconditionally it errors all 13 fixtures. Plan: the check runs only when `files == nil`
   (live dir scan); fixtures skip wrapper counting.
2. **Fixtures that need `forkLockedDo` define it.** P-Smuggle and N-OpenClosure each carry
   `func forkLockedDo(open func() (*os.File, error), fill func(*os.File) error, hold func()) error { return nil }`
   in the same fixture source, so `pkg.Scope().Lookup("forkLockedDo")` resolves.
3. **Mode assertion is umask-independent.** 0o755 wrappers: `perm&0o111 == 0o111`; `createTempForkLocked`:
   `perm == 0o600`; `createForkLocked`: `perm&0o111 == 0` and size 0. Content via `os.ReadFile`.
4. **AC4 mutant is uncompilable as written** (`wf(p, …)`: no `p` in scope). Use `wf("x", nil, 0o755)`.
5. **AC5 mutant needs an import**: `toolchain_pin_gate_test.go` does not import `syscall` (V5b); the
   drill adds it and the revert removes it.
6. **AC10's "N-OpenClosure is reported" is unobservable in the same run**: the live-violation `Fatalf`
   (l.316–318, kept first) stops the test before fixtures. AC10 fires on the live red only; a second
   observation (N-OpenClosure) needs the fixture loop to run before the live check — not done; the
   live red naming both real closures is sufficient (both are exempt only through that clause).
7. **Floors**: D3 floor stays `>= 12` write+fork dirs (design); measured 13 is recorded. `>= 26` dirs.
8. **Full gate**: final `go test ./... -count=1` (+ `verify_ail.sh`) per controller, not `verify_go.sh`.

## 2. Common gate vocabulary

```
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter229
export AILANG_BIN=$HOME/.pinned-ailang/ailang; $AILANG_BIN --version      # must print AILANG v0.41.0
FENCE: go vet ./host/verifygate/                                          # go build skips _test.go
FAST:  go test ./host/verifygate/ -count=1 -v -run '^(TestForkLocked.*|TestKernelRefusesExecOfWriterOpenFile|TestVerifygateTestWritesAreForkLocked|TestNoParallelWriteForkPackagesOutsideVerifygate)$'
RACE:  go test -race ./host/verifygate/ -count=1 -run 'TestForkLocked|TestVerifygateTestWritesAreForkLocked|TestNoParallel'
PKG:   go test ./host/verifygate/ -count=1      # ≈150 s;  COUNT: same with -json > ~/.ailang/state/w229_vg_<M>.json
AIL:   ./scripts/verify_ail.sh                  # needs AILANG_BIN exported
DRILL: H=$(shasum -a 256 <file>); apply mutant; run named test → must FAIL with the named text;
       git checkout -- <file>; git diff --quiet && [ "$(shasum -a 256 <file>)" = "$H" ]
```
Bank each drill's failing output at `~/.ailang/state/w229_drill_<AC>.log` (`/tmp` is wiped).

## M1 — wrapper bodies guarded by the after-Close probe (kills E1, E2, E3)

File `host/verifygate/forklocked_write_test.go`; imports gain `sync/atomic`.
- `var forkLockedProbe atomic.Pointer[func(stage string, f *os.File)]` with a doc comment: test-only;
  runs on the RLock-holding goroutine; **must not fork** (R1 deadlock).
- `forkLockedDo` (defer `RUnlock` unchanged): after fill/hold, `if p := forkLockedProbe.Load(); p != nil { (*p)("filled", f) }`;
  then `err := f.Close()`; then the same call with `"closed"`; `return err`. Fill-error path unchanged.
- New serial `TestForkLockedWrappersHoldLockThroughClose` (no `t.Parallel`, no `t.Run`): one `t.TempDir()`
  root; probe ignores `f.Name()` outside it; per stage records the stage name and checks
  `filled`: `!syscall.ForkLock.TryLock()` (else `Unlock` + record "not read-held") and `f.Stat().Size()==len(want)`;
  `closed`: `!TryLock()` (else `Unlock` + record "released") and `errors.Is(statErr, os.ErrClosed)`.
  `Store(&fn)` + `t.Cleanup(func(){ forkLockedProbe.Store(nil) })`. Drives, each in its own subdir:
  `writeFileForkLocked(p, data, 0o755)`; `copyFileForkLocked(src, dst, "rel", 0o755)` with `src` from
  `os.Open` of a file written by `writeFileForkLocked` **before** the probe is installed;
  `createTempForkLocked(dir, "p*", data)`; `createForkLocked(p)`;
  `forkLockedWrite(p, O_WRONLY|O_CREATE|O_TRUNC, 0o755, fill, nil)`. Resets the stage log per call.
  Failure texts (`t.Errorf`, so all wrappers report): `"%s: fork-lock probe stages = %v, want [filled closed]"`,
  `"%s: ForkLock not read-held during fill"`, `"%s: ForkLock released before Close returned"`,
  `"%s: probe at closed stage saw an open file"`, plus content/mode mismatches (§1.3).
- `TestForkLockedWriteBlocksConcurrentFork` l.123: `0o644` → `0o755`, nothing else.

**Gate M1**: FENCE; FAST (2 `--- PASS` ForkLocked tests, kernel `--- SKIP`, scan PASS); RACE (rc 0, no
`WARNING: DATA RACE`); COUNT → **151 pass + 1 skip, 0 fail**; AIL rc 0.
**Drills** (on the M1 commit, `-run '^TestForkLockedWrappersHoldLockThroughClose$' -count=1 -timeout 60s`):
- **AC3 (E3)**: in `forkLockedDo` delete `defer syscall.ForkLock.RUnlock()` and insert
  `syscall.ForkLock.RUnlock()` immediately before `err := f.Close()` → fires
  `ForkLock released before Close returned` for all five; then `-run '^TestForkLockedWriteBlocksConcurrentFork$'`
  and `-run '^TestVerifygateTestWritesAreForkLocked$'` on the same mutant stay **PASS** (only D1 sees E3).
- Probe half of **AC1 (E1)**: replace `createTempForkLocked`'s body with raw `os.CreateTemp` + `Write` +
  `Close` → fires `createTempForkLocked: fork-lock probe stages = [], want [filled closed]`.
- Probe half of **AC2 (E2)**: in `forkLockedWrite`, `if mode&0o111 != 0 { f, err := os.OpenFile(path, flag, mode); …fill; return f.Close() }`
  → fires `writeFileForkLocked: fork-lock probe stages = [], want [filled closed]`.
Revert + byte identity after each.

## M2 — Uses-based, identity-matched scan with the narrowed exemption (kills E5b, E8b, E6b, F6)

Same file.
- **Refactor** `createTempForkLocked` (V10): open closure = `func() (*os.File, error) { return os.CreateTemp(dir, pattern) }`;
  `name = f.Name()` becomes the first statement of `fill`. Cleanup/return logic unchanged.
- **Doc comment** on `forkLockedDo` (verbatim from design M2): "`scanForkLockedWrites` hardcodes this
  function's first argument (`Args[0]`) as the exempt open closure. Do not reorder the signature."
- **`scanForkLockedWrites`** per D2: `types.Info{Uses, Defs}`; keep `conf.Check` and the error path.
  1. `pkg.Scope()` lookups: `forkLockedDo`, `TestKernelRefusesExecOfWriterOpenFile`, the 5 wrappers.
  2. AST pass per file: collect `exempt map[*ast.Ident]bool` — for each `CallExpr` whose Fun ident
     `Uses`-resolves to the `forkLockedDo` object, if `Args[0]` is a `*ast.FuncLit` whose last stmt is a
     `ReturnStmt` with one `CallExpr` result whose Fun Sel resolves to `os.OpenFile`/`os.CreateTemp`, mark
     that Sel ident. Record FuncDecl ranges `[Pos,End)` with names, and the kernel-control body range
     (identity via `info.Defs[fd.Name] == kernelObj`).
  3. Iterate **all** `info.Uses`; banned = `*types.Func`, `Pkg().Path()=="os"` ∧ name ∈ {WriteFile, OpenFile,
     Create, CreateTemp, NewFile}, or `"syscall"` ∧ name ∈ {Open, Openat, Creat}. Skip if `exempt[id]` or
     inside the kernel-control body. `where` = enclosing FuncDecl name or `package-level initializer`.
     Message unchanged: `raw <pkg>.<Name> in <where>: use a fork-locked wrapper`.
  4. Import bans: aliased/dot `"os"` (unchanged text) + any import of `"io/ioutil"` or
     `"golang.org/x/sys/unix"` (`banned import %s evades the write scan`).
  5. Live mode only (§1.1): `wrapperCalls` += each `Uses` entry equal to a wrapper object outside
     `wrapperFile`; scan error if any wrapper object is nil or its `Pos()` file ≠ `wrapperFile`.
  6. `sort.Strings(rep.violations)`. Delete `rawWriteAllowlist` and `wrapperFuncs`.
- **Fixtures** (append to the `positives` table; same `want` prefix matching): P-FuncValue, P-RawFd (both
  lines), P-PkgInit (want text also checked for `package-level initializer`), P-NamedWrapper, P-Ioutil
  (ImportSpec line), P-RootWrite, P-Smuggle (`want` = the `os.WriteFile` line; additionally assert the
  `os.OpenFile` opener line is **absent**: `"fixture Smuggle: exempt opener %s reported"`). N-OpenClosure
  joins the negative check (zero violations). Sources exactly as the design table + §1.2.
- Contingency (D2): if the first live run reports a pre-existing hit, STOP, report it, and choose
  identity-exemption or a follow-up row — never weaken the ban silently. (V5c predicts zero.)

**Gate M2**: FENCE; FAST; RACE; COUNT → **151 pass + 1 skip, 0 fail** (no new top-level test);
live floors unchanged (≥9 files, ≥21 wrapper sites, identity-counted); AIL rc 0.
**Drills** (on the M2 commit, `-run '^TestVerifygateTestWritesAreForkLocked$' -count=1`):
- **AC1 scan half**: E1 body (as M1) → `raw os.CreateTemp in func createTempForkLocked`.
- **AC2 scan half**: E2 branch (as M1) → `raw os.OpenFile in func forkLockedWrite`.
- **AC4 (E5b)**: append `func h() { wf := os.WriteFile; _ = wf("x", nil, 0o755) }` to
  `ail_binary_gate_test.go` → names `ail_binary_gate_test.go:<line>: raw os.WriteFile in func h`.
- **AC5 (E8b)**: add `"syscall"` import + P-RawFd body as `func h()` in `toolchain_pin_gate_test.go` →
  names both the `syscall.Open` and `os.NewFile` lines.
- **AC6 (E6b)**: in the scan, `if encl == "" { continue }` before reporting → `fixture PkgInit: … not reported`.
- **AC7 (F6)**: skip reporting when the enclosing FuncDecl name is `forkLockedWrite` → `fixture NamedWrapper: … not reported`.
- **AC10**: delete the exempt-ident collection (step 2's marking) → live red naming the
  `forkLockedWrite` (l.~56) and `createTempForkLocked` opener lines (§1.6).
- **AC12**: insert `_ = os.WriteFile(filepath.Join(dir, "y"), nil, 0o644)` as the first statement of
  `createTempForkLocked`'s open closure → names that line `in func createTempForkLocked`; the
  `os.CreateTemp` line is **not** in the output (grep the banked log).
Re-run the M1 drills' probe halves is not required (M1 code untouched except the refactor; FAST covers it).

## M3 — R2 tripwire + CI SKIP fail

- **`TestNoParallelWriteForkPackagesOutsideVerifygate`** (serial) in the same file. Split into
  `loadGoDirs(roots ...string) (map[string]map[string][]byte, error)` (walk, `os.ReadFile`, keys with the
  `../../` prefix trimmed) and pure `parallelWriteForkViolations(dirs) (viol []string, nDirs, nWriteFork int, err error)`
  (parse-only `parser.ParseFile`; selector `X` ident name match per D3; `Parallel` = `SelectorExpr`
  call with `len(Args)==0` in a `_test.go`; skip dir `host/verifygate`; sorted). Message:
  `<dir>/<file>:<line>: t.Parallel in a write+fork package outside the fork-locked scan`.
  Live: `loadGoDirs("../../host", "../../cmd")`; `nDirs < 26` → `"sweep saw %d Go dirs, want >= 26"`;
  `nWriteFork < 12` → `"sweep saw %d write+fork dirs, want >= 12"`; any violation → fatal list.
  Known positive `{"zz/x_test.go": t.Parallel + os.WriteFile + exec.Command}` must be named; known
  negative = same without `t.Parallel()` → zero. (Fixture source as strings; no banned AST use.)
- **`ci.yml` l.200–204**: keep the l.196–199 comment (append "row 142: PASS loop, fail-loud linux
  precondition"), change `run: >-` + body to the design's `run: |` block verbatim (5-test list).
- Local dry-run of the step body (known positives, record rc in the PR):
  `ruby -ryaml -e 's=YAML.load_file(".github/workflows/ci.yml")["jobs"]["go-verify"]["steps"].find{|x|(x["name"]||"").start_with?("Fork-locked")};File.write(File.expand_path("~/.ailang/state/w229_step.sh"),s["run"])'`
  (PyYAML is **not** installed on the rig; ruby's YAML is — measured)
  then `RUNNER_TEMP=$(mktemp -d) bash ~/.ailang/state/w229_step.sh` → **rc 1** `::error::this step must run on linux`;
  `RUNNER_OS=Linux RUNNER_TEMP=$(mktemp -d) bash …` → **rc 1** `::error::TestKernelRefusesExecOfWriterOpenFile did not PASS on linux (SKIP?)` (darwin SKIP; the two earlier listed tests PASS).

**Gate M3**: FENCE; FAST (4 PASS + kernel SKIP); RACE; COUNT → **152 pass + 1 skip, 0 fail**; AIL rc 0;
`./scripts/check_no_personal_email.sh` rc 0; `ruby -ryaml -e 'YAML.load_file(".github/workflows/ci.yml")'` rc 0;
then **final** `go test ./... -count=1` (AILANG_BIN exported) rc 0.
**Drills** (on the M3 commit, `-run '^TestNoParallelWriteForkPackagesOutsideVerifygate$'`):
- **AC9 (R2)**: add `t.Parallel()` as the first line of one test in `host/broker/broker_test.go` → fires
  `host/broker/broker_test.go:<line>: t.Parallel in a write+fork package outside the fork-locked scan`.
- Floor drill: change the live root `"../../host"` to `t.TempDir()` → fires `sweep saw <n> Go dirs, want >= 26`.
- Positive drill: drop `Parallel` from the matcher (e.g. match `"Parallel_"`) → known positive not named.
Revert + byte identity after each.

## 5. CI-only acceptance (PR time; the rig has no linux runtime)

- **AC8**: scratch PR commit `if true { t.Skip("drill") }` atop `TestKernelRefusesExecOfWriterOpenFile`
  → step fails with `::error::TestKernelRefusesExecOfWriterOpenFile did not PASS on linux (SKIP?)`;
  revert commit; `git diff <pre-drill-sha> HEAD --stat` empty. Record job ids.
- **AC11 / first linux run**: PR-head job `go host build + test gate` log (V11's `gh api … | grep`,
  names widened to all 5) shows a top-level `--- PASS:` for each; `-race` leg green.
- Contingency (glm r2, verbatim in design M3): kernel control SKIP/FAIL on the runner → HALT, escalate
  finding 5; never weaken the loop.
- Re-fetch `origin/dev` before merge (row 135/138 follow-ups may touch `ci.yml`).

## 6. Design discrepancies for the controller (the design is not edited)

1. D2 wrapper-identity scan error applies to fixtures too as written → would error every fixture (§1.1).
2. AC4 mutant references undefined `p` (§1.4); AC5 omits the needed `syscall` import (§1.5).
3. AC10 claims N-OpenClosure is reported in the same run, but the live `Fatalf` fires first (§1.6).
4. V9's "12 write+fork dirs" undercounts the D3 AST sweep, which finds 13 (`host/boundary` value use);
   the floor of 12 still holds.
5. M3 gate names `verify_go.sh`; the controller's gate is `go test ./...` + `verify_ail.sh` (§1.8).
6. M1 "content and mode match" is umask-sensitive as literally read; plan pins umask-free checks (§1.3).
