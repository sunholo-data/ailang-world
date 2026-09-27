# w-kernel-exports-carry-no-contract — prove the shared kernel laws and ratchet PUB011

**Status**: Planned — design and working prototype, iteration 198 (designer). Controller owns commits. No publish requested or performed.
**Item**: charter row 94, clause 1, deterministic kernel.
**Scope**: five new proven laws; three measured, in-code exemptions; monotone PUB011 ratchet. PUB012 package-mode test discovery and PUB020 usage documentation are explicitly outside this item.

## 0. Problem, measured

The pinned v0.41.0 baseline reports 19 exports, 11 verified contracts and eight uncontracted exports (V1). The prototype reports 16 verified and three uncontracted (V2). These are capability proofs, not an attempt to turn every quality badge green. The package report explicitly says inline tests are not discovered in package mode (V1, paired package/directory control V18); do not adopt `--strict` as this item's CI gate.

Why is this not a package? These laws specify the existing frozen kernel surface. Extend its proof obligations in canonical `world/`, then project the package; do not introduce a second semantic implementation. The projection script copies the four allowlisted modules, and the gate compares their hashes (V3).

## 1. Per-export decisions

No new preconditions on the five accepted contracts. All five signatures explicitly declare `! {}`. Each `ensures` below is the exact implemented text. The negative specification control is **`ensures { false }`**, replacing only that function's ensures; it must produce that function's `counterexample`, with another identity still `verified` in the same invocation (V4). Implementation mutations supply the stronger evidence that these laws catch plausible bugs.

| Export | Verdict on v0.41.0 | Bug killed / limitation | Evidence |
|---|---|---|---|
| contracts.proposalMatchesWorld | verified | digest-only identity comparison accepts a different algorithm | V4 |
| contracts.verificationMatchesProposal | verified | digest-only proposal identity comparison | V4 |
| contracts.commitAllowed | verified, after contracting its callees | removing `v.accepted` from the body | V4 |
| logepoch.renderRef | exempt; exact law errors | specification interpolation calls unsupported `$builtin.show`; body itself encodes | V5 |
| logepoch.cacheKey | exempt; skipped | body calls unencodable `renderRef` | V5 |
| transitions.plan | verified | substituting `w.logHead` for `w.stateRoot` | V4 |
| transitions.verify | verified, after contracting its callee | setting rejected branch `accepted: true` | V4 |
| transitions.commit | exempt; skipped | body calls unencodable `applyRevision`, even with nonnegative revision precondition | V5 |

### 1.1 Exact accepted specifications

**`proposalMatchesWorld`** — `requires`: none.

```text
ensures { result == (w.stateRoot.algo == p.inputWorld.algo && w.stateRoot.digest == p.inputWorld.digest) }
```

**`verificationMatchesProposal`** — `requires`: none.

```text
ensures { result == (p.proposalHash.algo == v.proposalHash.algo && p.proposalHash.digest == v.proposalHash.digest) }
```

**`commitAllowed`** — `requires`: none.

```text
ensures { result == (v.accepted && w.stateRoot.algo == p.inputWorld.algo && w.stateRoot.digest == p.inputWorld.digest && p.proposalHash.algo == v.proposalHash.algo && p.proposalHash.digest == v.proposalHash.digest) }
```

**`plan`** — `requires`: none.

```text
ensures { result.inputWorld == w.stateRoot && result.proposalHash == proposalHash && result.transitionFn == transitionFn && result.goal == goal }
```

**`verify`** — `requires`: none.

```text
ensures { result.proposalHash == p.proposalHash && result.accepted == (w.stateRoot.algo == p.inputWorld.algo && w.stateRoot.digest == p.inputWorld.digest) }
```

The two identity predicates state field identity independently of `sameRef`; `commitAllowed` states the combined authorization law independently of the shared predicate calls. These are not verbatim body copies. `plan` deliberately states only provenance and goal preservation, not every default field; it does not claim that confidence, evidence or capability defaults are proven. `verify` proves proposal binding and acceptance iff the input world matches; reason strings remain outside this contract.

### 1.2 Exact rejected specifications and in-code exemptions

These are **probes, not contracts left attached merely to lower PUB011**. `scripts/test_kernel_contracts.py` constructs each in an isolated worktree-local directory, runs ai-check, and removes it. No new `.ail` file under `design_docs/` escapes the module census.

