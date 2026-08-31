# w-canary-fence-passes-a-gutted-canary — the repo's version-agnostic positive-arm miscompile detector can be reduced to a no-op with every gate green, because the fence counts a TOKEN instead of an ASSERTION

**Status**: Planned
**Date**: 2026-08-30
**Queue item**: 49, `w-canary-fence-passes-a-gutted-canary` (clause-2, evaluator-found at P42,
controller-reproduced first-party before adoption)
**Estimated**: ≤0.2 day (one test function's known-positive control replaced with a structural
AST shape assertion, ~85–100 changed lines + four stdlib imports in `host/verifygate/toolchain_pin_gate_test.go`;
**no production code, no `host/store` assertion, no `run.sh`, no `ci.yml`, no `scripts/` edit**)
**Designer**: `pi:ollama/deepseek-v4-flash:0731-cloud` (design-doc-creator, iteration 138;
Fable rotation entry unavailable on weekly capacity)
**Revision**: round-2 narrow-refinement carve-out. Quorum rounds 1 and 2 both BLOCKED with all
three external reviewers PRESENT and `absent_reviewers` empty. Round 1's concrete fixes were
applied in the single designer revision. Round 2's two rejecting reviewers supplied the same
concrete top-level-assertion fix; it is applied verbatim as M13/V34 without a third quorum under
the ratified carve-out. No design-direction objection or human decision remains.
**Toolchain boundary**: every command below was run first-party in this worktree at `0bbb1a9`
(clean tree; porcelain 0 re-checked after every mutation arm), shell `zsh`,
`PATH=/opt/homebrew/bin:$PATH`, darwin/arm64, `go version` = `go1.26.6 darwin/arm64` (V2). The
full host gate was baselined with the pinned released binary `AILANG_BIN=/tmp/ailang-v0300/ailang`
(reports `AILANG v0.30.0`, V16). No AILANG (`.ail`) source is written or changed by this design;
the pinned `ailang` binary is exercised only in the gate baselines, never in the new test.

