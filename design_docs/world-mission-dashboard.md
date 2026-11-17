# World mission dashboard — iteration 399 (2026-11-17)

- **Done this fire:** row 243 **V3 design** written and reviewed: [`w-world-address-bar-verbs-v3.md`](planned/w-world-address-bar-verbs-v3.md) r2.
  - `ailang-worldd episode finish <ep> --preset <p> --confirm-on-surface`: after the commit, the CLI stores a verb request (the exact prompt lines, the command, the finish scope, title and full body), prints a card link and waits up to 10 minutes.
  - You approve on `/surface/confirm/<id>` with your passkey or SSH key. The daemon verifies it and records the finish approval as decided by that key. The CLI re-checks your signature over its own bytes, then opens the PR.
  - Without the flag the typed `y` is unchanged. Publishing stays typed-only (R4); enrolment stays typed (R3). The address bar still runs no mutating verb until V4.
  - Commit-then-approve, because a preview can never match the real prompt (the commit sha depends on the clock). `episode start` moved to V4 (mints are already key-approvable).
  - Six sprints, ≈ 4,950 raw lines. Sprints 1–2 (the pure law; a card page tested against fixtures) need no answer. Sprints 3–6 (credentials schema 9, the routes, the POST, the CLI flag) wait for D-WORLD-118.
  - Quorum blocked twice (rounds on the store schema, then on the fixes' own claims); the controller measured every premise and applied the reviewers' fixes. Independent audit (sonnet) **PASS 87**, 0 blocking.
- **OPEN decision: D-WORLD-118** — may `episode finish`'s PR confirmation be a verified enrolled key? **A (recommended):** yes, draft or ready. **B:** draft PRs only. **C:** no; V3 stops. **Default:** nothing confirms a finish from the surface and no store schema changes; the loop builds sprints 1–2 only.
- **Next loop work:** row 243 V3 sprint 1 planner (the pure law, ≈ 935 raw lines).
- **Attended, for Mark when convenient:** the V1/V2 no-terminal check (QUICKSTART §19); D-WORLD-118.
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Release rows 150 and 151 are not routable.
- **Routing this fire:** designer opus via Agent (r0 21 min, r1 9 min); evaluator sonnet via Agent (2 min); quorum gemini-4-argon, oc-glm-5-3, oc-kimi-k3 (sol absent on quota both rounds). Codex over ration. Metered $0.98 (quorum).
- **Instrument notes:** self-notices still return rc 2 (`watermark channel unavailable`; the #574/#524 watermark files are absent).
- **Bookkeeping thread:** #574.
