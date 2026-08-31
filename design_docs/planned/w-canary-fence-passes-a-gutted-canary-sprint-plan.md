# Sprint plan — `w-canary-fence-passes-a-gutted-canary` (queue row 49)

**Design doc:** `design_docs/planned/w-canary-fence-passes-a-gutted-canary.md`
**Planning base:** `origin/dev` = `0bbb1a96603fa75279b8f9f55e9d2fe922fb6a2c` (`0bbb1a9`)
**Duration:** <=0.2 day (about 1.75 hours)
**Risk:** Low implementation risk; medium verification-discipline risk
**Persistent implementation surface:** one existing test file, approximately +85 to +100 / -3 lines
**Dependencies:** no code dependency; the ratified quorum-round-2 narrow-refinement carve-out is binding
**Milestones:** 3, all required

## Outcome and frozen scope

Replace the `strings.Count(src, "stateRoot") >= 2` positive-arm control in
`TestCanaryDeclaresPositiveArmOnly` with a Go AST shape fence. The fence must accept exactly one
`TestToolchainCanary` declaration containing exactly one canonical assertion as a **direct statement
of `TestToolchainCanary.Body.List`**:

```go
if rows[0].field != "stateRoot" {
	t.Fatalf(...)
}
```

The outer assertion must itself be direct in the function body; its `t.Fatalf` must be a direct
expression statement in that assertion's body. Recursive discovery of the assertion is forbidden.
That last constraint is M13, added by quorum round 2 under the ratified narrow-refinement carve-out.

Only `host/verifygate/toolchain_pin_gate_test.go` may be persistently modified. Add `fmt`, `go/ast`,
`go/parser`, and `go/token`; add `canaryAssertionShapeProblems`, `isRowsField`, `isStateRootLit`, and
`isTFatalfCall`; replace only the old `stateRoot` count clause. Do not modify production code,
`host/store/toolchain_canary_test.go`, `.github/workflows/ci.yml`, `run.sh`, `go.mod`, `scripts/`, or
any `.ail` file. The canary file is a temporary mutation venue only and must be restored after every
arm. The existing `GOTOOLCHAIN` zero-needle and `POSITIVE ARM ONLY` marker clauses remain
byte-unchanged.

## Planner corrections and clarifications

1. The design's blanket sentence that “M10/M11/M12/M13 are the quorum objection-1 additions” is
   inaccurate for M13. M10-M12 came from round 1 objection 1; M13 came from the quorum-round-2
   objection and is admitted only by the ratified narrow-refinement carve-out. This plan uses the
   latter provenance.
2. The design header still calls itself “round-1” even though its M13/V34 material records the
   round-2 refinement. Execution must use the final M13-bearing text, not the stale header label.
3. The design's `~50 LOC` estimate understates its own code sketch: the helper, three predicates,
   comments, imports, and call site are realistically about +85 to +100 / -3 changed lines. The
   schedule remains <=0.2 day because the code is local, stdlib-only test machinery.
4. The planning tree is not porcelain-empty because the supplied design doc is untracked. The
   tracked implementation baseline is nevertheless exact: `git diff --name-only origin/dev --`
   produced 0 lines before these two planning artifacts were created. Do not claim global
   porcelain 0 while the design/plan artifacts are present.
5. `.ailang/state/**` is ignored by `.gitignore`. The required sprint JSON is still created at the
   requested path; a controller that wants to commit it will need an explicit force-add. The
   executor must not perform that git write.

No premise, systemic claim, or acceptance criterion required redesign. The systemic audit counted
23 `strings.Count` controls across four verifygate files and classified this clause as the only one
guarding another file's assertion rather than a pin, marker, or identity. The local repair is
therefore sufficiently systemic for this sprint.

## Current status and recent velocity

Nothing from this design is implemented at the planning base: the old token-count clause is still
present and all four proposed helpers/imports are absent. The recent seven-day window represented
by `5fd6fb3..origin/dev` contains 57 commits, including 17 implementation-shaped commits, and
49 files with 17,810 insertions plus 917 deletions. That churn is documentation-heavy and is not a
safe LOC/day implementation forecast. The useful local signal is completion of the six adjacent
P42-P47 implementation items during 2026-08-27/28. This sprint is far smaller than those items;
its <=0.2-day duration is governed by thirteen serial mutation restorations and the full gates,
not by typing roughly 95 changed test lines.

## Baseline ledger at `origin/dev` `0bbb1a9`

Measured first-party in the assigned worktree on 2026-08-30 with zsh,
`PATH=/opt/homebrew/bin:$PATH`, `go1.26.6 darwin/arm64`, and
`/tmp/ailang-v0300/ailang` reporting `AILANG v0.30.0`:

