# Sprint plan — w-dispatch-lever-parser-false-reds-on-valid-yaml (queue row 55)

- **Design doc:** `design_docs/planned/w-dispatch-lever-parser-false-reds-on-valid-yaml.md` (quorum-cleared round 2, carve-out applied)
- **Worktree / branch:** `.wt-iter149-row55` on `sprint/w-dispatch-lever-parser` @ `95c17be`
- **Base commit all baselines measured at:** `95c17be` (design commit on top of `234d9da`; the design commit touches only `design_docs/planned/*.md`, so every code baseline is `234d9da`'s)
- **Size:** ~0.2d · 3 milestones, each independently committable and testable
- **Scope class:** BRITTLENESS, not soundness. Not re-scoped.
- **Constraints honoured:** no `go.mod`/`go.sum` change, no new dependency, `tools/launchd/*` untouched, no `.ail` change.

---

## 0. Planner's verification pass — what I re-derived, and what I refuted

**Rule:** the design doc is a claim, not a source. Every file:line, grep count and behavioural
row below was re-measured first-party in this worktree by the planner before any AC was written.

### 0.1 CONFIRMED (re-derived first-party, doc was right)

| Doc claim | Planner's command | Result |
|---|---|---|
| helper `onBlockTriggerKeys` at L14–84; caller `TestEveryWorkflowDeclaresDispatchLever` at L112–136 | `cat -n host/verifygate/dispatch_lever_gate_test.go` | CONFIRMED (`func` L14, closing `}` L84; `func Test…` L112, `}` L136 — file is 136 lines). *(The controller's task brief said the helper is "L14-83"; the DOC's "L14–84" is the correct one. The doc wins.)* |
| `l == "on:"` at L25; `TrimLeft(…, " ")` at L43 and L50; L128 control `strings.Contains(src, "on:")`; Glob at L113; ANTI-VACUITY FLOOR at L117; one call site at L131 | same | ALL CONFIRMED, exact |
| L110–111 comment quoted verbatim in §1(e) | `sed -n '110,111p'` | CONFIRMED byte-for-byte |
| V1 quoted `"on":` Fatals | throwaway `TestZZIter149Probe`, `go test ./host/verifygate/ -run TestZZIter149Probe -v` | ``instrument failure: V1_quoted.yml has no top-level `on:` trigger block`` |
| V2 flow style Fatals | same | identical Fatalf |
| V3 tab indent returns `keys=[]`, **no** Fatal | same | `RESULT V3_tab keys=[]`, subtest PASS |
| V4 canonical control | same | `keys=[push workflow_dispatch]` |
| V5 real `ci.yml` control | same, `os.ReadFile(repoRoot/.github/workflows/ci.yml)` | `keys=[push pull_request workflow_dispatch]` |
| V6 scalar arm cascades (2 messages) | same | ``V6_scalar.yml: `workflow_dispatch:` has scalar value "garbage"…`` **and** `keys=[push]` (caller then emits the second) |
| V21 `on: push` / `on: [push, …]` Fatal with the misleading message | same | both `has no top-level \`on:\` trigger block` |
| V22 inline + own-line comment both green at base | same | both `keys=[push workflow_dispatch]` |
| V9 the L128 control | same probe: `strings.Contains` on 3 fixtures | `quoted+runs-on=true` · `quoted+NO-runs-on=false` · `real ci.yml=true` — **exactly as V9 states** |
| V7 / V8 / V10 / V11 / V13 / V16 / V17 / V18 / V19 / V23(a) | see §0.3 for the raw numbers | ALL CONFIRMED, including every file:line in V13's 7-read-site set |
| Baseline `AILANG_BIN=… go test ./host/verifygate/` rc=0 | re-run by the planner, not deferred to M2 | **rc=0 in 114.9s** (doc/brief said ~175s; the number is not load-bearing, the rc is) |
| Baseline `go test ./host/verifygate/` without `AILANG_BIN` | `go test ./host/verifygate/ -count=1` | **rc=1, exactly 17 `--- FAIL`, all `AILANG_BIN is unset`** — confirms every test AC must set it |
| `go build ./...` at base | `go build ./...` | rc=0, worktree porcelain 0 after probe deletion |

Probe file `host/verifygate/zz_iter149_probe_test.go` was deleted; `git status --porcelain | wc -l` → `0`.

### 0.2 PREMISES REFUTED — defects this plan repairs

Five. Three are defects the **doc introduces**; they are fixed in the milestones below. Two are
pre-existing properties of HEAD and go to the queue (§6), not to a milestone.

---

**R-1 (doc-introduced, SEVERE) — AC4's build-clean fence is VACUOUS: `go build ./...` does not
compile `_test.go` files.**

The doc's AC4 reads *"every mutant builds (`go build ./...` rc=0 on the mutated tree — a
non-building mutant proves nothing)"*. The entire mutation table lives in a `_test.go` file, which
`go build` never compiles. Measured:

```
$ printf '\nfunc zzBroken() { var x int = "not an int"; _ = x }\n' >> host/verifygate/dispatch_lever_gate_test.go
$ go build ./... ; echo rc=$?
rc=0                       # <-- a type error in the mutated file, and the fence reads GREEN
$ go vet ./host/verifygate/
vet: host/verifygate/dispatch_lever_gate_test.go:138:31: cannot use "not an int" … as int value
$ go test -count=1 -run '^$' ./host/verifygate/
host/verifygate/dispatch_lever_gate_test.go:138:31: cannot use "not an int" …
FAIL  github.com/sunholo-data/ailang-world/host/verifygate [build failed]
```

So AC4's stated purpose — *catch a mutant that does not build* — is exactly what it cannot do.
**Fix:** this plan's build-clean fence is `go test -count=1 -run '^$' ./host/verifygate/` → rc=0
(compiles the test binary, runs nothing; measured green at base in 4.0s, prints
`ok … [no tests to run]`). `go vet ./host/verifygate/` is the equivalent alternative (8.7s
repo-wide at base, rc=0). `go build ./...` is retained only as a non-test-code fence, never as the
mutant fence.

---

**R-2 (doc-introduced, SEVERE) — MUT-E as specified has an EMPTY red set: nothing can reach the
L128 control with a fixture.**

MUT-E's expected red is *"table case: quoted key via the GATE-path wrapper arm, using the
`runs-on:`-FREE fixture"*. But the L128 control lives inside
`TestEveryWorkflowDeclaresDispatchLever`, whose input is `filepath.Glob(repoRoot/.github/workflows/*)`,
and `repoRoot` is a package-level `var` computed from `runtime.Caller` — not injectable:

```
$ sed -n '26,29p' host/verifygate/ail_binary_gate_test.go
var (
    repoRoot = findRepoRoot()
```

The only way a fixture reaches L128 is to write a file into `.github/workflows/` — which
`toolchain_pin_gate_test.go:209-211` forbids (measured, verbatim):

```go
if !slices.Equal(workflowFiles, []string{"ci.yml"}) {
    t.Errorf("workflow files=%v, want exactly [ci.yml]; a second workflow may carry unscanned toolchain pins", workflowFiles)
}
```

and which the doc's own §5 correctly rules out. So as written MUT-E reds **nothing** — it is an
undetectable mutant dressed as evidence for the shared needle set.

**Fix (M1 deliverable):** extract the L128 control into a pure predicate
`srcDeclaresOnKey(src string) bool`, call it from L128, and give it its own table test
`TestOnBlockControlNeedle`. Then MUT-E has a real, single-subtest red set. This is the *minimum*
change that makes the doc's own MUT-E executable; it is not scope growth.

---

**R-3 (doc-introduced) — MUT-A and MUT-E collide if the needle set is mutated as a slice; and
AC2's `"on":` needle FAILS ON A CORRECT IMPLEMENTATION.**

*(a) Collision.* §3.2 makes `acceptedOnForms` a slice shared by `matchTopLevelOn` and the L128
control. MUT-E is labelled *"(shared-slice mutant)"*. Deleting the quoted needles from a **shared**
slice narrows BOTH sites, so MUT-E's red set would strictly contain MUT-A's and the two mutants
would no longer be separable — the mutation table would be measuring one edit twice. **Fix:** the
mutation protocol below forbids editing `acceptedOnForms` itself. MUT-A inserts a guard *inside
`matchTopLevelOn`'s loop*; MUT-E inserts a guard *inside `srcDeclaresOnKey`*. Both leave the slice
literal intact, so the two red sets are disjoint (enumerated in §4).

