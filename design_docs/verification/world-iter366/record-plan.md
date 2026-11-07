# WORLD-ITER366-ADMIN — record plan

Design: [record-design.md](record-design.md), G1–G8; immutable planning receipt: [record-state.json](record-state.json).
One milestone: **M1_ADMIN_RECORD_366 — prepare and verify the mandatory administrative record**.
Estimate: one short executor session (30–60 minutes); zero product/test LOC. Main risk: accidental history or authority drift.
Registry reuse: **none — administrative bookkeeping**, no package capability or registry access.
Controller's explicit ADMIN quorum exception applies; no human decision is awaited for this record.
Planner checkout HEAD: `adb4fb80b9d7600ad9118aa15a29fb2e4a4e7c89` (detached, includes committed inputs).
Canonical product base: `1ad7c8fdc80bdc6e1d5532fb17eb6b7406809dd4`; never substitute planner HEAD for this baseline.

## Authority to carry unchanged

Row256 remains **PARKED needs-human-review; independent MiniMax FAIL 85; LOCAL_CRITERION UNMET; draft PR #515**.
D-WORLD-105/110/111/112 remain OPEN: defaults C (233 parked/fail closed), A (255 as built), A (exactly the two named baseline failures, no additional failures), and no autonomous credential change.
D111 and D112 are independent human holds; sanitization does not undo exposure or resolve112. Publish-counter/net6 remain unattributed; nonrecurrence is not attribution.
Conditional order: **250 → 255 → 256 → 257 → 201 → 243**; no autonomous skip256. Release150/151 take priority only when routable; attended A1 and188 N2d/N4a remain waits.
Rows201/243 use fakes/test keys; real-key exits remain attended. Clauses1–7 MET and goal unchanged are inherited assessments.
Recommended record wording: **Progress**: 1.0: clauses 1–7 MET (inherited assessment); goal unmoved.

## Exact ownership and work sequence

Only these five existing administrative outputs may change:

1. `design_docs/world-mission.md`: insert one newest STATUS366 line+blank, retain365/364, remove exact363 line+blank for transfer; replace only the unique `Current route (iteration 365)` with366. No other charter edit; final6,345 lines.
2. `design_docs/world-mission-status-archive.md`: preserve every old prefix byte; append the exact363 pair (209 bytes, SHA256 `47c176790ec63fad66aeaf56ba117b3fe2f09d8ab2e07946e32616f349123afe`);730→732 lines.
3. `design_docs/world-mission-log.md`: preserve complete prior prefix; append one truthful366 entry headed exactly `## 366 — 2026-11-07 — ADMIN record only; row 256 remains parked [ADMIN]`. Entries39→40: **NO rotation; rotate only above40**.
4. `design_docs/world-mission-index.md`: insert exactly `| 366 | 2026-11-07 | ADMIN record only; row 256 remains parked [ADMIN] |` before365, preserving every earlier byte/row/order including361.
5. `design_docs/world-mission-dashboard.md`: concise truthful366 summary, at most40 physical lines, retaining parks, defaults, conditional order and evidence limits.

Executor-owned NEW outputs are `design_docs/verification/world-iter366/executor-proof.json` AND `design_docs/verification/world-iter366/execution-state.json`; require both and `.snap/M1/` absent before creation. Never overwrite another role's artifact.
Original record-state.json is an IMMUTABLE planning receipt: milestone passes=false forever; no executor progress-field edits. Design, plan, state metadata and input map stay immutable.
Retain the explicit freeze map in state. Controller freezes the consumed original planning-state bytes at execution launch using an external hash receipt; no cyclic hash in state.
First verify individually frozen input hashes against state; the five editable files are pre-edit baselines, subsequently verified by exact expected transformations. All other frozen inputs remain byte-exact.
Parse base-files.json programmatically; use its6,442 baseline paths/hashes for preservation, never print the full manifest or earlier evidence collections.
Then edit the five outputs, run only the short read-only comparisons below, and write executor-proof.json with actual commands/comparisons, return codes, diagnostics, positive controls, findings and provenance by AC name. At M1 completion create execution-state.json with actual preparation passes, preparation completion timestamp and notes: distinguish prepared G1–G7 from pending controller-owned G8 independent review/CI/delivery; preparation passes never claims G8 completion.
Use existing tools and short inline byte comparisons; no novel long scripts, acceptance diagnostic gates, suites, product source/tests/gates, shared state, credentials, network, external messaging or nested agents. No TMPDIR/GOTMPDIR overrides or seeded guards.
All baseline files outside the five, log archive, ALL earlier evidence and every existing366 artifact remain exact. Allow genuinely NEW concurrent controller evidence untouched; never freeze/hash/copy the whole366 directory.
Raw worktree-census.txt is local-only, reconstructed from tracked worktree-census-original.json original_text with original_sha256 verification. Preserve original bytes and original whitespace failure; no repair or waiver.

