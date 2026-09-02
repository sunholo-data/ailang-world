# SPRINT PLAN — w-canary-fence-blind-to-a-skipped-canary (queue row 56)

**Design doc**: `design_docs/planned/w-canary-fence-blind-to-a-skipped-canary.md` (quorum-cleared,
round-2 narrow-refinement carve-out)
**Worktree**: `/tmp/wt-row56`, branch `sprint/w-canary-fence-blind-to-a-skipped-canary`, based on
`origin/dev` at `d75b9c3e9cba01508a5f6c2b516f89019e36bdb1`
**Planner lane**: `opus fail-closed:env-pin`
**Toolchain boundary**: every command in this plan was run FIRST-PARTY by the planner in this
worktree at `d75b9c3`, shell `zsh`, `PATH=/opt/homebrew/bin:$PATH`, `go version` =
`go1.26.6 darwin/arm64`. No `.ail` source is written or changed, so no `ai-check` / `verify_ail.sh`
gate applies to this sprint.
**Size**: ~0.3d. ONE file changed: `host/verifygate/toolchain_pin_gate_test.go`. No production code,
no `host/store` edit (read-only by path), no `run.sh`, no `ci.yml`, no `scripts/`, no new package,
no new import, no new Go dependency.

---

## 0. Base facts the sprint must not re-litigate

### 0.1 `go build` IS NOT A COMPILE FENCE FOR A `_test.go`

Every assertion in this sprint lives in a `_test.go`. `go build ./...` returns rc=0 with a hard type
error inside a `_test.go` (measured by this mission 2026-09-01; memory
`go-build-is-not-a-compile-fence-for-test-files`). **No milestone gate, acceptance criterion or
mutation arm in this plan uses `go build`.** The compile fence is:

- `go test -count=1 -run '^$' <pkg>` — the binding fence, valid on EVERY arm; or
- `go vet <pkg>` — valid on every arm EXCEPT `M-RETURN` (see 0.2).

### 0.2 `go vet ./host/store/` is rc=1 on the `M-RETURN` arm BY DESIGN

Re-measured first-party this session (see 4.1, arm M-RETURN): `go vet ./host/store/` → **rc=1**,
`host/store/toolchain_canary_test.go:22:2: unreachable code`. That is the `unreachable` ANALYZER
firing, not a compile failure — the same tree gives `go test -count=1 -run '^$' ./host/store/`
**rc=0** and the canary itself `--- PASS`. **`go vet ./host/store/` is therefore NOT an acceptance
gate on the M-RETURN arm.** Use the `go test -run '^$'` compile fence there. (This is the doc's V34
superseding V18; independently reconfirmed here.)

### 0.3 `./scripts/verify_go.sh` is NOT a sprint gate

It is RED on this rig for two reasons that are not World failures:
(a) the FLEET-owned driver copy is ~759 diff-lines behind fleet HEAD and the drift gate reds on that
by design (`D-WORLD-DRIVER-1`: that red means "the fleet must commit", never "absorb it");
(b) row 58's known flake `TestHandlerTimeoutKillsTheWholeProcessGroup`
(`exec_started=false forked=false`).
**All sprint gates are scoped to the two packages this sprint touches: `./host/verifygate/` and
`./host/store/`.** A gate that is already red on untouched `dev` is a broken gate, not a finding.

---

## 1. The change — four edits, one file

`host/verifygate/toolchain_pin_gate_test.go`. Line numbers re-derived first-party at `d75b9c3`
(`grep -n`), NOT transcribed from the design doc:

| # | Location @`d75b9c3` | Edit |
|---|---|---|
| E1 | `:377` `	f, err := parser.ParseFile(fset, "", src, 0)` | → `parser.ParseFile(fset, "", src, parser.ParseComments)`. **MANDATORY**: under mode `0` the parser discards comments, `f.Comments` is empty, and A7 can never fire (proven by the AC4b differ-test, 4.2). |
| E2 | before `:421` | declare `var problems []string` (the helper has NO accumulator today — confirmed by reading the source: it returns a fresh `[]string{…}` at each of three early failure points) |
| E3 | before `:421` | append the A5 / A6 / A7 blocks |
| E4 | `:421` `	return nil` | → `	return problems` |

