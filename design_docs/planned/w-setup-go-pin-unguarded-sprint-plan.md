# Sprint plan — `w-setup-go-pin-unguarded` (queue row 41, iteration 128)

**Design doc**: [`w-setup-go-pin-unguarded.md`](w-setup-go-pin-unguarded.md) — 714 lines, quorum-cleared,
landed on `dev` at `74c47d5`.
**Planner**: `claude:opus-5` (sprint-planner), 2026-08-26.
**Base**: `74c47d5683ae5b86b19e1908da11ceccb3bb1c93`, tree clean (`git status --porcelain | wc -l` → `0`).
**Base currency check**: `git diff --name-only fd490ca 74c47d5` → **exactly one path**,
`design_docs/planned/w-setup-go-pin-unguarded.md`. Every `Vnn` row the doc measured at `fd490ca`
therefore describes an *unchanged* code tree at `74c47d5`; the planner re-ran the load-bearing ones
anyway (§Planner first-party verification) and they all reproduced.
**Sprint worktree**: `/Users/voightkampff/dev/sunholo-data/.wt-world-iter128`.
**Platform**: `darwin/arm64`, `go version go1.26.6 darwin/arm64` (auto-selected by `go.mod:3`).
**Estimate**: ~0.15 day. **Risk**: LOW for the production surface (one new test file + ~10 lines of
shell), MEDIUM for the drill (18 arms, one of which is network-bound — see D1).

---

## 0. The shape of this sprint in one paragraph

Two production artifacts change: one **new** Go test file (`host/verifygate/toolchain_pin_gate_test.go`,
two tests, ~170 LOC) and ~10 lines of **existing** shell
(`design_docs/verification/w-race-gate-blindspot/run.sh`). Nothing else. The whole difficulty of the
sprint is **ordering**: AC2 demands that Test B be *recorded RED* before the `run.sh` edit and *recorded
GREEN* after, with Test A green at both readings. A sprint that writes the test and the edit together
cannot discharge AC2 no matter how green its final tree is — the doc says so explicitly
(`w-setup-go-pin-unguarded.md:492`: "A sprint reporting only the final green has NOT discharged AC2").
That single constraint is what fixes the milestone split below at **two**, and it is the only reason
this ~0.15-day item is decomposed at all.

## 1. Milestone split — two stages, ONE commit

| # | Name | Files touched | Doc ACs closed | Ends with |
|---|---|---|---|---|
| **M1** | **RED stage** — author both tests; record Test B failing and Test A passing | `host/verifygate/toolchain_pin_gate_test.go` (**CREATE**) | AC2 **(i)**, AC5, and Test A's half of AC1/AC4 | Test B **FAILING** on purpose, `run.sh` byte-untouched |
| **M2** | **GREEN stage** — apply the `run.sh` edit; record both tests passing; run the drill | `design_docs/verification/w-race-gate-blindspot/run.sh` (**MODIFY**) | AC2 **(ii)**, AC1, AC3, AC4, AC5, AC6 | full green tree + 18 recorded mutation arms |

**Why two and not three.** The drill (AC3) is not a third milestone: it changes no production artifact
that survives, it runs on the M2 tree, and inventing a milestone for it would misrepresent a
verification pass as a delivery. Three milestones would be padding.

**Why two and not one.** AC2's evidence is *two recorded runs of the same command against two different
trees*. There is no way to obtain the first reading once the `run.sh` edit exists. The stages are
therefore mandatory, and their boundary is the `run.sh` edit.

### ONE commit, not two — and why that is not a contradiction

The executor performs **no git write operations at all** (house rule, §7). The controller builds a
**single** commit from the final M2 tree. This is deliberate and it is what keeps the charter's
"nothing lands red" rule and AC2's "record a red" requirement from colliding: **M1's red is a recorded
*reading*, never a committed *state*.** If M1 were committed on its own, that commit would carry a
failing Go test into CI. It must not be.

**Consequence for the executor**: the M1 → M2 boundary is a *reporting* boundary. Snapshot M1's two
command outputs verbatim into the implementation report **before** touching `run.sh`. There is no
second chance.

---

## 2. Milestone M1 — the RED stage

### Files
- **CREATE** `host/verifygate/toolchain_pin_gate_test.go`, `package verifygate`.
- **TOUCH NOTHING ELSE.** In particular `run.sh` must remain byte-identical; M1's exit condition asserts
  it (`shasum -a 256`).

### What to write
Exactly the two tests specified at `w-setup-go-pin-unguarded.md:147-192` (Test A, six steps) and
`:194-234` (Test B, seven steps), plus the **two verbatim DECLARED-RESIDUAL doc comments** at `:303-315`
and `:320-328`. The doc is prescriptive here; the plan does not restate the spec. Read those spans
before writing a line.

Six things the planner measured that the executor will otherwise re-derive the hard way:

1. **`repoRoot` already exists** as a package-level var at `host/verifygate/ail_binary_gate_test.go:27`
   (`repoRoot = findRepoRoot()`, `findRepoRoot` at `:31`). Reuse it. Do **not** redeclare it, and do
   **not** add a second `findRepoRoot`.
2. **Neither test may call `requirePinned`** (`ail_binary_gate_test.go:39`) or read `AILANG_BIN`.
   That is AC4, and the reason is measured: the package's shim-arm tests `t.Fatal` when `AILANG_BIN`
   is unset (doc V14), so a new test that borrowed `requirePinned` would red AC4 while AC1 stayed green
   on a rig that happens to export the variable.
3. **`filepath.Glob` returns directory-prefixed paths.** Test A step 6 must map every match through
   `filepath.Base` before comparing to `[]string{"ci.yml"}`. Without it the assertion fails
   unconditionally, at base, forever (doc V36). This is the single most likely way to lose an hour.
4. **The precedent to mirror** is `TestZ3PinDeclaredOnceAndInstalledInBothJobs`,
   `host/verifygate/ail_binary_gate_test.go:668`. Its known-positive-control block is `:675-681`; its
   table-driven count block is `:682-695`; its line-exact repair rationale is `:697-700`. Copy the
   *structure*, and duplicate the control needles rather than sharing a helper — the doc's Conflict
   Surface requires that one helper's edit cannot blind both tests.
5. **No name collisions.** The planner searched `host/verifygate/` for `pinValues`,
   `TestGoToolchainPinsAgreeAndMatchJobList`, `TestMiscompileInstrumentProbesPinnedToolchain` and
   `toolchain_pin_gate`: **0 files with hits for each**, against a same-scope known-positive control
   (`repoRoot` → **2 files**). The package has three files today:
   `ail_binary_gate_test.go`, `evidence_manifest_gate_test.go`, `module_manifest_gate_test.go`.
   The third is **row 43's fenced file — do not open it.**
6. **The base numbers Test A's constants must equal**, all re-measured by the planner at `74c47d5`:
   job set `{ailang-verify, go-verify}` at `ci.yml:17` / `:98`; `GOTOOLCHAIN` keyed lines **2**
   (`:21`, `:102`); `go-version` keyed lines **2** (`:28`, `:109`); `uses: actions/setup-go@v5` **2**;
   `go-version-file` **0**; `.github/workflows/` contains exactly **`ci.yml`**; `go.mod:3` is
   `go 1.26.6` with **0** `^toolchain ` lines.

