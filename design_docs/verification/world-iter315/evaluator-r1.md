# Evaluator report: iteration 315, row 150 design amendment (c36a4306, doc r5)

Evaluator: sonnet (agent-tool, alias-pin), independent of the Opus designer and controller.
Scratch builds and arms: `/tmp/world-iter315-eval/` (checkers built from the banked `.txt` exactly per §4.4; source sha256 matched the doc's pins).

## Score: 77 / 100 -> PASS (threshold 70). No BLOCKING findings.

| Criterion | Max | Score |
|---|---|---|
| Correctness of the rules vs evidence, both directions | 30 | 24 |
| Verification / non-vacuity of the evidence rows | 25 | 20 |
| Doc consistency and completeness | 20 | 13 |
| Carve-out / process honesty | 15 | 11 |
| Scope discipline | 10 | 9 |

## What reproduced (independent re-runs)
- V64: `p0check check r1-rects.json` over the 11 candidates: output byte-identical to banked `v64-p0-positive.txt` (`r3-accept=9 r2-accept=8`). Seeds 1-2 only pass after the block-box conversion, as stated.
- V66: `phonecheck` over V50 + 9 stores x 2 passes: byte-identical to banked `v66-phone-positive.txt` (`r3-accept=19 r2-accept=1`).
- V67: rebuilt `v67-phoneneg` from the banked source: 13 candidates, 2 accepted (shift-2, shift-4), 11 refused. Matches.
- V65: regenerated `p0check mutate` arms: seed 5 gives 16 refuse + 3 accept, eval-p1 subset 2 refuse + 2 accept = 18 refused / 5 accepted. Matches.
- V72: `v64-rects-r1c-seed5.json` sha `ba7c83ce...`, V50 sha `999125ff...94be8`, both source sha256 pins match. Accept/refuse controls print as banked, rc 0 for both verdicts.
- V69: ancestry rc 0 for 4e6e89dc, eacd90e7, dba643a9; reversed control rc 1. `git diff --name-only 4e6e89dc..09ff2821 -- host/workbench` = 0 files, `git log` 0 commits; positive range control `aeaa957d` lists render.go and render_test.go; `host/daemon/workbench*` 0, `host/pins` 0; nothing in the capture scripts, host/workbench, pins after dba643a9. daemon.go diff has no still/workbench lines. Claim holds.
- V73: eval-p1 layout roster-proven 3435.5625..3507.5625 (72.0), rects h 73; seed 5 96.0 / h 97 = r1-rects.json. Holds.
- The retained iter305 probe compiles at HEAD (`go test -run '^$' ./cmd/ailang-worldd` with it copied in: ok; file removed again).
- quorum-r3/r4 json are byte copies of the `~/.ailang/state/mission-quorum-world/...` files; verdicts as the doc states.

## (a) Code vs prose
p0check matches (a1)-(a4), background, (b) clause by clause, including the free-gap set (only `approvals-chain`, `roster-tested`) and `gc >= 0`. phonecheck matches (i), (ii), clip >= 387, thumb X0 <= 8 and X1 <= 384 (X1 is exclusive), with these divergences:
1. NON-BLOCKING. Prose: thumb is "a luma 120-200 run of at least 100 px". Code (`scrollbar()`): span x1-x0 >= 100 over any pixels in the row, not a contiguous run. Arm `scattered-pixels-no-thumb` (two mid-gray pixels at x=3 and x=150 in an otherwise blank bottom strip) -> `r3 ACCEPT`. Impractical on real pages (needs ink at x <= 8) but the checker is weaker than the sentence.
2. NON-BLOCKING. (ii) is a stand-in (std-dev >= 4.0 = MinRectStdDev only), not the M7b gate (min size, colors); the doc says "the M7b gate inside those rectangles" and relies on `marketingcheck check` separately. Fine, but say so.
3. NON-BLOCKING. Prose "ends at x <= 384" vs code X1 exclusive: a thumb whose last pixel is 383 (X1=384) is accepted, X1=385 refused. State the convention.

## (b) My own negative arms (not in the designer's set)
Base: seed-5 rects (= R1) unless noted. SHOULD = my judgement of the right verdict.

P0 (`/tmp/world-iter315-eval/neg/`, `neg2/`):
| Arm | SHOULD | Result |
|---|---|---|
| roster-attested y+24 alone (line step where not allowed) | refuse | REFUSE (a3 gap 39, next row -9) |
| roster-claimed y+24 alone | refuse | REFUSE (a3) |
| light timeline-first y+24 | refuse | REFUSE (a3 gap 124 vs 100) |
| graph-svg h+24 | refuse | REFUSE (a4) |
| provenance payload h+24 | refuse | REFUSE (a4) |
| packet stack head `validation` y+24 | refuse | REFUSE (a2) |
| roster-tested/attested/claimed h+24 with top fixed (overlaps next row) | refuse | REFUSE x3 (a3 gap -9) |
| tested h+48 with attested y+24 only | refuse | REFUSE (a3) |
| `validation` h+240 overlapping `approvals-chain` by 1 px | refuse | REFUSE (a3 gap -1) |
| `approvals-chain` h = 1 px sliver (dh -24) | refuse | REFUSE (a4 min 24) |
| roster-tested w-3 / w-2 | refuse / accept | REFUSE (a1) / ACCEPT |
| background y+3 / bottom +2 | refuse | REFUSE / REFUSE (via (b), baseline bottom sits at H-40) |
| all four roster rows y+240 (below fold) | refuse | REFUSE (b) |
| all four roster rows y+24 (a body paragraph wrapped above) | accept | ACCEPT (by design) |
| `approvals-chain` y+48 / y-48 (a rect that would sit on a different paragraph) | arguably refuse | ACCEPT x2 |
| `validation` h+216 (inflated into blank space, no overlap) | arguably refuse | ACCEPT |
| compounding drift: every roster gap +2 and height +2 (row ~8-10 px off R1) | accept per rule text | ACCEPT |
| extra junk region / regions reversed in array | ignore | ACCEPT / ACCEPT (name-matched) |

Phone (`/tmp/world-iter315-eval/phn/`, from seed-5 pass-A):
| Arm | SHOULD | Result |
|---|---|---|
| short thumb starting at x=20 (x 20-300) | refuse | REFUSE |
| thumb x 9-330 / x 8-330 | refuse / accept | REFUSE / ACCEPT (boundary exact) |
| thumb ending X1=384 / X1=385 | accept / refuse | ACCEPT / REFUSE |
| title band rightmost ink 386 / 387 | refuse / accept | REFUSE / ACCEPT |
| thumb moved above the bottom-20 band | refuse | REFUSE |
| dark thumb (luma 100) | refuse | REFUSE |
| short thumb x 6-130 (large overflow) | accept (defect present) | ACCEPT |
| short thumb x 6-90 (< 100 px) | accept if a real wide-overflow store | REFUSE (see (c)) |
| scattered gray pixels, no thumb | refuse | ACCEPT (divergence 1) |

Mutation check of p0check (my own): neutering single clauses and running both arm sets. Against the designer's `mutate` arm set, mutants `a1-off`, `a3 gc>=0 (free gaps)`, `a4 min-height` and `a4 non-paragraph equality` all SURVIVE (V65 has no arm isolating them); with my arms they are killed (`N8` kills a1, `validation h+240` kills the free-gap gc>=0, `approvals-chain h=1` kills min-height, `graph/payload h+24` kill the else branch). The non-free `gc>=0` is redundant (R1 non-free gaps are 15 or 100, so |e|<=2 already forces gc>=13); harmless.

## (c) Legitimately different stores wrongly refused (constructible?)
- Nothing in the evidence lets me construct a refused-but-legit store beyond the rule's disclosed scope. The nine stores contain exactly ONE real non-baseline variation (eval-p1: PROVEN row 72 vs 96, the last row, so no gap rule is exercised); seeds 3-6, p2, p3 are identical to R1 to the pixel. So the whole-line machinery for gaps and mid-stack wraps rests on synthetics only (the doc's "Residuals" paragraph discloses this; good).
- A phone store whose overflow is so large that the thumb is < 100 px (content > ~1500 px wide) would be refused as "no thumb". Observed content widths are ~455-478 px, so unlikely; unmeasured.
- A wrap or new row ABOVE the stack head on the packet page (the object page renders the full timeline first, so `validation` top depends on every timeline entry above it) is refused by design and escalated; the release store is a different mint, so P3 may well hit this. Disclosed in §4.4/§8.

