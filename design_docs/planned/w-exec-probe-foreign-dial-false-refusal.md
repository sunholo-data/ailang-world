# w-exec-probe-foreign-dial-false-refusal — judge arm 6's forward legs on the child's own line, so a local loopback dialer no longer refuses startup

- **Status**: Planned (designed iteration 359, r1). Queue row 250, routed by D-WORLD-108 = B
  (attended 2026-11-06: rows 250, then 255 → 257). The row says it "needs its own design and
  quorum"; this is that design. It builds on row 242 (`w-exec-probe-cross-profile-net6-hoist.md`,
  whose §6 R1/R2 first measured this residual) and reuses the reverse leg's vocabulary
  (`w-exec-probe-reverse-leg-pins.md`, row 230: child / foreign / silent, foreign retried, silent
  sticky). It does not touch the reverse leg's retry or reset semantics (row 233, parked,
  D-WORLD-105).
- **Item**: Clause 3 (MET; availability hygiene on the Workspace.Exec fail-closed startup probe).
  It changes what arm 6's two forward legs count, so it is a structure change on a fail-closed gate
  and needs its own quorum.
- **Designer lane**: `claude:claude-opus-5-5` via Agent opus.
- **Estimated**: ~0.5 day. Production `exec_probe.go` +197/−33 lines (about half is comments and the
  classifier). One new test file, `host/broker/exec_probe_net6_dial_test.go` (299 lines, four
  tests). Two existing tests need small edits: three lines in `exec_probe_test.go` and one regex in
  `exec_probe_net6_test.go`. Four names go into the `tests=` list at `ci.yml:247`. The byte-exact
  prototype is `design_docs/verification/world-iter359/designer-proto.patch.txt` (V13).
- **Measurement base**: `1835d8be` (= origin/dev), worktree `sprint/world-iter359-foreign-dial`.
  `exec_probe.go` sha256 prefix `19de369381579387`, which is row 242's landed prototype (V1).
- **Instruments**: `go1.26.6 darwin/arm64`. That matches `go.mod`'s `go 1.26.6`, so no toolchain
  switch was needed. Also node v26.12.1, git 2.54.0 and 16 cores (V1). The rig was shared, with load
  averages of 2.9–9.8 at the timed runs. The Sunholo start page (`serve.py`, pid 45693) was running
  throughout (V3, V4). No Linux tool, `bwrap` or real srt is available locally
  (`WORLD_EXEC_SRT_NODE_MODULES` is unset), so every Linux and real-srt statement here is
  **UNMEASURED**. Prototypes ran on sibling `git archive` copies of `1835d8be` under
  `.designer-scratch-iter359/{head,proto,applycheck}`, never under `/tmp`, and were deleted afterwards. Raw outputs and scripts
  are banked as `design_docs/verification/world-iter359/designer-*.txt` (`$V` below).

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

On each accept, World reads at most one line's worth of bytes, one byte at a time, within
`execProbeNet6Wait` = 500 ms:

| Bytes read | Class | Counted as | Why |
|---|---|---|---|
| exactly the control line | control | — (World's own; before every window) | proves the listener accepts **and** the classifier reads |
| exactly the child line | **child** | a reach: the leg **fails at once** | the child, or a relay of the child's bytes, reached World |
| a byte that cannot begin either line (decided at that byte) | **foreign** | the window is **inconclusive**: retry | positive evidence that another protocol is speaking. The start page's `GET /` is foreign at its first byte (V3) |
| nothing: EOF, reset or deadline | **silent** | a reach: the leg **fails at once** | **the subtle case.** A proxy that connects upstream and then relays nothing, a child whose write is blocked, and a port scanner all look the same. Nothing proves it is not the child, so it fails closed |
| a strict prefix of a line, then EOF or deadline | **silent** | a reach: the leg **fails at once** | a truncated child line is not evidence of a foreign dialer |
| an accept still pending at `window()`'s bound | **silent** | a reach | unreachable by construction (the classifier's own deadline is shorter), and fails closed if it ever happens |

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

`net6Leg` judges one leg, attempt by attempt:

1. If the token is not `WORLD-PROBE-REFUSED …`, or the window has any **child** or **silent**
   accept, the leg **fails**. The text names the token, the attempt and the counts:
   `…: token "…" on attempt N; accepts in its window: child C, silent S, foreign F (want a refusal
   and no accept from the child or silent; <run>)`.
2. If the window has no accept at all, the leg **passes**.
3. Otherwise every accept in the window was foreign, so the window is **inconclusive**. Unless this
   was attempt `execProbeNet6Attempts` (3), the leg re-runs its own op alone, with a fresh window:
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

### 3.5 The patch (the planner pins these bytes: V13)

`exec_probe.go`, by function:

- **`execProbeScript`**: `connect(port,nonce)` writes `dial(nonce)` with `s.end` on connect.
  `proxy(port,nonce)` writes it through a 200 and passes the nonce through to `connect` when there is
  no proxy. Both ops parse `"<port>:<nonce>"`. The doc comment says so.
- **New constants** `execProbeDialPrefix`, `execProbeNet6Wait` (500 ms) and
  `execProbeNet6Attempts` (3), each with its reason.
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

`ProbeExecProfilesEach`, the row-242 hoist, the reverse leg's judgement, the arm order and every
other arm are untouched. The `newProbeListener` call site stays where T2 of row 242 pins it.

## 4. The by-construction argument

**Claim A (no new false pass).** Arm 6's raw or proxy leg passes only on an attempt whose window
held **no accept at all** and whose token is a refusal. That is HEAD's pass condition (`accepts`
still 1 after the control, with a refusal token), applied to that attempt's window.

