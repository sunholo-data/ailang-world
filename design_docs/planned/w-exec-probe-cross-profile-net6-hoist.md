# w-exec-probe-cross-profile-net6-hoist — bind every profile's arm-6 listener before the fan-out, so cross-profile port exclusion holds by construction

- **Status**: Planned (designed iteration 340, r1). Queue row 242. The successor row named in
  `design_docs/planned/w-exec-probe-reverse-leg-pins.md` §3.4 (row 230, the parent: its R2 row,
  V14–V18 and P4). It does not touch the reverse leg's retry or reset semantics, which row 233
  (`w-exec-probe-holder-exit-reset.md`, parked, D-WORLD-105 OPEN, default C) owns.
- **Item**: Clause 3 (the Workspace.Exec fail-closed startup probe). This is a structure change on a
  fail-closed startup gate, so it needs its own quorum.
- **Estimated**: ~0.25 day. Production +23/−18 lines in `host/broker/exec_probe.go`, with no
  behaviour change on any failure path (§3). One new test file,
  `host/broker/exec_probe_net6_test.go` (216 lines, two tests), and two names added to the `tests=` list at
  `ci.yml:247`. The byte-exact prototype is `design_docs/verification/world-iter340/designer-proto.patch.txt`
  (V11).
- **Measurement base**: `e22e0794` (= origin/dev), worktree `sprint/world-iter340-net6-hoist`.
  `exec_probe.go` sha256 prefix `96f33172667d27de` (V1).
- **Instruments**: `go1.26.6 darwin/arm64`, 16 cores, node v26.12.1, git 2.54.0, darwin ephemeral
  range 49152–65535 (V1). The rig was shared. Load averages at the timed runs read 5.0–6.2 (V1, V7,
  V9), and a live `ailang-worldd serve` (pid 79864) was running throughout (V13). No Linux tool is
  available locally, so every Linux statement here is **UNMEASURED**. Prototypes ran on sibling
  `git archive` copies of `e22e0794` (`.scratch-iter340-designer-{head,proto,eqhead,eqproto}`),
  never under `/tmp`. All four copies were deleted afterwards. Raw outputs and scripts are banked as
  `design_docs/verification/world-iter340/designer-*.txt`.

---

## 1. What row 242 says, and what this session measured

