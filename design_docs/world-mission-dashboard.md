# World mission dashboard — iteration 372 (2026-11-09)

- **Landed this fire:** row 257 `w-exec-cache-seed-at-bind` — PR #531, merge `72079fbe`, CI 3/3 on the PR head (attempt 2, after a 0-job start-up failure) and on the merge. Judge sonnet PASS 91, zero blocking.
- **What it does:** the exec GOCACHE seed is copied when the episode's exec handler is bound, on its own 60 s clock, into staging unique to the daemon process and attempt, then published by one rename with a `.world-seed-complete` marker. A cache without a valid marker is copied again. Exec calls only check, join, or re-attempt (cap 3). `episode report` shows a live `exec_cache` line.
- **Attended today (PR #527):** D-WORLD-110 A, 111 B, 112 A (six keys rotated), 113 B. Records 367–371 abandoned; iteration 366 was canonical until this record.
- **Row 256:** taken by an attended session (rebasing #515, measuring the D-WORLD-111 = B attribution criterion). The loop does not pick, edit or rerun it.
- **OPEN decisions:** D-WORLD-105 only (default C: row 233 parked, fail closed).
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Release rows 150 (attended A1) and 151 (row 188 N2d/N4a) not routable.
- **Next loop row:** 201 `w-world-approvals-by-key`, then 243 — fakes and test keys only (D-WORLD-109 = A); the real-key exit runs stay attended.
- **Routing this fire:** controller opus; designer opus (Agent); quorum glm-5.3 / kimi-k3 / gpt-6.1-sol, $0.63; planner opus (Agent); executor haiku (Agent, three runs); judge sonnet (Agent). Codex unusable at preflight (probe rc 75).
- **Instrument notes:** the head-SHA check-runs poll is blind to `pull_request` runs — poll the run id. Gate-0 self-notice script rc 2 until a directive is marked on #524.
- **Bookkeeping thread:** #524 (rotates Monday 2026-11-16 07:00).
