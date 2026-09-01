# Sprint plan — `w-racecontrol-floor-bump-disarms-the-race-control` (queue row 48)

**Design doc:** [`design_docs/planned/w-racecontrol-floor-bump-disarms-the-race-control.md`](w-racecontrol-floor-bump-disarms-the-race-control.md)
(quorum rounds 1–3 BLOCKED at full strength; direction settled by the **attended human ruling
`D-WORLD-28`, Mark Edmondson, 2026-09-01**; round-3 objections were completeness/verification only
and were closed under the narrow-refinement carve-out. **The direction is RATIFIED and this plan
does not re-open it.**)
**Sprint base:** `dev` @ `8bb9214`, working tree clean except the design doc itself (porcelain 1,
that one path, at planning entry AND exit)
**Planner:** mission-control SPRINT-PLANNER, lane `opus fail-closed:env-pin`, in the repo working tree
**Estimate:** 0.15 day (≈1.35 h) · **4 milestones, 4 commits** · **Risk:** low
(base ≈1.0 h + ≈0.35 h for the controller-authorised D1 P1-presence needle)
**Scope:** exactly three files —
`scripts/verify_go.sh`, `host/verifygate/toolchain_pin_gate_test.go`,
`design_docs/verification/w-race-gate-blindspot/racecontrol/go.mod`. No new file, no new package, no
new CI job, no `.ail`, **no edit to `tools/launchd/*` (frozen core)**, no edit to `ci.yml`, `run.sh`,
`repro/`, or the root `go.mod`.

---

## 0. The one-paragraph version

The race-detector known-positive control at `scripts/verify_go.sh:229` is disarmed by a blanket
floor bump in `design_docs/verification/w-race-gate-blindspot/racecontrol/go.mod`. `D-WORLD-28`
ratified a three-part fix: **P1** a runtime fail-closed floor check in `verify_go.sh` after the
miscompile deny-list `esac` at `:224`; **P2** binding the race leg at `:229` to
`GOTOOLCHAIN="$ACTIVE_GO"`; **P3** a static gate
`TestRaceControlFloorStaysBelowRootToolchain` in `host/verifygate` asserting
`racecontrol floor <= root floor` plus a count=1 pin on the P2 needle; plus a `LOAD-BEARING` fence
comment in `racecontrol/go.mod` with the `go 1.22` directive byte-unchanged.
**I landed all four parts in this working tree, ran every arm the doc names (M1–M10) plus five arms
the doc does not name, and restored the tree byte-identical** (§3, §5). Four of the doc's own arms
and claims did not survive that (§4): **M6 is a vacuous arm** (0 tests run), **M10 as written is
killed by the pre-existing miscompile deny-list, not by P1** (its P1-removed control is a
byte-identical red), the doc's **"P1 presence needle"** is asserted in Conflict Surface but
implemented nowhere and the hole is measurable, and **AC7's sha256 clause is not a computable
assertion**. Repairs for all four are measured below with sole-killer controls. **The controller AUTHORISED
D1 mid-plan (§4.3): the P1 presence needle is IN SCOPE**, designed as six narrow *semantic*
assertions rather than one literal count, landed, and killed by seven named mutants (M12–M18) with
a green cosmetic control — measurements in §3.8. Everything else the doc claims — V2, V3, V6,
V9–V16, V18, M1–M5, M7, M8, and all four V14 runtime arms — **reproduced exactly, first-party, at
`8bb9214`**.

---

## 1. Environment — load-bearing, on EVERY command

```sh
export PATH=/opt/homebrew/bin:$PATH        # else rc=127 on go/gh/git/ailang
export AILANG_BIN=/tmp/ailang-v0300/ailang # AILANG v0.30.0 (commit e37b370)
```

Rig facts the executor must not rediscover:

- The shell may be `zsh`: **`${PIPESTATUS[0]}` expands to EMPTY there.** Never use it. Use
  `cmd > /tmp/out 2>&1; echo "rc=$?"` and read the file.
- Quote every glob: `--include='*.go'`, never bare.
- `scripts/verify_go.sh:17` is `cd "$(dirname "$0")/.."`. **A patched COPY executed from `/tmp`
  resolves the repo root to `/`** and dies in the tracked-binary hygiene gate long before P1. Every
  prototype copy in §3 lives in `scripts/` and is `rm`-ed afterwards.
- **`go build ./host/verifygate/` is rc=1 at base** (`no non-test Go files in …/host/verifygate` —
  it is a test-only package). Measured this session. The correct compile gate for that package is
  **`go vet ./host/verifygate/`** (rc=0 at base and after). No acceptance criterion in this plan has
  the `go build ./host/verifygate/` shape.
- `./scripts/verify_go.sh` has been observed **FLAKY** on this rig. It ran **rc=1** for me on the
  landed tree, with the failing set `{host/broker TestHandlerTimeoutKillsTheWholeProcessGroup}` under
  the `-race` leg — a **flake**, re-run 3/3 rc=0 in isolation (§3.7). One red from it is a datum, never
  a verdict; always name the failing set.

---

## 2. Lane split — what the executor may run, and what it must NOT

The executor is **`codex:gpt-5.6-sol` in a WORKSPACE-WRITE SANDBOX**. Two hard consequences:

### 2.1 Sockets

`host/daemon/daemon.go:634` calls `net.Listen("tcp", addr)` (verified this session). A
workspace-write sandbox cannot bind loopback sockets, so **any gate that reaches `host/daemon` will
fail in a way indistinguishable from a regression.** The same applies to the `exec.Command`
subprocess tests in `host/broker`, `host/capsule`, `host/archive` and to
`host/verifygate/ail_binary_gate_test.go` (which shells out to `verify_ail.sh`).

| Gate | Lane | Why |
|---|---|---|
| `go test ./host/verifygate/ -run '^TestRaceControlFloorStaysBelowRootToolchain$' -count=1 -v` | **executor** | pure `os.ReadFile` + string compare; no socket, no subprocess, no `AILANG_BIN` |
| `go vet ./host/verifygate/`, `gofmt -l host/verifygate/` | **executor** | static |
| `grep` / `awk` / `shasum` assertions | **executor** | static |
| the standalone P1 battery (§3.5) | **executor** | `awk` in a `mktemp -d`; no repo `go.mod` |
| `go test ./host/verifygate/ -count=1` (whole package) | **CONTROLLER** | `ail_binary_gate_test.go` forks `verify_ail.sh` |
| `go test ./... -count=1` | **CONTROLLER** | reaches `host/daemon`'s `net.Listen` |
| `./scripts/verify_go.sh` (and every `zz_vg_*.sh` prototype run) | **CONTROLLER** | full gate; sockets + subprocesses + `-race` |
| `M9` / `M10a'` / `M10b` runtime arms (§3.4) | **CONTROLLER** | they *are* `verify_go.sh` runs |

**If the executor runs a controller-lane gate and it reds, that red is uninterpretable.** Do not
report it, do not "fix" it — hand it to the controller.

### 2.2 No git writes

The executor makes **no** git write operation: no `add`, `commit`, `stash`, `checkout`, `mv`, `rm`,
`restore`. Instead:

- **Cumulative `.snap/M<k>/` snapshots.** After each milestone, copy the milestone's touched files
  into `.snap/M<k>/<same relative path>`; `M<k>`'s snapshot is the *cumulative* state, so the
  controller reconstructs one commit per milestone from an uncommitted tree.
- Measured this session: **`.snap/` is NOT gitignored** — `git check-ignore -v .snap/M1/x` → **rc=1**
  against the same-call control `git check-ignore -v .ailang/state/x` → **rc=0** (`.gitignore:3:**/.ailang/`).
  The instrument fires; the rc=1 is a measurement, not a silence.
- Dot-prefixed directories are skipped by the Go tool and by `gofmt`, so `.snap/` cannot red a gate.
- **All restores use `cp` from `/tmp/rc48_backup/`, verified by `shasum -a 256`.** `git checkout --` is
  forbidden.
- **M8's arm uses `mv`, never `git mv`** (the doc writes `git mv`; that is a git write and is
  replaced here).

### 2.3 Backups the executor takes FIRST, before touching anything

```sh
mkdir -p /tmp/rc48_backup
cp scripts/verify_go.sh                                                     /tmp/rc48_backup/verify_go.sh
cp go.mod                                                                   /tmp/rc48_backup/root.go.mod
cp design_docs/verification/w-race-gate-blindspot/racecontrol/go.mod        /tmp/rc48_backup/racecontrol.go.mod
cp host/verifygate/toolchain_pin_gate_test.go                               /tmp/rc48_backup/toolchain_pin_gate_test.go
shasum -a 256 /tmp/rc48_backup/*
```

**Expected, exactly (measured at `8bb9214` this session):**