| # | Claim | Verdict at `e22e0794` | Evidence |
|---|---|---|---|
| 1 | Exclusion holds by statement order: every profile binds `net6` in `probe()` (`exec_probe.go:486`) before its own Run A (`:496`), and profiles fan out at `:253` | **CONFIRMED.** `newProbeListener` has exactly one call site, `:486`, inside `probe()`, which runs inside the goroutine started at `:253`. No test calls `probe` directly | V2 |
| 2 | No current test guards the ordering | **CONFIRMED (re-run of the controller's premise).** A 2 s sleep immediately before `:486` keeps `go test -count=1 -run ExecProbe ./host/broker/` at rc 0. The unmutated control is also rc 0 in the same harness | V4 |
| 3 | The hazard needs a *slow pre-bind step* | **NEW: production already has one, and only for some profiles.** For a profile that declares `protected`, `NewExecHandler` (`:408`) runs the protected baseline, which runs Git (`handlers_exec.go:200` → `exec_protected_record.go:36,40`) **before** that profile's `net6` bind. An unprotected profile runs no Git there. So the per-profile pre-bind time is not uniform. With that Git held open until the first pick (no sleep inserted, only the existing `worktreeGit` seam), the new test T1 reds at HEAD 10 of 10 times: the unprotected profile picks while the protected profile's `net6` does not yet exist | V5, V3 |
| 4 | Is #3 live in production? | **Not on this rig. It is an exposure, not an incident.** The real preflight Git took 8.5 / 11.1 / 14.0 ms (min / median / max, 30 runs, V9). The parent's worst measured margin was 216.9 ms, for unprotected profiles only (parent V16). So a real Git would have to run about 20× slower than measured before the order inverts. The margin argument therefore rests on an unmeasured product: protected-profile Git latency multiplied by load | V9; parent V16, V17 |

## 2. The mechanism being removed

A reverse-leg pick is `freeLoopbackPort` (`:678`): `bind(0)`, read the port, close it. `bind(0)` never
returns a port a live listener holds (parent V15: 0 hits in 200,000 picks, against 768 for the same
ports once released). So a pick cannot land on a `net6` that **already** exists. The hazard is the
other order. Profile X picks port P and closes it. Later, profile Y's `net6` `bind(0)` returns P.
X's child then binds P for its listen (the reverse leg), and World dials it. Now Y's counted
listener and X's reverse leg share a port, and parent V18 shows that both profiles refuse
(`accepts 9` for the owner; SILENT for the poller).

At HEAD, every profile's bind sits behind its own pre-bind work: `randomHex`, two `os.Mkdir`,
`NewExecHandler` (and, for protected profiles, Git), the cache `MkdirAll`s and two `controlWrite`s.
Every pick sits behind Run A, a fixed 200 ms sleep, and Run B. Exclusion holds only while every
profile's pre-bind work finishes before any other profile's first pick.

## 3. Design

### 3.1 The structural change (the planner pins these bytes: V11)

`ProbeExecProfilesEach` binds all N listeners after the shared scratch is planted and before the
first `go` statement. It hands listener `i` to profile `i`'s goroutine and closes all N in one
parent-level `defer`, which runs after `wg.Wait()`:

```go
	errs = make([]error, len(profiles))
	// Every profile's arm-6 net6 listener is bound here, before any profile
	// starts (row 242): bind(0) never returns a held port (row 230 V15), so
	// no profile's reverse-leg pick can return a port a net6 listener takes
	// later, whatever a profile does before its Run A. A bind failure is
	// carried in the listener and judged by its owner's arm-6 control leg.
	// All of them close here, after every profile's probe has returned.
	net6 := make([]*probeListener, len(profiles))
	defer func() {
		for _, l := range net6 {
			l.close()
		}
	}()
	for i := range profiles {
		net6[i] = newProbeListener()
	}
	var wg sync.WaitGroup
	for i, p := range profiles {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs[i] = probeExecProfile(ctx, cfg, env, p, net6[i])
		}()
	}
	wg.Wait()
	return errs, nil
```

There are three other edits. `probeExecProfile` gains a final parameter `net6 *probeListener` and
passes it on as `pp.probe(ctx, env, net6, worktree, sibling, id)`. `probe` gains the parameter
`net6 *probeListener` after `env`. In `probe`, the lines `net6 := newProbeListener()` and
`defer net6.close()` (`:486–487`) are deleted. `netCtl := net6.control()` stays **at the same
point**, before Run A, and the comment above it becomes the four-line comment in the patch. Row
230's P4 comment at the fan-out (`:247–252`) is replaced by the block comment above. Nothing else
changes. Go 1.26 (`go.mod`) gives each loop iteration its own `i`, so the closure captures are
per-profile, as they are today at `:255`.

### 3.2 Who owns each listener on every path

`newProbeListener` never returns nil and never panics. On a failed `net.Listen` it returns a listener
with `err` set and `ln` nil (`:632–638`). `close` is then a no-op (`:671–675`), so the parent's
close loop is total.

| Path | Listeners bound | Closed by | When | Measured |
|---|---|---|---|---|
| `cfg.Sandbox == nil` (`:235`) | none | — | — | read: the return precedes the bind loop |
| `plantExecProbe` fails (`:240`) | none | — | — | read: the bind loop follows `return nil, err`. **UNMEASURED dynamically** |
| `newProbeListener` fails for profile k while 0..k−1 succeed | k carries `err`; the others are live | the parent `defer` (k's close is a no-op) | after `wg.Wait()` | V8 (failure injected into every bind) |
| `probeExecProfile` early return: `randomHex` (`:398`), `os.Mkdir` (`:404`), `NewExecHandler` (`:412`), `MkdirAll` (`:426`) | that profile's, unused | the parent `defer` | after `wg.Wait()` | V8 (the `NewExecHandler` return, via a missing `read_root`): `listeners-left=0` |
| Normal path | all | the parent `defer` | after `wg.Wait()`, so after every owner's Run C, Run D and arm 8 | V6 (T1's after-return check), V7 (MUT-4 reds it) |
| A goroutine panics | — | the process dies, as today (no `recover` anywhere in the path) | — | read |

No listener closes before its owner's Run C finishes: the only close is ordered after `wg.Wait()`,
which comes after every owner's `wg.Done()`. MUT-6 (close before `wg.Wait()`) reds both T1 and the
existing pass-through test (V7).

### 3.3 Fail-closed semantics: unchanged on every path

| Failure | Before (HEAD) | After (hoist) |
|---|---|---|
| Profile k's bind fails | `errs[k]` is an `*ExecProbeError` whose `arm6-network` names `control: World's live listener did not accept an unsandboxed connect: <err>` first, followed by the raw, proxy and reverse legs | **identical** (V8: normalised texts byte-equal; the arm set and the `N of 7` count are equal) |
| An early return in `probeExecProfile` | a plain wrapped `error` (for example `broker: exec probe for profile "badroot": broker: exec handler: read_root …`), with no arm list | **identical** (V8) |
| `plantExecProbe` fails | `(nil, err)` | identical (the bind loop is never reached) |
| A pass, or any arm failure | `nil`, or `*ExecProbeError` with arms in `execProbeOrder` | identical; the arm code is untouched |

`errs[i]` stays `profiles[i]`'s own result. `ExecProbeError`, `ExecProbeArmFailure` and the arm
order are untouched. `ProbeExecProfilesEach` and `ProbeExecProfiles` keep their signatures. So
`ExecToolProbe` (`exec_tools.go:24,134`) and the one production caller (`exec_startup.go:235`) do
not change. Only the unexported `probeExecProfile` and `probe` gain a parameter. Nothing outside
`exec_probe.go` calls either of them (V2).

**One deliberate change is lifetime.** Today a listener lives from its own bind, after the pre-bind
work, until `probe` returns. After the hoist it lives from before any profile starts until all
profiles finish. In process, this only strengthens exclusion. Its only cross-process effect is in
§6 R1.

## 4. The by-construction argument and the mutation property

**Claim.** At every reverse-leg pick by any profile, every profile's `net6` socket is bound and
held, or failed to bind. So no `net6` bind can happen after any pick.

**Proof.** (1) All `newProbeListener` calls complete in the parent before its first `go` statement,
and in the Go memory model a `go` statement happens-before the goroutine it starts. (2) Picks happen
only inside those goroutines: `pick()` (`:582`) is reached only from `reverseLeg`, which is reached
only from `probe` → `probeExecProfile` → the goroutine (V2). (3) The only close is the parent's
`defer`, which runs after `wg.Wait()`, and `wg.Wait()` happens-after every `wg.Done()`. So no `net6`
is created after a pick, and none is released before its owner finishes. With parent V15 (`bind(0)`
never returns a held port), no pick equals any live `net6` port, and no later `net6` can take a
picked port. ∎ The only assumption left is V15, the same kernel property HEAD relies on.

**The row's mutant, "move the bind back after a slow pre-bind step".** After the hoist it splits into
two kinds of edit:

- **Adding a slow step anywhere is harmless by construction.** A sleep in the parent before the bind
  loop (MUT-1) delays every bind, but no goroutine exists yet. A sleep at the old bind site in `probe`
  (MUT-1b, the controller's mutant translated to the new structure) delays nothing that matters,
  because `net6` already exists. T1 stays green on both, as it should (V7).
- **Reintroducing the hazard takes moving the bind.** The bind call has to land after the fan-out:
  back into `probe` (MUT-2), into the goroutine (MUT-3), into `probeExecProfile` after
  `NewExecHandler` (MUT-3b), or into a helper in another file called from the goroutine (MUT-7).
  **T2 reds all four** (V7, V10). T2 is a source pin. It requires exactly one `newProbeListener` call
  in the package's non-test files, inside `ProbeExecProfilesEach`, outside any function literal, and
  lexically before that function's first `go` statement. **T1 reds MUT-2 and MUT-3b**, where a
  production slow step (the protected Git) precedes the moved bind. T1 does **not** red MUT-3 or
  MUT-7, where the moved bind is the goroutine's first statement and nothing slow precedes it; T2 is
  their only guard. If anyone then adds a slow step before such a moved bind, that combined edit is
  MUT-3b's shape, and T1 reds too.

**Why T1 reaches production wiring, not a seam (iteration 339's lesson).** T1 injects no
listener. It calls the public `ProbeExecProfilesEach` with two real profiles, the pass-through shim
and a logging probe client that `exec`s World's own node. It uses only two seams that already
exist. `worktreeGit` holds the protected profile's real `NewExecHandler` Git open until the first
pick. `cfg.reversePort` records a snapshot at each pick and then returns the production
`freeLoopbackPort()`. Each snapshot lists the process's own listening 127.0.0.1 TCP sockets,
read from its file descriptors (`getsockname`, `SO_TYPE`, `getpeername` = ENOTCONN, because darwin
returns "protocol not available" for `SO_ACCEPTCONN`; V12). The `net6` ports are taken from the
`connect <port>` argument that each profile's real Run A passed to the client, not from the code
under test. T1 asserts five things:

1. exactly two distinct `net6` ports (MUT-5 reds this);
2. the first pick landed while the protected profile's Git had started and had not finished (the
   non-vacuity guard);
3. both ports were live listeners at every pick;
4. no listener the call created outlives it (MUT-4 reds this);
5. the existing `assertProbeCleaned`.

Positive control: T1 reds at HEAD 10/10. Each red's snapshot contains the picking profile's own
`net6` port (for example `live map[58850:true]`), which shows the instrument sees a production
listener when one exists (V5). Negative control: 20/20 green on the prototype (V6).

## 5. Test plan and mutation table

- **T1** `TestExecProbeNet6BoundBeforeAnyPick`. Behavioural, as above. It takes about 2.7–3.0 s
  (V6). It needs node and skips by name without one (`probeNode`), like the other shim tests.
- **T2** `TestExecProbeNet6BindSiteIsBeforeTheFanOut`. The source pin above. It is pure and cannot
  skip. It refuses vacuity: it fails if it parsed fewer than 100 declarations. The package already
  uses `go/parser` source pins (`handlers_parallel_guard_test.go`, V2).
- Both names join the verbose `tests=` list at `ci.yml:247`, so that a Linux skip cannot read as a
  pass (the family's rule, parent §4).
- Existing tests are unchanged. In each mutant row, the pass-through pin (`accepts 2` / `accepts 3`)
  was run as its own named test.

All rows below were measured on the prototype by `designer-run-mutants.sh.txt` and
`designer-mutate.py.txt` (V7) or by the follow-up runs (V10). Each mutant was applied to a pristine
copy and each named test was run on its own with `-count=1`. Then the `-run ExecProbe` subset was
run. Both copies were restored and verified with `cmp` afterwards.

| Mutant | Edit (on the hoisted prototype unless stated) | T1 | T2 | pass-through `TP` | `-run ExecProbe` | Expected → observed |
|---|---|---|---|---|---|---|
| CONTROL | none | 0 | 0 | 0 | 0 | green → green |
| HEAD (positive control) | unhoisted `e22e0794` + the new test file | **1** | **1** | — | 0 (without T1 and T2) | red → red; also 10/10 red at `-count=10` (V5) |
| HEAD-SLEEP (row's own, as the controller wrote it) | HEAD + `time.Sleep(2s)` before `:486` | **1** | **1** | — | 0 (without T1 and T2) | the old suite blind → it is (V4); the new tests red → red |
| MUT-1 | `time.Sleep(2s)` in the parent before the bind loop | 0 | 0 | 0 | 0 | harmless by construction → green |
| MUT-1b | `time.Sleep(2s)` at the old bind site in `probe` | 0 | 0 | 0 | 0 | harmless by construction → green |
| MUT-2 | `net6 = newProbeListener(); defer net6.close()` back in `probe` | **1** | **1** | 0 | **1** | red → red (`pick 0: net6 port … was not a live listener`) |
| MUT-3 | bind moved into the goroutine as its first statement | 0 | **1** | 0 | **1** | T2 red, T1 blind → as predicted |
| MUT-3b | bind moved into `probeExecProfile` after `NewExecHandler` | **1** | **1** | 0 | **1** | red → red |
| MUT-4 | parent never closes | **1** | 0 | 1, then 0 in 3 of 3 reruns | **1** | T1 red (`listeners left … [57297 57298]`) → red; the `TP` red is not attributable (§6 R2) |
| MUT-5 | every profile gets `net6[0]` | **1** | 0 | 0 | **1** | red (`want two distinct ports`) → red |
| MUT-6 | all closed before `wg.Wait()` | **1** | 0 | **1** (3 of 3 reruns) | **1** | red → red; `TP` misses `accepts 2` / `accepts 3` |
| MUT-7 | bind moved to a helper in another file, called from the goroutine | 0 | **1** | — | — | T2 red (the all-files scan) → red (V10) |

Every mutant that reintroduces the hazard (MUT-2, -3, -3b, -7) reds T2, and every one with a slow
pre-bind step also reds T1. Every slow-step-only mutant (MUT-1, -1b) is harmless by construction and
stays green.

## 6. Residuals and risks

- **R1, the cross-process residual (out of scope; unchanged in kind).** Another process's `bind(0)`
  can return our pick's port, and another process's dial can land on our `net6`. Neither is in-process,
  so the hoist cannot remove either. **The hoist does change one window, slightly.** A foreign dial
  that lands on a profile's `net6` between its bind and its `control()` now inflates the count before
  `control()` runs, so `control()` sees 2 and refuses (`accept count 2 after one connect, want 1`):
  fail-closed, never a false pass. That window was about zero at HEAD, since `control()` immediately
  follows the bind. After the hoist it equals the profile's pre-bind work: milliseconds, plus about
  11 ms of Git for protected profiles (V9). That is small next to the existing bind → `acceptsC`
  window of Run A + 200 ms + Run B + Run C + 200 ms, a span of seconds (parent V16). The tail (a
  finished profile's `net6` held until the slowest profile finishes) has no failure mode: that
  owner's counts are already sampled. **UNMEASURED:** the foreign-dial rate itself, on darwin or
  Linux.
- **R2, the pass-through `accepts 4` flake (controller-filed; not this row).** The controller saw
  `TestExecProbePassThroughShimFailsEveryConfinementArm` red 2 times in 43 runs with `accepts 4`
  where the pin wants `accepts 3`. **The hoist does not change it.** The test has one profile, so no
  cross-profile pick exists. The extra accept lands between the `acceptsA` and `acceptsC` samples,
  and the hoist leaves that window and both samples untouched. It moves only the bind earlier (see
  R1, which would surface as a *control* failure, not `accepts 4`). This session saw one unexplained
  `TP` red in the MUT-4 row, which cannot cause it: the leak happens after the test's only call
  returns, in a run of that test alone. Three reruns were green. The detail was overwritten before
  it was read, so **its text is UNMEASURED** and it is only consistent with the controller's rate
  (V7 note).
- **R3, T2 is lexical.** It pins where the call is written, not when it runs. It cannot see a call
  made through a function value. For example, `f := newProbeListener` in the parent followed by
  `f()` in the goroutine has no `CallExpr` on the identifier, so the call-site count drops to 0 and
  T2 **reds**. A raw `net.Listen` for `net6` written inside `probe` would bypass `newProbeListener`
  entirely. T2 would stay green (the parent's call is still there), but T1 reds it, as MUT-2 shows
  (the ignored parent listener is the same shape). These are reasoned, not measured, apart from
  MUT-2.
- **R4, T1's descriptor scan.** It is capped at 4096 descriptors or `RLIMIT_NOFILE`. A missed
  descriptor can only cause a false *red*. A false green would need a pending loopback dial whose
  local port equals a not-yet-bound `net6` port: negligible, and not measured. **Linux UNMEASURED:**
  `SO_TYPE` and `getpeername` → ENOTCONN on a listener are standard there. CI's verbose list (§5)
  is the first Linux run.
- **R5, holding listeners longer.** At most N loopback ports are held for the whole probe instead of
  each owner's span. N is the profile count. The full-suite comparison (V14) found no failure that the hoist could explain. Each side had one load-timed difference, and the proto one was green on rerun.

## 7. Milestones

- **M1 (single, ~0.25 day).** Apply `designer-proto.patch.txt` (V11; `git apply` reproduces both
  files byte-exact). Add T1 and T2 to `ci.yml:247`. Prove, one at a time: T1 and T2 red at the parent
  commit (HEAD row); MUT-2, MUT-3b and MUT-6 red T1; MUT-3 and MUT-7 red T2; MUT-1 and MUT-1b green
  on both. Then run `./scripts/verify_ail.sh` (no `.ail` changes) and `go build ./... && go test ./...`
  under the pinned `AILANG_BIN`.

## 8. Files

| File | Change |
|---|---|
| `host/broker/exec_probe.go` | +23/−18 in total. `ProbeExecProfilesEach` binds and closes; the `probeExecProfile` and `probe` signatures and the call change; `probe` loses its bind, its close and the row 230 P4 comment; the P4 comment at the fan-out is replaced |
| `host/broker/exec_probe_net6_test.go` | new: T1, T2 and `loopbackListenPorts` |
| `.github/workflows/ci.yml` | two names appended at `:247` |

## 9. What this item is NOT doing

It does not change reverse-leg retry or reset semantics (row 233, D-WORLD-105), the hello wait,
the pick function, arm texts or the arm order. It does not touch the cross-process residual (R1), the
`accepts 4` flake (R2, the controller's own row) or any `.ail` code. It adds no seam: T1 uses the
existing `worktreeGit` and `cfg.reversePort`.

## 10. Verification Log

`$V` = `design_docs/verification/world-iter340`.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Base, instruments, load | `git rev-parse --short HEAD`; `shasum -a 256 host/broker/exec_probe.go`; `$V/designer-instruments.txt`; `uptime` at each timed run (`$V/designer-mutants.txt`, `$V/designer-stability.txt`) | `e22e0794`; `96f33172667d27de`; go1.26.6 darwin/arm64, node v26.12.1, git 2.54.0, 16 cores, range 49152–65535; load 5.01 (mutants), 5.64 (stability) |
| V2 | The bind site, the fan-out, the pick and the callers at HEAD | `$V/designer-greps.txt` (exact greps inside) | `:486 net6 := newProbeListener()`, `:487 defer net6.close()`, `:488 netCtl`, `:496 Run A`, `:537 reverseLeg`, `:582 pick()`, `:253 go func()`, `:255 probeExecProfile`, `:258 wg.Wait()`. `newProbeListener`: 1 call site (`:486`) plus the definition (`:632`, the known positive). `ProbeExecProfilesEach` callers: `exec_startup.go:235`, `exec_tools.go:134`. `.probe(ctx` in tests: none, against `probeProtected(` in 2 files (the positive). Early returns at `:398 :404 :412 :426`. `go/parser` tests: 4 files. `t.Parallel()` files: 0 |
| V3 | Protected profiles run Git before their bind | same file, protected section | `handlers_exec.go:200 captureProtectedBaseline(ctx, "", …)`; `exec_protected_record.go:36 worktreeGit()`, `:40 protectedGitPreflight`; `worktree_delta.go:72 var worktreeGit` |
| V4 | The controller's premise, re-run | `$V/designer-run-mutants.sh.txt` HEAD and HEAD-SLEEP rows: `go test -count=1 -run ExecProbe -skip TestExecProbeNet6 ./host/broker/` → `$V/designer-mutants.txt` | `HEAD … rc=0 :: 0 FAIL lines`; `HEAD-SLEEP … rc=0 :: 0 FAIL lines` |
| V5 | T1 and T2 red at HEAD (positive control) | the HEAD rows of `$V/designer-mutants.txt`; `go test -count=10 -run '^TestExecProbeNet6BoundBeforeAnyPick$'` → `$V/designer-stability.txt` | rc 1 and rc 1; `PASS=0 FAIL=10`; first failure `pick 0: net6 port 58856 was not a live listener (slow profile in flight=true; live map[58850:true])` |
| V6 | The prototype is green and stable | `go test -count=20 -run '^TestExecProbeNet6' -v` → `$V/designer-stability.txt` | `rc=0`, T1 `PASS=20 FAIL=0`, T2 `PASS=20`, T1 2.67–2.97 s |
| V7 | Mutation table | `$V/designer-run-mutants.sh.txt` + `$V/designer-mutate.py.txt` → `$V/designer-mutants.txt`; reruns → `$V/designer-mutant-reruns.txt` | as §5. `# restored both copies`. MUT-4 `TP` reruns ok ×3; MUT-6 `TP` FAIL ×3 with `does not report "a raw connect … accepts 2"` |
| V8 | Failure-path equivalence | `$V/designer-equiv_test.go.txt` in HEAD and proto copies, each with the scratch-only injection shown in `$V/designer-equiv.txt` | both rc 0; `IDENTICAL normalised ZZ lines: 8`; bind-fails `arm6-network: control: World's live listener did not accept an unsandboxed connect: zz injected bind failure; …`; `listeners-left=0` in all 4 cases |
| V9 | Real preflight Git latency | `$V/designer-git-preflight.txt` (30× `git -c core.fsmonitor=false config --null --get-regexp '^filter\.'` in a non-repo dir) | `rc=1 min=8.5ms median=11.1ms max=14.0ms`, load 6.21 |
| V10 | T2 scans all package files; MUT-7 | `$V/designer-t2-allfiles.txt` | CONTROL rc 0; MUT-2, -3, -3b rc 1; HEAD rc 1; MUT-7 T2 rc 1 (`bindNet6Late (not before a go statement) at …zz_mut7.go`), T1 rc 0 |
| V11 | The patch is byte-exact | `git apply $V/designer-proto.patch.txt` on an `e22e0794` archive; `cmp` against both prototype files (`$V/designer-t2-allfiles.txt`, last line) | `APPLY-OK`; proto `exec_probe.go` sha256 `19de369381579387`, test `400faf0ffc25a3f0` |
| V12 | The darwin listener-detection instrument | `$V/designer-fdscan-probe.txt` | `SO_ACCEPTCONN=0 err=protocol not available SO_TYPE=1 … getpeername err=socket is not connected`; `loopbackListenPorts = map[50208:true] (contains the live listener)` |
| V13 | A live daemon on the rig | `pgrep -fl "ailang-worldd serve"` in `$V/designer-git-preflight.txt` | pid 79864 |
| V14 | Full suite, proto against HEAD, both without `AILANG_BIN` (so environmental reds are expected in both) | proto: `go build ./... && go vet ./host/broker/ && go test -count=1 ./...`; HEAD: `go test -count=1 -skip TestExecProbeNet6 ./...` → `$V/designer-full-suite.txt` | 79 failing top-level tests in each. `AILANG_BIN` messages: 119 and 119. Only in proto: `TestDoctorQuickIsFast` (`quick took 2.17s` against a 2 s bound, at load 5.8; it is green `-count=3` on a rerun, and `cmd/ailang-worldd` has 0 references to the probe). Only in HEAD: `TestWorktreeSnapshotRefusesAConfiguredFilter`. **The pinned-`AILANG_BIN` gate is UNMEASURED here; the executor runs it** |
