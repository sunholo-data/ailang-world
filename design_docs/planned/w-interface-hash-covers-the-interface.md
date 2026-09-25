# w-interface-hash-covers-the-interface — carry upstream's interface identity v2 beside the manifest-coverage `interfaceHash`, and stop calling v1 an interface digest

- Status: **planned** · Date: 2026-09-25 · Designer: claude-opus-5-5 (iter-188) · Base commit: `1a1d929` · Owning queue row: **27** · Scope class: kernel correctness (clause-1) · Verify profile: go-host (no `.ail` change)

## §1 Problem

Row 27, restated with this iteration's own measurements. Every claim here has a row in §12.

- **P1. `interfaceHash` covers the manifest, not the interface.** `pkgproj.InterfaceHash` (`host/pkgproj/pkgproj.go:87`) hashes `name`, `edition`, the optional `ailang` bound, the sorted export module names and the sorted `effects.max`. It never opens a `.ail` file (V1).
- **P2. The row's stimulus leaves it unchanged.** Add `| ProbeArm(HashRef)` to the exported `Evidence` ADT in a scratch copy of `packages/world-core`, and `hash_v1` stays `sha256:d16cc882…4083`. Upstream's `hash_v2` moves `b25fe031…` → `0a46f29a…`, and the signature count goes 53 → 54 [`ailang pkg quality --json`, v0.41.0] (V4, V5).
- **P3. The honesty defect has three sites.**
  - `docs/SELF_MOD_PUBLISH.md:84-88` shows `interfaceHash` to the attended operator as one of "the identity of what you are about to publish permanently" (V14).
  - `host/pkgproj/pkgproj.go:37-40` says v1 is "the shape of what the package EXPORTS" (V1).
  - Item 17's §8.3 (`design_docs/implemented/w-validated-proven-evidence-boundary.md:1052`) opens with the false sentence *"The `world/core` interface hash changes because a public ADT changes"* (V15).
- **P4. Upstream already built the missing digest.**
  - Upstream `internal/pkg/hasher.go` `InterfaceHash` is the same manifest-only function (V2).
  - `internal/pkg/hasher_v2.go` `InterfaceHashV2` (prefix `sha256:ifacev2:`) writes `iface.HashProjection` of each exported module, then the same manifest fields. It returns the hash and a sorted, de-duplicated signature list.
  - It first shipped in tag **v0.35.0**, via commits `db32c6c64` and `a223e7274` (V2).
  - The pinned v0.41.0 exposes it via `ailang pkg quality --json`, key `.interface`. `signatures` there is a **count**, not the list (V3).
  - The registry already records a v2 for our published package: `world/core@0.1.0` serves `interface_hash_v2 = sha256:ifacev2:b25fe031…`, which equals our pinned local computation (V7).

## §2 Decision: arm (c), upstream's v2. No human A/B is needed.

| Arm | Verdict | Measured reason |
|---|---|---|
| (a) change `InterfaceHash` to cover the interface | **Rejected** | v1 re-implements the registry's `interface_hash` field. Reconcile compares it to the served `d16cc882…` (V7, V13). Any change to its bytes makes every reconcile of the immutable, already-published 0.1.0 resolve `conflict`. It would also break step 7's CLI cross-check, which compares v1 against `ailang publish --dry-run`'s "Interface hash" line (V12). |
| (b) rename only (`manifestHash`) | **Rejected** | This is dishonest by omission: there would still be no interface digest anywhere World looks. It also renames a key that is frozen in three wire formats (payload, approval scope, golden). |
| (c) keep v1, honestly labelled, and add a real interface digest | **Chosen** | The digest is **upstream's `InterfaceHashV2`**, not one World invents. The registry already serves it for 0.1.0 (V7), so reconcile gains a fourth independent check for free. |

**Why no A/B for Mark.** The row expected an (a)-vs-(b) compatibility call. Measurement removes it:
- (a) is not a trade-off. It is a guaranteed `conflict` against an immutable published artifact.
- (b) does not satisfy the row.
- (c) leaves every recorded 0.1.0 artifact byte-identical (§3 D6).

What remains is one future-facing consequence, named in D3 for the controller's judgement. It does not block this row.

## §3 Design

### D1 — Obtain v2 at the host boundary from the pinned binary (Q1)

A new file `host/pkgproj/iface.go` holds two functions and one type.

**`ParseQualityInterface(out []byte) (InterfaceIdentity, error)`** is pure. It refuses, never defaults. It rejects:
- non-JSON input;
- a `schema` other than `ailang.package-quality/v1`;
- a missing `interface` section;
- `hash_v2` that is absent or does not match `^sha256:ifacev2:[0-9a-f]{64}$`;
- `hash_v1` that is absent or does not match `^sha256:[0-9a-f]{64}$`.

