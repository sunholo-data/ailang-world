# World iteration 231 — independent sprint evaluator r1 (row 143)

VERDICT: PASS 94/100

Judged: `git diff 95575aa..HEAD` on `sprint/w-verifygate-forklock-guard-residuals-2`
(worktree ~/dev/sunholo-data/.wt-world-iter230-forklock2; HEAD df7ac90; M1 e0ccf3f, M2 613f94c, M3 df7ac90).
Design: design_docs/planned/w-verifygate-forklock-guard-residuals-2.md (r3). Plan: ...-sprint-plan.md.
Binary: AILANG_BIN=~/.pinned-ailang/ailang -> `AILANG v0.41.0`; go1.26.6 darwin/arm64.
All evidence below is judge-run. The executor's logs were not used. Harness: ~/.ailang/state/w231_eval/mut.sh.
Raw outputs: ~/.ailang/state/w231_eval/run{1,2,3}.txt, race.log.

## Scope
- numstat: 1 code file `host/verifygate/forklocked_write_test.go` +39/-2; 2 planned design docs. No production code,
  no .ail, no ci.yml, no scripts, no tools/launchd.
- ci.yml l.207 `tests=` names all five fork-lock tests. The diff adds no top-level `func Test`, so no ci.yml change is needed.
- Base anchors re-checked at 95575aa: l.337 bannedOSFuncs, l.469 identity check, l.565 `}` slice close (the plan's
  correction of the design's l.566 is right), l.647 writeSel, l.697 exemption, l.746 EOF.

## AC check
- AC1 (live CopyFS escape): KILLED. Observed `ail_binary_gate_test.go:763: raw os.CopyFS in func evalCopyFS230: use a fork-locked wrapper`.
  With N10a also applied, the live scan goes silent and only the fixture fires. So before this row the escape
  went undetected, which confirms V5.
- AC2/AC3 (N10 scan/tripwire): KILLED at l.590 / l.765 (texts below).
- AC4 (N3) and AC5 (N5): KILLED at l.590 / l.589. Base controls: the same mutants at pristine 95575aa give rc=0 (survive),
  so the new fixtures are what kills them.
- AC6 (N6): KILLED at l.773 `known host positive not named: []`. The live sweep stays green, which shows only the fixture
  catches it. Base control rc=0.
- AC7: with `if true {` the live red fires with exactly 10 `t.Parallel in a write+fork` lines (counted, matches V14). The
  fixture red cannot appear in the same run because the live Fatalf comes first. The design's "fixture AND live" premise
  is false; the plan already flags this in §1.3/§7.2. Judge drill T4 turns off the live assertion and shows the negative
  fixture firing: l.781 `exempt host/verifygate fixture flagged: [host/verifygate/x_test.go:10: ...]`.
- AC8: all five fork-lock tests pass plain (4 PASS + darwin SKIP) and under -race. The full verifygate package passes under -race.
  Floors untouched.

## Mutation table (each mutant applied, single targeted test run, reverted; sha256 of both touched files byte-identical after every run:
forklocked_write_test.go 7d12cd18…1c99, ail_binary_gate_test.go 8ee31341…ad96)

S = ^TestVerifygateTestWritesAreForkLocked$, T = ^TestNoParallelWriteForkPackagesOutsideVerifygate$

