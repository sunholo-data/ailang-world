# Sprint plan — `w-broker-base-flake` M1 (queue item 16, clause-3)

**Item**: queue item 16 — *localise the `host/broker` 0.76% base flake, de-align the
exec-assessment race, and prove the gate without sampling it*.
**Status**: PLANNED · **NO SPLIT** · one milestone, **5 tasks / 5 commits**
**Scope**: **M1 ONLY.** M2 (the post-reap re-sweep) is decision-gated by the design doc's §6 and
is NOT planned, NOT implemented, and NOT authorized by anything in this sprint.
**Design doc**: [`design_docs/planned/w-broker-base-flake.md`](w-broker-base-flake.md)
(quorum-cleared over 3 rounds; resolution A human-ratified 2026-08-13T06:12:23Z, `option A`)
**Base**: `dev` @ `002b13f`, clean tree (`git status --porcelain` → 0 lines, P0)
**Planner**: mission-control iteration 80, opus sprint-planner, first-party measurement on this rig
**THE CONTROLLER MAKES ALL COMMITS.** The executor has **no git write permission** — no
`git add`/`commit`/`push`/`checkout`/`stash`/`restore`/`worktree`, no `gh pr`. **Every restore in
this sprint is `cp` from `/tmp/w16_backup/`**, never `git checkout -- <file>`: the mutated file is
uncommitted by construction, so `git checkout --` would delete the executor's work and the sha256
assertion would then *report* the disaster instead of preventing it.

**Headline price: ≈0.7 day, i.e. 1.4× the doc's ~0.5 day.** The driver is not new scope. It is
that **the design doc's §5.4 recorder code block, applied verbatim, makes AC3 — the doc's only
forced-failure proof, the arm that replaces the queue row's `-count=20` coin flip — structurally
incapable of firing.** That is measured, not argued (§0.2 C1), and it is a one-line fix that the
executor must be told before it writes the file rather than after AC3 comes back green.

**The whole mechanism was PROTOTYPED, RUN, and MUTATED end-to-end this session** on the real tree
under `cp` backups, restored byte-identical (P1–P12, §7). The seam, the delegating recorder, the
warm-up, the marker fixture, the AST guard, and four of the nine mutation arms were all executed.
Nothing in T1–T4 is exploratory.

---

## 0. Planner's first-party verification

Everything below was re-measured at `002b13f` before any planning. Command/observed table in §7.

### 0.1 Design-doc and controller premises CONFIRMED

| Premise | Verdict |
|---|---|
| `handlers.go:100` is the sole group-kill site; `cmd.Cancel` is the seam point | **CONFIRMED** (P2) |
| the seam applies cleanly, +sha `8419874…` → `666a5d4…` — **byte-identical to N5's recorded hashes** | **CONFIRMED** (P4) — the doc's N5 row reproduces exactly |
| `TestHandlerTimeoutKillsTheWholeProcessGroup` at `handlers_test.go:729`, `_, _, _ = session.Invoke` at `:740`, 2s assertion at `:743` | **CONFIRMED, exact** (P2) |
| base flake: 20/20 `--- PASS` isolated `-race` runs today | **CONFIRMED** (P3) — reproduces W4 |
| a `-run` naming no test prints `PASS` at rc=0 (V-G/W4) | **CONFIRMED first-party** (P3, P5) — the discipline holds on every AC below |
| `handlerSession` injects any `Handler` and returns the `*handlerRecordingStore` | **CONFIRMED** (P2) |
| 0 `t.Parallel` across 46 `_test.go` under `host/`+`cmd/`, grep rc=1; 0 `.Parallel` *selectors* in `host/broker` | **CONFIRMED** (P6) — R4 reproduced, and widened to the guard's actual grammar |
| no `wantFileCount`-class file pin on `host/broker` (only `host/boundary/allowlist_world_test.go:1163-1165`) | **CONFIRMED** (P6) — R5 reproduced, which is what makes two new `.go` files here safe |
| the two `host/broker` AST boundary gates are green at base | **CONFIRMED** (P5): exactly 2 `--- PASS:` |
| `warm-up` guard works and is cheap; warm exec moves the fork ahead of the deadline | **CONFIRMED** (P8): warm-up exec 215–244 ms, then `forked` present on every timed run with elapsed ~104–107 ms |
| MUT-SEAM-BYPASS is lethal | **CONFIRMED** (P11): elapsed stays green at 103 ms, kill-record `count=0` reds |
| MUT-BOUND-LOOSE composes as specified | **CONFIRMED** (P10): with the bound at 10 s the neutered-kill one-shot goes **green** at 5.16 s |
| MUT-GUARD-PARALLEL is lethal and vet stays clean | **CONFIRMED** (P12): guard reds naming `handlers_test.go:730:4`, `go vet` rc=0 |
| both verify gates green at base | **CONFIRMED** (P1): `verify_ail.sh` rc=0 (`4 required identities, 14 named tests`), `verify_go.sh` rc=0 in 4:57 |

### 0.2 Seven corrections to the design doc — each MEASURED, each changing the plan

#### C1 — **BLOCKING. §5.4's recorder code block makes AC3/P1 VACUOUS.** The load-bearing one.

The doc's §5.4 recorder is written as:

```go
rec.errno = syscall.Kill(-pgid, syscall.SIGKILL) // the real kill
```

