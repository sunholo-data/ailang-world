# w-effect-journal — Closing the dispatch→record window at effect granularity (the EFFECT journal)

**Status**: Planned (NEW-DOC, queue row 4c, designer rotation iter-36 — design quorums at pick)
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
  The full delta is in Decision 2: **four new methods, two changed methods, zero removed, zero
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
| Intent is durable BEFORE dispatch; outcome AFTER the record | The entire point of the item — the ordering IS the window closure; frozen as a law in the sketch | this doc | high |
| Denied effects write NO journal rows | No dispatch → no window → nothing to journal; the denial record (landed M3 discipline) is already durable and complete; journal growth stays 2 rows per *dispatched* effect | this doc | low |
| `PendingIntents` gains a `semantic_id` filter; `PendingEffectIntents` is a separate method | Without the filter, commit recovery decodes effect payloads as commit intents and errors on every store with a pending effect (V9: `PendingIntents` decodes every row via `decodeJournalIntent`) | this doc | medium |
| CF-N-2 and CF-N-3 are acceptance criteria with their own mutations | Inherited carry-forwards die by evidence, not by prose | charter | low |

### Design Freeze (the sprint must not renegotiate these)

- [x] `host/store/schema.sql` byte-for-byte unchanged; **no DDL change, no migration machinery**.
- [x] The `host/store` method delta is EXACTLY Decision 2's table: `AppendEffectIntent`,
  `AppendEffectOutcome`, `GetEffectReceipt`, `PendingEffectIntents` (new);
  `validateIntent` + `PendingIntents` (changed, as specified); nothing else.
- [x] The ordering law: **no handler dispatch without a prior durable effect intent**; the
  outcome append happens after the effect record object is durable; both directions tested and
  mutation-pinned.
