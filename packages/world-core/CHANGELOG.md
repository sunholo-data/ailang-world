# Changelog — world/core

## 0.1.0 — 2026-09-21

**The first publish of AILANG World's pure semantic core** — the four frozen modules the
world graph is defined in, projected deterministically from their canonical sources in this
repository and published through World's own propose → verify → commit pipeline (clause 7,
controlled self-modification proven once).

- **`world/types`** — the immutable, content-addressed world-graph vocabulary.
- **`world/contracts`** — `isValidNextWorld` and the validity predicates, Z3-proven.
- **`world/transitions`** — `applyRevision` and the pure transition functions; a transition
  is a value, so replay reconstructs state from pure transitions plus recorded effect
  results rather than from a log of mutations.
- **`world/logepoch`** — `sameRef` / `servesEntry`, the log-epoch semantics that let the
  interpreter version be pinned and transition functions be content-addressed (DESIGN.md
  open question 11), decided before the log format froze.

Effects: **none**. All 19 exported functions are pure, and the package declares an empty
effect ceiling — the core is a pure semantic kernel by construction, and the broker is the
only path to the outside world (clause 3).

Contracts: **11 verified, 0 refuted** under the pinned compiler. The projection is
byte-identical to its canonical sources, asserted by sha256 per module, and the published
tarball's identity is fixed by a committed ready-packet golden that the readiness gate
compares field by field before any publish is permitted.
