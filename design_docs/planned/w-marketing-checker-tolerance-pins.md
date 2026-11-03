# w-marketing-checker-tolerance-pins: one boundary pair per unpinned row-231 tolerance

Status: Planned (iteration 338), **r1**.
**Planner-Lane:** mission routing default
- **Item**: queue row 240 (`design_docs/world-mission.md:6312`). Sources: the row-231 judge
  (`verification/world-iter328/evaluator-sprint-r1.md`, findings 1–5) and the row-238 judge's addendum
  (`verification/world-iter336/evaluator-r1.md`, survivors 3–4).
- **Kind**: release-hygiene. It moves no clause. No `.ail` changes. **Test-only**: no production
  tolerance changes.
- **Base**: `2c51c581` (= `origin/dev`), branch `sprint/world-iter338-tolerance-pins`; go1.26.6,
  AILANG v0.41.0, tool v0.54.15 (V01).
- **Estimated**: ~0.2d, 1 executor run. One new test file (~110 lines); nothing else changes.
- **Evidence**: `design_docs/verification/world-iter338/designer-*.txt`. Every measurement ran on a
  `git archive` scratch copy (`../.scratch-iter338-designer`), never in the worktree. Each mutant was
  applied to one line and restored from `git show HEAD:<file>`, then checked with `cmp`.
  Runner: `designer-mut.sh.txt`.

## 1. Problem (measured at the base)

Line numbers in the row and in both judge reports are stale. The row says `phone.go:157`; the line is
now `:168` (V02). All line numbers below are re-derived at `2c51c581`. Each row is one line-addressed
replacement, then `go test -count=1 ./scripts/marketingcheck/` on the full package. All rows ran in
one round, with an unmutated control first (V03, `designer-mutants-base.txt`).

| ID | Site | Mutant | rc | Red tests |
|---|---|---|---|---|
| control | none | none | 0 | none |
| P-X0 | `position.go:150` background left | `posAbs(c.X-b.X) <= 2` → `<= 200` | **0** | none |
| P-Y | `position.go:150` background top | `posAbs(c.Y-b.Y) <= 2` → `<= 200` | **0** | none |
| P-R | `position.go:150` background right | `posAbs(c.X+c.W-b.X-b.W) <= 2` → `<= 200` | **0** | none |
| P-B | `position.go:150` background bottom | `posAbs(c.Y+c.H-b.Y-b.H) <= 2` → `<= 200` | **0** | none |
| P-A4 | `position.go:101` (a4) non-paragraph height | `posAbs(dh) > 2` → `posAbs(dh) > 3` | **0** | none |
| P-A3 | `position.go:129` (a3) non-free gap | `posAbs(e) <= 2` → `posAbs(e) <= 3` | **0** | none |
| P-GC | `position.go:129` | `gc >= 0 && posAbs(e) <= 2` → `posAbs(e) <= 2` | **0** | none (§2.5) |
| F-E | `phone.go:168` (i) entry | `phoneNear(s.FirstEntry, phoneR1Entry, 4)` → `…, 40)` | **0** | none |
| M-POS | `main.go:51` | `case "position":` → `case "position-unrouted":` | 1 | `RunbookMatchesHelp/position-lines` only |
| M-PH | `main.go:53` | `case "phone":` → `case "phone-unrouted":` | 1 | `RunbookMatchesHelp/phone-lines` only |
| D-EXP | `docs/QUICKSTART.md:1682` | `$DEMO_DIR/rects.json` → `design_docs/verification/world-iter312/r1-rects.json` | 1 | `rects-value` **and** `position-lines` |
| P-NEG | `position.go:150` (positive control) | `posAbs(c.Y-b.Y) <= 2` → `<= -1` | 1 | 50 of 51 `PositionCorpus` subtests, Parity, `PositionCLI/accept-rc`, `position-lines` |

