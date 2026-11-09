# w-exec-cache-seed-at-bind — clone the exec GOCACHE seed when the episode's exec handler is bound, not inside the first exec call

- **Status**: Planned (designed iteration 372, **r3**: r1 BLOCKED 3 of 3 at quorum round 1; r2
  fixes the two defects named there, marked **[r2]**; r2 BLOCKED 3 of 3 at round 2, and **r3** is
  the controller's narrow-refinement carve-out applying the three round-2 reviewers' own fixes,
  marked **[r3]**, §11). Queue row 257, in D-WORLD-108's scope after
  row 256. Row 256 (`w-daemon-startup-under-load`, draft PR #515) is held by an attended session;
  this design avoids every file #515 changes except one hunk in `host/daemon/workspace.go`, which
  does not overlap #515's hunks (§8).
- **Item**: Clause 5 (MET; agent UX). It moves one subprocess (the seed copy) out of a call's
  budget, makes the copy atomic, **[r2]** makes readiness depend on a completion marker the copy
  step writes (not on existence), and shows its outcome in `episode report`. It adds no record
  field, no class, no route and no store schema change. **[r3]** D-257-1 (§3.5) is settled as **A**:
  bind only at the first dispatched call that wires the exec tool; no GET gains a side effect.
- **Designer lane**: `claude:claude-opus-5-5` via Agent opus.
- **Estimated**: ~0.75 day. No prototype patch was written; the test plan (§5) is unrun. Production
  code: about +180/−40 lines, most of it in one new file `host/broker/exec_seed.go`.
- **Measurement base**: `ad13d835` (= origin/dev) for r1; **[r2]** the worktree was rebased onto
  `dc907675`, which changed only the charter, so every line number below still holds. Worktree
  `/Users/voightkampff/.ailang/state/world-iter372-row257`, branch `sprint/world-iter372-row257`.
  sha256 prefixes: `handlers_exec.go` `996c5ee11daed4b9`, `workspace.go` `1a178f8bc1e39f89`,
  `session_self.go` `6a98846db138b6d3`, `episode_report.go` `64128681d7bffa13` (V1).
- **Instruments**: `go1.26.6 darwin/arm64`, 16 cores, APFS (the seed, `/tmp` and the state dir
  are on one volume, `/dev/disk3s5`). Pinned `AILANG_BIN=$HOME/.pinned-ailang/ailang`
  (`AILANG v0.41.0`) for the load run. TMPDIR and GOTMPDIR were left at the rig default. The clone
  measured is the real `ailang-world` seed (226 MB, 2569 files), read only; scratch copies went to
  `/tmp/iter372-designer/` and were deleted after each sample. No Linux host was used, so every
  Linux statement is **UNMEASURED**.

---

## 1. What row 257 says, and what this session measured

Row 257: the GOCACHE seed is cloned on the episode's first exec call, inside that call's budget.
`prepareCache` clones the operator's seed under `execSeedBudget(ctx)`, the remaining call budget.
Under load the copy competes with the command it serves, and a failure repeats on every call.
Build: clone at bind or at `episode start`, report a failure there with the row-255 reason, keep
the per-call path as a check that the copy exists, measure the clone at bind under
package-parallel load, and say what bound it gets.

What the code does today (V2, V3):

1. The daemon has no bind step before the agent's first call. The coordinator calls
   `BinderFor(episode, grants)` inside every `Dispatch` (`host/coordinator/coordinator.go:316-317`).
   The production binder builds a fresh registry per call (`host/daemon/daemon.go:917`), and
   `workspaceTools.execHandler` builds the episode's `ExecHandler` once and caches it
   (`host/daemon/workspace.go:468`, construction at `:496`). **[r2]** V18 checks the trigger: the
   `Bind` at `coordinator.go:316-317` sits after the `if effectful` block, so it runs on every
   dispatched call, effectful or not; and `registry()` calls `execHandler` on every call
   (`workspace.go:450`) unless the workspace tools are off or the episode worktree does not
   resolve (`:429-436`, which return an empty registry). `execHandler` builds an `ExecHandler` only
   when exec is configured and the episode has a profile (`:469-487`). So "bind" is **the first
   dispatched call of any tool that wires the exec tool**: the first call of the episode, when the
   worktree resolves and a profile applies.
6. **[r2]** Nothing creates the project cache tree before the first exec call (V17). The only
   `MkdirAll` under `<state>/exec-cache/` is in `prepareCache` (`handlers_exec.go:532`, `:543`).
   `NewExecHandler` creates only the settings dir (`:189`, `<state>/exec/`), and `provisionEpisode`
   (`cmd/ailang-worldd/session_new.go:396`) never names `exec-cache`. So a copy started at bind or
   from a warm GET must create its own parent tree.
2. `NewExecHandler` does not touch the cache (`host/broker/handlers_exec.go:150-205`).
   `Execute` calls `prepareCache(ctx)` at `:308`, under the handler context the coordinator made at
   `host/coordinator/effectful.go:311` (`HandlerCap` = 10 s, `:47`). The profile's `timeout_ms` is
   9000 for `ailang-world` (V4).
3. `prepareCache` (`:528-562`) checks each cache with `os.Lstat(dst)`. If it is missing it runs
   `/bin/cp -cR seed dst` (or `cp -a --reflink=auto` off darwin) through `runBounded` with
   `execSeedBudget(ctx)` (`:565-573`), the call's remaining time. On error it removes `dst` and
   returns `FailSeedCache`. The next call starts the whole copy again from zero, under its own
   9 s.
4. A copy that dies part-way can leave a partial `dst`. `cp -cR` copies what it can and exits 1
   (V5: a seed with one unreadable file leaves `dst` holding the other two). Today `RemoveAll(dst)`
   covers the exit-1 case, but a daemon killed during the copy, or a `RemoveAll` that fails,
   leaves a partial `dst`, and every later call accepts it because the check is only `Lstat`.
5. The startup probe builds a scratch `ExecHandler` (`host/broker/exec_probe.go:452`) and never
   seeds: it makes the cache dirs itself (`:462-470`, "never seeded by the probe"). So the copy
   must not move into `NewExecHandler`, or every serve start would clone each profile's seed
   inside #515's startup budget.

Measured clone times (V6–V9), one sample = one `cp -cR` of the real 226 MB seed to a fresh path.
Each interval also includes one `python3` interpreter start (the timer), so these are upper
bounds; that overhead was not separated (**UNMEASURED**).

| Condition | N | min | median | max | load1 range |
|---|---|---|---|---|---|
| Idle rig | 15 | 0.345 s | 0.351 s | 0.472 s | 3.11–3.12 |
| `go test ./host/... -count=1 -p 8`, pinned AILANG, full run (rc 0, 309 s), one sample every ~6 s | 40 | 0.346 s | 0.385 s | 0.482 s | 3.99–8.34 (summed `ps %cpu` 38–685) |
| Same, first 80 s only, AILANG_BIN unset (tests failed fast; supplementary) | 15 | 0.348 s | 0.379 s | 0.628 s | 2.95–5.48 |
| 48 busy `yes` loops on 16 cores (3× oversubscribed) | 12 | 0.686 s | 1.014 s | 1.146 s | 29.9–56.6 |

Package-parallel `go test` at `-p 8` did not load this 16-core rig enough to slow the clone
(median 0.385 s against 0.351 s idle). CPU oversubscription did: about 3× at load ~30–57. The
incident regime (load 110–130, row 256) is **UNMEASURED**; I did not reproduce it.

## 2. The mechanism being changed

Two things are wrong, and the row names the first:

- **Wrong budget.** The copy runs on the call's clock. Under load a copy that needs more than the
  call's remaining time is killed, deleted and restarted by the next call, so it may never finish.
- **Wrong check.** `Lstat(dst)` accepts a partial copy left by a crash. That is a fail-open gap
  the row's "keep the per-call path as a check that the copy exists" depends on closing.
  **[r2]** Making new copies atomic is not enough: a partial `dst` written by the code before this
  change survives the upgrade's own restart. So existence alone must never mean ready. Ready means
  `dst` holds completion evidence that only a finished copy step writes (§3.1a).

## 3. Design

### 3.1 The seed state lives on the handler, and the copy runs on its own clock

A new file `host/broker/exec_seed.go` holds one small state machine per `ExecHandler`:

```go
type execSeed struct {
	mu       sync.Mutex
	attempts int           // copies started, ≤ execSeedMaxAttempts
	running  chan struct{} // non-nil while a copy runs; closed when it ends
	started  time.Time     // of the running or last attempt
	took     time.Duration // of the last finished attempt
	by       string        // "bind" or "call": who started the last attempt
	last     *HandlerFailure // nil after a success
	ok       bool          // the last attempt succeeded
}

const (
	execSeedBindBudget  = 60 * time.Second // §3.4
	execSeedMaxAttempts = 3                // §3.3
)
```

`(*ExecHandler).StartSeed()` is exported. **[r2]** If every seeded cache is *verified* (§3.1a),
it marks the state ready and spawns nothing, so a daemon restart over a verified `exec-cache` dir
costs no copy. A cache that exists but is not verified counts as missing.
Otherwise it starts attempt 1 in a goroutine and returns at once. The goroutine runs under
`context.WithTimeout(context.Background(), h.seedBudget)` (`h.seedBudget` defaults to
`execSeedBindBudget`; tests set it). It never uses a caller's context, so no call's deadline can
kill it. This root is the ctx-less binder's, like `NewExecHandler`'s baseline capture, and gets
one row in `host/store/context_roots_test.go` (`"host/broker/handlers_exec.go|NewExecHandler|Background"`
is the precedent at `:80`).

One attempt seeds every seeded cache in `sortedKeys` order. For each one it runs the steps
below. **[r2]** Any step's error ends the attempt with
`*HandlerFailure{Class: FailSeedCache, Detail: name + ": " + what + " (attempt k of 3)"}`, after
`RemoveAll(<staging>)`. No error path marks the cache ready.

0. **[r2]** `os.MkdirAll(<cache>, 0o700)`, where `<cache>` is `<state>/exec-cache/<ep>/<project>`.
   This also creates `<state>/exec-cache/<ep>/`, the staging path's parent. The mode matches
   `prepareCache`'s. It is idempotent, and nothing else creates this tree before the first exec
   call (§1 item 6). Error: `what` = `cache dir cannot be created`.
1. **[r2]** If the seed itself has a top-level `.world-seed-complete`, fail with `what` = `the seed
   holds World's marker; remove it from the seed`. Otherwise a copy that died part-way could carry
   a copied marker into `dst` (see the legacy case in step 2).
2. **[r2]** Classify `<cache>/<name>`. If it is absent, go on. If it is verified (§3.1a), this
   cache is done; skip to the next one. If it is present but not verified, it is **legacy or
   unverified**: a partial copy from the old code, a copy whose marker the agent removed, or
   anything else. Remove it with `RemoveAll`. If the removal fails, `what` = `an unverified copy
   cannot be removed`; the attempt fails and the cache is not ready.
3. **[r3]** Create this attempt's own staging directory, exclusively and uniquely:
   `<state>/exec-cache/<ep>/.<project>.<name>.seeding-<incarnation>-<k>`, where `<incarnation>` is
   a value fixed once per daemon process (its pid and start time in nanoseconds) and `<k>` is the
   attempt number. It is made with `os.Mkdir`, which fails if the path exists, so no other copy —
   an earlier attempt, a second handler, or a `cp` left running by a dead daemon — can own it. If
   `Mkdir` fails, `what` = `a staging directory cannot be created`. Before that, the attempt removes
   abandoned staging directories of this `<project>.<name>` from other incarnations or earlier
   attempts (bounded cleanup: one `RemoveAll` per directory, at most 8 per attempt; a removal
   failure is logged and does not fail the attempt, because staging is never published by anyone
   but its own attempt).
4. Run `cp -cR seed <staging>/copy` (darwin) or `cp -a --reflink=auto seed <staging>/copy` through
   `runBounded`, as today. **[r3]** The attempt publishes only after this, its own child, has
   exited 0 (`runBounded` waits for the child). Error: `what` = `seedWhat(err)`.
5. **[r2]** Only after `cp` returns 0, write the marker `<staging>/copy/.world-seed-complete` (§3.1a)
   with `os.WriteFile(…, 0o600)`. Error: `what` = `the completion marker cannot be written`.
6. `os.Rename(<staging>/copy, <cache>/<name>)`, then `RemoveAll(<staging>)`. Error: `what` = `the
   copy cannot be published`.

**[r3] Lifecycle ownership.** A staging directory belongs to exactly one attempt of one daemon
incarnation. Only that attempt writes the marker into it and only that attempt renames out of it.
A `cp` child that outlives a killed daemon (R4) keeps writing only into its own incarnation's
staging directory, which the new incarnation never publishes; the new incarnation's cleanup
removes it (bounded as in step 3), and if the old child is still writing, removal can fail or
leave residue, which is logged and retried at the next attempt. There is no one-copy disk bound
(R4, rewritten).

The staging path is a sibling of the project cache dir, not inside it. The sandbox may write
only the worktree and `<state>/exec-cache/<ep>/<project>` (`host/broker/exec_settings.go:104,115`),
so the agent cannot see or race the staging copy. Rename within one directory tree on one volume
is atomic, so `<cache>/<name>` appears with its marker in one step, and only after a copy that
returned 0.

#### 3.1a [r2] Completion evidence: the marker

- **What it is.** A regular file `.world-seed-complete` at the top of the published copy, whose
  bytes are exactly `world-exec-seed v1 <name>\n` (for example `world-exec-seed v1 GOCACHE\n`).
  `v1` is the format version. A later format bumps it, and an older binary treats a marker it does
  not recognise as unverified and re-seeds (fail-closed).
- **Verified** means: `Lstat(<cache>/<name>)` is a directory, `Lstat` of the marker is a regular
  file (not a symlink) of at most 64 bytes, and its bytes are exactly the expected line. A pure
  function `seedMarkerValid(name string, data []byte) bool` decides the bytes.
- **Why inside the copy.** The rename publishes the copy and its marker in one step, so there is no
  moment when `dst` has a marker that a crash could leave without its copy, or the reverse.
- **Crash behaviour.** Each window leaves a state the next attempt handles without trusting it:

  | Crash point | What is left | Next attempt |
  |---|---|---|
  | during `cp` (step 4) | partial staging, no `dst` | **[r3]** step 3 removes the abandoned staging directory (bounded) and copies into a new unique one |
  | after `cp` returned 0, before the marker (between 4 and 5) | full staging without a marker, no `dst` | the same: staging is never promoted, only re-made |
  | after the marker, before the rename (between 5 and 6) | full staging with a marker, no `dst` | the same |
  | after the rename | `dst` with a valid marker | verified; ready, no copy |
  | during step 2's `RemoveAll` of an unverified `dst` | part of an unverified `dst` | still unverified (it never had a valid marker); step 2 removes the rest |
  | power loss after the rename, before the marker's data reached the disk | `dst` with an empty or short marker | the bytes do not match; unverified; removed and copied again |

  Staging is never trusted, so no `fsync` is needed for safety. The power-loss row is reasoned
  from the rule; I did not test power loss (**UNMEASURED**).
- **Trust.** The marker is inside `<cache>/<name>`, which the sandbox can write
  (`exec_settings.go:115`). So the agent can delete, change or forge it, and it can change the
  cache's contents after publication. That only affects its own episode's cache, which it can
  already write freely. The seed is never written (AC2.8 at `handlers_exec_srt_test.go:125-127`),
  and other episodes' caches are outside its write paths. The marker proves one thing: World's copy
  step finished for this `dst`. It does not prove the cache is unchanged since then, and the design
  does not need that.
- **Go ignores it.** V19: a copy of the real seed with the marker added, used as `GOCACHE` for a
  `go build`, builds (rc 0) and the marker is still there afterwards.

`ExecHandlerConfig` gains `Log io.Writer` (nil = no log). Each finished attempt writes one line:
`ailang-worldd: exec cache for episode "<ep>" (<project>): GOCACHE seeded at bind in 0.38s` or
`... failed: <failure line>`.

### 3.2 The per-call path: a check, then join, never spawn on the call's clock

`prepareCache(ctx)` keeps its first loop (World's own `home`, `tmp`, `pycache`, `npm` dirs, and
unseeded caches made empty). Then, for the seeded caches:

1. **The check.** **[r2]** If every seeded `<cache>/<name>` is verified (§3.1a: a directory with
   a valid marker), return nil. A `dst` that exists without a valid marker is treated as missing
   and **[r3]** falls through to step 2, so the join cannot be skipped while a copy is running; a
   later attempt removes it before copying (§3.1 step 2). This is the only work the
   call does when bind succeeded: two `Lstat`s and one read of at most 64 bytes, and no
   subprocess.
2. **Join.** If a copy is running (started at bind or by an earlier call), wait for it, or for
   `ctx.Done()`, whichever is first. If the copy finishes, go back to 1. If the call's context ends
   first, refuse with `seed-cache: GOCACHE: still copying since bind (4.2s); retry`. The copy keeps
   running on its own clock; the next call finds it done.
3. **Re-attempt.** If no copy is running and `attempts < 3` (bind failed, or the agent deleted its
   GOCACHE or its marker), start the next attempt in the background, exactly as `StartSeed` does, and join it
   as in 2. The call waits within its own budget, but the copy is not bound by it.
4. **Give up.** If `attempts == 3` and the last one failed, refuse at once with the last line plus
   `; 3 of 3 attempts failed: fix the seed, then restart serve`. Nothing is spawned.
   **[r3]** If `attempts == 3` and the last attempt succeeded but the cache is again missing or
   unverified, refuse at once with the same give-up line (no spawn).

If `StartSeed` was never called (the test fixtures and any future caller), step 3 starts attempt 1
from the call. So there is one copy path, and it never runs on a call's clock.

Every refusal above is a `*HandlerFailure` with the row-255 class `seed-cache`, so it reaches the
effect record, the MCP result and `why` exactly as row 255 built (no change there). The class set
stays closed at 10.

### 3.3 Why re-attempt, and why three

The brief asks what the per-call path does after a bind failure. The options were:

- **Refuse forever** until restart. Fail-closed, but a copy that timed out under a transient load
  spike would disable `workspace-exec` for the whole episode.
- **Re-copy inside each call** (today). The row's complaint: the failure repeats on every call.
- **Re-attempt off the call's clock, capped** (chosen). A transient failure recovers without a
  restart. A deterministic one (seed missing, unreadable file) stops after three attempts with a
  line saying so. Each attempt is logged once and shown in the report (§3.6), so the retry is not
  silent. The cap is per handler, so at most 3 × 60 s of copying per episode × profile.

Three: one bind attempt plus two retries is enough to ride out a spike that lasts one copy, and
small enough that a permanent fault stops within the first few calls. No measurement justifies 3
over 2 or 5; it is a judgement (§6 R3).

### 3.4 The bound: 60 seconds

The copy no longer blocks any call except an exec call that needs the cache, and that call waits
only within its own budget (§3.2). So a longer bound costs only this: an exec call during a stuck
copy is refused "still copying" until the bound ends. A shorter bound brings back the row's
failure. 60 s is 59× the loaded median (1.014 s at load ~42) and 52× the loaded max (1.146 s), and
171× the idle median (V6–V9). Under the incident's load 110–130 the clone time is **UNMEASURED**;
if the clone scaled with CPU oversubscription as it did at 3× (V9), it would still be well under
60 s, but that is an inference, not a measurement. A full (non-reflink) copy on Linux is also
**UNMEASURED**. The bound is a constant, not a flag: no operator has asked to tune it.

### 3.5 When bind happens, and D-257-1

Today bind is the first dispatched call that wires the exec tool (§1 item 1, **[r2]** reworded
after V18). This design calls `StartSeed` right
after `NewExecHandler` succeeds in `workspaceTools.execHandler` (`workspace.go:496-506`). That alone
takes the copy off the exec call's clock and, when the agent's first call is not an exec (a read,
a check), starts the copy before the first exec. **[r2]** It does so only when that first call
wires the exec tool, which V18 shows is every dispatched call once the worktree resolves and a
profile applies.

To start it **before the agent runs**, the daemon needs a request from the operator side first.
No POST reaches the daemon between the mint and the agent's first call: `episode start` mints
into the credentials store directly and only reads `GET /v1/health`; `episode launch` reads
`GET /v1/session/self` with the session token just before it runs the agent
(`cmd/ailang-worldd/episode_launch.go:350`). The `/v1` route table is frozen
(`host/daemon/daemon.go:1054-1058`).

**D-257-1 — settled [r3] as A.** Round 2 (kimi-k3) found B's load-bearing premise unverified:
nothing in this doc shows that the freshly minted token is written to the session file and
resolvable by `d.resolver.ResolveContext` at the instant after `provisionEpisode` returns. The
reviewer's fix was "until that row exists, demote the recommendation from B to A", and r3 takes it.
The build therefore does **A**. B stays below as a recorded option for a later row, with its
missing premise named; it is not built.

- **A.** Bind only at the first dispatched call that wires the exec tool (**[r2]** wording, V18).
  No route or GET changes. The first exec call can
  still have to wait for a copy started by a read a moment earlier.
- **B (not built; [r3] premise unverified).** Also start the copy from an authenticated `GET /v1/session/self`. The
  handler resolves the token (`host/daemon/session_self.go:46-62`; **[r2]** V20 quotes `:49`
  `ResolveContext`, the 401 on `out.Denied` at `:58-60`, and `b := out.Success` at `:63`), then
  calls `d.workspace.warmExec(episode)`, which binds the exec handler through the same
  `episodeRoot` and `execHandler` path that `registry()` uses, and calls `StartSeed`. **[r2]** The
  attempt creates its own parent tree (§3.1 step 0), so a warm GET on a fresh episode, before any
  `prepareCache`, can seed (test D3). It is
  idempotent (a bound handler with a started or finished copy does nothing), it needs a valid
  session token, and it writes no store entry and spends no budget. `episode launch` already makes
  this GET before it runs the agent. `episode start` makes one too, right after the mint (both
  start paths: `cmd/ailang-worldd/episode.go` after `provisionEpisode` at `:343`, and
  `episode_preset.go` after `provisionEpisode` at `:301`, **[r2]** on the non-launch branch only,
  before `printLaunch` at `:311`; the `--launch` branch goes to `launchEpisode` at `:304`, whose
  launch path already makes this GET, V22). It is bounded at 3 s, read back from the session
  file as launch does (`episode_launch.go:330`). **[r2]** The 3 s bound is a new constant
  `episodeWarmTimeout = 3 * time.Second` in `cmd/ailang-worldd/episode.go`, used as the
  `http.Client.Timeout` of the one GET, next to `getDaemonHealth`'s 5 s client at `:396` (V23).
  3 s is below the health GET's 5 s because the warm GET only starts the copy and returns. An unreachable daemon or a refused GET prints one
  `exec cache: not started (<why>); the first call will start it` line and does not fail the
  start.
- **C.** A new `POST /v1/episodes/{episode}/prepare`. Cleanest HTTP semantics, but it changes the
  frozen table, which needs a ledger decision of its own.

**[r3]** With A, `warmExec`, the `session_self.go` change, the `episode start` warm GET and
`episodeWarmTimeout` are all dropped from the build.

### 3.6 Where a bind failure is visible

- **The exec call** (unchanged carrier): the `seed-cache` line in the record, the MCP result and
  `why` (§3.2).
- **The daemon log**: one line per attempt (§3.1).
- **`episode report`**: `EpisodeReport` gains
  `ExecCache []ReportExecCache json:"exec_cache,omitempty"`, one entry per seeded cache of the
  episode's bound handler: `{"project","cache","state","attempts","by","took_ms","failure"}`,
  with `state` one of `copying`, `ready`, `failed`, `gave-up`. `handleEpisodeReport`
  (`host/daemon/episode_report.go:692`; **[r2]** V21 quotes `:692`, `readEpisodeStart` at `:709`,
  the NotFound at `:721` and `writeJSON` at `:724`) fills it from `d.workspace.execCacheStatus(episode)`. It is
  live daemon state, not store state: the report's other fields come from the store, and this one
  says so in its doc comment and in the text renderer. After a daemon restart it is empty until the
  next bind. The text renderer (`cmd/ailang-worldd/episode_report.go`) prints one line per entry,
  e.g. `exec cache   ailang-world GOCACHE failed at bind (attempt 1 of 3): seed-cache: GOCACHE: timed out after 1m0s`.
  The NotFound guard (`:720-721`) also passes when `ExecCache` is non-empty.
- **[r3]** No session-self carrier: D-257-1 = A.

### 3.7 Pure core

`execSeedNext(state, allVerified, now) action` decides steps 1–4 of §3.2 as a pure function of
the state and two inputs (no I/O, no clock read); `seedWhat`, the detail builder and **[r2]**
`seedMarkerValid` and `seedMarkerBytes(name)` are pure. They are
tested with no process run. This is Go host code at the effect boundary (S2); no `.ail` is touched,
so S1's Z3 contracts do not apply, and no package can hold it, because the copy is a host
subprocess (S3).

## 4. By-construction argument

- **No call runs against a missing cache.** `Execute` reaches `runBounded` for srt only after
  `prepareCache` returns nil (`handlers_exec.go:308` before `:329`). `prepareCache` returns nil only
  from §3.2 step 1, which requires every seeded `<cache>/<name>` to be verified.
- **No call runs against a partial copy World made.** **[r2]** Step 1 accepts a `dst` only with a
  valid marker (§3.1a). World writes that marker only in §3.1 step 5, after `cp` returned 0, into
  staging, and publishes it only by the rename in step 6. A killed copy leaves only staging, which
  the sandbox cannot reach and the next attempt removes first. A `dst` from the code before this
  change has no marker, because that code never wrote one, and step 1 of §3.1 refuses a seed that
  holds one, so that `dst` cannot carry a copied marker into readiness. It is unverified, so it is
  removed and copied again, or, if removal fails, the attempt fails with `seed-cache`. Every
  crash window ends in a state the next attempt does not trust (§3.1a table).
- **[r2] The bound tree exists before the copy writes into it.** §3.1 step 0 creates
  `<state>/exec-cache/<ep>/<project>` inside the attempt, so the attempt needs no earlier
  `prepareCache`.
- **No copy runs on a call's clock.** The only `runBounded` call for `/bin/cp` is in the attempt
  goroutine, under `h.seedBudget` from a background root. A call's context is used only to stop
  *waiting*.
- **Retries are bounded and visible.** `attempts` only grows, under `mu`; step 3 needs
  `attempts < 3`; each attempt logs once and sets the report state.
- **One copy at a time.** `running` is set and cleared under `mu`; a second caller joins it.
  `seedMu` is replaced by `execSeed.mu`, which is never held while waiting.
- **[r3] No two copies share a staging tree.** Each attempt's staging directory is created with an
  exclusive `os.Mkdir` under a name unique to its daemon incarnation and attempt (§3.1 step 3), and
  only that attempt writes its marker and renames out of it. A surviving `cp` from a dead daemon
  can write only into its own incarnation's directory, which nobody publishes.
- **[r3] The episode→handler cache is safe to share.** It is reached from two goroutine sources:
  dispatch (`registry()` → `execHandler`, and `Dispatch` holds no coordinator-wide lock, so
  dispatches can overlap) and the report fill (`execCacheStatus`). The map is already guarded:
  `execMu` "guards execH only" (`workspace.go:123-126`, V24), and `execHandler` takes it at
  `:490-491`. All access goes through one mutex-protected accessor
  `(*workspaceTools).boundExecHandler(episode)`, which `execHandler`'s cache lookup and
  `execCacheStatus` both use; `execCacheStatus` copies the handler pointer under `execMu` and
  reads the seed state under `execSeed.mu`, never holding both.

## 5. Test plan and mutation table

All in `host/broker/exec_seed_test.go` unless named; they use the existing `newExecFixture` with a
`caches.GOCACHE.seed` profile (as `failure_test.go:155-218` does). A test seam
`h.seedCommand func(seed, dst string) (string, []string)` lets two tests run a slow copy
(`/bin/sh -c 'sleep 0.5; exec /bin/cp -cR "$1" "$2"'`).

| Id | Test | What it asserts | Mutant it kills |
|---|---|---|---|
| **S1** | `TestSeedAtBindFirstCallSpawnsNoCopy` | `StartSeed`, wait for the copy; record `Spawns()`; one `Execute` adds exactly 1 (srt); `<cache>/GOCACHE` holds the seed's files; no staging path is left | `StartSeed` a no-op (the call copies: delta 2) |
| **S2** | `TestSeedAtBindFailureReachesTheCall` | missing seed: after `StartSeed`, the state is `failed`; `Invoke` through a live broker session gives `seed-cache: GOCACHE: the copy exited 1 (attempt 2 of 3)` in the record and `EffectFailedError.Reason`; srt never spawned | per-call path skips the check and runs srt |
| **S3** | `TestSeedAttemptsAreCapped` | missing seed: calls 1 and 2 start attempts 2 and 3; call 3 spawns nothing and its line ends `3 of 3 attempts failed: fix the seed, then restart serve` | cap removed (call 3 spawns a 4th copy) |
| **S4** | `TestSeedPartialCopyIsNeverAccepted` | a seed with one mode-000 file: the copy exits 1, `<cache>/GOCACHE` does not exist, staging is removed; and a planted staging dir from a "crash" is removed and replaced by the next attempt | copy straight into `dst` with no rename and no `RemoveAll` (call 2 accepts the partial `dst`) |
| **S5** | `TestSeedOutlivesTheCallThatWaits` | slow seam; an `Execute` with a 100 ms context refuses `seed-cache: GOCACHE: still copying since bind (…)`; the copy then finishes; the next `Execute` adds 1 spawn (srt only) | the copy bound to the caller's context (it is killed; call 2 copies again) |
| **S6** | `TestSeedHasItsOwnBound` | `h.seedBudget = 1ms`: the state is `failed` with `timed out after 1ms (attempt 1 of 3)` | bound read from a caller's context (no deadline: 30 s default) |
| **S7** | `TestSeedNextIsPure` | table over `execSeedNext`: all four steps and the cap edge | any branch swapped |
| **S8** [r2] | `TestStartSeedReusesOnlyAVerifiedCache` | a `<cache>/GOCACHE` published by a previous attempt (directory plus valid marker; a restart): `StartSeed` spawns nothing and the state is `ready`, and an `Execute` adds 1 spawn (srt). Second arm: the same directory with a marker of the wrong version (`world-exec-seed v0 GOCACHE`) is removed and copied again (1 copy) | `StartSeed` always copies (arm 1); marker bytes not checked (arm 2) |
| **S9** [r2] | `TestLegacyPartialCacheIsRemovedAndReseeded` | a marker-less partial `<cache>/GOCACHE` holding one stray file (the old code's crash artifact): `StartSeed` removes it, the copy is published with a marker, the stray file is gone; and, in a second fixture with `StartSeed` never called, the first `Execute` succeeds with a `Spawns()` delta of 2 (one copy, then srt), after which the stray file is gone and the marker is present | readiness by `Lstat` alone (delta 1: srt runs against the planted dir, and the stray file survives) |
| **S10** [r2] | `TestSeedCrashWindowsAreNeverTrusted` | three planted states with no `dst`: partial staging; full staging without a marker; full staging with a marker. Each: `StartSeed` removes staging and publishes a fresh copy. A fourth: `dst` with an empty marker file (power loss) is removed and copied again. A fifth: an unverified `dst` that cannot be removed (its parent made 0o500) gives `seed-cache: GOCACHE: an unverified copy cannot be removed (attempt 1 of 3)` and the state is `failed`, never `ready` | staging promoted when it holds a marker; a removal error treated as ready |
| **S11** [r2] | `TestSeedHoldingTheMarkerIsRefused` | a seed with a top-level `.world-seed-complete`: the attempt fails with `the seed holds World's marker; remove it from the seed`, and `<cache>/GOCACHE` does not exist | step 1 removed (the copy is published with the operator's marker) |
| **S12** [r2] | `TestSeedMarkerValidIsPure` | table over `seedMarkerValid`: the exact line; a trailing byte; v0 and v2; another cache's name; empty; 65 bytes | any comparison loosened |
| **D1** | `TestEpisodeReportShowsBindSeedFailure` (`host/daemon/exec_seed_test.go`, rig of `health_exec_tool_status_test.go:101`) | a bound episode whose seed is missing: `GET /v1/episodes/<ep>/report` has `exec_cache[0]` with `state: failed` and the line, with no committed call | report omits the field |
| **D3** [r2, r3] | `TestBindSeedsAFreshEpisode` (same file) | a fresh episode worktree and an empty state dir with no `exec-cache` at all; the bind path (`workspaceTools.execHandler`, as `registry()` calls it) with no prior `prepareCache` and no fixture-made dirs: the state becomes `ready` and `<state>/exec-cache/<ep>/<project>/GOCACHE` holds the seed's files and the marker | §3.1 step 0 removed (the copy fails with ENOENT, state `failed`) |
| **D4** [r3] | `TestEpisodeReportRacesFirstDispatch` (same file; run under `-race`) | a `GET /v1/episodes/<ep>/report` issued concurrently with the episode's first dispatched call: no panic, no race report, and a well-formed `exec_cache` view (empty or one entry) | `execCacheStatus` reading `execH` without `execMu` |
| **S13** [r3] | `TestSeedRestartIgnoresASurvivingOldCopy` (`host/broker/exec_seed_test.go`) | a controlled "old copy" process (a shell loop) keeps writing files into an earlier incarnation's staging directory while a new handler seeds: the new attempt's staging directory and the published `<cache>/GOCACHE` never contain the old process's files; the published copy is verified; after the old process stops, the next attempt's cleanup removes the abandoned directory, and no more than 8 abandoned directories are removed per attempt | a fixed staging path shared across incarnations (the old process's files reach the published cache); cleanup unbounded or absent |
| **C1** | `TestEpisodeReportPrintsExecCache` (`cmd/ailang-worldd/episode_report_test.go`) | the text renderer prints the exec cache line; with no entries it prints nothing new | renderer drops the line |
| edited | `TestExecFailureClassesReachTheRecord` (`failure_test.go:208-218`) | "seed-cache exited" expects ` (attempt 1 of 3)`; "seed-cache timed out" uses `seedBudget = 1ms` instead of `expiredCtx` (an expired call context now means "still copying", not a copy timeout) | — |
| edited | `TestProductionContextRoots` (`host/store/context_roots_test.go`) | +1 allowed root for the seed goroutine | — |

## 6. Residuals and risks

- **R1** The incident's load 110–130 is **UNMEASURED**. If a 226 MB clone ever needs more than
  60 s, all three attempts time out and `workspace-exec` is refused for the episode with a line that
  says why. That is fail-closed, not silent.
- **R2** Linux (`cp -a --reflink=auto` on a filesystem without reflinks is a full copy) is
  **UNMEASURED**.
- **R3** The cap of 3 is a judgement (§3.3).
- **R4 [r3]** A daemon killed during a copy can leave the `cp` child running until it ends. It
  writes only into its own incarnation's staging directory, which is never published (§3.1 step 3,
  Lifecycle ownership). Abandoned staging directories are removed at the next attempt's start, at
  most 8 per attempt; a removal that fails while an old child is still writing is logged and
  retried at a later attempt. So disk use is one live staging copy per episode × profile plus any
  abandoned directories not yet removed; this design does **not** claim a one-copy bound.
- **R5** The report's `exec_cache` is live state and is empty after a restart until the next bind
  (§3.6). Making it durable would need a store entry, which is out of scope.
- **R6 [r3]** None: D-257-1 = A, so no GET gains a side effect. The cost is that the copy starts at
  the agent's first dispatched call, not before the agent runs.
- **R7** An agent that deletes its GOCACHE **[r2]** or only its marker spends an attempt: the next
  exec call finds the cache unverified, and the attempt removes it and copies again. That is
  intended, and the cap still holds. An agent that forges the marker, or empties the cache but
  keeps the marker, gets the cache it made. That only affects its own episode (§3.1a, Trust).
- **R8 [r2]** The marker sits in the published cache, so a tool that lists the cache dir sees one
  extra dotfile. `go build` ignores it (V19). Other caches (npm, pycache) are not seeded today; a
  seeded one would need the same check before it ships.
- **R9 [r2]** An operator seed that holds `.world-seed-complete` at its top level is refused
  (§3.1 step 1) and the exec tool stays refused until the seed is fixed.

## 7. Milestones

- **M1 (broker).** `exec_seed.go` (state, `StartSeed`, attempt, `execSeedNext`, status), the
  `prepareCache` rewrite, `ExecHandlerConfig.Log`, `execSeedBudget` removed, **[r2]** the marker
  (`seedMarkerBytes`, `seedMarkerValid`) and steps 0–2 of §3.1. Tests S1–S12 and the two edited
  tests.
  **[r3]** Plus S13 and the unique staging directory of §3.1 step 3.
  `go test ./host/broker/ -run 'Seed|LegacyPartialCache|TestExecFailureClassesReachTheRecord' -count=1` and
  `go test ./host/store/ -run TestProductionContextRoots -count=1`.
- **M2 (daemon and CLI).** **[r3]** The `workspace.go` hunk (which also routes the cache lookup
  through `boundExecHandler`), `host/daemon/exec_seed.go` (`boundExecHandler`, `execCacheStatus`),
  `exec_cache` on the report, the report text line. Tests D1, D3, D4, C1.
  `go test ./host/daemon/ ./cmd/ailang-worldd/ -run 'ExecCache|BindSeed|ReportRacesFirstDispatch' -count=1`
  and `go test -race ./host/daemon/ -run 'ReportRacesFirstDispatch' -count=1`.
- Then the controller's full `./scripts/verify_go.sh` and `./scripts/verify_ail.sh` (unchanged
  `.ail`).

## 8. Files

| File | Change | In #515? |
|---|---|---|
| `host/broker/exec_seed.go` | new: the seed state machine, **[r2]** the marker and steps 0–2 | no |
| `host/broker/handlers_exec.go` | `ExecHandlerConfig.Log`, `seed`/`seedBudget`/`seedCommand` fields, `prepareCache` rewritten, `execSeedBudget` removed | no |
| `host/broker/exec_seed_test.go`, `host/broker/failure_test.go` | new tests; two TX rows edited | no |
| `host/store/context_roots_test.go` | +1 allowlist row | no |
| `host/daemon/workspace.go` | **one hunk** in `execHandler` (`:489-506`): **[r3]** cache lookup through `boundExecHandler`, pass `Log: w.errLog` and call `eh.StartSeed()` | **yes**, but #515's hunks are at `:217` and `:266-284`, 200+ lines away; a rebase over #515 applies cleanly. If it ever conflicts, keep both: #515's hunks change exclusion logging, this one the handler construction. |
| `host/daemon/exec_seed.go` | new: **[r3]** `boundExecHandler`, `execCacheStatus` (methods on `*workspaceTools`, kept out of `workspace.go` on purpose) | no |
| `host/daemon/episode_report.go` | `ExecCache` field and fill, NotFound guard | no |
| `host/daemon/exec_seed_test.go` | D1, **[r2]** D3, **[r3]** D4 | no |
| `cmd/ailang-worldd/episode_report.go` (+ test) | text line | no |

Not touched: **[r3]** `host/daemon/session_self.go`, `cmd/ailang-worldd/episode.go`,
`cmd/ailang-worldd/episode_preset.go` (D-257-1 = A), `host/broker/exec_probe.go`, `host/broker/exec_startup.go`, `host/daemon/daemon.go`,
`cmd/ailang-worldd/doctor*.go`, `docs/QUICKSTART.md` (no new operator command or flag; the report
line explains itself).

## 9. What this is NOT doing

- Not changing the startup probe, profile admission, the registry bootstrap or any startup
  deadline (row 256).
- Not adding a failure class, a record field or a `/v1` route (C is named, not chosen).
- Not making the seed copy faster, and not changing how the operator builds the seed.
- Not making the report's `exec_cache` durable (R5).
- **[r2]** Not checking a published cache's contents against the seed. The marker proves World's
  copy step finished; what the agent does to its own cache after that is its business (§3.1a).
- Not touching `.ail` code, `tools/launchd/*` or the shared skills.

## 10. Verification Log

Commands ran in the worktree unless they say otherwise. Line numbers are at `ad13d835`.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Base and instruments | `git log --oneline -1`; `shasum -a 256 host/broker/handlers_exec.go host/daemon/{workspace,session_self,episode_report}.go \| cut -c1-16`; `go version`; `sysctl -n hw.ncpu`; `$HOME/.pinned-ailang/ailang --version` | `ad13d835 Merge pull request #520 …`; `996c5ee11daed4b9`, `1a178f8bc1e39f89`, `6a98846db138b6d3`, `64128681d7bffa13`; `go1.26.6 darwin/arm64`; `16`; `AILANG v0.41.0` |
| V2 | Bind is the first effectful call | `grep -n` over `coordinator.go`, `daemon.go`, `workspace.go`, `effectful.go` | `coordinator.go:316-317 bound, err := transitionreg.Bind(… c.cfg.Binder(call.EpisodeID, grants))`; `daemon.go:917 reg := d.workspace.registry(episodeID)   // a fresh map per call`; `workspace.go:468 func (w *workspaceTools) execHandler`, `:496 eh, err := broker.NewExecHandler(…`; `effectful.go:311 hctx, cancelHandler := context.WithTimeout(ctx, handlerBudget(ctx))`, `:47 HandlerCap = 10 * time.Second` |
| V3 | The copy runs on the call's clock; the check is `Lstat` | `cat -n host/broker/handlers_exec.go` | `:308 if err := h.prepareCache(ctx); err != nil {`; `:538 if _, err := os.Lstat(dst); err == nil {`; `:548 args := []string{"-cR", seed, dst} // APFS clonefile`; `:553 runBounded(ctx, handlerBounds{execTimeout: execSeedBudget(ctx)}, …`; `:556 _ = os.RemoveAll(dst)`; `:565-573 execSeedBudget` returns `time.Until(dl)` |
| V4 | The `ailang-world` profile and seed | `python3 -I -c '…json.load(…ailang-world.json)…'`; `du -sm $SEED`; `find $SEED -type f \| wc -l`; `df -h $SEED /tmp` | `timeout_ms 9000`; `caches {'GOCACHE': {'seed': '/Users/voightkampff/.ailang-world-exec/seeds/ailang-world/go-build'}}`; `226`; `2569`; both on `/dev/disk3s5` |
| V5 | `cp -cR` leaves a partial copy on error | seed `f1 f2 f3`, `chmod 000 seed/f2`, `/bin/cp -cR seed dst; echo rc=$?; ls dst` | `cp: seed/f2: Permission denied`, `rc=1`, `f1 f3` |
| V6 | Idle clone | `/tmp/iter372-designer/clone_n.sh idle 15 $SEED /tmp/iter372-designer` (timer: `python3 time.time()` before and after `/bin/cp -cR`) | `idle: N=15 min=0.345s median=0.351s max=0.472s rc_nonzero=0 load1 min=3.11 median=3.12 max=3.12` |
| V7 | Clone during `go test -p 8`, pinned | `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/... -count=1 -p 8 &`, one clone every ~6 s while it ran | `loaded2: N=40 min=0.346s median=0.385s max=0.482s rc_nonzero=0 load1 min=3.99 median=4.505 max=8.34`; `elapsed=309s`; `gotest rc=0`; summed `ps -A -o %cpu` 38–685 |
| V8 | Same, first 80 s, AILANG_BIN unset (supplementary) | `go test ./host/... -count=1 -p 8 &`, 15 clones 3 s apart from t=20 s | `loaded: N=15 min=0.348s median=0.379s max=0.628s rc_nonzero=0 load1 min=2.95 median=3.35 max=5.48`; `gotest rc=1` (pinned-interpreter tests failed without `AILANG_BIN`; not a load property) |
| V9 | Clone under CPU oversubscription | 48 × `/usr/bin/yes > /dev/null`, 40 s warm-up, 12 clones 2 s apart, spinners killed (`pgrep -x yes \| wc -l` → `0`) | `cpu48: N=12 min=0.686s median=1.014s max=1.146s rc_nonzero=0 load1 min=29.91 median=41.89 max=56.57` |
| V10 | The probe never seeds | `sed -n 430,475p host/broker/exec_probe.go` | `:452 h, err := NewExecHandler(…`; `:466 dirs = append(dirs, filepath.Join(h.cache, name)) // never seeded by the probe` |
| V11 | The sandbox writes only the worktree and the project cache | `grep -n Cache host/broker/exec_settings.go` | `:27 Cache string // <state>/exec-cache/<ep>/<project>: writable, readable`; `:104 allowRead := dedupe(append([]string{spec.Worktree, spec.Cache}, …`; `:115 AllowWrite: []string{spec.Worktree, spec.Cache},` |
| V12 | No POST reaches the daemon before the agent; launch GETs session self; the table is frozen | `sed -n 300,356p cmd/ailang-worldd/episode.go`; `sed -n 330,352p cmd/ailang-worldd/episode_launch.go`; `sed -n 1050,1104p host/daemon/daemon.go` | `episode.go:314 health, err := fetchDaemonHealth(hostPort)`, `:343 provisionEpisode(…)`, `:348 printLaunch(…)`; `episode_launch.go:350 c.doRequest(ctx, http.MethodGet, "/v1/session/self", nil, …"Bearer " + token)`; `daemon.go: "The thirteen /v1 patterns below are the complete frozen v1 machine table"` |
| V13 | #515's hunks in shared files | `git fetch origin sprint/world-iter362-row256`; `git diff --stat ad13d835...FETCH_HEAD -- host cmd docs`; `git diff ad13d835...FETCH_HEAD -- host/daemon/workspace.go` | 11 files (`doctor.go`, `doctor_hosts.go`, `QUICKSTART.md`, `exec_probe.go`, `exec_startup.go`, `daemon.go`, `workspace.go` +3 tests); `workspace.go` hunks `@@ -217,7` and `@@ -266,11 +266,21` only |
| V14 | The row-255 seed tests this changes | `sed -n 140,260p host/broker/failure_test.go` | `:208 {name: "seed-cache exited", …` → `"seed-cache: GOCACHE: the copy exited 1"`; `:211 {name: "seed-cache timed out", …` with `expiredCtx`, → `"seed-cache: GOCACHE: timed out after 1ms"` |
| V15 | Context-root precedent | `sed -n 60,85p host/store/context_roots_test.go` | `:80 "host/broker/handlers_exec.go|NewExecHandler|Background": 1,` |
| V16 | Scratch removed, tree clean | `ls /tmp/iter372-designer/`; `git status --short` | only scripts and text outputs, no `dst-*`; only this doc listed (re-checked at r2: same) |
| V17 [r2] | Where the project cache tree is created today | `grep -rn 'exec-cache' --include='*.go' host cmd \| grep -v _test`; `grep -n MkdirAll host/broker/handlers_exec.go`; `grep -n 'MkdirAll\|exec-cache' cmd/ailang-worldd/session_new.go`; `grep -rn 'func provisionEpisode' cmd/ailang-worldd/*.go` | `exec-cache` appears in production code only as `handlers_exec.go:178` (the path join), `exec_startup.go:156` and `exec_tools.go:131` (forbidden roots), `exec_probe.go:460` (probe cleanup), `exec_settings.go:27` (comment), `failure.go:31` (class name). `MkdirAll` in `handlers_exec.go`: `:189` (`filepath.Dir(h.settingsPath)`, the `exec/` dir), `:532` (`filepath.Join(h.cache, d)` inside `prepareCache`), `:543` (`dst`, unseeded caches, inside `prepareCache`). `session_new.go`: only `:492 os.MkdirAll(filepath.Dir(r.outPath), 0o700)` (the session file); `provisionEpisode` is `session_new.go:396`. So only `prepareCache` creates `<state>/exec-cache/<ep>/<project>` |
| V18 [r2] | What triggers the exec handler's construction | `sed -n 270,320p host/coordinator/coordinator.go`; `grep -rn 'execHandler(\|\.registry(' host/daemon/*.go \| grep -v _test`; `sed -n 425,452p host/daemon/workspace.go` | `coordinator.go`: `effectful := false` … `if effectful { … seedGrants … }`, then, outside that block, `:316-317 bound, err := transitionreg.Bind(… c.cfg.Binder(call.EpisodeID, grants))`; `daemon.go:917 reg := d.workspace.registry(episodeID)`; `workspace.go:429-436` `if !w.enabled() { return broker.Registry{} }` and `epRoot, ok := w.episodeRoot(episodeID); if !ok { return broker.Registry{} }`; `:450 reg[broker.EffectWorkspaceExec] = w.execHandler(episodeID, epRoot)`. No other caller of `execHandler(` |
| V19 [r2] | A marker inside GOCACHE does not disturb Go | copy the real seed to `/tmp/iter372-designer/v19/gc`; `printf 'world-exec-seed v1 GOCACHE\n' > gc/.world-seed-complete`; a one-file module; `GOCACHE=…/gc go build`; run it; `cat` the marker; then `rm -rf` the scratch dir; `ls $SEED/.world-seed-complete` | `build rc=0`; `hi`; `world-exec-seed v1 GOCACHE`; `ls -a gc` shows `.world-seed-complete`, `00`, `01`; the real seed has no marker (`No such file or directory`) |
| V20 [r2] | Session self resolves the token before anything else | `sed -n 40,70p host/daemon/session_self.go` | `:46 func (d *Daemon) handleSessionSelf(…)`; `:49 out, err := d.resolver.ResolveContext(ctx, r.Header.Get("Authorization"), time.Now().Unix())`; `:58 if out.Denied != nil {`; `:60 writeAPIError(…, http.StatusUnauthorized)`; `:63 b := out.Success`; `:64 rep, known, err := buildEpisodeReport(ctx, d.intents, d.reads, b.EpisodeID)` |
| V21 [r2] | The report handler and its NotFound guard | `sed -n 685,725p host/daemon/episode_report.go` | `:692 func (d *Daemon) handleEpisodeReport(…)`; `:709 start, startErr := d.readEpisodeStart(ctx, episode)`; `:720 if !ok && start == nil && rep.Finish == nil {`; `:721 writeAPIError(w, "NotFound", "episode "+episode+" has no committed calls", …)`; `:724 writeJSON(w, http.StatusOK, rep)` |
| V22 [r2] | The preset start path's branches | `sed -n 290,312p cmd/ailang-worldd/episode_preset.go` | `:301 if code := provisionEpisode(req, out, errOut, env); code != exitOK {`; `:304 return launchEpisode(ep, name, dir, headless, out, errOut, …)` inside `if launching {`; `:311 printLaunch(out, launchInfo{…})` |
| V23 [r2] | Where the 3 s bound would live | `sed -n 394,400p cmd/ailang-worldd/episode.go`; `grep -n 'episodeWarmTimeout' cmd/ailang-worldd/*.go` | `:395 func getDaemonHealth(addr string) …`, `:396 client := &http.Client{Timeout: 5 * time.Second,`; `episodeWarmTimeout` does not exist yet (no match): the build adds it to `episode.go` |

| V24 [r3] | The episode→handler map is mutex-guarded; dispatch is not serialized (controller-measured) | `grep -n "execMu\|execH " host/daemon/workspace.go`; `grep -n "sync.Mutex\|\.Lock()" host/coordinator/coordinator.go`; `grep -n "func (c \*Coordinator) Dispatch" host/coordinator/*.go` | `workspace.go:123 // execMu guards execH only. It is never held with mu, so exec`, `:125 execMu sync.Mutex`, `:126 execH  map[string]episodeTool // episode id -> constructed exec handler`, `:490 w.execMu.Lock()`, `:491 defer w.execMu.Unlock()`, `:504-505` lazy map init; `coordinator.go`: no `sync.Mutex` or `.Lock()` match; `coordinator.go:267 func (c *Coordinator) Dispatch(…)` |
| V25 [r3] | The rebase onto `dc907675` changed only the charter (controller-measured) | `git diff --stat ad13d835 dc907675`; `shasum -a 256 host/broker/handlers_exec.go host/daemon/{workspace,session_self,episode_report}.go \| cut -c1-16` at `dc907675` | `design_docs/world-mission.md \| 9 +++++----`, `1 file changed`; `996c5ee11daed4b9`, `1a178f8bc1e39f89`, `6a98846db138b6d3`, `64128681d7bffa13` (identical to V1) |

## 11. Quorum log

- **Round 1** (iteration 372, artifact `mission-quorum-world/w-exec-cache-seed-at-bind-2026-11-09T17-13-15Z.json`,
  author `claude:claude-opus-5-5`, Anthropic seat benched; seated `oc-glm-5-3`, `oc-kimi-k3`,
  `gpt6-1-sol`; $0.294): **BLOCKED 3/3.** Surfaces: readiness-by-existence accepts a legacy partial
  `dst` (glm, gpt); the attempt never creates its parent tree (kimi). Revised by the designer → r2.
- **Round 2** (artifact `…-2026-11-09T17-26-20Z.json`, same seats, $0.337): **BLOCKED 3/3.**
  Surfaces (all different from round 1): concurrent access to the episode→handler map from the new
  report reader (glm); option B's token-resolvable premise unverified (kimi); a fixed staging path
  shared with a `cp` that outlives a dead daemon (gpt). No objection disputes the design direction,
  and each carries a concrete reviewer-authored `proposed_fix`.
- **r3 = controller narrow-refinement carve-out** (Gate 2; first use ratified at iteration 95):
  applied the three round-2 fixes. glm: V24 (map guard, dispatch not serialized), the §4 accessor
  bullet, test D4 under `-race`, V25 (rebase diff and shas). kimi: B demoted to A ("until that row
  exists, demote the recommendation from B to A"), the §3.2 cap-edge sentence, and step 1 reworded
  to fall through to step 2. gpt: exclusively created unique staging per incarnation and attempt,
  publish only after its own child exits 0, lifecycle ownership and bounded cleanup (8 per
  attempt), the one-copy disk claim withdrawn (R4), restart test S13. The cleanup count of 8 is the
  controller's choice of a bound the reviewer required but did not number. No re-quorum (the
  protocol allows one); the sprint judge checks the build against this r3.
