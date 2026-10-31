# w-exec-probe-reverse-leg-pins — bound arm 6's hello wait from both sides, pin the child listener's drain, and settle the cross-profile and collision-rate residuals of row 229

- **Status**: Planned (designed iteration 329, r2: quorum r1 fixes §11; quorum r2 narrow-refinement carve-out §12). Queue row 230. An amendment to
  `design_docs/implemented/w-exec-probe-reverse-leg-port-race.md` (row 229, the parent; its §4.1,
  §4.2, §5 R2 and R4) and a sibling of `design_docs/implemented/w-exec-probe-reset-holder-false-pass.md`
  (row 232).
- **Item**: Clause 3 (the Workspace.Exec fail-closed startup probe), hygiene: two mutants survive
  the parent's suite (iteration-314 judge O8 and O10), and two residuals were never measured.
- **Estimated**: ~0.25 day. Production about +8/−3 LOC in `host/broker/exec_probe.go` (one named
  constant, one derived message, two comments). Tests about +90 LOC in
  `host/broker/exec_probe_test.go`. One `ci.yml` list edit (§7).
- **Measurement base**: `e5c1867f` (worktree `sprint/world-iter329-exec-probe-reverse-leg-pins`,
  clean Go tree, V1).
- **Instruments**: `go1.26.6 darwin/arm64`, 16 cores, node v26.12.1 (V1). Load means 32 spinning
  `yes` on 16 cores, with the PIDs killed by a trap and 0 alive afterwards. The load average read
  15.08–24.42 at the start of the loaded runs (V12, V17). There is no Linux tool locally (V1). Every
  Linux statement below is either sourced from CI (V20) or marked **UNMEASURED** and handed to
  AC-LINUX-F3 (§3.3). The scratch sources, runners and raw outputs are checked in under
  `design_docs/verification/world-iter329/` as `.txt` files. Scratch tests were copied into
  `host/broker/` as `zz_iter329_*_test.go`, run, and deleted. After the deletion,
  `git status --porcelain -- host cmd world` printed 0 lines and `exec_probe.go` hashed to its
  base sha256 prefix `390c74e841f26a12` (V8, V9, V13).

---

## 1. The four items, and what this session measured

