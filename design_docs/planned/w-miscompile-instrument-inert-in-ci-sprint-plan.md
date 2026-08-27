# Sprint plan — `P44` (queue row 44, `w-miscompile-instrument-inert-in-ci`)

**Design doc**: [`w-miscompile-instrument-inert-in-ci.md`](w-miscompile-instrument-inert-in-ci.md)
(790 lines; round 1 BLOCKED 3/3 → revision 1; round 2 BLOCKED 3/3 on one shared objection →
narrow-refinement carve-out, `gpt5-6-sol`'s `proposed_fix` applied verbatim, demanded verification
measured as **V19**). **§Quorum round 2 is the most recent state of the design and wins over any
earlier section.** **THE DESIGN DOC WINS** over this plan everywhere except the five places named
in §1, each of which carries the planner's own first-party measurement of why the doc's literal
text would red — or fail to red — a *correct* landing.

**Planner**: `claude-opus-5`, mission-control iteration 133.
**Worktree**: `/Users/voightkampff/dev/sunholo-data/.planner-wt-iter133`, detached at **`aa0ab05`**
(= `80c2bd2` + the design-doc commit; the three files this sprint touches are **byte-identical**
between the two commits — planner-measured, §0.2).
**Estimate**: **~0.35 day.** Three files modified, none created. **~115 changed lines**
(`run.sh` +57/−6 measured on the planner's candidate; `ci.yml` +7/−4; `toolchain_pin_gate_test.go`
+~45 appended plus ~3 reworded doc-comment lines). **Zero existing assertion bodies touched.**
**Risk**: **MEDIUM.** Not on the typing — every proposed byte has been built and exercised by the
planner (§0.2). The risk is entirely in **(a)** the one property whose failure is unrecoverable
without a revert (the next push to `dev` must stay green — §4 RD-1, §7 R1), **(b)** three
acceptance steps that are green-at-base while measuring nothing, and **(c)** one acceptance step
the design writes in a form that **contradicts its own shipped mechanism** (§1 D3).

---

## 0. Rules of engagement for the executor (read first — these override habit)