- [x] Effect invocation IDs are `effect:<episodeInvocationID>:<ordinal>` (ordinal decimal ≥ 0);
  the commit appender REJECTS the `effect:` prefix; the effect appender REQUIRES the full shape.
  No other ID form enters the effect journal.
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
  effect namespace, two effects collide only at identical `(episodeID, ordinal)` — which is the
  broker's own sequencing bug if it happens, and surfaces as `DuplicateInvocationError` with
  differing bytes, the landed discipline.
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
| `AppendEffectIntent(id, EffectIntent)` | **new** | mirrors `AppendIntent`'s tx discipline (idempotent on identical bytes, `DuplicateInvocationError` on differing bytes, object + journal row in ONE tx, `nextJournalSeqTx`); validates the `effect:` shape, `Ordinal >= 0`, non-empty `Effect`/`Scope`/`EpisodeID`, non-zero `RequestRef`, and that `id` == the derived `effect:<EpisodeID>:<Ordinal>` |
| `AppendEffectOutcome(id, EffectOutcome)` | **new** | mirrors `AppendOutcome`: requires a durable effect intent, rejects a second outcome, requires `Status ∈ {succeeded, failed}` and non-zero `RecordRef` |
| `GetEffectReceipt(id)` | **new** | the same three-state receipt law (`not-started` / `indeterminate` / `resolved`) over the effect payloads; rejects non-`effect:` IDs with a structured error |
| `PendingEffectIntents(limit, fromIndex...)` | **new** | `PendingIntents`'s keyset paging, filtered to `o.semantic_id = 'world/effect-intent/v1'`, same `MaxPendingIntentsPage` bound |
| `validateIntent` | **changed** | adds: reject `InvocationID` with the reserved `effect:` prefix (structured `InvocationMismatchError`) — the commit-side half of P3's disjointness |
| `PendingIntents` | **changed** | query gains `AND o.semantic_id = 'world/journal-intent/v1'` — behaviour-preserving on every existing store (no effect intents exist anywhere yet, V14), and REQUIRED once they do (V9: today's query decodes every pending row as a commit intent and would error on the first effect payload) |

`GetReceipt` is deliberately **unchanged**: with disjoint namespaces a commit ID can never
resolve an effect row. Defense-in-depth (rejecting `effect:` IDs there too) is a MAY, decided at
sprint time; the disjointness tests are the MUST.

## Decision 3 — The ID namespace

`effect:<episodeInvocationID>:<ordinal>`. Prefix-based (not infix) so the check is one
`strings.HasPrefix`; the episode ID is embedded verbatim so recovery findings group by episode
with zero extra lookups; the ordinal is the broker's per-episode dispatch counter, making the
derivation deterministic and collision-free by construction within an episode. The derivation
lives in ONE exported function (`store.EffectInvocationID(episodeID string, ordinal int64)`) used
by both the broker and the tests — never re-derived inline (the mirror-drift lesson). Its
canonical text form gets named golden tests (S1's canonical-text rule); promotion of the
derivation into the `.ail` sketch is deliberately NOT promised here because it needs string
stdlib verification under the fluency protocol at sprint time (Open Decision OD2).

## Decision 4 — The rewired broker pipeline (allowed arms only)

```
decide → (denied? → record denial, return — UNCHANGED, no journal)
       → PutObject(request payload)                     [durable, true]
       → AppendEffectIntent                             [durable — THE WINDOW CLOSES HERE]
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

- **files**: `host/store/journal.go` (+~230 — the two codecs, four new methods, `validateIntent`
  namespace check, `PendingIntents` filter, `EffectInvocationID`), `host/store/journal_test.go`
  (+~260 — golden bytes ×2, disjointness both directions, duplicate/orphan-outcome discipline,
  receipt three-state walk, `PendingIntents` cross-contamination, `PendingEffectIntents` paging),
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

- **files**: `host/broker/broker.go` (+~90 — intent-before-dispatch, ordinal counter per
  session/episode, outcome-after-record; denied arm untouched), `host/broker/recover.go` (+~60 —
  the effect-granularity page walk + widened finding type), `host/broker/broker_test.go` /
  `recover_test.go` (+~280 — the **crash-window test**: fault-injected store loses everything
  after the handler runs but before record/outcome; reopen; `GetEffectReceipt` =
  `indeterminate`; `Recover` surfaces it with the counting-probe registry asserted at ZERO
  dispatches; plus ordering tests, denied-no-journal, failed-arm-journals-outcome),
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
| `host/store/journal.go` | +~230 | codecs, 4 new methods, 2 changed methods |
| `host/store/journal_test.go` | +~260 | new-surface tests |
| `host/store/recover_test.go` | +~5 | CF-N-3 mirror row |
| `design_docs/sketches/storejournal.ail` | +~15 | LAW 5 + CF-N-3 row (re-verified live) |
| `host/broker/broker.go` | +~90 | pipeline rewiring |
| `host/broker/recover.go` | +~65 | effect-granularity walk + CF-N-2 justification |
| `host/broker/{broker,recover,episode}_test.go` | +~330 | window proof, recovery, rows |
| `bench/BASELINE.md`, `README.md` | ~15 | re-measure + note |

~1,010 LOC. **Byte-unchanged**: `host/store/schema.sql`, `host/store/store.go` (bindCommitIntentTx
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
7. **AC7** — **the window**: the crash-simulation test (fault-injected store, handler ran,
   record/outcome lost) yields `indeterminate` on reopen, is surfaced by `Recover`, with the
   counting-probe registry at ZERO dispatches.
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
| AC7 window | `MUT-INTENT-AFTER-DISPATCH` (form: in `host/broker/broker.go`, move the `AppendEffectIntent` call after the handler dispatch) — PRODUCTION | crash-window test (receipt reads `not-started` where `indeterminate` is required) | success-path pipeline tests (discriminating: the happy path cannot tell the orders apart) |
| AC3 | `MUT-NAMESPACE-DROP` (form: delete the `effect:` rejection in `validateIntent`, `host/store/journal.go`) — PRODUCTION | commit-side disjointness test | all pre-existing journal tests |
| AC4 | `MUT-PENDING-FILTER-DROP` (form: remove the `semantic_id` predicate from `PendingIntents`' SQL) — PRODUCTION | cross-contamination test (decode error on the planted effect intent) | `PendingEffectIntents` tests |
| AC5 | `MUT-OUTCOME-ORPHAN` (form: remove the intent-existence check in `AppendEffectOutcome`) — PRODUCTION | orphan-outcome test | duplicate-outcome test |
| AC2 | `MUT-CODEC-REORDER` (form: swap two wire fields in the effect-intent codec) — PRODUCTION | golden-bytes test | commit-codec goldens |
| AC9 | `MUT-DENIED-JOURNALED` (form: write an effect intent on the denied arm) — PRODUCTION | denied-no-journal test | denial-record tests |
| AC7/recovery | `MUT-EFFECT-REDISPATCH` (form: in the effect-granularity walk, invoke the registry's handler for an indeterminate finding) — PRODUCTION | counting-probe zero-dispatch test | surfacing tests |
| AC11 | `MUT-RECOVERY-UNBOUNDED` (form: replace the bounded page loop's condition with `true` in `host/broker/recover.go` — NOTE the family: iter-35 found TWO forms redding by different mechanisms; name which form ran) — PRODUCTION, must compile | never-draining-fake bound test | single-page and multi-page tests |
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
| V25 | The designer's V16 sketch baseline reproduces | `ai-check -timeout 15s sketches/storejournal.ail` on the pinned binary; `ailang test --format json` (JSON parsed after the leading `→ Running tests…` line) | `verify {verified:7, counterexample:0, skipped:0, errors:0}`, all 7 named identities `verified`; **`len(tests[]) == 30`**, `passed_tests == 37`, `failed_tests == 0` — **identical to V16**, two numbers reported separately, no `skipped_tests` assertion |

## Open Decisions (escalated with recommended defaults — the sprint proceeds on the defaults)

- **OD1 — `GetReceipt` defense-in-depth.** Should `GetReceipt` structurally reject `effect:` IDs
  (beyond namespace disjointness making them unreachable)? **Default: yes, one guard + one
  test** — cheap, and it turns a silent empty result into a structured error.
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
