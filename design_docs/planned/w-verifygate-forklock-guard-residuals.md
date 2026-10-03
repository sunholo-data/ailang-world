# w-verifygate-forklock-guard-residuals

**Status**: PLANNED · **Queue row**: 142 · **Clause**: clause-2 (MET — hygiene; rule (e) position)
**Author lane**: claude:claude-opus-5-5 (iteration 229 designer) · **Date**: 2026-10-03 · **Base**: `origin/dev` `2d3a255`
**Parent**: [w-verifygate-etxtbsy](../implemented/w-verifygate-etxtbsy.md) (fix `382fd7e`; judge
verdict [evaluator-r1](../verification/world-iter228/evaluator-r1.md), PASS 88, mutants E1, E2, E3, E5b, E8b, E6b SURVIVED).
**Scope**: `host/verifygate/forklocked_write_test.go` plus one `ci.yml` step. Test-only. No production
code, no `.ail`, no `tools/launchd/*`, and no change to gate-script semantics. Estimate ~0.5 d.

## Problem

The iter-228 fix is correct, but its guard has measured holes. The judge's mutants survive in
five places:

1. **The wrapper bodies have no guard.** `rawWriteAllowlist` exempts the whole of `forkLockedWrite`
   and `createTempForkLocked` by name (V3). Test 1 drives only `forkLockedWrite` at 0o644, and
   its `hold` seam runs before `Close` (V4). So the following mutants survive:
   - E1: `createTempForkLocked` writes without the lock.
   - E2: the lock is skipped for executable modes.
   - E3: `RUnlock` runs before `Close`.
2. **The scan sees direct calls only** (V5). These mutants survive:
   - E5b: `wf := os.WriteFile` inside a new helper.
   - E8b: `syscall.Open` + `os.NewFile` inside a new helper.
3. **No fixture covers package-level initializers** (E6b).
4. **The CI verbose step passes on a SKIP** of the Linux kernel control (V6).
5. **Residual R2.** Other packages write and exec scripts but sit outside the scan (V9).

Judge finding 6 (bare-name matching) is fixed here because it costs little. Finding 5 (a kernel
that drops the write-deny rule) is not changed: that red is correct and loud.

## Verification Log