i.e. the test-side wrapper **re-implements** the kill. But `killGroup` is a package-level `var`,
and the test **replaces** it. So the production seam body is dead code for the whole duration of
the test — and MUT-KILL-NEUTER / P1 / AC3 mutates *exactly that dead body*
(`§7 P1`: "change `syscall.Kill(-pgid, …)` to `syscall.Kill(pgid, …)` **inside the §5.4 `killGroup`
seam**"; `§8` row 1: "in `handlers.go`").

Measured, one binary, one mutant, two recorder shapes (P9):

| recorder shape | MUT-KILL-NEUTER applied to `handlers.go` | verdict |
|---|---|---|
| doc §5.4 verbatim (re-implements the kill) | `--- PASS` at **105 ms**, `survived=false` | **mutation invisible — AC3 cannot fire** |
| delegating (`orig := killGroup` captured, `rec.errno = orig(pgid)`) | `--- FAIL` at **5.168 s**, `survived=true`, `count=1`, `errno=nil`, `offset=103 ms` | **lethal, and reproduces §6's full H1 signature** |

The fix is one line and is **binding on T4**: capture `orig := killGroup` *before* the swap and
delegate to it. Consequences beyond AC3, all good: the delegating shape is what makes §7 P1's
extended claim true — the forced child-only kill reproduced the complete decision-condition
signature (`*HandlerTimeoutError` reachable via `errors.As`, kill count 1, errno nil, offset
≈103 ms, `survived` present, elapsed 5.168 s inside §2's 5.0–5.5 s natural-death band). With the
doc's shape, P1's "proves §6's decision condition is reachable and truthfully reported" would have
been asserted about a run that never left the healthy path.

Note the shape of this defect for the record: the doc *says* "a recording wrapper that still
performs the real kill" (§5.4, prose) — which is the delegating design — and then writes a code
block that duplicates the kill instead of wrapping it. The prose is right; the code is wrong; the
code is what an executor types.

#### C2 — §5.1's pass-path assertion, taken literally, reds every healthy run.

§5.1 says the test "explicitly asserts that `err` is a `*HandlerTimeoutError`". It is not.
`Session.invoke` returns `&EffectFailedError{… cause: err}` (`broker.go:284`), whose `Unwrap`
returns the cause (`broker.go:140-142`). Measured on a healthy run (P8):

```
err=*broker.EffectFailedError   errors.As(&*HandlerTimeoutError)=true   err.(*HandlerTimeoutError)=false
```

So the assertion must be `errors.As`, and the plan writes it that way (T4). An executor that typed
the doc's sentence as a type assertion would have shipped a test that fails 100% of the time —
loudly, so it would have been caught, but it would have burned an execution round.

#### C3 — **AC2 cannot fail for the reason it claims. The §5.4 mutex is not `-race`-provable.**

AC2 says it "**is the deliberate other half of N5**… it can fail (an unsynchronised recorder reds
under `-race`)", discharging `gemini-3-1-pro`'s round-2 objection. Measured (P7):

- an **unsynchronised** recorder (no mutex, fields read immediately after `Invoke` returns), built
  `-race`, **20 isolated runs: 0 `DATA RACE`, 20/20 `--- PASS`**;
- **known-positive control in the same binary**: a deliberate two-goroutine increment →
  `WARNING: DATA RACE`, `--- FAIL`. The detector is live.

The source says why, and it is not luck: `os/exec.(*Cmd).watchCtx` calls `c.Cancel()` and *then*
sends on `resultc` (`$GOROOT/src/os/exec/exec.go:805-820`), while `Wait` **receives** from that
channel (`:930`). That send/receive is a happens-before edge, so every write the `Cancel` hook
performed is visible to the goroutine returning from `Wait` — hence to the test, which reads only
after `Invoke` (and therefore `Wait`) returned. There is no reachable race in the shape under test.

**What the plan does with this.** It **keeps the mutex** — it is ratified quorum material, it costs
two lock operations, and it is genuine defence-in-depth against a future refactor that reads the
record while `Invoke` is still running (e.g. any `WaitDelay`, or a second kill site). It **does
not** let AC2 claim to prove it. AC2's honest content is written into the criteria as exactly what
it can decide: 20 counted `--- PASS:` lines from the recorder-swapped test under `-race`, with the
per-run marker and kill-record assertions observed in the log. Leaving AC2's original wording in
place would have put a gate in this sprint whose stated null case cannot occur — the S6 defect
this document exists to remove, for the *third* time in its own history (the row's `-count=20`,
round 3's cold-majority arm, and now this).

#### C4 — the pass path must NOT assert `survived` absent, and the doc never says so.

A strictly stronger assertion is available and free: on a healthy run `survived` is never written
(it needs the full 5 s), so `survived` absent could be asserted unconditionally. **It must not be**,
because MUT-BOUND-LOOSE (§8) requires the composed `MUT-KILL-NEUTER + 2s→10s` one-shot to go
**GREEN** — that green *is* the proof the 2 s bound is load-bearing. Measured (P10): with the bound
at 10 s and the kill neutered, the delegating test passes at 5.16 s **with `survived=true`**. Add a
`survived`-absent assertion and that arm reds instead, and AC3's specified verdict becomes
unreachable. The trade — a weaker per-run assertion bought for a lethal bound-mutation — is real
and is nowhere stated in the doc. The plan takes the doc's side (no `survived` assertion on the
pass path; `survived` is printed in the diagnosis) and records the trade here.

#### C5 — §5.3's perturbation claim is false for one of the three markers.

"the marker writes sit **after** `sleep 5 &`, so the fork's timing relative to `cmd.Start` is
byte-for-byte unmoved" — but `exec_started` is written **before** `sleep 5 &` in the doc's own
fixture block. The fork is therefore delayed by one file creation. Immaterial in practice (P8: with
the warm-up, `forked` is present on every run and elapsed is 104–107 ms against a 100 ms deadline,
so the fork lands with the margin §5.5 claims), but the stated discipline is inaccurate and a
reader relying on "byte-for-byte unmoved" would be relying on something untrue.

#### C6 — the seam is **+7/−1**, not +5/−1.

Measured `git diff --stat` after applying the seam whose sha matches N5's recorded `666a5d4…`
exactly: `1 file changed, 7 insertions(+), 1 deletion(-)` (P4). The doc states "+5/−1" in five
places. Since the sha is identical, N5 measured the same text and mislabelled the count. Cosmetic,
but AC6 pins a hash and a reviewer counting lines against the doc would find a discrepancy.

#### C7 — §5.6's CI budget baseline is stale by 2.5×.

§5.6 prices the new arms "against the race leg's 600 s watchdog with `host/broker` at 76.9 s
critical path (W9)". That 76.9 s is a **comment inside `scripts/verify_go.sh:110-111`** dated to
`7550ee9`, not a fresh measurement. Measured today at `002b13f` (P1): `host/broker` = **90.2 s**
in the plain leg and **193.4 s** in the race leg; the 600 s watchdog wraps the whole race leg
(`verify_go.sh:113-126`). The conclusion survives — ~3 s of committed pure-Go sleeps plus ~0.25 s
of warm-up per leg against a 600 s cap on a 193 s critical path is comfortable — but the margin is
≈3×, not ≈8×, and the plan states the measured number so the next budget decision starts from a
true one.

#### C8 — §5.2's store seam is underspecified: `handlerSession` hardcodes the store.

`handlerSession` constructs `&handlerRecordingStore{base: openTestStore(t)}` internally, so there is
no "timing variant" to swap in without either bypassing `handlerSession` (calling `newSession`
directly, duplicating its grant/registry setup) or extending the recording store. `objectStore` *is*
an interface (`broker.go:37-43`), so both are possible. **Plan decision** (T3): add one nil-default
hook field `onStoreCall func(op string)` to `handlerRecordingStore` and fire it at the top of
`PutObject`, `AppendNextEffectIntent`, `AppendEffectOutcome` — the three calls §5.2 names. The
diagnosis helper sets it on the store `handlerSession` already returns. Cost: 4 lines in
`handlers_test.go`, no new constructor, no duplicated session wiring, and every existing caller of
`handlerSession` is unaffected because the hook is nil.

Related: `repositoryRoot(t)` already exists in this package (`invoke_boundary_test.go:198`) and
returns the **repo** root. The guard needs the **package** directory, so it resolves its own
`runtime.Caller(0)` + `filepath.Dir` and must not redeclare that name.

#### C9 — the doc never says where the warm-up helper and its unit test live.

§5.5/AC5(a) require a test-owned warm-up helper with an injectable runner plus
`TestWarmUpRunsExactlyOnceUnderABoundedContext`, but assign them to no file. **Plan decision**:
both go in `handlers_stall_diag_test.go`, holding the new-file count at **two** — which is what the
conflict-surface analysis (§10, W11/R5) actually measured. A third new file would still be safe
(P6 re-confirms no file-count pin on `host/broker`) but it would be outside what the doc checked.

#### C10 — observation, not a defect: `W16_WARM` is a probe session's name.

The committed fixture's guard variable is `W16_WARM` — "W16" is the design session's probe label
(`/tmp/w16probe*`), meaningless in the landed tree. The plan keeps the doc's literal so the landed
code and the doc agree, and records the observation for whoever renames it later. If the executor
prefers a self-describing name it must change the doc's §5.3 block in the same commit.

---

## 1. Milestone decision — NO SPLIT

One milestone, `M1`, five tasks, five commits. The doc draws exactly one internal boundary worth
splitting on (§12's fallback: land §5.1–5.4+5.6, demote §5.5) and that boundary is a **fallback
under time pressure, not a plan** — the doc says so, and says the demotion also costs AC5(a). It
is not taken: §5.5's warm-up is what makes AC3 deterministic (cold, the same mutation leaked only
2/8, W7), so demoting it would put the sprint's one forced-failure proof back on a coin flip.

The five tasks are sequenced by dependency, not by size:

| # | Task | Files | LOC | Why here |
|---|---|---|---|---|
| T1 | the `killGroup` seam | `host/broker/handlers.go` | +7/−1 | nothing else can be written until the seam exists; measured gate-clean (P4/P5) |
| T2 | the no-`t.Parallel` AST guard | `handlers_parallel_guard_test.go` (NEW) | ~70 | must land **no later than** the first test that swaps the package global, since it is what makes the `t.Cleanup` restore safe |
| T3 | diagnosis helper, 2 attribution arms, warm-up helper + unit test, store hook | `handlers_stall_diag_test.go` (NEW), `handlers_test.go` (+4) | ~190 | T4 calls all of it |
| T4 | the modified group-kill test | `handlers_test.go` | ~+85 | consumes T1–T3 |
| T5 | mutation campaign + full AC sweep | none (one-shots, all restored) | 0 | verdicts only mean something once T1–T4 are all in |

Total ≈ **350 LOC, of which ≈343 are test-only**; the entire production delta is T1's +7/−1.

**Sizing: ≈0.7 day** against the doc's ~0.5. The delta is not scope, it is: the C1 correction (the
doc's stated 0.5 d assumed the recorder was "the probe's two-assignment pattern, N4" — it is not,
it is a delegating wrapper whose design the doc got wrong); nine mutation rows, five of them
one-shots with sha-verified restores; and AC1 running `verify_go.sh` twice at **4:57 measured
each** (P1). Deliberately NOT in the estimate: M2, and any sampling of the 0.76%.