### M1 exit gates (all EXECUTOR)

| Gate | Command | Base reading at `74c47d5` | Required at M1 end |
|---|---|---|---|
| M1-a | `gofmt -l host/verifygate/` | rc=0, **0 lines** | rc=0, 0 lines |
| M1-b | `go vet ./host/verifygate/` | rc=0, no output | rc=0, no output |
| M1-c | `go build ./...` | rc=0, **no output** | rc=0, no output |
| M1-d | `go test ./host/verifygate/ -run 'TestGoToolchainPinsAgreeAndMatchJobList' -count=1 -v` | rc=0 but `[no tests to run]` — the test does not exist | **rc=0, exactly 1 `=== RUN`, 1 `--- PASS`** |
| M1-e | `go test ./host/verifygate/ -run 'TestMiscompileInstrumentProbesPinnedToolchain' -count=1 -v` | rc=0 but `[no tests to run]` | **rc=1, exactly 1 `=== RUN`, 1 `--- FAIL`, failure text naming the first missing piece** |
| M1-f | `go test ./host/verifygate/ -run 'TestZ3PinDeclaredOnceAndInstalledInBothJobs' -count=1 -v` | rc=0, `--- PASS (0.00s)` | rc=0, `--- PASS` — the same-package neighbour must not move |
| M1-g | `shasum -a 256 design_docs/verification/w-race-gate-blindspot/run.sh` | `b8b19bc3…` is the *git blob*; take the file digest at sprint start and compare | **unchanged from sprint start** |
| M1-h | `git status --porcelain` | 0 lines | **exactly one line**, `?? host/verifygate/toolchain_pin_gate_test.go` |

**M1-e is the AC2(i) evidence.** Record its output *verbatim*. The expected failure names one or more
of: the floor token `go1.26.6` absent from `KNOWN_GOOD`; no `PINNED=` line; `saw_pinned_ok` site count
below 3. All three are the measured base state (doc V6, V34; planner-reproduced §8 rows P6/P7).

**Capture rule.** `go test` exit codes must be taken without a pipe:
`go test … > /tmp/out 2>&1; rc=$?`. `${PIPESTATUS[0]}` is silently empty in `zsh`; do not use it.

---

## 3. Milestone M2 — the GREEN stage

### Files
- **MODIFY** `design_docs/verification/w-race-gate-blindspot/run.sh` — and only that file.

### The edit
Verbatim from `w-setup-go-pin-unguarded.md:241-265`. Landing sites, re-measured by the planner at
`74c47d5` against a full numbered read of the 94-line script:

| Site | Base line | Change |
|---|---|---|
| `KNOWN_GOOD` | `:25` `KNOWN_GOOD="go1.25.6 go1.24.9"` | prepend `go1.26.6`; add the `PINNED="go1.26.6"` line with its two-line comment |
| flags | `:27-:29` `saw_bad=0` / `saw_good=0` / `ran=0` | add `saw_pinned_ok=0` |
| `probe()` case | `:50` `case "$out" in`, `:51` `OK*)`, `:52` `BUG*)`, `:53` `esac` | the `OK*)` arm at **`:51`** gains `; [ "$tc" = "$PINNED" ] && saw_pinned_ok=1` |
| fourth guard | after the `saw_good` block `:87-:91`, before the RESULT banner `:92-:94` | the 6-line `INSTRUMENT FAILURE: the PINNED toolchain` block |
| banner | after `:94` | one `pinned toolchain ($PINNED) reported OK  : yes` line |

Two shell facts the planner checked in the numbered read so the executor does not have to guess:

- `run.sh:20` is `set -uo pipefail` — **no `-e`**. The `OK*)` arm therefore cannot abort `probe()` when
  `[ "$tc" = "$PINNED" ]` is false, and `probe()`'s explicit `return 0` at `:54` normalises the arm's
  exit status anyway. The doc's one-liner is safe as written.
- The `SKIPPED` path is `:34-:38` and `return 0`s **without touching any flag**. That is precisely the
  hole the fourth guard closes, and it is what makes AC6's guard-trip rehearsal work.

### M2 exit gates

| Gate | Lane | Command | Base reading at `74c47d5` | Required at M2 end |
|---|---|---|---|---|
| G1 | executor | `go build ./...` | **rc=0, no output** (planner-measured) | rc=0, no output |
| G2 | executor | `gofmt -l host/verifygate/` | **rc=0, 0 lines** | rc=0, 0 lines |
| G3 | executor | `go vet ./host/verifygate/` | **rc=0, no output** | rc=0, no output |
| G4 | executor | `go test ./host/verifygate/ -run 'TestGoToolchainPinsAgreeAndMatchJobList\|TestMiscompileInstrumentProbesPinnedToolchain' -count=1 -v` | **rc=0 with `testing: warning: no tests to run` / `ok … [no tests to run]` — 0 `=== RUN` lines** | **rc=0, exactly 2 `=== RUN`, 2 `--- PASS`** (AC1) |
| G5 | executor | `go test ./host/verifygate/ -run 'TestNoSuchToolchainPinTestZZZ' -count=1 -v` | **rc=0, `testing: warning: no tests to run`, `ok … [no tests to run]`** | **UNCHANGED** — still `[no tests to run]`. This is AC1's vacuity control: it must keep saying so *after* the sprint |
| G6 | executor | `env -u AILANG_BIN go test ./host/verifygate/ -run '<A>\|<B>' -count=1` | **rc=0, `ok … [no tests to run]`** | **rc=0, `ok`, no `[no tests to run]`** (AC4) |
| G7 | executor | `go test ./host/verifygate/ -run 'TestZ3PinDeclaredOnceAndInstalledInBothJobs' -count=1 -v` | **rc=0, `--- PASS (0.00s)`, `ok … 0.186s`** | rc=0, `--- PASS` |
| G8 | executor | `bash -n design_docs/verification/w-race-gate-blindspot/run.sh` | **rc=0, no output** | rc=0, no output — the cheap syntax gate on the edited shell |
| G9 | executor | `git status --porcelain` | **0 lines** | **exactly two lines**: `?? host/verifygate/toolchain_pin_gate_test.go` and ` M design_docs/verification/w-race-gate-blindspot/run.sh` |
| G10 | **controller** | `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh` | **rc=1 — RED AT BASE, twice, both times only `cmd/ailang-worldd TestCLIRealSubprocessEpisode`; CI green on this exact commit. See §4.3** | rc=0, **or** rc=1 with that single package/test and the standalone control green — CONDITIONED gate |
| G10b | **controller** | `AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -count=1` | **rc=0, `ok … 53.723s`** — deterministic, no socket, the package this sprint changes | rc=0 |
| G10c | **controller** | `go test ./cmd/ailang-worldd/ -race -count=1 -run 'TestCLIRealSubprocessEpisode'` | **rc=0 on 3/3 consecutive runs** (`ok … 3.086s / 2.452s / 4.001s`) — G10's discharge control | rc=0 |
| G11 | **controller** | `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` | **rc=0**; `✓ world package gate PASSED: 9/9 steps`; `✓ verify gate PASSED: 11 required identities verified, 40 named tests pass` | rc=0, **byte-identical banner** — this sprint writes no `.ail`, so any movement in either number is a scope breach |
| G12 | **controller** | `./design_docs/verification/w-race-gate-blindspot/run.sh` | **UNINFORMATIVE UNDER SANDBOX** — see §5 | rc=0, pinned line in banner (AC6 i) |
| G13 | **controller** | AC6's guard-trip rehearsal (§6) | mechanism absent at base: `grep -c 'PINNED\|saw_pinned' run.sh` → **0** with same-file control `grep -c 'KNOWN_' run.sh` → **4** | rc=1, `INSTRUMENT FAILURE: the PINNED toolchain (go1.99.9) never reported OK`, **no RESULT banner** |