Measurements were taken in `/Users/voightkampff/dev/sunholo-data/.wt-world-iter229` at `2d3a255`,
on go1.26.6 darwin/arm64 with `AILANG_BIN` set to the v0.41.0 pin.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Module Go version | `grep -n 'go 1' go.mod`; `go version` | `3:go 1.26.6`; `go1.26.6 darwin/arm64` |
| V2 | `TryLock` exists and observes a held `RLock` | `go doc sync.RWMutex.TryLock`; a `/tmp` probe (deleted) that takes `ForkLock.RLock`, then `TryLock`, then `RUnlock`, then `TryLock`. Run plain and `-race` | `func (rw *RWMutex) TryLock() bool` (Go ≥1.18). `tryLockFailsWhileRLocked=true tryLockSucceedsAfterRUnlock=true`, the same with and without `-race` |
| V3 | The allowlist is by name and exempts whole bodies | `cat -n host/verifygate/forklocked_write_test.go` | l.197–201 `rawWriteAllowlist` = {`forkLockedWrite`, `createTempForkLocked`, `TestKernelRefusesExecOfWriterOpenFile`}. l.284 checks `rawWriteAllowlist[encl]` with `encl` = FuncDecl name. l.291 `wrapperFuncs[fo.Name()]` is also by bare name (finding 6) |
| V4 | Test 1 covers one wrapper, 0o644, hold before Close | same | l.123 `forkLockedWrite(path, …, 0o644, nil, func(){…})`. In `forkLockedDo`: l.38 `defer RUnlock`, l.49–51 `hold()`, l.52 `return f.Close()`. Nothing observes the lock *after* `hold` |
| V5 | The scan looks at call `Fun` only | same | l.264–296 `ast.Inspect` matches only `*ast.CallExpr`. A banned func used as a **value** is never visited. The banned set is `os`.{WriteFile, OpenFile, Create, CreateTemp} (l.183) only |
| V5a | Raw-fd routes that exist | `GOOS=darwin\|linux go doc syscall.{Open,Openat,Creat}`; `go doc os.NewFile`; `go doc os.Root \| grep 'func (r \*Root)'`; `go doc io/ioutil \| grep func` | `syscall.Open` exists on both. `Openat` and `Creat` are **linux only** ("no symbol" on darwin). `os.NewFile` exists. `*os.Root` has `Create`, `OpenFile`, `WriteFile` (the method objects have `Pkg()=="os"`, so the name set already matches them). `io/ioutil` has `WriteFile`, `TempFile` (they wrap `os`, so they evade an `os`-only scan) |
| V5b | Today's imports | `grep -rln '"golang.org/x/sys/unix"\|"io/ioutil"' --include='*_test.go' .`; `grep -ln '"syscall"' host/verifygate/*.go` | no hits. `syscall` is imported by `mission_config_gate_test.go` and `forklocked_write_test.go` only. `x/sys` is `// indirect` in go.mod |
| V5c | Zero live hits under the widened ban (quorum r1, glm/kimi). Each zero is paired with a known positive for the same instrument | (i) `grep -n 'syscall\.' host/verifygate/*.go \| grep -v forklocked_write_test.go`. (ii) `grep -rnE 'os\.NewFile\|syscall\.(Open\|Openat\|Creat)\b\|os\.Root\|OpenRoot\|ioutil' host/verifygate/`; control: the same pattern `grep -cE` on this doc. (iii) value-use probe `P='[=(,] *os\.(WriteFile\|OpenFile\|Create\|CreateTemp)([^a-zA-Z(]\|$)'`, `grep -rnE "$P" host/verifygate/`; controls: `printf 'wf := os.WriteFile\n' \| grep -cE "$P"`, and `grep -rnE "$P" host cmd` | (i) only `mission_config_gate_test.go:43 &syscall.SysProcAttr{Setpgid: true}` and `:55 syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)`, neither a banned name; control: the same grep on `forklocked_write_test.go` hits l.37–38 `syscall.ForkLock.RLock/RUnlock`. (ii) no output, rc=1; control: `10`. (iii) no output, rc=1; controls: `1`, and `host/boundary/allowlist_world_test.go:380:var rawWrite = os.WriteFile` (real code, outside the scan). **Instrument note:** the probe as first written (`…[^a-zA-Z(]` with no `\|$`) returned `0` on the `printf` control, because an end-of-line value use has no trailing char. The zero was only accepted after the `\|$` fix made the control fire |
| V5d | Wrapper signatures that M1 drives (quorum r1, kimi) | `grep -nE '^func (writeFileForkLocked\|copyFileForkLocked\|createForkLocked\|createTempForkLocked\|forkLockedWrite)\(' host/verifygate/forklocked_write_test.go` | l.55 `forkLockedWrite(path string, flag int, mode os.FileMode, fill func(*os.File) error, hold func()) error`; l.60 `writeFileForkLocked(path string, data []byte, mode os.FileMode) error`; l.69 `copyFileForkLocked(src *os.File, dst, rel string, mode os.FileMode) error`; l.91 `createTempForkLocked(dir, pattern string, data []byte) (string, error)`; l.113 `createForkLocked(path string) error`. `forkLockedDo` is l.36 `(open func() (*os.File, error), fill func(*os.File) error, hold func()) error` |
| V5e | CI runners (quorum r1, gemini) | `grep -n 'runs-on' .github/workflows/*.yml` | `20: runs-on: ubuntu-latest`, `118: runs-on: ubuntu-latest`. There is no macOS CI job |
| V6 | The CI step passes on SKIP | `grep -n -A4 'Fork-locked write tests' .github/workflows/ci.yml` | l.200–204: a bare `go test … -v -run '^(…)$'` with no assertion on the output. A SKIP exits 0 |
| V6a | Pipe hazard in this file | `sed -n 93,100p .github/workflows/ci.yml` | l.93 "NO PIPE. `cmd … \| grep -q` under `set -o pipefail` is a RACE" (SIGPIPE 141, measured 3/40). Other steps use `set -euo pipefail` |
| V6b | Dry run of M3's step body at `2d3a255` (before M1–M3 exist), extracted from this doc | `RUNNER_TEMP=$(mktemp -d) bash body.sh`; then with `RUNNER_OS=Linux` | unset: rc=1 `::error::this step must run on linux`. `Linux`: `go test` `ok`, then rc=1 `::error::TestForkLockedWrappersHoldLockThroughClose did not PASS on linux (SKIP?)`: a listed test that never ran also fails the loop. After M1–M3 land, darwin fails on the kernel control instead (M3) |
| V7 | Baseline of the fast tests | `go test ./host/verifygate/ -run 'TestForkLocked\|TestVerifygateTestWritesAreForkLocked\|TestKernelRefuses' -count=1 -v` | `PASS TestForkLockedWriteBlocksConcurrentFork (0.26s)`, `SKIP TestKernelRefusesExecOfWriterOpenFile`, `PASS TestVerifygateTestWritesAreForkLocked (1.96s)`, `ok 2.492s` |
| V8 | A closed file is detectable | `/tmp` probe (deleted): `CreateTemp`, `Close`, `f.Stat()` | `file already closed isErrClosed=true` (`errors.Is(err, os.ErrClosed)`) |
| V9 | R2: packages that write and fork, and their parallelism | `grep -rln 't\.Parallel()' --include='*.go' host cmd`. Then, per dir over all `*.go`, count `os\.(WriteFile\|OpenFile\|Create\|CreateTemp)\(` × `exec\.Command(Context)?\(\|os\.StartProcess\(\|syscall\.(ForkExec\|StartProcess)\(` × `t\.Parallel()` | `t.Parallel` appears **only** in `host/verifygate/{mission_config,module_manifest}_gate_test.go`. 12 write+fork dirs: `cmd/ailang-worldd` w22 e4, `cmd/world-publish` w11 e1, `host/archive` w8 e3, `host/broker` w25 e8, `host/capsule` w10 e3, `host/coordinator` w1 e3, `host/daemon` w9 e2, `host/pkgproj` w5 e7, `host/replay` w5 e1, `host/runbook` w1 e2, `host/store` w3 e2, `host/verifygate` (par 10). **Parallel outside verifygate: 0** |
| V9a | Cost of a parse-only sweep | `find host cmd -name '*.go' \| wc -l`; time of `gofmt -l host cmd`; dir counts | 272 files; 0.07 s; 26 Go dirs; 178 `_test.go` |

