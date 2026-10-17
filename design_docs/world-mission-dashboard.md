# World mission dashboard — iteration 277 (2026-10-17)

- **This fire:** row 156 (`w-exec-protected-test-paths`) went from no doc to a design, a 19-milestone plan and its first milestone (PR #298). A profile may declare `protected` test paths. They are refused to World's write tools and denied in the srt sandbox (macOS both pattern forms, Linux directories only), and every exec result records whether the protected files still equal the base commit.
- **Key measured findings:** the row's example `**/*.test.ts` protects only one directory level in AILANG's matcher; listing matches as literal paths cannot stop new test files being created (6/6 written); AILANG's write check can be bypassed through a symlink (upstream ticket owed at landing).
- **M0a:** `m0_matrix.py --protected-self-test/--protected-only`, contained before any legacy side effect; CI runs the self-test. Independent judge sonnet r1 **PASS 82**, zero blocking.
- **Quorum:** two rounds, one external seat present (gemini-4-argon); OpenAI and both Ollama seats unreachable. r2's fix applied verbatim (carve-out).
- **Latest GitHub release:** `v1.0.0-rc.1` (pre-release, 2026-10-16). Proof/compiler pin AILANG 0.41.0; tool pin 0.52.1.
- **Clauses:** **7 of 7 MET.** Before `v1.0.0` (D-WORLD-81): finish row 156, then row 150 (9/17), then row 151.
- **Parked on Mark:** **D-WORLD-82** (row 156 design decisions D-156-1..4; defaults A, A, A, B; the plan follows the defaults).
- **Next:** row 156 M0b (Linux measurement on CI), then M1a.
- **Routing/cadence:** designer Agent opus (glm/kimi lanes on the Ollama weekly limit), planner and executor codex gpt-6.1-sol (recipe), judge Agent sonnet. 4h throttle.
- **Metered:** $0.076 (quorum). No shared skill, harness or routing-policy edit.
- **Evidence:** log entry 277; `design_docs/verification/world-row156-design-2026-10-17/`; role artifacts in `~/.ailang/state/world-iter277/`.
