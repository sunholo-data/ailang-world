# w-marketing-position-checker: commit the P0/P3 position and phone predicates as CI-gated `marketingcheck` subcommands

- **Status**: r1, iteration 328 (unattended): the one mandated revision after quorum r1 on r0
  (`9abc25b5`) was BLOCKED. r1 applies fixes A, B and C (§10) and changes nothing else.
- **Item**: queue row 231 (V02). The P0/P3 position predicate (`p0check`) and the phone predicate
  (`phonecheck`) exist only as banked controller programs,
  `design_docs/verification/world-iter315/v64-p0check-main.go.txt` and `v66-phonecheck-main.go.txt`.
  On capture day they are built in a scratch directory and trusted through two sha256 pins
  (orchestration §4.4 and its V72–V74). This row makes them `marketingcheck position` and `marketingcheck
  phone`, with the V64–V67 corpora as CI fixtures. It also takes the cheap carried findings and
  files the rule changes as successors (§2.6, §7).
- **Release-hygiene** row. It moves no clause, and row 150 does not depend on it (V02). A1 has not
  run yet (V30), so the orchestration doc's §4.4 invocations switch to the binary (§2.7).
- **Estimated**: about 1 executor day in 3 milestones. Production: about +400 LOC in two new files
  (the prototype is 196 + 207 lines, V28) plus +4 lines in `main.go`. Tests: about +700 LOC.
  Fixtures: 34 JSON files and 10 skeleton JSON files, all text (§2.3). Prototype sources: V10/V11.
- **Measurement base**: `cec48f69af6da7fe4817d9a1bf3e9f01b0274cb0` (origin/dev), worktree branch
  `sprint/world-iter328-marketing-position-checker` (V01).
- **Instruments**: `go1.26.6 darwin/arm64` inside the repo (GOTOOLCHAIN auto) and `go1.26.4`
  outside any module (V01). `AILANG_BIN=$HOME/.pinned-ailang/ailang` for the doc tests (V21). The
  banked programs, the prototype port, the arm generators and every mutant were built in
  throwaway copies under `/tmp/w328d` (`git archive HEAD`), never in the worktree. Every evidence
  file is `design_docs/verification/world-iter328/vNN-*.txt`, and every program used is banked
  beside its output as `.txt`. **Row numbering**: V01–V33 are this doc's rows (§9; V31–V33 added in r1). V40 and
  V50–V74 are the orchestration doc's rows; the corpora keep their orchestration names
  (the V64–V67 corpora).

## 1. Problem

P0 item (4) and P3 read two verdicts that no CI job has ever run. The banked sources:

- are built from `.txt` copies in `/tmp` on capture day (§4.4 lines 236–237, V28);
- exit 0 whatever the verdict, so the verdict is read from a text column (V16–V17: `p0check` on a
  REFUSE arm, rc 0);
- accept vacuously in two ways. A baseline whose `graph.png` content list is empty gives
  `same r3 ACCEPT`, rc 0. An invocation with zero candidates prints `TOTAL candidates=0`, rc 0, for
  both programs (V16);
- isolate only part of their own clauses. Four p0check mutants survive the whole V65 arm set
  (`(a1)`, the free-gap `gc ≥ 0`, the 24 px minimum, non-paragraph height; 0 flips each, V12). Two
  phonecheck clauses, (ii) and `dims`, are not isolated by any V67 arm (V11 control, V14).

`marketingcheck` itself is CI-gated. Its package ran in both legs of the base commit's CI run
(`ok … scripts/marketingcheck 26.907s` plain, `157.042s` race, run 41146815448), only through
`verify_go.sh`'s `go test ./...` (no CI step names a `TestMarketing*` test; V19). So a test added to
the package is CI-gated with no workflow change. A `t.Skip` in it would not be visible, because the
leg runs without `-v` (V19). §2.4 rules that out for every CI-path test.

## 2. Decisions

### 2.1 Surface and exit-code contract

Two subcommands join the switch at `main.go:36` (V28). The names `position` and `phone` collide
with nothing in the package (V27). Neither has a measured reason to change.

```
marketingcheck position --baseline RECTS LABEL=RECTS...
marketingcheck phone LABEL=PNG...
```

- **`position`** is the r3 `check` mode of `p0check`, ported verbatim (§2.2). It loads the baseline
  with the package's `LoadRectangles`, which already refuses an empty content list, a missing shot
  or a wrong schema (`png.go:174`, V27). It reuses `Rect`/`Rectangles` and takes the image width and
  height from the `Shots` table (`png.go:49`). **[r1]** Its values equal p0check's: 1440 wide and
  1100/1100/3700/2100/3000 tall, both lists printed from the files (V32). The V10 parity
  prototype was built **in-package**. Alone in an empty directory it does not compile
  (`undefined: Rect`, `undefined: Rectangles`, rc 1). As `scripts/marketingcheck/position.go` in a
  `git archive HEAD` tree it builds (rc 0), and its `shotByName` (`:41`) reads `png.go`'s `Shots`
  (V32). So V10's parity ran against the package's real table. A geometry change that no fixture
  notices is pinned separately by `TestMarketingPositionGeometryMatchesOracle` (§3 M1, §4 P14).
- **`phone`** is the r3 verdict of `phonecheck`. The V58 scan functions are copied unchanged. The
  reference-image argument is dropped, because it fed only the withdrawn r2 field and the
  `CONTROL` line. The new real-corpus test takes over the `CONTROL` line's job: it asserts exact R1
  header and entry on all 10 real images (§2.3).
- **r2 is dropped from both.** §4.4 calls it informational, and r3 withdraws it.
- **Output**: one line per label, `LABEL ACCEPT|REFUSE …`, then one indented line per failing
  clause. The position reason strings are byte-identical to p0check's `r3:` lines (V10). The run
  ends with `TOTAL candidates=N accept=A refuse=R`. `position` also prints `regions=17` per label
  (12 content plus 5 backgrounds, V16–V17).
- **Exit codes** follow the package's existing `check`/`docs` convention (V27):
  - **0** when every label is ACCEPT;
  - **1** when any label is REFUSE, or an input is refused (`bad-rectangles`, `zero-regions`,
    `decode-failed`, named on stderr as `marketingcheck: <code> <name>: <detail>`);
  - **2** on a usage error: no candidates, a token that is not `LABEL=PATH`, a duplicate or empty
    label, or a missing `--baseline`;
  - `--help` exits 0.

  The prototype measured each case (V16–V17). Capture day reads **the exit code and the label's
  line together**, and a disagreement is BLOCKED (§2.5).

### 2.2 Behaviour preservation: a parity oracle that stays in CI

The banked sources become the **oracle**. Two CI tests copy each `.txt` into `t.TempDir()` as
`*.go` and check its sha256 against the §4.4 pin (`0a9a6795…`, `31576922…`; both re-verified,
V03). They build it with `go build -o <bin> <file>.go`. File mode needs no `go.mod` and costs
2.04 s on a cold cache and 0.34 s on a warm one (V26). Then they run oracle and subcommand on the
same files and compare per label:

- **position**: the r3 verdict **and** the r3 reason lines. On the prototype, 44 inputs (V64 11, V65
  23, and 10 new arms) gave 0 verdict differences and 0 reason differences. The same harness, given
  the prototype with `(a1)` neutered, reports 1 verdict and 5 reason differences, so it can see a
  difference (V10).
- **phone**: the r3 verdict and each clause's pass/fail. On 59 inputs (V50, 9 stores × 2 passes,
  desktop-light, 12 phoneneg, 23 judge and 4 new arms), the prototype differs in exactly one
  place, the named divergence below. The `(ii)`-dropped control shows 1 verdict and 3 clause
  differences (V11).

