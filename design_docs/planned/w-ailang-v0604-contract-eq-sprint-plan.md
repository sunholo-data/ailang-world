# Sprint plan — w-ailang-v0604-contract-eq (charter row 260, iteration 393)

**Planner-Lane:** opus (unattended). **Design:** [w-ailang-v0604-contract-eq.md](w-ailang-v0604-contract-eq.md)
r2 (final; §9 carve-out). **Decision:** `D-WORLD-116` OPEN, default **A** applied: World ships its own
fix-forward on a fresh branch, never touches PR #565, never publishes.
**Base:** `cbcf80fb` on `sprint/world-iter393-row260` (origin/dev `842b9331` + record 392 + the two
design commits). **Plan artifacts:** `design_docs/verification/world-iter393/plan/` (`M1.patch`,
`M2.patch`, `gates.sh`, `gates-planner-run.txt`, `mutants-planner.txt`).

**Goal.** CI's `ailang-code verify gate` installs `releases/latest` = AILANG v0.60.4, which
type-checks contract clauses (`TC_CONTRACT_001`). Make the kernel's four failing `ensures` (S1–S4)
type-check on v0.60.4 while staying 16/16 `verified` on v0.41.0 (frozen/package compiler), move
the package candidate to `world/core@0.1.2` (0.1.1 is published and immutable), and record the
review of v0.60.4 in the release manifest.

**Binaries** (never the `-dirty` PATH `ailang`):

```bash
export V604=/tmp/ailang-v0604-i393/ailang            # AILANG v0.60.4, commit 563c6c9
export PIN=$HOME/.pinned-ailang/ailang               # AILANG v0.41.0, commit 24ee108
export V05411=/tmp/i393-dl-v05411/bin/ailang         # AILANG v0.54.11 (shasum -c OK, re-checked by the planner)
```

## 1. How the executor applies this sprint

The patches were produced by `git diff` in the planner's scratch tree and proven by `git apply`
into a second fresh worktree at `cbcf80fb`, where every gate below ran green (§5).

```bash
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter393-row260
P=$W/design_docs/verification/world-iter393/plan
cd $W && git apply --check $P/M1.patch && git apply $P/M1.patch
# review, run M1 acceptance (§2), commit M1 (controller)
git apply --check $P/M2.patch && git apply $P/M2.patch
# review, run M2 acceptance (§3) + full gates, commit M2 (controller)
```