**G4 is AC2(ii).** Record it verbatim beside M1-e. The pair *is* the acceptance criterion.

---

## 4. The gate list — every base reading, measured

Every command below was **run by the planner on the unmodified `74c47d5` tree** in the sprint worktree,
shell `zsh`, exit codes captured without a pipe. Nothing here is transcribed from the design doc.

### 4.1 `AILANG_BIN` is mandatory, and the bare form is red by design

```
env -u AILANG_BIN ./scripts/verify_go.sh   →  rc=1
✗ AILANG_BIN is unset — host/replay tests would t.Skip() silently and this gate would be false-green.
  Export the pinned released binary, e.g. AILANG_BIN=/tmp/ailang-v0300/ailang
```

That red is the script refusing loudly, not a repo defect. **Every gate line that invokes either verify
script therefore carries `AILANG_BIN=/tmp/ailang-v0300/ailang`.** The planner confirmed the path exists
and is the pinned release:

```
/tmp/ailang-v0300/ailang --version  →  rc=0,  "AILANG v0.30.0",  "Commit: e37b370"
```

If that path is ever absent, the controller must re-materialise the pinned v0.30.0 release before
running G10/G11 — it must **not** substitute a dev build, and must **not** drop the prefix (that
converts a strict gate into an unsatisfiable one, not a lenient one).

### 4.2 `go build ./...` — MEASURED, not assumed

```
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter128
go build ./... > /tmp/out 2>&1; rc=$?
→ rc=0, output 0 lines
```

**Green at base.** Fast, no network beyond the already-populated module cache, no socket. Executor gate.

### 4.3 `verify_go.sh` — **RED AT BASE on this machine**, and the red is not this sprint's (D2, D5)

This is the gate the controller must understand before trusting any result from it. The planner ran
`AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh` **twice on the unmodified `74c47d5` tree**.
**Both runs exited rc=1. Both failed in the same package, on the same test, with two different
signatures.**

**Run 1 (cold cache) — rc=1, in the PLAIN `go test ./...` leg:**
```
── go test ./... -count=1
--- FAIL: TestCLIRealSubprocessEpisode (30.02s)
    cli_test.go:128: build subprocess binary: signal: killed
FAIL	github.com/sunholo-data/ailang-world/cmd/ailang-worldd	30.764s
…
ok  	github.com/sunholo-data/ailang-world/host/verifygate	259.610s
VERIFY_GO_RC=1
```

**Run 2 (warm cache) — the plain leg PASSED; rc=1 in the `-race -timeout 8m` leg:**
```
── go test ./... -count=1          (all packages ok)
── go test ./... -count=1 -race -timeout 8m
WARNING: DATA RACE
Read at 0x00c0001513b0 by goroutine 24:
  bytes.(*Buffer).String()
  …cmd/ailang-worldd.TestCLIRealSubprocessEpisode()  cli_test.go:176
Previous write at 0x00c0001513b0 by goroutine 27:
  bytes.(*Buffer).grow() … io.Copy … os/exec.(*Cmd).writerDescriptor.func1()
  …cmd/ailang-worldd.TestCLIRealSubprocessEpisode()  cli_test.go:139
--- FAIL: TestCLIRealSubprocessEpisode (11.30s)
    cli_test.go:176: daemon announcement timed out; stderr=
    testing.go:1712: race detected during execution of test
FAIL	github.com/sunholo-data/ailang-world/cmd/ailang-worldd	12.905s
…
ok  	github.com/sunholo-data/ailang-world/host/verifygate	138.043s
VERIFY_GO_RC=1
```

**Diagnosis, read from the source rather than guessed** (`cmd/ailang-worldd/cli_test.go`):

- `:123` wraps an **in-test `go build`** in a hard `context.WithTimeout(…, 30*time.Second)`. **That is
  run 1.** A cold `GOCACHE` blows the budget and `exec` reports `signal: killed`. The planner had, in
  the minutes before, probed mutations **M8** (`go.mod` floor → `go 1.27.0`, which *fetches and switches
  to a different compiler*) and **M9**, each forcing a full stdlib rebuild and evicting the cache. The
  corroborating symptom in the same transcript is `host/verifygate 259.610s`.
- `:137-138` is `var daemonErr bytes.Buffer; cmd.Stderr = &daemonErr`, so `os/exec` starts a copier
  goroutine that writes into `daemonErr` until `cmd.Wait()` returns — and `Wait()` is parked in the
  goroutine at `:143`. `:169-177` is a `select` with a **5-second** `time.After` at `:175`, and its
  timeout branch at `:176` calls `daemonErr.String()` **while that copier is still writing**. **That is
  run 2**: an unsynchronised read of an `exec.Cmd.Stderr` buffer, latent on the happy path and
  observable only when the announcement times out — which happens under the load of a full `-race`
  sweep (`host/broker` alone took **233.505 s** in that transcript).

**Three controls that fix the attribution:**

| Control | Result |
|---|---|
| the failing test alone, warm, no `-race` | `--- PASS: TestCLIRealSubprocessEpisode (2.90s)`, rc=0 — vs a killed 30.00 s budget cold |
| the failing test alone, warm, **`-race`**, 3 consecutive runs | rc=0, rc=0, rc=0 (`ok … 3.086s / 2.452s / 4.001s`) |
| **CI on this exact commit** — `gh run list --branch dev --limit 5` | `74c47d5` **success**; the four commits before it (`fd490ca`, `1cc8cf4`, `699f592`, `592a221`) all **success** |

**Verdict: `verify_go.sh` is red at base on this darwin/arm64 machine under full-sweep `-race` load, green
in CI, and green standalone 3/3. It measures the repo — specifically a pre-existing unsynchronised
buffer read in `cmd/ailang-worldd`'s own test — not this sprint.** Row 41 touches no Go production code
and adds one test file in `host/verifygate`, a package with **zero** socket binds and zero daemon
coupling.

**The gate is REPAIRED rather than handed over red, in two parts:**

- **G10 is CONDITIONED, not deleted.** `verify_go.sh` still runs, and its result is still a sprint
  verdict — *except* that a red whose **only** failing package is `cmd/ailang-worldd` with one of the two
  signatures above (`cli_test.go:128 build subprocess binary: signal: killed`, or `cli_test.go:176
  daemon announcement timed out` + `race detected`) is the base flake. The controller discharges it by
  running the standalone control `go test ./cmd/ailang-worldd/ -race -count=1 -run
  'TestCLIRealSubprocessEpisode'` (base: 3/3 rc=0) and **reporting both readings**. It is never silently
  re-run, and any *other* failing package is a genuine sprint red.
