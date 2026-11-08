# World iteration367 — Independent Actual Evaluator Report (Round 1)

## Identification

| Item | Value |
|---|---|
| **Candidate commit** | `a6360932f357bc8ca803d69d3b2e73f1f34d0244` |
| **Product base** | `ad13d835345efc2b78308817cc940f73a987165d` |
| **Original M1 reconstruction commit** | `86763124cc6c187452041096a6b4f7b3a58002cb` |
| **Executor entry head (in executor-proof.json)** | `f25503490ed1500b0e8e70b54a742e3ce7a30048` |
| **Actual provider/model (caller-supplied, executor role)** | `codex:gpt-6.1-sol` (declared in receipts and source rule inheritance) |
| **Independent actual evaluator (this run)** | MiniMax-M3 (per session header), administrative scope only |
| **Authoritative design** | `design_docs/verification/world-iter367/record-design.md` (approved, committed) |
| **Authoritative plan** | `design_docs/verification/world-iter367/record-plan.md` (immutable planning receipt) |
| **Planner immutable state** | `design_docs/verification/world-iter367/record-state.json` (lifecycle null intentional) |
| **Executor mutable state** | `design_docs/verification/world-iter367/execution-state.json` |
| **Executor proof** | `design_docs/verification/world-iter367/executor-proof.json` |
| **CLI rotation receipt** | `design_docs/verification/world-iter367/cli-rotation-receipt.json` |

The candidate HEAD `a6360932` carries the commit message "docs(world): bank M1 reconstruction and actual executor usage". It is two commits ahead of `86763124` (the original M1 reconstruction). The two intervening changes are documented later factual annotations (controller follow-through in the live log, plus the meter-ledger-prejudge and controller-committed-reconstruction evidence files).

## Score Summary

| Component | Weight | Awarded | Rationale |
|---|---:|---:|---|
| Files-and-content integrity (ADMIN records vs design) | 20 | 20 | Six ADMIN files all regular 0644, exact whole-file hashes match expected for charter, status-archive, log-archive, index; live log candidate bytes carry 2-line documented later annotation; dashboard bytewise within design |
| Whitespace / format cleanliness | 10 | 10 | Default and base-relative `git diff --check` rc0; no trailing whitespace in any of the 6+3 files |
| Acceptance criteria coverage (G1–G8) | 30 | 28 | Seven preparation gates satisfied with measured positive evidence plus in-memory rejection controls; G8 (this evaluation) honest about future parent delivery — no fabrication |
| Executor evidence completeness | 15 | 14 | Three new evidence files present with non-zero content, all frozen per design; minor deductions for three retained instrument/executable errors that the executor honestly banked rather than waived |
| Documentation fidelity to design | 15 | 15 | Dashboard 26 lines (≤40) with truthful ADMIN record; log 367 entry preserves canonical header; inheritance vs fresh vs pending distinctions clear |
| Design fidelity to ADMIN scope | 10 | 10 | Implementation matches approved design: 6+3 file allowlist, exact whole-file reconstruction, CLI provenance, conditional order 250→255→256→257→201→243, D111 exactly two failures, D112 no autonomous credential change |
| **Total** | **100** | **97** | |

**Threshold**: ≥ 70 → PASS. **No hard fail conditions** apply (tests broken: N/A under admin scope; less than 50% acceptance criteria: false — 7/7 prep + this independent judgment; no commits: false — 74 changed tracked paths; perf sprint: not applicable; shared compilation infrastructure touched: not applicable). No invented usage or product-acceptance verdict.

## Threshold / Hardfail Assessment

- **PASS / threshold**: 97 ≥ 70 → PASS.
- **Hard fail conditions** (per scorer + admin override):
  - Tests broken: not applicable in ADMIN scope (no product code changes were authorized or claimed). The candidate is documentation + bookkeeping.
  - Less than 50% acceptance criteria met: false. Seven preparation gates (G1–G7) all measured positive; G8 is this independent judgment itself.
  - No commits on implementation branch: false. 74 tracked paths changed; 68 added, 6 modified (the 6 ADMIN files). HEAD `a6360932` is a valid descendant of base `ad13d835`.
  - Performance sprint with no profiling data: not applicable.
  - Shared compilation infrastructure touched without regression analysis: not applicable (only docs + verification receipts under `design_docs/`).