- The right-edge mutant (P-R) was not in the row. It survives too.
- The four just-inside background mutants (`<= 2` → `<= 1`) also survive at the base. Only the new
  arms kill them (V08). The background is pinned in neither direction.
- Parity cannot see any of these. Both banked oracles carry the same tolerances, and the seed stores'
  backgrounds equal R1's exactly (V05).
- `TestMarketingPositionCLI/dispatch` and `TestMarketingPhoneCLI/dispatch` call `run` with no
  arguments and expect rc 2. With the case deleted, rc is still 2 (usage), so neither subtest is red
  under M-POS or M-PH (V03, V04).

## 2. Design

### 2.1 Placement: one new file, `scripts/marketingcheck/tolerance_test.go`

The arms do **not** go into `positionNewArms` or `phoneArms`. Those tables carry byte pins
(`positionNewArmSHA256`, "the banked generator's output"), fixed counts (51, 45, 63, V02), and oracle
parity. The phone table is also replayed on the real corpus. Adding arms there would move all of
those, and still pin nothing new, because the oracle has the same tolerance. Row 246 edits those
same tables in `position_test.go`. A separate file has no textual conflict with it (§5). The file
reuses the existing helpers and adds three test functions. Its reference implementation is the
prototype `designer-tolerance-pins-proto_test.go.txt`, which passed 15/15 at the base (V07). The
executor ports it with the names below.

### 2.2 Position: `TestMarketingPositionTolerancePins` (12 subtests)

- **Builder.** Read `testdata/position/stores/r1c-seed5.json`, the same source
  `writePositionNewArms` uses. `json.Unmarshal` gives a fresh `Rectangles` for each subtest. Apply
  one in-memory mutation, then call `PositionJudge(loadPositionBaseline(t), d)`.
  Seed 5's geometry equals R1's in every shot (control: `r1-seed1` differs, V05). So each offset
  below is exactly the offset from the baseline. Background edits copy the `ShotRects` value, edit
  `.Background`, and store it back in the map.
- **Assertion.** `regions == 17`, and `reasons` equals the listed slice exactly (empty means ACCEPT).
  This is stricter than `positionVerdict`, because it names the region and the clause.

| Subtest | Mutation of seed 5 | Want |
|---|---|---|
| `bg-x0+2` / `bg-x0+3` | workbench-light background `X += k; W -= k` (right edge fixed) | ACCEPT / `workbench-light.png/background: edge >2px` |
| `bg-y+2` / `bg-y+3` | background `Y += k; H -= k` (bottom fixed) | ACCEPT / same line |
| `bg-right-2` / `bg-right-3` | background `W -= k` | ACCEPT / same line |
| `bg-bottom-2` / `bg-bottom-3` | background `H -= k` | ACCEPT / same line |
| `a4-graph-svg-h+2` / `-h+3` | `graph.png` content[0] `H += k` | ACCEPT / `graph.png/graph-svg: (a4) height 277 vs 274` |
| `a3-timeline-gap+2` / `+3` | workbench-light `header.H -= 1`; `timeline-first.Y += k-1` | ACCEPT / `workbench-light.png/timeline-first: (a3) gap after header 103 vs R1 100` |

**How each arm reaches its line.**
- *Background.* Each pair moves exactly one of the four edges by k and leaves the other three at
  offset 0. So at k = 3 only that edge's term on `:150` is false. Shrinking keeps `posInside` true, so
  no `(b)` line appears. The gutter is `x=2 y=343 w=24 h=717` (V05).
- *(a4).* `graph.png` holds one content region (`graph-svg`), so no successor gap or (a5) term
  follows it. Under P-A4, `h+3` therefore turns to ACCEPT and does not just change its tag.
