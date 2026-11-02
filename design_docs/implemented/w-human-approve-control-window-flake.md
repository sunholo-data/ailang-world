# w-human-approve-control-window-flake: AC8's probe reads a slow second dispatch as "max = 1"

- **Status**: Implemented (iteration 334, PR pending). Was: Planned. Designed iteration 334, unattended. **r3** (controller narrow-refinement
  carve-out after quorum r2 BLOCKED, reviewers' verbatim fixes, §10). **r2**: r1 after quorum r1 BLOCKED
  (§10). The `late-second` stagger is now anchored on the first probe entry. The stagger checks are
  labelled empirical. CONTROL's probe wiring is stated. The full gate has been run (V16).
- **Item**: queue row 235, a test defect in row 176's AC8 (`8511badd`, 2026-10-14, PR #261). The
  production serialization is correct. The fix is test-only: one file.
- **As landed (iteration 334, judge NB3)**: the landed diff is the r2 prototype (+126/−17, sha
  `bb75420d0713`) plus §2.5's timeout seam and `TestConcurrencyProbeNamesAMissedArrival`: **+177/−17**,
  file sha `877013b8…` (plan diff `8dee6916…`). Line numbers in §4 and V17 are the prototype's. On the
  landed file the killers measured by the executor and the judge are `:372` (`unheld`), `:375`
  (serialized `max`), `:392` (M3, with the mutant's +4 lines) and `:422`/`:426` (the handshake test).
- **Estimated**: ~0.1 day. About +126/−17 LOC in `host/daemon/decisionpacket_daemon_test.go`, 0
  production. Prototype diff: V10.
- **Measurement base**: `388001c2` (worktree `sprint/world-iter334-human-approve-control-window`, V1).
- **Instruments**: `go1.26.6 darwin/arm64`, 16 cores; `AILANG_BIN=$HOME/.pinned-ailang/ailang`
  (v0.41.0) for Go runs. Everything ran in throwaway `git archive HEAD` copies, never in the
  worktree: base harness `/tmp/world-iter334/designer/base`, r2 prototype
  `/tmp/world-iter334/designer-r2/proto2`, full gate `/tmp/world-iter334/designer-r2/gate`.
  Evidence is in `/tmp/world-iter334/designer-r2/`: `runmut.sh` (applies one substitution, runs,
  restores, and prints the sha256 before, edited and restored), `r2run.sh`/`r2run.txt`,
  `proto-r2.diff`, plus the r1 `soak.sh`/`soak.txt`/`proto.diff`.

---

## 1. The defect

`TestHumanApproveIsSerializedAcrossEpisodes` (`host/daemon/decisionpacket_daemon_test.go:226`)
dispatches two `Human.Approve` invocations on two episodes at once, twice. The first pass wraps a
`concurrencyProbe` in `serializedHandler` and wants `max == 1`. The second pass (CONTROL) uses the
bare probe and wants `max == 2`. The probe's `Execute` waits for a second entrant like this (`:214-218`):

```go
select {
case <-p.second:
case <-time.After(200 * time.Millisecond):
case <-ctx.Done():
}
```

**Two holes come from that 200 ms.**

- **H1, the CI red (the row).** CONTROL needs both dispatches inside `Execute` together. If the
  second reaches `Execute` more than 200 ms after the first, the first has already left, and CONTROL
  reads `max = 1`. CI run 40980017500 attempt 1, go race gate (V3):
  `--- FAIL: TestHumanApproveIsSerializedAcrossEpisodes (1.72s)` …
  `decisionpacket_daemon_test.go:267: unwrapped probe max = 1, want 2`. Attempt 2 was green.
  - **The mechanism reproduces on demand (V7, V9). This is empirical, not a scheduling
    guarantee.** The controller's probe: a 0 ms or 150 ms stagger gives `max=2`; 250 ms or 400 ms
    gives `max=1`. My base harness puts a 300 ms sleep before the second `Dispatch` in today's test.
    Result: 20/20 FAIL plain and 20/20 under `-race`, every one at
    `:270 unwrapped probe max = 1, want 2`. Harness lines are base +3 from the 3 inserted lines, so
    this is base `:267`, the CI line. No load and no Linux are needed. Load widens the pre-Execute path, which is why no load is needed to reproduce the mechanism (V9); load can also deschedule a goroutine past its timer — the R2 residual — which is why AC-1 and AC-3 are labeled empirical, not gates.
- **H2, a vacuity hole the row does not name.** The serialized arm has the same window. The first
  holder leaves `Execute` on its own 200 ms timer. If the second dispatch reaches the mutex only
  after that, `max` is 1 **whether or not the mutex exists**. So on exactly the slow runs that broke
  CONTROL, the serialized arm checks nothing. Measured (V9): with the mutex deleted from
  `serializedHandler` (row 176's mutant MD3) **and** the 300 ms stagger, the serialized arm passes
  20/20 (harness `:263 serialized probe max = 1`, base `:260`, the log line after the assertion).
  Only CONTROL fails. On the CI red itself the serialized arm passed (V3), so that run gave no
  evidence that the mutex works.

## 2. The decision: wait on events, a held-mutex oracle, an anchored replay of the gap

### 2.1 Probe: return on the second arrival, bounded, and say why if not

`concurrencyProbe` gains `entries`, `unheld`, `missed`, `inside2` (renamed from `second`),
`entered`, `arrived` and `held`. Only the **first** entry waits, and that entry closes `entered`.

- **What it waits for.** In CONTROL it waits for `inside2`, which closes when `cur` reaches 2 inside
  the probe. In the serialized arm it waits for `arrived`, which the outer gate closes (§2.2). The
  wait returns as soon as the event fires, so a healthy run pays nothing for the bound.
- **The two arms are wired differently. This is binding on the executor.**
  - **CONTROL** gets its own `concurrencyProbe` with **`held == nil`** and `arrived == nil`.
  - When `held` is nil, the `TryLock` check and the `unheld` counter are off. Nothing holds
    `humanMu` on the unwrapped path, so a non-nil `held` there would fail every CONTROL entry.
  - The `unheld == 0` assertion runs **only in the serialized arm** (measured prototype `:366`, final diff `:367`). CONTROL
    asserts `missed` (`:379`/`:380`) and `max == 2` (`:382`/`:383`) only.
- **Bound: `probeWaitBound = 10 * time.Second`.** `dispatchPair`'s whole-pair bound becomes
  `probeWaitBound + 5*time.Second` (was a literal 5 s, `:249`). Both sit under
  `boundedTestContext`'s 30 s (V5). So the bound that fires first is the probe's own, and its
  message names the cause. The whole test takes 0.43 s idle today (V8) and took 1.72 s on the CI
  red (V3).
- **When the bound expires, the test says so. It never reads as `max = 1`.** The probe records
  `missed`, built from the constant:
  `fmt.Sprintf("second Human.Approve dispatch did not arrive within %s", probeWaitBound)`, or
  "test context ended …". Each arm checks `missed` first:
  `serialized arm measured nothing: <why> (dispatch errors [..])` or
  `CONTROL measured nothing: …`. The dispatch errors are printed because a second dispatch that
  fails before `Execute` is the likeliest real cause. Measured (V11, M5): setting the constant to
  100 ms fails 3/3 with
  `:363: serialized arm measured nothing: second Human.Approve dispatch did not arrive within 100ms (dispatch errors [<nil> <nil>])`.

### 2.2 Serialized arm: an outer arrival gate plus a held-mutex check

The seam is real (V4). `d.binder` copies `d.human` into a fresh registry on every call
(`daemon.go:917-918`). `d.human` is a field "so AC8 can wrap a probe" (`:437-438`). The test already
replaces it with `serializedHandler{mu: sh.mu, h: probe}` (`:231`), where `sh.mu` is
`&d.humanMu` (`daemon.go:807`; wiring asserted at `:180-181`, V15). The fix adds one more layer
**outside** the serialization:

```go
gate := &arrivalGate{second: make(chan struct{})}   // counts arrivals, closes second at 2
probe := &concurrencyProbe{inside2: make(chan struct{}), entered: make(chan struct{}), arrived: gate.second, held: sh.mu}
gate.h = serializedHandler{mu: sh.mu, h: probe}
d.human = gate
```

1. The first dispatch takes the mutex and enters the probe. It waits until the **second has
   reached the gate**. The gate is outside the mutex, so the second can get there while the first
   is still inside.
2. The first then holds `serializedHold = 50 ms` longer, ending early if `inside2` fires, and
   releases.
3. **Each serialized-arm entry checks that `held` is held**: `p.held.TryLock()` must fail. If it
   succeeds, the probe unlocks at once and counts the entry in `unheld`. The arm asserts
   `unheld == 0` **before** `max == 1`:
   `serialized probe ran %d of %d executions without the daemon mutex held`.

**What this establishes, and what it does not.**
- It establishes two facts:
  - The second dispatch **arrived** at the handler boundary while the first was still active
    inside `Execute`.
  - The daemon mutex was **held** for every probe entry (`TryLock` failed each time).
- It does **not** observe the second dispatch calling `Lock`. A `sync.Mutex` exposes no waiters.
  The second may still have been on its way from the gate into `serializedHandler.Execute`
  (`daemon.go:928-931`) when the first released.
- The serialized arm therefore no longer depends on its timer to be meaningful. Today it can pass
  with no second arrival at all (H2). Now it fails by name unless the second has arrived.

**The deterministic mutation oracle is `TryLock`.** Consider any `serializedHandler` that does not
hold `mu` around the inner `Execute`: the mutex deleted (M1), an unlock before the call (M1b), or the
bare probe (M2). In each, `TryLock` succeeds inside the probe on the first entry, whenever the
second arrives. One thing could make `TryLock` fail falsely: another goroutine holding `humanMu`.
The only other code that takes it is the finish route (`episode_finish.go:158`), and this file
never posts there (V6). `max == 1` stays as a second, independent signal (M6/M6b, V11). `TryLock`
already has a test precedent in this repo (`host/verifygate/forklocked_write_test.go:183`, V6).

**No commit before `Execute` can strand the second dispatch.** Before the handler runs, the broker
appends only an episode-scoped effect intent (`broker.go:252-261`), and the two episodes differ.
The world-head commit (`coordinator.go:399`) comes after `Execute`. The first now holds until the
second is at the gate, so nothing commits before the second has loaded its world. The world-head
conflict that follows lands after `Execute`, and `requireOneCommitted` already accepts it
(`:285-286`). Every prototype arm logged that conflict (V10).

### 2.3 Regression arm: replay the CI gap, anchored on the first entry

The test becomes two subtests over one helper, `humanApproveSerializedArms(t, stagger)`. Each
subtest gets a fresh `newHandlerDaemon`, so the episode names do not collide.

| Subtest | Second dispatch starts | Purpose |
|---|---|---|
| `together` | at once | today's case |
| `late-second` | after the first dispatch has **entered the probe** (`<-probe.entered`), then `lateSecond = 300 * time.Millisecond` | the CI gap: ≥250 ms read as `max = 1` (V7) |

**[r2] The anchor.** `dispatchPair` takes the probe's `entered` channel. In `late-second`, the
second goroutine waits for it before it sleeps `lateSecond`. That wait is bounded by
`probeWaitBound` and ctx, and a timeout becomes the second dispatch's error
(`first dispatch did not enter the probe within …`). The gap is therefore a **lower bound measured
from the first probe entry**, not from goroutine start. The scheduler can no longer delay the
first dispatch past the sleep.

**What stays empirical.** The fixed test's outcome does not depend on the gap: it waits on events,
up to 10 s. Two checks do depend on it:
- AC-1's base reproduction;
- AC-3, the kill of mutant M3 (the old 200 ms timer).

Both are **empirical, scheduling-sensitive checks, not deterministic gates**. **[r3]** M3 can
survive a run whenever the first goroutine is late relative to the second dispatch's arrival, in
any of three ways (gpt6-1-sol r2):
- it is descheduled between closing `entered` and creating its 200 ms timer (delayed timer
  creation), so the timer starts late;
- its timer has fired, but it is descheduled for more than 100 ms before it decrements `cur`;
- the timer and the arrival are both ready when it selects, and the arrival wins (Go picks
  uniformly among simultaneously ready cases).

In each case the second dispatch must also get through its whole path into the probe in time.

Measured: 60/60 killed, including 20/20 at a load average of 13.9 (V11).

**[r3] Controllable timeout seam: ADOPTED for the named failure only (gpt6-1-sol r2's verbatim fix,
§2.5).** M3's kill stays empirical. **[r2, kept for M3] Barrier or controllable-timer seam for the
old window: not adopted.** It would make M3's
kill deterministic. That needs a timer seam in the probe and a two-phase handshake, about 30 more
lines. They would test the test's own window, not a production property. The production property is
"one `Human.Approve` at a time", and its deterministic oracle is `TryLock`. The fixed probe has no
window left to test: its only timer is the 10 s failure bound. The anchor is small (one channel and
one bounded wait) and closes the ordering gap the reviewer named. The residual is stated here and in
R2.

Wall time measured idle: today 0.43–0.44 s, r2 prototype 0.78–0.80 s, so **+0.36 s** (V8, V10).
`-run '^TestHumanApproveIsSerializedAcrossEpisodes$'` still selects both subtests.

### 2.5 [r3] Named-failure handshake (a test-only timeout seam)

gpt6-1-sol r2, verbatim: "Add a test-only controllable timeout seam to `awaitSecond`. For the
named-failure acceptance test, keep the arrival channel blocked, explicitly trigger the timeout,
and wait for acknowledgement that `missed` was recorded before allowing dispatch completion. Retain
the real 10 s bound in normal runs. Replace M5's deterministic claim with evidence from this
handshake-driven test."

Shape (the planner fixes the names):
- `concurrencyProbe` gains a test-only `timeout <-chan time.Time` field. `awaitSecond` selects on it
  when it is non-nil and on `time.After(probeWaitBound)` otherwise, so every existing arm keeps the
  real 10 s bound (`timeout` is nil there).
- A new deterministic test (for example `TestConcurrencyProbeNamesAMissedArrival`) drives
  `concurrencyProbe.Execute` directly on one goroutine. The arrival channel is never closed. The
  test sends on `timeout` itself, then waits (bounded) for `Execute` to return, which happens only
  after `missed` is recorded. It asserts `missed` equals
  `fmt.Sprintf("second Human.Approve dispatch did not arrive within %s", probeWaitBound)`, and that
  the probe recorded no second entry.
- Mutants it must kill, each deterministically: (i) `awaitSecond` stops recording `missed` on the
  timeout branch; (ii) the timeout branch's text drops the constant; (iii) `awaitSecond` ignores the
  `timeout` field (the test then blocks until its own bound and fails by name, never hangs).

### 2.4 Alternatives

- **(a) A longer fixed window (for example 2 s) in both arms.** It narrows H1 but does not close it.
  The serialized arm would cost 2 s per run, because that timer releases it. H2 stays. Rejected.
- **(b) The row's candidate: separate budgets, with the serialized arm keeping its 200 ms timer.**
  It fixes H1 and leaves H2: V9's MD3 + stagger survives the serialized arm. Rejected.
- **(c) An instrumented `sync.Locker` in `serializedHandler`, so the test can see waiters.** That
  would show that the second called `Lock`. It needs a production type change (`mu *sync.Mutex` to
  an interface). `TryLock` already kills the mutex mutants deterministically without it. Rejected.
  S2 holds: no production code changes.

## 3. Risk

- **Test-only.** Every non-test file is untouched. The wiring check `sh.mu != &d.humanMu` stays in
  `TestDaemonInvocationMintsDecisionPacket` (`:180-181`, V15). This row does not touch it.
- **Any regression that stops the second dispatch before `Execute` now fails by name within 10 s.**
  Today such a run would pass the serialized arm on its timer. That is the intent.
- **When `TryLock` fails on a held mutex, nothing changes.** When it succeeds, the probe unlocks at
  once, before anything else runs in that entry.
- **`time.Sleep(lateSecond)` is the injected gap, not a synchronisation.** The ordering comes from
  the `entered` anchor.

## 4. Mutation table

Each mutant was applied with `runmut.sh` and the file was restored byte-identical (V9, V11). The
**killer** is the assertion that fires. Lines are as measured on the r2 prototype
(`9862b1662185`):
- serialized arm: `:363` `missed`, `:366` `unheld`, `:369` `max`;
- CONTROL: `:379` `missed`, `:382` `max`.

The final diff (`bb75420d0713`, V10) differs only in the `arrivalGate` comment, which is one line
longer (§2.2 wording). So each line is one higher there: `:364/:367/:370/:380/:383`. **[r3] Measured
on the final sha, not derived (controller, V17):** M1 and M7 fire at `:367`; M3 fires at `:387`
(`:383` plus the mutant's +4 lines).

| ID | Mutant | Base `388001c2` (today's test) | r2 prototype (plain ×20 / `-race` ×20) | Sole killer | Kind |
|---|---|---|---|---|---|
| A | 300 ms sleep before the second `Dispatch`, today's probe | 20/20 FAIL plain and race (`:270 unwrapped probe max = 1, want 2`) | n/a: replaced by `late-second` | CONTROL | empirical |
| A2 | A + MD3 | serialized arm **passes** 20/20; only `:270` fails | n/a | (shows H2) | empirical |
| M1 | MD3: delete `s.mu.Lock()` + `defer s.mu.Unlock()` (`daemon.go:929-930`) | killed idle 20/20 (`:258 max = 2`), survives the serialized arm under A's gap (A2) | 20/20, 20/20; both subtests | `:366` `unheld` ("ran 2 of 2 executions without the daemon mutex held") | **deterministic** |
| M1b | unlock before the inner call | not measured | race 20/20; both subtests | `:366` | **deterministic** |
| M2 | serialized arm uses the bare probe (`gate.h = probe`) | n/a | 20/20, 20/20; both subtests | `:366` | **deterministic** |
| M3 | old 200 ms timer (`awaitSecond` → `select{inside2; After(200ms); ctx}`) | (this *is* today's probe) | 20/20, 20/20, and race 20/20 under 32 burners; `late-second` only | `late-second` `:382` (`:386` with the mutant's +4 lines) | **empirical** (§2.3) |
| M5 | `probeWaitBound` = 100 ms | n/a | plain 3/3, `late-second` | `:363` "… did not arrive within 100ms" | **[r3] empirical** (scheduling-dependent: the arrival can win; gpt6-1-sol r2). The deterministic named-failure evidence is §2.5's handshake test (AC-4) |
| M7 | **[r3]** `serializedHandler.Execute` locks a package-level `otherMu` instead of `s.mu` (`daemon.go:929-930`) | n/a | **controller, final sha `bb75420d0713`, race ×20: rc=1, 20/20 top-level FAIL, 40 subtest FAILs, all `:367 serialized probe ran 2 of 2 executions without the daemon mutex held`** (V17) | `:367` `unheld` (not `max`) | **deterministic** |
| M7b | **[r3]** test-side rewiring `serializedHandler{mu: otherMu, h: probe}` with `held: otherMu` | n/a | **controller, final sha, race ×20: rc=0, 20/20 PASS (40 subtests)** — survives every assertion in this file (V17) | none in this file: the production binding is pinned by `:180-181` (V15) | boundary (R3) |
| M6 | M1 + `TryLock` check disabled | n/a | race 20/20; both subtests | `:369` (`serialized probe max = 2, want 1`), the secondary signal | empirical |
| M6b | M6 + `serializedHold = 0` (r1 measurement) | n/a | race 20/20 (r1 prototype) | serialized `max` | empirical |

## 5. Acceptance criteria

Run in `host/daemon` unless noted. Every command is bounded (`-timeout`).

- **AC-1 (base reproduction; EMPIRICAL, scheduling-sensitive).** On a copy of the base tree,
  insert `if i == 1 { time.Sleep(300 * time.Millisecond) }` before `_, errs[i] = c.Dispatch`
  (`:242`). Run
  `go test [-race] -timeout 600s ./host/daemon -run '^TestHumanApproveIsSerializedAcrossEpisodes$' -count=20 -v`,
  then restore.
  - Want: FAIL with `unwrapped probe max = 1, want 2`. Measured: 20/20 plain and 20/20 race (V9).
  - This shows the mechanism. It is not a gate: the sleep is not anchored on the first entry.
- **AC-2 (mutex removal; DETERMINISTIC).** Apply M1, M1b and M2 one at a time to the fixed tree.
  Use the AC-1 command at `-count=20`, plain and `-race`.
  - Want: 20/20 FAIL, both subtests, killer
    `serialized probe ran 2 of 2 executions without the daemon mutex held`.
  - Measured (V11). Sole killer: the `unheld` assertion.
- **AC-3 (old window; EMPIRICAL, scheduling-sensitive).** Apply M3, same command.
  - Want: FAIL in `late-second` only, killer `unwrapped probe max = 1, want 2`; `together` passes.
  - Measured: 20/20 plain, 20/20 race, and 20/20 race under 32 burners (V11).
  - Residual: a run can survive when the first goroutine is late (§2.3, three ways).
  - **[r3] Adjudication (oc-kimi-k3 r2, verbatim):** a lone M3 survival whose log shows the second dispatch entering only after the first goroutine resumed (>100 ms deschedule past its expired timer) is the R2 residual, not a regression — re-run once at -count=20; a second survival fails the row. AC-1 is a demonstration on the base tree, not a gate on the fix.
- **AC-4 (named failure; [r3] DETERMINISTIC via §2.5's handshake).** Run the §2.5 test
  `-count=200` plain and `-race`. Want: 200/200 PASS. Apply §2.5's mutants (i)–(iii) one at a time:
  each must fail every run by name, and (iii) must fail on the test's own bound, never hang. M5 (bound
  set to 100 ms in the `late-second` arm) stays as an empirical demonstration only (measured 3/3, V11).
- **AC-5 (soak).** Run
  `go test -race -timeout 900s ./host/daemon -run '^TestHumanApproveIsSerializedAcrossEpisodes$' -count=200`,
  then the same with 32 `yes > /dev/null` burners running.
  - Want: `ok`, 0 FAIL. Measured: 200/200 idle and 200/200 at a load average up to 35 (V12).
- **AC-6 (package and full gate).**
  - `go vet ./host/daemon/`, `go test ./host/daemon/ -count=1` and the same with `-race` (V12).
  - `go build ./... && go test -timeout 1800s ./...` and `./scripts/verify_ail.sh` with the
    released binaries. Measured on base + the final diff: full suite rc=0 (315 s, 35 `ok`, 0 FAIL)
    and `verify_ail.sh` rc=0 (11 s) (V16).
  - CI on the fix head: all 3 jobs green, with `ok  github.com/sunholo-data/ailang-world/host/daemon`
    in the race gate.
- **AC-7 (wall time).** One idle `-v` run shows the whole test under 1.2 s (r2 prototype 0.78–0.80 s,
  V10).

## 6. Files touched, conflict surface

| File | Change |
|---|---|
| `host/daemon/decisionpacket_daemon_test.go` | **Imports:** `fmt` is added. **Constants:** `probeWaitBound`, `serializedHold`, `lateSecond`. **Types:** `arrivalGate`. **`concurrencyProbe`:** gains `entries`, `unheld`, `missed`, `inside2` (renamed from `second`), `entered`, `arrived`, `held` and `awaitSecond`. **Field `held`: set only in the serialized arm (`held: sh.mu`). CONTROL's probe is a separate instance with `held == nil` and `arrived == nil`, so its `TryLock` check and `unheld` counter are off. `unheld == 0` is asserted only in the serialized arm, and CONTROL asserts `missed` and `max` only.** **Test:** 2 subtests over `humanApproveSerializedArms(t, stagger)`. **`dispatchPair`:** takes `entered`, anchors the stagger on it, and bounds the pair at `probeWaitBound + 5*time.Second`. `requireOneCommitted` is unchanged |

Conflicts:
- `concurrencyProbe` has no other users (4 matches, all in this file, V4).
- No charter row names this file (V14, 0 matches). The only row that names the test is 235 itself.
- The only branches with commits on the file that are not in base are the pre-squash branches of
  rows 172 and 176. Their work landed as the squash merges `d36b148d` and `8511badd` (V14).
- `w-marketing-capture.md` cites `:24,40,69` of this file read-only. The added `fmt` import shifts
  those lines by +1.

**S1/S3/S4/S5/S7:** n/a (test-only, host, no `.ail`, no usage surface). **S6:**
- M1, M1b and M2 show the serialized arm can no longer pass without the mutex.
- M3 shows (empirically) that `late-second` fails on the old window.
- §2.5's handshake test shows that a missing arrival fails by name (AC-4); M5 is an empirical demonstration.

## 7. Scope

- **Linux measurement (iteration 318's throwaway branch): not needed.** The fixed test's verdict
  depends on events, not on how long the path before `Execute` takes, up to the 10 s bound.
  `late-second` replays the CI gap on every platform. CI on the fix head (AC-6) is the Linux check.
- **Out of scope: row 200's sibling flake.** `TestWorktreeSnapshotRefusesAConfiguredFilter`
  (`host/broker`) is a different test with a different mechanism (V13). It stays on row 200.
- **Sweep.** `time.After(200 * time.Millisecond)` appears once in `host/daemon` tests, in this probe
  (V4).

## 8. Residuals

- **R1. The bound is still a number.** If the second dispatch took more than 10 s to reach `Execute`,
  the test would fail, by name. The whole test took 1.72 s on CI and at most 1.00 s at a load
  average of 35 (V12).
- **R2. [r2] M3's kill (AC-3) and AC-1 are empirical.** The anchor makes the gap a lower bound from
  the first entry. One nondeterminism remains: the first goroutine could be descheduled for more
  than 100 ms after its 200 ms timer fires, while the second gets all the way into the probe. Then
  M3 would survive one run. Measured 60/60 killed (V11). The fixed test does not carry this
  residual. The deterministic mutation oracle for the production property is `TryLock` (AC-2).
- **R3. `TryLock` shows that the mutex is held, not who holds it, and not that the second called
  `Lock`.** **[r3, oc-glm-5-3 r2 verbatim]** TryLock shows only that the mutex passed as held is held around Execute. M7 is caught by unheld, not by max. A consistent rewiring of both the wrapper mutex and held (M7b) is caught by no assertion in this file — max == 1 survives it because the substitute still serializes. The binding of serializedHandler to &d.humanMu in production New is pinned by :180-181 (V15), and the type's s.mu.Lock around Execute is mutex-parametric, so the composition still covers the production property; no claim is made that max catches wrong-mutex mutants. (M7, M7b measured: V17.)

## 9. Verification Log

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Base and tools | `git -C <wt> log --oneline -1`; `git status --porcelain`; `go version`; `grep '^go ' go.mod`; `~/.pinned-ailang/ailang --version`; `sysctl -n hw.ncpu` | `388001c2 Merge pull request #456 …`; clean; `go1.26.6 darwin/arm64`; `go 1.26.6`; `AILANG v0.41.0`; `16` |
| V2 | Origin of the test | `git blame -L199,271 --porcelain host/daemon/decisionpacket_daemon_test.go \| grep -E '^[0-9a-f]{40}' \| cut -c1-8 \| sort \| uniq -c`; `git log -1 8511badd` | `73 8511badd`; `2026-10-14 feat(row176): daemon binds Human.Approve, serialized daemon-wide (D-WORLD-73 = A) (#261)` |
| V3 | The CI red | `gh api repos/…/actions/runs/40980017500`; `…/attempts/1/jobs`; `gh api --allow-escape-sequences …/jobs/123817474361/logs` (saved `/tmp/world-iter334/designer/job-race-a1.log`) | head `e15fc11f`, run now `attempt=2 success`; attempt 1: `go race gate … failure`, other 2 `success`; log `:424 --- FAIL: TestHumanApproveIsSerializedAcrossEpisodes (1.72s)`, `:425 …:260: serialized probe max = 1`, `:427 …:267: unwrapped probe max = 1, want 2`, `FAIL host/daemon 379.848s` |
| V4 | Probe, seam and wiring as cited | `sed -n 199,271p host/daemon/decisionpacket_daemon_test.go`; `sed -n 434,439p;807p;915,931p host/daemon/daemon.go`; `grep -rn 'serializedHandler{' host/`; `grep -rn 'time.After(200' host/daemon/*_test.go \| wc -l`; `grep -rn concurrencyProbe host cmd --include='*.go' \| wc -l` | `:216 case <-time.After(200 * time.Millisecond):`; `:231 d.human = serializedHandler{mu: sh.mu, h: probe}`; `:242 _, errs[i] = c.Dispatch(ctx, calls[i])`; `:249 time.After(5 * time.Second)`; `daemon.go:807 d.human = serializedHandler{mu: &d.humanMu, …}`; `:918 reg[broker.EffectHumanApprove] = d.human`; `:929-930 s.mu.Lock() / defer s.mu.Unlock() // MD3 target`; `serializedHandler{` 2 sites; `time.After(200` count `1`; `concurrencyProbe` `4`, all in this file |
| V5 | Test context bound | `sed -n 10,15p host/daemon/m7b_deadline_test.go` | `context.WithTimeout(context.Background(), 30*time.Second)` |
| V6 | `humanMu` lockers; `TryLock` precedent; writes before `Execute` | `grep -rn 'humanMu.Lock\|MD3 target' host --include='*.go'`; `grep -n finish host/daemon/decisionpacket_daemon_test.go` (control: `grep -c finish host/daemon/episode_finish_test.go`); `grep -rln TryLock host cmd`; `sed -n 240,261p host/broker/broker.go`; `grep -n 'Store.Commit' host/coordinator/coordinator.go` | lockers: `episode_finish.go:158` and `daemon.go:929` only; `finish` here only at `:26` (a plan JSON key) and `:250` (a message), with no route call (control `42`); `TryLock`: `host/verifygate/forklocked_write_test.go` (`:183`); `:252/:255 Append…EffectIntent(ctx, s.episodeID, …)` then `:261 handler.Execute`; `coordinator.go:399 c.cfg.Store.Commit` in `appendAndCommit` |
| V7 | Controller's mechanism measurement (cited, not re-run) | `/tmp/world-iter334/mechanism-measure.txt`; probe `/tmp/world-iter334/zz_iter334_stagger_probe_test.go.txt` | today's probe: `0s max=2`, `150ms max=2`, `250ms max=1`, `400ms max=1`; probe waiting up to 10 s or ctx: `0s`, `400ms`, `2s` all `max=2` |
| V8 | Base wall time | `(cd base && go test ./host/daemon -run '^TestHumanApproveIsSerializedAcrossEpisodes$' -count=5 -v \| grep '^--- ')` | 5 PASS, 0.43–0.44 s |
| V9 | Base harness: H1, H2 (r1, base unchanged) | `runmut.sh base <test> A-base-stagger 20 {0,1} '<sleep 300ms for i==1>'`; MD3 applied to `base/…/daemon.go` (sha `e6b70d7e3f3f→3737b1ed8187`), unstaggered ×20, then `A2` ×20, then restored (`e6b70d7e3f3f`) | A: `rc=1 PASS=0 FAIL=20` plain (21.4 s) and race (23.1 s), all `:270 unwrapped probe max = 1, want 2`; sha `53a708a5264f→2c4d47bc3e89→53a708a5264f`. MD3 unstaggered: 20 FAIL `:258 serialized probe max = 2, want 1`. A2: 20 FAIL, every one at `:270`, with `:263 serialized probe max = 1` logged, so the serialized arm passed 20/20 |
| V10 | **[r2]** Prototype | `proto2` = fresh archive + r2 span; `gofmt -l` (empty); `go vet ./host/daemon/`; `diff -u base proto2` (`proto-r2.diff`); `-count=5 -v` | measured version: `+125/−17`, vet ok, test-file sha `9862b1662185`. **Final** (the `arrivalGate` comment reworded to §2.2, one line longer): `+126/−17`, `gofmt -l` empty, vet ok, sha `bb75420d0713`; final ×20 `-race` in `gate2.txt` (V16). Measured version ×5: whole 0.78–0.80 s, `together` 0.09 s, `late-second` 0.69–0.71 s; every arm logs the world-head conflict, `gate arrivals = 2` |
| V11 | **[r2]** Prototype mutants (`r2run.txt`) | `runmut.sh proto2 <file> <id> <n> {0,1} '<sub>'` for M1, M2, M3 (plain and race ×20), M1b (race ×20), M5 (×3), M6 (race ×20, M1 applied by hand, restored); M3 race ×20 under 32 `yes` burners | M1, M2: `rc=1 PASS=0 FAIL=20` plain and race, 40 subtest FAILs each, all `:366 serialized probe ran 2 of 2 executions without the daemon mutex held`; M1b race the same. M3: 20/20 plain (18.0 s), 20/20 race (20.7 s), 20/20 race loaded (loadavg `5.38`→`13.87`), `late-second` only, `:386 unwrapped probe max = 1, want 2`. M5: 3/3 `:363 serialized arm measured nothing: second Human.Approve dispatch did not arrive within 100ms (dispatch errors [<nil> <nil>])`. M6: race 20/20, both subtests, `:369 serialized probe max = 2, want 1`. Every restore matched (`9862b1662185`, `e6b70d7e3f3f`); burners gone after (`pgrep -x yes` 0) |
| V12 | **[r2]** Soak and package (`r2run.txt`) | `go test -race -timeout 900s … -count=200 -v`, idle, then with 32 burners; `go test -race ./host/daemon -count=1`; r1 `soak.txt` for package plain | idle: `rc=0`, 200 PASS / 0 FAIL, `ok … 185.121s`, slowest 0.97 s. Loaded (loadavg `4.75`→`35.25`): `rc=0`, 200 PASS / 0 FAIL, `ok … 187.257s`, slowest 1.00 s. Package race (measured version): `ok … 309.210s`, rc=0, 311 s. Package plain (r1 prototype): `ok … 253.811s` |
| V13 | Row 200's sibling, different test | `grep -n 'FAIL' /tmp/world-iter334/job455.log`; `grep -n '^200\. ' design_docs/world-mission.md` | `:1036 --- FAIL: TestWorktreeSnapshotRefusesAConfiguredFilter (0.02s)`, `FAIL host/broker 211.934s`; row `200. … w-filter-refusal-control-flake` |
| V14 | **[r2]** Who else touches the file | `grep -n 'decisionpacket_daemon_test' design_docs/world-mission.md`; control `grep -n 'TestHumanApproveIsSerializedAcrossEpisodes' design_docs/world-mission.md`; `grep -rln 'decisionpacket_daemon_test' design_docs/`; `git log --all --oneline ^388001c2 -- host/daemon/decisionpacket_daemon_test.go`; `git branch -a --contains <sha>` | charter: **0** matches; control: 1, `6307:235.` only. Docs: row 172/176 verification and plan files, `w-marketing-capture.md` (`:24,40,69` cited read-only) and this doc. Branch commits not in base: `618e27cb` (row 176, `sprint/row176-daemon-human-approve-binding`) and `167f3e3c` (row 172, `design/w-decision-packet-production-projection`), both pre-squash branches; their work landed as squash merges (`8511badd` #261 and `d36b148d` #249 in `git log -- host/daemon/decisionpacket_daemon_test.go`) |
| V15 | **[r2]** The mutex wiring check | `sed -n 178,182p host/daemon/decisionpacket_daemon_test.go` | `sh, ok := d.human.(serializedHandler)` / `if !ok \|\| sh.mu != &d.humanMu {` / `t.Fatalf("Human.Approve handler = %T, want serializedHandler with daemon mutex", d.human)` |
| V16 | **[r2]** Full gate on base + `proto-r2.diff` | `git archive HEAD` into `designer-r2/gate`; `patch -p0 < proto-r2.diff` (sha `bb75420d0713`, the final diff); `git init` + one commit in the copy, because `host/pins` runs `git ls-files`; `AILANG_BIN=$HOME/.pinned-ailang/ailang`; `go build ./... && go test -timeout 1800s ./...`; then, with AILANG_BIN unset, `PATH=~/.ailang/state/ailang-release-v0.55.0/x:$PATH`, `WORLD_PKG_AILANG_BIN=~/.ailang/state/ailang-release-v0.41.0/x/ailang`, `./scripts/verify_ail.sh` (`fullgate.out`, `verify_ail.out`) | **First attempt** (no `.git` in the copy; `fullgate.out`): rc=1 in 337 s, 34 ok, only `host/pins TestReleaseLiteralsExplainedByManifest: scan_test.go:256: git ls-files: exit status 128`, which is the missing repo, not the diff. **Git-initialised copy** (`gate2.txt`, `fullgate2.out`): `go build ./... && go test ./...` **rc=0, 315 s**, 35 `ok`, 0 FAIL, 2 no-test; `ok host/daemon 308.931s`, `ok host/pins 0.276s`. `verify_ail.sh` **rc=0, 11 s**: AILANG_BIN unset, `ailang` on PATH = v0.55.0, `✓ compiler pinned by exact bytes: AILANG v0.41.0`, `✓ world package gate PASSED: 9/9`, `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass`; it also prints `ℹ UNRECOGNISED RELEASE: legs 1-2 resolved v0.55.0` (expected set v0.41.0/v0.33.0/v0.30.0), an informational line from the prescribed PATH. Final file ×20 `-race`: 20 PASS, `ok … 19.856s` |
| V17 | **[r3, controller]** M7/M7b, and M1/M3 on the final sha (`/tmp/world-iter334/controller-m7.txt`, banked as `verification/world-iter334/controller-m7.txt`) | In `designer-r2/gate` (sha test `bb75420d0713`, daemon `e6b70d7e3f3f`): M7 = `daemon.go` `s.mu.Lock()/defer s.mu.Unlock()` → `otherMu.Lock()/defer otherMu.Unlock()` + `var otherMu sync.Mutex`; M7b = test `held: otherMu` and `serializedHandler{mu: otherMu, …}` (+ `_ = sh`); each `go test -race ./host/daemon -run '^TestHumanApproveIsSerializedAcrossEpisodes$' -count=20 [-v]`; M1/M3 via `runmut.sh gate … 20 1` | M7: rc=1, 40 subtest FAILs, all `:367 serialized probe ran 2 of 2 executions without the daemon mutex held`; daemon restored `e6b70d7e3f3f`. M7b attempt 1: build failure `declared and not used: sh` (instrument failure, not a verdict); attempt 2: rc=0, 20 top-level PASS, 40 subtest PASS, `ok … 20.088s`; test restored `bb75420d0713`. M1-final: 20/20 FAIL, 40 × `:367`. M3-final: 20/20 FAIL, `late-second` only, 20 × `:387 unwrapped probe max = 1, want 2` |
| V18 | **[r3]** `requireOneCommitted` accepts the world-head conflict (base citation `:285-286`) | `sed -n 278,290p host/daemon/decisionpacket_daemon_test.go` (worktree, base `388001c2`) | `func requireOneCommitted(t *testing.T, errs [2]error) {` … `case store.IsConflict(err):` / `t.Logf("world-head conflict (expected for one of two concurrent commits): %v", err)` / `default: t.Error(err)` |
| V19 | **[r3]** Scope of `boundedTestContext` in the fixed file (oc-kimi-k3 r2) | `grep -n boundedTestContext` and `sed -n 322,326p` on `designer-r2/gate/host/daemon/decisionpacket_daemon_test.go` (`bb75420d0713`) | `:325 ctx := boundedTestContext(t)` inside `dispatchPair`: every dispatch pair gets a fresh 30 s context, so the probe's 10 s bound precedes the context per pair; no stacking across the 4 pairs |

## 10. Quorum log

**Round 1** (artifact `~/.ailang/state/mission-quorum-world/w-human-approve-control-window-flake-2026-11-02T12-26-15Z.json`,
2026-11-02T12:26:15Z, $0.2409): **BLOCKED**. All 3 seats were present:
- `oc-glm-5-3` PASS;
- `oc-kimi-k3` PASS;
- `gpt6-1-sol` REJECT.

The controller's in-session reading was a pass. Each fix below is the reviewer's own.

- **gpt6-1-sol (REJECT, the blocker).** Quoted: "Sleeping 300 ms before the second Dispatch does not
  guarantee that its Execute entry occurs more than 200 ms after the first: the scheduler can delay
  the first dispatch until after that sleep." Also: "The outer arrival gate also establishes
  arrival, not that the second dispatch attempted the mutex." **Correct on both counts.** Applied:
  - (a) `late-second` now starts its stagger after the first probe entry (`entered`, §2.3).
    Re-measured in V10–V12.
  - (b) AC-1, AC-3, M3's row and R2 are labelled empirical, and the residual is stated. AC-2
    (`TryLock`) is the deterministic oracle.
  - (c) §2.2 now says the gate establishes arrival while the first is active, plus the mutex held
    throughout, and that it does not observe the second calling `Lock`. Every "proves contention"
    claim is removed.
  - The barrier or controllable-timer seam is not adopted. The reason is in §2.3: it would test the
    test's own window, not a production property.
- **oc-glm-5-3 (PASS) and oc-kimi-k3 (PASS): CONTROL's `held` was left implicit.** Both noted that
  an executor could wire `held: sh.mu` into CONTROL and red the suite. Applied: §2.1 and §6 now
  state that CONTROL's probe is a separate instance with `held == nil`, that the check and counter
  are off there, and that `unheld == 0` runs only in the serialized arm.
- **oc-glm-5-3: the failure text and two unlogged premises.** Applied:
  - `missed` is built with `fmt.Sprintf(…, probeWaitBound)`. M5 now edits the constant, and the
    message reads `within 100ms` (V11).
  - V14 records who touches the file. V15 records the `:180-181` wiring check.
- **oc-kimi-k3: the full gate was unexecuted.** Applied: V16 runs `go build ./... && go test ./...`
  and `verify_ail.sh` on base + the r2 diff. AC-6 cites it.

**Round 2** (artifact `~/.ailang/state/mission-quorum-world/w-human-approve-control-window-flake-2026-11-02T12-59-22Z.json`,
2026-11-02T12:59:22Z, $0.3083): **BLOCKED**. All 3 seats present (`absent_reviewers` = `[]`):
- `oc-glm-5-3` REJECT: R3's "max == 1 would still catch a wrong-mutex mutant" was unmeasured and
  false; also §1's "Load only stretches …", an unlogged `:285-286` citation, derived line numbers.
- `oc-kimi-k3` PASS: add an adjudication rule for a lone empirical M3 survival; verify the scope of
  `boundedTestContext` in the fixed file.
- `gpt6-1-sol` REJECT: M5/AC-4 still treated a scheduling-dependent timeout as deterministic; M3's
  residual list was incomplete (delayed timer creation; simultaneously ready cases).

**Surface tracking (gate-2 rule (a), round 3 on):** r1 and r2 objections both land on one surface —
how strongly the doc labels its oracles (deterministic vs empirical) and what its residuals claim.
No reviewer disputes the direction (event-driven waits, a `TryLock` oracle, an anchored replay).
Not a split signal.

**r3 = controller narrow-refinement carve-out** (gate-2's carve-out: every remaining blocking
objection carries a concrete reviewer-authored `proposed_fix` and none disputes the direction). The
controller applied the reviewers' verbatim fixes and measured what they asked to be measured, with no
third designer run:
- glm: R3 replaced with the reviewer's text; M7 and M7b measured (V17: M7 killed by `unheld` 20/20,
  M7b survives 20/20 as predicted); §1's sentence replaced with the reviewer's text; V18 added; M1/M3
  re-measured on the final sha (V17).
- gpt: §2.5 adds the test-only timeout seam and handshake test verbatim; AC-4 now cites it; M5
  relabelled empirical; §2.3 lists the three late-first-goroutine cases. The seam test and its
  mutants are executor deliverables, measured at Gate 3 and by the judge.
- kimi (pass): AC-3 carries the adjudication text verbatim; V19 records the context scope.