### Guard-logic mutants
| # | Mutant | Test | Result | Observed |
|---|---|---|---|---|
| N10a | drop `"CopyFS": true` from bannedOSFuncs | S | KILLED | :590 fixture CopyFS: fixture_CopyFS.go:5\|raw os.CopyFS in func h not reported (got []) |
| N10b | drop `"CopyFS": true` from writeSel["os"] | T | KILLED | :765 known CopyFS positive not named: [] |
| N10c | live `func evalCopyFS230() { _ = os.CopyFS("dst", os.DirFS("testdata")) }` appended to ail_binary_gate_test.go | S | KILLED | :527 raw (non-fork-locked)… / ail_binary_gate_test.go:763: raw os.CopyFS in func evalCopyFS230 |
| N10c+N10a | live escape + ban removed | S | KILLED (fixture only, live silent) | :590 fixture CopyFS … not reported |
| N3 | `info.Uses[fid] != doObj` → `fid.Name != "forkLockedDo"` | S | KILLED | :590 fixture LocalDo: fixture_LocalDo.go:9\|raw os.OpenFile in func h not reported (got []) |
| N5 | drop `&& (fo.Name() == "OpenFile" \|\| fo.Name() == "CreateTemp")` | S | KILLED | :589 fixture CreateOpener: fixture_CreateOpener.go:8\|raw os.Create in func h not reported (got []) |
| N6 | `dir != "host/verifygate"` → `!strings.HasPrefix(dir, "host/")` | T | KILLED | :773 known host positive not named: [] |
| AC7 | exemption → `if true {` | T | KILLED | :727 parallel tests in write+fork packages… (10 lines) |
| O1 | opener check → `fo.Name() != "WriteFile"` (widened opener set) | S | KILLED | :590 fixture CreateOpener … not reported |
| O2 | identity → `info.Uses[fid] == nil` (any resolved callee) | S | KILLED | :590 fixture LocalDo … not reported |
| O3 | identity → `info.Uses[fid] == nil \|\| info.Uses[fid].Name() != "forkLockedDo"` (object-name compare) | S | KILLED | :590 fixture LocalDo … not reported |
| O4 | CopyFS moved from writeSel["os"] to forkSel["os"] (misclassified) | T | KILLED | :765 known CopyFS positive not named: [] |
| O5 | exemption → `!strings.HasPrefix(dir, "host/verifygate")` | T | SURVIVED | ok (also exempts host/verifygate/* subdirs and host/verifygatex) |
| O6 | opener exemption: drop `fo.Pkg() != nil && fo.Pkg().Path() == "os" &&` | S | SURVIVED (equivalent) | ok. Only banned os/syscall objects reach `exempt[id]`, and no banned syscall name is OpenFile/CreateTemp |
| O7 | CopyFS moved from writeSel["os"] to writeSel["syscall"] | T | KILLED | :765 known CopyFS positive not named: [] |
| O8 | exemption → `dir != "host/verifygate" && !strings.HasPrefix(dir, "host/z")` | T | KILLED | :773 known host positive not named: [] |

Guard mutants: 16 run, 14 killed. 1 survivor is equivalent (O6), 1 is minor and real (O5).

### Test-side and vacuity probes (each weakens or changes a shipped fixture, sometimes paired with a guard mutant)
| # | Probe | Test | Result | Meaning |
|---|---|---|---|---|
| T1 | LocalDo without `doDecl` + N3 | S | SURVIVED | the shipped doDecl is load-bearing: without it doObj is nil and N3 goes unseen. Shipped fixture is correct |
| T1b | LocalDo without `doDecl`, alone | S | SURVIVED (green) | same: the variant is still green, so only the shipped shape pins N3 |
| T2 | `host/zz` key → `hostzz` + N6 | T | SURVIVED | the `host/` prefix of the shipped key is load-bearing. Shipped key is correct |
| T3 | negative key `host/verifygate` → `host/verifygatex`, alone | T | KILLED | :781 exempt host/verifygate fixture flagged: [host/verifygatex/…]. The exact-name compare is live |
| T4 | AC7 + live assertion turned off (`if false {`) | T | KILLED | :781 exempt host/verifygate fixture flagged: [host/verifygate/x_test.go:10 …]. The new negative is reachable and fires |
| T5 | bodyCopyFS uses os.WriteFile instead of CopyFS + N10b | T | SURVIVED | the shipped CopyFS line is what pins writeSel. Not vacuous |
| T6 | LocalDo want text dropped to the `:9` prefix + N3 | S | KILLED | :590 fixture LocalDo: fixture_LocalDo.go:9 not reported. The text is redundant but harmless |
| T7 | CreateOpener returns os.OpenFile, alone | S | KILLED | :590 … raw os.Create in func h not reported. The fixture text pins the opener under test |

Base controls (mutant on pristine 95575aa): N3 rc=0, N5 rc=0, N6 rc=0, so all three survived before this row.

Total: 24 mutant/probe runs + 3 base controls. Guard mutants 14/16 killed (1 equivalent).

## Vacuity
Every new fixture fails when the code it protects is removed:
- CopyFS scan fixture ← N10a
- zzcopyfs ← N10b/O4/O7
- LocalDo ← N3/O2/O3
- CreateOpener ← N5/O1
- host/zz ← N6/O8
- host/verifygate negative ← T3/T4

All four new Fatalf texts were observed firing, so none is unreachable. The `len(x) != 1 ||` count halves are not
needed to produce a red: an empty slice would panic on the `[0]` index anyway. They are harmless.

## Findings
BLOCKING: none.

NON-BLOCKING:
1. O5 survives. Mutant text: `if !strings.HasPrefix(dir, "host/verifygate") {` at l.706.
   - Effect: it would exempt a future `host/verifygate/<sub>` (or `host/verifygatex`) package from the tripwire.
   - Why it matters: the full scan covers only `.` (the package itself), so such a package would be covered by neither guard.
   - Today no such directory exists, so on the present tree it behaves like the original.
   - Fix: a one-line positive keyed `host/verifygate/zz`.
2. Design evidence drift: the design doc (planned/, r3) still has two wrong statements.
   - AC8 cites `AILANG_BIN=~/.pinned-ailang-tools/v0.52.1/ailang`. The repo pin is v0.41.0, and v0.52.1
     reds the full suite at base.
   - AC7 claims "fixture AND live red" in one run, which cannot happen because the live Fatalf fires first.
   - The plan measures and flags both in §1.1/§1.3/§7. Fold them into the design (or its implemented copy) at landing.
3. O6 is an equivalent mutant: the opener package check is unreachable-redundant. Informational only.

## Gates (judge-run, v0.41.0)
- `go vet ./host/verifygate/` rc 0
- the five fork-lock tests `-v`: 4 PASS + SKIP (darwin, by design), rc 0
- same four tests `-race` rc 0
- `go build ./...` rc 0
- `go test ./host/verifygate/ -count=1 -race` rc 0 (351.5 s), 0 `DATA RACE`
- `gofmt -l host/verifygate/` empty
- `check_no_personal_email.sh` rc 0
- `git status --porcelain` empty at finish (asserted)

## Score
- Correctness/AC coverage: 39/40. All ACs met. The AC7 design premise is wrong, but the plan corrected it.
- Test strength: 27/30. 14/16 guard mutants killed, all design-named mutants killed, 1 minor real survivor (O5).
- Scope/discipline: 15/15. One test file, no ci.yml/production change, no new top-level test.
- Evidence integrity: 13/15. Plan claims re-measured true (anchors, kill lines, 10-count, base survivals). The design
  doc's stale binary and AC7 premise are not yet corrected in the doc itself.
- Total: 94/100. PASS.

## Round 2

Judged: `git diff df7ac90..d5d569e` only. The commit adds an M4 test fixture (+8 lines in forklocked_write_test.go) and
edits the design doc (AC7, AC8, Residuals). Pre-run sha256 of forklocked_write_test.go is 05c23043…0bb2. Every mutant
was reverted to that exact sha. ail_binary_gate_test.go is unchanged at 8ee31341…ad96.

### (a) Mutants (T = ^TestNoParallelWriteForkPackagesOutsideVerifygate$; raw output in ~/.ailang/state/w231_eval/run_r2.txt)
| # | Mutant | Result | Observed |
|---|---|---|---|
| R2a | (r1 O5) `if !strings.HasPrefix(dir, "host/verifygate") {` | KILLED | :789 known host/verifygate/zz positive not named: [] |
| R2b | `if !strings.HasPrefix(dir, "host/verifygat") {` | KILLED | :789 known host/verifygate/zz positive not named: [] |
| R2c | `if !strings.Contains(dir, "verifygate") {` | KILLED | :789 known host/verifygate/zz positive not named: [] |
| R2d | new fixture's `len(sv) != 1 \|\|` dropped + R2a | KILLED | panic: index out of range [0] with length 0. The count check is redundant but harmless |
| R2e | `if !strings.HasSuffix(dir, "verifygate") {` | SURVIVED | ok. This would exempt any other directory named `verifygate`, e.g. `cmd/verifygate`. `find` shows only `./host/verifygate` exists, so it behaves the same on today's tree. Not blocking |
| R2f | probe: sv key changed to `host/zz` + R2a | SURVIVED (expected) | the shipped key `host/verifygate/zz` is what catches R2a, so the fixture is not vacuous |

Finding 1 from r1 is FIXED: the O5 mutant is now KILLED.

### (b) Doc corrections
- AC7 is accurate. My r1 run measured the live red with exactly 10 lines and the fixture firing only when the live
  assertion was turned off (T4, :781).
- AC8 now names `~/.pinned-ailang/ailang`, which is v0.41.0 and is the repo pin.
- The Residuals note is accurate.

Finding 2 from r1 is FIXED.

### (c) Gates (v0.41.0)
- `go vet ./host/verifygate/` rc 0.
- Five fork-lock tests plain: rc 0 (4 PASS + darwin SKIP).
- Same five with `-race`: rc 0, 0 `DATA RACE`.
- `gofmt -l host/verifygate/` empty.
- `git status --porcelain` empty at the end (checked).

### Findings r2
BLOCKING: none.
NON-BLOCKING: R2e (`HasSuffix(dir, "verifygate")`) survives. It is a cosmetic gap that behaves the same on today's tree.
O6 stays equivalent.

Score r2 is 97/100:
- correctness 40
- test strength 28
- scope 15
- evidence 14

VERDICT r2: PASS 97/100
