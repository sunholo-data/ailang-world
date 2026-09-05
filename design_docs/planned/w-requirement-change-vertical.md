# W-REQUIREMENT-CHANGE-VERTICAL: One Goal, Changed Evidence, One Governed Action

**Status:** Proposed integration design — queue admission requested by Mark, 2026-09-05; implementation not approved; quorum not run.
**Author:** Astra.
**Target:** Companion to existing World clause-5 work; no additional release threshold.
**Priority:** P1 candidate; dependency-gated.
**Estimated:** 4–7 engineering days after prerequisites; broader UI and upstream work excluded.
**Parents:** [HUMAN-SURFACE](../HUMAN-SURFACE.md), [SCENARIOS](../SCENARIOS.md), existing queue item 7 (`w-approval-inbox`).
**Dependencies:** [Evidence applicability](w-evidence-applicability.md); production session/projection/inbox work for interactive deployment; production proof authority for a PROVEN display.
**Cross-project:** AILANG semantic repair packet and lifecycle pilot; the first local deterministic drill can use recorded fixtures.

## Problem and duplicate/coverage decision

World already specifies the human workbench and decision packets. This design adds a bounded integration episode and measurable acceptance to that existing work. It creates neither a new inbox, a second protocol nor another renderer application.

A successful initial proposal does not demonstrate that a human can change its goal, understand which evidence became stale, safely resume after a restart and inspect the eventual action. This episode exercises that chain. The read-only workbench exists; the live full interaction remains dependent on the charter's session authority, projections and approval-inbox items.

## Goals and non-goals

Make one requirement amendment understandable and enforceable through World's existing state and action boundaries. The model receives the same selected revision and evidence applicability that the human sees.

Non-goals: new DecisionPacket fields; new core evidence grades; arbitrary rollback or near-O(1) branch API; new live email/GitHub/publish handlers; changing the resident non-inferiority floor; automatically passing World's real-question value gate with a synthetic drill.

## The episode

Use a software-engineering task with an independently authored oracle: change a configuration-selection function from “use the first matching entry” to “use the last matching entry.” Preserve the invariant that malformed inputs return a structured failure and never cause an outward action. This is a task specification, not an assertion about an existing AILANG API.

1. Record R1, S1 and passing evidence. Create the proposal and its ordinary existing decision packet.
2. Human amends to R2. The authorized selection transition records R2; R1's evidence remains inspectable but is stale for the active goal.
3. A resident gets an AILANG repair packet for the selected source and requirement identities. It proposes S2; independently authored tests evaluate R2 and the retained failure invariant.
4. The pane shows old/new criterion, changed result, affected evidence, actual check verdicts and the exact requested action. An unchanged dependency may be displayed as reused only with matching captured identity and policy.
5. Restart the episode before a decision. Rebuild from stored references, revalidate authority/evidence, and display the same pending outcome with its deadline.
6. Authorize one action. Initial drill uses an inert counting handler inside a confined test host; production later reuses an EXISTING supported broker effect, with no new outward integration in this sprint.
7. Ask why the final action was permitted, which requirement it satisfied, and why the first approval no longer applied. Each answer must link to recorded objects or explicitly state unavailable.

## State and authority rules

Reuse `DecisionPacket/v1`, `TimeoutPolicy`, the existing proposal/evidence references, grants, broker receipts and HUMAN-SURFACE decision verbs. An amendment creates a new requirement revision and proposal/packet; it never edits the old packet's meaning. Defer records a new bounded deadline under existing laws. Expiry never manufactures human approval; ExecuteIfGranted requires separately established current authority.

The applicability report is a snapshot, not permission. Before committing a selected revision or dispatching an action, the domain coordinator must establish that the selected requirement/proposal/authority identities still match. A check performed earlier in the UI is insufficient.

**First integration has one serialized episode executor.** It serializes amendments, approvals and dispatch admission for this domain; external callers cannot bypass it with raw claimed bindings. Capture the selected world/requirement identity and use the store's existing expected-head compare-and-append boundary when recording the decision. Admission and durable action intent must have an explicit linearization point with respect to later amendments. An amendment ordered before that point invalidates the old action; one ordered after dispatch admission cannot retroactively promise the remote action never occurred.

