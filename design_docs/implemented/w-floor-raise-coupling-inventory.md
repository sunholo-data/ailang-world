# w-floor-raise-coupling-inventory — six files move on every floor raise, and the map lives in a commit message

**Status**: Planned
**Date**: 2026-08-27
**Queue item**: 43, `w-floor-raise-coupling-inventory` (clause-2, sprint-surfaced at `P6.V`,
6th site controller-first-party, queued iter-127)
**Estimated**: ~0.1 day (two verbatim-proposed documentation blocks — one new §S8 in
`design_docs/coding-standards.md`, one comment block at the head of `scripts/verify_ail.sh`;
one ~60 LOC static enforcement test in a NEW file `host/verifygate/floor_raise_inventory_test.go`;
one one-line stale-anchor repair in `docs/SELF_MOD_PUBLISH.md`; **zero executable lines of any
gate changed, zero thresholds moved, zero existing test assertions touched**)
**Designer**: `claude-fable-5` (design-doc-creator, iteration 132)
**Revision**: 1 (2026-08-27) — quorum round 1 BLOCKED at full strength; all three objections
applied, one in a discriminating form after its literal fix measured VACUOUS (§Quorum round 1,
V26–V28)
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
> repair is a MAP, not a mechanism: publish the inventory verbatim — Tier 1 ONLY; the
> un-rehearsed Tier-2 enumeration was struck from the published homes at quorum round 1 — in
> two durable homes
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
  At revision the probe was extended to the ninth (site-1) needle and the exactly-once marker
  fatal, eight arms, both homes — including the arm proving the naive site-1 needle vacuous
  (V26, V27).
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
- **The published homes carry Tier 1 ONLY** (quorum round 1, objection 2). The Tier-2 site
  list is asserted-complete (one designer's code reads, V22/V25), not verified-complete, and
  publishing an un-rehearsed enumeration as authoritative reproduces the row's own defect one
  level up. Each home carries the one-line Tier-2 forward reference — but the `interfaceHash`
  MUST-move branch STAYS in both homes: it is a statement about a digest's mechanism (P5,
  V14), not an enumeration of sites, and the reader must still learn not to "fix" a moving
  `interfaceHash` on a module-set change nor a still one on a Tier-1 raise.

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

### Tier 2 — a raise that changes the HASHED MANIFEST FIELDS (exports / effects / name /
### edition / AILANG bound) or the PACKAGED-MODULE CENSUS (DEFERRED SCOPE — NOT PUBLISHED)

Asserted-complete only: enumerated from one designer's code reads (V22/V25 prove the cited
anchors exist; they cannot prove no uncited anchor exists), never from a rehearsed raise —
exactly the evidence grade the published map must not carry (quorum round 1, objection 2:
Tier 1 is backed by three independent sweeps plus a rehearsed commit; Tier 2 has none of
these). This list stays HERE, in the design doc, as the seed for the rehearsal-gated
follow-up row (§Follow-up queue row); the two published homes carry Tier 1 plus a one-line
forward reference only.

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

The pointer line deliberately repeats the phrase `FLOOR-RAISE COUPLING INVENTORY`, so every
extractor — AC2's awk and the test alike — must anchor on the `# ── `-styled marker lines,
never on the bare phrase: an unanchored awk range re-opens at the pointer and runs to EOF
(measured: 32 lines ending in the script tail, vs the marker-bounded block when anchored,
V28; quorum round 1, objection 3).

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
#                                                      step 9/9 reds printing a diff against the
#                                                      committed golden; replace the golden with
#                                                      the diff's new line and re-run
#   5. docs/SELF_MOD_PUBLISH.md                        digest-table rows for contentHash and
#                                                      tarballSHA256 (host/runbook binds doc↔golden)
#   6. host/verifygate/module_manifest_gate_test.go    the pristine-control marker string —
#                                                      hand-maintained BY DESIGN; deriving it from
#                                                      EXACT_TOTAL_VERIFIED would make the control
#                                                      vacuous (row 43's evaluator refutation)
# interfaceHash does not move for .ail byte changes or packaged-module changes that leave
# the hashed manifest fields unchanged. It MUST move when a hashed manifest field changes,
# including name, edition, the optional AILANG bound, exports, or effects. Do not "fix" the
# third digest on a Tier-1 raise. A raise that changes the hashed manifest fields (or the
# packaged-module census) touches additional sites beyond these six; that inventory is
# deferred to a future item pending a first-party rehearsal.
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

