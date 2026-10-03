# Executor report — row 142 (w-verifygate-forklock-guard-residuals), iteration 229

Lane: claude-sonnet-5-5. Worktree `.wt-world-iter229`, branch `sprint/w-verifygate-forklock-guard-residuals`.
Pin: `AILANG v0.41.0`. Commits: M1 `000ef12`, M2 `6a34737`, M3 `3c68ece`. Evidence logs: `~/.ailang/state/w229_*`.
Touched: `host/verifygate/forklocked_write_test.go`, `.github/workflows/ci.yml` (one step). No production code, no `.ail`, no `tools/launchd/*`.

## Gates (rc captured directly, no pipes)

| Gate | M1 | M2 | M3 |
|---|---|---|---|
| gofmt -l host/verifygate | empty | empty | empty |
| go vet ./host/verifygate/ | rc 0 | rc 0 | rc 0 |
| FAST (named tests, -v) | 2 PASS, kernel SKIP, scan PASS | same | 4 PASS (+ NoParallel), kernel SKIP |
| -race subset | rc 0, 0 DATA RACE | rc 0, 0 | rc 0, 0 |
| go test ./host/verifygate/ -count=1 | rc 0; 151 pass + 1 skip, 0 fail | rc 0; 151 + 1 | rc 0; 152 + 1 (as the plan expected) |
| verify_ail.sh | rc 0 | rc 0 | rc 0 |
| check_no_personal_email.sh | - | - | rc 0 |
| ci.yml YAML parse (ruby) | - | - | rc 0 |
| go test ./... -count=1 | - | - | rc 0 (26 ok packages) |

Live scan: zero violations on first run at M2 (no contingency hit); floors (>=9 files, >=21 wrapper sites, identity-counted) hold.
M3 sweep floors (>=26 dirs, >=12 write+fork dirs) hold; zero live violations.

## CI step body, run locally (extracted from ci.yml via ruby YAML)

1. `RUNNER_TEMP=$(mktemp -d) bash w229_step.sh` (RUNNER_OS unset): output `::error::this step must run on linux`, rc 1.
2. `RUNNER_OS=Linux ...`: the 4 other tests PASS, `--- SKIP: TestKernelRefusesExecOfWriterOpenFile`, then
   `::error::TestKernelRefusesExecOfWriterOpenFile did not PASS on linux (SKIP?)`, rc 1.

## Mutation drills (each reverted; `git diff --quiet` clean after every revert)

| Drill | Mutant | Result | Failing line |
|---|---|---|---|
| AC3 (E3) | RUnlock before Close (deferred one removed) | KILLED (5/5 wrappers); `ConcurrentFork` and scan tests stayed PASS | `forklocked_write_test.go:278: writeFileForkLocked: ForkLock released before Close returned` (same for all five) |
| AC1 probe (E1) | createTempForkLocked raw CreateTemp+Write+Close | KILLED | `createTempForkLocked: fork-lock probe stages = [], want [filled closed]` |
| AC2 probe (E2) | forkLockedWrite exec-mode raw branch | KILLED | `writeFileForkLocked: fork-lock probe stages = [], want [filled closed]` (also copyFileForkLocked, forkLockedWrite) |
| AC1 scan | same E1 body | KILLED | `forklocked_write_test.go:108: raw os.CreateTemp in func createTempForkLocked: use a fork-locked wrapper` |
| AC2 scan | same E2 branch | KILLED | `forklocked_write_test.go:73: raw os.OpenFile in func forkLockedWrite: use a fork-locked wrapper` |
| AC4 (E5b) | `wf := os.WriteFile` helper in ail_binary_gate_test.go | KILLED | `ail_binary_gate_test.go:763: raw os.WriteFile in func h: use a fork-locked wrapper` |
| AC5 (E8b) | syscall.Open + os.NewFile helper (+ syscall import) in toolchain_pin_gate_test.go | KILLED | `toolchain_pin_gate_test.go:1780: raw syscall.Open in func h` and `:1781: raw os.NewFile in func h` |
| AC6 (E6b) | skip uses with no enclosing FuncDecl | KILLED | `fixture PkgInit: fixture_PkgInit.go:5\|package-level initializer not reported (got [])` |
| AC7 (F6) | by-name exemption for `func forkLockedWrite` | KILLED | `fixture NamedWrapper: fixture_NamedWrapper.go:5\|raw os.OpenFile in func forkLockedWrite not reported (got [])` |
| AC10 | delete exempt-ident marking | KILLED | live red names `:110 raw os.CreateTemp in func createTempForkLocked` and `:72 raw os.OpenFile in func forkLockedWrite` (N-OpenClosure not observable, plan override 6) |
| AC12 | raw `os.WriteFile` first stmt of createTempForkLocked's open closure | KILLED | `forklocked_write_test.go:110: raw os.WriteFile in func createTempForkLocked`; the `os.CreateTemp` opener is absent from the log (grep count 0) |
| AC9 (R2) | `t.Parallel()` in host/broker/broker_test.go | KILLED | `host/broker/broker_test.go:41: t.Parallel in a write+fork package outside the fork-locked scan` |
| Floor | live root `../../host` -> `t.TempDir()` | KILLED | `sweep saw 2 Go dirs, want >= 26` |
| Positive | matcher `"Parallel_"` | KILLED | `known positive not named: []` |

No mutant survived. AC8 and AC11 (linux runner) are CI-only, not executed here.

## Deviations / notes

- Followed all six plan overrides. Wrapper-identity check runs only on the live scan (fixtures skip it); P-Smuggle and
  N-OpenClosure fixtures define their own `forkLockedDo`; mode checks are umask-independent; AC4 used `wf("x", nil, 0o755)`;
  AC5 added the `syscall` import.
- Fixture `want` entries support a `file:line|text` form (text checked against the reported line) so P-FuncValue, P-PkgInit,
  P-NamedWrapper etc. assert the `where` text; an `absent` list asserts the exempt Smuggle opener (line 10) is not reported.
- N-OpenClosure and a plain-wrapper negative run as a map of negatives (zero violations each).
- The scan reports at the identifier position (`id.Pos()`), not `call.Pos()`; same line in every case.
- Drill AC6 first attempt accidentally applied a no-op edit; it was reverted unrun and the drill redone properly. No effect on results.
- Floor drill: the `t.TempDir()` replacement yields 2 dirs (cmd only), message `sweep saw 2 Go dirs`.