1. **NO git write operations.** No `add`, `commit`, `stash`, `checkout`, `branch`, `worktree`,
   `mv` (git's), `restore`, `reset`. Git **reads** are required and allowed: `git show`,
   `git status`, `git diff`. The controller commits; the sprint hands off an **uncommitted** tree.
2. **`git status --porcelain` is NEVER empty during this sprint** and that is correct by
   construction. Do not treat a non-empty porcelain as a defect and do not "clean" it.
3. **Restores during the mutation drill use `cp` from `/tmp/p44-backup/`, NEVER
   `git checkout --`.** The whole deliverable is uncommitted; `git checkout --
   design_docs/verification/w-race-gate-blindspot/run.sh` would **delete this sprint's largest
   edit** with no undo. §5.1 is the only restore recipe.
4. **`export AILANG_BIN=/tmp/ailang-v0300/ailang` before any whole-package Go test.**
   Planner-measured, and this is finding **D1**: `go test ./host/verifygate/ -count=1` is **rc=1
   at base** with **17 `--- FAIL`** lines, every one reading `AILANG_BIN is unset`. With the export
   it is **rc=0 in 48.4 s**, 31 `--- PASS`, 0 FAIL. The design doc's AC7 omits the export. **A red
   whose failure text says `AILANG_BIN is unset` is NOT this diff.** Do not "fix" it.
5. **Never read `rc` alone on a `-run`-scoped Go test.** A `-run` selector matching nothing exits
   **0** with `[no tests to run]` — measured at base for AC1 (§0.1). Charter row 33 exists for
   this. Every scoped run below is judged by **counting `=== RUN` / `--- PASS` / `--- FAIL` from
   `-v` output**, with a nonsense-pattern control in the same breath.
6. **Exit codes through pipes lie, and `${PIPESTATUS[0]}` is EMPTY under zsh.** Capture directly —
   `cmd > /tmp/f 2>&1; rc=$?` — then read the file back. In every two-arm comparison print both
   `rc` values side by side and assert they DIFFER where they should.
7. **An empty search result is a claim, not a fact.** Every grep whose answer is a zero carries a
   known-positive control **in the same call**.
8. **`./scripts/verify_go.sh` legitimately emits exactly 2 `WARNING: DATA RACE` lines** from its
   own known-positive control in `design_docs/verification/w-race-gate-blindspot/racecontrol/`,
   under the banner `── race-detector known-positive control` (`verify_go.sh:226`, `:229`, `:232`).
   Three prior iterations have had to attribute these. **They are not a finding. Do not report
   them.**
9. **Sandbox honesty.** The executor runs under a `workspace-write` sandbox that DENIES loopback
   socket binds. Nothing in this sprint binds a socket, so no step is expected to be
   `UNINFORMATIVE UNDER SANDBOX` on that account — but `run.sh` needs **outbound** network only
   when a toolchain is absent from `GOMODCACHE`. On this rig all seven are cached and a warm
   attended run is **3.3 s** (§0.2). If a probe prints `SKIPPED (toolchain unavailable: …)` for a
   toolchain other than the deliberately-fake `go1.99.99`, **report the arm as UNINFORMATIVE
   UNDER SANDBOX and stop** — do not report a verdict. The controller re-runs it outside.
10. **Rows 31, 45, 48, 49, 50 and 51 are separate open items on neighbouring surfaces. Do not
    touch them.** Row 45 edits `normalizeToolchainPin` near the **top** of
    `host/verifygate/toolchain_pin_gate_test.go`; this sprint appends at the **end** — disjoint
    hunks, but flagged here for whichever lands second. Do not refresh any `ci.yml:172` citation:
    the census is **row 31's**, and this sprint's own `ci.yml` edit renumbers the file again
    (design §Deferred Scope; row 43's P7: *line numbers rot, literals do not*). The only
    exceptions are the ~3 comment lines inside the doc-comment M2 must touch anyway.
11. Do not touch `tools/launchd/*`, any `.ail` file, `scripts/verify_go.sh`,
    `scripts/verify_ail.sh`, `host/store/toolchain_canary_test.go`,
    `design_docs/verification/w-race-gate-blindspot/repro/**` or `.../racecontrol/**`. **No `.ail`
    byte changes: this is a single-gate (Go) sprint and `./scripts/verify_ail.sh` need not run.**

---

## 0.1 The baselined acceptance table — every command run on the PRISTINE tree at `aa0ab05`

`VERIFIED BY CONTROLLER` rows were re-derived first-party by the planner and agreed exactly.
Rows marked `PLANNER` are the planner's own measurements. **Read the "how it is READ" column: three
of these are green at base while measuring nothing, and one is red at base for a reason that is not
this diff.**

| AC | Command (verbatim) | **BASE result at `aa0ab05`** | Expected POST result | How it is READ |
|---|---|---|---|---|
| **AC1** | `go test ./host/verifygate/ -run '^TestMiscompileInstrumentStepIsGatedInCI$' -count=1 -v` | **rc=0**, `ok … 0.172s [no tests to run]`, `=== RUN`=**0**, `--- PASS`=**0** — *green while measuring nothing* | `=== RUN`=1, `--- PASS`=1, `--- FAIL`=0 | **Counted lines, never rc.** Paired control: `-run '^TestNoSuchWiringGateZZZ$' -v` → `no tests to run` count **2**, **rc=0** (base-measured) |
| **AC2** | block-scoped `awk` (§4 M2) then `grep -c 'continue-on-error'` on the block; same call `grep -c 'design_docs/verification/w-race-gate-blindspot/run.sh' .github/workflows/ci.yml` | block count **1**; path control **1**; (repo-wide `continue-on-error` in `ci.yml` = **1**, recorded not asserted) | block count **0**; path control **1** | Two numbers in one call. The **path control is the known-positive** proving the awk found the step and the grep is live. A `continue-on-error` on any OTHER step must NOT fail this — round-2 R1 forbids a repo-wide ban |
| **AC3** | `diff <(git show aa0ab05:host/verifygate/toolchain_pin_gate_test.go \| awk '/^func TestMiscompileInstrumentProbesPinnedToolchain/,/^}/') <(awk '…' host/verifygate/toolchain_pin_gate_test.go)`; same for `TestReproModuleFloorStaysBelowKnownBadToolchains`; then the 4-test scoped run | both diffs **empty by identity**; scoped run **rc=0**, `=== RUN`=**4**, `--- PASS`=**4**, `--- FAIL`=**0**, `ok … 0.208s` | both diffs still **empty**; `=== RUN`=4, `--- PASS`=4 | Diff emptiness + counted lines. **Green at base by identity** — its teeth are drill arm **X7** |
| **AC4(a)** | `./design_docs/verification/w-race-gate-blindspot/run.sh` | **rc=0**, `host: darwin/arm64   default toolchain: go1.26.4`, 4× `BUG: Field="" want "stateRoot"`, 3× `OK`, `-N`→OK / `-l`→BUG, banner `RESULT: reproduction confirmed, and both controls fired.` **7/7 probes ran, 0 SKIPPED.** Warm wall time **3.3 s** | **rc=0** and the **byte-identical** darwin banner | rc captured directly + `grep -c` on the banner. Planner proved the darwin banner output is **byte-identical** between base and candidate on the P8 profile (§0.2) |
| **AC4(b)** | guard-trip — **the design's text is WRONG here; see §1 D3.** Corrected form in §4 M1 | n/a (mutation) | rc=1 naming the **refuse arm**, `INSTRUMENT FAILURE: no verified platform contract for darwin/arm64` | both rc printed side by side; `grep -c` on the named text; sha256 restore |
| **AC4(c)** | **referenced by the doc's mutation table but never DEFINED in §Acceptance criteria.** Defined here, §4 M1 | n/a (mutation) | rc=1 naming `INSTRUMENT FAILURE (PLATFORM ALARM)` | same |
| **AC5** | `mv design_docs/verification/w-race-gate-blindspot/repro/main.go{,.MUT}` → run `run.sh` → restore | rc=1 via floor 1 (`no toolchain ran at all`). Planner-measured **on the candidate** (floor 1 is byte-identical to base, so base behaves the same): rc=1, floor-1 text **1**, coverage-floor text **0** | identical: rc=1, floor-1 text 1, coverage-floor text 0 | rc + two `grep -c`s. The point is *which* floor fires: floor 1 must fire **before** the new coverage floor |
| **AC6** | `gh run view <run> --log \| grep -c 'RESULT: linux/amd64 clean'` on the landing PR's check and the first post-merge `dev` run | on run `33069999131`: step conclusion `success` **with** `INSTRUMENT FAILURE`=1, `RESULT:`=0, control banner=1 (the contradiction this item removes) | step `success`, linux RESULT sentence=1, `INSTRUMENT FAILURE`=0, control=1 | **CONTROLLER-OWNED — the executor CANNOT close this.** No PR exists at sprint time and the executor makes no git writes. Report it **NOT ATTEMPTED (controller-owned)**, never a verdict |
| **AC7** | `bash -n …/run.sh`; `stat -f '%Sp' …/run.sh`; `gofmt -l host/verifygate/`; `go vet ./host/verifygate/`; **`AILANG_BIN=/tmp/ailang-v0300/ailang` `go test ./host/verifygate/ -count=1`** | rc=**0**; `-rwxr-xr-x`; **empty**; rc=**0**; **rc=0 in 48.4 s** WITH the export — and **rc=1 with 17 `--- FAIL`/`AILANG_BIN is unset` WITHOUT it (D1)** | all identical, plus `--- PASS` ≥ 32 (31 base + the new test) | rc + counted lines. **Never run the package test without the export** |
| **RD-1** | the constraint-1 replay — §4 M1, planner-defined and planner-measured | n/a (the revised floors do not exist at base) | linux profile → **rc=0** + the linux RESULT sentence; four negative controls RED | **The one property whose failure is unrecoverable without a revert.** BLOCKING gate on M1 |

---

## 0.2 What the planner measured on top, and why each is load-bearing

Everything below was measured first-party in this worktree at `aa0ab05`, today.

- **The proposed `run.sh` was BUILT, not imagined.** The planner generated the exact bytes the
  design's §M1 specifies (all six edits) into `/tmp/p44-cand/run.sh` — **156 lines, +57/−6 against
  base, `bash -n` rc=0** — and ran it, mutated it and replayed it. Every claim in this plan about
  post-landing behaviour is a reading off those bytes, not a prediction. **The executor must
  regenerate these bytes itself and re-run the same steps; do not copy a result from this plan.**
- **RD-1, the constraint-1 property, RE-DERIVED.** The design discharges "the next push must not
  red `dev`" at P17 by *replaying* P5's measured linux profile. The planner did the replay against
  the candidate's own floor block: `ran=7 ran_bad=4 bad_expected=4 saw_bad=0 saw_good=1
  saw_pinned_ok=1 host_pair=linux/amd64 expect_defect=0` → **rc=0**, printing
  `RESULT: linux/amd64 clean — no KNOWN-BAD toolchain reproduced here, matching / the iteration-46
  AC6 measurement; all 4 known-bad and 7 total / probes ran, and the known-good and pinned
  (go1.26.6) toolchains both reported OK.` Three negative controls in the same harness: a
  3-of-4 partial → **rc=1** coverage floor; a linux `saw_bad=1` → **rc=1** PLATFORM ALARM; and the
  P8 darwin profile → **rc=0** with output **byte-identical to base**.
- **The platform probe was exercised on all five arms** with a `uname` shim on `PATH`:
  `Linux`/`x86_64` → `linux/amd64 expect_defect=0` (**this is the ubuntu-24.04 runner, P17's three
  in-run reports**); `Darwin`/`arm64` → `darwin/arm64 expect_defect=1`; `Linux`/`aarch64` →
  **refuse, rc=1**; `Darwin`/`x86_64` → **refuse, rc=1**; `FreeBSD`/`riscv64` → **refuse, rc=1**.
  Override immunity confirmed: with `GOOS=windows GOARCH=amd64` exported the block still computes
  `darwin/arm64` on this rig (P16's mechanism claim, re-derived).
- **Every bind of the two EXISTING static tests survives the candidate**, re-derived as shell
  equivalents: exact-shebang count **1**; column-0 `KNOWN_BAD="`/`KNOWN_GOOD="`/`PINNED="` **1**
  each; `saw_pinned_ok` sites **3**; literal `INSTRUMENT FAILURE: the PINNED toolchain` **1**;
  `KNOWN_BAD`/`KNOWN_GOOD` values byte-unchanged. The new top-level names (`ran_bad`,
  `bad_expected`, `host_os`, `host_arch`, `host_pair`) collide with **no** name
  `shellAssignmentValues` looks for (row-50 deference holds).
- **The new test's own run.sh assertions re-derived against the candidate**: `uname -s` **1**,
  `uname -m` **1**; **executable** `go env GOOS`/`go env GOARCH` after stripping from the first
  `#` → **0**, with the known-positive control in the same call showing the token IS present at
  candidate `:38` **inside a comment** — which is precisely the round-1 form's failure (V19 arm A)
  and precisely what round 2's scoped form allows. `host_pair` + `go env` on one code line → **0**.
- **Three drill arms rehearsed end-to-end against the candidate on this darwin rig**, each with its
  landed-proof and a sha256 restore: **X3** (`mv repro/main.go`) → rc=1, floor-1 text **1**,
  coverage-floor text **0** (floor 1 fires first, as designed); **X6** (append `go1.99.99` to
  `KNOWN_BAD`) → rc=1 with `INSTRUMENT FAILURE: 4 of 5 KNOWN-BAD probes completed`; **X5**
  (`chmod -x`) → **rc=126**. `go/version.IsValid("go1.99.99")` → **true** and the oldest KNOWN_BAD
  stays `go1.26.0`, so X6's static layer correctly stays GREEN as the doc predicts.
- **The base commit question is closed.** `git diff --stat 80c2bd2 aa0ab05` → **one file, the
  design doc, +790**. The three sprint files are byte-identical across the two commits, so every
  `git show 80c2bd2:…` in the design doc is runnable **verbatim**; this plan uses `aa0ab05` for
  clarity and the two are interchangeable here.
- **No other artefact binds what this sprint moves.** Repo-wide: nothing in `*.go`/`*.sh`/`*.yml`
  reads the step **name** (only prose in `world-mission.md`, the log, and
  `w-race-gate-blindspot.md:208` — all true history, all left alone).
  `TestZ3PinDeclaredOnceAndInstalledInBothJobs` (`ail_binary_gate_test.go:668`) reads `ci.yml` for
  `Z3_VER:`/`Z3_SHA:` counts, the two literal install lines and the controls `ailang-verify:` /
  `go-verify:` / `./scripts/verify_go.sh` — **all untouched**. `TestNoCIStepOrScriptReachesThe`
  `PublishEntrypoint` (`runbook_stageb_test.go:329`, AC30) scans `ci.yml` + `scripts/*.sh` for
  `world-publish`: the proposed comment and step contain **none** of that token, and its
  known-positive control `verify_go.sh` in `ci.yml` still fires. It does **not** scan
  `design_docs/**`, so `run.sh` is unconstrained by it.
- **`TestNoRigAbsolutePaths`** (`ail_binary_gate_test.go:552`) globs `host/verifygate/*.go` and
  errors on `/tmp/ailang`, `/Users/`, `/home/runner/` (needles assembled at runtime). The file M2
  appends to **is scanned**. The design's sketch contains none — **keep it that way**: cite
  provenance as `P16`/`P17`/`V19` and the doc path, **never** a rig path.
- **P15 re-derived**: `scripts/verify_go.sh:258` is `go test ./... -count=1` from the repo root,
  with the `-race` leg at `:262`. A `host/verifygate` red **is** a CI-gate red.
- **The new test needs no new imports**: `os` (`:5`), `path/filepath` (`:6`), `strings` (`:9`) are
  already in the file's import block. `repoRoot` is a package-level var at
  `ail_binary_gate_test.go:27` in `package verifygate`.
- Cost centres on this rig: attended `run.sh` **3.3 s warm** (all 7 toolchains present in
  `GOMODCACHE`); `go build ./...` **1.7 s**; scoped verifygate test **~0.2 s**; whole
  `host/verifygate` package **48.4 s**.

---

## 1. Design-doc defects the planner found — five. **Two would have shipped a false verdict.**

Contradicting the designer here is the loop working, not failing.

### D1 — **AC7's whole-package clause is RED AT BASE; the design omits the required export** · severity **HIGH**

The doc's AC7 reads `go test ./host/verifygate/ -count=1 → rc=0` and its Base note says only
*"full-package run happens at sprint time (it runs the verify_ail shim arms; ~60 s)"* — i.e. the
designer **did not run it**. Measured verbatim at `aa0ab05`:

```
go test ./host/verifygate/ -count=1        →  rc=1,  17 × "--- FAIL",
                                              every message: "AILANG_BIN is unset — the shim arms
                                              need the pinned released delegate to run the real gate"
AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -count=1
                                           →  rc=0 in 48.4 s, 31 --- PASS, 0 --- FAIL, 0 --- SKIP
```

**Concrete failure avoided**: an executor running AC7 as written reports a 17-failure RED against
a correct landing, and the plausible next move — "make the failing tests skip" — is the exact
false-green class this whole mission line exists to kill (those tests `t.Fatalf` **on purpose**
rather than skip). **Resolution**: §0 rule 4 and every AC7/M3 invocation in this plan carry
`AILANG_BIN=/tmp/ailang-v0300/ailang`. The scoped runs (AC1, AC3) need **no** export — they are
pure text scans — and are left unexported so their lane-independence stays visible.

### D2 — **AC6 cannot be closed by the executor at all** · severity **HIGH (process)**

AC6 requires reading *"the landing PR's check run AND the first dev run after merge"*. At sprint
time there is no PR, no merge, and the executor is forbidden every git write (§0 rule 1). The
design files it as an acceptance criterion alongside six the executor **can** run.

**Concrete failure avoided**: an executor either fabricates a verdict from the *pre-landing* run
`33069999131` (which is the base contradiction, not the post state) or reports the sprint
incomplete. **Resolution**: AC6 is marked **CONTROLLER-OWNED** in §0.1 and §4 M3. The executor
reports it `NOT ATTEMPTED (controller-owned)` and nothing else. Its executor-side proxy — the
thing AC6 would have measured about the *script* — is **RD-1**, which the executor **does** run.

### D3 — **AC4(b)'s guard-trip is a revision-0 leftover that contradicts the shipped mechanism, and asserts text that cannot appear** · severity **HIGH** · MEASURED

AC4(b) says: *"mutate `if [ "$host_pair" = "darwin/arm64" ]` to `"darwin/amd64"` (every platform
in clean-mode), run → rc=1 whose output names the PLATFORM ALARM text"*.

Two things are wrong. First, **there is no `if [ "$host_pair" = … ]` anywhere in the shipped
design** — quorum R1's fail-closed `case` replaced it, and `case` has no "clean-mode default":
that is the whole point of the refuse arm. Second, and decisive, the planner **ran it**:

```
mutation:  darwin/arm64) expect_defect=1 ;;   →   darwin/amd64) expect_defect=1 ;;
landed:    grep -c 'darwin/amd64) expect_defect=1' → 1 (was 0);  'darwin/arm64) …=1' → 0 (was 1)
result:    rc=1
           "INSTRUMENT FAILURE: no verified platform contract for darwin/arm64"
           grep -c 'PLATFORM ALARM'                    → 0      ← the doc predicts ≥ 1
           grep -c 'no verified platform contract'     → 1
```

**Concrete failure avoided**: an executor asserting the PLATFORM ALARM text under AC4(b) reds a
**correct** landing, and the natural repair — editing the shipped `case` until the alarm text
appears — would dismantle R1's fail-closed arm. The design's own §Named RED mutations row **M2**
already states the correct outcome (*"attended darwin RED via the fail-closed `*)` arm …
rehearsed as AC4(c)"*); AC4(b) simply was not re-derived when the mechanism changed. **This is
the design doc committing, in its acceptance list, the exact defect its §"What the round-2
revision deliberately STOPPED claiming teeth over" paragraph names: a cell that no longer
predicts the right result is worse than no cell.**

**Resolution.** Split into two rehearsals, both planner-measured on the candidate:

- **AC4(b) — the refuse arm.** Mutation as the doc writes it; assert **rc=1** and
  `no verified platform contract for darwin/arm64` = 1, and **`PLATFORM ALARM` = 0**.
- **AC4(c) — the PLATFORM ALARM, newly DEFINED here** (the doc references AC4(c) in mutation row
  M2 but never defines it). The alarm is reachable on darwin only by flipping the darwin arm's
  **polarity**, not its label:

```
mutation:  darwin/arm64) expect_defect=1 ;;   →   darwin/arm64) expect_defect=0 ;;
landed:    grep -c 'darwin/arm64) expect_defect=0' → 1 (was 0)
result:    rc=1
           "INSTRUMENT FAILURE (PLATFORM ALARM): a KNOWN-BAD toolchain reported BUG on
            darwin/arm64 — measured NOT affected (iteration-46 AC6). …"
           grep -c 'PLATFORM ALARM' → 1
```

### D4 — **the design's network-flake framing does not describe the LOCAL rehearsals** · severity **LOW**

§Declared residual 6 and §Quorum round 1 price the gate on network dependence ("six of seven
probes"). True of CI. **Not true of this rig**: all seven toolchains are already in `GOMODCACHE`
(`golang.org/toolchain@v0.0.1-go1.{24.9,25.6,26.0,26.3,26.5,26.6}.darwin-arm64` + the default
`go1.26.4`), and a warm attended `run.sh` is **3.3 s**, not minutes. **Resolution**: the drill's
darwin arms are cheap — run them all, do not skip on cost. Only drill arm **X6** (`go1.99.99`)
needs an outbound round-trip, and it SKIPs identically whether the answer is
`toolchain not available` or a DNS failure, so its coverage-floor RED is robust to a
network-denied sandbox. §0 rule 9 is the honest reading rule if any *other* probe SKIPs.

### D5 — **the `ci.yml` comment forward-references a path that does not exist yet** · severity **LOW, land it anyway**

The proposed comment ends `# Design: design_docs/implemented/w-miscompile-instrument-inert-in-ci.md`
while the doc is at `design_docs/planned/…` until the controller moves it on landing. Nothing
binds the string (measured: no test reads it). **Resolution**: land it **exactly as written** —
this matches the P43 precedent and the controller's move closes it. Flagged so the executor does
not "correct" it to `planned/` and so the controller knows the move is load-bearing for one
comment line. *(Trivial adjacent nit, non-blocking: the new test's
`filepath.Join(repoRoot, miscompileReproducerPath)` joins a slash-path directly where
`module_manifest_gate_test.go:54` uses `filepath.FromSlash`. Identical on darwin and linux; the
repo has both idioms. Land the doc's form.)*

### Not defects — recorded so nobody downstream re-derives them as findings

- **AC1's base rc=0 with `[no tests to run]`** is the sprint's sharpest fact, not an anomaly.
- **AC3's diff is empty at base by identity.** By design; its teeth are drill arm X7.
- **`verify_go.sh`'s two `WARNING: DATA RACE` lines** are its own known-positive control (§0 rule 8).
- **`Observatory: NNNMB (warn threshold: 200MB)`** on stderr from every `ailang` invocation is
  normal.
- **`GOTOOLCHAIN` appears twice in `ci.yml`** and must stay twice — it is
  `TestGoToolchainPinsAgreeAndMatchJobList`'s cardinality pin, untouched here.

---

## 2. Milestone → acceptance-criterion mapping

The design defines two milestones (M1, M2). This plan keeps their IDs and numbers, and adds two
non-editing milestones (M3 sweep, M4 drill) matching house practice. **Diffed against the doc's AC
list: every AC1–AC7 is claimed, AC6 is reassigned to the controller with its reason stated (D2),
and AC4 is split into (a)/(b)/(c) with (c) newly defined (D3). No other disagreement found.**

| Milestone | Files | ACs it CLOSES | ACs it re-confirms |
|---|---|---|---|
| **M1** — `run.sh` learns the platform truth (still non-gating; zero CI risk) | `design_docs/verification/w-race-gate-blindspot/run.sh` | **RD-1**, **AC4(a)(b)(c)**, **AC5**, AC7's `bash -n` + exec-bit clauses | AC3 (untouched by construction) |
| **M2** — make the refusal binding, salt the flag's ground | `.github/workflows/ci.yml`, `host/verifygate/toolchain_pin_gate_test.go` | **AC1**, **AC2**, **AC3** | AC7 |
| **M3** — full acceptance sweep, one recorded pass | none | **AC7** (all clauses) | AC1–AC5, RD-1 |
| **M4** — mutation drill X1–X8 + green control | none net (every arm restored) | — | the drill's own arms |
| *(post-landing)* | — | **AC6** — **CONTROLLER-OWNED** (D2) | — |

**Why M1 precedes M2, and why RD-1 sits inside M1.** M1 changes only a step that is still
swallowed by `continue-on-error: true`, so it carries **zero CI risk**. M2 removes the swallow.
**RD-1 — the proof that the revised floors exit 0 on the measured linux profile — must be GREEN
before M2 is started.** Doing M2 first, or skipping RD-1, is the one ordering whose failure mode
is a red `dev` recoverable only by revert.

---

## 3. The exact file touch set

| # | Path | Action | Changed lines | Milestone |
|---|---|---|---|---|
| 1 | `design_docs/verification/w-race-gate-blindspot/run.sh` | MODIFY — 6 edits (§4 M1) | **+57 / −6** (planner-measured on the candidate) | M1 |
| 2 | `.github/workflows/ci.yml` | MODIFY — the miscompile step block only | **+7 / −4** | M2 |
| 3 | `host/verifygate/toolchain_pin_gate_test.go` | MODIFY — append 1 const + 1 test func at EOF; reword ~3 lines **inside** the `TestMiscompileInstrumentProbesPinnedToolchain` doc-comment | **+~45 / −3** | M2 |

**Nothing else.** No file created, none deleted. Zero `.ail`. Zero existing assertion bodies.

---

## 4. Milestones, ordered

### M1 — `run.sh` learns the platform truth

Six edits, all specified verbatim in the design doc §Milestones M1. **Transcribe from the design
doc's fenced blocks; do not retype from memory.** The file indents with **TABs** — every existing
block body uses `\t`, and the new `case` arms and floor bodies must too.

1. **After `ran=0` (`:32`)**, insert the declarations + platform-probe + polarity blocks. Keep the
   two new counter assignments at **column 0** (row-50 deference); the `case`-internal
   `host_os=`/`host_arch=`/`expect_defect=` arms are necessarily TAB-indented, like every other
   conditional body in the file.
2. **In `probe()`, immediately after `local tc="$1" expect="$2" out rc bin` (`:35`)**, add
   `[ "$expect" = "BAD" ] && bad_expected=$((bad_expected + 1))`. A SKIP still counts here — that
   is the point.
3. **Immediately after `ran=$((ran + 1))` (`:46`)**, add
   `[ "$expect" = "BAD" ] && ran_bad=$((ran_bad + 1))`.
   *(`set -e` is absent — `:20` is `set -uo pipefail` — so the `&&` one-liner form is safe.
   Planner-verified: `bash -n` rc=0 and both counters read correctly at runtime.)*
4. **Replace the banner read at `:61`** with
   `echo "host: $host_pair   default toolchain: $(go version | awk '{print $3}')"`. This removes
   the file's last read of the overridable channel and prints the identical string on both
   verified platforms.
5. **Replace `:84-89`** (the unconditional `saw_bad` floor) with the three-floor block: coverage
   floor, then `expect_defect==1 && saw_bad==0`, then `expect_defect==0 && saw_bad!=0`. **Floors
   `:80-83` and `:90-100` stay byte-identical.**
6. **Split the banner at `:101-104`** per platform. The darwin branch's four lines are
   **byte-identical** to today's (planner-verified: base tail and candidate tail produce identical
   output on the P8 profile), so `OUTPUTS.md` §5's transcript stays truthful. The linux branch is
   deliberately a **different** sentence so the mission's historical greps
   (`RESULT: reproduction confirmed` = 0 on pre-landing runs) stay meaningful.

Then, in order:

```bash
# ---- hygiene first: a syntax error makes every reading below meaningless
bash -n design_docs/verification/w-race-gate-blindspot/run.sh; echo "bash -n rc=$?"   # 0 (base: 0)
stat -f '%Sp' design_docs/verification/w-race-gate-blindspot/run.sh                    # -rwxr-xr-x
```

#### RD-1 — **the constraint-1 re-derivation. BLOCKING. Do not start M2 until this is green.**

> The change must not red `dev` on the next push. The design discharges this at P17 by replaying
> the measured linux profile through the revised floors. **RE-DERIVE it here; do not inherit it.**
> This is the only property in the sprint whose failure is unrecoverable without a revert.

```bash
R=design_docs/verification/w-race-gate-blindspot/run.sh
TAIL=$(grep -n '^if \[ "\$ran" -eq 0 \]; then' "$R" | cut -d: -f1)
[ -n "$TAIL" ] || { echo "RD-1 INSTRUMENT FAILURE: floor block not located"; exit 1; }

replay() {  # $1=outfile, rest = variable assignments
  out=$1; shift
  { echo '#!/usr/bin/env bash'; echo 'set -uo pipefail'; echo 'PINNED="go1.26.6"'
    for v in "$@"; do echo "$v"; done
    tail -n +$TAIL "$R"; } > /tmp/p44-rd1.sh
  bash /tmp/p44-rd1.sh > "$out" 2>&1; echo $?
}

# (a) THE LOAD-BEARING ARM — P5's measured linux/amd64 profile, verbatim
rc_a=$(replay /tmp/rd1-a.out ran=7 ran_bad=4 bad_expected=4 saw_bad=0 saw_good=1 \
               saw_pinned_ok=1 'host_pair="linux/amd64"' expect_defect=0)
# (b) P8's measured darwin/arm64 profile
rc_b=$(replay /tmp/rd1-b.out ran=7 ran_bad=4 bad_expected=4 saw_bad=1 saw_good=1 \
               saw_pinned_ok=1 'host_pair="darwin/arm64"' expect_defect=1)
# (c) NEGATIVE CONTROL — one known-bad probe SKIPPED on linux
rc_c=$(replay /tmp/rd1-c.out ran=6 ran_bad=3 bad_expected=4 saw_bad=0 saw_good=1 \
               saw_pinned_ok=1 'host_pair="linux/amd64"' expect_defect=0)
# (d) NEGATIVE CONTROL — the defect SPREADS to linux
rc_d=$(replay /tmp/rd1-d.out ran=7 ran_bad=4 bad_expected=4 saw_bad=1 saw_good=1 \
               saw_pinned_ok=1 'host_pair="linux/amd64"' expect_defect=0)

echo "RD-1 rc: a=$rc_a b=$rc_b c=$rc_c d=$rc_d"     # EXPECT  a=0 b=0 c=1 d=1
grep -c 'RESULT: linux/amd64 clean'      /tmp/rd1-a.out   # 1   ← dev stays GREEN
grep -c 'INSTRUMENT FAILURE'             /tmp/rd1-a.out   # 0
grep -c 'RESULT: reproduction confirmed' /tmp/rd1-b.out   # 1
grep -c 'KNOWN-BAD probes completed'     /tmp/rd1-c.out   # 1
grep -c 'PLATFORM ALARM'                 /tmp/rd1-d.out   # 1
```

**Read `a=0 b=0` AND `c=1 d=1` together.** `a=0` alone is equally consistent with a floor block
that cannot fail at all; the two negative controls are what make `a=0` a measurement. If
`a` is anything but 0, **STOP — do not proceed to M2** and report; landing M2 on top of it reds
`dev` on the next push.

**Second half of RD-1 — the platform probe resolves to the LISTED arm on the real runner.** P17
establishes from run `33069999131`'s own log (`Image: ubuntu-24.04`; setup-go cache key
`setup-go-Linux-x64-ubuntu24-…`; `AILANG v0.30.0 on Linux/x86_64`; `go1.26.6 linux/amd64` ×7) that
`uname -s`/`uname -m` there give `Linux`/`x86_64`. Re-derive the mapping locally with a shim:

```bash
d=$(mktemp -d)
printf '#!/bin/sh\ncase "$1" in\n-s) echo Linux ;;\n-m) echo x86_64 ;;\nesac\n' > "$d/uname"
chmod +x "$d/uname"
HS=$(grep -n '^# Platform probe (row 44' "$R" | cut -d: -f1)
HE=$(grep -n 'no verified platform contract' "$R" | head -1 | cut -d: -f1)
{ echo 'set -uo pipefail'; sed -n "${HS},$((HE+1))p" "$R"
  echo 'echo "host_pair=$host_pair expect_defect=$expect_defect"'; } > /tmp/p44-head.sh
PATH="$d:$PATH" bash /tmp/p44-head.sh; echo "rc=$?"   # host_pair=linux/amd64 expect_defect=0, rc=0
bash /tmp/p44-head.sh; echo "real-kernel rc=$?"       # host_pair=darwin/arm64 expect_defect=1, rc=0
# fail-closed controls — each must rc=1 with the refuse text
for pair in "Linux aarch64" "Darwin x86_64" "FreeBSD riscv64"; do
  set -- $pair
  printf '#!/bin/sh\ncase "$1" in\n-s) echo %s ;;\n-m) echo %s ;;\nesac\n' "$1" "$2" > "$d/uname"
  PATH="$d:$PATH" bash /tmp/p44-head.sh; echo "  $pair rc=$?"    # rc=1 each
done
# override immunity (P16's mechanism claim)
GOOS=windows GOARCH=amd64 bash /tmp/p44-head.sh                  # still darwin/arm64
```

#### AC4 — the attended darwin legs

```bash
# AC4(a).  base: rc=0, 4x BUG, 7/7 ran, 0 SKIPPED, 3.3 s warm
./design_docs/verification/w-race-gate-blindspot/run.sh > /tmp/p44-ac4a.txt 2>&1; echo "rc=$?"  # 0
grep -c 'RESULT: reproduction confirmed, and both controls fired\.' /tmp/p44-ac4a.txt   # 1
grep -c 'BUG: Field="" want "stateRoot"' /tmp/p44-ac4a.txt                              # 5 (4 probes + the -l control)
grep -c 'SKIPPED' /tmp/p44-ac4a.txt                                                     # 0  ← else UNINFORMATIVE (rule 9)
diff <(tail -4 /tmp/p44-ac4a.txt) <(git show aa0ab05:design_docs/verification/w-race-gate-blindspot/run.sh \
  | tail -4 | sed "s/^echo \"//; s/\"$//; s/\$PINNED/go1.26.6/g") && echo DARWIN_BANNER_BYTE_IDENTICAL
```

Take the drill backup **now** (§5.1) — AC4(b)/(c) are mutations of the deliverable.

```bash
# AC4(b) — the REFUSE arm (D3-corrected).  Land-proof BEFORE reading any result.
sed -i '' 's|^\tdarwin/arm64) expect_defect=1 ;;|\tdarwin/amd64) expect_defect=1 ;;|' "$R"
echo "landed: $(grep -c 'darwin/amd64) expect_defect=1' "$R") (want 1, was 0) / \
$(grep -c 'darwin/arm64) expect_defect=1' "$R") (want 0, was 1)"
./"$R" > /tmp/p44-ac4b.txt 2>&1; echo "rc=$?"                                  # 1  (AC4(a) rc was 0)
grep -c 'no verified platform contract for darwin/arm64' /tmp/p44-ac4b.txt     # 1
grep -c 'PLATFORM ALARM' /tmp/p44-ac4b.txt                                     # 0  ← D3: NOT the alarm
cp /tmp/p44-backup/run.sh "$R"; shasum -a 256 "$R" /tmp/p44-backup/run.sh      # two equal hashes

# AC4(c) — the PLATFORM ALARM (newly defined; polarity flip, not label swap)
sed -i '' 's|^\tdarwin/arm64) expect_defect=1 ;;|\tdarwin/arm64) expect_defect=0 ;;|' "$R"
echo "landed: $(grep -c 'darwin/arm64) expect_defect=0' "$R") (want 1, was 0)"
./"$R" > /tmp/p44-ac4c.txt 2>&1; echo "rc=$?"                                  # 1
grep -c 'INSTRUMENT FAILURE (PLATFORM ALARM)' /tmp/p44-ac4c.txt                # 1
cp /tmp/p44-backup/run.sh "$R"; shasum -a 256 "$R" /tmp/p44-backup/run.sh      # two equal hashes
```

#### AC5 — floor 1 still fires first

```bash
mv design_docs/verification/w-race-gate-blindspot/repro/main.go{,.MUT}
(cd design_docs/verification/w-race-gate-blindspot/repro && GOTOOLCHAIN=go1.26.6 go build . 2>&1 | head -1)
#   land-proof: "no Go files in …/repro"
./"$R" > /tmp/p44-ac5.txt 2>&1; echo "rc=$?"                        # 1
grep -c 'no toolchain ran at all' /tmp/p44-ac5.txt                  # 1  ← floor 1
grep -c 'KNOWN-BAD probes completed' /tmp/p44-ac5.txt               # 0  ← coverage floor must NOT fire first
mv design_docs/verification/w-race-gate-blindspot/repro/main.go{.MUT,}
git status --porcelain design_docs/verification/w-race-gate-blindspot/repro/   # empty
```

**M1 is DONE when RD-1 is green on all four arms and the fail-closed shim controls all rc=1.**
Snapshot → `.snap/M1/design_docs/verification/w-race-gate-blindspot/run.sh`.

---

### M2 — make the refusal binding, and salt the flag's ground

**(a) `.github/workflows/ci.yml`** — replace lines `167-172` (the two comment lines, the `- name:`
line, and the `continue-on-error: true` line) with the design's §M2 block verbatim. `timeout-minutes: 15`
and the `run:` line are unchanged. Do **not** renumber or refresh any `ci.yml:172` citation (§0 rule 10).

**(b) `host/verifygate/toolchain_pin_gate_test.go`** — two appends at **EOF** plus one in-comment
reword:

- `const miscompileReproducerPath = "design_docs/verification/w-race-gate-blindspot/run.sh"`
- `func TestMiscompileInstrumentStepIsGatedInCI(t *testing.T)` — the design's §M2 sketch
  **verbatim**, including its doc-comment and DECLARED RESIDUAL paragraph. **This is round-2's
  scoped form**: it inspects only the miscompile step's own block for `continue-on-error`, and
  rejects only **executable** uses of `go env GOOS`/`go env GOARCH` in `run.sh` (comments strip at
  the first `#`). The round-1 repo-wide form **reds on its own documentation** — measured at V19
  arm A and re-derived by the planner (§0.2). Do not "simplify" it back.
- Reword the ~3 stale lines **inside** `TestMiscompileInstrumentProbesPinnedToolchain`'s
  doc-comment at `:186-190` — the `non-gating, continue-on-error: true, ci.yml:172` clause and the
  now-discharged residual sentence *"flipping ci.yml:172 to gating is the named follow-up in
  Deferred Scope, paired with OD-1"*. These are **comment bytes only**; AC3's function-body diff
  extracts from `^func …` to `^}` and therefore stays empty.
- **No new imports** (`os`, `path/filepath`, `strings` already present). **No absolute path
  anywhere in the file** — `TestNoRigAbsolutePaths` scans it (§0.2).

```bash
# AC2 — flag gone from THIS step; step still present.   base: 1 / 1
awk -v p='design_docs/verification/w-race-gate-blindspot/run.sh' '
  /^[[:space:]]*- name:/ { if (has) exit; buf=$0"\n"; next }
  { buf=buf $0 "\n"; if (index($0,p)) has=1 }
  END { if (has) printf "%s", buf }' .github/workflows/ci.yml > /tmp/p44-step.txt
wc -l < /tmp/p44-step.txt                                              # >0, else the awk found nothing
grep -c 'continue-on-error' /tmp/p44-step.txt                          # 0   (base: 1)
grep -c 'design_docs/verification/w-race-gate-blindspot/run.sh' .github/workflows/ci.yml   # 1  KP, firing
grep -c 'continue-on-error' .github/workflows/ci.yml                   # 0   (base: 1) — RECORD, do not assert
grep -c 'timeout-minutes' .github/workflows/ci.yml                     # KP for that zero, must be >=1

# AC1 — the wiring test RUNS and PASSES.   base: rc=0, `=== RUN`=0
go test ./host/verifygate/ -run '^TestMiscompileInstrumentStepIsGatedInCI$' -count=1 -v \
  > /tmp/p44-ac1.txt 2>&1; echo "rc=$?"
grep -c '^=== RUN'  /tmp/p44-ac1.txt      # 1   (base: 0)
grep -c '^--- PASS' /tmp/p44-ac1.txt      # 1   (base: 0)
grep -c '^--- FAIL' /tmp/p44-ac1.txt      # 0
# paired nonsense control, same call shape — proves the counter can read a zero
go test ./host/verifygate/ -run '^TestNoSuchWiringGateZZZ$' -count=1 -v > /tmp/p44-nc.txt 2>&1; echo "rc=$?"
grep -c 'no tests to run' /tmp/p44-nc.txt # 2, and rc=0.  THIS is why rc is never the gate.

# AC3 — the two coupled function BODIES are byte-stable, and the four tests are green
for fn in TestMiscompileInstrumentProbesPinnedToolchain TestReproModuleFloorStaysBelowKnownBadToolchains; do
  diff <(git show aa0ab05:host/verifygate/toolchain_pin_gate_test.go | awk "/^func $fn/,/^}/") \
       <(awk "/^func $fn/,/^}/" host/verifygate/toolchain_pin_gate_test.go) && echo "$fn BODY BYTE-STABLE"
done
go test ./host/verifygate/ -run 'TestMiscompileInstrumentProbesPinnedToolchain|TestReproModuleFloorStaysBelowKnownBadToolchains|TestGoToolchainPinsAgreeAndMatchJobList|TestCanaryDeclaresPositiveArmOnly' \
  -count=1 -v > /tmp/p44-ac3.txt 2>&1; echo "rc=$?"
grep -c '^=== RUN' /tmp/p44-ac3.txt        # 4  (base: 4)
grep -c '^--- PASS' /tmp/p44-ac3.txt       # 4  (base: 4)
```

**Closes: AC1, AC2, AC3.** Snapshot → `.snap/M2/` (all three deliverable paths).

---

### M3 — the full acceptance sweep, recorded in one pass

Re-run AC1–AC5 and RD-1 from above, then AC7 in full:

```bash
bash -n design_docs/verification/w-race-gate-blindspot/run.sh; echo "rc=$?"      # 0   (base: 0)
stat -f '%Sp' design_docs/verification/w-race-gate-blindspot/run.sh              # -rwxr-xr-x
gofmt -l host/verifygate/ | wc -l                                                # 0   (base: 0)
go vet ./host/verifygate/; echo "vet rc=$?"                                      # 0   (base: 0)
go build ./...; echo "build rc=$?"                                               # 0   (base: 0, 1.7 s)
# D1: the export is MANDATORY here.  ~50 s.
AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -count=1 -v > /tmp/p44-pkg.txt 2>&1
echo "pkg rc=$?"                                                                 # 0
grep -c '^--- FAIL' /tmp/p44-pkg.txt        # 0
grep -c 'AILANG_BIN is unset' /tmp/p44-pkg.txt   # 0   ← if >0 you forgot the export (D1), not a regression
grep -c 'rig-absolute path' /tmp/p44-pkg.txt     # 0   ← TestNoRigAbsolutePaths on the edited file
grep -c '^--- PASS' /tmp/p44-pkg.txt        # >= 32  (base: 31)
```

**AC6 is CONTROLLER-OWNED (D2).** Report it `NOT ATTEMPTED (controller-owned)`. Do not read run
`33069999131` and call it AC6 — that run is the **base** contradiction.

**This sweep is the green control for the whole drill. Capture it before touching anything in M4.**

---

### M4 — the mutation drill (§5) and hand-off

Run §5 in full, then leave the tree **uncommitted** with `.snap/M1/`, `.snap/M2/` present and
report: the three touched paths, the M3 sweep table with every observed value beside its
expectation, the nine drill arms and the X1-C scoping control, each with its landed-proof and observed RED lines, the post-drill
green control, and the two sha256s equal to `POST_SPRINT.sha`.

---

## 5. The mutation drill — 9 RED arms (X1–X9), 1 scoping GREEN control (X1-C), 1 post-drill green control

Arms are labelled **X1–X8** here to avoid colliding with milestone names; each maps 1:1 to the
design doc's §Named RED mutations rows **M1–M8**. Venue: `verify_go.sh:258` runs
`go test ./... -count=1` from the repo root inside the gated `go-verify` job, so a `verifygate`
RED **and** a step-exit RED are both CI-gate reds (P15, re-derived).

Every arm shows three things **in this order**:

1. **The mutation LANDED** — a grep/count that **moved**, printed *before* any result is read.
   *A mutation that did not apply is a false green* (iteration-131's lesson).
2. **The tree still builds** — `go build ./...` rc=0 **and** `bash -n …/run.sh` rc=0 where the arm
   touches the script. A mutation that breaks the build reds for the wrong reason and proves nothing.
3. **The expected RED**, judged by **counted `--- FAIL` lines and the named message**, or by a
   directly-captured `rc` — never through a pipe.

### 5.1 The house restore recipe — read before arm X1

**Take the backup AFTER M2, from the POST-SPRINT tree.** Backing up from `aa0ab05` restores the
sprint away.

```bash
mkdir -p /tmp/p44-backup/gh
cp design_docs/verification/w-race-gate-blindspot/run.sh /tmp/p44-backup/run.sh
cp .github/workflows/ci.yml                              /tmp/p44-backup/gh/ci.yml
cp design_docs/verification/w-race-gate-blindspot/repro/go.mod /tmp/p44-backup/repro-go.mod
shasum -a 256 /tmp/p44-backup/run.sh /tmp/p44-backup/gh/ci.yml /tmp/p44-backup/repro-go.mod \
  > /tmp/p44-backup/POST_SPRINT.sha
cat /tmp/p44-backup/POST_SPRINT.sha     # record these three hashes in the report
```

**Restore after EVERY arm — `cp`, never `git checkout --`:**

```bash
cp /tmp/p44-backup/run.sh       design_docs/verification/w-race-gate-blindspot/run.sh
cp /tmp/p44-backup/gh/ci.yml    .github/workflows/ci.yml
cp /tmp/p44-backup/repro-go.mod design_docs/verification/w-race-gate-blindspot/repro/go.mod
chmod +x design_docs/verification/w-race-gate-blindspot/run.sh     # X5 clears the bit; cp does not restore it
shasum -a 256 design_docs/verification/w-race-gate-blindspot/run.sh \
              .github/workflows/ci.yml \
              design_docs/verification/w-race-gate-blindspot/repro/go.mod    # compare to POST_SPRINT.sha
git status --porcelain    # non-empty by construction; assert only the EXPECTED paths appear
```

`git checkout -- design_docs/verification/w-race-gate-blindspot/run.sh` would **delete M1
entirely**, because nothing is committed. There is no undo. Do not type it.

### 5.2 The arms

Shorthand: `T() { go test ./host/verifygate/ -run '^TestMiscompileInstrumentStepIsGatedInCI$' -count=1 -v 2>&1; }`
and `S() { ./design_docs/verification/w-race-gate-blindspot/run.sh > /tmp/p44-arm.out 2>&1; echo $?; }`

| # | doc row | Exact edit | Landed-proof (printed BEFORE the result) | Expected RED |
|---|---|---|---|---|
| **X1** | M1 | insert `        continue-on-error: true` under the renamed step in `ci.yml` | `grep -c 'continue-on-error' .github/workflows/ci.yml` **0 → 1**; block-scoped count **0 → 1** | `T \| grep -c '^--- FAIL'` = **1**, message `re-introduces "continue-on-error" in the miscompile step` with the offending line number. **MEASURED at V19, not predicted** |
| **X1-C** | *(scoping control — NOT a red arm)* | plant `continue-on-error: true` on the UNRELATED `worldd benchmark smoke gate` step | `grep -c 'continue-on-error' .github/workflows/ci.yml` **0 → 1**; **block-scoped count stays 0** | **GREEN**: `T \| grep -c '^--- PASS'` = **1**. This is the load-bearing control — without it X1's red is equally consistent with a repo-wide ban that round-2 R1 forbids. MEASURED at V19 |
| **X2** | M2 | `\tdarwin/arm64) expect_defect=1 ;;` → `\tlinux/amd64) expect_defect=1 ;;` | `grep -c 'linux/amd64) expect_defect=1' run.sh` **0 → 1**; `grep -c 'darwin/arm64) expect_defect=1'` **1 → 0** | attended darwin: `S` → **rc=1**, `no verified platform contract for darwin/arm64` = 1 (the fail-closed arm). Linux-side prediction is profile-replayed: re-run RD-1's replay with `expect_defect=1 saw_bad=0 host_pair=linux/amd64` → **rc=1** `(or GOOD NEWS)`. **A polarity lie now fails on BOTH platforms** |
| **X3** | M3 | `mv repro/main.go repro/main.go.MUT` | `(cd repro && GOTOOLCHAIN=go1.26.6 go build .)` → rc≠0, `no Go files in …/repro` (quote it) | `S` → **rc=1**, `no toolchain ran at all` = 1, `KNOWN-BAD probes completed` = **0**. **Planner-rehearsed on the candidate: exactly this** |
| **X4** | M4 | `repro/go.mod`: `go 1.22` → `go 1.27` | `grep '^go ' repro/go.mod` prints `go 1.27` | **TWO layers**: `go test ./host/verifygate/ -run '^TestReproModuleFloorStaysBelowKnownBadToolchains$' -count=1 -v` → `--- FAIL`=1 naming the floor; and `S` → **rc=1** via floor 1 (`ran==0`), every probe SKIPping. Restore from `/tmp/p44-backup/repro-go.mod` |
| **X5** | M5 | `chmod -x run.sh` | `stat -f '%Sp' run.sh` → `-rw-r--r--` | `S` → **rc=126** (planner-measured), **and** `T2() { go test ./host/verifygate/ -run '^TestMiscompileInstrumentProbesPinnedToolchain$' -count=1 -v; }` → `--- FAIL`=1, `is not executable`. **Restore with `chmod +x`** |
| **X6** | M6 | append ` go1.99.99` to `KNOWN_BAD=` | `grep -c 'go1.99.99' run.sh` **0 → 1** | `S` → **rc=1**, `INSTRUMENT FAILURE: 4 of 5 KNOWN-BAD probes completed` (planner-measured verbatim). **The static layer stays GREEN by design** — `go/version.IsValid("go1.99.99")`=true and oldest KNOWN_BAD is still `go1.26.0`, so `TestReproModuleFloorStaysBelowKnownBadToolchains` passes; record that green, it is the point (these teeth are behavioural) |
| **X7** | *(new — AC3's teeth)* | delete the `if count := strings.Count(src, "saw_pinned_ok"); count < 3` guard's **body line** from `TestMiscompileInstrumentProbesPinnedToolchain` | `grep -c 'saw_pinned_ok site count' host/verifygate/toolchain_pin_gate_test.go` **1 → 0** | **AC3's body diff becomes NON-EMPTY.** Restore. Without this arm AC3's empty diff is green-by-identity and measures nothing |
| **X8** | M8 | `case "$(uname -s)" in` → `case "$(echo Darwin)" in` | `grep -c 'uname -s' run.sh` **1 → 0** | `T \| grep -c '^--- FAIL'` = **1**, `run.sh no longer reads the kernel via "uname -s"`. **MEASURED at V19 (M7b).** Note the behavioural layer alone would MISS this on a darwin rig — the mutant still resolves `darwin` — which is why the static pin is load-bearing |
| **X9** | M7 | `host_pair="$host_os/$host_arch"` → `host_pair="$(go env GOOS)/$(go env GOARCH)"` | `grep -c 'host_pair="\$(go env GOOS)' run.sh` **0 → 1** | `T \| grep -c '^--- FAIL'` = **1**, `executable use of "go env GOOS"` (+`GOARCH`). **MEASURED at V19 (M7a)** |

*(Nine rows for eight doc mutations: X1-C is the doc's scoping control, and X7 is the planner's
addition giving AC3 teeth it otherwise lacks. The doc's M7/M8 map to X9/X8.)*

**Green control**: after the final restore, re-run **all of M3**. Every AC green, the package run
clean, and the three sha256s equal to `POST_SPRINT.sha`. **The control is also run BEFORE arm X1**
(M3 provides it) — a drill whose control runs only at the end cannot distinguish "the arms
restored correctly" from "the arms never applied".

### 5.3 Not mutated, and declared

- **The new test file's own assertions.** No arm edits `TestMiscompileInstrumentStepIsGatedInCI`;
  every red is produced from the production side.
- **`scripts/verify_go.sh`, `scripts/verify_ail.sh`, `host/store/toolchain_canary_test.go`,
  `racecontrol/`.** Out of scope — a mutation there tests rows 48/49, not this one.
- **`normalizeToolchainPin`** (top of the same file). Row 45's surface. Do not touch it.
- **A step-level `if:` that never evaluates true.** Declared residual 2: the wiring test cannot
  see it, no actionlint runs in this repo (P41's V18), and no behavioural floor backstops it.
  Declared, not drilled.

---

## 6. What this sprint explicitly does NOT do

- **No third CI job, no darwin runner.** Option B is priced in the design and deferred; it would
  couple into `TestGoToolchainPinsAgreeAndMatchJobList`'s job-set cardinality.
- **No repo-wide `GOOS`/`GOARCH` ban.** Round 2 blocked 3/3 on exactly that; the scoped form is the
  ratified fix. Re-introducing the repo-wide ban reds the tree **on arrival** (V19 arm A).
- **No `ci.yml:172` citation sweep.** 25 lines across 8 files; that census is **row 31's**, and this
  sprint's own edit renumbers `ci.yml` again. The only exception is the ~3 comment lines inside the
  doc-comment M2 must touch anyway, which would otherwise become *false residual statements*, not
  merely stale numbers. The historical step-name quotes at `w-race-gate-blindspot.md:208` and
  `bench/BASELINE.md` stand as true history.
- **No `KNOWN_BAD` membership-length bind.** Declared latent narrowing, scoped out, partly covered
  behaviourally by the coverage floor.
- **No `.ail` byte changes, no `./scripts/verify_ail.sh` run, no `./scripts/build_world_package.sh`.**
  Single-gate (Go) sprint.
- **No widening to `go test ./...`.** The narrowest gate that can fail for this diff is
  `go test ./host/verifygate/ -count=1`. Prefer it. `verify_go.sh` runs the wide form in CI; that
  is the CI leg's business, not the executor's acceptance instrument.
- **No git writes, no PR, no AC6 verdict.**

---

## 7. Velocity and risk

**Velocity.** ~115 changed lines across 3 files, 0 created. Measured cost centres on this rig:
attended `run.sh` **3.3 s warm**; whole `host/verifygate` package **48.4 s** (the sprint's long
pole, run ~3 times); `go build ./...` **1.7 s**; scoped verifygate test **~0.2 s**; RD-1's four
replays **< 1 s total**. The acceptance sweep is ~1 min; the ten-arm drill is ~10 min (five arms
run the 3 s script, three run the 48 s package). **The estimate's risk is not in the typing** —
every proposed byte has been generated and exercised. Budget the 0.35 day in RD-1, in the AC4
split, and in reading counted lines rather than exit codes. Reference band: row 43 (`ecfb62d`)
landed a comparable three-file item with a 5-milestone plan inside 0.1 d; this one is larger by the
`run.sh` rewrite and the behavioural drill.

**Risks.**

| # | Risk | Severity | Mitigation |
|---|---|---|---|
| **R1** | **THE CONSTRAINT THAT OUTRANKS EVERYTHING: the change reds `dev` on the next push, recoverable only by revert.** The linux leg has never reached floors 3–4, so nothing in the repo's history tells you what they do there | **CRITICAL** | **RD-1 is a BLOCKING gate inside M1, re-derived not inherited**: P5's measured linux profile replayed through the *executor's own* floor block → rc=0, with two negative controls proving the block can still fail. M2 (the flip) is not started until RD-1 is green. Planner has already run it once against candidate bytes: **rc=0**. Second half: the `uname` shim proves the ubuntu-24.04 runner maps into the LISTED `linux/amd64` arm, not the refuse arm |
| **R2** | **AC4(b) asserted as the design writes it** → false RED against a correct landing, and the "repair" dismantles R1's fail-closed arm | HIGH | §1 D3, measured. AC4(b) asserts the **refuse** text and `PLATFORM ALARM`=0; the alarm gets its own newly-defined AC4(c) |
| **R3** | **AC7's package run without `AILANG_BIN`** → 17 unrelated FAILs read as a regression | HIGH | §0 rule 4, §1 D1. The failure text `AILANG_BIN is unset` is the tell; the sweep greps for it explicitly |
| **R4** | **AC1 read as rc** — the naive `-run` form is rc=0 at base with zero tests run | HIGH | Every scoped run counts `=== RUN`/`--- PASS`/`--- FAIL` from `-v`; AC1 ships a paired nonsense-pattern control whose rc is also 0 |
| **R5** | **The round-1 repo-wide token ban re-introduced** by an executor "simplifying" the scoped test | MEDIUM | §6; V19 arm A measured the arrival-red; X1-C is the boundary control |
| **R6** | **`git checkout --` reflex during the drill** deletes the uncommitted M1 | MEDIUM | §5.1 states the consequence plainly; `cp` from `/tmp/p44-backup/` only, plus `chmod +x` after X5 |
| **R7** | **A toolchain SKIPs under a network-denied sandbox**, turning a behavioural arm into a silent wrong answer | MEDIUM | §0 rule 9: any unexpected `SKIPPED` ⇒ report **UNINFORMATIVE UNDER SANDBOX**, never a verdict. All 7 toolchains are cached today; only X6's fake `go1.99.99` needs the network, and it SKIPs either way |
| **R8** | **TAB/space corruption** when transcribing the `case` arms and floor bodies | MEDIUM | The file indents with TABs throughout. `bash -n` is the first gate after every M1 edit; a space/TAB mix does not break `bash` but breaks nothing either — the real detector is the RD-1 replay, which fails loudly if a floor is mis-nested |
| **R9** | **A rig-absolute path in the appended test comment** reds `TestNoRigAbsolutePaths` | LOW | Forbidden in M2; M3's sweep greps `rig-absolute path` on the package output |
| **R10** | **Row 45 lands into the same file** | LOW | Disjoint hunks (top vs EOF). Flagged for whichever lands second; no action here |

---

## 8. Open questions for the human

Neither blocks execution.

- **The `ci.yml` comment's forward reference** to `design_docs/implemented/…` (D5) is only true
  after the controller moves the doc on landing. Land as written; the move is load-bearing for
  one comment line.
- **Declared residual 6's price is now live**: after this lands, a *single* transient known-bad
  toolchain-fetch SKIP reds `dev`. The design prices this against 42/42 observed probes and 0
  SKIPs over 6 runs (rule-of-three bound ≈7%/probe at 95%), a ~47 s step under a 15-minute
  ceiling, and an idempotent re-run. **Accepted by the design; surfaced here because the first
  time it fires it will look like a regression and it is not.**

---

**SPRINT_PLAN_PATH**: `design_docs/planned/w-miscompile-instrument-inert-in-ci-sprint-plan.md`
**SPRINT_JSON_PATH**: `.ailang/state/sprint_w-miscompile-instrument-inert-in-ci.json`
**SPRINT_JSON_PATH (TRACKED COPY -- byte-identical)**: `design_docs/planned/sprint_w-miscompile-instrument-inert-in-ci.json`

> **Note for the controller**: `.gitignore:3` is `**/.ailang/`, so the `.ailang/state/` copy is
> **untracked and will not be committed**. The P43 pair was committed as
> `design_docs/planned/sprint_w-floor-raise-coupling-inventory.json` (`ecfb62d`) and moved to
> `implemented/` on landing (`9d81fc6`). Both copies here are byte-identical (sha256 verified);
> commit the `design_docs/planned/` one.