**The one deliberate verdict change, (6).** Phone (ii) becomes `CheckRegion` (`png.go:281`, the
M7b region gate: bounds, 64×16 minimum, std-dev ≥ 4.0, ≥ 2 colours). Before, it was the banked
stand-in, std-dev ≥ 4.0 only. Given (i), the two differ only for a header band 15 px tall (both
edges 4 px inside R1's 23 px, `MinRectH` 16; V29). The arm `header-15px` makes that case:
- the banked program ACCEPTs it;
- `phone` REFUSEs it with `(ii) rect-too-small phone.png: header` (V11, V14).

This change matches the instrument to C8′'s own words ("the M7b content-variance gate inside
those rectangles", V28), as r5's run fix did. **Re-measured**: all 9 iteration-312 stores (both
passes), V50, the seed-5 desktop control and every V67 and judge arm keep their banked verdict
(V11: 58 of 59 equal). The parity test carries a **divergence ledger**, exactly `{header-15px:
oracle ACCEPT, phone REFUSE}`. Any other difference is a red, and so is this entry if it stops
differing. A later rule change (§7 S-a) must edit the ledger, so no verdict changes silently.

### 2.3 Fixtures

**Position: JSON, checked in as produced.** These go under `scripts/marketingcheck/testdata/position/`:

| Files | Produced by (recorded in `SHA256SUMS` comments) | Pin |
|---|---|---|
| `stores/<store>.json` ×9 | the V60 block-box probe (`r1-rects-probe_test.go.txt:228`, `json.MarshalIndent`, V25), copied via `world-iter315/v64-rects-<store>.json` | sha256 equal to the store's `rects.json` line in `world-iter315/v64-inputs-sha256.txt` (9/9 equal today, V05). r1-seed1/2 are the superseded Range-box files, kept as must-REFUSE |
| `stores/r1-seed{1,2}-blockconverted.json` | banked `p0check fromlayouts <v64-layouts> <v64-rects> <out>`; cmp-identical to the banked files (V06) | sha256 of `world-iter315/v64-rects-*-blockconverted.json` |
| `v65/<arm>.json` ×19, `v65-p1/<arm>.json` ×4 | banked `p0check mutate world-iter315/v64-rects-{r1c-seed5,eval-p1}.json <dir>`. Deterministic: identical sha lists over 2 runs × 2 toolchains (V23) | the V23 sha256 lists |

The baseline is read in place from `design_docs/verification/world-iter312/r1-rects.json` (sha256
`ba7c83ce…`, V03), the file capture day uses. Git stores identical blobs once, so the
four stores equal to R1 cost nothing. The ten new arms (§2.6 (5) and the permissive pins) are
**in-test Go mutations** of `stores/r1c-seed5.json`, ported from the banked generator
`v13-isoarms-main.go.txt`. They are not files. Their banked-p0check verdicts are V13.

**Phone: why no PNG.** `verify_go.sh:240` runs `tracked_binary_gate.sh`, which classifies a blob
as binary if it contains a NUL byte. Every PNG does. The gate admits PNGs only flat under
`website/static/img/v1/` (V18). Measured on a throwaway commit of the base tree:
- control: rc 0, `0 other binary blobs among 5159 tracked files`;
- a real phone.png under `scripts/marketingcheck/testdata/phone/`: rc 1, `not-allowlisted`;
- the same PNG in a subdirectory of the allowance: rc 1, `not-allowlisted`;
- the same PNG flat in the release directory: `marketingcheck docs --assets` refuses it as
  `extra-shot` (V18).

The repo tracks 0 PNGs (V18). A base64 or run-length copy of a capture would get past the gate by
form alone. That is the same move as the classifier steering in the gate's own history
(`tracked_binary_gate.sh:13–23`, V28), so it is rejected.

**Phone: skeleton fixtures.** A skeleton is the scan's own geometry of one real capture: band
rectangles, thumb rectangle and colour, title right, page colour, plus the source sha256 and the
generating command. **[r1]** The source is recorded by file basename, so the bytes do not depend
on the host path, and the skeleton also carries `render_pix_sha256`. Each is about 2.4 KB of JSON
(V15, V31). A committed renderer (`phone_test.go`) draws
it back as an `image.RGBA`:
- page colour;
- each band as alternating ink columns, with its end columns inked;
- the header inked through `title_right`;
- the thumb filled.

Measured on the prototype (V15):
- **[r1]** all 10 renders scan to exactly their real images' header, entry, band count, thumb and
  title-right, and all 10 ACCEPT on both `phone` and the oracle (V31);
- all 39 arm generators (12 phoneneg, 23 judge, 4 new) give the **same banked verdict on the
  render as on the real image** (0 differences).

So arms derived from the render stand in for arms derived from the capture. **[r1, fix A]** M2
commits **one skeleton per distinct real image, 10 in all**:
`testdata/phone/skeleton-<corpus name>.json` for the 9 iteration-312 stores' pass A and V50, about
23.6 KB in total (V31). That puts all 10 real geometries in CI: thumb X1 runs from 324 to 340.
Seed 5 pass A (thumb X1 333) stays the base the arms are generated from. They are written by
`WORLD_PHONE_CORPUS=$HOME/.ailang/state/world-iter328-phone-corpus WORLD_PHONE_SKELETON_WRITE=1
go test -count=1 -run '^TestMarketingPhoneRealCorpus$' ./scripts/marketingcheck`. They are
checked in as written and never edited by hand.

**Transformations, each named and pinned.**
- JSON: probe `MarshalIndent` → `cp` → git → `json.Unmarshal`. The files carry 0 CR bytes and end
  in LF. There is no `.gitattributes`, `core.autocrlf` is unset, and `check-attr` reports nothing
  (V25). An eol rewrite or re-encode changes the bytes, and the sha256 pins catch it, so no
  `.gitattributes` is added (that would be a repo-wide policy file outside this row).
- PNG: Chrome PNG → `png.Decode`. All 11 corpus images decode as `*image.RGBA` (V24). Then
  scan → skeleton JSON → git → `Unmarshal` → renderer → `png.Encode` → TempDir file → the
  subcommand's `png.Decode`. The skeleton records `render_pix_sha256`, the sha256 of the rendered
  `Pix`. The test asserts that sha **and** that decode(encode(render)).Pix equals the render. It pins
  pixels, not encoder bytes. Encoder bytes did match across go1.26.4 and go1.26.6 (V09: phoneneg
  outputs 12/12 equal to the iteration-315 pins), but nothing guarantees that across a toolchain
  bump, and a byte pin would red with no verdict at stake.

**Real corpus (not in CI).** The 10 distinct real phone images (pass B equals pass A for every
store, V08) and the seed-5 desktop shot were in `/tmp/world-iter312`, which macOS wipes on boot.
This design copied them, with `cp -p`, to `~/.ailang/state/world-iter328-phone-corpus/` with a
`SHA256SUMS` file. All 11 hashes equal the iteration-315 pins (V24). M2 commits that list as
`testdata/phone/corpus-SHA256SUMS`. `TestMarketingPhoneRealCorpus` (env `WORLD_PHONE_CORPUS`)
checks every hash, then asserts:
- the 10 images ACCEPT, with header == R1 and entry == R1 exactly, and thumb X1 per store as orchestration V66
  lists;
- desktop REFUSE;
- every arm regenerated from the real seed 5 has its V15 column-1 verdict, except `header-15px`;
- oracle parity with the same one-entry ledger;
- **[r1]** a fresh skeleton derived from each of the 10 real images is byte-equal to the committed
  skeleton of the same name. Exactly 10 comparisons must run. On the prototype this check gives
  `CHECK images=10 byte-equal=10 geometry-equal=10`, rc 0. Each of three controls makes it fail:
  one edited skeleton byte (rc 1, `byte-equal=9`), a missing corpus directory (rc 1, `BLOCKED`)
  and one changed `SHA256SUMS` line (rc 1, `hash=false` on that image) (V31).