**Byte-unchanged and NOT to be touched**: the three existing early returns
(func-decl-count `!= 1` @`:392`-ish, assertion-if-count `!= 1`, direct-`t.Fatalf`-count `!= 1`), the
helpers `isRowsField` / `isStateRootLit` / `isTFatalfCall`, and in
`TestCanaryDeclaresPositiveArmOnly` (@`:456`) the `GOTOOLCHAIN` zero-needle and the
`POSITIVE ARM ONLY` marker. Because the three early returns stay, the new checks run only when the
shape checks passed — that is the design's chosen semantics.

**Imports**: `fmt`, `go/ast`, `go/parser`, `go/token`, `strings` are ALL already imported
(planner-verified, `sed -n '1,15p'`). No new import.

**Traversal policies are DELIBERATELY DIFFERENT and both are load-bearing:**
- **A5 (skip calls) DESCENDS into `*ast.FuncLit`** — `t.Skip` on the outer `t` inside a closure
  genuinely `runtime.Goexit`s the test. Proven RED by the `M-SKIP-FUNCLIT` arm.
- **A6 (returns) STOPS at `*ast.FuncLit`** (`return false` from the Inspect callback) — a `return`
  inside a closure exits the literal, not the test. Counting it would be a FALSE RED of exactly the
  class this mission landed yesterday (row 55). Proven by the `M-RETURN-FUNCLIT` MUST-STAY-GREEN
  control plus the AC4d(iii) differ-test.

The planner applied all four edits in situ, ran every arm below, and restored the file
sha256-byte-identically (`f0667a280df795a1…` → `f0667a280df795a1…`). The patch is known to compile,
`gofmt`-clean, and behave exactly as tabulated. **Measured size: +53 lines** (706 → 759).

---

## 2. Milestones

Five checkpoints. **All five collapse into ONE commit** — see 2.6.

### M1 — widen the parse mode (E1 only)

Change `:377` to `parser.ParseComments`. Nothing else.

**Gate (planner-baselined, all rc=0 with E1 applied and no other edit):**
```
gofmt -l host/verifygate/                                       # must print nothing
go vet ./host/verifygate/                                       # rc=0
go test -count=1 -run '^$' ./host/verifygate/                   # rc=0  (COMPILE FENCE, not go build)
go test ./host/verifygate/ -run '^TestCanaryDeclaresPositiveArmOnly$' -count=1 -v   # rc=0, 1 --- PASS
```
The last line IS **AC4c**: widening the mode changes what the parser *keeps*, not what it *accepts*,
so the four retained P49 shape assertions must still pass. Planner-measured under the full patch:
`--- PASS`, rc=0.

### M2 — A5, the skip zero-needle (E2 + E4 + the A5 block)

Adds `var problems []string`, the A5 `ast.Inspect` that DESCENDS into func literals matching
`t.Skip` / `t.Skipf` / `t.SkipNow` selectors on receiver ident `t`, and flips `return nil` →
`return problems`. Selector set is CLOSED by construction: `go doc testing.T` exposes exactly those
three skip methods.

**Gate:** M1's four commands, all still green, PLUS the four A5 red arms of 4.1
(`M-SKIP`, `M-SKIPF`, `M-SKIPNOW`, `M-SKIP-FUNCLIT`) each producing fence rc=1 with the substring
`t.Skip/t.Skipf/t.SkipNow call count=1, want 0`.

### M3 — A6, the return zero-needle with the `*ast.FuncLit` stop

**Gate:** M1's four commands green, PLUS:
- red arm `M-RETURN` → fence rc=1, substring `return statement count=1, want 0`;
- **must-stay-green** arm `M-RETURN-FUNCLIT` → fence **rc=0, `--- PASS`** (a red here is a
  row-55-class false positive and M3 is NOT done);
- **AC4d(iii) differ-test**: run `M-RETURN-FUNCLIT` once against the shipped A6 and once with A6's
  `if _, ok := n.(*ast.FuncLit); ok { return false }` deleted; assert with
  `[ "$rc_shipped" -ne "$rc_descend" ]`, not by eye. Planner-measured: `rc_shipped=0`,
  `rc_descend=1` → DIFFER, PASS.

### M4 — A7, the build-constraint zero-needle

Walks `f.Comments` for `//go:build` / `// +build` prefixes. Depends on M1 (E1).