> **Thesis:** row 42's `TestCanaryDeclaresPositiveArmOnly` guards the canary's marker clause with
> `strings.Count(src, "stateRoot") >= 2` — a known-positive control that counts **occurrences of a
> string** rather than the presence of an **assertion**. Deleting `host/store/toolchain_canary_test.go`'s
> three real assertion lines (`if rows[0].field != "stateRoot" { t.Fatalf(...) }`, `:52-53`) and
> replacing them with a comment that still names `stateRoot` once drops the needle count **3 → 2**,
> which still satisfies `>= 2`. Measured first-party this session (V5): the fence `--- PASS`, the
> canary `TestToolchainCanary` `--- PASS` now asserting nothing, and the scoped suite fails only on
> the documented `AILANG_BIN`-unset condition — **the gutted canary is green everywhere**. The
> artifact being protected is the repo's **version-agnostic positive-arm detector** for the
> darwin/arm64 array-literal miscompilation (`verify_go.sh:214-216` calls it "the version-agnostic
> detector for any version that miscompiles the shape"; the known-bad runtime arm lives in the
> nested `design_docs/verification/w-race-gate-blindspot/repro` module + `run.sh`, and the CI wiring
> is pinned by `TestMiscompileInstrumentProbesPinnedToolchain`/`TestMiscompileInstrumentStepIsGatedInCI` —
> see the census in V31). **The fix:** replace the token-presence count with a **structural Go AST
> shape assertion** — parse the canary file and require that `TestToolchainCanary` exists exactly
> once, contains exactly one **top-level** `if` whose condition is the binary `!=` comparison
> `rows[0].field != "stateRoot"` (left operand structurally `rows[0].field`, operator `!=`, right
> operand the string literal `"stateRoot"`), and that the **outer** if-body contains exactly one
> **direct** `t.Fatalf` expression statement. A comment that merely names `stateRoot` is not an
> assertion and cannot satisfy the shape; neither can a constant-vs-constant no-op
> (`if "stateRoot" != "stateRoot"`), a wrong left operand (`if rows[0].n != "stateRoot"`), or a
> Fatalf hidden under `if false` inside the body. **The generalisable half (the row's own):** *a
> known-positive control that counts a TOKEN proves the file still mentions the subject, never that
> it still tests it* — the same distance between prose and code the mission recorded at iter-124
> (`Authorization` = 1, in a comment).

## Quorum round-1 block — the two blocking objections and how this revision applies them

Round-0 was BLOCKED. All three external reviewers were PRESENT; `absent_reviewers` is empty. The
two blocking objections are applied verbatim in substance; neither is argued.

**Objection 1 (gpt5-6-sol) — the loose shape accepts a green no-op canary.** The round-0 AST fence
accepted `if "stateRoot" != "stateRoot" { t.Fatalf(...) }` — it satisfies NEQ + right literal +
descendant Fatalf but is always false. Recursive `ast.Inspect` also accepted a Fatalf hidden under
`if false`. **Required fix (applied):** bind the exact canonical assertion skeleton. The left
operand must structurally be `rows[0].field`; operator `!=`; right operand the string literal
`"stateRoot"`; the **outer** if-body must contain a **direct** `t.Fatalf` expression statement, not
merely a descendant call. OD-1/default-NO is removed (the left operand is now required). RED
mutations are added for constant-vs-constant (M10), `other != "stateRoot"` (M11), and Fatalf nested
under `if false` (M12). Explicit exact-once anti-vacuity floors are retained. Measured: the old
loose shape returns `SHAPE OK` on M10 and M12 (V33); the revised shape REDs all three (V27–V29).

**Objection 2 (oc-glm-5-2) — the "only first-party detector" claim was unsupported and over-broad.**
The doc's bare claim was measured by the controller rather than forwarded as an unverified premise.
Controller census at origin/dev `0bbb1a9` found additional runtime reproducer/instrument code in
`design_docs/verification/.../run.sh` and `TestMiscompileInstrumentProbesPinnedToolchain`/CI wiring.
The precise measured source is `scripts/verify_go.sh:214-216` (recorded verbatim in V30). The
negative exclusion in the round-0 census was too broad because the nested runtime reproducer
deliberately lives under `design_docs`. **Required fix (applied):** the thesis no longer claims
"only first-party detector." It is narrowed to the exact supported statement: this canary is the
repo's **version-agnostic positive-arm detector**, while the nested module/`run.sh` is the
**known-bad runtime arm** and the CI wiring is pinned separately. The census is added with explicit
scope and limitation (V31).

## The finding in one paragraph

`TestCanaryDeclaresPositiveArmOnly` (`host/verifygate/toolchain_pin_gate_test.go:367`) reads the
canary file by path and runs three checks: a `stateRoot` count `>= 2` (`:374`, the known-positive
control), a `GOTOOLCHAIN` zero-needle (`:377`), and a `POSITIVE ARM ONLY` marker (`:380`). The
first check is the defect: it counts **tokens**, so a canary whose real assertion has been deleted
and replaced by a comment naming `stateRoot` once still passes (count 3 → 2, still `>= 2`). The
repro is exact and first-party (V5): with the three real lines at `host/store/toolchain_canary_test.go:52-53`
replaced by `// stateRoot is the expected field value.`, the fence `--- PASS`, the canary `--- PASS`
(now asserting nothing), and the scoped `go test ./host/verifygate/ ./host/store/ -count=1` sweep
fails **only** on the documented `AILANG_BIN`-unset condition (the module-manifest shim tests), not
on the gutting. The gutted mutant **builds and typechecks** (`go vet` rc=0, `go build` rc=0, V6) and
passes **all three** of the old checks (V7: `GOTOOLCHAIN`=0, marker=1, `stateRoot`=2) — so the only
thing that can catch it is a check that inspects the **assertion shape**, not the token count. The
repair is one structural AST assertion in the same test, replacing the `:374` count clause. This
does **not** falsify row 42's claim — row 42's Decision scopes Test B to *"the canary's assertion
has moved out of this file"* plus the known-bad-arm fence, and its residual section explicitly says
*"Test B fences one token"* (`w-canary-control-does-not-survive-a-floor-raise.md:261`). This item
closes that residual.

## Premises

Each premise is one or more Verification Log rows; a claim without a row does not appear here.

- **P1 — the fence is green on the pristine tree**: `TestCanaryDeclaresPositiveArmOnly` → `--- PASS`
  (V3); the canary `TestToolchainCanary` → `--- PASS` (V4).
- **P2 — the defect reproduces exactly as filed**: gutting the three assertion lines at
  `toolchain_canary_test.go:52-53` into a comment naming `stateRoot` once drops the count 3 → 2;
  the fence `--- PASS`, the canary `--- PASS`, and the scoped suite fails only on the documented
  `AILANG_BIN`-unset condition (V5). Restore was sha256-byte-identical (`a23cfa79…`), porcelain 0.
- **P3 — the gutted mutant is a real, buildable program**: `go vet ./host/store/` rc=0 and
  `go build ./host/store/` rc=0 with the gutted file in place (V6) — the mutant is not a syntax
  error; it is a valid Go test that asserts nothing.
- **P4 — the gutted mutant passes every OLD check**: `GOTOOLCHAIN` count 0 (zero-needle green),
  `POSITIVE ARM ONLY` present (marker green), `stateRoot` count 2 (old `>= 2` green) (V7). The RED
  can therefore only come from a check that inspects the assertion **shape** — this is what makes
  each mutation arm provably downstream of the new assertion rather than suite-wide.
- **P5 — the AST shape check is small, stdlib, and discriminates**: `go/parser`, `go/ast`,
  `go/token`, `fmt` are all stdlib (V8). A ~50-line helper parses the canary and requires
  `TestToolchainCanary` exactly once, one `if rows[0].field != "stateRoot"` assertion if-stmt, and
  one **direct** `t.Fatalf` expression statement in that if-body. Prototype verdicts: pristine
  `SHAPE OK` (V9); gutted RED (V10); function-missing RED (V11); Fatalf-missing RED (V12);
  ambiguous (two assertions) RED (V13); prose/comment-only RED (V14); inverted `==` RED (V15);
  Fatalf-moved-outside-if RED (V16); constant-vs-constant no-op RED (V27); wrong-left-operand RED
  (V28); Fatalf-nested-under-`if false` RED (V29).
- **P6 — the full host gate is green on pristine origin/dev with the pinned binary**:
  `AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.26.6 go test ./... -count=1` → rc=0 (V17);
  `./scripts/verify_ail.sh` → rc=0 (V18); scoped `go test ./host/verifygate/ ./host/store/ -count=1`
  with `AILANG_BIN` → rc=0 (V19); `go build ./...` rc=0 (V20); `gofmt -l host/verifygate/ host/store/`
  → 0 files (V21); `go vet ./host/verifygate/ ./host/store/` rc=0 (V22).
- **P7 — the systemic surface is enumerated with a live control**: the only gate that reads
  `host/store/toolchain_canary_test.go` is `TestCanaryDeclaresPositiveArmOnly`
  (`toolchain_pin_gate_test.go:368`); `verify_go.sh` and `ci.yml` do not read the canary file (V23,
  control `verify_go.sh`/`run.sh` referenced in Go → grep live). The canary file's `stateRoot`
  occurrences are exactly 3, at `:32` (the array literal), `:52` (the comparison), `:53` (the
  Fatalf message) (V24). Similar token-count positive controls under `host/verifygate` number 23
  `strings.Count` call-sites across 4 test files (V25); the `:374` clause is the only one that
  guards an **assertion** rather than a pin/needle/marker, and the only one whose subject is a
  different file's test body. The controller census (V31) scopes the detector claim: the canary is
  the version-agnostic positive-arm detector; the nested module/`run.sh` is the known-bad runtime
  arm; the CI wiring is pinned separately.

### Design Freeze

- **One test function modified, no test file created**: the known-positive control inside
  `TestCanaryDeclaresPositiveArmOnly` (`toolchain_pin_gate_test.go:374-376`) is replaced by a call
  to a new package-level helper `canaryAssertionShapeProblems(src string) []string`; the helper and
  its three tiny predicates (`isRowsField`, `isStateRootLit`, `isTFatalfCall`) are added to the
  same file. Four stdlib imports added: `fmt`, `go/ast`, `go/parser`, `go/token`. No new module
  dependency.
- **The `GOTOOLCHAIN` zero-needle (`:377`) and the `POSITIVE ARM ONLY` marker (`:380`) are
  retained byte-unchanged** — they are separate fences for separate purposes (a re-added in-module
  bad arm; the required label). Only the `:374` token-count clause is replaced.
- **No production code, no `host/store` assertion, no `run.sh`, no `ci.yml`, no `scripts/` edit**:
  the canary file is read-only (by path) from `host/verifygate`, exactly as today.
- **The shape check is structural, not textual**: it parses the canary with `go/parser` and walks
  the AST. A comment, a string literal in an array, or prose that merely names `stateRoot` is not
  an assertion and cannot satisfy the shape (P5, V10/V14).
- **The shape is bound to the exact canonical assertion skeleton** (quorum objection 1): the left
  operand must structurally be `rows[0].field` (a `SelectorExpr` on an `IndexExpr` with index `0`),
  the operator `!=`, the right operand the string literal `"stateRoot"`, and the **outer** if-body
  must contain a **direct** `t.Fatalf` expression statement — not a descendant call. This rejects
  the constant-vs-constant no-op, a wrong left operand, and a Fatalf hidden under `if false`
  (V27–V29).
- **Anti-vacuity is built in**: the target function and each assertion component must be found
  **exactly once**; a count of 0 (missing) or >1 (ambiguous) is a loud RED with a named message
  (P5, V11/V12/V13). The parser itself is a floor: a parse error is a RED, never a silent pass.
- **No control derives its expectation from the value it checks**: the shape is authored against
  the known-good assertion text, not generated from the canary's own contents.

## Decision — a structural Go AST shape assertion replaces the token-presence count

The queue row offers two candidate mechanisms: (a) a shape assertion — "require the comparison
against `"stateRoot"` and a `t.Fatalf` in the same function" — and (b) a behavioural fence — "have
`host/verifygate` execute the canary and require it to red on a seeded-wrong value." The task brief
directs: *"Prefer a structural Go AST check over fragile text matching if it stays small."* **This
design chooses (a) with a Go AST check**, and rejects (b) — see Alternatives.

### The new shape assertion

Replace the `:374-376` clause in `TestCanaryDeclaresPositiveArmOnly` with:

```go
if problems := canaryAssertionShapeProblems(src); len(problems) > 0 {
    for _, p := range problems {
        t.Errorf("%s: %s", canaryPath, p)
    }
    t.Fatalf("instrument failure: %s no longer asserts the miscompile shape", canaryPath)
}
```

and add the helper (same file, package `verifygate`):

```go
// canaryAssertionShapeProblems parses the canary test source and reports any deviation from the
// required assertion SHAPE: the function TestToolchainCanary must exist exactly once, contain
// exactly one top-level `if` whose condition is the binary `!=` comparison rows[0].field != "stateRoot"
// (left operand structurally rows[0].field, operator !=, right operand the string literal
// "stateRoot"), and that if-body must contain exactly one DIRECT t.Fatalf expression statement
// (not merely a descendant call). Empty result = shape holds. This is a structural AST check, not
// a token count: a comment that merely names "stateRoot" is not an assertion, and neither is a
// constant-vs-constant no-op, a wrong left operand, or a Fatalf hidden under `if false` (queue row
// 49; quorum objection 1). Each component is required exactly once so an ambiguous or missing shape
// is loud, never silently accepted.
func canaryAssertionShapeProblems(src string) []string {
    var problems []string
    fset := token.NewFileSet()
    f, err := parser.ParseFile(fset, "", src, 0)
    if err != nil {
        return []string{fmt.Sprintf("parse error: %v", err)}
    }
    var funcs []*ast.FuncDecl
    for _, d := range f.Decls {
        if fd, ok := d.(*ast.FuncDecl); ok && fd.Name.Name == "TestToolchainCanary" {
            funcs = append(funcs, fd)
        }
    }
    if len(funcs) != 1 {
        return []string{fmt.Sprintf("TestToolchainCanary func decl count=%d, want exactly 1", len(funcs))}
    }
    var assertions []*ast.IfStmt
    for _, st := range funcs[0].Body.List {
        is, ok := st.(*ast.IfStmt)
        if !ok {
            continue
        }
        if be, ok := is.Cond.(*ast.BinaryExpr); ok && be.Op == token.NEQ &&
            isRowsField(be.X) && isStateRootLit(be.Y) {
            assertions = append(assertions, is)
        }
    }
    if len(assertions) != 1 {
        return []string{fmt.Sprintf("top-level `rows[0].field != \"stateRoot\"` assertion if-stmt count=%d, want exactly 1", len(assertions))}
    }
    body := assertions[0].Body
    var directFatalf int
    for _, st := range body.List {
        es, ok := st.(*ast.ExprStmt)
        if !ok {
            continue
        }
        if c, ok := es.X.(*ast.CallExpr); ok && isTFatalfCall(c) {
            directFatalf++
        }
    }
    if directFatalf != 1 {
        return []string{fmt.Sprintf("direct t.Fatalf expression statement in assertion body count=%d, want exactly 1", directFatalf)}
    }
    return problems
}