**`QueryInterface(dir, manifest, ailangBin) (InterfaceIdentity, error)`** is the effect:
- It runs `ailang pkg quality --json --no-run .` with `cmd.Dir = dir` and `cmd.Env = childenv.Scrubbed(os.Environ())`.
- It reads **stdout only**, the same discipline and reasoning as `CrossCheck` (`pkgproj.go:214-229`).
- It is bounded by the caller, exactly as `CrossCheck` is. The only production caller, step 7, already runs under `run_bounded 120`.
- It accepts **exit 0 or 2**. Upstream documents 2 as "one or more gates". Under `--no-run`, v0.41.0 *always* exits 2 with a spurious `PUB015 _smoke.ail failed` gate, even on the pristine package (V8, §11 O1). Any other exit is refused even when stdout parses (the arm that kills `MUT-RC1-ACCEPTED`).
- It then applies the **free cross-check**: if the binary's `hash_v1` ≠ `pkgproj.InterfaceHash(manifest)`, it refuses with both values named. This is a second, independent check of the Go re-implementation of v1 and of `worldCoreManifest`, alongside step 7's dry-run prefix compare.

The `InterfaceIdentity` type is `{V1, V2 string; Signatures int}`. `Signatures` is recorded, not gated. A `signatures ≥ 1` guard was prototyped and removed because no real input can reach it on its own. The only v2-unbuilt output the tool produces also lacks `hash_v2`, so the v2 guard always fires first, and the mutation survived (§6).

**Why `--no-run`.** Identity should not depend on executing package code. v2 is identical with and without it [V8: same `b25fe031…`, 53].

**Priced alternatives:**

| Source of v2 | Cost | Verdict |
|---|---|---|
| Pinned binary, `pkg quality --json` | 0.28–0.65 s per call [go test timings, V17]. The binary is already mandatory in CI (`AILANG_BIN` and `WORLD_PKG_AILANG_BIN`, both pinned to v0.41.0, V11) and already spawned by step 7. | **Chosen** |
| Registry's served `interface_hash_v2` | Exists only **after** publish. Using it as the expectation would make reconcile compare the served value with itself, which is vacuous. | Rejected as a source. Kept as reconcile's **observation** (D4). |
| Re-implement `iface.HashProjection` in Go | Upstream's projection lives in `internal/`, like v1's (DD-1). It would be a second drifting copy of compiler internals. | Rejected |

**Failure behaviour**, all loud and all before any value is used:

| Condition | Result |
|---|---|
| Binary absent | `exec` error |
| Non-JSON output | parse error |
| Missing key | named refusal |
| Exit 1 | refusal |
| v1 disagreement | refusal naming both values |

Tests fail with `t.Fatal`, never skip, when `AILANG_BIN` is unset. This is the `requirePinned` pattern (`host/verifygate/ail_binary_gate_test.go:39`), not the skipping `pinnedBinary` pattern of `host/replay`.

### D2 — v2 joins the ready packet, generated by the gate

- `pkgproj.ReadyPacket` gains `InterfaceHashV2 string` with JSON key `interfaceHashV2`. Alphabetically it goes between `interfaceHash` and `package`, which is the gate's `sort_keys` order.
- `ReadyPacketFields` and `Field()` gain the key.
- `RecomputeReadyPacket(dir, manifest, compilerVersion, interfaceHashV2)` takes v2 as a **parameter**. This follows the precedent the file already sets for `compilerVersion`: "the one field that cannot be recomputed locally is … a PARAMETER" (`readypacket.go:26-30`).
- `scripts/verify_world_package.sh` step 7's helper also calls `pkgproj.QueryInterface` and emits `interfaceHashV2=`. Step 9's python requires it and writes it.
- The committed golden gains one key. Every existing byte of it is unchanged.

The Go parser is the **only** parser of the tool's output, because the shell consumes the helper's `key=value` lines.

### D3 — `world-publish`'s fence gets a second source for v2

- `cmd/world-publish/fences.go` gains `frozenInterfaceHashV2 = "sha256:ifacev2:b25fe031…e8d8"`. This is the **registry-recorded** v2 of immutable `world/core@0.1.0` (V7).
- `requireUndriftedPacket` passes it to `RecomputeReadyPacket`, so R-PACKET-DRIFT compares the golden's v2 against a constant. This is the same two-source argument `frozenCompilerVersion` is stated with.
- `TestWorldCoreManifestMatchesTheCommittedGolden` instead recomputes v2 **live** via `QueryInterface` on the real projection. That binds golden ↔ tree, while the constant binds golden ↔ published 0.1.0.