- *(a3).* `timeline-first` follows a non-paragraph head across a **non-free** gap (`gb = 100`).
  `header.H -= 1` (dh −1, inside (a4)) lowers the header bottom to 322. Then `Y += 1` gives
  `gc = 102` (e = +2) and `Y += 2` gives `gc = 103` (e = +3).
  - The −1 height residual is what removes the (a5) masking that the row-238 judge reported. Under
    P-A3, drift is −1 + 3 = +2 ≤ `positionDriftTol`, so the mutant ACCEPTs.
  - Without the compensation the mutant would only swap `(a3)` for `(a5)`.
  - At the base, okGap is false, so drift stays −1 and the only line is the (a3) line.

### 2.3 Phone: `TestMarketingPhoneEntryTolerancePins` (2 subtests, `entry+4`, `entry+5`)

- **Builder.** `loadPhoneSkeleton(t, testdata/phone/skeleton-r1c-seed5-pass-a-phone.json)`. This is
  the V58 scan geometry of the real seed-5 capture, which `TestMarketingPhoneSkeleton` proves
  renders back to R1. Add `dy` to `Y0`/`Y1` of `Bands[phoneEntryStart : phoneEntryStart+6]` (bands
  9–14). Then `renderPhoneSkeleton` → `writePhonePNG` (into `t.TempDir()`) → `judgePhoneFile`, the
  subcommand's decode path.
- **Reach.** `ScanPhone`'s `FirstEntry` is the union of exactly those six bands, and band 0 (the
  header) is untouched. The free rows are 30 above band 9 and 27 below band 14 (V06), so a shift of
  ≤ 5 neither merges nor reorders bands.
- **Assertion.** `sc.Header == phoneR1Header`, and `sc.FirstEntry == phoneR1Entry` shifted by `dy`.
  `dy = 4` gives `phoneVerdict(reasons) == "ACCEPT"`.
  `dy = 5` gives `"REFUSE (i)"`, and the one reason is
  `(i) header {25 45 370 68} entry {41 433 364 615} vs R1 {25 45 370 68} {41 428 364 610} ±4` (V07).
  The existing `v67/shift-5` arm moves the header as well, and the header's own tolerance refuses it.
  That is why F-E survives today.

### 2.4 Routes: `TestMarketingRunRoutesPositionPhone` (in scope)

Two asserting calls through `run`, both expecting rc 0:
- `run([]string{"position", "--baseline", positionBaseline, "a=" + <seed-5 store path>})`;
- `run([]string{"phone", "a=" + <seed-5 skeleton render written as PNG>})`.

A deleted case falls through to usage and gives rc 2. A swapped case also fails: `runPhone` rejects
`--baseline`, and `runPosition` requires it. Both give rc 2 (V08: M-POS, M-PH, M-SWAP, M-SWAP2).
- **Why in scope.** It costs ~10 lines, and it closes an S6 gap: the two existing `dispatch`
  subtests cannot fail on their own null case.
- **What it does not change.** The existing subtests stay as they are. They still pin the
  no-argument rc and `usageText`.

### 2.5 `gc >= 0` (P-GC): equivalent on the R1 baseline. No pin.

- **When the mutant can differ.** The mutant removes a conjunct, so it can only turn a refusal into
  an acceptance. That happens only when `gc < 0` and `|gc - gb| <= 2`, which needs `gb <= 1`.
- **R1's gaps.** R1's non-free gaps are `100, 100, 15, 15, 15`, minimum 15 (V05). With
  `|e| <= 2`, `gc >= 13 > 0`. Free gaps take the other branch and are out of the row.
- **Sweep measurement (V09).** The sweep moved each non-free region's top through
  `gc = -5 … gb+5`: 300 candidates.
  - It ran under the base judge and under P-GC. The 300 R1 lines are byte-identical.
  - The same run has a labelled instrument-health control: a copy of R1 with `gb = 0`. That control
    differs on exactly `gc = -2, -1`, so the sweep can see the mutant where it is not equivalent.
