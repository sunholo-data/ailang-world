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

---
*Changes to this document are ratification-class (human gate). It is deliberately short —
every sprint reads it; token cost is a standing tax.*
