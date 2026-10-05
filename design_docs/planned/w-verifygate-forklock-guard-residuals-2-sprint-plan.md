# Sprint plan — w-verifygate-forklock-guard-residuals-2 (iteration 231, row 143)

**Design**: `design_docs/planned/w-verifygate-forklock-guard-residuals-2.md` (r3, `5f5228e`). **Branch**:
`sprint/w-verifygate-forklock-guard-residuals-2`, base `95575aa` (= `origin/dev`). **Scope**:
`host/verifygate/forklocked_write_test.go` only. Test-only: no production code, no `.ail`, no
`tools/launchd/*`, no `scripts/*`, no `ci.yml` change. **Size**: 3 milestones, ~0.25 d, one executor session.
Commit each milestone separately; every drill reverts to that milestone's commit and proves byte identity.

**Planner prototype.** The whole M1–M3 change was implemented in a scratch copy of `95575aa`
(`git archive HEAD` → `~/.ailang/state/w231_plan/scratch`), every mutant in §4 was run against it, and
the survival controls were run against the pristine base. The diff is banked at
`~/.ailang/state/w231_plan/impl.patch` (sha256 `fde5b058…006e058e`); `git apply --check` of it in the
worktree is rc 0. The executor may apply it as-is or type it from §2 — the two are byte-identical.
Drill logs: `~/.ailang/state/w231_plan/drill_<mutant>.log`.

## 0. Re-measured at `95575aa` (planner, 2026-10-05, go1.26.6 darwin/arm64)

Every line cite in this plan was re-read with `awk 'NR>=a&&NR<=b{printf "%d: %s\n",NR,$0}'` on the file at `95575aa`.

| Anchor | Line | Text |
|---|---|---|
| `bannedOSFuncs` | 337 | `var bannedOSFuncs = map[string]bool{"WriteFile": true, …, "NewFile": true}` |
| `doObj` lookup | 415 | `doObj := pkg.Scope().Lookup("forkLockedDo")` |
| identity check (N3 target) | 469 | `if !ok \|\| info.Uses[fid] != doObj {` |
| opener-name check (N5 target) | 488–489 | `… fo.Pkg().Path() == "os" &&` / `(fo.Name() == "OpenFile" \|\| fo.Name() == "CreateTemp") {` |
| live-scan red (fires before fixtures) | 526–527 | `if len(live.violations) > 0 { t.Fatalf("raw (non-fork-locked) file writes in verifygate tests:\n%s", …)` |
| `doDecl` | 538 | one line, trailing `\n` |
| P-Smuggle | 563–564 | last entry of `positives` |
| `positives` slice close | **565** | `	}` (design says l.566 — off by one; l.566 is `for _, p := range positives {`) |
| `not reported` | 581 | `t.Fatalf("fixture %s: %s not reported (got %v)", …)` |
| `exempt opener … reported` | 587 | |
| N-OpenClosure | 596 | |
| `writeSel["os"]` | 647 | inside `writeSel :=` l.646–649 |
| dir exemption (N6/AC7 target) | 697 | `if dir != "host/verifygate" {` |
| live tripwire red (fires before fixtures) | 717–718 | `if len(viol) > 0 { t.Fatalf("parallel tests in write+fork packages outside the fork-locked scan:\n%s", …)` |
| `body` / `pos` / prefix assertion / `neg` | 729 / 730 / 735–736 / 738 | |
| `known negative flagged` | 744 | |
| function close | 746 | `}` (last line of file) |

`grep -n 'CopyFS' host/verifygate/forklocked_write_test.go; echo rc=$?` → no output, `rc=1` (control
`grep -c NewFile` same file → `6`).

## 1. Plan overrides and corrected premises (see §7 for the controller list)

