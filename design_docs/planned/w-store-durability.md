# w-store-durability — What the Kernel Guarantees Across a Crash (CF-B-2 + the durable intent/outcome journal + commit receipts)

**Status**: Planned (NEW-DOC, queue item 4b, iter-25)
**Date**: 2026-07-28
**Charter clause**: clause-1 (deterministic kernel)
**Verified against**: **`AILANG v0.30.0`** — the pinned released binary at `/tmp/ailang-v0300/ailang`
(`AILANG v0.30.0`, commit `e37b370`, built 2026-07-19 — no `-dirty` suffix; re-confirmed
first-party this session). Every `.ail` claim in this doc was checked live on that binary with z3
on PATH and the contracts proven in `verify.results[]` (never a bare exit code); the checked
artifact is [`sketches/storejournal.ail`](../sketches/storejournal.ail) (Appendix A), created and
verified WITH this doc — **7/7 contracts Z3-`verified`, 0 errors, 25 named tests pass** (`len(tests[])`
= 25; `passed_tests`/`total_tests` read **32** because they ALSO count the 7 contract-derived
properties — the landed `verify_ail.sh` correction D-B, controller-remeasured iter-25 and
re-measured live after the round-1 revision added LAW 6 `intentBindsCommit`), and
the full `verify_ail.sh` sweep re-run green at 10 modules with the 4-identity / 14-test `world/`
totals unperturbed.
**Traces to**: [DESIGN.md](../DESIGN.md) §1 (replay "is not aspirational — a pure transition plus
recorded effect results reconstructs any historical state bit-for-bit"), §7/§9 (recording as the
replay input), §15 (the store's layer), §14 (boundaries); charter clause 1 + guardrails
([world-mission.md](../world-mission.md) queue item 4b, the verbatim source of both halves)
**Depends on**: [w-world-library-m1.md](../implemented/w-world-library-m1.md) (the store this
hardens) and [w-worldd-m2.md](../implemented/w-worldd-m2.md) (the single-writer daemon + the
CF-B-2 carry-forward). **Feeds**: [w-effect-broker-m3.md](w-effect-broker-m3.md) (PARKED —
its round-2 objection 2A is half (ii) of this item) and
[w-mcp-projection.md](w-mcp-projection.md) (BLOCKED — its prerequisite 3, the commit-boundary
contract, is Decision 4 of this doc).
**Estimated**: ~1.5–2 days (the queue band's honest middle; the pre-committed overflow cut line
is in the Milestones section)

> **Scope note — one question, three symptoms.** This item exists because the same
> kernel-durability question arrived from three independent directions inside two iterations:
> **(i) CF-B-2** — `store.Commit` durably persists log entries (and, measured this session,
> worlds) that no reader can ever load, and the daemon's refusal at the REST boundary is a
> boundary patch over a kernel defect; **(ii) the broker effect journal** — `gpt5-6-sol`'s
> `w-effect-broker-m3` round-2 objection: a handler can execute a REAL external effect (an
> `FS.Write`, a `Git.Commit`, a **paid** `Model.Infer`) and die before its record is durable,
> leaving replay unable to distinguish "never executed" from "executed but record lost";
> **(iii) the commit-boundary contract** — the same reviewer, reviewing `w-mcp-projection` from a
> different direction: an HTTP context can expire while an atomic commit is in flight, so absent
> an atomic "not-started versus committed" contract a caller can observe cancellation on a commit
> that succeeded. Independent corroboration from multiple directions is why this is real work
> rather than bookkeeping. The answer to all three is ONE design: the kernel validates what it
> writes (Decision 1), tells the truth about what it already wrote (Decision 2), and provides a
> durable intent/outcome journal with a three-state receipt law (Decisions 3–5) that both the
> commit path and the future broker consume. What this item is NOT: the broker itself (M3,
> parked), any REST/CLI projection of receipts (the M2 route table stays frozen), or log-repair
> tooling (Decision 2 rules repair out on immutability grounds).

---

## Motivation

Clause 1 of the ratified bar promises *"an immutable, content-addressed world store with an
append-only transition log and proven deterministic replay."* Two measured facts break the spirit
of that promise today:

**First, the append-only log can contain entries that are not readable — and the store keeps
functioning around the hole.** Re-probed first-party this session (Premise Log V2, full
transcript): committing an entry whose `PrevEntryHash` is the zero `hashref.HashRef` succeeds
(`err = <nil>`), durably persists `prev_entry_hash_ref = ""` (TEXT NOT NULL is satisfied by
`''`), advances the selected head to the new world — whose `LogHead` now addresses an entry that
`GetLogEntry` can never return — and then accepts a perfectly legal follow-on commit that chains
entry 2 onto the poisoned entry 1 and reads back fine. The log grows a **permanent hole
mid-chain** with readable entries on both sides, no detection, and no recovery path. The blast
radius reaches the REST surface: both `GET /v1/log/{index}` and the bounded range loop
`GET /v1/log` serve `store.GetLogEntry`, and the range handler returns 500 for the **whole
range** the moment the loop touches a poisoned index (code-verified, V9) — a durable
denial-of-read on the log surface. And the class is **wider than the charter row states**: this
session's probe (V3) confirmed a zero `TransitionFn` poisons its entry identically, and a zero
`World.LogHead` produces a world row `GetWorld` cannot load — `Commit` validates **none** of the
seven ref fields it renders to TEXT (V4), so every one of them is a CF-B-2 waiting for a caller.

**Second, nothing in the kernel can say what happened across a crash.** The parked broker design
names its own crash window honestly (its Decision 3: the handler executes before its record is
durable) and defers the fix here; the `w-mcp-projection` carve-out found the same gap on the
commit path and marked its premise row `Commit-boundary contract` **UNVERIFIED — PREREQUISITE**.
No landed API exposes a durable pre-execution intent, a stable invocation/idempotency ID, or a
queryable receipt — so today a crashed process cannot ask the store "did my commit land?" or "did
my effect run?", and any answer it invents risks either lying ("not committed" when it committed)
or silently re-executing a paid, non-idempotent effect.

The daemon already refuses zero refs at its boundary (`parseRef`, with the reasoning written out
at `handlers.go:155-165` — which itself calls the store asymmetry "a real M1 defect, recorded as
a carry-forward"). The charter's objection stands: that is enforcement in the wrong layer. An
embedded writer — a test, a future broker, any Go consumer — bypasses it entirely. Authority as
enforcement, not convention, was exactly the rationale for the ratified single-writer lock; this
item applies the same principle to what the writer is allowed to write.

## Premises (hard constraints — each verified in the Premise Verification Log)

- **P1 — the defect is the CLASS, not the field.** `store.Commit` renders every ref via
  `HashRef.String()`, which returns `""` for the documented-invalid zero value, and no write path
  validates (V4, V6). The fix must cover every ref column on every write path
  (`Commit`, `PutObject`, `PutWorld`, `SetRegistryHead`, `SelectHead`, `PutVerifyResult`), not
  re-patch `PrevEntryHash` alone — the systemic-issue rule: one unified fix, no future bugs in
  this area.
- **P2 — the genesis lenience stays exactly one field wide.** The kernel DEFINES the zero
  `ObservedHead` as the genesis marker (`Commit`'s own doc comment), the daemon's
  `parseGenesisRef` mirrors it, and `TestGenesisRefLenienceIsExactlyOneField` pins the blast
  radius (V8). This design keeps that single exception byte-for-byte and must leave that test
  green: `ObservedHead` is compared, never persisted — every ref that is *written* must be valid.
- **P3 — kernel changes are ratification-class.** This item reopens the landed M1 `host/store`
  (behavioral change to `Commit` et al.) and makes the first-ever change to `schema.sql`
  (byte-unchanged through M2 — additive journal table). Both are presented as named ARMS with
  trade-offs (Decisions 1 and 3) for Mark's attended ratification, the `w-worldd-m2` arm-A
  precedent. The sprint does not start before that ratification.
- **P4 — schema migration is designed, not assumed.** The schema is applied on every writer
  `Open` via `db.Exec(schemaSQL)` with `CREATE TABLE IF NOT EXISTS`, and read-only handles skip
  schema application (V11). An **additive new table** therefore reaches every existing store on
  its next writer open with zero data migration; the DDL of **existing** tables is byte-unchanged
  (SQLite cannot `ALTER TABLE … ADD CHECK`, so retrofitting constraints onto `log_entries` would
  fork old and new stores' definitions — rejected; Go-layer validation is the single enforcement
  for existing tables, and the new journal table carries SQL constraints from birth).
- **P5 — already-written poisoned rows get a stated policy, not silence.** Decision 2:
  reject-on-read stays (it is correct fail-loud behavior), a bounded detection sweep names the
  holes, and repair is ruled out in writing — with the immutability argument, not a shrug.
- **P6 — determinism and zero-cloud are untouched.** No wall-clock enters any canonical payload
  (intent/outcome objects carry caller-supplied logical time, the broker-doc P5 discipline); the
  journal adds zero Go dependencies; `TestDaemonDependencyAllowlist` and the store's existing
  import surface are unchanged.
- **P7 — every wait and allocation is bounded, and the KERNEL owns the bound** (Standing Rule 6,
  in the design as in the shell). Sharpened in quorum round 2 (`gpt5-6-sol`): a caller-supplied
  limit is not a bound. Every paged API requires `1 <= limit <= Max…` against the kernel constants
  `MaxPendingIntentsPage` / `MaxIntegrityScanPage` and returns `InvalidLimitError` otherwise; the
  detection sweeps are paginated by per-table keyset cursors (Decision 2); the daemon's startup
  integrity check pages to completion or to a fixed total-row/time budget, and on exhaustion emits
  `integrity_scan_incomplete` with its continuation cursor — never a message implying a full scan;
  crash-injection tests poll a
  **captured pid** under a `date +%s`-style deadline (the iter-24 scar, honoured in test code:
  `kill -0 "$pid"` semantics via `os.Process` + deadline, never a name-pattern poll).
- **P8 — the verify gates extend without edits.** Verified from the script this session (V14):
  `verify_ail.sh`'s manifest keys required identities to `world/` modules only, its exact totals
  (4 identities / 14 tests) are `world/`-scoped, and the module count is dynamic — the new sketch
  moved the sweep 9 → 10 modules with totals unperturbed (V15, live run). **Honest limitation,
  stated**: Leg 2 executes `ailang test` on `world/` only, so the sketch's 25 named tests are NOT
  gate-run in CI; the sprint's verify_commands run them explicitly (the `effectbroker.ail`
  M3.A pattern), and the sketch's contracts ARE swept by Leg 1's per-module ai-check.

## High-Impact Decisions

| Decision | Why High Impact | Chosen By | Deadline | Change Cost |
|----------|-----------------|-----------|----------|-------------|
| Validate every written ref at the kernel write path (ARM V1), not on read (ARM V2) | Ends the CF-B-2 class for every field and every caller; the daemon's boundary check becomes defense-in-depth instead of the only wall | **ratification (Mark)** | SD.A | medium |
| Existing poisoned rows: reject-on-read stays + bounded detection sweep + documented-unrecoverable | Repair = rewriting committed history; a silent skip = grade laundering; detection makes the hole operable | this doc (D2) | SD.A | low |
| Durable intent/outcome journal as an additive SQLite table + content-addressed payload objects (ARM J1), not a zero-schema registry-head chain (ARM J2) | One substrate answers the broker journal AND the commit receipt; SQL gives atomic append, bounded pending-scan, and at-most-one-outcome enforcement | **ratification (Mark)** | SD.B | high |
| `Commit` gains an optional stable invocation ID, binds in-tx to the intent's commit-defining fields, and writes the receipt outcome in the SAME transaction | Makes "not-started versus committed" atomic by construction AND makes a receipt for an operation other than the recorded intent unrepresentable (round-1 objection 1) — the CF-D-2 contract, achievable without two-phase commit | **ratification (Mark, same packet)** | SD.B | high |
| Receipt law is three-state (`not-started` / `indeterminate` / `resolved`) with the never-lie rule, Z3-proven in the sketch | The exact contract `w-mcp-projection` AC13 is blocked on; a definitive "not committed" while the outcome is unknown is structurally unexpressible | sketch (Appendix A) | SD.B | high |
| Recovery NEVER auto-re-executes an indeterminate intent; retries require explicit per-handler reconciliation | The silent-duplicate-execution risk (a paid `Model.Infer` run twice) is the named failure mode of half (ii) | this doc (D5) | SD.C | high |
| Crash-injection = real subprocess kills at named points (the `writer_lock_test` pattern), not in-process simulation | An in-process "crash" cannot prove SQLite's rollback + journal atomicity across process death | this doc (D6) | SD.C | medium |
| Journal cost enters the day-1 baseline (2 new benchmark names in the hardcoded manifest) | The §13.5 budget must price the new write-path tax before anything depends on it | this doc (D7) | SD.C | low |

### Design Freeze (the sprint must not renegotiate these)

- [ ] Every ref field the store WRITES is validated before any transaction begins; a zero or
  unparseable ref returns a structured `InvalidRefError{Op, Field}` and leaves the store
  untouched. `ObservedHead` remains the single, kernel-defined zero-legal COMPARED field (P2).
- [ ] Existing tables' DDL in `schema.sql` is byte-unchanged; the ONLY schema change is the
  additive `journal` table (Decision 3's DDL), applied via the landed `CREATE TABLE IF NOT
  EXISTS` mechanism on writer open.
- [ ] Poisoned rows already written are never rewritten, deleted, or silently skipped; the
  detection sweep reports them; reads keep failing loudly with the current structured errors.
- [ ] The journal is append-only with gapless sequence numbers (Law 4); intent and outcome
  payloads are content-addressed store objects (`world/*` semantic IDs per Decision 3); at most
  one intent and one outcome per invocation ID, enforced by the schema.
- [ ] `GetReceipt` answers in exactly the three states of the sketch's `receiptState` law and can
  never answer `not-started` while a durable intent exists (`mayReportNotStarted`, Z3-proven).
- [ ] `Commit` with an `InvocationID` loads the durable intent inside its own transaction and
  compares the commit-defining fields (invocation ID, planned `NextWorld.Ref`,
  `Entry.EntryHash`, `ObservedHead`, `PrevEntryHash`, `TransitionFn`, `TransitionRef`,
  `Interpreter` — the round-2 list) against the actual request BEFORE any mutation; any
  mismatch returns structured `InvocationMismatchError` with the store untouched; a `Commit` of
  an already-resolved ID returns the existing receipt idempotently ONLY when those fields match,
  else the mismatch error (`intentBindsCommit`, Z3-proven — quorum round 1).
- [ ] Detection is per-table with keyset cursors: `ScanUnreadableLog(fromIndex, limit)` over the
  integer log index, `ScanUnreadableWorlds(afterRef, limit)` over the ref TEXT in lexicographic
  order — never `OFFSET`, never implicit rowids (quorum round 1).
- [ ] Every paged API enforces `1 <= limit <= Max…` against a KERNEL constant
  (`MaxPendingIntentsPage`, `MaxIntegrityScanPage`) and returns `InvalidLimitError` outside it;
  the daemon's startup integrity check pages to completion or to a fixed total-row/time budget and
  emits `integrity_scan_incomplete` with a continuation cursor when the budget is exhausted — a
  truncated scan never reports as a clean one (quorum round 2).
- [ ] No recovery path dispatches, re-executes, or retries anything. The kernel surfaces
  indeterminate intents; resolution is an explicit consumer action gated by `retryAllowed`
  (Z3-proven).
- [ ] The M2 REST route table, error envelope, D7 timeout constants, and CLI verb set are
  byte-untouched; the only daemon delta is the bounded startup integrity warning (Decision 2).
- [ ] `host/replay/**`, `host/{hashref,canon,archive,registry}/**`, `world/**`, `go.mod`,
  `go.sum`: byte-unchanged.
- [ ] The sketch `sketches/storejournal.ail` is the frozen law; the Go mirror is pinned by a
  drift test over every law and every `tests[]` boundary case (the `worlddapi`/`effectbroker`
  precedent).

---

## Decision 1 — Validate-on-write at the kernel (ratification ARMS V1/V2)

**The measured defect class (all first-party this session, V2–V4).** `Commit` renders seven ref
fields to TEXT with `HashRef.String()` — `EntryHash`, `TransitionFn`, `Interpreter`,
`PrevEntryHash`, `TransitionRef`, and the world row's `Ref`/`StateRoot`/`LogHead` — and inserts
whatever it gets; `PutObject`/`PutWorld`/`SetRegistryHead`/`SelectHead`/`PutVerifyResult` share
the pattern. The zero value renders as `""`, satisfies every `TEXT NOT NULL`, and every reader
(`GetLogEntry`, `GetWorld`, `SelectedHead`, …) parses with `hashref.Parse`, which — correctly —
rejects it. Write-anything + read-strictly is the asymmetry; three distinct poison shapes were
reproduced live (zero prev → unreadable entry; zero transitionFn → unreadable entry; zero world
logHead → unreadable world, with the selected head pointing at it), and the controller's
independent eight-field matrix (V23) closed out the class measurement: **all eight** zero-field
commits return `commit_err=<nil>`; **seven** produce an unreadable row, and the eighth —
`NextWorld.Ref` — commits and reads back *fine*, the empty-string world ref becoming the
selected head (degenerate but readable).

**ARM V1 (recommended) — reject on write, everywhere.** One unexported helper,
`validateRef(op, field string, ref hashref.HashRef) error`, returning a structured
`*InvalidRefError{Op, Field}`; applied to **every** ref a write path persists, before any
transaction begins (the `verifyObject` placement precedent). `Commit` additionally validates the
object refs it inserts (`Hash` is already content-verified; `InterfaceHash` is not — included).
The check is `hashref.Parse(ref.String())` round-trip or equivalently `ref.IsZero()` plus
canonicality — implementation may use the cheaper form, but the TEST asserts rejection of both
the zero value and a hand-constructed non-canonical value if the type ever permits one. The
daemon's boundary check is retained unchanged and becomes defense-in-depth (two independent
layers, the F-floor pattern). The stale `handlers.go` comment block that records the asymmetry as
an open defect is corrected in the same milestone to cite the kernel fix (the CF-B-1
stale-comment precedent).

**ARM V2 (rejected, presented for the ratification record) — support the zero on read.** Make
`GetLogEntry` (and symmetrically `GetWorld` etc.) return the zero `HashRef` for empty TEXT. Why
rejected: (a) it blesses a value the type documents as invalid and every other layer rejects;
(b) M1's own genesis convention seeds entry 0's `PrevEntryHash` from the genesis world's
`LogHead` — a real content address (V10) — so a zero prev is *never* legal in a written entry and
lenience would legitimize only corruption; (c) read-lenience cannot distinguish "legal absent"
from "torn write"; (d) it must either bless all eight fields' zeros (a world whose `StateRoot` is
"" becomes readable-but-meaningless) or leave the class inconsistently half-fixed; (e) — the
sharpest, measured, argument (V23) — the eighth field is fatal to V2 on its own terms: a zero
`NextWorld.Ref` commits AND reads back without error, the empty-string ref becoming the selected
head, so a read-side fix never even *observes* a failure to be lenient about — that poison shape
is only closable at the write. The charter
row phrased the kernel decision as exactly this fork; the recommendation is V1.

**Why is this not a package? (S3)** Validation of the kernel's own write path cannot live outside
the kernel: an external validator is precisely the boundary patch this item exists to end — any
embedded caller bypasses it by construction. The kernel grows by one helper, one structured error
type, and zero new concepts.

## Decision 2 — Already-written poisoned rows: reject-on-read + detection, repair ruled out

Constraint: a fix that only prevents new bad writes leaves existing holes unreadable forever. The
policy, stated explicitly:

- **Reject-on-read STAYS.** The current structured parse errors are correct fail-loud behavior; a
  read that silently skipped or zero-filled a hole would hide corruption — grade laundering at
  the kernel layer.
- **Repair is ruled out, with the argument written down.** (a) There is no canonical `LogEntry`
  encoding anywhere in the repo — `host/canon` canonicalizes source text only, and nothing
  derives or verifies `EntryHash` against header bytes (negative existence, V12) — so the "true"
  prev value of a poisoned row is not recomputable from the row itself; (b) even where the chain
  invariant suggests the value (entry N's prev "should" be entry N−1's `EntryHash`), rewriting a
  committed row violates the append-only guarantee clause 1 is built on, and any chained
  reference to the row's current bytes would be silently invalidated. Rewriting history to make
  it prettier is the one move this kernel must never learn.
- **Quarantine (moving/flagging rows) is the same move** — a mutation of committed log state —
  and is rejected for the same reason.
- **Detection is the deliverable — one sweep per table.** Round 1 (gemini-3-1-pro) rejected the
  earlier single API `ScanUnreadable(fromIndex, limit)`: `log_entries` is ordered by an integer
  `index`, but `worlds` is keyed by a content-addressed `TEXT` ref, so one integer `fromIndex`
  cannot safely paginate `worlds` without unstable `OFFSET`s or implicit SQLite rowids — which a
  `VACUUM` can renumber, violating the determinism axiom. **That API is replaced** (explicit
  substitution) by the reviewer's pair: `store.ScanUnreadableLog(fromIndex int64, limit int)`
  pages `log_entries` by its integer `index`; `store.ScanUnreadableWorlds(afterRef string,
  limit int)` pages `worlds` by its ref TEXT primary key with an explicit keyset cursor,
  **ordered lexicographically on the ref TEXT** (`WHERE ref > afterRef ORDER BY ref LIMIT n`) —
  stated, not implied: lexicographic TEXT order is stable and deterministic, uses no `OFFSET`
  and no rowid, and is resumable across writes and `VACUUM`. Both sweeps are bounded, read-only,
  attempt the standard parse on each row, and return the failing index/ref with the failing
  field.
- **Kernel-enforced page bounds, and a startup scan that cannot silently truncate** (round 2,
  `gpt5-6-sol`, applied verbatim under the narrow-refinement carve-out). *Replaced claim, stated
  rather than deleted*: round 1's text said the sweeps and `PendingIntents` were bounded because
  they "respect the caller's limit" — the reviewer correctly refuted that, since *"merely stopping
  at the supplied limit does not bound allocation or query work"* when the limit itself is
  arbitrary and zero/negative/oversized behaviour is undefined. The contract is now:
  - Kernel constants **`MaxPendingIntentsPage`** and **`MaxIntegrityScanPage`**. Every paged API
    (`PendingIntents`, `ScanUnreadableLog`, `ScanUnreadableWorlds`) requires
    `1 <= limit <= Max…` and returns a structured **`InvalidLimitError`** otherwise. The kernel,
    not the caller, owns the ceiling.
  - `ailang-worldd serve`'s startup integrity check **pages with a fixed page size until either
    completion or a fixed total-row/time budget is reached** — it is not "one bounded sweep",
    which could silently miss every hole after page 1.
  - If the budget is exhausted, the daemon emits a **distinct `integrity_scan_incomplete`
    warning carrying the continuation cursor and the counts scanned so far — never a message
    implying the store was fully scanned.** A partial scan that reads as a clean bill of health
    is the grade-laundering failure mode this whole item exists to refuse.

  It does not refuse to serve either way: a store with one historic hole still serves every
  readable row, and the poisoned index already fails loudly on access. The REST range-read behavior is deliberately UNCHANGED — a range that
  crosses a hole keeps failing as today, because a skip-and-continue response shape would be a
  frozen-contract change that hides exactly what the sweep exists to surface.
- **Population honesty**: no committed `.db` fixture exists in the repo (V20), the daemon has
  refused the REST path since M2.B, and every landed test uses real hashes — so the *known*
  poisoned population is zero and this policy is expected to matter for embedded misuse and
  future forensics, not for any store we ship. That is an expectation, not a proof; the sweep is
  what makes it checkable per-store.

## Decision 3 — The durable intent/outcome journal (ratification ARMS J1/J2)

**The shape (reviewer's own terms, both sources).** A caller that is about to do something whose
consequence must survive its own death — dispatch an external effect, or drive a commit whose
caller may vanish mid-flight — first appends a durable **intent**; when the consequence is known,
it appends exactly one **outcome** bound to the same invocation ID. Between the two, the
invocation is **indeterminate** — a first-class, queryable state, never a guess.

**ARM J1 (recommended) — additive `journal` table + content-addressed payload objects.**

```sql
-- Appended to schema.sql (additive; existing tables byte-unchanged — P4):
CREATE TABLE IF NOT EXISTS journal (
    seq           INTEGER PRIMARY KEY,             -- gapless, assigned in-tx (Law 4)
    kind          TEXT    NOT NULL CHECK (kind IN ('intent','outcome')),
    invocation_id TEXT    NOT NULL CHECK (invocation_id <> ''),
    object_ref    TEXT    NOT NULL CHECK (object_ref <> ''),
    UNIQUE (invocation_id, kind)                   -- ≤1 intent, ≤1 outcome per invocation
);
```

The row is the durable, ordered index; the payload (what was intended / what happened) is an
ordinary content-addressed object (`world/journal-intent/v1` / `world/journal-outcome/v1`)
written **in the same transaction** as its row — so a journal row without its object, or an
object without its row, is unrepresentable. The CHECK constraints give the new table the
write-validity floor at the SQL layer from birth (defense-in-depth beneath the Go validation —
new tables can carry constraints; existing ones cannot, P4). Gapless `seq` is assigned as
`max(seq)+1` inside the transaction — safe precisely because the ratified single-writer lock
makes the store's writer unique across processes.

API (all on `*store.Store`, all bounded):

| Method | Contract |
|--------|----------|
| `AppendIntent(id string, intent JournalIntent) (seq int64, ref HashRef, err error)` | idempotent by ID: same ID + byte-identical intent object → returns the existing row (no-op); same ID + different bytes → structured `DuplicateInvocationError`. Validation per Decision 1 applies |
| `AppendOutcome(id string, outcome JournalOutcome) (seq int64, ref HashRef, err error)` | structured error if no intent exists for `id` (the sketch's `corrupt` arm is unrepresentable by construction) or an outcome already does |
| `GetReceipt(id string) (Receipt, ok bool, err error)` | the three-state law: intent? outcome? → `not-started` / `indeterminate` / `resolved`, mirroring `receiptState` (drift-tested against the sketch) |
| `PendingIntents(limit int) ([]PendingIntent, error)` | intents without outcomes, oldest first. `1 <= limit <= MaxPendingIntentsPage` — outside that range returns `InvalidLimitError`; the KERNEL owns the ceiling, not the caller (round-2 fix, P7) |

**ARM J2 (rejected, presented for the record) — zero-schema: objects + a registry-head chain**
(the broker doc's P2 pattern: journal head in `epoch_registry_heads`, each record embedding the
previous ref). Why rejected: (a) atomicity of object-write + head-move needs a new store method
anyway, so J2 saves no kernel-API surface — it only trades a table for a convention; (b)
"intents without outcomes" becomes either an unbounded chain walk or a `semantic_id` table scan
with no defined order — both worse under Standing Rule 6 than one indexed SQL query; (c) at-most-
one-outcome-per-invocation cannot be *enforced*, only asserted; (d) per-invocation names would
make `epoch_registry_heads` an unbounded per-request namespace, stretching a table whose landed
tenants are named registries (`world/epoch-registry/v1`, and the broker doc's planned
`world/approvals/v1`). The zero-schema instinct is right for *content* — which is why payloads
stay content-addressed objects — and wrong for the *index* that recovery must scan.

**Growth**: the journal is append-only like the log; a journal row is two TEXT refs and two small
TEXT fields (≈150 bytes + payload objects). No compaction in this item (OD4).

**Why is this not a package? (S3)** The journal's entire value is that it is in the SAME SQLite
file and the SAME transactions as the store state it witnesses — "did my commit land?" is only
answerable atomically from inside the commit's own transaction. A package journal beside the
store would reintroduce the exact torn-write window it exists to close. The *consumers* —
per-handler reconciliation logic, broker recovery policy — are NOT kernel and stay out
(Decision 5, M3, packages).

## Decision 4 — Commit receipts: the atomic "not-started versus committed" contract (CF-D-2)

`store.Commit`'s struct gains one optional field: `InvocationID string` (zero value = today's
behavior, byte-compatible with every landed caller). When set:

1. The caller appends the intent first — `AppendIntent(id, intent)` where the intent object
   records the planned `NextWorld.Ref`, `Entry.EntryHash`, `ObservedHead`, **`PrevEntryHash`,
   `TransitionFn`, `TransitionRef`, `Interpreter`**, and caller-supplied logical time. Those
   seven plus the invocation ID are the **commit-defining fields**, and their canonical encoding
   is fixed by the intent codec's golden bytes. *The four entry refs were added in quorum round 2
   (`gemini-3-1-pro`), applied verbatim under the narrow-refinement carve-out; the round-1 list of
   three is REPLACED, not extended silently.* The reviewer's ground is exact and load-bearing:
   **premise V12 records that the kernel does not verify `EntryHash` against the entry's
   contents**, so with only `NextWorld.Ref`/`EntryHash`/`ObservedHead` compared, a caller could
   append intent A and then `Commit` an entry with mutated `PrevEntryHash`/`TransitionFn`/
   `Interpreter` **without changing `EntryHash`** — and the receipt would falsely claim the
   original intent succeeded. That is the semantic-equality contract of round 1 defeated through
   the gap round 1 did not close. Until a canonical entry encoding exists (Deferred Scope), the
   binding must name every unverified ref field explicitly. **This is the
   acceptance point**: before it, cancellation guarantees no durable mutation; after it, the
   only truthful answers are indeterminate or resolved.
2. **`Commit` binds to the intent before it mutates anything** (quorum round 1, gpt5-6-sol —
   this step replaces the round-0 text, which established atomic *presence* of commit + outcome
   but not semantic *equality* among the intent, the actual `Commit` arguments, and the
   outcome): inside its existing single transaction, `Commit` loads the durable intent for `id`
   and compares the canonical commit-defining fields — planned `NextWorld.Ref`,
   `Entry.EntryHash`, `ObservedHead`, `PrevEntryHash`, `TransitionFn`, `TransitionRef`,
   `Interpreter`, invocation ID — against the actual request (the round-2 field list; see step 1
   for why the shorter list was unsound). Any mismatch
   returns a structured `InvocationMismatchError{ID, Field}` and **leaves the store untouched**;
   a receipt can therefore never claim resolution for an operation different from the recorded
   intent (`intentBindsCommit`, sketch LAW 6, Z3-proven). Reuse of an **already-resolved** ID is
   specified, not left open: `Commit` returns the existing receipt idempotently when the fields
   match, and the mismatch error when they do not — a resolved invocation is never silently
   re-resolved. A set `InvocationID` with **no** durable intent is likewise a structured error
   (intent-first is the protocol, step 1).
3. Only then does `Commit` write the outcome row + outcome object **inside that same
   transaction** (between today's steps 5 and 6). Consequence, by construction: *the commit is
   durable if and only if its receipt is* — there is no window in which the world advanced but
   the receipt is missing, or vice versa. This is what makes the CF-D-2 contract achievable
   without two-phase commit: the store IS the receipt store, so atomicity is one SQLite
   transaction, not a distributed protocol.
4. A crashed caller (or an HTTP handler whose context expired mid-commit — the
   `w-mcp-projection` scenario) recovers by `GetReceipt(id)`: `resolved` → the commit landed,
   outcome carries the new head; `indeterminate` → the intent is durable but no outcome is —
   and for the commit consumer specifically, reconciliation is **deterministic and total**:
   `Commit`'s atomicity means "no outcome" ⇒ the transaction did not commit ⇒ the store is
   untouched; the reconciler verifies (does `NextWorld.Ref` from the intent exist? does the
   entry?) and appends a reconciling outcome recording "not executed". No guessing, no lying,
   and the *mechanism* — not just the claim — is what the crash-injection tests exercise.
5. The caller must never surface a definitive "not committed" from a bare error/cancellation
   while the receipt is unqueried — the sketch's `mayReportNotStarted` law, and AC5's mutation.

This is the reviewer's contract from `w-mcp-projection` Decision 6, made concrete: the stable
invocation/idempotency ID is the journal's `invocation_id`; the "defined point of no return" is
the intent append; the "queryable durable receipt" is `GetReceipt`. When this lands, that doc's
premise row `Commit-boundary contract` can flip to VERIFIED citing `store.AppendIntent` /
`Commit.InvocationID` / `GetReceipt` and the SD.C crash tests — un-blocking prerequisite 3 of
the clause-6 item.

**Deliberate non-goal**: no REST/CLI projection of receipts. The M2 route table is frozen;
projecting receipts is the clause-6 item's business once it unblocks (its AC13 already specifies
the test shape). This item's consumers are embedded (tests now, broker + projection later).

## Decision 5 — Recovery law: indeterminate fails closed; retries require reconciliation

The kernel surfaces facts; it never acts on them:

- `PendingIntents` + `GetReceipt` are the ONLY recovery surface. There is no `Retry`, no
  `Redispatch`, no auto-resolution anywhere in `host/store`.
- The consumer contract, frozen here for M3 (and proven over a probe consumer in SD.C): on
  finding an indeterminate intent, a consumer must surface a structured
  **`IndeterminateEffectError`** (the reviewer's name) to its caller and may proceed to
  re-execution ONLY after an explicit, per-handler reconciliation establishes the effect did not
  happen — the sketch's `retryAllowed(indeterminate, reconciled)` law, Z3-proven.
- **Per-handler reconciliation table (the frozen contract M3 implements; normative for future
  handlers):**

| Consumer / handler | Reconciliation rule | Auto-retry after reconcile? |
|--------------------|--------------------|----------------------------|
| `store.Commit` (this item) | deterministic: query the store for the intent's planned refs; atomicity makes absence proof of non-execution | yes — reconcile is total, so retry is safe |
| `FS.Write` (M3) | compare target content hash against the intent's payload hash | yes if absent/differs; write is idempotent by content |
| `Git.Commit` (M3) | probe the repo for the intended commit (ref/sha) | yes if absent |
| `Model.Infer` (M3) | **none exists** — a paid inference leaves no reliably queryable trace | **never** — resolution is explicit (operator or policy), recorded as an outcome of kind `abandoned` |
| `Human.Approve` (M3) | query the approval-request object's presence | yes if absent |

The table is here, in the kernel item, because it is the *contract* half of half (ii): under
answer (b) this doc carries it in full; under answer (a) the implementation of the four handler
rows merges into M3 while the substrate and the law stay here (see the gating section).

## Decision 6 — Crash-injection method (real process death, named kill points)

The `writer_lock_test.go` subprocess pattern, extended: a helper binary (test-built, the M2.C
real-subprocess precedent) opens the store and advances to a **named stop point** —
`after-intent`, `after-external-effect` (the probe consumer touches a scratch file),
`mid-commit-before-outcome` (arranged by helper logic, not sleeps), `after-outcome` — then either
exits or is SIGKILLed by the parent. The parent then reopens the store and asserts the receipt
law: which state, what `PendingIntents` returns, that the probe effect's presence/absence matches
the reconciliation rule, and that nothing was auto-re-executed (counting probe). All waits poll
the **captured** `os.Process` under a deadline; no name-pattern polling, no unbounded loop (P7 —
the iter-24 rules, encoded in test code).

## Decision 7 — Cost: the journal tax enters the day-1 baseline

Two new benchmarks in `host/daemon/bench_test.go` (the harness home), two new names in
`scripts/bench_worldd.sh`'s **hardcoded manifest**, rows in `bench/BASELINE.md` re-measured in
ONE invocation (the iter-22 discipline): `BenchmarkJournalAppend` (intent append, the new
write-path floor) and `BenchmarkCommitWithReceipt` (full commit + in-tx receipt vs the landed
`BenchmarkStoreCommit` — the marginal receipt tax must be visible as a diffable pair). Initial
targets to validate and record, honest about enforcement (CI asserts mechanism, thresholds live
on the dev rig baseline, the D6 precedent): journal append p95 ≤ 10 ms; commit-with-receipt p95
within +20% of bare store-commit p95.

## Decision 8 — What stays OUT (and where it goes)

No broker, no handlers, no REST/CLI surface, no repair tooling, no `EntryHash`-recomputation
integrity layer (a canonical entry encoding does not exist to verify against — V12 — and
inventing one is its own log-format-adjacent design), no WAL/fsync tuning, no schema version
pragma, no journal compaction. Each has a named destination in Deferred Scope.

---

## The gating relationship — what shrinks if Mark answers (a)

Charter item 4 (`w-effect-broker-m3`) is PARKED on a binary question whose recorded default is
**(b)** (keep M3's scope, weaken its claim, journal lands here). As of this iteration there is
still no answer. **This doc is designed for default (b) and carries both halves.** The cut, if
the answer is (a) — half (ii) merges into M3 — is along milestone lines, not a rewrite:

- **SD.A (half i) and SD.B (substrate + commit receipts) stay here unchanged.** The queue row's
  iter-23 approximation "this item shrinks to CF-B-2 alone" predates iter-24's fold-in of the
  commit-boundary prerequisite (CF-D-2): the commit receipt is consumed by `w-mcp-projection`
  independently of the broker, and its atomicity argument (Decision 4) lives inside
  `store.Commit` — it cannot move to M3 without M3 reopening the store, which is exactly what
  M3's own P4 forbids. This doc is the durable record of that precision.
- **SD.C splits**: the generic parts (commit-path crash injection, probe-consumer recovery proof,
  the frozen law + reconciliation table) stay; the four *handler-specific* reconciliation
  implementations and broker-side crash injection merge into M3's milestones, which would then
  depend on SD.B's landed substrate. M3's Decision 3 "honest ordering limitation" paragraph and
  its Deferred-Scope journal row are superseded either way — under (a) by M3's own revision,
  under (b) by this item.
- Either answer can be ratified in the SAME attended comment as this doc's kernel arms
  (Decision 1 V1/V2, Decision 3 J1/J2, Decision 4's `Commit` extension) — one human touchpoint
  unparks two items.

## Milestones (each independently CI-green and mergeable; ~1.5–2d total)

### SD.A — Repro fixture + validate-on-write + detection sweep (~0.5d) — half (i), in scope under BOTH answers

- **files**: `host/store/durability_repro_test.go` (~150 — **the committed repro fixture, the
  first deliverable per the charter's ghost-close rule**: reproduces this doc's V2/V3 probe
  shapes and asserts the post-fix behavior — structured `InvalidRefError` per field, store
  untouched by diffing row counts + head before/after), `host/store/store.go` (+~70 —
  `validateRef`, `InvalidRefError`, application at every write path per Decision 1),
  `host/store/scan.go` (~110 — `ScanUnreadableLog` + `ScanUnreadableWorlds`, bounded,
  read-only, keyset pagination per Decision 2) + `scan_test.go` (~140 — fixture store with
  hand-inserted poisoned rows via raw SQL in BOTH tables, proving detection without the public
  API being able to create one anymore, plus cursor-resumability cases for each sweep),
  `host/daemon/daemon.go` (+~20 — bounded startup sweeps, one per table, + named warning; zero
  route/constant changes) + `daemon_test.go` (+~40),
  `host/daemon/handlers.go` (comment-only: the :155-165 block updated to cite the kernel fix)
- **acceptance_checks**: AC1–AC4, AC10 (daemon part), AC12; `TestGenesisRefLenienceIsExactlyOneField`
  and every landed `host/store` + `host/replay` suite green unmodified
- **verify_commands**: `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh` ·
  `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh`
- **ci_green_boundary**: no journal exists yet; nothing depends on SD.B

### SD.B — Journal substrate + commit receipts (~0.5–0.75d) — needs the ratified arms

- **files**: `host/store/schema.sql` (+~10 — the Decision 3 DDL, additive only),
  `host/store/journal.go` (~240 — types, deterministic codec for intent/outcome objects
  (golden-bytes test, the broker-doc P5 discipline), `AppendIntent`/`AppendOutcome`/
  `GetReceipt`/`PendingIntents`, gapless seq in-tx), `host/store/store.go` (+~45 —
  `Commit.InvocationID` + in-tx intent load/field-compare with `InvocationMismatchError`
  (Decision 4 step 2) + in-tx outcome write), `host/store/journal_test.go` (~320 — receipt
  tri-state drift test against the sketch's `receiptState`/law tests, **intent-binding drift
  test mirroring `intentBindsCommit` on every `tests[]` boundary case** — `AppendIntent(id, A)`
  then `Commit(id, B)` → mismatch error + store untouched; resolved-ID re-commit idempotent
  only on field match — idempotent re-append,
  duplicate-ID rejection, gapless seq, bounded `PendingIntents`, golden bytes, migration: open a
  pre-journal fixture store → table appears, existing DDL byte-identical via
  `sqlite_master` comparison)
- **acceptance_checks**: AC5–AC8, AC9 (sketch), AC13 (migration), AC15 (intent binding)
- **verify_commands**: both gates + explicit
  `cd design_docs && /tmp/ailang-v0300/ailang ai-check -timeout 5s sketches/storejournal.ail &&
  /tmp/ailang-v0300/ailang test --format json sketches/storejournal.ail` (P8's honest
  limitation: Leg 2 does not run sketch tests — the sprint runs them explicitly)
- **ci_green_boundary**: journal complete and consumable; crash proofs are SD.C's

### SD.C — Crash injection + recovery proof + bench + close-out (~0.5d)

- **files**: `host/store/crash_test.go` (~260 — Decision 6's helper + named kill points +
  captured-pid bounded waits), `host/store/recover_test.go` (~120 — probe consumer:
  `IndeterminateEffectError` surfacing, `retryAllowed` gate, commit-path deterministic
  reconciliation, Model.Infer-style never-retry proof over the probe), `host/daemon/bench_test.go`
  (+~50), `scripts/bench_worldd.sh` (+2 manifest names), `bench/BASELINE.md` (all rows
  re-measured in ONE invocation), `README.md` (+~8)
- **acceptance_checks**: AC10–AC11, AC14; doc → `implemented/` with every box checked
- **verify_commands**: full sweep — both gates + `./scripts/bench_worldd.sh --smoke` + the
  explicit sketch test run
- **ci_green_boundary**: LANDS the item

**Pre-committed overflow cut line**: if the sprint runs past ~2d, the cut is the **probe-consumer
recovery proof's handler-table breadth** (`recover_test.go` keeps the commit-path + never-retry
cases and defers the FS/Git-shaped probe cases to M3, which owns those handlers anyway under
either answer). SD.A and the journal atomicity + crash-injection core are NOT cuttable — they are
the item's substance.

## Files to Create/Modify (aggregate)

| File | Est. LOC | Change |
|------|---------:|--------|
| `design_docs/sketches/storejournal.ail` | 163 | **already created + verified with this doc, re-verified after the round-1 revision** (7/7 Z3-proven, 25 named tests + 7 contract properties = `passed_tests` 32) |
| `host/store/store.go` | +~115 | validateRef + InvalidRefError + write-path application; `Commit.InvocationID` + in-tx intent binding (`InvocationMismatchError`) + in-tx outcome |
| `host/store/schema.sql` | +~10 | additive `journal` table (existing DDL byte-unchanged) |
| `host/store/journal.go` | ~240 | new: journal types, codec, four API methods |
| `host/store/scan.go` | ~125 | new: bounded `ScanUnreadableLog` + `ScanUnreadableWorlds` (keyset cursors, Decision 2) + the `Max…` page constants and `InvalidLimitError` guard (round 2) |
| `host/store/durability_repro_test.go` | ~150 | new: **the committed CF-B-2 repro fixture** |
| `host/store/{scan,journal,crash,recover}_test.go` | ~840 | new |
| `host/daemon/daemon.go` + `daemon_test.go` | +~80 | startup sweeps (one per table) paging to completion or a fixed budget + per-hole warning + the `integrity_scan_incomplete` truncation warning and its multi-page test (round 2); no route/constant change |
| `host/daemon/handlers.go` | ~0 (comment) | stale CF-B-2 carry-forward comment corrected |
| `host/daemon/bench_test.go` | +~50 | two journal benchmarks |
| `scripts/bench_worldd.sh` | +2 | manifest names |
| `bench/BASELINE.md` | ~8 | journal rows (full re-measure) |
| `README.md` | +~8 | operator note (startup sweep warning) |

Estimated ~1,750 LOC. **Byte-unchanged**: `host/replay/**`, `host/{hashref,canon,archive,registry}/**`,
`world/**`, `cmd/**` (CLI verbs untouched), `scripts/verify_ail.sh`, `scripts/verify_go.sh`,
`go.mod`, `go.sum`, `.github/**` (both CI jobs already sweep the new code; bench-smoke already
runs the manifest script).

## Conflict Surface (every landed behaviour this design could collide with)

- **vs the ratified single-writer lock (M2.A).** The journal writes through the SAME writer
  handle inside the store; no second write path, no new `Open` variant. Gapless in-tx `seq`
  assignment is CORRECT only because the writer lock makes the writer unique — stated as a
  dependency, not an accident. A future multi-writer design would have to redesign Law 4;
  that is the tripwire's purpose.
- **vs the frozen M2 REST/CLI surface.** Zero route changes, zero verb changes, zero error-
  envelope changes, D7 constants untouched; the range-read's fail-on-hole behavior deliberately
  preserved (Decision 2). The ONLY daemon delta is the bounded per-table startup sweeps +
  warning and a stale comment correction. The frozen `TestBoundedWaitsAndBodyLimit` and route-table tests must
  pass unmodified.
- **vs the genesis lenience (M2.B, `parseGenesisRef`).** `ObservedHead` stays the single
  zero-legal COMPARED field; `TestGenesisRefLenienceIsExactlyOneField` green unmodified is an
  acceptance criterion. Decision 1 validates only what is WRITTEN — the daemon's carefully-argued
  boundary reasoning survives intact, minus its "carry-forward" clause.
- **vs the replay engine (M1.M5).** `host/replay/**` byte-unchanged. Replay reads
  `log_entries`/`objects`/`worlds` — the journal table is invisible to it; validate-on-write
  strictly STRENGTHENS its input (a store that replays today still replays, since every
  replayable store already contains only parseable refs — an unreadable entry was never
  replayable). The replay-doubling proof is untouched.
- **vs `w-effect-broker-m3` (PARKED).** This item IS the store-hardening destination that doc's
  Decision 3 and Deferred Scope point at, and its round-2 objection 2A's fix shape (intent /
  journal head / outcome / `IndeterminateEffectError` / per-handler rules / crash injection) is
  implemented here at the substrate layer in the reviewer's own terms. On unpark, M3 consumes
  `AppendIntent`/`AppendOutcome` around its dispatch (under (a) it implements the handler table
  rows; under (b) it cites this item's landed claim). M3's P2 "zero schema change" premise
  remains true OF M3 — the schema change is this item's, ratified separately.
- **vs `w-mcp-projection` (BLOCKED, prerequisite 3).** Decision 4 is the verified commit-boundary
  contract its premise row demands; the row flips to VERIFIED citing the SD.B/SD.C API + tests.
  Its other two prerequisites (upstream `#498` seam; clause-3 registry) are untouched here.
- **vs the store's object immutability + `PutObject` content verification.** Intent/outcome
  payloads are ordinary content-addressed objects; nothing is ever rewritten; the journal table
  references objects by ref exactly as `Transition.evidence` does. No "update" API is added
  anywhere.
- **vs `epoch_registry_heads` tenancy.** Untouched — ARM J2 was rejected partly to keep
  per-invocation names OUT of the named-registry table (its landed tenants stay
  `world/epoch-registry/v1` + M3's planned `world/approvals/v1`).
- **vs `verify_ail.sh`'s exact-totals gate.** Verified from the script (V14): required
  identities/tests are `world/`-keyed, totals are `world/`-scoped, module count dynamic. The new
  sketch moved the sweep 9 → 10 modules with 4/14 unperturbed — observed live (V15), and
  iteration logs asserting "EXACTLY 9 modules" must read 10 from this commit forward (the
  broker doc's note, now actual).
- **vs `schema.sql`'s history.** Byte-unchanged through M2 was a per-item choice, not a standing
  freeze; this is the FIRST post-M1 schema change and is treated with full weight: additive-only,
  migration mechanism verified (V11), existing-DDL byte-identity asserted by test (AC13), and the
  change is inside the ratification packet.

## Systemic-Issue Audit

Is this a one-off patch? It is the opposite move, twice over. (1) The charter row reported ONE
field (`PrevEntryHash`); the pre-design audit probed the pattern and found the CLASS (every ref
column, both tables, V3/V4) — the fix is one helper on every write path, not a `PrevEntryHash`
special case, so no future "zero `Interpreter`" bug report exists to be filed. (2) Three
independently-reported symptoms (CF-B-2, the broker crash window, the commit-boundary contract)
are answered by ONE substrate (journal + receipt law) instead of three bespoke mechanisms — the
anti-pattern this avoids is a broker-private journal, a projection-private receipt store, and a
store-private validity check drifting apart. The OS-gravity-well check: the kernel grows by one
additive table, one validation helper, four journal methods, and two bounded per-table scans — every
tempting consumer behavior (retry policy, handler reconciliation, REST projection, repair
tooling) has a named non-kernel destination.

## Deferred Scope

| Item | Belongs to | Boundary |
|------|-----------|----------|
| Broker dispatch wrapped in intent/outcome; the four handler reconciliation implementations | `w-effect-broker-m3` (under (a): its M3 milestones; under (b): its first post-unpark consumer change) | substrate + frozen contract table here; consumers there |
| REST/CLI projection of receipts + `IndeterminateEffectError` over the wire | `w-mcp-projection` / a daemon follow-up doc (extends the frozen route table via its own quorum) | embedded API only in this item |
| Canonical `LogEntry` encoding + `EntryHash` derivation/verification | its own kernel design (log-format-adjacent; touches D1/epoch decisions) | out — no such encoding exists today (V12), and inventing one is not crash-durability |
| Log repair / quarantine tooling | ruled out on immutability grounds (Decision 2), not deferred | detection only |
| Journal compaction / retention | future design when size data exists (OD4) | append-only, measured |
| Schema version pragma / migration framework | future design at the SECOND schema change | this change needs only the landed `IF NOT EXISTS` mechanism (V11) |
| WAL mode / fsync tuning | perf follow-up if the D7 baseline rows demand it | baseline rows land first |

## Acceptance Criteria

- [ ] **AC1 — the committed repro fixture exists and is load-bearing**:
  `durability_repro_test.go` reconstructs the V2 probe shape (zero `PrevEntryHash`) and the V3
  class shapes (zero `TransitionFn`; zero world `LogHead`) and asserts each now fails with
  structured `InvalidRefError` naming the field, with the store byte-untouched (head + row
  counts unchanged).
- [ ] **AC2 — whole-class coverage**: a table-driven test zeroes each ref field of each write
  path (`Commit`'s seven, `PutObject`'s two, `PutWorld`'s three, `SetRegistryHead`,
  `SelectHead`, `PutVerifyResult`) individually and asserts per-field rejection.
- [ ] **AC3 — the genesis exception survives exactly**: zero `ObservedHead` on a genesis commit
  still succeeds; `TestGenesisRefLenienceIsExactlyOneField` green unmodified.
- [ ] **AC4 — detection without creation**: `ScanUnreadableLog` and `ScanUnreadableWorlds` each
  find a raw-SQL-inserted poisoned row in their own table (the public API can no longer create
  one) and report index/ref + failing field; `ScanUnreadableWorlds` resumed from a mid-scan
  `afterRef` covers exactly the remaining refs in lexicographic order (no overlap, no skip); the
  daemon startup sweeps log the named warning on a fixture store and serve normally.
- [ ] **AC5 — the receipt law, non-vacuously**: `GetReceipt` drift-tested against the sketch's
  `receiptState` on all four boolean combinations (the `corrupt` arm asserted unrepresentable —
  `AppendOutcome` without intent is a structured error); after `AppendIntent` alone, the answer
  is `indeterminate` and NEVER `not-started` (`mayReportNotStarted` mirrored).
- [ ] **AC6 — commit-receipt atomicity**: crash injection at `mid-commit-before-outcome` yields
  a store with NEITHER the new world NOR the outcome; a completed commit yields BOTH; no
  interleaving exists in which exactly one is durable.
- [ ] **AC7 — idempotency semantics**: `AppendIntent` same-ID + same-bytes is a no-op returning
  the original seq; same-ID + different-bytes is `DuplicateInvocationError`; at most one outcome
  per ID enforced by the schema (raw-SQL duplicate attempt fails the UNIQUE constraint).
- [ ] **AC8 — gapless append**: journal `seq` values are exactly 1..N with no gaps after any
  green test run including rollback-inducing failures (Law 4 mirrored).
- [ ] **AC9 — the sketch is gate-real**: `sketches/storejournal.ail` — 7 contracts `verified` in
  `verify.results[]` with z3 present (never the silent-skip exit-0), **25** named tests
  (`len(tests[])`, NOT `passed_tests` — which is 32 because it also counts the 7 contract-derived
  properties; gate on the named list, correction D-B) via the
  explicit test run; full `verify_ail.sh` sweep green at 10 modules, 4/14 `world/` totals
  unperturbed.
- [ ] **AC10 — bounded everything, with the KERNEL owning the ceiling** (extended in quorum
  round 2, `gpt5-6-sol`, applied verbatim): every paged API (`PendingIntents`,
  `ScanUnreadableLog`, `ScanUnreadableWorlds`) is exercised with **zero, negative, `Max…`, and
  `Max…+1`** limits — the three interior-invalid cases return `InvalidLimitError` and the
  `Max…` case succeeds; each paginates via its keyset cursor (`fromIndex` / `afterRef`) against
  over-populated fixtures; a **multi-page startup test** proves the daemon either traverses the
  whole store or emits the distinct `integrity_scan_incomplete` warning carrying its
  continuation cursor and scanned counts — never a clean-looking message over a truncated scan;
  every crash-test wait polls a captured pid under a deadline (code-reviewed + the deadline test
  pattern).
- [ ] **AC11 — never auto-re-execute**: the probe-consumer recovery proof shows an indeterminate
  intent surfacing `IndeterminateEffectError` with the probe's counting handler at ZERO
  re-dispatches; the commit-path reconciler resolves deterministically; the Model.Infer-shaped
  probe is never retried even after reconciliation is offered (`retryAllowed` gate mirrored).
- [ ] **AC12 — frozen surfaces byte-stable by diff, not claim**: `host/replay/**`,
  `host/{hashref,canon,archive,registry}/**`, `world/**`, `cmd/**`, `scripts/verify_{ail,go}.sh`,
  `go.mod`, `go.sum`, `.github/**`; daemon route table, error envelope, and D7 constants
  asserted unchanged by the landed tests.
- [ ] **AC13 — migration proven**: opening a pre-journal fixture store with the new binary
  creates the journal table and leaves every pre-existing table's `sqlite_master` DDL
  byte-identical; `OpenReadOnly` on an un-migrated store fails structurally on journal reads
  (documented), and normally after one writer open.
- [ ] **AC14 — the cost is priced**: both journal benchmarks in the hardcoded smoke manifest;
  `bench/BASELINE.md` re-measured in one invocation with the new rows and the
  commit-vs-commit-with-receipt pair visible.
- [ ] **AC15 — intent binding is field-level, not ID-level** (round-1 objection 1): after
  `AppendIntent(id, A)`, a `Commit` carrying `id` but differing in ANY commit-defining field
  (planned `NextWorld.Ref`, `Entry.EntryHash`, `ObservedHead`, **`PrevEntryHash`,
  `TransitionFn`, `TransitionRef`, `Interpreter`** — the round-2 list) returns structured
  `InvocationMismatchError` with the store byte-untouched (head + row counts + journal
  unchanged); `Commit(id, A)` succeeds; a repeat `Commit(id, A)` after resolution returns the
  existing receipt idempotently; a repeat with different fields is the mismatch error, never a
  second receipt. The Go drift test mirrors the sketch's `intentBindsCommit` on every `tests[]`
  boundary case — field-level binding proven, not invocation-ID equality alone. **The
  `EntryHash`-preserving case is a REQUIRED row** (round 2, `gemini-3-1-pro`): mutate only
  `PrevEntryHash`/`TransitionFn`/`Interpreter` while holding `EntryHash` byte-identical — since
  V12 records that the kernel never derives `EntryHash` from entry contents, this is the exact
  drift a three-field binding would wave through.
- [ ] Both CI jobs green on every milestone PR and every dev merge.
- [ ] The kernel arms (V1/V2, J1/J2, `Commit.InvocationID`) are ratified by Mark, attended,
  BEFORE SD.A/SD.B execute — recorded in the charter STATUS like the single-writer arm A.

## Non-Vacuity — the named RED mutation for every gate (S6)

Per the iter-20/22 guardrail: a mutation's result is only evidence if the mutation itself is
valid — each must COMPILE and change observable behaviour before its result is scored; refuted
mutations go to the Ruled-out ledger, not the bin.

| Gate | Named mutation that must turn it RED |
|------|--------------------------------------|
| AC1 repro fixture | **MUT-ZERO-PREV-REACCEPTED** — remove `validateRef` from `Commit`'s `PrevEntryHash` → the fixture's rejection assertion reds AND its follow-up read-back probe reproduces the original hole |
| AC2 class coverage | **MUT-FIELD-SKIP** — exempt exactly ONE field (`Interpreter`) from validation → only the per-field table case for that field reds, proving coverage is per-field, not incidental |
| AC3 genesis exception | **MUT-GENESIS-WIDEN** — route `ObservedHead` through strict `validateRef` too → the landed genesis tests red (the exception is load-bearing in BOTH directions) |
| AC4 detection | **MUT-SWEEP-SILENT** — make a sweep swallow parse failures and return empty (applied to `ScanUnreadableLog` and `ScanUnreadableWorlds` in turn — two runs, per-table coverage) → that sweep's scan test + the daemon warning test red |
| AC5 receipt law | **MUT-RECEIPT-LIE** — return `not-started` when an intent exists without an outcome → drift test against the sketch's `receiptState`/`mayReportNotStarted` reds |
| AC6 atomicity | **MUT-SPLIT-TX** — write the outcome row in a separate transaction AFTER `Commit`'s tx commits → the `mid-commit` crash injection observes a committed world with no receipt → red |
| AC7 idempotency | **MUT-DUP-INTENT** — accept same-ID different-bytes as a silent overwrite → red at the duplicate-rejection test AND the golden-bytes/content-address check |
| AC8 gapless seq | **MUT-SEQ-GAP** — assign `seq` from a process-local counter that survives a rolled-back tx → gapless assertion reds after an induced rollback |
| AC10 bounds | **MUT-UNBOUNDED-PENDING** — ignore `limit` in `PendingIntents` → bounded test against the over-populated journal reds |
| AC10 kernel ceiling (round 2) | **MUT-CALLER-OWNS-LIMIT** — drop the `1 <= limit <= Max…` guard so an oversized caller limit is honoured → the `Max…+1` case reds (it returns rows instead of `InvalidLimitError`), proving the ceiling is kernel-enforced rather than caller-advisory |
| AC10 startup completeness (round 2) | **MUT-SCAN-SILENT-TRUNCATE** — on budget exhaustion emit the ordinary completion message instead of `integrity_scan_incomplete` → the multi-page startup test reds; the mutation is behaviour-changing by construction (a different warning string on a fixture that provably exceeds one page), satisfying the iter-20/22 valid-mutation guardrail |
| AC11 never-retry | **MUT-AUTO-RETRY** — on recovery, re-dispatch indeterminate intents automatically → counting-probe assertion reds (dispatch count > 0) and the Model.Infer-shaped case reds independently |
| AC13 migration | **MUT-DDL-DRIFT** — "improve" an existing table's DDL in `schema.sql` (add a CHECK to `log_entries`) → the `sqlite_master` byte-identity comparison against the pre-journal fixture reds |
| AC14 bench manifest | **MUT-BENCH-DROP** — drop either new benchmark name → `bench_worldd.sh`'s hardcoded-manifest gate reds (the landed V27/B1 closure) |
| AC9 sketch gates | **MUT-LAW-WEAKEN** — flip `retryAllowed`'s body to `true` → its `ensures` produces a Z3 counterexample (gate parses `verify.counterexample > 0` as failure, never the exit code); flip a `receiptState` arm → named test reds in the explicit run |
| AC15 intent binding (round-2 field list) | **MUT-INTENT-NARROW-BIND** — restore the round-1 three-field compare (`NextWorld.Ref`/`EntryHash`/`ObservedHead`) and drop the four entry refs → the `EntryHash`-preserving AC15 row reds while every other AC15 row stays green, proving the four added fields are load-bearing individually and not decorative |
| AC15 intent binding | **MUT-INTENT-UNBOUND** — weaken `Commit`'s in-tx comparison to invocation-ID equality alone (drop the field compare) → AC15's mismatch case reds: `Commit(id, B)` after `AppendIntent(id, A)` succeeds and writes a receipt claiming A resolved as B; the drift test against `intentBindsCommit`'s mismatch rows reds independently |

## Quorum verification log

### Round 1 — 2026-07-28 (iter-25): `gpt5-6-sol` REJECT + `gemini-3-1-pro` REJECT, controller PASS → BLOCKED

Both reviewers present; neither disputed the design direction — both objections were
completeness/determinism defects with reviewer-authored fixes, accepted without re-litigation.
Before the round's fixes were applied, the controller re-measured the doc's claims first-party
and made two corrections that this revision preserves: (1) the "26/26 named tests" claim was
factually wrong — `len(tests[])` was 20; `passed_tests` read 26 only because it also counts
contract-derived properties (correction D-B) — corrected throughout; (2) premise row **V23**
(the controller's independent zero-ref matrix over every `Commit` ref field) was added, and this revision now
cites it in Decision 1's arm comparison (rejection ground (e) of ARM V2).

**Objection 1 — `gpt5-6-sol` (verbatim):**

> **Strongest objection**: The public API does not bind `Commit.InvocationID` to the durable
> intent's planned commit fields. A caller can `AppendIntent(id, A)` and then call `Commit` with
> the same ID but different entry/world/head data B; the design says `Commit` writes a resolved
> outcome but specifies no in-transaction comparison with A. The resulting receipt can therefore
> claim resolution for an operation different from the recorded intent, defeating the journal's
> truthfulness and deterministic reconciliation contract.
>
> **Catch**: Decision 4, AC6, and the `outcomeMatchesIntent` sketch law establish atomic presence
> of commit and outcome, but not semantic equality among the intent, actual `Commit` arguments,
> and outcome. Behavior for reusing an already-resolved invocation ID is also unspecified.
>
> **Proposed fix**: Revise Decision 4 to require that `Commit` load the intent inside its
> transaction and compare its canonical planned `NextWorld.Ref`, `Entry.EntryHash`,
> `ObservedHead`, and other commit-defining fields against the actual request before any
> mutation. A mismatch must return structured `InvocationMismatchError` and leave the store
> untouched. If the invocation is already resolved, return the existing receipt idempotently
> only when those fields match; otherwise return the mismatch error. Add an acceptance criterion
> and RED mutation covering `AppendIntent(id, A)` followed by `Commit(id, B)`, and strengthen
> `outcomeMatchesIntent` plus the Go drift test to prove field-level binding rather than
> invocation-ID equality alone.

**Applied** (the replacement is explicit, not silent — Decision 4's round-0 step 2 established
atomic *presence* only and is superseded in place):

- **Decision 4 step 2** rewritten in the reviewer's terms: in-tx intent load + canonical
  field-compare (planned `NextWorld.Ref`, `Entry.EntryHash`, `ObservedHead`, invocation ID —
  **widened again in round 2**, see below) before any mutation; mismatch → structured `InvocationMismatchError{ID, Field}`, store
  untouched; already-resolved ID → idempotent existing receipt only on field match, else the
  mismatch error; missing intent for a set ID → structured error. Design Freeze gained the
  matching box.
- **AC15** (new) — the `AppendIntent(id, A)` + `Commit(id, B)` criterion, incl. the
  resolved-ID reuse cases and the drift test over `intentBindsCommit`'s `tests[]` boundary rows.
- **MUT-INTENT-UNBOUND** (new RED mutation) — weaken the compare to ID equality alone → AC15 reds.
- **Sketch strengthened, not just prose**: LAW 6 `intentBindsCommit` added — a total-theorem
  `ensures` predicate over the commit-defining fields (8 string parameters, exact canonical-text
  equality), Z3-`verified` live on the pinned binary; LAW 5 `outcomeMatchesIntent` retained as
  the identity half with its comment marking it insufficient alone. Re-measured counts (all
  first-party, this session): `verified` 6 → **7**, `len(tests[])` 20 → **25**, `passed_tests`
  26 → **32** — propagated to the header, P8, Files table, AC9, V15, Appendix A.
- **SD.B** files updated: `store.go` +~45 (intent load/compare), `journal_test.go` ~320
  (binding drift test); AC15 added to SD.B's acceptance checks; A11 lists the new error type.

**Objection 2 — `gemini-3-1-pro` (verbatim):**

> **Strongest objection**: Decision 2 defines a single API `ScanUnreadable(fromIndex, limit)` to
> paginate a detection sweep over BOTH `log_entries` and `worlds`. Since `log_entries` is ordered
> by an integer `index` but `worlds` is keyed by a content-addressed `TEXT` hash, a single
> integer `fromIndex` cannot safely paginate `worlds` without relying on unstable `OFFSET`s or
> implicit SQLite `rowid`s (which can mutate during a `VACUUM`). This violates the deterministic
> behavior axiom and breaks stable pagination.
>
> **Proposed fix**: Modify the API to separate the sweeps:
> `ScanUnreadableLog(fromIndex int64, limit int)` and
> `ScanUnreadableWorlds(afterRef string, limit int)`. Update the startup sweep in SD.A to
> iterate both independently.

**Applied verbatim** (the single `ScanUnreadable` API is replaced — explicit substitution,
stated at each site):

- **Decision 2**'s detection bullet now defines the pair, with `ScanUnreadableWorlds`'s cursor
  ordering stated explicitly: **lexicographic on the ref TEXT** (`WHERE ref > afterRef ORDER BY
  ref LIMIT n`) — stable, deterministic, no `OFFSET`, no rowid, resumable across `VACUUM`.
  Design Freeze gained the matching box.
- **SD.A** files (`scan.go` ~110, `scan_test.go` ~140 with per-table poison + cursor-resume
  cases) and the daemon startup sweep now iterate both sweeps independently
  (`daemon.go` +~20).
- **AC4** (per-table detection + lexicographic resume proof), **AC10** (each sweep paginates
  via its own keyset cursor), **MUT-SWEEP-SILENT** (applied to each sweep in turn), the Files
  table, P7, the Conflict Surface daemon bullet, and the Systemic-Issue Audit ("two bounded
  per-table scans") all updated.

**Post-revision gate**: sketch re-verified on the pinned v0.30.0 binary (`ai-check` 7/7
`verified`, 0 counterexamples; `ailang test --format json` 25 named / 32 total, 0 failed) and
the full `./scripts/verify_ail.sh` sweep re-run green — transcripts in Appendix A and V15.

## Axiom Compliance

| Axiom | Score | Justification |
|-------|-------|---------------|
| A1: Determinism | +2 | No wall-clock in any canonical payload; deterministic codec with golden bytes; validation makes the readable-store invariant total |
| A2: Replayability | +2 | Closes the "executed but record lost" hole that made effectful replay unable to trust its own inputs; every commit/effect attempt becomes durably detectable |
| A3: Effect Legibility | +1 | Intent/outcome objects are typed, content-addressed statements of what was attempted and what happened — including across death |
| A4: Explicit Authority | +1 | No authority change; the never-retry law removes an implicit authority (silent re-execution) no one granted |
| A5: Bounded Verification | +1 | Receipt lookup is one indexed query; pending-scan and sweep are limit-bounded |
| A6: Safe Concurrency | +1 | Builds on (and documents dependence on) the enforced single-writer model; no new concurrent surface |
| A7: Machines First | +2 | The durability law is a Z3-proven, compiler-checked artifact (6 contracts) with a pinned Go mirror, not prose |
| A8: Minimal Syntax | 0 | No language surface |
| A9: Cost Visibility | +2 | The journal tax enters the hardcoded bench manifest + baseline the milestone it exists; the paid-effect double-execution risk is priced at exactly never |
| A10: Composability | +2 | One substrate consumed by three announced consumers (commit path, broker, projection) without rework; consumers stay out of the kernel |
| A11: Structured Failure | +2 | `InvalidRefError`, `DuplicateInvocationError`, `InvocationMismatchError`, `IndeterminateEffectError` (consumer contract), three-state receipts — no boolean mush, no lie states |
| A12: System Boundary | +1 | Law in the sketch, mechanism in `host/store`, policy (retry/reconcile) explicitly pushed OUT to consumers |

**Net Score: +17** ✅ — hard axioms A1/A3/A4/A7 all ≥ 0.

## Premise Verification Log (this session, 2026-07-28 — first-party unless marked CLAIM)

| # | Claim | Evidence (exact command / file:line) | Status |
|---|-------|--------------------------------------|--------|
| V1 | Pinned binary is clean released v0.30.0 | `/tmp/ailang-v0300/ailang --version` → `AILANG v0.30.0`, commit `e37b370`, built 2026-07-19, no `-dirty` | **VERIFIED** |
| V2 | CF-B-2 reproduces exactly as the charter + iter-25 controller measured | Throwaway `host/store` probe test, this session (deleted after; the committed fixture is SD.A's): zero `HashRef.String()` = `""`/`IsZero=true` · `Commit(zero PrevEntryHash) err=<nil>` · `GetLogEntry(1) ok=false err="store: log entry 1 prevEntryHash: hashref: empty hashref text"` · raw SQL `count=1, prev_entry_hash_ref=""` · `SelectedHead()` = the NEW world ref · `GetWorld(head)` ok with `LogHead` addressing the unreadable entry · follow-on legal commit `err=<nil>`, `GetLogEntry(2) ok=true` — the permanent mid-chain hole | **VERIFIED** |
| V3 | The class is wider than the charter row: zero `TransitionFn` and zero world `LogHead` poison identically | same probe: zero-`TransitionFn` commit `err=<nil>`, `GetLogEntry(3)` → `"transitionFn: hashref: empty hashref text"`; zero-`LogHead` world commit `err=<nil>`, `GetWorld` → `"log head: hashref: empty hashref text"` | **VERIFIED** |
| V4 | `Commit` writes all refs via `String()` with zero validation; every write path shares the pattern | read `host/store/store.go:575-657` (insert at `:631-640` renders `h.PrevEntryHash.String()` et al. verbatim); `PutObject`/`PutWorld`/`SetRegistryHead`/`SelectHead`/`PutVerifyResult` read likewise (`:282`, `:333`, `:434`, `:550`, `:470`) | **VERIFIED** |
| V5 | Every reader parses strictly | `store.go:375-429` (`GetLogEntry` parses all five ref columns, structured error per field); `GetWorld` `:346-371`; `SelectedHead` `:523` | **VERIFIED** |
| V6 | The zero `HashRef` is documented invalid and renders as `""` | `host/hashref/hashref.go:59-82` ("The zero value is invalid"; `String()` returns `""` for it; `IsZero` at `:75`) | **VERIFIED** |
| V7 | `TEXT NOT NULL` is satisfied by `''` on the prev column | `host/store/schema.sql:42` + V2's raw-SQL read-back of the persisted `""` | **VERIFIED** |
| V8 | The daemon refuses zero refs at the boundary; the genesis lenience is exactly one field | `host/daemon/handlers.go:466` (`parseRef` on `prevEntryHash`), `:140-172` (`parseGenesisRef` + the written-out reasoning, incl. "That store asymmetry is a real M1 defect, recorded as a carry-forward"); `TestGenesisRefLenienceIsExactlyOneField` named there | **VERIFIED** (code read) |
| V9 | One poisoned index fails the WHOLE range read | `handlers.go` `handleLogRange`: the loop calls `GetLogEntry`, and any `err != nil` returns 500 for the entire request (read this session — the directive asked for independent confirmation, given) ; single-entry `GET /v1/log/{index}` likewise 500s (`handleLogEntry`) | **VERIFIED** (code read) |
| V10 | The genesis convention seeds entry 0's prev from a REAL hash | `host/store/store_test.go:103` (`PrevEntryHash: genesis.LogHead`) — zero prev is never legal in a written entry | **VERIFIED** |
| V11 | Additive schema migration mechanism already exists | `store.go:222-241`: `openSQLite(..., applySchema)` runs `db.Exec(schemaSQL)` on writer open; schema is all `CREATE TABLE IF NOT EXISTS`; read-only handles skip application | **VERIFIED** |
| V12 | No canonical `LogEntry` encoding / `EntryHash` derivation exists (negative existence — repair premise) | `host/canon/` contains `source.go` only (source-text canonicalization); repo grep for `EntryHash` consumers outside store round-trips: none | **VERIFIED** |
| V13 | No journal/intent/receipt/invocation-ID machinery exists (negative existence) | repo-wide grep `journal\|intent\|invocation\|receipt\|idempoten` over `host/ cmd/ world/` non-test Go: only unrelated "idempotent" prose comments in archive/canon/daemon/registry | **VERIFIED** |
| V14 | `verify_ail.sh` totals are `world/`-scoped and the module count is dynamic; Leg 2 does NOT run sketch tests | read `scripts/verify_ail.sh`: `REQUIRED_VERIFIED` keys 4 `world/` identities; `EXACT_TOTAL_VERIFIED=4` summed over `world/*` only (`case "$mod" in world/*)`); `EXACT_TOTAL_TESTS=14` from `test --format json world/`; only `checked -eq 0` is fatal on count | **VERIFIED** |
| V15 | The sketch verifies and perturbs nothing — **re-measured live after the round-1 revision added LAW 6 `intentBindsCommit`** | `cd design_docs && ai-check -timeout 5s sketches/storejournal.ail` → `check.passed: true`, `verify: {verified: 7, counterexample: 0, skipped: 0, errors: 0}`, all seven named in `results[]`; `ailang test --format json` → `failed: 0, passed: 32, skipped: 0, total: 32` — **but `len(tests[])` = 25 named tests; `passed_tests` also counts the 7 contract-derived properties (correction D-B; pre-revision figures were 6 contracts / 20 named / 26 passed, controller-remeasured iter-25)**; full `./scripts/verify_ail.sh` → PASS, "4/4 required world/ identities verified across **10** module(s)", 14 tests, re-run after the revision | **VERIFIED** |
| V16 | The store is the enforced single writer (Law-4 safety basis) | `w-worldd-m2` A1 landed + its cross-process proof; `store.go` `Open` writer-lock path present | **VERIFIED** (code present; cross-process proof inherited from the landed suite, not re-run) |
| V17 | The iter-25 controller's measured CF-B-2 output | directive text | **CLAIM — but independently corroborated by V2/V3** (my probe agrees on every line) |
| V18 | The broker round-2 objection 2A text + fix shape (half ii's source) | read `w-effect-broker-m3.md` Quorum verification log (objection + proposed fix verbatim in-doc) | **VERIFIED** (doc read; the quorum event itself is inherited record) |
| V19 | The commit-boundary contract + UNVERIFIED premise row (CF-D-2's source) | read `w-mcp-projection.md` Decision 6, AC13, premise row `Commit-boundary contract`, mutations `MUT-COMMIT-BOUNDARY-LIE`/`MUT-DROP-DEADLINE` | **VERIFIED** (doc read) |
| V20 | No committed `.db` fixture exists (poisoned-population premise) | `find . -name "*.db" -not -path "./.git/*"` → empty | **VERIFIED** |
| V21 | Go build/test of the new code | no Go exists at doc time; LOC/structure are estimates | **CLAIM** — the sprint's gates are the check |
| V22 | `verify_go.sh` green in this worktree at doc time | **UNVERIFIED in this sandbox** — not run this session (no Go changes exist to gate; the last first-party green is the iter-24 controller record) | **CLAIM** |
| V23 | **The full EIGHT-field matrix, measured** — CONTROLLER, iter-25, independently of the designer: a table-driven probe zeroing each ref field of `Commit` in isolation, then reading back. **Corrected in round 2**: this row originally said "seven-field matrix … the seventh, `NextWorld.Ref`" while listing seven POISONED fields before it — an off-by-one in the controller's own row, caught by `gemini-3-1-pro` and fixed here rather than quietly. There are **eight** ref fields; seven poison, one does not | `commit_err=<nil>` for **all eight**. Poison confirmed per field (7): `TransitionFn` → `GetLogEntry err="…transitionFn: hashref: empty hashref text"` · `Interpreter` → `"…interpreter: …"` · `EntryHash` → `"…hash: …"` · `TransitionRef` → `"…transitionRef: …"` · `PrevEntryHash` → `"…prevEntryHash: …"` · `NextWorld.LogHead` → `GetWorld(head) err="…log head: …"` · `NextWorld.StateRoot` → `GetWorld(head) err="…state root: …"`. The **eighth**, `NextWorld.Ref`, commits and reads back **"fine"** — an empty-string world ref becomes the selected head: degenerate-but-readable, hence the one field a read-side fix (ARM V2) could never catch, and the sharpest single argument for ARM V1 | **VERIFIED** |

**Upstream findings: none.** The sketch verified clean on first authoring by following the
already-documented v0.30.0 disciplines (total-theorem ensures per U3; no contracted-callee calls
per U1; string-projection outputs per U2) — recorded here because it is the first sketch in four
for which the patterns held without surfacing a new toolchain defect. The round-1 addition of
LAW 6 (`intentBindsCommit`, an 8-parameter string-equality contract) also verified on the first
run — the widest-arity contract proven in this repo to date, same disciplines.

## Open Decisions (escalated with recommended defaults — the sprint proceeds on the defaults)

- **OD1 — invocation-ID form.** Opaque caller-supplied string vs a kernel-derived content hash of
  the intent. **Default: opaque caller string** (uniqueness schema-enforced; a content-derived ID
  would make two identical legitimate attempts collide) — the intent OBJECT is content-addressed
  regardless, so integrity is not the ID's job.
- **OD2 — daemon behavior on a store with historic holes.** Warn-and-serve vs refuse-to-start.
  **Default: warn-and-serve** (Decision 2's argument: readable rows stay servable, the hole
  already fails loudly on access; refusing would turn one historic hole into a total outage).
- **OD3 — schema version pragma.** **Default: not now** — one additive table rides the landed
  `IF NOT EXISTS` mechanism; a version pragma is warranted at the second schema change, as its
  own small design.
- **OD4 — journal retention.** **Default: append-only forever, measured** — a journal row is
  ~150 bytes + payload objects; revisit with real size data, as a design, not a sprint tweak.

## Appendix A — the verified sketch (landed with this doc as `design_docs/sketches/storejournal.ail`)

Verified this session on the pinned binary — **re-run live after the round-1 revision added
LAW 6 `intentBindsCommit`**: `check.passed: true` · **7/7 contracts `verified`**
(0 counterexamples, 0 errors, 0 skips — z3 on PATH, results enumerated) · **25 named tests pass**
(0 failed, 0 skipped; `passed_tests`/`total_tests` report 32 because they also count the 7
contract-derived properties — gate on `len(tests[])`, correction D-B). 163 lines. The file in `sketches/` is the artifact; it is not duplicated here
byte-for-byte to avoid the two-copies drift the compiler-checked-docs rule exists to prevent —
the sweep checks the file, and this appendix records the verification transcript:

```json
{"check": {"passed": true, "error_count": 0},
 "verify": {"available": true, "verified": 7, "counterexample": 0, "skipped": 0, "errors": 0,
  "results": [["intentBindsCommit","verified"],["writableRefText","verified"],
              ["isIndeterminate","verified"],["mayReportNotStarted","verified"],
              ["retryAllowed","verified"],["journalSeqNext","verified"],
              ["outcomeMatchesIntent","verified"]]}}
```

```json
{"failed_tests": 0, "passed_tests": 32, "skipped_tests": 0, "success": true, "total_tests": 32}
```

The seven laws: `writableRefText` (the write-validity floor — no ref column may persist the zero
rendering), `isIndeterminate` (durable intent + no durable outcome), `mayReportNotStarted` (the
never-lie rule: a definitive "not started/committed" only when NO intent exists),
`retryAllowed` (re-execution only when not indeterminate or explicitly reconciled),
`journalSeqNext` (gapless append from 1), `outcomeMatchesIntent` (exact-match invocation
identity binding), `intentBindsCommit` (round-1 addition: field-level commit binding — a commit
resolves ONLY the intent whose commit-defining fields it matches, so `AppendIntent(id, A)` +
`Commit(id, B)` is structurally a mismatch, never a receipt); plus the tests-covered canonical
projection `receiptState` (`not-started` /
`indeterminate` / `resolved` / `corrupt`), whose arms restate the proven laws — the
`decideLabel` string-projection pattern, applied because `tests[]` cannot express ADT expected
values on v0.30.0 (U2) and canonical text forms take named tests per S1.

### Round 2 — 2026-07-28 (iter-25): `gpt5-6-sol` REJECT + `gemini-3-1-pro` REJECT, controller PASS → BLOCKED → **NARROW-REFINEMENT CARVE-OUT APPLIED**

Both reviewers present (`$0.157032`). Both round-1 objections were confirmed applied; the two new
objections are again **completeness/determinism defects with concrete, reviewer-authored
`proposed_fix` prose, and neither disputes the design DIRECTION** — path, arms, milestone shape,
ratification framing and the (a)/(b) gating section all stand. Both limbs of the mission-control
narrow-refinement carve-out therefore hold, and the controller applied the reviewers' **verbatim**
fixes as a bounded 2nd revision. **No third round** — the carve-out is one bounded revision, not a
re-litigation. This SATISFIES the objections; it is not force-passing (Standing rule 2 still
forbids proceeding over a contested DIRECTION, and none was contested).

**Objection 3 — `gpt5-6-sol` (verbatim):**

> **Strongest objection**: The claimed bounded-wait/allocation contract is not actually
> specified: `PendingIntents(limit int)`, `ScanUnreadableLog(..., limit int)`, and
> `ScanUnreadableWorlds(..., limit int)` accept arbitrary caller-controlled limits with no
> kernel-enforced maximum or defined behavior for zero, negative, or oversized values. Merely
> stopping at the supplied limit does not bound allocation or query work, so P7, AC10, and the A5
> compliance claim are unsupported and violate the bounded-waits axiom.
>
> **Catch**: The startup integrity check is also described as only "one bounded sweep" per table,
> which can silently miss poisoned rows after the first page. No fixed page size, total scan
> budget, continuation behavior, or explicit truncation warning is specified.
>
> **Proposed fix**: Add kernel constants such as `MaxPendingIntentsPage` and
> `MaxIntegrityScanPage`; require `1 <= limit <= Max...` and return a structured
> `InvalidLimitError` otherwise. Specify that startup scans page with a fixed page size until
> either completion or a fixed total-row/time budget is reached. If the budget is exhausted, emit
> a distinct `integrity_scan_incomplete` warning containing the continuation cursor and counts,
> never implying the store was fully scanned. Extend AC10 with zero, negative, maximum, and
> maximum-plus-one cases, an overlarge caller-limit RED mutation, and a multi-page startup test
> proving either complete traversal or the explicit incomplete-scan warning.

**Applied verbatim.** The superseded claim is stated rather than deleted: round 1 called the APIs
bounded because they "respect the caller's limit", which the reviewer refuted — a caller-supplied
limit is not a bound. Now: Decision 2 gains the kernel-constants / `InvalidLimitError` /
fixed-page-size / total-budget / `integrity_scan_incomplete` block; `PendingIntents`'s API-table
row states the range contract; **P7** is rewritten around "the KERNEL owns the bound"; the Design
Freeze gains a matching box; **AC10** is extended with the zero / negative / `Max…` / `Max…+1`
cases and the multi-page startup test; two new RED mutations — **MUT-CALLER-OWNS-LIMIT** and
**MUT-SCAN-SILENT-TRUNCATE** — join the table; the Files table re-prices `scan.go` and
`daemon.go`.

**Objection 4 — `gemini-3-1-pro` (verbatim):**

> **Strongest objection**: Decision 4 and AC15's intent binding only compares NextWorld.Ref,
> Entry.EntryHash, and ObservedHead against the actual request. Because Premise V12 explicitly
> states the kernel does not verify EntryHash against the entry's contents, a caller can submit an
> intent but call Commit with mutated fields (e.g., PrevEntryHash, TransitionFn, Interpreter)
> without changing the EntryHash. Commit will accept this silently mismatching payload, and the
> resulting receipt will falsely claim the exact original intent succeeded, fundamentally breaking
> the semantic equality contract established in Objection 1.
>
> **Catch**: Premise V23 asserts a 'seven-field matrix' and calls NextWorld.Ref 'the seventh', but
> explicitly lists 7 poisoned fields before it (TransitionFn, Interpreter, EntryHash,
> TransitionRef, PrevEntryHash, NextWorld.LogHead, NextWorld.StateRoot), making NextWorld.Ref the
> eighth field.
>
> **Proposed fix**: Update Decision 4 (steps 1 and 2) and AC15 to expand the 'commit-defining
> fields' to explicitly include all unverified LogEntry ref fields (PrevEntryHash, TransitionFn,
> TransitionRef, Interpreter), ensuring the actual committed entry cannot drift from the recorded
> intent.

**Applied verbatim** — and note what this objection is: it does not challenge objection 1's fix,
it finishes it through the gap objection 1 left open. Decision 4 steps 1 and 2 and the Design
Freeze box now name all seven commit-defining ref fields; **AC15** carries the
`EntryHash`-preserving row as REQUIRED (mutate `PrevEntryHash`/`TransitionFn`/`Interpreter` while
holding `EntryHash` byte-identical — the exact drift a three-field binding waves through, and it
is only possible *because* V12 records that the kernel never derives `EntryHash` from entry
contents); and **MUT-INTENT-NARROW-BIND** proves the four added fields are load-bearing
individually by restoring the three-field compare and requiring exactly that row to red.

**The catch was a defect in the CONTROLLER's own premise row, and it is corrected, not buried.**
V23 said "seven-field matrix … the seventh, `NextWorld.Ref`" while listing seven *poisoned* fields
before it — eight fields, not seven. The row now reads **eight ref fields: seven poison, one
(`NextWorld.Ref`) is degenerate-but-readable**, with the correction and its source stated inline.
An independent reviewer catching an arithmetic error in the controller's own first-party evidence
is the quorum working exactly as intended; the lesson is recorded in the iteration log rather than
being quietly patched.

**Controller's first-party re-measurement after this revision** (never trusting the designer's or
its own prior report): `ai-check` → `check.passed: true`, **7/7 contracts `verified`**, 0
counterexamples, 0 errors; `ailang test --format json` → `len(tests[]) = 25` named tests,
`passed_tests = 32` (25 named + 7 contract-derived properties), 0 failed;
`AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` → **PASS**, 4/4 required `world/`
identities across **10** modules, 14/14 named `world/` tests. Diff is doc + sketch only.

---
**Document created**: 2026-07-28 (mission iteration 25, designer role). **Ratification needed
before sprint**: Decision 1 arm (V1 recommended), Decision 3 arm (J1 recommended), Decision 4's
`Commit.InvocationID` extension — one attended packet, ideally alongside the parked (a)/(b)
answer for `w-effect-broker-m3` (either answer is compatible with this design; the milestone cut
is pre-drawn above).