Without the env the test skips. That is the one CI-path skip, and §2.5 makes the capture-day step
run it.

### 2.4 Anti-vacuity

- A baseline with any shot's content empty is refused by `LoadRectangles` (rc 1,
  `bad-rectangles graph.png: content rectangles required`). The banked program ACCEPTs that case
  (V16).
- A candidate with the right schema but no content rectangle in the five 1440 px shots gets rc 1,
  `zero-regions` (V16).
- A wrong-schema or unparsable candidate gets rc 1, `bad-rectangles`.
- Zero candidates gets rc 2. Both banked programs give rc 0 here (V16).
- `position` counts the regions it compared and refuses a count of 0 with `zero-regions`. That
  guard is the backstop if `LoadRectangles` is ever loosened.
- `phone`: a non-PNG gets rc 1 `decode-failed`. Zero bands is a named REFUSE reason
  (`zero-bands: no ink rows read`, rc 1) on the blank arm, besides (i)–(iii′) (V16).
- No CI-path test calls `t.Skip`. The only skip is `TestMarketingPhoneRealCorpus` without its env.

### 2.5 Capture-day control step (replaces §4.4's source pins and orchestration V72's controls)

Run from the repository root at the capture commit. Record that commit's green CI run id. Any
deviation below is BLOCKED and escalated before a capture-store verdict is read. No step is piped
(a pipe hides the exit code).
1. `WORLD_PHONE_CORPUS=$HOME/.ailang/state/world-iter328-phone-corpus go test -count=1 -timeout 20m -v -run '^TestMarketing(Position|Phone)' ./scripts/marketingcheck > "$LOG" 2>&1; echo rc=$?`
   gives `rc=0`. Then `grep -c '^--- SKIP' "$LOG"` gives `0`, and `grep -c '^--- PASS: TestMarketingPhoneRealCorpus' "$LOG"` gives `1`.
   **[r1, fix C]** The corpus is host-local, unversioned disk state. If the directory is missing,
   or any file's sha256 differs from `testdata/phone/corpus-SHA256SUMS`, the step is BLOCKED and
   escalated to Mark; it is never skipped.
   - The only rebuild is `cp -p`: the 9 store images from
     `/tmp/world-iter312/<store>/capture/pass-a/phone.png` (volatile, wiped on boot) and V50 from
     `~/.ailang/state/world-iter305/probe-seeded/pass-a/phone.png` (durable). The result must
     match the committed hashes.
   - No other durable copy of the iteration-312 images exists (V33).
   - A fresh seed produces different bytes (D-WORLD-75), so a lost corpus cannot be regenerated.
   - With the corpus lost, CI's real-geometry coverage is the 10 committed skeletons alone.
     Whether to capture on that basis is Mark's call.
2. **Position accepts**: `position --baseline <R1> accept-seed5=…/stores/r1c-seed5.json
   accept-eval-p1=…/stores/eval-p1.json` gives rc 0 and two ACCEPT lines. eval-p1 is a real
   store that differs from R1 (`roster-proven` dh −24, orchestration V73), so this control is **non-identity**.
   r2's identity-like rule refuses it (V05).
3. **Position refuse**: `position --baseline <R1> refuse=…/v65/neg-attested+14.json` gives rc 1 and
   `refuse REFUSE` with `packet.png/roster-attested: (a3) gap after roster-tested 29 vs R1 15` (V16–V17).
4. **Phone accepts**: `phone accept-v50=<V50> accept-seed2=<corpus>/r1-seed2-pass-a-phone.png`
   gives rc 0. Seed 2's thumb X1 is 340 against V50's 329. No image is a reference any more, so
   neither accept is an identity comparison.
5. **Phone refuse**: `phone refuse=<capture day's pass-a/workbench-light.png>` gives rc 1.

### 2.6 Carried findings (1)–(7) and NB2