1. **The controller's `AILANG_BIN` is wrong for this repo (FALSE PREMISE, design AC8 and planner brief).**
   `/Users/voightkampff/.pinned-ailang-tools/v0.52.1/ailang` is AILANG v0.52.1. The repo pin is **v0.41.0**
   (`host/verifygate/ail_binary_gate_test.go:50–59`, `mission_config_gate_test.go:447`), served by
   `/Users/voightkampff/.pinned-ailang/ailang` (`--version` → `AILANG v0.41.0`, Commit `24ee108`).
   Measured at the pristine base with v0.52.1 exported: `go test ./... -count=1` **rc 1** — 30 `--- FAIL`
   lines in `host/coordinator`, `host/pkgproj`, `host/verifygate`, every one a pin refusal
   (`… is not the pinned v0.41.0 interpreter` / `pinned delegate … unavailable or wrong (never skip)` /
   `pinned binary unusable`; an awk sweep for a FAIL whose message is not the v0.52.1 refusal printed
   nothing); `./scripts/verify_ail.sh` **rc 1** at `World package step 9/9 … ✗ wrong compiler version: AILANG v0.52.1`.
   That red measures the binary, not the repo. Every gate in this plan uses
   `AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang`.
2. **New fixtures get their own failure texts.** The design's AC6/AC7 quote the existing texts
   (`known positive not named` / `known negative flagged`). Reusing them would make a red ambiguous
   about which fixture fired. Plan texts: `known CopyFS positive not named: %v`,
   `known host positive not named: %v`, `exempt host/verifygate fixture flagged: %v`. All measured (§4).
3. **AC7's fixture red is unobservable in the same run as its live red (FALSE PREMISE, design AC7
   "the `host/verifygate` negative fixture reds AND the live test reds").** The live `t.Fatalf`
   (l.717–718, kept first) stops the test before any fixture runs, so under the `true` mutant only the
   live red (10 lines) appears. The fixture's kill was measured separately by also neutralising the
   live assertion (`if false && len(viol) > 0`) in scratch — drill `AC7fx`. Same class as row 142's §1.6.
   N10's live escape (AC1) likewise fires the live red at l.526–527, never the CopyFS fixture.
4. **Slice anchor**: the `positives` slice closes at **l.565**, not l.566 (§0). Insert after l.564.
5. **`bodyCopyFS` is a literal const**, not a `%s` template: it inlines `\tt.Parallel()\n` so the
   `t.Parallel()` call is line 10, matching `body`'s positive. Measured: prefix `zzcopyfs/x_test.go:10` passes.
6. **AC1 escape name** is `evalCopyFS231` (this iteration), not the design's `evalCopyFS230`. Cosmetic.

## 2. Exact edits (all in `host/verifygate/forklocked_write_test.go`, cites at `95575aa`)

### M1 — `os.CopyFS` banned in both guards (kills N10)

- **l.337** → `var bannedOSFuncs = map[string]bool{"WriteFile": true, "OpenFile": true, "Create": true, "CreateTemp": true, "NewFile": true, "CopyFS": true}`
- **l.647** → `		"os":      {"WriteFile": true, "OpenFile": true, "Create": true, "CreateTemp": true, "NewFile": true, "CopyFS": true},`
- **after l.564** (P-Smuggle's `want` line; before the slice's `}` at l.565):
  ```go
  		{"CopyFS", "package p\n\nimport \"os\"\n\nfunc h() { _ = os.CopyFS(\"dst\", os.DirFS(\"src\")) }\n",
  			[]string{"fixture_CopyFS.go:5|raw os.CopyFS in func h"}, nil},
  ```
- **after l.745** (the `neg` check's closing `}`; before the function's `}` at l.746) — the tripwire positive,
  in its **own single-key map and own call** (design r3 / V16: `parallelWriteForkViolations` appends a
  violation for every `t.Parallel` of every write+fork key in the map, so a second key in `pos` would make
  `len(pv) == 2` and red the existing l.735 assertion):
  ```go

  	// Each fixture below is its own single-key map and its own call: the
  	// function reports every callsite across all keys of one map.
  	// os.CopyFS is a write (not host-prefixed, so only writeSel is exercised).
  	const bodyCopyFS = "package zz\n\nimport (\n\t\"os\"\n\t\"os/exec\"\n\t\"testing\"\n)\n\nfunc TestX(t *testing.T) {\n\tt.Parallel()\n\t_ = os.CopyFS(\"dst\", os.DirFS(\"src\"))\n\t_ = exec.Command(\"x\")\n}\n"
  	cv, _, _, err := parallelWriteForkViolations(map[string]map[string][]byte{"zzcopyfs": {"x_test.go": []byte(bodyCopyFS)}})
  	if err != nil {
  		t.Fatal(err)
  	}
  	if len(cv) != 1 || !strings.HasPrefix(cv[0], "zzcopyfs/x_test.go:10: t.Parallel in a write+fork package") {
  		t.Fatalf("known CopyFS positive not named: %v", cv)
  	}
  ```

