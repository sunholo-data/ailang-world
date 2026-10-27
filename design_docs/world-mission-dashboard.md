# World mission dashboard — iteration 309 (2026-10-27)

- **This fire:** row 150. M7c2b is BUILT. `scripts/capture_marketing.sh` now runs the whole capture in one command: Chrome version, a supervised `ailang-worldd serve` (frozen config, both proof identities, the operator's tool path), routes, two six-shot passes, a refusal if Chrome changed mid-run, diagnostics, and the verifier. The daemon is stopped on every exit path.
- **New flags:** `--worldd-bin`, `--bind` (127.0.0.1 only, refused before anything starts), `--marketingcheck-bin`. `marketing_capture.py daemon` now takes exactly one of `--ready-file` or `--ready-url`.
- **Controller amendments:** A1, daemon lifetime 660 s derived from the production constants, so pass B cannot outlive the daemon. A2, a `bind-refused` subtest.
- **Judge:** sonnet PASS 83/100, zero blocking; the plan's 3 named mutants were killed.
- **Controller:** the executor's sandbox could not bind loopback, so every gate was re-run outside it (green, including `-race`). The controller's drill killed 8 of 12 further mutants. Two survivors were closed by test-only subtests (`single-pass`, `daemon-exits`).
- **Proven against fakes only.** The real-Chrome drill is M7c4.
- **Roles (Agent tool):** executor codex gpt-6.1-sol through an Agent supervisor running the pinned recipe; evaluator sonnet. Designer and planner not needed.
- **Clauses:** 1–7 MET; goal unmoved. Release order (D-WORLD-81) is 156 → 150 → 151.
- **Row 156:** 22/24; D-WORLD-86 holds M5c/M5d.
- **Row 150:** 19/23 built. Pending: M7c3, M7c4, M8a, M8b. Controller steps: R1, P0, P3. Attended: A1.
- **Next:** M7c3 (frozen public inventory guard), then M7c4.
- **Carried (judge, non-blocking):** A1's lifetime and `ready()`'s stop branch are untested; missing/unknown-flag and port refusals have no subtest, and shell usage errors exit 1, not 2. Also carried from 307: `pass` creates `--out-dir` before it validates routes.
- **Cosmetic residuals (for Mark, if wanted):** the phone title is clipped with horizontal overflow, the provenance payload clips horizontally, and the tall image proportions are unusual.
- **Parked on Mark:** D-WORLD-86. A: drop the alias pre-checks and rely on pin+M5e (recommended); B: defence in depth. Default is HOLD.
- **Cost:** metered $0.00; codex and Anthropic quota only.
- **Record/evidence:** world-mission-log.md entry 309; verification/world-iter309/.