| # | Finding (row 231, V02) | Decision | Pin (test row) and the mutant it kills |
|---|---|---|---|
| (1) | inside the two free gaps, `approvals-chain` ±48 px and `validation` +216 px into blank page are accepted | **Deferred** to S-a (§7). Tightening it amends §4.4 (a3)/(a4), a rule of a quorum-carved doc that needs its own round and a re-measure of all 9 stores. Here the current verdicts are **pinned**: `perm-approvals-y+48`, `perm-approvals-y-48`, `perm-validation-h+216` ACCEPT on both the oracle and the prototype (V13, V10) | `TestMarketingPositionCorpus/perm-*` turns red on any silent tightening; S-a must flip them deliberately and update the ledger |
| (2) | compounding +2 px per gap | **Deferred** to S-a, same reason. Pinned by `perm-roster-drift-2px`: each roster gap +2 and height +2 puts PROVEN 14 px below R1. ACCEPT (V13) | as (1) |
| (3) | both capture-day accept controls are identity | **In scope**: §2.5 steps 2 and 4 | an identity-only position checker (r2's rule) refuses `accept-eval-p1` (V05), so the control step is BLOCKED. In CI, `TestMarketingPositionCorpus/eval-p1` kills the same mutant |
| (4) | P0 item (1)'s probe logs `repeat_equal`/`script_tags`, and there is no command | **In scope, text**: §4.4 item (1) gets orchestration V40's command and a count-asserted reading (§2.7). **Deferred** to S-b: making the probe self-asserting and committing it. It mints its own store and needs Chrome, and it lives in `cmd/ailang-worldd`, a different conflict surface | reading rule measured on the banked orchestration-V40 output: 4 pass lines, 0 fail lines, and the DEFAULT control shows 4 hits (V22). A probe that logs `repeat_equal=false` fails the `= 0` count |
| (5) | V65 isolates none of (a1), free-gap `gc ≥ 0`, the 24 px minimum or non-paragraph height | **In scope**: 5 isolating arms plus 1 boundary arm, each failing exactly one clause on the banked oracle (V13). Phone gains `header-solid-ink` ((ii) alone), `tall-845` (dims alone) and `header-15px` (§2.2), because V67 isolates neither (ii) nor dims (V14) | §4: P1–P4 survive all of V65 (0 flips) and are killed only by their arm. F1–F3 likewise (V12, V14) |
| (6) | (ii) is a std-dev stand-in | **In scope**: (ii) := `CheckRegion` (§2.2) | F3 (stand-in restored) is killed only by `header-15px` (V14) |
| (7) | a genuine thumb split by one off-range pixel left of its midpoint is refused | **Deferred** to S-a. Bridging gaps rewrites (iii′)'s "run" (a rule change), and the 19 real images are clean runs (V08). Pinned: `B1-thumb-split-px150` and `B1c-…px8` REFUSE, `B1b-…px250` ACCEPT (V15) | `TestMarketingPhoneArms/judge/B1-*`. A silent gap-bridging change flips them |
| NB2 | nothing pins WHICH file `WORLD_MARKETING_RECTS` names in §13. Reverting it to R1 keeps both doc tests green | **In scope**: reproduced first-party. With the revert, all 19 subtests of the two tests still pass, while the `--rectangls` control reds (V21). New subtests in §3 M3 | D1: the revert reds `rects-value` and `position-lines` |

### 2.7 Exact document changes (M3)

**Orchestration doc** (`design_docs/planned/w-marketing-capture-orchestration.md`). Line numbers are
at base (V28).
- **:220** (P0 item (1)): append the command. It runs in a scratch copy (`git archive HEAD`),
  never the worktree: `cp design_docs/verification/world-iter305/iter305-seeded-still-chrome-probe_test.go.txt cmd/ailang-worldd/zz_iter305_probe_test.go && WORLD_CHROME=<Chrome> PROBE_OUT=<dir> go test -count=1 -timeout 540s -v -run '^TestIter305SeededStillCaptureProbe$' ./cmd/ailang-worldd > <log> 2>&1`.
  It passes only if all of these hold: rc 0; one `--- PASS: TestIter305SeededStillCaptureProbe`;
  `grep -cE 'STILL .* repeat_equal=true script_tags=0 '` = 4; and
  `grep -cE 'STILL .* (repeat_equal=false|script_tags=[1-9])'` = 0 (V22).
- **:223**: "`p0check` and `phonecheck`" → "`marketingcheck position` and `marketingcheck phone`".
- **:235**, the final sentence only: "Exact controller invocations (r4 … V71). Both checkers are built
  from their banked source in a scratch directory outside the repo, never inside the worktree:" →
  "Exact controller invocations (row 231; r4/V71 before it). Both checks are subcommands of the
  committed `scripts/marketingcheck`, run from the repository root at a commit whose CI run is
  green:".
- **:236** → "P0 (a): `go run ./scripts/marketingcheck position --baseline design_docs/verification/world-iter312/r1-rects.json <label>=<store-rects.json>` … The verdict is `<label> ACCEPT` or `<label> REFUSE` with one indented line per failing clause, read together with the exit code (§2.1 of row 231's design). r2 is no longer printed."
- **:237** → "Phone: `go run ./scripts/marketingcheck phone <label>=<capture pass-a/phone.png>`; (i), (ii), (iii′), with (ii) = `CheckRegion`. No reference image is an argument."
- **:238** (source pins and control step) → §2.5's five steps. The V50 and seed-5 hashes stay, in full.
- **:239** ("Why the committed checker is not extended") → "**Committed (row 231).** The banked
  `.txt` programs remain the parity oracle that CI builds from the pinned sources; a change to
  either source, or any verdict difference outside the one-entry ledger, is a red."
- **:248** (P3) → the two `go run ./scripts/marketingcheck …` lines with `release=`.
- **§7: no change.** §7 names no checker and no `.txt` (0 hits against 7 in §4.4, V28). C8′ (ii)
  already reads "the M7b content-variance gate inside those rectangles", which is now literally
  what runs.
- **§9**: append "### Row 231": the invocations switch to the committed binary, the rules are
  unchanged except (6)'s named header-15px divergence, and the evidence is this doc's V-rows.

**QUICKSTART §13** (`docs/QUICKSTART.md`, heading :1663 unchanged; `CheckDocs` refuses its absence, `main.go:853`, V28).
- After the inputs block (:1678–1683), add a P3 paragraph and a bash block with three
  `go run ./scripts/marketingcheck position --baseline design_docs/verification/world-iter312/r1-rects.json …`
  lines:
  - `accept-seed5=`/`accept-eval-p1=` on the testdata stores, which must give rc 0;
  - `refuse=…/v65/neg-attested+14.json`, which must give rc 1;
  - `release="$WORLD_MARKETING_RECTS"`.
- After the capture block (:1697–1703), add `go run ./scripts/marketingcheck phone
  refuse="$CAPTURE_OUT/pass-a/workbench-light.png"` (rc 1) and `… phone
  release="$CAPTURE_OUT/pass-a/phone.png"` (rc 0). The phone accept controls use images held
  outside the repo, so they stay in orchestration §4.4 only, and §13 points there.
- :1720–1723 ("Use the sha256-pinned `.txt` checkers … exit 0 either way.") → "Both checks are
  `marketingcheck` subcommands, tested in CI against the iteration-315 corpora. Each exits 0 only
  when every label is ACCEPT, 1 when any label is REFUSE or an input is refused, and 2 on a usage
  error; read the exit code and the label's line together."
- `usageText` (`main.go:22`) gains the two usage lines. That keeps the existing
  `marketingcheck-lines` subtest, which needs a usage line containing `marketingcheck <sub> ` and
  every documented flag (`docs_test.go:438–456`, V27), green on the new lines.

## 3. Milestones

Every milestone runs `gofmt -l scripts/marketingcheck` (which must print nothing), `go vet
./scripts/marketingcheck`, and `go test -run '^$' ./...` (`go build` does not compile test
files). Each test row's sole killer is in §4.

**M1: `position`** (`position.go`, `position_test.go`, `testdata/position/**`, `main.go`).
- `TestMarketingPositionFixtures`:
  - every testdata file's sha256 equals `SHA256SUMS`;
  - each store's sha256 appears in `world-iter315/v64-inputs-sha256.txt`;
  - the block-converted files match the banked files;
  - the baseline is `ba7c83ce…`.
- `TestMarketingPositionCorpus`: one subtest per label, giving the verdict and the failing clause's
  tag (`(a1)`…`(b)`, `missing`, `background`):
  - V64: the two Range-box stores REFUSE, the 9 others ACCEPT;
  - V65: 18 REFUSE and 5 ACCEPT;
  - new: the 5 iso arms REFUSE, `acc-tested-w-2` and the 4 `perm-*` arms ACCEPT;
  - 44 rows, 19 ACCEPT (V10).
- `TestMarketingPositionParity`: the oracle (§2.2) over the same 44 inputs, with an empty ledger. It
  asserts it compared 44 labels.
- `TestMarketingPositionCLI`: exit codes 0/1/2 and every §2.4 position case.
- **[r1, fix B]** `TestMarketingPositionGeometryMatchesOracle`. It reads the oracle's `heights`,
  `order` and `inside()` width from `world-iter315/v64-p0check-main.go.txt` and asserts that
  `positionShots` equals `order` and that each shot's `Shots` width × height equals the oracle's.
  It fails if any of the three is not found. Prototype: `v32-position-geom-test.go.txt`, green
  unmutated (V32).
- Acceptance: `go test -count=1 -timeout 20m -v -run '^TestMarketingPosition' ./scripts/marketingcheck`
  gives 5 `--- PASS` lines, 0 `--- SKIP` and 0 `--- FAIL`. `go run ./scripts/marketingcheck
  position --help` exits 0.

**M2: `phone`** (`phone.go`, `phone_test.go`, `testdata/phone/**`).
- `TestMarketingPhoneSkeleton`: for each of the **10** skeletons, the render's scan equals the
  recorded geometry, `render_pix_sha256` matches, the PNG round trip is lossless, and the render
  ACCEPTs (V31).
- `TestMarketingPhoneArms`: in-test Go ports of the 12 phoneneg, 23 judge and 4 new generators on
  the seed-5 render, plus both renders. 41 rows, 12 ACCEPT (V14 baseline; per-arm list V15 column
  3).
- `TestMarketingPhoneParity`: the oracle over the same 41 PNGs written to TempDir, with ledger
  `{header-15px}`.
- `TestMarketingPhoneCLI`: §2.4's phone cases.
- `TestMarketingPhoneRealCorpus` (§2.3).
- Acceptance:
  - CI form: `go test -count=1 -timeout 20m -v -run '^TestMarketingPhone' ./scripts/marketingcheck`
    gives 4 PASS and exactly 1 SKIP (`TestMarketingPhoneRealCorpus`).
  - With `WORLD_PHONE_CORPUS=$HOME/.ailang/state/world-iter328-phone-corpus`: 5 PASS and 0 SKIP.
    The executor banks this output.

**M3: documents and NB2** (`docs/QUICKSTART.md`, orchestration doc, `docs_test.go`).
`TestMarketingRunbookMatchesHelp` gains three subtests:
- `rects-value` checks four things:
  - §13 has exactly one `export WORLD_MARKETING_RECTS=` command, and its value is
    `$DEMO_DIR/rects.json`;
  - the capture line passes `--rectangles "$WORLD_MARKETING_RECTS"`;
  - the `release=` position line passes `release="$WORLD_MARKETING_RECTS"`;
  - that line's `--baseline` is the R1 path, which is not the candidate.
- `position-lines` executes every documented `position` line from the repo root (`cmd.Dir =
  "../.."`, a binary built once):
  - `accept*` labels must give ACCEPT and `refuse*` labels REFUSE, with the matching rc;
  - the `release` line runs twice, with the export's expanded path holding `stores/eval-p1.json`
    (rc 0, ACCEPT) and then `v65/neg-attested+14.json` (rc 1, REFUSE).
- `phone-lines` executes the two `$CAPTURE_OUT` lines, with the seed-5 render (ACCEPT) and
  `tall-845` (REFUSE) written to `$CAPTURE_OUT/pass-a/`.

Acceptance: `go test -count=1 -timeout 20m -v -run 'TestMarketingRunbookMatchesHelp|TestMarketingDocsLinks' ./scripts/marketingcheck`
gives 0 FAIL, with the three new subtests listed as PASS. The executor runs §13's new lines
verbatim (S7) and banks the output.

**Repo gate (CI runs both; nothing lands red):**
- `./scripts/verify_ail.sh`;
- `go build ./... && go test -count=1 -timeout 20m ./...`;
- `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh`. Its race leg is
  `-timeout 8m` per package (`verify_go.sh:355`). The package took 157 s under race on CI (V19),
  and the new phone work adds about 3 s under race (V26).

## 4. Mutation table

Each mutant is applied alone. The killer is the assertion that fires. The "Kills" column names
every red subtest of that test function; where one subtest is named alone, it is the sole killer.
The parity test of the same milestone also reds wherever an arm's verdict moves, because it runs
the same arm.

| ID | Mutant | Kills (expected red) | Evidence |
|---|---|---|---|
| P1 | `(a1)` check removed | `TestMarketingPositionCorpus/iso-a1-tested-w-3`: want REFUSE (a1), got ACCEPT | V12: 0 V65 flips, 1 new |
| P2 | `c.H >= 24` removed | `…/iso-minh-approvals-h-1` | V12 |
| P3 | non-paragraph `abs(dh) > 2` removed | `…/iso-nonpara-graph-svg-h+24` and `…/iso-nonpara-payload-h+24` (the clause's two arms) | V12 |
| P4 | free-gap `gc >= 0` removed | `…/iso-freegap-validation-h+240` (gap −1 vs R1 239) | V12 |
| P5 | `(a2)` removed | `…/v65/neg-packet-all+14`, `neg-payload+14`, `v65-p1/neg-packet-all+14` | V12 |
| P6 | `(a3)` removed | 6 V65 arms plus `iso-freegap` | V12 |
| P7 | `(b)` removed | `…/v65/neg-proven-grows-6-lines-below-fold` | V12 |
| P8 | background edge check removed | `…/v65/neg-background-over-content` | V12 |
| P9 | baseline loaded with plain `json.Unmarshal` (no `LoadRectangles`) | `TestMarketingPositionCLI/baseline-empty-content`: want rc 1 `bad-rectangles`, got rc 0 ACCEPT | V16 (banked behaviour) |
| P10 | zero-candidate guard removed | `TestMarketingPositionCLI/no-candidates`: want rc 2 | V16 |
| P11 | `return 1` on REFUSE → `return 0` | `TestMarketingPositionCLI/refuse-rc` | V17 |
| P12 | one reason string changed (`(a3) gap after` → `(a3) gap`) | `TestMarketingPositionParity`: reason difference. Corpus tags still match, so only parity reds | V10 control |
| P13 | one byte appended to `testdata/position/v65/neg-attested+14.json` | `TestMarketingPositionFixtures`: sha256 mismatch | — |
| P14 | **[r1]** `Shots` `graph.png` 2100 → 2101 (verdict-neutral) | `TestMarketingPositionGeometryMatchesOracle`: `graph.png: Shots 1440x2101, oracle 1440x2100`. It is the sole killer among the position tests: parity over the 44 inputs shows 0 verdict and 0 reason differences. Outside them, the existing `TestMarketingShotGeometry/graph.png` and `TestMarketingPNGContent/canonical-shots` also red, but they pin `Shots` to their own tables, not to the oracle | V32 |
| P15 | **[r1]** `Shots` `packet.png` 3700 → 3660 | `TestMarketingPositionGeometryMatchesOracle`: `packet.png: Shots 1440x3660, oracle 1440x3700`. Not verdict-neutral: parity also reds (19 verdict and 43 reason differences, via `packet.png/background: (b)`) | V32 |
| F1 | (ii) removed | `TestMarketingPhoneArms/new/header-solid-ink` (also `header-15px`) | V14 |
| F2 | dims removed | `…/new/tall-845` | V14 |
| F3 | (ii) back to the std-dev stand-in | `…/new/header-15px` only; also the parity ledger entry stops differing | V14, V11 |
| F4 | `X0 <= 8` removed | `…/judge/thumb-x9` (+4 others) | V14 |
| F5 | thumb-present removed | `…/v67/no-scrollbar` (+5) | V14 |
| F6 | `X1 <= 384` removed | `…/v67/full-thumb` (+2) | V14 |
| F7 | clip removed | `…/v67/title-unclipped` (+3) | V14 |
| F8 | (i) removed | `…/v67/shift-5` (+3) | V14 |
| F9 | longest run → first-to-last span (pre-r5) | `…/judge/scattered-pixels-no-thumb` (+2 split arms) | V14 |
| F10 | renderer stripe period 2 → 3 | `TestMarketingPhoneSkeleton`: `render_pix_sha256` mismatch | — |
| F11 | `zero-bands` reason removed | `TestMarketingPhoneCLI/blank`: reason missing | V16 |
| D1 | §13 export reverted to the R1 file (NB2) | `TestMarketingRunbookMatchesHelp/rects-value`; also `position-lines` (the second `release` run ACCEPTs) | V21 (survives today) |
| D2 | the §13 `release=` line names the R1 path | `…/rects-value` (candidate equals baseline) | — |
| D3 | the §13 refuse control's path mistyped | `…/position-lines`: rc 1 with `bad-rectangles`, no `refuse REFUSE` line | — |

The rows with "—" are mechanical. The executor measures them in the sprint.

## 5. Acceptance criteria

- **AC-1 (parity).** Both parity tests pass in CI, comparing 44 position labels and 41 phone labels.
  The phone ledger holds exactly `header-15px`. With the real corpus, phone parity also passes over
  the 59 V11 inputs with the same ledger.
- **AC-2 (verdicts unchanged on real stores).** With `WORLD_PHONE_CORPUS`, the real-corpus test
  shows 10/10 real phone images ACCEPT. In CI, the 7 block-box stores and the 2 block-converted files ACCEPT,
  and the 2 Range-box files REFUSE. Both match the banked V64/V66 outputs (V05, V08).
- **AC-3 (non-vacuous).** §2.4's six cases give their stated rc and code. A grep of the new test
  files for `t.Skip` finds only `TestMarketingPhoneRealCorpus`. That grep is labelled
  instrument-health only.
- **AC-4 (mutations).** Every §4 row is red with its mutant applied, and green when restored. The
  executor runs them in a throwaway copy.
- **AC-5 (docs).** M3's subtests pass, D1–D3 red, and §2.7's lines changed as written.
  `git diff --stat` touches only §6's files.
- **AC-6 (repo gate).** §3's gate lines, plus CI green on the PR.

## 6. Files touched, conflict surface

| File | Change |
|---|---|
| `scripts/marketingcheck/main.go` | two `case`s, two `usageText` lines |
| `scripts/marketingcheck/position.go`, `phone.go` | new |
| `scripts/marketingcheck/position_test.go`, `phone_test.go` | new |
| `scripts/marketingcheck/docs_test.go` | 3 subtests in `TestMarketingRunbookMatchesHelp` |
| `scripts/marketingcheck/testdata/position/**` | 34 JSON files and `SHA256SUMS` (new directory; none exists today, V27) |
| `scripts/marketingcheck/testdata/phone/**` | 10 skeleton JSON files (one per distinct real image) and `corpus-SHA256SUMS` |
| `docs/QUICKSTART.md` | §13 only |
| `design_docs/planned/w-marketing-capture-orchestration.md` | §4.4 :220/:223/:235/:236–239/:248 and §9 |

Open PRs at measurement: none. `gh pr list --state open` gives `[]`, and the merged control lists
#439 and #438 (V20). Nothing touches `go.mod`, `tools/launchd/*`, `.github/` or the binary gate.

**Standards.**
- **S1/S2/S4/S5**: n/a. No `.ail` is written, and the tool is Go at the host/tooling boundary that
  checks Go-host artifacts.
- **S3**: `scripts/` is neither `world/` nor `host/`. The predicates check release-image geometry,
  not World semantics, so there is no package to publish.
- **S6**: §2.4.
- **S7**: `--help`, the §13 runnable lines, and the fixtures as examples.

## 7. Successors (named, filed by the controller)

- **S-a `w-marketing-predicate-tolerances`**, for (1), (2) and (7). It amends §4.4 (a3)/(a4)
  (bound the free-gap line steps and the cumulative drift) and the (iii′) run definition (bridge a
  1 px gap, keep the scattered-pixel refusal). It needs the orchestration doc's own quorum round,
  a re-measure of all 9 stores and the R1 seeds, and a ledger edit in both parity tests. The
  `perm-*` and `B1*` pins mark exactly what it may flip.
- **S-b `w-marketing-still-probe-commit`**, the second half of (4). Commit the iter305 still probe
  as a `WORLD_CHROME`-gated test in `cmd/ailang-worldd` that `t.Errorf`s on `repeat_equal=false`
  or `script_tags>0`. It would also run item (1) on the same store as items (2)–(3) (iteration-315
  judge NB4).

## 8. Residuals (unmeasured or not covered)

- **CI never sees a real Chrome phone image.** The 10 skeletons stand in. Their fidelity is
  measured on 39 arms (V15) and 10 images (V31), not proven in general. The real-corpus test covers the real images
  only when the capture-day step runs it.
- **`verify_go.sh` runs without `-v`**, so a future `t.Skip` in these tests would pass unseen
  (V19). AC-3 is the only guard, and it is a grep. A named-PASS manifest for the package (like CI's
  `-v` steps for other packages) is not added here.
- **[r1, fix C] The real corpus is host-local.** It lives only in
  `~/.ailang/state/world-iter328-phone-corpus/`, unversioned disk state on this host. CI never
  reads it: CI's coverage of real geometry is the 10 committed skeletons. If the directory is lost
  or a hash differs, capture-day step 1 is BLOCKED and escalated (§2.5). It can be rebuilt only
  while `/tmp/world-iter312` survives, and no durable copy of those images exists (V33). This row
  adds no new storage mechanism.
- **Untested paths.** The `perm-*` pins record permissiveness and do not judge it. P0 still never
  measures a stack-head move on a real store (orchestration §9 r3 residual, unchanged).

## 9. Verification Log

Evidence files are `design_docs/verification/world-iter328/`. "(+ctl)" marks a paired positive or
negative control in the same file.

| # | Claim | Command | Observed |
|---|---|---|---|
| V01 | base, branch, toolchains | `git rev-parse HEAD`; `git branch --show-current`; `go version` in the repo copy and in `/tmp` (`v01-base-toolchain.txt`; V01–V09 written by `v01-v09-battery.sh.txt`) | `cec48f69af6d…`; `sprint/world-iter328-marketing-position-checker`; `go1.26.6` (repo), `go1.26.4` (`/tmp`) |
| V02 | the row, its carried items, the judge findings, NB2 | `grep -n '^231\. ' design_docs/world-mission.md`; greps over `world-iter315/evaluator-r{1,2}.md`, `world-iter316/evaluator-r*.md` (`v02-row-and-findings.txt`) | `6300:231. … w-marketing-position-checker`, items (1)–(7) and NB2 as quoted in §2.6. The judge rows hit: `approvals-chain` y+48 / `validation` h+216 ACCEPT, the cumulative drift ACCEPT, "four … mutants … SURVIVE", "introduces a fragility", NB2 "SURVIVED" |
| V03 | §4.4 pins re-verified | `shasum -a 256` of both sources, the seed-5 rects, R1 and V50; `grep -o <hash> orchestration \| uniq -c` (`v03-source-pins.txt`) | `0a9a6795…927c`, `31576922…0b4c`, `ba7c83ce…3333` (seed 5 = R1), `999125ff…94be8`. Each full hash appears once in the doc. (+ctl) the V74 tamper hash in full: 0 hits |
| V04 | the banked programs build from `.txt` | copy to `main.go`, `go mod init`, `go build` ×5; `cmp` to the `.txt`; go1.26.6 file-mode build (`v04-banked-builds.txt`) | 5 binaries; every source `cmp`-identical; file mode `rc=0` |
| V05 | V64 reproduced | banked `p0check check r1-rects.json` over the 11 banked candidates; `diff` to `v64-p0-positive.txt`; sha of banked vs `/tmp/world-iter312/*/rects.json` (`v05-v07-p0-reproduction.txt`) | BYTE-IDENTICAL; r3 ACCEPT 9/9 (Range-box ×2 REFUSE); `eval-p1 r3 ACCEPT r2 REFUSE`; 9/9 banked == originals |
| V06 | block conversion | `p0check fromlayouts` on banked layouts, `cmp`; `convcheck` on 7 stores (same file) | both `cmp`-identical; `CONVCHECK rows=28 mismatches=0` |
| V07 | V65 reproduced | `p0check mutate` (seed 5, eval-p1) then `check`; `diff` to `v65-p0-negative.txt` (same file) | BYTE-IDENTICAL; 18 REFUSE, 5 ACCEPT; `TOTAL candidates=23 r3-accept=5` |
| V08 | V66 reproduced; PNG pins | banked `phonecheck <V50> …` over V50 + 9 × 2; `diff`; sha of all 18 PNGs vs `v64-inputs-sha256.txt` (`v08-v09-phone-reproduction.txt`) | BYTE-IDENTICAL; 19 `r3 ACCEPT`; 18/18 pinned (+ctl all-zero hash: 0). Pass A == pass B per store |
| V09 | V67 reproduced; encoder stability | `phoneneg` and the v67b judge generator on seed 5 pass A; sha vs `v67-synthetic-sha256.txt`; `phonecheck` over the 14; `diff` (same file) | 12/12 phoneneg sha equal (go1.26.4 build against the iteration-315 pins) and the scattered arm `c57b01b1…` equal; BYTE-IDENTICAL; `TOTAL candidates=14 r3-accept=2` |
| V10 | position parity | `parity_position.sh` (`v10-parity-position.sh.txt`) on prototype `v10-prototype-position.go.txt` (`v10-position-parity.txt`) | `PARITY position: inputs=44 banked-accept=19 verdict-diffs=0 reason-diffs=0`. (+ctl) `POS_MUT=a1`: `verdict-diffs=1 reason-diffs=5` |
| V11 | phone parity; the one divergence | `parity_phone.sh` (`v11-parity-phone.sh.txt`) on `v11-prototype-phone.go.txt`; arm `v11-h15-main.go.txt` (`v11-phone-parity.txt`) | `inputs=59 banked-accept=30 verdict-diffs=1`: only `new/header-15px ACCEPT→REFUSE (ii) rect-too-small`; `clause-level diffs=1`. (+ctl) `PH_MUT=ii`: `header-solid-ink REFUSE ACCEPT`, 3 clause diffs |
| V12 | position mutation matrix | `posmatrix.sh` (`v12-posmatrix.sh.txt`, `v12-position-mutation-matrix.txt`) | a1, minh, nonpara, freegc: `flips(v65+v65p1)=0`, each killed only by its new arm; a2 3, a3 6 (+iso-freegap), b 1, bg 1 |
| V13 | new P0 arms on the banked oracle | `isoarms v64-rects-r1c-seed5.json` (`v13-isoarms-main.go.txt`); `p0check check` (`v13-new-p0-arms.txt`) | each iso arm fails exactly one clause (`(a1) 106..1316 vs 106..1319`; `(a3) gap after validation -1 vs R1 239`; `(a4) height 1 vs 25 not whole lines`; `(a4) height 298 vs 274`; `(a4) height 84 vs 60`); `acc-tested-w-2` and the 4 `perm-*` arms `r3 ACCEPT`; drift arm offsets `roster-tested dy 2 … roster-proven dy 14`, `dh 2` each |
| V14 | phone mutation matrix (render arms) | `phmatrix.sh` (`v14-phmatrix.sh.txt`) with arms from `v14-newarms-main.go.txt` (`v14-phone-mutation-matrix.txt`) | baseline `12 accept / 29 refuse of 41`. `ii`: header-15px, header-solid-ink; `standin`: header-15px only; `dims`: tall-845 only; x0/present/x1/clip/i/span each ≥1 kill |
| V15 | skeleton fidelity (**re-run in r1** on the final skeleton format) | `skel derive/render` (`v31-skel-main.go.txt`; the r0 program `v15-skel-main.go.txt` is kept as history); prototype on renders and real images; arms regenerated from the r1 seed-5 render; `skelcmp.sh` (`v15-skelcmp.sh.txt`, `v15-phone-skeleton.txt`) | renders: header `{25 45 370 68}`, entry `{41 428 364 610}`, thumb `{6 833 333 841}` / `{6 833 340 841}`, titleRight 389, equal to the real images; skeletons 2366/2365 bytes; `arms: 39`, `real-vs-render (banked) diffs: 0`, `render banked-vs-committed diffs: 1` (header-15px). The r1 seed-2 render is cmp-identical to r0's (V31), and the V14 matrix's mutant lines re-run identical on it |
| V16 | vacuity | the banked programs and the prototype on `vac/*.json` and zero candidates (`v16-v17-vacuity-exit-codes.txt`) | banked: empty-content baseline `same r3 ACCEPT` rc 0; zero candidates rc 0 (both). Prototype: `bad-rectangles graph.png: content rectangles required` rc 1; `zero-regions` rc 1 ×2; zero candidates rc 2 (both); blank `zero-bands` REFUSE rc 1; non-PNG `decode-failed` rc 1 |
| V17 | exit codes | same file | banked REFUSE arm rc 0. Prototype REFUSE rc 1, `accept … ACCEPT regions=17` rc 0, V50 `ACCEPT` rc 0 |
| V18 | why no PNG | `gate-demo.sh` (`v18-gate-demo.sh.txt`): the gate with `verify_go.sh`'s flags on a scratch commit; `marketingcheck docs --assets` (`v18-png-fixture-gate.txt`) | control rc 0 (`0 other binary blobs among 5159`); testdata PNG rc 1 `not-allowlisted`; allowance subdir rc 1 `not-allowlisted`; flat fixture `extra-shot` rc 1 (+ctl empty dir `missing-shot`); `:22 BINARY iff … NUL`; tracked PNGs 0 (+ctl JSON 1255) |
| V19 | CI coverage | `gh run view 41146815448 --job … --log \| grep scripts/marketingcheck`; greps of `.github`, `verify_go.sh` (`v19-ci-coverage.txt`) | plain `26.907s`, race `157.042s`; `WORLD_CHROME` 0 hits (+ctl AILANG_BIN 15); `TestMarketing\|scripts/marketingcheck` in CI/scripts 0 hits (+ctl 4); `verify_go.sh:339` `go test ./... -count=1`, `:355` `-timeout 8m`; the Chrome drill skips without `WORLD_CHROME` (`:54`) |
| V20 | open PRs | `gh pr list --state open --json number,headRefName,files` (`v20-open-prs.txt`) | `[]` rc 0 (+ctl merged #439 files=3, #438 files=28), 2026-10-31T00:31:22Z |
| V21 | NB2 | revert `QUICKSTART.md:1681` to R1; `go test -run 'TestMarketingRunbookMatchesHelp\|TestMarketingDocsLinks' -v` (`v21-nb2-reproduction.txt`) | unmutated rc 0; mutant rc 0, every subtest PASS (survives); (+ctl) `--rectangls` rc 1, `docs_test.go:404/409` |
| V22 | carried (4) | copy the probe in, `go test -run '^TestIter305SeededStillCaptureProbe$' -v` without env; greps over the probe and its orchestration-V40 output (`v22-still-probe.txt`) | `--- SKIP … probe env unset`, `ok`; `:249 t.Logf("STILL … repeat_equal=`; read rule 4 / 0 (+ctl DEFAULT 4); `--- PASS` 1 |
| V23 | mutate determinism | 2 runs × {go1.26.4, go1.26.6} of `p0check mutate`, `cmp` the sha lists (`v23-mutate-determinism.txt`) | all identical; 19 + 19 arm hashes listed |
| V24 | corpus preserved | `cp -p` to `~/.ailang/state/world-iter328-phone-corpus/`; `shasum`; check against `v64-inputs-sha256.txt`; `dtype` (`v24-dtype-main.go.txt`) prints `%T` of `png.Decode` (`v24-phone-corpus-preserved.txt`) | 11 files; 11/11 pinned; all 11 `*image.RGBA <nil>` |
| V25 | JSON transformations | CR count, last byte; `ls .gitattributes`; `git config core.autocrlf`; `check-attr -a`; `grep MarshalIndent` (`v25-json-transformations.txt`) | 13 files CR-bytes=0, last-byte 0a (+ctl CRLF sample 1); no `.gitattributes`; autocrlf unset (rc 1); check-attr empty (+ctl scratch repo `text: unset`); probe `:228 json.MarshalIndent` |
| V26 | cost | `/usr/bin/time -p` on the prototype phone over 39 PNGs plain and `-race`; a cold and a warm file-mode oracle build (`v26-cost.txt`) | 0.18 s; 3.02 s; 2.04 s cold, 0.34 s warm |
| V27 | CLI surface | `sed -n '22,52p' main.go`; runCheck returns; collision grep; png.go type and func lines; `docs_test.go:438–456`; test list; `ls testdata` (`v27-cli-surface.txt`) | 7 subcommands, rc 0/1/2 convention; collisions none (+ctl `case "docs"` 1); `Rect`/`Rectangles`/`LoadRectangles`/`CheckRegion` present; no `testdata` dir |
| V28 | lines to change | greps of QUICKSTART and orchestration for the checker names; §7 vs §4.4 counts; headings (`v28-doc-lines.txt`) | QUICKSTART :1722–1723; orchestration :223, :235–238, :248; §7 0 hits (+ctl §4.4 7); C8′ "(ii) the M7b content-variance gate inside those rectangles"; §13 heading :1663; `main.go:853` `heading absent:` refusal; switch at `main.go:36`; gate `:13 HISTORY`, `:17 STEERABLE`, `:22 BINARY iff`; prototype 196 + 207 lines |
| V29 | (ii) equivalence | `grep` / `sed` of phonecheck `ii :=` and png.go `CheckRectBounds`/`rectStats`/`CheckRegion`/thresholds (`v29-ii-equivalence.txt`) | stand-in `ii := c.HeaderSD >= 4.0 && c.EntrySD >= 4.0` (`:205`); `MinRectH 16`; the R1 header 23 ±4 per edge gives a minimum of 15 |
| V30 | A1 not run | orchestration `:4`; charter greps; `ls website/static/img/v1/` (`v30-a1-status.txt`) | "the controller steps P0 and P3 and the attended A1 remain"; `A1 LANDED\|captured\|PASS` 0 (+ctl `attended A1` 5); v1 dir holds `captions.md` only |
| V31 | **[r1, fix A]** 10 skeletons: the real-corpus check is satisfiable | `skel derive` on each `*-phone.png` of `~/.ailang/state/world-iter328-phone-corpus`; `skel check <corpus> skel10`; three controls; `marketingcheck phone` (prototype) and the banked oracle on the 10 renders (`v31-skel-main.go.txt`, `v31-ten-skeletons.txt`) | 10 files, 23649 bytes; `CHECK images=10 byte-equal=10 geometry-equal=10` rc 0, every line `hash=true … pix+roundtrip=true`. Controls: (1) one byte edited (`cmp -l`: `2193 71 70`) → `byte-equal=9`, rc 1; (2) missing corpus → `BLOCKED: corpus SHA256SUMS: … no such file or directory`, rc 1; (3) one SHA256SUMS line changed → `hash=false` on r1-seed2, rc 1. Renders: `TOTAL candidates=10 accept=10 refuse=0`, rc 0; oracle 10× `-> r3 ACCEPT`; thumb X1 324–340 |
| V32 | **[r1, fix B]** how the V10 prototype was built; the geometry values; the geometry pin and its mutants | `v32-geometry.sh.txt`: `go build -o x p.go` on the prototype alone; `go build ./scripts/marketingcheck` in a `git archive HEAD` tree with the prototype in-package; `sed -n '49,53p' png.go`; `sed -n '57,58p;93p' v64-p0check-main.go.txt`; a python comparison; prototype test `v32-position-geom-test.go.txt` unmutated and under two `Shots` mutants, plus parity and the existing PNG tests under each (`v32-geometry-premise.txt`) | alone: `./p.go:34:18: undefined: Rect`, `./p.go:38:31: undefined: Rectangles`, rc 1; in-tree rc 0; prototype `:41 shot, _ := shotByName(s)`, `png.go:55 func shotByName`. Shots: light/dark 1440×1100, packet 1440×3700, graph 1440×2100, provenance 1440×3000 (phone 390×844); p0check `:57` heights 1100/1100/3700/2100/3000, `:93` `c.X+c.W <= 1440`; `ALL EQUAL`. Test: unmutated `ok`. graph 2101: `FAIL … graph.png: Shots 1440x2101, oracle 1440x2100`, parity `verdict-diffs=0 reason-diffs=0`, `TestMarketingShotGeometry/graph.png` and `TestMarketingPNGContent/canonical-shots` FAIL. packet 3660: `FAIL … packet.png: Shots 1440x3660, oracle 1440x3700`, parity `verdict-diffs=19 reason-diffs=43`. Restored: `png.go` cmp-identical, `ok` |
| V33 | **[r1, fix C]** where the corpus can be rebuilt from | `ls` the corpus; `find ~/.ailang/state -name '*.png' -path '*312*' \| wc -l`; `ls /tmp/world-iter312/*/capture/pass-a/phone.png`; `ls -la` V50; `grep D-WORLD-75` (`v33-corpus-rebuild-sources.txt`) | 11 PNGs + SHA256SUMS; durable iteration-312 copies 0 (+ctl world-iter305 PNGs 36); `/tmp` originals 9 (present today); V50 56838 bytes in `~/.ailang/state/world-iter305`; `QUICKSTART.md:1714` "A fresh seed produces a different store and different bytes" |

## 10. Quorum verification log

**Round 1** on r0 (`9abc25b5`). Artifact `verification/world-iter328/quorum-r1.json`, a byte copy of
`~/.ailang/state/mission-quorum-world/w-marketing-position-checker-2026-10-31T00-42-30Z.json`:
iso_ts `2026-10-31T00:42:30Z`, $0.2120, **BLOCKED**.
- Seats: `gpt6-1-sol` REJECT, `oc-glm-5-3` REJECT, `gemini-4-argon` PASS.
- Claude was benched as the author vendor; no seat was absent.
- Neither reject disputes the direction.

r1 is the one mandated revision. It changes only what is listed here.

- **A. gpt6-1-sol (REJECT).**
  - Verbatim: "The mandatory real-corpus acceptance test is unsatisfiable as specified. §2.3 requires
    a fresh skeleton derived from each of 10 distinct real phone images to be byte-equal to a
    committed skeleton, but M2 commits only two skeletons. Because each skeleton includes its
    source-image sha256, the other eight images cannot match either committed fixture. This blocks
    the capture-day control step rather than providing the claimed verification."
  - **Fix:** commit one skeleton per distinct real image, 10 in all, recorded by basename. That
    covers all 10 real geometries in CI, and byte-equality is per image against its own skeleton.
    This is the recommended route; the two-image restriction was not taken.
  - Proven on the prototype, with three failing controls, and render fidelity re-run on all 10
    (V31, V15). §2.3, §3 M2, §6, §8.
- **B. oc-glm-5-3 (REJECT).**
  - Verbatim: "§2.1's load-bearing port premise — that `position` may take image geometry from
    png.go's `Shots` table (`png.go:49`) because its values are 'the same … as p0check's constants,
    V27' — is asserted, not verified … the doc never states whether the V10 prototype was built
    inside the git-archive tree importing scripts/marketingcheck/png.go or standalone with inlined
    geometry …"
  - **Fix:** the controller's measurements were re-run first-party (V32). The prototype is
    in-package, and the two value lists are printed from the files and are equal.
  - New test row `TestMarketingPositionGeometryMatchesOracle` (§3 M1), with mutation rows P14 and
    P15 (§4). P14 is the verdict-neutral case glm names: parity shows 0 differences, and the new
    test is the sole killer among the position tests.
  - §2.1 now states how the prototype was built.
- **C. gemini-4-argon (PASS, non-blocking).**
  - Verbatim: "The real phone corpus is not checked into git due to tracked_binary_gate.sh and
    instead resides solely in a host-local state directory
    (~/.ailang/state/world-iter328-phone-corpus/), leaving CI reliant on synthetic skeleton renders
    and making capture-day step 1 dependent on unversioned host disk state."
  - **Fix:** §2.5 step 1 and §8 now say so plainly. A missing directory or a differing hash is
    BLOCKED and escalated. The only rebuild is `cp -p` from `/tmp/world-iter312` while it survives
    plus V50 from `~/.ailang/state/world-iter305`. No durable copy of the iteration-312 images
    exists (V33). CI's real-geometry coverage is the 10 skeletons.
  - No storage mechanism is added. Gemini's catch about a `go build` timeout is not in the r1
    directive and is not applied.