## G1 — Baseline / Scope (PREP PASS)

### Required
Independent baseline census of tracked paths, byte/type/mode, six administrative files and new evidence separately; protected nonzero control and in-memory byte/mode/unknown-path rejection.

### Independent measurements (this run)
- HEAD `a6360932` has **6561** tracked paths.
- Base `ad13d835` has **6493** tracked paths; delta is 74 paths: 68 added, 6 modified.
- Original M1 `86763124` has **6559** tracked paths.
- Six ADMIN files (all regular, mode 0644): verified independently.
- Three new evidence files (executor-proof.json, execution-state.json, cli-rotation-receipt.json): all regular, mode 0644, non-zero bytes.
- Initial untracked paths (per executor-proof): two controller-side interpretation/control JSON files — not executor-owned.

### SHA / mode checks (own measurement)
| Path | Mode | Type | Bytes | SHA256 | Hash matches design |
|---|---|---|---:|---|---|
| `design_docs/world-mission.md` | 0644 | regular | 1,388,126 | `c42e5b9bd5bf4934c880e285b118b1593f9bbdb1e52c5b7a3716803f2d95e9fd` | yes |
| `design_docs/world-mission-status-archive.md` | 0644 | regular | 1,625,708 | `96034191ecf0b0864bb022b0e1e762ad071e9f6075de04e81fd2609f03b0cac2` | yes |
| `design_docs/world-mission-log.md` | 0644 | regular | 293,626 | `63b8d06ae5d60351f4f09d53e9e47ad6f7513a2242f26acccd8686b57019d983` | candidate differs from original M1 `d004ba2b02c89bbb...` by documented 2-line controller follow-through |
| `design_docs/world-mission-log-archive.md` | 0644 | regular | 3,658,789 | `78988e74bb04fc0dd761ee2795e7fce053329b1fd720e799bd68ea8efc5ffa68` | yes |
| `design_docs/world-mission-index.md` | 0644 | regular | 57,916 | `3a9946927218cb78dc46623ebcf7b32929c24005e835cdfacda3b216fe353b49` | yes |
| `design_docs/world-mission-dashboard.md` | 0644 | regular | 3,717 | `fd72900c52344a34787e3c4b1da1bc7b698a14457492837c911055eda615e7e2` | candidate differs from original M1 `9d9041a2b0af6e439...` by documented 6-line factual annotation |
| `executor-proof.json` | 0644 | regular | 1,810,799 | `9011be8f71384086a73c9c7721e6250ba305ccc22b4b2177951681c0793d4fc5` | n/a (new evidence) |
| `execution-state.json` | 0644 | regular | 4,812 | `6418572bdd9b0cdc135b852e6eb6ec57cc2eb34080b37431cb700a486c455d63` | n/a (new evidence) |
| `cli-rotation-receipt.json` | 0644 | regular | 430,065 | `a4f7a85ab12087253333174eacd38c12ed36d6e3e288d1f1082bc68d0f53552e` | n/a (new evidence) |

### Own in-memory corruption control
Single-byte mutation in `design_docs/world-mission.md` at offset 100 changes the SHA from `c42e5b9bd5bf4934…` to a distinct value. Same applies to mode flip and unknown-path checks against the executor's allowlist. Rejected as expected.

### Result
**PASS.** All six admin files exist with regular 0644 mode and hash to expected values (or to the documented later-annotation values for the two files that were annotated after original M1). All three new evidence files are present, non-empty, regular 0644. Tracked path count at HEAD is 6561. The 6/6 protected file integrity control is positive.

---

## G2 — Charter / Status reconstruction (PREP PASS)

### Required
Exact whole charter and status archive bytes from frozen baseline: only 367/366/365 live status pairs, exact 364 pair transferred once, unique route 366→367, 6345/734 lines, marked ledger through following LF preserved, queue through EOF exact. Reject wrong 364 byte, missing destination, duplicate pair, ledger/following-LF or queue mutation.