**`interfaceHash` does not move for `.ail` byte changes or packaged-module changes that
leave the hashed manifest fields unchanged. It MUST move when a hashed manifest field
changes, including `name`, `edition`, the optional AILANG bound, exports, or effects.**
So do not "fix" the third digest on a Tier-1 raise. A raise that changes the hashed manifest
fields (or the packaged-module census) touches additional sites beyond these six; that
inventory is deferred to a future item pending a first-party rehearsal.

Recipe, all six in the SAME commit: edit 1 and 3 by hand → `./scripts/build_world_package.sh`
(2) → run the pinned gate; step 9/9 reds printing a diff against the committed golden;
replace the golden with the diff's new line and re-run (4) *(step 9/9's red-arm diff format
verified via the V25 mechanism, not re-executed this session — see the V16 note)* → copy the two moved digests into the runbook table (5) → update the marker (6) →
re-run to green; `go test ./host/runbook/` binds 5↔4. Enforced by
`host/verifygate/floor_raise_inventory_test.go`.
```

### (c) `host/verifygate/floor_raise_inventory_test.go` — one static test, `TestFloorRaiseInventoryNamesEveryCoupledFile(t)`

Logic prototyped and measured this session (V20, extended at revision by V26/V27): read
`scripts/verify_ail.sh`; locate the block by the STYLED marker literals
`# ── FLOOR-RAISE COUPLING INVENTORY` and `# ── END FLOOR-RAISE COUPLING INVENTORY`, and
`t.Fatalf` unless each occurs EXACTLY ONCE with begin preceding end (**the known-positive
control: an empty, vanished, duplicated, or misordered enumeration fails loudly, never
passes as zero — and a bounded extractor can never silently run wide, V27 arms g/h**). The
bare phrase must not be the anchor: the pointer line above `REQUIRED_VERIFIED` repeats it by
design (V28). Assert the block contains each of NINE needles — the eight shared needles
`packages/world-core/world/`, `REQUIRED_VERIFIED`, `EXACT_TOTAL_VERIFIED`,
`world_package_ready_packet.golden.json`, `SELF_MOD_PUBLISH.md`,
`module_manifest_gate_test.go`, `interfaceHash`, `does not move for`, plus site 1's
enumerated-row literal `#   1. world/<module>.ail` — then read
`design_docs/coding-standards.md`, extract the text bounded between the `## S8` heading and
the next `##` heading (or EOF), assert the `## S8` heading exists, and assert the eight
shared needles plus site 1's table-row literal `` | 1 | `world/<module>.ail` `` **within
that bounded extract only** (so a future §S9 reusing a needle term cannot satisfy it), so the two hand-authored homes cannot drift apart silently. Site 1's needle is
enumerated-row-anchored BECAUSE THE BARE PATH WAS MEASURED VACUOUS: `world/<module>.ail` is
a strict substring of site 2's row (`packages/world-core/world/<module>.ail`), so a fixture
containing ONLY site 2's row satisfies the naive needle (`grep -c` → 1, V26; probe arms b/e
PASS with the site-1 row deleted, in BOTH homes, V27), while the row-anchored form scores 0
there (rc=1, control firing on the real row) and reds naming the missing needle (V26, V27
arms c/f). The pre-revision text's declared residual — "site 1 has no distinctive greppable
token" — was FALSE for the proposed text (7 occurrences of the literal in this doc's
pre-revision blocks, controller-measured and reproduced, V26); what is true is that the
BARE token is non-discriminating, which the row anchoring repairs. No `AILANG_BIN`, no
subprocess; reuses `repoRoot` only.

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

- **CREATE** `host/verifygate/floor_raise_inventory_test.go` (~60 LOC; logic prototyped V20,
  extended at revision V26/V27; name collision-free V19).
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
- **Future Tier-2 raises** no longer edit the block's site list — since quorum round 1 the
  block carries no Tier-2 enumeration, only the stable forward-reference line; the census
  values formerly cited there now live only in this doc's §Tier 2. Value drift inside the
  block's remaining prose still cannot red the test (the needles are file-name- and
  row-shaped, not value-shaped) — declared residual, same class as prose everywhere.

**PLANNER-SURFACED, CONTROLLER-CONFIRMED (round 3).** Two collisions this section did not name:

