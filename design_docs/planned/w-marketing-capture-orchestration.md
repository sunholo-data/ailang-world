# W-Marketing-Capture-Orchestration — the rest of row 150, written against what is built

**Status:** NEW r0 (iteration 305 designer). This is the as-built split of `w-marketing-capture` for M7c–M8. It is pending quorum.
**Target:** World v1.0.0, queue row 150 / D-WORLD-64, D-WORLD-69 (A), D-WORLD-75 (A), D-WORLD-81 (A)
**Priority:** P1 release · **Estimated:** 8 executor runs (≈196 executor-minutes from the plan's per-run estimates). After those come one controller rehearsal and one attended release capture. See §6 for the inventory.
**Created:** 2026-10-26 (iteration 305)
**Base:** `4e6e89dc44c9ee9771b415b3cfdc23fadfaadff1` (origin/dev). Every row in §8 was measured at this commit. Parent r4 base `c21f303` and iteration-304 HEAD `039bfa77` are both ancestors (V32).
**Planner-Lane:** mission routing default. It applies **only to the pending milestones in §5**, and only after this document's own quorum PASS. This document does not authorize, and never authorized, any landed code (§2).
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
| GLM (a): re-anchor at HEAD; landed/pending per milestone, with SHA and executed evidence | Header `Base:`. §3 inventory, where every landed milestone has its PR, merge SHA, judge verdict and pinning tests. Those tests were **run at this base** (V35). |
| GLM (b): re-run V24 against the real flag; re-read V4/V6 at HEAD | V37 (V4), V38 (V6/V26), V39 (V14/V24 static) and V40 (real `WorkbenchStill` repeat-GETs on the **seeded demo store**, every capture route, default-mode control). |
| GLM (c): actual run inventory | §6. |
| GLM (d), Kimi (2): real-content two-pass on Chrome 156 | **Done:** V41. All six shots are byte-identical across two isolated passes on Chrome 156.0.8078.47, on the real seeded store, through the committed supervisor. |
| Kimi (1): landing authority | §2 landing provenance, V33 and V34. Stated plainly: no full-seat PASS existed when any milestone landed. |
| Kimi: precondition before the attended run | §4.4 precondition P0. |
| Kimi (3): stale self-labelling | This document is r0. §7 lists the parent's rounds and claims no pass. |
| Gemini catch: rectangles must match the landed M5 render | V42/V43. Measuring that render showed the roster is **not in the viewport** at the planned 1440×1400. §4.1 changes the shot geometry (new M7c0). |
| Gemini fix: make dual-commit provenance explicit | One base. Every new row is at `4e6e89dc`. Parent rows V1–V31 are cited only where V37–V40 re-measure or supersede them. |

**New finding not raised by any reviewer (V42).** Three of the six planned shots do not contain their required content at HEAD: packet, graph and provenance. The object and entry pages render the full timeline first, and `--screenshot` captures only the viewport. The parent's PD1 resolution (iteration 266) fixed this for the light/dark shots only. §4.1 is the fix.

## 2. Landing provenance (answers Kimi)

This section records how each landed milestone reached `dev`. **This document does not authorize any of them, and nothing here retroactively turns a BLOCKED round into a pass.**

The authorities that were actually in force, measured from the record:

1. **D-WORLD-69 = A** (Mark, attended 2026-10-09; `world-mission.md:988`). It retains row 150's full showcase "each through its normal design/quorum/plan/judge gates".
2. **The narrow-refinement carve-out, iteration 266.** Quorum r3 (`w-marketing-capture-2026-10-14T21-13-09Z.json`) BLOCKED 3/3 with glm, kimi and sonnet present. Every objection carried a concrete fix and none disputed direction. The controller wrote r4 under the shared skill's carve-out, which "routes straight to sprint-planner without a third quorum round" (`world-mission.md:1144-1147`). That routing is what the parent's "after design quorum/freeze" meant at iteration 266. **r4 itself was never quorumed.** (V34)
3. **Iterations 267–276** (M1b–M5a) recorded "Quorum: none owed (plan r4 already routed in 266)" or equivalent (V34).
4. **Iterations 298–301** (M5b, M5c, M6, M7a) ran no pick-time quorum. Iterations 300 and 301 say "Quorum-at-pick: not re-run — the design carries its r1–r4 record from iteration 266" (V34).
5. **Iteration 303** (M7b) ran a fresh pick-time quorum, `w-marketing-capture-2026-10-23T03-32-39Z`. It read PROCEED with Gemini the only voter: GLM and Kimi unreachable (HTTP 429), Sonnet absent on quota. The iteration-304 evaluator called that "a pass with a named hole" (V34).
6. **Every milestone** had its own independent judge verdict, recorded in charter row 150 and in the table below, plus PR CI.

**Plain statement.** No round with all seats present has ever passed this design (r1, r2, r3, r5 and r6 BLOCKED; r4 not quorumed; the 2026-10-23 pick-time PROCEED had one voter). M1a–M7b landed on the authorities above. That is recorded here as fact, not excused. This document authorizes only the pending milestones in §5, and only after its own quorum PASS. The attended capture additionally needs P0 (§4.4) and Mark's attended mint.

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
| M2a | LANDED | `TestDemoSeedAuthenticatedCommits`, `TestDemoSeedPreflightRefusals`, `TestDemoSeedHelpPreparation` (cmd/ailang-worldd/demo_test.go) | PASS ×3 |
| M2b | LANDED | `TestDemoSeedReadTestCheck`, `TestDemoSeedRejectsToolReportRefusals` | PASS ×2 with tool pin v0.54.15. **FAIL ×2 with the plan's v0.52.1** (V36) |
| M3a | LANDED | `TestDemoProofHTTPProducerRoundTrip`, `TestDemoProofSourceHashRefusal` (demo_proof_test.go) | PASS ×2 |
| M3b | LANDED | `TestDemoProofWrongKeyCannotProve`, `TestDemoProofNoVerifiedSetRefuses`, `TestDemoProofUncertainCommitStops` | PASS ×3 |
| M4a | LANDED | `TestDemoSeedApprovalPacket` | PASS (v0.54.15) |
| M4b | LANDED | `TestDemoSeedWhyNamedLinks`, `TestDemoSeedManifestReferences` | PASS ×2 (v0.54.15) |
| M5a | LANDED | `TestWorkbenchPacketRosterAnchoredBounds` (host/daemon/workbench_packet_evidence_test.go) | PASS |
| M5b | LANDED | `TestWorkbenchPacketRosterTestedAndProven`; `TestRenderPacketEvidenceFourRows` (host/workbench/render_test.go) | PASS, PASS |
| M5c | LANDED | `TestDemoSeedRosterKindMismatchRefuses`; `TestDemoSeedEndToEnd` (the plan's M8b gate test) | PASS ×2 (v0.54.15) |
| M6 | LANDED | `TestServeWorkbenchStillFlag` (cmd/ailang-worldd/main_test.go), `TestWorkbenchStillAllRegions` (host/daemon/workbench_still_test.go), `TestWorkbenchStillDefaultLiveScript` (host/workbench/render_test.go) | PASS ×3 |
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

The cause is the page structure. Object routes render the full timeline before the inspector (`render.go:219` timeline, `:229` inspector, `:248` world graph, `:260` provenance walk). `supportedWorkbenchQuery` refuses `from` on object routes (`workbench.go:83-88`), so the timeline cannot be paged away. The parent §2 table's required content ("validated packet … four-grade evidence roster", "selected entry and checked graph targets", "invocation payload…") is therefore **not in those images**. The light/dark/phone shots are unaffected.

**Measured candidates (V43):**

| Candidate | Result | Disposition |
|---|---|---|
| URL fragment `#payload-label` (the page's only `id`, V48) | Uniform image, 1 sampled colour, 9,069 B, and **byte-identical across passes** | Rejected. Blank output that still passes two-pass equality. It shows why the M7b nonblank/content gate is load-bearing. |
| Text fragment `#:~:text=…` (packet, graph) | Uniform image, same as above | Rejected |
| Taller viewport, packet at 1440×4400 | 509,026 B. The validated packet, `approvals chain: pending` and all four roster rows (TESTED · ATTESTED · CLAIMED/demo-authored · PROVEN) are in the image. Byte-identical across passes | **Adopted** |
| Taller viewport, graph and provenance at 1440×4400 | 252,869 B and 418,054 B, byte-identical. The graph section and the provenance payload are in the image | **Adopted** (heights to be measured; see below) |
| Still-only layout reorder (inspector first) | Not built | Rejected. It breaks M6's invariant that still HTML equals default HTML minus the script tag (`workbench_still_test.go:82`), and it is a product layout change outside this row |
| Allowing `from` on object routes | Not built | Rejected for this row. It is a workbench query-grammar change; file it separately if wanted |

**Rule.** Each shot's height is the smallest multiple of 100 px that is ≥ (bottom edge of the shot's last required region on the frozen seeded store + 40 px). It is never more than 4,400 px, and every PNG stays under the parent §3 1 MiB cap.
- The controller measures the bottom edges in the rehearsal. The instrument is the committed probe `iter305-tall-viewport-probe` plus a pixel scan for the region boxes. The result is recorded as a Verification row before M7c0 lands.
- Approximate values read from 2.2×-downscaled views (V43; **not acceptance values**): packet ≈3,600, graph ≈2,100, provenance ≈3,000.
- Width stays 1440. Light, dark and phone stay at 1440×1100 and 390×844.
- `Shots` (`scripts/marketingcheck/png.go:49-53`) is the single source of dimensions. The capture-manifest route and dimension rows follow it, so routes and sizes cannot drift.
- Tall images are honest but not the aesthetic Mark may want. Their only alternative is a product layout change, which this document does not make. The tall geometry discharges 150(c) as written. If Mark wants different framing, that is a new row and does not reopen 150(b).

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
- **P1.** The rehearsal (§5, R1) has passed on the same Chrome major as the capture day.
- **P2.** Chrome's `--version` is read before pass A and after pass B, and the two must be equal (parent §2, carried into M7c2b).

## 5. Pending milestones

Each milestone is one executor run of ≤30 min. Files, interfaces, acceptance and mutants are those of `w-marketing-capture-sprint-plan.md` §M7c1–§M8b (iteration-304 planner, re-checked at this base in V44), **as amended by the column on the right**. Where the plan and this table disagree, this table governs.

| Run | Purpose (plan section) | Amendments from this document |
|---|---|---|
| **M7c0** (new) | Shot geometry. Files: `scripts/marketingcheck/png.go` (`Shots` heights for packet, graph and provenance only) and `png_test.go` | Inputs are the controller's measured heights (§4.1 rule; a V-row must exist before the run). Test `TestMarketingShotGeometry`: each tall shot is listed, has height ≤4,400, and a fixture of that size passes structure + cap checks. Mutant **M7c0-height**: restore packet to 1400; the geometry assertion "packet shot must cover the measured roster bottom" must red. Dependency: R0 (the controller's measurement). Its place before M7c1 is free; it shares no file with M7c1. |
| **M7c1** | Protected-file observation comparator (plan §M7c1, unchanged files) | `CompareObservations` iterates over `before` records **looked up by path in a `map[string]ObservationRecord` built from `after`**, never by position. That way mutant M7c1-set reaches `observation-digest` for the missing record instead of panicking (iteration-304 evaluator). |
| **M7c2a** | Python `chrome-version`/`ready`/`pass` modes plus Go `routes` (plan §M7c2a) | `routes` rows carry the **M7c0** dimensions from `Shots`. The `version` subtest keeps the two real strings (155.0.8059.41, 156.0.8078.47). |
| **M7c2b** | Shell end-to-end with fake Chrome and a fake daemon (plan §M7c2b) | Unchanged. The serve argv also names `--tool-ailang-bin` (the v0.54.15 pin is the operator's input; no version is hard-coded in the shell). |
| **M7c3** | Frozen public inventory guard (plan §M7c3) | Unchanged. The seeded store's log is contiguous (indexes 0–9, V40), so `handleLogRange`'s stop-at-gap does not truncate it. |
| **M7c4** | Real Chrome drill, plus inventory wiring into the shell (plan §M7c4) | File list **adds** a local rig helper inside `cmd/ailang-worldd/marketing_chrome_test.go` (§4.3); `demo_test.go` is not edited. **Completion is defined below.** |
| **M8a** | Docs, captions, manifest schema (plan §M8a) | File list **adds `scripts/marketingcheck/png_test.go`**: its `cli/unknown` case passes `docs` and expects rc 2 (`png_test.go:800`), which M8a's `docs` subcommand changes. The environment block uses tool pin v0.54.15. Captions say "static capture of recorded history" and disclose the tall framing. |
| **M8b** | Final repository gates and protected-path observation (plan §M8b) | Its `TestMarketingCaptureFakeChromeSixShots` gate runs with **`-timeout 300s`**, not 120 s; the test spans two shell runs of ≤120 s each. |

**M7c4 completion semantics.** M7c4 is COMPLETE when all of the following hold:
(a) its PR lands with PR CI green;
(b) with `WORLD_CHROME` unset the drill prints SKIP, and with `WORLD_CHROME=/nonexistent` it FAILs;
(c) an unset or unloadable `WORLD_MARKETING_RECTS` FAILs with Chrome set (a review check, since that guard is test code);
(d) `TestMarketingCaptureFakeChromeSixShots/inventory` passes, and mutant M7c4-inventory-wiring is killed;
(e) an independent judge PASSes.

Completion does **not** include a real-Chrome pass, and it never marks C8/C9 met. M8a depends on M7c4 COMPLETE. M8a does not depend on R1.

**Controller-owned steps (not executor milestones):**
- **R0, geometry measurement.** Tall captures of packet, graph and provenance on a test-minted seeded store; bottom edges recorded; heights per the §4.1 rule. Instrument: `iter305-tall-viewport-probe`. Runs before M7c0.
- **R1, rehearsal.** After M7c4: measure the content rectangles over the M7c0 geometry, where the packet rectangles cover the validated header and all four roster rows (Gemini's catch). Write `WORLD_MARKETING_RECTS`. Run `WORLD_CHROME=… WORLD_MARKETING_RECTS=… go test -run '^TestMarketingChromeSixShots$' ./cmd/ailang-worldd -count=1 -timeout 660s -v`. Record the hashes. Also set the provisional `png.go` content thresholds from these real renders (`png.go:27`), with a measured V-row; the content thresholds are not the executor's to choose.
- **A1, attended release capture.** The parent's attended block (plan "Release boundary"), run with tool pin v0.54.15, after P0–P2. Then the controller verifies, promotes and commits the six PNGs, and the committed-blob gate runs after that commit.

## 6. Run inventory (replaces "eight executor runs")

| Kind | Items | Status |
|---|---|---|
| Executor runs landed | M1a, M1b, M2a, M2b, M3a, M3b, M4a, M4b, M5a, M5b, M5c, M6, M7a, M7b | 14 landed (§2). M2a and M5b each took more than one fire (M2a "run2"; M5b iterations 296–298) |
| Executor runs pending | M7c0, M7c1, M7c2a, M7c2b, M7c3, M7c4, M8a, M8b | 8 runs. Plan estimates 20+28+28+25+25+25+25 = 176 min plus M7c0 ≈20 min, about 196 min [plan estimates, not measured] |
| Controller steps | R0, R1, P0 | 3. Each is bounded by its command timeout (probe ≤540 s; drill ≤660 s) |
| Attended | descriptor publish + one `session mint` (Mark types the confirmations), then A1 | 1 session |

Row 150 will then have 22 executor milestones in total (the plan's 21 plus M7c0).

## 7. Acceptance

The parent's C1–C11 and 150(a)–(c) are unchanged, with three additions:
- **C8′.** Each shot's required regions lie inside its image at the M7c0 geometry. The R1 rectangles are measured inside those bounds, and the M7b content gate passes on every region.
- **C12.** P0 is recorded within 24 h before A1 and shows 6/6 identical, nonblank shots.
- **C13.** No milestone in §5 starts before this document's quorum PASS. M7c0 does not start before R0's V-row exists.

If the M7 fallback triggers, 150(b) remains open and the row is reported blocked, not passed (parent, Kimi r5 text).

## 8. Risks

- **Tall framing aesthetics.** Honest content, unusual image proportions. This is a product choice for Mark, routed as a new row if raised. It is not a blocker.
- **Chrome auto-update.** The rig went from 155 to 156 between iterations 266 and 304. P0 and P2 bound this; equality is promised only within one run.
- **Scrollbars in the image.** Headless shots show a vertical scrollbar, and provenance also shows a horizontal one (wide payload `pre`). Both are deterministic across passes (V41). They are cosmetic and disclosed in captions; no CSS change is made here.
- **Graph labels truncated** at 1440 width (SVG text overflows the node box) (V43 image). Cosmetic and pre-existing (row 149). Not changed here.
- **Tool-pin drift.** The pin moved once already (§4.2). The seed tests fail loudly on mismatch, which is correct.
- **Carried shell-hup flake** (M7a). It was green here at load 2 (V35). It is red under load per iteration 304. Executors report it and do not edit `supervisor_test.go`.
- **Fresh seeds differ.** The three probe runs produced three different heads (V41/V43), as D-WORLD-75 = A expects. Equality is claimed only between the two passes over one frozen store.

## 9. Review status

This is **r0** of a new document. It was split from `w-marketing-capture` after that document's r5/r6, with the rules' SPLIT disposition recorded in iteration 304.

Parent round history, by artifact:
- r1 `quorum/w-marketing-capture-2026-10-08T13-05-07Z.json` — BLOCKED.
- r2 `…2026-10-08T13-11-05Z.json` — BLOCKED.
- r3 `…2026-10-14T21-13-09Z.json` — BLOCKED 3/3 (glm, kimi, sonnet).
- r4 — the controller's carve-out revision; not quorumed.
- Pick-time `~/.ailang/state/mission-quorum-world/w-marketing-capture-2026-10-23T03-32-39Z.json` — PROCEED with one voter (Gemini).
- r5 `world-iter304/quorum-r5.json` (`…2026-10-26T07-35-08Z`) — BLOCKED.
- r6 `world-iter304/quorum-r6.json` (`…T07-44-46Z`) — BLOCKED.

**None of these is a pass, and none is retroactively made one here.** The parent's r5 dispositions are the parent's. This document re-measures their subjects at HEAD (V37–V40) rather than carrying "fixed" labels.

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
| V43 (framing candidates) | `iter305-viewport-candidates-probe_test.go.txt` and `iter305-tall-viewport-probe_test.go.txt`, run the same way (outputs `…-output.txt`) | Fragment shots (`#payload-label`, `#:~:text=approvals%20chain`, `#:~:text=World%20graph`): verdict ok, **1 sampled colour, 9,069 / 7,129 B, byte-identical** — blank yet equal. Tall: packet 1440×4400 509,026 B `7a356978…9e4`, identical; graph 1440×4400 252,869 B `c6f025ea…27cb`, identical; provenance 1440×4400 418,054 B `ec2a680a…b790`, identical. Inspected: the packet image holds the validated packet, pending status and four roster rows (display ≈1265–1615 of 2000 → ≈2,780–3,550 px); the graph section ends ≈2,050 px; the provenance payload ends ≈2,900 px. Approximate (2.2× downscaled reading); R0 measures exactly. The three probe runs had heads `751e3331…`, `439d4b1f…`, `874fc47a…` (fresh seeds differ; D-WORLD-75). |
| V44 (plan references at base) | `sed -n 13p;36p;73p scripts/marketingcheck/main.go`; `sed -n 22p;49,53p png.go`; `grep -n '^func BuildCaptureManifest\|^func CheckPass\|^func CheckScheme\|provisional' png.go`; `sed -n 23p;800p png_test.go`; `sed -n 22,36p;80,85p scripts/capture_marketing.sh`; `sed -n 204,206p host/store/context_roots_test.go`; `git grep -n 'func newHandlerDaemon\|func workbenchLiveREST\|func testCommit\|func authHeader' -- host/daemon`; `handleLogRange` `!ok → break` | All of the iteration-304 plan's line references hold: `main.go:13` check gate, `:36` refusal format, `:73` O_EXCL; `png.go:22` refuse, `:49-53` Shots, `:27` "provisional until M7c", `:272` CheckScheme, `:323` CheckPass, `:462` BuildCaptureManifest; `png_test.go:23` wantCode, `:800` `cli/unknown` with `docs`; the capture help lists the ten flags and the tail refuses (rc 1); the allow-list has `main.go`/`png.go`; seams at handlers_test.go:21/47/78 and workbench_live_test.go:74; handlers.go:656-657 stops at a gap. Sizes [wc -l]: shell 85, Python 241, main.go 83, png.go 512. |
| V45 (M7c4 rig) | `sed -n 1727,1830p cmd/ailang-worldd/demo_test.go`; `sed -n 548,651p`; `sed -n 20,40p;740,756p host/daemon/daemon.go` | `return r, m, rig.token` (:1830). `demoSeRig{url, db, token, ws, proofKey, log}` holds no daemon; its `t.Cleanup` closes it (:646-648). `store.Open` failure → "another process already holds writer authority" (:744-751). The probe's inlined rig (V40) seeded, closed and reopened the db three times without error. |
| V46 (absent tests, with control) | `git grep -c 'func TestMarketingProtectedObservationComparator\|func TestMarketingCapturePassFakeChrome\|func TestMarketingCaptureFakeChromeSixShots\|func TestMarketingCaptureInventoryMutation\|func TestMarketingChromeSixShots\|func TestMarketingDocsLinks\|func TestMarketingRunbookMatchesHelp' -- '*.go' \| wc -l`; control `git grep -c 'func TestMarketingTwoPassInventory' -- '*.go'` | 0 matches for the seven pending test names; control 1 (`png_test.go:1`). |
| V47 (instruments) | `"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" --version </dev/null`; `go version`; `uptime` | `Google Chrome 156.0.8078.47`; `go1.26.6 darwin/arm64`; load averages 1.65–2.67 during runs. |
| V48 (anchors) | `grep -c 'id="' host/workbench/render.go`; control `grep -c 'aria-label="' …` | 1 (`payload-label`, render.go:239); control 20. No section has an id, so no fragment anchor targets the packet or graph. Fragments also measured blank (V43). |
| V49 (worktree hygiene) | `git -C <worktree> status --porcelain` after each probe and after deleting `zz_iter305_*_test.go` | Only `?? design_docs/verification/world-iter305/` before this document was written. No probe `.go` file remains in the source tree. |
