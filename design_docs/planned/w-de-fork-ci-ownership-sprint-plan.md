# Sprint plan — w-de-fork-ci-ownership (DE-FORK CI ownership repair)

**Design doc**: `design_docs/planned/w-de-fork-ci-ownership.md` (revision 2 + controller-applied
round-3 refinements; AC6/M9 added from the restored `gpt6-astra` reviewer)
**Sprint id**: `w-de-fork-ci-ownership`
**Base**: worktree `/tmp/world-iter167`, branch `mission/world-iter167-defork-ci`, at
`565d0b257dc6f28fcd25479fccbfd334b1e420ea` (= `origin/dev`). `git status --porcelain` shows one
line, `?? design_docs/planned/w-de-fork-ci-ownership.md` (the design doc, not yet committed).
**Direction**: attended ruling `D-WORLD-34` (answered A, Mark Edmondson, 2026-09-07) — World may
retire the local `--driver-fleet-check` diagnostic and its fixtures/invocations; **"Preserve World
application verification"** and **"no blanket deletion of product checks"** bind as hard as the
grant. **Not re-litigable.** This plan proposes no alternative to the retirement and found no
reason to want one.
**Estimate**: ~0.35 d, 5 milestones. Scope is the design's `Files to Modify/Create` list; §8
records the one place I believe that list is incomplete and the two places I believe an
acceptance criterion needs a correction rather than a relaxation.

**Base file pins** (executor: assert these before starting, and restore to them after every
real-tree mutation arm):

| file | sha256 at base |
|---|---|
| `.github/workflows/ci.yml` | `d78fa829b951f458b582e3965967895b1ceb01e24eb9f450a76a482d787eb0a4` |
| `scripts/verify_go.sh` | `2e80cca964a6b56a9f2b4fca06d59d707c236c6fe6f9865ecf85c70ec84cb33d` |
| `host/verifygate/toolchain_pin_gate_test.go` | `9ac3916f068118eecb5b87e7e5f2bdf7ffecd1970e708bdf29068119862b7c37` |

---

## 0. What the planner actually ran (measure BEFORE plan)

Everything in §2 was executed by me in this worktree at `565d0b2`, this iteration, with
`export PATH=/opt/homebrew/bin:$PATH` and `export AILANG_BIN=$HOME/.pinned-ailang/ailang`
(reports `AILANG v0.30.0`). **No production file was mutated for any probe**; every fixture lives
under `/tmp/iter167-baseline/`. The fleet checkout `~/dev/sunholo-data/ailang` was **not read** —
where a probe needed `--driver-fleet-check` I pointed `AILANG_FLEET_REPO` at an absent directory
(`/tmp/definitely-not-a-repo-9f3`) and measured the absent-fleet refusal branch instead. The
shared main checkout was not touched.

Runs: full Go suite ×2 (once under deliberate parallel load, once unloaded), the focused
`host/verifygate` package, `verify_go.sh` ×5 (no-arg, unknown flag, non-dash arg,
`--mission-config-check`, `--evidence-manifest-check`), `--driver-fleet-check` ×1 against an
absent fleet, `verify_ail.sh` ×1, the isolated flake test ×1, a `go env GOVERSION`-shimmed
deny-list probe ×1, the standalone `-race` leg ×1, seven ledger-validator fixture arms, and the
retirement/wiring greps with positive and negative controls.

I carried over **nothing** from the design's Verification Log that was cheap to re-run. V7, V12,
V17, V18, V19, V20, V21, V22 were all **re-measured first-party** and all reproduced. **V11 is
the only row I carried over rather than re-ran** (it reads the forbidden fleet checkout); I
substituted the absent-fleet branch, which is sufficient for every claim this sprint makes.

---

## 1. PLANNER FINDINGS — read these before executing

### Finding 1 (BLOCKING for M7/M8 as written) — the step-6 sequence test is a **banner** parser, so it cannot see a deleted *command* under a retained banner

Solution Design step 6 specifies "parse the default flow … for the ordered banner/guard
sequence" and AC9 says M7 ("delete the plain `go test ./...` line") and M8 ("delete the `-race`
leg") must each red it. Measured, the default flow's banners are:

```
309: echo "── AILANG_BIN=$AILANG_BIN ($ver)"
311: echo "── tracked-binary hygiene gate"
347: echo "── mission routing + decision-ledger gate"        <- retired block
440: echo "   ✓ toolchain floor gate: ..."
442: echo "── race-detector known-positive control"
453: echo "── go build ./..."
457: echo "── focused host/evidence named-manifest gate (37 exact top-level tests)"
472: echo "── go test ./... -count=1"
478: echo "── go test ./... -count=1 -race -timeout 8m"
492: echo "✓ go gate PASSED: build clean, plain and race tests pass with pinned AILANG_BIN ($ver)"
```

The banner text `── go test ./... -count=1` and the **command** `go test ./... -count=1` are two
different lines (472 and 474). A mutation that deletes line 474 and keeps 473-472 removes a
product check while leaving the banner sequence **exactly intact** — the step-6 test stays green.
That is precisely the "delete until green" outcome AC9 exists to make unreachable, and as written
M7/M8 can be satisfied by deleting only the cosmetic anchor.

**Correction (does not relax AC9; strengthens the instrument):** the step-6 sequence test asserts
an ordered list of **(banner, executable anchor)** pairs, not banners alone. For the two legs
that matter it additionally requires:

- exactly one line whose trimmed text is exactly `go test ./... -count=1` **at column 0**
  (the executed plain-suite leg, distinct from the `echo "── …"` banner and from
  `go test -json ./host/evidence -count=1`);
- exactly one occurrence of the literal `cmd = ["go", "test", "./...", "-count=1", "-race", "-timeout", "8m"]`
  inside the python watchdog heredoc, plus exactly one `sys.exit(124)` and one
  `-race leg timed out after 600s`;
- exactly one line whose trimmed text is exactly `go build ./...`;
- exactly one `check_evidence_manifest "$evidence_json" 37`.

M7 and M8 are then re-specified in §4 to delete the **command**, not the banner, and each has the
sequence test as its sole killer. The banner-only variants are kept as extra arms (M7b/M8b) so
both halves are proven live.

### Finding 2 (M1 has TWO killers, and that is fixable by where the banner is printed)

As the design writes it, `check_mission_config` would print both the `── ` banner and the row
count. Then M1 (delete the default call) and M3 (hollow the function body) are killed by the same
set, and neither mutation discriminates wiring from behaviour.

**Prescription** (implementation freedom the design explicitly grants): print
`echo "── World mission-input gate (decision ledger)"` **at the call site**, in the default flow,
and print the row count + non-certification disclosure **from inside `check_mission_config`**.
Then:

- M1 (delete the call) reds the step-6 sequence test (banner anchor absent) **and** AC3's
  default-flow integration arm. Two killers by construction — one wiring, one behavioural — and
  the executor records BOTH.
- M3 (hollow the body) reds AC1 (no nonzero row count, no disclosure) and AC2's bad-input arms,
  and **does not** red the step-6 sequence test. That asymmetry is the proof the two tests measure
  different things.

M1 is flagged in §4 as lacking a sole killer. This is by design, not an accident, and the plan
says which killer belongs to which claim so a green cannot be attributed to the wrong one.

### Finding 3 (BLOCKING for M9 as written) — running the M9 arm against the real script turns the mutation kill into a *timeout*, which the design forbids counting as anything

AC6 requires `AILANG_BIN=<pinned> ./scripts/verify_go.sh --driver-fleet-check` to be nonzero with
**all prerequisites satisfied**. Under M9 (refusal removed) that invocation falls through and runs
the **entire product gate** — measured below at several minutes. Inside the mandatory
`context.WithTimeout(…, 10*time.Second)` helper that is a deadline expiry, and the design's own
rule says a deadline expiry is an instrument failure that "never satisfies a refusal". It equally
must not be recorded as M9's kill: "the test timed out" is not "the test detected the mutation".

**Correction (does not relax AC6):** split AC6's refusal proof into two arms with an explicit
*early*-stop discriminator.

- **AC6-real (unmutated only, cheap, real tree):** `AILANG_BIN=<pinned> ./scripts/verify_go.sh --driver-fleet-check`
  must be **nonzero**, must **not** contain `✓ go gate PASSED`, must contain a diagnostic naming
  `--driver-fleet-check`, and — the load-bearing addition — must **not** contain the string
  `── AILANG_BIN=`. That last assertion proves the refusal fired *before* the binary preflight,
  which is the whole claim, and it makes the run terminate in well under a second. Same for one
  arbitrary unknown flag (`--zzq-unknown-mode-9f3`).
