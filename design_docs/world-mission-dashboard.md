# World mission dashboard — iteration 385 (2026-11-12)

- **Landed this fire:** row 201 A3a sprint 2, "the routes" (PR #556, merge `d1027477`). CI was 3/3 on the head at attempt 2, after a known `host/broker` flake (row 259), and 3/3 on the merge. Judge sonnet PASS 92, zero blocking.
  - Read routes `GET /surface/requests.json` and `GET /surface/requests/{id}`, behind a Host check. Their JSON carries no challenge, artifact, token or message hash.
  - The surface origin set `O`: loopback at the arriving port, `serve --surface-origin`, and the origins of enrolled passkeys. Origin, `Sec-Fetch-Site` and media-type guards.
  - `POST /surface/approve/{id}` runs the guards, then the shared judge at `now`, then the bind check, then records one `surface_approvals` row through a non-waiting writer. It never mints; `approve <id> --await` does.
  - The design's fallback cut fired at 1,995 lines. The SSHSIG POST arm moves to sprint 3, and an `sshsig` POST is a 422 naming `approve <id> --signature`.
  - Sprint 1's carried test gap is closed: a test now fails if the minter judges at the request's creation time instead of the poll's `now`.
- **Next loop work:** A3a sprint 3, "the page and the app":
  - the card page and `approve.js`; *Needs you*, the decision list and the opt-in notification;
  - live-Chrome XA1–XA5;
  - Appendix C's SSHSIG arm (RT6, RT13);
  - the judge's J36 test gap (no test covers an `ActivePasskeys` error on the card route).
- **Attended, pending Mark:** A1b exit run (QUICKSTART §15); A2 real-passkey exit run (A2 §10, QUICKSTART §16).
- **OPEN decisions:** D-WORLD-105 (default C), D-WORLD-114 (default A, applied).
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Release rows 150/151 not routable.
- **Routing this fire:** planner opus (Agent, 77 min), executor haiku (Agent ×2, 13 min), evaluator sonnet (Agent, 9 min); no designer, no quorum. Codex unusable at preflight (rc 75). Metered $0.00.
- **Instrument notes:** the self-notice script reports `watermark channel unavailable` until a directive is marked on #524. Row 259's broker flake reached instance 4 (PR #556 attempt 1).
- **Bookkeeping thread:** #524 (rotates Monday 2026-11-16 07:00).
