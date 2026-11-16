# Iteration 396 — independent audit of `w-world-address-bar-verbs-v2.md` r1

- **Judge**: evaluator `sonnet` via the Agent tool (resolver `agent-tool sonnet declared:alias-pin`), read-only,
  in its own detached worktree at `39bb1ff7`. It ran for 2.8 min and used 102,419 tok. Generator ≠ judge: the
  designer was opus, the judge sonnet, and the quorum was cross-vendor (Google, Z-AI, Moonshot).
- **Verdict**: **PASS 84/100**, 0 blocking.
  - Evidence 28/30.
  - Consistency 19/25.
  - Soundness and safety 23/25.
  - Buildability 14/20.

## W-rows re-run (all reproduce)

- Required rows: W2, W3, W4, W5, W10, W11, W12, W19 and W39–W44.
- Further rows: W1, W6–W9, W15, W16, W18, W20–W26, W27 (approximate), W29–W31, W34, W35, W36 (truncated)
  and W37.
- Negative rows W19, W22, W23, W25, W26, W43 and W44 each carry a known-positive control in the same command.

## Non-blocking findings

The controller reproduced findings 1, 2, 3 and 5 first-party. The doc's §13.2 hands all of them to the planner.

1. **`sum` and `whyTarget` are used outside the walk.** `sum` (`why.go:774`) is called by
   `demo_seed.go:730,801,829,998,1064`. `whyTarget{kind:…}` (`why.go:297`) is built at
   `demo_seed.go:658`. Neither appears in the V2a adapter list.
2. **The two moving test files are not "package line only".**
   - `why_protected_test.go:52,59,74` uses cmd's rig helpers.
   - `TestWhyProtectedLine` mixes pure subtests with an `e2e-exec` subtest.
   - `why_input_test.go:31,58` calls `elide(p)`.
3. **`c.get` collides with an existing method.** `cli.go:94` already defines
   `func (c *client) get(path string) (int, string, error)`.
4. **VS2/T1's pinned `d.*` set is incomplete.** It omits `whyGetter`, the seven handler method values and any
   once field. This fails loudly under VS1's `ast.Inspect` + `DeepEqual` shape.
5. **A stale reason can survive into a result.** `host/verbs/lookup.go:99`'s "The app cannot run it yet." reason
   is copied into `a.Reason`. The route must clear it on a `result`.
6. **The 1 MiB cap counts raw bytes.** The JSON-escaped body can be larger. It is still bounded; RV2 should say
   so.
7. **Smaller items.**
   - The "8 non-test call sites" wording in §2.1.
   - VS2a pins direct imports only.
   - The audit's rough W27 re-run paired at `75 228`, against the doc's `56 209`.

## Candidates checked by the judge's own method

- **(a) The page's 15 s budget: TRUE.**
  - `app.js:484` sets `VERB_ROUTE_BUDGET_MS = 15000`.
  - `:512-519` arm a timer that sets `job.late` and calls `job.ctl.abort()`, and `fetch` takes `{signal}`.
  - On a late answer the page says *the verb route did not answer in time*.
  - On the server, `readCtx` derives from `r.Context()`.
- **(b) VS2's pinned `d.*` set: INCONSISTENT** (finding 4). It fails loudly, never silently.

## Soundness

- **Read-only: holds.** The in-process getter reaches only the seven unprotected GET handlers. The mux is private,
  and an unknown path is an error, never a 404.
- **Bounded: holds.** The limits are a line of ≤1,024 bytes, `--scan` ≤ 5000, a walk of ≤10 s, ≤8 MiB per read
  and an answer of ≤1 MiB.
- **No silent fallback: holds.** Rows (b) and (c) come before (f), so a deadline or a read failure never prints
  "not found" or ✗, and `why` never reaches `address`.
- **"Byte-identical" claims: consistent** with §2.1's health carve-out.
- **Amendments A1–A5:** each one is needed against the parent's text (parent lines 246, 355-356, 365 and 382).
- **§9 dispositions of the V1 judge's notes:** all reasonable.