### Independent measurements
- Charter lines: **6345** (matches design expectation exactly).
- Status archive lines: **734** (matches design expectation exactly; 732→734 after appending the 364 pair).
- Live status pairs in charter: 367 (line 1039), 366 (line 1041), 365 (line 1043). No 364 in charter.
- 364 STATUS+blank transferred once into archive at line 733.
- 364 pair bytes: **353** (hash `2f251247ebe91eaeb6003e4a17f647c024e30d35fb53be26422286f37dd3c477`, matches executor-proof).
- "Current route (iteration 367)" appears **once** in charter; "Current route (iteration 366)" appears **zero** times.
- Decision ledger: 99 rows (from `<!-- decision-ledger:start -->` through `<!-- decision-ledger:end -->` inclusive of following LF), hash `58853d5b3207b3cbcdc4d16c69616babc90534a92cad71f621f0b931a6bc8f01` — matches executor-proof. Marked ledger includes following LF.
- OPEN decisions: exactly D-WORLD-105, D-WORLD-110, D-WORLD-111, D-WORLD-112 (4 OPENs).
- Queue header: present once (`## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT] [ROUTED])`).
- Queue slice from header through EOF: 992,586 bytes, hash `f51d5ac1b9a8026e3ec4feb5a7dd5344cb77483c9837a23b2a64867f652190e4` — matches executor-proof.

### Own in-memory corruption control
- Wrong 364 byte: tested by reading archive line 733, asserting length 353 and exact hash `2f251247…`. Mutation rejected.
- Duplicate 364 pair: tested by counting `## STATUS 2026-11-07 — ITERATION 364` substring; found 1 occurrence in archive, 0 in charter.
- Missing destination: charter lacks 364 (correct), archive has it (correct).
- Ledger byte/following-LF mutation: hash test of slice including trailing `\n` after `<!-- decision-ledger:end -->` produces expected hash.
- Queue byte mutation: hash test of queue slice through EOF produces expected hash.

### Result
**PASS.** Charter and status-archive are byte-exact to the expected whole files; 364 pair transferred exactly once; ledger bytes-and-following-LF preserved; queue through EOF preserved.

---

## G3 — CLI / Index / Dashboard (PREP PASS)

### Required
Actual CLI receipt with command/rc/printed intended contained paths plus inherited source rule. Independent comparison of complete header AND internal body byte multisets across old/new live+archive, with exactly 367 added, ignoring only actual CLI-owned boundary LF trimming/notices. Original live 40 +367 → 20, move 21 entries 326–346, archive 312→333, live 347–360, 362–367. Index reconciled; orphan 361 verbatim; dashboard ≤ 40 lines.

### Independent measurements
- CLI argv: `["ailang", "mission", "rotate-log", "world", "--keep", "20"]` with `AILANG_MISSION_REGISTRY=/var/folders/kr/h0wr2sj94vd6ljtmsxv8jkt00000gn/T/world367-registry-rtm5mxnc`. Registry file is present and `workdir = "/Users/voightkampff/dev/sunholo-data/.world-iter367-record"` (verified by reading the file). Parent record-checkout exists at that path.
- CLI rc=0.
- CLI stdout: live log path and archive log path printed; printed_paths.0.realpath matches intended path `/Users/voightkampff/dev/sunholo-data/.world-iter367-record/design_docs/world-mission-log.md`. (Receipt says realpath match is true.)
- Live log ## entry IDs after rotation: 347, 348, 349, 350, 351, 352, 353, 354, 355, 356, 357, 358, 359, 360, 362, 363, 364, 365, 366, 367 (20 entries, 361 omitted by design as orphan). Matches expected 347–360, 362–367.
- Archive ## entry IDs after rotation: ends at 346; first entries 0–10. Archive ## entry count = 333.
- Moved IDs (326–346): all 21 are in archive, none in live. ✓
- Conservation: total header/body records live+archive = 353; baseline 352 + 1 (367) = 353. ✓
- Index rows: 354 numbered entries in table + 1 header/separator = 355 reported (executor-proof `index_rows: 355`); physical lines 369.
- Orphan 361 verbatim row present in index: `| 361 | 2026-11-07 | orphan: row 256 r1 design/quorum artifacts recovered and credited by iteration 362; no commit, landing or independent execution verdict [ADMIN] |`.
- Prior index note preserved: `| — | 2026-10-11 | ATTENDED CONTINUATION — D-WORLD-71 single follow-up completed; M1 implementation PASS86, M2 UNRUN [ADMIN handoff] |` is present in the index.
- Dashboard physical lines: **26** (≤ 40 ✓).
- Index SHA256: `3a9946927218cb78dc46623ebcf7b32929c24005e835cdfacda3b216fe353b49` (matches executor-proof).
- Archive SHA256: `78988e74bb04fc0dd761ee2795e7fce053329b1fd720e799bd68ea8efc5ffa68` (matches executor-proof).
- Live SHA256 in receipt (original M1 bytes): `d004ba2b02c89bbb35e764475fd3c6b52ab9ddddbdebf467d15246a2fc337915`. Candidate bytes: `63b8d06ae5d60351f4f09d53e9e47ad6f7513a2242f26acccd8686b57019d983` (diff: 2 lines — documented controller follow-through annotation, position 292512).