- **Why no pin.** `runPosition` accepts any `--baseline`. But §13's only release line passes R1,
  and `rects-value` pins that. A pin would need a synthetic baseline that the checker never sees in
  production. §2's fixture rule excludes that. P-GC is therefore recorded as **equivalent on R1**.
  It stays in §3 as an expected survivor, so a later change to R1's gaps shows up.

### 2.6 The §13 export observation: out of scope (it no longer reproduces)

The row-231 judge saw D1 (the export reverted to the R1 file) red only `rects-value`. At `2c51c581`
it reds both `rects-value` and `position-lines`. The second comes from `exportPath`'s
"outside $DEMO_DIR" guard (`docs_test.go:569`) (V03, `designer-mutants-base.txt`). Nothing is left
to do.

### 2.7 Not in scope

- Sign-loss mutants (`posAbs(dh) > 2` → `dh > 2`, and so on). The pairs here are one-sided (+k),
  and the row asks for boundary arms. Row 246 already covers the drift sign.
- The free-gap branch's `gc >= 0`.
- Any parity ledger entry. No oracle input changes.

## 3. Mutation plan (all measured in V08, `designer-mutants-proto.txt`)

Apply each mutant as **one line-addressed replacement**. `posAbs(c.X-b.X) <= 2` and
`posAbs(c.X+c.W-b.X-b.W) <= 2` also occur on `:90` ((a1)). `phone.go:168` contains `, 4)` twice. Use,
for example, `sed -i '' '150s/posAbs(c.X-b.X) <= 2/posAbs(c.X-b.X) <= 200/' scripts/marketingcheck/position.go`.
Restore after each mutant with `git checkout -- <file>`. Run the full package each time.

| # | Mutant (file:line, old → new) | Must red | Measured |
|---|---|---|---|
| T1 | `position.go:150` `posAbs(c.X-b.X) <= 2` → `<= 200` | `PositionTolerancePins/bg-x0+3` | sole killer |
| T2 | `:150` `posAbs(c.Y-b.Y) <= 2` → `<= 200` | `bg-y+3` | sole |
| T3 | `:150` `posAbs(c.X+c.W-b.X-b.W) <= 2` → `<= 200` | `bg-right-3` | sole |
| T4 | `:150` `posAbs(c.Y+c.H-b.Y-b.H) <= 2` → `<= 200` | `bg-bottom-3` | sole |
| T5–T8 | `:150` each of the four terms `<= 2` → `<= 1` | `bg-x0+2`, `bg-y+2`, `bg-right-2`, `bg-bottom-2` | each sole |
| T9 | `:101` `posAbs(dh) > 2` → `posAbs(dh) > 3` | `a4-graph-svg-h+3` | sole |
| T10 | `:101` `posAbs(dh) > 2` → `posAbs(dh) > 1` | `a4-graph-svg-h+2` | also `Corpus/new/bnd-nonpara-drift-2/-3`, Parity |
| T11 | `:129` `posAbs(e) <= 2` → `posAbs(e) <= 3` | `a3-timeline-gap+3` | sole |
| T12 | `:129` `posAbs(e) <= 2` → `posAbs(e) <= 1` | `a3-timeline-gap+2` | also `Corpus/new/perm-roster-drift-2px`, Parity |
| T13 | `phone.go:168` `phoneNear(s.FirstEntry, phoneR1Entry, 4)` → `…, 40)` | `PhoneEntryTolerancePins/entry+5` | sole |
| T14 | `phone.go:168` same → `…, 3)` | `entry+4` | also `PhoneArms/v67/shift-4`, Parity |
| T15 | `main.go:51` `case "position":` → `case "position-unrouted":` | `RunRoutesPositionPhone` | also `RunbookMatchesHelp/position-lines` |
| T16 | `main.go:53` `case "phone":` → `case "phone-unrouted":` | `RunRoutesPositionPhone` | also `phone-lines` |
| T17 | `main.go:52` `runPosition(args[1:], os.Stdout, os.Stderr)` → `runPhone(…)` | `RunRoutesPositionPhone` | also `position-lines` |
| T17b | `main.go:54` `runPhone(args[1:], os.Stdout, os.Stderr)` → `runPosition(…)` | `RunRoutesPositionPhone` | also `phone-lines` |
| T18 | `:129` `gc >= 0 && posAbs(e) <= 2` → `posAbs(e) <= 2` | **none**: expected rc 0 (§2.5) | survives, equivalent on R1 |
| C | `:150` `posAbs(c.Y-b.Y) <= 2` → `<= -1` (positive control) | 50 Corpus subtests and more | rc 1 |