// isRowsField reports whether e is structurally rows[0].field: a SelectorExpr whose X is an
// IndexExpr on the identifier `rows` with the integer literal index 0, and whose Sel is `field`.
func isRowsField(e ast.Expr) bool {
    sel, ok := e.(*ast.SelectorExpr)
    if !ok {
        return false
    }
    if sel.Sel.Name != "field" {
        return false
    }
    ix, ok := sel.X.(*ast.IndexExpr)
    if !ok {
        return false
    }
    id, ok := ix.X.(*ast.Ident)
    if !ok || id.Name != "rows" {
        return false
    }
    bl, ok := ix.Index.(*ast.BasicLit)
    return ok && bl.Kind == token.INT && bl.Value == "0"
}

func isStateRootLit(e ast.Expr) bool {
    bl, ok := e.(*ast.BasicLit)
    return ok && bl.Kind == token.STRING && bl.Value == `"stateRoot"`
}

func isTFatalfCall(c *ast.CallExpr) bool {
    sel, ok := c.Fun.(*ast.SelectorExpr)
    if !ok {
        return false
    }
    id, ok := sel.X.(*ast.Ident)
    if !ok {
        return false
    }
    return id.Name == "t" && sel.Sel.Name == "Fatalf"
}
```

**Why the shape is exactly this.** The canary's real assertion is `if rows[0].field != "stateRoot" {
t.Fatalf(...) }` (`:52-53`). The four load-bearing components are: (1) the function that carries it,
(2) the `rows[0].field != "stateRoot"` comparison (the miscompile detector's check — a `==` would
invert the assertion and fire on the wrong condition; a constant-vs-constant `"stateRoot" !=
"stateRoot"` is always false and asserts nothing; a wrong left operand such as `rows[0].n` checks
the wrong field), and (3) the **direct** `t.Fatalf` in the same if-body (the failure action — a
Fatalf elsewhere in the function, e.g. the `len(rows)` or `len(Field)` guards, is not the
assertion's action; a Fatalf hidden under `if false` inside the body is a descendant call that
never fires). Requiring each **exactly once** makes the check reject both a missing shape (count 0)
and an ambiguous shape (count >1) loudly. The `!=` operator is required (not any comparison) so an
inverted assertion is caught (V15). The left operand is bound to `rows[0].field` (not any `!=
"stateRoot"`) so a constant-vs-constant no-op and a wrong-field comparison are caught (V27, V28) —
this is the quorum objection-1 fix that replaces the round-0 OD-1 default-NO. The Fatalf is
required as a **direct expression statement in the outer if-body** (not merely a descendant in the
function) so a Fatalf moved out of the guard or hidden under `if false` is caught (V16, V29) — the
function legitimately contains other `t.Fatalf` calls (`:45`, `:50`), so a function-wide Fatalf
count would be vacuous.

**Why AST and not text matching.** A text match such as `strings.Contains(src, "rows[0].field !=
\"stateRoot\"")` is fragile to whitespace/formatting and can be satisfied by a comment or a string
literal elsewhere; it is the same token-presence class as the defect. `go/parser`/`go/ast` are
stdlib (P5, V8), the helper is ~50 LOC, and it inspects the actual statement tree — a comment is
not a statement, so prose cannot satisfy it (V14). This is the "small structural check" the brief
prefers.

### What the gate CANNOT catch — declared residual

- **The shape is not the behaviour.** The AST check proves the canary *contains* the assertion
  shape; it cannot prove the assertion *fires* on a miscompiling toolchain. The runtime proof lane
  is the nested repro module + `run.sh`'s `saw_bad` floor (row 42's Test A and the reproducer),
  which is darwin-only and attended/local while `ci.yml:172` `continue-on-error: true` stands
  (row 44). This item does not change that lane.
- **A semantically-equivalent but differently-shaped assertion evades the shape.** If a future
  editor rewrites the assertion to a different AST shape (e.g. `if got := rows[0].field; got !=
  "stateRoot" { t.Fatalf(...) }` — an `IfStmt` with an `Init`), the `Cond` is still a `rows[0].field
  != "stateRoot"` binary expr and the Fatalf is still a direct statement in the body, so it passes —
  which is correct, because the assertion is still real. But a rewrite that changes the comparison
  to a helper call (`if !isStateRoot(rows[0].field) { t.Fatalf(...) }`) would RED the shape even
  though the assertion is real — a false positive that forces a deliberate, reviewed shape change.
  Likewise a legitimate rename of the local `rows` or the `field` selector would RED the shape
  (the left operand is now bound, per objection 1). Declared, not patched: the shape is pinned to
  the current canonical form on purpose, and a future refactor that wants a different shape must
  update the fence in the same commit (the same contract row 42's M7 placement arm enforces).
- **The `GOTOOLCHAIN` zero-needle and `POSITIVE ARM ONLY` marker remain token fences** over the
  same file; they are retained unchanged (Design Freeze) and carry row 42's own residual (a
  re-added bad arm avoiding the literal `GOTOOLCHAIN` string slips the zero-needle). This item
  does not claim to close that residual.
- **Prose is unguarded.** The canary's doc comment is prose; the `POSITIVE ARM ONLY` marker is the
  one prose artifact given a needle (retained). No new prose is bound.

## Alternatives rejected

1. **Behavioural fence — execute the canary and require it to red on a seeded-wrong value** (the
   queue row's "better" option). Rejected as the primary mechanism: a behavioural test of the
   assertion *logic* requires the logic to be extracted into a testable function or duplicated in
   `host/verifygate`, and it would test a **copy** of the assertion, not the canary file's actual
   assertion — so it does not directly protect the canary's shape. A behavioural test that runs the
   *miscompile* is darwin/arm64-only (the miscompile does not reproduce on linux/amd64 — row 44's
   measured fact), so on CI it is green whether or not the assertion exists — the exact vacuity
   class this item exists to remove. The AST check is static, lane-independent, and inspects the
   real file. Rejected as primary; noted as a possible future complement.
2. **Text matching on the assertion line** (`strings.Contains(src, "rows[0].field != \"stateRoot\"")`).
   Same token-presence class as the defect; fragile to formatting; satisfiable by a comment or a
   string literal elsewhere. Rejected.
3. **A regex over the source** for the assertion pattern. Fragile to whitespace and to the same
   comment/string-literal evasion; a regex is text matching with extra steps. Rejected.
4. **Raise the `stateRoot` count floor** (e.g. `>= 3`). The gutted mutant drops to 2, so `>= 3`
   would catch *this* mutant — but it is a token count, so a future gutting that retains three
   `stateRoot` tokens (e.g. a comment naming it three times) slips it, and it would also RED a
   legitimate refactor that legitimately reduces the count. It does not assert shape. Rejected.
5. **Move the canary's assertion into a shared helper both files call.** Touches `host/store`
   production/test code (out of scope per the brief) and changes the canary's shape, which is the
   very thing the fence protects. Rejected.

## Ordering

Gated on nothing. Neighbours named and not absorbed: **row 42** (`w-canary-control-does-not-survive-a-floor-raise.md`)
is the parent — this item closes its declared residual *"Test B fences one token"* (`:261`) and
reuses its Test B home and read-by-path pattern; **row 44** owns `run.sh`/`ci.yml` inertness — this
item changes no runtime behaviour of the instrument; **row 45** (`GOTOOLCHAIN` normalizer) and
**row 48** (`racecontrol` floor bump) are untouched. The sprint applies the single test-function
edit; no other file changes.

## Files to Create/Modify

- **MODIFY** `host/verifygate/toolchain_pin_gate_test.go` (+~50 LOC: the `:374-376` clause replaced
  by the shape-assertion call; the `canaryAssertionShapeProblems` helper + `isRowsField` +
  `isStateRootLit` + `isTFatalfCall`; imports `fmt`, `go/ast`, `go/parser`, `go/token` added). No
  name collisions at base (V26: no existing `canaryAssertionShapeProblems`/`isRowsField`/
  `isStateRootLit`/`isTFatalfCall` in the package).
- **No other files.** `host/store/toolchain_canary_test.go` (read-only, by path), `run.sh`,
  `ci.yml`, root `go.mod`, `scripts/*`, `racecontrol/` — untouched.

## Conflict Surface

This is a Go test/gate change, not an AILANG parser/typechecker/codegen change, so the Conflict
Surface is not mechanically mandatory; it is included because the gate is shared machinery.

- **`TestCanaryDeclaresPositiveArmOnly` itself** — the only test edited. Its `GOTOOLCHAIN` zero-needle
  (`:377`) and `POSITIVE ARM ONLY` marker (`:380`) are retained byte-unchanged; only the `:374`
  clause is replaced. The sprint re-confirms both retained clauses still pass on the post-sprint
  tree (AC2).
- **`TestToolchainCanary` (`host/store/toolchain_canary_test.go`)** — the subject, read-only by
  path. Its assertions are untouched, so its green under the pinned toolchain (V4) is unaffected;
  AC5's vet/gofmt covers the edit.
- **`TestMiscompileInstrumentProbesPinnedToolchain` and `TestReproModuleFloorStaysBelowKnownBadToolchains`**
  — siblings in the same file, untouched; they read `run.sh`/`repro/go.mod`, not the canary file.
  The new helper is package-private and name-collision-free (V26).
- **`go/parser`/`go/ast`/`go/token`/`fmt`** — stdlib, no new dependency; no conflict with the
  existing `go/version` import.
- **`verify_go.sh` / `ci.yml`** — do not read the canary file (V23); no interaction. The canary is
  the version-agnostic positive-arm detector; `verify_go.sh:214-216` and the CI wiring are separate
  arms (V30, V31).

## Systemic-Issue Audit

Is "a known-positive control that counts a token instead of an assertion" a pattern? Census of
token-count positive controls under `host/verifygate`: **23** `strings.Count` call-sites across 4
test files (V25). Each was read and classified: the `ail_binary_gate_test.go` needles (`:415`,
`:471`, `:480`, `:514`, `:692`, `:721`) and `floor_raise_inventory_test.go`/`module_manifest_gate_test.go`
needles count **pins, markers, and manifest identities** — artifacts whose *presence* is the
property being guarded, so a token count is the honest bound there. The `toolchain_pin_gate_test.go`
needles (`:150`, `:170`, `:174`, `:307`, `:407`) likewise count pin/needle/marker occurrences. The
**one** control that guards an **assertion** — a property that is only real when the code *does*
something, not when it *mentions* something — is the `:374` `stateRoot` count over the canary file.
That is the single instance of the defect class, and it is the one this item fixes. The
generalisable rule the row records — *a token count proves the file still mentions the subject,
never that it still tests it* — is the same distance between prose and code the mission recorded at
iter-124 (`Authorization` = 1, in a comment). The fix is correctly local: the other 22 controls
guard presence, not behaviour, and are not this defect.

**Controller census with explicit scope and limitation (quorum objection 2).** The round-0 claim
that this canary is the mission's "only first-party detector" was over-broad and is withdrawn. The
controller measured the detector surface at origin/dev `0bbb1a9` with:

```
rg -n -i 'miscompil|array[- ]literal|toolchain canary|ToolchainCanary' --glob '*.go' --glob '*.sh' --glob '*.yml' --glob '!design_docs/**' .
```

which found `scripts/verify_go.sh`, `host/store/toolchain_canary_test.go`, and the verifygate
tests/wiring (`TestMiscompileInstrumentProbesPinnedToolchain`, `TestMiscompileInstrumentStepIsGatedInCI`,
`miscompileReproducerPath`). The negative exclusion `!design_docs/**` was too broad because the
nested runtime reproducer deliberately lives under `design_docs/verification/w-race-gate-blindspot/`
(`repro` module + `run.sh`). **Narrowed thesis (V31):** this canary is the repo's
**version-agnostic positive-arm detector** for the array-literal miscompile; the nested module +
`run.sh` is the **known-bad runtime arm**; the CI wiring is pinned separately. The precise
two-line citation for the version-agnostic claim is `scripts/verify_go.sh:214-216` (V30). This item
protects only the positive-arm canary's assertion shape; it does not claim to be the only detector.

## Deferred Scope

- **Row 44 `w-miscompile-instrument-inert-in-ci`** — `ci.yml:172` `continue-on-error: true` discards
  `run.sh`'s exit code; the runtime proof lane is darwin-only. Named, not absorbed; this item's
  static shape check runs in any lane.
- **Row 42's `GOTOOLCHAIN` zero-needle residual** — a re-added bad arm avoiding the literal
  `GOTOOLCHAIN` string slips the zero-needle; retained unchanged, not re-scoped here.
- **A behavioural complement** (execute the canary's assertion logic against a seeded-wrong value)
  — rejected as primary (Alternatives 1); if taken later it is a separate item, not folded in.

## Acceptance Criteria

Each AC carries its vacuity self-test and its **observed result on the unmodified tree at `0bbb1a9`**
(Verification Log rows cited). The sprint re-runs each on the post-sprint tree.

- **AC1 — the shape assertion exists, RUNS, and passes on the post-sprint tree, in run-existence
  form.** `GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run '^TestCanaryDeclaresPositiveArmOnly$'
  -count=1 -v` → rc=0 with exactly 1 top-level `=== RUN` and 1 `--- PASS`; a paired nonsense pattern
  (`-run TestNoSuchCanaryFenceZZZ`) prints `[no tests to run]`, proving the instrument says so rather
  than passing vacuously. **Base @0bbb1a9: the verbatim command → `--- PASS` (V3); the nonsense form
  is the vacuity self-test.**
- **AC2 — the retained fences stay green.** On the post-sprint tree, the `GOTOOLCHAIN` zero-needle
  and `POSITIVE ARM ONLY` marker still pass (they are byte-unchanged); the scoped
  `GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run '^TestCanaryDeclaresPositiveArmOnly$' -count=1`
  → rc=0. **Base: both retained clauses green (V3).**
- **AC3 — the recorded gutted mutant REDs the new fence, and builds/typechecks.** Sprint evidence:
  apply the exact gutting (delete `:52-53`, replace with `// stateRoot is the expected field value.`),
  run `go vet ./host/store/` (rc=0, the mutant is a valid program) and the scoped fence
  (`go test ./host/verifygate/ -run '^TestCanaryDeclaresPositiveArmOnly$' -count=1 -v`) → **rc=1,
  `--- FAIL`**, message names `rows[0].field != "stateRoot"` assertion if-stmt count=0. Restore
  sha256-byte-identical (`a23cfa79…`), porcelain 0. **Base: the mutant builds (V6) and passes all
  OLD checks (V7); the prototype REDs it (V10).**
- **AC4 — a prose/comment token does not satisfy the fence.** The gutted mutant's comment names
  `stateRoot` once (count 2) yet the shape check REDs it (V10/V14) — a comment is not an assertion.
  **Base: prototype-proven (V14).**
- **AC5 — hygiene**: `go vet ./host/verifygate/ ./host/store/` rc=0 and `gofmt -l host/verifygate/
  host/store/` prints nothing. **Base: both green (V21, V22).**
- **AC6 — the full host gate stays green on the post-sprint tree.**
  `AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.26.6 go test ./... -count=1` → rc=0 and
  `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` → rc=0. **Base: both green (V17, V18).**

Explicitly rejected as an AC: "the full verify gate is green" alone — a package-wide `ok` can print
while the named test never ran (the `[no tests to run]` form is rc=0); AC1's run-existence form is
the binding version.

## Non-Vacuity — named RED mutation for every added assertion

Production side mutated (the canary file) — never the test helper. Assertion coverage:
A1-func-exists-once←M3/M8, A2-top-level-assertion-if-exists-once←M1/M2/M5/M6/M10/M11/M13, A3-direct-fatalf-in-body-exists-once←M4/M7/M12,
A4-parse-floor←M9. Each arm is **downstream of its named assertion, not merely suite-wide**: the
gutted mutant passes every OLD check (V7), so the only thing that can RED it is the new shape
assertion; and each arm's failure message names the specific component that deviated. M10/M11/M12/M13
are the quorum objection-1 additions: the old loose shape returned `SHAPE OK` on M10 and M12 (V33),
so these arms prove the bound shape is necessary.

| # | Exact edit (to `host/store/toolchain_canary_test.go`) | Expected RED (named assertion) | Shape |
|---|---|---|---|
| M1 | **THE RECORDED GUTTED MUTANT**: delete `:52-53` (`if rows[0].field != "stateRoot" { t.Fatalf(...) }`), replace with `// stateRoot is the expected field value.` | `canaryAssertionShapeProblems`: `rows[0].field != "stateRoot"` assertion if-stmt count=0 | **threat-shaped: the filed repro (V5)** — builds/typechecks (V6), passes all OLD checks (V7), only the shape catches it (V10) |
| M2 | Replace the assertion with a prose comment naming `stateRoot` (no code) | same: assertion if-stmt count=0 | prose/comment token must not satisfy the fence (V14) |
| M3 | Rename `TestToolchainCanary` → `TestToolchainCanaryRenamed` | same: func decl count=0 | anti-vacuity: the target function must exist exactly once (V11) |
| M4 | Remove the `t.Fatalf` from the assertion body, keep the comparison (`_ = rows[0].field`) | same: direct t.Fatalf expression statement in assertion body count=0 | the failure action is part of the assertion (V12) |
| M5 | Duplicate the assertion (two `rows[0].field != "stateRoot"` ifs) | same: assertion if-stmt count=2 | **ambiguous shapes rejected loudly** (V13) |
| M6 | Invert the comparison to `== "stateRoot"` | same: assertion if-stmt count=0 (Op != NEQ) | an inverted assertion fires on the wrong condition (V15) |
| M7 | Move the `t.Fatalf` outside the if-body (keep comparison) | same: direct t.Fatalf expression statement in assertion body count=0 | the Fatalf must be the assertion's action, not a stray call (V16) |
| M8 | Delete the whole `TestToolchainCanary` function | same: func decl count=0 | the canary itself must exist (V11's stronger form) |
| M9 | Corrupt the file so it no longer parses (e.g. truncate mid-statement) | same: `parse error` | **the parser/instrument floor fires** — a parse error is a RED, never a silent pass |
| M10 | **Constant-vs-constant no-op**: `if "stateRoot" != "stateRoot" { t.Fatalf(...) }` | same: assertion if-stmt count=0 (left operand not `rows[0].field`) | **quorum objection 1**: always-false comparison must not satisfy the shape (V27; old loose shape accepted it, V33) |
| M11 | **Wrong left operand**: `if rows[0].n != "stateRoot" { t.Fatalf(...) }` | same: assertion if-stmt count=0 (left operand not `rows[0].field`) | **quorum objection 1**: the comparison must check the field, not another member (V28) |
| M12 | **Fatalf nested under `if false`**: `if rows[0].field != "stateRoot" { if false { t.Fatalf(...) } }` | same: direct t.Fatalf expression statement in assertion body count=0 | **quorum objection 1**: a descendant Fatalf that never fires must not satisfy the shape (V29; old loose shape accepted it, V33) |
| M13 | **Entire assertion nested under `if false`**: `if false { if rows[0].field != "stateRoot" { t.Fatalf(...) } }` | `top-level rows[0].field != "stateRoot" assertion if-stmt count=0` | **quorum round-2 objection, applied verbatim under the narrow-refinement carve-out**: recursive discovery would accept a canonical assertion whose ancestor makes it unreachable; only a direct function-body statement counts |

Green control for all arms: the unmutated post-sprint tree passes AC1/AC2/AC5, and every arm ends
restored sha256-byte-identical with `git status --porcelain` empty — the recipe V5 already ran once
at base (`a23cfa79…` before and after).

## Resolved Decisions

Round-0 carried two open decisions with controller defaults. Quorum round-1 closed them; no OPEN
human decision remains — these are concrete reviewer-authored fixes, not a direction dispute.

- **RD-1 — the comparison's left operand is bound to `rows[0].field`.** *Resolved: YES (mandated by
  quorum objection 1).* The round-0 default-NO is removed. The left operand must structurally be
  `rows[0].field` (a `SelectorExpr` on an `IndexExpr` with index `0`), so a constant-vs-constant
  no-op and a wrong-field comparison are rejected (M10, M11). The cost — a legitimate rename of
  the local `rows` or the `field` selector REDs the fence — is accepted and declared as a residual
  (a shape change must update the fence in the same commit).
- **RD-2 — the `t.Fatalf`'s message is not required to contain `"stateRoot"`.** *Resolved: NO.* The
  message text is prose; the assertion's action is the Fatalf call itself, and the comparison
  already pins the `"stateRoot"` token. Requiring the message would re-introduce a token count
  inside the shape. No reviewer raised this; it is closed as a concrete decision, not left open.

## Axiom Compliance

| Axiom | Score | Justification |
|-------|-------|---------------|
| A1: Determinism | +1 | The AST check is deterministic — same source, same verdict; no runtime, no environment |
| A2: Replayability | 0 | No state written |
| A3: Effect Legibility | 0 | No effects |
| A4: Explicit Authority | 0 | No capability changes |
| A5: Bounded Verification | +1 | Static, bounded parse; no subprocess, no network, no solver |
| A6: Safe Concurrency | 0 | No concurrency impact |
| A7: Machines First | +1 | A structural check replaces fragile text matching; cheaper for AI to reason about |
| A8: Minimal Syntax | 0 | No new syntax |
| A9: Cost Visibility | 0 | No resource changes |
| A10: Composability | 0 | Localized to one test function |
| A11: Structured Failure | +1 | Each shape deviation is a named, specific failure message |
| A12: System Boundary | 0 | No boundary changes |

**Net Score: +4** ✅ Proceed to implementation. **Hard Violation Check**: A1 (no nondeterminism),
A3 (no hidden effects), A4 (no ambient access), A7 (not optimizing for human convenience over
machine analysis) — none violated.

## Implementation Milestones

**Phase 1 — the shape helper** (~0.5h): add `canaryAssertionShapeProblems`, `isRowsField`,
`isStateRootLit`, `isTFatalfCall`, and the four imports to `host/verifygate/toolchain_pin_gate_test.go`;
verify the package compiles (`go build ./host/verifygate/`).

**Phase 2 — wire the assertion** (~0.25h): replace the `:374-376` clause in
`TestCanaryDeclaresPositiveArmOnly` with the shape-assertion call; run AC1/AC2 scoped.

**Phase 3 — mutation drill** (~0.5h): run M1–M13, each restored sha256-byte-identical, porcelain 0
between arms; confirm each REDs the named assertion and the pristine control stays green.

**Phase 4 — full gate** (~0.5h): AC5 hygiene, AC6 full host gate + `verify_ail.sh` with the pinned
binary.

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| A future legitimate refactor changes the assertion shape and REDs the fence (false positive) | Medium | Declared residual; the shape is pinned to the canonical form on purpose; a shape change must update the fence in the same commit. Binding the left operand to `rows[0].field` (RD-1) adds a rename of `rows`/`field` to this class — accepted and declared |
| The AST check is more code than the token count it replaces | Low | ~50 LOC, four stdlib imports, no new dependency; the brief explicitly prefers a small structural check |
| A semantically-equivalent but differently-shaped assertion evades the shape | Low | The essential components (function, `rows[0].field != "stateRoot"`, direct Fatalf in body) are the assertion's real skeleton; a rewrite that keeps them passes correctly |
| The full host gate is slow (~76s verifygate) | Low | AC1/AC2/AC3 run scoped; AC6 is the final-tree gate, not per-arm |

## Verification Log

All rows run first-party by the designer at `0bbb1a9` (clean tree), shell `zsh`,
`PATH=/opt/homebrew/bin:$PATH`, darwin/arm64, 2026-08-30. KP = known-positive control carried in
the same call. V9–V16 and V27–V29/V33 ran the Decision's helper verbatim as a standalone prototype
in `/tmp/astproto2` (parsing the real canary file and the mutant copies), deleted before the
porcelain checks; the sprint re-runs all arms against the landed test. V30–V32 are controller
census/citation rows recorded verbatim.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Worktree is `0bbb1a9`, clean; every mutation arm below ended restored | `git rev-parse HEAD && git status --porcelain \| wc -l` (re-checked after every arm) | `0bbb1a96603fa75279b8f9f55e9d2fe922fb6a2c`, `0` |
| V2 | Toolchain boundary | `go version` at repo root | `go1.26.6 darwin/arm64` |
| V3 | Pristine fence green | `GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run '^TestCanaryDeclaresPositiveArmOnly$' -count=1 -v` | rc=0, 1 `=== RUN`, 1 `--- PASS` |
| V4 | Pristine canary green | `GOTOOLCHAIN=go1.26.6 go test ./host/store/ -run '^TestToolchainCanary$' -count=1 -v` | rc=0, 1 `--- PASS` |
| V5 | **THE REPRO: gutted canary is green everywhere** | `cp` backup; delete `:52-53`, replace with `// stateRoot is the expected field value.`; fence + canary + scoped suite; `cp` restore + sha256 + porcelain | `stateRoot` count 3→2; fence `--- PASS`; canary `--- PASS`; scoped `go test ./host/verifygate/ ./host/store/ -count=1` fails **only** on the documented `AILANG_BIN`-unset module-manifest shim tests (store `ok`); restore `a23cfa79419ae69136e62981d5bd0c8ea68cdf2e154ba973fb99a3d3c8b47bfd` byte-identical, porcelain 0 |
| V6 | The gutted mutant builds/typechecks | gutted file in place; `GOTOOLCHAIN=go1.26.6 go vet ./host/store/`; `go build ./host/store/` | vet rc=0; build rc=0 — the mutant is a valid Go test that asserts nothing |
| V7 | The gutted mutant passes every OLD check | `grep -c` per needle on the gutted copy | `GOTOOLCHAIN` **0** (zero-needle green), `POSITIVE ARM ONLY` **1** (marker green), `stateRoot` **2** (old `>= 2` green) — only a shape check can catch it |
| V8 | `go/parser`/`go/ast`/`go/token`/`fmt` are stdlib | `go doc go/parser`, `go doc go/ast` | both present, stdlib — no new dependency |
| V9 | Prototype: pristine shape holds | `/tmp/astproto2/astproto2 <pristine canary>` | `SHAPE OK`, rc=0 |
| V10 | Prototype: gutted mutant REDs | same on `/tmp/astproto2/m1_gutted.go` | `PROBLEM: \`rows[0].field != "stateRoot"\` assertion if-stmt count=0, want exactly 1`, rc=1 |
| V11 | Prototype: function-missing REDs | same on `/tmp/astproto2/m3_rename.go` (renamed func) | `PROBLEM: TestToolchainCanary func decl count=0, want exactly 1`, rc=1 |
| V12 | Prototype: Fatalf-missing REDs | same on `/tmp/astproto2/m4_nofatalf.go` (comparison kept, Fatalf removed) | `PROBLEM: direct t.Fatalf expression statement in assertion body count=0, want exactly 1`, rc=1 |
| V13 | Prototype: ambiguous REDs | same on `/tmp/astproto2/m5_ambig.go` (two assertions) | `PROBLEM: \`rows[0].field != "stateRoot"\` assertion if-stmt count=2, want exactly 1`, rc=1 |
| V14 | Prototype: prose/comment-only REDs | same on `/tmp/astproto2/m2_prose.go` (comment naming stateRoot, no code) | `PROBLEM: \`rows[0].field != "stateRoot"\` assertion if-stmt count=0, want exactly 1`, rc=1 |
| V15 | Prototype: inverted comparison REDs | same on `/tmp/astproto2/m6_invert.go` (`==` instead of `!=`) | `PROBLEM: \`rows[0].field != "stateRoot"\` assertion if-stmt count=0, want exactly 1`, rc=1 |
| V16 | Prototype: Fatalf moved outside if-body REDs | same on `/tmp/astproto2/m7_move.go` | `PROBLEM: direct t.Fatalf expression statement in assertion body count=0, want exactly 1`, rc=1 |
| V17 | Full host gate green on pristine origin/dev | `AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.26.6 go test ./... -count=1` | rc=0, all 19 packages `ok` (incl. verifygate 76.449s, store 8.405s) |
| V18 | `verify_ail.sh` green on pristine | `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` | rc=0, `verify gate PASSED: 11 required identities verified, 40 named tests pass` |
| V19 | Scoped suite green on pristine with pinned binary | `AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ ./host/store/ -count=1` | rc=0, both `ok` |
| V20 | `go build ./...` green | `GOTOOLCHAIN=go1.26.6 go build ./...` | rc=0 |
| V21 | gofmt clean | `gofmt -l host/verifygate/ host/store/ \| wc -l` | `0` |
| V22 | go vet clean | `GOTOOLCHAIN=go1.26.6 go vet ./host/verifygate/ ./host/store/` | rc=0 |
| V23 | Systemic: only one gate reads the canary file | `grep -rn "toolchain_canary_test" --include="*.go" .` (excl. design_docs); KP `grep -rn "verify_go.sh" --include="*.go" host/` and `grep -rn "run.sh" --include="*.go" host/verifygate/` | exactly one Go hit: `toolchain_pin_gate_test.go:368`; KP: `verify_go.sh` referenced in `evidence_manifest_gate_test.go`, `run.sh` in `toolchain_pin_gate_test.go` — grep live. `verify_go.sh`/`ci.yml` do not read the canary file |
| V24 | `stateRoot` occurrences in the canary file | `grep -n "stateRoot" host/store/toolchain_canary_test.go` | exactly 3: `:32` array literal, `:52` comparison, `:53` Fatalf message |
| V25 | Token-count positive controls under `host/verifygate` | `grep -rh "strings.Count" host/verifygate/*_test.go \| wc -l`; per-file `grep -rc` | **23** call-sites across 4 files (ail_binary 8, floor_raise 4, module_manifest 4, toolchain_pin 7); the `:374` clause is the only one guarding an assertion |
| V26 | No name collision for the new helper/predicates | `grep -h "^func " host/verifygate/*_test.go \| grep -i "canaryAssertionShape\|isRowsField\|isStateRootLit\|isTFatalfCall"` | no hits — names are free |
| V27 | Prototype: constant-vs-constant no-op REDs | same on `/tmp/astproto2/m10_const.go` (`if "stateRoot" != "stateRoot" { t.Fatalf(...) }`) | `PROBLEM: \`rows[0].field != "stateRoot"\` assertion if-stmt count=0, want exactly 1`, rc=1 — **quorum objection 1** |
| V28 | Prototype: wrong left operand REDs | same on `/tmp/astproto2/m11_other.go` (`if rows[0].n != "stateRoot" { t.Fatalf(...) }`) | `PROBLEM: \`rows[0].field != "stateRoot"\` assertion if-stmt count=0, want exactly 1`, rc=1 — **quorum objection 1** |
| V29 | Prototype: Fatalf nested under `if false` REDs | same on `/tmp/astproto2/m12_nested.go` (`if rows[0].field != "stateRoot" { if false { t.Fatalf(...) } }`) | `PROBLEM: direct t.Fatalf expression statement in assertion body count=0, want exactly 1`, rc=1 — **quorum objection 1** |
| V30 | **Precise two-line citation for the version-agnostic claim** | `sed -n '214,216p' scripts/verify_go.sh` | line 214 `# This deny-list is the measured set: go1.26.0-go1.26.5 on darwin/arm64.`; line 215 `# Future go1.26.6 or go1.27.x versions are not covered here; the canary in this`; line 216 `# gate is the version-agnostic detector for any version that miscompiles the shape.` — **quorum objection 2** |
| V31 | **Controller census of the detector surface, with scope/limitation** | `rg -n -i 'miscompil\|array[- ]literal\|toolchain canary\|ToolchainCanary' --glob '*.go' --glob '*.sh' --glob '*.yml' --glob '!design_docs/**' .` | hits in `scripts/verify_go.sh` (`:216`, `:220-221`), `host/store/toolchain_canary_test.go`, and `host/verifygate/toolchain_pin_gate_test.go` (`TestMiscompileInstrumentProbesPinnedToolchain` `:233`, `TestMiscompileInstrumentStepIsGatedInCI` `:400`, `miscompileReproducerPath` `:385`). The `!design_docs/**` exclusion is too broad — the nested runtime reproducer deliberately lives under `design_docs`. **Narrowed thesis:** the canary is the version-agnostic positive-arm detector; the nested module/`run.sh` is the known-bad runtime arm; CI wiring is pinned separately — **quorum objection 2** |
| V32 | Existing identifier `canaryPath` confirmed (Gemini) | `grep -n "canaryPath" host/verifygate/toolchain_pin_gate_test.go` | present at `:368` (`canaryPath := filepath.Join(repoRoot, "host", "store", "toolchain_canary_test.go")`) and used in the `:375`/`:378`/`:381` messages — the new shape-assertion call reuses it |
| V33 | **The old loose shape accepts the two no-op mutants (why the bound shape is necessary)** | loose prototype (any left operand + recursive Fatalf) on `/tmp/astproto2/m10_const.go` and `m12_nested.go` | both `SHAPE OK`, rc=0 — the round-0 shape would have let the constant-vs-constant no-op and the `if false`-hidden Fatalf through; the revised shape REDs them (V27, V29) — **quorum objection 1** |
| V34 | **Round-2 carve-out: the assertion itself must be top-level** | structural inspection of the revised helper plus sprint mutation M13 | The helper iterates only `funcs[0].Body.List`; a nested assertion is not enumerated. M13 is a mandatory sprint arm and must build/typecheck, leave the old token checks and canary green, and RED the new fence before this milestone can land. No pre-sprint execution result is claimed. |

## Related Documents

- [`../implemented/w-canary-control-does-not-survive-a-floor-raise.md`](../implemented/w-canary-control-does-not-survive-a-floor-raise.md)
  — row 42, the parent: built `TestCanaryDeclaresPositiveArmOnly` and its `:374` token-count
  control; its residual *"Test B fences one token"* (`:261`) is the exact gap this item closes.
- [`../implemented/w-race-gate-blindspot.md`](../implemented/w-race-gate-blindspot.md) — where the
  canary and the nested-module pattern were born; the canary is the artifact this fence protects.
- [`../implemented/w-setup-go-pin-unguarded.md`](../implemented/w-setup-go-pin-unguarded.md) — row
  41, built the file this item extends and the helpers it reuses.
- [`../implemented/w-miscompile-instrument-inert-in-ci.md`](../implemented/w-miscompile-instrument-inert-in-ci.md)
  — row 44, the runtime lane caveat this item declares rather than fixes.
- [`../planned/w-racecontrol-floor-bump-disarms-the-race-control.md`](../planned/w-racecontrol-floor-bump-disarms-the-race-control.md)
  — row 48, the sibling nested-module exposure; distinct (different module, different property).
- `design_docs/world-mission.md` queue row 49 (this item) and rows 42, 44, 45, 48.
