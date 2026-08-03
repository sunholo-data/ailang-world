# w-bench-load-confound — make the benchmark baseline's comparability conditions machine-recorded, and make the A/B against the parent commit the only mechanically-valid form of a cost claim

**Status:** **PARKED `needs-human-review` on OD-6** (2026-08-03, iteration 42). Two quorum rounds,
both BLOCKED; round 2 went `gemini-3-1-pro` **pass** / `gpt5-6-sol` **reject**. The remaining
objection disputes whether the design's central claim holds, and its own proposed fix branches
two ways — so the narrow-refinement carve-out does **not** apply and was **not** applied.
**No milestone may be routed to a planner or an executor until OD-6 is answered** (see *Open
decision for the human*). The design below is what will be built *if* OD-6 resolves toward
branch A; BC.A is very nearly unaffected either way.
**Item:** `w-bench-load-confound` (queue item 4f, the only unblocked head item; charter row at
`design_docs/world-mission.md:1568`) + carry-forward **CF-K-2** (toolchain as a condition of
comparability), which folds in here.
**Clause:** clause-1
**Date:** 2026-08-03
**Iteration:** 42
**Author:** `claude:claude-fable-5` (rotation designer); every premise row below was re-measured
first-party by this designer at `c1e6125` unless explicitly marked otherwise.
**Estimate:** 0.35–0.6 day total across two milestones.

**Revision round 2 (2026-08-03).** Both round-1 objections are adopted with their reviewers'
fixes; the design direction is unchanged. (1) `gpt5-6-sol`: the recorder now runs every `go`
invocation under a hardcoded wall-clock deadline with process-group kill — mirrored verbatim
from `verify_ail.sh`'s `run_bounded` (V26, P23) — because `go test -timeout` bounds only the
test binary's own clock, not the compile step or a wedged child (D1 *Bounded execution*, AC1,
`MUT-REC-STALL`, P25). (2) `gemini-3-1-pro`: `go env GOVERSION` is directory-sensitive under
this rig's live `GOTOOLCHAIN=auto` (re-measured first-party, P24), so tree-dependent probes now
run against `--dir` via `go -C`; a caller-cwd probe would write the variant's toolchain into
BOTH sections of a pair and R4 would compare a value to itself — a gate that cannot fail, inside
CF-K-2's remedy, and load-bearing the moment OD-1's floor change lands (D1 step 3, AC6,
`MUT-AB-FLOOR-SPLIT`, `MUT-PROBE-CALLER-DIR`). The D1 execution order is re-derived (rig-global
probes first) and D4's determinism invariant restated against it; the two fixes interact (the
`--dir`'d `go env` can trigger a toolchain fetch, so it too is bounded) and that interaction is
handled in D1. The three round-1 residuals — the ubuntu refusal as an unmeasured expectation,
R4's in-file-only parent linkage, the recorder as itself a competing process — are unchanged and
remain labelled as residuals, not upgraded to claims.

---

## Problem statement

`bench/BASELINE.md` makes a comparability promise it has no mechanism to keep. Line 5 says
*"Later sprints diff against this file on the same development rig"*, and lines 7–8 reason
explicitly about gate honesty (*"Noise-gating a shared runner would be a dishonest gate (S6)"*)
— but the file considered CI noise and not the standing local confound: **the development rig is
shared with the sibling V1 mission**, whose eval suite runs `ollama` + `llama-server` (and
generically-named temp-path `go-build` binaries) at up to 100% CPU on a schedule this loop does
not control and cannot see.

Iteration 39 paid for this. With V1's suite live (`load averages: 5.22 4.99 5.91`),
`BenchmarkBrokerFSRead` p95 read **4.529 ms** against the idle-rig M3.C row of **0.7472 ms** —
a **6.06× apparent regression** that the sprint was about to bank as the effect journal's cost.
The identical invocation on the pre-MJ.C parent `b485ead`, under the same load, minutes apart,
read **4.523 ms**: the true cost was **+0.13% — zero**. The confound was caught only because the
ratio looked implausible enough to warrant a control.

**The gap is provenance and automation, not literally "no load number exists anywhere."**
`BASELINE.md:18` records the Go toolchain and `BASELINE.md:22` records the measurement-time load
— but both were **typed by the controller by hand, after the fact**. Nothing in
`scripts/bench_worldd.sh` (49 lines, `--smoke` only) or in `host/daemon/bench_test.go` produces
any condition mechanically, so the next measurement omits them **by default**, and a hand-typed
condition can be wrong, stale, or absent with nothing going red. CF-K-2 sharpens the stakes:
iteration 40 proved go1.26.0–1.26.5 **miscompile** a shape present at `host/store/scan.go:74`
and `:112` on darwin/arm64, so a toolchain change invalidates every number in the file — the
toolchain is a comparability condition of exactly the same kind as the load, and it too is
recorded only by hand today.

This item makes two things mechanical:

1. **the conditions record** — emitted by the harness at measurement time, never typed; and
2. **the policy** — a cost claim is valid **only** as a same-rig A/B against the parent commit,
   and the file's structure gate goes RED when a measurement lands without its paired control.

It deliberately does **not** re-derive the amortisation numbers (see *Out of scope*): those are
pinned to a toolchain that parked decision **OD-1** exists to change, and banking numbers under
an unresolved condition is precisely the defect this item exists to fix.

## Premise Verification Log

All rows measured by this designer at `c1e6125` (clean tree, verified) on 2026-08-03, on the
development rig (darwin/arm64, zsh, `PATH` prefixed with `/opt/homebrew/bin`), except where
marked **charter-cited**. Empty/negative rows carry a known-positive control in the same call.
Rows **P23–P25** were measured during revision round 2 (2026-08-03, same rig, same `c1e6125`
checkout; tree state: only this document untracked).

