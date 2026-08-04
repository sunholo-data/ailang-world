# w-bench-load-confound — make the benchmark baseline's comparability conditions machine-recorded, and make the single-session A/B against the parent commit the only form in which benchmark evidence enters the file ~~— the only mechanically-valid form of a cost claim~~ *(title clause SUPERSEDED, round 5: the grammar validates EVIDENCE, not claims — see `4f/OD-8`)*

**Status:** **ITEM COMPLETE 2026-08-05 (iteration 50).** Both milestones landed — **BC.A′**
(`0b72019`, iteration 48) and **BC.B′** (`d357474`, iteration 49) — and controller pass **`C2b`**
discharged the last five named mutations, so the census is **16 of 16** (re-derived by command,
**P48**: the "12" carried in earlier records was a transcription, not a measurement). `AC6` is
discharged with its regime clause honoured and its vacuity arm measured (**P45**); `AC7`'s three
re-recording mutations are discharged against an honest same-session control arm (**P46**). One
doc prediction was **REFUTED** by the run that was supposed to confirm it (**P47**) and is
corrected in place rather than absorbed. Non-blocking `4f/CF-R-3` is carried out of the item.
Prior header follows.

~~**ROUTABLE — `4f/OD-8` ANSWERED (attended, 2026-08-04, charter commit `ea5e405`):
EVIDENCE, NOT CLAIMS.**~~ Mark took alternative (1) with (3) as bookkeeping — the OD-6 stamp's
wording is amended to *"mechanically complete, contemporaneous, tamper-evident **evidence** for a
cost claim"*, and claim VALIDATION (N≥3 paired runs, noise handling) is deliberately out of 4f's
scope and becomes a **named requirement of `w-agent-floor-m4`'s experimental design**. Milestones
**BC.A′/BC.B′ are ROUTABLE**; nothing in the mechanism changes. Full ratification text in the
*Open decision for the human* section. Premise rows **P9/P22 were stale and are superseded by
P37–P39** (controller-measured at routing time, iteration 48) — read the two CONTROLLER NOTEs
before *Deferred* before implementing AC6.

<details><summary>Superseded status text (round 5, iteration 47) — kept per the no-prescience convention</summary>

~~**PARKED — awaiting `4f/OD-8`**~~ (round-5 quorum, 2026-08-04, iteration 47:
BLOCKED, two rejects, both correct, both adopted in revision round 5 below. The second
reject's scope limb — R1–R6 never determine whether the observed delta supports a cost claim
— is TRUE **by ratified design** (OD-6 excluded the reviewer's third limb), but it exposes a
gap between the OD-6 stamp's wording ("mechanically valid cost claims") and what branch A
delivers: mechanically complete, contemporaneous, tamper-evident **evidence** for a cost
claim, not a validated claim. The controller is raising that wording gap as **`4f/OD-8`** —
controller-authored bookkeeping, not resolved in this document. **NOT ROUTABLE until the
human answers `4f/OD-8`.**) Underlying ratification state: OD-6 = **BRANCH A** (attended
stamp 2026-08-03) — see the *Open decision for the human* section, kept in place and marked
RATIFIED per mission convention; nothing in round 5 reopens it. The design below is the
**branch-A design**: single-session interleaved `--record-pair`, pair ID, control-reuse
rejection. ~~Milestones BC.A′/BC.B′ are routable to a planner/executor.~~ **SUPERSEDED
(round 5): routability is suspended pending `4f/OD-8`.**

</details>

**Item:** `w-bench-load-confound` (queue item 4f; charter row at
`design_docs/world-mission.md:1568`) + carry-forward **CF-K-2** (toolchain as a condition of
comparability), which folds in here; this revision also discharges **CF-M-1** (D4's specific
refusal marker) and **CF-M-2** (deadline arithmetic controller-re-run cold — P27).
**Clause:** clause-1
**Date:** 2026-08-04 (original 2026-08-03)
**Iteration:** 47 (original 42)
**Author:** `claude:claude-fable-5` (rotation designer, both the original and this branch-A
revision); premise rows P1–P25 were measured first-party by this designer at `c1e6125`; rows
P26–P31 were measured at `61348b9` (this revision) — first-party except where marked
**controller-measured (iteration 47)**.
**Estimate:** ~0.6–0.9 day total across two milestones (re-sized by the human at ratification;
the original 0.35–0.6 day sizing is SUPERSEDED — branch A reworks the recorder interface, so
the original BC.A is not a head start; see Milestones).

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

**Revision round 3 — the BRANCH-A revision (2026-08-04, iteration 47, post-ratification).**
The round-2 reject (`gpt5-6-sol`) held: R4 compared hardware and toolchain but neither load nor
temporal pairing, and accepted a stale control — so an idle variant could be paired with a
loaded control and the grammar would bless it. That scope judgment parked as OD-6; the human
ratified **branch A** (attended stamp 2026-08-03, verbatim in the OD-6 section below). This
revision implements exactly the ratified scope: (A1) one bounded capture session
`--record-pair --variant <dir> --control <dir>` replaces the two independent `--record` calls;
(A2) both benchmark binaries are **prebuilt** before any measurement, moving compilation outside
the measured window; (A3) a unique **pair ID** in both sections, an interleaved fixed leg order
**control/variant/variant/control** within the one session, and per-leg start/end timestamps +
load + competing-process snapshots; (A4) R4 additionally requires pair-ID identity and rejects
control reuse across pairs. The reviewer's third limb — a measured acceptance rule for
within-pair load divergence — is **excluded by the ratification** (the recommendation the human
ratified took neither a threshold nor an acceptance rule) and the exclusion now has measured
support: P29 shows ~1.4× within-condition spread at essentially constant load, so no data on
this rig can derive a defensible divergence threshold, and `BASELINE.md:7-8` already names
load-threshold gating a dishonest gate (S6). Consequence stated plainly here and in *What these
gates CANNOT fail*: ~~an interleaved single-session pair BOUNDS load skew — by making both
roles share one rig episode, so the confound is bounded by the session's own duration — but it
does not ELIMINATE it.~~ **SUPERSEDED (round 5, `gpt5-6-sol`): dimensionally wrong — duration
bounds temporal separation of the legs, not the magnitude of load divergence; the corrected
statement lives in D2, D5 and *What these gates CANNOT fail*.** Every sentence of the round-2 design that claimed more (the "correct
under any load" clause in D2 and Out-of-scope) is marked SUPERSEDED in place, per this
mission's convention of not rewriting decided text to look prescient. The prior deadline
arithmetic is also SUPERSEDED: prebuilding splits the bounded populations into a compile that
can legitimately take minutes and measurement legs that measured 6–8 s under load (P27), so the
single `REC_BENCH_TIMEOUT_S=600` is replaced by `REC_PREBUILD_TIMEOUT_S` + `REC_LEG_TIMEOUT_S`
(D1), discharging **CF-M-2**. `gemini-3-1-pro`'s round-2 non-blocking catch (**CF-M-1**, the
specific `hw.ncpu` refusal marker in D4's CI assertion) is folded in and discharged.

**Revision round 4 (2026-08-04, iteration 47, same iteration — the sanctioned quorum-fix
pass).** The branch-A revision's quorum BLOCKED on one reject (`gemini-3-1-pro`): the
`role: control` section recorded no counterpart commit, so the `pair_id` derivation was
recomputable only AFTER pairing — **pair-scoped** — while `MUT-PAIR-ID-SPLIT`'s and
`MUT-CONTROL-REUSE`'s predictions were **section-scoped**, and nothing anywhere defined what
the checker does with a block it cannot pair or in what order the R4 limbs run (P32; the
grep for any evaluation-order statement returned nothing against a live `R4`×59 control).
The reviewer's first proposed fix is adopted: both sections now record the pair's two commit
endpoints (`pair_variant_commit`/`pair_control_commit`, symmetric — D1 step 6), **R4c is
section-local** (derivation + role-binding) and runs BEFORE any pairing, the checker's
evaluation order is stated in D3 and frozen, and an unpairable block is R4d's **named RED,
never a skip**. Four prediction texts are re-derived against the fixed mechanism
(`MUT-PAIR-ID-SPLIT`, `MUT-PAIR-TWO-SESSIONS`, `MUT-CONTROL-REUSE`'s silence note,
`MUT-CLAIM-NONPARENT`'s dual-fire mechanism) — the unsupported R4b prediction is struck
SUPERSEDED in place. Two further controller measurements land: the prebuild's path form
fails outside the module, so D1 step 7's `-C` is load-bearing, not a spelling (P33), and a
complete four-leg interleaved session sampled a 27 s measured window (P34 — a SAMPLE, per
the iter-46 spine). Direction unchanged; no scope change; production/test-probe tally
unchanged at 13/15.

**Revision round 5 (2026-08-04, iteration 47, same iteration — the final bounded revision
pass).** The round-4 revision's quorum BLOCKED with TWO rejects; both are right, both are
adopted; the second reject's scope limb is DEFERRED to `4f/OD-8`, not resolved here (Status
above; Quorum verification log below). (1) `gemini-3-1-pro`: the bounded-execution coverage
enumerated only the `go` invocations, the four legs, and `$AILANG_BIN --version` — the
`sysctl`/`ps` probes, every `python3` helper, and every `git` invocation were **unbounded**,
so the "every individual wait is bounded … sum of its parts by construction" sentence was
FALSE as written (P35). This is the SECOND time inside this document that the remedy for
unbounded waits grew unbounded waits — the round-2 self-catch on `$AILANG_BIN --version` was
the first — and the recurrence is now named in *What these gates CANNOT fail*, because an
unfalsifiable summary sentence over a partially-covered mechanism is exactly the defect this
item exists to eliminate. Fixed: **every external-binary invocation in the session runs
through `run_bounded`**; a fourth deadline class **`REC_UTIL_TIMEOUT_S=20`** covers the
cheap utilities, justified from first-party sub-second measurements (P36: worst member
0.050 s → a 400× explicitly-labelled engineering margin); the sum-of-parts sentence is
struck SUPERSEDED and restated to what actually holds, with its one named residual (the
helper's own interpreter startup); **`MUT-REC-UTIL-STALL`** (the `MUT-REC-STALL` shape
pointed at a newly-covered invocation) proves the widened coverage discriminates.
(2) `gpt5-6-sol`, limb (a): the doc's own "honest" statement — that the pair *"bounds load
skew by the session's own duration"* — was **dimensionally wrong**: duration bounds the
TEMPORAL SEPARATION of the legs, not the MAGNITUDE of load divergence; a competitor can
spike inside the session, hit only the variant legs, and nothing in the design says
otherwise. Every site is struck SUPERSEDED and replaced with what interleaving actually
buys — monotonic rig drift affects both roles roughly equally and largely cancels; a load
episode shorter than the session, or correlated with leg boundaries, neither cancels nor is
detected — and the reviewer's three concrete pass-anyway scenarios are folded verbatim in
substance into *What these gates CANNOT fail*. Limb (b) — that R1–R6 validate provenance
and session structure but never whether the delta supports a cost claim — is TRUE, is the
ratified OD-6 scope (branch A, third limb excluded), and newly exposes that the
ratification stamp's wording promised more than branch A can deliver; that wording gap is
the human's (`4f/OD-8`), and this revision's only action on it is a sweep making the
document's own voice never claim claim-validity anywhere (title, Status, problem
statement, D2). Tally after this round: 13 production discriminators / **16** named
mutations.

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
2. **the policy** — benchmark numbers enter the file's evidence **only** as a same-rig A/B
   against the parent commit, captured as ONE interleaved single-session pair (branch A,
   ratified), and the file's structure gate goes RED when a measurement lands without its
   same-session control. (Round-5 wording note: the grammar makes the pair the only
   admissible **evidence** form; whether an observed delta supports a cost **claim** is not a
   question R1–R6 answer — see *What these gates CANNOT fail* and `4f/OD-8`.)

It deliberately does **not** re-derive the amortisation numbers (see *Out of scope*): those are
pinned to a toolchain the repo has since moved off ~~that parked decision **OD-1** exists to
change~~ (**corrected iteration 48 — OD-1 is DISCHARGED, P38**), and re-deriving them requires
the very recorder this item is building, so the follow-up is blocked on **BC.A′ landing**, not on
a human. Banking numbers under a condition you have not recorded is precisely the defect this
item exists to fix.

## Premise Verification Log

All rows measured by this designer at `c1e6125` (clean tree, verified) on 2026-08-03, on the
development rig (darwin/arm64, zsh, `PATH` prefixed with `/opt/homebrew/bin`), except where
marked **charter-cited**. Empty/negative rows carry a known-positive control in the same call.
Rows **P23–P25** were measured during revision round 2 (2026-08-03, same rig, same `c1e6125`
checkout; tree state: only this document untracked). Rows **P26–P31** were measured during the
branch-A revision (2026-08-04, iteration 47) at `61348b9` (clean tree, `git status --porcelain`
empty, verified same call) — by this designer first-party except the two rows marked
**controller-measured (iteration 47)**, whose commands and outputs are recorded verbatim from
the mission controller's first-party runs this iteration; the designer independently re-ran the
cheap corroborating halves (the greps) and they reproduced. Rows **P32–P34** were added during
revision round 4 (2026-08-04, iteration 47, same worktree state): P33/P34 are
**controller-measured (iteration 47)** verbatim; P32's corroborating greps were re-run
first-party by this designer against the round-3 text before it was edited, and reproduced the
controller's figures exactly. Rows **P35–P36** were added during revision round 5 (2026-08-04,
iteration 47), both measured **first-party by this designer** — P35 against the round-4 text
BEFORE any round-5 edit (so the absence it records is the absence the quorum rejected), P36 on
the live rig at load 4.48. Rows **P45–P48** were added at **close-out time, iteration 50** (controller pass `C2b`), all
**controller-measured first-party** at `de80792` (clean tree, `git status --porcelain` empty,
verified in the same session), outside any sandbox, each carrying its own control in the same
call. **P47 REFUTES a prediction of this document** and **P48 corrects a count this document's
own records carried**; **P47b** applies in place the P3/P4 supersession markers that iteration 49
recorded only inside P42. Rows **P42–P44** were added at **execution time, iteration 49**, all three **controller-measured first-party** at `5ff281b` (clean tree, verified in the same session), each carrying its own control in the same call; **P42 supersedes the line citations of P3 and P4** (their claims stand). Rows **P37–P39** were added at **routing time, iteration 48**, after
the human ratified OD-8 — all three **controller-measured first-party** at `ea5e405` (clean tree,
`git status --porcelain` empty, verified in the same session), each carrying its own control in
the same call. They supersede **P9** and **P22**, which had gone stale against HEAD while this
document sat parked (Gate-2 rule 3b(vi); see the two CONTROLLER NOTEs before *Deferred*).

