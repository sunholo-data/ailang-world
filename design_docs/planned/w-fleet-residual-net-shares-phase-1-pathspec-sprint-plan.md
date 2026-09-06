# Sprint plan — `w-fleet-residual-net-shares-phase-1-pathspec`

- **Design:** [`w-fleet-residual-net-shares-phase-1-pathspec.md`](w-fleet-residual-net-shares-phase-1-pathspec.md), quorum round 2: PROCEED, 3/3 present and pass.
- **Queue:** World iteration 160, row 64, clause-2, HARNESS.
- **Sprint ID:** `w-fleet-residual-net-shares-phase-1-pathspec`.
- **Duration:** ~0.2d.
- **Milestones:** exactly one, `M1`.
- **Risk:** low. The production change is diagnostic wording and an adjacent comment; most of the cost is the self-contained regression and mutation proof.
- **Dependency:** none.

## 1. Outcome and fixed boundary

Make `check_driver_fleet()` report exactly what its existing pathspec measures. A successful comparison must name:

1. the World-tracked comparison set under `tools/launchd` plus the exact file `scripts/mission_decisions.sh`;
2. the separately checked `REQUIRED_FLEET_PATHS`, rendered from the array rather than repeated as a second hard-coded inventory; and
3. the phase-3 residual boundary, including the statement that files outside it are **unenumerated (not zero)**.

Keep the two existing `git ls-tree ... HEAD -- tools/launchd scripts/mission_decisions.sh` pathspecs byte-for-byte in meaning. Keep phase-1 accounting, phase-2 required checks, phase-3 eligibility/counting, per-path warning, return codes, callers, ownership, and skip/refusal behavior unchanged. The required list remains a separate mechanism; do not merge it into the pathspec or infer a broader inventory from it.

This sprint does **not** discover, own, compare, certify, list, or count other fleet paths. `scripts/mission_decisions_v2.sh` is a synthetic outside-boundary witness, not a proposed fleet member. No shared fleet ownership/inventory is introduced.

## 2. Exact baseline and verification facts

### Repository state at planner entry

- Detached planner worktree: `/Users/voightkampff/dev/sunholo-data/.planner-wt-world-iter160`.
- Exact HEAD: `f8cd08871c43f8c4347307f5b02c4f4defb45b40` (`docs: design precise fleet residual scope reporting`).
- `git status --porcelain=v1 -uall`: empty before the two plan artifacts were written.
- Source baseline: `d81ac424460e181c2a36a2ed11894415edca8048`. The only `d81ac4…f8cd08` changes are the design and its two quorum records; `git diff d81ac4… f8cd08… -- scripts/verify_go.sh host/verifygate` is empty.
- Live `scripts/verify_go.sh` SHA-256: `4614363ed51839a8b8a8a190429838b6d15315c91792ac585f08f3ab781b444a`, matching design V15.
- Platform measured by planner: Darwin/arm64, Go `go1.26.6 darwin/arm64`, `/bin/bash` 3.2.57.

### Controller-supplied baseline facts

These were measured by the controller before this planning run and are recorded, not re-labelled as planner executions:

- pinned AILANG v0.30.0 `verify_ail.sh`: rc=0;
- Go 1.26.6 `go build ./...`: rc=0;
- Go 1.26.6 `go vet ./...`: rc=0;
- Go 1.26.6 full `go test ./...`: rc=0;
- isolated `verify_go.sh --driver-fleet-check`: rc=1 on the existing three FLEET-owned driver differences;
- full `verify_go.sh`: rc=1 on those same existing three FLEET-owned driver differences.

The last two reds are the honest fleet-drift baseline. Do not claim or manufacture a full-green verifier, absorb the fleet-owned files, or weaken the gate. Socket denials observed under a sandbox are **UNINFORMATIVE UNDER SANDBOX**.

### Planner measurements at `f8cd088…`