| Function | Proposed requires | Proposed ensures (rejected) | Wrong-ensures control |
|---|---|---|---|
| renderRef | none | `ensures { result == "${r.algo}:${r.digest}" }` | `ensures { false }` is counterexample, rc=1 |
| cacheKey | none | `ensures { result == "${h.transitionFn.algo}:${h.transitionFn.digest}@${h.interpreter.algo}:${h.interpreter.digest}" }` | `ensures { false }` still skipped, rc=0; **not** a counterexample |
| commit | `requires { w.revision >= 0 }` | full Applied/Denied law below | `ensures { false }` still skipped, rc=0; **not** a counterexample |

```text
ensures { match result {
  Applied(next, t) => v.accepted && w.stateRoot == p.inputWorld
    && p.proposalHash == v.proposalHash
    && next.revision == w.revision + 1
    && next.stateRoot == outputWorld && next.logHead == nextLogHead
    && t.proposalHash == p.proposalHash && t.inputWorld == p.inputWorld
    && t.outputWorld == outputWorld && t.transitionFn == p.transitionFn,
  Denied(reason) => not(v.accepted && w.stateRoot == p.inputWorld
    && p.proposalHash == v.proposalHash)
} }
```

The commit law would catch accepting an unaccepted verification, denying a valid commit, failing to advance revision, or substituting the wrong output references. It is **not proven**. Do not duplicate/inline the runtime predicates to make it prove. The renderer law would catch a missing colon; the cache law would catch omitting or substituting the interpreter identity. No proof credit is claimed for either. `ensures { true }` would solve none of these requirements.

The exact in-code exemption comments are:

```text
-- UNCONTRACTABLE exact text law on v0.41.0 (iteration 198): interpolation
-- in ensures errors "unsupported application: $builtin.show"; std/string.concat
-- also errors in ensures. The BODY does encode (false ensures refutes).
-- Keep named inline text tests; reprobe the specification encoder on upgrade.
```

```text
-- UNCONTRACTABLE on v0.41.0 (iteration 198): exact text ensures is skipped
-- UNENCODABLE_TYPE: calls user function "renderRef" that is not SMT-encodable
-- in this context. Named inline tests pin both artifact references and delimiter.
```

```text
-- UNCONTRACTABLE on v0.41.0 (iteration 198): the full Applied/Denied law
-- is skipped, UNENCODABLE_TYPE: calls user function "applyRevision" that is
-- not SMT-encodable in this context, even with w.revision >= 0. Keep the
-- shared predicates; applyRevision is proven separately and the package smoke
-- exercises commit. Reprobe on compiler upgrade; see w-kernel-exports-carry-no-contract.
```

These exemptions are compiler/context-specific, not a claim that the mathematical functions are inherently uncontractable. A separate renderer probe using `std/string.concat` also errors in the ensures encoder (V6). Reprobe on compiler upgrades. The existing named renderer/cache tests and package commit smoke continue to run (V7); they are not substitutes for a universal commit proof.

## 2. Preserve shared-predicate anti-drift

The runtime bodies of all eight exports are unchanged (V8). `verify` calls `proposalMatchesWorld`; `commit` calls `commitAllowed`, which calls both shared identity predicates (V8). Compile order matters: contracting `proposalMatchesWorld` and `verificationMatchesProposal` first makes the downstream laws encodable (V4). Keep that shared architecture. Do not reopen the obsolete blanket Proposal-ADT explanation: the two Proposal identity proofs and their counterexamples now run (V4).

## 3. PUB011 ratchet

The prototype appends a bounded `pkg quality --json` invocation to `scripts/verify_world_package.sh`, already reached by `scripts/verify_ail.sh` from CI (V9). `scripts/check_world_contract_quality.py` validates schema, package identity, successful compilation, nonzero proven contracts and zero errors/counterexamples/skips, then requires `contracts.uncontracted_exports == UNCONTRACTED_EXPORTS` (currently 3). It deliberately does not gate unrelated badges.

