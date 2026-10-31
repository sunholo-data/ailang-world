# w-exec-probe-holder-exit-reset — a reset within one hold of a foreign answer, in a cleanly exited run, is that holder's exit, so it is retried, not refused

- **Status**: Planned (designed iteration 330, r1: quorum r1 fixes, §11). Queue row 233. An amendment to
  `design_docs/implemented/w-exec-probe-reverse-leg-port-race.md` (row 229, the reverse leg's
  §4.2 classification) and the successor named by
  `design_docs/implemented/w-exec-probe-reset-holder-false-pass.md` (row 232) §3 reading 3 and §8 R5.
  A sibling of `design_docs/planned/w-exec-probe-reverse-leg-pins.md` (row 230).
- **Item**: Clause 3 (MET; availability hygiene). When a foreign chatty holder exits while the
  reverse leg's poll is still dialling, the dial it cuts reads zero bytes and a reset. Today that
  dial is SILENT, so the attempt is a sticky refusal of startup, not INCONCLUSIVE (a retry).
- **Estimated**: ~0.25 day. Production about +35/−15 LOC in `host/broker/exec_probe.go` (one
  small enum, the classifier returns it, the poll's switch gains one timed case, `stop` takes
  `runClean`, `reverseLeg` passes `r.err == nil`). Tests about +150 LOC (one new table test; two
  call sites pass `true`). One `ci.yml` list edit. One QUICKSTART clause (§7).
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
| Linux | **MEASURED on a CI runner (`7.0.0-1017-azure`, 4 vCPU, node v22.23.4), and far worse than darwin.** Real srt, hold 1000 ms (Linux's real shape: the child binds in its own netns and holds): **0/100** idle, **54/100** under 8 `yes`, every one `F…FS` (a reset after 18–19 foreign answers, the last connected dial). The node mux stand-in (`mux-server.mjs`) gave 0/100 idle and **49/100** loaded, the same shape. Under r0's rule all of them are INCONCLUSIVE: **54 → 0** and **49 → 0**. **[r1]** The run was repeated with dial start times under the final rule: **49 → 0** (real srt) and **82 → 0** (stand-in), with the largest gap 94.4 ms of the 1000 ms window (V22) | V11 |
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

### 3.1 The rule [r1: tightened]

> **A zero-byte RESET records nothing if, and only if, (i) a dial that read a foreign answer
> STARTED less than `execProbeListenMS` (1000 ms) before the reset was observed, and (ii) the
> child's run then exited cleanly (`r.err == nil`).** Otherwise it is silent, as today. A reset is
> ECONNRESET at connect, on the hello write, or on the read, with zero bytes read. A zero-byte
> deadline (or any other zero-byte error) is silent wherever it falls, and the child's nonce
> dominates.

r0's rule was "any reset after any foreign answer". Quorum r1 (§11, objection 3) showed that r0
could not exclude a reset that came from the child's own listener. Conditions (i) and (ii) are
what make the proof in §3.3 go through. They cost nothing on the measured data: the largest gap
from the last foreign dial's start to an exit reset was **94.4 ms** (Linux, 131 resets) and
**57.7 ms** (darwin), against a 1000 ms bound, and every one of those runs exited cleanly (V21,
V22).

`judgeReverse` and `reverseSeen` (three bools) are **not** changed. The rule lives in the poll,
which is the only place that knows the order and timing of the dials, and in the poll's `stop`,
which `reverseLeg` calls with `r.err == nil` once the run has returned.

### 3.2 Code shape (r1 prototype measured: V23, V24)

```go
type dialSeen int // dialNone, dialChild, dialForeign, dialSilent, dialReset (as r0)

func classifyReverseDial(c net.Conn, want string) dialSeen // as r0: zero bytes + ECONNRESET (write or read) → dialReset

func pollLoopback(port int, want string, dial func(addr string) (net.Conn, error)) func(runClean bool) reverseSeen {
	…
	var exitReset bool        // poll goroutine only; read by stop after <-finished
	var lastForeign time.Time // the START of the latest dial that read a foreign answer
	…
		d := dialNone
		start := time.Now() // taken before the dial: never later than the real SYN
		if c, err := dial(addr); err == nil {
			d = classifyReverseDial(c, want)
			_ = c.Close()
		} else if dialErrIsSilent(err) { // ECONNRESET only (V19)
			d = dialReset
		}
		switch {
		case d == dialChild:
			child.Store(true)
		case d == dialForeign:
			foreign.Store(true)
			lastForeign = start
		case d == dialReset && !lastForeign.IsZero() &&
			time.Since(lastForeign) < time.Duration(execProbeListenMS)*time.Millisecond:
			exitReset = true // row 233: a foreign holder's exit, provisionally (§3.3)
		case d == dialReset || d == dialSilent:
			silent.Store(true)
		}
	…
	return func(runClean bool) reverseSeen {
		close(done)
		<-finished
		return reverseSeen{child: child.Load(), foreign: foreign.Load(),
			silent: silent.Load() || (exitReset && !runClean)}
	}
}

// reverseLeg: seen := poll(r.err == nil)
```

The prototype (`$V/proto-r1.diff.txt`) used the names `revDial*`. The sprint uses `dial*`, because
the test file already has a `revDialErr` helper (V14). The two test call sites of `pollLoopback`
(`exec_probe_test.go:829`, `:908`) pass `true`. `dialErrIsSilent` is unchanged. It matches
ECONNRESET only, and T8 pins that (V19), so no other connect error can become `dialReset`. No
refusal text changes.

### 3.3 Soundness: the accepted set does not grow on any execution that observes the child

The probe is a fail-closed startup gate. The property that must not weaken is P20: a listener the
child binds is unreachable from the host, so a real reach must refuse.

**Lemma (no observation of the child is ever excused).** Suppose that in one attempt the child
binds a host-reachable listener on `127.0.0.1:p` at time `B`, and the run exits cleanly. Then no
dial whose connection reached that listener is excused.

*Proof.*

1. **The listener lives at least `L = execProbeListenMS` = 1000 ms.** The script binds, prints
   `WORLD-PROBE-BOUND`, then runs `setTimeout(()=>{v.close();r(null)},+ms)` (`:94`, V17). Its only
   other exits are an uncaught error or a signal, and either one makes srt, and so `pp.run`, report
   a non-zero status. `r.err == nil` (condition ii) therefore means the listener was closed by its
   timer, at `C ≥ B + L`.
2. **While it listens, the child's listener never resets a connection.** Every connection it
   accepts gets `c.resume(); c.end(nonce)` (`:92`), which is the nonce and then a FIN. A dial that
   reaches it therefore reads the nonce (CHILD), or, if the nonce is late, zero bytes at the
   deadline (SILENT, which the rule never excuses). The kernel resets a connection queued on a
   listener only when that listener closes (Linux `tcp_abort_on_overflow` = 0 on the CI runner,
   V22: an overflow drops the SYN, it does not reset it). So a reset that comes from the child's
   listener is observed at `e ≥ C ≥ B + L`.
3. **A foreign answer means the child was not reachable on `127.0.0.1:p` at that moment.** Say a
   dial started at `s_f` and connected at `a_f ≥ s_f` to a listener that answered with something
   other than the nonce. Then the child's listener was not routable at `a_f`:
   - the child binds `127.0.0.1` without `reusePort` (V17), so it cannot share that address with
     another listener on darwin (V6, V7) or on Linux (V11);
   - darwin lets a wildcard listener share the port, but a `127.0.0.1` dial then always reaches
     the `127.0.0.1` listener, which is the child (20/20, V8);
   - World's client never answers anything but the nonce (`:92`; row 232 §4(a)).

   Because the child holds the port from `B` to `C`, it follows that `a_f < B` or `a_f ≥ C`.
4. **Close the cases.**
   - If `a_f ≥ C`: the foreign dial came after the child closed, and the poll is sequential, so any
     later reset dial also connects after `C` and cannot have reached the child.
   - If `a_f < B`: then `s_f ≤ a_f < B`, so `e − s_f > (B + L) − B = L`. Condition (i) fails, and
     the reset is SILENT.
   - **Measured, not true times.** The poll records `s_obs ≤ s_f`, because `start` is taken before
     `dial`. It reads `time.Since` at `e_obs ≥ e`. So `e_obs − s_obs ≥ e − s_f > L`, and a scheduler
     stall can only widen the measured gap. A stall never excuses a reset; it can only refuse.
     (The clock is Go's monotonic one.) ∎

**Theorem (the accepted set).** Take any execution of the reverse leg in which at least one
reverse-leg dial reached a host-reachable child listener. By the lemma, every such observation is
recorded as today (CHILD, or SILENT for a deadline or a reset), and the observations that are not
the child's are recorded as today or more leniently. `judgeReverse` gives CHILD priority over
everything, and SILENT priority over FOREIGN. So:

- an attempt that saw the nonce is REACHED, as today;
- an attempt that saw the child only through a deadline or a reset is SILENT, as today.

The rule cannot make that attempt INCONCLUSIVE, so it cannot lead to a later attempt's PASS.
**The leg's verdict on every execution that observes the child is unchanged.**

**Exactly what does change.** Executions in which no dial reached a child listener, and in which
every silent observation of an attempt is a reset satisfying (i) and (ii). Today such an attempt is
SILENT. Under the rule it is INCONCLUSIVE, so the leg retries on a fresh port, and three
INCONCLUSIVE attempts still refuse (`:621`, `:111`). By the lemma, the excused reset is **not** the
child's. Today's refusal of that attempt rests only on a non-child listener's exit, which is the
class row 229 deliberately retries (foreign).

**gpt's case, answered.** "A reachable child observed only through resets after a foreign answer"
cannot have those resets excused (lemma, case `a_f < B`). Such an execution is SILENT under the
rule, exactly as today. The rule does change executions with a reachable but **never-observed**
child, where today's refusal came from a foreign holder's exit by chance. In those, the child was
not observed under either rule. Today's refusal hangs on whether the holder's exit happened to cut
a dial: 49–54 of 100 collided attempts on loaded Linux, 0–2 of 100 elsewhere (V10, V11, V21, V22).
The rule does not change what the probe observes of the child. It changes only how a non-child
listener's exit is attributed.

**What the proof relies on (stated, with evidence):**

- (a) the poll is one goroutine and dials sequentially (`:745-765`, V3);
- (b) the child cannot co-hold `127.0.0.1:p` with another listener, and darwin routes to the more
  specific address (V6–V8, V11);
- (c) the script's listener closes only through its timer, unless the process fails (`:92-94`,
  V17), and a failure reaches `r.err` (srt's `process.exit(code ?? 0)` and `exit(1)` for a signal,
  `cli.js:428-447`, V5);
- (d) a listening socket does not reset queued connections before it closes. Linux
  `tcp_abort_on_overflow` = 0 was measured on the runner (V22). On darwin this is the BSD
  behaviour and was **not separately measured**.

If (b) or (d) failed, the cost would be a co-held child answering the nonce (REACHED), or an early
reset counted as silent. Both fail closed, except a (d) failure that resets within 1000 ms of a
foreign answer. That one is listed as R6.

| Attempt shape (poll order; `clean` = the run exited 0) | Today | After | Why |
|---|---|---|---|
| any dial reads the nonce | REACHED | REACHED | the child is first in `judgeReverse` (`:152`) |
| `tok` is not a probe token | NO-TOKEN | NO-TOKEN | checked before silent (`:154`) |
| reset(s) only, no foreign (`S`, T7 `reset-holder`, T9 `connect-reset`) | SILENT | SILENT | no foreign anchor |
| reset **before** the first foreign answer (`SF`, `CF`) | SILENT | SILENT | no foreign anchor yet |
| a zero-byte **deadline** anywhere (`FD`, `FSD`) | SILENT | SILENT | the rule never excuses a deadline |
| foreign, then a reset ≥ 1000 ms after the foreign dial **started** (`FP`, `LQ`) | SILENT | SILENT | condition (i) |
| foreign, then a reset < 1000 ms after, but the run did **not** exit cleanly (`FS`, unclean) | SILENT | SILENT | condition (ii) |
| foreign, then reset(s) < 1000 ms after the latest foreign dial started, clean exit (`FS`, `FC`, `FW`, `fS`, `GGFS`) | SILENT | **INCONCLUSIVE** | the changed row; by the lemma the reset is not the child's |
| foreign, then a non-reset connect error (`FA` ECONNABORTED, `FT` timeout) | INCONCLUSIVE | INCONCLUSIVE | `dialErrIsSilent` is ECONNRESET-only (V19): records nothing, as today |
| connect error only (`A`, `T`) | PASS | PASS | as today |
| foreign only, or `EADDRINUSE` with nothing | INCONCLUSIVE | INCONCLUSIVE | — |

### 3.4 Alternatives (r1: re-evaluated)

- **(A) Order-free `{foreign, silent}` → INCONCLUSIVE.** Rejected. It excuses `SF` and resets
  of any age. That is M31, killed by T12 `SF`, `CF`, `FP` and `LQ` (V24).
- **(B) Any silent after foreign, deadlines included.** Rejected. A deadline is a live, silent
  peer, which is row 229 R1's slow child. That is M32, killed by `FD` and `FSD`.
- **(r0) Any reset after any foreign answer.** Superseded. It could not exclude a child reset that
  follows a foreign answer from before `B` (quorum r1, objection 3). M37, which drops the window,
  is r0's rule, and it is killed by `FP` and `LQ`.
- **(Controller a) Excuse only when the child's token is `WORLD-PROBE-REFUSED <errno>`.** Rejected.
  On Linux the correct child binds in its own netns and reports `WORLD-PROBE-BOUND` (row 229 §4.3
  table, row "Linux, correct, port collided with srt's mux": BOUND → foreign). So (a) would excuse
  nothing on Linux, where the false refusals are: 49/100 and 54/100 loaded (V11, V22). The timed
  window gives the same guarantee for every token.
- **(Controller b) Adjacency plus a bound from the constants.** Rejected as stated. The longest
  gap between dial starts is `200 + 500 + 50 = 750 ms` (`execProbeHelloWait`, the dial timeout,
  `execProbePollEvery`). That is shorter than the 1000 ms hold, but it does not exclude a child
  reset adjacent to a foreign answer. For example, a foreign answer that arrives 450 ms after its
  connect, then the next dial connecting 200 ms later, gives `B` just after `a_f` and the child's
  close-time reset ≈ `a_f + 1000`, inside a 500 ms read. Measuring the gap directly (condition i)
  closes that case exactly, and every observed exit reset was adjacent anyway (V21, V22).
- **(C) "The reset must be the attempt's last connected dial."** Not needed for soundness, given
  the lemma. Recorded only.
- **(D) Count EPIPE as a reset.** Not adopted. It was never observed. The one write failure seen
  was ECONNRESET (`W`, V11). EPIPE therefore stays silent, which fails closed.

## 4. Test plan

| Test | Kind | Needs node? | Kills | Measured? |
|---|---|---|---|---|
| **T12** `TestExecProbeReversePollHolderExitAfterForeign` (new; 18 rows) | scripted dial seam. No socket. Sleeps are used only as lower bounds (rows `FP`, `LQ`, `GGFS`) | no | M30–M40 (M33 also by T7/T9, M36 also by T8) | **yes**, r1 prototype: control 52 PASS lines; every mutant red (V24) |
| T8 `TestExecProbeReverseDialErrIsSilent` (unchanged) | pure | no | M36 (`op-aborted`) | V19, V24 |
| T1, T7, T9 (T9 and T10 pass `true` to `stop`) | as today | as today | M33 (T7 `reset-holder`, T9 `connect-reset`) | V24 |
| T2–T6, T10, T11 (unchanged) | as today | as today | regression | green on the r1 prototype (V23) |

**T12's scripted alphabet** (`$V/zz_iter330_t12r1_test.go.txt`):

- `N` nonce; `F` HTTP 400; `f` clean EOF;
- `S` read ECONNRESET; `W` write ECONNRESET then EOF; `D` `os.ErrDeadlineExceeded`;
- `C` connect ECONNRESET; `A` connect ECONNABORTED; `T` a connect timeout (`os.ErrDeadlineExceeded`);
- `P` a read reset whose dial blocks ≥ `execProbeListenMS`;
- `L`/`G` a foreign answer whose dial blocks ≥ 600 ms; `Q` a read reset whose dial blocks ≥ 450 ms.

The dialer closes `all` once every step has been handed out. The test then calls
`stop(clean)` and asserts the exact `reverseSeen` and `judgeReverse("WORLD-PROBE-BOUND", seen)`.

**Why the sleeps cannot flake a correct tree.** Each sleep is a lower bound on elapsed time:

- `FP`: the gap is ≥ 1000 ms, so it must be SILENT for any scheduling.
- `LQ`: the gap is ≥ 600 + 450 ms from L's start, so SILENT.
- `GGFS`: the correct anchor is F's start, and its gap is only F + 50 ms + S, which is ≈ 50 ms.
  That leaves about 950 ms of scheduling margin before a correct tree could redden. The first-G
  anchor (M40) is ≥ 1250 ms.

| Rows | `clean` | Want `reverseSeen` | Want verdict |
|---|---|---|---|
| `FS`, `FC`, `FW`, `fS`, `GGFS`, `FA`, `FT` | true | `{foreign}` | INCONCLUSIVE |
| `FD`, `SF`, `CF`, `FSD`, `FP`, `LQ` | true | `{foreign, silent}` | SILENT |
| `FS-unclean` | **false** | `{foreign, silent}` | SILENT |
| `S` | true | `{silent}` | SILENT |
| `FSN` | true | `{child, foreign}` | REACHED |
| `A`, `T` | true | `{}` | PASS |

**AC-MUT.** Apply M30–M40 one at a time. Each reds its named killer, and the tree is restored
byte-identical. **AC-CI**: T12 joins step 8's `tests=` list (`ci.yml:247`, V16). **AC-STABLE**:
`go test ./host/broker -run '^TestExecProbeReverse' -count=20` passes 20/20, and T12 `-count=200`
under 32 `yes` gives 0 FAIL. **AC-GATE**: `gofmt -l` is empty, `go vet ./host/broker/` is clean,
`go build ./... && go test ./...` passes, and `./scripts/verify_ail.sh` passes. **AC-LINUX**: the fix PR's linux
step 8 shows `--- PASS: TestExecProbeReversePollHolderExitAfterForeign`.

## 5. Mutation table (continues the parent's M1–M21, row 232's M22–M29, the judge's O1–O10, row 230's K1–K9)

All were run against the r1 prototype (V24). The locals the rule uses are kept referenced so
that each mutant compiles.

| ID | Mutant | Killer(s) |
|---|---|---|
| M30 MUT-EXIT-RESET-SILENT | the excuse removed (**today's behaviour**) | T12 `FS`, `FC`, `FW`, `fS`, `GGFS`, `FSN` |
| M31 MUT-EXIT-RULE-ORDER-FREE | any reset is excused if the attempt **ever** saw foreign (option A) | T12 `SF`, `CF`, `FP`, `LQ` |
| M32 MUT-DEADLINE-EXCUSED | a deadline inside the window is excused too (option B) | T12 `FD`, `FSD` |
| M33 MUT-EVERY-RESET-EXCUSED | a reset is excused with no foreign anchor | T12 `S`, `SF`, `CF`, `FP`, `LQ`; T7 `reset-holder`; T9 `connect-reset` |
| M34 MUT-CONNECT-RESET-NOT-RESET | a connect ECONNRESET maps to `dialSilent` | T12 `FC` |
| M35 MUT-WRITE-RESET-NOT-RESET | only the read's ECONNRESET is a reset | T12 `FW` |
| **M36 MUT-DIALERR-WIDENED** [r1, quorum objections 1–2] | `dialErrIsSilent` also matches ECONNABORTED | T8 `op-aborted` (`dialErrIsSilent(… software caused connection abort) = true, want false`); T12 `A` (`seen {silent:true} verdict 3; want {} verdict 0`) |
| **M37 MUT-WINDOW-DROPPED** [r1] (= r0's rule) | condition (i) removed | T12 `FP`, `LQ` |
| **M38 MUT-ANCHOR-AT-FOREIGN-END** [r1] | the window is measured from the foreign dial's END, not its start | T12 `LQ` |
| **M39 MUT-CLEAN-GATE-DROPPED** [r1] | condition (ii) removed (`stop` ignores `runClean`) | T12 `FS-unclean` |
| **M40 MUT-ANCHOR-AT-FIRST-FOREIGN** [r1] (sound, but loses availability) | the window is measured from the FIRST foreign dial | T12 `GGFS` |
| M41 (wiring, review-only) | `reverseLeg` passes `true` instead of `r.err == nil` | **no deterministic killer.** It is one line, checked by review. A reverse-leg fixture would need a real exit reset timed against the run, which is wall-clock. Recorded as the single unpinned line |

## 6. Milestones

- **M1** (one milestone): `dialSeen`, the classifier's return, the poll's timed case and `stop(runClean)`,
  `reverseLeg`'s `poll(r.err == nil)`, T12, the two test call sites, the `ci.yml:247` gain, and the
  QUICKSTART clause. Prove M30–M40 red one at a time; review M41's one line.

## 7. Files

| File | Change | LOC (est.) |
|---|---|---|
| `host/broker/exec_probe.go` | `dialSeen` (near `:686`); `classifyReverseDial` returns it (`:687-720`, doc comment updated); `pollLoopback`'s locals, loop body (`:737-759`), `stop(runClean bool)` (`:767-771`) and doc comment (`:727-732`); `reverseLeg` `:600` `seen := poll(r.err == nil)` | +35 / −15 |
| `host/broker/exec_probe_test.go` | `scriptConn`, `scriptedDial`, T12; T9 `:829` and T10 `:908` call `stop(true)` | +150 / −2 |
| `.github/workflows/ci.yml` | `:247` `tests=` gains `TestExecProbeReversePollHolderExitAfterForeign` | ±1 |
| `docs/QUICKSTART.md` | §10 arm 6 (`:1168-1170`): after "(or that resets the host's dial)", add "; a reset within a second of a foreign answer on the same port, in a run that exited cleanly, is that process exiting, and is retried" | ±1 |

**S7** applies: the refusal-text paragraph describes which resets refuse. **S2/S3**: this is a
repair inside the host's effect boundary, with no new surface. It is not a package, for the same
reason as rows 229, 232 and 230. **S1/S4/S5** do not apply: there is no `.ail`. No dependency is
added.

## 8. Conflict surface, risks and residuals

**Conflict surface.** Three hunks of `exec_probe.go` (`:686-720`, `:737-771`, `:600`) and their
comments. Arms 1–5, 7 and 8, `judgeReverse`, `reverseSeen`, every constant and every refusal text
are untouched. `pollLoopback` has one production call site (`:598`, stopped at `:600`) and two test
call sites (`:829`, `:908`; V14). Its returned `stop` gains the `runClean` parameter.

- **R1, the rate after the fix is 0 for the observed signature; the collision itself is untouched.**
  In every hold-1000 run, darwin and Linux, real srt and stand-in alike, each exit silent was a
  reset after foreign, within 94.4 ms of the foreign dial's start, in a clean run: darwin load 2 → 0
  and 1 → 0, Linux load 54 → 0 and 49 → 0, Linux stand-in 49 → 0 and 82 → 0 (V10, V11, V21, V22). So
  the rule removes the false refusal it targets. It does not make the collision rarer (row 229 R4 watch), and a foreign holder
  that is silent **before** it answers still refuses.
- **R2, what the change does accept [r1].** It accepts no execution that observes the child (§3.3
  lemma and theorem). It does retry, rather than refuse, an attempt whose only silent evidence is a
  non-child listener's exit reset. If the child is reachable but never observed in that attempt,
  today's refusal came from that coincidence. This is stated, not hidden (§3.3, "gpt's case").
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
- **R6, darwin's "no reset before close" (§3.3 reliance d) is not separately measured [r1].** Linux
  `tcp_abort_on_overflow` = 0 was measured (V22). On darwin it rests on the BSD behaviour. A
  violation would matter only if a live child listener reset a host dial within 1000 ms of a
  foreign answer, in a run that exited cleanly. That cannot happen in a correct sandbox, where the
  host cannot reach the listener at all.
- **R7, one transient red during the r1 prototype's suite run [r1].**
  `TestExecProbePassThroughShimFailsEveryConfinementArm` failed once on its expected text: the proxy
  leg showed `accepts 4`, where the test wants 3 (`$V/out-proto-r1-suite.txt`). That is arm 6's
  World-owned listener, a path this change does not touch. It did not reproduce: 5/5 PASS on the
  prototype and 5/5 at base (V23). Recorded as the row 230 R2 cross-process accept class, not
  attributed to this rule.
- **R4, EPIPE and other write errors stay silent** (option D, not adopted). If CI ever shows an
  exit with that shape, the R4 trigger of row 230 §3.5 (`answered nothing to World's hello` in a
  non-test log) fires, and the class is re-measured first.
- **Carried forward unchanged:** row 232 R6 (darwin EINVAL) and row 230 R1 (the cross-profile
  ordering). That successor's hoist does not interact with this rule.

## 9. What this item is NOT doing

- No change to `judgeReverse`, `reverseSeen`, the attempt count, the hold, the cadence, the hello
  wait or any refusal text. `execProbeListenMS` is now also the excuse window, and the lemma
  depends on that: whoever changes the hold changes the window with it.
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
| V12 | **[r0] The r0 prototype, and the whole ExecProbe suite on it** | prototype applied (`$V/proto-exec_probe.diff.txt`); `$V/zz_iter330_t12_test.go.txt` copied in; `go vet ./host/broker/`; `go test ./host/broker/ -run '^(TestExecProbe\|TestZZIter330)' -count=1 -v` → `$V/out-proto-suite.txt`; then `git checkout -- host/broker/exec_probe.go`, test deleted, `shasum -a 256 host/broker/exec_probe.go`, `git status --porcelain` | vet clean; `22 --- PASS`, `0 --- FAIL`, `1 --- SKIP` (`TestExecProbeArm8`, real srt needs CI's env); T1–T11 and the T12 prototype all PASS; `ok … 36.160s`. Restored `96f33172667d27de`; porcelain shows only `?? design_docs/verification/world-iter330/` |
| V13 | **[r0] T12 r0 prototype against M30–M35** (superseded by V24) | `bash $V/run-mutants.sh.txt` → `$V/out-mutants.txt` (the step-8 reverse-leg tests T1, T2, T7, T8, T9, T10 plus T12) | `CONTROL (prototype): 44 PASS lines`. M30: `FC FS FSN FW fS` FAIL, `dials FS: seen {… foreign:true silent:true} verdict 3; want {… silent:false} verdict 4`. M31: `CF SF` FAIL, `dials SF: … verdict 4; want … silent:true verdict 3`. M32: `FD FSD` FAIL. M33: `…SilentHolderRefusesSticky/reset-holder`, `…PollRecordsConnectReset/connect-reset`, T12 `CF S SF` FAIL. M34: `FC` FAIL. M35: `FW` FAIL. `restored-to-prototype a3568e76a6709599` |
| V14 | Census | `grep -rn classifyReverseDial --include='*.go' host cmd \| grep -v host/broker/exec_probe.go \| wc -l`; `grep -rn 'pollLoopback(' --include='*.go' host cmd`; `grep -rn revDialErr --include='*.go' host cmd \| wc -l` | `0`; `exec_probe.go:598` (`reverseLeg`), `:733` (def), `exec_probe_test.go:829` (T9), `:908` (T10); `12` (the test helper, so the sprint avoids the `revDial*` prefix) |
| V15 | QUICKSTART's arm-6 text | `sed -n 1164,1172p docs/QUICKSTART.md` | `:1168-1170` "`reverse leg: a listener on … answered nothing to World's hello` means a process that never answers (or that resets the host's dial) held the port, and World cannot tell it from the child: find it." |
| V16 | CI step 8's list | `grep -n -e 'TestExecProbe' .github/workflows/ci.yml` | `:247` lists T1–T11 (`…TestExecProbeReverseHelloWaitBounds TestExecProbeChildListenerDrainsTheHello`) |
| V17 | World's child binds positionally (no `reusePort`) | `grep -n "v.listen(+port,'127.0.0.1'" host/broker/exec_probe.go` | `:94 v.listen(+port,'127.0.0.1',()=>{say('WORLD-PROBE-BOUND');…` |
| V18 | The stand-in rate this row inherited | read `design_docs/verification/world-iter318/designer-probe-outputs.txt` via row 232 V5/V6 | mux-shaped stand-in, 1000 ms, faithful: idle `pre=3 post=3`; load `silent-postconnect 13`, `pre=13 post=13` |
| V19 | **[r1, objections 1–2] `dialErrIsSilent` matches ECONNRESET only, and T8 pins it.** VERIFIED BY CONTROLLER at `a6af0cd4`, then re-run first-party | `sed -n 722,725p host/broker/exec_probe.go`; `grep -n 'dialErrIsSilent(' host/broker/exec_probe.go`; `go test ./host/broker/ -run '^TestExecProbeReverseDialErrIsSilent$' -count=1 -v` → `$V/out-t8-base.txt` | `:725 func dialErrIsSilent(err error) bool { return errors.Is(err, syscall.ECONNRESET) }`; callers `:725` (def) and `:757` only. T8: `op-reset`, `bare-reset`, `wrapped-reset` (want true) and `nil`, `op-refused`, `real-refused`, `op-einval`, `op-aborted`, `addrnotavail`, `timeout`, `deadline-exceeded` (want false): 11 `--- PASS`, `ok … 0.350s` |
| V20 | [r1] The constants and script lines the lemma and option (b) use | `grep -n '^const execProbeListenMS\|^const execProbeHelloWait\|^const execProbePollEvery\|DialTimeout("tcp", addr, 200\|net.createServer(c=>{c.resume();c.end' host/broker/exec_probe.go` | `:92 … net.createServer(c=>{c.resume();c.end('WORLD-PROBE-NONCE '…`; `:108 execProbeListenMS = 1000`; `:128 execProbeHelloWait = 500 * time.Millisecond`; `:131 execProbePollEvery = 50 * time.Millisecond`; `:735 … net.DialTimeout("tcp", addr, 200*time.Millisecond)`. With V17 (`:94` the timer close) |
| V21 | **[r1] darwin, real srt, hold 1000 ms, 32 `yes`, with dial start times** (`$V/srtexit-r1.go.txt`: every observation carries the `time.Now()` taken before its dial; each silent prints `sinceLastForeignDialStart` and `cleanExit`; rule E = §3.1) | `$V/run-darwin-r1.sh.txt` → `$V/out-darwin-r1.txt` | `load: { 4.35 5.96 9.01 }` → `{ 50.25 27.44 17.50 }`. `CLASSES map[F:2152 R:11 S:1 f:1]`; `cycle 79 silent S at exit-3.8ms sinceLastForeignDialStart=57.7ms cleanExit=true`; `STICKY-SILENT-CYCLES today=1 … ruleD=0 ruleE(timed+clean,r1)=0`; `yes-alive-after=0` |
| V22 | **[r1] Linux, the same with start times, plus `tcp_abort_on_overflow`.** Throwaway branch `probe/world-iter330-holderexit-r1` (`8b2d0fc`; only the workflow and probe files), run 41186735782, job 124495072930. Log fetched with `gh api --allow-escape-sequences …/jobs/124495072930/logs` (61,711 bytes), ANSI stripped → `$V/linux-probe-r1-run41186735782.txt`. Then `git push origin --delete probe/world-iter330-holderexit-r1`; `git ls-remote origin 'refs/heads/probe/world-iter330-*' \| wc -l` | `abort_on_overflow=0`. IDLE srt: `F:1999`, `today=0 … ruleE=0`. Load `0.96 …` → `11.48 …`. **LOAD srt: `CLASSES map[F:1870 R:9 S:49 f:1]`, `silentAfterForeign=49 connectedAfterFirstSilent=0`, `today=49 … ruleE(timed+clean,r1)=0`.** LOAD stand-in: `F:1900 R:7 S:81 W:1`, `today=82 … ruleE=0`. Over all 131 silent lines (130 `S`, 1 `W`): `max sinceLastForeignDialStart 94.4ms`, `cleanExit=true` 131/131. `yes-alive-after=0`, `PROBE-DONE`. Remote probe branches after deletion: `0` (control `refs/heads/dev`: `1`) |
| V23 | **[r1] The r1 prototype and the suite on it** | `$V/proto-r1.diff.txt` applied; `go vet ./host/broker/`; `$V/zz_iter330_t12r1_test.go.txt` copied in; `go test ./host/broker/ -run '^(TestExecProbe\|TestZZIter330)' -count=1 -v` → `$V/out-proto-r1-suite.txt`; then `TestExecProbePassThroughShimFailsEveryConfinementArm` `-count=5` at base and on the prototype; restore with `git checkout`, delete the scratch test, `shasum -a 256 host/broker/exec_probe.go` | vet clean. Suite: 21 `--- PASS`, 1 `--- FAIL`, 1 `--- SKIP` (`TestExecProbeArm8`). The FAIL was `TestExecProbePassThroughShimFailsEveryConfinementArm`, `… does not report "a connect through srt's proxy environment: token \"WORLD-PROBE-CONNECTED\", accepts 3"` (it saw `accepts 4`; R7). Re-runs: base `5 PASS`, prototype `5 PASS, 0 FAIL`. T12's 18 rows PASS (`4.24s`). Restored `96f33172667d27de` |
| V24 | **[r1] T12 r1 against M30–M40** | `python3 -I $V/run-mutants-r1.py.txt` → `$V/out-mutants-r1.txt` (T1, T2, T7, T8, T9, T10 and T12; the rule's locals kept referenced so each mutant compiles) | `CONTROL (r1 prototype): 52 PASS lines`. M30: T12 `FC FS FSN FW GGFS fS`. M31: `CF FP LQ SF`. M32: `FD FSD`. M33: T7 `reset-holder`, T9 `connect-reset`, T12 `CF FP LQ S SF`. M34: `FC`. M35: `FW`. **M36: T8 `op-aborted` (`dialErrIsSilent(dial tcp: connect: software caused connection abort) = true, want false`) and T12 `A` (`seen {… silent:true} verdict 3; want {…} verdict 0`).** M37: `FP LQ`. M38: `LQ`. M39: `FS-unclean`. M40: `GGFS`. `restored-to-r1-prototype 9cdd698ea2ff3284` |

## 11. Quorum log

**Round 1** (artifact `~/.ailang/state/mission-quorum-world/w-exec-probe-holder-exit-reset-2026-10-31T08-49-28Z.json`,
$0.3150): **BLOCKED**, 3/3 reject. Answered in r1:

| # | Objection (seat) | r1 answer |
|---|---|---|
| 1 | `oc-glm-5-3`: the rule maps all of `dialErrIsSilent`'s matches to `dialReset`, but no V-row shows that function's body | V19 shows the body (`errors.Is(err, syscall.ECONNRESET)`, `:725`) and its only caller (`:757`), and re-runs T8's 11 rows first-party (the controller had measured them first). §3.2 states that the match set is ECONNRESET only |
| 2 | `oc-kimi-k3`: no T12 row or mutant tries a NON-ECONNRESET connect error after a foreign answer, so a too-broad helper would pass | T12 gains `FA` (ECONNABORTED after foreign → `{foreign}`, as today), `A` (ECONNABORTED alone → `{}` PASS), `FT` and `T` (a connect timeout). **M36** widens `dialErrIsSilent` to ECONNABORTED and is red on T8 `op-aborted` and T12 `A` (V24). `FA` alone cannot kill M36, because an excused and a recorded-nothing dial look the same after a foreign answer. `A` is the row that discriminates |
| 3 | `gpt6-1-sol`: a foreign answer names neither the owner of a later connection nor the cause of its reset; r0 admitted a child seen only through resets after a foreign answer; "a different execution already escapes" is no proof | **The rule is tightened (§3.1), and the argument is replaced by a proof (§3.3).** A reset is excused only if a foreign dial **started** less than `execProbeListenMS` before the reset was observed, **and** the run exited cleanly. The lemma: the child's listener lives ≥ 1000 ms (its timer; a clean exit rules out an early death), it resets nothing before it closes, and no foreign answer can arrive while it holds `127.0.0.1:p`. So any child reset comes > 1000 ms after the start of every earlier foreign dial, and is never excused. Host timing only widens the measured gap. The theorem: every execution that observes the child keeps today's verdict. What does change is stated exactly: attempts whose only silent evidence is a non-child exit reset. **Availability is unchanged from r0 on the banked data:** darwin load 1 → 0, Linux load 49 → 0, Linux stand-in 82 → 0, largest gap 94.4 ms, every run clean (V21, V22). The controller's candidate (a), token-gating, is rejected because the Linux token is BOUND (row 229 §4.3), which would buy 0 of the 49. Candidate (b), adjacency plus a constant bound, is rejected with a counter-example (750 ms < 1000 ms does not exclude it) in favour of measuring the gap directly (§3.4). New mutants M37 (window dropped, = r0), M38 (end anchor), M39 (clean gate dropped) and M40 (first-foreign anchor) are all red (V24). M41 (`reverseLeg` passes `true`) is the one review-only line |
