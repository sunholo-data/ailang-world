# w-marketing-predicate-tolerances: bound the free gaps and the drift, and bridge a 1 px thumb split

Status: Planned (iteration 336), **r3**. r2 was the one designer revision after quorum r1 (`world-iter336/quorum-r1.json`, 2026-11-02T20:17:47Z) was BLOCKED (kimi REJECT, sol REJECT, gemini PASS); its fixes are in §8. Quorum r2 (`world-iter336/quorum-r2.json`, 20:29:21Z, same three seats, none absent) was BLOCKED again (kimi REJECT, sol REJECT, gemini PASS). Both remaining objections carry reviewer-authored fixes and neither disputes the design direction, so r3 is the controller's narrow-refinement carve-out, which applies those fixes verbatim (§9).
**Planner-Lane:** mission routing default
- **Item**: queue row 238 (`design_docs/world-mission.md:6310`, V21), row 231's successor S-a
  (`w-marketing-position-checker.md` §7). It takes carried findings (1), (2) and (7).
- **Kind**: release-hygiene, and it moves no clause. It tightens the P0/P3 position predicate and
  loosens the phone thumb's run definition by exactly one pixel.
- **Base**: `37100a432e1b8af9d234396cb8fee7ad627ecbc3`, branch
  `sprint/world-iter336-marketing-predicate-tolerances` (V01). Toolchain `go1.26.6 darwin/arm64`,
  which equals `go.mod` and all three `setup-go` pins in `ci.yml` (V01).
- **Estimated**: 2 executor runs. M1 is code and tests, about 250 changed lines (V12). M2 is the
  orchestration text. **Review**: no content-verified precedent exists (V27: row 231's design quorum
  stored only the doc path, and 0 of the 16 lines its M3 added to orchestration §4.4 appear verbatim
  in either reviewed version of its design; control 16 of 16 in the orchestration doc). Row 231's
  own design also says a §4.4 rule change "needs its own round". M2 therefore requires a one-page
  controller sign-off diffing the pasted A1–A7 against this design's §4 before commit, and the
  orchestration doc's own quorum round on the amended text (§5 M2). The orchestration quorum
  artifacts are dated 10-26 and 10-28 only (V26).
- **Evidence**: `design_docs/verification/world-iter336/designer-*.txt`. Every program ran in a
  `git archive` copy outside the worktree (`$SCRATCH` = `.scratch-iter336-designer/`). The worktree
  gained only this file and `.txt` files (V19). The r1 evidence files are kept, and the `designer-r2-*`
  files supersede them where both exist.

## 1. Problem (measured at the base)

Row 231 committed the r3 predicates unchanged. It pinned three known-permissive verdict groups rather than tightening them:

| Finding | Pin at the base (V06) | What the base does (V09, V11, V15) |
|---|---|---|
| (1) free gaps unbounded | `new/perm-approvals-y+48`, `new/perm-approvals-y-48`, `new/perm-validation-h+216`: ACCEPT | A free gap may move by any whole number of 24 px lines. That is ±2 lines for approvals, and −9 lines when the validation box grows into the gap. |
| (2) compounding drift | `new/perm-roster-drift-2px`: ACCEPT | Each gap and each height stays within ±2 px, yet the roster tops end up +2/+6/+10/+14 px off R1. The judge's gap-only construction (each packet gap +2) gives +4/+6/+8/+10, which is the charter's "8–10 px". It is also ACCEPT. |
| (7) split thumb | `judge/B1-thumb-split-px150`, `judge/B1c-thumb-split-px8`: REFUSE thumb; `judge/B1b-thumb-split-px250`: ACCEPT | The longest contiguous luma 120–200 run decides. A one-pixel split left of the midpoint leaves the right piece longest, so X0 is 151 or 9 and fails X0 ≤ 8. |

The parity tests compare the checker to the banked oracles. The oracle pins `0a9a6795…` and `31576922…` are unchanged (V02). At the base the position ledger is an empty `map[string]bool`, and the phone ledger hard-codes one entry, `header-15px` (V06). Every flip below must therefore become a ledger entry.

## 2. Decisions

### 2.1 (a3): a free gap may move by at most one whole line

The two free gaps are `validation` → `approvals-chain` (R1 239 px) and `approvals-chain` →
`roster-tested` (R1 95 px) (V18). In the template (`render.go:273`, V18):
- the first gap holds the `createdAt`, `policy`, `proposal`, `approval request` and `committed at` paragraphs;
- the second gap holds the two `evidence:` paragraphs.

**Measured:** in all 9 accepted stores both gaps equal R1 exactly (e = 0, k = 0) (V08).

**Rule.** Let k = round(e/24). A free gap passes when it is ≥ 0, |e − 24k| ≤ 2 and **|k| ≤ 1**.
- This is the smallest bound that keeps the one move the template is known to produce: one paragraph wrapping or unwrapping by one line.
- `v65/syn-pos-createdAt-wraps` (k = +1) stays ACCEPT. The bound-0 mutant reds it (V16).
- A store that wraps two more lines in one gap is now **refused and escalated**, not absorbed. This fails closed. No measured store does it, and `bnd-createdAt-wraps-2-lines` pins the case.

I rejected an asymmetric bound for the second gap (k ∈ {0, +1}). It rests on template reasoning that the rectangles cannot check. Instead `perm-approvals-y±24` are pinned as a named residual (§6).

### 2.2 (a5), new: cumulative drift ≤ 2 px below the stack head

(a3) and (a4) each allow ±2 px per term, and those leftovers add up. **Rule (a5)**, using the quorum
r1 wording (sol), applied to every region after a shot's stack head:
- For each earlier region, the height residual is dh − 24·round(dh/24) if it is a packet paragraph, and dh otherwise.
- For each gap through the current region's incoming gap, the gap residual is e − 24·round(e/24) for either designated body-text gap, and e otherwise.
- Here dh and e are signed differences from R1.
- Include only terms whose own (a4) or (a3), respectively, passed.
- The sum must be within ±2 px. A missing region ends the walk for that shot; "missing" already refuses.

