# Sprint plan — `TR.C` (queue item 11 `w-transition-registry`)

**Milestone**: `TR.C` — *the binding gate*. The LAST milestone of item 11.
**Status**: PLANNED · **NO SPLIT** · **TR.C ONLY** (TR.A and TR.B are landed and out of scope)
**Design doc**: [`design_docs/planned/w-transition-registry.md`](w-transition-registry.md)
**Sibling plans**: [`…-tra-sprint-plan.md`](w-transition-registry-tra-sprint-plan.md) (TR.A1+TR.A2, landed),
[`…-trb-sprint-plan.md`](w-transition-registry-trb-sprint-plan.md) (TR.B1+TR.B2, landed)
**Base**: `dev` @ `2d5a346`, clean tree, CI green both jobs
**Planner**: mission-control iteration 75, opus lane (`derive-planner-lane.sh` absent → fail-closed to
opus, reason token `missing-script`), first-party measurement on this rig
**Executor**: sandboxed worktree, **NO git write permission**.
**THE CONTROLLER MAKES ALL COMMITS.** The executor never runs `git commit`, `git add`, `git push`,
`git checkout`, `git stash`, `git restore`, or `gh pr`. **Restores are `cp` from
`/tmp/trc_backup/`**, never `git checkout -- <file>`: in a sprint worktree the mutated file is
uncommitted by construction, so `git checkout --` deletes the executor's work and the sha256
assertion then reports the disaster rather than preventing it.

**Headline price: ≈1.25 days, i.e. 2.5× the doc's 0.5 day.** The driver is not LOC (≈760); it is
**23 mutations / 46 arms**. §6 prices it and §6.1 explains why a split is the wrong answer here.

---

## 0. Planner's first-party verification of every premise

Every controller premise was re-measured at `2d5a346` before anything was designed. **All seven
survived.** Nothing in the controller's brief is refuted; two rows are *sharpened* (§0.1). Six
measurements the controller did not make are in §0.2, and five of them change the plan.

The full command/observed-output table is §8 (Verification Log). This section states only the
conclusions.

### 0.1 Controller premises — confirmed, with two sharpenings

| Premise | Verdict |
|---|---|
| `.Invoke(` production sites = exactly 3, at `publish_op.go:135,162,279`; test control 90 | **CONFIRMED** (V1) |
| exported `broker.NewSession` has zero production callers | **CONFIRMED, sharpened** — see below |
| `broker.Session`/`broker.NewSession` outside `host/broker` in production = 0; control 55 | **CONFIRMED** (V3) |
| AC11 at base `count=1`; `invoke_boundary_test.go` absent; AC5 control on the same package = 2 | **CONFIRMED** (V4, V5) |
| hold set AC5=2, AC6=3, AC7=3 | **CONFIRMED** (V5) |
| freshness sweep `b0f323a`→HEAD ex-`design_docs/` = 19 files, 4 of them in `host/broker/`, control 31 | **CONFIRMED** (V6) |
| three existing AST-based tests; the doc's placement wins unless a concrete defect is shown | **CONFIRMED — and a concrete defect in the ALTERNATIVE was measured** (§0.2(vii), V7) |

**Sharpening 1.** "The only production hit is its own definition at `broker.go:87`" is true only for
the *package-qualified exported* reading. The literal grep `NewSession(` also returns
`func NewReplaySession(` (`broker.go:92`) and three calls to the **unexported** `newSession(` at
`publish_op.go:132,159,276`. That is not a refutation — it is the firing positive control for the
instrument, and it matters for TR.C because **`newSession` is unexported and therefore
unreachable from outside `host/broker` by construction**, which is why TR.C only has to gate the two
*exported* constructors (V2).

**Sharpening 2.** The design doc's own V25 row records the `_test.go` control at **83**. It is now
**90** (V1). TR.B added test-side `.Invoke(` calls. The doc row is stale, not wrong; no action.

### 0.2 Seven measurements the controller did not make

#### (i) The doc's `MUT-BINDING-FOURTH-INVOKE` is ONE named mutation covering TWO distinct refusal branches. This is iteration 74's defect shape, verbatim.

A fourth production `Session.Invoke`:

- **inside** `host/broker` trips the **exemption-count** assertion (3 → 4);
- **outside** `host/broker` trips the **outside-is-clean** assertion.

These are different code paths in the detector, different assertions, and different failure
messages. One named mutation exercises exactly one of them and the plan-reader cannot tell which.
Iteration 74 shipped `equalRequirements`' second branch uncovered for precisely this reason
("doc, plan and executor each named one mutation"). **Split into `MUT-BINDING-FOURTH-INVOKE-IN`
(M1) and `MUT-BINDING-FOURTH-INVOKE-OUT` (M2).**

#### (ii) `MUT-BINDING-IF-FALSE` is INVISIBLE for five of the six detector branches, and the doc's stated RED observation cannot see it.

The doc's required observation is *"known three-site detector control sees zero instead of exactly
3"*. That is an observation about the **inside-broker** detector only.

**Measured (V8, V9): the repository is clean outside `host/broker` — the outside-broker detector
returns ZERO findings at base.** Therefore neutering *any* outside-broker branch with `if false &&`
changes nothing observable: 0 findings before, 0 findings after. A green suite containing a fully
disarmed outside-broker detector is indistinguishable from a green suite containing a working one.

**This is the single most important gap in TR.C's design**, and it is exactly the class this item
keeps shipping. The fix is structural, not additional vigilance: the committed test must carry a
**hermetic synthetic control suite** — a known-positive source per detector branch and a set of
known-negative sources — driven through the *same* detector core the repo walk uses. Then
`if false &&` on any branch reds its own control. §4 mandates one `if false` arm **per branch**
(M11–M15), each with a named control subtest as its RED observable.

#### (iii) The doc's "any `Invoke` selector call" is ambiguous, and one plausible reading is RED AT BASE.

`cmd/world-publish/main.go:367` calls **`broker.InvokeAttendedPublish(...)`** (V10). A detector that
matches the selector by *prefix* or *substring* flags it and TR.C reds on the pristine tree — a gate
that measures the repository rather than the change (rule 3e).

`InvokeAttendedPublish` is the sanctioned SM.D0 attended entrypoint; it is a package-level function
that constructs its session internally. It is **not** a `Session.Invoke` call and must not be gated.

**Freeze**: the detector matches on **exact equality of the selector identifier** with the assembled
token `"Invo"+"ke"`. `NEG-invokeattended` is a committed known-negative control (V11 row 7): measured
`findings=0`.

#### (iv) The doc requires the exemption IDENTITIES to be asserted but names no mutation that moves an identity while holding the count.

`MUT-BINDING-RAISE-COUNT` moves the count. Nothing in the doc's table moves an identity. A detector
that asserts `count==3` and *logs* the identities is green under an attacker who relocates one of the
three exempt calls into a new helper — the count is still 3. **This is the count-vs-identity mirror of
"guard the helper, miss the call site".** Added as **M3 `MUT-BINDING-MOVE-SITE`**, plus **M10
`MUT-BINDING-COUNT-INEQUALITY`**, which weakens the count assertion to `< 3` and *composes* with M1:
the identity assertion must still red, or it is decorative.

#### (v) The file enumeration must be a filesystem walk, not `go list` — and the reason is measured, not stylistic.

