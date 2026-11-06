# w-exec-probe-foreign-dial-false-refusal — judge arm 6's forward legs on the child's own line, so a local loopback dialer no longer refuses startup

- **Status**: Planned (designed iteration 359, **r3**: r2 after quorum round 1 BLOCKED, then the controller's narrow-refinement carve-out applying quorum round 2's verbatim fixes; §11). Queue
  row 250, routed by D-WORLD-108 = B (attended 2026-11-06: rows 250, then 255 → 257). The row says
  it "needs its own design and quorum"; this is that design. It builds on row 242
  (`w-exec-probe-cross-profile-net6-hoist.md`, whose §6 R1/R2 first measured this residual) and
  reuses the reverse leg's rule (`w-exec-probe-reverse-leg-pins.md`, row 230: child / foreign /
  silent, foreign retried, silent sticky). **[r2]** It reuses that rule as **code**, not only as
  vocabulary: the line reader, the per-attempt judge, the read deadline and the attempt bound are
  each one shared function or constant used by both legs (§3.6). It does not change the reverse
  leg's verdicts, retry or reset semantics (row 233, parked, D-WORLD-105); §3.6 shows the
  refactor is verdict-preserving and the 12 reverse-leg tests stay green (V17).
- **Item**: Clause 3 (MET; availability hygiene on the Workspace.Exec fail-closed startup probe).
  It changes what arm 6's two forward legs count, so it is a structure change on a fail-closed gate
  and needs its own quorum.
- **Designer lane**: `claude:claude-opus-5-5` via Agent opus.
- **Estimated**: ~0.5 day. **[r2]** Production `exec_probe.go` +237/−53 lines (about half is
  comments; r1 was +197/−33, and the growth is the shared `readProbeLine` and `judgeLeg` with
  their doc comments). One new test file, `host/broker/exec_probe_net6_dial_test.go` (305 lines,
  four tests). Two existing tests need small edits: three lines in `exec_probe_test.go` and one
  regex in `exec_probe_net6_test.go`. Four names go into the real-srt verbose step's `tests=` list
  on `ci.yml:247`, the line after the one holding `TestExecSrtStartupProbePasses` (`:246`), in
  the same step (§5, V15, V16). The byte-exact prototype is
  `design_docs/verification/world-iter359/designer-proto.patch.txt` (r2, 821 lines, V13).
- **Measurement base**: `1835d8be` (= origin/dev), worktree `sprint/world-iter359-foreign-dial`.
  `exec_probe.go` sha256 prefix `19de369381579387`, which is row 242's landed prototype (V1).