- **G10b is ADDED as the deterministic blast-radius gate**, because a conditioned gate alone is not a
  gate. `AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -count=1` — **base: rc=0,
  `ok github.com/sunholo-data/ailang-world/host/verifygate 53.723s`** (re-derived; a first run measured
  49.367 s). This is the package the sprint actually changes, it binds no socket, it is not flaky, and
  it exercises the whole package including the shim arms that need the pinned binary. **It is a
  CONTROLLER gate only because it reads `AILANG_BIN`** — an executor may run it and must label the
  result `UNINFORMATIVE UNDER SANDBOX` if the binary is unreachable.

**And the cache rule stands independently (D2):** any `verify_go.sh` run that follows a
toolchain-switching mutation arm or a `run.sh` execution **must be preceded by a cache-warming
`go build ./...`**, and the drill must place the toolchain-switching arms **last** (§6.3).

Also measured about `verify_go.sh`, because it decides the lane:

```
grep -n "go test" scripts/verify_go.sh
→ :245 go test -json ./host/evidence -count=1
  :256/:258 go test ./... -count=1
  :262 go test ./... -count=1 -race -timeout 8m
```
It runs `go test ./...` twice, and `./...` includes packages that **bind loopback sockets** — see §5.

### 4.4 `verify_ail.sh`

Controller gate, run for charter compliance rather than blast radius: **this sprint writes no `.ail`
source and touches no `.ail` file**, so the expected result is "byte-for-byte the base banner". The
planner measured `grep -c "ai-check" scripts/verify_ail.sh` → **13** (the script is the AILANG gate) and
found **0** hits for `listen|127.0.0.1|localhost|curl|wget` in it. Base reading: §8, row V-G11.

### 4.5 The narrow Go gates — all green at base, all fast

| Command | rc | Literal output |
|---|---|---|
| `gofmt -l host/verifygate/` | 0 | 0 lines |
| `go vet ./host/verifygate/` | 0 | (none) |
| `go test ./host/verifygate/ -run '<A>\|<B>' -count=1 -v` | 0 | `testing: warning: no tests to run` / `PASS` / `ok … 0.331s [no tests to run]` |
| `go test ./host/verifygate/ -run 'TestNoSuchToolchainPinTestZZZ' -count=1 -v` | 0 | `testing: warning: no tests to run` / `ok … 0.175s [no tests to run]` |
| `env -u AILANG_BIN go test ./host/verifygate/ -run '<A>\|<B>' -count=1` | 0 | `ok … 0.186s [no tests to run]` |
| `go test ./host/verifygate/ -run 'TestZ3Pin…' -count=1 -v` | 0 | `=== RUN` / `--- PASS (0.00s)` / `ok … 0.186s` |
| `bash -n design_docs/verification/w-race-gate-blindspot/run.sh` | 0 | (none) |
| `git status --porcelain` | 0 | 0 lines |

**The first three rows of that table are the AC1 trap made concrete.** `go test -run <name>` on a test
that does not exist exits **0** and prints `ok`. A gate phrased as "the command greens" is green at base
and measures nothing. Every gate above that names a test therefore also names its **`=== RUN` count**.

### 4.6 A gate deliberately NOT on the list

`env -u AILANG_BIN go test ./host/verifygate/ -count=1` (the whole package, no `-run`) is **red at base**
— the package's shim-arm tests `t.Fatal` on an unset `AILANG_BIN` by design (doc V14). It measures the
package's environment contract, not this sprint. It is excluded rather than handed over red. The
package-wide form that *is* meaningful, `AILANG_BIN=/tmp/ailang-v0300/ailang go test
./host/verifygate/ -count=1`, is already covered inside G10 (`verify_go.sh` runs `go test ./...`), so
listing it separately would double-count rather than add coverage.

---

## 5. Sandbox lane analysis — asked per gate, not assumed

The executor runs under `codex exec --sandbox workspace-write`, which **denies loopback socket binds**,
denies network, and denies writes outside the workspace. Each gate was interrogated on those three axes.

**Measured, with a same-scope known-positive control:**

```
grep -rn "net\.Listen\|httptest\.NewServer\|ListenAndServe" --include='*.go' host/ cmd/
→ host/broker/registry_reconcile_test.go:128   httptest.NewServer
  host/broker/registry_publish_test.go:110     httptest.NewServer
  host/daemon/daemon.go:634                    net.Listen("tcp", addr)
  host/daemon/handlers_test.go:420             httptest.NewServer
  host/daemon/daemon_test.go:561               httptest.NewServer
  cmd/ailang-worldd/cli_test.go:238            net.Listen("tcp", "127.0.0.1:0")
KP control, same scope: grep -rn "^func Test" --include='*_test.go' host/ | wc -l → 403
```

Zero of those hits are in `host/verifygate/`. That is the whole lane split:

- **`./scripts/verify_go.sh` runs `go test ./...`, which reaches `host/broker`, `host/daemon` and
  `cmd/ailang-worldd` — all of which bind loopback sockets. G10 is therefore a CONTROLLER gate.** An
  executor attempt is permitted but its result is `UNINFORMATIVE UNDER SANDBOX`, never pass and never
  fail. (`cmd/ailang-worldd/cli_test.go:132` also binds `127.0.0.1:0` in the subprocess daemon.)
- **`go test ./host/verifygate/…` binds nothing and needs no network.** All of G1–G9 are EXECUTOR gates.
- **`run.sh` FETCHES TOOLCHAINS FROM THE NETWORK** (`run.sh:17-18` says so in its own header comment;
  `probe()` at `:34` runs `GOTOOLCHAIN="$tc" go build`). Under the sandbox every probe reports `SKIPPED`,
  `ran` stays 0, and the script exits 1 at `:77-:80`. **Any gate that executes `run.sh` is
  UNINFORMATIVE UNDER SANDBOX. G12 and G13 are CONTROLLER gates, run outside the sandbox, and are
  never an executor obligation.** The executor's obligation regarding `run.sh` is *text only*: apply the
  edit, and pass `bash -n` (G8).
- **`verify_ail.sh` (G11)**: no socket, no network hits found; it is a controller gate for charter
  compliance and commit hygiene, not for sandbox reasons.

**Blanket rule for the executor's report**: a gate the executor could not obtain is reported as
`UNINFORMATIVE UNDER SANDBOX`, with the literal denial message. Never as a pass. Never as a fail.

---

## 6. The mutation drill (AC3) — 18 arms, split by lane

### 6.1 House recipe (applies to every arm without exception)

1. **Back up first.** `mkdir -p .snap/backup` and `cp` each mutable path into it *before the first
   mutant*, recording `shasum -a 256` for each: `.github/workflows/ci.yml`, `go.mod`,
   `design_docs/verification/w-race-gate-blindspot/run.sh`.
2. **Apply the exact edit** named in the doc's table (`w-setup-go-pin-unguarded.md:583-602`).
3. **Assert the mutant LANDED by querying the system's own view, not the file bytes you just wrote.**
   For `ci.yml` pin arms that is the occurrence count plus a line read-back (the doc's own V10/V11
   recipe: `1.26.6` occurrences **4 → 3**, and `sed -n '28p'` reads back `go-version: '1.25.6'`). For
   `go.mod` arms it is `grep -n '^go \|^toolchain' go.mod`. For the mode arm it is `git ls-files -s`
   plus `test -x`. A mutant you did not confirm landed is a claim, not a fact.
