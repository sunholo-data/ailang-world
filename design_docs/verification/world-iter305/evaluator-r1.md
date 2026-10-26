# Evaluator r1, iteration 305 (sonnet; worktree .eval-world-iter305 @ f354e5e5)
Score: A 50/60 + B 31/40 = 81/100. Verdict: PASS (>=70). No BLOCKING findings.

## Part A: M7c1
Gates (mine): go test -count=1 ./scripts/marketingcheck ok (16.5s); -race ok (33s); go vet clean; gofmt -l empty;
TestProductionGoSurface (host/store) ok; TestSubprocessSinkGate (host/verifygate) ok. Diff touches exactly main.go + observation_test.go.
Check CLI preserved: runCheck(args[1:]) with fs.Parse(args); mutating either side red in TestMarketingTwoPassInventory.
Non-vacuity: with main.go at 4e6e89dc the test file fails to compile (undefined Observe/ObservationRecord), i.e. red.

Plan's three named mutants: all red on the named assertion
- M7c-observation (SHA256 check -> false): same-line-count-change "want Code=observation-digest, got <nil>"
- M7c-line-observation (Lines check -> false): digest-equal-line-change "want Code=observation-lines, got <nil>"
- M7c1-set (guard -> false): deleted "want observation-set, got observation-digest tools/launchd/testdata/deep/y.md: sha256 ... became " (no panic)

Diff-derived mutants (test = TestMarketingProtectedObservationComparator; K=killed, S=survived)
K glob *.ail->*.al (set-coverage); K drop go.mod / go.sum / mission.md (set-coverage); K launchd root changed (set-coverage)
K lines+1, lines counts \r, sha truncated (set-coverage + python-parity); K sort removed (set-coverage; needed compile fix first)
K empty-guard off (empty); K empty code renamed (empty); K slices.Equal off (deleted, added)
K digest-off (same-line-count-change); K lines-off (digest-equal-line-change); K digest refusal Name != path (same-line-count-change, cli)
K O_EXCL->O_TRUNC in writeObservation (cli (ii)); K observe root-required off, observe usage rc, observe io-error rc
K compare requires --after, usage rc, skip CompareObservations, refusal format swapped (cli)
K dispatch observe/compare swapped or given unsliced args; K check parse args[1:] / dispatch unsliced (TestMarketingTwoPassInventory)
S duplicate-guard removed ( || len(Compact)!=len ): the "duplicate" row only adds a dup on one side, so slices.Equal already refuses. Both sides carrying the same duplicate is untested. LOW.
S set guard reduced to len-only: a same-count swapped path (delete X, add Z) still refuses via zero-record digest but with the wrong code; untested. LOW.
S missing-path tolerance (world dir absent, launchd dir absent, stat ErrNotExist) removed: no test on a tree lacking these; spec says "a missing path contributes nothing" and the Python heredoc tolerates it. LOW-MEDIUM.
S IsRegular filter removed (dir named x.ail, symlink): equivalent for dirs in WalkDir path, but a dir named world/x.ail would make Observe error where Python skips. LOW.
S `[]` vs null: records := []ObservationRecord{} -> nil survives (observe of empty root printing null untested). LOW.
S order (lines check before digest when both differ) and lines refusal Name: refusal code/name with both differing untested. LOW.
S observation-io prefix, observe NArg / --help rc / parse-error rc, compare --before-required / NArg / help rc, run() no-arg and unknown-arg rc 2 (usage rc not asserted), out file mode 0600, trailing newline, close-error propagation. LOW (CLI minutiae; the spec lists rc 2 for none/unknown arg, which no test asserts).
Equivalent: afterByPath miss-skip (guarded by set check), launchd IsDir skip (covered by IsRegular).
Two plan-ladder mutants hit compile errors first (regular-check, closeErr) and were re-done; none counted as kills.
Net: ~35 mutants, 25 killed, 0 survivors on a core comparison branch (digest, lines, set equality, empty); survivors are edge-of-spec gaps.

