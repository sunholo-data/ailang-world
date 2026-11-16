# World mission dashboard — iteration 395 (2026-11-16)

- **Done this fire:** row 243 **V1 sprint LANDED** (PR #582, merge `fb4d8b89`; PR and merge CI 3/3 first try): operator verbs in the World app's address bar.
  - `host/verbs`: one grammar for the CLI and the page (`Split`/`Quote`, the verb table, `Lookup`, the shared `ParseInterleaved` and `episode report` flags).
  - A read-only `GET /surface/verb.json`. The bar runs `episode report <e>` with the CLI's own bytes, opens row 201's card for `approve <id|ref>`, and shows the exact CLI line for every other verb.
  - No silent fallback: only a schema-valid address answer navigates. The page waits at most 15 s, and a superseded ask is dropped.
  - 1,909 changed lines in 8 sub-milestones. Judge (sonnet) **PASS 91/100**, 0 blocking. Its 8 non-blocking notes are handed to V2.
- **Next loop work:** row 243 **V2** child design (`why` in the bar). V3 (verb-approval law), V4 (host runner / herdr hand-off) and V5 follow.
- **Attended, for Mark when convenient:** the V1 no-terminal exit run (design §10).
- **OPEN decisions:** none. On 2026-11-16 Mark ruled D-WORLD-105 = C, 114 = A, 115 = A, 116 = A and 117 = A, and published `world/core` 0.1.2.
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Release rows 150 and 151 are not routable.
- **Routing this fire:** planner opus (Agent, 76 min), executor haiku (Agent, 3 runs, 9 min), evaluator sonnet (Agent, own worktree, 55 min). Codex unusable at preflight (probe rc 75). Metered $0.
- **Instrument notes:**
  - Self-notices still rc 2 (`watermark channel unavailable`; the #574/#524 watermark files are absent).
  - Dev's push CI at `db2b7423` was red on the row-259 flake (`TestSeedRestartIgnoresASurvivingOldCopy`, instance 6). Dev's tip is green.
- **Bookkeeping thread:** #574.
