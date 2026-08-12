# Sprint plan — `w-ail-gate-module-pin` (queue item 12, clause-1)

**Item**: queue item 12 — *pin the Leg-1 module SET in `scripts/verify_ail.sh` itself*.
**Status**: PLANNED · **NO SPLIT** · one milestone, **5 tasks / 5 commits**
**Design doc**: [`design_docs/planned/w-ail-gate-module-pin.md`](w-ail-gate-module-pin.md) (`d201a1e`,
quorum-cleared over 2 rounds, narrow-refinement carve-out applied to round 2)
**Base**: `dev` @ `c53db58`, clean tree (`git status --porcelain` empty, V0)
**Planner**: mission-control iteration 77, opus sprint-planner, first-party measurement on this rig
**Executor lane**: `codex:gpt-5.6-sol` under `--sandbox workspace-write`.
**THE CONTROLLER MAKES ALL COMMITS.** The executor has **no git write permission** — no `git
commit`/`add`/`push`/`checkout`/`stash`/`restore`, no `gh pr`. The controller reconstructs the 5
commits from cumulative `.snap/M<k>/` snapshots. **Restores are `cp` from `/tmp/w12_backup/`**,
never `git checkout -- <file>`: the mutated file is uncommitted by construction, so
`git checkout --` deletes the executor's work and the sha256 assertion then reports the disaster
rather than preventing it.

**Headline price: ≈0.75 day, i.e. 1.5× the doc's ~0.5 day** (§9 of the doc already flags itself as
"top of the band"). The driver is not the sweep — it is that **two of the doc's specified
mechanisms do not work as written on this rig** and one refusal branch is missing from §6. §6
prices it; §6.1 refuses the split and says why each seam fails structurally.

**The entire script change was PROTOTYPED, RUN, and MUTATED end-to-end this session** (V6–V16).
The frozen block is reproduced verbatim in §3/T1. Nothing in T1 is exploratory.

---

## 0. Planner's first-party verification, and what it REFUTED

Every controller premise was re-measured at `c53db58` before anything was planned. **All controller
premises survived** (§0.1). But the design doc did not: **five of its claims are corrected by
measurement** (§0.2), and two of them would have produced a *silently wrong implementation* — the
doc's own null-case arm cannot fire against a faithful reading of §4.2.

Full command/observed-output table in §8.

### 0.1 Controller premises — all CONFIRMED

| Premise | Verdict |
|---|---|
| Baseline `verify_ail.sh` at `c53db58` → rc=0, `4/4 … across 11 module(s)`, `14 required named tests`, `9/9 steps` | **CONFIRMED** (V1) |
| `scripts/verify_ail.sh` line citations `:167 :176 :233 :238 :243` are live | **CONFIRMED, exact** (V2) |
| `host/verifygate/` contains exactly one file; `runGate` hardwires `repoRoot` for both the script path and `cmd.Dir` (`:52-56`) → `runGateAt` genuinely required | **CONFIRMED** (V3, V4) |
| `LEG1_MODULES` unallocated (0 hits) with a same-call known-positive control | **CONFIRMED** — 0 hits rc=1; control `EXACT_TOTAL_VERIFIED` → 4 hits rc=0 (V5) |
| local `go` is 1.26.4 → `verify_go.sh` FATALs without `GOTOOLCHAIN=go1.25.6`; base condition | **CONFIRMED** — `verify_go.sh: FATAL: active toolchain go1.26.4 miscompiles host/store/scan.go` (V17) |
| `host/broker` ~18% base flake is out of scope | **ACCEPTED** — and *not reachable from this sprint*: no task runs `./host/broker`; only AC6's whole-repo `verify_go.sh` can meet it |

### 0.2 Five corrections to the design doc — each MEASURED, each changing the plan

#### (i) **REFUTED: §4.2's write-then-guard ordering kills the script before its own null-case message.** This is the load-bearing one.

§4.2 step 2 prescribes: *"write the actual set with `printf '%s\0' "${mods[@]}"` and the expected
set with `printf '%s\0' "${LEG1_MODULES[@]}"`, pass both through `LC_ALL=C sort -z`, **assert both
files non-empty** … then `cmp -s`"* — i.e. **write, then guard**.

`verify_ail.sh:30` sets `set -uo pipefail`. The rig's `/usr/bin/env bash` is **GNU bash
3.2.57(1)** (V7). Under bash 3.2 + `set -u`, `"${arr[@]}"` on an **empty** array is an
**unbound-variable ABORT**, not an empty expansion. Measured inside a script of exactly this shape
(V8), with the non-empty array as the same-shape known-positive control:

```
count=0
--- now the printf the doc prescribes ---
/tmp/w77_setu.sh: line 6: A[@]: unbound variable          <- script rc=1, next line NEVER REACHED
```
```
count=2
REACHED_AFTER_PRINTF rc=0                                 <- control: same script shape, non-empty
0000000    x  \0   y  \0
```

So under a faithful §4.2 implementation, `LEG1_MODULES=()` produces a **bash unbound-variable
error**, not `✗ … empty`. **`MUT-EMPTY-ALLOWLIST` as §5.4/§6 specify it — *"requires an early loud
`empty` failure"* — would FAIL against a correct-to-the-doc implementation.** The S6 branch would
be unreachable, and its committed arm would be asserting on a message the script can never print.

**Fix, adopted and measured (V11):** the two non-empty guards move **before** the `printf` writes
and read `${#arr[@]}`, which **is** `set -u`-safe on an empty array (`count=0` printed fine, V8).
Measured after the fix: `LEG1_MODULES=()` →
`✗ LEG1_MODULES allowlist is empty — the membership compare would pass vacuously; failing loudly`,
rc=1. The guard is on the ARRAY rather than the FILE; that is strictly earlier and strictly
stronger (a file can only be empty if its array was, or if `printf` itself failed).

**Executor consequence:** `MUT-EMPTY-ALLOWLIST` must **empty** the array (`LEG1_MODULES=()`), never
delete its declaration — deleting it makes `${#LEG1_MODULES[@]}` itself an unbound-variable abort,
which is loud but is not the designed message and would let an arm that asserts only `rc != 0` pass
for the wrong reason.

#### (ii) **§6's mutation table is INCOMPLETE (rule 3j): the `diff -u` diagnostic is a separately neuterable branch.**

§6 lists six mutations covering set-equality, the two null cases, the isolation control, and the
neutering. **None of them touches the diagnostic.** Measured (`MUT-DIAG-SILENT`, V13): replace the
`diff -u …` line with `:` and run the stray mutation —

```
rc=1
✗ swept .ail module set differs from the LEG1_MODULES allowlist — an intentional
  module add/remove must edit LEG1_MODULES in scripts/verify_ail.sh in the SAME commit
--- does output name the offending path? (expect 0) --- 0
```

The gate still **refuses**, still prints the message §4.2 quotes, and every §6 row still passes —
while the half a human actually acts on (*which* path is offending) is silently dead. That is
**rule 3k's exact hazard**, and the item's own thesis at one level down: a refusal that names
nothing is a guard, not a gate. **Added as branch B4 with its own mutation (M7) and its own
composed self-satisfaction arm (M9).**

#### (iii) **§8's conflict surface is INCOMPLETE: `TestNoRigAbsolutePaths` scans the new file.**

`host/verifygate/ail_binary_gate_test.go:TestNoRigAbsolutePaths` globs **`host/verifygate/*.go`**
and `t.Errorf`s on any file containing the assembled needles `"/tmp"+"/ailang"`, `"/Us"+"ers/"`,
`"/home"+"/runner/"`. The new `module_manifest_gate_test.go` lands **inside that glob**. Negative
existence with two same-call controls (V14): `TestNoRigAbsolutePaths` appears **0** times in the
design doc; controls `TestInScriptControl` → **4**, `TestCorpusPredicateMatchesShellAnchors` → **2**.

**Executor consequence:** the new file must derive every path from `repoRoot` and `t.TempDir()` and
must contain **no** rig-absolute literal — not in code, not in a comment. `t.TempDir()` on macOS
returns `/var/folders/…`, so no literal is needed; a helper reaching for `os.MkdirTemp("/tmp", …)`
would red a **landed** criterion.

#### (iv) **§5.2's "four-item copy set" is a DIRECTORY copy that drags gitignored compile caches.**

