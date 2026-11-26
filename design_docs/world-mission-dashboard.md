# World mission dashboard — iteration 437 (2026-11-26)

- **Done this fire:** row 243 **V4 sprint 5e3 LANDED** (PR #674, merge `22bb7804`). Merge CI went 3/3 on attempt 2. The judge gave **PASS 79** with 0 blocking findings.
  - **File publication has landed, after two deferrals.** A signed start's task file, its mint token and Claude's MCP bundle are written under the episode lock after one step reservation. Each file is written privately with exclusive no-replace.
  - Each file is read back independently before the step is marked done. Anything uncertain is resolved by the attended `ailang-worldd service runner reconcile-files ACTION-ID`.
  - Sessions minted from a #671 standing allowance are refused.
  - 1,650 lines, inside the 1,650 exception the directive allowed for this one unit.
  - Nothing dispatches yet: publication has no production caller.
- **Judge's open points (non-blocking):** 23 of its own 32 mutants survive their named tests (9 killed). They are mostly store admission checks that no `host/store` test exercises; the controller reproduced one (row expiry). None of these paths has a production caller, and the two directory-swap windows fail safe to `uncertain`. These become S5e4's pre-caller pins.
- **Instrument defect caught:** the plan's gate script put a `git` wrapper on PATH, and tests with a scrubbed environment broke it. It faked 43 extra `host/broker` base failures. The controller fixed it in scratch and re-recorded the base.
- **CI headroom warning:** `cmd/ailang-worldd` under `-race` took 477.6 s and 469.0 s on two slow CI runners. The bound is 480 s; the two prior dev runs took 220 s and 302 s, and this sprint added about 51 s. It is recorded on row 267, so new pins go in `host/store`.
- **Before your next `service upgrade` of the home service:** back up `credentials.db`.
- **DECISIONS FOR MARK:** none open (ledger 107 rows, 0 OPEN).
- **Next loop work:** **V4 S5e4**, the pre-caller pins plus the #667 GC skip, before S5f's first production caller:
  - judge-437's store predicates, as `host/store` tests
  - the `reconcile-mint` after files-reservation refusal (N33/N34)
  - judge-436's five store survivors
  - the #667 GC skip for operator worktrees

  Row 271 (14 local-only e2e tests expecting 10 se-tools) is still small.
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Goal unmoved; row 243 serves D-WORLD-94.
- **Routing this fire:** planner and executor `codex:gpt-6.2-sol` (430k + 78k tok). The judge was Agent `sonnet` (133k tok), a different vendor. Metered $0.
- **Dev CI:** green at the start (`eaa202b1`). After the merge (`22bb7804`), attempt 1 was red on row 267 plus one load-sensitive test; attempt 2 was 3/3.
- **Open PRs:** #617 is attended UI work that the loop does not touch.
- **Bookkeeping thread:** #635 (next rotation Monday 2026-11-30 07:00 local).
