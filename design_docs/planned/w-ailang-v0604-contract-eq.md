# w-ailang-v0604-contract-eq — make the kernel's contracts type-check under AILANG v0.60.4

**Status**: Planned, r1. Designer, iteration 393 (unattended), with a working prototype in a scratch clone.
r1 answers quorum round 1 (§9). The controller owns commits. Nothing was published and nothing
was requested of the registry.
**Item**: charter row 260 (`w-ailang-v0604-contract-eq`). Decision `D-WORLD-116` is OPEN with
default **A** applied: World builds its own complete fix-forward on a fresh branch, never touches
PR #565, and never publishes.
**Base**: worktree `sprint/world-iter393-row260` at `d0c9c4ce`, which is origin/dev `842b9331` plus
the iteration-392 record commit.
**Evidence**: `design_docs/verification/world-iter393/designer/`, cited below as `ev/…`.

## 0. Problem, measured

CI's `ailang-code verify gate` runs legs 1–2 on `releases/latest`. The "Install released ailang
binary" step is `base="https://github.com/sunholo-data/ailang/releases/latest/download"`
(`ci.yml:43`, V17a). `releases/latest` is now **v0.60.4**, which type-checks `requires`/`ensures`
clauses (ailang#2283, diagnostic `TC_CONTRACT_001`). The kernel uses `==` in four `ensures`
clauses on types that have no `Eq` instance:

| # | site (base `d0c9c4ce`) | function | `==` on |
|---|---|---|---|
| S1 | `world/types.ail:43` | `gradeOf` | `EvidenceGrade` (an enum) |
| S2 | `world/types.ail:173` | `timeoutOutcome` | `TimeoutOutcome` (an enum) |
| S3 | `world/transitions.ail:29` | `plan` | `HashRef` (`{algo, digest}` record alias; 3 comparisons) |
| S4 | `world/transitions.ail:48` | `verify` | `HashRef` (`result.proposalHash == p.proposalHash`) |

**v0.60.4 reports only the first failing clause in each module.** At base it shows S1 alone (V1).
S2 appears only when S1 is fixed (MUT-2), and S4 appears only when S3 is fixed (V4). The charter
row and the iteration-392 evidence name three sites; S4 is the fourth. With PR #565 applied,
`verify` would have been the next red after `plan`.

**Full v0.60.4 failure set (V1, a per-file `ai-check` sweep of all 32 gated or packaged `.ail`
files; the 7 `host/evidence/testdata/aicheck` fixtures are excluded on purpose):**
- 8 files fail, all with `TC_CONTRACT_001` from the four sites:
  - `world/{types,contracts,transitions}.ail`;
  - their `packages/world-core/world/*` copies;
  - `packages/world-core/_smoke.ail` and `host/replay/testdata/transition_fixture.ail`, through imports.
- 24 files pass:
  - the 10 sketches, the 10 se-tools modules and both `world-demo` files;
  - `world/logepoch.ail` and its package copy.
- No other new diagnostic exists. With S1–S4 fixed, `verify_ail.sh` runs to completion on v0.60.4 (V5):
  - Leg 1: 24 modules, 16/16 `world/` identities;
  - Leg 2c: 138 tests;
  - Leg 2: 40 tests;
  - Leg 2b: 503 tests;
  - Leg 3 (package leg, v0.41.0): steps 1–8 pass; step 9 fails only on the expected identity change.

The only legs-1–2 failures are therefore S1–S4. Leg 3 needs the identity regeneration in §2.

## 1. Decision — fix shape per site

| site | chosen fix | rejected alternative, with the measured reason |
|---|---|---|
| S1 `EvidenceGrade` | `deriving (Eq)` on the type (+1 line) | Restating `gradeOf`'s law without `==` on the result would need a nested `match result {…}` in the `ensures`: a larger and less readable specification. `EvidenceGrade` has no existing equality predicate that `deriving` would duplicate. `deriving (Eq)` compiles on all three binaries (V2). **It adds an exported `Eq[EvidenceGrade]` instance, which neither interface hash records** (§2, V18). |
| S2 `TimeoutOutcome` | `deriving (Eq)` on the type (+1 line) | Same reason as S1. It adds an exported `Eq[TimeoutOutcome]` instance, which is also unhashed. |
| S3 `plan` | restate the `ensures` field by field (`.algo` and `.digest` per HashRef), **on one line** | (B) `HashRef … deriving (Eq)` in `world/logepoch.ail` also works on all three binaries and moves no interface hash (V6). Adding an instance is not disqualifying in itself: S1/S2 add two. B is rejected for two reasons specific to `HashRef`. First, `HashRef` already has a deliberate, exported, proven identity, the field-equality predicate `sameRef` (`logepoch.ail:70–79`, D3), and B would add a second, parallel equality to it. Second, B changes `logepoch.ail`, which passes today. V15 measures B and A as equivalent, so neither choice is forced on soundness grounds. (C) A new nominal type for the digest pair changes every constructor site of `HashRef` (`_smoke.ail`, the replay fixture, the inline tests), which is the largest blast radius. The charter row suggested (B)/(C); the field-wise form adds the fewest exported instances (2, not 3) and keeps logepoch byte-identical. |
| S4 `verify` | the same field-wise restatement of `result.proposalHash`, on one line | Same as S3. `verify`'s own `accepted` conjunct and `applyRevision` already use this idiom. |
| S3/S4 alt (D) | — | **Call the exported `sameRef` in the `ensures`** (`sameRef(result.inputWorld, w.stateRoot) && …`, importing `sameRef`). Measured (V21): v0.60.4 verifies `plan` (3/3), but **v0.41.0, the frozen package compiler, marks `plan` `skipped`** ("SMT declaration closure could not resolve a required sort (cannot encode ensures clause …)"). That drops a `REQUIRED_VERIFIED` identity on the package lane, so D is rejected. Field access and `==` on `string` fields are the form both CI compilers prove. |

S2's line number: `world/types.ail:173` at base; MUT-2 reports `types:174` because S1's added
`deriving (Eq)` line sits above it.

**Each restated clause must stay on ONE line.** `scripts/test_kernel_contracts.py`, which is CI
step "Kernel contract and PUB011 mutation controls", builds its `ensures { false }` control by
replacing the text from `ensures {` up to the end of that **line** (lines 52–55). A multi-line
`ensures` for `plan` or `verify` leaves a dangling fragment. The script then dies on
`PAR_UNEXPECTED_TOKEN at world/transitions.ail:30:3` (V8, red). The single-line form passes (V8, green).

Exact replacement text for both `world/transitions.ail` and its projection:

```
ensures { result.inputWorld.algo == w.stateRoot.algo && result.inputWorld.digest == w.stateRoot.digest && result.proposalHash.algo == proposalHash.algo && result.proposalHash.digest == proposalHash.digest && result.transitionFn.algo == transitionFn.algo && result.transitionFn.digest == transitionFn.digest && result.goal == goal }
ensures { result.proposalHash.algo == p.proposalHash.algo && result.proposalHash.digest == p.proposalHash.digest && result.accepted == (w.stateRoot.algo == p.inputWorld.algo && w.stateRoot.digest == p.inputWorld.digest) }
```

**The restatement proves the same law; this is measured, not assumed (V15).**
- **The semantics of `==`.** For a record with derived `Eq`, a probe module proves both
  directions of `(a == b) ⇔ (a.algo == b.algo && a.digest == b.digest)` (`eqByFields` and
  `fieldsByEq` are `verified`). It also produces counterexamples when only `digest`, or only
  `algo`, differs under a `result == a` spec. Same results on v0.41.0 and v0.60.4.
- **The kernel itself.** The original `==` clauses (via candidate B's `Eq[HashRef]`) and A's
  field-wise clauses were run against the same 7 body mutants of `plan`/`verify` on v0.41.0 and
  v0.60.4. The mutants are MUT-5/5a/5b/6/6a/7/7a; 5a and 5b change one field only. The kill sets
  are identical: 7/7 counterexamples on each of the four candidate × binary pairs, and the
  unmutated control is `verified` throughout.

No function body changes.

**Compiler compatibility of the chosen shape** (V2, V3; per-file `ai-check -timeout 5s`; 30/30
files show `check.passed=true`):

| binary | role in this repo | `world/` contracts verified | delta vs base |
|---|---|---|---|
| v0.60.4 | CI legs 1–2 (`releases/latest`) | 16/16 (types 6, contracts 5, transitions 3, logepoch 2) | base: 2 (only logepoch type-checks) |
| v0.41.0 | package leg, `go-verify` `AILANG_BIN`, frozen compiler | 16/16 | none |
| v0.54.15 | `ci.yml:407` "Install PINNED v0.54.15 ailang tool binary (floor grader $TOOL)", job `go host build + test gate`; used by `scripts/floor` to grade banked benchmark solutions, never on `world/*.ail` (V19) | 16/16 at base AND candidate (6/5/3/2) | none (V19) |
| v0.30.0 | installed by no CI lane (V17c) | 13; `plan`/`verify`/`proposalMatchesWorld`/`verificationMatchesProposal` are `error` and `commitAllowed` is `skipped` | **none**: base is identical (V3) |

`REQUIRED_VERIFIED` and `EXACT_TOTAL_VERIFIED=16` are unchanged. No identity is added or lost,
so the S8 floor-raise inventory is **not** triggered (constants, verifygate marker and
`REQUIRED_TESTS` are untouched).

## 2. Package identity — `world/core` must become 0.1.2

`world/core@0.1.1` is PUBLISHED: the banked served record says `"version": "0.1.1"`,
`"published_at": "2026-10-02T07:58:43Z"` (V17b). Version immutability is a premise **cited
from the repo, not measured**:
- `docs/SELF_MOD_PUBLISH.md:6–7` says "a public publish is **immutable** (`registry-validator`
  409s on an existing version)";
- `host/broker/registry_publish.go:22–23` says "The registry is immutable and the publisher
  cannot recall a version".

No first-party record of an observed 409 exists: the banked 0.1.1 publish README has none (V17b).
Republishing 0.1.1 to test it would be an irreversible public write, so it is not probed. Any byte
change to the four packaged modules therefore needs a new version. PR #565 regenerates the identity **at 0.1.1**, a candidate that could never
be published.

Measured with v0.41.0 through `verify_world_package.sh` step 7/9 (`pkgproj.CrossCheck`/`QueryInterface`) and `publish --dry-run`:

| digest | base (0.1.1) | prototype (0.1.2) | moves? |
|---|---|---|---|
| `contentHash` | `sha256:47351707…a2780` | `sha256:359a82fde07d8e9004786448da37ea3fb07cdc256682efc0fa57240df3c4a72c` | **yes**. Depends only on the module bytes. It did not move when the manifest version, the release kind, the manifest `description`, the CHANGELOG date or the CHANGELOG text changed; it did move when one comment line was appended to a module (V7, V10, V16). The executor must reproduce this exact value |
| `interfaceHash` (v1) | `sha256:d16cc882…4083` | same | no |
| `interfaceHashV2` | `sha256:ifacev2:b25fe031…e8d8` | same (53 signatures) | no. The `frozenInterfaceHashV2` fence in `cmd/world-publish/fences.go:101` stays |
| `tarballSHA256` | `sha256:07d61809…f4a3` | `sha256:b97d0770…64e5` (prototype CHANGELOG) | **yes**. Depends on CHANGELOG/toml text, so the executor's value will differ |
| `tarballBytes` | 10610 | 10852 (prototype CHANGELOG) | **yes**, same dependency (V16: 10851–10876 across CHANGELOG/manifest variants) |

**Release kind: `feature` (decided; unchanged from 0.1.1).**

**S1/S2 add two exported instances, `Eq[EvidenceGrade]` and `Eq[TimeoutOutcome]`, and neither
interface hash records them.** Both `interfaceHash` and `interfaceHashV2` are byte-identical
between base and prototype (V7), and also with a third instance added in candidate B (V6). The
addition is still real (V18): a consumer module whose body is `gradeOf(e) == TESTED` fails to
type-check against the base kernel (`No instance for Eq[EvidenceGrade]`, v0.41.0) and passes
against the candidate, on v0.41.0 and v0.60.4.

So 0.1.2 adds a capability. `feature` is the honest kind, and it is the kind every end-to-end
measurement (V5, V7, AC1's green run) used.

Why the blind spot is acceptable for this release:
- The addition is purely additive. No exported signature changes: 53 signatures, v2 unchanged.
- The kernel's own consumers (`_smoke.ail`, the replay fixture) compile unchanged against 0.1.2 on all three binaries (V3).

The cost is that the interface hashes cannot tell 0.1.1 from 0.1.2 even though code written for
0.1.2 may not compile against 0.1.1. The CHANGELOG `## 0.1.2` section must name the two new
instances. The PR should note the hash gap as an upstream observation: interface hashes do not
cover derived instances. Whether to file it in `sunholo-data/ailang` is the controller's call.
This row does not work around it.

`[release].kind` accepts `security|fix|feature|breaking` (V10, v0.41.0 dry-run). Keeping `feature`
leaves `verify_world_package.sh:136`'s frozen `"release": {"kind":"feature"}` unchanged.

## 3. Files

Prototype diff (kind `feature`, docs prose not yet updated): `ev/prototype-candA-final.diff`,
13 files, +39/−23 (V7).

**M1 — kernel fix, projection and 0.1.2 candidate (one atomic commit; any subset reds a gate):**
1. `world/types.ail` adds `deriving (Eq)` to `EvidenceGrade` and `TimeoutOutcome` (+2). Credit PR #565 for these two lines in the commit body.
2. `world/transitions.ail` restates the `plan` and `verify` `ensures` on one line each (±2).
3. `packages/world-core/world/{types,transitions}.ail`: regenerate with `./scripts/build_world_package.sh`, never by hand (step 3/9).
4. `packages/world-core/ailang.toml`: `version = "0.1.2"`. `kind` stays `"feature"`. Extend the `[release]` comment: "0.1.2 keeps `feature` (row 260: adds `Eq` on two exported enums; no signature changed)".
5. `packages/world-core/CHANGELOG.md`: a new `## 0.1.2 — <date>` section on top. It must name the two added `Eq` instances and the restated `plan`/`verify` clauses. `release_history_test.go` requires newest section == frozen version, and PUB001 requires it to be non-empty. Leave `## 0.1.1`/`## 0.1.0` byte-identical.
6. `scripts/verify_world_package.sh`: step 4 `want` `"version":"0.1.2"` (the `release` kind entry is unchanged); step 7 regex `Publishing world/core@0\.1\.2\.\.\.`; step 9 packet `"version":"0.1.2"`; extend the explanatory comment at lines 127–129.
7. `scripts/world_package_ready_packet.golden.json`: replace with the step-9 diff's `+` line.
8. `cmd/world-publish/fences.go`: `frozenPackageVersion = "0.1.2"`, `attendedPhrase = "publish world/core@0.1.2 irreversibly"`, and the version comment. `frozenInterfaceHashV2` is **unchanged**.
9. `docs/SELF_MOD_PUBLISH.md`:
   - the `contentHash`/`tarballSHA256` table rows, which `host/runbook` binds to the golden;
   - the prose naming the unpublished candidate (lines ~99–104), the confirmation phrase (~229) and the fresh-store name (`world-0.1.2.db`, ~169).
   - The 0.1.1 PUBLISHED paragraph stays as history.
10. `host/pkgproj/testdata/dryrun_v0.41.0.txt` and `dryrun_v0.54.11.txt`: **re-capture** as the real stdout of `env -u AILANG_REGISTRY_API_KEY <bin> publish --dry-run` in `packages/world-core`, as the file header requires. Use v0.41.0 and the released v0.54.11, sha256-checked (V9: only the two version lines and the digests change). Do not hand-edit them; that is what #565 did.
11. `host/pkgproj/dryrun_shapes_test.go:117`: the "short hash no dots" case's full and 17-hex `contentHash` literals.
12. `tools/attended/self_mod_publish.sh:3`: a comment only ("0.1.2 is the row-260 candidate"). The version is read from the golden.

Not changed (measured green as-is):
- `scripts/verify_ail.sh`;
- `host/verifygate/*`;
- `host/replay/testdata/*`: replay goldens are byte-identical, and `host/replay` passes on v0.41.0;
- `cmd/world-publish/fences_test.go`, `release_history_test.go`, `wiring_test.go`;
- `host/broker/*`, `host/runbook/*`;
- the sketches and se-tools.

**M2 — record the review of v0.60.4 (separate commit, after M1 is green):**
13. `scripts/testdata/ailang_release_observed.txt`: add `v0.60.4` with a one-line provenance comment.
14. `pins/ailang-releases.json` and 15. `host/pins/ailang-releases.json`, kept byte-identical (`TestEmbeddedManifestMatchesCanonical`): append one release row, `v0.60.4` / `AILANG v0.60.4` / commit `563c6c9`, holding only the `observed` role, with no assets.
16. `host/pins/manifest_test.go:44`: the restated observed list becomes `"v0.41.0 v0.33.0 v0.30.0 v0.60.4"`.
- Add only the observed file and `host/pins` reds twice (`TestObservedReleasesFileMatchesManifest`, `TestReleaseLiteralsExplainedByManifest`); add the JSON row and `TestCurrentManifestIsTodaysPins` reds. All four files together are green (V11 on re-serialized JSON; **V20 on the exact hand-appended bytes below**).
- **The M2 edit, verbatim (measured green in V20; copy, do not retype).** In `pins/ailang-releases.json`, the last release object (`"release": "v0.30.0"`) closes with `    }` followed by `  ]` and `}`. Replace that closing `    }` with `    },` and insert this block before `  ]`:
  ```
      {
        "release": "v0.60.4",
        "version": "AILANG v0.60.4",
        "roles": [
          {"role": "observed", "since": "2026-11-16", "evidence": "upstream releases/latest on CI PATH from 2026-11-13 (type-checks contract clauses, ailang#2283); reviewed by row 260 (w-ailang-v0604-contract-eq)"}
        ]
      }
  ```
  Then `cp pins/ailang-releases.json host/pins/ailang-releases.json`. `git diff --stat` for the two JSON files is `7 +++++++` each, 0 deletions (V20). Every other byte is unchanged.
- Optional, same edit: the header's "v0.30.0 is RETAINED … replay goldens are still v0.30.0-scoped" is stale. `host/replay/replay_test.go:25–30` and `ci.yml` record that the goldens were re-verified on v0.41.0 on 2026-09-22, and no CI lane installs v0.30.0 (V17c). Correct the sentence or leave it for a follow-up, but do not delete the `v0.30.0` entry: it is still the instrument's control token.

Why M2 belongs to this row: the notice exists to say "the primary `.ail` gate has switched
compiler — review it". This row is that review. The notice is non-fatal, so M2 cannot block M1's green.

Estimated diff: M1 ≈ +55/−30 over 12 files (13 with item 12), most of it identity and docs.
M2 ≈ +20/−1 over 4 files.

## 4. Acceptance criteria (runnable; repo root; never the `-dirty` PATH `ailang`)

```bash
export V604=/tmp/ailang-v0604-i393/ailang          # released v0.60.4, asset sha256 checked
export PIN=$HOME/.pinned-ailang/ailang             # v0.41.0
```
- **AC1** CI's lane, locally: `AILANG_BIN=$V604 WORLD_PKG_AILANG_BIN=$PIN env -u AILANG_REGISTRY_API_KEY ./scripts/verify_ail.sh` exits 0 and prints:
  - `✓ 16/16 required world/ identities verified across 24 module(s)`;
  - `✓ all 138 execargs…`, `✓ all 40 required named tests pass`, `✓ all 503 se-tools…`;
  - `✓ world package gate PASSED: 9/9`, `PUB011 ratchet passed: uncontracted_exports=3`.
  - Before M2 the `ℹ UNRECOGNISED RELEASE … 'v0.60.4'` notice prints; after M2 it is absent.
- **AC2** Pin lane: the same command with `AILANG_BIN=$PIN` exits 0 with the same counts.
- **AC3** Per-file, both CI compilers: for each of `world/{types,contracts,transitions,logepoch}.ail`, `$V604 ai-check -timeout 5s <f>` and `$PIN ai-check -timeout 5s <f>` give `check.passed=true`, with verified counts of exactly 6/5/3/2 and no other status.
- **AC4** `AILANG_BIN=$PIN python3 scripts/test_kernel_contracts.py` exits 0 (`world contract history controls PASSED`).
- **AC5** Go: `go build ./... && go vet ./cmd/world-publish/... ./host/pkgproj/... ./host/runbook/... && AILANG_BIN=$PIN WORLD_PKG_AILANG_BIN=$PIN env -u AILANG_REGISTRY_API_KEY go test -count=1 ./cmd/world-publish/... ./host/broker/... ./host/pkgproj/... ./host/runbook/... ./host/replay/... ./host/verifygate/... ./host/pins/...` prints 7 `ok`.
  - Known unrelated load flakes in `host/broker` appeared in one full run (`TestNotPublishedDoesNotReAuthorizeTheConsumedApproval`, `TestWorktreeSnapshotRefusesAConfiguredFilter`). Both pass `-count=3` in isolation at base and at the candidate, and a full re-run was `ok` (V12). Re-run once before treating a red as real.
- **AC6** Identity:
  - the golden's `contentHash` is `sha256:359a82fde07d8e9004786448da37ea3fb07cdc256682efc0fa57240df3c4a72c`. This value is executor-independent: V16 measured it unchanged under CHANGELOG date and text edits and under manifest edits. A different value means the module bytes differ from §1's text;
  - `interfaceHash` and `interfaceHashV2` are byte-identical to base;
  - `"version":"0.1.2"`, and `packages/world-core/ailang.toml` `[release] kind = "feature"`;
  - `(cd packages/world-core && env -u AILANG_REGISTRY_API_KEY $PIN publish --dry-run)` prints `Publishing world/core@0.1.2...`, `contracts: ✓ 16/16 verified`, `changelog_section true` and `✓ no gates`.
- **AC7** History is not rewritten: `git diff origin/dev -- host/broker/testdata host/replay/testdata` is empty, and the CHANGELOG `## 0.1.1`/`## 0.1.0` sections are byte-identical.
- **AC8** CI: all three jobs of `.github/workflows/ci.yml` green on the PR and on the `dev` merge — `ailang-code verify gate`, `go host build + test gate` (it also installs v0.54.15 for the floor grader, V19) and `go race gate (go test -race, own job)`; read the commit's full check set, not a name list. The verify job installs `releases/latest`, so it exercises v0.60.4 (or newer, in which case record the version it resolved).
- **AC9** The PR body says it supersedes #565 (credited for the two `deriving` lines) and that #565 was not touched.

## 5. Mutation and negative controls

Measured under v0.60.4 with `ai-check -timeout 5s` on the prototype (`ev/mutations-v0604.txt`). Each
mutant reds `verify_ail.sh` Leg 1 (`check.passed != true`, or a required identity not `verified`):

| id | mutant | result |
|---|---|---|
| MUT-1 | drop `deriving (Eq)` from `EvidenceGrade` | `TC_CONTRACT_001` @ types:43: KILLED |
| MUT-2 | drop `deriving (Eq)` from `TimeoutOutcome` | `TC_CONTRACT_001` @ types:174: KILLED |
| MUT-3 | `plan` ensures back to `result.inputWorld == w.stateRoot` | `TC_CONTRACT_001` @ transitions:29: KILLED |
| MUT-4 | `verify` ensures back to `result.proposalHash == p.proposalHash` | `TC_CONTRACT_001` @ transitions:48: KILLED |
| MUT-5 | `plan` body `inputWorld: w.logHead` | `plan` counterexample: KILLED |
| MUT-6 | `plan` body `proposalHash: transitionFn` | `plan` counterexample: KILLED |
| MUT-7 | `verify` rejected branch `accepted: true` | `verify` counterexample: KILLED |
| MUT-5a/5b/6a/7a | one-field-only body changes (`algo` or `digest` of `inputWorld`); `transitionFn: proposalHash`; `verify` accepted branch `proposalHash: p.inputWorld` | counterexample on v0.41.0 and v0.60.4, for both A and B: KILLED (V15) |
| MUT-8 | drop `result.inputWorld.algo == …` from `plan`'s ensures (weakening) | still `verified`: **SURVIVES**. The identity gate cannot see a weaker specification, and the reviewer must diff the clause. Accepted limitation; it predates this row. |
| MUT-9 | split `plan`'s ensures over several lines | `test_kernel_contracts.py` `PAR_UNEXPECTED_TOKEN`: KILLED by AC4 (V8) |
| MUT-10 | skip the reprojection (edit `world/` only) | step 3/9 projection mismatch: KILLED |
| MUT-11 | stale golden / doc digest / dry-run fixture | step 9/9 diff; `host/runbook`; `host/pkgproj` `TestParseDryRun*`: KILLED (V9) |
| MUT-12 | version left at 0.1.1 in any one of manifest, step 4/7/9, fences | step 4 or 7 or 9, or `release_history_test`: KILLED |

## 6. Out of scope

- Publishing 0.1.2. The attended hand-off is below, and the loop never publishes (D-WORLD-109).
- Touching, commenting on, or closing PR #565.
- Pinning CI legs 1–2 to v0.60.3 (D-WORLD-116 option C).
- Moving the frozen or package compiler off v0.41.0.
- Rewriting `test_kernel_contracts.py`'s `commit` exemption probe (line 75). It uses record `==` and passes because the script asserts v0.41.0. When the frozen compiler moves to ≥ v0.60.x, that probe will hit `TC_CONTRACT_001` and must be restated field-wise. Note it in the PR.

## 7. Attended release hand-off (Mark, after M1 merges)

The procedure is `docs/SELF_MOD_PUBLISH.md` Stage B, unchanged except for the version. Every value comes from the golden.

```bash
export WORLD_STORE="$HOME/.ailang/world/world-0.1.2.db"   # fresh store, as D-WORLD-45 = A did for 0.1.1
export WORLD_COMPILER="$HOME/.pinned-ailang/ailang"         # v0.41.0, the frozen compiler
./tools/attended/self_mod_publish.sh preflight
./tools/attended/self_mod_publish.sh approve      # type: publish world/core@0.1.2 irreversibly
./tools/attended/self_mod_publish.sh dry-run
./tools/attended/self_mod_publish.sh live         # prints the --live command for Mark to paste
```

Before `live`, compare the served-candidate digests against the four rows of the runbook table at
full length. After publishing, bank the served metadata as
`design_docs/verification/world-attended-<date>-publish-0.1.2/`, as was done for 0.1.1.

## 8. Verification Log

Binaries:
- v0.60.4 is `/tmp/ailang-v0604-i393/ailang` (`AILANG v0.60.4`, commit `563c6c9`, release asset sha256 checked by the controller).
- v0.41.0 is `$HOME/.pinned-ailang/ailang`; v0.30.0 is `$HOME/.pinned-ailang/ailang.v0.30.0`.
- v0.54.11 is `/tmp/i393-dl-v05411/bin/ailang`: downloaded by `gh release download v0.54.11`, and `shasum -a 256 -c` printed `OK`.
- Z3 4.16.0.

Scratch clones:
- `/tmp/i393-designer-base` (base `d0c9c4ce`);
- `/tmp/i393-designer-scratch` (prototype, branch `candA-012-1line`);
- `/tmp/i393-designer-candB` (alternative B).

The real worktree's code was not edited.

| V | command | output | file |
|---|---|---|---|
| V0 | controller: `AILANG_BIN=$V604 ./scripts/verify_ail.sh` at base | rc 1; `✗ world/contracts.ail: check.passed != true` after 10 sketches | `world-iter393/baseline-verify_ail-v0604.log` |
| V1 | `sweep.sh <base> $V604` (per-file `ai-check`, module-root cwd) | 8 fail `TC_CONTRACT_001` (all at types:43), 24 pass | `ev/sweep-baseline-v0604.txt` |
| V2 | probe modules `rec`/`enumeq`/`named`/`recplain` on three binaries | record-alias `deriving (Eq)`, enum `deriving (Eq)` and a nominal type all `verified` on v0.60.4/v0.41.0/v0.30.0; `==` on a record with no Eq is a type error in a body on all three | `ev/probe-deriving-3bin.txt` |
| V3 | per-file `ai-check` of `world/*`, base vs prototype, v0.41.0 and v0.30.0 | v0.41.0 16/16 both; v0.30.0 13 + `error`×4 + `skipped`×1, identical at base and prototype | `ev/aicheck-world-baseline-v0410-v0300.txt`, `ev/aicheck-final-3bin.txt` (30/30 `passed=True`) |
| V4 | enums fixed, `plan` restated; `$V604 ai-check world/transitions.ail` | new `TC_CONTRACT_001` at transitions:56 (`verify`, S4) | session output; MUT-4 reproduces it at :48 |
| V5 | `AILANG_BIN=$V604 WORLD_PKG_AILANG_BIN=$PIN ./scripts/verify_ail.sh` (S1–S4 fixed, golden not yet moved) | legs 1/2c/2/2b green (24 modules, 16/16, 138, 40, 503); Leg 3 steps 1–8 ✓; step 9 golden diff only | `ev/verify_ail-candA-v0604-pkg-v0410.log` |
| V6 | alternative B (`HashRef … deriving (Eq)`), three binaries plus package gate | 16/16 on v0.60.4/v0.41.0; `interfaceHashV2` unchanged; contentHash `2664ce1d…` | `ev/aicheck-candB-hashref-deriving-3bin.txt`, `ev/verify_world_package-candB-v0410.log` |
| V7 | final prototype: `verify_ail.sh` as AC1 | rc 0; `✓ verify gate PASSED: 16 … 40 … 503`; golden contentHash `359a82fd…`; interface v1/v2 unchanged | `ev/verify_ail-final-v0604-pkg-v0410.log`, `ev/prototype-candA-final.diff` |
| V8 | `AILANG_BIN=$PIN python3 scripts/test_kernel_contracts.py`: multi-line vs single-line ensures | multi-line: `AssertionError … PAR_UNEXPECTED_TOKEN at world/transitions.ail:30:3`; single-line: `world contract history controls PASSED` | `ev/test_kernel_contracts-candA-multiline-v0410-RED.log`, `ev/test_kernel_contracts-candA-1line-v0410.log` |
| V9 | `go test ./host/pkgproj/...` before fixture re-capture; re-capture with v0.41.0 and v0.54.11 | 3 FAIL (`TestParseDryRun*`, `TestVerifyWorldPackageDryRunCheck…`); after re-capture `ok`. The re-captured fixtures differ from base only in the 2 version lines and the digests | session output; `ev/gotest-final-v0410.log` |
| V10 | `publish --dry-run` (v0.41.0) with kind `fix`/`patch`/`feature` | `fix` and `feature`: `✓ no gates`; `patch`: `kind must be one of security\|fix\|feature\|breaking`; `fix` gives contentHash `359a82fd…` (unchanged), 10855 bytes | `ev/release-kind-probe-v0410.txt` |
| V11 | M2 probe: observed file, then + pins JSON, then + manifest_test | 2 FAIL, then 1 FAIL, then `ok host/pins`; `host/verifygate -run 'Observed\|Release\|Notice\|Expected'` `ok` | `ev/prototype-observed-v0604-m2.diff` |
| V12 | AC5 Go set on the prototype | 6 `ok` + `host/broker` FAIL on 2 load-flaky tests; both pass `-count=3` at base and candidate; full `host/broker` re-run `ok` (195.7 s); `host/pins` `ok` | `ev/gotest-final-v0410.log`, `ev/broker-flake-base-vs-cand.txt`, `ev/gotest-broker-rerun-v0410.log` |
| V13 | `WORLD_PKG_AILANG_BIN=$V604 ./scripts/verify_world_package.sh` (informational) | steps 1–8 ✓ on v0.60.4 (check, 9 named skips, smoke, dry-run, tar); step 9 refuses `wrong compiler version: AILANG v0.60.4` by design | `ev/verify_world_package-final-under-v0604-informational.log` |
| V14 | `gh pr view 565` / `gh pr diff 565` (read only) | OPEN, head `44e6c8a9`, 7 files: both `deriving` lines, golden, doc table, two dry-run fixtures **hand-edited**, shapes test; version stays **0.1.1**; no `transitions.ail` change | `world-iter393/controller/v14-pr565-excerpt.txt` (banked in r2: state, head, 7-file list, the four `deriving` hunks) |
| V15 | (a) `ai-check` of probe module `receq` (record `H` with `deriving (Eq)`) on v0.41.0 and v0.60.4. (b) `mut_body.py <bin>` (7 body mutants of `plan`/`verify` plus an unmutated control) in clone A (field-wise) and clone B (original `==` with `HashRef … deriving (Eq)`), on v0.41.0 and v0.60.4 | (a) both binaries: `eqByFields` and `fieldsByEq` `verified` (`==` ⇔ field conjunction, both directions); `digestDiffers` and `algoDiffers` give `counterexample`; control `same` `verified`. (b) all 4 candidate × binary pairs: 7/7 mutants give `counterexample` on the mutated function, other functions `verified`, control all `verified`: identical kill sets | `ev/v15-record-eq-probe.txt`, `ev/v15-receq.ail.txt` (source; `.txt` so the design_docs Leg-1 sweep does not enumerate it), `ev/v15-body-mutants-A-vs-B.txt` |
| V16 | Copy of the prototype package, module bytes fixed, `publish --dry-run` (v0.41.0) after each change: CHANGELOG date, CHANGELOG text appended, manifest `description`, restore, CTRL: one comment line appended to `world/types.ail`, restore. Plus full-length values via the released v0.54.11 dry-run | contentHash prefix `359a82fde07d8e900` in every non-control row; tarball 10852/10851/10876/10861/10852 bytes; CTRL contentHash `a34fc93ca3b47305f`, restored back to `359a82fd`; v0.54.11 full value `sha256:359a82fde07d8e9004786448da37ea3fb07cdc256682efc0fa57240df3c4a72c` before and after a CHANGELOG date and text change | `ev/v16-contenthash-invariance-v0410.txt` |
| V17 | `grep -n` at `d0c9c4ce`: (a) `ci.yml` install URLs; (b) immutability sources and the banked 0.1.1 record; (c) `v0.30.0` in `.github/workflows/*.yml` and `scripts/*` | (a) `ci.yml:43` `releases/latest/download` (verify job); every other install is `releases/download/v0.41.0` (`:85`, `:178`, `:480`) or `v0.54.15` (`:410`). (b) quotes in §2; served record `"version": "0.1.1"`, `"published_at": "2026-10-02T07:58:43Z"`; no 409 in the banked README. (c) no `releases/download/v0.30.0` in `ci.yml`; every `v0.30.0` hit in the workflows and scripts is a comment, apart from the `verify_ail.sh:91` release-verdict control token (a string, not a binary) | `ev/v17-premises-grep.txt` |
| V18 | consumer probe `import world/types (…); export func isTested(e: Evidence) -> bool { gradeOf(e) == TESTED }`, ai-check in clones base and A | base: v0.41.0 `No instance for Eq[EvidenceGrade]`; v0.60.4 fails earlier on S1. Candidate: `passed=True` on both. Interface hashes are identical across the two states (V7) | `ev/v18-consumer-eq-probe.txt` |
| V19 | controller: `grep -n -B8 -A14 'v0.54.15' .github/workflows/ci.yml`; `grep` of `scripts/floor/*.py`; `~/.pinned-ailang-tools/v0.54.15/ailang ai-check -timeout 5s world/{types,contracts,transitions,logepoch}.ail` in the base clone and the candidate clone | v0.54.15 is the floor grader's `$TOOL` (step "Install PINNED v0.54.15 ailang tool binary (floor grader $TOOL)", job `go host build + test gate`; `scripts/floor/grade.py:50`), which grades banked benchmark solutions, not `world/`. Measured anyway: `passed=True`, verified 6/5/3/2 at base and at candidate | `world-iter393/controller/v19-v05415-kernel-aicheck.txt` |
| V20 | controller: fresh worktree at `d0c9c4ce`; the exact §3 M2 hand-append, `cp` to `host/pins`, observed file + `manifest_test.go:44`; `go test -count=1 ./host/pins/...`; `AILANG_BIN=$PIN go test -count=1 -run 'Observed\|Release\|Notice\|Expected' ./host/verifygate/...`; `AILANG_BIN=$V604 verify_ail.sh \| head` | diffstat 16+/1−, the JSON files additions only; `ok host/pins`; `ok host/verifygate` (78.9 s; the same subset without `AILANG_BIN` FAILs by design, "AILANG_BIN is unset — … Never skip"); the `UNRECOGNISED RELEASE` notice is absent | `world-iter393/controller/v20-m2-handappend.txt` |
| V21 | controller: candidate `world/` copy with `plan`'s ensures as `sameRef(…) && sameRef(…) && sameRef(…) && result.goal == goal`; `ai-check -timeout 5s world/transitions.ail` on v0.60.4 and v0.41.0 | v0.60.4: `passed` true, 3 `verified`; v0.41.0: `passed` true, `plan` **`skipped`** ("SMT declaration closure could not resolve a required sort (cannot encode ensures clause …)"), the other two `verified` | `world-iter393/controller/v21-sameref-in-ensures.txt` |

`sweep-candA-v0604.txt`, `sweep-candA-v0.41.0.txt` and `sweep-candA-v0.30.0.txt` record the
intermediate multi-line prototype. It has the same compile and verify verdicts, and was replaced
because of V8.

The probe scripts are banked under `ev/tools/` (`sweep.sh`, `mut.py` for MUT-1..8, `mut_body.py` for V15).

## 9. Quorum log

**r1** (`design_docs/verification/world-iter393/quorum/r1.json`, 2026-11-16T05:57:48Z). Synthesis
was BLOCKED; this is the one mandated revision. No seat disputed the direction (field-wise
restatement, enum `deriving`, 0.1.2 candidate).

| seat | verdict | objection surface |
|---|---|---|
| gemini-4-argon | PASS | The release-kind fork (`fix` vs `feature`) left M1's file list and step 4's frozen `want` undetermined. M2's JSON edit must not reformat the file. |
| oc-glm-5-3 | REJECT | (1) "Field equality ⇔ record equality" was asserted, not measured, and MUT-5..7 alone cannot show equivalence. (2) contentHash invariance under CHANGELOG bytes was inferred. (3) Registry immutability and "no CI lane uses v0.30.0" had no V-row, and the ci.yml premise had no banked citation. |
| oc-kimi-k3 | REJECT | (1) Inconsistency: S3 rejected B for adding an `Eq` instance while S1/S2 add two, and the `fix` recommendation rested on "adds no capability". (2) `fix` was never run end to end. (3) The same premise citations as glm. |

**What r1 changed (measurements first, prose after):**
- V15 measures the equivalence (record `==` ⇔ field conjunction, both directions, on both CI compilers), and A and B kill the same 7 body mutants on both compilers. §1's claim now cites V15.
- V16 measures contentHash invariance under CHANGELOG date, CHANGELOG text and manifest edits, with a module-byte positive control. AC6 keeps the exact value and cites V16.
- The release kind is decided: `feature`. V18 shows the instance addition is real and invisible to both interface hashes. §2 states the blind spot and why it is acceptable. §1 restates B's rejection on HashRef-specific grounds (`sameRef` duplication, logepoch churn), not "instances are bad". Files items 4–6 are unconditional. AC6 asserts the kind. No kind=`fix` run was needed, since every end-to-end measurement (V5, V7) used `feature`.
- V17 banks `ci.yml:43` (`releases/latest`), the first-party immutability sources (and says plainly that no observed 409 is banked), and the v0.30.0 grep. The M2 JSON edit is now append-only by hand.

**r2** (`quorum/r2.json`). Synthesis BLOCKED again: gemini-4-argon PASS, oc-glm-5-3 REJECT,
oc-kimi-k3 REJECT. No absent seat. The objections moved to new, narrow surfaces (rule (a) from
round 3 on is recorded here early):

| seat | verdict | objection surface |
|---|---|---|
| gemini-4-argon | PASS | — |
| oc-glm-5-3 | REJECT | V17a's fifth install site (v0.54.15, `ci.yml:410`) undisposed; AC8's "both jobs" unenumerated; the `sameRef` reuse question; V14 not banked |
| oc-kimi-k3 | REJECT | M2's hand-append bytes never measured (V11 re-serialized); S2 line 173 vs MUT-2 174 |

Every remaining objection carries a concrete reviewer-authored fix and none disputes the
direction, so the controller applied the **narrow-refinement carve-out** (Gate 2) instead of a
third designer round: r2 is the reviewers' fixes, each backed by a controller measurement —
V19 (the v0.54.15 lane disposed: floor grader only, and 16/16 at base and candidate anyway; the
§1 table row; AC8 names the three jobs), V20 (the exact hand-append measured green; the snippet is
now verbatim in §3), V21 (`sameRef` in `ensures` is `skipped` on v0.41.0, so it cannot replace
the field-wise form), V14 banked, and the S2 line-number note in §1. No re-quorum follows; the doc
routes to the planner, and the routing evidence records the carve-out.