**Gate:** M1's four commands green, PLUS:
- red arm `M-BUILDTAG` → fence rc=1, substring `build constraint count=1, want 0`;
- **AC4b differ-test**: run `M-BUILDTAG` twice against the landed helper — once with `:377` at
  `parser.ParseComments`, once reverted to `0` — and assert
  `[ "$rc_comments" -ne "$rc_mode0" ]`. Planner-measured: `rc_comments=1`, `rc_mode0=0` → DIFFER,
  PASS. **If these are equal, A7 is vacuous and M4 is not done, whatever the shipped arm reported.**
  Restore `:377` byte-identically afterwards.

### M5 — hygiene and full acceptance sweep

Run the whole table in §3 on the post-sprint tree, then `git status --porcelain` must show only the
two design docs and the one modified `.go` file.

### 2.6 Commit collapsing — ONE commit

M1–M5 are **verification checkpoints inside a single commit**, not separate commits, because the
intermediate states are not independently shippable:

- **M1 alone** is a widening with no assertion behind it — dead surface, and a reviewer would
  correctly ask why the mode changed. It is only justified by M4.
- **M4 alone without M1** is a VACUOUS assertion wearing a gate's clothes: `f.Comments` is empty
  under mode `0`, so A7 can never fire (AC4b measures exactly this).
- M2 and M3 are each individually compilable and green, but splitting a ~53-line addition to one
  helper across three commits buys nothing and breaks the "a shape change updates the fence in the
  same commit" contract row 49 established.

So: **one commit touching one file**, gated at five checkpoints on the way there.

---

## 3. Acceptance table — every command BASELINED ON THE PRISTINE WORKTREE

Every rc below was produced by the planner running the verbatim command at `d75b9c3` with a clean
tree. A gate already red at base would be a broken gate; there are none.

| AC | Verbatim command | **rc AT BASE (pristine `d75b9c3`)** | Observed at base | Required post-sprint |
|---|---|---|---|---|
| AC1 | `go test ./host/verifygate/ -run '^TestCanaryDeclaresPositiveArmOnly$' -count=1 -v` | **rc=0** | exactly 1 `=== RUN`, 1 `--- PASS`, `ok …/host/verifygate 0.313s` | rc=0, exactly 1 `=== RUN` + 1 `--- PASS` |
| AC1-vac | `go test ./host/verifygate/ -run 'TestNoSuchCanaryFenceZZZ' -count=1 -v` | **rc=0** | `testing: warning: no tests to run` / `ok … [no tests to run]` | unchanged — this is the vacuity self-test: it proves a package-wide `ok` can print while the named test never ran, so AC1's run-existence form (1 `=== RUN`) is the binding one |
| AC2 | `go test ./host/store/ -run '^TestToolchainCanary$' -count=1 -v` | **rc=0** | 1 `--- PASS`, `ok …/host/store 0.196s` | rc=0, 1 `--- PASS` (the canary is read-only; this must be unchanged) |
| AC4c | AC1, run with `:377` at `parser.ParseComments` | n/a at base (mode is `0`) | — | rc=0, `--- PASS` — planner-measured under the full patch: **PASS** |
| AC5a | `go vet ./host/verifygate/ ./host/store/` | **rc=0** | no output | rc=0 |
| AC5b | `gofmt -l host/verifygate/ host/store/` | **rc=0** | **0 lines** | 0 lines |
| CF-v | `go test -count=1 -run '^$' ./host/verifygate/` | **rc=0** | `ok … [no tests to run]` | rc=0 — THE compile fence for the edited package |
| CF-s | `go test -count=1 -run '^$' ./host/store/` | **rc=0** | `ok … [no tests to run]` | rc=0 — THE compile fence for every canary mutant |
| — | `./scripts/verify_go.sh` | **NOT A GATE** | rc=1 on pristine dev for two non-World reasons (§0.3) | not run as a sprint gate |

Note for the executor: the FIRST `go test`/`go vet` against `./host/store/` on a cold cache took
>120s in this worktree (planner measured 18s warm, timeout at 120s cold). **Warm the cache once with
`go test -count=1 -run '^$' ./host/store/` before the mutation drill** or arms will look like
failures when they are timeouts.

---

## 4. Non-vacuity — the mutation table, ALL ARMS PLANNER-MEASURED IN SITU

Production side mutated is always the **canary** (`host/store/toolchain_canary_test.go`), never the
test helper. All arms are ADDITION-shaped: they add a neutering form to a pristine canary. A removal
proves a check FIRES; only an addition proves it LOOKS.