The executor must still **reproduce** the generated parts rather than trust the patch (§2 "regenerate
check"): the projection, the golden and both dry-run fixtures are tool output, and re-running the
tool on the patched tree must leave `git status` unchanged.

## 2. M1 — kernel fix + projection + 0.1.2 candidate + docs (ONE atomic commit)

Any subset of these files reds a gate (MUT-10/11/12), so they land together. Commit body credits
PR #565 for the two `deriving (Eq)` lines and says #565 was not touched.

| # | file | change | numstat |
|---|---|---|---|
| 1 | `world/types.ail` | `deriving (Eq)` on `EvidenceGrade` and `TimeoutOutcome` | +2 −0 |
| 2 | `world/transitions.ail` | `plan`/`verify` `ensures` restated field-wise, **one line each**, byte-equal to design §1 | +2 −2 |
| 3 | `packages/world-core/world/types.ail` | `./scripts/build_world_package.sh` output | +2 −0 |
| 3 | `packages/world-core/world/transitions.ail` | same | +2 −2 |
| 4 | `packages/world-core/ailang.toml` | `version = "0.1.2"`; `[release]` comment line "0.1.2 keeps `feature` (row 260: …)"; `kind` stays `feature` | +2 −1 |
| 5 | `packages/world-core/CHANGELOG.md` | new `## 0.1.2 — 2026-11-16` section: names `Eq[EvidenceGrade]`, `Eq[TimeoutOutcome]` as **added** instances, the field-wise restatement, credits PR #565; 0.1.1/0.1.0 byte-identical | +17 −0 |
| 6 | `scripts/verify_world_package.sh` | step 4 `want` version, step 7 regex `0\.1\.2`, step 9 packet version, comment at :128–130 | +5 −4 |
| 7 | `scripts/world_package_ready_packet.golden.json` | the step-9 diff's `+` line (generated) | +1 −1 |
| 8 | `cmd/world-publish/fences.go` | `attendedPhrase`, `frozenPackageVersion = "0.1.2"`, version comment; `frozenInterfaceHashV2` comment (value **unchanged**) | +9 −5 |
| 9 | `docs/SELF_MOD_PUBLISH.md` | contentHash/tarballSHA256 table rows; candidate prose (~99–106); store name `world-0.1.2.db`; phrase/scope prose (~226–229); line 8 consumes `0.1.2`; 0.1.1 PUBLISHED paragraph kept as history | +19 −15 |
| 10 | `host/pkgproj/testdata/dryrun_v0.41.0.txt` | re-captured (real stdout, below) | +4 −4 |
| 10 | `host/pkgproj/testdata/dryrun_v0.54.11.txt` | re-captured (real stdout, below) | +5 −5 |
| 11 | `host/pkgproj/dryrun_shapes_test.go:117` | "short hash no dots" literals → `359a82fde07d…` / 17-hex `359a82fde07d8e900` | +1 −1 |
| 12 | `tools/attended/self_mod_publish.sh:2–4` | comment only: 0.1.1 published 2026-10-02, 0.1.2 is the row-260 candidate | +3 −2 |

**M1 total: 14 files, +74 −42** (`git diff --numstat` on the patched tree). The design estimated
+55/−30 over 12–13 files; the excess is docs prose (item 9) and the `frozenInterfaceHashV2` comment.

**Generated parts and their generating commands** (all real tool output; re-run and expect no diff):

```bash
./scripts/build_world_package.sh                         # items 3 (prints "projected 4 modules")
(cd packages/world-core && env -u AILANG_REGISTRY_API_KEY $PIN    publish --dry-run) > host/pkgproj/testdata/dryrun_v0.41.0.txt
(cd packages/world-core && env -u AILANG_REGISTRY_API_KEY $V05411 publish --dry-run) > host/pkgproj/testdata/dryrun_v0.54.11.txt
AILANG_BIN=$PIN WORLD_PKG_AILANG_BIN=$PIN env -u AILANG_REGISTRY_API_KEY ./scripts/verify_world_package.sh   # step 9 must say "equals committed golden"
```

The v0.41.0 dry-run was run twice and was byte-identical (deterministic, including the `(0.1s)`
timing field). Both fixtures differ from base only in the version lines and the tarball/content
digests, as design V9 predicted. Measured identity (planner, v0.41.0 and v0.54.11 agree):

| digest | base 0.1.1 | M1 0.1.2 |
|---|---|---|
| contentHash | `sha256:47351707…a2780` | `sha256:359a82fde07d8e9004786448da37ea3fb07cdc256682efc0fa57240df3c4a72c` (= design AC6) |
| interfaceHash / V2 | `d16cc882…4083` / `ifacev2:b25fe031…e8d8` | unchanged |
| tarballSHA256 | `07d61809…f4a3` | `sha256:be02038ea1e5dd77d73a686a67be6440f945abab1d7a84303a6cea41114d5004` |
| tarballBytes | 10610 | 11030 |

**M1 acceptance** (all ran on the planner's proof tree with only M1 applied; `bash gates.sh <tree>
<logs> m1` with `GATES_SKIP_VERIFY_GO=1`, 452 s, `FAILS: none`):

| AC | command (repo root) | result on proof tree (M1 only) |
|---|---|---|
| AC1 | `AILANG_BIN=$V604 WORLD_PKG_AILANG_BIN=$PIN env -u AILANG_REGISTRY_API_KEY ./scripts/verify_ail.sh` | rc 0; 16/16 across 24 modules, 138, 40, 503, 9/9, PUB011=3; `UNRECOGNISED RELEASE … 'v0.60.4'` **prints** (expected before M2) |
| AC2 | same with `AILANG_BIN=$PIN` | rc 0, same counts, no notice |
| AC3 | `$V604`/`$PIN ai-check -timeout 5s world/{types,contracts,transitions,logepoch}.ail` | 8/8: `passed=True`, verified 6/5/3/2, only status `verified` |
| AC4 | `AILANG_BIN=$PIN python3 scripts/test_kernel_contracts.py` | rc 0, `world contract history controls PASSED` |
| AC5 | `go build ./...`; `go vet ./cmd/world-publish/... ./host/pkgproj/... ./host/runbook/...`; `go test -count=1` × 7 packages with `AILANG_BIN=$PIN WORLD_PKG_AILANG_BIN=$PIN` | rc 0 / rc 0 / 7/7 `ok` (broker 199 s, verifygate 156 s; no re-run needed) |
| AC6 | golden contentHash = design value; interface v1/v2 = `git show cbcf80fb:` golden; version 0.1.2; kind `feature`; `$PIN publish --dry-run` lines | 13/13 OK |
| AC7 | `git diff cbcf80fb -- host/broker/testdata host/replay/testdata` empty; CHANGELOG from `## 0.1.1` to EOF byte-equal to base; newest section `## 0.1.2 — ` | OK |
| fmt | `gofmt -l` on the 2 changed Go files | empty |

## 3. M2 — record the v0.60.4 review (separate commit, after M1 is green)

| # | file | change | numstat |
|---|---|---|---|
| 13 | `scripts/testdata/ailang_release_observed.txt` | 3-line provenance comment (release, commit `563c6c9`, ailang#2283, row 260) + `v0.60.4` entry | +4 −0 |
| 14 | `pins/ailang-releases.json` | the design §3 verbatim block hand-appended after the `v0.30.0` object (`    }` → `    },`); no re-serialization | +7 −0 |
| 15 | `host/pins/ailang-releases.json` | `cp pins/ailang-releases.json host/pins/` (byte-identical) | +7 −0 |
| 16 | `host/pins/manifest_test.go:44` | observed list `"v0.41.0 v0.33.0 v0.30.0 v0.60.4"` | +1 −1 |

**M2 total: 4 files, +19 −1.** Cumulative M1+M2: 18 files, +93 −43.

**M2 acceptance:** `go test -count=1 ./host/pins/...` `ok`; `cmp pins/ailang-releases.json
host/pins/ailang-releases.json`; AC1 again with the notice **absent**; then the full `gates.sh … m2`
(§5), which includes `./scripts/verify_go.sh` (plain leg = CI's "go build + test gate").
The optional stale-header fix (design §3 "v0.30.0 is RETAINED … v0.30.0-scoped") is **not** in
M2; it is left for a follow-up so M2 stays append-only.

## 4. Mutation table (anchored on the diff)

MUT-1..12 are the design's (measured by the designer on the prototype; the kernel bytes in M1 are
identical to the prototype's, so they carry over). MUT-13..21 are the planner's, measured on the
patched scratch tree (`plan/mutants-planner.txt`; each mutant applied, its gate run, restored).

| id | mutant (diff line it attacks) | killing gate | result |
|---|---|---|---|
| MUT-1/2 | drop either `deriving (Eq)` (item 1) | AC1/AC3 v0.60.4 `TC_CONTRACT_001` | KILLED (design) |
| MUT-3/4 | `plan`/`verify` ensures back to record `==` (item 2) | AC1/AC3 v0.60.4 | KILLED (design) |
| MUT-5..7, 5a/5b/6a/7a | body mutants of `plan`/`verify` | AC3 counterexample, both compilers | KILLED (design V15) |
| MUT-8 | drop one conjunct of `plan`'s ensures | — | **SURVIVES** (accepted, predates row; reviewer diffs the clause against design §1 byte-for-byte) |
| MUT-9 | split an ensures over two lines | AC4 `PAR_UNEXPECTED_TOKEN` | KILLED (design V8) |
| MUT-10 | skip reprojection (item 3) | step 3/9 | KILLED (design) |
| MUT-11 | stale golden / doc digest / dry-run fixture (items 7, 9 table, 10) | step 9; `host/runbook`; `host/pkgproj` | KILLED (design V9) |
| MUT-12 | 0.1.1 left in any ONE of manifest, step 4/7/9, fences (items 4, 6, 8) | step 4/7/9, `release_history_test` | KILLED (design) |
| MUT-13/13b | one byte of CHANGELOG prose / of the toml comment changed after the golden was regenerated (items 4, 5) | step 9 `✗ ready packet differs byte-for-byte from golden` | KILLED ×2 |
| MUT-14 | `[release] kind = "fix"` (item 4) | step 4 `✗ manifest is not the exact frozen structure` | KILLED |
| MUT-15 | shapes-test literal left at the 0.1.1 contentHash (item 11) | `TestParseDryRunLeadingWhitespaceAndRefusals` `control: case "short hash no dots" did not mutate the fixture` | KILLED |
| MUT-16 | `attendedPhrase` left at `0.1.1` (item 8, phrase only) | `TestAttendedRehearsalNamesTheCandidate` `does not contain "publish world/core@0.1.2 irreversibly"` | KILLED |
| MUT-17/17b | `## 0.1.2` header removed, so `## 0.1.1` is newest (item 5) | `TestChangelogCarriesTheCandidate…` `newest CHANGELOG section is "0.1.1"`; step 7 dry-run `✗ PUB001 … no non-empty ## 0.1.2 section` | KILLED ×2 |
| MUT-18 | M2: the two JSON copies diverge (item 15) | `TestEmbeddedManifestMatchesCanonical` | KILLED |
| MUT-19 | M2: `manifest_test.go` not updated (item 16) | `TestCurrentManifestIsTodaysPins` `observed releases [… v0.60.4]` | KILLED |
| MUT-20 | M2: JSON row added, observed file not (item 13) | `TestObservedReleasesFileMatchesManifest` | KILLED |
| MUT-21 | the whole version left at 0.1.1 **consistently** (manifest, step 4/7/9, fences, CHANGELOG header, golden) — PR #565's shape | — | **SURVIVES** by construction: no test refuses the *published* 0.1.1 (only 0.1.0 via the preserved record). Caught today only by AC6's `version 0.1.2` assert in `gates.sh`. Follow-up row: bank the 0.1.1 served record and extend the "already published" check. |

**Named untested lines** (no gate reads them; reviewed by eye): the comments in `fences.go`
(version history + `frozenInterfaceHashV2`), `verify_world_package.sh:128–130`,
`self_mod_publish.sh:2–4`, the observed-file comment lines, and the `SELF_MOD_PUBLISH.md` prose
outside the digest table (lines 8, 15–16, 99–106, 156–158, 172, 226–229). The CHANGELOG and toml
prose are **byte**-bound by the golden (MUT-13) but their *content* is reviewed, not tested.

MUT-13..20 were run by `/tmp/i393-planner-tools/mutants.sh` (banked inside
`plan/mutants-planner.txt`): 10/10 KILLED, each with a nonzero rc **and** the named message; the tree's
`git diff cbcf80fb --binary | shasum` was `f06545bf…` before and after (restored). The unmutated
control for each is the green gate run in §7. MUT-21 was not run: it is argued from code
(`release_history_test.go:67` and `host/broker/release_candidate_test.go:23` refuse only `0.1.0`).

## 5. Gate script

`design_docs/verification/world-iter393/plan/gates.sh <tree> <logdir> <m1|m2>` — bash,
`set -uo pipefail`, no `set -e`; each leg's rc is the command's own (`> log 2>&1; rc=$?`, never
through a pipe); one `LEG <name> rc=<n> <secs>s` line per leg to `<logdir>/legs.txt`; exits 1 if any
leg failed. It refuses to start unless `$V604` reports `AILANG v0.60.4` and `$PIN` reports
`AILANG v0.41.0`. Legs: AC1 (notice expected in `m1`, absent in `m2`), AC2, AC3 (8 per-file
checks), AC4, gofmt on changed Go files, AC5 build / vet / 7 packages (`host/broker` re-run once on
failure, announced in `legs.txt` with the failing test names — design AC5), AC6 (13 asserts; base
interface hashes from `git show cbcf80fb:…golden.json`), AC7, `M2_pins_row` (m2 only), and
`./scripts/verify_go.sh` with `VERIFY_GO_LEGS=plain` (the repo's full Go gate, as iteration 392's
`gates.sh` ran it; skippable with `GATES_SKIP_VERIFY_GO=1`).

Planner runs on the proof tree, banked as `plan/gates-planner-run.txt`: see §7.

## 6. Risks

- **R1 `host/broker` load flakes** (design V12). Re-run once; a second red is real. `verify_go.sh`
  has no re-run, so a flake there needs a manual re-run of the leg, recorded as such.
- **R2 `releases/latest` moves past v0.60.4** before CI runs. AC8 then exercises the newer release;
  record the resolved version. A newer release outside the observed set prints the non-fatal notice
  again; that is a new review, not a red.
- **R3 regenerated bytes differ on the executor's machine.** contentHash is module-bytes-only and
  must equal design AC6 exactly. tarballSHA256/bytes depend on CHANGELOG/toml bytes; if the
  executor edits any prose, it must regenerate the golden, both fixtures and the doc table together
  (MUT-11/13), never hand-edit.
- **R4 CHANGELOG is public registry text** once published (Mark's attended step). It names PR #565;
  the controller may reword the credit, but any byte change re-triggers R3.
- **R5 MUT-21 / MUT-8 survive** (accepted; §4). The reviewer must diff the two `ensures` lines
  against design §1 and confirm `0.1.2` everywhere.
- **R6 stale text outside scope**: `host/broker/release_candidate_test.go:11` comment ("names the
  world/core@0.1.1 CANDIDATE"), `website/docs/guides/{attended-steps,troubleshooting}.md` (0.1.1
  phrase/store as history) and the observed-file header ("Two entries", "v0.30.0-scoped"). Design
  §3 lists `host/broker/*` as unchanged; left for a follow-up and named in the PR.
- **R7** `test_kernel_contracts.py:75`'s `commit` probe uses record `==`; it passes only because the
  script asserts v0.41.0 (design §6). Note in the PR.

## 7. Planner evidence and premises found false

**Runs (planner, this iteration).** Scratch `/tmp/i393-planner` (detached `cbcf80fb`): the
designer's prototype diff applied, then design items 4/6/8/9/12 prose added, CHANGELOG rewritten,
projection re-run (no change), both fixtures re-captured, golden regenerated from step 9, doc table
tarball row updated; `verify_world_package.sh` rc 0 (9/9, PUB011=3). `M1.patch` = `git diff --cached`
there; M2 edited on top, `M2.patch` = `git diff`. Proof: fresh worktree `/tmp/i393-planner-proof`
at `cbcf80fb`; `git apply --check M1.patch` OK, applied; gates run 1 (`m1`); `git apply --check
M2.patch` OK, applied; gates run 2 (`m2`). Both trees' `git diff cbcf80fb --binary` hash equal
(`f06545bf…`). The four kernel/projection `.ail` files are byte-identical to the designer's
prototype, so MUT-1..12 carry over; both `ensures` lines grep byte-equal to design §1.

| run | legs | result |
|---|---|---|
| 1 `m1` (M1 only, `GATES_SKIP_VERIFY_GO=1`) | AC1 (notice=yes) 11 s, AC2 25 s, AC3 8/8, AC4, gofmt (2 files), build, vet, 7 pkgs 7/7 (409 s), AC6 13/13, AC7 | **all rc 0**, 452 s |
| 2 `m2` (M1+M2, full) | AC1 (notice=no), AC2, AC3 8/8, AC4, gofmt (3 files), build, vet, 7 pkgs 7/7 (407 s; broker 201 s, no re-run), AC6 13/13, AC7, M2_pins_row | **all rc 0** |
| 2 `m2` | `verify_go_plain` | **rc 1** at 346 s: `scripts/marketingcheck` `TestMarketingSupervisorFakeDaemonCleanup/shell-hup` ("hangup trap must exit 129 … signal: killed", "cleanup exceeded bound: 28.99s"). Untouched by this diff. Isolated `-run` → PASS (5.3 s). Manual full re-run of the leg → **rc 0**, 345 s, 62 `ok`, 0 FAIL, `✓ go gate PASSED (plain)` |

All logs: `plan/gates-planner-run.txt`. Mutants: `plan/mutants-planner.txt`.

**Premises found false / corrected:**
1. **Design §3 item 14 vs its own verbatim snippet.** Item 14 says the v0.60.4 row carries
   "commit `563c6c9`"; the verbatim M2 block (V20, "copy, do not retype") has no `commit` key, and
   neither do the other observed-only rows (v0.33.0, v0.30.0). Minimal fix: the verbatim block
   (measured green) ships; the commit is recorded in the observed-file provenance comment.
2. **Design AC5 names `host/broker` as the only known load flake.** The full Go gate
   (`verify_go.sh`, CI's "go build + test gate") also flaked once in `scripts/marketingcheck`
   (above), a package this diff does not touch. Not a direction change: the executor re-runs a
   red `verify_go_plain` leg once, records both runs, and treats a second red as real.
3. **Design §3 estimate** (M1 ≈ +55/−30, 12–13 files) is low: M1 is +74/−42 over 14 files because
   item 9's prose and the `frozenInterfaceHashV2` comment are real lines. Not a direction change.

**Executor notes.** Never hand-edit the golden, the fixtures, the doc digest rows or the projection;
if any CHANGELOG/toml byte changes, regenerate all four (R3). Run `gates.sh` with `V604`/`PIN`
pointing at the binaries above; the script refuses other versions. Use mode `m1` after M1 and `m2`
after M2. The commit bodies and the PR must carry design AC9 (supersedes #565, credited, untouched)
and the R6/R7 notes, plus the upstream observation that interface hashes do not record derived
instances (design §2).