| # | Claim | How verified (command) | Observed | Verdict |
|---|---|---|---|---|
| P1 | The harness records NO conditions | `grep -nEi "load\|uptime\|sysctl\|GOVERSION\|toolchain\|nproc\|benchstat\|git rev-parse\|hw\.\|ncpu" scripts/bench_worldd.sh > /tmp/f1out 2>&1; echo "rc=$?"` with same-call control `grep -c -i "benchmark" scripts/bench_worldd.sh` | `rc=1`, zero output; control = **15** — the instrument sees the file, the absence is real | CONFIRMED |
| P2 | The harness is 49 lines and has exactly one mode | `wc -l scripts/bench_worldd.sh` + full first-party read | `49`; usage accepts only `--smoke`; the mode runs `go test -bench . -benchtime 1x -run '^$' ./host/daemon/` and asserts a hardcoded 10-name manifest (lines 15–26, 37–47) | CONFIRMED |
| P3 | CI runs the smoke gate, and it is a NAME gate, not a measurement gate | First-party read of `.github/workflows/ci.yml` (whole file) | Job `go-verify`, step `worldd benchmark smoke gate` at ~~`ci.yml:88-89`~~ **`:101-102`** runs `./scripts/bench_worldd.sh --smoke`. No numbers recorded, no thresholds evaluated anywhere in the workflow. **CITATION SUPERSEDED BY P42 (iteration 49), MARKER APPLIED HERE (iteration 50) — see P47b**: iteration 49 recorded the repair inside P42's own row and left this row reading `:88-89` under a `CONFIRMED` verdict, so a reader arriving at P3 got the stale citation with no signal. Re-measured at `de80792`: `:102` | CONFIRMED IN CLAIM, CITATION REPAIRED — read, not grepped, per the controller's F2 truncation warning |
| P4 | CI checkouts are shallow (no `fetch-depth`) | Same read | Both jobs use bare `actions/checkout@v4` (`ci.yml:13`, ~~`:53`~~ **`:55`**) with no `fetch-depth` key. **CITATION SUPERSEDED BY P42, MARKER APPLIED HERE (iteration 50) — same defect as P3.** Re-measured at `de80792`: `:13` and `:55` | CONFIRMED IN CLAIM, CITATION REPAIRED — constrains the checker to in-file structure (see D3) |
| P5 | The comparability promise and the S6 honesty line are real | First-party read of `bench/BASELINE.md:5-8` | `:5` "Later sprints diff against this file on the same development rig"; `:7-8` "Noise-gating a shared runner would be a dishonest gate (S6)" | CONFIRMED |
| P6 | The conditions ARE recorded today, but BY HAND | First-party read of `bench/BASELINE.md:16-26` | ~~`:18` `Go: go version go1.26.4 darwin/arm64`~~ **STALE — SUPERSEDED BY P40 (iteration 48)**: `:18` reads **`go1.25.6`** since `f19acac`, the same commit that staled P9 and P22. `:22` `Rig load at measurement: load averages: 5.22 4.99 5.91` — **still exact** (control leg, unchanged at HEAD). **The claim this row exists to make is UNAFFECTED and is if anything strengthened**: the line is hand-typed prose produced by nothing, which is exactly why it can drift silently against the toolchain it purports to describe | SUPERSEDED IN VALUE, CONFIRMED IN CLAIM — see P40 |
| P7 | The confound is a standing condition, live right now | `uptime; sysctl -n hw.ncpu hw.model; ps -Ao pid=,ppid=,pcpu=,comm= \| sort -k3 -rn \| head -8` at 2026-08-03 07:27 CEST | `load averages: 3.93 3.20 3.03` on 16-core Mac16,9; top consumers: **100.0%** `/var/folders/.../go-build200642988/b001/exe/solution` (pid 28606), **99.6%** `node` (pid 90462), **77.0%** `ollama` | CONFIRMED |
| P8 | The generic temp-path binary CAN be attributed one level up, mechanically | `ps -o comm= -p 28600` (ppid of the 100%-CPU `solution` binary) | `go` — the parent comm resolves; a `go`-spawned test binary is identifiable as such without guessing which mission owns it | CONFIRMED — D1's parent-comm field is feasible |
| P9 | Toolchain state and the floor | `go version; go env GOVERSION GOOS GOARCH; sed -n '1,5p' go.mod` | ~~`go1.26.4 darwin/arm64`; `go.mod:3` = `go 1.26.4`~~ **STALE — SUPERSEDED BY P37 (iteration 48)**: `go.mod:3` now reads **`go 1.25.6`** (landed `f19acac`, iteration 46). The rest of the row survives and is the reason the staleness is narrow: the `go` directive is a **floor** (iter-40, measured), so `go env GOVERSION` in this repo still reads **`go1.26.4`** and every *recorded* toolchain value in this design is unchanged. ~~lowering the floor is **OD-1, parked**~~ → **OD-1 DISCHARGED (P38)** | SUPERSEDED — see P37 |
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
| P22 | OD-1 is parked and constrains scope | **charter-cited** (`world-mission.md`, iter-40/41 status rows; controller F8) | ~~OD-1 = lower the `go.mod` floor 1.26.4 → 1.25.6; awaiting Mark; this loop may not decide it~~ **STALE — SUPERSEDED BY P38 (iteration 48): `4e/OD-1` is DISCHARGED.** Mark ratified it and it landed at iteration 46 (`f19acac`). Nothing in this item is constrained by it any more; the *Deferred* limb (iii) is now blocked on BC.A′ existing, not on a human | SUPERSEDED — see P38 |
| P23 | A bounded-execution precedent exists in-repo, and it is process-group-killing | Full first-party read of `scripts/verify_ail.sh:44-46,54-74` | `run_bounded` (V26): python3 `Popen(..., start_new_session=True)` puts the child in its own process group; on expiry `os.killpg(..., SIGKILL)` kills the WHOLE group, a named `✗ TIMEOUT after Ns: <cmd>` goes to stderr, exit **124**; deadlines are hardcoded constants (`GATE_LEG_TIMEOUT_S=120`, `GATE_TEST_TIMEOUT_S=180`), not env knobs; every binary invocation in both gate legs runs through it | CONFIRMED — D1 **mirrors this helper verbatim** rather than inventing a third form |
| P24 | `go env GOVERSION` is directory-sensitive; the switching mechanism is live; the recorder's `-C` form works; the AC6 fixture needs no network | Temp module with a `go 1.26.5` floor: `(cd "$d" && go env GOVERSION GOTOOLCHAIN)` · same in the repo · same in module-less `/private/tmp` · `go -C "$d" env GOVERSION` · `go -C <repo> test -c ./host/daemon/` · `ls "$(go env GOMODCACHE)/golang.org"` | Temp module → **`go1.26.5`**, `GOTOOLCHAIN=auto`; repo (floor `go 1.26.4`) → **`go1.26.4`**; no-module dir → `go1.26.4`; `go -C` env form → `go1.26.5` rc=0; `go -C … test` form compiles OK; toolchains **go1.25.6 AND go1.26.5 already cached** in GOMODCACHE | CONFIRMED — reproduces the controller's measurement first-party; a caller-cwd probe records the wrong tree's toolchain, and both floors the AC6 fixture uses are cached (no download) |
| P25 | The benchmark deadline is derivable from measurement, and 120 s would be too small | `uptime`; `time ( go test -c -o … ./host/daemon/ )` cold then warm; iteration-weighted sums computed from the recorded `BASELINE.md` rows (python3) | At load 2.87: cache-cold compile **128.85 s wall** (1.64 s user / 4.67 s system, 4% CPU — cache-cold and rig-contended); warm re-compile 0.174 s; recorded idle full-run binary time 3.433 s (`BASELINE.md:276`); loaded iteration-weighted sum **4.217 s** vs idle **1.924 s** = 2.19×; worst observed p95 inflation 6.06× (iter-39) | CONFIRMED as a measurement — but its **consequence is SUPERSEDED** by branch A: the 128.85 s belongs to a compile that A2 moves OUT of the measured window, so the single `REC_BENCH_TIMEOUT_S=600` bounded two populations with different physics. Replaced by the P27-derived structure in D1 (CF-M-2 discharged) |
| P26 | `go test -c` produces a standalone, cwd-independent, re-runnable benchmark binary — A2 is feasible | **controller-measured (iteration 47):** `go test -c -o "$BIN" ./host/daemon/` → rc=0, 6 s wall (warm), 17,023,394 B; `( cd /Users/voightkampff && "$BIN" -test.bench 'BenchmarkHeadRead\|BenchmarkBrokerFSRead' -test.benchtime 20x -test.run '^$' )` → rc=0 with p50/p95 reported; a SECOND invocation of the SAME binary from a THIRD cwd (`/tmp`) → rc=0. Designer re-ran the corroborating grep first-party at `61348b9`: `grep -nE "AILANG_BIN\|testdata\|\.\./\|os\.Getwd\|Getenv\|filepath\.Abs" host/daemon/bench_test.go` → rc=1, zero output; same-call control `grep -cE "b\.ResetTimer\|testing\.B"` → **20** | Repeat invocation from arbitrary cwd works — the interleaving depends on this; no cwd/env/testdata dependency exists in the benchmark file. On a compiled binary the flags are `-test.bench`/`-test.benchtime`/`-test.run`, NOT the `go test` spellings — the recorded `invocation:` must reflect what was executed | CONFIRMED |
| P27 | A prebuilt-binary measurement leg is seconds, not minutes, even on a loaded rig — a SAMPLE, not a bound | **controller-measured (iteration 47):** five full-leg runs of `"$BIN" -test.bench . -test.benchtime 200x -test.run '^$'` (all ten benchmarks) on a LOADED rig (`sysctl -n vm.loadavg` `{ 6.25 4.99 4.86 }`…`{ 5.58 4.97 4.86 }`; competitors ≥25% pcpu: `ollama` 140.4%, a `go-build/.../exe/solution` 97.6%, `node` 97.5%) | Leg wall times: **8, 6, 6, 7, 6 s**. Prebuild: 6 s warm, 1 s hot. Five runs on one loaded rig at one commit establish an ORDER OF MAGNITUDE (seconds, not minutes) and nothing about the tail — the iter-46 spine: a range you stopped measuring at is just a wider single number. Used ONLY to size the leg deadline's species; the margin over it is labelled engineering margin, not measurement (D1) | CONFIRMED — discharges **CF-M-2** (the deadline arithmetic is now controller-re-run, cold, this iteration) |
| P28 | The current script has NO record/pair/check mode — all of D1/D3 is net-new surface, nothing migrates | `grep -n -- '--record\|record-pair\|--check-claims' scripts/bench_worldd.sh` → rc=1, zero output; same-call control `grep -c 'bench'` → **7**; `wc -l` → **49** (re-run first-party at `61348b9`; matches P2 and the controller's iteration-47 run) | rc=1 with live control; the script still implements only `--smoke` | CONFIRMED |
| P29 | Within-condition noise is ~1.4× at essentially constant load — no within-pair load-divergence threshold is derivable from data on this rig | **controller-measured (iteration 47):** across the P27 runs plus the P26 spot check, at load averages varying only 5.54–6.39, `BenchmarkBrokerFSRead` p95 read 5.100, 4.088, 3.530, 3.721, 3.833, 4.225 ms | **~1.4× spread with conditions held roughly constant** (reference: idle M3.C baseline 0.7472 ms; iter-39 loaded artefact 4.529 ms = the 6.06× confound). A divergence rule would have to separate a real regression from noise that is itself 1.4× at constant load — the measured reason the ratified scope excludes the reviewer's third limb, alongside `BASELINE.md:7-8` (S6) | CONFIRMED |
| P30 | `verify_ail.sh`'s `run_bounded` is real and as described — the mirror source is unchanged | First-party re-read of `scripts/verify_ail.sh:44-80` at `61348b9` | `run_bounded()` at `:61-74`: python3 `Popen(..., start_new_session=True)`, `os.killpg(..., SIGKILL)` of the whole group on expiry, named `✗ TIMEOUT after Ns` to stderr, exit 124; hardcoded `GATE_LEG_TIMEOUT_S=120` / `GATE_TEST_TIMEOUT_S=180` at `:44-45`, not env knobs | CONFIRMED — the "mirrored, not sourced or extracted" argument stands unchanged |
| P31 | `bench/BASELINE.md` contains NO conditions block today — schema `bench-conditions/2` can be required exclusively, no schema-1 migration exists | `grep -n 'bench-conditions' bench/BASELINE.md` → rc=1, zero output; same-call controls: `grep -c '^Benchmark'` → **30** (the P14 value) and `python3 -c 'import secrets; print(len(secrets.token_hex(16)))'` → **32** (the pair-nonce source exists) | The round-2 schema/1 never shipped (PR #31 landed this DOC only); the checker needs no legacy-schema branch, and python3 `secrets` is available for the pair nonce | CONFIRMED |
| P32 | The round-3 text's control section recorded NO counterpart commit, and NO evaluation order was stated for R4a…R4f — the `pair_id` derivation was PAIR-scoped while three mutation predictions were SECTION-scoped, and the unpairable block was undefined | **found by quorum round 4 (`gemini-3-1-pro`, reject), measured by the controller, corroborating greps re-run first-party by this designer against the round-3 text:** extract the `schema: bench-conditions/2`…`conditions_sha256` block and grep `commit\|variant` → only `commit:`/`parent:` plus role/leg-order comments, no counterpart-commit field; sweep for any evaluation-order statement (`evaluation order\|evaluated (first\|before\|after)\|before any pairing`…) → **rc=1, zero hits**; same-call known-positive control: `grep -c 'R4'` → **59** | The reviewer's premise is TRUE; its stated mechanism ("mathematically impossible to evaluate") is OVERSTATED for a well-formed pair — the checker reads the whole file, grouping by `pair_id` supplies both commits, so happy-path R4c was evaluable as written. The REAL defect is one level in: under `MUT-PAIR-ID-SPLIT` and `MUT-CONTROL-REUSE` the checker cannot FORM the pair, so it has no `variant_commit` for the orphaned section, the predicted R4c REDs could not fire as specified, and "cannot evaluate" was one implementer-reading away from a silent skip — a rule whose failure mode is undefined for exactly the input the mutation produces | CONFIRMED — fixed this round: symmetric `pair_variant_commit`/`pair_control_commit` in BOTH sections (D1 step 6), section-local R4c, frozen evaluation order (D3), R4d's unpairable named RED |
| P33 | The prebuild's PATH form FAILS across worktrees; only the `go -C` form compiles the control tree — D1 step 7's `-C` is load-bearing, not an equivalent spelling | **controller-measured (iteration 47):** `go test -c -o "$BIN" /Users/voightkampff/dev/sunholo-data/.wt-world-iter47-ctl/host/daemon/` → **rc=1**, `directory ../.wt-world-iter47-ctl/host/daemon outside main module or its selected dependencies`, `FAIL … [setup failed]`; same compile via `go -C /Users/voightkampff/dev/sunholo-data/.wt-world-iter47-ctl test -c -o … ./host/daemon/` → **rc=0**, 2 s, 17,023,394 B (the rc=0 run is the same-call known-positive control for the rc=1 form) | Extends P24 from `go env`/`go test` to `go test -c` and to a SECOND worktree. An executor who "simplifies" the `-C` away gets a setup failure that reads like a broken tree | CONFIRMED |
| P34 | A complete four-leg interleaved session measured a 27 s window end-to-end — a SAMPLE, not a bound | **controller-measured (iteration 47):** two prebuilt binaries (variant at `61348b9`, control at its parent `f19acac`; control-tree prebuild via `go -C`: 2 s), legs in the frozen C/V/V/C order, each `-test.bench . -test.benchtime 200x -test.run '^$'`; load averages 4.39 → 4.87 across the session | Leg walls **7, 7, 6, 7 s**; TOTAL measured window **27 s**. One session on one loaded rig at one commit pair — it establishes that a whole branch-A session's measured window is TENS OF SECONDS, which is the fact `REC_LEG_TIMEOUT_S`'s margin reasons about, and it establishes nothing about the tail (the iter-46 spine: record the sample and its spread, never an interval you assume bounds anything) | CONFIRMED |
| P35 | The round-4 text bounded NONE of the `sysctl`/`ps`/`python3`/`git`/`date` invocations — the round-5 reject's premise is exactly right, and the "sum of its parts" sentence was false as written | **first-party, against the round-4 text before any round-5 edit:** collect every bounding-related line (`grep -nE 'run_bounded\|bounded runner\|bounded execution\|Bounded execution'` → **15** lines), then filter those lines for `sysctl\|python3\|\bgit\b\|\bps\b\|\bdate\b`; same-call known-positive control `grep -c 'sysctl'` → **20** | The filter returned exactly TWO lines — the premise rows P23/P30, which describe `verify_ail.sh`'s helper and mention python3 only as the HELPER's implementation language. **Zero** of the 15 bounding lines covers a recorder `sysctl`, `ps`, `git`, `date`, or python3 hashing/nonce invocation; the covered enumeration (round-4 D1) listed only `go env`, prebuilds, legs, `$AILANG_BIN --version`. The control (sysctl × 20) proves the instrument sees the file — the absence is real | CONFIRMED — fixed this round: every external invocation bounded, `REC_UTIL_TIMEOUT_S`, `MUT-REC-UTIL-STALL` |
| P36 | Every utility-class invocation is sub-second on this rig, even loaded — 20 s is a ≥400× margin, and reusing the 120 s probe constant would be indefensible padding | **first-party, load 4.48 (uptime in same call):** python3 `subprocess` timing of each member: `sysctl -n vm.loadavg`/`hw.ncpu`/`hw.model` · `ps -Ao pid=,ppid=,pcpu=,comm=` · `date -u` · `git status --porcelain` · `git rev-parse HEAD`/`HEAD^`/`--git-dir` · `python3 -c` secrets nonce · `python3` hashlib SHA-256 over a 17,825,792 B file (the prebuilt-binary size class, P26) | `0.005 / 0.004 / 0.003` s · `0.044` s · `0.004` s · `0.016` s · `0.011 / 0.012 / 0.011` s · `0.042` s · `0.050` s — **worst member 0.050 s** (the 17 MB hash). One run per member on one loaded rig: a SAMPLE establishing the class's order of magnitude (tens of milliseconds), per the iter-46 spine — used only to size the deadline's species; `REC_UTIL_TIMEOUT_S=20` is an explicitly-labelled **400× engineering margin** over the worst sampled member, not a measured tail | CONFIRMED |
| P37 | **The AC6 floor-straddle exists at HEAD and needs only ONE throwaway commit, because `go -C <dir> env GOVERSION` resolves each tree's own floor** (supersedes P9; re-measures P24's mechanism at today's HEAD) | **controller-measured (iteration 48)**, `ea5e405`, clean tree: `git worktree add --detach <repo>/../.bench-probe-iter48 HEAD` (**sibling of the repo, never `/tmp`**) → `grep '^go ' $WT/go.mod; go -C "$WT" env GOVERSION` → python3 in-place edit of the floor `1.25.6`→`1.26.5` with an `assert s2 != s` **mutation-applied control** → re-read both → `grep '^go ' go.mod; go -C "$(pwd)" env GOVERSION` in the main tree as the untouched control → `git worktree remove --force` | Control leg: floor `go 1.25.6` → **`go1.26.4`**. Mutated leg: floor `go 1.26.5` → **`go1.26.5`**. Mutation control printed `mutation applied` (an unmatched regex would have aborted, per the iter-45 spine). Main tree after: `go 1.25.6` / `go1.26.4` — **unchanged**. So HEAD is a valid AC6 *control* as-is: the fixture is `variant` = one commit raising the floor to `1.26.5`, `control` = HEAD (its parent), and the `go1.26.4` vs `go1.26.5` straddle `MUT-AB-FLOOR-SPLIT` predicts is genuine. Toolchain cache re-listed in the same session: **6 cached** (`1.24.9`, `1.25.6`, `1.26.0`, `1.26.2`, `1.26.3`, `1.26.5`) — `1.26.5` present, so no network (P24 holds); `1.26.4` is the *installed* system Go, not a cache entry | CONFIRMED |
| P38 | `4e/OD-1` is DISCHARGED, so nothing in this item is blocked on a human toolchain decision (supersedes P22) | **controller-measured (iteration 48)**: `sed -n '1759,1762p' design_docs/world-mission.md` (charter item 4e row) · `sed -n '1,6p' go.mod` | Charter reads `**ITEM COMPLETE 2026-08-04 (iter-46) … 4e/OD-1 AND 4e/OD-2 BOTH DISCHARGED**`; `go.mod:3` = `go 1.25.6`, i.e. the floor OD-1 was *about* has already moved. The *Deferred* limb (iii) is therefore blocked on **BC.A′ existing**, not on Mark | CONFIRMED |
| P41 | **AC6's known-positive CANNOT FIRE under the toolchain regime this repository requires, and the only straddle achievable on this rig is between two deny-listed compilers.** Five quorum rounds assumed the fixture was viable; nobody ran it | **controller-measured (iteration 48)**, prompted by the planner's DD-1 and then carried past it: (1) `for tc in go1.26.5 go1.25.6 go1.26.4; do GOTOOLCHAIN=$tc go test ./host/store/ -run TestToolchainCanary; done`; (2) a variant worktree with floor `1.26.5` and the repo as control, each read under BOTH `GOTOOLCHAIN=auto` and `GOTOOLCHAIN=go1.25.6`; (3) a `toolchain go1.24.9` / `toolchain go1.25.6` directive pair under `auto`, with a no-directive third tree as control; (4) `sed -n '35,41p' scripts/verify_go.sh` | (1) canary **FAILS on go1.26.5 AND on go1.26.4**, `ok` only on go1.25.6 — so the rig's DEFAULT toolchain reds the canary **by iteration-46 design**, and `verify_go.sh:40` deny-lists **`go1.26.0`–`go1.26.5`** entire. (2) Under `auto`: variant `go1.26.5`, control `go1.26.4` — a real straddle, **both deny-listed**. Under `GOTOOLCHAIN=go1.25.6` (the regime the gate requires and CI pins): **BOTH trees report `go1.25.6`** — `go env GOVERSION` returns the PIN, not the tree's requirement, so the straddle **silently collapses** and `MUT-AB-FLOOR-SPLIT` would come back **GREEN**: a vacuous pass in the one AC that exists to prove the probe is not vacuous. (3) The `toolchain` directive does **not** rescue it — all three trees returned `go1.26.4`, because a `toolchain` directive, like a floor, only ever forces **UPWARD**; downward selection is impossible from `go.mod` alone. The only cached toolchain above the local `go1.26.4` is `go1.26.5`, which is deny-listed, and the control is always the local toolchain — so **no deny-list-free straddle exists on this rig** | CONFIRMED — AC6 amended below; **this constrains only AC6, not the mechanism** |
| P40 | **The stale-premise sweep, done with the instrument instead of the name-list — and the record of why the first pass found two of three** | **controller-measured (iteration 48)**, after the PLANNER refuted the controller's claim that P9/P22 were the only stale rows: `git diff --name-only c1e6125..HEAD -- ':!design_docs'` (the doc's OLDEST measurement base, which rows P1–P22 declare) with the control `git diff --name-only 61348b9..HEAD -- ':!design_docs'` (the base rows P26+ declare) · `for rev in c1e6125 61348b9 HEAD; do git show "$rev:bench/BASELINE.md" \| sed -n '18p'; done` · `git log --oneline -3 -- bench/BASELINE.md` | From `c1e6125`: **8 non-doc files changed** (`ci.yml`, `bench/BASELINE.md`, `go.mod`, four `host/store/*`, `scripts/verify_go.sh`). From `61348b9`: **ZERO** — which is precisely why sweeping from the wrong base reads as "nothing is stale". `BASELINE.md:18` = `go1.26.4` at `c1e6125`, `go1.25.6` at `61348b9` and at HEAD; `git log` puts the change in **`f19acac`** — so P6, P9 and P22 are **one drift event with three faces**, not three findings. Control leg `:22` unchanged at both revisions, so the instrument is reading real content. **Honest tally: the controller's first pass found 2 of 3, because it swept the rows iteration 47 had NAMED rather than the commit that caused them; the third was found by the planner and confirmed here.** The durable instrument is the diff from the doc's own **oldest** declared measurement base — a document with rows measured at several bases is only as fresh as its oldest one | CONFIRMED |
| P39 | Branch A is still entirely unbuilt at HEAD — nothing from any earlier milestone shipped (re-confirms P28 at `ea5e405`) | **controller-measured (iteration 48)**: `grep -c -- "<pat>" scripts/bench_worldd.sh` for `--record-pair`, `--check-claims`, `--record ` with the **known-positive control** `--smoke` in the same call | `--record-pair` **0**, `--check-claims` **0**, `--record ` **0**, control `--smoke` **2** — the instrument sees the file, so the three zeros are measurements rather than a failed grep. Nothing to salvage, nothing to collide with | CONFIRMED |
| P42 | **The freshness sweep found a FOURTH face of the RG.A drift, and iteration 48's sweep could not have caught it because it swept the NAMED rows rather than the CHANGED files** (supersedes the citation halves of **P3** and **P4**; the claims themselves stand) | **controller-measured (iteration 49)** at `5ff281b`, clean tree: swept from the OLDEST declared base — `git diff --name-only c1e6125..HEAD -- ':!design_docs'` → **9 files**, against `git diff --name-only ea5e405..HEAD -- ':!design_docs'` → **1 file**; then `grep -n 'bench_worldd.sh --smoke' .github/workflows/ci.yml` at HEAD, at `ea5e405` and at `c1e6125`; and `grep -n 'actions/checkout@v4' .github/workflows/ci.yml`. Instrument control in the same call: `git diff --stat c1e6125..HEAD -- scripts/bench_worldd.sh` → **+414**, a file known to have changed | HEAD **`:101-102`**; `ea5e405` **`:102`**; `c1e6125` **`:89`**. So **P3's `ci.yml:88-89` was ALREADY STALE at `ea5e405`** — the base P37–P39 were measured at — and iteration 48's sweep missed it. `P4`'s second bare checkout moved `:53` → **`:55`**. Both rows' CLAIMS survive unchanged (the smoke gate is a NAME gate; both checkouts are bare with no `fetch-depth`); only their line citations rotted. The **sprint planner** had independently measured `:101-102` into the plan, so the plan was fresher than the doc it plans | CONFIRMED — P3/P4 citations REPAIRED here; the general lesson is that a document degrades in the rows nobody has a reason to re-read, so the sweep must key on the DIFF, never on the list of rows someone already suspected |
| P43 | **`--check-claims` was ADVERTISED but not implemented at `5ff281b`, and the advertisement was indistinguishable from a typo** | **controller-measured (iteration 49)**: `grep -c -- '--check-claims' scripts/bench_worldd.sh` then `./scripts/bench_worldd.sh --check-claims; echo rc=$?` with same-call control `./scripts/bench_worldd.sh --definitely-not-a-mode; echo rc=$?` | grep = **1**, and the single hit is the `usage()` line at `:9`. Both invocations print the SAME usage text and exit **rc=2** — the advertised mode and a garbage flag are byte-identical in behaviour | CONFIRMED — a grep count is not a capability check; BC.B′ makes the advertisement true |
| P44 | **`MUT-CLAIM-TOOLCHAIN-SPLIT`'s literal is a NO-OP against a pair recorded under this repository's own pinned regime — the mutation as written would pass vacuously** (`4f/CF-S-3`; same class as P41) | **controller-measured (iteration 49)** on the real acceptance pair: both conditions blocks record `goversion: go1.25.6` (the session ran under `GOTOOLCHAIN=go1.25.6`, what `verify_go.sh` requires and CI pins job-wide), so the doc's *"edit the control block's `goversion` to `go1.25.6`"* changes nothing. Substituted `go1.26.4`, resealed `conditions_sha256`, and asserted the replacement count = 1 and a differing file sha256 before believing the result | With the doc's literal: no textual change, gate stays GREEN. With `go1.26.4`: **R4e REDs `toolchain mismatch inside claimed A/B pair`** and no other limb fires. The evaluator found the same defect independently (NB-3) | CONFIRMED — the mutation spec must say *a value that DIFFERS from the variant's recorded `goversion`*, never a frozen literal; a fixture literal written under the old `go1.26.4`-default assumption goes vacuous the moment the regime changes |
| P45 | **`AC6` DISCHARGED end-to-end on a REAL floor-straddling pair, and the vacuity arm the regime clause warns about was measured in the same session rather than trusted** (`MUT-AB-FLOOR-SPLIT` + `MUT-PROBE-CALLER-DIR`; `4f/CF-S-2` honoured) | **controller-measured (iteration 50)** at `de80792`, clean tree, outside any sandbox. Fixture per the iteration-48 re-pointing: `control` = detached worktree at HEAD, `variant` = detached worktree at HEAD + one never-pushed commit raising the `go.mod` floor `1.25.6`→`1.26.5`, both **siblings of the repo, never `/tmp`**, the floor edit applied through a python3 `assert s2 != s` mutation-applied control. Straddle read BOTH ways before recording: `GOTOOLCHAIN=auto go -C <tree> env GOVERSION` and `GOTOOLCHAIN=go1.25.6 go -C <tree> env GOVERSION`. Then one bounded `--record-pair` session invoked `GOTOOLCHAIN=auto` **explicitly**, emission appended to `bench/BASELINE.md` under a label naming the deny-listed toolchains, `--check-claims` run before AND after the append, `bench/BASELINE.md` sha256-restored between arms. `MUT-PROBE-CALLER-DIR` then dropped `-C "$dir"` from both `go env` probes (2 occurrences, replacement count asserted, recorder sha256 `9d5e62ce`→`4d4016d7`) and re-recorded | **Straddle under `auto`: variant `go1.26.5`, control `go1.26.4`. Under `GOTOOLCHAIN=go1.25.6`: BOTH `go1.25.6`** — P41 reproduced first-party, so CF-S-2 is a measured precondition here and not an inherited warning. The recorded pair carries `goversion: go1.26.5` / `go1.26.4`, one shared `pair_id`, parent edge `variant.parent == control.commit`, interleave control `1/4`+`4/4` / variant `2/4`+`3/4`. **Control arm (rule 3d): the gate on the un-appended file is `rc=0`, `✓ PASSED`.** With the pair appended: **`rc=1`, EXACTLY ONE violation, `✗ toolchain mismatch inside claimed A/B pair`** — R4e alone, no other limb. `MUT-PROBE-CALLER-DIR`: both sections now record the CALLER's `go1.26.4` despite the variant tree's 1.26.5 floor, and the gate **GREENS a genuinely cross-toolchain pair** (`rc=0`, `✓ PASSED`) — RED at review, exactly as specified. Both files restored byte-identically (`shasum -c` OK, `git status --porcelain` empty) | CONFIRMED — the known-positive fires because of the `-C` placement, not by luck: removing the mechanism (the `-C`) flips the outcome, and removing the regime (`auto`) flips it the other way. Two arms, two directions, one session |
| P46 | **`AC7`'s three re-recording mutations DISCHARGED against an honest same-session control arm, and the doc's named silencing attempt relocates the RED rather than removing it** (`MUT-PAIR-TWO-SESSIONS`, `MUT-PAIR-SEQUENTIAL`, `MUT-PAIR-INLINE-BUILD`) | **controller-measured (iteration 50)** at `de80792`. Second fixture, deliberately **same-toolchain** so only the intended limb can fire: `variant` = worktree at `de80792`, `control` = worktree at its parent `a63293f`, both recorded under `GOTOOLCHAIN=go1.25.6` (the repo's required regime; `auto` is scoped to AC6's session alone). TWO real bounded sessions A and B on the same commits. Control arm first: session A's WHOLE pair appended → gate must green. `MUT-PAIR-TWO-SESSIONS` = splice A's control half with B's variant half (both halves individually pristine), then the doc's silencing attempt = retype the variant's `pair_id` to the control's AND recompute its `conditions_sha256` so R1's digest limb cannot mask R4b/R4c. `MUT-PAIR-SEQUENTIAL` = 4 anchored edits making the harness honestly un-interleaved (`leg_roles` → `control control variant variant`, `leg_order` string, and BOTH `emit_role` leg assignments). `MUT-PAIR-INLINE-BUILD` = prebuild skipped, each leg run as `go -C <dir> test -bench . -benchtime 200x -run '^$' ./host/daemon/`, `invocation:` honestly recording that spelling. Every mutation asserted its anchor count and a differing sha256 before its result was believed; `bench/BASELINE.md` sha256-restored between every arm | **Control arm: `rc=0`, `✓ PASSED`, `2 well-formed pairs`** — so every RED below is the splice/edit, not the fixture. `MUT-PAIR-TWO-SESSIONS`: **`rc=1`, exactly 2 violations, `✗ unpairable conditions block — no counterpart shares its pair_id` on BOTH spliced sections**, and **R4b did NOT fire** — the round-4 supersession validated on a real pair. Silencing attempt: **`rc=1`, 3 violations** — `✗ pair_id does not derive from recorded session fields` (R4c, on the retyped section only) **+** `✗ pair identity mismatch — sections are not from one session` (R4b) **+ a third the doc never named: `✗ legs not interleaved`**, because two sessions recorded 36 s apart occupy disjoint time windows and cannot fake an interleave. `MUT-PAIR-SEQUENTIAL`: **`rc=1`, exactly ONE violation, `✗ legs not interleaved — control legs are not outermost`** on an otherwise well-formed pair (control `1/4`+`2/4`, variant `3/4`+`4/4`). `MUT-PAIR-INLINE-BUILD`: **`rc=1`, 6 violations — `✗ invocation is not a prebuilt-binary invocation` on BOTH sections (R1) plus R2's orphan cascade on all four raw blocks**, since an invalid conditions block cannot authorize its numbers. Recorder restored byte-identically after each (`shasum -c` OK) | CONFIRMED — and the unnamed third RED in the silencing arm is the strongest of the three: independently-recorded halves are separated in TIME, so R4f catches them even if every identity field is forged |
| P47 | **REFUTED — `MUT-PAIR-INLINE-BUILD`'s stated SECONDARY OBSERVABLE does not fire on this rig, and a reviewer who looked for it instead of the rule would have cleared a session that genuinely folded compilation into the measured window** | **controller-measured (iteration 50)**, same run as P46. The doc predicts *"leg-1 elapsed jumps from seconds to a compile-bearing figure"*. Compared per-leg `elapsed_s` of the honest session A against the inline-build session, then took the control that separates *"compilation is cheap"* from *"compilation never happened"*: the recorder's own `prebuild_elapsed_s` field, which prices exactly the compile that the mutation folds in | Honest session A legs: **7, 7, 7, 7 s**. Inline-build legs: **8, 8, 9, 8 s** — a **1.1–1.3×** move. Control: `prebuild_elapsed_s` = **2, 2** (session A), **1, 1** (session B), **1, 1** (the AC6 session, whose variant tree compiled under a different toolchain entirely). So a full compile of these trees costs **1–2 s** against 7–9 s legs, i.e. the folded-in compilation sits **inside the ~1.4× within-condition noise this document already measured (P-row on within-condition spread)**. The prediction is false here not because compilation is absent but because the Go build cache makes it nearly free | **REFUTED — corrected in the mutation bullet rather than absorbed.** The transferable point is this document's own spine turned on itself: **a secondary observable that a cache can erase is not evidence, and it is worse than no observable, because it offers a reviewer a plausible reason to stop looking.** `R1`'s rule-based RED is what has teeth, and it fired on both sections. Same class as P41/P44 — a fixture expectation frozen under conditions that have since changed |
| P47b | **The iteration-49 repair of P3/P4 was recorded where the FINDER was looking, not where a READER would look — so both rows still served a stale citation under a `CONFIRMED` verdict** | **controller-measured (iteration 50)**: freshness sweep from the doc's OLDEST declared base per Gate-2 rule 3b(vi-b) — `git diff --name-only c1e6125..HEAD -- ':!design_docs'` → **9 files**, with the instrument control `git diff --stat c1e6125..HEAD -- scripts/bench_worldd.sh` → **+652** (a file known to have changed; it read +414 at iteration 48, so the control also proves the sweep is seeing today's tree). Then re-read the cited lines at HEAD and at `5ff281b`, and re-read rows P3/P4 themselves | The sweep's *finding* was correct at iteration 49 and is unchanged at HEAD (`:102` for the smoke gate, `:13`/`:55` for the checkouts — BC.B′'s +21 ci.yml lines land after the smoke step and moved nothing). But **P42's row said "P3/P4 citations REPAIRED here" while rows P3 and P4 were never touched** — no `SUPERSEDED` marker, unlike P6/P9/P22 which carry one. A reader arriving at P3 reads `:88-89` and `CONFIRMED` | CONFIRMED — markers now applied in place. **A supersession that lives only in the superseding row is not a repair; it is a note to the person who already knew.** The convention this document already uses for P6/P9/P22 is the correct one, and it was skipped precisely because the finder's attention was on the row being written |
| P48 | **The mutation census was a TRANSCRIPTION, not a measurement — the denominator carried in two charter records and one queue row is wrong** | **controller-measured (iteration 50)**, re-derived by command at the moment of writing per Gate-2 rule 3b(v-b): `jq -r '.mutation_tally' .ailang/state/sprints/w-bench-load-confound.plan.json` and `jq -r '.mutation_tally.evidence[], .mutation_tally.harness[]'`; cross-read against `awk '/^- \*\*`MUT-/{c++} END{print c}'` over this document | Plan census: **16 total** (6 HARNESS + 10 EVIDENCE). BC.A′ owns 3 (`MUT-REC-SILENT-DEFAULT`, `MUT-REC-STALL`, `MUT-REC-UTIL-STALL`, discharged iteration 48); **BC.B′ owns 13**, not 12. The document itself *defines* only 10 in bulleted form — the other 6 are specified inside AC prose — which is why a reader counting bullets gets a third number. Iteration 49 discharged 8 of BC.B′'s 13; this pass discharges the remaining 5 → **16 of 16, item-wide** | CONFIRMED — the "8 of 12" in the iteration-49 STATUS and queue row was a quantity quoted without a command in that session. Nobody was misled about *which* mutations remained (they were named individually and correctly); the arithmetic was simply unverifiable, and `8 + 5 = 13 > 12` was visible on the page |