`cp -R world/ design_docs/sketches/` copies the `.ailang/cache/` trees: **95 files**, not 4 items
(V9). Those caches are gitignored (`.gitignore:3 **/.ailang/`, `git check-ignore` rc=0 on
`world/.ailang/cache`, rc=1 on `world/types.ail` — control in the same call) and the **live** cache
contains entries for modules that **no longer exist** (`world___stray`,
`world__transitionregistry`) — i.e. the arms would inherit a stale, untracked, unreviewed input.

Measured (V10, V16): an **`.ail`-file-scoped** copy set is **13 files** (11 modules +
`scripts/verify_ail.sh` + `scripts/testdata/ailang_release_observed.txt`), the pristine control is
green at `11 module(s)`, and the cost is identical (**1.465 s** vs 1.638 s wall). **Freeze the copy
set as file-scoped, never `cp -R` of a directory.**

#### (v) **The `diff -u` header is undiagnosable under process substitution.** (rule 3k)

The measured prototype's first output was:

```
--- /dev/fd/63	2026-08-12 18:32:51
+++ /dev/fd/62	2026-08-12 18:32:51
```

A human cannot tell which side is the allowlist and which is the sweep. `diff -u --label` is
supported by this rig's `Apple diff (based on FreeBSD diff)` **and** by GNU diff (V15). Adopted and
measured (V16):

```
--- expected: LEG1_MODULES in scripts/verify_ail.sh
+++ actual:   .ail files swept under ROOTS
+world/_stray_manifest_probe.ail
```

#### (vi) Not a refutation — a trap for the arm author. **The Leg-1 BANNER contains the substring `ai-check`.**

§5.4 requires the stray arm to assert the output does **not** contain an `ai-check` of the probe.
A bare-substring anti-assertion measures the banner `── Leg 1: ai-check required-check manifest`,
not the sweep. Measured on the stray log (V12): bare `ai-check` → **1**; the sweep's own line form
`^   ai-check ` → **0**; same needle on the pristine log (control) → **11**. **Anchor the
anti-assertion on `\n   ai-check ` (three leading spaces), and pair it with the pristine control's
11 in the same test** — an absence assertion whose control never fires is decoration.

#### (vii) Scope, honestly re-measured against the prototype.

| | doc | measured prototype |
|---|---|---|
| `scripts/verify_ail.sh` changed lines | ~40–50 | **+66 / −4** (V6) |
| `host/verifygate/module_manifest_gate_test.go` | ~230–300 LOC | **~345 LOC** (4 committed arms, not 2–3; §6) |

---

## 1. What this sprint is, and what it explicitly does NOT do

**In scope — three files:**

| File | Purpose |
|---|---|
| `scripts/verify_ail.sh` | **MODIFIED.** `LEG1_MODULES` identity allowlist + Leg-1 restructure to *enumerate → compare → consume*. +66/−4 lines. |
| `host/verifygate/module_manifest_gate_test.go` | **NEW.** `package verifygate`. Isolated-root harness + **four** committed refusal arms. ~345 LOC. |
| `design_docs/verification/w-ail-gate-module-pin/mutations.md` | **NEW.** The 10-mutation / 20-arm transcript, same form as `trb2-mutations.md`. |

**Explicitly NOT in this sprint** — §10 of the doc is binding, and every line of it holds:

1. **No `EXACT_TOTAL_MODULES` count literal.** Defeated by the measured add-one-delete-one mutant
   (V8 of the doc; re-measured here as M3/AC4), redundant under set equality, and an item-13 tax.
2. **No CI workflow edit.** Job 1 runs the script at `ci.yml:96` exporting only
   `WORLD_PKG_AILANG_BIN` (`:93`); job 2 exports `AILANG_BIN` (`:144`). The membership compare is
   binary-independent, so no lane divergence is introduced (V18).
3. **No `host/boundary` edit** (its `wantFileCount = 1` landmine is scoped to that directory).
4. **Do not touch** `REQUIRED_VERIFIED`, `EXACT_TOTAL_VERIFIED`, `EXACT_TOTAL_TESTS`, the `:233`
   zero-guard, Leg 2, Leg 3, or `verify_world_package.sh`.
5. **No env-overridable knobs.** Gate policy stays hardcoded per `verify_ail.sh:12-14`.
6. **No sweeping** of `packages/world-core/*.ail` or `host/replay/testdata/transition_fixture.ail`.
7. **No committed test may mutate the live checkout.** This is what BLOCKED the doc's first draft:
   `verify_go.sh:108` runs `go test ./... -count=1` with **no `-p 1`**, and `host/boundary`'s
   `enumerateAIL` walks *and reads* the live `world/` on every run. Every committed arm runs a
   **real copied** script in a private `t.TempDir()` root. Only the sequential one-shot ACs
   (AC2–AC4) and the one-shot `MUT-NEUTER-CMP` touch the live tree, never concurrently with
   `go test`.
8. **Do not fix the `host/broker` ~18% base flake** (queue item 16). No task in this sprint runs
   `./host/broker`; only AC6's whole-repo `verify_go.sh` can meet it. **A `host/broker` red is not
   this sprint's regression** — read *which test* failed, rerun, attribute, move on.

**The `:233` zero-guard becomes shadowed, and that is correct.** After the restructure, branch B1
(`${#mods[@]} -eq 0`) fires before the sweep can leave `checked` at 0, so `:233` is unreachable in
the empty case. §10 forbids touching it; it stays as **defense-in-depth** and becomes reachable
again under `MUT-NEUTER-CMP`. Say this in the commit message so a reviewer does not read it as dead
code introduced by this change.

---

## 2. Acceptance criteria — the doc's AC1–AC7, plus the hold set

The doc's AC1–AC7 stand and are adopted verbatim, with three amendments forced by §0.2:

| AC | Amendment |
|---|---|
| **AC2** | The anti-assertion needle is `\n   ai-check ` (3 leading spaces), **not** the bare token `ai-check` — the banner contains it (§0.2(vi), V12). Pair with the pristine control's 11 occurrences. |
| **AC5** | Add: the arm log must show the `--label`ed diff header (`expected: LEG1_MODULES…` / `actual: .ail files swept…`), which is what makes the RED attributable to a human (§0.2(v)). |
| **AC7** | Unchanged, and now cheap: the copy set is file-scoped, so no arm ever opens a live directory for writing. |

**New AC8 — the landed criteria this change lands next to must not move.** Not in the doc; forced
by §0.2(iii).

### Hold set — re-measured at every task exit and at T5

| criterion | **base at `c53db58`** | required after |
|---|---|---|
| `verify_ail.sh` totals | rc=0, `4/4 … across 11 module(s)`, `14 required named tests`, `9/9 steps` (V1) | **4 / 11 / 14 UNMOVED** |
| `go test ./host/verifygate/` | **rc=0**, `ok … 15.714s`, **19** test names (V19) | rc=0, **23** names, < 30 s |
| `go vet ./host/verifygate/` | **rc=0** (V19) | rc=0 |
| `TestNoRigAbsolutePaths` | PASS | PASS — **the new file carries no rig-absolute literal** (V14) |
| `TestInScriptControl` anchor `v0.30.0|OK` | count **1** | **1** — prototype measured **1** (V20) |
| `TestEmptyExpectedReleaseSetFailsLoudly` fixture-read anchor | count **1** | **1** — prototype measured **1** (V20) |
| `TestCorpusPredicateMatchesShellAnchors` 2 regex literals | count **1** each | **1** each — the new block contains neither (V20) |
| `grep -c 'world/types.ail' scripts/verify_ail.sh` | **1** | **2** — expected and harmless; no committed test anchors on it (V20) |
| `host/boundary` `wantFileCount` | **1** (in `host/boundary` only; **0** in `host/verifygate`) | untouched (V21) |
| `GOTOOLCHAIN=go1.25.6 ./scripts/verify_go.sh` | rc=0 | rc=0 |

---

## 3. Task breakdown — 5 tasks, 5 commits, one milestone

Every task exits on the **full hold set** (§2) plus its own gate. Every Bash call starts
`export PATH=/opt/homebrew/bin:$PATH`; every `go` invocation carries `GOTOOLCHAIN=go1.25.6`; every
`verify_ail.sh` invocation carries `AILANG_BIN=/tmp/ailang-v0300/ailang` with z3 on PATH.