**Premise corrections.**
- (a) The allowlist does not exempt `forkLockedDo`, because `forkLockedDo` holds no banned call.
  The banned calls sit in the open closures inside `forkLockedWrite` and `createTempForkLocked`
  (V3). That is why narrowing the exemption to those closures is a clean cut.
- (b) Widening the ForkLock discipline to the 11 serial packages means more than 100 write sites
  (V9). That is well over 0.5 d, so R2 gets a tripwire, not a widening (decision D3).

## Quorum r1 dispositions

Round 1 blocked 3/3 ([artifact](../verification/world-iter229/w-verifygate-forklock-guard-residuals-2026-10-03T19-28-41Z.json)).
Each premise was measured before revising.

- **gemini-3-1-pro — "the CI grep breaks the macOS workflow": REFUTED** (V5e: both jobs are
  `ubuntu-latest`). We do **not** add a `RUNNER_OS` skip guard: it would make any future non-Linux
  leg silently skip the very check, which is the SKIP hole this row closes. The step instead gains a
  fail-loud linux precondition (M3).
- **oc-glm-5-3 / oc-kimi-k3 — "zero live violations is asserted, not measured": APPLIED** as V5c,
  with a known positive per zero, and a contingency clause in D2.
- **oc-kimi-k3 — wrapper signatures unmeasured: APPLIED** as V5d; M1's drive list matches it.
- **oc-kimi-k3 — N-OpenClosure could read as a local shadow: APPLIED** (fixture reworded, M2).
- **oc-glm-5-3 (secondary) — the CI step checks only the kernel control: APPLIED.** It now
  asserts a top-level `--- PASS:` line for every test in its `-run` list (M3, AC8, AC11).

