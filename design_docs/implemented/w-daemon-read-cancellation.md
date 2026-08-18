# w-daemon-read-cancellation — a bounded elapsed-time contract for the daemon's read path

- **Status**: **PARKED `needs-human-review`** (iteration 86) — two quorum rounds, both `blocked`,
  both with BOTH external reviewers PRESENT (`absent_reviewers: []` each round, so neither verdict
  is a pass-with-a-hole). Round 2's gemini reject was a PREMISE objection and is **DISCHARGED by
  controller measurement** — both claims measured TRUE, recorded as V28/V29, design unchanged.
  Round 2's gpt5-6-sol reject is a **DIRECTION dispute about the item's SCOPE BOUNDARY**, which
  forecloses the narrow-refinement carve-out (Gate 2: a carve-out requires every remaining objection
  to carry a concrete fix AND not dispute the direction). It needs a one-word human A/B — see §13.
- **Item**: queue item 18, `w-daemon-read-cancellation`, clause-3 (Standing rule 6: every wait is bounded)
- **Filed**: 2026-08-13, iteration 82, on first-party controller measurement
- **Designed**: 2026-08-14, iteration 86
- **Revised**: 2026-08-14, iteration 86, quorum round 2 — two reviewer rejects plus one controller
  finding; §12 records the disposition of every objection, with measurements
- **Estimated**: 1.5 days (three milestones, each independently CI-green; re-priced in round 2
  for the cancel-release gate, the ratchet test, and the two context fold-ins)
- **Measurement base**: `6fd26f0` (dev, clean; the Go tree is byte-identical to the gate-proven
  merge commit `aaada20` — V22)
- **Prerequisite of**: item 14 `w-workbench-read-only` (Mark ratified option B attended 2026-08-14:
  item 14 is DEFERRED BEHIND this item; its declared residual WB-R1 is discharged HERE, not there)
- **Files changed by implementation**: enumerated in §8; this document changes only itself

Every present-tense codebase statement is backed by a command and observed output in §11.
Controller-supplied facts (the six measurements in the queue row) were re-run first-party at this
base before use; two of them are refined below (V2, V7, V14) — the refinements strengthen the
finding, they do not weaken it.

---

## 1. Problem

Standing rule 6 says every wait is bounded. The daemon has a real, ratified, non-vacuous gate for
that rule — `TestBoundedWaitsAndBodyLimit` (host/daemon/daemon_test.go:202) pins six timeout
constants as literals and asserts all four `http.Server` waits are wired and non-zero (V10). And
the rule is still unenforced on every one of the seven `GET /v1/*` routes, because **a gate's
coverage is a property of the layer it observes**. The D7 gate observes the transport.
`http.Server.WriteTimeout` is a deadline on the connection's write side; it cannot cancel a
goroutine blocked inside a store call. The wait this item bounds happens one layer below the only
layer the green gate can see — which is precisely why three iterations of bounded-waits discipline
never noticed it.

What is measured at the base (each zero carries a same-scope known-positive control, §11):

- `host/store/store.go` contains zero `context.Context` (control: 14 `func (s *Store)` methods in
  the same file) (V1). All five read getters **on the daemon read path** are context-free — the
  queue row lists four; first-party measurement adds `GetRegistryHead:628`, which
  `GET /v1/registry/{name...}` calls, so the row's list was one getter short (V2). The store
  itself has a **sixth** getter, `GetVerifyResult:773`, deliberately outside this enumeration:
  its only production caller is `host/replay/replay.go:191` and no daemon handler reaches it
  (V25) — the count is five because the daemon read path is five, not because the enumeration
  stopped early. (It follows its caller into the follow-on item, §10.)
- Each getter blocks in `s.db.QueryRow` (5 call sites, V3) against a pool holding **exactly one
  connection** — `db.SetMaxOpenConns(1)` at store.go:297 (V8). One slow read therefore queues
  every other read AND the commit path behind it, with no deadline anywhere: the connection-pool
  wait in `database/sql` is unbounded unless the caller supplies a context.
- No handler consults the request context: `r.Context()` appears zero times in
  host/daemon/handlers.go (control: 9 `http.ResponseWriter` in the same file) and zero times in
  host/daemon/daemon.go — `handleHead` discards the request as `_` (V4).
- The production DSNs carry no `busy_timeout`; the only occurrence in host/store is a test-input
  string at writer_lock_test.go:609 (V6). Refinement (V14): SQLite's default busy timeout is 0,
  so today a cross-process lock collision fails **immediately** with SQLITE_BUSY rather than
  hanging — the defect at the lock layer is not an unbounded wait but the *absence of any declared
  policy*: a transient lock burst becomes an instant 500 instead of a bounded retry.
- The same pass carried a second, independent finding: the JSON handlers pass raw internal error
  text to the client — 11 `err.Error()` occurrences in handlers.go, six of them on the `Internal`
  (500) branches at handlers.go:220, 247, 277, 325, 351, 419 (V5). §2.6 decides
  sanitize-vs-expose here so item 14's renderer inherits a decision, not a defect.

The blast radius is all seven GET routes (route table daemon.go:461-468, V9): `/v1/health` touches
no store; the other six read through the five getters. `POST /v1/commit` shares the same
starved single-connection pool but is the write path; its elapsed-time contract is declared out of
scope in §10 (DR-1), not silently dropped.

## 2. The design question, settled

### 2.1 Question

What is the smallest change that gives all existing GET routes one bounded elapsed-time contract —
a request-scoped deadline that actually reaches the blocking call, an explicit timeout status, and
a test that reds when the propagation is removed — without growing the kernel or a new package?

### 2.2 Decision — context-aware store reads, by signature change, not by variant methods

The five read getters change to context-first signatures and route through
`s.db.QueryRowContext(ctx, …)`:

```go
func (s *Store) GetObject(ctx context.Context, ref hashref.HashRef) (Object, bool, error)
func (s *Store) GetWorld(ctx context.Context, ref hashref.HashRef) (World, bool, error)
func (s *Store) GetLogEntry(ctx context.Context, index int64) (LogEntry, bool, error)
func (s *Store) GetRegistryHead(ctx context.Context, name string) (hashref.HashRef, bool, error)
func (s *Store) SelectedHead(ctx context.Context) (hashref.HashRef, bool, error)
```

The alternative — keeping the context-free getters and adding `*Ctx` variants — is rejected on
this mission's own recurring lesson (iterations 77 and 85): a recogniser's coverage is a property
of its input grammar, and any gate that must *forbid the unbounded spelling* while both spellings
compile is a grep pretending to be a type system. What deleting the context-free spelling buys is
stated precisely, because round 1 of this document over-claimed it (§12, B1): the compiler forces
every call site to **state** its context, so the unbounded case is explicit and greppable at each
site — it does **not** make an unbounded read unrepresentable, because a deadline-free
`context.Background()` remains a valid argument and is deliberately passed at 11 production sites
after this item (below). The enforced complement is the ratchet test (§2.8): the 11 deadline-free
sites are enumerated and pinned, so the set can shrink but never silently grow. The measured cost
is 22 non-test call sites across 7 files plus 64 test call sites (V15) — mechanical,
compile-driven, priced in §9.

Non-daemon callers split three ways, by measurement (V26), not by package name:

- **`host/transitionreg` (3 sites) — folded in, real context.** `ReadSnapshot` and `Publish`
  already take a `context.Context` (transitionreg.go:70, 221) and the three getter sites
  (:74, :89, :233) sit directly inside them; they pass the caller's ctx through. Zero API change.
- **`host/broker/broker.go` (2 sites) — folded in, real context.** Both sites (:361, :396) are
  inside unexported `invokeReplay`, whose caller `Session.invoke` already holds a ctx and drops
  it at broker.go:173. `invokeReplay` gains a ctx parameter — unexported, zero API change.
- **`host/broker/approve.go` (8), `host/registry` (2), `host/replay` (1) — `context.Background()`
  verbatim, pinned.** Today's exact behavior made *visible* at the call site. Threading real
  deadlines here requires exported-API changes behind operator entry points and per-path deadline
  policy decisions (V27); that residue is the named follow-on item in §10 (DR-2), not a note.

Both mechanisms the queue row named are used, at the layer where each is true:

1. **Context is the elapsed-time bound.** It bounds the `database/sql` single-connection pool
   wait (the principal unbounded wait, V8) and interrupts an in-flight statement: modernc.org/
   sqlite v1.54.0 wires `ctx.Done()` to `sqlite3_interrupt` at statement level
   (stmt.go:105,295 — V14). No other mechanism bounds the pool wait at all.