**Blocking design-freeze mechanism:** inventory the existing journal/claim/store transaction APIs before planning the live path. If they cannot atomically bind the relevant current selection and consumed authority to durable intent, the live path remains unsupported. Propose the smallest justified host change in a separate reviewed refinement with a concurrency regression; never simulate atomicity through two adjacent calls. The deterministic drill may proceed, but it cannot be labeled production completion.

A restart after durable intent without an outcome displays indeterminate and uses existing reconciliation. No automatic redispatch, no local promise of exactly-once remote delivery, no cancellation-success claim solely from a cancellation request.

## Pane contract

Extend the existing workbench/approval projection, preserving HUMAN-SURFACE's five zoom levels and grounded-link conventions. Required visible fields:

- Selected criterion and owner; original request link.
- New vs previous proposed result.
- Evidence kind, verdict, scope, applicability and unavailable reasons.
- Exact requested action, authority scope and remaining resource/attention budgets.
- Pending deadline and timeout policy, with explicit state after expiry/restart.
- Links from the action receipt back to proposal, requirement selection and decision.

Grounded prose must not hide missing graph edges. A plausible narrative does not satisfy an unresolved reference. The renderer cannot synthesize PROVEN from an incoming grade string; it consumes only the resolved view permitted by the existing proof authority design. Existing generic grade labels alone are not that proof-to-renderer integration.

## Why this is a package / host boundary

The domain episode, requirement references and repair-packet adapter belong in an extension alongside world-requirements, not `world/` kernel modules. Thin host integration connects existing broker/transition registry and renderer. Each new host responsibility must justify why it cannot live in the package, per S3. Do not add a second auth/session subsystem while row 39 owns that work.

Suggested files after integration inventory: domain extension under `packages/`; focused adapter module under `host/` only where necessary; fixtures in the existing workbench/runbook test suites; an episode guide under `docs/`. No shared driver, compiler, DecisionPacket schema or `packages/world-core` manual edits.

## Phases / stopping points

| Phase | Result | Gate |
|---|---|---|
| RV1 | Recorded, deterministic episode and inert-handler drill | Can start after applicability design lands; cannot claim resident usability |
| RV2 | Existing inbox/projection consumes the episode; session-bound resident can amend/review | Existing session/projection/inbox prerequisites and current-selection admission mechanism verified |
| RV3 | Restart/adversarial drills, then real operator trial | Required checks pass; exact allowed real action separately authorized |

The integration is not complete until RV2/RV3 prerequisites are met. Keep partial states named; do not claim the deterministic fixture delivers the interactive product.

## Acceptance criteria and required failure controls

| ID | Acceptance | Mutation / adversarial case |
|---|---|---|
| RV1 | Selecting R2 makes R1 evidence visibly stale and prevents R1 action admission | `RV-IGNORE-REVISION`: changed criterion, same source |
| RV2 | S2 passes changed requirement and preserved invariant under independent oracle | `RV-WEAKEN-ORACLE`: hardcoded output; malformed-input violation |
| RV3 | Decision packet/authority cannot authorize changed artifacts | `RV-REUSE-APPROVAL`: change target/source after approval; zero handler calls |
| RV4 | Amendment/admission race obeys recorded ordering | `RV-SPLIT-CHECK-INTENT`: pause between check and intent, admit amendment, assert stale action refused |
| RV5 | Restart reconstructs state; unproven receipts remain unproven until validated | `RV-RESTART-TRUST`: serialized seal/grade injection |
| RV6 | Expired/deferred packet follows existing bounded laws without inventing approval | `RV-TIMEOUT-APPROVES`, `RV-FREE-DEFER`: explicit timeout and escalation-budget controls |
| RV7 | Durable intent with no outcome displays indeterminate; reconciliation dispatches zero | `RV-RETRY-UNKNOWN`: lost-response fixture and counting handler |
| RV8 | Three provenance answers identify actual requirement/proposal/receipt references | `RV-MISSING-LINK`: missing object must render unavailable, never plausible prose |
| RV9 | Model and human views identify the same selected revision and applicability | `RV-DIVERGENT-VIEWS`: update one projection only |
| RV10 | Null episode and zero-check oracle are refused; working control reaches one admitted inert action | `RV-EMPTY-GREEN`: no checks/no actions cannot masquerade as completion |

