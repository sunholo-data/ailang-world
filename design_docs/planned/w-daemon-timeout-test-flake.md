# w-daemon-timeout-test-flake — the design said "already-expired"; the test spelled it "1 ns", and in Go "1 ns" is a timer

- **Status**: Planned
- **Item**: queue row 19, `w-daemon-timeout-test-flake`, clause-2 (filed iter-93,
  controller-measured)
- **Estimated**: ~0.25 day (the row's estimate holds — §10)
- **Measurement base**: `0b368e1` (dev HEAD, clean tree), 2026-08-19
- **Instruments**: local `go1.26.4 darwin/arm64` (the controller's M-rows, this session's
  D-rows D1–D14, and the pre-round-2 D17 — text-only sed/grep, toolchain-independent; the
  iteration-94 controller rows D15–D16 and D18 ran under `GOTOOLCHAIN=go1.25.6`);
  gate toolchain `go1.25.6` via `GOTOOLCHAIN` (`verify_go.sh` denylists
  go1.26.0–1.26.5, D10); pinned `ailang` v0.30.0 at `/tmp/ailang-v0300/ailang` (gate runs
  only — **no `.ail` file is touched by this item**)
- **Files touched**: `host/daemon/read_deadline_test.go` (~+35/−2: one named constant, two
  stimulus-site edits, one new ~20-line mechanism-pin test), `host/daemon/handlers.go`
  (**comment-only**, ~+18/−0, zero behavioural bytes — the §4 LIMITATION block; D11/D13 prove
  no gate is moved by it)

This doc is built on the iteration-93→94 controller measurements (cited as M1–M4, provenance:
VERIFIED BY CONTROLLER at HEAD `0b368e1`, darwin/arm64, go1.26.4,
`AILANG_BIN=/tmp/ailang-v0300/ailang`), on first-party re-derivations and extensions made in
this design session (D1–D14, §12), and on two controller measurements run in answer to quorum
round 1's two premise objections (D15–D16, VERIFIED BY CONTROLLER, iteration 94 — §12; what
each objection was and what its measurement changed in this doc is recorded in §13), on a
pre-round-2 CONTROLLER-found defect in AC5 itself, measured in both arms and re-derived
first-party as D17 (§13), and on a controller measurement run in answer to quorum round 2's
one remaining objection — the reviewer's own proposed_fix's second arm, executed (D18, §12;
round 2 is logged in §13). Where the queue row is wrong it is corrected plainly (§1);
where this doc's chosen arm narrows what a test exercises, the narrowing is named, priced, and
carried by an explicit in-code limitation with a named successor (§3, §4).

---

## 1. Problem, corrected from the queue row

