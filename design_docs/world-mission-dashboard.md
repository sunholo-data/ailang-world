# World mission dashboard — iteration 314 (2026-10-28)

- **This fire:** dev was RED at Gate 1, so it preempted the queue: the arm 6 reverse-leg false refusal is fixed as row 229 (PR #403).
- **The red:** `origin/dev` `2769cff9` failed `TestExecSrtMCPGrantAndBudget` with "reverse leg: the host reached the listener the child bound…". A re-run on the identical tree was green, so it is a flake. The same false refusal could refuse daemon startup in production.
- **Cause:** the probe closed the port it picked, then counted ANY connect as the sandboxed child. On Linux the child is in its own netns, so a reach was always some other host process; srt's own mux binds an ephemeral port at every start.
- **Fix:** a fresh port and nonce per attempt, plus a hello on every dial. The child's nonce means refuse; other bytes or EOF mean foreign, so retry (≤3), then refuse "inconclusive"; silence means a sticky refusal. A leg passes only if no dial connected.
- **Measured first:** live against the pinned srt 0.0.78 mux, the hello gets an HTTP 400 in ≤3 ms (3/3). That is what made the quorum's last objection carve-out-eligible.
- **Quorum:** r1 and r2 BLOCKED (kimi, glm); r3 applied the reviewers' verbatim fixes under the carve-out. gpt6-1-sol was absent both rounds (OpenAI API 429, no credits). Spend $0.69.
- **Judge:** sonnet PASS 92/100, zero blocking. Two surviving mutants (hello wait, `c.resume()`) became row 230.
- **CI:** PR head 3/3 green, and all 7 new tests plus the real-srt tests PASS on Linux. The merge commit is polled before LANDED.
- **Roles (Agent tool):** designer opus ×3 runs (one doc; the third run is FLAGGED over diet), planner opus, executor haiku, evaluator sonnet.
- **Clauses:** 1–7 MET; goal unmoved. Release order (D-WORLD-81): 150 → 151 → v1.0.0.
- **Row 150:** 22/23 executor milestones built (M8b left); P0, P3 and A1 remain.
- **Next:** row 150's C8′(iii)/P0 store-dependent-geometry amendment (doc r3, designer + quorum) → M8b. Row 230 sits below the release rows.
- **Harness:** row 144's ticket was resolved by the fleet (`d95ead6`), and the reply is acked.
- **Parked on Mark:** nothing — the ledger has 0 OPEN rows.
- **Cost:** metered $0.69 (quorum); Anthropic quota otherwise. Codex is over its ration.
- **Record/evidence:** world-mission-log.md entry 314; verification/world-iter314/.
