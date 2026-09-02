# w-canary-fence-blind-to-a-skipped-canary — the row-49 AST fence proves the canary *contains* its assertion but cannot see that the assertion never *runs*, because `t.Skip()` as the first statement is invisible to all three of the fence's existing assertions, so the fence stays green over a canary that asserts nothing

**Status**: Planned
**Date**: 2026-09-03
**Queue item**: 56, `w-canary-fence-blind-to-a-skipped-canary` (clause-2, evaluator-found as P49
adversarial arm A10, controller-reproduced first-party before adoption)
**Estimated**: ~0.3 day (three reachability checks folded into the existing AST pass in
`host/verifygate/toolchain_pin_gate_test.go`; **no production code, no `host/store` assertion, no
`run.sh`, no `ci.yml`, no `scripts/` edit, no new package**)
**Designer**: `pi:ollama/deepseek-v4-flash:0731-cloud` (design-doc-creator, iteration 150)
**Toolchain boundary**: every command below was run first-party in this worktree at `d75b9c3`
(clean tree; porcelain 0 re-checked after every mutation arm), shell `zsh`,
`PATH=/opt/homebrew/bin:$PATH`, darwin/arm64, `go version` = `go1.26.6 darwin/arm64` (V2). No
AILANG (`.ail`) source is written or changed by this design; the change is pure Go test code, so no
ailang check/ai-check gate applies.

> **Thesis:** row 49's `TestCanaryDeclaresPositiveArmOnly` (`host/verifygate/toolchain_pin_gate_test.go:456`)
> reads the canary source and calls `canaryAssertionShapeProblems` (`:375`), which makes exactly
> three AST assertions over `host/store/toolchain_canary_test.go`: `TestToolchainCanary` exists
> exactly once, exactly one top-level `if rows[0].field != "stateRoot"` assertion if-stmt, and
> exactly one **direct** `t.Fatalf` expression statement in that if-body. **None of these looks at
> reachability.** Measured first-party this session (V9): inserting `t.Skip("MUTANT: reachability
> probe")` as the **first statement** of `TestToolchainCanary`'s body leaves every shape the fence
> checks unchanged — the func-decl count, the top-level `rows[0].field != "stateRoot"` if-stmt, and
> the direct `t.Fatalf` in its body are all still exactly one — so the mutant compiles
> (`go vet ./host/store/` rc=0), the canary itself reports `--- SKIP` (rc=0), and the fence
> `--- PASS` (rc=0). **GREEN over a canary that
> asserts nothing.** The fix is cheap in the shape P49 already built: extend the SAME AST pass with
> three **statically-visible reachability** checks — a zero-needle on `t.Skip`/`t.Skipf`/`t.SkipNow` calls in
> the body, a zero-needle on `return` statements in the body (an early return), and a zero-needle on
> `//go:build`/`// +build` constraints on the file. **Scope discipline the row must keep:** this
> closes *statically visible* reachability only. It does NOT make the fence a proof that the
> assertion fires on a miscompiling toolchain — that lane is the nested repro module plus `run.sh`'s
> `saw_bad` floor, which is darwin-only and attended while `ci.yml:172` `continue-on-error: true`
> stands (row 44). This row does not grow into that one.

## The finding in one paragraph

`TestCanaryDeclaresPositiveArmOnly` (`host/verifygate/toolchain_pin_gate_test.go:456`) reads the
canary file by path and calls `canaryAssertionShapeProblems` (`:375`), which parses the source and
asserts the *shape* of the assertion: `TestToolchainCanary` exists exactly once, exactly one
top-level `if rows[0].field != "stateRoot"` if-stmt, exactly one direct `t.Fatalf` in that if-body.
The shape is a *containment* property — it proves the canary **contains** the assertion skeleton,
never that the assertion **runs**. A `t.Skip()` placed as the first statement of the body is a
statement the shape pass does not enumerate (it only counts top-level if-stmts and the direct
Fatalf), so it is invisible to all three of the fence's existing assertions and the fence stays green. The repro is
exact and first-party (V9): with `t.Skip("MUTANT: reachability probe")` inserted as the first
statement, the mutant compiles (`go vet ./host/store/` rc=0), the canary reports `--- SKIP` (rc=0),
and the fence `--- PASS` (rc=0) — a canary that asserts nothing is green everywhere. The same
blindness holds for the sibling runtime-neutering forms the judge named: `t.SkipNow()`, an early
`return`, and a build tag all pass the current fence (V10); the `t.Skipf` form is the
controller-measured hole the widened A5 closes (V27). The repair is three reachability
zero-needles folded into the same AST pass: no `t.Skip`/`t.Skipf`/`t.SkipNow` call in the body, no `return`
statement in the body, no `//go:build`/`// +build` constraint on the file. This does **not** falsify
row 49's claim — row 49's own "What the gate CANNOT catch — declared residual" section
(`w-canary-fence-passes-a-gutted-canary.md:314-316`) states this exact class up front: *"The shape
is not the behaviour. The AST check proves the canary contains the assertion shape; it cannot prove
the assertion fires."* This item closes a **disclosed residual**, it does not fix a defect P49 hid.

## Premises

Each premise is one or more Verification Log rows; a claim without a row does not appear here.

- **P1 — the fence is green on the pristine tree**: `TestCanaryDeclaresPositiveArmOnly` → `--- PASS`
  (V7); the canary `TestToolchainCanary` → `--- PASS` (V8).
- **P2 — the defect reproduces exactly as filed**: inserting `t.Skip("MUTANT: reachability probe")`
  as the first statement of `TestToolchainCanary`'s body leaves every shape the fence checks
  unchanged — the func-decl count, the top-level `rows[0].field != "stateRoot"` if-stmt, and the
  direct `t.Fatalf` in its body are all still exactly one (V9); the mutant compiles
  (`go vet ./host/store/` rc=0), the canary reports
  `--- SKIP` (rc=0), and the fence `--- PASS` (rc=0) (V9). Restore was sha256-byte-identical
  (`a23cfa79…`), porcelain 0.
- **P3 — the current fence is blind to all five neutering forms**: `t.Skip`, `t.Skipf`, `t.SkipNow`, an early
  `return`, and a `//go:build` constraint each pass the current fence (rc=0) (V10; the `t.Skipf` hole is
  controller-measured in V27). The RED can
  therefore only come from a check that inspects **reachability**, not shape.
- **P4 — the canary file carries no build constraint today**: its first three lines are
  `package store`, blank, `import "testing"` (V4); a zero-needle on `//go:build`/`// +build` is
  green at base.
- **P5 — the reachability zero-needles are green at base**: `t.Skip` count 0 (known-positive control
  in the same call: `t.Fatalf` count 3) (V5); `GOTOOLCHAIN` count 0 (control: `POSITIVE ARM ONLY`
  count 1) (V6); the prototype reports `skip=0 ret=0 buildtag=0` on the pristine canary (V11).
- **P6 — the extended checks discriminate**: the prototype REDs `t.Skip` (V12), `t.SkipNow` (V13),
  an early `return` (V14), and a `//go:build` constraint (V15); the widened A5 also REDs `t.Skipf`
  (V27); it does NOT catch a `t.Skip` via an
  alias `s := t.Skip; s(...)` (V16) or a helper `helperSkip(t)` (V17) — the declared residuals.
- **P7 — every mutant is a real, buildable program** — verified by `go vet` rc=0 for the
  `t.Skip`/`t.Skipf`/`t.SkipNow`/`//go:build` mutants and by `go test -count=1 -run '^$'` rc=0 for the
  `M-RETURN` mutant, whose `go vet` rc=1 is an ANALYZER finding (`unreachable code`), not a compile
  failure (**V34**; V32, V21). The RED the sprint asserts comes from the fence, not from a compile
  error, for all five. **V18's `rc=0 for all four` cell is SUPERSEDED and is wrong on its `m_return`
  arm** — see V34, which re-ran the contradiction as a third, isolated measurement rather than
  adjudicating between the two earlier rows.
- **P8 — the change surface is one test file, no new imports**: `fmt`, `go/ast`, `go/parser`,
  `go/token`, `strings` are all already imported in `toolchain_pin_gate_test.go` (V22); the three
  checks fold into the existing `canaryAssertionShapeProblems` pass.

