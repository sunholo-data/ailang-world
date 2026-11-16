# World mission dashboard — iteration 396 (2026-11-16)

- **Done this fire:** row 243 **V2 child design r1 DONE** ([`w-world-address-bar-verbs-v2.md`](planned/w-world-address-bar-verbs-v2.md)). It puts `why <target>` in the World app's address bar.
  - The walk moves to `host/why` behind a `Getter`, so the CLI and the daemon call one implementation. `why.Head` replaces the forked `headState`.
  - V1's `GET /surface/verb.json` serves it in-process through a private mux over the seven `/v1` GET handlers. It is read-only.
  - Every failure is a refusal that carries the CLI line, never navigation or "not found". The failures are the 10 s read deadline, a read failure, a malformed answer, and an answer over the 1 MiB cap.
  - Every Enter owns the bar, which fixes V1 judge note 1.
  - The work is split into V2a (the seam, ≈ 970 lines, no visible change) and V2b (the bar, ≈ 930 lines).
  - Quorum round 1 BLOCKED 3/3, round 2 **PROCEED 3/3**. The independent audit (sonnet) gave **PASS 84**, 0 blocking. Its 7 planner inputs are in §13.2.
- **Next loop work:** the row 243 **V2a** sprint planner. After that come V2b, V3 (the verb-approval law, with a ledger ask), V4 (host runner and herdr hand-off) and V5.
- **Attended, for Mark when convenient:** the V1 no-terminal exit run (parent design §10).
- **OPEN decisions:** none.
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Release rows 150 and 151 are not routable.
- **Routing this fire:**
  - Designer: opus via Agent (r0 13.9 min, r1 3.0 min).
  - Quorum: gemini-4-argon, oc-glm-5-3 and oc-kimi-k3, $0.81 metered.
  - Evaluator: sonnet via Agent, read-only, 2.8 min.
  - Codex was over ration and unusable at preflight (probe rc 75).
- **Instrument notes:** self-notices still return rc 2 (`watermark channel unavailable`), because the #574/#524 watermark files are absent. Dev CI at `200ea5d5` is 3/3 green.
- **Bookkeeping thread:** #574.
