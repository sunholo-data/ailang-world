# Evaluator r1 - row 150 iter 306 (M7c0 + R0 V59)
Worktree /Users/voightkampff/dev/sunholo-data/.eval-world-iter306 @5eb7be62. Score 96/100. Verdict PASS.

## Scope
`git diff eacd90e7..5eb7be62 --stat`: exactly the 5 listed files (orchestration doc +1 row, two r0 output txts, png.go +2/-2, png_test.go +48/-1). Nothing else changed. png.go: only the 3 heights (packet 3700, graph 2100, provenance 3000); light/dark/phone untouched. png_test.go: "math" import, canonical list updated at the 3 entries, constants + shotRuleHeight + TestMarketingShotGeometry, byte-for-byte the plan's code except the V-row cite filled in as V59. No do-not-touch lines changed.

## B. Acceptance commands (my worktree)
- gofmt -l scripts/marketingcheck: no output. go vet: rc 0.
- TestMarketingShotGeometry -v: PASS packet.png, graph.png, provenance.png, unchanged; ok.
- TestMarketingPNGContent -v: PASS incl. canonical-shots, positive-packet.png/-graph.png/-provenance.png; ok.
- Whole package -v: rc 0, 149 "--- PASS", no FAIL, ok 18.15 s.
- -race ./scripts/marketingcheck, 3 runs: run1 FAIL TestMarketingSupervisorFakeDaemonCleanup/interrupted-shot ("[Errno 1] Operation not permitted"); run2 ok 45 s; run3 FAIL TestMarketingSupervisorFakeChromeTimeoutKillsGroup. Different supervisor tests fail on different runs, none touch png.go/Shots. The same interrupted-shot failure also appeared once in a non-race run with png.go reverted to base, so it is environmental/load (the plan's carried flake class), not caused by M7c0. NON-BLOCKING, but "-race rc 0" is not reproducibly obtainable here.

## C. Non-vacuity (png.go-only mutants; each restored by cp and cmp byte-identical)
- Revert png.go to eacd90e7: TestMarketingShotGeometry/packet,graph,provenance RED and TestMarketingPNGContent/canonical-shots RED (plus the unrelated env flake).
- M7c0-height (packet 3700->1400): `png_test.go:834: packet height must equal the rule value from measuredPacketBottom: got {Name:packet.png Width:1440 Height:1400 Scheme:light}, want 1440x3700` - matches spec exactly; canonical-shots also reds (`canonical shot packet.png: got ... want ...`).
- graph 2100->2200 and 2100->2000: KILLED (ShotGeometry/graph + canonical). provenance 3000->3100: KILLED. packet 3700->3600: KILLED. packet width 1440->1400: KILLED.
- dark height 1100->1200: KILLED (ShotGeometry/unchanged, canonical, light-scheme, dark-scheme). light height 1100->1000: KILLED (unchanged, canonical, uniform, mostly-gray, ...). phone 844->900: KILLED (unchanged, canonical). dark scheme->light: KILLED.
- graph scheme light->dark: killed ONLY by canonical-shots (geometry test does not compare Scheme for tall shots). Acceptable; the canonical list covers it.
- Test-side reasoning (not run): rule +40 margin change or a wrong constant would red the equality because heights are pinned by value in canonical-shots as well.
All 11 production mutants killed.

## D. R0 record V59
Re-derived from the two output files: seed1 manifest sha 0a1ef461..., head 8c4b8aff...; seed2 d2fa161c..., a90ced43...; approval-entry=9 both; B = 3564.56 / 2042.14 / 2881.17 in both, delta 0.00; margins 135.44 / 57.86 / 118.83; blocked=false; SHOT verdict ok x3 per seed; PASS 11.15 s / 16.65 s; rc 0. Row matches the files. Phone scrollW 689/698 in the files is mentioned as unusable instrument (consistent with V51's 701). The PNG byte sizes and the "viewed bottom 900 px" claim are not verifiable from the committed files (PNGs not committed) - NON-BLOCKING, controller-attested. Chrome --version here = 156.0.8078.48, matches the row.
My independent run (seed3, rc 0, PASS 13.91 s, PROBE_OUT eval-seed3): manifest sha 91219752...; BOTTOM packet 3564.56 -> 3700; graph 2042.14 -> 2100; provenance 2881.17 -> 3000; CONFIRM deltas 0.00, margins 135.44/57.86/118.83; SHOT ok x3; light need 613.05; phone scrollW 700; SEED head sha256:278bfef6e0c27e9958492817c0244b60883bf0ef8a020b63f7c68f1926fa79b1. Probe file deleted; git status clean.
Probe discrimination: it uses querySelector (first match) + getBoundingClientRect().bottom injected via a loopback proxy at scroll 0, measured at 4400 and again at the rule height with a delta; a wrong bottom could arise from (a) a selector matching a different/earlier element (e.g. first pre.payload), (b) a missing element: b(nil) returns -1 so rule gives height 0, which is not flagged as blocked and would only show as bottom=-1 in the log - a latent hole, (c) measuring the section box rather than the last visible content (section bottom is conservative: includes padding), (d) dump-dom timing before layout settles (deferred script; the CONFIRM rerun and 3 independent seeds agree to 0.00). The 4400 vs rule-height comparison makes viewport-dependent layout detectable. Adequate for this use; V52 pixel check backs it.

## E. Other observations (NON-BLOCKING)
1. In shipped tests, the `want > 4400` branch and the `bottom > H-40` branch are never exercised with current constants (3700 max; rule guarantees bottom+40<=H). They are guards, not tested arms; ±40 is actually enforced only through the equality with the rule. Fine per spec, but a small table-driven test of shotRuleHeight at boundaries (e.g. 3560.00 -> 3600, 3560.01 -> 3700, 4360.01 -> 4500 > cap) would pin the rule itself, since a rule-function mutant (e.g. +40 -> +30) is only caught because canonical-shots hard-codes the heights.
2. The constants comment has one very long line (about 140 chars: "...(iteration 306, two fresh seeds on Chrome 156.0.8078.47->.48, delta 0.00)"); gofmt-clean, cosmetic. The comment says "Chrome 156.0.8078.48" which is correct.
3. Planned doc: the M7c0 plan section status still says "pending"; the controller owns that bookkeeping.
4. Dead-but-spec-mandated `ok`-only check: `!ok || ...` Fatalf prints zero Shot on missing name; fine.

## Scoring
Tests pass 19/20 (package green; -race flaky on unrelated supervisor tests), lint 10/10, acceptance criteria 30/30, code quality 14/15, documentation (V59 row, citations, constants) 14/15, design fidelity 10/10, conditional categories n/a (no internal/ paths). Total 97 -> reported 96 after the unreproducible -race rc. PASS (>=70), no BLOCKING findings.

EVALUATION_RESULT: pass
EVALUATION_SCORE: 96/100
EVALUATION_ROUND: 1