- **`TestNoRigAbsolutePaths` (`host/verifygate/ail_binary_gate_test.go:553`) globs
  `host/verifygate/*.go`** and `t.Errorf`s on `/tmp/ailang`, `/Users/` and `/home/runner/`.
  The new `floor_raise_inventory_test.go` lands inside that glob, so a provenance comment
  naming the pinned binary's rig path would red the whole package. Confirmed first-party:
  the glob, the assembled-needle list and the `len(entries)==0` instrument floor are all at
  `:554-:575`. **The new file must contain no absolute rig path in code OR comment**, and the
  acceptance sweep runs the whole `host/verifygate` package rather than only the new test.
- **The pointer comment lands INSIDE a Python heredoc.** §(a) places it "directly above the
  `REQUIRED_VERIFIED = {` heredoc at `:274`"; measured, the heredoc opens at `:270`
  (`python3 - "$mod" "$tmp_json" <<'PY'`) and `:274` is *inside* it. The pointer must
  therefore be a column-0 Python comment, matching the existing one at `:273` — a shell-style
  indented `#` there is Python source, not a shell comment.

## Deferred Scope

- **Row 44** (`run.sh` inert in CI) and **row 49** (token-counting controls prove mention,
  not testing) — named, untouched.
- **A Tier-2 rehearsal AND the Tier-2 publication that waits on it** — see §Follow-up queue
  row. The §Tier 2 list in this doc is asserted-complete (code reads, V22/V25), never
  rehearsed, and is therefore NOT published to either durable home (quorum round 1,
  objection 2).
- **The red-arm of the regeneration recipe** (step 9/9 printing the new golden in its diff):
  not executed this session — it requires mutating the tree, which this role is forbidden to
  do. It is labelled UNVERIFIED-THIS-SESSION in the recipe row (V16 note), with provenance:
  `P6.V` and its amend executed exactly this path twice, first-party, two days ago (V9), and
  `w-evidence-grade-mapping-sprint-plan.md` V13 records the same recipe with a golden-regen
  helper and a passing byte-identity control.

## Follow-up queue row (controller to file)

**`w-floor-raise-tier2-inventory`** — publish the module/export-set-raise coupling
inventory, gated on a first-party rehearsal. This is not scope creep: the quorum-round-1
deferral creates the obligation, and the controller files the row. Closing predicate: a
first-party REHEARSED module/export-set raise (or removal) executed in a worktree — every
touched site recorded from the actual red-gate sequence and the commit diff, completeness
backed by the same evidence grade Tier 1 carries here (independent sweeps plus a rehearsed
commit, the V5/V6/V7 + V9 shape), and — only for a rehearsed HASHED-MANIFEST change — `interfaceHash`
measured to MOVE (the V14 counter-arm, run live); a packaged-module-census-only rehearsal
must instead show `interfaceHash` STABLE (V29 arm 1) — and only then the Tier-2 enumeration promoted from this doc's §Tier 2 into both
durable homes, replacing the forward-reference line. Until that row closes, the published
map stays Tier-1-only and the forward reference is the only Tier-2 claim either home makes.

## Acceptance Criteria

Each AC carries its base observation on the unmodified tree at `476069d`, run this session.

- **AC1 — the test exists, RUNS, and passes, in run-existence form.**
  `env -u AILANG_BIN go test ./host/verifygate/ -run
  '^TestFloorRaiseInventoryNamesEveryCoupledFile$' -count=1 -v` → rc=0 with exactly 1
  `=== RUN` and 1 `--- PASS`; a paired nonsense pattern (`-run 'TestNoSuchInventoryZZZ'`)
  prints `[no tests to run]`. **Base: the verbatim command → `ok … [no tests to run]`, rc=0
  (V18)** — the naive form is green at base measuring nothing; the `=== RUN` clause is the
  binding form and is red at base (0 of 1).
