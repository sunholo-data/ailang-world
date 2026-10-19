# w-workspace-exec-long-run-lane — long `workspace-exec` runs as a recorded start / collect pair, on a snapshot, inside the frozen D7 deadlines

**Status**: DRAFT for attended review (2026-10-19). Queue row 145. No quorum yet. Every OPEN decision (§8) carries a recommended default.
- **Direction:**
  - Row 145 was re-tagged an evidence row, attended on 2026-10-19 (Mark: "is our 9 second timeout too harsh?" → "yeah sounds good"). Its cost is measured as friction #15 (PR #311).
  - D-140-4 = A (D-WORLD-62, 2026-10-03) named this lane as the follow-up for runs longer than the handler cap.
- **Extends:**
  - row 140, `design_docs/planned/w-workspace-exec-toolchain-effect.md`: the §4.3 profile, §4.4 settings, §4.5 handler and §4.7 replay;
  - row 184, `w-toolchain-worktree-delta.md` (`worktree_delta`);
  - row 156, `w-exec-protected-test-paths.md` (the protected-path proof; IN-SPRINT, M2 on `dev`).
- **Base:** ailang-world `origin/dev` `78f2860`.
- **Rig:** macOS 26.6.2 arm64, node v26.0.0, srt **0.0.78** (the row-140 pin), tool binary AILANG v0.54.15 (`~/.pinned-ailang-tools/v0.54.15/ailang`). **Linux is not measured.**

**Clauses**:
- 3: a long run is a brokered, budgeted, recorded effect like any other;
- 4: the floor's World arm and the dogfood corpus need full suites;
- the north star: the mission loop runs its own gate through World (row 196).

**Estimate**: ≈ 2 executor days (M0–M5), each milestone ≤ ~30 min of executor work and ≤ ~250 LOC.

**Probe conditions**:
- Every V-row was run first-party on 2026-10-19, or read at the cited `file:line` on `78f2860`.
- Probe scripts and outputs are banked in `design_docs/verification/world-row145-design-2026-10-19/` (`q*.sh` / `q*.py` and their `*.out.txt`).
- The running trial daemon (127.0.0.1:47330) was only reached by an unauthenticated `GET /v1/health` from inside srt (Q1). Nothing was stopped or written.
- Episode worktrees were only read; the Q4 clone went to a scratch dir that was then removed.
- No full test suite was run for this doc. Durations come from the trial graders' own logs.
- This doc claims no `.ail` text.

## 1. Why (measured)

- **P1 — World's own gate does not fit in 9 s, by more than an order of magnitude.**
  - In episode w190, the unsandboxed graders took:
    - `go test ./...`: **244.0 s**;
    - the v0.54.15 tool gate (`-p 1`, 5 packages): **518.5 s**;
    - `verify_ail.sh`: **25.9 s**;
    - `go vet`: 1.4 s (V1).
  - w139t: 243.8 s, 486.7 s and 25.1 s respectively (V1).
  - 14 of 29 Go packages each take more than 9 s on their own, from `host/daemon` at 242.6 s down to `host/store` at 13.9 s (V2).
  - The row text says `verify_ail.sh` takes "minutes". The measurement is 25–26 s, but that is still 3× the cap.