| # | Row 230 says | Verdict at `e5c1867f` | Evidence |
|---|---|---|---|
| 1 | F2: `execProbeHelloWait` (500 ms) is pinned by no test, and 5 ms keeps every test green | **CONFIRMED for the whole `host/broker` package, not only `-run ExecProbe`.** Control and the 5 ms mutant give the same 363 verdict lines (333 PASS, 0 FAIL, 30 SKIP), both `ok`. The 30 SKIPs are the real-srt and tool-binary tests, so they are unmeasured locally. **New: there is also no deterministic pin in the other direction.** At 60 s, the reverse-leg subset V4 ran (`-run '^(TestExecProbeReverse|TestExecProbeForeignListener)'`) reds only through T6/`silent`, a wall-clock path (§3.1); the whole package was not run at 60 s **[r2]** | V2, V3, V4 |
| 2 | F3: `c.resume()` "avoids a RST on close with an unread hello" and is pinned by no test | **Survival CONFIRMED** on the whole package: 363 identical verdict lines with `c.resume()` removed. **The stated mechanism is REFUTED on darwin.** Without `resume`, node still reads the 18-byte hello out of the kernel into its own buffer (`readableLength` 18), so no unread kernel data is left to cause a reset. What `resume` changes is node's server-side socket lifecycle: `end`/`close` fire 20/20 with it and 0/20 without it. **No host-visible difference was found**: the nonce arrived 20/20 in both variants and the child exited in the same time; under the production poll, REACHED 40/40 per variant, idle and loaded. On Linux this is **UNMEASURED** | V5, V6, V7, V10, V11, V12 |
| 3 | R2: one profile's poll can land on another profile's live counted listener and raise its accept count to 2 | **Refuted at HEAD by measurement; enforcement residual open [r2].** The hazard does not occur in-process today, but nothing enforces that; the structural fix (hoist `net6`'s bind into `ProbeExecProfilesEach` before the fan-out) is a tracked successor row, and P4 comments the invariant. Profiles do run concurrently. But (a) every profile's `net6` is bound before its own Run A and held through Run C, (b) `freeLoopbackPort` never returned a port a live listener holds (0 of 200,000 picks; the paired control, with the same ports released, gave 768), and (c) in 40 runs of 8 concurrent profiles (10 iterations × {noop, passthrough} × {idle, load}), every Run A started at least 216.9 ms before any profile's first pick (0 collisions). **If forced, the consequence is real and worse than the row says:** the owner's proxy leg refuses with `accepts 9`, *and* the poller's own leg refuses as SILENT, not foreign. There is a cross-process residual of the same class as R4 | V14, V15, V16, V17, V18 |
| 4 | R4: no Linux measurement of the in-window collision rate exists | **Status unchanged: still none.** Watch item only. In the 8 failed CI runs since 2026-10-28, `inconclusive after` appears 0 times. The only `answered nothing to World's hello` lines are T6's deliberate `t.Logf`. Runs re-run to green were not scanned | V20 |

## 2. The mechanisms

**The hello wait has two jobs, so it needs two bounds.** `classifyReverseDial` sets one deadline,
`time.Now().Add(execProbeHelloWait)` (`exec_probe.go:662`). It writes the hello and reads the answer
under that deadline. The read ends in one of three ways:

- the nonce line, which is **child**;
- other bytes or a clean EOF, which is **foreign**;
- zero bytes at the deadline or on a reset, which is **silent** (`:680-689`, V23).

The wait therefore decides two things:

- **Floor: a legitimate answer must arrive inside the wait.** The answer that matters for
  availability is srt's mux, the identified in-window foreign holder (parent V7). It answers the
  hello with HTTP 400 (parent V23). If that answer arrives after the deadline, the dial reads as
  silent, and silent is a **sticky refusal** (parent §4.2 row 3), so startup is refused falsely.
  For World's own child, the floor does not change safety, since REACHED and SILENT both refuse. It
  does keep the refusal text truthful.
- **Ceiling: silence must be observed while the holder still holds.** The child holds its listener
  for `execProbeListenMS` = 1000 ms (`:99`), and the poll dials every 50 ms (`:732`). Suppose a dial
  connects just after the bind and the peer then sends nothing. If the wait outlives the hold, the
  holder's exit closes the socket, the read sees EOF, and the dial is classed **foreign**, which is
  retried, instead of **silent**, which refuses.
  - Measured with a 60 s wait: T6/`silent`'s stand-in holder exits at its 1000 ms hold. Its attempt
    was retried and passed, so arm 6 carried no reverse-leg text and T6 went red (V4).
  - For a breached child that accepts and never answers, the same conversion weakens P20: a silent
    reach becomes a retry.
  - The ceiling has a second, liveness reason. The poll ignores `ctx`, and `stop()` waits for an
    in-flight classify (`:736-741`, V24). With a 60 s wait, each T7 silent row took 60 s and still
    passed (V4). Inside the daemon's 9 s startup context (parent V16), that is a 51 s overrun.
  - The ceiling is therefore **`execProbeListenMS` less one poll period, 950 ms**: the first
    connected dial lands at most one poll period after the bind, and its deadline must fall inside
    the hold. The child's process exits 74–152 ms after the hold ends (V6), which is extra slack the
    bound does not count on.

**Today's numbers.** 500 ms lies inside [250 ms, 950 ms], 2× the floor below and 450 ms under the
ceiling. The **floor of 250 ms** is set above the slowest legitimate answer measured: srt's mux
answered in at most 229.6 ms of 2,000 dials at load, and World's child in at most 160.1 ms
(V11, V12). The row's "~100 ms" is below that measured worst case, so the floor is raised, and this
is recorded as a deviation from the row text.

**`c.resume()`'s real effect.** In node, an accepted socket begins reading as soon as it is created,
whether or not anyone consumes the data. Without `resume`, the hello therefore sits in node's
stream buffer, not the kernel's (`readableLength` 18 at 100 ms, 20/20; V5). Two consequences
follow:

- The reset the parent worried about needs unread data **in the kernel** when the socket closes,
  and that data is not there in either variant.
- The stream is never consumed, so `end` and `close` never fire and the accepted socket stays open
  until the child process exits (0/20 vs 20/20, V5).

That is socket hygiene in the child. It is not a host-visible verdict. The buffering is node's
stream layer, which is the same code on Linux, so the reading predicts no Linux difference either.
That prediction is **UNMEASURED** and is the subject of AC-LINUX-F3.

**The cross-profile reach.** `ProbeExecProfilesEach` starts one goroutine per profile
(`:224-229`). In `probe()`, `net6 := newProbeListener()` (`:455`) comes before Run A (`:465`) and
Run B (`:496`). The reverse leg's first `pick()` (`:551`) comes after both, and `net6` stays bound
until `probe()` returns (`:456`). A port held by a live listener is never returned by `bind(0)`
(V15).

So profile X's pick can equal profile Y's `net6` port only if Y binds `net6` **after** X's pick.
That needs Y's setup to lag X's Run A and Run B together. Across 40 runs of 8 profiles (10 iterations
× {noop, passthrough} × {idle, load}), the lag never came closer than 216.9 ms to that point.
Idle (20 runs) the worst margin was 216.9 ms with the noop shim (the fastest Run A and Run B) and
572.6 ms with the pass-through shim. Loaded (20 runs) it was 280.5 / 933.8 ms (V16, V17). The
runtime re-check path is single-flight (`exec_tools.go:344-346`), and no `host/broker` or
`host/daemon` test calls `t.Parallel()` (V19).

The exclusion therefore holds in-process, **but by measured ordering, not by a lock** (R1, §8).
Across processes, for example package binaries run in parallel by the race gate's `go test ./...`,
the ordering does not apply. A collision there needs another process's `bind(0)` to return exactly
the port just freed: about 1 in 16,384 per bind on darwin (V15; range 49152–65535). The Linux rate
is unmeasured. That is the R4 class.

## 3. Design

### 3.1 F2: pin the hello wait from both sides (production: P1, P2)

- **P1, a named poll period.** `const execProbePollEvery = 50 * time.Millisecond` sits beside
  `execProbeHelloWait`. `pollLoopback`'s `time.After(50 * time.Millisecond)` (`:732`) becomes
  `time.After(execProbePollEvery)`. Behaviour is unchanged. The constant exists so that T10's
  ceiling derives from code, not from a literal in a test.
- **P2, the SILENT text names the wait it used.** Today `:582` hard-codes `within 500 ms`. It
  becomes `… within %d ms; …` with `execProbeHelloWait.Milliseconds()`. At 500 the output is
  byte-identical, so T7, T6 and QUICKSTART §10 (which never quotes the number, V21) are unchanged.
  Measured: with the wait at 400 ms and today's literal, the refusal still says `within 500 ms`
  (V13).
- **P3, comments only.** `execProbeHelloWait`'s doc comment states the two bounds and why. The
  comment above `execProbeScript` (`:62-67`) states what `c.resume()` does, and does not do,
  per §2.
- **P4, comments only [r2].** At `net6 := newProbeListener()` (`:455`) and at the fan-out
  `go func()` (`:226`): every profile binds `net6` before its own Run A, so no profile's reverse-leg
  pick can return another profile's live `net6` port (`bind(0)` never returns a held port, V15);
  moving the bind later, or adding a slow step before it, reopens the cross-profile collision (§3.4,
  R1). Behaviour-identical, so no mutant.

**T10 shape** (deterministic: a simulated arrival, no sleep, no socket):

```go
// helloConn's peer answer ARRIVES at a fixed instant on the test's clock: Read hands it over only
// when the classifier's deadline is at or after that instant, else reports the deadline passed.
type helloConn struct { net.Conn; arrival time.Time; answer []byte; deadline time.Time; sets int; wrote string }
// SetDeadline/SetReadDeadline record; Write appends; Read: deadline set && deadline < arrival →
// os.ErrDeadlineExceeded; else one byte of answer, then io.EOF.

const execProbeHelloFloor = 250 * time.Millisecond // > the 229.6 ms worst mux answer at load (V11)

func TestExecProbeReverseHelloWaitBounds(t *testing.T) {
	// rows: slow-child (want+"\n" → {child}), slow-mux ("HTTP/1.1 400 Bad Request\r\n" → {foreign})
	// The ceiling derives from BOTH production constants, never a test literal: shorten the hold
	// or lengthen the poll period and the ceiling tightens with it (K9).
	listenWindow := time.Duration(execProbeListenMS) * time.Millisecond
	before := time.Now()
	conn := &helloConn{arrival: before.Add(execProbeHelloFloor), answer: []byte(row.answer)}
	seen := pollLoopback(40000, want, firstDialReturns(conn, refused))() // dials once before stop
	after := time.Now()
	// 1. seen == row.want and conn.wrote == execProbeHello            (floor)
	// 2. conn.sets > 0 && !conn.deadline.IsZero() &&
	//    !conn.deadline.After(after.Add(listenWindow - execProbePollEvery))   (ceiling)
}
```

**Why it cannot flake (requirement c).** The classifier computes `deadline = Now_c + wait`, where
`before ≤ Now_c ≤ after`. Check 1 holds whenever `Now_c + wait ≥ before + floor`, and that is
true for any `wait ≥ floor`, whatever the scheduler does. Check 2 holds whenever
`Now_c + wait ≤ after + 950 ms`, and that is true for any `wait ≤ 950 ms`, whatever the scheduler
does. Load can only make a **mutant** survive one run, never redden a correct tree:

- A stall of ≥ (floor − wait) between `before` and `Now_c` hides an undersized wait.
- A stall between `Now_c` and `after` hides an oversized one.

At load ~16 on 16 cores, the 250 ms prototype passed `-count=3000`, both rows, at load average
15–19 (V17). An earlier one-row prototype with a 100 ms floor also passed at 16.78 (V12).
The exact-boundary kills, 249 ms and 951 ms, are measured idle. Under extreme stalls they may
survive a run, as the bullets above say.

### 3.2 F3: the structural pin, chosen by measurement (T11, P3)

**Choice: a structural check.** A Linux-CI arm is not used: on the evidence, no platform shows a
host-visible effect, so an arm could not be non-vacuous.

- The candidate behavioural arm is a real child reached under the hello. **It already exists**: T4,
  `TestExecProbeReverseLegRealReachSurvivesAForeignAttempt`, runs World's node client unsandboxed
  through the pass-through shim. It is in CI step 8's PASS-guarded list, so it cannot skip on Linux
  (V22).
- T4 stays green with `c.resume()` removed on darwin (V3). The production-poll probe gave REACHED
  in every cycle in both variants, idle and loaded (V10, V12). The mechanism the line was said to
  guard is absent at the node layer (V5).
- A new Linux arm would therefore be predicted vacuous on its own platform. That is exactly what
  requirement (2) forbids.

**T11 `TestExecProbeChildListenerDrainsTheHello`** asserts two things:

- `execProbeScript` contains exactly one `net.createServer(`;
- that handler is `net.createServer(c=>{c.resume();c.end('WORLD-PROBE-NONCE '+nonce+'\n')})`.

**S6 label:** this is a structural pin of the child's socket hygiene, as P3's comment documents. It
does **not** discharge a behavioural claim, and the doc makes none for it. It is killed by K6
(V9). It is pure, so it runs and passes on every platform, with no skip.
**Acknowledged trade [r2, quorum r2 kimi]:** T11 pins, by exact source-string match, a line that
V5/V10/V12 show is host-visibly inert on the only measured platform. It guards against silent
removal of deliberate child socket hygiene, not against a measured failure; if AC-LINUX-F3 also
shows no difference, that remains its whole value, and a future reviewer may delete T11 together
with `c.resume()` by an explicit decision rather than by accident.

### 3.3 AC-LINUX-F3: the controller's Linux measurement (non-blocking)

Use a throwaway branch with row 232's V14 workflow pattern, on **`ubuntu-24.04`** (the `go-verify`
image, `ci.yml:121`; *not* `ubuntu-latest`). Copy in this session's scratch sources unchanged and
run three things:

1. `node f3-node-socket-probe.js.txt resume 20` and `… noresume 20`.
2. `node f3-child-exit-probe.js.txt host/broker/exec_probe.go resume 20 1000`, and the same with
   `noresume`.
3. `zz_iter329_f3_test.go.txt` as a package test: `TestZZIter329RealChildUnderPoll` with
   `ZZ_F3_VARIANT=resume|noresume ZZ_F3_CYCLES=40`, and `TestZZIter329ChildAnswerLatency`. Run both
   idle and under 2×nproc `yes`. Then run step 8's test list with the F3 mutant applied, **with
   `TestExecProbeChildListenerDrainsTheHello` (T11) excluded** from that run, or with the throwaway
   branch cut from the pre-M2 commit, where T11 does not exist yet. T11 reds structurally under
   the mutant on every platform, so leaving it in would make "a step-8 red" certain rather than a
   behavioural signal.

**Decision rule.** If `noresume` differs on Linux in any host-visible way, T4, already in step 8,
is the Linux arm that kills K6, and §3.2's rationale is amended. Host-visible means any of:

- a cycle that is not REACHED;
- a client error before the nonce;
- a child exit more than 100 ms later than `resume`'s maximum;
- a step-8 red from any test other than T11 (T11 is excluded from that run).

If `noresume` shows no difference, record "no difference" and T11 stands as the only pin. **Either
way the sprint's code is the same.** AC-LINUX-F3 settles a claim, not a code path. It is
**non-blocking**: the controller may run it before or after the merge, and a missing run leaves §1
row 2's Linux cell UNMEASURED.

### 3.4 Item 3: no code; the verdict and its trigger

**[r2, quorum r2 kimi option (b)]** The verdict is *refuted at HEAD by measurement; enforcement residual open*, with V14–V18 as the evidence. The in-process exclusion rests on statement order (R1), and V18 shows a forced collision refuses BOTH profiles, so the controller files a **successor row** for the structural fix: hoist the per-profile `net6` bind out of `probe()` (`:455`) into `ProbeExecProfilesEach` ahead of the goroutine fan-out (`:224-229`), handed into `probeExecProfile` and closed by a parent-level defer on its early-return paths (`:369-401`), so port exclusion holds by construction; that row's design must show a mutation that moves the bind back after a slow pre-bind step either reds a test or is impossible by construction. This sprint adds only **P4**: a comment at `:455` and at `:226` stating the ordering invariant the code relies on. The cross-process residual joins R4's watch (§3.5).

The forced reproduction (V18) shows what a collision would print. Profile Y prints
`a connect through srt's proxy environment: token "WORLD-PROBE-REFUSED ECONNREFUSED", accepts 9`,
and the poller X prints `reverse leg: a listener on 127.0.0.1:<Y's net6> accepted the host's dial
and answered nothing to World's hello`. `net6` accepts and closes at once, so World's hello meets a
reset and the dial reads silent. Both profiles refuse.

**Trigger:** in a real-srt or multi-profile CI log, an `accepts N` with N ≥ 2 on a raw or proxy
leg. T6's and the pass-through tests' deliberate lines are excluded: those tests run the connect
unsandboxed.

### 3.5 Item 4: watch item (no infrastructure)

The status is unchanged since iteration 314: there is no Linux measurement of the in-window
collision rate. V20 scanned the 8 failed CI runs since the parent landed (2026-10-28) and found no
instance of either signature outside deliberate test logs.

**Trigger, extended by §3.4.** Any of the following, in a CI log where it is not a test's
deliberate `t.Logf` or assertion, means re-measuring the collision rate on Linux *first*, before
any retry. Use the parent's V14-style probe, adding srt's mux and a parallel `go test ./...` as
binders.