| Check | Observation |
|---|---|
| `bash -n scripts/verify_go.sh` | rc=0 |
| `go test ./host/verifygate -run '^$' -count=1 -timeout=60s` | rc=0; package compiles; `[no tests to run]` is expected for this compile fence |
| `go test ./host/verifygate -run '^TestDriverFleetResidualScope$' -count=1 -timeout=60s -v` | rc=0 but prints `warning: no tests to run`; expected pre-implementation red condition is the missing named RUN/PASS identity, not the process rc |
| Live synthetic baseline | rc=0, old broad success wording present, new checked-set/boundary disclosure absent |
| Live synthetic outside addition | rc=0, 3 matches, zero residual warnings, zero sibling mentions, new boundary disclosure absent |
| Live synthetic inside addition | rc=0, helper warning exactly once, old unqualified count-one summary present, new qualified summary absent |

The synthetic repositories used local identity, disabled hooks/signing and system/global Git config, cleared inherited `GIT_*`, `CI`, and `AILANG_BIN`, used a synthetic HOME and explicit synthetic `AILANG_FLEET_REPO`, and bounded every subprocess to 10 seconds. They were automatically removed. No real source, HOME, fleet, or Git metadata was mutated.

### Seven-day velocity signal

The bounded seven-day history contains 15 non-merge commits touching `scripts/verify_go.sh` or `host/verifygate`, with 2,074 additions and 162 deletions across four unique files. That is churn, not a defensible LOC/day rate. The closest landed comparator, row 54 (`14036ee`), changed two files by +213/-1. The present design estimates roughly 15 production lines plus one 140–180-line test. Keeping one ~0.2d milestone is consistent with both the design and recent harness cadence; splitting it would add bookkeeping without an independently landable unit.

## 3. Single milestone

### M1 — Honest residual boundary and durable synthetic regression (~0.2d)

**Landing files:**

- modify `scripts/verify_go.sh` only inside `check_driver_fleet()` reporting and adjacent phase-3 commentary;
- create `host/verifygate/driver_fleet_scope_gate_test.go` with one named top-level regression, `TestDriverFleetResidualScope`;
- update the design document only with implementation/verification evidence after the code and drills are complete.

No `.ail`, `tools/launchd`, charter, mission log, workflow, `go.mod`, `go.sum`, shared skill, fleet checkout, or HOME path changes. No new dependency. The controller handles GitHub import and all Git writes.

#### Work slices

| Slice | Work | Estimate |
|---|---|---:|
| T1 | Change the success, required-path, boundary, residual-summary wording and adjacent phase-3 comment. Render an empty required array as `(none)` using Bash-3.2-safe guarded expansion. | 0.04d |
| T2 | Add the self-contained Go fixture and its baseline, outside, inside, required-missing, and empty-required-list subtests. Reuse `copyGateFile`; all repositories and mutations live below `t.TempDir()`. | 0.08d |
| T3 | Run the compile fence, focused test, health controls, and production wording-reversion drill; restore by captured bytes and hashes. | 0.05d |
| T4 | Run scoped checks and append exact implementation evidence to the design. | 0.03d |

## 4. Named acceptance criteria

### AC1 — Baseline disclosure is complete

`TestDriverFleetResidualScope/baseline_disclosure` creates committed World and fleet trees containing identical minimal:

- `tools/launchd/control.sh`;
- `tools/launchd/lib/pin-root.sh`; and
- `scripts/mission_decisions.sh`.

Run the copied live script with `bash scripts/verify_go.sh --driver-fleet-check`. Require rc=0, a named `3 files match fleet HEAD` success, the checked-set clause, the phase-3 boundary/unenumerated clause, and exactly one line equal after whitespace trimming to:

```text
explicit required paths (checked separately): tools/launchd/lib/pin-root.sh
```

Reject `SKIPPED`, `AILANG_BIN is unset`, and zero-comparable output. **Expected pre-change red reason:** the current script returns rc=0 but lacks all new disclosures and still says `tracked copy is current (untracked fleet additions not certified)`.