"Sole" means the full-package FAIL list held only that new subtest. The control ran green (rc 0) in
the same round.

## 4. Milestone M1 (test-only) and acceptance

M1 adds `scripts/marketingcheck/tolerance_test.go` with §2.2–§2.4. Every command runs under a copy of
`world-iter336/run_bounded.sh.txt`, with these limits: package tests 600 s, vet/build 600 s,
`go test ./...` 2400 s, `verify_ail.sh` 1200 s. Export
`AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang`.

1. `go test -count=1 -v -run '^TestMarketing(PositionTolerancePins|PhoneEntryTolerancePins|RunRoutesPositionPhone)$' ./scripts/marketingcheck/`
   → rc 0; 14 `--- PASS` subtest lines (12 + 2), 3 `--- PASS: TestMarketing` top-level lines, and 0 `--- SKIP`.
2. `go test -count=1 ./scripts/marketingcheck/` → rc 0.
3. Every §3 row reproduces as listed. T18 stays rc 0. The positive control C is rc 1.
4. `gofmt -l scripts/marketingcheck` prints nothing, and `go vet ./scripts/marketingcheck/` gives rc 0.
5. `go build ./... && go test ./...` → rc 0; `./scripts/verify_ail.sh` → rc 0.
6. `git diff --name-only 2c51c581 -- . ':!design_docs'` lists exactly
   `scripts/marketingcheck/tolerance_test.go`.
   `git diff --quiet 2c51c581 -- scripts/marketingcheck/position.go scripts/marketingcheck/phone.go scripts/marketingcheck/main.go`
   → rc 0, so the production files are byte-unchanged.

A red in step 2 only at `TestMarketingCaptureFakeChromeSixShots/version-changed` with `(120.0…s)`
is the row-247 flake. It was seen once in this design's measurements (V10). Re-run once and record
it. Do not edit `capture_test.go`.

## 5. Conflict surface and residuals

- **Files.** One new test file. No `.ail`, `go.mod`, `.github/`, `tools/launchd/*`, `host/`,
  production `scripts/marketingcheck/*.go`, testdata or doc change.
