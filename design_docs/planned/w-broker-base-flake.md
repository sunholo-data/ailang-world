# w-broker-base-flake — localise the 0.76% leak, de-align the exec-assessment race, and prove the gate without sampling it

- **Status**: Planned
- **Item**: queue item 16, `w-broker-base-flake`, clause-3 (filed iter-73)
- **Estimated**: ~0.5 day for M1 (this doc's landable scope); M2 is decision-gated (§6, §12)
- **Measurement base**: `b2c3f89` (dev HEAD, clean tree), 2026-08-13
- **Instruments**: pinned `ailang` **v0.30.0** at `/tmp/ailang-v0300/ailang` (`AILANG v0.30.0`, W-V);
  `go1.25.6` via `GOTOOLCHAIN` for every build/test (local `go` is 1.26.4, which `verify_go.sh`'s
  own denylist rejects, W9); the controller's process probe `/tmp/pgprobe/main.go` plus two
  per-run variants derived from it (W5–W7)
- **Files touched (M1)**: `host/broker/handlers_test.go` (modify one test, ~+70 lines),
  `host/broker/handlers_stall_diag_test.go` (NEW, ~150 LOC: diagnosis helper + two committed
  attribution arms), and `host/broker/handlers.go` (**+5/−1 lines**: the `killGroup` seam of
  §5.4, behaviour-identical, proven gate-clean by one-shot N5). Round 1 claimed zero production
  bytes; the revision trades that property away deliberately — see the Quorum verification log.
  The one production *mutation* in this doc (MUT-KILL-NEUTER) remains a one-shot with a `cp`
  backup, restored and sha-verified.

This doc is built on the iteration-78 controller measurements (cited as V-A..V-H) and on
first-party re-derivations and extensions made in this design session (cited as W1..W16, §14).
Where the queue row is now known wrong, this doc says so plainly and overrides it (§1). Where
this doc **disagrees with the controller's own reading** of its measurements, that is stated
too (§4) — two of the controller's conclusions are narrowed by new measurements taken here,
and one new mechanism-shaping fact was found that neither the row nor the controller had
(§3: the fixture's first-exec latency on darwin is ≈ the test's own deadline).

---

## 1. Problem, corrected from the queue row

`host/broker/handlers_test.go:729` `TestHandlerTimeoutKillsTheWholeProcessGroup` writes a
`#!/bin/sh` fixture whose body forks a 5-second grandchild (`sleep 5 &` + `wait`), runs it
through a `GitHandler` with `ExecTimeout: 100ms`, and asserts `session.Invoke(...)` returns in
under 2s (W1). The queue row (filed iter-73, on `TR.B`-planner measurements) is wrong in three
load-bearing places, and right in two:

| Row claim | Status at `b2c3f89` | Evidence |
|---|---|---|
| fails ~18% (2 of 11 isolated runs) | **WRONG at head**: 1 failure in 132 controller runs ≈ **0.76%**; my own 20/20 `-race` isolated runs passed | V-B; W4 |
| mechanism = "process-group kill timing under parallel load" | **REFUTED as stated**: 1,987 probe runs of the exact kill mechanics — 0 leaks, 0 `ESRCH`, 0 nil-`Process`, incl. 12-way parallel arms. But see §3–§4: the probe sampled only the *warm-fixture* regime | V-D; W5–W8 |
| prove the fix with `-count=20` that "reds before" | **VACUOUS as prescribed**: at p=0.0076, a `-count=20` run reds with probability 1−(1−p)²⁰ ≈ **14%** — a proof that fails to prove ~86% of the time is exactly the coin-flip gate S6 forbids. Overridden (§7), the same call iter-76 made against item 12's own row | V-C; S6 |
| the flake is real | **CONFIRMED**: reproduced once at head, `--- FAIL: … (5.28s)` under `-race` | V-A |
| without `AILANG_BIN` the package is red 100% | **CONFIRMED and correct behaviour**: `episode_test.go:167-169` `t.Fatal`s on an unset `AILANG_BIN` (never a skip). Every acceptance command below exports the pin | V-H; W10 |

One more defect the row never mentions, visible in one line of source: **the test throws away
its best evidence at the exact moment it matters.** `handlers_test.go:740` reads
`_, _, _ = session.Invoke(...)` — the returned error, which distinguishes a
`*HandlerTimeoutError` from a store failure from anything else, is discarded, and the failure
message then asserts a *cause* ("the kill missed the forked grandchild") that the elapsed-only
assertion cannot distinguish (W1). The two sibling timeout tests in the same file
(`TestGitHandlerTimeoutWritesFailureRecord`, `TestModelHandlerTimeoutWritesFailureRecord`,
W14) capture and type-check their errors; this one does not.

## 2. What the one failure's shape already tells us

The test times all of `Invoke`, whose live path is: decide → hash → `PutObject(request)` →
`AppendNextEffectIntent` → `GitHandler.Execute` (`os.MkdirTemp` + `runBounded` +
`os.RemoveAll`) → failure `putRecord` + `AppendEffectOutcome` — the store being a real SQLite
database via the pure-Go `modernc.org/sqlite` driver, opened `:memory:` by `openTestStore`
(W3). The controller is right (V-F) that a 5.28s elapsed is *formally* consistent with a stall
anywhere on that path — the rule-3i shape.

But the elapsed **value** is not mute. Both ways of forcing the leak measured this session
land in a tight band set by the grandchild's natural exit:

- a grandchild deliberately escaped from the process group (`set -m`): 20/20 leaks at
  **5.099–5.179s** (W6);
- the kill deliberately narrowed to the direct child only (the MUT-KILL-NEUTER analogue):
  5/5 leaks at **5.137–5.172s**, and with the full proposed fixture, **5.155s with the
  `survived` marker present** (W7, W16).

V-A's 5.28s sits exactly one `-race`-overhead above that band. A SQLite or `MkdirTemp` stall
has **no mechanism coupling its duration to the fixture's `sleep 5` parameter**; a held stdout
pipe does — `io.ReadAll` returns when the last writer dies, and the grandchild's natural death
is at ~5.15s. The elapsed value is therefore a *fingerprint of the pipe being held until the
grandchild exited on its own*, which re-implicates the kill **window** (not the kill *logic*,
which V-D exonerated) and demotes the store. §4 records this as a disagreement with V-F's
suspect ordering; §5's instruments are designed so the next failure settles it with evidence
rather than argument.

## 3. New finding: the fixture's first-exec latency equals the test's deadline (darwin)

Measured this session, first-party, three independent ways (W6–W8, W15):

1. A **freshly written** `#!/bin/sh` script's first exec on this darwin machine costs
   **~101–223ms** (8 cold runs: 223.4, then 101–103ms each); the **same file's** subsequent
   execs cost **~3–12ms** (W8).
2. The real test writes a **fresh fixture inode on every run** (`writeExecutable` →
   `os.WriteFile` into a fresh `t.TempDir()`, `handlers_test.go:118-125`, W1) — so on darwin
   the fixture usually has not even forked its grandchild when the 100ms deadline expires.
   Directly observed: with the kill deliberately broken to child-only, **cold fixtures leak
   only 2/8** (the other 6 end at ~102ms — nothing existed to survive) while **warm fixtures
   leak 5/5** (W7).
3. The controller's probe reuses **one fixture file across all 1,987 runs** (`/tmp/pgprobe`
   takes the script path as `os.Args[1]`, W5) — it therefore sampled the ~3ms **warm** regime
   exclusively, where the fork completes ~90ms before the kill and no interesting race exists.