## Findings

### BLOCKING
None.

### NON-BLOCKING
1. Stale text (should be fixed before routing, trivial edits):
   - `grep -n "This is \*\*r" design_docs/planned/w-marketing-capture-orchestration.md` -> line 330: "This is **r4** ... quorum pending. It is the one permitted revision". The doc is r5, the r4 round has BLOCKED.
   - C13 (line 309): M8b "and the P0/P3 rules that r3/r4 amend, are gated on this document's r3/r4 quorum disposition". Nothing states what the r5 carve-out does to that gate (the Planner-Lane line only implies it via "the carve-out recorded in §9"). M8b's start condition is therefore ambiguous.
   - Header: "Gemini passed r1 and r2" (r3 round: Gemini REJECT; r4 round: Gemini PASS).
   - §5 R1 status: "Whether that discharges C14's C8' half is for the controller to record once r3's quorum disposition permits it".
   - §10 preamble: "design_docs/verification/world-iter315/ (V64-V71, plus quorum-r3.json)" omits V72 output and `quorum-r4.json`.
   - §8 line 322 and the §9 r3 table (line 384) still say "72 px against 96" unqualified; GLM's fix asked for the single stored value (97/73) to be cited in "the r3 table, §4.4 and V64's message". r5 fixed §4.4 and added V73 only; V73 and §4.4 do say the 72/96 figures are block-box readings, so the number is not ambiguous, but the fix is partly applied.