2. **`busy_timeout` is the lock-layer policy, not the bound.** `writeDSN`/`readOnlyDSN`
   (host/store/writer_lock.go:176,187 — the queue row misplaced these in store.go, V7) gain
   `_pragma=busy_timeout(2000)`. The driver applies `_pragma` per physical connection at open,
   busy_timeout deliberately first (driver sqlite.go:217-231, V14), so the pragma survives pool
   connection recycling — which an `Exec("PRAGMA …")` after `Open` would not. Value 2000 ms:
   below the 10 s request deadline, so the context remains the outer bound; large enough to ride
   out a writer's commit burst instead of 500ing instantly.

### 2.3 The deadline: one new D7 constant, wired as a field

```go
// readDeadline bounds the elapsed time of every store read a GET handler
// performs (D7 addendum, w-daemon-read-cancellation). It must stay well below
// writeTimeout: the 503 must be writable inside the connection's write window.
readDeadline = 10 * time.Second
```

`Daemon` gains `readDeadline time.Duration`, set by `New` from the constant — the same
field-not-constant idiom `drainTimeout` already uses at daemon.go:243, and for the same reason:
the wiring is assertable and the value is shrinkable in tests, which is the only way the timeout
branch can be exercised without a ten-second test. `TestBoundedWaitsAndBodyLimit`'s constants
table (daemon_test.go:205-217) gains a seventh literal row.

Handlers install it through one helper, so removal is a one-line mutation with one named test:

```go
func (d *Daemon) readCtx(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), d.readDeadline)
}
```

All six store-reading GET handlers (`handleHead` in daemon.go — its `_ *http.Request` becomes
`r` — plus the five in handlers.go) call it, **and every handler must `defer cancel()`
immediately after calling `readCtx(r)`** — round 1 never stated this obligation (§12, A1). A
`context.WithTimeout` whose CancelFunc is never called leaks its timer and its context subtree
until the deadline fires, once per store-reading request. The mandate is a gate, not advice:
`TestReadCtxCancelledAfterHandler` (§2.5, layer 4) asserts on every route that the context the
getter received is already cancelled the moment `ServeHTTP` returns — with `readDeadline` left
at its 10 s default, only the handler's `cancel()` can produce that observation within test
time — and MU11 is the mutation that removes one `defer cancel()` and must red it. Deriving from
`r.Context()` rather than `context.Background()` is load-bearing: a client that disconnects
releases the store connection immediately instead of holding it for the residual deadline.
`handleLogRange`'s bounded loop passes the same context to every iteration, so a page that goes
slow mid-loop times out as a whole (nothing has been written yet — `writeJSON` runs once at the
end — so the 503 replaces the page cleanly).

Interplay with the frozen transport constants, stated so nobody re-derives it wrong:
`WriteTimeout` (30 s) starts after the request headers are read; a read that consumes the full
10 s deadline leaves 20 s of write window for the 503. The CLI client's `DefaultClientTimeout`
(30 s) is above the server deadline, so a slow read surfaces as the server's explicit 503, never
as a client-side context error racing it.

### 2.4 The explicit timeout status: 503, class `Timeout`, and the sketch moves with it

A read that exceeds the deadline returns **HTTP 503** with the existing `writeAPIError` envelope
and a new class:

```json
{"error": {"class": "Timeout", "message": "read deadline (10s) exceeded"}}
```

Why 503 and not 504: 504 is a gateway's statement about an upstream server; this daemon is the
origin. The standard library's own `http.TimeoutHandler` answers 503 for exactly this shape. Why
not `http.TimeoutHandler` itself: it bounds the *response*, then abandons the handler goroutine —
which stays blocked in the store, still holding the one connection. It is the D7 lesson again: it
observes the transport while the wait happens below it. Rejected.

Classification is by the context, not by error text: after a failed store call the handler checks
`ctx.Err() != nil` first and only then falls through to the 500 branch. This is deliberate — the
driver's interrupt path can surface as a SQLITE_INTERRUPT error that does not wrap
`context.DeadlineExceeded`, so string- or `errors.Is`-only classification would misfile real
timeouts as 500s. `errors.Is(err, context.DeadlineExceeded)` is checked as well, but the context
is the authority.

