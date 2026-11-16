# Judge report, iteration 393, row 260 w-ailang-v0604-contract-eq (sprint head 9e64b719)

Binaries confirmed by --version: v0.60.4 (563c6c9), v0.41.0 (24ee108), v0.54.11 (35a49b8), v0.54.15.

## Acceptance re-derived
- verify_ail.sh, AILANG_BIN=v0.60.4, WORLD_PKG_AILANG_BIN=v0.41.0: rc=0, 16 identities, 40 named tests, 503 se-tools tests, world package gate 9/9, PUB011 uncontracted_exports=3.
- verify_ail.sh, AILANG_BIN=v0.41.0: rc=0, same counts.
- test_kernel_contracts.py (v0.41.0): rc=0.
- go test -count=1 with AILANG_BIN=v0.41.0: cmd/world-publish, host/pins, host/pkgproj, host/runbook, host/broker, host/daemon, cmd/ailang-worldd, host/replay all ok. host/verifygate red ONCE under 8-package parallel load (TestMissionConfigDefaultFlowIntegration, "default-flow fixture timed out", 10 s deadline); green alone in 8.4 s. Load flake, not a sprint defect.
- Baseline red: first run without AILANG_BIN gave many "AILANG_BIN unset" reds (my env error, the gates refuse to skip). With the pin there is no deterministic red at the head. Base cbcf80fb was not separately run; the only red seen is the load flake above, which a changed-byte diff cannot cause.

## Claims
- ensures restatement: HashRef is `{ algo: string, digest: string }` (world/logepoch.ail:13), exactly two fields, so `.algo == && .digest ==` for each HashRef is logically identical to the original `==`; the plan/verify law is unchanged and the idiom matches applyRevision. The `deriving (Eq)` additions are purely additive. Contracts 16 verified / 0 refuted on both compilers.
- Regeneration: ./scripts/build_world_package.sh yields byte-identical packaged copies (git status clean). Fresh `publish --dry-run` with v0.41.0 and v0.54.11 are byte-identical to the two committed fixtures. Golden and runbook digest rows are bound by verify_world_package step 9 and host/runbook, both green, and mutated rows are killed. interfaceHash (d16cc882...) and interfaceHashV2 (b25fe031...) unchanged in diff; contentHash 359a82fd... is what the dry-run prints.
- Nothing published: registry packages/world/core/0.1.2/metadata.json returns 404 (0.1.1 returns 200); version 0.1.2 is new.
- PR #565: headRefOid 44e6c8a9, updatedAt 2026-11-12T21:44:47Z, untouched.

## Scope
The diff matches design section 3 items 1-16 exactly (18 non-design files). Nothing outside the list. The two deriving lines are credited to #565 in the M1 commit body.

## Mutants
33 total incl. controls: 27 killed, 0 real survivors, 4 survived as equivalent or unpinned prose (m23 equivalent, m28 comment, m26 and m27 docs prose), 2 controls survived as required. Every tree reverted; git status --porcelain empty.

## Non-blocking findings
1. docs/SELF_MOD_PUBLISH.md confirmation-phrase and fresh-store-name prose (m26, m27) are not pinned by any gate; they could drift back to 0.1.1 unnoticed (fences_test pins the code phrase, not the doc). Documentation-only.
2. test_kernel_contracts.py does not catch a consistent weakening of an ensures in both copies; the golden does (m30), so the bytes are still guarded.
3. host/verifygate TestMissionConfigDefaultFlowIntegration has a 10 s fixture deadline that flakes when run in parallel with the other heavy packages.

## Score
Correctness 38/40, Tests 18/20, Security 10/10 (no publish, fences moved with the version, interface hash constant intact), Standards 14/15, Scope 14/15 = 94/100. PASS.
