# World mission dashboard — iteration 241 (2026-10-07)

- Outcome: row 153 LANDED (PR #227 → `42a2511`); closes hygiene row 161.
- What changed: plan-phase timeout is a typed `PhaseTimeoutError` (A2A names the phase, retryable; MCP wire unchanged until row 159 / ailang#1602).
- Plan cap derived: invokeDeadline − HandlerCap − HandlerHeadroom = 4 s (was an unmeasured 2 s).
- Root cause: every capsule run cold-compiled the transition + std (5–7× slower); now runs copy a compile-cache template built in the publication check; cold runs are labelled and red in test rigs.
- Evidence: throttled-proxy plan timeouts 4/4 → 0/12; CI 5/5 independent attempts green on the PR head.
- Judge: opus r1 80 (1 blocking: machine-dependent ratio test) → r2 92, zero blocking.
- Quorum: r1/r2 BLOCKED at N−1 (gpt6-1-sol unreachable) → r3 narrow-refinement carve-out, reviewers' fixes verbatim.
- Folded in: iterations 239/240 unjudged drafts (Codex-controller lane park, PR #226) now recorded; row 164 scoped to Codex-controller fires.
- 1.0: clauses 1/2/3/6/7 MET; 4/5 UNMET.
- Next: row 152 → 93 → 114 → 139.
- New rows: 165 (capsule re-hashes 100 MB interpreter per run), 166 (broker head/tail timeout flake).
- Pending human: none (ledger ZERO OPEN).
- Dev base at Gate 1: `990a6da`, 2/2 green.
- Metered $0.49 (quorum); Anthropic subscription roles.