One negative expectation is stated as an expectation, not a fact: **the recorder is expected to
refuse on ubuntu-latest** because `sysctl -n vm.loadavg` / `hw.ncpu` are BSD names absent from
Linux's sysctl namespace. This was **not measured on CI** from this rig. The design does not
depend on the expectation being right: the CI step asserts the refusal, so if ubuntu ever
satisfies the probes the step itself goes RED loudly (see D4) — the failure direction is loud
either way.

## Conflict Surface

| Existing surface or precedent | Collision / reuse question | Resolution |
|---|---|---|
| `scripts/bench_worldd.sh` `--smoke` + its hardcoded 10-name manifest | The only existing mode; CI depends on it (`ci.yml:88-89`) | **Unchanged byte-for-byte in behavior.** New modes (`--record-pair`, `--check-claims`) are added beside it; `usage()` documents all three (S7). The name manifest is not touched. All new surface is net-new (P28: no `--record`-family mode exists to migrate) |
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

### D1. The pair recorder — one bounded session, prebuilt binaries, interleaved legs, all-or-nothing

`scripts/bench_worldd.sh` gains a mode (the round-2 single-role `--record <role> [--dir]` is
**SUPERSEDED** — it is exactly the two-independent-recordings shape the round-2 reject refuted,
and it never shipped, P28):

```
./scripts/bench_worldd.sh --record-pair --variant <dir> --control <dir>
```

Both directory arguments are **mandatory** — there is no default, no single-role form, and no
way to record half a pair. The instrument is held constant across the pair by construction: one
script invocation, from one checkout, drives both worktrees (the control commit may predate the
recorder's existence — true for the very first pair — which is why the control tree is measured
*by* the variant tree's recorder, never by its own).

Execution order is fixed and load-bearing (AC1 and D4 depend on it; re-derived for branch A):

1. **Argument validation** — a missing/unknown flag exits 2 via `usage()`. Each of `--variant`
   and `--control` must exist and satisfy `git -C <dir> rev-parse --git-dir` — both `git`
   validations under the bounded runner with the utility deadline (round 5, P35) — else a named
   exit 1. No comparison between the two dirs happens here — the parent-edge check is git
   state (step 5), deliberately AFTER the rig probes, so D4's off-rig determinism is unchanged.
2. **Rig-global probes**, each of which must produce non-empty output or the mode **exits 1
   naming the exact probe that failed, having emitted NO fences, partial or otherwise**. Fixed
   order: `sysctl -n hw.ncpu`; `sysctl -n hw.model`; `sysctl -n vm.loadavg` (shape `{ a b c }`,
   P11); `ps -Ao pid=,ppid=,pcpu=,comm=`; `date -u`; `AILANG_BIN` set + executable +
   `--version` capture — the `--version` execution runs under the bounded runner with the 120 s
   probe deadline, because it executes an env-supplied binary and every such wait is bounded (a
   wedged pin binary is a stall like any other). **Round 5: the same rule covers the cheap
   probes** — each `sysctl`, the `ps` snapshot, and `date -u` run under the bounded runner with
   the utility deadline (`REC_UTIL_TIMEOUT_S=20`), because a wedged kernel interface or
   filesystem stalls the recorder exactly like a wedged binary; through round 4 these were
   unbounded (P35). These read kernel/env state and are
   cwd-independent (the sysctl trio, `ps`, and `date` read the rig; `AILANG_BIN` is
   env-derived, not tree-derived). They run FIRST so that an off-rig runner refuses
   deterministically at `sysctl -n hw.ncpu` — before any `go` invocation (which on a
   floor-raised tree could attempt a toolchain download, P24) and before any git check (D4's
   determinism). A recorder that emits an empty or partial block when a probe is unavailable
   would be a **silent fallback**; the null case is a named non-zero exit, checked in CI (D4).
3. **Session stamp**: `session_utc` := one `date -u` reading taken here, shared verbatim by
   both sections; `pair_nonce` := 32 hex chars from python3 `secrets.token_hex(16)` (P31) —
   both invocations under the bounded runner with the utility deadline (round 5). An
   empty or failed nonce read is a named refusal — never a default, never a zero nonce.
4. **Tree-dependent probes, per role against its own dir**: `go -C <dir> env GOVERSION GOOS
   GOARCH` (the `-C` form, verified P24), run under the bounded runner with the 120 s probe
   deadline, once per role. `GOVERSION` is **directory-sensitive** on this rig — P24 measured
   `GOTOOLCHAIN=auto` live and a `go.mod` floor selecting the toolchain — and the control tree
   is a *different commit*, which can carry a *different floor*; **OD-1 is exactly such a
   change**. Probing in the caller's cwd would write the variant's toolchain into BOTH
   sections and R4's `goversion` comparison would compare a value to itself — the CF-K-2
   bypass, unchanged from round 2 (`MUT-PROBE-CALLER-DIR` still demonstrates it). The recorder
   **records each tree's toolchain honestly and does not refuse a cross-toolchain pair** —
   REDing that pair is the checker's job (R4's toolchain limb), and AC6's real floor-split
   fixture depends on being able to record one.
