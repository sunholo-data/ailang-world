# w-decision-packet-production-projection — the coordinator mints `world/decision-packet/v1` on the effectful commit that requests `Human.Approve`, and the workbench object page validates and projects it

- **Status**: **planned** (designed iteration 255, designer role, unattended mission stage; quorum round 1 BLOCKED 3/3 on completeness/verification only → revision §0.5; round 2 BLOCKED 3/3 → controller-authorised narrow-refinement carve-out applying the reviewers' own fixes, §0.6)
- **Queue row**: 172, `w-decision-packet-production-projection` · release prerequisite (moves no clause) · on the attended `D-WORLD-69` = **A** path 26 → 105 → **172** → 150 (rows 26 `0390aa5` and 105 `9d16a07` have landed). Owns the packet residual of row 150.
- **Measurement base**: `e869bfe4ec4aa0523bbe059494b7bf62d762a6e9` (V1); every V-row below was run first-party in this worktree.
- **Instrument**: `$HOME/.pinned-ailang/ailang`, AILANG **v0.41.0**, sha256 `1a67b0146858450182f48082299956ea5b08cdf30131979b188b339fcedb5b9f` (V24). The one `.ail` this design proposes was checked, Z3-verified, inline-tested and mutation-probed with that binary; the checked artifact is banked as `design_docs/verification/world-row172-design/decisionpacket.ail.txt` with its log `sketch-check.log` (`.txt` so the exact Leg-1 module manifest is not perturbed before M1 lands it, V23).
- **Estimate**: **~3.3 d**, 4 milestones, ~410 production LOC (Timeline). One sprint; the rest of the packet lifecycle is named follow-up work (Future Work).
- **Scope**: `design_docs/sketches/` (one new sketch), `scripts/verify_ail.sh` (two manifest lines), `host/coordinator` (codec, law mirrors, producer hook, validator), `host/broker` (three small pure exports), `host/daemon` (Human.Approve binding, REST reservation, object-page projection), `host/workbench/render.go` (one view + template block), `docs/QUICKSTART.md` + its executed test. **No** `world/*.ail` change (the S8 floor does not move), **no** store schema change, **no** new route or query key.

## §0 The row's premises, re-measured at base (do not trust, verify)

### 0.1 "No decision-packet producer or read projection" — TRUE

`DecisionPacket|decision-packet|decision_packet` reads **0** non-test Go files under `host`/`cmd`; the control `world/approval-request/v1` reads **1** (`host/broker/approve.go`) in the same call (V2). `world/types.ail:154-167` declares the seven-field record under semantic ID `world/decision-packet/v1`, logical times only (V3). No byte encoding was ever frozen: item 15 froze the record and its laws, not a wire form (V3, V28) — so this design defines the v1 wire bytes, as the first producer must. (The row-150 review's quoted range is `independent-review.md:17`, V28.)

### 0.2 "Linking a stored proposal" — PARTIALLY FALSE: no host path stores a `Proposal`

`world/types.ail` declares `Proposal`, but **0** Go or `.ail` files under `host cmd world packages` name a `world/proposal*` semantic ID, and **0** Go files name any `proposal/v*` ID, against a firing control (`world/effect-plan/v1`: 1 non-test file; `approval-request/v`: 2 matches) (V4). The only stored, committed object that IS a proposal in the propose → verify → commit sense is the effectful invocation's **plan object** (`world/effect-plan/v1`): the pure plan phase's proposed effect, hashed BEFORE any effect runs (`effectful.go:293`) and committed with the v2 record (`plan.go:121,132`, V8). This design binds `proposalHash` to that object (D1) and does not invent a `world/proposal/v1` — a new kernel wire type with a `list[Evidence]` and a contested `confidence` float (HUMAN-SURFACE §7 point 5) is not a sprint-sized side effect of this row.

### 0.3 "The normal commit path produces approval requests" — FALSE for the daemon today

`world/approval-request/v1` objects are minted only by `HumanHandler.Execute`'s `Human.Approve` branch (`approve.go:103-121`). Its only production construction is `mintAttendedApproval` (`publish_op.go:131-138`), reached only from the attended, TTY-fenced `cmd/world-publish` (`main.go:270`) (V5). The daemon's production binder registers only the workspace family: `d.binder` passes `d.workspace.registry(episodeID)` (`daemon.go:780-782`), and `Human` reads **0** in non-test `host/daemon` Go against a firing `EffectWorkspaceExec` control of 2 (V6). So a coordinator-side producer would have **zero** production reach unless the daemon binds `Human.Approve`; this design makes that binding an explicit, analysed deliverable (D4), not an assumption.

### 0.4 Citation drift from the row-150 review, corrected

The row-150 review cites the daemon coordinator's wall-clock `Now` at `daemon.go:694`; at base it is **`daemon.go:711`** (`Now: func() int64 { return time.Now().Unix() }`), while `coordinator.Config.Now` is the injected seam (`coordinator.go:52`, required at `:78`) (V7). This design reads no clock: the packet's times are the coordinator's one logical `now` (D3). Whether the daemon's tick source stays Unix seconds is **not** this row's decision — the review's H2 says a new production clock knob affecting invocation/expiry authority warrants its own freeze (V28); row 150 owns that question.

### 0.5 Quorum round 1 → revision

Round 1 was BLOCKED 3/3 (`gemini-3-1-pro`, `oc-glm-5-3`, `oc-kimi-k3`); no objection disputed the direction. The common gap: the validator's anchor hops had no §10 rows. Every fix below was re-measured first-party at `e869bfe` (V35–V43); no measurement contradicted the controller's figures.

- **Anchor hops unverified (all three):** `LogEntry` shape → V35; the record-v2 field names (`TransitionFn`, `Plan`, `Effects`; unexported, so decoded in-package) → V36, used verbatim in D1/D5/D6; the readStore method set ⊇ PacketReader → V37 (now cited in D5 in place of V14/V17).
- **`oc-glm-5-3` — would an extra commit object break replay/reconcile?** No: V38 shows that both paths read the record only. Added AC3b `TestReplayReconcileToleratesPacketObject` and mutation MP8; D1 is unchanged.
- **Determinism of code 7's scan:** V39 (`ORDER BY hash_ref`); D5 now says "the first 100 by `hash_ref` order".
- **M4's route:** V40.
- **`oc-kimi-k3` — the REST-forgery residual stated conditionally:** rewritten in the measured direction in D5, from V41 (the REST entry's `transitionRef` is caller-supplied).
- **`oc-glm-5-3` — D2 is an irreversible freeze:** now flagged next to D4, change cost **high**, with V42 for the item-15 content-addressed-wire-type discipline.
- **D4 (controller):** the authority analysis is now stated in Risks; the flag is kept.

### 0.6 Quorum round 2 → narrow-refinement carve-out

Round 2 was BLOCKED 3/3. No objection disputed the direction, and each carried a concrete fix. This is the controller-authorised, bounded second revision: it applies those reviewers' own fixes and nothing beyond them. Each premise was re-measured first-party at `955fbc0`, whose code equals `e869bfe` (V43).

- **`gemini-3-1-pro` — D4's authority analysis cited `validatePublishApproval` without a row.** Added V44, which quotes its decision and scope enforcement; D4 and the Risks row now cite it.
- **`oc-glm-5-3` — the TR.C claim was asserted, not verified.** Added V45, which quotes the census mechanism and its invariant, and rewrote the §9 TR.C cell. Per the controller's ruling, the pre-existing blindness (the daemon's `Workspace.Exec` binding has the same shape) is Future Work #5, not absorbed here.
- **`oc-glm-5-3` — the M4 quickstart precedent was unrowed.** Added V46.
- **`oc-glm-5-3` — the ≤ 80-read figure was unsubstantiated.** Added V47: 80 object reads plus one registry-head read. A5 and D7 now say "≤ 80 object reads + 1 head read".
- **`oc-kimi-k3` — content-identical requests collapse, so a requestRef-only code 7 refuses genuine packets.** The premise is real (V10; the approvals head appends per call, V20). Code 7 is now scoped to the `(requestRef, proposalHash)` pair: the new sketch law `samePair` and the renamed fact `pairUniqueOk` were re-checked (V24: 9/9 Z3-VERIFIED, `Passed: 46`, three mutants `VIOLATION`). Also added: the D5 subsection "Request-identity collapse interacts with code 7"; AC5 arm 7b; the two-approvals arm in AC1/AC3; MP8b; a Risks row on non-monotonicity; and the new code-7 sentence.

## Axiom Compliance

| Axiom | Score | Justification |
|-------|-------|---------------|
| A1: Determinism | +1 | The packet is a pure function of the commit's facts (logical `now`, the authorizing grant's expiry, the plan ref, the request ref); the validator's verdict is a pure ordered law (`packetRefusal`) over host-measured facts, so one store state renders one page |
| A2: Replayability | 0 | Replay is untouched; the packet is a committed object outside the recomposed output. Declared residual: replay does not re-derive the packet (its deadline is the grant's expiry, which effect records do not carry) |
| A3: Effect Legibility | +1 | No new effect; an existing landed handler is bound and serialized at the daemon boundary; all laws are pure `! {}` |
| A4: Explicit Authority | +1 | A packet confers no authority (it fronts a request; decisions stay `DecideApproval`'s); `Human.Approve` still needs a declared effect AND a live grant (broker law unchanged); REST may no longer write packet objects |
| A5: Bounded Verification | +1 | Validation is ≤ 2 indexed commit lookups, ≤ 6 point reads, one `RecentApprovals` call (≤ 40 heads; ≤ 80 object reads + 1 registry-head read, V47), one 101-row packet page with ≤ 100 point reads, and two row-105 grade views — all under the request's `readDeadline`; every bound is disclosed on the page |
| A6: Safe Concurrency | +1 | The daemon serializes `Human.Approve` in-process, closing the approvals-head lost-update race the binding would otherwise expose (V20) |
| A7: Machines First | +1 | Refusals are the codes 1–7 of a Z3-verified law, mapped to frozen sentences |
| A8: Minimal Syntax | 0 | No language surface |
| A9: Cost Visibility | +1 | The non-packet object page adds zero reads; the packet page's worst case is stated (D7) |
| A10: Composability | +2 | Reuses the plan object, the v2 record, `broker.DecodeRecord`, `RecentApprovals`, `ObjectCommits`, `ObjectsBySemanticID`, row 105's `objectGradeView`, and the effectplan sketch/drift pattern |
| A11: Structured Failure | +1 | Every refusal arm has a forged-input test and a named mutation (§8) |
| A12: System Boundary | +1 | Laws in the AILANG sketch (Z3), effects and bytes in Go, kernel `world/` untouched |

**Net +12 → move forward.** Hard axioms A1/A3/A4/A7 all ≥ 0.

## Problem Statement

**Current state (measured at `e869bfe`):** the packet type exists only as a declaration (V2, V3); no stored object is a `Proposal` (V4); the daemon cannot mint an approval request (V6); the workbench's decisions pane renders the approvals chain (`workbench.go:582`) and must not be renamed into packets (row-149 V15/V16, V17). A credentialed REST commit can store an object under ANY semantic ID — `decodeCommit` copies `object.SemanticID` verbatim (`handlers.go:847`) and reserves none (V15) — so "an object labelled `world/decision-packet/v1`" proves nothing. What REST cannot do is move a registry head: `SetRegistryHead` has exactly two production callers, the broker's approvals append and the package registry, and `commitRequest` carries no registry field (V16).

**Impact:** row 150 cannot show a real packet, item 7 has no packet to enforce timeouts on, and any packet a client fabricates today would be indistinguishable from a produced one.

## Goals

**Primary goal:** every committed effectful invocation whose single planned effect is an allowed, successful `Human.Approve` commits exactly one `world/decision-packet/v1` object in the SAME commit as its plan and record; the workbench object page of that packet renders it as **validated** only when its bytes, carrying commit, proposal, request, approvals-chain anchor, lifecycle fields and uniqueness all check — and otherwise renders the first failing check's sentence.

**Success metrics:** (1) AC2: a daemon invocation through the normal coordinator path produces a packet whose fields equal the request's logical `now`, the authorizing grant's expiry, the plan ref and the request ref. (2) AC5: seven forged-input arms each render their own refusal. (3) AC6: evidence on a validated packet comes only through row 105's `objectGradeView`. (4) AC7: REST cannot store a packet object.

## High-Impact Decisions

| Decision | Why High Impact | Chosen By | Change Cost |
|----------|-----------------|-----------|-------------|
| **D1 Producer = the coordinator's effectful commit**, gated on the plan's one effect being `Human.Approve` with status OK; the packet rides in `Commit.Objects` | The only path where a stored proposal (plan) and a broker-minted request coexist under one logical `now` (V8, V10); atomic with the record, so a packet exists iff its invocation committed | agent (V4, V8) | low |
| **D2a `proposalHash` = the plan object ref** | No `Proposal` object exists (V4); the plan is the pure, pre-effect proposal (V8) | agent | low (a `/v2` packet may later point at a real `Proposal`) |
| **D3 Lifecycle fields: `createdAt` = the invocation's `now`; `deadlineAt` = the authorizing `Human.Approve` grant's `ExpiresAt`; `policy` = `Cancel`; `escalationsRemaining` = 0** | No default is invented: the deadline is the effect's existing strict TTL (HUMAN-SURFACE:161-162, V25); `Cancel` is the one policy whose Z3-proven outcome is the ratified "typed rejected-timeout" unconditionally (`types.ail:174`, V25); a live grant makes the schedule well-formed by a proven law (`mintedScheduleWellFormed`, V24) | agent, citing ratified text | low |
| **D2b The v1 canonical byte form and `packetHash` rule (D2 below) — an IRREVERSIBLE freeze** | The first stored packet fixes these bytes forever: content addresses and `packetHash` depend on them, and amendments are `/v2`, never in-place (item-15 discipline, `w-decision-lifecycle-freeze.md:149-151`, V42) | agent — **flagged for controller/quorum confirmation**, like D4 | **high** (irreversible once a packet is stored) |
| **D4 The daemon binds `Human.Approve` (only), serialized by a daemon-held mutex; REST refuses `world/decision-packet/v1` objects** | Without the binding the producer has no production reach (V6); without serialization the approvals head can lose an append (V20); without the reservation REST can forge the label (V15) | agent — **flagged for controller/quorum confirmation** (an effect-surface widening; see Risks) | low |
| **D5 Validation is relational, anchored in two writer-exclusive facts**: the packet's first carrying commit must be its record's first carrying commit, and its request must be in the broker-only approvals chain | REST can forge objects and commits but cannot re-be the first carrier of an existing record nor move the approvals head (V14, V16) | agent | med |
| **D6 Projection on the existing object page** (no new route or query key), evidence via row 105 | Row 105's grade path is the only PROVEN path (V18); no new key keeps the page grammar unchanged | agent | low |

**Design freeze:** D1–D6 rest on measurements in §10; two decisions are flagged for explicit confirmation: D2b's byte-form freeze (irreversible) and D4's effect-surface widening (the controller may route it as a `D-WORLD` ask; the default if unanswered is to land M1 and M3's validator and hold M2's daemon binding).

## Solution Design

### D1 — Who mints the packet, and when

In `dispatchEffectful` (`effectful.go:273`), after the single planned effect has run (L1: at most one effect per v1 plan, `effectplan.go:32`, V9) and BEFORE `planEffectInvocation` (`:333`):

```text
if len(pl.Effects) == 1 && pl.Effects[0].Requirement.Effect == broker.EffectHumanApprove && results[0].Status == StatusOK:
    req      := broker.EffectRequest{Effect, Scope, Cost of the requirement, Now: now}    # the exact request the session decided
    grant,ok := broker.AuthorizingGrant(grants, req)        # decideOver's choice over the seeded grants (V13); the session copied them (V32)
    if !ok: return unrecorded(errPacketGrant)               # structurally unreachable: same law, same inputs (no mutation claimed)
    request, err := broker.ParsePendingApproval(results[0].Output)   # the handler's pending wire (approve.go:121,277-278, V10)
    if err != nil: return unrecorded(err)
    obj, err := mintPacket(now, grant.ExpiresAt, planRef, request, writtenBy(surface))
    if err != nil: return unrecorded(err)
    p = withPacket(p, obj)                                   # appended to Objects and Commit.Objects
```

Denied or failed `Human.Approve` effects mint no request and therefore no packet. The packet's carrying commit is the invocation's own `store.Commit`, so `commit_objects` records it at that entry (`store.go:1196`, V14). `dispatchEffectful` gains the `grants` parameter it is called beside (`coordinator.go:311-334`, V32). Post-effect failures go through the existing `unrecorded` path, so an executed `Human.Approve` is never silently orphaned.

### D2 — Canonical bytes and `packetHash`

Payload: compact JSON, members in the frozen field order, refs as canonical `HashRef` text (`hashref.String`/`Parse`, V31), ints as JSON integers, `policy` as the constructor name:

```json
{"packetHash":"sha256:…","proposalHash":"sha256:…","requestRef":"sha256:…","createdAt":1700,"deadlineAt":1760,"escalationsRemaining":0,"policy":"Cancel"}
```

`packetHash` is the self-hash the type comment cites (precedent `Proposal.proposalHash`): `SumSHA256` of the canonical **body** — the same object minus `packetHash`, same member order. The stored object's content address is `SumSHA256(payload)`; its interface hash is `SumSHA256("world/decision-packet/v1")` and provenance `coordinator:<surface>`, the coordinator's `object()` convention (V8). Decoding is strict: `DisallowUnknownFields`, one value, re-encode byte-equal, both refs parse, `policyCode(policy) >= 0`, `packetHash` recomputes. A semantic-ID constant `DecisionPacketV1` joins `plan.go`'s constants. **This byte form is an irreversible freeze (D2b):** once one packet is stored, any change is a `world/decision-packet/v2`, which is the discipline item 15 fixed for content-addressed wire types (`w-decision-lifecycle-freeze.md:149-151`, V42). The bytes above are therefore part of the ratification surface of this design, not an implementation detail.

### D3 — Where each lifecycle field comes from (no invented default, no clock)

| Field | Source | Why it is not invented |
|---|---|---|
| `createdAt` | the invocation's one logical `now` (`effectful.go:284`), which is also the request's `now` (`approve.go:112`) and the intent's `LogicalTime` | the item-15 logical-time idiom; the packet code reads no clock (V7, V10) |
| `deadlineAt` | `ExpiresAt` of the grant the broker law selects for this exact request | HUMAN-SURFACE:161-162: "every Decision Packet and `Human.Approve` effect MUST carry a strict TTL"; the effect's only landed TTL is its grant expiry, and `capabilityLive` (`decide.go:101`) admitted the request only if `0 <= now < ExpiresAt` — so `deadlineAt > createdAt` holds by the Z3-proven `mintedScheduleWellFormed` (V24, V25) |
| `policy` | `Cancel` | HUMAN-SURFACE:162-163: on expiry "the packet deterministically resolves to a typed rejected-timeout"; `timeoutOutcome(Cancel, _, _) = ResolveRejectedTimeout` is Z3-proven (`types.ail:174`, V25). `EscalateBounded`/`ExecuteIfGranted` need an escalation budget or an independent-authority fact no landed surface supplies |
| `escalationsRemaining` | 0 | under `Cancel` the budget is never read; 0 is the only value `producerFieldsOk` admits |

### D4 — Production reach, serialization, and the REST reservation

- **Binding.** `d.binder` becomes: `reg := d.workspace.registry(episodeID); reg[broker.EffectHumanApprove] = d.human; return broker.OpenBinder(...)`. `registry()` returns a fresh map per call (`workspace.go:338-362`, V6), and its exact-key tests call `registry()` directly, never `d.binder` (0 test calls of `d.binder`, V21), so they stay green. `Human.PollApproval` is NOT bound (minimal surface). `DecideApproval` stays operator-only; the daemon gains no decision route.
- **Authority analysis.** Binding does not grant: a skill must declare `Human.Approve` and the caller must hold a live, budgeted `Human.Approve` grant with the exact scope (`Decide`, `decide.go:48`). A pending request authorizes nothing; publish consumption requires an operator decision AND a publish-shaped scope (`validatePublishApproval`, `approve.go:486`, V44: it refuses unless the decision is `approve` (`:513`) and the request is `Human.Approve` (`:547`) whose scope parses as a publish-approval scope (`:558`) authorizing `Registry.Publish` (`:562`) for exactly this package (`:572`) and stamping its hashes (`:584`), and it refuses expired approvals (`:613`)).
- **Serialization.** `appendApprovalHead` reads the head (`approve.go:196`) then sets it (`:211`) without compare-and-set (V20). The writer lock makes the daemon the store's only writer process (`writer_lock.go`, V20), but the coordinator runs different episodes concurrently (the lock is per episode, `effectful.go:91-93`). `d.human` is therefore `serializedHandler{mu: &d.humanMu, h: broker.NewHumanHandler(d.store)}`. The broker-side CAS fix is a candidate row, not absorbed here.
- **REST reservation.** `handleCommit` refuses (400, `BadRequest`, `objects[i].semanticId: world/decision-packet/v1 is minted only by the invocation coordinator`) any commit object with that semantic ID. Only this ID is reserved here; the general passthrough (V15) is a candidate row.

### D5 — The validator (`coordinator.ValidatePacket`)

`ValidatePacket(ctx, r PacketReader, obj store.Object) (PacketResult, error)`, where `PacketReader` is the subset `{GetObject, GetLogEntry, ObjectCommits, ObjectsBySemanticID, GetRegistryHead}` of the daemon's 11-method `readStore` (`daemon.go:501-513`, V37), so the daemon passes `d.reads` unchanged. `GetLogEntry` returns `store.LogEntry{Header, EntryHash, TransitionRef}` (`store.go:136`, V35). The record is decoded into the package's own unexported `recordV2` (`plan.go:39-50`, V36): the validator lives in package `coordinator`, so it decodes the record directly and needs no new export. Each step measures one boolean, then `packetRefusal` (the Go mirror of the Z3-verified law) picks the first failure. Store errors return `err` and never become a code.

| Code | Fact (measured in this order; a step whose prerequisite failed is false) | Reads |
|---|---|---|
| 1 bytes | D2's strict decode; interface hash and `packetHash` recompute | 0 |
| 2 carried | `ObjectCommits(packet, 0, 1)` → first entry E; `GetLogEntry(E).TransitionRef` = R and R is a `world/invocation-record/v2` that decodes into `recordV2`; `ObjectCommits(R, 0, 1)` first == E | 2 lookups + entry + R |
| 3 proposal | `R.Plan == proposalHash.String()` and that object is `world/effect-plan/v1` | 1 |
| 4 request | `R.Effects` has exactly one ref; `broker.DecodeRecord` (`record.go:58`) → `Human.Approve`, allowed, not failed; its result (`world/effect-result/v1`, the handler output verbatim, `broker.go:311,319`, V11) parses to `requestRef`; the request object is `world/approval-request/v1` with effect `Human.Approve` and the record's scope | 3 |
| 5 anchored | `requestRef` is among `RecentApprovals(ctx, r, 40)` summaries (V17); a malformed chain is false with its ref named | ≤ 80 object + 1 head (V47) |
| 6 lifecycle | `lifecycleOk(createdAt, request.now, deadlineAt, escalationsRemaining, policyCode)` | 0 |
| 7 unique | no OTHER object among `ObjectsBySemanticID(DecisionPacketV1, "", 101)` — the first 100 **by `hash_ref` order** (`ORDER BY hash_ref`, `objects_by_semantic_id.go:19-20`, V39), so the scan is deterministic — decodes to the same `(requestRef, proposalHash)` PAIR (`samePair`) | 1 page + ≤ 100 |

**Why the anchors hold (measured direction).** A REST commit's entry `transitionRef` is caller-supplied (`logEntryRequest.TransitionRef`, `handlers.go:116`; parsed at `:820`; stored at `:860`, V41), so a REST forger CAN point a new entry at a real record R. That forgery **fails code 2**: R carries its invocation ID, so its bytes are unique to one invocation, and R's first carrier in the append-only `commit_objects` membership is the genuine coordinator commit (V14), never the forger's later entry. To pass code 2 a forger must therefore fabricate a NEW, label-forged `world/invocation-record/v2`, plus its `world/effect-plan/v1`, `world/effect-record/v1` and `world/effect-result/v1`, all through REST's semantic-ID passthrough (V15, Future Work #2). Even then, code 5 requires a `requestRef` that the broker-only approvals chain holds (V16). Under the pair-scoped code 7, naming a real request together with a fabricated plan no longer collides (see the next subsection). **Declared residual:** a request minted by `world-publish` has no genuine packet, so before M2 a REST forger that fabricates that whole label-forged chain (record v2, plan, effect record, effect result and packet) naming such a request passes codes 1–7. After M2 the REST reservation refuses the packet label itself, which closes this; what remains is direct DB tampering, or a store written before M2.

**Request-identity collapse interacts with code 7 (round 2, `oc-kimi-k3`).** An approval request's identity is its content, `{effect, scope, cost, requester, now}` (`approve.go:110-112`, V10), and `appendApprovalHead` appends a new head on every call even for an identical `requestRef` (`approve.go:194-213`, V20). So two genuine invocations that request approval at the same logical `now` with the same scope, cost and requester share ONE `requestRef`, and each mints its own packet. A requestRef-only uniqueness check would refuse both genuine packets. Code 7 therefore compares the `(requestRef, proposalHash)` pair: a packet is unique iff no other stored packet decodes to the same `requestRef` AND the same `proposalHash`. **The forgery argument, re-derived for the weaker law:** a forger who names a real request with a fabricated plan no longer trips code 7. Resistance therefore rests on three checks. Code 2 requires the packet to be first carried with a record whose `Plan` is `proposalHash`; a real record cannot be re-carried first (V14, V41), and a fabricated one needs label-forged objects (V15). Code 3 binds `proposalHash` to that record's `Plan`. Code 5 requires the broker-only approvals chain (V16). Against REST this is unchanged since V16: REST cannot move the chain, so the residual stays the one declared above. (Derived from D2's field set: two collapsed invocations whose plans are ALSO byte-identical and whose grants share an expiry mint byte-identical packets, which is one object. Its first carrier validates it under code 2.)

**Bounded states are not verdicts.** If code 5 fails and `RecentApprovals` reported `Truncated`, or code 7 passed on a truncated scan, the sentence says "not established within the bound" rather than refuting or affirming membership or uniqueness (D6).

### D6 — The projection (object page, no new key)

When `?object=<ref>` names an object whose `SemanticID == DecisionPacketV1`, the existing object branch (`workbench.go:473-483`) additionally calls `d.packetView(ctx, object)` and sets `ObjectView.Packet` (a new `*workbench.PacketView`; `ObjectView` has no such slot today, `render.go:46-58`, V19). Every other object page is byte-identical and does zero extra reads.

| State | Frozen sentence (asserted exactly by tests) |
|---|---|
| validated | `decision packet validated: committed with its invocation record at entry #E; proposal, approval request and approvals-chain anchor agree; lifecycle fields are the v1 producer's` (+ ` (duplicate check bounded to the first 100 stored packets)` when that scan truncated) |
| code 1 | `not a validated decision packet: the bytes are not a canonical world/decision-packet/v1 encoding whose packetHash recomputes` |
| code 2 | `not a validated decision packet: it was not committed together with its invocation record` |
| code 3 | `not a validated decision packet: proposalHash is not the plan of the invocation record committed with it` |
| code 4 | `not a validated decision packet: requestRef is not the approval request that invocation's Human.Approve effect returned` |
| code 5 | `not a validated decision packet: its approval request is not in the broker's approvals chain` — or, when the 40-head read truncated, `not a validated decision packet: its approval request is not among the newest 40 approval heads, so its approvals-chain membership is not established` — or `… the approvals chain is malformed at <ref>` |
| code 6 | `not a validated decision packet: its lifecycle fields disagree with its request's logical time or are not the v1 producer's` |
| code 7 | `not a validated decision packet: another stored packet names the same approval request and the same proposal` |

A **validated** panel additionally renders: `createdAt`/`deadlineAt` as logical integers with the fixed note `logical times; this host emits no Timeout transition (item 7 owns the emitter), so this page does not state whether the deadline has passed`; `policy Cancel · escalationsRemaining 0`; links to the proposal (`R.Plan`), the request, and entry #E; the request's status as the approvals chain reports it, labelled `approvals chain: pending | approve by <who> | deny by <who>` (approvals keep their name); and **evidence** — two grade lines, `proposal <ref>` and `transition function <R.TransitionFn>`, the latter using `R.TransitionFn`; each line is exactly `d.objectGradeView(ctx, ref)` (`workbench_grade.go:142`, V18). PROVEN can appear only through row 105's validator. A **refused** panel renders the sentence alone: no links, status or evidence are claimed for a packet that does not validate.

### D7 — Bounds and cost

Worst case per packet page, all under the request's existing `readDeadline`: 2 `ObjectCommits` lookups, 1 log entry, 6 point reads (record, plan, effect record, result, request, plus the page's own), ≤ 80 object reads + 1 registry-head read inside `RecentApprovals` (V17, V47), 1 packet page of 101 rows with ≤ 100 point reads, and 2 × row 105's bound (≤ 1 page + 100 validations each). REST-written objects are capped by `maxCommitBytes` (8 MiB, `daemon.go:118`, V15); `GetObject` is the existing workbench read surface.

### D8 — Pure core vs host boundary (S1/S2/S3)

The laws live in `design_docs/sketches/decisionpacket.ail` (169 lines, checked: 9/9 contracts Z3-VERIFIED, 37 named inline rows with the harness reporting `Passed: 46 Failed: 0`, and three mutants all `VIOLATION`, V24): `policyCode`, `scheduleWellFormed` (a restatement — a sketch cannot import `world/`), `mintable`, `mintedScheduleWellFormed`, `producerFieldsOk`, `sameRef`, `packetRefusal` (fact 7 = `pairUniqueOk`), `samePair` (the code-7 pair comparison), `lifecycleOk`. Go mirrors live in `host/coordinator/decisionpacket.go`, and `decisionpacket_drift_test.go` evaluates every sketch row through them using the effectplan drift parser (`effectplan_drift_test.go:71`, same package, V27). The gate pins the sketch by name: one `LEG1_MODULES` line plus one `REQUIRED_VERIFIED` entry, outside the `world/` total, the precedent the effectplan and execargs sketches set (V23). **Why not a package?** The Go host must decide at commit and page time without spawning an interpreter; the AILANG sketch is the specification, and the drift test keeps the two from diverging. Hashing, JSON and store reads are host-boundary Go.

The headline law, as checked:

```ailang
-- P2: every packet the producer mints has a well-formed schedule:
-- createdAt = now, deadlineAt = the grant's expiry, escalationsRemaining = 0.
export func mintedScheduleWellFormed(now: int, grantExpiresAt: int) -> bool ! {}
ensures { (not (now >= 0 && now < grantExpiresAt)) || result }
tests [
  ((0, 1), true),
  ((3, 100), true)
]
{
  scheduleWellFormed(now, grantExpiresAt, 0)
}
```

(Two measured fluency facts: an `ensures` that calls another user function was SKIPPED as unencodable, so the condition is inlined; and `0 - 1` is not a legal inline-test literal while `-1` is — V24.)

### D9 — Failure table

| Failure | Observable |
|---|---|
| `Human.Approve` denied/failed | no request, no packet; the invocation commits as today |
| pending output unparseable / grant not found / mint error after the effect ran | `EffectsUnrecordedError` via the existing `unrecorded` path (fail closed, effect records retained) |
| REST commit carrying a packet object | 400, nothing written |
| store/timeout error during validation | existing `writeWorkbenchStoreError` classification (503/500), never a sentence |
| any of codes 1–7 | that code's frozen sentence; no links, status or evidence |

## Implementation Plan (milestones)

- **M1 — laws, codec, drift (~0.6 d; ~120 prod LOC).** Land the banked sketch as `design_docs/sketches/decisionpacket.ail` with the two `verify_ail.sh` manifest lines; `host/coordinator/decisionpacket.go` (`DecisionPacketV1`, encode/decode, `packetHash`, the eight mirrors, `mintPacket`); `decisionpacket_drift_test.go`; codec tests. Gate: `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh` green; `go test ./host/coordinator -run 'Packet|Drift' -count=1` green.
- **M2 — producer + reach + reservation (~1.0 d; ~90 prod LOC).** The D1 hook and the `grants` parameter; `broker.AuthorizingGrant`, `broker.ParsePendingApproval`, `broker.DecodeApprovalRequest`; the daemon binding with `serializedHandler`; the REST reservation. Tests: AC1–AC3, AC7, AC8. Gate: `go test ./host/coordinator ./host/broker ./host/daemon -count=1`.
- **M3 — validator + projection (~1.1 d; ~200 prod LOC).** `ValidatePacket`, `PacketView` + template block, `d.packetView`, and the object-page branch. Tests: AC4–AC6 including every forged arm. Gate: `go test ./host/... -count=1`, and `TestWorkbenchObjectGrade` green with its file unmodified.
- **M4 — S7 usage + re-drill (~0.6 d).** A `docs/QUICKSTART.md` subsection (skill declaring `Human.Approve` → session with the grant → invoke → find the packet via `GET /v1/objects/by-semantic-id/world/decision-packet/v1` → open its object page), executed by a quickstart test (the `quickstart_exec_test.go` precedent, which reads `docs/QUICKSTART.md` and checks a section's commands against `serve --help`, V46); the full §8 mutation re-drill. Gate: both pinned verify gates green (controller-run outside any sandbox).

### Files

**New:** `design_docs/sketches/decisionpacket.ail`; `host/coordinator/{decisionpacket.go, decisionpacket_drift_test.go, decisionpacket_test.go}`; `host/daemon/{workbench_packet.go, workbench_packet_test.go}`. **Modified:** `scripts/verify_ail.sh` (2 lines); `host/coordinator/effectful.go` (~30 LOC) + `coordinator.go` (1 call site); `host/broker/{decide.go, approve.go}` (~35 LOC of pure exports) + their tests; `host/daemon/{daemon.go, handlers.go, workbench.go}` (~35 LOC); `host/workbench/render.go` (~25 LOC); `docs/QUICKSTART.md` + `cmd/ailang-worldd/quickstart_packet_test.go`.

## Examples

**Before:** `GET /v1/objects/by-semantic-id/world/decision-packet/v1` returns an empty page on every store; no daemon invocation can request approval.

**After (illustrative):** an episode with grant `{Human.Approve, scope "deploy:demo", ExpiresAt 1760, Budget 1}` invokes a skill whose plan requests `Human.Approve` at logical `now` 1700. Entry #7 carries input, output, plan `P`, record `R` and packet `K = {…,"proposalHash":P,"requestRef":Q,"createdAt":1700,"deadlineAt":1760,"escalationsRemaining":0,"policy":"Cancel"}`. `/workbench?object=K` renders the validated sentence, links to P, Q and #7, `approvals chain: pending`, and two grade lines (each the unconfigured constant or a row-105 grade). A REST client POSTing an object labelled `world/decision-packet/v1` gets a 400.

## Non-Goals

No Timeout emitter, DEFER or escalation (item 7); no `EscalateBounded`/`ExecuteIfGranted` minting; no decision route on the daemon; no `world/proposal/v1`; no packet list page or back-links from request/plan pages; no change to the daemon's tick source (row 150 / H2); no renaming of the approvals pane; no grade computed outside row 105; no store schema change; no `world/*.ail` change.

## Future Work — candidate queue rows (pre-existing defects found, not absorbed)

1. **Approvals-head append is not compare-and-set** (`approve.go:196,211`, V20): serialized here only inside the daemon; the broker should CAS.
2. **REST commit accepts every coordinator/broker semantic ID** (`handlers.go:847`, V15): this row reserves only the packet ID; records, plans and approvals objects remain label-forgeable.
3. **Content-identical approval requests collapse** (request identity = `{effect,scope,cost,requester,now}`, `approve.go:110-112`, V10): two same-`now` same-scope requests share one ref.
4. Item 7 follow-ons: the Timeout emitter and the other two policies with their declaration surface; a packet list and back-links; replay re-derivation of the packet (needs the grant expiry recorded).
5. **TR.C is blind to Registry-map insertions** (found by row 172; the daemon's existing Workspace.Exec binding is the same shape) — V45, `daemon.go:780-782`.

## Risks & Mitigations

| Risk | Impact | Mitigation |
|---|---|---|
| D4's `Human.Approve` binding is judged an authority-class change | Med | **Analysis (stated, not deferred):** the binding stays within the existing authority model. `Human.Approve` still needs a skill that DECLARES it AND a live, budgeted, exact-scope grant (`Decide`, `decide.go:48`, V13). Grants come only from attended session minting. A request never decides: decisions stay operator-only (`DecideApproval`), and publish consumption also needs an `approve` decision and a publish-shaped, hash-stamping scope (`validatePublishApproval`, V44). The flag is kept for confirmation; if refused, M1 and M3's validator land and M2's binding is held |
| Uniqueness is non-monotone in store growth (round 2, `oc-kimi-k3`) | Med | With a requestRef-only code 7, a LATER genuine invocation whose request collapses onto an earlier one (V10) would flip the earlier, validated packet to refused. Validation would then depend on unrelated future writes. The pair scope resolves this for collapsed requests with distinct plans (AC5 7b, the AC1/AC3 two-approvals arm). Disclosed residual: a second genuine packet with the same pair but a different `deadlineAt` (a different authorizing grant) still flips both to code 7 |
| D2b's byte form proves wrong after packets are stored | Med | It is irreversible by design (amendments are `/v2`, V42); hence the flag and the high change cost |
| Pending requests on the daemon can never be decided there | Med | Honest: the panel shows the chain's `pending` and states that no Timeout is emitted; deciding is item 7's scope |
| The deadline is the grant expiry, coarser than a per-decision TTL | Low | It is the ratified "strict TTL" the effect already carries; a declared per-packet TTL is a `/v2` question |
| Daemon `now` is wall-clock seconds, so packets differ across runs | Med (row 150) | Unchanged by design (V7, V28); row 150 owns its clock decision |
| Validation residual (world-publish requests, pre-M2 stores) | Low | Declared in D5; REST is closed for the packet ID after M2 |

## §7 Acceptance criteria

**AC1 — the producer mints exactly one packet, and only for an allowed, successful `Human.Approve`.** `TestProducerMintsPacketForApprovalOnly`: a real coordinator, real store, a binder registering `broker.NewHumanHandler`; arms OK → 1 packet, denied (expired grant) → 0, a non-approval effect → 0, a refusal plan → 0, each counted through `ObjectsBySemanticID`. Two-approvals arm (round 2): two allowed approvals in one test store at the same `now` and scope (same requester and cost, so V10 collapses them onto one `requestRef`), from plans that differ in their effect ID → 2 packets. Killed by MP1.

**AC2 — the packet's fields are the commit's facts.** `TestProducerPacketBindsInvocationFacts`: `createdAt` == the request object's `now` == the intent's `LogicalTime`; `deadlineAt` == the authorizing grant's `ExpiresAt` (two live grants with distinct expiries, the test asserting the one `broker.Decide` allows first); `proposalHash` == the v2 record's `plan`; `requestRef` == the pending output's ref; `policy` `Cancel`, budget 0; `packetHash` recomputes. Killed by MP2–MP4, MP6, MP7.

**AC3 — the packet is carried by the invocation's own commit.** `TestProducerPacketIsCarriedByTheInvocationCommit`: `ObjectCommits(packet)` == `[E]` == `ObjectCommits(record)`'s first entry, plus a daemon end-to-end arm (`TestDaemonInvocationMintsDecisionPacket`: real daemon, fake interpreter, a skill declaring `Human.Approve`, a session holding the grant). The two-approvals store from AC1 is also checked: both packets are carried by their own invocation commits, and both render the validated sentence on their object pages. Killed by MP5, MD2 and MP8b.

**AC3b — replay and reconcile tolerate the packet object.** `TestReplayReconcileToleratesPacketObject`: after AC3's commit, `Coordinator.Replay` succeeds and its recomposed output equals the committed output byte-for-byte. A resent task ID reconciles through `committed()` with the same `Result` as a packet-free invocation of an otherwise identical skill (the `Output` bytes are equal). Both paths read the record and the objects it names, never the commit's object list (V38). Killed by MP8.

**AC4 — a produced packet validates and renders the validated panel.** `TestWorkbenchPacketValidated`: the AC3 store; the packet's object page contains the validated sentence exactly once, the three links, `approvals chain: pending`, two grade lines equal to `objectGradeView` on the proposal and the transition function, and the logical-time note. Non-packet object pages are byte-identical to base, with zero added `ObjectsBySemanticID(DecisionPacketV1…)` calls (counting wrapper). Killed by MV10–MV12.

**AC5 — every forged or mismatched packet renders its own refusal.** `TestWorkbenchPacketRefusals`, each arm committed through `store.Commit` directly (the test seam REST no longer offers): (1) non-canonical whitespace bytes; (1b) a wrong `packetHash`; (2) a later commit re-carrying the real record with a forged packet; (3) `proposalHash` naming a different stored plan, and (3b) a missing object; (4) `requestRef` naming a different chain request; (5) a request object never appended to the chain; (5b) 41 newer heads → the "not established" sentence; (6) `createdAt` = request `now` + 1, and (6b) `policy` `EscalateBounded`; (7) a second packet naming the same request AND the same `proposalHash` (the same pair) → the code-7 sentence; (7b) a second packet naming the same request BUT a different `proposalHash` → BOTH render validated. Each refusal arm asserts its exact sentence and the absence of links. Killed by MV1–MV9.

**AC6 — no PROVEN is invented.** `TestWorkbenchPacketEvidenceOnlyThroughRowOneOhFive`: on a configured daemon (row 105 wiring) with a subject-bound envelope for the transition function, that line renders `PROVEN`, and the proposal line renders row 105's no-evidence reason; on an unconfigured daemon both lines render row 105's constant. Killed by MV10.

**AC7 — REST cannot store a packet object.** `TestRESTCommitRefusesDecisionPacketObjects`: 400 with the exact message; the store's packet page stays empty; a control commit with another semantic ID still succeeds. Killed by MD1.

**AC8 — `Human.Approve` is serialized across episodes.** `TestHumanApproveIsSerializedAcrossEpisodes`: `serializedHandler` around a fake inner handler that records its maximum concurrent entries and waits up to 200 ms for a second entrant; two episodes dispatch at once; the asserted max is 1. Killed by MD3.

**AC9 — the laws and their mirrors cannot drift.** `TestDecisionPacketSketchRows` evaluates every inline row of the sketch through the Go mirrors and fails on any law without a mirror (or the reverse); `verify_ail.sh` requires the sketch's eight identities by name. Killed by MS1–MS3.

## §8 Mutation table

| ID | Neutered edit | Killer | Predicted failure |
|---|---|---|---|
| MP1 | producer gate drops `results[0].Status == StatusOK` | AC1 | the denied arm returns `EffectsUnrecordedError` instead of committing with 0 packets |
| MP2 | `mintPacket(now+1, …)` | AC2 | `createdAt != request now` |
| MP3 | `deadlineAt = now + 1` instead of `grant.ExpiresAt` | AC2 | `deadlineAt = 1701, want 1760` |
| MP4 | `proposalHash = in.Hash` (input) | AC2 | `proposalHash != record.plan` |
| MP5 | packet stored with `PutObject` instead of `Commit.Objects` | AC3 | `ObjectCommits(packet) = []` |
| MP6 | policy token `"EscalateBounded"` | AC2, AC9 | field assertion; drift `producerFieldsOk` |
| MP7 | `packetHash` hashed over the full payload | AC2, AC4 | recompute mismatch; the produced packet refuses with code 1 |
| MP8 | the replay/reconcile path enumerates the commit's objects and treats any object the record does not name as divergence | AC3b | `Replay` returns `diverged("record", …)` / reconcile an `IntegrityError` for the packet-carrying invocation |
| MV1 | drop the canonical re-encode check | AC5 (1) | whitespace arm renders validated |
| MV2 | drop the `packetHash` recompute | AC5 (1b) | wrong-hash arm passes code 1 |
| MV3 | drop `ObjectCommits(R) first == E` | AC5 (2) | re-carry forgery validates |
| MV4 | drop `R.plan == proposalHash` | AC5 (3, 3b) | mismatched/missing proposal validates |
| MV5 | drop result→`requestRef` equality | AC5 (4) | other-request arm validates |
| MV6 | drop the approvals-chain membership check | AC5 (5) | unanchored request validates |
| MV7 | `lifecycleOk` mirror drops `createdAt == requestNow` | AC5 (6), AC9 | arm validates; drift row `((10, 11, 20, 0, 0), false)` reds |
| MV8 | drop the duplicate scan | AC5 (7) | duplicate validates |
| MP8b | revert the pair scope to a requestRef-only comparison (Go `samePair` mirror ignores `proposalHash`) | AC5 (7b), the AC1/AC3 two-approvals arm, AC9 | 7b renders the code-7 sentence instead of validated; drift row `samePair` 2 (`…p1…p2 → false`) reds; sketch side measured as `VIOLATION samePair` and `samePair_test_2` failing (V24, mutant C) |
| MV9 | treat a truncated miss as definitive | AC5 (5b) | sentence says "not in the chain" |
| MV10 | `objectGradeView(ctx, packetRef)` for the transition-function line | AC4, AC6 | the PROVEN line is missing |
| MV11 | a refused panel still renders links | AC5 | links are present on a refused arm |
| MV12 | the packet branch runs for every object | AC4 | the counting wrapper sees 1 call on a non-packet page |
| MD1 | REST reservation removed | AC7 | 200 and a stored packet |
| MD2 | `d.binder` without `Human.Approve` | AC3 (daemon arm) | `EffectsUnsupportedError` |
| MD3 | `serializedHandler` mutex removed | AC8 | max concurrent entries = 2 |
| MS1 | sketch `mintedScheduleWellFormed` body budget `-1` | verify_ail.sh Leg 1 | `VIOLATION mintedScheduleWellFormed` (measured, V24) |
| MS2 | sketch `packetRefusal` lifecycle code `6 → 0` | verify_ail.sh Leg 1 | `VIOLATION packetRefusal` (measured, V24) |
| MS3 | Go `policyCode` maps `"cancel"` to 0 | AC9 | drift row `("cancel", -1)` reds |

Structurally unreachable branches (the `AuthorizingGrant` miss after an allowed dispatch; `ParsePendingApproval` on the landed handler's own output) are stated, not claimed in this table.

## §9 Conflict Surface

| Neighbor | Collision and ordering |
|---|---|
| `verify_ail.sh` Leg-1 exact module manifest + `REQUIRED_VERIFIED` (V23) | +1 module line, +1 sketch entry in the SAME commit as the sketch; `EXACT_TOTAL_VERIFIED=16` counts `world/` only and does not move; S8's six-file floor is untouched (no `world/` change) |
| TR.C broker binding boundary (`invoke_boundary_test.go:246+`, `wantCount = 3`, V22, V45) | TR.C pins "no raw broker binding outside host/broker, and exactly the 3 `publish_op.go` exemptions inside it" (`outside_broker_is_clean` requires zero findings; `inside_broker_exemption` requires `wantCount = 3` with fixed identities). Its census enumerates every non-test `.go` file in the repository (`enumerateProductionGoFiles`: `filepath.WalkDir`, skipping `.git`, nested modules and `_test.go`), AST-parses each, and matches only `*ast.SelectorExpr` nodes named `Invoke` (any receiver), or `Session` / `NewSession` / `NewReplaySession` on a `host/broker` import, plus dot-imports. This is verified, not asserted (V45), so the daemon's `Registry` insertion is outside the scan. The same blindness already holds for the daemon's existing `Workspace.Exec` binding (Future Work #5) |
| Workspace registry exact-key tests (`workspace_test.go:203-213`, V21) | Untouched: the binding is added in `d.binder`, not in `registry()` |
| Coordinator replay and reconcile (`effectful.go:354+`, `coordinator.go:202+`, V38) | The packet is outside the recomposed output and is not referenced by the v2 record; replay and `committed()` read the record and the objects it names only (`ObjectCommits` appears 0 times in non-test `host/coordinator`), pinned by AC3b. Commits for other effects are unchanged |
| `why` (`why.go:181,665`, V26) | Walks the record and does not enumerate commit objects; unchanged |
| Row 105 grade projection (`workbench_grade.go`, V18) | Consumed unchanged; `TestWorkbenchObjectGrade` and its file untouched |
| Row 149 decisions pane (`workbench.go:582-585`, V17) | Unchanged; daemon-minted requests now appear there, as approvals |
| `host/evidence`, `host/store`, `world/`, `host/boundary`, AC10 subprocess census | No file changes and no new subprocess site |
| Row 150 | Consumes this row's packet page (`/workbench?object=<packet>`); keeps its clock and capture decisions |

## §10 Verification Log

All commands were run from this worktree root at `e869bfe` on 2026-10-13 with `/usr/bin/grep` (this rig's default `grep` is ugrep). Negative rows carry a same-call positive control.

| ID | Claim | Command (and control) | Observed |
|---|---|---|---|
| V1 | Base, clean worktree | `git rev-parse HEAD; git status --short` | `e869bfe4ec4aa0523bbe059494b7bf62d762a6e9`; empty before this doc |
| V2 | No packet producer; approvals positive control | `grep -rlE 'DecisionPacket\|decision-packet\|decision_packet' host cmd --include='*.go' \| grep -v _test.go \| wc -l`; control `grep -rl 'world/approval-request/v1' … \| wc -l` | `0`; control `1` = `host/broker/approve.go` |
| V3 | The frozen type and its laws | `sed -n 152,167p world/types.ail`; `grep -n "wellFormedSchedule\|timeoutOutcome" world/types.ail` | `:154-158` comment ("LOGICAL times … no field reads a wall clock. Amendments are /v2"), `:159-167` `export type DecisionPacket = { packetHash, proposalHash, requestRef, createdAt: int, deadlineAt: int, escalationsRemaining: int, policy: TimeoutPolicy }`; `:172` `timeoutOutcome`, `:196` `wellFormedSchedule`; no wire encoding is declared |
| V4 | No stored `Proposal` | `grep -rlE '"world/proposal' host cmd world packages --include='*.go' --include='*.ail' \| wc -l`; `grep -rn 'proposal/v' host cmd --include='*.go' \| wc -l`; controls `grep -rl '"world/effect-plan/v1"' host cmd … \| grep -v _test \| wc -l`, `grep -rn 'approval-request/v' … \| wc -l` | `0`, `0`; controls `1`, `2` |
| V5 | Approval requests: one mint site, one production caller | `grep -rn "EffectHumanApprove\b\|NewHumanHandler(\|newHumanHandler(\|MintAttendedApproval(" host cmd --include='*.go' \| grep -v _test.go` | `approve.go:17,57,61,103,547,550`; `publish_op.go:125,131,135,136,138`; `cmd/world-publish/main.go:270` — no daemon site |
| V6 | The daemon binds no `Human` effect | `grep -rn "Human" host/daemon --include='*.go' \| grep -v _test.go \| wc -l`; control `… EffectWorkspaceExec … \| wc -l`; `sed -n 780,782p host/daemon/daemon.go`; `sed -n 338,362p host/daemon/workspace.go` | `0`; control `2`; `broker.OpenBinder(d.store, episodeID, grants, d.workspace.registry(episodeID))`; `registry()` returns `broker.Registry{}` or `make(broker.Registry, …)` holding the AILANG family + `Workspace.Exec` only |
| V7 | Daemon `Now` is wall-clock; the coordinator takes it injected | `grep -n "time.Now().Unix" host/daemon/daemon.go`; `grep -n "Now       func() int64\|cfg.Now == nil" host/coordinator/coordinator.go` | `daemon.go:711 … Now: func() int64 { return time.Now().Unix() }`; `coordinator.go:52`, `:78` |
| V8 | One logical `now`; the plan is hashed before effects and committed with the record | `grep -n "now := c.cfg.Now()\|planRef := hashref.SumSHA256\|requestEffect(hctx, bound, e, now)\|planEffectInvocation(world" host/coordinator/effectful.go`; `grep -n "in, out, po := \|in, out, po, rec}" host/coordinator/plan.go`; `sed -n 80,86p host/coordinator/plan.go` | `:284 now := c.cfg.Now() // one logical time`, `:293 planRef`, `:308 requestEffect(…, now)`, `:333 planEffectInvocation`; `plan.go:121,132` objects `{in, out, po, rec}`; `object()` sets `InterfaceHash = SumSHA256(semanticID)` and `Provenance = by` |
| V9 | L1: at most one effect per v1 plan | `grep -n "maxPlanEffectsV1  = 1\|EffectPlanV1 = " host/coordinator/effectplan.go` | `:30 EffectPlanV1 = "world/effect-plan/v1"`, `:32 maxPlanEffectsV1 = 1 // L1` |
| V10 | The request records `now`; the handler returns the pending wire | `grep -n "requestPayload := mustApprovalJSON\|Requester: input.Requester, Now: req.Now\|return pendingBytes(requestObj.Hash)\|Status: \"pending\", RequestRef" host/broker/approve.go` | `:110`, `:112 … Now: req.Now`, `:121`, `:278` `{"status":"pending","requestRef":…}` |
| V11 | The effect result is the handler output, referenced by the record | `grep -n "resultObject(\|ResultRef:" host/broker/broker.go`; `sed -n 111,113p host/broker/record.go` | `broker.go:311 resultObj := resultObject(result)`, `:319 … ResultRef: resultObj.Hash`; `resultObject(payload)` = `brokerObject(EffectResultV1, payload)` |
| V12 | `EffectRecord` fields | `sed -n 19,30p host/broker/record.go` | `Effect, Scope, Cost, BudgetBefore, BudgetAfter, Allowed, Failed, Denial, RequestRef, ResultRef`; `DecodeRecord` exported at `:58` |
| V13 | Broker law surface | `grep -n "func Decide\|func decideOver\|func capabilityLive" host/broker/decide.go` | `:48 Decide` (exported), `:66 decideOver` (unexported), `:101 capabilityLive … now >= 0 && now < c.ExpiresAt` |
| V14 | Commit membership is recorded and readable | `grep -n "INSERT OR IGNORE INTO commit_objects\|const objectCommitsSQL\|func (s \*Store) ObjectCommits" host/store/*.go \| grep -v _test`; `grep -n "ObjectCommits(ctx" host/daemon/daemon.go` | `store.go:1196` insert inside `Commit`; `object_commits.go:20` (`ORDER BY entry_index`), `:26`; `daemon.go:511` in `readStore` |
| V15 | REST passes semantic IDs through; none reserved; body cap | `sed -n 840,850p host/daemon/handlers.go`; `grep -c "decision-packet\|reservedSemantic" host/daemon/handlers.go`; control `grep -c "SemanticID" …`; `grep -n "maxCommitBytes *=" host/daemon/*.go` | `:847 … SemanticID: object.SemanticID,`; `0`; control `8`; `daemon.go:118 maxCommitBytes = 8388608` |
| V16 | REST cannot move registry heads | `grep -rn "SetRegistryHead(\|CompareAndSetRegistryHead(" host cmd --include='*.go'` (non-test, non-store, non-interface); `sed -n 79,87p host/daemon/handlers.go \| grep -ci registry` | callers: `approve.go:211`, `transitionreg.go:276`, `registry.go:168`; `commitRequest` registry fields `0` |
| V17 | `RecentApprovals` bounds, reused by the workbench | `grep -n "maxHeads = min(maxHeads, 40)\|p.Truncated = !head.IsZero()\|if len(rows) > 20" host/broker/approvals_view.go`; `sed -n 582,592p host/daemon/workbench.go` | `:61`, `:101`, `:103`; `workbench.go:582 RecentApprovals(ctx, d.reads, 40)`, `:584` malformed → `Unavailable` text |
| V18 | Row 105's grade view and its call site | `grep -n "func (d \*Daemon) objectGradeView\|gradeCandidateScanLimit *=" host/daemon/workbench_grade.go`; `grep -n "objectGradeView(ctx, ref)" host/daemon/workbench.go` | `:142`, `:18 = 100`; `workbench.go:473` |
| V19 | `ObjectView` has no packet slot | `grep -n "type ObjectView struct\|References       \*ReferenceView" host/workbench/render.go`; `grep -c Packet …`; control `grep -c Grade …` | `:46`–`:57`; `0`; control `19` |
| V20 | Approvals head: read-then-set, no CAS; single writer process; per-episode locks | `grep -n "GetRegistryHead(ctx, ApprovalsV1)\|SetRegistryHead(ctx, ApprovalsV1" host/broker/approve.go`; `grep -n "exclusive writer lock" host/store/writer_lock.go`; `sed -n 91,93p host/coordinator/effectful.go` | `:196` read, `:211` set; `writer_lock.go:13` "exclusive lock on \"<canonical-db>.writer.lock\""; `lockEpisode` keys `episodeLocks` by `episodeID` |
| V21 | Registry exact-key tests call `registry()`, never `d.binder` | `sed -n 200,214p host/daemon/workspace_test.go`; `grep -rn "d.binder(\|\.binder(\"" host/daemon/*_test.go \| wc -l` | `registry(ep)` empty / `registry("ep1")` = `wantWorkspaceEffects`; `0` |
| V22 | TR.C exemption count | `grep -n "const wantCount = 3" host/broker/invoke_boundary_test.go` | `:274`; the identities are `publish_op.go` only |
| V23 | Gate manifest and the sketch precedent | `grep -n "^EXACT_TOTAL_VERIFIED=\|^LEG1_MODULES=(\|\"design_docs/sketches/effectplan.ail\": {" scripts/verify_ail.sh`; `sed -n 152,158p scripts/verify_ail.sh` | `:174`, `:331`, `:403 =16`; "Sketches are excluded from the world/ total"; `ROOTS` sweeps `design_docs` — so a stray `.ail` there reds Leg 1b (why the fixture is `.ail.txt`) |
| V24 | The sketch checks, verifies, tests, and all three mutants fail (pinned binary; re-run for round 2's pair scope) | scratch copy `sketches/decisionpacket.ail` (banked `…/decisionpacket.ail.txt`, 169 lines): `ailang check`, `ailang verify`, `ailang test`, `ailang ai-check`, three mutants; `shasum -a 256 ~/.pinned-ailang/ailang` | `✓ No errors found! rc=0`; `9 exported functions: 9 verified`; `Passed: 46 Failed: 0` (37 named inline rows incl. 3 `samePair`; round 1 was `8 verified`, `Passed: 42`); ai-check `verify.errors 0`, `verified 9 of 9`; mutant C (`samePair` body reverted to requestRef-only) → `VIOLATION samePair` and `✗ samePair_test_2`, `Passed: 45 Failed: 1`; mutant A `VIOLATION mintedScheduleWellFormed` (cex `now=0, grantExpiresAt=1`), mutant B `VIOLATION packetRefusal`; first draft: an `ensures` calling `mintable` was `SKIPPED` ("calls user function … not SMT-encodable"), and `0 - 1` test literals failed ("expected literal expression") — both fixed; sha `1a67b014…b9f`. Full log `sketch-check.log` |
| V25 | Ratified TTL and rejected-timeout text; Cancel's proven outcome | `sed -n 161,163p design_docs/HUMAN-SURFACE.md`; `grep -n "Cancel => ResolveRejectedTimeout," world/types.ail` | "every Decision Packet and `Human.Approve` effect MUST carry a strict TTL. If the human does not decide before expiration, the packet deterministically resolves to a typed rejected-timeout"; `:174`, `:180` |
| V26 | `why`'s coordinator flag is structural | `grep -n "coordinator.RecordV1 && obj.SemanticID != coordinator.RecordV2\|ch.Coordinator = true" cmd/ailang-worldd/why.go` | `:181`, `:665` |
| V27 | Drift-test precedent reads its sketch | `grep -n "sketches\", \"effectplan.ail\"" host/coordinator/effectplan_drift_test.go` | `:71` |
| V28 | Clock is row 150's question; the freeze defined no bytes | `grep -n "clock knob" design_docs/verification/world-row150-design/independent-review.md`; `grep -n -i "canonical\|wire\|encod" design_docs/implemented/w-decision-lifecycle-freeze.md` | `:17` "A new production clock knob affecting ordinary invocation/expiry authority warrants the proposed freeze"; the freeze's hits are Z3 encodability (`:142-144`, `:442`), the canonical/projected `types.ail` copy and package golden (`:400`, `:534`, `:571`) and "content-addressed wire type" discipline (`:151`) — none defines packet payload bytes |
| V29 | Handler/Registry types the binding uses | `grep -n "type Registry\|type Handler interface" host/broker/broker.go` | `:24 Handler{ Execute(ctx, EffectRequest, []byte) }`, `:36 Registry map[string]Handler` |
| V30 | New identifiers are unallocated | `for tok in DecisionPacketV1 ValidatePacket AuthorizingGrant ParsePendingApproval DecodeApprovalRequest PacketView mintPacket packetRefusal serializedHandler decisionpacket workbench_packet; do … grep -rn "$tok" host/ cmd/ scripts/ \| wc -l; done`; control `RecentApprovals` | all `0`; control `13` |
| V31 | HashRef text API | `grep -n "func SumSHA256\|func (h HashRef) String\|func Parse(" host/hashref/hashref.go` | `:171`, `:79`, `:127` |
| V32 | The session copies grants; the coordinator holds the seeded grants | `sed -n 122,125p host/broker/broker.go`; `grep -n "grants, err = c.seedGrants\|c.cfg.Binder(call.EpisodeID, grants))\|return c.dispatchEffectful(" host/coordinator/coordinator.go` | `grants: append([]Capability(nil), grants...)`; `:311`, `:317`, `:334` |
| V33 | world-publish's times are operator-supplied logical values | `grep -n "\"now\", 0, \"logical time of this act\"\|RequestedAt: opts.now" cmd/world-publish/main.go` | `:157`, `:572` |
| V34 | Existing coordinator tests never request `Human.Approve` | `grep -rn "Human" host/coordinator/*_test.go \| wc -l` | `0` — the producer hook cannot change any existing arm |

| V35 | `LogEntry` shape (the code-2 hop) | `sed -n 134,140p host/store/store.go` | `:136 type LogEntry struct { Header LogHeader; EntryHash hashref.HashRef; TransitionRef hashref.HashRef }` |
| V36 | Record v2 fields, unexported, in package `coordinator` | `sed -n 38,51p host/coordinator/plan.go` | `:39 type recordV2 struct` with `InvocationID invocationId, EpisodeID episodeId, SkillID skillId, TransitionFn string transitionFn, Interpreter interpreter, SemanticsEpoch int64 semanticsEpoch, Input input, Output output, Plan string plan, Effects []string effects` (closing `}` at `:50`) |
| V37 | `readStore` ⊇ PacketReader | `sed -n 494,513p host/daemon/daemon.go`; `awk '/type readStore interface/,/^}/' host/daemon/daemon.go \| grep -c "ctx context.Context"`; `grep -n "six" host/daemon/daemon.go` | `:501-513`, **11** methods incl. `GetObject`, `GetLogEntry`, `GetRegistryHead`, `ObjectsBySemanticID`, `ObjectCommits`; observation only: the doc comment at `:487`/`:493` still says "six" (stale comment, not touched by this row) |
| V38 | Replay/reconcile never enumerate commit objects | `grep -n "ObjectCommits" host/coordinator/*.go \| grep -v _test \| wc -l`; control `grep -c ObjectCommits host/daemon/daemon.go`; `grep -n "Objects" host/coordinator/*.go \| grep -v _test`; `sed -n 202,258p host/coordinator/coordinator.go`; `sed -n 354,480p host/coordinator/effectful.go` | `0`; control `1`; `Objects` only at `plan.go:77,151,154` (the builder); `committed()` reads record → output → world by ref; `Replay` reads record → source/input/plan by the record's refs → effect records → recomposes and compares output bytes (`bytes.Equal(outBytes, committed.OutputBytes)`) — no commit object list is read |
| V39 | The semantic-ID scan is deterministic | `sed -n 18,20p host/store/objects_by_semantic_id.go` | "hash_ref is the PRIMARY KEY, so it is a unique total order"; `:19-20 … FROM objects WHERE semantic_id = ? AND hash_ref > ? ORDER BY hash_ref LIMIT ?` |
| V40 | M4's lookup route exists | `grep -n "by-semantic-id" host/daemon/daemon.go` | `:903 mux.HandleFunc("GET /v1/objects/by-semantic-id/{name...}", d.handleObjectsBySemanticID)` |
| V41 | A REST entry's `transitionRef` is caller-supplied | `grep -n "TransitionRef" host/daemon/handlers.go` | `:116 TransitionRef string \`json:"transitionRef"\`` (`logEntryRequest`); `:820 transitionRef, err := parseRef(r.Entry.TransitionRef, "entry.transitionRef")`; `:860 … TransitionRef: transitionRef` |
| V42 | The content-addressed wire-type discipline D2b honours | `sed -n 149,151p design_docs/implemented/w-decision-lifecycle-freeze.md` | "amendments after freeze are NEW versions (the reserved semantic ID `world/decision-packet/v1` becomes `/v2`), never in-place edits — the same discipline as the frozen `LogHeader` and every content-addressed wire type" |
| V43 | Revision base: code unchanged since `e869bfe` | `git rev-parse --short HEAD; git log --oneline e869bfe..HEAD; git diff --stat e869bfe HEAD` | HEAD `955fbc0` = `e869bfe` + the controller's commit of this design alone (3 files: this doc, `decisionpacket.ail.txt`, `sketch-check.log`) — every V35–V42 row observes the `e869bfe` code |

| V44 | What publish consumption enforces (D4's authority premise; round 2, `gemini-3-1-pro`) | `grep -n "func validatePublishApproval" host/broker/approve.go`; `sed -n 486,619p host/broker/approve.go`; `grep -n 'decision.Decision != "approve"\|request.Effect != EffectHumanApprove\|ParsePublishApprovalScope(request.Scope)\|scope.Effect != EffectRegistryPublish\|scope.Publish != wantPublish\|"tarball hash", scope.TarballSHA256\|req.Now > scope.ExpiresAt' host/broker/approve.go` | `:486 func validatePublishApproval(ctx context.Context,`; the decision object must be `world/approval-decision/v1` with `decision.Decision != "approve"` → `ErrPublishApprovalDenied` (`:513`); its request must be `Human.Approve` (`:547`) with cost `PublishCost`; `scope, err := ParsePublishApprovalScope(request.Scope)` (`:558`); `scope.Effect != EffectRegistryPublish` refused (`:562`); `scope.Publish != wantPublish` refused (`:572`); the scope must stamp the manifest, tarball (`:584`), content and interface hashes; logical-time ordering, then `req.Now > scope.ExpiresAt` → `ErrPublishApprovalExpired` (`:613`) |
| V45 | TR.C's census mechanism and invariant (round 2, `oc-glm-5-3`) | `sed -n 246,280p host/broker/invoke_boundary_test.go`; `sed -n 110,160p …`; `sed -n 45,108p …` | `TestRegistryDispatchBindingBoundary` (`:247`): `enumerateProductionGoFiles(root)` (`:124`) walks the repository with `filepath.WalkDir`, skipping `.git`, nested `go.mod` modules and `_test.go`, and keeps `.go`; `parseAndScan` (`:110`) = `parser.ParseFile` + `scanFile`, which runs `ast.Inspect` over `*ast.SelectorExpr` only: `Sel.Name == "Invoke"` (any receiver) → `invoke-call`; on an identifier that is a local name of the `host/broker` import, `Session` → `session-type`, `NewSession` → `ctor-live`, `NewReplaySession` → `ctor-replay`; a dot-import → `dot-import`. Invariant: subtest `outside_broker_is_clean` requires `len(outside) == 0`, and `inside_broker_exemption` requires `wantCount = 3` with identities `publish_op.go|mintAttendedApproval|invoke-call: 2`, `publish_op.go|invokeAttendedPublish|invoke-call: 1`. A `Registry` map literal or index assignment matches no detector |
| V46 | The M4 quickstart-test precedent | `ls cmd/ailang-worldd/quickstart_exec_test.go; grep -n "QUICKSTART" cmd/ailang-worldd/quickstart_exec_test.go` | file present; `:73 guide, err := os.ReadFile("../../docs/QUICKSTART.md")`, `:69` "section10ServeLines returns QUICKSTART §10's … serve", `:118` "serve flag %s is not documented by serve --help" — it executes and checks a QUICKSTART section |
| V47 | `RecentApprovals`'s read count (A5/D7) | `grep -n "GetObject\|GetRegistryHead(ctx, ApprovalsV1)\|len(heads) < maxHeads\|readBody(row\.\|at most 80" host/broker/approvals_view.go`; `grep -n "maxHeads = min\|len(rows) > 20" …` | 1 `GetRegistryHead` (`:63`); the head loop `len(heads) < maxHeads` (`:69`) with `maxHeads = min(maxHeads, 40)` (`:61`) does one `GetObject` per head (`:74`) → ≤ 40; rows are capped at 20 (`:103`), each with one request `readBody` (`:122`) and at most one decision `readBody` (`:132`), each one `GetObject` (`:114`) → ≤ 40. Total: ≤ 80 object reads + 1 head read, matching the function's own comment at `:55` ("at most 80 object reads in total"). The figure stands; the doc now adds the 1 head read |

**Sweep note:** every "the code does / does not" sentence above cites a row; the negative rows V2, V4, V6, V15, V16, V19, V21, V30, V34 and V38 carry same-call controls or are exact enumerations. Every validator anchor hop is now a row: entry → `TransitionRef` (V35), record-v2 fields (V36), reader methods (V37), commit membership (V14), effect record and result (V11, V12), approvals chain (V16, V17), scan order (V39); and so are the replay/reconcile independence (V38), the REST entry surface (V41) and the freeze discipline (V42). Round 2 adds: publish-consumption enforcement behind D4 (V44), TR.C's actual census (V45), the M4 precedent (V46), the `RecentApprovals` read count (V47), and the re-checked pair-scoped sketch (V24).
