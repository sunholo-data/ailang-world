# World mission dashboard — iteration 384 (2026-11-11)

- **Landed this fire:** row 201 A3a sprint 1 (PR #554, merge `b44b5f7d`, CI 3/3 on the head at the first attempt and 3/3 on the merge). Judge sonnet PASS 91, zero blocking.
  - Credentials schema 7 `surface_approvals`, with guards that write 0 rows on any refusal.
  - One shared judge, `host/authority/judge.go`, which the CLI approve paths now call with unchanged refusal texts.
  - The CLI minter `approve <id> --await`: it re-judges a recorded passkey or SSH signature approval against the enrolled key and the request's own bytes, and mints with no terminal.
  - Nothing writes the table in production until sprint 2's POST route.
- **Decision applied:** D-WORLD-114 had no answer, so its default A was built (the CLI mints; the Tauri half waits for row 188 N3). It stays OPEN for Mark to confirm or veto.
- **Next loop work:** A3a sprint 2, "the routes": the read routes, the surface origin set and the Host/Origin/`Sec-Fetch-Site` guards, then the POST. It also carries the judge's test gap: no await test checks the judge runs at `now`.
- **Attended, pending Mark:** A1b exit run (QUICKSTART §15); A2 real-passkey exit run (A2 §10, QUICKSTART §16).
- **OPEN decisions:** D-WORLD-105 (default C), D-WORLD-114 (default A, applied).
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Release rows 150/151 not routable.
- **Routing this fire:** planner opus (Agent, 69 min), executor haiku (Agent ×2, 14 min), evaluator sonnet (Agent, 31 min); no designer, no quorum. Codex unusable at preflight (rc 75). Metered $0.00.
- **Instrument notes:** the self-notice script reports `watermark channel unavailable` until a directive is marked on #524. `commit_ms.sh` deleted the uncommitted plan files once (recovered byte-identical); the banked helper now refuses instead.
- **Bookkeeping thread:** #524 (rotates Monday 2026-11-16 07:00).