| # | Claim | How verified (command) | Observed | Verdict |
|---|---|---|---|---|
| P1 | The harness records NO conditions | `grep -nEi "load\|uptime\|sysctl\|GOVERSION\|toolchain\|nproc\|benchstat\|git rev-parse\|hw\.\|ncpu" scripts/bench_worldd.sh > /tmp/f1out 2>&1; echo "rc=$?"` with same-call control `grep -c -i "benchmark" scripts/bench_worldd.sh` | `rc=1`, zero output; control = **15** — the instrument sees the file, the absence is real | CONFIRMED |
| P2 | The harness is 49 lines and has exactly one mode | `wc -l scripts/bench_worldd.sh` + full first-party read | `49`; usage accepts only `--smoke`; the mode runs `go test -bench . -benchtime 1x -run '^$' ./host/daemon/` and asserts a hardcoded 10-name manifest (lines 15–26, 37–47) | CONFIRMED |
| P3 | CI runs the smoke gate, and it is a NAME gate, not a measurement gate | First-party read of `.github/workflows/ci.yml` (whole file) | Job `go-verify`, step `worldd benchmark smoke gate` at `ci.yml:88-89` runs `./scripts/bench_worldd.sh --smoke`. No numbers recorded, no thresholds evaluated anywhere in the workflow | CONFIRMED — read, not grepped, per the controller's F2 truncation warning |
| P4 | CI checkouts are shallow (no `fetch-depth`) | Same read | Both jobs use bare `actions/checkout@v4` (`ci.yml:13`, `:53`) with no `fetch-depth` key | CONFIRMED — constrains the checker to in-file structure (see D3) |
| P5 | The comparability promise and the S6 honesty line are real | First-party read of `bench/BASELINE.md:5-8` | `:5` "Later sprints diff against this file on the same development rig"; `:7-8` "Noise-gating a shared runner would be a dishonest gate (S6)" | CONFIRMED |
| P6 | The conditions ARE recorded today, but BY HAND | First-party read of `bench/BASELINE.md:16-26` | `:18` `Go: go version go1.26.4 darwin/arm64`; `:22` `Rig load at measurement: load averages: 5.22 4.99 5.91` — prose typed by the controller, produced by nothing | CONFIRMED |
| P7 | The confound is a standing condition, live right now | `uptime; sysctl -n hw.ncpu hw.model; ps -Ao pid=,ppid=,pcpu=,comm= \| sort -k3 -rn \| head -8` at 2026-08-03 07:27 CEST | `load averages: 3.93 3.20 3.03` on 16-core Mac16,9; top consumers: **100.0%** `/var/folders/.../go-build200642988/b001/exe/solution` (pid 28606), **99.6%** `node` (pid 90462), **77.0%** `ollama` | CONFIRMED |
| P8 | The generic temp-path binary CAN be attributed one level up, mechanically | `ps -o comm= -p 28600` (ppid of the 100%-CPU `solution` binary) | `go` — the parent comm resolves; a `go`-spawned test binary is identifiable as such without guessing which mission owns it | CONFIRMED — D1's parent-comm field is feasible |
| P9 | Toolchain state and the floor | `go version; go env GOVERSION GOOS GOARCH; sed -n '1,5p' go.mod` | `go1.26.4 darwin/arm64`; `go.mod:3` = `go 1.26.4`. The `go` directive is a **floor** (iter-40, measured); only `GOTOOLCHAIN` selects exactly, and lowering the floor is **OD-1, parked** | CONFIRMED |
| P10 | `benchstat` is NOT installed | `command -v benchstat; echo "rc=$?"` with same-call control `command -v go` | `rc=1`; control `/opt/homebrew/bin/go` | CONFIRMED — this design uses no benchstat |
| P11 | The recorder's probes exist and their output shapes are known | `sysctl -n vm.loadavg` · `sysctl -n hw.model hw.ncpu` · `date -u +%Y-%m-%dT%H:%M:%SZ` | `{ 2.96 3.20 3.06 }` · `Mac16,9` / `16` · `2026-08-03T05:29:25Z` | CONFIRMED — D1 parses these exact shapes |
| P12 | `python3` is present, and the repo already requires it loudly | `command -v python3` + first-party read of `scripts/verify_ail.sh` | `/opt/homebrew/bin/python3`; `verify_ail.sh` exits 1 with a named message when python3 is absent ("python3 is REQUIRED … fails the script LOUDLY") and states it is present on macOS dev machines and ubuntu-latest | CONFIRMED — the checker may embed python3 by precedent |
| P13 | `bench/BASELINE.md` holds exactly 3 raw benchmark blocks, and naive detection overcounts | python3 fence-parse counting fenced blocks containing a line matching `^Benchmark[A-Za-z_/]+(-[0-9]+)?\s`, vs. the naive `awk` ns/op count | Fence-parse: **3** blocks, opening at lines **222, 245, 264**. Naive ns/op count: **4** — the extra "block" is prose (`78.55 ns/op` at `:109`). The in-fence `^Benchmark` regex is the correct detector | CONFIRMED — with its own false-positive control |
| P14 | Those 3 blocks hold 30 benchmark lines (3 × 10 rows) | `grep -c '^Benchmark' bench/BASELINE.md` | `30` | CONFIRMED |
| P15 | Parent-commit resolution works at record time on this rig | `git rev-parse HEAD^; git cat-file -t b485ead` | `507f3e4…`, rc=0; `commit` | CONFIRMED — the recorder can emit `parent:` from git, first-party |
| P16 | Reference sweep for the Conflict Surface | `grep -rln "bench_worldd" .` and `grep -rln "BASELINE.md" .` (git dir excluded) | `bench_worldd`: `ci.yml`, `bench/BASELINE.md`, 4 implemented docs, charter+log+archive, the script itself. `BASELINE.md` in **code**: only `host/daemon/handlers.go:295`, which is a **comment** | CONFIRMED |
| P17 | `scripts/verify_ail.sh` and its manifest are out of this item's blast radius | Full first-party read of `verify_ail.sh` | It sweeps only `.ail` files under its two ROOTS; this item adds no `.ail` file and touches no world/ module | CONFIRMED |
| P18 | The pinned AILANG binary is present and is the pin | `/tmp/ailang-v0300/ailang --version` | `AILANG v0.30.0`, commit `e37b370` | CONFIRMED — and this item touches **no `.ail` file** (P17), so the pin appears here only as a recorded rig condition |
| P19 | Baseline repo state | `git rev-parse HEAD; git status --porcelain` | `c1e6125c…`; empty status | CONFIRMED |
| P20 | Adjacent parked docs touch the same surfaces | First-party reads: `w-race-gate-blindspot.md:139,172,184,229`; `w-mcp-projection.md:350-351` | The race doc (parked on OD-1/OD-2) flags `BASELINE.md` comparability as *"Item 4f owns the mechanism"* (`:184`) and proposes `ci.yml` Go-version/`-race` changes (`:139`, `:172`); the MCP doc explicitly *"does not rewrite `bench/BASELINE.md`"* | CONFIRMED — see Conflict Surface |
| P21 | Both sha256 CLIs exist locally, but tool divergence is avoidable | `command -v shasum sha256sum` | `/usr/bin/shasum`, `/sbin/sha256sum` | CONFIRMED — design decision: both recorder and checker hash via **python3 `hashlib`**, so darwin/linux CLI divergence cannot enter the gate |
| P22 | OD-1 is parked and constrains scope | **charter-cited** (`world-mission.md`, iter-40/41 status rows; controller F8) | OD-1 = lower the `go.mod` floor 1.26.4 → 1.25.6; awaiting Mark; this loop may not decide it | ACCEPTED AS CHARTER STATE |
| P23 | A bounded-execution precedent exists in-repo, and it is process-group-killing | Full first-party read of `scripts/verify_ail.sh:44-46,54-74` | `run_bounded` (V26): python3 `Popen(..., start_new_session=True)` puts the child in its own process group; on expiry `os.killpg(..., SIGKILL)` kills the WHOLE group, a named `✗ TIMEOUT after Ns: <cmd>` goes to stderr, exit **124**; deadlines are hardcoded constants (`GATE_LEG_TIMEOUT_S=120`, `GATE_TEST_TIMEOUT_S=180`), not env knobs; every binary invocation in both gate legs runs through it | CONFIRMED — D1 **mirrors this helper verbatim** rather than inventing a third form |
| P24 | `go env GOVERSION` is directory-sensitive; the switching mechanism is live; the recorder's `-C` form works; the AC6 fixture needs no network | Temp module with a `go 1.26.5` floor: `(cd "$d" && go env GOVERSION GOTOOLCHAIN)` · same in the repo · same in module-less `/private/tmp` · `go -C "$d" env GOVERSION` · `go -C <repo> test -c ./host/daemon/` · `ls "$(go env GOMODCACHE)/golang.org"` | Temp module → **`go1.26.5`**, `GOTOOLCHAIN=auto`; repo (floor `go 1.26.4`) → **`go1.26.4`**; no-module dir → `go1.26.4`; `go -C` env form → `go1.26.5` rc=0; `go -C … test` form compiles OK; toolchains **go1.25.6 AND go1.26.5 already cached** in GOMODCACHE | CONFIRMED — reproduces the controller's measurement first-party; a caller-cwd probe records the wrong tree's toolchain, and both floors the AC6 fixture uses are cached (no download) |
| P25 | The benchmark deadline is derivable from measurement, and 120 s would be too small | `uptime`; `time ( go test -c -o … ./host/daemon/ )` cold then warm; iteration-weighted sums computed from the recorded `BASELINE.md` rows (python3) | At load 2.87: cache-cold compile **128.85 s wall** (1.64 s user / 4.67 s system, 4% CPU — cache-cold and rig-contended); warm re-compile 0.174 s; recorded idle full-run binary time 3.433 s (`BASELINE.md:276`); loaded iteration-weighted sum **4.217 s** vs idle **1.924 s** = 2.19×; worst observed p95 inflation 6.06× (iter-39) | CONFIRMED — worst measured legitimate end-to-end ≈ 129 s compile + ~25 s run ≈ **155 s**; `REC_BENCH_TIMEOUT_S=600` is ~3.9× that (arithmetic in D1) |