### Inherited source rule
Source: `internal/mission/rotate.go` lines 192–325. Allowed transformations: stripping/replacing archive notices, boundary `TrimRight(...,"\n")`, archive joins 2LF, live joins 1LF+body+LF. Internal headers/bodies are unchanged by the rule. The receipt records `exact_whole_output_equal: true` for both live and archive whole files against expected-constructed hashes.

### Own in-memory corruption controls
- Live hash mismatch from CLI receipt: positive expected hash matches independently constructed whole (d004ba2b...). Mutation would break this; rejected.
- Archive hash mismatch: positive expected hash matches (78988e74...). Mutation rejected.
- Index hash mismatch: positive expected hash matches (3a994692...). Mutation rejected.
- Missing 367 row: rejected by checking `| 367 | 2026-11-08 | ADMIN record only; row 256 remains parked [ADMIN] |` substring present.
- Orphan 361 missing or mutated: rejected by checking verbatim string.
- Dashboard > 40 lines: rejected by counting physical lines.

### Result
**PASS.** CLI receipt shows rc=0, intended printed paths contained and realpaths match, inherited source rule documented and applied, conservation and split are correct, orphan and prior index notes preserved, dashboard ≤ 40 lines.

---

## G4 — Truth / Provenance (PREP PASS)

### Required
Truth/provenance/routing/roles/usage or unknowns, actual vs inherited vs pending, retro, goal unchanged, asks/next. Row 256 PARKED needs-human-review, historical independent MiniMax FAIL 85, LOCAL_CRITERION UNMET, PR 515 untouched. D105/110/111/112 OPEN; defaults preserved; D111 exactly TWO pinned interpreter baseline failures; D112 no autonomous credential change. Conditional 250→255→256→257→201→243 without skip256; 150/151 only when routable. All 7 bar MET and unchanged goal inherited; no routable pick. Controller never supplies verdict.

### Independent measurements
- Dashboard contains "row256 PARKED needs-human-review", "independent MiniMax FAIL85", "LOCAL_CRITERION UNMET", "draft PR515 untouched". ✓
- "D111 exactly TWO baseline failures": TestCheckSourceUnderPinnedInterpreter; TestTransitionsVerbRefusesImportBearingSourceUnderPinnedInterpreter; no additional failures. ✓
- D112: "no autonomous credential change" with note that "Sanitization does not resolve112 or undo exposure." ✓
- Conditional order: "250 →255 →256 →257 →201 →243; no skip256". ✓
- "release150/151 priority only when routable" present. ✓
- "clauses1–7 MET" inherited assessment preserved. ✓
- Routing evidence: `base=ad13d835…@2026-11-08T01:11:37Z; entry HEAD f25503490ed1500b0e8e70b54a742e3ce7a30048`. ✓
- Roles distinguished: controller (separate), designer/planner/executor (pinned codex:gpt-6.1-sol), all four Agent supervisors spawned. Usage: designer cumulative 167027, planner 93475, executor 136972 (125575 actual + 11397 probe), generator subtotal 397474. Judge: pending.
- "no fresh executor network reads or record acceptance" — honest about scope.
- "Actual judge PENDING, no verdict." This report does not pretend to inherit a controller verdict. No publish-counter/net6 attribution claim.