- **AC2 — the script block landed, is cleanly bounded, and names every site.**
  `awk '/^# ── FLOOR-RAISE COUPLING INVENTORY/,/^# ── END FLOOR-RAISE COUPLING INVENTORY/' scripts/verify_ail.sh`
  emits the block bounded by its own markers — asserted as **first line = the begin marker,
  last line = the END marker, and the extracted block `diff`s EMPTY against §(a)'s block**,
  never against a transcribed line count. **Re-derived at `50c3b91`: 28 lines** (V30). The
  earlier "25" was measured by V28 against the PRE-REVISION block and went stale when
  revision 1 and the round-2 carve-out added lines — this doc committing its own thesis'
  defect, caught by the sprint planner. A count is a claim about the text that existed when
  it was taken; AC2 therefore asserts a property, not a number.
  containing all nine script-home needles of §(c), and
  `grep -c 'FLOOR-RAISE COUPLING INVENTORY' scripts/verify_ail.sh` → 3 (begin marker, END
  marker, pointer line — re-derived on the block-carrying fixture at revision, V28, not
  transcribed: the count is 3 precisely BECAUSE the pointer line exists). The awk range MUST
  stay anchored to the `# ── ` styling: the unanchored form re-opens at the pointer line and
  runs to EOF (32 lines vs the bounded 25 on the fixture, V28; quorum round 1, objection 3).
  **Base: 0 occurrences (rc=1), with the same-file known-positive control `grep -c 'Leg 1'`
  → 6 firing in the same call (V17)** — the zero is a measurement, not a dead grep.
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
already ran the load-bearing arms at design time against fixtures (V20; extended at revision
by V26/V27, which also measured the arm the naive site-1 needle CANNOT red); the sprint re-runs
every arm against the landed test, each restored byte-identically (sha256, house recipe),
porcelain 0 after every arm.

| # | Exact edit | Expected RED | Design-time status |
|---|---|---|---|
| M1 | delete the `docs/SELF_MOD_PUBLISH.md` line from the script block | test reds naming the missing needle | **RUN on fixture: `FAIL inventory block omits "docs/SELF_MOD_PUBLISH.md"` (V20 arm c)** |
| M2 | delete the entire block (or only its END marker) | test reds through the marker fatal, never a silent zero | **RUN: real base file → instrument RED; END-marker-only deletion → instrument RED (V20 arms a, d)** |
| M3 | delete §S8 from `coding-standards.md` | test reds on the `## S8` heading assertion | sprint-run (same `Contains` mechanism V20 measured) |
| M4 | edit `does not move for` in the block to `does move for` | test reds on the exact-phrase needle | sprint-run |
| M5 | delete one needle from §S8 while leaving the script block intact | test reds on the cross-home clause — the homes may not drift apart | sprint-run |
| M6 | change `EXACT_TOTAL_VERIFIED=11` → `12` (the AC5 teeth) | AC5's non-comment diff is non-empty AND the gate reds (`expected exactly 12 … got 11`) | sprint-run; the gate branch is the long-established `:324` refusal |
| M7 | delete ONLY the `#   1. world/<module>.ail` row from the script block, site 2's row intact | test reds naming the missing site-1 row needle — and the NAIVE bare-path needle must be shown NOT to red here (that green is the measured vacuity) | **RUN on fixture: anchored → `FAIL inventory block omits "#   1. world/<module>.ail"`; naive bare path under the SAME arm → PASS, i.e. vacuous as the reviewer proposed it (V26, V27 arms b/c)** |
| M8 | delete ONLY site 1's table row from §S8, site 2's row intact | test reds naming the missing §S8 site-1 row needle | **RUN on fixture: anchored → `FAIL S8 omits` the row literal; naive bare path under the SAME arm → PASS, vacuous (V27 arms e/f)** |

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
(porcelain 0 re-checked, V1). Revision rows V26–V28 ran the same way — `/tmp/inv_rev132/`
fixtures and probe, zero tree writes, porcelain re-checked after the revision pass (only this
doc, already untracked, appears).

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
| V26 | REVISION — the reviewer-proposed site-1 needle is VACUOUS as written: bare `world/<module>.ail` is a strict substring of site 2's row | `/tmp` fixture holding ONLY site 2's row: `grep -c 'world/<module>\.ail'`, then anchored `grep -c '^#   1\. world/<module>\.ail'` on the same fixture, then the anchored form on a site-1-bearing fixture (KP), then `grep -c 'world/<module>\.ail'` on the pre-revision doc (C·, controller measured 7) | naive → `1` rc=0 — fires on site 2 ALONE; anchored → `0` rc=1 on the same fixture; KP anchored on the real site-1 row → `1` rc=0; pre-revision doc → `7` with same-file KP `grep -c 'EXACT_TOTAL_VERIFIED'` → `17` firing in the same call |
| V27 | REVISION — the discriminating ninth needle reds exactly where the naive form stays green, in BOTH homes, and the extractor fails loudly on missing/duplicated markers | `/tmp/inv_rev132/inv_probe3.go` (the §(c) logic verbatim, incl. the exactly-once styled-marker fatal), 8 arms: (a) full script fixture + anchored needle, (b) site-1 row deleted + NAIVE needle, (c) same deletion + anchored needle, (d) full §S8 fixture + anchored row needle, (e) §S8 minus site-1 row + NAIVE, (f) same + anchored, (g) END marker deleted, (h) block duplicated | a `PASS`; **b `PASS` — site 1 DELETED and the naive needle stays green: the vacuity, measured**; c `FAIL inventory block omits "#   1. world/<module>.ail"`; d `PASS`; **e `PASS` — same vacuity in the §S8 home**; f `FAIL S8 omits` the row literal; g `FAIL instrument: END marker count=0, want 1`; h `FAIL instrument: begin marker count=2, want 1` |
| V28 | REVISION — the unanchored awk range re-opens at the pointer line and runs to EOF; the styled-anchor form bounds exactly; AC2's occurrence count re-derived, not transcribed | 38-line `/tmp` fixture = 3 header lines + the proposed 25-line block + `set -uo pipefail`, middle lines, the pointer line, a `REQUIRED_VERIFIED` heredoc and 2 tail lines; old-AC2 awk vs `awk '/^# ── FLOOR-RAISE COUPLING INVENTORY/,/^# ── END FLOOR-RAISE COUPLING INVENTORY/'`; `grep -c`/`grep -n 'FLOOR-RAISE COUPLING INVENTORY'` | unanchored → 32 lines, LAST LINE THE SCRIPT TAIL (`echo "tail line 2"`) — over-wide, silently; anchored → exactly 25, first line the begin marker, last the END marker; occurrence count → `3` (begin `:4`, END `:28`, pointer `:32` of the fixture) |