## Part B: design doc r2
V56 re-run by me: rc 0, 9 top-level PASS (49 incl. subtests), no FAIL/SKIP; tool binary AILANG v0.54.15. Holds.
V57 re-derived: png.go:20-35 thresholds; png_test.go:359-360, 396, 438-441, 444-455 cites exact; exactly 3 Test funcs; TwoPassInventory (:548) pins no threshold. Holds. (Minor: :396 is the count==32 branch, the literal 32 is at :395 region; fine.)
V58 rebuilt in /tmp/ev305 and run on probe-seeded pass-a/b: output equals the recorded one (a, b PASS; timeline-only, blank, desktop-light FAIL; passes byte-identical, cmp confirmed). Discrimination: the wrong-crop arms are coarse. timeline-only fails because band COUNT drops below entryStart+6 (index-located entry), not because of position; blank fails on no bands. I added shift arms (content moved down 2/5/10 px): 2 px PASS, 5 and 10 px FAIL, so the +/-4 px rule does hold, but the doc never exercises it; Gemini's ">4 px shift" catch is only claimed, not demonstrated. Baseline is pass-a itself, so pass-b agreement is expected (doc admits this).
V42 viewed: packet.png (1440x1400) = header + timeline entries 0-4 only, no packet/roster; graph.png (1440x1100) = timeline (entry 9) + inspector with no world graph. Claim holds.
Authority honesty: Status, Planner-Lane, plain statement, C13, section 9 all say "carve-out, not a quorum PASS"; Gemini passing r1/r2 noted. No false PASS claim found.

Carve-out conditions
(a) Each r2 blocking objection (GLM, Kimi) had a concrete proposed_fix: held. Gemini's (pass) too.
(b) None disputes direction: held, with a caveat below. Kimi's is a self-contradiction in acceptance text, not a direction dispute.
Fidelity per edit
- GLM 1 (V56, annotate rows, weaken 1.4): faithful, command run as quoted (controller added AILANG_BIN, sensible).
- GLM 3 (copy quorum JSON, cite): faithful.
- GLM 2 (V57, name test update, owner, gate): faithful to the finding; the OWNER (M8a executor) and gate wording are controller choices. NON-BLOCKING judgment content.
- GLM 4 (reword 4.4, delete 4.1 implication): reviewer gave direction, controller authored the text. Accurate to V43. Minor.
- Kimi C8' / R1 phone bullet: C8' text is Kimi's verbatim; the V56->V58 renumbering is controller judgment (disclosed in table row 5). C12 "matching edit" authored by controller. The V58 scan itself (luma 110, band index 9..14, std-dev 4.0, scrollbar heuristics) is controller-designed, so "applied in the reviewers' own words" slightly overstates; Kimi's (iii) "identical versus the V50 reference image" is not a stored image (V50 was visual inspection) so (iii) is checked against pass-a.
- Gemini catch: doc says "met by V58's +/-4 px rule"; the proposed fix asked for explicit tightness wording in 5 R1, which R1 does not restate. Partial.
Conclusion: the carve-out is defensible and honestly recorded, but Kimi's fix needed controller design work, nearer the "controller judgment" line than the rule intends.

## Findings
BLOCKING: none.
NON-BLOCKING
1. observation_test.go: missing-path tolerance untested (mutants world-notexist / walk-notexist / stat-notexist survive). MEDIUM-LOW.
2. observation_test.go "duplicate" row: duplicate-guard clause killable only with a both-sides duplicate; len-only guard survives. LOW.
3. observation_test.go cli: usage rc 2 for no/unknown args, observe/compare --help rc, observation-io prefix, out mode, trailing newline unasserted. LOW.
4. main.go Observe: IsRegular on a directory named *.ail diverges from Python (error vs skip) and is untested. LOW.
5. w-marketing-capture-orchestration.md section 9 r2 table row 5/6: V58 scan and "reviewers' own words" overstate verbatim application; Gemini fix only referenced, not restated in R1; V58 has no >4 px shift arm (my shift arms show it works).
6. V58 (iii) reference is pass-a, not a stored V50 image.
