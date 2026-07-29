# w-effect-journal — Closing the dispatch→record window at effect granularity (the EFFECT journal)

**Status**: Planned — **QUORUM-CLEARED via the charter's narrow-refinement carve-out** (NEW-DOC,
queue row 4c, designer rotation iter-36). **r0** authored by `claude:claude-fable-5`; **r1** fixed
quorum round-1's blocking finding (resumption-ordinal collision); **r2** applied quorum round-2's
blocking fix VERBATIM as a bounded controller revision (the split allocator became one
transactional `AppendNextEffectIntent`) plus the round-2 pass-vote's concrete `GetReceipt` guard.
See the Quorum verification log at the end.
**Date**: 2026-07-29
**Charter clause**: clause-3 (explicit authority end-to-end — the last open window of objection 2A)
**Verified against**: **`AILANG v0.30.0`** — the pinned released binary at `/tmp/ailang-v0300/ailang`
(`AILANG v0.30.0`, commit `e37b370`, built 2026-07-19, no `-dirty` — `--version` run first-party
this session, V1). This doc authors **no new `.ail` at design time**; its one sketch change (MJ.A)
extends the already-landed, already-verified `design_docs/sketches/storejournal.ail`, whose
CURRENT baseline was measured first-party this session (V16): **7/7 contracts `verified`** by name
(`intentBindsCommit`, `writableRefText`, `isIndeterminate`, `mayReportNotStarted`, `retryAllowed`,
`journalSeqNext`, `outcomeMatchesIntent`), **`len(tests[]) == 30`**, `passed_tests == 37`
(reported separately; the gate is on `len(tests[])` only). The sprint re-measures after the edit
and records the real numbers — never quotes these.
**Traces to**: [DESIGN.md](../DESIGN.md) §1 ("a pure transition plus recorded effect results
reconstructs any historical state"), §7–§9 (effects / evidence), §14 (boundaries); charter clause 3
and the M3.D ratification (`c26b27d`, Mark, attended, 2026-07-28) — where this item's `host/store`
kernel reopen is **pre-ratified IN PRINCIPLE, and the design is NOT pre-approved**
**Depends on**: [w-effect-broker-m3.md](../implemented/w-effect-broker-m3.md) (**LANDED**, all
five milestones — the broker pipeline this item rewires and the option-(i) recovery path this item
extends) and [w-store-durability.md](../implemented/w-store-durability.md) (the commit journal
whose laws, receipt discipline, and DDL this item reuses without touching)
**Inherits**: **CF-N-2** (`maxRecoveryPages` unjustified magic number) and **CF-N-3**
(`retryAllowed(false, true)` untested) — both folded in as acceptance criteria, not footnotes
**Estimated**: **~1.5d** (top of the queue row's 1–1.5d band — honest ONLY because Path A below
avoids the DDL; Path B would honestly cost ~2.5–3d, see the correction note)

> **Scope note — the honest claim this item will be allowed to make at close.** Today, if the
> process dies between a handler's dispatch and its record write, an external effect (an
> `FS.Write`, a `Git.Commit`, a paid `Model.Infer`) HAS HAPPENED with no durable trace, and
> replay cannot distinguish "never executed" from "executed, record lost". M3's option (i) made
> every **commit** crash-detectable and said, verbatim, that the effect-granularity window
> stays open. This item closes it **at the granularity objection 2A named**: before any handler
> is dispatched, a durable per-effect intent exists; after this item, an effect with no durable
> intent was NEVER dispatched, and an effect whose intent has no outcome is INDETERMINATE and
> fail-closed — never silently absent, never auto-re-executed. **One residual stays open and is
> stated, not hidden** (Decision 5): a crash between the effect *record* write and the *outcome*
> append leaves an indeterminate receipt whose record already exists — fail-closed, resolvable by
> deterministic reconciliation, and strictly better than today, but "indeterminate" there means
> "executed and recorded, bookkeeping incomplete", not "unknown". No claim in this doc may state
> or imply that every crash ambiguity is eliminated.

---

## Correction note — THE QUEUE ROW'S COSTING CLAIM IS REFUTED (verified first-party, V4–V6)

The charter's queue row and the M3 doc both carry the M3 planner's costing of option (iii):
*"Cheaper than assumed — no schema change (the commit shape lives in the payload codec, not the
DDL, so AC13's `sqlite_master` gate stays green)."* The iter-36 controller reported this refuted;
per this mission's method I **re-ran all three probes myself** rather than citing the controller,
and the refutation stands:

- **The premise is true.** The commit shape (eight refs) lives in the payload codec
  (`host/store/journal.go:104-143`), not the DDL (V2).
- **The conclusion does not follow**, because the journal table's *kind vocabulary* and
  *cardinality* live in the DDL (`host/store/schema.sql:81-87`, V3):
  - **P1** — `INSERT … kind='effect-intent'` → `CHECK constraint failed: kind IN
    ('intent','outcome')` (V4). A new kind label is DDL-rejected.
  - **P2** — a second `kind='intent'` row for the same `invocation_id` →
    `UNIQUE constraint failed: journal.invocation_id, journal.kind` (V5). N effects need N
    **distinct** invocation IDs, or a DDL change.
  - **P3 — the sharp one** — a widened `CHECK` re-applied to an EXISTING store returns **rc=0,
    no error**, and `sqlite_master` shows the journal DDL **unchanged**; the new-kind insert then
    fails on the OLD constraint (V6). Cause: every statement in `schema.sql` is `CREATE TABLE IF
    NOT EXISTS`, a no-op on an existing table, and `store.Open` applies the schema by one
    `db.Exec(schemaSQL)` (V13). Combined with **V7** — zero migration machinery of any kind in
    `host/` (no `user_version`, no `ALTER TABLE`, no migrate, no schema version) — **any DDL
    change ships fail-OPEN**: new stores get the new schema, every existing store silently keeps
    the old one, and nothing detects the disagreement. A schema change that was never applied is
    indistinguishable from one that was — the iter-35 mutation lesson, at the schema layer.
- **V8** — the one DDL-drift gate, `TestPreJournalMigrationPreservesExistingDDL`
  (`host/store/journal_test.go:340`), `delete(after, "journal")`s the journal table out of its
  comparison and never exercises the upgrade-an-existing-store path — so it would not catch P3
  either.

**AND THE GATE THE CLAIM CITED HAS NO TEETH WHERE IT WAS CITED — measured by the controller with
three compiling mutations, not read (V21–V23).** The costing claim's safety net was "AC13's
`sqlite_master` gate stays green". It does stay green — but that says nothing, because:

- **V21** — `MUT-JOURNAL-DDL-WIDEN` (form: widen the `kind` CHECK in `host/store/schema.sql`,
  pattern asserted to match exactly once, diff printed, `go build`/`go vet` rc=0) leaves
  **`go test ./...` rc=0 across all 10 packages, zero FAIL**. *Nothing in this repository guards
  the journal table's DDL.*
- **V22** — the gate's OWN named mutation, `MUT-DDL-DRIFT` (`log_entries` CHECK), does red — but
  **read the message, not the exit code** (the iter-34 rule): it reds with `pre-journal schema
  source drifted: sha256=7358d876…`, i.e. the **sha256 source pin on `journal_test.go:344`**,
  which `t.Fatalf`s before a database is ever built. The `sqlite_master` comparison never runs.
- **V23** — decisive: `MUT-DDL-COMPARE-DEAD` (form: `if !reflect.DeepEqual(after, before)` →
  `if false && !reflect.DeepEqual(after, before)`; both variables stay used so the mutant must
  still compile — `go vet` rc=0) leaves the **whole `host/store` package green**. The
  `sqlite_master` byte-identity comparison contributes **zero discrimination today**.

The gate's real teeth are a source-text sha256 pin over the pre-journal *prefix* of `schema.sql`.
That genuinely protects pre-existing tables — the protection is real, just delivered by a
different mechanism than the one it is documented as having. It does **not** cover the journal
table (edits there are below the pinned prefix, V21) and does **not** cover the upgrade path
(V6). This is the **fifth instance** of this mission's signature shape — *a gate no production
change could fail* — and the first found at PICK time, in a gate inherited from an item that is
already COMPLETE and was judged with zero blocking findings. Per the charter's iter-32 process
fix, inherited gates get audited too; this is what that audit found.

So "no schema change" is not a discount that comes for free; it is a **design constraint this
item must actively satisfy**. The design below satisfies it (Path A), which is the only reason
the ~1.5d estimate survives. **And no premise in this doc may cite the `sqlite_master` gate as
evidence that the journal DDL is protected** — V21–V23 are why.

## Motivation

The broker's pipeline today (M3 Decision 3): decide → debit → **dispatch handler** → put result
object → **write effect record**. Everything durable happens AFTER the external world has already
been touched. M3.D's option (i) anchored the *episode commit* (intent once world+entry are built),
which makes the commit crash-detectable — but between one effect's dispatch and its record, the
only witness is process memory. `host/broker/recover.go:55-59` states this open window in its own
doc comment and names this queue item as its owner.

