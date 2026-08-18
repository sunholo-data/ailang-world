# Sprint plan — `w-daemon-read-cancellation` (queue item 18)

**Item**: queue item 18, `w-daemon-read-cancellation` — a bounded elapsed-time contract for the daemon's read path
**Status**: PLANNED · three milestones · **1.5 days** · three commits, M1 → M2 → M3
**Design doc**: [`w-daemon-read-cancellation.md`](w-daemon-read-cancellation.md) (711 lines, designed iter-86, revised in quorum round 2)
**Ratification**: **`D-WORLD-18` RESOLVED — option A, "ship the scoped item as designed"**, ratified attended by Mark 2026-08-17. Doc §13: *"A unparks this doc to sprint-planner as written."* **There is no designer round and no re-litigation of scope.** The one surviving round-2 objection (gpt5-6-sol's scope-boundary dispute) is settled; its remedy is the named follow-on `w-bounded-waits-operator-and-write-paths`, not this sprint.
**Base**: `dev` @ `d31362f`, both gates green (measured, §0)
**Planner**: mission-control iteration 91, opus sprint-planner, first-party measurement on this rig
**Executor**: sandboxed worktree, **NO git write permission**. **THE CONTROLLER MAKES ALL COMMITS.** The executor never runs `git commit`, `git add`, `git push`, `git checkout`, `git stash`, `git restore`, or `gh pr`. Restores are `cp` from a backup — see §7.

---

## 0. Planner's first-party verification of every premise

Every load-bearing present-tense claim in the design doc was re-run on this rig at `d31362f`
before it entered this plan. **Every zero carries a known-positive control in the same call**, per
the standing rule that an empty result is a claim, not a fact. zsh traps observed and avoided:
`--include='*.go'` is quoted at every site (unquoted it glob-expands and the command never runs);
`${PIPESTATUS[0]}` is never read (silently empty in zsh) — rc is captured directly.

### 0.1 The design doc's Verification Log, re-derived at HEAD

| Doc row | Claim | Verdict at `d31362f` | Command → observed output |
|---|---|---|---|
| V1 | `host/store/store.go` has zero context plumbing | **CONFIRMED** | `grep -c "context.Context" host/store/store.go` → **0**; control **in the same call** `grep -c "func (s \*Store)" host/store/store.go` → **14** |
| V2 | **five** context-free read getters on the daemon read path | **CONFIRMED, line numbers exact** | `grep -n 'func (s \*Store) Get\|func (s \*Store) Selected' host/store/store.go` → `GetObject:467`, `GetWorld:522`, `GetLogEntry:551`, `GetRegistryHead:628`, `GetVerifyResult:773`, `SelectedHead:802` — none takes a context. MU4a–e's cited sites 467/522/551/628/802 are byte-correct |
| V3 | each getter blocks in a context-free `QueryRow` | **CONFIRMED** | `grep -c 's.db.QueryRow(' host/store/store.go` → **5** |
| V4 | no handler consults the request context | **CONFIRMED, both controls fire** | `grep -c 'r.Context()' host/daemon/handlers.go` → **0** (control `http.ResponseWriter` → **9**); `grep -c 'r.Context()' host/daemon/daemon.go` → **0** (control `mux.HandleFunc` → **8**) |
| V5 | 11 `err.Error()` echoes, six on `Internal` branches | **CONFIRMED, all 11 line numbers exact** | `grep -n 'err.Error()' host/daemon/handlers.go` → Internal at **220, 247, 277, 325, 351, 419**; BadRequest at **215, 242, 392, 399, 410**. Count **11** |
| V6 | no production `busy_timeout`; only hit is a test input | **CONFIRMED** | `grep -rn 'busy_timeout' host/store/*.go \| grep -v _test \| wc -l` → **0**; same-scope control (drop the `-v _test`) → **1**, `writer_lock_test.go:609` `{":memory:?_pragma=busy_timeout(1000)", true}` |
| V7/V23 | DSN builders live in `writer_lock.go`, not `store.go` | **CONFIRMED** | `grep -n 'func resolveDSN\|func writeDSN\|func readOnlyDSN' host/store/*.go` → `writer_lock.go:120, 176, 187`; zero hits in `store.go` |
| V8 | one physical connection serializes every read and write | **CONFIRMED** | `grep -n 'SetMaxOpenConns' host/store/store.go` → **297**: `db.SetMaxOpenConns(1)` |
| V9 | 7 GET + 1 POST, frozen route table | **CONFIRMED** | `daemon.go:461-468`, eight `mux.HandleFunc` registrations; `handleHead` at `daemon.go:497` calls `d.store.SelectedHead()` |
| V10 | the D7 gate pins six literals and observes the transport only | **CONFIRMED** | `TestBoundedWaitsAndBodyLimit` at `daemon_test.go:202`; constants table rows `readHeaderTimeout 5s`, `readTimeout 30s`, `writeTimeout 30s`, `idleTimeout 120s`, `defaultClientTimeout 30s`, `shutdownTimeout 10s` |
| V11 | the Go envelope is a frozen mirror of the sketch | **CONFIRMED** | `handlers.go:16-18` — *"Class names and status codes mirror ApiError/httpStatus in the frozen, checked design_docs/sketches/worlddapi.ail exactly"* |
| V12 | the sketch edit moves **no** gate pin | **CONFIRMED — and PROVEN END-TO-END, see §0.3** | `EXACT_TOTAL_VERIFIED=10` (**line 323**, not 311), `EXACT_TOTAL_TESTS = 39` (**line 363**, not 350, and it has **SPACES around `=`**), `LEG1_MODULES` at :140 with `design_docs/sketches/worlddapi.ail` at :145 |
| V13 | the exact `.ail` edit verifies on the pinned binary | **REPRODUCED FIRST-PARTY** | see §0.3 |
| V15 | migration ripple 22 non-test / 86 total | **CONFIRMED EXACTLY** | quoted-include grep over `host cmd` → **TOTAL 86**, **non-test 22**, **test 64**. Non-test by file: `broker/approve.go` 8, `broker/broker.go` 2, `daemon/daemon.go` 1, `daemon/handlers.go` 5, `registry/registry.go` 2, `replay/replay.go` 1, `transitionreg/transitionreg.go` 3 |
| V16 | stdlib `context` cannot red the boundary gate | **CONFIRMED, and the near-miss checked** | `forbiddenImportPrefixes` (`allowlist_world_test.go:61`) = `cloud.google.com/`, `github.com/Azure/`, `github.com/aws/`, `…/host/registry`, `net/http/httptest`, `net/http/httputil`. The single `"context"` hit in that file is **its own import at line 8**, not a list entry. `wantFileCount = 1` at :1163 |
| V17 | no daemon error-log surface exists today | **CONFIRMED** | `log.` calls in non-test `host/daemon` → **0**; same-scope control `fmt.` → **16**. `ListenAnnouncePrefix` at `daemon.go:140` |
| V18 | base gates are green | **CONFIRMED AT HEAD** | `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` → **rc=0**, *"✓ verify gate PASSED: 10 required identities verified, 39 named tests pass"*, *"world package gate PASSED: 9/9"*. `GOTOOLCHAIN=go1.25.6 go build ./...` → rc=0 |
| V20 | AC2/AC9's vacuous `-run` baseline | **REPRODUCED** | `go test ./host/store -run '…' -v -count=1` → **rc=0**, `=== RUN` count **0**, `testing: warning: no tests to run` |
| V21 | AC3's baseline | **REPRODUCED EXACTLY** | `go test ./host/daemon -run '…6 selectors…' -v -count=1` → rc=0, **5** `=== RUN` lines, all from `TestBoundedWaitsAndBodyLimit` and its four subtests; the five new test functions enumerate **0** |
| V24 | `writeAPIError` is the single funnel | **CONFIRMED** | `grep -c 'writeAPIError' host/daemon/handlers.go` → **22** |
| V25 | the store has SIX getters; `GetVerifyResult` is off the daemon path | **CONFIRMED** | `GetVerifyResult:773`; zero daemon call sites (its only non-test caller is `host/replay/replay.go:153`) |
| V26 | five non-daemon sites already have a context in reach | **CONFIRMED** | `transitionreg.go:70 ReadSnapshot(ctx context.Context)`, `:221 Publish(ctx context.Context, …)`; `broker.go:155 Session.Invoke`, `:163 s.invoke`, **`:173` drops the ctx** calling `s.invokeReplay(req, decision)` (`:352`) |
| V27 | `Store.Commit` has exactly one production caller | **CONFIRMED** | `grep -rn '\.Commit(' --include='*.go' host cmd \| grep -v '_test\.go'` → **9** hits, of which **1** is `d.store.Commit` (`handlers.go:413`) and **8** are `database/sql` `tx.Commit()` in `host/store/{journal,store}.go` |
| §2.8 | the ratchet's **11** residual deadline-free store-read sites | **CONFIRMED EXACTLY, 8/2/1** | `grep -v '_test\.go' <all-sites> \| grep -cE 'approve\.go\|registry/registry\.go\|replay/replay\.go'` → **11**. Sites: `approve.go:169,195,225,234,251,261,495,522` (8), `registry/registry.go:127,135` (2), `replay/replay.go:153` (1). **This is a different scope from the doc's "10 production `context.Background()` literals" — do not conflate the two counts** |

### 0.2 Base-condition drift since the doc's measurement base `6fd26f0`

`git diff --name-only 6fd26f0..d31362f` → **9 files**, of which **1** is a `.go` file
(`host/verifygate/ail_binary_gate_test.go`) and one is `scripts/verify_ail.sh`. Both moved in
iteration 89's stderr-merge repair. **Consequences for this sprint, all benign but all stated:**

1. **V22 is now FALSE as written.** The doc says the base's Go tree is byte-identical to the
   gate-proven merge `aaada20`. At HEAD one Go file differs. It is in `host/verifygate` — a
   package this item does not touch — and the full suite is green anyway (§0.4), so the claim's
   *conclusion* survives; only its stated mechanism (byte-identity) does not. Re-derived directly
   instead of inherited.
2. **The doc's §8 gate-pin line numbers are STALE.** `EXACT_TOTAL_VERIFIED` is at **line 323**
   (doc says 311); `EXACT_TOTAL_TESTS` is at **line 363** (doc says 350). The *values* (10 and 39)
   and the disposition (**not moved**) are both correct. An executor grepping by line number will
   read the wrong line.
3. **`EXACT_TOTAL_TESTS = 39` has spaces around `=`** (it is inside the embedded python, not
   shell). A shell-shaped `grep 'EXACT_TOTAL_TESTS='` **misses it**. Carried trap from item 15.

### 0.3 The `.ail` edit and AC7, PRE-REGISTERED — proven twice, not argued

The doc's V13 was a scratch probe. This plan proves the same edit **against the real repo's gate**
and restores the tree byte-identically, so AC7's pass condition is a pre-registered value rather
than an expectation.

**Probe A — isolated module copy** (`/tmp/w18probe/sketches/`, the layout the sketch's own
`import sketches/logepoch` requires; **the wrong layout is an instrument failure that reads as a
kill on every arm** — my first attempt at `design_docs/sketches/` returned `check.passed=False`
with `LDR001 module not found: sketches/logepoch`, caught only because the *pristine control* was
run first):

| Arm | `check.passed` | `verify.errors` | `cex` | `passed_tests` |
|---|---|---|---|---|
| pristine control | **True** | 0 | 0 | **18** |
| after the three-anchor edit | **True** | 0 | 0 | **19** |

The edit is exactly three anchors and nothing else (diff reproduced verbatim in §3.2, task T2.5).

**Probe B — the real repo, the real gate.** `design_docs/sketches/worlddapi.ail` backed up by
`cp`, edited, `verify_ail.sh` run, restored:

```
PRE  sha256 = 3a83a1dd05e427f47452c32ac3ad1b05524d9b902c00a4173555baf126c9ef2b
POST sha256 = cb7f7f89b098d26fa5e683e873ff4eee6b080abe5df2a2b4969a680e7200a615   (landed-proof: PRE != POST)
./scripts/verify_ail.sh  ->  rc=0
   ✓ 10/10 required world/ identities verified across 11 module(s)
   ✓ all 39 required named tests pass (failed_tests=0)
   ✓ world package gate PASSED: 9/9 steps performed non-zero work
   ✓ verify gate PASSED: 10 required identities verified, 39 named tests pass
RESTORED sha256 = 3a83a1dd...  == PRE   (byte-identical restore verified)
```

**AC7 is therefore pre-registered: `10 / 39 / 9-of-9`, rc=0, unchanged.** If the executor's run
reports anything else, the edit is not the one probed here — **STOP and report**, do not adjust a pin.

One precision the doc does not state and an executor will trip over: the sketch runner's
`passed_tests` moves **18 → 19**, but `len(tests[])` in the same JSON is **16 in both arms** — the
new row lands in `passed_tests`, not in the `tests[]` array. `verify_ail.sh` pins `len(tests[])`
(line 363), and Leg 2 runs `ailang test world/` anyway, so the sketch is doubly invisible to it.
**Never quote `passed_tests` at a gate.**

### 0.4 Base gate state at `d31362f`, measured this iteration

```
GOTOOLCHAIN=go1.25.6 go build ./...                  -> rc=0
GOTOOLCHAIN=go1.25.6 go test ./... -count=1          -> rc=0, 17 packages ok, 0 FAIL
AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh -> rc=0 (10 / 39 / 9-of-9)
/tmp/ailang-v0300/ailang --version                   -> AILANG v0.30.0, Commit e37b370
```

`host/broker` passed in **50.8 s** on this run. That is **one** observation against a recorded
~18% base flake: it is **not** evidence the flake is gone, and a red broker leg after this diff is
**not** automatically this diff's fault. Attribute by shape against the base rate.

### 0.5 Findings that CHANGE the plan

#### (i) **PLAN-AFFECTING DOC DEFECT — the `readStore` seam is FIVE methods, not six. The doc says "six" twice.**

Design doc §2.5 layer 3: *"a **six-value** read interface `reads readStore`"*. Design doc §2.7:
*"The single new Go surface is a **six-method** unexported interface inside `host/daemon`."*

Measured:

```
$ grep -v '_test\.go' <all getter call sites> | grep 'host/daemon'
host/daemon/handlers.go:218  d.store.GetWorld
host/daemon/handlers.go:245  d.store.GetObject
host/daemon/handlers.go:275  d.store.GetLogEntry
host/daemon/handlers.go:323  d.store.GetLogEntry     <- handleLogRange, the SAME getter
host/daemon/handlers.go:349  d.store.GetRegistryHead
host/daemon/daemon.go:497    d.store.SelectedHead
ROUTE-LEVEL CALL SITES: 6      DISTINCT GETTERS: 5
```

**Six routes, six call sites, FIVE distinct getters** — `GetLogEntry` serves both `/v1/log/{index}`
and `/v1/log`. The doc's own §2.5 is internally inconsistent with itself: the same paragraph that
says "six-value" also says `blockingStore` *"overrides **all five** getters"*, and §2.5 layer 4
says `recordingStore` *"overrides **all five** getters"*. The "six" is the **route** count leaking
into a **method** count — the iteration-86 class ("a count is only true inside the scope it was
taken in") recurring in the same document that records it.

**Resolution, binding on the executor**: `readStore` has **exactly five methods**
(`GetObject`, `GetWorld`, `GetLogEntry`, `GetRegistryHead`, `SelectedHead`). **Do not invent a
sixth.** The two candidate sixths are both wrong and both harmful:
- `Commit` — the write path; §8/§10 declare it untouched (DR-1). Putting it behind a "reads"
  interface widens the seam past the item's ratified scope.
- `GetVerifyResult` — off the daemon read path entirely (V25); adding it makes the fake's
  coverage claim false and pulls the follow-on item's work forward.

Also **out of the seam**: `d.store.ScanUnreadableLog` (`daemon.go:392`) and
`d.store.ScanUnreadableWorlds` (`daemon.go:414`). These are store reads on the daemon, but they
live in the **integrity-report** path, not on any `/v1` GET route. They keep using `d.store`
directly. A `readStore` that absorbs them is out of scope and reds nothing that matters.

**AC-level consequence**: no acceptance criterion counts the interface's methods, so this defect
cannot red a gate — which is exactly why it must be settled here rather than discovered at
compile time by an executor guessing at a sixth name.

#### (ii) **A live inter-item collision on `host/store/store.go`: `D-WORLD-19` is OPEN.**

`D-WORLD-19` (opened iteration 90, unanswered) asks whether item 17's tranche 1 may extend
`host/store` with a bounded object read — *"putting a second item into `host/store` while item 18
is queued to bound that same package."* Its arm A would change `GetObject` (`store.go:467`) to
`OpenObject(ref) (io.ReadCloser, error)` or add a `maxBytes` bound. **That is the same function
M1 changes to `GetObject(ctx context.Context, …)`.**

This does **not** block this sprint — item 18 is ratified and routable, item 17 is parked — but it
is a real merge-conflict surface. Recorded in §8 as risk R5. If `D-WORLD-19` resolves A while this
sprint is in flight, item 17's tranche must rebase onto item 18's signature, not the reverse:
item 18 landed first and its ratification is older.

#### (iii) The doc's non-blocking finding #1 is still true and still one word.

`daemon.go:456` reads *"The seven patterns below"* above **eight** registrations. Verified at HEAD.
The doc says the sprint *may* fix it in passing but **must not count it as scope**. This plan
takes it: one word, in M2's `daemon.go` diff, mentioned here so it is not mistaken for drift.

---

## 1. Scope

### 1.1 What this sprint delivers

The design doc **exactly as written**, in its own §9 milestone structure: a request-scoped
deadline on the daemon's six store-reading GET routes, reaching the blocking `database/sql` call;
an explicit `503`/`Timeout` status mirrored in the frozen sketch; a lock-layer `busy_timeout`
policy; sanitized 500 bodies with the detail moved to a new `Config.ErrorLog`; and a ratchet test
pinning the 11 residual deadline-free store reads so the set can shrink but never grow.

### 1.2 What this sprint does NOT do — ratified deferrals, do not "fix"

Every line below is a **ratified** deferral (`D-WORLD-18` arm A). Doing any of it fails the sprint.

- **DR-1**: `POST /v1/commit` keeps only its transport bounds. `Store.Commit` keeps its signature.
- **DR-2**: `broker/approve.go` (8), `registry/registry.go` (2), `replay/replay.go` (1) receive
  `context.Background()` **verbatim** — the visible spelling of today's behavior. **Do not thread a
  real deadline through them.** Doing so shrinks the AC9 pin and reds AC9.
- **No store-boundary guard** rejecting nil/deadline-free contexts. It is the follow-on's closing
  move, landable only when the ratchet reads zero. Landing it now breaks all 11 sites (V27).
- No `GetVerifyResult` context migration (V25 — follow-on shape (b)).
- No SSE, no `http.ResponseController`, no streaming route, no retry above `busy_timeout`.
- No new route, no route-table edit, no store schema change, no gate-pin move.

### 1.3 FROZEN — the plan does not touch these, and the executor must not either

| Path / surface | Why |
|---|---|
| `tools/launchd/*` | frozen core (CLAUDE.md); FLEET-owned under `D-WORLD-DRIVER-1`. **Currently MODIFIED/UNTRACKED in the working tree on purpose.** |
| `scripts/verify_go.sh` | frozen core; **currently MODIFIED on purpose** (same fleet bundle) |
| `CLAUDE.md` | frozen core; **currently MODIFIED on purpose** (same fleet bundle) |
| `world/*.ail`, `packages/world-core/**` | this item edits **zero** `world/` files |
| the frozen `/v1` route table (`daemon.go:461-468`) | no route added or removed |
| the commit envelope, `Store.Commit`, any store schema | DR-1 |
| every gate pin (`EXACT_TOTAL_VERIFIED`, `EXACT_TOTAL_TESTS`, `LEG1_MODULES`, `verify_world_package.sh`, `wantFileCount`) | proven not-moved in §0.3 |

**`git add -A` IS FORBIDDEN.** The working tree carries a deliberate uncommitted fleet bundle
(`CLAUDE.md`, `scripts/verify_go.sh`, `tools/launchd/mission-control.sh`,
`tools/launchd/derive-planner-lane.sh`, `tools/launchd/test_mission_routing.sh`,
`tools/launchd/testdata/`). Do not stage, commit, revert, or edit any of it. The executor holds
no git write permission at all; this note exists so the **controller** does not sweep it in either.

**Do not run `scripts/verify_go.sh` as this sprint's gate.** Its driver-drift gate reds while that
bundle diverges from HEAD, and that red means *"the fleet must commit"*, never *"absorb it"*. This
sprint's Go gate is `go build ./... && go test ./...` run directly (§5).

---

## 2. Acceptance criteria → milestone → named mutation

The standing rule: **a guard is not a gate until something reds when you remove it.** The doc's §6
supplies 12 mutation rows. Mapped below per-milestone, with every AC that lacks one named
explicitly rather than quietly.

| AC | What it asserts | Closed by | Named mutation(s) that red it | Status |
|---|---|---|---|---|
| **AC1** | `context.Context` count in `store.go` ≥ 5, control ≥ 14 | M1 | **MU13 (planner-added)** | was **UNCOVERED** in the doc — see §2.1 |
| **AC2** | store-layer tests exist and pass (≥ 7 `=== RUN`) | M1 | MU4a, MU4b, MU4c, MU4d, MU4e, MU5, MU6 | covered (7 arms) |
| **AC3** | daemon-layer tests exist and pass (5 new funcs enumerate) | M2 (all but sanitize), M3 (sanitize) | MU1, MU2, MU3, MU7, MU8, MU9, MU11 | covered (7 arms) |
| **AC4** | `r.Context()` ≥ 1 in both `handlers.go` and `daemon.go` | M2 | **MU2** (semantic form: `r.Context()` → `context.Background()`) | covered — but see §2.1 |
| **AC5** | `err.Error()` count is exactly 5, and they are the BadRequest sites | M3 | MU7 | covered |
| **AC6** | `busy_timeout` ≥ 1 in non-test `host/store` | M1 | MU5, MU6 | covered |
| **AC7** | `verify_ail.sh` rc=0 with the SAME totals (10 / 39 / 9-of-9) | M2 | MU10 is **declared NOT CI-red** (honest one-sided pin); **MU14 (planner-added)** is the non-vacuity control | was **UNCOVERED** — see §2.1 |
| **AC8** | `go build ./... && go test ./... -count=1` rc=0 | M1, M2, M3 (each independently) | **MU13 (planner-added)** proves the compiler is the guard on the 86-site migration | was **UNCOVERED** — see §2.1 |
| **AC9** | ratchet pins {approve.go: 8, registry.go: 2, replay.go: 1, else 0} | M1 | MU12 | covered |

### 2.1 The three ACs the design doc leaves without a mutation, and what this plan does about each

The doc's §6 covers ACs 2, 3, 5, 6 and 9 well. **AC1, AC4, AC7 and AC8 have no §6 row.** Stated
plainly rather than smoothed over, with a remedy for each:

- **AC1 — no doc mutation.** AC1 is a grep snapshot. Its real guard is the **compiler**: once the
  five getters are context-first, a getter reverting to a context-free signature cannot build
  against 86 call sites that pass one. That is a stronger guard than a grep and it is currently
  unnamed. **Added as MU13** (§4.3): revert `GetObject`'s signature to context-free, assert
  `go build ./...` rc **≠ 0** with an error naming `GetObject`, restore, assert rc=0 again. MU13
  discharges AC1 and AC8 together.
- **AC4 — covered in substance, not by name.** MU2 (`r.Context()` → `context.Background()` in
  `readCtx`) is precisely AC4's mechanism and reds `TestDaemonReadDisconnect`. AC4's *grep* form
  has no mutation and needs none: it is a snapshot whose persistent form is MU2's test.
  **Mapped, not added.**
- **AC7 — MU10 is declared not-CI-red by the doc, honestly.** Leg 2 runs `ailang test world/`, so
  the sketch's inline test row is outside both gate legs (V12, re-derived §0.1). The doc's enforced
  guard for that drift direction is MU9's Go-side mirror test — correct, but it leaves **AC7's own
  instrument** (Leg 1's ai-check sweep of the sketch) unproven. **Added as MU14** (§4.3): introduce
  a deliberate syntax error in `worlddapi.ail`, assert `verify_ail.sh` reds naming that file,
  restore byte-identically, assert rc=0. MU14 proves Leg 1 actually sweeps the file the sprint
  edits — i.e. that AC7 can fail. Without it, AC7 is a gate nobody has shown can red.
- **AC8 — the gate is its own guard, plus MU13.** No further arm needed.

### 2.2 Mutation coverage per milestone (the commit-time obligation)

| Milestone | Mutation arms that must be run and killed before its commit | Count |
|---|---|---|
| M1 | MU4a, MU4b, MU4c, MU4d, MU4e, MU5, MU6, MU12, **MU13** | **9** |
| M2 | MU1, MU2, MU3, MU8, MU9, MU11, **MU14**, (MU10 = declared, run locally, recorded as NOT-CI-red) | **7 + 1 declared** |
| M3 | MU7 | **1** |
| **total** | | **17 arms + 1 declared** |

---

## 3. Milestone breakdown, day by day

Three commits. **Each milestone is independently CI-green and committable.** Order is strict:
M1 → M2 → M3. M2 depends on M1's context-first getters; M3 depends on M2's `read_deadline_test.go`
file existing.

### 3.1 M1 — store layer (0.5 d, ~520 LOC) → greens **AC1, AC2, AC6, AC8, AC9**

Daemon behavior is **unchanged** after M1: a request context without a deadline is exactly today's
bound. That is what makes M1 independently landable.

| # | Task | Files | ~LOC | Closes |
|---|---|---|---|---|
| T1.1 | Five getters become context-first; `QueryRow` → `QueryRowContext(ctx, …)` at all 5 sites | `host/store/store.go` (:467, :522, :551, :628, :802) | 20 | AC1 |
| T1.2 | `writeDSN`/`readOnlyDSN` inject `_pragma=busy_timeout(2000)` **only when the caller's DSN did not set one** | `host/store/writer_lock.go` (:176, :187) | 30 | AC6 |
| T1.3 | Migrate the **22 non-test** call sites: daemon 6 → `r.Context()` is M2's job, so **M1 passes `context.Background()` at the 6 daemon sites** and M2 replaces them (see trap below); transitionreg 3 → the ctx `ReadSnapshot`/`Publish` already hold; `invokeReplay` gains a `ctx context.Context` first parameter fed from `Session.invoke` (`broker.go:173`); approve 8 + registry 2 + replay 1 → `context.Background()` **verbatim** | `host/daemon/{daemon,handlers}.go`, `host/transitionreg/transitionreg.go`, `host/broker/{broker,approve}.go`, `host/registry/registry.go`, `host/replay/replay.go` | 30 | AC8 |
| T1.4 | Migrate the **64 test** call sites (mechanical, `context.Background()`) | 18 files across `host/{store,broker,daemon,transitionreg,registry}` — distribution in §3.1.1 | 64 | AC8 |
| T1.5 | `TestReadGettersHonorContext` — 5 subtests, real `:memory:` store, sole pool connection held via `s.db.Conn(…)`, already-expired ctx, 2 s watchdog whose **red path closes the occupying `*sql.Conn`** and then drains the result channel under a second 2 s bound | `host/store/context_read_test.go` (new) | 150 | AC2 |
| T1.6 | `TestProductionDSNSetsBusyTimeout` — `PRAGMA busy_timeout` **readback from the live connection** on both the write and read-only handle | same file | 60 | AC2/AC6 |
| T1.7 | `TestReadRetriesUnderTransientExclusiveLock` — file-backed DB, `BEGIN EXCLUSIVE` from a second raw driver conn, released at 300 ms, getter carries a 5 s ctx, bound assertion **≤ 3 s**. **PRE-REGISTERED OUTCOME REQUIRED** (below) | same file | 90 | AC2 |
| T1.8 | `TestNoNewDeadlineFreeStoreReads` — the §2.8 ratchet. Scans non-`_test` `.go` sources of the caller packages, rooted via `runtime.Caller`, pins **{approve.go: 8, registry.go: 2, replay.go: 1, all else: 0}** | same file | 80 | AC9 |
| T1.9 | Mutation sweep: MU4a–e, MU5, MU6, MU12, MU13 (9 arms) | — | 0 | §2.2 |

**T1.7's pre-registered outcome (the MU15 lesson — record which form landed).** The doc
deliberately does **not** claim whether the driver's busy sleep is context-interruptible. The
executor **must record which of these two the measurement shows**, in the milestone note:
- **(a) retry-wins**: the read returns the row at ≈300 ms once the lock releases (busy_timeout's
  retry loop won).
- **(b) interrupt-wins**: the read returns a context/interrupt error before 2000 ms.
Both satisfy the ≤ 3 s bound. **Neither is a failure.** What *is* a failure is not recording which
one happened. Under MU5 (injection deleted) the read must fail **instantly** with `SQLITE_BUSY`.

**T1.3's trap, stated because it is the one place M1 and M2 can silently disagree.** M1 must leave
the six daemon sites compiling and behaviour-identical. Passing `context.Background()` there is
correct for M1 **and it makes the AC9 ratchet read 11 + 6 = 17 if the scanner's scope includes
`host/daemon`.** Two ways out, and the plan takes the second: (a) widen the M1 pin then shrink it
in M2 — two edits to the same pin, a moving target; (b) **`readCtx` is a one-line helper landed in
M1**, returning `r.Context()` with **no** timeout (`context.WithCancel`), so the six daemon sites
pass `d.readCtx(r)` from M1 onward and never spell `context.Background()`. M2 then changes only
the helper's body to `context.WithTimeout(…, d.readDeadline)`. **This keeps the AC9 pin at 11 in
both milestones, gives MU1 its exact one-line mutation site, and is behaviour-identical in M1**
(a cancel-only context derived from the request has exactly today's bound). The `defer cancel()`
mandate applies from M1.

#### 3.1.1 The 64 test call sites, by file (measured — an executor that migrates 63 gets a red build, not a silent pass)

`store/store_test.go` 10 · `broker/handlers_test.go` 9 · `daemon/handlers_test.go` 7 ·
`broker/broker_test.go` 7 · `broker/publish_op_test.go` 5 · `broker/episode_test.go` 4 ·
`store/durability_repro_test.go` 3 · `broker/recover_test.go` 3 ·
`transitionreg/transitionreg_test.go` 2 · `store/recover_test.go` 2 · `store/crash_test.go` 2 ·
`registry/registry_test.go` 2 · `broker/registry_publish_test.go` 2 · `broker/approve_test.go` 2 ·
`store/writer_lock_test.go` 1 · `store/journal_test.go` 1 · `daemon/daemon_test.go` 1 ·
`broker/handler_error_repro_test.go` 1 — **18 files, 64 sites.**

### 3.2 M2 — daemon deadline + explicit status (0.75 d, ~610 LOC) → greens **AC3 (minus sanitize), AC4, AC7**

| # | Task | Files | ~LOC | Closes |
|---|---|---|---|---|
| T2.1 | `readDeadline = 10 * time.Second` constant + `Daemon.readDeadline time.Duration` field set by `New` (the `drainTimeout` field-not-constant idiom, `daemon.go:238-243, 343`) | `host/daemon/daemon.go` | 15 | AC3 |
| T2.2 | `readCtx` body becomes `context.WithTimeout(r.Context(), d.readDeadline)`; **every one of the six call sites keeps its `defer cancel()` immediately after** (§2.3 mandate) | `host/daemon/handlers.go` | 10 | AC3/AC4 |
| T2.3 | The **five-method** `readStore` seam (§0.5(i)); `New` wires it to the same `*store.Store`; the six read handlers read through `d.reads`. `handleHead`'s `_ *http.Request` becomes `r`. Fix `daemon.go:456`'s "seven patterns" → "eight" (one word, non-scope) | `host/daemon/daemon.go`, `host/daemon/handlers.go` | 60 | AC3/AC4 |
| T2.4 | Timeout classification: after a failed store call check **`ctx.Err() != nil` FIRST**, then fall through to `Internal`. `errors.Is(err, context.DeadlineExceeded)` checked as well, but **the context is the authority** — the driver's interrupt can surface as `SQLITE_INTERRUPT` which does not wrap `context.DeadlineExceeded`. Emit `503` / class `Timeout` through `writeAPIError` | `host/daemon/handlers.go`, `host/daemon/daemon.go` | 55 | AC3 |
| T2.5 | The `.ail` sketch edit — **exactly three anchors, verbatim, no other change** | `design_docs/sketches/worlddapi.ail` | 3 | AC7 |
| T2.6 | `TestBoundedWaitsAndBodyLimit` gains the **seventh** literal row `{"readDeadline", readDeadline, 10 * time.Second}` plus a wiring assertion that `New`'s daemon carries a non-zero `d.readDeadline` | `host/daemon/daemon_test.go` (:205-217 table) | 15 | AC3 |
| T2.7 | `TestDaemonReadDeadline/real-store-expired-deadline` — real seeded `:memory:` store, `d.readDeadline` shrunk to 1 ns, all six store-reading routes answer 503 class `Timeout`. **No watchdog, and none claimed** (both arms answer in µs) | `host/daemon/read_deadline_test.go` (new) | 90 | AC3 |
| T2.8 | `blockingStore` — embeds the real store, **overrides ALL FIVE getters** (A2's verbatim fix; a one-getter fake lets five routes 200 instantly and spuriously red the assertion). Every override selects on a test-owned `escape` channel and returns a **distinct sentinel** when it closes. `TestDaemonReadDeadline/blocking-store` drives all six routes table-driven at a 50 ms deadline; 2 s watchdog's red path is `t.Error` **then `close(escape)`** | same file | 130 | AC3 |
| T2.9 | `TestDaemonReadDisconnect` — cancel the request mid-block, `readDeadline` left at its **10 s default**, assert the store call unblocks within 2 s. This is the arm that discriminates `r.Context()` from `context.Background()` | same file | 60 | AC3/AC4 |
| T2.10 | `TestReadCtxCancelledAfterHandler` — `recordingStore` (embeds real store, all five getters record-and-delegate), all six routes serve a normal **200**, and immediately after `ServeHTTP` returns the recorded ctx must already report `context.Canceled`. `readDeadline` at 10 s default. No watchdog needed | same file | 90 | AC3 |
| T2.11 | `TestTimeoutStatusMirrorsSketch` — replay the sketch's `(Timeout(_), 503)` vector on the Go side (the `TestIsLoopbackHostMirrorsSketchPredicate` idiom) | same file | 40 | AC3/AC7 |
| T2.12 | Mutation sweep: MU1, MU2, MU3, MU8, MU9, MU11, MU14; MU10 run locally and **recorded as NOT-CI-red** | — | 0 | §2.2 |

**Watchdog discipline (B2, non-negotiable, every arm).** A watchdog that reds must also **release**
the blocked call and confirm the blocked goroutine exited before cleanup. A bare `t.Error` from a
watchdog goroutine unblocks nothing, and with `SetMaxOpenConns(1)` (V8) a parked getter holds the
sole connection — a deferred `s.Close()` behind it converts a clean red into a **suite-wide hang**.
Layer 1 releases by closing the occupying `*sql.Conn`; layer 3 by `close(escape)`; layers 2 and 4
are measured unable to block and claim no watchdog.

### 3.3 M3 — sanitize + log surface + docs (0.25 d, ~135 LOC) → greens **AC5** and the full set

| # | Task | Files | ~LOC | Closes |
|---|---|---|---|---|
| T3.1 | `internalErrorMessage = "internal store failure"`; `Config.ErrorLog io.Writer` (**nil → `os.Stderr`**); a one-line-per-error writer carrying the route and the verbatim error | `host/daemon/daemon.go`, `host/daemon/handlers.go` | 40 | AC5 |
| T3.2 | Sweep the **six** `Internal` branches at `handlers.go:220, 247, 277, 325, 351, 419` → `internalErrorMessage` + errLog line. **Leave the five BadRequest sites (215, 242, 392, 399, 410) EXACTLY as they are** | `host/daemon/handlers.go` | 15 | AC5 |
| T3.3 | `TestInternalErrorsAreSanitized` — seeds a sentinel-bearing error and asserts **two writes separately**: sentinel **absent** from the response body, sentinel **present** in the ErrorLog buffer. A mutation passing one still reds the other | `host/daemon/read_deadline_test.go` | 70 | AC5 |
| T3.4 | QUICKSTART: one short paragraph in the serve section (reads answer 503 class `Timeout` after 10 s; 500 bodies are generic and the detail is on the daemon's stderr). **QUICKSTART is executed-verbatim-maintained (S7) — re-execute the walkthrough before the commit** | `docs/QUICKSTART.md` (90 lines; serve section at :10-19, 409 demo at :69-76) | 10 | S7 |
| T3.5 | Mutation sweep: MU7 | — | 0 | §2.2 |

**AC5's exact pass condition**: `grep -c 'err.Error()' host/daemon/handlers.go` → **exactly 5**
(11 − 6), and the five survivors are the BadRequest sites. **Over-sanitizing a 400 fails AC5 just
as hard as leaving a 500 echoing.** The 400s echo parse failures of the *client's own input* and
contain no server state; stripping them guts the API's debugging affordance to protect nothing.
The 409 keeps its machine-readable two-head body; the 503's message names only the deadline.

**Never route error text to `announce`.** `ListenAnnouncePrefix` (`daemon.go:140`) is a one-line
protocol whose consumers read exactly one line, and iteration 28 measured extra announce lines
deadlocking `Run` against an `io.Pipe` (V17).

---

## 4. Mutation discipline

### 4.1 Protocol — every arm, no exceptions

```bash
# 0. backup ONCE, before any mutation, OUTSIDE the repo
cp <file> /tmp/w18_bak/<file basename>
shasum -a 256 <file>                      # PRE

# 1. anchor count: assert the token count equals the number in the table BEFORE editing.
#    A count that differs is INSTRUMENT FAILURE: stop, re-derive, do not proceed.
# 2. apply the exact edit.
# 3. shasum -a 256 <file>                 # POST; POST != PRE or the edit DID NOT LAND.
#    A "fix"/"mutation" that never applied prints the same green as one that worked.
# 4. compile the mutant:
#      production-code mutants:  GOTOOLCHAIN=go1.25.6 go build ./...
#      TEST-ONLY mutants:        GOTOOLCHAIN=go1.25.6 go vet ./host/...
#    `go build ./...` DOES NOT COMPILE _test.go AT ALL — mutation-builds on a test-only edit
#    are vacuous. A mutant that does not compile is INSTRUMENT FAILURE, never "survived".
# 5. KILL arm: the -run-scoped command from §4.2, with its expected rc.
#    ASSERT THE `=== RUN` ENUMERATION, NOT rc ALONE. A -run selector matching nothing exits 0
#    with "no tests to run" (V20, observed live at base).
# 6. INVERSE arm: the same package with -skip on the named selector, expected rc=0. This proves
#    YOUR test is the killer and not a bystander.
# 7. restore:  cp /tmp/w18_bak/<basename> <file>     # NEVER `git checkout --`
#    shasum -a 256 <file>  -> MUST equal PRE, byte-identical.
```

### 4.2 The doc's 12 arms, assigned

| ID | Mutation (file : site) | Milestone | `-run` selector that must red | Expected |
|---|---|---|---|---|
| MU1 | `readCtx`: `context.WithTimeout(r.Context(), d.readDeadline)` → `context.WithCancel(r.Context())` | M2 | `TestDaemonReadDeadline/real-store-expired-deadline` | rc≠0; mutant writes 200 + world body — a write the timeout branch can never produce |
| MU2 | `readCtx`: `r.Context()` → `context.Background()` | M2 | `TestDaemonReadDisconnect` | rc≠0; the unblock signal arrives only at the 10 s deadline, after the 2 s watchdog reds **and releases** |
| MU3 | classifier: `if ctx.Err() != nil` → `if false && ctx.Err() != nil` | M2 | `TestDaemonReadDeadline/blocking-store` | rc≠0; body class reads `Internal`, produced only by the 500 branch |
| MU4a–e | per getter: `QueryRowContext(ctx, …)` → `QueryRowContext(context.Background(), …)` at `store.go` **467 / 522 / 551 / 628 / 802** | M1 | `TestReadGettersHonorContext/<getter>` ×5 | rc≠0 each; the getter parks in the pool wait, watchdog reds and closes the occupying conn |
| MU5 | delete the `_pragma=busy_timeout(2000)` injection (`writer_lock.go:176/187` region) | M1 | `TestReadRetriesUnderTransientExclusiveLock` | rc≠0; read fails **instantly** with SQLITE_BUSY instead of returning the row |
| MU6 | injection value `2000` → `0` | M1 | `TestProductionDSNSetsBusyTimeout` | rc≠0; the **PRAGMA readback from the live connection** (driver state, not the Go literal) |
| MU7 | restore `err.Error()` on one Internal branch (`handlers.go:220`) | M3 | `TestInternalErrorsAreSanitized` | rc≠0 on the body assertion; **both** writes asserted separately |
| MU8 | `New`: `readDeadline` field wired to `0` | M2 | `TestBoundedWaitsAndBodyLimit` | rc≠0. **Declared WIRING-ONLY** by the doc (reads a value set alongside its mechanism) and paired with MU1, whose observable is the mechanism's own write. Do not re-classify it as a semantic kill |
| MU9 | Go timeout branch: `503` → `500` | M2 | `TestTimeoutStatusMirrorsSketch` | rc≠0; status checked against the sketch's replayed vector |
| MU10 | sketch arm `Timeout(_) => 503` → `=> 500` | M2 | **NOT CI-RED — declared, not pretended.** Leg 2 tests `world/` only (V12, re-derived §0.1); the sketch's inline rows are outside both gate legs | Run `AILANG_BIN=… ailang test <sketch>` **locally** and record that the 19th row reds. The enforced guard for this drift direction is MU9 |
| MU11 | delete one handler's `defer cancel()` after `readCtx(r)` (any of the six sites) | M2 | `TestReadCtxCancelledAfterHandler` | rc≠0; that route's recorded `ctx.Err()` is nil at assert time. `context.Canceled` has no other writer inside test time with `readDeadline` at 10 s |
| MU12 | add a getter call passing `context.Background()` in production code (scratch site in `handlers.go`) | M1 | `TestNoNewDeadlineFreeStoreReads` | rc≠0; the mutated file's scanned count exceeds its pin |

### 4.3 Two arms this plan ADDS, because three ACs had none (§2.1)

| ID | Mutation | Milestone | Kill command | Expected |
|---|---|---|---|---|
| **MU13** | revert `GetObject`'s signature at `store.go:467` to context-free (drop the `ctx context.Context` parameter and the `ctx` argument to `QueryRowContext`) | M1 | `GOTOOLCHAIN=go1.25.6 go build ./...` | **rc ≠ 0**, with at least one error naming `GetObject`. Then restore and assert rc=0. This is **AC1's and AC8's** only named guard: the compiler, not a grep, is what forbids the context-free spelling once 86 sites pass one |
| **MU14** | insert a deliberate syntax error in `design_docs/sketches/worlddapi.ail` (e.g. delete the closing `)` of the new `Timeout(string)` constructor) | M2 | `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` | **rc ≠ 0**, with Leg 1 naming `worlddapi.ail`. Then restore **byte-identically by sha256** and assert rc=0 with `10 / 39 / 9-of-9`. This is **AC7's** non-vacuity control: it proves Leg 1 actually sweeps the file the sprint edits, i.e. that AC7 **can** fail |

---

## 5. Acceptance commands, exactly as the executor must run them

Every command below is run from the **worktree root** with
`export GOTOOLCHAIN=go1.25.6 AILANG_BIN=/tmp/ailang-v0300/ailang`.

```bash
# AC1  baseline 0/14 (measured at HEAD)  -> pass: first >= 5, control still >= 14
grep -c "context.Context" host/store/store.go; grep -c "func (s \*Store)" host/store/store.go

# AC2  baseline rc=0 with 0 === RUN lines (vacuous)  -> pass: rc=0 AND >= 7 === RUN
go test ./host/store -run 'TestReadGettersHonorContext|TestProductionDSNSetsBusyTimeout|TestReadRetriesUnderTransientExclusiveLock' -v -count=1 > /tmp/ac2.out 2>&1; echo "rc=$?"
grep -c '=== RUN' /tmp/ac2.out

# AC3  baseline rc=0 with exactly 5 === RUN (all from TestBoundedWaitsAndBodyLimit)
#      -> pass: rc=0 AND === RUN lines for ALL FIVE new test functions AND the extended D7 row
go test ./host/daemon -run 'TestDaemonReadDeadline|TestDaemonReadDisconnect|TestReadCtxCancelledAfterHandler|TestInternalErrorsAreSanitized|TestTimeoutStatusMirrorsSketch|TestBoundedWaitsAndBodyLimit' -v -count=1 > /tmp/ac3.out 2>&1; echo "rc=$?"
grep '=== RUN' /tmp/ac3.out

# AC4  baseline 0 and 0 (controls 9 and 8)  -> pass: >= 1 in EACH file
grep -c 'r.Context()' host/daemon/handlers.go; grep -c 'r.Context()' host/daemon/daemon.go

# AC5  baseline 11  -> pass: EXACTLY 5, and they are lines 215/242/392/399/410's descendants
grep -c 'err.Error()' host/daemon/handlers.go; grep -n 'err.Error()' host/daemon/handlers.go

# AC6  baseline 0 (same-scope control: 1 test hit at writer_lock_test.go:609)  -> pass: >= 1
grep -rn 'busy_timeout' host/store/*.go | grep -v _test | wc -l

# AC7  PRE-REGISTERED: rc=0 with 10 / 39 / 9-of-9 (proven under the real edit, plan section 0.3)
./scripts/verify_ail.sh; echo "rc=$?"

# AC8
go build ./... && go test ./... -count=1; echo "rc=$?"

# AC9  baseline rc=0 with 0 === RUN (vacuous)  -> pass: rc=0 AND 1 === RUN AND pins {8,2,1,0}
go test ./host/store -run 'TestNoNewDeadlineFreeStoreReads' -v -count=1 > /tmp/ac9.out 2>&1; echo "rc=$?"
grep -c '=== RUN' /tmp/ac9.out
```

**The trap every test-running AC is written around**: a `-run` selector matching nothing exits **0**
with `testing: warning: no tests to run` — observed live at base for both AC2 and AC9. **Never read
the exit code alone.** Count `=== RUN`.

---

## 6. Estimates and velocity

| Milestone | Impl LOC | Test LOC | Total | Days |
|---|---|---|---|---|
| M1 store layer | ~145 | ~380 | ~525 | **0.50** |
| M2 daemon deadline + status | ~145 | ~465 | ~610 | **0.75** |
| M3 sanitize + log + docs | ~65 | ~70 | ~135 | **0.25** |
| **total** | **~355** | **~915** | **≈1270** | **1.50** |

**Velocity check — the doc's 1.5 d is defensible, and this plan AGREES with it.** Measured repo
references: `VL.B` priced 515 LOC at 0.5 d (**~1030 LOC/day**); `TR.A` priced 2630 LOC at 2 d
(**~1315/day**, and its own planner flagged that as *above* velocity); item 16 (`d9712dd`) landed
**335 Go insertions across 4 files** in a ~0.5 d band. This plan's ≈1270 LOC over 1.5 d is
**≈845 LOC/day** — comfortably inside the measured band and *below* the two references that were
themselves flagged as stretch. The doc's quarter-day uplift over round 1 (for the cancel-release
gate, the widened five-getter fake, and the ratchet) is priced correctly.

**Where the plan's own work reduces risk against the estimate**: §0.3 pre-registers AC7 end-to-end
(the sketch edit's gate behaviour is a known value, not a discovery), and §0.5(i) settles the
five-vs-six seam before an executor stalls on it. Both were open-ended stalls, not 0.05 d lines.

**Test-to-impl ratio is 2.6:1**, high but correct for this item: the item's entire thesis is that
the *existing* gate could not see the layer where the wait happens, so the new gates are the
deliverable. Do not compress the test bands to hit the estimate — compress M2's route table
instead and report.

---

## 7. Execution protocol

- **Worktree**: a **SIBLING of the repo**, never `/tmp` —
  `/Users/voightkampff/dev/sunholo-data/.wt-iter91`. Branch `sprint/w-daemon-read-cancellation`
  off `d31362f`.
- **The executor has NO git write permission.** No `git commit`, `git add`, `git push`,
  `git checkout`, `git stash`, `git restore`, `gh pr`. **THE CONTROLLER MAKES ALL COMMITS**, one
  per milestone, in order.
- **Restores are `cp` from a backup taken outside the repo.** `git checkout --` destroys
  uncommitted work and is forbidden even for mutant cleanup.
- **`git add -A` is forbidden** (§1.3 — the deliberate uncommitted fleet bundle).
- **Every milestone must be independently green** before its commit:
  `go build ./... && go test ./... -count=1` **and** `./scripts/verify_ail.sh`, both with the
  pinned env. Nothing lands red.
- **Never validate `.ail` against a `-dirty` dev build.** The PATH `ailang` is a dev build and is
  acceptable **only** for `ailang messages send`. The gate binary is
  `/tmp/ailang-v0300/ailang` = AILANG v0.30.0, commit `e37b370` (verified live this iteration).
- **`ai-check` is cwd-sensitive.** Run it from the repo root with a repo-root-relative module path.
  The wrong cwd reports `LDR001 module not found` — which reads as a **kill on every arm**. Always
  run a pristine known-positive control first (§0.3 caught exactly this).
- **`GOTOOLCHAIN=go1.25.6` is a BASE CONDITION of this rig, not a regression.** A bare
  `go test ./...` without `AILANG_BIN` exits rc=1 with ~8 failures — also a base condition.
- **Do not run `scripts/verify_go.sh`** (§1.3).

---

## 8. Risks

| # | Risk | Severity | Mitigation |
|---|---|---|---|
| **R1** | Executor invents a **sixth** `readStore` method from §2.5/§2.7's wrong count, most likely `Commit` — widening the seam past the ratified scope and putting the write path behind a "reads" interface | **HIGH** without §0.5(i) | §0.5(i) fixes the count at **five** and names both wrong candidates and the two out-of-seam `ScanUnreadable*` calls |
| **R2** | A watchdog reds without releasing the blocked getter; with `SetMaxOpenConns(1)` the deferred `s.Close()` parks behind it and the whole suite **hangs** instead of failing | **HIGH** | §3.2's B2 rule: every watchdog names its release mechanism; layer 1 closes the occupying `*sql.Conn`, layer 3 closes `escape`; result channels drained under a second bound before cleanup |
| **R3** | The 86-site migration lands a **different** deadline-free count than 11, silently reding AC9 — or worse, an executor "helpfully" threads a real deadline through approve/registry/replay, shrinking the pin | **MEDIUM** | §1.2 names DR-2 as ratified; §3.1's T1.3 trap keeps daemon sites on `readCtx` from M1 so the pin is 11 in **both** milestones; MU12 proves the scanner fires |
| **R4** | `host/broker`'s recorded **~18% base flake** reds M1's `go test ./...` and is attributed to the diff (both broker files change mechanically in M1) | **MEDIUM** | Attribute by shape against the base rate, never by one run. One green baseline run was observed this iteration (50.8 s) and is **not** evidence the flake is gone |
| **R5** | **`D-WORLD-19` is OPEN and its arm A edits `host/store/store.go:467` `GetObject`** — the same function M1 changes (§0.5(ii)) | **MEDIUM** | Does not block: item 18 is ratified and routable, item 17 is parked. If `D-WORLD-19` resolves A mid-flight, item 17's tranche rebases onto item 18's signature — item 18 landed first and its ratification is older. **Controller decision if it resolves before M1 commits** |
| **R6** | The `.ail` edit is applied differently from the probed three anchors and moves a pin | **LOW after §0.3** | The edit is pre-registered with its sha256 and its gate output; **any other totals = STOP and report**, never adjust a pin |
| **R7** | A mutation "passes" because it never landed; or a test-only mutant is compiled with `go build` (which skips `_test.go` entirely) and scored vacuous | **MEDIUM** | §4.1 mandates the sha256 landed-proof at step 3 and `go vet` for test-only mutants at step 4 |
| **R8** | QUICKSTART edited but not re-executed, breaking the S7 executed-verbatim guarantee | **LOW** | T3.4 makes re-execution part of the task, before the M3 commit |

**No risk on this list wants a controller decision before execution**, with one conditional: **R5**,
and only if `D-WORLD-19` is answered while the sprint is in flight.

---

## 9. Handoff

- **Plan**: `design_docs/implemented/w-daemon-read-cancellation-sprint-plan.md`
- **Progress JSON**: `sprint_w-daemon-read-cancellation.json` (repo root)
- **Milestones**: `M1_STORE_CONTEXT` (0.50 d) → `M2_DAEMON_DEADLINE` (0.75 d) → `M3_SANITIZE_AND_DOCS` (0.25 d)
- **Total**: 1.5 days, ≈1270 LOC, **risk medium**
- **Carried forward, out of scope, ratified**: `POST /v1/commit`'s pool wait (DR-1); real deadlines
  in approve/bootstrap/replay (DR-2); `GetVerifyResult`'s context; the store-boundary
  nil/deadline-free reject — all four are the named follow-on
  **`w-bounded-waits-operator-and-write-paths`**, whose four-part acceptance shape is in design doc
  §10 and which carries **two policy questions that return to Mark before it routes** (what bounds
  an attended approval *designed* to wait on a human; `Commit`'s atomicity vs a deadline).
- **Unblocks on landing**: queue item **14** `w-workbench-read-only`, whose declared residual WB-R1
  is discharged here.

**SPRINT_PLAN_PATH**: `design_docs/implemented/w-daemon-read-cancellation-sprint-plan.md`
**SPRINT_JSON_PATH**: `sprint_w-daemon-read-cancellation.json`
