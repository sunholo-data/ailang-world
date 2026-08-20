# Sprint plan — `w-archive-stderr-in-manifest` (queue row 21, clause-2)

**Item**: queue row **21** of `design_docs/world-mission.md` — *stderr in the version manifest:
the first PERSISTED instance of the stderr-merge class*. Locate the row with
`grep -nE '^21\. ' design_docs/world-mission.md`, never by line number (the charter prepends a
STATUS stamp per iteration; at planning time the row printed at `:3530`, and that number will
rot).
**Design doc (specification of record)**:
[`design_docs/planned/w-archive-stderr-in-manifest.md`](w-archive-stderr-in-manifest.md) — 808
lines, quorum-clean over 2 rounds + one recovered-reviewer solo re-run, doc commit `6a811e1`.
**Status**: PLANNED · **NO SPLIT** · Milestone A is the whole item · **6 tasks / 1 commit**.
**Base**: branch `dev` @ **`7d79ad3`**, working tree clean (`git status --porcelain` empty, V0).
The doc was authored against `b5ddf0e`; the two commits since are `design_docs/`-only, and every
premise below was **re-measured first-party at `7d79ad3`** in this planning session.
**Pins, on every command**: `export PATH=/opt/homebrew/bin:$PATH`,
`GOTOOLCHAIN=go1.25.6` (ambient `go1.26.4` is DENIED by `scripts/verify_go.sh:119` and reds
`host/store/toolchain_canary_test.go`), `AILANG_BIN=/tmp/ailang-v0300/ailang` (v0.30.0).
**Shell is zsh**: no `${PIPESTATUS[0]}` (bash-only, expands EMPTY); quote every `--include='*.go'`;
brace every `"${rev}:host/…"` (`:h`/`:p` are zsh history modifiers and `host/` starts with one).
**Headline price: 1.0 day — the top of the doc's own 0.5–1 d band, not above it.** The driver is
not the code delta (~55 lines). It is that **two of the doc's specified test mechanisms do not
work as written on this rig** (§0.2 C1, C2, each measured 3/3) and that **the doc's own
load-bearing round-2 decision — the CONDITIONAL self-heal — has no acceptance criterion and no
mutation** (§0.2 C4). §2 prices it; §6 refuses the split.

---

## 0. Planner's first-party verification — what survived and what was REFUTED

### 0.1 Design-doc premises re-measured at `7d79ad3` — all CONFIRMED

| Premise | Command (run this session) | Observed | Verdict |
|---|---|---|---|
| P1 stream split | `/tmp/ailang-v0300/ailang --version >/tmp/i98_stdout.txt 2>/tmp/i98_stderr.txt; echo rc=$?; wc -c < each` | rc=0; stdout **168 B**, line 1 `AILANG v0.30.0`; stderr **63 B** = `2026/08/20 01:13:12 Observatory: 314MB (warn threshold: 200MB)` | **CONFIRMED**, and the drift is re-witnessed: 301 MB (18th) → 308/309 MB (19th) → **314 MB** (20th). Separate files, never `2>&1`. |
| P2 the probe merges | `sed -n '384p;391p' host/archive/archive.go` | `:384 cmd := exec.Command(execPath, "--version")`; `:391 out, err := cmd.CombinedOutput()` | **CONFIRMED, exact lines** |
| P6 five non-test `exec.Command*` sites | `grep -rn --include='*.go' -E 'exec\.Command(Context)?\(' host/ cmd/ \| grep -v '_test.go'` | exactly 5: `broker/handlers.go:93`, `archive/archive.go:384`, `capsule/capsule.go:154`, `pkgproj/pkgproj.go:213`, `replay/replay.go:327`. Zero in `cmd/`. | **CONFIRMED** |
| P6b/P6c closure + control | `grep -rn --include='*.go' '"os/exec"' host/ cmd/ tools/ \| grep -v '_test.go'`; and `grep -rn --include='*.go' -E 'os\.StartProcess\|syscall\.Exec\b\|exec\.Command' host/ cmd/ \| grep -v '_test.go'` | same 5 files import `os/exec`; the alternation call returns exactly the 5 `exec.Command` lines — **the control alternation fires on all five in the same invocation over the same scope**, so the two zeros are measured zeros | **CONFIRMED** |
| P8 pkgproj merges | `sed -n '213p;219p;221p' host/pkgproj/pkgproj.go` | `:219 out, err := cmd.CombinedOutput()`, `:221` interpolates merged `out` | **CONFIRMED** |
| P9b the one polluted artifact | `ls -ld /tmp` (symlink → `private/tmp`); `find /private/tmp -maxdepth 6 -name '*.artifacts' -type d`; control `test -d /private/tmp/world-demo.db.artifacts` | 1 tree; control FIRES; two files under `…/sha256/e9746fef…3fb5/`; `ls /private/tmp/world-demo.db` → **No such file** | **CONFIRMED** — REAL path, per the doc's own recorded `find`-declines-a-symlinked-root instrument failure |
| P9b the manifest is polluted | `python3 -c "import json;print(repr(json.load(open('…/manifest.json'))['version']))"` | `'2026/08/18 21:02:42 Observatory: 301MB (warn threshold: 200MB)\nAILANG v0.30.0\nCommit: e37b370\nFull:   e37b370d1d7a9c4e7136b319e38bec4d5f2bd9a0\nBuilt:  2026-07-19T09:27:00Z\n\nThe AI-First Programming Language\nCopyright (c) 2025-2026\n'` | **CONFIRMED, byte-verbatim** — this is the AC4a fixture string |
| P10 idempotent path never probes | read `archive.go:249-266` | `:266 return ref, nil` under the comment "Identical bytes already archived: idempotent no-op success", **before** Step 5's probe at `:289` | **CONFIRMED** |
| P11 base gates, 3 packages | `GOTOOLCHAIN=go1.25.6 AILANG_BIN=… go build ./...` then `go test ./host/archive/... ./host/pkgproj/... ./host/daemon/... -count=1` | build **rc=0**; `ok host/archive 2.368s`, `ok host/pkgproj 0.534s`, `ok host/daemon 3.232s` | **CONFIRMED at `7d79ad3`** |
| P11b `verify_ail.sh` | `AILANG_BIN=… ./scripts/verify_ail.sh` | **rc=0** — `world package gate PASSED: 9/9 steps performed non-zero work`, `verify gate PASSED: 10 required identities verified, 39 named tests pass` | **CONFIRMED GREEN at `7d79ad3`** |
| P12 fixtures stay green | read `archive_test.go:22-32`, `daemon_test.go:37-41`, `:583` | both `fakeInterpreter` helpers write the version to **stdout only** (`printf '%s' '<version>'`), no stderr write anywhere in either helper; `daemon_test.go:583` asserts verbatim equality | **CONFIRMED** |
| P13 pkgproj exec seam untested | `grep -n 'CrossCheck\|func Test' host/pkgproj/pkgproj_test.go` | 3 test funcs (`:14`, `:32`, `:65`), **zero** `CrossCheck` references | **CONFIRMED** |
| P17 zero bounded exec in `host/archive`, control fires | `grep -n -E 'context\|Timeout\|CommandContext' host/archive/archive.go host/capsule/capsule.go host/replay/replay.go host/broker/handlers.go` — **check and control in ONE call over ONE scope** | `archive.go`: exactly **1** hit, `:56`, prose in a comment. Control fires hard: `replay.go:49/51/323/327` (`execTimeout = 60s`), `broker/handlers.go:14/90/93/127` (`handlerExecTimeout = 30s`), `capsule.go:29/117/152/154` (`capsuleExecTimeout = 60s`) | **CONFIRMED** |
| P18 `store.Open` creates | `grep -n 'Open opens' host/store/store.go` | `:218 // Open opens (or creates) the SQLite database at path` | **CONFIRMED** |
| `dryRunLine` shape T3 must match | `sed -n '181p' host/pkgproj/pkgproj.go` | `^  (Tarball\|Content hash\|Interface hash): (?:([0-9]+) bytes \()?((?:sha256:)[0-9a-f]{17})\.\.\.?\)?$` — **exactly 17 hex**, two leading spaces | **CONFIRMED**; the doc's proposed stderr line `  Tarball: 999 bytes (sha256:<17 hex>...)` matches |

