# World mission dashboard — iteration 393 (2026-11-16)

- **Done this fire:** row 260 LANDED ([#573](https://github.com/sunholo-data/ailang-world/pull/573), merge `10e6b5dc`). Judge (sonnet) PASS 94/100, zero blocking.
  - The kernel's contracts type-check under AILANG v0.60.4 again: four `TC_CONTRACT_001` sites (two enums gain `deriving (Eq)`, credited to #565; `plan`/`verify` restate `HashRef` equality field-wise).
  - **Dev's verify gate is green** on the head and the merge. The merge also landed record 392 (#572), so row 201 A4 sprint 4 counts as LANDED.
  - v0.60.4 is recorded as an observed release; the gate's "UNRECOGNISED RELEASE" notice is gone.
- **Ready for Mark:** a `world/core` **0.1.2** candidate (contentHash `359a82fd…`, interface hashes unchanged, kind `feature`). Not published. **D-WORLD-117** asks whether to publish (runbook: design §7 / `docs/SELF_MOD_PUBLISH.md` Stage B). The published 0.1.1 fails `TC_CONTRACT_001` under v0.60.4.
- **PR #565** (fleet account, not this mission's) is superseded by #573 and was left untouched.
- **Next loop work:** row 243 (address-bar verbs, new design r0, then quorum).
- **Attended, pending Mark:** exit runs for A1b (QUICKSTART §15), A2 (§16), A3a (§17), A4 (§18 / design §10); the 0.1.2 publish.
- **OPEN decisions:** D-WORLD-105 (default C), D-WORLD-114 (A, applied), D-WORLD-115 (A, applied), D-WORLD-116 (A, built this fire), **D-WORLD-117 (new)**.
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Release rows 150 and 151 are not routable.
- **Routing this fire:** designer opus (Agent, r0 + one revision), quorum ×2 ($0.52 metered), planner opus (Agent, 32 min), executor haiku (Agent ×2), evaluator sonnet (Agent, 17 min). Codex unusable at preflight (probe rc 75).
- **Capacity:** fires paused 2026-11-13 18:32 → 2026-11-16 02:31 (Anthropic weekly limit, codex quota unreadable, pi demoted); two 11-13 fires crashed at gate 3 (Cloudflare 403, weekly limit). Anthropic reset 11-16 06:00.
- **Instrument notes:**
  - Merge CI: `go host build + test gate` red on attempt 1 (row 259 flake, instance 5) and attempt 2 (runner could not download go1.26.5 for the race-gate repro); attempt 3 went 3/3.
  - Judge non-blocking: the runbook's publish phrase is unpinned prose (controller-reproduced; fails safe).
  - The self-notice script still reports `watermark channel unavailable`.
- **Bookkeeping thread:** #574 (rotated today from #524).