"Earlier region" includes the stack head. In the prototype, every region starts with `hRes := dh`, and only a packet paragraph replaces it with the line remainder. `setPrev()` also runs on the stack-head branch (`i == 0`). So a non-paragraph head's full signed dh enters the next region's sum (V23). While every term passes, the sum equals (top − R1 top) − (head top − R1 head top) − 24 × (whole lines above).

**Why failed terms are left out.** Such a candidate is already REFUSE, so the verdict is the same either way. I measured three variants on all 51 inputs (V10, V25). Each adds (a5) lines only to candidates that are already refused:
- **Full sum, failed gaps and heights both included:** adds lines to **10** arms.
- **Gaps-only sum, failed gaps included** (the mutant `M238-a5-accumulate-failed`): reds **11** corpus subtests.
- **Common to both:** the same 9 arms.
- **Only in the full sum:** `v65/neg-tested-height+14`. Its single failed term is a height, which the gaps-only sum skips.
- **Only in the gaps-only sum:** the raw Range-box stores `r1-seed1` and `r1-seed2`. In both, every packet height is 6 px short and every gap 6 px long. These cancel in the full sum and are uncancelled in the gaps-only sum (V25).

**Measured:** the maximum drift is 0 in every shot of all 9 stores (V08). (a5) refuses:
- `perm-roster-drift-2px` at +6/+10/+14;
- the gap-only construction at +4…+10 (V09, V15).

The boundary pairs isolate (a5). Every term in them passes (a3) and (a4) on its own, so only (a5) can decide:
- `bnd-drift-sum-2` / `-3` (packet: roster tops +1 per gap): ACCEPT / REFUSE.
- **`bnd-nonpara-drift-2` / `-3`** (light shot): header (the non-paragraph stack head) H+2, `timeline-first` Y+2 or Y+3. The sum is 2 + 0 or 2 + 1, so ACCEPT / REFUSE (a5). Under the mutant that zeroes non-paragraph height residuals, `-3` becomes ACCEPT and goes red (V24).

### 2.3 (iii′): a run may bridge exactly one 1 px gap

**Rule.** In each bottom-20 px row, list the maximal luma 120–200 segments.
- A **run** is one segment, or two consecutive segments separated by exactly one off-range pixel. A run may contain at most one bridge.
- The leftmost longest run is taken, measured from its first to its last in-range pixel.
- These are unchanged: ≥ 100 px, X0 ≤ 8, X1 ≤ 384 (exclusive), and the first-qualifying-row rule.

| Arm (render and real seed 5) | r3 | r4 | Why |
|---|---|---|---|
| `B1-thumb-split-px150` | REFUSE thumb (X0 151) | **ACCEPT** (6..333) | 6..149 and 151..332 are bridged |
| `B1c-thumb-split-px8` | REFUSE thumb (X0 9) | **ACCEPT** (6..333) | 6..7 and 9..332 are bridged, so the start 6 is kept |
| `B1b-thumb-split-px250` | ACCEPT (6..250) | ACCEPT (6..333) | The left piece was already longest. Only X1 moves, and parity compares verdict and clause booleans, so no ledger entry |
| new `B1d-thumb-two-splits-px100-px200` | REFUSE | REFUSE | It needs two bridges; the best single-bridge run is 101..332 |
| new `B1e-thumb-gap-2px-at-150` | REFUSE | REFUSE | A 2 px gap is not bridged |
| new `B6-alternating-pixels-no-thumb` | REFUSE | REFUSE | With 1 px dots, a run is at most 2 segments and 3 px long |
| `scattered-pixels-no-thumb` | REFUSE | REFUSE | The pixels are 147 px apart |
| new `B7-thumb-two-equal-bridged-runs` (r3, kimi r2) | REFUSE thumb (unbridged segments 6..60, 62..138, 200..259, 261..332; the longest is 77 px, under 100) | **ACCEPT selecting the left run** (6..139) | Paint luma 250 at x 61, x 139..199 and x 260 on rows 833..841. That leaves two bridged runs of equal span, 6..138 and 200..332 (133 px each), and the leftmost-longest rule (strict `>` on length) keeps the left one, so X0 = 6. **Not yet measured**: the r3 verdict column and the X0/X1 values are derived from the code, and M1 must measure them first. If the measured verdict or rectangle differs, M1 stops and reports. It does not re-pin |

Kimi's r2 tie-break arm needs mutant **M238-run-rightmost-tie** (`>` → `>=` in the run-length comparison, so a later equal run wins). It must red `TestMarketingPhoneArms/judge/B7-thumb-two-equal-bridged-runs`, since X0 becomes 200 and the arm refuses. The oracle implements the r3 rule, so B7 is expected to be a REFUSE→ACCEPT parity difference. If so, the phone ledger becomes **4** entries (`header-15px`, B1, B1c, B7), and every "ledger=3" acceptance line and amendment string in this doc reads "ledger=4" / "phone 4". M1 measures this and records it. A different measurement stops M1.

On all 10 real images and the desktop control the thumb rectangle is byte-identical. That means the
10 skeletons and `phoneRealThumbX1` are unchanged (V11, V12 `CHECK images=10 byte-equal=10`). The
rectangle changes on 8 of 97 inputs, and all 8 are synthetic split arms (V11).

### 2.4 New in-test arms (no new fixture files; `testdata/position/SHA256SUMS` stays at 34)

`positionNewArms` grows from 10 to 17 mutations of `stores/r1c-seed5.json`. The sha256 values are the generator's output (V09, V24) and go into `positionNewArmSHA256`:

| Arm | Mutation | sha256 | Verdict |
|---|---|---|---|
| `bnd-drift-sum-2` | attested Y+1, claimed Y+2, proven Y+2 | `65ce8108706c685153f4d0c2b5a98104d7961fabe9ae395da28b7ac8269210f1` | ACCEPT |
| `bnd-drift-sum-3` | attested Y+1, claimed Y+2, proven Y+3 | `da4b9cf1199daeadc7e497b747b1499d687b82e7fc13104b38006c029185a520` | REFUSE (a5) |
| `bnd-createdAt-wraps-2-lines` | approvals-chain **and every roster row** Y+48 | `032e1b0bcd66019af0cd1aa70a60c7afcb51fae75c7e977a8dc9762089183b88` | REFUSE (a3) |
| `perm-approvals-y+24` | approvals-chain only Y+24 (roster fixed) | `79601062a36ca3ea4a827a18e863ff54f446753308c41d198caa1835d9af4895` | ACCEPT (residual pin) |
| `perm-approvals-y-24` | approvals-chain only Y−24 (roster fixed) | `1860035bd84a026039dfdb1021756d5c0b01ee8d93fa289f3f2a8f38d981e7ef` | ACCEPT (residual pin) |
| `bnd-nonpara-drift-2` | workbench-light `header` H+2, `timeline-first` Y+2 | `400d3c1987ce6eb5993b7136dcac101fa6436cf0be57dff27cfe52899b40d98b` | ACCEPT |
| `bnd-nonpara-drift-3` | workbench-light `header` H+2, `timeline-first` Y+3 | `44e4cbd0d7f0e5ca47e317c5561c03a49aeeab3a8fd07e4f4957a954f9679ed9` | REFUSE (a5) |

Phone gains the three arms of §2.3. The arms go from 39 to 42, the CI inputs from 41 to 44, and the real-corpus parity inputs from 59 to 62.
Counts: `positionInputs` goes from 44 to 51, a verdict table of 51 rows with 19 ACCEPT. `phoneWant` has 44 rows with 14 ACCEPT (V12).

### 2.5 The parity ledgers (oracle sources and pins unchanged)

Each ledger maps a label to its expected oracle verdict, checker verdict and evidence. An entry that stops differing is a red, and so is any difference outside the ledger. Every entry must also be visited.

**Position: 8 entries.** All of the oracle's reason lines must still appear, and the added lines must be exactly these (V12):

| Label | Oracle | Position | Lines position adds |
|---|---|---|---|
| `new/perm-approvals-y+48` | ACCEPT | REFUSE | `packet.png/approvals-chain: (a3) gap after validation 287 vs R1 239`; `packet.png/roster-tested: (a3) gap after approvals-chain 47 vs R1 95` |
| `new/perm-approvals-y-48` | ACCEPT | REFUSE | `packet.png/approvals-chain: (a3) gap after validation 191 vs R1 239`; `packet.png/roster-tested: (a3) gap after approvals-chain 143 vs R1 95` |
| `new/perm-validation-h+216` | ACCEPT | REFUSE | `packet.png/approvals-chain: (a3) gap after validation 23 vs R1 239` |
| `new/perm-roster-drift-2px` | ACCEPT | REFUSE | `packet.png/roster-attested: (a5) top 3337 drifts +6 px from R1 3331 after whole-line steps`; `packet.png/roster-claimed: (a5) top 3381 drifts +10 px from R1 3371 after whole-line steps`; `packet.png/roster-proven: (a5) top 3449 drifts +14 px from R1 3435 after whole-line steps` |
| `new/bnd-drift-sum-3` | ACCEPT | REFUSE | `packet.png/roster-proven: (a5) top 3438 drifts +3 px from R1 3435 after whole-line steps` |
| `new/bnd-createdAt-wraps-2-lines` | ACCEPT | REFUSE | `packet.png/approvals-chain: (a3) gap after validation 287 vs R1 239` |
| `new/bnd-nonpara-drift-3` | ACCEPT | REFUSE | `workbench-light.png/timeline-first: (a5) top 426 drifts +3 px from R1 423 after whole-line steps` |
| `v65/neg-tested-at-proven-row` | REFUSE | REFUSE | `packet.png/roster-tested: (a3) gap after approvals-chain 263 vs R1 95`. This is reason-only: e = 168, which is 7 lines, was free under r3 |

**Phone: 3 entries.** They are the same in `TestMarketingPhoneParity` (44 inputs) and in the real-corpus parity (62 inputs):
- `new/header-15px`: oracle ACCEPT, phone REFUSE containing `(ii) rect-too-small phone.png: header`. This entry is unchanged.
- `judge/B1-thumb-split-px150`: oracle REFUSE, phone ACCEPT.
- `judge/B1c-thumb-split-px8`: oracle REFUSE, phone ACCEPT.

The oracle `.txt` files, their pins (V02) and the orchestration pin-registry line are **not changed**. The §4 text must not repeat either hash, because `checkPinRegistry` requires exactly one occurrence of each.

## 3. Before/after verdicts (r3 = committed base, r4 = prototype; V09, V11, V12, V14, V24)