## Decisions

**D1 — E3 is observed by an after-Close probe that calls `TryLock`, inside the writer's goroutine.**
This removes the race from the check.

`forkLockedDo` gains a test-only seam:

```go
var forkLockedProbe atomic.Pointer[func(stage string, f *os.File)]
```

- `forkLockedDo` calls the probe with stage `"filled"` after `fill`/`hold` and before `Close`.
- It calls the probe with stage `"closed"` **after `f.Close()` returns**.
- The `RUnlock` stays deferred.
- The probe runs on the goroutine that holds `RLock`, so no cross-goroutine timing is involved.
  `!ForkLock.TryLock()` is a deterministic "some reader holds it" (V2). If `TryLock` succeeds, the
  probe `Unlock`s at once and records the failure.

E3 (`RUnlock` moved before `Close`) makes the `"closed"` stage's `TryLock` succeed. So the
check "the lock is held at and after Close" is exact.

The probe also checks that the file was really closed: `f.Stat()` must satisfy `errors.Is(_, os.ErrClosed)` (V8). This
proves that the stage order is real. At `"filled"` the probe checks that `f.Stat().Size()` equals
the expected length, so the data was written inside the lock.

Concurrency safety rests on four points:
- `atomic.Pointer` means no data race under `-race`.
- Only a **serial** test sets the probe. Top-level parallel tests are paused while serial tests
  run (parent V8).
- The probe ignores any file outside its own `t.TempDir()`.
- The probe **must not fork**. A fork would deadlock on `ForkLock.Lock` against the caller's own
  `RLock`. This is stated in the doc comment.

Rejected alternatives:
- A concurrent fork racing `Close`: the result depends on timing.
- Asserting from `hold` only: `hold` precedes `Close` (V4), which is the gap E3 uses.

**D2 — The scan becomes Uses-based and identity-matched, with the exemption narrowed to
`forkLockedDo`'s open closure.** It iterates over **every** `info.Uses` entry, not only call `Fun`s.

A use is banned if its object is either of these:
- `*types.Func` with `Pkg().Path()=="os"` and a name in {WriteFile, OpenFile, Create, CreateTemp, **NewFile**}.
  This name set also covers the `*os.Root` methods (V5a).
- `*types.Func` with `Pkg().Path()=="syscall"` and a name in {**Open, Openat, Creat**}.

A use is exempt **only** if its position lies in one of two places:
- (a) a `*ast.FuncLit` that is `Args[0]` of a call whose `Fun` resolves to **the package-scope
  object** `pkg.Scope().Lookup("forkLockedDo")`;
- (b) the body of the FuncDecl whose `info.Defs` object is the package-scope
  `TestKernelRefusesExecOfWriterOpenFile`.

The scan also changes in these ways:
- The by-name `rawWriteAllowlist` is deleted (this fixes finding 6).
- Every `forkLockedDo(func() (*os.File, error) {…}, …)` call is genuinely locked, so the exemption
  is semantic. Nothing is exempt because of the name of the function it sits in.
- `wrapperFuncs` counting compares `info.Uses[id]` with the package-scope wrapper objects. It
  fails with a scan error if any wrapper is missing or is not declared in `forklocked_write_test.go`.
- The import bans gain `"io/ioutil"` and `"golang.org/x/sys/unix"` (zero today, V5b), next to the
  existing aliased/dot `"os"` ban.
- Violations are sorted, because map iteration is unordered.
- **Contingency** (glm): V5c measures zero live hits today. If M2's first live run nevertheless
  surfaces a non-exempt hit, that use is exempted by package-scope object identity, like the kernel
  control, or the syscall-name ban is split to a follow-up row. The scan must not red on
  pre-existing legitimate code, and the executor must report the hit and the choice made.
