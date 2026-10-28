# w-exec-probe-reverse-leg-port-race — arm 6's reverse leg counts ANY listener on its port as the child's, so a port another host process takes refuses startup

- **Status**: Planned
- **Item**: iteration 314, controller-filed from the red CI run on `origin/dev` `2769cff9`
  (run 40699055042 attempt 1, job 122853572046, step 8). Clause 3 (the Workspace.Exec effect's
  fail-closed startup probe), as an availability defect.
- **Estimated**: ~0.5–0.75 day [r3]. About +95/−20 LOC in `host/broker/exec_probe.go`, about
  +260 LOC of tests, a 3-line QUICKSTART prose edit and a 1-line CI test-list edit (§10).
- **Measurement base**: `2769cff9` (worktree `sprint/world-iter314-arm6-reverse-leg`, clean
  tree), 2026-10-28.
- **Instruments**: local `go1.26.6 darwin/arm64` and node v26.12.1 (`/opt/homebrew/bin/node`).
  CI logs were read with `gh api` (V3, V4). srt 0.0.78 source was read from an installed copy
  whose `cli.js` sha256 equals World's pin (V7). **No real srt runs locally**:
  `WORLD_EXEC_SRT_NODE_MODULES` is unset, so every real-srt test SKIPs by name (V10). The Linux
  claims rest on CI logs, on M0's committed Linux matrix and on reading srt's source. None of
  them was re-run on Linux in this session.
- **Files touched**: `host/broker/exec_probe.go` (the probe script's `listen` op, `pollLoopback`,
  a new `reverseLeg` and `judgeReverse`, and an unexported port-picker seam);
  `host/broker/exec_probe_test.go` (seven new tests, T1–T7 [r3]); `docs/QUICKSTART.md` §10 arm 6
  prose (S7); `.github/workflows/ci.yml` step 8 test list (all seven new tests added, so any SKIP
  fails CI [r3]).
  **No `.ail` file is touched** (V15). `host/daemon` is unchanged: the failing daemon test was
  the victim, and no daemon test asserts reverse-leg text (V14).

---

## 1. Problem, and what the controller's hypothesis got right and wrong

`TestExecSrtMCPGrantAndBudget` (host/daemon) went red on Linux CI because `New` refused the
daemon at the exec startup probe:

```
exec startup probe failed for profile "proj" (1 of 7 arms): arm6-network: reverse leg: the host
reached the listener the child bound on 127.0.0.1:44109
```

The same tree passed when the job was re-run (V3, V4). Only one instance of this signature has
been observed against a real sandbox (controller scan of 8 failed runs). Every other
"reached the listener" line in CI logs comes from deliberate pass-through cases: the
`t.Logf` at `exec_startup_test.go:499` (V20) and the assertion at `exec_probe_test.go:197`.

