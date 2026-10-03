# w-verifygate-etxtbsy

**Status**: DRAFT (iteration 228 designer), small defect fix, no feature. Base `origin/dev` `c2476ff`.
**Scope**: `host/verifygate` test code only (a test-only package: 8 `.go` files, all `_test.go`, V2).
No production code, no `.ail`, no `tools/launchd/*`, no gate script semantics change.
**Chosen option**: **C′**: every test-side file write in `host/verifygate` goes through one helper
that holds `syscall.ForkLock.RLock()` from open to close, so no fork in the test process can
inherit a write fd. A total-ban AST scan enforces it, with known-positive fixtures and mutations.

## Quorum log

**Round 1** was BLOCKED: all 3 present reviewers rejected, and `gpt6-1-sol` was unreachable.
The controller measured each objection, and revision 1 applies the results:
- **Refuted objection (gemini), "forks take ForkLock.RLock, so the polarity is reversed".**
  This is false. Forks take `ForkLock.Lock()` (V9a), so a writer holding `RLock` excludes them.
- **Refuted objection (kimi), attribution of :564.** :564 is the `runGateAt` caller line (V1a);
  the failing exec is :150.
- **Accepted (glm), liveness.** The starvation sentence is rewritten with V9b and V9, and test 1
  gains a 5 s bound after release (AC1b).
- **Accepted (glm), helper imports.** V2a records that no imported helper writes files for the tests.
- **Housekeeping.** V14 is relabelled as the pre-fix baseline, and V13a substantiates option D's rejection.

**Round 2** was BLOCKED. Under the narrow-refinement carve-out, the reviewers' own fixes are
applied where they hold, and the rest are refuted by measurement:
- R2 objection (gemini), "forkpipe2.go has no `hasWaitingReaders` branch" → **refuted
  (measurement)**: it is at l.26/28 (declaration) and l.52 (`if hasWaitingReaders(&ForkLock)`); see V9.
