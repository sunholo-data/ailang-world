# w-exec-probe-holder-exit-reset — a reset that FOLLOWS a foreign answer on the same attempt's port is that holder's exit, so it is retried, not refused

- **Status**: Planned (designed iteration 330, r0). Queue row 233. An amendment to
  `design_docs/implemented/w-exec-probe-reverse-leg-port-race.md` (row 229, the reverse leg's
  §4.2 classification) and the successor named by
  `design_docs/implemented/w-exec-probe-reset-holder-false-pass.md` (row 232) §3 reading 3 and §8 R5.
  A sibling of `design_docs/planned/w-exec-probe-reverse-leg-pins.md` (row 230).
- **Item**: Clause 3 (MET; availability hygiene). When a foreign chatty holder exits while the
  reverse leg's poll is still dialling, the dial it cuts reads zero bytes and a reset. Today that
  dial is SILENT, so the attempt is a sticky refusal of startup, not INCONCLUSIVE (a retry).
- **Estimated**: ~0.25 day. Production about +25/−12 LOC in `host/broker/exec_probe.go` (one
  small enum, the classifier returns it, the poll's switch gains one case). Tests about +100 LOC
  (one new table test). One `ci.yml` list edit. One QUICKSTART clause (§7).
- **Measurement base**: `a6af0cd4` (worktree `sprint/world-iter330-exec-probe-holder-exit-reset`;
  `host/broker/exec_probe.go` sha256 prefix `96f33172667d27de`, V1).
- **Instruments**: `go1.26.6 darwin/arm64`, 16 cores, node v26.12.1 (V1). Real srt is the pinned
  `@anthropic-ai/sandbox-runtime` 0.0.78, installed outside the repo at
  `~/.ailang/state/world-iter330-srt/`; its `cli.js` hash equals the pin (V4). Load means 32
  spinning `yes` on 16 cores on darwin and 2×nproc = 8 on the 4-vCPU Linux runner. The runner
started and killed them, with 0 alive afterwards (V10, V11).
  Linux was measured on a CI runner from a throwaway branch (V11). Probe sources, runners and raw
  outputs are under `design_docs/verification/world-iter330/` as `.txt` files (`V=` below). Scratch
  tests were copied into `host/broker/` as `zz_iter330_*_test.go`, run, and deleted. The prototype
  production change was applied, measured, and reverted with `git checkout`; the base hash was
  re-checked afterwards (V12).

---

## 1. The problem, and what this session measured

The reverse leg (arm 6) runs the child's `listen` op under srt and polls the chosen loopback port
from the host until the child's run returns (`exec_probe.go:598-600`, V3). Each connected dial is
classified (`classifyReverseDial`, `:692-720`): the nonce line is **child**, other bytes or a clean
EOF is **foreign**, and zero bytes otherwise (the deadline, a reset, a failed hello write) is
**silent**. A connect-time ECONNRESET is also silent (`:757`, row 232). `judgeReverse` (`:150-162`)
puts silent ahead of foreign, so an attempt that saw both is SILENT, a sticky refusal (V2).

srt's mux is the one identified in-window foreign binder. It binds `127.0.0.1:0` per srt start
(V5) and exits with srt (V5). The poll stops only after srt exits (V3). So when the chosen port
collides with the mux, the poll sees HTTP 400 (foreign) for the whole run, and the dial in flight
when srt exits can be cut with a zero-byte reset (silent). That one dial turns a retryable
collision into a refused startup.

| Row 233 premise | Verdict at `a6af0cd4` | Evidence |
|---|---|---|
| Today `{foreign, silent}` is judged SILENT | **CONFIRMED** (controller, then first-party): foreign-only → 4 (`revInconclusive`), foreign+silent → 3 (`revSilent`), silent-only → 3 | V2 |
| The R5 sequence (foreign answers, then a zero-byte reset as the mux exits) happens with **real** srt | **CONFIRMED on darwin under load; not seen idle.** Real srt 0.0.78, hold 1000 ms: **0/100** cycles idle, **2/100** under load. Hold 0 (darwin's real child exits at once with EPERM): **0/100** idle, **5/100 (each a lone `S`: see the hold-0 caveat in R5)** under load. In the hold-1000 shape every silent was a **read-side reset (`S`) after ≥19 foreign answers, ≤8 ms before srt's exit, and the last connected dial of its cycle**. The hold-0 silents were lone dials (`seq S`); that is an artefact of the probe starting its poll late (R5) | V9, V10 |
| Rate | Real srt is **lower** than the row-232 stand-in (3/100 idle, 13/100 load; V18). The stand-in over-states darwin idle | V9, V10, V18 |
| Linux | **MEASURED on a CI runner (`7.0.0-1017-azure`, 4 vCPU, node v22.23.4), and far worse than darwin.** Real srt, hold 1000 ms (Linux's real shape: the child binds in its own netns and holds): **0/100** idle, **54/100** under 8 `yes`, every one `F…FS` (a reset after 18–19 foreign answers, the last connected dial). The node mux stand-in (`mux-server.mjs`) gave 0/100 idle and **49/100** loaded, the same shape. Under the rule all of them are INCONCLUSIVE: **54 → 0** and **49 → 0** | V11 |
| `SO_REUSEPORT` might let the child co-hold the mux's port | **REFUTED for the mux and World's child, on darwin; and REFUTED on Linux too on Linux.** srt's dist has 0 `reusePort` uses; the mux binds positionally. On the live real-srt mux port, a Go default bind, a Go `SO_REUSEPORT` bind and the child's exact node `listen` shape all fail `EADDRINUSE`, while the paired control (two `SO_REUSEPORT` sockets) co-binds. Node refuses `reusePort` on darwin with `ENOTSUP` | V5, V6, V7, V11 |

## 2. The mechanism

srt's CLI exits from the child's `exit` handler with `process.exit` (`cli.js:428-447`, V5). The
kernel then closes the mux's listening socket and every socket it had accepted. A dial the poll has
in flight at that instant ends one of three ways, all observed (V9, V10):

- the mux had already answered: `F` (HTTP 400), **foreign**;
- the mux had accepted and read the hello but not answered: a FIN, so a clean EOF `f`, **foreign**
  (95/100 idle hold-0 cycles ended this way);
- the connection was still in the listen queue, or the hello was still unread in the kernel: an
  RST, so zero bytes and ECONNRESET on the read, `S`, today **silent**.

Load stretches the window between the handshake and the mux's answer. That is why, in the
hold-1000 shape, `S` appeared only under load: darwin 0 → 2 and Linux 0 → 54 per 100 (V9–V11). After the exit, the next dial is refused (`R`), or the poll has already
been stopped.

## 3. Design

### 3.1 The rule

> **A zero-byte RESET observed after the attempt's poll has already recorded a foreign answer
> records nothing.** A reset means ECONNRESET at connect, on the hello write, or on the read, with
> zero bytes read. Everything else is unchanged: a reset BEFORE the first foreign answer is silent,
> a zero-byte DEADLINE (or any other zero-byte error) is silent wherever it falls, and the child's
> nonce dominates.

`judgeReverse` is **not** changed. The rule lives in the poll, which is the only place that knows
the order of the dials.

**Ordering needs no new field.** `reverseSeen` is three bools and carries no order. It does not
need to. The poll goroutine is the only writer of `child`, `foreign` and `silent`, and it dials
sequentially (`:745-765`, V3). So "has this attempt already seen a foreign answer" is just
`foreign.Load()` read on that goroutine, in program order. No lock, no timestamp, no fourth bool.

### 3.2 Code shape (prototype measured: V12, V13)

```go
// dialSeen is what one reverse-leg dial observed.
type dialSeen int

const (
	dialNone    dialSeen = iota // the dial failed and proves no listener (ECONNREFUSED, EINVAL, a timeout, …)
	dialChild                   // the child's nonce line
	dialForeign                 // other bytes, or a clean EOF after a written hello
	dialSilent                  // zero bytes at the deadline, or another zero-byte error
	dialReset                   // zero bytes, then ECONNRESET at connect, on the hello write or on the read
)

func classifyReverseDial(c net.Conn, want string) dialSeen // body as today; the last arms become:
	// case werr == nil && errors.Is(rerr, io.EOF):                              return dialForeign
	// case errors.Is(werr, syscall.ECONNRESET) || errors.Is(rerr, syscall.ECONNRESET): return dialReset
	// default:                                                                  return dialSilent

// pollLoopback's loop body:
	d := dialNone
	if c, err := dial(addr); err == nil {
		d = classifyReverseDial(c, want)
		_ = c.Close()
	} else if dialErrIsSilent(err) {
		d = dialReset
	}
	switch {
	case d == dialChild:
		child.Store(true)
	case d == dialForeign:
		foreign.Store(true)
	case d == dialReset && foreign.Load():
		// Row 233: the foreign holder's exit. Records nothing; a deadline still refuses.
	case d == dialReset || d == dialSilent:
		silent.Store(true)
	}
```

The prototype (`$V/proto-exec_probe.diff.txt`, +34/−15) used the names `revDial*`. The sprint
uses `dial*`, because the test file already has a `revDialErr` helper (V14). `dialErrIsSilent`
keeps its name and meaning: a connect-time ECONNRESET proves a listener completed the handshake,
and that dial is now `dialReset`, so it is silent unless a foreign answer came first. No refusal
text changes. `classifyReverseDial` has no caller outside `exec_probe.go` (V14).

### 3.3 Soundness: which verdicts change, and why none of them can pass a breach

The probe is a fail-closed startup gate. The property that must not weaken is P20 (V46): a
listener the child binds is unreachable from the host, so a real reach must always refuse.

**Exactly one verdict changes.** Take one attempt, with child token `tok` and the poll's dials in
order. Today's `reverseSeen` is `{child:false, foreign:true, silent:true}`, and **every** silent
dial is a reset (connect, write or read ECONNRESET, zero bytes) that comes **after** the attempt's
first foreign answer. That attempt is SILENT today (sticky refusal). Under the rule it is
`{foreign:true}`, which is INCONCLUSIVE (a retry on a fresh port). That holds for `tok` =
`WORLD-PROBE-BOUND`, any `WORLD-PROBE-REFUSED <errno>`, and `EADDRINUSE`.

Nothing else changes:

| Attempt shape (poll order) | Today | After | Why unchanged / changed |
|---|---|---|---|
| any dial reads the nonce | REACHED | REACHED | the child is first in `judgeReverse` (`:152`) |
| `tok` is not a probe token | NO-TOKEN | NO-TOKEN | checked before silent (`:154`) |
| reset(s) only, no foreign (T7 `reset-holder`, row 232's connect reset) | SILENT | SILENT | no foreign answer precedes it |
| reset **before** the first foreign answer (`SF`, `CF`) | SILENT | SILENT | order is checked; an order-free rule is mutant M31 |
| a zero-byte **deadline** anywhere (`FD`, `DF`, `FSD`) | SILENT | SILENT | a deadline is a listener that is alive and not answering, the slow-child case row 229 R1 closed; excusing it is M32 |
| foreign, then reset(s) only (`FS`, `FC`, `FW`, `fS`) | SILENT | **INCONCLUSIVE** | the changed row |
| foreign only, or `EADDRINUSE` with nothing | INCONCLUSIVE | INCONCLUSIVE | — |
| nothing connected, sandbox token | PASS | PASS | — |

**Why the changed row cannot pass a breach.**

1. **No attempt the rule touches can pass itself.** Row 229's invariant (§4.3: an attempt passes
   only when *zero* dials connected) is untouched. The changed attempt had connected dials, so it
   is INCONCLUSIVE, never PASS. Three INCONCLUSIVE attempts refuse with
   `reverse leg: inconclusive after 3 attempts` (`:621`, `execProbeReverseAttempts = 3`, `:111`).
   The rule can delay a refusal into the next attempt. It cannot turn a refusal into a pass.
2. **The excused reset cannot be the child's listener on a co-held port.** For that, the child would
   have to hold a host-reachable listener on a port a foreign listener answered on, in the same
   attempt.
   - *Co-holding with srt's mux* (the only identified in-window binder) is refused. On the live
     real-srt mux port, the child's exact `listen` shape gets `EADDRINUSE`, and so do a Go default
     bind and a Go `SO_REUSEPORT` bind. The control, two `SO_REUSEPORT` sockets, does co-bind (V6,
     V7). srt never sets `reusePort` (V5). World's child binds positionally (`:94`), and node
     refuses `reusePort` on darwin anyway (`ENOTSUP`, V6). In a breach with a shared namespace the
     child's token is then `EADDRINUSE`, and that attempt is INCONCLUSIVE today, whatever the poll
     saw.
   - *Co-holding with some other foreign listener on darwin* is possible when one side binds the
     wildcard: Go's default sockets co-bind `0.0.0.0:p` with `127.0.0.1:p` in either order (V6).
     But every dial to `127.0.0.1:p` then reaches the `127.0.0.1` listener, 20/20 in both orders
     (V8). World's child binds `127.0.0.1`, so while it holds, the poll reaches the child and reads
     the nonce: REACHED.
   - *Linux* is stricter than darwin. Wildcard and `127.0.0.1` co-binds are refused in both orders, and so is every mixed `SO_REUSEPORT` pair. Only two sockets that **both** set `SO_REUSEPORT` co-bind (Go and node `reusePort:true`, V11). On the live real-srt mux port, all three binds fail `EADDRINUSE` (V11). On Linux the correct child also sits in its own netns, so the host cannot reach it at all (row 229 §4.3).
3. **A sequential child (foreign exits, then the child binds) is still caught.** Its listener
   answers the nonce (REACHED) or, if slow, reaches the deadline (silent, still refused under the
   rule). World's child resets host dials only **as it closes** (row 232 §3 reading 1). In every
   cycle where the child was reachable, the nonce was seen too (100/100 idle and loaded; row 232
   V5, V6). So the only breach the rule could hide is one that the poll observes **solely** through
   close-time resets after a foreign answer, for a whole 1000 ms hold at a 50 ms cadence. Such a
   breach already passes today in any attempt without a collision. The rule adds no new class of
   miss. It only adds the precondition "a foreign answer came first", and the leg then retries.
4. **The probe client is unchanged and still never answers foreign** (row 232 §4(a)). The rule
   depends on that assumption no more than today's code does: a child that answered foreign would
   already be retried today.

**What the rule relies on (stated, not enforced):** (a) the poll is single-goroutine and
sequential (the code at `:745-765`); (b) the child cannot co-hold a host-reachable `127.0.0.1`
listener with a foreign `127.0.0.1` listener without `SO_REUSEPORT` on both (V6, V7; Linux V11);
(c) World's probe client never sets `reusePort` (`:94`). If (b) ever failed, the co-held child
would receive dials and read the nonce, as in V8, so the failure would still be REACHED.

### 3.4 Rejected alternatives (measured where it mattered)

- **(A) Order-free:** any `{foreign, silent}` is INCONCLUSIVE, with no struct or poll change.
  Rejected. It also excuses a genuinely silent holder that later answered (`SF`), and a slow
  child whose first dials timed out. M31 is this rule and is killed by T12 `SF` and `CF` (V13).
- **(B) Any silent after foreign is foreign** (deadlines too). Rejected. A deadline means the peer
  is alive and silent, which is row 229 R1's slow-child case. Every exit-time silent observed with
  real srt was a reset (V9, V10), so (B) buys no availability over the chosen rule. M32 is (B) and
  is killed by T12 `FD` and `FSD` (V13).
- **(C) The chosen rule, plus "the reset must be the attempt's LAST connected dial".** Every observed
  exit reset was the last connected dial (`connectedAfterFirstSilent=0`, V9, V10, V11), so (C)
  would be sound too. It needs a pending state resolved at `stop()`, about 6 more LOC, and it
  defends only against a foreign holder that keeps resetting **after** it has answered. Point 3
  shows that holder cannot hide a breach. Not adopted. It is recorded here so a reviewer can ask
  for it.
- **(D) Count EPIPE on the hello write as a reset too.** Not adopted. EPIPE was never observed.
  The one failed write seen (Linux, `W`) was `write: connection reset by peer` followed by EOF, which
  is ECONNRESET and is covered by the rule (V11). There was no `O` class anywhere. EPIPE therefore
  stays silent, which fails closed.

## 4. Test plan

| Test | Kind | Needs node? | Kills | Measured? |
|---|---|---|---|---|
| **T12** `TestExecProbeReversePollHolderExitAfterForeign` (new) | deterministic: scripted dial seam, no sleep, no socket | no | M30, M31, M32, M34, M35; M33 together with T7 and T9 | **yes**, prototype: control green, each mutant red (V13) |
| T1 `TestExecProbeReverseLegJudge` (unchanged) | pure | no | `judgeReverse` must stay byte-identical: row `bound-silent-foreign` still wants SILENT, because the rule lives in the poll | green on the prototype (V12) |
| T7 `TestExecProbeReverseLegSilentHolderRefusesSticky` (unchanged) | real holders | no | M33, through `reset-holder` | M33 red (V13) |
| T9 `TestExecProbeReversePollRecordsConnectReset` (unchanged) | injected dial | no | M33, through `connect-reset` | M33 red (V13) |
| T2–T6, T8, T10, T11 (unchanged) | as today | as today | regression | all green on the prototype (V12) |

**T12 shape.** The prototype is `$V/zz_iter330_t12_test.go.txt`. `scriptConn` is a `net.Conn`
whose `Write` returns a scripted error and whose `Read` hands over a scripted answer one byte at a
time, then a scripted error (`io.EOF` when none). `scriptedDial(steps, want)` returns the steps in
order: `N` nonce, `F` HTTP 400, `f` clean EOF, `S` read `ECONNRESET`, `W` write `ECONNRESET` then
EOF, `D` `os.ErrDeadlineExceeded`, `C` connect `ECONNRESET` (`revDialErr`). After the last step it
returns ECONNREFUSED, and it closes an `all` channel once every step has been handed out. The test
starts `pollLoopback(40000, want, dial)`, waits on `all` (with a 10 s fatal guard), then calls
`stop()`. That is deterministic: the poll is sequential, and `stop()` waits for the in-flight
classify (row 230 V24). Each row asserts the exact `reverseSeen` **and**
`judgeReverse("WORLD-PROBE-BOUND", seen)`:

| Row | Want `reverseSeen` | Want verdict |
|---|---|---|
| `FS`, `FC`, `FW`, `fS` | `{foreign}` | INCONCLUSIVE |
| `FD`, `SF`, `CF`, `FSD` | `{foreign, silent}` | SILENT |
| `S` | `{silent}` | SILENT |
| `FSN` | `{child, foreign}` | REACHED |

The sprint names the rows by shape (`foreign-then-read-reset`, …). The letters here match the
prototype's output in V13.

**AC-MUT.** Apply M30–M35 one at a time. Each reds its named killer, and the tree is restored
byte-identical (sha256 of both Go files). **AC-CI**: T12 joins step 8's `tests=` list
(`ci.yml:247`, V16), since the family's rule is that every reverse-leg test is PASS-guarded and
cannot skip on Linux. **AC-STABLE**: `go test ./host/broker -run '^TestExecProbeReverse' -count=20`
passes 20/20, and T12 `-count=500` under 32 `yes` gives 0 FAIL. **AC-GATE**: `gofmt -l` is
empty, `go vet ./host/broker/` is clean, `go build ./... && go test ./...` passes, and
`./scripts/verify_ail.sh` passes. No `.ail` is touched. **AC-LINUX**: the fix PR's linux/amd64
step 8 shows `--- PASS: TestExecProbeReversePollHolderExitAfterForeign`.

## 5. Mutation table (continues the parent's M1–M21, row 232's M22–M29, the judge's O1–O10, row 230's K1–K9)

| ID | Mutant | Killer, and the assertion that fires | Run? |
|---|---|---|---|
| M30 MUT-EXIT-RESET-SILENT | the rule's case removed (**today's behaviour**) | T12 `FS`, `FC`, `FW`, `fS`, `FSN`: `seen {… foreign:true silent:true} verdict 3; want {… foreign:true silent:false} verdict 4` | V13 red |
| M31 MUT-EXIT-RULE-ORDER-FREE | a reset is excused if the attempt **ever** saw foreign (option A) | T12 `SF`, `CF`: `verdict 4; want … silent:true verdict 3` | V13 red |
| M32 MUT-DEADLINE-EXCUSED | a deadline after foreign is excused too (option B) | T12 `FD`, `FSD` | V13 red |
| M33 MUT-EVERY-RESET-EXCUSED | a reset records nothing even with no foreign before it | T12 `S`, `SF`, `CF`; T7 `reset-holder`; T9 `connect-reset` | V13 red |
| M34 MUT-CONNECT-RESET-NOT-RESET | a connect-time ECONNRESET maps to `dialSilent` | T12 `FC` | V13 red |
| M35 MUT-WRITE-RESET-NOT-RESET | only the read's ECONNRESET counts as a reset | T12 `FW` | V13 red |

## 6. Milestones

- **M1** (one milestone): `dialSeen`, the classifier's return, the poll's switch, T12, the `ci.yml:247`
  gain, and the QUICKSTART clause. Prove M30–M35 red one at a time.

## 7. Files

| File | Change | LOC (est.) |
|---|---|---|
| `host/broker/exec_probe.go` | `dialSeen` (near `:686`); `classifyReverseDial` returns it (`:687-720`, doc comment updated); `pollLoopback`'s loop body (`:747-759`) and its doc comment (`:727-732`) | +25 / −12 |
| `host/broker/exec_probe_test.go` | `scriptConn`, `scriptedDial`, T12 | +100 |
| `.github/workflows/ci.yml` | `:247` `tests=` gains `TestExecProbeReversePollHolderExitAfterForeign` | ±1 |
| `docs/QUICKSTART.md` | §10 arm 6 (`:1168-1170`): after "(or that resets the host's dial)", add "; a reset right after a foreign answer on the same port is that process exiting, and is retried" | ±1 |

**S7** applies: the refusal-text paragraph describes which resets refuse. **S2/S3**: this is a
repair inside the host's effect boundary, with no new surface. It is not a package, for the same
reason as rows 229, 232 and 230. **S1/S4/S5** do not apply: there is no `.ail`. No dependency is
added.

## 8. Conflict surface, risks and residuals

**Conflict surface.** Two hunks of `exec_probe.go` (`:686-720`, `:745-759`) and their comments. Arms
1–5, 7 and 8, `judgeReverse`, `reverseLeg`, every constant and every refusal text are untouched.
`pollLoopback` has one production call site (`:598`) and two test call sites (`:829`, `:908`;
V14). Its signature is unchanged.

- **R1, the rate after the fix is 0 for the observed signature; the collision itself is untouched.**
  In every hold-1000 run, darwin and Linux, real srt and stand-in alike, each exit silent was a
  reset after foreign: darwin load 2 → 0, Linux load 54 → 0, Linux stand-in 49 → 0 (V10, V11). So
  the rule removes the false refusal it targets. It does not make the collision rarer (row 229 R4 watch), and a foreign holder
  that is silent **before** it answers still refuses.
- **R2, the sticky guarantee is narrowed by exactly one shape** (§3.3), and the narrowing is to a
  retry, not to a pass. A breach visible **only** as close-time resets after a foreign answer gets
  another attempt. §3.3 point 3 argues that this breach is already invisible to today's code in any
  uncollided attempt.
- **R3, platform residuals.** darwin's wildcard co-bind (V6) is handled by routing (V8).
  Linux co-binding needs `SO_REUSEPORT` on both sockets (V11). Neither srt (V5) nor World's child (V17) sets it. A future child script that added `reusePort` would still be caught by routing: the kernel spreads connections across the group, so the child would receive dials and answer the nonce. That is argued, not measured, and so it is a review point for any edit to `execProbeScript`'s `listen` op.
- **R5, an exit reset that is the attempt's FIRST connected dial stays a sticky refusal.** The
  rule needs a foreign answer first, and nothing can tell that lone reset apart from a resetting
  holder (T7 `reset-holder`). The hold-0 runs show such lone resets often: darwin load 5/100, Linux
  idle 45/99, Linux load 66/99. **These numbers are a probe artefact, not a production rate.** The
  probe starts polling only after the child prints the proxy port (darwin) or after `ss` finds it
  (Linux). With hold 0 the child exits at once, so the first dial often lands on the exit itself.
  Production starts the poll **before** srt (`:598` before `:599`, V3), so it dials across the
  mux's whole life. On Linux that life spans the child's 1000 ms hold, and ≥18 foreign answers came
  first in every loaded cycle (V11). On darwin the correct child refuses its bind at once (`EPERM`),
  so the mux lives only from srt's start through the proxy op to the exit. How many dials land in
  that window is **UNMEASURED**. If it is often zero or one, a darwin collision with the mux can
  still refuse. The direction is fail-closed. Trigger: the row-230 §3.5 R4 trigger
  (`answered nothing to World's hello` in a non-test log). Measurement path: drive the production
  `listen` child under real srt with a poll started before srt, which needs the mux port known in
  advance (not possible with srt's `bind(0)`). Successor candidate only if the trigger fires.
- **R4, EPIPE and other write errors stay silent** (option D, not adopted). If CI ever shows an
  exit with that shape, the R4 trigger of row 230 §3.5 (`answered nothing to World's hello` in a
  non-test log) fires, and the class is re-measured first.
- **Carried forward unchanged:** row 232 R6 (darwin EINVAL) and row 230 R1 (the cross-profile
  ordering). That successor's hoist does not interact with this rule.

## 9. What this item is NOT doing

- No change to `judgeReverse`, `reverseSeen`, the attempt count, the hold, the cadence, the hello
  wait or any refusal text.
- No new listener options, no port range, no lock.
- No claim about holders other than srt's mux beyond §3.3's argument. Each excused reset needs a
  foreign answer first on the same port.

## 10. Verification Log

The base for every row is `a6af0cd4`. `V=design_docs/verification/world-iter330`,
`S=~/.ailang/state/world-iter330-srt/node_modules/@anthropic-ai/sandbox-runtime`,
`P=~/.ailang/state/world-iter330-srt/probe`. `AILANG_BIN=$HOME/.pinned-ailang/ailang` is set for
every `go test`.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Base and tooling | `git rev-parse HEAD`; `git status --porcelain`; `shasum -a 256 host/broker/exec_probe.go`; `go version`; `node --version`; `sysctl -n hw.ncpu` | `a6af0cd4…`; clean apart from this item's evidence dir; `96f33172667d27de…`; `go1.26.6 darwin/arm64`; `v26.12.1`; `16` |
| V2 | **Premise: today `{foreign, silent}` is SILENT** | VERIFIED BY CONTROLLER (scratch test at `a6af0cd4`), then first-party: `$V/zz_iter330_judge_test.go.txt` copied in, `go test ./host/broker/ -run '^TestZZIter330JudgePremise$' -count=1 -v`, then deleted → `$V/out-judge-premise.txt` | `foreign-only => 4`, `foreign+silent => 3`, `silent-only => 3`, control `control-none => 0`; `revPass=0 revReached=1 revNoToken=2 revSilent=3 revInconclusive=4`; `ok … 0.331s` |
| V3 | Today's code: the judge, the classifier, the poll, the ordering | `grep -n 'func judgeReverse\|case seen.silent:\|case seen.foreign \|\|type reverseSeen\|func classifyReverseDial\|case werr == nil && errors.Is(rerr, io.EOF):\|return false, false // silent\|func pollLoopback\|var child, foreign, silent atomic.Bool\|} else if dialErrIsSilent(err) {\|poll := pollLoopback\|r := pp.run(ctx, argv, false\|seen := poll()\|inconclusive after\|const execProbeReverseAttempts' host/broker/exec_probe.go` | `:111 execProbeReverseAttempts = 3`; `:134 type reverseSeen struct{ child, foreign, silent bool }`; `:150 func judgeReverse`; `:156 case seen.silent:` before `:158 case seen.foreign \|\| tok == "WORLD-PROBE-REFUSED EADDRINUSE":`; `:598 poll := pollLoopback(…)`, `:599 r := pp.run(…)`, `:600 seen := poll()`; `:621 … inconclusive after %d attempts …`; `:692 func classifyReverseDial(c net.Conn, want string) (child, foreign bool)`; `:716` EOF arm; `:719 return false, false // silent`; `:733 func pollLoopback`; `:737` the three atomics; `:757 } else if dialErrIsSilent(err) {`. The loop (`:745-765`) dials, classifies, stores and waits in one goroutine |
| V4 | Real srt is the pin | `npm install --no-audit --no-fund "@anthropic-ai/sandbox-runtime@0.0.78"` in `~/.ailang/state/world-iter330-srt`; `shasum -a 256 $S/dist/cli.js`; `jq -r .version $S/package.json`; `jq -r '.version, .cli_sha256' design_docs/verification/world-row140-m0/pin.json` | `3c3092bd26b3924046f38d793c716b513dde619cf792ec50b181e2c7cd40d96e`, `0.0.78`; pin `0.0.78`, `3c3092bd…d96e`: equal |
| V5 | **srt source: the mux binds positionally on `127.0.0.1:0`, never `reusePort`, and srt exits with the child** | `grep -rn -c reusePort --include='*.js' $S/dist` (summed); control `grep -rn -c 'listen(' …` (summed); `grep -rln 'SO_REUSEPORT\|exclusive: *false' --include='*.js' $S/dist \| wc -l`; `grep -n "mux.server.listen(p, '127.0.0.1')" $S/dist/sandbox/sandbox-manager.js`; `grep -n "const portRange = getPlatform() === 'windows'" …`; `grep -n "const \[lo, hi\] = range ?? \[0, 0\]\|doListen(range ? port : 0)" $S/dist/sandbox/listen-in-range.js`; `grep -n "if (!isWindows) {\|opts.httpServer.listen(path" $S/dist/sandbox/mux-proxy.js`; `grep -n "child.on('exit'\|process.exit(code ?? 0)" $S/dist/cli.js` | `reusePort-hits 0`; control `listen(-hits 9`; `0` files. `sandbox-manager.js:433 … listenInRange(mux.server, p => mux.server.listen(p, '127.0.0.1'), portRange, …)`; `:703` portRange is set only on Windows; `listen-in-range.js:10 const [lo, hi] = range ?? [0, 0];`, `:38 doListen(range ? port : 0);`. The HTTP backend is a Unix socket off Windows (`mux-proxy.js:111`, `:118 opts.httpServer.listen(path, …)`), so the mux is srt's only TCP binder there. `cli.js:428 child.on('exit', …)` → `:447 process.exit(code ?? 0);` |
| V6 | **Co-bind controls, darwin** (`$V/cobindctl.go.txt` in `$V/srtexit.go.txt`'s package; constants `$V/reuse-consts.go.txt`) | `cd $P && ./srtexit -mode cobindctl` → `$V/out-cobindctl-darwin.txt` | Go: `default-then-default`, `default-then-reuseport`, `reuseport-then-default` → `bind: address already in use`; **`reuseport-then-reuseport -> CO-BOUND`** (the positive control); `wildcard-default-then-loopback-default -> CO-BOUND`; `loopback-default-then-wildcard-default -> CO-BOUND`. Node: `reusePort first=0 second=0 -> second EADDRINUSE`; `first=1 second=1 -> first ENOTSUP`; `first=0 second=1 -> second ENOTSUP` |
| V7 | **Co-bind on a live REAL srt mux port** (cycle 1 of V9) | `$V/srtexit.go.txt` `coBindReport` while srt runs | `COBIND go-default 127.0.0.1:54255 -> … address already in use`; `COBIND go-SO_REUSEPORT 127.0.0.1:54255 -> … address already in use`; `COBIND node-child-shape EADDRINUSE` (node `net.createServer().listen(port,'127.0.0.1')`, the child's shape) |
| V8 | darwin routes a `127.0.0.1` dial to the `127.0.0.1` co-holder, not the wildcard one | `$V/route.go.txt`: `go run .` → `$V/out-route.txt` | `first=0.0.0.0 second=127.0.0.1 … -> map[second:127.0.0.1:20]`; `first=127.0.0.1 second=0.0.0.0 … -> map[first:127.0.0.1:20]` |
| V9 | **Real srt, darwin, idle.** srt runs `/bin/sh -c 'echo "PX:$HTTP_PROXY"; sleep H'` with `$V/srt-settings.json.txt`; the probe reads the mux port from `PX:`, runs pollLoopback's exact loop (200 ms dial, hello, 500 ms read, 50 ms cadence) and stops it after srt exits. Classes: `F` bytes, `f` clean EOF, `S` read reset, `C` connect reset, `W` write fail, `D` deadline, `R` refused | `$V/run-darwin.sh.txt` → `$V/out-darwin.txt` (`=== IDLE hold=1 n=100`, `=== IDLE hold=0 n=100`) | hold 1: `CLASSES map[F:2183 R:4 f:1]`; `anySilent=0`; `STICKY-SILENT-CYCLES today=0`. hold 0: `CLASSES map[F:5 f:95]`; `anySilent=0`; `today=0` |
| V10 | **Real srt, darwin, under 32 `yes`** | same runner (`=== LOAD …`) | `load: { 6.37 9.08 6.57 }` at start, `{ 46.12 31.51 17.55 }` at end. hold 1: `CLASSES map[F:2151 R:8 S:2 f:2]`; `anySilent=2 silentAfterForeign=2 silentBeforeForeign=0 multiSilent=0 connectedAfterFirstSilent=0 deadlineSilentAfterForeign=0`; `STICKY-SILENT-CYCLES today=2 … ruleD(reset-after-foreign)=0`; the two cycles `seq FFFFFFFFFFFFFFFFFFFS` (`S at exit-7.6ms`) and `FFFFFFFFFFFFFFFFFFFFFS` (`exit-3.3ms`). hold 0: `CLASSES map[F:33 R:1 S:5 f:62]`; `silentOnly=5`; `today=5 … ruleD=5`; each `seq S` (a lone dial, `exit-4.9` to `-20.7ms`). `yes-alive-after=0`. The runner was rebuilt mid-run to add Linux port discovery (a branch dead on darwin); the banked source is the final one |
| V11 | **Linux** (throwaway branch `probe/world-iter330-holderexit`, orphan commit `b410b9d7`, workflow `$V/linux-probe-workflow.yml.txt`) | Pushed, polled to `completed success` (run 41185322456, job 124490964931), log fetched with `gh api --allow-escape-sequences "repos/sunholo-data/ailang-world/actions/jobs/124490964931/logs"` (82,169 bytes raw), ANSI stripped → `$V/linux-probe-run41185322456.txt`. Branch then deleted: `git push origin --delete probe/world-iter330-holderexit`; `git ls-remote origin 'refs/heads/probe/world-iter330-*' \| wc -l` → `0` (control `refs/heads/dev` → `1`). On Linux the probe finds srt's host mux port by pid (`ss -Hltnp`), since the child's `HTTP_PROXY` names the in-netns bridge (`PX:***localhost:3128`) | `Linux … 7.0.0-1017-azure … x86_64`, `4`, `go1.26.6 linux/amd64`, node `v22.23.4`; srt sha256 = pin. COBINDCTL: Go `default-then-default`, `default-then-reuseport`, `reuseport-then-default` → in use; **`reuseport-then-reuseport -> CO-BOUND`**; `wildcard-default-then-loopback-default` and `loopback-default-then-wildcard-default` → **in use**; node `0/0 -> second EADDRINUSE`, **`1/1 -> CO-BOUND`**, `0/1 -> second EADDRINUSE`. Live mux: `COBIND go-default … in use`, `go-SO_REUSEPORT … in use`, `node-child-shape EADDRINUSE`. IDLE srt hold 1: `CLASSES map[F:1999]`, `today=0`. IDLE srt hold 0: `CLASSES map[C:1 R:54 S:43 W:1]`, `total=99 anyForeign=0 silentOnly=45`, `today=45 … ruleD=45`; `W` = `write: connection reset by peer r=EOF`; 1 `NO-PORT` cycle. IDLE stand-in: `F:2000`, `today=0`. Load `0.56 0.42 0.19` before the burners. **LOAD srt hold 1: `CLASSES map[F:1910 R:12 S:54]`; `anySilent=54 silentAfterForeign=54 silentBeforeForeign=0 multiSilent=0 connectedAfterFirstSilent=0 deadlineSilentAfterForeign=0`; `today=54 ruleA=0 ruleB=0 ruleD=0`**; shapes 53× `F{19}S`, 1× `F{18}S`. LOAD srt hold 0: `CLASSES map[F:1 R:24 S:66 f:8]`, `silentOnly=66`, `today=66 … ruleD=66`. **LOAD stand-in: `F:1942 R:4 S:49 f:2`; `silentAfterForeign=49`; `today=49 … ruleD=0`.** `load-end: 10.53 7.54 3.48`; `yes-alive-after=0`; `PROBE-DONE` |
| V12 | **The prototype, and the whole ExecProbe suite on it** | prototype applied (`$V/proto-exec_probe.diff.txt`); `$V/zz_iter330_t12_test.go.txt` copied in; `go vet ./host/broker/`; `go test ./host/broker/ -run '^(TestExecProbe\|TestZZIter330)' -count=1 -v` → `$V/out-proto-suite.txt`; then `git checkout -- host/broker/exec_probe.go`, test deleted, `shasum -a 256 host/broker/exec_probe.go`, `git status --porcelain` | vet clean; `22 --- PASS`, `0 --- FAIL`, `1 --- SKIP` (`TestExecProbeArm8`, real srt needs CI's env); T1–T11 and the T12 prototype all PASS; `ok … 36.160s`. Restored `96f33172667d27de`; porcelain shows only `?? design_docs/verification/world-iter330/` |
| V13 | **T12 prototype against M30–M35** | `bash $V/run-mutants.sh.txt` → `$V/out-mutants.txt` (the step-8 reverse-leg tests T1, T2, T7, T8, T9, T10 plus T12) | `CONTROL (prototype): 44 PASS lines`. M30: `FC FS FSN FW fS` FAIL, `dials FS: seen {… foreign:true silent:true} verdict 3; want {… silent:false} verdict 4`. M31: `CF SF` FAIL, `dials SF: … verdict 4; want … silent:true verdict 3`. M32: `FD FSD` FAIL. M33: `…SilentHolderRefusesSticky/reset-holder`, `…PollRecordsConnectReset/connect-reset`, T12 `CF S SF` FAIL. M34: `FC` FAIL. M35: `FW` FAIL. `restored-to-prototype a3568e76a6709599` |
| V14 | Census | `grep -rn classifyReverseDial --include='*.go' host cmd \| grep -v host/broker/exec_probe.go \| wc -l`; `grep -rn 'pollLoopback(' --include='*.go' host cmd`; `grep -rn revDialErr --include='*.go' host cmd \| wc -l` | `0`; `exec_probe.go:598` (`reverseLeg`), `:733` (def), `exec_probe_test.go:829` (T9), `:908` (T10); `12` (the test helper, so the sprint avoids the `revDial*` prefix) |
| V15 | QUICKSTART's arm-6 text | `sed -n 1164,1172p docs/QUICKSTART.md` | `:1168-1170` "`reverse leg: a listener on … answered nothing to World's hello` means a process that never answers (or that resets the host's dial) held the port, and World cannot tell it from the child: find it." |
| V16 | CI step 8's list | `grep -n -e 'TestExecProbe' .github/workflows/ci.yml` | `:247` lists T1–T11 (`…TestExecProbeReverseHelloWaitBounds TestExecProbeChildListenerDrainsTheHello`) |
| V17 | World's child binds positionally (no `reusePort`) | `grep -n "v.listen(+port,'127.0.0.1'" host/broker/exec_probe.go` | `:94 v.listen(+port,'127.0.0.1',()=>{say('WORLD-PROBE-BOUND');…` |
| V18 | The stand-in rate this row inherited | read `design_docs/verification/world-iter318/designer-probe-outputs.txt` via row 232 V5/V6 | mux-shaped stand-in, 1000 ms, faithful: idle `pre=3 post=3`; load `silent-postconnect 13`, `pre=13 post=13` |

## 11. Quorum log

r0: not yet reviewed.