One negative expectation is stated as an expectation, not a fact: **the recorder is expected to
refuse on ubuntu-latest** because `sysctl -n vm.loadavg` / `hw.ncpu` are BSD names absent from
Linux's sysctl namespace. This was **not measured on CI** from this rig. The design does not
depend on the expectation being right: the CI step asserts the refusal, so if ubuntu ever
satisfies the probes the step itself goes RED loudly (see D4) — the failure direction is loud
either way.

## Conflict Surface

| Existing surface or precedent | Collision / reuse question | Resolution |
|---|---|---|
| `scripts/bench_worldd.sh` `--smoke` + its hardcoded 10-name manifest | The only existing mode; CI depends on it (`ci.yml:88-89`) | **Unchanged byte-for-byte in behavior.** New modes (`--record`, `--check-claims`) are added beside it; `usage()` documents all three (S7). The name manifest is not touched |
| `host/daemon/bench_test.go` | Emits `p50_ms`/`p95_ms` via `b.ReportMetric`; could the recorder need to parse or change it? | **Untouched.** The recorder captures `go test` combined output verbatim as an opaque artifact; it parses nothing from it. Only code reference to `BASELINE.md` is a comment (`handlers.go:295`, P16) |
| `.github/workflows/ci.yml` `go-verify` job | Gains two cheap steps (structure gate + off-rig refusal assertion); iter-39 recorded CI-runtime cost as a real, permanent tax | Both steps are seconds (text parsing; a probe refusal that exits before any benchmark runs). The `ailang-verify` job is untouched. Shallow checkout (P4) is **kept** — the checker is designed to need no git history |
| `w-race-gate-blindspot.md` (parked, OD-1/OD-2) | Also proposes `ci.yml` edits (Go pin, `-race` leg) and its AC7 expects `BASELINE.md` to record a toolchain change | No logical conflict: this item builds the **mechanism** its AC7 will use (its `:184` says "Item 4f owns the mechanism"). Merge adjacency only — different `ci.yml` steps, different `BASELINE.md` sections |
| `w-mcp-projection.md` (planned) | References `BASELINE.md` | Explicitly does not rewrite it (P20); if it later appends measurements, it inherits this grammar — intended |
| `bench/BASELINE.md` historical content (3 raw blocks, P13) | The new pairing rule would RED all three pre-4f blocks | Explicit, **enumerated** legacy escape hatch: exactly 3 `legacy-unconditioned` markers, count-pinned in the checker (the hardcoded-manifest precedent from `verify_ail.sh` and the M2.A bench-name gate). A 4th marker is RED |
| `scripts/verify_ail.sh` + required-check manifest | Is it touched? | **NOT touched** — verified by full read (P17); no `.ail` file is added or changed |
| `go.mod` floor / `GOTOOLCHAIN` | Every toolchain question routes here | **Untouched — that is OD-1, parked (P22).** This item *records* `go env GOVERSION`; it selects nothing and asserts no version value |
| `scripts/verify_go.sh` | Already owns the loud `AILANG_BIN` guard | Reused as division of labor: `verify_go.sh` **gates** the pin; the recorder **records** `$AILANG_BIN --version` verbatim and gates only on the variable being set and executable |
| `tools/launchd/*`, `~/.ailang/state/mission-v1*` | Frozen core | Untouched |

### Why this is not a package

This item instruments the **host benchmark harness** (a shell script around `go test -bench`)
and gates the **structure of a Markdown evidence file** in CI. An AILANG package cannot read
`ps`/`sysctl`, invoke `go test`, or fail a GitHub Actions job. No kernel API, world type, policy,
or reusable package behavior is introduced; `world/` is untouched.

## Design

### D1. The conditions record — machine-emitted, refuse-loudly, paste-without-retyping

`scripts/bench_worldd.sh` gains a mode:

```
./scripts/bench_worldd.sh --record <variant|control> [--dir <worktree>]
```