### 0.2 EIGHT corrections to the design doc — each MEASURED, each changing the plan

#### C1 — **REFUTED, load-bearing: T2's deadline arm as worded FAILS against a correct implementation.**

The doc (§Decision 3, T2 fourth arm; AC7) prescribes: *"point the archive at a fake interpreter
whose script **sleeps far past the bound (e.g. 10 s, 50×) before answering**, call `Archive()`,
and assert (a) it returns within a deterministic upper bound on the measured wall clock
(e.g. **< 5 s**)"*.

Prototyped end-to-end this session (`/tmp/i98probe`, pinned `go1.25.6`) with **exactly** the
shipped shape — `exec.CommandContext(ctx, path, "--version")`, `cmd.Stdout/Stderr =
&bytes.Buffer`, 200 ms deadline — against a `#!/bin/sh` fake running `sleep 10; printf …`:

```
A_sleep_then_answer_noWaitDelay      elapsed=10.351s  deadlineExceeded=true  err=... (orig signal: killed)
A_sleep_then_answer_noWaitDelay      elapsed=10.155s  deadlineExceeded=true   [rerun 2]
A_sleep_then_answer_noWaitDelay      elapsed=10.153s  deadlineExceeded=true   [rerun 3]
```

**3/3: `cmd.Run()` returns after 10.15–10.35 s under a 200 ms deadline.** Cause:
`exec.CommandContext` SIGKILLs only the **direct child** (`sh`); the `sleep` **grandchild**
inherits the write end of the stdout pipe, and `cmd.Wait()` blocks on the output-copy goroutines
until that grandchild exits. Assertion (b) (`errors.Is(err, context.DeadlineExceeded)`) is TRUE
throughout — so **assertion (a) reds against the CORRECT implementation**. Written as the doc
words it, AC7 is not a tooth, it is a **false red the executor cannot make green**.

Two remedies measured in the same program:

```
B_exec_sleep_noWaitDelay             elapsed=203ms   deadlineExceeded=true   [fixture-side: `exec sleep 10`, no grandchild]
A_sleep_then_answer_WaitDelay100ms   elapsed=303ms   deadlineExceeded=true   [production-side: cmd.WaitDelay]
```

**Adopted: the fixture-side remedy (B).** The blocking fake becomes
`if [ "$1" = "--version" ]; then exec /bin/sleep 20; fi` — `exec` replaces the shell, so the
process the context kills IS the sleeper and no grandchild holds the pipe. Absolute `/bin/sleep`
because `cmd.Env = childenv.Scrubbed(os.Environ())` and a PATH-less lookup is one more thing that
can fail for the wrong reason. **Consequence the doc's wording must lose**: the fake can no longer
"answer after the sleep", so the unbounded-mutation red is *"waits the full 20 s, then returns
rc=0 with EMPTY stdout"* — still a deterministic assertion-red (wall bound + `DeadlineExceeded`
both fail), still not a hang, but the doc's phrase "succeeds after the fake's finite sleep" is
replaced by §5 M3's measured wording.

#### C2 — **REFUTED: a 200 ms shrunk `probeTimeout` is a flake on this rig.**

In the same program the **known-positive control** (an instantly-answering fake, which MUST
succeed) reported `elapsed=202ms  err=DEADLINE EXCEEDED` on its first run — a false red. Isolated
re-measurement, 4 consecutive probes of the same script:

```
instant#0 elapsed=108ms  stdout="AILANG v0.30.0" deadline=false err=<nil>   <- first-ever exec of a freshly written script
instant#1 elapsed=9ms    ...
instant#2 elapsed=8ms    ...
instant@10s elapsed=7ms  ...
```