2. Row 231: §4.4 and the §9 r5 table say the committed checker "is filed as queue row 231". At c36a4306 the charter's last row is 230 (`grep -n -E '^23[01]\.' design_docs/world-mission.md` -> only 230 exists at line 6296). The commit touches only the doc and verification files. The claim becomes true only when the controller files it with the iteration record; make sure that happens, or reword to "to be filed".
3. Capture-day controls are weaker than they look: the P0 accept control is seed 5 = R1 itself (identity), and the phone accept control is V50 vs V50 (identity, thumb 329 = the reference). A checker that accepts only the baseline passes both. Kimi's own fix named seed-5 pass-A for phone (non-identity, thumb 333); the controller replaced it with V50, justified as "no PNG can be banked and /tmp is volatile", but V50 lives in `~/.ailang/state/` (also outside the repo) and is required as the checker's first argument anyway. Also only one refuse arm per checker, so a neutered (a1)/(a2)/(a4)/(b)/clip clause is not caught by the control (source sha pins do bound that). Use eval-p1 rects as a second accept control (non-identity).
4. P0 item (1) names the iter305 probe's "still half", but that probe is one monolithic test (`TestIter305SeededStillCaptureProbe`, needs `WORLD_CHROME` and `PROBE_OUT`, mints its own seeded store, then runs the Chrome half); it LOGS `repeat_equal=` and `script_tags=` (`t.Logf`) rather than asserting. It is runnable (compiles at HEAD) and the values are readable from the log, but no command line is given and "asserting byte-identical HTML" is not what the instrument does; also (2), (3) and (1) run on different disposable stores although (3) says "the same store". Spell out the command and that the controller reads `repeat_equal=true script_tags=0` per route.
5. V65 arm set leaves mutants alive (see mutation check). The rules are fine; the evidence for each clause being independently exercised is not. Add the four isolating arms above to the corpus when row 231 commits the checker.
6. Residual permissiveness inherent to the design (not contradicted by the text, but the doc's "wrong crop fails (a)" should say "other than whole lines in the two body-text gaps"): a region can slide by whole 24 px lines inside the two free gaps and still pass (`approvals-chain` +/-48, `validation` inflated by 216 px); per-element +/-2 tolerances compound (roster rows ~8-10 px off R1 accepted). r2's absolute +/-2 did not have this. The content gate and P3's bottom bound limit the harm.
7. Charter guardrail (world-mission.md ~1152): a fix applied under the carve-out must acquire its own AC and a named mutation before the doc routes. The §9 r5 table names none per fix. Partly implicit: C12 now names item (1); V72's accept/refuse controls play the mutation role for fix 1; V73 is a reading. State it explicitly.
8. Minor: the V50 sha256 and seed-5 rects sha are truncated or absent in §4.4 (only the two source hashes are given in full there); full values are in V72 and `v64-inputs-sha256.txt`.

## Carve-out honesty (item 3)
- Kimi (r4 REJECT): fix offered two routes. r5 takes the second ("keep the controller-program route") with Kimi's three parts: (1) sha256 per `.txt` verified before every build: applied in §4.4 with both full hashes, both verified by me; (2) mandatory control step, BLOCKED on mismatch before the verdict is read: applied, with `refuse=<attested+14>` and desktop-light exactly as Kimi wrote; adaptation: phone accept = V50 instead of seed-5 pass-A (named controller-authored, see finding 3); (3) a statement of why the committed checker is not extended: applied as a bullet (not a "subsection"), named controller-authored. Its reasons are checkable: D-WORLD-81 does order "finish 150, then 151" (`world-mission.md:1000`); "the iteration-315 brief told the designer to prefer no milestone change" is a mission brief, not a reviewer text, and is labelled as such. Row 231 is a controller-invented element, named as such, but not yet in the charter (finding 2).
- GLM (r4 REJECT): checklist applied essentially verbatim (retain branch); two controller additions, both small: the file path of the retained probe, and "exactly per the §4.4 invocations below" for "r4 invocations". C12 aligned. The height fix is applied in §4.4 and V73, partially in the r3 table / §8 (finding 1).
- Gemini (PASS) catch: r3 column and row-149 coupling were already in r4; V72 adds the exit-code fact. Correct.
- Controller-invented words in r5: the "Why the committed checker is not extended" paragraph, the V50 substitution, the row 231 filing, "(ratified at first use in iteration 95)". All but the last are named controller-authored. I could not confirm the iteration-95 ratification claim from the charter text I searched (the carve-out is cited in many places and the guardrail at line 1152 describes it as a routing rule), so treat that parenthetical as unverified.
- Is any applied fix a disputed design direction? Kimi's objection is about the conflict surface of an uncommitted checker, which touches design direction, but Kimi explicitly offered "if the controller-program route is retained instead" with three concrete additions, and r5 applies those. Nobody disputes the r3 rules. The carve-out is legitimate on this record, with the caveat that Kimi's PREFERRED route (commit the predicate) was not taken; that choice is the controller's and is disclosed.

## Scope (item 5)
No executor milestone is changed: the M8b row text in §5 is identical to before except for the "PENDING" label; r3/r4/r5 each say "M8b is unchanged". r4 rewrote M7c0-M8a status from PENDING to LANDED and replaced the §2/§3/§6 inventories (reviewer-requested status corrections, from git log; verified plausible), which is bigger than "amend two rules" but not a milestone change. Row 231 appears in two places (§4.4, §9 r5 table) consistently as the committed-checker follow-up, "not taken here".

## Worktree
Restored; `git status --porcelain | wc -l` = 0 (see final message).
