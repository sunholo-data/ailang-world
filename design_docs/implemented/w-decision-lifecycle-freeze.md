# w-decision-lifecycle-freeze — the typed timeout-policy set and the frozen decision-packet schema

- **Status**: Planned — awaiting pick-time quorum. **This document PROPOSES; Mark RATIFIES.**
  Both deliverables are open ratification points in `design_docs/HUMAN-SURFACE.md` §7 (point 1:
  the typed finite set of timeout policies; point 3: decision-packet schema freeze timing).
  Nothing below is written as though the human choice is already made; §2 states the two asks.
- **Item**: queue item 15, `w-decision-lifecycle-freeze`, clause-5, filed 2026-08-11 (attended, Mark)
- **Estimated**: 1.0 day; at the top of the queued ~1 day band (§9)
- **Measurement base**: `bc8f193`, 2026-08-14
- **Instrument**: `/tmp/ailang-v0300/ailang`, AILANG v0.30.0 (`e37b370`) — the pinned released
  binary. The PATH `ailang` is a `-dirty` dev build and was not used for any claim here.
- **Files changed by implementation**: `world/types.ail`, `packages/world-core/world/types.ail`,
  `scripts/verify_ail.sh`, `host/verifygate/module_manifest_gate_test.go`, and
  `scripts/world_package_ready_packet.golden.json` — all five gate pins are listed in §8; an
  earlier item omitted pin 5 from its conflict surface and red-lighted CI.
- **Design result**: a three-constructor `TimeoutPolicy` ADT, a frozen v1 `DecisionPacket` world
  type, and five Z3-proven pure lifecycle laws (`timeoutOutcome`, `timeoutFiredLegally`,
  `validEscalation`, `validDefer`, `wellFormedSchedule`) — validated end-to-end on an isolated
  copy of the real `world/` tree in both rounds (V-P7 round 1; V-P13 the full revised set),
  including the four mutations Z3 is structurally blind to (V-P8, V-P14).
- **Round-2 revision (2026-08-14)**: round-1 quorum BLOCKED — both reviewers present, both
  REJECT. Disposition, each objection MEASURED rather than adopted or dismissed:
  `gemini-3-1-pro`'s two-unverified-premises objection is procedurally right (§11 had no rows for
  them) and substantively refuted — both premises hold on first-party measurement (V-H, V-I; the
  §3 citation is corrected to name the function line). `gpt5-6-sol`'s objection — the laws never
  connect the DEADLINE to the OUTCOME, and `validDefer` never forces the decrement — is CORRECT,
  and the design is fixed, not the prose: two new proven laws (`timeoutFiredLegally`,
  `validEscalation`) close both gaps. The reviewer's proposed fix AS WRITTEN is Z3-unencodable on
  the pinned binary and fails silently (V-P10) — it was rebuilt in the int-only form that
  verifies (V-P11, V-P13). A third verifier limitation measured en route (V-P12) routes upstream
  with the other two.

Every present-tense codebase claim is backed by a command, observed output, and a same-scope
known-positive control in §11. Rows the controller supplied were re-run first-party at `bc8f193`
before being relied on; each is labelled `VERIFIED BY CONTROLLER` with this designer's re-run
recorded. Probes were isolated under `/tmp/probe-item15/` and `/tmp/iso-item15/`; they are
measurements, not implementation artefacts, and touched no repo file.

---

## 1. Problem

HUMAN-SURFACE.md §3.1 is a **binding principle** applied verbatim from two independent round-2
reviewers, and it is entirely unimplemented. Every decision packet MUST carry a ledger-recorded
creation time, a decision deadline, and a timeout policy from a typed finite set; DEFER MUST
create a new bounded deadline and MUST NOT park indefinitely; at deadline the system emits an
explicit Timeout transition and follows the declared policy; **silence MUST NEVER synthesize
approval or rejection**; and replay MUST reproduce deadline behaviour deterministically from
ledger time. Without this, the approval inbox (item 7) can wedge on a human exactly the way this
loop's own slots wedged on a background agent — Standing Rule 6 restated at the UX layer.

Measured at `bc8f193` (V-D): there is no `DecisionPacket` anywhere in `world/` or `host/`
(0 files; controls fire), no escalation concept (0 non-test files; control 24), and no decision
deadline (`DecisionDeadline|deadlineAt|…` → 0 non-test hits; control fires). The lifecycle is
genuinely unbuilt, not partially built. Naive greps for "Timeout"/"Defer"/"expire" return
Go-idiom noise (`context` deadlines, the `defer` keyword, cache expiry) — none of it is
lifecycle machinery, and none is cited as prior art here.

**The central tension, which the queue row does not mention** (V-A): the world ledger records
NO time. `host/store/schema.sql` has zero time/timestamp/created-at columns across all 8
`CREATE TABLE`s, and the append-only `log_entries` table (schema.sql:36-45) carries exactly the
six frozen `LogHeader` fields plus `entry_hash_ref` and `transition_ref`. So §3.1's
"ledger-recorded creation time" has no column to live in. §3 resolves this without touching the
frozen log.

## 2. The two ratification asks (Mark decides; this doc proposes)

### 2.1 Ask 1 — the typed finite set of timeout policies (§7 point 1)

**Proposed set — exactly the three members the round-2 reviewers named, no additions, no
removals:**

```ailang
export type TimeoutPolicy
  = Cancel
  | EscalateBounded
  | ExecuteIfGranted
```

Mapping to the reviewers' verbatim phrases, with the one liberty taken stated:

| Constructor | Reviewer phrase (§3.1, verbatim) | Rename rationale |
|---|---|---|
| `Cancel` | "cancel" | none — verbatim |
| `EscalateBounded` | "remain safely unexecuted with bounded escalation" | the phrase is a sentence; the constructor keeps its two load-bearing words. The packet REMAINS UNEXECUTED while escalations occur; the bound is `escalationsRemaining` (§2.3). |
| `ExecuteIfGranted` | "execute only if authority was already independently granted" | "only if" and "already independently granted" are enforced by the law's `independentAuthority` input (§2.3), which the host derives from the LANDED broker check (`capabilityLive`, decide.go:101) at the recorded logical time — never minted by the timeout itself. |

**Why exactly three, and why nullary.** No member is added: every candidate fourth policy
examined ("park indefinitely", "auto-approve at deadline") is EXPLICITLY FORBIDDEN by §3.1's own
text (DEFER must not park indefinitely; silence must never synthesize approval). No member is
removed: each of the three covers a §3.1 clause the others cannot (resolve now / stay safe but
keep asking, boundedly / proceed under authority that already exists). The constructors are
nullary because the escalation bound is packet STATE that DEFER decrements, not policy identity —
keeping the set closed and every payload out of it. This makes the set total and finite in the
strongest available sense: a match over it with an exact contract is Z3-proven exhaustive under
the pinned binary (V-P1), and a missing arm is a verifier ERROR, not a silent acceptance (V-P9).

**The resolution semantics proposed for ratification with the set** (these are part of the ask —
they say what each policy DOES at deadline, which the bare names do not):