### Design Freeze

- **One test function's helper extended, no test file created**: `canaryAssertionShapeProblems`
  (`toolchain_pin_gate_test.go:375`) gains three reachability checks appended after its existing
  direct-Fatalf check; its final `return nil` becomes `return problems`. The sprint must declare
  `var problems []string` ahead of the appended checks (the helper currently has no such variable —
  it returns a fresh `[]string{…}` on each early failure and `nil` at the end, V28). The three
  EXISTING early returns (func-decl count, assertion-if count, direct-Fatalf count) are left
  byte-unchanged, so an early return still short-circuits and the new checks only run when the shape
  checks passed. No new helper functions, no new imports, no new package.
- **The three existing shape assertions (func-exists-once, top-level assertion if, direct Fatalf)
  are retained byte-unchanged** — they are the containment proof; the three new checks add the
  reachability proof.
- **The `GOTOOLCHAIN` zero-needle and the `POSITIVE ARM ONLY` marker are retained byte-unchanged** —
  separate fences for separate purposes.
- **No production code, no `host/store` assertion, no `run.sh`, no `ci.yml`, no `scripts/` edit**:
  the canary file is read-only (by path) from `host/verifygate`, exactly as today.
- **The reachability checks are structural, not textual**: they walk the AST of
  `TestToolchainCanary`'s body (`ast.Inspect`) and the file's comment groups, so a `t.Skip` in a
  comment or string literal is not a call and does not red; a `t.Skip` in a *different* function in
  the file does not neuter `TestToolchainCanary` and does not red. Nested function literals are
  handled per check (Decision, *Nested function literals*): A5 descends into them (a `t.Skip` on
  the outer `t` inside a closure still skips the test, V30); A6 stops at `*ast.FuncLit` (a `return`
  inside a closure exits the closure, not the test, V30).
- **Anti-vacuity is built in**: each reachability check is a zero-needle (count must be 0); a count
  >0 is a loud RED with a named message. The parser itself is a floor: a parse error is a RED, never
  a silent pass.

## Decision — three statically-visible reachability zero-needles fold into the existing AST pass

