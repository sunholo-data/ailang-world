# World Coding Standards (ratified: Mark + coordinator, attended, 2026-07-24)

Binding on every agent writing code in this repo — designers, executors, evaluators, and
attended sessions alike. The evaluator scores against these; the CI gate enforces what it can.
Origin: the M1-milestone-1 retrospective — 285 LOC of kernel AILANG shipped using **none** of
AILANG's distinguishing features (0 contracts, 0 inline tests, vacuously-green Z3 gate), because
these standards existed only as unwritten house style. See `w-m1-ailang-hardening.md` for the
retrofit.

## S1 — Z3-verifiable pure core

Every **pure function in `world/` with a meaningful invariant** carries `requires`/`ensures`
contracts, written FIRST as the specification. Executable predicate functions (the
`contracts.ail` pattern) complement contracts as shared runtime checks — they never replace
them: a predicate nobody proves is decoration. Where Z3 cannot encode a function (string
interpolation/`show`, etc.), **named inline `tests [(in, exp)]`** are the fallback — every
canonical text form (renderers, keys) must have one or the other. The verify gate asserts a
**manifest of named check identities** (see the hardening doc) — never bare exit codes, never
aggregate counts alone: a green gate must be non-vacuous by construction.

## S2 — Effects at the boundary only

`world/` kernel modules are pure (`! {}`). Effects live in the host layer (Go) and, later, in
effect-handler extensions behind the broker. A `world/` function that wants an effect is a
design smell: what it actually wants is to *return a description* the host executes.

## S3 — Slim kernel, package-first (the self-modification enabler)

`world/` holds ONLY the frozen kernel surface: core types, transition laws, epoch registry,
contract predicates. Everything else — new domains, policies, tools, projections — lands as
**AILANG packages** (registry-published, version-pinned, cascade-updatable), because granular
packages are what make controlled self-modification viable: a World behavior change is a
package version bump through propose→verify→commit (DESIGN.md §14), never a kernel edit. Every
proposed addition to `world/` or `host/` must answer, in its design doc: **"why is this not a
package?"** Kernel growth without that answer is the OS gravity well (DESIGN.md §12.4) — the
named failure mode of this project.


### S3 ledger — language split & kernel size (dated; drift must be visible)

Measured per attended review (repro: `wc -l` over `world/*.ail`, `design_docs/**/*.ail`,
non-test vs test `.go` under `host/ cmd/`). The ratified architecture (DESIGN §19 q1) expects a
small `.ail` semantic core over a Go host, with the `.ail` share GROWING as: the MCP projection
registry (.ail `@route` modules), policies, domain transitions, and clause-7 extension packages
land — and as upstream Z3 encoder limits lift (ADT-record sort). A FALLING `.ail` share or a
kernel that grows without "why is this not a package?" answers is drift against S3.

| Date | `.ail` core (world/) | `.ail` checked docs | Go host (non-test) | Go tests | Note |
|---|---|---|---|---|---|
| 2026-08-04 | 434 | 809 | 6,773 | 10,506 | Post M1–M3+4b/4c: the Go-host build phase, as planned. All 4 core modules carry Z3-proven contracts; the broker's decision law is authored in `.ail` (7/7 verified) and transcribed to Go under a byte-verbatim drift test — AILANG is the semantic source of truth even where Go executes. Tests exceed production 1.55:1. |

*Watch-item (not yet built): S3 is evaluator-scored, not machine-checked — a kernel-surface
manifest gate (pin `world/` exports the way the required-check manifest pins contracts) would
make kernel growth mechanically loud. Candidate hardening item when kernel churn justifies it.*

## S4 — Compiler-checked docs

Every `.ail` snippet in `design_docs/` ships as a checkable file swept by the CI gate. A doc
claiming "this compiles" without a checked artifact is a defect.

## S5 — AILANG fluency protocol (before writing a line of `.ail`)

The M1 gap was discoverability, not capability. Before authoring or reviewing AILANG:
1. Load the language reference: the `ailang-docs` MCP server (`.mcp.json`) or `ailang prompt`.
2. Explore before inventing: `ailang builtins list` / stdlib search — check what exists before
   writing helpers.
