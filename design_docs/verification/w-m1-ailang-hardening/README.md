# Verification fixtures — `w-m1-ailang-hardening`

These fixtures make the design doc's empirical claims **auditable and re-runnable**, per the
design-quorum objection (gpt5-6-sol, iter-11 re-quorum) that "V1–V22 cite ephemeral scratch modules
… with no repository paths, exact commands, captured outputs, hashes, or artifacts." They were added
by the controller as part of the **NARROW-REFINEMENT CARVE-OUT** application (ratified for the world
mission at the M1 GO; see `world-mission-status-archive.md`).

## How to run

```bash
AILANG_BIN=/tmp/ailang-v0300/ailang ./run.sh   # writes OUTPUTS.md
```

Fixtures are committed as `fixtures/*.ail.txt` (NOT `.ail`) on purpose: `scripts/verify_ail.sh`
sweeps `*.ail` under `design_docs/`, and several fixtures are *designed to fail* (encoder error,
counterexample). `run.sh` copies each to a temp dir as `.ail`, runs the pinned binary, and captures
exit codes + JSON into `OUTPUTS.md` (committed). Binary identity is recorded in OUTPUTS.md header
(v0.30.0, commit `e37b370`, sha256 `e9746fef…`).

## What each fixture backs (see OUTPUTS.md for the captured JSON)

| Fixture | V-rows | Captured result |
|---|---|---|
| `verified_baseline.ail.txt` | V1, V2, V5, V9, V17 | exit 0; `verify.verified == 5` (`applyRevision` + the 4 contracts.ail predicates, all `"verified"`); the uncontracted sum-returning `commit` is **absent** from `results[]` (V9); identity strings are **bare** (`"applyRevision"`, not module-qualified) → V17 confirmed. |
| `encoder_limits.ail.txt` | V8, V10-silent, **V3 (refined)** | exit **0** despite `verify.errors == 1`: `planLike` → `unknown constant mk_Proposal (Int (Seq Int))` (V8 float/empty-list mis-sort, exact); the silent-error class (V10) confirmed — a Z3 encoding error does **not** change the exit code. **See discrepancy D-A below.** |
| `counterexample.ail.txt` | V10 | exit **1**; `verify.counterexample == 1`; model `x=0, result=(+ x 2)`. Pins the exit-code table: verified→0, counterexample→1, encoding-error→0. |
| `inline_tests.ail.txt` | V12, V18, **leg-2 gate** | exit 0; 6 named inline tests all pass (`renderRef_test_1/2`, `sameRef_test_1/2`, `servesEntry_test_1/2`) with `tests[].name` identities (V18); a contract + tests coexist on `sameRef`/`servesEntry` (V12). **See discrepancy D-B below.** |

## Discrepancies the fixtures surfaced (the auditability objection was correct)

Running the fixtures caught two inaccuracies in the doc's first-draft claims. Both are corrected in
the doc (§ "Post-fixture corrections", iter-11); recorded here for the audit trail.

- **D-A — V3 was overstated.** The doc claimed "a contracted function may not call ANY user
  function." But `callsUserFn` (whose body calls the user function `sameRef`, and whose `ensures`
  restates `sameRef`'s field-equality) reports **`"verified"`** in `encoder_limits.ail.txt` — the
  encoder DOES handle a call to a user function whose body is itself SMT-encodable (field equality).
  V3's true, narrower claim: a contract calling a user function whose body is **unencodable**
  (interpolation, e.g. the original `sameRef`/`renderRef`) errors; calling an encodable-bodied one
  can verify. **The design decision to inline the predicate bodies still stands** — it is the
  strictly-safe choice (independent of callee encodability) and it is what `verified_baseline` proves
  — but the justification is corrected from "cannot call any user function" to "inlining removes the
  callee-encodability dependency."
- **D-B — leg-2 secondary count fixed from `passed_tests` to `len(tests[])`.** `ailang test
  --format json` counts passing **contract-derived properties** inside `passed_tests`
  (`inline_tests.ail.txt`: `passed_tests == 7` = 6 inline tests + 1 passing `servesEntry_property`;
  `total_tests == 8`, `skipped == 1`). So an exact `passed_tests == 8` secondary assertion would be
  FLAKY (it moves when a contract's property generator status changes). The gate's leg-2 secondary
  check is therefore corrected to assert on the count of named inline tests — `len(.tests[]) == 8`
  and every required name present with `status == "pass"` — never on `passed_tests`. `failed_tests
  == 0` remains a hard check.

## Claims verified on the PRODUCTION tree at sprint time (relabeled as acceptance checks)

Per the reviewer's own fallback ("Claims not backed by those artifacts must be relabeled as
assumptions and made acceptance checks"), the following are verified by the sprint against the REAL
`world/*.ail` modules (not scratch) and their captured JSON committed in the sprint PR: the exact
7-verified manifest across the three real modules, the 8 named inline tests in `world/logepoch`, and
the two negative gate tests (strip `applyRevision`'s contract → fail; delete both `renderRef` tests →
fail). These are the acceptance criteria in the design doc.
