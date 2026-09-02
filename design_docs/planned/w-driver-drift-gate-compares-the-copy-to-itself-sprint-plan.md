# Sprint plan — `w-driver-drift-gate-compares-the-copy-to-itself`

- **Design doc:** [`design_docs/planned/w-driver-drift-gate-compares-the-copy-to-itself.md`](w-driver-drift-gate-compares-the-copy-to-itself.md) (494 lines, quorum round 2 full strength)
- **Queue row:** `design_docs/world-mission.md` row **54** — clause-2, ~0.3d, gated on nothing
- **Planner:** mission-control iteration 148, `claude-opus-5` sprint-planner
- **Planned at:** `dev` = `f34cb6d00d5013f12e31823d160070d530aea8d5` (re-derived)
- **Executor worktree:** `/Users/voightkampff/dev/sunholo-data/.wt-world-iter148`, branch `sprint/w-iter148-driver-drift`, at `7d3aa00ef880657cc8a769f6726bb700aa54b5e0`
- **Fleet source measured at:** `~/dev/sunholo-data/ailang` HEAD `722e19c77bd7c5c4b7a8a47c9a31909e5387cc9d`
- **Milestones:** ONE (`MS1`). **Risk:** medium — the risk is not the code volume (~145 changed lines in one file), it is `set -e` semantics and AC vacuity, and the planner measured both.

---

## 0. Platform statement (G6)

Every measurement in this plan was taken on **darwin/arm64** (`Darwin arm64`, `go1.26.6 darwin/arm64`, `zsh` as the tool shell). The script's `#!/usr/bin/env bash` resolves to **GNU bash 3.2.57(1)-release** here — the Apple-shipped bash. The CI leg (`.github/workflows/ci.yml:20` and `:101`) is **`ubuntu-latest`**, i.e. bash 5.x.

**Single-platform greens.** Everything below is a darwin/arm64 green. The following have *per-OS* meaning and are explicitly NOT certified on ubuntu by this plan:

| Item | Why it is per-OS |
|---|---|
| Empty-array + `set -u` (`P8`) | bash 3.2 errors `A[@]: unbound variable`; bash 5.x does not. Planner could not run bash 5 (no `/opt/homebrew/bin/bash`, `ls` rc=1). |
| The CI loud-skip branch firing | Depends on GitHub Actions exporting `CI=true` and on `$HOME/dev/sunholo-data/ailang` being absent on the runner. Neither can be observed from this rig; both are inferred. **Labelled unestablished (§7).** |
| `TestCLIRealSubprocessEpisode` flake (`P3`) | Observed 1-of-2 on darwin under full-suite load. Behaviour on ubuntu unknown. |

---

## 1. Milestone breakdown

**ONE milestone.** The design touches exactly one file and the success-line relabel (§4.3) is a single line inside the same edit; splitting it would be theatre.

| ID | Name | Files | Changed lines | Depends on |
|----|------|-------|---------------|------------|
| **MS1** | Fleet-comparison arm in `scripts/verify_go.sh`: `check_driver_fleet()` + `REQUIRED_FLEET_PATHS` + `--driver-fleet-check` isolated mode + main-flow `if`-form call site + working-tree success-line relabel | `scripts/verify_go.sh` (MODIFY) | ~**+144 / ~1 modified** | none |