### T1 — `scripts/verify_ail.sh`: the identity allowlist and the Leg-1 restructure (+66/−4)

**This block is FROZEN. It was written, run, and mutated seven ways this session (V6–V16).** Type
it as given; deviations are re-litigating measured decisions.

**(a) The allowlist**, inserted immediately after the `ROOTS=( … )` block (they describe the same
sweep and must stay adjacent):

```bash
# Exact Leg-1 module manifest (identity allowlist, NOT a count). An intentional module
# add/remove is a ONE-LINE edit here, in the SAME commit. Repo-relative paths, matching
# the sweep's $mod key.
LEG1_MODULES=(
  design_docs/sketches/effectbroker.ail
  design_docs/sketches/logepoch.ail
  design_docs/sketches/storejournal.ail
  design_docs/sketches/transitions.ail
  design_docs/sketches/worlddapi.ail
  design_docs/sketches/worldkernel.ail
  design_docs/sketches/worldtypes.ail
  world/contracts.ail
  world/logepoch.ail
  world/transitions.ail
  world/types.ail
)
```

**(b) The NUL-aware, shell-quoting diagnostic formatter**, placed next to `run_bounded`:

```bash
# NUL-aware diagnostic formatter. Reads a NUL-delimited set file and prints ONE
# shell-quoted token per line, so `diff -u` stays line-oriented and a pathological path
# (newline, |, space, glob char) is RENDERED safely rather than PARSED.
_nul_quoted_lines() { # $1=NUL-delimited file
  local _p
  while IFS= read -r -d '' _p; do printf '%q\n' "$_p"; done < "$1"
}
```

**(c) Two more temp files and the extended trap**, beside the existing pair at `:156-158`:

```bash
tmp_mods_actual="$(mktemp -t verify_ail_mods_a.XXXXXX)"
tmp_mods_expected="$(mktemp -t verify_ail_mods_e.XXXXXX)"
trap 'rm -f "$tmp_json" "$tmp_test_json" "$tmp_mods_actual" "$tmp_mods_expected"' EXIT
```

**(d) Leg 1, restructured** — replaces the body from `total_verified=0` to the loop's
`done < <(find …)`:

```bash
total_verified=0
checked=0

# ── Leg 1a — enumerate ONCE into parallel indexed arrays (NUL end to end; no delimiter
# is ever embedded in a record, so no path can be mis-parsed by the gate that must reject it).
bases=(); rels=(); mods=()
for entry in "${ROOTS[@]}"; do
  base="${entry%%|*}"
  tree="${entry#*|}"
  if [ "$tree" = "." ]; then searchdir="$base"; else searchdir="$base/$tree"; fi
  [ -d "$searchdir" ] || continue
  while IFS= read -r -d '' f; do
    bases+=("$base")
    rels+=("${f#"$base"/}")      # module path relative to its base (what ai-check wants)
    mods+=("${f#./}")            # repo-relative path (manifest key), normalized (gemini catch)
  done < <(find "$searchdir" -name '*.ail' -print0 | sort -z)
done

# ── Leg 1b — MEMBERSHIP COMPARE, before any ai-check runs.
# Guards come BEFORE the printf writes and read ${#arr[@]}: under `set -u` (line 30) with
# bash 3.2 (the rig's /usr/bin/env bash), "${arr[@]}" on an EMPTY array is an unbound-variable
# ABORT, so a write-then-guard order would kill the script before its own null-case message
# could ever print. ${#arr[@]} is set -u-safe on an empty array.
if [ "${#mods[@]}" -eq 0 ]; then
  echo "✗ swept .ail enumeration was empty — the membership compare would pass vacuously; failing loudly" >&2
  exit 1
fi
if [ "${#LEG1_MODULES[@]}" -eq 0 ]; then
  echo "✗ LEG1_MODULES allowlist is empty — the membership compare would pass vacuously; failing loudly" >&2
  exit 1
fi
printf '%s\0' "${mods[@]}"         | LC_ALL=C sort -z > "$tmp_mods_actual"
printf '%s\0' "${LEG1_MODULES[@]}" | LC_ALL=C sort -z > "$tmp_mods_expected"
if ! cmp -s "$tmp_mods_expected" "$tmp_mods_actual"; then
  echo "✗ swept .ail module set differs from the LEG1_MODULES allowlist — an intentional" >&2
  echo "  module add/remove must edit LEG1_MODULES in scripts/verify_ail.sh in the SAME commit" >&2
  diff -u --label "expected: LEG1_MODULES in scripts/verify_ail.sh" \
          --label "actual:   .ail files swept under ROOTS" \
          <(_nul_quoted_lines "$tmp_mods_expected") <(_nul_quoted_lines "$tmp_mods_actual") >&2
  exit 1
fi
echo "   ✓ swept .ail module set equals the LEG1_MODULES allowlist (${#mods[@]} modules)"

# ── Leg 1c — consume the SAME arrays by index. cwd / run_bounded / absolute-temp semantics
# are untouched; "what is compared" and "what is checked" are the same array objects.
for i in "${!mods[@]}"; do
  base="${bases[$i]}"
  rel="${rels[$i]}"
  mod="${mods[$i]}"
  checked=$((checked + 1))
  echo "   ai-check $mod"
  … existing body, unchanged, re-indented one level out of the old `while` …
done
```

The old loop body moves verbatim into `for i in "${!mods[@]}"` — **only its indentation changes**.
`( cd "$base" && run_bounded … )`, the 124-TIMEOUT branch, the python manifest parse, and the
`case "$mod" in world/*)` accumulation are byte-identical.

**Exit gate**: `bash -n scripts/verify_ail.sh` rc=0; `AILANG_BIN=… ./scripts/verify_ail.sh` rc=0
with `4/11/14` unmoved **and** the new line `✓ swept .ail module set equals the LEG1_MODULES
allowlist (11 modules)` present. Hold set green (the four source-anchor tests in §2 especially).

### T2 — `host/verifygate/module_manifest_gate_test.go`: the isolated-root harness (~165 LOC)

Four helpers. **No rig-absolute path literal anywhere in the file** (§0.2(iii)).

- **`newIsolatedGateRoot(t) string`** — builds `<t.TempDir()>/iso` and copies, **file by file**
  (never `cp -R` of a directory, §0.2(iv)):
  1. `scripts/verify_ail.sh` (mode 0755) — placing it at `<iso>/scripts/verify_ail.sh` makes
     `<iso>` the repo root **by construction**, because the script does `cd "$(dirname "$0")/.."`
     (`:31`);
  2. `scripts/testdata/ailang_release_observed.txt` (the script FATALs on a missing/empty one,
     `:102-106`);
  3. every `world/*.ail` (4);
  4. every `design_docs/sketches/*.ail` (7).

  **Assert the copy landed: exactly 13 files under `<iso>`, of which exactly 11 are `.ail`.**
  Measured base 13/11 (V16). A copy helper that silently copied nothing would otherwise be caught
  only by the pristine control, one step too late.

- **`runGateAt(t, root string, env map[string]string) (int, string)`** — `runGate` (`:52-56`)
  hardwires `repoRoot` for **both** the script path and `cmd.Dir` (V4), so it cannot be reused.
  `runGateAt` is `runGate` with `root` substituted in both places; its env-hygiene blocklist
  (`AILANG_BIN`, `WORLD_PKG_AILANG_BIN`, `AILANG_SHIM_*`, `AILANG_Z3_PATH`) is reused as-is. Do
  **not** refactor the committed `runGate` — `TestReleaseChangeNotice` and friends depend on it.

- **`requirePristineControl(t, root)`** — runs the copied gate on the **unmutated** root and
  requires the Leg-1 success line `✓ 4/4 required world/ identities verified across 11 module(s)`.
  Also returns the output so an arm can take its `   ai-check ` count as the **control** for its own
  absence assertion (base 11, V12). Called by **every** arm **before** its mutation.

- **`mutateCopiedScript(t, root, old, new string)`** — reads `<iso>/scripts/verify_ail.sh`, asserts
  `strings.Count(src, old) == 1` (the `TestInScriptControl` discipline, `:405-413`), writes the
  mutant back into the **isolated** root. Unlike `TestInScriptControl`, it never `os.CreateTemp`s
  inside the live `scripts/` directory — this sprint is zero-live-write by design.