Consequences, in increasing order of importance:

- **The probe's exoneration is regime-scoped.** V-D remains valid for what it measured —
  `runBounded`'s kill logic is sound in steady state — but the real test lives in the cold
  regime, where fork, exec-assessment completion, and the kill all land within a few
  milliseconds of each other. This is the faithfulness divergence the controller asked to be
  checked (§4).
- **The test is partially vacuous per-run on darwin.** In the majority of local cold runs the
  deadline expires before the grandchild exists, so the group-kill-vs-grandchild property the
  test exists to guard is never exercised — the run passes because there was nothing to kill.
  A green that does not run its own scenario is S6's exact target. (Linux CI has no such
  assessment latency, so CI runs likely do exercise it — but nothing currently *observes*
  which mode any run was in.)
- **A refined mechanism hypothesis — labeled as such.** `kill(-pgid)` signals the group
  members that exist when the syscall runs; a child forked concurrently with the kill can miss
  the sweep while inheriting the dying group. With darwin's ~100ms assessment latency
  *aligning* the fixture's fork with the 100ms deadline, the fork-vs-kill race window — normally
  unreachable — is sampled on every darwin run, at a width plausibly matching the measured
  0.76% (and widening under the parallel load `TR.B` ran under, matching its 18%). This is a
  HYPOTHESIS in the row's sense — nothing below assumes it; §5's instruments are chosen to
  confirm or refute it on the next occurrence. Production handlers are outside this regime:
  the default `ExecTimeout` is 30s (`handlers.go:14`), so real fixtures fork ~29.9s before any
  kill.

## 4. Where this doc disagrees with the controller

Invited explicitly by the design brief; both disagreements are narrowings, not reversals.

1. **V-D/V-B regime scope.** The probe is mechanically faithful to `runBounded`
   (`CommandContext`, `Setpgid`, group-kill `Cancel`, `StdoutPipe`, `Stderr=Stdout`, capped
   `ReadAll`, `Wait` — read line-by-line, W5) with three divergences: it **reuses one fixture
   inode** (the material one, §3), it sets neither `cmd.Dir` nor `cmd.Env` (the real path sets
   both, `handlers_git.go:52-69`), and its cap is 1MB vs 8MB+1 (immaterial). The 1,987-run
   exoneration therefore covers the warm regime only; the cold regime was measured here at
   8-run scale and behaves categorically differently (W7–W8).
2. **V-F's suspect ordering.** "The store/SQLite path is now the LEADING suspect" does not
   survive the elapsed fingerprint (§2): both forced leak modes reproduce V-A's elapsed to
   within race overhead, and no store mechanism couples a stall's duration to the fixture's
   sleep parameter. The store stays on the suspect list (the instruments time it anyway), but
   the leading suspects are the two that the fingerprint fits: a kill that never reached the
   pipe-holder, or a kill that fired after the holder's natural death. V-F's core claim — that
   the *assertion* cannot distinguish these — is correct and is what §5 fixes.

V-G (the vacuous `-test.run` green) was reproduced first-party — a bogus test name prints
`no tests to run` / `PASS` / rc=0 (W4) — and is adopted as a hard rule on every acceptance
criterion below: a `-run`-selected command counts `--- PASS:`/`--- FAIL:` lines, never rc.

## 5. M1 — the localisation trap (this item's landable scope)

Design rule: every instrument must (a) not perturb the timing it measures, (b) leave a durable
artifact when CI catches the 0.76% (the go-test failure log, which both `verify_go.sh` legs
and the CI job preserve), and (c) be non-vacuous — each instrument has a mutation that kills
it (§8).

### 5.1 Capture the error

`_, _, _ = session.Invoke(...)` becomes a captured `(result, ref, err)`. The test explicitly
asserts that `err` is a `*HandlerTimeoutError` on the pass path (preventing silent non-timeout
early returns), and the failure branch prints `%+v` of the error alongside elapsed. A
`*HandlerTimeoutError` says `runBounded`'s deadline machinery ran to completion; anything else
immediately relocates the fault. Zero timing cost. (Pass-path wording per the round-1 quorum's
gemini fix, adopted verbatim.)

### 5.2 Phase decomposition, entirely test-owned

The test already owns both seams it needs — no production change:

- **Handler seam**: `handlerSession(t, effect, scope, handler)` takes any `Handler` (W3); wrap
  the `GitHandler` in a timing shim recording `Execute` enter/exit (`time.Now()` twice).
- **Store seam**: the session's store is already the test-local `handlerRecordingStore`
  wrapper (W3); a timing variant records a monotonic timestamp per store call
  (`PutObject`, `AppendNextEffectIntent`, `AppendEffectOutcome`).

Together these decompose `Invoke` elapsed into *pre-handler broker work* / *Execute window
(MkdirTemp + runBounded + RemoveAll)* / *post-handler store writes*, attributing a stall to
one of the three without touching `broker.go` or `handlers.go`. Cost per boundary: one
`time.Now()` — nanoseconds against a 2s bound.

### 5.3 Fixture markers: non-vacuity proof, and the honest limit of `survived`