The **cold first exec of a freshly written script costs ~108 ms** here (~30 ms warm, 7–9 ms hot),
and under the load of the just-finished 10 s probe it crossed 200 ms. The doc's "e.g. 200 ms"
leaves **1.9×** headroom over a measured 108 ms cold start. **Adopted: the shrunk bound is
`1 * time.Second`** (9× over the measured cold exec), the blocking fake sleeps **20 s** (20× the
bound), and the wall-clock assertion is **< 8 s** — which an unbounded probe (20 s) cannot
satisfy and a bounded one (~1 s) clears by 8×. The control failure is recorded rather than
discarded: it did its job, and a zeroed/false control is the documented "instrument broken"
signal this mission runs on.

#### C3 — **SCOPED: the doc's production claim "never an unbounded hang" is true only for a single-process interpreter.**

Decision 2 states *"the new startup execution costs at most `probeTimeout`, never an unbounded
hang"*. Per C1 that is false for any child that forks a grandchild inheriting stdout (measured:
10.15 s at a 200 ms bound). It **is** true for the artifact actually in scope — the archived
`ailang` is a single compiled Go binary with no forked children, so `CommandContext`'s kill
closes the pipes. **The sprint must state the bound with that scope attached** and must not write
the unscoped sentence into a comment. The `cmd.WaitDelay` hardening that closes the general case
is measured (303 ms) and **DECLARED AS OWED, NOT ABSORBED**: it is +1 production line, it is not
in the design doc, and it **leaks the grandchild** (measured: `pgrep -fl '^sleep 10'` → pid 80813
alive after the program exited; the planner reaped it — iteration-97's rule that a load
experiment proves its own teardown). Route it as a follow-up row, not as sprint scope.

#### C4 — **GAP, load-bearing: the CONDITIONAL self-heal has no acceptance criterion and no mutation.**

Round 2's objection R2-a was **UPHELD** and rewrote Decision 2 from an unconditional heal to
*"probe ONLY when `!strings.HasPrefix(m.Version, \"AILANG v\")` … a healthy or already-healed
artifact performs **zero** process executions"*. That is the doc's own headline round-2 change.
**No AC and no mutation row in the doc asserts it.** Walked the list: AC1 (T1 exists/passes),
AC2/AC3 (the two `CombinedOutput` mutations), AC4a/AC4b (the heal *fires* on a polluted
manifest), AC5 (gates), AC6 (audit re-enumeration), AC7 (the bound). **An UNCONDITIONAL heal
passes every one of them.** Under this mission's rule 3e that is a vacuous-by-omission
criterion, and under the per-branch mutation protocol the new `if !HasPrefix(...)` is a refusal
branch with no neutering mutation. **Added: AC8, T2 arm 5 (an execution-counting fake), and
mutations M4a/M4b** (§5) — note `if false && cond` neuters the *heal*, and its dual
`if true || cond` is required to neuter the *skip*, because a skip cannot be neutered by
falsifying its guard.

#### C5 — **MEASURED: no existing fixture can demonstrate the "zero executions" property, and every one of them is classified POLLUTED by the new predicate.**

`grep -rn 'fakeInterpreter(' host/archive/ host/daemon/` → the version strings in play are
`"FakeAILANG v9.9.9\n"` (`archive_test.go:46`, `daemon_test.go:550`), `"v1\n"`
(`archive_test.go:114`, `:169`) and `"AILANG v0.30.0\nCommit: e37b370\n"` (`daemon_test.go:684`).
`strings.HasPrefix(v, "AILANG v")` is **FALSE for the first four**. Consequences the executor
must handle, not discover:

1. `archive_test.go:117/122` archives the SAME bytes twice with version `"v1\n"` — i.e. it takes
   the idempotent path. After the change that path will **spawn a probe where it spawned none**.
   The fake is deterministic, so fresh stdout == stored version → compare equal → no rewrite →
   the test stays green. **Verify it, do not assume it**: this is the single most likely place
   for an existing green to move.
2. AC8's "healthy artifact ⇒ zero executions" arm therefore needs a **new** fixture whose version
   starts with `AILANG v`. It cannot reuse any existing one.
3. `strings` is **not imported** in `host/archive/archive.go` today —
   `grep -n 'strings\.\|"strings"' host/archive/archive.go` → **0 hits, rc=1**, with the
   known-positive control `grep -c 'strings\.' host/daemon/daemon.go` → **2** in the same breath.
   The heal adds that import.

#### C6 — **CONSTRAINT: `New`'s signature must not change.** `grep -rn --include='*.go' 'archive\.New(' host/ cmd/` → **9 call sites** across `host/daemon` (1: `daemon.go:440`), `host/broker` (**4**: `episode_test.go:61`, `registry_publish_test.go:1245/1253/1275`), `host/capsule` (1), `host/replay` (3) — the per-package breakdown sums to the headline 9 only with broker at 4; an earlier draft read 3 and did not sum, caught by the iteration-98 evaluator and reproduced per-package by the controller. `Archive` today is `struct{ root string }` and `New(storeDBPath string) *Archive`. The doc's "field set from the constant by `New`" is therefore the *only* affordable wiring: add `probeTimeout time.Duration` to the struct, set it to `probeTimeout` in `New`, and leave all 9 callers untouched. Adding a parameter or an options struct would turn a 55-line change into a 5-package one.

#### C7 — **T3 fixture traps, named so the executor does not pay for them.**

- The generated fake `ailang` script **must live OUTSIDE `packageDir`**. `CrossCheck` computes
  `ContentHash(packageDir)` and `CreateTarball(packageDir)` over that directory; a fake dropped
  inside it changes the very hashes the fake is supposed to echo back.
- It must be given as an **absolute path**, because `cmd.Dir = packageDir` (`pkgproj.go:214`).
- The fake must write its stderr line and its stdout block **sequentially from one shell**. Under
  `CombinedOutput()` both fds are the *same* pipe, so sequential shell writes cannot interleave
  mid-line; if they could, a split line would fail the `(?m)^…$` anchors, `parseDryRun` would see
  ONE `Tarball` line, and **the M2 red would silently not fire**.
