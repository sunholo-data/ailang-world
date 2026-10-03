---
title: Effects, capabilities and budgets
sidebar_position: 3
description: How the effect broker decides, records and accounts for every brokered effect, and how grants, scopes and per-episode budgets work.
---

# Effects, capabilities and budgets

## The broker

The **effect broker** is World's in-process authority boundary. It has no REST route and no
CLI verb of its own: the daemon's invocation coordinator sends each planned effect to it. For
every request the broker:

1. **Decides** against the session's grants: effect name, scope, liveness (expiry) and budget.
2. **Records.** A denied request writes an immutable effect record with the denial label and
   returns; nothing is journaled as spend and nothing is debited. An allowed request appends a
   durable **effect intent**, debits the grant, runs the handler, then writes the **outcome** and
   the effect record.
3. **Returns** the result, or a typed denial or failure.

Denial labels are fixed strings:

| Label | Meaning |
|---|---|
| `denied:effect-name` | No grant for this effect name |
| `denied:scope` | A grant exists, but for a different scope |
| `denied:expired` | The grant is not live (now is outside `0 ≤ now < expiresAt`) |
| `denied:budget` | The grant's remaining budget is below the request's cost |

An effect record (`world/effect-record/v1`) carries `effect`, `scope`, `cost`,
`budgetBefore`, `budgetAfter`, `allowed`, `failed`, `denial`, `requestRef` and `resultRef`.
Requests (`world/effect-request/v1`) and results are stored as content-addressed objects as
well, so a record can be resolved all the way down. In replay mode the broker reads these
records and never calls a handler.

If a handler fails, the broker records `failed` and the budget stays debited. A handler that
cannot tell whether its effect happened records an **indeterminate** outcome; treat an
indeterminate receipt as fail-closed and never retry it automatically.

## Capabilities (grants)

A grant is `EFFECT=SCOPE:BUDGET`, for example `Workspace.Read=worktree:50`. It is minted into a
session at a terminal (see [Sessions and episodes](sessions-and-episodes.md)). Each grant
expires with its session.

A grant matches a request on **effect name, scope, liveness and `cost ≤ remaining budget`**.
Cost is not part of a grant's identity. Scope is compared by exact string equality. The
software-engineering tools use the symbolic scope `worktree`, which the daemon binds to the
episode's own worktree directory.

Two different triples are checked at two different points:

- **Visibility and admission**: a transition's `access` (`{effect, scope, cost}`) against the
  grants. A session sees exactly the transitions whose `access` its grants satisfy, on the A2A
  card and in MCP `tools/list`, and may call only those.
- **Execution**: each planned request must exactly equal one of the descriptor's
  `declaredEffects`, and then pass the broker's decision against the grants.

The eight `se-tools` transitions each declare `access` at cost 0 and one effect at cost 1, so
one grant per effect name covers both. Six grants cover all eight tools, because `ailang-write`
and `ailang-edit` share `Workspace.Write` and the two searches share `Ailang.Discover`. A grant
with budget 0 still shows the tool, but every call is recorded as `denied:budget`.

### Effect names in the code today

| Effect | Handler | Status |
|---|---|---|
| `Workspace.Read`, `Workspace.Write`, `Ailang.Check`, `Ailang.Run`, `Ailang.Discover`, `Ailang.CLI` | `AilangToolHandler` (the coding tools) | Served when `serve` has both `--workspace-root` and `--tool-ailang-bin` |
| `FS.Read`, `FS.Write` | Direct file handler, scope = one canonical absolute file path | Library only; no production caller |
| `Git.Commit`, `Model.Infer`, `Human.Approve`, `Human.PollApproval` | Broker handlers | Library only; not served by the daemon |
| `Registry.Publish` | Package publish handler | Used only by the attended `world-publish publish` |

When a transition declares an effect the daemon has no handler for, the coordinator refuses it
before the plan or any effect runs (rule R8). Over A2A the message is
`transition declares an effect this daemon has no handler for`.

## Budgets

Budgets are counted in **calls** for the coding tools (each tool call costs 1). They are not
wall time, tokens or money.

Budget **persists per episode**. Before an effectful dispatch the coordinator sums the cost of
every executed effect intent for the episode, grouped by effect and scope, and seeds each grant
with `budget − spent` (never below 0). So:

- Budget survives restarts of the daemon.
- Budget survives minting a new session for the **same episode**: a fresh credential does not
  reset spend. Use a new episode for a fresh budget.
- Denied requests are not spend.
- Same-episode calls are serialised by a per-episode lock, so two concurrent calls cannot both
  spend the last unit.

:::note Two layers of budget for `ailang-run`
One `ailang-run` call costs one `Ailang.Run` unit no matter what the program does inside. The
program's own file operations are bounded by the AILANG policy's `[budgets] FS=1000` and
`timeout_ms`, not by the broker, and are not individually journaled (residual R-SE-11 in
[Residuals](../security/residuals.md)).
:::

## Receipts

A commit sent with `"invocationId": "rest:<id>"` records its intent first, and
`GET /v1/receipts/rest:<id>` then answers `resolved` (with `resultRef`), `not-started` or
`indeterminate`. Only the `rest:` namespace is exposed; tool-call and effect journals are not
reachable through this route. See the [HTTP API](../reference/http-api.md#get-v1receiptsid).
