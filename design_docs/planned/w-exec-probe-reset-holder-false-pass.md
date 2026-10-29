# w-exec-probe-reset-holder-false-pass — a dial that fails AT CONNECT with ECONNRESET is dropped, so a resetting holder can pass the reverse leg

- **Status**: Planned (iteration 318, designer, **r2**). Queue row 232. An amendment to
  `design_docs/implemented/w-exec-probe-reverse-leg-port-race.md` (the parent; §4.2, §6, §7).
- **Item**: Clause 3 (the Workspace.Exec fail-closed startup probe), a conformance defect that
  reddened `dev` itself. Run 40833618366, job 123313206642 ("go race gate", linux/amd64) at
  `df5a0671`: `--- FAIL: TestExecProbeReverseLegSilentHolderRefusesSticky/reset-holder`,
  `exec_probe_test.go:664: Why = "", picker calls 1; want the sticky silent refusal on attempt 1,
  calls 1` (V1). This is instance 3. PR #402 did not touch `host/broker/exec_probe*.go`.
- **Estimated**: ~0.25 day. About +12/−3 LOC in `host/broker/exec_probe.go`, about +75 LOC of
  tests, a two-name CI list edit and one QUICKSTART clause (§7).
- **Measurement base**: `df5a0671` (worktree `sprint/world-iter318-row232`, clean, V2).
- **Instruments**: `go1.26.6 darwin/arm64` (16 cores), node v26. Scratch probes built with local
  go1.26.4 (`net/fd_unix.go` byte-identical, V9), run from `/tmp/w318-designer/`; [r2] sources,
  runners and raw outputs are checked in as `design_docs/verification/world-iter318/designer-*`.
  Load = 32 `yes` on 16 cores, own PIDs killed, 0 alive after (V6, V7). **[r2] Linux was
  measured by the controller** on a CI runner (`7.0.0-1017-azure x86_64`, nproc 4, go1.26.6; V14).

---

## 1. The defect

`pollLoopback` (`exec_probe.go:694`) classifies a dial only when it connects:
`if c, err := net.DialTimeout("tcp", addr, 200*time.Millisecond); err == nil { … }` (`:705`).
A dial **error** records nothing: not child, not foreign, not silent (V3). The parent design says
"zero bytes at the deadline, **or a reset with zero bytes**: set `seen.silent`" (§4.2). The code
honours that only when the reset arrives after `connect` returned (`classifyReverseDial`, `:660`).

Go's non-blocking connect reads `SO_ERROR` once the socket is writable. If the peer accepted and
reset the connection before that read, `connect` itself returns
`os.NewSyscallError("connect", ECONNRESET)` (`fd_unix.go:127-131`, `default:` arm; V9). The
reset and the connect error are the same TCP event, and only Go's read timing separates them.