- The stderr line must be `  Tarball: 999 bytes (sha256:` + **exactly 17 lowercase hex** + `...)`
  — verified against `dryRunLine` at `pkgproj.go:181`.

#### C8 — **`go build ./...` is a real mutant-build check here, but it is not sufficient.**

`go build ./...` does **not** compile `_test.go`. All four mutations in §5 land in **non-test**
files (`archive.go`, `pkgproj.go`), so the doc's proof obligation is non-vacuous as written —
but a mutation that happens to break *test* compilation would still report "builds". The
mutant-build proof is therefore **three commands**, not one:
`go build ./...` **and** `go vet ./host/archive/ ./host/pkgproj/` **and**
`go test -run '^$' -count=1 ./host/archive/ ./host/pkgproj/` (compiles the test binaries without
running a single test).

---

## 1. Scope — what lands, what deliberately does not

**Lands (Milestone A, one commit) — ~377 lines total, of which ~290 are test:**

| File | Change | ~lines |
|---|---|---|
| `host/archive/archive.go` | `probeTimeout` constant; `probeTimeout` field on `Archive` set by `New`; `probeVersion` → method on `*Archive`, `exec.CommandContext` + separate `bytes.Buffer`s, `DeadlineExceeded` labelled and wrapped; CONDITIONAL self-heal on the idempotent path; `strings` import; doc comments (drop "combined output", add the *scoped* bound sentence per C3) | ~75 (50 probe + 25 heal) |
| `host/archive/archive_test.go` | stderr-emitting `fakeInterpreter` variant; T1; T2 arms 1–5 (heal, convergence, fail-loud, deadline, **zero-exec**) | ~200 |
| `host/pkgproj/pkgproj.go` | `CrossCheck` exec → `.Output()`; error path interpolates `out` + `ee.Stderr`; doc comment records the caller-supplied-bound obligation (P16) | ~12 |
| `host/pkgproj/pkgproj_test.go` | T3 (generated fake CLI, stdout dry-run block + one regex-matching stderr line) | ~90 |
| `host/daemon/daemon.go` | **comment only** — `HealthResponse.InterpreterVersion` doc (`:380-381`) says stdout | ~2 |

**Does NOT land — Design Freeze, restated so nobody "improves" it:**

- `host/broker/handlers.go:93` merges stderr into stdout **deliberately**
  (`cmd.Stderr = cmd.Stdout`, `:110-113`) so `classifyPublisherResult` can see markers on either
  stream. **Named non-change.** A commit that "fixes" it is a regression.
- `host/capsule/capsule.go:154` and `host/replay/replay.go:327` are already correct; `replay` is
  the model Decision 1 copies. No edit.
- The archived executable's **bytes, mode (0o555) and content address**. Only the sidecar moves.
- `tools/launchd/*` (frozen core, FLEET-owned). No `.ail` file. No new file, no new package.
- `cmd.WaitDelay` (C3) and reviewer option (2) — rollout safety for stores with existing epochs —
  are **DECLARED OWED**, routed as follow-up rows, not absorbed.

---

## 2. Day-by-day

**Budget: 1.0 day.** The doc's estimate is 0.5–1 d; this plan sits at the **top** of that band and
says so in those words rather than rounding down. Delta drivers: +0.15 d for the C1/C2 fixture
rework (measured, not guessed — the prototype is already written, which is why it is 0.15 and not
0.3), +0.15 d for AC8/T2-arm-5/M4a/M4b that the doc omits, +0.10 d for the five-mutation sweep
with landed-proof and restore controls, −0.10 d because the bounded-probe shape was prototyped
end-to-end this session and is reproduced verbatim in §3.

### Day 1 — morning (~3.5 h)

| Slot | Task | Output |
|---|---|---|
| 1 | **T1 — `host/archive/archive.go`** (§3): constant, field, `New` wiring, method conversion, `CommandContext` + buffers, labelled errors, `strings` import. **No heal yet.** | `go build ./...` rc=0; `go test ./host/archive/... -count=1` still `ok` (existing 7 `Archive(` call sites unchanged in behaviour) |
| 2 | **T2 — the conditional self-heal** on the idempotent path (§3.2), inserted immediately before `archive.go:266`'s `return ref, nil` | `go test ./host/archive/... ./host/daemon/... -count=1` `ok`; **explicitly re-check `archive_test.go:117/122`** (C5.1 — the one existing test whose idempotent path now spawns a probe) |
| 3 | **T3 — `host/pkgproj/pkgproj.go`** `.Output()` conversion + error-path `ee.Stderr` + the caller-bound doc comment | `go test ./host/pkgproj/... -count=1` `ok` (P13: the seam has zero coverage at base, so this leg cannot self-check — T5 is what proves it) |

### Day 1 — afternoon (~3.5 h)

| Slot | Task | Output |
|---|---|---|
| 4 | **T4 — `host/archive/archive_test.go`**: stderr `fakeInterpreter` variant; **T1-test** (AC1); **T2 arms 1–5** (AC4a, convergence, fail-loud, AC7 deadline with the C1/C2 shape, AC8 zero-exec) | `go test ./host/archive/ -run 'Stderr\|Heal\|Deadline\|NoProbe' -count=1 -v` — every named test **RUNS** (not "no tests to run") and passes |
| 5 | **T5 — `host/pkgproj/pkgproj_test.go`**: T3 with the C7 fixture rules | `go test ./host/pkgproj/ -run 'CrossCheckStderr' -count=1 -v` runs and passes |
| 6 | **T6 — the mutation sweep** (§5): M1, M2, M3, M4a, M4b, each with the 6-step protocol and its restore control; then AC6's re-enumeration; then the full gates | evidence block for the Verification Log; `verify_ail.sh` rc=0; `verify_go.sh` rc=0 |

### Day 1 evening / Day 2 spillover (only if T5's fixture fights back)