- `Cancel` at deadline → the packet resolves to the typed rejected-timeout (§3's "deterministically
  resolves to a typed rejected-timeout", verbatim).
- `EscalateBounded` at deadline with `escalationsRemaining > 0` → one escalation is emitted with a
  NEW bounded deadline and the counter decremented; the packet remains unexecuted.
- `EscalateBounded` at deadline with `escalationsRemaining <= 0` → the packet resolves to the typed
  rejected-timeout. **Exhaustion resolves; it never parks** — this is the "MUST NOT park
  indefinitely" clause made total.
- `ExecuteIfGranted` at deadline with live independent authority → the packet proceeds under THAT
  authority (an existing grant, live at the recorded logical time per the ratified broker law).
  Nothing is synthesized; the authority pre-exists the timeout.
- `ExecuteIfGranted` at deadline with NO live independent authority → the packet resolves to the
  typed rejected-timeout. **Silence never synthesizes approval** — and this arm is the one Z3
  cannot defend alone (§5, V-P8), so it carries a named runtime killer.

### 2.2 Ask 2 — the decision-packet schema freeze and its timing (§7 point 3)

**Proposed frozen v1 schema** (kernel world type; §4 has the exact code):

```ailang
export type DecisionPacket = {
  packetHash: HashRef,            -- self hash (precedent: Proposal.proposalHash)
  proposalHash: HashRef,          -- the proposal awaiting decision
  requestRef: HashRef,            -- the ApprovalRequestV1 object this packet fronts
  createdAt: int,                 -- LOGICAL creation time, caller-supplied (§3)
  deadlineAt: int,                -- logical decision deadline; deadlineAt > createdAt
  escalationsRemaining: int,      -- the DEFER/escalation budget; >= 0
  policy: TimeoutPolicy
}
```

**References, not embeds — deliberately.** P1 describes the packet as proposal + evidence +
authority + budget + reversibility. The v1 kernel record carries the lifecycle spine and LINKS
to the rest: `proposalHash` reaches `Proposal`, which already carries `evidence`,
`requiredCaps` and `confidence` — **measured, not assumed** (V-H quotes the full declaration
from `world/types.ail` with field-anchored counts and a firing control; round 1 stated this
without a row and a reviewer rightly flagged it — its conditional, "amend the v1 schema if the
fields are missing", does NOT trigger); `requestRef` reaches the landed approval machinery
(`ApprovalRequestV1`, approve.go:20), whose `approval_claims` table has single-use enforcement
but **no deadline, no policy, no creation time** (V-C) — the packet supplies exactly those three,
without altering that table. Two reasons, one principled and one measured: composition-by-link is
the P2/P4 idiom (renderers dereference), and embedding `evidence: list[Evidence]` would put a
`list[ADT]` inside the record — the shape measured unencodable by the pinned verifier. In fact
the measurement here is STRICTER than the repo's recorded limitation: a record containing even a
**bare** ADT field fails Z3 encoding (`unknown sort`, V-P2) — see §5 for what that forces.

**Freeze timing — the question §7.3 flags, answered with two options:**

- **Option A (RECOMMENDED): freeze v1 NOW, with this item.** The type lands in `world/types.ail`
  in this item's sprint; amendments after freeze are NEW versions (the reserved semantic ID
  `world/decision-packet/v1` becomes `/v2`), never in-place edits — the same discipline as the
  frozen `LogHeader` and every content-addressed wire type. Reasons: (1) item 7 is charged to
  build "to HUMAN-SURFACE.md", and HUMAN-SURFACE's own header says the §7 outputs "are ratified
  inputs to that work — implement to them"; without a frozen type the inbox will improvise a host
  struct that later fights the kernel type — precisely the coordinator-DB-row pattern §8.1 says
  to replace. (2) §3.1 is already binding and fully determines the minimal field set, so waiting
  buys no information. (3) Content addressing makes an early freeze cheap to amend: the cost of
  being wrong is a version bump, not a migration. (4) The §8 premise row "Decision packets carry
  a ledger-recorded creation time, deadline and typed timeout policy — NOT BUILT" cannot move
  until a type exists.
- **Option B: ratify the policy set now, defer the record freeze to item 7's design.** Stated
  honestly rather than strawmanned: it lets real inbox usage inform the field set. Its cost: the
  §7.1 set cannot be proven *followed at deadline* without the fields the law reads
  (`deadlineAt`, `escalationsRemaining`), so §3.1's replay clause stays untestable until item 7 —
  and item 7 is parked on item 5's chain with no ETA, so the freeze would wait on an unrelated
  MCP blocker.

**Round-2 note on this A/B.** The reviewer who caught the deadline gap opened with "do not
choose Option A **as written**" — and its own fix offers law/schema expansion as the alternative
to deferral, so the DIRECTION (freeze the packet in this item) was not disputed by either
reviewer. The expansion is now done and measured (§4.2, V-P13): the deadline obligation and the
escalation-decrement obligation are proven laws, so the specific defect that made
Option-A-as-written unacceptable — a frozen packet whose lifecycle laws could not see WHEN a
timeout fired — no longer exists. The recommendation REMAINS Option A, on the revised design:
reasons (1)–(4) are untouched, and the revision strengthens (2), since the field set the laws
read (`deadlineAt`, `escalationsRemaining`) is now load-bearing in five proofs rather than
three; Option B's cost grows correspondingly.

The one-word ask to Mark, when this reaches ratification: **§7.1 = the three-constructor set with
§2.1's resolution semantics (yes/amend); §7.3 = Option A or Option B.**

## 3. The central tension resolved — "ledger-recorded" without a clock in the ledger

§3.1 demands a ledger-recorded creation time; the world log has no time column (V-A); the packet
becomes a world type. The resolution is to design WITH the ratified logical-time idiom that
already exists in this repo **twice**, rather than inventing a third mechanism or adding a time
column:

- `host/store/journal.go:26-27`, verbatim: *"LogicalTime is supplied by the caller; journal
  payloads never read a wall clock."* `LogicalTime int64` sits INSIDE the content-addressed
  payload of `JournalIntent`, `JournalOutcome`, `EffectIntent` and `EffectOutcome` (V-B).
- `host/broker/decide.go:15-18,101`: `Capability.ExpiresAt int64` with
  `capabilityLive(c, now) = now >= 0 && now < c.ExpiresAt` and the refusal label
  `denied:expired`; `host/broker/approve.go:369-371`: `PublishApprovalScope.ExpiresAt`, *"the
  LAST logical time at which the approval may be used"* (V-B). `validatePublishApproval` — the
  function at approve.go:485, its ordering block at :594-616 under the comment *"Logical time,
  in both directions."* — already enforces ORDERING RELATIONS between recorded logical times in
  both directions: four explicit comparisons, each with its named error, quoted in V-I. (Round 1
  cited the block's line range as if it were the function's, and cited the enforcement claim
  without a row; both are corrected — the claim itself holds.)

**So "ledger-recorded creation time" means: `createdAt` is a caller-supplied logical instant
carried INSIDE the content-addressed packet payload, which the ledger references** — exactly as
`JournalIntent.LogicalTime` is ledger-recorded today. The world log stays time-free; the frozen
six-field `LogHeader` is untouched; `host/store/schema.sql` is untouched. Deadlines are logical
comparisons (`now >= deadlineAt`), the same shape as `capabilityLive`.

**Replay determinism, precisely — corrected in round 2.** The Timeout transition, when the host
emits it, RECORDS the logical `now` at which it fired inside its own content-addressed payload.
Replay re-applies recorded transitions; it never re-derives "did the deadline pass?" from a
clock. Round 1 stopped at one auditor identity —
`recordedOutcome == timeoutOutcome(packet.policy, packet.escalationsRemaining, recordedAuthorityFact)`
— and claimed §3.1's replay clause was satisfied "by construction". **That was overstated, and
the round-1 reviewer who rejected it was right**: the identity is insensitive to WHEN the
timeout fired. An outcome emitted BEFORE the deadline satisfies it; so does an escalation that
never decrements its budget. The revised design closes both with two further proven laws
(§4.2), so the auditor obligation set on a recorded Timeout is now:

1. `recordedOutcome == timeoutOutcome(policy, escalationsRemaining, recordedAuthorityFact)` —
   the outcome is the one the policy dictates (round 1's law, unchanged);
2. `timeoutFiredLegally(deadlineAt, recordedNow)` — the timeout fired at or after its deadline,
   so an early timeout is a detectable lie, not a compatible history;
3. on every recorded rebind (a DEFER or an `EscalateBounded` escalation):
   `validEscalation(oldBudget, newBudget, recordedNow, newDeadlineAt)` — the budget decremented
   by exactly one and the new deadline is strictly future. (`validDefer` remains the GUARD —
   "may a defer happen now" — the reviewer's catch was precisely that a guard alone never forces
   the decrement on the recorded pair.)

All three are deterministic functions of ledger data; no wall-clock ordering exists anywhere in
the semantics. **What replay still cannot check, named rather than implied** (§5 has the full
coverage map): that `recordedAuthorityFact` is HONEST — it is host-derived (from
`capabilityLive` at the recorded time) and no pure law can verify a bool input's provenance;
and that a timeout which SHOULD have fired DID fire — omission is a liveness property of the
item-7 emitter, invisible to any per-record audit. Both are host obligations, listed as
residuals in §5, not covered by the proofs.

**What this item explicitly DEFERS, and to where** (stated visibly, per the item's charge):

| Deferred machinery | Belongs to | Why not here |
|---|---|---|
| The host emitter — the sweep that observes `now >= deadlineAt` on a pending packet and appends the explicit Timeout transition | item 7's implementation (the inbox is the first surface that materializes packets) | it is effectful host machinery over a surface that does not exist yet; freezing the law it must obey is this item |
| Escalation delivery/notification (P5 batching) | item 7 + the future notifier | same |
| DEFER's "records the required evidence" clause — the transition-boundary check that a defer carries a non-empty evidence ref | item 7's transition wiring | a pure kernel law cannot validate a `HashRef` referent (the item-13 boundary finding, inherited); `validDefer` governs the REBOUND bounds, which are what "bounded" means |
| Enforcement of a monotonic logical-time tick source | host follow-on (with the emitter) | today every recorded `Now` is caller-supplied by design (journal.go:27, decide.go:22-23); this item freezes what the values MEAN, not who ticks them |
| `context.Context` plumbing / bounded store reads | **item 18 — explicitly not this item** | separately queued; scope discipline |

## 4. Proposed change

### 4.1 Placement: extend `world/types.ail`; no new module, no package

`TimeoutPolicy`, `TimeoutOutcome`, `DecisionPacket` and the five laws are core type vocabulary —
the packet "becomes a world type — kernel-adjacent" is §7.3's own words. Item 13 set the
precedent: `world/types.ail` already owns the `Evidence`/`EvidenceGrade` semantic surface with a
proven eliminator co-located. Extending it keeps `LEG1_MODULES` (pin 1) and the whole package
module inventory (four modules, six tar entries) unchanged; a new `world/decisions.ail` would
move pin 1, the manifest-gate module count, and every package allowlist for zero semantic gain.

**Why is this not a package?** The same answer that carried item 13 through quorum: every surface
must share ONE interpretation of what silence means at a deadline. A package could version-skew
the timeout law — two surfaces disagreeing about whether silence approved something is the
cardinal §3.1 violation. The policy set is a closed kernel law, not a domain policy. Domain-level
DEFAULTS (what deadline lengths, how many escalations a given goal class gets) ARE package/host
material and are not frozen here — only the vocabulary and the resolution law are kernel.

### 4.2 Exact AILANG code (validated verbatim on an isolated copy of the real tree — V-P7 round 1, V-P13 the full revised set)

Appended to `world/types.ail` (existing declarations byte-untouched):

```ailang
-- The typed finite set of timeout policies (HUMAN-SURFACE §3.1 / §7 point 1).
-- PROPOSED pending ratification; the three members are the round-2 reviewers'
-- candidates. Nullary by design: the escalation bound is packet state.
export type TimeoutPolicy
  = Cancel
  | EscalateBounded
  | ExecuteIfGranted

-- The typed outcome of the explicit Timeout transition at deadline. Silence
-- never synthesizes approval: no outcome constructor approves anything, and
-- ExecuteUnderPriorAuthority requires authority that already existed.
export type TimeoutOutcome
  = ResolveRejectedTimeout
  | EscalateWithNewDeadline
  | ExecuteUnderPriorAuthority

-- The frozen v1 decision packet (HUMAN-SURFACE §7 point 3). References, not
-- embeds: proposalHash reaches evidence/caps; requestRef reaches the landed
-- approval objects. createdAt/deadlineAt are LOGICAL times (journal.go:26-27
-- idiom); no field reads a wall clock. Amendments are /v2, never in-place.
export type DecisionPacket = {
  packetHash: HashRef,
  proposalHash: HashRef,
  requestRef: HashRef,
  createdAt: int,
  deadlineAt: int,
  escalationsRemaining: int,
  policy: TimeoutPolicy
}

-- The total timeout law: what the explicit Timeout transition must do, per
-- policy, given the packet's escalation budget and the host-derived fact of
-- whether independent authority is live at the recorded logical time.
export func timeoutOutcome(policy: TimeoutPolicy, escalationsRemaining: int, independentAuthority: bool) -> TimeoutOutcome ! {}
ensures { result == match policy {
  Cancel => ResolveRejectedTimeout,
  EscalateBounded => if escalationsRemaining > 0 then EscalateWithNewDeadline else ResolveRejectedTimeout,
  ExecuteIfGranted => if independentAuthority then ExecuteUnderPriorAuthority else ResolveRejectedTimeout
} }
{
  match policy {
    Cancel => ResolveRejectedTimeout,
    EscalateBounded => if escalationsRemaining > 0 then EscalateWithNewDeadline else ResolveRejectedTimeout,
    ExecuteIfGranted => if independentAuthority then ExecuteUnderPriorAuthority else ResolveRejectedTimeout
  }
}

-- DEFER rebound law: a defer is valid only with a strictly-future new deadline
-- and remaining escalation budget. DEFER MUST NOT park indefinitely.
export func validDefer(now: int, newDeadline: int, escalationsRemaining: int) -> bool ! {}
ensures { result == (newDeadline > now && escalationsRemaining > 0) }
{
  newDeadline > now && escalationsRemaining > 0
}

-- Packet schedule well-formedness: creation time is a valid logical instant,
-- the deadline is strictly after creation, and the budget is non-negative.
export func wellFormedSchedule(createdAt: int, deadlineAt: int, escalationsRemaining: int) -> bool ! {}
ensures { result == (createdAt >= 0 && deadlineAt > createdAt && escalationsRemaining >= 0) }
{
  createdAt >= 0 && deadlineAt > createdAt && escalationsRemaining >= 0
}

-- Timeout deadline law: an explicit Timeout transition is legal only at or
-- after the packet's deadline, judged against the RECORDED logical now.
-- This is what connects the deadline to the outcome: an early timeout is a
-- detectable lie under replay, not a compatible history.
export func timeoutFiredLegally(deadlineAt: int, recordedNow: int) -> bool ! {}
ensures { result == (recordedNow >= deadlineAt) }
{
  recordedNow >= deadlineAt
}

-- Escalation/DEFER rebind law: any recorded rebind must decrement the budget
-- by EXACTLY one and set a strictly-future new deadline. validDefer guards
-- "may a defer happen now"; this law validates the recorded before/after
-- pair, so a non-decrementing escalation is a detectable lie under replay.
export func validEscalation(oldEscalationsRemaining: int, newEscalationsRemaining: int, recordedNow: int, newDeadlineAt: int) -> bool ! {}
ensures { result == (oldEscalationsRemaining > 0 && newEscalationsRemaining == oldEscalationsRemaining - 1 && newDeadlineAt > recordedNow) }
{
  oldEscalationsRemaining > 0 && newEscalationsRemaining == oldEscalationsRemaining - 1 && newDeadlineAt > recordedNow
}
```

**Why two laws and not one.** The obvious single-law form — fold the outcome check into one
`validTimeout(…, outcome: TimeoutOutcome) -> bool` conjunction — does NOT typecheck on the
pinned binary: `No instance for Eq[TimeoutOutcome] in scope`, and the suggested workaround
`import std/prelude` is itself unavailable (`IMP012_UNSUPPORTED_NAMESPACE`, V-P12). ADT equality
inside a general boolean expression needs an `Eq` instance; the top-level `result == match …`
postcondition form does not — which is exactly why `timeoutOutcome` verifies and the combined
law cannot. The split is a measured constraint, not a stylistic choice (§5, limitation 3).

Plus five private test adapters (`outcomeCode`, `deferCode`, `scheduleCode`, `firedCode`,
`escalationCode`) whose exact bodies are in the V-P7/V-P13 probes: single-`int` case
dispatchers, because two v0.30.0 harness limitations were MEASURED, not assumed —
multi-argument `tests [...]` rows fail to parse (V-P5), and tuple-valued test inputs are
collected but fail at runtime with "no pattern matched" in BOTH tuple-pattern styles (V-P5).
Nineteen named identities result: `outcomeCode_test_1..6` (one per semantic branch of the
timeout law), `deferCode_test_1..3`, `scheduleCode_test_1..3`, `firedCode_test_1..3` (early
fire refused / fires exactly at deadline / fires after), `escalationCode_test_1..4` (valid
rebind / zero budget / non-decrementing budget / non-future new deadline).

The laws take bare parameters, not the `DecisionPacket` record, because that is the shape the
pinned verifier can prove (§5) — the round-1 reviewer's proposed `validTimeout(packet, …)`
signature was built and run verbatim, and it is Z3-unencodable with `check.passed` still `true`
and rc 0 (V-P10), the silent failure mode this whole document is built around. The record
itself is declared but carries no contract; its declaration does not poison the module (V-P13:
6 identities verified, `errors=0`, with the record present).

### 4.3 Gate pins moved in the same commit

1. Pin 1 — `scripts/verify_ail.sh:135` `LEG1_MODULES`: **unchanged** (no module added); listed
   because the conflict surface must show it was considered, not because it moves.
2. Pin 2 — `scripts/verify_ail.sh:310` `EXACT_TOTAL_VERIFIED`: 5 → **10** (re-measured for
   round 2 on the isolated tree: `world/types.ail` reports 6 verified — `gradeOf` + the five
   laws — atop the 4 in other modules; V-P13).
3. Pin 3 — `scripts/verify_ail.sh:333` `REQUIRED_TESTS` (python set): add the nineteen
   identities.
4. Pin 4 — `scripts/verify_ail.sh:342` `EXACT_TOTAL_TESTS` (python, spaces around `=`): 20 →
   **39** (observed on the isolated tree, V-P13 — NOT transcribed from round 1's 32; observed
   output authoritative at implementation).
5. Pin 5 — `host/verifygate/module_manifest_gate_test.go:128` marker string: `"✓ 5/5 required
   world/ identities verified across 11 module(s)"` → `"✓ 10/10 … across 11 module(s)"` (module
   count stays 11; only the identity count moves). The `:208` comment region's "11 modules"
   sentence stays true.
6. Also in `verify_ail.sh` (~:262-267): `REQUIRED_VERIFIED["world/types.ail"]` moves from
   `{"gradeOf"}` to `{"gradeOf", "timeoutOutcome", "timeoutFiredLegally", "validEscalation",
   "validDefer", "wellFormedSchedule"}`.

`EXACT_TOTAL_MODULES` does not exist (V-E; grep 0, control 4) — the 11 is an allowlist
cardinality. Package projection: run `./scripts/build_world_package.sh`, require byte equality of
canonical and projected `types.ail`, regenerate `scripts/world_package_ready_packet.golden.json`
through `host/pkgproj` (never hand-edited). Four modules and six tar entries are unchanged;
content/interface/tarball hashes all move.

## 5. What the proof proves — and the one thing it structurally cannot

The exact contracts prove: `timeoutOutcome` is TOTAL over the three-constructor policy set and
returns exactly the ratified resolution for every (policy, budget, authority) combination;
`timeoutFiredLegally` is exactly the deadline comparison against the recorded logical now;
`validEscalation`, `validDefer` and `wellFormedSchedule` are exactly their stated conjunctions
(the escalation law's decrement clause is the exact `oldEscalationsRemaining - 1`, not an
inequality). Non-vacuity is measured in both directions and for both rounds' laws: a
deliberately false postcondition yields a real counterexample (V-P4), a body-only mutation of
the NEW deadline law draws `counterexample=1` (V-P15), and a missing match arm with the
contract present is a verifier ERROR even though v0.30.0 ACCEPTS non-exhaustive ADT matches at
check time with rc 0 (V-P9) — the contract is the only totality guard, which is why every
future `TimeoutPolicy` constructor added without an arm reds Leg 1 rather than silently falling
through.

**The measured blind spot, and its named killers.** Z3 cannot see a consistent lie: mutating the
`ExecuteIfGranted`/no-authority arm to `ExecuteUnderPriorAuthority` in BOTH contract and body —
i.e., making silence synthesize execution, the exact §3.1 cardinal violation — still prints
`verified` (V-P8, shas `1321b29d…` → `2c031275…`; re-landed on the revised round-2 tree in
V-P14, sha `fda7f30b…` → `dca1eead…`, same result). The mutations were LANDED and measured, not
predicted: on every consistent lie in V-P14 — early timeout, non-decrementing escalation,
non-future rebind deadline, synthesized execution — Z3 stayed green (`verified=6, cex=0,
errors=0`) and exactly one named test failed. Proof plus tests, never proof alone; every
semantic branch and every conjunction clause has a runtime case for this reason.

**The round-2 coverage map — which lifecycle violation each law forecloses, and what stays
host-owned.** The round-1 reviewer named four ways a recorded history could lie while remaining
"compatible with all three proofs". With the revised law set:

| Violation | Foreclosed by | Measured killer (V-P14) |
|---|---|---|
| Early timeout (outcome recorded before the deadline) | `timeoutFiredLegally` — pure, proven | `firedCode_test_1`, sole failure; Z3 green on the lie |
| Non-decrementing escalation | `validEscalation` (exact `- 1`) | `escalationCode_test_3`, sole failure |
| Missing / non-future new deadline | `validEscalation` (`newDeadlineAt > recordedNow`) | `escalationCode_test_4`, sole failure |
| Fabricated authority | **NO PURE LAW — declared host residual** | none can exist: `independentAuthority` is a host-derived bool INPUT, and a pure law cannot verify its provenance. The obligation sits with the item-7 emitter — derive the fact from the LANDED `capabilityLive` check (decide.go:101) at the recorded logical time, and record the capability ref alongside the outcome so an auditor can re-derive it. Stated here so no reader believes the proofs cover it. |

**Three verifier limitations measured here, routed UPSTREAM (rule: no local workarounds):**

1. **A record containing a bare ADT field is Z3-unencodable** — `unknown sort 'Packet'`,
   `verify.errors=1`, while `check.passed` stays `true` and rc stays 0 (V-P2; flat-int-record
   control verifies, V-P3). This is STRICTER than the repo's recorded limitation ("the failing
   shape is a RECORD containing `list[ADT]`", iter-79): the list is not needed; the ADT field
   alone fails. The iter-79 lesson recurs at finer grain — a recorded limitation is a claim at
   the granularity that sufficed for the case that produced it. Consequence for this design: the
   laws take bare parameters, and no contract is attempted on `DecisionPacket` itself.
   Implementation files this as a `sunholo-data/ailang` issue + `ailang messages send
   mission-control` note.
2. **The inline-test harness cannot execute tuple-valued inputs**: identities are collected, then
   every case fails "no pattern matched", in both nested-match and direct tuple-pattern forms
   (V-P5). Same routing. The `caseId` dispatcher is the honest in-repo shape, not a workaround of
   the language (it uses only ratified constructs).
3. **ADT equality in a general boolean expression requires an `Eq` instance the stdlib cannot
   supply on v0.30.0** — `outcome == match policy { … }` inside a conjunction fails typecheck
   with `No instance for Eq[TimeoutOutcome] in scope`, while the top-level `result == match …`
   postcondition form verifies; and the error's own suggested workaround, `import std/prelude`,
   fails with `IMP012_UNSUPPORTED_NAMESPACE` (V-P12). This is what forces §4.2's three-law
   split instead of one combined `validTimeout` law. Same routing (measured round 2).

The proofs do NOT establish: that any `HashRef` in a packet resolves (inherited item-13
boundary limitation); that the host emitter fires at the right moment, or at all — omission is
emitter liveness (deferred, §3); that `independentAuthority` is honestly derived (the declared
host residual in the coverage map above); or that logical time is globally monotonic (deferred,
§3). None of these is silently claimed.

## 6. Mutation table

Every refusal/resolution branch has a neutering mutation; deterministic oracles only. "Both"
means contract AND body (the consistent-lie form); "body" means body only. Each mutation lands
with pre/post sha256 and is restored byte-identical — a mutation that never ran and a mutation
that did not red share an exit code (the item-13 V26 lesson).

| ID | Mutation | Expected RED (named observable) | Control |
|---|---|---|---|
| MU1 | `Cancel` arm → `EscalateWithNewDeadline` (body) | `verify.counterexample >= 1` on `timeoutOutcome` | restore → `verified` |
| MU2 | `ExecuteIfGranted`/no-auth arm → `ExecuteUnderPriorAuthority` (**both** — silence synthesizes execution) | Z3 stays GREEN (measured, V-P8/V-P14); `outcomeCode_test_5` fails, alone | restore → 39 named tests pass |
| MU3 | `EscalateBounded` exhaustion arm → `EscalateWithNewDeadline` (**both** — unbounded park) | `outcomeCode_test_3` fails | `outcomeCode_test_2` stays green |
| MU4 | `EscalateBounded` live arm → `ResolveRejectedTimeout` (**both**) | `outcomeCode_test_2` fails | `outcomeCode_test_3` stays green |
| MU5 | Delete one arm from `timeoutOutcome` body (contract kept) | `verify.errors == 1`, status `error` (V-P9 shape); `check.passed` stays `true` — rc is NOT the oracle | restore → `verified`, `errors=0` |
| MU6 | Add a 4th `TimeoutPolicy` constructor with no arm | `verify.errors == 1` (totality guard) | add explicit arm → `verified` |
| MU7 | `validDefer`: drop `newDeadline > now` (**both**) | `deferCode_test_2` fails | `deferCode_test_3` stays green |
| MU8 | `validDefer`: drop `escalationsRemaining > 0` (**both**) | `deferCode_test_3` fails | `deferCode_test_2` stays green |
| MU9 | `wellFormedSchedule`: `deadlineAt > createdAt` → `>=` (**both**) | `scheduleCode_test_2` fails | `scheduleCode_test_1/3` stay green |
| MU10 | Remove `timeoutOutcome` from `REQUIRED_VERIFIED` | AC3's manifest inspection fails; the named identity leaves Leg-1 output | restore literal |
| MU11 | Edit canonical `world/types.ail` only (skip projection) | package gate SHA-equality step fails | rebuild projection |
| MU12 | Rebuild projection, keep old golden | package gate golden byte-compare fails | recompute via `pkgproj` |
| MU13 | Leave pin 5's marker at `5/5` | `TestModuleManifestGate` reds the CI go-verify job | update marker |
| MU14 | `timeoutFiredLegally`: accept an early fire (**both** — the round-1 reviewer's "early timeout") | Z3 stays GREEN (measured, V-P14); `firedCode_test_1` fails, alone | restore → 39 named tests pass |
| MU15 | `validEscalation`: exact decrement → `<=` (**both** — non-decrementing escalation) | Z3 stays GREEN (measured, V-P14); `escalationCode_test_3` fails, alone | `escalationCode_test_1/2/4` stay green |
| MU16 | `validEscalation`: `newDeadlineAt >` → `>=` (**both** — missing/non-future new deadline) | Z3 stays GREEN (measured, V-P14); `escalationCode_test_4` fails, alone | `escalationCode_test_1/2/3` stay green |
| — | Fabricated authority (the round-1 reviewer's fourth mutation) | **NO mutation exists that a pure gate can red** — `independentAuthority` is a host-derived input (§5's coverage map); recorded as a declared item-7 host residual, not a pretended kill | n/a |

## 7. Acceptance criteria

Each AC names its observable, states how it fails AND how it passes, and the observable is
downstream of the mechanism it claims. None requires producing a value the instrument cannot
produce.

1. **Policy-set shape.** `TimeoutPolicy` has exactly the three constructors of §2.1 (or exactly
   the amended set Mark ratifies, if amended). Observable: the type declaration plus MU6's
   guard. Fails on any extra/missing/renamed constructor relative to the ratified text.
2. **Packet-schema shape.** `DecisionPacket` is exactly §4.2's seven fields; no `list[…]` of any
   ADT enters the record; existing exported types are byte-untouched. Fails on any field
   add/drop/retype relative to the ratified schema.
3. **Named proofs.** `REQUIRED_VERIFIED["world/types.ail"]` gains exactly
   `timeoutOutcome`, `timeoutFiredLegally`, `validEscalation`, `validDefer`,
   `wellFormedSchedule`; each reports `verified` with `verify.errors=0`, `counterexample=0`;
   `EXACT_TOTAL_VERIFIED` moves 5 → 10 (secondary guard). Fails if any identity is missing from
   `verify.results[]` or non-`verified`.
4. **Named runtime pins.** The nineteen identities of §4.2 are added to `REQUIRED_TESTS`;
   `EXACT_TOTAL_TESTS` moves 20 → 39, subject to observed runner JSON at implementation
   (`len(tests[])`, never `passed_tests`, which also counts contract-derived properties —
   observed 43 vs 39 on the round-2 probe tree, V-P13). Fails on a missing/failing identity or
   total drift.
5. **The consistent-lie discriminators.** Land MU2, MU14, MU15 and MU16, each with pre/post
   sha256; require Z3 to stay green on every one AND the named test to be the SOLE failure
   (`outcomeCode_test_5`, `firedCode_test_1`, `escalationCode_test_3`,
   `escalationCode_test_4` respectively); restore byte-identical; require 39/39 after each
   restore. This AC fails if the tests are ever weakened to the point a §3.1 violation — silence
   synthesizing execution, an early timeout, an undecremented budget, a non-future rebind —
   survives; detecting exactly that is why it exists.
6. **Totality guard RED/control.** Land MU5 and MU6 on an isolated copy; require
   `verify.errors==1` with `check.passed` still `true` in both (rc-blindness documented, never
   used as the verdict); controls restore to `verified`/`errors=0`.
7. **Per-branch neutering.** Every MU1–MU9 and MU14–MU16 lands (sha-verified), reds ONLY its
   named observable, and restores byte-identical. Fails if any mutation survives or reds a
   different observable. The fabricated-authority row is NOT an AC target — it has no pure
   killer by construction and is carried as a named residual, not a gate.
8. **Pin 5.** The manifest-gate marker reads `10/10 … across 11 module(s)` and the CI go-verify
   job is green on the PR head. Fails on a stale marker (this is the omission that red-lighted
   an earlier item's CI).
9. **Projection + golden.** `./scripts/build_world_package.sh` run; canonical/projected
   `types.ail` byte-equal; golden regenerated through `host/pkgproj` (content/interface/tar
   hashes move; four modules and six tar entries unchanged). Fails on hand-authored JSON or
   unchanged hashes.
10. **Full pinned gates.** `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` rc=0
    with the named output; `GOTOOLCHAIN=go1.25.6 ./scripts/verify_go.sh` green (the
    `GOTOOLCHAIN` requirement is a base condition of the rig, not a regression).
11. **Scope.** The diff touches ONLY the five files in the header. Fails if it touches
    `host/store`, `host/daemon`, `host/broker`, `schema.sql`, or adds any host emitter — that
    machinery is deferred per §3's table (emitter → item 7; context plumbing → item 18).
12. **Upstream routing.** The three verifier limitations of §5 are filed as
    `sunholo-data/ailang` issues with an `ailang messages send mission-control` note,
    referencing V-P2/V-P5/V-P12 evidence. Fails if any is silently absorbed as local lore.

Post-ratification acts that are NOT sprint ACs (they are controller/human acts, listed so they
are not lost): flip HUMAN-SURFACE §7.1/§7.3 to CLOSED with the ratified text; update the two §8
premise rows this item moves; reconcile the item-15/item-7 charter rows (§10).

## 8. Conflict surface

### 8.1 The five gate pins (fact E — all five listed, including the unmoved one)

- `scripts/verify_ail.sh:135` `LEG1_MODULES` (bash array compared as a set, own anti-vacuity
  floor at :227, loud diff at :234-236): **unchanged** — no module is added.
- `scripts/verify_ail.sh:262-267` `REQUIRED_VERIFIED`: `world/types.ail` set gains five names.
- `scripts/verify_ail.sh:310` `EXACT_TOTAL_VERIFIED`: 5 → 10.
- `scripts/verify_ail.sh:333` `REQUIRED_TESTS` (PYTHON set): +19 identities.
- `scripts/verify_ail.sh:342` `EXACT_TOTAL_TESTS` (PYTHON, spaces around `=` — shell-shaped greps
  miss it): 20 → 39.
- `host/verifygate/module_manifest_gate_test.go:128` marker string `5/5 … 11 module(s)` →
  `10/10 … 11 module(s)`; `:208` region's 11-module sentence remains true.

### 8.2 Package projection and golden

- `scripts/build_world_package.sh`: four-module allowlist unchanged; regenerates the `types.ail`
  copy wholesale.
- `packages/world-core/world/types.ail`: gains the addition; must stay byte-identical to
  canonical.
- `scripts/verify_world_package.sh`: four modules/exports, six tar entries, SHA equality, golden
  byte-compare — all inventory pins unchanged, all hashes move.
- `scripts/world_package_ready_packet.golden.json`: regenerated, never hand-edited.

### 8.3 Documents that change at RATIFICATION (not in the sprint)

- `design_docs/HUMAN-SURFACE.md` §7 points 1 and 3 (open → closed with ratified text), and two §8
  premise rows ("Decision packets carry…" NOT BUILT → BUILT for the type/laws; the deterministic-
  timeout row stays NOT BUILT until the item-7 emitter — this doc does not overclaim it).
- `design_docs/world-mission.md`: the item-15 row, and the item-7 row per §10.

### 8.4 Adjacent machinery deliberately NOT touched

`host/store/schema.sql` (no time column added — §3), `host/store/journal.go` (idiom reused, not
modified), `host/broker/approve.go`/`decide.go` (the landed approval/capability surfaces are the
packet's link targets, unchanged), `world/contracts.ail`/`transitions.ail`/`logepoch.ail`
(byte-untouched; their required identities and tests are regression guards via pins 2–4).

## 9. Scope and pricing

| Work | Time |
|---|---:|
| `world/types.ail` addition (types + 5 contracts + 5 adapters), pinned-binary green | 0.32 d |
| Gate pins (six edits across two files), MU1–MU16 with shas and controls | 0.40 d |
| Package projection + golden, full pinned gates, three upstream filings | 0.28 d |
| **Total** | **1.00 d** |

Still fits the queued ~1d band, now at its top edge — re-priced, not quietly absorbed. The
round-2 delta is bounded because it was MEASURED before pricing: the two added laws and both
adapters are already written and green on the isolated tree (V-P13), and all three new
mutations already landed, killed and restored there (V-P14) — the sprint increment is gate-pin
arithmetic, re-execution on the real tree, and one more upstream filing, not new design work.
The emitter (multi-day, item-7-shaped) and store plumbing (item 18) remain excluded; if quorum
demands either, the honest answer is a milestone split with re-measurement, not silent growth.

## 10. The item-7 relationship — a measured row-level discrepancy, not inherited

Item 15's row says it "blocks item 7". Item 7's OWN row (charter line ~2266) states its park
chain as `TR.A2 → TR.B → TR.C → item 5 P6.B → item 7` and does not name item 15 (V-G). Both are
ratified charter state and they disagree. What this design actually creates: a **content
prerequisite** — item 7 must construct packets conforming to the frozen schema and policy set
(HUMAN-SURFACE's header: the §7 outputs "are ratified inputs to that work — implement to them")
— but NOT a scheduling gate: landing this item does not unpark item 7 (item 5's chain still
holds), and not landing it would leave item 7 with no ratified packet schema to build to.
Proposed reconciliation at ratification: add item 15 to item 7's prerequisite chain as a content
input, or strike "blocks" from item 15's row in favour of "ratified input to". This doc asserts
neither side as already true.

## 11. Verification Log

Controller-supplied facts are labelled `VERIFIED BY CONTROLLER`; every one was **re-run
first-party by this designer** at `bc8f193` on 2026-08-14 with the same-scope controls shown.
Probe rows V-P1..P9 are this designer's own round-1 work; rows V-H, V-I and V-P10..P15 are the
round-2 revision's, run the same way (controller round-2 measurements re-run first-party before
being relied on, including the numbers the controller itself supplied). All probes ran against
the pinned `/tmp/ailang-v0300/ailang` (v0.30.0, `e37b370`); probe files lived under `/tmp/`
(`/tmp/probe-item15/`, `/tmp/iso-item15/`, `/tmp/iso-item15-r2/`), never in the repo. Semantic
verdicts come from JSON fields, never process rc.

| ID | Claim | Exact command | Observed output |
|---|---|---|---|
| V-A | `VERIFIED BY CONTROLLER` — the world ledger records no time; `log_entries` carries exactly the six frozen header fields + entry hash + transition ref. | `grep -ciE "timestamp\|created_at\|wall_clock" host/store/schema.sql; grep -cinE '\btime\b' host/store/schema.sql; grep -c "CREATE TABLE" host/store/schema.sql`; `sed -n '36,45p'` | `0`; `0`; control `8` (file read); log_entries columns as §1 states. |
| V-B | `VERIFIED BY CONTROLLER` — the logical-time idiom exists twice: journal payloads carry caller-supplied `LogicalTime` inside content-addressed bytes; broker expiry is logical (`capabilityLive`), approvals carry a logical `ExpiresAt`. | `sed -n '26,27p' host/store/journal.go; grep -n "LogicalTime" host/store/journal.go; grep -n "ExpiresAt\|capabilityLive" host/broker/decide.go; sed -n '367,371p' host/broker/approve.go` | journal.go:26-27 verbatim as quoted; `LogicalTime` on all four payload types; decide.go:18/55/101; approve.go:369 "LAST logical time". |
| V-C | `VERIFIED BY CONTROLLER` — approval machinery exists and is narrower than §3.1: `approval_claims` has no deadline, policy, or creation time. | `sed -n '92,96p' host/store/schema.sql` | exactly `approval_ref PK, request_ref, invocation_id UNIQUE`. |
| V-D | `VERIFIED BY CONTROLLER` — the lifecycle is unimplemented: no packet type, no escalation, no decision deadline. | `grep -rilE "decisionpacket\|decision_packet" world/ host/ \| wc -l` (+ controls: `Proposal` in `world/`, `func ` in `host/`); `grep -rilE "escalat" … \| grep -v _test \| wc -l` (control `Approval`); `grep -rn "DecisionDeadline\|decisionDeadline\|deadlineAt\|DeadlineAt" --include='*.go' . \| grep -v _test \| wc -l` (control `ApprovalRequestV1  =`) | `0` / controls `6`, `79`; `0` / control `24`; `0` / control `1`. All instruments fire. |
| V-E | `VERIFIED BY CONTROLLER` — the five pins are as fact E states; `EXACT_TOTAL_MODULES` does not exist; `world/` holds exactly four modules. | `sed -n '132,150p;225,240p;255,300p;305,345p' scripts/verify_ail.sh; sed -n '125,132p;205,212p' host/verifygate/module_manifest_gate_test.go; ls world/; grep -c "EXACT_TOTAL_MODULES" scripts/verify_ail.sh; grep -c "EXACT_TOTAL_VERIFIED" scripts/verify_ail.sh` | 11-entry `LEG1_MODULES` at :135 with empty-allowlist floor; `REQUIRED_VERIFIED` dict at ~:262 (`types.ail: {"gradeOf"}`); `EXACT_TOTAL_VERIFIED=5` at :310; `REQUIRED_TESTS` at :333; `EXACT_TOTAL_TESTS = 20` at :342; marker at gate test :128; `contracts.ail logepoch.ail transitions.ail types.ail`; `0` / control `4`. |
| V-P1 | The three-constructor policy law with an exact match contract (including if-arms) verifies non-vacuously. | `/tmp/probe-item15/policy.ail`; `AILANG_RELAX_MODULES=1 /tmp/ailang-v0300/ailang ai-check -timeout 5s policy.ail` → JSON | `check true`; `verified=1, counterexample=0, errors=0`; `timeoutOutcome verified`. |
| V-P2 | **A record containing a BARE ADT field is Z3-unencodable, silently at rc level** — stricter than the repo's recorded `list[ADT]` limitation. | `packetrec.ail` (record with `policy: Policy` field, contract on the record param); same ai-check | `check true`; `errors=1`; reason: `Z3 error … unknown sort 'Packet'`. |
| V-P3 | Same-shape control: the flat int record verifies. | `packetint.ail` (identical contract, no ADT field) | `verified=1, errors=0`, `wellFormedFlat verified`. |
| V-P4 | Non-vacuity control: a false postcondition yields a real counterexample. | `policyfalse.ail` (`ensures result == ResolveRejectedTimeout` over the honest body) | `counterexample=1`, status `counterexample`. |
| V-P5 | v0.30.0 inline tests cannot take multiple arguments (parse fail) NOR tuple inputs (collected, then runtime fail) — both tuple-pattern styles. | `policytests.ail` (4-tuple rows) → `ailang test --format json`; `tupletests.ail` / `tupletests2.ail` (tuple input, nested + direct patterns) | multi-arg: single `parse fail` test; tuple: `policyCase_test_1/2` collected, both `fail`, message `no pattern matched in match expression`. |
| V-P6 | The single-int `caseId` dispatcher executes all six branches AND coexists with the verified contract in one file. | `casetests.ail` → `ailang test --format json` + `ai-check` | `outcomeCode_test_1..6` all `pass` (6/0); same file `timeoutOutcome verified, errors=0`. |
| V-P7 | **Round 1's §4.2 addition (three laws), applied to an isolated copy of the REAL `world/` tree, is green end-to-end; baseline control reproduces the live pins.** Historical record — the full FIVE-law set is validated in V-P13, whose numbers govern the pins. | `cp world/*.ail /tmp/iso-item15/world/`; sha before `2cf5b004…` (= the exact file item 13 landed); append §4.2 + adapters → sha `1321b29d…`; `ai-check world/types.ail`; `ailang test --format json world/`; pristine-copy baseline in `/tmp/iso-base` | baseline: `len(tests)=20, failed=0` (reproduces `EXACT_TOTAL_TESTS=20`). Post-addition: `verified=4, errors=0` (`timeoutOutcome`, `validDefer`, `wellFormedSchedule`, `gradeOf`); `len(tests)=32, failed=0`; the twelve round-1 identities (`outcomeCode_test_1..6`, `deferCode_test_1..3`, `scheduleCode_test_1..3`) exactly as named. NOTE: `passed_tests=34 ≠ 32` — it counts contract-derived properties, reaffirming the gate's `len(tests[])` choice. The uncontracted `DecisionPacket` record does NOT poison the module. |
| V-P8 | **The consistent-lie mutation (silence synthesizes execution, contract AND body) is Z3-GREEN and killed by exactly one named test.** Landed, not predicted. | On the V-P7 tree: python replace asserting exactly 2 occurrences; sha `1321b29d…` → `2c031275…`; `ai-check` + `test`; restore from backup; sha re-check | ai-check: `verified=4, errors=0, cex=0` — `timeoutOutcome` says `verified` ON THE LIE. Tests: sole failure `outcomeCode_test_5`, `failed_tests=1`. Restored byte-identical (`1321b29d…`). |
| V-P9 | A missing match arm WITH the contract present is a verifier error while `check.passed` stays true and rc stays 0 — the contract is the only totality guard. | `policy_mut.ail` (EscalateBounded arm deleted from body) → ai-check JSON | `check true`; `verify.errors=1`; `timeoutOutcome status=error`. |
| V-G | `VERIFIED BY CONTROLLER` — the item-15/item-7 row discrepancy. | `grep -n "w-decision-lifecycle-freeze" design_docs/world-mission.md` + `sed` of both rows (item 15 at ~:2911, item 7 at ~:2266) | item 15's row: "blocks item 7, gated on nothing"; item 7's row: park chain `TR.A2 → TR.B → TR.C → item 5 P6.B → this item`, no mention of item 15. Both quoted in §10. |
| V-H | **Round 2** — `Proposal` carries `evidence`, `requiredCaps` and `confidence` (§2.2's references-not-embeds premise; refutes the round-1 objection's doubt, confirms its procedure). | `awk '/export type Proposal/,/^}/' world/types.ail`; `grep -c 'evidence:' world/types.ail` (+ `requiredCaps:`, `confidence:`; control `proposalHash:`) | Full nine-field declaration verbatim, including `evidence: list[Evidence]`, `requiredCaps: list[Capability]`, `confidence: float`; counts `2`/`1`/`1`; control `3` — instrument fires. |
| V-I | **Round 2** — `validatePublishApproval` enforces ordering relations between recorded logical times in both directions (§3's premise; corrects round 1's citation). | `grep -n 'func validatePublishApproval' host/broker/approve.go`; `sed -n '590,620p'`; `grep -n 'Logical time, in both directions'` | Function at `:485`; comment at `:594`; four explicit comparisons: `request.Now > req.Now` → `ErrPublishApprovalMalformed`; `decision.Now` outside `[request, publish]` → malformed; `scope.ExpiresAt < request.Now` → malformed ("expiry precedes its own request time"); `req.Now > scope.ExpiresAt` → `ErrPublishApprovalExpired` (at `:612-615`). |
| V-P10 | **Round 2** — the round-1 reviewer's proposed fix, taken LITERALLY (`validTimeout(packet, …)` with `packet: DecisionPacket`), is Z3-unencodable and fails SILENTLY. | `/tmp/probe-item15/r2fix_verbatim.ail` → `AILANG_RELAX_MODULES=1 /tmp/ailang-v0300/ailang ai-check -timeout 5s` | `check.passed=true`, rc 0; `verify.verified=0, errors=1`; `validTimeout -> error`: `Z3 error … unknown sort 'DecisionPacket'`. The remedy as written exhibits the exact failure mode this document is built around. |
| V-P11 | **Round 2** — the objection is satisfiable in int-only form: `timeoutOutcome` + `timeoutFiredLegally` + `validEscalation` verify together, no ADT equality, no record parameter. | `/tmp/probe-item15/r2fix_split.ail` → same ai-check | `check.passed=true`; `verified=3, cex=0, errors=0`; all three `verified`. |
| V-P12 | **Round 2** — the single-law conjunction form fails typecheck (`Eq` instance), and the suggested `std/prelude` workaround is unavailable — the measured constraint forcing the law split. | `/tmp/probe-item15/r2fix_bare.ail`; `/tmp/probe-item15/r2fix_prelude.ail` → same ai-check | bare: `check.passed=false`, `No instance for Eq[TimeoutOutcome] in scope. Equality (==, !=) needs an Eq instance` at the `outcome ==` site; prelude: `IMP012_UNSUPPORTED_NAMESPACE … namespace imports not yet supported`. |
| V-P13 | **Round 2** — the FULL revised §4.2 (five laws, five adapters), applied to a fresh isolated copy of the real `world/` tree, is green end-to-end; baseline reproduces the live pins. | `cp world/*.ail /tmp/iso-item15-r2/world/` (baseline `types.ail` sha `2cf5b004…`, identical to round 1's); baseline `ailang test --format json world/`; apply round-1 addition (`1321b29d…`) + two laws + two adapters → sha `fda7f30b…`; `ai-check world/types.ail`; `ailang test --format json world/` | Baseline: `len(tests)=20, failed=0`. Post-addition: `verified=6, cex=0, errors=0` (`gradeOf` + all five laws); `len(tests)=39, failed=0`; the seven new identities exactly `firedCode_test_1..3`, `escalationCode_test_1..4`. NOTE `passed_tests=43 ≠ 39` — contract-derived properties again; `len(tests[])` remains the oracle. |
| V-P14 | **Round 2** — all four consistent-lie mutations (the reviewer's early timeout, non-decrementing escalation, missing new deadline; plus the V-P8 re-land) are Z3-GREEN and each is killed by exactly one named test. Landed and restored, not predicted. | On the V-P13 tree, each via python replace asserting exactly 2 occurrences, sha'd, then `ai-check` + `test`, then restored from backup with sha re-check: early (`fda7f30b…` → `2930abeb…`); non-decrement (→ `f155f9c2…`); stale deadline (→ `3102a461…`); V-P8 re-land (→ `dca1eead…`) | Every mutation: `verified=6, cex=0, errors=0` — Z3 verified THE LIE. Sole test failures respectively: `firedCode_test_1`; `escalationCode_test_3`; `escalationCode_test_4`; `outcomeCode_test_5`. All four restores byte-identical (`fda7f30b…`). |
| V-P15 | **Round 2** — the new deadline contract is non-vacuous: a body-only mutation (contract kept honest) draws a real counterexample. | On the V-P13 tree: replace ONLY the body occurrence of `recordedNow >= deadlineAt` with `>` (sha `b8062d9e…`); `ai-check`; restore | `verified=5, counterexample=1, errors=0`; `timeoutFiredLegally -> counterexample`. Restored byte-identical (`fda7f30b…`). |

## 12. What this item is NOT doing

- It does not implement the host Timeout-transition emitter, the escalation notifier, or any
  sweep — deferred to item 7 with the packet's first materialized surface (§3).
- It does not add a time column, touch `host/store/schema.sql`, the frozen `LogHeader`, or the
  journal codecs.
- It does not do `context.Context` plumbing or bounded store reads — that is item 18.
- It does not enforce a monotonic logical-time tick source (deferred with the emitter, §3).
- It does not validate any `HashRef` referent (inherited item-13 boundary limitation, stated).
- It does not — and purely CANNOT — verify the provenance of the `independentAuthority` fact;
  that is the declared host residual of §5's coverage map, owed by the item-7 emitter.
- It does not decide domain defaults (deadline lengths, escalation budgets per goal class).
- It does not modify `Evidence`, `EvidenceGrade`, `gradeOf`, or any existing exported type.
- It does not add a module, package, tar member, Go API, effect, or renderer surface.
- It does not close HUMAN-SURFACE §7 points 1 or 3 — it puts a decidable proposal in front of
  Mark; only ratification closes them.

## Related

- `design_docs/HUMAN-SURFACE.md` — §3.1 (the binding principle implemented), §7 points 1 and 3
  (the two ratification asks), §8 premise rows this item moves.
- `design_docs/coding-standards.md` — S1 (contracts-first, named identities), S2 (the emitter is
  host work), S3 (why-not-a-package answered in §4.1), S5 (fluency protocol followed: reference
  loaded via the pinned binary's `ailang prompt`; every syntax claim probe-verified), S6 (honest
  gates; rc never the oracle).
- `design_docs/DESIGN.md` §1, §14 — kernel boundary; a lifecycle change after freeze is a
  version bump through propose→verify→commit, never an in-place edit.
- `design_docs/implemented/w-evidence-grade-mapping.md` — the co-location precedent, the
  proof-plus-tests discipline, and the `HashRef`-validation boundary this item inherits.
- `world/types.ail` — the module extended; `host/store/journal.go`, `host/broker/decide.go`,
  `host/broker/approve.go` — the ratified logical-time idiom designed with.
- `scripts/verify_ail.sh`, `host/verifygate/module_manifest_gate_test.go`,
  `scripts/build_world_package.sh`, `scripts/verify_world_package.sh` — the five pins and the
  projection gates.
- Upstream: three `sunholo-data/ailang` issues to file at implementation
  (record-with-bare-ADT-field Z3 sort error silent at rc level; tuple-valued inline-test inputs
  fail at runtime; `Eq`-instance requirement for ADT equality in general boolean expressions
  with the `std/prelude` namespace import unsupported), per §5.