T7's `reset-holder` (`SetLinger(0)` then `Close`, `exec_probe_test.go:494-496`) runs behind the
`/bin/sh` fixture client, which exits at once. So the attempt usually makes exactly **one** dial.
When that dial fails at connect, `judgeReverse` sees `reverseSeen{}` and returns PASS (row 5).
Re-measured this session (V4): against an accept-then-RST listener, 43 of 2000 sequential dials
idle and 79 of 2000 under load failed at connect with `connection reset by peer`. Under 4 spinning
dialers for 45 s, 96,566 of 588,731 did, with no other error class (V7). [r2] **On Linux CI**
(V14), the same holder shape with one sequential dialer (like T7's lone dial) gave ECONNRESET at
connect on 1,742 of 644,611 dials idle and on **147,968 of 343,060 under load**. Also under load,
T7 `reset-holder` at base reproduced the CI signature at least 7 times in `-count=300`.

**Not a breach (unchanged from iteration 317, re-argued in §4).** It is a gap against the
parent's classification, and it can make a resetting holder invisible on a short attempt.

## 2. Classification of dial errors (the decision)

The rule: a failed dial counts as an observation only if its error **proves the handshake
completed with a listener on the port**. Exactly one class does that: ECONNRESET.

| `net.DialTimeout` result on loopback | Class | Evidence |
|---|---|---|
| `nil` (connected) | unchanged: `classifyReverseDial` gives child, foreign or silent | parent §4.2 |
| `ECONNRESET`, matched with `errors.Is(err, syscall.ECONNRESET)` through `*net.OpError` → `*os.SyscallError` → `Errno` | **SILENT** (sticky refusal) | Never seen with nothing listening: 0 of 2,013,696 dials under load on darwin (V7), and **0 of 7,325,743 on Linux**, idle and loaded (V14). Seen with a resetting listener on both platforms (V4, V7; Linux 1,742 / 342 / 147,968 / 1,184 across the four rst runs, V14) and with a listener that closes while dials are in flight (V5, V6). The same event as the post-connect zero-byte reset the parent already classes silent |
| `ECONNREFUSED` | none | Nothing listens. That covers every dial before the child binds and after it closes, a macOS `EPERM` child, and a Linux netns child. Counts: 2,013,574 of 2,013,696 no-listener dials on darwin (V7) and 7,325,743 of 7,325,743 on Linux (V14). Against the resetting holder it occurred **0 times on Linux** (V14), so the holder does not surface as a class this fix drops. Classing it silent would refuse every attempt (M25) |
| `EINVAL` (`connect: invalid argument`, darwin only) | **none** | **Seen with NOTHING listening:** 122 of 2,013,696 under load (V7). So it is not evidence of a listener. Classing it silent would add false refusals on a loaded Mac (M27). It is also seen around a listener's close under load (V6), and so stays a known unclassified residual (R6). Linux showed 0 (V14) |
| timeout (`net.Error.Timeout()`, `os.ErrDeadlineExceeded`) | none | It is ambiguous: a full backlog, or a scheduler stall. Darwin: 0 in about 2.6 M loopback dials under load (V7). [r2] **Linux: timeouts DO occur against the resetting holder under 4 spinning dialers**, 143 idle and 179 loaded, and 0 with one sequential dialer (V14). They stay "none", because classing them silent is not proven safe (with nothing listening, no timeout was seen, but that is not a proof). This is the residual reason T7 is not strictly deterministic (§6). The child's listener is dialled about 20 times per attempt (V15), so one lost dial cannot hide a reach |
| `ECONNABORTED` | none | Not observed in any run (V5–V7). Pinned by a T8 row so that a change has to be deliberate |
| local-side errors (`EADDRNOTAVAIL`, `EMFILE`, `EACCES`/`EPERM` from a host firewall) | none | They describe the dialling side, not a peer |

**darwin and linux.** `syscall.ECONNRESET` is per-GOOS (`Errno(0x36)` on darwin/arm64,
`Errno(0x68)` on linux/amd64, V9), so `errors.Is` matches the platform's own value. The matcher
was run on darwin over synthesised and real errors, and a scratch program using it
cross-compiles for linux/amd64 (V10). [r2] **Measured on Linux CI (V14), not read from
memory.** On Linux the resetting holder's connect-time error is ECONNRESET (`0x68`,
`errors.Is` true). ECONNRESET never appears with nothing listening, and ECONNREFUSED never
appears against the holder. So on Linux, too, a connect-time ECONNRESET means a listener
completed the handshake, and the fix catches the T7 dial that today's code drops.
The assertion text in the CI failure (`Why = "", picker calls 1`, V1) does **not** identify the
errno. Any dropped dial error, or no dial at all, prints the same text. The errno is identified
by the Linux probe (V14), and the fixed head's own Linux CI must confirm it (AC-LINUX, §6).

### 2.1 Code shape

```go
// dialErrIsSilent reports whether a failed reverse-leg dial proves a listener accepted the
// handshake and reset it (row 232): only ECONNRESET. ECONNREFUSED (nothing listens), EINVAL
// (seen on darwin with nothing listening), timeouts and local errors record nothing.
func dialErrIsSilent(err error) bool { return errors.Is(err, syscall.ECONNRESET) }
```