*(b) AC2's needle.* AC2 greps `'errTabIndent\|errUnhandledOnForm\|"on":'`. The third alternative
matches the five bytes `"on":`. A conventional Go interpreted string literal for that form is
written `"\"on\":"`, whose bytes are `"`,`\`,`"`,`o`,`n`,`\`,`"`,`:`,`"` — the needle does **not**
occur. Measured:

```
$ printf 'x := "\\"on\\":"\ny := `"on":`\n' > /tmp/gostr.txt
$ grep -n '"on":' /tmp/gostr.txt
2:y := `"on":`          # only the BACKTICK raw-string form matches; the escaped form does not
```

So AC2 silently requires a raw-string literal. A correct implementation using the escaped form
reds AC2. **Fix:** AC2's needle set drops the `"on":` literal and uses only compiler-enforced
identifiers (`errTabIndent`, `errUnhandledOnForm`, `matchTopLevelOn`, `srcDeclaresOnKey`); quoted-form
coverage is proven by the named table subtests plus MUT-A, which is where it belongs.
*(The rig's `grep` is `ugrep 7.8.4`; BRE `\|` alternation works — measured `grep -c 'alpha\|beta'` → 2.
Note `grep -c` counts matching LINES, not hits, so the AC is restated as per-needle presence.)*

---

**R-4 (doc-introduced, minor but destructive if executed literally) — §3.4 item 1's edit range
`291-293` swallows a TRUE sentence.**

Measured (`sed -n '288,295p' design_docs/planned/w-ci-recovery-lever-absent.md`):

```
291  derived set, and EXACTLY ONE attributed message per defect (the `t.Errorf` names the file
292  and the absent lever; it never cascades — P5's precedent). It reuses the sibling's import set
293  (P15), so no imports change.
```

The wrong claim ends mid-L292 at `P5's precedent).`. L292–293 then carry a **separate, true**
sentence about the import set. V11's "sentence spans 291–293" is wrong; a wholesale replacement of
291–293 deletes a correct statement. **Fix:** M3's edit is pinned to the clause
`and EXACTLY ONE attributed message per defect (… it never cascades — P5's precedent).` and must
leave `It reuses the sibling's import set (P15), so no imports change.` intact. AC8 gains a control
for exactly this.

---

**R-5 (doc-introduced, minor) — "byte-for-byte message preservation" is asserted by NO acceptance
criterion.**

§3.2 promises the wrapper *"preserves TODAY'S messages byte-for-byte for the existing arms (so
row-47's mutation evidence for those arms stays valid)"*. No AC checks it, and no test in the repo
asserts those strings today — so the promise cannot fail, and a refactor that silently reworded
`instrument failure: %s has no top-level \`on:\` trigger block` would pass every AC in the doc.
**Fix:** M1 gains a small pure formatter and `TestOnBlockFailureMessagesUnchanged` (2 subtests),
and AC1's name list includes it.

### 0.3 Raw re-measurements backing §0.1 (planner-run, this worktree)

```
V7/V13  grep -rn 'workflows"' --include='*.go' .
        -> 7 read sites, 3 files, at EXACTLY the lines V13 names:
           toolchain_pin_gate_test.go:111,201,501 · ail_binary_gate_test.go:669
           dispatch_lever_gate_test.go:113 · runbook_stageb_test.go:339,361
V16     grep -rn onBlockTriggerKeys --include='*.go' .  -> 3 hits (L11 doc, L14 def, L131 call)
V17     grep -ci yaml go.mod -> 0   (control modernc -> 4)
        grep -rln yaml --include='*.go' . -> 1 file, the gate itself
V8      go.mod: 1 direct require (modernc.org/sqlite v1.54.0) + 9 `// indirect`
V23(a)  go list -deps ./... | wc -l -> 257 ; | grep -ci yaml -> 0 ; modernc -> 30 ; encoding/json -> 1
V18     dispatch_lever_gate_test.go: 'skips dotfiles'->0  'invisible'->0
        control 'nested subdirectory'->1 · 'a hidden file'->1 · 'the Glob is case-sensitive'->1
V19     same file: 'anti-vacuity floor' (cs) -> 0 ; 'anti-vacuity' (-ci) -> 1
        repo *.go: cs -> 2 (evidence_manifest_gate_test.go) ; -ci total -> 22