For concurrency tests, a positive control admits an unamended action and a controlled competing amendment changes the selected revision before admission. More timing repetitions cannot replace the explicit interleaving. All mutation names are proposed controls, not reported successful tests.

## Human evaluation

After deterministic controls pass, use counterbalanced task variants with the same operator: existing chat/log workflow versus World pane. Record correct diagnosis, incorrect approvals, time to identify stale evidence, recovery time and interruptions. Report small-sample results descriptively. Do not retell three synthetic fixture questions as the charter's three REAL operational questions; collect those separately during genuine usage.

The AILANG lifecycle pilot owns model/context comparisons. World retains its existing resident non-inferiority and real provenance value criteria; this design adds a scenario that can generate evidence for them, not replacement thresholds.

## High-impact decisions / design freeze

- [ ] Mark approves the first software task and inert-first scope as a refinement of existing inbox work.
- [ ] Existing packet/grade semantics preserved; additive presentation reconciled with HUMAN-SURFACE.
- [ ] Session authority/projection/inbox prerequisites verified at current code, not stale headers.
- [ ] Admission/intent linearization mechanism specified and tested before live execution.
- [ ] Independent oracle and operator-trial protocol frozen; real outward action authorized separately.

Designer may choose layout and fixture organization. Unresolved host integration blocks RV2, not permission to weaken freshness or authority. No new release-gating unit is asserted.

## Verification log / related work

Read `world/types.ail` (frozen packet and timeout/defer laws), `host/transitionreg/bind.go` (descriptor pin checks and confined invoker), `host/broker/{broker,approve}.go` (publish-bound approval and durable intent), `host/store/schema.sql` (selected-head compare/append intent), and `host/daemon/{daemon,workbench}.go` (registered read-only workbench). Existing `TestRecoverCountingProbeDispatchesZeroHandlers` was read and freshly run in the prior review; this supports only that specific recovery behavior.

The exact new admission/intent integration has NOT been proven by this review and is explicitly a design-freeze item, not “confirmed.” No new `.ail` syntax or absence claim is needed. Related designs: [applicability](w-evidence-applicability.md), [frozen lifecycle](../implemented/w-decision-lifecycle-freeze.md), [read-only workbench](../implemented/w-workbench-read-only.md), [HUMAN-SURFACE](../HUMAN-SURFACE.md). This is an acceptance/integration child of item 7, not a duplicate approval-inbox project. No implementation or paid evaluation accompanies it.

## Axiom compliance

Directional design assessment, not implementation approval. No hard violation proposed on A1/A3/A4/A7.

| Axiom | Score | Design constraint |
|---|---:|---|
| A1 Determinism | +1 | Identity-bound inputs; deterministic mechanical results |
| A2 Replayability | +1 | Preserve inputs and outcomes needed to reproduce the decision |
| A3 Effect legibility | 0 | Existing effect semantics unchanged |
| A4 Explicit authority | +1 | Metadata and model output never mint execution authority |
| A5 Bounded verification | +1 | Explicit limits and checks with named refusal outcomes |
| A6 Safe concurrency | 0 | Sequential first version; snapshot identity checked |
| A7 Machines first | +1 | Structured, versioned artifacts with explicit unavailable states |
| A8 Minimal syntax | +1 | No new language syntax |
| A9 Cost visibility | +1 | Record tool/model costs and failure overhead |
| A10 Composability | +1 | Reuse existing compiler, evidence and protocol boundaries |
| A11 Structured failure | +1 | Unknown, incomplete and stale cannot masquerade as success |
| A12 System boundary | +1 | Separate claims, verification, permission and action |

**Net +10.** Re-score if implementation scope changes.
