---
title: Known residuals
sidebar_position: 3
description: The honest list of known gaps in World's authority and confinement story, each with its owner and, where filed, the upstream issue.
---

# Known residuals

A residual is a known gap that has been measured, named and given an owner instead of being
hidden. This page lists the ones an operator or security reviewer should know about. The
software-engineering residuals (`R-SE-*`) come from §9 of
`design_docs/planned/w-software-engineering-domain.md`.

## Software-engineering tools (`R-SE-*`)

| ID | Residual | Owner |
|---|---|---|
| R-SE-1 | **Overhead.** Each call costs a lock, one or two capsule runs, a subprocess (measured 0.03–0.19 s) and a SQLite commit; same-episode calls serialise; `ailang-run` is capped at 10 s (the policy's own timeout is 8 s) where the reference lane allows 120 s, so long programs can time out. The clause-4 risk is retired: in the 2026-10-15 floor FINAL the World arm's median overhead was −17.9% (Claude Code) and −5.7% (codex), inside the +25% bound. | Closed for 1.0 by row 93 (D-WORLD-77); the 10 s cap remains |
| R-SE-2 | **Ergonomics.** No grep, glob or list tools beyond `ailang-cli tree`, and results carry `world` provenance fields. The floor held regardless (69/69 per agent in both arms). | Row 93 measured it; no open owner |
| R-SE-3 | **Arm parity.** Both benchmark arms must use the same AILANG version; the `examples-search` corpus differs from the reference lane's checkout scan. Held for the floor FINAL by its pre-registration. | Row 93 (closed, D-WORLD-77) |
| R-SE-4 | **Cache writes inside the worktree.** `run --policy` ignores `AILANG_CACHE_DIR` and writes its compile cache into `<worktree>/.ailang/`. Integrity is held by `.ailang/**` in the deny list. | Upstream `sunholo-data/ailang`: the supervised worker should honour `AILANG_CACHE_DIR` |
| R-SE-5 | **No `Git.Commit` tool.** Git commit runs repo config and hooks in a worktree the agent writes. | A later git-effects row, which must disable hooks and config |
| R-SE-6 | **At-least-once under retry.** A retried MCP call creates a new task and may repeat an effect; every execution is recorded. | Row 108's MCP retry-channel residual |
| R-SE-7 | **Worktree state is not in World's store.** Files are reconstructible only from write and edit payloads. Since row 184, each `ailang-run`, `ailang-check`, `ailang-cli` and `workspace-exec` result records a `worktree_delta` (paths, sizes and sha256 of the files the call created, modified or deleted; gitignored paths excluded), but not the file contents. | A future snapshot-effect row |
| R-SE-8 | **Tool-binary drift.** Every confinement fact is measured on one hash-pinned binary; the daemon refuses others. If a future binary fails the matrix, stay on the last one that passed; never widen. | Re-proved by the AC4.1 matrix on each new binary |
| R-SE-9 | **Two file-effect namespaces.** `FS.*` (per exact file, direct Go I/O, no production caller) and `Workspace.*` (per worktree, policy-layer confinement) coexist. | A later effects-consolidation row |
| R-SE-10 | **Single-effect plans** and no conflict retry in v1. | A later row |
| R-SE-11 | **Inner effects of `ailang-run` are policy-governed, not brokered.** One `Ailang.Run` per run regardless of what the program does; bounded by the policy's `FS=1000` budget and timeout; not individually journaled or replayable. Owner-ratified (D-SE-3 = A). | A later row brokering inner effects (needs a live bridge or an AILANG effect-trace export) |
| R-SE-12 | **Forgeable limit line.** A program can print its own `policy-result:` line and exit 3, so a reported `limit` is not authenticated. Confinement is unaffected. | Upstream: the supervisor should report on a channel the program cannot write |
| R-SE-13 | **`ailang-cli` writes outside the deny list.** Measured on v0.51.0: `fmt --write` ignored `fs_deny_write`, and `test` wrote a transient `_namedtest_body_<random>.ail` into the tested file's directory, deny-listed ones included. World refuses `fmt --write` in the plan and the handler. On v0.52.1, `policy-tool` refuses `fmt --write` on a deny-listed path and `test` no longer writes there (V75). | World re-audits on each new tool binary (`TestAilangCLIAdmittedFlagsAreTheAuditedSet`) |
| R-SE-14 | **Replay has no production caller.** `Coordinator.Replay` and the replay binder are exercised only by tests; no daemon route or CLI verb runs a replay. | A later row adding a replay verb |
| R-SE-15 | **Case-sensitive deny list on a case-insensitive volume.** Measured on v0.51.0 and fixed in World by rendering every fold variant (198 patterns). On v0.52.1 the deny matcher folds case upstream (`#1559`, V72), so the folded list is now redundant and kept as defence in depth. | Fixed upstream in v0.52.1; World re-checks on each new tool binary |
| R-156 | **An agent can edit the tests that judge it.** `workspace-exec` runs the operator's profiled test command against whatever test files are in the worktree, and the agent can write those files. A floor or grading run that trusts `workspace-exec` output can be gamed. Ruled a blocker for `v1.0.0` (D-WORLD-81). | Queue row 156 (protected test paths), next |
| R-106-3b | **No live capsule-to-broker bridge.** Transitions plan one effect round and one finish round; no data-dependent multi-step effects inside one call. | A future row, upstream feature request if needed |

## Platform-wide

| Residual | Owner |
|---|---|
| **Unauthenticated reads.** Every `GET /v1/*` route and `/workbench` passes without a session (session-authority residual R1). Only `POST /v1/commit`, the agent card, `/a2a/` and `/mcp/` require one. | A follow-up queue row (a predicate change in the middleware) |
| **The commit route checks only that a session exists**, not its grants, and does not record the commit through the session's effect ledger. | Session-authority residual R3 (deferred to a later consumer row) |
| **`session new`, mint and revoke need the daemon stopped** because they open the store for writing (`session list` does not). A live mint-then-commit needs a restart. | Recorded usability gap (value-demonstration finding F-5) |
| **A content-hash mismatch on `POST /v1/commit` answers `500 internal store failure`**, not a 4xx. | Proposed row `w-commit-content-mismatch-is-a-client-error` (finding F-3) |
| **Indeterminate effect receipts.** A crash after the effect record write and before the outcome append leaves a record that needs deterministic reconciliation before the receipt resolves. Treat `indeterminate` as fail-closed; never retry automatically. | Effect journal Decision 5 (stated, not hidden) |
| **No memory or CPU limits, and no containers or VMs** beneath the semantic checks. Only `workspace-exec` runs inside a sandbox with no network (srt); the AILANG tools rely on the policy layer. | Future milestone (container or microVM isolation) |
| **Registry namespace ownership.** The public registry does not check the vendor prefix against the publishing key. | Upstream `sunholo-data/ailang#633` |

## Upstream asks from the software-engineering row

Filed in `sunholo-data/ailang` while landing row 134: `#1547`, `#1548`, `#1551`–`#1554`,
`#1557`–`#1559`. The charter records that most are fixed in AILANG v0.52.1. Three are named
in the decision ledger:

- `#1557` — admit `Declassify` in a policy. Shipped in v0.52.1; row 135 (stdin, argv and
  per-call capabilities for `ailang-run`) builds on it.
- `#1558` — port-scoped loopback in policies. Shipped in v0.52.1; a row-135 `Net` run reaches
  only the operator's `--run-net-allow` pairs.
- `#1559` — case-folded `fs_deny_write` matching (R-SE-15). Folded upstream on v0.52.1 (V72).

World does not work around language gaps locally. When upstream fixes land, World adopts them by
moving to a new measured tool binary, behind a re-run of the confinement matrix.