- The `where` text is the enclosing FuncDecl name, or `package-level initializer`.
- `Openat`/`Creat` resolve only when the scan runs on linux (V5a). A linux-tagged file that uses
  them already fails the darwin type-check loudly, so the name set needs no build-tag logic.

With the narrowed exemption, E1 and E2 are also killed **statically**: a raw `os.CreateTemp`/`os.OpenFile`
in a wrapper body is outside the open closure. D1 then covers the bypasses the scan cannot see,
such as E3 and a skipped lock in `forkLockedDo`.

**D3 — R2 gets a tripwire, not a widening.** A new `TestNoParallelWriteForkPackagesOutsideVerifygate`
runs a parse-only sweep of every `*.go` under `../../host` and `../../cmd` (0.07 s, V9a), grouped by
directory. A directory is a violation if all of the following hold:
- it is not `host/verifygate`;
- a `_test.go` calls `X.Parallel()` with no arguments;
- any `.go` (test or production, because a test may write a script that production code execs)
  has a write selector (`os.{WriteFile,OpenFile,Create,CreateTemp,NewFile}`, `syscall.{Open,Openat,Creat}`);
- any `.go` has a fork selector (`exec.{Command,CommandContext}`, `os.StartProcess`,
  `syscall.{ForkExec,StartProcess,Exec}`).

The sweep follows the scan's own pattern:
- **Floors:** at least 26 dirs and at least 12 write+fork dirs (V9), so an empty sweep cannot pass.
- **A known positive.**
- **A known negative.**

It is syntactic, matching on the selector's `X` ident name, so an aliased import evades it. That is
recorded as R2′. A tripwire only has to fire on the ordinary way someone adds parallelism.

## Milestones

**M1 — Wrapper bodies guarded (kills E1, E2, E3).** File: `host/verifygate/forklocked_write_test.go`.
- Add `forkLockedProbe` (D1) and its two call points in `forkLockedDo`. Imports gain `sync/atomic`.
- New serial test `TestForkLockedWrappersHoldLockThroughClose`. It installs the probe
  (with `t.Cleanup` restoring `nil`) and drives each wrapper under its own `t.TempDir()`, with the exact signatures of V5d:
  - `writeFileForkLocked(p, data, 0o755)`
  - `copyFileForkLocked(src, dst, "rel", 0o755)`, with `src` from `os.Open` (read-only, not banned)
  - `createTempForkLocked(dir, "p*", data)`
  - `createForkLocked(p)`
  - `forkLockedWrite(p, O_WRONLY|O_CREATE|O_TRUNC, 0o755, fill, nil)`
- For each call it asserts that the probe stages are exactly `[filled closed]` and that the content
  and mode match.
- Failure texts:
  - `"<wrapper>: fork-lock probe stages = %v, want [filled closed]"` (bypass: E1, E2)
  - `"<wrapper>: ForkLock not read-held during fill"`
  - `"<wrapper>: ForkLock released before Close returned"` (E3)
  - `"<wrapper>: probe at closed stage saw an open file"`
- Test 1 switches to 0o755 (the mode that matters), with no other change.
- Gate: `go vet ./host/verifygate/` and
  `go test ./host/verifygate/ -run '^TestForkLocked' -count=1 -v`, plus the same with `-race`.

**M2 — Scan hardening (kills E5b, E8b, E6b, finding 6).** Same file. Rework
`scanForkLockedWrites` per D2: `types.Info` gains `Defs`, the scan iterates over Uses, the
exemption is the open closure plus the kernel-control identity, the banned set and import bans
are extended, and the output is sorted. The existing 5 positives and the negative stay. New
fixtures:

| Fixture | Source shape | Must report |
|---|---|---|
| P-FuncValue (E5b) | `func h() { wf := os.WriteFile; _ = wf("x", nil, 0o755) }` | the `os.WriteFile` line, `func h` |
| P-RawFd (E8b) | `func h() { fd, _ := syscall.Open("x", syscall.O_CREAT\|syscall.O_WRONLY, 0o755); f := os.NewFile(uintptr(fd), "x"); f.Close() }` | the `syscall.Open` line **and** the `os.NewFile` line |
| P-PkgInit (E6b) | `var _ = os.WriteFile("x", nil, 0o644)` at package level | that line, `package-level initializer` |
| P-NamedWrapper (F6) | `func forkLockedWrite() { _, _ = os.OpenFile("x", 0, 0) }` (the old by-name allowlist shape) | that line |
| P-Ioutil | `import "io/ioutil"` + `ioutil.WriteFile(…)` | the ImportSpec line |
| P-RootWrite | `func h(r *os.Root) { _ = r.WriteFile("x", nil, 0o755) }` | that line |
| N-OpenClosure | a package-level `forkLockedDo` defined in the fixture package, called with an open closure `func() (*os.File, error) { return os.OpenFile("x", 0, 0) }` | **nothing** (the exemption is live; the callee is the fixture's package-scope object, so identity matching applies) |

Live floors stay (at least 9 files, at least 21 wrapper sites) and are now counted by identity.
The live scan must stay at zero violations, which proves that the narrowed exemption covers both
real open closures (`forkLockedWrite` l.56, `createTempForkLocked` l.93) and the kernel
control. Gate: `go vet`, plus `go test ./host/verifygate/ -run '^TestVerifygateTestWritesAreForkLocked$' -count=1 -v`.

**M3 — CI SKIP fail and R2 tripwire.**
- Add `TestNoParallelWriteForkPackagesOutsideVerifygate` (D3) to the same file, with:
  - floors of 26 dirs and 12 write+fork dirs;
  - known positive: a synthetic dir `{"x_test.go": t.Parallel + os.WriteFile + exec.Command}` must be named;
  - known negative: the same dir without `t.Parallel` must not be named.
- Replace the `ci.yml` step body at l.200–204, and add the new test to its `-run` list. Use no pipe (V6a):

```yaml
        run: |
          set -euo pipefail
          # Assumes linux: fail loud, never skip (a skip guard reopens the SKIP hole).
          [ "${RUNNER_OS:-}" = Linux ] || { echo "::error::this step must run on linux"; exit 1; }
          tests="TestForkLockedWriteBlocksConcurrentFork TestForkLockedWrappersHoldLockThroughClose TestKernelRefusesExecOfWriterOpenFile TestVerifygateTestWritesAreForkLocked TestNoParallelWriteForkPackagesOutsideVerifygate"
          log="$RUNNER_TEMP/forklock-verbose.log"
          rc=0
          go test ./host/verifygate/ -count=1 -p 1 -timeout 120s -v \
            -run "^(${tests// /|})\$" >"$log" 2>&1 || rc=$?
          cat "$log"
          [ "$rc" -eq 0 ] || exit "$rc"
          for t in $tests; do
            grep -q -- "^--- PASS: $t " "$log" \
              || { echo "::error::$t did not PASS on linux (SKIP?)"; exit 1; }
          done
```

- The `-run` regex is built from the same `$tests` list the loop checks, so the two cannot drift.
  `grep -q` reads a file, and the step has no pipe at all (V6a).
- Local known positives with no edit needed, both recorded with their rc in the PR: on the rig,
  the body with `RUNNER_TEMP=$(mktemp -d)` and `RUNNER_OS` unset **must exit 1** on the linux
  precondition; with `RUNNER_OS=Linux` it **must exit 1** on
  `::error::TestKernelRefusesExecOfWriterOpenFile did not PASS on linux (SKIP?)`, because darwin
  SKIPs that test (V7).
- Gate: `verify_go.sh` and `verify_ail.sh` with `AILANG_BIN` set, plus `check_no_personal_email.sh`.

## Acceptance (load-bearing, each mutant → firing assertion; coding-standards S6)

| AC | Mutant (scratch worktree, reverted) | Assertion that fires |
|---|---|---|
| AC1 (E1) | `createTempForkLocked` does `os.CreateTemp` + `Write` + `Close` directly, with no `forkLockedDo` | scan: `raw os.CreateTemp in func createTempForkLocked`; **and**, with the scan disabled, `createTempForkLocked: fork-lock probe stages = [], want [filled closed]` |
| AC2 (E2) | `forkLockedWrite`: `if mode&0o111 != 0 {` raw `os.OpenFile` / fill / Close `}` | scan: `raw os.OpenFile in func forkLockedWrite`; probe: `writeFileForkLocked: … stages = []` |
| AC3 (E3) | `forkLockedDo`: replace `defer RUnlock` with `RUnlock()` immediately before `f.Close()` | `ForkLock released before Close returned` (test 1 and the scan stay green, which shows that only D1 sees E3) |
| AC4 (E5b) | add `func h() { wf := os.WriteFile; _ = wf(p, nil, 0o755) }` to `ail_binary_gate_test.go` | live scan names that file:line, `func h` |
| AC5 (E8b) | add the P-RawFd body as a helper in `toolchain_pin_gate_test.go` | live scan names the `syscall.Open` and `os.NewFile` lines |
| AC6 (E6b) | the scan skips uses with no enclosing FuncDecl (`if encl == "" { continue }`) | `fixture PkgInit: … not reported` |
| AC7 (F6) | restore the by-name exemption for FuncDecls named `forkLockedWrite` | `fixture NamedWrapper: … not reported` |
| AC8 (CI SKIP) | the kernel test skips unconditionally (`if true { t.Skip(…) }`) on a scratch PR commit | the step fails with `::error::TestKernelRefusesExecOfWriterOpenFile did not PASS on linux (SKIP?)`. The loop gives the same failure for a SKIP of **any** listed test. The darwin runs of the same body are the free local known positives (M3) |
| AC9 (R2) | add `t.Parallel()` to one test in `host/broker` | `host/broker/<file>:<line>: t.Parallel in a write+fork package outside the fork-locked scan` |
| AC10 (exemption live) | delete the open-closure exemption | the live scan reds on `forkLockedWrite`/`createTempForkLocked` and N-OpenClosure is reported |
| AC11 (regression) | n/a | all prior ACs of the parent stay green. The fast tests run in under 5 s locally (baseline 2.5 s, V7). CI shows a top-level `--- PASS: ` line for each of the five tests in the step, asserted by the step's loop |

## Risks and residuals

- **R1. Deadlock hazard.** Nothing inside a `ForkLock.RLock` region may fork: not `fill`, `hold`,
  the probe, or the open closure. A fork there calls `ForkLock.Lock` against the goroutine's own
  `RLock`, so it hangs until `-timeout`. Test 1's 5 s bound catches a lost wake-up, but not a
  self-deadlock inside a probe. The probe therefore does only `TryLock`/`Unlock`/`Stat`.
- **R2′. The tripwire is syntactic.** An aliased `exec`/`os` import, or a fork reached through
  another package, evades it. Widening the full ForkLock scan to the 11 packages (V9) stays a
  follow-up row, filed only if one of them needs parallelism.
- **R3. `TryLock` false-negative direction.** A stray goroutine from an earlier test holding
  `RLock` would make the after-Close `TryLock` fail even under E3, which hides the mutant. The
  correct direction cannot false-fail: a reader is always held at both stages. Serial scheduling
  (parent V8) makes a stray holder unlikely.
- **R4. Linux-only names.** `Openat` and `Creat` are matched only when the scan runs on linux (V5a).
  The rig cannot exercise P-RawFd's linux variants, so P-RawFd uses `syscall.Open` (both OSes).
- **R5. Judge finding 5 is unchanged.** A kernel that drops write-deny reds CI via the kernel
  control. That is correct and loud.

## Conflict surface

- `host/verifygate/forklocked_write_test.go`: every change except the CI step.
- `.github/workflows/ci.yml`: the one step at l.200–204.

Re-fetch `origin/dev` before merging: attended row-135 or row-138 follow-ups may touch `ci.yml`.

## Non-goals

- Widening the ForkLock discipline outside `host/verifygate` (R2′; V9 lists 11 packages).
- Any change to production code, gate scripts, or the 21 migrated write sites.
- Retry-on-ETXTBSY (parent option A stays rejected).
- Kernel-behaviour drift detection (finding 5).
