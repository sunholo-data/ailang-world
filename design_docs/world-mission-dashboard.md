# World mission dashboard — iteration 398 (2026-11-17)

- **Done this fire:** row 243 **V2b LANDED** (PR #587, merge `c251132f`; CI 3/3 on head and 3/3 on the merge). Row 243 **V2 is complete**: you can type `why <target>` in the World app's address bar.
  - The bar's answer is byte-identical to `ailang-worldd why` (VE1, and a test that compares against the real CLI path rather than `why.Run`). `--json` works too.
  - It runs in-process through a private mux over the seven `/v1` GET handlers. That is read-only: no new route, no POST, no store write, no new daemon field.
  - Every failure is a refusal that carries the CLI line, never navigation: the 10 s deadline, a read failure, a cancelled walk, an answer over 1 MiB, and `-`/`--result`/`--expect-store`, which only the terminal can run.
  - The page now treats every Enter as owning the bar (an empty or third submit supersedes a pending ask), and refusals with a CLI line open the panel. Live Chrome XV2 drives it; QUICKSTART §19 documents it.
  - 1,255 changed lines of a 1,400 cap, in 5 sub-milestones, each equal to the plan.
  - Judge (sonnet, own worktree) **PASS 91/100**, 0 blocking. Notes: the CLI line's quoting is unpinned (controller reproduced), and three cosmetic page strings are unpinned.
- **Next loop work:** row 243 **V3** child design (the verb-approval law for mutating verbs). It will need a ledger ask. Designer = rotation (next after opus is `codex:gpt-6.1-sol`; check the resolver). After V3 come V4 (host runner and herdr hand-off) and V5.
- **Attended, for Mark when convenient:** the V1 no-terminal exit run (parent design §10), which can now include `why`.
- **OPEN decisions:** none.
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Release rows 150 and 151 are not routable.
- **Routing this fire:**
  - Planner: opus via Agent (61 min).
  - Executor: haiku via Agent (2 runs, 4.7 min, no stop).
  - Evaluator: sonnet via Agent (24 min, concurrent with the 14-min gate sweep).
  - Codex was unusable at preflight (probe rc 75, quota unreadable). Metered $0.
- **Instrument notes:** self-notices still return rc 2 (`watermark channel unavailable`), because the #574/#524 watermark files are absent. CI never runs the real-stack `why` tests (`WORLD_TOOL_AILANG_BIN` unset); the plan's gates and the judge did, and they passed.
- **Bookkeeping thread:** #574.