- **AC6-M9 (fixture tree):** the same invocation against a copied `scripts/verify_go.sh` in a
  synthetic root with a valid `AILANG_BIN` shim and a controlled stop immediately after the
  preflight banner. Unmutated: no `── AILANG_BIN=` in the output. With M9 applied to the fixture
  copy: `── AILANG_BIN=` **is** present and the arm reds — cleanly, inside 10 s, with no timeout.

### Finding 4 (cost control AC6 needs and does not state) — the "no-arg default flow still reaches `✓ go gate PASSED`" positive control must not be a second full-gate run

AC6 lists as a positive control "the no-arg default flow still reaches `✓ go gate PASSED`". Run
inside a Go test that is impossible: the 10-second rule forbids it and the wall cost is minutes.
**Discharge it from AC7(2)'s recorded log** — the same real full-gate run, same head, recorded
with rc and wall time, is the control. The plan states this so the evaluator sees a deliberate
discharge rather than a dropped control. The other two AC6 positive controls
(`--mission-config-check` and `--evidence-manifest-check` still dispatch) are cheap and stay
inside the harness.

### Finding 5 (DECLARED RESIDUAL, measured) — the required refusal is scoped to first arguments beginning `-`, and a non-dash argument still fails open

Design step 4 and AC6 both scope the new refusal to "a first argument beginning `-` that matches
no recognized mode". Measured at base:

```
./scripts/verify_go.sh driver-fleet-check   -> rc=127, output BYTE-IDENTICAL to the no-arg run
./scripts/verify_go.sh --zzq-unknown-mode-9f3 -> rc=127, output BYTE-IDENTICAL to the no-arg run
```

Both forms are silently discarded today. The design's refusal closes the second and leaves the
first: after the repair, `./scripts/verify_go.sh driver-fleet-check` will run the whole product
gate and **exit 0**, which is the same fail-open class `gpt6-astra` objected to, one character
away. No World consumer passes a non-dash first argument (measured: `ci.yml:166` passes none;
the only argument-passing consumer is `evidence_manifest_gate_test.go`, which passes
`--evidence-manifest-check` plus two positional operands **after** the flag), so nothing breaks
if the refusal is widened.

**Proposed correction (widening, not relaxing):** refuse **any** unrecognized first argument, not
only dash-prefixed ones. AC6's arms gain one more: `AILANG_BIN=<pinned> ./scripts/verify_go.sh driver-fleet-check`
is nonzero and does not contain `── AILANG_BIN=`. If the executor declines the widening it must
say so and record this residual verbatim; it must not silently ship the narrow form as if it
closed the class.

### Finding 6 (relieves a real fear) — the ~195-line deletion does **not** red `p1NeedleSet`'s `verify_go.sh: FATAL: == 3` assertion

