# World mission dashboard — iteration 374 (2026-11-10)

- **Landed this fire:** row 201 sprint 2 (`w-world-approvals-by-key`, A0 M4–M6, plus M3b): PR #535, merge `ca8a4c17`, CI 3/3 on the PR head and on the merge. Judge sonnet PASS 90, zero blocking.
- **What it does:** a mint can now be requested without a terminal and approved later by id.
  - `ailang-worldd episode start <ep> --request` stores a pending request in credentials store schema 4 (`mint_requests`). It opens no terminal and grants nothing.
  - `ailang-worldd approve <id>` checks the request (Decode + Admit) before any terminal opens. It re-derives the request from the world and refuses on any difference. It bounds the y/N by the request's remaining life, then mints through `provisionEpisode`.
  - Consuming the request and inserting the session row happen in one transaction. The clock is read once inside it.
  - `approve --list` shows pending requests. Existing verbs are unchanged.
- **Operator note:** store v4 means a running older daemon refuses the store until `serve` restarts on the new build (as at v3). QUICKSTART §14 has the two commands. It is unverified until Mark's first attended run.
- **Next loop work:** row 201 A1b: approver enrol/list/revoke, `approve --signature`/`--sign-with`, and the `approvers` table at schema 5. It needs a child design r0. Then row 243. Fakes and test keys only (D-WORLD-109 = A); Mark's real-key exit run stays attended.
- **Row 256:** taken by an attended session; the loop does not pick it.
- **OPEN decisions:** D-WORLD-105 only (default C: row 233 parked, fail closed).
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Release rows 150 (attended A1) and 151 (row 188 N2d/N4a) are not routable.
- **Routing this fire:** controller opus; no designer (the r3 design covers M4–M6); planner opus (Agent); executor haiku (Agent, five runs); judge sonnet (Agent). Metered $0. Codex was unusable at preflight (probe rc 75).
- **Instrument notes:**
  - Export `AILANG_BIN=~/.pinned-ailang/ailang` in any shell that runs `host/daemon` tests. Without it, `TestDaemonStartReportsMissingTemplates` reds by design.
  - Run `tracked_binary_gate.sh` only with `verify_go.sh`'s flags.
  - The Gate-0 self-notice script returns rc 2 until a directive is marked on #524.
- **Bookkeeping thread:** #524 (rotates Monday 2026-11-16 07:00).