### M2 — exemption identity and opener name pinned (kills N3, N5)

- **after M1's `CopyFS` entry** in `positives`:
  ```go
  		// P-LocalDo: a local forkLockedDo is not the package-scope wrapper, so its
  		// open closure gets no exemption.
  		{"LocalDo", "package p\n\nimport \"os\"\n\n" + doDecl + "\nfunc h() {\n\tforkLockedDo := func(open func() (*os.File, error), fill func(*os.File) error, hold func()) error { return nil }\n\t_ = forkLockedDo(func() (*os.File, error) { return os.OpenFile(\"x\", 0, 0) }, nil, nil)\n}\n",
  			[]string{"fixture_LocalDo.go:9|raw os.OpenFile in func h"}, nil},
  		// P-CreateOpener: only os.OpenFile / os.CreateTemp are exempt openers.
  		{"CreateOpener", "package p\n\nimport \"os\"\n\n" + doDecl + "\nfunc h() {\n\t_ = forkLockedDo(func() (*os.File, error) { return os.Create(\"x\") }, nil, nil)\n}\n",
  			[]string{"fixture_CreateOpener.go:8|raw os.Create in func h"}, nil},
  ```
  Line numbers (design V13) confirmed by the passing run: LocalDo `:9`, CreateOpener `:8`.

### M3 — tripwire directory exemption pinned (kills N6; pins AC7)

- **after M1's tripwire block** (same function, still before its final `}`), each its own map and call:
  ```go
  	// The exemption is exactly host/verifygate: a sibling host package is named...
  	hv, _, _, err := parallelWriteForkViolations(map[string]map[string][]byte{"host/zz": {"x_test.go": []byte(fmt.Sprintf(body, "\tt.Parallel()\n"))}})
  	if err != nil {
  		t.Fatal(err)
  	}
  	if len(hv) != 1 || !strings.HasPrefix(hv[0], "host/zz/x_test.go:10: t.Parallel in a write+fork package") {
  		t.Fatalf("known host positive not named: %v", hv)
  	}
  	// ...and host/verifygate itself is not.
  	ev, _, _, err := parallelWriteForkViolations(map[string]map[string][]byte{"host/verifygate": {"x_test.go": []byte(fmt.Sprintf(body, "\tt.Parallel()\n"))}})
  	if err != nil {
  		t.Fatal(err)
  	}
  	if len(ev) != 0 {
  		t.Fatalf("exempt host/verifygate fixture flagged: %v", ev)
  	}
  ```
- Existing `zz` positive (l.730) and no-Parallel negative (l.738) untouched. The `host/verifygate` key
  cannot overlay the real package: fixtures are `parser.ParseFile`'d from memory in a separate call
  (design V17); measured — the implemented tree's live tripwire stays at 0 violations.

After all three: `gofmt -l host/verifygate/` empty (measured); diff = 1 file, +39/−2 (measured).

## 3. Common gate vocabulary

```
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter230-forklock2
export AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang; $AILANG_BIN --version   # must print AILANG v0.41.0 (§1.1)
FENCE: go vet ./host/verifygate/                                   # go build skips _test.go
TWO:   go test ./host/verifygate/ -run 'TestVerifygateTestWritesAreForkLocked|TestNoParallelWriteForkPackagesOutsideVerifygate' -count=1 -v
TWO-R: same with -race
FIVE:  go test ./host/verifygate/ -run 'TestForkLocked|TestVerifygateTestWritesAreForkLocked|TestNoParallel|TestKernelRefuses' -count=1 -v
FIVE-R: same with -race                                           # assert no 'WARNING: DATA RACE'
FULL:  go build ./... && go test ./... -count=1
AIL:   ./scripts/verify_ail.sh
EMAIL: ./scripts/check_no_personal_email.sh                        # present in scripts/ (checked)
DRILL: H=$(shasum -a 256 <file>); apply mutant; go test ./host/verifygate/ -run '^<Test>$' -count=1 > log 2>&1; rc=$?
       (never pipe the go test whose rc you cite); grep the log for the expected text;
       git checkout -- <file>; git diff --quiet && [ "$(shasum -a 256 <file>)" = "$H" ]
```
Bank each drill log at `~/.ailang/state/w231_exec/drill_<mutant>.log` (`/tmp` is wiped).