- **P2 — Other projects' "fast" suites do not fit either.** platform-source's `test-fast` ran 3,913 tests in 46.1 s (`real 51.6 s`) (V3).
- **P3 — The cap shapes agent behaviour, and its cost shows up as a second round.**
  - w190 made 91 `workspace-exec` calls, each a narrow `-run` filter. 4 of them still timed out (V4).
  - It listed as "could not run here" exactly the tests behind its two red graders:
    - the verifygate gate arms (>9 s);
    - `TestSeToolsManifestPublishes`;
    - the cliwalk/demo e2e tests.
  - Fixing those took a second World round (friction #15, V5).
- **P4 — Half of what w190 could not run is not about time.**
  - The cliwalk/demo and daemon MCP e2e tests bind a loopback port, which srt denies under World's settings. 3 of w190's 91 exec results show `bind: operation not permitted` (V4, V6).
  - A long-run lane alone does not make them runnable (§4.9, D-145-7).
- **P5 — The frozen deadline chain leaves no room for a synchronous fix.**
  - The deadlines are:
    - `HandlerCap` 10 s with 6 s headroom (`host/coordinator/effectful.go:47-48`);
    - `invokeDeadline` 20 s and `writeTimeout` 30 s (`host/daemon/daemon.go:99-100`);
    - `ExecMaxTimeoutMS = HandlerCap − 1 s` (`host/broker/exec_profile.go:30-33`).
  - Raising them was option B of D-140-4 and was rejected. A 244 s suite would not fit even a 30 s write timeout.
  - So the run must outlive the call that starts it.

## 2. Goals and non-goals

**Goals**
- G1: an agent can run a profile command for up to an operator-set wall clock (minutes, not seconds) without any World call exceeding the D7 deadlines.
- G2: a long run is as confined as a short one: the same srt settings shape, no network, the read fence, the deny names, row 156's protected paths.
- G3: it is as recorded as a short one. The start and every collect commit entries; the final result carries the row-140 `execResult` fields plus `worktree_delta` and the protected-path proof. Replay never re-executes.
- G4: it is bounded:
  - per-run wall clock (operator);
  - per-episode concurrency;
  - a process-wide pool that does not starve short calls;
  - output head and tail, as now;
  - a grant budget.
- G5: a daemon restart leaves no orphan process and no unexplained run.

**Non-goals**
- Network for long runs (D-140-2 stands).
- Streaming partial output to the agent (§9 R-145-3).
- Changing `HandlerCap`, `invokeDeadline` or `writeTimeout`.
- Loopback, unless D-145-7 rules otherwise.

## 3. Verification log

| V | what | evidence |
|---|---|---|
| V1 | Grader wall times: w190 `go-vet` 1 384 ms, `tool-gate-v0.54.15` 518 493 ms, `go-test-all` 243 961 ms, `verify-ail` 25 907 ms; w139t 19 617 / 486 693 / 243 790 / 25 130 ms | `~/.ailang/state/trial-2026-10-16/evidence/ailang-world/episodes/{w190,w139t}.json` `graders[].duration_ms` |
| V2 | 29 package lines in w190's `go test ./...`; 14 take > 9 s: daemon 242.6, broker 164.0, verifygate 156.1, ailang-worldd 110.7, capsule 87.0, transitionreg 74.1, coordinator 72.6, world-publish 72.0, replay 57.8, pkgproj 54.4, archive 46.6, evidence 29.8, projection 15.9, store 13.9 s | `q3_w190_go_test_all_packages.out.txt` (copied from the grader's stdout) |
| V3 | platform-source `test-fast`: `9 failed, 3913 passed, 50 skipped … in 46.11s`, `real 51.61` | `~/.ailang/state/trial-2026-10-16/trial/tps-e-base-testfast.txt` (tail) |
| V4 | w190: 91 `workspace-exec` calls, 91 results; 4 `timed_out:true`; 3 `bind: operation not permitted` | `q2_w190_exec.py` → `q2_w190_exec.out.txt` |
| V5 | w190's own summary lists `TestSeToolsManifestPublishes`, the confinement matrix, the verifygate gate arms and `verify_ail.sh` as unrunnable; round 2 was needed for exactly those count anchors | w190 transcript `result` event; friction #15 (PR #311) |
| V6 | srt 0.0.78 under World-shaped settings with `allowLocalBinding:false`: `bind` → `PermissionError [Errno 1]`; connecting to 127.0.0.1:47330 fails | `q1_run.sh` → `q1_run.out.txt` |
| V7 | Same settings with `allowLocalBinding:true`: bind ok, a connect to its own loopback server works (200) — **and an unauthenticated connect to the World daemon at 127.0.0.1:47330 succeeds**; external HTTPS still fails | `q1_run.out.txt` |
| V8 | The daemon's GET routes (`/v1/log`, `/v1/log/{index}`, `/v1/objects/{ref}`, …) pass through unauthenticated (declared residual R1); only `POST /v1/commit` is protected | `host/daemon/daemon.go:965-971` (`isProtected`) |
| V9 | World renders `network: {allowedDomains: [], deniedDomains: []}` and never sets `allowLocalBinding` (Go zero value false) | `host/broker/exec_settings.go:44-47, 116` |
| V10 | The short handler: verify the stack → `prepareCache` → `beginWorktreeDelta` → `runBounded` with `execTimeout = min(profile timeout, ctx deadline)` → head/tail capture → `execResult` | `host/broker/handlers_exec.go:210-297` |
| V11 | Children run in their own process group (`Setpgid: true`), killed as a group on the deadline (`syscall.Kill(-pgid, SIGKILL)`) | `host/broker/handlers.go:144, 167` |
| V12 | `procbound.MaxOutstanding = 8` is process-wide: admitted children not yet reaped; a 9th is refused at once (`ErrCleanupBacklog`) | `host/procbound/procbound.go:26-33` |
| V13 | Capture keeps 8 KiB head + 56 KiB tail per stream; the runaway kill is 64 MiB per stream (`DefaultExecMaxOutputBytes`) | `host/broker/handlers.go:271-272`; `handlers_exec.go:76-78` |
| V14 | The broker allows cost 0 (`withinEffectBudget: cost >= 0 && cost <= budget`) | `host/broker/decide.go:113-115` |
| V15 | `workspace-exec` declares one `Workspace.Exec@worktree` effect, cost 1; the plan cannot see the profile | `packages/se-tools/transitions.json` (`workspace-exec`); `packages/se-tools/se_tools/exec.ail` header |
| V16 | AG-UI emits only committed entries; the CUSTOM name table is closed (`world.entry.committed`); extensions need a new golden | `host/agui/agui.go:1-16` |
| V17 | No MCP progress notifications exist in the daemon (no `notifications/progress` or `progressToken`) | `grep` over `host/daemon`, `host/projection` (non-test): 0 hits |
| V18 | APFS clone (`cp -c -R`) of episode worktree w190 (73 MB, 4 856 files): 0.78 / 0.75 / 0.75 s; all 4 856 files present | `q4_clone.sh` → `q4_clone.out.txt` |
| V19 | `execWorldEnv` already pins `TMPDIR` (and HOME, caches) for exec children | `host/broker/exec_profile.go:52` |
| V20 | Replay returns the recorded output and never re-executes (row 140 §4.7) | `w-workspace-exec-toolchain-effect.md:334-338` |