The role argument is **mandatory** — there is no default role. `--dir` names the git worktree to
measure (default: the script's own repo root); it exists because the **instrument must be held
constant across an A/B pair** while the code-under-test varies: the control run at the parent
commit is driven by the *variant* worktree's recorder (the parent commit may predate the
recorder's existence — true for the very first pair).

Execution order is fixed and load-bearing (AC1 and D4 depend on it; re-derived in revision
round 2):

1. **Argument validation** — unknown/missing role exits 2 via `usage()`. The `--dir` value
   (default: the script's own repo root) must exist and satisfy `git -C <dir> rev-parse
   --git-dir`, else a named exit 1. The resolved `--dir` is the probe target in **both** the
   defaulted and the explicit case — one code path, not two.
2. **Rig-global probes**, each of which must produce non-empty output or the mode **exits 1
   naming the exact probe that failed, having emitted NO conditions block, partial or
   otherwise**. Fixed order: `sysctl -n hw.ncpu`; `sysctl -n hw.model`; `sysctl -n vm.loadavg`
   (shape `{ a b c }`, P11); `ps -Ao pid=,ppid=,pcpu=,comm=`; `date -u`; `AILANG_BIN` set +
   executable + `--version` capture — the `--version` execution runs under the bounded runner
   with the 120 s probe deadline, because it executes an env-supplied binary and every such
   wait is bounded (a wedged pin binary is a stall like any other). These read kernel/env state
   and are cwd-independent
   (classification: the sysctl trio, `ps`, and `date` read the rig; `AILANG_BIN` is env-derived,
   not tree-derived). They run FIRST so that an off-rig runner refuses deterministically at
   `sysctl -n hw.ncpu` — before any `go` invocation (which on a floor-raised tree could attempt
   a toolchain download, P24) and before the dirty-tree check (D4's determinism). A recorder
   that emits an empty or partial block when a probe is unavailable would be a **silent
   fallback**; this mission's axioms forbid it, so the null case is a named non-zero exit,
   checked in CI (D4).
3. **Tree-dependent probes, against `--dir`**: `go -C <dir> env GOVERSION GOOS GOARCH` (the
   `-C` form, verified P24), run under the bounded runner below with the 120 s probe deadline.
   `GOVERSION` is **directory-sensitive** on this rig — P24 measured `GOTOOLCHAIN=auto` live
   and a `go.mod` floor selecting the toolchain — and the control run measures a *different
   commit*, which can carry a *different floor*; **OD-1 is exactly such a change**. Probing in
   the caller's cwd (the round-1 design) would write the variant's toolchain into BOTH sections
   of a pair, and R4's `goversion` comparison would compare a value to itself: a gate that
   cannot fail, inside CF-K-2's remedy — the shape iteration 40 shipped and the quorum caught
   there too. `GOOS`/`GOARCH` derive from the selected toolchain and caller env, not from tree
   content, but they ride the same `-C` invocation and are recorded verbatim from it.
4. **Git state** (against `--dir`, via `git -C <dir>`): `git status --porcelain` must be
   **empty** — a measurement of a dirty tree is not attributable to a commit, so the mode
   refuses. `commit` := `git -C <dir> rev-parse HEAD`, `parent` := `git -C <dir> rev-parse
   HEAD^` (P15), both full SHAs.
5. **Measurement**: capture `load_before` + `competing_before`; run
   `go -C <dir> test -bench . -benchtime 200x -run '^$' ./host/daemon/` under the bounded
   runner with the 600 s benchmark deadline (`go -C … test` verified compiling this package,
   P24), capturing combined output verbatim (non-zero rc → exit 1, output shown; deadline
   expiry → the named TIMEOUT refusal below); capture `load_after` + `competing_after`; record
   elapsed seconds. The recorded `invocation:` line is the command **verbatim as executed,
   including its `-C <dir>` token** — the two sections of a pair therefore differ in that token
   by construction, and R4 deliberately compares toolchain-identity fields, never the
   invocation string. Everything in the invocation other than `-C <dir>` is hardcoded in the
   script, so it cannot vary between the sections of a pair without a review-visible script
   edit.

**Bounded execution (Standing Rule 6).** `go test -timeout` bounds only the test binary's own
clock — not the compile step, not a wedged child process, not the `go` tool itself — so **every
invocation of a non-trivial external binary above — both `go` invocations and the
`$AILANG_BIN --version` capture — runs through a `run_bounded` helper mirrored verbatim from
`scripts/verify_ail.sh:61-74`** (the V26 helper, read first-party, P23: python3
`start_new_session` puts the child in its own process group; on expiry the WHOLE group gets
SIGKILL — killing only the direct child would orphan the compiled test binary — a named
`TIMEOUT after Ns` goes to stderr, and the exit is 124). **Mirrored, not sourced or extracted**:
sourcing `verify_ail.sh` would execute the gate on load, and extracting a shared lib would edit
`verify_ail.sh` — which P17 and AC5 pin as untouched — to move 13 lines. The copy carries a
provenance comment naming `verify_ail.sh:61-74`/V26, and divergence between mirror and original
is a review-visible change (Design Freeze). Two hardcoded deadlines, **policy, not env knobs**
(the D3 precedent):

- **`REC_PROBE_TIMEOUT_S=120`** on the probe-stage binaries (`go env` and the
  `$AILANG_BIN --version` capture) — same species and value as
  `GATE_LEG_TIMEOUT_S` (P23). It exists because the two fixes interact: moving the probe into
  `--dir` means `GOTOOLCHAIN=auto` may fetch an uncached toolchain over the network when first
  probing a floor-raised tree (P24) — a wait that must itself be bounded. (The AC6 fixture
  avoids the fetch entirely: both floors it uses are already cached, P24.)
- **`REC_BENCH_TIMEOUT_S=600`** on the benchmark run — chosen from measurement (P25), recorded
  here: the worst *measured legitimate* end-to-end on this rig is ≈ 155 s (a cache-cold,
  rig-contended compile took **128.85 s wall** at load 2.9 — so a 120 s ceiling would clip a
  legitimate run — plus a loaded run bounded by ~25 s: the 4.2 s loaded iteration-weighted sum
  with iter-39's worst observed 6.06× inflation applied). 600 s is ~3.9× that and ~2× a
  doubled-load extrapolation, while still bounding a wedge inside ten minutes. Raising it is a
  review-visible constant edit, never a knob.

A deadline expiry is a **named non-zero refusal that emits ZERO fences** — the same class as a
probe refusal (`bench run TIMEOUT (>600s) — no conditions emitted`). A stall can never produce
a partial or `unknown`-filled conditions block. Proven by `MUT-REC-STALL` (AC1).

**Competing-process capture** (the P7/P8 hard case): from one `ps -Ao pid=,ppid=,pcpu=,comm=`
snapshot, keep rows with `pcpu ≥ 25.0`, sorted descending, capped at 8, one line each:
`pcpu=<x> pid=<p> comm=<verbatim path> parent=<comm of ppid>`. The generic temp-path `go-build`
binary is handled by **recording, not guessing**: the verbatim path plus the parent's comm (P8
measured it resolving to `go`) is what is mechanically observable; naming the owning *mission*
is not, and the recorder does not pretend otherwise. If the parent has already exited,
`parent=?` is written — an explicit recorded unknown, distinguishable from an absent field. If
no process meets the threshold, the line is the explicit `competing_before: none>=25%` — an
explicit negative, never an empty field.

**Output**: one paste-ready section printed to stdout AND written to a `mktemp` file whose path
is printed — the human appends it with `cat "$file" >> bench/BASELINE.md` (or a labelled section
paste), retyping nothing. The section is:

````
```bench-conditions
schema: bench-conditions/1
role: variant
utc: 2026-08-03T05:29:25Z
commit: <full sha>
parent: <full sha>
tree: clean
goversion: go1.26.4
goos_goarch: darwin/arm64
ncpu: 16
hw_model: Mac16,9
ailang_pin: AILANG v0.30.0 (e37b370) via $AILANG_BIN
load_before: 2.96 3.20 3.06
competing_before: pcpu=100.0 pid=28606 comm=/var/folders/.../exe/solution parent=go
competing_before: pcpu=99.6 pid=90462 comm=node parent=<...>
invocation: go -C /private/tmp/world-wt-4f test -bench . -benchtime 200x -run '^$' ./host/daemon/
elapsed_s: <n>
load_after: <a b c>
competing_after: none>=25%
output_sha256: <hex>
conditions_sha256: <hex>
```
```text
<the go test combined output, verbatim>
```
````

**Integrity fields**, both computed via python3 `hashlib` (P12, P21 — no shasum/sha256sum CLI
divergence): `output_sha256` = SHA-256 of the raw benchmark output (the lines between the
` ```text ` fences, joined with `\n`, trailing `\n`); `conditions_sha256` = SHA-256 of every
block line from `schema:` through `output_sha256:` inclusive, each terminated with `\n`. These
are **tamper-evidence, not cryptographic provenance** (see *What these gates CANNOT fail*):
their job is to make retyping, truncation, and post-hoc number-editing RED, and — applying the
iteration-41 lesson that *a gate whose failure message hands you the silencing value is not a
gate* — **no failure message anywhere in this design prints an expected digest**. Re-greening
requires re-running the recorder, not pasting a value out of the error text.

### D2. The policy — (ii)-as-A/B, written into the file and cashed out mechanically

The charter's stated preference is adopted and its reasoning recorded in `BASELINE.md`'s header:
**an A/B against the parent commit is MANDATORY for any cost claim.** The rejected alternative —
refuse-to-record above a load threshold — is out of scope with reasons (see *Out of scope*):
`BASELINE.md:7` already calls noise-gating a shared runner dishonest, and a threshold gate
assumes an idle rig that a two-mission rig does not guarantee, whereas the A/B is correct under
any load (it is what actually caught the 6.06× artefact).

"Mandatory" cashed out:

- **The artifact that carries a claim** is a **pair** of recorder emissions in
  `bench/BASELINE.md`: one `role: variant` section at the measured commit and one
  `role: control` section whose `commit` equals the variant's `parent`. New numbers cannot
  enter the file's evidence in any other form: **every** non-legacy raw benchmark block must be
  bound to a conditions block, and **every** variant must have its control. Appending a
  measurement without its parent-commit control is structurally RED — which is the A/B mandate
  expressed as file grammar rather than as prose exhortation.
- **What checks it**: `./scripts/bench_worldd.sh --check-claims` (D3).
- **What goes RED, where**: the `worldd bench claim-structure gate` step in CI's `go-verify` job
  (every push to `dev` and every PR), and the same command locally. F2/P3 establish that today's
  CI evaluates no bench structure at all; this adds the missing leg without touching the smoke
  gate.
- **The A/B procedure** (added to `BASELINE.md`'s header, executed-verbatim per S7 — by the
  controller, since the executor sandbox denies the loopback binds five benchmarks need, the
  standing `<CONTROLLER-MEASURED>` precedent): record the variant at HEAD; `git worktree add
  /tmp/bench-control HEAD^`; record the control via
  `./scripts/bench_worldd.sh --record control --dir /tmp/bench-control`; append both sections;
  run `--check-claims`; commit.

Prose in the header additionally states the policy's meaning for readers: a delta between
invocations recorded under different conditions is **indicative only**; the pair is the only
load-independent signal. (Prose cannot be machine-policed; the grammar above polices the
evidence. The residual is stated honestly below.)

### D3. The structure checker — in-file, history-free, null-case-loud

`./scripts/bench_worldd.sh --check-claims` parses **`bench/BASELINE.md` only** (hardcoded path,
no path argument, no env knob — the `verify_ail.sh` "policy is not env-overridable" precedent).
Implementation: bash + embedded python3 (P12 precedent), python3 absence fails loudly. Rules,
each with a named failure:

- **R1 — block validity.** Every ` ```bench-conditions ` block must contain every required key
  from D1's schema, non-empty, `tree: clean`, and a `conditions_sha256` that recomputes. The
  mismatch message names the block's line number and **does not print the expected digest**.
- **R2 — no orphan numbers.** Every fenced block containing a line matching
  `^Benchmark[A-Za-z_/]+(-[0-9]+)?\s` (the P13-verified detector, immune to the `ns/op`-in-prose
  false positive) must be immediately preceded — within the 5 lines above its opening fence —
  by either the close of a valid conditions block or a `<!-- legacy-unconditioned: pre-4f … -->`
  marker.
- **R3 — output binding.** For each conditions block, the adjacent raw block's content must hash
  to `output_sha256`. Editing one digit of a recorded number is RED; the message does not print
  the expected hash.
- **R4 — the A/B mandate.** Every `role: variant` block must have a `role: control` block in the
  file with `control.commit == variant.parent`, and the pair's `goversion`, `goos_goarch`,
  `ncpu`, and `hw_model` must be **identical** — the CF-K-2 tooth: a cross-toolchain "A/B" is
  not a comparison and REDs by name (`toolchain mismatch inside claimed A/B pair`). Parent
  linkage is checked **in-file** (recorder-emitted `parent:` vs control `commit:`), not via git
  history — deliberate: CI checkouts are shallow (P4) and squash-merged branch SHAs may become
  unreachable later; the git edge is true at record time by construction (`git rev-parse HEAD^`,
  P15), and what CI re-verifies forever is the recorded structure. This boundary is declared in
  *What these gates CANNOT fail*.
- **R5 — enumerated legacy.** Exactly **3** legacy markers exist (the P13 blocks at lines
  222/245/264 get them in BC.B). `!= 3` is RED in both directions: a 4th marker (the "just mark
  it legacy" evasion) and a deleted marker both fail. Widening the pin is a review-visible
  checker edit, never a paste.
- **R6 — the checker's own null case (S6).** Zero raw benchmark blocks in the file, or an
  unreadable/empty file, is RED (`bench/BASELINE.md contains no benchmark evidence — refusing a
  vacuous pass`). A checker that greens on an empty file is the null-case defect this repo has
  now found seven times.

### D4. CI wiring — two steps, both cheap, one asserting the refusal path

Appended to `go-verify` after the smoke gate (`ailang-verify` untouched):

1. **`worldd bench claim-structure gate`**: `./scripts/bench_worldd.sh --check-claims`. Pure
   text parsing; sub-second.
2. **`worldd bench recorder refuses off-rig`**: runs `--record variant` and **asserts it exits
   non-zero AND stderr names a failed probe** (explicit `if`-negation, not `!`-with-`set -e`
   ambiguity; greps the named `probe FAILED` marker). On ubuntu the BSD sysctl probes are
   expected absent, so this exercises the D1 refusal path — the anti-silent-fallback property —
   as a standing CI gate, before any benchmark would run (seconds, not minutes). If a future
   runner ever satisfies every probe, the step REDs loudly and a human looks; the expectation
   failing is itself loud, never silent (see the note under the Premise Verification Log).

Determinism re-derived for the revised D1 order: **rig-global probes run before tree-dependent
probes, git state, and the benchmark**, so on ubuntu the refusal fires at the FIRST probe in the
fixed order — `sysctl -n hw.ncpu` — before any `go` invocation (no toolchain fetch can start on
CI) and before the dirty-tree check (earlier steps' downloaded tarballs cannot affect which
error fires). The step's assertion is unchanged: exit non-zero AND a named failed probe on
stderr. The `--dir` change does not weaken this: `--dir` defaults to the checkout root, which
exists and is a git dir on CI, so step 1 passes and the refusal still lands in step 2,
deterministically.

### D5. `BASELINE.md` content changes (BC.B)

- Header gains: the policy paragraph and A/B procedure (D2), and the statement that the
  machine-emitted conditions block is the only valid form for new measurement conditions.
- The three historical raw blocks (P13) each gain their `legacy-unconditioned` marker line.
- The amortisation section's existing note gains one sentence naming the **post-OD-1 follow-up**
  (see *Out of scope / Deferred*) instead of the current open-ended "when queue item 4f lands a
  load gate" — this item deliberately lands **no load gate**, and the pointer must not imply one.
- One real recorder **pair** (variant at the BC.B merge candidate, control at its parent),
  produced by the controller, is appended under a section explicitly labelled: *mechanism
  acceptance run — conditioned on go1.26.4, superseded by the first post-OD-1 measurement; NOT a
  milestone performance baseline*. This gives deliverable (i) a real in-file instance, gives
  R1–R4 a live pair to validate forever, and gives the named mutations their targets — while the
  label keeps it from being read as banked performance numbers under the unresolved OD-1
  condition.

## Acceptance Criteria

Every AC names a concrete mutation that makes it RED and states what makes that RED hard to
silence by pasting a value back.

### AC1 — the recorder refuses loudly; it never emits a partial block

Owned by **BC.A**. With any single probe unavailable, `--record <role>` exits non-zero, names
the probe, and emits **zero** `bench-conditions` fences (checked by grepping the captured
output). With a dirty `--dir` tree it refuses naming the git state. **With the benchmark
invocation stalled past the deadline, it exits non-zero via the named TIMEOUT refusal, emits
zero fences, and leaves no orphaned child process running** (the process-group kill, D1). Probe
checks precede git checks precede the bounded benchmark (D1 order).

Named RED mutation: **`MUT-REC-SILENT-DEFAULT`** edits `scripts/bench_worldd.sh` (**HARNESS**)
to make a failed probe yield the literal `unknown` instead of refusing. Detection: the AC1
verification run (a PATH-shadowed `sysctl` stub returning empty, rc=0) now **emits a block
instead of refusing** — RED at BC.A review; and permanently RED in CI once BC.B's off-rig step
lands, because `--record` on ubuntu would start succeeding. Hard to silence: there is no value
to paste — the gate demands live probe output, and the CI step's assertion is on the *refusal
behavior*, not on any literal.

Named RED mutation: **`MUT-REC-STALL`** edits `scripts/bench_worldd.sh` (**HARNESS**) in ONE
edit: the benchmark invocation is replaced by a stub that writes its own child's PID to a temp
file, spawns that child sleeping 3600 s, and then sleeps itself — a wedge with a grandchild,
the exact Standing-Rule-6 stall shape — and the same edit lowers `REC_BENCH_TIMEOUT_S` from
600 to 5. (Waiting out the real 600 s deadline proves nothing extra about the kill path and
would breach the mutation ceiling; the constant is deliberately not an env knob, so the lowered
value is part of the sha256-restored mutation itself, not a runtime override. That the shipped
constant equals 600 is asserted by reading the landed script at review.) Detection: the run
exits non-zero with the named `TIMEOUT` on stderr, the captured output contains **zero**
`bench-conditions` fences, and after exit `kill -0 <recorded child pid>` fails — the group kill
reached the *grandchild*, not just the direct child. Hard to silence: the assertion is on
refusal behavior and on process death; there is no value to paste.

### AC2 — the emission is complete, and its integrity fields recompute independently

Owned by **BC.A**. A real `--record variant` emission on the dev rig (controller-performed —
sandbox loopback denial, the standing precedent) contains every D1 key, non-empty, with
`competing_*` lines carrying verbatim paths and `parent=` attribution for temp-path binaries
(P8), and both hashes recompute under an **independent** python3 reimplementation (not the
script's own function).

Named RED mutation: **`MUT-EDIT-RAW-NUMBER`** edits one digit of one `p95_ms` value in the
emitted raw block (**EVIDENCE**). The independent recompute mismatches `output_sha256` — and
after BC.B, `--check-claims` R3 is RED. Hard to silence: the failure message prints no expected
hash (the iter-41 lesson, applied); the legitimate green path is re-running the recorder.

### AC3 — the A/B mandate has structural teeth

Owned by **BC.B**. `--check-claims` is green on the post-BC.B file (3 legacy markers + 1 real
pair) and RED under each of:

- **`MUT-CLAIM-NONPARENT`** (**EVIDENCE**): edit the acceptance variant block's `parent:` to the
  grandparent SHA **and recompute `conditions_sha256`** — isolating R4 from R1, proving the
  parent-linkage check itself discriminates. Expected message names the pair, not any SHA to
  paste; the green paths are recording a real control at the true parent or a review-visible
  checker edit.
- **`MUT-CLAIM-TOOLCHAIN-SPLIT`** (**EVIDENCE**): edit the control block's `goversion` to
  `go1.25.6`, recompute its checksum → R4 REDs with the toolchain-mismatch message. This is
  CF-K-2's tooth: numbers from two toolchains can never satisfy the claim grammar.
- **`MUT-CLAIM-ORPHAN-NUMBERS`** (**EVIDENCE**): append a new ` ```text ` block containing one
  well-formed `Benchmark…` line with no conditions block → R2 REDs. This is the "developer
  writes the words 'A/B done'" evasion: words attach no conditions block, so any *numbers*
  offered as evidence go RED; words without numbers are review's province and are declared as
  such below.

### AC4 — the legacy escape hatch is enumerated, not open

Owned by **BC.B**. Exactly 3 `legacy-unconditioned` markers pass.

Named RED mutation: **`MUT-LEGACY-FOURTH`** (**EVIDENCE**): add a 4th marker above the
AC3 orphan block → R5 REDs (`expected exactly 3 legacy markers, found 4`). Hard to silence: the
only green paths are deleting the marker (restoring AC3's RED) or editing the checker's
hardcoded count — a review-visible policy change, per the `verify_ail.sh` precedent.

### AC5 — null cases are loud, and the blast radius is exactly three files

Owned by **BC.B**. Named RED mutation: **`MUT-CHECK-NULL`** (**EVIDENCE**): truncate
`bench/BASELINE.md` to empty (backup taken in the same command, byte-identical restore proven by
sha256 both sides — the iter-38/39 mutation discipline) → R6 REDs rather than passing vacuously.

Scope assertions, verified at review: the implementation changes **only**
`scripts/bench_worldd.sh`, `bench/BASELINE.md`, and `.github/workflows/ci.yml`. Untouched:
`host/daemon/bench_test.go` and its 10-name smoke manifest, `scripts/verify_ail.sh` and its
required-check manifest (P17), `scripts/verify_go.sh`, `go.mod` (OD-1, P22), all `.ail` files,
`world/`, `host/` production code, `tools/launchd/*`. `--smoke` output and semantics are
byte-compatible (CI's existing step must pass unmodified).

### AC6 — the toolchain probe measures the measured tree, and a REAL cross-toolchain pair REDs

Owned by **BC.B** (the probe behavior it exercises lands in BC.A). This is objection 2's
reviewer-required probe, and it runs against the exact condition OD-1 will create. The fixture
is a **throwaway worktree** — the iter-40/41 measurement pattern (charter-cited, P22 context):
the real `go.mod` is never touched, the floor edits exist only as never-pushed commits under
`/tmp`, and the worktree is removed afterwards.

Fixture: `git worktree add /tmp/bench-floorsplit HEAD` (detached), then two throwaway commits —
**A** edits the floor to `go 1.26.5`, **B** on top of A reverts it to `go 1.26.4`. Both
toolchains are already cached on this rig (P24), so no network is involved. Record `control` at
A and `variant` at B via `--record … --dir /tmp/bench-floorsplit` (two bounded recorder runs);
append the pair to `bench/BASELINE.md` under the AC5 backup/sha256-restore discipline.

- **`MUT-AB-FLOOR-SPLIT`** (**EVIDENCE**, the known-positive): B's parent IS A, so R4's parent
  linkage is genuinely satisfied — and `--check-claims` must RED with exactly
  `toolchain mismatch inside claimed A/B pair` (`go1.26.4` vs `go1.26.5`), **not** the
  nonparent message. Unlike the hand-edited `MUT-CLAIM-TOOLCHAIN-SPLIT` (which isolates R4's
  comparison logic), this proves on **real emissions** that the `--dir` probe records each
  tree's own toolchain end-to-end and that R4's tooth bites on a genuine floor-straddling pair.
- **`MUT-PROBE-CALLER-DIR`** (**HARNESS**): edit the recorder to drop `-C <dir>` from the
  `go env` probe — restoring the caller-cwd form the round-1 design shipped — and re-record
  `control` at A (one bounded run). Its section now records the CALLER's `go1.26.4` despite the
  tree's 1.26.5 floor, and `--check-claims` **greens a genuinely cross-toolchain pair** — RED
  at review, because the known-cross-toolchain fixture must RED. This is the vacuity probe for
  the probe placement itself: it reproduces the round-1 bug and shows the fix, not luck, is
  what makes this AC's known-positive fire. sha256-restored.

Mutation ceilings (the DG.A precedent): each named mutation is one run plus one restored
baseline run, 120-second ceiling per invocation — **except the three recorder-invoking runs in
AC6, which are instead bounded by the recorder's own hardcoded deadlines (D1): the bound holds
by construction and a stall cannot exceed 600 s per invocation** — no retry loops; every
mutation states its edited file and its HARNESS/EVIDENCE classification above — none is
presented as proof of a kernel property.

## Milestones

### BC.A — the conditions recorder (~0.2 day)

Modify only `scripts/bench_worldd.sh`:

- add `--record <variant|control> [--dir]` per D1 (rig-global-then-tree-dependent probe order;
  `go -C`/`git -C` against the resolved `--dir` for every tree-dependent invocation; the
  mirrored `run_bounded` helper with its provenance comment and the two hardcoded deadlines;
  loud refusals including the named TIMEOUT; competing-process capture with parent attribution;
  verbatim output capture; python3-hashlib integrity fields; paste-ready emission to stdout +
  mktemp file);
- update `usage()` for all modes (S7);
- leave `--smoke` byte-compatible.

Owns **AC1** and **AC2**. Independently CI-green: no CI change lands here; the existing smoke
step passes unmodified. Acceptance evidence: the shadowed-probe refusal transcript, the dirty-
tree refusal transcript, the `MUT-REC-STALL` timeout transcript (named TIMEOUT, zero fences,
dead grandchild), and one controller-performed real emission with the independent hash
recompute.

### BC.B — the claim-structure gate, the policy, and the wired file (~0.3 day)

- add `--check-claims` per D3 (R1–R6, hardcoded path, no knobs);
- `bench/BASELINE.md`: policy header + A/B procedure, 3 legacy markers, the amortisation
  pointer correction, and the labelled controller-recorded acceptance pair (D5);
- `.github/workflows/ci.yml`: the two `go-verify` steps per D4.

Owns **AC3**, **AC4**, **AC5**, **AC6**. Independently CI-green: lands with the file already
satisfying the grammar it enforces. Executes and records all seven BC.B named mutations
(AC3's three, AC4's one, AC5's one, AC6's two — the AC6 pair via the throwaway floor-split
worktree, removed afterwards) plus re-runs of AC1/AC2's checks against the landed harness.

## Out of scope, with reasons

- **The threshold / refuse-on-load option from (ii)** — REJECTED, not deferred.
  `BASELINE.md:7-8` already calls noise-gating a shared runner a dishonest gate (S6), and a
  threshold gate assumes an idle rig will eventually happen — on a rig shared with a mission
  whose schedule this loop cannot see (P7), that is not guaranteed, so the gate would either
  waste iterations waiting or get its threshold quietly raised until vacuous. The A/B is correct
  under any load and is the mechanism that actually caught the 6.06× artefact. The conditions
  block makes load **visible**; nothing in this design gates on its value.
- **(iii) — re-deriving the amortisation section** — DEFERRED, blocked on **OD-1**. The section
  is pinned to M3.C idle-rig numbers and labelled so. Re-deriving it now would bank a fresh set
  of numbers conditioned on go1.26.4 — a toolchain that OD-1 exists to change, on a compiler
  release line proven to miscompile landed durability code (iter-40). **Banking numbers under an
  unresolved condition is precisely the defect this item exists to fix**, so (iii) becomes the
  named follow-up *"post-OD-1 amortisation re-derivation"* — first clean-rig invocation on the
  ratified toolchain, recorded through this item's recorder, superseding the acceptance pair —
  and not a milestone here. The D5 pointer edit writes exactly this into `BASELINE.md`.
- **Any change to the `go.mod` floor or toolchain selection** — that **is** OD-1, parked for
  Mark (P22). This design records `GOVERSION`; it never selects or asserts one.
- **Benchstat, statistical machinery, more samples** — P10: not installed; nothing here needs
  it; adding a pinned dependency is a separate decision with its own doc if a future item wants
  distributional comparisons.
- **Changing `host/daemon/bench_test.go`** — the instrument's measurement side is not this
  item's defect; the recording side is.

## Design Freeze

The executor may not quietly change these invariants:

- `--record` emits **all-or-nothing**: any probe failure, dirty tree, or failed benchmark run
  produces a named non-zero exit and zero emitted fences. No `unknown`, no partial block, no
  env-var escape hatch.
- Execution order is frozen: rig-global probes (first: `sysctl -n hw.ncpu`) → tree-dependent
  probes (`go -C <dir> env …`) → git state (`git -C <dir> …`) → bounded benchmark. D4's CI
  determinism depends on the rig-global-first half; R4's cross-toolchain tooth depends on the
  tree-dependent-against-`--dir` half. **Moving any tree-dependent probe to the caller's cwd is
  the CF-K-2 bypass** (`MUT-PROBE-CALLER-DIR` demonstrates it) and is a policy change requiring
  review.
- Every invocation of a non-trivial external binary (`go env`, `go test`, `$AILANG_BIN
  --version`) runs through the `run_bounded` mirror (provenance comment: `verify_ail.sh:61-74`,
  V26, P23) with hardcoded deadlines `REC_PROBE_TIMEOUT_S=120` (probes) /
  `REC_BENCH_TIMEOUT_S=600` (benchmark); neither is env-overridable; expiry SIGKILLs the whole
  process group and is a named non-zero refusal emitting **zero** fences; behavioral divergence
  between the mirror and its `verify_ail.sh` original is a review-visible change.
- The role argument is mandatory; there is no default role.
- `competing_*` capture records verbatim paths and parent comm; it never maps a process to a
  mission name.
- Both integrity hashes are computed with python3 `hashlib` in both the recorder and the
  checker; no shasum/sha256sum CLI dependency enters the gate.
- **No failure message in recorder or checker prints an expected digest or an expected SHA** —
  the iter-41 re-green lesson is a frozen property, not a style choice.
- `--check-claims` reads the hardcoded `bench/BASELINE.md` path; no path argument, no env knob;
  the legacy-marker count (3) and the raw-block detector regex are hardcoded policy.
- Raw-block detection is the in-fence `^Benchmark…` regex (P13), never an `ns/op` substring
  count.
- R4's toolchain-identity comparison covers `goversion`, `goos_goarch`, `ncpu`, `hw_model` —
  removing any field from the comparison is a policy change requiring review.
- `--smoke` and its 10-name manifest are byte-compatible; `bench_test.go`, `verify_ail.sh`,
  `verify_go.sh`, `go.mod`, and all `.ail` files are untouched.
- The acceptance pair lands with its "NOT a milestone performance baseline / superseded
  post-OD-1" label verbatim; the executor does not present its numbers as a baseline.
- Every named mutation records its edited file, classification, single-run-plus-restore
  discipline, and sha256-verified restoration.

## What these gates CANNOT fail

These limits are part of the design, not residual fine print:

- **They cannot prove a number was not fabricated.** The hashes are tamper-evidence binding the
  recorded conditions, the recorded output, and the file together; anyone who can run python3
  can forge all of them coherently. The protection against invented values remains what it has
  been for six milestones: the provenance-honesty culture (`<CONTROLLER-MEASURED>`, refusal to
  invent) plus review. This gate makes *silent drift and careless retyping* red, not fraud.
- **They cannot re-verify the git parent edge at CI time.** R4 checks recorder-emitted structure;
  the edge is git-true at record time by construction, but a squash-merged branch SHA may be
  unreachable later and CI checkouts are shallow (P4). A coherent forged pair with a fictitious
  parent SHA passes R4. Stated, accepted, bounded by review.
- **They cannot police prose.** A cost claim written only in words, citing no numbers, attaches
  no evidence and trips nothing here. What the grammar guarantees is narrower and real: no
  *benchmark numbers* enter the file's evidence without machine-recorded conditions and a
  parent-commit control.
- **They cannot make absolute numbers comparable across loads** — nothing here gates on load.
  The conditions block makes incomparability visible; the A/B is the only load-independent
  signal; both facts are written into the file header rather than solved.
- **`ps` %CPU is a decaying average sampled at two instants**, not a profile of the run; a
  competitor that starts and stops between the two captures is invisible. `load_before`/
  `load_after` bound, but do not narrate, the interval.
- **They cannot see the V1 mission's schedule.** That unobservability is the problem statement;
  the recorder observes the rig, it does not coordinate with the sibling loop.
- **The recorder is this-rig-specific by design** (BSD sysctl names, darwin `ps` semantics). On
  any other platform it refuses loudly — that refusal is a feature and is what CI asserts; a
  future Linux measurement rig needs a probe extension, reviewed.
- **They cannot fix resolution or sample-count limits** — the `BenchmarkBrokerDecide`
  one-tick bound and the "ratio near 1.0 runs three times" rule stand unchanged; this item
  records conditions, it does not improve the clock.
- The CI off-rig step asserts refusal on today's ubuntu-latest; it proves the refusal path
  executes, not that every conceivable probe failure is caught on every future runner.
- **The 600 s deadline cannot distinguish a wedge from a run legitimately slower than 600 s.**
  Such a run is refused loudly — the safe direction — and the ceiling is a hardcoded constant
  chosen from measurement (P25: ~3.9× the worst measured legitimate end-to-end), raised only by
  a review-visible edit, never by a knob. A rig that legitimately exceeds it has changed enough
  that the constant *should* be re-derived in review.
- **The 120 s probe bound can refuse a legitimate first-ever toolchain download** on a slow
  network. The refusal names the probe either way; the remedy is pre-caching the toolchain,
  which this rig already does for every floor on either side of OD-1 (P24).

## Open decision for the human

> Numbering note (controller, iteration 42): this is **OD-6**, continuing the mission-wide
> sequence. OD-1/OD-2 belong to `w-race-gate-blindspot` (item 4e); OD-3/OD-4/OD-5 belong to
> `w-ddl-gate-teeth` (item 4d). This section is controller-authored bookkeeping around the
> designer's document; it changes no design text.

### OD-6 — **The pair records the load and then never reads it. Does this item grow to make the A/B contemporaneous, or does it stop claiming the A/B validates a cost claim?**

**The objection, `gpt5-6-sol`, quorum round 2 (verdict: reject), verbatim in the log entry.**
R4 accepts any control block whose `commit` equals `variant.parent`. It compares `goversion`,
`goos_goarch`, `ncpu`, `hw_model` — and nothing else. So an idle variant may be paired with a
control recorded days earlier under heavy load, or vice versa, and the file grammar blesses it
as a mechanically valid cost claim. The reviewer's catch is the sharp part: *the iteration-39
anecdote shows one near-contemporaneous pair happened to cancel load; it does not establish that
arbitrary independently appended runs do.* The document's central claim — that this A/B form is
*"correct under any load"* — is therefore unsupported as written.

**Controller verification (first-party, not a restatement of the reviewer).** Read against the
document's own text at the revised revision:

- R4's full constraint set is `control.commit == variant.parent` plus identical `goversion`,
  `goos_goarch`, `ncpu`, `hw_model`. Nothing more.
- A search for any temporal, load-comparability, pair-identity or control-reuse constraint across
  the whole document returns **nothing** — the only hits are the words "reuse"/"collision" in the
  Conflict Surface table's own column headers.
- **Known-positive control, same call:** the conditions schema *does* carry `utc`,
  `load_before` and `load_after` (schema lines in D1's example emission).

So the defect is not that the confound is unrecorded. **The recorder captures the load; the gate
that blesses the pair never reads it.** The conditions block makes incomparability visible to a
human, while the machine check that decides "this is a valid cost claim" is blind to precisely
the variable this item exists to control. That is this mission's signature shape — a gate that
cannot fail on the property that matters — and this is its **third** appearance inside this
document's own evolution (round 1's caller-cwd toolchain probe, kept as `MUT-PROBE-CALLER-DIR`;
now this). The reviewer is right on the facts.

**Why this parks rather than being fixed headlessly.** The reviewer's `proposed_fix` offers two
branches and they deliver different items:

- **Branch A — make the pair contemporaneous by construction.** Replace the two independent
  `--record` calls with one bounded `--record-pair --variant <dir> --control <dir>` that
  prebuilds both binaries, stamps a unique pair ID, and executes an interleaved fixed order
  (control/variant/variant/control) inside a single capture session; R4 then matches the pair ID
  and rejects control reuse. This keeps the item's promise. It also replaces the design's central
  interface, roughly doubles the measurement work per claim, and grows the estimate beyond the
  0.25–0.5 day the charter row sized — a scope decision.
- **Branch B — keep the mechanism and weaken the claim.** State that the pair is *mechanically
  complete evidence*, not a *mechanically valid cost claim*, and require human review before any
  cost assertion. Smaller, arguably the more honest reading of what the grammar can guarantee —
  but it changes what item 4f delivers relative to its charter row, which promised the A/B
  *"becomes MANDATORY for any cost claim"*.

The reviewer's third limb — *"define a measured acceptance rule for excessive within-pair load
divergence"* — the loop should **not** attempt in either branch: no data exists on this rig to
derive a defensible threshold, and inventing one would be the noise-gate this item already
rejects (`BASELINE.md:7-8`, S6), wearing a new name.

Choosing between A and B is a judgment about scope and about what the item promises, not a defect
with a single reviewer-authored resolution — so the narrow-refinement carve-out does not apply,
and Standing Rule 2 forbids proceeding over a contested central claim.

**Controller recommendation: branch A, bounded.** Take the pair ID, the single-session interleaved
capture, the per-leg timestamps and load snapshots, and R4's control-reuse rejection; take
**neither** the load-divergence acceptance rule nor any threshold. Then state honestly, in *What
these gates CANNOT fail*, that a contemporaneous interleaved pair bounds but does not eliminate
load skew. Rationale: branch A keeps the charter row's promise, and the interleaving is the part
that actually earns the load-independence claim — iteration 39's pair worked *because* its legs
were minutes apart under the same load, which is the property branch A makes structural instead
of accidental. Estimated at ~0.6–0.9 day, i.e. the item outgrows its original sizing; that is the
cost of the promise, and it belongs to the human, not to me.

**BC.A is very nearly branch-independent.** The recorder's probes, bounded execution, refusal
semantics and integrity fields are needed under both branches; only its invocation surface
changes (`--record` vs `--record-pair`). If the queue needs forward motion before OD-6 is
answered, BC.A is the routable half — but that too is the human's call, since branch A would
rework the interface BC.A ships.

## Deferred

- **Post-OD-1 amortisation re-derivation** (item (iii)) — first clean-toolchain invocation,
  recorded via `--record` as a proper pair, supersedes the acceptance pair and re-derives the
  amortisation ratios from same-conditions rows. Blocked on **OD-1**; named in `BASELINE.md` by
  the D5 pointer edit.
- **`w-race-gate-blindspot` AC7 hookup** — when OD-1 resolves and the race doc's toolchain pin
  lands, its condition-change record uses this item's mechanism (its `:184` already assigns
  ownership here). No action in this item.

---

## Quorum verification log (pick-time quorum, iteration 42)

Reviewers `gpt5-6-sol` (OpenAI) and `gemini-3-1-pro` (Google) — both **present in both rounds**
(`absent_reviewers: []`), reject-by-default synthesis. Controller `claude-opus-5` voted pass in
both rounds. Artifacts: `.ailang/state/mission-quorum/w-bench-load-confound-2026-08-03T05-39-46Z.json`
and `…T05-57-51Z.json`. `metered=$0.2104` total ($0.0882 + $0.1222).

**Round 1 — BLOCKED, both reviewers reject. Both objections were right, and both were adopted.**

| Reviewer | Objection | Disposition |
|---|---|---|
| `gpt5-6-sol` | The recorder runs `go test -bench` with no wall-clock bound — Standing Rule 6. `go test -timeout` bounds only the test binary's clock, not the compile step or a wedged child | **ADOPTED.** `run_bounded` mirrored from `verify_ail.sh:61-74`; timeout = named refusal, zero fences; `MUT-REC-STALL` proves the kill reaches a grandchild |
| `gemini-3-1-pro` | `go env GOVERSION` probed in the caller's directory, not `--dir`; Go's toolchain switching makes it tree-dependent, so R4 would compare the variant's toolchain to itself | **ADOPTED.** Tree-dependent probes now run via `go -C <dir>`; `MUT-AB-FLOOR-SPLIT` + `MUT-PROBE-CALLER-DIR` |

The controller **measured the second premise rather than forwarding it** (P24, and re-verified the
fix's mechanism with a control in the same call: `go -C <module floor 1.26.5> env GOVERSION` →
`go1.26.5` while plain `go env GOVERSION` in this repo → `go1.26.4`, `GOTOOLCHAIN=auto` live).
The measurement made the objection **sharper than filed**: OD-1 is a parked proposal to lower this
repo's floor, so the first A/B pair straddling it is exactly the cross-toolchain pair R4 exists to
red — the round-1 design would have greened it.

**Round 2 — BLOCKED (1 pass / 1 reject). The remaining objection is OD-6.**

| Reviewer | Verdict | Objection |
|---|---|---|
| `gemini-3-1-pro` | **pass** | Non-blocking: D4's CI assertion greps a generic `probe FAILED` marker, so a mutation bypassing only the `sysctl` probes could still exit non-zero via the `AILANG_BIN` check and masquerade as a green step. Fix: assert the **specific** `hw.ncpu` failure text. Recorded as **CF-M-1** for the implementer — a real improvement, not applied here because the doc is parked |
| `gpt5-6-sol` | **reject** | R4 compares hardware and toolchain but neither load nor temporal pairing, and accepts a stale control — so an idle variant can be paired with a loaded one and pass. → **OD-6** |

**A self-catch worth recording**, because it is the trap iteration 40 named and the revision
directive warned about: while re-reading its own fix, the designer found that its first draft
bounded only the `go` invocations and left `$AILANG_BIN --version` — an env-supplied binary —
**unbounded inside the remedy for unbounded waits**. It is now bounded in D1 step 2, the deadline
bullet and the Design Freeze. *A revision is not a smaller change than an original.*

---

*Verification Log rows: 25 (P1–P25). Named mutations: 10 (`MUT-REC-SILENT-DEFAULT`,
`MUT-REC-STALL`, `MUT-EDIT-RAW-NUMBER`, `MUT-CLAIM-NONPARENT`, `MUT-CLAIM-TOOLCHAIN-SPLIT`,
`MUT-CLAIM-ORPHAN-NUMBERS`, `MUT-LEGACY-FOURTH`, `MUT-CHECK-NULL`, `MUT-AB-FLOOR-SPLIT`,
`MUT-PROBE-CALLER-DIR`). **This document authorizes no implementation: it is PARKED on OD-6, and
no milestone may be routed to a planner or an executor until the human answers it.***
