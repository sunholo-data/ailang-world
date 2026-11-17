# World mission dashboard — iteration 397 (2026-11-17)

- **Done this fire:** row 243 **V2a LANDED** (PR #585, merge `2d530cbe`; CI 3/3 on head and on the merge, first try). This is the seam for `why` in the World app's address bar.
  - The CLI's `why` walk now lives in `host/why` behind an injected `Getter` (`Run`, `ParseTarget`, `Head`, `Elide`, `Sum`). V2b's daemon will call the same function in-process.
  - The CLI keeps a thin front and a compatibility layer. `headState` and `elide` are one-line delegations, so there is one `/v1/head` reader.
  - Behaviour-neutral: `ailang-worldd why` prints the same bytes with the same exits. A characterisation golden generated at the base pins 11 paths no test drove before. The bar still answers `why` with its CLI line.
  - 1,299 changed lines of a 1,300 cap. Design §2.2 (`WhyFlags`, VA1) moved to V2b, measured at 111 lines.
  - Judge (sonnet, own worktree) **PASS 90/100**, 0 blocking. It also ran the real-stack why chains that CI skips (`WORLD_TOOL_AILANG_BIN` unset in CI); they passed.
- **Next loop work:** the row 243 **V2b** sprint planner (the bar runs `why`; ≈ 1,041 lines with §2.2 folded in). After that come V3 (the verb-approval law, with a ledger ask), V4 (host runner and herdr hand-off) and V5.
- **Attended, for Mark when convenient:** the V1 no-terminal exit run (parent design §10).
- **OPEN decisions:** none.
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Release rows 150 and 151 are not routable.
- **Routing this fire:**
  - Planner: opus via Agent (88 min).
  - Executor: haiku via Agent (6 min, one run, no stop).
  - Evaluator: sonnet via Agent (53 min).
  - Codex was over ration and unusable at preflight (probe rc 75). Metered $0.
- **Instrument notes:** self-notices still return rc 2 (`watermark channel unavailable`), because the #574/#524 watermark files are absent. Dev CI at `2d530cbe` is 3/3 green.
- **Bookkeeping thread:** #574.