---

## 2. Task T1 — the `killGroup` seam (`handlers.go`, +7/−1)

Exact text, prototyped this session and sha-verified against N5 (P4). Insert immediately above the
`// runBounded is the one timeout and allocation surface` comment:

```go
// killGroup is the cancellation kill boundary, a package-level seam so the
// timeout tests can observe the kill's time, target and errno directly.
var killGroup = func(pgid int) error {
	return syscall.Kill(-pgid, syscall.SIGKILL)
}
```

and change the one line inside `cmd.Cancel` (`handlers.go:100`):

```go
-		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
+		return killGroup(cmd.Process.Pid)
```

The nil-`Process` guard **stays in `Cancel`**. `handlers.go:116`'s `cmd.Process.Kill()` (the
overflow path's direct-child kill) is **out of scope and untouched** (N1/P2).

**Post-T1 verification (all measured on the seamed tree this session, P4/P5):**
- `shasum -a 256 host/broker/handlers.go` → `666a5d420ce2bd11bb69ce93aada9f27838de919e39bd453160b67c7e4642d70`
  (base was `841987499da53f115584c477d9306193fa87130e484d1255d853dae6e10b191c`). **This is the
  post-M1 committed hash AC6 pins.**
- `go build ./host/broker/` rc=0, `go vet ./host/broker/` rc=0.
- `go test ./host/broker/ -run 'TestRegistryDispatchBindingBoundary|TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential' -v -count=1`
  → exactly **2** `--- PASS:` lines.