## 4. Design

### 4.1 Two new tools, one new effect

- **`workspace-exec-start`** `{command, args?}`:
  - the same argument grammar as `workspace-exec` (V15), resolved against the same profile;
  - only commands the profile marks long-capable (§4.3) are admitted;
  - returns `{run_id, status:"running", snapshot_digest, started_at_ms, cap_ms}` within the handler cap (§4.4).
- **`workspace-exec-collect`** `{run_id, wait_ms?, cancel?}`:
  - returns `{run_id, status, elapsed_ms}`, where `status` is one of `running`, `finished`, `timed_out`, `cancelled` or `abandoned`;
  - once terminal, it also returns the full row-140 `execResult` (V10) with `worktree_delta` and the protected-path proof (§4.6);
  - `wait_ms` (0–8 000, default 0) long-polls *inside* the handler budget, never past it;
  - `cancel:true` kills the run's process group (V11) and returns `cancelled`.
- **Effect:**
  - Both tools run on a new effect, **`Workspace.ExecLong`**, scope `worktree`, with its own grant (D-145-1).
  - Start costs 1. Collect costs 0 (allowed, V14) but **still commits an entry**, World's rule for every call.
  - `run_id` is minted by the handler: 128 random bits, hex. It is bound to the episode, so a run id from another episode refuses.
- **Why not reuse `workspace-exec`:**
  - Its record and descriptor promise a complete result in one call (V15).
  - A separate effect lets the operator grant short calls freely and long runs sparingly.

### 4.2 Records and replay

