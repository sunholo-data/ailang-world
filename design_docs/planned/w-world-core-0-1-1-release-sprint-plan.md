# Sprint plan — w-world-core-0-1-1-release (iteration 211)

**Planner-Lane:** opus-required
**Authority:** `D-WORLD-41` RESOLVED A (attended, Mark, 2026-09-28): release the proof-hardened
package as `world/core@0.1.1` (proof hardening, no exported signature change), keep
`[release] kind = "feature"`; the loop prepares and rehearses the candidate, **a human runs the
irreversible publish**.
**Design:** [w-kernel-exports-carry-no-contract.md](w-kernel-exports-carry-no-contract.md) §4 and
§5 row M6 ("attended release preparation, separate ≤150 LOC slices after decision … M6 is not
permission to publish"). M1–M5 landed as `4b5b77b`.
**Base:** `origin/dev` `52447c3`, worktree `.wt-world-iter211`, branch
`sprint/w-world-core-0-1-1-release`. The planner made no git writes and published nothing.
**Prototype:** the whole change is prototyped, green and mutation-tested in the worktree;
diff saved at `~/.ailang/state/mission-world-iter211-evidence/planner_prototype.diff`
(+233 / −21 lines, 11 files). **Recommendation: the executor starts FROM the prototype**, reviews
it against this plan, and re-runs every acceptance command; it is not a skeleton to rewrite.

Command prefix for everything below (the pin is mandatory; the key is stripped so no registry
write can authenticate):

```sh
export AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_PKG_AILANG_BIN=$HOME/.pinned-ailang/ailang
export TMPDIR=$HOME/.ailang/state/mission-world-iter211-evidence/tmp   # outside the repo
alias E='env -u AILANG_REGISTRY_API_KEY'
```

## 1. Verification Log (planner, measured this iteration)

| # | Command | Observed |
|---|---|---|
| V1 | `curl -o /dev/null -w '%{http_code}' https://storage.googleapis.com/ailang-registry/packages/world/core/{0.1.0,0.1.1,0.2.0}/metadata.json` (read-only GET, key stripped) | `200` / `404` / `404`; 0.1.1 body is GCS `NoSuchKey` XML. Positive control = the 0.1.0 `200` in the same call. |
| V2 | GCS object listing `?prefix=packages/world/core/` | exactly two objects, both `0.1.0/` (`metadata.json` 8143 B, `package.tar.gz` 9785 B, 2026-09-21T19:39:38Z). **0.1.1 is unused.** |
| V3 | `shasum -a 256` of served 0.1.0 metadata vs `host/broker/testdata/metadata_world_core_0.1.0.json` | identical `3bc30977…59fd` — the preserved fixture is still byte-equal to what the registry serves. |
| V4 | `git grep -n '0\.1\.0' HEAD -- . ':!design_docs'` | **60 hits / 32 files** (+86 in `design_docs/`, historical, untouched). Classified in §2. `git grep '0\.1\.1'` = 13 hits (one must move, §2). The dry-run regex `0\.1\.0` (escaped) is invisible to the literal grep; found by reading `verify_world_package.sh:214`. |
| V5 | Baseline at `52447c3`: `E go vet ./...`; `E go test -race -count=1 -v ./cmd/world-publish/... ./host/broker/... ./host/pkgproj/... ./host/runbook/...`; `E ./scripts/verify_world_package.sh`; `world-publish packet` | vet rc=0; tests rc=0, **445** `=== RUN`, 4× `ok`; package gate rc=0 (9/9, PUB011=3); packet rc=0 names `world/core@0.1.0`, tarball `d532d923…`, 10127 B. |
| V6 | Probe: copy of package with only `version = "0.1.1"` → `ailang publish --dry-run` and `pkg quality --json` | dry-run **rc=1** `✗ PUB001 CHANGELOG.md has no non-empty '## 0.1.1' section`; quality rc=2, `gates=[PUB001]`. **The pinned compiler's own dry-run (gate step 7) enforces PUB001** — no new check needed for it. |
| V7 | Prototype package gate before golden refresh | rc=1 at step 9; diff shows only `tarballBytes 10127→10610`, `tarballSHA256 d532d923…→07d61809…`, `version 0.1.0→0.1.1`. **contentHash (`47351707…`), interfaceHash (`d16cc882…`) and interfaceHashV2 (`b25fe031…`) unchanged** — contentHash covers the `.ail` modules only. |
| V8 | Golden regenerated **with the repo's own generator**: the `+{…}` line that `verify_world_package.sh` step 9 prints (`diff -u GOLDEN tmp_ready`), written verbatim; gate re-run | rc=0: `✓ world package gate PASSED: 9/9 steps`, `PUB011 ratchet passed: uncontracted_exports=3`. |
| V9 | Prototype: `E go vet ./...`; `E go test -race -count=1 -v` over the four packages + `./host/verifygate/...` | vet rc=0; rc=0, **589** `=== RUN` with verifygate (**449** over the V5 four-package set = 445 + 4 new tests), 5× `ok`. |
| V10 | Prototype: `E ./scripts/verify_ail.sh` (PATH carries the pin) | rc=0: `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass`. |
| V11 | Prototype: `E go test -count=1 ./...` | run 1 rc=1: `TestQueryInterfaceReturnsWhileADescendantHoldsStdout` `descendant pid not recorded` (3 s budget under full parallel load; version-independent); isolated 3/3 `ok`; run 2 rc=0, 24× `ok`. **Pre-existing load flake, not caused by this change** (same class as the parked row-114 timing budget). |
| V12 | M1-only boundary (two new test files moved out): `E go vet ./...`; `E go test -count=1` over the five packages | vet rc=0; all 5 `ok`. M1 is a green commit boundary on its own. |
| V13 | Headless rehearsal, prototype binary, fresh temp store (`rehearsal/headless.log`) | `packet` rc=0 `ready packet for world/core@0.1.1`, all fields EQUAL; `approve` prints scope `…/version:0.1.1#…tarball=sha256:07d61809…` then `STOP fence=tty reason=no-controlling-terminal` rc=3; `publish --dry-run` `STOP fence=tty` rc=3; `reconcile` rc=0, 0 indeterminate. |
| V14 | `world-publish reconcile --probe` (read-only GETs) + `ailang publish --dry-run` in `packages/world-core` | `state=not-published package=world/core@0.1.1 absent=3 … each with a firing same-pass control`; dry-run rc=0 `Publishing world/core@0.1.1...`, `release: kind feature changelog_section true`, `16/16 verified`, `✓ no gates`, `Tarball: 10610 bytes`, `Dry run complete. Tarball ready but not uploaded.` |
| V15 | In-process attended rehearsal (`TestAttendedRehearsalNamesTheCandidate`, fence stack satisfied by the injected probe) | rc=0: `REHEARSAL — no request of any kind was made.`, `package world/core@0.1.1`, approval scope `…/version:0.1.1#…`, tarball `sha256:07d61809…`. |
| V16 | Copy of the operator's real store `~/.ailang/world/world.db` (sha256 `19122b31…` before and after — original never opened) → `world-publish reconcile --store <copy>` | **rc=3 `STOP fence=store reason=unopenable … has user_version 2 with application schema; binary requires 4; refusing to modify`.** No migration exists in `host/store`. See §6 finding F2. |
| V17 | `ls -la ~/.ailang/world/` | `.last_approval_ref` (2026-09-21) exists: a CONSUMED 0.1.0-scoped ref. |

## 2. Inventory — where 0.1.0 is frozen (60 hits outside `design_docs/`)

**(a) MUST MOVE to 0.1.1 — 15 sites** (13 literal `0.1.0` hits, the escaped dry-run regex, and one
`0.1.1` literal that becomes wrong once the golden says 0.1.1):

| Site | Change |
|---|---|
| `packages/world-core/ailang.toml:3` `version`, `:10` comment | `0.1.1`; comment cites D-WORLD-41 and the kept `kind = "feature"` |
| `scripts/verify_world_package.sh:133` (step 4 frozen manifest), `:214` (step 7 `^Publishing world/core@0\.1\.0\.\.\.$`), `:303` (step 9 packet `version`) | `0.1.1` in all three; independent literals, kept that way (two-source rule) |
| `scripts/world_package_ready_packet.golden.json` | regenerated (V8): `version 0.1.1`, `tarballBytes 10610`, `tarballSHA256 sha256:07d6180986739249bc4e5adaf3eb2bf44ce59821861a86652c31149bc09df4a3` |
| `cmd/world-publish/fences.go:66` `attendedPhrase`, `:71` `frozenPackageVersion` | `publish world/core@0.1.1 irreversibly`; `"0.1.1"` |
| `cmd/world-publish/fences_test.go:397` R-PACKET-VERSION trigger | literal `"version":"0.1.0"` → derived from `frozenPackageVersion` (else the trigger silently becomes a no-op — MUT-9) |
| `host/pkgproj/readypacket_test.go:153` version mutation `p.Version = "0.1.1"` | `p.Version += "-x"` (else mutant == base — MUT-10) |
| `tools/attended/self_mod_publish.sh:2` header, `:224` banner | banner version READ from the golden (like `compilerVersion` already is) |
| `docs/SELF_MOD_PUBLISH.md:8`, `:83–88`, tarball digest row, step 6 store, step 7 phrase note | candidate wording; `tarballSHA256` row; fresh-store export (F2, recommended pending D-WORLD-45) |
| `packages/world-core/CHANGELOG.md` | NEW `## 0.1.1 — 2026-09-30` section above the untouched `## 0.1.0` |

**(b) MUST STAY 0.1.0 — 19 hits (preserved published history / captured tool output):**
`host/broker/testdata/metadata_world_core_0.1.0.json` (2), `host/broker/testdata/CAPTURED.txt` (1),
`host/broker/registry_reconcile_v2_test.go` (2, pins the served record), `CHANGELOG.md` `## 0.1.0`
(1, byte-equal to the registry's recorded notes), `host/broker/registry_reconcile_test.go:51,59`
(2, captured `NoSuchKey` body from when 0.1.0 was absent), `host/pkgproj/testdata/*quality*.json`
(9, `pkg quality` captures "as produced" — regenerating would rewrite captured evidence),
`cmd/world-publish/fences.go:91` (1, `frozenInterfaceHashV2` correctly names the registry-recorded
0.1.0 v2 hash; comment annotated — V7/V14 show the candidate's v2 is identical, so the value stays),
`scripts/verify_world_package.sh:127` (1, quotes the historical PUB001 refusal; annotated).

**(c) UNRELATED — 28 hits:** `host/daemon/daemon.go:68` (daemon's own semver),
`host/projection/projection_test.go` (2, agent card), `host/broker/registry_publish_test.go` (2,
fake-registry fixture version), `approve_test.go` (2), `publish_op_test.go` (1),
`registry_reconcile_test.go` (9 more, URL-grammar arms), `registry_publish.go:65`, `approve.go:326`
(scope-grammar examples), `host/pkgproj/iface_test.go:29`, `readypacket_test.go:230`
(version-ignoring interface test), `cmd/world-publish/wiring_test.go:571` (version-irrelevant
literal), `fences.go:69` (grammar example), `docs/SELF_MOD_PUBLISH.md:189–190` (scope grammar),
`sprint_w-decision-lifecycle-freeze.json` (3). Plus 86 `design_docs/` hits: history, untouched.

**Counts: must-move 15 · must-stay 19 · unrelated 28** (13 + 19 + 28 = 60 literal hits).

## 3. Milestones (commit order; each ≤150 changed LOC)

### M1 — move the release identity to 0.1.1 as ONE unit (~93 changed lines, 9 files)

Files: `packages/world-core/ailang.toml`, `packages/world-core/CHANGELOG.md`,
`scripts/verify_world_package.sh`, `scripts/world_package_ready_packet.golden.json`,
`cmd/world-publish/fences.go`, `cmd/world-publish/fences_test.go`,
`host/pkgproj/readypacket_test.go`, `tools/attended/self_mod_publish.sh`, `docs/SELF_MOD_PUBLISH.md`.

Why one unit: manifest, CHANGELOG, golden, gate literals, publisher fence and runbook digest are
cross-checked (MUT-1/3/4/11/12); any subset reds a gate. Regenerate the golden ONLY via V8 (the
step-9 `+{…}` line); never hand-type a hash. Keep `[release] kind = "feature"` (D-WORLD-41). The
CHANGELOG `## 0.1.1` date is part of the tarball bytes: changing it later means regenerating the
golden and runbook row (the drift fences enforce this).

Acceptance (all rc=0 unless stated):
1. `E ./scripts/verify_world_package.sh` → last lines `✓ world package gate PASSED: 9/9 steps performed non-zero work` and `PUB011 ratchet passed: uncontracted_exports=3`.
2. `E ./scripts/verify_ail.sh` (pin first on PATH) → `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass`.
3. `E go vet ./...` → rc=0, no output.
4. `E go test -race -count=1 ./cmd/world-publish/... ./host/broker/... ./host/pkgproj/... ./host/runbook/... ./host/verifygate/...` → 5× `ok`.
5. `go build -o "$TMPDIR/wp" ./cmd/world-publish && E "$TMPDIR/wp" packet` → `ready packet for world/core@0.1.1`, `tarballSHA256 sha256:07d6180986739249bc4e5adaf3eb2bf44ce59821861a86652c31149bc09df4a3`, `EQUAL to the committed golden in every field`.
6. `(cd packages/world-core && E "$AILANG_BIN" publish --dry-run)` → `Publishing world/core@0.1.1...`, `changelog_section true`, `✓ no gates`, `Tarball ready but not uploaded`.
7. `git diff HEAD -- host/broker/testdata host/pkgproj/testdata` → empty; `git diff HEAD -- packages/world-core/CHANGELOG.md | grep '^-[^-]'` → empty (history only appended).

### M2 — guard the published history (~73 LOC, new `cmd/world-publish/release_history_test.go`, first half)

`changelogSections` + `TestChangelogCarriesTheCandidateAndPreservesPublishedHistory`: newest
CHANGELOG section == `frozenPackageVersion` and non-empty; `frozenPackageVersion` ≠ the preserved
record's version; CHANGELOG `## 0.1.0` body == the registry-recorded `quality.release.notes` in
`metadata_world_core_0.1.0.json` (instrument guards: ≥2 sections parsed, non-empty notes).
Acceptance: `E go test -count=1 -v -run TestChangelogCarries ./cmd/world-publish/` → `--- PASS`;
MUT-2b and MUT-7 red it.

### M3 — rehearsal and reconcile-conflict arms (~88 LOC)

`TestAttendedRehearsalNamesTheCandidate` (append to `release_history_test.go`; drives
`publish --dry-run` with the satisfied probe; asserts the REHEARSAL line, phrase, package, approval
scope `version:0.1.1`, golden tarball) and new `host/broker/release_candidate_test.go`:
`TestReconcileCandidateAgainstPublished010IsConflict` (candidate config vs served 0.1.0 doc →
`conflict`, "served document identifies world/core@0.1.0, want world/core@0.1.1") and
`TestReconcileCandidateBytesAsVersion010IsDigestConflict` (candidate hashes reconciled AS 0.1.0 →
digest `conflict`; same-test control: `cfg01()` still `succeeded-reconciled`).
Acceptance: `E go test -count=1 -v -run 'Candidate' ./cmd/world-publish/ ./host/broker/` → three
`--- PASS`; then repeat M1 acceptance 1–4 at the final boundary.

### M4 — bookkeeping (controller, not executor)

Charter row 94 M6 status → "0.1.1 candidate prepared + rehearsed; awaiting attended publish";
no claim that the registry serves 0.1.1.

## 4. Mutation ledger (all run against the prototype by `mutate.py`, restored by byte copy; `shasum -c` of all 11 files OK afterwards)

| ID | Mutation | Named killer | Result |
|---|---|---|---|
| MUT-1 | manifest `version` left `0.1.0` (golden/gate say 0.1.1) | package gate **step 4** `✗ manifest is not the exact frozen structure` | KILLED |
| MUT-1b | same, Go side | `TestWorldCoreManifestMatchesTheCommittedGolden` (reds first at `tarballBytes 10608 vs 10610` — the toml is in the tarball) | KILLED |
| MUT-2 | CHANGELOG `## 0.1.1` heading removed | package gate **step 7** dry-run `✗ PUB001 … no non-empty '## 0.1.1' section` | KILLED |
| MUT-2b | same, Go side | `TestChangelogCarriesTheCandidateAndPreservesPublishedHistory` (instrument arm: 1 section parsed) | KILLED |
| MUT-3 | golden left at the 0.1.0 packet | package gate **step 9** `✗ ready packet differs byte-for-byte from golden` | KILLED |
| MUT-4 | `frozenPackageVersion` left `"0.1.0"` | `world-publish packet` → `STOP fence=packet reason=version` rc=3 | KILLED |
| MUT-5 | approval scope drops the version (a 0.1.0 grant would authorize 0.1.1) | `TestPublishApprovalRefusalSetWithALandedPositiveControl` (wrong-scope arm) | KILLED |
| MUT-6 | preserved 0.1.0 fixture `tarball_hash` edited | `TestReconcilePublished010OverFourDigests` → `conflict` | KILLED |
| MUT-7 | historical `## 0.1.0` notes rewritten (`11 verified` → `16 verified`) | `TestChangelogCarries…` `published history was rewritten` | KILLED |
| MUT-8 | `attendedPhrase` left at 0.1.0 | `TestAttendedRehearsalNamesTheCandidate` | KILLED |
| MUT-9 | R-PACKET-VERSION trigger keeps literal `"0.1.0"` | `TestEveryRefusalBranchStops…/R-PACKET-VERSION` | KILLED |
| MUT-10 | pkgproj version mutation keeps `"0.1.1"` | `TestReadyPacketEqualNamesTheFirstDifferingField` `Equal reported EQUAL after mutating version` | KILLED |
| MUT-11 | runbook `tarballSHA256` row left at `d532d923…` | `TestRunbookDigestsAppearVerbatimInTheCommittedGolden` | KILLED |
| MUT-12 | step-7 dry-run identity regex left at `0\.1\.0` | package gate **step 7** `✗ dry-run missing package identity` | KILLED |
| MUT-13 | candidate reconciled as 0.1.0 treated as success | `TestReconcileCandidateBytesAsVersion010IsDigestConflict` | KILLED |
| MUT-14 | helper banner version | no test reads the banner (only the runbook AC30 scan, which admits the helper by name) | NOT GUARDED — accepted: display-only, value is read from the golden and `die`s if empty |

15/15 guarded mutations killed. Executor: re-run MUT-1…MUT-13 (script and logs in the evidence dir).

## 5. The attended publish — for a HUMAN, never run by the loop

At a real terminal on the rig, in a checkout of the landed commit:

```bash
export WORLD_STORE="$HOME/.ailang/world/world-0.1.1.db"   # fresh store: RECOMMENDED, pending D-WORLD-45 — see F2
./tools/attended/self_mod_publish.sh preflight   # packet must read world/core@0.1.1, 0 indeterminate
./tools/attended/self_mod_publish.sh approve     # type: publish world/core@0.1.1 irreversibly
./tools/attended/self_mod_publish.sh dry-run     # REHEARSAL — no request
./tools/attended/self_mod_publish.sh live        # PRINTS the command below; paste it yourself
"$WORLD_BIN" publish --live --store "$WORLD_STORE" \
  --registry-origin https://storage.googleapis.com/ailang-registry \
  --publisher "$HOME/.pinned-ailang/ailang" --credential-file "$HOME/.config/ailang/registry.key" \
  --approval-ref <ref minted by approve> --now 2 --expires 1000 < /dev/tty
```

Then `./tools/attended/self_mod_publish.sh reconcile` must read `succeeded-reconciled` for 0.1.1.
Run `approve` fresh: the saved `~/.ailang/world/.last_approval_ref` is a consumed 0.1.0 ref (V17);
it names no object in the fresh store and its scope is `version:0.1.0`, so `publish --live`
would refuse it before any POST — a wasted step, not a hazard.

## 6. Findings and false premises

- **F1 (directive premise).** "Run the attended rehearsal/dry-run path" cannot be done headless:
  `world-publish publish --dry-run` and `approve` sit behind the controlling-terminal fence and
  STOP (`fence=tty`, rc=3) in this loop by design (V13). Substitutes measured: the read-only
  `packet` verb, the compiler's `publish --dry-run`, `reconcile --probe` (GET-only), and the
  in-process rehearsal with the injected probe (V14–V15).
- **F2 (new, blocks the documented attended path).** The operator's store
  `~/.ailang/world/world.db` (which holds the 0.1.0 publish record) is schema v2; the current binary
  requires v4 and refuses to open it (V16). The runbook/helper default would STOP at
  `fence=store reason=unopenable`. Decision **pending D-WORLD-45 (Mark)**: it affects provenance continuity of the attended
publish, so it is not the loop's call. RECOMMENDED option: publish 0.1.1 from a FRESH store file
(`world-0.1.1.db`); alternative: migrate the legacy store first (future queue row). The legacy
store must never be modified or deleted. The runbook presents both and exports the fresh-store
path as the recommended default; helper default unchanged (it fails loudly, not wrongly).
- **F3 (design §4 premise).** "Updating the local golden … changed sources hit the local drift
  fence" still holds, but **contentHash does NOT move** with the manifest/CHANGELOG (V7): the
  release delta is visible only in `tarballSHA256`/`tarballBytes`/`version`. The approval scope
  binds tarball + version, so it is still exact.
- **F4 (hidden coupling).** Two tests used literal version strings that silently break when the
  golden moves (`fences_test.go:397`, `readypacket_test.go:153`); a `0.1.0`-only grep misses the
  second. Both now derive or suffix (MUT-9/10).
- **F5.** `world-publish transitions` (row 106) reuses `attendedPhrase`, so its operator now types
  `publish world/core@0.1.1 irreversibly`. Tests read the constant; behaviour-neutral, noted only.
- **F6.** Full-suite load flake `TestQueryInterfaceReturnsWhileADescendantHoldsStdout` (V11) —
  unrelated, row-114 class.

One OPEN question requires Mark: D-WORLD-45 (F2 store choice). Version and kind were ruled (D-WORLD-41).