`scripts/verify_go.sh` carries **17** `verify_go.sh: FATAL:` literals at base; the retirement
removes 9 of them (`:178, :221, :225, :231, :237, :255, :363, :368`, plus the `check_driver_fleet`
region's), leaving 8 file-wide. `toolchain_pin_gate_test.go:1370` asserts a count of exactly 3 —
but `p1FloorGateRegion` (`:1305-1328`) delimits its `gate` string by the sentinels
`# --- P1 (queue row 48` and `echo "   ✓ toolchain floor gate:`, both of which survive untouched.
The count is region-scoped, not file-scoped. **Measured green at base and structurally unaffected.**
Recorded because a planner or executor could reasonably (and wrongly) treat this as a blocker and
start weakening it.

### Finding 7 (relieves a second fear) — removing `launchd-drivers` does not red the Z3 job assertion

`TestZ3PinDeclaredOnceAndInstalledInBothJobs` (`ail_binary_gate_test.go:669-712`) counts
**install lines** (`installs != 2`), not jobs, and its known-positive control list names only
`ailang-verify:`, `go-verify:` and `./scripts/verify_go.sh`. The `launchd-drivers` job installs no
Z3. Unaffected. Likewise `runbook_stageb_test.go:363` counts `verify_go.sh` occurrences in
`ci.yml`; the surviving `go-verify` job keeps that count nonzero (measured: 4 mentions, 1 of them
the invocation at `:166`).

### Finding 8 (confirms V21 and upgrades it from "flake" to "load-triggered") — I reproduced the announce timeout on the first attempt, by loading the machine on purpose

Run 1 was executed **concurrently** with the focused `host/verifygate` run. It failed exactly as
V21 records. Run 2, unloaded, was clean. The isolated rerun passed in 0.80 s. This matters for
§6(a): the flake is not random, it is a resource-contention artifact, so the executor's *reporting
protocol* is the control, and the executor must **not** run the AC7 legs concurrently with each
other or with anything else.

### Finding 9 (no action; recorded) — every consumer claim in the design reproduced

Repo-wide, `mission-control.sh` → 1 code site (`verify_go.sh:349`); `test_mission_routing.sh` → 2
(`verify_go.sh:348`, `ci.yml:227`); `test_mission_stall.sh` → 1 (`ci.yml:226`);
`derive-planner-lane.sh` → 1 (`verify_go.sh:349`); `mission-template.plist` → 0. The retirement
grep pattern returns **26** hits in exactly two files (`verify_go.sh` 20, the deleted fixture 6)
and zero elsewhere including `cmd/`. Negative control (nonsense literal) → 0. I found **no
additional consumer** and **no reason to dispute** the design's scope.

---

## 2. BASELINE for every acceptance criterion — measured by me, at `565d0b2`, before any change

Legend: **RED@base** = the AC fails at base and this repair turns it green. **GREEN@base** = the
AC already holds and the criterion is a *regression guard*, not a repair — for each of those the
row says what it is actually protecting.

### AC1 — new mission-config mode on real World input

| # | command | rc | exact marker at base | verdict |
|---|---|---|---|---|
| 1.1 | `AILANG_BIN=<pin> ./scripts/verify_go.sh --mission-config-check` | **127** | `── mission routing + decision-ledger gate` then `/bin/bash: tools/launchd/test_mission_routing.sh: No such file or directory` | **RED@base** — the mode does not exist; the flag is *discarded* and the default flow runs and dies at the deleted routing test |
| 1.2 | same, with `AILANG_BIN` **unset** | **1** | `✗ AILANG_BIN is unset — host/replay tests would t.Skip() silently and this gate would be false-green.` | **RED@base** — AC1's independence claim ("still succeeds with AILANG_BIN unset") is false at base, and this exact line is what the new dispatcher arm must jump ahead of |
| 1.3 | `/bin/bash scripts/mission_decisions.sh --check --file design_docs/world-mission.md` (the validator the new mode must call) | **0** | `decision ledger valid: 22 rows` | **GREEN@base** — the validator works; AC1 is protecting *that the gate calls it*, not that it exists. The test asserts **nonzero**, never `22` |

### AC2 — fixture refusal family (all seven arms, measured against the real validator)

Fixtures built from the charter's real `<!-- decision-ledger:start/end -->` block (22 rows) under
`/tmp/iter167-baseline/fx/`.

| arm | command (`/bin/bash scripts/mission_decisions.sh --check --file …`) | rc | exact diagnostic | verdict |
|---|---|---|---|---|
| valid (mandatory control) | `fx/valid.md` | **0** | `decision ledger valid: 22 rows` | GREEN@base |
| missing charter | `fx/does-not-exist.md` | **1** | `mission_decisions: unreadable doc: …/does-not-exist.md` | GREEN@base |
| absent ledger block | `fx/noblock.md` | **1** | `expected exactly one decision-ledger block (start=0 end=0)` | GREEN@base |
| zero-row ledger | `fx/zerorow.md` | **1** | `decision ledger has no rows` | GREEN@base |
| duplicate ID | `fx/dupid.md` | **1** | `duplicate decision ID: D-WORLD-5` | GREEN@base |
| invalid status | `fx/badstatus.md` | **1** | `invalid status for D-WORLD-5: MAYBE` | GREEN@base |
| empty required field | `fx/emptyfield.md` | **1** | `empty decision/evidence field for D-WORLD-5` | GREEN@base |
| missing **validator** | `/bin/bash <absent path>` | **127** | `/bin/bash: …: No such file or directory` | GREEN@base (this is the shell's, not the validator's; the new function must **propagate** it, not skip) |

**All eight are GREEN@base at the validator level, and AC2 is RED@base as a whole** — because at
base no harness runs any of them through the gate, and `check_mission_config` does not exist. What
AC2 is actually protecting is therefore: (i) that the new mode *propagates* each of these eight
distinct failures instead of swallowing them, and (ii) — the part with real teeth — that the new
code does **not** mock the validator into success. The eight diagnostics above are the exact
strings the new fixture arms must assert on; a fixture that reproduces the rc but not the
diagnostic is not evidence.

### AC3 — default-flow integration

| # | command | rc | exact marker at base | verdict |
|---|---|---|---|---|
| 3.1 | `AILANG_BIN=<pin> ./scripts/verify_go.sh` (valid ledger, real tree) | **127** | last three lines: `✓ 0 binary blobs among 327 tracked files` / `── mission routing + decision-ledger gate` / `/bin/bash: tools/launchd/test_mission_routing.sh: No such file or directory` | **RED@base** — there is no reached-stage marker after the mission gate; the flow stops at `:348`, before `go build` |
| 3.2 | invalid ledger aborts before Go | **n/a** | unmeasurable at base: the flow never reaches a ledger check at all, so "aborts before Go" is trivially and vacuously true. Recorded as **UNMEASURABLE@base** | see §8 |
| 3.3 | wall time of 3.1 | — | **1 s** | the base gate spends 1 s and asserts nothing about Go |

### AC4 — CI and default-flow wiring

| # | command | result at base | verdict |
|---|---|---|---|
| 4.1 | `grep -n '^  [a-z-]*:' .github/workflows/ci.yml` (under `jobs:`) | `18: ailang-verify:`, `99: go-verify:`, **`202: launchd-drivers:`** — **3 jobs** | **RED@base** (AC4 wants exactly 2) |
| 4.2 | `grep -n 'tools/launchd' .github/workflows/ci.yml` | **3 active commands**: `:226 /bin/bash tools/launchd/test_mission_stall.sh`, `:227 /bin/bash tools/launchd/test_mission_routing.sh`, `:231 for f in tools/launchd/*.sh; do /bin/bash -n "$f" \|\| exit 1; done` | **RED@base** |
| 4.3 | `grep -n 'tools/launchd' scripts/verify_go.sh` | **10 hits**, of which the active default-flow ones are `:348`, `:349`, `:361`, `:366` (plus `:122`, `:140`, `:175`, `:203`, `:218`, `:242`, `:244` in the retired fleet region) | **RED@base** |
| 4.4 | `grep -n 'verify_go.sh\|verify_ail.sh' .github/workflows/ci.yml` | `:97 run: ./scripts/verify_ail.sh`, `:166 run: ./scripts/verify_go.sh` (both present) | **GREEN@base** — protects against the repair deleting a product job while deleting the driver job |
| 4.5 | is there an existing test asserting 4.1-4.3? | **No.** `grep -rn 'tools/launchd' host cmd --include='*_test.go'` hits only the deleted fixture and one comment at `toolchain_pin_gate_test.go:131` | **RED@base** — AC4's assertion is entirely new; nothing in the suite can see it today |

### AC5 — Go pin assertions stay non-vacuous

| # | command | rc | marker | verdict |
|---|---|---|---|---|
| 5.1 | `AILANG_BIN=<pin> go test ./host/verifygate -run '^TestGoToolchainPinsAgreeAndMatchJobList$' -count=1 -v` | **0**, 5 s wall | `=== RUN TestGoToolchainPinsAgreeAndMatchJobList` present, `--- PASS` | **GREEN@base** |
| 5.2 | control list at `toolchain_pin_gate_test.go:123` | contains the literal `"launchd-drivers:"` | **must be edited or M-removal reds as `instrument failure`, not as an inventory mismatch** (V5, re-read) | GREEN@base |
| 5.3 | `wantJobs` `:163` / `wantGoPinnedJobs` `:164` | `{ailang-verify, go-verify, launchd-drivers}` / `{ailang-verify, go-verify}` | after the repair the two sets **coincide**; the design requires both be kept as independent sets anyway | GREEN@base |

AC5 is GREEN@base and stays green. It is protecting the *edit*: the danger is not that the pin
test fails, it is that an executor makes it pass by weakening it (dropping `wantGoPinnedJobs`,
collapsing per-job attribution into a whole-file count, or deleting the `:123` control entry
rather than updating it). M5/M6 are what prove it was not weakened — see §4 and §6(b).

### AC6 — retirement is complete and the retired mode is REFUSED

| # | command | rc | exact marker at base | verdict |
|---|---|---|---|---|
| 6.1 | `grep -rn -E 'check_driver_fleet\|driver-fleet-check\|AILANG_FLEET_REPO\|REQUIRED_FLEET_PATHS\|driver_fleet_scope' host scripts .github \| wc -l` | — | **26** (`scripts/verify_go.sh` 20, `host/verifygate/driver_fleet_scope_gate_test.go` 6) | **RED@base** (AC6 wants 0) |
| 6.2 | same instrument, needle `--mission-config-check` (positive control) | — | **0** | **RED@base** — the control that makes 6.1's future zero believable does not yet exist |
| 6.3 | `wc -l host/verifygate/driver_fleet_scope_gate_test.go` | — | **293** lines, present | **RED@base** (must be deleted) |
| 6.4 | `AILANG_BIN=<pin> ./scripts/verify_go.sh --zzq-unknown-mode-9f3` vs the no-arg run | **127 / 127** | `diff -q` **silent — BYTE-IDENTICAL**; the argument is DISCARDED and the default flow runs | **RED@base** — this is `gpt6-astra`'s fail-open, reproduced first-party |
| 6.5 | `AILANG_BIN=<pin> ./scripts/verify_go.sh driver-fleet-check` (non-dash) vs no-arg | **127 / 127** | `diff -q` **silent — BYTE-IDENTICAL** | **RED@base**, and **not closed by the design as written** — Finding 5 |
| 6.6 | `AILANG_FLEET_REPO=/tmp/definitely-not-a-repo-9f3 ./scripts/verify_go.sh --driver-fleet-check` | **1** | `verify_go.sh: FATAL: DRIVER DRIFT (D-WORLD-DRIVER-1) — fleet source /tmp/definitely-not-a-repo-9f3 is absent; the fleet-comparison arm cannot run, so driver currency is NOT certified` | **GREEN@base as a "mode exists" control** — proves `--driver-fleet-check` IS dispatched today, so its post-repair refusal is a real change of behaviour and not a no-op |
| 6.7 | `./scripts/verify_go.sh --evidence-manifest-check` (positive control: real modes dispatch) | **2** | `usage: ./scripts/verify_go.sh --evidence-manifest-check JSON EXACT` | **GREEN@base**, must stay byte-identical |
| 6.8 | `grep -c '✓ go gate PASSED' <6.4 output>` | — | **0** | GREEN@base **but vacuously** — the gate dies at `:348` long before the banner. **After the repair this arm is only meaningful together with the `── AILANG_BIN=` absence assertion of Finding 3** |
| 6.9 | `git diff --name-only HEAD -- tools/launchd \| wc -l` | — | **0** | GREEN@base; control is `git diff --name-only` nonempty at landing |

### AC7 — complete product surface

| # | command | rc | wall | marker | verdict |
|---|---|---|---|---|---|
| 7.1a | `AILANG_BIN=<pin> go test ./... -count=1` — **run 1, deliberately under parallel load** | **1** | **135 s** | 18 pkgs `ok`; `--- FAIL: TestCLIRealSubprocessEpisode (6.14s)` / `cli_test.go:196: daemon announcement timed out; stderr=` / `FAIL github.com/sunholo-data/ailang-world/cmd/ailang-worldd 10.622s` | **the V21 flake, reproduced first-party on the first attempt** |
| 7.1b | isolated rerun `go test ./cmd/ailang-worldd -run '^TestCLIRealSubprocessEpisode$' -count=1 -v` | **0** | **2 s** | `--- PASS: TestCLIRealSubprocessEpisode (0.80s)` | matches V21's 0.86 s |
| 7.1c | `AILANG_BIN=<pin> go test ./... -count=1` — **run 2, unloaded** | **0** | **84 s** | **19 packages `ok`**, zero FAIL lines | **GREEN@base** — the product suite is healthy; slowest: `host/broker` 132.9 s and `host/verifygate` 126.4 s **under load**, both well under their loaded selves when run alone |
| 7.2 | `AILANG_BIN=<pin> ./scripts/verify_go.sh` | **127** | **1 s** | dies at `── mission routing + decision-ledger gate` → `/bin/bash: tools/launchd/test_mission_routing.sh: No such file or directory`. **No `── go build ./...`, no `go test` output, no `✓ go gate PASSED`** | **RED@base — this is the sprint's headline red.** The Go product gate is *suspended*, not failing |
| 7.3 | `AILANG_BIN=<pin> ./scripts/verify_ail.sh` | **0** | **6 s** | `✓ world package gate PASSED: 9/9 steps performed non-zero work` / `✓ verify gate PASSED: 11 required identities verified, 40 named tests pass` | **GREEN@base** — protects against collateral damage to the AIL lane |
| 7.4 | `gh api repos/sunholo-data/ailang-world/commits/565d0b2…/check-runs` | — | — | `go host build + test gate` **failure**; `launchd drivers (bash 3.2)` **failure**; `ailang-code verify gate` **success** | **RED@base (2 of 3)** |
| 7.5 | `AILANG_BIN=<pin> go test ./host/verifygate -count=1` | **0** | **129 s under load** (`ok … 127.668s`) | zero FAIL | **GREEN@base** |
| 7.6 | standalone `-race` leg, `go test ./... -count=1 -race -timeout 8m`, run under the *same* 600 s python `subprocess`+`killpg` watchdog `verify_go.sh:480-490` uses | **0** | **194 s** | **19 packages `ok`**, zero `WARNING: DATA RACE`, zero FAIL; slowest `host/broker` 170.3 s, `host/verifygate` 89.4 s. **406 s of headroom under the script's own 600 s watchdog** | **GREEN@base when run directly**, but **unreachable through `verify_go.sh` at base** (7.2 dies at `:348`). This is the dominant cost of 7.2 |

### AC8 — documentation

| # | check | at base | verdict |
|---|---|---|---|
| 8.1 | `grep -n "drift gate" CLAUDE.md` | `:29` — "…and `verify_go.sh`'s drift gate reds while the working-tree driver diverges from HEAD — that red means 'the fleet must commit'…" | **RED@base** — the clause describes live enforcement that this sprint retires; it becomes false at landing |
| 8.2 | design doc status line | `Status: Revision 2 … No acceptance criterion is met; nothing has been executed.` | RED@base, updated at landing |
| 8.3 | independent evaluator | not yet spawned | RED@base |

### AC9 — every surviving product check is proven to still run

Baseline is 7.2's log. **At base the log contains 2 of the 10 required anchors and stops.**

| required anchor (AC9's order) | present at base? |
|---|---|
| `── AILANG_BIN=` with `(AILANG v0.30.0` | **YES** (`── AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang (AILANG v0.30.0`) |
| `── tracked-binary hygiene gate` + `✓ 0 binary blobs among N tracked files`, N>0 | **YES** (`✓ 0 binary blobs among 327 tracked files`) |
| new mission-config banner + nonzero row count | **NO** (does not exist) |
| `✓ toolchain floor gate:` | **NO** (source `:440`, never reached) |
| `── race-detector known-positive control` + `WARNING: DATA RACE` | **NO** (source `:442`) |
| `── go build ./...` | **NO** (source `:453`) |
| `── focused host/evidence named-manifest gate (37 exact top-level tests)` + `✓ all 37 required top-level evidence tests passed exactly once` | **NO** (source `:457`) |
| `── go test ./... -count=1` | **NO** (source `:472`) |
| `── go test ./... -count=1 -race -timeout 8m` | **NO** (source `:478`) |
| `✓ go gate PASSED` | **NO** (source `:492`) |

AC9's deny-list sub-arm, measured: with a `go` shim on PATH whose `go env GOVERSION` returns
`go1.26.1` (everything else `exec`s the real `go`), `./scripts/verify_go.sh` is **rc=127** and
`grep -c miscompiles` is **0** — the deny-list at `:390` is **unreachable at base**. So this
sub-arm is **RED@base for the reason AC9 exists**, and its post-repair green (`verify_go.sh: FATAL:
active toolchain go1.26.1 miscompiles host/store/scan.go's`) is a genuine restoration, not a
recorded no-op. Shim recipe, verbatim, is in §7.

---

## 3. Milestones

Each is independently committable and ends at a stated green gate. **MS2 is the one that makes CI
green. MS3, MS4 and MS5 are the controls that stop "delete until green" from being the outcome** —
and the plan deliberately orders MS3 *before* MS4 so the new World check exists in the tree before
the ~195-line deletion lands.

### MS1 (~0.03 d) — retire the obsolete CI job and update the inventory pin
**Files**: `.github/workflows/ci.yml` (delete lines **201-231**, the blank line + the whole
`launchd-drivers:` job through EOF — 31 lines), `host/verifygate/toolchain_pin_gate_test.go`
(remove `"launchd-drivers:"` from the `:123` control list; `wantJobs` `:163` →
`{"ailang-verify","go-verify"}`; update the stale `:131-138` comment about the third job —
**keep `wantGoPinnedJobs` as a separate slice** even though the two now coincide, and keep the
`:164-168` cross-check that every Go-pinned job is in `wantJobs`).
**Does NOT**: touch `scripts/verify_go.sh`, touch anything under `tools/launchd/`.
**Green gate**: `AILANG_BIN=<pin> go test ./host/verifygate -run '^TestGoToolchainPinsAgreeAndMatchJobList$|^TestZ3PinDeclaredOnceAndInstalledInBothJobs$' -count=1 -v`
rc=0 **with both `=== RUN` lines present in `-v` output** (a bare `-run` is rc=0 on a typo — see
§5), plus `AILANG_BIN=<pin> go test ./host/runbook -count=1` rc=0.
**Gates**: AC5 (first half), part of AC4.
**Note**: this milestone alone turns the `launchd drivers (bash 3.2)` check red → **absent**. It
does **not** fix `go-verify`.

### MS2 (~0.05 d) — add `check_mission_config` and wire it in place of the retired local-copy block — **THIS IS THE MILESTONE THAT MAKES CI GREEN**
**File**: `scripts/verify_go.sh` only.
1. Define `check_mission_config()` beside `check_evidence_manifest` (before the dispatcher). It
   runs `/bin/bash scripts/mission_decisions.sh --check --file design_docs/world-mission.md` with
   the **exact** charter path; it `bash -n`-checks the validator first; it fails (never skips) if
   either the validator or the charter is missing or unreadable; it propagates the validator's
   rc; it prints the validator's row-count line and then the explicit scope disclosure —
   *World decision-ledger validated; centralized runtime currency is NOT certified here.*
2. Add the dispatcher arm `--mission-config-check` **after** `--evidence-manifest-check` and
   **before** the `AILANG_BIN` preflight; reject extra arguments (`[ "$#" -ne 1 ]` → usage, exit 2).
3. Replace lines **347-384** (`echo "── mission routing…"` through the `check_driver_fleet` default
   invocation and its `fleet_rc` handling) with, in this order: `echo "── World mission-input gate
   (decision ledger)"` **at the call site** (Finding 2), then `check_mission_config`. Everything
   from `:386` (the deny-list comment) onward is untouched, byte for byte.
**Green gate**: `AILANG_BIN=<pin> ./scripts/verify_go.sh` reaches `✓ go gate PASSED` (rc=0);
`/bin/bash scripts/verify_go.sh --mission-config-check` rc=0 **with `AILANG_BIN` unset**;
`bash -n scripts/verify_go.sh` rc=0.
**Gates**: AC1, AC3 (behavioural half), AC7(2), AC9 (behavioural half).
**After MS1+MS2 both CI checks are green.** Stopping here would be exactly the "delete until
green" outcome the ruling forbids; MS3-MS5 are what make it unreachable.

### MS3 (~0.10 d) — the World-owned regression tests (**control milestone #1**)
**File**: `host/verifygate/mission_config_gate_test.go` (new, ~220 lines).
Arms, all with `context.WithTimeout(context.Background(), 10*time.Second)` + `exec.CommandContext`,
each checking `ctx.Err()` **before** interpreting a nonzero exit, and each bounding output-pipe
cleanup so an inherited pipe cannot keep the test hung after cancellation:
- **AC1 arms**: real mode against the real charter → rc=0, **nonzero** row count (regex, never the
  literal 22), disclosure string present, with `AILANG_BIN` and every `AILANG_*`/fleet variable
  cleared from the child environment.
- **AC2 arms**: the eight fixtures of §2/AC2, each asserting **rc AND the exact diagnostic**, in
  isolated copied-script roots carrying the **real** `scripts/mission_decisions.sh` (never a mock),
  plus the mandatory valid-fixture control in the same harness.
- **Blocking control** (AC2's last sentence): a fixture whose validator is replaced by a script
  that sleeps past the deadline → the helper must classify it explicitly as
  `instrument failure: context deadline exceeded`, and the test must assert that string, so a
  timeout can never be counted as one of the seven refusals.
- **`bash -n` control re-homed from the deleted fixture (V8)**: run `/bin/bash -n` over the
  **copied** `scripts/verify_go.sh` in the fixture root; rc must be 0. Base-measured: rc=0.
- **AC3 default-flow arm**: fixture root, valid pinned-binary shim, controlled stop after a chosen
  later stage; assert the mission banner and row count appear **before** the stop marker for the
  valid ledger, and that an invalid ledger stops **before** any `go` invocation.
- **AC4 wiring arms** (read the REAL `ci.yml` and REAL `verify_go.sh`): exactly two jobs; both
  `./scripts/verify_go.sh` and `./scripts/verify_ail.sh` still invoked; **zero** active
  (non-comment) commands referencing `tools/launchd/`; the parsed active-command set asserted
  **non-empty** first (an empty parse must be an instrument failure, not a pass) and the two
  `verify_*.sh` invocations used as the known-positive anchors. Each failure message **names the
  consumer file and line**.
- **AC6 arms**: the retirement grep with its `--mission-config-check` positive control and a
  nonsense-literal negative control; the deleted fixture's absence; the refusal arms of Finding 3
  (real-tree cheap arm + fixture arm), including the `── AILANG_BIN=` absence discriminator; the
  `--evidence-manifest-check` and `--mission-config-check` dispatch controls.
**Green gate**: `AILANG_BIN=<pin> go test ./host/verifygate -run '^TestMissionConfig' -count=1 -v`
rc=0 **with every expected `=== RUN` subtest name present**, and `go vet ./host/verifygate` rc=0
**read before any test verdict**.
**Gates**: AC1, AC2, AC3, AC4, AC6.

### MS4 (~0.05 d) — retire the fleet-comparison region and delete its fixture (**the deletion, done last of the code changes**)
**Files**: `scripts/verify_go.sh` (delete `:112-123` preamble+config, `:125-259`
`check_driver_fleet`, `:270-278` the `--driver-fleet-check` dispatcher arm) and
`host/verifygate/driver_fleet_scope_gate_test.go` (DELETE, 293 lines). In the **same** edit add
the unknown-first-argument refusal immediately after the two surviving mode arms and **before**
the `AILANG_BIN` preflight — nonzero exit, diagnostic **naming the argument**, and per Finding 5
covering **any** unrecognized first argument, not only dash-prefixed ones. **No stub, no alias, no
`exit 0` "retired" arm.**
**Green gate**: retirement grep = 0 with its positive control firing; MS3's AC6 arms green;
`AILANG_BIN=<pin> ./scripts/verify_go.sh` still reaches `✓ go gate PASSED`;
`AILANG_BIN=<pin> go test ./host/verifygate -count=1` rc=0.
**Gates**: AC6.
**Ordering rationale**: MS3 lands the replacement controls *before* MS4 removes 195 lines, so at no
commit does the tree contain the deletion without the check that replaces it.

### MS5 (~0.04 d) — product-check preservation sequence test + docs (**control milestone #2**)
**Files**: `host/verifygate/mission_config_gate_test.go` (the step-6 test), `CLAUDE.md` (one
clause), `design_docs/planned/w-de-fork-ci-ownership.md` (status + execution evidence).
The sequence test asserts the ordered **(banner, executable anchor)** list of §Finding 1 against
the real `scripts/verify_go.sh`, with the retired block **absent** and the mission-config call
**present**. Known-positive control, per the design: the same parser must first find the V20
(a)-(k) list **on the base text with the retired block still present** — embed the base text as a
testdata fixture (`git show 565d0b2:scripts/verify_go.sh`) so the parser is proven to read a real
file rather than to be trivially satisfied.
`CLAUDE.md:29`: replace the drift-gate-as-live-enforcement clause with a pointer to `D-WORLD-34`.
No other CLAUDE.md edit.
**Green gate**: the full §7 block.
**Gates**: AC8, AC9, AC7.

---

## 4. Mutation plan, made executable

**Protocol for every arm** (non-negotiable, from this loop's scar tissue):
1. Back up the target with `cp` to a path **outside the repo and outside `/tmp`**; record `shasum -a 256`.
2. Apply the mutation.
3. Read `go vet ./host/verifygate` rc **BEFORE** any test verdict — a non-compiling mutant reads
   exactly like a green.
4. Run the named killer command. Record **rc and the exact message text**, not just rc.
5. Run a *second*, unrelated in-scope command to show the mutation is not globally red.
6. Restore from the `cp` backup. **Never `git checkout -- <file>`** — the working tree carries
   uncommitted sprint work, and `git checkout --` restores to HEAD, silently reverting it.
7. Re-assert `shasum -a 256` equals the pre-mutation value. Only then move to the next arm.
8. `-run` filters are **rc=0 on a typo**: every arm's verdict must be read from `-v` output with
   the expected `=== RUN <name>` line present, never from a bare `-run` rc.

Real-tree arms (M4, M5, M6, M7, M8) mutate `.github/workflows/ci.yml` or `scripts/verify_go.sh`
under this protocol. Fixture-tree arms (M1, M2, M3, M9) mutate copies inside `t.TempDir()`.

| # | exact edit that lands the mutation | command that must turn red | expected sole killer (assertion / arm name) | sole? |
|---|---|---|---|---|
| **M1** | In the fixture copy of `verify_go.sh`, delete **only** the default-flow line `check_mission_config` (keep the `── World mission-input gate` banner echo, keep the function definition) | `go test ./host/verifygate -run '^TestMissionConfigDefaultFlowIntegration$' -count=1 -v` and `-run '^TestVerifyGoProductCheckSequence$' -count=1 -v` | **TWO killers by construction**: (a) `MissionConfigDefaultFlowIntegration/valid_ledger` — the row-count line is absent before the stop marker; (b) `VerifyGoProductCheckSequence` — the *executable anchor* `check_mission_config` is missing from the ordered pair list | **NO — flagged, Finding 2.** Deliberate: (a) is the behavioural claim, (b) the wiring claim. Executor records both, and confirms (a) reds when (b) is skipped |
| **M2** | Append `\|\| true` to the validator invocation inside `check_mission_config` | `go test ./host/verifygate -run '^TestMissionConfigRejectsBadLedger$' -count=1 -v` | `MissionConfigRejectsBadLedger` — **all six** invalid-input subtests flip to rc=0 | **NO — killed by six sibling subtests of one test.** That is one killer at test granularity; acceptable. Also reds AC3's invalid-ledger arm (integration), which the executor records as a second, expected killer |
| **M3** | Replace `check_mission_config`'s body with `return 0` (keep the function and the call site) | `go test ./host/verifygate -run '^TestMissionConfigRealCharter$\|^TestMissionConfigRejectsBadLedger$' -count=1 -v` | `MissionConfigRealCharter` — no nonzero row count, no disclosure string | **YES for the sequence test's purposes** — and the discriminating fact the executor must record is that `VerifyGoProductCheckSequence` stays **GREEN** under M3. If it reds, the banner was put inside the function and Finding 2's separation was not implemented |
| **M4a** | Real tree: re-insert `/bin/bash tools/launchd/test_mission_routing.sh` into `scripts/verify_go.sh`'s default flow | `go test ./host/verifygate -run '^TestActiveWiringHasNoDriverConsumers$' -count=1 -v` | `ActiveWiringHasNoDriverConsumers/verify_go_sh` — message must name `scripts/verify_go.sh:<line>` | **NO** — also reds `AILANG_BIN=<pin> ./scripts/verify_go.sh` (rc=127). Expected and recorded; the *test* is the killer, the gate red is corroboration |
| **M4b** | Real tree: add a step `run: /bin/bash tools/launchd/test_mission_stall.sh` to the `go-verify` job in `ci.yml` | same command | `ActiveWiringHasNoDriverConsumers/ci_yml` — message must name `.github/workflows/ci.yml:<line>` | **YES** |
| **M5** | Real tree: delete `GOTOOLCHAIN: go1.26.6` from the `ailang-verify` job and add a second one to `go-verify` (whole-file count stays 2) | `go test ./host/verifygate -run '^TestGoToolchainPinsAgreeAndMatchJobList$' -count=1 -v` | the **per-job attribution** arm: `ci.yml: job "ailang-verify" carries 0 GOTOOLCHAIN line(s), want 1` **and** `job "go-verify" carries 2 …, want 1`. The whole-file-vs-per-job-sum arm stays green (2 == 2), which is the point | **YES** |
| **M6** | Real tree: append a well-formed unclassified job to `ci.yml` (`  probe-job:` / `    name: probe` / `    runs-on: ubuntu-latest` / `    steps:` / `      - run: 'true'`) — no Go pin, no Z3 install | same command | the **exact inventory** arm: `ci.yml: enumerated jobs=[ailang-verify go-verify probe-job], want [ailang-verify go-verify]` | **YES** — and the executor must confirm `TestZ3PinDeclaredOnceAndInstalledInBothJobs` stays green (installs still 2), which is what proves M6 is an *inventory* kill and not a blanket red |
| **M7** | Real tree: delete the **command** line `go test ./... -count=1` at `:474` (keep the banner at `:472` and the `go version` at `:473`) | `go test ./host/verifygate -run '^TestVerifyGoProductCheckSequence$' -count=1 -v` | `VerifyGoProductCheckSequence` — the executable anchor `go test ./... -count=1` at column 0 has count 0, want 1 | **YES** — **only if Finding 1's correction is implemented.** With banners alone this arm is a **silent pass**: measured, the banner and the command are different lines |
| **M7b** | Real tree: delete the **banner** `echo "── go test ./... -count=1"` (keep the command) | same command | same test, banner-anchor arm | YES |
| **M8** | Real tree: delete the whole `-race` leg (`:478-490`: banner, `go version`, and the python heredoc) | same command | `VerifyGoProductCheckSequence` — banner anchor and the `cmd = ["go", "test", …, "-race", …]` / `sys.exit(124)` / `-race leg timed out after 600s` anchors all count 0 | **YES** |
| **M8b** | Real tree: delete only `sys.exit(124)` + the 600 s timeout branch from the heredoc (keeping the leg) | same command | same test, watchdog-anchor arm | YES — proves the watchdog is pinned, not just the leg |
| **M9** | Fixture copy of `verify_go.sh` (Finding 3): delete the unknown-first-argument refusal block, keeping both real mode arms; fixture carries a valid `AILANG_BIN` shim reporting `AILANG v0.30.0` and a controlled stop immediately after the preflight banner | `go test ./host/verifygate -run '^TestRetiredModeIsRefused$' -count=1 -v` | `RetiredModeIsRefused/driver_fleet_check` — **the discriminator is the presence of `── AILANG_BIN=` in the output**, which is absent when the refusal is live and present when M9 removes it. Secondary assertions in the same arm: rc becomes 0-or-controlled-stop rather than the refusal rc, and the diagnostic naming `--driver-fleet-check` disappears | **YES** — provided the arm is fixture-based. **Run against the real script it is killed by a 10-second deadline expiry, which the design forbids counting as a refusal and which the executor must equally not count as a kill** |

**Mutations I believe lack a sole killer**: **M1** (two, by deliberate construction — Finding 2)
and **M2** (six sibling subtests plus an integration arm — one killer at test granularity, several
at subtest granularity). Neither is a missing test; both are *over*-covered rather than under-,
and the plan says which assertion carries which claim so the executor cannot attribute a green to
the wrong one. **No mutation in this set is killed by nothing** once Finding 1's correction lands;
**M7 as the design writes it IS killed by nothing**, which is why Finding 1 is marked BLOCKING.

---

## 5. Explicit finite budgets

Anchored on CI's own ceilings, measured at base: `ci.yml:165` `timeout-minutes: 25` on the
`./scripts/verify_go.sh` step of `go-verify`; `ci.yml:175` `timeout-minutes: 15` on the compiler
reproducer; `ci.yml:182` and `:186` `timeout-minutes: 2` on the two bench steps; `verify_go.sh:480-490`
the internal **600 s** python `-race` watchdog (SIGKILL to the process group, `sys.exit(124)`).
`ci.yml:212`'s `timeout-minutes: 10` belongs to the job MS1 deletes.

| gate | measured at base | **budget (deadline)** | on expiry |
|---|---|---|---|
| `go test ./... -count=1` | 84 s unloaded / 135 s loaded | **8 min** | FAILURE. Record rc, wall and the full log. Not a retry-in-place |
| `go test ./host/verifygate -count=1` | 129 s under load | **8 min** | FAILURE |
| `-race` leg alone | **194 s, rc=0, 19 pkgs ok, 0 data races** (measured) | **600 s** — the script's own watchdog; do **not** raise it. 406 s headroom at base | the watchdog SIGKILLs the process group and `sys.exit(124)`; that IS the failure, never a retry |
| real `./scripts/verify_go.sh` end to end | unreachable at base (dies in 1 s) | **20 min** — deliberately **under** CI's 25-min step ceiling, so a local green implies a CI green with 5 min of headroom | FAILURE |
| `./scripts/verify_ail.sh` | 6 s | **5 min** | FAILURE |
| every new Go fixture/mutation subprocess helper | — | **10 s** (`context.WithTimeout`, `exec.CommandContext`), checked via `ctx.Err()` **before** interpreting a nonzero exit | explicit `instrument failure: context deadline exceeded` — **never** a refusal, **never** a mutation kill |
| whole executor sprint wall | — | **90 min** of gate time | report partial results honestly; do not drop an AC to fit |

**There is no `timeout(1)` on this rig.** Bound long shell gates by wrapping in the same
`python3` `subprocess` + `p.wait(timeout=N)` + `os.killpg` pattern `verify_go.sh:480` already uses,
or by `date +%s` before/after with the deadline asserted after the fact. Do not invent a
`timeout`/`gtimeout` call — it will not exist and the failure will look like the gate's.

Derived cost of `verify_go.sh` end to end, from measured parts: preflight+hygiene ~2 s, the new
mission gate <1 s, race-detector control + `go build ./...` + the focused `host/evidence` 37-test
manifest ~30-60 s, plain full suite **84 s**, `-race` leg **194 s** ⇒ **≈ 5-6 min**, against a
**20 min** budget and CI's **25 min** step ceiling. Total AC7 gate time, sequential and unloaded:
84 + ~330 + 6 + 129 ≈ **9-10 min**, comfortably inside the 90-min envelope. **Run them
sequentially.** Finding 8: running them concurrently is what manufactures the AC7 flake — I
reproduced it that way on the first attempt.

---

## 6. Risk register

| # | risk | disposition |
|---|---|---|
| **(a)** | **`TestCLIRealSubprocessEpisode` load flake** (V21; reproduced by me first-party at 7.1a: rc=1, 135 s, `cli_test.go:196: daemon announcement timed out; stderr=`; isolated rerun PASS in 0.80 s; unloaded full suite rc=0, 84 s, 19 pkgs). It is **load-triggered, not random** (Finding 8). | **Prevention first**: the executor runs the AC7 legs **strictly sequentially**, nothing else on the rig. **If it fires anyway, the protocol is fixed and "rerun until green" is NOT available**: (1) record the red — rc, wall, and the verbatim `cli_test.go:196` line — in the evidence log as `AC7(1) run 1`; (2) run `go test ./cmd/ailang-worldd -run '^TestCLIRealSubprocessEpisode$' -count=1 -v` in isolation and record it as a **diagnostic**, explicitly labelled *not* an AC pass; (3) rerun the **full** command once, as `AC7(1) run 2`, and **report both results in the sprint record**; (4) a second red is a **FAILURE of AC7** — stop, do not run a third. At most **one** disclosed rerun, and the disclosure is mandatory whether or not the rerun is green. Reporting only the green run is a fabricated result. It is a pre-existing flake outside this sprint's scope: a candidate queue row, **not** a change here. |
| **(b)** | **Weakening the job-inventory pin while updating it** (AC5/M5/M6). The edit must remove three literals; the tempting shortcuts are to drop `wantGoPinnedJobs` (now equal to `wantJobs`), to collapse per-job attribution into the whole-file count, or to delete rather than update the `:123` control entry. Any of those turns a live gate into a tautology while all tests stay green. | **Bounded diff + mutation proof.** The `toolchain_pin_gate_test.go` diff must be **≤ ~12 changed lines** and must not delete `wantGoPinnedJobs`, the `:164-168` membership cross-check, the per-job loop, or the whole-file-equals-per-job-sum arm. M5 and M6 are the proof, and §4 names their *distinct* killers precisely so "both mutations red" cannot be reported as evidence when the truth is "everything is red". Evaluator instruction: read the diff of this file line by line before accepting AC5. |
| **(c)** | **Removing the `launchd drivers (bash 3.2)` check name strands a required status check.** | **DISCHARGED, measured.** `gh api repos/sunholo-data/ailang-world/branches/dev/protection` → **404 `Branch not protected`** (re-measured by me this iteration). `dev` has no branch protection and no required-status-check list, so a disappearing check name cannot block a merge. Named here rather than left silent. Separately: the charter's Gate-3b reads the `ailang-code verify gate` check **by name** (`D-WORLD-35`) and that job is untouched. |
| **(d)** | **Over-deleting a product gate while removing the adjacent retired block.** `:347-384` (retired) and `:386` onward (deny-list, floor, race control, build, evidence, suites) are contiguous. An off-by-a-few deletion silently removes a product check. | MS4 is ordered **after** MS3 and MS5's sequence test exists to catch exactly this; plus the executor asserts `shasum -a 256` on the surviving `:386-492` region by extracting it before and after and diffing (`sed -n '386,492p'` at base vs the corresponding post-repair region), and AC9's real log carries all ten anchors. |
| **(e)** | **A success-returning stub for the retired mode.** | Forbidden explicitly (design step 4). AC6's grep (0 hits, positive control firing) plus the Finding-3 refusal arms. M9 proves the refusal is live. |
| **(f)** | **`p1NeedleSet`'s `verify_go.sh: FATAL: == 3` reds on the 195-line deletion.** | **DISCHARGED, measured** (Finding 6): the count is delimited by surviving sentinels, not file-wide. If it reds anyway, that is a **real** signal — the deletion crossed into the P1 floor gate. **Do not relax it.** |
| **(g)** | **The `--driver-fleet-check` retirement removes the only `bash -n scripts/verify_go.sh` control in the suite** (V8). | MS3 re-homes it into the new fixture. Base-measured rc=0. Evaluator must confirm the new fixture actually runs `bash -n` over the **copied** `verify_go.sh`, not over the validator only. |
| **(h)** | **Non-dash first argument still fails open after the repair** (Finding 5). | Proposed widening in MS4. If the executor implements the narrow dash-only form, it must record the residual verbatim and must not claim the fail-open class is closed. |
| **(i)** | **Fleet checkout contamination.** The retired arm reads `~/dev/sunholo-data/ailang`. | The retirement removes the only reader. Until MS4 lands, no command in this plan sets or defaults `AILANG_FLEET_REPO` at a real fleet path; where a probe needed it, point it at `/tmp/definitely-not-a-repo-9f3`. `tools/launchd/*` is frozen: `git diff --name-only <base>..<head> -- tools/launchd` must be **empty** at landing, with `git diff --name-only` nonempty as the control. |
| **(j)** | **Evidence fabrication under time pressure.** Several gates are minutes long and the temptation is to assert a green from an earlier run. | Every AC7/AC9 claim must be backed by a **file** in `/tmp/iter167-evidence/` containing the command, rc, wall time and full output. §7 names them. A claim with no artifact is treated as UNMEASURED, not as a pass. |

---

## 7. Final verification block — the executor runs this **verbatim**, at the end, sequentially

```zsh
set -u
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN=$HOME/.pinned-ailang/ailang
cd /tmp/world-iter167
E=/tmp/iter167-evidence; mkdir -p "$E"

run() {   # run <name> <budget_seconds> -- <cmd...>
  name=$1; budget=$2; shift 3
  s=$(date +%s)
  "$@" > "$E/$name.log" 2>&1
  rc=$?
  e=$(date +%s); w=$((e-s))
  printf 'name=%s rc=%d wall=%ds budget=%ds verdict=%s\n' \
    "$name" "$rc" "$w" "$budget" \
    "$( [ $w -le $budget ] && echo within || echo OVER_BUDGET_FAILURE )" \
    | tee "$E/$name.rc"
}

# --- V0 hygiene ---
run V0-vet          120 -- go vet ./...
run V0-fmt          120 -- gofmt -l host cmd
run V0-bashn         60 -- /bin/bash -n scripts/verify_go.sh

# --- AC7, the complete product surface, IN THIS ORDER, NOTHING ELSE RUNNING ---
run AC7-1-fullsuite 480 -- go test ./... -count=1
run AC7-2-verify_go 1200 -- ./scripts/verify_go.sh
run AC7-3-verify_ail 300 -- ./scripts/verify_ail.sh
run AC7-5-verifygate 480 -- go test ./host/verifygate -count=1

# --- AC9: every surviving product check appears, IN ORDER, in AC7-2's real log ---
python3 - "$E/AC7-2-verify_go.log" > "$E/AC9-sequence.txt" 2>&1 <<'PY'
import re,sys
log=open(sys.argv[1],errors="replace").read()
anchors=[
 "── AILANG_BIN=", "(AILANG v0.30.0",
 "── tracked-binary hygiene gate", "✓ 0 binary blobs among",
 "── World mission-input gate", "decision ledger valid:",
 "✓ toolchain floor gate:",
 "── race-detector known-positive control", "WARNING: DATA RACE",
 "── go build ./...",
 "── focused host/evidence named-manifest gate (37 exact top-level tests)",
 "✓ all 37 required top-level evidence tests passed exactly once",
 "── go test ./... -count=1",
 "── go test ./... -count=1 -race -timeout 8m",
 "✓ go gate PASSED",
]
pos=-1; ok=True
for a in anchors:
    i=log.find(a,pos+1)
    print(("OK  " if i>pos else "FAIL"), i, repr(a)); ok = ok and i>pos
    if i>pos: pos=i
n=re.search(r"✓ 0 binary blobs among (\d+) tracked files",log)
rows=re.search(r"decision ledger valid: (\d+) rows",log)
print("tracked_files=",n.group(1) if n else "ABSENT","(must be >0)")
print("ledger_rows=",rows.group(1) if rows else "ABSENT","(must be >0)")
ok = ok and n and int(n.group(1))>0 and rows and int(rows.group(1))>0
print("AC9_SEQUENCE:", "PASS" if ok else "FAIL")
sys.exit(0 if ok else 1)
PY
echo "AC9-sequence rc=$?" | tee -a "$E/AC9-sequence.txt"

# --- AC9 deny-list sub-arm: shimmed GOVERSION must REFUSE with `miscompiles` ---
SH=$E/goshim; mkdir -p "$SH"; REALGO=$(command -v go)
printf '#!/bin/sh\nif [ "$1" = "env" ] && [ "$2" = "GOVERSION" ]; then echo go1.26.1; exit 0; fi\nexec %s "$@"\n' "$REALGO" > "$SH/go"
chmod +x "$SH/go"
PATH=$SH:$PATH ./scripts/verify_go.sh > "$E/AC9-denylist.log" 2>&1
echo "AC9-denylist rc=$? miscompiles_hits=$(grep -c miscompiles "$E/AC9-denylist.log")" | tee "$E/AC9-denylist.rc"
# REQUIRED: rc != 0 AND miscompiles_hits >= 1   (base: rc=127, hits=0 — unreachable)

# --- AC6: retirement complete, with positive and negative controls ---
{
  echo "retired_hits=$(grep -rn -E 'check_driver_fleet|driver-fleet-check|AILANG_FLEET_REPO|REQUIRED_FLEET_PATHS|driver_fleet_scope' host scripts .github | wc -l | tr -d ' ')   # want 0 (base: 26)"
  echo "poscontrol_hits=$(grep -rn -- '--mission-config-check' host scripts .github | wc -l | tr -d ' ')   # want >0 (base: 0)"
  echo "negcontrol_hits=$(grep -rn -- 'zzq-nonsense-literal-9f3' host scripts .github | wc -l | tr -d ' ')   # want 0"
  echo "fixture_present=$(ls host/verifygate/driver_fleet_scope_gate_test.go 2>/dev/null | wc -l | tr -d ' ')   # want 0 (base: 1)"
} | tee "$E/AC6-retirement.txt"

# --- AC6: the retired mode is REFUSED, with prerequisites SATISFIED ---
for arg in --driver-fleet-check --zzq-unknown-mode-9f3 driver-fleet-check; do
  ./scripts/verify_go.sh "$arg" > "$E/AC6-refuse${arg}.log" 2>&1; rc=$?
  echo "$arg rc=$rc names_arg=$(grep -c -- "$arg" "$E/AC6-refuse${arg}.log") saw_preflight=$(grep -c '── AILANG_BIN=' "$E/AC6-refuse${arg}.log") saw_PASSED=$(grep -c '✓ go gate PASSED' "$E/AC6-refuse${arg}.log")"
done | tee "$E/AC6-refusal.txt"
# REQUIRED for --driver-fleet-check and --zzq-unknown-mode-9f3: rc!=0, names_arg>=1,
#   saw_preflight=0 (refused BEFORE the binary preflight), saw_PASSED=0.
# `driver-fleet-check` (non-dash): same if Finding 5's widening was implemented; if not,
#   record the residual verbatim and do NOT claim the class is closed.

# --- AC6/AC1: the surviving modes still dispatch ---
env -u AILANG_BIN /bin/bash scripts/verify_go.sh --mission-config-check > "$E/AC1-missioncfg-nobin.log" 2>&1
echo "AC1 nobin rc=$? rows=$(grep -c 'decision ledger valid:' "$E/AC1-missioncfg-nobin.log") disclosure=$(grep -ci 'not certified' "$E/AC1-missioncfg-nobin.log")" | tee "$E/AC1.rc"
./scripts/verify_go.sh --evidence-manifest-check > "$E/AC6-emc.log" 2>&1
echo "emc rc=$? (want 2, usage line unchanged)" | tee -a "$E/AC6-retirement.txt"

# --- AC4/AC6: frozen-core and scope controls ---
{
  echo "launchd_in_diff=$(git diff --name-only 565d0b2..HEAD -- tools/launchd | wc -l | tr -d ' ')   # want 0"
  echo "diff_control=$(git diff --name-only 565d0b2..HEAD | wc -l | tr -d ' ')   # want >0"
  echo "jobs=$(grep -cE '^  [a-z0-9-]+:$' .github/workflows/ci.yml)  # inspect: want exactly the 2 product jobs"
  echo "launchd_in_ci=$(grep -c 'tools/launchd' .github/workflows/ci.yml)   # want 0 (base: 3)"
  echo "launchd_in_verify_go=$(grep -c 'tools/launchd' scripts/verify_go.sh)   # want 0 (base: 10)"
  echo "verify_scripts_in_ci=$(grep -cE '\./scripts/verify_(go|ail)\.sh' .github/workflows/ci.yml)   # want >=2"
} | tee "$E/AC4-wiring.txt"

# --- AC7(4): remote CI on the exact head, read from the API not the badge ---
HEAD_SHA=$(git rev-parse HEAD)
gh api "repos/sunholo-data/ailang-world/commits/$HEAD_SHA/check-runs" \
  --jq '.check_runs[] | "\(.name)\t\(.status)\t\(.conclusion)"' | tee "$E/AC7-4-checkruns.txt"
# REQUIRED: exactly 2 rows; both `success`; `launchd drivers (bash 3.2)` ABSENT.
```

### Evidence artifacts the executor must leave behind

Under `/tmp/iter167-evidence/`, each `.rc` carrying **name, rc, wall time, budget and
within/OVER_BUDGET verdict**, each `.log` the full output:

| artifact | proves |
|---|---|
| `AC7-1-fullsuite.log` + `.rc` (and `AC7-1-fullsuite-run2.*` **only if** risk (a) fired, plus `AC7-1-isolated-flake.log`) | AC7(1), and the flake disclosure protocol |
| `AC7-2-verify_go.log` + `.rc` | AC7(2) **and all of AC9** — this single log is the behavioural proof that every surviving product check ran |
| `AC7-3-verify_ail.log` + `.rc` | AC7(3), no collateral damage to the AIL lane |
| `AC7-5-verifygate.log` + `.rc` | AC7(5) |
| `AC9-sequence.txt` | the ordered-anchor check over the real log, anchor by anchor with its byte offset |
| `AC9-denylist.log` + `.rc` | the deny-list still refuses (rc≠0, `miscompiles` ≥ 1) — **unreachable at base** |
| `AC6-retirement.txt`, `AC6-refusal.txt`, `AC6-refuse*.log` | 26→0 with a live positive control; the retired mode refused before the preflight |
| `AC4-wiring.txt` | 3→2 jobs, 3→0 and 10→0 launchd references, both product scripts still invoked, empty `tools/launchd` diff with a nonempty control |
| `AC1.rc`, `AC1-missioncfg-nobin.log` | the new mode succeeds with `AILANG_BIN` unset, nonzero rows, disclosure present |
| `mutations/M<N>.{pre.sha256,vet.rc,killer.log,killer.rc,control.log,post.sha256}` — **six files per arm, M1…M9 plus M4a/M4b/M7b/M8b** | each mutation's kill, its `go vet` rc read first, its unrelated-command control, and byte-exact restore |
| `V0-vet.rc`, `V0-fmt.log`, `V0-bashn.rc` | compile fence (`go build ./...` is **not** one — it is rc=0 with a type error in a `_test.go`) |
| `AC7-4-checkruns.txt` | remote CI on the exact head, from the API |

---

## 8. Acceptance criteria I believe need correction (stated, never silently narrowed)

1. **AC9 / M7 / M8 — the step-6 test as specified cannot kill M7.** Banners and commands are
   different lines (measured). **Proposed correction: assert ordered (banner, executable anchor)
   pairs** — Finding 1. This *strengthens* the criterion; nothing is dropped. If the executor
   declines it, M7 will report a green that proves nothing and AC9's non-vacuity argument
   collapses.
2. **AC6 / M9 — the M9 arm cannot be run against the real script inside the mandatory 10 s helper**
   without turning the kill into a deadline expiry, which the design itself says can never count as
   a refusal. **Proposed correction: fixture-based M9 arm with a controlled stop, and a
   `── AILANG_BIN=` absence discriminator on both arms** — Finding 3. Nothing is relaxed: the
   real-tree arm still runs with prerequisites satisfied.
3. **AC6's "no-arg default flow still reaches `✓ go gate PASSED`" positive control is discharged by
   AC7(2)'s recorded log**, not by a second full-gate run inside a Go test — Finding 4. Stated so
   it reads as a deliberate discharge and not a dropped control.
4. **AC6's refusal is one character short of closing the class it names** (non-dash first arguments
   are still discarded — measured). **Proposed correction: widen to any unrecognized first
   argument** — Finding 5.
5. **AC3's "invalid ledger aborts before Go" has no meaningful baseline**: at base the flow never
   reaches a ledger check, so the property is vacuously true. Recorded as **UNMEASURABLE@base**
   rather than as a green. Its post-repair measurement is real and stands.

**Is the design's `Files to Modify/Create` list complete?** I found **one gap**, and I am naming it
rather than expanding scope: the list does not mention **`design_docs/planned/w-de-fork-ci-ownership-sprint-plan.md`**
(this file) or the charter/log records the controller writes at landing. Those are the controller's,
not the executor's. Beyond that the list is **complete and correct** — I searched repo-wide for
every other consumer (Findings 6, 7, 9) and found none: no `cmd/` reference, no other test pinning
retired-block literals, no test that job-counts `ci.yml`, and no branch-protection dependency. **The
executor must not add a file to the list.** If it believes one is needed, it says so and stops.

---

## 9. Out of scope (deferred, named)

- **Cleaning up the `tools/launchd/` residue** (`README.md`, `derive-planner-lane.sh`,
  `mission-template.plist`, `test_mission_stall.sh`, `testdata/planner-lane/`, and the **tracked**
  `mission-control.sh.tmp.astra`). Frozen core, fleet-owned. After this sprint it has **zero World
  consumers** (measured) and is the fleet's to remove. The sprint diff must contain **no path under
  `tools/launchd/`**.
- **Fixing `TestCLIRealSubprocessEpisode`'s 5-second announce deadline.** Pre-existing,
  load-triggered, reproduced twice (V21 and Finding 8). A **candidate queue row**, not a change here.
- **Any attestation about the installed shared runtime, launchd state, central registry, routing or
  watchdog behaviour.** Under `D-WORLD-34` those belong to the shared AILANG driver. The new mode's
  disclosure exists precisely to say World is *not* certifying them.
- **Restoring or vendoring the deleted driver/routing tests in any form.**

**SPRINT_PLAN_PATH**: `design_docs/planned/w-de-fork-ci-ownership-sprint-plan.md`