- **The start record:**
  - holds the request `{command, args}`;
  - holds the output `{run_id, snapshot_digest, argv, profile, sandbox, cap_ms}`.
- **A collect record:**
  - holds the request `{run_id, wait_ms, cancel}`;
  - holds the output: the status object.
  - The terminal collect also carries `start_entry` (the start's entry ref), so a walk joins the two.
  - Exactly one collect per run carries the terminal result. Later collects return the same terminal status with `result_entry` (a ref) instead of repeating the streams.
- **Replay** returns each recorded output and never re-executes or re-polls (V20). A walk (`why`) on the terminal collect reaches the start through `start_entry`.

### 4.3 Profile and operator caps

- **Per command:**
  - A command gains an optional `long_timeout_ms`. Without it, the command is not long-capable.
  - `timeout_ms` keeps its ≤ 9 000 meaning for `workspace-exec`.
- **New `serve` flags** (operator):
  - `--exec-long-max-ms` (default 1 800 000 = 30 min): the ceiling every `long_timeout_ms` must be under, refused at startup otherwise.
  - `--exec-long-max-concurrent` (default 2): the process-wide pool (§4.5).
- The profile digest covers the new field, as every profile field already is.

### 4.4 Start inside the cap: snapshot, then detach

- **The steps** are row 140's (V10) up to the spawn:
  1. gate the profile and grammar;
  2. verify srt, `argv0` and the settings digests;
  3. `prepareCache`;
  4. **snapshot the worktree** with an APFS clone into `<cache>/runs/<run_id>/tree`;
  5. render that run's srt settings, with the worktree paths rebased onto the snapshot;
  6. spawn under `Setpgid`, without waiting;
  7. return.
- **Timing:** the clone took 0.75 s on a 73 MB / 4 856-file World worktree (V18). That is well inside a 10 s handler, but it is a single-repo measurement: M0 measures platform-source and the ailang compiler first.
- **Why a snapshot:**
  - The agent keeps editing while a suite runs. On the live worktree, a 244 s run would test a moving tree.
  - Its `worktree_delta` would mix the agent's own edits with what the run wrote.
  - The snapshot fixes the tested tree at the moment of start (`snapshot_digest`), and the run's writes stay in the snapshot.
  - **A long run therefore never changes the live worktree.** That is stricter than the short lane, and simpler to prove.
- **Off APFS** (Linux, or a non-APFS volume), the clone falls back to a full copy, capped by a size limit. Over the limit, start refuses with the measured size. M0 measures Linux `cp --reflink=auto`.

### 4.5 Bounds and concurrency

- **Wall clock:** the run's deadline is `long_timeout_ms` from spawn. On expiry, the process group is killed (V11) and the run is `timed_out`.
- **Output:** the same head/tail capture and 64 MiB runaway kill (V13), held in memory until collected. Terminal results are written to the run registry (§4.7), so a later collect survives a restart.
- **Per episode:** one running long run (D-145-4). A second start refuses with the running run's id.
- **Process-wide pool:**
  - `--exec-long-max-concurrent`, counted separately.
  - A long run must not hold one of `procbound`'s 8 slots for minutes (V12), or six long runs would refuse every short tool call in the daemon.
  - M2 gives long runs their own admission counter and keeps procbound for the reap guarantee.

### 4.6 `worktree_delta` and the protected-path proof

- **`worktree_delta`:** the delta is taken on the **snapshot**, between start and the end of the run (row 184's probe, rebased). It reports what the run wrote, which the agent never sees in its worktree.
- **Protected paths (row 156):**
  - The protected set is rendered for the snapshot path. The proof at termination is that every protected file in the snapshot is byte-identical to the start.
  - The snapshot itself is cloned from the worktree, so the start state is whatever the agent had. Row 156's write refusals already keep the agent from editing protected files there.
  - This design follows row 156's final rendering. If 156's API moves, M3 follows it.

### 4.7 Restart, cleanup, cancellation

- **Run registry:** `<state>/exec-runs/<ep>/<run_id>.json` (0600), written at spawn. It holds:
  - `{pgid, pid, process start time, snapshot path, deadline, start_entry}`;
  - the terminal result, once there is one.
- **On `serve` startup:** for every non-terminal registry entry:
  - kill the group if `pid` still has the recorded start time (a guard against PID reuse);
  - mark it `abandoned` with the reason `daemon restarted`;
  - delete the snapshot.

  The process outlives the daemon today, because it is in its own group (V11). This sweep is the only thing that stops a run the daemon is no longer watching. M0 measures that orphan survival as a premise.
- **Collect** of an `abandoned` run returns that status. There is no partial output.
- **Snapshots** are deleted at the terminal collect, or 24 h after the run ends, whichever comes first. An episode's revocation deletes all of them.

### 4.8 How the agent waits

- **Polling:** the agent polls `collect` with `wait_ms ≤ 8 000`. A 244 s suite costs about 30 collect entries; a `wait_ms` of 0 costs one entry per poll.
- **Descriptor:** it tells agents to start, keep working, then collect with `wait_ms`.
- **Push:** neither MCP progress (V17) nor an AG-UI event (V16) exists. AG-UI streams committed entries only, so a push would need the daemon to commit an entry no agent call made. That is a new kind of writer, deferred to D-145-6.

### 4.9 Loopback is a separate question (D-145-7)

- **Measured** (V6–V8): srt's `allowLocalBinding:true` does make `httptest`-style tests work.
- **But** it also lets the sandboxed command connect to **every** loopback service, including the World daemon. The daemon's read routes are unauthenticated (residual R1).
- So a test command in episode A could read every episode's log entries and object payloads (file contents written through World), plus anything else listening on loopback on the operator's machine.
- That is a confinement regression against D-140-1. It is not a profile tweak.

## 5. Premises

| premise | status |
|---|---|
| A long run of World's suite takes minutes (P1) | measured (V1, V2) |
| APFS clone is sub-second at World's size | measured (V18); other repos and Linux in M0 |
| A child survives daemon death | **assumed** from `Setpgid` (V11); measured in M0 |
| Collect at cost 0 is allowed by the broker | measured (V14) |
| Loopback opt-in exposes the daemon | measured (V7, V8) |

## 6. Milestones (sketch)

- **M0** — Premises:
  - clone time on platform-source and the ailang compiler;
  - Linux reflink and copy;
  - orphan survival after a `kill -9` of `serve`;
  - srt with a cwd inside the exec cache.

  Banked like Q1–Q4.
- **M1** — The profile `long_timeout_ms` and the `--exec-long-*` flags, with startup refusals and digest coverage. Pure and table-tested.
- **M2** — The broker: the run table, the separate admission pool, the start handler (snapshot, rebased settings, detached spawn), and collect/cancel with `wait_ms`. Fake-node tests for every status; real srt for one run.
- **M3** — `worktree_delta` and the protected proof on the snapshot; the terminal result shape.
- **M4** — The run registry and the startup sweep: tests for a kill-9 restart, an abandoned run, and a PID-reuse guard.
- **M5** — The se-tools transitions `workspace-exec-start.ail` / `-collect.ail`, the descriptors, the `Workspace.ExecLong` wiring (all the enumeration sites w190 listed), the docs, and the `session new` preset. Publish phrase: Mark's.
- **M6 (attended)** — Re-measure: an ailang-world episode runs `go test ./...` and `verify_ail.sh` itself through the lane, and its graders agree with what it ran. A pass needs no operator-only grader except the loopback e2e tests (P4).

## 7. Acceptance criteria

- **AC1:** no World call in the lane exceeds the D7 deadlines. `TestProjection_NoDeadlineTampering` and the coordinator's budget pins are unchanged.
- **AC2:** a 244 s `go test ./...` runs through start/collect.
  - Its terminal result has `exit_code`, digests and head/tail streams.
  - The live worktree is byte-identical before and after.
  - `worktree_delta` reports only what the run wrote, on the snapshot.
- **AC3:** confinement:
  - a long run cannot write the live worktree, read outside the fence, reach the network, or (unless D-145-7 = B) bind loopback;
  - each is tested with a negative control.
- **AC4:** bounds:
  - the wall-clock kill → `timed_out`;
  - a second start in an episode refuses;
  - the pool limit refuses at `max+1`;
  - short `workspace-exec` calls still run while the pool is full.
- **AC5:** restart: `kill -9` of `serve` mid-run, then restart, gives `abandoned`, no surviving process group, and the snapshot removed.
- **AC6:** replay of start and collect entries returns the recorded outputs with zero spawns.
- **AC7:** a run id from another episode refuses at collect.

## 8. Decisions for Mark (OPEN, each with a recommended default)

- **D-145-1 — Effect and budget.**
  - **A (recommended):** a new effect, `Workspace.ExecLong`, with its own grant. Start costs 1; collect costs 0 and is recorded.
  - B: reuse `Workspace.Exec`, with start costing K (e.g. 10) and collect costing 0.
  - C: reuse `Workspace.Exec` at cost 1 for every start and collect.
- **D-145-2 — Tool shape.**
  - **A (recommended):** two tools, `workspace-exec-start` and `workspace-exec-collect` (with `wait_ms` and `cancel`).
  - B: three tools (cancel separate).
  - C: one tool with a `mode` argument.
- **D-145-3 — Which commands may run long.**
  - **A (recommended):** only those whose profile entry sets `long_timeout_ms` ≤ `--exec-long-max-ms` (default 30 min).
  - B: any profile command, at one operator-wide cap.
- **D-145-4 — Which tree a long run sees.**
  - **A (recommended):** an APFS snapshot at start (V18, 0.75 s), so the live worktree is never written by a long run.
  - B: the live worktree, with a tree digest at start and end and a `tree_changed_during_run` flag.
  - C: the live worktree, with World's write tools refused while a run is active.
- **D-145-5 — Concurrency.**
  - **A (recommended):** one long run per episode, and a process-wide pool of 2 outside procbound's 8 (V12).
  - B: N per episode (operator flag).
- **D-145-6 — Completion signal.**
  - **A (recommended):** polling only in v1 (`collect` with `wait_ms` ≤ 8 s).
  - B: also a daemon-written "run finished" entry, surfaced as a new AG-UI CUSTOM name for the workbench and desktop (rows 188 and 198). This needs a new writer identity and a golden (V16).
- **D-145-7 — Loopback for test commands (friction #4).**
  - **A (recommended for this row):** keep it denied, and file a separate row. That row's premise is V7/V8: opening loopback exposes the World daemon's unauthenticated reads (every episode's records and payloads) and any other local service. Its candidate fixes:
    1. close residual R1 (authenticate the GET routes), then allow a per-command `allow_local_binding` opt-in;
    2. per-episode network namespaces (Linux only);
    3. tests that need loopback stay operator-run graders.
  - B: allow a per-command opt-in now, accepting the cross-episode read exposure on a single-operator machine.
  - C: allow it for long runs only (no different in exposure from B).

## 9. Risks and residuals

- **R-145-1:** a snapshot of a very large worktree (node_modules, venvs in-tree) may not clone fast off APFS. Start refuses with the measured size rather than blocking (§4.4). M0 sizes the real projects.
- **R-145-2:** polling commits about 30 entries per 4-minute run. That is acceptable for the log, and row 198's episode report should group a run's collects.
- **R-145-3:** no partial output while running. An agent waiting on a hung suite sees only `elapsed_ms` until the cap. A follow-up could return the current tail at collect, unrecorded beyond the status.
- **R-145-4:** the loopback tests (P4) stay operator-only until D-145-7's separate row lands. M6 measures how much of World's gate that still leaves out.

## 10. Dependants

- **Row 196** (loop executors through World) cannot run the loop's own gate through World without this lane: P1, 244–518 s. Row 145's charter text names it the "Gate for R1 step 3 (the loop on World)".
- **Row 198** (`episode report`) should count a long run once: start, then terminal collect.
- **Row 188** (desktop) is where D-145-6 B's push signal would show.
