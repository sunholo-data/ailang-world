# w-ail-gate-module-pin — pin the Leg-1 module set in `verify_ail.sh` itself

- **Status**: Planned
- **Item**: queue item 12, `w-ail-gate-module-pin`, clause-1
- **Estimated**: ~0.5 day (held; see §9)
- **Measurement base**: `304120b` (dev HEAD, clean tree), 2026-08-12
- **Instruments**: pinned `ailang` **v0.30.0** (commit `e37b370`) at `/tmp/ailang-v0300/ailang`;
  **z3 4.16.0** at `/opt/homebrew/bin/z3` (present for every run below — a solverless run is a
  false green and would invalidate every rc reported here)
- **Files touched**: `scripts/verify_ail.sh` (+66/−4 lines),
  `host/verifygate/module_manifest_gate_test.go` (NEW, ~270 LOC incl. the isolated-root
  helper — grown from the first draft's ~150–200; see §9)

Every codebase claim in this doc was re-measured first-party at `304120b` (Verification Log,
§11). Where a claim inherited from the queue row or iteration 71 was re-derived, the doc says
so; one inherited assumption was **refuted** by measurement (V7b) and the design changed because
of it. **This is the second draft**: a 2-reviewer quorum blocked the first on §5's committed
arms mutating the live checkout — §5.1 records the block and the repair, and rows V18–V24 are
the revision-session measurements (V18–V20 re-derive the controller's confirmation facts
first-party).

---

## 1. Problem

`scripts/verify_ail.sh` Leg 1 sweeps every `.ail` module under two roots
(`ROOTS=("design_docs|." ".|world")`, `verify_ail.sh:127-130`) and counts them into `$checked`.
That count is compared against **nothing except zero**: `$checked` occurs in exactly four
places — init (`:167`), increment (`:176`), the `-eq 0` vacuity guard (`:233`), and the
success-line print (`:243`) (V1). The known-positive control in the same file is `:238-239`:
`total_verified` **is** exactly pinned (`EXACT_TOTAL_VERIFIED=4`), so the script demonstrably
can pin a total and simply does not pin this one (V2). The charter's thrice-repeated
"totals stay 4/11/14" is therefore enforced for the 4, secondarily for the 14, and **the 11 is
decorative**.

Measured consequences at `304120b`, all reproduced live this session:

- **A stray module passes.** A valid `world/_stray.ail` (leaf module, no contracts, no tests)
  landed and the gate printed `✓ 4/4 required world/ identities verified across 12 module(s)`
  and `✓ verify gate PASSED`, rc=0 (V5).
- **A deleted leaf module passes.** Removing `design_docs/sketches/storejournal.ail` (imported
  by nothing, importing nothing) produced `… across 10 module(s)` and PASSED, rc=0 (V7).
- **Add-one-delete-one restores the count.** Both mutations composed produce
  `… across 11 module(s)` — the success line is **byte-identical to the pristine baseline's** —
  and the gate PASSES, rc=0 (V8).

Iteration 71 repaired only the *acceptance command* (a one-shot grep of the printed total on a
tree that no longer exists). The item's residual is exactly the guard-vs-gate distinction: a
guard is not a gate until something reds in CI when you remove it. This design moves the pin
into `verify_ail.sh` itself and commits the RED proof as a `host/verifygate` test.

## 2. The design question, settled: identity allowlist, not a count pin

The queue row says "an exact-total assertion mirroring `:239`". This doc deliberately does
**not** implement that shape, for three measured reasons:

1. **The count pin is defeated by a measured mutation.** V8 is the decisive row: with the wrong
   module set, the count is 11 and the success line is indistinguishable from the baseline. An
   `EXACT_TOTAL_MODULES=11` assertion passes V8's mutant. This is not hypothetical — both halves
   and their composition were run against the real gate this session.
2. **The script's own header forbids the aggregate form.** `verify_ail.sh:9-10`: the gate
   "asserts a HARDCODED manifest of proven-contract and passing-test IDENTITIES (not aggregate
   counts), with exact totals as SECONDARY checks only." A bare module count would be the one
   aggregate-primary check in a file whose header exists to warn against exactly that.
3. **The sibling leg already implements the identity form.** `verify_world_package.sh:86-96`
   (Leg 3 step 2/9) builds the expected `.ail` path list from a hardcoded array, `find`s the
   actual set, asserts **both** enumerations non-empty (so an empty compare cannot pass
   vacuously, `:93-94`), and `cmp -s`'s them with a `diff -u` on mismatch — observed live in the
   baseline run: `✓ find observed exactly 5 expected .ail files; no unexpected .ail file exists`
   (V9). `packages/world-core` is pinned by identity while Leg 1's 11 modules are pinned by
   nothing. This design ports that pattern, not a new one.

**What the count then is: redundant, and deliberately not added as a literal.** `cmp -s` over
two sorted full enumerations entails the count; a separate `EXACT_TOTAL_MODULES` would be a
third literal that must move in the same commit (the item-13 worry) while adding zero detection
power. The `$checked` print at `:243` stays as human-facing output — after this change it is
guarded by the set compare rather than decorative. Note the contrast with `EXACT_TOTAL_VERIFIED`
(V2): that count is a *secondary* defense over a set (`verify.results[]`) whose full membership
is **not** enumerated anywhere (only the required subset is), so the count adds coverage there.
Here the set **is** fully enumerated, so a count adds nothing.

## 3. Allowlist coverage: both swept roots (all 11 modules)

The pin covers everything Leg 1 sweeps: the 7 sketches
(`design_docs/sketches/{effectbroker,logepoch,storejournal,transitions,worlddapi,worldkernel,worldtypes}.ail`)
plus the 4 core modules (`world/{contracts,logepoch,transitions,types}.ail`) (V6).

The sketches' empty required-sets (`:121`) exist so a future *contracted* sketch cannot perturb
the **verified total** — that reasoning is about the contract axis and does not extend to
membership. On the membership axis the sketches are exactly as exposed as `world/`: V7 is a
*sketch* deletion passing silently, and a stray dropped in `design_docs/sketches/` traverses the
identical `find` in the identical loop as the `world/` stray in V5. Pinning `world/` only would
leave half the measured defect open.

### Out-of-scope `.ail` files (the F7 question, answered)