| enumeration | files | what it misses |
|---|---:|---|
| filesystem walk from module root, skipping nested modules and `_test.go` | **39** | — |
| `go list ./...` `.GoFiles` | **38** | `host/store/writer_lock_other.go` |

`writer_lock_other.go` carries `//go:build !unix` (V13). `go list` never reports it on darwin **or**
on linux/CI, so a `go list`-driven enumeration leaves a real production file **permanently
unscanned** — a perfect place to park a `broker.Session` reference. `go/parser` does not evaluate
build constraints, so the walk reads it. `go list` is retained only as a **superset cross-check**
(walk ⊇ go list), which is the one instrument that is *not* a hand-maintained list and that catches
"the walk silently missed a whole new package".

#### (vi) `design_docs/verification/w-race-gate-blindspot/{racecontrol,repro}` are SEPARATE Go modules.

Each has its own `go.mod` (V12). They are out-of-module scratch reproducers containing
`package main`. A naive repo-root walk parses them. They are clean today (V8) — but they are
uncontrolled scratch space, and a gate whose green depends on scratch files staying clean is a gate
with a foreign dependency. **The walk skips any directory containing its own `go.mod`.** Measured:
`skipped_nested_modules=2`.

Because the skip is *not* observable against the live repo today (removing it keeps the gate green —
those two files contain no broker reference), it is proven by a **hermetic synthetic control**: a
`t.TempDir()` tree with a nested `go.mod`, asserting the walker does not descend. A branch whose only
proof would be a live mutation that cannot red is a branch that must be proven synthetically.

#### (vii) The alternative placement (`host/boundary`) is a MEASURED defect. The doc's placement stands.

The brief asked for a justified decision. `host/boundary` is semantically the natural home — it is
the package whose doc comment says it *"holds executable tests for dependencies that cross host
package boundaries"*, and its `repoRoot`/overlay/mutation machinery could be reused in-package.

**It is nonetheless disqualified, and the reason is measured, not aesthetic.** The landed
`TestBoundaryASTWriteGuard` pins `const wantFileCount = 1` and enumerates **every** `.go` file in the
directory, `_test.go` included (`allowlist_world_test.go:1157,1163-1166`). Adding any second file to
`host/boundary` reds it. Probe (V7, artifact removed, `git status --porcelain` empty afterwards):

```
--- FAIL: TestBoundaryASTWriteGuard (0.00s)
    allowlist_world_test.go:1165: host/boundary contains 2 .go files, want 1:
      ["allowlist_world_test.go" "zz_trc_placement_probe_test.go"]
```

Placing TR.C in `host/boundary` would **move a landed criterion**, which §1 forbids absolutely.
And the reuse argument is weaker than it looks: `host/boundary`'s helpers are unexported and live in
a `_test.go` file, so they are unimportable from any other package — TR.C in `host/broker` has no
reuse to forgo, because there is none available across the package line either way.

**Decision: `host/broker/invoke_boundary_test.go`, `package broker`, exactly as the doc's Files table
and AC11 say.** No divergence to flag; the doc was right, and now for a reason a reviewer can re-run.

---

## 1. What TR.C is, and what this plan explicitly does NOT do

**In scope** — three files:

| File | Purpose |
|---|---|
| `host/broker/invoke_boundary_test.go` | **NEW.** The whole milestone. `package broker`. Repository-wide `go/parser`/`go/ast` binding assertion + exact three-call legacy exemption + hermetic synthetic control suite. ~760 LOC. |
| `design_docs/planned/w-transition-registry.md` | AC11 zero-tolerance activation (delete the `count -eq 1` arm) + the four doc/plan divergences of §2. ~15 lines. |
| `design_docs/verification/w-transition-registry/trc-mutations.md` | **NEW.** The 23-mutation / 46-arm transcript, same form as `trb2-mutations.md`. |

**ZERO production LOC.** TR.C changes no dispatch: V25–V27 of the design doc, re-measured here as
V1–V3, show there is no coordinator/daemon→broker path to change. P6.B builds the first
registry-mediated path *under* this boundary.

**Explicitly NOT in this sprint** — if the executor touches any of these, the milestone is wrong:

1. **No `.ail` change.** `verify_ail.sh` totals stay **4 / 11 / 14 UNMOVED**.
2. **Do not touch `scripts/verify_ail.sh`** — queue item 12, deliberately out of scope.
3. **Do not repair the `host/broker` base flake** (`TestHandlerTimeoutKillsTheWholeProcessGroup`,
   ~18%) — queue item 16. Skip it; never silence it.
4. **No store table, no DDL, no REST route, no CLI verb, no projection package, no MCP/A2A.**
5. **No production `.go` edit that survives the sprint.** Mutation arms edit production files; every
   one is restored from `/tmp/trc_backup/` and the restore is asserted by sha256.
6. **TR.C MUST NOT MOVE ANY EARLIER MILESTONE'S CRITERION.** AC1=3, AC2=3, AC3=4, AC4=2, AC5=2,
   AC6=3, AC7=3, AC8=2, AC10 build+count=1, `AC-INVOKE3` n=3/p=3, `AC-VET` rc=0. AC11 alone moves,
   1 → 2. **If your change adds or moves a production `Invoke` call site, TR.C reds itself.**
   TR.C adds no production code, so the only realistic way to breach this is to forget to restore a
   mutant — which is why §4's protocol asserts the restore.

---

## 2. AC reconciliation and the four doc/plan divergences

TR.C closes exactly one criterion, **AC11**, and activates it. Everything else is a hold.

### The four divergences found (reported, not smoothed over)

| # | Divergence | Resolution |
|---|---|---|
| **D1** | `MUT-BINDING-FOURTH-INVOKE` is one mutation for two branches (§0.2(i)). | Split into **M1** (in-broker, trips the count) and **M2** (outside, trips outside-is-clean). Doc row amended at T6. |
| **D2** | The doc gates *"either exported session constructor"* — **two** constructors, two detector branches — and its mutation table names **no** mutation for either. | Added **M4 `MUT-BINDING-CTOR-LIVE`** and **M5 `MUT-BINDING-CTOR-REPLAY`**. Same shape as D1: one phrase, two branches. |
| **D3** | `MUT-BINDING-IF-FALSE`'s stated RED observation only covers the inside-broker detector; five outside-broker branches are unobservable at base (§0.2(ii)). | One `if false` arm per branch, **M11–M15**, each red-ing a named hermetic control subtest. |
| **D4** | *"any `Invoke` selector call"* is ambiguous; a prefix/substring reading is RED at base on `broker.InvokeAttendedPublish` (§0.2(iii)). | Frozen as exact selector-identifier equality; `NEG-invokeattended` committed as a known-negative control. |

Two smaller notes, recorded and fixed in the same T6 doc edit:

- **D5** The doc's Files table omits the mutation transcript and the AC11 activation edit. TR.A and
  TR.B both delivered both; TR.C does too.
- **D6** The design doc's V25 `_test.go` control reads 83; it is now 90 (§0.1 sharpening 2).

### AC11 — base and delivered

| | value |
|---|---|
| **Base at `2d5a346`** | `count=1`, rc=0. The single name is the known control `TestReplayReturnsRecordedBytesWithoutDispatch`; `TestRegistryDispatchBindingBoundary` is absent, `host/broker/invoke_boundary_test.go` does not exist. (V4) |
| **Same-call known-positive** | AC5's name-set on the **same package** with the **same `-list` instrument** returns **2** (V5). The instrument fires. |
| **Delivered** | `count=2`, both tests PASS, and the base-tolerant `test "$count" -eq 1 ||` arm is **deleted** (T6). |

