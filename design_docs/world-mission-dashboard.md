# World mission dashboard — iteration 394 (2026-11-16)

- **Done this fire:** row 243 design **r2 DONE** — [`w-world-address-bar-verbs.md`](planned/w-world-address-bar-verbs.md), operator verbs in the World app's address bar with one Go grammar shared by the CLI and the page.
  - Parent with phases V1–V5. **V1 is specified to build** (1 sprint ≈ 1,750 lines): `host/verbs`, a read-only `GET /surface/verb.json`, `episode report` run from the bar, `approve <id|ref>` opening row 201's cards, the CLI line shown for start/finish/launch/why. No page-driven state change.
  - V2 (`why`), V3 (verb-approval law), V4 (host runner / herdr hand-off), V5 (exit-run readiness) each need a child design.
  - Quorum blocked twice (r1: silent fallback, shadow set; r2: bare `episode`, unbounded fetch), controller carve-out r2, independent audit (sonnet) **CARVE-OUT VALID 84/100**; its one blocking finding was fixed.
- **Next loop work:** row 243 V1 sprint planner.
- **Ready for Mark:** `world/core` 0.1.2 candidate on dev (D-WORLD-117: publish?).
- **Attended, pending Mark:** exit runs for A1b (QUICKSTART §15), A2 (§16), A3a (§17), A4 (§18 / design §10); the 0.1.2 publish.
- **OPEN decisions:** D-WORLD-105 (default C), D-WORLD-114 (A, applied), D-WORLD-115 (A, applied), D-WORLD-116 (A, built), D-WORLD-117 (pending). No new ask this fire.
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Release rows 150 and 151 are not routable.
- **Routing this fire:** designer opus (Agent; r0 10 min, r1 5 min), quorum ×2 ($0.675 metered, glm/kimi/sol present), evaluator sonnet (Agent, audit). Planner/executor not spawned (design-only fire). Codex unusable at preflight (probe rc 75).
- **Instrument notes:**
  - Self-notices still rc 2 (`watermark channel unavailable`; #574/#524 watermark files absent).
  - Transcript-scan notice for the 393 fire: 2 hits (one key name), in another project's session files.
- **Bookkeeping thread:** #574.