| Controller claim | Status | Evidence |
|---|---|---|
| Intermittent, not a regression: parents green, attempt 2 green on the identical tree | **CONFIRMED** | V3, V4. `exec_probe.go` is byte-unchanged since `885af2d1` (V1) |
| On Linux the child binds inside its own netns, so the host cannot have reached the *child* | **CONFIRMED, with a stronger in-run witness.** M0 measured `BOUND; host reached=False` on Linux and `REFUSED EPERM` on macOS (V9). The failing run itself also reported **only** the reverse leg ("1 of 7 arms", one leg). So the raw connect and the proxy connect to World's counted listener were both refused, with accepts still 1, in that same run. The network namespace held while the host "reached" something on 127.0.0.1:44109 | V3, V2 (legs are judged independently and joined; exec_probe.go:466-488) |
| A TOCTOU port race: another host-namespace process bound `revPort` between `freeLoopbackPort`'s Close and the poll | **CONSISTENT, and a concrete in-window binder is identified.** srt 0.0.78 starts its mux proxy (HTTP+SOCKS) on **`127.0.0.1` port 0 in the host namespace** on every non-Windows start (V7). Run C's own srt therefore binds one ephemeral host port *inside the poll window*. That listener accepts and sends nothing, which matches the shape the poll counts as a reach. The binder in the CI instance is not proven, because the log cannot name it. The fix in §4 does not depend on which binder it was | V7, V11 |
| "`go test ./...` runs package binaries in parallel, which makes a collision plausible" | **REFUTED for the observed instance.** Step 8 runs `go test ./host/broker/ ./host/daemon/ -p 1`, so the two package binaries run one after the other. Neither package calls `t.Parallel()` (count 0; the same grep finds 3 files elsewhere in the repo). No other test binary was binding ports in that step. The in-window binders that remain are srt itself (V7), the test process's own goroutines, and the runner's other processes | V5, V6 |
| The same false refusal can happen in production | **CONFIRMED by reading.** `Run` builds the daemon under a 9 s startup context, and the probe is fail-closed. A failed `--exec-profile` refuses startup. A failed directory profile is silently left out ("excluded" with the probe's reason). Production uses the same srt with the same mux bind | V16, V18, V7 |
| The probe client's listener sends nothing on accept | **CONFIRMED.** `net.createServer(c=>c.destroy())` | V2 |

TCP self-connect was considered as a third mechanism and ruled out as the primary cause. Go's
TCP dialer detects a self-connect and retries it (V8). The fix below would classify one
correctly anyway, because a self-connected socket never returns the child's nonce.

## 2. The mechanism

Today the reverse leg (exec_probe.go:455-485, V2) works like this:

1. `revPort := freeLoopbackPort()`: listen on `127.0.0.1:0`, read the port, **close it**.
2. `poll := pollLoopback(revPort)` dials `127.0.0.1:revPort` every 50 ms until stopped. **Any
   dial that connects** sets `reached=true`.
3. Run C (srt, no trampoline) runs the client's `proxy` op, then `listen revPort:1000`.
4. If the token is BOUND or REFUSED and `reached` is true, arm 6 fails with "the host reached
   the listener the child bound".

The port is free only at the instant of step 1. Everything after that treats a connect on that
number as proof that it was the child's listener. On Linux the child's listener is in another
namespace and *cannot* be what the host reached. On macOS the child's bind is refused, so there
is no child listener to reach at all. On a correct sandbox, then, **every** reach is a
non-child listener, and today's code turns each one into a refusal.

**Deterministic reproduction at base (V11).** A stand-in srt passes every run through except
Run C (the argv carrying `listen`). For Run C it binds a silent listener on that port, which is
the shape srt's own mux has, and prints the Linux-correct tokens `WORLD-PROBE-REFUSED
ECONNREFUSED` / `WORLD-PROBE-BOUND` **without running the child**. With that stand-in,
`ProbeExecProfiles` reported `reverse leg: the host reached the listener the child bound on
127.0.0.1:60489` **3 of 3 times on darwin**. The paired control, the same stand-in without the
bind, never reported a reverse leg (0 of 3). The listener reached was provably not the child's,
because the child never ran. This is the CI signature, produced on demand without srt timing.

**[r3] The reproduction now uses a CHATTY stand-in.** V11's stand-in was a *silent* listener,
`net.createServer(c=>{})`. Under the r3 design (§4.2) a holder that stays silent after World's
hello is a **silent holder**, and a silent holder refuses by design. So V11's shape no longer
reproduces the *flake*; it now becomes the sticky-refusal test (T6 `silent` row, T7). srt's mux
is not silent once spoken to. Sent `WORLD-PROBE-HELLO\n`, it answers `HTTP/1.1 400 Bad Request`
and closes (V23). The deterministic reproduction of the CI flake therefore uses a chatty
stand-in shaped like the mux: it says nothing until it receives bytes, then answers HTTP 400 and
closes. At base that stand-in also produced `reverse leg: the host reached the listener the child
bound on 127.0.0.1:<port>`, **3 of 3 times on darwin** (V25), because today's poll never speaks
first.

**A sibling hazard, not fixed here (R2).** `ProbeExecProfilesEach` probes profiles concurrently
(V19). Profile X's poll can land on profile Y's *live counted* listener if Y binds it after X
picked `revPort`. In that case X's reach is foreign, which this fix handles, but Y's accept count
reaches 2 and Y's raw or proxy leg refuses falsely. Identity on X's side cannot undo Y's count.
It is a different leg and needs ≥2 profiles; §11 names it as a successor.

## 3. Options

The property that must not weaken is V46/P20: **a listener the child binds is unreachable from
the host**, so a real reach must always refuse. The bug is that "something answered on the port"
is being used as a stand-in for "the child's listener answered".

**(A) Child identity on accept, a hello discriminator, and a bounded retry on a fresh port.
CHOSEN [r3].** The probe client's listener writes `WORLD-PROBE-NONCE <nonce>\n` on accept, with a
fresh 128-bit nonce for each attempt passed through argv. On every connected dial the host first
writes the fixed hello `WORLD-PROBE-HELLO\n`, then reads with a 500 ms deadline. A dial is
classified as one of three things:
- **child**: it read exactly the nonce line. This is a real reach and refuses.
- **foreign**: it read any other bytes, or got an immediate clean EOF. This is retryable.
- **silent holder**: it read zero bytes after the hello, either at the deadline or on a reset.
  This refuses, and the refusal is **sticky across the leg**: one silent holder on any attempt
  refuses the leg with a distinct message, whatever a later attempt would show.

An attempt with a foreign dial and no child or silent dial is INCONCLUSIVE. It is retried on a
fresh port, up to 3 attempts in all. If every attempt is inconclusive, the leg refuses with a
second distinct message.

**Sensitivity [r3].** Today any connected dial refuses. After r3, the only connected dials that
do not refuse are those whose peer *actively answered the hello* with non-nonce bytes or a clean
close. World's own probe client never answers that way: its listener writes the nonce on
accept, or, if that answer is slow or lost, the dial is silent and refuses. The probe client is
World's verified node running a fixed script (row 140 §4.6). So for a listener the child binds,
the leg is **exactly as sensitive to a real breach as it is today**. Retrying applies only to
listeners that prove they are not the child by answering the hello: srt's mux answers HTTP 400
(V23).

**(B) Bounded retry, no identity.** This treats every reach as suspect and retries, refusing only
if all 3 attempts reach. It gets three things wrong:
- **It weakens detection.** A breach that is reachable on attempt 1 but not on attempt 2 (an
  intermittently reachable child listener, or a breach that loses the 50 ms poll race inside a
  1 s listen window) would *pass*. Today one reach refuses; (B) needs three.
- **It cannot tell the operator** whether the refusal was a breach or three collisions. Both
  print the same text.
- **It cannot tell a real reach on a collided port from a collision.** With a shared namespace,
  a foreign holder makes the child's bind fail EADDRINUSE. (B) then sees "reached, refused",
  the same as a correct macOS sandbox.

**(C1) World holds `revPort` itself for the whole of Run C, so nothing else can take it.** This
closes the race completely, but it stops measuring the property. The host's dial would always
reach World's own listener, so the signal has to move to the child's bind errno (BOUND means a
separate namespace, EADDRINUSE means a shared one). That inverts V46 from a reachability
observation into an inference. It also depends on macOS reporting Seatbelt's EPERM before an
in-use port's EADDRINUSE, which is **not measured** (M0 bound to an unheld port). And a
`SO_REUSEPORT` bind would co-bind silently. Rejected: it trades a direct observation for an
unmeasured one.

**(C2) The child binds port 0 and prints the port; the host dials once after BOUND.** This is
M0's own shape (V9). It narrows the window but does not close it, because the child's port
number can still be held by a host process. It also needs Run C's stdout read mid-run, and
`runBounded` captures only until exit (exec_probe.go:286-289). Rejected: more plumbing, and the
race remains.

**(C3) Pick `revPort` outside the ephemeral range.** No `bind(0)` on the host, srt's mux
included, can then land on it. But it needs per-platform sysctl reads, it still collides with
fixed-port services, and like (B) it still cannot classify a reach. It could be combined with
(A) later, but it adds nothing to (A)'s correctness. Rejected.

(A) is the only option that *classifies* a reach. Every other option either keeps "any reach is
the child's" (B, C2, C3) or stops observing reachability at all (C1).

## 4. The chosen design (one milestone, independently landable)

### 4.1 The probe client (the `execProbeScript` constant)

The `listen` op takes `<port>:<ms>:<nonce>`, and its server ends each accepted connection with
the nonce line:

```js
function listen(port,ms,nonce){return new Promise(r=>{const v=net.createServer(c=>{c.resume();c.end('WORLD-PROBE-NONCE '+nonce+'\n')});
// …error/listen/BOUND/close unchanged…
else if(op==='listen'){const q=x.split(':');t=await listen(q[0],q[1],q[2])}
```

The nonce comes from the existing `randomHex(16)` (crypto/rand, exec_probe.go:311). It is fresh
for every attempt, so a listener left over from attempt *n* cannot answer for attempt *n+1*.
`c.resume()` [r3] drains the host's hello, so unread input cannot turn the child's close into a
reset. A hello-then-nonce exchange against a node listener of exactly this shape returned the
nonce 500/500 times, the slowest under 1 ms on darwin (V26).

### 4.2 The host side

- **`pollLoopback(port, want string) func() reverseSeen` [r3].** On each connected dial it first
  writes `WORLD-PROBE-HELLO\n`, then sets a 500 ms read deadline and reads at most
  `len(want)+1` bytes up to `\n`. It classifies the dial as follows:
  - the line equals `want+"\n"`: set `seen.child`;
  - ≥1 byte that is not that line, or a clean EOF with zero bytes: set `seen.foreign`;
  - zero bytes at the deadline, or a reset with zero bytes: set `seen.silent`.

  The reset case follows the reviewer's "zero bytes after the hello" branch. It also covers a
  child listener closing with a dial still queued in its backlog, which can arrive as a reset.
  Such a dial therefore refuses; it is never retried. Dial cadence and timeout are unchanged
  (50 ms, 200 ms). Measured: the plain nonce read answered 500/500 dials, the slowest under 1 ms
  (V13), and so did the hello-then-nonce exchange (V26), both on darwin. srt's mux, sent the
  hello, answered `HTTP/1.1 400 Bad Request` and closed within 0–3 ms, 3 of 3 times. Sent
  nothing, it stayed silent past 1.5 s, 2 of 2 times (V23, darwin).
- **`judgeReverse(tok string, seen reverseSeen) revVerdict`** is pure. It is the whole decision
  for one attempt, applied in this order:

| # | Condition | Verdict | Why |
|---|---|---|---|
| 1 | `seen.child` | **REACHED**: refuse with today's text, `reverse leg: the host reached the listener the child bound on 127.0.0.1:<port>` | Only the child knows the nonce, so a real reach always wins, whatever else happened |
| 2 | token neither `WORLD-PROBE-BOUND` nor `WORLD-PROBE-REFUSED …` | **NO-TOKEN**: refuse with today's text, `reverse leg: no token from the child's listener …`. **Not retried** | MUT-PROBE-RC-ONLY: a missing client must keep naming itself |
| 3 | `seen.silent` [r3] | **SILENT**: refuse, **sticky**: `reverse leg: a listener on 127.0.0.1:<port> accepted the host's dial and answered nothing to World's hello within 500 ms; World cannot tell it from the child's listener, so it refuses`. **Not retried** | A silent peer could be the child's listener answering slowly. Refusing keeps today's "one connected dial refuses" guarantee for every listener that does not prove it is foreign |
| 4 | `seen.foreign`, **or** token `WORLD-PROBE-REFUSED EADDRINUSE` | **INCONCLUSIVE**: retry on a fresh port | Something World did not plant holds the port and answered the hello as nothing the child would say. Seen from the child's side, that is EADDRINUSE |
| 5 | otherwise (BOUND or REFUSED, zero dials connected) | **PASS** | |

- **`(pp *execProber) reverseLeg(ctx, proxyPort) (first execProbeRun, why string)`.** Attempt 1
  is Run C exactly as today (`proxy`, then `listen`), so the happy path spawns no extra process
  and the proxy leg is still judged on `first.token(0)`. Attempts 2 and 3 run `listen` alone,
  without the trampoline and under the same settings. `execProbeReverseAttempts = 3`. REACHED,
  NO-TOKEN and SILENT end the leg at once, on whichever attempt they occur. That makes the silent
  refusal sticky: no later attempt can run to overturn it [r3]. When every attempt is
  inconclusive, the leg refuses with:
  `reverse leg: inconclusive after 3 attempts: a listener the child did not bind held each port
  (127.0.0.1:<p1> …, …)`, which lists each port and whether the attempt saw a foreign reach or
  EADDRINUSE. If the picker fails on attempt 1, Run C still runs with `proxy` alone, so the
  proxy leg is judged, and the reverse leg refuses with today's `no free port` text.
- **Seam.** An unexported field `ExecProbeConfig.reversePort func() (int, error)` (nil means
  `freeLoopbackPort`) is copied into `execProber.pickPort`. `ExecProbeConfig.Client` already
  sets the precedent of a test-steerable probe input (`host/broker/exec_probe.go:104-106`:
  "Client is the probe client run inside the sandbox; "" = Sandbox.Node … Tests point it at a
  missing path.", V24). Package tests are the only code that can set an unexported field.

### 4.3 Why a real breach is still refused, on both platforms

**Invariant [r3]: an attempt passes only when *zero* dials connected and the child gave a sandbox
token.** A connected dial has three possible outcomes:
- **child** (the nonce): REACHED, refuse.
- **silent** (zero bytes after the hello): SILENT, refuse, sticky for the leg.
- **foreign** (it answered the hello with non-nonce bytes or a clean close): INCONCLUSIVE.
  World retries on a fresh port. The leg can then pass, but only on a *later* attempt in which
  no dial connected at all; otherwise it refuses.

No connected dial ever makes *its own* attempt pass. A foreign listener can delay the leg, or
refuse it through exhaustion, but it cannot make the attempt it answered pass. The listener of
World's probe client can only ever produce child or silent, never foreign (§3, sensitivity), so a
listener the child binds always refuses. Because only a nonce match refuses as REACHED, a
hostile local process that learns the nonce from `ps` and echoes it can only cause a *refusal*.
That is fail-closed, so the nonce does not need to be secret.

| Platform and state | Child token | Host sees | Verdict |
|---|---|---|---|
| Linux, correct (own netns) | BOUND | nothing | PASS, 1 attempt |
| Linux, correct, port collided with srt's mux or another chatty holder (the CI case) | BOUND | foreign (HTTP 400, V23) | INCONCLUSIVE, then a fresh port, then PASS |
| any platform, port collided with a holder that stays silent after the hello [r3] | any | silent | **REFUSED, silent** (sticky; distinct text). This is today's outcome for this case, now named |
| Linux, shared netns (breach) | BOUND (host loopback) | nonce | **REFUSED, reached.** The raw and proxy legs also refuse on their own |
| Linux, breach, port collided | REFUSED EADDRINUSE | foreign | INCONCLUSIVE, then a fresh port, then nonce, then **REFUSED, reached** |
| macOS, correct (Seatbelt) | REFUSED EPERM | nothing | PASS, 1 attempt |
| macOS, correct, port collided | REFUSED EPERM (or EADDRINUSE; the order is unmeasured) | foreign | INCONCLUSIVE, then retry, then PASS. Either errno order is handled |
| macOS, bind allowed (breach) | BOUND | nonce | **REFUSED, reached** |
| pass-through shim (the positive control) | BOUND | nonce | **REFUSED, reached** (TestExecProbePassThroughShimFailsEveryConfinementArm) |
| every attempt collided | any | foreign ×3 | **REFUSED, inconclusive** (distinct text) |

**Budget.** Retries happen only on a collision. Today the whole probe takes 2.47 s on Linux CI
(V17), inside the daemon's 9 s startup context (V16). Each extra attempt costs one srt start plus
the 1 s listen. A run whose budget is spent gets no token (exec_probe.go:278-283), which is
NO-TOKEN and so a refusal. Exhausting the budget therefore fails closed.

## 5. Residuals (named and priced)

- **R1, a slow intermittent breach. CLOSED in r3.** r1 classed a reachable child listener that
  did not produce the nonce within 500 ms as *foreign*, so a slow breach followed by a clean
  attempt passed. r1's deferred tightening ("refuse any attempt whose foreign reaches timed out
  mid-read") would have *restored the CI flake*: srt's mux is silent until spoken to (V23, and
  V11's silent stand-in shows the same), so refusing on silence would refuse every mux
  collision. r3 removes the dilemma with the hello. The mux answers the hello with HTTP 400 and
  is retried. A slow or silent peer, which includes a slow child, refuses as SILENT on whichever
  attempt it occurs, and no later attempt can overturn that. The price is availability: a
  collision with a foreign holder that stays silent after the hello still refuses startup, as
  every collision does today, but now with a text that names it. srt's mux, the only in-window
  binder identified (V7), is not such a holder (V23).
- **R2, a cross-profile reach on a live listener** (§2). Not fixed here; it is a successor.
- **R3, EADDRINUSE becomes retry instead of pass.** Today `REFUSED EADDRINUSE` with no reach
  passes. The new design retries it. This is strictly stricter. A host where something always
  holds the chosen port refuses with the inconclusive text instead of passing vacuously.
- **R4, independence of attempts [r2].** 3 attempts suffice only if per-attempt collisions are
  near-independent. The per-srt-start mux `bind(0)` argument (V7) supports independence by
  reading, but no Linux measurement of the in-window collision rate exists. If exhaustion
  refusals appear in CI, this assumption is the first thing to re-measure.
- **Scope of V23 [r3].** The mux's answer to the hello was measured on **darwin only**. The
  claim that Linux behaves the same rests on source reading. `mux-proxy.js` peeks the first
  byte, and any byte that is not 0x04/0x05 is piped to the HTTP backend, whose `clientError`
  handler writes the 400. Only the backend's *transport* branches on Windows; it is a Unix
  socket elsewhere. If Linux CI ever shows `reverse leg: a listener … answered nothing to World's
  hello`, this is the assumption to re-measure.

## 6. Mutation table

Each mutant is applied alone. The killer is the named test, read as counted `--- FAIL:` lines.
A holder is **chatty** if it reads bytes and answers (`HTTP/1.1 400 Bad Request` then close,
srt's mux shape, V23). It is **silent** if it accepts and never writes.

| ID | Mutant | Killer, and the assertion that fires |
|---|---|---|
| M1 MUT-REV-FOREIGN-IS-REACH | a connected dial counts as REACHED whatever it read (today's behaviour) | T1 row (BOUND, foreign) wants INCONCLUSIVE. T2: Why is non-empty `reverse leg: the host reached …`, the CI signature. T6 `chatty` row: arm 6's Why contains `reverse leg` |
| M2 MUT-REV-NO-RETRY | the first INCONCLUSIVE refuses | T2: Why is `reverse leg: inconclusive after 1 …` instead of a pass, and picker calls are 1, not 2 |
| M3 MUT-REV-ANY-BYTES | any non-empty read is the nonce | T2 rows `forged-nonce`, `garbage` and `http-400` refuse as reached |
| M4 MUT-REV-PREFIX-NONCE | only the `WORLD-PROBE-NONCE ` prefix is compared | T2 row `forged-nonce` refuses as reached |
| M5 MUT-REV-EXHAUST-PASSES | exhaustion returns `""` | T3: wants a non-empty Why |
| M6 MUT-REV-UNBOUNDED | no attempt bound (loop until pass) | T3: picker calls are not exactly 3 (it hangs to the test's 60 s context and reds there) |
| M7 MUT-REV-EXHAUST-AS-REACH | exhaustion reuses the reached text | T3: Why must NOT contain `the host reached the listener the child bound` |
| M8 MUT-REV-EADDRINUSE-PASSES | `REFUSED EADDRINUSE` with no reach passes (today's behaviour) | T1 row (EADDRINUSE, none). T3 row `child-eaddrinuse` passes instead of refusing |
| M9 MUT-REV-NEVER-REAL | the nonce is never matched (comparison inverted, or the script stops writing it) | T4: Why is not the reached text on port 2 (it becomes the inconclusive or silent text). **P1**, the existing positive control, loses its `reverse leg: the host reached` substring |
| M10 MUT-REV-SAME-PORT | a retry reuses the attempt-1 port | T4: the chatty holder is still there, so exhaustion text replaces the reached text on port 2 |
| M11 MUT-REV-FOREIGN-DOMINATES | (child and foreign) is judged INCONCLUSIVE | T1 row (BOUND, child+foreign) wants REACHED |
| M12 MUT-REV-NOTOKEN-RETRIED | NO-TOKEN is judged INCONCLUSIVE | T1 row ("", none) wants NO-TOKEN. T5: picker calls are 3, not 1, and the Why lacks `no token` |
| M13 MUT-REV-ALWAYS-RETRY | every attempt runs even after a PASS | T5 clean row: picker calls are 3, not 1 |
| M14 MUT-REV-SEAM-ONLY | `probe()` keeps the old inline leg and only the seam path is fixed | T6 `chatty` row: through `ProbeExecProfiles`, arm 6's Why contains `reverse leg` |
| M15 MUT-REV-HELLO-DROPPED [r3] | `pollLoopback` never writes the hello | T2 row `http-400`: the holder stays silent, so the leg refuses SILENT instead of passing. T6 `chatty` row: arm 6's Why contains `reverse leg` |
| M16 MUT-REV-SILENT-RETRIED [r3] | SILENT is judged INCONCLUSIVE (retry) | T1 row (BOUND, silent) wants SILENT. T7 row `silent-attempt-1`: it passes on attempt 2 instead of refusing, with picker calls 2, not 1 |
| M17 MUT-REV-SILENT-NOT-STICKY [r3] | a silent attempt is forgotten when a later attempt is clean (for example SILENT records only "retry", and a later PASS returns `""`) | T7 row `silent-after-chatty`: wants the silent text naming port 2 and picker calls 2. The mutant passes on attempt 3 |
| M18 MUT-REV-EOF-IS-SILENT [r3] | a clean zero-byte EOF is classed silent | T2 row `clean-close` refuses SILENT instead of passing |
| M19 MUT-REV-RESET-IS-FOREIGN [r3] | a zero-byte reset is classed foreign | T7 row `reset-holder` (SO_LINGER 0 close) passes on attempt 2 instead of refusing SILENT |
| M20 MUT-REV-SILENT-AS-REACH-TEXT [r3] | SILENT reuses the reached text | T7: Why must NOT contain `the host reached the listener the child bound`, and must contain `answered nothing to World's hello`. T6 `silent` row: the same |
| M21 MUT-REV-SILENT-DOMINATES-CHILD [r3] | (child and silent) is judged SILENT | T1 row (BOUND, child+silent) wants REACHED |

## 7. Acceptance criteria

Every `-run` verdict is read from counted `--- PASS:` / `--- FAIL:` / `--- SKIP:` lines, never
from rc alone, because a `-run` that matches nothing is rc 0.

**Test table [r3]:**

| Test | Needs node? | Kills |
|---|---|---|
| T1 `TestExecProbeReverseLegJudge` | no (pure) | M1, M8, M11, M12, M16, M21 |
| T2 `TestExecProbeReverseLegChattyForeignHolderIsRetried` | no (`/bin/sh` fixture, V24) | M1, M2, M3, M4, M15, M18 |
| T3 `TestExecProbeReverseLegExhaustionRefusesDistinctly` | no (`/bin/sh` fixture) | M5, M6, M7, M8 |
| T4 `TestExecProbeReverseLegRealReachSurvivesAForeignAttempt` | yes | M9, M10 |
| T5 `TestExecProbeReverseLegNoTokenIsNotRetried` | no (`/bin/sh` fixture) | M12, M13 |
| T6 `TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse` | yes (node-run srt stand-in) | M1, M14, M15, M20 |
| T7 `TestExecProbeReverseLegSilentHolderRefusesSticky` | no (`/bin/sh` fixture) | M16, M17, M19, M20 |
| P1 `TestExecProbePassThroughShimFailsEveryConfinementArm` (existing, unchanged) | yes | M9 |

- **T1 `TestExecProbeReverseLegJudge`** (pure table, no I/O). It holds §4.2's order in 14 rows:
  - (BOUND, none) → PASS; (REFUSED EPERM, none) → PASS.
  - (BOUND, foreign) → INCONCLUSIVE; (REFUSED EPERM, foreign) → INCONCLUSIVE;
    (REFUSED EADDRINUSE, none) → INCONCLUSIVE.
  - (BOUND, child) → REACHED; (BOUND, child+foreign) → REACHED; (BOUND, child+silent) →
    REACHED; ("", child) → REACHED.
  - (BOUND, silent) → SILENT; (BOUND, silent+foreign) → SILENT; (REFUSED EPERM, silent) →
    SILENT.
  - ("", none) → NO-TOKEN; ("", foreign) → NO-TOKEN.
- **T2 `TestExecProbeReverseLegChattyForeignHolderIsRetried`** [r3] (an `arm8Fixture`-style
  `/bin/sh` client and srt; it runs with node absent, V24). Rows are {`linux-bound`,
  `darwin-eperm`} × {`http-400` (reads the hello, answers `HTTP/1.1 400 Bad Request`, closes:
  the mux shape), `forged-nonce` (a wrong 32-hex nonce in the right format), `garbage`,
  `clean-close` (accepts, closes with zero bytes)}. The picker returns, on attempt 1, a port held
  by a test-owned Go listener of that behaviour, and on attempt 2 a free port. Each row wants
  `why == ""` and exactly 2 picker calls.
- **T3 `TestExecProbeReverseLegExhaustionRefusesDistinctly`** (`/bin/sh` fixture). Rows:
  `chatty-every-time` (BOUND, every port held by an `http-400` holder) and `child-eaddrinuse`
  (free ports, the client prints `WORLD-PROBE-REFUSED EADDRINUSE`). Each row wants the Why to
  start with `reverse leg: inconclusive after 3 attempts`, to NOT contain `the host reached the
  listener the child bound`, and exactly 3 picker calls.
- **T4 `TestExecProbeReverseLegRealReachSurvivesAForeignAttempt`** (needs node; `probeNode`
  skips by name; the pass-through shim with World's real node client). The picker gives an
  `http-400`-held port on attempt 1 and a free port on attempt 2. Attempt 1: the unsandboxed
  child gets EADDRINUSE and the holder answers the hello with 400, so the attempt is
  INCONCLUSIVE. Attempt 2: the child binds and writes the nonce. The test wants Why ==
  `reverse leg: the host reached the listener the child bound on 127.0.0.1:<attempt-2 port>`.
  This is the "a real reach still refuses" proof: a collision on attempt 1 does not mask a
  breach on attempt 2.
- **T5 `TestExecProbeReverseLegNoTokenIsNotRetried`** (`/bin/sh` fixture). Row `silent-client`
  (no listen token) wants a Why containing `reverse leg: no token from the child's listener` and
  1 picker call. Row `clean` (BOUND, free port) wants `""` and 1 picker call.
- **T6 `TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse`** (needs node; the production
  entry point `ProbeExecProfiles`). An srt stand-in binds a holder on the `listen` port only when
  the argv also carries `proxy`, that is, on attempt 1 = Run C, as srt's own mux would. It prints
  the Linux-correct tokens without running the child, and passes every other run through. Rows:
  - `chatty` [r3]: the holder answers any bytes with HTTP 400, then closes. This is the
    deterministic **flake reproduction**. It wants arm 6's Why to contain no `reverse leg` text.
    **Base: RED, 3 of 3** at `2769cff9` (V25).
  - `silent` (V11's stand-in, kept): it wants arm 6's Why to contain `reverse leg: a listener on
    127.0.0.1:` and `answered nothing to World's hello`. This is the sticky refusal through the
    production entry point.
  - `control` (no holder): it wants no `reverse leg` text.

  Neither stimulus is timing-dependent: the holder keeps the port for the whole run, and the
  poll dials every 50 ms.
- **T7 `TestExecProbeReverseLegSilentHolderRefusesSticky`** [r3] (`/bin/sh` fixture). Rows:
  - `silent-attempt-1`: a silent holder, then free ports. It wants the Why to start with
    `reverse leg: a listener on 127.0.0.1:<p1> accepted the host's dial and answered nothing`,
    no reached text, and 1 picker call.
  - `silent-after-chatty`: `http-400` on attempt 1, silent on attempt 2, free on attempt 3. It
    wants the silent text naming `<p2>` and 2 picker calls.
  - `reset-holder`: accepts, then closes with SO_LINGER 0, so it resets with zero bytes. It
    wants the silent text and 1 picker call.
- **P1, the positive control, unchanged and green:** `TestExecProbePassThroughShimFailsEveryConfinementArm`
  keeps asserting `reverse leg: the host reached` (exec_probe_test.go:197). It passes now only
  through the nonce path, because the pass-through child binds host loopback and answers with
  its nonce. `TestExecProbeMissingClientFailsTheRefusalArms` stays green unchanged (arm 6 is
  still refused with no token).
- **AC-CI [r3].** **All seven** new tests, T1–T7, are added to step 8's `tests=` list in
  `ci.yml`. That step fails unless each listed test prints `--- PASS:` (ci.yml:254-257, V5).
  "Never skips" is therefore a checked gate, not an assumption. Any SKIP on the Linux runner,
  whether from absent node or a hidden environment dependence, fails CI, and M1–M21 stay
  enforced there (S6). V24 measured the `/bin/sh` fixture machinery that T2, T3, T5 and T7
  reuse. With node absent from `PATH` and `WORLD_EXEC_SRT_NODE_MODULES=/nonexistent`, it printed
  5 `--- PASS` lines and 0 `--- SKIP` lines. T4 and T6 are node-needing and listed anyway.
- **AC-REAL.** `TestExecSrtStartupProbePasses`, plus step 8 in full on Linux CI with real srt,
  stays green. That covers the Linux correct path: BOUND, no reach, PASS on 1 attempt.
- **AC-STABLE.** `go test ./host/broker -run '^TestExecProbeReverseLeg|^TestExecProbeForeignListener' -count=20 -v`
  must show every new test PASS 20/20 with zero FAIL. This is a stability read, not a rate
  claim: no stimulus here is timing-dependent.
- **AC-GATE.** `go build ./... && go test ./...` is green, and so is `./scripts/verify_ail.sh`.
  The latter is a no-op for this item, since no `.ail` changes, but CI runs it.
- **AC-MUT.** Each of M1–M21 is applied alone and reds exactly its named killer. The tree is
  restored byte-identical (sha256 of both Go files before and after).

**What this cannot see.** None of these tests runs the *CI* binder, srt's real mux, inside a
probe. They reproduce the class with a stand-in shaped like the mux as V23 measured it (chatty
on input, silent otherwise). The mux measurement itself is darwin-only (§5). The CI instance's
actual binder stays unproven (§1). The fix does not depend on it for safety: an unidentified
silent binder refuses, as today.

## 8. Coding standards

- **S2, effects at the boundary.** All of this lives in the Go host's probe, which is already
  the effect boundary: sockets, subprocesses, srt. `judgeReverse` is split out as a pure
  function, so the decision is table-tested without I/O (T1). `world/` is untouched.
- **S3, "why is this not a package?"** This is a repair inside an existing host component, the
  fail-closed startup probe of the brokered `Workspace.Exec` effect, which row 140 placed in the
  host. It adds no kernel surface and no new capability. Moving the probe out is not in scope.
- **S6, honest gates.** Every load-bearing claim is discharged by a named mutant and its single
  firing assertion (§6). The flake reproduction (T6 `chatty`) is red at base by measurement
  (V25), with a paired control. **All seven** new tests are in the CI list that fails on a
  missing `--- PASS` [r3], so no SKIP, node-related or not, passes silently. A connected dial
  never makes its own attempt pass. Silent and nonce dials refuse, and only a dial that answered
  the hello as foreign is retried (§4.3), so the leg has no vacuous-green path.
- **S7, usage docs.** The probe's operator-facing text gains two new refusals,
  `reverse leg: inconclusive …` and [r3] `reverse leg: a listener on … answered nothing to
  World's hello …`. QUICKSTART §10's arm 6 line (docs/QUICKSTART.md:1146-1148,
  V21) gains one sentence. It says a reach counts only when the child's listener answers with
  World's nonce, that a port another host process holds is retried on a fresh port up to 3
  times, and what the two refusals mean. Inconclusive means something on the host keeps taking
  the chosen ports, so retry, or find the process. Silent [r3] means a process that never answers
  held the port, and it cannot be told apart from the child, so find the process. The §10 refusal table's 51-row census is unchanged,
  because the arm's name is unchanged.
- **S1, S4, S5:** not applicable, since there is no `.ail` (V15).

## 9. Conflict surface

`exec_probe.go` has not changed since `885af2d1` (V1). The edit is confined to:
- the `listen` op of `execProbeScript`;
- `probe()` lines 452-488, the Run C and reverse-leg block;
- `pollLoopback`;
- `ExecProbeConfig` and `execProber`.

Arm 8 (`probeProtected`) and arms 1–5 and 7 are untouched.

**Caller census [r2].** The two helpers this item changes or re-homes are `pollLoopback`, whose
signature and semantics change, and `freeLoopbackPort`, which becomes the default port picker.
`grep -rn 'pollLoopback\|freeLoopbackPort' --include='*.go' .` returns exactly 6 lines, all in
`host/broker/exec_probe.go`:
- `:455`, the call `revPort, err := freeLoopbackPort()`;
- `:460`, the call `poll := pollLoopback(revPort)`;
- `:551/:552`, `freeLoopbackPort`'s doc comment and definition;
- `:561/:563`, `pollLoopback`'s doc comment and definition.

Each helper has exactly one call site, and both are inside the reverse leg. The positive control
in the same run, `grep -rn 'func runBounded\|\.runBounded(' --include='*.go' host | wc -l`,
returns `9`, and `host/`, `cmd/` and `world/` all exist, so the grep covered them (V22). The
reverse leg is therefore the only caller of both helpers. "Untouched" for arms 1–5, 7 and 8 is
now verified by this census rather than asserted. My re-run matches the controller's census
line for line.

The row 140 design doc's §4.6 wording
("a listener the child binds must be unreachable from the host") stays true and needs no
amendment. [r3] The reviewer's fix branched on whether srt's mux answers garbage input. It does,
with HTTP 400 (V23), so the strong branch holds. P20 is not narrowed, and §4.6 is not amended.

## 10. Sizing

| File | Change | LOC (est.) |
|---|---|---|
| `host/broker/exec_probe.go` | the script's listen op (+nonce, +`c.resume()`); `pollLoopback` writes the hello, reads, and classifies child / foreign / silent [r3]; `judgeReverse` (pure, 5 verdicts); `reverseLeg` (attempt loop, sticky silent, two distinct messages); seam field; `probe()` calls `reverseLeg` | +95 / −20 |
| `host/broker/exec_probe_test.go` | T1–T7; a holder helper (http-400 / forged / garbage / clean-close / silent / reset); a counting picker; T6's node stand-in in two shapes | +260 |
| `docs/QUICKSTART.md` | §10 arm 6 sentences | +3 |
| `.github/workflows/ci.yml` | step 8 `tests=` gains all seven, T1–T7 [r3] | +1 / −1 |

About 0.5–0.75 executor day, including the 21-mutant pass.

## 11. What this item is NOT doing

- **R2, cross-profile reaches on a live counted listener.** This is a successor candidate,
  `w-exec-probe-cross-profile-poll`. Options there include a process-wide registry of World's
  own probe-listener ports that the poll consults before each dial, or a hello on the poll's
  dials that the live listener does not count. It needs ≥2 concurrently probed profiles. It was
  not measured.
- Tightening the reverse leg's accepted refusal codes to `EPERM`/`EACCES` (only EADDRINUSE
  changes here).
- Changing how the raw and proxy legs count accepts.
- Any srt change. srt's mux binding port 0 is correct srt behaviour. The defect is World's
  inference.

## 12. Verification Log

All commands were run from the worktree at `2769cff9` unless a row says otherwise. The r3 rows
V23–V26 ran at `b225976a`, whose Go tree is byte-identical to `2769cff9`; that commit changed
only this doc. Scratch test files (V11–V13, V25, V26) were created under
`host/broker/zz_iter314_*_test.go`, run, and **deleted**. `git status --porcelain` showed no Go
change afterwards.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | `exec_probe.go` is unchanged since `885af2d1` | `git log --format='%h %ad %s' --date=short -- host/broker/exec_probe.go \| head -5`; `git diff 885af2d1 HEAD -- host/broker/exec_probe.go \| wc -l` | top: `885af2d1 2026-10-20 feat(row 156 M4c)…`; diff lines `0` |
| V2 | Today's reverse leg, poll and client listener | Read `host/broker/exec_probe.go` | `:80` `net.createServer(c=>c.destroy())`; `:455-464` freeLoopbackPort → pollLoopback → Run C with `listen <port>:<ms>`; `:481-485` no-token else `reached` → `reverse leg: the host reached the listener the child bound on 127.0.0.1:%d`; `:551-559` Listen `127.0.0.1:0`, read port, `ln.Close()`; `:573-577` any `DialTimeout` success → `reached.Store(true)`; `:466-488` legs judged separately and joined with `"; "` |
| V3 | The CI failure text and its context | `gh api repos/sunholo-data/ailang-world/actions/jobs/122853572046/logs --allow-escape-sequences \| grep -n 'TestExecSrtMCPGrantAndBudget\|TIMING the startup probe'` | `:614 === RUN TestExecSrtMCPGrantAndBudget`; `:615 exec_e2e_srt_test.go:211: New: daemon startup failed at config: … (1 of 7 arms): arm6-network: reverse leg: the host reached the listener the child bound on 127.0.0.1:44109`; `:616 --- FAIL: TestExecSrtMCPGrantAndBudget (2.63s)` |
| V4 | Attempt 1 failed, attempt 2 passed | `gh api …/runs/40699055042/attempts/1/jobs --jq …`; same for `attempts/2` | attempt 1: `122853572046 go host build + test gate failure 1`; attempt 2: `122935140505 go host build + test gate success 2` |
| V5 | Step 8 runs packages serially and guards on PASS lines | `grep -n 'Row 140 Workspace.Exec real-srt' -A30 .github/workflows/ci.yml` | `:250 go test ./host/broker/ ./host/daemon/ -count=1 -p 1 -timeout 600s -v`; `:254-256` a `grep -q -- "^--- PASS: $t "` loop, `::error::$t did not PASS on linux (SKIP?)` |
| V6 | No `t.Parallel` in the two packages (negative, with a control) | `cat host/daemon/*_test.go host/broker/*_test.go \| grep -c 't.Parallel()'`; `grep -rl 't.Parallel()' --include='*_test.go' . \| wc -l` | `0`; control `3` (e.g. `host/verifygate/mission_config_gate_test.go`) |
| V7 | Pinned srt 0.0.78 binds its mux on host `127.0.0.1:0` per start; the Linux host bridges are Unix sockets | On `~/.ailang/state/mission-world-iter241/srt/node_modules/@anthropic-ai/sandbox-runtime`: `grep '"version"' package.json`; `shasum -a 256 dist/cli.js`; `sed -n 699-720p dist/sandbox/sandbox-manager.js`; `sed -n 433p …`; `sed -n 10,11p dist/sandbox/listen-in-range.js`; `grep -n LISTEN dist/sandbox/linux-sandbox-utils.js`; `grep -n 'ExecSandboxRelease *=\|ExecSandboxCLISHA256 *=' host/broker/*.go` | `0.0.78`; `3c3092bd…d96e` = World's `ExecSandboxCLISHA256` (exec_sandbox.go:29); "Other platforms bake the actual ephemeral port … so they keep using port 0", `portRange = … windows ? … : undefined`; `listenInRange(mux.server, p => mux.server.listen(p, '127.0.0.1'), portRange, …)`; `const [lo, hi] = range ?? [0, 0]; let port = lo;`; `:917/:954 UNIX-LISTEN` (host), `:1066-1067 TCP-LISTEN:3128/1080` (inside the namespace) |
| V8 | Go's dialer retries a TCP self-connect | `grep -n selfConnect "$(go env GOROOT)/src/net/tcpsock_posix.go"` (go1.26.6) | `:111 for i := 0; i < 2 && … (selfConnect(fd, err) \|\| spuriousENOTAVAIL(err)); i++ {`; `:124 func selfConnect` |
| V9 | M0: Linux binds in its own netns and is unreachable; macOS refuses | `sed -n 80,81p …/out/linux-0.0.78-run3-both-modes.txt`; `sed -n 38,39p …/out/macos-0.0.78.txt`; `sed -n 604-622p …/m0_matrix.py` | Linux `observed 'WORLD-PROBE-BOUND'` / `'WORLD-PROBE-BOUND 38327; host reached=False'`; macOS `'WORLD-PROBE-REFUSED EPERM'` / `'…; host reached=None'`; M0 used `listen(0)` plus one host connect after BOUND |
| V10 | The existing probe tests are green locally, and the real-srt tests skip by name | `go test ./host/broker -run 'TestExecProbe' -count=1 -v` | `--- PASS` for ArmNames, PassThroughShim… (1.98s), NoopShim…, ZeroShim…, FailingToolchain…, MissingClient…, Arm8Wired, Arm8FailsUnsandboxed (2 linux rows SKIP by design), Arm8Absent, ProtectedSample; `SKIP TestExecProbeArm8: real srt absent: WORLD_EXEC_SRT_NODE_MODULES is unset …` |
| V11 | **Deterministic reproduction at base**: a non-child listener on the reverse port produces the CI signature | A scratch test: `newProbeFixture` + `fakeSrtTree` with a node stand-in that, when the argv has `listen`, binds `net.createServer(c=>{})` on that port and prints `WORLD-PROBE-REFUSED ECONNREFUSED` / `WORLD-PROBE-BOUND` without running the child (else pass-through); a control row without the bind. `go test ./host/broker -run TestScratchIter314 -count=3 -v` | arm row `REVERSE-LEG-REACHED=true` 3/3; control `false` 3/3. Arm 6 Why: `… reverse leg: the host reached the listener the child bound on 127.0.0.1:60489`. **[r3] reclassified:** this *silent* stand-in becomes the sticky-refusal stimulus (T6 `silent`, T7); the flake reproduction is V25's chatty stand-in |
| V12 | (Instrument context only, darwin, not Linux evidence) an immediate `bind(0)` rarely returns the just-freed port | A scratch loop of 20,000 × (`freeLoopbackPort`; `net.Listen("127.0.0.1:0")`; compare); `sysctl net.inet.ip.portrange.first net.inet.ip.portrange.last` | `NEXT-BIND-SAME-PORT 0 of 20000`; range `49152`–`65535` |
| V13 | A nonce line from a node listener is read promptly | A scratch test: a node `net.createServer(c=>c.end('WORLD-PROBE-NONCE abc123\n'))`; 500 × (DialTimeout 200 ms, 500 ms read deadline, `ReadString('\n')`) | `NONCE-READ ok=500/500 max=648µs` |
| V14 | No daemon test asserts reverse-leg text | `grep -rn "reached the listener\|reverse leg" host cmd docs \| grep -v '^host/broker/exec_probe.go:'` | only `host/broker/exec_probe_test.go:189` (comment) and `:197` (assertion), which double as the known-positive control |
| V15 | No `.ail` involvement | `grep -rl 'arm6\|reverse leg\|WORLD-PROBE\|ExecProbe' --include='*.ail' . \| wc -l`; control `grep -rl 'Exec' --include='*.ail' . \| wc -l`; control `grep -rl 'WORLD-PROBE' --include='*.go' . \| wc -l` | `0`; controls `7` and `3` |
| V16 | 9 s daemon startup context; probe runs clamp to it | `grep -n '9\*time.Second' host/daemon/daemon.go`; read exec_probe.go:277-283 | `:1248 d, err := startBounded(ctx, cfg, 9*time.Second, New)`; `budget = min(budget, time.Until(dl))`; `budget <= 0` → `the startup budget is spent` |
| V17 | Linux CI probe duration | Attempt-1 log, `grep -n 'TIMING the startup probe'` | `:540 exec_probe_srt_test.go:34: TIMING the startup probe took 2.466536014s` |
| V18 | Production consequence of a failed profile | `sed -n 236,258p host/broker/exec_startup.go` | a failed `--exec-profile` returns the error (startup refused); a failed directory profile is appended to `st.Excluded` and deleted from `st.Profiles` |
| V19 | Profiles are probed concurrently | Read exec_probe.go:174-183 | one goroutine per profile, `wg.Wait()` |
| V20 | `exec_startup_test.go:499` is a log, not an assertion | `sed -n 499p host/daemon/exec_startup_test.go` | `t.Logf("%s: %v", id, err)` |
| V21 | QUICKSTART's arm 6 prose (the S7 edit site) | `sed -n 1140,1160p docs/QUICKSTART.md` | `:1146-1148` "`arm6-network`: World opens its own live `127.0.0.1` listener … a port the child binds must be unreachable from the host." |
| V22 | **[r2]** Caller census: the reverse leg is the only caller of `pollLoopback` and `freeLoopbackPort`. Run at worktree HEAD `72eba94b` | `grep -rn 'pollLoopback\|freeLoopbackPort' --include='*.go' .`; control `grep -rn 'func runBounded\|\.runBounded(' --include='*.go' host \| wc -l`; `ls -d host cmd world` | Exactly 6 lines, all in `host/broker/exec_probe.go`: `:455 revPort, err := freeLoopbackPort()`, `:460 poll := pollLoopback(revPort)`, `:551 // freeLoopbackPort is a port nothing listens on …`, `:552 func freeLoopbackPort() (int, error) {`, `:561 // pollLoopback dials 127.0.0.1:port every 50 ms …`, `:563 func pollLoopback(port int) func() bool {`. Control `9`. `cmd host world` all exist. This is identical to the controller's census (VERIFIED BY CONTROLLER) |
| V23 | **[r3]** srt 0.0.78's mux answers a non-SOCKS hello with HTTP 400 and closes, but stays silent if sent nothing. **Scope: darwin only**; Linux rests on source reading, since the mux path is the same on every non-Windows platform | VERIFIED BY CONTROLLER (darwin, node v26.12.1), then re-run first-party: `shasum -a 256 $S/dist/cli.js \| cut -c1-12` and `perl -e 'alarm 60; exec @ARGV' node /tmp/w314/muxprobe/probe.mjs $S`, with `S=~/.ailang/state/mission-world-iter241/srt/node_modules/@anthropic-ai/sandbox-runtime`. The script (sha256 prefix `28c3311502f8`) imports `createMuxProxyServer` + `createHttpProxyServer` from `$S`, listens the mux on `127.0.0.1:0`, and dials it 3× sending `WORLD-PROBE-HELLO\n` and 2× sending nothing (1500 ms wait). Source: `sed -n 76,100p $S/dist/sandbox/mux-proxy.js`; `grep -n '0x04\|0x05\|DEFAULT_FIRST_BYTE' …/mux-proxy.js`; `grep -n clientError -A10 $S/dist/sandbox/http-proxy.js` | `3c3092bd26b3` (= pin). Controller: `HELLO` 3/3 `close after 0-3ms bytes=28 first="HTTP/1.1 400 Bad Request"`; `SILENT` 2/2 `deadline after 1500ms bytes=0`. First-party re-run: `HELLO  : close after 3ms / 1ms / 0ms bytes=28 first="HTTP/1.1 400 Bad Request"`; `SILENT : deadline after 1500ms bytes=0` ×2; rc 0. Source: `SOCKS_FIRST_BYTES = new Set([0x04, 0x05])`, `DEFAULT_FIRST_BYTE_TIMEOUT_MS = 10000`, `client.once('readable', … peek … SOCKS_FIRST_BYTES.has(peek[0]) ? handleSocksConnection : dispatchHttp)`; `server.on('clientError', …)` writes the 400 when `socket.writable` ("Destroy only after the flush"). `W` = 0x57 is not 0x04/0x05 |
| V24 | **[r3, glm]** The `arm8Fixture`-style `/bin/sh` client + fake-srt probe path runs with **node absent** and no real srt: 0 SKIP. The `ExecProbeConfig.Client` precedent, file:line | `env -i HOME=$HOME PATH=/usr/bin:/bin sh -c 'command -v node; command -v go'` (control: node absent); `ls /usr/bin/node /bin/node`; then `env -i HOME=$HOME PATH=/usr/bin:/bin WORLD_EXEC_SRT_NODE_MODULES=/nonexistent /opt/homebrew/bin/go test ./host/broker -run '^TestExecProbeArm8FailsUnsandboxed$/^(pass-through\|missing-client\|forged-token\|darwin-refuses-all)$' -count=1 -v`, then count `--- PASS` / `--- SKIP`. `grep -n 'Client  *string' host/broker/exec_probe.go`; `grep -n 'func arm8Fixture\|t.Setenv(execNodeEnv' host/broker/exec_probe_test.go` | node: not found (rc 1), and `ls` says *No such file or directory* for both paths. Test: rc 0, `PASS=5 SKIP=0` (the parent and 4 subtests, each driving `pp.run` through the sh srt shim). `Client string` at `host/broker/exec_probe.go:106`, with its doc at `:104-105` ("… Tests point it at a missing path."). `arm8Fixture` at `exec_probe_test.go:256` sets `execNodeEnv` to `/bin/sh` at `:258`. Deviation from the reviewer's verbatim command: `go` is invoked by absolute path because `PATH=/usr/bin:/bin` has no `go` on this machine. An earlier run of the whole `TestExecProbeArm8FailsUnsandboxed` + `TestExecProbeArm8Absent` under the same `env -i` showed `PASS=12 SKIP=2`. The 2 SKIPs are the **host-GOOS-gated** `linux-ebusy-passes` / `linux-exdev-passes` rows (`exec_probe_test.go:352-353`), not an environment skip, so the counted run excludes them |
| V25 | **[r3]** Flake reproduction at base with a **chatty** (mux-shaped) stand-in | A scratch test: like V11, but the stand-in's holder is `net.createServer(c=>{c.once('data',()=>c.end('HTTP/1.1 400 Bad Request\r\nConnection: close\r\n\r\n'))})`. `perl -e 'alarm 300; exec @ARGV' go test ./host/broker -run TestScratchR3 -count=3 -v` | `CHATTY-BASE REVERSE-LEG-REACHED=true reverse leg: the host reached the listener the child bound on 127.0.0.1:53132` (then `:53659`, `:54186`): 3/3 |
| V26 | **[r3]** Hello-then-nonce against a node listener of §4.1's shape | Same scratch run: node `net.createServer(c=>{c.resume();c.end('WORLD-PROBE-NONCE abc123\n')})`; 500 × (DialTimeout 200 ms, write `WORLD-PROBE-HELLO\n`, 500 ms read deadline, `ReadString('\n')`) | `HELLO-NONCE ok=500 other=0`, max `495µs` / `612µs` / `515µs` over the 3 runs |

## 13. Quorum log

**Round 1** (2026-10-28; artifact
`~/.ailang/state/mission-quorum-world/w-exec-probe-reverse-leg-port-race-2026-10-28T16-32-14Z.json`):
**BLOCKED**.

| Reviewer | Verdict |
|---|---|
| `oc-glm-5-3` | PASS |
| `gemini-4-argon` | PASS |
| `oc-kimi-k3` | REJECT |
| `gpt6-1-sol` | absent (OpenAI 429, no credits) |

- **The objection (`oc-kimi-k3`).** The reviewer did not dispute the design's direction. Its
  objection was that §9's claim that arms 1–5, 7 and 8 are untouched had no caller census for
  the two helpers whose signatures or semantics change, `pollLoopback` and `freeLoopbackPort`.
  **REAL**: the claim was asserted, not measured.
- **Fix applied in r2.** This applies the reviewer's proposed fix verbatim, with measured values:
  - a §9 "Caller census" paragraph and Verification Log row V22, re-run first-party and matching
    the controller's census (6 lines, one call site per helper, both in the reverse leg;
    control 9);
  - residual R4 in §5, on the independence of attempts and the missing Linux measurement of the
    in-window collision rate;
  - this log.

  Nothing else changed.

**Round 2** (2026-10-28; artifact
`~/.ailang/state/mission-quorum-world/w-exec-probe-reverse-leg-port-race-2026-10-28T16-37-10Z.json`):
**BLOCKED**.

| Reviewer | Verdict |
|---|---|
| `gemini-4-argon` | PASS |
| `oc-glm-5-3` | REJECT |
| `oc-kimi-k3` | REJECT |
| `gpt6-1-sol` | absent (429) |

- **`oc-kimi-k3`: the doc's central safety claim contradicts its design.** **REAL.**
  - §3's "exactly as sensitive" and §4.3's "Neither is PASS" were false as written: a
    connected, non-nonce dial was INCONCLUSIVE, and a later attempt could PASS. R1 conceded a
    slow-breach weakening.
  - R1's deferred tightening would have restored the flake, because the mux is silent until
    spoken to.
  - Proposed fix: write a hello on every dial and classify child / foreign / silent holder,
    with a silent holder refusing stickily. This was gated on measuring whether the mux answers
    garbage input.
  - **Controller measurement (V23, darwin, pinned srt):** hello → `HTTP/1.1 400`, then close,
    3/3; nothing sent → silent past 1.5 s, 2/2. Source reading agrees. The reviewer's first
    branch therefore applies: "answer/close → the discriminator works and the strong invariant
    stands untouched".
- **`oc-glm-5-3`: "T1, T2, T3, T5 never skip" was asserted, not measured.** **REAL.** Only T4
  and T6 were CI-guarded, so a hidden node or srt dependence in the sh fixtures could have
  dropped the M1–M14 killers on Linux while every gate stayed green. The `ExecProbeConfig.Client`
  precedent also had no file:line.
- **Carve-out r3: reviewers' verbatim fixes applied.**
  - **kimi:**
    - §3 (A) rewritten around the hello discriminator. "Exactly as sensitive" is kept and is
      now true, with its premise stated.
    - §4.1 adds `c.resume()`. §4.2 adds the hello, the three-way classification and the sticky
      SILENT verdict. A zero-byte reset counts as silent, under the reviewer's "zero bytes after
      the hello" branch.
    - §4.3's invariant replaces "Neither is PASS" with the exact three-outcome statement.
    - R1 is closed, the deferred-tightening note is replaced, and the V23 scope is stated
      (darwin-only).
    - §2 and V11 reclassify the silent stand-in. The flake reproduction is now V25's chatty
      stand-in.
    - Tests: T2 becomes the chatty-holder retry test, T6 gains `chatty` and `silent` rows, and
      T7 (sticky silent refusal) is new. Mutants M15–M21 are added.
    - Rows V23, V25 and V26.
  - **glm:**
    - AC-CI lists all seven new tests in step 8's `tests=` list.
    - V24 records the node-less run (0 SKIP) and the `Client` file:line. The `go` binary was
      invoked by absolute path, because `PATH=/usr/bin:/bin` has no `go` on this machine.
  - **Not applied (outside both proposed fixes):** kimi's catch also observed that a leg that
    passes after inconclusive attempts emits nothing ("no-silent-fallback"). Neither proposed
    fix asked for a change there, so none was made.

## Related

- `design_docs/planned/w-workspace-exec-toolchain-effect.md`: row 140's §4.6 startup probe, the
  M0 platform tokens (V45, V46), P20, and the mutants MUT-PROBE-NET-VACUOUS and
  MUT-PROBE-RC-ONLY. Both mutants stay killed: P1, and T5 with the existing missing-client test.
- `design_docs/verification/world-row140-m0/`: the M0 matrix and its Linux and macOS outputs
  (V9).
- `design_docs/planned/w-daemon-timeout-test-flake.md`: the shape template; an earlier flake
  fixed by moving the stimulus off a race.