Six `.ail` files exist outside Leg 1's roots (V6). They are **deliberately not swept**, and this
is correct scoping rather than a second gap:

- `packages/world-core/*.ail` (5 files) — already pinned **by identity** by Leg 3 step 2/9
  (exact allowlist, V9) *and* by content: step 3/9 asserts SHA-256 equality of each package
  module against its canonical `world/` source. Adding them to Leg 1 would double-gate bytes
  already gated twice.
- `host/replay/testdata/transition_fixture.ail` — declares
  `module host/replay/testdata/transition_fixture` and is a **replay input**: the archived
  released binary executes it during replay, driven by `host/replay/replay.go` and exercised by
  `host/replay/replay_test.go` in the go-verify job (V16). It already has a gate — a broken
  fixture reds `go test ./...`. Sweeping it into Leg 1 would make the primary `.ail` gate depend
  on a Go-test fixture and add a third source-root semantics (`host/...` module paths) to a
  ROOTS mechanism that currently has two. If a stray-`.ail`-under-`host/` pin is ever wanted,
  it is a separate item; this one pins what Leg 1 sweeps.

## 4. Proposed change to `scripts/verify_ail.sh`

### 4.1 The allowlist

One bash array in the hardcoded-gate-policy section (near `ROOTS`, which it must stay adjacent
to — they describe the same sweep):

```bash
# Exact Leg-1 module manifest. An intentional module add/remove is a ONE-LINE edit here,
# in the SAME commit. Repo-relative paths, matching the sweep's $mod key.
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

`LEG1_MODULES` is unallocated today (grep: 0 hits repo-wide; control `ROOTS` = 3 hits in the
same file, V15). Gate policy stays hardcoded and non-env-overridable, per the `:12-14` contract.

### 4.2 Enumerate once, compare before executing anything

Leg 1 is restructured from *enumerate-and-check-in-one-loop* to *enumerate, compare, then
check* — a single enumeration feeds both the membership compare and the ai-check sweep:

1. Walk `ROOTS` with the existing `find … -name '*.ail' -print0` expression, reading each
   NUL-terminated result into **indexed bash arrays** (`bases[]`, `rels[]`, `mods[]`) — parallel
   arrays, never a delimiter-joined record. **No `base|rel|mod` triple file exists.**
   *(This replaces the first draft's line-delimited triple file — see §4.2a.)*
2. **Membership compare, before any `ai-check` runs**, kept **NUL-delimited end to end**. First
   assert `${#mods[@]}` and `${#LEG1_MODULES[@]}` are non-zero, then write the actual set with
   `printf '%s\0' "${mods[@]}"` and the expected set with `printf '%s\0' "${LEG1_MODULES[@]}"`
   and pass both through `LC_ALL=C sort -z`. The guards must precede the writes: bash 3.2 plus
   `set -u` aborts on `"${arr[@]}"` for an empty array, while `${#arr[@]}` is safe. These are
   the S6 null-case guards, mirroring
   `verify_world_package.sh:93-94` — without them a broken `find` plus an emptied array compare
   equal and pass vacuously); then `cmp -s`. Mismatch diagnostics are emitted through a
   **NUL-aware formatter** into `diff -u --label` for display only, with paths
   shell-quoted via `printf '%q'`), so a pathological path is *rendered* safely rather than
   *parsed*. The message on mismatch:

   ```
   ✗ swept .ail module set differs from the LEG1_MODULES allowlist — an intentional
     module add/remove must edit LEG1_MODULES in scripts/verify_ail.sh in the SAME commit
   ```

   Forgetting the allowlist edit therefore fails loudly, **naming the exact file and variable
   to edit**, with the offending path visible as a `+`/`-` line in the diff.

   **`LC_ALL=C` on both sorts, stated at its honest severity** (gemini-3-1-pro's suggestion,
   adopted): this fixes no live bug. Both sides of the compare are sorted in the same
   environment within one run, so a locale difference shifts both orderings identically and
   `cmp -s` still matches. What it buys is byte-determinism of the `diff -u` output across
   machines and locales, and immunity to a future edit that hardcodes one side's order (e.g. a
   pre-sorted allowlist literal compared without re-sorting). It costs one env prefix, so it is
   adopted — as hygiene, not as a fix.
3. The existing sweep loop consumes the **same arrays** by index
   (`for i in "${!mods[@]}"; do base="${bases[$i]}"; rel="${rels[$i]}"; mod="${mods[$i]}"; …`),
   preserving the `cd "$base"` / `run_bounded` / absolute-temp-path semantics untouched.
   `$checked`, the `:233` zero-guard, and the `:243` print are unchanged. Because the compare
   and the sweep read the *same array objects* rather than two parses of one serialisation,
   "what is compared" and "what is checked" cannot desync — which is strictly stronger than the
   first draft's shared temp file.

### 4.2a Why not a `base|rel|mod` triple file — the round-2 objection, applied verbatim

Round 2 passed `gemini-3-1-pro` and was rejected by `gpt5-6-sol` on the first draft's
line-delimited `base|rel|mod` temp file: a delimiter that can occur in the data is not a safe
encoding for a gate whose whole job is to reject *unexpected* paths, and "no current repo path
contains `|`" is a statement about today's paths — precisely the "the shape space of what it
refuses was never enumerated" failure this repo recorded as its most recent spine. Its fix is
adopted **verbatim** under the narrow-refinement carve-out (concrete reviewer-authored fix; the
design DIRECTION — identity allowlist, compare-first, isolated arms — was not disputed by either
reviewer): keep discovered paths NUL-delimited end to end via indexed arrays, `printf '%s\0'`,
`LC_ALL=C sort -z`, and a NUL-aware/shell-quoted formatter for diagnostics.

**The controller measured the objection before applying it, and reports the result honestly
(V25, V26): the two exploits gpt5-6-sol names do NOT reproduce — the *class* it names does.**
Both arms were run in isolated temp trees with a firing control (file-creation count asserted
== 2 before any verdict was read):

- **`|` in a filename** (`world/types.ail|junk.ail` alongside `world/types.ail`): column
  extraction yields an EXTRA row (`junk.ail`), so the compare **fails** and the pin **holds**.
  The claim that the truncated set could "satisfy the manifest comparison" is refuted.