- `reverse leg: inconclusive after 3 attempts`;
- `reverse leg: a listener on … answered nothing to World's hello`;
- §3.4's `accepts N≥2`.

## 4. Test plan

Each test is one row, and each kills named mutations (§5). "Measured" means the mutation was run
against this session's prototype of the test, at `e5c1867f`.

| Test | Kind | Needs node? | Kills | Measured? |
|---|---|---|---|---|
| **T10** `TestExecProbeReverseHelloWaitBounds` (new; rows `slow-child`, `slow-mux`) | deterministic, injected dial + simulated arrival | no | K1 (5 ms), K2 (249 ms), K3 (951 ms), K4 (60 s), K5 (no deadline), K9 (hold halved) | **yes**: V8, V25 (control green; each mutant red with the named assertion) |
| **T11** `TestExecProbeChildListenerDrainsTheHello` (new) | structural, S6-labelled | no | K6 (`c.resume()` removed) | **yes**: V9 |
| **T7** `TestExecProbeReverseLegSilentHolderRefusesSticky` (extended: also want `within <execProbeHelloWait.Milliseconds()> ms;`) | existing fixture | no | K7 (the SILENT text not derived from the wait) | **yes, on today's code**: V13. With wait 400 ms and the literal kept, it is red; at 500 it is green. After P2, K7 is the planner's revert-P2-and-set-400 mutant, **to be re-proven by the executor** |
| **T6** `TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse` (unchanged) | wall-clock (stand-in holder lifetime) | yes | K4 as a side effect only | yes: V4 (60 s → `silent` red). Not relied on: it is timing-bound, and T10 is the pin |
| **T4** (unchanged; CI step 8) | real node child, pass-through | yes | K6 on Linux **only if** AC-LINUX-F3 shows a difference | darwin: survives K6 (V3). Linux: UNMEASURED |

**AC-MUT.** Apply K1–K7 and K9 one at a time. Each reds its named killer, and the tree is restored
byte-identical (sha256 of both Go files). **AC-GATE**: `gofmt -l` is empty, `go vet ./host/broker/`
is clean (clean at base, V1), `go build ./... && go test ./...` passes, and
`./scripts/verify_ail.sh` passes. No `.ail` is touched. **AC-CI**: T10 and T11 join step 8's
`tests=` list (`ci.yml:247`). Both are pure and cannot skip, but the family's rule is "every
reverse-leg test is PASS-guarded" (parent §6), so the list stays complete. **AC-STABLE**:
`go test ./host/broker -run '^TestExecProbeReverse|^TestExecProbeChildListener' -count=20` passes
20/20. Then T10 `-count=3000` under 32 `yes` gives 0 FAIL, with the prototype's V17 as the
baseline.

## 5. Mutation table (continues the parent's M1–M21, row 232's M22–M29, the judge's O1–O10)

| ID | Mutant | Killer, and the assertion that fires | Run? |
|---|---|---|---|
| K1 (= O8) | `execProbeHelloWait = 5 * time.Millisecond` | T10 both rows: `seen {silent:true}`, want `{child}` / `{foreign}` | V8 red |
| K2 | `= 249 * time.Millisecond` (just under the floor) | T10 both rows, the same assertion | V8 red |
| K2c (control) | `= 250 ms` / `= 950 ms` (the bounds themselves) | none: T10 stays green, which proves the bounds are inclusive | V8 green / green |
| K3 | `= 951 * time.Millisecond` (just over the ceiling) | T10: `hello deadline at stop+950.9ms … want … no later than stop+950ms` | V8 red |
| K4 | `= 60 * time.Second` | T10 ceiling assertion (deterministic); T6/`silent` (wall-clock) | V8 red; V4 red |
| K5 | the `SetDeadline` line removed from `classifyReverseDial` | T10: `sets 0` | V8 red |
| K6 (= O10) | `c.resume();` removed from the listen op | T11: `the listen op's connection handler is not …` | V9 red; the full package stays green (V3) |
| K7 | the SILENT text keeps a literal while the wait changes (500 → 400 ms) | T7 extension: `want it to contain "… within 400 ms;"` | V13 red (on today's code) |
| K9 | `execProbeListenMS = 500` (the hold halved; wait stays 500 ms) | T10 ceiling: the bound becomes 500 − 50 = 450 ms, and a 500 ms wait exceeds it: `hello deadline at stop+499.9ms … want … no later than stop+450ms`. **This red is the intended coupling**: with a 500 ms hold, a 500 ms wait outlives the hold for a dial landing one poll period after the bind, so a silent holder's exit would read as foreign, which is retried (§2's ceiling). Whoever shortens the hold must shorten the wait too | V25 red (prototype; ceiling derived from `execProbeListenMS`) |
| K8 | `execProbePollEvery` not used by `pollLoopback` (the literal 50 ms kept) | **none**: behaviour-identical, so it is no mutant. P1's constant is checked by code review; T10 derives its ceiling from it | n/a |

## 6. Milestones

- **M1, F2.** P1 and P2, then T10 (with `helloConn` and `execProbeHelloFloor`), then the T7
  extension. Prove K1–K5, K7 and K9 red one at a time. A `ci.yml:247` gain of T10.
- **M2, F3 and record.** P3's two comments and P4's two ordering-invariant comments [r2], then T11. Prove K6 red. A `ci.yml:247` gain of T11. The
  controller runs AC-LINUX-F3 (§3.3, non-blocking) and records it under `design_docs/verification/world-iter329/`.
  Items 3 and 4 need no code. The controller's record carries §3.4 and §3.5's verdicts and
  triggers into the charter row's closing note. That is the controller's edit, not the sprint's.

M1 and M2 touch disjoint hunks and land independently.

## 7. Files

| File | Change | LOC (est.) |
|---|---|---|
| `host/broker/exec_probe.go` | P1, `execProbePollEvery` (new const near `:110`; use at `:732`). P2, `:582`'s SILENT text derived from `execProbeHelloWait`. P3, the doc comments at `:62-67` and `:109`. P4, the ordering-invariant comments at `:455` and `:226` | +12 / −3 |
| `host/broker/exec_probe_test.go` | `helloConn`, `execProbeHelloFloor`, T10, T11, T7's one extra assertion | +90 |
| `.github/workflows/ci.yml` | `:247` `tests=` gains `TestExecProbeReverseHelloWaitBounds TestExecProbeChildListenerDrainsTheHello` | ±1 |

`docs/QUICKSTART.md` is not touched: the user-visible text is byte-identical at 500 ms, and §10
never quotes the number (V21). **S7** therefore does not apply. **S2/S3**: this is a repair inside
the host's effect boundary with no new surface, the same answer as the parent §8 to "why not a
package?". **S1/S4/S5** do not apply, since there is no `.ail`. No dependency is added.

