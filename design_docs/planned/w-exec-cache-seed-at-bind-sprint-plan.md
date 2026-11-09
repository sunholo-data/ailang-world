# Sprint plan — w-exec-cache-seed-at-bind (iteration 372, row 257)

- **Design**: [w-exec-cache-seed-at-bind.md](w-exec-cache-seed-at-bind.md), **r3**, including §11. The
  [r3] text is binding: **D-257-1 = A**. Do **not** touch `host/daemon/session_self.go`,
  `cmd/ailang-worldd/episode.go`, `cmd/ailang-worldd/episode_preset.go`. Do **not** add `warmExec`,
  `episodeWarmTimeout` or any GET side effect.
- **Planner base**: `99fc663df6377fa1e767b5fcda6653e52dd79164` (= origin/dev `dc907675` + the design doc),
  worktree `/Users/voightkampff/.ailang/state/world-iter372-row257`, branch `sprint/world-iter372-row257`.
- **Executor**: `claude-haiku-5-5`. Every name, signature, string and assertion below is fixed. Where this
  plan and the design differ, this plan was checked against the tree and wins. The deviations are listed
  in §1.2.
- **Estimate**: about **1.25 days** of executor time across six milestones, about **1,300 LOC**
  (≈400 production, ≈900 test). The design says "+180/−40 production". That was low: §1.3 explains why.
- **Not touched**: `.ail` files, `tools/launchd/*`, shared skills, `docs/QUICKSTART.md`,
  `host/broker/exec_probe.go`, `host/broker/exec_startup.go`, `host/daemon/daemon.go`, the `/v1` route table.

---

## 1. Pre-flight: line references checked against the base

### 1.1 Checked (sha256 prefixes match design V1/V25: `996c5ee11daed4b9`, `1a178f8bc1e39f89`, `6a98846db138b6d3`, `64128681d7bffa13`)

| Design reference | Tree at `99fc663d` | Verdict |
|---|---|---|
| `handlers_exec.go` `prepareCache` `:528-562` | doc comment `:525-527`, `func` `:528`, `seedMu.Lock` `:529`, ends `:562` | correct |
| `handlers_exec.go` `execSeedBudget` `:565-573` | comment `:564`, func `:565-573` | correct |
| `handlers_exec.go` `seedMu` | field at `:141` (`seedMu sync.Mutex`). It is the **only** use of package `sync` in the file | correct; see M1b: drop the `"sync"` import (keep `"sync/atomic"`) |
| `handlers_exec.go` `Execute` → `prepareCache` `:308`, srt `runBounded` `:329` | `:308`, `:329` | correct |
| `handlers_exec.go` `MkdirAll` `:189`, `:532`, `:543` | same | correct |
| `workspace.go` `execHandler` `:468-506` | `func` at `:468`, **ends at `:509`**. `execMu.Lock` `:490`, lookup `:492`, `NewExecHandler` `:496-498`, lazy map `:504-506`, map write `:507`, `return h` `:508` | **wrong end line**: the function and the M2a hunk span `:490-509`, not `:489-506` |
| `workspace.go` `execMu` `:123-126` | comment `:123-124`, `execMu sync.Mutex` `:125`, `execH map[string]episodeTool` `:126` | correct |
| `failure_test.go` seed rows `~:208-218` | "seed-cache exited" `:208-210`, "seed-cache timed out" `:211-217`. `:218` is the time-budget-spent row | **off by one**: `:208-217` |
| `context_roots_test.go` `~:80` | `"host/broker/handlers_exec.go|NewExecHandler|Background": 1,` at `:80` | correct |
| `episode_report.go` `:692-724` | `handleEpisodeReport` `:692`, `readEpisodeStart` `:709`, guard `:720`, NotFound `:721`, `writeJSON` `:724` | correct |
| `newExecFixture` | `host/broker/handlers_exec_m2_test.go:91`; `execFixture` `:34`, `execFixtureOpts` `:43`, `f.call` `:144`, `f.stubCalls` `:185` | (the design gives no line) |
| `coordinator.go:316-317`, `daemon.go:917`, `effectful.go:47,311`, `exec_settings.go:104,115` | as quoted | correct |

The read-only checks `go build ./...` and `go vet ./host/broker/ ./host/daemon/` both pass at the base.

### 1.2 Planner decisions where the design leaves a choice or cannot be built literally

1. **`boundExecHandler` and re-entrancy.** §4 and §8 say that `execHandler`'s cache lookup goes
   "through `boundExecHandler`". `execHandler` already holds `w.execMu` from the lookup through
   construction to the map write (`:490-508`), and `sync.Mutex` is not re-entrant. A literal call
   would deadlock. So the one read of `w.execH` is factored into `execEntryLocked` (caller holds
   `execMu`). Both `execHandler` (under its existing lock) and `boundExecHandler` (which takes the
   lock itself) use it. `execCacheStatus` uses only `boundExecHandler`. This meets the [r3]
   intent: one guarded accessor, and the report never reads `execH` without `execMu`.
2. **`execSeedNext` has no `now` parameter.** §3.7 writes `execSeedNext(state, allVerified, now)`,
   but no branch reads the clock. The "still copying (4.2s)" elapsed time is formatted by a separate
   pure helper `seedStillCopyingDetail`.
