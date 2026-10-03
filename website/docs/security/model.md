---
title: Security model
sidebar_position: 1
description: Explicit authority in World — what the session, broker, capability and budget checks enforce, what AILANG's policy layer enforces for the coding tools, and where each check lives.
---

# Security model

World's security principle is **explicit authority**: no computation holds ambient authority,
and every effect needs a capability, a budget and a recorded decision. Security has two layers
and neither substitutes for the other:

- **Semantic**: sessions, capabilities, budgets, declared effects, recorded decisions.
- **Physical**: process confinement beneath the semantic checks.

This page says what is enforced where, today.

## The path of one tool call

```text
agent ── Authorization: Bearer ──▶ session resolver ── (episode, grants)
                                       │
            registry snapshot ∩ grants ▼
                               visible transitions ── tools/list, agent card
                                       │ call by ID
                                       ▼
                     pure transition (capsule, --caps "") ── effect plan
                                       │ plan laws; (effect,scope,cost) ∈ declaredEffects
                                       ▼
                     broker: effect name, scope, liveness, budget ── record
                                       │ allowed
                                       ▼
                     handler: policy-tool / run --policy in the episode worktree
                                       │ AILANG policy layer confines the subprocess
                                       ▼
                     commit one log entry (per-episode lock, compare-and-append)
```

## What is enforced, and by what

| Check | Enforced by | Where |
|---|---|---|
| Daemon reachable only locally | Loopback-only bind, no override | `serve` startup |
| No ambient publish credential | Refuse to start if `AILANG_REGISTRY_API_KEY` is set; tool children get a minimal environment | `ailang-worldd` entry, tool handler |
| One writer per store | Non-waiting OS lock on `<db>.writer.lock` | store `Open` |
| Who is calling | Bearer session credential, stored hashed, with expiry | session middleware and projection |
| Minting authority needs a human | `/dev/tty` and a confirmation | `session mint` |
| What a session can see and call | Registry descriptors filtered by `access` against live grants, fresh per request | MCP `tools/list`, A2A card and admission |
| Transitions cannot do I/O themselves | Transitions run in the capsule with `--caps ""`, a pinned and re-hashed interpreter, an FS jail, a scrubbed environment, a timeout and an output cap | capsule |
| A transition asks only for what it declared | Plan laws, and each request must equal a `declaredEffects` triple | coordinator, then the bound invoker again |
| Each effect is authorized and paid for | Effect name, exact scope, liveness, `cost ≤ budget`; denials recorded without debit | broker (decision law Z3-verified in a sketch, transcribed to Go under a drift test) |
| Budget cannot be reset by reconnecting | Spend summed per episode from durable intents | coordinator |
| Tool effects stay inside the worktree | AILANG's policy layer (`policy-tool`, `run --policy`) under a rendered restricted policy | the v0.52.1 tool binary |
| Agents cannot edit their own policy, cache or corpus | Startup refuses state or corpus inside the workspace root; policy and cache live under `<db-dir>` | daemon startup |
| Every outcome is evidence | Effect records, intents and outcomes; one log entry per call | broker, journal, store |
| Irreversible and registry writes need a human | Controlling-terminal fence, typed phrase, one-shot approval | `world-publish` |

## The AILANG policy layer does the confinement

World adds **no Go-side path confinement** for the coding tools. The handlers delegate to
AILANG's hardened policy layer, which was measured first-party on the exact tool binary:

- `policy-tool` serves read, write, edit, check, ai_check, the search ops and an allowlisted set
  of CLI ops, and refuses `run`.
- `ailang run --policy` runs programs under a supervised worker whose own file layer applies the
  same sandbox and protections to everything the program does.

Each episode gets a rendered policy at `<db-dir>/policies/<episode>.toml`:
`security_mode = "restricted"`, `allowed_caps = ["IO", "FS"]`, `fs_sandbox` = the episode
worktree, `timeout_ms = 8000`, an `fs_deny_write` list, `entry = "main"`, and
`[budgets] FS = 1000`. The full list of measured confinement facts is in
[Tool confinement](tool-confinement.md).

This is a deliberate trade, ratified by the owner: World inherits the AILANG lane's hardening
rather than reimplementing it, and it records the gaps that hardening still has as residuals.

## What is not provided

Stated plainly:

- **No network isolation, no memory or CPU limits, no containers.** The capsule and the tool
  handlers are a process-safety floor (pinned executable, FS jail, scrubbed environment,
  timeout, output cap), not full isolation. Docker is explicitly out of scope for the coding
  tools. Container or microVM isolation is future milestone work.
- **Reads are unauthenticated.** Every `GET /v1/*` route and `/workbench` is open to anything
  that can reach the loopback port.
- **The inner effects of `ailang-run` are not brokered individually.** One run is one
  `Ailang.Run` effect. What the program does inside is confined and budgeted by the AILANG
  policy, recorded only as the run envelope, and not reconstructible in replay (R-SE-11).
- **A run's reported `limit` is not authenticated.** A program can print its own
  `policy-result:` line (R-SE-12). This misleads only the agent about its own run.
- **Calls are at-least-once under retry.** MCP JSON-RPC IDs give no idempotency; a retried
  call creates a new task and may repeat an effect. Every execution is recorded (R-SE-6).
- **Namespace ownership in the public registry is not enforced** by the registry itself
  (`sunholo-data/ailang#633`).

See [Residuals](residuals.md) for the full list with owners.

:::warning A confined program can never be the World client
Restricted mode has no loopback grant and refuses `--net-*` flags under `--policy` (measured: a
run's `httpGet("http://127.0.0.1:7644/")` printed `E_NET_IP_BLOCKED`). The client of World's MCP
surface is always an outside harness such as pi or Claude Code.
:::