- `go test ./host/broker/ -race -run '^TestHandlerTimeoutKillsTheWholeProcessGroup$' -count=5 -v`
  → 5 `--- PASS:` (the unmodified test still passes over the seam).

## 3. Task T2 — `handlers_parallel_guard_test.go` (NEW, ~70 LOC)

One test, `TestBrokerTestsDoNotCallParallel`. **Prototyped and mutated this session** (P6, P12);
the prototype's measured behaviour is the specification:

- resolve the **package** directory from this file's own `runtime.Caller(0)` + `filepath.Dir`;
  `t.Fatal` if `runtime.Caller` fails or the directory cannot be read (**do not** reuse
  `repositoryRoot`, C8);
- `os.ReadDir` the package directory, take every entry whose name ends `_test.go` (no globbing);
- `parser.ParseFile` each, `t.Fatal` on any parse error naming the file;
- one `ast.Inspect` pass per file counting **two** selector names: `Parallel` (offenders, recorded
  as `file:line:col`) and `Helper` (the anti-vacuity known-positive);
- **anti-vacuity floor, all three `t.Fatal`** — zero enumerated `_test.go` files; the enumeration
  missing either anchor (`handlers_test.go` and the guard's own file, via `filepath.Base`); zero
  `Helper` selectors across the set;
- then `t.Fatalf` on any offender, naming every `file:line:col` and saying *why* (it would race the
  `t.Cleanup` restore of the `killGroup` package global).

**Measured base and mutant (P6, P12):**
- pristine: `--- PASS`, `enumerated=15 files, Helper selectors=62, offenders=[]` (13 committed
  `_test.go` + the two probe files present at measurement time; post-M1 the count is 15 again:
  13 + the two new files);
- MUT-GUARD-PARALLEL (`t.Parallel()` as the first statement of the group-kill test): `--- FAIL`,
  `offenders=[handlers_test.go:730:4]`, `go vet ./host/broker/` still rc=0 — so the red is the
  guard's own, not the compiler's;
- **also lethal on the method-value form** `var tt = t.Parallel` → `offenders=[handlers_test.go:730:13]`.
  This is the grammar-coverage claim §5.4 makes about why an AST walk beats a grep, executed rather
  than asserted.

Declared residual, unchanged from the doc: a reflection-spelled `MethodByName("Parallel")` is out
of grammar, and any unrelated selector named `Parallel` (e.g. `config.Parallel()`) is a false
positive. Measured false positives today: **0** (P6, `grep -rn '\.Parallel' host/broker/` rc=1 with
a same-call `.Helper` control returning 9 files).

## 4. Task T3 — `handlers_stall_diag_test.go` (NEW, ~190 LOC) + 4 lines in `handlers_test.go`

### 4.1 The store hook (in `handlers_test.go`, C8)

Add `onStoreCall func(op string)` to `handlerRecordingStore` and, at the top of `PutObject`,
`AppendNextEffectIntent` and `AppendEffectOutcome`:

```go
if s.onStoreCall != nil {
	s.onStoreCall("PutObject") // etc.
}
```

Nil for every existing caller, so no landed test changes behaviour.

### 4.2 `invokeWithStallDiagnosis` — the phase decomposition (§5.2)

A helper taking `(t, session, store, req, payload)` and returning a struct carrying: `err`,
total `elapsed`, the Execute enter/exit offsets (from a `Handler` timing shim wrapping the real
handler — `handlerSession` already accepts any `Handler`, P2), the per-store-call offsets (from
§4.1's hook), and a `String()`/`Format` method rendering **elapsed · pre-handler window · Execute
window · post-handler window · `%+v` of the error**. Every boundary costs one `time.Now()`.

### 4.3 The warm-up helper with an injectable runner (§5.5, AC5(a))

```go
type warmUpRunner func(ctx context.Context, path string) error

type warmUpCall struct {
	count        int
	hadDeadline  bool
	timeout      time.Duration
}

func warmUpFixture(t *testing.T, path string, run warmUpRunner) { /* bounded ctx, t.Fatal on error */ }
```

Production runner: `exec.CommandContext(ctx, path)` with `cmd.Env = []string{"W16_WARM=1"}` under
`context.WithTimeout(context.Background(), 30*time.Second)`; a non-zero exit or a timeout is
`t.Fatalf` (never a silent skip, §5.5). Measured cost (P8): **215–244 ms**, paid once.

`TestWarmUpRunsExactlyOnceUnderABoundedContext` injects a recording runner that performs no exec
and asserts, on **exact values, no timing**:
- `count == 1`;
- `ctx.Deadline()` returned `ok == true` (the bounded-context flag) — kills MUT-WARM-UNBOUNDED;
- the deadline is in the future / the timeout is non-zero.

### 4.4 The two committed attribution arms (§5.6)

- `TestStallDiagnosisAttributesHandlerWindow` — a pure-Go `HandlerFunc` that sleeps 1.5 s then
  returns `fmt.Errorf("stub: %w", ErrHandlerTimeout)`; requires ≥1.2 s attributed to the Execute
  window **and** the error text present in the rendered diagnosis.
- `TestStallDiagnosisAttributesStoreWindow` — an `onStoreCall` hook that sleeps 1.5 s on
  `AppendEffectOutcome`; requires ≥1.2 s in the post-handler window.

Neither arm reaches `runBounded`, so neither asserts on the kill record (§5.6). CI cost ≈3 s per
leg, against the measured 193.4 s `host/broker` race-leg critical path and a 600 s watchdog (C7).

## 5. Task T4 — the modified group-kill test (`handlers_test.go`, ~+85)

Order matters; this is the sequence the prototype ran (P8–P11).

1. `markdir := t.TempDir()`.
2. Fixture: **reuse `writeExecutable` unmodified** (five other tests share it, §13) — pass the
   body only, since `writeExecutable` already supplies `#!/bin/sh\nset -eu\n`:
   ```
   if [ "${W16_WARM:-}" = "1" ]; then exit 0; fi
   : > "<markdir>/exec_started"
   sleep 5 && : > "<markdir>/survived" &
   : > "<markdir>/forked"
   wait
   ```
   `<markdir>` is baked into the source as an absolute path (the handler's env is a fixed
   allowlist, so no variable can reach the script, W3).
3. `warmUpFixture(t, fake, execWarmUpRunner)` — §4.3.
4. `NewGitHandler(GitHandlerConfig{GitPath: fake, ExecTimeout: 100 * time.Millisecond})`, unchanged.
5. `session, rec := handlerSession(t, EffectGitCommit, scope, timingShim{handler})`.
6. **The kill recorder — the C1 shape, binding:**
   ```go
   rec := &killRecord{}          // mu, count, offset, pgid, errno
   orig := killGroup             // CAPTURE FIRST — this is what makes AC3 lethal
   t.Cleanup(func() { killGroup = orig })
   var invokeStart time.Time
   killGroup = func(pgid int) error {
   	rec.mu.Lock()
   	defer rec.mu.Unlock()       // held across the WHOLE write
   	rec.count++
   	rec.offset = time.Since(invokeStart)
   	rec.pgid = pgid
   	rec.errno = orig(pgid)      // NOT syscall.Kill(...) — see §0.2 C1
   	return rec.errno
   }
   ```
7. `invokeStart = time.Now()`; the timed `Invoke`, **error captured** (`_, _, invokeErr :=`).
8. Snapshot under the same lock into **field copies**, never `snap := *rec` (vet copylocks, §5.4).
9. Assertions — every one of these is a per-run postcondition, none is a threshold over a
   population:
   - `elapsed <= 2*time.Second` (bound **unchanged**, §13);
   - `errors.As(invokeErr, &timeout)` with `timeout *HandlerTimeoutError` — **`errors.As`, not a
     type assertion** (§0.2 C2);
   - `count == 1 && errno == nil && offset >= 100*time.Millisecond && pgid > 0`;
   - markers `exec_started` **and** `forked` both present (the non-vacuity signal, §5.3/§7 P3(b));
   - **no assertion on `survived`** (§0.2 C4) — it is printed, not asserted.
10. On any failure: print the full diagnosis — elapsed, the three-phase split, `%+v` of the error,
    all three marker bits, and the kill record (count / offset / pgid / errno). This is what §6's
    decision condition reads.
11. Update the test's leading comment: the landed text says *"Only elapsed is asserted. Checking
    that the grandchild is DEAD afterwards is vacuous"* — true of the old test, false of this one,
    and leaving it would be a doc-level lie inside the file the next reader opens.

**Measured on the prototype of exactly this shape (P8):** elapsed 104–107 ms, kill `count=1
offset=102–103 ms errno=nil`, markers `exec_started=true forked=true survived=false`,
`errors.As` true. 20/20 `--- PASS` under `-race` (P7, on the no-mutex variant; the mutexed variant
passed every run it was given).

## 6. Task T5 — mutation campaign and AC sweep

Protocol, non-negotiable, applied to **every** row: `cp` backup into `/tmp/w16_backup/` first ·
apply the mutant · **assert `shasum -a 256` MOVED** · **assert `go vet ./host/broker/` rc=0 on the
mutated tree** (`go build ./...` does not compile `_test.go` at all, so it cannot gate a test-side
mutant) · run the **`-run`-scoped** command and read its **counted `--- PASS:`/`--- FAIL:` lines,
never rc** · restore by `cp` · **assert the sha returned to its pre-mutation value** · re-run to
green. Prefer neutering (`if false && <cond>`) over deleting a block so every import stays used.

| # | Mutation | Where | `-run` scope | Required | Status |
|---|---|---|---|---|---|
| MUT-KILL-NEUTER | `-pgid` → `pgid` inside `killGroup` | `handlers.go` (production) | `^TestHandlerTimeoutKillsTheWholeProcessGroup$` `-count=5` | **5/5 `--- FAIL`**, each diagnosis showing the full H1 signature (count 1, errno nil, offset ≈100 ms, `survived` present, elapsed 5.0–5.5 s) | **PROTOTYPED P9**: 1/1 FAIL at 5.168 s with the exact signature; mutant sha `3dc74e2…`, vet rc=0 |
| MUT-BOUND-LOOSE | 2 s → 10 s, **composed on top of MUT-KILL-NEUTER** | `handlers_test.go` | same | the composed pair goes **GREEN** — that green is the kill signal | **PROTOTYPED P10**: `--- PASS` at 5.16 s, `survived=true` |
| MUT-SEAM-BYPASS | revert `Cancel` to the inline `syscall.Kill(-cmd.Process.Pid, …)`, orphaning the seam | `handlers.go` | same | `--- FAIL` on **kill-record count 0** while elapsed stays green | **PROTOTYPED P11**: FAIL at 103.6 ms, `count=0`, sha `45abc7c…`, vet rc=0 |
| MUT-GUARD-PARALLEL | `t.Parallel()` first statement of the group-kill test | `handlers_test.go` | `^TestBrokerTestsDoNotCallParallel$` | exactly 1 `--- FAIL:` naming `handlers_test.go:<line>` | **PROTOTYPED P12**: FAIL, `handlers_test.go:730:4`; method-value form also caught |
| MUT-ERR-DISCARD | restore `_, _, _ =` at the `Invoke` call | `handlers_test.go` | `TestStallDiagnosis` | handler-window arm reds: diagnosis lacks the error text | to run |
| MUT-DIAG-BLIND-EXEC | neuter the Execute enter/exit timestamps | `handlers_stall_diag_test.go` | `^TestStallDiagnosisAttributesHandlerWindow$` | exactly 1 `--- FAIL:` | to run |
| MUT-DIAG-BLIND-STORE | neuter the store-call timestamps | `handlers_stall_diag_test.go` | `^TestStallDiagnosisAttributesStoreWindow$` | exactly 1 `--- FAIL:` | to run |
| MUT-MARKER-DROP | fixture stops writing `forked`/`exec_started` | `handlers_test.go` | `^TestHandlerTimeoutKillsTheWholeProcessGroup$` | `--- FAIL` on the non-vacuity assertion | to run |
| MUT-WARM-DROP | `if false && …` on the warm-up runner call | `handlers_stall_diag_test.go` | `^TestWarmUpRunsExactlyOnceUnderABoundedContext$` | exactly 1 `--- FAIL:`, recorded count **0 ≠ 1** | to run |
| MUT-WARM-UNBOUNDED | helper passes `context.Background()` | `handlers_stall_diag_test.go` | same | exactly 1 `--- FAIL:`, deadline flag **unset** | to run |

Ten rows (the doc's nine plus MUT-BOUND-LOOSE counted separately from its partner). **No row's
verdict is a rate, a majority, or a threshold over a population.**

---

## 7. Acceptance criteria, each with its BASE-STATE result

All from repo root. **`export AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.25.6` before any
Go command** — `verify_go.sh` FATALs without them and `host/broker` is red 100% without
`AILANG_BIN` (a base condition, not a regression). **Every `-run` command's verdict is its counted
`--- PASS:` / `--- FAIL:` lines, never rc** — measured first-party this session: a `-run` naming no
test prints `testing: warning: no tests to run` + `PASS` at **rc=0** (P3, P5).

| AC | Command | **BASE at `002b13f`** | Required after M1 |
|---|---|---|---|
| **AC1** | `./scripts/verify_ail.sh`; `./scripts/verify_go.sh` | **rc=0 / rc=0** (P1). ail: `4 required identities verified, 14 named tests pass`, `9/9 steps`. go: 4:57 wall, `host/broker` 90.2 s plain + 193.4 s race | both rc=0 again |
| **AC2** | `go test -c -race -o /tmp/w16-broker-race.test ./host/broker`; 20 × `/tmp/w16-broker-race.test -test.run '^TestHandlerTimeoutKillsTheWholeProcessGroup$' -test.count=1 -test.v`, counting `--- PASS:` | **20 pass / 0 fail / 0 neither** on today's *unmodified* test (P3), reproducing W4; control `-test.run '^TestNoSuchTestName$'` → 0 `--- PASS:` lines at rc=0 | 20 × `--- PASS:` on the **modified** test, and every run's log showing `exec_started`+`forked` and the kill record (count 1, errno nil, offset ≥ deadline). **Verdict is the counted lines and the logged assertions — NOT a proof of the mutex (§0.2 C3)** |
| **AC3** | the MUT-KILL-NEUTER protocol of §6, including the composed MUT-BOUND-LOOSE pair and byte-identical restore | n/a (production seam does not exist at base). Prototyped: mutant reds 5.168 s with the full H1 signature; composed pair greens at 5.16 s; restore to `841987…` verified (P9–P11) | 5/5 `--- FAIL` then 5/5 `--- PASS` after restore; composed pair GREEN |
| **AC4** | `go test ./host/broker/ -run 'TestStallDiagnosis' -v -count=1` | **0 `--- PASS:` lines**, rc=0, `[no tests to run]` — the file does not exist (P5) | exactly **2** `--- PASS:` |
| | `go test ./host/broker/ -run 'TestRegistryDispatchBindingBoundary|TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential' -v -count=1` | **exactly 2 `--- PASS:`** (P5: 0.27 s, 1.46 s) — and **2 again on the seamed tree** (P5, reproducing N5) | exactly **2** `--- PASS:` |
| | `go build ./host/broker/`; `go vet ./host/broker/` | **rc=0 / rc=0** (P1) | rc=0 / rc=0 |
| **AC5(a)** | `go test ./host/broker/ -run '^TestWarmUpRunsExactlyOnceUnderABoundedContext$' -v -count=1` | **0 `--- PASS:` lines**, rc=0, `[no tests to run]` — the test does not exist (P5) | exactly **1** `--- PASS:`; then MUT-WARM-DROP and MUT-WARM-UNBOUNDED each → exactly **1 `--- FAIL:`** with `go vet ./host/broker/` rc=0 on the mutated tree, each restored sha-identical, each re-run → 1 `--- PASS:` |
| **AC5(b)** | *(no new command)* — names which of AC2's assertions carries the warm-up's proof | see AC2 | AC2's 20 logs each show `exec_started`, `forked`, and the single kill record |
| **AC6** | `git status --porcelain`; `shasum -a 256 host/broker/handlers.go` | **0 lines** (P0); base hash `841987499da53f115584c477d9306193fa87130e484d1255d853dae6e10b191c` | `git status --porcelain` identical before/after the whole T5 run; `handlers.go` hash identical to its **post-T1 committed** value `666a5d420ce2bd11bb69ce93aada9f27838de919e39bd453160b67c7e4642d70` — **not** the base hash |
| **AC7** | `go test ./host/broker/ -run '^TestBrokerTestsDoNotCallParallel$' -v -count=1` | **0 `--- PASS:` lines**, rc=0, `[no tests to run]` — the file does not exist (P5/P6) | exactly **1** `--- PASS:`; then MUT-GUARD-PARALLEL → exactly **1 `--- FAIL:`** naming `handlers_test.go`, `go vet` rc=0, restore byte-identical, re-run → 1 `--- PASS:` |

**Not an acceptance criterion, deliberately** (round-3 quorum): any cold-run majority, any
darwin-only arm, MUT-WARM-SKIP, the `set -m` escape fixture, and the queue row's `-count=20`. If a
gate's verdict is a threshold over a run population, it does not belong in this sprint.

### Hold set — must not move

`host/boundary` (`wantFileCount = 1` lives there and only there, P6 — M1 touches nothing in it) ·
every `.ail` file, `scripts/verify_ail.sh`'s `LEG1_MODULES`, `scripts/verify_go.sh`, `.github/`,
`tools/launchd/*` · `writeExecutable` (five other tests share it) · the 2 s assertion bound ·
`handlers.go:116`'s overflow-path `cmd.Process.Kill()` · `host/store/store.go` · the doc's §6
decision condition (this sprint arms it; it authorizes nothing).

**M1 touches no `.ail` file, so the five `world/*.ail` gate pins are not in play.** If any task
finds itself editing a `.ail`, stop and escalate — that is out of scope by construction.

---

## 8. Flake-attribution protocol for the executor (the ~0.76% is in the test being modified)

`TestHandlerTimeoutKillsTheWholeProcessGroup` has a ~0.76% base failure rate per isolated run
(~1.5% per full `verify_go.sh`, two legs). **A single red run of that test is not evidence of a
regression.** Post-T4 the two are distinguishable by the diagnosis the sprint just built:

- **base flake**: elapsed ≈5.0–5.5 s, `err` unwraps to `*HandlerTimeoutError`, and the kill record
  is populated (count 1). Re-measure **once, in isolation, at the arm level**, and attribute by the
  recorded signature. Record the diagnosis verbatim in the sprint log — it is the artifact §6 needs.
- **a regression this sprint caused**: anything else — a red on the marker assertion, a red on the
  kill-record count, `errors.As` false, a vet/build failure, or a red that reproduces on every run.

The committed test itself **never retries and never skips**. This protocol binds the operator's
re-measurement of an arm, not the gate.

---

## 9. Planner verification log

All rows measured first-party at `002b13f`, 2026-08-13, on this rig.
`export AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.25.6` for every Go row.
Repo writes: **one sha-verified one-shot on `host/broker/handlers.go` plus two throwaway probe
`_test.go` files, all removed; `git status --porcelain` = 0 lines and `handlers.go` restored to
`841987…` before this plan was written** (P0 re-run at the end).

| # | Claim | Command | Observed |
|---|---|---|---|
| P0 | base tree | `git rev-parse --short HEAD`; `git status --porcelain \| wc -l` | `002b13f`; **0** — re-run after all probes: **0**, `handlers.go` sha `841987499da5…` |
| P1 | both gates green at base; real CI budget | `./scripts/verify_ail.sh`; `time ./scripts/verify_go.sh` | ail **rc=0** (`✓ 4 required identities verified, 14 named tests pass`, `9/9 steps`); go **rc=0**, **4:56.89 wall**; `host/broker` **90.233 s** plain leg, **193.374 s** race leg; watchdog is 600 s over the whole race leg (`verify_go.sh:113-126`); the `76.9 s` the doc cites is a stale comment at `verify_go.sh:110-111` |
| P2 | code shape | read `handlers.go:75-135,1-30`, `handlers_test.go:110-135,720-760`, `broker.go:1-60,101-310`; `grep -n 'syscall.Kill\|cmd.Cancel'` | `cmd.Cancel` at `:96`, sole group kill at `:100`; `handlers.go:116` `cmd.Process.Kill()` is the overflow path; test at `:729` with `_, _, _ = session.Invoke` and `elapsed > 2*time.Second`; `handlerSession` at `:73` takes any `Handler`, hardcodes `&handlerRecordingStore{base: openTestStore(t)}`; `objectStore` is an interface (`broker.go:37`); `invoke` returns `&EffectFailedError{cause: err}` at `broker.go:284` with `Unwrap` at `:140` |
| P3 | AC2 base + the vacuous-green control | `go test -c -race -o /tmp/w16p-broker-race.test ./host/broker`; 20 × isolated `-test.run '^TestHandlerTimeout…$' -test.count=1 -test.v`; control `-test.run '^TestNoSuchTestName$'` | **pass=20 fail=0 neither=0**; control **rc=0 with 0 `--- PASS:` lines** and `testing: warning: no tests to run` — V-G/W4 reproduced, hence the counted-lines rule on every AC |
| P4 | the seam applies to N5's exact bytes | `cp` backup; apply seam; `shasum -a 256`; `git diff --stat` | `841987499da53f11…` → **`666a5d420ce2bd11…`** — identical to N5's recorded hashes; diff stat **`7 insertions(+), 1 deletion(-)`**, i.e. **+7/−1, not the doc's +5/−1** |
| P5 | every `-run` AC's base state | `go test ./host/broker/ -run <pat> -v -count=1` for `TestStallDiagnosis`, `^TestWarmUpRunsExactlyOnceUnderABoundedContext$`, `^TestBrokerTestsDoNotCallParallel$`, and the two boundary gates; `go build`/`go vet ./host/broker/` | first three: **rc=0, 0 `--- PASS:`, 0 `--- FAIL:`**, `[no tests to run]`; boundary gates **exactly 2 `--- PASS:`** (0.27 s, 1.46 s) at base **and 2 again on the seamed tree**; build rc=0, vet rc=0 |
| P6 | guard grammar census + no file-count pin | `grep -rn '\.Parallel' host/broker/` (rc); control `grep -rln '\.Helper' host/broker/`; `grep -rln --include='*_test.go' 't\.Parallel' host/ cmd/` (rc); control `t.TempDir`; `find host cmd -name '*_test.go' \| wc -l`; `grep -rn --include='*.go' wantFileCount host/`; guard prototype run | `.Parallel` **0 hits, grep rc=1**, control **9 files**; repo-wide `t.Parallel` **0, rc=1**, control **27 files**, **46** test files (R4 reproduced); `wantFileCount` hits **only** `host/boundary/allowlist_world_test.go:1163-1165` (R5 reproduced); guard prototype `--- PASS`, `enumerated=15 files, Helper selectors=62, offenders=[]` |
| P7 | **the mutex is not `-race`-provable** | unsynchronised recorder variant, `go test -c -race`, 20 isolated runs counting `DATA RACE`; known-positive control (two goroutines incrementing an int) in the **same** binary; read `$GOROOT/src/os/exec/exec.go:781-825,925-940` | **DATA_RACE=0, pass=20, fail=0**; control **`WARNING: DATA RACE` + `--- FAIL`** — the detector is live; source: `watchCtx` calls `c.Cancel()` then `resultc <- ctxResult{…}` (`:805-820`), `Wait` receives at `:930` → happens-before edge. **AC2's stated null case cannot occur** |
| P8 | the full modified-test shape works, warm-up included | prototype of §5 T4 (marker fixture, warm-up, delegating recorder, mutex, all assertions), seamed tree, `go test -run '^TestProbeRecorder' -v` | `--- PASS` ×2; warm-up **217–244 ms**; elapsed **104–107 ms**; kill `count=1 offset=102–103 ms pgid>0 errno=<nil>`; markers `exec_started=true forked=true survived=false`; `err=*broker.EffectFailedError`, **`errors.As`=true, bare type assertion=false** |
| P9 | **MUT-KILL-NEUTER is invisible to the doc's recorder and lethal to a delegating one** | mutate `killGroup`'s body in `handlers.go` (`-pgid`→`pgid`), sha `3dc74e2…`, `go vet` rc=0; run both recorder shapes | delegating: **`--- FAIL` 5.168 s**, `survived=true`, `count=1 errno=nil offset=103 ms` — the full H1 signature; doc-verbatim re-implementing: **`--- PASS` 105 ms**, `survived=false` — **the mutation is not observed at all** |
| P10 | MUT-BOUND-LOOSE composes as specified, and `survived` must not be asserted | keep MUT-KILL-NEUTER, change the prototype's bound 2 s→10 s, run delegating arm; restore probe sha-identical | **`--- PASS` at 5.160 s with `survived=true`** — green, exactly as §8 requires, and only because no `survived` assertion exists (§0.2 C4) |
| P11 | MUT-SEAM-BYPASS is lethal through the record, not the clock | re-seam (`666a5d4…`), then revert `Cancel` to the inline kill (sha `45abc7c…`), `go vet` rc=0, run delegating arm | **`--- FAIL` at 103.6 ms** — `elapsed` green, `count=0`, `offset=0s`: the kill-record assertion is what kills it |
| P12 | MUT-GUARD-PARALLEL is lethal in two grammars, and the red is the guard's | add `t.Parallel()` to the group-kill test (sha `f29c7cd…`), `go vet` rc=0, run guard; then the method-value form `var tt = t.Parallel`; restore `handlers_test.go` (`git diff --stat` empty) | call form: **1 `--- FAIL:`**, `offenders=[handlers_test.go:730:4]`; method-value form: **1 `--- FAIL:`**, `offenders=[handlers_test.go:730:13]`; `go vet` rc=0 in both — the grep a lesser guard would have used sees neither |

---

## 10. What this sprint explicitly does NOT do

- **Not** planning, implementing, or authorizing **M2** (the post-reap re-sweep). §6's decision
  condition is *armed* by this sprint and satisfied by nothing in it.
- **Not** re-opening resolution A. The seam lands; that is ratified mission state (`option A`,
  2026-08-13T06:12:23Z).
- **Not** reintroducing any sampling gate: no cold-run majority, no darwin-only criterion, no
  MUT-WARM-SKIP, no `set -m` escape arm, no `-count=20` before/after.
- **Not** claiming the mechanism is confirmed. §3's kill-then-fork race stays a hypothesis; this
  sprint's value survives its refutation.
- **Not** adding a retry or a skip, **not** moving the 2 s bound, **not** touching
  `writeExecutable`, `host/boundary`, any `.ail`, either verify script, CI, or the frozen core.
- **Not** removing the §5.4 mutex despite §0.2 C3 — it is ratified, free, and defensible as
  defence-in-depth. Only the *claim about what proves it* is corrected.

## 11. Risks

| Risk | Likelihood | Mitigation |
|---|---|---|
| the executor writes the doc's §5.4 recorder verbatim and reports AC3 as failing | **high without this plan** | §0.2 C1 and §5 step 6 state the delegating shape as binding, with the measured PASS/FAIL pair |
| the ~0.76% base flake reds a T5 arm and is read as a regression | ~1.5% per full gate run | §8's attribution protocol; the diagnosis this sprint builds is the discriminator |
| a mutant that does not compile reds "in the predicted direction for the wrong reason" | medium | every row asserts `go vet ./host/broker/` rc=0 on the mutated tree before its verdict is read |
| `git checkout --` used to restore an uncommitted mutant, destroying the executor's work | low, catastrophic | all restores are `cp` from `/tmp/w16_backup/`, with a sha assertion after |
| AC1's two `verify_go.sh` runs eat ~10 min of the budget | certain | priced into the 0.7 day; run them once at the start of T5 and once at the end, never per task |
| the guard false-positives on some future unrelated `X.Parallel()` | 0 today (P6) | declared residual, same as the doc; the failure message names the file:line so the diagnosis is one line long |

## 12. Open question for the human — none blocking

The sprint is fully specified and needs no decision to proceed. One item is worth a human's eye
**after** landing, and is recorded rather than acted on: §0.2 C3 means a blocking round-2 objection
was satisfied with a mutex that guards a race Go's `os/exec` contract forecloses, and the AC
written to prove it cannot fail. Nothing in this sprint depends on resolving that — the mutex stays
— but the quorum protocol's habit of adopting a reviewer's `proposed_fix` *verbatim* without
measuring the objection's premise has now produced one gate that could not fail (this) and, one
round earlier, two that could only pass by luck (round 3's own confession). That is a pattern about
the loop, not about this item.

---

**SPRINT_PLAN_PATH**: `design_docs/planned/w-broker-base-flake-sprint-plan.md`
**SPRINT_JSON_PATH**: `.ailang/state/sprints/w-broker-base-flake.plan.json`
