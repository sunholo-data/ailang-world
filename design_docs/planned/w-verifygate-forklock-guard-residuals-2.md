# w-verifygate-forklock-guard-residuals-2

**Status**: PLANNED · **Queue row**: 143 · **Clause**: clause-2 (MET — hygiene; below every routable UNMET-clause row, rule (e))
**Author lane**: pi:ollama/glm-5.3:cloud (iteration 230 designer) · **Date**: 2026-10-04 · **Base**: `95575aa`
**Parent**: [w-verifygate-forklock-guard-residuals](../implemented/w-verifygate-forklock-guard-residuals.md) (landed `2da63de`, #198;
judge verdict [evaluator-r1](../verification/world-iter229/evaluator-r1.md), PASS 92, zero blocking; new mutants N3, N5, N6, N10, N15b SURVIVED).
**Revision**: r1.
**Scope**: `host/verifygate/forklocked_write_test.go` only. Test-only. No production code, no `.ail`, no `tools/launchd/*`, no
gate-script semantics, and no `ci.yml` change (V10: this row adds no top-level test name, so the CI step's list cannot drift). Estimate ~0.25 d.

## Problem

Row 142's guard is correct, but the iter-229 judge killed 11 of 16 new mutants and measured five
survivors. Three are test holes in the two guards; two are design residuals this row records
rather than fixes. Source of truth: the judge's mutation table and non-blocking findings 1–4
([evaluator-r1](../verification/world-iter229/evaluator-r1.md)). Line cites are this worktree at `95575aa`.

1. **N10 — `os.CopyFS` evades both guards** (finding 1). It writes files, keeps exec bits, and
   opens its fds inside package `os` with no ForkLock (`go doc`: files are created with "mode
   0o666 plus any execute permissions from the source", V3). It is in neither `bannedOSFuncs`
   (l.337) nor the tripwire's `writeSel["os"]` (l.647), and zero live uses exist today (V5), so
   it is now the cheapest live escape: the judge staged `os.CopyFS(dst, os.DirFS("testdata"))`
   in `ail_binary_gate_test.go` and everything stayed green.
2. **N6 — the tripwire's directory exemption is unpinned** (finding 2). The exemption is the
   bare compare `dir != "host/verifygate"` (l.697), and the only positive fixture uses the
   synthetic dir `zz` (l.730). A one-line mutant (`!strings.HasPrefix(dir, "host/")`) silently
   neuters the R2 tripwire for every host package, and nothing fires.
3. **N3/N5 — the exemption's identity and opener-name checks are unpinned** (finding 3). The
   scan exempts an open closure only when the callee is the package-scope `forkLockedDo` object
   (l.415, l.469 `info.Uses[fid] != doObj`) and only when the returned opener is
   `os.OpenFile`/`os.CreateTemp` (l.489). Every `return os.` in the file is one of those two
   names (V8), so no fixture pins a local/shadowed `forkLockedDo` callee (must REPORT) or a
   closure returning `os.Create` (must REPORT — only the two named openers are exempt).
4. **N15b — D1 is a two-point sample** (finding 4). The probe observes the lock at "filled"
   (l.62, before `Close` at l.64) and at "closed" (l.66, after), never during `Close` (V6), so a
   lock dropped only across `Close` survives. The judge accepts this for the realistic E3 shape:
   recorded as a known limit of D1 in Residuals. No test.
5. **N13 — a fork inside `fill` is caught only by `-timeout`** (the R1 deadlock; the judge's N13
   drill panicked at 30 s). Recorded as design residual R1. No watchdog/timeout test in this row.

The design widens no further: no ban beyond `CopyFS`, no exemption change, no wrapper change,
no CI change, no new top-level test. Both guards' existing tests deepen in place.

## Milestones

**M1 — `os.CopyFS` banned in both guards (kills N10).** File: `host/verifygate/forklocked_write_test.go`.
Tests touched (names, no new top-level test): `TestVerifygateTestWritesAreForkLocked`,
`TestNoParallelWriteForkPackagesOutsideVerifygate`.
- Add `"CopyFS": true` to `bannedOSFuncs` (l.337) and to `writeSel["os"]` (l.647).
- Scan fixture, appended to the positives slice after the `Smuggle` entry (l.563–565, before the
  slice's `}` at l.566):
  `{"CopyFS", "package p\n\nimport \"os\"\n\nfunc h() { _ = os.CopyFS(\"dst\", os.DirFS(\"src\")) }\n", []string{"fixture_CopyFS.go:5|raw os.CopyFS in func h"}, nil}`.
  It type-checks under the scan's source importer (`os.DirFS` comes from the same `"os"` import;
  no `io/fs` identifier appears).
- Tripwire positive, added to the fixture block before the test's closing `}` (l.746), with a
  `bodyCopyFS` const = the `body` template (l.729) with the write line replaced by
  `_ = os.CopyFS("dst", os.DirFS("src"))` (the parallel call stays at l.10), keyed `zz` —
  deliberately NOT host-prefixed, so this fixture isolates the `writeSel` claim from M3's
  dir-exemption claim. Expect exactly 1 violation prefixed
  `zz/x_test.go:10: t.Parallel in a write+fork package` (message l.699).
- Acceptance (mutant → firing assertion):
  - **AC1 (N10, live escape)**: add `func evalCopyFS230() { _ = os.CopyFS("dst", os.DirFS("testdata")) }`
    to `ail_binary_gate_test.go` (scratch worktree) → the live scan reds
    `ail_binary_gate_test.go:<line>: raw os.CopyFS in func evalCopyFS230`. Today it is silent (V5).
  - **AC2 (N10, scan)**: drop `"CopyFS": true` from `bannedOSFuncs` →
    `fixture CopyFS: fixture_CopyFS.go:5|raw os.CopyFS in func h not reported (got [])` (l.581).
  - **AC3 (N10, tripwire)**: drop `CopyFS` from `writeSel["os"]` → the `zz` CopyFS fixture reds
    `known positive not named: []` (l.736).
- Gate: `go vet ./host/verifygate/`; `go test ./host/verifygate/ -run
  'TestVerifygateTestWritesAreForkLocked|TestNoParallelWriteForkPackagesOutsideVerifygate' -count=1 -v`,
  plain and `-race`. The live scan must stay at zero violations (V5) and the floors must not move
  (no live `CopyFS` use changes `filesSeen`/`wrapperCalls`/`nWriteFork`).

**M2 — exemption identity and opener name pinned (kills N3, N5).** Same file, same positives
slice (after M1's `CopyFS` entry). Test touched: `TestVerifygateTestWritesAreForkLocked`.
- **P-LocalDo (N3)**: `{"LocalDo", "package p\n\nimport \"os\"\n\n" + doDecl + "\nfunc h() {\n\tforkLockedDo := func(open func() (*os.File, error), fill func(*os.File, error), hold func()) error { return nil }\n\t_ = forkLockedDo(func() (*os.File, error) { return os.OpenFile(\"x\", 0, 0) }, nil, nil)\n}\n", []string{"fixture_LocalDo.go:9|raw os.OpenFile in func h"}, nil}`.
  The local callee is a `*types.Var`, so the identity check (l.469) refuses the exemption; the N3
  bare-name mutant (`fid.Name != "forkLockedDo"`) exempts it and the fixture reds.
- **P-CreateOpener (N5)**: `{"CreateOpener", "package p\n\nimport \"os\"\n\n" + doDecl + "\nfunc h() {\n\t_ = forkLockedDo(func() (*os.File, error) { return os.Create(\"x\") }, nil, nil)\n}\n", []string{"fixture_CreateOpener.go:8|raw os.Create in func h"}, nil}`.
  `os.Create` is banned (l.337) and only `OpenFile`/`CreateTemp` are exempt openers (l.489); the
  N5 mutant (any `os` opener exempt) reds.
- The existing N-OpenClosure (l.596) and P-Smuggle (l.563) stay unchanged and green: the
  exemption stays live and stays narrow (only the returned `os.OpenFile` opener is exempt).
- Acceptance:
  - **AC4 (N3)**: apply the judge's N3 mutant at l.469 →
    `fixture LocalDo: fixture_LocalDo.go:9|raw os.OpenFile in func h not reported (got [])`.
  - **AC5 (N5)**: apply the judge's N5 mutant at l.489 (drop the name check) →
    `fixture CreateOpener: fixture_CreateOpener.go:8|raw os.Create in func h not reported (got [])`.
- Gate: `go vet`; the scan test plain and `-race`.

**M3 — tripwire directory exemption pinned (kills N6).** Same file, fixture block of
`TestNoParallelWriteForkPackagesOutsideVerifygate` (before l.746).
- **Positive keyed `host/zz`** (finding 2): the existing `body` (l.729) with
  `fmt.Sprintf(body, "\tt.Parallel()\n")`, keyed `"host/zz"` → exactly 1 violation prefixed
  `host/zz/x_test.go:10: t.Parallel in a write+fork package`. Under N6
  (`!strings.HasPrefix(dir, "host/")`) it is silently exempt, so the fixture reds — and the live
  sweep stays green under that mutant either way (V7), which is exactly why the fixture is
  load-bearing.
- **Negative keyed `host/verifygate`**: the same parallel body keyed `"host/verifygate"` →
  0 violations. It pins the exemption at fixture level; replacing the l.697 guard with `true`
  reds this fixture (`known negative flagged: [host/verifygate/x_test.go:10 …]`, l.744) and the
  live test (the 10 real `t.Parallel()` calls, V7).
- The existing `zz` positive (l.730) and no-Parallel negative (l.738) stay untouched.
- Acceptance:
  - **AC6 (N6)**: apply the judge's N6 mutant at l.697 → `known positive not named: []` on the
    `host/zz` fixture; the live test stays green, proving only the fixture sees it.
  - **AC7 (exemption liveness)**: replace l.697 with `true` → the `host/verifygate` negative
    fixture reds AND the live test reds with 10 real violations (V7).
  - **AC8 (regression)**: every row-142 AC stays green; the six iter-228 survivors stay dead;
    floors unchanged (`filesSeen` ≥ 9, `wrapperCalls` ≥ 21, `nDirs` ≥ 26, `nWriteFork` ≥ 12);
    all five fork-lock tests PASS plain (baseline 6.3 s, V11) and `-race`;
    `verify_go.sh`, `verify_ail.sh` (no `.ail` touched) and `check_no_personal_email.sh` pass
    with no gate-script change.

## Verification Log

Measured in `/Users/voightkampff/dev/sunholo-data/.wt-world-iter230-forklock2` at `95575aa`,
on go1.26.6 darwin/arm64.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Base, branch, clean tree, parent landing | `git rev-parse --short HEAD; git branch --show-current; git status --porcelain \| wc -l; git log --oneline -3` | `95575aa`; `sprint/w-verifygate-forklock-guard-residuals-2`; `0` lines dirty; parent landed as `2da63de` ("test(row142) … (#198)") directly beneath HEAD |
| V2 | Toolchain | `grep -n 'go 1' go.mod; go version` | `3:go 1.26.6`; `go1.26.6 darwin/arm64` |
| V3 | `os.CopyFS` exists and keeps exec bits | `go doc os.CopyFS` | `func CopyFS(dir string, fsys fs.FS) error`; "Files are created with mode 0o666 plus any execute permissions from the source" — matches judge N10 verbatim |
| V4 | The two ban sets and the exemption, and that `CopyFS` is absent | `grep -n 'var bannedOSFuncs\|var bannedSyscallFuncs' host/verifygate/forklocked_write_test.go`; `grep -n 'writeSel :=\|dir != "host/verifygate"' <same>`; `grep -n 'CopyFS' <same>; echo rc=$?`; control: `grep -c 'NewFile' <same>` | l.337 `bannedOSFuncs` = {WriteFile, OpenFile, Create, CreateTemp, NewFile} (no `CopyFS`); l.338 `bannedSyscallFuncs`; l.646–648 `writeSel` os set = the same five names; l.697 `if dir != "host/verifygate" {`; `CopyFS` grep: no output, `rc=1`; control `6` — the zero is the symbol's absence, not a dead instrument |
| V5 | Zero live `CopyFS` uses, so the widening reds nothing that exists | `grep -rn 'CopyFS' --include='*.go' host cmd; echo rc=$?`; `grep -rn 'CopyFS' host/verifygate; echo rc=$?`; control: `grep -rln 'NewFile' --include='*.go' host cmd` | both CopyFS greps: no output, `rc=1`; control lists 17 files (incl. `forklocked_write_test.go`, `toolchain_pin_gate_test.go`), so the same instrument finds a present symbol |
| V6 | D1 samples the lock at two points only (N15b) | `sed -n '55,68p' host/verifygate/forklocked_write_test.go` | l.61–62 probe `"filled"`; l.64 `err = f.Close()`; l.65–66 probe `"closed"`. No probe call between l.62 and l.66 — the lock is unobserved during `Close` |
| V7 | `t.Parallel()` lives only in `host/verifygate`, so the l.697 exemption is load-bearing | `grep -rln 't\.Parallel()' --include='*.go' host cmd`; then `grep -n 't\.Parallel()'` on each hit | only `host/verifygate/{mission_config,module_manifest}_gate_test.go` and this file. Real calls: 1 (mission_config l.246) + 9 (module_manifest l.255–549) = 10; this file's hits (l.641, l.728, l.730) are two comments and the fixture string literal. Outside verifygate: 0 |
| V8 | The identity and opener-name checks (N3/N5 targets) are unpinned by fixtures | `grep -n 'Lookup("forkLockedDo")\|info.Uses\[fid\] != doObj\|fo.Name() == "OpenFile" \|\| fo.Name() == "CreateTemp"\|info.Defs\[fd.Name\] == kernelObj' <file>`; `grep -n 'return os\.' <file>` | l.415 `doObj := pkg.Scope().Lookup("forkLockedDo")`; l.469 `!ok \|\| info.Uses[fid] != doObj`; l.489 `(fo.Name() == "OpenFile" \|\| fo.Name() == "CreateTemp")`; l.455 kernel identity `info.Defs[fd.Name] == kernelObj`. Every `return os.` in the file (l.72, l.110, l.563, l.596) returns `OpenFile` or `CreateTemp`, and every fixture callee is the package-scope `doDecl` — no fixture exercises a non-package-scope callee or a non-exempt opener |
| V9 | Fixture inventory and insertion anchors | `grep -n 'const doDecl\|"RootWrite"\|"Smuggle"\|"OpenClosure"\|not reported\|exempt opener\|const body\|pos := map\|neg := map\|known positive not named\|known negative flagged\|want >= 26\|want >= 12\|t.Parallel in a write+fork' <file>` | doDecl l.538; positives RootWrite l.559, Smuggle l.563–565, slice closes l.566; failure texts l.581 (`not reported`) and l.587 (`exempt opener … reported`); negative OpenClosure l.596; tripwire `body` l.729, `pos` keyed `zz` l.730, prefix assertion `zz/x_test.go:10` l.735–736, `neg` l.738, `known negative flagged` l.744; floors l.721 (`>= 26`) / l.724 (`>= 12`); message l.699 |
| V10 | The CI verbose step needs no change | `grep -n 'Fork-locked' .github/workflows/ci.yml`; `sed -n '196,216p' .github/workflows/ci.yml` | step `Fork-locked write tests, verbose (w-verifygate-etxtbsy linux gate)` at l.201; its `tests=` list names all five fork-lock tests — including both tests this row deepens — and the loop asserts a top-level `--- PASS: ` line per test. M1–M3 add no top-level test name, so the list cannot drift and no `ci.yml` edit is made |
| V11 | Live baseline of the five tests at the base | `go test ./host/verifygate/ -run 'TestForkLocked\|TestVerifygateTestWritesAreForkLocked\|TestNoParallel\|TestKernelRefuses' -count=1 -v` | `PASS TestForkLockedWriteBlocksConcurrentFork (0.26s)`; `PASS TestForkLockedWrappersHoldLockThroughClose (0.00s)`; `SKIP TestKernelRefusesExecOfWriterOpenFile` (darwin, as designed); `PASS TestVerifygateTestWritesAreForkLocked (5.76s)`; `PASS TestNoParallelWriteForkPackagesOutsideVerifygate (0.10s)`; `ok 6.296s` |

**Instrument note.** A repo-wide `grep -rn 'os.CopyFS' --include='*.go' . | head -3` was also run
and is NOT cited: the pipe makes `$?` report `head`'s exit, not grep's (the parent's V6a hazard).
Every zero above is an un-piped grep whose rc was echoed.

## Conflict surface

- `host/verifygate/forklocked_write_test.go`: all three milestones — two map literals (l.337,
  l.647) and the fixture lists of the two guard tests. No production symbol, no wrapper, no gate.
- `.github/workflows/ci.yml`: NO change (V10). Re-check the step before merging only if a sibling
  row edits it first.
- The executor re-measures every line cite if anything lands between this doc and execution;
  the cites are from `95575aa`.

## Residuals

- **R1 (carried from the parent; judge-confirmed by N13).** A fork inside the `RLock` region
  (`fill`, `hold`, the probe, or an open closure) self-deadlocks on `ForkLock.Lock` and is caught
  only by `go test`'s `-timeout` (30 s locally, 120 s in the CI step). N13 measured exactly that.
  No watchdog or timeout test is added in this row; the probe's MUST-NOT-fork doc comment stays
  the guard.
- **R6 (new, judge N15b) — D1 is a two-point sample.** The probe observes the lock at "filled"
  and "closed" only (V6); a mutant that releases and re-takes the lock across `Close` survives.
  The judge accepts this for the realistic E3 shape. Recorded as a known limit of D1; no test.
- Carried unchanged from row 142: **R2′** the tripwire is syntactic (an aliased import evades
  it); **R3** a stray `RLock` holder can hide E3-class mutants (false-negative direction only);
  **R4** `syscall.Openat`/`Creat` match only on linux; **R5** a kernel that drops write-deny reds
  CI via the kernel control — correct and loud.

## Quorum

Round results pending (artifacts to land under `design_docs/verification/world-iter230/`).

- **r1**: _pending_ — designer `pi:ollama/glm-5.3:cloud`; reviewers per the routing table.
- **r2** (only if r1 blocks): _pending_.