`pollLoopback(port int, want string, dial func(addr string) (net.Conn, error))` gains the seam
parameter. `nil` means `net.DialTimeout("tcp", addr, 200*time.Millisecond)`. Its `err != nil`
branch becomes `if dialErrIsSilent(err) { silent.Store(true) }`. `reverseLeg` passes `nil`
(`:566`). Cadence, timeouts, `classifyReverseDial`, `judgeReverse` and every refusal text are
unchanged. The SILENT text ("answered nothing to World's hello") still holds, because a reset
peer answered nothing.

## 3. False-refusal risk, measured

**Ordering (read).** The poll is **not** stopped before the child closes its listener.
`reverseLeg` starts the poll (`:566`), blocks in `pp.run` until the client process exits
(`:567`), and only then stops the poll (`:568`). The client's `listen` op calls `v.close()` at
`+ms` inside the process (`:84`). So the poll dials straight across the close, and the close is
measurable.

**Probe** (`design_docs/verification/world-iter318/designer-close-probe.go.txt`, sha256 prefix
`0b5f821ff9e8`, run from `/tmp/w318-designer/probe/`). It frees a port
the way `freeLoopbackPort` does, then starts `pollLoopback`'s exact dial, hello and classify loop
before it spawns the process (50 ms cadence, 200 ms dial, 500 ms hello wait). The loop stops only
after the process exits. `-mode child` runs `execProbeScript`'s `listen` body verbatim
(`:82-84`). `-mode mux` is an srt-mux-shaped chatty holder (HTTP 400 on any bytes, parent V23)
whose process exits after the hold. `-tight 4` adds 4 spinning dialers to amplify the close.
Columns are cycles with ≥1 silent observation under today's rule and under the fix (V5, V6).

| Shape | Run | Cycles | Child nonce seen | Silent today | Silent with the fix |
|---|---|---|---|---|---|
| World's child, 1000 ms hold, faithful | idle | 100 | 100 | 2 | 2 |
| World's child, 1000 ms hold, faithful | load | 100 | 100 | 0 | 0 |
| World's child, 100 ms, tight 4 | idle / load | 200 / 200 | 200 / 200 | 167 / 179 | 185 / 191 |
| mux-shaped holder, 1000 ms, faithful | idle | 100 | 0 | **3** | **3** |
| mux-shaped holder, 1000 ms, faithful | load | 100 | 0 | **13** | **13** |
| mux-shaped holder, 100 ms, tight 4 | idle / load | 200 / 200 | 0 | 198 / 194 | 198 / 195 |

**Readings.**
1. **World's own child listener does reset host dials as it closes** (at and after connect), so
   iteration 317's "the child's listener does not reset" is false at close-time. Harmless: in a
   correct sandbox it is unreachable from the host (Linux netns, macOS `EPERM`; parent §4.3). When
   reachable (a breach, the pass-through control) the nonce was seen in **every** cycle, and
   `judgeReverse` row 1 returns REACHED whatever silent says (parent M21). **The fix converts no
   production pass into a false refusal on World's own listener.**
2. **The fix's increment is ≈0 at the real cadence:** 0 of 400 faithful cycles changed. Under
   amplification, 18 and 12 of 200 child cycles changed (idle, load), and 0 and 1 of 200 mux
   cycles did. Production dials every 50 ms (V15), not continuously.
3. **A pre-existing hazard, not caused by the fix (R5).** A *foreign chatty* holder whose process
   exits mid-poll leaves a zero-byte reset, so today the attempt is SILENT (sticky refusal), not
   INCONCLUSIVE (retry): 3/100 idle and 13/100 loaded, from the post-connect path alone. srt's mux
   (parent V7) exits with srt while the poll is live, so on a mux collision (itself rare) about
   3–13% of collided attempts refuse startup. Unchanged here; §8 proposes a row.

**Mitigation: none needed.** ECONNRESET-only keeps P20: no connected dial passes its own attempt,
and a connect-time reset now refuses like its post-connect twin.

