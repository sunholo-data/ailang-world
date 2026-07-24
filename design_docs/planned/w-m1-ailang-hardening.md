# w-m1-ailang-hardening — Z3 Contracts + Inline Tests for the M1 Kernel

**Status**: Planned
**Date**: 2026-07-24
**Charter clause**: clause-1
**Mission**: ailang-world
**Queue item**: `w-m1-ailang-hardening` (Mark-directed, queued iter-10; preempts w-world-library-m1 M4)
**Verified against**: **released `AILANG v0.30.0`** (`/tmp/ailang-v0300/ailang`, commit `e37b370`,
built 2026-07-19) + **Z3 4.16.0** — the same pinned artifact as the M1 sprint. Every contract and
test shape in this doc was run through that exact binary before being specified (see the
Verification Log). NOT verified against the PATH dev build. Claim V5 was found FALSE against
production types on 2026-07-24 (executor iter-12 + controller iter-13); the achievable proven set
is 4, and upstream `sunholo-data/ailang#477` tracks the ADT-in-record encoder gap.
**Traces to**: [w-world-library-m1.md](w-world-library-m1.md) (M1 milestone 1 shipped the four
modules this doc retrofits); [world-mission.md](../world-mission.md) queue item; memory
`ailang-feature-discoverability-gap`
**Depends on**: w-world-library-m1 milestones 1–3 (LANDED); no dependency on M4–M6
**Estimated**: ~0.5–1 day
**Priority**: P0 (cheaper now than after M4/M5 build on these types)

> **Scope note.** Retrofit only. Three deliverables: (1) proven `requires`/`ensures` contracts on
> the int/bool/record invariants of `world/{contracts,transitions,logepoch}.ail`, (2) inline
> `tests [(in,exp)]` on the pure functions Z3 cannot encode (the string renderers), (3) a
> **non-vacuous** verify gate that asserts a **hardcoded required-check manifest** of
> proven-contract and passing-test *identities* (not aggregate counts) from the
> `ai-check` / `ailang test` JSON. **OUT of scope, deliberately:** effects (M1 is intentionally
> pure `! {}`), package extensions (frozen core, clause 7 later), any Go host change (M2/M3 are
> tested and evaluator-passed).

---

## Problem Statement