The fix is the one the M3 planner scoped as option (iii): **effect-shaped intents in the store's
journal** — a durable "I am about to dispatch effect X" written BEFORE the handler runs, and a
durable outcome after the record exists. The commit journal cannot carry them (its intent shape
requires six commit refs that are structurally unknowable pre-dispatch, V2/V12); this item adds an
effect-shaped payload beside it, in the same table, under the same laws.

**The non-negotiable inherited from the ratification record (`c26b27d`), restated verbatim in
force:** never synthesize placeholder or dummy refs to satisfy `validateIntent`. Every field of
the effect intent below is TRUE at write time; nothing is invented to pass validation.

## Premises (hard constraints — each verified in the Premise Verification Log)

- **P1 — `host/store/schema.sql` is byte-for-byte unchanged.** Path A's entire point, and the
  premise is enforced by `git diff --exit-code` on the file (AC1) — **not** by the
  `sqlite_master` gate, which V21–V23 measure as having no teeth on the journal table. What
  byte-identity buys, decisively: **no migration story is needed**, because V6/V7 prove this repo
  cannot ship one safely today. (`TestPreJournalMigrationPreservesExistingDDL`'s sha256-pinned
  pre-journal prefix `43a9c80b…` also stays untouched — stated as a fact, not leaned on as a
  guarantee.)
- **P2 — the kernel reopen is exactly scoped, and every touched method is named.** M3 Design
  Freeze bullet 8 said "zero `host/store` method changes"; this item is the ratified exception.
  The full delta is in Decision 2: **five new methods, two changed methods, zero removed, zero
  changed signatures on existing methods.** Anything beyond that list is out of scope for the
  sprint without returning here.
- **P3 — invocation-ID namespaces are disjoint by construction, enforced in BOTH appenders.**
  Effect IDs live in a reserved namespace (`effect:` prefix, Decision 3); the commit appender
  rejects it, the effect appender requires it. V14: no existing caller, test or production, uses
  the reserved prefix, so the commit-side narrowing breaks nobody.
- **P4 — the effect journal never lies.** An effect intent asserts only pre-dispatch truths
  (what, where, at what cost, from which episode, in which order). No commit-shaped fields, no
  placeholder refs, no fields whose value is a guess.
- **P5 — recovery reports; it never acts.** The landed no-dispatch / no-auto-resolve /
  fail-closed discipline of `broker.Recover` extends unchanged to effect granularity. The
  `retryAllowed` law (sketch LAW 3) governs re-execution; automatic re-execution of an
  indeterminate effect remains structurally forbidden.
- **P6 — every wait and allocation is bounded.** The effect journal adds **no new subprocess or
  exec surface** (so the `Setpgid`/`Kill(-pid)` rule is not newly triggered — stated, not
  skipped); the recovery walk stays bounded, and CF-N-2's bound becomes justified and tested
  rather than magic.
