# World mission dashboard — iteration 237, 2026-10-06
Latest release: world/core@0.1.1 published attended; no release this run.
World 1.0: clauses 1/2/3/6/7 MET; 4/5 UNMET.
Critical path (D-WORLD-65): 136 LANDED → 141 → 152 → 153 → 93 → 114 → 139 → release.
Row 136 LANDED this fire: PR #219 → d83baad. MCP refusals logged as `mcp refusal` with tested record refs; `mcp:` ids + `coordinator:mcp`.
Upstream ask: sunholo-data/ailang#1602 (typed JSON-RPC error / isError from mcphttp). Row 159 watches it (Gate-1 predicate).
New rows: 159 (typed MCP error on the wire, gated on #1602), 160 (se-tools e2e never run in CI).
Next pick: row 141 w-workspace-project-layouts.
Parked on Mark: nothing — ledger ZERO OPEN.
Judge: opus Agent, PASS 92 → r2 96, zero blocking — design_docs/verification/world-iter237/.
Routing: planner Agent opus, executor Agent sonnet, evaluator Agent opus (openrouter/codex/ollama over ration). $0 metered.
Quorum: OpenAI seat (gpt6-1-sol) dead on API credits (429) — World quorums run at N−1.
Gate binary: v0.41.0 pin at ~/.pinned-ailang/ailang (not the v0.52.x tool pin).
Bookkeeping issue #202; full memory in charter/log/status archive.
