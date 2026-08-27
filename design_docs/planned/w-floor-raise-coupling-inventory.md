# w-floor-raise-coupling-inventory — six files move on every floor raise, and the map lives in a commit message

**Status**: Planned
**Date**: 2026-08-27
**Queue item**: 43, `w-floor-raise-coupling-inventory` (clause-2, sprint-surfaced at `P6.V`,
6th site controller-first-party, queued iter-127)
**Estimated**: ~0.1 day (two verbatim-proposed documentation blocks — one new §S8 in
`design_docs/coding-standards.md`, one comment block at the head of `scripts/verify_ail.sh`;
one ~50 LOC static enforcement test in a NEW file `host/verifygate/floor_raise_inventory_test.go`;
one one-line stale-anchor repair in `docs/SELF_MOD_PUBLISH.md`; **zero executable lines of any
gate changed, zero thresholds moved, zero existing test assertions touched**)
**Designer**: `claude-fable-5` (design-doc-creator, iteration 132)
**Toolchain boundary**: every command below was run first-party in this worktree at `476069d`
(clean tree; porcelain 0 re-checked after every arm), shell `zsh`, darwin/arm64, 2026-08-27.
The pinned released binary `/tmp/ailang-v0300/ailang` (`AILANG v0.30.0`) drove the full gate
run (V16); no `.ail` source is written or changed by this design. Probe programs ran from
`/tmp` against `host/pkgproj` inside this module, writing nothing into the tree (V13, V14, V20).

> **Thesis:** the `P6.V` floor raise (10 → 11, PR #96) proved that adding one identity to
> `REQUIRED_VERIFIED` moves SIX files — the law module, its projection copy, both constants in
> `verify_ail.sh`, the ready-packet golden, the runbook digest table, and the verifygate
> pristine-control marker — and the only complete enumeration of that set is the commit message
> of `699f592` (V9). Three roles each found a different subset: the design doc's file table
> named 2 (V10), the sprint planner found 3 more, and the 6th surfaced only when the
> controller's out-of-sandbox gate re-run met a genuine rc=1. The queue row itself now carries
> first-party proof that the coupling is live: it quotes the marker as `✓ 10/10 …`, which was
> true at iter-127 and is stale today — the landed literal is `✓ 11/11 required world/
> identities verified across 11 module(s)` at `module_manifest_gate_test.go:128` (V4). The
> repair is a MAP, not a mechanism: publish the inventory verbatim in two durable homes
> (`coding-standards.md` §S8 and a comment block at the head of `verify_ail.sh`), with the
> regeneration recipe beside it — executed this session in its no-op form, byte-reproducing
> (V15, V16) — and `interfaceHash`'s non-movement stated with its measured mechanism (V12–V14)
> so nobody "fixes" the third digest. Explicitly refused: deriving the verifygate control from
> `EXACT_TOTAL_VERIFIED` — a control derived from the value it checks is vacuous by
> construction (§Non-Goals).

## The finding in one paragraph

The verify gate is complete by construction — every coupled site reds *individually* when it
drifts (the projection byte-check at step 3/9, the golden `cmp -s` at step 9/9, the AC28
doc↔golden digest binding, the marker `Contains` on gate output) — but **no artifact
enumerates the set**, so a floor-raiser discovers the sites by iterating against red gates,
one at a time, in whatever order their sandbox lets them see. `P6.V` paid that cost in full:
the executor's sandbox could not run the package gate, so the 6th site (the verifygate
control string) surfaced only in the controller's mandatory out-of-sandbox re-run, labelled
`UNINFORMATIVE UNDER SANDBOX` by the executor — honestly, and correctly under its own rules
(row 43, V11). The generalisation the row records: *a gate that is complete by construction
can still have a coupling surface that is complete only by memory, and the memory lives in
whoever last raised the floor.* The memory currently lives in `699f592`'s commit message
(V9), which the next floor-raiser will not read; the `w-mcp-projection` file table that
should have carried it named 2 rows for `P6.V` and is historical now (V10). This item moves
the map to where the next raiser's hands will be: the head of the file that holds both
constants, and the standards document every sprint reads.

## Premises

Each premise is one or more Verification Log rows; a claim without a row does not appear here.

- **P1 — the six-file touch set is real and first-party**: `git show --stat 699f592` lists
  exactly 7 files, and 6 of them (all but the sprint plan) are the coupling surface; the
  commit message enumerates them with reasons, and the amend commit `1e8a018` demonstrates the
  cascade a second time on a comment-only law change: `contentHash` and `tarballSHA256` moved,
  `tarballBytes` 8362 → 8806, `interfaceHash` byte-identical in the diff (V9).
- **P2 — the row's own transcription is stale, which is the point**: the ratified queue row
  quotes `✓ 10/10 …` (world-mission.md:4403, V11); the landed literal is `✓ 11/11 required
  world/ identities verified across 11 module(s)` at `module_manifest_gate_test.go:128` (V4).
  A value transcribed from a document is a claim about that document, not about the repo —
  the coupling site moved underneath its own queue row in one iteration.
