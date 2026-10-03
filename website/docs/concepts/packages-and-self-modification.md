---
title: Packages and self-modification
sidebar_position: 6
description: How World ships its own behaviour as AILANG packages, what the world/core and se-tools packages are, and the attended fences around publishing.
---

# Packages and self-modification

World changes itself through the same machinery it offers everyone else. Its behaviour lives
in versioned **AILANG packages**, and the kernel stays deliberately small. Every proposed
addition to the kernel or host has to answer "why is this not a package?" (coding standard
S3).

## Two kinds of publishing

There are two different publish operations. Do not confuse them.

| | Transition registry publish | Package publish |
|---|---|---|
| Command | `world-publish transitions` | `world-publish approve`, then `world-publish publish --live` |
| Writes to | The local store's `world/transition-registry/v1` head | The public AILANG package registry |
| Reversible | A later revision can replace a descriptor; history stays | **No.** A published version is immutable |
| Fences | Daemon stopped, controlling terminal, typed confirmation | Daemon not involved; controlling terminal, typed confirmation, a one-shot approval minted in a separate invocation, a golden ready-packet, a credential file |
| Runbook | [QUICKSTART §6 and §9](../getting-started/coding-tools.md) | `docs/SELF_MOD_PUBLISH.md` |

Both are attended: an agent cannot complete either. See
[Attended steps](../guides/attended-steps.md).

## `world/core`

`packages/world-core/` is World's pure semantic core as an AILANG package: the modules
`world/types`, `world/contracts`, `world/transitions` and `world/logepoch`. It declares no
effects (`[effects] max = []`). It is a deterministic projection of the canonical sources,
rebuilt with `./scripts/build_world_package.sh`, never edited by hand.

- `world/core@0.1.0` was published on 2026-09-21. That publish was clause 7's "controlled
  self-modification proven once".
- `world/core@0.1.1` was published on 2026-10-02 (attended). It is proof hardening with no
  exported signature change: 16 contracts verified, 0 refuted.

Published versions are immutable. The next release needs a new candidate, a new golden
ready-packet, and a new attended run.

## `se-tools`

`packages/se-tools/` holds the eight software-engineering transitions and their checked-in
publish manifest, `packages/se-tools/transitions.json`. Each module is self-contained
(`module se_tools/read` and so on, importing only `std/*`), exports
`main(input: string) -> string`, and carries a contracted path predicate plus named inline
tests. It is published into a store's transition registry with `world-publish transitions`, not
to the public package registry.

## The boundaries

Self-modification is scoped. These are out of scope and are never changed through World's own
pipeline:

1. **The mission-loop machinery** (the shared driver and the mission-control skill).
2. **The AILANG compiler.** Language gaps found here are filed upstream in
   `sunholo-data/ailang`; World never forks or works around the compiler.
3. **The live daemon binary.** Replacing a running `ailang-worldd` is an attended operation.

Because worlds are immutable, rolling back a bad change is a roll forward: a new transition
that re-pins the earlier version, with the failed episode kept as evidence.

:::note What is not enforced
The `world/` package namespace is a convention World enforces in its own gate, not registry
namespace ownership. The registry currently authorizes publishers with a shared key that is not
checked against the vendor prefix (tracked upstream as `sunholo-data/ailang#633`).
:::