```
ab782f11db0f7f259f73dd55a58eaf5a30b871bb79bd98bacbe964d50efc025b  racecontrol.go.mod
7a2983617bb9fc33747f664564fe8d8ab54fc3a177ec4dfb8c61b29ba79a7e52  root.go.mod
c71da93ed287fd690bb61a677d19f63c614ebafd2ec353508b7bc24109fae79b  toolchain_pin_gate_test.go
27eab122f4b15ac1febe0fb3aed9886d900a03d6da65377878b08843f337cd2b  verify_go.sh
```

`go.mod` and `racecontrol/go.mod` are **mutated only inside restored arms** and must be byte-identical
to the above at sprint exit. `scripts/verify_go.sh` and `toolchain_pin_gate_test.go` are the sprint's
own deliverables and change permanently.

---

## 3. What I measured before writing this plan

I landed P1 + P2 + P3 + the fence in this working tree, ran every arm, and restored. All three
protected files came back byte-identical (§5). **No expected red in this plan is predicted; every
one below is a reading I took.**

### 3.1 Base facts (doc rows V1, V3, V9–V12) — all CONFIRMED

| Reading | Measured |
|---|---|
| `git rev-parse HEAD` | `8bb9214d5bab8e58963b4f6e2506c2fd6f420f28` |
| module census (`find . -name go.mod -not -path './.git/*'`) | exactly **3**: `./go.mod`, `…/racecontrol/go.mod`, `…/repro/go.mod` |
| root `go env GOTOOLCHAIN GOVERSION` | `auto` / `go1.26.6`; `go version go1.26.6 darwin/arm64` |
| `GOTOOLCHAIN=local go env GOVERSION` | **`go1.26.4`** (below root floor `1.26.6`) |
| inside `racecontrol/`: `go env GOTOOLCHAIN GOVERSION` | `auto` / **`go1.26.4`** — nested auto-selection differs, exactly the hole P2 closes |
| `verify_go.sh` anchors | `:217 ACTIVE_GO=`, `:218 case`, `:224 esac`, `:229 go run -race .`, `:232 grep -q 'WARNING: DATA RACE'` |
| AC1 vacuity trap at base | `-run '^TestRaceControlFloorStaysBelowRootToolchain$'` → `ok … [no tests to run]`, **rc=0** — the naive rc-only form is green-at-base measuring nothing |
| name collision | `grep -rn "TestRaceControlFloorStaysBelowRootToolchain" host/` → **0** |
| V11 sibling selector (single-pipe) | **5** `=== RUN`, **5** `--- PASS`, rc=0 |
| V10 CI linkage | `ci.yml:166 run: ./scripts/verify_go.sh`; `verify_go.sh:258 go test ./... -count=1`; `go list ./...` contains `github.com/sunholo-data/ailang-world/host/verifygate` |

### 3.2 The finding, re-measured (doc rows V2, V6, V13) — all CONFIRMED

With `ACTIVE_GO` captured **at the repo root** (`go1.26.6`), from `racecontrol/`:

| Arm | Command | Measured |
|---|---|---|
| base | `GOTOOLCHAIN=go1.26.6 go run -race .` (floor `go 1.22`) | rc=1, **2** `WARNING: DATA RACE` |
| M1 disarm | floor → `go 1.27.0`, same command | rc=1, **0** races, `go: go.mod requires go >= 1.27.0 (running go 1.26.6; GOTOOLCHAIN=go1.26.6)` |
| M1 under `auto` | floor → `go 1.27.0`, `GOTOOLCHAIN=auto` | rc=1, **2** races — **the auto rescue** |
| M2 equality | floor → `go 1.26.6`, `GOTOOLCHAIN=go1.26.6` | rc=1, **2** races — `<=` is the right bound |
| V13 ARM2 | floor `go 1.26.6`, `GOTOOLCHAIN=local` | rc=1, **0** races, `… (running go 1.26.4; GOTOOLCHAIN=local)` — DISARMED |
| restore | `cp` + `shasum` | `ab782f11…`, post-restore control re-fires rc=1 / **2** races |

**Note for the executor (new, first-party):** the doc says the disarm's "only concealment is a warm
toolchain cache". That is CONFIRMED and I can now name the cache — `$(go env GOMODCACHE)/golang.org`
holds `toolchain@v0.0.1-go1.{24.9,25.6,26.0,26.2,26.3,26.5,26.6,27.0}.darwin-arm64`, which is why
even a `go 1.27.0` floor is rescued under `auto` on this rig. **Every AC3/M1 runtime reading must
therefore be taken under `GOTOOLCHAIN="$ACTIVE_GO"` (or an explicit `go1.26.6`), never under `auto`,
or the arm reads GREEN for the wrong reason.**

### 3.3 P3 static-lane arms, run against the LANDED test

Selector, on every row: `GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run
'^TestRaceControlFloorStaysBelowRootToolchain$' -count=1 -v`.

| Arm | Exact edit | Measured | Verdict |
|---|---|---|---|
| **control** | none (post-sprint tree) | rc=0, `=== RUN`=**1**, `--- PASS`=**1** | GREEN ✅ |
| **M1** | `racecontrol/go.mod` `go 1.22` → `go 1.27.0` | rc=1, RUN=1, FAIL=1, `racecontrol module floor "go1.27.0" is above the root module floor "go1.26.6"` | as doc ✅ |
| **M2** | → `go 1.26.6` | rc=0, RUN=1, PASS=1 | GREEN control, as doc ✅ |
| **M3** | → `go banana` | rc=1, RUN=1, FAIL=1, `instrument failure: racecontrol module floor "gobanana" is not a valid Go version` | as doc ✅ |
| **M4** | delete the `go 1.22` line | rc=1, RUN=1, FAIL=1, `…/racecontrol/go.mod: found 0 line(s) beginning with "go ", want exactly 1` | as doc ✅ |
| **M5** | root `go 1.26.6` → `go 1.20` | rc=1, RUN=1, FAIL=1, `racecontrol module floor "go1.22" is above the root module floor "go1.20"` | as doc ✅ (the module still loads under `go 1.20`; I checked because it need not have) |
| **M6** | root → `go banana` | **rc=1, RUN=0, PASS=0, FAIL=0**, `go: errors parsing go.mod: go.mod:3: invalid go version 'banana': must match format 1.23.0` | **REFUTED — see §4.1** |
| **M6′** *(repair)* | root → `go 1.26.6 // pin` | rc=1, RUN=1, FAIL=1, `instrument failure: root module floor "go1.26.6 // pin" is not a valid Go version` | **use this** ✅ |
| **M7** | revert P2 (`GOTOOLCHAIN="$ACTIVE_GO" go run -race .` → `go run -race .`) | rc=1, RUN=1, FAIL=1, `…/scripts/verify_go.sh: execution-binding needle "GOTOOLCHAIN=\"$ACTIVE_GO\" go run -race ." count=0, want 1` | as doc ✅ |
| **M8** | `mv racecontrol/go.mod racecontrol/go.mod.moved` | rc=1, RUN=1, FAIL=1, `open …/racecontrol/go.mod: no such file or directory` | as doc ✅ (with a **side effect**: while moved, `go list ./...` gains **1** racecontrol entry — restore immediately) |
| **M11** *(unnamed by the doc)* | delete the whole P1 block from `verify_go.sh`, keep P2 | **rc=0, RUN=1, PASS=1** | **silent GREEN — see §4.3** |

### 3.4 P1 runtime-lane arms (CONTROLLER lane), run against the LANDED block

Method: three patched **copies** of the landed `verify_go.sh` placed in `scripts/` (never `/tmp` —
line 17), each with `echo "HARNESS: reached end of race leg"; exit 0` inserted after the
`WARNING: DATA RACE` arming `fi`, so the arm terminates at the point of interest instead of running
the whole 4-minute gate. Deleted afterwards (`rm -f scripts/zz_vg_*.sh`).

| Arm | Script | Env | rc | `reached end of race leg` marker | `WARNING: DATA RACE` | FATAL text |
|---|---|---|---|---|---|---|
| **A1** GREEN | P1 present, deny live | ambient | **0** | **1** | **2** | — (prints `✓ toolchain floor gate: go1.26.6 >= root module floor go1.26.6`) |
| **A2 = M9** | P1 present, **deny NEUTERED** | `GOTOOLCHAIN=local` | **1** | **0** | **0** | `verify_go.sh: FATAL: active toolchain go1.26.4 is BELOW the root module floor go1.26.6;` |
| **A3 = discriminating CONTROL** | **P1 REMOVED**, deny NEUTERED | `GOTOOLCHAIN=local` | **0** | **1** | **2** | — |
| **A4** | P1 present, deny **LIVE** | `GOTOOLCHAIN=local` | **1** | **0** | **0** | `verify_go.sh: FATAL: active toolchain go1.26.4 miscompiles host/store/scan.go's` |

**M9's sole killer is P1**: A3 is the byte-identical scenario without the P1 block and it sails
through with the control armed. A4 shows the two guards' messages are lexically distinguishable, so
an arm's red is always attributable.