V10     w-ci-recovery-lever-absent.md:424 'skips dotfiles' ; :425 'invisible' ; block 422-427
V11     'EXACTLY ONE attributed message' -> 1 (L291) ; 'never cascades' -> 2 (L292, L503)
AC7     'is a directory' -ci -> 0 in BOTH target files  (so AC7's presence needle CAN fail)
ci.yml  one file in .github/workflows/ ; bare `on:` at col 0 ; 2x `runs-on:` ;
        `workflow_dispatch:` has NO value  (=> MUT-D leaves the gate green — control holds)
```

**Nothing else in the doc was refuted.** Specifically checked and found sound: the D-WORLD-30 quote
(`design_docs/world-mission.md:740`, verbatim), the transferable/not-transferable split, the charter
row-55 text (`design_docs/world-mission.md:4829`), §3.3's per-shape dispositions, the `oc-glm-5-2`
flow-guard specification, §5's Conflict Surface (all 7 sites, incl. `host/runbook`), R1–R5, and the
needle audit's reasoning about zero-needles vs presence-needles.

---

## 1. Where the plan overrides the design doc

| Topic | Doc says | This plan says | Winner |
|---|---|---|---|
| Mutant build fence | `go build ./...` rc=0 | `go test -count=1 -run '^$' ./host/verifygate/` rc=0 | **PLAN** (R-1, measured) |
| MUT-E reachability | gate-path wrapper arm with a fixture | requires extracting `srcDeclaresOnKey`; fixture-driven table case | **PLAN** (R-2, measured) |
| MUT-E mechanism | "shared-slice mutant" | guard inside `srcDeclaresOnKey`; `acceptedOnForms` never edited | **PLAN** (R-3a) |
| AC2 needles | includes literal `"on":` | identifiers only | **PLAN** (R-3b, measured) |
| Row-47 edit range | replace 291–293 | replace the clause only; keep the import sentence | **PLAN** (R-4, measured) |
| Message preservation | promised, unchecked | `TestOnBlockFailureMessagesUnchanged` + AC1 | **PLAN** (R-5) |
| helper line span | L14–84 | (brief said L14-83) | **DOC** |
| Everything else | — | — | **DOC** |

---

## 2. Milestones

### M1 (0.10d) — pure parse core, shared needle set, the three hardenings

**Files:** `host/verifygate/dispatch_lever_gate_test.go` only. Stdlib imports may be added
(`errors`, `fmt`); **no** `go.mod`/`go.sum` change.

**Deliverables**

1. `var acceptedOnForms = []string{…}` — the three accepted column-0 forms (bare, double-quoted,
   single-quoted). **This literal is never edited by any mutant** (R-3a).
2. `func matchTopLevelOn(line string) (rest string, ok bool)` — column-0-anchored, exact accepted
   forms only, consuming `acceptedOnForms`. Preserves P14's rationale (a trimmed comparison would
   anchor on any nested `on:`). Feeds BOTH the block finder and the duplicate counter.
3. `func srcDeclaresOnKey(src string) bool` — the extracted L128 known-positive control, consuming
   `acceptedOnForms` **directly, not via `matchTopLevelOn`** (this is what keeps MUT-A and MUT-E
   disjoint). L128 becomes `if !srcDeclaresOnKey(src) { t.Fatalf(…) }` with its message unchanged.
4. Sentinels `errNoOnBlock`, `errDuplicateOnBlock`, `errTabIndent`, `errUnhandledOnForm`, carried by
   a typed error that `Unwrap()`s to the sentinel (so `errors.Is` works) and whose `Error()` renders
   the user-facing text.
5. `func parseOnBlockTriggers(src string) (keys []string, scalarValued map[string]string, err error)`
   — no `*testing.T`, no I/O. Implements §3.3: quoted forms accepted; flow mapping accepted **only**
   when comment-stripped, trimmed `rest` both starts `{` and ends `}` (`oc-glm-5-2`'s verbatim
   guard) else `errUnhandledOnForm`; any `\t` in a scanned line's indentation prefix →
   `errTabIndent`; any other non-empty `rest` → `errUnhandledOnForm`; the scalar-value rule applies
   in block and flow, with inline `#` stripped before judging (the `eb215c3` fix, V22).
6. `func onBlockFailureMessage(path string, err error) string` — pure formatter; the wrapper
   `t.Fatalf("%s", onBlockFailureMessage(path, err))`.
7. `func onBlockTriggerKeys(t *testing.T, path, src string) []string` — unchanged signature and
   unchanged behaviour on every arm that exists today; emits one `Errorf` per `scalarValued` entry
   (the cascade is KEPT deliberately, §3.3 row d).
8. Three table tests with **pinned subtest names** (the mutation table enumerates these by name):

   `TestOnBlockTriggerParserShapes`:
   `canonical_block` · `real_ci_yml` · `quoted_double` · `quoted_single` · `flow_mapping` ·
   `flow_scalar_violation` · `flow_unclosed` · `tab_indent` · `scalar_on` · `sequence_on` ·
   `duplicate_mixed` · `no_on_block` · `block_scalar_violation` · `inline_comment_value`

   `TestOnBlockControlNeedle`:
   `control_quoted_no_runs_on` (true — base `strings.Contains(src,"on:")` is **false** here, V9) ·
   `control_real_ci_yml` (true) · `control_absent` (false — anti-vacuity for the needle itself)

   `TestOnBlockFailureMessagesUnchanged`:
   `msg_no_on_block` · `msg_duplicate_on_block` — assert the rendered strings are byte-identical to
   today's L31 and L41 format strings.

**Fence:** M1 must NOT touch `dispatch_lever_gate_test.go` L108–111 (M3 owns that comment, and
AC6's Go-site base counts are measured against HEAD).

- **AC1** — `AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -count=1 -v -run
  'TestOnBlockTriggerParserShapes|TestOnBlockControlNeedle|TestOnBlockFailureMessagesUnchanged|TestEveryWorkflowDeclaresDispatchLever'`
  → rc=0 **AND** the `-v` output contains a `--- PASS:` line for **all four** names.
  *(Base timing for the existing arm alone: 0.19s test, 19.8s wall incl. build.)*
  *Still passes if false:* rc=0 on zero matched tests — which is precisely why the four-name
  `--- PASS` check is the load-bearing half. A refactor handling none of the three shapes still
  passes `TestEveryWorkflowDeclaresDispatchLever` (ci.yml is canonical, V5) — the table subtests are
  what red; a MISSING subtest passes everything, which is why AC2 pins names and M2 pins mutants.
- **AC2 (case-presence floor, repaired per R-3b)** — each of `errTabIndent`, `errUnhandledOnForm`,
  `matchTopLevelOn`, `srcDeclaresOnKey`, `parseOnBlockTriggers` returns `grep -c <ident>
  host/verifygate/dispatch_lever_gate_test.go` ≥ 1, **and** each of the 14 + 3 + 2 subtest names
  above returns ≥ 1. All needles are Go identifiers or string literals this sprint itself writes, so
  casing is compiler-/author-enforced (case-sensitive is correct here).
  *Still passes if false:* grep proves text, not behaviour — AC2 exists only to make the M2 red sets
  addressable by name; behaviour is AC4's job.
- **AC3 (no-dependency fence)** — `git diff --stat -- go.mod go.sum` → empty (base: empty, measured).
  *Still passes if false:* nothing; any dependency addition touches `go.mod`.
- **AC3b (test-compile fence, replaces the doc's `go build` fence)** — `go test -count=1 -run '^$'
  ./host/verifygate/` → rc=0 (base: `ok … [no tests to run]`, 4.0s) **and** `go build ./...` → rc=0.
  *Still passes if false:* `go build ./...` alone would pass a test file that does not compile
  (R-1, measured) — the `go test -run '^$'` half is what fails.
- **AC3c (comment fence)** — `git diff -U0 -- host/verifygate/dispatch_lever_gate_test.go` in M1's
  commit contains no hunk overlapping L108–111 of the base file.
  *Still passes if false:* nothing mechanical if the hunk header drifts; the evaluator reads the
  diff. Cheap belt-and-braces so AC6's base counts survive M1.

### M2 (0.05d) — mutation run + full-package regression

**Mutation protocol (binding):**
- Neuter by **inserting a guard** (`if false && …`, or an early `continue`) that keeps every
  identifier and import in use. **Never delete code, and never edit `acceptedOnForms`** (R-3a).
- After each mutant: `go test -count=1 -run '^$' ./host/verifygate/` → rc=0 (build-clean, R-1),
  then `AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -count=1 -v -run
  'TestOnBlockTriggerParserShapes|TestOnBlockControlNeedle|TestOnBlockFailureMessagesUnchanged|TestEveryWorkflowDeclaresDispatchLever'`
  and record the **exact set** of `--- FAIL:` subtest names.
- Revert the mutant (`git checkout -- <file>` is forbidden by the controller's rules for this
  worktree; restore from a copy taken before mutating) and confirm porcelain returns to M1's state.
- `-count=1` on every run so no cached verdict is read.
- No `-skip <arm> && rc=0` criterion is used anywhere: every mutant is judged on an ENUMERATED
  failing set, because every mutant reaches ≥1 named subtest and several reach more than one.

**Mutation table (blast radius reasoned per arm; four tests are the complete consumer set — V16
proves `onBlockTriggerKeys` has exactly one production call site, and Go test helpers are not
importable across packages).**

| Mutant | Edit (guard-inserted, must still compile) | EXPECTED FAIL set (exact `--- FAIL:` subtest names) | MUST STAY GREEN |
|---|---|---|---|
| **MUT-A** quoted-form | inside `matchTopLevelOn`'s loop over `acceptedOnForms`, `if form != acceptedOnForms[0] { continue }` | `…ParserShapes/quoted_double`, `…/quoted_single`, `…/duplicate_mixed` | every other subtest; **all of `TestOnBlockControlNeedle`** (proves the control is fed by the slice, not by `matchTopLevelOn`); `TestEveryWorkflowDeclaresDispatchLever` (ci.yml is bare `on:`) |
| **MUT-B** flow-parse | flow-mapping branch behind `if false && …` (falls through to `errUnhandledOnForm`) | `…ParserShapes/flow_mapping`, `…/flow_scalar_violation` | `…/flow_unclosed` (already `errUnhandledOnForm` — proves MUT-B is specific, not blanket); all other tests |
| **MUT-C** tab-detect | the `\t`-in-indent check behind `if false && …` | `…ParserShapes/tab_indent` (gets `keys=[]`, nil err; wants `errTabIndent`) | all other subtests; all other tests |
| **MUT-D** scalar-value | the value-judging branch behind `if false && …` (key appended, no violation recorded) | `…ParserShapes/block_scalar_violation`, `…/flow_scalar_violation` | `…/inline_comment_value` (must stay NON-violating — proves MUT-D's red is the check vanishing, not the comment-strip); `TestEveryWorkflowDeclaresDispatchLever` (measured: ci.yml's `workflow_dispatch:` carries no value, so the gate is green either way) |
| **MUT-E** control-drift | inside `srcDeclaresOnKey`, `if form != acceptedOnForms[0] { continue }` — narrows the L128 control back to the bare form only. **`acceptedOnForms` itself untouched**, so `matchTopLevelOn` is unaffected and MUT-E's set stays disjoint from MUT-A's | `…ControlNeedle/control_quoted_no_runs_on` — and ONLY that. **Its fixture is load-bearing: `runs-on:`-FREE.** A fixture carrying `runs-on:` satisfies the narrowed needle incidentally (V9, measured: `quoted+runs-on=true`) and the mutant walks straight through | `…ControlNeedle/control_real_ci_yml` and `…/control_absent`; **all** `…ParserShapes` subtests; `TestEveryWorkflowDeclaresDispatchLever` (ci.yml carries both bare `on:` and `runs-on:`) |

- **AC4** — for each of MUT-A…MUT-E: (i) `go test -count=1 -run '^$' ./host/verifygate/` rc=0 on
  the mutated tree; (ii) the observed `--- FAIL:` subtest set equals the EXPECTED FAIL set above
  **exactly** — no more, no fewer; (iii) every name in MUST STAY GREEN appears as `--- PASS:`.
  *Still passes if false:* a mutant that reds EVERYTHING passes a naive "the mutant reds" check —
  clause (ii)'s *exactly* and clause (iii) are what fail it. A mutant that does not compile passes
  the doc's `go build ./...` fence — clause (i) is what fails it (R-1). A mutant whose target case
  does not exist (MUT-E before `srcDeclaresOnKey` is extracted) reds nothing and would pass a
  "no unexpected failures" reading — clause (ii)'s equality in BOTH directions is what fails it.
- **AC5** — on the final, fully reverted tree: `AILANG_BIN=/tmp/ailang-v0300/ailang go test
  ./host/verifygate/ -count=1` → rc=0. **Base measured first-party by the planner: rc=0 in 114.9s**,
  so this AC is red-able and is not vacuous.
  *Still passes if false:* a test DELETED by the diff passes AC5 — AC1's four-name `--- PASS` list
  and AC2's subtest-name presence floor cover that. A hang does not pass (go test's 10m default).
- **AC5b (blast-radius fence, in lieu of a `verify_go.sh` AC)** — `go build ./...` rc=0 and
  `go vet ./...` rc=0 (base: both rc=0, vet 8.7s). **No AC of the form `verify_go.sh rc=0` exists**:
  that script is rc=1 at base on this rig for two reasons owned elsewhere (the `host/broker` flake
  `TestHandlerTimeoutKillsTheWholeProcessGroup`, queue row 58; and the row-54 fleet-drift arm, 759
  diff-lines behind). The change is confined to one `_test.go` file in package `verifygate`, and Go
  test helpers are not importable across packages (V16), so `go test ./host/verifygate/` IS the
  complete blast radius — a whole-suite set-comparison would buy nothing here at ~0.2d and would
  import two unrelated red arms.
  *Still passes if false:* `go vet ./...` would not catch a behavioural regression in another
  package — but no other package can reach this code, which is the measured claim it fences.

### M3 (0.05d) — the three prose corrections, one commit

1. `design_docs/planned/w-ci-recovery-lever-absent.md` — replace **only** the clause
   `and EXACTLY ONE attributed message per defect (… it never cascades — P5's precedent).` (starts
   L291, ends mid-L292) with the measured truth: the scalar arm emits TWO messages (V6),
   deliberately, as defense-in-depth; the absent-lever arm alone is one. **Leave
   `It reuses the sibling's import set (P15), so no imports change.` intact** (R-4). L503's
   historical record stays.
2. Same file, Residual 3, **L422–427** — rewrite per V7: Glob `*` DOES enumerate dotfiles and
   case-mismatched names; a nested subdirectory is a LOUD `t.Fatal` (`is a directory`), not
   invisible. What remains genuinely invisible: a workflow outside `.github/workflows/` entirely —
   labelled as inference from the Glob path literal at `dispatch_lever_gate_test.go:113`.
3. `host/verifygate/dispatch_lever_gate_test.go:110-111` — rewrite the wrong clauses **in the
   comment's own words** (V18: the doc's phrasings grep 0 here): "a hidden file" is seen; "the Glob
   is case-sensitive" is irrelevant to a bare `*`; the nested-subdirectory case is a loud
   `is a directory` Fatal, not unseen.

- **AC6 (absence of the measured-wrong claims, per site; every removal needle has a measured base
  count ≥1 in ITS OWN file, so each zero can fail)** —
  doc site: `grep -c 'skips dotfiles' design_docs/planned/w-ci-recovery-lever-absent.md` → **0**
  (base **1**, L424, measured).
  Go site: `grep -c 'a hidden file' host/verifygate/dispatch_lever_gate_test.go` → **0** AND
  `grep -c 'the Glob is case-sensitive'` (same file) → **0** (base **1** each, measured).
  Known-positive control in the same call, per file: `grep -ci 'anti-vacuity' <file>` ≥ 1 in EACH
  (measured: doc **9**, Go file **1** at L117 written `// ANTI-VACUITY FLOOR:`), read as PRESENCE,
  never as an exact count — the case-sensitive lowercase form greps **0** in the Go file (measured),
  which is the defect this control was repaired for.
  *Still passes if false:* a rewording that keeps the wrong claim without those literal phrases
  ("Glob omits hidden files") — hence AC7.
- **AC7 (human-readable assertion, evaluator-checkable)** — the rewritten Residual-3 block and the
  rewritten code comment each contain BOTH corrected facts: dotfiles/case-mismatches ARE enumerated;
  a nested subdirectory is a LOUD `is a directory` Fatal. Mechanical floor:
  `grep -ci 'is a directory'` ≥ 1 in **each** of the two files. **Base measured: 0 in BOTH files**,
  so this presence needle can fail. `-ci` so a future re-caser cannot make it read zero.
  *Still passes if false:* the phrase present but negated — accepted residual; the evaluator reads
  the sentence. The grep is the floor, not the assertion.
- **AC8** — `grep -c 'EXACTLY ONE attributed message' design_docs/planned/w-ci-recovery-lever-absent.md`
  → **0** (base **1**). Control A: `grep -c 'never cascades'` (same file) → **exactly 1**
  (base **2**: L292 removed, L503's historical record survives). Control B (new, per R-4):
  `grep -c 'It reuses the sibling.s import set' <same file>` → **exactly 1** — the true sentence
  adjacent to the removed clause must survive.
  *Still passes if false:* deleting L503 too gives 0/0 — Control A's *exactly 1* fails that.
  Replacing the whole 291–293 block deletes the import sentence — Control B fails that (this is the
  case the doc's own edit range would have caused, R-4).

---

## 3. Day plan

One session, ~0.2d, three commits:

| Step | Work | Wall |
|---|---|---|
| 1 | M1 implementation + 19 named subtests; AC1–AC3c | ~0.10d |
| 2 | M2 five mutants, each: mutate → compile-fence → enumerate FAIL set → restore; then AC5/AC5b | ~0.05d (full-package run 115s, dominates) |
| 3 | M3 three prose edits; AC6–AC8 | ~0.05d |

Commit boundaries: M1, M2 (evidence only — the tree returns to M1's state, so M2's commit records
the mutation log; if it produces no file change it folds into M1's record), M3.

---

## 4. Success metrics

- 19 new named subtests across 3 new tests; 5 mutants each with an exact, disjoint enumerated red set.
- `AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -count=1` rc=0 (base rc=0, 114.9s).
- `go build ./...` rc=0, `go vet ./...` rc=0, `go test -count=1 -run '^$' ./host/verifygate/` rc=0.
- `git diff --stat -- go.mod go.sum` empty.
- Three prose sites corrected with measured-base removal needles at 0 and presence controls ≥1.

## 5. Risks

| Risk | Mitigation |
|---|---|
| Implementer writes `srcDeclaresOnKey` in terms of `matchTopLevelOn` | M1 deliverable 3 forbids it explicitly; if it happens, MUT-A's red set grows by `control_quoted_no_runs_on` and AC4's exact-set clause fails — the AC catches the deviation rather than hiding it |
| Flow parser over-reaches into general YAML | `oc-glm-5-2`'s verbatim guard (must start `{` **and** end `}`) is a hard deliverable; `flow_unclosed` is the named subtest that proves refusal, and it is MUT-B's MUST-STAY-GREEN |
| M1 and M3 both touch `dispatch_lever_gate_test.go` | AC3c fences M1 off L108–111; AC6's Go-site base counts are measured against HEAD and re-checkable at any point before M3 |
| `-count=1` omitted and a cached verdict is read | written into the mutation protocol and into AC1/AC4/AC5 |

## 6. Refer to the queue (pre-existing at HEAD — NOT milestones, NOT absorbed)

The doc's §5b says "nothing is filed". Two pre-existing defects surfaced during this planning pass
that the doc did not have; neither is in row 55's scope and neither grows this sprint.

1. **`go build ./...` is used repo-wide as a "the tests still compile" fence, and it cannot be one.**
   Measured first-party (R-1): a type error inside a `_test.go` file leaves `go build ./...` at rc=0.
   `git grep -l 'go build \./\.\.\.'` returns **20+ tracked files**, including `CLAUDE.md`'s hard-rule
   verify gate and many `design_docs/implemented/*sprint-plan.md` mutant fences. `scripts/verify_go.sh`
   itself is NOT affected (it runs `go build ./...` **and** `go test ./... -count=1`, and the latter
   compiles test files) — so this is an **acceptance-criteria instrument-class defect**, not a CI
   hole. Suggested row: audit every AC/doc that uses `go build ./...` as a test-compile fence and
   replace it with `go vet` or `go test -run '^$'`. ~0.1d, clause-2.
2. **The L128 known-positive control is not column-0 anchored, so `runs-on:` satisfies it.** V9,
   re-measured first-party: a workflow with a quoted `"on":` trigger block and any `runs-on:` line
   passes the control on a file whose trigger block the control cannot see. This sprint broadens the
   needle set but does NOT anchor the control at column 0, so the weakness survives. Low value:
   after M1 the helper accepts quoted forms, which removes the only measured consequence. File only
   if the queue wants completeness — flagged, not recommended.
