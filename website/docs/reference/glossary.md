---
title: Glossary
sidebar_position: 5
description: Short definitions of the terms used across the AILANG World documentation.
---

# Glossary

**Access.** The `{effect, scope, cost}` requirement on a transition descriptor. A session sees
and may call a transition only if its grants satisfy it.

**AILANG.** The deterministic, effect-typed programming language World uses as its transaction
language. Lives in `sunholo-data/ailang`; World consumes released binaries only.

**Archive.** The content-addressed store of binaries next to the world store. `--ailang-bin`
and `--tool-ailang-bin` are archived and re-hashed before use.

**Attended.** An operation that requires a human at a controlling terminal: `session mint`,
`world-publish approve`, `publish`, `transitions`.

**Broker (effect broker).** The in-process authority boundary that decides, journals, executes
and records every brokered effect.

**Budget.** The amount a grant may spend, in the effect's cost units (calls for the coding
tools). Persisted per episode.

**Calling convention.** Every transition exports `main(input: string) -> string`. Effectful
transitions are called in a plan phase and, optionally, a finish phase.

**Capsule.** The subprocess floor a transition runs in: pinned interpreter, `--caps ""`, FS jail,
scrubbed environment, timeout, output cap.

**Capability.** See *grant*.

**Clause.** One of the seven checkable end-states of the World 1.0 bar in the mission charter.

**Commit.** A compare-and-append that adds objects, a world row and one log entry, accepted only
if the observed head is still the selected head.

**Coordinator (invocation coordinator).** The host component that runs a transition call: bind,
plan, effect, finish, commit.

**Credential ID.** The hash of a session token that the store keeps; used to revoke.

**Declared effects.** The exact `{effect, scope, cost}` triples a transition's plan may request.

**Effect.** A named kind of interaction with the outside world, such as `Workspace.Read` or
`Ailang.Run`. The only nondeterminism in World.

**Effect plan.** The `world/effect-plan/v1` object a pure transition returns to request effects.
At most one effect per plan in v1.

**Effect record.** The immutable `world/effect-record/v1` object written for every allowed,
denied or failed brokered effect.

**EffectsUnrecorded.** The error returned when an effect ran and was recorded but the call did
not commit. Names the effect record refs over A2A.

**Epoch (semantics epoch).** A number the epoch registry assigns to an interpreter release; log
entries and descriptors carry it so history is replayed under the right semantics.

**Episode.** The unit of work a session belongs to. Names the budget ledger and, for the coding
tools, the worktree `<workspace-root>/<episode>`.

**Genesis.** The first commit of a store, with an empty `observedHead`.

**Grant.** `EFFECT=SCOPE:BUDGET` minted into a session, with an expiry. Matches a request on effect
name, exact scope, liveness and `cost ≤ budget`.

**Head (selected head).** The current world ref.

**Indeterminate.** A receipt state meaning the outcome is unknown. Fail-closed: never retry
automatically.

**Interpreter pin.** The archived interpreter's content hash, written into every log entry and
descriptor. Reported as `interpreter_ref` by `/v1/health`.

**Invocation record.** The `world/invocation-record/v1` (pure) or `/v2` (effectful) object a log
entry's `transitionRef` points at.

**Journal.** Durable intents and outcomes for commits and effects; the source of receipts and
budget spend.

**Log.** The append-only sequence of transition entries, each with a six-field header.

**Object.** A content-addressed blob with `hash`, `interfaceHash`, `semanticId` and `provenance`.
Never modified.

**PIN.** The interpreter release that runs transitions: AILANG v0.41.0.

**Policy (episode policy).** The restricted AILANG policy rendered per episode at
`<db-dir>/policies/<episode>.toml`, which confines the coding tools.

**`policy-tool`.** The AILANG binary's typed tool interface, run by the handlers for every coding
tool except `ailang-run`.

**Provenance walk.** Answering "why did this happen" by following content-addressed references
through the read API.

**Receipt.** The three-state outcome of a `rest:` commit: `resolved`, `not-started`,
`indeterminate`.

**Registry (transition registry).** The published, pinned set of transitions at
`world/transition-registry/v1`.

**Replay.** Re-deriving a recorded invocation from its pure transition plus recorded effect
results, without calling any handler. Implemented in the library; no operator verb yet.

**Residual.** A known, measured gap with a named owner.

**Scope.** The string an effect is authorized over. Compared exactly. The coding tools use
`worktree`.

**Self-contained module.** A transition source that imports only `std/*`, so it type-checks in an
empty scratch root.

**Semantic ID.** An object's type-like name, for example `world/effect-record/v1`. Queryable with
`GET /v1/objects/by-semantic-id/…`.

**Session.** A 64-hex bearer credential bound to one episode, a set of grants and an expiry.

**SSE.** Server-Sent Events. `/mcp/` replies as `event: message` with the JSON on a `data:` line.

**STOP.** `world-publish`'s refusal: `STOP fence=<name>`, exit 3, nothing happened.

**TOOL.** The AILANG release the coding tools run: v0.51.0.

**Transition.** A pure, pinned AILANG program and the only way World's state changes.

**TTY fence.** The requirement that a process have a controlling terminal, and for
`world-publish` that stdin be that terminal.

**Workbench.** The read-only HTML operator renderer at `/workbench`.

**World.** An immutable state value: `{ref, revision, stateRoot, logHead}`.

**Worktree.** The git worktree an episode's coding tools are confined to.