**Named consequence (controller's call; it does not block this row).**
- The existing fence tests already run the real golden through the real fence. `MUT-FENCE-CONST` was killed by `TestEveryRefusalBranchStopsWithItsExactLineAndHasAPositiveControl` and `TestTheEntrypointReachesTheProductionRefusal` (§6).
- So a **future exported-interface change** to `packages/world-core` will red CI with `STOP fence=packet reason=drift` on `interfaceHashV2` until a version-bump design exists.
- Contract-only and body-only edits do **not** trip it. v2 is unchanged under both (V6), so S8's Tier-1 floor-raise recipe is unaffected.
- **Recommendation: accept.** This is the row's re-break surfacing at the one tool that asserts 0.1.0's identity. S3 already says a behaviour change is a package version bump.
- The alternative (the fence skipping v2) was rejected: it would leave a golden field that no Go fence checks.

### D4 — Reconcile compares four digests, and an absent v2 is not success

In `host/broker/registry_reconcile.go`:
- `RegistryMetadata` gains `InterfaceHashV2 *string` with JSON key `interface_hash_v2`. It is a pointer so that JSON `null` and an absent key are distinguishable from a value.
- `ReconcileConfig` gains `ExpectedInterfaceV2 string`. It is **not** added to `PublishHashes`, so the frozen wires stay untouched (D6).
- **R6** now refuses an empty `ExpectedInterfaceV2`, with the message "requires all four expected digests".
- **P4** runs after P3's three-arm compare:
  - served v2 `nil` → `ReconcileConflict`, detail `served document carries no interface_hash_v2; expected …`;
  - mismatch → `ReconcileConflict`, detail `served interface-v2 hash … does not match the expected …`;
  - match → `succeeded-reconciled`, detail "all four expected digests".
- `cmd/world-publish` `reconcileConfigFor` sets `ExpectedInterfaceV2: packet.InterfaceHashV2`.

**Why absent v2 is `conflict`, not success and not `probe-unavailable`.** The document is present and decodes, and P1 already owns "did not decode". A present document that cannot vouch for an expected digest is a disagreement with the reviewed identity. Upstream stores v2 as `omitempty` (V10), so a registry that stops serving it produces exactly this absence. It is real in the wild: the control `sunholo/auth@0.4.1` is schema v1 with **no** `interface_hash_v2` key (V9). That does not affect the control's role, because the classifier only checks that the control is well-formed JSON (`registry_reconcile.go:391`).

### D5 — Honesty edits

- **`docs/SELF_MOD_PUBLISH.md`:**
  - relabel the row `interfaceHash (manifest coverage only)`;
  - add the row `interfaceHashV2 (exported interface) | sha256:ifacev2:b25fe031…`;
  - "These three digests" → four;
  - §4's field list gains `interfaceHashV2`.
- **`host/pkgproj/pkgproj.go:37-41`:** v1 covers the manifest's export *module names*, never module contents; the interface digest is `QueryInterface`'s v2. The package comment notes that v1 is unchanged at upstream `origin/dev` (V2).
- **`scripts/verify_ail.sh:49`:** the comment names v1 as manifest coverage and points to `interfaceHashV2`.
- **§8.3** of `w-validated-proven-evidence-boundary.md` becomes: *"The `world/core` interface identity v2 (`interfaceHashV2`) changes because a public ADT changes; the manifest-coverage `interfaceHash` (v1) does not — measured `d16cc882 → d16cc882` at iter-81 and again at row 27 (iter-188). Content and tarball hashes …"*, followed by a dated correction note citing row 27. The row deferred this edit specifically to the implementing change (§10.12).
- **`coding-standards.md` S8 is not edited.** It is ratification-class, and its sentence is true of v1 as named. It needs no new inventory row, because v2 does not move on a Tier-1 raise (V6). §11 O3 records an optional clarifying line to propose to Mark.

### D6 — What does NOT change (the 0.1.0 publish path's recorded artifacts stay verifiable)

- `world/registry-publish-request/v1`: `publishPayloadWire`, `publishPayloadFields`, `DecodePublishPayload` (`DisallowUnknownFields`);
- the approval scope grammar: `PublishApprovalScope`, `FormatPublishApprovalScope`, and its strict parser (`approve.go:360-410, 580-590`);
- `publishResultWire`, `PublishHashes`, `recomputePublishHashes`, `comparePublishHashes`;
- `frozenPackageVersion`, `frozenCompilerVersion`, `attendedPhrase`.

The prototype diff touches **none** of `registry_publish.go` or `approve.go`; it touches only `registry_publish_test.go`'s AC10(a) driver map (§8). A recorded durable request or approval for 0.1.0 therefore decodes and validates exactly as before. Binding v2 into the payload or scope belongs to the next version's publish design (§11 O2), because `world-publish` can publish nothing but 0.1.0.

**Z3 / pure core.** No `.ail` file in `world/` or `packages/world-core/` changes. No contract, predicate or inline test is added or removed. `REQUIRED_VERIFIED`/`EXACT_TOTAL_VERIFIED` do not move.

**S3, "why is this not a package?"** This is host-side publication verification. `pkgproj`'s own package comment gives the circularity argument: code that computes the hashes authorizing publication cannot itself be a published package.

## §4 Files

| File | Change |
|---|---|
| `host/pkgproj/iface.go` (new, ~80 LOC) | `InterfaceIdentity`, `ParseQualityInterface`, `QueryInterface` |
| `host/pkgproj/iface_test.go` (new) | T1–T7 (§6) |
| `host/pkgproj/testdata/gen_quality_fixtures.sh` (new) + 3 generated + 3 derived fixtures | §7 |
| `host/pkgproj/readypacket.go` / `_test.go` | field, frozen list, `Field`, parameter; Equal-table row |
| `host/pkgproj/pkgproj.go` | comments only (D5) |
| `host/broker/registry_reconcile.go` / `_test.go` | D4; existing fixtures gain `interface_hash_v2` and `ExpectedInterfaceV2`; R6 text "four" |
| `host/broker/registry_reconcile_v2_test.go` (new) + `testdata/metadata_*.json` (captured) | T8–T11 |
| `host/broker/registry_publish_test.go` | AC10(a) driver `drivePkgprojQueryInterface` for the new subprocess site |
| `cmd/world-publish/fences.go`, `main.go`, `wiring_test.go` | D3, `reconcileConfigFor`, T12 |
| `host/runbook/runbook_stageb_test.go` | T13: a v2 digest guard (the existing regex is blind to it, V16) |
| `scripts/verify_world_package.sh`, `scripts/world_package_ready_packet.golden.json` | D2 |
| `docs/SELF_MOD_PUBLISH.md`, `scripts/verify_ail.sh`, `design_docs/implemented/w-validated-proven-evidence-boundary.md` | D5 |

## §5 Acceptance criteria

- **AC1 (re-break).** T1 is green, and `MUT-V2-IS-V1` reds it.
- **AC2 (negative control).** T2 is green, and `MUT-V2-IS-SOURCE` reds it.
- **AC3.** `AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_PKG_AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_world_package.sh` prints `9/9` and exits 0 with the new golden.
- **AC4.** The same gate with step 9's python emitting no `interfaceHashV2` exits 1 at `ready packet differs byte-for-byte from golden` (`MUT-GATE-PY-OMITS-V2`).
- **AC5 (live, read-only).** `go run ./cmd/world-publish reconcile --store <scratch.db> --registry-origin https://storage.googleapis.com/ailang-registry --probe` prints `state=succeeded-reconciled … detail="served metadata matches all four expected digests"`.
  - The prototype printed exactly this (V18).
  - Base `1a1d929` prints "all three" (V18 control).
  - T8 pins the same outcome offline, against the **captured** served document.
- **AC6.** `go vet ./... && AILANG_BIN=… go test ./...` exits 0, and `./scripts/verify_ail.sh` is green.
- **AC7.** `git diff --stat` shows no change to `host/broker/registry_publish.go` or `host/broker/approve.go`.
- **AC8.** §8.3's first sentence no longer claims that v1 moves. `grep -c 'interface hash changes because a public ADT changes' design_docs/implemented/w-validated-proven-evidence-boundary.md` returns 0. As a control, the same grep against `git show 1a1d929:<path>` returns 1.

## §6 Test plan and mutation matrix: every row was RUN against the prototype

The prototype is a copy of the repo at `1a1d929` under `~/.ailang/state/iter188-designer/proto`. Runner: `mutate.sh` applies one perl substitution, runs the named tests, and restores the file (V19).

| # | Test | Mutation killed | Result |
|---|---|---|---|
| T1 | `TestInterfaceV2MovesWhenAnExportedADTGainsAConstructor`: copy the package, insert `  \| ProbeArm(HashRef)` after `ProofReceipt`; v2 moves, v1 does not, signatures +1 | `MUT-V2-IS-V1` (v2 derived from v1) | **KILLED** (T1 + fixture test) |
| T2 | `TestInterfaceV2IgnoresACommentOnlyEdit`: append a comment; identity unchanged | `MUT-V2-IS-SOURCE` (v2 := content hash) | **KILLED** |
| T3 | `TestQueryInterfaceRefusesAV1Disagreement`: manifest `Edition:"2"` | `MUT-NO-V1-CROSSCHECK` | **KILLED** |
| T4 | `TestQueryInterfaceRefusesAnInfraExitEvenWithValidJSON`: fake binary prints the real fixture, exits 1 | `MUT-RC1-ACCEPTED` | **KILLED** |
| T1 | (same) | `MUT-RC2-REFUSED` (exit 2 treated as failure) | **KILLED** |
| T5 | `TestParseQualityInterfaceFixtures`: pristine = `b25fe031…`/53; ctor ≠ pristine; v2-unbuilt refuses **with the v2-absent text** | `MUT-ABSENT-V2-OK` | **KILLED**. It SURVIVED first, while the test asserted only `err != nil` and the signatures guard produced a bystander error. |
| T6 | `TestParseQualityInterfaceDerivedNegatives`: each derived fixture asserts its own guard's text | `MUT-V2-SHAPE-UNCHECKED`, `MUT-V1-GUARD-OFF`, `MUT-SCHEMA-UNCHECKED` | **KILLED ×3**. All three SURVIVED before T6 existed. |
| — | the signatures ≥ 1 guard | `MUT-SIGS-UNCHECKED` | SURVIVED. **Guard removed from the design** (D1), because no real input reaches it on its own. |
| T7 | `ReadyPacket` Equal table gains an `interfaceHashV2` row | `MUT-PACKET-FIELD-ALIAS` (`Field` returns v1) | **KILLED** |
| T8 | `TestReconcilePublished010OverFourDigests`: captured live document → success, "four" | — (positive arm; the anchor for T9 and T10) | green |
| T9 | `TestReconcileV2MismatchIsConflict`: expected = the ctor v2 `0a46f29a…` | `MUT-RECON-NO-V2-COMPARE` | **KILLED** |
| T10 | `TestReconcileAbsentV2IsNotSuccess`: `null` and absent, derived from the captured doc; instrument: the auth fixture has no v2 key | `MUT-RECON-ABSENT-OK` | **KILLED** |
| T11 | `TestReconcileRefusesEmptyExpectedV2` | `MUT-RECON-R6-V2` | **KILLED** |
| T12 | `TestReconcileConfigCarriesThePacketsInterfaceV2` | `MUT-WP-DROP-EXPECTED-V2` | **KILLED**. It SURVIVED the whole `cmd/world-publish` suite before T12 existed. |
| — | existing fence tests, real golden | `MUT-FENCE-CONST` (last nibble flipped) | **KILLED** by 2 existing tests (D3) |
| T13 | `TestRunbookInterfaceV2DigestAppearsVerbatimInTheCommittedGolden` | `MUT-RUNBOOK-ROW-DROPPED` | **KILLED** |
| AC4 | package gate step 9 | `MUT-GATE-PY-OMITS-V2` | **KILLED** (gate rc=1; control rc=0) |

**Tally (final prototype): 17 mutations run, 17 killed, 0 survivors.** The first pass had 6 survivors:
- 5 were fixed by a stricter assertion or a new arm: `ABSENT-V2-OK`, `V2-SHAPE`, `V1-GUARD`, `SCHEMA`, `WP-DROP`;
- 1 was fixed by deleting an unreachable guard: `SIGS-UNCHECKED`.

The executor re-runs this matrix against the real change. It is a claim until then.

## §7 Parser fixture discipline (Q4)

- **Generated, checked in as produced.** `testdata/gen_quality_fixtures.sh` copies `packages/world-core` into `mktemp -d`, applies the recorded edit, runs `$AILANG_BIN pkg quality --json <dir>`, and writes stdout unmodified. It produces three fixtures:
  - `quality_world_core_pristine.json` (`b25fe031…`, 53);
  - `quality_world_core_probe_ctor.json` (`0a46f29a…`, 54);
  - `quality_world_core_v2_unbuilt.json`. A type-error export gives `compile.ok=false`, `signatures:0`, and **no `hash_v2` key** (V10).

  The script prints the temp path. Paths embedded in the unbuilt fixture's error text are temp paths, never a home directory.
- **Captured, checked in as served.** `host/broker/testdata/metadata_world_core_0.1.0.json` (8143 bytes) and `metadata_sunholo_auth_0.4.1.json`, fetched with `curl -sS` from the public bucket (V7, V9). The capture command is recorded in a header-free sibling `testdata/CAPTURED.txt`.
- **Derived negatives** exist only where no tool output can reach a guard. Each is a named `jq` transform of the pristine fixture, kept in the generator:
  - `.interface.hash_v2 = .interface.hash_v1`;
  - `del(.interface.hash_v1)`;
  - `.schema = "ailang.package-quality/v2"`.

  None is hand-typed.
- **Anti-vacuity.** A parse of non-empty input that yields no `hash_v2` is an error naming the missing field, never `""` with success. T5 proves this on the **real** v2-unbuilt output. That is the exact shape upstream emits (`omitempty`, V10), so a naive `json.Unmarshal` into a `string` field would have returned `""`.

## §8 Conflict surface

- **AC10(a)** (`registry_publish_test.go:1133`) enumerates every production `exec.Command`. The new site made it red until a driver was added: the prototype measured `N = 6 production subprocess launch sites in 6 files` (V17). The driver is required.
- **S8 floor-raise inventory:** golden row 4 now carries `interfaceHashV2`, which does not move on a contract-only raise (V6). No new inventory row is needed.
- **Existing reconcile tests:** `metadataDocument`, `reconcileCfg` and three R6 message assertions change. The prototype diff is +15/−9.
- **Future `world/` ADT or export changes:** these red the fence tests (D3). A queued row that adds a constructor must plan a version bump.
- **Unrelated `interfaceHash`, OUT of scope.** A world-graph object's interface ref is a different concept:
  - `host/daemon/handlers.go` `objectResponse`;
  - `host/broker/record.go:118` (`SumSHA256(semanticID)`);
  - `host/store`, `host/transitionreg`, `host/evidence/validator.go`, `host/workbench`;
  - `world/logepoch.ail`.

  Evidence is in V13.
- Row 29 (`ReadObject`'s zero-`InterfaceHash` invariant) concerns that other concept too, so there is no overlap.

## §9 Non-goals

- No change upstream. No vendored or forked hasher.
- No re-publish of 0.1.0 and no new version.
- No change to the publish payload, the approval scope grammar or the result wire.
- No renaming of `interfaceHash` in any frozen wire.
- No `.ail` or Z3 change.
- No edit to `coding-standards.md`.
- The world-graph object `interfaceHash` is out of scope.
- The pre-existing tension that a floor raise moves the "0.1.0" golden's `contentHash` away from the published bytes is not addressed here (§11 O4).

## §10 Milestones (≈0.8 d total)

- **M1 — `pkgproj` v2 query and parser (0.35 d).**
  - Work: `iface.go`, `iface_test.go` (T1–T6), the generator plus 6 fixtures, and the AC10(a) driver.
  - Accept: `cd <repo> && AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/pkgproj/ ./host/broker/ -run 'InterfaceV2|QueryInterface|ParseQuality|EverySubprocessSite' -count=1 -v` shows all PASS. Re-run `MUT-V2-IS-V1`, `-IS-SOURCE`, `-ABSENT-V2-OK`, `-NO-V1-CROSSCHECK`, `-RC1-ACCEPTED`, `-V2-SHAPE`, `-V1-GUARD`, `-SCHEMA`; each must red.
- **M2 — packet, gate, fence, reconcile (0.3 d).**
  - Work: D2, D3, D4, T7–T13, and the `SELF_MOD_PUBLISH.md` v2 row. The golden is regenerated **by the gate's own red diff** (the S8 recipe), never hand-edited.
  - Accept: AC3, AC4, AC5, AC7, and the full `go vet ./... && go test ./...`. Re-run the remaining 9 mutations.
- **M3 — honesty (0.15 d).**
  - Work: D5 (runbook prose, `pkgproj.go` comments, `verify_ail.sh:49`, §8.3 plus a dated correction note).
  - Accept: AC8, AC6, `./scripts/verify_ail.sh`.
  - The controller files §11 O1 upstream.

## §11 Open items and declared residuals

- **O1 (upstream, for the controller to file on `sunholo-data/ailang` with a mission-control note).** `ailang pkg quality --json --no-run` on a package whose smoke passes emits `PUB015 _smoke.ail failed` and exits 2 (V8). "Not run" is reported as "failed". The design does not depend on it: D1 accepts exit 2 because gates are orthogonal to identity. This is not a workaround.
- **O2.** Binding v2 into the publish payload and approval scope is deferred to the next version's publish design. For 0.1.0, the registry itself refuses PUB005 on publisher/validator v2 skew (V2). That is **only** when the header is sent: a v2 build failure is shadow-mode, meaning it is logged and not refused.
- **O3 (optional, ratification-class, Mark).** S8 could add: "`interfaceHashV2` does not move on a contract-only or body-only raise (measured iter-188)."
- **O4.** The golden is labelled 0.1.0 but floor raises move its `contentHash` away from the published bytes. That is pre-existing and orthogonal. Candidate for a new row.
- **Controller correction to M1.** The controller wrote that the validator "refuses PUB005 on disagreement". Measured: it refuses only when the publisher **sends** `X-Interface-Hash-V2` and the value differs; otherwise v2 is shadow-mode (V2). Also, `pkg quality`'s `signatures` is an integer count; the set appears only in registry metadata as `interface_signatures`.

## §12 Verification Log

All `ailang` runs used `~/.pinned-ailang/ailang`, which reports `AILANG v0.41.0 / Commit: 24ee108`. PATH `ailang` is `v0.42.0-7-g265e6fb7a-dirty` and was **not** used. `S=~/.ailang/state/iter188-designer`; each scratch copy is `cp -R packages/world-core $S/<name>`, and the repo copy was never mutated (`git status --short | wc -l` → 0 at the end).

| V | Claim | Command | Observed |
|---|---|---|---|
| V1 | v1 hashes the manifest only, and the comment calls it "the shape of what the package EXPORTS" | `sed -n 30,110p host/pkgproj/pkgproj.go` | `InterfaceHash` writes `name:`, `edition:`, `ailang:`, `export:<mod>`, `effect:` and nothing else. `ContentHash` is the only `.ail` walker. Lines 37-40 contain that sentence. |
| V2 | Upstream v1 is identical; v2 exists; tag; validator behaviour | `cd ~/dev/sunholo-data/ailang && git fetch origin dev; git show origin/dev:internal/pkg/hasher.go \| grep -A18 '^func InterfaceHash'`; `git show origin/dev:internal/pkg/hasher_v2.go \| grep -n 'ifacev2\|HashProjection'`; `git tag --contains db32c6c64 \| sort -V \| head -1`; `git show origin/dev:cmd/registry-validator/main.go \| sed -n 290,320p` | Same five fields and order (origin/dev `b5502791c`, 2026-09-24). `interfaceHashV2Prefix = "sha256:ifacev2:"`, `iface.HashProjection(j)` at l.41. First tag is `v0.35.0`. Validator: `if claimed := r.Header.Get("X-Interface-Hash-V2"); claimed != "" && claimed != v2Hash` → PUB005 skew; `v2Err` is only logged. |
| V3 | Pinned binary exposes both hashes; pristine values | `(cd $S/pristine && $B pkg quality --json . > ../q_pristine.json); jq -c '.interface' q_pristine.json` | rc=0. `{"source":"server","hash_v1":"sha256:d16cc882…4083","hash_v2":"sha256:ifacev2:b25fe031…e8d8","signatures":53}` |
| V4 | Stimulus moves v2 only | `sed -i '' '29a\  \| ProbeArm(HashRef)' ctor/world/types.ail`, then quality | rc=0, `compile.ok=true` (contracts 10/11: one ensures no longer holds over the widened ADT). v1 `d16cc882…`, v2 `sha256:ifacev2:0a46f29a089fe1add745412d3e95d5b0093c9bd7e05177b6f3aa8daabf14c30e`, 54 |
| V5 | Negative control: comment only | `echo '-- iter188 negative control: comment only' >> comment/world/types.ail`, then quality | v1 `d16cc882…`, v2 `b25fe031…`, 53 (both unchanged) |
| V6 | A floor-raise-style contract or a body edit does not move v2; a new export does (positive control) | perl: add `requires { now >= 0 }` to `validDefer`; swap its body conjuncts; append `export func probeArity`; quality on each | contract: `b25fe031…`/53; body: `b25fe031…`/53; newfn: `sha256:ifacev2:d35379ad…75af`/54. v1 `d16cc882` in all three. |
| V7 | Published 0.1.0 serves v1 and v2 equal to local | `curl -sS -o $S/metadata_world_core_0.1.0.json -w '%{http_code} %{size_download}' https://storage.googleapis.com/ailang-registry/packages/world/core/0.1.0/metadata.json; jq -c '{interface_hash,interface_hash_v2,content_hash,tarball_hash}'` | `200 8143`. `interface_hash` = `d16cc882…`, `interface_hash_v2` = `sha256:ifacev2:b25fe031…e8d8`, content `0c8c6061…` and tarball `44fc9fab…` equal the golden |
| V8 | `--no-run` exits 2 with spurious PUB015 and the same v2; it works offline with a scrubbed env | `env -i HOME=$HOME PATH=/usr/bin:/bin HTTPS_PROXY=http://127.0.0.1:9 AILANG_REGISTRY=http://127.0.0.1:9 AILANG_REGISTRY_VALIDATOR=http://127.0.0.1:9 $B pkg quality --json --no-run .` (pristine); control without `--no-run` | `rc=2`, gates `[{"code":"PUB015","level":"gate","msg":"_smoke.ail failed"}]`, v2 `b25fe031…`/53. Control run: `rc=0`, gates `[]`, same v2. |
| V9 | The control package serves no v2 key | `curl … sunholo/auth/0.4.1/metadata.json; jq -c '{has_v2:has("interface_hash_v2"),schema}'` | `200`; `{"has_v2":false,"schema":"ailang.package-metadata/v1"}`. Positive control on the world/core doc: `{"has_v2":true,"schema":"ailang.package-metadata/v2"}` |
| V10 | A v2-unbuilt quality run omits `hash_v2`; upstream uses omitempty | Append `export func probeBroken(x: int) -> int ! {} { x ++ "s" }`, quality `--no-run`; `git show origin/dev:internal/pkg/quality.go \| grep -n hash_v2` | `rc=2`, `.interface` = `{"source":"server","hash_v1":"sha256:d16cc882…","signatures":0}`, gates PUB000+PUB015. `HashV2 string \`json:"hash_v2,omitempty"\`` at quality.go:150 |
| V11 | CI installs the pinned v0.41.0 as both `AILANG_BIN` and `WORLD_PKG_AILANG_BIN`; step 7 already spawns it | `grep -n 'AILANG_BIN\|v0.41.0' .github/workflows/ci.yml`; `grep -n 'run_bounded\|go run' scripts/verify_world_package.sh` | ci.yml:97 asserts v0.41.0, :99 exports `WORLD_PKG_AILANG_BIN`, :159/:161 `AILANG_BIN`. Gate :244 is `run_bounded 120 … go run "$tmp_helper" "$AILANG_BIN"` |
| V12 | Step 7 cross-checks v1 against the dry-run "Interface hash" line | `sed -n 180,300p host/pkgproj/pkgproj.go` | The `dryRunLine` regex includes `Interface hash`, and `Compare` has an `interface` arm |
| V13 | v1 users vs. the unrelated object `interfaceHash` | `git grep -n "InterfaceHash\|interfaceHash" -- ':!design_docs' \| grep -v _test.go`; `sed -n 36,45p host/daemon/handlers.go`; `sed -n 112,120p host/broker/record.go` | Package sites: pkgproj.go/readypacket.go; broker registry_publish.go (92,128,140,151,202,218,603,614,732), registry_reconcile.go (162,263,365), approve.go (368,378,401,585); cmd/world-publish main.go (426,504,546), fences.go:85; SELF_MOD_PUBLISH.md (66,87); verify_world_package.sh; golden. Object sites: handlers.go `objectResponse.InterfaceHash` from `store.Object`; record.go:118 `InterfaceHash: hashref.SumSHA256([]byte(semanticID))`; store/transitionreg/evidence/workbench/logepoch as listed in §8. |
| V14 | The runbook presents v1 as publish identity | `sed -n 60,95p docs/SELF_MOD_PUBLISH.md` | "These three digests are the identity of what you are about to publish permanently", table row `interfaceHash` `d16cc882…` |
| V15 | §8.3's false sentence | `awk '/^### 8.3/{p=1} p{print; n++} n>6{exit}' design_docs/implemented/w-validated-proven-evidence-boundary.md` | Line 1052: "The `world/core` interface hash changes because a public ADT changes." |
| V16 | The runbook digest regex is blind to v2 | `printf 'sha256:ifacev2:<64 b>\nsha256:<64 a>\n' \| grep -cE 'sha256:[0-9a-f]{64}'`; the prototype with v2 in the golden but not in the doc | `1` (only the plain digest). The prototype's existing `host/runbook` passed **ok** with v2 absent from the doc; the new T13 then red, `names v2 digests []` |
| V17 | Prototype: tests, timings, AC10(a) | `cd $S/proto && AILANG_BIN=$B go test ./host/pkgproj/ -run '…' -v`; `go test ./host/broker/` before and after the driver; `go vet ./... && go test ./...` | T1 0.58s, T2 0.52s, T3 0.28s. Broker before the driver: `AC10(a) re-derived N = 6 … host/pkgproj/iface.go:66 exec.Command` → FAIL; after: `ok`. Full suite: vet rc=0, test rc=0, 21 `ok`, 69 s |
| V18 | Live reconcile: prototype four digests, base three | `go run ./cmd/world-publish reconcile --store $S/empty.db --registry-origin https://storage.googleapis.com/ailang-registry --probe`, in proto and then in the repo at `1a1d929` | proto: `state=succeeded-reconciled package=world/core@0.1.0 … detail="served metadata matches all four expected digests"`. base: `… detail="served metadata matches all three expected digests"` |
| V19 | Mutation matrix | `$S/mutate.sh <ID> <file> '<perl>' <pkg> '<re>'` per §6 row; gate mutation by perl on step 9 python, then `verify_world_package.sh` | As tabulated in §6. Gate mutant: `rc=1`, `✗ ready packet differs byte-for-byte from golden`; unmutated control `rc=0` |
| V20 | Package gate green with the v2 golden | `AILANG_BIN=$B WORLD_PKG_AILANG_BIN=$B ./scripts/verify_world_package.sh` (proto) | `rc=0`, `✓ canonical JSON equals committed golden byte-for-byte`, `✓ world package gate PASSED: 9/9` |