The doc names the honest ceiling: if `CreateTarball`'s hashing turns out to be order- or
path-sensitive in a fresh temp dir, T5 is the slot that overruns. **Spillover rule**: T5 slips,
T6 does not — the mutation sweep is the non-vacuity proof and is never the thing that gets cut.
If T5 cannot be made deterministic in +2 h, STOP and report; do not weaken the assertion to
"CrossCheck returns some error".

**One commit at the end of T6.** Message names the queue row, the doc, and the five mutations
with their measured reds.

---

## 3. The frozen implementation shapes

### 3.1 The bounded probe (`host/archive/archive.go`)

```go
// probeTimeout bounds the ELAPSED TIME of one "<interpreter> --version"
// execution. The probe now also runs on the idempotent path (the conditional
// self-heal below), which is reached on every daemon start, and nothing else
// bounds that sequence: daemon.New()'s startup block (daemon.go:431-462)
// carries no deadline, and host/archive had zero bounded execution before this
// constant. 10 s is ~2 orders of magnitude above the measured probe (rc=0 in
// well under 1 s for 168 bytes of stdout) and far below the verify gate's
// enclosing budgets (-timeout 8m race under a 600 s outer wait,
// verify_go.sh:150-153), so the labelled KindExecFailure always fires first.
//
// SCOPE OF THE BOUND (measured, not assumed): CommandContext SIGKILLs the
// direct child. For the archived ailang -- a single compiled binary with no
// forked children -- that closes the output pipes and Run() returns at the
// deadline. A hypothetical interpreter that forked a child inheriting stdout
// could hold Wait() open for that child's lifetime; closing that general case
// needs cmd.WaitDelay and is deliberately out of scope here.
const probeTimeout = 10 * time.Second

type Archive struct {
	// root is "<store>.db.artifacts": the parent of the "interpreters" tree.
	root string
	// probeTimeout is the per-probe wall-clock bound, seeded from the constant
	// by New. It is a field rather than a bare constant for exactly one reason:
	// tests shrink it (the daemon.go readDeadline idiom, daemon.go:284-290), and
	// that is the only way the deadline arm exercises the SHIPPED branch without
	// a ten-second test.
	probeTimeout time.Duration
}

func New(storeDBPath string) *Archive {
	return &Archive{root: storeDBPath + artifactsSuffix, probeTimeout: probeTimeout}
}

func (a *Archive) probeVersion(execPath string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.probeTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, execPath, "--version")
	// Decision 4 of w-self-mod-vertical: PRESERVED VERBATIM.
	cmd.Env = childenv.Scrubbed(os.Environ())
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("%s --version: timed out after %v: %w (stdout: %q, stderr: %q)",
				execPath, a.probeTimeout, context.DeadlineExceeded, stdout.String(), stderr.String())
		}
		return "", fmt.Errorf("%s --version: %w (stdout: %q, stderr: %q)",
			execPath, err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}
```

`New` keeps its one-parameter signature (C6). Callers of `probeVersion` become `a.probeVersion`
— today there is exactly one, `archive.go:289`.

### 3.2 The CONDITIONAL self-heal — inserted immediately before `archive.go:266`

```go
		// Identical bytes already archived: idempotent no-op success. The temp
		// file is discarded by the deferred cleanup.
		//
		// One exception, and it is CONDITIONAL on purpose: a sidecar written by
		// a build that merged the interpreter's stderr into its version string
		// is polluted forever, because this path never re-probes. Heal it -- but
		// only when the stored string fails the well-formedness test, so a
		// healthy or already-healed artifact performs ZERO process executions
		// and this really is the no-op the comment above claims.
		//
		// The predicate is POSITIVE-SHAPE, not a denylist for "Observatory":
		// the log line is one instance of a class (any stderr chatter), and
		// matching the known pollutant would leave the next one undetected.
		if !strings.HasPrefix(existing.Version, "AILANG v") {
			fresh, err := a.probeVersion(finalPath)
			if err != nil {
				return hashref.HashRef{}, &ReplayError{
					Kind:   KindExecFailure,
					Ref:    ref,
					Path:   finalPath,
					Detail: "cannot obtain --version from archived interpreter while healing its sidecar",
					Err:    err,
				}
			}
			if fresh != existing.Version {
				existing.Version = fresh
				if err := a.writeManifest(ref, existing); err != nil {
					return hashref.HashRef{}, err
				}
			}
		}
		return ref, nil
```

Note the two nested conditions are **not** redundant: the outer one is the zero-exec guarantee
(AC8/M4b), the inner one is convergence — no rewrite churn when the probe agrees (T2 arm 2).
Both get their own mutation.

### 3.3 `CrossCheck` (`host/pkgproj/pkgproj.go:219`)