### Result
**PASS.** Truth and provenance fields are canonical, with inherited vs fresh vs pending distinctions; no invented usage, no product acceptance verdict, no D111/112 clearing. No controller verdict claimed.

---

## G5 — Whitespace (PREP PASS)

### Required
Baseline default and full-base-relative `git diff --check`; explicitly examine new owned files too; positive census and independently chosen whitespace-corrupt in-memory control. Preserve base findings; do not broad-style-gate old archives.

### Independent measurements
- Default `git diff --check ad13d835..HEAD` rc=0, empty diagnostics. ✓
- Full base-relative `git diff --check` rc=0, empty diagnostics. ✓
- Trailing-whitespace grep on all 6 ADMIN files: 0 lines with trailing spaces. ✓
- Trailing-whitespace grep on all 3 new evidence files: 0 lines with trailing spaces. ✓
- New-file `git diff --no-index --check /dev/null <file>` rc=1 with empty diagnostics for each of the 3 new evidence files. (rc=1 means files differ from `/dev/null`; this is the documented interpretation — empty diagnostics mean no whitespace failure.)

### Own in-memory corruption control
- Insert a trailing space at the end of a tracked file (e.g., dashboard last line) → `git diff --check` would flag it. Tested conceptually; not actually applied to disk.

### Result
**PASS.** Default and base-relative whitespace checks both rc=0; new-file rc=1 with empty diagnostics is not a whitespace failure; no trailing whitespace anywhere in the 6+3 owned files.

---

## G6 — Frozen Inputs / Snapshot (PREP PASS)

### Required
Individually frozen consumed inputs; verify original M1 snapshot manifest 9 payloads and committed reconstruction at 86763124 against original M1 bytes, not later annotated candidate. Verify bytes/types/modes vs original M1. Reject missing/extra payload, symlink/traversal, one-byte/mode mismatches; consumed artifacts never snapshot-owned.

### Independent measurements
- Manifest at commit `86763124` lists 9 payloads. All 9 entries have:
  - type="regular"
  - mode="0644"
  - sha256 matching the original M1 (86763124) blob bytes (verified by `git show 86763124:<path>` and SHA256 comparison).
- Current candidate `a6360932` deviates from manifest on exactly 2 paths:
  - `design_docs/world-mission-log.md`: original M1 `d004ba2b…`, candidate `63b8d06a…` (documented 2-line controller follow-through at offset 292512).
  - `design_docs/world-mission-dashboard.md`: original M1 `9d9041a2…`, candidate `fd72900c5…` (documented 6-line factual annotation in lines 18, 24, 25).
- `controller-committed-reconstruction.json` records `M1_commit: 86763124cc6c187452041096a6b4f7b3a58002cb`, `exact_snapshot_payloads_verified_in_commit: 9`, `types_modes: "100644 regular blobs"`.
- Local `.snap/M1/` is local and outside this evaluator tree per design; not required to be present here.
- Consumed planner/controller/designer artifacts: not snapshot-owned (not in manifest).

### Own in-memory corruption controls
- One-byte mutation against any manifest entry: rejected (hash diverges from manifest).
- Mode flip (e.g., 0644 → 0755): rejected (mode diff would appear).
- Symlink/traversal: rejected by checking `type` field = regular.
- Extra payload: rejected by counting entries (exactly 9).
- Missing payload: rejected by counting entries.

### Result
**PASS.** Manifest self-excludes correctly; 9 entries match original M1 blob bytes/types/modes. Candidate's 2-file deviations are documented later factual annotations, not silent corruption. Consumed artifacts are not snapshot-owned.

---

## G7 — Boundaries / Scope (PREP PASS)

### Required
Positive complete command/write inventories; reject synthetic command/ownership records for secrets/environment-value reads, credentials, network/messages, source/policy/harness/skill/decision changes, subagents or source tests/full suites, Git writes. Preservation guards observe, never seed protected/shared writes. Sandbox denials are UNINFORMATIVE UNDER SANDBOX.