5. **Git state, per role** (`git -C <dir>`; every `git` invocation in this step under the
   bounded runner with the utility deadline — round 5, P35: a wedged filesystem during a
   `git status` is a stall like any other, and it was one of the invocations the round-4 text
   left unbounded): `git status --porcelain` must be **empty** in
   BOTH trees — a measurement of a dirty tree is not attributable to a commit. `commit` :=
   `rev-parse HEAD` and `parent` := `rev-parse HEAD^` (P15) per tree, full SHAs. Then the
   session's own parent-edge check: **refuse, named, unless `control.commit ==
   variant.parent`** — a pair that is not variant-vs-parent cannot be recorded at all, so the
   in-file edge R4 later re-checks is git-true at record time by construction.
6. **Pair ID**: `pair_id` := SHA-256 (python3 `hashlib`, P12/P21) of the canonical byte string
   `bench-pair/2\n<session_utc>\n<variant_commit>\n<control_commit>\n<pair_nonce>\n`. Written
   verbatim into BOTH sections — **and so are the derivation's two commit inputs, as the
   session-level fields `pair_variant_commit` and `pair_control_commit`, identical verbatim in
   both sections** (round-4 fix, P32: the round-3 control section recorded no variant commit,
   so its ID was recomputable only after pairing — pair-scoped where the mutation predictions
   were section-scoped). The fields are **symmetric — both sections carry both — rather than
   the control-only field the reviewer literally proposed**, for a stated reason a reader
   would otherwise trip on: an asymmetric schema forces a role-conditional derivation rule
   (variant recomputes from `commit`+`parent`, control from `variant_commit`+`commit` — two
   rules where one suffices), whereas the pair's endpoints are session-level facts exactly
   like `session_utc` and `pair_nonce`, so both sections record them identically, one
   role-independent derivation rule covers both, R4b polices their cross-section identity,
   and R4c binds them to each section's own `commit`/`parent` (variant: `commit ==
   pair_variant_commit` AND `parent == pair_control_commit`; control: `commit ==
   pair_control_commit`) — the redundancy is not waste, it is the consistency check. The ID
   is therefore **recomputable by the checker from ANY ONE section's own fields, before any
   pairing** (R4c, section-local), so retyping it — or splicing sections from two different
   sessions, which necessarily carry different nonces and session stamps — breaks
   recomputation rather than merely a checksum; and it is unique per session by the nonce,
   which is what makes "this control belongs to exactly this variant" a file-grammar fact
   instead of a hope.
7. **Prebuild, both roles, before ANY measurement** (A2): `go -C <dir> test -c -o
   <session-tmpdir>/bench.<role> ./host/daemon/` under the bounded runner with the prebuild
   deadline, once per role, compile order control-then-variant (fixed, immaterial to
   measurement — nothing is being timed yet). Per role, record `prebuild_elapsed_s` and
   `binary_sha256` (python3 `hashlib` over the binary file). Compilation — the 128.85 s
   cache-cold population of P25 — is thereby OUTSIDE the measured window by construction; a
   leg never compiles anything (P26: the compiled binary is standalone and cwd-independent).
   **The `-C` form is not an equivalent spelling of a path argument** (P33): passing the
   control tree's package as an absolute path from the variant checkout — `go test -c -o
   "$BIN" <control-tree>/host/daemon/` — fails rc=1 with `directory … outside main module
   or its selected dependencies`, `FAIL … [setup failed]`, while the identical compile via
   `go -C <control-tree> test -c -o … ./host/daemon/` succeeds (2 s warm, P33). An executor
   who "simplifies" the `-C` away gets a setup failure that reads like a broken tree; the
   `-C` is load-bearing here for the same reason it is in step 4.
8. **Interleaved measurement legs, fixed order control/variant/variant/control** (A3) — the
   order is frozen policy (Design Freeze). Per leg: capture `legN_start_utc`, `legN_load`
   (`sysctl -n vm.loadavg`), `legN_competing` (the ps snapshot rule below); run
   `( cd <session-tmpdir> && ./bench.<role> -test.bench . -test.benchtime 200x -test.run '^$' )`
   under the bounded runner with the leg deadline — note the compiled-binary flag spellings
   `-test.bench`/`-test.benchtime`/`-test.run` (P26), and the recorded `invocation:` line is
   the command **verbatim as executed**, binary path and `-test.*` flags, never the `go test`
   spelling the session did not run; capture combined output verbatim; capture `legN_end_utc`,
   `legN_elapsed_s`, `legN_output_sha256`. The per-leg captures around the run — `date -u`
   (start/end), `sysctl -n vm.loadavg`, the `ps` snapshot, and the python3 hash of the leg
   output — each run under the bounded runner with the utility deadline (round 5, P35/P36).
   A non-zero leg rc → named exit 1 with the output
   shown; a deadline expiry → the named TIMEOUT refusal. **Either way the whole session emits
   ZERO fences** — a session that fails on leg 3 of 4 produces nothing, not a half-pair
   (staging below). Everything in the leg invocation other than the binary path is hardcoded
   in the script, so it cannot vary between legs or roles without a review-visible script edit.
9. **Emission, only after leg 4 completes**: all fields and all four raw outputs accumulate in
   a `mktemp -d` session directory during steps 3–8; only when every leg has succeeded and
   every hash is computed does the recorder assemble and print the two sections. There is no
   code path that emits a fence before step 9 — that is how all-or-nothing survives
   interleaving: not by cleanup-on-failure, but by construction.

**Bounded execution (Standing Rule 6).** `go test -timeout` bounds only the test binary's own
clock — not the compile step, not a wedged child process, not the `go` tool itself — so
~~every invocation of a non-trivial external binary above — the `go env` probes, both
prebuilds, all four legs, and the `$AILANG_BIN --version` capture — runs through a
`run_bounded` helper~~ **SUPERSEDED (round 5, `gemini-3-1-pro`): that enumeration was the
defect — everything NOT on it (`sysctl`, `ps`, `date`, every `git`, every `python3` helper)
was unbounded (P35), so a wedged filesystem or kernel lock during a `git` or `sysctl`
invocation could stall the recorder indefinitely while the text claimed every wait was
bounded.** The round-5 rule is **closed-world: EVERY external-binary invocation in the
session, with no exception, runs through a `run_bounded` helper mirrored verbatim from
`scripts/verify_ail.sh:61-74`** — the `sysctl` trio and the `ps`/`date` captures (step 2 and
per leg), both step-1 `git rev-parse --git-dir` validations, every step-5
`git status`/`git rev-parse`, every `python3` helper (the step-3 nonce read and every
`hashlib` computation in steps 6–9), the per-role `go env` probes, both prebuilds, all four
legs, and the `$AILANG_BIN --version` capture (the V26 helper, re-read first-party
this revision, P30: python3 `start_new_session` puts the child in its own process group; on
expiry the WHOLE group gets SIGKILL — killing only the direct child would orphan a spawned
grandchild — a named `TIMEOUT after Ns` goes to stderr, and the exit is 124). This is the
**second** unbounded wait to grow inside this document's own remedy for unbounded waits —
round 2's draft left `$AILANG_BIN --version` unbounded and was self-caught; rounds 3–4 left
the utility class unbounded and were quorum-caught — so the recurrence is recorded in *What
these gates CANNOT fail* and the closed-world rule is frozen (Design Freeze), with
`MUT-REC-UTIL-STALL` (AC1) as the discriminator that a newly-covered invocation actually
refuses instead of stalling. **Mirrored, not
sourced or extracted**: sourcing `verify_ail.sh` would execute the gate on load, and extracting
a shared lib would edit `verify_ail.sh` — which P17 and AC5 pin as untouched — to move 13
lines. The copy carries a provenance comment naming `verify_ail.sh:61-74`/V26, and divergence
between mirror and original is a review-visible change (Design Freeze).

**The deadline structure is restructured for branch A, and this discharges `4f/CF-M-2`** ("P25's
deadline arithmetic is designer-measured, not controller-re-run cold"): the controller re-ran
the arithmetic cold this iteration (P26/P27), and what the re-run showed is not a different
number but a different SHAPE — the round-2 `REC_BENCH_TIMEOUT_S=600` bounded ONE window
containing TWO populations with different physics: a compile that legitimately took 128.85 s
cache-cold and can take longer still when `GOTOOLCHAIN=auto` fetches a toolchain over the
network (P24), and a measurement run that, once the binary is prebuilt, sampled at 6–8 s per
full leg on a loaded rig (P27). A2 separates the populations, so the deadlines separate too.
**Four** hardcoded deadlines (round 5 adds the utility class), **policy, not env knobs** (the
D3 precedent); each is per-invocation of `run_bounded`. ~~The session has no single umbrella
knob because every individual wait is bounded, and the session's worst case is therefore the
sum of its parts by construction.~~ **SUPERSEDED (round 5)** — as written through round 4 that
sentence was FALSE: the `sysctl`/`ps`/`python3`/`git` invocations were bounded by nothing
(P35), so "every individual wait" quietly meant "some waits". What holds of the FIXED design:
every external-binary invocation runs through `run_bounded` (the closed-world enumeration
above), so the session's worst case is the sum of its per-invocation deadlines — ≈48 min
fully wedged (2 × 600 s prebuilds + 4 × 120 s legs + 3 × 120 s probe-class invocations +
~40 utility-class invocations × 20 s ≈ 800 s), against tens of seconds for a real session
(P27/P34) — **with one named residual**: `run_bounded` is itself a python3 wrapper, and the
wrapper's own interpreter startup is the one wait the helper cannot police, because it IS the
policing mechanism; that residual is inherited verbatim from `verify_ail.sh` (V26/P30) and is
stated in *What these gates CANNOT fail* rather than hidden behind "by construction":

- **`REC_PROBE_TIMEOUT_S=120`** on the probe-stage binaries (`go env` per role and the
  `$AILANG_BIN --version` capture) — unchanged from round 2; same species and value as
  `GATE_LEG_TIMEOUT_S` (P23/P30). A first-ever probe of a floor-raised tree may fetch a
  toolchain (P24) — bounded like every other wait. (The AC6 fixture avoids the fetch: both
  floors it uses are already cached, P24.)
- **`REC_UTIL_TIMEOUT_S=20`** (round 5) on every cheap utility invocation: the `sysctl` trio,
  every `ps` snapshot, every `date -u` reading, every `git` invocation (steps 1 and 5), and
  every `python3` helper — the nonce read and every `hashlib` computation, including hashing
  the ~17 MB prebuilt binaries. Justification: every member of the class was measured
  first-party at ≤ 0.050 s on a loaded rig (P36; the worst member is the 17 MB hash), so 20 s
  is an **explicitly-labelled 400× engineering margin** over the worst sampled member — a
  margin over a sample, not a measured tail. The 120 s probe constant is deliberately NOT
  reused: nothing in this class legitimately takes minutes (no network, no toolchain fetch,
  no compile — kernel reads, local git metadata, local hashing), so a member exceeding 20 s
  is a wedged rig and refusing loudly is the safe direction, while a 120 s ceiling would be
  unjustifiable padding that quadruples the wedged-session worst case for no measured reason.
  The class has ~40 members per session, so even fully wedged it adds ≈13 min, not hours; the
  refusal names the exact invocation that expired.
- **`REC_PREBUILD_TIMEOUT_S=600`** on each `go test -c` prebuild. Justification: the worst
  *measured legitimate* compile on this rig is **128.85 s** (cache-cold, rig-contended, P25 —
  that measurement stands; only its round-2 consequence is superseded); 600 s is ~4.7× that,
  an **explicitly-labelled engineering margin** covering a possible in-band toolchain fetch
  whose duration was NOT measured (network-dependent; pre-caching is the remedy either way),
  while still bounding a wedged compile inside ten minutes.
- **`REC_LEG_TIMEOUT_S=120`** on each of the four measurement legs. Justification: the leg
  population is compile-free by construction (step 7), and its only measurements are P27's
  five-run loaded-rig SAMPLE — 6, 6, 6, 7, 8 s — and P34's complete four-leg interleaved
  session (7, 7, 6, 7 s legs; 27 s total measured window, same species, same loaded rig),
  which together establish an order of magnitude
  (seconds), **not a bound**; the iter-46 spine forbids dressing either sample up as one. 120 s
  is therefore an **explicitly-labelled engineering margin**: ~15× the worst sampled leg, the
  same species and value as `GATE_LEG_TIMEOUT_S` (P30), small enough that four wedged legs
  cannot exceed eight minutes, and honest about what it is — a margin over a sample, not a
  measured tail. If a legitimate leg ever exceeds it, the refusal is loud and the constant is
  re-derived in review (the safe direction; see *What these gates CANNOT fail*).
- ~~`REC_BENCH_TIMEOUT_S=600`~~ — **SUPERSEDED** (round-2 shape; retired with the single-window
  measurement it bounded).

A deadline expiry at ANY stage is a **named non-zero refusal that emits ZERO fences** — the
same class as a probe refusal (e.g. `bench leg 3/4 TIMEOUT (>120s) — no conditions emitted`).
A stall can never produce a partial pair, a half-session, or an `unknown`-filled block. Proven
by `MUT-REC-STALL` (AC1), whose kill must still reach a grandchild.

**Competing-process capture, per leg** (the P7/P8 hard case): from one
`ps -Ao pid=,ppid=,pcpu=,comm=` snapshot taken at leg start, keep rows with `pcpu ≥ 25.0`,
sorted descending, capped at 8, one line each:
`legN_competing: pcpu=<x> pid=<p> comm=<verbatim path> parent=<comm of ppid>`. The generic
temp-path `go-build` binary is handled by **recording, not guessing**: the verbatim path plus
the parent's comm (P8 measured it resolving to `go`) is what is mechanically observable; naming
the owning *mission* is not, and the recorder does not pretend otherwise. If the parent has
already exited, `parent=?` is written — an explicit recorded unknown, distinguishable from an
absent field. If no process meets the threshold, the line is the explicit
`legN_competing: none>=25%` — an explicit negative, never an empty field. Four snapshots per
session (one per leg) replace round 2's before/after pair: the leg is the unit of exposure now.

**How the two legs per role are reported — no averaging, ever.** Each role's section is
followed by BOTH of its legs' raw outputs, verbatim, in leg order. The recorder computes **no
summary statistic** — no mean of two legs, no "representative" leg, no reconciled number. This
is the iter-46 spine applied: two legs are a sample of size two, and collapsing them into one
number destroys the only within-session noise reading the pair produces for free. What a reader
is entitled to conclude: the leg-adjacent comparisons (leg 1 control vs leg 2 variant; leg 3
variant vs leg 4 control) are the tightest-window A/B signals this rig can produce; the
within-role spread (leg 2 vs leg 3, leg 1 vs leg 4) is the session's own noise floor, read
against P29's ~1.4× constant-load spread; and an effect smaller than that within-role spread is
**not resolved by this pair** — the file shows the reader that fact instead of hiding it in an
average. Prose in `BASELINE.md`'s header states exactly this reading rule (D5).

**Output**: one paste-ready emission — BOTH sections, each followed by its two raw leg blocks —
printed to stdout AND written to a `mktemp` file whose path is printed; the human appends it
with `cat "$file" >> bench/BASELINE.md`, retyping nothing. Schema `bench-conditions/2` (the
version moves because the block shape changed: pair/session/leg fields added,
`load_before`/`load_after`/`elapsed_s`/`output_sha256` replaced by per-leg fields; no
`bench-conditions/1` block exists anywhere to migrate, P31 — the checker recognises `/2` only).
Per role, the emission is:

````
```bench-conditions
schema: bench-conditions/2
role: control                      # or: variant
pair_id: <64 hex — recomputable, R4>
pair_nonce: <32 hex>
session_utc: 2026-08-04T05:29:25Z  # identical in both sections — one session stamp
pair_variant_commit: <full sha>    # identical in both sections — the pair's variant endpoint (round 4, P32)
pair_control_commit: <full sha>    # identical in both sections — the pair's control endpoint; R4c binds both to this section's own commit/parent
commit: <full sha>
parent: <full sha>                 # variant.parent must equal control.commit (step 5, R4)
tree: clean
goversion: go1.26.4
goos_goarch: darwin/arm64
ncpu: 16
hw_model: Mac16,9
ailang_pin: AILANG v0.30.0 (e37b370) via $AILANG_BIN
prebuild_elapsed_s: <n>
binary_sha256: <hex>
invocation: <session-tmpdir>/bench.control -test.bench . -test.benchtime 200x -test.run '^$'
leg_order: control,variant,variant,control   # identical in both sections — frozen policy
leg1_seq: 1/4                      # this role's first leg position (variant: 2/4)
leg1_start_utc: <utc>
leg1_end_utc: <utc>
leg1_elapsed_s: <n>
leg1_load: 2.96 3.20 3.06
leg1_competing: pcpu=100.0 pid=28606 comm=/var/folders/.../exe/solution parent=go
leg1_output_sha256: <hex>
leg2_seq: 4/4                      # this role's second leg position (variant: 3/4)
leg2_start_utc: <utc>
leg2_end_utc: <utc>
leg2_elapsed_s: <n>
leg2_load: <a b c>
leg2_competing: none>=25%
leg2_output_sha256: <hex>
conditions_sha256: <hex>
```
```text
<this role's FIRST leg raw output, verbatim>
```
```text
<this role's SECOND leg raw output, verbatim>
```
````

**Integrity fields**, all computed via python3 `hashlib` (P12, P21 — no shasum/sha256sum CLI
divergence): each `legN_output_sha256` = SHA-256 of that leg's raw output (the lines between
its ` ```text ` fences, joined with `\n`, trailing `\n`); `binary_sha256` = SHA-256 of the
prebuilt binary's bytes; `conditions_sha256` = SHA-256 of every block line above the
`conditions_sha256:` line itself, each terminated with `\n`; and `pair_id` is recomputable from
each section's OWN `session_utc` + `pair_variant_commit` + `pair_control_commit` + `pair_nonce`
— section-locally, no counterpart consulted (step 6, checked by R4c before any pairing). These are
**tamper-evidence, not cryptographic provenance** (see *What these gates CANNOT fail*): their
job is to make retyping, truncation, splicing, and post-hoc number-editing RED, and — applying
the iteration-41 lesson that *a gate whose failure message hands you the silencing value is not
a gate* — **no failure message anywhere in this design prints an expected digest, ID, or
nonce**. Re-greening requires re-running the recorder, not pasting a value out of the error
text.

### D2. The policy — (ii)-as-A/B, written into the file and cashed out mechanically

The charter's stated preference is adopted and its reasoning recorded in `BASELINE.md`'s header:
**an A/B against the parent commit, captured as one interleaved single-session pair, is
MANDATORY for any cost claim.** The rejected alternative — refuse-to-record above a load
threshold — is out of scope with reasons (see *Out of scope*): `BASELINE.md:7` already calls
noise-gating a shared runner dishonest, and a threshold gate assumes an idle rig that a
two-mission rig does not guarantee.

> **RATIFIED WORDING (`4f/OD-8`, attended 2026-08-04 — binding on the header text BC.B′ writes
> into `BASELINE.md`).** "Mandatory" is a rule about **admissibility of evidence**, never about
> the **validity of a claim**. The header must say, in substance: *a pair recorded this way is
> mechanically complete, contemporaneous, tamper-evident **evidence for** a cost claim; the
> checker enforces that evidence is admissible, and it determines nothing about whether the
> observed delta supports the claim.* Whether a delta is real — repetitions, noise bounds — is
> **out of scope here by ratification** and is an owned requirement of `w-agent-floor-m4`'s
> experimental design (N≥3 paired runs + noise handling). Any phrasing of the form "the only
> mechanically valid form of a cost claim" is **forbidden in the shipped file**; it is the exact
> over-claim `gpt5-6-sol` blocked on across two quorum rounds and the human corrected.

> ~~"…whereas the A/B is correct under any load (it is what actually caught the 6.06×
> artefact)."~~ **SUPERSEDED (round 2 → branch A).** This was the sentence the round-2 reject
> refuted: an *independently recorded* pair is not correct under any load — R4 compared neither
> load nor temporal pairing, so an idle variant could pass against a loaded stale control.
> ~~What is true, and all this design now claims: an interleaved single-session pair BOUNDS
> load skew — both roles experience the same rig episode, so the confound is bounded by the
> session's own duration — but does not ELIMINATE it.~~ **SUPERSEDED (round 5,
> `gpt5-6-sol`): dimensionally wrong — the session's duration bounds the TEMPORAL SEPARATION
> of the legs, not the MAGNITUDE of load divergence; a competitor can spike inside the
> session, hit only the variant legs, and nothing in this design says otherwise.** What is
> true, and all this design now claims: **the interleaving makes the two roles' legs adjacent
> in time and alternating in order (C/V/V/C), so a MONOTONIC drift in rig conditions affects
> both roles roughly equally and largely cancels in the comparison; a load episode shorter
> than the session, or one correlated with leg boundaries, does NOT cancel and is NOT
> detected.** The per-leg load and process snapshots give a HUMAN reader visible evidence of
> divergence to judge; the CHECKER mechanically enforces nothing about load — by ratified
> design (the excluded third limb; P29 shows why no threshold can adjudicate it). Iteration
> 39's rescue worked *because* its legs happened to be minutes apart under the same load —
> branch A makes the adjacency structural instead of accidental, which narrows the window for
> divergence without bounding its size.

"Mandatory" cashed out:

- **The artifact that carries a claim** is a **single-session pair** of recorder emissions in
  `bench/BASELINE.md`: one `role: variant` section and one `role: control` section carrying the
  SAME `pair_id`, with `control.commit == variant.parent`, each with both of its interleaved
  legs. New numbers cannot enter the file's evidence in any other form: **every** non-legacy
  raw benchmark block must be bound to a conditions block, and **every** variant must have its
  same-session control. Appending a measurement without its pair — or with a control borrowed
  from another session — is structurally RED, which is the A/B mandate expressed as file
  grammar rather than as prose exhortation.
- **What checks it**: `./scripts/bench_worldd.sh --check-claims` (D3).
- **What goes RED, where**: the `worldd bench claim-structure gate` step in CI's `go-verify` job
  (every push to `dev` and every PR), and the same command locally. F2/P3 establish that today's
  CI evaluates no bench structure at all; this adds the missing leg without touching the smoke
  gate.
- **The A/B procedure** (added to `BASELINE.md`'s header, executed-verbatim per S7 — by the
  controller, since the executor sandbox denies the loopback binds five benchmarks need, the
  standing `<CONTROLLER-MEASURED>` precedent): `git worktree add ../.bench-control HEAD^`;
  `./scripts/bench_worldd.sh --record-pair --variant . --control ../.bench-control` (ONE
  invocation, one session, four legs); append the emitted file; run `--check-claims`;
  `git worktree remove ../.bench-control`; commit. There is no two-invocation form to fall
  back to — the recorder has no single-role mode (D1).
  **PATH CORRECTED (iteration 48):** ~~`/tmp/bench-control`~~ — the shared mission-control skill
  forbids worktrees under `/tmp`, because a `/tmp`-rooted checkout fails CWD-resolving tests for
  its *location* and CI cannot reproduce the red. The procedure written into `BASELINE.md` is
  executed verbatim by future readers (S7), so it must not teach the forbidden path. Use a
  **sibling of the repo** here and in AC6.

Prose in the header additionally states the policy's meaning for readers: a delta between
invocations recorded under different conditions is **indicative only**; ~~the pair is the only
load-independent signal~~ **SUPERSEDED (round 5): no signal here is load-independent** — the
pair is the only form whose conditions are machine-recorded and whose legs are contemporaneous
by construction; what its delta means under those recorded conditions remains the reader's
judgment (*What these gates CANNOT fail*). (Prose cannot be machine-policed; the grammar above
polices the evidence. The residual is stated honestly below.)

### D3. The structure checker — in-file, history-free, null-case-loud

`./scripts/bench_worldd.sh --check-claims` parses **`bench/BASELINE.md` only** (hardcoded path,
no path argument, no env knob — the `verify_ail.sh` "policy is not env-overridable" precedent).
Implementation: bash + embedded python3 (P12 precedent), python3 absence fails loudly.