**M10 (the doc's version) and its repairs:**

| Arm | Root `go.mod` edit | `go env GOVERSION` under the mutant | P1 script | P1-REMOVED control | Verdict |
|---|---|---|---|---|---|
| **M10 (doc)** | delete the `go 1.26.6` line | **`go1.26.4`** | rc=1, marker=0, FATAL = **deny-list's** `miscompiles host/store/scan.go's` | **rc=1, byte-identical red** | **REFUTED — §4.2** |
| **M10c** | → `go banana` | **`go1.26.4`** | rc=1, FATAL = **deny-list's** | (same) | **REFUTED — §4.2** |
| **M10a′** *(repair)* | TAB-indent the `go 1.26.6` directive | `go1.26.6` | rc=1, marker=0, races=0, FATAL = `root go.mod has 0 column-0 'go ' lines, want exactly 1;` | **rc=0, marker=1, races=2** | **use this** ✅ sole killer |
| **M10b** *(repair)* | duplicate the `go 1.26.6` line | `go1.26.6` | rc=1, marker=0, races=0, FATAL = `root go.mod has 2 column-0 'go ' lines, want exactly 1;` | **rc=0, marker=1, races=2** | **use this** ✅ sole killer |

### 3.5 The P1 comparator battery (doc rows V15, V16, V18) — 15/15 CONFIRMED

Run against the block **as landed**, extracted verbatim into a standalone harness over a `mktemp -d`
whose `go.mod` is written per arm (`ACTIVE_GO` supplied as `$1`).

| Case | `ACTIVE_GO` | root floor | rc | branch |
|---|---|---|---|---|
| equality | `go1.26.6` | `go1.26.6` | 0 | `✓ toolchain floor gate` |
| above (major .Z) | `go1.27.0` | `go1.26.6` | 0 | green |
| above (patch) | `go1.26.7` | `go1.26.6` | 0 | green |
| below, deny-listed | `go1.26.4` | `go1.26.6` | **1** | `is BELOW the root module floor` |
| **below, NOT deny-listed** | **`go1.25.6`** | `go1.26.6` | **1** | `is BELOW the root module floor` — the class the deny-list structurally cannot catch |
| ordering trap (low) | `go1.9` | `go1.26.6` | **1** | below |
| ordering trap (high) | `go1.10` | `go1.9` | 0 | green — `go1.10 >= go1.9` |
| 2-vs-3 component | `go1.26` | `go1.26.0` | 0 | green |
| malformed A | `devel +abc` | `go1.26.6` | **1** | `cannot order toolchain tokens` |
| malformed B | `go1.26.6rc1` | `go1.26.6` | **1** | `cannot order toolchain tokens` |
| malformed C | `gobanana` | `go1.26.6` | **1** | `cannot order toolchain tokens` |
| malformed floor | `go1.26.6` | `go banana` | **1** | `cannot order toolchain tokens` |
| floor-read: 0 lines | `go1.26.6` | *(no `go` line)* | **1** | `has 0 column-0 'go ' lines` |
| floor-read: 2 lines | `go1.26.6` | duplicated | **1** | `has 2 column-0 'go ' lines` |
| floor-read: TAB-indented | `go1.26.6` | `\tgo 1.26.6` | **1** | `has 0 …` — shell reader and `moduleGoFloor` agree on column-0 semantics |

**V18 extraction, re-measured:** the RAW `go.mod` token `1.26.6` matches
`^go1\.[0-9]+(\.[0-9]+)?$` **0** times; the `go`-prefixed form matches **1**. The
`ROOT_FLOOR="go$(awk …)"` normalisation is load-bearing exactly as `gemini-3-1-pro` objected.
**V18 lexical fail-open control, re-measured in `bash`:** `[[ "go1.9" < "go1.26.6" ]]` is **FALSE**
and `[[ "go1.9" < "go1.10" ]]` is **FALSE** — a naive string compare would let a toolchain seventeen
minor versions below the floor through. This is why the block uses `awk`.

### 3.6 The fence (doc row V8, AC7) — CONFIRMED with one correction

Landed the doc's Example-2 fence (8 comment lines appended to the existing 3-line header, above
`module`). Measured on the fenced file:

- `grep -c 'LOAD-BEARING'` → **1** (base: **0**)
- `grep -c '^go 1.22$'` → **1**, at line **14**
- `awk '/^go /{n++} END{print n+0}'` → **1** (`moduleGoFloor` still finds exactly one floor)
- `git diff --unified=0 -- <file> | grep '^-' | grep -v '^---' | wc -l` → **0** (no line deleted or
  modified — this is the *real* "directive byte-unchanged" assertion; see §4.4)
- needle collision: `grep -cF 'GOTOOLCHAIN="$ACTIVE_GO" go run -race .'` on the fenced file → **0**
- post-fence: `TestRaceControlFloorStaysBelowRootToolchain` rc=0 RUN=1 PASS=1; the race control still
  fires rc=1 / **2** races.

### 3.7 Whole-tree gates on the landed tree (CONTROLLER lane)

| Gate | Measured |
|---|---|
| `go vet ./host/verifygate/` | **rc=0** |
| `gofmt -l host/verifygate/` | **empty** |
| `go test ./host/verifygate/ -count=1 -v` | **rc=0**, `=== RUN`=**53**, `--- PASS`=**35**, `--- FAIL`=**0**, `--- SKIP`=**0** (base, restored tree: RUN=**52**, PASS=**34**) |
| `./scripts/verify_go.sh` | **rc=1** — failing set **`{host/broker TestHandlerTimeoutKillsTheWholeProcessGroup}`** under the `-race` leg. Re-run in isolation **3/3 rc=0**. `host/broker` is untouched by this sprint. **FLAKE, not a regression.** The gate printed `✓ toolchain floor gate: go1.26.6 >= root module floor go1.26.6` at log line 23 and the race-control leg fired **2** `WARNING: DATA RACE`, i.e. P1 and P2 both behaved. |

### 3.8 The P1 presence needle (D1) — designed, landed, and killed (7 arms + 1 green control)

Selector on every static row: `GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run
'^TestRaceControlFloorStaysBelowRootToolchain$' -count=1 -v`. Runtime rows use the same patched-copy
harness as §3.4. Design rationale and the per-needle gutting analysis are in §4.3.

**Static lane — each mutant reds on its own named needle:**

| Arm | Exact edit to `scripts/verify_go.sh` | Landing assertion | Measured |
|---|---|---|---|
| control | none | — | rc=0, RUN=1, PASS=1 |
| **M12** | delete the whole P1 block | sentinel count 1→**0** | rc=1, RUN=1, FAIL=1 — `P1 block sentinel "# --- P1 (queue row 48" count=0, want 1: the D-WORLD-28 fail-closed block is absent or duplicated` (**P1a**) |
| **M13** | delete **only** the `if [ "$root_go_lines" -ne 1 ]; then … fi` branch | guard 1→**0**, in-block FATALs 3→**2** | rc=1, RUN=1, FAIL=1 — `P1 floor-read count guard ("[ \"$root_go_lines\" -ne 1 ]") count=0, want 1: a refusal branch of the D-WORLD-28 floor gate was removed or reworded` (**P1b**) |
| **M14** | swap the comparator operands | correct call 1→**0**, swapped 0→**1** | rc=1, RUN=1, FAIL=1 — `P1 comparator call "go_version_ge \"$ACTIVE_GO\" \"$ROOT_FLOOR\"" count=0, want 1` (**P1d**) |
| **M15** | `ROOT_FLOOR="go$(awk …)"` → `ROOT_FLOOR="go1.26.6"` | hardcoded 0→**1**, derived 1→**0** | rc=1, RUN=1, FAIL=1 — `P1 block contains hardcoded Go version literal(s) [go1.26.6]` (**P1e**) |
| **M16** | below-floor branch `exit 1 ;;` → `exit 0 ;;` | in-block `exit 0` 2→**3** | rc=1, RUN=1, FAIL=1 — `P1 block has 2 \`exit 1\` statements, want 3 (one per refusal)` (**P1c**) |
| **M18** | move the whole block **below** the race leg | sentinel line 226→**231**, race leg 264→**229** | rc=1, RUN=1, FAIL=1 — `P1 floor gate is out of order (deny-list@11399, P1@11963, race leg@11915)` (**P1f**) |
| **M17 — GREEN CONTROL** | reword a **comment** word inside the block (`never hardcoded.` → `(never hardcoded, see D-WORLD-28).`) | reworded 0→**1** | **rc=0, RUN=1, PASS=1** — the set is not "any edit reds"; it is not row 49's token count |

**Runtime lane — why the needle is the sole killer, measured.** Same six variants, run as patched
copies with the `HARNESS: reached end of race leg` marker:

| Variant | **AMBIENT** (what CI runs: `ACTIVE_GO=go1.26.6`) | **M9 conditions** (deny neutered, `GOTOOLCHAIN=local`, `ACTIVE_GO=go1.26.4`) |
|---|---|---|
| correct block | rc=0 marker=1 races=2 | rc=1 marker=0 races=0, `is BELOW the root module floor` |
| M13 (branch deleted) | **rc=0 marker=1 races=2** | rc=1 — but on the *below-floor* branch; the arm cannot tell M13 from correct |
| **M14 (operands swapped)** | **rc=0 marker=1 races=2** | **rc=0 marker=1 races=2 — FAIL OPEN.** Even the hostile arm goes green |
| M15 (floor hardcoded) | **rc=0 marker=1 races=2** | rc=1 — identical to correct today; only a future root-floor bump would diverge |
| **M16 (`exit 0`)** | **rc=0 marker=1 races=2** | **rc=0 marker=0 — FALSE GREEN.** The gate prints its own FATAL and then exits **success**, having skipped the race control and everything after it |
| M18 (block moved) | **rc=0 marker=1 races=2** | rc=1 — but the race control was already invoked under the unvetted toolchain |

**Under the ambient toolchain — the only condition CI ever runs — all six gutted variants are
byte-identical in behaviour to the correct block.** The runtime lane is blind to every one of them in
CI. For **M14** and **M16** the static needle is the sole killer *in the entire sprint*, hostile arm
included.

**Two residuals, measured not asserted** (carried into §9.1): a mutant that preserves every byte of
text but severs the dispatch (`case "$floor_rc" in` → `case "0" in`, **R1**) and a mutant that inverts
the comparator's verdicts inside the `awk` program (**R2**) both **PASS** the static needle set
(rc=0, RUN=1, PASS=1) — a static scan cannot see either. Both are caught by the **M9 runtime arm**,
whose expected red *vanishes* under them (measured: R1 and R2 each give `rc=0 marker=1 races=2` under
M9 conditions, where the correct block gives rc=1/marker=0).

---

## 4. Refutations — corrections that are BINDING on the executor

### 4.1 REFUTED — M6 is a vacuous arm; it runs **zero** tests

**Doc claim** (Non-Vacuity table, M6): *root `go.mod` `go 1.26.6` → `go banana` ⇒ RED
`root-floor-validity floor: root floor "gobanana" invalid`.*

**Measured:**

```
rc=1  === RUN=0  --- PASS=0  --- FAIL=0
go: errors parsing go.mod:
go.mod:3: invalid go version 'banana': must match format 1.23.0
```

Go's own module loader rejects the mutated root `go.mod` before the test binary is built, so
`TestRaceControlFloorStaysBelowRootToolchain` never runs. Under the doc's **own AC10 rule** — a
go-test arm's evidence is a *counted* `=== RUN`/`--- PASS`, never a bare rc — an arm with `=== RUN`=0
is not evidence. As written, M6 leaves the test's `!version.IsValid(rootFloor)` branch with **no
killer**.

**Repair, measured (M6′):** root `go 1.26.6` → **`go 1.26.6 // pin`**. Go accepts a trailing comment,
so the module loads and the test runs; `moduleGoFloor`'s `canonicalizeVersionPin` yields the token
`"go1.26.6 // pin"`, which `version.IsValid` rejects:

```
rc=1  === RUN=1  --- PASS=0  --- FAIL=1
    toolchain_pin_gate_test.go:574: instrument failure: root module floor "go1.26.6 // pin" is not a valid Go version
```

M3 (the *racecontrol*-side `go banana`) is **not** affected: `racecontrol/go.mod` is a nested module
that `go test ./host/verifygate/` never loads. Measured: M3 runs and reds on its named branch.

### 4.2 REFUTED — M10 as written is killed by the miscompile deny-list, not by P1

**Doc claim** (Non-Vacuity table, M10; Risks table; AC11): *delete the `go 1.26.6` line from the root
`go.mod` ⇒ the P1 floor-read check FATALs before `:229`*, and *"M9 reds the below-floor refusal branch
solely; M10 reds the unreadable-floor instrument branch solely — together they cover both of P1's
branches with no overlap."*

**Measured — the opposite:** deleting the root `go` directive changes which toolchain
`GOTOOLCHAIN=auto` selects. The root `go.mod` floor is *what pushes `auto` up to `go1.26.6`*; remove
it and `go env GOVERSION` falls back to the local base:

```
col0 'go ' lines = 0     go env GOVERSION = go1.26.4
M10 (P1 present):   rc=1 marker=0  FATAL: active toolchain go1.26.4 miscompiles host/store/scan.go's
M10 (P1 REMOVED) :  rc=1 marker=0  FATAL: active toolchain go1.26.4 miscompiles host/store/scan.go's
```

The P1-removed control is a **byte-identical red**. This is precisely the failure the doc's own
V14/A3 control shape exists to catch, and M10 fails it: **the pre-existing miscompile deny-list
supplies M10's red, not P1.** `go banana` in the root `go.mod` (M10c) fails the same way for the same
reason. M10's "the race leg is unreached" clause is satisfied by both — which is exactly why
unreached-ness alone is not attribution.

**Repairs, both measured with sole-killer controls** (§3.4 table): **M10a′** TAB-indent the root
directive, and **M10b** duplicate it. Both keep `go env GOVERSION = go1.26.6` (deny-list passes),
both red on P1's own attributed message, and **both P1-removed controls are rc=0 / marker=1 / races=2**.

**Additional correction to the same claim:** P1 has **three** refusal branches, not two — floor-read
count≠1, below-floor ordering, and malformed-token (`exit 2`). The **malformed-token branch is not
reachable through the live script on this rig by any root-`go.mod` mutation**, because every `go`
directive Go itself rejects also changes `go env GOVERSION` into the deny-listed base (above). It is
reachable only from the `ACTIVE_GO` side (`devel`, `go1.26.6rc1`), which cannot be installed here. Its
killer is therefore the **standalone battery of §3.5**, and the plan labels it as such rather than
banking it as a runtime-lane arm.

### 4.3 REFUTED, and now FIXED — the doc's "P1 presence needle" was asserted and implemented nowhere

**Doc claim** (Conflict Surface): *"The new test scans two needles: the P2 execution-binding substring
`GOTOOLCHAIN="$ACTIVE_GO" go run -race .` (count=1) **and a P1 presence needle that the failure-open
deletion reds**."*

No such needle appeared in the Solution Design's component 2, in any AC, or in the M1–M10 table, and
there was no mutant for it. **Measured (M11):** with P2 and P3 intact and the entire P1 block deleted
from `scripts/verify_go.sh`:

```
static gate: rc=0  === RUN=1  --- PASS=1
runtime lane: rc=0  marker=1  races=2
```

Deleting the block that `D-WORLD-28` mandates was **green in both lanes**.

#### 4.3.1 Disposition — D1 AUTHORISED by the controller

The controller authorised implementation, on the record, for three reasons that belong in this plan
rather than in a chat log:

1. **The doc already promises the needle.** Shipping without it would not defer scope; it would leave
   the design doc asserting something false about its own deliverable. Implementing closes an
   inconsistency; declaring it residual would require *editing the doc to withdraw a claim* — strictly
   more work and strictly worse.
2. **P2 already has exactly this protection.** M7 pins the execution-binding needle, with the stated
   rationale *"the machinery that binds execution must not silently revert."* P1 and P2 land in the
   same commit, in the same file, for the same reason. Protecting the edit the human did **not** have
   to rule on while leaving unprotected the one they **did** rule on is an asymmetry with no argument
   behind it.
3. **This mission has already paid for this class twice**, and both are in the charter: **row 49** (a
   token count blind to shape-gutting — 22 of 23 mutations walked through it) and **row 59** (a static
   grep cannot prove an assertion is live).

#### 4.3.2 The constraint, and how the needle is designed to not become a third instance

A bare `strings.Count(src, "some literal") == 1` **is** row 49's defect. The needle set therefore
asserts **load-bearing semantics** of the block — things a mutant cannot delete without deleting the
behaviour — as **six narrow assertions**, not one broad one. The block is first split into its
comparator half (`go_version_ge`'s body) and its gate half, so an assertion about one cannot be
satisfied by text in the other.

| # | Assertion | What it is load-bearing for | What a mutant would have to do to satisfy it **while** gutting the block |
|---|---|---|---|
| **P1a** | the block's sentinel comment and its `✓` success line each occur exactly once, sentinel before success | block PRESENCE and delimitation | keep both markers and gut only the interior — then P1b–P1e fire |
| **P1b** | the gate half contains, once each: the guard expression `[ "$root_go_lines" -ne 1 ]`, the stem `is BELOW the root module floor`, the stem `cannot order toolchain tokens`; and **exactly 3** `verify_go.sh: FATAL:` refusals | the three *distinct* failure modes exist and are separately attributed | keep all three branch texts but sever their dispatch — that is residual **R1** (§9.1), which the M9 runtime arm catches |
| **P1c** | the gate half contains **exactly 3** `exit 1`, and **zero** `exit 0` | every refusal actually refuses | print a FATAL and continue → the `exit 1` count drops → red. There is no way to keep 3 exits and fail open textually |
| **P1d** | the comparator half contains all of `exit 0`, `exit 1`, `exit 2`; and the gate half calls `go_version_ge "$ACTIVE_GO" "$ROOT_FLOOR"` **exactly once, in that operand order** | the three-way contract (>= / < / not-a-release-version) and the **direction** of the comparison | swapping the operands inverts the gate into a fail-open that **even the M9 arm cannot see** (measured, §3.8) — so P1d has no substitute anywhere else in the sprint |
| **P1e** | the block contains **zero** concrete `go1.<d>.<d>` literals, and does contain `awk '/^go /{print $2; exit}' go.mod` | the floor is READ, never hardcoded (a literal goes stale on the next root-floor bump) | hardcode today's value → the negative assertion fires. A negative assertion has no "satisfy it with decoration" move |
| **P1f** | byte offsets satisfy `deny-list case < P1 sentinel < P2 needle` | P1 runs **after** the deny-list and **before** the race control | move the block below the race leg → red. Position cannot be faked by adding text elsewhere |

Two of these (P1c, P1e) are **negative** assertions — "this must NOT appear" — which is the specific
shape row 49's token count lacked. P1f is **positional**. P1d pins an **operand order**. None of the
six is satisfiable by a comment.

#### 4.3.3 Measured: the needles kill, and they do not over-kill

Full readings in §3.8. Summary: **M12** (whole-block deletion) → P1a. **M13** (delete one refusal
branch, leave the rest) → P1b. **M14** (swap comparator operands) → P1d. **M15** (hardcode the floor)
→ P1e. **M16** (below-floor branch `exit 1` → `exit 0`) → P1c. **M18** (move the block below the race
leg) → P1f. **M17** (reword a comment word inside the block) → **GREEN** — the anti-brittleness
control that proves the set is not "any edit reds".

**And the finding that settles why the needle is not optional:** under the **ambient** toolchain — the
condition CI actually runs under — **all six gutted variants are `rc=0, marker=1, races=2`, byte-for-byte
indistinguishable from the correct block.** The runtime lane is completely blind to every one of them
in CI. For **M14** and **M16** even the deliberately hostile M9 arm goes green (`rc=0 marker=1 races=2`
and `rc=0 marker=0` respectively — M16 exits **success** after printing its own FATAL). For those two,
the static needle is the **sole killer in the entire sprint**.

### 4.4 REFUTED — AC7's sha256 clause is not a computable assertion

**Doc claim** (AC7): *"`shasum -a 256 …/racecontrol/go.mod` → the post-comment hash is the base hash
plus the comment block."*

Cryptographic hashes do not compose; there is no operation that turns `ab782f11…` "plus a comment
block" into a checkable value, and the AC as written cannot be executed. **Replacement, measured
(§3.6):** `git diff --unified=0 -- <file> | grep '^-' | grep -v '^---' | wc -l` → **0**, together with
`grep -n '^go 1.22$'` → line **14** and `awk '/^go /{n++} END{print n+0}'` → **1**. Those three are
executable and jointly assert "only `//` lines were added; the directive did not move or change".

### 4.5 Smaller corrections (not blocking, but the executor should not trip on them)

| # | Doc text | Measured |
|---|---|---|
| a | *Files to Modify:* `scripts/verify_go.sh` **(+~8 LOC net)** | the verbatim P1 block is **35 lines** (`:225–:259` after insertion) plus **1** altered line at the old `:229` (new `:264`). A ~4× understatement; the executor's diff review must expect +35/−0 and +1/−1. |
| b | AC9 / V10: *`grep -n 'go test ./\.\.\.' scripts/verify_go.sh` → **a** match at `:258`* | **3** matches: `:5` (header comment), `:256` (`echo`), `:258` (the executable line). Bind the AC to "≥1 match whose line is not a comment and not an `echo`", or grep `'^go test \./\.\.\. -count=1$'` → exactly **1** at `:258`. |
| c | AC2: *"M9/M10 red P1's refusal and instrument floors in the runtime lane **(V13)**"* | V13 is the `racecontrol` ARM row; it contains no P1 measurement. The correct citations are V14 (M9) and V16 (branches). Citation error only. |
| d | Non-Vacuity M8: *`git mv racecontrol/go.mod racecontrol/go.mod.moved`* | `git mv` is a git write; the executor is forbidden from git writes. Use `mv`. Measured equivalent, same red. |
| e | AC3 / Example 1 attack sketch runs `GOTOOLCHAIN=go1.26.6 go run -race .` after the bump | correct as written, and **load-bearing**: under `GOTOOLCHAIN=auto` this rig's warm toolchain cache (which contains `go1.27.0`) makes the same bump fire **2** races. Never take an M1 runtime reading under `auto`. |

### 4.6 What the doc got RIGHT (measured, not assumed)

V1 census, V2, V3, V6, V8, V9, V10, V11, V12, V13 all reproduced verbatim. The whole V14 battery
(A1/A2/A3/A4) reproduced verbatim, including the discriminating control that makes M9 attributable.
V15/V16/V18 reproduced 15/15, including the `go1.25.6` case the deny-list cannot catch, the `go1.9`
vs `go1.10` ordering trap, the `go` normalisation, and both lexical fail-open controls. M1, M2, M3,
M4, M5, M7, M8 all red (or green, for M2) on exactly the branch and with exactly the message the doc
names. The exact P1 shell block in the doc's Implementation Plan is **correct as written** — I landed
it verbatim and could not measure a defect in it; every one of its branches carries a distinct
message and no two branches share one.

---

## 5. Land-and-restore discipline — what I did, and what the executor repeats

Every arm in §3 followed: `shasum -a 256` before → land the mutant → **assert the mutant LANDED**
(old literal count → 0, new literal count → 1, or the equivalent file-existence reading) → measure →
restore by `cp` → `shasum -a 256` byte-identical → re-run the pristine control after every batch.

**Planner exit state, measured:**

```
27eab122f4b15ac1febe0fb3aed9886d900a03d6da65377878b08843f337cd2b  scripts/verify_go.sh                                        ✅
7a2983617bb9fc33747f664564fe8d8ab54fc3a177ec4dfb8c61b29ba79a7e52  go.mod                                                      ✅
ab782f11db0f7f259f73dd55a58eaf5a30b871bb79bd98bacbe964d50efc025b  design_docs/verification/…/racecontrol/go.mod               ✅
c71da93ed287fd690bb61a677d19f63c614ebafd2ec353508b7bc24109fae79b  host/verifygate/toolchain_pin_gate_test.go                  ✅
git status --porcelain → only ` M design_docs/planned/w-racecontrol-floor-bump-disarms-the-race-control.md`
                          (pre-existing at planning entry; not authored by me)
scripts/zz_vg_*.sh → removed
post-restore control → `-run '^TestRaceControlFloorStaysBelowRootToolchain$'` back to `[no tests to run]`
```

---

## 6. Decisions the controller owns before M1 starts

| # | Decision | Default | Evidence |
|---|---|---|---|
| **D1** | The P1 presence needle (§4.3) | **AUTHORISED by the controller — IMPLEMENTED.** Six semantic assertions (P1a–P1f), seven mutants (M12–M18), one green cosmetic control (M17) | M11 measured rc=0 / RUN=1 / PASS=1 with P1 deleted; §3.8 for the needle arms |
| **D2** | Adopt **M6′**, **M10a′**, **M10b** in place of the doc's M6 and M10 | **adopt** — the doc's versions are measurably not arms | §4.1, §4.2 |
| **D3** | Amend AC7's sha256 clause to the diff-shape check | **amend** | §4.4 |

D2 and D3 are corrections of measured defects and do **not** touch the ratified `D-WORLD-28`
direction. **D1 is settled: authorised and implemented** — it enlarges P3 only, and enforces the
ruling's own block rather than re-opening its design.

---

## 7. Milestones

Every go-test acceptance criterion below binds a **counted** number of top-level `=== RUN` and
`--- PASS` lines. **A bare rc=0 is never evidence** — `go test -run <selector-that-matches-nothing>`
exits **0** printing `ok … [no tests to run]` (measured at base this session, and again with the
nonsense selector `TestNoSuchRaceControlFloorTestZZZ`). Every count below is a reading I took on the
landed tree, not a prediction.

Counting recipe used everywhere (zsh-safe):

```sh
GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run '<SEL>' -count=1 -v > /tmp/out.txt 2>&1; echo "rc=$?"
echo "RUN=$(grep -c '^=== RUN' /tmp/out.txt) PASS=$(grep -c '^--- PASS' /tmp/out.txt) FAIL=$(grep -c '^--- FAIL' /tmp/out.txt)"
```

---

### Milestone M1 — P3: the static gate (executor lane)

**File:** `host/verifygate/toolchain_pin_gate_test.go` · **+116 lines, −0** (measured; appended.
**No new import** — `go/version`, `os`, `path/filepath`, `regexp`, `strings`, `testing` are all already
imported at `:3–15`, `regexp` at `:11`)

**Do:** append `TestRaceControlFloorStaysBelowRootToolchain`. Shape (I landed and measured exactly
this):

1. `raceFloor := moduleGoFloor(t, filepath.Join(repoRoot, "design_docs", "verification", "w-race-gate-blindspot", "racecontrol", "go.mod"))`
2. `if !version.IsValid(raceFloor) { t.Fatalf("instrument failure: racecontrol module floor %q is not a valid Go version", raceFloor) }`
3. `rootFloor := moduleGoFloor(t, filepath.Join(repoRoot, "go.mod"))`
4. `if !version.IsValid(rootFloor) { t.Fatalf("instrument failure: root module floor %q is not a valid Go version", rootFloor) }`
5. `if version.Compare(raceFloor, rootFloor) > 0 { t.Fatalf(<names the exact consequence: the control refuses before it can fire and verify_go.sh FATALs "the race detector is not armed" for the wrong reason>) }`
6. read `filepath.Join(repoRoot, "scripts", "verify_go.sh")` **once**; **P2 needle** —
   ``const p2Needle = `GOTOOLCHAIN="$ACTIVE_GO" go run -race .` ``; record `p2At := strings.Index(...)`;
   `if got := strings.Count(src, p2Needle); got != 1 { t.Fatalf(…count=%d, want 1…) }`
7. **P1 needle set (D1, authorised — rationale and gutting analysis in §4.3, measurements in §3.8).**
   Add a `p1FloorGateRegion(t, src) (comparator, gate string, start int)` helper that locates the block
   by its sentinel `"# --- P1 (queue row 48"` and its success line ``echo "   ✓ toolchain floor gate:``
   (each **exactly once**, sentinel first — else `t.Fatalf`), then splits the region at the first line
   that is exactly `}` into the comparator half and the gate half. Then assert:
   - **P1b** in the *gate* half, count==1 each: `[ "$root_go_lines" -ne 1 ]`,
     `is BELOW the root module floor`, `cannot order toolchain tokens`; and
     `strings.Count(gate, "verify_go.sh: FATAL:") == 3`
   - **P1c** `strings.Count(gate, "exit 1") == 3` **and** `strings.Count(gate, "exit 0") == 0`
   - **P1d** the *comparator* half contains all of `exit 0`, `exit 1`, `exit 2`; and the gate half
     contains ``go_version_ge "$ACTIVE_GO" "$ROOT_FLOOR"`` **exactly once** (operand order is the gate's
     direction — swapped, it accepts every below-floor toolchain)
   - **P1e** `regexp.MustCompile(`+"`"+`go1\.[0-9]+\.[0-9]+`+"`"+`)` finds **zero** matches in
     `comparator+gate` (verified premise: the comparator's own `^go1\.[0-9]+(\.[0-9]+)?$` grammar does
     **not** match this pattern — `[` is not a digit), and the gate half contains
     ``awk '/^go /{print $2; exit}' go.mod``
   - **P1f** `denyAt := strings.Index(src, `+"`"+`case "$ACTIVE_GO" in`+"`"+`)`; require
     `denyAt < p1At && p1At < p2At`, with an instrument-failure `t.Fatalf` if `denyAt < 0`

   Every message must name the consequence, not just the count — these are the strings M12–M18 assert
   against in §3.8.

**Notes:** `repoRoot` is a **package-level variable** set in `ail_binary_gate_test.go:27`
(`findRepoRoot()`), not a function — the doc's "reusing `repoRoot`" is right, but do not write
`repoRoot()`. The test performs no `requirePinned`, reads no `AILANG_BIN`, opens no socket.

> **M1 runs BEFORE P2/P1 land, so steps 6–7 are RED on arrival** (needle count=0, sentinel count=0). That is correct and is
> the point: M1's own AC is stated at the end of **M2**, when both halves exist. If the executor
> prefers a green boundary at M1, land M2 first — the plan tolerates either order, but the counts
> below are the *post-M2* readings.

**Verification (executor lane), after M2 lands:**

| Command | Expected |
|---|---|
| `gofmt -l host/verifygate/` | **empty** |
| `go vet ./host/verifygate/` | **rc=0** |
| AC1: `-run '^TestRaceControlFloorStaysBelowRootToolchain$' -count=1 -v` | **rc=0, RUN=1, PASS=1, FAIL=0** |
| AC1 vacuity control: `-run 'TestNoSuchRaceControlFloorTestZZZ'` | **rc=0** and output contains `[no tests to run]` — **treat as RED if it ever shows RUN≥1** |
| AC5: `env -u AILANG_BIN GOTOOLCHAIN=go1.26.6 go test … -run '^TestRaceControlFloorStaysBelowRootToolchain$' -count=1 -v` | **rc=0, RUN=1, PASS=1** |
| collision: `grep -c 'func TestRaceControlFloorStaysBelowRootToolchain' host/verifygate/toolchain_pin_gate_test.go` | **1** |
| regression: `-run 'TestReproModuleFloorStaysBelowKnownBadToolchains\|TestCanaryDeclaresPositiveArmOnly\|TestMiscompileInstrumentProbesPinnedToolchain\|TestMiscompileInstrumentStepIsGatedInCI\|TestGoToolchainPinsAgreeAndMatchJobList' -count=1 -v` | **rc=0, RUN=5, PASS=5** (single-pipe alternation — an escaped `\|` is a *literal pipe* to Go's regexp and runs **0** tests at rc=0) |

**Snapshot:** `mkdir -p .snap/M1/host/verifygate && cp host/verifygate/toolchain_pin_gate_test.go .snap/M1/host/verifygate/` → **Commit 1**.

---

### Milestone M2 — P1 + P2 in `scripts/verify_go.sh` (executor edits; CONTROLLER verifies)

**File:** `scripts/verify_go.sh` · **+35 / −0** (P1) and **+1 / −1** (P2)

**P1:** insert the design doc's **exact** block, verbatim and unmodified, immediately after the
miscompile deny-list `esac` at `:224`. I landed it byte-for-byte and could not measure a defect;
do not "improve" it. After insertion it occupies **`:225`–`:259`** (`:225` blank, `:226–229` comment,
`:230–239` `go_version_ge()`, `:240–245` the count check, `:246` `ROOT_FLOOR=`, `:247–249` the guarded
call, `:250–258` the `case`, `:259` the `✓` line). Offsets measured on the landed copy.

**P2:** at the old `:229` (new `:264`), change

```
race_control_output="$(cd design_docs/verification/w-race-gate-blindspot/racecontrol && go run -race . 2>&1)"
```
to
```
race_control_output="$(cd design_docs/verification/w-race-gate-blindspot/racecontrol && GOTOOLCHAIN="$ACTIVE_GO" go run -race . 2>&1)"
```

**Do NOT touch** `:217` (`ACTIVE_GO=`), `:218–224` (the deny-list), or the arming/FATAL block that
follows the race leg.

**Verification — executor lane (static only):**

| Command | Expected |
|---|---|
| `bash -n scripts/verify_go.sh; echo "rc=$?"` | **rc=0** |
| `grep -cF 'GOTOOLCHAIN="$ACTIVE_GO" go run -race .' scripts/verify_go.sh` | **1** |
| `grep -c 'racecontrol && go run -race \.' scripts/verify_go.sh` | **0** (the bare form is gone) |
| `grep -c 'toolchain floor gate' scripts/verify_go.sh` | **1** |
| `grep -c '# --- P1 (queue row 48' scripts/verify_go.sh` | **1** (P1a sentinel) |
| `awk '/# --- P1 \(queue/,/toolchain floor gate/' scripts/verify_go.sh \| grep -c 'verify_go.sh: FATAL:'` | **3** (P1b) |
| `awk 'NR>=240 && NR<=259' scripts/verify_go.sh \| grep -c 'exit 1'` / `… 'exit 0'` | **3** / **0** (P1c) |
| `grep -c 'go_version_ge "$ACTIVE_GO" "$ROOT_FLOOR"' scripts/verify_go.sh` | **1** (P1d operand order) |
| `awk '/# --- P1 \(queue/,/toolchain floor gate/' scripts/verify_go.sh \| grep -cE 'go1\.[0-9]+\.[0-9]+'` | **0** (P1e — no hardcoded floor) |
| `grep -n 'ACTIVE_GO=$(go env GOVERSION)' scripts/verify_go.sh` | **`217:`** (unmoved) |
| `sed -n '218,224p' scripts/verify_go.sh \| shasum -a 256` | must equal the same range from `/tmp/rc48_backup/verify_go.sh` (deny-list untouched) |
| `awk 'END{print NR}' scripts/verify_go.sh` | **311** (base **276** + 35) |
| AC1 re-run (§M1 table) | **rc=0, RUN=1, PASS=1** — the needle now exists |

**Verification — CONTROLLER lane (sandbox-hostile; executor must NOT run these):**

| Arm | Expected |
|---|---|
| `./scripts/verify_go.sh` | prints `   ✓ toolchain floor gate: go1.26.6 >= root module floor go1.26.6` before the `── race-detector known-positive control` header; the race leg emits **2** `WARNING: DATA RACE`. **Treat rc=1 as a datum and NAME the failing set** — `{host/broker TestHandlerTimeoutKillsTheWholeProcessGroup}` under the `-race` leg is the known flake (3/3 rc=0 in isolation). |
| A1 GREEN (patched copy in `scripts/`, marker + `exit 0` after the arming `fi`) | **rc=0, marker=1, races=2** |
| **M9** = A2 (deny-list case body neutered to `: ;;`, `GOTOOLCHAIN=local`) | **rc=1, marker=0, races=0**, FATAL `active toolchain go1.26.4 is BELOW the root module floor go1.26.6;` |
| **M9's discriminating control** = A3 (same, P1 block deleted) | **rc=0, marker=1, races=2** — this is what makes M9 attributable to P1 alone |
| A4 (deny-list live, `GOTOOLCHAIN=local`) | **rc=1, marker=0**, FATAL `miscompiles host/store/scan.go's` — distinguishable message |
| **M10a′** (root directive TAB-indented, P1 script, ambient) | **rc=1, marker=0, races=0**, FATAL `root go.mod has 0 column-0 'go ' lines, want exactly 1;` |
| **M10a′ control** (same, P1 deleted) | **rc=0, marker=1, races=2** |
| **M10b** (root directive duplicated, P1 script, ambient) | **rc=1, marker=0, races=0**, FATAL `root go.mod has 2 column-0 'go ' lines, want exactly 1;` |
| **M10b control** (same, P1 deleted) | **rc=0, marker=1, races=2** |
| **DO NOT RUN** the doc's M10 (delete the root `go` line) or M10c (`go banana`) as P1 arms | both are killed by the deny-list; their P1-removed controls are byte-identical reds (§4.2) |

All patched copies live in **`scripts/`** (never `/tmp` — line 17) and are removed with
`rm -f scripts/zz_vg_*.sh` before the snapshot. Root `go.mod` restores by `cp` from
`/tmp/rc48_backup/root.go.mod`; assert `7a2983617bb9fc33747f664564fe8d8ab54fc3a177ec4dfb8c61b29ba79a7e52`.

**Executor-lane substitute for the P1 branch battery** (§3.5) — this one *is* sandbox-safe, because it
runs the extracted block over a `mktemp -d`, never the repo:

| Case | Expected rc / branch |
|---|---|
| `go1.26.6` vs floor `go1.26.6` | 0 / `✓ toolchain floor gate` |
| `go1.27.0`, `go1.26.7` vs `go1.26.6` | 0 / green |
| `go1.26.4`, **`go1.25.6`**, `go1.9` vs `go1.26.6` | **1** / `is BELOW the root module floor` |
| `go1.10` vs floor `go1.9` | 0 / green (ordering trap) |
| `go1.26` vs floor `go1.26.0` | 0 / green |
| `devel +abc`, `go1.26.6rc1`, `gobanana`; or floor `go banana` | **1** / `cannot order toolchain tokens` |
| floor file with 0 / 2 / TAB-indented `go ` lines | **1** / `has 0 …` / `has 2 …` / `has 0 …` |
| `bash -c '[[ "go1.9" < "go1.26.6" ]]'` | **rc=1 (FALSE)** — the fail-open the block exists to avoid |

**Snapshot:** `mkdir -p .snap/M2/scripts .snap/M2/host/verifygate && cp scripts/verify_go.sh .snap/M2/scripts/ && cp host/verifygate/toolchain_pin_gate_test.go .snap/M2/host/verifygate/` → **Commit 2**.

---

### Milestone M3 — the `LOAD-BEARING` fence (executor lane)

**File:** `design_docs/verification/w-race-gate-blindspot/racecontrol/go.mod` · **+9 / −0**

**Do:** append the doc's Example-2 fence (8 `//` lines + one `//`-only separator) to the **existing**
3-line header, immediately above `module ailang-world/verification/race-detector-control`. The
`go 1.22` line is **not moved and not edited**; it ends at file line **14**.

**Verification (executor lane):**

| Command | Expected |
|---|---|
| `grep -c 'LOAD-BEARING' <file>` | **1** (base **0** — the fence is AC7-red at base by design) |
| `grep -c '^go 1.22$' <file>` | **1** |
| `grep -n '^go 1.22$' <file>` | **`14:go 1.22`** |
| `awk '/^go /{n++} END{print n+0}' <file>` | **1** (`moduleGoFloor` still sees exactly one floor) |
| **AC7 replacement (§4.4):** `git diff --unified=0 -- <file> \| grep '^-' \| grep -v '^---' \| wc -l` | **0** — nothing deleted or modified, only `//` lines added |
| needle collision: `grep -cF 'GOTOOLCHAIN="$ACTIVE_GO" go run -race .' <file>` | **0** |
| AC1 re-run | **rc=0, RUN=1, PASS=1** |
| CONTROLLER: `(cd <racecontrol> && GOTOOLCHAIN=go1.26.6 go run -race .)` | **rc=1, 2× `WARNING: DATA RACE`** |

**Snapshot:** cumulative into `.snap/M3/` → **Commit 3**.

---

### Milestone M4 — mutation battery, hygiene, closure (mixed lanes)

**Files:** none permanently. Every edit here is a restored arm.

**Static-lane arms (executor).** For each: `shasum -a 256` before → land → **assert landed** → run the
AC1 selector → `cp` restore → `shasum -a 256` byte-identical → re-run the pristine control.

| Arm | Edit | Landing assertion | Expected |
|---|---|---|---|
| M1 | `racecontrol` `go 1.22` → `go 1.27.0` | old=0, new=1 | rc=1, RUN=1, FAIL=1, `…"go1.27.0" is above the root module floor "go1.26.6"` |
| **M2 (GREEN control)** | → `go 1.26.6` | old=0, new=1 | **rc=0, RUN=1, PASS=1** |
| M3 | → `go banana` | old=0, new=1 | rc=1, RUN=1, FAIL=1, `racecontrol module floor "gobanana" is not a valid Go version` |
| M4 | delete the `go 1.22` line | col0 `go ` count=0 | rc=1, RUN=1, FAIL=1, `found 0 line(s) beginning with "go ", want exactly 1` |
| M5 | root `go 1.26.6` → `go 1.20` | old=0, new=1 | rc=1, RUN=1, FAIL=1, `"go1.22" is above the root module floor "go1.20"` |
| **M6′** (replaces M6) | root → `go 1.26.6 // pin` | `grep -n '^go '` shows the comment | rc=1, **RUN=1**, FAIL=1, `root module floor "go1.26.6 // pin" is not a valid Go version` |
| M7 | revert P2 in `verify_go.sh` | needle=0, bare=1 | rc=1, RUN=1, FAIL=1, `execution-binding needle … count=0, want 1` |
| M8 | **`mv`** (not `git mv`) `racecontrol/go.mod` → `go.mod.moved` | file absent | rc=1, RUN=1, FAIL=1, `open …/racecontrol/go.mod: no such file or directory`; **restore immediately** (`go list ./...` gains a racecontrol entry while moved) |
| **M12** *(D1)* | delete the whole P1 block from `verify_go.sh` | sentinel 1→0 | rc=1, RUN=1, FAIL=1, `P1 block sentinel … count=0, want 1` (**P1a**) |
| **M13** *(D1, single-branch)* | delete **only** `if [ "$root_go_lines" -ne 1 ]; then … fi` | guard 1→0, in-block FATALs 3→2 | rc=1, RUN=1, FAIL=1, `P1 floor-read count guard … count=0, want 1` (**P1b**) |
| **M14** *(D1)* | swap the comparator operands | correct call 1→0, swapped 0→1 | rc=1, RUN=1, FAIL=1, `P1 comparator call … count=0, want 1` (**P1d**) |
| **M15** *(D1)* | `ROOT_FLOOR="go$(awk …)"` → `ROOT_FLOOR="go1.26.6"` | hardcoded 0→1, derived 1→0 | rc=1, RUN=1, FAIL=1, `contains hardcoded Go version literal(s) [go1.26.6]` (**P1e**) |
| **M16** *(D1)* | below-floor branch `exit 1 ;;` → `exit 0 ;;` | in-block `exit 0` 2→3 | rc=1, RUN=1, FAIL=1, `has 2 \`exit 1\` statements, want 3` (**P1c**) |
| **M17 — GREEN CONTROL** *(D1)* | reword a **comment** word inside the P1 block | reworded 0→1 | **rc=0, RUN=1, PASS=1** — mandatory: without it the set is unfalsifiable |
| **M18** *(D1)* | move the whole P1 block **below** the race leg | sentinel line > race-leg line | rc=1, RUN=1, FAIL=1, `P1 floor gate is out of order (deny-list@…, P1@…, race leg@…)` (**P1f**) |
| **R1 / R2 — residual probes** *(D1)* | `case "$floor_rc" in`→`case "0" in`; and invert `if(x>y) exit 0; if(x<y) exit 1}` | 1 each | **rc=0, RUN=1, PASS=1 — the static set does NOT catch these.** Record the PASS; their killer is the M9 runtime arm (§9.1) |

**Runtime-lane arms (CONTROLLER):** M9 + its A3 control, M10a′ + control, M10b + control, A1, A4 —
all with the expected readings in the M2 table above.

**D1 runtime-blindness arms (CONTROLLER, and they are the point of the needle):** run each of M13,
M14, M15, M16, M18, R1, R2 as a patched copy under **(i)** the ambient toolchain and **(ii)** M9
conditions. Expected, measured: **ambient → rc=0, marker=1, races=2 for all seven, identical to the
correct block**; under M9 conditions **M14, R1, R2 → rc=0 marker=1 races=2 (fail open)** and
**M16 → rc=0 marker=0 (false green: it exits SUCCESS after printing its own FATAL)**. Record these —
they are what makes the static needle non-optional rather than defensive decoration.

**Hygiene + closure (executor):**

| Command | Expected |
|---|---|
| `gofmt -l host/verifygate/` | **empty** |
| `go vet ./host/verifygate/` | **rc=0** (AC6 — **not** `go build ./host/verifygate/`, which is rc=1 at base) |
| AC9-a: `grep -n 'run: ./scripts/verify_go.sh' .github/workflows/ci.yml` | **`166:`**; control — `sed -n '163p'` is the step name `go build + test gate (replay tests run against pinned AILANG_BIN)` |
| AC9-b (**corrected**, §4.5b): `grep -cE '^go test \./\.\.\. -count=1$' scripts/verify_go.sh` | **1** (the doc's looser pattern matches **3** lines, two of them a comment and an `echo`) |
| AC9-c: `go list ./... \| grep -c verifygate` | **1** |
| AC8: `find . -name go.mod -not -path './.git/*' \| wc -l` | **3** |
| AC10: nonsense selector `-run 'TestNoSuchRaceControlFloorTestZZZ'` | rc=0 **with `[no tests to run]`** — RED if it ever runs a test |
| final: `git status --porcelain` | only the three deliverable paths + `.snap/` |
| final: `shasum -a 256 go.mod design_docs/verification/w-race-gate-blindspot/racecontrol/go.mod` | `7a298361…` and **the post-fence hash from M3**, not `ab782f11…` (the fence is a deliverable; the *root* `go.mod` must be `7a298361…`) |
| **CONTROLLER final:** `go test ./host/verifygate/ -count=1 -v` | rc=0, **RUN=53, PASS=35** (base measured this session: RUN=**52**, PASS=**34**), FAIL=0, SKIP=0 |
| **CONTROLLER final:** `./scripts/verify_go.sh` | green, or rc=1 **with the failing set named** and shown to be outside `{scripts/verify_go.sh, host/verifygate, racecontrol}` |

**Snapshot:** cumulative into `.snap/M4/` → **Commit 4**.

---

## 8. Acceptance criteria — non-vacuity audit

| AC | Binds a counted `=== RUN`/`--- PASS`? | Red at base? | Green after? | Notes |
|---|---|---|---|---|
| AC1 | **yes** (RUN=1, PASS=1) | **yes** — base gives `[no tests to run]`, RUN=**0** | yes | the doc's own repair of the rc-only form |
| AC2 | yes (via M1–M8 + M9/M10 readings) | yes | yes | citation corrected (§4.5c) |
| AC3 | yes (M1: RUN=1, FAIL=1) | n/a (mutation arm) | red-on-mutant | runtime half must use `GOTOOLCHAIN="$ACTIVE_GO"`, never `auto` (§4.5e) |
| AC4 | yes (M2: RUN=1, PASS=1) | n/a | green-on-mutant | equality arms the control (2 races, measured) |
| AC5 | **yes** (RUN=1, PASS=1) | yes | yes | `env -u AILANG_BIN`; no socket, no subprocess |
| AC6 | n/a (vet/gofmt) | **green at base — declared** | green | `go vet`, **not** `go build ./host/verifygate/` |
| AC7 | n/a (text) | **yes** (`LOAD-BEARING`=0 at base) | yes | sha256 clause replaced (§4.4) |
| AC8 | n/a | **green at base — declared** | green | census closure, measures nothing about this sprint |
| AC9 | n/a | **green at base — declared** | green | linkage existence; sub-grep corrected (§4.5b) |
| AC10 | **the rule itself** | — | — | obeyed by every go-test row in §7 |
| AC11 | yes (marker + race counts) | n/a | measured | M10 replaced by M10a′/M10b (§4.2); the `exit 2` branch is a **standalone-battery** killer, not a runtime arm |
| **AC12** *(new, D1)* | **yes** (each mutant RUN=1, FAIL=1; M17 RUN=1, PASS=1) | **yes** — the needle set does not exist at base | yes | the P1 block is pinned by **six** semantic assertions with **seven** named mutants and **one green cosmetic control**; two residual probes (R1/R2) are recorded as PASSing and routed to the runtime lane (§9.1) |

**Declared vacuity residuals:** AC6, AC8 and AC9 are green on the unmodified tree and therefore
measure the sprint's edit only in the negative (they must not *become* red). The doc declares this
for AC6 and AC9; AC8 is added here. **No AC in this plan is satisfiable by a `go test` run that
executed zero tests.**

---

## 9. Declared residuals

### 9.1 The one class a static scan genuinely cannot hold (measured, not asserted)

The P1 needle set (D1, §4.3) catches **structural** gutting: deletion, single-branch removal, operand
inversion, hardcoding, fail-open exits, and reordering — seven named mutants, §3.8. It **cannot** catch
**semantic** neutering that preserves every byte of text:

- **R1** — `case "$floor_rc" in` → `case "0" in`: all three branches present, dispatch dead.
  **Measured: static set PASSES (rc=0, RUN=1, PASS=1).**
- **R2** — invert the comparator's verdicts inside the `awk` program (`if(x>y) exit 1; if(x<y) exit 0`).
  **Measured: static set PASSES.**

Both are caught by the **M9 runtime arm**, whose expected red *vanishes* under them (measured: R1 and
R2 each give `rc=0 marker=1 races=2` under M9 conditions where the correct block gives rc=1/marker=0).
So the sprint catches them. **What nothing catches is R1/R2 landing later, in ongoing CI** — because CI
only ever runs the ambient condition, where a correct P1 and a neutered P1 are behaviourally identical
(§3.8). Closing that would require a CI job that forces a below-root-floor `ACTIVE_GO`, which this rig
cannot supply (every locally available below-floor base is itself deny-listed). **Declared, not
papered over: a needle that pretended to cover R1/R2 would be row 49's defect wearing this sprint's
name.**

### 9.2 The rest

1. **P1's malformed-token (`exit 2`) branch has no live-script killer on this rig** (§4.2). Every root
   `go.mod` mutation that produces an unorderable token also flips `go env GOVERSION` into the
   deny-listed base, so the deny-list fires first; and the `ACTIVE_GO`-side tokens (`devel`,
   `go1.26.6rc1`) cannot be installed here. **Its killer is the standalone battery of §3.5 and nothing
   else** — a weaker instrument than a full-script arm, and stated plainly rather than banked as one.
2. **The static gate cannot know the live `ACTIVE_GO`** — the doc's own declared residual; the chain
   `floor(racecontrol) <= floor(root) <= ACTIVE_GO` is P3 (static) + P1 (runtime) + P2 (binding).
3. **A re-spelled execution binding escapes the P2 needle** (computed variable, indirection). Row 42's
   declared residual, inherited. The same bound applies to P1d's call-site needle.
4. **`verify_go.sh` is flaky on this rig.** A single red is a datum; the failing set must always be
   named.

---

## 10. Time

| Task | Milestone | Est. | Lane |
|---|---|---|---|
| T1 — append the P3 test | M1 | 0.15 h | executor |
| T2 — insert P1 verbatim + P2 one-liner | M2 | 0.15 h | executor |
| T3 — P1 runtime arms + controls (9 script runs) | M2 | 0.30 h | **controller** |
| T4 — fence + AC7 | M3 | 0.05 h | executor |
| T5 — static mutation battery (8 arms, land-and-restore) | M4 | 0.20 h | executor |
| **T5b — D1 needle battery: M12–M18 + R1/R2 (9 arms)** | M4 | 0.15 h | executor |
| **T5c — D1 runtime-blindness arms (7 variants × 2 conditions)** | M4 | 0.20 h | **controller** |
| T6 — hygiene, AC8/AC9/AC10, final gates | M4 | 0.15 h | mixed |
| **Total** | | **≈1.35 h** | the doc's ~0.1 d + the authorised D1 needle |

---

**SPRINT_PLAN_PATH**: `design_docs/planned/w-racecontrol-floor-bump-disarms-the-race-control-sprint-plan.md`
**DESIGN_DOC_PATH**: `design_docs/planned/w-racecontrol-floor-bump-disarms-the-race-control.md`