The queue row names the mechanism: *"a zero-needle on `t.Skip` / `t.SkipNow` — plus the sibling
runtime-neutering forms the judge named (`t.SkipNow()`, an early `return`, a build tag) — folds into
the SAME AST pass that already walks `TestToolchainCanary`'s body."* **This design does exactly
that.** It extends `canaryAssertionShapeProblems` with three checks, all scoped to the function body
(or, for the build tag, the file's comment groups), all zero-needles, all structural.

**The A5 selector set is CLOSED by construction against `go doc testing.T`'s three skip methods.**
`testing.T` exposes exactly `Skip(args ...any)`, `Skipf(format string, args ...any)`, and
`SkipNow()` (V27) — a matcher over a fixed stdlib surface is justified by enumerating that surface,
not by naming the forms someone happened to think of. A5 therefore matches all three; there is no
fourth skip form on `testing.T` to miss. The pristine control (V27) proves the widened matcher does
not over-match: it still counts 0 on the unmutated canary.

### The three new checks

Append to `canaryAssertionShapeProblems` after its existing direct-Fatalf check, and change the
final `return nil` to `return problems`. The block below is byte-for-byte what V32 patched into the
REAL helper in situ (with `:377` at `parser.ParseComments`): `go vet ./host/verifygate/` rc=0,
`gofmt -l` empty, pristine fence `--- PASS`, then restored byte-identical (V33):

```go
	// A5 — reachability: no t.Skip / t.Skipf / t.SkipNow call anywhere in the body, INCLUDING
	// inside nested func literals. A t.Skip on the outer `t` inside a closure still Goexits
	// TestToolchainCanary (V30: the assertion after it never runs), so descending is not merely
	// conservative — it is correct. The selector set is CLOSED by construction against
	// `go doc testing.T`'s three skip methods (Skip, Skipf, SkipNow) — see V27.
	var problems []string
	skipCalls := 0
	ast.Inspect(funcs[0].Body, func(n ast.Node) bool {
		if c, ok := n.(*ast.CallExpr); ok {
			if sel, ok := c.Fun.(*ast.SelectorExpr); ok {
				if id, ok := sel.X.(*ast.Ident); ok && id.Name == "t" &&
					(sel.Sel.Name == "Skip" || sel.Sel.Name == "Skipf" || sel.Sel.Name == "SkipNow") {
					skipCalls++
				}
			}
		}
		return true
	})
	if skipCalls != 0 {
		problems = append(problems, fmt.Sprintf("t.Skip/t.Skipf/t.SkipNow call count=%d, want 0 (a skipped canary asserts nothing)", skipCalls))
	}

	// A6 — reachability: no early return in the body. The assertion is the last statement,
	// so any return before it neuters it. Traversal STOPS at *ast.FuncLit: a return inside a
	// nested func literal exits the literal, not TestToolchainCanary (V30), so counting it
	// would false-red a canary whose assertion still runs (V31, the row-55 class).
	returns := 0
	ast.Inspect(funcs[0].Body, func(n ast.Node) bool {
		if _, ok := n.(*ast.FuncLit); ok {
			return false
		}
		if _, ok := n.(*ast.ReturnStmt); ok {
			returns++
		}
		return true
	})
	if returns != 0 {
		problems = append(problems, fmt.Sprintf("return statement count=%d, want 0 (an early return neuters the assertion)", returns))
	}

	// A7 — reachability: no build constraint on the file. A build tag can exclude the
	// canary from the build entirely, so the assertion never runs.
	buildTags := 0
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			if strings.HasPrefix(c.Text, "//go:build") || strings.HasPrefix(c.Text, "// +build") {
				buildTags++
			}
		}
	}
	if buildTags != 0 {
		problems = append(problems, fmt.Sprintf("build constraint count=%d, want 0 (a build tag can exclude the canary from the build)", buildTags))
	}

	return problems
```

**Why these three and no more.** The row's scope discipline is explicit: *"this closes *statically
visible* reachability only."* The three forms are the ones the judge named and the ones that are
statically visible in the AST without dataflow or interprocedural analysis: a direct `t.Skip`/
`t.Skipf`/`t.SkipNow` call, a direct `return`, and a file-level build constraint. Each is a zero-needle
because the canary currently contains none of them (P4, P5). Each is structural (AST) rather than
textual so a `t.Skip` in a comment or string literal does not false-red, and a `t.Skip` in a
different function in the file does not false-red (it does not neuter `TestToolchainCanary`).

**Why the checks are scoped to the function body.** A `t.Skip` in a helper or another test function
in the same file does not prevent `TestToolchainCanary`'s assertion from running; only a neutering
form *inside* `TestToolchainCanary`'s body (or a build tag that excludes the whole file) does. The
build-tag check is file-level because a constraint is a file property.

**Nested function literals — A5 and A6 deliberately take DIFFERENT traversal policies (quorum
round 1, gpt5-6-sol's second catch).** `ast.Inspect(funcs[0].Body, …)` descends into `*ast.FuncLit`
nodes, so an undifferentiated walk would count a `return` or a `t.Skip` that sits inside a closure.
The reviewer offered two dispositions — document the conservative descend policy explicitly, or
stop traversal at `*ast.FuncLit` — and the two checks do not want the same one, because the two
forms do different things at runtime. Measured in a throwaway module (V30): `func() { return }()`
followed by `t.Fatalf("REACHED…")` → the Fatalf prints and the test `--- FAIL`s, i.e. a
closure-nested `return` exits the *literal* and the assertion after it still runs;
`func() { t.Skip(…) }()` followed by the same Fatalf → `--- SKIP` and the Fatalf line never prints,
i.e. a closure-nested `t.Skip` on the outer `t` DOES neuter the test (`t.Skip` is `runtime.Goexit`
on that `t`, whichever frame calls it). Therefore: **A6 stops at `*ast.FuncLit`** (reviewer option
2) — counting a closure-nested `return` would false-red a canary whose assertion still fires, which
is exactly the row-55 class (a gate red on valid code) this mission landed one iteration ago;
**A5 descends** (reviewer option 1, stated here as the explicit policy) — and it is not merely
conservative, it is correct, since the skip genuinely reaches `TestToolchainCanary`. Both
dispositions are zero-needles TODAY and neither can regress the pristine tree: the canary file
contains zero func literals (`grep -c 'func('` → `0`; AST `FuncLit` count 0, V29), so on the
pristine canary the descend and stop policies read identical counts (A5 0/0, A6 0/0, V29). The
known-positive-adjacent control that the body scoping is real: the file DOES contain exactly one
`return` — `func (s canaryString) String() string { return string(s) }` at `:59`, outside
`TestToolchainCanary` — so the file-level `ReturnStmt` count is 1 while the body-scoped A6 count
must read 0, and does (V29). The policy split is load-bearing, not cosmetic: with the real helper
patched in situ, `M-RETURN-FUNCLIT` (`func() { return }()` as the first statement) leaves the canary
`--- PASS` and the shipped fence `--- PASS` (V32), while the same mutant under an A6 that descends
REDs the fence with `return statement count=1` — a false red over a canary that still asserts
(V33). `M-SKIP-FUNCLIT` (`func() { t.Skip(…) }()`) makes the canary `--- SKIP` and the shipped
fence REDs it with `t.Skip/t.Skipf/t.SkipNow call count=1` (V32); a stop-at-FuncLit A5 would have
read 0 and MISSED it (V31).

**Why the build-tag check reads `f.Comments` — AND WHY THE SPRINT MUST CHANGE THE PARSE MODE IN THE
SAME COMMIT.** `go/parser` populates `f.Comments` **only** under `parser.ParseComments`; a
`//go:build` or `// +build` directive is a comment and appears there (V15). **CONTROLLER CORRECTION
(V25, measured 2026-09-02, and it refutes this paragraph's first draft): the helper does NOT already
parse that way.** `host/verifygate/toolchain_pin_gate_test.go:377` reads
`parser.ParseFile(fset, "", src, 0)` — mode **`0`**, so `f.Comments` is **empty**, and check A7 as
drafted **can never fire**: it would append no problem for a canary carrying `//go:build ignore`,
i.e. it would be a vacuous assertion wearing a gate's clothes. Measured in two arms with a
known-positive control (V25): on the pristine canary, mode `0` yields **0** comments while
`ParseComments` yields **17**; on the same file prefixed with `//go:build ignore`, mode `0` still
yields **0** while `ParseComments` yields **18**. The designer's V15 prototype ran as a standalone
program that did not inherit the helper's mode, which is why the prototype fired and the integration
would not — a same-shape-different-scope instrument failure.
**Therefore A7 is conditional on a fourth, mandatory edit:** change `:377` to
`parser.ParseFile(fset, "", src, parser.ParseComments)`. That mode is already used elsewhere in this
repo (`cmd/world-publish/wiring_test.go:60`, `host/broker/invoke_boundary_test.go:112,:317`), so it
is not novel here. The parser does not *respect* build constraints (it parses regardless), so the
existing shape checks A1–A4 still run unchanged under the wider mode — the sprint must prove that
with a pristine-control arm, not assume it. **AC and Non-Vacuity consequence:** the A7 mutant arm
(`//go:build ignore`) is only a valid non-vacuity proof if it is run *after* the mode change; run it
in both modes and require the verdicts to DIFFER, so the arm cannot pass for the wrong reason.

### What the gate CANNOT catch — declared residual

- **The shape is not the behaviour, and neither is static reachability.** The three new checks prove
  the canary *contains* the assertion shape and that no *statically visible* neutering form precedes
  it. They do not prove the assertion *fires* on a miscompiling toolchain. The runtime proof lane is
  the nested repro module + `run.sh`'s `saw_bad` floor (row 42's Test A and the reproducer), which
  is darwin-only and attended/local while `ci.yml:172` `continue-on-error: true` stands (row 44).
  This item does not change that lane.
- **A helper function that skips (`helperSkip(t)`)** — a call to a helper that internally calls
  `t.Skip`. The AST of `TestToolchainCanary`'s body shows a call to `helperSkip`, not a `t.Skip`
  call, so A5 does not red (V17). Detecting it requires interprocedural analysis. Declared: the cost
  (a whole-program call graph over the canary file) is not worth it for a ~0.3d item, and a helper
  that skips is a deliberate, reviewed shape change that the fence's own contract (update the fence
  in the same commit) already governs.
- **`t.Skip` reached via an alias (`s := t.Skip; s(...)`)** — the call is `s(...)`, an identifier
  call, not a `t.Skip` selector, so A5 does not red (V16). Detecting it requires alias analysis.
  Declared for the same cost reason.
- **A return/skip under a *dynamically* false condition** (e.g. `if rows[0].field == "stateRoot" {
  return }`). A6 catches *any* `return` in the body, including one under a statically-false guard
  (`if false { return }`) — a conservative RED that forces removal of dead code. But a return under
  a condition that is *dynamically* true exactly when the assertion would pass requires dataflow
  analysis. Declared: it is also not a real threat to the miscompile detector — if the toolchain
  miscompiles, `rows[0].field != "stateRoot"`, the guard is false, and the assertion still fires.
- **A closure-nested `return` is NOT a residual — it is a non-threat by language rule.** A6 stops at
  `*ast.FuncLit`, so a `return` inside a func literal in the body is not counted; that is correct,
  not a gap, because a `return` in a literal cannot exit `TestToolchainCanary` (V30: the assertion
  after it still fires; V32: the canary still `--- PASS`es on `M-RETURN-FUNCLIT`). Nothing the stop
  leaves uncounted can neuter the assertion.
- **A closure-nested `t.Skip` is a RED, not a residual.** A5 descends into func literals, so
  `func() { t.Skip(…) }()` reds (V32). If a closure-nested skip is ever written on a *different*
  `*testing.T` (e.g. a subtest's `t` shadowing the outer one), A5 still reds it — a conservative red
  on a shape the canary has no reason to contain (V29: zero func literals today). Declared as the
  explicit conservative policy, per the reviewer's option 1.
- **Other exit forms the judge did not name** — `runtime.Goexit()`, `os.Exit(0)`, a `panic`
  swallowed by a deferred `recover`. Not skips, not returns, not build tags; outside this row's named
  forms and not added here (scope discipline). Declared, not silently assumed away.
- **The whole file being deleted** — already handled by the existing `os.ReadFile` error path in
  `TestCanaryDeclaresPositiveArmOnly` (`t.Fatal(err)`). Covered, not re-scoped.
- **The test being renamed** — already handled by the existing func-decl-count==1 assertion (A1).
  Covered, not re-scoped.
- **A semantically-equivalent but differently-shaped assertion** — inherited from row 49's residual
  (e.g. an `IfStmt` with an `Init`, or a helper-call comparison) and unchanged by this item.

## Alternatives rejected

1. **A textual zero-needle on `t.Skip` over the whole file** (`strings.Count(src, "t.Skip") == 0`).
   Simpler, but fragile: it would red on a `t.Skip` in a comment or string literal, and it would red
   on a `t.Skip` in a *different* function that does not neuter `TestToolchainCanary`. The AST check
   is precise about what is a call and where it lives. Rejected in favour of the structural check.
2. **A behavioural fence — execute the canary and require it to red on a seeded-wrong value.** This
   is the row-44 lane (nested repro module + `run.sh`), which is darwin-only and attended while
   `ci.yml:172` `continue-on-error: true` stands. It is out of scope by the row's own scope
   discipline and by the task brief's scope fence. Rejected.
3. **Constant-folding to detect `if <never true> { return }`.** Requires evaluating conditions, which
   is beyond "statically visible reachability" and beyond a ~0.3d item. A6's conservative "any return
   in the body" already reds the dead-code form; the dynamically-false form is a declared residual.
   Rejected.
4. **A whole-file `t.Skip` zero-needle plus a separate build-tag check, leaving the early-return
   form unguarded.** The judge explicitly named the early `return` as a sibling form; leaving it
   unguarded would be a known hole in the same class. Rejected.

## Ordering

Gated on nothing. Neighbours named and not absorbed: **row 49** (`w-canary-fence-passes-a-gutted-canary.md`)
is the parent — this item closes its declared residual *"The shape is not the behaviour"*
(`:314-316`) and reuses its `canaryAssertionShapeProblems` pass; **row 44** owns `run.sh`/`ci.yml`
inertness — this item changes no runtime behaviour of the instrument; **row 42** owns the canary and
the nested-module pattern — untouched. The sprint applies the single helper extension; no other file
changes.

## Files to Create/Modify

- **MODIFY** `host/verifygate/toolchain_pin_gate_test.go` (+~40 LOC), in **two** places:
  1. `canaryAssertionShapeProblems` gains the three reachability checks; a `var problems []string`
     is **declared** ahead of them (the helper currently has no such variable — it returns a fresh
     `[]string{…}` on each early failure and `nil` at the end, V28), and its final `return nil` becomes
     `return problems`. The three EXISTING early returns (func-decl count, assertion-if count,
     direct-Fatalf count) are left **byte-unchanged** — an early return still short-circuits, so the
     new checks only run when the shape checks passed. (If the sprint prefers the new checks to run
     even when a shape check fails, it must say so and justify it; this design does not.)
  2. **`:377` `parser.ParseFile(fset, "", src, 0)` becomes
     `parser.ParseFile(fset, "", src, parser.ParseComments)`** — MANDATORY, not optional: under mode
     `0` the parser discards comments, so `f.Comments` is empty and check A7 can never fire (V25,
     two arms with a firing control). Shipping A7 without this edit ships a vacuous assertion.
  No new imports (V22), no new helper functions, no new package.
- **No other files.** `host/store/toolchain_canary_test.go` (read-only, by path), `run.sh`,
  `ci.yml`, root `go.mod`, `scripts/*`, the nested repro module — untouched.

## Conflict Surface

This is a Go test/gate change, not an AILANG parser/typechecker/codegen change, so the Conflict
Surface is not mechanically mandatory; it is included because the gate is shared machinery.

- **`canaryAssertionShapeProblems` itself** — the only helper extended. Its three existing shape
  assertions are retained byte-unchanged; the three new checks are appended. The sprint re-confirms
  the retained assertions still pass on the post-sprint tree (AC1).
- **`TestCanaryDeclaresPositiveArmOnly`** — the only test edited (indirectly, via its helper). Its
  `GOTOOLCHAIN` zero-needle and `POSITIVE ARM ONLY` marker are retained byte-unchanged.
- **`TestToolchainCanary` (`host/store/toolchain_canary_test.go`)** — the subject, read-only by
  path. Its assertions are untouched, so its green under the pinned toolchain (V8) is unaffected.
- **`go/ast`, `go/parser`, `go/token`, `fmt`, `strings`** — stdlib, already imported (V22); no new
  dependency, no conflict with the existing `go/version` import.
- **`verify_go.sh` / `ci.yml`** — do not read the canary file; no interaction. The canary is the
  version-agnostic positive-arm detector; `verify_go.sh` and the CI wiring are separate arms.

## Deferred Scope

- **Row 44 `w-miscompile-instrument-inert-in-ci`** — `ci.yml:172` `continue-on-error: true` discards
  `run.sh`'s exit code; the runtime proof lane is darwin-only. Named, not absorbed; this item's
  static reachability checks run in any lane.
- **The interprocedural/alias residuals** — a helper that skips and a `t.Skip` via an alias are
  declared residuals (see Decision); if taken later they are a separate item, not folded in.
- **Row 49's `GOTOOLCHAIN` zero-needle residual** — a re-added bad arm avoiding the literal
  `GOTOOLCHAIN` string slips the zero-needle; retained unchanged, not re-scoped here.

## Acceptance Criteria

Each AC carries its vacuity self-test and its **observed result on the unmodified tree at `d75b9c3`**
(Verification Log rows cited). The sprint re-runs each on the post-sprint tree. Per base fact A,
**no AC is written against `verify_go.sh rc=0`** — that gate is rc=1 on pristine dev (V19, V20), so
it is an unsatisfiable criterion; every AC is against `go test`/`go vet` on named packages.

- **AC1 — the extended fence exists, RUNS, and passes on the post-sprint tree, in run-existence
  form.** `go test ./host/verifygate/ -run '^TestCanaryDeclaresPositiveArmOnly$' -count=1 -v` →
  rc=0 with exactly 1 top-level `=== RUN` and 1 `--- PASS`; a paired nonsense pattern
  (`-run TestNoSuchCanaryFenceZZZ`) prints `[no tests to run]`, proving the instrument says so
  rather than passing vacuously. **Base @d75b9c3: the verbatim command → `--- PASS` (V7); the
  nonsense form is the vacuity self-test.**
- **AC2 — the canary still passes.** `go test ./host/store/ -run '^TestToolchainCanary$' -count=1 -v`
  → rc=0, 1 `--- PASS`. **Base: `--- PASS` (V8).**
- **AC3 — the recorded `t.Skip` mutant REDs the extended fence, and builds/typechecks.** Sprint
  evidence: apply the exact mutant (insert `t.Skip("MUTANT: reachability probe")` as the first
  statement), run `go vet ./host/store/` (rc=0, the mutant is a valid program) and the scoped fence
  (`go test ./host/verifygate/ -run '^TestCanaryDeclaresPositiveArmOnly$' -count=1 -v`) → **rc=1,
  `--- FAIL`**, message names `t.Skip/t.Skipf/t.SkipNow call count=1, want 0`. Restore sha256-byte-identical
  (`a23cfa79…`), porcelain 0. **Base: the mutant builds — `go vet` rc=0 AND `go test -count=1 -run '^$'` rc=0, both
  re-measured in the complete five-arm table (V34, superseding V18) — and passes the CURRENT fence
  (V9); the prototype REDs it (V12).**
- **AC4 — the `t.SkipNow`, `t.Skipf`, early-`return`, and build-tag mutants each RED the extended fence.**
  Sprint evidence: apply each mutant, `go vet ./host/store/` rc=0 — **except `M-RETURN`, where
  vet's own `unreachable` analyzer reds the mutant (rc=1, `toolchain_canary_test.go:22:2: unreachable
  code`, V32 — this contradicts V18's "rc=0"; vet exits 1 on any analyzer finding) and the compile
  fence is `go test -count=1 -run '^$' ./host/store/` rc=0 instead (V21's form; the package did
  compile and run: canary `--- PASS`, V32)** —, scoped fence → rc=1 with the
  named message (`t.Skip/t.Skipf/t.SkipNow call count=1`, `return statement count=1`, `build constraint
  count=1`). Restore byte-identical between arms. **Base: all five pass the CURRENT fence (V10); the
  prototype REDs them (V13, V14, V15); the `t.Skipf` mutant is the controller-measured hole (V27).**
- **AC4b — the M-BUILDTAG arm is proven non-vacuous ACROSS THE PARSE MODE, not merely red.** Because
  A7's ability to fire is entirely a property of the parse mode (V25), a red on M-BUILDTAG is only
  evidence if it *would not* have been red under the old mode. Sprint evidence: run the M-BUILDTAG
  arm twice against the landed helper — once with `:377` at `parser.ParseComments` (the shipped
  form) and once with it reverted to `0` — and require the two exit codes to **DIFFER**
  (`rc=1` vs `rc=0`), printed side by side, asserted with
  `[ "$rc_comments" -ne "$rc_mode0" ]` rather than eyeballed. If they are equal, A7 is vacuous and
  the milestone is not done, whatever the shipped arm reported. Restore `:377` byte-identical
  afterwards. **Base: not measurable at base — A7 does not exist yet; V25 is the pre-image
  (mode0 cannot distinguish arm 1 from arm 2, `ParseComments` can: 17 vs 18 comments).**
- **AC4c — widening the parse mode does not disturb the retained assertions.** `parser.ParseComments`
  changes what the parser *keeps*, not what it *accepts*, but that is a claim, so measure it: with
  `:377` at `ParseComments`, the four pre-existing P49 shape assertions (A1–A4) must still pass on
  the pristine canary (AC1), and at least one retained-assertion mutant from the parent doc must
  still RED. **Base: A1–A4 green under mode `0` (V7); the widened mode is the new variable.**
- **AC4d — the nested-func-literal policy is proven in BOTH directions, against the landed helper.**
  (i) `M-SKIP-FUNCLIT` (`func() { t.Skip("MUTANT: closure-nested skip on the outer t") }()` as the
  first statement): `go vet ./host/store/` rc=0; the canary itself → `--- SKIP` (the neutering is
  real); scoped fence → **rc=1, `--- FAIL`**, message `t.Skip/t.Skipf/t.SkipNow call count=1, want
  0`. (ii) `M-RETURN-FUNCLIT` (`func() { return }()` as the first statement): `go vet ./host/store/`
  rc=0; the canary itself → `--- PASS` (the assertion still fires); scoped fence → **rc=0,
  `--- PASS`** — a FALSE-POSITIVE control: a red here is a row-55-class defect and the milestone is
  not done. (iii) The `*ast.FuncLit` stop is load-bearing, measured the AC4b way: run
  `M-RETURN-FUNCLIT` once against the shipped A6 and once with A6's `*ast.FuncLit` early
  `return false` removed, and require the two exit codes to DIFFER (`rc=0` vs `rc=1`), asserted with
  `[ "$rc_shipped" -ne "$rc_descend" ]` rather than eyeballed. Restore both files
  sha256-byte-identical, porcelain 0. **Base: not measurable at base — A5/A6 do not exist yet; the
  pre-image is V29–V33 (prototype and in-situ against `d75b9c3`): (i) canary SKIP / fence rc=1,
  (ii) canary PASS / fence rc=0, (iii) rc=0 vs rc=1.**
- **AC5 — hygiene**: `go vet ./host/verifygate/ ./host/store/` rc=0 and `gofmt -l host/verifygate/
  host/store/` prints nothing. **Base: both green (V22, V23).**

Explicitly rejected as an AC: "the full verify gate is green" alone — `verify_go.sh` is rc=1 on
pristine dev (V19, V20), so it is unsatisfiable; and a package-wide `ok` can print while the named
test never ran (the `[no tests to run]` form is rc=0). AC1's run-existence form is the binding
version.

## Non-Vacuity — named RED mutation for every added assertion

Production side mutated (the canary file) — never the test helper. Assertion coverage:
A5-skip-zero-needle←M-SKIP/M-SKIPNOW/M-SKIPF/M-SKIP-FUNCLIT, A6-return-zero-needle←M-RETURN (its
`*ast.FuncLit` stop←M-RETURN-FUNCLIT as a must-stay-green control, load-bearing by V33),
A7-buildtag-zero-needle←M-BUILDTAG.
Each arm is **downstream of its named assertion, not merely suite-wide**: all five mutants pass the
CURRENT fence (V10), so the only thing that can RED them is the new reachability check; and each
arm's failure message names the specific neutering form. **All five are ADDITION-shaped** — they
*add* a neutering form to a pristine canary. This is the shape this file's own history keeps
re-earning: a removal proves the check FIRES, only an addition proves it LOOKS. The existing P49
mutants (M1–M13 in the parent doc) are the removal-shaped arms that prove the shape check FIRES; the
five below are the addition-shaped arms that prove the reachability checks LOOK.

| # | Exact edit (to `host/store/toolchain_canary_test.go`) | Expected RED (named assertion) | Shape |
|---|---|---|---|
| M-SKIP | **THE RECORDED REPRO**: insert `t.Skip("MUTANT: reachability probe")` as the first statement of `TestToolchainCanary`'s body | A5: `t.Skip/t.Skipf/t.SkipNow call count=1, want 0` | **threat-shaped: the filed repro (V9)** — builds/typechecks (V18), passes the CURRENT fence (V9), only the reachability check catches it (V12) |
| M-SKIPNOW | Insert `t.SkipNow()` as the first statement | A5: `t.Skip/t.Skipf/t.SkipNow call count=1, want 0` | sibling form the judge named (V13) |
| M-SKIPF | Insert `t.Skipf("MUTANT: reachability probe")` as the first statement | A5: `t.Skip/t.Skipf/t.SkipNow call count=1, want 0` | **controller-measured hole (V27): the as-drafted A5 missed it; the widened matcher catches it** |
| M-RETURN | Insert `return` as the first statement | A6: `return statement count=1, want 0` | sibling form the judge named (V14) |
| M-BUILDTAG | Prepend `//go:build ignore` (blank line after) to the file | A7: `build constraint count=1, want 0` | sibling form the judge named (V15) |
| M-SKIP-FUNCLIT | Insert `func() { t.Skip("MUTANT: closure-nested skip on the outer t") }()` as the first statement | A5: `t.Skip/t.Skipf/t.SkipNow call count=1, want 0` | **quorum round-1 arm (gpt5-6-sol): A5 descends into `*ast.FuncLit` because the skip genuinely reaches the test (canary `--- SKIP`, V32); a stop-at-FuncLit A5 would MISS it (A5-stop=0, V31)** |
| M-RETURN-FUNCLIT | Insert `func() { return }()` as the first statement | **NO RED — the fence must stay `--- PASS` (rc=0)**, and the canary itself stays `--- PASS` (V32) | **quorum round-1 FALSE-POSITIVE CONTROL (gpt5-6-sol): A6 stops at `*ast.FuncLit` because the return cannot exit the test (V30); under a descending A6 this arm reds with `return statement count=1` (V33) — the row-55 class. Its own non-vacuity is the AC4d(iii) differ-test** |

Green control for all arms: the unmutated post-sprint tree passes AC1/AC2/AC5, and every arm ends
restored sha256-byte-identical with `git status --porcelain` empty — the recipe V9 already ran once
at base (`a23cfa79…` before and after).

**Residual demonstration (not a RED — proves the declared residual is real).** M-ALIAS (`s := t.Skip;
s("MUTANT")`) and M-HELPER (`helperSkip(t)` with a `func helperSkip(t *testing.T) { t.Skip(...) }`)
both leave the extended fence green (prototype V16, V17) — the interprocedural/alias residuals are
real and declared, not silently assumed away.

## Implementation Milestones

**Phase 1 — the three reachability checks** (~0.3d, the whole item): append A5/A6/A7 to
`canaryAssertionShapeProblems`, change its final `return nil` to `return problems`; verify the
package compiles (`go vet ./host/verifygate/`); run AC1/AC2 scoped; run the M-SKIP/M-SKIPNOW/M-SKIPF/
M-RETURN/M-BUILDTAG drill plus the M-SKIP-FUNCLIT/M-RETURN-FUNCLIT pair (AC4d, including the (iii)
differ-test), each restored sha256-byte-identical, porcelain 0 between arms; run AC5
hygiene. One milestone — this is a ~0.3d item.

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| A future legitimate refactor adds a `return` or a `t.Skip` to the canary and REDs the fence (false positive) | Medium | Declared residual; the canary has no legitimate returns or skips today (P5), and a deliberate neutering is exactly what the fence exists to catch. A legitimate shape change must update the fence in the same commit (row 49's contract) |
| The build-tag check misses a constraint form | Low | Covers both modern (`//go:build`) and legacy (`// +build`) forms; the canary has neither today (V4). A future exotic constraint form is a declared residual |
| The AST checks are more code than a textual zero-needle | Low | ~40 LOC, no new imports, no new dependency; the structural form is precise about what is a call and where it lives (Alternatives 1) |
| The full host gate is slow | Low | AC1–AC4 run scoped; no full-suite AC is required (base fact A) |
| A6 descends into a nested func literal and reds on a `return` that cannot exit the test (false positive, the row-55 class) | Medium | Closed by design: A6 stops at `*ast.FuncLit` (V30 language fact; V32/V33 in situ both ways, AC4d(iii)). A5 deliberately does NOT stop, because a closure-nested `t.Skip` really skips (V30, V32) |

## Verification Log

All rows run first-party by the designer at `d75b9c3` (clean tree), shell `zsh`,
`PATH=/opt/homebrew/bin:$PATH`, darwin/arm64, 2026-09-03. KP = known-positive control carried in
the same call. V11–V17 ran the Decision's three checks verbatim as a standalone prototype in
`/tmp/reachproto` (parsing the real canary file and the mutant copies), deleted before the porcelain
checks; the sprint re-runs all arms against the landed test. V29–V33 were run by the round-1
revision pass (2026-09-02, same worktree at `d75b9c3`): V30 in a throwaway module
`/tmp/closureprobe`; V29/V31 via a standalone prototype `/tmp/reachproto2` that walks the body under
BOTH traversal policies side by side; V32/V33 by patching the REAL helper in situ (pristine
`toolchain_pin_gate_test.go` sha256 `f0667a280df795a1…`) and restoring it byte-identically.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Worktree is `d75b9c3`, clean; every mutation arm below ended restored | `git rev-parse HEAD && git status --porcelain \| wc -l` (re-checked after every arm) | `d75b9c3e9cba01508a5f6c2b516f89019e36bdb1`, `0` |
| V2 | Toolchain boundary | `go version` at repo root | `go1.26.6 darwin/arm64` |
| V3 | Fence and helper line numbers | `grep -n "func TestCanaryDeclaresPositiveArmOnly"` and `grep -n "func canaryAssertionShapeProblems"` in `host/verifygate/toolchain_pin_gate_test.go` | fence `:456`, helper `:375` |
| V4 | Canary file carries no build constraint | `head -3 host/store/toolchain_canary_test.go` | `package store`, blank, `import "testing"` — no `//go:build`/`// +build` |
| V5 | `t.Skip` zero-needle baseline, with KP in the same call | `grep -c "t.Skip" host/store/toolchain_canary_test.go`; KP `grep -c "t.Fatalf"` same file | `t.Skip` **0** (grep rc=1); KP `t.Fatalf` **3** (grep rc=0) — the zero is real, the pattern is live |
| V6 | `GOTOOLCHAIN` zero-needle baseline, with KP in the same call | `grep -c "GOTOOLCHAIN"`; KP `grep -c "POSITIVE ARM ONLY"` same file | `GOTOOLCHAIN` **0** (grep rc=1); KP `POSITIVE ARM ONLY` **1** (grep rc=0) |
| V7 | Pristine fence green | `go test -count=1 -v -run TestCanaryDeclaresPositiveArmOnly ./host/verifygate/` | rc=0, 1 `=== RUN`, 1 `--- PASS` |
| V8 | Pristine canary green | `go test -count=1 -v -run TestToolchainCanary ./host/store/` | rc=0, 1 `--- PASS` |
| V9 | **THE REPRO: skipped canary is green everywhere** | `cp` backup; insert `t.Skip("MUTANT: reachability probe")` as first statement; `go vet ./host/store/`; canary; fence; `cp` restore + sha256 + porcelain | sha256 `a23cfa79419ae69136e62981d5bd0c8ea68cdf2e154ba973fb99a3d3c8b47bfd` → `3d7c343ce10c047209d1d9fdb22d4fca7186d4ec222a25eb63068ebe002ff110`; vet rc=0; canary `--- SKIP` rc=0; fence `--- PASS` rc=0; restore `a23cfa79…` byte-identical, porcelain 0 |
| V10 | **The current fence is blind to all four neutering forms** | for each of `m_skip`, `m_skipnow`, `m_return`, `m_buildtag`: place mutant, run `go test -count=1 -run TestCanaryDeclaresPositiveArmOnly ./host/verifygate/`, restore | all four `ok` rc=0 — the current fence cannot see any of them |
| V11 | Prototype: pristine canary has no reachability neutering | `/tmp/reachproto` on the pristine canary | `skip=0 ret=0 buildtag=0` |
| V12 | Prototype: `t.Skip` REDs | same on `m_skip.go` | `skip=1` — A5 fires |
| V13 | Prototype: `t.SkipNow` REDs | same on `m_skipnow.go` | `skip=1` — A5 fires |
| V14 | Prototype: early `return` REDs | same on `m_return.go` | `ret=1` — A6 fires |
| V15 | Prototype: `//go:build` constraint REDs | same on `m_buildtag.go` | `buildtag=1` — A7 fires (the `//go:build ignore` comment is in `f.Comments`) |
| V16 | Prototype: `t.Skip` via alias NOT caught (residual) | same on `m_alias.go` (`s := t.Skip; s("MUTANT")`) | `skip=0` — declared residual, real |
| V17 | Prototype: helper that skips NOT caught (residual) | same on `m_helper.go` (`helperSkip(t)`) | `skip=0` — declared residual, real |
| V18 | ~~Every mutant is a real, buildable program~~ **SUPERSEDED BY V34 on its `m_return` arm** | for each of `m_skip`, `m_skipnow`, `m_return`, `m_buildtag`: place mutant, `go vet ./host/store/`, restore | ~~vet rc=0 for all four~~ — **the `m_return` cell is WRONG**. Replaced by the V34-confirmed result: `m_return`: vet **rc=1** (`unreachable code`), package compiles and runs via `go test -count=1 -run '^$' ./host/store/` rc=0 (V21 form). The other three arms stand (V34 re-measured them: rc=0). Row retained unedited-in-substance for provenance; cite **V34**, never this row |
| V19 | **Base fact A: `verify_go.sh` is rc=1 on pristine dev** | `./scripts/verify_go.sh` | rc=1, `✗ AILANG_BIN is unset — host/replay tests would t.Skip() silently and this gate would be false-green.` — the anti-false-green guard fires before any go test runs |
| V20 | **Base fact A: the driver-drift gate also reds** | `./scripts/verify_go.sh --driver-fleet-check` | rc=1, `FATAL: DRIVER DRIFT vs FLEET (D-WORLD-DRIVER-1)` — `derive-planner-lane.sh`, `mission-control.sh`, `test_mission_routing.sh` differ from fleet HEAD `c8c841e24e90ee6231884867fd2f09490140c7cc` (the premise's `193e3abae` is stale; the drift is real) |
| V21 | **Base fact B: `go build` is not a compile fence for a `_test.go` file** | throwaway module in `/tmp/buildfence` with a hard type error in `x_test.go`; `go build ./...`; `go vet ./...`; `go test -run '^$' ./...` | build rc=0 (the type error is invisible to build); vet rc=1; test rc=1 — every compile fence in this doc is `go vet`/`go test -run '^$'`, never `go build` |
| V22 | Imports needed are already present | `sed -n '1,14p' host/verifygate/toolchain_pin_gate_test.go` | `fmt`, `go/ast`, `go/parser`, `go/token`, `strings` all imported — no new imports |
| V23 | Hygiene at base | `go vet ./host/verifygate/ ./host/store/`; `gofmt -l host/verifygate/ host/store/ \| wc -l` | vet rc=0; gofmt `0` files |
| V25 | **CONTROLLER-MEASURED, AND IT REFUTES THE DESIGN'S OWN PREMISE: the helper parses with mode `0`, so `f.Comments` is empty and check A7 as drafted is VACUOUS.** Two arms, one known-positive control. | `grep -n "parser.ParseFile" host/verifygate/toolchain_pin_gate_test.go`; then a standalone program parsing the canary under mode `0` and under `parser.ParseComments`, on (arm 1) the pristine canary and (arm 2) the same file prefixed `//go:build ignore` | `:377 parser.ParseFile(fset, "", src, 0)`. Arm 1: mode0 **0** comments, ParseComments **17**. Arm 2: mode0 **0**, ParseComments **18**. So mode `0` cannot distinguish the A7 mutant from pristine; `ParseComments` can. Control that the mode is not exotic here: `grep -rn ParseComments --include='*.go'` → `cmd/world-publish/wiring_test.go:60`, `host/broker/invoke_boundary_test.go:112`, `:317`. **Consequence folded into the Decision: the sprint must also change `:377` to `parser.ParseComments`, and must run the A7 arm in BOTH modes requiring the verdicts to DIFFER.** |
| V26 | Controller's own fleet-HEAD sha in the brief was stale; V20's correction stands | `verify_go.sh` driver-drift arm re-run by the controller with `AILANG_BIN` exported | the drift is real and now **1092** diff-lines across the three driver files (control: an unchanged driver file diffs at **0**); the fleet HEAD literal moved, the drift did not |
| V24 | Parent doc declares this exact class | `grep -n "What the gate CANNOT catch\|The shape is not the behaviour" design_docs/implemented/w-canary-fence-passes-a-gutted-canary.md` | residual section at `:314`, the sentence at `:316` — this item closes a disclosed residual, not a defect P49 hid |
| V27 | **The A5 selector set is CLOSED by construction against `go doc testing.T`'s three skip methods, and the as-drafted A5 had a real hole the widened matcher closes.** | `go doc testing.T`; then the controller's four-arm probe (PRISTINE / M-SKIP / M-SKIPNOW / M-SKIPF), each arm placing the mutant, running the A5 matcher, restoring byte-identical (sha256 back to `a23cfa79419ae691`, porcelain clean) | `go doc testing.T` enumerates exactly **`Skip(args ...any)`, `Skipf(format string, args ...any)`, `SkipNow()`** — three methods, no more. Probe: PRISTINE sha=`a23cfa79419a` A5-as-drafted=0 A5-widened=0; M-SKIP sha=`9203903fe73c` A5-as-drafted=1 A5-widened=1; M-SKIPNOW sha=`8311cfeff59b` A5-as-drafted=1 A5-widened=1; **M-SKIPF sha=`c695c66416e7` A5-as-drafted=**0** A5-widened=1** — the as-drafted matcher missed `t.Skipf`; the widened matcher catches it and does not over-match on pristine. |
| V28 | **The helper's real source: it has NO `problems` accumulator — it returns a FRESH `[]string{...}` at each of three early failure points and `return nil` at the end, so the sprint MUST declare `var problems []string` before the appended checks.** | `sed -n '375,382p' host/verifygate/toolchain_pin_gate_test.go` | `func canaryAssertionShapeProblems(src string) []string {` / `fset := token.NewFileSet()` / `f, err := parser.ParseFile(fset, "", src, 0)` / `if err != nil { return []string{fmt.Sprintf("parse error: %v", err)} }` / `var funcs []*ast.FuncDecl` — confirmed: no `problems` variable exists; the three EXISTING early returns are left byte-unchanged (they still short-circuit), and the new checks run only after the shape checks pass. |
| V29 | **Pristine canary has NO func literals, and its only `return` is OUTSIDE the body** — both traversal policies are zero-needles today; the body scoping is real | `grep -c 'func(' host/store/toolchain_canary_test.go`; `grep -n return` same file; `/tmp/reachproto2` census on the pristine file | `0`; `59:func (s canaryString) String() string { return string(s) }` (the only hit, outside `TestToolchainCanary`); census `file: FuncLit=0 ReturnStmt=1 \| body: FuncLit=0 A6-descend=0 A6-stopAtFuncLit=0 A5-descend=0 A5-stopAtFuncLit=0` — file-level `ReturnStmt` 1, body-scoped 0 under either policy |
| V30 | **Runtime fact behind the policy split: a closure-nested `return` does NOT exit the enclosing test; a closure-nested `t.Skip` on the outer `t` DOES** | throwaway module `/tmp/closureprobe` (`go 1.26`): `TestReturnInFuncLit` = `func() { return }()` then `t.Fatalf("REACHED…")`; `TestSkipInFuncLit` = `func() { t.Skip("closure-nested skip on the outer t") }()` then `t.Fatalf("REACHED…")`; `go vet ./...`; `go test -count=1 -v ./...` | vet rc=0; test rc=1: `=== RUN TestReturnInFuncLit` / `x_test.go:9: REACHED: the assertion after a closure-nested return still fires` / `--- FAIL: TestReturnInFuncLit`; `=== RUN TestSkipInFuncLit` / `x_test.go:15: closure-nested skip on the outer t` / `--- SKIP: TestSkipInFuncLit` (its Fatalf line never printed) |
| V31 | Prototype census of the three revision mutants under BOTH policies — the policies split exactly on the closure arms | `/tmp/reachproto2` on `m_retlit.go` (`func() { return }()` first stmt), `m_skiplit.go` (`func() { t.Skip(…) }()` first stmt), `m_return.go` (bare `return` first stmt) | `m_retlit`: `body: FuncLit=1 A6-descend=1 A6-stopAtFuncLit=0 A5-descend=0 A5-stopAtFuncLit=0`; `m_skiplit`: `body: FuncLit=1 A6-descend=0 A6-stopAtFuncLit=0 A5-descend=1 A5-stopAtFuncLit=0`; `m_return`: `body: FuncLit=0 A6-descend=1 A6-stopAtFuncLit=1 A5-descend=0 A5-stopAtFuncLit=0` — descend-A6 false-reds `m_retlit`; stop-A5 misses `m_skiplit`; both policies agree on the bare `return` |
| V32 | **IN SITU: the Decision's exact code block (A5 descends, A6 stops at `*ast.FuncLit`, A7, `var problems`, `:377` → `ParseComments`) compiles against the REAL helper and behaves as specified on pristine + three mutants** | patch `toolchain_pin_gate_test.go` (sha → `e71c14a64d878eb4`); `go vet ./host/verifygate/`; `gofmt -l host/verifygate/`; per arm: place mutant, `go vet ./host/store/`, canary `go test -count=1 -v -run '^TestToolchainCanary$' ./host/store/`, fence `go test -count=1 -v -run '^TestCanaryDeclaresPositiveArmOnly$' ./host/verifygate/`, restore canary + sha | vet rc=0; gofmt empty. PRISTINE (`a23cfa79419ae691`): fence rc=0 `--- PASS`. `m_retlit` (`7cbabbb1378e2d7e`): vet rc=0; canary `--- PASS`; fence rc=0 `--- PASS` (correct: not a threat). `m_skiplit` (`13c85ea2e880f350`): vet rc=0; canary `--- SKIP`; fence rc=1 `toolchain_pin_gate_test.go:520: …/host/store/toolchain_canary_test.go: t.Skip/t.Skipf/t.SkipNow call count=1, want 0 (a skipped canary asserts nothing)` / `--- FAIL`. `m_return` (`2ff015dccd545c3c`): **vet rc=1 `host/store/toolchain_canary_test.go:22:2: unreachable code`** (contradicts V18's "rc=0"; the package still compiled and ran — canary `--- PASS`, the assertion neutered); fence rc=1 `… return statement count=1, want 0 (an early return neuters the assertion)` / `--- FAIL`. Canary restored `a23cfa79419ae691` after every arm |
| V33 | **The `*ast.FuncLit` stop in A6 is LOAD-BEARING: the as-drafted descending A6 false-reds `m_retlit`; verdicts DIFFER (rc=0 vs rc=1)**; then full restore | same patch with only A6's `if _, ok := n.(*ast.FuncLit); ok { return false }` removed (sha → `a4979d509cc51daa`); `go vet ./host/verifygate/`; pristine fence; `m_retlit` fence; restore both files; `git status --porcelain` | vet rc=0; pristine fence rc=0 `--- PASS`; `m_retlit` fence **rc=1** `toolchain_pin_gate_test.go:517: …: return statement count=1, want 0 (an early return neuters the assertion)` / `--- FAIL` — a red over a canary that still `--- PASS`es (V32) — vs V32's rc=0 for the same mutant. Restore: `toolchain_pin_gate_test.go` sha `f0667a280df795a1` (= pristine), canary `a23cfa79419ae691`; porcelain = `?? design_docs/planned/w-canary-fence-blind-to-a-skipped-canary.md` only (1 line) |
| V34 | **THE CONTRADICTION RESOLVED BY A THIRD, ISOLATED MEASUREMENT (round-2 carve-out, all three reviewers' objection): V18 is WRONG on `m_return`, V32 is RIGHT — and the complete five-arm `vet`/compile-fence table is now recorded rather than inferred.** Controller-run, first-party, at `d75b9c3`, `go1.26.6 darwin/arm64`, PRISTINE canary (no verifygate patch), so the reading is about the MUTANT and not about any prototype. | pristine control `go vet ./host/store/`; then per arm: `cp` restore, insert the mutant as the first statement of `TestToolchainCanary` (build-tag arm: prepend `//go:build ignore_canary` above `package store`), `shasum -a 256`, `go vet ./host/store/`, `go test -count=1 -run '^$' ./host/store/`, `cp` restore | PRISTINE `a23cfa79419ae691` vet **rc=0**. `m_skip` `9203903fe73c4e24` vet **rc=0**, fence-compile **rc=0**. `m_skipf` `8842b6bacc6243b2` vet **rc=0**, **rc=0**. `m_skipnow` `8311cfeff59be3ae` vet **rc=0**, **rc=0**. **`m_return` `2ff015dccd545c3c` vet **rc=1**, stderr `host/store/toolchain_canary_test.go:22:2: unreachable code`, compile-fence **rc=0`** — so the package COMPILES and the vet red is the unreachable-code ANALYZER, exactly as V32 recorded and contrary to V18. `m_buildtag` `43596b333f39ee1f` vet **rc=0**, compile-fence **rc=0** (`ok … [no tests to run]`, i.e. the tag excluded the file, which is the mutant's whole point). Restored `a23cfa79419ae691` byte-identical; `git status --porcelain` shows only the untracked design doc |


## Quorum verification log

**Round 1 — `iso_ts` 2026-09-02T11:35:17Z — verdict `blocked` — FULL STRENGTH.** Artifact:
`.ailang/state/mission-quorum/w-canary-fence-blind-to-a-skipped-canary-2026-09-02T11-35-17Z.json`
(`.synthesis.verdict` = `blocked`; `.synthesis.absent_reviewers` = `[]`; all three reviewers
`present: true`; `.synthesis.total_cost_usd` = `0.11536481999999999`, which equals the sum of the
three per-reviewer costs). The controller's in-session verdict was `pass` ("Direction endorsed; I am
not the independent eye"). All three reviewers rejected; none challenged the design direction (three
statically-visible reachability zero-needles folded into the same AST pass). Objections, where each
is now answered, and who applied the fix:

| Reviewer | Verdict | Cost (USD) | Objection (one line) | Where answered in this doc | Applied by |
|---|---|---|---|---|---|
| gpt5-6-sol | reject | 0.063645 | (a) the doc claimed the `t.Skip` mutation leaves the AST "byte-identical" — false: the mutation adds an `ExprStmt`/`CallExpr` and the two logged sha256s differ. (b) `ast.Inspect(funcs[0].Body, …)` descends into nested function literals, so A5/A6 may reject skips/returns that cannot exit or skip `TestToolchainCanary` — document that conservative policy explicitly or stop traversal at `*ast.FuncLit` | (a) Thesis and P2 now state the narrower property — the three counted shapes are unchanged, the sha256s differ (V9); the fence's own three `== 1` assertions passing on the mutant IS the counter evidence, no separate AST-counter row was added. (b) Decision → *Nested function literals* paragraph; the A5/A6 code block; Design Freeze structural bullet; AC4 (vet-rc correction) and AC4d; Non-Vacuity M-SKIP-FUNCLIT / M-RETURN-FUNCLIT; three new declared-residual bullets; Risks; V29–V33 | (a) controller, prior pass. (b) **this revision pass** — before it the doc had zero occurrences of "FuncLit", "func literal" or "nested function" |
| gemini-3-1-pro | reject | 0.025634 | A5 matched `Skip`/`SkipNow` but omitted `t.Skipf`, the same silent-green bypass | A5 selector matches `Skip \|\| Skipf \|\| SkipNow`; the problem string names all three; M-SKIPF arm in Non-Vacuity; V27 four-arm probe (as-drafted=0, widened=1 on M-SKIPF) | controller, prior pass |
| oc-glm-5-2 | reject | 0.02608582 | the helper body is the load-bearing premise and was asserted, not shown — variable names, return pattern, `funcs` population and parse mode all unverified | V28 (real source: no `problems` accumulator, fresh `[]string{…}` per early failure, `funcs` from `f.Decls`, `return nil` at the end); V25 (parse mode `0`); Design Freeze and Files to Modify state the `var problems []string` requirement; V32 now also proves the block compiles and runs against the real helper in situ | controller, prior pass (V28); the in-situ compile proof (V32) by this revision pass |

**Round 2 — `iso_ts` 2026-09-02T14:46:58Z — verdict `blocked` — FULL STRENGTH.** Artifact:
`.ailang/state/mission-quorum/w-canary-fence-blind-to-a-skipped-canary-2026-09-02T14-46-58Z.json`
(`.synthesis.absent_reviewers` = `[]`, cross-checked against
`[.reviewers[]|select(.present==false)]` = `[]`; control `has("synthesis")` true; total `$0.1568`).
All three reviewers rejected, and **all three raised the SAME single objection**: the doc carried a
contradiction between **P7/V18** (`go vet ./host/store/` rc=0 for the `M-RETURN` mutant) and
**V32/AC4** (rc=1, `unreachable code`), acknowledged but unresolved — so one of two same-setup
measurements was stale or fabricated and the doc did not say which.

| Reviewer | Verdict | Cost (USD) | Objection | Disposition |
|---|---|---|---|---|
| gpt5-6-sol | reject | 0.0930 | P7 is false by the doc's own later evidence; P7 also attributed `t.Skipf` coverage to V18, which never tested that arm | UPHELD and applied |
| gemini-3-1-pro | reject | 0.0387 | P7 and V18 retain a knowingly-false claim the author had already discovered in V32 | UPHELD and applied |
| oc-glm-5-2 | reject | 0.0252 | V18 and V32 directly contradict each other on the identical command, mutant, worktree and Go version; resolve it by a third measurement rather than asserting it away | UPHELD and applied |
| controller (in-session) | pass | — | — | overruled by the three above |

**Disposition: the NARROW-REFINEMENT CARVE-OUT** (ratified for this mission at iteration 13),
conditions checked BEFORE use and both met: (a) every remaining blocking objection carries a
concrete, reviewer-authored `proposed_fix`, and (b) none disputes the design DIRECTION — all three
dispute one factual premise cell. This is also a PREMISE objection, so per the shared skill's rule
3f the controller **measured it first-party rather than forwarding it**, which is precisely what
`oc-glm-5-2`'s `proposed_fix` asked for: **V34** re-ran all five mutants in isolation on a pristine
tree. **Verdict: V18 is WRONG on its `m_return` arm; V32 is RIGHT.** `m_return` is `vet rc=1`
(`unreachable code`) with `go test -count=1 -run '^$'` **rc=0** — a real, compilable program whose
vet red is an analyzer finding. The other four arms are `vet rc=0`, so V18 was wrong in exactly one
cell.

Fixes applied VERBATIM from the reviewers' own `proposed_fix` text (no controller-invented
resolution, no objection overridden): **P7** rewritten to `oc-glm-5-2`'s wording, extended with the
`t.Skipf` arm `gpt5-6-sol` correctly noted V18 never covered; **V18** marked SUPERSEDED with its
`m_return` cell replaced by the V34-confirmed result and retained for provenance; **AC3**'s
base-evidence citation moved off V18 onto V34; **V34** added as the third measurement. AC4 already
carried the `go test -run '^$'` compile fence for this arm from the revision round.

No round 3 was run: the carve-out routes straight to sprint-planner, which is what it exists for.

## Related Documents

- [`../implemented/w-canary-fence-passes-a-gutted-canary.md`](../implemented/w-canary-fence-passes-a-gutted-canary.md)
  — row 49, the parent: built `canaryAssertionShapeProblems` and its three shape assertions; its
  residual *"The shape is not the behaviour"* (`:314-316`) is the exact gap this item closes.
- [`../implemented/w-canary-control-does-not-survive-a-floor-raise.md`](../implemented/w-canary-control-does-not-survive-a-floor-raise.md)
  — row 42, the grandparent: built `TestCanaryDeclaresPositiveArmOnly` and the canary.
- [`../implemented/w-race-gate-blindspot.md`](../implemented/w-race-gate-blindspot.md) — where the
  canary and the nested-module pattern were born.
- [`../implemented/w-miscompile-instrument-inert-in-ci.md`](../implemented/w-miscompile-instrument-inert-in-ci.md)
  — row 44, the runtime lane caveat this item declares rather than fixes.