An increase fails immediately. A decrease also fails until its baseline is lowered in the same patch, so improvement cannot leave reusable slack. CI passes the event name explicitly to `resolve_base(event_name, env)` via `--event`, with separate PR-base and push-before environment values; there is no `||` selection chain. The `pull_request` branch selects `pull_request.base.sha`. Pushes trigger only on `dev`, so the `push` branch selects `event.before`, the dev tip before that push. Full checkout history is retained. The named `workflow_dispatch` branch selects `HEAD^` and prints `contract quality gate: workflow_dispatch fallback to HEAD^`; this fallback is logged, never silent (V19). These are local workflow-event simulations, not a triggered Actions run.

The new value must be nonnegative and no greater than the previous value. Only a verified base commit with no checker file bootstraps at three. An invalid or all-zero ref exits 1 with `contract quality gate: cannot resolve ratchet base <ref>` and no traceback; resolution failure is deliberately caught, never treated as bootstrap (V20). Missing event fields and unsupported events refuse. Environment values select history, not the numeric threshold (source inspection V22). Like all repository policy, this cannot resist a patch that deletes the gate itself.

Controls (V4, V10): real added uncontracted export changes the field from 3 to 4 and fails; pristine package succeeds; removed field and zero proofs fail; count 2 with baseline 3 fails pending a baseline update; changing baseline from historical 2 to current 3 fails. The five new required identities also prevent replacing a proven export with an exemption while keeping the count constant. An intentional exemption replacement needs review of the three in-code reasons.

History controls in `scripts/test_world_contract_history.py` also run from the kernel harness (V19–V21). Each asserts its result; printing alone cannot pass it.

| Control | Mutation killed |
|---|---|
| Simulated PR payload extraction, distinct base/before SHAs; assert PR selects base SHA | Selecting push-before or HEAD^ for a PR |
| Push selects the distinct before SHA | Selecting PR base or HEAD^ for a push |
| Dispatch selects HEAD^ and exact fallback output | Wrong fallback ref or removing the log line |
| CI expressions and explicit `--event` invocation asserted | Miswiring either SHA field or reverting to the old selection chain |
| All-zero and invalid refs: assert rc=1, exact stderr and empty stdout | Bootstrap on resolution failure, or restoring an uncaught traceback |
| Real root commit, asserted to have no checker: rc=0 and `ratchet monotonic: 3 -> 3` | Refusing legitimate bootstrap or changing its value |

## 4. Release and published-0.1.0 consequence

The old published-record fixture records content `sha256:0c8c60616e592dc01891e8bbb59350786f242a2f79a9eb2c587ae8b0ca2e00b9` and tarball `sha256:44fc9fab7be710f09b84274f744445db41a79df77c71cc43c0e952d75d97c27f` (V11). The seeded prototype moved the local golden to content `sha256:2883217cc3e308c44e2e57e0d132091cb6358507af2b28fed1d04cdf3f705e3f`, tarball `sha256:3d68abc753a748144dbaa00482308b5dd1e0adb861c7fd5042a3725ac45141c4`, 10060 bytes instead of 9785; both interface hashes stayed unchanged (V12). The delivered packet includes the additional rejected-commit smoke arm and has the final digests in §8. This is a new artifact, not the published artifact with extra annotations.

The package gate freezes manifest version 0.1.0 and release kind `feature`, checks the dry-run identity, and compares the golden. `cmd/world-publish` also freezes version 0.1.0, recomputes the packet and refuses packet drift; its live path requires attended fences (V13). Updating the local golden therefore does **not** prove that the public registry serves these new bytes.

Reconcile first requires an undrifted local packet, then uses that packet's expected hashes, not a historical durable intent. For a present record, a hash mismatch resolves `conflict` (V13–V14). Consequently, with the old golden retained, changed sources hit the local drift fence; with this prototype's refreshed golden and the published 0.1.0 record, the expected outcome is `conflict`. This consequence is inferred from the measured comparator and fixture, not a claim to have probed the live registry this session. Preserve the original fixture and historical release notes; do not rewrite published history to make reconciliation green.

Recommendation: treat this as a proof-hardening candidate for a new release. Add a new versioned CHANGELOG section explaining five added proofs and three explicit exemptions, retain the historical 0.1.0 section, validate the chosen `[release] kind`, and update the manifest freeze, dry-run identity, publisher version/confirmation/approval scopes, golden and release tests together. Keep the prototype version unchanged until that release work is authorized. `docs/SELF_MOD_PUBLISH.md` marks the refreshed digest rows as an unpublished rehearsal (V15).