3. **After a join, the call returns the joined attempt's own failure.** It does not start another
   attempt. Otherwise one call would run attempts 2 and 3 back to back, which contradicts S3 ("calls 1
   and 2 start attempts 2 and 3").
4. **The "still copying" line names who started the copy**: `since bind` or `since call` (the `by`
   field). The design's example `since bind (4.2s)` is the bind case.
5. **Give-up line when the last attempt succeeded** (the [r3] cap edge: attempts == 3, last succeeded,
   cache unverified again): Detail = `<name>: the published copy is missing or unverified (attempt 3 of 3)`
   plus the give-up suffix.
6. **Cleanup rule for abandoned staging.** Remove any entry of `<state>/exec-cache/<ep>/` whose name
   starts with `.<project>.<name>.seeding-` and is not this attempt's own name. Go in `os.ReadDir`
   order (lexical) and make at most 8 removal attempts. Each failure is logged and does not fail the
   attempt.
7. **Log line formats** (§3.1 says "`... failed: <failure line>`"):
   `ailang-worldd: exec cache for episode "<ep>" (<project>): <names> seeded at <by> in <s>s` and
   `ailang-worldd: exec cache for episode "<ep>" (<project>): seeding at <by> failed: <line>`.
8. **Report `failure` key** is always present (`json:"failure"`, `""` when none), matching §3.6's key list.
9. **S5 has two arms** (call-started and bind-started). Only the call-started arm can kill the mutant
   "the copy is bound to the caller's context", because a bind-started copy has no caller context.
10. **S6 and the edited TX timeout row use the slow seed command** as well as `seedBudget = 1ms`. With
    only a 1 ms budget, a tiny `cp` could finish first, and the test would be flaky.
11. **`TestExecHandlerOwnsNoProcessLifecycle`** also scans `exec_seed.go`. This keeps §4's "the only
    `/bin/cp` runs through `runBounded`" machine-checked. It is a 3-line edit.

### 1.3 Why the LOC estimate is higher than the design's

The design's +180/−40 counts only the state machine. The r2/r3 additions also need code: the marker,
verification, the 7-step attempt, unique staging, bounded cleanup, the status view, the daemon
accessor, the report field and the CLI line. Production is about 260 (`exec_seed.go`) − 45 + 25
(`handlers_exec.go`) + 90 (daemon) + 20 (cmd) ≈ **400**. Fourteen broker tests plus three daemon
tests with their rigs come to about **900** test LOC. Six milestones keep each one at or near ~300 LOC.

---

## 2. Shared definitions (the executor copies these exactly)

### 2.1 `host/broker/exec_seed.go` — names, constants, signatures

```go
package broker

// Row 257 (design_docs/planned/w-exec-cache-seed-at-bind.md r3): the exec
// cache seed is copied on its own clock, started at bind (StartSeed) or by a
// call that finds no copy, never under a call's deadline. A copy is published
// only by an atomic rename of a staging tree that holds World's completion
// marker; readiness is the marker, never existence.

const (
	execSeedBindBudget  = 60 * time.Second // §3.4
	execSeedMaxAttempts = 3                // §3.3
	seedMarkerName      = ".world-seed-complete"
	seedMarkerMax       = 64 // bytes
	seedCleanupMax      = 8  // abandoned staging dirs removed per attempt (§3.1 step 3)
	seedGiveUpSuffix    = "; 3 of 3 attempts failed: fix the seed, then restart serve"
)

// The live seeding states (§3.6).
const (
	ExecSeedCopying = "copying"
	ExecSeedReady   = "ready"
	ExecSeedFailed  = "failed"
	ExecSeedGaveUp  = "gave-up"
)

// seedIncarnation names this daemon process in staging paths: pid and start
// time in nanoseconds (§3.1 step 3).
var seedIncarnation = strconv.Itoa(os.Getpid()) + "-" + strconv.FormatInt(time.Now().UnixNano(), 10)

type execSeed struct {
	mu       sync.Mutex
	attempts int             // copies started, ≤ execSeedMaxAttempts
	running  chan struct{}   // non-nil while a copy runs; closed when it ends
	started  time.Time       // of the running or last attempt
	took     time.Duration   // of the last finished attempt
	by       string          // "bind" or "call": who started the last attempt
	last     *HandlerFailure // nil after a success
	ok       bool            // the last attempt succeeded (or StartSeed found every cache verified)
}

type seedAction int

const (
	seedReady  seedAction = iota // every seeded cache is verified: return nil
	seedJoin                     // a copy is running: wait for it or the call's ctx
	seedStart                    // start attempt attempts+1 in the background, then join
	seedGiveUp                   // the cap is reached: refuse, spawn nothing
)

type seedView struct {
	attempts int
	running  bool
}

// ExecSeedStatus is one seeded cache's live seeding state (row 257 §3.6).
type ExecSeedStatus struct {
	Project  string
	Cache    string
	State    string // ExecSeedCopying | ExecSeedReady | ExecSeedFailed | ExecSeedGaveUp
	Attempts int
	By       string
	Took     time.Duration
	Failure  string // the HandlerFailure line; "" when none
}
```

Pure functions (no I/O, no clock):

```go
func seedMarkerBytes(name string) []byte                          // []byte("world-exec-seed v1 " + name + "\n")
func seedMarkerValid(name string, data []byte) bool               // len(data) <= seedMarkerMax && bytes.Equal(data, seedMarkerBytes(name))
func execSeedNext(s seedView, allVerified bool) seedAction        // §2.3
func seedStateName(attempts int, running, ok, failed bool) string // §2.4
func seedAttemptDetail(name, what string, k int) string           // name + ": " + what + fmt.Sprintf(" (attempt %d of %d)", k, execSeedMaxAttempts)
func seedStillCopyingDetail(name, by string, elapsed time.Duration) string
	// name + ": still copying since " + by + " (" + elapsed.Round(100*time.Millisecond).String() + "); retry"
func seedStagingPrefix(project, name string) string               // "." + project + "." + name + ".seeding-"
func seedStagingName(project, name, incarnation string, k int) string // seedStagingPrefix(project, name) + incarnation + "-" + strconv.Itoa(k)
func defaultSeedCommand(seed, dst string) (string, []string)
	// darwin: "/bin/cp", []string{"-cR", seed, dst}; else "/bin/cp", []string{"-a", "--reflink=auto", seed, dst}  (hostGOOS())
```

Methods (I/O):

```go
func seedCacheVerified(dst, name string) bool        // §2.5
func (h *ExecHandler) seededCaches() []string        // sortedKeys(h.profile.Caches) filtered to Seed != ""
func (h *ExecHandler) seedsVerified() (string, bool) // first unverified seeded name, or ("", true); no seeded caches → ("", true)
func (h *ExecHandler) StartSeed()                    // exported; §2.6
func (h *ExecHandler) startAttemptLocked(by string) chan struct{} // caller holds h.seed.mu
func (h *ExecHandler) runSeedAttempt(k int, by string, done chan struct{}) // the goroutine; the ONE context root
func (h *ExecHandler) seedAttempt(ctx context.Context, k int) *HandlerFailure // every seeded cache in order; first failure ends it
func (h *ExecHandler) seedOne(ctx context.Context, k int, name, seed string) *HandlerFailure // steps 0–6
func (h *ExecHandler) cleanupStaging(epDir, name, own string)
func (h *ExecHandler) awaitSeed(ctx context.Context) error // §2.7, called at the end of prepareCache
func (h *ExecHandler) SeedStatus() []ExecSeedStatus       // §2.8
```

Imports of `exec_seed.go`: `bytes`, `context`, `errors`, `fmt`, `io` (only if used), `os`,
`path/filepath`, `strconv`, `sync`, `time`. It must **not** import `os/exec`, `syscall` or
`host/procbound` (TestExecHandlerOwnsNoProcessLifecycle scans it from M1b).

### 2.2 Failure-line strings (exact; class is always `FailSeedCache` = `seed-cache`)

| Step | `what` (Detail = `seedAttemptDetail(name, what, k)`) | `Err` |
|---|---|---|
| 0 `os.MkdirAll(h.cache, 0o700)` | `cache dir cannot be created` | the MkdirAll error, wrapped `fmt.Errorf("broker: exec handler: seed cache %s: %w", name, err)` |
| 1 the seed has a top-level marker (`os.Lstat(filepath.Join(seed, seedMarkerName))` returns `err == nil`) | `the seed holds World's marker; remove it from the seed` | `fmt.Errorf("broker: exec handler: seed %s holds %s", seed, seedMarkerName)` |
| 2 unverified `dst` cannot be removed | `an unverified copy cannot be removed` | wrapped RemoveAll error |
| 3 `os.Mkdir(staging, 0o700)` fails | `a staging directory cannot be created` | wrapped Mkdir error |
| 4 `runBounded` fails | `seedWhat(err)` (existing, `failure.go:177`) | `fmt.Errorf("broker: exec handler: seed cache %s from %s: %w", name, seed, err)` (today's text) |
| 5 marker `os.WriteFile(…, 0o600)` fails | `the completion marker cannot be written` | wrapped |
| 6 `os.Rename` fails | `the copy cannot be published` | wrapped |

The per-call lines:
- Join timed out: `&HandlerFailure{Class: FailSeedCache, Detail: seedStillCopyingDetail(name, by, time.Since(started)), Err: fmt.Errorf("broker: exec handler: seed cache %s: %w", name, ctx.Err())}`
  → e.g. `seed-cache: GOCACHE: still copying since bind (4.2s); retry`.
- Give-up, last failed: `&HandlerFailure{Class: FailSeedCache, Detail: last.Detail + seedGiveUpSuffix, Err: last.Err}`
  → e.g. `seed-cache: GOCACHE: the copy exited 1 (attempt 3 of 3); 3 of 3 attempts failed: fix the seed, then restart serve`.
- Give-up, last succeeded: Detail `name + ": the published copy is missing or unverified (attempt 3 of 3)" + seedGiveUpSuffix`, `Err: errors.New("broker: exec handler: seed cache " + name + " is unverified after 3 attempts")`.

Expected full lines used by tests:
- `seed-cache: GOCACHE: the copy exited 1 (attempt 1 of 3)` (missing seed, attempt 1)
- `seed-cache: GOCACHE: timed out after 1ms (attempt 1 of 3)`
- `seed-cache: GOCACHE: an unverified copy cannot be removed (attempt 1 of 3)`
- `seed-cache: GOCACHE: the seed holds World's marker; remove it from the seed (attempt 1 of 3)`

### 2.3 `execSeedNext`

```go
switch {
case allVerified:
	return seedReady
case s.running:
	return seedJoin
case s.attempts < execSeedMaxAttempts:
	return seedStart
default:
	return seedGiveUp
}
```

### 2.4 `seedStateName`

```go
switch {
case running:
	return ExecSeedCopying
case ok:
	return ExecSeedReady
case failed && attempts >= execSeedMaxAttempts:
	return ExecSeedGaveUp
case failed:
	return ExecSeedFailed
}
return "" // nothing started
```

### 2.5 `seedCacheVerified(dst, name)`

1. `fi, err := os.Lstat(dst)`. Return false unless `err == nil && fi.IsDir()`.
2. `m := filepath.Join(dst, seedMarkerName)`, then `mi, err := os.Lstat(m)`. Return false unless
   `err == nil && mi.Mode().IsRegular() && mi.Size() <= seedMarkerMax`.
3. `data, err := os.ReadFile(m)`. Return false on error.
4. `return seedMarkerValid(name, data)`.

### 2.6 `StartSeed()`

```go
func (h *ExecHandler) StartSeed() {
	if len(h.seededCaches()) == 0 {
		return
	}
	_, verified := h.seedsVerified() // I/O outside the lock
	h.seed.mu.Lock()
	defer h.seed.mu.Unlock()
	if h.seed.running != nil || h.seed.attempts > 0 || h.seed.ok {
		return // idempotent
	}
	if verified {
		h.seed.ok, h.seed.by = true, "bind" // a restart over a verified cache: no copy
		return
	}
	h.startAttemptLocked("bind")
}
```

`startAttemptLocked(by)`: `h.seed.attempts++`, `k := h.seed.attempts`,
`done := make(chan struct{})`, `h.seed.running, h.seed.started, h.seed.by = done, time.Now(), by`,
`go h.runSeedAttempt(k, by, done)`, `return done`.

`runSeedAttempt(k, by, done)`, in this exact order:
1. `ctx, cancel := context.WithTimeout(context.Background(), h.seedBudget)`. This is the **only**
   `context.Background()` in the file, and it must feed `WithTimeout` directly
   (TestProductionContextRoots).
2. `start := time.Now()`, `f := h.seedAttempt(ctx, k)`, `cancel()`, `took := time.Since(start)`.
3. If `h.log != nil`, write the log line (§1.2 item 7). On success:
   `fmt.Fprintf(h.log, "ailang-worldd: exec cache for episode %q (%s): %s seeded at %s in %.2fs\n", h.episode, h.profile.Project, strings.Join(h.seededCaches(), ","), by, took.Seconds())`.
   On failure:
   `fmt.Fprintf(h.log, "ailang-worldd: exec cache for episode %q (%s): seeding at %s failed: %s\n", h.episode, h.profile.Project, by, f.Line())`.
   (Add `strings` to the imports.)
4. `h.seed.mu.Lock()`. Set `took`, `last = f`, `ok = f == nil`, `running = nil`. Then `h.seed.mu.Unlock()`.
5. `close(done)`.

Step 3 runs before steps 4 and 5, so a test that waits for `done` can read the log buffer without a race.

`seedAttempt(ctx, k)`: `for _, name := range h.seededCaches() { if f := h.seedOne(ctx, k, name, h.profile.Caches[name].Seed); f != nil { return f } }; return nil`.

`seedOne(ctx, k, name, seed)`: let `epDir := filepath.Dir(h.cache)`, `dst := filepath.Join(h.cache, name)`,
`own := seedStagingName(h.profile.Project, name, h.incarnation, k)`, `staging := filepath.Join(epDir, own)`.
- **0** `os.MkdirAll(h.cache, 0o700)`.
- **1** seed marker check (table §2.2).
- **2** If `seedCacheVerified(dst, name)`, return nil (this cache is done). Else if `os.Lstat(dst)` returns
  `err == nil`, then `os.RemoveAll(dst)`, and on error fail as in step 2 of the table.
- **3** `h.cleanupStaging(epDir, name, own)`, then `os.Mkdir(staging, 0o700)`.
- **4** `bin, args := h.seedCommand(seed, filepath.Join(staging, "copy"))`, `h.spawns.Add(1)`,
  `_, err := runBounded(ctx, handlerBounds{execTimeout: h.seedBudget}, handlerCommand{path: bin, args: args, dir: epDir, env: execHostEnv()})`.
- **5** `os.WriteFile(filepath.Join(staging, "copy", seedMarkerName), seedMarkerBytes(name), 0o600)`.
- **6** `os.Rename(filepath.Join(staging, "copy"), dst)`, then `_ = os.RemoveAll(staging)`.
- Every failure from step 3 onward calls `_ = os.RemoveAll(staging)` before it returns. No error path sets `ok`.

`cleanupStaging(epDir, name, own)`: `ents, err := os.ReadDir(epDir)`. If err != nil, return (ENOENT
is normal). `removed := 0`. For each `e` in order: skip it unless
`strings.HasPrefix(e.Name(), seedStagingPrefix(h.profile.Project, name)) && e.Name() != own`. If
`removed == seedCleanupMax`, break. `removed++`. If `err := os.RemoveAll(filepath.Join(epDir, e.Name())); err != nil && h.log != nil`,
log `ailang-worldd: exec cache for episode %q (%s): abandoned staging %s not removed: %v\n`.

### 2.7 `awaitSeed(ctx)`, the per-call path (§3.2)

```go
func (h *ExecHandler) awaitSeed(ctx context.Context) error {
	for {
		name, verified := h.seedsVerified()
		h.seed.mu.Lock()
		act := execSeedNext(seedView{attempts: h.seed.attempts, running: h.seed.running != nil}, verified)
		var done chan struct{}
		switch act {
		case seedReady:
			h.seed.mu.Unlock()
			return nil
		case seedGiveUp:
			f := h.giveUpLocked(name)
			h.seed.mu.Unlock()
			return f
		case seedStart:
			done = h.startAttemptLocked("call")
		case seedJoin:
			done = h.seed.running
		}
		by, started := h.seed.by, h.seed.started
		h.seed.mu.Unlock()
		select {
		case <-done:
		case <-ctx.Done():
			return &HandlerFailure{Class: FailSeedCache, Detail: seedStillCopyingDetail(name, by, time.Since(started)),
				Err: fmt.Errorf("broker: exec handler: seed cache %s: %w", name, ctx.Err())}
		}
		if _, ok := h.seedsVerified(); ok {
			return nil
		}
		h.seed.mu.Lock()
		last := h.seed.last
		h.seed.mu.Unlock()
		if last != nil {
			c := *last
			return &c // the joined attempt's own failure (§1.2 item 3)
		}
		// The joined attempt succeeded but the cache is unverified again: decide again (bounded by the cap).
	}
}
```

`giveUpLocked(name)` returns the give-up failure from §2.2, using `h.seed.last` when it is non-nil.

### 2.8 `SeedStatus()`

Lock `h.seed.mu` and read `attempts, running != nil, ok, last, by, took`, then unlock.
`state := seedStateName(attempts, running, ok, last != nil)`. If `state == ""`, return nil. Otherwise
return one `ExecSeedStatus` per `h.seededCaches()` name, each with `Project: h.profile.Project`, that
`Cache`, and the shared state/attempts/by/took. `Failure` is `last.Line()` when `last != nil`.

---

## 3. Milestones

Each milestone ends with: `gofmt -l <its files>` prints nothing, `go vet <its packages>` passes, and
its gate commands pass. Use `AILANG_BIN=$HOME/.pinned-ailang/ailang` on every `go test`. Set neither
TMPDIR nor GOTMPDIR. The **controller** (not the executor) runs `./scripts/verify_go.sh` and
`./scripts/verify_ail.sh` with `AILANG_BIN=$HOME/.pinned-ailang/ailang` after M2b. No `.ail` changes,
so verify_ail must pass unchanged.

### M1a — pure core (design M1, part 1) · ~70 prod + ~90 test LOC

**Files**: new `host/broker/exec_seed.go` (only the §2.1 constants, `execSeed`, `seedAction`, `seedView`,
`ExecSeedStatus`, `seedIncarnation`, and the pure functions: `seedMarkerBytes`, `seedMarkerValid`,
`execSeedNext`, `seedStateName`, `seedAttemptDetail`, `seedStillCopyingDetail`, `seedStagingPrefix`,
`seedStagingName`, `defaultSeedCommand`). New `host/broker/exec_seed_test.go`.

**Tests**
- **S7 `TestSeedNextIsPure`**: a table over `execSeedNext(seedView{a, r}, v)`:
  (0,false,true)→seedReady; (3,false,true)→seedReady; (1,true,true)→seedReady; (1,true,false)→seedJoin;
  (3,true,false)→seedJoin; (0,false,false)→seedStart; (2,false,false)→seedStart; (3,false,false)→seedGiveUp;
  (4,false,false)→seedGiveUp. A second table over `seedStateName`:
  (1,true,false,false)→copying; (0,false,true,false)→ready; (1,false,true,false)→ready;
  (1,false,false,true)→failed; (2,false,false,true)→failed; (3,false,false,true)→gave-up;
  (0,false,false,false)→"". Also assert `seedAttemptDetail("GOCACHE","the copy exited 1",2) == "GOCACHE: the copy exited 1 (attempt 2 of 3)"`,
  `seedStillCopyingDetail("GOCACHE","bind",4210*time.Millisecond) == "GOCACHE: still copying since bind (4.2s); retry"`,
  and `seedStagingName("proj","GOCACHE","1-2",3) == ".proj.GOCACHE.seeding-1-2-3"`.
  *Kills*: any branch of `execSeedNext` swapped (for example join and start), and the cap edge `<` → `<=`.
- **S12 `TestSeedMarkerValidIsPure`**: `string(seedMarkerBytes("GOCACHE")) == "world-exec-seed v1 GOCACHE\n"`.
  `seedMarkerValid("GOCACHE", x)` is true only for `x = "world-exec-seed v1 GOCACHE\n"`. It is false for:
  `…GOCACHE\nx` (a trailing byte); `…GOCACHE` (no newline); `world-exec-seed v0 GOCACHE\n`;
  `world-exec-seed v2 GOCACHE\n`; `world-exec-seed v1 OTHER\n`; `""`; `nil`; `strings.Repeat("x", 65)`.
  Also: a name of 45 bytes `n` gives `len(seedMarkerBytes(n)) == 65` and `seedMarkerValid(n, seedMarkerBytes(n)) == false` (the size cap).
  *Kills*: `bytes.HasPrefix` for `bytes.Equal`, the size cap dropped, the version not compared.

**Gate**
```
go vet ./host/broker/
AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/broker/ -run '^(TestSeedNextIsPure|TestSeedMarkerValidIsPure)$' -count=1 -v
```

### M1b — state machine, `prepareCache` rewrite, context root (design M1, part 2) · ~210 prod + ~150 test LOC

**Files**
- `host/broker/exec_seed.go`: add `seedCacheVerified`, `seededCaches`, `seedsVerified`, `StartSeed`,
  `startAttemptLocked`, `runSeedAttempt`, `seedAttempt`, `seedOne`, `cleanupStaging`, `awaitSeed`,
  `giveUpLocked`, `SeedStatus`, all exactly as in §2.
- `host/broker/handlers_exec.go`:
  - `ExecHandlerConfig` gains `Log io.Writer // nil = no log; one line per finished seed attempt (row 257)` after `MaxOutputBytes`. Add the `"io"` import.
  - `ExecHandler`: **delete** `seedMu sync.Mutex`. Add `episode string`, `log io.Writer`, `seed execSeed`,
    `seedBudget time.Duration // execSeedBindBudget; tests set it`,
    `seedCommand func(seed, dst string) (string, []string) // defaultSeedCommand; tests set it`, and
    `incarnation string // seedIncarnation; tests set it`. Remove the `"sync"` import (keep `"sync/atomic"`).
  - `NewExecHandler`: in the `h := &ExecHandler{…}` literal, also set `episode: cfg.Episode, log: cfg.Log,
    seedBudget: execSeedBindBudget, seedCommand: defaultSeedCommand, incarnation: seedIncarnation`.
    **Do not** call `StartSeed` here (the probe builds scratch handlers, §1 item 5).
  - Replace `prepareCache` (`:525-562`) with: a doc comment that names row 257; the existing World-dirs
    loop (`home`, `tmp`, `pycache`, `npm`, unchanged lines and `FailExecCache` text); then
    `for _, name := range sortedKeys(h.profile.Caches) { if h.profile.Caches[name].Seed == "" { MkdirAll(dst, 0o700) with the existing FailExecCache "<name>: cannot be created" } }`;
    then `return h.awaitSeed(ctx)`. No lock in `prepareCache`.
  - **Delete** `execSeedBudget` (`:564-573`) and its comment.
- `host/store/context_roots_test.go`: after the `NewExecHandler` row at `:80`, add
  ```go
  // Row 257: the exec cache seed copy runs on its own clock, started at bind
  // inside the ctx-less BinderFor or by a call that found no copy; one
  // attempt is bounded by execSeedBindBudget (60 s) fed directly from this
  // root. /bin/cp only, no Store call.
  "host/broker/exec_seed.go|runSeedAttempt|Background": 1,
  ```
- `host/broker/failure_test.go` (TX rows `:208-217`):
  - "seed-cache exited": want → `"seed-cache: GOCACHE: the copy exited 1 (attempt 1 of 3)"`.
  - "seed-cache timed out": replace `mutate: … expiredCtx{f.h}` with
    `mutate: func(_ *testing.T, f *execFixture) Handler { f.h.seedBudget = time.Millisecond; f.h.seedCommand = slowSeedCommand; return f.h }`,
    and want → `"seed-cache: GOCACHE: timed out after 1ms (attempt 1 of 3)"`. Keep `expiredCtx` (time-budget-spent still uses it).
- `host/broker/handlers_exec_test.go` (`TestExecHandlerOwnsNoProcessLifecycle`, after the handlers_exec.go check at `:88-90`):
  `if hits := forbidden("exec_seed.go"); len(hits) != 0 { t.Fatalf("exec_seed.go owns process-lifecycle code: %v", hits) }`.
- `host/broker/exec_seed_test.go`: add the helpers and S1–S3.

**Test helpers** (in `exec_seed_test.go`):
```go
// slowSeedCommand runs the default copy after 0.5 s, on every OS.
func slowSeedCommand(seed, dst string) (string, []string) {
	bin, args := defaultSeedCommand(seed, dst)
	return "/bin/sh", append([]string{"-c", `sleep 0.5; exec "$0" "$@"`, bin}, args...)
}
// makeSeed writes <dir>/ab/entry = "seeded\n" under a fresh canonicalTempDir and returns the seed dir.
func makeSeed(t *testing.T) string
// seedFixture is newExecFixture with caches.GOCACHE.seed = seed.
func seedFixture(t *testing.T, seed string) *execFixture
// waitSeedIdle waits (≤ 20 s, else t.Fatal) until h.seed.running is nil, re-checking after each close.
func waitSeedIdle(t *testing.T, h *ExecHandler)
// seedState returns the single SeedStatus entry (t.Fatal if len != 1).
func seedState(t *testing.T, h *ExecHandler) ExecSeedStatus
// stagingLeft lists the entries of filepath.Dir(f.cache) whose name contains ".seeding-".
func stagingLeft(t *testing.T, f *execFixture) []string
// invokeLine runs one Workspace.Exec "test a" through a live broker session (as TestExecFailureClassesReachTheRecord does) and returns EffectFailedError.Reason ("" on success).
func invokeLine(t *testing.T, h Handler) string
```

**Tests**
- **S1 `TestSeedAtBindFirstCallSpawnsNoCopy`**: `f := seedFixture(t, makeSeed(t))`. Set `var log bytes.Buffer; f.h.log = &log`.
  Call `f.h.StartSeed()`, then `waitSeedIdle`. Assert `seedState` = {State "ready", Attempts 1, By "bind", Failure ""}.
  Save `n := f.h.Spawns()` (== 1). `f.mustRun(t, "test", "a")`, then assert `f.h.Spawns() == n+1`
  and `f.stubCalls(t) == 1`. Assert `<cache>/GOCACHE/ab/entry` == "seeded\n", the marker bytes ==
  `"world-exec-seed v1 GOCACHE\n"`, `stagingLeft` is empty, and `log.String()` contains
  `exec cache for episode "ep1" (proj): GOCACHE seeded at bind in `.
  *Kills*: `StartSeed` made a no-op (the delta is 2).
- **S2 `TestSeedAtBindFailureReachesTheCall`**: the seed is `filepath.Join(canonicalTempDir(t), "no-such-seed")`.
  Set the log buffer, then `StartSeed`, then `waitSeedIdle`. Assert the state is {failed, 1, bind,
  `seed-cache: GOCACHE: the copy exited 1 (attempt 1 of 3)`}, and the log contains
  `seeding at bind failed: seed-cache: GOCACHE: the copy exited 1 (attempt 1 of 3)`.
  `invokeLine(t, f.h)` == `seed-cache: GOCACHE: the copy exited 1 (attempt 2 of 3)` (the record's
  `Failure` equals it too, as in the TX loop). `f.stubCalls(t) == 0`.
  *Kills*: the per-call path skips the check (`prepareCache` returns nil before `awaitSeed`): srt runs and there is no failure.
- **S3 `TestSeedAttemptsAreCapped`**: missing seed, then `StartSeed` and `waitSeedIdle`. Call 1 is `f.call(t,"test","a")`;
  `errors.As` gives `*HandlerFailure` whose `Line()` == `…the copy exited 1 (attempt 2 of 3)`. Call 2 gives `…(attempt 3 of 3)`.
  Save `n := f.h.Spawns()`. Call 3's `Line()` == `seed-cache: GOCACHE: the copy exited 1 (attempt 3 of 3); 3 of 3 attempts failed: fix the seed, then restart serve`,
  and `f.h.Spawns() == n`. The state is gave-up with 3 attempts.
  *Kills*: the cap removed (call 3 spawns a 4th copy).

**Gate**
```
go vet ./host/broker/ ./host/store/
AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/broker/ -run '^(TestSeedAtBindFirstCallSpawnsNoCopy|TestSeedAtBindFailureReachesTheCall|TestSeedAttemptsAreCapped|TestSeedNextIsPure|TestSeedMarkerValidIsPure|TestExecFailureClassesReachTheRecord|TestExecHandlerOwnsNoProcessLifecycle|TestExecHandlerProjectRootAndCaches|TestFailureLineBoundsAndSanitises)$' -count=1 -v
go test ./host/store/ -run '^TestProductionContextRoots$' -count=1
AILANG_BIN=$HOME/.pinned-ailang/ailang go test -race ./host/broker/ -run '^TestSeed' -count=1
```
`TestProductionContextRoots` is named explicitly because this diff adds a production context root (charter rule).

### M1c — budget, join, partial copy, seed-marker tests (design M1, part 3) · test only, ~170 LOC

**File**: `host/broker/exec_seed_test.go`.

- **S4 `TestSeedPartialCopyIsNeverAccepted`**: the seed is `makeSeed` plus `f2` (content "x") with
  `os.Chmod(f2, 0o000)` (`t.Cleanup` chmods it back to 0o644). Run `StartSeed` and `waitSeedIdle`. The
  state is failed and its Failure starts with `seed-cache: GOCACHE: the copy exited `.
  `os.Lstat(<cache>/GOCACHE)` gives ENOENT and `stagingLeft` is empty. Then plant a "crash" staging dir
  `filepath.Join(filepath.Dir(f.cache), seedStagingName("proj","GOCACHE","dead-1",1), "copy", "partial")`
  as a file. Chmod f2 to 0o644 and call `f.mustRun(t,"test","a")` (attempt 2, by call). Assert the
  planted dir is gone, `seedCacheVerified(<cache>/GOCACHE,"GOCACHE")`, no `partial` file under GOCACHE,
  and `stagingLeft` is empty.
  *Kills*: copying straight into `dst` with no rename and no `RemoveAll` (the partial `dst` exists after attempt 1).
- **S5 `TestSeedOutlivesTheCallThatWaits`**:
  - *Arm "call"* (no `StartSeed`): `f.h.seedCommand = slowSeedCommand`.
    `ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)`, then
    `f.h.Execute(ctx, …"test","a")`. The `*HandlerFailure` Line has prefix
    `seed-cache: GOCACHE: still copying since call (` and suffix `); retry`. Then `waitSeedIdle`, and
    the state is ready with 1 attempt. Save `n`, call `f.mustRun`, then assert `Spawns() == n+1`.
  - *Arm "bind"*: a new fixture with the slow seam. `StartSeed`, then the same 100 ms Execute; the line
    has prefix `seed-cache: GOCACHE: still copying since bind (`. Then `waitSeedIdle` and ready.
  *Kills*: the copy bound to the caller's context (`runSeedAttempt` given the call's ctx as parent:
  the arm-"call" copy dies, the state is failed and the next call copies again, so the delta is 2).
- **S6 `TestSeedHasItsOwnBound`**: set `f.h.seedBudget = time.Millisecond` and `f.h.seedCommand = slowSeedCommand`.
  `StartSeed`, `waitSeedIdle`. The state is failed with Failure == `seed-cache: GOCACHE: timed out after 1ms (attempt 1 of 3)`.
  *Kills*: the bound read from anything but `h.seedBudget` (for example `handlerExecTimeout`: the 0.5 s copy succeeds).
- **S11 `TestSeedHoldingTheMarkerIsRefused`**: `makeSeed`, then write `<seed>/.world-seed-complete` with
  the **valid** bytes `world-exec-seed v1 GOCACHE\n` (so the mutant would publish a cache that verifies).
  Save `n := f.h.Spawns()`, `StartSeed`, `waitSeedIdle`. Failure ==
  `seed-cache: GOCACHE: the seed holds World's marker; remove it from the seed (attempt 1 of 3)`,
  `<cache>/GOCACHE` gives ENOENT, and `Spawns() == n` (no cp ran).
  *Kills*: step 1 removed (the copy is published and is ready).

**Gate**
```
AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/broker/ -run '^(TestSeedPartialCopyIsNeverAccepted|TestSeedOutlivesTheCallThatWaits|TestSeedHasItsOwnBound|TestSeedHoldingTheMarkerIsRefused)$' -count=1 -v
AILANG_BIN=$HOME/.pinned-ailang/ailang go test -race ./host/broker/ -run '^TestSeed' -count=1
```

### M1d — restart, legacy, crash-window and survivor tests (design M1, part 4) · test only, ~280 LOC

**File**: `host/broker/exec_seed_test.go`. Helper: `plantCache(t, dir string, marker []byte, extra string)`
creates `dir/ab/entry` = "planted\n", plus `dir/<extra>` = "x" if extra != "", plus `dir/.world-seed-complete` = marker if marker != nil.

- **S8 `TestStartSeedReusesOnlyAVerifiedCache`**:
  - *Arm 1*: `plantCache(<cache>/GOCACHE, seedMarkerBytes("GOCACHE"), "")`, then `n := Spawns()`,
    `StartSeed()`. With no wait needed, assert `Spawns() == n` and the state is {ready, Attempts 0, By "bind"}.
    `mustRun` gives `Spawns() == n+1`, and `ab/entry` is still "planted\n".
  - *Arm 2*: a new fixture; plant with marker `world-exec-seed v0 GOCACHE\n`. `StartSeed`, `waitSeedIdle`.
    Assert `Spawns() == n+1`, the state is ready with 1 attempt, `ab/entry` == "seeded\n", and the marker is valid.
  *Kills*: `StartSeed` always copies (arm 1 delta 1), and marker bytes not checked (arm 2 keeps "planted").
- **S9 `TestLegacyPartialCacheIsRemovedAndReseeded`**:
  - *Arm bind*: `plantCache(<cache>/GOCACHE, nil, "stray")`, `StartSeed`, `waitSeedIdle`. The state is
    ready, `stray` is gone, the marker is valid, and `ab/entry` == "seeded\n".
  - *Arm call* (a new fixture, no `StartSeed`): the same plant, `n := Spawns()`, `mustRun`. Then
    `Spawns() == n+2` (the copy, then srt), `stray` is gone and the marker is valid.
  *Kills*: readiness by `Lstat` alone (`seedCacheVerified` → `os.Lstat(dst) == nil`: the delta is 1 and the stray file survives).
- **S10 `TestSeedCrashWindowsAreNeverTrusted`**: subtests, each with a new fixture. `os.MkdirAll(filepath.Dir(f.cache))`
  first. `old := filepath.Join(filepath.Dir(f.cache), seedStagingName("proj","GOCACHE","dead-1",1))`.
  - "partial staging": only `old/copy/planted-only` = "x" (no `ab/entry`, no marker; write it directly, not with `plantCache`).
  - "full staging, no marker": `plantCache(old+"/copy", nil, "planted-only")`.
  - "full staging, marker": `plantCache(old+"/copy", seedMarkerBytes("GOCACHE"), "planted-only")`.
  For each of the three: `n := Spawns()`, `StartSeed`, `waitSeedIdle`. Assert `Spawns() == n+1`, `old`
  is absent, GOCACHE is verified with `ab/entry` == "seeded\n", and there is **no** `planted-only` in GOCACHE.
  - "power loss": `plantCache(<cache>/GOCACHE, []byte{}, "")` (an empty marker). `StartSeed`, `waitSeedIdle`:
    1 copy, ready, the marker is valid and `ab/entry` == "seeded\n".
  - "unremovable": `plantCache(<cache>/GOCACHE, nil, "")`, then `os.Chmod(f.cache, 0o500)` (`t.Cleanup` → 0o700).
    `StartSeed`, `waitSeedIdle`. The state is **failed** (never ready) with Failure ==
    `seed-cache: GOCACHE: an unverified copy cannot be removed (attempt 1 of 3)`.
  *Kills*: staging promoted when it holds a marker (the "full staging, marker" case shows `planted-only`),
  and a removal error treated as ready ("unremovable" reports ready, or a different line).
- **S13 `TestSeedRestartIgnoresASurvivingOldCopy`**:
  - *Subtest "survivor"*: `epDir := filepath.Dir(f.cache)`, `MkdirAll(epDir)`.
    `old := filepath.Join(epDir, seedStagingName("proj","GOCACHE","old-1",1))` (computed by the
    production helper, so a fixed-name mutant puts the old process in the new attempt's path). Start
    `exec.Command("/bin/sh","-c", `i=0; while :; do mkdir -p "$0/copy"; : > "$0/copy/old-$i"; i=$((i+1)); sleep 0.01; done`, old)`.
    `t.Cleanup` kills it and Waits. Poll (≤ 5 s) until `old/copy/old-0` exists. `StartSeed`, `waitSeedIdle`.
    The state is ready, GOCACHE is verified, and a walk of GOCACHE finds **no** file whose name has
    prefix `old-`. Then kill and Wait the loop, then `os.MkdirAll(old+"/copy")` and write
    `old/copy/old-final` (residue). Remove `<cache>/GOCACHE/.world-seed-complete` (the agent deleted
    its marker) and call `mustRun` (attempt 2, by call). Assert `old` is absent (Lstat ENOENT) and the marker is valid again.
  - *Subtest "bounded cleanup"*: plant 10 dirs `seedStagingName("proj","GOCACHE", fmt.Sprintf("old%02d",i), 1)` for i=0..9,
    plus `.proj.OTHER.seeding-x-1`. `StartSeed`, `waitSeedIdle`. The remaining `.proj.GOCACHE.seeding-*`
    entries are exactly the names for `old08` and `old09`, and `.proj.OTHER.seeding-x-1` still exists.
    The state is ready.
  *Kills*: a fixed staging path shared across incarnations (old files reach GOCACHE, or `Mkdir` fails),
  cleanup absent (`old` survives), and cleanup unbounded (fewer than 2 remain).

**Gate**
```
AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/broker/ -run '^(TestStartSeedReusesOnlyAVerifiedCache|TestLegacyPartialCacheIsRemovedAndReseeded|TestSeedCrashWindowsAreNeverTrusted|TestSeedRestartIgnoresASurvivingOldCopy)$' -count=1 -v
AILANG_BIN=$HOME/.pinned-ailang/ailang go test -race ./host/broker/ -run '^TestSeed|^TestStartSeed|^TestLegacyPartialCache' -count=1
# design §7's M1 gate, verbatim:
AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/broker/ -run 'Seed|LegacyPartialCache|TestExecFailureClassesReachTheRecord' -count=1
go test ./host/store/ -run '^TestProductionContextRoots$' -count=1
AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/broker/ -count=1 -timeout 900s
```

### M2a — daemon: bind, accessor, report field (design M2, part 1) · ~90 prod + ~190 test LOC

**Files**
- `host/daemon/workspace.go`, **one hunk** in `execHandler` (`:490-509`):
  - `:492` `if cached, ok := w.execH[episodeID]; ok && …` → `if cached, ok := w.execEntryLocked(episodeID); ok && …`.
  - `:498` `OperatorHome: w.exec.operatorHome, MaxOutputBytes: w.exec.maxOutputBytes})` → `…, MaxOutputBytes: w.exec.maxOutputBytes, Log: w.errLog})`.
  - After the `if err != nil {…}` block (`:499-502`) and before `h = eh`, insert:
    ```go
    	// Row 257: the exec cache seed copy starts at bind, on its own clock.
    	eh.StartSeed()
    ```
  - Nothing else in this file changes. #515's hunks are at `:217` and `:266-284`, so there is no overlap (§6).
- New `host/daemon/exec_seed.go`:
  ```go
  package daemon

  // Row 257 (w-exec-cache-seed-at-bind r3 §3.6, §4): the bound exec handler's
  // live seed state for the episode report. Kept out of workspace.go on purpose.

  // execEntryLocked is the one read of w.execH; the caller holds w.execMu.
  func (w *workspaceTools) execEntryLocked(episodeID string) (episodeTool, bool) {
  	t, ok := w.execH[episodeID]
  	return t, ok
  }

  // boundExecHandler is the mutex-protected accessor: the episode's bound
  // *broker.ExecHandler, or nil. It holds execMu only for the map read.
  func (w *workspaceTools) boundExecHandler(episodeID string) *broker.ExecHandler {
  	if w == nil {
  		return nil
  	}
  	w.execMu.Lock()
  	t, ok := w.execEntryLocked(episodeID)
  	w.execMu.Unlock()
  	if !ok {
  		return nil
  	}
  	eh, _ := t.h.(*broker.ExecHandler)
  	return eh
  }

  // execCacheStatus is the report's exec_cache: live daemon state, empty after
  // a restart until the next bind. execMu and the seed's own mutex are never
  // held together.
  func (w *workspaceTools) execCacheStatus(episodeID string) []ReportExecCache {
  	eh := w.boundExecHandler(episodeID)
  	if eh == nil {
  		return nil
  	}
  	var out []ReportExecCache
  	for _, s := range eh.SeedStatus() {
  		out = append(out, ReportExecCache{Project: s.Project, Cache: s.Cache, State: s.State, Attempts: s.Attempts,
  			By: s.By, TookMS: s.Took.Milliseconds(), Failure: s.Failure})
  	}
  	return out
  }
  ```
- `host/daemon/episode_report.go`:
  - In `EpisodeReport` (after `Notes`, `:67`):
    ```go
    	// ExecCache (row 257) is the bound exec handler's seed state, one entry
    	// per seeded cache. It is LIVE daemon state, not store state: empty after
    	// a daemon restart until the episode's next bind.
    	ExecCache []ReportExecCache `json:"exec_cache,omitempty"`
    ```
  - A new type after `ReportEffect`:
    ```go
    // ReportExecCache is one seeded exec cache's live state (row 257 §3.6).
    type ReportExecCache struct {
    	Project  string `json:"project"`
    	Cache    string `json:"cache"`
    	State    string `json:"state"` // copying | ready | failed | gave-up
    	Attempts int    `json:"attempts"`
    	By       string `json:"by"` // bind | call
    	TookMS   int64  `json:"took_ms"`
    	Failure  string `json:"failure"`
    }
    ```
  - `handleEpisodeReport`: after the `rep.Finish, err = …` block (`:715-719`) and before the guard, insert
    `rep.ExecCache = d.workspace.execCacheStatus(episode)`. Change the guard at `:720` to
    `if !ok && start == nil && rep.Finish == nil && len(rep.ExecCache) == 0 {`.
- New `host/daemon/exec_seed_test.go`, with a rig built from `health_exec_tool_status_test.go:18-63`
  (fake srt tree, `broker.ArchiveExecSandbox`, profile JSON):
  ```go
  // newSeedDaemon: canonical base; <base>/state, <base>/workspace/<ep> for each ep, <base>/home; a profile
  // "proj" whose caches.GOCACHE.seed = seed; d := newHandlerDaemon(t); then
  // d.workspace.root, d.workspace.stateDir = root, stateDir; d.workspace.bindingFor = nil;
  // d.workspace.exec = &workspaceExec{profile: p, byProject: map[string]*broker.ExecProfile{"proj": p}, sandbox: sb, operatorHome: home}.
  // A t.Cleanup registered after the dirs exist waits (≤ 20 s) until no bound handler's SeedStatus is "copying".
  func newSeedDaemon(t *testing.T, seed string, episodes ...string) (d *Daemon, root, stateDir string)
  func waitDaemonSeed(t *testing.T, d *Daemon, ep string) []ReportExecCache // polls execCacheStatus every 10 ms, ≤ 20 s, until no entry is "copying"
  ```

**Tests**
- **D1 `TestEpisodeReportShowsBindSeedFailure`**: seed `<base>/no-such-seed`, `ep1`.
  `d.workspace.execHandler("ep1", filepath.Join(root,"ep1"))`, then `waitDaemonSeed`.
  `rec := requestRecorder(t, d, http.MethodGet, "/v1/episodes/ep1/report", nil)` gives 200 with
  `total_calls` 0, and `exec_cache` == `[{project:"proj", cache:"GOCACHE", state:"failed", attempts:1, by:"bind", took_ms:<any>, failure:"seed-cache: GOCACHE: the copy exited 1 (attempt 1 of 3)"}]`.
  *Kills*: the report omits the field (404 NotFound), and `eh.StartSeed()` removed from workspace.go (no entries, so 404).
- **D3 `TestBindSeedsAFreshEpisode`**: a seed with `ab/entry` = "seeded\n". Assert `os.Lstat(<state>/exec-cache)`
  gives ENOENT. Then `execHandler("ep1", …)` returns a `*broker.ExecHandler`. `waitDaemonSeed` gives one
  entry {ready, 1, bind}. `<state>/exec-cache/ep1/proj/GOCACHE/ab/entry` == "seeded\n", and
  `…/GOCACHE/.world-seed-complete` == `world-exec-seed v1 GOCACHE\n`.
  *Kills*: §3.1 step 0 removed (the staging Mkdir fails ENOENT and the state is failed).
- **D4 `TestEpisodeReportRacesFirstDispatch`**: episodes `ep1`…`ep8`, a small seed. Use a start channel
  and a `sync.WaitGroup`. Goroutine A calls `d.workspace.execHandler(ep, root/ep)` for each ep in order
  (what `registry()` calls on a first dispatch). Goroutine B, for 8 rounds, GETs `/v1/episodes/<ep>/report`
  for each ep. Each response must be 404 with class NotFound, or 200 whose `exec_cache` has length 0
  or 1 and, if 1, `cache == "GOCACHE"` and `state` in {copying, ready}. Failures go to a channel; there
  is no `t.Fatal` inside goroutines. Close start, Wait (bounded: `select` on a done channel against
  `time.After(60*time.Second)`), then `waitDaemonSeed` for each ep.
  *Kills (under `-race`)*: `execCacheStatus` reading `w.execH[episodeID]` without `execMu`.

**Gate**
```
go vet ./host/daemon/
AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/daemon/ -run '^(TestEpisodeReportShowsBindSeedFailure|TestBindSeedsAFreshEpisode|TestEpisodeReportRacesFirstDispatch)$' -count=1 -v
AILANG_BIN=$HOME/.pinned-ailang/ailang go test -race ./host/daemon/ -run 'ReportRacesFirstDispatch' -count=1
AILANG_BIN=$HOME/.pinned-ailang/ailang go test -race ./host/daemon/ -run '^(TestEpisodeReportShowsBindSeedFailure|TestBindSeedsAFreshEpisode)$' -count=1
AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/daemon/ -count=1 -timeout 900s
```

### M2b — CLI text line (design M2, part 2) · ~20 prod + ~35 test LOC

**Files**
- `cmd/ailang-worldd/episode_report.go`:
  ```go
  // execCacheLine is one exec_cache entry of the report table (row 257).
  func execCacheLine(e daemon.ReportExecCache) string {
  	line := fmt.Sprintf("exec cache   %s %s %s", e.Project, e.Cache, e.State)
  	if e.Attempts > 0 {
  		line += fmt.Sprintf(" at %s (attempt %d of 3)", e.By, e.Attempts)
  	}
  	if e.State == "ready" && e.Attempts > 0 {
  		line += fmt.Sprintf(" in %dms", e.TookMS)
  	}
  	if e.Failure != "" {
  		line += ": " + e.Failure
  	}
  	return line
  }
  ```
  In `printReport`, as the **last** statement (after the `Uncommitted` block):
  `for _, e := range r.ExecCache { fmt.Fprintln(w, execCacheLine(e)) }`.
- `cmd/ailang-worldd/episode_report_test.go`: **C1 `TestEpisodeReportPrintsExecCache`**.
  - (a) `rep := sampleReport()`, `rep.ExecCache = []daemon.ReportExecCache{{Project:"ailang-world", Cache:"GOCACHE", State:"failed", Attempts:1, By:"bind", Failure:"seed-cache: GOCACHE: timed out after 1m0s (attempt 1 of 3)"}}`.
    `printReport` output contains the line
    `exec cache   ailang-world GOCACHE failed at bind (attempt 1 of 3): seed-cache: GOCACHE: timed out after 1m0s (attempt 1 of 3)\n`.
  - (b) {State "ready", Attempts 1, By "call", TookMS 380} gives `exec cache   ailang-world GOCACHE ready at call (attempt 1 of 3) in 380ms`.
  - (c) {State "ready", Attempts 0, By "bind"} gives exactly `exec cache   ailang-world GOCACHE ready`.
  - (d) `printReport(sampleReport())` contains no `exec cache` (nothing new without entries).
  *Kills*: the renderer drops the line.

**Gate**
```
go vet ./cmd/ailang-worldd/
AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./cmd/ailang-worldd/ -run '^(TestEpisodeReportPrintsExecCache|TestEpisodeReportTablePrintsEveryClass|TestEpisodeReportTablePrintsTheReadsLine)$' -count=1 -v
# design §7's M2 gates, verbatim:
AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/daemon/ ./cmd/ailang-worldd/ -run 'ExecCache|BindSeed|ReportRacesFirstDispatch' -count=1
AILANG_BIN=$HOME/.pinned-ailang/ailang go test -race ./host/daemon/ -run 'ReportRacesFirstDispatch' -count=1
```

### Controller-only, after M2b
```
go build ./... && go vet ./host/broker/ ./host/daemon/ ./host/store/ ./cmd/ailang-worldd/
go test ./host/store/ -run '^TestProductionContextRoots$' -count=1
AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh
AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh
git status --short   # only the files in §5
```

---

## 4. Mutation table (controller and judge run these; apply one at a time, run the named test, expect FAIL, restore, `cmp`)

| Id | Mutant (single edit) | Must fail |
|---|---|---|
| MUT-START-NOOP | `StartSeed` body → `return` | S1, D1, D3 |
| MUT-SKIP-CHECK | `prepareCache`: `return h.awaitSeed(ctx)` → `return nil` | S2 |
| MUT-NO-CAP | `execSeedNext`: `s.attempts < execSeedMaxAttempts` → `true` | S3, S7 |
| MUT-DIRECT-DST | `seedOne` step 4 copies to `dst` (not `<staging>/copy`); steps 5–6 skipped; no `RemoveAll` on error | S4 |
| MUT-CALLER-CTX | `runSeedAttempt` takes a `parent context.Context` used in place of `context.Background()`; `awaitSeed` passes `ctx`, `StartSeed` passes `context.Background()` | S5 (arm call) |
| MUT-BUDGET | `runSeedAttempt`/`seedOne`: `h.seedBudget` → `handlerExecTimeout` | S6 |
| MUT-SWAP | `execSeedNext`: the `seedJoin` and `seedStart` results swapped | S7 |
| MUT-ALWAYS-COPY | `StartSeed`: `if verified {…}` block deleted | S8 arm 1 |
| MUT-MARKER-PREFIX | `seedMarkerValid`: `bytes.Equal` → `bytes.HasPrefix(data, []byte("world-exec-seed v"))` | S8 arm 2, S12 |
| MUT-LSTAT-READY | `seedCacheVerified` → `_, err := os.Lstat(dst); return err == nil` | S9, S10 power loss |
| MUT-PROMOTE | `cleanupStaging`: an entry whose `copy/` holds a valid marker is renamed to `dst` instead of removed | S10 full staging, marker |
| MUT-REMOVE-OK | `seedOne` step 2: a RemoveAll error → `return nil` | S10 unremovable |
| MUT-NO-STEP1 | `seedOne` step 1 deleted | S11 |
| MUT-NO-SIZECAP | `seedMarkerValid`: `len(data) <= seedMarkerMax &&` removed | S12 |
| MUT-FIXED-STAGING | `seedStagingName` ignores `incarnation` and `k` (returns the prefix minus the trailing `-`) | S13 survivor |
| MUT-NO-CLEANUP | `cleanupStaging` body → `return` | S13 survivor |
| MUT-UNBOUNDED | `cleanupStaging`: the `removed == seedCleanupMax` break deleted | S13 bounded cleanup |
| MUT-NO-REPORT | `handleEpisodeReport`: the `rep.ExecCache = …` line deleted | D1 |
| MUT-NO-STEP0 | `seedOne` step 0 deleted | D3 |
| MUT-UNLOCKED | `execCacheStatus`: `w.boundExecHandler(episodeID)` → a direct `w.execH[episodeID]` read with no lock | D4 under `-race` |
| MUT-NO-LINE | `printReport`: the `ExecCache` loop deleted | C1 |
| MUT-NO-BIND | `workspace.go`: `eh.StartSeed()` deleted | D1, D3 |

---

## 5. Files (final `git status --short` after M2b)

| File | Milestone | Change |
|---|---|---|
| `host/broker/exec_seed.go` | M1a, M1b | new |
| `host/broker/exec_seed_test.go` | M1a–M1d | new: S1–S13 |
| `host/broker/handlers_exec.go` | M1b | `Log`, new fields, `seedMu`/`sync` import removed, `prepareCache` rewritten, `execSeedBudget` removed |
| `host/broker/failure_test.go` | M1b | two TX rows edited |
| `host/broker/handlers_exec_test.go` | M1b | the guard also scans `exec_seed.go` |
| `host/store/context_roots_test.go` | M1b | +1 row |
| `host/daemon/workspace.go` | M2a | one hunk `:490-509` |
| `host/daemon/exec_seed.go` | M2a | new: `execEntryLocked`, `boundExecHandler`, `execCacheStatus` |
| `host/daemon/episode_report.go` | M2a | `ExecCache`, `ReportExecCache`, fill, guard |
| `host/daemon/exec_seed_test.go` | M2a | new: D1, D3, D4 |
| `cmd/ailang-worldd/episode_report.go` | M2b | `execCacheLine`, print loop |
| `cmd/ailang-worldd/episode_report_test.go` | M2b | C1 |

---

## 6. #515 overlap

Row 256 (draft PR #515, attended) changes `host/daemon/workspace.go` only at `@@ -217,7` and
`@@ -266,11 +266,21` (design V13). This sprint's single hunk is in `execHandler` at `:490-509`, more than
200 lines away, so a rebase in either order applies cleanly. If it ever conflicts, keep both: #515
changes exclusion logging and this sprint changes handler construction. No other file is shared
(#515 also touches `exec_probe.go`, `exec_startup.go`, `daemon.go`, `doctor*.go`, QUICKSTART, which
this sprint does not).

---

## 7. Risks (design §6, plus what planning found)

- **R1** The incident load (110–130) is UNMEASURED. If a clone needs more than 60 s, three attempts
  time out and exec is refused with a line that says why (fail-closed).
- **R2** Linux full copies (no reflink) are UNMEASURED. CI exercises only the small test seeds.
- **R3** The cap of 3 is a judgement.
- **R4** A `cp` that survives a killed daemon writes only into its own incarnation's staging. Cleanup
  takes at most 8 directories per attempt, and the design claims no one-copy disk bound.
- **R5** `exec_cache` is live state, empty after a restart until the next bind.
- **R7/R8/R9** As in the design: an agent spends an attempt by deleting its marker; one extra dotfile
  in GOCACHE; a seed that holds the marker is refused.
- **P1 (planning)** The marker cap of 64 bytes means a cache name over 44 bytes can never verify. Such
  a cache is re-seeded until it gives up (fail-closed, but the line says "fix the seed"). The cache-name
  grammar (`execEnvNamePattern`) has no length bound. Real names are short (`GOCACHE`). Not fixed here; see the report.
- **P2 (planning)** If the seed's top directory has no owner-write bit, `cp` preserves the mode and the
  marker write fails with `the completion marker cannot be written` (fail-closed). Today's code had no
  write into the copy, so this is a new way to fail on read-only seeds.
- **P3 (planning)** Same-process cleanup can remove a staging dir of a *second* handler for the same
  episode × project. Today that cannot happen: one handler per episode × project, and profiles are
  immutable after startup. If it did, the exclusive `Mkdir` or the publish would fail the attempt
  closed, never publish a mix.
- **P4 (timing)** S5 and S6 rely on a 0.5 s `sleep` seam against 100 ms and 1 ms bounds, which leaves
  wide margins. D4 is bounded at 60 s.

---

## 8. Acceptance criteria (the judge scores against these)

1. `host/broker/exec_seed.go` exists with every §2.1 name and signature. `execSeedBindBudget == 60*time.Second`, `execSeedMaxAttempts == 3`.
2. `prepareCache` runs no `/bin/cp` and holds no lock. The seed copy's only `runBounded` call is in `seedOne`, reached only from `runSeedAttempt`, under `context.WithTimeout(context.Background(), h.seedBudget)`.
3. `seedMu` and `execSeedBudget` are gone, and `NewExecHandler` does not call `StartSeed`.
4. Readiness is the marker: `seedCacheVerified` requires a directory and a regular marker file of ≤ 64 bytes equal to `world-exec-seed v1 <name>\n`.
5. Each attempt copies into an exclusively created staging dir `.<project>.<name>.seeding-<incarnation>-<k>`. It writes the marker only after `cp` exits 0 and publishes only by `os.Rename`. It removes at most 8 abandoned staging dirs.
6. Every seed failure is a `*HandlerFailure{Class: "seed-cache"}` whose line matches §2.2 exactly, and the failure class set stays at 10.
7. Attempts are capped at 3. The fourth request refuses with the give-up suffix and spawns nothing.
8. A call whose context ends while a copy runs refuses with `still copying since <bind|call> (<d>); retry`, and the copy goes on.
9. `TestProductionContextRoots` passes with exactly one new row `"host/broker/exec_seed.go|runSeedAttempt|Background": 1`.
10. `workspace.go` changes in one hunk inside `execHandler` (lookup via `execEntryLocked`, `Log: w.errLog`, `eh.StartSeed()`). No change to `session_self.go`, `episode.go`, `episode_preset.go`, `daemon.go`, `exec_probe.go` or `exec_startup.go`.
11. `execCacheStatus` reads the handler map only through `boundExecHandler` (under `execMu`) and never holds `execMu` and the seed mutex together. D4 passes under `-race`.
12. `GET /v1/episodes/{ep}/report` carries `exec_cache` (omitted when empty). The NotFound guard passes when `exec_cache` is non-empty, and the CLI prints one `exec cache   …` line per entry and nothing when there are none.
13. Tests S1–S13, D1, D3, D4 and C1 exist under the names in §3 and pass. The two TX rows are edited as specified.
14. Each §4 mutant, applied alone, makes its named test(s) fail.
15. `gofmt -l` is clean on every touched file. `go build ./...`, `go vet` on touched packages, `./scripts/verify_go.sh` and `./scripts/verify_ail.sh` (pinned `AILANG_BIN`) pass.
16. `git status --short` lists only the §5 files. No `.ail`, `tools/launchd/*`, QUICKSTART or skill changes.