| Input | r3 | r4 |
|---|---|---|
| stores `r1-seed1`, `r1-seed2` (superseded Range-box files, kept as must-REFUSE) | REFUSE (a1),(a2),(a3),(a4) | same |
| `r1-seed1-blockconverted`, `r1-seed2-blockconverted`, `r1b-seed3`, `r1b-seed4`, `r1c-seed5` (= R1, `ba7c83ce…`), `r1c-seed6`, `eval-p1` (`roster-proven` dh −24), `eval-p2`, `eval-p3` | ACCEPT ×9 | ACCEPT ×9 |
| R1 baseline `r1-rects.json` as a candidate | ACCEPT | ACCEPT |
| `new/perm-approvals-y+48`, `new/perm-approvals-y-48`, `new/perm-validation-h+216` | ACCEPT | REFUSE (a3) |
| `new/perm-roster-drift-2px` | ACCEPT | REFUSE (a5) |
| new `bnd-drift-sum-2` / `bnd-drift-sum-3` | ACCEPT / ACCEPT | ACCEPT / REFUSE (a5) |
| new `bnd-nonpara-drift-2` / `bnd-nonpara-drift-3` | ACCEPT / ACCEPT | ACCEPT / REFUSE (a5) |
| new `bnd-createdAt-wraps-2-lines` | ACCEPT | REFUSE (a3) |
| new `perm-approvals-y+24`, `perm-approvals-y-24` | ACCEPT | ACCEPT |
| the other 29 V65/V65-p1/iso/acc arms | as pinned | same verdict and tags. Only `neg-tested-at-proven-row` gains a line, an (a3) line |
| 10 real phone images (9 stores' pass A and V50; thumb X1 324–340) | ACCEPT | ACCEPT, with an identical thumb |
| seed-5 desktop control | REFUSE (i),dims,thumb | same |
| `judge/B1-…px150`, `judge/B1c-…px8` (render and real seed 5) | REFUSE thumb | ACCEPT |
| `judge/B1b-…px250` (render and real seed 5) | ACCEPT | ACCEPT |
| new `B1d`, `B1e`, `B6` (render and real seed 5) | REFUSE thumb | REFUSE thumb |
| the other 36 phone arms (render and real seed 5) | as pinned | same |

**Capture-day control step under r4** (V14). The §4.4 control text needs no change:
- **Step 2:** `accept-seed5`, `accept-eval-p1` ACCEPT, rc 0.
- **Step 3:** `refuse REFUSE` with the same `(a3) gap after roster-tested 29 vs R1 15` line, rc 1.
- **Step 4:** V50 and seed 2 ACCEPT (thumbs 329 and 340), rc 0.
- **Step 5:** the desktop shot, rc 1.

## 4. Orchestration amendment (M2 pastes this text verbatim; this document's quorum reviews it)

**Pre-paste check (hard stop).** In `design_docs/planned/w-marketing-capture-orchestration.md`,
`grep -cF '<anchor>'` must give exactly `1` for each of the seven anchors A1–A7 below. The values
at the base are recorded in V20. Any other count stops M2. The executor then re-bases, re-runs V20
and asks the controller. It never pastes by line number.

**A1, (a3).** Replace the whole bullet whose line contains `the gap may instead differ from the R1 gap by a whole number of 24 px body lines, ±2 px;` with:

> - (a3) every later region's gap (its top minus the previous region's **measured** bottom) is ≥ 0 and within ±2 px of the R1 gap. Across the two packet gaps whose intervening content is body-text paragraphs (`validation` → `approvals-chain`, `approvals-chain` → `roster-tested`), the gap may instead differ from the R1 gap by **one** whole 24 px body line, up or down, ±2 px (row 238: a move of two or more lines in one gap is refused and escalated, not absorbed);

**A2, (a4).** Leave unchanged. Insert the following as a new bullet immediately after the bullet whose line contains `Every other region's height is within ±2 px of R1;`:

> - (a5) (row 238) cumulative drift: for every region after a shot's stack head, walking the regions in R1 order, the sum of the following residuals is within ±2 px. For each earlier region, the height residual is dh − 24·round(dh/24) if it is a packet paragraph, and dh otherwise. For each gap through the current region's incoming gap, the gap residual is e − 24·round(e/24) for either designated body-text gap, and e otherwise. Here dh and e are signed differences from R1. Include only terms whose own (a4) or (a3), respectively, passed. "Earlier region" includes the stack head, and a missing region ends the walk for that shot. While every term passes, each region therefore lies within ±2 px of its R1 top plus the stack head's measured shift plus the whole-line steps above it, so ±2 px per term cannot compound;

**A3, Committed bullet.** Replace `or any verdict difference outside the one-entry ledger, is a red.` with `or any verdict or reason difference outside the ledgers (row 238: position 8 entries, phone 3; each fixes the oracle verdict, the checker verdict and the differing reason), is a red.`

**A4.** Replace `moved by anything other than whole lines in the two body-text gaps` with `moved by anything other than one whole line in either body-text gap, drifted more than 2 px in sum below its stack head (row 238)`.

**A5, (iii′).** Replace `(a luma 120–200 run of at least 100 px)` with `(a luma 120–200 run of at least 100 px, where a run is one contiguous segment or two segments separated by exactly one off-range pixel, with at most one such bridge per run, measured from its first to its last in-range pixel; row 238)`.

**A6, §7 C8′.** Replace `and packet paragraph heights and the two body-text gaps in whole 24 px lines ±2 px (§4.4, r3)` with `packet paragraph heights in whole 24 px lines ±2 px, the two body-text gaps by at most one whole 24 px line ±2 px, and a cumulative drift of at most ±2 px below each stack head (§4.4, r3; row 238 (a3)/(a5))`.

**A7, §9.** Insert the following immediately before the line `## 10. Verification Log` (that is, after the Row 231 subsection):

> ### Row 238 (iteration 336, `w-marketing-predicate-tolerances`)
>
> Carried findings (1), (2) and (7) are taken. (a3) bounds each body-text gap to one whole line. The new (a5) bounds the cumulative drift below a stack head to ±2 px, with non-paragraph heights counting in full. The (iii′) run bridges exactly one 1 px gap, once. The `perm-approvals-y±48`, `perm-validation-h+216` and `perm-roster-drift-2px` pins flip to REFUSE; `B1-thumb-split-px150` and `B1c-thumb-split-px8` flip to ACCEPT. All 9 iteration-312 stores, the R1 baseline, the 10 real phone images and the capture-day controls keep their verdicts. The parity ledgers name every difference from the unchanged oracles (position 8, phone 3). This amendment was reviewed as that design's §4, as row 231's was. Evidence: that design's V-rows, under `design_docs/verification/world-iter336/`.

## 5. Milestones