**Decision for Mark:** A — 0.1.1 proof-hardening release, retain measured `kind="feature"` for the added proof capability (**recommended**); B — 0.2.0 capability release with `kind="feature"`; while unanswered, the loop tests the local candidate, preserves published 0.1.0 history, and never publishes.

## 5. Milestones (compile order, each ≤150 changed implementation LOC)

| Milestone | Deliverable / approximate LOC | Acceptance observations | Mutation killed |
|---|---|---|---|
| M1 | two shared identity contracts and obsolete-comment removal, ~25 canonical LOC | Run pinned ai-check on contracts; read both named verified results; run V4 controls | replace either `sameRef` body with digest-only equality → named counterexample |
| M2 | commitAllowed, plan, verify laws; three exemptions, ~45 canonical LOC | Run V4; read five verified and five false-ensures counterexamples; read three exemption diagnostics | remove accepted check, misbind plan input, accept mismatch → named counterexamples |
| M3 | contract/quality control harness, ~115 Python LOC | Run `python3 scripts/test_kernel_contracts.py`; read all controls and final success | false postconditions, actual body bugs, added export, missing field, zero proof |
| M4 | floor inventory and package test skip pins, ~40 LOC plus generated projection | Regenerate projection; update BOTH REQUIRED_VERIFIED and EXACT_TOTAL_VERIFIED to 16; refresh golden and runbook digests; update pristine marker; run full gate and Go inventory/runbook tests | drop a required ensures → missing identity; stale golden → step 9 diff; stale marker → pristine test failure |
| M5 | PUB011 parser/history check and CI wiring, ~80 LOC | Run full gate; V4/V10 controls; CI compares previous baseline; observe 3 and 16 | added export → 4 rejected; increased baseline → monotonicity assertion |
| M6 | attended release preparation, separate ≤150 LOC slices after decision | Read new manifest/CHANGELOG; run dry-run/package/fence tests; compare against preserved 0.1.0 record | mismatched packet → drift; mismatched served digest → conflict |

M1–M5 are a single floor-raise landing unit: no intermediate commit may leave the six-file inventory inconsistent. The six sites are canonical modules, generated projection, verify_ail constants, ready-packet golden, SELF_MOD_PUBLISH digest rows and module_manifest_gate_test marker. The seventh affected gate is the package property's named skip allowance: this prototype measures nine no-generator skips, including the five newly contracted exports, and pins those names (V7). Do not treat them as executed properties. Static proofs and 40 named examples are the evidence.

M6 is not permission to publish. Its acceptance is inspection and dry-run observation. The controller commits; an attended operator separately authorizes any irreversible publication.

## 6. Conflict surface and limitations

Coordinate edits to all three canonical modules, projection, gate scripts, digest golden, runbook and pristine marker. No production Go behavior changes are proposed. The remaining compiler gaps belong upstream; this item records reproducible probes without copying predicates into callers. This prototype does not prove host persistence, world determinism in all contexts, renderer injectivity, or a full commit law.

### 6.1 Relation to `D-WORLD-37` (quorum r2 objection, reviewer fix (a) applied verbatim by the controller)

Quorum r2 (`claude-sonnet-5@claude-p`) asked that this section "name D-WORLD-37, state its current status, and show the accepted ensures clauses are policy-neutral." Measured by the controller at `12a0270`:

- **Status: RESOLVED, not open.** The ledger row reads `| D-WORLD-37 | RESOLVED |` (ruled A, attended; recorded in `9efe2f4`). `mission_decisions.sh --open` lists only `D-WORLD-38`, `-39`, `-40`.
- **Different surface.** `D-WORLD-37` governs the Go host's durable store: context-aware `Commit` in `host/store/store.go:874` (`func (s *Store) Commit(c Commit) error`), with finite budgets, a cancellation cutoff and reconciliation of uncertain outcomes. This item touches only the pure AILANG kernel functions `world/transitions.ail` `plan`/`verify`/`commit` and `world/contracts.ail`. No Go file changes.
- **Policy-neutral by construction.** Every accepted `ensures` states what the **unchanged** body already computes (V8: the diff changes declarations, `! {}` effect rows, contracts and comments only; no body changes). A contract that restates existing kernel semantics cannot pre-commit a store-durability policy. `D-WORLD-40`'s bound table (timeouts, the three-outcome contract) is likewise a host/store matter, and none of this item's ensures mention time, cancellation or durability.