`out, err := cmd.CombinedOutput()` → `out, err := cmd.Output()`; the `:221` error return becomes
`fmt.Errorf("ailang publish --dry-run: %w: stdout: %s", err, out)` plus `ee.Stderr` via an
`errors.As(err, &ee)` on `*exec.ExitError` (`ExitError.Stderr` is populated when `c.Stderr` is
nil — verified in the pinned toolchain's own `$GOROOT/src/os/exec/exec.go:1003`). The doc comment
above `CrossCheck` gains: *callers must supply the wall-clock bound; the sole production caller
does, via `run_bounded 120` in `scripts/verify_world_package.sh:183` (SIGKILL on the process
group, exit 124).*

---

## 4. Acceptance criteria — per-AC base status, with the command that establishes it

All commands: `export PATH=/opt/homebrew/bin:$PATH`, `GOTOOLCHAIN=go1.25.6`,
`AILANG_BIN=/tmp/ailang-v0300/ailang`, run from the repo root. Every base status below was
**measured this session at `7d79ad3`**.

| AC | Criterion | Base status | Command that establishes the base | Can it FAIL at base? |
|---|---|---|---|---|
| **AC1** | T1 exists and **runs** (not "no tests to run") and passes: `go test ./host/archive/ -run 'Stderr' -count=1 -v` | **the test does not exist** | check `grep -c 'FAKE-STDERR-MARKER' host/archive/archive_test.go` → **0, rc=1**; control **in the same file** `grep -c 'fakeInterpreter' host/archive/archive_test.go` → **5, rc=0** | **YES** — and note `-run` with no match is green by construction, so AC1's tooth is carried by AC2, not by its own exit code. The `-v` "=== RUN" line is the non-vacuity proof that the test actually ran. |
| **AC2** | Named RED mutation M1 (restore `CombinedOutput()` at the probe): mutant builds, mutation proven landed, **T1 AND T2 arm 1 both FAIL**, revert clean | **not runnable** (T1/T2 absent) — but the defect it guards is **MEASURED PRESENT**: `sed -n '391p' host/archive/archive.go` → `out, err := cmd.CombinedOutput()` | as left | **YES at sprint end** — this is the decisive red; the queue row's own criterion is exactly this |
| **AC3** | Named RED mutation M2 (restore `CombinedOutput()` in `CrossCheck`): T3 fails with `duplicate dry-run Tarball line`; same protocol | **not runnable** (T3 absent; the seam has ZERO coverage) | `grep -n 'CrossCheck\|func Test' host/pkgproj/pkgproj_test.go` → 3 test funcs, **0** `CrossCheck` refs; and `sed -n '219p' host/pkgproj/pkgproj.go` → `CombinedOutput()` | **YES at sprint end** |
| **AC4a** *(primary)* | T2 arm 1 seeded with the **REAL** polluted manifest bytes: heal rewrites `version` to the clean stdout | **FAILS** — the fixture string is polluted and the idempotent path never re-probes | `python3 -c "…json…['version']"` on `/private/tmp/world-demo.db.artifacts/interpreters/sha256/e9746fef…3fb5/manifest.json` → begins `2026/08/18 21:02:42 Observatory: 301MB (warn threshold: 200MB)\n`; and `sed -n '266p' host/archive/archive.go` → `return ref, nil` before Step 5 | **YES** — hermetic, reproduces on any rig |
| **AC4b** *(live, secondary)* | `ailang-worldd serve --db /tmp/world-demo.db --ailang-bin /tmp/ailang-v0300/ailang`; then `curl -s localhost:<port>/v1/health` shows `interpreter_version` beginning `AILANG v0.30.0`, and the on-disk manifest's `version` no longer contains `Observatory` | **FAILS** | same manifest read as AC4a; store creation is available (`grep -n 'Open opens' host/store/store.go` → `:218` "(or creates)") | **YES.** **BOUNDED CLAIM**: this witnesses the ONE orphaned tree measured at `/private/tmp/world-demo.db.artifacts`. It claims nothing about stores outside the four searched roots — the doc's declared residual. Do not write "rollout safe" anywhere. |
| **AC5** | `./scripts/verify_ail.sh` rc=0; `go build ./...` rc=0; `go test ./host/archive/... ./host/pkgproj/... ./host/daemon/... -count=1` ok; `./scripts/verify_go.sh` rc=0 | **GREEN at base** | `verify_ail.sh` → **rc=0**, `10 required identities verified, 39 named tests pass`; build **rc=0**; `ok host/archive 2.368s`, `ok host/pkgproj 0.534s`, `ok host/daemon 3.232s`; full `verify_go.sh` re-run this session (§7) | **NO — and this is declared, not hidden.** A gate green at base measures the repo, not the change. AC5 is a **regression guard, NOT a tooth**; the change-specific teeth are AC1–AC4 and AC7–AC8. Per the doc's own scoping the sprint neither launders nor absorbs any pre-existing full-gate red. |
| **AC6** | Post-change re-enumeration `grep -rn --include='*.go' -E 'exec\.Command(Context)?\(' host/ cmd/ \| grep -v '_test.go'` returns the **same five** sites, with 1 and 2 quoted as reading via `.Output()`/buffers | **sites 1 and 2 read via `CombinedOutput()`** | the grep above → 5 lines, plus `sed -n '391p' host/archive/archive.go` and `sed -n '219p' host/pkgproj/pkgproj.go` | **YES** — stated as a **positive re-enumeration**, never as "0 `CombinedOutput` remaining" (that negative form is what this mission ruled out) |
| **AC7** | T2 arm 4: shrunk `probeTimeout` (**1 s**, per C2), fake `exec /bin/sleep 20`, `Archive()` returns in **< 8 s** wall, error is `KindExecFailure` with `errors.Is(err, context.DeadlineExceeded)` true | **not runnable**, and the defect is **MEASURED PRESENT** | `grep -n -E 'context\|Timeout\|CommandContext' host/archive/archive.go host/capsule/capsule.go host/replay/replay.go host/broker/handlers.go` — **one call, one scope**: `archive.go` = 1 hit (`:56`, comment prose); control fires with `execTimeout = 60s` / `handlerExecTimeout = 30s` / `capsuleExecTimeout = 60s` in the three siblings | **YES at sprint end.** Fixture shape is **NOT** the doc's — see C1/C2; the doc's shape reds against a correct implementation (measured 10.15 s at a 200 ms bound, 3/3) |
| **AC8** *(NEW — §0.2 C4)* | T2 arm 5: a fake whose script **appends one byte to a counter file on every invocation** and whose version starts `AILANG v`. Archive once (counter == 1). Call `Archive()` again on the same bytes. Assert (a) counter is **still 1** — zero executions on the idempotent path — and (b) the sidecar `version` is byte-unchanged. | **not runnable** (no conditional heal exists; the idempotent path returns at `:266`) | `grep -n 'strings\.\|"strings"' host/archive/archive.go` → **0 hits, rc=1**, with the same-call-class control `grep -c 'strings\.' host/daemon/daemon.go` → **2**; and the read of `archive.go:249-266` showing the bare `return ref, nil` | **YES at sprint end.** Without AC8, an UNCONDITIONAL heal passes AC1–AC7 unchanged — which is precisely the vacuity rule 3e forbids. |

---

## 5. The mutation protocol — PER BRANCH, five mutations

**Protocol, applied to every row (6 steps; results go in the sprint's Verification Log):**

1. Back up the target file first: `cp host/archive/archive.go /tmp/w21_backup/archive.go` (and
   likewise `pkgproj.go`). **Restores are `cp` from `/tmp/w21_backup/`, NEVER
   `git checkout -- <file>`** — the file is uncommitted by construction during the sweep, so
   `git checkout --` deletes the work and the restore control then *reports* the disaster instead
   of preventing it.
2. Apply the mutation. **Neuter with `if false && <cond>` (or `if true || <cond>` where the
   branch under test is a SKIP), never by deleting the block** — a deleted block can stop
   compiling, and "the mutant doesn't build" would masquerade as "the guard fired".
3. **Prove the mutant builds** — three commands, per C8:
   `go build ./...` rc=0 **and** `go vet ./host/archive/ ./host/pkgproj/` rc=0 **and**
   `go test -run '^$' -count=1 ./host/archive/ ./host/pkgproj/` rc=0.
4. **Prove the mutation landed**: `git diff --name-only` lists **exactly** the one file, and the
   diff hunk is quoted verbatim in the log.
5. **Show the named test FAILING**, with the failure text quoted — not merely a non-zero rc.
6. **Restore and prove clean**: `cp /tmp/w21_backup/<file> <file>`; `git status --porcelain`
   empty for that file; **restore control** — re-run the same test and show it GREEN again. A
   restore that is not proven green is a mutation you are still carrying.

| # | Branch under test | Mutation (compiles by construction) | Expected red (measured shape) | Which AC |
|---|---|---|---|---|
| **M1** | `probeVersion` returns stdout ONLY | in `probeVersion`, replace the buffer read with `out, err := cmd.CombinedOutput()` and `return string(out), nil` | **T1**: both assertions fail — `FAKE-STDERR-MARKER` appears in `m.Version`, and the full-value equality breaks. **T2 arm 1**: the heal writes the merged string, equality with clean stdout fails | AC2 |
| **M2** | `CrossCheck` parses stdout ONLY | in `CrossCheck`, `cmd.Output()` → `cmd.CombinedOutput()` | **T3** fails: `parseDryRun` sees TWO `Tarball` lines → `duplicate dry-run Tarball line`; the test asserts success | AC3 |
| **M3** | the probe is BOUNDED | `exec.CommandContext(ctx, execPath, "--version")` → `exec.Command(execPath, "--version")` (leave `ctx`/`cancel` in place so the mutant still compiles; `_ = ctx` if vet complains) | **T2 arm 4** fails on BOTH arms: the probe waits the fake's full **20 s** (> the 8 s wall assertion) and then returns rc=0 with **EMPTY stdout**, so `errors.Is(err, context.DeadlineExceeded)` is false. Deterministic ~20 s **assertion**-red, not a hang. *(This is the corrected shape — the doc's "succeeds after the fake's finite sleep" assumed a fake that answers; per C1 such a fake defeats the bound entirely.)* | AC7 |
| **M4a** *(NEW)* | the heal FIRES on a polluted sidecar | `if false && !strings.HasPrefix(existing.Version, "AILANG v")` | **T2 arm 1 / AC4a** fails: the polluted fixture is returned unhealed | AC4a |
| **M4b** *(NEW)* | the heal is SKIPPED on a healthy sidecar | `if true \|\| !strings.HasPrefix(existing.Version, "AILANG v")` — the **dual** form; `if false &&` cannot neuter a skip | **T2 arm 5 / AC8** fails: the counter reads **2**, not 1 — a process was executed on the no-op path | AC8 |

**Not mutations, and labelled as such so no decorative arm is written**: T2 arm 2 (convergence —
no rewrite when the probe agrees) is covered by M4b's counter only indirectly; assert it directly
by comparing the sidecar's `mtime`/bytes across two `Archive()` calls. T2 arm 3 (fail-loud on an
unexecutable archived binary) is a **fixture** arm, asserted directly, with no mutation.