| Gate | Exact command | Observed base result |
|---|---|---|
| B1 provenance | `git rev-parse HEAD && git rev-parse origin/dev && git diff --name-only origin/dev --` | both revisions `0bbb1a96603fa75279b8f9f55e9d2fe922fb6a2c`; tracked diff 0 lines |
| B2 formatting | `gofmt -l host/verifygate/ host/store/` | empty, rc=0 |
| B3 vet | `GOTOOLCHAIN=go1.26.6 go vet ./host/verifygate/ ./host/store/` | rc=0 |
| B4 build | `GOTOOLCHAIN=go1.26.6 go build ./...` | rc=0 |
| B5 named fence | `GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run '^TestCanaryDeclaresPositiveArmOnly$' -count=1 -v` | rc=0; exactly one `=== RUN`; one `--- PASS` |
| B6 vacuity control | `GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run '^TestNoSuchCanaryFenceZZZ$' -count=1 -v` | rc=0; `[no tests to run]`; zero `=== RUN` |
| B7 named canary | `GOTOOLCHAIN=go1.26.6 go test ./host/store/ -run '^TestToolchainCanary$' -count=1 -v` | rc=0; exactly one `=== RUN`; one `--- PASS` |
| B8 scoped packages | `AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ ./host/store/ -count=1` | rc=0; both packages `ok` |
| B9 full host gate | `AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.26.6 go test ./... -count=1` | rc=0; all packages `ok` |
| B10 AIL gate | `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` | rc=0; 11/11 identities, 40 named tests, world package 9/9 |

Base sha256 ledger:

| Path | SHA-256 | Role |
|---|---|---|
| `host/verifygate/toolchain_pin_gate_test.go` | `fc8efffc641b19e93897d9ddd925cb5faea6b82d239f3fbf96e323c3384f3614` | sole persistent implementation edit |
| `host/store/toolchain_canary_test.go` | `a23cfa79419ae69136e62981d5bd0c8ea68cdf2e154ba973fb99a3d3c8b47bfd` | read-only mutation venue |
| `.github/workflows/ci.yml` | `aed8e186fb57036eb6b03509cbb668850d577d46e6cc68b30e7a4c042108ed85` | forbidden |
| `design_docs/verification/w-race-gate-blindspot/run.sh` | `b80109aa57882a3d261757fd268888cf84ea0326180bb34969ce2ce9d149d0bb` | forbidden |
| `go.mod` | `7a2983617bb9fc33747f664564fe8d8ab54fc3a177ec4dfb8c61b29ba79a7e52` | forbidden |

## Milestones

### M1 — implement the exact AST recognizer (~0.08 day, ~80-90 test LOC)

In `host/verifygate/toolchain_pin_gate_test.go`, add the four stdlib imports and the four
package-private helpers from the design decision. `canaryAssertionShapeProblems` must:

1. parse the supplied source with `parser.ParseFile` and return a named `parse error` problem;
2. enumerate top-level declarations and require exactly one function named
   `TestToolchainCanary`;
3. enumerate **only** `funcDecl.Body.List` to find the canonical assertion—do not use
   `ast.Inspect`, recursive walking, or descendant collection for assertion discovery;
4. require exactly one direct `*ast.IfStmt` whose condition is a `*ast.BinaryExpr` with
   `token.NEQ`, left operand structurally `rows[0].field`, and right operand exactly the Go string
   literal `"stateRoot"`;
5. enumerate only that if-statement's `Body.List` and require exactly one direct expression
   statement calling selector `t.Fatalf`; and
6. return exact-count problems for zero and greater-than-one cases.

M1 gates:

```bash
export PATH=/opt/homebrew/bin:$PATH
gofmt -w host/verifygate/toolchain_pin_gate_test.go
test -z "$(gofmt -l host/verifygate/ host/store/)"
GOTOOLCHAIN=go1.26.6 go vet ./host/verifygate/ ./host/store/
GOTOOLCHAIN=go1.26.6 go build ./...
```

Acceptance: all commands rc=0; no dependency outside the standard library; helper names are unique;
source review confirms assertion discovery iterates `funcs[0].Body.List` directly. M13 cannot be
deferred to the mutation drill: the implementation structure that kills it is part of M1.

### M2 — replace the token control and prove the pristine positive arm (~0.04 day, ~8-10 test LOC)

Replace only the old `strings.Count(src, "stateRoot") < 2` clause with a call to
`canaryAssertionShapeProblems(src)`. Report each named problem and terminate the test with the
design's instrument-failure message. Keep the `GOTOOLCHAIN` and `POSITIVE ARM ONLY` clauses exactly
as they were.

M2 gates:

```bash
export PATH=/opt/homebrew/bin:$PATH
GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ \
  -run '^TestCanaryDeclaresPositiveArmOnly$' -count=1 -v
GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ \
  -run '^TestNoSuchCanaryFenceZZZ$' -count=1 -v
GOTOOLCHAIN=go1.26.6 go test ./host/store/ \
  -run '^TestToolchainCanary$' -count=1 -v
AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.26.6 \
  go test ./host/verifygate/ ./host/store/ -count=1
```

Acceptance: the two real named runs each show exactly one `=== RUN` and one `--- PASS`; the nonsense
selector shows zero `=== RUN` and `[no tests to run]`; both scoped packages pass; the deliverable
file contains zero `strings.Count(src, "stateRoot")` calls; and the unchanged retained checks still
contain exactly these conditions:

```go
if count := strings.Count(src, "GOTOOLCHAIN"); count != 0 {
if !strings.Contains(src, "POSITIVE ARM ONLY") {
```

### M3 — discharge M1-M13 and run the landing gates (~0.08 day, 0 persistent LOC)

Create one backup of the **post-M2** canary in `/tmp/w-canary-fence-backup/`, record its SHA-256,
and use copies from that backup for every arm. Never restore with `git checkout`, `git restore`, or
`git reset`: those commands risk deleting the uncommitted implementation. After every arm, copy the
backup back, require SHA-256 equality with the post-M2 canary, and re-run the pristine named fence.
At every milestone boundary, `git diff --name-only origin/dev --` may name only
`host/verifygate/toolchain_pin_gate_test.go`; planning artifacts are untracked/ignored state, not
implementation scope.

Common mutation gate (all M1-M13):

```bash
export PATH=/opt/homebrew/bin:$PATH
GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ \
  -run '^TestCanaryDeclaresPositiveArmOnly$' -count=1 -v
```

Every arm must return rc!=0, print exactly one top-level `=== RUN`, print `--- FAIL`, and contain the
row's named problem text. A package-wide red or a compile error in `host/store` does not substitute
for the named fence red. M1, M2, M3, M4, M5, M6, M7, M8, and M9 prove the original four assertion
floors; M10-M12 bind quorum-round-1 objection 1; M13 alone binds the quorum-round-2 direct-body
refinement.

## Mandatory mutation matrix

| Arm | Exact temporary edit to `host/store/toolchain_canary_test.go` | Required named fence result | Additional requirement |
|---|---|---|---|
| M1 | Replace the canonical assertion's two lines with exactly `// stateRoot is the expected field value.` | rc!=0; problem contains ``top-level `rows[0].field != \"stateRoot\"` assertion if-stmt count=0, want exactly 1`` | `go vet ./host/store/` rc=0; named canary rc=0; old token facts remain `stateRoot`=2, `GOTOOLCHAIN`=0, marker=1 |
| M2 | Replace the same assertion with exactly `// The stateRoot token remains as prose; no executable check remains.` | same canonical assertion count=0 problem | named canary rc=0; comment/prose cannot satisfy the AST fence; old token floor remains green |
| M3 | Rename only `func TestToolchainCanary` to `func TestToolchainCanaryRenamed` | problem contains `TestToolchainCanary func decl count=0, want exactly 1` | `go vet ./host/store/` rc=0; proves exact target-function anti-vacuity |
| M4 | Keep the canonical `if`, but replace its direct `t.Fatalf(...)` with `_ = rows[0].field` | problem contains `direct t.Fatalf expression statement in assertion body count=0, want exactly 1` | source parses/builds; comparison alone is insufficient |
| M5 | Duplicate the entire canonical assertion immediately after itself | problem contains canonical assertion `count=2, want exactly 1` | source parses/builds; ambiguous shapes red loudly |
| M6 | Change only the canonical condition's `!=` to `==` | canonical assertion count=0 problem | fence red is required even though the mutated canary itself also reds on the good toolchain |
| M7 | Leave an empty canonical `if` body and move the same `t.Fatalf(...)` to the next direct function-body statement | direct-Fatalf count=0 problem | a stray function-level Fatalf cannot satisfy the assertion action; do not use the canary's runtime red as evidence |
| M8 | Delete the entire `TestToolchainCanary` declaration and its body; leave `canaryString` intact | function declaration count=0 problem | stronger deletion form of the function-existence floor |
| M9 | Append the exact incomplete declaration `func m9ParseFloor(` at EOF | problem contains `parse error:` | the verifygate named test must red from its parser floor; a `host/store` compile failure is expected but not credited |
| M10 | Replace the canonical condition with `"stateRoot" != "stateRoot"`; keep the direct Fatalf body | canonical assertion count=0 problem | always-false constant comparison rejected because the left operand is not `rows[0].field` |
| M11 | Replace the canonical condition with `rows[0].n != "stateRoot"`; keep the direct Fatalf body | canonical assertion count=0 problem | wrong selector rejected; `host/store` need not compile because `n` and string are incomparable—the named AST fence red is the evidence |
| M12 | Replace the canonical body with `if false { t.Fatalf("Field=%q want %q", rows[0].field, "stateRoot") }` | direct-Fatalf count=0 problem | source must parse/build; descendant Fatalf is not direct; pristine named canary remains green |
| M13 | Replace the direct canonical assertion with `if false { if rows[0].field != "stateRoot" { t.Fatalf("Field=%q want %q", rows[0].field, "stateRoot") } }` | canonical **top-level** assertion count=0 problem | source must parse/build; old token facts remain green and named canary passes; proves only a direct `TestToolchainCanary.Body.List` statement counts |