The daemon's two read-deadline timeout tests are intermittently RED at base, on a clean tree,
with the failure shape *"a timed-out read answered 200, want the sketch's 503 for Timeout"*.
This makes `go test ./...` — the command every "nothing lands red" gate in this repo runs —
non-deterministic, and a future controller will be tempted to misattribute the red to their own
diff (the row's own framing, and the item-16 class arriving in `host/daemon`).

The row is right about the class, the shape, and the proof standard, and wrong in one
load-bearing place:

| Row claim | Status at `0b368e1` | Evidence |
|---|---|---|
| the affected test is `TestTimeoutStatusMirrorsSketch` | **UNDER-COUNTED: there are TWO tests, sharing ONE stimulus.** `grep -n "time.Nanosecond" host/daemon/read_deadline_test.go` → `:254` (inside `TestDaemonReadDeadline/real-store-expired-deadline`) and `:526` (inside `TestTimeoutStatusMirrorsSketch`). Both are item-18 M2 pins on the 503/`Timeout` contract; a fix that repairs one and not the other leaves the gate exactly as non-deterministic as today. Every acceptance criterion below is scoped so a one-site repair cannot satisfy it (§8 AC5). | M1, D2 |
| fails ~1 in 20 / ~1 in 10 | **CONFIRMED at higher counts**: `real-store-expired-deadline` 4/500 (~0.8%); `TestTimeoutStatusMirrorsSketch` 13/2000 runs, 14 route-level `answered 200` (~0.65%; one run failed on TWO routes). Route-agnostic — failures landed on `world`, `object`, and `log range` across runs. | M2 |
| timing-driven, pre-existing, not M3's doing | **CONFIRMED** — the mechanism is in §2 and it is as old as the tests themselves. | M3, D6 |
| "prove the fix with `-count` high enough that the measured rate would have shown" | **ADOPTED, with the arithmetic the row asks for** (§8 AC3/AC4). The row learned item 16's `-count=20` lesson: at p̂=0.0065 a 20-run sample stays green 87.8% of the time, and M2's own first sample was 0/20 — proof of nothing (0.992²⁰ ≈ 85.2% green at the real-store arm's rate). | M2, §8 |
| "Do NOT weaken the assertion — a 200 here means the timeout contract genuinely did not hold" | **HONOURED — and answered rather than dodged.** Both tests' assertion bodies are byte-unchanged by this item. What the 200 says about the *production* contract is a real question and gets its own section (§3) and an in-code LIMITATION (§4) instead of silence. | §3 |

## 2. The mechanism: why `1 ns` flakes and why `≤ 0` cannot

The shipped read path (D5; the `readCtx` body quoted below and the 10 s package constant
`readDeadline` are pinned verbatim by D16 — quorum round 1 correctly found the doc asserting
both with no row): all six store-reading GET routes — five handlers in `handlers.go`
plus `handleHead` in `daemon.go:589`, the file split that was iteration 93's whole lesson — run

```go
ctx, cancel := d.readCtx(r)      // context.WithTimeout(r.Context(), d.readDeadline)
defer cancel()
v, ok, err := d.reads.GetX(ctx, ...)
if err != nil {
    if timedOut(ctx, err) { writeReadTimeout(...); return }
    d.writeInternalError(w, r, err); return
}
```

`timedOut` is consulted **only on the error path**. The tests shrink `d.readDeadline` to
`1 * time.Nanosecond` — a **positive** duration. `context.WithTimeout` is
`WithDeadline(parent, time.Now().Add(timeout))` (go1.25.6 `context.go:703`), and
`WithDeadline` branches on `dur := time.Until(d)`:

- `dur > 0` (the shipped stimulus): the context is returned **live** and a `time.AfterFunc`
  timer is armed (the lines immediately below `:645`, D6). Cancellation happens on a timer goroutine ~microseconds
  later. If the `:memory:` SQLite read completes first — and it can, the captured 200 bodies are
  real data — then `err == nil`, the classifier never runs, and the route answers 200 with the
  real body. That is the 0.65–0.8% (M2, M3).
- `dur <= 0`: `WithDeadline` **cancels synchronously at construction** —
  `c.cancel(true, DeadlineExceeded, cause)` before returning, no timer, no goroutine
  (`context.go:645`, identical in go1.25.6 and go1.26.4, D6). `ctx.Err()` is
  `context.DeadlineExceeded` before the store call begins.

With an already-cancelled context the read cannot succeed: `database/sql`'s connection
acquisition opens with `select { default: case <-ctx.Done(): return nil, ctx.Err() }`
(`sql.go`, `func (db *DB) conn`, D7), and the Go spec guarantees the `default` case is taken
**only when no other case can proceed** — a closed `Done` channel always wins. So the getter
deterministically returns an error whose chain is `context.DeadlineExceeded` (the store wraps
with `%w`, D8), `timedOut` deterministically returns true (its first arm reads the same
`ctx.Err()`), and the route deterministically answers 503/`Timeout`. This is a
language-semantics claim paired with an empirical arm: M4 measured the changed stimulus at
**0 failures in 2,500 runs** across both tests.

The candidate fix M4 measured — `-1 * time.Nanosecond` at both sites — is therefore not a
tuning of the race; it moves the stimulus onto a categorically different branch of the runtime,
one with no timer to race.

## 3. The design question: is the 200 a test bug or a production bug?

The 200 the flake produces is real data from a read that genuinely succeeded. The daemon's
deadline is enforced by context propagation into `database/sql`, which checks at connection
acquisition (D7) and — for the one shape measured — via the driver's mid-flight
interrupt (D15: the pinned `modernc.org/sqlite v1.54.0` interrupted a CPU-bound read that
provably runs ≥20 s uncancelled, 2.13 ms after a 300 ms deadline). D15 verifies prompt
cancellation only for one CPU-bound query using modernc.org/sqlite v1.54.0 on darwin/arm64.
It does not establish a general bound for lock-blocked or other driver waits; therefore this
item makes no general production bounded-wait claim. There is no post-read check. So the
production residual is two things, not one: (i) **a read that COMPLETES SUCCESSFULLY before
the cancellation is observed answers 200 even though it overran the deadline**; (ii) **a
LOCK-blocked read is bounded by busy_timeout, not by readDeadline** (D18: under a 300 ms
context deadline a lock-blocked read terminated at 2.042929083 s — bounded, but by
busy_timeout's ~2 s, the deadline exceeded 6.8× while the busy-retry loop ran, and the error
still surfacing as `context deadline exceeded`, so the wire class would be Timeout while the
TIMING was governed by a different mechanism entirely). Today the composition is still safe
only because busy_timeout (2 s) < readDeadline (10 s) — an ordering NOTHING in the code
asserts; an ordering, not a guarantee. Both residuals are about the **production contract**,
and the row forbids papering over them. Three arms:

**ARM A — test-only: make the stimulus deterministically already-expired.** Chosen. Reasons,
in order of weight:

1. **The row prescribes it.** The SCOPE sentence — *"make the stimulus deterministic rather
   than racing (the deadline must be observed before the read can answer)"* — is the `dur <= 0`
   branch described in English: with construction-time cancellation, the deadline IS observed
   (by `database/sql`'s acquisition check) before any read can begin. The assertion is not
   weakened; it is byte-unchanged.
2. **It is a restoration, not a narrowing.** The item-18 design doc's own words for this arm
   are *"with an **already-expired deadline** against a real store, the correct arm answers 503
   in microseconds and every mutated arm answers 200 in microseconds — neither can block"*, and
   it claimed the arm **Deterministic** (D9). The implementation spelled "already-expired" as
   `1 ns`, which in Go means "expires 1 ns from now" — a future deadline, timer-armed. The
   ratified design's determinism claim was true of the mechanism it named and false of the
   spelling that landed. The strongest evidence that this was a spelling error and not a
   choice: the **same sprint's store-layer test got it right** —
   `host/store/context_read_test.go:121` constructs its expired context as
   `context.WithDeadline(context.Background(), time.Now().Add(-time.Hour))`, expired at
   construction (D9). ARM A makes the daemon-layer file say what its sibling and its design doc
   already say.
3. **What it narrows, priced honestly.** The 1 ns stimulus occasionally (0.65–0.8%) sampled
   "the deadline fires while a real read is in flight". Post-change these two tests never
   sample that shape. But (a) a 0.8% sample was never *coverage* — a gate that exercises its
   scenario one run in ~130 is the per-run vacuity S6 names, the same argument that killed item
   16's `-count=20`; (b) the shape "deadline fires DURING a read" **is** covered
   deterministically, at the read seam, by the existing `blocking-store` subtest — a getter
   parked until `ctx.Done()` fires at 50 ms, through the production classifier, with an
   interrupt-shaped error that discriminates `timedOut`'s arms, which the real-store arm cannot
   do anyway (§5, D4, D8); and (c) the store layer covers "expiry during a real pool wait"
   deterministically with a real store (`TestReadGettersHonorContext`'s occupied-connection
   design, D9). What is genuinely uncovered — a real read that completes
   successfully before the cancellation is observed, answering 200 despite overrunning its
   deadline — cannot be forced deterministically by any committed test, because you cannot make
   a real `:memory:` read reliably slow. That residual, together with D18's (the deadline does
   not govern lock waits; busy_timeout does), is exactly the LIMITATION of §4.

**ARM B — production: after a successful read, answer 503 if the read's context has expired.**
Rejected here, for reasons that are measurements rather than taste:

1. **B does not fix the flake.** A post-read `ctx.Err()` check is itself timer-dependent: under
   the current 1 ns stimulus the `AfterFunc` may not have run when the check executes, so the
   handler can still read `ctx.Err() == nil` after a fast read and answer 200 (D6 — `Err()`
   reports what `cancel` set; nothing computes expiry lazily). B is deterministic only if the
   stimulus is expired at construction — i.e. B **requires A anyway**; "B alone" is not a
   reachable design point, only C is.
2. **The behaviour change is real and adverse at the shipped deadline.** At the shipped 10 s
   (`readDeadline`, D16), the only reads the new branch converts are those that completed
   within the cancellation-observation lag — measured at 2.13 ms past the deadline (D15) — so
   the client has already waited ~the deadline, the result is real and current, and a 503
   discards work the client will immediately re-request at full price. The deadline's ratified
   purpose is Standing rule 6 — a bound on the *wait* — and D15 measures a prompt mid-flight
   bound for exactly one shape (one CPU-bound query, one driver version, one platform; it
   establishes no general bound for lock-blocked or other driver waits, and on the lock path
   D18 measures the bound that fires being busy_timeout, not the deadline);
   "no request whose deadline passed ever sees 200" is a *stronger, different* contract.
3. **The new branch is untestable without a contract-violating fake.** With an
   expired-at-construction stimulus the read errors at acquisition, so the post-read branch is
   dead to every deterministic committed stimulus in the repo. The only way to exercise
   "succeeded despite expiry" on demand is a store wrapper that ignores its context and
   delegates with a live one — a fixture modelling precisely the misbehaving store item 18
   existed to eliminate. Buildable, but it means the branch's only gate rides on a fake
   violating the seam's contract — a G5-shaped hazard that deserves its own reviewed design,
   not a rider.
4. **Authority.** `timedOut`'s error-path-only placement is a shipped, evaluator-scored,
   three-iteration-old design statement (its own doc comment says so). Reversing it inside a
   test-flake row scoped "make the stimulus deterministic" is a direction change smuggled past
   the quorum that reviewed item 18.

**ARM C — both.** Rejected with B: it inherits B's items 2–4 in full and blows the 0.25 d band
(seven production sites — six handlers plus the log-range loop — each needing mutation
coverage, plus the fake of B.3).

**Did D15 change this choice? Re-argued, not defended.** The strongest honest case for ARM B
was never the flake (B.1: B does not fix it). It was the possibility the quorum's first
objection named: that the mid-flight bound did not exist, in which case B's post-read check
would have been the only enforcement a long in-flight read ever met, and "the deadline bounds
the wait" would have been false. D15 closes that possibility for the one shape it measures —
a CPU-bound query on the pinned driver and platform, interrupted 2.13 ms after the deadline —
and no further: it establishes no general bound for lock-blocked or other driver waits. For
that measured shape, the only thing B adds is the single conversion B.2 prices — turning the
narrow completes-before-cancellation 200 into a 503 that discards completed work — while B.3's
contract-violating fixture and B.4's authority problem are unchanged by the measurement. D15
therefore makes ARM B *less* necessary, not more, on the shape it covers. ARM A stands.

**Did D18 change this choice? Re-tested, not defended.** D18 (quorum round 2 — §12, §13)
measured the lock-blocked path the reviewer's proposed_fix demanded: the read TERMINATED —
every path measured terminates within a stated bound — but the bound that fired was
busy_timeout (~2 s), not the 300 ms request deadline, which was exceeded 6.8× while the
busy-retry loop ran; the error still surfaces as `context deadline exceeded`, so the wire
class would be Timeout while the TIMING was governed by a different mechanism entirely.
Re-tested against this, ARM B does not become more attractive: B's post-read check runs only
when the read SUCCEEDS (err nil, by B's own definition), and D18's lock-blocked read surfaces
an error — B's branch cannot run on that path and changes nothing D18 measured. What D18
changes is what this doc may CLAIM (§3 above, §4.3), and it names a residual for the
controller to file (§11): today the composition is still safe only because busy_timeout (2 s)
< readDeadline (10 s) — an ordering NOTHING in the code asserts; an ordering, not a
guarantee. ARM A stands, re-tested against D18.

**The mandated consequence of choosing A**: the daemon does NOT enforce the deadline's *status
contract* on a read that completes before the cancellation is observed — such a read answers
200. No general bounded-wait claim is attached to that statement: D15 verifies prompt
cancellation only for one CPU-bound query using modernc.org/sqlite v1.54.0 on darwin/arm64,
and a LOCK-blocked read is bounded by busy_timeout, not by readDeadline (D18). Both residuals
are stated as a named, in-code LIMITATION (§4) with a named successor
item, **`w-daemon-late-read-503`**, which — if the stronger contract is ever wanted — would land
the post-read check at all seven sites with its own design, quorum, and the B.3 fixture
question answered properly. Filing that row is a controller/record action; this doc names the
identifier so the comment and any future row point at one string.

## 4. M1 — the change (one milestone, independently landable)

### 4.1 One named constant, both sites

In `host/daemon/read_deadline_test.go`, replacing the two literals at `:254` and `:526`:

```go
// expiredReadDeadline is the deterministic timeout stimulus for the read-
// deadline tests: any NON-POSITIVE duration makes context.WithTimeout take
// context.WithDeadline's `dur <= 0` branch, which cancels the context
// SYNCHRONOUSLY at construction — no timer, no goroutine, no race — so the
// store read is refused at connection acquisition and the 503/Timeout branch
// runs on every route, every run. A small POSITIVE duration (the previous
// `1 * time.Nanosecond`) is a FUTURE deadline: it arms a time.AfterFunc, and a
// fast read can complete before the timer goroutine runs, answering 200 with a
// real body (measured at base: ~0.65–0.8% of runs). Do not "shrink" this back
// to a positive value; TestExpiredReadDeadlineExpiresAtConstruction reds on
// the sign, and the design doc for w-daemon-timeout-test-flake holds the
// measurements.
const expiredReadDeadline = -1 * time.Nanosecond
```

with both stimulus lines becoming `d.readDeadline = expiredReadDeadline`. One constant, not two
edits, so the two tests cannot diverge again without moving a count an AC pins (§8 AC5).

### 4.2 The mechanism pin

A new test in the same file, through **production** `readCtx` (it wraps, it does not replace):

```go
// TestExpiredReadDeadlineExpiresAtConstruction pins the property every 503
// assertion in this file now rests on: the stimulus context is DEAD BEFORE any
// store read can begin. Two assertions, two mutations:
//   - the sign check kills the "shrink it back to a positive nanosecond"
//     mutation deterministically (no timing anywhere);
//   - the ctx.Err() check goes through the production readCtx, so a readCtx
//     that ignores d.readDeadline (or re-derives from the 10s constant) reds
//     here in one run.
func TestExpiredReadDeadlineExpiresAtConstruction(t *testing.T) {
	if expiredReadDeadline >= 0 {
		t.Fatalf("expiredReadDeadline = %s, want a negative duration — a positive value arms "+
			"a timer and re-creates the 200-vs-503 race this constant exists to remove",
			expiredReadDeadline)
	}
	d := newHandlerDaemon(t)
	d.readDeadline = expiredReadDeadline
	ctx, cancel := d.readCtx(httptest.NewRequest(http.MethodGet, "/v1/head", nil))
	defer cancel()
	if ctx.Err() == nil {
		t.Fatalf("readCtx under an already-expired deadline returned a LIVE context — the "+
			"stimulus must be expired at construction, before any store read can begin")
	}
	if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Fatalf("ctx.Err() = %v, want context.DeadlineExceeded", ctx.Err())
	}
}
```

All imports it needs are already in the file (D2's Read); `newHandlerDaemon` is the file's own
constructor idiom; the test name is unallocated (D13). No fake anywhere: the asserted value is
produced by the context package through the production helper.

### 4.3 The LIMITATION, in the code it limits

D15 verifies prompt cancellation only for one CPU-bound query using modernc.org/sqlite
v1.54.0 on darwin/arm64. It does not establish a general bound for lock-blocked or other
driver waits; therefore this item makes no general production bounded-wait claim. The
residual is two things, not one: (i) a read that completes before the cancellation is
observed answers 200; (ii) a LOCK-blocked read is bounded by busy_timeout, not by
readDeadline (D18). Appended to `timedOut`'s doc comment in `handlers.go` (comment-only;
token collision checked, D13 and the round-2 re-check in §13):

```go
// LIMITATION(w-daemon-late-read-503): this classifier is consulted ONLY on the
// error path — on a successful read err is nil and this function never runs.
// D15 verifies prompt cancellation only for one CPU-bound query using
// modernc.org/sqlite v1.54.0 on darwin/arm64. It does not establish a general
// bound for lock-blocked or other driver waits; therefore this item makes no
// general production bounded-wait claim. The residual is two things, not one:
//   (i) a read that COMPLETES SUCCESSFULLY before the cancellation is
//       observed answers 200 with the real body even though it overran the
//       deadline;
//  (ii) a read blocked on a LOCK is bounded by busy_timeout, NOT by the
//       request deadline (w-daemon-timeout-test-flake D18: under a 300ms
//       context deadline a lock-blocked read returned after 2.04s, governed
//       by busy_timeout's ~2s, the deadline exceeded 6.8x while the
//       busy-retry loop ran — and the error still surfaces as
//       deadline-exceeded, so the wire class is Timeout while the timing was
//       governed by a different mechanism entirely). Today that composition
//       is safe only because busy_timeout (2s) is shorter than the 10s
//       request deadline — an ORDERING nothing in this code asserts, not a
//       guarantee.
// Enforcing a stronger status contract means a post-read expiry check at all
// seven read sites and a decision that completed work must be discarded; that
// is the successor item named above, not a rider on a test fix. The
// read-deadline tests deliberately use an ALREADY-EXPIRED stimulus
// (expiredReadDeadline) precisely because a near-expiry one races the
// completes-before-cancellation window.
```

The comment deliberately avoids the literal `err.Error()` — item 18's AC5 grep counts exactly 5
of those in this file today, and AC6(b) proves the count is unmoved (the iteration-92 lesson: a
comment can move a raw-grep observable). It likewise spells neither of this item's own AC5
stimulus tokens, so the package-scope sweep still reads exactly its expected count, and it
avoids the bare token `readDeadline` so AC5's stripped-vs-raw instrument control on this file
(stripped 6 / raw 7, D17) is unmoved. It DOES spell `busy_timeout`: checked before landing —
every recorded grep counting that token is scoped to `host/store/*.go` (item 18's V6/AC6),
not to this file, and this item's own AC5/AC6(b) patterns do not match it; the check, with
numbers, is in §13.

## 5. What feeds the classifier (the G5 audit)

For every assertion this item touches, what else in the system can produce the asserted value?

- **Under the fixed stimulus, BOTH arms of `timedOut` are saturated.** The getter's error is
  the store's `%w`-wrap of `database/sql`'s `ctx.Err()` — i.e. its chain carries
  `context.DeadlineExceeded` (D7, D8) — while `ctx.Err()` is non-nil by construction. So the
  real-store arm **cannot discriminate** a single-arm neutering of `timedOut` (MU3's class):
  neuter either arm and the other still classifies. This is not a regression — the 1 ns
  stimulus had the identical property on its passing runs — and the discrimination burden
  stays where item 18 put it: the `blocking-store` subtest, whose `errStoreInterrupted`
  sentinel deliberately does NOT wrap `context.DeadlineExceeded` so only the `ctx.Err()` arm
  can classify it (D4/D8). Stated so nobody reads AC3's determinism as classifier coverage.
- **The mechanism-pin test** asserts a value (`ctx.Err()`) produced by the context package via
  production `readCtx`; the only test-supplied input is the field write — the same stimulus
  shape every arm in the file uses. Nothing echoes a fixture's value back to it.
- **No new fake is prescribed anywhere in this item.** The two stores the file already wraps
  (`blockingStore`, `recordingStore`) are untouched.

## 6. Mutation table

Landed-proof rule (inherited): sha256 of every mutated file must move before a verdict is read;
one-shots restore from `cp` backups, never `git checkout --`; **every mutant must pass
`go vet ./host/daemon` before its verdict is read** (test-file mutants are invisible to
`go build ./...` — the recorded iter-77 gotcha; the one production-file mutant also gets
`go build ./host/daemon/`). "The mutant does not compile" proves nothing.

| # | Mutation (exact edit, file) | Killed by — deterministically, no timing | Inverse arm (expected red SET, exactly) |
|---|---|---|---|
| MU-STIM-POSITIVE | `expiredReadDeadline = -1 * time.Nanosecond` → `= 1 * time.Nanosecond` (one-shot, `read_deadline_test.go`) | `-run '^TestExpiredReadDeadlineExpiresAtConstruction$'` → exactly 1 `--- FAIL:`, the sign Fatal. Zero timing dependence — this is the arm that makes a 0.8% race killable at `-count=1`. | `-skip` the killer, package `-count=1`: expected rc=0. Residual, priced: the mutant IS the base flake, so ≈1.44% of inverse runs red with the 200-want-503 shape — that red is the mutation's own phenotype, attributed by shape, re-measured once (the committed tests never retry; the operator's re-measurement of a one-shot arm is not a test retry). |
| MU-DEADLINE-DETACH | `readCtx` body: `d.readDeadline` → `readDeadline` (the 10 s package constant; production one-shot, `handlers.go`; compiles — both identifiers in scope, pinned by D16: the constant at `daemon.go:128`, the field-consuming `readCtx` body at `handlers.go:269-271`) | the same `-run` → exactly 1 `--- FAIL:`, the "returned a LIVE context" Fatal (field −1 ns ignored, 10 s future deadline → live ctx). | `-skip` the killer: **multiply-killed, and the red set is the assertion** — `TestDaemonReadDeadline/real-store-expired-deadline` and `TestTimeoutStatusMirrorsSketch` red deterministically (every route 200: the M2 failure shape at 100% instead of <1%), `normal-deadline-answers-200` and everything else green. Any red outside that set is unexplained and fails the arm. |
| MU-SITE-REVERT | the `TestTimeoutStatusMirrorsSketch` stimulus line → `d.readDeadline = 1 * time.Nanosecond` (bypassing the constant; one-shot, `read_deadline_test.go`) | **AC5's package-scope sweep**: the comment-stripped count reads 2 — the constant plus the reverted line, one positive — against an expected 1 → red. Stated plainly: **no committed test can kill this mutant** — a `-count=1` run of the weakened test is green 99.35% of the time. This row is iteration 93's lesson expressed as a mutation: the only per-commit detector for a re-introduced sub-1% race is a static sweep, so the sweep is an AC and its scope is the whole package, not a file. | no killer test to skip; the vet arm and byte-identical restore still apply. Package `-count=1` expected green (the mutant is the base flake at one site, ≈0.65%). |

Not in the table, and why: the §4.3 comment has no behaviour to mutate — a mutation whose only
observable is a grep of prose adds nothing beyond AC6(a) itself. And no mutation targets "the
suite still catches a removed deadline": that is item 18's shipped MU1, whose kill (200 on
every route) is deterministic at base and stays deterministic here — re-proven in passing by
MU-DEADLINE-DETACH's inverse arm.

## 7. Determinism proof standard (the G6 arithmetic, once, used by the ACs)

Measured base rates (M2): real-store arm p̂ = 4/500 = 0.008; mirrors-sketch p̂ = 13/2000 =
0.0065. Detection probability at count n is 1−(1−p̂)ⁿ:

| Arm | n | P(≥1 failure at base rate) |
|---|---|---|
| real-store | 20 | **14.8%** — and M2's actual first sample was 0/20: proof of nothing |
| real-store | 500 | 98.2% (M2's measuring run: 4 failures) |
| real-store | **1000** (AC3) | **99.97%** |
| mirrors-sketch | 20 | 12.2% |
| mirrors-sketch | **2000** (AC4) | **99.9998%** (M2's measuring run: 13 failures) |

So AC3/AC4 run at counts under which the base tree reds with ≥99.97% probability — the ACs can
fail, and a green at head is information. The high counts prove the **measured mechanism at its
measured rate** is gone; the determinism claim itself does not rest on sampling but on the
`dur <= 0` branch (D6), the acquisition refusal (D7), and the committed mechanism pin (AC2) —
which is why no `-count` loop is committed to CI (§9).

## 8. Acceptance criteria

All from repo root; `export AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.25.6` on every
command (`verify_go.sh` FATALs without the pin and denylists the local go1.26.4 — base
conditions, D10). Every `-run` command's verdict is its **counted** `--- PASS:`/`--- FAIL:`
lines, never rc alone (a `-run` matching nothing prints a vacuous rc=0 `PASS` — the inherited
V-G rule). `grep` rc semantics per G1: 1 = no match, 2 = no such path; counts are read from
printed output, never from a pipe that eats the exit code. Each AC states its scope and what it
cannot see.

- **AC1 — both gates green, before and after.** `./scripts/verify_ail.sh` rc=0 and
  `./scripts/verify_go.sh` rc=0 on the pristine tree and again after the change. Base reading,
  stated honestly: green at `0b368e1` *minus this item's own flake* — each full `verify_go.sh`
  runs both affected tests in both legs, so the base gate reds on exactly these two tests with
  probability ≈2.9% (per leg 1−(0.992·0.9935) ≈ 1.44%; two legs ≈ 2.87%). A base red matching
  the 200-want-503 shape is attributed AT BASE (rule 3d) and the gate re-run once; at head that
  shape is extinct and any red here is real. *Cannot see: anything test-specific — that is
  AC2–AC5's job.*
- **AC2 — the mechanism pin exists and passes.**
  `go test ./host/daemon -run '^TestExpiredReadDeadlineExpiresAtConstruction$' -v -count=1` →
  exactly 1 `--- PASS:` line. Base: **0** such lines — the test does not exist, and the
  counted line (not rc) is the verdict for exactly that reason. *Cannot see: whether the two
  big tests actually use the constant it pins — AC5 carries that.*
- **AC3 — real-store arm deterministic at a count the base could not survive.**
  `go test ./host/daemon -run '^TestDaemonReadDeadline$/^real-store-expired-deadline$' -count=1000 -v`
  → exactly 1000 counted `--- PASS: TestDaemonReadDeadline/real-store-expired-deadline` lines
  and zero `--- FAIL` lines. Base: 4/500 measured (M2); at that rate this AC reds with ≈99.97%
  probability (§7). *Cannot see: the second test (AC4), the sibling arms
  (`normal-deadline-answers-200`, `blocking-store` — exercised by AC1's full gates), or a
  future re-introduction of the race (AC5).*
- **AC4 — mirrors-sketch deterministic.**
  `go test ./host/daemon -run '^TestTimeoutStatusMirrorsSketch$' -count=2000 -v` → exactly 2000
  `--- PASS:` lines, zero `--- FAIL`. Base: 13/2000 measured (M2); P(red at base) ≈ 99.9998%
  (§7). *Cannot see: same blindnesses as AC3, mirrored.*
- **AC5 — the stimulus class is gone at PACKAGE scope, and wider — counted with comments
  STRIPPED, because §4.1's own comment spells the token.** (a)
  `sed 's://.*::' host/daemon/*.go | grep -c 'Nanosecond\|Microsecond'` → **exactly 1**: the
  `expiredReadDeadline` definition line, whose matched text contains the leading `-` (the glob
  IS the recursive scope: `host/daemon` has no subdirectory, D17). Base: **2** — the two
  `= 1 * time.Nanosecond` stimulus lines, which are code and survive the strip (D17; D2/D3
  locate them). The count MOVES, 2 → 1, and each side is a different code state. (b) The same
  stripped count at the wider scope:
  `find host cmd -name '*.go' -exec sed 's://.*::' {} + | grep -c 'Nanosecond\|Microsecond'`
  → **exactly 1**; base: **2**, measured at that scope, not assumed from (a); the head arm is
  composed — the changed package reads 1 and all of `host/ cmd/` outside `host/daemon` reads
  0 (D17). Instrument control, same strip, proving the strip removes comment hits WITHOUT
  eating code hits: on `handlers.go`, stripped `readDeadline` = 6 vs raw = 7 — the one comment
  mention vanishes, the six code hits survive (D17). Secondary reading, RECORDED AND DECLARED
  NOT THE GATE: the raw form (`grep -rn 'Nanosecond\|Microsecond' host/daemon/`) reads **2 at
  base AND 2 at head** — §4.1's comment quotes the retired spelling in prose, so the raw count
  cannot move, cannot fail for the right reason, and would red any future editor who merely
  mentions the token; that is exactly the class iteration 92 ("a comment can quote a future
  code state") and iteration 93's E-row ("AC5's grep counts comments") already paid for, which
  is why the stripped form is the gate. The scope is written down
  because iteration 93's headline gate passed on a leaking tree by being scoped to one file
  while the seventh site lived in another; here the claim is package-wide and repo-`host/`-wide,
  so the AC is measured at both. *Cannot see, declared:* a future racy stimulus spelled at
  `Millisecond` scale or through a variable, as before; and — new, bought by the strip — the
  strip is TEXTUAL, not syntactic: it deletes everything right of the first `//` on a line,
  including real code after a `//` inside a string literal, so a racy stimulus spelled in a
  region the strip discards is invisible to the gate (the raw secondary reading above is the
  only record that would show it, and it is not gating). The 50 ms at `:271` is deliberately
  outside the pattern — a positive deadline is safe there because the parked getter cannot
  return until `ctx.Done()` fires (D4). The residual — no committed guard enforces this
  forever — is priced and accepted in §10, not silent.
- **AC6 — the limitation is in the code, and the comment edit moved nothing else.** (a)
  `grep -c 'LIMITATION(w-daemon-late-read-503)' host/daemon/handlers.go` → 1 (base: 0, D13;
  same-call known-positive control: `grep -c 'writeReadTimeout' host/daemon/handlers.go` → 7,
  base and head). (b) `grep -c 'err.Error()' host/daemon/handlers.go` → **still 5** — the live
  count behind item 18's AC5. This arm is base-green **by design and says so**: it is the
  paired no-regression control proving a comment-only edit cannot move a recorded raw-grep
  observable (the iteration-92 class); its only reachable failure mode is this item's own edit.
  (c) `go build ./host/daemon/` rc=0. *Cannot see: the comment's truthfulness — that is prose;
  its behavioural content is exactly the uncovered residual §3 names, which is why the
  successor item exists.*
- **AC7 — the mutation table executed in full** (§6): every mutant vets (and the production one
  builds) before its verdict; every kill read as counted `--- FAIL:` lines from the named
  killer; every inverse arm observes exactly its stated red set with no unexplained reds; every
  sha256 moves and restores byte-identical.
- **AC8 — hygiene.** `git status --porcelain` empty after the full AC run; sha256 of both
  touched files equal to their committed post-M1 values (the one-shots of AC7 must leave no
  trace).

## 9. Conflict surface

- **Files touched**: `host/daemon/read_deadline_test.go` (test-only), `host/daemon/handlers.go`
  (comment-only). No new files, no new packages, no new exported identifiers.
- **`.ail` / `verify_ail.sh`**: **no `.ail` file is touched**; the gate pins (10 identities /
  39 named tests / 9-of-9 steps per the item brief) are unmoved. This doc contains no `.ail`
  snippet, so S4's checked-docs sweep gains nothing to sweep.
- **The frozen sketch** `design_docs/sketches/worlddapi.ail` is read (not written) by
  `sketchHTTPStatusVectors`; untouched, and the sketch-parsing arms of
  `TestTimeoutStatusMirrorsSketch` are stimulus-independent — they run identically at any
  deadline.
- **Nothing pins `handlers.go`** by name, hash, line or count outside `host/broker`'s own
  driver map (D11 — the filtered grep returned 0 with the unfiltered broker hits as the control
  that the instrument reads), and the only test in `host/daemon` that reads source text reads
  the sketch, not `.go` files (D12) — so the comment edit can move no gate. The one recorded
  raw-grep observable over this file (`err.Error()` count = 5) is pinned unmoved by AC6(b).
- **`verify_go.sh` / CI runtime**: the high-count runs of AC3/AC4 are operator acceptance
  commands, **not committed loops** — CI cost is one added ~millisecond unit test. The race
  leg's 600 s watchdog and driver drift gate are untouched (D10).
- **`tools/launchd/*`** (FLEET-owned) untouched; no skill files touched; queue/record edits are
  the controller's, not this sprint's.
- **S3, answered**: why is this not a package? It is ~35 test lines and a comment inside the
  existing `host/daemon` package, adding no surface, no file, and no kernel growth — there is
  nothing to package.
- **Sibling flake rows 20/21**: not touched and **no shared root cause is claimed** — none was
  demonstrated by command. Row 20 is a wall-clock-vs-output-cap race in `host/capsule`; row 21
  is a `CombinedOutput` stderr merge in `host/archive`; neither involves a context deadline.
- **`t.Parallel`**: the new test is serial like every other test in `host/` (the broker item's
  census); nothing here perturbs that invariant.

## 10. Sizing

The row says ~0.25 d. **It holds**: the change is one constant, two one-line site edits, one
~20-line unit test, one comment block; the cost centers are the two high-count proof runs
(minutes of wall clock, run once by the executor) and three one-shot mutations, every one of
which was either measured this session or by the controller (M4 ran the exact changed stimulus
at 500+2000; the mutations are literal reverts of measured states).

What is deliberately NOT built, so the decision is visible: a committed AST guard against future
positive sub-deadline stimuli (the shape of `host/broker`'s `t.Parallel` guard, ~60 LOC plus
its own mutation). Declined because the class has exactly two instances, both in one file,
both landed by one sprint, and AC5's two-scope sweep plus the constant's own comment covers the
current tree; the broker precedent guarded a *mutable production seam*, a standing hazard —
this guards a test idiom. If a third racy-stimulus site ever appears, that guard is the named
follow-up and this paragraph is its ≥2-instance ledger. This is a declared residual, not an
oversight.

## 11. What this item is NOT doing

- **Not** changing production behaviour: zero behavioural bytes in `handlers.go`/`daemon.go`
  (the LIMITATION is a comment; ARM B/C rejected on the record, §3).
- **Not** weakening, moving, or re-wording either test's assertions — the 503/`Timeout`
  contract pins are byte-unchanged; only the stimulus moves.
- **Not** adding a retry or a skip anywhere; the operator's once re-measurement of a one-shot
  arm (§6) binds the operator, not the committed gate.
- **Not** committing a `-count` loop to CI (§7's last sentence).
- **Not** touching rows 20/21, the store-layer context test (it already spells expiry
  correctly, D9), `blockingStore`/`recordingStore`, any `.ail` file, `verify_*.sh`, CI
  workflows, or the frozen driver.
- **Not** claiming the late-read 200 is impossible after this item — the opposite: §4.3 writes
  that limitation into the code, named, with its successor.

New residual named as a candidate queue row for the controller to file (not designed here):
the deadline does not govern lock waits — busy_timeout does (D18) — and nothing asserts
busy_timeout < readDeadline.

## 12. Verification Log

M-rows: VERIFIED BY CONTROLLER at HEAD `0b368e1`, darwin/arm64, go1.26.4,
`AILANG_BIN=/tmp/ailang-v0300/ailang`; commands as given, re-derivable, not silently widened.
D1–D14: measured first-party in this design session, 2026-08-19, at the same HEAD; D17:
measured first-party in the pre-round-2 revision pass (same HEAD, same discipline),
re-deriving the controller's iteration-94 finding that AC5's raw form was defective (§13);
repo writes this session: this doc only. D15–D16: VERIFIED BY CONTROLLER (iteration 94, same HEAD
`0b368e1`, `GOTOOLCHAIN=go1.25.6`, `AILANG_BIN=/tmp/ailang-v0300/ailang`, darwin/arm64), run
in answer to quorum round 1's two premise objections (§13); D18: VERIFIED BY CONTROLLER
(iteration 94, same HEAD and pins, `modernc.org/sqlite v1.54.0`), run in answer to quorum
round 2's one remaining objection — the second arm of the reviewer's own proposed_fix,
executed (§13); all recorded with their commands and NOT widened beyond what was measured.

| # | Claim | Command | Observed |
|---|---|---|---|
| M1 | TWO affected tests, one stimulus | `grep -n "time.Nanosecond" host/daemon/read_deadline_test.go` | `:254` (inside `TestDaemonReadDeadline`, subtest `real-store-expired-deadline`), `:526` (inside `TestTimeoutStatusMirrorsSketch`) |
| M2 | base failure rates; small samples vacuous | `go test ./host/daemon -run 'TestDaemonReadDeadline/real-store-expired-deadline' -count=500 -v`; `… -run 'TestTimeoutStatusMirrorsSketch' -count=2000 -v` | rc=1, 4× `status = 200, want 503` in 500 (~0.8%); rc=1, 13 `--- FAIL` / 14 `answered 200` in 2000 (~0.65%; one run failed two routes); failures route-agnostic (world, object, log range); a first 0/20 sample proved nothing |
| M3 | mechanism: positive 1 ns arms a timer; the 200 bodies are real | code path read + captured bodies | `readCtx` → `WithTimeout(parent, 1ns)` → `time.Until > 0` → `AfterFunc` timer; store read can complete first → `err == nil` → classifier never consulted → 200 with real data (`{"items":[]}`, full object, full world) |
| M4 | the non-positive stimulus, measured in both arms | both sites set to `-1 * time.Nanosecond`; same two commands; `go vet ./host/daemon` | 500/500 and 2000/2000 green, rc=0 both, vet rc=0; file restored from `cp` backup, sha256 byte-identical (`59c8f03df154edd1…`), `git status --porcelain` empty |
| D1 | instrument pins | `git rev-parse --short HEAD`; `go version`; `GOTOOLCHAIN=go1.25.6 go env GOROOT` | `0b368e1`; `go1.26.4 darwin/arm64`; gate toolchain resolves to `…toolchain@v0.0.1-go1.25.6.darwin-arm64` (present, no download needed) |
| D2 | M1 re-derived; site containment; file imports cover §4.2 | same grep as M1; Read of `read_deadline_test.go` in full | identical hits `:254`/`:526`; `:254` sits in the `real-store-expired-deadline` subtest, `:526` before the six-route loop of `TestTimeoutStatusMirrorsSketch`; `httptest`, `errors`, `context` already imported; `newHandlerDaemon` is the file's constructor idiom |
| D3 | class sweep at the widest scope: no third site | `grep -rn "Nanosecond\|Microsecond" --include='*.go' host/ cmd/` | exactly the 2 known hits (which are the same-call known-positive control); rc=0 |
| D4 | every `readDeadline` writer in tests; the 50 ms site is deterministic | `grep -rn "readDeadline = " host/daemon/*_test.go`; Read of the `blocking-store` subtest and `blockingStore.block` | assignments at `:254`, `:271` (`50 * time.Millisecond`), `:526`; all other hits are assertions. `:271` is safe at any positive value: the getter parks in `select { <-ctx.Done() / <-escape }` and cannot return early — deadline-DURING-read, covered deterministically |
| D5 | six production `readCtx` sites; both tests drive all six | `grep -n "readCtx(" host/daemon/*.go`; Read `seedReadRoutes` (`:49-57`) and `daemon.go:588-600` | call sites: `handlers.go:314` (world), `:347` (object), `:383` (log entry), `:438` (log range), `:472` (registry), **`daemon.go:589` (`handleHead`)** — the route table, not a file, is the enumeration, so both tests inherit full-width coverage incl. iteration 93's seventh-site file split |
| D6 | the `dur <= 0` branch cancels synchronously; positive durations arm a timer; identical across both toolchains | `grep -n "dur <= 0"` + sed of `$GOROOT/src/context/context.go` for go1.26.4 AND the go1.25.6 toolchain; `WithTimeout` at `:703` | `if dur <= 0 { c.cancel(true, DeadlineExceeded, cause); return … }` before any timer; else `c.timer = time.AfterFunc(dur, …)`; `WithTimeout(parent, t)` = `WithDeadline(parent, time.Now().Add(t))` — both toolchains, line `:645` in each |
| D7 | `database/sql` refuses an expired ctx deterministically at acquisition | sed of `func (db *DB) conn` in `$GOROOT/src/database/sql/sql.go` | opens with `select { default: case <-ctx.Done(): return nil, ctx.Err() }`; Go spec: `default` is taken only when no other case can proceed, so a closed `Done` always wins — language guarantee, paired with M4's 2,500 clean runs as the empirical arm |
| D8 | store getters wrap with `%w` → both `timedOut` arms saturated under the fixed stimulus; arm discrimination lives in `blocking-store` | sed of `store.GetWorld` (`store.go:523`); Read of `errStoreInterrupted` (`read_deadline_test.go:62-79`) | `fmt.Errorf("store: get world %q: %w", …)` — chain carries `DeadlineExceeded`; the blocking sentinel deliberately does NOT wrap it ("that is what makes MU3 … a real kill"), so §5's classifier analysis is the shipped design's own |
| D9 | the ratified design prescribed "already-expired" and claimed the arm deterministic; the sibling store test spells it correctly | sed `design_docs/implemented/w-daemon-read-cancellation.md:236-243`; grep+sed `host/store/context_read_test.go:121-122` | doc: "with an **already-expired deadline** against a real store, the correct arm answers 503 in microseconds … Deterministic … **neither can block**"; store test: `context.WithDeadline(context.Background(), time.Now().Add(-time.Hour))` — expired at construction, the §4.1 mechanism, already in-repo from the same sprint |
| D10 | gate base conditions | grep of `scripts/verify_go.sh` | AILANG_BIN v0.30.0 exact-token gate; toolchain denylist go1.26.0–1.26.5 with "e.g. GOTOOLCHAIN=go1.25.6" in its own message (`:113-121`); race leg under a 600 s watchdog; driver drift gate against `tools/launchd/` |
| D11 | nothing outside `host/broker` pins `handlers.go` (negative claim + control) | `grep -rn 'handlers\.go' host/ scripts/ .github` over `*_test.go`,`*.sh`,`*.yml`, then filtered of `host/broker` | filtered: 0 lines (rc=1); unfiltered hits exist in `host/broker` (its driver map — a different package's `handlers.go`) — the control that the instrument reads |
| D12 | the only source-reading test in `host/daemon` reads the sketch, not `.go` | `grep -rn "ReadFile\|WalkDir\|parser.ParseFile" host/daemon/*_test.go` | one hit: `read_deadline_test.go:452` — `sketchHTTPStatusVectors` reading `worlddapi.ail`; no gate can be moved by a `.go` comment |
| D13 | token and name allocation for §4.2/§4.3, with controls | `grep -c "err.Error()" host/daemon/handlers.go`; `grep -rn "LIMITATION(" host/ cmd/ scripts/`; `grep -rc "writeReadTimeout" host/daemon/handlers.go`; `grep -rn "TestExpiredReadDeadline" host/` | 5 (the item-18 AC5 live count AC6(b) pins); `LIMITATION(` absent (rc=1) with the same-call `writeReadTimeout`=7 control matching; test name unallocated (rc=1) |
| D14 | the queue row's text is what §1 corrects | Read `design_docs/world-mission.md` row 19 | row names only `TestTimeoutStatusMirrorsSketch`; SCOPE sentence: "make the stimulus deterministic rather than racing (the deadline must be observed before the read can answer)"; "Do NOT weaken the assertion"; ~0.25 d |
| D15 | **the pinned driver interrupts a read ALREADY IN FLIGHT — the wait is bounded mid-flight, by measurement** (quorum round 1, objection 1; VERIFIED BY CONTROLLER, iteration 94). Driver: `modernc.org/sqlite v1.54.0` (`go.mod:5`; imported at `host/store/store.go:32` as the pure-Go CGo-free driver) | throwaway `host/store/zz_probe_ctx_test.go` (written, run, DELETED; `git status --porcelain` verified clean afterwards): file-backed DB, recursive CTE of ~200,000,000 row-steps via `db.QueryRowContext`; CONTROL arm 20 s budget, MID-FLIGHT CANCEL arm 300 ms deadline (cancellation arrives long after the query is certainly executing); the probe FAILS LOUDLY if the "slow" query completes in under 500 ms, so a degenerate query could not have produced a false pass | control: elapsed `20.001040542s`, err = `context deadline exceeded` — the instrument check proving the query is genuinely long-running, without which the cancel arm proves nothing; cancel: elapsed `302.130667ms`, err = `context deadline exceeded` — the in-flight read was interrupted **2.13 ms after the deadline**. NOT covered, so this row cannot be over-quoted (rule 3b(ii)): one driver version, one platform (darwin/arm64), one query shape (CPU-bound); it says nothing about a read blocked on a LOCK rather than on CPU — `busy_timeout` territory, item 18's ground |
| D16 | `readDeadline` is the 10 s package constant and `readCtx` is `context.WithTimeout(r.Context(), d.readDeadline)` — the row load-bearing for BOTH §3 ARM B.2 (the shipped-deadline behaviour argument) and §6 MU-DEADLINE-DETACH's claim that the mutant compiles with both identifiers in scope (quorum round 1, objection 2; VERIFIED BY CONTROLLER, iteration 94) | `grep -n "readDeadline = 10 \* time.Second" host/daemon/daemon.go`; `sed -n '269,271p' host/daemon/handlers.go`; same-file known-positive control `grep -c "func (d \*Daemon)" host/daemon/handlers.go` | `128:	readDeadline = 10 * time.Second`; `handlers.go:269-271` = `func (d *Daemon) readCtx(r *http.Request) (context.Context, context.CancelFunc) { return context.WithTimeout(r.Context(), d.readDeadline) }`; control → 8 |
| D17 | AC5's raw grep is vacuous against §4.1 (2 at base AND at head — it cannot move); the comment-stripped form moves 2 → 1 at BOTH scopes; the strip removes comment hits without eating code hits; no other §8 AC is comment-sensitive | scratch copy of `host/daemon` with §4.1 applied verbatim (const + comment block; both stimulus lines → `d.readDeadline = expiredReadDeadline`) and the §4.3 block appended to its `handlers.go`; raw `grep -rn 'Nanosecond\|Microsecond'` and stripped `sed 's://.*::' *.go \| grep -c 'Nanosecond\|Microsecond'` at package scope; wider-scope base `find host cmd -name '*.go' -exec sed 's://.*::' {} + \| grep -c …` plus the head composition `-not -path 'host/daemon/*'`; control `sed 's://.*::' handlers.go \| grep -c readDeadline` vs raw; `grep -c 'err.Error()'` on the scratch `handlers.go`; `find host/daemon -type d`; cased-token grep of the §4.2 block | raw: base 2 (both scopes), scratch head **2** (the §4.1 comment line + the const — AC5's old "exactly 1" was unreachable); stripped: base 2 (both scopes), scratch head **1** (the const only), rest-of-`host/ cmd/` 0 → wide head 1; control stripped 6 / raw 7 — comment hit removed, code hits kept; `host/daemon` has no subdirectory, so the `*.go` glob equals the old recursive scope; §4.2's block spells neither cased token (grep 0 — its prose says "nanosecond", lowercase, outside the pattern), so its landing adds no hit; AC6(b) reads **5** on the scratch `handlers.go` carrying §4.3 (base 5 — unmoved), AC6(a) token 1 with `writeReadTimeout` control 7 |
| D18 | **the lock-blocked read path IS bounded — but by `busy_timeout`, NOT by the request deadline; the context deadline did not cut the wait short and was exceeded 6.8× while the error still surfaced as deadline-exceeded** (quorum round 2, `gpt5-6-sol`'s remaining objection — the second arm of its own proposed_fix, RUN; VERIFIED BY CONTROLLER, iteration 94, HEAD `0b368e1`, `GOTOOLCHAIN=go1.25.6`, darwin/arm64, `modernc.org/sqlite v1.54.0`) | throwaway `host/store/zz_probe_lock_test.go` (written, run, DELETED; `git status --porcelain` verified clean afterwards): file-backed DSN with `_pragma=busy_timeout(2000)&_pragma=journal_mode(delete)`; a writer holds a write transaction for 10 s (longer than BOTH bounds, so whichever fires is identifiable); a second connection then runs `QueryRowContext` under a 300 ms context deadline; CONTROL arm: an unlocked read on the same reader, same call — the probe FATALs if this control fails | CONTROL: elapsed `448.458µs`, err = nil, v = `"x"` — the reader works; the blocked arm is not measuring a broken fixture. LOCK-BLOCKED arm: elapsed `2.042929083s`, err = `context deadline exceeded` — the read TERMINATED, so the path IS bounded, but by **busy_timeout (~2 s)**, NOT by the 300 ms request deadline: the context deadline did NOT cut it short and was exceeded **6.8×** while the busy-retry loop ran, and the error still SURFACES as `context deadline exceeded`, so the wire class would be Timeout while the TIMING was governed by a different mechanism entirely. NOT covered, so this row cannot be over-quoted (rule 3b(ii)): one driver version, one platform, one pragma configuration (`busy_timeout(2000)`, `journal_mode(delete)`), one lock shape (a held write transaction); it establishes no general bound for other blocking mechanisms and makes no general production bounded-wait claim |

## Related

- Queue row 19 — `design_docs/world-mission.md` (corrected in §1: two tests, not one)
- `design_docs/implemented/w-daemon-read-cancellation.md` — item 18, whose §"already-expired"
  prescription this item restores (D9) and whose 503/`Timeout` contract these tests pin
- `design_docs/implemented/w-broker-base-flake.md` — the direct precedent: same class (base
  flake in a `host/` package), source of the counted-lines rule, the one-shot discipline, the
  attribute-by-shape protocol, and the S6 argument against sampling gates
- `design_docs/coding-standards.md` S6 (honest gates — the per-run vacuity argument of §3),
  S3 (answered in §9), S4 (no `.ail` snippet shipped because none is needed)
- `host/daemon/handlers.go:249-298` — the `readCtx`/`timedOut`/`writeReadTimeout` cluster;
  `host/daemon/daemon.go:589` — `handleHead`, the sixth read site
- Successor (named, not filed): **`w-daemon-late-read-503`** — the ARM B contract, if ever
  wanted, with the §3 B.3 fixture question answered under its own quorum

## 13. Quorum verification log

- **Round 1 — BLOCKED.** Both reviewers rejected; `absent_reviewers: []` — a full-strength
  quorum with no hole, so both objections stand as real. Both were classified **PREMISE**
  objections (each asserts a codebase claim the doc carried no evidence for); **neither
  disputes ARM A's direction**, so this revision makes the document carry the evidence and
  does not reopen the decision.
- **Both premises were measured by the controller (iteration 94) before this revision** —
  run, not forwarded:
  - **`gpt5-6-sol`**: the doc relied on an unverified — and as written, overstated — claim
    that the SQLite/context path bounds the WAIT via a mid-flight driver interrupt; D6–D8
    cover context construction, pre-acquisition refusal, and error wrapping, but nothing
    covered cancellation once a read is in flight. **Correct: no row existed.** The
    measurement is D15, and it is favourable — now the strongest evidence the doc has. What
    it changed: §2's citation, §3's opening (the residual is restated as precisely and only
    the completes-before-cancellation 200, with the unbounded-wait reading explicitly marked
    refuted), ARM A 3(c) and ARM B.2 (the bound is a measurement, not an inference), a new
    "Did D15 change this choice?" paragraph re-arguing rather than defending ARM A (verdict:
    B is now *less* necessary), and the §4.3 LIMITATION comment rewritten around the precise
    residual with D15's non-coverage stated so the row cannot be over-quoted.
  - **`gemini-3-1-pro`**: the doc asserted `readDeadline` is a 10 s package constant and
    `readCtx` is `context.WithTimeout(r.Context(), d.readDeadline)` — load-bearing for §3
    ARM B.2 and §6 MU-DEADLINE-DETACH — with zero verification. **Correct: no row existed,
    and the premise itself is true.** The measurement is D16, with a same-file known-positive
    control. What it changed: §2, §3 ARM B.2, and the MU-DEADLINE-DETACH row now cite D16 as
    the row carrying both identifiers, and the row itself names its two consumers.
- **ARM A is unchanged.** Objection 1's measurement narrows the residual ARM A accepts and
  weakens the remaining case for ARM B; objection 2's measurement confirms the code is what
  the doc said it was. No AC number moved: AC5(a) base = 2, AC5(b) base = 2 at `host/ cmd/`
  scope, AC6(b) = 5, AC6(a) = 0 with control 7, AC2's pin test absent at base — all as
  re-derived first-party by the controller before this pass.
- **Pre-round-2 controller touch-up — AC5 as revised was UNSATISFIABLE and VACUOUS at once,
  found by the CONTROLLER before round 2, not by a reviewer.** The revision's own §4.1
  comment spells `time.Nanosecond` in prose, so AC5's raw grep read **2 at base and 2 at
  head**: it could never reach its demanded "exactly 1" (unsatisfiable against this doc's own
  §4.1 — the AC would red the milestone for a non-defect), and it could never move (vacuous —
  a green would carry no information). The controller measured both arms on a scratch
  application of §4.1 before this fix; the fix is iteration 92's own precedent (AC4''s
  comment-stripped grep), itself measured in both arms — stripped 2 at base → 1 at head, at
  BOTH scopes, with a stripped-vs-raw control proving the strip removes comment hits without
  eating code hits — and re-derived first-party as D17. Recorded plainly because it is the
  worst part: this document QUOTES both governing lessons — iteration 92's "a comment can
  quote a future code state" and iteration 93's E-row "AC5's grep counts comments" — applies
  them to `handlers.go` in §4.3, and then reproduced the identical defect in its own headline
  sweep, in the AC sitting beside the comment that breaks it. Citing a class is not
  protection; only running the instrument against the proposed tree is — and the designer,
  who self-reported the risk ("adversarial reviewers may catch that"), declined to run it;
  the controller ran it. Point-5 sweep of every other §8 AC for the same class (D17): AC6(b)
  measured unmoved at 5 against a scratch `handlers.go` carrying the §4.3 comment (which
  avoids the token by construction — the doc already said so; now it is measured); AC6(a)
  greps FOR a comment, so comment-sensitivity is its function, not a defect; AC1/AC2/AC3/AC4
  count gate rc and test-output lines, AC7 counts vet/build/sha256/`--- FAIL:` lines, AC8
  reads porcelain — none greps source for a token this doc's comments spell. **None other
  changed.**
- **Round 2 — REJECTED by one reviewer; carve-out invoked.** `gemini-3-1-pro`: **PASS**
  (flipped from reject — its premise objection is satisfied by D16); `gpt5-6-sol`: **REJECT**,
  ONE remaining blocking objection; `absent_reviewers: []` — again a full-strength quorum. The
  controller classified the objection **COMPLETENESS/ATTRIBUTION, not DIRECTION**: it does not
  dispute ARM A, which is what makes the narrow-refinement carve-out available. ARM A did not
  move in this pass. The objection: the doc claimed the production deadline bounds the wait
  "regardless" and proposed an in-code statement that it "therefore bounds the WAIT", while
  D15 verifies only one CPU-bound SQLite query on one driver version and platform and
  explicitly excludes lock-blocked reads — the doc generalized beyond its own evidence on the
  mission-critical bounded-waits axiom. **The objection was CORRECT.** The reviewer's
  proposed_fix was applied VERBATIM: its limitation sentence ("D15 verifies prompt
  cancellation only for one CPU-bound query using modernc.org/sqlite v1.54.0 on darwin/arm64.
  It does not establish a general bound for lock-blocked or other driver waits; therefore this
  item makes no general production bounded-wait claim.") now stands word-for-word in §3, in
  §4.3's prose, and in the in-code comment; the over-claims the catch names — "the wait itself
  stays bounded regardless", "the deadline therefore bounds the WAIT", and the "precisely and
  only" residual framing — are deleted from §3 and §4.3, and the residual is restated as TWO
  things (§3, §4.3). The fix's second arm — verification rows exercising lock
  contention/busy-timeout behaviour under the pinned gate toolchain and driver — was **RUN by
  the controller** (iteration 94), recorded as **D18**, and its result **STRENGTHENED the
  objection rather than dissolving it**: the lock-blocked read terminated at 2.042929083 s
  under a 300 ms context deadline — bounded, but by busy_timeout (~2 s), not by the deadline,
  which was exceeded 6.8× while the busy-retry loop ran, the error still surfacing as
  `context deadline exceeded`. The controller's finding, stated as found: the fix's "expand
  the design or block" trigger does NOT fire — every path measured terminates within a stated
  bound (2.04 s, by busy_timeout) — but the doc's prior claim that the DEADLINE bounds the
  wait is FALSE on the lock-blocked path; today the composition is still safe only because
  busy_timeout (2 s) < readDeadline (10 s), an ordering NOTHING in the code asserts — an
  ordering, not a guarantee. That residual is named at the end of §11 as a candidate queue
  row for the controller to file; it is not designed here. The arm choice was re-tested
  against D18 (§3, "Did D18 change this choice?"): ARM B's post-read check runs only on a
  SUCCESSFUL read and D18's lock-blocked read surfaces an error, so D18 does not make ARM B
  more attractive; ARM A stands.
- **Round-2 token re-check of the revised §4.3 comment, with numbers** (the mandate: keep the
  comment free of `err.Error()`, `Nanosecond`, `Microsecond`, and of `busy_timeout` if any
  recorded grep counts it; re-checked against AC5 and AC6(b)): the revised comment spells
  `Nanosecond` **0** times and `Microsecond` **0** times (its durations are `300ms`, `2.04s`,
  `~2s`, `2s`, `10s`), so AC5(a) still reads head **1** and AC5(b) head **1**, both bases
  unchanged at **2**; it contains **0** occurrences of literal `err.Error()`, so AC6(b) stays
  **5**; it spells `busy_timeout` **3** times — checked: this doc's recorded greps count
  `Nanosecond\|Microsecond` (AC5), `LIMITATION(w-daemon-late-read-503)` / `writeReadTimeout`
  (AC6(a)) and `err.Error()` (AC6(b)), none of which matches it, and the only recorded greps
  in the repo counting `busy_timeout` (item 18's V6 and its sprint plan's AC6) are scoped to
  `host/store/*.go`, not `host/daemon/` — measured this pass:
  `grep -rn 'busy_timeout' host/store/*.go | grep -v _test` is untouched by a
  `host/daemon/handlers.go` comment by scope. The comment also avoids the bare token
  `readDeadline`, so AC5's stripped-vs-raw instrument control on `handlers.go` (stripped 6 /
  raw 7, D17) is unmoved. No AC number moved in this revision: AC5(a) base 2 → head 1, AC5(b)
  base 2 → head 1, AC6(a) = 1 with control 7, AC6(b) = 5, AC2 absent at base — all as before.
- **Observation, left and not acted on (out of this pass's mandate):** the header's
  `handlers.go` estimate "~+18/−0" now under-counts the §4.3 comment (~+27 lines after this
  revision); the edit remains comment-only, zero behavioural bytes, and D11/D13's "no gate
  moved" analysis is unchanged.
