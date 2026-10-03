# Sprint plan — w-verifygate-etxtbsy (iteration 228)

**Design**: `design_docs/planned/w-verifygate-etxtbsy.md` (option C′). **Branch**: `sprint/w-verifygate-etxtbsy`,
base `42f3dc1` (design commit on `c2476ff`). **Scope**: `host/verifygate/*_test.go` + one `ci.yml` step.
No production code, no `.ail`, no `tools/launchd/*`, no `scripts/*` semantics change.
**Size**: 3 milestones, ~1 executor session. Commit each milestone separately (mutation drills
revert to the milestone commit and prove byte identity against it).

## 0. Re-measured at HEAD `42f3dc1` (planner, 2026-10-03, go1.26.6 darwin/arm64, 16 cores)

| Fact | Command | Output |
|---|---|---|
| Pinned binary | `~/.pinned-ailang/ailang --version` | `AILANG v0.41.0`, commit `24ee108` (= CI's pin: ci.yml l.77–108/134 downloads release v0.41.0 to `$HOME/.pinned-ailang/ailang`; `verify_go.sh` l.154–181 rejects unset `AILANG_BIN` or a version token ≠ `v0.41.0`) |
| verifygate is 8 files, all `_test.go` | `ls host/verifygate/*.go \| wc -l`; `… \| grep -v _test` | `8`; empty |
| 21 banned write calls | `grep -n 'os\.WriteFile(\|os\.OpenFile(\|os\.Create(\|os\.CreateTemp(' host/verifygate/*.go` | 21 lines: ail_binary 376, 400, 464(CT), 563, 569(CT); module_manifest 42(OF), 279, 289, 315, 349, 392, 432, 478; evidence_manifest 75, 113, 134, 197; mission_config 111; toolchain_pin 484, 504, 529(Create). Per file W/OF/C/CT = ail 3/0/0/2, evid 4/0/0/0, mission 1/0/0/0, module 7/1/0/0, toolchain 2/0/1/0; dispatch_lever, floor_raise_inventory, subprocess_sink 0 → **matches V11** |
| Baseline pass events | `AILANG_BIN=… go test -json -count=1 ./host/verifygate/` → count `Action=pass` with `Test` | **148** pass (57 top-level), 0 fail, 0 skip → matches the `e6064b6` figure. Artifact `~/.ailang/state/w228_vg_base_plain.json` |
| Local verifygate timing (pre-fix) | same run; `go test -count=1 -race -timeout 8m ./host/verifygate/` | plain **146.8 s** (wall 147 s); race **145.0 s** (wall 146 s). CI pre-fix baseline (V14): plain 179.3 s, race 215.1 s |

## 1. Doc premises found false or under-specified (plan overrides)

1. **Core signature cannot express `CreateTemp`.** `forkLockedWrite(path, flag, mode, fill, hold)` is
   path-based; `os.CreateTemp` generates the name. Plan: an internal core
   `forkLockedDo(open func() (*os.File, error), fill func(*os.File) error, hold func()) error`
   holds the lock; `forkLockedWrite` is a closure over `os.OpenFile`; `createTempForkLocked` a
   closure over `os.CreateTemp`. The scan allowlist is the set of enclosing **FuncDecls**
   {`forkLockedWrite`, `createTempForkLocked`, `TestKernelRefusesExecOfWriterOpenFile`} (calls
   inside func literals attribute to the enclosing FuncDecl).
2. **"keeps its error text"** is not automatic: the core returns one undifferentiated error. Plan:
   `copyFileForkLocked` tracks whether `fill` ran (closure bool + captured `io.Copy` error) and
   returns `create copy target %s: %v` (open failed) vs `copy %s: %v` (fill failed), byte-identical
   to `module_manifest_gate_test.go` l.44/48. `open copy source` (l.36) stays outside the lock.
3. **AC4 file floor is inconsistent**: design text says `< 9`, AC table says `files < 8`. After M1
   there are 9 files; **floor = 9**.
4. **AC4 mutant "drop the wrapper floor to 0" is not a detecting mutant** (it silences the floor).
   Replaced by two mutants that must fire (§M3).
5. **AC2 line `:315` will move** after M2 migrates the file; the mutant targets
   `mutateCopiedScript`'s write call wherever it lands, and the expected message names that line.
6. Local timing cannot discharge AC7 (16-core rig, 145 s race); AC7 is CI-only (2-core runner).

## 2. Common gate vocabulary

```
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter228
export AILANG_BIN=$HOME/.pinned-ailang/ailang          # must print AILANG v0.41.0
FENCE:  go vet ./host/verifygate/ && go test -run '^$' ./host/verifygate/   # go build ./... skips _test.go
NEW:    go test ./host/verifygate/ -count=1 -v -run '^(TestForkLockedWriteBlocksConcurrentFork|TestKernelRefusesExecOfWriterOpenFile|TestVerifygateTestWritesAreForkLocked)$'
RACE:   go test -race ./host/verifygate/ -count=1 -run '^(TestForkLockedWriteBlocksConcurrentFork|TestKernelRefusesExecOfWriterOpenFile|TestVerifygateTestWritesAreForkLocked)$'
COUNT:  go test -json -count=1 ./host/verifygate/ > ~/.ailang/state/w228_vg_<M>.json  (+ the §0 python counter)
FULL:   ./scripts/verify_go.sh && ./scripts/verify_ail.sh        # both need AILANG_BIN exported
DRILL:  apply mutant → run named test (must FAIL with the named message) → git checkout -- <file>
        → git diff --stat (must be EMPTY) → shasum -a 256 <file> equals the pre-mutant hash.
```
Bank every drill's failing output under `~/.ailang/state/w228_drill_<AC>.log` (`/tmp` is wiped).

## M1 — wrappers and mechanism tests

**New file** `host/verifygate/forklocked_write_test.go` (package `verifygate`; stdlib only):
- `forkLockedDo(open, fill, hold)` — `syscall.ForkLock.RLock(); defer syscall.ForkLock.RUnlock()`,
  then `open()`, `fill(f)` (if non-nil), `hold()` (if non-nil), `f.Close()`; on fill error close
  and return it. Doc comment: nothing in the region may fork (deadlock: `acquireForkLock` →
  `ForkLock.Lock`); `hold` is a test-only seam, `nil` in all production-of-tests callers.
- `forkLockedWrite(path string, flag int, mode os.FileMode, fill func(*os.File) error, hold func()) error`.
- `writeFileForkLocked(path string, data []byte, mode os.FileMode) error` — `O_WRONLY|O_CREATE|O_TRUNC`
  (exact `os.WriteFile` semantics; `os.WriteFile` does not chmod an existing file, neither does this).
- `copyFileForkLocked(src *os.File, dst, rel string, mode os.FileMode) error` — `O_CREATE|O_EXCL|O_WRONLY`, `io.Copy`, error texts per §1.2.
- `createTempForkLocked(dir, pattern string, data []byte) (string, error)` — on fill/close error removes the file and returns `""`.
- `createForkLocked(path string) error` — `O_RDWR|O_CREATE|O_TRUNC`, 0o666, empty fill (the `os.Create` shape).
- `TestForkLockedWriteBlocksConcurrentFork` (serial, all OS): temp path; `hold` closes `inside`
  then waits on `release`; after `<-inside` a goroutine runs `exec.Command("/bin/sh","-c","exit 0").Start()`
  and sends its error on `started`. Assert no receive within 250 ms → else
  `t.Fatal("Start returned while a fork-locked write fd was open")`; `close(release)`; assert receive
  within 5 s with nil err → else `t.Fatal("Start did not return within 5 s of release")`; `Wait()` the cmd.
- `TestKernelRefusesExecOfWriterOpenFile` (serial): `runtime.GOOS != "linux"` → `t.Skip("ETXTBSY is
  a Linux kernel rule; darwin permits exec of a writer-open file (design V7)")`. Else raw
  `os.OpenFile(path, O_CREATE|O_WRONLY|O_TRUNC, 0o755)`, write `#!/bin/sh\nexit 0\n`, keep fd open,
  `exec.Command(path).Start()` must satisfy `errors.Is(err, syscall.ETXTBSY)` else
  `t.Fatalf("expected ETXTBSY exec of a writer-open file, got %v", err)`; then `Close`, and the same
  exec must `Run()` nil. (Serial tests never overlap top-level `t.Parallel()` tests, so the raw fd is safe.)

**Gate M1**: FENCE; NEW (expect test 1 `--- PASS`, test 2 `--- SKIP` on darwin); RACE; COUNT
(expect **149 pass + 1 skip**, 0 fail).
**Drills** (on the M1 commit; run only `-run '^TestForkLockedWriteBlocksConcurrentFork$' -timeout 60s`):
- **AC1**: delete the `RLock`/`RUnlock` pair in `forkLockedDo` → fires `Start returned while a fork-locked write fd was open`. Revert, byte identity.
- **AC1b**: delete only the deferred `RUnlock` → fires `Start did not return within 5 s of release` (must fail, not hang to `-timeout`). Revert, byte identity.

## M2 — migrate all 21 write sites (assertions, arms, error texts unchanged)

| File | Sites → wrapper |
|---|---|
| `ail_binary_gate_test.go` | 376 (`writeExecutable`), 400, 563 → `writeFileForkLocked`; 464–475 and 569–580 (`CreateTemp`+`WriteString`+`Close`) → `name, err := createTempForkLocked(dir, pattern, []byte(mutant))`, keep `defer os.Remove(name)` and the following `Chmod` |
| `module_manifest_gate_test.go` | 42 `copyGateFileErr` body → `copyFileForkLocked(in, dst, rel, mode)` (keeps `open copy source` + `MkdirAll` outside); 279, 289, 315 (`mutateCopiedScript`), 349, 392, 432, 478 → `writeFileForkLocked` |
| `evidence_manifest_gate_test.go` | 75, 113, 134, 197 → `writeFileForkLocked` |
| `mission_config_gate_test.go` | 111 (`writeFile`) → `writeFileForkLocked(path, raw, mode)` |
| `toolchain_pin_gate_test.go` | 484 (`writeExecutableAt`), 504 → `writeFileForkLocked`; 529 (`assertWritableSentinel`) → `createForkLocked(sentinel)`, keep the `t.Fatalf` text, drop the now-unneeded `file.Close` block |

Remove imports that become unused (`io` in module_manifest only if no other use; check with FENCE).
**Gate M2**: FENCE; instrument-health control (not load-bearing):
`grep -n 'os\.WriteFile(\|os\.OpenFile(\|os\.Create(\|os\.CreateTemp(' host/verifygate/*.go` → only
`forklocked_write_test.go` lines (forkLockedWrite, createTempForkLocked, test 2); `grep -c` of
wrapper calls outside that file = **21**; COUNT → **149 pass + 1 skip, 0 fail** (no arm lost);
FULL rc 0; record verifygate plain/race seconds from the FULL log.

## M3 — go/types scan guard + ci.yml verbose step

**Add to** `forklocked_write_test.go`:
- `scanForkLockedWrites(dir string, files map[string][]byte) (scanReport, error)` — when `files` is nil,
  enumerates **all** `*.go` in `dir` via `os.ReadDir`; parses with `go/parser` (ParseComments not needed);
  `types.Config{Importer: importer.ForCompiler(fset, "source", nil)}.Check("verifygate", …)` with
  `Info.Uses`; a type error → returned error. Violations: (a) any `ImportSpec` for `"os"` with
  `Name != nil` (alias or dot); (b) every `*ast.CallExpr` whose callee ident/selector resolves in
  `Uses` to a `*types.Func` with `Pkg().Path()=="os"` and `Name()` ∈ `bannedWriteFuncs`
  {WriteFile, OpenFile, Create, CreateTemp}, unless the enclosing FuncDecl is in the §1.1 allowlist;
  a banned call with no enclosing FuncDecl is a violation. Report entries are `file:line: msg`.
  Also counts `filesSeen` and `wrapperCalls` (calls to the 4 wrappers + `forkLockedWrite` outside
  `forklocked_write_test.go`).
- `TestVerifygateTestWritesAreForkLocked` (serial; fixtures in a table loop, **no `t.Run`** so the
  event count is +1): (1) live scan of `"."`: error → fatal; `filesSeen < 9` →
  `t.Fatalf("scan enumerated %d files, want >= 9", …)`; `wrapperCalls < 21` →
  `t.Fatalf("scan found %d wrapper call sites, want >= 21", …)`; any violation → fatal listing them.
  (2) five known-positive fixture packages (string sources passed via `files`), one per banned func
  plus `import osw "os"` + `osw.WriteFile(...)` (must yield both the ImportSpec and the call
  violation); each must be reported at its file:line, else `t.Fatalf("fixture %s: %s not reported", …)`.
  (3) known negative: a wrapper call + `os.ReadFile` → zero violations.

**`.github/workflows/ci.yml`**: insert after the `go build + test gate` step (l.191–194), before the
row-24 step, with a comment citing this design and `verify_go.sh` having no `-v` (V16):
```
      - name: Fork-locked write tests, verbose (w-verifygate-etxtbsy linux gate)
        timeout-minutes: 3
        run: >-
          go test ./host/verifygate/ -count=1 -p 1 -timeout 120s -v
          -run '^(TestForkLockedWriteBlocksConcurrentFork|TestKernelRefusesExecOfWriterOpenFile|TestVerifygateTestWritesAreForkLocked)$'
```
**Gate M3**: FENCE; NEW (2 PASS + 1 SKIP on darwin); RACE; COUNT → **150 pass + 1 skip, 0 fail**;
FULL rc 0; `./scripts/check_no_personal_email.sh` rc 0; `python3 -c 'import yaml,sys;yaml.safe_load(open(".github/workflows/ci.yml"))'` rc 0.
**Drills** (on the M3 commit; run only `-run '^TestVerifygateTestWritesAreForkLocked$'`):
- **AC2**: change `mutateCopiedScript`'s `writeFileForkLocked(path, []byte(mutant), 0o755)` back to
  `os.WriteFile(...)` → fires naming `module_manifest_gate_test.go:<that line>`. Revert, byte identity.
- **AC3**: remove `"CreateTemp"` from `bannedWriteFuncs` → fires `fixture …: CreateTemp … not reported`. Revert, byte identity.
- **AC4a**: change the live scan dir `"."` to `t.TempDir()` → fires `scan enumerated 0 files, want >= 9`.
  **AC4b**: rename the counted wrapper name set to an unused name → fires `… 0 wrapper call sites, want >= 21`. Revert each, byte identity.
- **AC8**: in `mission_config_gate_test.go` add import `osx "os"` and
  `func aliasedEvasion(p string) error { return osx.WriteFile(p, nil, 0o644) }` → fires naming the
  aliased ImportSpec line **and** the call line. Revert, byte identity.

## CI-only acceptance (PR time; rig has no Linux runtime, V7)

- **AC5** (labelled scratch drill): on a throwaway commit, move test 2's `Close` before the
  first exec → CI verbose step shows `expected ETXTBSY …` FAIL; then a revert commit; final PR
  head `git diff <pre-drill-sha> HEAD --stat` empty. Record job ids.
- **AC6**: "go host build + test gate" and the new verbose step green on PR head **and** merge
  commit; verbose log shows `--- PASS: TestKernelRefusesExecOfWriterOpenFile` (not SKIP) → fill V18
  in the design doc with job id + log line.
- **AC7**: from the go-gate job log, `-race` leg (`go test ./... -count=1 -race -timeout 8m`) completes
  < 600 s (no `timed out after 600s`); verifygate race time within ±15% of 215 s (183–247 s);
  record plain/race verifygate seconds next to the local 146.8 s / 145.0 s pre-fix figures.
- A green CI run is not proof of the fix (1 red in 6 runs); AC1–AC5 carry the proof.

## Risks

- Re-fetch `origin/dev` before merge (row 135 follow-ups touch verifygate tests); re-run the §0 grep
  on the merged tree — any new raw write is caught by the scan anyway.
- AC1b leaves a goroutine stuck in `Start`; run it isolated with `-timeout 60s`, never in the full package.