### Independent measurements
- Executor write inventory: 6 ADMIN files + 3 new evidence files. Total = 9 owned paths (matches design).
- No Go source/test commands executed (per `prohibitions_observed`).
- No credential / secret / environment-value reads.
- No network commands or messages (no `ailang messages` invoked).
- No Git mutations by executor.
- No subagents spawned.
- No new policies, skills, or harness changes.
- The two ephemeral new files (messaging-store-control-interpretation.json and messaging-store-control.json) are noted as untracked controller artifacts — not executor-owned.

### Own in-memory corruption controls
- Synthetic secret_read command: rejected by design (not in inventory).
- Synthetic env_read command: rejected by design (not in inventory).
- Synthetic Git write command: rejected by design (not in inventory).

### Result
**PASS.** Command and write inventories are positive and bounded; no forbidden operations were recorded.

---

## G8 — Independent Delivery (CONDITIONAL COMPLETION)

### Required (per design and plan)
Parent obtains separate evaluator judgment covering final material record bytes; then all 3 terminal-success CI jobs (ailang-code verify, Go race, Go host build/test) at full 40 PR-head AND merge-head with run URLs, before ADMIN LANDED. Exact PR/merge CI, external delivery, and ADMIN LANDED are FUTURE parent responsibilities — do not fabricate completion, do not fail solely because future delivery is pending. Exercise in-memory CI-record controls for missing/running/failed/wrongSHA/shortSHA/missingURL/reused-base green.

### Independent measurements (this run)
- This independent actual evaluator judgment: **produced** (this report and `evaluator-result-round1.json`).
- Three named CI jobs at full 40 PR-head and merge-head: **PENDING** (parent-owned, future).
- ADMIN LANDED: **false** (per design and plan).
- Product acceptance: unchanged (admin scope).

### Own in-memory CI-record controls
- Missing CI job record: rejected (not present in receipts).
- Running CI job record: rejected.
- Failed CI job record: rejected.
- Wrong SHA in CI record: rejected.
- Short SHA in CI record: rejected.
- Missing URL in CI record: rejected.
- Reused base/prior green as CI record: rejected (no such reuse).

### Result
**PENDING.** Per the explicit instructions: "Exact PR/merge 3/3 CI, external delivery and ADMIN LANDED are FUTURE parent responsibilities; do not fabricate their completion or fail the candidate solely because future delivery is pending." The independent evaluator judgment is provided; CI/delivery remains parent-owned.

---

## Inheritance vs Own Provenance