- R2 objections (gemini) → **applied fixes**: the finite-bound liveness sentence (10 arms per V5); the scan enumerates every `*.go`, not only `*_test.go`.
- R2 objections (glm) → **applied fixes**: `O_EXCL` and "keeps its error text" are substantiated by V15 (glm's option b); "no `-v`" by V16; the `hasWaitingReaders` relief is scoped to forkpipe2 (Linux), and the darwin wait is bounded by the finite number of in-process forks.
- R2 objection (kimi), a syntactic selector match is evadable → **applied fix**: kimi's
  `go/types` spec, the aliased/dot import ban, the no-FuncDecl rule, the 5th fixture and AC8.
  The importer was probed (V17).
- R2 objection (kimi), the forkpipe2 path is only argued → **applied fix**: V18 (pending, filled at PR time).

## Problem

`c2476ff` CI run 37132037850 **attempt 1** went red on job "go host build + test gate"
(job 111228885768) with
`--- FAIL: TestModuleManifestEmptyEnumerationFailsLoudly (30.11s)` and
`module_manifest_gate_test.go:564: start isolated verify gate: fork/exec …/001/iso/scripts/verify_ail.sh: text file busy`
(V1). Attempt 2 of the same run was re-run green at 15:35Z (V3), so this is an **intermittent**
defect on dev, not a deterministic red. The previous five dev runs were green (V3).

Mechanism (golang/go#22315, quoted in the Go toolchain's own source, V6). Goroutine A writes the
script (`os.WriteFile`/`os.OpenFile`, both `O_CLOEXEC`) and closes it. Goroutine B, a different
parallel test, calls `exec.Cmd.Start` while A's write fd is open. The forked child gets a copy of
the whole fd table, so it inherits A's write fd. `O_CLOEXEC` closes that fd only when the child
reaches its own `execve`. Until then the inode has a writer, so A's later `execve` of that file
returns `ETXTBSY` on Linux. macOS does not enforce this (V7: exec of a file that still has a
writer open succeeds on darwin), so the rig cannot reproduce the failure.

It became reachable in `e6064b6` (PR #192). That commit added `t.Parallel()` to 9 arms in
`module_manifest_gate_test.go` and 1 in `mission_config_gate_test.go` (V4, V5). Top-level parallel tests run
"in parallel with (and only with) other parallel tests" (V8), so before that commit no two
`verifygate` tests could overlap and no fork could race a write.

**Premises corrected:** (a) the failure was in the **plain** leg `go test ./... -count=1`, not
under `-race`. The log has no `-race -timeout` leg line because the plain leg failed first (V1).
(b) The job is not stuck red: attempt 2 is green (V3). (c) The count is 10 `t.Parallel()` calls
in two files, not ~10 in one.

## Verification Log

All commands were run from `/Users/voightkampff/dev/sunholo-data/.wt-world-iter228` at `c2476ff`,
with go1.26.6 darwin/arm64, unless noted otherwise.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Failure text and leg | `grep -n "── go test\|FAIL\|text file busy" ~/.ailang/state/w228_ci_red.log`; `grep -n -- "-race -timeout" …` | l.377 `── go test ./... -count=1`; l.401 `--- FAIL: TestModuleManifestEmptyEnumerationFailsLoudly (30.11s)`; l.403 `fork/exec /tmp/TestModuleManifestEmptyEnumerationFailsLoudly1783032982/001/iso/scripts/verify_ail.sh: text file busy`; l.405 `FAIL …/host/verifygate 191.410s`; the `-race -timeout` grep returns nothing (that leg never ran) |
| V1a | Attribution of :564 and the 30.11 s | `sed -n 555,566p host/verifygate/module_manifest_gate_test.go` | l.558 `t.Parallel()`; l.560 `requirePristineControl(t, root)`; l.561 `mutateCopiedScript(...)`; l.564 `rc, out := runGateAt(t, root, …)`. `runGateAt` (l.140) calls `t.Helper()` (l.141), so its `t.Fatalf` reports the caller line 564. The failing `execve` is therefore `runGateAtErr:150` on the `verify_ail.sh` that l.561 rewrote immediately before. Most of the 30.11 s is `requirePristineControl` (one shared full gate run, before the mutation). |
| V2 | verifygate is test-only | `ls host/verifygate/*.go \| wc -l`; `ls *.go \| grep -v _test` | `8`; empty. Control: the same `ls` lists 8 `_test.go` files |
| V2a | No imported helper writes on the tests' behalf | `go list -f '{{.GoFiles}} {{.TestImports}}' ./host/verifygate/`; per-file `import (…)` blocks | `GoFiles=[]` (no same-package non-test helpers). TestImports are **stdlib only**: bytes context crypto/sha256 encoding/json errors fmt go/ast go/build/constraint go/parser go/token go/types go/version io io/fs os os/exec path/filepath regexp runtime slices sort strconv strings sync syscall testing time. No repo package is imported. `io` appears only as `io.Copy` into an fd opened by `os.OpenFile` (`copyGateFileErr`), which the ban already covers; `t.TempDir`/`os.MkdirTemp` create directories and leave no fd open |
| V3 | Run history | `gh run list --branch dev --workflow ci.yml -L 6`; `gh api …/runs/37132037850/attempts/1/jobs` | attempt 1: `go host build + test gate failure 111228885768`; run now `attempt:2`, both jobs success at `2026-10-03T15:35:31Z`; e6064b6, 78db602, 901b87c, 0b18b94, fe835ed `success` |
| V4 | Parallel arms are only in verifygate | `grep -rln "t\.Parallel()" --include='*_test.go' .` | `host/verifygate/mission_config_gate_test.go`, `host/verifygate/module_manifest_gate_test.go` (the hits are the known positive) |
| V5 | Where `t.Parallel` is called | `grep -n "t.Parallel()" module_manifest_gate_test.go`; `grep -n -B8 t.Parallel mission_config_gate_test.go` | module_manifest l.264, 341, 387, 420, 456, 467 (subtest), 498, 517, 558; mission_config l.246 (`TestMissionConfigBlockingTimeoutClassified`) |
| V6 | Toolchain precedent | `grep -rn ETXTBSY $(go env GOROOT)/src` | `cmd/go/internal/test/test.go:1648` "spurious ETXTBSY errors on Unix platforms (see https://go.dev/issue/22315)" with an unbounded retry around `cmd.Run`; `base/base.go:218` `RunStdin` with a bounded retry (3 tries, `100ms<<try`); `cmd/internal/script/cmds.go:457` retries `cmd.Start` |
| V7 | macOS is benign | `/tmp/w228probe` probe that execs a script while its `O_WRONLY` fd is still open | `GOOS=darwin writer-open exec: out="ran\n" err=<nil> isETXTBSY=false`. **Linux could not be measured here**: `docker/podman/colima/limactl/orb not found` |
| V8 | Parallel semantics | `go doc testing.T.Parallel` | "run in parallel with (and only with) other parallel tests" |
| V9 | Every fork takes ForkLock for writing | `grep -n ForkLock $GOROOT/src/syscall/exec_unix.go forkpipe2.go`; `go doc syscall.ForkLock` | `exec_unix.go:199 acquireForkLock()` in `forkExec`; `forkpipe2.go:31` "ForkLock is exported and we've promised that during a fork we will call ForkLock.Lock"; `grep -n hasWaitingReaders $(go env GOROOT)/src/syscall/forkpipe2.go` → `26:// hasWaitingReaders reports whether any goroutine is waiting`, `28:func hasWaitingReaders(rw *sync.RWMutex) bool`, `52: if hasWaitingReaders(&ForkLock) {`. So on Linux, when readers are waiting, that branch lets them through between overlapping forks |
| V9a | Forks take the WRITE lock (polarity) | `grep -n "ForkLock.Lock()\|standard library should never" $GOROOT/src/syscall/forkpipe{,2}.go`; `//go:build` lines | `forkpipe2.go:45 ForkLock.Lock()` (build: dragonfly/freebsd/**linux**/netbsd/openbsd/solaris); overlapping forks share that one write lock via the `forking` counter, and l.75 re-takes `ForkLock.Lock()`; l.62–63 says "the standard library should never take a read lock on ForkLock"; `forkpipe.go:25 ForkLock.Lock()` (build: aix/**darwin**). **So a reader holding `RLock` excludes every fork on both CI and the rig.** |
| V9b | RWMutex writer preference | `go doc sync.RWMutex` | "If any goroutine calls RWMutex.Lock while the lock is already held by one or more readers, concurrent calls to RWMutex.RLock will block until the writer has acquired (and released) the lock" |
| V10 | RLock blocks a concurrent fork | `/tmp/w228probe` prototype: hold `ForkLock.RLock`, `exec.Command("/bin/sh","-c","exit 0").Start()` from a goroutine, check after 250 ms; run with and without `-race` | `holdLock=true startReturnedWhileWriteOpen=false startErr=<nil>`; `holdLock=false startReturnedWhileWriteOpen=true` (both plain and `-race`) |
| V11 | Write call sites in verifygate tests | per-file `grep -c "os\.WriteFile(" / "os\.OpenFile(" / "os\.Create(" / "os\.CreateTemp("` | ail_binary 3/0/0/2, evidence_manifest 4/0/0/0, mission_config 1/0/0/0, module_manifest 7/1/0/0, toolchain_pin 2/0/1/0, the other 3 files 0. **21 sites** |
| V12 | Exec call sites | `grep -n "exec.Command" host/verifygate/*_test.go` | real: module_manifest:150, evidence_manifest:142,160, ail_binary:58,66,296,479,584, toolchain_pin:493,593,615, mission_config:40; subprocess_sink:234–302 are fixture **string literals** (so a grep can't act as the guard) |
| V13 | Budget numbers | `git show e6064b6` | verifygate `go test -v` 238 s → 104 s; GOMAXPROCS=2 race leg 465 s → 312 s; the cap is `p.wait(timeout=600)` (`verify_go.sh:328`); l.45 "passing test events before, 148 after: +1 test, +3 subtests" |
| V13a | Why D is not viable | `git show -s 78db602` | l.9–11 "a budget overrun, not a test failure: verify_go.sh's -race leg exceeded 600s on the 2-core runner (every package passed; host/verifygate alone took 423s in the plain leg after #190 …)" |
| V14 | CI baseline, **UNFIXED dev** (attempt-2 re-run of c2476ff, job 111231727567) | `gh run view --job 111231727567 --log` | verifygate plain `179.329s`, race `215.062s`, `✓ go gate PASSED` |
| V15 | `copyGateFileErr` already uses `O_EXCL`; every dst is fresh | `sed -n '31,51p' host/verifygate/module_manifest_gate_test.go`; `grep -n "copyGateFile\(Err\)\?(" host/verifygate/*_test.go`; `sed -n 50,80p evidence_manifest_gate_test.go` | l.42 `out, err := os.OpenFile(dst, os.O_CREATE\|os.O_EXCL\|os.O_WRONLY, mode)`; error texts `open copy source %s`, `create copy target %s`, `copy %s` (l.36/44/48). Call sites: module_manifest l.65, 68, 84 (in `buildIsolatedGateRoot`, root = `t.TempDir()/iso` via l.55, or `os.MkdirTemp` for the shared control), l.26 (the `copyGateFile` wrapper); evidence_manifest l.79 (root = `filepath.Join(t.TempDir(), "iso")`, l.58). Every dst is a fresh path under a per-test root |
| V16 | verify_go.sh runs go test without `-v` | `grep -n 'go test' scripts/verify_go.sh` | l.299 `go test -json ./host/evidence -count=1`; l.312 `go test ./... -count=1`; l.321/325 race leg `go test ./... -count=1 -race -timeout 8m`. There is no `-v`; l.5 and l.10 are comments |
| V17 | `go/types` over verifygate is practical | 30-line probe (`/tmp`, deleted): parse every `host/verifygate/*.go`, `types.Config{Importer: …}.Check`, count `Uses` resolving to `os.{WriteFile,OpenFile,Create,CreateTemp}`; plus a fixture `import osw "os"; var _ = osw.WriteFile(…)` | `importer=source files=8 typeErr=<nil> bannedResolved=21 elapsed=599ms`; `importer=default … bannedResolved=21 elapsed=1.435s`. The fixture gives `bannedResolved=1` with both importers. 21 matches V11 (known positive). **Chosen: `importer.ForCompiler(fset,"source",nil)`**, which needs no `go list` subprocess. Aliased-import grep `^\s+[A-Za-z_.]+ "os"$` over `*.go` → rc=1 (none today); control: plain `"os"` imports = 8 |
| V18 | forkpipe2 path behaves as argued (**PENDING**, filled at PR time) | new ci.yml verbose step | test 1 `--- PASS` on linux/amd64 (job id + log line to be recorded) |

**Write-then-exec enumeration (V11 × V12, read line by line):**

| Site | Written by | Exec'd how | In a parallel test? |
|---|---|---|---|
| `module_manifest:150` `runGateAtErr` | `copyGateFileErr:42` (`OpenFile` 0755); `mutateCopiedScript:315` | direct `execve` | **YES**. This is the failing **exec** (V1a: l.561 mutates, l.564 → :150 execs). The racing **fork** was some other parallel test's `Start`; the log cannot tell which. C′ covers every candidate, because every fork in the process takes the write lock (V9a) |
| `evidence_manifest:142,160` (verify_go.sh) | l.75, l.79 (`copyGateFile`), l.197 | direct | no |
| `ail_binary:479`, `:584` | `CreateTemp` l.464/569 + `WriteString` + `Close` + `Chmod` | direct | no |
| `ail_binary:376` `writeExecutable` shim | `WriteFile` 0755 | indirect: verify_ail.sh execs it as `AILANG_BIN` | no |
| `toolchain_pin:484` `writeExecutableAt` | `WriteFile` 0755 | indirect: fake tools exec'd by a bash fixture | no |
| `mission_config` (verify_go.sh, mission_decisions.sh) | `writeFile:113` 0755 | `/bin/bash <script>` (interpreter, no execve of the file) | the l.246 arm, but no execve |

Outside verifygate, `host/{broker,archive,capsule,pkgproj,daemon,replay,transitionreg}` and
`cmd/world-publish` tests also write 0755/0700 scripts and exec them. None of those packages has a
parallel test (V4), so they are out of scope (Residual R2).

## Design: options

- **(A) Retry `Start` on `errors.Is(err, syscall.ETXTBSY)` with bounded backoff.** Rejected as the
  primary fix, but this is the fallback if C′ is refused. It has toolchain precedent (V6). It only
  covers a *direct* `execve` made from Go. It cannot cover the indirect sites (`ail_binary:376`,
  `toolchain_pin:484`), where bash is the process that execs a file this process wrote. Those
  are serial today, and a future `t.Parallel()` on them would bring the same flake back, out of
  reach of any Go-side retry. It also suppresses the symptom rather than removing the cause, and
  it leaves a non-zero residual when the backoff runs out. `exec.Cmd` cannot be restarted, so the
  helper would need a `func() *exec.Cmd` constructor at every exec site.
- **(B) Run copied scripts as `bash <script>`.** Rejected. The `mode-bit` arm (l.284) only
  compares digests through `pristineControlCovers` and never execs, so B would not break it.
  But B silently drops the copied script's `#!/usr/bin/env bash` and its exec bit from what the
  gate exercises. It also leaves both indirect shim sites uncovered, and the direct sites in
  `ail_binary` and `evidence_manifest` would all need the same change.
- **(C) Package-level RWMutex.** Rejected as stated. A private mutex only excludes forks that
  also take it, so every exec site *and* every write site would need it: a two-sided protocol
  with twice the places to miss.
- **(C′) `syscall.ForkLock.RLock()` around each write: CHOSEN.** The standard library already
  takes `ForkLock` for writing around every fork in the process (V9), whether the fork comes from
  `exec.Cmd`, `os.StartProcess` or anything else. So only the write side needs the protocol.
  While any fork-locked write fd is open, no fork can start. Once the fd is closed, no child can
  hold it. The cause is removed, with no timing involved, and indirect execs are covered too.
  The cost is that a fork waits out an in-flight write of a small file (microseconds).
  **Liveness** is guaranteed in both directions, and each wait is bounded by one small file
  write or one fork-to-exec:
  - **A fork cannot be starved by a chain of writes.** Once a fork's `Lock` is waiting, new
    `RLock` calls block until the fork has acquired and released the lock (V9b).
  - **A write cannot be starved by a stream of forks.** On Linux (`forkpipe2.go`, the CI path),
    `acquireForkLock`'s `hasWaitingReaders` branch lets waiting readers through between
    overlapping forks (V9). That relief is specific to forkpipe2. On the rig (darwin,
    `forkpipe.go`), `acquireForkLock` is a plain `ForkLock.Lock()` (V9a), and the wait is bounded
    by the finite number of in-process forks. On both, the stream of forks is strictly bounded by
    the finite number of concurrent parallel tests (10 arms per V5), so the write proceeds once
    the finite set of overlapping forks completes.
- **(D) Revert the `t.Parallel()` calls.** Rejected. That undoes the 238 s → 104 s and
  465 s → 312 s speed-up (V13). Before that speed-up, `78db602` had to revert a feature because
  the race leg overran its 600 s cap on the 2-core runner, with verifygate alone at 423 s in the
  plain leg (V13a).

### Mechanism (new file `host/verifygate/forklocked_write_test.go`, package `verifygate`)

- `forkLockedWrite(path string, flag int, mode os.FileMode, fill func(*os.File) error, hold func()) error`
  is the core. It takes `syscall.ForkLock.RLock()`, then `OpenFile`, `fill`, the optional `hold`
  seam (tests only; `nil` everywhere else), then `Close`, then `RUnlock` (deferred). Nothing
  inside the locked region forks: `fill` writes bytes only. A fork inside the region would
  deadlock (`acquireForkLock` → `ForkLock.Lock`), and a doc comment says so.
- There are thin wrappers for each call shape used today. `writeFileForkLocked(path, data, mode)`
  replaces `os.WriteFile`. `copyFileForkLocked(src, dst, mode)` uses `O_EXCL`, replaces the body of
  `copyGateFileErr`, and keeps its error text. It can keep `O_EXCL`, because every call site
  writes a fresh path under a per-test root (V15). `createTempForkLocked(dir, pattern, data) (string, error)`
  replaces `CreateTemp`, `WriteString` and `Close`, and the caller still does `Chmod` (`Chmod` opens
  no fd). `createForkLocked` covers the `os.Create(sentinel)` case.
- The rule is total: **every** file-creating or writing call in verifygate tests goes through a
  wrapper, including writes of non-executables such as `.ail` mutations. The scan cannot know
  which paths will be exec'd later, and the lock costs almost nothing.

### Mechanical guard: `TestVerifygateTestWritesAreForkLocked` (AST, not grep, per V12)

The scan enumerates **all** `*.go` files in `host/verifygate` with `os.ReadDir`, test and
non-test alike. It parses them with `go/parser` and type-checks them as one package with
`go/types`, using `importer.ForCompiler(fset, "source", nil)`. That works because the files import
the standard library only (V2a), and V17 measures it.

A call is a violation if its callee resolves through `types.Info.Uses` to a `*types.Func` with
`Pkg().Path() == "os"` and a `Name` in {`WriteFile`, `OpenFile`, `Create`, `CreateTemp`}. Further rules:
- The scan also fails on any aliased or dot `ImportSpec` for `"os"`.
- A banned call with no enclosing `FuncDecl`, for example a package-level `var` initializer, is
  a violation, never a skip.
- The allowlist (the `forkLocked*` wrappers and the Linux kernel control below) is matched on
  the attributed enclosing function.
- A type-check error fails the scan.

It follows the fixture pattern already used by `subprocess_sink_gate_test.go`, and has three parts:
- **Floors.** It fails if fewer than 9 files are enumerated (the 8 of V2 plus the new wrapper
  file), or if fewer than 21 wrapper call sites are found (V11), so an empty enumeration cannot pass.
- **Known positives.** Five fixture sources must each be reported with file:line: one per banned
  call, plus `import osw "os"; osw.WriteFile(...)`, which must be flagged both as an aliased
  import and as a resolved banned call.
- **Known negative.** A wrapper call and an `os.ReadFile` must not be reported.

### Proof the fix works when the rig cannot reproduce the flake

1. `TestForkLockedWriteBlocksConcurrentFork` (all OS, serial). In `hold`, it signals `inside`
   and waits on `release`. A goroutine then runs `exec.Command("/bin/sh","-c","exit 0").Start()`.
   After 250 ms the test asserts that `Start` has **not** returned. It then closes `release` and
   asserts that `Start` returns with nil error **within 5 s of the release**, so a liveness
   regression such as a deadlock or a lost wake-up fails loudly instead of hanging until `-timeout`. This is the same shape as the V10 prototype,
   measured on darwin both plain and with `-race`.
2. `TestKernelRefusesExecOfWriterOpenFile` (Linux-only real-kernel known positive; on other OSes
   it skips with the reason). It opens a 0755 script `O_WRONLY` and keeps the fd open, without
   the wrapper, so it is on the scan allowlist. It then asserts that `exec.Command(path).Start()`
   satisfies `errors.Is(err, syscall.ETXTBSY)`, and that the same exec succeeds after `Close`.
   That proves the kernel premise on the CI runner. A kernel that stops returning ETXTBSY makes
   this test fatal, not a silent pass. **Not verified locally** (no Linux runtime, V7).
3. Both tests run in a new verbose CI step,
   `go test ./host/verifygate/ -count=1 -v -run '^(TestForkLockedWriteBlocksConcurrentFork|TestKernelRefusesExecOfWriterOpenFile|TestVerifygateTestWritesAreForkLocked)$'`.
   The step exists because `verify_go.sh` runs without `-v`, so a skip and a pass look the same
   there. The precedent is the row 24 verbose step.

## Milestones

- **M1: wrappers and mechanism tests.** Add `forklocked_write_test.go` with the core,
  4 wrappers, test 1 and test 2. Gate: `go vet ./host/verifygate/` and
  `go test ./host/verifygate/ -run '^(TestForkLockedWrite|TestKernelRefuses)' -count=1 -v`
  (test 2 skips on darwin, which is expected).
- **M2: migrate all 21 write sites** (V11) to the wrappers. Assertions, arms and error texts stay
  as they are. Gate: full `AILANG_BIN=~/.pinned-ailang/ailang ./scripts/verify_go.sh` rc 0, with
  the verifygate test-event count unchanged from 148 (the `e6064b6` figure) plus the new tests.
- **M3: AST scan guard and the CI verbose step.** Add `TestVerifygateTestWritesAreForkLocked`
  with its floors and fixtures, and add the `ci.yml` step from item 3. Gate: `verify_go.sh`,
  `verify_ail.sh`, `check_no_personal_email.sh`.

## Acceptance (load-bearing, each with mutant → firing assertion; coding-standards S6)

| AC | Claim | Mutant | Assertion that fires |
|---|---|---|---|
| AC1 | The lock prevents the fork | delete the `RLock`/`RUnlock` in `forkLockedWrite` | test 1 "Start returned while a fork-locked write fd was open" |
| AC2 | No raw write is left | revert `mutateCopiedScript:315` to `os.WriteFile` | scan test names `module_manifest_gate_test.go:315` |
| AC3 | Scan covers every banned call | remove `CreateTemp` from the banned set | the CreateTemp known-positive fixture is "not reported" |
| AC4 | Scan cannot pass vacuously | point the scan at an empty dir, or drop the wrapper floor to 0 | the floor assertion fires (files < 8, or wrapper sites < 21) |
| AC5 | Kernel premise is real on CI | close the writer before exec in test 2 | "expected ETXTBSY" (**CI/Linux only**, a labelled drill on a scratch PR commit, then reverted) |
| AC6 | Linux green | n/a | "go host build + test gate" and the new verbose step are green on the PR head **and** on the merge commit, and the verbose log shows test 2 `--- PASS` (not SKIP) |
| AC1b | The lock does not deadlock | make `forkLockedWrite` skip its `RUnlock` | test 1 "Start did not return within 5 s of release" |
| AC8 | Aliased imports cannot evade the scan | rename the `os` import in one verifygate file (e.g. `osx "os"`) and use `osx.WriteFile` | scan fires naming the aliased `ImportSpec` and the call site file:line |
| AC7 | Budget holds | n/a | race leg finishes under 600 s; verifygate race time within ±15% of the **pre-fix** 215 s baseline (V14, measured on unfixed dev) |

A green CI run alone is **not** proof of the fix: the flake was 1 red in the 6 observed runs
(V3), and a single green would also be likely without any fix. AC1–AC5 carry the proof. AC6 and
AC7 are regression fences.

## Files / Conflict Surface

- `host/verifygate/forklocked_write_test.go` (new)
- `host/verifygate/{module_manifest,evidence_manifest,ail_binary,mission_config,toolchain_pin}_gate_test.go`
  (write sites only)
- `.github/workflows/ci.yml` (one verbose step)

This could conflict with any attended branch that touches the verifygate tests (row 135 follow-ups).
Re-fetch `origin/dev` before merging.

## Residuals

- **R1.** `syscall.ForkLock` covers forks in *this* process only. Writes made by other processes,
  such as gate scripts writing files that a later bash execs, are a different class and are not
  in play here: no gate script writes and then execs its own output in an overlapping process.
- **R2.** Write-then-exec in serial packages outside verifygate (listed under the enumeration).
  These are exposed only if a background goroutine forks while the test writes. They have not
  been observed and are not routed here. A follow-up row would widen the scan to `host/` and
  `cmd/` test files.
- **R3.** AC5 and test 2 can only be checked on Linux CI. The rig cannot drill them.
- **R4.** Test 1's 250 ms window matters only in the mutant direction: a `Start` slower than
  250 ms on a loaded machine could hide a mutant. The correct direction cannot flake, because the
  fork physically cannot happen while the lock is held.