---

## 6. Split refusal

**NO SPLIT.** Milestone A is the whole item, one commit. Every candidate seam was checked and each
lands a red or a vacuous tree:

- *"archive fix first, pkgproj second"* — two independent files, but the sprint's own non-vacuity
  proof (§5) is a single sweep over one backup dir and one clean-tree assertion; splitting it
  lands the pkgproj `.Output()` change **with no committed red test** for one commit's duration.
  That is the guard-not-a-gate failure this repo has paid for repeatedly.
- *"code first, tests second"* — lands the behaviour change with zero committed teeth. Refused.
- *"tests first, code second"* — lands a RED tree. `go test ./...` is CI's gate; nothing lands red.
- *"bounded probe first, self-heal second"* — the heal is what puts the probe on the startup path;
  landing the heal without the bound is round 1's objection 1 shipped deliberately. Landing the
  bound alone is fine but is 15 lines and no observable behaviour, i.e. a commit with no tooth.
- *"daemon.go comment separately"* — a two-line comment. Not a milestone.

---

## 7. Base-gate evidence at `7d79ad3` (planner, first-party)

| Gate | Command | Result |
|---|---|---|
| `.ail` gate | `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` | **rc=0** — `✓ world package gate PASSED: 9/9 steps performed non-zero work`; `✓ verify gate PASSED: 10 required identities verified, 39 named tests pass`. Its own first output line is `2026/08/20 01:17:28 Observatory: 314MB (warn threshold: 200MB)` — the stderr chatter this item stops merging into data, live in the base run. |
| Go build | `GOTOOLCHAIN=go1.25.6 go build ./...` | **rc=0** |
| Touched packages | `go test ./host/archive/... ./host/pkgproj/... ./host/daemon/... -count=1` | `ok host/archive 2.368s`, `ok host/pkgproj 0.534s`, `ok host/daemon 3.232s` |
| Full Go gate | `GOTOOLCHAIN=go1.25.6 AILANG_BIN=… ./scripts/verify_go.sh` | see §7.1 — measured in this session on an **uncontaminated** rig (planner reaped its own prototype orphans before starting; `pgrep -fl '^sleep 10'` → none) |