- **P7 — the gates extend without edits.** New Go code is swept by `verify_go.sh`/CI; the sketch
  edit is to an ALREADY-swept module, so `verify_ail.sh`'s module count stays 11 and the 4/4
  identity + 14 world-test totals cannot be perturbed (the required manifests key on `world/`
  modules only — the M3 doc's V16, mechanism unchanged since). The bench manifest gains **zero
  new names** (Decision 6); existing rows are re-measured because the pipeline's cost changes.

## High-Impact Decisions

| Decision | Why High Impact | Chosen By | Change Cost |
|----------|-----------------|-----------|-------------|
| **Path A**: DDL-free — per-effect synthetic invocation IDs + effect-shaped payload codecs, reusing `kind='intent'/'outcome'` | The only path that closes the window WITHOUT shipping a fail-open schema change (V6/V7); keeps the estimate inside the band | this doc (quorum at pick) | high |
| Reserved ID namespace `effect:<episodeID>:<ordinal>`, enforced in both appenders | ID-space disjointness is what lets one table carry two payload shapes without ambiguity; unenforced, a collision corrupts both journals' receipt laws | this doc | medium |
| **The ordinal is minted by the STORE inside the appending transaction (`AppendNextEffectIntent`) — the broker never holds one** | Two successive quorum findings, both accepted: r1 (gemini-3-1-pro) — a transient counter restarting at 0 after a crash GUARANTEES `DuplicateInvocationError` against the pre-crash `effect:<episodeID>:0` row, bricking the episode; r2 (gpt5-6-sol) — deriving it durably but appending in a SEPARATE operation just moves the collision to a TOCTOU. Only an in-transaction mint closes both | quorum r1 + r2 | high |
| Intent is durable BEFORE dispatch; outcome AFTER the record | The entire point of the item — the ordering IS the window closure; frozen as a law in the sketch | this doc | high |
| Denied effects write NO journal rows | No dispatch → no window → nothing to journal; the denial record (landed M3 discipline) is already durable and complete; journal growth stays 2 rows per *dispatched* effect | this doc | low |
| `PendingIntents` gains a `semantic_id` filter; `PendingEffectIntents` is a separate method | Without the filter, commit recovery decodes effect payloads as commit intents and errors on every store with a pending effect (V9: `PendingIntents` decodes every row via `decodeJournalIntent`) | this doc | medium |
| CF-N-2 and CF-N-3 are acceptance criteria with their own mutations | Inherited carry-forwards die by evidence, not by prose | charter | low |

### Design Freeze (the sprint must not renegotiate these)

- [x] `host/store/schema.sql` byte-for-byte unchanged; **no DDL change, no migration machinery**.
- [x] The `host/store` method delta is EXACTLY Decision 2's table (**r2**):
  `AppendNextEffectIntent`, `AppendEffectOutcome`, `GetEffectReceipt`, `PendingEffectIntents`
  (new — **four**, since r2 folded the ordinal derivation INTO the appender's transaction and
  removed both `AppendEffectIntent` and `NextEffectOrdinal`); `validateIntent`, `PendingIntents`
  and `GetReceipt` (changed, as specified); nothing else.
- [x] The ordering law: **no handler dispatch without a prior durable effect intent**; the
  outcome append happens after the effect record object is durable; both directions tested and
  mutation-pinned.
- [x] Effect invocation IDs are `effect:<episodeInvocationID>:<ordinal>` (ordinal decimal ≥ 0);
  the commit appender REJECTS the `effect:` prefix; the effect appender REQUIRES the full shape.
  No other ID form enters the effect journal. **The ordinal is minted by the STORE, inside the
  same transaction that appends the intent (`AppendNextEffectIntent`) — never by the broker, never
  from an in-memory counter, never read outside a transaction.** There is no ordinal cache
  (Decision 3; quorum r1 found the fresh-zero collision, quorum r2 found that a split
  read-then-append merely moved it to a TOCTOU).
- [x] No placeholder/dummy refs, ever — every effect-intent field is true at write time
  (ratification record `c26b27d`).
- [x] Recovery never dispatches, never auto-resolves, never appends an outcome, never
  re-executes — extended verbatim to effect granularity.
- [x] The broker's record/result object shapes, the three-arm `recordConsistent` law, Replay
  mode, and `host/replay/**` are byte-untouched.
- [x] Denied effects journal nothing; succeeded AND failed effects journal intent + outcome.

## Decision 1 — Path A over Path B (the F3-informed choice)

**Path A (CHOSEN) — DDL-free.** One synthetic invocation ID per effect, derived from the episode's
invocation ID and the effect's ordinal position, reusing the existing `kind='intent'/'outcome'`
vocabulary and the existing `UNIQUE (invocation_id, kind)` cardinality. V5's P2b probe confirms
the mechanism live: with today's DDL, `('intent','inv-1')` and `('intent','inv-1/effect/0')`
coexist. What Path A costs, stated plainly:

- **ID-space semantics**: the invocation-ID namespace becomes structured. The commit appender's
  accepted domain NARROWS (it now rejects `effect:`-prefixed IDs) — a compatibility break in
  principle, a no-op in fact (V14: zero existing callers use the prefix; the only production
  journal writers are the M3.D episode driver pattern and tests, enumerated in V14).
- **Collision question, answered**: can an effect ID collide with a commit ID? Only if a commit
  caller chooses an `effect:`-prefixed ID — which the commit appender now structurally rejects,
  in both `AppendIntent` and (via the same `validateIntent`) anything built on it. Within the
  effect namespace, two effects collide only at identical `(episodeID, ordinal)`. **r0 called
  that "the broker's own sequencing bug if it happens" — the quorum showed that framing hid a
  GUARANTEED case: a resumed broker whose transient counter restarts at 0 collides with the
  pre-crash `effect:<episodeID>:0` row on its first new dispatch.** Decision 3's durable
  derivation removes that case by construction; the surviving collision surface really is a
  sequencing bug, and it still fails loudly as `DuplicateInvocationError` on differing bytes,
  the landed discipline.
- **Two journal rows per dispatched effect** — journal growth is linear in dispatched effects
  (denied effects add zero rows). The seq-keyset paging and `MaxPendingIntentsPage` bound (V9)
  absorb this without change.

**Path B (REJECTED) — new kind labels (or a new table).** Requires the DDL change that P1 rejects
and P3 punishes: per V6, the edit is **silently dropped on every existing store**, and per V7
there is no migration machinery to make it otherwise. Taking Path B honestly means designing and
shipping a schema-versioning + fail-LOUD-on-mismatch mechanism first (a `user_version` pin
checked at `Open`, an upgrade step, tests for the un-upgraded-store path) — a ratification-class
kernel change of its own, ~1–1.5d on top, for **no semantic gain over Path A**: the payload
codec, not the kind label, is what distinguishes effect rows from commit rows either way. Path B
becomes right only if a future item needs kinds the codec cannot express; it would then inherit
the migration design as its first milestone. Recorded here as the rejected alternative with that
reason.

**Why is this not a package? (S3)** For the same reason the commit journal was not: this is the
durability floor BENEATH the layer packages live in. A package's effects are what the journal
makes crash-honest; the journal cannot itself be subject to propose→verify→commit. The kernel
delta is deliberately minimal (P2) and everything above it — reconciliation policies, journal
compaction, per-handler reconcilers — is package-lane deferred scope.

## Decision 2 — The kernel reopen, exactly scoped

New payload codecs (payload-codec change, exactly as the queue row's true premise says —
`world/effect-intent/v1`, `world/effect-outcome/v1`, built by the landed `journalObject` /
canonical-JSON machinery, golden-bytes-tested like the commit codecs):

```go
// All fields are TRUE before dispatch. No commit-shaped refs exist here — that
// is the point, and the reason validateIntent's six-ref rule does not apply.
type EffectIntent struct {
    InvocationID string          // "effect:<episodeID>:<ordinal>" — the synthetic ID itself
    EpisodeID    string          // the episode's (commit-journal) invocation ID
    Ordinal      int64           // 0-based dispatch order within the episode
    Effect       string          // e.g. "FS.Write" — the decided request's effect name
    Scope        string          // the decided request's scope
    Cost         int64           // the decided request's cost (the debit that will stand)
    RequestRef   hashref.HashRef // content address of the request payload bytes, put first
    LogicalTime  int64           // caller-supplied; journal payloads never read a wall clock
}

type EffectOutcome struct {
    InvocationID string
    Status       string          // "succeeded" | "failed" — denied never journals
    RecordRef    hashref.HashRef // the effect record object (which carries everything else)
    LogicalTime  int64
}
```

| Method | New/Changed | Exact change |
|--------|-------------|--------------|
| `AppendNextEffectIntent(episodeID string, intentWithoutID EffectIntent) (id string, ordinal int64, err error)` | **new** (**r2**, quorum round-2 blocking fix — the reviewer's own signature, adopted verbatim) | **ONE transaction does all of it**: acquire the store's write serialization → derive the next ordinal from ALL matching durable intents (resolved, indeterminate and pending alike — the r1 candidate-set, now computed INSIDE the tx) → validate exhaustion (`ordinal == math.MaxInt64` → structured `OrdinalExhaustedError`, never a silent wrap) → construct the canonical ID `effect:<episodeID>:<ordinal>` and payload → insert the object + journal row → commit. The ordinal never exists outside the transaction, so no read→allocate→append interleaving is possible. Also validates non-empty `Effect`/`Scope`/`episodeID` and a non-zero `RequestRef`, and keeps `AppendIntent`'s duplicate discipline for the constructed row |
| ~~`AppendEffectIntent(id, EffectIntent)`~~ | **REMOVED in r2** | superseded — a caller-supplied ordinal is exactly the split allocation the round-2 objection names. The broker no longer chooses effect IDs; the store mints them |
| ~~`NextEffectOrdinal(episodeID)`~~ | **REMOVED in r2** | superseded — its derivation moved INSIDE `AppendNextEffectIntent`'s transaction. Exposing it as a standalone read is what made the TOCTOU expressible |
| `AppendEffectOutcome(id, EffectOutcome)` | **new** | mirrors `AppendOutcome`: requires a durable effect intent, rejects a second outcome, requires `Status ∈ {succeeded, failed}` and non-zero `RecordRef` |
| `GetEffectReceipt(id)` | **new** | the same three-state receipt law (`not-started` / `indeterminate` / `resolved`) over the effect payloads; rejects non-`effect:` IDs with a structured error |
| `PendingEffectIntents(limit, fromIndex...)` | **new** | `PendingIntents`'s keyset paging, filtered to `o.semantic_id = 'world/effect-intent/v1'`, same `MaxPendingIntentsPage` bound |
| `validateIntent` | **changed** | adds: reject `InvocationID` with the reserved `effect:` prefix (structured `InvocationMismatchError`) — the commit-side half of P3's disjointness |
| `PendingIntents` | **changed** | query gains `AND o.semantic_id = 'world/journal-intent/v1'` — behaviour-preserving on every existing store (no effect intents exist anywhere yet, V14), and REQUIRED once they do (V9: today's query decodes every pending row as a commit intent and would error on the first effect payload) |

`GetReceipt` gains ONE guard (**r2** — `gemini-3-1-pro`'s round-2 fix, adopted verbatim although
it voted pass): `strings.HasPrefix(id, "effect:")` → structured `InvocationMismatchError`, matching
`validateIntent`'s strictness. Its reason is concrete — without it, a caller mistakenly passing an
`effect:` ID reaches `decodeJournalIntent` on an effect payload and surfaces a codec failure
instead of a clean boundary rejection. This promotes OD1 from a sprint-time MAY to a MUST; OD1 is
struck from Open Decisions accordingly.

## Decision 3 — The ID namespace

`effect:<episodeInvocationID>:<ordinal>`. Prefix-based (not infix) so the check is one
`strings.HasPrefix`; the episode ID is embedded verbatim so recovery findings group by episode
with zero extra lookups. The derivation lives in ONE exported function
(`store.EffectInvocationID(episodeID string, ordinal int64)`) used by both the broker and the
tests — never re-derived inline (the mirror-drift lesson). Its canonical text form gets named
golden tests (S1's canonical-text rule); promotion of the derivation into the `.ail` sketch is
deliberately NOT promised here because it needs string stdlib verification under the fluency
protocol at sprint time (Open Decision OD2).

**The ordinal (r1 — this replaces r0's refuted claim).** r0 froze *"the ordinal is the broker's
per-episode dispatch counter, making the derivation deterministic and collision-free by
construction within an episode"* — **FALSE across a crash boundary**, as the round-1 quorum
(gemini-3-1-pro) correctly found: Replay is byte-untouched and Recovery only reports, so nothing
re-initialized the transient counter, and a resumed broker's first dispatch at ordinal 0 is a
guaranteed `DuplicateInvocationError` (differing bytes) against the pre-crash intent —
permanently bricking the episode. The correction:

- **r1's correction (superseded in form by r2, kept because its REASON still governs the
  candidate set):** the ordinal must come from durable state, never a fresh in-memory zero.
  A fresh episode gets 0 because the journal holds nothing for it; a resumed episode gets
  `1 + max(journaled ordinal)`, past every ID the pre-crash process claimed, including any
  indeterminate one. One scan per episode start, bounded by journal size (P6); whether it uses
  the `UNIQUE (invocation_id, kind)` index via a prefix range is an implementation detail the
  sprint may take or leave — the V27 probe verified the semantics, not the plan.
- **r2 (quorum round 2, `gpt5-6-sol`) — THE SPLIT ALLOCATOR IS GONE; there is no in-memory
  counter at all.** r1 kept a per-process cache incremented per dispatch, with
  `NextEffectOrdinal` and `AppendEffectIntent` as two separate operations. The reviewer's
  objection, accepted: *"`NextEffectOrdinal(episodeID)` and `AppendEffectIntent` occur in separate
  operations, so two episode starts or broker instances can both observe the same maximum and mint
  the same effect ID … the resumption fix merely replaces the restart collision with a TOCTOU
  collision."* Its fix is applied **verbatim**: one transactional store operation,
  `AppendNextEffectIntent`, which acquires the store's write serialization, derives the next
  ordinal from all matching durable intents, validates exhaustion, constructs the canonical ID and
  payload, and inserts the object plus journal row **before releasing the transaction**. The
  broker-maintained ordinal cache is REMOVED (the reviewer's condition for keeping it — separately
  enforced and verified exclusive per-episode ownership — is not something this item establishes,
  so the cache goes). The ordinal is therefore never held outside a transaction, and the
  read→allocate→append sequence cannot interleave.
- **What the landed single-writer lock does and does not cover — stated, and NOT used to soften
  the fix (V28).** One limb of the objection is refuted by landed code: two broker *processes*
  cannot share a store as writers, because `store.Open` takes a non-waiting exclusive lock and the
  second caller gets `*WriterAlreadyActive` (proven cross-process by the landed A1 test). The other
  limb is REAL and unguarded: `host/store/store.go` and `journal.go` contain **no mutex of any
  kind**, so two goroutines in ONE process — two concurrent episodes — can interleave freely. The
  reviewer's fix is adopted in full regardless: a transactional allocator is correct under both
  limbs, and narrowing a fix to the part of an objection one can refute is how a real defect
  survives a review.
- **Alternative considered — the reviewer's own proposal**, deriving from restored episode state
  (`int64(len(episode.History))` or records counted during Replay): **rejected, because it has a
  residual collision in exactly the crash this item exists to close.** History counts *records*,
  and an indeterminate effect is an intent whose record was lost — so after the AC7 crash,
  `len(History)` returns the indeterminate intent's own ordinal and the first new dispatch
  collides with it. The journal sees every claimed ordinal, including indeterminate ones; the
  history does not. The substance of the reviewer's fix (durable derivation, never a fresh zero)
  is adopted; the source is the journal itself.

## Decision 4 — The rewired broker pipeline (allowed arms only)

```
episode start (fresh OR resumed) → NOTHING to initialize: the broker holds no ordinal (r2, Decision 3)
decide → (denied? → record denial, return — UNCHANGED, no journal)
       → PutObject(request payload)                     [durable, true]
       → AppendNextEffectIntent(episodeID, intent)      [durable, ordinal minted IN-TX — THE WINDOW CLOSES HERE]
       → debit → dispatch handler
       → put result object (success) / none (failed)    [landed M3 discipline, unchanged]
       → write effect record (three-arm, unchanged)
       → AppendEffectOutcome(status, recordRef)         [receipt → resolved]
       → return result/error + record ref
```

Crash truth-table after this item (the receipt law doing its job per effect):

| Crash point | Durable state | Receipt | Meaning |
|---|---|---|---|
| before intent | nothing | `not-started` | **never dispatched** — and that answer is now TRUE by the ordering law |
| intent ↔ outcome | intent | `indeterminate` | dispatch may have happened; fail-closed, surfaced by recovery, never re-executed without the `retryAllowed` gate |
| record ↔ outcome | intent + record | `indeterminate` | executed AND recorded; bookkeeping incomplete — the stated residual (Decision 5) |
| after outcome | intent + outcome | `resolved` | fully accounted |
| **resumption (after ANY of the above)** | every journaled intent for the episode, resolved or not | next ordinal derived and appended in ONE tx by `AppendNextEffectIntent` | **no collision is possible, and none is expressible**: every ID the pre-crash process claimed is a journal row, so every journaled ordinal ≤ max and the next mint is strictly fresh — and because the derivation never leaves the transaction, two concurrent allocations cannot observe the same max. A fresh-zero counter was the r1 quorum-refuted defect; a split derive-then-append was the **r2** one |

The debit ordering (intent before debit) means a crash between intent and debit loses only
in-memory ledger state — which the record stream already reconstructs (landed M3 Decision 3), and
an indeterminate intent tells recovery the debit's status is exactly as unknown as the dispatch.

## Decision 5 — Recovery at effect granularity + the stated residual

`broker.Recover` additionally pages `PendingEffectIntents` (same bounded walk, same cursor
discipline) and surfaces one finding per indeterminate effect —
`IndeterminateEffectError{InvocationID, EpisodeID, Ordinal, Effect, Scope}` (the existing type
gains the effect-shaped fields; its two commit-planned fields stay for commit findings) — never
dispatching, never resolving, never appending. The `retryAllowed` law is unchanged and now has a
per-effect subject, which is what it was always specified over (`storejournal.ail` LAW 3).

**The residual, stated:** an indeterminate effect whose `RecordRef`-bearing outcome was lost but
whose record object WAS written is resolvable by deterministic reconciliation (the record exists,
content-addressed; a reconciler may append the outcome it proves). That reconciler is
**deferred scope** — it is an act, and this item's recovery only reports. Until it lands,
such effects stay indeterminate: fail-closed, honest, strictly more knowledge than today.

## Decision 6 — CF-N-2, CF-N-3, and the bench posture

- **CF-N-2** (`maxRecoveryPages = 1 << 20`, `host/broker/recover.go:10`, V10): the constant gets
  (a) a justification comment deriving the bound (`maxRecoveryPages × MaxPendingIntentsPage` ≈
  10⁹ intents — orders of magnitude above any store this substrate can produce, while keeping
  Standing-rule-6's guarantee that a corrupted cursor cannot loop forever), and (b) a TEST of the
  exceeded-bound error path via a never-draining fake store (the landed `recoveryStore` interface
  makes this injectable without production changes), pinned by `MUT-RECOVERY-UNBOUNDED`.
- **CF-N-3** (`retryAllowed(false, true)` untested, V11 — three sites: the law's `tests[]` in
  `storejournal.ail` (rows cover only `(F,F),(T,F),(T,T)`), the Go application at
  `host/broker/recover.go:46`, and the test-local mirror at `host/store/recover_test.go:26`): the
  `(false, true)` row is added at ALL THREE sites, and its non-vacuity is proven by
  `MUT-RETRY-XOR` (below) — a mutant that ONLY that row catches (truth-table in the Non-Vacuity
  section).
- **Bench**: `scripts/bench_worldd.sh`'s hardcoded manifest gains **zero names**. The existing
  `BenchmarkBrokerFSRead` row now measures the full pipeline INCLUDING the two journal appends —
  the effect journal's price lands inside the row that already exists. All `bench/BASELINE.md`
  rows are re-measured in ONE invocation (the landed discipline) and the delta on the FS.Read row
  is reported as this item's cost, visibly (A9), with a short note in `BASELINE.md`.

## Milestones (each independently CI-green and mergeable; ~1.5d total)

### MJ.A — The kernel reopen + the law (~0.5d)

- **files**: `host/store/journal.go` (+~255 — the two codecs, four new methods incl.
  `AppendNextEffectIntent` (in-tx mint + exhaustion check), `validateIntent` namespace check,
  the `GetReceipt` namespace guard, `PendingIntents` filter,
  `EffectInvocationID`), `host/store/journal_test.go`
  (+~290 — golden bytes ×2, disjointness both directions, duplicate/orphan-outcome discipline,
  receipt three-state walk, `PendingIntents` cross-contamination, `PendingEffectIntents` paging,
  in-tx ordinal minting over fresh/resolved/indeterminate/adversarial-suffix populations, the
  two-goroutine concurrent-allocation test and the `MaxInt64` exhaustion test — AC7b),
  `design_docs/sketches/storejournal.ail` (+~15 — **LAW 5** `effectDispatchLawful(dispatched,
  intentDurable) -> bool`, total-theorem `ensures result == (dispatched == false ||
  intentDurable)` with positive AND negative `tests[]` rows, plus CF-N-3's `((false, true),
  true)` row on LAW 3 — both edits are shapes this sketch already proves; **verified live before
  landing**: `AILANG_BIN=/tmp/ailang-v0300/ailang` `ai-check` with contract identities asserted
  from `verify.results[]` (expect 8 named, was 7 — V16), and `test --format json` with
  `len(tests[])` and `passed_tests` recorded as two separate numbers, gated on `len(tests[])`
  (baseline 30); no assertion on `skipped_tests`, per premise V14 of the hardening doc)
- **acceptance**: AC1–AC6, AC8; `schema.sql` byte-unchanged by `git diff --exit-code`
- **verify**: `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh` ·
  `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` (11 modules, 4/4 identities,
  14 world tests unperturbed)
- **ci_green_boundary**: nothing consumes the new methods yet; the broker is untouched

### MJ.B — The rewired pipeline + the crash-window proof + recovery (~0.5d)

- **files**: `host/broker/broker.go` (+~95 — intent-before-dispatch via `AppendNextEffectIntent`,
  which returns the minted ID; the broker keeps NO ordinal state at all (r2, Decision 3);
  outcome-after-record; denied arm untouched), `host/broker/recover.go` (+~60 —
  the effect-granularity page walk + widened finding type), `host/broker/broker_test.go` /
  `recover_test.go` (+~300 — the **crash-window test**: fault-injected store loses everything
  after the handler runs but before record/outcome; reopen; `GetEffectReceipt` =
  `indeterminate`; `Recover` surfaces it with the counting-probe registry asserted at ZERO
  dispatches; **then — the r1 quorum-mandated arm — the resumed broker dispatches a NEW effect
  on the same episode and it SUCCEEDS: fresh ID past the indeterminate ordinal, no
  `DuplicateInvocationError`, intent + outcome journaled**; plus ordering tests,
  denied-no-journal, failed-arm-journals-outcome),
  `host/broker/episode_test.go` (+~40 — the landed episode now shows per-effect receipts
  resolved end-to-end alongside its commit anchor)
- **acceptance**: AC7 (the window closure itself), AC9, AC10
- **ci_green_boundary**: recovery still reports-only; replay untouched

### MJ.C — CF-N-2/CF-N-3 discharge + bench re-measure + close-out (~0.25–0.5d)

- **files**: `host/broker/recover.go` (CF-N-2 comment), `host/broker/recover_test.go`
  (never-draining-fake bound test; the `(false,true)` rows), `host/store/recover_test.go`
  (mirror row), `bench/BASELINE.md` (full re-measure, one invocation, delta note), `README.md`
  (+~5), doc → `implemented/` with every box checked **only after CI is observed green**
- **acceptance**: AC11, AC12, AC13; the full Non-Vacuity table demonstrated RED per the iter-35
  instrument discipline (below)
- **verify**: both gates + `./scripts/bench_worldd.sh --smoke` + `go test ./...` zero skips

**Pre-committed overflow cut line**: if the honest total exceeds ~1.5d, the cut is **the
episode_test extension and the BASELINE delta note → follow-up (~0.25d)**; the kernel reopen, the
ordering law, the crash-window proof, and both carry-forward discharges are NOT cuttable — they
are the item.

## Files to Create/Modify (aggregate)

| File | Est. LOC | Change |
|------|---------:|--------|
| `host/store/journal.go` | +~255 | codecs, 4 new methods (incl. `AppendNextEffectIntent`, r2), 3 changed methods |
| `host/store/journal_test.go` | +~290 | new-surface tests |
| `host/store/recover_test.go` | +~5 | CF-N-3 mirror row |
| `design_docs/sketches/storejournal.ail` | +~15 | LAW 5 + CF-N-3 row (re-verified live) |
| `host/broker/broker.go` | +~95 | pipeline rewiring + durable ordinal init |
| `host/broker/recover.go` | +~65 | effect-granularity walk + CF-N-2 justification |
| `host/broker/{broker,recover,episode}_test.go` | +~350 | window proof (incl. post-recovery dispatch), recovery, rows |
| `bench/BASELINE.md`, `README.md` | ~15 | re-measure + note |

~1,090 LOC (was ~1,010 in r0; the delta is the r1 ordinal-derivation method + its tests + the
post-recovery-dispatch arm, absorbed inside the ~1.5d estimate — the method mirrors existing
query shapes and is NOT on the overflow cut line: it is correctness, not polish).
**Byte-unchanged**: `host/store/schema.sql`, `host/store/store.go` (bindCommitIntentTx
and the Commit path are untouched — V12), `host/replay/**`, `host/capsule/**`,
`host/{hashref,canon,archive,registry}/**`, `host/daemon/**`, `cmd/**`, `world/**`,
`scripts/**`, `.github/**`, `go.mod`, `go.sum`.

## Conflict Surface (MANDATORY — every landed behaviour this design could collide with)

- **vs the commit journal's laws and callers.** `AppendIntent`/`AppendOutcome`/`GetReceipt`/
  `bindCommitIntentTx` semantics are unchanged; `validateIntent` NARROWS its domain by the
  reserved prefix — V14 enumerates all callers (broker `recover.go`, `episode_test.go`,
  `bench_test.go`, store tests) and none uses it. `bindCommitIntentTx` (V12) compares eight
  commit fields against a decoded commit intent; with disjoint namespaces it can never see an
  effect payload. **Collision risk named**: a future caller minting commit IDs with an `effect:`
  prefix now gets a structured rejection — that is the design working, and this doc is its
  record.
- **vs `PendingIntents`' existing consumers.** The only consumer is `broker.Recover` (V14). The
  `semantic_id` filter is behaviour-identical on every store that exists today (zero effect
  intents anywhere) and diverges exactly when it must (first pending effect intent). The
  cross-contamination test pins both directions.
- **vs `TestPreJournalMigrationPreservesExistingDDL`.** `schema.sql` untouched ⇒ the sha256
  prefix pin (V8) cannot red; the journal-table DDL it deliberately excludes is also untouched.
  This item ADDS what that test lacks in spirit: the effect-journal tests run against stores
  created with today's unmodified schema — which is the whole population, per Path A.
- **vs the broker's three-arm record law and Replay mode.** Record shapes, `recordConsistent`,
  Replay's zero-dispatch discipline: byte-untouched. The journal writes are NEW durability
  around the existing record write, not a change to it. Replay never reads the journal (it reads
  records); recovery reads the journal and never replays. No cycle.
- **vs M3.D's commit-boundary anchoring.** Unchanged and still necessary: the episode's commit
  intent is what binds the finished world; effect intents cover the per-effect window UNDER it.
  The episode test shows both operating on one episode.
- **vs the single-writer lock.** All new writes go through the injected `*store.Store` in the
  same transactions discipline; no second handle, no new writer.
- **vs `verify_ail.sh` totals.** Sketch module count stays 11 (edit, not addition); required
  manifests key on `world/` only (P7). The sketch's own contract count changes 7→8 and
  `len(tests[])` 30→(measured) — recorded at sprint time, never quoted.
- **vs the frozen log format / registry heads.** Untouched: the effect journal lives entirely in
  the `journal` + `objects` tables; zero log entries, zero head moves (the M3 Decision-7
  discipline holds — the broker still never writes log entries).

## Systemic-Issue Audit

Is this a one-off? No — it is the THIRD instance of one pattern: durable-intent-before-act
(commit journal SD.B/SD.C; commit-boundary anchoring M3.D; effect granularity here). This design
deliberately reuses the SAME table, the SAME receipt law, the SAME retry gate, and the SAME
paging discipline rather than minting a parallel mechanism — the unified-solution rule applied.
The one gap it does NOT fix (and names): the repo still has no schema-migration story (V7). That
is left deliberately unbuilt — building it without a consumer would be speculative kernel growth
(S3); Path B's rejection records where it becomes necessary. **Flagged for the charter's risk
register: any future item whose design requires a DDL change must budget the migration mechanism
as its first milestone.**

## Deferred Scope

| Deferred | Where it goes |
|----------|---------------|
| Deterministic reconciler for indeterminate-with-record effects | package-lane follow-up (`w-effect-reconcile`), consuming `GetEffectReceipt` + record objects |
| Per-handler reconciliation contracts (`Git.Commit` idempotency probes etc.) | same follow-up; the SD.C frozen contract is its spec |
| Journal compaction / archival of resolved effect rows | future store item, needs a growth measurement first |
| Schema versioning + fail-loud migration machinery | first milestone of whichever future item first needs a DDL change (named in the audit above) |
| **Repairing the inert DDL gate (V21–V23)** — give `TestPreJournalMigrationPreservesExistingDDL` teeth over the journal table and the upgrade-an-existing-store path, or retire the dead `sqlite_master` comparison and document the sha256 pin as the real mechanism | NOT this item (it is a pre-existing defect in landed `w-store-durability` code, and this item's AC1 deliberately does not depend on that gate). Raised by the iter-36 controller as its own charter queue row so it dies by evidence rather than by prose |
| Promoting `EffectInvocationID` derivation into the sketch | OD2 |
| Operator CLI/REST surface over effect receipts | clause-5/6 items; S7 applies when a user-facing surface exists (none does here — this is host-internal API, documented in godoc + this doc) |

## Acceptance Criteria

1. **AC1** — `git diff --exit-code host/store/schema.sql` (byte-unchanged). This `git diff` IS
   the gate; `TestPreJournalMigrationPreservesExistingDDL` staying green unmodified is recorded
   alongside it but is **not** evidence for AC1 (V21: a widened journal CHECK passes the entire
   repo suite). Citing a green run from an inert gate is the defect this doc exists downstream of.
2. **AC2** — golden-bytes tests for BOTH new codecs, committed.
3. **AC3** — namespace disjointness both directions: commit appender rejects `effect:` IDs;
   effect appender rejects non-conforming IDs; derived-ID mismatch rejected.
4. **AC4** — `PendingIntents` returns zero effect intents with a pending effect intent planted;
   `PendingEffectIntents` returns zero commit intents symmetrically.
5. **AC5** — outcome discipline: `AppendEffectOutcome` without a durable intent → structured
   error; second outcome → `DuplicateInvocationError`; idempotent identical-bytes re-append.
6. **AC6** — `GetEffectReceipt` walks all three states over a real store lifecycle.
7. **AC7** — **the window AND the resumption** (expanded in r1 per the quorum): the
   crash-simulation test (fault-injected store, handler ran, record/outcome lost) yields
   `indeterminate` on reopen, is surfaced by `Recover`, with the counting-probe registry at ZERO
   dispatches; **and the resumed broker then dispatches a NEW effect on the same episode
   successfully — no `DuplicateInvocationError`, its ID strictly past the indeterminate
   ordinal, its receipt walking to `resolved`**.
7b. **AC7b (r2, quorum round-2 fix)** — **allocation atomicity**: (i) two concurrent
    `AppendNextEffectIntent` calls for the SAME episode, run from two goroutines against one
    store, yield two DISTINCT ordinals and two durable rows — never a duplicate, never a lost
    write; (ii) an episode whose journal already holds ordinal `math.MaxInt64` yields a structured
    `OrdinalExhaustedError` rather than wrapping or silently reusing. Both are the reviewer's own
    named tests ("two concurrent allocations for the same episode, plus a structured error test
    for `MaxInt64` exhaustion").
7c. **AC7c (r2)** — the premise log carries evidence for the transaction/serialization behaviour
    the allocator relies on (the reviewer's third ask): the tx boundary is asserted by test, not
    by prose, and V28's split finding (cross-process covered by the writer lock, in-process NOT)
    is restated wherever the atomicity claim is made.
8. **AC8** — sketch re-verified live on the pinned binary: 8 named contracts in
   `verify.results[]` all `verified`; `len(tests[])` and `passed_tests` reported as two numbers,
   gate on `len(tests[])` ≥ baseline+new rows; no `skipped_tests` assertion.
9. **AC9** — ordering: denied effects journal NOTHING; failed effects journal intent + outcome
   (`Status: "failed"`); both tested.
10. **AC10** — the landed episode test passes with per-effect receipts `resolved` AND its
    commit anchor intact.
11. **AC11** — CF-N-2: justified bound + never-draining-fake test of the exceeded-pages error.
12. **AC12** — CF-N-3: `((false,true), true)` present at all three sites; `MUT-RETRY-XOR`
    demonstrated red on ONLY that row.
13. **AC13** — full gates: `verify_go.sh`, `verify_ail.sh` (11 modules / 4 identities / 14 world
    tests), bench smoke, `go test ./...` zero skips, protected-path `git diff --exit-code`,
    CI observed green before any box is checked. **No grep-based honest-claim gate** (the AC19
    lesson: it was unsatisfiable and monotonically increasing in honesty); the honest-claim
    check is the quorum reading the Scope note and Decision 5 against the code.

## Non-Vacuity — the named RED mutation for every gate (S6)

Per the iter-35 charter process fix, every mutation run must: assert its pattern matched
**exactly once**, print the applied diff, **confirm the mutant compiles** (a build failure
wearing a red gate's clothes proves nothing), name the FORM (a mutation name denotes a FAMILY),
report the reds it must NOT cause, and revert from a `cp` backup verified by hash.

| Gate | Named mutation (form) | Must red | Must stay green |
|------|----------------------|----------|-----------------|
| AC7 window | `MUT-INTENT-AFTER-DISPATCH` (form: in `host/broker/broker.go`, move the `AppendNextEffectIntent` call after the handler dispatch) — PRODUCTION | crash-window test (receipt reads `not-started` where `indeterminate` is required) | success-path pipeline tests (discriminating: the happy path cannot tell the orders apart) |
| AC7 resumption (r1, re-aimed by r2) | `MUT-ORDINAL-ZERO-RESUME` (form: in `host/store`'s `AppendNextEffectIntent`, replace the in-tx derived ordinal with the constant `0` — the r1 form named `broker.go`'s initializer, which r2 DELETED, so the mutation moves with the logic) — PRODUCTION, must compile | ONLY the new post-recovery-dispatch assertion (`DuplicateInvocationError` on the resumed episode's first new dispatch) | ALL pre-existing crash-window assertions (indeterminate receipt, `Recover` surfacing, zero dispatches) AND every fresh-episode pipeline test — 0 is the CORRECT initialization for a fresh episode, which is exactly why this mutant is discriminating: only the resumption arm can tell the derivation from the constant |
| AC3 | `MUT-NAMESPACE-DROP` (form: delete the `effect:` rejection in `validateIntent`, `host/store/journal.go`) — PRODUCTION | commit-side disjointness test | all pre-existing journal tests |
| AC4 | `MUT-PENDING-FILTER-DROP` (form: remove the `semantic_id` predicate from `PendingIntents`' SQL) — PRODUCTION | cross-contamination test (decode error on the planted effect intent) | `PendingEffectIntents` tests |
| AC5 | `MUT-OUTCOME-ORPHAN` (form: remove the intent-existence check in `AppendEffectOutcome`) — PRODUCTION | orphan-outcome test | duplicate-outcome test |
| AC2 | `MUT-CODEC-REORDER` (form: swap two wire fields in the effect-intent codec) — PRODUCTION | golden-bytes test | commit-codec goldens |
| AC9 | `MUT-DENIED-JOURNALED` (form: write an effect intent on the denied arm) — PRODUCTION | denied-no-journal test | denial-record tests |
| AC7/recovery | `MUT-EFFECT-REDISPATCH` (form: in the effect-granularity walk, invoke the registry's handler for an indeterminate finding) — PRODUCTION | counting-probe zero-dispatch test | surfacing tests |
| AC11 | `MUT-RECOVERY-UNBOUNDED` (form: replace the bounded page loop's condition with `true` in `host/broker/recover.go` — NOTE the family: iter-35 found TWO forms redding by different mechanisms; name which form ran) — PRODUCTION, must compile | never-draining-fake bound test | single-page and multi-page tests |
| AC7b atomicity (r2) | `MUT-ORDINAL-SPLIT-TX` (form: in `AppendNextEffectIntent`, compute the next ordinal in its OWN transaction/query and commit it before opening the append transaction — i.e. restore r1's split allocator; NOTE the FAMILY: a variant that merely moves the read one statement earlier inside the SAME tx is NOT this mutation and proves nothing, so the record must name which form ran) — PRODUCTION, must compile | the two-goroutine concurrent-allocation test (duplicate ordinal / duplicate row) | every single-threaded allocation test, the resumption test, and the exhaustion test — a split allocator is still CORRECT serially, which is exactly what makes this mutant discriminating |
| AC7b exhaustion (r2) | `MUT-ORDINAL-WRAP` (form: delete the `math.MaxInt64` exhaustion check so the increment wraps) — PRODUCTION | the `OrdinalExhaustedError` test | all ordinary allocation tests |
| AC12 | `MUT-RETRY-XOR` (form: `retryAllowed` body `!indeterminate \|\| reconciled` → `(!indeterminate) != reconciled` in `host/broker/recover.go:46`) — PRODUCTION | ONLY the `((false,true), true)` row — truth table: `(F,F)`→T✓, `(T,F)`→F✓, `(T,T)`→T✓, `(F,T)`→**F✗** — proving the new row is the only non-vacuous witness | all three pre-existing rows |

## Axiom Compliance

| Axiom | Score | Justification |
|-------|-------|---------------|
| A1: Determinism | +2 | Logical time caller-supplied; canonical JSON codecs; deterministic ID derivation; golden bytes |
| A2: Replayability | +2 | The exact gap between "recorded" and "replayable-with-confidence" closes: absence of intent now PROVES non-execution |
| A3: Effect Legibility | +2 | Every dispatched effect leaves a durable pre-dispatch statement; the audit trail has no memory-only phase |
| A4: Explicit Authority | +1 | Authority law untouched; the journal makes its exercise crash-honest |
| A5: Bounded Verification | +1 | Bounded page walks; CF-N-2's bound justified and tested |
| A6: Safe Concurrency | 0 | Single-writer model inherited; no new concurrent surface |
| A7: Machines First | +2 | The ordering law is sketch-frozen and Z3-provable; receipts are typed, structured states |
| A8: Minimal Syntax | 0 | No language surface |
| A9: Cost Visibility | +1 | The journal's runtime price lands inside the existing bench row, delta reported |
| A10: Composability | +1 | Reconciler/compaction stack on the receipt API without kernel changes |
| A11: Structured Failure | +2 | Indeterminacy is a typed state with a typed finding, never an absence |
| A12: System Boundary | +1 | Durability in the kernel; policy (reconciliation) deferred to packages |

**Net Score: +15** ✅ — hard axioms A1/A3/A4/A7 all ≥ 0.

## Premise Verification Log (live evidence, this session, 2026-07-29, worktree at `0f2afad`)

Every row executed first-party this session unless marked **CLAIM**. The iter-36 controller
published its own measurements; per the mission's method, every load-bearing one was RE-RUN here
rather than cited — all reproduced, none refuted.

| # | Claim | Command | Actual result |
|---|-------|---------|---------------|
| V1 | Pinned binary is clean released v0.30.0 | `/tmp/ailang-v0300/ailang --version` | `AILANG v0.30.0`, commit `e37b370`, built `2026-07-19T09:27:00Z`, no `-dirty` |
| V2 | `validateIntent` requires six non-zero commit refs + non-empty ID; `ObservedHead` optional-if-zero | read `host/store/journal.go:210-236` | confirmed verbatim — the six-field slice + optional ObservedHead branch; hence a pre-dispatch COMMIT intent for an effect is structurally impossible and the effect shape must be its own payload |
| V3 | Journal DDL is `kind IN ('intent','outcome')` + `UNIQUE (invocation_id, kind)` | read `host/store/schema.sql:81-87` | confirmed byte-for-byte as quoted in the correction note |
| V4 | **P1**: new kind label DDL-rejected | fresh DB from repo `schema.sql`; `INSERT … kind='effect-intent'` (sqlite3 3.51.0 CLI) | `Error: stepping, CHECK constraint failed: kind IN ('intent','outcome') (19)` |
| V5 | **P2 + P2b**: one intent per ID; distinct ID clears it | second `kind='intent'` same ID; then ID `'inv-1/effect/0'` | `UNIQUE constraint failed: journal.invocation_id, journal.kind (19)`; distinct-ID insert OK — Path A's mechanism confirmed live |
| V6 | **P3**: DDL widening silently dropped on an existing store | store from today's schema → widen CHECK in a copy (literal asserted to occur exactly once) → re-apply to the EXISTING file → read `sqlite_master` → retry insert; control: fresh DB from widened schema | re-apply rc=0 **no error**; journal DDL **unchanged** (`kind IN ('intent','outcome')`); new-kind insert fails on the OLD constraint; fresh DB accepts it — **fail-open confirmed** |
| V7 | **Negative existence**: no migration machinery anywhere in `host/` | `grep -rn 'user_version\|ALTER TABLE\|migrate\|Migrate\|schema_version' host/` | **0 hits** |
| V8 | The DDL-drift gate excludes the journal table and never tests the upgrade path; pins the pre-journal prefix by sha256 | read `host/store/journal_test.go:340-371` | confirmed: `delete(after, "journal")`; sha256 pin `43a9c80b…`; fresh-DB + reopen only — no existing-store DDL-edit path |
| V9 | Receipt law is per-invocation via `journalRowFor(id, kind)`; `PendingIntents` decodes EVERY pending row via `decodeJournalIntent` and pages by seq keyset under `MaxPendingIntentsPage = 1000` | read `host/store/journal.go:364-470, 17-18` | confirmed — including the decode-every-row behaviour that forces the `semantic_id` filter (Decision 2) |
| V10 | **CF-N-2**: the magic bound | read `host/broker/recover.go:10` | `const maxRecoveryPages = 1 << 20` — that literal, that line, no justification comment |
| V11 | **CF-N-3**: `(false,true)` untested at all three sites | read `host/broker/recover.go:46-48`, `host/store/recover_test.go:26`, `design_docs/sketches/storejournal.ail` LAW 3 `tests[]` | confirmed: sketch rows are exactly `((F,F),T), ((T,F),F), ((T,T),T)` — `(F,T)` absent; both Go mirrors match |
| V12 | `bindCommitIntentTx` byte-compares all eight commit fields inside the tx | read `host/store/store.go:800-830` | confirmed — eight-field table incl. `ObservedHead`; decodes the intent payload as a commit intent (disjoint namespaces keep effect payloads out of its reach) |
| V13 | `store.Open` applies the schema by one `db.Exec(schemaSQL)` (`//go:embed schema.sql`), `applySchema=false` only for read-only handles | read `host/store/store.go:37-38, 244-263` | confirmed — the mechanism that makes V6 the live upgrade behaviour |
| V14 | **Negative existence**: no caller uses an `effect:`-style ID; journal-API callers enumerated | `grep -rn 'PendingIntents\|GetReceipt\|AppendIntent\|AppendOutcome' host/ cmd/ --include='*.go' -l`; grep of `InvocationID` literals | callers: `host/broker/{recover.go,recover_test.go,episode_test.go}`, `host/daemon/bench_test.go`, `host/store/{journal,journal_test,recover_test,crash_test}` — none in `cmd/`; ID literals found: `"golden"`, `"recover-…"` — zero use any reserved prefix |
| V15 | Probe-harness caveat, labelled | — | V4–V6 ran on system sqlite3 CLI 3.51.0; the store uses `modernc.org/sqlite`. The constraints under test are SQLite SQL semantics (CHECK / UNIQUE / `CREATE TABLE IF NOT EXISTS` no-op), identical across bindings; extended error codes differ (CLI `(19)` vs the controller's driver-reported `(275)/(2067)`). The controller's driver-side run and this CLI run agree on every behaviour |
| V16 | `storejournal.ail` current baseline on the pinned binary, z3 on PATH (`/opt/homebrew/bin/z3`) | `cd design_docs && AILANG_BIN=… ai-check -timeout 15s sketches/storejournal.ail`; `ailang test --format json …` | `verify: {verified: 7, counterexample: 0, skipped: 0, errors: 0}`, all 7 named identities `verified` (listed in the header); **`len(tests[]) == 30`**, `passed_tests == 37`, `failed_tests == 0` — the two numbers reported separately per the standing rule |
| V17 | Duplicate/coverage gate: no planned doc covers this topic | `ls design_docs/planned/`; `grep -ril 'effect.intent\|effect journal\|dispatch.*record window' design_docs/planned/` | `planned/` holds `w-log-epoch-decision.md`, `w-mcp-projection.md`; grep: **no overlap** |
| V18 | Charter facts: pre-ratified-in-principle scope, forbidden dummy refs, CF-N-2/CF-N-3 inheritance, queue row text | read `design_docs/world-mission.md` (queue row 4c ~1140-1148; ratification block ~858-910) | confirmed verbatim, incl. "design doc + quorum at pick as usual" and the explicit placeholder-ref prohibition |
| V19 | LOC and day estimates | **CLAIM** — no code exists at doc time | the sprint's gates are the check; estimates follow M3's measured ~1.0× planner multiplier |
| V20 | The MJ.A sketch edit verifies on v0.30.0 | **CLAIM** — the edit is NOT authored at design time (stated in the header); LAW 5 is a bool-only total-theorem ensures and a `tests[]` row, both shapes this exact sketch already proves 7 of / 30 of | MJ.A verifies live BEFORE landing and records the real numbers; a failure there is a design defect to bring back here, not to patch around |

**Rows V21–V25 were measured by the ITERATION-36 CONTROLLER**, outside the designer's session, in
a scratch worktree at `0f2afad` (removed afterwards; every mutation applied under an
exactly-once assertion with the diff printed, the mutant confirmed to COMPILE, and reverted from a
`cp` backup verified byte-identical by `sha256`, per the iter-35 instrument discipline).

| # | Claim | Command | Actual result |
|---|-------|---------|---------------|
| V21 | **Nothing in the repo guards the journal table's DDL** | `MUT-JOURNAL-DDL-WIDEN` (form: `kind` CHECK widened to admit `'effect-intent','effect-outcome'` in `host/store/schema.sql`; 1/1 pattern match; diff printed) → `go build` rc=0, `go vet` rc=0, then `AILANG_BIN=… go test ./... -count=1` | **rc=0, 10/10 packages `ok`, 0 FAIL** — the mutant compiles and survives the entire suite. Reverted; `schema.sql` sha256 `13893a29…` restored byte-identical |
| V22 | The gate's own named `MUT-DDL-DRIFT` reds by a DIFFERENT mechanism than documented | `MUT-DDL-DRIFT` (form: `CHECK (written_by <> '')` added to `log_entries`; 1/1 match; `go build` rc=0) → `go test -run TestPreJournalMigrationPreservesExistingDDL` | **FAIL**, message `journal_test.go:345: pre-journal schema source drifted: sha256=7358d876…` — the **source-text pin**, which `t.Fatalf`s at line 345, *before* any DB is built. The `sqlite_master` comparison at line 368 never executes |
| V23 | The `sqlite_master` byte-identity comparison is **inert** | `MUT-DDL-COMPARE-DEAD` (form: `if !reflect.DeepEqual(after, before)` → `if false && !reflect.DeepEqual(...)`, keeping both vars used so it must compile; 1/1 match; `go vet` rc=0) → `go test ./host/store/ -count=1` | **`ok`, package fully green** — the comparison contributes zero discrimination. `journal_test.go` reverted, sha256 `9fe7484b…` byte-identical |
| V24 | **Path A's mechanism holds on the REAL driver** — closes the designer's own V15 harness caveat | throwaway test in `host/store` on `modernc.org/sqlite` via `store.Open`: insert `kind='intent'` for `episode-1`, then for `effect:episode-1:0` and `effect:episode-1:1` | all three accepted; **3 `intent` rows coexist**. The designer measured this on the sqlite3 CLI; it now also holds on the binding production uses. Probe removed |
| V25b | **The r1 fix's crux holds on the REAL driver** — controller re-run of the designer's V27, which ran on the sqlite3 CLI | throwaway test in `host/store` on `modernc.org/sqlite` via `store.Open`: plant `effect:episode-1:0` RESOLVED (intent+outcome), `effect:episode-1:1` INDETERMINATE (intent only), plus adversarial `effect:episode-10:7` and non-decimal `effect:episode-1:abc`; scan `kind='intent' AND invocation_id LIKE 'effect:episode-1:%'` | returns `[effect:episode-1:0, effect:episode-1:1, effect:episode-1:abc]` — **the RESOLVED ordinal 0 IS seen** (max=1 → next=2), and `effect:episode-10:7` does **not** leak (the `:` terminator discriminates). Confirms the designer's reason for rejecting the reviewer's `len(episode.History)` variant. **Scope honesty**: in the same probe `PendingIntents` returned empty, but that is a PROBE ARTIFACT — the rows were planted with an `object_ref` having no `objects` row, so the method's JOIN drops them. It neither supports nor refutes V9; V9 stands on its code read. Probe removed |
| V25 | The designer's V16 sketch baseline reproduces | `ai-check -timeout 15s sketches/storejournal.ail` on the pinned binary; `ailang test --format json` (JSON parsed after the leading `→ Running tests…` line) | `verify {verified:7, counterexample:0, skipped:0, errors:0}`, all 7 named identities `verified`; **`len(tests[]) == 30`**, `passed_tests == 37`, `failed_tests == 0` — **identical to V16**, two numbers reported separately, no `skipped_tests` assertion |

**Rows V26+ were measured by the DESIGNER in the r1 revision pass** (this session, 2026-07-29,
this worktree at `9be4140`), verifying the quorum-blocking fix's mechanism first-party.

| # | Claim | Command | Actual result |
|---|-------|---------|---------------|
| V26 | **No existing store method can derive the resumption ordinal** — a new method is required (hence the Decision 2 / freeze / aggregate updates) | read `host/store/journal.go` in full (471 lines) | confirmed: `PendingIntents` (`journal.go:427-470`) filters `NOT EXISTS … kind='outcome'`, so RESOLVED effect intents vanish from its result set (and `PendingEffectIntents` mirrors it); `journalRowFor`/`GetReceipt` (`journal.go:364-423`) are strictly per-`(invocation_id, kind)`; no other method queries `journal` by ID prefix. A pending-only derivation undercounts by every resolved effect |
| V27 | **The `NextEffectOrdinal` candidate-set semantics hold on the real DDL, including the adversarial cross-episode suffix, and the pending-only view provably misses resolved ordinals** | fresh DB from repo `schema.sql` (sqlite3 CLI 3.51.0 — SQL-semantics probe, same V15 caveat and V24 driver-parity precedent); rows: commit intent `ep-1`; effect intents ordinals 0 (with outcome → RESOLVED) and 1; adversarial `effect:ep-1:x:0` from an episode ID extending the prefix; then the exact-prefix `substr` query vs the pending-only query | prefix query returns exactly `{effect:ep-1:0, effect:ep-1:1, effect:ep-1:x:0}` — the resolved row present, the adversarial suffix `x:0` present-but-non-decimal (skipped by the strict decimal parse, so max ordinal = 1, next = 2); pending-only query returns `{effect:ep-1:1, effect:ep-1:x:0}` — **resolved ordinal 0 MISSING**, proving both the new method's necessity and the collision in any pending-only (or, analogously, record-history-based) derivation |
| V28 | **The round-2 objection is HALF refuted and HALF real — and the fix is applied in full anyway** (controller, iter-36) | `grep -n 'LOCK_EX\|LOCK_NB\|WriterAlreadyActive' host/store/writer_lock*.go host/store/store.go`; `grep -n 'sync\.\|Mutex\|RWMutex' host/store/store.go host/store/journal.go` | **Cross-process limb REFUTED**: `store.Open` takes a non-waiting exclusive lock and a second writer gets `*WriterAlreadyActive` (landed A1, proven cross-process), so "two broker instances sharing a store" cannot arise. **In-process limb REAL**: **zero** mutex hits in `store.go`/`journal.go` — two goroutines in one process can interleave a split read→allocate→append freely. The reviewer's transactional allocator is adopted in full rather than narrowed to the surviving limb |
| V29 | **Doc size silently degraded its own review** (controller, iter-36 — process finding) | quorum round-1 artifact `w-effect-journal-2026-07-29T05-50-34Z.json` | `gpt5-6-sol` absent, `absent_reason: budget`, `error: estimated cost $0.1005 (doc ~13952 input tok) exceeds cap $0.1000 (pre-flight refusal, zero spend)`. Round 1 therefore ran ONE-EYED, and the reviewer it lost is the one that later found the TOCTOU. The quorum degraded correctly and named the absentee — it was not a silent pass — but the cap was raised to $0.25 for round 2 and both reviewers were present |

## Open Decisions (escalated with recommended defaults — the sprint proceeds on the defaults)

- ~~**OD1 — `GetReceipt` defense-in-depth.**~~ **CLOSED in r2** — promoted from a MAY to a MUST
  (Decision 2), applying `gemini-3-1-pro`'s round-2 fix verbatim. No longer an open decision.
- **OD2 — sketch-side ID derivation.** Promote `EffectInvocationID`'s canonical text form into
  `storejournal.ail` with named tests? **Default: Go-side goldens only in this item**; promotion
  waits for a fluency-protocol check of string stdlib support on v0.30.0 (never from memory —
  the S5 rule), recorded as a possible follow-up, not a silent drop.
- **OD3 — episode driver's journal use in production.** The episode driver lives in test code
  today (V14); the first production driver (M4/clause-6 wiring) must adopt intent-before-dispatch
  from birth. **Default: this doc is the binding record**; a note in `host/broker/broker.go`'s
  package comment points here.
- **OD4 — outcome+record atomicity.** Should `AppendEffectOutcome` optionally carry the record
  object into its own tx, closing the Decision-5 residual entirely? **Default: no in this item**
  — it widens the kernel API mid-reopen for a residual that is already fail-closed and
  reconcilable; escalate to its own small design if the residual is ever observed in practice.

## Related Documents

- [w-effect-broker-m3.md](../implemented/w-effect-broker-m3.md) — the broker this rewires; its
  Decision-3 corrected supersession note and M3.D section are this item's origin story
- [w-store-durability.md](../implemented/w-store-durability.md) — the commit journal + receipt
  law this item extends to a second payload shape
- [w-m1-ailang-hardening.md](../implemented/w-m1-ailang-hardening.md) — the verifier facts
  binding on the MJ.A sketch edit (two-numbers rule, no `skipped_tests` assertion, premise V14)
- `design_docs/sketches/storejournal.ail` — the law this item edits (LAW 5 + CF-N-3 row)
- Charter queue row 4c + ratification `c26b27d` — scope, pre-ratification-in-principle, and the
  placeholder-ref prohibition

## Quorum verification log

**Round 1** (`2026-07-29T05-50-34Z`, cap $0.10/reviewer, metered **$0.034**) — **BLOCKED**.
`gemini-3-1-pro` **reject**: the in-memory per-episode ordinal counter is never re-initialized
after a crash, so a resumed broker's first dispatch collides with the durable
`effect:<episodeID>:0` and bricks the episode with `DuplicateInvocationError`. `gpt5-6-sol`
**ABSENT — `budget`**: pre-flight refusal, *"estimated cost $0.1005 (doc ~13952 input tok) exceeds
cap $0.1000"*, **zero spend**; the quorum degraded to N−1 and named the absentee (never a silent
pass). Controller verdict: pass.

**r1 revision** (same Fable designer lane, bounded, rc=0). The objection was accepted as a real
defect. The fix derives the ordinal from durable journal state rather than a fresh zero — and the
designer **declined the reviewer's own `len(episode.History)` variant with a reason**: history
counts *records*, and an indeterminate effect is precisely an intent whose record was lost, so
after exactly the AC7 crash that variant returns the indeterminate intent's own ordinal and still
collides. r0's frozen claim *"collision-free by construction within an episode"* was marked FALSE
across a crash boundary and corrected everywhere it was restated. Controller re-verified the crux
on `modernc.org/sqlite` rather than the designer's CLI harness (**V25b**).

**Round 2** (`2026-07-29T06-00-47Z`, cap raised to $0.25 so the round-1 absentee could run,
metered **$0.130**; **both reviewers present**) — **BLOCKED**, but only one vote.
`gemini-3-1-pro` **PASS** (its round-1 objection resolved), with a concrete non-blocking fix: add a
`strings.HasPrefix(id, "effect:")` guard to `GetReceipt`. `gpt5-6-sol` **reject**: the r1 fix is
**not atomic** — `NextEffectOrdinal` and `AppendEffectIntent` are separate operations, so two
allocations can observe the same maximum; *"the resumption fix merely replaces the restart
collision with a TOCTOU collision."* Controller verdict: pass.

**r2 — the charter's NARROW-REFINEMENT CARVE-OUT, applied by the controller.** Both limbs are
satisfied: the remaining blocking objection (a) carries a concrete **reviewer-authored**
`proposed_fix`, and (b) disputes no design DIRECTION — it endorses Path A, the journal shape,
intent-before-dispatch and the ID namespace, and objects only to the *atomicity* of one method.
That is the "determinism" limb the carve-out names. The reviewer's fix was applied **verbatim**:
a single transactional `AppendNextEffectIntent(episodeID, intentWithoutID) (id, ordinal, err)`;
the broker-side ordinal cache **removed**; `MaxInt64` exhaustion given a structured error; and the
reviewer's two named tests added as **AC7b**, with premise-log evidence for the transaction
behaviour as **AC7c**. `gemini-3-1-pro`'s `GetReceipt` guard was adopted too, closing **OD1**.

**One thing the controller measured and did NOT use to soften the fix (V28).** Half of the round-2
objection is refuted by landed code: two broker *processes* cannot share a store as writers
(`store.Open`'s non-waiting exclusive lock → `*WriterAlreadyActive`, proven cross-process). The
other half is real and unguarded — there is **no mutex anywhere** in `store.go`/`journal.go`, so
two goroutines in one process can interleave. The transactional allocator is adopted in full
regardless. Narrowing a fix to the part of an objection you can refute is how a real defect
survives a review, and this mission has paid for that lesson in the opposite direction already.

**Not force-passed.** Standing rule 2 still forbids proceeding over a contested design DIRECTION;
nothing here was contested at that level, and no objection was overridden — each was *satisfied*.
Next stop is the sprint-planner, not the executor: this doc is a design, and its milestones are
still unplanned.