**Execution rule (r3, gpt6-1-sol r2, verbatim):** "Every executor command runs under a process-group deadline. Use a checked-in or explicitly verified runner that terminates the command and its descendants on expiry, records a nonzero timeout result, and stops the milestone. Set explicit limits for build/vet/full-suite verification, each mutation run, and each capture-day control. `go test -timeout` remains an additional test-binary limit, not the outer deadline."
- **Runner:** `design_docs/verification/world-iter336/run_bounded.sh.txt`, copied to a scratch path as `run_bounded.sh`. It is the `run_bounded()` idiom from `scripts/verify_ail.sh:214`: a new session, then `os.killpg(…, SIGKILL)` on expiry, exit 124. It is verified by V28, where a hanging grandchild is killed in 3 s, rc is 124, and the next `&&` step is skipped; the controls give rc 7 and rc 0, passed through.
- **Invocation:** `bash "$RB" <limit_s> "$OUT" <cmd…>; rc=$?`. Any rc of 124 stops the milestone with no continuation to M2.
- **Limits:**
  - `go vet`, `gofmt -l` and `go build ./...`: 600 s each.
  - Full `go test ./...` (with `AILANG_BIN` set to the pinned binary): 2400 s.
  - `./scripts/verify_ail.sh`: 1200 s.
  - The `TestMarketing(Position|Phone)` acceptance run: 1500 s, with `-timeout 20m` inside it.
  - Each mutation run: 600 s.
  - Each capture-day control (steps 2–5): 300 s.

**M1: predicates and tests** (one executor run).
- **Files:** `scripts/marketingcheck/position.go`, `phone.go`, `position_test.go`, `phone_test.go`.
- **Reference patch:** `designer-r2-patch.diff.txt` (4 files, +210/−39, V12).
- **Also update by hand:** the stale count comments ("the 44 position inputs", "44-row verdict table … 19 ACCEPT", "41-row … 12 ACCEPT", "39 arms"), and the r3-only file-header comments of `position.go` and `phone.go`.

Acceptance. Check every exit code, and pipe nothing:
1. `WORLD_PHONE_CORPUS=$HOME/.ailang/state/world-iter328-phone-corpus go test -count=1 -timeout 20m -v -run '^TestMarketing(Position|Phone)' ./scripts/marketingcheck > "$LOG" 2>&1; echo rc=$?` prints `rc=0`. In the log, `grep -c '^--- SKIP'` is 0, and the log contains exactly:
   - `PARITY position: inputs=51 verdict-diffs=0 reason-diffs=0 ledger=8`
   - `PARITY phone: inputs=44 diffs=0 ledger=3`
   - `PARITY phone: inputs=62 diffs=0 ledger=3`
   - `CHECK images=10 byte-equal=10`
2. All of these give rc 0: `go vet ./scripts/marketingcheck`, `gofmt -l scripts/marketingcheck` (with empty output), `go build ./... && go test ./...`, and `./scripts/verify_ail.sh`.
3. `git diff --exit-code 37100a43 -- design_docs/verification/world-iter315/ scripts/marketingcheck/testdata/` gives rc 0.
4. Capture-day steps 2–5 give the V14 outputs.

Mutations. Apply each one alone. Each must red **every** subtest listed below, which is the complete list measured in `designer-r2-mutants.txt` (V16). Every listed mutant also reds the parent `TestMarketingPositionCorpus` or `TestMarketingPhoneArms` and the parent parity test:

| Id | Mutation | Must-red subtests |
|---|---|---|
| M238-a3-unbounded | drop `&& posAbs(posLines(e)) <= positionFreeGapLines` | `TestMarketingPositionCorpus/new/perm-approvals-y+48`, `TestMarketingPositionCorpus/new/perm-approvals-y-48`, `TestMarketingPositionCorpus/new/perm-validation-h+216`, `TestMarketingPositionCorpus/new/bnd-createdAt-wraps-2-lines`; `TestMarketingPositionParity` (ledger `v65/neg-tested-at-proven-row` and the four above) |
| M238-a3-bound2 | `positionFreeGapLines = 2` | `TestMarketingPositionCorpus/new/perm-approvals-y+48`, `TestMarketingPositionCorpus/new/perm-approvals-y-48`, `TestMarketingPositionCorpus/new/bnd-createdAt-wraps-2-lines`; `TestMarketingPositionParity` |
| M238-a3-bound0 | `positionFreeGapLines = 0` | `TestMarketingPositionCorpus/v65/syn-pos-createdAt-wraps`, `TestMarketingPositionCorpus/new/iso-minh-approvals-h-1`, `TestMarketingPositionCorpus/new/perm-approvals-y+24`, `TestMarketingPositionCorpus/new/perm-approvals-y-24`; `TestMarketingPositionParity` |
| M238-a5-off | (a5) never evaluated (`if chain && false`) | `TestMarketingPositionCorpus/new/perm-roster-drift-2px`, `TestMarketingPositionCorpus/new/bnd-drift-sum-3`, `TestMarketingPositionCorpus/new/bnd-nonpara-drift-3`; `TestMarketingPositionParity` |
| M238-a5-tol3 | `positionDriftTol = 3` | `TestMarketingPositionCorpus/new/bnd-drift-sum-3`, `TestMarketingPositionCorpus/new/bnd-nonpara-drift-3`; `TestMarketingPositionParity` |
| M238-a5-tol1 | `positionDriftTol = 1` | `TestMarketingPositionCorpus/new/bnd-drift-sum-2`, `TestMarketingPositionCorpus/new/bnd-nonpara-drift-2`; `TestMarketingPositionParity` |
| **M238-a5-nonpara-zero** | `hOK, hRes := true, dh` → `hOK, hRes := true, 0` (non-paragraph height residual zeroed) | `TestMarketingPositionCorpus/new/bnd-nonpara-drift-3`; `TestMarketingPositionParity` (ledger `new/bnd-nonpara-drift-3`) |
| M238-a5-accumulate-failed | a failed gap's residual is also added | `TestMarketingPositionCorpus/stores/r1-seed1`, `…/stores/r1-seed2`, `…/v65/neg-attested+14`, `…/v65/neg-dark-timeline-first+24`, `…/v65/neg-light-header+40`, `…/v65/neg-light-header-whole-view`, `…/v65/neg-rows-reordered-tested-claimed`, `…/v65/neg-tested-64x16-glyphs`, `…/v65/neg-tested-at-attested-row`, `…/v65/neg-tested-at-proven-row`, `…/v65-p1/neg-tested-at-attested-row` (each `…` = `TestMarketingPositionCorpus`; 11 subtests); `TestMarketingPositionParity` |
| M238-run-nobridge | never bridge | `TestMarketingPhoneArms/judge/B1-thumb-split-px150`, `TestMarketingPhoneArms/judge/B1c-thumb-split-px8`; `TestMarketingPhoneParity` (both ledger entries) |
| M238-run-gap2px | bridge gaps ≤ 2 px | `TestMarketingPhoneArms/judge/B1e-thumb-gap-2px-at-150`; `TestMarketingPhoneParity` |
| M238-run-chain-bridges | unlimited bridges per run | `TestMarketingPhoneArms/judge/B1d-thumb-two-splits-px100-px200`, `TestMarketingPhoneArms/judge/B6-alternating-pixels-no-thumb`; `TestMarketingPhoneParity` |
| M238-run-rightmost-tie (r3) | the run-length comparison `>` becomes `>=` (a later equal run wins) | `TestMarketingPhoneArms/judge/B7-thumb-two-equal-bridged-runs`; `TestMarketingPhoneParity`. Derived from the code; M1 measures it first |
| F4 (row 231) | `X0 <= 8` removed | `TestMarketingPhoneArms/judge/B5-two-pieces-rows-differ`, `TestMarketingPhoneArms/judge/short-thumb-x20`, `TestMarketingPhoneArms/judge/thumb-x9`, `TestMarketingPhoneArms/judge/B1d-thumb-two-splits-px100-px200`, `TestMarketingPhoneArms/judge/B1e-thumb-gap-2px-at-150`; `TestMarketingPhoneParity` (B1/B1c no longer kill it) |
| F9 (row 231) | first-to-last span | `TestMarketingPhoneArms/judge/scattered-pixels-no-thumb`, `TestMarketingPhoneArms/judge/B1d-thumb-two-splits-px100-px200`, `TestMarketingPhoneArms/judge/B1e-thumb-gap-2px-at-150`, `TestMarketingPhoneArms/judge/B6-alternating-pixels-no-thumb`; `TestMarketingPhoneParity` |