So neither M2 nor any other milestone here is gated on `D-WORLD-37` or `D-WORLD-40`. Only M6 (the release) waits, and it waits on this doc's own release decision in §4.

## 7. Verification Log

Commands below run from the worktree. Every command uses this exact common prefix (including source-inspection commands):

```sh
export AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang
export PATH="$(dirname "$AILANG_BIN"):$PATH"
export TMPDIR="$PWD/.iteration198"
export GOCACHE="$PWD/.iteration198/go-cache"
```

`.iteration198` is disposable evidence/cache storage inside the worktree, not a delivery artifact. Create it with `mkdir -p .iteration198`. The retained controls create and remove their own scratch directory inside this worktree. Negative observations always accompany a known-positive identity, package or comparison arm in the same command.

| ID | Exact command after prefix | Observed output / interpretation |
|---|---|---|
| V0 | `cat design_docs/coding-standards.md; "$AILANG_BIN" prompt` | S1–S8 read; pinned reference explains module paths, effect rows, interpolation and match syntax. Read before editing `.ail`. |
| V1 | `"$AILANG_BIN" pkg quality --json packages/world-core` (before edits) | rc=0; schema `ailang.package-quality/v1`; compile.ok=true; style.exported_funcs=19; contracts verified=11, total=11, uncontracted_exports=8, errors=0, counterexample=0; smoke passed=true; tests notes `inline test blocks are not yet discovered in package mode (m-package-test-discovery)`; PUB011/PUB012/PUB020 badges. Compile/smoke are the positive controls for the undiscovered-tests report. |
| V2 | `"$AILANG_BIN" pkg quality --json packages/world-core` (after projection) | rc=0; verified=16, total=16, uncontracted_exports=3; errors=0, skipped=0, counterexample=0; compile.ok=true. |
| V3 | `cat scripts/build_world_package.sh; ./scripts/build_world_package.sh` | Allowlist types/contracts/transitions/logepoch; copy into stage and replace projection; `allowlisted modules: iterated=4 wc-l=4`; `projected 4 modules into packages/world-core/world`. V7 supplies byte-identity observation. |
| V4 | `python3 scripts/test_kernel_contracts.py` | rc=0. Each of proposalMatchesWorld, verificationMatchesProposal, commitAllowed, plan, verify: `verified; rc=0`, false-ensures `counterexample; rc=1`, body mutation `counterexample; rc=1`. Every ai-check invocation asserts another verified identity. Real added export: `uncontracted_exports=4, ratchet=3`; pristine `PUB011 ratchet passed: uncontracted_exports=3`. Missing field, zero proofs, lower count and raised baseline killed. |
| V5 | `python3 scripts/test_kernel_contracts.py` | renderRef law: rc=1, `encoding error: cannot encode ensures clause: let body: let value: unsupported application: $builtin.show([$tmp1])`; false law: counterexample. cacheKey: rc=0, skipped, `calls user function "renderRef" that is not SMT-encodable in this context`. commit: rc=0, skipped, `calls user function "applyRevision" that is not SMT-encodable in this context`. Both latter false laws also skipped. Other named identities verify in every invocation. |
| V6 | `python3 scripts/test_kernel_contracts.py` | Scratch-only replacement of renderRef ensures with `result == concat([r.algo, ":", r.digest])` plus std/string import: check.passed=true; renderRef error `unsupported application: std/string.concat([$tmp3])`; sameRef and servesEntry verified. No contract accepted from these trials. |
| V7 | `./scripts/verify_ail.sh` | Final rc=0: `16/16 required world/ identities verified across 11 module(s)`; `all 40 required named tests pass`; `4/4 projection hashes equal their canonical sources`; `exactly 9 skips, all no-generator, all named`; plan → verify → commit smoke completed; `world package gate PASSED: 9/9 steps performed non-zero work`; `PUB011 ratchet passed: uncontracted_exports=3`; `verify gate PASSED: 16 required identities verified, 40 named tests pass`. |
| V8 | `git diff -- world/contracts.ail world/transitions.ail world/logepoch.ail; rg -n 'sameRef\(\|proposalMatchesWorld\(\|commitAllowed\(\|verificationMatchesProposal\(' world/contracts.ail world/transitions.ail` | Diff changes declarations/contracts/comments only, preserving bodies and tests; positive call-site enumeration shows sameRef in both identity bodies, both predicates in commitAllowed, proposalMatchesWorld in verify, commitAllowed in commit. |
| V9 | `rg -n 'verify_ail|ratchet|fetch-depth' .github/workflows/ci.yml; tail -6 scripts/verify_world_package.sh; cat scripts/check_world_contract_quality.py` | CI calls `./scripts/verify_ail.sh`; new full-history monotonic step; package gate invokes bounded quality and checker; checker exact count 3, prior-baseline nonincrease and fail-closed schema/proof checks. |
| V10 | `python3 scripts/test_kernel_contracts.py; python3 scripts/check_world_contract_quality.py --base HEAD` | Controls rc=0; `ratchet monotonic: 3 -> 3`; simulated prior 2/current 3 reports `raised baseline killed: ratchet increase forbidden`. Verified HEAD bootstrap check rc=0, `ratchet monotonic: 3 -> 3`. |
| V11 | `sed -n '1,48p' host/broker/registry_reconcile_v2_test.go; cat host/broker/testdata/metadata_world_core_0.1.0.json` | cfg01 pins old content/tar hashes quoted in §4; public-record fixture identifies world/core 0.1.0 and matching interface hashes. This is a local recorded observation, not a fresh registry GET. |
| V12 | `./scripts/verify_ail.sh` (before golden refresh, paired with final V7 run) | Initial run reached step 9: `ready packet differs byte-for-byte from golden`, old/new JSON gives §4 hashes and 9785→10060 bytes; both interface hashes identical. Final V7 comparison passes using refreshed golden. |
| V13 | `sed -n '444,528p' cmd/world-publish/main.go; sed -n '58,92p' cmd/world-publish/fences.go; sed -n '297,337p' cmd/world-publish/fences.go; sed -n '345,390p' host/broker/registry_reconcile.go; cat scripts/verify_world_package.sh` | runReconcile calls requireUndriftedPacket then passes packet hashes; frozen version 0.1.0; phrase `publish world/core@0.1.0 irreversibly`; packet mismatch refuses `drift`; comparator sets `ReconcileConflict` on mismatched hashes and `ReconcileSucceededReconciled` on equality; manifest freeze kind=feature/version=0.1.0 and dry-run identity present. |
| V14 | `go test ./host/broker -run 'Test(ReconcilePublished010OverFourDigests\|ReconcileV2MismatchIsConflict\|ResolvePresentRefusalsP1P3)$' -count=1` | rc=0; `ok github.com/sunholo-data/ailang-world/host/broker 0.333s`; exact recorded metadata succeeds and mismatch arms produce conflict. |
| V15 | `sed -n '78,104p' docs/SELF_MOD_PUBLISH.md; go test ./host/runbook/...` | Rehearsal warning and new content/tar rows visible; rc=0, `ok github.com/sunholo-data/ailang-world/host/runbook` |
| V16 | `go vet ./...` | rc=0, no diagnostics; paired with compiling/running packages in V17 as positive tool-health control. |
| V17 | `go test ./host/verifygate/... ./host/pkgproj/... ./cmd/world-publish/... ./host/runbook/...` | rc=0; verifygate `ok` (54.663s), pkgproj `ok` (15.919s), world-publish `ok` (3.895s), runbook `ok` (1.717s). No sandbox-denied socket result occurred. |
| V18 | `"$AILANG_BIN" pkg quality --json packages/world-core; "$AILANG_BIN" test --format json world/` | Package tests files=0/passed=0 with explicit discovery-gap notes; directory mode returns 40 named tests, failed_tests=0. Known-positive inline-test discovery control paired with package-mode negative observation. |
| V19 | `sed -n '1,36p' .github/workflows/ci.yml; python3 scripts/test_world_contract_history.py` | rc=0. Workflow shows push branches `[dev]`, pull_request, workflow_dispatch, fetch-depth 0, explicit event and separate SHA expressions. Simulated payload extraction: `pull_request: base=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa; PASS`; `push: base=bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb; PASS`; `workflow_dispatch: base=HEAD^; PASS`; exact captured line `contract quality gate: workflow_dispatch fallback to HEAD^`; `CI event/base wiring: PASS`; `world contract history controls PASSED`. Function called directly with simulated event name and extracted env; no workflow triggered. |
| V20 | `python3 scripts/test_world_contract_history.py; python3 scripts/check_world_contract_quality.py --base 0000000000000000000000000000000000000000` | Harness rc=0: both all-zero and `nosuchref` assert rc=1, exact `contract quality gate: cannot resolve ratchet base <ref>`, empty stdout, no traceback. Positive historical arm: `pre-checker base: ratchet monotonic: 3 -> 3; rc=0; PASS`. Direct all-zero CLI invocation exits rc=1 and prints `contract quality gate: cannot resolve ratchet base 0000000000000000000000000000000000000000`, no traceback. |
| V21 | `./scripts/verify_ail.sh > .iteration198/revision-verify-ail.log 2>&1` then `tail -14 .iteration198/revision-verify-ail.log`; independently `python3 scripts/test_kernel_contracts.py > .iteration198/revision-kernel-controls.log 2>&1` then `cat .iteration198/revision-kernel-controls.log` | Revision reruns both rc=0 (exit status retained before reading logs). `world package gate PASSED: 9/9 steps performed non-zero work`; `PUB011 ratchet passed: uncontracted_exports=3`; `verify gate PASSED: 16 required identities verified, 40 named tests pass`. Kernel output includes prior compiler/mutation controls, `kernel contract controls PASSED`, raised-baseline refusal, all V19/V20 history arms and `world contract history controls PASSED`. No sandbox denial occurred. |
| V22 | `sed -n '44,100p' scripts/check_world_contract_quality.py; tail -3 scripts/test_kernel_contracts.py; git diff --check` | rc=0. Source shows named event branches, missing-base assertion, unsupported-event refusal, deliberate CalledProcessError conversion for resolution, verified-commit-only bootstrap, and CLI error handler. Harness ends with `history_controls()`. Diff check emits no diagnostics. |