`.github/workflows/ci.yml` is **NOT** touched (C2's "only if genuinely required" — it is not: `V17`/`:166` already runs `./scripts/verify_go.sh`, re-derived).

### Day plan (~0.3 day)

| Slice | Work | ~Cost |
|---|---|---|
| T1 | Insert `AILANG_FLEET_REPO`, `REQUIRED_FLEET_PATHS`, `check_driver_fleet()` after line **109** (the `}` closing `check_evidence_manifest`). Apply the **P4 enumerator fix**. | 0.05d |
| T2 | Insert the `--driver-fleet-check` isolated-mode block at line **119** (after the `--evidence-manifest-check` block's `fi` at 118, before the `AILANG_BIN` gate at 120). Insert the main-flow `if`-form call site after the working-tree success line. Relabel that line `(working-tree arm)`. | 0.05d |
| T3 | Run the **13-criterion battery** (§4). Each isolated-mode arm is sub-second; AC8 is ~1s. | 0.05d |
| T4 | Run the **10-arm mutation drill** (§5) with landed/parses/red-set/blast-radius per arm. | 0.10d |
| T5 | Final clean re-run + integrity check (sha256 of `tools/launchd/*` unchanged, porcelain scoped to driver paths empty). | 0.05d |

---

## 2. Planner findings — what my own measurement refuted or shrank

Ordered by how much they change the work. **On these, this plan wins over the design doc.** Everywhere else the design doc wins.

### P1 — HIGH, REFUTES the doc. Every row of the doc's AC baseline table is wrong, and it makes five ACs half-vacuous.

The doc's §2 "AC baselines" table says, for AC1–AC7: *"Gate exits 0, prints `working tree matches HEAD`"*. That is the baseline of the **full** `./scripts/verify_go.sh`, not of the command each AC actually runs.

Measured on the unmodified tree, exit codes captured without a pipe (`cmd > /tmp/out 2>&1; rc=$?`), all seven arms:

```
$ AILANG_FLEET_REPO=<fixture> ./scripts/verify_go.sh --driver-fleet-check
rc=1
✗ AILANG_BIN is unset — host/replay tests would t.Skip() silently and this gate would be false-green.
  Export the pinned released binary, e.g. AILANG_BIN=/tmp/ailang-v0300/ailang
```

`--driver-fleet-check` is not yet a recognised mode, so `$1` falls past the `--evidence-manifest-check` branch (line 111) into the `AILANG_BIN` gate (line **120**) and the script exits **1**. `working tree matches HEAD` is never printed; the string is absent from all seven baseline captures.

**Consequence.** AC1, AC3, AC4, AC5, AC6 all assert `exit != 0`. That clause is **SATISFIED AT BASE — for the wrong reason.** It is a vacuous half: it would remain satisfied against an implementation that placed the isolated mode *after* the `AILANG_BIN` gate, or that returned 1 unconditionally. Only the string clause discriminates. AC2 and AC7 assert `exit == 0` and genuinely FAIL at base (rc=1) — those two are sound.

**Repair (binding).** Every AC that asserts `exit != 0` gains a same-call negative clause: **output must NOT contain `AILANG_BIN is unset`**. Paired with the known-positive that the base run *does* contain it, this makes the non-zero attributable to the arm. Applied in §4.

### P2 — HIGH, REFUTES the doc. AC8 as written cannot fire: `touch` does not dirty git.

The doc's AC8 is `touch tools/launchd/mission-control.sh && AILANG_BIN=<pinned> ./scripts/verify_go.sh`, baseline *"the existing arm fires"*.

Measured, with a known-positive control on the **same path** in the **same call**:

```
touch tools/launchd/mission-control.sh
git status --porcelain -- tools/launchd/ scripts/mission_decisions.sh   ->  []          # EMPTY
printf '\n# probe\n' >> tools/launchd/mission-control.sh
git status --porcelain -- tools/launchd/ scripts/mission_decisions.sh   ->  [ M tools/launchd/mission-control.sh]
```

`touch` changes only mtime; git refreshes the stat cache, finds identical content and reports nothing. So AC8 as written runs the full script against a **clean** tree, the working-tree arm passes, and AC8's assertions **fail**. It is not the pass-at-base regression guard the doc claims — it is a broken instrument. Repaired in §4 (AC8).

### P3 — HIGH. AC8's chosen gate is not reliably green at base; it is a 1-of-2 coin flip on a 3-minute run.

| Run | Command | rc | Wall | Cause |
|---|---|---|---|---|
| #1 | `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh` (clean tree) | **1** | 179s | `--- FAIL: TestCLIRealSubprocessEpisode (8.66s) … daemon announcement timed out` |
| #2 | identical | **0** | 216s | — |
| control | `go test ./cmd/ailang-worldd/ -count=1 -run TestCLIRealSubprocessEpisode` | 0 | 1s | passes in isolation |

This is the mission's own known load-flaky wall-clock test (charter row 19's carry note: *"read WHICH test failed, never the exit code"*). A gate that is red at base half the time measures the repo, not the change.

**Repair (binding).** AC8 is re-scoped to the driver gate itself, which is reached **before** the go legs. Measured both arms:

| Arm | rc | Wall | `FATAL: DRIVER DRIFT (D-WORLD-DRIVER-1)` | `go build ./...` |
|---|---|---|---|---|
| dirty (real content change) | 1 | **1s** | present | **absent** |
| clean (control, same file) | 0 | 216s | absent | present |

Asserting *both* — the FATAL present **and** `go build ./...` absent — proves the script aborted **at the gate**, is immune to the downstream flake, and is 180× cheaper.

### P4 — MEDIUM. A real fail-open in §4.1: Phase 1's enumerator is `git ls-files` while everything else in the arm is HEAD-based.

The arm's whole claim is that it compares the **committed** copy. Phase 3 asks `git cat-file -e "HEAD:$path"`. The doc's own domain measurement `V32` used `git ls-tree -r --name-only HEAD`. But Phase 1's code enumerates with `git ls-files` (index + worktree). On a clean tree the two agree — planner re-derived, `diff` of the two sorted lists rc=0, both 6 paths — so the divergence is invisible today.

Measured hole, in a **throwaway clone** (`git clone . /tmp/p148/clone`; the repo and worktree indices were never touched):

```
# control, unmutated clone, real fleet:
  tools/launchd/mission-control.sh (local c8c92a5c… != fleet 253154ad…)   <- compared

$ git rm --cached tools/launchd/mission-control.sh      # in HEAD, out of index
  ls-files count 6 -> 5 ; ls-tree HEAD count still 6
# re-run: grep 'mission-control.sh (local' -> rc=1  (SILENTLY DROPPED)
# unclassified count unchanged at 41; compared = 5, so the liveness control does NOT fire
```

Phase 1 loses the driver's main file; Phase 3 skips it too (HEAD still has it, so `cat-file -e` succeeds and it `continue`s); nothing reports the loss. In the main flow the existing working-tree arm catches a staged deletion first, but the **isolated `--driver-fleet-check` mode has no such protection** and this is the mode all seven ACs drive.

**Repair (binding, one line).** Phase 1 enumerates with:

```bash
done < <(git ls-tree -r --name-only HEAD -- tools/launchd scripts/mission_decisions.sh)
```

Measured: closes the hole (the file is compared again in the mutated clone) **and** is behaviour-identical on the clean tree — full-output `diff` between fixed and unfixed harness rc=0.

### P5 — MEDIUM. No acceptance criterion exercises the GREEN path. §4.3 specifies the success line and nothing tests it.

An implementation that never emits `✓ fleet-comparison arm: …` passes all eight of the doc's ACs. Planner constructed the green path (clone with `tools/launchd/lib/pin-root.sh` committed locally, plus a controlled fleet carrying the identical blob) and measured:

```
   ✓ fleet-comparison arm: 7 tracked frozen-core files match fleet HEAD 894dccc0… — tracked copy is current (untracked fleet additions not certified)
rc=0
```

`compared` = **7**, exactly as §4.3 predicts. Added as **AC9**.

### P6 — MEDIUM. A SURVIVING MUTANT: nothing tests the `unclassified` counter.

Residual 4's "41 unclassified" is a load-bearing claim. Planner mutated `unclassified=$((unclassified + 1))` → `+ 0)`: mutant landed (sha differs), parses (`bash -n` rc=0), and **killed nothing** — all 11 criteria in the pre-repair battery passed. Repaired by **AC13**, which asserts the count line `1 unclassified fleet-only paths not certified`; measured to kill the mutant single-arm.

### P7 — MEDIUM. One mutant is NOT single-arm, so "everything else still passes" is invalid for it.

Doc row *"Neuter MISSING_IN_FLEET (Phase 1)"* is presented as killed by AC6. Measured red set is **{AC4, AC6}** — two arms. Mechanism: with the `[ -z "$fleet_blob" ]` branch neutered, the empty-fleet fixture no longer `continue`s, so `compared` reaches 6, the `compared -eq 0` liveness refusal never fires, and AC4 reads `DRIVER DRIFT vs FLEET` instead of `0 comparable driver files`. Every other mutant measured single-arm. Recorded in the table at §5.

### P8 — LOW, latent, per-OS. `set -u` + an empty `REQUIRED_FLEET_PATHS` is a hard failure on bash 3.2.

```
bash 3.2.57  set -euo pipefail; A=(); for x in "${A[@]}"; ...   -> rc=1, "A[@]: unbound variable"
bash 3.2.57  (control) A=(one); same construct                  -> rc=0, prints "one" / REACHED
```

The set has one member today so this is latent, not live. But the day the fleet lands `lib/pin-root.sh` and someone "cleans up" the now-satisfied entry, `verify_go.sh` hard-fails on the **rig** with an opaque error while staying green on ubuntu bash 5. **Recommendation (non-blocking):** a comment forbidding an empty set, or `${REQUIRED_FLEET_PATHS[@]+"${REQUIRED_FLEET_PATHS[@]}"}`. Bash 5 behaviour could not be measured here (`/opt/homebrew/bin/bash` absent) — see §7.

### P9 — LOW. AC2's third clause is a tautology.

AC2 asserts *output does not contain `driver is current`*. §4.3 removed that bare phrase from the entire design; planner grepped all 13 measured output captures — **0 hits**, with a same-pattern known-positive control on a synthesized file (rc=0). No implementation of this design can violate it. **Repair:** assert instead that the SKIP output does not contain **`tracked copy is current`** — a string the implementation genuinely emits on the green path (measured in P5), so the negative has a real positive counterpart. Measured: green path emits it, skip path does not.

### P10 — INFO / OPERATIONAL. **The moment this lands, `./scripts/verify_go.sh` on the rig is RED.** The controller must be told, not surprised.

The faithful harness against the **real** fleet checkout:

```
rc=1
41 × "⚠ unclassified fleet-only path (not tracked, not required): …"
FATAL: DRIVER DRIFT vs FLEET (D-WORLD-DRIVER-1) — the committed copy differs from fleet HEAD 722e19c7…:
  tools/launchd/derive-planner-lane.sh   (local 52674b04… != fleet 918291f9…)
  tools/launchd/mission-control.sh       (local c8c92a5c… != fleet 253154ad…)
  tools/launchd/test_mission_routing.sh  (local 56c9541e… != fleet b46dbb0e…)
```

And when the fleet lands those three, the next run reds again on `MISSING LOCALLY: tools/launchd/lib/pin-root.sh`. §8 of the doc declares this intended (*"that red is the point"*), and CI stays green (loud skip, measured rc=0 on the harness). But CLAUDE.md's hard rule is *"nothing lands red"*, so state it plainly: **the landing gate for this PR is CI, which is green; the rig red is the designed signal that the fleet must commit.** The executor must NOT "fix" it, must NOT edit or absorb the driver (C1), and must NOT weaken the arm to make the rig green.

### P11 — LOW. Base-commit and porcelain trap for the executor.

`dev` is `f34cb6d`; the worktree is at `7d3aa00`. `git diff --name-only 7d3aa00 f34cb6d -- scripts/verify_go.sh tools/launchd/ scripts/mission_decisions.sh .github/workflows/ci.yml` is **empty** — the code baseline is identical, `f34cb6d` adds only the 494-line design doc. But the worktree's `git status --porcelain` is **not** empty at entry: `?? design_docs/planned/w-driver-drift-gate-compares-the-copy-to-itself.md`. An executor asserting a clean porcelain will misread it. Scope integrity checks to `-- tools/launchd/ scripts/mission_decisions.sh` and to sha256.

### P12 — INFO. The C1 mutation carve-out is needed exactly ONCE, and not by any mutation row.

Every mutation row either mutates `scripts/verify_go.sh` (World's own file, in scope) or a **/tmp controlled fleet repo** — including the addition-shaped arm, which the doc phrases as mutating the "fleet repo". Planner measured it working entirely in `/tmp`, touching neither `tools/launchd/` nor `~/dev/sunholo-data/ailang`. The only place a frozen-core file is temporarily modified is **AC8**, and it must restore from a `cp` backup (never `git checkout --`, because the executor's own work is uncommitted by construction). Planner executed this and verified the restore: sha256 `1632b31b0403411986237709d0d7a92102bb22cbb2a7b4fbc51c020cbd0440a8` before and after, porcelain empty after.

### P13 — INFO, G2 clean bill. Every figure in the design doc re-derived and AGREES.

Nothing in the doc's own numbers was refuted. Re-derived first-party at `7d3aa00` against fleet `722e19c7`:

| Doc claim | Re-derived | Agrees |
|---|---|---|
| 6 local driver paths (`V8`, `V32`) | 6 (both `ls-files` and `ls-tree`) | ✅ |
| 48 fleet paths (`V33`) | 48 | ✅ |
| 42 missing locally (`V34`) | 42 | ✅ |
| 0 missing in fleet (`V35`) | 0 | ✅ |
| union 48 (`V36`) | 48 | ✅ |
| 11 commits behind (`V6`) | 11, newest `aca6908bd feat(mission-driver): controller fallback is now a chain…` | ✅ |
| 705 differing lines (`V4`+`V5`) | 615 + 90 = 705 | ✅ |
| 550 / 1075 line counts (`V2`,`V3`) | 550 / 1075 | ✅ |
| blobs `c8c92a5c…` / `253154ad…` (`V21`,`V22`) | identical | ✅ |
| `lib/pin-root.sh` fleet blob `5902a60c…` (`V28`) | identical | ✅ |
| `pin-root.sh` is in the 42 (`V38`) | line 13 of the comm output | ✅ |
| line numbers 111 / 120 / 190 / 200 / 205 / 212 (`V24`,`V41`) | identical; `fi` of the evidence block at 118 | ✅ |
| 41 unclassified (§7 residual 4) | 41, counted from the harness against the real fleet | ✅ |
| §4.3's predicted `compared` = 7 | 7 | ✅ |
| `V20` per-path MATCH/DIFFER verdicts | identical, all 6 | ✅ |

**Also confirmed by first-party execution: the design's §4.1/§4.2 code is correct as written under bash 3.2 + `set -euo pipefail`.** Planner built a byte-faithful harness of §4.1 and §4.2 and ran it: all eight refusal branches fire, the CI loud skip returns 0 through both call sites, and the main-flow site reaches its end with `fleet_rc=2`. **CB1 and CB1-bis are genuinely fixed.** The only code change this plan asks for beyond the doc is P4's one-line enumerator swap.

`${CI:-}` is a brand-new dependency: `grep -rn --include='*.sh' --include='*.yml' 'CI:-\|CI=true' scripts .github tools` → rc=**1** (no hits), with a same-paths same-flags known-positive control (`AILANG_BIN`) → rc=0, 43 hits. Confirms the doc's "no existing consumer" and Residual 6.

---

## 3. What to write, and where

All in `scripts/verify_go.sh`. **Ship the design doc's §4.1 and §4.2 text verbatim except for the single P4 line.** The whole drill is keyed on these exact strings — a re-wording silently voids every arm:

- `FATAL: DRIVER DRIFT vs FLEET (D-WORLD-DRIVER-1)`
- `enumerated 0 comparable driver files`
- `World-tracked paths MISSING IN FLEET`
- `REQUIRED fleet paths MISSING LOCALLY`
- `fleet-comparison arm SKIPPED (fleet checkout absent at`
- `driver currency NOT certified here`
- `fleet source $AILANG_FLEET_REPO is absent`
- `✓ fleet-comparison arm: $compared tracked frozen-core files match fleet HEAD $fleet_head — tracked copy is current (untracked fleet additions not certified)`
- `⚠ unclassified fleet-only path (not tracked, not required): $path`
- `⚠ $unclassified unclassified fleet-only paths not certified (see above)`

**Insertion points (re-derived line numbers on the unmodified file, 311 lines):**

| What | Where |
|---|---|
| `AILANG_FLEET_REPO=` default, `REQUIRED_FLEET_PATHS=(…)`, `check_driver_fleet() { … }` | after line **109** (`}` closing `check_evidence_manifest`), before line 111 |
| `if [ "${1:-}" = "--driver-fleet-check" ]; then … fi` | at line **119** — after the `--evidence-manifest-check` block's `fi` (line 118), **before** the `AILANG_BIN` gate (line 120). This placement is what makes the mode free of `AILANG_BIN`. |
| main-flow `if`-form call site (`fleet_rc`) | immediately after the working-tree success line (line **212**) |
| success-line relabel → `…working tree matches HEAD (working-tree arm)` | line **212** |

**`set -e` discipline (C3 — the plan's highest-risk line).** Both call sites use the `if` form. Do **not** copy the neighbouring `check_evidence_manifest "$2" "$3"; exit $?` shape at lines 116–117: that bare call is benign only because its exit code happens to be preserved, and it is exactly the pattern that burned quorum rounds 1 and 2 here. Verified working in the harness:

```bash
rc=0
if check_driver_fleet; then :; else rc=$?; fi
[ "$rc" -eq 2 ] && rc=0
exit "$rc"
```

---

## 4. Acceptance criteria — with the measured baseline as part of each criterion

Thirteen criteria. **AC1–AC7 and AC9–AC13 drive the isolated `--driver-fleet-check` mode** (sub-second each, no `AILANG_BIN`). **AC8 uses the full script** because the working-tree arm lives only in the main flow — but it is scoped so it exits at the gate in ~1s.

Fixtures (all in `/tmp`, all built from World's own `HEAD`; the real fleet checkout is never modified):

- `F_DIFF` — World's 6 tracked paths, with `# fleet-only drift` appended to `mission-control.sh`
- `F_EMPTY` — `git init` + one empty commit
- `F_PIN` — World's 6 paths **plus** `tools/launchd/lib/pin-root.sh`
- `F_NOTRACK` — World's 6 paths **minus** `tools/launchd/test_mission_routing.sh`
- `F_ADD` — a green-path fleet **plus** a fleet-only `tools/launchd/lib/new-helper.sh`
- `G` — a throwaway clone of the worktree with `tools/launchd/lib/pin-root.sh` committed locally, plus `GFLEET` matching it exactly (this is the only way to reach the green path today)
- `F_NOPIN`, `F_PINDIFF` — `G`'s tree with `pin-root.sh` removed / altered in the fleet

| AC | Command (from the worktree unless noted) | Assert | **BASELINE on unmodified tree (measured)** |
|----|---|---|---|
| **AC1** | `AILANG_FLEET_REPO=$F_DIFF ./scripts/verify_go.sh --driver-fleet-check` | `rc != 0` **and** contains `DRIVER DRIFT vs FLEET` **and** names `tools/launchd/mission-control.sh` **and** does NOT contain `AILANG_BIN is unset` ← **P1 repair** | rc=**1**, output = `✗ AILANG_BIN is unset …`. `DRIVER DRIFT vs FLEET` absent (rc=1). **The `rc != 0` half is satisfied at base for the wrong reason; the new negative clause is what makes it fail at base.** |
| **AC2** | `CI=true AILANG_FLEET_REPO=/nonexistent ./scripts/verify_go.sh --driver-fleet-check` | `rc == 0` **and** contains `SKIPPED` **and** does NOT contain `tracked copy is current` ← **P9 repair** | rc=**1**, `SKIPPED` absent. Genuinely FAILS at base. |
| **AC3** | `env -u CI AILANG_FLEET_REPO=/nonexistent ./scripts/verify_go.sh --driver-fleet-check` | `rc != 0` **and** contains `fleet source` and `absent` **and** NOT `AILANG_BIN is unset` | rc=**1** from the `AILANG_BIN` gate; `fleet source` absent. |
| **AC4** | `AILANG_FLEET_REPO=$F_EMPTY ./scripts/verify_go.sh --driver-fleet-check` | `rc != 0` **and** contains `0 comparable driver files` **and** NOT `AILANG_BIN is unset` | rc=**1** from the `AILANG_BIN` gate; `0 comparable driver files` absent. |
| **AC5** | `AILANG_FLEET_REPO=$F_PIN ./scripts/verify_go.sh --driver-fleet-check` | `rc != 0` **and** contains `MISSING LOCALLY` **and** names `tools/launchd/lib/pin-root.sh` **and** NOT `AILANG_BIN is unset` | rc=**1** from the `AILANG_BIN` gate; `MISSING LOCALLY` absent. |
| **AC6** | `AILANG_FLEET_REPO=$F_NOTRACK ./scripts/verify_go.sh --driver-fleet-check` | `rc != 0` **and** contains `MISSING IN FLEET` **and** names `tools/launchd/test_mission_routing.sh` **and** NOT `AILANG_BIN is unset` | rc=**1** from the `AILANG_BIN` gate; `MISSING IN FLEET` absent. |
| **AC7** | `CI=true AILANG_FLEET_REPO=/nonexistent ./scripts/verify_go.sh --driver-fleet-check` | `rc == 0` (CB1 guard: the loud skip must never redden) | rc=**1**. Genuinely FAILS at base. |
| **AC8** *(repaired, P2+P3)* | `cp tools/launchd/mission-control.sh /tmp/p148-bak/mc` ; `printf '\n# ac8 probe\n' >> tools/launchd/mission-control.sh` ; `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh` ; **restore with `cp` from the backup** | `rc != 0` **and** contains `FATAL: DRIVER DRIFT (D-WORLD-DRIVER-1)` **and** does NOT contain `go build ./...` (proves it aborted AT the gate) — plus the clean control run has the complementary pair | **PASSES at base** and must still pass after: dirty rc=**1** in **1s**, FATAL present, `go build ./...` absent; clean control rc=0/216s, FATAL absent, `go build` present. The doc's literal `touch` form does **not** pass at base — measured empty porcelain (P2). |
| **AC9** *(new, P5)* | from `G`: `AILANG_FLEET_REPO=$GFLEET ./scripts/verify_go.sh --driver-fleet-check` | `rc == 0` **and** contains `tracked copy is current` | rc=**1** (`AILANG_BIN` gate); success line absent. |
| **AC10** *(new, addition-shaped)* | from `G`: `AILANG_FLEET_REPO=$F_ADD ./scripts/verify_go.sh --driver-fleet-check` | `rc == 0` **and** contains `unclassified fleet-only path (not tracked, not required): tools/launchd/lib/new-helper.sh` **and** contains `tracked copy is current` | rc=**1** (`AILANG_BIN` gate); neither string present. |
| **AC11** *(new)* | from `G`: `AILANG_FLEET_REPO=$F_NOPIN ./scripts/verify_go.sh --driver-fleet-check` | `rc != 0` **and** `MISSING IN FLEET` naming `tools/launchd/lib/pin-root.sh` **and** NOT `AILANG_BIN is unset` | rc=**1** from the `AILANG_BIN` gate; string absent. |
| **AC12** *(new)* | from `G`: `AILANG_FLEET_REPO=$F_PINDIFF ./scripts/verify_go.sh --driver-fleet-check` | `rc != 0` **and** `DRIVER DRIFT vs FLEET` naming `tools/launchd/lib/pin-root.sh` **and** NOT `AILANG_BIN is unset` | rc=**1** from the `AILANG_BIN` gate; string absent. |
| **AC13** *(new, P6)* | from `G`: `AILANG_FLEET_REPO=$F_ADD ./scripts/verify_go.sh --driver-fleet-check` | `rc == 0` **and** contains `1 unclassified fleet-only paths not certified` | rc=**1** (`AILANG_BIN` gate); string absent. |

**No AC in this plan is red at base for a reason unrelated to the change.** AC1–AC7 and AC9–AC13 are red at base *by design* (the fix is absent). AC8 is green at base and is a regression guard. The two ACs the doc got right (AC2, AC7) were already non-vacuous; the five it got half-vacuous (AC1, AC3, AC4, AC5, AC6) are repaired by the `AILANG_BIN is unset` negative clause; AC8 is repaired twice over (P2 + P3).

**Battery result on a byte-faithful harness of §4.1/§4.2 (planner-executed):** all **13** PASS.

---

## 5. Mutation plan — 10 arms, blast radius measured

Every arm below was **executed by the planner** against a byte-faithful harness of §4.1/§4.2 and the 13-criterion battery. Per G5, each row records: that the mutant **LANDED** (sha256 differs from base `675d231c1191`), that it still **PARSES** (`bash -n` rc=0 — a mutant that cannot parse is not a guard that fired), the **measured red set**, and the **blast radius**.

Neuter with `if false && <cond>`, never by deleting the block. Restore `scripts/verify_go.sh` from a `cp` backup after every arm.

> **Mutation-engine warning, learned the hard way.** The planner's first attempt used `awk sub(pat,rep)`, which treats the pattern as a regex — all eight mutants produced byte-identical files and the "did it land" guard caught it. Use literal `index()`/`substr()` replacement, or `sed` with everything escaped, and **always** verify sha256 changed before believing a red or a green.

| # | Mutation | Killed by | Which WRITE the assertion reads | Blast radius (MEASURED red set) |
|---|---|---|---|---|
| **M1** | Phase 1 DIFFERING: `if false && [ "$local_blob" != "$fleet_blob" ]` (1st occurrence) | AC1 | the `differing` string, appended **only** inside the `!=` branch; the FATAL `printf`s that variable | **{AC1}** — single-arm. "Everything else still passes" is valid. |
| **M2** | Phase 1 MISSING_IN_FLEET: `if false && [ -z "$fleet_blob" ]` (1st) | AC6 | the `missing_in_fleet` string, appended only inside the `-z` branch | **{AC4, AC6}** — **TWO arms (P7).** Neutering removes the `continue`, so `compared` reaches 6 on the empty fleet, the liveness refusal never fires, and AC4 reads `DRIVER DRIFT vs FLEET`. "Everything else still passes" is **INVALID** for this row. |
| **M3** | Phase 2 MISSING_LOCALLY: `if false && ! git cat-file -e "HEAD:$path"` | AC5 | `missing_locally`, written only in the required-absent branch | **{AC5}** — single-arm. |
| **M4** | Phase 2 MISSING_IN_FLEET: `if false && [ -z "$fleet_blob" ]` (2nd) | **AC11** | `missing_in_fleet`, Phase 2 branch. AC6 cannot reach it (AC6's missing path is World-tracked, Phase 1) — this is why AC11 exists | **{AC11}** — single-arm. |
| **M5** | Phase 2 DIFFERING: `if false && [ "$local_blob" != "$fleet_blob" ]` (2nd) | **AC12** | `differing`, Phase 2 branch. AC1 cannot reach it | **{AC12}** — single-arm. |
| **M6** | Path-liveness: `if false && [ "$compared" -eq 0 ]` | AC4 | the `0 comparable driver files` FATAL, written only in that branch | **{AC4}** — single-arm. |
| **M7** | CI loud skip: `if false && [ -n "${CI:-}" ]` | AC2 | the `SKIPPED` line; with the branch dead the arm falls through to the rig refusal and rc becomes 1 | **{AC2}** — single-arm. Note AC7 survives it in the battery only because AC7 duplicates AC2's rc clause; keep both, AC7 is the CB1 guard by name. |
| **M8** | Rig typed refusal: `if false && [ -z "${CI:-}" ]` | AC3 | the `fleet source … absent` FATAL, written only in that branch; with it dead the function reaches `return 0` | **{AC3}** — single-arm. |
| **M9** | **Counter mutant (P6):** `unclassified=$((unclassified + 0))` | **AC13** | the count line `⚠ N unclassified …` reads `$unclassified`; the per-path `⚠` report is a *sibling* write and survives, so AC10 alone cannot see it | **{AC13}** — single-arm. **This mutant SURVIVED the doc's own 9-row table**; AC13 is what kills it. |
| **M10** | **ADDITION-shaped (proves the arm LOOKS, not merely FIRES):** Phase 3 enumerates the **local** tree instead of the fleet — `done < <(git ls-tree -r --name-only HEAD -- …)` (drop `-C "$AILANG_FLEET_REPO"`) | **AC10** + AC13 | the unclassified report is produced only by Phase 3 enumerating the **fleet** tree; a local enumeration yields nothing fleet-only | **{AC10, AC13}** — two arms, both Phase-3 observers. This is the arm the doc's round-1 mutant missed. |

**All 10 mutants: LANDED (sha256 differs), PARSED (`bash -n` rc=0), 2 diff lines each.** No mutant produced an unparseable file, so no red in this table is a syntax error masquerading as a guard.

**A removal-only table would have been hollow here.** M10 is the addition-shaped arm and it is the only one that can distinguish "Phase 3 exists" from "Phase 3 looks at the right tree".

---

## 6. Executor constraints

- **Git writes are FORBIDDEN.** No `add`, `commit`, `stash`, `checkout`, `restore`, `reset`, `branch`, `worktree`. Reads (`show`, `status`, `diff`, `log`, `ls-files`, `ls-tree`, `rev-parse`, `cat-file`) are required and allowed. The controller commits.
- **`tools/launchd/*` is FROZEN CORE (C1).** Never edit, absorb, vendor or sync the driver. The one temporary modification is **AC8**, which must:
  1. `cp` the file to a `/tmp` backup **first**;
  2. append a real content change (`touch` does nothing — P2);
  3. restore with `cp` from that backup, **never `git checkout --`** (the executor's own work is uncommitted by construction);
  4. assert the sha256 matches `1632b31b0403411986237709d0d7a92102bb22cbb2a7b4fbc51c020cbd0440a8` and `git status --porcelain -- tools/launchd/ scripts/mission_decisions.sh` is empty afterwards.
- **No mutation row needs the C1 carve-out (P12).** Every fleet-side fixture is a `/tmp` repo. Never modify `~/dev/sunholo-data/ailang`.
- **Files forbidden:** `tools/launchd/*` (except AC8's temporary probe), `scripts/verify_ail.sh`, `scripts/mission_decisions.sh`, any `.ail` file, `go.mod`/`go.sum`, `.github/workflows/ci.yml` (not required — P-note in §1).
- **Do not weaken the existing arm (C4).** Lines 200–212's `driver_tracked` path-liveness control and `driver_drift` working-tree diff are untouched; only the success-line *text* at 212 gains ` (working-tree arm)`.
- **Message text is FROZEN.** See §3's string list.
- **Assert substrings, never whole lines and never absolute paths.** The FATAL messages interpolate `$AILANG_FLEET_REPO` and `$fleet_head`, which are rig-specific.
- **Exit-code discipline.** Capture without a pipe: `cmd > /tmp/out 2>&1; rc=$?`. `${PIPESTATUS[0]}` is bash-only and silently empty in zsh. Never `|| echo 0` inside `$(...)`. `grep`'s rc 1 = no match, 2 = no such file; `| wc -l` throws that away. Quote glob-shaped flag values (`--include='*.sh'`).
- **The rig red is expected (P10).** After MS1 lands, `./scripts/verify_go.sh` on the rig reds with three DIFFERING driver paths. Do not fix it, do not touch the driver, do not soften the arm. CI is the landing gate and it loud-skips.
- **Environment:** `AILANG_BIN=/tmp/ailang-v0300/ailang` (the *file*; `/tmp/ailang-v0300` alone is a directory). `export PATH=/opt/homebrew/bin:$PATH` before any `gh`. No `timeout(1)` on this box.

---

## 7. What the planner could NOT establish

Labelled as unestablished, not assumed:

1. **That the CI loud-skip actually fires on `ubuntu-latest`.** It requires GitHub Actions to export `CI=true` *and* `$HOME/dev/sunholo-data/ailang` to be absent on the runner. `grep -n "CI=" .github/workflows/ci.yml` → rc=1 (no explicit set), so the arm depends on the runner's implicit default. Neither leg was observed from this rig. **This is the single assumption on which "CI stays green" rests.** If it is wrong, CI reds on the rig-refusal branch. Cheapest mitigation available to the executor: none locally — it is observable only on the first CI run of the PR. Flagging it so a CI red on the merge is read correctly and not mistaken for a code defect.
2. **bash 5.x behaviour of the empty-array construct (P8).** No bash ≥4 on this rig.
3. **Whether `TestCLIRealSubprocessEpisode` also flakes on ubuntu.** Observed 1-of-2 on darwin under load; not measured on Linux. AC8's repair makes this moot for this sprint.
4. **The long-run flake rate.** Two full-suite runs is a small sample; 1-of-2 is a lower bound on the nuisance, not a rate.
5. **Whether `~/.ailang-driver-pin/world` participates in what launchd actually executes.** The doc's `V13`/`V15` establish the pin worktree exists and that nothing in this repo references it; the plan does not depend on it either way, and the arm does not measure it.

---

## 8. Success metrics

- MS1 green: all **13** acceptance criteria pass; AC8 still passes.
- All **10** mutation arms executed, each proven landed (sha256) and parsing (`bash -n` rc=0), each with its measured red set recorded — including M2's two-member set and M9/M10's Phase-3 pair.
- `scripts/verify_go.sh` passes `bash -n` rc=0 at every checkpoint.
- `tools/launchd/` sha256 set unchanged end-to-end; `git status --porcelain -- tools/launchd/ scripts/mission_decisions.sh` empty at exit.
- CI green on the PR (both jobs, SHA-addressed on the merge commit).
- The rig red is **declared**, named in the PR body, and attributed to the fleet — never fixed locally.