---

## 3. Task breakdown — 6 tasks, 6 commits, one milestone

Every task exits on the **full hold set** (§5) plus its own gate. Every task's commands carry
`export PATH=/opt/homebrew/bin:$PATH` and `GOTOOLCHAIN=go1.25.6`.

### T1 — the production-file walker and its anti-vacuity assertions (~150 test LOC)

Deliver `enumerateProductionGoFiles(root) ([]string, walkStats, error)` inside
`invoke_boundary_test.go`, plus the `enumeration` subtest.

**Walker contract** (frozen):
- root = module root, located by `runtime.Caller(0)` → `filepath.Join(dir, "..", "..")`, never cwd;
- skip `.git`;
- skip any directory **other than root** that contains its own `go.mod`, counting the skips;
- skip `*_test.go`, counting the skips;
- collect every remaining `*.go` as a repo-relative slash path; sort.

**`enumeration` subtest — five independent anti-vacuity assertions, each its own `t.Fatalf`:**

| # | assertion | base value | mutation it exists for |
|---|---|---|---|
| E1 | `len(files) != 0` — *fail loudly*, never "nothing to check, pass" | 39 | M16 |
| E2 | `len(files) >= 30` (floor; measured 39) | 39 | M17 |
| E3 | every path in `requiredAnchors` was walked | 7/7 | M17 |
| E4 | no walked path ends `_test.go`, **and** `skippedTests > 0` | 0 / 44 | M18 |
| E5 | walked ⊇ `go list ./...` `.GoFiles`, and `len(goList) >= 30` | 39 ⊇ 38 | M19 |

`requiredAnchors` is a hand-maintained list of **seven** paths and that is deliberate — it is an
assertion *about* the walk, not the enumeration itself (which the doc requires be unmaintained):
`host/broker/publish_op.go`, `host/broker/broker.go`, `host/broker/confined.go`,
`host/transitionreg/bind.go`, `cmd/world-publish/main.go`, `cmd/ailang-worldd/main.go`, and
**`host/store/writer_lock_other.go`** — the last one included precisely so that a future rewrite of
the enumeration onto `go list` reds here (§0.2(v)).

E5's `go list` invocation mirrors `host/boundary`'s landed `goListDeps` helper: `exec.LookPath("go")`,
`exec.CommandContext` with a 90 s timeout, `cmd.Dir = root`. CI precedent exists and is green.

**Exit gate**: `enumeration` subtest PASSes; `t.Logf` prints `walked=39 skipped_tests=44
skipped_nested_modules=2 golist=38`. Hold set green.

### T2 — the detector core, assembled tokens, and import resolution (~200 test LOC)

Deliver `scanFile(rel string, tree *ast.File, fset *token.FileSet, insideBroker bool) []finding`.