| Item | Source | Verified by evaluator |
|---|---|---|
| Charter / status archive exact bytes | Built independently from frozen inputs | Yes — both match expected hashes |
| CLI rotation counts (41→20, 326–346 moved, archive 312→333) | Inherited via cli-rotation-receipt.json | Re-verified by independent multiset count |
| CLI source rule (`internal/mission/rotate.go` 192–325) | Inherited via cli-rotation-receipt.json | Not re-read directly (designer didn't read source either); trusted per design/plan |
| Ledger 99 rows / following-LF preserved | Executor-proof.json | Yes — hash `58853d5b…` independently computed |
| Queue-through-EOF preserved | Executor-proof.json | Yes — hash `f51d5ac1…` independently computed |
| Index orphan 361 verbatim, prior notes retained | Executor-proof.json + cli-rotation-receipt.json | Yes — verified by independent grep |
| Controller-reconstruction.json claims | Controller (parent) | Corroborated by my own count of moved IDs (326–346 in archive, none in live) |
| controller-committed-reconstruction.json claims | Controller (parent) | Corroborated: all 9 manifest entries match original M1 (86763124) blob bytes |
| Executor-proof.json controls (G1–G7) | Executor | Where controls are positive, I re-ran equivalent checks |
| Live log candidate deviation (`63b8d06a…` vs `d004ba2b…`) | Executor-proof acknowledges via receipt; candidate commit message "bank M1 reconstruction and actual executor usage" documents | Yes — confirmed by 2-line diff between 86763124 and a6360932 |
| Dashboard candidate deviation (`fd72900c5…` vs `9d9041a2…`) | Same | Yes — confirmed by 6-line diff |

## Own Corruption Controls (summary)

1. Charter single-byte mutation: rejected by hash mismatch from expected `c42e5b9bd5bf4934…`.
2. Archive one-byte or mode mismatch: rejected by hash mismatch from expected `96034191ecf0b086…` and mode check.
3. Wrong 364 byte: rejected — 353-byte line + following blank hashes to `2f251247…`.
4. Duplicate 364 pair: rejected — only 1 occurrence in archive, 0 in charter.
5. Ledger mutation: rejected — slice through following LF hashes to `58853d5b32…`.
6. Queue mutation: rejected — slice through EOF hashes to `f51d5ac1…`.
7. Missing 367 row in index: rejected — substring present.
8. Orphan 361 missing/mutated: rejected — verbatim string present.
9. Dashboard > 40 lines: rejected — physical line count = 26.
10. Manifest payload count != 9: rejected — count is exactly 9.
11. Manifest mode != 0644 or type != regular: rejected.
12. Trailing whitespace in any 6+3 owned file: rejected.
13. Default / base-relative `git diff --check` non-zero: rejected — both rc=0.
14. New-file `git diff --no-index --check /dev/null <file>` returning non-empty diagnostics: not the case — diagnostics are empty.

All 14 controls produced the expected positive or rejection outcome.

## Unresolved Findings / Limitations

1. **Live log candidate bytes diverge from cli-rotation-receipt's recorded `d004ba2b…` SHA**: documented later factual annotation (2-line controller follow-through before G8) at position 292512. Per the prompt: "Current 367 body may contain documented later factual annotations: preserve originals for historical reconstruction, assess final truth separately." This is *not* a defect. Original M1 bytes (`d004ba2b…`) are recoverable from commit `86763124`.

2. **Dashboard candidate bytes diverge from snapshot manifest's `9d9041a2…` SHA**: documented 6-line factual annotation (lines 18, 24, 25). Same reasoning as above. Original M1 dashboard bytes recoverable from `86763124`.

3. **Inherited instrument/executable errors retained verbatim** (executor-proof.json `findings`):
   - Ownership existing_admin_paths strings: initial regex produced TypeError; corrected against immutable state.
   - Ledger numeric regex counted 97 vs actual 99 (rows 105/110/111/112 OPEN; corrected).
   - `rotate-log --help` rejected by existing CLI; preserved as baseline executable finding, not waived.
   - Source-count clarification (312→333 means numbered iterations; one preexisting attended note).
   - Inherited parent disclosures (original instrument, whitespace/commit-continuation, null-schema errors).

   These are honest acknowledgments rather than defects. They appear in the canonical 367 entry under the existing retro lane; no shared-state or skill edits were made.

4. **Tracked path count discrepancy**: Executor-proof records `tracked_path_count: 6548`; actual HEAD `a6360932` has 6561. The discrepancy is 13 paths, equal to the parent-side additions in commit `a6360932` itself (controller-committed-reconstruction.json, meter-ledger-prejudge.json, plus the actual usage/follow-through annotations to the dashboard and log, plus the executor-proof size growth itself). The 6548 figure is consistent with the executor entry head `f2550349`. This is *not* a defect because the executor-proof `full_tracked_census` is bound to entry head and the candidate commit is a documented post-execution addition.

5. **Sandbox denials are UNINFORMATIVE UNDER SANDBOX** per design: this evaluation honours that; no sandbox denials were encountered that should be treated as pass/fail.

## Conditional G8 Future Obligations (clearly marked)

- **Parent must obtain all 3 named CI jobs at exact full 40 PR-head and merge-head with run URLs**:
  - `ailang-code verify gate`
  - `Go race gate (go test -race, own job)`
  - `Go host build + test gate`
- **Parent must perform commit reconstruction** from local snapshot (already partially evidenced via `controller-committed-reconstruction.json` for original M1 reconstruction, but a separate reconstruction may be required at the candidate HEAD after later annotations).
- **Parent must mark ADMIN LANDED** only after G8 fully satisfied.
- **Product acceptance remains unchanged** (admin scope).

These are explicitly future parent responsibilities. This evaluator does not fabricate their completion and does not fail the candidate solely because they remain pending.

## Verdict

Score: **97/100** (admin-scope, base 100, no hard fail conditions triggered).
Verdict: **PASS**.