### Baseline at pristine `95575aa` (planner-measured; a gate red here measures the repo, not the change)

| Gate | Result |
|---|---|
| FENCE | rc 0 |
| TWO | rc 0 — both PASS (3.04 s / 0.06 s) |
| TWO-R | rc 0 |
| FIVE | rc 0 — 4 PASS + `TestKernelRefusesExecOfWriterOpenFile` SKIP (darwin, by design) |
| FIVE-R | rc 0 |
| `go build ./...` | rc 0 |
| FULL, `AILANG_BIN` = v0.52.1 (controller's) | **rc 1** — 30 pin refusals only (§1.1) |
| AIL, `AILANG_BIN` = v0.52.1 | **rc 1** — step 9/9 `wrong compiler version: AILANG v0.52.1` |
| FULL, `AILANG_BIN` = v0.41.0 | **rc 0** — 26 `ok` packages, 0 `--- FAIL` |
| AIL, `AILANG_BIN` = v0.41.0 | **rc 0** — `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 365 se-tools named tests pass` |
| EMAIL | rc 0 |

Implemented tree (scratch, all three milestones): FENCE rc 0; TWO rc 0; FIVE rc 0 (4 PASS + SKIP);
FIVE-R rc 0, 0 `DATA RACE`; FULL (v0.41.0) rc 0, 26 `ok` packages, 0 `--- FAIL` (scratch is a non-worktree copy; the executor re-runs FULL in the worktree).

## 4. Test plan — mutant → killing assertion → observed failure (all run by the planner)

Each mutant was applied to the implemented scratch tree, the one named test run with `-count=1`, the
failing line copied from the log, then `git checkout -- .` + `git diff --quiet` ("reverted-clean").
Line numbers in observed text are the **implemented** file's (they shift by the inserted lines).
"Base control" = the same mutant on pristine `95575aa` (no new fixtures): rc 0 there proves the new
fixture is what kills it.

| Mutant | Edit (implemented tree) | Test run | Killing assertion | Observed (verbatim) | Base control | Status |
|---|---|---|---|---|---|---|
| N10 scan | drop `"CopyFS": true` from `bannedOSFuncs` | `^TestVerifygateTestWritesAreForkLocked$` | P-CopyFS (`not reported`) | `forklocked_write_test.go:590: fixture CopyFS: fixture_CopyFS.go:5\|raw os.CopyFS in func h not reported (got [])` | n/a (CopyFS absent at base) | MEASURED, rc 1 |
| N10 tripwire | drop `"CopyFS": true` from `writeSel["os"]` | `^TestNoParallelWriteForkPackagesOutsideVerifygate$` | `zzcopyfs` positive | `forklocked_write_test.go:765: known CopyFS positive not named: []` | n/a | MEASURED, rc 1 |
| N10 live escape | append `func evalCopyFS231() { _ = os.CopyFS("dst", os.DirFS("testdata")) }` to `ail_binary_gate_test.go` | `^TestVerifygateTestWritesAreForkLocked$` | live scan (l.526–527) | `forklocked_write_test.go:527: raw (non-fork-locked) file writes in verifygate tests:` / `ail_binary_gate_test.go:763: raw os.CopyFS in func evalCopyFS231: use a fork-locked wrapper` | rc 0 (silent — the judge's N10 reproduced) | MEASURED, rc 1 |
| N3 | l.469 → `if !ok \|\| fid.Name != "forkLockedDo" {` | `^TestVerifygateTestWritesAreForkLocked$` | P-LocalDo | `forklocked_write_test.go:590: fixture LocalDo: fixture_LocalDo.go:9\|raw os.OpenFile in func h not reported (got [])` | rc 0 (survives) | MEASURED, rc 1 |
| N5 | l.488–489 → `… fo.Pkg().Path() == "os" {` (name check dropped) | `^TestVerifygateTestWritesAreForkLocked$` | P-CreateOpener | `forklocked_write_test.go:589: fixture CreateOpener: fixture_CreateOpener.go:8\|raw os.Create in func h not reported (got [])` (589: the mutant deletes a line) | rc 0 (survives) | MEASURED, rc 1 |
| N6 | l.697 → `if !strings.HasPrefix(dir, "host/") {` | `^TestNoParallelWriteForkPackagesOutsideVerifygate$` | `host/zz` positive | `forklocked_write_test.go:773: known host positive not named: []` — the live sweep stays green under N6 (it reached the fixtures) | rc 0 (survives) | MEASURED, rc 1 |
| AC7 | l.697 → `if true {` | `^TestNoParallelWriteForkPackagesOutsideVerifygate$` | live tripwire (l.717–718) | `forklocked_write_test.go:727: parallel tests in write+fork packages outside the fork-locked scan:` + exactly 10 lines: `host/verifygate/mission_config_gate_test.go:246` and `host/verifygate/module_manifest_gate_test.go:{255,332,378,411,447,458,489,508,549}`, each `: t.Parallel in a write+fork package outside the fork-locked scan` | not run by the planner (design V14 measured the same 10-line live red at base) | MEASURED, rc 1 |
| AC7fx | AC7 + live assertion neutralised (`if false && len(viol) > 0`) — scratch-only, to observe the fixture | same | `host/verifygate` negative | `forklocked_write_test.go:781: exempt host/verifygate fixture flagged: [host/verifygate/x_test.go:10: t.Parallel in a write+fork package outside the fork-locked scan]` | n/a | MEASURED, rc 1 |

Nothing in this table is UNMEASURED. Not re-drilled here (design AC8 regression set: row-142 ACs and
the iter-228 survivors): covered by TWO/FIVE staying green on the implemented tree; the executor does
not re-run those drills (no row-142 code is touched).

## 5. Milestone gates and drills

**Gate M1**: FENCE; TWO; TWO-R. Live scan and live tripwire at zero violations; floors unchanged.
**Drills M1** (on the M1 commit): N10 scan, N10 tripwire, N10 live escape — expected texts per §4
(line numbers will differ on the M1-only commit: before M2/M3 land, `not reported` is l.583 and the
`zzcopyfs` check's Fatalf sits right after l.745; match on the text, cite the observed line).
**Gate M2**: FENCE; TWO; TWO-R. **Drills M2**: N3, N5 per §4.
**Gate M3**: FENCE; FIVE; FIVE-R; FULL; AIL; EMAIL (with v0.41.0 `AILANG_BIN`, §1.1). **Drills M3**: N6, AC7,
AC7fx per §4. Revert + byte identity after every drill; `git status --porcelain` empty before each commit
except the milestone's own diff.

## 6. Design residuals carried (no test this row)

R1 (fork inside the RLock region caught only by `-timeout`), R6 (D1 is a two-point probe sample, N15b),
R2′/R3/R4/R5 carried from row 142 — unchanged, as the design records.

## 7. Design discrepancies for the controller (the design is not edited)

1. **Gate binary (false premise)**: design AC8 and the planner brief cite
   `AILANG_BIN=/Users/voightkampff/.pinned-ailang-tools/v0.52.1/ailang`; the repo pin is v0.41.0 at
   `/Users/voightkampff/.pinned-ailang/ailang`. With v0.52.1 the full suite and `verify_ail.sh` are red at the
   pristine base (§1.1, measured).
2. **AC7 (false premise)**: the fixture red and the live red cannot both appear in one run; the live
   `Fatalf` fires first (§1.3). Fixture kill measured via AC7fx.
3. **Anchor**: `positives` closes at l.565, not l.566 (§0).
4. **Brief path**: the prior-art plan lives at `design_docs/implemented/w-verifygate-forklock-guard-residuals-sprint-plan.md`, not `planned/`.
5. New fixtures use distinct failure texts (§1.2) — a deliberate plan choice, not a design change.