## Acceptance criteria (all pending; names also appear in state)

- **AC-G1-SCOPE**: read-only Git census confirms the two full base identities above and records actual executor HEAD separately; canonical baseline census nonzero,6,442 paths and5/5 allowlist controls. Audit tracked diff and untracked paths/types separately; only enumerated owned writes allowed, with existing local census and NEW untouched controller evidence recognized.
- **AC-G2-TRANSFORM**: construct expected charter/archive/index bytes in memory from canonical base and compare whole outputs. STATUS366/365/364 present, exact363 pair transferred once, route365→366 once; archive old prefix and every old index byte preserved. Entire marked ledger (start through end marker plus following newline) and unique Queue-header-through-EOF suffix remain byte-exact (designer O7 hashes).
- **AC-G3-COUNTS**: log has full old prefix and exactly one canonical366 heading,40 entries with no rotation; log archive unchanged. Charter6,345 lines; status archive732; dashboard≤40 physical lines.
- **AC-G4-PRESERVE**: recheck all individually frozen immutable inputs and every protected baseline path/hash/type/mode; input drift is a finding. Verify editable baseline hashes before edits and exact expected outputs afterward. Preserve all earlier evidence; snapshot and controller reconstruction match the explicit payload manifest byte-for-byte.
- **AC-G5-TRUTH**: record actual role/model provenance, degradation, runtime/token/cost usage or explicit unknowns, retro, unchanged goal/asks/next. Attribute Gate0/1, clauses and base CI green to inherited controller receipts; no invented role probe, acceptance, landing or product completion. ADMIN success never accepts256 or clears111/112.
- **AC-G6-WHITESPACE**: run both `git diff --check` and `git diff --check 1ad7c8fdc80bdc6e1d5532fb17eb6b7406809dd4 --`; explicitly inspect NEW owned text outputs using `git diff --no-index --check /dev/null <path>`. Record rc and diagnostics separately: new-file rc1 denotes content differences, not automatically a whitespace error. Pair empty diagnostics with nonempty path controls. Report inherited base reds without repair or waiver; whole old archive whitespace style is not an acceptance gate.
- **AC-G7-BOUNDARY**: no child test verdict or Go suite acceptance. Inherited pristine baseline charter rc0 is provenance only. Any required `go test ./host/verifygate -run '^TestMissionConfigRealCharter$' -count=1` belongs to controller outside sandbox, only if authorized after side-effect inspection and pristine baseline. Sandbox denial is **UNINFORMATIVE UNDER SANDBOX**, never a product verdict.
- **AC-G8-DELIVERY**: controller obtains independent final-record judgment covering all material changes, reconstructs/commits from snapshot, and verifies all3 terminal-success CI jobs at exact full40 PR-head AND merge-head with run URLs: ailang-code verify, Go race, Go host build/test. Controller owns external delivery and landing. Prior/base/product greens do not substitute; no future completion claimed.

## Cumulative M1 handoff

After preparing the record, proof and execution state, snapshot final full bytes of exactly the five administrative paths above plus executor-proof.json AND execution-state.json into `.snap/M1/`, mirroring their repository-relative paths.
Create exactly one `.snap/M1/manifest.json`: one entry per seven payload files with repository path, snapshot path, type, mode and SHA256. Manifest excludes itself to avoid a cyclic hash; no other payload is allowed.
Exclude consumed original design/plan/planning-state, all receipts/controller files, raw census and earlier evidence. Verify all seven snapshot copies against live payload before handoff; controller verifies again after reconstruction.
No child Git writes (including add/commit/stash/checkout); controller reconstruction commit plus manifest suffices. Bulky `.snap/M1/` remains local-only, never tracked. Original planning milestone passes=false forever; execution preparation does not claim delivery.
Hand the prepared snapshot/proof/execution-state and consumed plan/state directly to controller; no user approval, automatic external handoff or further product action is requested.