**M2: orchestration text** (one executor run, after M1 passes; M1 and M2 land together or not at all, because M1 changes the rules §4.4 defines).
- **Files:** `design_docs/planned/w-marketing-capture-orchestration.md` only.
- **Steps:** run the §4 pre-paste check, then paste A1–A7.
- **Controller sign-off (r3, kimi r2):** before commit, the controller diffs each pasted A1–A7 block in the orchestration doc against this design's §4 text and banks the result as `world-iter336/m2-signoff.txt`. Any delta is a hard stop.
- **Orchestration quorum (r3, reinstated):** after the sign-off, the controller runs `ailang design-quorum design_docs/planned/w-marketing-capture-orchestration.md --author claude:claude-opus-5-5` on the amended doc. A BLOCKED round with a direction objection means neither M1 nor M2 lands, and the row parks. The Gate-2 one-revision and narrow-refinement rules apply as for any doc.

Acceptance:
- `go test -count=1 -timeout 300s -run '^TestMarketing(Position|Phone)Parity$' ./scripts/marketingcheck` gives rc 0, including both `pin-registry` subtests.
- `grep -cF` gives `1` for each of `(a5) (row 238) cumulative drift`, `exactly one off-range pixel`, `### Row 238 (iteration 336` and `position 8 entries, phone 3`.
- `grep -cF` gives `0` for `outside the one-entry ledger`.

Mutation **M238-doc-pin-dup**: paste either oracle hash into the Row 238 subsection. `TestMarketingPositionParity/pin-registry` or `TestMarketingPhoneParity/pin-registry` must go red.

## 6. Conflict surface, concurrent work, residuals, premises checked

- **Conflict surface.** The change is Go only, under `scripts/marketingcheck`, plus one Markdown file.
  - Nothing else changes: no `.ail`, `go.mod`, `.github/`, `tools/launchd/*`, fixtures, oracle sources or `main.go`.
- **Concurrent-work sweep (V22, controller-measured).** I found no concurrent work on these files:
  - There are 0 open PRs on `sunholo-data/ailang-world`.
  - No origin branch changes `scripts/marketingcheck` or the orchestration doc relative to `37100a43`. The control shows the instrument does see changes: `c6a12dc9^..origin/dev` lists 15 files.
  - The last commits on those paths are row 231's `03fe5a1a`, `9ddab47c` and `c6a12dc9`.
  - Three NEXT/IN-SPRINT rows name the files:
    - 150 is the attended release row, waiting on Mark's A1 capture.
    - 238 is this row.
    - 240 (`w-marketing-checker-tolerance-pins`) is queued **after** 238 in D-WORLD-106 order, so it runs sequentially on top of this.

  If a concurrent change appears before M1 or M2 executes, this mission re-bases and re-runs V20, and M2's pre-paste check enforces that for the orchestration doc.
- **Residual (pinned, not fixed).** Moving only `approvals-chain` by ±24 px is ACCEPT (`perm-approvals-y±24`). The rectangles cannot tell it apart from a real `createdAt` unwrap combined with an `evidence:` wrap. (a4) heights stay unbounded in line count: `syn-pos-proven-5-lines` and `v65-p1/neg-proven-grows-6-lines-below-fold` stay ACCEPT, and only (b) caps growth.
- **Premises checked.**
  - (i) The charter's "8–10 px" is the judge's gap-only construction (V15). The pinned arm also grows heights and reaches +14. Both are refused.
  - (ii) Of the 11 store files, 2 are superseded Range-box captures that already REFUSE. The 9 judged stores include the 2 block-converted ones.
  - (iii) The `perm-*`/`B1*` pins mark every **verdict** flip among the existing arms, but not every parity difference. `v65/neg-tested-at-proven-row` gains a reason line (V09). Both ledgers also change shape (V06).
  - (iv) (a4)'s text is unchanged; the drift bound is the new (a5).

## 7. Verification Log (base `37100a43`; files under `design_docs/verification/world-iter336/`)