## 8. Prototype manifest and results

The delivered M1–M5 candidate adds six asserted Q3 quality rejections and an unaccepted-verification package smoke arm. Package step 6 requires its `denied:verification contract failed` line; the `if true` commit mutant fails there. The smoke byte change moves the local ready packet to content `sha256:473517079249f3959b5aba9727f62b40130624c9915612a47736fdc1dae2a780`, tarball `sha256:d532d923dc0d681cb0782d292cdce1d60d2a41b987e7fcfe6c3fef854f089cb9`, 10127 bytes. The runbook digest rows match; both interface hashes remain unchanged. The full AILANG gate reports 16/16 identities, 40 named tests, package 9/9, and PUB011=3. M6 remains an attended release decision; no publish or git write was performed.

## 9. Quorum log

Author `codex:gpt-6-astra` (OpenAI benched). Seats: `claude-sonnet-5@claude-p`, `gemini-3-1-pro`, `oc-glm-5-3`; `oc-kimi-k3` as replacement. **`oc-glm-5-3` and `oc-kimi-k3` were ABSENT (`unreachable`) in both rounds**: the Ollama Cloud weekly gauge is at 100% (`ailang mission quota`, 2026-09-27). So every verdict below is from two external reviewers, not three.

- **r1 BLOCKED 2/2** (artifact `w-kernel-exports-carry-no-contract-2026-09-27T01-47-34Z.json`). Both objections were on ONE surface, the CI history step. sonnet said the `workflow_dispatch` → `HEAD^` fallback was never exercised; gemini said invalid-base refusal was untested. The controller measured both TRUE in substance: an invalid base exited rc=1 only through an uncaught `CalledProcessError` traceback. Revision (astra, 197 s): `resolve_base(event, env)`, a logged dispatch fallback, a deliberate refusal message, and six history controls (V19–V22).
- **r2 BLOCKED 1/1, gemini PASS.** sonnet said open ruling `D-WORLD-37` governs the same `commit`/`verify` functions. The controller measured the premise **FALSE** (§6.1: it is RESOLVED, and it is about Go `host/store.Store.Commit`). Per the **narrow-refinement carve-out**, the reviewer's own fix (a) was applied verbatim as §6.1. No objection disputed the direction.