**Assembled tokens** (the doc's constraint — no literal scan needle anywhere in the file):

```go
invokeSel  := "Invo" + "ke"
sessType   := "Sess" + "ion"
ctorLive   := "New" + "Session"
ctorReplay := "New" + "Replay" + "Session"
brokerPath := "github.com/sunholo-data/ailang-world/host/" + "broker"
```

**Import resolution** — per file, resolve the *local* name bound to `brokerPath` from
`tree.Imports`: default `broker`; explicit alias → that alias; `_` → not referenceable; `.` → a
finding in its own right (a dot-import makes `Session` a bare identifier and defeats selector
matching entirely — that is an evasion, not a style choice).

**Six detector branches** (this list IS the rule-3j enumeration; it is anchored to the diff TR.C
writes, not to the doc's decision list):

| id | branch | matcher |
|---|---|---|
| R1 | `Invoke` selector call | `*ast.CallExpr` → `*ast.SelectorExpr` with `Sel.Name == invokeSel` (**exact equality** — D4) |
| R2 | live constructor | `CallExpr` → `SelectorExpr` `<local>.NewSession` |
| R3 | replay constructor | `CallExpr` → `SelectorExpr` `<local>.NewReplaySession` |
| R4 | `Session` type exposure | any `SelectorExpr` `<local>.Session` (covers param, field, var, return, composite literal) |
| R5 | dot-import of `host/broker` | `ImportSpec` with `Name.Name == "."` |
| R17 | alias resolution | R2–R4 must fire through a non-default local name |

R1 is receiver-type-blind by design: it flags **any** `x.Invoke(...)`, not only
`Session.Invoke(...)`, because resolving receiver types needs `go/packages` type information this
walk deliberately does not load. That is a **fail-closed over-approximation** and it is recorded as a
**stated limitation** in the file's doc comment, in the same form as `host/boundary`'s two stated
limitations — a future unrelated `Invoke` method outside `host/broker` reds this gate and the
message must say so.

Every finding carries `{path, line, kind, enclosingFunc}`. `enclosingFunc` is resolved from the
file's `*ast.FuncDecl` spans; it is what makes T3's identity assertion possible.

**Exit gate**: `go vet ./host/...` rc=0; hold set green. (No assertion yet — T3 and T4 consume this.)

### T3 — the real-repository gate (~150 test LOC)

Two subtests of `TestRegistryDispatchBindingBoundary`:

**`outside_broker_is_clean`** — every walked file **not** under `host/broker/` yields zero findings.
On failure the message names `path:line kind=… fn=…` for every finding. Base: **0 findings** (V8).

**`inside_broker_exemption`** — files under `host/broker/` yield findings of kind `invoke-call` only,
and that finding set must equal, **as a set of identities and as an exact count**, the frozen
exemption:

| file | enclosing function | n |
|---|---|---:|
| `host/broker/publish_op.go` | `mintAttendedApproval` | 2 |
| `host/broker/publish_op.go` | `invokeAttendedPublish` | 1 |
| | **total** | **3** |

Measured at base by a standalone `go/parser` prototype run against this exact tree (V8):

```
INSIDE host/broker findings=3
  host/broker/publish_op.go:135 kind=invoke-call fn=mintAttendedApproval
  host/broker/publish_op.go:162 kind=invoke-call fn=mintAttendedApproval
  host/broker/publish_op.go:279 kind=invoke-call fn=invokeAttendedPublish
OUTSIDE host/broker findings=0
```

This reproduces the design doc's frozen exemption set **exactly**, including the two unexported
enclosing functions (`mintAttendedApproval` at `publish_op.go:127`, `invokeAttendedPublish` at
`:257` — V14). The doc's carve-out is confirmed against the live tree, not taken on trust.

**Two assertions, not one** (D-note for §0.2(iv)): the count assertion is `== 3` (exact, never `<=`
or `>=`), and the identity assertion is set-equality over `{file, enclosingFunc, kind}` with per-
identity multiplicity. Either alone is bypassable; M3 and M10 are the arms that prove it.

**Exit gate**: AC11 `count=2`, both tests PASS with the tolerant arm still in place. Hold set green.

### T4 — the hermetic synthetic control suite (~180 test LOC)

This is the task §0.2(ii) exists for. Twelve control sources, parsed from **in-memory strings**
through `parser.ParseFile(fset, name, src, parser.ParseComments)` and fed to the **same `scanFile`**
the repo walk uses. Subtest `detector_controls`, one `t.Run` per case, all measured by the prototype
(V11):

| case | source shape | required | measured |
|---|---|---|---|
| `POS-invoke-outside` | `s.Invoke(...)` outside broker | 1 × `invoke-call` | ✅ 1 |
| `POS-ctor-live` | `broker.NewSession(...)` | 1 × `ctor-live` | ✅ 1 |
| `POS-ctor-replay` | `broker.NewReplaySession(...)` | 1 × `ctor-replay` | ✅ 1 |
| `POS-session-type` | `func f(s *broker.Session)` | 1 × `session-type` | ✅ 1 |
| `POS-alias-session` | `import bk "…/broker"; *bk.Session` | 1 × `session-type` | ✅ 1 |
| `POS-dot-import` | `import . "…/broker"` | 1 × `dot-import` | ✅ 1 |
| `NEG-invokeattended` | `broker.InvokeAttendedPublish(...)` | **0** | ✅ 0 |
| `NEG-bound-request` | `b.Request(...)` on a `*broker.BoundInvoker` | **0** | ✅ 0 |
| `NEG-prose-only` | a comment naming `Session.Invoke` and `broker.NewSession` | **0** | ✅ 0 |
| `NEG-unrelated-selector` | `x.Ping()` with no broker import | **0** | ✅ 0 |
| `WALK-nested-module` | `t.TempDir()` tree with a nested `go.mod` | walker does not descend | §0.2(vi) |
| `WALK-unparseable` | `t.TempDir()` file with a syntax error | walker/scan returns an **error**, test FAILS loudly | R14 |

`NEG-prose-only` is the committed proof of the doc's "a text scanner cannot distinguish code from
prose". Its live counterpart is measured: across the **same 39 files**, a substring scanner sees
**26** occurrences of the token; the AST detector sees **3** (V9). That 26-vs-3 gap is M20's kill
signal.

`WALK-unparseable` is the doc-implicit branch nobody names: a walker that `continue`s past a parse
error is a walker that can be blinded by one unparseable file. It must **fail**, never skip.

**Exit gate**: all twelve controls PASS. AC11 `count=2`. Hold set green.

### T5 — the mutation sweep: 23 mutations, 46 arms (~0 LOC, the bulk of the time)

§4 in full. Transcript at `design_docs/verification/w-transition-registry/trc-mutations.md`.

### T6 — AC11 zero-tolerance activation and the doc repairs (the merge criterion)

The milestone is **not done** without this. TR.A and TR.B both carried the identical step.

- **T6.a** Edit AC11 in the design doc: delete `test "$count" -eq 1 || {` and its closing `}`,
  require `-eq 2`, and run both tests unconditionally.
- **T6.b** Machine check that no tolerant arm survives in AC11, with a known-positive control in the
  same call (grep the AC1–AC10 block, which legitimately has none left either — so the control must
  be a *deliberately constructed* string, not a sibling AC; see §5).
- **T6.c** Record **M23 `MUT-DELETE-TR-C-TEST`** RED: rename `TestRegistryDispatchBindingBoundary`,
  re-run the **activated** AC11, require rc≠0 with `count=1`; restore; re-run, require rc=0
  `count=2`.
- **T6.d** Apply the four divergence repairs D1–D4 and the two notes D5–D6 to the design doc's
  mutation table and Files table (§2).
- **T6.e** Final gates: `./scripts/verify_ail.sh` (4/11/14), `./scripts/verify_go.sh`,
  `go vet ./host/...`, `AC-INVOKE3`, and the complete hold set re-measured (§5).

---

## 4. Mutation discipline

### 4.1 Protocol — every arm, no exceptions

For each mutation:

1. **Back up** the target to `/tmp/trc_backup/` with `cp`. Record `sha256` **before**.
2. Apply the mutant. **Neuter with `if false && <cond>`, never by deleting a block** — deletion
   breaks the build, and *"the mutant does not compile"* wears the same exit code as *"the guard
   fired"*.
3. **Assert the mutant LANDED**: `sha256` after ≠ `sha256` before. Record both.
4. **Assert the mutant BUILDS**: `GOTOOLCHAIN=go1.25.6 go build ./...` → **rc=0**, and
   `GOTOOLCHAIN=go1.25.6 go vet ./host/...` → rc=0. **Do not read any test result before this
   passes.**
5. **Kill arm**, scoped with `-run` to the *named subtest* (rule 3i): require rc≠0 **and** record the
   **FAIL line and subtest name**, never the exit code alone. An rc=1 in exactly the direction you
   predicted may be the `host/broker` base flake (queue item 16, ~18%).
6. **Inverse arm**, same mutant: re-run with the new test excluded
   (`-skip 'TestRegistryDispatchBindingBoundary$'`, plus the flake skip) and require **rc=0**. This
   is what proves your test is the killer rather than a bystander. Where the inverse arm is
   unsatisfiable by construction (a shared mechanism co-detected by a landed test), **record that
   phrase**; do not weaken the co-detector.
7. **Restore** with `cp` from `/tmp/trc_backup/`. **Never `git checkout --`.** Assert `sha256` equals
   the recorded before-value.
8. **Re-run the kill arm** and require rc=0. A restore that did not restore is otherwise invisible.

**Every broker command carries `AILANG_BIN=/tmp/ailang-v0300/ailang`** (v0.30.0, commit `e37b370`,
V15). Without it `host/broker` is 100% red on `TestEpisodeLiveReplayThreeArmsAndEvidence` — a red
that measures the environment.

**Every whole-package broker run carries `-skip 'TestHandlerTimeoutKillsTheWholeProcessGroup$'`.**

### 4.2 Rule 3j — the branch inventory, anchored to TR.C's OWN diff

The unit of mutation is the **branch**, not the milestone. TR.C's entire deliverable is a refusal, so
the enumeration is the complete list of ways `TestRegistryDispatchBindingBoundary` can refuse.

**The rule-3j first cut is broken in two ways this repo has measured, and both bite TR.C:**

1. `grep -c 'return .*fmt.Errorf(.*%w'` is blind to every non-wrapping refusal. TR.C's refusals are
   `t.Fatalf`/`t.Errorf`, so even the repaired
   `grep -cE 'return .*(fmt\.Errorf|errors\.New)\('` returns **0**. The correct cut for a *test*
   deliverable is `grep -cE '\bt\.(Fatal|Fatalf|Errorf)\('`, and it is a **floor, not the
   enumeration**.
2. **Ordinary `git diff` OMITS UNTRACKED FILES.** `invoke_boundary_test.go` is **new**, so any
   rule-3j cut anchored to `git diff` returns **0** on this sprint. Use
   `git diff --no-index /dev/null host/broker/invoke_boundary_test.go`. Verify the instrument fires
   before banking its output: a 0 you cannot explain is not a 0.

**Cut instrument for T5, with its known-positive control in the same call:**

```bash
export PATH=/opt/homebrew/bin:$PATH
cd <worktree>
# the cut — untracked-safe
git diff --no-index /dev/null host/broker/invoke_boundary_test.go \
  | grep -cE '^\+.*\bt\.(Fatal|Fatalf|Errorf)\('
# known-positive control, SAME instrument, SAME kind of path, a file known to have many
git diff --no-index /dev/null host/boundary/allowlist_world_test.go \
  | grep -cE '^\+.*\bt\.(Fatal|Fatalf|Errorf)\('     # base: 46
# known-negative control, SAME instrument, scoped to the SAME directory
git diff --no-index /dev/null host/broker/broker.go \
  | grep -cE '^\+.*\bt\.(Fatal|Fatalf|Errorf)\('     # base: 0
```

The enumeration below is by **reading the file**, with the grep only as a floor.

### 4.3 The 23 mutations / 46 arms

**Group A — production mutants. Prove the gate catches the real threat.**

| # | ID | Mutation (compiling) | Branch | Required RED observation (scoped `-run`) |
|---|---|---|---|---|
| M1 | `MUT-BINDING-FOURTH-INVOKE-IN` | add a 4th `.Invoke(` **inside** `host/broker` (`confined.go`, reachable dead code so it compiles and vets) | R1/count | `…/inside_broker_exemption`: count=4 want 3 |
| M2 | `MUT-BINDING-FOURTH-INVOKE-OUT` | add an `.Invoke(` in `host/transitionreg/bind.go` | R1/outside | `…/outside_broker_is_clean` names `host/transitionreg/bind.go:N kind=invoke-call` |
| M3 | `MUT-BINDING-MOVE-SITE` | move the `publish_op.go:162` call out of `mintAttendedApproval` into a new unexported helper — **total stays 3** | R7 identity | `…/inside_broker_exemption`: identity set mismatch; the **count assertion stays green** |
| M4 | `MUT-BINDING-CTOR-LIVE` | `broker.NewSession(...)` in `host/transitionreg/bind.go` | R2 | outside-clean, `kind=ctor-live` |
| M5 | `MUT-BINDING-CTOR-REPLAY` | `broker.NewReplaySession(...)`, same file | R3 | outside-clean, `kind=ctor-replay` |
| M6 | `MUT-BINDING-SESSION-TYPE` | `func trcProbe(*broker.Session)` in `bind.go` | R4 | outside-clean, `kind=session-type` |
| M7 | `MUT-BINDING-ALIAS-IMPORT` | alias the broker import in `bind.go` to `bk`, use `*bk.Session` | R17 | outside-clean, `kind=session-type` **via the alias** |
| M8 | `MUT-BINDING-DOT-IMPORT` | dot-import `host/broker` in `cmd/world-publish/main.go` | R5 | outside-clean, `kind=dot-import` |

M7 is not redundant with M6: a detector that hardcodes the identifier `broker` instead of resolving
the import passes M6 and fails M7. That is the alias branch, and it has its own arm.

**Group B — test-file mutants. Prove the gate is not self-satisfying.**

| # | ID | Mutation | Branch | Required RED observation |
|---|---|---|---|---|
| M9 | `MUT-BINDING-RAISE-COUNT` | pin the exemption count 3 → 4 | R6 | count assertion RED (sees 3, wants 4) — proves the pin is exact |
| M10 | `MUT-BINDING-COUNT-INEQUALITY` | `!= 3` → `< 3`, **composed with M1** | R6+R7 | the **identity** assertion must still RED. If green, the identity assertion is decorative and T3 is wrong. Inverse: M10 **alone** must be rc=0 |
| M11 | `MUT-BINDING-IF-FALSE-INVOKE` | `if false &&` on R1 | R1 | `…/detector_controls/POS-invoke-outside`: findings=0 want 1 |
| M12 | `MUT-BINDING-IF-FALSE-CTOR-LIVE` | `if false &&` on R2 | R2 | `…/POS-ctor-live` RED |
| M13 | `MUT-BINDING-IF-FALSE-CTOR-REPLAY` | `if false &&` on R3 | R3 | `…/POS-ctor-replay` RED |
| M14 | `MUT-BINDING-IF-FALSE-SESSION-TYPE` | `if false &&` on R4 | R4/R17 | `…/POS-session-type` **and** `…/POS-alias-session` RED |
| M15 | `MUT-BINDING-IF-FALSE-DOT-IMPORT` | `if false &&` on R5 | R5 | `…/POS-dot-import` RED |
| M16 | `MUT-BINDING-EMPTY-WALK` | walker returns an empty slice | E1 | `…/enumeration`: "walked ZERO production .go files" — **must FAIL, never pass vacuously** |
| M17 | `MUT-BINDING-WALK-SKIP-BROKER` | walker skips `host/broker` | E2/E3 | anchor-missing RED naming `host/broker/publish_op.go`; count also 0≠3 |
| M18 | `MUT-BINDING-WALK-INCLUDE-TESTS` | drop the `_test.go` exclusion | E4 | `…/enumeration` names a `_test.go` path; also inside-count explodes far past 3 |
| M19 | `MUT-BINDING-GOLIST-DEAD` | feed the superset check an empty `go list` slice | E5 | `…/enumeration`: "go list enumerated 0 packages" — a vacuous superset check |
| M20 | `MUT-BINDING-TEXT-SCANNER` | replace the AST core with `strings.Contains(src, invokeSel)` | R1 + prose | inside count = **26** ≠ 3 (V9). This is the committed proof that the AST is load-bearing, and the arm that catches "a detector that finds itself" |
| M21 | `MUT-BINDING-NEG-CONTROL-DEAD` | make the detector emit a finding for every file | R16 | all four `NEG-*` controls RED (findings>0). Without this, M11–M15's positive controls could be satisfied by an always-firing detector |
| M22 | `MUT-BINDING-PARSE-SWALLOW` | `continue` on parse error instead of returning it | R14 | `…/detector_controls/WALK-unparseable` RED |
| M23 | `MUT-DELETE-TR-C-TEST` | rename `TestRegistryDispatchBindingBoundary` **after** T6.a | AC11 inventory | **activated** AC11 rc≠0 with `count=1`; the known control alone is no longer accepted |

**Totals: 23 mutations, 46 arms** (one kill + one inverse each; M10's inverse is the composed-vs-
alone pair).

### 4.4 Two arms flagged UNCERTAIN, to be resolved by the executor and reported

- **M1's inverse arm.** A fourth in-broker `.Invoke(` in `confined.go` may be co-detected by a landed
  broker test if the executor places it on a live path. **Place it on dead-but-referenced code** so
  the inverse arm (`-skip` the new test → rc=0) is satisfiable. If it proves unsatisfiable, record
  the co-detector by name and do not weaken it.
- **M20's build.** Swapping the AST core for a text scan changes the helper's signature. Keep the
  signature and change only the body, or the mutant fails step 4 and the arm proves nothing.

---

## 5. Acceptance commands, as the executor must run them

All baselined on the **pristine tree at `2d5a346`** this session; **the base result is part of the
criterion** (rule 3e). **No command contains `rg`.** Every command starts
`export PATH=/opt/homebrew/bin:$PATH` and carries `GOTOOLCHAIN=go1.25.6`.

### AC11 — the one criterion TR.C closes

```bash
# BASE FORM (as the doc has it today). Base observed: count=1, rc=0.
export PATH=/opt/homebrew/bin:$PATH
count=$(GOTOOLCHAIN=go1.25.6 go test ./host/broker -list 'Test(ReplayReturnsRecordedBytesWithoutDispatch|RegistryDispatchBindingBoundary)$' 2>/dev/null | grep -c '^Test' || true)
echo "AC11 count=$count"
test "$count" -eq 1 || { test "$count" -eq 2 && GOTOOLCHAIN=go1.25.6 go test ./host/broker -run 'Test(ReplayReturnsRecordedBytesWithoutDispatch|RegistryDispatchBindingBoundary)$' -count=1; }
echo "AC11 rc=$?"
```

```bash
# ACTIVATED FORM, written by T6.a. This is what must be green at merge.
export PATH=/opt/homebrew/bin:$PATH
count=$(GOTOOLCHAIN=go1.25.6 go test ./host/broker -list 'Test(ReplayReturnsRecordedBytesWithoutDispatch|RegistryDispatchBindingBoundary)$' 2>/dev/null | grep -c '^Test' || true)
test "$count" -eq 2 && GOTOOLCHAIN=go1.25.6 AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/broker -run 'Test(ReplayReturnsRecordedBytesWithoutDispatch|RegistryDispatchBindingBoundary)$' -count=1
```

### AC-ACTIVATED — T6.b's tolerant-arm check, with a control in the same call

```bash
# The AC11 block must no longer contain a base-tolerant arm.
export PATH=/opt/homebrew/bin:$PATH
D=design_docs/planned/w-transition-registry.md
tolerant=$(grep -c 'test "\$count" -eq 1 ||' "$D")          # required: 0 after T6.a
control=$(grep -c 'test "\$count" -eq 2 ' "$D")             # known-positive: >0, the instrument fires
echo "tolerant=$tolerant control=$control"
[ "$tolerant" -eq 0 ] && [ "$control" -gt 0 ]
```

`grep` exit 1 is "no match" and 2 is "no such file"; `grep -c` collapses both into a count, so the
control in the same call is what distinguishes "no tolerant arm" from "wrong path".

### The hold set — re-measured at every task exit and at T6.e

| criterion | command scope | **base observed at `2d5a346`** | required after TR.C |
|---|---|---|---|
| AC1 | `./host/transitionreg -list` (3 names) | `count=3` | **3** |
| AC2 | `./host/transitionreg -list` (3 names) | `count=3` | **3** |
| AC3 | `./host/transitionreg -list` (4 names) | `count=4` | **4** |
| AC4 | `./host/store -list` (2 names) | `count=2` | **2** |
| **AC5** | `./host/broker -list` (2 names) | **`count=2`** | **2** |
| **AC6** | `./host/transitionreg -list` (3 names) | **`count=3`** | **3** |
| **AC7** | `./host/transitionreg -list` (3 names) | **`count=3`** | **3** |
| AC8 | `./host/replay` | `count=2`, PASS | unchanged |
| AC9 | `./scripts/verify_ail.sh` | rc=0, `modules11=1 tests14=1 steps9=1` | **4/11/14 UNMOVED** |
| AC10 | `go build ./...` + focused | build rc=0 | unchanged |
| **AC11** | `./host/broker -list` (2 names) | **`count=1`** | **2, activated** |
| **`AC-INVOKE3`** | grep over `host/ cmd/` | **`n=3 p=3 t=90`, rc=0** | **n=3 p=3** |
| `AC-VET` | `go vet ./host/...` | **rc=0**, zero findings | rc=0 |
| `AC-NOFLAKE` | whole broker package, flake skipped | **rc=0**, `ok host/broker 35.752s` | rc=0, **< 60 s** |

```bash
# AC-INVOKE3 — the standing guard. If TR.C moves this, TR.C reds itself.
export PATH=/opt/homebrew/bin:$PATH
n=$(grep -rn '\.Invoke(' --include='*.go' host/ cmd/ | grep -v _test.go | wc -l | tr -d ' ')
p=$(grep -rn '\.Invoke(' --include='*.go' host/ cmd/ | grep -v _test.go | grep -c 'host/broker/publish_op.go')
t=$(grep -rn '\.Invoke(' --include='*.go' host/ cmd/ | grep -c _test.go)   # known-positive control
echo "n=$n p=$p t=$t"; [ "$n" -eq 3 ] && [ "$p" -eq 3 ] && [ "$t" -gt 0 ]
```

```bash
# AC-VET — copylocks and friends live OUTSIDE go test's default vet subset. Base: rc=0.
export PATH=/opt/homebrew/bin:$PATH
GOTOOLCHAIN=go1.25.6 go vet ./host/...
```

```bash
# AC-NOFLAKE — whole broker package, serial, base flake skipped, AILANG_BIN set.
export PATH=/opt/homebrew/bin:$PATH
GOTOOLCHAIN=go1.25.6 AILANG_BIN=/tmp/ailang-v0300/ailang \
  go test ./host/broker -skip 'TestHandlerTimeoutKillsTheWholeProcessGroup$' -count=1
```

```bash
# AC9 — the AILANG gate. TR.C touches no .ail; 4/11/14 must not move.
export PATH=/opt/homebrew/bin:$PATH
out=$(AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh 2>&1); rc=$?
m=$(printf '%s\n' "$out" | grep -c '4/4 required world/ identities verified across 11 module(s)')
t=$(printf '%s\n' "$out" | grep -c 'all 14 required named tests pass')
s=$(printf '%s\n' "$out" | grep -c 'world package gate PASSED: 9/9 steps')
echo "rc=$rc modules11=$m tests14=$t steps9=$s"
[ "$rc" -eq 0 ] && [ "$m" -eq 1 ] && [ "$t" -eq 1 ] && [ "$s" -eq 1 ]
```

---

## 6. Estimate — and why the doc's 0.5 day is 2.5× low

| Task | impl LOC | test LOC | notes |
|---|---:|---:|---|
| T1 walker + 5 anti-vacuity assertions | 0 | 150 | incl. the `go list` superset cross-check |
| T2 detector core + import resolution + assembled tokens | 0 | 200 | 6 branches, alias/dot handling, enclosing-func spans |
| T3 real-repo gate (outside-clean + exemption count **and** identities) | 0 | 150 | 2 subtests |
| T4 hermetic control suite | 0 | 180 | 12 cases, 2 of them `t.TempDir()` walker cases |
| T5 mutation sweep | 0 | 0 | **23 mutations / 46 arms** + transcript (~250 md lines) |
| T6 activation + doc repairs D1–D6 | 0 | 15 | design-doc lines, not Go; 1 delete arm, full gates |
| explanatory comments (this repo's landed style) | 0 | ~65 | `allowlist_world_test.go` is ~40% comment |
| **total** | **0** | **≈760** | 695 task LOC + 65 comment LOC |

**Reference velocity**: `VL.B` priced 515 LOC at 0.5 day → **~1000 LOC/day**. TR.A was priced at
2630 LOC over 2 days (~1315/day), its planner called that above velocity and recommended a split;
**the split held**. TR.B was priced at 1740 LOC in the doc's 1 day (~1740/day), was split, and
**that split held too**.

**TR.C's LOC alone is ~0.75 day.** The doc's 0.5 day is therefore already 1.5× low on code — but
that is not the real gap.

**The real gap is the sweep.** TR.B1 carried 16 mutations / 23 arms inside 0.75 day and TR.B2 carried
17 / 20 inside 1 day — call it **~20–25 arms/day** with production code to write alongside. TR.C
carries **46 arms**, and eight of them (M1–M8) mutate *production* files in `host/broker`,
`host/transitionreg`, and `cmd/world-publish`, each requiring a `go build ./...` + `go vet` gate
before any test result may be read, and a sha256-asserted `cp` restore afterwards. At TR.B's observed
rate that is ~0.5 day of sweep on its own.

### Honest price: **≈1.25 days. 2.5× the doc's 0.5 day.**

### 6.1 Verdict: **NO SPLIT.** And the number is still owed.

Both prior milestones of this item were split on the planner's own measurement. TR.C should not be,
and the reason is structural rather than a judgement call:

| candidate seam | why it fails |
|---|---|
| gate first, sweep + activation second | Lands a refusal-only deliverable whose non-vacuity is unproven. That is precisely the failure mode this item has hit three times. Rejected outright. |
| `Invoke` half (R1/R6/R7) first, constructor/type/import half (R2–R5/R17) second | Both halves live in **one test file** and **one top-level test**. The second milestone would edit the first's file, and the `MUT-DELETE-TR-C-TEST` activation arm can only fire once. Two milestones, one AC, one file — the split has no seam to cut along. |
| walker+detector first, real-repo gate second | The walker is not independently mergeable: alone it asserts nothing, so milestone 1 would ship a green that measures nothing. |

TR.C is one AC, one file, one test, and one activation. **Run it whole at ~1.25 days.** If the
controller needs it inside one day, the only honest lever is the arm count — and cutting arms is
what burned TR.A2, TR.B1, and TR.B2 in succession. **Do not cut the sweep.**

---

## 7. Execution protocol and risks

### Protocol

- **Worktree**: a sibling of the repo, e.g. `/Users/voightkampff/dev/sunholo-data/.wt-iter75`.
  **Never under `/tmp`** — cwd-relative path tests then fail for the location rather than the code.
  TR.C is *especially* sensitive: `repoRoot` is resolved from `runtime.Caller(0)`, so a relocated
  checkout is fine, but a `/tmp` checkout can collide with `os.TempDir()`-based assertions.
- **Every** Bash call starts `export PATH=/opt/homebrew/bin:$PATH` (else `go`/`gh`/`node` are rc=127).
- **Every** `go` invocation carries `GOTOOLCHAIN=go1.25.6`. Without it `scripts/verify_go.sh`
  FATALs with *"active toolchain go1.26.4 miscompiles host/store/scan.go"* — that is the script
  refusing an unpinned toolchain, a **base condition, not a regression**. Do not chase it.
- **Every whole-package `./host/broker` run carries `AILANG_BIN=/tmp/ailang-v0300/ailang`** and
  `-skip 'TestHandlerTimeoutKillsTheWholeProcessGroup$'`.
- **`go vet ./host/...` after every task and every mutation arm.** A green `go test` is not a green
  `go vet`.
- **zsh, not bash**: `${PIPESTATUS[0]}` expands EMPTY (zsh spells it `${pipestatus[1]}`, 1-indexed)
  — capture with `cmd > /tmp/out 2>&1; echo "rc=$?"`. Quote glob-shaped flag values
  (`--include='*.go'`), or zsh aborts with `no matches found` and you read 0 hits from a command that
  never ran. **Brace any variable followed by a colon** (`"${rev}:host/x"`) — this repo's Go all lives
  under `host/`, a history-modifier letter, so `git show "$rev:host/…"` silently reads `.ost/…` and
  returns a plausible zero. zsh does not word-split unquoted variables: use arrays and assert
  `${#ARR[@]}`.
- **`rg` is not a binary** — it is a harness-injected shell function, absent under `env -i` and in
  CI, used in 0 of `ci.yml` and all six `scripts/*.sh`. Never in a committed command or an AC.
- **`git diff` omits untracked files.** `invoke_boundary_test.go` is new. Use
  `git diff --no-index /dev/null <file>` (§4.2).
- **Restores are `cp` from `/tmp/trc_backup/`**, never `git checkout -- <file>`.
- **Never touch** `~/.ailang/state/mission-v1*` or the V1 checkout.
- **SANDBOX CAVEAT.** A gate verdict obtained inside a `workspace-write` sandbox is
  **UNINFORMATIVE — neither a pass nor a fail**: loopback binds are denied there, which both invents
  failures and hides real ones. Report sandbox results as "sandbox, uninformative"; the controller
  re-runs every gate outside the sandbox and that run is the verdict.

### Risks

| # | Risk | Assessment |
|---|---|---|
| **R1** | **The doc prices TR.C at 0.5 day; I price it at ≈1.25 days** (§6). | **DECISION WANTED**, but the recommendation is unambiguous: **run whole, no split** (§6.1). The only lever inside one day is the arm count, and cutting arms is this item's recurring failure. |
| **R2** | **`if false &&` on any outside-broker branch is invisible without the synthetic controls** (§0.2(ii)). If T4 is descoped, M11–M15 all become unfalsifiable and TR.C ships a possibly-disarmed detector wearing a green. | **T4 IS THE MILESTONE.** It is not optional polish. If time pressure arrives, cut nothing else first. |
| **R3** | R1 is receiver-type-blind: any future unrelated `x.Invoke(...)` outside `host/broker` reds this gate. | **Accepted and documented** as a stated limitation in the file header, fail-closed by design. Resolving receiver types needs `go/packages`, which is a dependency and a runtime cost this gate deliberately does not take. The RED message must say so, so a future sprint is not sent hunting a phantom. |
| **R4** | The `host/broker` base flake (~18%, `TestHandlerTimeoutKillsTheWholeProcessGroup`) makes 46 arms a coin-flip factory if read by exit code. | Mitigated: every arm is `-run`-scoped to a named subtest, and attribution is by **FAIL line**, never rc. Queue item 16; **not TR.C's to fix and TR.C must not silence it**. |
| **R5** | A mutant left behind by a failed restore would breach §1's "no earlier criterion moves". | Mitigated by §4.1 steps 7–8: sha256-asserted `cp` restore **and** a re-run of the kill arm requiring rc=0. `AC-INVOKE3` at T6.e is the independent backstop — a stray production `.Invoke(` moves `n` off 3. |
| **R6** | E5's `go list` subprocess adds a toolchain dependency to a unit test. | Precedent is landed and green in CI: `host/boundary/allowlist_world_test.go` calls `go list -deps` four times per run. Same helper shape, 90 s timeout, `cmd.Dir = root`. |
| **R7** | The two nested scratch modules under `design_docs/verification/` are uncontrolled space that a future sprint could fill with anything. | Mitigated by the nested-`go.mod` skip (§0.2(vi)), proven by a hermetic `t.TempDir()` control rather than a live mutation that cannot red. |

---

## 8. Verification Log

Every row is a command actually run on **2026-08-12** at **`2d5a346`**, clean tree, in the main
checkout (not a sandbox). Empty/negative results carry a **known-positive control scoped to the same
path in the same call**.

| ID | Codebase claim | Command | Observed output |
|---|---|---|---|
| V0 | inspected revision, clean tree | `git log --oneline -1; git status --porcelain` | `2d5a346 mission(world) iter 74: TR.B2 LANDED …`; **empty** status |
| V1 | **controller P1 confirmed**: production `.Invoke(` = exactly 3 | `grep -rn '\.Invoke(' --include='*.go' host/ cmd/ \| grep -v _test.go` + `\| wc -l` + `--include='*_test.go' … \| wc -l` | `publish_op.go:135,162,279` → **3**; known-positive `_test.go` control → **90** |
| V2 | **controller P2 confirmed + sharpened**: exported `broker.NewSession` has 0 production callers; `newSession` is unexported | `grep -rn 'NewSession(\|NewReplaySession(\|newSession(' --include='*.go' host/ cmd/ \| grep -v _test.go` | definitions at `broker.go:87,92,101`; the only *calls* are to the **unexported** `newSession` at `broker.go:88,98` and `publish_op.go:132,159,276`. Test control → **30** |
| V3 | **controller P3 confirmed**: `broker.Session`/`broker.NewSession` outside `host/broker` in production = 0 | `grep -rn 'broker\.Session\|broker\.NewSession' --include='*.go' host/ cmd/ \| grep -v _test.go \| grep -v '^host/broker/'` + control `grep -rn 'broker\.' … \| wc -l` | rc=1, **0 hits**; known-positive control in the same directory scope → **55** production `broker.` references, so the instrument fires |
| V4 | **controller P4 confirmed**: AC11 base | the exact AC11 command from the doc | `AC11 count=1`, `rc=0`; `-list` prints only `TestReplayReturnsRecordedBytesWithoutDispatch`; `test -e host/broker/invoke_boundary_test.go` → rc=1 (absent) |
| V5 | **controller P5 confirmed**: same-package positive control + hold set | AC5/AC6/AC7 exact commands | **AC5 count=2** (same package, same `-list` instrument as AC11 → the instrument fires), **AC6 count=3**, **AC7 count=3** |
| V6 | **controller P6 confirmed**: freshness sweep | `git diff --name-only b0f323a HEAD -- . ':(exclude)design_docs' \| wc -l`; `… \| grep '^host/broker/'`; control `git diff --name-only b0f323a HEAD \| wc -l` | **19** files; the 4 under `host/broker/` are `broker.go`, `broker_test.go`, `confined.go`, `decide.go` (TR.B's own landed work); control incl. docs **31** |
| V7 | **`host/boundary` placement is a MEASURED defect** — `TestBoundaryASTWriteGuard` pins `wantFileCount = 1` over **all** `.go` incl. `_test.go` | write `host/boundary/zz_trc_placement_probe_test.go` (`package boundary`), `go test ./host/boundary -run 'TestBoundaryASTWriteGuard$' -count=1`, remove, `git status --porcelain` | **rc=1**: `allowlist_world_test.go:1165: host/boundary contains 2 .go files, want 1: ["allowlist_world_test.go" "zz_trc_placement_probe_test.go"]`. After removal, `git status --porcelain` **empty** |
| V8 | **the detector design reproduces the doc's frozen exemption set exactly**, and the repo is clean outside `host/broker` | standalone `go/parser` prototype (stdlib only, `/tmp/trcprobe`) run against this tree | `walked=39 skipped_tests=44 skipped_nested_modules=2`; `INSIDE host/broker findings=3` — `publish_op.go:135 fn=mintAttendedApproval`, `:162 fn=mintAttendedApproval`, `:279 fn=invokeAttendedPublish`; **`OUTSIDE host/broker findings=0`** |
| V9 | **AST vs text discriminator, measured over the SAME 39 files** | same prototype, `strings.Count(src, "Invo"+"ke")` alongside the AST scan | text scanner → **26**; AST detector → **3**. M20's kill signal |
| V10 | **D4**: a prefix/substring reading of "any `Invoke` selector call" is RED AT BASE | `grep -rn 'Invoke' --include='*.go' host/ cmd/ \| grep -v _test.go` | `cmd/world-publish/main.go:367: broker.InvokeAttendedPublish(...)`, plus 22 comment/prose lines naming `Session.Invoke` |
| V11 | **all 12 hermetic controls behave as designed**, including the D4 known-negative | prototype `syn` mode over in-memory sources | `POS-invoke-outside 1 invoke-call`; `POS-ctor-live 1 ctor-live`; `POS-ctor-replay 1 ctor-replay`; `POS-session-type 1 session-type`; `POS-alias-session 1 session-type`; `POS-dot-import 1 dot-import`; **`NEG-invokeattended 0`**; `NEG-bound-request 0`; `NEG-prose-only 0`; `NEG-unrelated-selector 0` |
| V12 | **nested modules exist under `design_docs/`** | `find . -name go.mod -not -path './.git/*'` | `./go.mod`, `./design_docs/verification/w-race-gate-blindspot/racecontrol/go.mod`, `./design_docs/verification/w-race-gate-blindspot/repro/go.mod` |
| V13 | **`go list` would leave a real production file permanently unscanned** | `comm` of the walk set (nested modules excluded) against `go list -f '{{range .GoFiles}}…'` `./...` | walk **39**, `go list` **38**; `GOLIST-ONLY:` **none**; `WALK-ONLY: host/store/writer_lock_other.go`. `head -3` of that file → `//go:build !unix` |
| V14 | the two exempt enclosing functions are the doc's, and they are unexported | `grep -n '^func ' host/broker/publish_op.go` | `127:func mintAttendedApproval(...)` (contains :135, :162); `257:func invokeAttendedPublish(...)` (contains :279); exported wrappers `MintAttendedApproval` :123 and `InvokeAttendedPublish` :247 contain no `.Invoke(` |
| V15 | pinned binary identity | `/tmp/ailang-v0300/ailang --version` | `AILANG v0.30.0`, `Commit: e37b370` |
| V16 | `AC-INVOKE3` base | the §5 `AC-INVOKE3` block | `n=3 p=3 t=90`, **rc=0** |
| V17 | `AC-VET` base — non-vacuous and attributable | `GOTOOLCHAIN=go1.25.6 go vet ./host/...` | **rc=0**, zero findings |
| V18 | AC9 base — the AILANG totals TR.C must not move | the §5 AC9 block | `rc=0 modules11=1 tests14=1 steps9=1` |
| V19 | hold set AC1–AC4, AC10 base | the doc's exact AC1/AC2/AC3/AC4 `-list` commands; `go build ./...` | **3 / 3 / 4 / 2**; build **rc=0** |
| V20 | `AC-NOFLAKE` base and its cost budget | `go test ./host/broker -skip 'TestHandlerTimeoutKillsTheWholeProcessGroup$' -count=1` with `AILANG_BIN` | **rc=0**, `ok github.com/sunholo-data/ailang-world/host/broker 35.752s` |
| V21 | module identity and package inventory | `head -3 go.mod`; `go list ./...` | `module github.com/sunholo-data/ailang-world`, `go 1.25.6`; **17** packages |
| V22 | production vs test file split (the numbers E1/E2/E4 pin) | `find . -name '*.go' -not -name '*_test.go' -not -path './.git/*' \| wc -l`; without the `_test` exclusion | **41** raw / **85** total; **39** after the nested-module skip; **44** `_test.go` skipped |
| V23 | `host/broker` is `package broker` for internal tests, so the doc's placement compiles as specified | `head -1 host/broker/registry_publish_test.go host/broker/broker_test.go` | both `package broker` |
| V24 | no `vendor/`; two build-constrained files only | `test -d vendor`; `grep -rln '//go:build' --include='*.go' .` | vendor rc=1 (absent); `host/store/writer_lock_unix.go`, `host/store/writer_lock_other.go` |

No `.ail` source is proposed, so S5 pinned-binary source validation is not applicable; V18
nevertheless runs the pinned binary over every existing `.ail` module.

---

## 9. Handoff

- **Sprint plan**: `design_docs/planned/w-transition-registry-trc-sprint-plan.md` (this file)
- **Sprint JSON**: `.ailang/state/sprints/w-transition-registry-trc.plan.json`
- **Design doc**: `design_docs/planned/w-transition-registry.md`
- **Base**: `dev` @ `2d5a346`
- Neither artifact is committed by the planner. **The controller commits.**

SPRINT_PLAN_PATH: design_docs/planned/w-transition-registry-trc-sprint-plan.md
SPRINT_JSON_PATH: .ailang/state/sprints/w-transition-registry-trc.plan.json