| V | Command | Observed |
|---|---|---|
| V01 | `git rev-parse HEAD; go version; grep -n '^go ' go.mod; grep -n "go-version: '1.26.6'" .github/workflows/ci.yml` (`designer-env.txt`) | `37100a43…`; `go1.26.6 darwin/arm64`; `3:go 1.26.6`; ci.yml lines 38, 130, 432 |
| V02 | `shasum -a 256` of both oracle `.txt`, and a count of each pin in the orchestration doc (`designer-env.txt`) | `0a9a6795…927c` and `31576922…0b4c`; each pin occurs once |
| V03 | `shasum -a 256 r1-rects.json stores/r1c-seed5.json` | both give `ba7c83ce…3333` |
| V04 | `shasum -a 256 -c corpus-SHA256SUMS` inside the corpus dir, plus the V50 sha | 11 × `OK`, rc 0; V50 `999125ff…4be8` |
| V05 | M1 acceptance 1 at the base (`designer-baseline-tests.txt`) | rc 0, 0 SKIP; `PARITY position: inputs=44 verdict-diffs=0 reason-diffs=0`; phone `inputs=41 diffs=0 ledger=1`; real `inputs=59 diffs=0 ledger=1` |
| V06 | `grep -n 'ledger :='` and the `perm-`/`B1` verdict rows (`designer-env.txt`) | position `map[string]bool{}`; phone `{"new/header-15px": true}`; the 4 perm arms ACCEPT; B1 and B1c REFUSE thumb; B1b ACCEPT |
| V07 | prototype run with r3 options, compared against the committed functions (`designer-proto-*.txt`, `designer-r2-proto.txt`) | position `CONTROL … 0 of 51 differ`; phone `0 of 97 differ` |
| V08 | `GEOM` lines (`designer-proto-position.txt`) | all 9 stores give `approvals-chain:e=0,k=0 roster-tested:e=0,k=0`, and every `maxTopDrift=0`; eval-p1 gives `roster-proven:dh=-24` |
| V09 | verdict table and `ARMSHA` lines (`designer-r2-proto.txt`) | the 9 stores and R1 ACCEPT under both; 4 perm flips; the 7 shas of §2.4; `neg-tested-at-proven-row` gets +1 (a3) line |
| V10 | `FULL-VARIANT` lines (`designer-proto-position.txt`) | `FULL vs ORTHOGONAL (a5): accept/refuse differs on 0 of 49` (0 of 51 in r2, `designer-r2-proto.txt`) |
| V11 | `TestIter336ProtoPhone` with `WORLD_PHONE_CORPUS` (`designer-proto-phone.txt`) | 10 real ACCEPT, thumb unchanged; B1 and B1c FLIP; B1b `{6 833 250 841} -> {6 833 333 841}`; `THUMB-RECT changed … 8 of 97` |
| V12 | M1 acceptance 1 on the r2 patched copy (`designer-r2-tests.txt`; patch `designer-r2-patch.diff.txt`, diffstat 4 files +210/−39) | rc 0, 0 SKIP; exactly the four §5 lines (`inputs=51 … ledger=8`) |
| V13 | `go test -count=1 -timeout 600s ./scripts/marketingcheck`, patched (r1) copy and base copy (`designer-r4-tests.txt`) | both `ok`. r2 changes only `position_test.go`; V12 covers it |
| V14 | capture-day steps 2–5 on the patched copy (`designer-r4-controls.txt`) | rc 0 / 1 / 0 / 1, outputs as in §3. `position.go` is identical in r1 and r2 (`designer-r2-patch.diff.txt`) |
| V15 | gap-only drift arm, base vs patched (`designer-r4-controls.txt`) | base `ACCEPT`; r4 `REFUSE` (a5) +4, +6, +8, +10 |
| V16 | `mutate.py` × 13, each run `-run '^TestMarketing(Position\|Phone)(Corpus\|Arms\|Parity)$'` (`designer-r2-mutants.txt`) | every mutant gives rc 1 and reds exactly the §5 lists |
| V17 | row 231's F4 and F9 on the patched copy (same file) | both rc 1, as listed in §5 |
| V18 | `sed -n 273p host/workbench/render.go \| grep -o '<p>…'`; the R1 packet gaps (`designer-template-and-targets.txt`) | paragraph order as in §2.1; gaps 239, 95, 15, 15, 15 |
| V19 | `git -C <worktree> status --porcelain` at hand-off | only `?? design_docs/planned/w-marketing-predicate-tolerances.md` and `?? design_docs/verification/world-iter336/` |
| V20 | `grep -cF '<anchor>'` for A1–A7, plus controls (`designer-r2-m2-anchors.txt`) | 1 for each of the 7 anchors. Controls: `(a3)` gives 6 and a misspelled anchor gives 0 |
| V21 | `sed -n 6310p design_docs/world-mission.md` (`designer-env.txt`) | the row text quoted in the brief |
| V22 | controller sweep, verbatim (`designer-V22-sweep-controller.txt`; source sha256 `5f23408e…92ce`) | open PRs: none listed. Branch diffs since base: `0`; control: `15 files`. Last commits `03fe5a1a`, `9ddab47c`, `c6a12dc9`. Rows naming the files: 150, 238, 240 (status-filter control: 4 rows without the filter) |
| V23 | `grep -n` / `sed -n` of the residual code in the r2 `position.go` (`designer-r2-code.txt`) | `90: hOK, hRes := true, dh`; `92: hRes = dh - posLines(dh)*positionLine` (paragraphs only); `109: if i == 0 {` … `113: setPrev()`; `134: drift += prevRes`. Control: M238-a5-nonpara-zero edits line 90 and reds V24's arm |
| V24 | `TestIter336ProtoR2` `NONPARA` lines (`designer-r2-proto.txt`) and the M238-a5-nonpara-zero run (`designer-r2-mutants.txt`) | `bnd-nonpara-drift-3 r3=ACCEPT r4=REFUSE (a5) nonpara-zero-mutant=ACCEPT` with reason `workbench-light.png/timeline-first: (a5) top 426 drifts +3 px from R1 423`; `-2` gives ACCEPT/ACCEPT/ACCEPT. Mutant: rc 1, `--- FAIL: TestMarketingPositionCorpus/new/bnd-nonpara-drift-3`. The oracle gives ACCEPT on both: the parity passes with its ledger entry (V12) |
| V25 | `TestIter336ProtoR2` `A5LINES`, `FULL adds`, `GAP-ONLY adds` and `RAW` lines (`designer-r2-proto.txt`) | the full sum adds lines on 10 arms and the gaps-only sum on 11, with the split shown in §2.2. Raw seeds: every packet `(a4) height … vs …` is −6 and every `(a3) gap` is +6 |
| V26 | `ls ~/.ailang/state/mission-quorum-world/ \| grep w-marketing-capture-orchestration`; `git log -1 03fe5a1a` (`designer-r2-quorum-precedent.txt`) | 4 artifacts, dated 2026-10-26 and 2026-10-28 only. Control: `w-marketing-position-checker-2026-10-31T00-42-30Z/…51-18Z`. `03fe5a1a` is 2026-10-31 03:26 |