### AC2 — Outside-boundary additions are unenumerated, not zero

`TestDriverFleetResidualScope/outside_boundary_unenumerated` commits only `scripts/mission_decisions_v2.sh` to the synthetic fleet. First require `git ls-tree -r --name-only HEAD` to contain it. Then require rc=0, 3 matches, zero per-path residual warnings, zero filename mentions, and the complete phase-3 boundary sentence containing `files outside this boundary are unenumerated (not zero), not certified by phase 3`.

This assertion is the primary row-64 regression killer. **Expected pre-change red reason:** comparison rc remains 0 and the sibling remains silent, but the complete boundary/unenumerated sentence is absent.

### AC3 — Same-boundary positive control still reports exactly one

`TestDriverFleetResidualScope/inside_boundary_positive_control` keeps the sibling addition and commits `tools/launchd/new-helper.sh`. Assert both additions are in fleet HEAD. Require rc=0, 3 matches, zero sibling mentions, exactly one preserved per-path warning naming the helper, exactly one summary saying `1 unclassified fleet-only paths within the phase-3 boundary not certified (see above)`, and the same boundary disclosure as AC2.

AC2 and AC3 run under the same top-level test invocation. **Expected pre-change red reason:** the per-path warning fires, but the current count-one summary lacks `within the phase-3 boundary` and the new disclosure is absent. The existing warning is the positive instrument-health control proving phase 3 actually enumerates its stated domain.

### AC4 — Separate required-path refusal remains intact

`TestDriverFleetResidualScope/required_path_missing` uses a fresh case with `tools/launchd/lib/pin-root.sh` committed in fleet but omitted from World, while the other two comparable paths remain equal. Require World HEAD to omit and fleet HEAD to contain the required file. Require rc=1, `REQUIRED fleet paths MISSING LOCALLY:`, and `tools/launchd/lib/pin-root.sh (REQUIRED by World, absent locally)`. Reject success, skip, `AILANG_BIN is unset`, zero-comparable, missing-in-fleet, difference, and accounting-failure outputs.

This is a compatibility/attribution control, not the row-64 regression killer. The design’s V15 already measured the live branch with these results.

### AC5 — The required-list renderer is healthy when the list is empty

`TestDriverFleetResidualScope/empty_required_list` edits only the verifier copy inside that subtest’s `t.TempDir()` so `REQUIRED_FLEET_PATHS` is empty. Require the edit to match exactly one array block, `bash -n` rc=0, comparison rc=0, 3 matches through phase 1, exactly one `explicit required paths (checked separately): (none)` line, and no `unbound variable` output.

This is a positive future-shape control demanded by the quorum handoff. It does not alter the production required list or merge it into either pathspec. The planner’s disposable prototype measured rc=0, one `(none)` line, 3 matches, and no unbound-variable error under Bash 3.2.

### AC6 — Production wording reversion kills the claim; scoped gates pass after restore

After the ordinary named test is green, capture the exact post-implementation bytes and SHA-256 of both landing source files before any drill. Copy the implementation to a disposable temporary checkout/repository. In that copy only, replace exactly one new production reporting block with the base success/summary block and remove the new required/scope disclosure lines; do not change pathspecs, comparison logic, fixtures, dispatch, assertions, or expected strings.

Run the ordinary named test against the copied mutant. Its process must go red specifically at `outside_boundary_unenumerated`: the underlying script command remains rc=0 with 3 matches, while the independent assertion fails because the phase-3 boundary/`unenumerated (not zero)` text is missing. A syntax error, missing fixture, timeout, skipped test, changed rc, or compile failure does not count. A committed mutation subtest may additionally expect this oracle failure, but it cannot replace the ordinary named-test red demonstration.

Restore from captured pre-drill bytes, never from HEAD. Recompute both hashes and require exact equality to their pre-drill values; rerun the focused test green with named RUN/PASS and no skips.

Planner prototype pins, measured only on the disposable exact edit described in §6:

- live pre-implementation verifier: `4614363ed51839a8b8a8a190429838b6d15315c91792ac585f08f3ab781b444a`;
- prototype post-implementation/pre-drill verifier: `8c997533a305491cf9cf2c6bf24c110bdd52b4dedac5059f9b89f1e90cc8a3ea`;
- prototype wording-reversion mutant: `0e63e575cfba94073b54d0ff7e1b66539ce61666e7145da020ad2c77c1fe0461`;
- restored prototype verifier: `8c997533a305491cf9cf2c6bf24c110bdd52b4dedac5059f9b89f1e90cc8a3ea`.

The executor must record its own exact post-implementation pre-drill hashes for `scripts/verify_go.sh` and `host/verifygate/driver_fleet_scope_gate_test.go`; the prototype pins are evidence, not permission to force a byte-identical implementation.

## 5. Test and compile fence

Run from the executor worktree with `PATH=/opt/homebrew/bin:$PATH`; bound the focused test to 60 seconds:

```bash
bash -n scripts/verify_go.sh
go test ./host/verifygate -run '^$' -count=1 -timeout=60s
go test ./host/verifygate -run '^TestDriverFleetResidualScope$' -count=1 -timeout=60s -v
```

The second command is the **test compilation fence**. The third is accepted only with `=== RUN   TestDriverFleetResidualScope`, named subtest RUN/PASS lines, final PASS, and no `warning: no tests to run` or SKIP. A zero exit with no named RUN is red.

Then run the repository’s scoped syntax/format/static checks appropriate to the two code files. Record the controller’s full-suite results separately; do not turn the known fleet-drift rc=1 into an implementation failure or a false full-green claim. The final evidence must say which command ran, its rc, named test output, mutation red reason, and both before/after restore hashes.

## 6. Planner reproduction of the load-bearing assertions

The planner copied the actual `scripts/verify_go.sh`, applied only the specified reporting/comment edit in a temporary synthetic World, and created fresh committed World/fleet repositories. Every subprocess had a 10-second timeout. Results:

| Arm | Measured prototype result |
|---|---|
| Baseline | rc=0; 3 matches; all three disclosures present; required member line exactly once |
| Outside addition | rc=0; 3 matches; zero sibling mentions; zero residual warnings; unenumerated disclosure present |
| Inside addition | rc=0; helper warning exactly once; qualified summary exactly once; zero sibling mentions; disclosure present |
| Missing required locally | rc=1; required FATAL and exact path present; success/skip/binary-gate wrong outputs absent |
| Empty required array | rc=0; 3 matches; `(none)` line exactly once; no unbound-variable output |
| Wording-reversion mutation | exact-one replacement; rc=0 and 3 old-style matches remain; boundary disclosure absent; AC2 oracle false for the intended reason |
| Restore | SHA-256 returned to `8c9975…a3ea`; outside arm rc=0 and disclosure present again |

This establishes that the proposed assertions discriminate the wording defect while the same-boundary warning and required-path refusal remain live. It is prototype evidence, not a claim that implementation or its durable Go test already exists.

## 7. Risks, residuals, and parked judgment

- Human-readable output has external consumers the tracked repository cannot enumerate. The design measured no tracked executable exact-text parser; row 64 authorizes the wording change. Do not claim universal absence.
- Exact punctuation may cause brittle tests. Assert the complete semantic clauses required by AC1–AC3, but do not pin indentation or dynamic fleet HEAD.
- The required array may later become empty. AC5 makes the Bash-3.2 behavior explicit now.
- Wider residual discovery and a shared fleet inventory remain deferred. Discovering evidence that the fixed pathspec itself must change is genuine new judgment: stop, record it in the report, and return it to the controller rather than expanding this sprint.
- Any new source file beyond the one named Go regression, any production behavior beyond reporting/comments, or any ownership change is out of scope.

No human confirmation is required for this already-authorized sprint. No message handoff or Git write is part of the planner lane; the controller imports and commits the two artifacts.