**Evaluation order (round 4 — stated because it was previously stated nowhere, P32, and it is
load-bearing for the mutation predictions):** the checker evaluates **R6** (file-level null
case) first; then the per-block rules **R1 → R2 → R3**; then **R4c, which is SECTION-LOCAL and
runs BEFORE any pairing** — each conditions block is checked against its OWN fields with no
counterpart consulted; then **grouping by `pair_id`, which is itself a rule, not plumbing** —
a block that cannot be placed into a well-formed pair is RED via R4d's named unpairable case,
never skipped; then the pair-level limbs **R4a/R4b/R4e/R4f** on each well-formed pair; then
**R5**. The checker does not stop at the first failure: every rule that can evaluate is
evaluated, every failure is reported with its named message, and the exit is non-zero if any
rule REDded (this is what makes dual-fire predictions like `MUT-CLAIM-NONPARENT`'s
observable). A pair-level limb that cannot evaluate because its pair never formed is **not a
pass and not a silent skip** — the block is already RED via R4d, and "cannot evaluate" is
never "nothing to report" (the same null-case-loud discipline R6 applies to an empty file).
Rules, each with a named failure:

- **R1 — block validity.** Every ` ```bench-conditions ` block must declare
  `schema: bench-conditions/2` (the only schema that ever shipped, P31 — any other value is
  RED by name), contain every required key from D1's schema, non-empty, `tree: clean`, a
  `conditions_sha256` that recomputes, and an `invocation:` line matching the
  **prebuilt-binary grammar** — a binary path plus `-test.bench`/`-test.benchtime`/`-test.run`
  flags, never a `go … test` spelling (`invocation is not a prebuilt-binary invocation` — the
  in-file tooth against folding the compile back into the measured window, and `binary_sha256`
  must be present for the same reason). Mismatch messages name the block's line number and
  **print no expected digest**.
- **R2 — no orphan numbers.** Every fenced block containing a line matching
  `^Benchmark[A-Za-z_/]+(-[0-9]+)?\s` (the P13-verified detector, immune to the `ns/op`-in-prose
  false positive) must belong to a valid conditions block's **raw-block group** — each
  conditions block owns exactly the TWO ` ```text ` blocks that follow it (leg order), the
  first opening within 5 lines of the conditions block's close — or carry a
  `<!-- legacy-unconditioned: pre-4f … -->` marker. A conditions block followed by fewer or
  more than two raw blocks is RED (`conditions block does not own exactly two leg outputs`).
- **R3 — output binding, per leg.** For each conditions block, raw block 1 must hash to
  `leg1_output_sha256` and raw block 2 to `leg2_output_sha256`. Editing one digit of one
  recorded number in either leg is RED; the message names block and leg, never the expected
  hash.
- **R4 — the A/B mandate, pair-identity edition.** Pair-level rules, each with its own named
  failure, none printing an expected value:
  - **R4a — parent edge.** Every `role: variant` block must have a `role: control` block with
    `control.commit == variant.parent` (`A/B pair is not variant-vs-parent`). Checked
    **in-file** (recorder-emitted fields), not via git history — deliberate and unchanged: CI
    checkouts are shallow (P4), squash-merged SHAs may become unreachable, and the edge is
    git-true at record time by construction (D1 step 5, P15). Declared in *What these gates
    CANNOT fail*. (Round 4 note: R4a is now IMPLIED by R4c's role-binding plus R4b's shared
    endpoint fields; it is retained as the direct pair-level statement of the mandate, with
    its own named message — redundant limbs that fail independently are defense in depth,
    and `MUT-CLAIM-NONPARENT` still targets this one directly.)
  - **R4b — pair identity.** The variant and its control must carry the **same `pair_id`**,
    the same `session_utc`, the same `pair_nonce`, the same `pair_variant_commit`, the same
    `pair_control_commit`, and the same `leg_order` (`pair identity mismatch — sections are
    not from one session`). ~~This is what kills the round-2 reject's scenario: a stale
    control recorded in an earlier session necessarily carries a different pair ID, so it
    cannot be attached to a new variant.~~ **SUPERSEDED (round 4 — mechanism attribution
    corrected, P32):** a stale control with a DIFFERENT pair ID never reaches R4b, because no
    pair forms — R4d's unpairable case REDs it first. R4b's own job is the retype path:
    sections FORCED under one `pair_id` whose other identity fields betray two sessions.
  - **R4c — pair-ID recomputation and endpoint binding, SECTION-LOCAL (round 4, P32: runs
    BEFORE any pairing — see the evaluation order).** Two clauses, each evaluable on a lone
    block with no counterpart consulted: **(i) derivation** — the section's `pair_id` must
    recompute as SHA-256 of
    `bench-pair/2\n<session_utc>\n<pair_variant_commit>\n<pair_control_commit>\n<pair_nonce>\n`
    (D1 step 6) from that SECTION's own recorded fields (`pair_id does not derive from
    recorded session fields`); **(ii) role-binding** — the endpoint fields must bind the
    section's own commits: `role: variant` requires `commit == pair_variant_commit` AND
    `parent == pair_control_commit`; `role: control` requires `commit == pair_control_commit`
    (`pair endpoints do not bind this section's own commits`). A retyped ID, or a section
    whose commit/session fields were edited around a kept ID, breaks one of the two clauses —
    not just a checksum — and it breaks them **on that section alone**, which is what makes
    `MUT-PAIR-ID-SPLIT`'s prediction derivable from the mechanism rather than asserted.
  - **R4d — pair cardinality: control reuse AND the unpairable null case.** File-wide, every
    `pair_id` value must appear on **exactly two** conditions blocks, exactly one
    `role: variant` and exactly one `role: control` (`control reused across pairs` /
    `duplicated pair member`). Two variants citing one control, a duplicated control section,
    or three blocks sharing an ID are all RED. And — **stated as a rule, not a remark
    (round 4, P32)** — a conditions block whose `pair_id` appears on NO other block is RED
    with its own named message (`unpairable conditions block — no counterpart shares its
    pair_id`), **never skipped, never a vacuous pass**: an implementer who reads "the checker
    cannot form this block's pair" as "nothing to evaluate" has built the silent pass this
    limb exists to forbid. Combined with R4b/R4c, "reuse" is a pure file-grammar fact — no
    git history, no timestamps trusted beyond what the hashes bind.
  - **R4e — toolchain/hardware identity (the CF-K-2 tooth, unchanged).** The pair's
    `goversion`, `goos_goarch`, `ncpu`, and `hw_model` must be **identical**
    (`toolchain mismatch inside claimed A/B pair`). A cross-toolchain "A/B" is not a
    comparison; the recorder records each tree honestly (D1 step 4) and this limb REDs it.
  - **R4f — interleave structure.** Both sections' `leg_order` must be the frozen
    `control,variant,variant,control`; the four legs' `legN_seq` positions must be exactly
    {control: 1/4, 4/4; variant: 2/4, 3/4}; and the recorded `legN_start_utc` values must be
    strictly increasing in seq order with each leg's `end_utc ≥ start_utc`
    (`legs not interleaved — control legs are not outermost`). This REDs a sequentially
    captured pair (control/control/variant/variant) that records its timestamps honestly; a
    harness that FORGES timestamps is fraud, bounded by review like every hash here (see
    *What these gates CANNOT fail*).