- **newline in a filename**: 2 real files produce a **4-line** triple file and a garbled `mod`
  column (two empty rows). The compare again **fails** — fail-safe, not exploitable.

So the objection is **over-stated as an exploit and correct as a defect**. What is genuinely
established is that the line-delimited encoding *corrupts* on such paths: the sweep loop's own
`while IFS='|' read -r base rel mod` mis-parses the same rows, so the gate would run against
garbled paths and its `diff -u` diagnostic would print blank/partial lines — a real robustness
and diagnosability defect, and the first draft's "prevents drift" claim for the shared temp file
is false in exactly this case. The fix costs nothing (`find -print0` already emits NUL; keeping
it NUL removes the parse entirely) and is strictly better, so it is adopted on those grounds
rather than on the exploit's. Recording this distinction is the point: a reviewer's objection is
a claim too, and one that is right for a different reason than stated is worth landing *with*
the correction rather than laundered into agreement.

**Why compare-first rather than compare-after:** a stray `.ail` file is untrusted input. Today
the gate runs `ai-check` — which shells out to Z3 and parses the file — on anything the sweep
finds. An allowlist gate that executes the candidate before checking membership is backwards;
compare-first refuses the stray without ever handing it to the toolchain, fails in seconds
instead of after a full sweep, and makes the committed RED arms cheap (§5.4). The single shared
enumeration is what prevents drift between "what is compared" and "what is checked" — the
rejected alternative (a second, independent `find` before the untouched loop) reintroduces
exactly the two-instruments-disagreeing class this repo keeps paying for.

**Restructure risk, named:** this touches Leg 1's hardened loop (cwd discipline, `run_bounded`
out-file semantics). Mitigation is AC1: the full gate must be green with totals `4/11/14` and
the module-count print unchanged both before and after the edit, same tree, same binary.

## 5. Where the non-vacuity proof lives: committed arms in `host/verifygate`, each in a private isolated root

A one-shot acceptance command is a proof about a tree that no longer exists — item 12's residual
is precisely that distinction, so the RED mutations land as a **committed test file**,
`host/verifygate/module_manifest_gate_test.go` (does not exist today, V15; the package has no
file-count pin — `wantFileCount` greps to 0 hits there, positive control `host/boundary/allowlist_world_test.go:1163` = 1, V10). The file reuses `requirePinned` (loud `t.Fatal` when
`AILANG_BIN` is unset — never a skip) and the env-hygiene of `runGate` via a new
root-parameterized `runGateAt(t, root, env)` — the committed `runGate` itself is hardwired to
the live checkout (`exec.Command(repoRoot/scripts/verify_ail.sh)`, `cmd.Dir = repoRoot`, V24)
and cannot be reused as-is. Every arm runs the **real copied** `scripts/verify_ail.sh`, not a
reimplementation of its logic.

### 5.1 First draft BLOCKED: why the arms must not touch the live tree

The first draft's arms mutated the live checkout (probe written to the real `world/`, deletion
of the real sketch, `t.Cleanup` restore) and argued that skipping `t.Parallel()` made this safe.
Both quorum reviewers — gpt5-6-sol and gemini-3-1-pro, independently, across two providers —
blocked on the same defect: `go test ./...` runs **packages** concurrently, and "no
`t.Parallel()`" governs only ordering *within* `host/verifygate`. The premise is true by
construction: the Go gate runs `go test ./... -count=1` with no `-p 1` (`verify_go.sh:108`,
V18), so GOMAXPROCS-many package test binaries execute at once. And the collision is not
theoretical: `host/boundary/allowlist_world_test.go` **enumerates the live `world/` tree and
then reads every file it enumerated** on every run (`enumerateAIL` walks
`filepath.Join(root, "world")` at `:197-199`; `checkAILGroup` reads each enumerated path at
`:293`; any error becomes `t.Fatal` at `:775-777`; `root` is the live checkout via
`runtime.Caller`, `:80-87` — all V19). A stray `.ail` created and then removed by a
`host/verifygate` arm opens a window in which `host/boundary` enumerates the file and then
fails to read it — ENOENT → `t.Fatal` in a package that has nothing to do with this change.
Nondeterministic CI, exactly as both reviews said. (One review attributed the race to
`host/broker/invoke_boundary_test.go`; that file never touches `.ail` — zero `.ail` references,
walker skips non-`.go` at `:149` — V20. The real collider both enumerates *and* reads, which is
stricter than either review claimed.) The reviewers proposed the same fix, adopted here
verbatim: private isolated roots. Do not re-propose live-tree mutation from committed tests.

### 5.2 The isolated root, and why a file-scoped copy is sufficient