The M1 flagship AILANG showcase (`world/{contracts,transitions,logepoch,types}.ail`, 285 LOC,
PR #2) shipped using **none** of AILANG's distinguishing features:

- **0 Z3 contracts.** Contracts 1–4 exist only as *decorative* `bool` predicates — executable,
  but nothing machine-proves they hold. Contract 4 ("the resulting world increments revision by
  exactly one") is checked by no one: `isValidNextWorld` is exported and never called.
- **0 inline tests.** `ailang test --format json world/` on the current tree reports
  `total_tests: 0` (verified 2026-07-24). The string renderers `renderRef`/`cacheKey` — the
  canonical on-disk text forms the whole content-addressing story rests on — have **no machine
  check of any kind**: they use string interpolation → builtin `show`, which has no SMT encoding,
  so Z3 skips them *and* no tests exist.
- **A vacuous gate.** `scripts/verify_ail.sh` runs `ai-check` on every module and passes on
  exit code alone. Baseline run (pinned binary, 2026-07-24): all four world modules report
  `verified=0, skipped=0, errors=0, results=[]`, exit 0. The gate proves module *count*, not
  verification. Worse: `ai-check` **exits 0 even when a contract produces a Z3 encoding
  error** (`verify.errors > 0` — empirically confirmed), so naively adding contracts would not
  even make failures visible without parsing the JSON.

Root cause was discoverability (see `ailang-feature-discoverability-gap`), now partially
corrected (`.mcp.json` + upstream issue `sunholo-data/ailang#476`). This doc closes the gap in
the code itself.

## Goals

1. `ai-check` reports `status == "verified"` for **every identity in a required-check
   manifest**: `(transitions, applyRevision)`, `(contracts, isValidNextWorld)`,
   `(logepoch, sameRef)`, and `(logepoch, servesEntry)` — 4 in total, with the exact total
   (`== 4`) as a *secondary* check only.
2. `ailang test` passes **all 14 named inline tests** (`renderRef_test_1/2`, `sameRef_test_1/2`,
   `cacheKey_test_1/2`, `servesEntry_test_1/2`, `proposalMatchesWorld_test_1/2`,
   `verificationMatchesProposal_test_1/2`, `commitAllowed_test_1/2` — per-test identities
   confirmed available in v0.30.0 JSON, V18/V25), covering every interpolation-based function
   and every Proposal-taking predicate whose designated machine check is a test.
3. `scripts/verify_ail.sh` **fails loudly** when any required identity is missing, errored,
   counterexampled, or skipped, or when any required named test is missing or failing — with
   the manifest and totals **hardcoded** (no env override can weaken the gate; only
   `AILANG_BIN`, the binary path, stays configurable) and CI-safe.

### Non-goals

- No effects, no capability enforcement, no package extensions (charter: frozen core).
- No Go host changes. No new modules beyond the four that shipped in M1 milestone 1.
- No attempt to prove the string renderers in Z3 — v0.30.0 cannot encode interpolation; tests
  are their designated check.

---

## Empirical Verification Log (pinned v0.30.0 + Z3 4.16.0, 2026-07-24)

Every language claim below was executed on `/tmp/ailang-v0300/ailang` against `/tmp` scratch
modules mirroring the real `world/*` types and module graph (multi-module, repo-root-relative
paths — same shape as production). **V-numbers are cited throughout the doc.**

| # | Claim | Evidence |
|---|-------|----------|
| V1 | Contract syntax: `requires {…}` / `ensures {…}` go **between** the effect annotation and the body; `ai-check` reports `verify.results[].status == "verified"` | scratch `applyRevision` → `"status": "verified"` |
| V2 | Record-field postconditions prove: `ensures { result.revision == w.revision + 1 && result.stateRoot.algo == outputWorld.algo && … }` (5-conjunct) **verified**, including with `World`/`HashRef` imported cross-module | multi-module scratch mirroring `world/logepoch` + `world/types` + `world/transitions` |
| V3 | **A contracted function may not call ANY user function** — not cross-module, not same-module, not even an SMT-encodable one. Z3 fails with `unknown constant sameRef (HashRef HashRef)` → `status: "error"`. The v0.30.0 encoder never inlines callees. Provable contracts require **fully-inlined bodies** | body `sameRef(…)` + inlined ensures → `error`; same body inlined → `verified` |
| V4 | `implies` is **not** AILANG syntax (`PAR_UNEXPECTED_TOKEN` at the keyword). Exact-equality form `ensures { result == (<full condition>) }` is the provable shape and entails the implication corollaries | parse-error transcript; exact form verified on all 4 predicates |
| V5 | **SUPERSEDED 2026-07-24 — FALSE against production types.** Only `isValidNextWorld` verifies. The three Proposal-taking predicates Z3-error with `unknown sort 'Proposal'` because production `Proposal` contains `evidence: list[Evidence]`, where `Evidence` is an ADT; `ai-check` exits 0 silently | executor iter-12 + controller iter-13 pinned-binary reproductions against scratch mirroring the real production types |
| V6 | `sameRef` with a **field-equality body** + `ensures { result == (left.algo == right.algo && left.digest == right.digest) }` → `"verified"` (its current interpolation body is unencodable) | scratch logepoch mirror |
| V7 | Any function whose body reaches interpolation (`renderRef`, `cacheKey`, current `sameRef`) is **unencodable**; contracted callers get `UNENCODABLE_TYPE` skip or Z3 error. Tests are the only machine check | scratch3 skip transcript: `calls user function "sameRef" that is not SMT-encodable` |
| V8 | `plan` **cannot** carry a contract in v0.30.0: the encoder mis-sorts literals in record construction — `confidence: 1.0` → SMT `Int` vs declared `Real`; `expectedEffects: []` → `(Seq Int)` vs `(Seq String)` → `unknown constant mk_Proposal` error | scratch `planlike` error transcript |
| V9 | An **uncontracted** function returning a sum type and calling predicates + the proven helper (`commit`'s retrofit shape) produces **no Z3 obligation, no error** — absent from `verify.results` entirely | scratch `commitlike` |
| V10 | `ai-check` exit codes: all-verified → **0**; counterexample → **1**; parse error → **1**; **Z3 encoding error (`verify.errors>0`) → 0 (SILENT)**. The gate must parse JSON, not trust exit codes | four exit-code transcripts |
| V11 | `ai-check` has **no** `--json` flag and needs none: JSON is its native output (`--help` lists only `-relax-modules`, `-timeout`, `-verify-recursive-depth`). JSON shape: `.check.passed`, `.verify.{verified,counterexample,skipped,errors,results[]}` | `--help` output + every transcript above |
| V12 | Inline test syntax verified: single record arg `tests [(({ algo: "sha256", digest: "abc" }), "sha256:abc")]`; two record args `((({…}), ({…})), true)`; nested `HashRef`-in-`LogHeader` literals; `list[string]` literals in `EpochRecord`. **Contracts + tests coexist** on one function (order: signature → `! {}` → `requires`/`ensures` → `tests […]` → body) | scratch9/scratch10: `check: true`, tests pass |
| V13 | `ailang test`: has `--format json` (fields `total_tests`, `passed_tests`, `failed_tests`, `skipped_tests`, `success`); failing test → exit 1; **0 tests found → exit 1** (`--allow-skips` to override); directory mode `ailang test world/` aggregates across modules | `--help` + run transcripts |
| V14 | Contract-derived property tests over **record-typed parameters skip** (`no generator for parameter rec: EpochRecord`) — expected noise; passing inline tests alongside skipped properties still exit 0 | scratch9: 4 pass + 1 prop-skip → exit 0 |
| V15 | Baseline (real repo, pinned binary): all four world modules `verified=0, errors=0, results=[]`, exit 0; `ailang test` finds 0 tests. Today's gate is fully vacuous | baseline transcript 2026-07-24 |
| V16 | Negative-existence: sketches are self-contained (`design_docs/sketches/*.ail` import `sketches/*`, never `world/*`) — the `sameRef` retrofit cannot affect them; the Go host (`host/`) shares no code with the `.ail` modules (no FFI in M1), so no Go behavior changes | grep of `sameRef\|renderRef` across the repo |
| V17 | `ai-check` emits **bare function names** in `verify.results[].function` (e.g. `"sameRef"`, no module qualifier). `verify_ail.sh` checks one module per invocation, so a manifest keyed by **(module file, bare name)** is unambiguous | scratch `world/logepoch` retrofit mirror (contracts+tests exactly as D3): `results: [sameRef, servesEntry]`, `verified: 2`, exit 0 |
| V18 | `ailang test --format json` **exposes per-test identities**: `tests[].name` follows `<function>_test_<N>` with `status: "pass"/"fail"` (+ `error` detail on fail); contract-derived properties are listed separately in `properties[]` with `status: "skip"`. A failing test still emits the **full JSON** alongside exit 1 | 8-test scratch run (all 8 names present, pass) + deliberate-fail run (`renderRef_test_1` → `"fail"`, `"expected WRONG, got sha256:abc"`, exit 1, complete JSON) |
| V19 | `ailang test` prints a human banner (`→ Running tests in …`) on **stdout ahead of the JSON** (stderr empty) — the test-leg parser must strip lines before the first `{`. `ai-check` stdout is pure JSON (every transcript above parsed directly) | bash stream-separation run: banner + JSON both in the stdout capture, stderr empty |
| V20 | **NT1 analog (silent identity vanish):** deleting `sameRef`'s `ensures` makes the identity **disappear from `results[]` entirely** — `verified` 2→1, `errors: 0`, exit 0. Invisible to exit codes and maskable by any aggregate floor; only a required-identity check catches it | scratch mutation: ensures line stripped → `results: [servesEntry]`, exit 0 |
| V21 | **NT2 analog (aggregate floor blind spot):** deleting both `renderRef` tests → `passed_tests: 6, failed_tests: 0`, exit 0 — an aggregate floor of 6 **still passes**; the named-test manifest fails (both `renderRef_test_*` absent) | scratch mutation: tests block removed → 6 names, no `renderRef_test_*`, exit 0 |
| V22 | Directory mode `ailang test --format json world/` merges all modules into **one JSON preserving `tests[].name`** (names stay bare — no module qualifier) and exits 0 when a test-free module sits alongside tested ones; a test-free module **alone** exits 1 but still emits valid JSON (`total_tests: 0`, `success: false`) | two-module scratch (`logepoch` + test-free `types`): dir-mode exit 0 with all 8 names; `types.ail` solo → exit 1, valid JSON |
| V23 | A contract on `proposalMatchesWorld(w: World, p: Proposal)` with an inlined field-equality body + exact `ensures` reports function `status: "error"`, reason `Z3 error ... Invalid constant declaration: unknown sort 'Proposal'`, `verify.errors: 1`, process exit 0 (SILENT). `verificationMatchesProposal` and `commitAllowed` are the same class because both take `Proposal` | controller iter-13 pinned-binary scratch mirroring the REAL types |
| V24 | A contract on `isValidNextWorld(w: World, next: World, outputWorld: HashRef, nextLogHead: HashRef)` with an inlined field-equality body + exact `ensures` reports `status: "verified"`, `verify.errors: 0`, exit 0; `World`/`HashRef` contain no ADT | controller iter-13 pinned-binary scratch mirroring the REAL types |
| V25 | Inline `tests [(in,exp)]` on a Proposal-taking predicate run and pass using a full 9-field `Proposal` literal (`evidence: []`, `confidence: 1.0`, `requiredCaps: []`, `expectedEffects: []`): `proposalMatchesWorld_test_1/2` both `status: "pass"`, `failed_tests: 0`. With all three Proposal predicates tests-only (no `ensures`) plus the proven `isValidNextWorld` contract, `ai-check` reports `verified: 1, errors: 0, counterexample: 0` — a clean gate | controller iter-13 pinned-binary scratch mirroring the REAL types |
| V26 | **Bounded-execution semantics (re-quorum objection).** `ai-check -timeout` is a **per-function Z3 timeout** (default 5s — `--help`: "Per-function Z3 timeout"), **NOT** a process-wide wall-clock bound; `ailang test` has **no** timeout flag at all. Neither leg self-bounds its wall-clock, so a solver/runner/parse hang would block CI indefinitely — the gate must wrap BOTH invocations in an OS-level bounded subprocess that kills the process group on expiry and fails loudly (Standing Rule 6) | controller iter-13 `ai-check --help` + `test --help` on the pinned binary |

---

## Solution Design

### D1 — The sum-type resolution: `applyRevision`, a proven helper `commit` composes

**The key design question.** `commit` returns `CommitResult = Applied(World, Transition) |
Denied(string)`. You cannot write `ensures { result.revision == w.revision + 1 }` on a
sum-typed result, and v0.30.0 offers no per-constructor postcondition. Additionally, `commit`'s
body calls `commitAllowed` — a user function — so by V3 *any* contract on `commit` itself would
Z3-error regardless of the result type.

**Resolution (empirically verified, V2 + V9):** extract the next-world construction into a pure
helper carrying Contract 4's core as a Z3-proven postcondition, and have `commit` call it:

```ailang
-- world/transitions.ail (NEW helper, above commit)
-- Contract 4's core, machine-proven: the next world increments revision by
-- exactly one and carries the supplied output-state root and next-log-head.
export func applyRevision(w: World, outputWorld: HashRef, nextLogHead: HashRef) -> World ! {}
requires { w.revision >= 0 }
ensures {
  result.revision == w.revision + 1
    && result.stateRoot.algo == outputWorld.algo
    && result.stateRoot.digest == outputWorld.digest
    && result.logHead.algo == nextLogHead.algo
    && result.logHead.digest == nextLogHead.digest
}
{
  { revision: w.revision + 1, stateRoot: outputWorld, logHead: nextLogHead }
}
```

`commit` keeps its exact signature and `CommitResult` return; only the `Applied` world literal
is replaced by the helper call:

```ailang
  if commitAllowed(w, p, v)
  then Applied(
    applyRevision(w, outputWorld, nextLogHead),
    { proposalHash: p.proposalHash, inputWorld: p.inputWorld,
      outputWorld: outputWorld, transitionFn: p.transitionFn, evidence: [] }
  )
  else Denied("verification contract failed")
```

The increment invariant is Z3-proven at the helper; `commit` (uncontracted) composes it and
produces no Z3 obligation and no error (V9). This is the strongest provable shape v0.30.0
admits for a sum-typed commit.

### D2 — contracts.ail: one proven predicate; Proposal predicates get inline tests

Only `isValidNextWorld` gets a Z3 contract. Its `World`/`HashRef` parameters contain no ADT, so
its body is rewritten as inlined field equality and its `ensures` restates the **exact** semantics
(`ensures { result == (<full condition>) }`, V4/V24). The inlining rationale remains V3/D-A: a
contract that must be Z3-encodable takes the strictly-safe path independent of callee encodability.

The three Proposal-taking predicates — `proposalMatchesWorld`,
`verificationMatchesProposal`, and `commitAllowed` — keep their **existing shared-predicate
bodies unchanged**, including calls to `sameRef` and to each other. They gain 2 inline `tests`
each and carry **no `ensures`**. A contract on any of them reaches production `Proposal`, whose
`evidence: list[Evidence]` transitively contains a user ADT, and fails with
`unknown sort 'Proposal'` while `ai-check` exits 0 silently (V23; upstream
`sunholo-data/ailang#477`). Their machine check is therefore the six named inline tests (V25).

**The drift objection, answered.** contracts.ail's stated design is "verify and commit call the
SAME predicates so policy cannot drift." That anti-drift shared-predicate architecture is
preserved for all three Proposal predicates: their bodies remain shared-predicate calls, and the
callers (`verify`, `commit`) remain unchanged. Only `isValidNextWorld` is inlined because its
contract is Z3-encodable and proven; where the v0.30.0 encoder cannot form the obligation, inline
tests provide the machine check without duplicating predicate logic.

### D3 — logepoch.ail: `sameRef` becomes structural (and proven); renderers get tests

`sameRef` currently compares canonical *text* forms (`renderRef(l) == renderRef(r)`), which is
unencodable (V7). Retrofit to field equality with a proven contract (V6):

```ailang
export func sameRef(left: HashRef, right: HashRef) -> bool ! {}
ensures { result == (left.algo == right.algo && left.digest == right.digest) }
tests [
  ((({ algo: "sha256", digest: "abc" }), ({ algo: "sha256", digest: "abc" })), true),
  ((({ algo: "sha256", digest: "abc" }), ({ algo: "sha256", digest: "def" })), false)
]
{
  left.algo == right.algo && left.digest == right.digest
}
```

This is a **deliberate, tiny semantic strengthening**, called out in the Conflict Surface:
text-form equality conflates refs whose fields straddle a `:` (e.g. `{algo:"a", digest:"b:c"}`
vs `{algo:"a:b", digest:"c"}` both render `"a:b:c"`). Field equality is the documented intent
("Structural equality") and the finer relation. No producer emits colon-bearing algos (`sha256`
only, per the settled epoch decision), so no observable behavior changes on real data.

`renderRef` and `cacheKey` keep their interpolation bodies (they ARE the canonical text form)
and get inline tests — their only possible machine check (V7). `servesEntry` is pure int
equality, so it gets both a proven exact `ensures` and a test (V5-pattern + V12).

### D4 — Per-function retrofit table

Machine-check legend: **Z3** = `ai-check` reports `verified` (empirically confirmed on the
pinned binary); **tests** = inline `tests [(in,exp)]` under `ailang test`; **—** = none possible
in v0.30.0 (reason given).

| Function | Module | `requires` / `ensures` | Inline tests | Machine check |
|---|---|---|---|---|
| `applyRevision` **(new)** | transitions | `requires { w.revision >= 0 }`; `ensures { result.revision == w.revision + 1 && result.stateRoot.algo == outputWorld.algo && result.stateRoot.digest == outputWorld.digest && result.logHead.algo == nextLogHead.algo && result.logHead.digest == nextLogHead.digest }` | — (record params ⇒ property skips, V14) | **Z3 ✅** (V2) |
| `commit` | transitions | none — sum-typed result + calls user fns (V3); composes the proven helper | — (Proposal/Verification literals heavy; stretch only) | via `applyRevision` + `commitAllowed` |
| `plan` | transitions | none — encoder mis-sorts `1.0`/`[]` literals in record construction (V8) | stretch | — (documented v0.30.0 limitation) |
| `verify` | transitions | none — body calls shared predicate (V3) | stretch | — |
| `proposalMatchesWorld` | contracts | none (Proposal contains ADT `evidence` ⇒ `unknown sort`, V23/#477); existing shared-predicate body unchanged | 2: `proposalMatchesWorld_test_1`, `proposalMatchesWorld_test_2` | **tests only** (V25) |
| `verificationMatchesProposal` | contracts | none (Proposal contains ADT `evidence` ⇒ `unknown sort`, V23/#477); existing shared-predicate body unchanged | 2: `verificationMatchesProposal_test_1`, `verificationMatchesProposal_test_2` | **tests only** (V25) |
| `commitAllowed` | contracts | none (Proposal contains ADT `evidence` ⇒ `unknown sort`, V23/#477); existing shared-predicate body unchanged | 2: `commitAllowed_test_1`, `commitAllowed_test_2` | **tests only** (V25) |
| `isValidNextWorld` | contracts | exact ensures: `result == (next.revision == w.revision + 1 && stateRoot field-eq outputWorld && logHead field-eq nextLogHead)`; inlined field-eq body | — | **Z3 ✅** (V24) |
| `sameRef` | logepoch | exact ensures over field equality; body retrofitted from text-eq (D3) | 2 (equal / differing digest) | **Z3 ✅ + tests** (V6, V12) |
| `renderRef` | logepoch | none — interpolation (V7) | 2: `({algo:"sha256",digest:"abc"}) → "sha256:abc"`, `({algo:"sha1",digest:"00"}) → "sha1:00"` | **tests only** |
| `cacheKey` | logepoch | none — interpolation via `renderRef` (V7) | 2: full `LogHeader` literal → `"sha256:aa@sha256:bb"`; distinct interpreter digest → distinct key | **tests only** (V12: nested literals confirmed) |
| `servesEntry` | logepoch | `ensures { result == (rec.epoch == h.semanticsEpoch) }` | 2: matching epoch → `true`, mismatched → `false` | **Z3 ✅ + tests** |
| `types.ail` (all) | types | none — type declarations only, no functions | — | type-check (unchanged) |

**Totals: 4 Z3-proven contracts, 14 inline tests.** The gate (D5) requires each of these **by
identity, not by count**: `(transitions, applyRevision)`, `(contracts, isValidNextWorld)`,
`(logepoch, sameRef)`, and `(logepoch, servesEntry)`, plus all 14 required passing test names,
with the exact totals (`== 4`, `== 14`) as secondary checks only.
Required invariants and renderer coverage are not interchangeable counts (V20/V21 demonstrate
both maskings empirically). A benign refactor that renames or moves one of these functions must
update the hardcoded manifest in the same PR — a deliberate, reviewable act rather than a
silent gate weakening.

### D5 — Non-vacuous gate: required-check manifest in `scripts/verify_ail.sh`

Current gate: per-module `ai-check`, trusts exit codes, passes on module count alone (V15) —
and exit codes are **insufficient** because a Z3 encoding error exits 0 (V10). Aggregate
floors are also insufficient: a dropped contract *vanishes silently* from `results[]` (V20)
and a deleted test pair still clears a floor of 6 (V21) — required invariants and test coverage
are not interchangeable counts. The gate is therefore a **required-check manifest**: hardcoded
per-module identity sets that must verify/pass, with exact totals as secondary checks. Retrofit
the script (same file, same ROOTS mechanism, same AILANG_BIN contract):

**Gate policy — hardcoded constants, deliberately NOT env-overridable.** Only `AILANG_BIN`
remains configurable, and it selects the binary *path*, not gate strength. There are no
`MIN_*` knobs: production CI cannot lower the requirements via environment variables.
Keyed by **(module file, bare function name)** because `ai-check` emits bare names in
`verify.results[].function` (V17) and the script checks one module per invocation:

```python
REQUIRED_VERIFIED = {
    "world/transitions.ail": {"applyRevision"},
    "world/contracts.ail":   {"isValidNextWorld"},
    "world/logepoch.ail":    {"sameRef", "servesEntry"},
    "world/types.ail":       set(),
}
REQUIRED_TESTS = {   # logepoch (8) + contracts (6)
    "renderRef_test_1", "renderRef_test_2", "sameRef_test_1", "sameRef_test_2",
    "cacheKey_test_1", "cacheKey_test_2", "servesEntry_test_1", "servesEntry_test_2",
    "proposalMatchesWorld_test_1", "proposalMatchesWorld_test_2",
    "verificationMatchesProposal_test_1", "verificationMatchesProposal_test_2",
    "commitAllowed_test_1", "commitAllowed_test_2",
}
EXACT_TOTAL_VERIFIED = 4   # secondary check, world/ modules only
EXACT_TOTAL_TESTS    = 14  # secondary check
GATE_LEG_TIMEOUT_S   = 120 # hardcoded wall-clock cap per ai-check module leg (Rule 6, V26)
GATE_TEST_TIMEOUT_S  = 180 # hardcoded wall-clock cap for the directory-mode test leg (V26)
```

**Leg 0 — bounded execution (hardcoded deadlines, non-env-overridable, V26).** `ai-check
-timeout` bounds only individual Z3 queries and `ailang test` has no timeout at all (V26), so a
solver/runner/parse hang would block CI indefinitely — a bounded-waits violation (Standing Rule
6). Every binary invocation in **both** legs runs through a hardcoded-deadline wrapper that starts
the child in its own process group and, on expiry, kills the **whole group**, prints a named
`✗ TIMEOUT …` to STDERR, and exits `124`. The deadlines are constants (`GATE_LEG_TIMEOUT_S`,
`GATE_TEST_TIMEOUT_S`) — no environment variable can extend or disable them (only `AILANG_BIN`
stays configurable). Portable form (python3 is already required by the parser; `start_new_session`
→ process-group kill):

```bash
run_bounded() {  # $1=timeout_s  $2=out_file  $3..=cmd ;  exit 124 + named msg on expiry
  local t="$1" out="$2"; shift 2
  python3 - "$t" "$out" "$@" <<'PY'
import os, signal, subprocess, sys
t = int(sys.argv[1]); out = sys.argv[2]; cmd = sys.argv[3:]
with open(out, "wb") as f:
    p = subprocess.Popen(cmd, stdout=f, stderr=subprocess.STDOUT, start_new_session=True)
    try:
        sys.exit(p.wait(timeout=t))
    except subprocess.TimeoutExpired:
        os.killpg(os.getpgid(p.pid), signal.SIGKILL)
        sys.stderr.write("✗ TIMEOUT after %ds: %s\n" % (t, " ".join(cmd))); sys.exit(124)
PY
}
```

A `124` from `run_bounded` is **always fatal** (it is the one exit code that is NOT advisory —
distinct from the Z3-error-exits-0 and counterexample-exits-1 classes the JSON parse handles).
Keep `ai-check`'s per-function `-timeout` at its 5 s default as an inner belt so individual Z3
queries also self-bound.

**Leg 1 — ai-check manifest.** For each module, capture the JSON to a temp file and parse
(python3 — present on macOS and ubuntu-latest runners; no jq dependency). Exit codes are
**advisory** (`|| true` / RC captured) — the JSON parse is authoritative, since a Z3
encoding error exits 0 (V10) and a vanished identity exits 0 (V20) — **except a `124` TIMEOUT
(Leg 0), which is fatal**:

- **every** swept module (world/ + sketches): fail if `.check.passed != true`,
  `.verify.errors > 0`, or `.verify.counterexample > 0`;
- **world/ modules**: fail unless every name in `REQUIRED_VERIFIED[module]` appears in
  `.verify.results[]` with `status == "verified"` — a missing identity, `error`,
  `counterexample`, or `skip` for a required name each fail with a message naming the
  (module, function) pair;
- accumulate `total_verified` over **world/ modules only**; after the loop assert
  `total_verified == 4` exactly (secondary). Sketches are excluded from the total and carry
  empty required sets — a future contracted sketch can neither mask a required identity
  (requirements are per-module) nor perturb the exact total.

```bash
total_verified=0                                # explicit init before the loop
for mod in "${ROOTS[@]}"; do
  mod=${mod#./}                                 # normalize path so keys match the manifest (gemini catch)
  run_bounded "$GATE_LEG_TIMEOUT_S" "$tmp_json" "$AILANG_BIN" ai-check -timeout 5s "$mod"; rc=$?
  [ "$rc" -eq 124 ] && { echo "✗ ai-check TIMEOUT on $mod (>${GATE_LEG_TIMEOUT_S}s)"; exit 1; }
  # other exit codes advisory (V10/V20) — the JSON parse below is authoritative
  mod_verified=$(python3 - "$mod" "$tmp_json" <<'PY'
  # loads REQUIRED_VERIFIED (hardcoded above); validates check/errors/counterexample
  # for all modules and required identities for world/ modules; prints the module's
  # verified count to STDOUT on success; on any failure, prints the named-identity
  # message to STDERR (objection-2 fix, gemini re-quorum) and exits 1 — so the message
  # bypasses the $() stdout capture and stays visible when the `|| exit 1` fires.
PY
  ) || exit 1
  case "$mod" in world/*) total_verified=$((total_verified + mod_verified));; esac
done
if [ "$total_verified" -ne 4 ]; then
  echo "✗ expected exactly 4 proven world/ contracts, got $total_verified"; exit 1
fi
```

**Leg 2 — named inline tests.** v0.30.0's test JSON **does expose per-test identities**
(`tests[].name`, V18) — so the gate takes the stronger of the two paths the quorum
anticipated: required *named* tests, not per-module counts. One directory-mode run (names
are preserved and merged across modules, V22):

```bash
run_bounded "$GATE_TEST_TIMEOUT_S" "$tmp_test_json" "$AILANG_BIN" test --format json world/; rc=$?
[ "$rc" -eq 124 ] && { echo "✗ ailang test TIMEOUT (>${GATE_TEST_TIMEOUT_S}s)"; exit 1; }
# other exit codes advisory; JSON parse authoritative
# python3 parser:
#  - strip non-JSON prefix lines before the first "{" (stdout banner, V19)
#  - fail unless EVERY name in REQUIRED_TESTS appears in .tests[] with status == "pass"
#    (message to STDERR names each missing/failing test identity — objection-2 fix)
#  - fail if .failed_tests > 0
#  - assert len(.tests[]) == 14 exactly (secondary) — NOT .passed_tests (see fixture
#    correction D-B: .passed_tests also counts contract-derived PROPERTIES and would be
#    flaky; the count of NAMED inline tests is stable)
```

The `|| true` (objection-2 fix) interacts with the manifest by design: a failing test exits 1
but still emits full JSON (V18), so the parser reports *which* identity failed; a crashed or
zero-test run yields absent/failing JSON fields, which the parser fails loudly. Property-test
**skips are tolerated** (record-param generators don't exist in v0.30.0, V14); assertions
are on named `tests[]` entries and `failed_tests`, never on `skipped_tests` or
`properties[]`. Test names are bare, not module-qualified (V22). Two test-bearing modules now
exist (`world/logepoch` and `world/contracts`) with no collision because the seven function
names are distinct; if a future collision arises, leg 2 must switch to the documented
per-module-run escape hatch (recorded as a comment in the script).

**Negative acceptance tests** (one-shot, recorded in the PR description — not permanent CI
jobs): **NT1** — a scratch copy of `world/` with `applyRevision`'s contract stripped must
fail the gate naming the missing `(transitions, applyRevision)` identity (V20 confirms the
aggregate-invisible vanish this catches). **NT2** — a scratch copy with both `renderRef`
tests removed must fail the gate naming the missing `renderRef_test_1`/`renderRef_test_2`
identities even though the other 12 required tests still pass (V21 confirms an aggregate floor
would miss the removed identities). **NT3** — pointing `AILANG_BIN` at a mock that `sleep`s
longer than `GATE_LEG_TIMEOUT_S` must make the ai-check leg terminate at the deadline with a
named `✗ ai-check TIMEOUT` message and a nonzero exit — and no environment variable extends it
(the bounded-waits teeth, V26). **NT4** — the same mock against the directory-mode test leg must
be terminated at `GATE_TEST_TIMEOUT_S` with a named `✗ ailang test TIMEOUT` and nonzero exit.
Each mutation fails independently.

**CI safety:** `AILANG_BIN` stays the only knob (CI exports its checksum-verified install);
the manifest and exact totals are constants in the script; no hardcoded machine paths;
python3 invoked as `python3 -` heredoc (no new files) and its absence fails the script
loudly rather than passing vacuously. Sketches under `design_docs/` are still swept by leg 1
(must type-check with `errors == 0`) and are exempt from identity requirements and totals.

---

## Conflict Surface

This change touches AILANG source parsed by the pinned toolchain and the CI gate script. The
honest enumeration:

1. **Contract/tests syntax position vs the existing grammar.** Contracts occupy the slot
   between the effect annotation and the body; `tests […]` sits after `ensures` and before the
   body. Both positions verified parse-clean on v0.30.0 **in combination** (V12). Risk: the
   current M1 functions mostly omit `! {}` (e.g. `plan`, `verify`, `commit` today have bare
   `-> T {` bodies). Contracted functions will gain an explicit `! {}` (the verified shape,
   V1); uncontracted functions are left byte-identical. No other grammar position is touched.

2. **The ADT return of `commit` (the load-bearing constraint).** No postcondition can be
   attached to `CommitResult` (sum type), and no contract can live on `commit` at all because
   its body calls user functions (V3). The design routes the invariant through `applyRevision`
   (D1). **What is and isn't proven must be stated honestly:** Z3 proves *the helper*
   increments the revision and proves `isValidNextWorld`; `commitAllowed` is tests-only because
   its `Proposal` parameter reaches an ADT (V23). That `commit`'s `Applied` arm *uses* the helper
   is enforced by code review + the `isValidNextWorld` proven predicate remaining available to
   the Go host — not by Z3. A future AILANG with per-constructor postconditions could close this
   gap; out of scope here.

3. **No-user-calls encoding rule vs the shared-predicate architecture.** Inlining field
   equality is limited to `isValidNextWorld`, whose contract must be Z3-encodable (V24).
   `proposalMatchesWorld`, `verificationMatchesProposal`, and `commitAllowed` keep their
   shared-predicate bodies unchanged, preserving contracts.ail's anti-drift architecture; their
   machine check is 2 inline tests each, with no `ensures` because Proposal contracts Z3-error
   (V23/V25). The `verify`/`commit` → shared-predicate call structure is unchanged.

4. **`sameRef` semantic change (text-eq → field-eq).** Deliberate strengthening (D3). Differs
   only for refs whose `algo` contains `:` — impossible for producer-generated refs (`sha256`
   fixed by the settled epoch decision) but *representable* by the type. Downstream check:
   sketches are self-contained and don't import `world/logepoch` (V16); the Go host shares no
   code with the `.ail` modules (V16); `world/contracts.ail` keeps its shared-predicate call sites.
   `renderRef`'s output format is untouched — the on-disk/in-log canonical
   text form does not change, so nothing the Go host has persisted is affected.

5. **v0.30.0 encoder limitations bound the scope.** `plan` cannot be contracted (literal-sort
   encoder errors, V8); `verify` cannot (calls predicates, V3); and a record parameter that
   transitively contains a user ADT cannot be encoded (`unknown sort 'Proposal'`, V23), tracked
   upstream as `sunholo-data/ailang#477`. None is silently dropped: the table records the reason,
   and the mission's upstream channel (`sunholo-data/ailang#476` thread) remains the escalation
   path for the plan/verify gaps (function-call inlining, float/empty-list literal sorts, record
   generators for property tests).

6. **Gate change blast radius.** `verify_ail.sh` gains JSON parsing, a hardcoded
   required-check manifest, and a second leg. Failure modes considered: (a) Z3-error-exits-0
   is caught by `.verify.errors > 0` (V10), and the *silent identity vanish* — a dropped
   contract disappearing from `results[]` with exit 0 (V20) — is caught by the per-identity
   requirement, which no aggregate count could see; (b) exit codes are now **advisory**
   (`|| true`), so the JSON parse is the authoritative gate — a crashed run or zero-test run
   yields absent/unparseable JSON or missing required identities, both of which fail loudly;
   leg 2 runs only on `world/` (sketches carry no tests and are not swept by leg 2); (c) the
   test runner's stdout banner precedes the JSON (V19) — the parser strips non-JSON prefix
   lines before decoding; (d) property skips must not fail the gate (V14) — asserted on named
   `tests[]` entries and `failed_tests` only, never `skipped_tests`; (e) python3 absence would
   break the gate — it is present on macOS dev machines and ubuntu-latest CI runners; the
   script fails loudly if `python3` is missing rather than passing vacuously; (f) **manifest
   coupling is deliberate**: renaming or moving a required function breaks the gate until the
   hardcoded manifest is edited in the same PR — visible, reviewable friction in place of the
   silent weakening that downward-overridable floors permitted; sketches can never mask a
   required identity (requirements are per-(module, function) and the `== 4` total is scoped
   to `world/`); (g) bare test names could in principle collide across future test-bearing
   modules (V22) — `logepoch` and `contracts` now both carry tests with no collision; a future
   collision forces leg 2 to the documented per-module-run fallback, noted in a script comment.

7. **Programs that MUST still work after this change** (regression fixtures, all existing):
   - `world/types.ail` — byte-identical, still type-checks (`ai-check` exit 0).
   - `world/logepoch.ail`, `world/contracts.ail`, `world/transitions.ail` — `ai-check` exit 0
     with `errors == 0, counterexample == 0` and the new `verified` counts;
     `world/contracts.ail` also carries 6 passing inline tests.
   - `host/...` Go tests and the CI `go-verify` job — untouched files, must stay green.
   - `design_docs/sketches/*.ail` — still swept by leg 1, still pass with 0 contracts.
   - CI workflow's `AILANG_BIN` export — unchanged contract, no workflow edit required unless
     the workflow pins the script's output format (it checks exit code only).

---

## Implementation Plan (single sprint, ~0.5–1 day)

**Phase 1 — logepoch (~1h).** D3: `sameRef` field-eq body + ensures + 2 tests; `renderRef` +
`cacheKey` tests; `servesEntry` ensures + tests. Run `ai-check` + `ailang test` on the module.

**Phase 2 — contracts (~1h).** D2: `isValidNextWorld` → inlined body + exact ensures; the three
Proposal-taking predicates keep their shared-predicate bodies, gain 2 inline tests each, and carry
no `ensures`. Expect `verified: 1` (`isValidNextWorld`), `errors: 0`, plus 6 passing inline tests
(`proposalMatchesWorld`/`verificationMatchesProposal`/`commitAllowed`, 2 each).

**Phase 3 — transitions (~1h).** D1: add `applyRevision`, rewire `commit`'s `Applied` arm.
Expect `verified: 1` and **no** entry for `commit` in `verify.results` (V9).

**Phase 4 — gate (~1–2h).** D5: retrofit `verify_ail.sh` (Leg 0 `run_bounded` hardcoded-deadline
wrapper, required-check manifest, JSON parse with banner-strip, test leg, exit-codes-advisory
except the fatal `124` TIMEOUT, path normalization `mod=${mod#./}`). Run the negative tests
NT1/NT2 (manifest teeth) **and NT3/NT4 (bounded-waits teeth — sleeping `AILANG_BIN` mock)**,
then full gate + CI.

Executor note (mission requirement): load the version-locked syntax first (`ailang-docs` MCP
`prompt_get` or `ailang prompt`) before touching `.ail`.

## Acceptance Criteria

- [ ] `scripts/verify_ail.sh` green on the pinned binary with **every required identity
      verified**: `(transitions, applyRevision)`, `(contracts, isValidNextWorld)`,
      `(logepoch, sameRef)`, `(logepoch, servesEntry)` — each present in `ai-check`
      `verify.results[]` (bare names, V17) with `status == "verified"`, and world/-scoped
      total exactly **4** (secondary check).
- [ ] `ailang test --format json world/` green with **all 14 required named tests passing**:
      `renderRef_test_1/2`, `sameRef_test_1/2`, `cacheKey_test_1/2`, `servesEntry_test_1/2`,
      `proposalMatchesWorld_test_1/2`, `verificationMatchesProposal_test_1/2`, and
      `commitAllowed_test_1/2` each `status == "pass"` in `tests[]` (V18/V25),
      `failed_tests == 0`, and `len(tests[]) == 14` exactly (secondary — NOT `passed_tests`,
      which also counts passing contract properties → flaky; fixture correction D-B); property
      skips tolerated (V14).
- [ ] `verify.errors == 0` and `verify.counterexample == 0` on every module (the silent-error
      class, V10, is now gate-fatal).
- [ ] **Gate policy is hardcoded:** no environment variable can lower, bypass, or shrink the
      required-identity manifest or the exact totals; `AILANG_BIN` (binary path) is the only
      configurable knob.
- [ ] **Bounded waits (Standing Rule 6, V26):** both gate legs run through the hardcoded-deadline
      `run_bounded` wrapper (Leg 0) that starts the child in its own process group and SIGKILLs the
      whole group on expiry; NT3/NT4 (a sleeping `AILANG_BIN` mock) prove each leg terminates at
      `GATE_LEG_TIMEOUT_S`/`GATE_TEST_TIMEOUT_S` with a named TIMEOUT failure + nonzero exit, and
      that no environment variable can extend or disable the deadline.
- [ ] **Negative test NT1 (required identities have teeth):** running the retrofitted gate
      against a scratch copy of `world/` with `applyRevision`'s contract stripped exits
      non-zero naming the missing `(transitions, applyRevision)` identity — even though 3
      other functions still verify (the aggregate-invisible vanish, V20). Recorded once in
      the PR description; not a permanent CI job.
- [ ] **Negative test NT2 (named tests have teeth):** running the gate against a scratch copy
      with both `renderRef` tests removed exits non-zero naming the missing
      `renderRef_test_1`/`renderRef_test_2` identities — even though 12 tests still pass and
      `failed_tests == 0` (V21). Each mutation fails independently of the other. Recorded
      once in the PR description; not a permanent CI job.
- [ ] CI green end-to-end (ai-check gate + go-verify job untouched).
- [ ] No diffs outside `world/{logepoch,contracts,transitions}.ail` and
      `scripts/verify_ail.sh`. `world/types.ail` byte-identical.

## Quorum verification log + carve-out application (iter-11)

- **Designer:** `claude:claude-fable-5` (rotation, subscription). **Controller:** opus.
- **Round 1 quorum** (gpt5-6-sol reject, gemini-3-1-pro pass) → BLOCKED. Objection: the
  aggregate-floor gate (`MIN_VERIFIED=6`) was semantically weak (a dropped contract still cleared
  the floor). **Gate-mandated Fable revision** rewrote D5 to the hardcoded required-check manifest.
- **Round 2 re-quorum** (gpt5-6-sol reject, gemini-3-1-pro reject) → BLOCKED on two NEW **narrow,
  direction-preserving** objections with concrete `proposed_fix`: (gpt5-6-sol) the V-log wasn't
  auditable — commit a verification fixture directory; (gemini-3-1-pro) Leg-1 `$()` capture swallows
  the python error — route it to stderr.
- **NARROW-REFINEMENT CARVE-OUT applied** (ratified for the world mission at the M1 GO, attended;
  see `world-mission-status-archive.md`). Both objections satisfy (a) concrete verbatim
  `proposed_fix` and (b) no design-direction dispute → the controller made this bounded 2nd revision
  applying the reviewers' VERBATIM fixes:
  1. **gpt5-6-sol (auditability):** committed `design_docs/verification/w-m1-ailang-hardening/`
     (fixtures + `run.sh` + captured `OUTPUTS.md`, pinned binary sha256 `e9746fef…`). Claims not
     reproducible on scratch are relabeled as sprint-acceptance checks on the production tree (the
     reviewer's own fallback clause).
  2. **gemini-3-1-pro (stderr):** Leg-1 python error routed to stderr; leg-2 likewise.

### Post-fixture corrections (evidence-backed, from `OUTPUTS.md`)

Running the reviewer-demanded fixtures caught two first-draft inaccuracies — precisely the value the
auditability objection was raised to capture. Both are corrected above; the design DIRECTION is
unchanged.

- **D-A — V3 refined (not a design change).** `callsUserFn` (a contract whose body calls the user
  function `sameRef`) **verified** in `encoder_limits.ail.txt` — so a contract CAN call a user
  function when the callee's body is SMT-encodable. V3's "may not call ANY user function" is
  corrected to: a contract calling an **unencodable-bodied** callee (interpolation) errors; an
  encodable-bodied one can verify. The decision to **inline** the predicate bodies still stands — it
  is the strictly-safe choice independent of callee encodability, and `verified_baseline.ail.txt`
  proves the inlined predicates verify.
- **D-B — leg-2 secondary count is `len(tests[]) == 8`, not `passed_tests == 8`.** `--format json`
  counts passing contract-derived **properties** in `passed_tests` (`inline_tests.ail.txt`:
  `passed_tests == 7`, `total_tests == 8`), so `passed_tests` is flaky. D5 + the acceptance criteria
  now assert the count of named inline `tests[]`. `failed_tests == 0` stays a hard check.

### iter-13 revision (7->4 descope, empirically forced)

- Executor iter-12 and controller iter-13 falsified V5 on 2026-07-24 against production types
  using the pinned `/tmp/ailang-v0300/ailang` binary: a contract on any Proposal-taking predicate
  errors with `unknown sort 'Proposal'` because `Proposal.evidence` contains the `Evidence` ADT,
  while `ai-check` exits 0 silently. Upstream `sunholo-data/ailang#477` tracks this encoder gap.
- Controller iter-13 reproduced that `isValidNextWorld` verifies cleanly on the real
  `World`/`HashRef` shape and that inline Proposal-predicate tests run and pass with full
  production-shaped literals. The achievable scope is therefore 4 Z3-proven contracts:
  `applyRevision`, `isValidNextWorld`, `sameRef`, and `servesEntry`.
- `proposalMatchesWorld`, `verificationMatchesProposal`, and `commitAllowed` move to tests-only:
  their shared-predicate bodies remain unchanged, they carry no `ensures`, and they gain 2 inline
  tests each. New totals: **4 verified / 14 tests**.
- The design DIRECTION is unchanged: proven contracts wherever the encoder allows, inline tests as
  the machine check everywhere else, and a hardcoded non-vacuous required-check manifest gate.

### iter-13 re-quorum + narrow-refinement carve-out (2nd revision)

- **Re-quorum** on the descoped doc (`gpt5-6-sol`, `gemini-3-1-pro`; controller verdict pass;
  metered $0.095): **gemini-3-1-pro PASS** (only non-blocking notes: future JSON-schema coupling of
  `len(tests[]) == 14`; a `mod=${mod#./}` path-normalization catch — the latter folded into Leg 1).
  **gpt5-6-sol REJECT**, one blocking objection: the gate invokes `ai-check` and `ailang test`
  with **no enforced wall-clock deadline**, so a solver/runner hang blocks CI indefinitely —
  a bounded-waits (Standing Rule 6) violation. Concrete `proposed_fix`: hardcoded non-overridable
  deadlines on both legs; establish `ai-check -timeout` semantics; wrap in a process-group-killing
  subprocess timeout; add sleeping-mock negative tests.
- **NARROW-REFINEMENT CARVE-OUT applied** (already ratified for the world mission, iter-11 M1 GO,
  attended). The objection (a) carries a concrete reviewer-authored `proposed_fix` and (b) does NOT
  dispute the design DIRECTION — it is a determinism/robustness completeness fix (bounded waits).
  Controller made this bounded 2nd revision applying the reviewer's VERBATIM fix:
  1. **V26** establishes empirically that `ai-check -timeout` is a **per-function** Z3 timeout (not
     process-wide) and `ailang test` has no timeout — so both legs need an OS-level bound.
  2. **Leg 0 `run_bounded`** wrapper: hardcoded `GATE_LEG_TIMEOUT_S`/`GATE_TEST_TIMEOUT_S`
     (non-env-overridable), `start_new_session` process group, SIGKILL-on-expiry, named TIMEOUT to
     STDERR, exit `124` (the one fatal, non-advisory code). Both legs route through it. Controller
     pre-validated the mechanism on 4 cases (fast→rc0 output captured; 30 s hang→rc124 killed at
     the 2 s deadline; nonzero exit preserved ≠ timeout; process-group kill leaves no leaked child).
  3. **NT3/NT4** (sleeping `AILANG_BIN` mock) prove each leg terminates at its deadline; an
     acceptance criterion added.
- Not a force-pass (Standing rule 2): the design direction was uncontested; the fix is the
  reviewer's own. Routed straight to sprint-planner after this revision (no 3rd quorum).

## Related Documents

- [w-world-library-m1.md](w-world-library-m1.md) — parent M1 kernel design (milestones 1–3
  LANDED). This doc hardens its milestone-1 output; it does not overlap M4–M6 scope.
- [w-log-epoch-decision.md](w-log-epoch-decision.md) — SETTLED; fixes `sha256` + tagged
  `HashRef`, which is why the `sameRef` strengthening (D3) is observationally safe.
- Mission queue item + memory `ailang-feature-discoverability-gap` — root-cause context and
  the upstream escalation channel for the encoder gaps found here (V3, V8).

*(The skill's neural related-doc search is unavailable in this repo — no `ailang docs search`
corpus; the two docs above are the complete planned/ set and were reviewed directly.)*