3. **Verify every syntax claim against the pinned released binary** (`ailang check` /
   `ai-check` / `ailang test`) before it enters a doc or a prompt — never from memory; a wrong
   syntax claim propagates further than no claim.
4. Contracts are the point of the language: if your `.ail` uses no contracts, no inline tests,
   and no effect rows, justify why in the PR — the default assumption is you missed them.

## S6 — Honest gates

A gate that can pass vacuously will. Every gate must fail loudly on its own null case: zero
modules found, zero obligations proven, zero tests discovered, `verify.errors > 0` swallowed
by a zero exit code. When you find such a case (as `w-m1-ailang-hardening` did with
`ai-check`'s exit code), fix the gate AND file the upstream issue.

### Load-bearing criteria require mutations

A criterion that claims an assertion is load-bearing must be discharged by a **MUTATION** that
makes the test red, never by a static count of the assertion's own text. This binds both
design-doc acceptance criteria and sprint-plan test-plan rows: each load-bearing claim must name
the mutant and the single assertion that fires when it lands. A grep may appear only as an
explicitly labelled instrument-health control; it cannot discharge the load-bearing claim.

## S7 — Usage surfaces ship usage docs (the discoverability rule; added 2026-07-28, Mark)

Every user-facing surface (CLI verb, REST route, package API) lands WITH: (a) `--help` /
self-description; (b) a runnable quickstart covering the full happy path **including payload
construction** (`docs/QUICKSTART.md` or equivalent, executed-verbatim before commit); (c) a
working example or generator. **Knowledge that exists only in test files is a defect** — this
rule's third-instance origin: the M1 fluency gap, the ailang-feature gap, and the worldd commit
schema that even the coordinator had to reverse-engineer from `cli_test.go`. "The tests show
how" is not documentation; the evaluator scores it.

## S8 — The floor-raise coupling inventory (added 2026-08-27, row 43)

*A gate that is complete by construction can still have a coupling surface that is complete
only by memory, and the memory lives in whoever last raised the floor.* Raising the
verified-identity floor touches **six files**; at `P6.V` three roles each found a different
subset, and the full set existed only in commit `699f592`'s message. The map:

| # | File | What moves |
|---|---|---|
| 1 | `world/<module>.ail` | the new contract (the law itself) |
| 2 | `packages/world-core/world/<module>.ail` | projection copy — regenerate with `./scripts/build_world_package.sh`, never hand-edit |
| 3 | `scripts/verify_ail.sh` | BOTH constants: `REQUIRED_VERIFIED` and `EXACT_TOTAL_VERIFIED` |
| 4 | `scripts/world_package_ready_packet.golden.json` | `contentHash`, `tarballSHA256`, `tarballBytes` |
| 5 | `docs/SELF_MOD_PUBLISH.md` | the `contentHash` and `tarballSHA256` digest-table rows |
| 6 | `host/verifygate/module_manifest_gate_test.go` | the pristine-control marker string — hand-maintained; deriving it from `EXACT_TOTAL_VERIFIED` would make the control vacuous (S6) |

**`interfaceHash` does not move for `.ail` byte changes or packaged-module changes that
leave the hashed manifest fields unchanged. It MUST move when a hashed manifest field
changes, including `name`, `edition`, the optional AILANG bound, exports, or effects.**
So do not "fix" the third digest on a Tier-1 raise. A raise that changes the hashed manifest
fields (or the packaged-module census) touches additional sites beyond these six; that
inventory is deferred to a future item pending a first-party rehearsal.

Recipe, all six in the SAME commit: edit 1 and 3 by hand → `./scripts/build_world_package.sh`
(2) → run the pinned gate; step 9/9 reds printing a diff against the committed golden;
replace the golden with the diff's new line and re-run (4) *(step 9/9's red-arm diff format
verified via the V25 mechanism, not re-executed this session — see the V16 note)* → copy the two moved digests into the runbook table (5) → update the marker (6) →
re-run to green; `go test ./host/runbook/` binds 5↔4. Enforced by
`host/verifygate/floor_raise_inventory_test.go`.

---
*Changes to this document are ratification-class (human gate). It is deliberately short —
every sprint reads it; token cost is a standing tax.*