- **P3 — no seventh site on the swept patterns**: `identities verified` over
  `*.go|*.sh|*.yml` → exactly 3 hits (test:128, verify_ail.sh:328, :405);
  `EXACT_TOTAL_VERIFIED|REQUIRED_VERIFIED` same scope → 7 hits, all `scripts/verify_ail.sh`
  (V5, V6, reproducing the controller's V-G). A third, independent pattern — the three
  *current digest values themselves*, repo-wide minus `.git` — finds exactly two live homes
  (the golden and `SELF_MOD_PUBLISH.md:85-87`) plus historical records (sprint JSONs, archived
  status stamps, implemented sprint plans) that record past digests and are NOT coupled (V7).
  Counts are true only inside these scopes; see §Inventory for the scope statement.
- **P4 — the regeneration recipe byte-reproduces at base**: `./scripts/build_world_package.sh`
  is idempotent (sha256 of all 6 projected files identical before/after, porcelain 0, V15),
  and the full pinned gate run is green end-to-end: `✓ 11/11 …`, 40 named tests, 9/9 package
  steps, canonical JSON equal to the golden byte-for-byte, rc=0 (V16).
- **P5 — `interfaceHash`'s inputs are seven manifest strings, and nothing else**:
  `pkgproj.InterfaceHash` (pkgproj.go:87-105) hashes `name:`, `edition:`, optional `ailang:`,
  sorted `export:` lines, sorted `effect:` lines; it never opens a file (V12). A probe
  computing it from those strings alone — zero `.ail` bytes read — reproduces the committed
  golden value `sha256:d16cc882…` exactly (V13). The counter-probe: appending one export
  moves it to `sha256:0dcf526f…` (V14). So non-movement under a P6.V-shaped raise and
  MANDATORY movement under a module-set-changing raise are the same measured mechanism.
- **P6 — the enforcement design works before it is written**: a standalone probe carrying the
  proposed test's verbatim logic ran four arms — real `verify_ail.sh` at base → loud
  instrument-failure RED (markers absent; an empty enumeration cannot pass), fixture with the
  proposed block → PASS, fixture minus one site → RED naming `docs/SELF_MOD_PUBLISH.md`,
  fixture minus the END marker → loud RED (V20). The proposed test name does not collide and
  its naive `-run` form is `[no tests to run]` rc=0 at base — the AC1 vacuity trap (V18, V19).
- **P7 — a live stale anchor of exactly the class this item repairs**:
  `docs/SELF_MOD_PUBLISH.md:39` says `verify_ail.sh` "invokes verify_world_package.sh at
  :224"; the call site is `:403` (V3, V17). Line numbers rot; literals do not — which is why
  every inventory row below carries both, and why the runbook repair replaces the number with
  a positional literal rather than a fresher number.

### Design Freeze

- **Zero executable-line changes to any gate.** The `verify_ail.sh` edit is a comment block
  plus one pointer comment line above `REQUIRED_VERIFIED`; the non-comment-line diff against
  `476069d` must stay empty (AC5; base-measured empty, V23). `EXACT_TOTAL_VERIFIED` stays 11,
  `EXACT_TOTAL_TESTS` stays 40, the marker at test:128 is byte-untouched.
- **The comment block must not contain the token `world-publish`.** AC30
  (`runbook_stageb_test.go:341-381`) scans every `scripts/*.sh` and FATALS on any line
  carrying it (V24). It also must not quote the marker string verbatim — the inventory names
  the site, it does not mint a fourth copy of the literal.
- **One new test file, no existing test edited.** `floor_raise_inventory_test.go` is static
  text-scanning only: no `requirePinned`, no `AILANG_BIN`, runs in any lane (the V14 pattern
  from row 42's doc, re-measured here as V18/V21).
- **The two homes are hand-authored duplicates deliberately, and the test binds both.** The
  test's own needle list is a third independently-authored copy — an expectation, like the
  marker itself, never derived from either document (§Non-Goals rule applied to this item's
  own mechanism).
- **`coding-standards.md` is ratification-class** (its own footer). The authorization is the
  ratified queue row itself, which names the file as the durable home (world-mission.md:4403,
  V11); the §S8 text is proposed VERBATIM below so the PR ratifies exactly what lands.
- **`packages/world-core/**` is never hand-edited** — regeneration only (P4); this sprint
  does not touch it at all.

## The inventory (the load-bearing artifact)

Scope statement, per the mission's counting rule: **six sites on the patterns
`identities verified`, `EXACT_TOTAL_VERIFIED|REQUIRED_VERIFIED` (over `*.go|*.sh|*.yml`), and
the three current digest literals (repo-wide minus `.git`, historical records excluded)** —
V5, V6, V7. My own third pattern (digest values) found no seventh coupled site; it did find
the class of *historical records* that carry old digests legitimately and must never be
"updated" (V7). Line numbers below are at `476069d` and will rot; the literals will not.

### Tier 1 — every verified-identity floor raise (the P6.V shape: new identity in an existing module)

| # | File | Anchor (literal · line @476069d) | What moves | Why |
|---|---|---|---|---|
| 1 | `world/<module>.ail` (P6.V: `world/contracts.ail`) | the new `ensures`-bearing contract | the law itself | the identity must exist and prove before anything pins it |
| 2 | `packages/world-core/world/<module>.ail` | step 3/9 byte-identity check · `verify_world_package.sh:98-112` | projection copy of the same bytes | `verify_ail.sh:403` calls the package gate, which binds `world/*.ail` to the published-package projection by SHA-256 equality (V3, V25) |
| 3 | `scripts/verify_ail.sh` | `REQUIRED_VERIFIED = {` · `:274` AND `EXACT_TOTAL_VERIFIED=11` · `:323` | BOTH constants — the identity entry and the exact total | the manifest is identity-keyed with the total as a secondary check; editing one without the other reds the gate (V3) |
| 4 | `scripts/world_package_ready_packet.golden.json` | single-line canonical JSON · step 9/9 `cmp -s` at `verify_world_package.sh:244` | `contentHash`, `tarballSHA256`, `tarballBytes` — and NOTHING else on this shape | `ContentHash` walks `.ail` bytes and the tarball contains them; the module's new bytes move both (V9, V12) |
| 5 | `docs/SELF_MOD_PUBLISH.md` | digest table · `:85-87` | the `contentHash` and `tarballSHA256` rows; the `interfaceHash` row does NOT change | AC28 (`runbook_stageb_test.go:244`) requires exactly 3 distinct doc digests, each present verbatim in the golden, with a negative control (V24) |
| 6 | `host/verifygate/module_manifest_gate_test.go` | `const marker = "✓ 11/11 required world/ identities verified across 11 module(s)"` · `:128` | the identity count (first `11`); the trailing `11 module(s)` moves only on Tier 2 | the pristine known-positive control for five mutation tests; hand-maintained BY DESIGN — see §Non-Goals (V4) |

Note on site 4/5 breadth: **any** byte change to a projected `world/*.ail` module — a new
named test, a comment — triggers sites 2, 4, 5 (the `1e8a018` amend proved this on a
comment-only change, V9). Sites 3 and 6 are specific to raises of the verified-identity
floor. The test floor (`EXACT_TOTAL_TESTS=40`, `REQUIRED_TESTS`) lives entirely inside
`verify_ail.sh` and is out of this row's scope.

### Tier 2 — a raise that changes the module/export set (additional sites; enumerated from code, V22/V25, not from a rehearsed raise)

- `scripts/verify_ail.sh` — `LEG1_MODULES=(` · `:140` (the exact `.ail` allowlist)
- `host/verifygate/module_manifest_gate_test.go` — the isolated-copy census
  `files != 13 || ailFiles != 11` · `:85-86`, and the marker's `across 11 module(s)` tail
- `scripts/build_world_package.sh` — `ALLOWLIST` and `EXPECTED_MODULE_COUNT=4` · `:10`
- `scripts/verify_world_package.sh` — `MODULES` · `:33`, `EXPORTS` · `:34`, the frozen
  manifest `want` dict · `:119-122`, the expected tar entries · `:202`
- `packages/world-core/ailang.toml` — the exports list (and through it the golden's
  `exports` array)
- **and `interfaceHash` MOVES, by design** — exports are an input (P5, V14). On this shape a
  *non-moving* `interfaceHash` is the defect.

## Decision — publish the map in two durable homes, bind both with one static test, repair one stale anchor

### (a) `scripts/verify_ail.sh` — comment block at the head, verbatim

Inserted after the existing header comment (after `:29`, before `set -uo pipefail`), plus one
pointer line directly above the `REQUIRED_VERIFIED = {` heredoc at `:274` reading
`# Before adding an identity here: read the FLOOR-RAISE COUPLING INVENTORY at the head of this file.`

```
# ── FLOOR-RAISE COUPLING INVENTORY — charter row 43 ──────────────────────────
# Raising the verified-identity floor (a new identity in REQUIRED_VERIFIED below)
# touches SIX files. No single gate enumerates the set — each site only reds
# individually, in discovery order. Edit all six in the SAME commit:
#   1. world/<module>.ail                              the new contract (the law)
#   2. packages/world-core/world/<module>.ail          regenerate: ./scripts/build_world_package.sh
#                                                      (step 3/9 byte-identity; never hand-edit)
#   3. scripts/verify_ail.sh                           BOTH constants: REQUIRED_VERIFIED and
#                                                      EXACT_TOTAL_VERIFIED
#   4. scripts/world_package_ready_packet.golden.json  contentHash, tarballSHA256, tarballBytes;
#                                                      take the new line from step 9/9's diff
#   5. docs/SELF_MOD_PUBLISH.md                        digest-table rows for contentHash and
#                                                      tarballSHA256 (host/runbook binds doc↔golden)
#   6. host/verifygate/module_manifest_gate_test.go    the pristine-control marker string —
#                                                      hand-maintained BY DESIGN; deriving it from
#                                                      EXACT_TOTAL_VERIFIED would make the control
#                                                      vacuous (row 43's evaluator refutation)
# interfaceHash does NOT move on this shape of raise: host/pkgproj.InterfaceHash hashes
# manifest fields only (name, edition, ailang bound, sorted exports, sorted effects), never
# .ail bytes. Do not "fix" the third digest. It MUST move only when the export/module set
# changes — and that shape touches MORE sites: LEG1_MODULES below, the 13/11 census and
# "module(s)" count in module_manifest_gate_test.go, build_world_package.sh's allowlist,
# verify_world_package.sh's MODULES/EXPORTS/frozen manifest, packages/world-core/ailang.toml.
# Recipe + rationale: design_docs/coding-standards.md §S8.
# Enforced by: host/verifygate/floor_raise_inventory_test.go.
# ── END FLOOR-RAISE COUPLING INVENTORY ───────────────────────────────────────
```

(The block deliberately never contains the token the AC30 scanner kills on, and never quotes
the marker literal — Design Freeze.)

### (b) `design_docs/coding-standards.md` — new §S8, verbatim

Appended after §S7, before the ratification footer. A new top-level section rather than an
S6 subsection: S6 states a *principle* (honest gates) that is stable; this is an
*operational map* with a different lifecycle — it must be re-verified on every floor raise —
and it needs its own anchor (`§S8`) for the script comment and future queue rows to cite.

```markdown
## S8 — The floor-raise coupling inventory (added 2026-08-27, row 43)

*A gate that is complete by construction can still have a coupling surface that is complete
only by memory, and the memory lives in whoever last raised the floor.* Raising the
verified-identity floor touches **six files**; at `P6.V` three roles each found a different
subset, and the full set existed only in commit `699f592`'s message. The map:

| # | File | What moves |
|---|---|---|
| 1 | `world/<module>.ail` | the new contract (the law itself) |
| 2 | `packages/world-core/world/<module>.ail` | projection copy — regenerate with `./scripts/build_world_package.sh`, never hand-edit |
| 3 | `scripts/verify_ail.sh` | BOTH constants: `REQUIRED_VERIFIED` and `EXACT_TOTAL_VERIFIED` |
| 4 | `scripts/world_package_ready_packet.golden.json` | `contentHash`, `tarballSHA256`, `tarballBytes` |
| 5 | `docs/SELF_MOD_PUBLISH.md` | the `contentHash` and `tarballSHA256` digest-table rows |
| 6 | `host/verifygate/module_manifest_gate_test.go` | the pristine-control marker string — hand-maintained; deriving it from `EXACT_TOTAL_VERIFIED` would make the control vacuous (S6) |

**`interfaceHash` does NOT move** on this shape of raise — `host/pkgproj.InterfaceHash`
hashes manifest fields only, never `.ail` bytes — so do not "fix" the third digest. It MUST
move only when the export/module set changes, and that shape touches more sites (the
allowlists and censuses; full list in the inventory block at the head of
`scripts/verify_ail.sh`).

Recipe, all six in the SAME commit: edit 1 and 3 by hand → `./scripts/build_world_package.sh`
(2) → run the pinned gate; step 9/9 reds printing the golden diff whose `+` line is the new
golden (4) → copy the two moved digests into the runbook table (5) → update the marker (6) →
re-run to green; `go test ./host/runbook/` binds 5↔4. Enforced by
`host/verifygate/floor_raise_inventory_test.go`.
```

### (c) `host/verifygate/floor_raise_inventory_test.go` — one static test, `TestFloorRaiseInventoryNamesEveryCoupledFile(t)`

Logic prototyped and measured this session (V20): read `scripts/verify_ail.sh`; locate the
`FLOOR-RAISE COUPLING INVENTORY` … `END FLOOR-RAISE COUPLING INVENTORY` markers and
`t.Fatalf` if either is absent or misordered (**the known-positive control: an empty or
vanished enumeration fails loudly, never passes as zero**); assert the block contains each of
eight needles — `packages/world-core/world/`, `REQUIRED_VERIFIED`, `EXACT_TOTAL_VERIFIED`,
`world_package_ready_packet.golden.json`, `SELF_MOD_PUBLISH.md`,
`module_manifest_gate_test.go`, `interfaceHash`, `does NOT move` — then read
`design_docs/coding-standards.md`, assert the `## S8` heading exists, and assert the same
eight needles in that section, so the two hand-authored homes cannot drift apart silently.
Declared residual: site 1 (`world/<module>.ail`) has no distinctive greppable token — a
`world/` needle matches everything and would be vacuous — so it is carried by prose in both
homes and by needle-8's coupling (an identity lands via `REQUIRED_VERIFIED`). No
`AILANG_BIN`, no subprocess; reuses `repoRoot` only.

### (d) `docs/SELF_MOD_PUBLISH.md:39` — one-line stale-anchor repair

Replace `./scripts/verify_ail.sh   # invokes verify_world_package.sh at :224` with
`./scripts/verify_ail.sh   # its final leg invokes verify_world_package.sh`. The line number
was two refactors stale (`:403` today, V3/V17) — the repair removes the rotting number
rather than refreshing it, keeping the literal that cannot rot. The token
`verify_world_package.sh` survives, so the AC30 control count and `runbook_commands_test.go:147`
are unaffected (V24); `./scripts/verify_world_package.sh` also stands verbatim at `:25`.

## Non-Goals

1. **Do NOT make the verifygate control derive its expected marker from
   `EXACT_TOTAL_VERIFIED`** (or from `REQUIRED_VERIFIED`'s cardinality, or from any value the
   gate computes). This was the controller's first instinct at `P6.V` and the evaluator
   refuted it, and the refutation is binding here: a control whose expectation is computed
   from the value it checks **cannot fail when that value is wrong** — it verifies that
   substitution works, not that the gate verifies anything. If `EXACT_TOTAL_VERIFIED` were
   silently corrupted to 7, a derived marker would expect `✓ 7/7 …`, observe it, and pass:
   vacuous by construction, which S6 forbids. The literal at `:128` is correct *because* it
   is decoupled and hand-maintained — its cost (one hand edit per raise) is its value (an
   independent observer). **What is missing is the map, not the mechanism.** This item adds
   zero derivation anywhere: both homes and the test's needle list are independently
   authored; the test compares, never generates.
2. **No executable gate changes.** No threshold, no assertion, no strength change; AC5 pins
   the non-comment-line diff of `verify_ail.sh` empty against `476069d`.
3. **No auto-generated inventory.** A generator's sweep scope would silently become the map
   (the same defect one level up); the hand map plus a needle test is the honest form.
4. **Not the test floor.** `EXACT_TOTAL_TESTS`/`REQUIRED_TESTS` raises live in one file and
   cascade through sites 2/4/5 only via `world/*.ail` bytes; noted in the inventory, not
   expanded.
5. **Not repairing `w-mcp-projection.md`'s historical table** — the row declares it
   historical; this inventory supersedes it as the living map.

## Alternatives rejected

1. **Derive the marker** — §Non-Goals 1. Rejected as mechanism, kept as the rule this doc
   carries forward (it already governed row 42's design).
2. **Single home** (only `coding-standards.md`, or only the script comment): the P6.V
   evidence is that discovery is role-local — the executor reads the script, the designer
   reads the standards. Two homes bound by one test costs one duplication and buys both
   audiences; drift between them is machine-caught (V20 arm c).
3. **A gate step that enumerates the six files** (e.g. verify_ail.sh asserting their mtimes
   move together): a heuristic over commit shapes, red on legitimate partial work, and an
   executable-line change this item forbids. Rejected.
4. **Refreshing the runbook's `:224` to `:403`**: mints the next stale number; P7 is the
   measured proof this class rots. Positional literal instead.

## Ordering

Gated on nothing. Neighbours named and not absorbed: **row 42** (landed, `58c8f7f`) enforced
the repro-module floor binding this inventory cites as a pattern-sibling; **row 44** owns
`run.sh` CI inertness; **row 49** owns the token-counting-control class. A floor-raiser's
obligation after this lands, one sentence: open `scripts/verify_ail.sh`, read the block, edit
all six in the same commit, and expect `interfaceHash` to hold still unless the module set
changed.

## Files to Create/Modify

- **CREATE** `host/verifygate/floor_raise_inventory_test.go` (~50 LOC; logic prototyped V20;
  name collision-free V19).
- **MODIFY** `scripts/verify_ail.sh` — comment block after `:29` + one pointer comment above
  `:274`; **zero non-comment lines changed** (AC5).
- **MODIFY** `design_docs/coding-standards.md` — §S8 appended verbatim from §(b).
- **MODIFY** `docs/SELF_MOD_PUBLISH.md` — the `:39` line only, per §(d).

No other files. `packages/world-core/**`, the golden, the marker, `world/*.ail`,
`verify_world_package.sh`, `build_world_package.sh`, `ci.yml` — untouched.

## Conflict Surface

- **AC30's zero-needle** (`runbook_stageb_test.go:341-381`) scans `scripts/*.sh` for
  `world-publish` and fatals on any hit; the block avoids the token (Design Freeze). Its
  known-positive control counts `verify_world_package.sh` in `verify_ail.sh` with a `< 1`
  bound (`:364-374`, V24) — the block only increases that count. Measured green at base
  (`go test ./host/runbook/` rc=0, V21).
- **The five module-manifest tests** copy `scripts/verify_ail.sh` into an isolated root and
  run it (`newIsolatedGateRoot`, `:48-89`): comment lines travel inert; the census counts
  files (13) and `.ail` files (11), not lines. The `mutateCopiedScript` anchors
  (`LEG1_MODULES=(` block, the `mods+=` line) are not inside the new block.
- **`host/runbook` AC28/AC29** bind `SELF_MOD_PUBLISH.md`'s digests and phrases; the `:39`
  edit touches neither (V17's repo-wide sweep found the line bound by nothing; V24 lists the
  bound phrases).
- **`coding-standards.md`** is read by every sprint role; §S8 appends and renumbers nothing.
- **Future Tier-2 raises** will edit the block itself (the `13/11` census values it cites);
  the test's needles are file-name-shaped, not value-shaped, so value drift inside the block
  cannot red the test — declared residual, same class as prose everywhere.

## Deferred Scope

- **Row 44** (`run.sh` inert in CI) and **row 49** (token-counting controls prove mention,
  not testing) — named, untouched.
- **A Tier-2 rehearsal** (actually adding a module and walking all ~11 sites): out of budget
  for a 0.1d documentation item; Tier 2 is enumerated from code reads (V22/V25) and labelled
  so.
- **The red-arm of the regeneration recipe** (step 9/9 printing the new golden in its diff):
  not executed this session — it requires mutating the tree, which this role is forbidden to
  do. It is labelled UNVERIFIED-THIS-SESSION in the recipe row (V16 note), with provenance:
  `P6.V` and its amend executed exactly this path twice, first-party, two days ago (V9), and
  `w-evidence-grade-mapping-sprint-plan.md` V13 records the same recipe with a golden-regen
  helper and a passing byte-identity control.

## Acceptance Criteria

Each AC carries its base observation on the unmodified tree at `476069d`, run this session.

- **AC1 — the test exists, RUNS, and passes, in run-existence form.**
  `env -u AILANG_BIN go test ./host/verifygate/ -run
  '^TestFloorRaiseInventoryNamesEveryCoupledFile$' -count=1 -v` → rc=0 with exactly 1
  `=== RUN` and 1 `--- PASS`; a paired nonsense pattern (`-run 'TestNoSuchInventoryZZZ'`)
  prints `[no tests to run]`. **Base: the verbatim command → `ok … [no tests to run]`, rc=0
  (V18)** — the naive form is green at base measuring nothing; the `=== RUN` clause is the
  binding form and is red at base (0 of 1).
- **AC2 — the script block landed and names every site.**
  `awk '/FLOOR-RAISE COUPLING INVENTORY/,/END FLOOR-RAISE COUPLING INVENTORY/' scripts/verify_ail.sh`
  emits a non-empty block containing all eight needles of §(c), and
  `grep -c 'FLOOR-RAISE COUPLING INVENTORY' scripts/verify_ail.sh` → 3 (begin, END line's
  repetition, pointer-free; the exact expected count is fixed by the landed text and recorded
  by the sprint). **Base: 0 occurrences (rc=1), with the same-file known-positive control
  `grep -c 'Leg 1'` → 6 firing in the same call (V17)** — the zero is a measurement, not a
  dead grep.
- **AC3 — §S8 landed.** `grep -c '^## S8' design_docs/coding-standards.md` → 1. **Base: 0
  (rc=1), control `grep -c '^## S6'` → 1 in the same call (V17).**
- **AC4 — the non-movement statement is in BOTH homes.** `grep -c 'interfaceHash'
  scripts/verify_ail.sh` ≥ 1 AND `grep -c 'interfaceHash' design_docs/coding-standards.md`
  ≥ 1. **Base: 0 hits across both files (rc=1), control `grep -n 'EXACT_TOTAL_VERIFIED=11'
  scripts/verify_ail.sh` → `:323` firing in the same call (V17).**
- **AC5 — gate strength untouched.**
  `diff <(grep -v '^[[:space:]]*#' scripts/verify_ail.sh) <(git show
  476069d:scripts/verify_ail.sh | grep -v '^[[:space:]]*#')` → empty, AND the full pinned
  gate re-runs green with the identical banner (`✓ 11/11 …`, 40 tests, 9/9, rc=0). **Base:
  diff empty by identity (V23); gate green with that exact banner (V16).** Green-at-base by
  design — this AC measures the sprint's constraint, and its teeth are M6.
- **AC6 — the runbook anchor repair.** `grep -c 'at :224' docs/SELF_MOD_PUBLISH.md` → 0 AND
  `grep -c 'its final leg invokes verify_world_package.sh' docs/SELF_MOD_PUBLISH.md` → 1 AND
  `go test ./host/runbook/ -count=1` → rc=0. **Base: 1 / 0 / rc=0 `ok … 3.077s` (V17, V21)**
  — the first two clauses red at base, the third green at base and re-run to prove the edit
  broke no binding.
- **AC7 — hygiene.** `go vet ./host/verifygate/` rc=0 and `gofmt -l host/verifygate/` empty.
  **Base: both green (V21).**

## Non-Vacuity — named RED mutation for every added assertion

All arms mutate the production side (the two documents), never the test. The prototype
already ran the load-bearing arms at design time against fixtures (V20); the sprint re-runs
every arm against the landed test, each restored byte-identically (sha256, house recipe),
porcelain 0 after every arm.

| # | Exact edit | Expected RED | Design-time status |
|---|---|---|---|
| M1 | delete the `docs/SELF_MOD_PUBLISH.md` line from the script block | test reds naming the missing needle | **RUN on fixture: `FAIL inventory block omits "docs/SELF_MOD_PUBLISH.md"` (V20 arm c)** |
| M2 | delete the entire block (or only its END marker) | test reds through the marker fatal, never a silent zero | **RUN: real base file → instrument RED; END-marker-only deletion → instrument RED (V20 arms a, d)** |
| M3 | delete §S8 from `coding-standards.md` | test reds on the `## S8` heading assertion | sprint-run (same `Contains` mechanism V20 measured) |
| M4 | edit `does NOT move` in the block to `does move` | test reds on the exact-phrase needle | sprint-run |
| M5 | delete one needle from §S8 while leaving the script block intact | test reds on the cross-home clause — the homes may not drift apart | sprint-run |
| M6 | change `EXACT_TOTAL_VERIFIED=11` → `12` (the AC5 teeth) | AC5's non-comment diff is non-empty AND the gate reds (`expected exactly 12 … got 11`) | sprint-run; the gate branch is the long-established `:324` refusal |

Green control for all arms: the unmutated post-sprint tree passes AC1–AC7.

## Open Decisions

- **OD-1 — §S8 as a top-level section vs an S6 subsection.** *Default: §S8*, for the
  lifecycle argument in §(b); an S6 subsection would bury an operational map inside a stable
  principle and give the script comment no clean anchor. Reviewer may overrule; the text
  moves unchanged either way.
- **OD-2 — should the test also pin the runbook repair (AC6's needles)?** *Default: NO.* A
  prose needle on a doc line that no behaviour depends on is ceremony (row 42's residual
  reasoning); AC6's greps at sprint time are the check, and AC28/AC30 already guard what
  matters in that file.

## Verification Log

All rows run first-party by the designer at `476069d` (clean tree, V1), shell `zsh`,
darwin/arm64, 2026-08-27. Rows reproducing a controller measurement from the iteration-132
brief are marked (C·x). KP = known-positive control in the same call. Probes V13/V14/V20 are
standalone `/tmp` programs run against this module; nothing was written into the tree
(porcelain 0 re-checked, V1).

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Worktree `476069d`, clean before and after every arm | `git rev-parse HEAD && git status --porcelain \| wc -l` (re-run after V15, V20) | `476069d3…`, `0`; `0` after every arm |
| V2 | The pinned binary is present and is the release | `ls -la /tmp/ailang-v0300/ailang && /tmp/ailang-v0300/ailang --version \| head -1` | 91,826,738 bytes; `AILANG v0.30.0` |
| V3 (C·V-A/V-B) | Call site and both constants, located | `grep -n "verify_world_package" scripts/verify_ail.sh`; `grep -n "EXACT_TOTAL_VERIFIED\|REQUIRED_VERIFIED" scripts/verify_ail.sh` | `:403` call; hits at `:274 :298 :323 :324 :325 :328 :405` — 7, all one file. Reproduces the brief |
| V4 (C·V-D) | The marker literal is `11/11`, so the queue row's `10/10` transcription is stale | `grep -n 'required world/ identities verified across' host/verifygate/module_manifest_gate_test.go` | `:128: const marker = "✓ 11/11 required world/ identities verified across 11 module(s)"` |
| V5 (C·V-G) | Sweep 1: 3 sites on `identities verified`, code+CI scope | `grep -rn "identities verified" --include='*.go' --include='*.sh' --include='*.yml' .` + KP `grep -c "identities verified" scripts/verify_ail.sh` | exactly 3 hits: test:128, verify_ail.sh:328, :405; KP → 2 (pattern fires in-scope) |
| V6 (C·V-G) | Sweep 2: all `EXACT_TOTAL_VERIFIED\|REQUIRED_VERIFIED` sites are one file | same grep, that pattern, same scope | 7 hits, all `scripts/verify_ail.sh` (the V3 lines) |
| V7 | Sweep 3 (my pattern): the three current digest VALUES live in exactly two coupled homes | `grep -rn "0c8c60616e592dc0\|d16cc88270ff4c4e\|d4ff710e4850cd70" --exclude-dir=.git .` | live: `scripts/world_package_ready_packet.golden.json:1`, `docs/SELF_MOD_PUBLISH.md:85-87`. All other hits are historical records (`sprint_w-decision-lifecycle-freeze.json`, status archive, implemented sprint plans, world-mission.md) — records of past packets, NOT coupled sites |
| V8 | Golden consumers read it at runtime; none embeds digests | `grep -rln "world_package_ready_packet.golden" --exclude-dir=.git .` cross-checked against V7 | consumers incl. `cmd/world-publish/*`, `host/pkgproj/readypacket*.go`, `host/runbook/runbook_stageb_test.go`, `scripts/verify_world_package.sh` — none carries a digest literal (V7 found none there) |
| V9 | THE SIX-FILE TOUCH SET, first-party: P6.V and its amend | `git show --stat --format='%H %s' 699f592`; `git show 1e8a018 -- docs/SELF_MOD_PUBLISH.md scripts/world_package_ready_packet.golden.json` | 7 files (6 + sprint plan); message enumerates all six with reasons and states `interfaceHash correctly does NOT move`. Amend diff: `contentHash` and `tarballSHA256` rows changed, the `interfaceHash` row (`d16cc882…`) byte-identical; message: `tarballBytes 8362 -> 8806, interfaceHash correctly unchanged` |
| V10 | "The design doc named 2": the historical table's P6.V rows | `sed -n '583,592p' design_docs/planned/w-mcp-projection.md` (and the merged row at `:344`) | exactly two P6.V rows: `world/*.ail` and `scripts/verify_ail.sh` |
| V11 | The queue row, verbatim, including its stale `10/10` quote and the durable-home directive | `sed -n '4403,4429p' design_docs/world-mission.md` | row 43 quotes `✓ 10/10 required world/ identities verified across 11 module(s)`; names `coding-standards.md` + a `verify_ail.sh` head comment as the homes; carries the derivation refusal |
| V12 | `InterfaceHash` mechanism: manifest fields only, no file reads | read `host/pkgproj/pkgproj.go` in full | `InterfaceHash` `:87-105` hashes `name:`/`edition:`/optional `ailang:`/sorted `export:`/sorted `effect:` lines; no `os.Open`/`ReadFile` in it. `ContentHash` `:47-85` walks `.ail` bytes (`file:%s\n` + data); `CreateTarball` `:107-169` tars `ailang.toml` + `.ail` (+`AGENT.md`/`assets/`). Version deliberately excluded (`:37-42` comment; `TestInterfaceHashIgnoresTheVersion`) |
| V13 | Non-movement MEASURED: the committed value derives from seven manifest strings alone | `/tmp/iface_probe.go`: `pkgproj.InterfaceHash` over the frozen manifest literal, zero `.ail` reads; `go run` from repo root | `sha256:d16cc88270ff4c4eaaa583e644d3ea30e2e4b2e36f95fd7108d920046cdb4083` — equal to the golden's committed `interfaceHash`. Corroborates charter `world-mission.md:3981` (PE.A's ADT-constructor change: content/tarball moved, interface held) |
| V14 | The counter-arm: an export-set change MOVES it | `/tmp/iface_probe2.go`: same manifest + appended export `world/newmod` | base `d16cc882…`, grown `sha256:0dcf526f…`, `moved=true` — Tier 2's mandatory movement is the same mechanism |
| V15 | Regeneration recipe, projection arm: idempotent at base | `shasum -a 256` all 6 package files → `./scripts/build_world_package.sh` → re-`shasum` → `diff`; porcelain | `allowlisted modules: iterated=4 wc-l=4`, `projected 4 modules`; diff empty (byte-identical); porcelain 0 |
| V16 | Regeneration recipe, verification arm: the full pinned gate is green and byte-reproduces the golden | `AILANG_BIN=/tmp/ailang-v0300/ailang WORLD_PKG_AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` | rc=0. `✓ 11/11 required world/ identities verified across 11 module(s)`; `✓ all 40 required named tests pass`; 9/9 package steps incl. `4/4 projection hashes equal`, `3 full hashes agree`, `canonical JSON equals committed golden byte-for-byte`. (The binary's `Observatory: 400MB` stderr warning appeared and broke nothing — the gate's stdout-only discipline held.) **The red-arm (step 9/9 diff printing a NEW golden) was NOT run — it requires tree mutation this role is forbidden; labelled UNVERIFIED-THIS-SESSION, provenance V9 (P6.V ran it twice)** |
| V17 | Every AC base state, with in-scope KPs | the grep/awk pairs quoted per-AC in §Acceptance Criteria | inventory marker 0 (rc=1) / KP `Leg 1`=6; `^## S8` 0 (rc=1) / KP `^## S6`=1; `interfaceHash` 0 hits in both homes (rc=1) / KP `EXACT_TOTAL_VERIFIED=11` → `:323`; `at :224` present at `docs/SELF_MOD_PUBLISH.md:39` and it is the ONLY repo hit for `invokes verify_world_package\|:224` outside design_docs |
| V18 | AC1's vacuity trap at base | `env -u AILANG_BIN go test ./host/verifygate/ -run '^TestFloorRaiseInventoryNamesEveryCoupledFile$' -count=1` | `ok … 0.311s [no tests to run]`, rc=0 — bare rc is vacuous at base; the `=== RUN` clause is the repair |
| V19 | No test-name collision; the package's 29 existing tests enumerated | `grep -h '^func Test' host/verifygate/*_test.go`; `grep -c 'FloorRaiseInventory' host/verifygate/*_test.go` | 29 `func Test` lines, none matching; count 0 in all four files (rc=1) with the non-empty listing as the same-call KP |
| V20 | THE PROTOTYPE: the proposed test logic, four arms, fixtures in `/tmp` | `/tmp/inv_probe.go` (verbatim proposed logic) against (a) real `scripts/verify_ail.sh`, (b) fixture + proposed block, (c) fixture minus the SELF_MOD line, (d) fixture minus END marker | (a) `FAIL instrument: … markers absent … (begin=-1 end=-1)` — base-RED, loud; (b) `PASS`; (c) `FAIL inventory block omits "docs/SELF_MOD_PUBLISH.md"`; (d) `FAIL instrument … (begin=29 end=-1)`. M1/M2 measured before the test exists |
| V21 | Hygiene + runbook bindings green at base | `go test ./host/runbook/ -count=1`; `go vet ./host/verifygate/`; `gofmt -l host/verifygate/ \| wc -l` | `ok … 3.077s` rc=0; vet rc=0; gofmt 0 |
| V22 | Tier-2 anchors exist where cited | `grep -n 'LEG1_MODULES=(' scripts/verify_ail.sh`; `grep -n 'files != 13' host/verifygate/module_manifest_gate_test.go`; `grep -n 'EXPECTED_MODULE_COUNT=4' scripts/build_world_package.sh`; `grep -n 'readonly MODULES\|readonly EXPORTS\|readonly GOLDEN\|readonly PACKAGE_DIR' scripts/verify_world_package.sh` | `:140`; `:85`; `:10`; `:16 :32 :33 :34` |
| V23 | AC5's base identity | `diff <(grep -v '^[[:space:]]*#' scripts/verify_ail.sh) <(git show 476069d:scripts/verify_ail.sh \| grep -v '^[[:space:]]*#') \| wc -l` | `0` |
| V24 | The runbook test bindings this item must not break, located | read `host/runbook/runbook_stageb_test.go` `:228-291` and `:330-382`; `grep -n 'verify_ail\|verify_world_package\|invokes' host/runbook/*.go` | AC28 at `:244`: regex `sha256:[0-9a-f]{64}` `:232`, golden-KP `want 3` `:259`, doc exact-3 `:263`, membership loop `:270-276`, distinctness `:279-285`, all-f negative control `:287-290`; AC30 `:341-381` scans `scripts/*.sh` + ci.yml for `world-publish`, fatal on any hit, KP `verify_world_package.sh`-in-`verify_ail.sh` bound `< 1` at `:364-374`; `runbook_commands_test.go:147` requires `./scripts/verify_world_package.sh` in Stage A (stands at doc `:25` regardless of the `:39` edit) |
| V25 (C·V-C) | The nine-step gate's structure and the step-3/step-9 anchors | read `scripts/verify_world_package.sh` in full | step banners 1/9–9/9; step 3/9 per-module SHA-256 equality `:98-112`; step 9/9 golden `cmp -s` `:244` with `diff -u` on failure; `GOLDEN` `:32`; helper computes `contentHash/interfaceHash/tarballSHA256/tarballBytes` via `pkgproj` `:163-188` |

## Related Documents

- [`../implemented/w-canary-control-does-not-survive-a-floor-raise.md`](../implemented/w-canary-control-does-not-survive-a-floor-raise.md)
  — row 42 (`P42`, PR #98 → `58c8f7f`): the pattern-sibling (a floor raise silently disarming
  a control); its Design Freeze already cites this row's derivation refusal as binding.
- [`w-mcp-projection.md`](w-mcp-projection.md) — `P6.V`'s parent; its file table (`:344`,
  `:587-588`) is the "named 2" historical record this inventory supersedes.
- `design_docs/planned/w-mcp-projection-p6v-sprint-plan.md` — the sprint that paid the
  discovery cost; commit `699f592` carries the only prior enumeration.
- `docs/SELF_MOD_PUBLISH.md` — site 5, and the home of the P7 stale anchor this item repairs.
- `design_docs/world-mission.md` — row 43 at `:4403` (the ratified directive, including the
  homes and the refusal); `:3981` (PE.A's first-party `interfaceHash` non-movement
  measurement, corroborated by V13).
- [`../implemented/w-evidence-grade-mapping-sprint-plan.md`](../implemented/w-evidence-grade-mapping-sprint-plan.md)
  — its V13 is the house precedent for golden regeneration with a byte-identity control.