4. **Assert it BUILDS**: `go build ./... > /tmp/out 2>&1; rc=$?` → rc=0. A mutant that does not compile
   reds every test for the wrong reason and proves nothing about the new guard.
5. **Run the named test with `-run`** and capture without a pipe. Record the **literal** failure text.
6. **Restore from the `cp` backup.** `cp .snap/backup/<file> <path>`.
   **NEVER `git checkout -- <path>`** — the executor's own uncommitted work (the new test file, the
   `run.sh` edit) is in this tree and `git checkout --` would delete it.
7. **Assert sha256 byte-identity** against the pre-arm digest, and `git status --porcelain` back to its
   expected two-line M2 state.

**Additions and mode arms need a restore that `cp` cannot do:**
- **M5, M7** edit `ci.yml` in place → normal `cp` restore.
- **M14** *creates* `.github/workflows/other.yml` → restore is `rm`, then confirm
  `ls -1 .github/workflows/` prints exactly `ci.yml` and `git status --porcelain` has no `??` line for it.
- **M12** is `chmod -x` → **sha256 will be byte-identical even while the mutant is live**, so the digest
  proves nothing here. The landing predicate is `test -x … ; echo $?` → 1, and the restore predicate is
  `chmod +x` plus `git ls-files -s …/run.sh` → mode `100755` and `git status --porcelain` clean of it.
  Say so in the report; do not let the trivially-passing sha256 stand in as the restore proof.

### 6.2 M1 and M2 are the point of the whole item

**M1** (`ci.yml:28` `'1.26.6'` → `'1.25.6'`) and **M2** (`ci.yml:109`, same) are the `P6.T` drill's
**recorded SURVIVORS** (`w-mcp-projection.md:726`), reproduced as survivors again at `fd490ca` by this
doc's designer (V10, V11). Nothing in the repo turns them red today. **After this sprint they must RED
`TestGoToolchainPinsAgreeAndMatchJobList`.** If either survives, the sprint has not delivered its
thesis, and that is a finding to escalate immediately rather than paper over — the doc's mutation
protocol is explicit that a SURVIVED mutant is a *result*, recorded with the commands that failed to
turn it red, never widened-scope away.

### 6.3 Lane assignment

| Arm | Target | Lane | Measured reason |
|---|---|---|---|
| M1, M2, M3, M4, M5, M6, M7 | `.github/workflows/ci.yml` text | **EXECUTOR** | text edit; detector is `go test ./host/verifygate/` — no socket, no network |
| M9 | `go.mod` += `toolchain go1.25.6` | **EXECUTOR** | planner-measured: `go build ./...` **rc=0**, `go test ./host/verifygate/ -run TestZ3Pin…` **rc=0**. The floor `go 1.26.6` still wins the selection, so **no toolchain is fetched** |
| M10 | base state (`KNOWN_GOOD` lacks the pin) | **EXECUTOR — no arm to run** | it *is* the M1-stage tree; discharged by AC2(i)'s recorded red (gate M1-e). Record it as "discharged by AC2(i), see M1-e", not as a separate mutation |
| M11, M12, M13, M15, M16, M17, M18 | `run.sh` text / mode | **EXECUTOR** | text or mode edit; detector is a static Go text scan. **None of these executes `run.sh`** |
| M14 | new `.github/workflows/other.yml` | **EXECUTOR** | file creation inside the workspace |
| **M8** | `go.mod` `go 1.26.6` → `go 1.27.0` | **CONTROLLER** | **planner-measured**: the mutant triggers `go: downloading go1.27.0 (darwin/arm64)` — a **network toolchain fetch**, then an **89 s** full `go build ./...` and a **49 s** `go test`, and it evicts `GOCACHE` (see D2 / §4.3). Under `workspace-write` the fetch is denied *and* the toolchain unpacks outside the workspace, so the executor cannot make this mutant build. It is UNINFORMATIVE UNDER SANDBOX |

**M8's controller preconditions and ordering — mandatory:**
1. Before running M8, confirm the toolchain is obtainable:
   `ls -d ~/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.27.0.darwin-arm64` → must exist (the planner
   warmed it during this planning session; it is present on this machine now). If it is absent and the
   network is unavailable, M8 is **UNINFORMATIVE**, recorded as such — not skipped silently, and not
   substituted with a different version without saying so.
2. **Run M8 LAST of all 18 arms.** It is the only arm that switches compilers.
3. **After restoring M8, run `go build ./... ` to re-warm `GOCACHE` before G10.** Skipping this is how
   §4.3's `TestCLIRealSubprocessEpisode` flake is manufactured.
4. M8 is the doc's only **double-red** arm: it must red **both**
   `TestGoToolchainPinsAgreeAndMatchJobList` (cross-file floor mismatch) **and**
   `TestMiscompileInstrumentProbesPinnedToolchain` (`KNOWN_GOOD` lacks `go1.27.0`, and `PINNED=` ≠ the
   new floor). Run both names and record both.

### 6.4 AC6's runtime arms — CONTROLLER only

Neither is a mutation in the AC3 table; both execute `run.sh` and are therefore outside the sandbox.

- **G12 / AC6(i)** — attended `run.sh` on the post-sprint tree. Expect rc=0, a line reading
  `go1.26.6   expect=GOOD  got: OK (rc=0)`, and the new banner line. This works on `darwin/arm64`
  because a KNOWN_BAD toolchain genuinely reports `BUG` there (doc V32). It would **not** work on
  `ubuntu-latest`, where the miscompilation does not reproduce and `run.sh` already exits 1 on 10 of
  the last 10 `dev` runs (doc V37) — that is **queue row 44's** defect, pre-existing at HEAD, and this
  sprint neither fixes nor worsens it. Do not "fix" it here.
- **G13 / AC6(ii)** — the guard-trip rehearsal. Temporarily set `PINNED="go1.99.9"` and place
  `go1.99.9` **first** in `KNOWN_GOOD`; run; expect the `SKIPPED` line, then
  `INSTRUMENT FAILURE: the PINNED toolchain (go1.99.9) never reported OK`, **rc=1, and no `RESULT:`
  banner at all**. Restore by `cp` from the backup and assert sha256 byte-identity. House recipe, same
  as any arm.

**Cost note the controller should price in:** `run.sh` probes 6 toolchains, and
`go1.26.4` is the one KNOWN_BAD entry **not** in the local module cache
(`~/go/pkg/mod/golang.org/` holds go1.24.9, go1.25.6, go1.26.0, go1.26.2, go1.26.3, go1.26.5, go1.26.6,
go1.27.0 — measured). Expect one network fetch, and expect `run.sh` to churn `GOCACHE`; re-warm with
`go build ./...` before G10.

---

## 7. Executor policy

- **NO git write operations.** No `add`, `commit`, `branch`, `checkout`, `stash`, `push`. Read-only git
  (`status`, `diff`, `show`, `rev-parse`, `ls-files`) is expected and required. **The controller builds
  the single commit.**
- **Restores are `cp` from `.snap/backup/`, never `git checkout --`** — the tree holds uncommitted work.
- **Exit codes without pipes**: `cmd > /tmp/out 2>&1; rc=$?`. `${PIPESTATUS[0]}` is empty in `zsh`.
- **Never `|| echo 0` inside `$(...)`** around `grep -c`: `grep -c` exits 1 on a legitimate zero and the
  fallback concatenates into `0\n0`.