| V27 | controller: `jq keys` on both row-231 quorum artifacts; `git show 3aaee4d0:`/`08dc31be:` the design; `git show 03fe5a1a -- <orchestration>` added lines, each `grep -qF` in both design versions and in the orchestration doc (`v27-precedent-content-check.txt`) | artifacts store `.doc` = path only; 0 of 16 added lines are in the design r1 reviewed by quorum r2, and 0 of 16 are in carve-out r2. Control: 16 of 16 are in the orchestration doc at `03fe5a1a`. Row 231 design r1 line 257: §4.4 rule change "needs its own round" |
| V28 | controller: `bash v28-demo.sh` over the banked runner (`v28-run-bounded-demo.txt`, runner sha256 `1635f8d7…`) | hanging `sh -c 'sleep 300 & …; wait'` with a 3 s limit: rc 124, elapsed 3 s, grandchild dead; `run_bounded 2 … sleep 30 && echo NEXT` gives `next step skipped, rc=124`. Controls: `exit 7` gives rc 7 and `true` gives rc 0 |

## 8. Quorum r1 dispositions (r2)

| # | Reviewer: objection (their terms) | Fix applied | Where |
|---|---|---|---|
| 1 | gpt6-1-sol (REJECT): §4's (a5) phrase "its whole height change (other regions)" can be read as zeroing non-paragraph residuals. Replace it with their text; verify that non-paragraph dh counts, including the stack head; add a boundary test and a non-paragraph-zero mutant | Their four sentences are pasted verbatim into §2.2 and §4 A2. The code was confirmed (V23). I added `bnd-nonpara-drift-2`/`-3`, whose header is the non-paragraph stack head. The mutant `M238-a5-nonpara-zero` reds `-3` (V24). There is a new ledger entry because the oracle gives ACCEPT | §2.2, §2.4, §2.5, §4 A2, §5 |
| 2a | oc-kimi-k3 (REJECT): the conflict surface is asserted, not verified | Added the controller's sweep as V22 and kimi's re-base sentence | §6, V22 |
| 2b | kimi: M2 depends on absolute line numbers 230/246/248/250/314 | Changed to unique-string anchors A1–A7 with a `grep -cF` = 1 pre-paste hard stop. All 7 measured 1, with controls | §4, §5 M2, V20 |
| 2c | kimi: ellipses in the mutation table; V10's 10 vs the mutant's 11 not reconciled | Every must-red subtest is enumerated from the r2 re-run. Reconciled: the two variants differ, 9 arms in common. Full adds `neg-tested-height+14`; gaps-only adds raw `r1-seed1`/`r1-seed2` | §5, §2.2, V16, V25 |
| 3 | gemini-4-argon (PASS, note): the ±24 arms should say approvals-chain only | The §2.4 cells now read "approvals-chain only … (roster fixed)", and the `bnd-createdAt` cell says the roster moves too | §2.4 |
| 4 | controller: M2 waited for a second quorum | Removed in r2, citing row 231's precedent (V26). **Reversed in r3**: V27 refutes that precedent by content (see §9) | header, §5 M2 |

## 9. Quorum r2 dispositions (r3, controller narrow-refinement carve-out)

Both r2 seats that rejected (`oc-kimi-k3`, `gpt6-1-sol`) gave a concrete `proposed_fix`, and neither disputes the design direction (the predicates, bounds and bridge rule). `gemini-4-argon` passed. `absent_reviewers` is `[]`. The r2 objections land on two surfaces that r1 did not raise: review provenance and command bounding. They do not localise onto one surface while a reviewer flips, so this is not a SPLIT signal.

| # | Reviewer: objection (their terms) | Fix applied (verbatim where they wrote text) | Where |
|---|---|---|---|
| 5 | oc-kimi-k3 (REJECT): the row-231 precedent is "only timestamp-verified, not content-verified" | Measured by the controller (V27): the artifacts store only the doc path, and 0 of the 16 added §4.4 lines are verbatim in either reviewed design version (control 16 of 16). The precedent FAILS. Applied kimi's branch (b): "No content-verified precedent exists; M2 therefore requires a one-page controller sign-off diffing the pasted A1–A7 against this design's §4 before commit". The orchestration doc's own quorum is also reinstated, as row 231 §7 and row 238 both require | header, §5 M2, V27 |
| 5b | kimi (optional): no arm exercises the leftmost-longest tie-break in the bridged state space | Added arm `B7-thumb-two-equal-bridged-runs` (ACCEPT selecting left) and mutant `M238-run-rightmost-tie`; M1 measures both before relying on them | §2.3, §5 |
| 6 | gpt6-1-sol (REJECT): acceptance runs `go vet`, `go build ./...` and `./scripts/verify_ail.sh` "without process-level deadlines" | Added gpt's execution rule verbatim, with a banked runner, explicit limits and V28, which shows a hanging child process being killed and the next step not running | §5, V28 |