The fixture body (composed by the test, which bakes an absolute marker directory into the
script source — `GitHandler`'s env is a fixed allowlist, so no env variable can reach it, W3)
becomes:

```sh
#!/bin/sh
set -eu
if [ "${W16_WARM:-}" = "1" ]; then exit 0; fi
: > "<markdir>/exec_started"
sleep 5 && : > "<markdir>/survived" &
: > "<markdir>/forked"
wait
```

Marker semantics, each verified live this session (W16), **corrected by the round-1 quorum**:

- `exec_started` + `forked` present, `survived` absent, elapsed ~103ms — healthy run, and
  **proof the run was non-vacuous** (the grandchild existed and was killed). This pair is
  asserted on every pass, converting §3's silent per-run vacuity into a visible signal.
- `survived` **present** after a >2s elapsed proves exactly one thing: **the sleeper reached
  its natural 5s completion before any effective kill**. The marker is written at
  sleep-completion in *both* the never-killed case (H1) and the case where the deadline
  machinery stalled ~5s and the kill fired only after the natural death (H2-late) — the file
  looks identical either way, and its mtime bounds nothing about when or whether
  `syscall.Kill` ran. Round 1 claimed the marker alone read "never killed"; that claim was
  wrong and is withdrawn. The H1-vs-H2 discrimination is carried by §5.4's kill-boundary
  record, not by this marker.
- `survived` **absent** after a >2s elapsed — the sleeper was killed *late*, somewhere between
  the 2s bound and its 5s natural exit (an H2 shape on any reading; §5.4's record then says
  when the kill actually fired).

Perturbation discipline: the marker writes sit **after** `sleep 5 &`, so the fork's timing
relative to `cmd.Start` is byte-for-byte unmoved from today's fixture; the writes race the
sleeper, not the fork. The `survived` writer changes the grandchild from `sleep` to a
subshell-plus-`sleep` — both inherit the group and both hold the pipe, and the equivalence
was measured: group kill 103ms clean, child-only kill 5.155s leak (W16), matching today's
fixture's behaviour in both arms (W7).

### 5.4 The kill-boundary recorder — M2's seam brought forward (resolution A)

The round-1 doc deferred a `killGroup` seam to M2 and let `survived` + `*HandlerTimeoutError`
authorize it. The quorum showed that evidence is ambiguous (§5.3, Quorum log): the one bit the
decision gate rested on cannot distinguish H1 from H2-late. Two honest resolutions existed —
(B) keep M1 production-free, narrow its claim to *localisation*, and pay a second trip through
a 0.76% event before mechanism selection; or (A) land the seam now and make the single captured
occurrence decisive. This doc takes **(A)**, for one reason argued in full: the event fires
roughly once per ~130 CI runs, so each firing must yield a mechanism-grade record — a trap that
answers *where* but not *why* (option B) converts one rare-event wait into two. The seam is
also exactly the seam M2 needs for its fault-injection arm, so nothing here is throwaway. The
price, stated plainly: **M1 stops being a zero-production-byte milestone**, which round 1
counted as a selling point. The trade is 6 lines of behaviour-identical production code against
an entire second quarter-of-CI-runs wait; this doc judges that cheap and says so on the record.

The production change in `handlers.go` (+5/−1, applied/gated/reverted as one-shot N5):

```go
// killGroup is the cancellation kill boundary, a package-level seam so the
// timeout tests can observe the kill's time, target and errno directly.
var killGroup = func(pgid int) error {
	return syscall.Kill(-pgid, syscall.SIGKILL)
}
```

with `cmd.Cancel`'s kill line becoming `return killGroup(cmd.Process.Pid)` — the nil-`Process`
guard stays in `Cancel`, and `handlers.go:100` is the sole group-kill site in the package's
production code (N1), so this one seam covers the whole boundary. N5 proved the change
gate-clean by measurement, not argument: build + vet green, both `host/broker` AST gates
`--- PASS`, the group-kill test 5/5 green under `-race`, restore byte-identical.

The modified test swaps in a **recording wrapper that still performs the real kill** (restored
via `t.Cleanup`; safe because the file has no `t.Parallel`, W12): it records invocation count,
monotonic offset from the `Invoke` start, target pgid, and returned errno — the exact evidence
set the blocking objection demanded. The recorder is the controller's own probe pattern
(`/tmp/pgprobe/main.go:31-41` captures `killErr` and `nilProc` at a cost of one assignment per
field, N4) moved behind the seam.

**Pass-path assertion, committed**: this test's deadline *always* expires, and the direct
child's only exit path before its 5s `wait` completes is the kill itself — so a healthy run
must show exactly one `killGroup` invocation, errno nil, at offset ≥ the 100ms deadline. Base
evidence that errno nil is the invariant and not a hope: 1,987 probe runs, 0 `ESRCH`, 0
nil-`Process` (V-D). A zero-count or non-nil errno on an elapsed-green run is itself a red —
which is also what makes the seam un-bypassable (MUT-SEAM-BYPASS, §8).

**Failure-path record**: the diagnosis prints the kill record (count, offset, pgid, errno)
alongside elapsed, phase split, error, and markers. §6's decision condition reads *this*, not
the marker alone.

### 5.5 Warm the fixture inode (the one behavioural change, honestly labeled)

Before the timed `Invoke`, the test execs the fixture once directly with `W16_WARM=1`
(exiting at the guard). Measured: the guarded first exec pays the ~104–269ms assessment; the
next exec of the same inode costs ~12ms (W15). Effects:

- On darwin, the fork now completes ~85ms **before** the deadline instead of racing it — the
  test *actually exercises the grandchild-kill path on every run* (today it mostly doesn't,
  §3), and the `forked`-marker assertion holds with ~8× margin.
- The aligned-constants regime — the only regime in which the flake has ever been observed —
  is removed from the test. If §3's hypothesis is right, the flake goes with it.

What this is **not** claimed to be: a diagnosis. If the 0.76% has a different mechanism, the
warm-up may not remove it — which is exactly why the §5.1–5.4 trap lands in the same commit
and stays armed. It is also neither a retry (the timed run happens once; the warm-up exec is
a different, unmeasured invocation whose only job is inode state) nor a skip (the assertion
is unchanged at 2s and the test runs everywhere it ran before).

### 5.6 Committed attribution arms — pure Go, deliberately NOT the subprocess escape

`handlers_stall_diag_test.go` extracts the decomposition into a helper
(`invokeWithStallDiagnosis`) used by the real test and by two committed arms that FORCE a
stall deterministically and assert the diagnosis attributes it to the right phase:

- **`TestStallDiagnosisAttributesHandlerWindow`** — a pure-Go `Handler` stub that sleeps 1.5s
  then returns a wrapped `ErrHandlerTimeout`; requires the failure diagnosis to place ≥1.2s in
  the Execute window and to contain the error text.
- **`TestStallDiagnosisAttributesStoreWindow`** — a store wrapper whose `AppendEffectOutcome`
  sleeps 1.5s; requires ≥1.2s in the post-handler window.

(The kill-record assertions of §5.4 live in the modified real test, not in these arms — the
arms' stub handler never reaches `runBounded`, so they exercise the diagnosis plumbing only.)

Why pure-Go stalls and not the `set -m` process-group-escape fixture (which forces the *real*
leak): the escape was measured **non-deterministic in exactly the regime a committed test
would run in** — 20/20 leaks on a warm fixture but **1/8 on fresh fixtures** (W6), because the
same first-exec latency of §3 usually lets the kill win before the escape exists. A committed
arm that reds 7 runs in 8 for environmental reasons would be this item manufacturing the very
flake class it was filed against. The escape fixture is used only in warmed one-shot form
(§7, AC5-adjacent) where the operator controls the regime. CI cost of the two arms: ~3s plain
leg, ~5–8s race leg, against the race leg's 600s watchdog with `host/broker` at 76.9s
critical path (W9) — stated so the budget change is a decision, not a drift.

## 6. The fix question, answered honestly: M2 is decision-gated, not landed

The row demands "fix or correctly bound it — never a retry or a skip". M1 *bounds* it: the
assertion stays at 2s, the rate at head is 0.76% per isolated run (~1.5% per full
`verify_go.sh`, which runs the test in both legs), the de-alignment (§5.5) removes the only
regime the flake was ever observed in, and any recurrence now writes a diagnosis instead of a
mystery. The *fix* — if §3's kill-then-fork race is confirmed — is a re-sweep of the process
group after `cmd.Wait` reaps the child (a surviving member keeps the pgid alive; a second
`kill(-pgid, SIGKILL)` collects any fork-race escapee, `ESRCH` meaning "nothing to collect").
Its deterministic before/after proof drives the `killGroup` seam — which M1 now lands (§5.4) —
because the race itself cannot be forced from outside; M2's remaining scope is the re-sweep
plus its arms.

The re-sweep is a behavioural change to the one frozen decision-pipeline package on an
unconfirmed mechanism. It does NOT land in this item. Round 1's decision condition —
`survived` present plus a `*HandlerTimeoutError` — was defeated by the quorum: that pair is
also produced when the deadline machinery stalls ~5s and the kill fires after the sleeper's
natural death (H2-late), a fault a post-reap re-sweep would do nothing for. The corrected
decision condition requires direct evidence from the kill boundary. **M2 (the re-sweep) is
authorized only when one captured diagnosis shows ALL of:**

1. `err` is a `*HandlerTimeoutError` (the deadline machinery ran to completion);
2. the §5.4 record shows **exactly one `killGroup` invocation, errno nil, at monotonic offset
   ≈ the 100ms deadline** — the kill completed successfully long before the sleeper's ~5s
   natural exit, so "kill fired late" and "kill never ran" are both excluded by measurement;
3. the Execute window ends in the ~5.0–5.5s natural-death band (§2's pipe-held fingerprint);
4. `survived` is present.

That conjunction proves a successfully-signalled group whose forked member outlived the sweep.
With this fixture the grandchild has no job-control mechanism for *leaving* the group
(`#!/bin/sh`, no `set -m`), so the surviving member holds the dying pgid — the exact state the
post-reap re-sweep collects, and M2's fault-injection arm re-verifies that through the same
seam. **If the diagnosis shows anything else** — zero kill invocations, non-nil errno, a late
kill offset, or `survived` absent — the follow-up is whatever that evidence names instead of
the re-sweep; **if nothing recurs in a quarter of CI runs**, the item closes as bounded with
the trap left armed. A doc that guessed the sweep now would be repeating the row's own error —
a prescription outliving its diagnosis, the exact shape iter-76 killed in item 12's row.

## 7. Replacing `-count=20`: proofs that cannot pass by luck

The row's before/after proof sampled a 0.76% event (V-C). Every proof below either forces the
failure deterministically or measures a precondition with a ~100%-vs-0% effect size. All run
with `AILANG_BIN=/tmp/ailang-v0300/ailang` exported and `GOTOOLCHAIN=go1.25.6`; every
`-run`-selected command asserts its `--- PASS:`/`--- FAIL:` lines per §4/V-G.

- **P1 — the gate observes the kill mechanism (MUT-KILL-NEUTER).** One-shot: `cp` backup of
  `handlers.go`, change `syscall.Kill(-pgid, …)` to `syscall.Kill(pgid, …)` **inside the §5.4
  `killGroup` seam** (same one-line child-only-kill mutation as round 1, relocated with the
  kill), sha256 must move; run the modified (warmed) test `-count=5`: **5/5 `--- FAIL`** at
  ~5.1–5.3s; restore byte-identical; `-count=5` green. Measured analogue this session: warm
  child-only kill leaked 5/5 at 5.137–5.172s and wrote `survived` (W7, W16). Warming is what
  makes this deterministic — cold, the same mutation leaked only 2/8 (W7), which is the row's
  `-count=20` failure in miniature. **This one-shot now also proves §6's decision condition is
  reachable and truthfully reported**: a child-only kill returns errno nil at the deadline while
  the grandchild survives, so each forced FAIL's diagnosis must show the full H1 signature —
  kill count 1, errno nil, offset ≈100ms, `survived` present, elapsed in the natural-death band.
  A diagnosis that cannot reproduce the signature under forced H1 conditions would never be
  trusted to report a spontaneous one.
- **P2 — the diagnosis tells the truth (committed).** The two §5.6 arms, every CI run.
- **P3 — the precondition proof (darwin one-shot, the before/after for §5.5).** At base, 20
  runs of the *unwarmed* test with markers: the majority must show `forked` absent-at-deadline
  (vacuous-mode; measured analogue ~6/8, W7/W8). After §5.5: 20 runs, **0** vacuous-mode. This
  measures the race *precondition* (fork-vs-deadline alignment) rather than the rare race
  *outcome* — deterministic in aggregate where `-count=20` was a coin flip. Recorded as a
  darwin-only AC: linux has no assessment latency and both arms would read 0.
- **P4 — the pass arm.** Pristine tree, modified test, 20 isolated `-race` runs, 20/20
  `--- PASS` (the weak direction, labeled as such; base measured 20/20 today, W4).

## 8. Mutation table

Landed-proof rule (inherited from item 12): sha256 of every mutated file must move before a
verdict is read; one-shots restore from `cp` backups, never `git checkout --`. Compile gates:
the `_test.go` changes need **`go vet ./host/broker/`, not `go build ./...`** (which does not
compile test files at all — the recorded iter-77 gotcha); the §5.4 seam in `handlers.go` is
production code and IS covered by `go build` (both were run green in one-shot N5 — vet is still
run because the bulk of M1 remains test-side).

| # | Mutation | Measured today (b2c3f89) | Required post-change | Observable that reds | Killed by |
|---|---|---|---|---|---|
| MUT-KILL-NEUTER | `-pgid` → `pgid` inside `killGroup` in `handlers.go` (one-shot, restored) | analogue: warm 5/5 leak @5.14–5.17s, `survived` written (W7, W16) | warmed test reds 5/5, diagnosis shows full H1 signature (§7 P1) | elapsed >2s + kill record + markers | AC3 (P1) |
| MUT-SEAM-BYPASS | revert `Cancel` to the inline `syscall.Kill(-cmd.Process.Pid, …)`, orphaning the seam | this IS the pre-M1 code shape (N1) | modified test reds: elapsed stays green but the recorder saw **zero** invocations, failing the §5.4 pass-path assertion | kill-record count ≠ 1 on a pass run | the modified test itself |
| MUT-ERR-DISCARD | restore `_, _, _ =` at the Invoke call | this IS today's code (W1) | handler-window arm reds: diagnosis lacks the error text | missing `%+v` error in output | committed P2 arm |
| MUT-DIAG-BLIND-EXEC | drop the Execute enter/exit timestamps from the helper | n/a | `TestStallDiagnosisAttributesHandlerWindow` reds | attribution ≠ Execute window | committed P2 arm |
| MUT-DIAG-BLIND-STORE | drop the store-call timestamps | n/a | `TestStallDiagnosisAttributesStoreWindow` reds | attribution ≠ store window | committed P2 arm |
| MUT-MARKER-DROP | fixture stops writing `forked`/`exec_started` | n/a | real test reds on its non-vacuity assertion | `forked` absent on a pass run | the modified test itself |
| MUT-WARM-SKIP | remove the `W16_WARM=1` warm-up exec | this IS today's regime | darwin: `forked`-at-deadline assertion reds in most runs (measured ~6/8 analogue, W7) | vacuous-mode signal fires | AC5 (P3, darwin one-shot) |
| MUT-BOUND-LOOSE | raise the test's 2s bound to 10s | n/a | MUT-KILL-NEUTER one-shot **stops redding** (5.15s < 10s) — the composed one-shot's green IS the kill signal, run as a pair | AC3's red arm goes green | AC3 protocol (one-shot pair) |

Each committed assertion reads its observable *through* the mechanism it guards (the
diagnosis text is produced by the timestamps; the markers by the fixture's own control flow),
not alongside it.

## 9. Acceptance criteria

All from repo root; `export AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.25.6`; every
`-run` command's verdict is its counted `--- PASS:`/`--- FAIL:` lines (rc alone accepted
nowhere — V-G/W4). Base-state results per rule 3e are stated inline.

- **AC1 — baseline green, both gates, before and after.** `./scripts/verify_ail.sh` rc=0 and
  `./scripts/verify_go.sh` rc=0 on the pristine tree (base: green at `b2c3f89`) and again
  after the change. (`verify_go.sh` FATALs without the exports above — base condition.)
- **AC2 — modified test stable at head.** Prebuilt `-race` binary, 20 isolated runs of
  `^TestHandlerTimeoutKillsTheWholeProcessGroup$`: 20× `--- PASS`, each run's log showing the
  markers assertion executed AND the §5.4 kill-record assertion (count 1, errno nil, offset ≥
  deadline) executed. Base: today's test measured 20/20 (W4); this re-runs it modified; the
  seam alone (recorder not yet swapped) measured 5/5 green under `-race` in one-shot N5.
- **AC3 — P1 executed** (MUT-KILL-NEUTER protocol of §8, including the MUT-BOUND-LOOSE pair
  and byte-identical restore).
- **AC4 — committed arms and boundary gates live.** `go test ./host/broker/ -run
  'TestStallDiagnosis' -v -count=1` → exactly 2 `--- PASS` lines (base: 0 — the file does not
  exist, W11); `go test ./host/broker/ -run 'TestRegistryDispatchBindingBoundary|
  TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential' -v -count=1` → exactly 2
  `--- PASS` lines (base with the seam applied: 2, N5); then `go build ./host/broker/` and
  `go vet ./host/broker/` both rc=0.
- **AC5 — P3 executed on darwin** with both arms' vacuous-mode counts recorded (base analogue:
  ~6/8 cold, 0 warm — W7/W15).
- **AC6 — live tree hygiene.** `git status --porcelain` identical before/after the full run;
  `shasum -a 256 host/broker/handlers.go` identical to its pre-AC3 value — which is the
  **post-M1 committed content including the §5.4 seam**, not the base `b2c3f89` hash (the seam
  is a landed change; only the AC3/P1 one-shot must restore byte-identically).

## 10. Conflict surface

- **`host/broker/invoke_boundary_test.go` (TR.B/TR.C AST gates)**: the production walker
  skips `_test.go` at `:145`/`:218` and asserts it skipped a nonzero number (`:223`) — a new
  test file only grows that count; the ≥30 floor (`:212`) counts production files, and the
  §5.4 seam adds no file, only lines to an existing one; the binding detector's complete kind
  set is `invoke-call`/`ctor-live`/`ctor-replay`/`session-type`/`dot-import` with the
  inside-broker exemption pinned by identity to `publish_op.go` (`:274-288`, N2) — a
  package-level `var` and a call to it match none of those kinds. Not argued but **measured**:
  the gate ran `--- PASS` against the seamed tree in one-shot N5.
- **`registry_publish_test.go:1076` `enumerateSubprocessSites`** matches only
  `exec.Command`/`exec.CommandContext` call sites in non-test `.go` under `host/ cmd/`
  (`:1099-1108`), pins the resulting **file set** to a five-file driver map (`:1181-1186`,
  `handlers.go` among them) and pins no line numbers (N3). The seam adds no `exec.Command*`
  site and no file; the new test file's warm-up `exec.Command` is in a `_test.go` — skipped
  (`:1084`). Measured green against the seamed tree in N5. Round 1 deferred this gate-check to
  a future M2 sprint; bringing the seam forward meant doing the check now, and it passed.
- **`host/boundary` `wantFileCount = 1`** (`allowlist_world_test.go:1163`) is scoped to
  `host/boundary`'s own `.go` files; nothing is added there. No file-count pin exists on
  `host/broker` (repo grep, W11).
- **`scripts/verify_ail.sh` / `LEG1_MODULES`** (landed by item 12 at `:135`): no `.ail` file
  is touched or added — the allowlist is unmoved (W11).
- **`verify_go.sh` race-leg watchdog (600s)**: +~5–8s from the committed arms against a
  76.9s package critical path (W9) — comfortably inside budget, stated per the no-silent-caps
  rule.
- **`episode_test.go` AILANG_BIN refusal** (V-H/W10): single-test `-run` commands above don't
  select it, but every AC exports the pin anyway so full-package runs stay meaningful; CI's
  go-verify job already exports it (`ci.yml:144`), job 1 deliberately does not (W10).
- **S3 "why is this not a package?"**: not applicable in the S3 sense — M1 adds no kernel or
  host *surface*: the seam is an unexported package-level var with production behaviour
  identical to the inline call it replaces, plus test instrumentation; M2, if authorized,
  changes one unexported function and will answer for itself.

## 11. Mutation-campaign guidance (the row's real worry, quantified)

The row's stakes were fake kills and falsified inverse arms in future clause-3 mutation work.
At head the numbers are: p ≈ 0.0076 per isolated run of this one test, ~1.5% per full
`verify_go.sh` (two legs). A 32-arm campaign reading package-level verdicts sees ≥1
flake-contaminated arm with probability ≈ 22% (single-leg arms) to 39% (both-legs arms) —
real, but bounded and now *attributable*: post-M1, a flake-red carries the full diagnosis
(elapsed ≈ 5.0–5.5s, phase split, error, `survived` bit, kill record — count/offset/pgid/errno). Protocol for campaign operators: a
red arm whose only failing test is this one with that fingerprint is re-measured **at the arm
level** once, in isolation, and attributed by diagnosis; the committed test itself never
retries and never skips — the interdiction the row rightly imposed binds the gate, not the
operator's re-measurement of an arm.

## 12. Sizing

The row said ~0.5d assuming diagnose→fix→`-count=20`. As re-scoped — M1 only — **~0.5 day
holds even with the §5.4 seam added**, because every mechanic was de-risked live this session:
the seam itself was applied, built, vetted, gate-tested and reverted in one-shot N5 (the
recorder swap that remains is the probe's two-assignment pattern, N4); the wrappers are test
plumbing over seams the file already owns (W3); the fixture body and both marker arms were
executed end-to-end (W16); the warm-up latency is measured (W15); the committed arms are
pure-Go sleeps; and the one-shot protocols were rehearsed via probe analogues (W6–W8). What
is deliberately NOT in the 0.5d: M2 (decision-gated, §6) and any chase of the 0.76% by brute
sampling. If the sprint planner prices the fixture work above the band, the honest fallback is
to land §5.1–5.4+5.6 (the trap, seam included) and demote §5.5 to the follow-up — the trap is
the part that must not slip, and without §5.4 it cannot authorize anything (§6).

## 13. What this item is NOT doing

- **Not** claiming the mechanism is confirmed. §3 is a hypothesis with new supporting
  measurements; the instruments decide, and the doc's value survives the hypothesis failing.
- **Not** landing a production *fix* (the re-sweep stays decision-gated on a captured
  diagnosis, §6). M1's only production change is the §5.4 seam — +5/−1 behaviour-identical
  lines whose entire job is observability; round 1's "zero production bytes" claim is
  deliberately given up, and the Quorum log records why.
- **Not** adding a retry or a skip anywhere, and **not** moving the 2s assertion bound.
- **Not** committing the `set -m` escape arm — measured 1/8 deterministic in the regime a
  committed test runs in (W6); it survives only as a warmed one-shot instrument.
- **Not** running the row's `-count=20` as a red-arm proof (14% coin flip, V-C) — superseded
  by §7's forced and precondition proofs.
- **Not** touching `writeExecutable` (five other tests share it, W1); the group-kill test
  composes its own fixture body.
- **Not** touching `.ail` files, `verify_ail.sh`, `verify_go.sh`, CI workflows, `host/boundary`,
  the frozen core, or anything item 12 landed.

## 14. Verification Log

All rows measured first-party, 2026-08-13, at `b2c3f89`. W-rows are the round-1 design
session; N-rows are the round-2 revision session (quorum response). Binary pins verified in
W-V. Controller rows V-A..V-H are cited above as second-party context; every V-row this
design *leans on* is re-derived or extended below (W4→V-B/V-G, W5→V-D, W7/W8→new, W10→V-H).
Probes live outside the repo (`/tmp/w16probe*`, `/tmp/pgprobe`). Repo writes across both
sessions: this doc, plus **one sha-verified one-shot on `host/broker/handlers.go` in round 2
(N5), restored byte-identical to its base hash before this doc was finished** — the same
one-shot discipline §8 imposes on the sprint (`cp` backup, sha moves, sha restores).

| # | Claim | Command | Observed |
|---|---|---|---|
| W-V | instrument pins | `go version`; `/tmp/ailang-v0300/ailang --version` | `go1.26.4 darwin/arm64` (hence `GOTOOLCHAIN=go1.25.6` on every gate command — verify_go.sh's own denylist rejects go1.26.0–1.26.5, W9); `AILANG v0.30.0` |
| W1 | test shape; error discarded; fresh fixture per run | Read `handlers_test.go:118-125,:729-747` | fixture = `#!/bin/sh\nset -eu\n` + body, fresh `t.TempDir()` file each run; `ExecTimeout: 100ms`; `_, _, _ = session.Invoke(...)` at `:740`; assertion `elapsed > 2*time.Second` at `:743`; failure text names the grandchild mechanism it cannot observe |
| W2 | runBounded mechanics; no post-reap sweep; no WaitDelay | Read `handlers.go:82-131` | `CommandContext` + `Setpgid` + `Cancel`=`syscall.Kill(-pid, SIGKILL)` (`:95-101`), `StdoutPipe`, `Stderr=Stdout`, `LimitedReader` cap+1, `Wait`; timeout mapped at `:121-122`; kill runs exactly once, nothing re-sweeps after reap |
| W3 | Invoke path = broker work + Execute + store writes; seams test-owned; store is SQLite | Read `broker.go:163-307`, `handlers_git.go:38-69`, `broker_test.go:15-27`, `store/store.go:4-5`; grep `handlerSession\|handlerRecordingStore` `handlers_test.go` | live path: decide→hash→`PutObject`→`AppendNextEffectIntent`→`Execute`(`MkdirTemp`+`runBounded`+`RemoveAll`)→`putRecord`+`AppendEffectOutcome`; `handlerSession` injects any `Handler` and wraps the store in `handlerRecordingStore{base: openTestStore(t)}`; store = `modernc.org/sqlite` (pure Go), `:memory:` in tests |
| W4 | base rate consistent with V-B; run-assertion discipline (re-derives V-G) | `GOTOOLCHAIN=go1.25.6 go test -c -race -o /tmp/w16-broker-race.test ./host/broker`; 20× isolated `-test.run '^TestHandlerTimeoutKillsTheWholeProcessGroup$' -test.count=1 -test.v`, counting `--- PASS:` lines; control: bogus `-test.run '^TestNoSuchTestName$'` | `pass=20 fail=0 neither=0` (each PASS line required, not rc); control printed `testing: warning: no tests to run` + `PASS` rc=0 — V-G's vacuous green, reproduced, hence §4's rule |
| W5 | probe faithfulness verdict (the brief's direct question) | Read `/tmp/pgprobe/main.go` against `handlers.go:82-131` | mechanics faithful (same 7 elements); divergences: **one fixture inode reused across all runs** (script path is `os.Args[1]` — material, §3), no `cmd.Dir`/`cmd.Env` (real path sets both), 1MB cap vs 8MB+1 (immaterial), plain binary not go-test harness |
| W6 | `set -m` escape leaks warm but not cold | escape fixture `set -eu; set -m; sleep 5 &; wait` via per-run probe: 20 runs warm; then 8 runs, each on a **freshly written** fixture file | warm: **20/20 leak**, elapsed 5.099–5.179s, killErr nil each; cold: **1/8 leak** (5.239s), 7 runs end 101–103ms |
| W7 | child-only kill (MUT-KILL-NEUTER analogue): warm deterministic, cold not | probe variant with `syscall.Kill(cmd.Process.Pid, …)` (positive-control diff: one line vs W5 source): 5 runs warm fixture; 8 runs fresh fixtures | warm: **5/5 leak** 5.137–5.172s; cold: **2/8 leak** (5.257s, 5.17s), 6 runs end 101–103ms — in ¾ of cold runs nothing existed to survive the narrowed kill |
| W8 | fresh-script first-exec latency ≈ the test's deadline (darwin) | python: 8 fresh trivial `#!/bin/sh` scripts exec'd once each; then same file 5 more times | cold: 223.4, 101.9, 102.8, 102.3, 103.4, 102.6, 101.1, 102.5 ms; warm: 10.2, 7.8, 3.6, 3.0, 2.9 ms |
| W9 | gate legs, watchdog, toolchain denylist, race-control | Read `scripts/verify_go.sh` (whole file) | AILANG_BIN exact-token gate (`v0.30.0`, dirty = rejection); tracked-binary gate; go1.26.0–1.26.5 FATAL denylist; race known-positive control; `go test ./... -count=1` then `-race -timeout 8m` under a 600s python killpg watchdog; comment records host/broker = 76.9s race-leg critical path |
| W10 | V-H re-derived: the AILANG_BIN red is a `t.Fatal`, and CI pins it | grep/sed `episode_test.go:166-170`, `handlers_test.go:200-209`, `.github/workflows/ci.yml` | `TestEpisodeLiveReplayThreeArmsAndEvidence` `t.Fatal`s on unset `AILANG_BIN` (`:169`); contrast: Model.Infer helper `t.Skip`s (`:205`) with a comment naming CI as the alarm; ci.yml exports `AILANG_BIN` in go-verify (`:144`) and only `WORLD_PKG_AILANG_BIN` in job 1 (`:93`) |
| W11 | conflict-surface census; proposed file unallocated | grep `wantFileCount` host/ (hits only `host/boundary/…:1163-1165`); grep `-n 'len(files)\|_test.go'` `invoke_boundary_test.go`; sed `registry_publish_test.go:1075-1110`; grep `LEG1_MODULES scripts/verify_ail.sh`; `ls host/broker/` | AST walker skips `_test.go` (`:145`,`:218`), floors ≥30 (`:212`), outside-broker ban (`:270`); `enumerateSubprocessSites` skips `_test.go` (`:1084`) but walks production `host/**` — M2 warning recorded; `LEG1_MODULES` at `verify_ail.sh:135` (untouched); no `handlers_stall_diag_test.go` exists (control: 26+ files listed, `handlers_test.go` present) |
| W12 | no `t.Parallel` in the file (negative-existence with in-file control) | `grep -c 't.Parallel' host/broker/handlers_test.go`; control `grep -c 't.Helper'` same file | 0 hits; control = 7 — the instrument reads the file |
| W13 | single-run cost (perturbation budget baseline) | `time` one isolated `-race` run of the test | `--- PASS: … (0.12s)`; process total 1.29s |
| W14 | siblings capture their errors; same bound shape, no fork | Read `handlers_test.go:286-326` | `TestGitHandlerTimeoutWritesFailureRecord` / `TestModelHandlerTimeoutWritesFailureRecord`: `sleep 5` direct child, 40ms bound, 2s assertion, error asserted via `assertHandlerFailureRecord(..., ErrHandlerTimeout)` |
| W15 | warm-up guard works and is cheap | 3 trials: fresh guarded fixture, exec with `W16_WARM=1`, then re-exec | first exec 269.2 / 106.1 / 104.0 ms (assessment paid); second 12.4 / 12.2 / 12.2 ms |
| W16 | full proposed fixture verified end-to-end under both kill shapes | proposed body (guard + baked marker dir + `survived` writer), warmed, through the faithful probe and the child-only variant; `ls` markers after each | group kill: 103ms, markers `exec_started forked`, no `survived`; child-only kill: **5.155s**, markers `exec_started forked survived` — the subshell grandchild preserves both arms' behaviour. *(Round 1 read this as a never-killed-vs-killed-late discriminator; the quorum showed the marker cannot carry that — §5.3 as revised. The measurement stands; the round-1 inference from it is withdrawn.)* |
| N1 | `cmd.Cancel`'s kill is the sole group-kill site in broker production code — one seam covers the boundary | `grep -n 'syscall\.Kill' host/broker/*.go` (non-test); `grep -n 'Process\.Kill'` same scope | exactly one hit, `handlers.go:100` `syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)` (the pattern demonstrably reads — it returned the known site); the only other kill is `handlers.go:116` `cmd.Process.Kill()`, the overflow path's direct-child kill, out of scope |
| N2 | the binding gate cannot match the seam | Read `invoke_boundary_test.go:247-313` | detector kinds are exactly `invoke-call`, `ctor-live`, `ctor-replay`, `session-type`, `dot-import` (incl. method-value/func-value forms); inside-broker exemption pinned to `wantCount=3`, identities all in `publish_op.go` (`:274-288`). A package-level `var killGroup` + a call to it is none of these kinds |
| N3 | the subprocess gate pins a file set, not lines, and matches only `exec.Command*` | Read `registry_publish_test.go:1076-1131,1181-1190` | matcher: `pkg.Name == "exec"` && sel ∈ {`Command`,`CommandContext`}, skips `_test.go` (`:1084`); driver map pins five files incl. `handlers.go` (`:1181-1186`) with `len(drivers) != len(files)` as the tripwire; site lines are `t.Logf`ed, never asserted |
| N4 | errno capture at the kill boundary costs one assignment per field | Read `/tmp/pgprobe/main.go:31-41` (the controller's 1,987-run instrument) | `cmd.Cancel` records `r.nilProc` on the nil-`Process` early return and `r.killErr = err` around the real `syscall.Kill` — the recorder-behind-the-seam pattern of §5.4, already proven at scale by V-D |
| N5 | the §5.4 seam is gate-clean, behaviour-preserving, and restorable — measured, not argued | one-shot: `cp` backup; apply seam (+5/−1); `shasum` before/after; `go build ./host/broker/`; `go vet ./host/broker/`; `go test -run 'TestRegistryDispatchBindingBoundary\|TestEverySubprocessSiteIsDriven…' -v -count=1` counting `--- PASS`; `go test -race -run '^TestHandlerTimeoutKillsTheWholeProcessGroup$' -count=5 -v`; restore; `shasum` | sha `8419874…` → `666a5d4…` (edit took — positive control), BUILD_OK, VET_OK, gates **2× `--- PASS`** (0.11s, 1.35s), group-kill **5/5 `--- PASS`** (0.12–0.14s each), restore sha `8419874…` byte-identical |
| N6 | nothing pins `handlers.go` by hash, line, or count anywhere in gates/scripts/CI | `grep -rn 'handlers\.go'` over `host/ scripts/ .github/` (`*_test.go`,`*.sh`,`*.yml`), sibling names excluded; control: same sweep style finds the known `wantFileCount` pin in `host/boundary/allowlist_world_test.go:1163` | sole hit is the N3 driver-map key (file-set membership, which the seam does not change); the control pin was found, so the instrument reads |

## Quorum verification log

**Round 1 (iter-78): BLOCKED.** Two-reviewer quorum, both present, no N-1 degrade.

- **`gpt5-6-sol` — BLOCKING, and correct.** The objection: the `survived` marker does not
  distinguish "never killed" (H1) from "kill fired after natural death" (H2-late) — the
  grandchild writes the marker at sleep-completion in both cases, the test-owned timestamps
  bound only the whole `Execute` window, and nothing establishes when or whether `syscall.Kill`
  ran. Yet round 1's §5.3 asserted the stronger "never killed" reading and §6 authorized M2's
  re-sweep from `survived` + `*HandlerTimeoutError` — evidence equally consistent with a
  ~5s deadline-machinery stall that a post-reap re-sweep would do nothing for. Checked against
  the round-1 fixture (`sleep 5 && : > "<markdir>/survived" &` — marker creation is coupled to
  the sleeper's completion, nothing else) and **conceded without reservation**: the defect was
  real and load-bearing, since that one bit carried the whole M2 decision gate. No refutation
  was attempted because the fixture's own source settles it.
- **Resolution taken: (A) — bring M2's `killGroup` seam forward into M1** (§5.4), so the test
  records the kill's invocation count, monotonic offset, target pgid, and errno directly at
  the boundary — the objection's own `proposed_fix` evidence list. The alternative (B), gpt5's
  fallback of narrowing M1 to localisation and deferring mechanism selection, was rejected
  because it converts one wait on a 0.76% event into two: a trap that fires once per ~130 CI
  runs and then cannot name the mechanism wastes its firing. Price paid and recorded: M1 is no
  longer zero-production-bytes (+5/−1 behaviour-identical lines in `handlers.go`). The seam
  was proven gate-clean by a sha-verified one-shot on the real tree — build, vet, both
  `host/broker` AST boundary gates, and the group-kill test 5/5 under `-race`, all green,
  restore byte-identical (N1–N3, N5–N6). §6's decision condition now requires the recorded
  kill signature (one invocation, errno nil, offset ≈ deadline) *plus* `survived` plus the
  natural-death-band elapsed; §7 P1 was extended so the forced child-only-kill one-shot must
  reproduce that full H1 signature, proving the gate reachable and truthfully reported; §8
  gained MUT-SEAM-BYPASS so reverting to the inline kill cannot pass silently.
- **`gemini-3-1-pro` — PASS**, one non-blocking fix: make §5.1's pass path explicitly assert
  `err` is a `*HandlerTimeoutError` (preventing silent non-timeout early returns) and print
  `%+v` on failure. **Adopted verbatim** in §5.1.

## Related

- Queue item 16 — `design_docs/world-mission.md` (the row this doc corrects and re-scopes)
- Iteration-78 controller measurements V-A..V-H (first-party probe `/tmp/pgprobe/main.go`)
- `design_docs/planned/w-ail-gate-module-pin.md` — the register model, and iter-76's precedent
  for overriding a queue row's own acceptance prescription
- `design_docs/coding-standards.md` S6 (honest gates) — the rule the `-count=20` prescription
  and the darwin per-run vacuity both fail
- `handlers.go:79-131` `runBounded`; `handlers_test.go:729` the test; the `killGroup` seam
  (§5.4, one-shot-proven in N5); M2's residual scope — the re-sweep (§6)