- **R5 — enumerated legacy.** Exactly **3** legacy markers exist (the P13 blocks at lines
  222/245/264 get them in BC.B′). `!= 3` is RED in both directions: a 4th marker (the "just mark
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
2. **`worldd bench recorder refuses off-rig`**: runs
   `--record-pair --variant . --control .` and **asserts it exits non-zero AND stderr carries
   the SPECIFIC first-probe marker** — the exact `probe FAILED: sysctl -n hw.ncpu` text, not a
   generic `probe FAILED` substring (explicit `if`-negation, not `!`-with-`set -e` ambiguity).
   **This is `4f/CF-M-1`, folded in and discharged here**: `gemini-3-1-pro`'s round-2 catch was
   that a generic marker lets a mutation that bypasses only the sysctl probes still exit
   non-zero later (e.g. on CI's unset `AILANG_BIN`) and masquerade as a green step — the
   specific marker makes the assertion discriminate WHERE the refusal fired, not merely THAT
   something failed. On ubuntu the BSD sysctl probes are expected absent, so this exercises the
   D1 refusal path — the anti-silent-fallback property — as a standing CI gate, before any
   prebuild or leg would run (seconds, not minutes). If a future runner ever satisfies every
   probe, the step REDs loudly and a human looks; the expectation failing is itself loud, never
   silent (see the note under the Premise Verification Log).

Determinism re-derived for the branch-A D1 order: **rig-global probes run before the session
stamp, tree-dependent probes, git state, prebuilds, and legs**, so on ubuntu the refusal fires
at the FIRST probe in the fixed order — `sysctl -n hw.ncpu` — before any `go` invocation (no
toolchain fetch can start on CI) and before any git check (earlier steps' downloaded tarballs
cannot affect which error fires). Passing `.` for both roles is deliberate and safe: D1 step 1
checks only that each dir exists and is a git dir (both true of the checkout root — and a
shallow CI checkout could not resolve `HEAD^` anyway, P4), no dir comparison happens at
validation, and the same-dir "pair" can never survive to a recording because on any rig that
passes the probes, step 5's parent-edge check REDs `control.commit == variant.commit ≠
variant.parent` by name. The refusal therefore lands at `hw.ncpu`, deterministically, and the
step's grep is pinned to exactly that marker.

### D5. `BASELINE.md` content changes (BC.B′)

- Header gains: the policy paragraph and A/B procedure (D2), the statement that the
  machine-emitted conditions block is the only valid form for new measurement conditions, the
  **leg reading rule** (D1: four legs, no averaging; leg-adjacent comparisons are the signal;
  within-role spread is the session's own noise floor; effects inside that spread are not
  resolved by the pair), and the honest statement verbatim — round-5 corrected: the round-3
  form, ~~*an interleaved single-session pair bounds load skew by the session's own duration;
  it does not eliminate it*~~, is **SUPERSEDED as dimensionally wrong** (duration bounds
  temporal separation, not divergence magnitude), and what the header carries instead is:
  *interleaving makes the roles' legs adjacent in time and alternating in order, so a
  monotonic drift in rig conditions largely cancels in the comparison; a load episode shorter
  than the session, or one correlated with leg boundaries, neither cancels nor is detected —
  the per-leg load and process snapshots are evidence for the reader's judgment, and nothing
  mechanical gates on load.*
- The three historical raw blocks (P13) each gain their `legacy-unconditioned` marker line.
- The amortisation section's existing note gains one sentence naming the **amortisation
  re-derivation follow-up** (see *Out of scope / Deferred*) instead of the current open-ended
  "when queue item 4f lands a load gate" — this item deliberately lands **no load gate**, and the
  pointer must not imply one. **WORDING CORRECTED (iteration 48, P38):** the follow-up is
  *"pending the first clean-toolchain pair recorded through `--record-pair`"* — ~~*"pending
  OD-1"*~~, which is discharged and would be a dangling reference the day this lands.
- One real recorder **pair** (variant at the BC.B′ merge candidate, control at its parent —
  ONE `--record-pair` session, four legs, both sections plus four raw blocks), produced by the
  controller, is appended under a section explicitly labelled: *mechanism acceptance run —
  conditioned on the recorded `GOVERSION`, superseded by the first amortisation re-derivation;
  NOT a milestone performance baseline* (**label corrected, iteration 48**: it must name the
  toolchain the run actually recorded rather than the literal `go1.26.4`, and must not cite a
  discharged decision). This gives deliverable (i) a real in-file instance, gives R1–R4 (all
  limbs, including recomputation and interleave) a live pair to validate forever, and gives the
  named mutations their targets — while the label keeps it from being read as a banked
  performance baseline.

## Acceptance Criteria

Every AC names a concrete mutation that makes it RED and states what makes that RED hard to
silence by pasting a value back.

### AC1 — the session refuses loudly; it never emits a partial pair

Owned by **BC.A′**. With any single probe unavailable, `--record-pair` exits non-zero, names
the probe, and emits **zero** `bench-conditions` fences (checked by grepping the captured
output). With a dirty tree in EITHER role's dir it refuses naming the dir and the git state.
With `control.commit != variant.parent` it refuses naming the parent-edge check (D1 step 5).
**With any leg stalled past the leg deadline — including leg 3 of 4, after two legs have
already succeeded — it exits non-zero via the named TIMEOUT refusal, emits zero fences (no
half-pair, no single-section emission), and leaves no orphaned child process running** (the
process-group kill; emission-only-after-leg-4, D1 step 9). **The same refusal contract covers
the utility class (round 5, P35): a wedged `git`/`sysctl`/`ps`/`python3`/`date` invocation
expires its 20 s deadline, exits non-zero with the named TIMEOUT naming that invocation, and
emits zero fences.** Probe checks precede the session
stamp precede git checks precede prebuilds precede legs (D1 order).

Named RED mutation: **`MUT-REC-SILENT-DEFAULT`** edits `scripts/bench_worldd.sh` (**HARNESS**)
to make a failed probe yield the literal `unknown` instead of refusing. Detection: the AC1
verification run (a PATH-shadowed `sysctl` stub returning empty, rc=0) now **emits fences
instead of refusing** — RED at BC.A′ review; and permanently RED in CI once BC.B′'s off-rig
step lands, because `--record-pair` on ubuntu would stop failing at the exact `probe FAILED:
sysctl -n hw.ncpu` marker the step greps (CF-M-1). Hard to silence: there is no value to paste
— the gate demands live probe output, and the CI step's assertion is on the *refusal marker*,
not on any literal the mutation can emit.

Named RED mutation: **`MUT-REC-STALL`** edits `scripts/bench_worldd.sh` (**HARNESS**) in ONE
edit: the **leg-3 invocation** is replaced by a stub that writes its own child's PID to a temp
file, spawns that child sleeping 3600 s, and then sleeps itself — a wedge with a grandchild,
the exact Standing-Rule-6 stall shape, placed mid-session so the refusal is tested AFTER two
legs have succeeded — and the same edit lowers `REC_LEG_TIMEOUT_S` from 120 to 5. (Waiting out
the real deadline proves nothing extra about the kill path; the constant is deliberately not an
env knob, so the lowered value is part of the sha256-restored mutation itself, not a runtime
override. That the shipped constant equals 120 is asserted by reading the landed script at
review.) Detection: the run exits non-zero with the named `TIMEOUT` on stderr naming leg 3, the
captured output contains **zero** `bench-conditions` fences — the two completed legs' data died
in staging, as designed — and after exit `kill -0 <recorded child pid>` fails: the group kill
reached the *grandchild*, not just the direct child. Hard to silence: the assertion is on
refusal behavior and on process death; there is no value to paste.

Named RED mutation: **`MUT-REC-UTIL-STALL`** (round 5) edits `scripts/bench_worldd.sh`
(**HARNESS**) in ONE edit: the **step-5 variant-tree `git status --porcelain` invocation** —
a member of the utility class that rounds 3–4 left unbounded (P35) — is replaced by the same
wedge-with-a-grandchild stub as `MUT-REC-STALL` (child PID recorded to a temp file, grandchild
sleeping 3600 s, child then sleeping itself), and the same edit lowers `REC_UTIL_TIMEOUT_S`
from 20 to 5 (the same not-an-env-knob discipline: the lowered constant is part of the
sha256-restored mutation, and the shipped value 20 is asserted by reading the landed script at
review). Detection: the run exits non-zero with the named `TIMEOUT` on stderr naming the
git-status invocation, the captured output contains **zero** `bench-conditions` fences, and
after exit `kill -0 <recorded child pid>` fails — the group kill reached the grandchild. This
is the discriminator for the round-5 coverage widening itself: it proves a NEWLY-covered
invocation refuses instead of stalling, so the closed-world enumeration is behavior, not
prose. Hard to silence: the assertion is on refusal behavior and process death; there is no
value to paste. sha256-restored.

### AC2 — the emission is complete, and its integrity fields recompute independently

Owned by **BC.A′**. A real `--record-pair` emission on the dev rig (controller-performed —
sandbox loopback denial, the standing precedent) contains BOTH sections and all FOUR raw leg
blocks, every D1 key non-empty, `legN_competing` lines carrying verbatim paths and `parent=`
attribution for temp-path binaries (P8), identical `pair_id`/`pair_nonce`/`session_utc`/
`pair_variant_commit`/`pair_control_commit`/`leg_order` across the two sections, `legN_seq`
at exactly {1/4, 4/4} (control) and
{2/4, 3/4} (variant) — and every integrity field recomputes under an **independent** python3
reimplementation (not the script's own functions): all four `legN_output_sha256`, both
`binary_sha256`, both `conditions_sha256`, and the `pair_id` derivation itself (D1 step 6) —
the derivation recomputed **twice, once from EACH section's own fields alone** (round 4:
section-locality is the fixed property, so the acceptance evidence exercises it per section,
never via the counterpart).

Named RED mutation: **`MUT-EDIT-RAW-NUMBER`** edits one digit of one `p95_ms` value in one
emitted leg block (**EVIDENCE**). The independent recompute mismatches that leg's
`legN_output_sha256` — and after BC.B′, `--check-claims` R3 is RED naming block and leg. Hard
to silence: the failure message prints no expected hash (the iter-41 lesson, applied); the
legitimate green path is re-running the recorder.

### AC3 — the A/B mandate has structural teeth

Owned by **BC.B′**. `--check-claims` is green on the post-BC.B′ file (3 legacy markers + 1 real
single-session pair) and RED under each of:

- **`MUT-CLAIM-NONPARENT`** (**EVIDENCE**): edit the acceptance variant block's `parent:` to the
  grandparent SHA **and recompute `conditions_sha256`** — isolating R4a from R1, proving the
  parent-edge check itself discriminates. (R4c also fires — **round-4 corrected mechanism,
  P32**: ~~the edited commit-family no longer derives the kept `pair_id`~~ **SUPERSEDED** —
  the derivation clause reads the untouched `pair_variant_commit`/`pair_control_commit`
  fields and still recomputes; what fires is R4c's **role-binding clause**, because the
  edited `parent` no longer equals the section's own `pair_control_commit` (`pair endpoints
  do not bind this section's own commits`), section-locally, before pairing. The predicted
  output is BOTH named failures — R4c's binding RED and R4a's pair-level edge RED, in that
  evaluation order; a run that shows only one is itself a finding.) Expected messages name
  the pair, never a SHA to paste.
- **`MUT-CLAIM-TOOLCHAIN-SPLIT`** (**EVIDENCE**): edit the control block's `goversion` to
  `go1.25.6`, recompute its checksum → R4e REDs with `toolchain mismatch inside claimed A/B
  pair`. This is CF-K-2's tooth: numbers from two toolchains can never satisfy the claim
  grammar. (Retained unchanged from round 2, as the ratified scope requires.)
- **`MUT-CLAIM-ORPHAN-NUMBERS`** (**EVIDENCE**): append a new ` ```text ` block containing one
  well-formed `Benchmark…` line with no conditions block → R2 REDs. This is the "developer
  writes the words 'A/B done'" evasion: words attach no conditions block, so any *numbers*
  offered as evidence go RED; words without numbers are review's province and are declared as
  such below.

### AC4 — the legacy escape hatch is enumerated, not open

Owned by **BC.B′**. Exactly 3 `legacy-unconditioned` markers pass.

Named RED mutation: **`MUT-LEGACY-FOURTH`** (**EVIDENCE**): add a 4th marker above the
AC3 orphan block → R5 REDs (`expected exactly 3 legacy markers, found 4`). Hard to silence: the
only green paths are deleting the marker (restoring AC3's RED) or editing the checker's
hardcoded count — a review-visible policy change, per the `verify_ail.sh` precedent.

### AC5 — null cases are loud, and the blast radius is exactly three files

Owned by **BC.B′**. Named RED mutation: **`MUT-CHECK-NULL`** (**EVIDENCE**): truncate
`bench/BASELINE.md` to empty (backup taken in the same command, byte-identical restore proven by
sha256 both sides — the iter-38/39 mutation discipline) → R6 REDs rather than passing vacuously.

Scope assertions, verified at review: the implementation changes **only**
`scripts/bench_worldd.sh`, `bench/BASELINE.md`, and `.github/workflows/ci.yml`. Untouched:
`host/daemon/bench_test.go` and its 10-name smoke manifest, `scripts/verify_ail.sh` and its
required-check manifest (P17), `scripts/verify_go.sh`, `go.mod` (**untouched because this item
selects no toolchain — the old "OD-1, parked" reason is discharged, P38**), all `.ail` files,
`world/`, `host/` production code, `tools/launchd/*`. `--smoke` output and semantics are
byte-compatible (CI's existing step must pass unmodified).

### AC6 — the toolchain probe measures each measured tree, and a REAL cross-toolchain pair REDs

> **DISCHARGED 2026-08-05 (iteration 50, controller pass `C2b`) — see P45.** Both mutations run
> outside any sandbox on the real floor-split fixture. `MUT-AB-FLOOR-SPLIT`: `rc=1`, **exactly one
> violation**, `✗ toolchain mismatch inside claimed A/B pair`, against a `rc=0` control arm on the
> un-appended file. `MUT-PROBE-CALLER-DIR`: the caller-cwd probe records `go1.26.4` for both trees
> and the gate **greens a genuinely cross-toolchain pair** — the round-1 bug reproduced, the fix
> shown to be what makes the known-positive fire. The regime clause below was **honoured and
> measured**, not assumed: under `auto` the straddle is `go1.26.5`/`go1.26.4`; under
> `GOTOOLCHAIN=go1.25.6` both trees read `go1.25.6` and the criterion would have passed vacuously.

Owned by **BC.B′** (the probe behavior it exercises lands in BC.A′). This is round-1
objection 2's reviewer-required probe. ~~it runs against the exact condition OD-1 will
create~~ **SUPERSEDED (iteration 48): that condition EXISTS at HEAD** — `4e/OD-1` is discharged,
the floor is `go 1.25.6`, and `go env GOVERSION` resolves to `go1.26.4` (P37/P38). The fixture is
a **throwaway worktree pair** — the iter-40/41 measurement pattern: the real `go.mod` is never
touched, the floor edit exists only as a never-pushed commit, and the worktrees are removed
afterwards.

**Fixture (re-pointed at the real straddle, iteration 48 — P37).** ~~two throwaway commits, A
raising the floor to `go 1.26.5` and B on top of A reverting it to `go 1.26.4`, in two worktrees
under `/tmp`~~ **SUPERSEDED on two counts.** (1) **HEAD is already a valid control**, so only ONE
throwaway commit is needed: `variant` = a detached worktree at HEAD with the floor edited to
`go 1.26.5` and committed (never pushed); `control` = a detached worktree at **HEAD itself**,
which is the variant's parent, so the parent-edge check passes by construction. Measured: floor
`1.25.6` → `go1.26.4`, floor `1.26.5` → `go1.26.5` — the genuine straddle this AC needs.
(2) **Neither worktree may live under `/tmp`** — the shared mission-control skill forbids it
outright (a `/tmp`-rooted checkout fails CWD-resolving tests for its *location*, producing a red
CI cannot reproduce). Place both as **siblings of the repo**, e.g.
`<repo>/../.bench-floorsplit-variant` and `<repo>/../.bench-floorsplit-control`. Both toolchains
are cached on this rig (P37 re-listing: 6 cached, `1.26.5` among them), so no network is
involved. Record the pair in ONE session:
`--record-pair --variant <…>/.bench-floorsplit-variant --control <…>/.bench-floorsplit-control`,
and append it to `bench/BASELINE.md` under the AC5 backup/sha256-restore discipline. The
two-commit form remains permissible and produces the same result; the one-commit form is
recommended because it is strictly less setup for identical evidence. **The predicted RED message
is unchanged.**

> **AC6 REGIME CLAUSE — BINDING, AND MEASURED RATHER THAN ASSUMED (iteration 48, P41).** Five
> quorum rounds accepted this fixture without anyone running it. It only works under one regime,
> and under the repository's *own* required regime it fails **green**:
>
> - **The fixture session MUST be invoked with `GOTOOLCHAIN=auto` explicitly.** Under
>   `GOTOOLCHAIN=go1.25.6` — the regime `scripts/verify_go.sh` requires and CI pins job-wide —
>   `go env GOVERSION` returns **the pin** for every tree, both sections record `go1.25.6`, and
>   `MUT-AB-FLOOR-SPLIT` **PASSES**. That is a vacuous known-positive in the one criterion whose
>   entire job is to prove the probe is not vacuous, so it must be stated as a precondition of the
>   run, not left to the invoking shell.
> - **Both members of the only achievable straddle are deny-listed compilers** (`go1.26.4`
>   control, `go1.26.5` variant; `verify_go.sh:40` deny-lists `go1.26.0`–`go1.26.5`), and there is
>   no deny-list-free alternative on this rig: floors and `toolchain` directives select only
>   UPWARD, the sole cached toolchain above the local `go1.26.4` is `go1.26.5`, and the control is
>   always the local toolchain. **This is acceptable and must be labelled, not silently absorbed**:
>   AC6 is a probe of *the recorder's per-tree `go -C` placement*, not a benchmark whose numbers
>   anyone will use. The planner measured that the `go1.26.5`-built `host/daemon` bench binary
>   still runs all ten rows `rc=0`, so the legs complete and a pair is emitted.
> - **The emitted acceptance-fixture pair MUST NOT be appended to `bench/BASELINE.md` as a
>   milestone baseline**, and its section label must name the deny-listed toolchains it was
>   recorded under. It is fixture evidence for AC6 and nothing else.
> - **The canary reds on the rig's default toolchain by design** (iteration 46). Every other
>   command in this sprint — builds, `go test`, `verify_go.sh`, the acceptance pair — runs under
>   `GOTOOLCHAIN=go1.25.6`. `GOTOOLCHAIN=auto` is scoped to the AC6 fixture session ALONE, and
>   the executor must not export it globally.

- **`MUT-AB-FLOOR-SPLIT`** (**EVIDENCE**, the known-positive): the pair is same-session (R4b/c/d
  genuinely satisfied) and variant-vs-parent (R4a genuinely satisfied) — and `--check-claims`
  must RED with exactly `toolchain mismatch inside claimed A/B pair` (`go1.26.4` vs
  `go1.26.5`), **not** any other limb's message. Unlike the hand-edited
  `MUT-CLAIM-TOOLCHAIN-SPLIT` (which isolates R4e's comparison logic), this proves on **real
  emissions** that the per-role `go -C` probe records each tree's own toolchain end-to-end and
  that R4e's tooth bites on a genuine floor-straddling pair.
- **`MUT-PROBE-CALLER-DIR`** (**HARNESS**): edit the recorder to drop `-C <dir>` from the
  `go env` probe — restoring the caller-cwd form the round-1 design shipped — and re-record
  the fixture pair (one bounded session). Both sections now record the CALLER's `go1.26.4`
  despite the control tree's 1.26.5 floor, and `--check-claims` **greens a genuinely
  cross-toolchain pair** — RED at review, because the known-cross-toolchain fixture must RED.
  This is the vacuity probe for the probe placement itself: it reproduces the round-1 bug and
  shows the fix, not luck, is what makes this AC's known-positive fire. sha256-restored.

### AC7 — a pair is single-session by grammar: identity, no reuse, real interleaving, no in-window compile

> **DISCHARGED 2026-08-05 (iteration 50) — `MUT-PAIR-ID-SPLIT` and `MUT-CONTROL-REUSE` at
> iteration 49; the three re-recording mutations at pass `C2b`, see P46/P47.** Every arm ran
> against an honest same-session control pair that greened first (`rc=0`, `2 well-formed pairs`),
> so each RED is attributable to the edit and not to the fixture. `MUT-PAIR-TWO-SESSIONS` → **2
> violations, R4d on both spliced sections, R4b silent** (the round-4 supersession, validated on a
> real pair); its named silencing attempt relocates to **R4c + R4b + an unnamed R4f**, because two
> sessions 36 s apart cannot fake an interleave. `MUT-PAIR-SEQUENTIAL` → **exactly one violation**,
> R4f, on an otherwise well-formed pair. `MUT-PAIR-INLINE-BUILD` → **R1 on both sections plus R2's
> orphan cascade**, and its stated secondary observable **REFUTED** (P47).

Owned by **BC.B′** (the behaviors it exercises land in BC.A′). This is the branch-A core — each
limb of the ratified scope gets its own discriminator, and each mutation predicts its exact
named failure:

- **`MUT-PAIR-ID-SPLIT`** (**EVIDENCE**): edit the acceptance control block's `pair_id` (one
  hex digit) **and recompute its `conditions_sha256`** — isolating R4 from R1. Predicted
  (**round-4 re-derived, P32**): **R4c REDs on the edited section, section-locally and BEFORE
  any pairing** (`pair_id does not derive from recorded session fields` — every derivation
  input is intact and the edited ID derives from none of them), **and R4d REDs BOTH sections
  as unpairable** (`unpairable conditions block — no counterpart shares its pair_id` — the
  edited ID and the variant's original ID each now appear on exactly one block).
  ~~R4b REDs `pair identity mismatch — sections are not from one session`~~ **SUPERSEDED
  (round 4)**: R4b is pair-level and NO PAIR FORMS under this mutation, so the round-3 text
  predicted a RED the specified mechanism could not produce — the reviewer's catch; the
  repair is R4c's section-locality plus R4d's named unpairable case, not a stronger R4b. No
  message prints an ID. The green path is re-running the session, not editing the other
  section to match — retyping the variant's ID to the edited value forms a pair but now
  fails R4c on BOTH sections, because the matched ID derives from neither section's fields.
- **`MUT-PAIR-TWO-SESSIONS`** (**EVIDENCE**, the round-2 reviewer's exact scenario): run the
  recorder
  TWICE on the same fixture commits (two real bounded sessions) and splice the FIRST session's
  control section with the SECOND session's variant section — a stale control attached to a
  fresh variant, each half individually pristine (R1/R3 green, and R4c green section-locally:
  each half derives from its own session's fields). Predicted (**round-4 re-derived, P32**):
  ~~R4b REDs~~ **SUPERSEDED — R4b never evaluates here, because no pair forms**; what fires is
  **R4d, REDding BOTH spliced sections as unpairable** (`unpairable conditions block — no
  counterpart shares its pair_id`) — the two sessions necessarily minted different
  `pair_nonce`/`session_utc`, so the IDs differ and neither section has a counterpart.
  Retyping one side's `pair_id` to force a pair merely relocates the RED: the retyped section
  fails R4c's derivation clause, and even past that, R4b REDs the differing
  `session_utc`/`pair_nonce`/endpoint fields. This is the mutation that
  proves independently-recorded halves can no longer form a claim — the round-2 objection,
  answered mechanically.
- **`MUT-CONTROL-REUSE`** (**EVIDENCE**): duplicate the acceptance pair's control section and
  its two raw blocks verbatim elsewhere in the file (a control serving twice). Predicted: R4d
  REDs `control reused across pairs / duplicated pair member` — the `pair_id` now appears on
  three blocks, two of them `role: control`. Hard to silence (**round-4 corrected, P32**):
  ~~giving the copy a fresh ID fails R4c (underivable) and R4b (no matching variant)~~
  **SUPERSEDED** — the round-3 wording said "underivable" without defining what the checker
  DOES with an underivable, unpairable block (the silent-skip invitation the round-4 quorum
  caught), and R4b cannot fire on a block that never forms a pair; giving the copy a fresh
  retyped ID now fails **R4c's derivation clause, section-locally** (a retyped ID derives
  from nothing) AND **R4d's unpairable named case** (the fresh ID has no counterpart) — two
  named REDs, neither skippable; deleting the copy is the restore.
- **`MUT-PAIR-SEQUENTIAL`** (**HARNESS**): edit the recorder's leg loop to run
  control/control/variant/variant — the un-interleaved shape branch A exists to forbid — while
  leaving emission and timestamping honest; re-record the fixture pair (one bounded session).
  Predicted: R4f REDs `legs not interleaved — control legs are not outermost` (the recorded
  `legN_seq`/start-time order shows control finishing before any variant leg ran). A harness
  that additionally forges its timestamps to fake interleaving is fraud-boundary, declared in
  *What these gates CANNOT fail*. sha256-restored.
- **`MUT-PAIR-INLINE-BUILD`** (**HARNESS**): edit the recorder to skip the prebuild (D1
  step 7) and run each leg as `go -C <dir> test -bench . -benchtime 200x -run '^$'
  ./host/daemon/` — folding compilation back into the measured window; re-record the fixture
  pair (one bounded session). Predicted: R1 REDs `invocation is not a prebuilt-binary
  invocation` on both sections (the honest recorder records the `go test` spelling it actually
  ran, and `binary_sha256` has nothing to hash). ~~Secondary observable at review: leg-1 elapsed
  jumps from seconds to a compile-bearing figure.~~ **STRUCK — REFUTED IN EXECUTION (P46/P47,
  iteration 50).** R1 fired exactly as written, on both sections, and dragged R2's orphan cascade
  with it (6 violations total). The *secondary observable did not*: honest legs measured 7,7,7,7 s
  against the inline-build's 8,8,9,8 s, because a warm Go build cache prices a full compile of
  these trees at **1–2 s** (measured via the recorder's own `prebuild_elapsed_s`) — inside this
  document's own within-condition noise. **Do not offer this as a review signal.** A reviewer who
  checked the timing instead of the rule would have cleared a session that really had folded
  compilation into the measured window. sha256-restored.

Mutation ceilings (the DG.A precedent): each named mutation is one run plus one restored
baseline run, 120-second ceiling per invocation — **except the recorder-invoking runs in AC6
and AC7 (`MUT-AB-FLOOR-SPLIT`'s fixture session, `MUT-PROBE-CALLER-DIR`, the two
`MUT-PAIR-TWO-SESSIONS` sessions, `MUT-PAIR-SEQUENTIAL`, `MUT-PAIR-INLINE-BUILD`), which are
instead bounded by the recorder's own hardcoded per-invocation deadlines (D1): the bound holds
because every external invocation is individually bounded (round 5), so a fully-wedged session
cannot exceed the sum of its bounded parts (2×600 s prebuild + 4×120 s legs + 3×120 s
probe-class + ~40×20 s utility-class ≈ 48 min; a real fixture session sits in
the tens of seconds, P27/P34)** — no retry loops; every mutation states its edited file and its
HARNESS/EVIDENCE classification above — none is presented as proof of a kernel property.

## Milestones

**Re-cut for branch A at the human's ratified sizing (~0.6–0.9 day total).** The charter's
standing observation that the original BC.A was "very nearly branch-independent" is **MOOT and
withdrawn**: branch A reworks the recorder's entire invocation surface (`--record <role>` →
`--record-pair`), its control flow (one session, prebuilds, four staged legs, end-only
emission) and its schema — so the original BC.A is **not a free head start**, and the original
~0.2/~0.3 sizing did not survive; pretending otherwise would be re-sizing by wishful
arithmetic. Nothing from the original milestones has shipped (P28), so nothing is thrown away
but the estimate.

### BC.A′ — the pair recorder (~0.35–0.5 day)

Modify only `scripts/bench_worldd.sh`:

- add `--record-pair --variant <dir> --control <dir>` per D1 (fixed 9-step session order:
  validation → rig-global probes → session stamp + nonce → per-role `go -C` probes → per-role
  git state + the parent-edge refusal → pair-ID derivation with the
  `pair_variant_commit`/`pair_control_commit` endpoints written into BOTH sections (round 4,
  P32) → both prebuilds with
  `binary_sha256` → four interleaved legs C/V/V/C with per-leg utc/load/competing/output/hash
  → end-only emission); the mirrored `run_bounded` helper with its provenance comment and the
  **four** hardcoded deadlines (`REC_PROBE_TIMEOUT_S=120`, `REC_PREBUILD_TIMEOUT_S=600`,
  `REC_LEG_TIMEOUT_S=120`, `REC_UTIL_TIMEOUT_S=20` — round 5: EVERY external-binary
  invocation, including `sysctl`/`ps`/`date`/`git`/`python3`, runs through the helper, P35);
  loud refusals including the named per-leg TIMEOUT; staging-dir
  all-or-nothing; python3-hashlib integrity fields; paste-ready two-section emission to
  stdout + mktemp file;
- update `usage()` for all modes (S7);
- leave `--smoke` byte-compatible.

Owns **AC1** and **AC2**. Independently CI-green: no CI change lands here; the existing smoke
step passes unmodified. Acceptance evidence: the shadowed-probe refusal transcript, the
dirty-tree refusal transcript, the non-parent-pair refusal transcript, the `MUT-REC-STALL`
leg-3 timeout transcript (named TIMEOUT, zero fences, dead grandchild), the
`MUT-REC-UTIL-STALL` git-status timeout transcript (round 5 — the same three assertions
against a utility-class invocation), and one
controller-performed real pair emission with the independent recompute of every hash and of
the pair-ID derivation.

### BC.B′ — the claim-structure gate, the policy, and the wired file (~0.25–0.4 day)

- add `--check-claims` per D3 (R1–R6 with R4's six limbs, hardcoded path, no knobs);
- `bench/BASELINE.md`: policy header + single-session A/B procedure + leg reading rule +
  bounds-not-eliminates statement, 3 legacy markers, the amortisation pointer correction, and
  the labelled controller-recorded acceptance pair (D5);
- `.github/workflows/ci.yml`: the two `go-verify` steps per D4, step 2 pinned to the specific
  `hw.ncpu` marker (CF-M-1).

Owns **AC3**, **AC4**, **AC5**, **AC6**, **AC7**. Independently CI-green: lands with the file
already satisfying the grammar it enforces. Executes and records all twelve BC.B′ named
mutations (AC3's three, AC4's one, AC5's one, AC6's two, AC7's five — the fixture sessions via
the throwaway floor-split worktrees, removed afterwards) plus re-runs of AC1/AC2's checks
against the landed harness.

## Out of scope, with reasons

- **The threshold / refuse-on-load option from (ii)** — REJECTED, not deferred.
  `BASELINE.md:7-8` already calls noise-gating a shared runner a dishonest gate (S6), and a
  threshold gate assumes an idle rig will eventually happen — on a rig shared with a mission
  whose schedule this loop cannot see (P7), that is not guaranteed, so the gate would either
  waste iterations waiting or get its threshold quietly raised until vacuous.
  ~~"The A/B is correct under any load and is the mechanism that actually caught the 6.06×
  artefact."~~ **SUPERSEDED (round 2 → branch A)** — the refuted over-claim, same as D2's;
  ~~the honest form: the single-session interleaved A/B **bounds** load skew by the session's
  own duration~~ **SUPERSEDED AGAIN (round 5, `gpt5-6-sol`): still dimensionally wrong —
  duration bounds temporal separation, not divergence magnitude.** The honest form: the
  single-session interleaved A/B makes the legs contemporaneous and alternating, so monotonic
  rig drift largely cancels in the comparison, and it is the descendant of the mechanism that
  caught the 6.06× artefact (whose legs
  were contemporaneous by luck, not by construction). The conditions block makes load
  **visible**; nothing in this design gates on — or bounds — its value.
- **The reviewer's third limb — a measured acceptance rule for excessive within-pair load
  divergence — EXCLUDED by the ratification.** The human ratified branch A *as recommended*,
  and the recommendation took neither a divergence threshold nor any acceptance rule
  (OD-6 section, verbatim). Two reasons, both now on the record: (1) a within-pair
  load-divergence threshold is `BASELINE.md:7-8`'s dishonest noise gate (S6) wearing a new
  name; (2) **P29 is measured evidence that no such threshold is derivable on this rig** —
  within-condition noise is ~1.4× with load held essentially constant (5.54–6.39), so a rule
  would have to separate a real regression from noise of the same magnitude as the thing it
  gates on, and no data exists here to set that line defensibly. What the interleaving buys
  instead is stated in D2 and *What these gates CANNOT fail* (round-5 corrected wording): the
  legs are adjacent in time and alternating in order, so monotonic rig drift largely cancels
  in the comparison — episodic or leg-correlated load divergence is neither cancelled,
  detected, nor adjudicated by a
  threshold nobody can justify. The per-leg load and process snapshots put any divergence in
  front of the reviewing human, which is where this design leaves it.
- **(iii) — re-deriving the amortisation section** — DEFERRED. ~~blocked on **OD-1**~~
  **REASON SUPERSEDED (iteration 48, P38): OD-1 is DISCHARGED and the floor has already moved to
  `go 1.25.6` (`f19acac`).** The deferral itself is unchanged and its new reason is stronger,
  because it is structural rather than political: re-deriving the section means recording a
  clean-toolchain pair **through `--record-pair`**, which BC.A′ has not built yet (P39 — zero
  occurrences at HEAD, `--smoke` control = 2). You cannot use the mechanism in the milestone that
  creates it. The section stays pinned to M3.C idle-rig numbers and labelled so; (iii) becomes
  the named follow-up *"amortisation re-derivation"* — first clean invocation recorded through
  this item's recorder, superseding the acceptance pair — and not a milestone here. The D5
  pointer edit writes exactly this into `BASELINE.md`, in the corrected wording.
- **Any change to the `go.mod` floor or toolchain selection** — ~~that **is** OD-1, parked for
  Mark (P22)~~ **still out of scope, for a plainer reason now that OD-1 is discharged: this item
  has no business selecting a toolchain.** This design *records* `GOVERSION`; it selects nothing
  and asserts no version value. The AC6 fixture's floor edit lives only in a throwaway,
  never-pushed worktree commit (P37) and never touches the real `go.mod`.
- **Benchstat, statistical machinery, more samples** — P10: not installed; nothing here needs
  it; adding a pinned dependency is a separate decision with its own doc if a future item wants
  distributional comparisons.
- **Changing `host/daemon/bench_test.go`** — the instrument's measurement side is not this
  item's defect; the recording side is.

## Design Freeze

The executor may not quietly change these invariants:

- `--record-pair` emits **all-or-nothing at SESSION granularity**: any probe failure, dirty
  tree, failed parent-edge check, failed prebuild, or failed/expired leg — including leg 3 of
  4 — produces a named non-zero exit and **zero** emitted fences. No `unknown`, no partial
  block, no half-pair, no single-section emission, no env-var escape hatch. Emission happens
  only after leg 4 (D1 step 9) — that ordering is itself frozen.
- The session order is frozen: validation → rig-global probes (first: `sysctl -n hw.ncpu`) →
  session stamp + nonce → per-role tree-dependent probes (`go -C <dir> env …`) → per-role git
  state + parent-edge refusal → pair-ID derivation → both prebuilds → four legs → emission.
  D4's CI determinism depends on the rig-global-first half; R4e's cross-toolchain tooth
  depends on the per-role-`-C` half. **Moving any tree-dependent probe to the caller's cwd is
  the CF-K-2 bypass** (`MUT-PROBE-CALLER-DIR` demonstrates it) and is a policy change
  requiring review.
- **The leg order `control,variant,variant,control` is frozen policy** — not configurable, not
  reorderable. Running roles sequentially is the un-interleaved shape the round-2 reject
  refuted (`MUT-PAIR-SEQUENTIAL` demonstrates it; R4f REDs it).
- **Both binaries are prebuilt before any leg runs** (D1 step 7); a leg invokes only the
  compiled test binary with `-test.*` flags. Folding compilation back into a leg is
  `MUT-PAIR-INLINE-BUILD`'s shape; R1's invocation grammar REDs it.
- **No summary statistic over legs, ever**: all four raw leg outputs reach the file verbatim;
  the recorder computes no mean, median, or "representative" value across legs (the iter-46
  spine, frozen).
- The pair ID is derived exactly as D1 step 6 specifies (`bench-pair/2` canonical string,
  python3 `hashlib`, nonce from python3 `secrets`) and is recomputed by R4c
  **section-locally**; changing the
  derivation string or dropping the nonce is a policy change requiring review.
- **The pair's commit endpoints are recorded in BOTH sections** (round 4, P32) as
  `pair_variant_commit`/`pair_control_commit`, identical verbatim across the pair, and the
  pair-ID derivation is **recomputable from any ONE section's own fields, with no pairing
  required first** (R4c clause i), with R4c clause ii binding the endpoints to each
  section's own `commit`/`parent`. Dropping either field, making the schema asymmetric
  (control-only), or moving the derivation back to pair-scope is a policy change requiring
  review.
- **The checker's evaluation order is frozen** (round 4, P32): R6 → R1 → R2 → R3 → R4c
  (section-local, BEFORE any pairing) → grouping by `pair_id` with R4d's cardinality and
  unpairable named RED → R4a/R4b/R4e/R4f on well-formed pairs → R5; every evaluable rule
  runs, every failure is reported, and **an unpairable conditions block is RED by name,
  never skipped** — "cannot evaluate" is never "nothing to report".
- **EVERY external-binary invocation in the session — `sysctl`, `ps`, `date`, every `git`,
  every `python3` helper, `go env`, `go test -c`, the four leg binaries, `$AILANG_BIN
  --version` — runs through the `run_bounded` mirror** (provenance
  comment: `verify_ail.sh:61-74`, V26, P23/P30) with hardcoded deadlines
  `REC_PROBE_TIMEOUT_S=120` / `REC_PREBUILD_TIMEOUT_S=600` / `REC_LEG_TIMEOUT_S=120` /
  `REC_UTIL_TIMEOUT_S=20`; none is
  env-overridable; expiry SIGKILLs the whole process group and is a named non-zero refusal
  emitting **zero** fences; behavioral divergence between the mirror and its `verify_ail.sh`
  original is a review-visible change. **The widened round-5 coverage is itself frozen
  policy: adding any external invocation OUTSIDE the helper is a policy change requiring
  review** — the covered-list has been wrong twice inside this document (round-2 self-catch,
  round-5 quorum catch, P35), which is exactly why the rule is now closed-world rather than
  an enumeration that can silently lag the code.
- Both `--variant` and `--control` are mandatory; there is no default, no single-role mode,
  and no two-invocation fallback.
- `competing_*` capture records verbatim paths and parent comm; it never maps a process to a
  mission name.
- All integrity hashes (per-leg outputs, binaries, conditions, pair-ID derivation) are
  computed with python3 `hashlib` in both the recorder and the checker; no shasum/sha256sum
  CLI dependency enters the gate.
- **No failure message in recorder or checker prints an expected digest, an expected SHA, an
  expected pair ID, or an expected nonce** — the iter-41 re-green lesson is a frozen property,
  not a style choice.
- `--check-claims` reads the hardcoded `bench/BASELINE.md` path; no path argument, no env knob;
  the legacy-marker count (3) and the raw-block detector regex are hardcoded policy.
- Raw-block detection is the in-fence `^Benchmark…` regex (P13), never an `ns/op` substring
  count.
- R4's limbs are all frozen: the parent edge (R4a), pair identity across
  `pair_id`/`session_utc`/`pair_nonce`/`pair_variant_commit`/`pair_control_commit`/
  `leg_order` (R4b), section-local pair-ID recomputation + endpoint role-binding (R4c), the
  exactly-two-blocks-per-ID cardinality **including the unpairable named RED** (R4d), the
  toolchain-identity comparison over
  `goversion`, `goos_goarch`, `ncpu`, `hw_model` (R4e), and the interleave-structure check
  (R4f) — removing any limb, any compared field, or the unpairable null case is a policy
  change requiring review.
- `--smoke` and its 10-name manifest are byte-compatible; `bench_test.go`, `verify_ail.sh`,
  `verify_go.sh`, `go.mod`, and all `.ail` files are untouched.
- The acceptance pair lands with its "NOT a milestone performance baseline / superseded by the
  first amortisation re-derivation" label verbatim (**wording corrected iteration 48 — the old
  "post-OD-1" spelling names a discharged decision, P38**); the executor does not present its
  numbers as a baseline.
- Every named mutation records its edited file, classification, single-run-plus-restore
  discipline, and sha256-verified restoration.

## What these gates CANNOT fail

These limits are part of the design, not residual fine print:

- **They cannot prove a number was not fabricated.** The hashes — including the pair ID and
  its nonce derivation, and including R4f's timestamps — are tamper-evidence binding the
  recorded conditions, the recorded outputs, and the file together; anyone who can run python3
  can forge all of them coherently, including a fictitious "interleaved" session that never
  ran. The protection against invented values remains what it has been for six milestones: the
  provenance-honesty culture (`<CONTROLLER-MEASURED>`, refusal to invent) plus review. This
  gate makes *silent drift, careless retyping, and honest-harness shortcuts* red, not fraud.
- **They cannot re-verify the git parent edge at CI time.** R4 checks recorder-emitted structure;
  the edge is git-true at record time by construction, but a squash-merged branch SHA may be
  unreachable later and CI checkouts are shallow (P4). A coherent forged pair with a fictitious
  parent SHA passes R4. Stated, accepted, bounded by review.
- **They cannot police prose.** A cost claim written only in words, citing no numbers, attaches
  no evidence and trips nothing here. What the grammar guarantees is narrower and real: no
  *benchmark numbers* enter the file's evidence without machine-recorded conditions and a
  parent-commit control.
- **They cannot determine whether an observed delta supports a cost claim — and a pair can
  PASS every rule while the delta means nothing** (round 5, `gpt5-6-sol`'s three concrete
  scenarios, on the record): **a competitor affecting only the variant legs**; **arbitrarily
  divergent per-leg load** between the roles; **an effect smaller than the recorded
  within-role spread**. R1–R6 establish provenance, session integrity, contemporaneity, pair
  identity, non-reuse, and toolchain/hardware identity — nothing more. Pair identity and the
  C/V/V/C order prove the legs are CONTEMPORANEOUS; they do not prove them COMPARABLE or
  load-independent. This is the ratified OD-6 scope (the excluded third limb, Out of scope,
  P29); whether the wider claim-validity question re-enters scope is **`4f/OD-8`**, the
  human's decision, not this document's.
- ~~**An interleaved single-session pair BOUNDS load skew; it does not ELIMINATE it.** Both
  roles share one rig episode, so the confound is bounded by the session's own duration.~~
  **SUPERSEDED (round 5, `gpt5-6-sol`): the sentence written as this document's honest
  statement was itself an over-claim — the session's duration bounds the TEMPORAL SEPARATION
  of the legs, not the MAGNITUDE of load divergence — the exact failure mode this document
  exists to eliminate, committed in the sentence that was supposed to be the safe one.** What
  is true: the interleaving makes the roles' legs adjacent in time and alternating in order,
  so a MONOTONIC drift in rig conditions affects both roles roughly equally and largely
  cancels in the comparison; a load episode shorter than the session — a competitor that
  starts or stops between legs, or one correlated with leg boundaries — does NOT cancel and
  is NOT detected, and nothing here gates on load (deliberately: the excluded third limb, Out
  of scope, P29). The per-leg load and process snapshots make any within-session shift
  visible to a HUMAN reader, who judges it; the CHECKER mechanically enforces nothing about
  load. The round-2 claim that the A/B form is "correct under any load" is
  SUPERSEDED and must not be restated anywhere.
- **A pair cannot resolve an effect smaller than its own within-role spread.** P29 measured
  ~1.4× within-condition noise at essentially constant load; the four-leg emission shows the
  reader that spread (leg 2 vs 3, leg 1 vs 4) instead of averaging it away, and a delta inside
  it is not evidence of anything — **plainly: an effect below that ~1.4× scale is not
  distinguishable by this mechanism at all** (the third round-5 scenario above), and the
  recorded second leg per role exists precisely so a reader SEES the spread instead of being
  told a single number. The file reports the sample; conclusions stay with the
  reader.
- **The bounding helper cannot police its own substrate, and the coverage list itself has now
  failed twice.** `run_bounded` is a python3 wrapper: the wrapper's own interpreter startup is
  the one wait it cannot bound, because it IS the bounding mechanism — a residual inherited
  verbatim from `verify_ail.sh` (V26/P30), stated here rather than absorbed into "by
  construction". And the enumeration of what is bounded has been wrong twice INSIDE THIS
  DOCUMENT: round 2's draft left `$AILANG_BIN --version` unbounded inside the remedy for
  unbounded waits (self-caught), and rounds 3–4 left `sysctl`/`ps`/`python3`/`git` unbounded
  while asserting every wait was bounded (quorum-caught, P35). The round-5 rule is
  closed-world — every external invocation goes through the helper, and `MUT-REC-UTIL-STALL`
  probes a member of the class that was missed — but a future edit adding an external call
  outside the helper is exactly the historical failure shape, and review must treat any new
  external invocation as unbounded until shown otherwise.
- **`ps` %CPU is a decaying average sampled at four instants** (one per leg start), not a
  profile of the run; a competitor that starts and stops between captures is invisible. The
  per-leg loads ~~bound~~ **sample (round 5 — same dimensional discipline)**, but do not
  narrate, each leg's interval.
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
- **The deadlines cannot distinguish a wedge from a legitimately slower run.** A prebuild past
  600 s or a leg past 120 s is refused loudly — the safe direction — and each ceiling is a
  hardcoded constant (prebuild: ~4.7× the worst measured compile, P25; leg: an
  explicitly-labelled engineering margin of ~15× over P27's five-run loaded SAMPLE, which is
  an order-of-magnitude reading, not a bound), raised only by a review-visible edit, never by
  a knob. A rig that legitimately exceeds them has changed enough that the constants *should*
  be re-derived in review.
- **The leg deadline is sized from a sample, and the design says so.** Five runs on one loaded
  rig at one commit (P27) establish seconds-not-minutes and nothing about the tail; if the
  tail bites, the failure is a named refusal, not a silently wrong number — which is the
  direction this mission accepts.
- **The 120 s probe bound can refuse a legitimate first-ever toolchain download** on a slow
  network. The refusal names the probe either way; the remedy is pre-caching the toolchain,
  which this rig already does for every floor on either side of OD-1 (P24).

## Open decision for the human — **ANSWERED: OD-6 RATIFIED, BRANCH A · OD-8 RATIFIED, EVIDENCE-NOT-CLAIMS**

> Numbering note (controller, iteration 42): this is **OD-6**, continuing the mission-wide
> sequence. OD-1/OD-2 belong to `w-race-gate-blindspot` (item 4e); OD-3/OD-4/OD-5 belong to
> `w-ddl-gate-teeth` (item 4d). This section is controller-authored bookkeeping around the
> designer's document; it changes no design text.

> **RATIFICATION (attended stamp 2026-08-03, recorded in the ratified charter at
> `design_docs/world-mission-status-archive.md:4`), verbatim:**
>
> > (3) **4f OD-6 = BRANCH A**: single-session interleaved `--record-pair` + pair ID +
> > control-reuse rejection (~0.6–0.9d, re-sized); mechanically valid cost claims — the
> > clause-4 floor's overhead measurement inherits this validity.
>
> The human took branch A as recommended: the recommendation included **neither** the
> reviewer's third limb (a measured within-pair load-divergence acceptance rule) **nor** any
> threshold, so that limb is **excluded from the ratified scope** — reasons in *Out of scope*
> (S6 + P29). Branch B's text is retained below, unedited, as the alternative that lost; per
> mission convention this section records the decision in place rather than rewriting the
> question to look prescient. The branch-A design is the body of this document (revision
> round 3, iteration 47).

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
load skew *(round-5 note on this pre-ratification phrasing, kept unedited per convention: the
"bounds load skew" form is itself SUPERSEDED — duration bounds temporal separation, not
divergence magnitude; the corrected statement is what *What these gates CANNOT fail* now
carries)*. Rationale: branch A keeps the charter row's promise, and the interleaving is the part
that actually earns the load-independence claim — iteration 39's pair worked *because* its legs
were minutes apart under the same load, which is the property branch A makes structural instead
of accidental. Estimated at ~0.6–0.9 day, i.e. the item outgrows its original sizing; that is the
cost of the promise, and it belongs to the human, not to me.

**BC.A is very nearly branch-independent.** The recorder's probes, bounded execution, refusal
semantics and integrity fields are needed under both branches; only its invocation surface
changes (`--record` vs `--record-pair`). If the queue needs forward motion before OD-6 is
answered, BC.A is the routable half — but that too is the human's call, since branch A would
rework the interface BC.A ships.

> **Post-ratification note (iteration 47):** the paragraph above is **MOOT** — branch A was
> ratified, it reworks exactly that interface, and the Milestones section re-cuts the work as
> BC.A′/BC.B′ at the human's ~0.6–0.9 day sizing. No pre-ratification BC.A work existed to
> salvage (P28). Kept unedited above per the no-prescience convention.

### OD-8 — **The ratification promises "mechanically valid cost claims". Branch A, built exactly as ratified, delivers mechanically complete, contemporaneous, tamper-evident EVIDENCE for a cost claim — not a validated one. Which did you mean?**

> Controller-authored bookkeeping (iteration 47), like the OD-6 section above. It changes no
> design text. **RAISED** 2026-08-04 at quorum round 5. **ANSWERED** 2026-08-04 (attended) — see
> the ratification box immediately below.

> **RATIFICATION (attended stamp 2026-08-04, recorded in the charter's STATUS block at
> `design_docs/world-mission.md`, commit `ea5e405`), verbatim:**
>
> > **Mark resolved 4f/OD-8: EVIDENCE, NOT CLAIMS.** The OD-6 stamp's wording is amended: branch
> > A delivers mechanically complete, contemporaneous, tamper-evident **EVIDENCE for a cost
> > claim** — claim VALIDATION (whether the observed delta statistically supports the claim:
> > repetitions, noise bounds) is deliberately out of 4f's scope and **folds into
> > `w-agent-floor-m4`'s experimental design** (which must specify N≥3 paired runs + noise
> > handling regardless; the floor's design doc gains this as a named requirement). 4f branch-A
> > milestones BC.A′/BC.B′ → ROUTABLE.
>
> **This is alternative (1) with (3) as bookkeeping** — the recommendation, and the objecting
> reviewer's own round-2 fallback. Three consequences, none of which is new design work:
> **(i)** the OD-6 ratification's phrase *"mechanically valid cost claims"* is superseded by
> *"mechanically complete, contemporaneous, tamper-evident evidence for a cost claim"* wherever
> this document or the charter relies on it — the title clause was already struck at round 5, and
> the D2 policy text must read the same way (BC.B′);
> **(ii)** `gpt5-6-sol`'s round-5 blocking limb is **resolved by the human, not overridden by the
> controller** — the reviewer disputed the *claim*, the human corrected the *claim*, and the
> mechanism it disputed is unchanged. The narrow-refinement carve-out was correctly NOT used;
> **(iii)** claim validation (N≥3 paired runs, noise handling) is now an owned, named requirement
> of `w-agent-floor-m4`, not an unowned gap. It is **out of scope here** by ratification.
>
> **Everything the round-5 quorum fixed stands**: the closed-world bounding rule (P35), the
> dimensional corrections to the skew/separation sentence, and the enumerated CANNOT-fail
> scenarios.

**How it arose.** `gpt5-6-sol` rejected the branch-A design on the same ground it rejected in
round 2, and this time the objection landed on the *ratification's own wording* rather than on the
mechanism: "R1–R6 validate provenance and session structure, but they never determine whether the
observed delta supports a cost claim. … Pair identity and C/V/V/C ordering prove contemporaneity,
not comparability or load independence." Its full text is preserved in the Quorum verification log.

**What is true, stated plainly.** The reviewer is right about the mechanism. R1–R6 establish
provenance, session integrity, contemporaneity, pair identity, non-reuse, and toolchain/hardware
identity — and **nothing about load**, by ratified design, because the ratification took branch A
*as recommended* and the recommendation excluded the reviewer's third limb (a measured within-pair
load-divergence acceptance rule). Three scenarios pass the grammar and are now enumerated in *What
these gates CANNOT fail*: a competitor affecting only the variant legs; arbitrarily divergent
per-leg load; an effect smaller than the ~1.4× within-role spread P29 measured at essentially
constant load. The round-5 revision also struck the doc's own "bounds load skew by the session's
duration" formulation as **dimensionally wrong** — duration bounds temporal *separation*, not skew
*magnitude*.

**Why this is NOT a re-ask of OD-6.** OD-6 asked *which branch*: build the contemporaneity
mechanism (A), or keep the old mechanism and weaken the claim (B). Mark answered **A**, and A is
what the document now specifies. OD-8 asks a question OD-6 did not: the ratification's sentence
promises *"mechanically valid cost claims"*, and branch A cannot produce that — so does the
**claim** get corrected to match the mechanism, or does the **mechanism** grow to match the claim?
A headless loop must not choose, because either answer changes what item 4f delivers and one of
them changes the charter's own words.

**Alternatives.**

1. **Correct the claim; ship branch A as designed.** The pair is mechanically complete,
   contemporaneous, tamper-evident *evidence*; the recorded per-leg loads and the second leg per
   role are what let a **human** judge comparability, and the checker enforces admissibility, not
   validity. Cost: the charter row's promise is met in evidence, not in validation, and the
   clause-4 floor's overhead measurement inherits *evidence* rather than *validity*. **No new
   design work; BC.A′/BC.B′ are routable immediately.**
2. **Grow the mechanism: add a comparability criterion** (the excluded third limb, in some
   defensible form — e.g. requiring the observed effect to exceed the measured within-role spread
   by a stated factor). Cost: real new design work, a re-quorum, and re-sizing beyond ~0.6–0.9 d.
   The controller's standing objection is unchanged and now has a measurement behind it: P29 puts
   within-condition noise at ~1.4× at essentially constant load, and `BASELINE.md:7-8` already
   calls noise-gating a shared runner a **dishonest gate (S6)**, so a threshold derived from this
   rig would be a number chosen rather than measured.
3. **Amend the ratification's wording in the charter** and otherwise proceed as (1) — the same
   mechanism, with the correction recorded where the promise was made rather than only in the doc.

**Recommendation: (1), optionally with (3) as bookkeeping** — and the strongest argument for it is
the objecting reviewer's own: `gpt5-6-sol`'s round-2 `proposed_fix` already named this outcome as
its fallback — *"If no defensible comparability rule can be derived, revise the policy to say the
pair is mechanically complete evidence—not a mechanically valid cost claim."* No defensible rule
has been derived, and P29 is the measurement saying why. Alternative 2 is the only one that writes
code, and the loop does not believe the threshold it would need is obtainable on this rig.

**Cost of deferring:** item 4f stays parked and the queue has no other routable item; BC.A′/BC.B′
are ~0.6–0.9 d and unblock the moment this is answered. **Every round-5 fix stands regardless of
the outcome** — the closed-world bounding rule, the dimensional corrections, and the enumerated
CANNOT-fail scenarios are improvements under either answer.

---

> **CONTROLLER NOTE (iteration 47) — TWO PREMISE ROWS IN THIS DOCUMENT ARE STALE AGAINST HEAD, AND
> FIVE QUORUM ROUNDS DID NOT CATCH IT.** Verified first-party at `61348b9`: `go.mod:3` now reads
> **`go 1.25.6`**, not the `go 1.26.4` recorded in **P9**, and **`4e/OD-1` is DISCHARGED**, not
> "parked / awaiting Mark" as **P22** records — Mark ratified it and it landed at iteration 46 in
> `f19acac` (`git show --stat f19acac -- go.mod` → 1 file changed). Three consequences for whoever
> picks this up after OD-8, none of them fatal and all of them cheaper to know than to rediscover:
> **(a)** the *Deferred* item below is blocked on a decision that no longer exists — limb (iii)'s
> stated reason is void and needs re-deriving or re-justifying, not inheriting; **(b)** AC6's
> floor-split fixture reasons about "the exact condition OD-1 will create", which now EXISTS at
> HEAD rather than being hypothetical, so the fixture should be re-pointed at the real straddle
> before it is built; **(c)** the impact is narrower than it looks, because the `go` directive is a
> **floor** — `go env GOVERSION` in this repo still reads `go1.26.4` under `GOTOOLCHAIN=auto`, so
> every *recorded* toolchain value in the design is unchanged. This is Gate-2's rule 3b(vi) in its
> own habitat: a long document is an instrument, its Verification Log is the control, and reviewers
> read for design soundness rather than for freshness against HEAD — so a doc parked across four
> iterations accumulates premises that were true when written and are not true now.

> **CONTROLLER NOTE (iteration 48) — THE THREE CONSEQUENCES ABOVE ARE NOW MEASURED, NOT
> FORECAST.** The note above was written from a *reading* of HEAD; it correctly named what needed
> re-checking and then, correctly, left the re-checking undone. Iteration 48 ran it. All commands
> and outputs are recorded as rows **P37–P39**; the operative results:
> **(a) — RE-JUSTIFIED, NOT VOID.** Limb (iii) is no longer blocked on OD-1 (discharged, P38), but
> it does not therefore come into scope: it is blocked on `--record-pair` **existing**, since its
> whole content is "record the first clean-toolchain pair *through this item's recorder*". The
> dependency changed identity, not direction — it now points at BC.A′ landing rather than at a
> human decision. Restated in *Deferred* below. **This is a correction of a stated reason, not a
> change of scope.**
> **(b) — THE STRADDLE IS REAL, AND IT MAKES AC6's FIXTURE CHEAPER, NOT DIFFERENT (P37).** Measured
> in a throwaway sibling worktree at HEAD: floor `go 1.25.6` → `go -C <dir> env GOVERSION` =
> **`go1.26.4`**; the same tree with the floor edited to `go 1.26.5` → **`go1.26.5`**; the main
> tree unchanged throughout (control). So the AC6 fixture no longer needs **two** throwaway
> commits (A raises the floor, B reverts it): **HEAD itself is a valid control**, because its own
> floor already resolves to a different toolchain than a `go 1.26.5` child would. One commit on
> top of HEAD (`variant`, floor → `1.26.5`), with HEAD as `control`, satisfies the parent-edge
> check by construction and produces the same genuine `go1.26.4` vs `go1.26.5` straddle
> `MUT-AB-FLOOR-SPLIT` requires. The two-commit form still works and remains permissible; the
> one-commit form is the recommended simplification and it is what the executor should build.
> The predicted RED message is unchanged.
> **(b2) — THE FIXTURE WORKTREES MUST NOT LIVE UNDER `/tmp`.** AC6's text says
> `git worktree add /tmp/bench-floorsplit`. The shared mission-control skill now forbids this
> outright: a `/tmp`-rooted checkout fails tests that resolve paths against the CWD *for the
> location rather than the code*, and CI never reproduces the red. Use a **sibling of the repo**
> (e.g. `<repo>/../.bench-floorsplit-<n>`), as iteration 48's own probe did. This overrides the
> path literals in AC6 and anywhere else in this document; nothing else about the fixture changes.
> **(c) — CONFIRMED, and it is why (a)/(b) are cheap:** every *recorded* `GOVERSION` value in this
> design is unchanged, because the `go` directive is a floor (P37's control leg).

## Deferred

- **Amortisation re-derivation** (item (iii)) — first clean-toolchain invocation, recorded via
  `--record-pair` as a proper single-session pair, supersedes the acceptance pair and re-derives
  the amortisation ratios from same-conditions rows. ~~Blocked on **OD-1**~~ **SUPERSEDED
  (iteration 48): OD-1 is DISCHARGED** (`f19acac`, iteration 46 — P38). It remains deferred for a
  different and stronger reason: it consumes `--record-pair`, which **BC.A′ has not built yet**
  (P28 — zero occurrences at HEAD, `--smoke` control = 2). So it is blocked on BC.A′ landing, a
  dependency this item creates rather than one a human must resolve. Named in `BASELINE.md` by the
  D5 pointer edit — whose wording must say *"pending the first recorded clean-toolchain pair"*,
  **not** *"pending OD-1"* (BC.B′).
- **`w-race-gate-blindspot` AC7 hookup** — ~~when OD-1 resolves and~~ **the race doc's toolchain
  pin has LANDED** (`f19acac`, iteration 46), so this hookup is unblocked on its stated condition
  and waits only on this item's mechanism existing. Its condition-change record uses this item's
  recorder (its `:184` already assigns ownership here). **No action in this item** — unchanged.

---

## Quorum verification log (pick-time quorum, iteration 42)

> **Note for the next reviewer (iteration 47): this document is the BRANCH-A REVISION,
> authored AFTER the human ratified OD-6 = branch A** (attended stamp 2026-08-03; verbatim
> quote in the OD-6 section). What changed relative to the round-2 text you may have seen:
> D1 is now a single-session `--record-pair` with prebuilt binaries, a nonce-derived
> recomputable pair ID, and four interleaved legs (C/V/V/C) with per-leg
> timestamps/load/process snapshots; R4 gained limbs R4b–R4d and R4f (pair identity,
> recomputation, control-reuse cardinality, interleave structure); the schema moved to
> `bench-conditions/2`; the single 600 s benchmark deadline was replaced by
> prebuild/leg deadlines (P26/P27 — discharging CF-M-2); D4's CI assertion is pinned to the
> specific `hw.ncpu` marker (discharging CF-M-1); AC7 adds five branch-A mutations; the
> "correct under any load" claims are struck SUPERSEDED in D2 and Out-of-scope; the
> reviewer's third limb (load-divergence acceptance rule) is explicitly excluded with
> measured support (P29). The two BLOCKED rounds below are the history that produced this
> revision, retained unedited. **Round 4, appended below, is the quorum ON this branch-A
> revision**; its one reject produced revision round 4 — the pair-scope fix (P32) — applied
> in place in this document. **Round 5, appended below that, is the quorum on the round-4
> text**: BLOCKED, two rejects, both adopted in revision round 5; the second objection's
> scope limb is DEFERRED to `4f/OD-8`, not resolved here, and the document is PARKED on it
> (see Status).

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

**Round 4 — the branch-A revision's quorum (iteration 47): BLOCKED, 1 reject. The catch is
real; its stated mechanism is overstated; both facts are on the record (P32).** (Numbering
continues the controller's mission-log sequence for this item; rounds 1–2 above were the
iteration-42 pick-time quorum on the pre-branch-A text.)

| Reviewer | Verdict | Objection |
|---|---|---|
| `gemini-3-1-pro` | **reject** | Verbatim catch: *"The control block's schema is missing an input field required by the `pair_id` hash derivation string."* Full objection: D1 step 6 derives `pair_id` from `bench-pair/2\n<session_utc>\n<variant_commit>\n<control_commit>\n<pair_nonce>\n` and R4c claims recomputation "from that pair's own recorded fields", but the `role: control` block recorded no `variant_commit` — only its own `commit` and `parent` — so, the reviewer argued, R4c is "mathematically impossible to evaluate on a control block as described", invalidating AC7's `MUT-PAIR-ID-SPLIT` prediction |

Verbatim `proposed_fix`: *"Add `variant_commit: <full sha>` to the `role: control` block's
emission schema in D1 so it possesses all fields needed to recompute the ID. Alternatively,
remove the commits from the canonical byte string (e.g., `bench-pair/2\n<session_utc>\n
<pair_nonce>\n`), since `conditions_sha256` already provides tamper-evidence for the commit
fields (though this alternative would require dropping R4c from the `MUT-CLAIM-NONPARENT`
prediction)."*

**The controller measured the premise before routing it (P32), and the measurement sharpened
the objection rather than forwarding it.** The premise is TRUE — the control schema carried no
counterpart commit (grep of the schema block, with a live control). But "mathematically
impossible" is OVERSTATED for a well-formed pair: the checker reads the whole file, grouping
by `pair_id` supplies both commits, and R4c's own text said "from that PAIR's own recorded
fields" — so happy-path R4c was evaluable as written. The REAL defect sat one level in: the
derivation was **pair-scoped while the mutation predictions were section-scoped**, and no rule
defined the unpairable block — under `MUT-PAIR-ID-SPLIT` the edited section cannot be paired,
so the predicted R4c RED could not fire as specified; `MUT-CONTROL-REUSE`'s own text already
wrote "R4c (underivable)" without defining what underivable DOES; and no evaluation order for
R4a…R4f existed anywhere (grep: zero hits; same-call control: `R4` × 59). A rule whose failure
mode is undefined for exactly the input its mutation produces is this mission's signature
shape — an implementer could read "cannot evaluate" as "skip", a silent pass.

**Disposition — ADOPTED, the FIRST branch of the proposed fix** (the alternative — dropping
the commits from the canonical string — would route around the circularity instead of removing
it: it weakens the ID's binding to the pair's endpoints and pays for it by dropping
`MUT-CLAIM-NONPARENT`'s R4c limb, as the reviewer's own parenthetical concedes). Applied, all
in place: symmetric `pair_variant_commit`/`pair_control_commit` in BOTH sections (D1 step 6 —
symmetric rather than the control-only field literally proposed, for the stated
reader-facing reason: one role-independent derivation rule instead of two, endpoints as
session-level facts policed by R4b, redundancy turned into R4c's binding check); R4c made
section-local (derivation + role-binding) and ordered BEFORE pairing; the checker's
evaluation order stated in D3 and frozen; R4d's unpairable case made a named RED, never a
skip; and every prediction naming R4c re-derived against the fixed mechanism —
`MUT-PAIR-ID-SPLIT` and `MUT-PAIR-TWO-SESSIONS` corrected (unsupported R4b predictions struck
SUPERSEDED), `MUT-CONTROL-REUSE`'s silence note corrected, `MUT-CLAIM-NONPARENT`'s dual-fire
retained with its mechanism corrected to R4c's binding clause.

**Round 5 — the round-4 revision's quorum (iteration 47): BLOCKED, 2 rejects. Both are
right: one is a mechanism defect (fixed), one is a claim-accuracy defect (fixed) carrying a
scope limb that is the human's (deferred to `4f/OD-8`, NOT resolved here).**

`gemini-3-1-pro`, **reject**, verbatim:

> "The document asserts that 'every individual wait is bounded, and the session's worst case
> is therefore the sum of its parts by construction' (D1), but it exclusively applies the
> `run_bounded` helper to `go` invocations and `$AILANG_BIN`. External binary invocations in
> Step 2 (`sysctl`, `ps`), Step 3/6 (`python3`), and Step 5 (`git status --porcelain`,
> `git rev-parse`) are left unbounded. A wedged filesystem or kernel lock during a `git` or
> `sysctl` invocation will stall the recorder infinitely, violating the 'bounded waits' axiom
> and falsifying the 'sum of its parts' mathematical bound."

**Disposition — ADOPTED in full (P35, reproduced first-party against the round-4 text before
editing: 15 bounding-related lines, zero covering a recorder
`sysctl`/`ps`/`git`/`date`/python3 invocation, control sysctl × 20).** This is the SECOND
instance of the same trap inside this document — round 2's draft left `$AILANG_BIN
--version` unbounded inside the remedy for unbounded waits (self-caught, recorded under
round 2 above); the remedy then grew a second set of unbounded waits, and the "sum of its
parts" sentence claiming otherwise was exactly the unfalsifiable summary this item exists to
eliminate. Fix, all in place: closed-world bounding (EVERY external-binary invocation
through `run_bounded` — D1 steps 1/2/3/5/8, the Bounded execution enumeration); a fourth
deadline class `REC_UTIL_TIMEOUT_S=20`, justified from first-party measurement (P36: worst
utility member 0.050 s → an explicitly-labelled 400× engineering margin, NOT the 120 s
probe constant, which nothing in the class can justify); the sum-of-parts sentence struck
SUPERSEDED and restated with its arithmetic (≈48 min fully wedged) and its one named
residual (the helper's own interpreter startup — *What these gates CANNOT fail*);
`MUT-REC-UTIL-STALL` (AC1) as the RED discriminator on a newly-covered invocation; the
recurrence itself recorded in *What these gates CANNOT fail*; coverage frozen (Design
Freeze). Tally: mutations 15 → 16 (the refusal family remains ONE production discriminator —
not padded — so 13 stands).

`gpt5-6-sol`, **reject**, verbatim:

> "The central promise remains mechanically unsupported: R1–R6 validate provenance and
> session structure, but they never determine whether the observed delta supports a cost
> claim. A pair can pass with a competitor affecting only the variant legs, arbitrarily
> divergent per-leg load, or an effect smaller than the recorded within-role spread. Pair
> identity and C/V/V/C ordering prove contemporaneity, not comparability or load
> independence. The repeated claim that session duration "bounds load skew" is also
> unverified and dimensionally wrong: duration bounds temporal separation, not the magnitude
> of load skew."

**Disposition — limb (a), the dimensional error: ADOPTED in full.** The phrase "bounds load
skew by the session's own duration" — written as this document's HONEST statement — was
itself an over-claim, the failure mode this item exists for, committed in the sentence that
was supposed to be the safe one. Every site (revision-3 preamble, D2's superseded block, D2's
closing prose, D5's header rule, both Out-of-scope bullets, *What these gates CANNOT fail*,
plus an annotation on the pre-ratification OD-6 phrasing) is struck SUPERSEDED and replaced
with what is derivable: adjacent alternating legs cancel MONOTONIC drift; episodic or
leg-correlated divergence is neither cancelled nor detected; the per-leg load/`ps` fields are
evidence for a HUMAN's judgment and the checker enforces nothing about load, by ratified
design. The reviewer's three pass-anyway scenarios are folded verbatim in substance into
*What these gates CANNOT fail*, with the third tied to P29's ~1.4× spread: an effect below
that scale is not distinguishable by this mechanism at all.

**Disposition — limb (b), "R1–R6 never determine whether the observed delta supports a cost
claim": TRUE, and DEFERRED to `4f/OD-8`, not resolved here.** This is the reviewer's round-2
position restated, and it is the question the human already answered: OD-6 ratified branch A
over branch B, explicitly excluding the third limb, so the mechanism cannot validate a cost
claim BY RATIFIED DESIGN. What the objection newly exposes is a gap between the
ratification's WORDING ("mechanically valid cost claims") and what branch A delivers —
mechanically complete, contemporaneous, tamper-evident EVIDENCE for a cost claim. That
wording gap is a decision about what item 4f promises; the controller is raising it as
**`4f/OD-8`** (controller-authored), and this revision's only action is the sweep making the
document's own voice claim exactly what the grammar establishes — provenance, session
integrity, contemporaneity, pair identity, non-reuse, toolchain/hardware identity — and
nothing more (title, Status, problem statement, D2; the ratification quote itself is kept
verbatim as the human's words, not adopted as the document's assertion). The objection is
preserved verbatim above because this mission keeps the arguments that lost — and this one,
twice filed, may yet win at OD-8.

---

*Verification Log rows: 36 (P1–P36; P26/P27/P29 controller-measured iteration 47, corroborating
greps designer-re-run; P32–P34 added in revision round 4 — P33/P34 controller-measured, P32's
greps designer-re-run against the round-3 text; P35–P36 added in revision round 5,
designer-measured first-party — P35 against the round-4 text before any round-5 edit). Named
mutations: 16 — the 10 originals retained (`MUT-REC-SILENT-DEFAULT`,
`MUT-REC-STALL` (reshaped to leg 3), `MUT-EDIT-RAW-NUMBER`, `MUT-CLAIM-NONPARENT`,
`MUT-CLAIM-TOOLCHAIN-SPLIT`, `MUT-CLAIM-ORPHAN-NUMBERS`, `MUT-LEGACY-FOURTH`, `MUT-CHECK-NULL`,
`MUT-AB-FLOOR-SPLIT`, `MUT-PROBE-CALLER-DIR`) plus 5 branch-A additions (`MUT-PAIR-ID-SPLIT`,
`MUT-PAIR-TWO-SESSIONS`, `MUT-CONTROL-REUSE`, `MUT-PAIR-SEQUENTIAL`, `MUT-PAIR-INLINE-BUILD`)
plus 1 round-5 addition (`MUT-REC-UTIL-STALL`);
classification: 6 HARNESS, 10 EVIDENCE.*

***Honest production-vs-test-probe tally*** *(counting convention stated because this mission has
mis-tallied it twice): a PRODUCTION discriminator is a check that lives in shipped code and fires
on every CI run forever; a TEST-SIDE probe fires once at review. Production discriminators —
counting R4's limbs individually because each has its own named failure and its own mutation:
R1 (with the schema pin, invocation grammar, and `binary_sha256` requirement), R2, R3, R4a, R4b,
R4c (round 4: section-local, two clauses — derivation + role-binding — counted as ONE limb, not
padded per clause), R4d (round 4: cardinality including the unpairable named RED — likewise ONE
limb), R4e, R4f, R5, R6 = **11 checker rules**, plus the recorder's all-or-nothing refusal
family (counted as ONE discriminator, not padded per refusal cause — round 5's widened
bounding coverage folds into this same family and does NOT move the count) and D4 step 2's
specific
`hw.ncpu` refusal assertion = **13 production discriminators**. Test-side probes: the **16 named
mutations** above, each one-run-plus-restore at review. The 16 prove the 13 can fail; the 13 are
what stand guard afterwards.*

~~***This branch-A revision authorizes implementation: OD-6 is RATIFIED (attended, 2026-08-03)
and milestones BC.A′/BC.B′ are routable to a planner and executor.***~~ **SUPERSEDED
(round 5): this document is PARKED on `4f/OD-8`** — the gap between the OD-6 stamp's wording
("mechanically valid cost claims") and what branch A delivers (mechanically complete,
contemporaneous, tamper-evident EVIDENCE for a cost claim) is the human's to resolve.
**NOT ROUTABLE until `4f/OD-8` is answered.** OD-6 itself remains RATIFIED (attended,
2026-08-03); nothing in round 5 reopens it, and the round-5 mechanism and claim-accuracy
fixes above stand regardless of how OD-8 resolves.