- **Quote glob-shaped flag values**: `--include='*.go'`. `zsh` aborts the command otherwise.
- **A search that found nothing is a claim.** Pair every empty result with a known-positive control
  scoped to the same path, in the same call.
- **Every count carries its scope** in the sentence that reports it.
- **Report the literal reading.** Never a predicted red as an observed one. A SURVIVED mutant is a
  finding, recorded with the commands that failed to turn it red.

---

## 8. Planner first-party verification (run at `74c47d5`, sprint worktree, `zsh`)

| # | Claim | Command | Observed |
|---|---|---|---|
| V-A | Base is `74c47d5`, clean | `git rev-parse HEAD`; `git status --porcelain \| wc -l` | `74c47d5683ae5b86b19e1908da11ceccb3bb1c93`; `0` |
| V-B | The only change since the doc's measurement base is the doc | `git diff --name-only fd490ca 74c47d5` | exactly `design_docs/planned/w-setup-go-pin-unguarded.md` |
| V-C | Toolchain boundary | `go version`; `go env GOVERSION` | `go1.26.6 darwin/arm64`; `go1.26.6` |
| V-D | `go.mod` floor, no `toolchain` directive | `grep -n "^go \|^toolchain" go.mod`; `c=$(grep -c "^toolchain " go.mod)` | `3:go 1.26.6`; count `0`, rc `1` (legitimate zero) |
| V-E | Job enumeration is exactly two | `awk '/^jobs:/{p=1;next} p && /^  [a-z0-9-]+:$/{print NR": "$0}' ci.yml` | `17:   ailang-verify:`, `98:   go-verify:` |
| V-F | Pin sites | `grep -n "^jobs:\|GOTOOLCHAIN\|go-version\|uses: actions/setup-go" ci.yml` | `16:jobs:`, `21`+`102` `GOTOOLCHAIN: go1.26.6`, `26`+`107` `setup-go@v5`, `28`+`109` `go-version: '1.26.6'` |
| V-G | Test A's control counts | `grep -c` per needle on `ci.yml` | `ailang-verify:` **1**; `go-verify:` **1**; `uses: actions/setup-go@v5` **2**; `./scripts/verify_go.sh` **1**; `go-version:` **2**; `go-version-file` **0** (rc=1, legitimate zero) |
| V-H | `ci.yml` length; workflows dir | `wc -l < ci.yml`; `ls -1 .github/workflows/` | `196`; exactly `ci.yml` |
| V-I | `run.sh` mode and lists | `test -x`; `git ls-files -s`; `sed -n '24,25p'` | `EXEC-yes`; `100755 b8b19bc3749f38d002d52b72a69fc044c8224ca0`; `KNOWN_BAD="go1.26.0 go1.26.3 go1.26.4 go1.26.5"`, `KNOWN_GOOD="go1.25.6 go1.24.9"` |
| V-J | Shebang bytes, with control | `head -1 run.sh \| od -c`; `grep -c '^#!/usr/bin/env bash$'`; KP `grep -c '^KNOWN_GOOD='` | `# ! / u s r / b i n / e n v   b a s h \n`; anchored **1**; KP **1** |
| V-K | The guard mechanism is ABSENT at base, against a firing control | `grep -c 'PINNED\|saw_pinned' run.sh`; KP `grep -c 'KNOWN_' run.sh` | **0** (rc=1); KP **4** |
| V-L | The pin is absent from `run.sh`, against a firing control | `grep -c "1\.26\.6" run.sh`; KP `grep -c "1\.25\.6" run.sh` | **0** (rc=1); KP **1** |
| V-M | `run.sh` parses | `bash -n run.sh > /tmp/out 2>&1; rc=$?` | rc=0, no output |
| V-N | `run.sh` shape, numbered read | `awk '{print NR": "$0}' run.sh` | 94 lines; `:20 set -uo pipefail`; lists `:24`/`:25`; flags `:27-:29`; `probe()` `:31-:55` with SKIPPED `:34-:38` and `case` `:50-:53`; three fail-loud blocks `:77-:80`, `:81-:86`, `:87-:91`; RESULT banner `:92-:94`; **no `PINNED`/`saw_pinned` identifier** |
| V-O | `go build ./...` | `go build ./... > /tmp/out 2>&1; rc=$?` | **rc=0, 0 lines of output** |
| V-P | AC1's base trap | `go test ./host/verifygate/ -run '<A>\|<B>' -count=1 -v` | rc=0, `testing: warning: no tests to run`, `ok … 0.331s [no tests to run]` — **0 `=== RUN` lines** |
| V-Q | AC1's nonsense control | `-run 'TestNoSuchToolchainPinTestZZZ' -count=1 -v` | rc=0, `ok … 0.175s [no tests to run]` |
| V-R | AC4's base form | `env -u AILANG_BIN go test ./host/verifygate/ -run '<A>\|<B>' -count=1` | rc=0, `ok … 0.186s [no tests to run]` |
| V-S | AC5 hygiene | `go vet ./host/verifygate/`; `gofmt -l host/verifygate/` | rc=0, no output; rc=0, **0 lines** |
| V-T | The precedent is green | `GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run 'TestZ3Pin…' -count=1 -v` | rc=0, `--- PASS (0.00s)`, `ok … 0.186s` |
| V-U | Bare `verify_go.sh` refuses loudly | `env -u AILANG_BIN ./scripts/verify_go.sh > /tmp/out 2>&1; rc=$?` | **rc=1**, `✗ AILANG_BIN is unset — host/replay tests would t.Skip() silently and this gate would be false-green.` |
| V-V | The pinned binary exists and is the release | `/tmp/ailang-v0300/ailang --version` | rc=0, `AILANG v0.30.0`, `Commit: e37b370` |
| V-W | No name collisions for the new identifiers, with control | `grep -rc` per name over `host/verifygate/`; KP `grep -rl "repoRoot" host/verifygate/ \| wc -l` | `pinValues` **0 files**, both test names **0 files**, `toolchain_pin_gate` **0 files**; KP **2 files** |
| V-X | `repoRoot`/`requirePinned` seats | `grep -n "^func findRepoRoot\|func requirePinned" host/verifygate/*.go`; read `:26-:29` | `ail_binary_gate_test.go:31` / `:39`; `var ( repoRoot = findRepoRoot(); pinned = os.Getenv("AILANG_BIN") )` at `:26-:29` |
| V-Y | **M8 probe** — the doc's `go 1.27.0` mutant is network-bound and cache-destroying | set `go.mod:3` → `go 1.27.0`; `go build ./...`; `go test ./host/verifygate/ -run TestZ3Pin…`; `cp` restore | **cold**: `go: downloading go1.27.0 (darwin/arm64)`, then a 2-minute wall-clock timeout on the test leg. **warm**: build rc=0 in **89 s**, test rc=0 in **49 s**. Restore sha256-equal (`7a2983…`), porcelain `0` |
| V-Z | **M9 probe** — the `toolchain` mutant is cheap and offline | append `toolchain go1.25.6` to `go.mod`; `go build ./...`; `go test …`; restore | landed (`grep -c '^toolchain '` → **1**); build rc=0 (48 s), test rc=0 (22 s); **no toolchain download**; restore sha256-equal, porcelain `0` |
| V-AA | `verify_go.sh` reaches socket-binding packages | `grep -n "go test" scripts/verify_go.sh`; `grep -rn "net\.Listen\|httptest\.NewServer\|ListenAndServe" --include='*.go' host/ cmd/`; KP `grep -rn "^func Test" --include='*_test.go' host/ \| wc -l` | `go test ./...` at `:258` and `:262 -race`; 6 socket sites in `host/broker`, `host/daemon`, `cmd/ailang-worldd`; **none in `host/verifygate/`**; KP **403** test functions |
| V-AB | The `verify_go.sh` base red is a cold-cache flake, not a repo defect | run 1: full `verify_go.sh`; control: the single failing test alone, warm | run 1 **rc=1**, `--- FAIL: TestCLIRealSubprocessEpisode (30.02s)  cli_test.go:128: build subprocess binary: signal: killed`, `host/verifygate 259.610s`. Control **rc=0**, `--- PASS: TestCLIRealSubprocessEpisode (2.90s)`. Budget is a hard `context.WithTimeout(…, 30*time.Second)` at `cli_test.go:123` |
| V-AC | Toolchain module cache inventory (prices `run.sh` and M8) | `ls -1d ~/go/pkg/mod/golang.org/toolchain@*/` | go1.24.9, go1.25.6, go1.26.0, go1.26.2, go1.26.3, go1.26.5, go1.26.6, **go1.27.0** (warmed by V-Y). **`go1.26.4` is NOT cached** — `run.sh` will fetch it |
| V-G10 | **`verify_go.sh` is RED AT BASE, twice, in one package only** | `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh > /tmp/out 2>&1; rc=$?`, run twice | **run 1 rc=1** (cold): plain leg, `--- FAIL: TestCLIRealSubprocessEpisode (30.02s) cli_test.go:128: build subprocess binary: signal: killed`. **run 2 rc=1** (warm): plain leg all `ok`, `-race` leg `WARNING: DATA RACE` read at `cli_test.go:176` vs write from `cli_test.go:139`, `--- FAIL … daemon announcement timed out; stderr=`, `race detected during execution of test`. **Every other package `ok` in both runs** |
| V-G10b | The deterministic blast-radius gate is green | `AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -count=1 > /tmp/out 2>&1; rc=$?` | **rc=0**, `ok github.com/sunholo-data/ailang-world/host/verifygate 53.723s` (a first run: 49.367 s) |
| V-G10c | The base red does not reproduce standalone under `-race` | `go test ./cmd/ailang-worldd/ -race -count=1 -run 'TestCLIRealSubprocessEpisode'`, three consecutive runs | **rc=0, rc=0, rc=0** — `ok … 3.086s`, `ok … 2.452s`, `ok … 4.001s` |
| V-G10d | CI is GREEN on this exact commit — the red is machine-local | `gh run list --branch dev --limit 5 --json displayTitle,conclusion,headSha,status` | `74c47d5` **success**; `fd490ca`, `1cc8cf4`, `699f592`, `592a221` all **success** |
| V-G10e | The race's mechanism, read from source not guessed | `awk 'NR>=136 && NR<=180' cmd/ailang-worldd/cli_test.go` | `:137-138` `var daemonErr bytes.Buffer; cmd.Stderr = &daemonErr`; `:143` `go func(){ waited <- cmd.Wait() }()`; `:175` `case <-time.After(5 * time.Second):`; `:176` `t.Fatalf("daemon announcement timed out; stderr=%s", daemonErr.String())` — an unsynchronised read of a live `exec.Cmd.Stderr` buffer, on the timeout path only |
| V-G11 | `verify_ail.sh` is GREEN at base | `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh > /tmp/out 2>&1; rc=$?` | **rc=0**; `✓ world package gate PASSED: 9/9 steps performed non-zero work`; `✓ verify gate PASSED: 11 required identities verified, 40 named tests pass` |
| V-G12 | The stderr-buffer race is not already filed anywhere in the mission's docs, with control | `grep -rc "daemonErr" --include='*.md' design_docs/`; KP `grep -c "w-floor-raise-coupling-inventory" design_docs/world-mission.md` | **0** design_docs files name `daemonErr`; KP **1** — the instrument sees strings that are there |
| V-G13 | The doc's proposed queue row 44 is NOT yet in the charter, with control | `grep -c "w-miscompile-instrument-inert-in-ci" design_docs/world-mission.md`; KP `grep -c "w-floor-raise-coupling-inventory" …` | **0** (rc=1); KP **1**. The charter's queue currently ends at row **43** |
| V-AD | The doc names no milestones (the AC cross-check's basis) | `grep -ci "milestone" <doc>`; KP `grep -c "Acceptance Criteria\|AC[1-6]" <doc>` | **0**; KP **26** |

---

## 9. AC cross-check — doc versus plan (controller-mandated)

**Method.** For each milestone in this plan, list the ACs the plan says it closes, and list the ACs the
DOC assigns to the corresponding stage. Diff. **The DOC wins where they disagree.**

**The doc's list is empty by construction.** `grep -ci "milestone"` over
`design_docs/planned/w-setup-go-pin-unguarded.md` → **0**, against a same-file known-positive control
(`grep -c "Acceptance Criteria\|AC[1-6]"` → **26**). The doc contains no milestone decomposition at all.
What it *does* impose is two ordering constraints, and those are the only things a diff can be taken
against:

- **AC2 (`:483-492`)** names two stages by name: "(i) test file applied, `run.sh` untouched → Test B
  FAILS …, Test A PASSES"; "(ii) the `run.sh` edit applied → both PASS".
- **AC3 (`:495`)** and **AC6 (`:514-527`)** both say "**on the post-sprint tree**", which places the
  drill and the runtime guard arms after stage (ii).
- **AC1, AC4, AC5** name no stage; they are final-tree properties.

**The diff:**

| Stage | ACs the DOC places here | ACs this PLAN's milestone section names | Delta |
|---|---|---|---|
| stage (i) / **M1** | AC2(i) | AC2(i), AC5, "Test A's half of AC1/AC4" | **plan is a superset** — see note 1 |
| stage (ii) / **M2** | AC2(ii), AC3, AC6 | AC2(ii), AC1, AC3, AC4, AC5, AC6 | **plan is a superset** — see note 2 |
| unassigned in doc | AC1, AC4, AC5 | (assigned to M2, with AC5 also gated at M1) | **plan assigns what the doc left unassigned** |

**Result: NO CONTRADICTION. The delta is entirely additive, and it moves no AC away from where the doc
puts it.** Specifically:

1. *Note 1* — the plan runs `gofmt`/`go vet` (AC5) and Test A's green (AC1/AC4 clauses) at M1 as well as
   at M2. The doc does not forbid this; AC5 is explicitly "green-at-base by DESIGN … listed so the
   sprint's final-tree gate list is complete" (`:509-513`), and AC1/AC4 are final-tree criteria that the
   plan merely *also* samples early. **AC5 and AC1 and AC4 are only CLOSED at M2**, on the final tree.
   The M1 occurrence is a fast-fail, not a closure. Stated here so the evaluator does not read the M1
   row as an early closure claim.
2. *Note 2* — AC1/AC4/AC5 land in M2 because M2 produces the final tree. The doc leaves them unassigned;
   assigning them to the only stage that can satisfy them is not a disagreement.

**Where the plan is narrower than the doc, and says so:** the doc's AC3 asserts "the eighteen named
mutations red the named tests" as a single undifferentiated obligation. This plan splits that obligation
across two lanes (**17 executor arms + M8 controller-only**, §6.3) on measured sandbox grounds. That is a
*routing* refinement, not a scope reduction: all 18 arms still run, all 18 results are still required,
and if M8 cannot be obtained even by the controller it is recorded `UNINFORMATIVE`, never as passed.

**Where the plan adds an obligation the doc does not have:** the cache-warming rule in §4.3/§6.3(3).
The doc could not have known it — it is a property of `cmd/ailang-worldd/cli_test.go:123` interacting
with a toolchain-switching mutation, and the planner discovered it by tripping over it.

---

## 10. Doc defects found by the planner

See the report; recorded here for the record.

- **D1 — M8 is not runnable in the executor lane, and the doc does not say so.** The doc's mutation
  table presents all 18 arms as one homogeneous set. Measured: M8 (`go.mod` → `go 1.27.0`) triggers
  `go: downloading go1.27.0 (darwin/arm64)` — a network toolchain fetch that `--sandbox workspace-write`
  denies, unpacking into `~/go/pkg/mod`, outside the workspace. **Not a design error** (the mutation is
  correct and the doc never claimed a lane), but it is a gap the plan must close, and it does: §6.3.
- **D2 — M8's arm poisons the next `verify_go.sh` run, and the doc does not warn.** Measured: after the
  M8 probe, `verify_go.sh` went **rc=1** on `TestCLIRealSubprocessEpisode: build subprocess binary:
  signal: killed` at exactly its 30 s budget, with `host/verifygate` taking **259.610 s** in the same
  transcript; the same test alone, warm, passes in **2.90 s**. Repair: run M8 last, re-warm with
  `go build ./...`, and never report that specific failure signature as a sprint red without a warm
  re-run. §4.3, §6.3.
- **D3 — `run.sh:51`, not `:50`, is the `OK*)` arm.** The doc's edit block says "the `OK*)` arm of
  `probe()`'s case (`:50–:53`)". `:50` is `case "$out" in`; `:51` is the `OK*)` arm; `:52` is `BUG*)`;
  `:53` is `esac`. The doc's range is the whole `case` construct and is not wrong, but an executor
  anchoring on "`:50`" edits the wrong line. Cosmetic; named so it costs no one a debug cycle.
- **D5 — NOT a doc defect: a REPO defect the planner tripped over, and the most valuable thing in this
  plan.** `cmd/ailang-worldd/cli_test.go` reads an `exec.Cmd.Stderr` buffer that a live `os/exec` copier
  goroutine is still writing to. `:137-138` binds `&daemonErr` as the subprocess's stderr; `:143` parks
  `cmd.Wait()` in another goroutine; `:176` (the 5 s timeout branch of the `select` at `:169-177`) then
  calls `daemonErr.String()` **before `Wait()` has returned**. The race detector reported it verbatim
  (read at `:176`, previous write from the copier started at `:139`). It is latent on the happy path,
  so it fires only when the announcement times out — i.e. under load. **Three sibling `t.Fatalf` sites
  at `:172`, `:174` and `:176` all read `daemonErr.String()`; `:174` is on the `<-waited` branch and is
  safe, `:172` and `:176` are not.** Nothing in `design_docs/` names `daemonErr` (**0** files, against a
  firing same-file control). **Recommended: file a new queue row `w-worldd-cli-stderr-buffer-race`.
  It is explicitly OUT OF SCOPE for row 41** — this sprint must not fix it, and the plan handles it by
  conditioning G10 and adding G10b/G10c (§4.3).
- **D6 — the doc's queue row 44 does not exist in the charter yet.** `w-setup-go-pin-unguarded.md:292`
  files `w-miscompile-instrument-inert-in-ci` as "queue row 44", but
  `grep -c "w-miscompile-instrument-inert-in-ci" design_docs/world-mission.md` → **0** against a firing
  same-file control (`w-floor-raise-coupling-inventory` → **1**); the charter's queue ends at row **43**.
  The doc's Deferred Scope and Test B residual comment both point at a row that is not written down.
  **Controller action, not executor action**: add it when this iteration's charter edit is made.
- **D4 — the doc's `wc -l` for `run.sh` is never stated**, and the edit block's "after `:87–:91`, before
  the RESULT banner" is the only positional anchor for the fourth guard. Planner-measured: `run.sh` is
  **94 lines**; the banner is `:92-:94`; the insertion point is between line 91 and line 92. Not a
  defect so much as a missing number; supplied.

**Nothing in the doc was found to be unimplementable.** Test A's six steps, Test B's seven steps, the
production edit and all six ACs are executable as written on this tree. D1/D2 are lane and sequencing
facts the doc had no way to measure; D3/D4 are anchor precision.

---

## 11. Scope fence (restated, because it is enforceable)

**Must not be touched by this sprint, at any stage, including drill arms:**
`host/store/toolchain_canary_test.go` · `scripts/verify_ail.sh` · `scripts/verify_go.sh` ·
`packages/world-core/` · `scripts/world_package_ready_packet.golden.json` ·
`docs/SELF_MOD_PUBLISH.md` · `host/verifygate/module_manifest_gate_test.go` ·
**`.github/workflows/ci.yml` as a production edit**.

`ci.yml` and `go.mod` are **mutated and restored** by the drill; they are not *edited*. The final
`git status --porcelain` must show **exactly two** paths (§G9). Anything else is a scope breach.

The `continue-on-error: true` question at `ci.yml:172` is **queue row 44**. Do not answer it here.
Rows 42 and 43 own `host/store/` and the floor-raise coupling inventory. Do not absorb them.

## 12. Definition of done

1. Both tests exist, both RUN (2 `=== RUN`), both PASS — AC1, with the nonsense control still saying
   `[no tests to run]`.
2. M1-e and G4 are both recorded verbatim in the implementation report — AC2, and **only** that pair
   discharges it.
3. All 18 mutation arms recorded with literal outputs, each restored byte-identically (or, for M12, with
   its mode restored and `git status` clean) — AC3, with M1/M2 explicitly confirmed as **now RED**.
4. `env -u AILANG_BIN` form green with 2 `=== RUN` — AC4.
5. `go vet` rc=0 and `gofmt -l` silent — AC5.
6. AC6's three controller results recorded: the attended `run.sh` rc=0 with the pinned line, the
   guard-trip rc=1 with no RESULT banner, and the byte-identical restore.
7. Controller gates, outside the sandbox, after a cache-warming `go build ./...`:
   `verify_ail.sh` rc=0 **with the banner byte-identical to base** (`9/9 steps`,
   `11 required identities verified, 40 named tests pass`); `go test ./host/verifygate/ -count=1` (G10b)
   rc=0; and `verify_go.sh` (G10) rc=0 — **or** rc=1 whose only failing package is `cmd/ailang-worldd`
   with one of the two §4.3 signatures, discharged by G10c green, **with both readings reported**.
8. `git status --porcelain` shows exactly the two expected paths; the controller builds **one** commit.