## 8. Conflict surface, risks and residuals

**Conflict surface.** The edit is three hunks of `exec_probe.go`:

- the constants block (`:104-110`);
- the SILENT `Sprintf` (`:582`);
- `pollLoopback`'s `select` (`:732`).

There are also two comments. `execProbeHelloWait` is used only in `exec_probe.go` (`:110`, `:662`; P2 adds `:582` as a known post-change use site;
plus comments), with `execProbeListenMS`'s census as the control (V2). No test references
`execProbeScript` or `c.resume` (V7). Arms 1–5, 7 and 8 are untouched.

- **R1, the cross-profile exclusion is measured, not enforced.** It rests on "every profile binds
  `net6` before any profile's first pick". The worst margin seen was 216.9 ms, in 40 runs (20 idle, 20 loaded; V16, V17).
  - **What could reopen it:** a future edit that moves `newProbeListener` later in `probe()`, or
    that adds a slow, stall-prone step before it.
  - **Not proposed:** a barrier across profiles, where every profile signals after binding `net6`
    and every reverse leg waits for all. Each early-return path in `probeExecProfile` (`:369-401`)
    would have to signal, or the call deadlocks. That is a larger concurrency surface than a
    residual nobody has observed justifies.
  - **[r2]** The structurally simpler alternative, hoisting the `net6` bind into
    `ProbeExecProfilesEach` before the fan-out (no signalling, no deadlock surface), is NOT rejected:
    it is the successor row named in §3.4, kept out of this hygiene sprint because it changes
    probe structure on a fail-closed startup gate and needs its own design and quorum.
  - The R4 trigger (§3.5) catches it if it ever happens.
- **R2, cross-process collisions** (another process's `bind(0)` returns the port a poll just
  freed). These are possible in the race gate, where `go test ./...` runs package binaries in
  parallel (`verify_go.sh:350`), and three test packages probe: `cmd/ailang-worldd` (5 files), `host/broker` (5) and `host/daemon` (2) (V19b). The darwin rate is about 1/16,384
  per bind (V15). The Linux rate is UNMEASURED. This is the R4 class.
- **R3, F3 on Linux is UNMEASURED.** §3.3's decision rule handles both outcomes. If Linux
  *does* reset, T4 already sits on the CI platform as the killer.
- **R4, the floor is an observation.** The 229.6 ms worst mux answer comes from 2,000 dials on one
  loaded Mac (V11). A slower CI runner could exceed 250 ms. The floor is a minimum for future
  edits, not the shipped value, and 500 ms keeps 2.2× headroom over the observed worst case. If the
  trigger in §3.5 fires with the SILENT text, compare it against this first.
- **R5, the boundary kills are idle-exact.** K2 and K3 can survive a single run under an extreme
  stall (§3.1). The direction is safe: a mutant survives, but a correct tree never reddens.
- **Not changed, carried forward:** the parent's R5 (a foreign holder's exit classed SILENT, row
  232 §8) and R6 (darwin EINVAL). The hello wait's ceiling makes R5's window no wider.

## 9. What this item is NOT doing

- No barrier or registry for cross-profile reaches (R1).
- No collision-rate infrastructure (item 4 is a watch).
- No change to `execProbeListenMS`, the poll cadence, the dial timeout, `judgeReverse` or any
  refusal text. P2 produces the same bytes.
- No Linux claim beyond what AC-LINUX-F3 measures.

## 10. Verification Log