**Instrument discipline carried from iteration 96/97, applied here:**
never `… | tail -N; echo rc=$?` (that prints `tail`'s rc over a `FAIL`); never `${PIPESTATUS[0]}`
(bash-only, EMPTY in zsh); always `cmd > /tmp/out 2>&1; echo "rc=$?"` and read the file.
`host/broker` is genuinely slow — ~85 s plain, ~184 s race, with two page-bound tests at ~220 s
each running 1,048,576 iterations. **A slow suite here is not a hang.**

### 7.1 Full-gate result — GREEN at `7d79ad3`, and an instrument failure caught in the act

```
✓ go gate PASSED: build clean, plain and race tests pass with pinned AILANG_BIN (AILANG v0.30.0 …)
verify_go rc=0
```

`grep -c 'FAIL' /tmp/i98_verify_go.txt` → **0**, with the known-positive control on the **same
file in the same call** `grep -c 'PASS\|ok  '` → **44**, so the zero is a measured zero.
Plain leg: `ok host/archive 9.008s`, `ok host/broker 79.613s`, `ok host/daemon 9.566s`,
`ok host/pkgproj 3.145s`. Race leg (`-timeout 8m`): `ok host/archive 10.592s`,
`ok host/broker 109.997s`, `ok host/daemon 12.900s`, `ok host/pkgproj 4.222s`. Also green in the
same run: the tracked-binary hygiene gate (0 blobs / 218 files), the routing + decision-ledger
gate (9 passed, 0 failed), the **driver drift gate** (6 tracked driver files, working tree matches
HEAD — no `D-WORLD-DRIVER-1` red), and the race-detector known-positive control (the deliberate
`WARNING: DATA RACE` in `design_docs/verification/w-race-gate-blindspot/racecontrol/main.go`,
which is the gate proving the detector is armed — **not** a failure).

**So the doc's P11b is now RESOLVED in the strong direction at HEAD**: `verify_go.sh` is GREEN at
base, not UNDETERMINED. The sprint has nothing pre-existing to launder or absorb, and AC5's
scoping clause is retained anyway because it costs nothing and survives either outcome.

**INSTRUMENT FAILURE, RECORDED (planner, this session — a fresh instance of the class this
mission keeps paying for).** The planner ran the gate as
`./scripts/verify_go.sh > f 2>&1; echo "verify_go rc=$?" >> f; grep -c '^FAIL' f >> f`. The
harness reported the job **"failed with exit code 1"** — because the exit code of a compound
command is the exit code of its **LAST** command, and `grep -c` exits **1** when it matches
nothing. The gate's own rc was 0 and its own verdict line was `✓ go gate PASSED`. This is the
`… | tail; echo rc=$?` artifact from iteration 96 arriving through a different door: **any
trailing command's rc will impersonate the gate's**. Generalised rule for the executor: put the
gate LAST in its own invocation, or capture `rc=$?` on the very next line and read nothing else
as the verdict. The mission's `${PIPESTATUS[0]}` habit does not help here either — zsh expands it
EMPTY.

---

## 8. Risks

| Risk | Likelihood | Mitigation |
|---|---|---|
| T5/T3's tarball fixture is path- or order-sensitive in a fresh temp dir | medium — the doc names it as the honest 1 d ceiling driver | C7's rules remove the two known traps (fake outside `packageDir`, absolute path). Hard stop at +2 h; report rather than weaken the assertion. |
| The shrunk `probeTimeout` reds spuriously under rig load | **measured real at 200 ms** (C2) | 1 s bound = 9× the measured 108 ms cold exec; 8 s wall assertion = 8× the bound. If it still flakes, raise the bound and the wall assertion together, keeping the 20 s sleep fixed — never by widening only the wall assertion, which would let an unbounded probe pass. |
| An existing test's green moves because the idempotent path now probes | **measured likely target**: `archive_test.go:117/122` (version `"v1\n"`, fails the predicate) | Slot 2 re-runs `./host/archive/... ./host/daemon/...` immediately after the heal lands, before any new test is written, so a moved green is attributed to the heal and not to T4. |
| The executor "also fixes" `host/broker/handlers.go:93` | low, but it is the doc's named non-change | §1 restates it; AC6's re-enumeration quotes site 3 unchanged. |
| The ambient `go1.26.4` leaks into one command | **measured real** — it bit iteration 97's controller for three suites | Every command in this plan carries `GOTOOLCHAIN=go1.25.6`; `host/store/toolchain_canary_test.go` is the backstop. |
| `rc=127` on any command | recurring | It is a **PATH gap**, not a broken toolchain and not a spent quota. `export PATH=/opt/homebrew/bin:$PATH` first. |

---

## 9. Handoff

- **Executor lane**: opus sprint-executor, worktree isolation.
- **Worktree**: a **sibling of the repo**, e.g. `/Users/voightkampff/dev/sunholo-data/ailang-world-w21`.
  **NEVER under `/tmp`** — `host/verifygate` and `host/boundary` derive `repoRoot` from
  `runtime.Caller` and copy live trees, so a relocated checkout reds for its location rather than
  for the code. (`/tmp` is additionally a symlink to `private/tmp` on this rig, which `find`
  silently declines as a traversal root.)
- **Commits**: ONE, at the end of T6, per the controller's granularity rule.
- **Do not touch**: `tools/launchd/*` (frozen core), `~/.ailang/state/mission-v1*`.
- **Companion JSON**: `.ailang/state/sprints/w-archive-stderr-in-manifest.plan.json`.

SPRINT_PLAN_PATH: design_docs/planned/w-archive-stderr-in-manifest-sprint-plan.md
SPRINT_JSON_PATH: .ailang/state/sprints/w-archive-stderr-in-manifest.plan.json