- **`requireLiveTreeUntouched(t)`** — asserts no `world/_stray*` exists and
  `design_docs/sketches/storejournal.ail`'s sha256 equals the value recorded at test start (AC7's
  in-test half). Called at the end of every arm.

**`requirePinned(t)`** is reused as-is (`:39-50`): a loud `t.Fatal` when `AILANG_BIN` is unset,
never a skip.

**Exit gate**: `go vet ./host/verifygate/` rc=0 — **`go build ./...` does not compile `_test.go` at
all**, so vet is the compile gate for a test-only file. Hold set green. (No assertion yet; T3
consumes these.)

### T3 — the four committed refusal arms (~180 LOC)

One arm per refusal branch (§4.2). Each: build an isolated root → **pristine control** → apply its
mutation → assert the RED → `requireLiveTreeUntouched`.

| arm | mutation | asserts |
|---|---|---|
| `TestModuleManifestRejectsStrayModule` | writes `<iso>/world/_stray_manifest_probe.ail`, fixed 3-line valid leaf | rc==1; output contains `LEG1_MODULES`, the probe path, and the `expected:`/`actual:` labels; output contains **zero** `\n   ai-check ` lines **while the pristine control had 11**; output does **not** contain `verify gate PASSED` |
| `TestModuleManifestRejectsDeletedModule` | deletes `<iso>/design_docs/sketches/storejournal.ail` | rc==1; the diff names `-design_docs/sketches/storejournal.ail`; live `storejournal.ail` sha256 unchanged |
| `TestModuleManifestEmptyAllowlistFailsLoudly` | `mutateCopiedScript` replaces the allowlist body with `LEG1_MODULES=()` (anchor `LEG1_MODULES=(` count==**1**, V20) | rc==1 **and** output contains `LEG1_MODULES allowlist is empty`; output does **not** contain `unbound variable` (the §0.2(i) regression guard) |
| `TestModuleManifestEmptyEnumerationFailsLoudly` | `mutateCopiedScript` neuters the actual-side enumeration (`mods+=("${f#./}")` → same + `; mods=()`, anchor count==1) | rc==1 **and** output contains `swept .ail enumeration was empty` |

**Why the deletion target must remain `storejournal.ail`, and isolation does not relax it.** The
doc's V7b measured that deleting `sketches/worldtypes.ail` reds via its **importer's** typecheck
(`sketches/transitions.ail: check.passed != true`) — so under `MUT-NEUTER-CMP` a non-leaf arm would
keep redding through an unrelated path and **mask the neutering**, surviving AC5 vacuously.
`storejournal` imports nothing and is imported by nothing; its red flows only through the compare.