A helper (`newIsolatedGateRoot(t)`) builds a private root under `t.TempDir()` and copies exactly
13 files: `scripts/verify_ail.sh`, `scripts/testdata/ailang_release_observed.txt` (the script
FATALs on a missing/empty expected-release file, `:102-106`), four `world/*.ail` files, and seven
`design_docs/sketches/*.ail` files. Directory copies are forbidden because they drag gitignored
`.ailang/cache` entries. The script does `cd "$(dirname "$0")/.."` (`:31`), so placing the copy
at `<iso>/scripts/verify_ail.sh` makes `<iso>` the repo root **by construction** — measured:
run from `$HOME`, the copy checked exactly its own 11 modules and never looked at the live tree
(V21). `cmd.Dir` is set to `<iso>` anyway (the reviewers' formulation; belt to the `$0`
resolution's braces). The live checkout is never written by any arm.

The copy can be this small because of compare-first (§4.2): a **refusing** arm stops at the
membership compare — it never reaches any `ai-check`, never reaches Leg 2, never reaches
Leg 3. The **pristine control** (§5.3) runs further: all of Leg 1 and Leg 2 against the copy
(measured ~1.7 s wall total, V21 — the modules are small), then stops at Leg 3 with the
deliberate, asserted rc=127 `./scripts/verify_world_package.sh: No such file or directory`
(`:302`) — Leg 3's own script and `packages/` tree are intentionally not copied, and no arm
needs them. Per arm: the pristine control reaches Legs 1–2 and the known Leg-3 stop; the
mutation run reaches only the compare.

### 5.3 The pristine-copy control — the isolation's own non-vacuity row

An isolated root that is missing a file reds for the WRONG reason, and an arm that reds for the
wrong reason proves nothing while looking exactly like a kill — the same error class V7b caught
in the first draft, one level up. So **every arm, in the same test, before applying its
mutation**, runs the copied script on the unmutated root and requires the Leg-1 success line
`✓ 4/4 required world/ identities verified across 11 module(s)` — a marker only reachable past
the membership compare, so its presence proves both that the copy is complete (all 11 modules
present, or the compare would have redded) and that the script actually executes against the
copied root. Only then is the mutation applied and the red asserted. Without this row the
committed-arm story is a vacuous pass wearing a green. Its own non-vacuity is MUT-ISO-INCOMPLETE
(§6): drop one file from the copy set and the control itself must fail.

### 5.4 The three arms

- **`TestModuleManifestRejectsStrayModule`** — builds an isolated root; pristine control
  (§5.3); writes a valid leaf probe `<iso>/world/_stray_manifest_probe.ail` (fixed 3-line
  content); runs the copied gate; requires rc=1, output containing `LEG1_MODULES` **and** the
  probe path, and output **not** containing `ai-check world/_stray_manifest_probe.ail` (pins
  the compare-first property) nor `verify gate PASSED`. Because compare-first refuses before
  any ai-check, the mutation run costs seconds. The defect this arm guards was reproduced
  inside an isolated copy this session — a stray written into the copy's `world/` sailed
  through today's gate at `12 module(s)` (V22) — so isolation measures the same thing V5
  measured live. Finally the arm asserts the live `world/` contains no `_stray_manifest_probe`
  file (the no-live-write assertion, R4's in-test half).
- **`TestModuleManifestRejectsDeletedModule`** — builds an isolated root; pristine control;
  deletes `<iso>/design_docs/sketches/storejournal.ail`; runs the copied gate; requires rc=1
  with the missing path named in the diff output; asserts the **live** `storejournal.ail`'s
  sha256 is unchanged from the value recorded at test start. No restore machinery exists —
  `t.TempDir()` cleanup replaces the first draft's in-memory-backup protocol entirely. **The
  deletion target must still be a leaf**, and isolation does not relax this: V7b measured that
  deleting `sketches/worldtypes.ail` reds via its importer's typecheck
  (`sketches/transitions.ail: check.passed != true`), so under MUT-NEUTER (compare removed) a
  non-leaf arm would *keep redding* through that unrelated path and mask the neutering — the
  arm would survive AC5 vacuously. `storejournal` (imports nothing, imported by nothing, zero
  required identities, V7b) is the one target whose red flows only through the compare.
- **`TestModuleManifestEmptyAllowlistFailsLoudly`** — the S6 null-case arm. Builds an isolated
  root; pristine control with the **unmutated** script; then overwrites
  `<iso>/scripts/verify_ail.sh` with a mutant whose allowlist enumeration is neutered
  (anchor-count asserted ==1 before replacing, the `TestInScriptControl` discipline,
  `ail_binary_gate_test.go:405`); runs it; requires an early loud `empty` failure rather than a
  pass or a Leg-1 red. Note one deliberate divergence from `TestInScriptControl`: that test
  `os.CreateTemp`s its mutant *inside the live `scripts/` directory* (harmless to the
  `host/boundary` walk — a dot-prefixed non-`.ail` file — but still a live write); this item's
  arms are zero-live-write by design, so the mutant lives in the isolated root. If the sprint
  runs tight, this is the one arm demotable to a one-shot mutation-table entry — the first two
  are mandatory, because the guard-vs-gate proof (MUT-NEUTER) lives in the stray arm.

The arms still do not call `t.Parallel()` — but after this revision that is defense-in-depth
(they share only the pinned binary and read-only live sources), not the load-bearing
protection. The load-bearing protection is the private root; the first draft had these two
confused, which is what the quorum caught.

Attribution note: each arm asserts the failure message **names its own mutation's path**, which
is what makes the red attributable, and the pristine control (§5.3) is the green half measured
in the same test against the same copied root. The live tree's green is separately asserted
every CI run by job 1 running the script itself and by the package's existing accept arms.

## 6. Mutation table

Landed-proof rule for every arm: record sha256 of the mutated file before and after, require
the hash to have **moved** before reading any verdict (a mutation that never landed reads
identically to one that did not red). For the committed arms the landed-proof hash is taken on
the **isolated root's** file — and its dual is the live-tree half: the LIVE counterpart's hash
must NOT move across the arm (§5.4's no-live-write assertions, AC7). One-shot mutations restore
from a `cp`/in-memory backup, never `git checkout --`. Shell has no compiler, so no "mutant
does not build" masking exists here — but the landed-proof requirement is the same discipline
by another route.

| # | Mutation | Measured today (304120b) | Required post-change | Observable that reds | Killed by |
|---|---|---|---|---|---|
| MUT-STRAY | add valid leaf `world/_stray….ail` **in the arm's isolated root** | **rc=0**, `12 module(s)`, PASSED — live (V5) and inside an isolated copy (V22) | rc=1 before any ai-check of the stray | `diff -u` `+world/_stray…` + message naming `LEG1_MODULES` | committed `TestModuleManifestRejectsStrayModule`; AC2 |
| MUT-DEL-LEAF | delete `design_docs/sketches/storejournal.ail` **in the isolated root** | **rc=0**, `10 module(s)`, PASSED (V7) | rc=1 | `diff -u` `-…storejournal.ail` | committed `TestModuleManifestRejectsDeletedModule`; AC3 |
| MUT-SWAP | both of the above composed | **rc=0**, `11 module(s)` — success line byte-identical to baseline (V8) | rc=1 | both paths in the diff | AC4 (one-shot; jointly covered by the two arms above) |
| MUT-NEUTER | neuter the **live** script's `cmp` (sha'd working-copy edit, restored; one-shot, sequential — never concurrent with `go test`). The arms COPY the live script into their isolated roots, so the neutering propagates into every arm | n/a (assertion doesn't exist yet) | committed stray arm **FAILS** | `go test -run TestModuleManifestRejectsStray` reds | AC5 — the proof this is a gate, not a guard |
| MUT-EMPTY-ALLOWLIST | empty `LEG1_MODULES` in a mutant script written into the isolated root | n/a | loud `✗ … empty` rc=1 (S6 null case), not a vacuous pass | the non-empty guard's own message | committed `TestModuleManifestEmptyAllowlistFailsLoudly` (or one-shot fallback) |
| MUT-ENUM-EMPTY | neuter the actual-side enumeration (mutant script in an isolated root) | n/a | loud `✗ swept .ail enumeration was empty` rc=1 | the actual-side non-empty guard | one-shot in sprint |
| MUT-DIAG-SILENT | neuter only the `diff -u --label` diagnostic inside the inequality branch | n/a | refusal remains rc=1 but the stray arm fails its offending-path assertion | the diagnostic's `+`/`-` path line | committed stray arm; mutation M7 |
| MUT-ISO-INCOMPLETE | omit one file (e.g. `storejournal.ail`) from an arm's copy set | n/a | the arm's **pristine control** (§5.3) fails — compare reds on the missing path *before any mutation is applied* | the required Leg-1 success line is absent; the compare diff names the un-copied file | one-shot in sprint — proves the control reads the compare rather than passing vacuously |

Each proposed assertion (set-equality, expected-non-empty, actual-non-empty) has at least one
mutation above that reds **through the observable that assertion reads** — the diff/message text
is produced *by* the compare mechanism, not set alongside it.

## 7. Acceptance criteria

All commands run from repo root with `export AILANG_BIN=/tmp/ailang-v0300/ailang`, z3 on PATH,
and are baselined against the pristine-tree measurements in the Verification Log (a criterion
about a behaviour the gate has been measured not to look at is listed with its measured-today
rc so both arms are visible). AC2–AC4 are **one-shot operator commands on the live tree, run
sequentially with nothing else executing** — the concurrency hazard that blocked the first
draft (§5.1) applies to committed tests inside `go test ./...`, which is why the committed arms
(AC5) run in isolated roots and AC7 exists. The one-shots stay on the live tree deliberately:
they prove the real installed gate reds, not only a copy of it.

- **AC1 — baseline preserved.** Pristine tree post-change:
  `./scripts/verify_ail.sh; echo rc=$?` → `rc=0`, output contains
  `✓ 4/4 required world/ identities verified across 11 module(s)` and
  `✓ verify gate PASSED: 4 required identities verified, 14 named tests pass` — the exact V4
  lines; totals `4/11/14` unmoved. (Fails if the loop restructure broke Leg 1.)
- **AC2 — stray reds (today: rc=0, V5).** `printf 'module world/_stray\n\nexport func strayId(x: int) -> int = x\n' > world/_stray.ail`;
  gate → `rc=1`, output contains `LEG1_MODULES` and `world/_stray.ail`, contains **neither**
  `ai-check world/_stray.ail` **nor** `verify gate PASSED`; `rm world/_stray.ail`; gate →
  `rc=0` again.
- **AC3 — leaf deletion reds (today: rc=0, V7).** `cp` backup, `rm design_docs/sketches/storejournal.ail`;
  gate → `rc=1` with the path named; restore from backup (sha256 equal to the backup's); gate →
  `rc=0`.
- **AC4 — the count-pin defeater reds (today: rc=0 with a baseline-identical success line, V8).**
  AC2's add + AC3's delete composed; gate → `rc=1` with **both** paths in the diff output.
- **AC5 — committed arms pass in isolation, and red when the mechanism is removed.**
  `GOTOOLCHAIN=go1.25.6 go test ./host/verifygate/ -run 'TestModuleManifest' -v` → all arms
  PASS, and each arm's log shows its pristine-control marker (the Leg-1 success line from its
  own isolated root, §5.3) **before** its mutation red — a run missing the control marker is a
  fail even if the verdict is green. Then neuter the **live** script's `cmp` (sha'd backup),
  same command → the stray arm **FAILS**, because the arms copy the live script and the
  neutering propagates into their isolated roots; restore byte-identical. Also
  `go vet ./host/verifygate/` → rc=0 — **`go build ./...` does not compile `_test.go` at all**,
  so vet is the compile gate for a test-only file.
- **AC6 — full gates green.** `./scripts/verify_ail.sh` rc=0 and
  `GOTOOLCHAIN=go1.25.6 ./scripts/verify_go.sh` rc=0 (the latter runs `go test ./...`, which
  includes the new arms; it FATALs without `GOTOOLCHAIN=go1.25.6` — base condition, not a
  regression).
- **AC7 — live checkout byte-unchanged across the full test run** (gpt5-6-sol's criterion —
  the row that would have caught the first draft's defect). Before a full
  `GOTOOLCHAIN=go1.25.6 go test ./... -count=1` (may be taken around AC6's `verify_go.sh`
  run): record `git status --porcelain` output, `shasum -a 256
  design_docs/sketches/storejournal.ail`, and assert no `world/_stray*` exists. After: the
  porcelain output is identical, the sha is identical, and `world/_stray*` still does not
  exist. The committed arms additionally assert their own halves of this in-test (§5.4), so a
  violation is caught per-arm, not only end-to-end.

## 8. Conflict surface

- **`scripts/verify_ail.sh` source is read by committed harness tests.** `TestInScriptControl`
  requires anchor `v0.30.0|OK` count==1; `TestCorpusPredicateMatchesShellAnchors` requires the
  two version-regex literals count==1; `TestEmptyExpectedReleaseSetFailsLoudly` reads the source
  too (V11). The new block must not duplicate any of those literals — it contains none of them.
  `world/types.ail` appears once in the script today (the python `REQUIRED_VERIFIED` dict) and
  will appear twice after the array lands; no committed test anchors on that string (V15).
- **`host/verifygate` accepts a new file.** No file-count pin exists there (V10; the
  `host/boundary` `wantFileCount = 1` landmine is scoped to `host/boundary` only — the fix
  there is always a new package, and is not needed here). `TestNoRigAbsolutePaths` does glob
  every `host/verifygate/*.go` file, including the new one, so all paths in it are assembled from
  `repoRoot` and `t.TempDir()`; no rig-absolute literal may appear even in a comment.
- **The repo-wide AST gate** (`host/broker/invoke_boundary_test.go`) floors the walked
  production-file count at ≥30 and pins no exact total (`:209-213`, V13), and the walker skips
  `_test.go` files anyway; the new test file must contain no `Invoke`/`NewSession`/
  `NewReplaySession` selectors or `broker.Session` mentions — it needs none. It never reads
  `.ail` files (V20), so it is NOT a live-tree collision surface.
- **`host/boundary` walks and reads the LIVE `world/` tree on every `go test ./...`**
  (`allowlist_world_test.go:197-199` enumerate, `:293` read, error → `t.Fatal` `:775-777`,
  V19) — this is the collision that blocked the first draft (§5.1). The committed arms never
  write the live tree, so no interaction remains. The walk is scoped to `world/` only (zero
  `design_docs` references in that file, V19), which means the deletion arm alone would not
  have collided *via this walk* — but both arms are isolated anyway: a design that argues one
  live mutation is safe is one refactor of `enumerateAIL`'s root away from re-opening the race.
- **CI** (`.github/workflows/ci.yml`): job 1 runs `./scripts/verify_ail.sh` (`:96`) with **no
  `AILANG_BIN` exported** (it exports only `WORLD_PKG_AILANG_BIN`), resolving `releases/latest`
  off PATH — the membership compare is binary-independent, so no lane divergence is introduced.
  Job 2 (go-verify) exports the pinned binary and installs z3, so the new arms run there with
  `requirePinned` satisfied (V12). No workflow edit is needed.
- **Queue item 13 (`w-evidence-grade-mapping`)** — the ergonomics this design was told to
  protect. Item 13 adds contracts and tests to the **existing** `world/types.ail`; that moves
  `EXACT_TOTAL_VERIFIED` and `EXACT_TOTAL_TESTS` but **not** `LEG1_MODULES`, because the module
  *set* is orthogonal to per-module contract/test content. The feared "third literal moving in
  the same commit" only materializes if item 13 adds a new module *file* — in which case the
  edit is one array line and forgetting it fails with a message naming `LEG1_MODULES` and
  `scripts/verify_ail.sh` verbatim.
- **Queue item 16 (`w-broker-base-flake`)** — no overlap: the new arms stop at the membership
  compare (no leg-3, no `host/broker` involvement), and the file set is disjoint. The ~18%
  `host/broker` flake can still red an AC6 `verify_go.sh` run for unrelated reasons; if it does,
  rerun and attribute — do not chase it inside this item.
- **No committed instrument asserts the module count today** (grep for
  `modules11|EXACT_TOTAL_MODULES|module(s)` across `scripts/ host/ .github/` hits only the
  `:243` print, which is the in-scope control firing, V14) — so nothing else needs updating
  when the pin lands, and iteration 71's one-shot AC command is superseded rather than
  contradicted.

## 9. Scope re-checked after the isolation rework: held at ~0.5 day, at the top of the band

Script: one array + an enumerate/compare/consume restructure of an existing loop, +66/−4 lines,
patterned line-for-line on `verify_world_package.sh:86-96` — unchanged from the first draft.
Test file: the isolation rework is real LOC — the `newIsolatedGateRoot` copy helper, the
`runGateAt` variant (`runGate` is live-hardwired, V24), per-arm pristine controls, and the
no-live-write assertions take the file to **~270 LOC**. The estimate is
re-checked against that growth rather than silently held: it stays at ~0.5 day, now at the top
of the band, for two measured reasons. First, the mechanics were de-risked live this session —
the 13-file `.ail`-scoped copy set is sufficient, the script self-roots from `$0`, and a pristine legs-1–2
run against the copy costs ~1.7 s wall (V21/V22) — so neither construction nor CI runtime is
exploratory. Second, the helper *replaces* first-draft machinery it deletes: the
in-memory-backup/`t.Cleanup` careful-restore protocol is gone entirely (`t.TempDir()` cleanup
subsumes it). If the sprint planner prices it above 0.5 day, the growth is honest mechanical
LOC, not hidden risk. No pins move (the set is today's 11), no CI edit, no docs beyond this
one. The one place scope could still shrink — committing MUT-EMPTY-ALLOWLIST as a third arm —
has a named demotion path (§5.4). The first draft's other fallback, compare-after-sweep (keep
the loop untouched, compare at the end), is now priced HIGHER than before and should not be
taken lightly: without compare-first the refusing arms and every pristine control pay a full
ai-check sweep, and the file-scoped copy stops being obviously sufficient; that trade is recorded
here so it is a decision, not a drift.

## 10. What this item is NOT doing

- **Not** adding an `EXACT_TOTAL_MODULES` count literal (§2 — defeated by V8, redundant under
  set equality, and the item-13 third-literal tax).
- **Not** sweeping `packages/world-core/*.ail` or `host/replay/testdata/transition_fixture.ail`
  into Leg 1 (§3 — the former is double-pinned by Leg 3 steps 2–3; the latter is a replay input
  gated by `host/replay`'s Go tests).
- **Not** touching `REQUIRED_VERIFIED`, `EXACT_TOTAL_VERIFIED`, `EXACT_TOTAL_TESTS`, the `:233`
  zero-guard, Leg 2, Leg 3, or `verify_world_package.sh`.
- **Not** adding env-overridable knobs (gate policy stays hardcoded per `:12-14`).
- **Not** editing `.github/workflows/ci.yml`, `host/boundary`, or anything in the frozen core.
- **Not** fixing the `host/broker` base flake (item 16 owns it).
- **Not** mutating the live checkout from any committed test — the first draft did (with
  restore) and was quorum-BLOCKED for it (§5.1). Only the sequential one-shot ACs (AC2–AC4)
  touch the live tree, with restore and byte-equality checks, never concurrently with
  `go test ./...`.

## 11. Verification Log

All rows measured first-party, 2026-08-12, at `304120b` (clean tree before and after — every
mutation below was restored and `git status --porcelain` re-verified; in the revision session
the porcelain baseline is the single untracked entry for this doc itself), with
`AILANG_BIN=/tmp/ailang-v0300/ailang` (v0.30.0, commit `e37b370`) and z3 4.16.0 at
`/opt/homebrew/bin/z3`. Rows re-derive the queue row's and the brief's claims rather than
inheriting them; V5's probe content (and hence sha) differs from iteration 71's — it is an
independent reproduction of the same defect, not a replay of the recorded one. **Rows V18–V24
are the revision-session measurements taken after the quorum block; V18–V20 are first-party
re-derivations of the controller's C1–C3 confirmation facts and are marked as such.**

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | `$checked` compared only against zero | Read `scripts/verify_ail.sh` | `:167` init, `:176` increment, `:233` `-eq 0`, `:243` print — no other comparison |
| V2 | script CAN pin totals (positive control) | Read `scripts/verify_ail.sh` | `:238` `EXACT_TOTAL_VERIFIED=4`, `:239` `-ne` compare; `EXACT_TOTAL_TESTS = 14` in Leg-2 python (compare at `:294`) |
| V3 | script byte-unchanged since the item was filed | `git diff --stat adfaa0b..HEAD -- scripts/verify_ail.sh` | empty; control same range over `design_docs/world-mission.md` → `1 file changed, 202 insertions(+), 9 deletions(-)` — instrument fires |
| V4 | baseline green + non-vacuous | `./scripts/verify_ail.sh > log; echo rc=$?` | `rc=0`; `✓ 4/4 required world/ identities verified across 11 module(s)`; `✓ verify gate PASSED: 4 required identities verified, 14 named tests pass` |
| V5 | stray module passes today | write `world/_stray.ail` (3 lines, sha256 `3c5bdecc…b064`; `world/*.ail` 4→5; standalone `ai-check`: `check.passed=True verify.errors=0 results=0`); run gate; `rm`; tree clean | `rc=0`; `… across 12 module(s)`; `verify gate PASSED` |
| V6 | module census | `find . -name '*.ail' -not -path './.git/*' \| sort` | 17 files: 7 sketches + 4 `world/` (the 11 swept) + 5 `packages/world-core/` + 1 `host/replay/testdata/` |
| V7 | leaf deletion passes today | `cp` backup; `rm design_docs/sketches/storejournal.ail`; run gate; restore; tree clean | `rc=0`; `… across 10 module(s)`; `verify gate PASSED` |
| V7b | **refutation** — non-leaf deletion reds for the WRONG reason | same protocol on `sketches/worldtypes.ail`; `grep -Hn '^import'` over all swept modules | `rc=1` at `✗ design_docs/sketches/transitions.ail: check.passed != true` (its importer, not any set pin). Import graph: `sketches/worldtypes` ← `sketches/transitions`; `sketches/logepoch` ← 3 sketches; `world/{logepoch,types,contracts}` ← `world/` siblings; **storejournal imports nothing, imported by nothing** — hence the arm choice in §5.4 |
| V8 | add-one-delete-one defeats a count pin today | V5's add + V7's delete composed; run gate; restore both; tree clean | `rc=0`; `✓ 4/4 required world/ identities verified across 11 module(s)` — byte-identical to V4's line; `verify gate PASSED` |
| V9 | sibling leg implements the identity pattern | Read `scripts/verify_world_package.sh:86-96`; V4 run output | expected list from array, `find` actual, both-non-empty guards `:93-94`, `cmp -s` + `diff -u`; live: `✓ find observed exactly 5 expected .ail files; no unexpected .ail file exists` |
| V10 | `host/verifygate` has no file-count pin (negative-existence) | `grep -rn 'wantFileCount' host/` | hits only `host/boundary/allowlist_world_test.go:1163-1165` (positive control); `ls host/verifygate/` → 1 file, 727 lines |
| V11 | harness tests anchor on script source | Read `ail_binary_gate_test.go:405-451,:502,:712-727` | `TestInScriptControl` requires `v0.30.0\|OK` count==1 (temp-copy-mutant pattern); `TestCorpusPredicateMatchesShellAnchors` requires 2 regex literals count==1; `TestEmptyExpectedReleaseSetFailsLoudly` reads source at `:504` |
| V12 | CI lanes | `grep -n verify_ail .github/workflows/ci.yml` + read `:80-110,:150` | job 1 runs the script at `:96`, exports `WORLD_PKG_AILANG_BIN` only (no `AILANG_BIN`); job 2 sets `GOTOOLCHAIN: go1.25.6`; z3 installed in both jobs (`:150` comment + `TestZ3PinDeclaredOnceAndInstalledInBothJobs`) |
| V13 | AST gate floors, doesn't pin, walked files | Read `host/broker/invoke_boundary_test.go:209-213` | `len(files)==0` fatal + `len(files)<30` fatal; no exact equality |
| V14 | nothing committed asserts the module count (negative-existence) | `grep -rn 'modules11\|EXACT_TOTAL_MODULES\|module(s)' scripts/ host/ .github/` | sole hit: `verify_ail.sh:243` (the print — in-scope control firing) |
| V15 | proposed names unallocated (negative-existence) | `grep -rn LEG1_MODULES scripts/ host/ .github/`; `ls host/verifygate/module_manifest_gate_test.go`; `grep -c 'world/types.ail' scripts/verify_ail.sh` | 0 hits (control `grep -c ROOTS` = 3); file does not exist; count = 1 |
| V16 | replay fixture is a Go-test input, not a Leg-1 candidate | `head host/replay/testdata/transition_fixture.ail`; `grep -rln transition_fixture host/` | declares `module host/replay/testdata/transition_fixture`, imports `world/*`; consumers `host/replay/replay.go`, `host/replay/replay_test.go` ("the archived released AILANG binary EXECUTES this source during replay") |
| V17 | instrument versions | `/tmp/ailang-v0300/ailang --version`; `command -v z3; z3 --version` | `AILANG v0.30.0` commit `e37b370`; `/opt/homebrew/bin/z3`, `Z3 version 4.16.0 - 64 bit` |
| V18 | *(re-derives controller C1)* the Go gate runs package tests CONCURRENTLY | `sed -n '100,113p' scripts/verify_go.sh`; `grep -n '\-p ' scripts/verify_go.sh` | `:108` `go test ./... -count=1`, and no `-p 1` anywhere (grep 0 hits; positive control in the same read: `-count=1` at `:108` and the `-race` rerun below it are visible), so Go's default package-level parallelism applies |
| V19 | *(re-derives controller C2)* `host/boundary` enumerates AND reads the LIVE `world/` on every run | grep/sed over `host/boundary/allowlist_world_test.go`; scope control on the SAME file: `grep -c 'design_docs'` vs `grep -c '\.ail'` | `enumerateAIL` `:197` walks `filepath.Join(root, "world")` `:199` collecting `*.ail`; `checkAILGroup` `:284` reads every enumerated file via `ov.readFile` `:293`, any error propagates to `t.Fatal` `:775-777`; `root` = the live checkout (`repoRoot`, `runtime.Caller`, `:80-87`). Scope control: 0 `design_docs` hits vs 8 `.ail` hits in the same file — the walk covers `world/` only (why the deletion arm alone would not collide via it; both arms isolated regardless, §8) |
| V20 | *(re-derives controller C3; refutes one review's attribution)* `host/broker` never reads `.ail` — the collider is `host/boundary` | `grep -c '\.ail' host/broker/invoke_boundary_test.go` (known-positive control: same pattern on `host/boundary/allowlist_world_test.go`) | 0 hits (control = 8); the broker walker skips `_test.go` (`:145`, `:218`) and non-`.go` (`:149`). gemini-3-1-pro's cited file does not collide; the real collider both enumerates and reads, which is stricter |
| V21 | isolated root works: `$0`-relative root resolution + the four-item copy set is sufficient | copied `scripts/verify_ail.sh`, `scripts/testdata/ailang_release_observed.txt` (`:102-106` FATALs without it), `world/`, `design_docs/sketches/` into a `mktemp -d` dir; ran `<copy>/scripts/verify_ail.sh` from `$HOME` with the pinned binary | all 11 ai-check lines name the COPY's modules; `✓ 4/4 required world/ identities verified across 11 module(s)`; `✓ all 14 required named tests pass (failed_tests=0)`; then rc=127 at `:302` `./scripts/verify_world_package.sh: No such file or directory` — the deliberate Leg-3 stop; ~1.7 s wall for legs 1–2 |
| V22 | the defect reproduces INSIDE the isolation | wrote AC2's 3-line stray to `<copy>/world/_stray.ail`; reran the copied script | `ai-check world/_stray.ail` ran; `✓ 4/4 … across 12 module(s)` — V5's silent acceptance, inside the copy; an isolated arm measures the same defect the live V5 measured |
| V23 | the isolated runs left the live checkout byte-unchanged | after V21+V22: `git status --porcelain`; `ls world/_stray*`; `shasum -a 256 design_docs/sketches/storejournal.ail` | porcelain = only this doc's own pre-existing untracked entry; no live `world/_stray*` (glob no-match); storejournal sha `adf0760b…5079` unchanged — the AC7 protocol, executed once already |
| V24 | committed `runGate` is live-hardwired — arms need a root-parameterized variant | Read `ail_binary_gate_test.go:52-56` | `exec.Command(filepath.Join(repoRoot, "scripts", "verify_ail.sh"))` with `cmd.Dir = repoRoot`; its env-hygiene blocklist is reusable as-is in `runGateAt` |

| V25 | round-2 objection arm 1 — `\|` in a filename does NOT defeat the pin | isolated temp tree, control asserted (2 files created); `world/types.ail` + `world/types.ail\|junk.ail`; build the line-delimited triples; `cut -d'\|' -f3 \| LC_ALL=C sort`; `cmp` against a 1-entry allowlist | extraction yields an EXTRA row `junk.ail`; `cmp` **fails** → pin HOLDS. gpt5-6-sol's stated exploit is **refuted**; the encoding is still mis-parsed |
| V26 | round-2 objection arm 2 — a newline in a filename corrupts the encoding but is fail-safe | same protocol, name = `a.ail<NL>contracts.ail`; **control: NUL-counted file count asserted == 2 before any verdict was read** (a first attempt created only 1 file and was DISCARDED as vacuous) | 2 real files → **4 lines** in the triple file; `mod` column garbled (2 empty rows); `cmp` **fails** → pin HOLDS. NUL-end-to-end control on the identical tree also correctly fails → the adopted fix is at least as strong and needs no parse |

## Quorum verification log

| Round | Reviewers | Verdict | Disposition |
|---|---|---|---|
| R1 `2026-08-12T11:52:17Z` | `gpt5-6-sol` reject · `gemini-3-1-pro` reject · controller pass | **BLOCKED** | Both rejected on the SAME defect: §5's arms mutated the live checkout while `go test ./...` runs packages concurrently. Controller confirmed the premise first-party AND corrected its attribution (`verify_go.sh:108` has no `-p 1`; the colliding package is `host/boundary` — `allowlist_world_test.go:197` enumerates and `:293` reads the live `world/` — **not** `host/broker`'s AST gate, which filters `.go` at `:149`). Routed to the designer for a full revision pass; fixes adopted from both reviewers' `proposed_fix` (isolated `t.TempDir()` root, real copied script, `cmd.Dir`, byte-unchanged criterion, `LC_ALL=C`). |
| R2 `2026-08-12T12:07:09Z` | `gpt5-6-sol` reject · `gemini-3-1-pro` **pass** · controller pass | **BLOCKED** (1 objection) | Sole remaining objection: the `base\|rel\|mod` line-delimited encoding is unsafe for a gate that must reject unexpected paths. Both carve-out limbs met — a concrete reviewer-authored `proposed_fix`, and no dispute of the design DIRECTION (a robustness/completeness objection about an encoding). **Narrow-refinement carve-out APPLIED** (long-ratified in this mission; first use OK'd iter-6/8, applied at iters 48/49/50/65): gpt5-6-sol's fix adopted **verbatim** — NUL-delimited end to end via indexed arrays, `printf '%s\0'`, `LC_ALL=C sort -z`, NUL-aware/shell-quoted diagnostics (§4.2, §4.2a). Controller measured the objection before applying it: its two named exploits are **refuted** (V25, V26) while the class it names is real; recorded as such rather than laundered into agreement. |

Both `presentCount` readings were **2 external reviewers present, `absent_reviewers` empty** — no
N−1 degrade, and the controller's own `--controller-verdict` was not load-bearing for either
synthesis. Metered: R1 $0.069134 + R2 $0.099862.

## Related

- Queue items 12 (this), 13 and 16 — `design_docs/world-mission.md`
- `design_docs/planned/w-m1-ailang-hardening.md` §D5 — the gate's origin
- `scripts/verify_world_package.sh:86-96` — the ported identity-allowlist pattern
- Iteration 71 log entry — the one-shot AC this item supersedes