The `APIError` envelope mirrors `ApiError`/`httpStatus` in the frozen, compiler-checked
`design_docs/sketches/worlddapi.ail` **exactly** (handlers.go:16-18, V11). An explicit timeout
class therefore requires touching one `.ail` file: `ApiError` gains `Timeout(string)`, `httpStatus`
gains the arm `Timeout(_) => 503`, and the inline test table gains
`((Timeout("read deadline exceeded")), 503)`. This exact edit was probed on a scratch copy against
the pinned v0.30.0 binary before entering this document: pristine control `check.passed=True,
verify.errors=0, cex=0`; edited file identical verdicts, inline tests 18 → 19, all pass (V13).
**It moves no gate pin** (V12): `EXACT_TOTAL_VERIFIED=10` counts `world/*` modules only
(verify_ail.sh's `case "$mod" in world/*)`), sketches carry empty required-identity sets by
design, Leg 2 runs `ailang test world/` so `EXACT_TOTAL_TESTS=39` cannot see a sketch test row,
and the module set is unchanged so `LEG1_MODULES` stands. The honest consequence — the sketch's
new test row is *not* CI-enforced — is handled by the repo's existing mirror idiom (§5): a Go test
replays the sketch vector and pins the 503 literal, exactly as
`TestIsLoopbackHostMirrorsSketchPredicate` pins the loopback predicate.

### 2.5 The test that reds when the propagation is removed (the non-negotiable)

Four layers, each observing a write only its mechanism can produce. Every watchdog in this
section obeys one rule, added in round 2 (§12, B2): **a watchdog that reds must also RELEASE the
blocked call and confirm the blocked goroutine exited before cleanup runs** — a bare `t.Error`
from a watchdog goroutine unblocks nothing, and a deferred `s.Close()` behind a still-parked
getter converts a clean red into a suite-wide hang. Each layer names its release mechanism.

1. **Store layer, real store, no fakes** — `TestReadGettersHonorContext` (white-box, package
   `store`): seed a `:memory:` store, occupy the single pool connection via `s.db.Conn(…)`, then
   call each of the five getters — in a goroutine writing to a result channel — with an
   already-expired context. Correct code returns in microseconds with a context error (the pool
   wait consults the context). Under the mutation that drops `ctx` inside a getter, the call
   blocks in the pool wait; the 2 s watchdog turns that into a deterministic red with a named
   message. Release: the watchdog's red path is `t.Error` **then close the occupying
   `*sql.Conn`**, which returns the sole connection to the pool; the parked getter acquires it,
   completes (result discarded), and its exit is confirmed by draining the result channel under
   a second 2 s bound — so the deferred `s.Close()` runs against a quiescent pool and a red
   subtest exits within ~4 s. One subtest per getter, so each of the five propagation sites has
   its own selector.
2. **Daemon layer, real store, no fakes** — `TestDaemonReadDeadline/real-store-expired-deadline`:
   a daemon over a real seeded `:memory:` store with `d.readDeadline` shrunk to 1 ns answers 503
   class `Timeout` on every store-reading route. Remove the deadline installation (`readCtx`
   mutation) and the store answers normally — 200 — which reds the 503 assertion. Deterministic,
   fake-free, and the observable (class `Timeout`) is written by no other branch in the codebase.
   No watchdog is needed and none is claimed: with an already-expired deadline against a real
   store, the correct arm answers 503 in microseconds and every mutated arm answers 200 in
   microseconds — neither can block.
3. **Daemon layer, blocking stimulus** — the six read handlers move from the concrete
   `d.store` field to a six-value read interface `reads readStore` that `New` wires to the same
   `*store.Store` (a seam that the production path passes through unchanged). The hang-shaped
   tests wrap — not replace — the real store: `blockingStore` embeds it and **overrides all five
   getters** to block until `ctx.Done()`, then return `ctx.Err()`. All five, because the six
   routes reach the store through five different getters: round 1's single-getter override would
   have let every other route fall through to the embedded real store, answer 200 instantly, and
   spuriously red the 503 assertion (§12, A2 — the reviewer's catch, applied verbatim). Release:
   every blocking override also selects on a test-owned `escape` channel and returns a distinct
   sentinel error when it closes; the 2 s watchdog's red path is `t.Error` **then
   `close(escape)`**, which releases the parked getter under any mutation (under MU2 the fake's
   context is Background-derived and would otherwise stay parked until the 10 s deadline), so
   the handler returns, `ServeHTTP` completes, and cleanup runs — and the sentinel can never be
   read as a pass. `TestDaemonReadDeadline/blocking-store` drives all six routes table-driven
   with a 50 ms deadline and asserts 503 within the watchdog; `TestDaemonReadDisconnect` cancels
   the request mid-block with the deadline left at its 10 s default and asserts the store call
   unblocks within 2 s — which is exactly the arm that discriminates `r.Context()` from
   `context.Background()` in `readCtx` (under that mutation the unblock arrives only at the 10 s
   deadline, and the watchdog reds — and releases — first).
4. **Cancel-release layer** — `TestReadCtxCancelledAfterHandler` (the §2.3 mandate's gate): a
   `recordingStore` (same wrap-not-replace seam — embeds the real store, overrides all five
   getters to record the context they receive and delegate) serves all six routes to a normal
   200; immediately after `ServeHTTP` returns, the recorded context must already report
   `context.Canceled`. `readDeadline` stays at its 10 s default, so within test time that
   observation is producible only by the handler's `defer cancel()`; remove one (MU11) and that
   route's recorded `ctx.Err()` is nil at assert time. Nothing blocks in this layer (the
   delegate is a real, unoccupied store), so no watchdog is needed.

The iteration-80 vacuity trap (a seam that REPLACES makes mutations of the replaced body
vacuous) is priced in, and was re-checked after round 2 widened the fake to all five getters:
the fake replaces store *bodies* only — all handler code in layer 3 runs production paths;
every mutation in §6 that targets handler code is exercised by layer 2's fake-free test as well
as layer 3's; every mutation targeting store internals (MU4a–e) is witnessed only in layer 1
against the real store; MU11's witness records-and-delegates to a real store; MU12's witness
reads the source tree, not any fake. No mutation's only witness runs behind the fake.

The lock-layer policy gets a real-contention test with **pre-registered outcomes** (the MU15
lesson: record which form landed): `TestReadRetriesUnderTransientExclusiveLock` (package `store`,
file-backed DB) holds `BEGIN EXCLUSIVE` from a second raw driver connection, launches a getter
with a 5 s context, and releases the lock after 300 ms. With `busy_timeout(2000)` the retry loop
wins and the read returns the row; with the injection removed the read fails instantly with
SQLITE_BUSY — red. Whether the driver's busy sleep is ctx-interruptible is deliberately NOT
claimed here; the test's bound assertion is ≤ 3 s (covering both the interrupt-wins and the
busy_timeout-expires outcome), and the sprint must record which outcome the measurement shows.
Termination is by construction rather than by watchdog: the `BEGIN EXCLUSIVE` lock is released
both by the 300 ms timer and by an unconditional deferred release at test exit, so no failure
path leaves it held; the getter carries its own 5 s context as the outermost bound; every arm —
retry-wins, interrupt-wins, injection-removed instant SQLITE_BUSY — returns through one of
those, and cleanup runs after the getter's result channel is drained.

### 2.6 Sanitize-vs-expose, decided: sanitize the six `Internal` branches, keep the rest

**Decision: the six 500 branches stop echoing `err.Error()`; every other branch keeps its
current text.** The 500 body message becomes the fixed constant `"internal store failure"`; the
full error is written to a new `Config.ErrorLog io.Writer` (nil → `os.Stderr`), one line with the
route and the verbatim error.

Rationale, not vibes:

- **Expose is wrong for 500s.** Store errors interpolate the display DSN path
  (`store: open %q`, store.go:292), schema state, and driver strings — environment detail about
  the daemon host. Loopback is locality, not authentication (the ratified item-14 language): any
  local process can hold the socket, and item 14 will pipe these bodies into rendered HTML. The
  party entitled to internals is the process owner, who owns stderr.
- **Sanitize is wrong for 400s.** The five `BadRequest` `err.Error()` sites (handlers.go:215,
  242, 392, 399, 410) echo parse failures of *the client's own input* — hashref text, JSON decode
  positions — and contain no server state. Stripping them would gut the API's main debugging
  affordance to protect nothing. They stay.
- 409 keeps its machine-readable two-head body (that is protocol, not leakage); 503's message
  names only the deadline constant.
- **Not to `announce`**: the announce writer is a one-line protocol (`ListenAnnouncePrefix`,
  daemon.go:140) whose consumers read exactly one line, and iteration 28 measured extra announce
  lines deadlocking `Run` against an `io.Pipe`. Error logging gets its own writer (V17).

Post-change, `grep -c 'err.Error()' host/daemon/handlers.go` is exactly **5** (11 minus the six
Internal branches) — AC5's snapshot.

### 2.7 Why is this not a package? (S3)

No new package and no kernel growth: the item edits existing host packages (`host/store`,
`host/daemon`) at their existing boundary, and the one `.ail` edit is to a compiler-checked
*sketch* (S4 artifact), not to `world/`. The single new Go surface is a six-method unexported
interface inside `host/daemon`. There is nothing here to publish, version, or cascade; making
"bounded reads" a package would mean exporting the daemon's private read seam, which is authority
widening, not modularity.

### 2.8 The deadline-free residue, pinned: a ratchet now, the boundary reject at zero

gpt5-6-sol proposed rejecting a nil or deadline-free context at the store boundary. Evaluated on
its merits and measured, not dismissed: the guard is the single highest-value part of the wider
fix, and it is also the one part that CANNOT land in this item — the moment M1 lands it would
break all 11 production call sites this item deliberately leaves deadline-free (approve path 8,
registry bootstrap 2, replay 1 — V26/V27), because each site's honest bound is a policy question
this item cannot answer (what bounds an attended approval that is *designed* to wait on a
human?). The guard is therefore the follow-on item's declared closing move (§10, shape (c)),
landable exactly when the ratchet below reads zero.

What IS cheap and in-item lands here: **`TestNoNewDeadlineFreeStoreReads`** (package `store`)
scans the production (non-`_test`) `.go` source of the caller packages for getter calls whose
context argument is `context.Background()` and pins the exact per-file counts — approve.go 8,
registry.go 2, replay.go 1, everything else 0. A new deadline-free store read anywhere in
host/ or cmd/ reds it (MU12); threading a real context through an existing site shrinks a pinned
count, a deliberate one-line test edit in the same diff. This converts DR-2 from a spelled
residue into an enforced no-new-unbounded-reads contract, and it makes the follow-on item's
progress mechanically observable (11 → 0). The test lives in `host/store`, not `host/boundary`,
because the boundary package's file census is pinned at `wantFileCount = 1` (V16) and this item
adds no file there.

## 3. Proposed change, per file

| File | Change |
|---|---|
| `host/store/store.go` | five getters become context-first; `QueryRow` → `QueryRowContext` (5 sites) |
| `host/store/writer_lock.go` | `writeDSN`/`readOnlyDSN` inject `_pragma=busy_timeout(2000)` when the caller's DSN did not set one |
| `host/store/context_read_test.go` (new) | `TestReadGettersHonorContext` (5 subtests), `TestProductionDSNSetsBusyTimeout` (PRAGMA readback on both handle kinds), `TestReadRetriesUnderTransientExclusiveLock`, `TestNoNewDeadlineFreeStoreReads` (the §2.8 ratchet) |
| `host/store/*_test.go` (existing) | mechanical: getter call sites gain `context.Background()` (V15: 64 test-site total across repo) |
| `host/daemon/daemon.go` | `readDeadline` constant + field; `Config.ErrorLog`; `reads readStore` seam wired in `New`; `handleHead` takes `r`, installs `readCtx`, classifies timeout |
| `host/daemon/handlers.go` | `readCtx` helper; five handlers install it; timeout classification (`ctx.Err()` first); six Internal branches sanitized to `internalErrorMessage` + `errLog` line |
| `host/daemon/daemon_test.go` | `TestBoundedWaitsAndBodyLimit` gains the `readDeadline` literal row + wiring assertion |
| `host/daemon/read_deadline_test.go` (new) | `TestDaemonReadDeadline` (real-store-expired-deadline + blocking-store arms), `TestDaemonReadDisconnect`, `TestReadCtxCancelledAfterHandler` (the §2.3 cancel gate), `TestInternalErrorsAreSanitized`, `TestTimeoutStatusMirrorsSketch` |
| `host/broker/approve.go`, `host/registry/registry.go`, `host/replay/replay.go` | mechanical: `context.Background()` at the 11 remaining getter call sites (V26/V27) — the visible spelling of today's behavior, pinned by the §2.8 ratchet; no semantic change |
| `host/broker/broker.go` | unexported `invokeReplay` gains a `ctx` parameter from `Session.invoke` (broker.go:173); its 2 getter sites pass the caller's real context (V26) |
| `host/transitionreg/transitionreg.go` | its 3 getter sites pass the ctx that `ReadSnapshot`/`Publish` already receive (V26); no signature change |
| `design_docs/sketches/worlddapi.ail` | `Timeout(string)` constructor, `Timeout(_) => 503` arm, one test row — probed verbatim in V13 |
| `docs/QUICKSTART.md` | one short paragraph in the serve section: reads answer 503 class `Timeout` after 10 s; 500 bodies are generic and the detail is on the daemon's stderr. QUICKSTART is executed-verbatim-maintained (S7), so the sprint re-executes the walkthrough before commit |

## 4. What this buys — and what it does not

Buys: one elapsed-time contract on all six store-reading GET routes (and `/v1/health` untouched,
having no wait to bound); every store read forced to state its context at the call site — the
unbounded case becomes explicit, greppable, and pinned by the §2.8 ratchet, though **not**
unrepresentable, since a deadline-free context stays valid at 11 enumerated production sites;
five of the 16 non-daemon sites upgraded to the caller's real context at zero API cost (V26); a
lock-layer retry policy that turns transient contention from an instant 500 into a bounded retry;
an explicit, machine-readable timeout class item 14's renderer can display truthfully; and 500
bodies that no longer leak paths into a future HTML surface.

Does not buy, declared (§10 makes each a named residual, not an omission): a bounded
`POST /v1/commit` (DR-1); real deadlines inside the approve/bootstrap/replay read paths (DR-2 —
they get the visible, ratchet-pinned `context.Background()` spelling); and no claim that the
driver's busy-wait sleep is context-interruptible until the M2 contention test measures which
pre-registered outcome holds (§2.5).

Stated plainly, per the round-2 precision objection (§12, B1): this item establishes (a) a
bounded elapsed-time contract on the daemon's GET surface, (b) an enforced floor under the
remaining deadline-free read set — it cannot grow — and (c) a named, acceptance-shaped queue
item for that set (§10). It does **not** establish repo-wide bounded waits: standing rule 6 is
discharged for the daemon read path only, and §10 states the exposure that remains.

## 5. Persistent non-vacuity

Every property that must survive implementation is attached to a committed Go test in the
`go test ./...` CI leg. One-shot greps appear only as AC snapshots.

| Persistent property | Gate | Positive/control arm |
|---|---|---|
| each getter's pool-wait consults the caller's context | `TestReadGettersHonorContext/<getter>` ×5 | unmutated call returns in µs; watchdog is the red path, so the test cannot hang green |
| deadline is actually installed on the request path | `TestDaemonReadDeadline/real-store-expired-deadline` | same routes answer 200 with a normal deadline in the sibling subtest |
| timeout renders 503 + class `Timeout`, not 500 | `TestDaemonReadDeadline` asserts class and code | `TestTimeoutStatusMirrorsSketch` replays the sketch vector (the mirror idiom, V11) |
| client disconnect releases the store connection | `TestDaemonReadDisconnect` | blocked-then-unblocked signal read from the wrapper's channel |
| busy_timeout reaches every physical connection | `TestProductionDSNSetsBusyTimeout` | `PRAGMA busy_timeout` readback == 2000 on write and read-only handles (driver state, not a Go constant) |
| busy retry is real, not decorative | `TestReadRetriesUnderTransientExclusiveLock` | lock released at 300 ms → row returned; injection removed → instant SQLITE_BUSY |
| 500 bodies carry no internal error text | `TestInternalErrorsAreSanitized` | sentinel appears in the ErrorLog buffer in the same test — the detail is moved, not destroyed |
| every handler releases its CancelFunc on return (no timer leak) | `TestReadCtxCancelledAfterHandler` | the same test's 200 assertion on each route — a cancel that fired *early* would error the delegated store call and red the 200 arm |
| the deadline-free read set cannot grow | `TestNoNewDeadlineFreeStoreReads` | the pinned counts are non-zero (11 total), so the scanner provably sees the existing sites — an empty or mis-rooted scan cannot pass |
| readDeadline stays pinned and non-zero | extended `TestBoundedWaitsAndBodyLimit` | six existing literal rows keep passing unchanged |
| sketch and Go cannot drift on the 503 | `TestTimeoutStatusMirrorsSketch` | Leg 1 still ai-checks the sketch (`check.passed`, `errors=0`) on every CI run |

## 6. Mutation / test-plan table

Per hard rule: each row names the single `-run` selector that must red, and states **which write
the assertion reads** — no row's observable is a value set alongside its mechanism, except MU8,
which is declared as wiring-only and paired with MU1.

| ID | Mutation (file : site) | Reds under | Which write the assertion reads |
|---|---|---|---|
| MU1 | `readCtx`: `context.WithTimeout(r.Context(), d.readDeadline)` → `context.WithCancel(r.Context())` (handlers.go, new helper) | `-run 'TestDaemonReadDeadline/real-store-expired-deadline'` | the HTTP status+class written by the timeout branch; mutated code writes 200 with a world body — a write the timeout mechanism can never produce |
| MU2 | `readCtx`: `r.Context()` → `context.Background()` | `-run 'TestDaemonReadDisconnect'` | the wrapper's unblock signal, closed only by ctx-done propagation; mutated code closes it at the 10 s deadline, after the 2 s watchdog reds |
| MU3 | timeout classifier: `if ctx.Err() != nil` → `if false && ctx.Err() != nil` (handlers.go) | `-run 'TestDaemonReadDeadline/blocking-store'` | the class string in the response body — mutated code writes `Internal`, which only the 500 branch produces |
| MU4a–e | one per getter: `QueryRowContext(ctx, …)` → `QueryRowContext(context.Background(), …)` (store.go:467/522/551/628/802 bodies) | `-run 'TestReadGettersHonorContext/<getter>'` (5 selectors) | the getter's return-within-watchdog carrying the context error — produced by ctx machinery inside `database/sql`, reachable no other way while the single connection is held |
| MU5 | delete the `_pragma=busy_timeout(2000)` injection (writer_lock.go:176/187 region) | `-run 'TestReadRetriesUnderTransientExclusiveLock'` | the successfully-read row after the 300 ms lock release — producible only by the driver's busy retry loop |
| MU6 | injection value 2000 → 0 | `-run 'TestProductionDSNSetsBusyTimeout'` | `PRAGMA busy_timeout` readback from the live connection — driver connection state, not the Go literal that set it |
| MU7 | restore `err.Error()` on an Internal branch (handlers.go:220; the test seeds a sentinel-bearing error) | `-run 'TestInternalErrorsAreSanitized'` | two writes, separately asserted: the response body (sentinel must be absent) and the ErrorLog buffer (sentinel must be present) — a mutation passing one still reds the other |
| MU8 | `New`: `readDeadline` field wired to 0 | extended `-run 'TestBoundedWaitsAndBodyLimit'` | **wiring-only** (reads the constructed field, a value set alongside); declared per the discriminator rule and paired with MU1, whose observable is the mechanism's own write |
| MU9 | Go timeout branch: `http.StatusServiceUnavailable` → `http.StatusInternalServerError` | `-run 'TestTimeoutStatusMirrorsSketch'` | the response status code checked against the sketch's replayed 503 vector |
| MU10 | sketch arm `Timeout(_) => 503` → `Timeout(_) => 500` (worlddapi.ail) | **NOT CI-red — declared, not pretended** (V12: Leg 2 tests `world/` only; sketch inline tests are outside both gate legs). Locally `ailang test sketches/worlddapi.ail` reds (19th row). The enforced guard for this drift direction is MU9's mirror test on the Go side — the same one-sided-pin shape as iteration 85's MU10, recorded honestly | n/a — the row exists to state the gate's true coverage boundary |
| MU11 | delete one handler's `defer cancel()` after `readCtx(r)` (handlers.go, any of the six installation sites) | `-run 'TestReadCtxCancelledAfterHandler'` | the recorded context's `Err()` read after `ServeHTTP` returns — `context.Canceled` is written only by the CancelFunc, and with `readDeadline` at its 10 s default no other writer exists inside test time |
| MU12 | add a getter call passing `context.Background()` in production code (scratch site in handlers.go) | `-run 'TestNoNewDeadlineFreeStoreReads'` | the per-file count the scanner reads from the source tree — a write the mutation makes by existing; the mutated file's count exceeds its pin |

Mutant hygiene carried from the drills: restore from `cp` backups, never `git checkout --`;
mutation builds on test-only edits use `go vet`, not `go build` (which skips `_test.go`).

## 7. Acceptance criteria

Baselines were **observed at `6fd26f0`** (commands and outputs in §11), not assumed. Known trap,
measured live at base: a `-run` selector matching nothing exits 0 with `no tests to run` (V20) —
so every test-running AC asserts the `=== RUN` enumeration, never the exit code alone.

**AC1 — the unbounded spelling is gone.**
`grep -c 'context.Context' host/store/store.go && grep -c 'func (s \*Store)' host/store/store.go`
Baseline: `0` / `14` (V1). Pass: first count ≥ 5, control still ≥ 14. Fails if any getter keeps a
context-free signature. Producible: five signatures plus five `QueryRowContext` calls each name
the type.

**AC2 — store-layer cancellation tests exist and pass.**
`GOTOOLCHAIN=go1.25.6 go test ./host/store -run 'TestReadGettersHonorContext|TestProductionDSNSetsBusyTimeout|TestReadRetriesUnderTransientExclusiveLock' -v -count=1`
Baseline: rc=0, `testing: warning: no tests to run` — the vacuous pass, observed (V20). Pass:
rc=0 AND ≥ 7 `=== RUN` lines covering the five getter subtests and both lock-layer tests. Fails
under any of MU4a–e, MU5, MU6. Producible: `:memory:`/file stores, `db.Conn`, and a second raw
driver connection generate every stimulus without sockets.

**AC3 — daemon-layer deadline tests exist and pass.**
`GOTOOLCHAIN=go1.25.6 go test ./host/daemon -run 'TestDaemonReadDeadline|TestDaemonReadDisconnect|TestReadCtxCancelledAfterHandler|TestInternalErrorsAreSanitized|TestTimeoutStatusMirrorsSketch|TestBoundedWaitsAndBodyLimit' -v -count=1`
Baseline: 5 `=== RUN` lines, all from the existing `TestBoundedWaitsAndBodyLimit`; the new test
functions enumerate 0 (V21). Pass: rc=0 AND `=== RUN` lines for all five new test functions AND
the extended D7 constants row. Fails under MU1, MU2, MU3, MU7, MU8, MU9, MU11. Producible:
`httptest.NewRecorder` + real `:memory:` stores; no socket binds, so the codex-sandbox loopback
restriction cannot mask a result.

**AC4 — handlers consult the request context.**
`grep -c 'r.Context()' host/daemon/handlers.go; grep -c 'r.Context()' host/daemon/daemon.go`
Baseline: `0` and `0`, controls 9 `http.ResponseWriter` (handlers.go) and 8 `mux.HandleFunc`
(daemon.go) (V4, V9). Pass: ≥ 1 in each file (the helper; `handleHead`). Fails if a handler is
wired around the deadline path.

**AC5 — the sanitize decision landed, and only where decided.**
`grep -c 'err.Error()' host/daemon/handlers.go`
Baseline: `11` (V5). Pass: exactly `5`, and the five survivors are the BadRequest sites
(handlers.go:215, 242, 392, 399, 410 at base). Fails if a 500 branch still echoes, or if a 400
branch was over-sanitized. Persistent form: `TestInternalErrorsAreSanitized`.

**AC6 — lock policy exists in production code.**
`grep -rn 'busy_timeout' host/store/*.go | grep -v _test | wc -l`
Baseline: `0`, known-positive same-scope control: the test-input hit at writer_lock_test.go:609
(V6). Pass: ≥ 1, in writer_lock.go's DSN builders. Persistent form: the PRAGMA readback test.

**AC7 — the `.ail` gate holds with identical pins.**
`AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh`
Baseline: rc=0 — `10 required identities verified, 39 named tests pass`, 11-module allowlist,
package gate 9/9, run fresh at this base (V18). Pass: rc=0 with the SAME totals after the sketch
edit — V12 is the measured reason the totals cannot move; if they do move, the sketch edit did
something this design did not price, and the sprint stops. This AC can fail (a malformed sketch
edit reds `check.passed` in Leg 1) and cannot pass identically-vacuously (the edit is swept by
Leg 1 on every run).

**AC8 — the Go gate holds.**
`GOTOOLCHAIN=go1.25.6 go build ./... && GOTOOLCHAIN=go1.25.6 go test ./... -count=1`
Baseline: build rc=0 observed at base (V18); the full test leg was green at `aaada20`, whose Go
tree is byte-identical to this base (V22). Caveat carried from the queue: `host/broker` has a
~18% base flake — a red broker leg is investigated against that base rate, not waved through and
not silently retried into green. Fails on any compile break in the 86-site migration or any new
test.

**AC9 — the deadline-free residue is pinned.**
`GOTOOLCHAIN=go1.25.6 go test ./host/store -run 'TestNoNewDeadlineFreeStoreReads' -v -count=1`
Baseline: rc=0, `no tests to run` — the vacuous shape (V20), which is why the pass condition
counts enumeration. Pass: rc=0 AND 1 `=== RUN` line AND the pinned per-file counts equal
{approve.go: 8, registry.go: 2, replay.go: 1, all else: 0}. Fails under MU12, and fails if the
M1 migration lands a different deadline-free count than this design declares — the pin is a
claim, checked, not a snapshot. Producible: a source scan of the caller packages' non-test
files, rooted via `runtime.Caller` so the result does not depend on the test's working
directory.

## 8. Conflict surface

Files touched: the rows of §3. Gate pins and frozen surfaces, each with its disposition:

| Pin / frozen surface | Disposition |
|---|---|
| `scripts/verify_ail.sh` `EXACT_TOTAL_VERIFIED=10` (line 311) | **not moved** — counts `world/*` only (V12); this item edits zero `world/` files |
| `scripts/verify_ail.sh` `EXACT_TOTAL_TESTS=39` (line 350) | **not moved** — Leg 2 runs `ailang test world/`; the sketch's 19th test row is invisible to it (V12, V13) |
| `LEG1_MODULES` 11-module allowlist | **not moved** — no `.ail` file added or removed; `sketches/worlddapi.ail` is already row 5 |
| Leg 3 `verify_world_package.sh` (9 steps, 4-export manifest, golden) | **untouched** — projects `world/` packages only; no `world/` edit |
| This item touches exactly ONE `.ail` file | `design_docs/sketches/worlddapi.ail` — the "likely touches no `.ail`" expectation is FALSE, for a measured reason: handlers.go:16-18 freezes the envelope as a mirror of the sketch, so an honest new error class must move both sides together (V11, V13) |
| `TestBoundedWaitsAndBodyLimit` (daemon_test.go:202) | **edited** — gains the `readDeadline` literal row and wiring assertion; the six existing rows are unchanged |
| `host/boundary/allowlist_world_test.go` (`wantFileCount = 1`, four protected groups, forbidden-prefix list) | **untouched** — stdlib `context` is not in `forbiddenImportPrefixes` (V16); no file added to `host/boundary`; no import of any forbidden prefix anywhere in this diff |
| `ListenAnnouncePrefix` one-line announce protocol (daemon.go:140) | **untouched** — error logging goes to the new `Config.ErrorLog`, never to announce (V17; iteration-28 deadlock precedent) |
| Frozen `/v1` route table (8 patterns, daemon.go:461-468) | **untouched** — no route added or removed; only handler internals change |
| `POST /v1/commit` semantics | **untouched** — `Commit` keeps its signature; DR-1 |
| `world/types.ail`, `world/*` | **untouched** |
| `docs/QUICKSTART.md` | edited (one paragraph); executed-verbatim rule applies (S7) |
| `host/broker` flake | its two files change mechanically; a green broker test after this diff is not evidence the ~18% flake is gone, and a red one is not automatically this diff's fault — attribute by shape against the base rate |

## 9. Milestones and pricing

Each milestone is independently CI-green and committable; order is M1 → M2 → M3.

- **M1 — store layer (0.5 d).** Context-first getters, `QueryRowContext`, busy_timeout DSN
  injection, the 86-call-site mechanical migration (daemon sites pass `r.Context()`;
  transitionreg passes its existing ctx; `invokeReplay` gains a ctx parameter; the 11 remaining
  sites pass `context.Background()`), `host/store/context_read_test.go` complete (including the
  contention test with its pre-registered outcome record, and the §2.8 ratchet). Green: AC1,
  AC2, AC6, AC8, AC9; daemon behavior unchanged (no deadline yet — a request context without a
  deadline is exactly today's bound).
- **M2 — daemon deadline + explicit status (0.75 d).** `readDeadline` constant/field/helper, six
  handlers wired **each with `defer cancel()` per the §2.3 mandate**, timeout classification,
  sketch `Timeout` arm (the V13-probed edit, verbatim), D7 test extension,
  `read_deadline_test.go` except the sanitize test — including the all-five-getter
  `blockingStore`, the escape-channel release plumbing, and `TestReadCtxCancelledAfterHandler`.
  Green: AC3 (minus sanitize), AC4, AC7.
- **M3 — sanitize + log surface + docs (0.25 d).** `internalErrorMessage`, `Config.ErrorLog`,
  the six-branch sweep, `TestInternalErrorsAreSanitized`, QUICKSTART paragraph + re-execution.
  Green: AC5 and the full set.

Total 1.5 d, the top of the row's 1–1.5 d estimate — the quarter-day over round 1 prices the
cancel-release gate, the widened fake with its release plumbing, and the ratchet. Not priced, and therefore a stop-and-return
if discovered necessary: any `world/` edit, any new route, any store schema change, any change to
`Commit`, any upstream `ailang` change.

## 10. What this item is NOT doing (declared residuals)

- **DR-1**: `POST /v1/commit` keeps only its transport bounds. The write path shares the starved
  single-connection pool, and bounding it interacts with the compare-and-append transaction's
  atomicity — a separate design decision, carried by the named item below, not absorbed.
- **DR-2**: the approve/bootstrap/replay read sites receive the *visible*, ratchet-pinned
  `context.Background()` spelling, not real deadlines (transitionreg and the broker replay path
  DO get real contexts — folded in, §2.2/V26).
- **DR-1 and DR-2 are one NAMED follow-on queue item, proposed for the queue at ratification:
  `w-bounded-waits-operator-and-write-paths`.** Its acceptance shape, declared here so the
  deferral is tracked work with a definition of done, not a §10 note: **(a)** `Store.Commit`
  becomes context-aware with a pinned transaction/pool-wait deadline added to the D7 constants
  table — it has exactly one production caller, handlers.go:413 (V27), so the mechanical cost is
  small and the real work is the atomicity decision; **(b)** the §2.8 ratchet's pinned counts go
  11 → 0 by threading contexts through `DecideApproval` (8 sites — an operator entry point whose
  deadline policy for an *attended approval designed to wait on a human* must be ratified, not
  defaulted), `registry.Bootstrap` (2 sites; sole production caller daemon.go:368, so the daemon's
  startup context is available), and `ReplayEntry`/`ReplayEpisode` (1 site — `GetVerifyResult`
  gains its context in the same move, closing the store's sixth getter, V25); **(c)** with the
  ratchet at zero, land the store-boundary guard rejecting nil/deadline-free contexts —
  gpt5-6-sol's fix, correct as the END state, breaking-by-measurement as an opening move (§2.8);
  **(d)** a test proving no production caller passes `context.Background()` to the store (the
  ratchet asserting zero). Measured full-scope cost that keeps this out of the current 1–1.5 d
  item: the 11 sites sit in 8 enclosing functions behind 3 exported entry points with 40
  entry-point test call sites to migrate (V27), plus the Commit atomicity decision and the
  approval-deadline policy decision — each a ratification-worthy call, not a mechanical edit.
- **Clause-3 exposure remaining after this item, stated plainly**: `POST /v1/commit`'s pool wait
  is still deadline-free — in practice a commit queued behind reads now waits at most ~10 s per
  read (reads time out), but behind another commit it is unbounded; the operator entry points
  (approval, bootstrap, replay) hold 11 deadline-free reads by declared, pinned policy. That set
  cannot grow (AC9) and its elimination has a named item with an acceptance shape.
- No SSE, no `http.ResponseController` write-deadline relaxation, no streaming route (the V28-era
  facts about those remain undisturbed).
- No retry policy above the busy_timeout — the daemon does not re-run timed-out reads.
- No claim that the driver's busy sleep is context-interruptible until M2's contention test
  records which pre-registered outcome held (§2.5).
- No change to the frozen route table, the commit envelope, `world/*`, `tools/launchd/*`, or any
  gate pin.

## 11. Verification Log

Measurement base: **`6fd26f0`** (branch dev, clean tree), 2026-08-14, darwin/arm64. Pinned
binary `/tmp/ailang-v0300/ailang` = AILANG v0.30.0 commit `e37b370` (verified live, V13); the
PATH `ailang` was not used. Controller-supplied rows from the queue item were re-run first-party
before use; rows V2, V7 and V14 refine them. Rows V25–V27 were added in quorum round 2 at the
same base (the tree did not move between rounds). Every negative/zero row carries a same-scope
known-positive control in the same call.

| ID | Claim | Command | Observed output |
|---|---|---|---|
| V1 | store has zero context plumbing; same-file control fires | `grep -c 'context.Context' host/store/store.go; grep -c 'func (s \*Store)' host/store/store.go` | `0`; control `14` |
| V2 | **five** context-free read getters, not four — `GetRegistryHead` is the queue row's omission | `grep -n 'func (s \*Store) \(GetObject\|GetWorld\|GetLogEntry\|SelectedHead\|GetRegistryHead\)' host/store/store.go` | `467 GetObject`, `522 GetWorld`, `551 GetLogEntry`, `628 GetRegistryHead`, `802 SelectedHead` — none takes a context |
| V3 | each getter blocks in a context-free QueryRow | `grep -c 's.db.QueryRow(' host/store/store.go` | `5` |
| V4 | no handler consults the request context; handleHead discards it | `grep -c 'r.Context()' host/daemon/handlers.go; grep -c 'http.ResponseWriter' host/daemon/handlers.go; grep -c 'r.Context()' host/daemon/daemon.go` | `0`; control `9`; `0` (daemon.go:473,496: `_ *http.Request`) |
| V5 | 11 raw `err.Error()` echoes; six are 500 branches | `grep -c 'err.Error()' host/daemon/handlers.go` + read of handlers.go | `11`; Internal at 220, 247, 277, 325, 351, 419; BadRequest at 215, 242, 392, 399, 410 |
| V6 | no production busy_timeout; the only hit is a test-input string | `grep -rn 'busy_timeout' host/store/*.go \| grep -v _test \| wc -l` + the test grep | `0`; writer_lock_test.go:609 is `{":memory:?_pragma=busy_timeout(1000)", true}` — an in-memory-DSN classification input, not a production pragma |
| V7 | DSN construction is centralized — **in writer_lock.go, not store.go** (queue-row citation corrected; conclusion unchanged) | `grep -n 'func resolveDSN\|func writeDSN\|func readOnlyDSN' host/store/*.go` | writer_lock.go:120, 176, 187; store.go:244/252/276/280 are the call sites |
| V8 | one physical connection serializes every read and write | Read store.go:289-310 | `db.SetMaxOpenConns(1)` at :297; comment: the compare-and-append transaction is "the sole serialization point" |
| V9 | route table: 7 GET + 1 POST; health touches no store; head reads SelectedHead in daemon.go | `grep -c 'mux.HandleFunc' host/daemon/daemon.go` + read :459-508 | `8`; registrations at :461-468; `handleHealth` :473 (no store call); `handleHead` :496-508 calls `d.store.SelectedHead()` |
| V10 | the D7 gate is real and observes the transport layer only | Read daemon.go:60-116, :514-522; daemon_test.go:185-235 | six literal constants pinned (5s/30s/30s/120s/30s/10s); all four `http.Server` fields wired in `newServer`; nothing in the test touches handler or store execution |
| V11 | the Go error envelope is a frozen mirror of the sketch | Read handlers.go:16-18; sed worlddapi.ail:83-113 | "Class names and status codes mirror ApiError/httpStatus in the frozen, checked design_docs/sketches/worlddapi.ail exactly"; sketch has 5 constructors, `httpStatus` with 4 test rows; the loopback predicate uses the same mirror idiom with a named replay test |
| V12 | the planned sketch edit moves NO gate pin | Read verify_ail.sh:118-123, 130-146, 236-244, 296-317, 331-350 | manifest comment: "Sketches carry EMPTY required sets and are excluded from the total"; `REQUIRED_VERIFIED` keys are 4 `world/` files; `total_verified` counted under `case "$mod" in world/*)`; `EXACT_TOTAL_VERIFIED=10`; Leg 2 runs `ailang test --format json world/`; `REQUIRED_TESTS` contains no sketch name; `EXACT_TOTAL_TESTS = 39`; `LEG1_MODULES` = 11 incl. `design_docs/sketches/worlddapi.ail` |
| V13 | the exact planned `.ail` edit verifies on the pinned binary | scratch copy of `sketches/{worlddapi,logepoch}.ail`; pristine `ai-check`; apply `Timeout(string)` ctor + `=> 503` arm + test row; `ai-check` + `test` | pristine control: `check.passed=True, errors=0, cex=0`; edited: `check.passed=True, errors=0, cex=0`; tests **18 → 19**, `19 passed, 0 failed`; binary reports `AILANG v0.30.0 / Commit: e37b370` |
| V14 | driver facts: per-connection `_pragma` (busy_timeout applied first); ctx → `sqlite3_interrupt` at statement level; **no driver default busy_timeout** ⇒ lock collision today is an instant SQLITE_BUSY, not a hang | `go.mod`; grep of `$GOMODCACHE/modernc.org/sqlite@v1.54.0` | `modernc.org/sqlite v1.54.0`; sqlite.go:217-231 sorts `busy_timeout` first in `applyQueryParams` (runs per new connection); `interruptOnDone` used at stmt.go:105, 295 and tx.go:71; grep for a busy default: no hit |
| V15 | migration ripple, measured | `grep -rn '\.GetObject(\|\.GetWorld(\|\.GetLogEntry(\|\.SelectedHead(\|\.GetRegistryHead(' --include='*.go' host cmd`, split by `_test` | non-test: **22** sites in 7 files (broker/approve.go 8, broker/broker.go 2, daemon/daemon.go 1, daemon/handlers.go 5, registry/registry.go 2, replay/replay.go 1, transitionreg/transitionreg.go 3); total incl. tests: **86** |
| V16 | stdlib `context` cannot red the boundary gate | grep of host/boundary/allowlist_world_test.go | the gate forbids a registry/HTTP/cloud *prefix list* (`forbiddenImportPrefixes`, :61) with an `fmt` stdlib control at :861; `context` is not in it |
| V17 | there is no daemon error-log surface today; announce is a one-line protocol | grep for logging in host/daemon + cmd/ailang-worldd; daemon.go:136-140 | zero `log.` calls in the daemon; cmd passes `stdout, stderr io.Writer`; `ListenAnnouncePrefix` consumers read exactly one line (iteration-28 deadlock precedent) |
| V18 | base gates are green | `GOTOOLCHAIN=go1.25.6 go build ./...`; `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` | build rc=0; verify gate rc=0: "10 required identities verified, 39 named tests pass", 11-module allowlist, package gate 9/9, compiler pinned by exact bytes |
| V19 | getters return wrapped errors, so a pool-wait ctx error surfaces through `%w` | Read store.go:467-480, 551-575 | each getter wraps the scan error: `fmt.Errorf("store: get log entry %d: %w", index, err)` shape |
| V20 | AC2's baseline, and the vacuous `-run` trap observed live | `go test ./host/store -run 'TestReadGettersHonorContext\|TestProductionDSNSetsBusyTimeout\|TestReadRetriesUnderTransientExclusiveLock' -v -count=1` | rc=0, `testing: warning: no tests to run`, `ok … [no tests to run]` — which is why every test AC counts `=== RUN` lines |
| V21 | AC3's baseline | same command shape against ./host/daemon with the five selectors | `5` `=== RUN` lines, all from the existing `TestBoundedWaitsAndBodyLimit` and its subtests; the four new test functions enumerate 0 |
| V22 | the base's Go tree equals the gate-proven merge commit | `git diff --name-only aaada20..6fd26f0` | 5 files, all under `design_docs/` — zero Go files; iteration 85's full `go test -count=1` green (17 packages) therefore describes this exact Go tree |
| V23 | write/read DSN builders are the correct injection seam | Read writer_lock.go:176-203 | `writeDSN` returns a bare path when `params` is empty, else `fileURI`; `readOnlyDSN` always builds a `file:` URI with `mode=ro`; both are the single seam every file-backed open passes through (store.go:252, 280) |
| V24 | writeAPIError is the single error-writing funnel the new class joins | `grep -c 'writeAPIError' host/daemon/handlers.go` | `22` |
| V25 | the store has SIX `Get*`/`Selected*` methods, not five; `GetVerifyResult` is off the daemon read path | `grep -n 'func (s \*Store) Get\|func (s \*Store) Selected' host/store/store.go`; `grep -rn '\.GetVerifyResult(' --include='*.go' host cmd \| grep -v _test` | GetObject:467, GetWorld:522, GetLogEntry:551, GetRegistryHead:628, **GetVerifyResult:773**, SelectedHead:802; `GetVerifyResult`'s sole production caller is host/replay/replay.go:191 — zero hits in host/daemon |
| V26 | five of the 16 non-daemon getter sites already have a context in reach | read of transitionreg.go:70,221 and broker.go:155-173,352 | `ReadSnapshot(ctx context.Context)` (:70) and `Publish(ctx context.Context, …)` (:221) enclose the three transitionreg sites (:74, :89, :233); `Session.Invoke(ctx, …)` (:155) → `s.invoke(ctx, …)` holds a ctx and **drops it** at :173 calling `invokeReplay(req, decision)` (:352), which encloses both broker.go sites (:361, :396) |
| V27 | full-scope cost of the reviewer's B3 proposal, and what its boundary guard would break | enclosing-function scan of approve.go/broker.go/registry.go/replay.go; caller greps for `DecideApproval(\|registry.Bootstrap(\|ReplayEpisode(\|NewReplaySession(`; `grep -rn '\.Commit(' --include='*.go' host cmd \| grep -v _test`; `grep -rn 'context.Background()' --include='*.go' host cmd \| grep -v _test` | the 11 residual sites sit in decideApproval:159, appendApprovalHead:194, walkApprovalHead:224 (×4), validatePublishApproval:485 (×2), Bootstrap:112 (×2), ReplayEntry:151 (×1) — 8 enclosing functions behind 3 exported entry points (`DecideApproval`, `Bootstrap`, `ReplayEntry`/`ReplayEpisode`); those entry points have ONE in-repo production caller (daemon.go:368, Bootstrap) and 40 test call sites; `Store.Commit` has exactly one production caller (handlers.go:413 — the other 8 non-test `.Commit(` hits are `database/sql` `tx.Commit()` internals); production `context.Background()` today: 10 sites in 8 files — so a store-boundary reject of deadline-free contexts would break 11 production sites the moment M1 lands |
| V28 | **[R2 gemini objection, measured by the CONTROLLER]** §2.3's claim that `handleLogRange` defers all writing to a single terminal `writeJSON` — i.e. that a mid-loop 503 replaces the page cleanly and cannot corrupt a partially-written response | full-function read of `handleLogRange` in `host/daemon/handlers.go`, then `grep -n 'writeJSON\|Flush\|WriteHeader\|w.Write'` over that function body only | **CLAIM TRUE.** The loop accumulates into `items := make([]logEntryResponse, 0, limit)` and the function body contains exactly ONE write on the success path — `writeJSON(w, http.StatusOK, logRangeResponse{Items: items})` as its final statement. Every error path (`from` parse, `limit` parse, the in-loop `Internal` branch) calls `writeAPIError` and `return`s immediately, so no path writes and then continues looping. Zero `Flush`/`WriteHeader`/`w.Write` calls in the function. A mid-loop deadline therefore replaces the whole page cleanly |
| V29 | **[R2 gemini objection, measured by the CONTROLLER]** §2.3's claim that the CLI client's `DefaultClientTimeout` is 30 s | `grep -rn 'DefaultClientTimeout' --include='*.go' .`; then read `host/daemon/daemon.go:105-135` | **CLAIM TRUE.** `defaultClientTimeout = 30 * time.Second` (`daemon.go:110`, D7 table), exported as `DefaultClientTimeout = defaultClientTimeout` (`daemon.go:133`) precisely so the CLI "cannot invent a second, unbounded deadline"; consumed at `cmd/ailang-worldd/cli.go:32` and pinned by `daemon_test.go:222-224` and `main_test.go:162`. So the 10 s read deadline sits inside a 30 s client bound, as §2.3 states |

> **Provenance note (rule 3b(v)):** V28 and V29 were added by the CONTROLLER at iteration 86, not by
> the designer, in response to round-2 gemini objections. Both objections were CORRECT as process
> — the two claims were load-bearing, present-tense and absent from this log, which the mission's
> premise-verification gate forbids — and both premises measured TRUE, so the design is unchanged.
> This discharges gemini's round-2 reject by measurement rather than by argument.

### Non-blocking repository findings

1. daemon.go:456 still says "The seven patterns below" above eight registrations — the same stale
   comment item 14's doc noted; this sprint may fix it in passing (one word), or leave it to
   item 14, but must not count it as scope.
2. The queue row's "SQLite lock acquisition is unbounded" is imprecise in an instructive
   direction (V14): with no busy_timeout the wait is zero, not unbounded — the failure is an
   *undeclared* policy (instant 500 under transient contention), and the unbounded waits live in
   the connection pool and statement execution. The row's remedy is unchanged; its mechanism
   sentence is corrected here so the next reader does not redesign against the wrong layer.

## 12. Quorum verification log (round 2)

Both external reviewers were present; total metered $0.0985. Disposition of every objection,
with classification and evidence. One objection was partially argued rather than conceded (B3);
none was ignored.

| # | Reviewer / source | Classification | Disposition | Evidence |
|---|---|---|---|---|
| A1 | gemini-3-1-pro: `readCtx` callers never told to release the CancelFunc — timer leak | narrow completeness — CORRECT | Fix applied verbatim (§2.3 now mandates `defer cancel()` immediately after `readCtx(r)`) and made enforceable: new layer-4 test `TestReadCtxCancelledAfterHandler` (§2.5) + MU11 + AC3 coverage. A guard became a gate: removing any handler's `defer cancel()` reds a named selector | controller measurement stands: round 1 had 6 `cancel` mentions, zero stating the obligation |
| A2 | gemini-3-1-pro: `blockingStore` overriding ONE getter contradicts driving all six routes — five routes would hit the embedded real store and answer 200 | narrow completeness — CORRECT | Fix applied verbatim: §2.5 layer 3 overrides **all five** getters. The iteration-80 wrap-not-replace analysis was re-checked after the widening and holds (stated at the end of §2.5: no mutation's only witness runs behind the fake) | the route table reads through five distinct getters (V9, V2) — a one-getter fake provably 200s on the other routes |
| B1 | gpt5-6-sol: "unrepresentable at the type level" / bounded-waits claims are FALSE while deadline-free contexts remain valid | narrow precision — CORRECT | Not defended. Every over-claiming sentence rewritten (§2.2, §4): the signature change makes the unbounded case *explicit and greppable*, not unrepresentable; §4 now states plainly what is and is not established. The honest gap is enforced rather than papered: §2.8 ratchet pins the 11 deadline-free sites | 11 production sites deliberately pass `context.Background()` after this item (V26/V27) — the strong claim was false by the doc's own numbers |
| B2 | gpt5-6-sol: watchdogs may red without terminating the blocked call, leaving cleanup hung | narrow completeness — CORRECT | Every watchdog in §2.5 now names its release mechanism and exit bound: layer 1 closes the occupying `*sql.Conn` (pool returns, getter completes, result channel drained before `s.Close()`); layer 3 closes a test-owned `escape` channel every blocking fake selects on (sentinel return, never readable as a pass); layer 2 measured as unable to block (no watchdog claimed); the contention test terminates by construction (deferred lock release + 5 s ctx) | a `t.Error` from a watchdog goroutine releases nothing; with `SetMaxOpenConns(1)` (V8) a parked getter holds the sole connection and `s.Close()` hangs behind it |
| B3 | gpt5-6-sol: 16 deadline-free production sites + unbounded `Commit` mean the bounded-waits axiom is not satisfied; proposed removing DR-1/DR-2 and bounding everything | DIRECTION — premise TRUE, remedy partially folded, boundary argued | **Folded in (measured cheap)**: 5 of 16 sites get the caller's REAL context at zero API change (transitionreg's ctx-bearing signatures; `invokeReplay` un-dropping `Session.invoke`'s ctx — V26); the §2.8 ratchet test pins the remaining 11 so the set cannot grow (AC9, MU12). **Argued (measured expensive)**: the store-boundary reject would break all 11 remaining sites on landing (V27); those sites sit behind 3 exported operator entry points with 40 test call sites and TWO ratification-worthy policy decisions (attended-approval deadline; Commit atomicity) — ≥ 1 additional day on a 1–1.5 d item. **Non-negotiable delivered**: the residue is a NAMED follow-on item with a four-part acceptance shape ending in the reviewer's own guard (§10), and the item's claims were re-worded to match what it delivers (B1) | V26, V27; §10's `w-bounded-waits-operator-and-write-paths` |
| C | controller: line 36's "all five read getters" stated as a STORE property is false — the store has six | narrow precision (enumeration-as-universal) — CORRECT | Claim re-scoped to the daemon read path (§1); `GetVerifyResult:773` named with its sole caller and its out-of-scope reason; new row V25 so the next reader sees the enumeration was deliberate; its context migration assigned to the follow-on item's shape (b) | V25: six methods; `GetVerifyResult` called only from replay.go:191, never from a daemon handler |

## 13. PARKED — the one-word decision for Mark

Two quorum rounds, both `blocked`, **both with both external reviewers present** (`absent_reviewers`
was `[]` in each round, so neither verdict is a pass-with-a-named-hole). Metered total **$0.2299**
(R1 $0.0985 + R2 $0.1314).

Round 2 left exactly **one** live objection. gemini's round-2 reject was a PREMISE objection and is
discharged above by controller measurement (V28/V29 — both premises TRUE, design unchanged). What
remains is `gpt5-6-sol`'s, and it is a dispute about the item's **scope boundary**, not about any
fact: its premise is TRUE and this document already recorded it. Gate 2's narrow-refinement
carve-out requires every remaining objection to carry a concrete fix **and not dispute the design
DIRECTION**; this one disputes the direction, so the controller may not resolve it. Standing rule 2:
never force through a guardrail.

**The dispute in one sentence.** The reviewer holds that an item filed under clause-3 / Standing
rule 6 ("every wait is bounded") cannot claim compliance while `POST /v1/commit` can wait
indefinitely for the single DB connection and 11 store-read sites deliberately carry
`context.Background()` — and that a follow-on item plus a non-growth ratchet do not close that gap.
The document's position, after folding in everything measurably cheap, is that the daemon's seven
GET routes are the item's honest unit, the residue is tracked rather than hidden, and the wider
scope carries two policy questions a design pass cannot settle alone.

| | **Option A — ship the scoped item** | **Option B — re-size to repo-wide bounded waits** |
|---|---|---|
| Scope | The 7 daemon GET routes; 5 of 16 non-daemon sites folded in at zero API cost; ratchet pins the remaining 11 | Additionally: bound `Store.Commit`, thread real deadlines through broker/registry/replay/transitionreg, add the store-boundary guard rejecting deadline-free contexts |
| Cost | **1.5 d**, as designed | **≥ 2.5 d** — the guard breaks all 11 residual sites the moment M1 lands (V27); they sit in 8 functions behind 3 exported entry points with **40** test call sites |
| Blocks item 14? | No — unblocks it at 1.5 d | Yes — item 14 waits a further ~1 d |
| Residual clause-3 exposure | Operator + write paths, tracked as a named follow-on item | None |
| Needs a human policy call first? | No | **Yes, twice**: what bounds an *attended approval* deliberately designed to wait on a human, and how `Commit`'s atomicity trades against a deadline |

**Reply `A` or `B` on the bookkeeping issue.** `A` unparks this doc to sprint-planner as written.
`B` returns it to the designer for a re-sized third round, and the two policy questions above come
back to you first.

Note the second policy question is genuinely load-bearing and is why the controller did not simply
adopt the reviewer's fix: bounding an approval that is *supposed* to wait on a human is a semantic
decision about the World's authority model, not a plumbing choice.

## Related

- [w-workbench-read-only](w-workbench-read-only.md) — item 14, ratified option B; its §3.3 names
  WB-R1 and the exact follow-on contract this document supplies
- [w-worldd-m2](../implemented/w-worldd-m2.md) — Decision 7, the transport-layer bounded-waits
  table this item extends downward
- [coding-standards](../coding-standards.md) — S2 (this is host-boundary work), S3 (§2.7), S6
  (MU10's honest one-sided-pin declaration), S7 (§3 QUICKSTART row)
- `design_docs/sketches/worlddapi.ail` — the frozen API shape; gains the `Timeout` class (V13)
