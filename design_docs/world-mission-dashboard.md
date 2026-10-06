# World mission dashboard — iteration 238, 2026-10-06
Latest release: world/core@0.1.1 published attended; no release this run.
World 1.0: clauses 1/2/3/6/7 MET; 4/5 UNMET.
Critical path (D-WORLD-65): 136 LANDED → 141 LANDED → 152 → 153 → 93 → 114 → 139 → release.
Row 141 LANDED this fire: PR #224 → 6eeff52. `--workspace-module-root` (subdirectory sandbox) + `--workspace-package-cache` (read-only registry snapshot); QUICKSTART §11.
Upstream asks: sunholo-data/ailang#1607 (pkg-docs egress under a policy), #1608 (read-only package root). #1602 still open (row 159).
New rows: 162 (pkg_docs unrecorded egress), 163 (lock `path` entries outside the sandbox — measure first).
CI: TestExecSrtMCP* time out in verify_go.sh's parallel -race step on 7/9 attempts since ~14:00Z (pass in the -p 1 srt step) = row 153's defect (attended filed it as hygiene row 161 mid-fire). dev HEAD 6eeff52 green on rerun.
Next pick: row 152 (default) — or row 153 if Mark answers D-WORLD-68 = A.
Parked on Mark: D-WORLD-68 (promote row 153 ahead of 152? A recommended; default B).
Judge: opus Agent, PASS 91 → r2 96, zero blocking — design_docs/verification/world-iter238/.
Routing: designer/planner/evaluator Agent opus, executor Agent sonnet (codex/ollama/openrouter over ration). $0.58 metered (quorum).
Quorum: gemini/glm/kimi seated; OpenAI seat (gpt6-1-sol) dead on API credits.
Gate binary: v0.41.0 pin at ~/.pinned-ailang/ailang (not the v0.52.x tool pin).
Bookkeeping issue #202; full memory in charter/log/status archive.