For M1, M2, M10, M12, and M13, explicitly record the old-control counts (`stateRoot >= 2`,
`GOTOOLCHAIN == 0`, marker present) before running the new fence. These are the threat-shaped arms
where proving the old test stayed green makes the new fence's contribution non-vacuous. For all
arms, require the restored SHA-256 to equal the post-M2 backup and the pristine named fence to pass
before starting the next arm.

## Final landing gates

Run in this order after all mutations are restored:

```bash
export PATH=/opt/homebrew/bin:$PATH
test -z "$(gofmt -l host/verifygate/ host/store/)"
GOTOOLCHAIN=go1.26.6 go vet ./host/verifygate/ ./host/store/
GOTOOLCHAIN=go1.26.6 go build ./...
GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ \
  -run '^TestCanaryDeclaresPositiveArmOnly$' -count=1 -v
GOTOOLCHAIN=go1.26.6 go test ./host/store/ \
  -run '^TestToolchainCanary$' -count=1 -v
AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.26.6 \
  go test ./host/verifygate/ ./host/store/ -count=1
AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.26.6 \
  go test ./... -count=1
AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh
```

Required results: empty gofmt output; vet/build rc=0; named tests each run exactly once and pass;
scoped and full suites rc=0; verify gate reports 11/11 identities, 40 named tests, and world package
9/9. Confirm the canary, CI workflow, `run.sh`, and `go.mod` equal their base SHA-256 values. The
only persistent implementation diff is `host/verifygate/toolchain_pin_gate_test.go`.

## Success metrics and handoff boundary

## Executor completion record — 2026-08-30

- [x] **M1** — exact non-recursive AST recognizer added; gofmt empty, scoped vet rc=0,
  `go build ./...` rc=0; four helper definitions, one direct `funcs[0].Body.List` iteration,
  and zero `ast.Inspect` calls. Commit `ca2ecd6`.
- [x] **M2** — old `stateRoot` token clause replaced; named fence and named canary each produced
  exactly one RUN and one PASS; the nonsense selector produced zero RUNs plus the no-tests warning;
  two scoped packages passed; the retired clause count is zero and each retained clause count is
  one. Commit `345a73a`.
- [x] **M3** — M1-M13 each changed the canary bytes, returned non-zero from exactly one named
  fence run, printed one FAIL and the row-specific problem, then restored by copy to
  `a23cfa79419ae69136e62981d5bd0c8ea68cdf2e154ba973fb99a3d3c8b47bfd`; the pristine named fence
  produced one RUN and one PASS after every restore. Survivors: **0 of 13**.
- [x] **Landing gates** — gofmt empty; scoped vet and full build rc=0; named fence and canary each
  one RUN/one PASS; **2** scoped packages and **19** full-suite packages passed; `verify_ail.sh`
  reported **11/11** identities, **40** named tests, and world package **9/9**.
- [x] **Scope and restore** — canary, CI workflow, runtime `run.sh`, and `go.mod` match the base
  SHA-256 ledger. The only persistent implementation diff is
  `host/verifygate/toolchain_pin_gate_test.go`.

**Measured deviation:** M8's two requirements are jointly impossible as written. Deleting the
entire `TestToolchainCanary` declaration leaves the file's sole `testing` import unused, so the
first exact arm returned `go vet ./host/store/` rc=1 with `"testing" imported and not used`.
The file restored byte-identically. The successful temporary M8 mutant deleted the target function
and its now-unused import; it then vetted rc=0 and the named fence failed with function-declaration
count 0. This is a temporary mutation-construction correction, not a persistent-scope or mechanism
change.

- One test file modified; no production or AILANG source changed.
- Structural fence requires the exact canonical comparison and direct Fatalf, exactly once.
- M1-M13 all red the named fence with their specified component message; no survivor.
- M13 proves a canonical assertion nested beneath `if false` is rejected.
- All pristine, hygiene, scoped, full-host, and AIL verification gates pass.
- No documentation/example file is required: this is a static host gate over an existing canary.
- Executor updates only each sprint JSON milestone's `passes`, `started`, `completed`, and `notes`
  fields. No commit, mission-record edit, or external handoff is part of this planning task.