- **Instruments**: `go1.26.6 darwin/arm64`. That matches `go.mod`'s `go 1.26.6`, so no toolchain
  switch was needed. Also node v26.12.1, git 2.54.0 and 16 cores (V1). The rig was shared, with load
  averages of 2.9–9.8 at the timed runs. The Sunholo start page (`serve.py`, pid 45693) was running
  throughout (V3, V4). No Linux tool, `bwrap` or real srt is available locally
  (`WORLD_EXEC_SRT_NODE_MODULES` is unset), so every Linux and real-srt statement here is
  **UNMEASURED**. Prototypes ran on sibling `git archive` copies of `1835d8be` under
  `.designer-scratch-iter359/{head,proto,applycheck}` (r1) and
  `.designer-scratch-iter359-r2/{proto,patchtree,applycheck}` (r2), never under `/tmp`, and were
  deleted afterwards. Raw outputs and scripts are banked as
  `design_docs/verification/world-iter359/designer-*.txt` (`$V` below; r2's as `designer-r2-*`).

---

## 1. What row 250 says, and what this session measured

| # | Claim | Verdict at `1835d8be` | Evidence |
|---|---|---|---|
| 1 | Arm 6 counts every accept on its `net6` listener and requires exactly 1 (the control) after each leg | **CONFIRMED.** The accept loop increments `accepts` and closes without reading (`exec_probe.go:651`). The legs require `acceptsA != 1` → fail (`:552`) and `acceptsC != 1` → fail (`:556`). The counts are cumulative: `acceptsC` also includes Run A's accepts | V2 |
| 2 | A foreign dial between bind and `acceptsC` refuses the profile; it is never a false pass | **CONFIRMED, now deterministically.** A dial from the test process during the raw leg's run gives `accepts 2` on both legs. A dial during the proxy leg's run gives `accepts 2` on the proxy leg (HEAD positive control, 3 of 3 foreign rows red). Every shape of extra accept refuses at HEAD (child line, silent, connect-close), so HEAD is fail-closed | V8 |
| 3 | The rig's start page dials every loopback listener about every 31.8 s | **CONFIRMED again today, with a new fact: it SENDS BYTES.** Over 70 s, 8 idle listeners each got 2 accepts. On every one (16 of 16), the first read returned `GET / HTTP/1.1\r\nHost: 127.0.0.1:<port>\r\nUser-Agent: SunholoHome/1` within 0–2 ms. The source is `services.py:59 probe()`, which sends an HTTP GET with `User-Agent: SunholoHome/1.0`. This fact decides the design: the rig's one known dialer is **positively identifiable as foreign by its bytes** | V3, V4 |
| 4 | No current test pins the foreign-accept case | **CONFIRMED.** The only test text about arm-6 accept counts is the pass-through test's `accepts 2` / `accepts 3` (`exec_probe_test.go:204–205`, the known positive). No test dials `net6` | V5 |
| 5 | The designer's 2–6 % false-refusal estimate per profile per probe | **Still UNMEASURED directly.** The window/period derivation in row 242's §6 R1 stands. This session saw the effect live once: an unsolicited dial landed inside a prototype test's window, and the prototype classified it **foreign** (V10) | row 242 V15; V10 |

## 2. The mechanism being changed

HEAD's arm 6 has three legs on World's live `net6` listener, plus the reverse leg:

- **control** (`control()`, `:659–674`): World dials unsandboxed, then waits for `accepts == 1`.
- **raw** (Run A, `connect <port>`, `:500`): the child dials and prints `WORLD-PROBE-CONNECTED` or
  `WORLD-PROBE-REFUSED <errno>`.
- **proxy** (Run C, `proxy <port>`, `:584`): the child sends an HTTP CONNECT through srt's proxy
  environment, or dials directly when there is none.

Each leg fails unless its token is a refusal **and** the host-side count is still 1. The count is
World's own evidence. It does not trust the child: a proxy that reaches the listener while the child
prints a refusal still fails. That property is the reason the token alone is not enough (§3.1).
But the count cannot tell *who* dialed, so any local process that dials the port inside the
0.56–1.79 s window (row 242 V15) refuses the profile.

## 3. Design

### 3.1 The direction: nonce-classified accepts, with a foreign-only window retried

The row offers two directions.

- **(a) "Judge the leg on the token the child prints"** is **rejected**. It drops the host-side
  evidence. A sandbox whose proxy connects upstream to `127.0.0.1:<port>` and then answers the child
  with a non-200, or tunnels and then closes, makes the child print `WORLD-PROBE-REFUSED HTTP-…`
  or `…ECLOSED`. Today that fails on `accepts`; under (a) it would pass. That is a fail-open change.
- **(b) "Count only accepts that carry a probe-issued nonce"** is **chosen, in its fail-closed
  form**. Counting *only* nonce accepts would be the mirror-image hole: every accept without the
  nonce would be ignored, including a silent proxy reach. So the design counts **every** accept and
  classifies each one by the first line it carries. Only a **provably** foreign accept stops
  counting as a reach, and even that is never a pass by itself: it makes the window inconclusive,
  and the leg runs again (§3.3). This is the reverse leg's rule from row 230, transposed: a foreign
  answer is retried, a silent one refuses for good, and the child always loses.

### 3.2 Classification (`classifyNet6Dial`, new)

Every probe client dial now writes one line as soon as it connects:
`WORLD-PROBE-DIAL <child nonce>\n`. The raw op writes it on `connect`. The proxy op writes it
through the tunnel after a 200 (`s.end(...)`, where HEAD calls `s.destroy()`). World's own control
dial writes `WORLD-PROBE-DIAL <control nonce>\n`. Both nonces are 16 random bytes in hex, made
per listener and so per profile per startup. The child nonce travels in the op argument
`"<port>:<nonce>"`, as the reverse leg's nonce does. The control nonce never leaves World.

**[r2]** On each accept, World reads at most one line's worth of bytes, one byte at a time, within
`execProbeHelloWait` = 500 ms, through the reverse leg's own reader, now factored out as
`readProbeLine` (§3.6). It stops at `\n`, at the longest expected line, at the deadline, or at the
first byte that begins neither expected line:

| Bytes read | Class | Counted as | Why |
|---|---|---|---|
| exactly the control line | control | — (World's own; before every window) | proves the listener accepts **and** the classifier reads |
| exactly the child line | **child** | a reach: the leg **fails at once** | the child, or a relay of the child's bytes, reached World |
| a byte that cannot begin either line (decided at that byte) | **foreign** | the window is **inconclusive**: retry | positive evidence that another protocol is speaking. The start page's `GET /` is foreign at its first byte (V3) |
| nothing: EOF, reset or deadline | **silent** | a reach: the leg **fails at once** | **the subtle case.** A proxy that connects upstream and then relays nothing, a child whose write is blocked, and a port scanner all look the same. Nothing proves it is not the child, so it fails closed |
| a strict prefix of a line, then EOF or deadline | **silent** | a reach: the leg **fails at once** | a truncated child line is not evidence of a foreign dialer |
| an accept still pending at `window()`'s bound, `execProbeWindowWait` (1.5 s) | **silent** | a reach | unreachable by construction: the bound exceeds the classifier's own read deadline by a second (the §3.5 invariant, pinned by TC), so a pending accept there is a stalled classifier. It fails closed if it ever happens |

A whole line that differs from both lines (for example a forged `WORLD-PROBE-DIAL 000…`) diverges at
its first wrong byte, so it is foreign, as the reverse leg treats a forged nonce. A relay that
*prepends* bytes (for example a PROXY-protocol header) also reads as foreign. §4 shows why that
still cannot pass.

### 3.3 Windows and the leg verdict (`net6Leg`, new)

The listener records each accept's class **in accept order**. `mark()` returns the accept count so
far, and `window(from)` waits until every accept since `from` is classified, then returns the counts
and where the next window starts. Windows are contiguous: the raw leg's starts at the `mark()` taken
right after `control()` returns, and each later window starts where the previous one ended. So every
accept from Run A's spawn until the proxy leg's last sample lies in exactly one window. Run B, the
reverse leg and every retry are covered, as HEAD's cumulative `acceptsC` covered them.

`net6Leg` judges one leg, attempt by attempt. **[r2]** The decision is the shared `judgeLeg`
(§3.6), called with the window's counts as booleans, the token test "is a refusal", and no
retried token; `judgeReverse` is the same function with the reverse leg's token test. So:

1. If the token is not `WORLD-PROBE-REFUSED …`, or the window has any **child** or **silent**
   accept (`judgeLeg` returns `revReached`, `revNoToken` or `revSilent`), the leg **fails**. The text names the token, the attempt and the counts:
   `…: token "…" on attempt N; accepts in its window: child C, silent S, foreign F (want a refusal
   and no accept from the child or silent; <run>)`.
2. If the window has no accept at all, the leg **passes**.
3. Otherwise every accept in the window was foreign, so the window is **inconclusive**. Unless this
   was attempt `execProbeAttempts` (3; **[r2]** the one bound both legs use), the leg re-runs its own op alone, with a fresh window:
   `connect <port>:<nonce>` with the trampoline, or `proxy <port>:<nonce>` without it, as in Run A
   and Run C. After the third inconclusive window the leg **fails**:
   `…: inconclusive after 3 attempts: every attempt's window held an accept that was not the
   child's (foreign F1, F2, F3)`.

`probe()` calls `net6Leg` for the raw leg right after Run A's 200 ms settle, before arms 1–4 are
judged and before Run B. That keeps the windows in order. It calls `net6Leg` for the proxy leg
after `reverseLeg` and its settle. The leg texts join the arm-6 failure exactly as HEAD's do, so
the control, raw, proxy and reverse legs are each still named.

### 3.4 The control stays meaningful

`control()` now writes the control line and waits up to 2 s for an accept **classified as the
control**. HEAD waited for `accepts == 1`. So the control proves three things: the listener
accepts, World can read a line from an accepted dial, and the classifier recognises an exact line.
A foreign accept that lands before the control (row 242 V16's window) no longer satisfies it, and it
no longer refuses anything: it falls before the raw window's `mark()`. That is provably not the
child, because no child of this probe exists before Run A is spawned (statement order in `probe`).

**[r3]** The 2 s is HEAD's control bound, unchanged (V20): HEAD's `control()` dials with
`net.DialTimeout(…, 2*time.Second)` and waits `for deadline := time.Now().Add(2 * time.Second)`
(`exec_probe.go:663`, `:668`), and the patch keeps that deadline line byte-for-byte as context. So
r3 neither shortens nor lengthens the control's wait, and the control cannot become a new source of
the false-refusal class this item removes.

### 3.5 The patch (the planner pins these bytes: V13)

`exec_probe.go`, by function:

- **`execProbeScript`**: `connect(port,nonce)` writes `dial(nonce)` with `s.end` on connect.
  `proxy(port,nonce)` writes it through a 200 and passes the nonce through to `connect` when there is
  no proxy. Both ops parse `"<port>:<nonce>"`. The doc comment says so.
- **[r2] Constants.** One new line constant, `execProbeDialPrefix` (`"WORLD-PROBE-DIAL "`). r1's
  `execProbeNet6Wait` and `execProbeNet6Attempts` are **gone**; each leg uses the shared one:
  - `execProbeHelloWait` (500 ms, unchanged value and name) is the deadline of **every** probe line
    read (`readProbeLine`). Its comment gains the forward use: there a late line reads as silent,
    which refuses, so the wait can delay a pass but never make one. Its row-230 floor (250 ms, a
    legitimate answer must be admitted) and ceiling (950 ms, pinned by
    `TestExecProbeReverseHelloWaitBounds`) still bind it. The name is kept because that test, the
    reverse leg's refusal text and row 230's design all cite it.
  - `execProbeAttempts` (3) replaces `execProbeReverseAttempts` (renamed; no test names it, V18) and
    bounds both legs' retries, so R2's odds cannot drift between the legs.
  - **new** `execProbeWindowWait = execProbeHelloWait + time.Second` (1.5 s), the bound in
    `window()` that r1 wrote inline. Reason: each classifier sets its read deadline when its
    goroutine starts, after the accept, so the extra second covers that goroutine's scheduling
    delay on a loaded rig (V1: loads up to 9.8). **Invariant: `execProbeWindowWait >
    execProbeHelloWait`.** It holds by the definition, and TC's first line asserts it, so a later
    edit that sets the bound to the deadline or below goes red (MUT-WINDOW-WAIT-NOT-ABOVE, §5).
    Under the invariant, an accept still pending at the bound is a stalled classifier, never a
    read that was still allowed to finish, and `window` counts it silent.
- **`probe`**: `from := net6.mark()` after `control()`. Run A's argv carries `net6.arg()`. The
  `acceptsA` and `acceptsC` reads and their two `if` blocks are replaced by
  `rawWhy := pp.net6Leg(…"connect", true…)` and `proxyWhy := pp.net6Leg(…"proxy", false…)`, appended
  to `legs` in the same order.
- **`reverseLeg(ctx, proxyArg string)`**: it took `proxyPort int`. Its only callers are `probe` and
  `runReverseLeg` in the test (V6).
- **`probeListener`**: the `accepts atomic.Int64` field becomes `ctl`, `child`, `mu` and `seen`.
  `newProbeListener` makes both nonces (a `randomHex` failure is carried in `err` and judged by the
  control leg, as a bind failure is) and classifies each accept in its own goroutine. New methods
  `arg`, `mark` and `window`; `control` is rewritten as in §3.4; `close` is unchanged.
- **New** `net6Class`, `net6Window`, `net6Leg` and `classifyNet6Dial`.
- **[r2] Shared with the reverse leg** (§3.6): new `readProbeLine`, which `classifyReverseDial`
  and `classifyNet6Dial` both call; new `judgeLeg`, which `judgeReverse` (now a one-line wrapper)
  and `net6Leg` both call.

`ProbeExecProfilesEach`, the row-242 hoist, the reverse leg's verdicts, the arm order and every
other arm are unchanged. The `newProbeListener` call site stays where T2 of row 242 pins it.

### 3.6 [r2] Reuse analysis against the reverse leg

The reverse leg in the same file already classifies loopback dials as child / foreign / silent and
retries a foreign attempt. r1 wrote a second stack beside it. r2 goes through it piece by piece.
A piece is **shared** unless there is a precise protocol difference.

| Reverse-leg piece (HEAD) | Forward-leg need | r2 |
|---|---|---|
| The read loop in `classifyReverseDial`: deadline, one byte at a time, stop at `\n` or at the line's length | the same: read one line from a loopback peer within a deadline | **Shared: `readProbeLine(c, hello, wants...)`.** Both classifiers call it. The reverse leg passes `execProbeHello` and one want. The forward leg passes `""` (no hello) and two wants (control, child) |
| The hello write before the read | none | **Not shared, and cannot be: a protocol difference.** On the reverse leg World is the dialer and must speak first, since the child's listener answers the hello. On the forward legs World is the listener and the child dials and speaks first. A hello would add a write into a relay that is not owed one. It is a parameter of the shared reader, not a second reader |
| `execProbeHelloWait` (500 ms) | a read deadline | **Shared.** r1's separate `execProbeNet6Wait` (also 500 ms) is deleted. A change to the deadline now reaches both legs, because there is one constant read inside the one reader. Its **value** is pinned by the reverse leg's floor/ceiling test (T10). The forward tests scale with it rather than pin it: TC's hold and elapsed bounds are written in terms of it. The forward leg's own bound on it is the §3.5 invariant. MUT-SHARED-WAIT-1500MS is red on T10 and green on the forward tests (§5), as this predicts |
| `execProbeReverseAttempts` (3) | a retry bound | **Shared**, renamed `execProbeAttempts`. The forward bound had no independent measurement: r1 chose 3 by the reverse leg's precedent. MUT-SHARED-ATTEMPTS-2 is red from both legs' tests |
| `judgeReverse(tok, seen)`: child wins, then a bad token, then silent, then foreign or EADDRINUSE retries, else pass | the same order over the window's classes | **Shared: `judgeLeg(tokOK, seen, retryTok)`.** `judgeReverse` keeps its signature (its table test `TestExecProbeReverseLegJudge` is unchanged) and becomes `judgeLeg(tok is BOUND or a refusal, seen, tok is EADDRINUSE)`. `net6Leg` calls `judgeLeg(tok is a refusal, {child>0, foreign>0, silent>0}, false)`. The two token tests differ because the ops differ: `listen` succeeds with `BOUND`, and its `EADDRINUSE` means another process held the fresh port. `connect`/`proxy` have no good token but a refusal, and no fresh port to retry onto. MUT-SHARED-JUDGE-SILENT-RETRIED is red from both legs' tests |
| The class decision after the read: exact line → child; **any** other bytes → foreign; zero bytes then a clean EOF after a written hello → foreign; else silent | exact control or child line → that class; bytes that begin neither line → foreign; **no bytes, EOF, or a strict prefix** → silent | **Not shared: two precise differences**, so each classifier keeps its own few-line `switch`. (i) **A clean EOF with no bytes.** On the reverse leg it follows World's hello. The child's listener always answers (`c.end('WORLD-PROBE-NONCE …')`), so a peer that closes cleanly is not the child. On the forward legs a connect-then-close is exactly HEAD's client and a relay that reaches and relays nothing. It may be the child, so it must be silent (TS `raw-scan`; MUT-SILENT-AS-FOREIGN). (ii) **A strict prefix.** On the forward legs the child's line can cross srt's proxy tunnel, which can cut it, so a prefix is silent (TC's two prefix rows; MUT-PREFIX-AS-FOREIGN). The reverse leg's foreign-on-any-bytes is row 230's ratified rule, and this item does not change reverse-leg verdicts (row 233, D-WORLD-105) |
| `reverseLeg`'s attempt loop | an attempt loop | **Not factored into a generic driver.** Only the bound and the judge are common, and both are now shared. The per-attempt work differs: the reverse leg picks a fresh port and nonce, spawns `listen` and polls; a forward leg re-runs `connect`/`proxy` against the **one** live listener and reads a window of its accept log. The failure texts differ too. A driver would take nearly the whole body as a closure and share one `for` line |
| `reverseSeen` and `revVerdict` | the same booleans and verdicts | **Shared as types** (`judgeLeg` takes and returns them). They keep their names so no row-230 test changes. `net6Window` keeps integer counts only for the failure text |

**The refactor does not change any reverse-leg verdict.** `readProbeLine` keeps `classifyReverseDial`'s
order (set the deadline, write the hello, read) and its bound (`len(want)+1` = `len(want+"\n")`).
It adds one early stop, at the first byte at which what was read begins no want. After that byte
the buffer is non-empty and can never equal `want+"\n"`, so `classifyReverseDial`'s `len(buf) > 0
→ foreign` branch gives the verdict it gave at HEAD, only sooner. The EOF branch needs zero bytes,
so the early stop cannot reach it. `judgeReverse` is the HEAD `switch` with its two token
expressions passed as arguments. All 12 reverse-leg tests on the `ci.yml:247` list pass on the r2
prototype (V17).

## 4. The by-construction argument

**Claim A (no new false pass).** Arm 6's raw or proxy leg passes only on an attempt whose window
held **no accept at all** and whose token is a refusal. That is HEAD's pass condition (`accepts`
still 1 after the control, with a refusal token), applied to that attempt's window.

*Proof.* **[r2]** `net6Leg` returns `""` in exactly one branch, `judgeLeg(...) == revPass`.
`judgeLeg` returns `revPass` only after its four earlier cases fail: `seen.child`, `!tokOK`,
`seen.silent`, `seen.foreign || retryTok`. `net6Leg` passes `retryTok = false` and
`seen = {child > 0, foreign > 0, silent > 0}`, so a pass means the token is a refusal and
`child`, `silent` and `foreign` are all 0. Every accept in the window was counted as one of the
three: `window` maps control and pending to silent. ∎

**Claim B (every escape HEAD fails, the design fails).** Take a sandbox that lets the child, or a
relay acting for it, reach `net6` on a leg.

- If the relay delivers the child's line, the accept is **child**, and the leg fails on that
  attempt.
- If it delivers nothing, or a prefix of the line, the accept is **silent**, and the leg fails on
  that attempt.
- If it delivers other bytes first (it rewrites or prepends), the accept is **foreign**, and the
  window is inconclusive. The leg then re-runs. A **deterministic** escape reaches again on every
  attempt, so after 3 attempts it fails "inconclusive".

The only escape that can pass is one that is **both** byte-altering **and** intermittent, reaching
on some attempts and not on the last one. HEAD passes any intermittent escape on a lucky single run
too, so this is not a new class. It does narrow the odds: an altering escape that reaches with
probability p per run passes with probability 1 − p³ against HEAD's 1 − p. This is the one place
the design is weaker than HEAD, and it is stated as R2. No known srt path alters tunnel bytes. srt's
CONNECT proxy relays bytes verbatim, so this case is reasoned, not measured.

**Claim C (the foreign dial no longer refuses).** A dialer that sends a non-line byte within
500 ms, and whose dials do not land in all 3 of a leg's consecutive windows, cannot refuse arm 6.
The start page sends `G` within 2 ms (V3). Its period is about 31.8 s. A leg's three windows together last
about 1–2.5 s (Run A or Run C, plus two retries of about 0.3 s each: V9; row 242 V15), so on this rig
the expected residual from it is 0. That figure is
reasoned from the measured period; no startup was counted.

## 5. Test plan and mutation table

New file `host/broker/exec_probe_net6_dial_test.go`. Every dial comes **from the test process**, so
the tests do not depend on the rig's start page. The client, `net6DialClient` (a `/bin/sh` script
in place of World's node), answers every op as a confining sandbox would. On each `connect` or
`proxy` op it writes the op's argument to `op-N` and **blocks until the test writes `go-N`**. The
test's goroutine sees `op-N`, makes the planned dial(s) to the port, with the nonce from the
argument when the shape is `child`, and then releases the client. Each dial therefore lands while
that op's run is in flight: inside its window, deterministically. The pass-through shim runs the
client unsandboxed, so every other arm passes and `err == nil` means arm 6 passed.

- **TF** `TestExecProbeNet6ForeignDialIsRetriedNotRefused`. Its rows are `none`, `raw`, `proxy`
  and `raw-twice-proxy` (two foreign dials on the first two raw attempts plus one on the proxy).
  Each dial is the start page's own `GET` shape. Arm 6 must pass, and each op must have run **at
  least** the planned number of times (1 + the number of foreign windows) and at most 3. A rig
  dialer can add a retry but never remove one (V10 is why the bound is `≥`). **RED at HEAD**
  (`accepts 2`, V8): this is the row's pin.
- **TS** `TestExecProbeNet6ChildOrSilentAcceptStillFails`. Rows: `raw-child`, `proxy-child` (a reach
  that the child's own token denies), `raw-silent` (held 700 ms with no bytes), `raw-scan` (connect
  and close), `proxy-silent`, `raw-foreign-and-child` (both in one window) and `raw-child-on-retry`
  (foreign on attempt 1, child on attempt 2). Each must fail arm 6 with the exact
  `on attempt N; accepts in its window: child C, silent S,` prefix and run its op exactly N times.
  A child or silent window is never retried.
- **TE** `TestExecProbeNet6ForeignEveryAttemptRefuses`. Foreign on all 3 raw attempts must refuse
  with `inconclusive after 3 attempts`, after exactly 3 runs.
- **TC** `TestExecProbeNet6DialClassifier`. Pure, over a real loopback pair. It covers ten byte
  shapes: control, child, child then more, an HTTP GET, a wrong nonce, a PROXY preamble then the
  child line (foreign), connect-close, held silent, a prefix then close, and a held prefix. It also
  checks that classification never outlives its deadline. **[r2]** Its first line asserts the §3.5
  invariant `execProbeWindowWait > execProbeHelloWait`.
- **[r2] The shared pieces are pinned from both call sites** by existing tests, with no new test
  name. The reverse side is pinned by `TestExecProbeReverseLegJudge` (`judgeLeg` through
  `judgeReverse`), `TestExecProbeReverseHelloWaitBounds` (`readProbeLine` and `execProbeHelloWait`
  through `classifyReverseDial`) and the reverse-leg exhaustion tests (`execProbeAttempts`). The
  forward side is pinned by TS, TC, `TP` and TE. The `MUT-SHARED-*` rows of the table below show
  the reader, the judge and the attempt bound going red from both sides. The deadline's value goes
  red from the reverse side only, as §3.6 explains.
- **Edited** `TestExecProbePassThroughShimFailsEveryConfinementArm` (`TP`). Its two leg strings
  become `token "WORLD-PROBE-CONNECTED" on attempt 1; accepts in its window: child 1, silent 0,`
  for the raw and the proxy legs. This pins, through World's real node client, that an unsandboxed
  child's dial is classified **child** (not silent or foreign) on both legs. It also removes the
  test's `accepts 4` flake: an unsolicited accept is now foreign, and foreign is not in `TP`'s
  strings.
- **Edited** `TestExecProbeNet6BoundBeforeAnyPick`. Its client-log regex becomes
  `` (?m) connect (\d+):[0-9a-f]{32}$ ``, because the connect argument now carries the nonce
  (V7).
- **Edited** `runReverseLeg`: `reverseLeg(ctx, "1:"+strings.Repeat("0", 32))`.
- **[r2]** `ci.yml`: TF, TS, TE and TC are appended to the **third** `tests=` line of the step
  "Row 140 Workspace.Exec real-srt tests, verbose (linux bubblewrap)" (`ci.yml:238–260`). That line
  is `:247`, which already ends with row 242's `TestExecProbeNet6BoundBeforeAnyPick` and
  `TestExecProbeNet6BindSiteIsBeforeTheFanOut`. The step's four `tests=` lines (`:245–248`) build
  **one** list: `:246` holds `TestExecSrtStartupProbePasses` (the real-srt startup probe) and `TP`,
  and `:247` holds the 12 reverse-leg tests. The step fails unless `WORLD_EXEC_SRT_NODE_MODULES` is
  set. It runs the whole list with `go test ./host/broker/ ./host/daemon/ -count=1 -p 1 -v -run …`,
  then requires `--- PASS: <t>` for **every** name, so a SKIP fails the step (V15, V16). So placing
  the four names on `:247` puts them on the same real-srt verbose list as
  `TestExecSrtStartupProbePasses`; the step asserts --- PASS for every name, so a SKIP fails it —
  their Linux PASS is measured only when the §7 gate runs **[r3, glm's verbatim text]**. Each name
  appears once; no other list is needed, because this step's `go test` is the CI run that is
  checked by name.

**[r2]** The mutation table was re-measured on the r2 prototype by `$V/designer-r2-run-mutants.sh.txt`
and `$V/designer-r2-mutate.py.txt` → `$V/designer-r2-mutants.txt` (V19; loads 3.2–9.4). Each mutant
was applied to a pristine copy of the r2 prototype's `exec_probe.go`, the tests were run with
`-count=1`, and the copy was restored and checked with `cmp`. The r1 mutants were re-anchored on
r2's code (the judge is now `judgeLeg`, the classifier's tail a `switch` on `s`) and ran against
`TP`, TF, TS, TE and TC. The new **shared-piece** mutants and CONTROL also ran against the 12
reverse-leg tests (column **Rev**). 1 = that test went red. MUT-NO-NONCE is dropped: r1 measured it
degenerating to MUT-PREFIX-AS-FOREIGN. The HEAD row is r1's measurement (V8). r2 changed none of
TF/TS/TE's assertions, only the bound's identifier in TF.

| Mutant | Edit | `TP` | TF | TS | TE | TC | Rev | Expected → observed |
|---|---|---|---|---|---|---|---|---|
| CONTROL | none | 0 | 0 | 0 | 0 | 0 | 0 (12 of 12 green) | green → green (17 of 17) |
| HEAD (positive control) | `1835d8be` + TF/TS/TE (TC needs new identifiers) | — | **1** (`accepts 2`) | **1** (texts only: HEAD refuses every row, fail-closed) | **1** | — | — | TF red → red (V8) |
| **MUT-SHARED-READ-SHORT** [r2] | `readProbeLine` reads one byte less than the longest want | **1** | **1** (all 4 rows: the control line is never whole) | **1** (4 child rows) | 0 | **1** (`control`, `child`, `child-then-more`) | **1** (`ReverseHelloWaitBounds/slow-child`: `seen {child:false foreign:true …}`; `ReverseLegRealReachSurvivesAForeignAttempt`) | red from **both** legs → red from both |
| **MUT-SHARED-JUDGE-SILENT-RETRIED** [r2] | `judgeLeg`'s `seen.silent` → `revInconclusive` | 0 | 0 | **1** (`raw-silent`, `raw-scan`, `proxy-silent`) | 0 | 0 | **1** (`ReverseLegJudge`, `ReverseLegSilentHolderRefusesSticky` ×3 rows, `ForeignListenerOnTheReversePortDoesNotRefuse/silent`) | red from **both** legs → red from both |
| **MUT-SHARED-ATTEMPTS-2** [r2] | `execProbeAttempts` = 2 | 0 | **1** (`raw-twice-proxy`) | 0 | **1** | 0 | **1** (`ReverseLegExhaustionRefusesDistinctly`: `chatty-every-time`, `child-eaddrinuse`) | red from **both** legs → red from both |
| **MUT-SHARED-WAIT-1500MS** [r2] | `execProbeHelloWait` = 1500 ms | 0 | 0 | 0 | 0 | 0 | **1** (`ReverseHelloWaitBounds` both rows; `ForeignListenerOnTheReversePortDoesNotRefuse/silent`) | red from the reverse leg's ceiling only → as predicted (§3.6: the forward tests scale with the constant, so the value is pinned once, by T10) |
| **MUT-WINDOW-WAIT-NOT-ABOVE** [r2] | `execProbeWindowWait = execProbeHelloWait` | 0 | 0 | 0 | 0 | **1** (`execProbeWindowWait 500ms must exceed the line-read deadline execProbeHelloWait 500ms`) | 0 | red (the §3.5 invariant) → red |
| MUT-RAW-COUNT | the leg fails on **any** accept (HEAD's rule inside the new structure) | 0 | **1** | **1** (`raw-child-on-retry`) | **1** | 0 | — | red → red |
| MUT-SILENT-AS-FOREIGN | the classifier's last `return net6Silent` → `net6Foreign` | 0 | 0 | **1** (`raw-silent`, `raw-scan`, `proxy-silent`) | 0 | **1** (4 rows) | — | red → red |
| MUT-PREFIX-AS-FOREIGN | a non-empty strict prefix → foreign | 0 | 0 | 0 | 0 | **1** (2 prefix rows) | — | red (TC only) → red |
| MUT-CHILD-AS-FOREIGN | the exact child line → foreign | **1** | 0 | **1** (4 child rows) | 0 | **1** | — | red → red |
| MUT-CHILD-AS-CONTROL | the exact child line → control | **1** | 0 | **1** (texts: arm 6 still refuses, with `silent 1`, 11 of 11) | 0 | **1** | — | red → red; still fail-closed by `window`'s mapping |
| MUT-TOKEN-ONLY | `net6Leg` ignores the window (the row's direction (a)) | **1** | **1** | **1** (all 7) | **1** | 0 | — | red → red |
| MUT-RETRY-CHILD | a child window is retried like a foreign one | 0 | 0 | **1** (4 child rows: arm 6 = `""`, **a false pass**) | 0 | 0 | — | red → red |
| MUT-RETRY-SILENT | a silent window is retried | 0 | 0 | **1** (3 silent rows) | 0 | 0 | — | red → red |
| MUT-EXHAUST-PASSES | the third inconclusive window passes | 0 | 0 | 0 | **1** | 0 | — | red → red |
| MUT-UNBOUNDED | no attempt bound | 0 | 0 | 0 | **1** | 0 | — | red → red |
| MUT-FROM-NOT-ADVANCED | windows never advance (`*from` kept) | **1** | **1** | 0 | 0 | 0 | — | red → red |
| MUT-WINDOW-AT-JUDGE | the window starts at judgement time, not before the run | **1** | **1** | **1** (all 7) | **1** | 0 | — | red → red |
| MUT-CLIENT-NO-LINE | the client's `connect` destroys instead of writing its line | **1** | 0 | 0 | 0 | 0 | — | red (`TP`'s `child 1`) → red |
| MUT-PENDING-AS-FOREIGN | `window` counts a still-pending accept as foreign | 0 | 0 | 0 | 0 | 0 | — | **SURVIVES.** Unreachable by construction: `window`'s bound `execProbeWindowWait` (1.5 s) exceeds the classifier's read deadline `execProbeHelloWait` (500 ms), and **[r2]** that invariant is now pinned by TC (MUT-WINDOW-WAIT-NOT-ABOVE), so nothing is pending at the bound. Killing it would need a seam that stalls the classifier; none is added |
| MUT-CONTROL-ANY-ACCEPT | `control()` returns on any accept, not on its own line | 0 | 0 | 0 | 0 | 0 | — | **SURVIVES.** It weakens only the control's proof that the classifier reads a line. That proof is pinned independently: `TP`'s `child 1, silent 0` (an unsandboxed child classified child through World's real node, on both legs) and TC. Under the mutant, a foreign accept between the hoisted bind and `control()` (row 242 V15: milliseconds) would satisfy the control. That still proves the listener accepts; only the classifier-read proof is lost, and `TP` and TC pin that proof |

Every r1 mutant reproduced r1's kill pattern on r2, test for test (V19 against V11). Every mutant
that would turn a reach into a pass (MUT-RETRY-CHILD, MUT-RETRY-SILENT, MUT-SILENT-AS-FOREIGN,
MUT-TOKEN-ONLY, MUT-EXHAUST-PASSES, MUT-UNBOUNDED, MUT-WINDOW-AT-JUDGE,
MUT-SHARED-JUDGE-SILENT-RETRIED) is killed by TS or TE. Every mutant that would bring the false
refusal back (MUT-RAW-COUNT, MUT-FROM-NOT-ADVANCED, MUT-TOKEN-ONLY, MUT-WINDOW-AT-JUDGE) is killed by
TF. **[r2]** The shared reader, the shared judge and the shared attempt bound are each red from a
reverse-leg test **and** a forward-leg test, so a fix to either leg's use reaches the other and is
checked there. The two survivors (MUT-PENDING-AS-FOREIGN, MUT-CONTROL-ANY-ACCEPT) survive again on
r2. They are stated above with their reasons, and no seam is added to kill them.

## 6. Residuals and risks

- **R1, the remaining false refusals (stated, not fixed).** A foreign dialer still refuses the
  profile if, inside a leg's window, it:
  - sends **nothing** within 500 ms (a port scan such as `nc -z`, a TCP health check, a client that
    waits for a banner);
  - sends a prefix of `WORLD-PROBE-DIAL ` and stalls; or
  - lands in all 3 consecutive windows of one leg, which needs a dialing period shorter than about
    0.5–0.8 s.

  Each of these is refused **closed**, with a text naming the silent accept or the inconclusive
  exhaustion, never as a pass. On this rig the one known dialer sends a `GET` within 2 ms (V3), so
  its contribution goes from row 242's estimated 2–6 % per profile per probe to 0 by Claim C.
  Silent dialers on other hosts are **UNMEASURED**.
- **R2, the one narrowing (Claim B).** A byte-altering **and** intermittent escape now has three
  tries to land a clean window instead of one: it passes with probability 1 − p³ where HEAD's
  passes with 1 − p. This is the one way the design is weaker than HEAD. No srt path is known to alter tunnel bytes, and a
  deterministic one still refuses (TE). The alternative, failing at once on foreign, would keep the
  start page's false refusals. The design takes the reverse leg's precedent (row 230: foreign →
  retry, bounded) over that.
  **[r3] The byte-fidelity premise is now verified, not asserted (V21).** The controller read the
  PINNED srt 0.0.78 source (`npm pack`, `cli.js` sha256 equal to `pin.json`): a CONNECT is
  TLS-terminated only when `network.tlsTerminate` is configured and the first bytes are a TLS
  ClientHello, routed to a MITM socket only when `network.mitmProxy` is configured, and sent to a
  parent proxy only when `network.parentProxy` is configured; otherwise it dials direct, writes the
  200, forwards `head`, and `pipe`s both directions with no transform (`http-proxy.js:355–372`; the
  SOCKS path is likewise `pipe` both ways). World's srt settings `network` struct has exactly four
  fields (`allowedDomains`, `deniedDomains`, `allowLocalBinding`, `allowAllUnixSockets`) and no
  `.go` file outside tests names any of those three keys (0; control `allowedDomains` 1). So under
  World's configuration the pinned srt relays tunnel bytes verbatim, and the byte-altering class R2
  describes has no instance in the pinned srt. R2 stays stated because a future srt pin or a
  configuration change could create one; §7's executor step re-checks it, and its fallback is the
  fail-at-first-foreign rule. This is a static source read, not a dynamic Linux measurement.
- **R3, latency.** A clean probe adds no wait: `window()` returns as soon as nothing is pending, and
  a foreign accept is classified at its first byte. Each inconclusive window adds one small run plus
  the 200 ms settle. That measured about 0.3 s per retry (`raw` 1.50 s against `none` 1.21 s, V9).
  A silent accept costs up to 500 ms before the leg refuses.
- **R4, real srt and Linux (UNMEASURED here).** Under real srt the raw connect is refused and the
  proxy's CONNECT to `127.0.0.1` is denied, so no accept lands and the window is empty: the same
  verdict as HEAD. The client's new `s.end(line)` only runs after a successful connect.
  **[r2]** `TestExecSrtStartupProbePasses` is already on the real-srt verbose list, on `ci.yml:246`
  (V15, the controller's measurement; r1 said `:247`, which is the next line of the same list). It
  is the first real-srt run of this patch, and the step fails if it SKIPs, so the patch cannot land
  green with the new client line unmeasured under real srt. It skips locally (`bwrap` is absent
  and `WORLD_EXEC_SRT_NODE_MODULES` is unset). That is why the Linux step is an acceptance gate in
  §7.
- **R5, a nonce in argv.** Another local process could read the child nonce from `ps` and send the
  child line, which only makes arm 6 **refuse** (child). It cannot make it pass. The control nonce is
  never in any argv.

## 7. Milestones

- **M1 (single, ~0.5 day).** Apply `$V/designer-proto.patch.txt`. `git apply` reproduces all five
  files byte-exact (V13). Then prove, each test run on its own:
  - TF is red at the parent commit (the HEAD row, V8) and green on the patch;
  - the mutants in §5 reproduce their kill pattern: at minimum MUT-RAW-COUNT, MUT-SILENT-AS-FOREIGN,
    MUT-TOKEN-ONLY, MUT-RETRY-CHILD, MUT-RETRY-SILENT and MUT-EXHAUST-PASSES, and **[r2]** the
    shared-piece mutants MUT-SHARED-READ-SHORT and MUT-SHARED-JUDGE-SILENT-RETRIED, each red from
    **both** a reverse-leg test and a forward-leg test;
  - the edited `TP` and `TestExecProbeNet6BoundBeforeAnyPick` are green;
  - **[r2]** the 12 reverse-leg tests on `ci.yml:247` are green on the patch (the refactor is
    verdict-preserving, §3.6).

  Then run `./scripts/verify_ail.sh` (no `.ail` changes) and `go build ./... && go vet ./host/broker/
  && go test ./...` under the pinned `AILANG_BIN`.

  **[r3] Stability, on the applied patch (gemini, verbatim):** run `go test -count=3 -run
  '^(TestExecProbePassThroughShimFailsEveryConfinementArm|TestExecProbeNet6.*)$' -v ./host/broker/`
  on the applied r2 patch so the multi-run stability proof from V12 is verified on the refactored
  code (the controller already ran it on the r2 prototype with the reverse-leg tests added: V22).

  **[r3] Executor verification (blocking; kimi's fix, verbatim):** with the pinned srt install,
  grep/inspect the proxy's CONNECT handler and record in the sprint Verification Log (a) whether
  bytes after a 200 are relayed verbatim in both directions, (b) whether any header/banner is
  prepended or appended upstream. If any alteration exists, the design reverts to
  fail-at-first-foreign (no retry on windows containing byte-bearing accepts; retry only nothing),
  i.e. reclassify any accept carrying bytes that begin neither expected line as silent. The
  controller's own reading of the pinned 0.0.78 source (V21) found no alteration; the executor
  confirms it against the pin the patch is built with and records the result.

  **[r2] Acceptance gate: the Linux real-srt CI step.** Real srt and `bwrap` are absent locally, so
  the item is **not accepted** until the PR's CI step "Row 140 Workspace.Exec real-srt tests, verbose
  (linux bubblewrap)" is green. That step asserts `--- PASS` for `TestExecSrtStartupProbePasses`
  (`:246`, the patched probe client under real srt) and for TF, TS, TE and TC (`:247`). A SKIP of any
  of them fails the step. The executor and evaluator record that step's result, with the run URL,
  in the sprint's Verification Log. A local green is not a substitute.

## 8. Files

| File | Change |
|---|---|
| `host/broker/exec_probe.go` | **[r2]** +237/−53: the client script's `connect`/`proxy`; `execProbeDialPrefix` and `execProbeWindowWait`; `execProbeReverseAttempts` → `execProbeAttempts`; `execProbeHelloWait`'s comment; `probe`'s arm-6 legs; the `reverseLeg` parameter; `probeListener` rewritten (nonces, classified accepts, `arg`/`mark`/`window`/`control`); new `net6Class`, `net6Window`, `net6Leg`, `classifyNet6Dial`; shared `readProbeLine` (out of `classifyReverseDial`) and `judgeLeg` (out of `judgeReverse`) |
| `host/broker/exec_probe_net6_dial_test.go` | new, 305 lines: TF, TS, TE, TC (whose first line pins the §3.5 invariant) and the `net6DialClient` rig |
| `host/broker/exec_probe_test.go` | +3/−3: `TP`'s two leg strings; `runReverseLeg`'s argument. No reverse-leg test changes |
| `host/broker/exec_probe_net6_test.go` | +1/−1: the connect-argument regex |
| `.github/workflows/ci.yml` | +1/−1: four names appended to the real-srt verbose step's `tests=` line `:247` (the list that also holds `TestExecSrtStartupProbePasses` on `:246`) |

## 9. What this item is NOT doing

It does not change the reverse leg's behaviour (verdicts, retry, reset, hello wait: rows 230, 232
and 233). **[r2]** It does refactor the reverse leg's code: the read loop moves into the shared
`readProbeLine`, the judge into `judgeLeg`, and the attempt constant is renamed. §3.6 shows this is
verdict-preserving and V17 shows the 12 reverse-leg tests green. It does not change the row-242
hoist, the arm order, the arm-6 failure texts of the control and reverse legs, or any other arm. It
does not make arm 6 tolerate silent foreign dialers (R1), and it adds no seam to production: the
tests use the existing `cfg.Client`. It does not touch `.ail` code, `tools/launchd/*` or the start
page.

## 10. Verification Log

`$V` = `design_docs/verification/world-iter359`. The scratch copies are
`.designer-scratch-iter359/{head,proto,applycheck}`.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Base, instruments, toolchain | `git rev-parse --short HEAD`; `shasum -a 256 host/broker/exec_probe.go`; `go version` against `grep '^go ' go.mod`; node, git, cores, `uptime`; `pgrep -fl homepage/serve.py`; srt and `bwrap` presence → `$V/designer-instruments.txt`, `$V/designer-greps.txt` | `1835d8be`; `19de369381579387` (row 242's landed proto); `go1.26.6` = `go 1.26.6` (no toolchain switch); node v26.12.1, git 2.54.0, 16 cores; `serve.py` pids 45688/45693; `WORLD_EXEC_SRT_NODE_MODULES` unset, `bwrap-absent` |
| V2 | HEAD's arm 6 counts every accept, cumulatively, and the client destroys on connect | `grep -n "accepts" host/broker/exec_probe.go`; `grep -n "s.destroy();r(.WORLD-PROBE-CONNECTED…"`; function-site grep → `$V/designer-greps.txt` | `:503 acceptsA := net6.accepts.Load()`, `:544 acceptsC := …`, `:552 … acceptsA != 1`, `:556 … acceptsC != 1`, `:651 l.accepts.Add(1)` (then `c.Close()`, no read), `:669 if l.accepts.Load() == 1` (control); `:500 "connect", strconv.Itoa(net6.port)`; script `:81` and `:89` call `s.destroy()` on connect and on the proxy's status line |
| V3 | The start page **sends a GET** on every dial | 8 idle listeners, 70 s, first read within 500 ms logged (`$V/designer-foreign-bytes.go.txt`) → `$V/designer-foreign-bytes.txt`; `grep -n "User-Agent\|def probe\|connection.request" services.py` → `$V/designer-greps.txt` | `TOTAL accepts=16` (2 per listener, at 18:00:05.959 and 18:00:37.781, a 31.8 s period); 16 of 16 first reads were `"GET / HTTP/1.1\r\nHost: 127.0.0.1:<port>\r\nUser-Agent: SunholoHome/1"`, in 0–2 ms; `services.py:59 def probe`, `:76 connection.request("GET", …`, `:77 "User-Agent": "SunholoHome/1.0"` |
| V4 | The dialer is still live today (the controller's premise) | `$V/controller-premise.txt` (controller, 16:53Z) | `unsolicited_accepts_per_listener=[2, 2, 2, 2, 2, 2, 2, 2]` in 75 s; `pgrep … wc -l` = 2 |
| V5 | No test pins a foreign accept on `net6` | `grep -rn "accepts [0-9]" host/ --include=*_test.go`; `grep -rln "newProbeListener\|probeListener" host/ --include=*_test.go`; `grep -rn "WORLD-PROBE-DIAL" host/ \| wc -l` with `WORLD-PROBE-NONCE` as the control → `$V/designer-greps.txt` | only `exec_probe_test.go:204–205` (the TP strings, the known positive); only `exec_probe_net6_test.go` (row 242's source pin); `WORLD-PROBE-DIAL` 0 against `WORLD-PROBE-NONCE` 8 |
| V6 | The changed functions' callers | `grep -rn "reverseLeg(\|\.control()\|newProbeListener(" host/ --include=*.go` → `$V/designer-greps.txt` | `reverseLeg`: `exec_probe.go:542` and `exec_probe_test.go:580` only; `control()`: `:493` only; `newProbeListener`: `:257` only |
| V7 | Row 242's T1 parses the connect argument, so it must change | `grep -n 'connect\|Atoi' host/broker/exec_probe_net6_test.go` | `:110 regexp.MustCompile(`(?m) connect (\d+)$`)` |
| V8 | HEAD positive control: TF is red at HEAD; HEAD refuses every TS shape | `$V/designer-head-positive-control.txt`: HEAD copy + TF/TS/TE (pre-V14 assertions), `go test -count=1 -run '^TestExecProbeNet6(ForeignDial\|ChildOrSilent\|ForeignEvery)' -v`, load 3.83 | TF `none` PASS; `raw`, `proxy` and `raw-twice-proxy` FAIL: `arm 6 refused on foreign dials alone: a raw connect … accepts 2 …; a connect through srt's proxy environment … accepts 2` (proxy row: proxy leg `accepts 2` only). TS: all 7 rows refuse at HEAD (`accepts 2`/`3`, red on text only). TE red (`accepts 2 … accepts 3`, one run) |
| V9 | The prototype is green; per-retry cost | proto copy, `go test -count=1 -run '^TestExecProbeNet6(ForeignDial\|ChildOrSilent\|ForeignEvery\|DialClassifier)' -v` (first run, before V14); `TP`, no-op and missing-client `-count=3` | all PASS; TF `none` 1.21 s, `raw` 1.50 s, `proxy` 1.47 s, `raw-twice-proxy` 2.12 s; ops logged as `connect <port>:<32 hex>`, the retry visible (`connect` ×2); `TP` 3 of 3 PASS with the new strings |
| V10 | A **live** rig dial inside a prototype window is classified foreign | the mutation script's first CONTROL run, load 9.80 → `$V/designer-rig-dial-in-window.txt` | `proxy-child`: `accepts in its window: child 1, silent 0, foreign 1`, with only one test dial (`child:ok`). The unsolicited accept read as foreign, and the arm still refused on the child. This is why V14 relaxed TS/TE's exact `foreign 0` and TF's exact run counts |
| V11 | Mutation table | `$V/designer-run-mutants.sh.txt` + `$V/designer-mutate.py.txt` → `$V/designer-mutants.txt` (summary; loads 1.9–7.5) | as §5; `# restored proto exec_probe.go (cmp ok)`; MUT-RETRY-CHILD's false pass verbatim: `arm 6 = "", want it to contain "a raw connect … on attempt 1; accepts in its window: child 1, silent 0,"` |
| V12 | Stability and the arm-6 subset on the **r1** prototype (r2 re-ran its tests once, `-count=1`, V17; the `-count=3` stability run was not repeated on r2) | `go test -count=3 -run '^(TestExecProbePassThroughShimFailsEveryConfinementArm\|TestExecProbeNet6.*)$' -v`; `go test -count=1 -run ExecProbe ./host/broker/`; `go vet ./host/broker/` → `$V/designer-stability.txt` | 3/3 PASS for each of the 7 tests (including row 242's T1/T2 with the new regex); `ok … 70.415s`; ExecProbe subset `ok … 46.682s`; `vet-ok` |
| V13 | **[r2]** The patch is byte-exact | r2 patch made by `git diff` in a `git archive 1835d8be` tree (`patchtree`) holding the r2 prototype's five files. Then `git apply $V/designer-proto.patch.txt` on a fresh `git archive` of `1835d8be` (`applycheck`), and `cmp` of every file against the prototype → `$V/designer-r2-apply.txt` | `APPLY-OK`; `cmp` same ×5; sha256 prefixes `exec_probe.go` `110fffc3f31a10d7`, `exec_probe_test.go` `ee59ba2cc5742855`, `exec_probe_net6_test.go` `8ec970891599257a`, `exec_probe_net6_dial_test.go` `8c7f9602c1105bdc` (305 lines), `ci.yml` `4efb18438d736bb4` (unchanged from r1); patch 821 lines, sha256 `7fb2537268498db0`; numstat `exec_probe.go` +237/−53, dial test +305, `exec_probe_test.go` +3/−3, `exec_probe_net6_test.go` +1/−1, `ci.yml` +1/−1. (r1: `43bfd807641f4fd1`, `f4f98eb35246b6f8`, patch `4793fa7d2d108932`, +197/−33) |
| V14 | The tests the mutants ran are the patch's tests (r1; **[r2]** the r2 mutants ran on `proto`, whose three test files `cmp` same as the r2 patch tree, V13) | `cmp` of the proto copy's three test files against the patch tree | `same` ×3. The assertion relaxation (V10) came before the final mutation run and the patch |
| V15 | **[r2, replaces r1's V15]** `TestExecSrtStartupProbePasses` is on the real-srt verbose list, and that step asserts PASS per name (glm's objection; controller's measurement at `4721da68`) | `$V/controller-r1-srt-ci.txt`: `grep -n -o "TestExecSrtStartupProbePasses" .github/workflows/ci.yml`; control `grep -n -o TestExecProbeNet6BoundBeforeAnyPick …`; `grep -nE "bwrap\|bubblewrap\|sandbox-runtime\|WORLD_EXEC_SRT\|srt" …`; `grep -n "func TestExecSrtStartupProbePasses" -r host` | `246:TestExecSrtStartupProbePasses`; control `247:TestExecProbeNet6BoundBeforeAnyPick`; `:238 - name: Row 140 Workspace.Exec real-srt tests, verbose (linux bubblewrap)`; `:243 [ -n "${WORLD_EXEC_SRT_NODE_MODULES:-}" ] \|\| { echo "::error::WORLD_EXEC_SRT_NODE_MODULES is unset"; exit 1; }`; `:226` exports the pinned srt install into `$GITHUB_ENV`; `host/broker/exec_probe_srt_test.go:25`. The step body (`ci.yml:238–260`, read in full this round) runs `go test ./host/broker/ ./host/daemon/ -count=1 -p 1 -timeout 600s -v -run "^(${tests// /\|})\$"`, then `grep -q -- "^--- PASS: $t "` for every name (`::error::$t did not PASS on linux (SKIP?)`). r1's "`ci.yml:247` list" was the right step and list but the wrong line |
| V16 | **[r2]** Where the four new names land | in `applycheck` (the r2 patch applied): `grep -n -o -E "TestExecProbeNet6(ForeignDialIsRetriedNotRefused\|ChildOrSilentAcceptStillFails\|ForeignEveryAttemptRefuses\|DialClassifier)\|TestExecSrtStartupProbePasses\|TestExecProbePassThroughShimFailsEveryConfinementArm\|TestExecProbeNet6BoundBeforeAnyPick" .github/workflows/ci.yml`; `grep -c TestExecProbeNet6DialClassifier …` → `$V/designer-r2-ci-placement.txt` | `246:TestExecSrtStartupProbePasses`, `246:TestExecProbePassThroughShimFailsEveryConfinementArm`, `247:TestExecProbeNet6BoundBeforeAnyPick` (controls), then `247:` ×4 for TF, TS, TE, TC; count 1 (one list); the step's `- name:` and four `tests=` lines plus the `go test` and `PASS: $t` lines |
| V17 | **[r2]** The r2 prototype is green, including all 12 reverse-leg tests (the refactor keeps their verdicts) | r2 proto copy (load 1.17), `go test -count=1 -timeout 400s -v -run '^(<the 12 reverse-leg names on ci.yml:247>\|TestExecProbePassThroughShimFailsEveryConfinementArm\|TestExecProbeNet6.*)$' ./host/broker/`; `go vet ./host/broker/`; `gofmt -l host/broker/` → `$V/designer-r2-green.txt` | 19 of 19 `--- PASS` (TF, TS, TE, TC, row 242's T1/T2, `TP`, and `TestExecProbeReverseLegJudge`, `…ChattyForeignHolderIsRetried`, `…ExhaustionRefusesDistinctly`, `…NoTokenIsNotRetried`, `…SilentHolderRefusesSticky`, `…RealReachSurvivesAForeignAttempt`, `TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse`, `TestExecProbeReverseDialErrIsSilent`, `…PollRecordsConnectReset`, `…HelloWaitBounds`, `…PollCadence`, `TestExecProbeChildListenerDrainsTheHello`); `ok … 35.461s`; `VET-OK`; `gofmt -l` empty |
| V18 | **[r2]** The shared names and their call sites; the rename touches no test | `$V/designer-r2-greps.txt`: at `1835d8be`, `grep -rn --include=*_test.go execProbeReverseAttempts host \| wc -l` with control `… execProbeHelloWait host`; `grep … "judgeReverse(\|classifyReverseDial(" host`; in `applycheck`, `grep -n "readProbeLine(\|judgeLeg(\|execProbeAttempts\|execProbeWindowWait\|execProbeNet6Wait\|execProbeNet6Attempts\|execProbeReverseAttempts" …` | `0` against the control's 2 hits (`exec_probe_test.go:668`, `:889`); `judgeReverse` called at `:610` and by its table test `exec_probe_test.go:452`. Patched tree: `readProbeLine` called at `:838` (forward) and `:899` (reverse); `judgeLeg` at `:177` (through `judgeReverse`) and `:675` (`net6Leg`); `execProbeAttempts` at `:618` (reverse) and `:683` (forward); `execProbeWindowWait` `:157`, used at `:772`, asserted at dial test `:254`. `execProbeNet6Wait`, `execProbeNet6Attempts` and `execProbeReverseAttempts` give **no** hits in the same call that matches the others |
| V19 | **[r2]** The mutation table on r2, including the shared-piece mutants | `$V/designer-r2-run-mutants.sh.txt` + `$V/designer-r2-mutate.py.txt` → `$V/designer-r2-mutants.txt` (loads 3.2–9.4); r2's edits to the r1 prototype are `$V/designer-r2-edit.py.txt` | as §5. CONTROL `pass=17 fail=0`; `MUT-SHARED-READ-SHORT` `pass=11 fail=6` (both legs); `MUT-SHARED-JUDGE-SILENT-RETRIED` `pass=13 fail=4` (both legs); `MUT-SHARED-ATTEMPTS-2` `pass=14 fail=3` (both legs); `MUT-SHARED-WAIT-1500MS` `pass=15 fail=2` (reverse only); `MUT-WINDOW-WAIT-NOT-ABOVE` `fail=1` (TC); the two survivors `rc=0 pass=5 fail=0`; `# restored r2 exec_probe.go (cmp ok)`; MUT-RETRY-CHILD's false pass verbatim: `arm 6 = "", want it to contain "a raw connect to World's live listener: token \"WORLD-PROBE-REFUSED ECONNREFUSED\" on attempt 1; accepts in its window: child 1, silent 0,"` |
| V20 | **[r3]** HEAD `control()` wait bound (glm r2) | at `1835d8be`: `sed -n 659,674p host/broker/exec_probe.go`; `grep -n "2 \* time.Second" $V/designer-proto.patch.txt` | `c, err := net.DialTimeout("tcp", l.ln.Addr().String(), 2*time.Second)` (`:663`) and `for deadline := time.Now().Add(2 * time.Second); …` (`:668`); in the r2 patch that deadline line appears only as an unchanged context line (patch line 373). HEAD's bound is 2 s and r3 keeps it |
| V21 | **[r3]** The pinned srt relays CONNECT-tunnel bytes verbatim under World's config (kimi r2) | `$V/controller-r2-srt-relay.txt`: `npm pack @anthropic-ai/sandbox-runtime@0.0.78`, extracted without install scripts, `cli.js` sha256 checked against `pin.json`; `sed` of `http-proxy.js:262–268, 300–306, 355–372`, `mux-proxy.js:72–75`; `grep -nE "network\.(tlsTerminate\|mitmProxy\|parentProxy)" sandbox-manager.js`; `sed -n 44,49p host/broker/exec_settings.go`; `grep -rnE "tlsTerminate\|mitmProxy\|parentProxy" --include=*.go host cmd \| grep -v _test.go \| wc -l` with control `allowedDomains` | digest equal to the pin; termination/MITM/parent-proxy routes are each gated on a config key; otherwise `upstream.write(head)`, `upstream.pipe(socket)`, `socket.pipe(upstream)`; World's `execSettingsNet` has 4 fields, none of the three keys; grep 0, control 1. Static read; not run on Linux |
| V22 | **[r3]** `-count=3` stability on the **r2** prototype, with the reverse-leg tests (gemini r2, kimi r2) | `$V/controller-r2-stability.txt`: r2 patch applied to `git archive 1835d8be`, go1.26.6, load 3.88; `go test -count=3 -timeout 600s -v -run "^(TestExecProbePassThroughShimFailsEveryConfinementArm\|TestExecProbeNet6.*\|TestExecProbeReverseLeg.*)$" ./host/broker/` | rc 0; 13 tests × 3 `--- PASS` each (TF, TS, TE, TC, row 242's T1/T2, TP, and 6 `TestExecProbeReverseLeg*`). The HEAD positive control on r2's test file is MUT-RAW-COUNT in V19 (HEAD's raw counting under r2's tests: TF red on `raw`, `proxy`, `raw-twice-proxy`; TS `raw-child-on-retry` red; TE red) |
| — | **UNMEASURED** | — | the full `go test ./...` under the pinned `AILANG_BIN` (the executor runs it; r2 ran `go build ./...` and `go vet ./host/broker/` only, V17); Linux; real srt (R4: the §7 acceptance gate is the CI step); whether TF/TS/TE/TC PASS rather than SKIP on the CI runner (the step asserts it); silent foreign dialers on other hosts (R1); a startup-level false-refusal count, before or after; r2's HEAD positive control was not re-run (r1's V8 stands, and r2 changed no TF/TS/TE assertion) |

## 11. Quorum log

**Round 1 (r1, `4721da68`): BLOCKED, $0.3073.** Seats: gemini-4-argon **PASS**, oc-glm-5-3
**REJECT**, oc-kimi-k3 **REJECT**. claude was benched as the author's vendor; no seat was absent.

- **glm: the real-srt premise was unlogged.** R4 said `TestExecSrtStartupProbePasses` is "already
  on the `ci.yml:247` list" and is the patch's first real-srt run. No Verification Log row
  supported that. V15's grep (`TestExecProbeNet6`) could not match that name, and its excerpt
  elided the middle of the line. If the premise were false, the new `s.end(line)` tunnel write
  would never be measured. **r2's answer:** the controller measured it
  (`$V/controller-r1-srt-ci.txt`). The premise is true, but the line is `:246`, not `:247`, in the
  same step and the same list, and the step asserts `--- PASS` per name. V15 now cites that
  measurement with its commands and outputs, and V16 records where the four new names land in the
  patched tree (`:247`, once each, in that same PASS-asserted list). R4, §5 and §8 name the exact
  line for each test. §7 makes the Linux real-srt step an explicit acceptance gate.
- **kimi: a parallel stack with no reuse analysis.** r1 transplanted the reverse leg's rule but
  wrote new types and a new classifier beside the reverse leg's. It had independently chosen
  constants: a 500 ms read deadline, and an unnamed 1.5 s `window()` bound that was absent from
  §3.5. That is a silent-divergence hazard on a fail-closed gate. **r2's answer** (§3.6) is a
  piece-by-piece reuse analysis, and the prototype is refactored to match:
  - **Shared:** the line reader (`readProbeLine`, out of `classifyReverseDial`), the judge
    (`judgeLeg`, out of `judgeReverse`), the read deadline (`execProbeHelloWait`; r1's
    `execProbeNet6Wait` is deleted) and the attempt bound (`execProbeAttempts`; r1's
    `execProbeNet6Attempts` is deleted).
  - **Not shared, each for a stated protocol difference:** the hello write (a parameter of the one
    reader), the post-read class decision (clean EOF and strict prefix differ by direction) and the
    attempt loop's body.
  - The 1.5 s bound is now `execProbeWindowWait`, with its reason, and the invariant
    `execProbeWindowWait > execProbeHelloWait` is pinned by TC.
  - The 12 reverse-leg tests stay green (V17). The `MUT-SHARED-*` mutants show the reader, the
    judge and the attempt bound red from both legs' tests. The deadline's value is red from the
    reverse leg's T10, which pins it for both (V19).
- **Non-blocking, folded in:** gemini's two surviving mutants stay stated with their reasons (§5);
  gemini's "verify `TestExecSrtStartupProbePasses` on CI" is now §7's acceptance gate; the
  controller's R2 note is kept plain (§4 Claim B, §6 R2: 1 − p³ against HEAD's 1 − p).

**Round 2 (r2, `50afe67e`): BLOCKED, $0.3735.** Same seats: gemini-4-argon **PASS**, oc-glm-5-3
**PASS**, oc-kimi-k3 **REJECT**; claude benched; no seat absent.

- **kimi (blocking): R2 is a quantified weakening (1 − p³ vs 1 − p) resting on an unverified
  premise** that no srt path alters tunnel bytes; the §7 CI step cannot see that class. Proposed
  fix: an executor verification of the pinned srt's CONNECT relay with a fail-at-first-foreign
  fallback, plus re-running the `-count=3` stability run and the HEAD positive control on r2 code.
  Its catch: "Verify, don't assert, the byte-fidelity premise … absent that, the quorum must
  explicitly ratify the regression."
- **glm (pass, closest concern):** the control's 2 s wait had no Verification Log row and might
  differ from HEAD's; and §5 stated the Linux PASS in the present tense.
- **gemini (pass):** re-run `-count=3` on the refactored code in M1.

**r3 — the controller's narrow-refinement carve-out** (gate-2 rule: after the one re-quorum, every
remaining blocking objection carries a concrete reviewer-authored fix and none disputes the
design's direction; kimi's objection is about verifying a premise, and its own text says
verification answers it). Applied, verbatim where the reviewers wrote text, after running the
measurements they named:
- kimi: the controller **verified** the premise first-party (V21: the pinned srt relays verbatim
  under World's config; static source read). kimi's executor step and fallback are in §7 verbatim;
  §6 R2 now cites V21 and keeps the 1 − p³ statement. The `-count=3` run on r2 code is V22 (13
  tests × 3 PASS), and the HEAD positive control on r2's tests is MUT-RAW-COUNT (V19, cited in V22).
- glm: V20 records HEAD's 2 s control bound, and §3.4 says it is unchanged (glm's option (a)); §5's
  sentence is replaced with glm's text.
- gemini: §7 carries its `-count=3` command verbatim.

No design direction changed and the prototype patch is unchanged from r2
(`designer-proto.patch.txt`, sha256 prefix `7fb2537268498db0`).