The base for every row is `e5c1867f`. `V=design_docs/verification/world-iter329`.
`AILANG_BIN=$HOME/.pinned-ailang/ailang` is set for every `go test`. Each mutant was applied with
`sed`, then reverted with `git checkout -- host/broker/exec_probe.go`, and the base sha256 prefix
`390c74e841f26a12` was re-checked.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Base, tooling, vet | `git rev-parse HEAD`; `git status --porcelain -- host cmd world \| wc -l`; `go version`; `sysctl -n hw.ncpu`; `node --version`; `which docker podman colima orb limactl lima nerdctl`; `go vet ./host/broker/` | `e5c1867f717c…`; `0`; `go1.26.6 darwin/arm64`; `16`; `v26.12.1`; every Linux tool `not found`; vet rc 0 |
| V2 | Today's constants, deadline, cadence, literal text; who uses the wait | `grep -n 'const execProbeListenMS\|const execProbeHelloWait\|…SetDeadline(time.Now().Add(execProbeHelloWait))\|case <-time.After(50 \* time.Millisecond):\|within 500 ms\|…200\*time.Millisecond)\|const execProbeRunBudget' host/broker/exec_probe.go`; `grep -rn execProbeHelloWait --include='*.go' .`; control `grep -rn execProbeListenMS --include='*.go' .` | `:96` RunBudget 8 s; `:99` `execProbeListenMS = 1000`; `:110` `500 * time.Millisecond`; `:582` `…World's hello within 500 ms…`; `:662` SetDeadline; `:704` 200 ms dial; `:732` 50 ms. The wait has 4 hits, all in `exec_probe.go` (`:109,:110,:657,:662`). **[r2]** `:109` and `:657` are comments: `sed -n 657p host/broker/exec_probe.go` → `// at most len(want)+1 bytes up to '\n' within execProbeHelloWait. The line` (the doc comment above `classifyReverseDial`); so the code uses are `:110` (decl) and `:662`, and P2 adds `:582`. Control: 3 hits (`:98,:99,:566`), so the grep works |
| V3 | **F2 and F3 survive the WHOLE package** | `bash $V/run-fullpkg-f2.sh.txt` (control, then wait 5 ms) and `bash $V/run-fullpkg-f3.sh.txt` (`c.resume();` removed); verdict lines diffed after stripping timings. Summaries: `$V/out/fullpkg-{control,f2-5ms,f3-noresume}.summary.txt` | control `ok … 235.024s`, 333 PASS / 0 FAIL / 30 SKIP. 5 ms `ok … 224.279s`, the same. noresume (`c.resume() occurrences after mutation: 0`) `ok … 221.705s`, the same. Both mutants show `same 363 verdict lines as control`. T4 `--- PASS` in all three. The SKIPs are the real-srt and tool-binary tests (`TestExecSrt*`, `TestExecProbeArm8`, `TestAilangRun*`, …) |
| V4 | A 60 s wait against the existing reverse-leg tests | wait `60 * time.Second`; `go test -count=1 -timeout 280s -run '^(TestExecProbeReverse\|TestExecProbeForeignListener)' -v ./host/broker/` → `$V/out-f2-60s-existing.txt` | rc 1 after 132 s. T7 `silent-attempt-1` PASS **60.03s**, `silent-after-chatty` PASS **60.21s**. `--- FAIL: …ForeignListenerOnTheReversePortDoesNotRefuse/silent`, `exec_probe_test.go:733: silent: arm 6 must carry the silent refusal, not the reached text: …` (no reverse-leg text). Every other test PASS |
| V5 | Node buffers the hello without `resume`; only the server-side lifecycle differs | `node $V/f3-node-socket-probe.js.txt resume 20`; `… noresume 20` | resume `serverEnd 20, serverClose 20, readableLength@100ms {"0":20}, clientNonce 20, clientErr {}`. noresume `serverEnd 0, serverClose 0, readableLength@100ms {"18":20}, clientNonce 20, clientEnd 20, clientErr {}` |
| V6 | The child's verbatim listen op (extracted from `exec_probe.go`): the nonce and exit timing in both variants | `node $V/f3-child-exit-probe.js.txt host/broker/exec_probe.go {resume,noresume} 20 1000` (the script refuses to run if no `c.resume();` is there to remove) | resume `nonce 20, other 0, err {}, BOUND 20, exit 1083–1119 ms`. noresume `nonce 20, other 0, err {}, BOUND 20, exit 1074–1152 ms` |
| V7 | No test references the script or `resume` | `grep -rn 'c\.resume' --include='*_test.go' host cmd \| wc -l`; `grep -rn execProbeScript --include='*_test.go' host cmd`; control `grep -rn 'revHolder(' --include='*_test.go' host cmd \| wc -l`; `grep -c 'c.resume()' host/broker/exec_probe.go` | `0`; no lines; control `2`; `1` in production |
| V8 | **T10 prototype vs K1–K5 and the bound controls** | `bash $V/run-f2-mutants.sh.txt` over `$V/zz_iter329_f2_test.go.txt` (sha256 prefix `f7d26a0ea751`) → `$V/out-f2-mutants.txt` | control PASS. 5 ms FAIL: `zz_iter329_f2_test.go:76: an answer 250ms after the hello: seen {… silent:true} … want {child:true …}` and the `{foreign:true}` row. 249 ms FAIL (the same). 250 ms PASS. 950 ms PASS. 951 ms FAIL `:79: hello deadline at stop+950.876541ms (sets 1); want … no later than stop+950ms`. 60 s FAIL `stop+59.999914875s`. no-deadline FAIL `(sets 0)`. `restored sha 390c74e841f26a12` |
| V9 | **T11 prototype vs K6** | `$V/zz_iter329_t11_test.go.txt` copied in; control, then `sed 's/net.createServer(c=>{c.resume();c.end(/net.createServer(c=>{c.end(/'` → `$V/out-t11-mutant.txt` | control `--- PASS`. Mutant (`0 occurrences left`) `zz_iter329_t11_test.go:18: the listen op's connection handler is not net.createServer(c=>{c.resume();c.end('WORLD-PROBE-NONCE '+nonce+'\n')})`, `--- FAIL`. `restored 390c74e841f26a12` |
| V10 | The production poll against World's real child, idle | `ZZ_F3_VARIANT={resume,noresume} ZZ_F3_CYCLES=40 go test -run '^TestZZIter329RealChildUnderPoll$' -v ./host/broker/` (`$V/zz_iter329_f3_test.go.txt`, prefix `c113a97f494a`) → `$V/out-f3-realchild-idle.txt` | resume `seen={child:true …}` 40/40. noresume `{child:true}` 39 + `{child:true … silent:true}` 1. Every cycle REACHED (child-first). The stray silent is row 232's close-time reset (its §3 reading 1), seen in both variants under load (V12) |
| V11 | Answer latency to the hello: child and srt's mux, idle and at load (the floor's basis) | `$V/mux-server.mjs.txt` on srt 0.0.78 (`cli.js` sha256 prefix `3c3092bd26b3`, equal to World's pin by V26) plus `TestZZIter329MuxAnswerLatency` and `TestZZIter329ChildAnswerLatency`, 2,000 dials each, idle then under `$V/run-load.sh.txt` → `$V/out-mux-latency.txt`; 500-dial child runs → `$V/out-latency-idle.txt`, `$V/out-load.txt` | idle: mux `max=7.059916ms` (p99 0.82 ms), child `max=2.791041ms`. Load (avg 24.42): mux `ok=2000 … p99=9.149333ms max=229.593916ms`, child `max=125.657417ms`. 500-dial runs: idle max 1.50 / 2.26 ms (resume / noresume); load max 125.26 / 160.07 ms. `bad=0` everywhere |
| V12 | V10 and V11 under load; the T10 prototype (100 ms floor) at load | `$V/run-load.sh.txt` wrapping the latency, real-child and `-test.count=3000 -test.run '^TestZZIter329HelloWait'` runs → `$V/out-load.txt` | `load averages: 16.78`. Real child, load: resume `{child}` 39 + `{child,silent}` 1; noresume 38 + 2; all REACHED. `-count=3000` `PASS`. `load runner rc=0 yes-alive-after=0` |
| V13 | The SILENT text is a literal | `$V/zz_iter329_msg_test.go.txt` (T7's fixture, `silent` holder) at control, then with wait `400 * time.Millisecond` → `$V/out-msg-mutant.txt` | control `--- PASS`. 400 ms: `Why = "reverse leg: a listener on 127.0.0.1:53407 … within 500 ms; …"; want it to contain "…within 400 ms;"`, `--- FAIL`. `restored 390c74e841f26a12; porcelain-go=0` |
| V14 | Profiles are probed concurrently; `net6` is bound before Run A and Run B, held through the reverse leg, and the pick comes after | `grep -n 'for i, p := range profiles {\|go func() {$\|errs\[i\] = probeExecProfile' host/broker/exec_probe.go`; `grep -n 'net6 := newProbeListener()\|defer net6.close()\|a := pp.run(ctx, argv, true…\|b := pp.run(ctx, \[\]string{"/bin/sh"\|c, revWhy := pp.reverseLeg\|port, err := pick()\|poll := pollLoopback(port\|r := pp.run(ctx, argv, false\|seen := poll()\|acceptsA :=\|acceptsC :=' …` | `:224` loop, `:226` `go func() {`, `:228` probeExecProfile. `:455` net6, `:456` defer close, `:465` Run A, `:467` acceptsA, `:496` Run B, `:506` reverseLeg, `:508` acceptsC, `:551` pick, `:567` poll, `:568` run, `:569` stop |
| V15 | **`bind(0)` never returns a live listener's port** (paired control) | `$V/zz_iter329_alloc_test.go.txt` (prefix `6d77810a767b`): hold 64 listeners, `freeLoopbackPort` ×200,000; release the same 64; ×200,000 again → `$V/out-alloc.txt`; `sysctl net.inet.ip.portrange.{first,last}` | `hits-on-LIVE-ports=0 (distinct ports 16290); same 64 released: hits=768 (distinct 16354)`. `first: 49152`, `last: 65535`. So about 768/(64×200,000) ≈ 1/16,667 per bind per port |
| V16 | **Ordering margin, idle**: every Run A starts before any first pick (8 profiles × 10 runs per shim) | `$V/zz_iter329_xprof_test.go.txt` (prefix `2c0652c45a99`; a logging shim that rewrites the client's connect/proxy target to port 1, so `net6` is never dialled by World's own legs). `ZZ_XP_MODE={passthrough,noop} ZZ_XP_N=8 ZZ_XP_ITERS=10 go test -run '^TestZZIter329CrossProfileOrdering$' -v` → `$V/out-xprof-idle.txt` | passthrough: margins 572.6–998.7 ms, `worst-margin=572.6ms collisions=0`, `profiles-with-accepts>=2=0` in every run. noop: 216.9–516.3 ms, `worst-margin=216.9ms collisions=0`, accepts≥2 0 in every run |
| V17 | The same at load; the T10 prototype (250 ms floor) `-count=3000` at load | `$V/run-load.sh.txt` wrapping noop/passthrough ×10 and `-test.count=3000 -test.run '^TestZZIter329HelloWait'` → `$V/out-xprof-load.txt` | `load averages: 15.08 19.21`. noop `worst-margin=280.5ms collisions=0`; passthrough `worst-margin=933.8ms collisions=0`. `-count=3000` `PASS`. `yes-alive-after=0` |
| V18 | **A forced collision's consequence** (the picker returns the first-logged `net6` port) | `ZZ_XP_MODE=forced ZZ_XP_N=2 ZZ_XP_ITERS=3 go test -run '^TestZZIter329CrossProfileOrdering$' -v` → `$V/out-xprof-forced.txt` | Each of the 3 runs: `pick∩net6=2 profiles-with-accepts>=2=1`. The owner: `arm6-network: a connect through srt's proxy environment: token "WORLD-PROBE-REFUSED ECONNREFUSED", accepts 9 …; reverse leg: a listener on 127.0.0.1:60308 … answered nothing to World's hello within 500 ms`. The other profile: `arm6-network: reverse leg: a listener on 127.0.0.1:60308 … answered nothing …` (SILENT, not foreign) |
| V19 | Callers, single flight, `t.Parallel`, probing packages | `grep -rn "ProbeExecProfiles\|probeExecProfile(" --include='*.go' host cmd \| grep -v _test.go`; `grep -n 'cfg.Probe(\|s.cfg.Probe = ProbeExecProfilesEach' host/broker/exec_tools.go`; `grep -n 's.flight = flight$\|if flight == nil {$\|A shared flight' host/broker/exec_tools.go`; `grep -rn 't.Parallel()' --include='*_test.go' host/broker host/daemon \| wc -l`; control `grep -rln 't.Parallel()' --include='*_test.go' . \| wc -l`; `grep -n 'go test ./... -count=1 -race' scripts/verify_go.sh` | Callers of `ProbeExecProfilesEach`: `exec_probe.go:201` (the `ProbeExecProfiles` wrapper), `exec_startup.go:235` (startup), and `exec_tools.go:134` `s.cfg.Probe = ProbeExecProfilesEach`, called at `:616` `errs, err := s.cfg.Probe(probeCtx, cfg, order)`. `probeExecProfile` has its only call at `:228`. `:106` "A shared flight serializes re-checks", `:344` `if flight == nil {`, `:346` `s.flight = flight`. `0`; control `3`. (r0 listed probing test files here from a `\| head`-truncated grep that dropped `host/daemon`; r1 moves the census to V19b.) `verify_go.sh:350` `go test ./... -count=1 -race` (no `-p`, so packages run in parallel) |
| V19b | **[r1]** Census of probe-bearing test packages (controller-measured at `e5c1867f`, re-run first-party) | `grep -rln 'ProbeExecProfiles\|ExecStartupConfig\|ExecProfileFlags\|exec-profile' --include='*_test.go' . \| xargs -n1 dirname \| sort \| uniq -c` | `5 cmd/ailang-worldd`, `5 host/broker`, `2 host/daemon`: three packages |
| V20 | CI status for R4 since the parent landed | `gh run list --repo sunholo-data/ailang-world --workflow ci.yml --status failure --limit 40 …`; for the 8 runs since 2026-10-28: `gh run view <id> --log-failed`, counting `inconclusive after`, `answered nothing to World's hello`, `the host reached the listener`, `accepts [2-9]`; context lines for the hits → `$V/out-ci-scan.txt` | `inconclusive=0` in all 8. `silent=1` only in 40851432834 and 40769031350, each `exec_probe_test.go:724: silent: arm 6 = "a raw connect … accepts 2 …"` (T6's deliberate `t.Logf`); their failure is row 232's `TestExecProbeReverseLegSilentHolderRefusesSticky/reset-holder`. The other 6: all 0. **Limitation:** runs whose later attempt went green are listed as success and were not scanned |
| V21 | QUICKSTART §10 never quotes the wait | `grep -c '500 ms' docs/QUICKSTART.md`; control `grep -c "answered nothing to World" docs/QUICKSTART.md` | `0`; control `1` (`:1168`) |
| V22 | CI step 8: image, list, PASS guard | `grep -n 'runs-on: ubuntu-24.04\|name: go host build + test gate\|TestExecProbeReverseLegRealReachSurvivesAForeignAttempt' .github/workflows/ci.yml`; `grep -n 'did not PASS on linux (SKIP?)' .github/workflows/ci.yml` | `:117` `name: go host build + test gate`, `:121` `runs-on: ubuntu-24.04`, `:247` the list carrying T4 and every reverse-leg test; `:257` `…$t did not PASS on linux (SKIP?)` |
| V23 | **[r1]** The classification switch reads as §2 claims | `sed -n 680,689p host/broker/exec_probe.go` (controller-measured, re-run first-party) | `switch {` / `case string(buf) == want+"\n":` / `return true, false` / `case len(buf) > 0:` / `return false, true` / `case werr == nil && errors.Is(rerr, io.EOF):` / `return false, true` / `}` / `return false, false // silent` / `}` |
| V24 | **[r1]** `stop()` waits for the in-flight poll goroutine | `sed -n 736,741p host/broker/exec_probe.go` (controller-measured, re-run first-party) | `return func() reverseSeen {` / `close(done)` / `<-finished` / `return reverseSeen{child: child.Load(), foreign: foreign.Load(), silent: silent.Load()}` / `}` / `}`. The goroutine closes `finished` only on return (`:710`, `defer close(finished)`), after any classify in progress |
| V25 | **[r1] K9: the hold halved against the T10 prototype**, whose ceiling is `time.Duration(execProbeListenMS)*time.Millisecond - 50*time.Millisecond` | `$V/zz_iter329_f2_test.go.txt` copied in; control; then `sed 's/^const execProbeListenMS = 1000$/const execProbeListenMS = 500/'` → `$V/out-k9-mutant.txt` | control `--- PASS`. Mutant: `zz_iter329_f2_test.go:79: hello deadline at stop+499.903833ms (sets 1); want a deadline no later than stop+450ms (listen window less one poll period)` (both rows), `--- FAIL`. `restored 390c74e841f26a12; porcelain-go=0` |
| V26 | **[r1] srt pin equality** behind V11's mux | `jq -r '.version, .cli_sha256' design_docs/verification/world-row140-m0/pin.json`; `shasum -a 256 $S/dist/cli.js` and `jq -r .version $S/package.json` (`S` as in V11); `sed -n 206,212p .github/workflows/ci.yml` | pin `0.0.78`, `3c3092bd26b3924046f38d793c716b513dde619cf792ec50b181e2c7cd40d96e`; local copy `3c3092bd26b3924046f38d793c716b513dde619cf792ec50b181e2c7cd40d96e 0.0.78`, equal. CI enforces it: `PIN=design_docs/verification/world-row140-m0/pin.json` … `[ "$GOT" = "$WANT" ] \|\| { echo "::error::srt cli.js digest drifted from pin.json"; exit 1; }` |

## 11. Quorum r1

Artifact: `design_docs/verification/world-iter329/quorum-r1.json`. **BLOCKED**: `gemini-4-argon`
PASS; `oc-glm-5-3` REJECT; `oc-kimi-k3` REJECT. Neither reject disputes the design's direction.
Both are evidence-consistency objections. The controller measured items 1–4 first-party at
`e5c1867f` before this revision, and each resolution cites that measurement.

| # | Objection (seat) | Resolution |
|---|---|---|
| 1 | §8 R2 "three packages probe" cites V19, which names two (glm, kimi) | The claim was true and the log was wrong: r0's census came from a `\| head`-truncated grep that dropped `host/daemon`. **Source: the controller's census**, re-run as V19b: `5 cmd/ailang-worldd`, `5 host/broker`, `2 host/daemon`. §8 R2 now names all three and cites V19b; V19's observed column no longer carries a census |
| 2 | "80 runs of 8 concurrent profiles" vs the logged 10 iterations × 4 cells (kimi) | **False in r0.** Source: the controller's count, 2 `ZZXP SUMMARY` lines (iters=10) in each of `out-xprof-idle.txt` and `out-xprof-load.txt`, so 40 runs. §1 row 3, §2 and R1 now say 40 (20 idle: worst 216.9 / 572.6 ms; 20 loaded: 280.5 / 933.8 ms). No second batch was run |
| 3 | V11's "= World's pin" is unlogged (glm) | True, now logged. **Source: the controller's `jq` on `pin.json`**, re-run as V26 with the local `cli.js` digest and the CI enforcement lines (`ci.yml:206-212`) |
| 4 | `:680-688` and `:736-739` are asserted without log rows (glm) | Both read as claimed. **Source: the controller's reads**, re-run as V23 (`:680-689`) and V24 (`:736-741`); §2's citations now point at them |
| 5 | T10's `listenWindow` is undefined, so a test literal could recreate the drift class (kimi ③) | §3.1's snippet now defines `listenWindow := time.Duration(execProbeListenMS) * time.Millisecond`. K9 (`execProbeListenMS = 500`) is added and **measured red**, not green as kimi's sketch expected: the ceiling becomes 450 ms < 500 ms (V25). §5 states that this red is the intended coupling |
| 6 | §3.3 step 3's mutant-applied step-8 run includes T11, so "a step-8 red" fires unconditionally (glm catch); the heading says both "before merge" and "non-blocking" (kimi ④) | Step 3 now excludes T11 from that run (or cuts the branch from the pre-M2 commit), and the decision rule's step-8 criterion excludes T11. The heading drops "before merge". One position, **non-blocking**, is stated in §3.3 and M2 |

Also fixed in passing: the header's "AC-LINUX-F3 (§4.3)" now reads §3.3.

## 12. Quorum r2 — narrow-refinement carve-out (controller-applied, iteration 329)

Quorum r2 (`design_docs/verification/world-iter329/quorum-r2.json`, $0.3281) BLOCKED with the same
seats: gemini-4-argon PASS, oc-glm-5-3 REJECT, oc-kimi-k3 REJECT, no absentees. Both rejects carry a
concrete reviewer-authored `proposed_fix` and neither disputes the design direction (glm: evidence
placement; kimi: the item-3 verdict's wording and a successor row, offering option (b) itself), so
the controller applied the reviewers' fixes under the narrow-refinement carve-out rather than a third
designer run. Each claim was checked first-party before applying.

| # | Objection (r2) | Controller check | Applied |
|---|---|---|---|
| 1 | glm: V19b was pasted into §1 row 4 (7 cells, evidence cell `V19b`), not into §10 | `grep -n '^\| 4 \|'` showed row 4 carrying the V19b cells and a stray `\| V20 \|` line; confirmed | glm fix (1) and (2) verbatim: row 4 restored with evidence `V20`; V19b inserted between V19 and V20 |
| 2 | glm catch: `:657` hit unclassified | `sed -n 657p host/broker/exec_probe.go` → a doc comment | glm fix (3): V2 and §8 classify `:657` as comment, list `:582` as P2's added use |
| 3 | glm catch: "the existing suite reds only through T6" generalised from a `-run` subset | V4's command is the reverse-leg subset; confirmed | §1 row 1 scoped to V4's subset; "whole package not run at 60 s" stated |
| 4 | kimi: item 3 closed as REFUTED with no enforcement and no successor | V18 (forced) refuses both profiles; R1 rests on statement order; confirmed | kimi option (b): verdict demoted to "refuted at HEAD by measurement; enforcement residual open"; successor row for the `net6` hoist; P4 ordering-invariant comments at `:455`/`:226`; R1 names the hoist as the successor, not a rejected option |
| 5 | kimi lesser flag: T11 pins a host-visibly inert line by string match | V5/V10/V12 as cited | §3.2 records the trade explicitly |

No production behaviour changed relative to r1; P4 is comments only.