## 4. What the fix relies on (stated, not enforced)

"World's own probe client never answers foreign" holds because of the fixed script
(`c.resume(); c.end('WORLD-PROBE-NONCE '+nonce+'\n')`, `:82`), which runs on World's verified
node. **No code enforces it** (iteration-317 judge caveat). The fix relies on (a) `judgeReverse`'s
child-first order (`:128-131`, parent M21); (b) the child's listener being unreachable in a
correct sandbox (parent §4.3, M0 V9); (c) ECONNRESET never arising with nothing listening: darwin
0 of 2,013,696 under load (V7), Linux 0 of 7,325,743 idle and loaded (V14).

If (c) failed on some host, the effect would be a refusal, never a pass, so it fails closed.

## 5. Mutation table (continues the parent's M1–M21)

| ID | Mutant | Killer, and the assertion that fires |
|---|---|---|
| M22 MUT-DIALERR-DROPPED | `pollLoopback` ignores dial errors (**today's code**) | T9 `connect-reset`: `seen` = `{}`, want `{silent}` |
| M23 MUT-RESET-NOT-SILENT | `dialErrIsSilent` returns false for ECONNRESET | T8 `op-reset`, `bare-reset`, `wrapped-reset`; T9 `connect-reset` |
| M24 MUT-RESET-NO-UNWRAP | `err == syscall.ECONNRESET` (no `errors.Is`) | T8 `op-reset` and `wrapped-reset` (row `bare-reset` stays green: the in-table control) |
| M25 MUT-REFUSED-IS-SILENT | ECONNREFUSED is silent | T8 `op-refused`, `real-refused`; T9 `refused-only`; T5 `clean` and every T2 row (free ports dial ECONNREFUSED, so the leg refuses SILENT instead of `""`) |
| M26 MUT-ANY-ERR-IS-SILENT | `err != nil` is silent | T8 `op-refused`, `op-einval`, `op-aborted`, `timeout`, `addrnotavail`; T5 `clean` |
| M27 MUT-EINVAL-IS-SILENT | EINVAL is silent | T8 `op-einval`; T9 `einval-then-refused` |
| M28 MUT-TIMEOUT-IS-SILENT | a timeout is silent | T8 `timeout` |
| M29 MUT-CONNECT-RESET-IS-FOREIGN | a connect-time reset sets `foreign` | T9 `connect-reset`: want exactly `{silent}`, not `{foreign}` (the leg would retry, then pass on attempt 2) |

## 6. Tests

| Test | Needs node? | Kills |
|---|---|---|
| T8 `TestExecProbeReverseDialErrIsSilent` | no (pure) | M23, M24, M25, M26, M27, M28 |
| T9 `TestExecProbeReversePollRecordsConnectReset` | no (injected dial) | M22, M23, M25, M27, M29 |
| T7 `TestExecProbeReverseLegSilentHolderRefusesSticky` (existing, unchanged) | no | parent M16, M17, M19, M20 |

- **T8**, a table over `dialErrIsSilent`. True: `op-reset` =
  `&net.OpError{Op:"dial", Net:"tcp", Err: os.NewSyscallError("connect", syscall.ECONNRESET)}`,
  `bare-reset`, `wrapped-reset` (`fmt.Errorf("x: %w", op-reset)`). False: `nil`, `op-refused`,
  `real-refused` (a real dial to a freed port, the known-positive shape control), `op-einval`,
  `op-aborted`, `addrnotavail`, `timeout` (a real `DialContext` with an expired context). V10 ran
  every shape except `addrnotavail`.
- **T9** calls `pollLoopback(port≠0, want, fake)` and stops it at once; the poll dials once before
  it checks stop (`:690-693`), so ≥1 dial is guaranteed with no timing. Rows: `connect-reset`
  (call 1 `op-reset`, later `op-refused`) wants exactly `{silent:true}`; `refused-only` and
  `einval-then-refused` want `{}`.
- **T7 `reset-holder` after the fix.** With one sequential dialer (T7's shape), every outcome
  against that holder was a connect (then zero bytes and a reset: silent) or a connect-time
  ECONNRESET, on darwin (V4) and Linux (V14 `-dialers 1`, idle and loaded). The fix makes both
  SILENT, so the row is **correct on every observed sequential outcome**. It is not strictly
  deterministic: a 200 ms connect timeout records nothing, and Linux produced 143/179 against this
  holder under 4 spinning dialers (V14). The deterministic proof is T9.
- **Red at base.** T9 cannot compile at base (the seam is new). M22 *is* today's behaviour on the
  new seam, and AC-MUT shows it red. T7 is not a reliable red on darwin: 0/300 this session at
  32 `yes` (V8), 5/300 for the iteration-317 controller, and 0/300 and 1/600 for its judge. [r2]
  **On Linux under load it is a reliable red at base:** at least 7 `--- FAIL: …/reset-holder` in
  `-count=300`, each with the CI signature (V14). That figure is a floor, counted from a
  `tail -30`.
- **AC-CI: yes.** T8 and T9 join step 8's `tests=` list (`ci.yml:247`), because that step fails
  unless each listed test prints `--- PASS:` (`:255-257`, V11). The parent listed all its
  reverse-leg tests, the pure T1 included, so a SKIP fails CI. The new tests follow the same rule.
- **AC-STABLE.** `go test ./host/broker -run '^TestExecProbeReverse' -count=20 -v`: every test
  passes 20/20. Then T7 `reset-holder` `-count=300` under 32 `yes`: 0 FAIL. That is a stability
  read, not the proof.
- **AC-LINUX [r2].** The fix PR's linux/amd64 CI must show `--- PASS:` for T7 `…/reset-holder`,
  T8 and T9 in both "go host build + test gate" (step 8's list plus the subtest line) and "go race
  gate". A red there blocks the merge: re-measure the Linux dial-error classes (V14's probe) and
  amend §2's table with the observed class before any retry.
- **AC-LINUX-LOAD [r2] (controller-run, before merge).** V14's workflow pattern on a throwaway
  branch at the **fixed** head: start 2×nproc `yes`; run `go test ./host/broker/ -run
  '^TestExecProbeReverseLegSilentHolderRefusesSticky$/^reset-holder$' -count=300 -v`; kill the
  `yes` PIDs. Count over the **whole** log, not a tail. Want 0 `--- FAIL: …/reset-holder` and
  exactly 300 `--- PASS: …/reset-holder` (so a run that matched nothing cannot pass). Base gives
  ≥7 FAIL (V14), which is the paired red.
- **AC-MUT.** Apply M22–M29 one at a time. Each reds exactly its named killer, and the tree is
  restored byte-identical (sha256 of both Go files). **AC-GATE**: `go build ./... && go test
  ./...` and `./scripts/verify_ail.sh` are green; no `.ail` is touched.

## 7. Files touched, docs, conflict surface

| File | Change | LOC |
|---|---|---|
| `host/broker/exec_probe.go` | `dialErrIsSilent`; `pollLoopback` dial seam and error branch; the `reverseLeg` call passes `nil` | +12 / −3 |
| `host/broker/exec_probe_test.go` | T8, T9 | +75 |
| `.github/workflows/ci.yml` | `:247` gains `TestExecProbeReverseDialErrIsSilent TestExecProbeReversePollRecordsConnectReset` | ±1 |
| `docs/QUICKSTART.md` | §10 arm 6 (`:1152`): "answered nothing to World's hello" gains "(or reset the host's dial)" | ±1 |

**S7:** one clause, since a resetting holder is now provably in the "answered nothing" class; the
refusal text is unchanged. **S2/S3:** as the parent (§8), a repair inside the host's effect
boundary with no new surface, so not a package. **S1/S4/S5:** n/a, no `.ail`.

**Caller census** (V12). `grep -rn 'pollLoopback\|classifyReverseDial\|judgeReverse' host/` → 13
lines: 11 in `exec_probe.go` and 2 in `exec_probe_test.go:451-452` (T1's `judgeReverse` calls).
`pollLoopback` has exactly **one** call site, `reverseLeg` at `:566`. `classifyReverseDial` has one
(`:706`). Positive control in the same session: the `runBounded` census is `9`. Only the reverse
leg changes; arms 1–5, 7 and 8 are untouched.

## 8. Residuals

- **R5, a foreign holder's exit is classed SILENT (pre-existing; §3 reading 3).** 3/100 idle,
  13/100 loaded per collided attempt (mux-shaped stand-in, darwin). Successor row
  `w-exec-probe-holder-exit-reset`; candidate: a zero-byte reset after a foreign answer *in the
  same attempt* is foreign, since the child cannot co-hold that port (EADDRINUSE) barring
  unmeasured `SO_REUSEPORT`. It needs real srt and Linux, so it is out of scope.
- **R6, darwin EINVAL** at a listener's close under load (104 dials in the amplified child run,
  V6) stays unclassified. It cannot be told apart from the no-listener EINVAL (V7), and the
  child's nonce dominates anyway.
- **R7, Linux (downgraded in r2; now measured, V14).** The Linux connect-time class is
  ECONNRESET, as on darwin. What remains is (i) the spinning-dialer timeouts (§2), which T7's
  sequential shape did not show, and (ii) one runner kernel (`7.0.0-1017-azure`). AC-LINUX and
  AC-LINUX-LOAD are the gates for both.

## 9. Verification Log

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | The CI failure on `dev` | VERIFIED BY CONTROLLER (`gh api` job 123313206642); excerpt `sed -n 400,422p /tmp/w318/race.txt` | `07:14:56 --- FAIL: …SilentHolderRefusesSticky/reset-holder (0.21s)`; `exec_probe_test.go:664: Why = "", picker calls 1; want the sticky silent refusal on attempt 1, calls 1`; `FAIL github.com/…/host/broker 227.232s`; the neighbouring packages are `ok` |
| V2 | Base and tooling | `git rev-parse HEAD; git status --porcelain \| wc -l`; `which docker podman colima orb limactl`; `go version` | `df5a0671dadd…`; `0`; no Linux tool found; `go1.26.6 darwin/arm64` |
| V3 | Today's poll, ordering and child close | `grep -n 'poll := pollLoopback\|r := pp.run(ctx, argv\|seen := poll()\|v.close();r(null)\|net.DialTimeout("tcp", addr' host/broker/exec_probe.go`; read `:655-728` | `:84 …setTimeout(()=>{v.close();r(null)},+ms)`; `:566 poll := pollLoopback(…)`; `:567 r := pp.run(…)`; `:568 seen := poll()`; `:705 if c, err := net.DialTimeout(…); err == nil {`, with no else branch |
| V4 | Connect-time reset re-measured (iteration 317's probe) | `/tmp/w318-designer/rst/rst` built from `design_docs/verification/world-iter317/row232-dial-reset-probe.go.txt` (sha256 prefix `1ef0d93ea139`), idle, then under load | idle `dials=2000 connected=1957 dial_errors=43`, all `…nect: connection reset by peer`; load `connected=1921 dial_errors=79`, the same error |
| V5 | Close-time classes, idle. Source `designer-close-probe.go.txt`; runner `designer-runidle.sh.txt`; raw output `designer-probe-outputs.txt` (all in `design_docs/verification/world-iter318/`) | `cd /tmp/w318-designer/probe && ./probe -mode child -n 100 -hold 1000`; `-mode mux -n 100 -hold 1000`; `-mode child -n 200 -hold 100 -tight 4`; `-mode mux -n 200 -hold 100 -tight 4` | child faithful: REFUSED 242, child 1901, silent-postconnect 2 (2 cycles), `CYCLES child-seen=100 silent(pre)=2 silent(post)=2`. mux faithful: foreign 1898, silent-postconnect 3, `pre=3 post=3`. child amp: ECONNRESET 115 (80 cycles), `child-seen=200 pre=167 post=185`. mux amp: ECONNRESET 1, `pre=198 post=198` |
| V6 | The same, under load. Runner `designer-runload.sh.txt` (32 `yes`, PIDs in `load.pids`, `trap kill`); output `designer-probe-outputs.txt` | as V5 | child faithful: REFUSED 784, child 1813, `pre=0 post=0`. mux faithful: silent-postconnect 13, `pre=13 post=13`. child amp: ECONNRESET 30 (22 cycles), EINVAL 104 dials over 81 ports, `child-seen=200 pre=179 post=191`. mux amp: ECONNRESET 1, EINVAL 8, `pre=194 post=195`. 0 `yes` PIDs alive afterwards |
| V7 | **No-listener and reset-holder controls under load.** Source `designer-ctl.go.txt` (sha256 prefix `94c2572b08aa`); runner `designer-runctl.sh.txt`; output `designer-probe-outputs.txt` | `/tmp/w318-designer/ctl/runctl.sh` (32 `yes`): `./ctl -mode none -secs 45`; `./ctl -mode rst -secs 45` (4 spinning dialers; errno via `errors.As(err, &syscall.Errno)`) | none: `connection refused 2013574`, `invalid argument 122`, with **no** reset and no timeout. rst: `connected 492165`, `connection reset by peer 96566`, nothing else. 0 `yes` PIDs alive afterwards |
| V8 | T7 at base under load, and the focused suite idle | `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/broker/ -run '^TestExecProbeReverseLegSilentHolderRefusesSticky$/^reset-holder$' -count=300` (inside V6's load); then `go test ./host/broker/ -run '^TestExecProbeReverseLeg' -count=1 -v` idle | load: `ok … 57.503s`, rc 0 (0/300 FAIL). idle: 6 `--- PASS`, 0 FAIL, 0 SKIP, rc 0 |
| V9 | Go's connect path and per-GOOS errnos | `sed -n 44,140p $(go env GOROOT)/src/net/fd_unix.go`; `diff` against go1.26.4's copy; `grep -n 'ECONNRESET *= Errno\|ECONNREFUSED *= Errno\|ECONNABORTED *= Errno' …/zerrors_{linux_amd64,darwin_arm64}.go` | `:72` and `:127-131`: a non-zero `SO_ERROR` → `os.NewSyscallError("connect", err)`; `fd_unix-identical`; darwin `ECONNRESET = Errno(0x36)`, `ECONNREFUSED 0x3d`, `ECONNABORTED 0x35`; linux `0x68`, `0x6f`, `0x67` |
| V10 | The `errors.Is` matcher on each shape; it compiles for Linux. Source `designer-is.go.txt` (sha256 prefix `b9496cec44ad`) | `cd /tmp/w318-designer/is && go run . && GOOS=linux GOARCH=amd64 go build -o /dev/null .` | `nil false`; `op-reset true`; `bare-reset true`; `op-refused false`; `real-refused false` (`connect: connection refused`); `op-einval false`; `op-aborted false`; `real-timeout false` (`i/o timeout`); `wrapped-fmt true`; `linux-amd64-compiles` |
| V11 | CI step 8 guard and list | `sed -n 236,262p .github/workflows/ci.yml` | `:247` lists the parent's T1–T7; `:251` `go test ./host/broker/ ./host/daemon/ … -run "^(${tests// /\|})\$"`; `:255-257` `grep -q -- "^--- PASS: $t "` … `did not PASS on linux (SKIP?)` |
| V12 | Caller census, with a control | `grep -rn 'pollLoopback\|classifyReverseDial\|judgeReverse' host/ \| wc -l`, then the lines; `grep -rn 'func runBounded\|\.runBounded(' --include='*.go' host \| wc -l` | `13`: `exec_probe.go:114,125,128,566,573,655,660,690,691,694,706`, `exec_probe_test.go:451,452`; control `9` |
| V13 | QUICKSTART edit site | `sed -n 1146,1153p docs/QUICKSTART.md` | `:1151-1153` "`reverse leg: a listener on … answered nothing to World's hello` means a process that never answers held the port …" |
| V14 | **[r2] Linux dial-error classes, and T7 at base under load.** VERIFIED BY CONTROLLER: throwaway branch `probe/world-iter318-linux-errno` (`cf93d0f5`), run 40842980481. Log: `design_docs/verification/world-iter318/linux-errno-probe-run40842980481.txt`. Probe: `linux-errno-probe.go.txt` (sha256 prefix `cdae00854dcb`), an accept-then-`SetLinger(0)`-`Close` listener (T7's shape) | `for m in none rst; do for d in 1 4; do /tmp/probe -mode $m -dialers $d -secs 30; done; done`, idle, then under 2×nproc `yes`; then under load, `go test ./host/broker/ -run '^TestExecProbeReverseLegSilentHolderRefusesSticky$/^reset-holder$' -count=300 2>&1 \| tail -30` | Log `:217-219`: `Linux … 7.0.0-1017-azure … x86_64`, `4`, `go1.26.6 linux/amd64`. Idle (`:226-236`): none d1 `connection refused(0x6f) … 1328238`; none d4 `4085689`; rst d1 `connected 642869`, `connection reset by peer(0x68) is-ECONNRESET=true 1742`; rst d4 `connected 1488248`, reset `342`, `timeout 143`. Load (`:245-255`): none d1 refused `638159`; none d4 `1273657`; rst d1 `connected 195092`, reset `147968`; rst d4 `connected 488125`, reset `1184`, `timeout 179`. No other class anywhere. T7 (`:264-279` and on): 7 `--- FAIL: …/reset-holder` in the tail, each `exec_probe_test.go:664: Why = "", picker calls 1`; `FAIL …/host/broker 5.693s`. Controller: dev's failed job re-ran green on attempt 2 |
| V15 | **[r2]** The cadence: a 1000 ms hold, a 50 ms poll and a 200 ms dial | `grep -n 'const execProbeListenMS\|time.After(50 \* time.Millisecond)\|200\*time.Millisecond' host/broker/exec_probe.go` | `:98 const execProbeListenMS = 1000`; `:705 … net.DialTimeout("tcp", addr, 200*time.Millisecond)`; `:719 case <-time.After(50 * time.Millisecond):`. So about 1000/50 = 20 dials per attempt |

## 10. Quorum log

**Round 1** (artifact `design_docs/verification/world-iter318/quorum/r1.json`): **BLOCKED**.

Seats: `oc-glm-5-3` REJECT, `oc-kimi-k3` REJECT, `gemini-4-argon` PASS, `gpt6-1-sol` absent
(unreachable).
- **Both rejects rested on one premise, and it was REAL:** the Linux connect-time errno, on the
  only platform where the defect was observed, was asserted from memory (old §2, R7). glm: the
  efficacy direction was unbounded (the reset could surface as ECONNREFUSED, which the fix
  drops), so it had to be a gate, not a residual; the `:98` cadence had no V-row. kimi: "the
  CI signature is byte-identical to darwin" proved nothing, since the text fits any dropped dial
  error. gemini (PASS): the scratch probes lived only in `/tmp`.
- **r2:** V14 (the controller's Linux CI measurement; stronger than the kernel-source reading glm
  asked for) refutes the ECONNREFUSED worry, with 0 against the holder. §1, the §2 table, the §2
  Linux paragraph and §4(c) are rewritten on it; the timeout row records the Linux
  spinning-dialer timeouts; kimi's wording replaces "byte-identical". §6 gains AC-LINUX and
  AC-LINUX-LOAD and records Linux T7 as a reliable red at base. R7 is downgraded. V15 is new. The
  scratch sources, runners and raw outputs are checked in as `designer-*.txt`, and V5–V7 and V10
  point at them. The T8 note is corrected (V10 did not run `addrnotavail`).
