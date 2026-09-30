# Changelog — world/core

## 0.1.1 — 2026-09-30

**Proof hardening of the pure semantic core. No exported signature changed** — the interface
hashes (v1 and v2, 53 signatures) are identical to 0.1.0, so every consumer of 0.1.0 compiles
unchanged against 0.1.1.

- **Five kernel exports now carry Z3-proven contracts:** `proposalMatchesWorld` and
  `verificationMatchesProposal` (exact field identity through `sameRef`), `commitAllowed`,
  `plan` and `verify`. Every added `ensures` restates what the unchanged body already
  computes; no function body changed.
- **Three exports are explicitly exempt, with the measured reason in the source:** the pinned
  compiler's specification encoder cannot state their exact laws — `renderRef` (string
  interpolation in `ensures`), `cacheKey` and `commit` (each calls a user function the encoder
  skips as not SMT-encodable). Named inline tests and the package smoke exercise them; they
  are recorded exemptions, not claimed proofs.
- **The smoke program gained a rejected-commit arm**: an unaccepted verification must be
  denied, so a commit that ignored the verdict fails the package's own smoke run.

Contracts: **16 verified, 0 refuted**, uncontracted exports **8 → 3**. Effects: still
**none**.

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
