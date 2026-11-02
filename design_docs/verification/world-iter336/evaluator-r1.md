# Evaluator r1 - row 238 w-marketing-predicate-tolerances (iteration 336)

Judged at 95179277 (diff 37100a43..95179277). Score 89/100, PASS (70 to pass).
Breakdown: correctness 33/35, test strength 19/25, spec fidelity 14/15, doc/code agreement 14/15, hygiene 9/10.

## Gates (measured)
- Targeted: `WORLD_PHONE_CORPUS=... go test -count=1 -v -run '^TestMarketing(Position|Phone)'` rc 0, SKIP 0, 146 PASS lines.
  PARITY position inputs=51 verdict-diffs=0 reason-diffs=0 ledger=8; PARITY phone inputs=45 diffs=0 ledger=4 and inputs=63 diffs=0 ledger=4; CHECK images=10 byte-equal=10.
- go vet clean; gofmt -l empty; go build ./... rc 0.
- Full `AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_PHONE_CORPUS=... go test -count=1 ./...` rc 0, 35 ok packages (marketingcheck 126 s).
- world-iter315 oracle sources: `git diff 37100a43..95179277 -- design_docs/verification/world-iter315/` is empty.
- Scope: changed non-evidence files are exactly position/phone(+_test).go, the two planned docs, sprint plan, sprint JSON, orchestration doc. Evidence dir 480K, 41 txt + 4 json, none >200K, no binaries.
- Real data: position over the 11 store files gives 9 ACCEPT, raw r1-seed1/r1-seed2 REFUSE (TOTAL accept=9 refuse=2). Phone over the 10 real images: 10 ACCEPT, thumbs unchanged (324..340, V50 329), titleRight 389.

## Mutation testing (against M1 diff; files restored, shasum verified OK after each)
KILLED (named subtests): free-gap bound 1->2 (bnd-createdAt-wraps-2-lines, perm-approvals-y+-48), 1->0 (iso-minh-approvals-h-1, perm-approvals-y+-24), drift tol 2->3 (bnd-drift-sum-3, bnd-nonpara-drift-3, Parity), 2->1 (bnd-drift-sum-2, bnd-nonpara-drift-2), (a5) gap-term sign flip, prevRes sign flip, gap term not accumulated, head residual dropped, okGap condition removed (13 reds incl. r1-seed1), hOK condition removed (this one did not compile: unused var, n/a), drift reset per region, paragraph residual raw dh, free gap residual raw e, signed-lines bound (positive only), negative-only bound, own-height residual added (current region), bridge gap ==2 -> ==3 / <=3 / ==1 / off / end=segs[i+1][0], chained second bridge, tie-break rightmost (B7), threshold 100->101 / 99, longest-selection collapsed. Phone parity and real-corpus tests also red on bridge mutants.
Parity ledger: tolerance mutants `posAbs(c.Y-b.Y) > 3` reds both Corpus and Parity; drift tol and bound mutants red Parity. Ledger entries each pin oracle verdict, position verdict and the exact added reason lines (position_test.go:594-613; phone ledger [3]string).

## SURVIVING mutants (non-blocking test gaps)
1. position.go `if posAbs(drift) > positionDriftTol` -> `if drift > positionDriftTol`: SURVIVES. No arm has cumulative NEGATIVE drift (every (a5) arm is positive: +3, +6, +14). A candidate drifting -3 in total (each term within its per-term bound) would be wrongly ACCEPTED if the abs were lost.
2. position.go posLines `math.Round` -> `math.Floor`: SURVIVES. No arm has a small negative paragraph/free-gap residual (e.g. dh=-1 or e=-1), so rounding mode is unpinned.
3. position.go (a4) non-paragraph `posAbs(dh) > 2` -> `> 3`: SURVIVES (pre-existing row-231 code; no +-3 height arm; also masked by (a5) giving REFUSE anyway).
4. position.go non-free (a3) `posAbs(e) <= 2` -> `<= 3`: SURVIVES (same: no +-3 gap arm, masked by (a5)).
Equivalent / not real survivors (analysed): `chain=false` in the missing branch (next region's pb-missing branch sets it), non-free gap residual rounding (|e|<=2 gives 0 lines), `gc >= 0` on free gaps (gb 239/95 minus 24 stays positive), (a3) line bound on non-free gaps (lines==0 whenever okGap).
Plan-unlisted mutants that DID die: (a5) sign errors, own-height residual added, bridge width off-by-one (==3, <=3, ==1), bound signed-only/negative-only.

## BLOCKING findings
None.

## NON-BLOCKING findings
1. Surviving mutants 1-4 above. Suggested follow-up arms: a drift of -3 (two -2 steps), dh=-1 / e=-1 residuals on packet paragraphs/free gaps, +-3 height and gap boundaries.
2. Orchestration doc says V81-V83 raw output is in controller-orch-r7-measurements.txt - exists. I re-ran V82 (0 files) and V83 (12 files, all design_docs/ or QUICKSTART) - agree.
3. The r7 carve-out makes P0 depend on a future judge content review of row 231's 16 pasted lines (C13); that is a downstream gate not satisfied by this evaluation.

## Doc/code agreement
- (a3): doc "one whole 24 px line up or down, +-2 px" == code posLineStep && |posLines|<=1 on exactly approvals-chain and roster-tested (posFreeGap). OK.
- (a5): doc residual definitions (paragraph dh-24*round, other dh; gap e-24*round for the two designated gaps, else e; only passed terms; head included; missing region ends walk) match position.go (chain flag, hOK, okGap). Bound +-2 == positionDriftTol. OK.
- (iii'): doc "one contiguous segment or two separated by exactly one off-range pixel, at most one bridge, first to last in-range pixel, X1 exclusive" == phone.go (gap==2 index delta, single pairing, leftmost longest). OK.
- Numbers: 8 position / 4 phone ledger entries, 51 / 45 / 63 inputs, 10 images agree with test output. V79 cites `phone.go:94` `== 2`: still line 94.

## Process claims
- quorum-r1: kimi reject, sol reject, gemini pass -> blocked; quorum-r2: same; orchestration r1 and r2: all three reject, none absent. Doc §9/V80/V81 verdicts and costs match the JSON.
- "No objection disputed row 238's predicates": fair for both orchestration rounds (objections are provenance/baseline/review-record; gemini r2 asks for a stale opener fix and a renderer-drift check). The design-round objections (sol: (a5) wording ambiguity; gemini r2: missing-region note) did touch predicate wording but were fixed in design r3, and the doc text now has "dh otherwise" and "a missing region ends the walk".
- Every r7 fix in the dispositions table has a corresponding artifact: §9 opener, V82, V83, V81 row, C13 precondition.