**Why arms 3 and 4 are COMMITTED rather than demoted (§6.2 answers the doc's §5.4 question).**

**Exit gate**: `go test ./host/verifygate/ -run 'TestModuleManifest' -v` → 4 PASS, and **each arm's
log shows its pristine-control marker before its mutation red** — a run missing a control marker is
a FAIL even if the verdict is green. `go test ./host/verifygate/` → **23** names, rc=0, < 30 s.
Hold set green.

### T4 — the mutation sweep: 10 mutations / 20 arms (~0 LOC, the bulk of the time)

§4 in full. Transcript at `design_docs/verification/w-ail-gate-module-pin/mutations.md`.

### T5 — live one-shot ACs, AC7, and the doc repairs (the merge criterion)

- **T5.a** AC2/AC3/AC4 on the **live** tree, sequentially, with nothing else running (§5).
- **T5.b** AC5's live half: neuter the **live** script's `cmp` (`cp` backup + sha256), rerun
  `go test ./host/verifygate/ -run TestModuleManifest` → **the stray and delete arms FAIL** (the
  arms copy the live script, so the neutering propagates); restore byte-identical; rerun → PASS.
- **T5.c** AC7: `git status --porcelain` + `shasum -a 256 design_docs/sketches/storejournal.ail` +
  `ls world/_stray*` recorded **before** a full `GOTOOLCHAIN=go1.25.6 go test ./... -count=1`, and
  identical **after**.
- **T5.d** AC6: `AILANG_BIN=… ./scripts/verify_ail.sh` rc=0 and
  `GOTOOLCHAIN=go1.25.6 ./scripts/verify_go.sh` rc=0. **If `host/broker` reds, read *which* test
  failed** — the ~18% base flake is queue item 16, not this sprint's regression. Rerun once and
  attribute; never silence.
- **T5.e** Apply divergence repairs **D1–D6** (§2 of this plan / §0.2) to the design doc: the §4.2
  guard ordering, §6's missing B4 row, §8's missing `TestNoRigAbsolutePaths`, §5.2's copy-set
  scoping, the `--label` adoption, and the scope numbers.
- **T5.f** Final hold set re-measured in full.

---

## 4. Mutation discipline

### 4.1 Protocol — every arm, no exceptions

1. **Back up** the target to `/tmp/w12_backup/` with `cp`; record `sha256` **before**.
2. Apply the mutant. **Neuter with `if false && <cond>` or an anchored one-line replacement, never
   by deleting a block.** For `LEG1_MODULES`, **empty** it — never delete the declaration
   (§0.2(i)).
3. **Assert the mutant LANDED**: `sha256` after ≠ before. Record both. *(Shell has no compiler, so
   there is no "the mutant does not build" masking here — the landed-proof hash is the same
   discipline by another route. For the Go test file, `go vet ./host/verifygate/` **is** the build
   gate and must be rc=0 before any verdict is read.)*
4. **Kill arm**, `-run`-scoped to the **named test**, requiring rc≠0 **and** the recorded FAIL line
   and message — never the exit code alone.
5. **Inverse arm**, same mutant: `go test ./host/verifygate/ -skip 'TestModuleManifest'` → rc=0.
   This is what proves the new arms are the killer and not a bystander.
6. **Restore** with `cp` from `/tmp/w12_backup/`. **Never `git checkout --`.** Assert `sha256`
   equals the recorded before-value.
7. **Re-run the kill arm** and require rc=0. A restore that did not restore is otherwise invisible.

### 4.2 Rule 3j — the branch inventory, anchored to THIS SPRINT'S OWN DIFF

The milestone's headline verb is **refuse**, so the unit of mutation is the **BRANCH**, not the
milestone. This enumeration was produced by **reading the prototype diff** (V6), not by
transcribing §6 — §6 was written before a line of the code existed, and it is **missing B4** (V13).

| id | refusal branch (new in this diff) | message it writes | mutation |
|---|---|---|---|
| **B1** | actual enumeration empty — `[ "${#mods[@]}" -eq 0 ]` | `✗ swept .ail enumeration was empty …` | **M6** |
| **B2** | expected allowlist empty — `[ "${#LEG1_MODULES[@]}" -eq 0 ]` | `✗ LEG1_MODULES allowlist is empty …` | **M5** |
| **B3** | set inequality — `! cmp -s "$tmp_mods_expected" "$tmp_mods_actual"` | `✗ swept .ail module set differs …` + `exit 1` | **M4** |
| **B4** | the diagnostic INSIDE B3 — `diff -u --label … _nul_quoted_lines …` | the `--label`ed `+`/`-` lines naming the offending path | **M7** (missing from §6) |

**Pre-existing, deliberately untouched** (§10): the `:233` `checked -eq 0` guard, now shadowed by
B1, retained as defense-in-depth and reachable again under M4.

**Rule-3j cut instrument, with both controls in the same call.** Two traps bite this sprint:
`git diff` **omits untracked files** (the new `_test.go`), and a shell refusal is neither
`fmt.Errorf` nor `t.Fatalf`:

```bash
export PATH=/opt/homebrew/bin:$PATH
# shell side — a TRACKED, MODIFIED file, so ordinary git diff is correct here
git diff -- scripts/verify_ail.sh | grep -cE '^\+.*exit 1'
# Go side — the file is NEW, so git diff returns 0; --no-index is required
git diff --no-index /dev/null host/verifygate/module_manifest_gate_test.go \
  | grep -cE '^\+.*\bt\.(Fatal|Fatalf|Errorf)\('
# known-positive control, SAME instrument, a file known to have many        (base: 46)
git diff --no-index /dev/null host/boundary/allowlist_world_test.go \
  | grep -cE '^\+.*\bt\.(Fatal|Fatalf|Errorf)\('
# known-negative control, SAME instrument, same directory                    (base: 0)
git diff --no-index /dev/null host/broker/broker.go \
  | grep -cE '^\+.*\bt\.(Fatal|Fatalf|Errorf)\('
```

**The greps are a FLOOR, not the enumeration.** The enumeration is the table above, by reading the
diff.

### 4.3 The 10 mutations / 20 arms — and the observable each one reads (rule 3i)

**Every "kills which mutation" row names the observable AND the write that produces it.** All five
Group-A/B observables below were **measured on the prototype this session** (V11–V13, V16).

**Group A — data mutations. Prove the gate catches the real threat.**

| # | ID | Mutation | Branch | RED observable — and **which write produces it** |
|---|---|---|---|---|
| M1 | `MUT-STRAY` | valid 3-line leaf `world/_stray_manifest_probe.ail` (iso; live for AC2) | B3+B4 | `+world/_stray_manifest_probe.ail` in the diff, **written by the `diff -u` call inside the `! cmp -s` branch** — downstream of the compare, not set alongside it. Measured rc=1 in **0.046 s** with **0** `   ai-check ` lines (control: 11). |
| M2 | `MUT-DEL-LEAF` | delete `design_docs/sketches/storejournal.ail` (iso; live for AC3) | B3+B4 | `-design_docs/sketches/storejournal.ail`, **same write**. Measured rc=1, 0.045 s. |
| M3 | `MUT-SWAP` | M1 + M2 composed — the **count stays 11** | B3+B4 | **both** paths in one diff, **same write**. Measured rc=1. **This is the row that kills `EXACT_TOTAL_MODULES`**: today this exact tree passes with a success line byte-identical to the baseline's. |

**Group B — script-branch neutering. One per branch (rule 3j).**

| # | ID | Mutation | Branch | RED observable |
|---|---|---|---|---|
| M4 | `MUT-NEUTER-CMP` | `if ! cmp -s …` → `if false && ! cmp -s …` on the **LIVE** script (one-shot, sequential, sha'd backup) | B3 | the committed **stray and delete arms FAIL** (they copy the live script). **Measured on the prototype: with the mutant, the stray sails through at `✓ 4/4 … across 12 module(s)` and `   ai-check world/_stray_manifest_probe.ail` runs** — i.e. neutering restores today's defect exactly. **This is AC5, the proof this is a gate and not a guard.** |
| M5 | `MUT-EMPTY-ALLOWLIST` | `LEG1_MODULES=()` in the **iso** copy (anchor count==1) | B2 | `✗ LEG1_MODULES allowlist is empty …`, **written by B2's own guard**. Measured rc=1. Arm additionally requires the output **not** to contain `unbound variable` — the §0.2(i) ordering regression. |
| M6 | `MUT-ENUM-EMPTY` | neuter the actual-side enumeration in the **iso** copy | B1 | `✗ swept .ail enumeration was empty …`, **written by B1's own guard**. Measured rc=1. |
| M7 | `MUT-DIAG-SILENT` | replace the `diff -u --label …` line with `:` in the **iso** copy | B4 | the stray arm's **path assertion** reds while rc and the message stay identical. **Measured: rc=1, message present, `grep -c _stray_manifest_probe` = 0.** Not in §6. |

**Group C — self-satisfaction mutations. Prove the arms are not decorative.**

| # | ID | Mutation | RED observable |
|---|---|---|---|
| M8 | `MUT-ISO-INCOMPLETE` | omit `storejournal.ail` from the copy set | **the pristine control fails before any mutation is applied** — the `across 11 module(s)` marker is absent and the compare names the un-copied file. **Measured: rc=1, marker count 0.** Proves the control reads the compare rather than passing vacuously. |
| M9 | `MUT-ARM-RC-ONLY` | weaken the stray arm to assert only `rc == 1`, **composed with M7** | the arm must go **GREEN**. If it stays red, the weakening did not take; if the un-weakened arm was already green under M7, the path assertion is decorative and T3 is wrong. Inverse: M9 **alone** must be rc=0. |
| M10 | `MUT-ARM-CONTROL-DEAD` | delete `requirePristineControl` from an arm, **composed with M8** | the arm must go **GREEN on an incomplete copy** — proving the control, not the mutation, is what makes the copy's completeness observable. Inverse: M10 alone must be rc=0. |

**Totals: 10 mutations, 20 arms** (one kill + one inverse each; M9/M10's inverses are the
composed-vs-alone pairs).

**Why every observable is downstream of its mechanism.** B1's and B2's messages are written by
their own guards. B3's `exit 1` and B4's diff lines are written inside the `! cmp -s` branch. The
pristine-control marker (`✓ 4/4 … across 11 module(s)`, `verify_ail.sh:243`) sits **after** the
compare **and** after all 11 `ai-check` runs, so its presence proves copy completeness *and*
execution — nothing sets it alongside. The one absence assertion (no `   ai-check ` lines) is paired
with its own firing control in the same test (11 in the pristine run, 0 after the mutation).

---

## 5. Acceptance commands, as the executor must run them

All baselined on the **pristine tree at `c53db58`** this session; **the base result is part of the
criterion** (rule 3e). **No command contains `rg`.**

```bash
# AC1 / AC6 / AC9-equivalent — the AILANG gate. 4/11/14 must not move.
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN=/tmp/ailang-v0300/ailang
./scripts/verify_ail.sh > /tmp/w12_ac1.log 2>&1; echo "rc=$?"        # base rc=0
grep -c '4/4 required world/ identities verified across 11 module(s)' /tmp/w12_ac1.log   # 1
grep -c 'all 14 required named tests pass' /tmp/w12_ac1.log                              # 1
grep -c 'world package gate PASSED: 9/9 steps' /tmp/w12_ac1.log                          # 1
grep -c 'swept .ail module set equals the LEG1_MODULES allowlist' /tmp/w12_ac1.log       # 1 AFTER T1
```

```bash
# AC2 — stray reds (today: rc=0 at 12 module(s), V5 of the doc). LIVE tree, sequential.
export PATH=/opt/homebrew/bin:$PATH; export AILANG_BIN=/tmp/ailang-v0300/ailang
printf 'module world/_stray\n\nexport func strayId(x: int) -> int = x\n' > world/_stray.ail
./scripts/verify_ail.sh > /tmp/w12_ac2.log 2>&1; echo "rc=$?"          # want 1
grep -c 'LEG1_MODULES'                 /tmp/w12_ac2.log   # >0
grep -c 'world/_stray.ail'             /tmp/w12_ac2.log   # >0
grep -c '^   ai-check '                /tmp/w12_ac2.log   # 0  <- the SWEEP's line form, not the banner
grep -c 'verify gate PASSED'           /tmp/w12_ac2.log   # 0
rm world/_stray.ail
./scripts/verify_ail.sh > /dev/null 2>&1; echo "restored rc=$?"        # want 0
```

```bash
# AC3 — leaf deletion reds (today: rc=0 at 10 module(s)).
export PATH=/opt/homebrew/bin:$PATH; export AILANG_BIN=/tmp/ailang-v0300/ailang
mkdir -p /tmp/w12_backup
cp design_docs/sketches/storejournal.ail /tmp/w12_backup/
before=$(shasum -a 256 design_docs/sketches/storejournal.ail | cut -d' ' -f1)
rm design_docs/sketches/storejournal.ail
./scripts/verify_ail.sh > /tmp/w12_ac3.log 2>&1; echo "rc=$?"          # want 1
grep -c 'storejournal.ail' /tmp/w12_ac3.log                            # >0
cp /tmp/w12_backup/storejournal.ail design_docs/sketches/
after=$(shasum -a 256 design_docs/sketches/storejournal.ail | cut -d' ' -f1)
echo "restore ok=$([ "$before" = "$after" ] && echo yes || echo NO)"
./scripts/verify_ail.sh > /dev/null 2>&1; echo "restored rc=$?"        # want 0
```

```bash
# AC4 — the count-pin defeater. Today: rc=0 with a success line BYTE-IDENTICAL to the baseline's.
# Compose AC2's add with AC3's delete, then require BOTH paths in ONE diff.
```

```bash
# AC5 — committed arms, and the guard-vs-gate proof.
export PATH=/opt/homebrew/bin:$PATH; export AILANG_BIN=/tmp/ailang-v0300/ailang
GOTOOLCHAIN=go1.25.6 go test ./host/verifygate/ -run 'TestModuleManifest' -v -count=1
GOTOOLCHAIN=go1.25.6 go vet ./host/verifygate/          # go build ./... does NOT compile _test.go
# then MUT-NEUTER-CMP on the LIVE script (cp backup + sha256) -> the stray and delete arms FAIL
# restore byte-identical -> rerun -> PASS
```

```bash
# AC8 (new) — the landed criteria next door must not move.
export PATH=/opt/homebrew/bin:$PATH
GOTOOLCHAIN=go1.25.6 go test ./host/verifygate/ -run 'TestNoRigAbsolutePaths$|TestInScriptControl$|TestCorpusPredicateMatchesShellAnchors$|TestEmptyExpectedReleaseSetFailsLoudly$' -v -count=1
grep -c 'v0\.30\.0|OK' scripts/verify_ail.sh                                   # want 1
grep -cF "grep -vE '^[[:space:]]*(#|\$)' scripts/testdata/ailang_release_observed.txt" scripts/verify_ail.sh  # want 1
grep -rc 'wantFileCount' host/verifygate/ail_binary_gate_test.go               # want 0 (control: host/boundary -> 3)
```

```bash
# AC7 — live checkout byte-unchanged across the full test run.
export PATH=/opt/homebrew/bin:$PATH; export AILANG_BIN=/tmp/ailang-v0300/ailang
git status --porcelain > /tmp/w12_porc_before
shasum -a 256 design_docs/sketches/storejournal.ail > /tmp/w12_sha_before
ls world/_stray* 2>/dev/null | wc -l                     # want 0
GOTOOLCHAIN=go1.25.6 go test ./... -count=1 > /tmp/w12_full.log 2>&1; echo "rc=$?"
git status --porcelain > /tmp/w12_porc_after
shasum -a 256 design_docs/sketches/storejournal.ail > /tmp/w12_sha_after
diff /tmp/w12_porc_before /tmp/w12_porc_after && echo "porcelain identical"
diff /tmp/w12_sha_before  /tmp/w12_sha_after  && echo "sha identical"
ls world/_stray* 2>/dev/null | wc -l                     # want 0
```

---

## 6. Estimate — and why the doc's ~0.5 day is 1.5× low

| Task | shell LOC | test LOC | notes |
|---|---:|---:|---|
| T1 script: allowlist + formatter + restructure | 66 | 0 | prototyped end to end this session; **zero exploratory risk** |
| T2 isolated-root harness (4 helpers) | 0 | 165 | `newIsolatedGateRoot`, `runGateAt`, `requirePristineControl`, `mutateCopiedScript`, `requireLiveTreeUntouched` |
| T3 four committed refusal arms | 0 | 180 | one per branch B1–B3 + the empty-enumeration null case |
| T4 mutation sweep | 0 | 0 | **10 mutations / 20 arms** + transcript (~180 md lines) |
| T5 live one-shots, AC7, doc repairs D1–D6 | 0 | 0 | ~25 design-doc lines |
| **total** | **66** | **345** | **≈411 changed lines** |

**Reference velocity.** TR.C (iteration 75) priced 760 LOC + 46 arms at 1.25 days and the estimate
held. This sprint is roughly **half** on both axes (411 lines, 20 arms) → ~0.6 day on pure
proportion. Two things push it up and one pushes it down:

- **up** — the arms are a *new harness*, not new cases in an existing one: `runGateAt`,
  the copy helper, and `mutateCopiedScript` are all first-of-kind in this package;
- **up** — the doc's §4.2 ordering and §5.2 copy-set are both wrong as written (§0.2(i),(iv)), so
  the executor is implementing a *corrected* spec and must not "fix" it back;
- **down** — the arms are **cheap to run**: a refusing arm is **0.046 s** and a pristine control is
  **1.465 s** (V16), versus TR.C's ~35 s per `host/broker` arm. The whole `host/verifygate` package
  is 15.7 s at base (V19) and should stay under 30 s.

### Honest price: **≈0.75 day. 1.5× the doc's ~0.5 day.**

The doc anticipated this: §9 says *"if the sprint planner prices it above 0.5 day, the growth is
honest mechanical LOC, not hidden risk."* That is correct as far as it goes, but it is not the
whole reason — **the extra quarter-day is the two corrected mechanisms and the missing B4 branch,
not LOC.**

### 6.1 Verdict: **NO SPLIT.** Each candidate seam fails structurally.

| candidate seam | why it fails |
|---|---|
| **T1 (script) first, arms second** | Lands the pin with **no committed non-vacuity proof** — a guard, not a gate. That is *precisely* item 12's residual: iteration 71 landed a one-shot acceptance command about a tree that no longer exists and the item came straight back. Splitting here re-commits the exact error the item exists to close. **Rejected outright.** |
| **stray arm first, deletion + null-case arms second** | All four arms live in **one new file** sharing **one** helper set (`newIsolatedGateRoot`, `runGateAt`, `mutateCopiedScript`). Milestone 2 would edit milestone 1's file for zero independent value, and **AC7** (live checkout byte-unchanged across a full `go test ./...`) can only be measured once over the whole run. No seam. |
| **script + arms first, sweep second** | A sweep-less milestone ships an unproven gate, and **M4 `MUT-NEUTER-CMP` *is* the deliverable's proof** (AC5). Deferring it defers the only evidence that this is a gate. |
| **`world/` half first, `design_docs/sketches/` half second** | Not a seam at all — the allowlist is **one array** compared by **one `cmp -s`**. Half an allowlist is a *weaker gate*, not a smaller one: the doc's §3 shows the sketch half is where the measured deletion defect (V7) actually lives. |

**Run it whole at ≈0.75 day.** If the controller needs it inside half a day, the only honest lever
is the arm count — and this repo's recorded failure mode, three iterations running, is exactly
that cut. **Do not cut arms 3 and 4.**

### 6.2 Answer to the doc's §5.4 question: `MUT-EMPTY-ALLOWLIST` **SHIPS as a committed arm** — and so does `MUT-ENUM-EMPTY`

**Price:** the `mutateCopiedScript` helper (~25 LOC) is needed by **both** and is shared; arm 3
≈ 25 LOC, arm 4 ≈ 20 LOC. Marginal cost of taking both rather than neither: **~70 LOC, ≈25 minutes,
+0.1 s of CI.**

**Recommendation: commit both. Do not take the demotion path.** Three reasons, in order of force:

1. **A one-shot is a proof about a tree that no longer exists — which is item 12's ENTIRE residual
   complaint about iteration 71.** Demoting the S6 null case to a one-shot re-commits, at smaller
   scale, the very error this item was filed to correct. The demotion path is internally
   inconsistent with the item's own thesis.
2. **MEASURED: the branch these arms guard is exactly the one the doc got wrong** (§0.2(i)). A
   committed arm that requires the output to contain `LEG1_MODULES allowlist is empty` **and not**
   `unbound variable` catches any future reordering of the guards on **every CI run**. A one-shot,
   run once against an already-correct implementation, could never catch a later reordering — and a
   later reordering is not hypothetical, it is the single most likely edit to this block.
3. **S6 is a binding standard** ("Every gate must fail loudly on its own null case"), and this gate
   now has **two** null cases (empty allowlist, empty enumeration), not one. §6 of the doc already
   lists both; only one was slated for committing. Committing one and one-shotting the other would
   leave the pair asymmetric for no measured reason.

### 6.3 Answer on the restructure: the doc is RIGHT, and I have a stronger argument than the one it gives

**The enumerate → compare → consume restructure is correct. I do not disagree with the doc, so no
adjudication is needed — but the doc's own justification is weaker than what is measurable.**

The doc argues compare-first on three grounds: untrusted input should not be executed; refusals are
fast; the copy set stays small. Measured:

- **speed**: refusal in **0.046 s** vs a full legs-1–2 sweep at **1.465–1.638 s** — a **~32×**
  difference per refusing arm (V11, V16). Real, but not decisive on its own.
- **copy-set sufficiency**: I could **not** reproduce this as a defect. A compare-after sweep of 10
  or 12 modules runs fine inside the same 13-file copy. **This claim of §4.2 is unsupported**, and
  T5.e should soften it rather than repeat it.
- **the decisive one, which the doc does NOT make.** An **invalid** stray under compare-after reds
  **for the wrong reason**. Measured on one tree with both scripts (V12):

  | script | rc | what a human reads |
  |---|---|---|
  | today's (compare-after-nothing) | 1 | `✗ world/_bad_stray.ail: check.passed != true` — a **parse** failure. Sends the reader to fix syntax. |
  | prototype (compare-first) | 1 | `+world/_bad_stray.ail` in the **manifest** diff. Sends the reader to `LEG1_MODULES`. |

  Under compare-after, every committed arm's correctness would **depend on its probe being valid
  AILANG** — and if the probe ever stopped parsing, the arm would keep passing while measuring
  something else entirely. That is the doc's own **V7b wrong-reason class one level up**, and it is
  the strongest argument for compare-first available. Compare-first removes the dependency outright.

**Verdict: implement the restructure as designed.** The doc wins as the reviewed artifact; the
delta is that its §4.2 "the four-item copy stops being obviously sufficient" clause is unsupported
by measurement and should be replaced with the wrong-reason argument above.

---

## 7. Execution protocol and risks

### Protocol

- **Work in place** in `/Users/voightkampff/dev/sunholo-data/ailang-world` on `dev`. The executor
  has **no git write permission**; the controller reconstructs the 5 commits from cumulative
  `.snap/M<k>/` snapshots. **Empty the target files before reconstructing commits from `.snap/`** —
  the final-tree hash cannot detect absorption.
- **Every** Bash call starts `export PATH=/opt/homebrew/bin:$PATH` (else `go`/`gh`/`node` are
  rc=127 — an environment defect, never a spent quota).
- **Every** `verify_ail.sh` invocation exports `AILANG_BIN=/tmp/ailang-v0300/ailang` **and needs z3
  on PATH**. The PATH `ailang` is a `-dirty` dev build the gate **REFUSES** by design; a solverless
  run is a false green.
- **Every** `go` invocation carries `GOTOOLCHAIN=go1.25.6`. Local `go` is **1.26.4**, so
  `verify_go.sh` FATALs with *"active toolchain go1.26.4 miscompiles host/store/scan.go"* — a
  **base condition, not a regression** (V17). Do not chase it.
- **`go vet ./host/verifygate/` after every task and every arm.** `go build ./...` **does not
  compile `_test.go` at all**, so a mutation-BUILDS check via `go build` is **vacuous** on this
  test-only file. Vet is the compile gate.
- **zsh, not bash**: `${PIPESTATUS[0]}` expands EMPTY (zsh spells it `${pipestatus[1]}`) — capture
  with `cmd > /tmp/out 2>&1; echo "rc=$?"`. Quote glob-shaped flag values (`--include='*.go'`) or
  zsh aborts with `no matches found` and you read 0 from a command that never ran. **Brace any
  variable followed by a colon** (`"${rev}:host/x"`). zsh does not word-split unquoted variables —
  use arrays and assert `${#ARR[@]}`. `echo` is not byte-faithful — use `printf '%s'` or `od -c`.
- **The script's own shell is bash 3.2.57** (`/usr/bin/env bash` on this rig). No `mapfile`, no
  associative arrays, and `"${arr[@]}"` on an empty array under `set -u` is an **abort** (§0.2(i)).
- **`rg` is not a binary** here — a harness-injected shell function, absent under `env -i` and in
  CI. Never in a committed command or an acceptance criterion.
- **`git diff` omits untracked files.** The new `_test.go` needs
  `git diff --no-index /dev/null <file>`; `scripts/verify_ail.sh` is tracked so ordinary `git diff`
  is correct for it. Verify the instrument fires before banking its output.
- **Restores are `cp` from `/tmp/w12_backup/`**, never `git checkout -- <file>`.
- **Never touch** `~/.ailang/state/mission-v1*` or the V1 checkout.
- **SANDBOX CAVEAT.** A gate verdict obtained inside a `workspace-write` sandbox is
  **UNINFORMATIVE — neither a pass nor a fail**: loopback binds are denied there, which both
  invents failures (`host/daemon`, `host/broker`) and hides real ones. Report every socket-touching
  result as **"sandbox, uninformative"**; the controller re-runs them outside and that run is the
  verdict. `verify_ail.sh` and `./host/verifygate` do **not** bind sockets and are informative under
  the sandbox — but `go test ./...` (AC7) is not.

### Risks

| # | Risk | Assessment |
|---|---|---|
| **R1** | The doc prices this at ~0.5 day; I price it at **≈0.75 day** (§6). | **DECISION WANTED**, recommendation unambiguous: **run whole, no split** (§6.1). The doc itself flags 0.5 as "top of the band". |
| **R2** | **An executor implementing §4.2 *faithfully* ships a broken null case** (§0.2(i)). | Mitigated by freezing the block verbatim in §3/T1 **and** by arm 3's `not-contains "unbound variable"` assertion. **This is the single most likely way this sprint ships a false green.** |
| **R3** | The new file reds a **landed** criterion via `TestNoRigAbsolutePaths` (§0.2(iii)). | Mitigated: AC8 runs that test explicitly at every task exit. All paths from `repoRoot` + `t.TempDir()`; **no rig-absolute literal, not even in a comment**. |
| **R4** | Touching Leg 1's hardened loop breaks cwd / `run_bounded` / absolute-temp semantics. | Mitigated by AC1 (totals `4/11/14` unmoved, same tree, same binary) **and** by the prototype: the loop body moves verbatim, only its indentation changes (V6). |
| **R5** | The `host/broker` ~18% base flake reds AC6's `verify_go.sh`. | Not this sprint's regression (queue item 16). No task runs `./host/broker` directly. Read **which** test failed; rerun once; attribute; **never silence**. |
| **R6** | A mutant left behind by a failed restore breaches AC7. | Mitigated by §4.1 steps 6–7: sha256-asserted `cp` restore **and** a re-run of the kill arm requiring rc=0. AC7's porcelain/sha/glob triple is the independent backstop. |
| **R7** | Item 13 (`w-evidence-grade-mapping`) later adds a module *file* and must edit two literals in one commit. | Accepted and named: item 13 adds contracts/tests to the **existing** `world/types.ail`, which moves `EXACT_TOTAL_VERIFIED`/`EXACT_TOTAL_TESTS` but **not** `LEG1_MODULES`. If it ever adds a file, the edit is one array line and forgetting it fails with a message naming `LEG1_MODULES` and `scripts/verify_ail.sh` verbatim. |

---

## 8. Verification Log

Every row is a command actually run on **2026-08-12** at **`c53db58`**, clean tree, in the main
checkout (not a sandbox), with `AILANG_BIN=/tmp/ailang-v0300/ailang` (**AILANG v0.30.0**, commit
`e37b370`) and **z3 4.16.0** at `/opt/homebrew/bin/z3`. Empty/negative results carry a
**known-positive control scoped to the same path in the same call**. The live checkout was
**never written**: every prototype and mutation ran in `mktemp -d` roots, and `git status
--porcelain` was empty before and after (V0, V22).

| ID | Claim | Command | Observed |
|---|---|---|---|
| V0 | inspected revision, clean tree | `git log --oneline -1; git status --porcelain` | `c53db58 mission(world) iter 76: item 12 DESIGNED …`; porcelain **empty** |
| V1 | **controller baseline CONFIRMED** | `AILANG_BIN=… ./scripts/verify_ail.sh` | **rc=0**; `✓ 4/4 required world/ identities verified across 11 module(s)`; `✓ all 14 required named tests pass (failed_tests=0)`; `✓ world package gate PASSED: 9/9 steps performed non-zero work`; `✓ verify gate PASSED`; **2.825 s** wall |
| V2 | doc's line citations are exact | `grep -n 'checked=0\|checked + 1\|checked" -eq 0\|EXACT_TOTAL_VERIFIED=4\|module(s)' scripts/verify_ail.sh` | `:167` init, `:176` increment, `:233` `-eq 0`, `:238` `EXACT_TOTAL_VERIFIED=4`, `:243` print — **all five exact** |
| V3 | `host/verifygate/` has exactly one file | `ls -la host/verifygate/`; `wc -l` | `ail_binary_gate_test.go` only, **727** lines |
| V4 | **`runGate` is live-hardwired** → `runGateAt` required | read `ail_binary_gate_test.go:52-56` | `exec.Command(filepath.Join(repoRoot, "scripts", "verify_ail.sh"))` with `cmd.Dir = repoRoot` |
| V5 | `LEG1_MODULES` unallocated (negative existence + control **in the same call**) | `grep -rn LEG1_MODULES scripts/ host/ .github/`; then `grep -rn EXACT_TOTAL_VERIFIED scripts/ host/ .github/` | **0 hits, rc=1**; control **4 hits, rc=0** at `verify_ail.sh:238,239,240,243` — the instrument fires |
| V6 | the prototype patch applies and parses | python anchored patch → `/tmp/w77_verify_ail_proto.sh`; `bash -n`; `diff -u` vs original | `bash -n` **rc=0**; **+66 / −4** lines (final, with `--label`) |
| V7 | **the rig's `bash` is 3.2.57** | `command -v bash; bash --version \| head -1` | `/bin/bash`, `GNU bash, version 3.2.57(1)-release (arm64-apple-darwin25)` |
| V8 | **REFUTES doc §4.2 ordering** — `"${arr[@]}"` on an empty array under `set -u` ABORTS | two scripts of identical shape, `set -uo pipefail`, empty vs non-empty array | empty: `count=0`, then `A[@]: unbound variable`, **rc=1, next line never reached**. Control (non-empty): `count=2`, `REACHED_AFTER_PRINTF rc=0`, `od -c` → `x \0 y \0`. `${#A[@]}` is safe in **both** |
| V9 | **REFUTES §5.2's "four-item copy"** — `cp -R` drags gitignored caches | `cp -R world/ design_docs/sketches/` into an iso; `find -type f \| wc -l`; `git check-ignore -v` | **95** files (11 `.ail`); `.gitignore:3 **/.ailang/` matches `world/.ailang/cache` (rc=0); control `world/types.ail` **not** ignored (rc=1). Live cache holds entries for **nonexistent** modules `world___stray`, `world__transitionregistry` |
| V10 | V21 of the doc CONFIRMED — `$0`-relative rooting works; copy is sufficient | ran the **unmodified** copied script from `$HOME` | all 11 ai-check lines name the COPY's modules; `✓ 4/4 … across 11 module(s)`; `✓ all 14 … tests pass`; then **rc=127** at `:302` `./scripts/verify_world_package.sh: No such file or directory` — the deliberate Leg-3 stop; **1.630 s** |
| V11 | **the prototype's four branches all fire, each with its own message** | prototype in an iso; four mutations, each restored and sha-checked | pristine → `✓ swept .ail module set equals the LEG1_MODULES allowlist (11 modules)` + `across 11 module(s)`, **1.638 s**. `LEG1_MODULES=()` → `✗ LEG1_MODULES allowlist is empty …` rc=1. enumeration neutered → `✗ swept .ail enumeration was empty …` rc=1. stray → rc=1 in **0.046 s**. delete → rc=1 in **0.045 s** |
| V12 | **Q3 DECIDER + the banner trap** | one iso tree with an **invalid** stray, run under TODAY'S script then the prototype; grep the stray log two ways | TODAY: rc=1 `✗ world/_bad_stray.ail: check.passed != true` (**wrong reason**), 0.998 s. PROTOTYPE: rc=1 `+world/_bad_stray.ail` in the manifest diff, 0.061 s. Banner trap: bare `ai-check` → **1**; `^   ai-check ` → **0**; same needle on the pristine log (control) → **11** |
| V13 | **B4 exists and is missing from §6** — `MUT-DIAG-SILENT` | replace the `diff -u` line with `:` (anchor count==1), apply the stray | **rc=1**, message present, `grep -c '_stray_manifest_probe'` → **0**. The refusal survives; the path-naming dies |
| V14 | **REFUTES §8's completeness** — `TestNoRigAbsolutePaths` scans the new file | read the test; `grep -c` the doc for it **and two controls in the same call** | globs `host/verifygate/*.go`, `t.Errorf`s on `"/tmp"+"/ailang"`, `"/Us"+"ers/"`, `"/home"+"/runner/"`. Doc mentions: **0**; controls `TestInScriptControl` **4**, `TestCorpusPredicateMatchesShellAnchors` **2** |
| V15 | `diff --label` is available on this rig | `diff -u --label … --label … a b`; `diff --version` | labels honoured, rc=1; `Apple diff (based on FreeBSD diff)` |
| V16 | **the corrected shape, end to end** — `.ail`-scoped copy + `--label` | 13-file copy; stray applied; then pristine | copy = **13 files**; stray → rc=1 with `--- expected: LEG1_MODULES in scripts/verify_ail.sh` / `+++ actual:   .ail files swept under ROOTS` / `+world/_stray_manifest_probe.ail`; pristine → `across 11 module(s)`, **1.465 s** |
| V16b | **MUT-NEUTER-CMP restores today's defect exactly** (AC5's mechanism) | `if ! cmp -s …` → `if false && ! cmp -s …`, stray in place | `✓ 4/4 … across **12** module(s)`; `   ai-check world/_stray_manifest_probe.ail` count **1** — the gate sails through, so the committed arms must fail |
| V16c | **MUT-ISO-INCOMPLETE kills the control** | omit `storejournal.ail` from the copy set, run pristine | **rc=1**; `across 11 module(s)` marker count **0**; the compare names the un-copied file — the control is non-vacuous |
| V16d | **the NUL pipeline holds on both round-2 exploit shapes, and renders them legibly** | iso with `world/types.ail\|junk.ail` and a path containing a newline | **rc=1**; diff shows `+$'world/a\nb.ail'` and `+world/types.ail\|junk.ail` — refusal fires, both paths **named**. Honest count: **6** real files (NUL-counted); the NUL→NL line count reads **7**, inflated by the embedded newline — the instrument, not the tree |
| V17 | **`verify_go.sh` base condition** | `go version`; `./scripts/verify_go.sh` **without** `GOTOOLCHAIN` | `go1.26.4 darwin/arm64`; rc=1, `verify_go.sh: FATAL: active toolchain go1.26.4 miscompiles host/store/scan.go` — the script REFUSING an unpinned toolchain |
| V18 | CI lanes; no workflow edit needed | `grep -n 'AILANG_BIN\|verify_ail.sh' .github/workflows/ci.yml` | job 1 runs the script at `:96` exporting only `WORLD_PKG_AILANG_BIN` (`:93`); job 2 exports `AILANG_BIN` (`:144`). The membership compare is binary-independent |
| V19 | `host/verifygate` base cost and inventory | `go vet ./host/verifygate/`; `go test … -count=1`; `-list '.*'` | vet **rc=0**; test **rc=0**, `ok … 15.714s`; **19** test names |
| V20 | the new block breaks no source anchor | anchor counts in the ORIGINAL and the PROTOTYPE, same call | `v0.30.0|OK` **1 → 1**; fixture-read anchor **1 → 1**; both `TestCorpusPredicate…` regex literals absent from the new block; `world/types.ail` **1 → 2** (expected; no committed test anchors on it); `LEG1_MODULES=(` → **1** (a clean mutant anchor) |
| V21 | `host/verifygate` has no file-count pin (negative existence + control) | `grep -rn 'wantFileCount' host/`; `grep -rc` on the verifygate file | verifygate **0**; control `host/boundary/allowlist_world_test.go` **3 hits** incl. `wantFileCount = 1` — the instrument fires |
| V22 | the planning session left the live checkout byte-unchanged | `git status --porcelain`; `git log --oneline -1` | porcelain **empty**; HEAD still `c53db58` |

No new `.ail` source is proposed, so S5's pinned-binary source validation is not applicable; V1,
V10 and V11 nevertheless run the pinned v0.30.0 binary over every existing `.ail` module.

---

## 9. Handoff

- **Sprint plan**: `design_docs/planned/w-ail-gate-module-pin-sprint-plan.md` (this file)
- **Sprint JSON**: `.ailang/state/sprints/w-ail-gate-module-pin.plan.json`
- **Design doc**: `design_docs/planned/w-ail-gate-module-pin.md` (`d201a1e`)
- **Base**: `dev` @ `c53db58`
- Neither artifact is committed by the planner. **The controller commits.**

SPRINT_PLAN_PATH: design_docs/planned/w-ail-gate-module-pin-sprint-plan.md
SPRINT_JSON_PATH: .ailang/state/sprints/w-ail-gate-module-pin.plan.json