### Quorum round 1 (2026-08-27, full strength — recorded so the next reader does not re-derive it)

All three reviewers present (`absent_reviewers: []`); verdict BLOCKED on two rejects plus one
pass-with-a-catch; the design DIRECTION was undisputed by any reviewer. All three objections
applied in this one protocol-mandated revision:

1. **`gpt5-6-sol`, REJECT (blocking)** — the test omitted site 1, so the durable inventory
   could lose the law-source row while `TestFloorRaiseInventoryNamesEveryCoupledFile` still
   passed; and the pre-revision declared residual ("site 1 has no distinctive greppable
   token") was FALSE for the proposed text — both blocks contain `world/<module>.ail`
   verbatim (controller-measured 7 occurrences, reproduced first-party in V26). **Applied —
   but the reviewer's proposed fix as literally written was itself measured VACUOUS**: bare
   `world/<module>.ail` is a strict substring of site 2's row, so the reviewer's own
   non-vacuity arm (delete only the site-1 row, require RED) stays GREEN under it (V26; V27
   arms b/e, both homes). Applied in a discriminating form instead: enumerated-row literals
   per home (`#   1. world/<module>.ail` in the script block; site 1's table row in §S8),
   measured to red under exactly that arm (V27 arms c/f; mutations M7/M8).
2. **`oc-glm-5-2`, REJECT (blocking)** — the published Tier-2 enumeration was
   asserted-complete, not verified-complete: Tier 1 has three independent sweeps plus a
   rehearsed commit (V5/V6/V7, V9); Tier 2 had one designer's code reads (V22 proves the
   cited anchors exist, not that no uncited anchor exists) — the row's own defect, one level
   up. **Applied as written**: both published homes now carry Tier 1 only plus the one-line
   forward reference; the `interfaceHash` MUST-move branch stays in both homes (a digest
   mechanism, not a site enumeration); publication of Tier 2 is rehearsal-gated via
   §Follow-up queue row.
3. **`gemini-3-1-pro`, PASS with a catch (non-blocking)** — AC2's unanchored awk range
   re-opens at the pointer line (which repeats the bare phrase by design) and runs to EOF.
   **Applied**: the awk anchored to the `# ── ` styling (measured bounded, V28), the test's
   extractor given the same anchoring plus an exactly-once marker fatal (V27 arms g/h), and
   AC2's expected occurrence count re-derived at 3 rather than transcribed.

| V29 | ROUND-2 CARVE-OUT — `interfaceHash` moves for HASHED MANIFEST FIELDS ONLY, never for a packaged-module-census change. Four arms, one variable each, run by the CONTROLLER (not the designer) at `6c34d27` | a Go probe calling `pkgproj.InterfaceHash` directly, arms: (1) packaged-but-unexported module change = hashed fields untouched, (2) an export appended, (3) an effect added, (4) exports permuted | BASE `sha256:d16cc882…` — **equal to the committed `interfaceHash` in `docs/SELF_MOD_PUBLISH.md`, which is the known-positive control proving the probe builds the real manifest**; arm 1 **UNMOVED** (gpt5-6-sol's premise CONFIRMED); arm 2 → `sha256:0dcf526f…` MOVED (reproduces the designer's V14); arm 3 → `sha256:df50680e…` MOVED; arm 4 **UNMOVED** (exports are sorted before hashing). Source read in the same pass: `host/pkgproj/pkgproj.go:87` hashes `name`, `edition`, the optional `ailang` bound, sorted `Exports.Modules` and sorted `Effects.Max` — the packaged-module set never reaches it |

## Quorum round 2 + the narrow-refinement carve-out (controller record)

**Round 2 ran at FULL STRENGTH — `absent_reviewers: []`, all three reviewers present —
and came back BLOCKED**, metered `$0.1484` (round 1 `$0.1246`; cumulative `$0.2730` of the
`$5` ceiling). Verdicts moved: `gemini-3-1-pro` **reject → pass**, `oc-glm-5-2`
**reject → pass**, `gpt5-6-sol` **reject → reject**, on a *different* surface.

**Disposition: the narrow-refinement carve-out, applied as a bounded 2nd revision by the
CONTROLLER using the reviewers' own verbatim text.** Its two conditions both hold: every
surviving blocking objection carries a concrete reviewer-authored `proposed_fix`, and none
disputes the design DIRECTION — round 2's single reject is a *precision* defect in one
published sentence, and the other two reviewers passed the direction outright. The
first-use-needs-Mark gate does not apply: World has applied this carve-out before
(iteration 44, `w-ddl-gate-teeth` DG.B), so it is already ratified for this mission.

**The surviving objection was MEASURED, not forwarded.** `gpt5-6-sol` objected that the
published *"`interfaceHash` … MUST move when the module/export set changes"* is over-broad,
because neither V12 nor V14 covers a module-set-only change. Read at
`host/pkgproj/pkgproj.go:87`, `InterfaceHash` hashes the manifest's `name`, `edition`,
optional AILANG bound, **sorted exports** and sorted effects — it never sees the packaged
module set. Confirmed on four arms (V29), with the probe's base reproducing the committed
digest as its known-positive control. **The reviewer is right**, and its verbatim
replacement sentence is now the text in both homes.

**One mechanical consequence the fix does not mention, and it would have shipped a red
sprint.** The test's shared-needle list and mutation arm M4 both keyed on the exact literal
`does NOT move`, which the reviewer's replacement sentence deletes. Applied blind, the
landed comment block would no longer contain the needle its own enforcement test asserts —
the test would red at sprint time for a reason unrelated to the design. Both the needle and
M4 are realigned to `does not move for`. This is a consequence of the reviewer's own fix,
not a controller-invented resolution.

**Non-blocking catches, both applied verbatim.** `gemini-3-1-pro`: the §S8 extraction in
§(c) was unbounded, so a future §S9 reusing a needle term could satisfy it — the test now
bounds its search between the `## S8` heading and the next `##` (or EOF).
`oc-glm-5-2`: site 4's recipe published a `+`-line diff-format assertion whose red arm was
never executed this session (V16) — softened in both homes to the verified mechanism, with
the evidence grade stated inline.

| V30 | ROUND-3 — §(a)'s block is **28** lines, not the 25 V28 recorded; the number went stale inside this very document when the revision passes added lines | `awk '/^# ── FLOOR-RAISE COUPLING INVENTORY/,/^# ── END FLOOR-RAISE COUPLING INVENTORY/' design_docs/planned/w-floor-raise-coupling-inventory.md | wc -l` at `50c3b91` | **28**. V28's 25 was correct against the pre-revision block and is a claim about that text. Surfaced by the sprint planner, confirmed first-party by the controller. AC2 no longer asserts any count: it asserts the marker bounds and a `diff`-empty against §(a). **This is the item's own generalisation firing on the item: a value transcribed rather than re-derived.** |

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