- **Concurrent work (V11).** There are 0 open PRs (control: merged #464/#463).
  - Row 246 (`w-a5-drift-sign-and-rounding-pins`) adds arms and ledger entries in
    `position_test.go`. This row does not touch that file, and no test name collides.
  - Row 247 is the capture flake noted in §4.
- **R1.** The pins are relative to R1 and to the seed-5 fixtures. If R1, the seed-5 store or the
  seed-5 skeleton changes, the expected reason strings in §2.2–§2.3 change with them. The existing
  `TestMarketingPositionFixtures` and `TestMarketingPhoneSkeleton` pins make such a change loud.

## 6. Verification Log (base `2c51c581`; files under `design_docs/verification/world-iter338/`)

| V | Command | Observed |
|---|---|---|
| V01 | `git rev-parse HEAD origin/dev`; `go version`; both pinned `ailang --version` (`designer-env.txt`) | `2c51c581…` twice; `go1.26.6 darwin/arm64`; `AILANG v0.41.0`; `AILANG v0.54.15` |
| V02 | `grep -n -F` of every site, each call with a known-positive control (`designer-sites.txt`) | `position.go` `101` dh, `129` gap, `150` background; X0/right also on `:90`; `phone.go:168` (row says `:157`); `main.go` `49` docs (control), `51` position, `52` runPosition, `53` phone; dispatch subtests `position_test.go:758`, `phone_test.go:642`; `QUICKSTART.md:1682` export; `posFreeGap` `:53–55`; counts `51`/`45`/`63`/`34` at `position_test.go:550`, `phone_test.go:513`, `:771`, `position_test.go:463` |
| V03 | §1 mutant round (`designer-mutants-base.txt`) | control rc 0 (39 s); P-X0, P-Y, P-R, P-B, P-A4, P-A3, P-GC, F-E each rc 0 with `FAIL: none`; M-POS rc 1 with `RunbookMatchesHelp/position-lines` only; M-PH rc 1 with `phone-lines` only; D-EXP rc 1 with `rects-value` and `position-lines`; P-NEG rc 1 (50 Corpus subtests, Parity, `PositionCLI/accept-rc`, `position-lines`). Every mutated file was restored byte-identically (cmp) |
| V04 | the M-POS and M-PH FAIL lists in V03, compared with the subtest names in V02 | neither lists `TestMarketingPositionCLI` or `TestMarketingPhoneCLI` (their `dispatch` subtests stay green). Control: the same lists do name the docs subtests |
| V05 | geometry dump of R1 and the stores (`designer-geometry.txt`) | gutter background `x=2 y=343 w=24` in all five shots; `graph.png` has 1 content region (`graph-svg h=274`); workbench `header y=114 h=209`, `timeline-first y=423`, `gb=100`; non-free gaps `[100, 100, 15, 15, 15]`, min 15; `r1c-seed5` geometry == R1 `True`; control `r1-seed1` `False` |
| V06 | seed-5 phone skeleton dump (same file) | header band 0 `{25 45 370 68}`; entry bands 9–14 `Y 428…610`; free rows 30 above band 9 and 27 below band 14 |
| V07 | `go test -count=1 -v -run '^TestProto'` with the prototype in scratch; `gofmt -l`; `go vet`; `diff -r` against a fresh archive (`designer-proto-run.txt`) | rc 0, 15 `--- PASS`; log `entry+4: entry {41 432 364 614} ACCEPT []` and `entry+5: … REFUSE (i) ["(i) header {25 45 370 68} entry {41 433 364 615} vs R1 … ±4"]`; gofmt empty, vet rc 0; diff lists only the prototype file |
| V08 | the §3 round with the prototype present (`designer-mutants-proto.txt`) | control rc 0; T1–T9, T11 and T13: rc 1 and the only reds are the named new subtest. T10 and T12 add the existing Corpus subtests listed in §3; T14 adds `PhoneArms/v67/shift-4` and Parity; T15–T17b red `TestProtoDispatch` and the docs subtest; T18 (P-GC) rc 0 |
| V09 | `TestProtoGapSweep` (`designer-gc-sweep_test.go.txt`) under the base judge and under P-GC; `diff` (`designer-gc-sweep.txt`) | `SWEEP cases=300 min-nonfree-gb=15`; both runs rc 0; `diff` shows only `ctl-gb0 gc=-2` and `gc=-1` (base refuses with (a3), mutant `[]`); 300 R1 lines identical |
| V10 | the first, unlogged D-EXP run (scratch `out/D-EXP.log`) | an extra `--- FAIL: TestMarketingCaptureFakeChromeSixShots/version-changed (120.01s)`, `rc=-1`, package `157.334s`. The logged re-run of the same mutant (V03) does not show it. This is row 247's flake, not a D-EXP effect |
| V11 | `gh pr list --state open`; control `--state merged --limit 2`; grep of rows 246/247 (`designer-flake-and-concurrency.txt`) | no open PRs; control `#464`, `#463`; row 246 = `w-a5-drift-sign-and-rounding-pins`, row 247 = `w-capture-fakechrome-deadline-flake` |