**Downstream-not-suite-wide proof — planner-measured at base, all SEVEN arms:** with the UNPATCHED
helper, every one of the seven mutants leaves the fence **rc=0**. So the only thing that can RED any
arm is the new reachability check, not some unrelated pre-existing assertion.

```
M-SKIP CURRENT-fence rc=0        M-RETURN CURRENT-fence rc=0        M-SKIP-FUNCLIT CURRENT-fence rc=0
M-SKIPF CURRENT-fence rc=0       M-BUILDTAG CURRENT-fence rc=0      M-RETURN-FUNCLIT CURRENT-fence rc=0
M-SKIPNOW CURRENT-fence rc=0
```

### 4.1 Arms — RED arms and MUST-STAY-GREEN arms are explicitly distinguished

All rows below were run by the planner with the four edits applied in situ. `vet` = `go vet
./host/store/`; `CF` = `go test -count=1 -run '^$' ./host/store/`; `canary` = `go test -count=1 -v
-run '^TestToolchainCanary$' ./host/store/`; `fence` = `go test -count=1 -v
-run '^TestCanaryDeclaresPositiveArmOnly$' ./host/verifygate/`.

| Arm | Class | Exact mutant text (pin it verbatim — see R1) | canary sha256(16) | vet | **CF** | canary result | fence rc | **EXACT expected `--- FAIL` set** |
|---|---|---|---|---|---|---|---|---|
| PRISTINE | control | — | `a23cfa79419ae691` | rc=0 | rc=0 | `--- PASS` | **0** | none; `--- PASS: TestCanaryDeclaresPositiveArmOnly` |
| M-SKIP | **RED (A5)** | insert `t.Skip("MUTANT: reachability probe")` as first stmt of the body | `3d7c343ce10c0472` | rc=0 | rc=0 | `--- SKIP` | **1** | exactly `--- FAIL: TestCanaryDeclaresPositiveArmOnly`; substrings `t.Skip/t.Skipf/t.SkipNow call count=1, want 0 (a skipped canary asserts nothing)` AND `instrument failure: … no longer asserts the miscompile shape` |
| M-SKIPF | **RED (A5)** | insert `t.Skipf("MUTANT: reachability probe %s", "x")` | `9c8778abdd21c8ae` | rc=0 | rc=0 | `--- SKIP` | **1** | same as M-SKIP |
| M-SKIPNOW | **RED (A5)** | insert `t.SkipNow()` | `8311cfeff59be3ae` | rc=0 | rc=0 | `--- SKIP` | **1** | same as M-SKIP |
| M-RETURN | **RED (A6)** | insert `return` | `2ff015dccd545c3c` | **rc=1** `…:22:2: unreachable code` — ANALYZER, not a compile error; **DO NOT GATE ON VET HERE** | **rc=0** | `--- PASS` (the assertion is neutered but the test still reports pass — this is the whole point) | **1** | exactly `--- FAIL: TestCanaryDeclaresPositiveArmOnly`; substring `return statement count=1, want 0 (an early return neuters the assertion)` |
| M-BUILDTAG | **RED (A7)** | prepend `//go:build ignore_canary` + a blank line above `package store` | `43596b333f39ee1f` | rc=0 | rc=0 | `ok … [no tests to run]` (the tag excluded the file — the mutant's whole point) | **1** | exactly `--- FAIL: TestCanaryDeclaresPositiveArmOnly`; substring `build constraint count=1, want 0 (a build tag can exclude the canary from the build)` |
| M-SKIP-FUNCLIT | **RED (A5 descends)** | insert `func() { t.Skip("MUTANT: closure-nested skip on the outer t") }()` | `13c85ea2e880f350` | rc=0 | rc=0 | `--- SKIP` (the skip really reaches the test) | **1** | same as M-SKIP. A stop-at-FuncLit A5 would read 0 and MISS this |
| M-RETURN-FUNCLIT | **MUST STAY GREEN** (A6's FuncLit stop; FALSE-POSITIVE control) | insert `func() { return }()` | `7cbabbb1378e2d7e` | rc=0 | rc=0 | `--- PASS` (the assertion still fires) | **0** | **NO `--- FAIL` AT ALL.** `--- PASS: TestCanaryDeclaresPositiveArmOnly`. A red here is a row-55-class defect and the milestone is NOT done |

**A mutant that reaches MORE arms than its row says is itself a finding** — e.g. if M-BUILDTAG also
emitted the A5 message, or if M-RETURN-FUNCLIT emitted anything. The planner observed exactly one
problem line per red arm and zero on the green arm. The executor must assert the same.

### 4.2 The two differ-tests — both planner-measured, both PASS

Neither is eyeballed; both use a shell `-ne` assertion.

**AC4b — A7 is non-vacuous ACROSS the parse mode.** A7's ability to fire is entirely a property of
the parse mode, so a red on M-BUILDTAG is evidence only if it would NOT have been red under the old
mode. With M-BUILDTAG in place, run the fence once at `parser.ParseComments` and once at `0`:
```
[ "$rc_comments" -ne "$rc_mode0" ] || { echo "A7 IS VACUOUS"; exit 1; }
```
Planner-measured: **`rc_comments=1`, `rc_mode0=0` → DIFFER, PASS.** Restore `:377` byte-identically.

**AC4d(iii) — A6's `*ast.FuncLit` stop is load-bearing.** With M-RETURN-FUNCLIT in place, run the
fence once against the shipped A6 and once with the `if _, ok := n.(*ast.FuncLit); ok { return
false }` guard deleted:
```
[ "$rc_shipped" -ne "$rc_descend" ] || { echo "THE FUNCLIT STOP IS NOT LOAD-BEARING"; exit 1; }
```
Planner-measured: **`rc_shipped=0`, `rc_descend=1` → DIFFER, PASS.** The descending variant emits
`return statement count=1, want 0 (an early return neuters the assertion)` over a canary that
`--- PASS`es — a textbook row-55 false red. Restore the helper byte-identically.

### 4.3 Restore discipline (the planner followed it; the executor must too)

Every arm: `cp` the pristine backup back, then `shasum -a 256` must read
`a23cfa79419ae69136e62981d5bd0c8ea68cdf2e154ba973fb99a3d3c8b47bfd` for the canary and
`f0667a280df795a1da16d0374a2ccb78ef0e0892d926f7a68817b1ea2fc8e70b` for the pre-sprint helper, and
`git status --porcelain` must be re-checked. The planner's final state is clean (§6).

### 4.4 Declared residuals — demonstrations, NOT reds

`M-ALIAS` (`s := t.Skip; s("MUTANT")`) and `M-HELPER` (`helperSkip(t)`) both leave the extended fence
GREEN. They are optional demonstrations that the declared interprocedural/alias residuals are real.
They are **not** acceptance gates and a green there is not a failure.

---

## 5. Refuted doc premises — with the measurement

The design doc is quorum-cleared, not correct. Everything it cites was re-derived first-party.

**Confirmed unchanged** (no refutation): fence at `:456`, helper at `:375`, `parser.ParseFile` at
`:377` with mode `0`, `return nil` at `:421`; the helper has NO `problems` accumulator and returns a
fresh `[]string{…}` at each of three early failure points; `funcs` is populated by scanning `f.Decls`
for a `*ast.FuncDecl` named `TestToolchainCanary`; `fmt`/`go/ast`/`go/parser`/`go/token`/`strings`
all already imported; the canary contains zero func literals and its only `return` is at `:59`
(`func (s canaryString) String() string { return string(s) }`), OUTSIDE `TestToolchainCanary`;
`ParseComments` already used at `cmd/world-publish/wiring_test.go:60` and
`host/broker/invoke_boundary_test.go:112,:317`. V34's five-arm vet/compile table reproduced exactly.

| # | Doc / brief claim | Measurement | Consequence for the sprint |
|---|---|---|---|
| **R1** | The brief's fact-2 sha table pins `M-SKIP` = `9203903fe73c4e24` and `M-SKIPF` = `8842b6bacc6243b2`, while the doc's Non-Vacuity table names the mutant texts `t.Skip("MUTANT: reachability probe")` and `t.Skipf("MUTANT: reachability probe")`. **These do not correspond.** | Planner ran four texts through the same insertion point: `t.Skip("MUTANT")` → `9203903fe73c4e24`; `t.Skip("MUTANT: reachability probe")` → `3d7c343ce10c0472`; `t.Skipf("MUTANT %s","x")` → `8842b6bacc6243b2`; `t.Skipf("MUTANT: reachability probe")` → `74485845cdbd561a`. | **The sha is a fingerprint of the mutant TEXT, not of the arm.** The executor must pin the exact text from §4.1 and derive its own shas. Copying a sha from the doc or the brief and asserting on it will red for the wrong reason. §4.1's sha column is planner-derived for §4.1's exact texts. |
| **R2** | Doc V32 cites the fence failure at `toolchain_pin_gate_test.go:520`; V33 cites `:517`. | With the planner's patch the same failures print at `:517` (shipped) and `:514` (A6-descend variant). | **Line numbers in the `--- FAIL` output are a patch-layout artifact** (they move with the comment-line count of the inserted block). Assert on the MESSAGE SUBSTRING only, never on the line number. §4.1's expected-FAIL column is substring-only by construction. |
| **R3** | Doc "Files to Create/Modify": `(+~40 LOC)`. | Planner's in-situ patch: 706 → **759 lines, +53**. | Cosmetic; the estimate is low by ~30%. Not a blocker, but the executor should not treat a +53 diff as scope creep. |
| **R4** | Doc P3 and AC4 say "all **five** pass the CURRENT fence (V10)", citing V10 — whose own command list ran only **four** arms (`m_skip`, `m_skipnow`, `m_return`, `m_buildtag`), never `m_skipf`, and never either FuncLit arm. | Planner measured the CURRENT (unpatched) fence against **all seven** arms: rc=0 on every one (§4). | The claim is TRUE but was OVER-CITED — V10 did not support it. It is now supported by a first-party seven-arm measurement, which is strictly stronger. No change to the design. |
| **R5** | Doc's Decision code block places `var problems []string` in the MIDDLE of the A5 comment block (5 comment lines, then the declaration, then `skipCalls := 0`). | Planner applied it that way: `gofmt -l` empty, `go vet` rc=0 — it is legal and formats clean. | Cosmetic only. Recommend hoisting `var problems []string` ABOVE the A5 comment so the comment describes the code it introduces. Either form passes every gate. |
| **R6** | Doc AC1 base evidence cites V7's command `go test -count=1 -v -run TestCanaryDeclaresPositiveArmOnly ./host/verifygate/` (unanchored `-run`), while AC1 itself specifies the anchored `'^…$'` form. | Both run; the planner baselined the **anchored** form (rc=0, exactly 1 `=== RUN`). | Use the anchored form everywhere. The unanchored pattern would also match any future `TestCanaryDeclaresPositiveArmOnlyXxx`, weakening the run-existence count. |

No refutation touches the design DIRECTION: three statically-visible reachability zero-needles folded
into the same AST pass, with A5 descending and A6 stopping at `*ast.FuncLit`. That direction is
confirmed correct by seven in-situ arms and two differ-tests.

---

## 6. Planner's final worktree state

```
$ git status --porcelain
?? design_docs/planned/w-canary-fence-blind-to-a-skipped-canary.md
?? design_docs/planned/w-canary-fence-blind-to-a-skipped-canary-SPRINT.md

$ shasum -a 256 host/verifygate/toolchain_pin_gate_test.go host/store/toolchain_canary_test.go
f0667a280df795a1da16d0374a2ccb78ef0e0892d926f7a68817b1ea2fc8e70b  host/verifygate/toolchain_pin_gate_test.go
a23cfa79419ae69136e62981d5bd0c8ea68cdf2e154ba973fb99a3d3c8b47bfd  host/store/toolchain_canary_test.go
```

Both source files restored byte-identically after all nine mutation/variant arms. Only the two
design docs are untracked; nothing tracked is modified.

---

## 7. Definition of done

1. One commit, one changed `.go` file (`host/verifygate/toolchain_pin_gate_test.go`), four edits.
2. Every row of §3 at its required post-sprint rc.
3. Every RED arm in §4.1 at fence rc=1 with its exact message substring, and exactly one problem line.
4. `M-RETURN-FUNCLIT` at fence rc=**0** — the must-stay-green control.
5. Both differ-tests in §4.2 asserted with `-ne`, both DIFFER.
6. Canary and helper restored sha256-byte-identical after every arm; final `git status --porcelain`
   shows only the intended change plus the design docs.
7. `./scripts/verify_go.sh` is NOT run as a gate (§0.3), and `go build` appears nowhere as a compile
   fence (§0.1).

---

## 8. ROUND 2 — the `sonnet` evaluator FAILED round 1 at 52/100, and both blocking findings were real

Round 1 (commit `8d7b110`) shipped A5/A6/A7 and passed every gate in §3 and every arm in §4.1.
The adversarial evaluator (`sonnet`; generator≠judge) scored it **52/100 FAIL** on two arms the
plan's mutation table never ran. **Both were reproduced first-party by the controller before being
actioned** (ghost discipline — an evaluator claim is a lead, not a fact), and both are now closed.

### 8.1 BLOCKING 1 (CONFIRMED) — `goto` jumps past the assertion and the fence stayed GREEN

The mutant (`V35`, canary sha256(16) `f7d6d640257d9f61`) inserts `goto End` above the assertion and
`End:` after it. Controller re-measurement, first-party:

| | |
|---|---|
| `go test -count=1 -run '^$' ./host/store/` | **rc=0** — it compiles |
| `go vet ./host/store/` | rc=1 `unreachable code` — **statically visible** |
| canary | **`--- PASS`** |
| extended fence (round 1) | **rc=0 `--- PASS`** — green over a canary that asserts nothing |

This is **worse than the `t.Skip` hole the row exists to close**: a skip at least prints `--- SKIP`,
while a `goto` reports a confident `--- PASS`. It sits squarely inside the row's own declared scope
("statically visible reachability"), and `scripts/verify_go.sh` does not run `go vet`, so nothing in
the enforced gate would have caught it either. **The evaluator's refutation of the round-1 thesis is
accepted: three checks did not close statically-visible reachability. Four do.**

**Fix — A8**: zero-needle on `*ast.BranchStmt` with `Tok == token.GOTO` in the body, **stopping at
`*ast.FuncLit`** (Go forbids a goto crossing a function boundary, so a closure-local goto cannot
skip the outer assertion; counting it would be a row-55-class false red — proven by `D2`).

### 8.2 BLOCKING 2 (CONFIRMED) — A7's byte-prefix match false-redded a benign doc comment

`strings.HasPrefix(c.Text, "// +build")` fires on `// +buildAlerts is an unrelated internal
codename…` (`V36`, sha256(16) `50f83fc840f7f68f`), which **Go does not read as a constraint**:

```
go/build/constraint.IsPlusBuild("// +buildAlerts is an unrelated internal codename")  ->  false
go/build/constraint.IsGoBuild  ("//go:buildFoo bar")                                  ->  false
go/build/constraint.IsGoBuild  ("//go:build ignore_canary")                           ->  true
go/build/constraint.IsPlusBuild("// +build linux")                                    ->  true
```

Measured on that mutant: `go vet` **rc=0** (Go's own `buildtags` analyzer is silent), the file is
**not excluded** (canary `--- PASS`), yet round-1 A7 reported **rc=1**. A textbook row-55 false red,
caused by a byte prefix standing in for a grammar.

**Fix**: match with `go/build/constraint` — Go's own parser — instead of `strings.HasPrefix`.
One new import (`go/build/constraint`); this supersedes §1's "no new import" line.

### 8.3 A THIRD NARROWING WAS WRITTEN, MEASURED, AND DELETED RATHER THAN SHIPPED (`V37`)

Alongside 8.2 the controller added a **position** guard (`cg.End() >= f.Package -> continue`), on
the theory that a grammar-valid `//go:build` quoted in prose *after* the package clause is inert and
should not red. Its differ-test "passed" (rc=0 with the guard, rc=1 without) — **and the pass is
vacuous.** Go rejects a misplaced `//go:build` outright, both after the package clause and inside a
function body:

```
host/store/toolchain_canary_test.go:21:2: misplaced //go:build comment
go vet rc=1 · go test -count=1 -run '^$' ./host/store/ rc=1 · the canary never runs
```

So the guard can only change the verdict on a tree **that does not build**, and a gate's reading on
an unbuildable tree is not a verdict. Shipping it would have added a branch **no arm can ever
reach** — the anti-vacuity-floor class this charter tracks in its Repo Profile watch-item, arriving
in the fix for a false red rather than in the gate it protects. **It was deleted.** The lesson is
that a differ-test only certifies a branch if the tree it differs on is one the toolchain accepts:
*pair every differ-test with the compile fence on the SAME tree, or a rejected tree will sell you a
dead branch.*

### 8.4 Round-2 arm table — 12 arms, every mutant compiles (`CF=0` on all 12)

| Arm | sha256(16) | vet | canary | fence rc | owner |
|---|---|---|---|---|---|
| M-SKIP | `3d7c343ce10c0472` | 0 | `--- SKIP` | **1** | A5 |
| M-SKIPF | `9c8778abdd21c8ae` | 0 | `--- SKIP` | **1** | A5 |
| M-SKIPNOW | `8311cfeff59be3ae` | 0 | `--- SKIP` | **1** | A5 |
| M-SKIP-FUNCLIT | `13c85ea2e880f350` | 0 | `--- SKIP` | **1** | A5 (descends) |
| M-RETURN | `2ff015dccd545c3c` | 1 † | `--- PASS` | **1** | A6 |
| M-BUILDTAG | `43596b333f39ee1f` | 0 | no test ran | **1** | A7 (modern) |
| M-BUILDTAG-PLUS | `bfc770f860a5b455` | 0 | no test ran | **1** | A7 (legacy) |
| **M-GOTO** | `f7d6d640257d9f61` | 1 † | `--- PASS` | **1** | **A8 (new)** |
| M-RETURN-FUNCLIT | `7cbabbb1378e2d7e` | 0 | `--- PASS` | **0** | must stay green |
| **M-GOTO-FUNCLIT** | `126b3fde0eab654c` | 0 | `--- PASS` | **0** | must stay green |
| **M-PROSE-DOC** | `50f83fc840f7f68f` | 0 | `--- PASS` | **0** | must stay green (8.2) |
| **M-GOBUILD-LOOKALIKE** | `4e96af10770a5197` | 0 | `--- PASS` | **0** | must stay green |

† `unreachable code` — the ANALYZER, not a compile failure; both arms are `CF=0`. Not a gate.

Each red arm produced **exactly one** `--- FAIL` and **exactly one** problem line.

### 8.5 Non-vacuity — each check neutered in turn, red set is EXACTLY the arms it owns

Neutering `if <counter> != 0` to `if false && <counter> != 0`:

| Neutered | Arms that went GREEN | Arms owned |
|---|---|---|
| A5 `skipCalls` | M-SKIP M-SKIPF M-SKIPNOW M-SKIP-FUNCLIT | identical |
| A6 `returns` | M-RETURN | identical |
| A7 `buildTags` | M-BUILDTAG M-BUILDTAG-PLUS | identical |
| A8 `gotos` | M-GOTO | identical |

No neuter moved an arm it does not own — so no check is carrying another's evidence.

### 8.6 Four differ-tests, all asserted with `-ne`

| | shipped | variant | verdict |
|---|---|---|---|
| **D1** A7 grammar vs `strings.HasPrefix` (on M-PROSE-DOC) | rc=0 | rc=1 | DIFFER — the grammar match is load-bearing |
| **D2** A8's FuncLit stop (on M-GOTO-FUNCLIT) | rc=0 | rc=1 | DIFFER — prevents a false red |
| **D3** A7 parse mode `ParseComments` vs `0` (on M-BUILDTAG) | rc=1 | rc=0 | DIFFER — A7 not vacuous |
| **D4** A6's FuncLit stop (on M-RETURN-FUNCLIT) | rc=0 | rc=1 | DIFFER — load-bearing |

### 8.7 Non-blocking evaluator findings, accepted as declared residuals

- `panic()` + deferred `recover()` neuters the assertion and leaves the fence green — the design doc
  declared it; the evaluator **measured** it first-party (canary `--- PASS`, fence `--- PASS`). Out
  of scope, correctly declared.
- The shape checks still `return` early, so A5–A8 do not run when a shape check fails. The evaluator
  confirmed this masks the *message*, never the *verdict*: with both a shape break and a skip landed,
  the fence still reds (rc=1) and merely reports the shape problem. Fail-closed, so it stands.
- A5's receiver is hardcoded to the ident `t`. Pre-existing: row 49's `isTFatalfCall` has the
  identical hardcoding, so a renamed receiver breaks the whole fence, not just the new checks. Not
  introduced here; not closed here.
- `M-ALIAS` (`s := t.Skip; s("x")`) and `M-HELPER` (`helperSkip(t)`) remain green by design.

### 8.8 Gates, round 2

`gofmt -l` **0 lines** · `go vet ./host/verifygate/ ./host/store/` **rc=0** ·
`go test -count=1 ./host/verifygate/ ./host/store/` **rc=0, 0 `--- FAIL`** with
`AILANG_BIN=/tmp/ailang-v0300/ailang` (`AILANG v0.30.0`) ·
`./scripts/verify_ail.sh` **rc=0** (11 required identities verified, 40 named tests pass) — run even
though no `.ail` changed, as the repo-wide control.