*Proof.* `net6Leg` returns `""` in exactly one branch, `w.foreign == 0` after the
`child > 0 || silent > 0 || token not a refusal` branch has been ruled out. So `child`, `silent`
and `foreign` are all 0, and every accept in the window was counted as one of the three: `window`
maps control and pending to silent. ∎

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
  checks that classification never outlives its deadline.
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
- `ci.yml:247`: TF, TS, TE and TC are appended to the verbose `tests=` list, so that a Linux skip
  cannot read as a pass (the family's rule).

The mutation table was measured by `$V/designer-run-mutants.sh.txt` and `$V/designer-mutate.py.txt`
on the prototype (V11). Each mutant was applied to a pristine copy of the prototype's
`exec_probe.go`, then `TP`, TF, TS, TE and TC were run with `-count=1`, and the copy was restored
and checked with `cmp`. 1 = that test went red.

| Mutant | Edit | `TP` | TF | TS | TE | TC | Expected → observed |
|---|---|---|---|---|---|---|---|
| CONTROL | none | 0 | 0 | 0 | 0 | 0 | green → green |
| HEAD (positive control) | `1835d8be` + TF/TS/TE (TC needs new identifiers) | — | **1** (`accepts 2`) | **1** (texts only: HEAD refuses every row, fail-closed) | **1** | — | TF red → red (V8) |
| MUT-RAW-COUNT | the leg fails on **any** accept (HEAD's rule inside the new structure) | 0 | **1** | **1** (`raw-child-on-retry`) | **1** | 0 | red → red |
| MUT-SILENT-AS-FOREIGN | the classifier's last `return net6Silent` → `net6Foreign` | 0 | 0 | **1** (`raw-silent`, `raw-scan`, `proxy-silent`) | 0 | **1** (4 rows) | red → red |
| MUT-PREFIX-AS-FOREIGN | a non-empty strict prefix → foreign | 0 | 0 | 0 | 0 | **1** (2 prefix rows) | red (TC only) → red |
| MUT-CHILD-AS-FOREIGN | the exact child line → foreign | **1** | 0 | **1** (4 child rows) | 0 | **1** | red → red |
| MUT-CHILD-AS-CONTROL | the exact child line → control | **1** | 0 | **1** (texts: arm 6 still refuses, with `silent 1`, 11 of 11) | 0 | **1** | red → red; still fail-closed by `window`'s mapping |
| MUT-NO-NONCE (as scripted) | a `WORLD-PROBE-DIAL ` prefix stops the early foreign verdict, and a non-line tail reads foreign | 0 | 0 | 0 | 0 | **1** (2 prefix rows) | **as measured it degenerates to MUT-PREFIX-AS-FOREIGN**, so it is reported, not counted as a separate kill |
| MUT-TOKEN-ONLY | `net6Leg` ignores the window (the row's direction (a)) | **1** | **1** | **1** (all 7) | **1** | 0 | red → red |
| MUT-RETRY-CHILD | a child window is retried like a foreign one | 0 | 0 | **1** (4 child rows: arm 6 = `""`, **a false pass**) | 0 | 0 | red → red |
| MUT-RETRY-SILENT | a silent window is retried | 0 | 0 | **1** (3 silent rows) | 0 | 0 | red → red |
| MUT-EXHAUST-PASSES | the third inconclusive window passes | 0 | 0 | 0 | **1** | 0 | red → red |
| MUT-UNBOUNDED | no attempt bound | 0 | 0 | 0 | **1** | 0 | red → red |
| MUT-FROM-NOT-ADVANCED | windows never advance (`*from` kept) | **1** | **1** | 0 | 0 | 0 | red → red |
| MUT-WINDOW-AT-JUDGE | the window starts at judgement time, not before the run | **1** | **1** | **1** (all 7) | **1** | 0 | red → red |
| MUT-CLIENT-NO-LINE | the client's `connect` destroys instead of writing its line | **1** | 0 | 0 | 0 | 0 | red (`TP`'s `child 1`) → red |
| MUT-PENDING-AS-FOREIGN | `window` counts a still-pending accept as foreign | 0 | 0 | 0 | 0 | 0 | **SURVIVES.** Unreachable by construction: `window`'s bound (1.5 s) exceeds the classifier's own read deadline (500 ms), so nothing is pending at the bound. Killing it would need a seam that stalls the classifier; none is added |
| MUT-CONTROL-ANY-ACCEPT | `control()` returns on any accept, not on its own line | 0 | 0 | 0 | 0 | 0 | **SURVIVES.** It weakens only the control's proof that the classifier reads a line. That proof is pinned independently: `TP`'s `child 1, silent 0` (an unsandboxed child classified child through World's real node, on both legs) and TC. Under the mutant, a foreign accept between the hoisted bind and `control()` (row 242 V15: milliseconds) would satisfy the control. That still proves the listener accepts; only the classifier-read proof is lost, and `TP` and TC pin that proof |

Every mutant that would turn a reach into a pass (MUT-RETRY-CHILD, MUT-RETRY-SILENT,
MUT-SILENT-AS-FOREIGN, MUT-TOKEN-ONLY, MUT-EXHAUST-PASSES, MUT-UNBOUNDED, MUT-WINDOW-AT-JUDGE) is
killed by TS or TE. Every mutant that would bring the false refusal back (MUT-RAW-COUNT,
MUT-FROM-NOT-ADVANCED, MUT-TOKEN-ONLY, MUT-WINDOW-AT-JUDGE) is killed by TF. The two survivors are
stated above with their reasons.

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
  tries to land a clean window instead of one. No srt path is known to alter tunnel bytes, and a
  deterministic one still refuses (TE). The alternative, failing at once on foreign, would keep the
  start page's false refusals. The design takes the reverse leg's precedent (row 230: foreign →
  retry, bounded) over that.
- **R3, latency.** A clean probe adds no wait: `window()` returns as soon as nothing is pending, and
  a foreign accept is classified at its first byte. Each inconclusive window adds one small run plus
  the 200 ms settle. That measured about 0.3 s per retry (`raw` 1.50 s against `none` 1.21 s, V9).
  A silent accept costs up to 500 ms before the leg refuses.
- **R4, real srt and Linux (UNMEASURED here).** Under real srt the raw connect is refused and the
  proxy's CONNECT to `127.0.0.1` is denied, so no accept lands and the window is empty: the same
  verdict as HEAD. The client's new `s.end(line)` only runs after a successful connect.
  `TestExecSrtStartupProbePasses` (already on the `ci.yml:247` list) is the first real-srt run of
  this patch. It skips locally (`bwrap` is absent and `WORLD_EXEC_SRT_NODE_MODULES` is unset).
- **R5, a nonce in argv.** Another local process could read the child nonce from `ps` and send the
  child line, which only makes arm 6 **refuse** (child). It cannot make it pass. The control nonce is
  never in any argv.

## 7. Milestones

- **M1 (single, ~0.5 day).** Apply `$V/designer-proto.patch.txt`. `git apply` reproduces all five
  files byte-exact (V13). Then prove, each test run on its own:
  - TF is red at the parent commit (the HEAD row, V8) and green on the patch;
  - the mutants in §5 reproduce their kill pattern: at minimum MUT-RAW-COUNT, MUT-SILENT-AS-FOREIGN,
    MUT-TOKEN-ONLY, MUT-RETRY-CHILD, MUT-RETRY-SILENT and MUT-EXHAUST-PASSES;
  - the edited `TP` and `TestExecProbeNet6BoundBeforeAnyPick` are green.

  Then run `./scripts/verify_ail.sh` (no `.ail` changes) and `go build ./... && go vet ./host/broker/
  && go test ./...` under the pinned `AILANG_BIN`.

## 8. Files

| File | Change |
|---|---|
| `host/broker/exec_probe.go` | +197/−33: the client script's `connect`/`proxy`; three constants; `probe`'s arm-6 legs; the `reverseLeg` parameter; `probeListener` rewritten (nonces, classified accepts, `arg`/`mark`/`window`/`control`); new `net6Class`, `net6Window`, `net6Leg`, `classifyNet6Dial` |
| `host/broker/exec_probe_net6_dial_test.go` | new: TF, TS, TE, TC and the `net6DialClient` rig |
| `host/broker/exec_probe_test.go` | +3/−3: `TP`'s two leg strings; `runReverseLeg`'s argument |
| `host/broker/exec_probe_net6_test.go` | +1/−1: the connect-argument regex |
| `.github/workflows/ci.yml` | +1/−1: four names appended to the `tests=` list at `:247` |

## 9. What this item is NOT doing

It does not change the reverse leg (retry, reset, hello wait: rows 230, 232 and 233), the row-242
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
| V12 | Stability and the arm-6 subset on the final prototype | `go test -count=3 -run '^(TestExecProbePassThroughShimFailsEveryConfinementArm\|TestExecProbeNet6.*)$' -v`; `go test -count=1 -run ExecProbe ./host/broker/`; `go vet ./host/broker/` → `$V/designer-stability.txt` | 3/3 PASS for each of the 7 tests (including row 242's T1/T2 with the new regex); `ok … 70.415s`; ExecProbe subset `ok … 46.682s`; `vet-ok` |
| V13 | The patch is byte-exact | `git apply $V/designer-proto.patch.txt` on a `git archive` of `1835d8be`, then `cmp` against the prototype files → `$V/designer-instruments.txt` | `APPLY-OK`; `exec_probe.go` `43bfd807641f4fd1`, `exec_probe_test.go` `ee59ba2cc5742855`, `exec_probe_net6_test.go` `8ec970891599257a`, `exec_probe_net6_dial_test.go` `f4f98eb35246b6f8`, `ci.yml` `4efb18438d736bb4`; patch sha256 `4793fa7d2d108932`; numstat `exec_probe.go` +197/−33 |
| V14 | The tests the mutants ran are the patch's tests | `cmp` of the proto copy's three test files against the patch tree | `same` ×3. The assertion relaxation (V10) came before the final mutation run and the patch |
| V15 | `ci.yml:247` is the list the names join | `grep -n "TestExecProbeNet6" .github/workflows/ci.yml` → `$V/designer-greps.txt` | `:247 tests="$tests TestExecProbeReverseLegJudge …`, ending `TestExecProbeNet6BindSiteIsBeforeTheFanOut"` |
| — | **UNMEASURED** | — | the full `go test ./...` under the pinned `AILANG_BIN` (the executor runs it); Linux; real srt (R4); silent foreign dialers on other hosts (R1); a startup-level false-refusal count, before or after |

## 11. Quorum log

(filled by the controller)
