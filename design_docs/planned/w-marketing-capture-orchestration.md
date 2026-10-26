# W-Marketing-Capture-Orchestration — the rest of row 150, written against what is built

**Status:** r2 (iteration 305). r0 designer, r1 designer revision after quorum r1 BLOCKED, r2 **controller narrow-refinement carve-out** after quorum r2 BLOCKED: every remaining objection carried a concrete reviewer fix and none disputed direction, so the fixes are applied in the reviewers' own terms (§9) and the doc routes to the planner without a third quorum. No round of this document has passed with all seats; Gemini passed r1 and r2.
**Target:** World v1.0.0, queue row 150 / D-WORLD-64, D-WORLD-69 (A), D-WORLD-75 (A), D-WORLD-81 (A)
**Priority:** P1 release · **Estimated:** 8 executor runs (≈196 executor-minutes from the plan's per-run estimates). After those come one controller rehearsal and one attended release capture. See §6 for the inventory.
**Created:** 2026-10-26 (iteration 305)
**Base:** `4e6e89dc44c9ee9771b415b3cfdc23fadfaadff1` (origin/dev). Every row in §10 was measured at this commit. Parent r4 base `c21f303` and iteration-304 HEAD `039bfa77` are both ancestors (V32).
**Planner-Lane:** mission routing default. It applies **only to the pending milestones in §5**, and only after this document's own quorum disposition permits it (a quorum PASS, or the Gate-2 narrow-refinement carve-out recorded in §9; r2 used the carve-out, not a PASS). This document does not authorize, and never authorized, any landed code (§2).
**Dependencies:** M1a–M7b landed (§3, V33). Tool pin is AILANG v0.54.15 (V36). Proof/interpreter pin is v0.41.0. System Chrome must be major ≥155 (V47). Fleet PNG allowance `2d14dbd` (parent V2). Attended descriptor publish and attended `session mint` are needed for the release store.
**Parent:** [`w-marketing-capture.md`](w-marketing-capture.md). It keeps the landed M1–M7b design (D1–D3, the seed sequence S1–S6, the PNG gate §3). This document does not restate those. It changes only what §4 says it changes.

## 1. Problem statement

Row 150 has built 14 of its milestones. These are the demo seed, the evidence roster, the still flag, the bounded supervisor and the PNG verifier. What is left is the orchestration:
- the shell that serves the frozen store and makes two six-shot passes;
- the frozen-inventory and protected-path guards;
- the real-Chrome drill;
- the docs;
- the attended release capture.

The parent design was re-quorumed in iteration 304 and BLOCKED twice. The mission's rules sent the remaining work to a new document anchored at HEAD (SPLIT). That document carries the r6 objections verbatim. They are reproduced below exactly as recorded in `design_docs/verification/world-iter304/quorum-r6.json` (`jq '.reviewers[] | {model, v: .result}'`).

### 1.1 GLM 5.3 (`oc-glm-5-3`), r6 REJECT, verbatim

> **strongest_objection:** The document's verification basis is split across two divergent commits without re-verification of its load-bearing premises at the state that will actually ship: rows V1–V28 were measured at/near base c21f303 at iteration 266, rows V29–V31 at HEAD 039bfa77, and the doc's own iteration-304 annotations state that M1–M6 and M7a/M7b have already landed — meaning the renderer, packet-evidence projection, and serve path that D2/D3 depend on have been modified by the doc's own milestones since V4 (grade resolution), V6 (routes/payload flag) and especially V14/V24 (exactly one script tag, no inline script, no nonce, empty varying-byte inventory) were measured. No verification row re-executes any of these at HEAD. The central still-mode determinism premise of D2 therefore rests on measurements of a pre-implementation renderer, and the doc simultaneously claims 'This revision changes design_docs only' and 'Proposed implementation paths below are future work, not files changed by this revision' while annotating that they are already changed — an internally falsified self-description that an executor or controller cannot safely act on.
>
> **catch:** Two specific verification gaps to fix before proceeding: (1) the blanket annotation 'M1–M6 ... have since been built' is asserted, not verified — only M1a (28ff4db3), M5a (bcd2d7cd), M7a (38ac4167/98f48b20) and M7b (fa04df00/02e12be2) carry SHAs; M2, M3, M4 and above all M6's --workbench-still flag have no executed row, no SHAs, and no cited C1–C6 acceptance results; (2) the §2 capture recipe's real-content byte-equality premise (V25) was measured on Chrome 155.0.8059.41, which V29 confirms is no longer installed (rig auto-updated to 156; only a synthetic data: page was measured on 156), so the product-behavior premise on the instrument that will actually be used is unmeasured and only covered by the still-pending M7/V21 gate.
>
> **proposed_fix:** Re-anchor the document as an as-built revision at HEAD 039bfa77 instead of a future-work plan frozen at c21f303: (a) update the Base/branch line and strike or correct the 'changes design_docs only' non-goal and the 'future work' framing of the milestone table, converting each M-row to landed/pending status with its commit SHA and executed acceptance evidence (test names/commands) or an explicit 'not yet verified' mark; (b) add verification rows executed at HEAD: V32 re-running the V24 probe against the real --workbench-still flag (serve the seeded or test-helper store, two GETs 2s apart on /workbench, ?object=<ref>, ?object=<ref>&payload=1, ?from&entry routes: assert 0 script tags and byte-identical HTML, plus default-mode script presence), V33 re-reading V4/V6 at HEAD (grade resolution now including the authored-claim case at workbench_packet_evidence.go:265-267, route and payload=1 behavior) to confirm the iteration-266 premises survived M2–M6; (c) replace the stale 'eight executor runs, each ≤30 min' estimate with the actual landed/pending run inventory; (d) optionally add a V34 row re-running the V25 real-content two-pass probe once on the current Chrome 156 before the attended release capture, so the §2 recipe's determinism premise is measured on the installed instrument rather than inferred from a superseded version.

### 1.2 Kimi K3 (`oc-kimi-k3`), r6 REJECT, verbatim

> **strongest_objection:** The doc's own gating claim is falsified by its own annotations. The header asserts "Planner-Lane: codex-ok for the milestones below after design quorum/freeze," yet the iteration-304 annotations record that M1a (28ff4db3), M5a (bcd2d7cd), M7a (38ac4167, PR #360) and M7b (fa04df00/02e12be2, PR #361) have already landed while every quorum round on record — r1, r2, r3, and the latest r5 (iteration 304: GLM REJECT, Kimi REJECT) — is BLOCKED, with no passing round anywhere. V30's verification is measurement of already-merged code, so quorum review has become retrospective decoration: the design is documenting a fait accompli instead of gating one. Either the "after design quorum/freeze" precondition was breached by those merges, or it has been silently voided; the doc annotates around this contradiction rather than resolving it. Secondary concrete gap: all real-content two-pass byte-stability evidence (V25) was measured on Chrome 155, which V29 confirms is no longer installed on the rig; the only Chrome-156 evidence is a synthetic data: page. Section 2's title "measured on 155 and 156" therefore overstates the instrument evidence for the only browser that can actually run M7 — V21 (the real product measurement) remains PENDING, so M7's feasibility premise currently rests on zero real-content measurements on any installed browser.
>
> **catch:** Three things to fix or verify: (1) Reconcile the Planner-Lane line with the landed M1a/M5a/M7a/M7b merges — name the authority under which code merged ahead of a passing quorum, or revert the claim that codex-ok applies only after freeze. (2) Re-measure the V25-style two-pass real-content capture on the currently installed Chrome 156.0.8078.47 before calling the capture recipe "measured" for the rig that will run the release; a 390x300 data: page with one prefers-color-scheme rule does not establish byte-stability of the multi-pane workbench HTML. (3) The document is labeled "REVISED r4 (iteration 266)" but contains and dispositions an iteration-304 r5 round — the self-identification is stale; this is an r6-candidate and should say so, with r5's dispositions clearly pending re-review, not retroactively "fixed" by the same party that authored the revision.
>
> **proposed_fix:** Replace the header Planner-Lane sentence with an explicit provenance statement, e.g.: "Planner-Lane: M1a (28ff4db3), M5a (bcd2d7cd), M7a (38ac4167, PR #360), M7b (fa04df00/02e12be2, PR #361) landed under sprint-plan authority <name row/decision record and date> prior to design quorum; all remaining milestones (M7c, M8, and the attended V21/M7 capture) are blocked until an r6 quorum PASS. This document does not itself authorize the landed code." Add verification-log row V32: `git log --oneline main --grep='row 150' --merges` plus `gh pr view 360 361 --json mergedAt,reviews` documenting which gate each merge passed. Downgrade section 2's heading to "Bounded system-Chrome capture (real-content evidence on Chrome 155; instrument check only on 156)" and add a precondition to M7: before the attended release run, the controller re-executes the V24/V25 probe (script-stripped repeat-GET plus two Chrome passes on a test-helper store) against the rig's current Chrome and records the hashes as V33; if that probe fails, M7 does not start and the M7 fallback escalation applies immediately. Finally, relabel the document status line to "REVISED r6-candidate (iteration 304 designer, post-r5 dispositions, BLOCKED pending r6 quorum)" and move the r5 disposition table under Review status with each item marked "pending re-review" rather than "fixed."

### 1.3 Gemini 4 Argon (`gemini-4-argon`), r6 PASS, verbatim

> **strongest_objection:** Cross-execution state drift: the document is anchored to a dual baseline (original r4 base `c21f303` vs iteration 304 HEAD `039bfa77` where M1–M6 and M7a/M7b have already landed), requiring inline annotations to reconcile 'proposed future work' with already-committed code (V29, V30). While V29–V31 thoroughly verify the landed state, the header's `Base: c21f30310e04941afd3ccfb6b8748177db7b06c8` remains technically behind the V29/V30 measurements.
>
> **catch:** Verify during M7c execution that the baseline content rectangles measured in the V25 test-helper rehearsal are updated to match the exact coordinates of the 4-row evidence roster and pending status rendered by the landed M5 (`bcd2d7cd`) and M2/M4 (`demo_seed.go`) code.
>
> **proposed_fix:** In the header, append `(r4 base; iteration 304 measurements V29–V30 taken at HEAD 039bfa77)` to the `Base:` line so the dual-commit provenance of the verification log is explicit at the top of the document.

### 1.4 Where each objection is answered

| Objection | Answer in this document |
|---|---|
| GLM (a): re-anchor at HEAD; landed/pending per milestone, with SHA and executed evidence | Header `Base:`. §3 inventory, where every landed milestone has its PR, merge SHA, judge verdict and pinning tests. Those tests were run at this commit; nine cmd tests first ran under the superseded v0.52.1 pin (V35), re-verified by V56 under the shipping v0.54.15 pin (9/9 PASS). |
| GLM (b): re-run V24 against the real flag; re-read V4/V6 at HEAD | V37 (V4), V38 (V6/V26), V39 (V14/V24 static) and V40 (real `WorkbenchStill` repeat-GETs on the **seeded demo store**, every capture route, default-mode control). |
| GLM (c): actual run inventory | §6. |
| GLM (d), Kimi (2): real-content two-pass on Chrome 156 | **Done for the `shot` path:** V41. All six shots are byte-identical across two isolated passes on Chrome 156.0.8078.47, on the real seeded store. Each shot was taken by `python3 scripts/marketing_capture.py shot …`, which is the committed M7a supervisor (PR #360) in its `shot` mode. `scripts/marketingcheck/supervisor_test.go` is that supervisor's test, not the supervisor (V55).<br>**Not in the probe path:**<br>• the supervisor's `daemon` lifecycle mode (the probe served the store in-process);<br>• the orchestration shell `scripts/capture_marketing.sh`, which refuses real capture until M7c2b (V55).<br>**Transfer argument:** M7c2a's `pass` mode calls the same `shot()` per row (plan §M7c2a), so Chrome receives the same argv. **Determinism on real Chrome through the full daemon + shell path is unmeasured** and is gated at R1, before any release capture (§7 C14). P0 and P2 re-measure it on capture day. |
| Kimi (1): landing authority | §2 landing provenance, V33 and V34. Stated plainly: no full-seat PASS existed when any milestone landed. |
| Kimi: precondition before the attended run | §4.4 precondition P0, plus the position assertion (r1). GLM r1's release-store containment is P3. |
| Kimi (3): stale self-labelling | This document is r0. §7 lists the parent's rounds and claims no pass. |
| Gemini catch: rectangles must match the landed M5 render | V42/V43. Measuring that render showed the roster is **not in the viewport** at the planned 1440×1400. §4.1 changes the shot geometry (new M7c0) using exact bottoms measured on two seeds (V51). R1 rectangles cover all six shots, with position assertions (§5). |
| Gemini fix: make dual-commit provenance explicit | One base. Every new row is at `4e6e89dc`. Parent rows V1–V31 are cited only where V37–V40 re-measure or supersede them. |

**New finding not raised by any reviewer (V42).** Three of the six planned shots do not contain their required content at HEAD: packet, graph and provenance. The object and entry pages render the full timeline first, and `--screenshot` captures only the viewport. The parent's PD1 resolution (iteration 266) fixed this for the light/dark shots only. §4.1 is the fix.

**r1 finding (V52).** The provenance shot also clips horizontally. The invocation record's JSON is one unbroken token run: the document is 3,715 px wide against a 1,425 px layout. The image therefore shows the record's first fields but not all of its stored input/plan/effect/output references. §4.1 narrows that shot's claim to match.

## 2. Landing provenance (answers Kimi)

This section records how each landed milestone reached `dev`. **This document does not authorize any of them, and nothing here retroactively turns a BLOCKED round into a pass.**

The authorities that were actually in force, measured from the record:

1. **D-WORLD-69 = A** (Mark, attended 2026-10-09; `world-mission.md:988`). It retains row 150's full showcase "each through its normal design/quorum/plan/judge gates".
2. **The narrow-refinement carve-out, iteration 266.** Quorum r3 (`w-marketing-capture-2026-10-14T21-13-09Z.json`) BLOCKED 3/3 with glm, kimi and sonnet present. Every objection carried a concrete fix and none disputed direction. The controller wrote r4 under the shared skill's carve-out, which "routes straight to sprint-planner without a third quorum round" (`world-mission.md:1144-1147`). That routing is what the parent's "after design quorum/freeze" meant at iteration 266. **r4 itself was never quorumed.** (V34)
3. **Iterations 267–276** (M1b–M5a) recorded "Quorum: none owed (plan r4 already routed in 266)" or equivalent (V34).
4. **Iterations 298–301** (M5b, M5c, M6, M7a) ran no pick-time quorum. Iterations 300 and 301 say "Quorum-at-pick: not re-run — the design carries its r1–r4 record from iteration 266" (V34).
5. **Iteration 303** (M7b) ran a fresh pick-time quorum, `w-marketing-capture-2026-10-23T03-32-39Z`. It read PROCEED with Gemini the only voter: GLM and Kimi unreachable (HTTP 429), Sonnet absent on quota. The iteration-304 evaluator called that "a pass with a named hole" (V34).
6. **Every milestone** had its own independent judge verdict, recorded in charter row 150 and in the table below, plus PR CI.

**Plain statement.** No round with all seats present has ever passed this design (r1, r2, r3, r5 and r6 BLOCKED; r4 not quorumed; the 2026-10-23 pick-time PROCEED had one voter). M1a–M7b landed on the authorities above. That is recorded here as fact, not excused. This document authorizes only the pending milestones in §5, and only once its own quorum disposition permits it: for r2 that is the Gate-2 narrow-refinement carve-out (§9), not a quorum PASS, and that is recorded here as fact. The attended capture additionally needs P0 (§4.4) and Mark's attended mint.

| Milestone | PR | Merge / squash commit (V33) | Merged (UTC) | Iteration | Gate at pick | Judge (charter row 150) |
|---|---|---|---|---|---|---|
| M1a | #264 | `28ff4db3` | 2026-10-14 22:15 | 266 | r4 carve-out (r3 BLOCKED) | sonnet PASS 88 |
| M1b | #266 | `52ae7ff9` | 2026-10-15 01:41 | 267 | none owed (r4 routed) | evaluator PASS 90 |
| M2a | #267 | `420c07e7` | 2026-10-15 05:59 | 268 | none owed | independent r3 PASS 90 |
| M2b | #273 | `07d1d9fb` | 2026-10-15 23:30 | 271 | none owed | sonnet FAIL 66 → PASS 86 |
| M3a | #275 | `1ad27432` | 2026-10-16 03:18 | 272 | none owed | sonnet PASS 82 → 90 |
| M3b | #277 | `c08729b6` | 2026-10-16 08:10 | 273 | none owed | sonnet PASS 88 → 93 |
| M4a | #283 | `7ed11523` | 2026-10-16 13:28 | 274 | none owed | sonnet PASS 91 |
| M4b | #289 | `acb3f54d` | 2026-10-16 17:17 | 275 | none owed | sonnet FAIL 66 → 86 → 86 |
| M5a | #295 | `bcd2d7cd` | 2026-10-16 20:24 | 276 | none owed | sonnet PASS 90 |
| M5b | #353 | `f4b2dfdb` (merge) | 2026-10-22 01:21 | 296–298 | none recorded | sonnet PASS 86 → 90 |
| M5c | #354 | `ebdcd7dc` (merge) | 2026-10-22 04:45 | 299 | none recorded | sonnet PASS 80 → 85 |
| M6 | #359 | `47b88a99` (merge) | 2026-10-22 08:18 | 300 | "not re-run" | sonnet PASS 93 |
| M7a | #360 | `c73c7904` (merge) | 2026-10-22 13:10 | 301 | "not re-run" | sonnet PASS 86 → 91 → 91 |
| M7b | #361 | `72dda61d` (merge) | 2026-10-23 04:23 | 302/303 | PROCEED, Gemini only voter | MiniMax PASS 93 |

The judge column is copied from charter row 150 (`world-mission.md:6179`). This document did not re-derive those scores. The charter labels M1a–M5a "LANDED" and M5b–M7b "BUILT". Both mean "merged to dev" (V33); the label difference is the charter's bookkeeping. Merge-commit CI was not re-read by this document.

## 3. As-built inventory at the base (answers GLM)

Executed at `4e6e89dc` on 2026-10-26, load average ≈2. Pins: `AILANG_BIN` = v0.41.0/24ee108. `WORLD_TOOL_AILANG_BIN` = **v0.54.15**/0384450 (V36). Raw output is in `design_docs/verification/world-iter305/v32-*.txt`.

| Milestone | Status | Tests that pin it (file) | Result at base [go test, V35] |
|---|---|---|---|
| M1a | LANDED | `TestWorldDemoEvidencePinnedManifest`, `TestWorldDemoLegacyGradeDrift` (host/coordinator/world_demo_test.go) | PASS, PASS |
| M1b | LANDED | `TestWorldDemoApprovalExactOnePlan` (same file) | PASS |
| M2a | LANDED | `TestDemoSeedAuthenticatedCommits`, `TestDemoSeedPreflightRefusals`, `TestDemoSeedHelpPreparation` (cmd/ailang-worldd/demo_test.go) | PASS ×3 (v0.52.1 run, V35; re-verified v0.54.15, V56) |
| M2b | LANDED | `TestDemoSeedReadTestCheck`, `TestDemoSeedRejectsToolReportRefusals` | PASS ×2 with tool pin v0.54.15. **FAIL ×2 with the plan's v0.52.1** (V36) |
| M3a | LANDED | `TestDemoProofHTTPProducerRoundTrip`, `TestDemoProofSourceHashRefusal` (demo_proof_test.go) | PASS ×2 (v0.52.1 run, V35; re-verified v0.54.15, V56) |
| M3b | LANDED | `TestDemoProofWrongKeyCannotProve`, `TestDemoProofNoVerifiedSetRefuses`, `TestDemoProofUncertainCommitStops` | PASS ×3 (v0.52.1 run, V35; re-verified v0.54.15, V56) |
| M4a | LANDED | `TestDemoSeedApprovalPacket` | PASS (v0.54.15) |
| M4b | LANDED | `TestDemoSeedWhyNamedLinks`, `TestDemoSeedManifestReferences` | PASS ×2 (v0.54.15) |
| M5a | LANDED | `TestWorkbenchPacketRosterAnchoredBounds` (host/daemon/workbench_packet_evidence_test.go) | PASS |
| M5b | LANDED | `TestWorkbenchPacketRosterTestedAndProven`; `TestRenderPacketEvidenceFourRows` (host/workbench/render_test.go) | PASS, PASS |
| M5c | LANDED | `TestDemoSeedRosterKindMismatchRefuses`; `TestDemoSeedEndToEnd` (the plan's M8b gate test) | PASS ×2 (v0.54.15) |
| M6 | LANDED | `TestServeWorkbenchStillFlag` (cmd/ailang-worldd/main_test.go), `TestWorkbenchStillAllRegions` (host/daemon/workbench_still_test.go), `TestWorkbenchStillDefaultLiveScript` (host/workbench/render_test.go) | PASS ×3 (`TestServeWorkbenchStillFlag`: v0.52.1 run, V35; re-verified v0.54.15, V56) |
| M7a | LANDED | `TestMarketingSupervisorFakeChromeSuccess`, `…NoFile`, `…TimeoutKillsGroup`, `…FakeDaemonCleanup` (scripts/marketingcheck/supervisor_test.go) | PASS ×4. The carried shell-hup flake did not appear at this load |
| M7b | LANDED | `TestMarketingPNGStructure`, `TestMarketingPNGContent`, `TestMarketingTwoPassInventory` (png_test.go) | PASS ×3. Package: 7/7 tests, 115 subtests |
| M7c0 | PENDING (new, §4.1) | — | not yet verified |
| M7c1 | PENDING | `TestMarketingProtectedObservationComparator` | absent at base (V46) |
| M7c2a | PENDING | `TestMarketingCapturePassFakeChrome` | absent |
| M7c2b | PENDING | `TestMarketingCaptureFakeChromeSixShots` | absent |
| M7c3 | PENDING | `TestMarketingCaptureInventoryMutation` | absent |
| M7c4 | PENDING | `TestMarketingChromeSixShots` | absent |
| M8a | PENDING | `TestMarketingDocsLinks`, `TestMarketingRunbookMatchesHelp` | absent |
| M8b | PENDING | final gates | not yet verified |

**Not re-verified at base:**
- The landed milestones' mutation ledgers (they are recorded in each milestone's judge report).
- Merge CI.
- The full repository suite (CLAUDE.md gate; not this document's to claim).
- `verify_ail.sh`.

## 4. Design changes for the remaining work

Everything in the parent's D1–D3 and §1/§3 stands as built. The four changes below are this document's own.

### 4.1 Shot geometry: put the required content in the image (new milestone M7c0)

**Measured (V42).** On the seeded demo store at HEAD, served with the real `WorkbenchStill`, the planned viewports captured only part of the required content:
- `packet.png` (1440×1400) and `provenance.png` (1440×1400) captured only the header and timeline entries 0–4.
- `graph.png` (1440×1100) captured the timeline and inspector; the World graph is below the fold.

The cause is the page structure. Object routes render the full timeline before the inspector (`render.go:219` timeline, `:229` inspector, `:248` world graph, `:260` provenance walk). `supportedWorkbenchQuery` refuses `from` on object routes (`workbench.go:83-88`), so the timeline cannot be paged away. The parent §2 table's required content ("validated packet … four-grade evidence roster", "selected entry and checked graph targets", "invocation payload…") is therefore **not in those images** (V53 executes those line citations).

The other three shots were inspected:
- **Light** meets its parent §2 row (parent line 227): header (world, revision 9, log head, stateRoot) and timeline entries 0–3 (V50, controller).
- **Dark** shows the same regions in dark tokens (V42).
- **Phone** meets its row (parent line 232): world/revision/log head and a readable timeline, entries 0–1 (V50, controller). But the phone page overflows horizontally: there is a horizontal scrollbar and the title is clipped to "AILANG World Workbenc". That is a disclosed cosmetic risk (§8), not a content failure.

**Measured candidates (V43):**

| Candidate | Result | Disposition |
|---|---|---|
| URL fragment `#payload-label` (the page's only `id`, V48) | Uniform image, 1 sampled colour, 9,069 B, and **byte-identical across passes** | Rejected. Blank output that still passes the supervisor's structural verdict and two-pass equality. The M7b content gate was never run on these images; V43 observed only the supervisor verdict, the sampled colour count and the hashes. |
| Text fragment `#:~:text=…` (packet, graph) | Uniform image, same as above | Rejected |
| Taller viewport, packet at 1440×4400 | 509,026 B. The validated packet, `approvals chain: pending` and all four roster rows (TESTED · ATTESTED · CLAIMED/demo-authored · PROVEN) are in the image. Byte-identical across passes | **Adopted**; exact height in V51 (3700) |
| Taller viewport, graph and provenance at 1440×4400 | 252,869 B and 418,054 B, byte-identical. The graph section and the provenance payload are in the image | **Adopted**; exact heights in V51 (graph 2100, provenance 3000) |
| Still-only layout reorder (inspector first) | Not built | Rejected. It breaks M6's invariant that still HTML equals default HTML minus the script tag (`workbench_still_test.go:82`), and it is a product layout change outside this row |
| Allowing `from` on object routes | Not built | Rejected for this row. It is a workbench query-grammar change; file it separately if wanted |

**Rule.** A shot's height is H = ⌈(B + 40) / 100⌉ × 100, where B is the bottom edge of its last required region on a frozen seeded store. B is in CSS px at device scale 1, taken from Chrome's own layout (`getBoundingClientRect().bottom` at scroll 0), and is the **maximum across at least two independently seeded stores**.
- **If B + 40 > 4,400 the shot is BLOCKED**: escalate to Mark, M7 does not start, and the height is never clamped.
- Every PNG stays under the parent §3 1 MiB cap.
- Width stays 1440. Light, dark and phone stay at 1440×1100 and 390×844 (V50/V51 show their required regions fit).

**Measured exactly (V51).** The values are identical on two independent seeds, so the delta across seeds is 0.00 px for every region, while the seeds' world heads differ:

| Shot | Required region (selector) | B (px) | H (rule) | Margin at H | PNG at H [stat] |
|---|---|---|---|---|---|
| packet.png | `section[aria-label="decision packet"]` (validated sentence → proof-envelope roster row) | 3564.56 | **3700** | 135.44 | 443,543 B |
| graph.png | `section[aria-label="world graph"]` | 2042.14 | **2100** | 57.86 | 213,531 B |
| provenance.png | `pre.payload` (the stored invocation record) | 2881.17 | **3000** | 118.83 | 284,616 B |

**The method was checked against pixels (V52):**
- The graph card's border row in the 1440×2100 image is at y = 2041. The page background starts at 2042, matching B = 2042.14.
- The packet box's bottom border sits about 2 px above B = 3564.56 in the 3700-tall image. Column 90 crosses its rounded corner, so the pixel scan reads 3559 there. That difference is inside the 40 px margin.
- Re-measuring each region at its new height H gave the same B (delta 0.00). The layout width is the same 1425 px because the page still scrolls vertically.
- The six images at these heights total about 1.18 MB [sum of stat sizes], under the 8 MiB total cap.

M7c0 transcribes these three B values into `png_test.go` constants (§5). Before M7c0 is dispatched, R0 re-confirms them on the then-current Chrome with two seeds.

**Provenance claim narrowed (V52).** At 1440 px the payload `pre` shows the start of the stored invocation record: `invocationId`, `episodeId`, `skillId`, `transitionFn`, `interpreter`, and the wrapped effect ref. Its later fields run past the right edge. The parent row's "invocation payload, stored input/plan/effect/output references" therefore becomes **"the stored invocation record, opened from the approval entry, with its leading fields legible; the complete reference set is the accompanying `why` JSON"**. The parent row already says the `why` JSON "is the actual walk".
- Widening the shot to 3,715 px was not adopted, because it is unreadable as an image.
- Breaking the payload's long tokens with CSS is a workbench renderer change (row 149's surface; `.payload{overflow:auto;white-space:pre-wrap}` at `render.go:198`). It is outside row 150.
- Recorded residual for Mark: if he wants the full record legible in the image, that is a new row. No OPEN row owns it today: row 171 is PARKED hygiene, so it is not a valid owner under D-WORLD-23 (i).

### 4.2 Tool pin correction

The workspace tool pin moved from v0.52.1 to **v0.54.15** at `4b686d8e` (row 139, 2026-10-17). That is after M1a–M5a and before M5b (V36). `daemon.ToolBinaryRelease` reads `pins.VersionLine(pins.RoleTool)` (`workspace.go:52`). Seed tests fail fast against v0.52.1 with "want release AILANG v0.54.15" (V35).

The following are stale and are superseded by this document:
- The parent's "v0.52.1" in V8, S3 and line 181.
- The plan's `export WORLD_TOOL_AILANG_BIN=…/v0.52.1/ailang` (plan line 15).

Every pending milestone and the attended capture use `$HOME/.pinned-ailang-tools/v0.54.15/ailang`. CI already installs it (`.github/workflows/ci.yml:383-401`).

### 4.3 The real drill needs its own rig helper

The plan's M7c4 says "Seed with `demoM5cRealSeed`, then stop the rig". That cannot be done at HEAD:
- `demoM5cRealSeed` returns `(worldReader, manifest, rig.token)` (`demo_test.go:1830`). It exposes neither the db path nor the daemon.
- Its daemon closes only at test cleanup.
- A second daemon on the same db fails on the cross-process writer lock (`daemon.go:744-751`).

The iteration-305 probe solved this by composing `newDemoApprovalRig`'s steps in its own file with the daemon handle retained. It then ran `demo seed`, closed the seed daemon, and reopened the db under `WorkbenchStill` (V40). M7c4 does the same in its new test file. It must not edit `demo_test.go`. The probe source (`design_docs/verification/world-iter305/iter305-seeded-still-chrome-probe_test.go.txt`) is the reference.

### 4.4 Preconditions on the attended capture (answers Kimi)

- **P0.** Within 24 h before the attended release capture, the controller runs the seeded-store probe on the rig's then-current Chrome. It does this on a disposable test-minted store, never the release store. The probe is `iter305-seeded-still-chrome-probe`, copied to `cmd/ailang-worldd/zz_*_test.go`, run with `-count=1`, then deleted. Once M7c4 lands, the M7c4 drill with measured rectangles replaces the probe. The controller records the Chrome `--version`, all six pass-A/pass-B SHA256s and the still-HTML repeat-GET result as a Verification row. **If any shot differs between passes, any route's still HTML varies, or any shot fails the content gate, M7 does not start.** The parent's M7 fallback then applies at once: the varying element is named with both hashes, escalated to Mark, 150(b) stays open, and the row is reported blocked.
- **P0 position assertion.** P0 and R1 do not rely on byte equality and nonblank content alone; V43's fragment images passed the supervisor's structural verdict only, and the M7b content gate was never run on them; they show that structural success plus byte equality admits a blank frame. For every required region of all six shots, the region's Chrome-layout rectangle on the store being captured is measured with the V51 method (GET-only still pages behind the measurement proxy). Two checks apply:
  - (a) the rectangle equals the R1 baseline within ±2 px on every edge;
  - (b) the rectangle lies wholly inside the image, with bottom ≤ H − 40.

  The M7b content-variance check then runs inside that rectangle. A busy but wrong crop fails (a); a region pushed below the fold fails (b).

  **Phone** (390 px; the dump-dom instrument cannot lay out at that width, V51) uses the V58 screenshot scan instead of (a)/(b): (i) phone header and first-timeline-entry rectangles extracted from the capture-day PNG by the V58 scan, matching the R1 phone baseline within ±4 px; (ii) the M7b content-variance gate inside those rectangles; (iii) the disclosed overflow defect unchanged versus the V50 reference image (title clip and horizontal-scrollbar geometry identical within ±2 px). The P0 layout assertion makes no claim about phone.
- **P1.** The rehearsal (§5, R1) has passed on the same Chrome major as the capture day.
- **P2.** Chrome's `--version` is read before pass A and after pass B, and the two must be equal (parent §2, carried into M7c2b).
- **P3. Release-store containment** (GLM's fix). The release store comes from Mark's attended mint, not from a test mint. Before pass A of A1, the controller runs a **GET-only** still probe on the frozen release store, served by `ailang-worldd serve --workbench-still` and measured through the V51 proxy:
  - every required bottom edge must be ≤ H − 40 for its shot;
  - the P0 position assertion must hold against the R1 baseline.

  The probe commits nothing; the parent's before/after public-inventory check brackets it. Two seeds already gave identical layout (V51), but the release store is a different mint, so this is checked on that store rather than assumed. **Any failure, including any B + 40 > 4,400, is BLOCKED and escalated to Mark. Heights are never clamped and nothing is recaptured silently.**

## 5. Pending milestones

Each milestone is one executor run of ≤30 min. Files, interfaces, acceptance and mutants are those of `w-marketing-capture-sprint-plan.md` §M7c1–§M8b (iteration-304 planner, re-checked at this base in V44), **as amended by the column on the right**. Where the plan and this table disagree, this table governs.

| Run | Purpose (plan section) | Amendments from this document |
|---|---|---|
| **M7c0** (new) | Shot geometry. Files: `scripts/marketingcheck/png.go` (`Shots` heights: packet 3700, graph 2100, provenance 3000; nothing else) and `scripts/marketingcheck/png_test.go` | **Constants:** `png_test.go` gains `measuredPacketBottom = 3564.56`, `measuredGraphBottom = 2042.14` and `measuredProvenanceBottom = 2881.17`. Their comment cites V51 and the R0 re-confirm row; a Go test cannot read a V-row, so the constants are the transcription.<br>**Test `TestMarketingShotGeometry`:** each tall `Shots` height **equals** ⌈(constant + 40)/100⌉×100 and is ≤4,400. Light, dark and phone are unchanged.<br>**Canonical list (V54):** `png_test.go:462` hard-codes the expected `Shots` list (packet and provenance 1440×1400, graph 1440×1100) and must be updated to the new heights. The generic 1440×1100 fixtures at `png_test.go:321-337` and `:366` are not tied to `Shots` and stay as they are. `fixture()`/`fixtureSet()` already size from `Shots` (`png_test.go:147,164`).<br>**Mutant M7c0-height:** set packet back to 1400 in `png.go`. It must red on `TestMarketingShotGeometry`'s equality assertion, "packet height must equal the rule value from measuredPacketBottom". The canonical-list subtest also reds, but the named kill is the equality assertion, not the fixture or cap checks.<br>**Dependency:** R0's re-confirm V-row. Its order relative to M7c1 is free; the two share no file. |
| **M7c1** | Protected-file observation comparator (plan §M7c1, unchanged files) | `CompareObservations` iterates over `before` records **looked up by path in a `map[string]ObservationRecord` built from `after`**, never by position. That way mutant M7c1-set reaches `observation-digest` for the missing record instead of panicking (iteration-304 evaluator). |
| **M7c2a** | Python `chrome-version`/`ready`/`pass` modes plus Go `routes` (plan §M7c2a) | `routes` rows carry the **M7c0** dimensions from `Shots`. The `version` subtest keeps the two real strings (155.0.8059.41, 156.0.8078.47). |
| **M7c2b** | Shell end-to-end with fake Chrome and a fake daemon (plan §M7c2b) | Unchanged. The serve argv also names `--tool-ailang-bin` (the v0.54.15 pin is the operator's input; no version is hard-coded in the shell). |
| **M7c3** | Frozen public inventory guard (plan §M7c3) | Unchanged. The seeded store's log is contiguous (indexes 0–9, V40), so `handleLogRange`'s stop-at-gap does not truncate it. |
| **M7c4** | Real Chrome drill, plus inventory wiring into the shell (plan §M7c4) | File list **adds** a local rig helper inside `cmd/ailang-worldd/marketing_chrome_test.go` (§4.3); `demo_test.go` is not edited. **Completion is defined below.** |
| **M8a** | Docs, captions, manifest schema (plan §M8a) | File list **adds `scripts/marketingcheck/png_test.go`**: its `cli/unknown` case passes `docs` and expects rc 2 (`png_test.go:800`), which M8a's `docs` subcommand changes. The environment block uses tool pin v0.54.15. Captions say "static capture of recorded history" and disclose the tall framing. Also applies the R1 threshold V-row to `png.go` and the `TestMarketingPNGContent` boundary literals (V57; owner and gate in §5 R1). |
| **M8b** | Final repository gates and protected-path observation (plan §M8b) | Its `TestMarketingCaptureFakeChromeSixShots` gate runs with **`-timeout 300s`**, not 120 s; the test spans two shell runs of ≤120 s each. |

**M7c4 completion semantics.** M7c4 is COMPLETE when all of the following hold:
(a) its PR lands with PR CI green;
(b) with `WORLD_CHROME` unset the drill prints SKIP, and with `WORLD_CHROME=/nonexistent` it FAILs;
(c) an unset or unloadable `WORLD_MARKETING_RECTS` FAILs with Chrome set (a review check, since that guard is test code);
(d) `TestMarketingCaptureFakeChromeSixShots/inventory` passes, and mutant M7c4-inventory-wiring is killed;
(e) an independent judge PASSes.

Completion does **not** include a real-Chrome pass, and it never marks C8/C9 met. M8a depends on M7c4 COMPLETE. M8a does not depend on R1.

**Controller-owned steps (not executor milestones):**
- **R0, geometry re-confirm.** Exact values already exist (V51). Before M7c0 is dispatched, the controller re-runs `iter305-region-bottoms-probe` on the then-current Chrome against **two** independently seeded test-minted stores. It records both seeds' B values, the per-region delta and the max across seeds as a V-row. If any max differs from the V51 constants, M7c0 transcribes the new max and the V-row says so. If any B + 40 > 4,400, the shot is BLOCKED and escalated to Mark.
- **R1, rehearsal.** After M7c4, measure the content rectangles for **all six shots** over the M7c0 geometry (GLM's fix), each from the V51 Chrome-layout method:
  - light and dark: the `nav[aria-label="world browser"]` header region and the first timeline entry;
  - packet: the validated sentence, `approvals chain: pending`, and each of the four roster rows (Gemini's catch);
  - graph: the `section[aria-label="world graph"]` SVG;
  - provenance: the `pre.payload` region;
  - phone: the header and the first timeline entry, extracted with the V58 scan from the R1 PNG (baseline stored with WORLD_MARKETING_RECTS), because the dump-dom instrument cannot open a 390 px window (V51).

  These rectangles become the P0/P3 position baseline. R1 also runs the full supervisor `daemon` + `capture_marketing.sh` path on real Chrome (C14). Write `WORLD_MARKETING_RECTS`. Run `WORLD_CHROME=… WORLD_MARKETING_RECTS=… go test -run '^TestMarketingChromeSixShots$' ./cmd/ailang-worldd -count=1 -timeout 660s -v`. Record the hashes. Also measure the provisional `png.go` content thresholds (`png.go:27`) on these real renders and record them in a V-row; the values are not the executor's to choose. **Owner and gate of the code change (GLM r2, V57):** `TestMarketingPNGContent` pins the current thresholds as literal boundaries (`png_test.go:359-360` MinRectW/H, `:396` MinColors, `:438-441` LightMinY/DarkMaxY, `:444-455` MinRectStdDev). M8a (executor) applies the R1 V-row's values to `png.go` and moves those boundary literals with them; its gate is `go test ./scripts/marketingcheck -count=1` green with each boundary arm still failing one step past the new value.
- **A1, attended release capture.** The parent's attended block (plan "Release boundary"), run with tool pin v0.54.15, after P0–P2. Then the controller verifies, promotes and commits the six PNGs, and the committed-blob gate runs after that commit.

## 6. Run inventory (replaces "eight executor runs")

| Kind | Items | Status |
|---|---|---|
| Executor runs landed | M1a, M1b, M2a, M2b, M3a, M3b, M4a, M4b, M5a, M5b, M5c, M6, M7a, M7b | 14 landed (§2). M2a and M5b each took more than one fire (M2a "run2"; M5b iterations 296–298) |
| Executor runs pending | M7c0, M7c1, M7c2a, M7c2b, M7c3, M7c4, M8a, M8b | 8 runs. Plan estimates 20+28+28+25+25+25+25 = 176 min plus M7c0 ≈20 min, about 196 min [plan estimates, not measured] |
| Controller steps | R0 (re-confirm), R1 (all-six rectangles, full-path C14), P0, P3 | 4. Each is bounded by its command timeout (probe ≤540 s, measured 11–15 s; drill ≤660 s) |
| Attended | descriptor publish + one `session mint` (Mark types the confirmations), then A1 | 1 session |

Row 150 will then have 22 executor milestones in total (the plan's 21 plus M7c0).

## 7. Acceptance

The parent's C1–C11 and 150(a)–(c) are unchanged, with three additions:
- **C8′.** For the five 1440 px shots, every required region lies inside its image at the M7c0 geometry and the P0 position assertion (±2 px of the R1 baseline, bottom ≤ H − 40) holds. For phone, whose 390 px layout the dump-dom instrument cannot measure (V51), the check is: (i) phone header and first-timeline-entry rectangles extracted from the capture-day PNG by the V58 scan, matching the R1 phone baseline within ±4 px; (ii) the M7b content-variance gate inside those rectangles; (iii) the disclosed overflow defect unchanged versus the V50 reference image (title clip and horizontal-scrollbar geometry identical within ±2 px). C8′ makes no claim that the P0 layout assertion applies to phone. The M7b content gate passes on every region of all six shots.
- **C12.** P0 is recorded within 24 h before A1. It shows 6/6 identical, nonblank shots with the position assertion holding for the five 1440 px shots and the V58 phone check holding for phone, and P3 holds on the release store before pass A (with the same phone carve-out).
- **C13.** No milestone in §5 starts before this document's quorum disposition permits it. For r2 that is the Gate-2 narrow-refinement carve-out recorded in §9 (not a quorum PASS). M7c0 does not start before R0's re-confirm V-row exists.
- **C14.** Before any release capture, R1 records two complete six-shot passes on real Chrome through the **full path**: `capture_marketing.sh` → `marketing_capture.py daemon` serving the store → `pass` → `shot`. All six pass-A/pass-B SHA256s must be equal, and the run must satisfy C8′. V41's `shot`-mode evidence does not discharge C14.

If the M7 fallback triggers, 150(b) remains open and the row is reported blocked, not passed (parent, Kimi r5 text).

## 8. Risks

- **Tall framing aesthetics.** Honest content, unusual image proportions. This is a product choice for Mark, routed as a new row if raised. It is not a blocker.
- **Chrome auto-update.** The rig went from 155 to 156 between iterations 266 and 304. P0 and P2 bound this; equality is promised only within one run.
- **Scrollbars in the image.** Headless shots show a vertical scrollbar, and provenance also shows a horizontal one (wide payload `pre`). Both are deterministic across passes (V41). They are cosmetic and disclosed in captions; no CSS change is made here.
- **Phone overflow** (V50). The phone page overflows horizontally and its title is clipped to "AILANG World Workbenc". This is cosmetic, disclosed in captions, and outside row 150's scope (renderer CSS).
- **Provenance horizontal clip** (V52). The claim is narrowed (§4.1). A legible full record is a residual for Mark.
- **Instrument limit** (V51). With `--window-size=390,844`, headless `--dump-dom` reported `innerWidth` 500, so that instrument cannot measure the 390 px phone layout. Phone regions come from the screenshot instead (V50, R1).
- **Release-store layout.** Two test seeds gave identical layout (V51). The release store is checked directly (P3), never assumed.
- **Graph labels truncated** at 1440 width (SVG text overflows the node box) (V43 image). Cosmetic and pre-existing (row 149). Not changed here.
- **Tool-pin drift.** The pin moved once already (§4.2). The seed tests fail loudly on mismatch, which is correct.
- **Carried shell-hup flake** (M7a). It was green here at load 2 (V35). It is red under load per iteration 304. Executors report it and do not edit `supervisor_test.go`.
- **Fresh seeds differ.** The three probe runs produced three different heads (V41/V43), as D-WORLD-75 = A expects. Equality is claimed only between the two passes over one frozen store.

## 9. Review status

This is **r1** of a new document. It was split from `w-marketing-capture` after that document's r5/r6, with the rules' SPLIT disposition recorded in iteration 304. r0 was committed as `867c3643`.

Parent round history, by artifact:
- r1 `quorum/w-marketing-capture-2026-10-08T13-05-07Z.json` — BLOCKED.
- r2 `…2026-10-08T13-11-05Z.json` — BLOCKED.
- r3 `…2026-10-14T21-13-09Z.json` — BLOCKED 3/3 (glm, kimi, sonnet).
- r4 — the controller's carve-out revision; not quorumed.
- Pick-time `~/.ailang/state/mission-quorum-world/w-marketing-capture-2026-10-23T03-32-39Z.json` — PROCEED with one voter (Gemini).
- r5 `world-iter304/quorum-r5.json` (`…2026-10-26T07-35-08Z`) — BLOCKED.
- r6 `world-iter304/quorum-r6.json` (`…T07-44-46Z`) — BLOCKED.

**None of these is a pass, and none is retroactively made one here.** The parent's r5 dispositions are the parent's. This document re-measures their subjects at HEAD (V37–V40) rather than carrying "fixed" labels.

This document's own rounds:
- **Quorum r1** on r0, artifact `design_docs/verification/world-iter305/quorum-r1.json` (iso_ts `2026-10-26T11:47:05Z`, $0.33, no absent reviewers): Gemini 4 Argon PASS; GLM 5.3 REJECT; Kimi K3 REJECT. Result: **BLOCKED**.
- **r1** is the one permitted revision; one re-quorum follows.

### r1 dispositions (pending re-review)

Each item below is the designer's response. None is marked fixed; the re-quorum decides.

| # | Reviewer: point (their terms) | Response in r1 | Where |
|---|---|---|---|
| 1 | Kimi: §1.4 says "through the committed supervisor" but the supervisor is "nowhere in the command path". Fix: say the supervisor was not in the probe path and either run a supervisor-path pass or gate it at R1 and bind that into §7. | **Partly refuted, partly applied.** The probe ran `python3 ../../scripts/marketing_capture.py shot …` (probe lines 260/270). That script **is** the committed M7a supervisor (PR #360); `supervisor_test.go` is its test (V55, controller). What was not in the path: the supervisor's `daemon` mode and `capture_marketing.sh`. §1.4 and V41 now say exactly that and give the transfer argument (the `pass` mode calls the same `shot()`). Full-path determinism on real Chrome is bound to R1 before any release capture as C14. | §1.4, §5 R1, §7 C14, V41, V55 |
| 2 | Kimi: "R0 must measure bottoms on >=2 independently seeded stores and record the delta … the §4.1 height rule takes the max across seeds". | **Applied, and measured now.** Two independent seeds (different world heads) give identical exact bottoms (delta 0.00). M7c0 heights are now concrete: packet 3700, graph 2100, provenance 3000. The method is checked against pixels. R0 shrinks to a two-seed re-confirm with the max and delta recorded. | §4.1, §5 R0, V51, V52 |
| 3 | Kimi: P0 needs "a position assertion (expected-region signature at expected coordinates)". | **Applied.** P0 and R1 measure each required region's Chrome-layout rectangle on the store being captured. It must equal the R1 baseline within ±2 px, lie inside the image, and have bottom ≤ H − 40; then the content gate runs inside it. | §4.4 P0, §7 C8′/C12 |
| 4 | Kimi: M7c0's test "would still pass when the M7c0-height mutant restores packet to 1400"; state where the measured-bottom constants live and which assertion compares them. | **Applied in the reviewer's terms.** The constants live in `png_test.go`, with a comment citing V51 and the R0 row. `TestMarketingShotGeometry` asserts height equality with the rule value. The mutant reds on that equality assertion, not on fixture or cap checks. | §5 M7c0 |
| 5 | GLM: "'The light/dark/phone shots are unaffected' — is asserted, not measured"; add V50 inspecting light and phone. | **Applied.** V50 (controller): light meets parent line 227; phone meets parent line 232; phone overflows horizontally (cosmetic, disclosed). §4.1 cites V50 instead of asserting "unaffected". Phone meets its row, so its height rule is not extended. | §4.1, §8, V50 |
| 6 | GLM: "Amend C8′ and R1 so rectangle measurement covers all six shots". | **Applied.** | §5 R1, §7 C8′ |
| 7 | GLM: release-store premise. Either verify structural identity or run "a release-store containment measurement before pass A of A1 … every required bottom edge <= shot height - 40"; "bottom edge + 40 exceeds the 4,400 px cap as blocked-and-escalated to Mark, never silently clamped". | **Applied (containment option).** P3 is a GET-only still probe on the frozen release store before pass A. The >4,400 case is BLOCKED and escalated, never clamped (§4.1 rule, P3, R0). | §4.1, §4.4 P3 |
| 8 | GLM (minor): the render.go:219/229/248/260 cites "appear in no V-row". | **Applied.** V53 executes them. | V53 |
| 9 | Gemini catch and fix: R0 records exact unscaled bounds, and fixtures "that hardcode 1400 or 1100 for packet, graph, or provenance must be updated or parameterized to read from Shots". | **Applied.** The exact bounds are V51. V54 (controller C3, re-read here): only the canonical-list assertion at `png_test.go:462` encodes the shot heights, so M7c0 updates it. The generic 1440×1100 fixtures are not shot-bound and stay. `fixture()` already reads `Shots`. | §5 M7c0, V54 |
| 10 | (r1 finding, no reviewer) The provenance shot clips horizontally. | **New.** The claim is narrowed to what is legible; the full record is a residual for Mark. Re-review requested. | §1, §4.1, §8, V52 |

- **Quorum r2** on r1, artifact `design_docs/verification/world-iter305/quorum-r2.json` (`w-marketing-capture-orchestration-2026-10-26T12-00-33Z`, $0.33, no absent reviewers): Gemini 4 Argon PASS; GLM 5.3 REJECT; Kimi K3 REJECT. Result: **BLOCKED**. Surfaces: r1 = light/phone inspection, seed invariance, supervisor-path wording, M7c0 test spec; r2 = the tool pin behind nine §3 results and threshold-test ownership (GLM), phone acceptance (Kimi). Gemini's r2 catch (tight phone rectangles) is met by V58's ±4 px rule.
- **r2 is a controller revision under the Gate-2 narrow-refinement carve-out.** Every r2 blocking objection carried a concrete `proposed_fix`; none disputed the design direction (Gemini passed every round). The fixes below are applied in the reviewers' own words; the controller executed the measurements they named. There is no third quorum round. This is not a quorum PASS and is not described as one.

### r2 dispositions (controller carve-out, iteration 305)

| # | Reviewer: fix (their terms) | Applied | Where |
|---|---|---|---|
| 1 | GLM: "add V56 re-running the nine unmarked tests under the shipping pin … recording rc and 9/9 PASS"; annotate the rows and weaken §1.4. | Ran the reviewer's exact command: rc 0, 9/9 PASS (V56). §1.4 and the four §3 rows now name both pins. | §1.4, §3, V56 |
| 2 | GLM: "add V57 reading TestMarketingPNGContent and TestMarketingTwoPassInventory … for hard-coded provisional png.go:27 thresholds; if pinned, amend §5 R1 to name the test update, its owner and its gate". | Pinned (V57). R1 measures; M8a owns the `png.go` + boundary-literal change; gate named. | §5 R1, M8a, V57 |
| 3 | GLM: "copy /tmp/world-iter305-quorum-r1.json to design_docs/verification/world-iter305/quorum-r1.json and cite the durable path in §9". | Done for r1 and r2. | §9 |
| 4 | GLM: "reword §4.4 to state that V43's fragment images passed the supervisor's structural verdict only and the M7b content gate was never run on them, and delete §4.1's implication that the current gate catches blanks". | Done, both places. | §4.1, §4.4 |
| 5 | Kimi: rewrite C8′ (text given), matching edit to C12, add and execute the phone-rectangle scan on both V41 passes and a deliberately wrong crop, amend R1's phone bullet. | C8′ and the R1 bullet use Kimi's text, with the scan numbered **V58** (Kimi called it V56; GLM's fix already took that number). C12 carries the phone carve-out. V58 executed: both passes PASS, timeline-only crop and blank frame FAIL. | §4.4 P0, §5 R1, §7, V58 |
| 6 | Gemini (PASS) catch: phone rectangles tight enough that a >4 px shift fails. | V58 uses ±4 px on header and first-entry rectangles. | V58 |

## 10. Verification Log (iteration 305, base `4e6e89dc`)

Raw outputs: `design_docs/verification/world-iter305/`. Probe PNGs are kept outside the repo in `~/.ailang/state/world-iter305/probe-{seeded,viewport,tall}/`; images are not committed (the tracked-binary gate allows PNGs only under `website/static/img/v1/`). Probe test files were placed in `cmd/ailang-worldd` only while running and then deleted. Their sources are kept as `.txt`.

| ID | Command or file:line | Observed |
|---|---|---|
| V32 | `git rev-parse HEAD`; `git merge-base --is-ancestor c21f30310e04941afd3ccfb6b8748177db7b06c8 HEAD`; same for `039bfa77` | `4e6e89dc44c9ee9771b415b3cfdc23fadfaadff1`; rc 0; rc 0. Parent r4 base and iteration-304 HEAD are both ancestors of this base. |
| V33 | `for n in 264 266 267 273 275 277 283 289 295 353 354 359 360 361 364; do gh pr view $n --json number,mergedAt,mergeCommit,title,headRefName; done`; `git merge-base --is-ancestor <sha> HEAD` for each merge SHA in §2 plus `acd65340`; `git log --oneline --first-parent origin/dev --grep 'row150\|row 150\|row-150'` | The 15 PRs merged at the times and commits in §2 (#364 = `acd65340`, iteration-304 docs). All 15 SHAs have rc 0. M1a–M5a are squash commits; M5b–M7b are merge commits. Kimi's suggested `--merges` grep alone would miss M1a–M5a, which is why the first-parent log was used too. |
| V34 | `sed -n 22722,23081p design_docs/world-mission-log-archive.md \| grep -ci quorum` = 18 (control: iteration-266 block 22674–22721 = 7), contexts printed with a Python regex scan; `world-mission-log.md` 1196–1376 `grep -ci quorum` = 0; 1377–1446 and 1447–1524 `grep -ci 'quorum-at-pick: not re-run'` = 1 each; 1525–1574 `grep -ci 'PROCEED degraded'` = 1; `jq '.synthesis' ~/.ailang/state/mission-quorum-world/w-marketing-capture-2026-10-23T03-32-39Z.json`; `world-mission.md:1141-1175` | Iteration 266: r3 BLOCKED → "narrow-refinement carve-out … r4 is the controller's bounded revision"; "an r4 re-quorum (the carve-out routes straight to the planner)" is in the ruled-out list. Iterations 271–276 each say "Quorum: none owed (plan r4 [already] routed in 266)". Iterations 267–268 say "R1–r3 remain BLOCKED; inherited iteration266 ratified narrow-refinement carve-out permits r4. No quorum PASS invented." (The other matches in that range belong to rows 114 and 182.) Iterations 298–299: no quorum mention. 300/301: "not re-run". 303: PROCEED degraded, Gemini only. The 2026-10-23 artifact: `verdict: proceed`, absent glm (unreachable), kimi (unreachable), sonnet (quota); cost $0.02414. |
| V35 | `AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_TOOL_AILANG_BIN=… go test -count=1 -v -run '^(…)$'` per package (exact regexes at the top of each `v32-*.txt`) | [go test] host/coordinator 3/3 PASS, rc 0, 8 s; host/daemon 3/3 PASS, rc 0, 77 s; host/workbench 2/2 PASS, rc 0, 1 s. cmd/ailang-worldd with tool `v0.52.1`: rc 1, 9 PASS / 7 FAIL, every failure `want release AILANG v0.54.15` (`v32-cmd.txt`). Re-run with `v0.54.15`: rc 0, 7/7 PASS, 45 s (`v32-cmd-toolpin.txt`). `go test ./scripts/marketingcheck -count=1 -timeout 170s -v`: rc 0, 18 s, 7/7 tests, 115 subtests PASS, 0 FAIL. `TestProductionGoSurface`: PASS. SKIP count 0 in every file. |
| V36 | `git grep -n 'ToolBinaryRelease *=' -- host`; `jq -c '.releases[] \| {release, roles:[.roles[]?.role]}' host/pins/ailang-releases.json`; `git show -s 4b686d8e`; `--is-ancestor 4b686d8e HEAD` / `… 28ff4db3`; both binaries `--version` (40 s Python bound) | `workspace.go:52 var ToolBinaryRelease = pins.VersionLine(pins.RoleTool)`. Roles: v0.41.0 interpreter/proof/…; **v0.54.15 tool**; v0.52.1 retired. `4b686d8e 2026-10-17 feat(row 139): workspace tool pin v0.52.1 → AILANG v0.54.15`, rc 0 (in HEAD), rc 1 (not in M1a). Binaries: `AILANG v0.41.0 / 24ee108`; `AILANG v0.54.15 / 0384450`. Stale v0.52.1 text: parent lines 57, 170, 181; plan lines 15, 19. |
| V37 (supersedes parent V4) | `grep -n 'func \|Proven\|Unavailable\|ValidateProof\|Resolve' host/daemon/workbench_grade.go`; `sed -n 258,298p host/daemon/workbench_packet_evidence.go`; `grep -n 'ValidateProof\|Resolve(sealed)\|"PROVEN"' …evidence.go` | Object grades are still PROVEN-only: `objectGradeView` (`workbench_grade.go:142`) → `ValidateProof` (:159) → `Resolve` (:164); everything else is `NewGradeUnavailable(…)` (:144, :155, :185). The packet roster is separate (`workbench_packet.go:93` → `packetRosterFromValidated`, evidence.go:123). Labels: `recorded-effect` → ATTESTED (:262-264); **`authored-claim` → CLAIMED, producer = author (:265-267)**; `test-invocation` → `packetTestedRow` (:268-272); `proof-envelope` → `packetProvenRow`, PROVEN only after `ValidateProof` (:498) and `Resolve(…).Proven()` (:504), label set at :508. The seeded packet page shows both subject-grade lines "GRADE UNAVAILABLE — stored proof envelopes exist; none is bound to this subject" and the four roster rows (V40, V43 image). V4's premise survives, and D3 is built as designed. |
| V38 (supersedes V6/V26) | `sed -n 44,60p` and `76,92p host/daemon/workbench.go`; `grep -n 'payload\|malformed' host/daemon/workbench.go`; `grep -n HandleFunc host/daemon/daemon.go` | Accepted keys: world, object, from, entry, payload, refsAfter, commitsAfter (:47-55). Object routes accept only object/payload/refsAfter/commitsAfter (:83-88), **so `from` cannot page an object page** (cause of V42). `payload` other than ""/0/1 → 400 "malformed payload flag" (:303-304); `payload=1` shows the payload (:421). Routes: `GET /workbench` (daemon.go:1045), `/workbench/live.js` (:1046), `/v1/objects/by-semantic-id/{name...}` (:1038), `/v1/log` (:1040), `/v1/head` (:1035). New since r4 and not used by row 150: `/surface*` (:1047-1056). |
| V39 (supersedes V14/V24 static half) | `grep -n '<script\|<style\|nonce' host/workbench/render.go`; `grep -c '<script' …` = 1; `git grep -n -i 'workbenchstill\|workbench-still\|\.Still\b\|Still:' -- 'host/*.go' 'cmd/*.go' ':!*_test.go'` | `render.go:190 <style>`; `:203 {{if .Still}}{{else}}<script src="/workbench/live.js" defer></script>{{end}}`. The only script emission is now conditional; there is no nonce and no inline script. Flag path: `main.go:452` (`--workbench-still`, default false) → `:544 WorkbenchStill` → `daemon.go:237` → `workbench.go:360 Page{Still: d.cfg.WorkbenchStill}`. |
| V40 (re-runs V24 on the real flag, seeded store) | `iter305-seeded-still-chrome-probe_test.go.txt` run as `cmd/ailang-worldd/zz_iter305_probe_test.go` with `WORLD_CHROME=<Chrome> PROBE_OUT=~/.ailang/state/world-iter305/probe-seeded go test -count=1 -timeout 540s -v -run '^TestIter305SeededStillCaptureProbe$' ./cmd/ailang-worldd` (both pins). Output `iter305-seeded-still-chrome-probe-output.txt` | rc 0, wall 15 s. Real `demo seed` S1–S6: manifest 77 refs, 4 why rows, approval entry 9. Default daemon, every capture route: `script_tags=1`. Still daemon (`WorkbenchStill: true`, same db): for `/workbench`, `?object=<packet>`, `?from=9&entry=9` and `?object=<invocation>&payload=1`, every result had `repeat_equal=true` (2 s apart), `script_tags=0`, `style_tags=1`, `equals_default_minus_tag=true`. The packet page contains TESTED, ATTESTED, CLAIMED, PROVEN, `approvals chain: pending` and `demo-authored`. The varying-byte inventory of the still HTML on the seeded store is empty. |
| V41 (Chrome 156 real content; answers GLM (d), Kimi (2)) | Same probe, second half: `python3 scripts/marketing_capture.py shot --chrome … --out <pass-a\|pass-b>/<name> --url <still server><route> --width W --height H --scheme light\|dark` for the six parent-§2 shots, two isolated passes; `--version` before and after; `GET /v1/head` before/after | `Google Chrome 156.0.8078.47` before and after (equal). 12/12 verdict `ok`, 0.46–0.81 s each. `HEAD before==after true`. **6/6 byte-identical across passes** [sha256]: light `5a3a9e8a…72a17`, dark `c2d31905…1341b`, packet `51343daf…c478`, graph `e300ff59…53d0`, provenance `c2d56f7b…ae40`, phone `999125ff…94be8`. Every PNG decoded at its exact dimensions with chunks `IHDR IDAT IEND` and 633–898 sampled colours. The dark image renders dark tokens (inspected); the probe's "corner" pixel is the light scrollbar track, not the page background. **This is real-content evidence on Chrome 156 for the parent's planned shot list.** V42 shows that list is incomplete. |
| V42 (new finding) | Read `pass-a/{workbench-dark,packet,graph,provenance}.png` from V41 | packet.png (1440×1400) and provenance.png show only the header and timeline entries 0–4; graph.png shows timeline + inspector; none shows its required region. workbench-dark shows the header and timeline in dark tokens, as required. |
| V43 (framing candidates) | `iter305-viewport-candidates-probe_test.go.txt` and `iter305-tall-viewport-probe_test.go.txt`, run the same way (outputs `…-output.txt`) | Fragment shots (`#payload-label`, `#:~:text=approvals%20chain`, `#:~:text=World%20graph`): verdict ok, **1 sampled colour, 9,069 / 7,129 B, byte-identical** — blank yet equal. Tall: packet 1440×4400 509,026 B `7a356978…9e4`, identical; graph 1440×4400 252,869 B `c6f025ea…27cb`, identical; provenance 1440×4400 418,054 B `ec2a680a…b790`, identical. Inspected: the packet image holds the validated packet, pending status and four roster rows (display ≈1265–1615 of 2000 → ≈2,780–3,550 px); the graph section ends ≈2,050 px; the provenance payload ends ≈2,900 px. Approximate (2.2× downscaled reading); **superseded by V51's exact values.** The three probe runs had heads `751e3331…`, `439d4b1f…`, `874fc47a…` (fresh seeds differ; D-WORLD-75). |
| V44 (plan references at base) | `sed -n 13p;36p;73p scripts/marketingcheck/main.go`; `sed -n 22p;49,53p png.go`; `grep -n '^func BuildCaptureManifest\|^func CheckPass\|^func CheckScheme\|provisional' png.go`; `sed -n 23p;800p png_test.go`; `sed -n 22,36p;80,85p scripts/capture_marketing.sh`; `sed -n 204,206p host/store/context_roots_test.go`; `git grep -n 'func newHandlerDaemon\|func workbenchLiveREST\|func testCommit\|func authHeader' -- host/daemon`; `handleLogRange` `!ok → break` | All of the iteration-304 plan's line references hold: `main.go:13` check gate, `:36` refusal format, `:73` O_EXCL; `png.go:22` refuse, `:49-53` Shots, `:27` "provisional until M7c", `:272` CheckScheme, `:323` CheckPass, `:462` BuildCaptureManifest; `png_test.go:23` wantCode, `:800` `cli/unknown` with `docs`; the capture help lists the ten flags and the tail refuses (rc 1); the allow-list has `main.go`/`png.go`; seams at handlers_test.go:21/47/78 and workbench_live_test.go:74; handlers.go:656-657 stops at a gap. Sizes [wc -l]: shell 85, Python 241, main.go 83, png.go 512. |
| V45 (M7c4 rig) | `sed -n 1727,1830p cmd/ailang-worldd/demo_test.go`; `sed -n 548,651p`; `sed -n 20,40p;740,756p host/daemon/daemon.go` | `return r, m, rig.token` (:1830). `demoSeRig{url, db, token, ws, proofKey, log}` holds no daemon; its `t.Cleanup` closes it (:646-648). `store.Open` failure → "another process already holds writer authority" (:744-751). The probe's inlined rig (V40) seeded, closed and reopened the db three times without error. |
| V46 (absent tests, with control) | `git grep -c 'func TestMarketingProtectedObservationComparator\|func TestMarketingCapturePassFakeChrome\|func TestMarketingCaptureFakeChromeSixShots\|func TestMarketingCaptureInventoryMutation\|func TestMarketingChromeSixShots\|func TestMarketingDocsLinks\|func TestMarketingRunbookMatchesHelp' -- '*.go' \| wc -l`; control `git grep -c 'func TestMarketingTwoPassInventory' -- '*.go'` | 0 matches for the seven pending test names; control 1 (`png_test.go:1`). |
| V47 (instruments) | `"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" --version </dev/null`; `go version`; `uptime` | `Google Chrome 156.0.8078.47`; `go1.26.6 darwin/arm64`; load averages 1.65–2.67 during runs. |
| V48 (anchors) | `grep -c 'id="' host/workbench/render.go`; control `grep -c 'aria-label="' …` | 1 (`payload-label`, render.go:239); control 20. No section has an id, so no fragment anchor targets the packet or graph. Fragments also measured blank (V43). |
| V49 (worktree hygiene) | `git -C <worktree> status --porcelain` after each probe and after deleting `zz_iter305_*_test.go` | r0: only `?? design_docs/verification/world-iter305/` before this document was written. r1: only the three new `iter305-region-bottoms-*` files under that directory. No probe `.go` file remains in the source tree. |
| V50 (controller, iteration 305; answers GLM) | Controller viewed `~/.ailang/state/world-iter305/probe-seeded/pass-a/workbench-light.png` and `phone.png` (V41 outputs) | **workbench-light (1440×1100):** header (world, revision 9, log head, stateRoot) and timeline entries 0–3 are visible, so it meets parent line 227. **phone (390×844):** world/revision/log head and a readable timeline, entries 0–1, are visible, so it meets parent line 232. Phone also overflows horizontally, with a horizontal scrollbar and the title clipped to "AILANG World Workbenc": a disclosed cosmetic risk, not a content failure. |
| V51 (exact bottoms, two seeds; answers Kimi (2), Gemini) | `iter305-region-bottoms-probe_test.go.txt` run as `cmd/ailang-worldd/zz_iter305_bottoms_test.go`, **twice** (two independent `demo seed` runs), with both pins and `WORLD_CHROME=<Chrome 156.0.8078.47> PROBE_OUT=~/.ailang/state/world-iter305/probe-bottoms-seed{1,2} go test -count=1 -timeout 540s -v -run '^TestIter305RegionBottomsProbe$' ./cmd/ailang-worldd`. Method: a measurement-only loopback proxy in front of the `WorkbenchStill` daemon drops the CSP header and appends one external script before `</body>`. The script writes each selector's `getBoundingClientRect()` into a `<pre>` after every section. Chrome is run as `--headless=new --dump-dom --window-size=W,H --force-device-scale-factor=1`. Outputs `iter305-region-bottoms-seed{1,2}-output.txt` | rc 0 and rc 0, 13 s and 12 s [go test]. Seeds differ (heads `44150b4a…` and `72490ae4…`; manifest SHA256 `9efa5602…` and `6a9971dc…`); approval entry 9 in both. **Identical in both seeds** (delta 0.00): packet B = 3564.56 → H 3700; graph B = 2042.14 → H 2100; provenance B = 2881.17 → H 3000. Re-measured at 1440×H, B is unchanged and the margins are 135.44, 57.86 and 118.83. Layout width is 1425 (client) in every confirm run, because the vertical scroll persists (scrollHeight 5242, 2551, 4666). Light at 1440×1100: needed bottom 613.05 (the first timeline entry) ≤ 1060. Each `supervisor shot` at the new heights returned verdict ok. **Instrument limit:** with `--window-size=390,844`, `--dump-dom` reported `innerWidth` 500, so its phone numbers are not 390-px layout and are not used; the phone verdict rests on V50. Dump-dom's inner height is the window height minus 87; viewport height does not move flow layout here, as the unchanged B values confirm. |
| V52 (method check against pixels; provenance clip) | `/tmp/iter305-crop/main.go` (Go `image/png`, outside the repo): `go run . <seed1>/packet.png out 3565`, `… graph.png out 2042`, `… provenance.png out 2881`. It prints column-90 pixels for rows B−6..B+3 and writes a crop with a red marker row at B, which was viewed | **Graph:** row 2041 = `#dedcd3` (card border), row 2042 = `#faf9f5` = background. That matches B = 2042.14 within 1 px. **Packet:** the marker sits about 2 px below the packet box's bottom border, and the roster's last row (proof-envelope … version AILANG v0.41.0) is above it. Column 90 reads the border at row 3559 because it crosses the rounded corner. **Provenance:** the marker sits just below the payload `pre`'s last line. The visible JSON shows `invocationId`, `episodeId`, `skillId`, `transitionFn` and `interpreter`, and is cut at the right edge; V51 records document `scrollWidth` 3715 against client 1425. Sizes at the new heights [stat]: seed 1 packet 443,543, graph 213,531, provenance 284,616; seed 2 443,011, 211,545, 284,459. All are 1440×H as decoded from IHDR. |
| V53 (answers GLM minor) | `grep -n 'aria-label="timeline"\|aria-label="inspector"\|aria-label="world graph"\|aria-label="provenance walk"' host/workbench/render.go`; `grep -n '{{define "packet"}}' host/workbench/render.go`; `grep -n '\.payload{' host/workbench/render.go` | `219:<section aria-label="timeline">`, `229:<section aria-label="inspector">`, `248:<section aria-label="world graph">`, `260:<section aria-label="provenance walk">`. The packet partial at `:273` opens with `<section aria-label="decision packet">` and renders inside the inspector (`:239`). `:198 .payload{overflow:auto;white-space:pre-wrap}`. That is the document order §4.1 cites. |
| V54 (controller, iteration 305 C3; re-read by the designer) | `sed -n 462p scripts/marketingcheck/png_test.go`; `grep -n '1400' scripts/marketingcheck/png_test.go scripts/marketingcheck/supervisor_test.go`; `sed -n 145,200p scripts/marketingcheck/png_test.go \| grep -n 'Width\|Height\|Rect('` | `:462` `expected := []Shot{… {"packet.png", 1440, 1400, "light"}, {"graph.png", 1440, 1100, "light"}, {"provenance.png", 1440, 1400, "light"} …}`. The only other `1400` is an x-coordinate in a rectangle-bounds case (`:357`). `fixture()` builds `image.Rect(0, 0, s.Width, s.Height)` (`:147`) and rectangles from `s.Width` (`:164`). The generic 1440×1100 cases at `:321-337` and `:366` are not tied to `Shots` (controller). |
| V55 (controller, iteration 305 C2; re-read by the designer) | `grep -n 'marketing_capture.py\|"python3", sup, "shot"' design_docs/verification/world-iter305/iter305-seeded-still-chrome-probe_test.go.txt`; `grep -n 'refused' scripts/capture_marketing.sh`; `grep -n 'def daemon\|shot(args' scripts/marketing_capture.py` | Probe `:260 sup := filepath.Join("..", "..", "scripts", "marketing_capture.py")` and `:270 exec.Command("python3", sup, "shot", …)`, i.e. the committed M7a supervisor (PR #360, V33) in `shot` mode. The shell's `:84` prints `six-shot pass is not implemented yet (M7c); capture refused`. The supervisor's `daemon()` (`:135`) calls the same `shot(args, stop)` (`:165`) as `shot` mode (`:233`), but it was not exercised by the probe (the store was served in-process). |
| V56 (controller, iteration 305; GLM r2 fix 1) | `AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang go test -count=1 -v -timeout 540s -run '^(TestDemoSeedAuthenticatedCommits\|TestDemoSeedPreflightRefusals\|TestDemoSeedHelpPreparation\|TestDemoProofHTTPProducerRoundTrip\|TestDemoProofSourceHashRefusal\|TestDemoProofWrongKeyCannotProve\|TestDemoProofNoVerifiedSetRefuses\|TestDemoProofUncertainCommitStops\|TestServeWorkbenchStillFlag)$' ./cmd/ailang-worldd` (`v56-cmd-shipping-pin.txt`); tool binary `--version` = `AILANG v0.54.15` | [go test] rc 0, 44 s wall, `--- PASS` count 9, no FAIL/SKIP lines. |
| V57 (controller, iteration 305; GLM r2 fix 2) | `sed -n 20,35p scripts/marketingcheck/png.go`; `grep -n '159\|160\|\b97\b\|\b96\b\|range \[\]int{7, 8}\|W: 63\|H: 15' scripts/marketingcheck/png_test.go`; control `grep -n '^func Test' png_test.go` = 3 | `png.go:27` "Content thresholds are provisional until M7c measures real renders": MinColors 32, ChromaticDelta 12, ChromaticFraction 0.01, LightMinY 160, DarkMaxY 96, MinRectW 64, MinRectH 16, MinRectStdDev 4.0. `TestMarketingPNGContent` (`png_test.go:307`) pins them as boundaries: `:359-360` (63×16 / 64×15 rejected), `:396` (32 colours rejected), `:438-441` (160/96 accepted, 159/97 rejected), `:444-455` (delta 7 rejected, 8 accepted). `TestMarketingTwoPassInventory` (`:548`) pins no threshold. |
| V58 (controller, iteration 305; Kimi r2 fix) | Scan source `v58-phonescan-main.go.txt` (Go `image/png`; ink = luma < 110 in x < 370, excluding the bottom 20 px; header = band 0; first entry = union of bands 9–14, i.e. `entry 0` through `select entry 0`; content = luma std-dev ≥ 4.0; scrollbar = gray run ≥ 100 px in the bottom 20 px; title right = rightmost ink in band 0). Wrong crops from `v58-mkwrong-main.go.txt` (timeline-only = rows shifted up 330 px; blank 390×844). Run: baseline pass-a phone.png vs pass-a, pass-b, timeline-only, blank, and desktop workbench-light (cross-layout control). Output `v58-phonescan-output.txt` | Baseline header `{25 45 370 68}`, first entry `{41 428 364 610}`, std-dev 104.2/51.7, scrollbar `{6 833 329 841}`, title right 389 (clipped at the 390 px edge). pass-a and pass-b: (i)(ii)(iii) true → **PASS**. Timeline-only: 13 bands, entry not found → (i)(ii)(iii) false → **FAIL**. Blank: 0 bands → **FAIL**. Desktop light: (i) false (entry x 90 vs 41) → **FAIL**. `FAILS 3`, rc 0. An earlier run of the same scan with the entry index set to 7 picked the stateRoot link underline and is superseded (index fixed from the measured band list). The two passes are byte-identical (V41), so their agreement is expected; the wrong crops are the discriminating arms. |
