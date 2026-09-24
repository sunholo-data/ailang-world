# w-workbench-unpinned-render-hunks — pin the five surviving workbench hunks and the fourteen more the full enumeration found next to them

- Status: **implemented** (iteration 185: judged PASS 95/100 + r2 98/100, zero blocking; see §13) · Date: **2026-09-24** · Row **34** (clause-5)
- Iteration: 185, designer role · Measurement base: `origin/dev` = `fd99840` (V0)
- Scope: **test-only.** Charter row 34 lists seven hunks. Re-measured here, two are already killed (H1, H2; §1.1) and five still survive (H3–H7). A rule-3n enumeration of every conditional and boolean hunk on the three surfaces those five sit on (`supportedWorkbenchQuery`, `workbenchHref`, and the grade/verdict template line `render.go:154`) executed **56** compiling mutants. **20 survive** the whole suite: the five charter hunks, 14 more that the suite does not kill today, and 1 equivalent mutant that no test can kill (§2e). This sprint adds tests that kill all 19 killable survivors, each by a named test. It changes no production code: no guard was measured dead or wrong (§2a).
- Revision 1 (quorum r1: astra REJECT upheld — counts; glm REJECT refuted by measurement; gemini PASS): every prose mutant total is corrected to agree with the transcripts and the §7 table. There are 56 mutants, not 55: 36 were killed before this sprint, and 55 non-equivalent mutants must be killed after it. The table and the drill were already right; only the totals in the text were off by one. The ancestry of `9574d08` and the absence of any in-flight overlap are now measured (V20). See §11.
- Revision 2 (quorum r2: gemini-3-1-pro PASS, oc-glm-5-2 PASS, gpt6-astra REJECT with a concrete `proposed_fix`; narrow-refinement carve-out, controller-applied, the reviewer's fix verbatim): `TestRenderGradeWithoutVerdictClaim` now also rejects `class="verdict-` and `aria-label="test verdict` anywhere in the rendered page; new mutation-control row **V10** (an unconditional PASS span appended after the grade paragraph) SURVIVES the r1 prototype and pristine `fd99840`, and is killed by both subtests (V21); I3 is scoped to the constructor-produced grade states the tests cover. Totals are now **57** mutants: **36** killed before, **21** survive before, **56** killed after, Q17 survives. See §11.
- Query grammar: **unchanged** (§4). The accepted-state set is correct; it was only unpinned.
- Estimate: **~0.25 day, ~90 LOC, tests only.** The figure comes from a prototype written, run and deleted at design time: +84 insertions across 2 test files at r1 (**+90** after the r2 negative assertion, measured by the planner: 37 in `workbench_test.go`, 53 in `render_test.go`), `go vet ./...` rc=0, `go test ./...` rc=0, and all 19 targeted mutants killed (V14, V15). Two milestones (§8).

---

## §1 Problem

The workbench is how a human checks provenance in a browser (clause 5). Row 34 records hunks
that no test pins. Such a hunk can be broken, and the whole suite still passes. Each fact below was
measured at `fd99840`, and §10 lists the command for each one.

### 1.1 The charter's seven hunks, re-measured

Instrument: `mutdrive.py` (V3). It applies each mutant as an exact single-occurrence string
replacement. Then it runs `go build ./...`, and a build failure is recorded as a build failure,
never as a kill. Then it runs `go test ./host/workbench ./host/daemon ./host/boundary -count=1`,
the classification arm the row itself uses. It restores the file by `cp` from a backup. After the
run, `git status --porcelain` was empty.

| ID | Anchor at `fd99840` | Mutant | Build | Suite | Status |
|---|---|---|---|---|---|
| H1 | `render.go:130` | `{{if .World.Available}}` → `{{if false}}` | rc=0 | rc=1, red: `TestWorkbenchRendersSeededWorldAndTimeline` | **closed** |
| H2 | `render.go:155` | `{{if .PayloadTruncated}}<p>truncated</p>{{end}}` → `{{if false}}…` | rc=0 | rc=1, red: `TestWorkbenchPayloadPreviewBound`, `TestWorkbenchPayloadPreviewBound/oversize` | **closed** |
| H3 | `render.go:154` | `aria-label="test verdict PASS"` → `aria-label="x"` | rc=0 | rc=0, empty red set | **open** |
| H4 | `render.go:98` | `if query != "" && query[0] == '?' {` → `if query != "" {` | rc=0 | rc=0 | **open** |
| H5 | `workbench.go:70` | `if len(query) != 2 {` → `if false && len(query) != 2 {` | rc=0 | rc=0 | **open** |
| H6 | `workbench.go:73` | `query["from"] != nil && query["entry"] != nil` → `\|\|` | rc=0 | rc=0 | **open** |
| H7 | `workbench.go:76` | `query["object"] != nil && query["payload"] != nil` → `\|\|` | rc=0 | rc=0 | **open** |

**Why H1 and H2 closed, and where their protection is.** The row was queued at iter-112 with
the render package as the test scope ("whole package rc=0"). Two `host/daemon` handler tests
landed the next day, 2026-08-23. Each one renders through the real template [V6]:

- H1 is killed by `TestWorkbenchRendersSeededWorldAndTimeline`, through
  `strings.Contains(body, head.String())` at `workbench_test.go:129`. The world ref reaches the
  page only inside `{{if .World.Available}}`. The assertion came from `5fd6fb3` (WB.C, #85).
- H2 is killed by `TestWorkbenchPayloadPreviewBound/oversize`, through
  `!strings.Contains(body, "truncated")` at `:289-290`. The assertion came from `e563339` (#87).

Both kills are **cross-package only**. Run with `go test ./host/workbench` alone, H1 and H2 both
still **survive** (rc=0) [V5]. This is the same situation the row already records for
`{{if .PayloadShown}}`. The row's own classification arm spans all three packages, so both
count as closed. The cross-package fact goes into §9 so that nobody investigates it again.

**Anchor drift.** The row gives `:69`/`:72`/`:75` and `render.go:101`. At `fd99840` the same
text is at `workbench.go:70`/`:73`/`:76` and `render.go:98` [V19]. The mutants are matched by
text, so the drift does not change any result.

### 1.2 What each open hunk means behaviourally (measured, not read)

I measured each one with a temporary probe test in each package (`zz_probe185_test.go`, run and
then **deleted** [V7–V10]). The daemon probe seeds genesis and one `testCommit`. It then sends
`GET /workbench` for **all 32 subsets** of the five-key vocabulary `{world, object, from, entry,
payload}`, using valid values, and records both `supportedWorkbenchQuery(q)` and the HTTP status.
I ran the probe on the pristine tree and again under each surviving mutant.

- **Pristine.** Exactly 5 of the 32 subsets are accepted: `∅`, `{world}`, `{object}`,
  `{from,entry}`, `{object,payload}`. All 5 return 200. The other **27** return `400` with
  `>unsupported workbench parameter combination<` [V7]. That is §2.2 of the original design,
  exactly.
- **H5** (cardinality guard neutered). **12** extra subsets become `200`, all with 3 or more keys.
  One is `world+object+payload`, the row's own example. Another is `world+from+entry`, which lets
  a timeline page carry a `world` that the timeline ignores [V8].
- **H6** (`from`/`entry` pair `&&`→`||`) opens **6** pairs: `from+world` (the row's
  `?world=W&from=1`), `from+object`, `entry+world`, `entry+object`, `from+payload` and
  `entry+payload` [V8].
- **H7** (`object`/`payload` pair `&&`→`||`) opens **6** pairs: `payload+world` (the row's
  `?world=W&payload=1`), `object+world`, `from+object`, `entry+object`, `from+payload` and
  `entry+payload` [V8].
- **H4.** `workbenchHref("object=abc")` returns `/workbenchobject=abc`, `"//evil.example/x"`
  returns `/workbench//evil.example/x`, and `"javascript:alert(1)"` returns
  `/workbenchjavascript:alert(1)`. In pristine code all three return `/workbench`. Through
  `Render`, a row with `SelectHref: "object=abc"` emits `<a href="/workbenchobject=abc"` where the
  pristine output is `<a href="/workbench"` [V9]. **No production writer reaches this branch.**
  Each of the five `Href`-bearing writers in `host/daemon/workbench.go` (`:141`, `:233`, `:302`,
  `:315`, `:329`) emits `?object=…` or `pageHref` → `?from=…`, and the template's only literal
  argument is `""` (`render.go:129`) [V11]. So the guard is defence in depth for a future
  view-model writer. It is a contract of the function, and it can be observed only by calling
  the function directly.
- **H3.** The PASS span still renders under the mutant, as
  `<span class="verdict-pass" aria-label="x">✓ verdict: PASS</span>`. The reason nothing notices
  is that **no test ever renders a passing verdict.** `t.Run("pass", …)` (`render_test.go:86`)
  checks only struct fields, while the `fail` sibling renders [V10].

### 1.3 The rule-3n enumeration: fifteen more survivors on the same three surfaces (fourteen killable, one equivalent)

The five charter hunks were each found one at a time by different evaluators (iter-112, 119, 121,
123). Rule 3n says to anchor the enumeration to the surface, not to the defect. So I mutated
**every** condition and boolean operand on the three surfaces. For each surface the mutants were:
condition forced to `true`/`false`, relational operator negated, `&&`/`||` swapped, each operand
dropped, each operand negated, and each `return` literal flipped. For the template line it was:
each `{{if}}` forced both ways, the `eq`→`ne` swap, and each attribute/text on the PASS span.
That gave **56 mutants** (34 grammar + 10 href + 10 grade line + H1/H2), **all of which compile.** **20 survive** (V3; §7 lists every row):

- **`supportedWorkbenchQuery`: 34 mutants, 11 survivors.** H5–H7 plus Q18 (`len != 2` → `return true`, which
  opens 16 subsets with 3 or more keys), Q19/Q20 (each operand of the `from`/`entry` pair
  dropped), Q23 (the pair condition forced `true`), Q26/Q27 (each operand of the
  `object`/`payload` pair dropped) and Q30 (`return true`) [V8]. The eleventh, Q17, is equivalent (§2e). The other 23 mutants are killed. All of them either break an *accepted* state or the
  single refused one-key state `?from=0`, and those are the only states the suite exercises.
- **`workbenchHref`: 10 mutants, 3 survivors.** H4, W5 (`if true {`) and W8 (fallback `return "/workbench"`
  → `return "/workbench" + query`). All three have H4's behaviour on a non-`?` input [V9]. The
  other 7 are killed, mostly because `query[0]` panics on the `""` argument that every page
  renders.
- **Grade/verdict line `render.go:154`: 10 mutants, 6 survivors.** H3, V2 (`{{if .Grade.HasVerdict}}` →
  `{{if true}}`), V4 (`{{if eq .Grade.Verdict "FAIL"}}` → `{{if true}}`), V6 (`class="verdict-pass"`
  → `class="verdict-fail"`), V7 (`✓ verdict: {{.Grade.Verdict}}` → `✓ verdict: `) and V9
  (`{{if .Grade.Available}}` → `{{if true}}`). V2 is the dangerous one: a PROVEN grade with **no**
  verdict renders `<span class="verdict-pass" aria-label="test verdict PASS">✓ verdict: </span>`.
  That is a claim of a passing test where none exists, a clause-5 honesty defect. V4 renders a
  PASS as `✗` / `verdict-fail` / `aria-label="test verdict FAIL"`. V9 renders an unavailable grade
  as `<p><span></span></p>`, dropping the `GRADE UNAVAILABLE — reason` text [V10].

**Mechanism.** No test calls `supportedWorkbenchQuery` or `workbenchHref` directly. Both have
**0** `_test.go` references, against 4 and 8 references overall [V13]. The handler tests observe
the grammar only through `TestWorkbenchRefusalBranches`. That test supplies one refused
combination (`?from=0`, which has one key and exits at `:67-68`) and five accepted ones. So the
whole `len(query) ≥ 2` refusal half is unexercised. The row noted this for pairs, and it holds
for triples too.

---

## §2 Decisions

### (a) Test-only. No production change.

The brief allows a production change only where a guard is measured to be dead or wrong. None is:

- Every guard in `supportedWorkbenchQuery` is **live**. Each surviving mutant changes the
  function's output on at least one reachable subset (V8), and pristine output matches the
  grammar exactly (V7).
- `workbenchHref`'s `query[0] == '?'` guard cannot be reached from *today's* production writers
  (V11), but it is not dead code. It is the function's documented behaviour on its whole input
  domain, and the template calls it on every `Href`-bearing field. Removing it would mean any
  future writer that forgets the `?` produces a broken link without any test failing. The
  existing tests cannot see it; the new test pins it.
- On `render.go:154`, every branch is reachable through `NewGradeView`/`NewGradeUnavailable`,
  which are the constructors the render package exports (V10).

### (b) The grammar: an exhaustive truth table at the function, plus three HTTP witnesses.

`supportedWorkbenchQuery` reads only which keys are present and how many there are (`:63-77`).
Unknown and duplicated keys are refused before it is called (`:149-158`). So its reachable domain
is exactly the **32 subsets** of the five-key vocabulary, and an exhaustive table is complete by
construction. It kills every non-equivalent mutant of the function, both the ones enumerated
here and any future edit, rather than only the ones someone thought of. The table has 32
subtests, each named by its sorted key set (`none`, `from+world`, …), so the red set for any
mutant names exactly the combinations it opens. The test also carries a self-check: `seen` must
equal `len(accepted)`, which stops a mistyped accepted-set key from silently turning into an
expected-refused row.

The truth table does not observe the wire. The `!supportedWorkbenchQuery → 400` wiring at `:159`
is already pinned by the existing `unsupported-combination` row (Q2 and Q13 are killed by it).
I still add the three HTTP rows the row's evaluators recommended, as behavioural witnesses in
`TestWorkbenchRefusalBranches`: `world+from`, `world+payload`, and the triple
`world+object+payload`. They keep a refused multi-key query on the HTTP path, so an inlining
refactor that removed the helper would still be caught at the handler. They also serve as second
killers for H5–H7, Q18, Q20, Q23, Q26 and Q30 (§7).

### (c) `workbenchHref`: a direct unit test.

Its survivors are observable only on inputs no production writer emits (§1.2), so the test calls
the function directly: `""`, `?object=abc`, `object=abc`, `//evil.example/x`,
`javascript:alert(1)`. The test is the sole killer of H4, W5 and W8.

### (d) The grade line: render the PASS verdict, the verdict-less grade and the unavailable grade.

- `TestGradeViewRequiresTestVerdict/pass` gains the render arm that its `fail` sibling already
  has. It asserts the **exact** span
  `<p><span>TESTED</span> <span class="verdict-pass" aria-label="test verdict PASS">✓ verdict: PASS</span></p>`
  and asserts that `aria-label="test verdict FAIL"` is absent. It is the sole killer of H3, V4,
  V6 and V7.
- The new `TestRenderGradeWithoutVerdictClaim` has two subtests. `no-verdict` checks that a
  PROVEN grade renders `<p><span>PROVEN</span></p>`, which kills V2, the false-PASS claim.
  `unavailable` checks that it renders `<p>GRADE UNAVAILABLE — no canonical host projection</p>`,
  which kills V9.

### (e) Q17 is equivalent. It is excluded from the kill obligation and recorded.

The Q17 mutant is `if len(query) != 2 {` → `if len(query) > 2 {`. Execution reaches `:70` only
when `len(query) ∉ {0, 1}` (`:64-69` return first), and there `≠ 2 ⇔ > 2`. The daemon probe
under Q17 gives a truth table identical to pristine on all 32 subsets (V16). No test can kill it,
so the drill (AC6) requires it to **survive**, and the prototype's full suite confirmed that it
does (V15).

### (f) Sole killers.

At subtest granularity the new tests are the **only** killers for H3, H4, W5, W8, V2, V4, V6, V7,
V9, V10 (the measured red set is that one test, plus its parent when the killer is a subtest), Q19 and Q27 (truth-table subtests only). For H5–H7,
Q18, Q20, Q23, Q26 and Q30, the truth table is the named killer, and the §2b HTTP witnesses
co-kill some of them. That co-killing is deliberate: it is defence in depth for the wiring, not
redundancy that hides a gap. §7 gives the measured red-set size for every row.

---

## §3 Design

The prototype is below. It was written, run and deleted at design time: `gofmt`-clean, vet
rc=0, `go test ./...` rc=0 (V14). The executor may rename things, but must keep every assertion
string byte-for-byte, because §7's kill results were measured against these exact strings.

### 3.1 `host/daemon/workbench_test.go` (+37)

Add three rows to `TestWorkbenchRefusalBranches`'s table, directly after `unsupported-combination`:

```go
		{"unsupported-pair-world-from", "/workbench?world=" + world + "&from=0", http.StatusBadRequest, "BadRequest", "unsupported workbench parameter combination", nil},
		{"unsupported-pair-world-payload", "/workbench?world=" + world + "&payload=1", http.StatusBadRequest, "BadRequest", "unsupported workbench parameter combination", nil},
		{"unsupported-triple-world-object-payload", "/workbench?world=" + world + "&object=" + object + "&payload=1", http.StatusBadRequest, "BadRequest", "unsupported workbench parameter combination", nil},
```

Add a new top-level test:

```go
// TestSupportedWorkbenchQueryTruthTable walks every subset of the five-key
// vocabulary. The function reads only key presence and count, so these 32
// subsets are its whole reachable domain: exactly five are accepted states.
func TestSupportedWorkbenchQueryTruthTable(t *testing.T) {
	keys := []string{"entry", "from", "object", "payload", "world"}
	accepted := map[string]bool{"none": true, "world": true, "object": true, "entry+from": true, "object+payload": true}
	seen := 0
	for mask := 0; mask < 1<<len(keys); mask++ {
		query := map[string][]string{}
		var names []string
		for i, key := range keys {
			if mask&(1<<i) != 0 {
				query[key] = []string{"x"}
				names = append(names, key)
			}
		}
		name := strings.Join(names, "+")
		if name == "" {
			name = "none"
		}
		if accepted[name] {
			seen++
		}
		t.Run(name, func(t *testing.T) {
			if got := supportedWorkbenchQuery(query); got != accepted[name] {
				t.Fatalf("supportedWorkbenchQuery(%s) = %v, want %v", name, got, accepted[name])
			}
		})
	}
	if seen != len(accepted) {
		t.Fatalf("accepted states enumerated = %d, want %d", seen, len(accepted))
	}
}
```

### 3.2 `host/workbench/render_test.go` (+47)

Append this to `t.Run("pass", …)` inside `TestGradeViewRequiresTestVerdict`, after the struct check:

```go
		var body bytes.Buffer
		if err := Render(&body, Page{Title: "test", Object: &ObjectView{Grade: view}}); err != nil {
			t.Fatal(err)
		}
		want := `<p><span>TESTED</span> <span class="verdict-pass" aria-label="test verdict PASS">✓ verdict: PASS</span></p>`
		if rendered := body.String(); !strings.Contains(rendered, want) || strings.Contains(rendered, `aria-label="test verdict FAIL"`) {
			t.Fatalf("rendered PASS verdict, want %q in %q", want, rendered)
		}
```

Add two new top-level tests:

```go
func TestRenderGradeWithoutVerdictClaim(t *testing.T) {
	proven, err := NewGradeView(GradePROVEN, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name  string
		grade GradeView
		want  string
	}{
		{"no-verdict", proven, `<p><span>PROVEN</span></p>`},
		{"unavailable", NewGradeUnavailable("no canonical host projection"), `<p>GRADE UNAVAILABLE — no canonical host projection</p>`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var body bytes.Buffer
			if err := Render(&body, Page{Title: "test", Object: &ObjectView{Grade: tc.grade}}); err != nil {
				t.Fatal(err)
			}
			rendered := body.String()
			if !strings.Contains(rendered, tc.want) {
				t.Fatalf("rendered grade, want %q in %q", tc.want, rendered)
			}
			for _, claim := range []string{`class="verdict-`, `aria-label="test verdict`} {
				if strings.Contains(rendered, claim) {
					t.Fatalf("grade without a verdict rendered a verdict claim %q in %q", claim, rendered)
				}
			}
		})
	}
}

func TestWorkbenchHrefAppendsOnlyQueryStrings(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"", "/workbench"},
		{"?object=abc", "/workbench?object=abc"},
		{"object=abc", "/workbench"},
		{"//evil.example/x", "/workbench"},
		{"javascript:alert(1)", "/workbench"},
	} {
		if got := workbenchHref(tc.in); got != tc.want {
			t.Errorf("workbenchHref(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
```

### 3.3 Invariants this sprint makes persistent

- **I1.** `supportedWorkbenchQuery` accepts exactly `{∅, {world}, {object}, {from,entry},
  {object,payload}}` over the 32 vocabulary subsets. Any change to the accepted set reds a named
  subtest.
- **I2.** `workbenchHref` returns `"/workbench"+q` exactly when `q` starts with `?`, and returns
  `"/workbench"` otherwise.
- **I3.** Scoped to the constructor-produced grade states the tests cover (`NewGradeView` with a
  PASS verdict, with no verdict, and `NewGradeUnavailable`): the page shows a PASS span only when
  there is a `Verdict == PASS`. For a grade with no verdict and for an unavailable grade, the
  rendered page contains **no** `class="verdict-` and **no** `aria-label="test verdict` anywhere
  (not merely the expected paragraph somewhere), and an unavailable grade shows its reason. The
  page's CSS names `.verdict-pass`/`.verdict-fail` as selectors only, so the `class="` form does not
  match it (V21). This is not claimed as an unrestricted rendering invariant.

---

## §4 Query grammar delta

**None.** `acceptedWorkbenchKeys` (`workbench.go:34-40`) and `supportedWorkbenchQuery`
(`:63-77`) stay byte-identical (AC5). The five accepted states are unchanged. The new tests assert
the grammar that is already there: the 27 refused subsets were already `400` (V7). They were
simply never tested.

---

## §5 Conflict Surface

| Touched | How | Consequence |
|---|---|---|
| `TestWorkbenchRefusalBranches` (`workbench_test.go:134-216`) | +3 table rows | Existing rows are unchanged. The loop body already asserts status, class token, message and security headers for each row. The earlier sprint `w-workbench-timeline-seam` treated this function as byte-identical (its AC5). That constraint was specific to that sprint and has no force here: that sprint has landed (`9574d08` is an ancestor of the base `fd99840`, V20). No in-flight sprint touches `TestWorkbenchRefusalBranches`: there are 0 open PRs, and the only other worktree is `.wt-world-iter182`, which is row 92 and has already landed as `157d6f9` (V20) |
| `TestGradeViewRequiresTestVerdict/pass` (`render_test.go:86-95`) | +8 lines, a render arm | Uses the same `Render`/`bytes.Buffer` idiom as `/fail` (`:68-84`) |
| `TestWorkbenchViewFieldsAllRender` | not touched | No view-model field is added or removed |
| Production `render.go`, `workbench.go` | **not touched** (AC5) | Every row-34 line anchor stays where it is |
| Rows 35/38 (landed at `9574d08`) | none | Their tests and mutation rows are unaffected, and the §7 drill runs the whole three-package arm |
| `.ail`, `/v1` routes, `readStore` | none | — |

Files the implementation changes: `host/daemon/workbench_test.go` and
`host/workbench/render_test.go`. Nothing else.

---

## §6 Acceptance criteria

`export AILANG_BIN=~/.pinned-ailang/ailang` (v0.41.0, V1) applies to every command.

- **AC1** `go vet ./...` → rc=0. (`go build` does not compile `_test.go` files, so it is not an
  acceptable fence here.)
- **AC2** `go test ./... -count=1` → rc=0.
- **AC3** `./scripts/verify_ail.sh` → rc=0, ending `✓ verify gate PASSED: 11 required identities
  verified, 40 named tests pass` (unchanged from V17). `git diff --name-only origin/dev -- '*.ail'
  | wc -l` → `0`.
- **AC4 — the new tests exist and pass.**
  `go test ./host/workbench ./host/daemon -count=1 -v -run 'TestSupportedWorkbenchQueryTruthTable|TestWorkbenchRefusalBranches|TestGradeViewRequiresTestVerdict|TestRenderGradeWithoutVerdictClaim|TestWorkbenchHrefAppendsOnlyQueryStrings'`
  → rc=0 and 0 `--- FAIL` lines, with a `--- PASS` line for each of:
  `TestSupportedWorkbenchQueryTruthTable` and all **32** of its subtests
  (`grep -c -- '--- PASS: TestSupportedWorkbenchQueryTruthTable/'` → `32`);
  `TestWorkbenchRefusalBranches/unsupported-pair-world-from`,
  `…/unsupported-pair-world-payload`, `…/unsupported-triple-world-object-payload`;
  `TestGradeViewRequiresTestVerdict/pass`; `TestRenderGradeWithoutVerdictClaim/no-verdict`,
  `…/unavailable`; `TestWorkbenchHrefAppendsOnlyQueryStrings`. (Prototype: 59 `--- PASS` lines
  for this command, V14.)
- **AC5 — production code is byte-identical.** `git diff --quiet origin/dev --
  host/workbench/render.go host/daemon/workbench.go` → rc=0. `git diff --name-only origin/dev --
  host/` → exactly `host/daemon/workbench_test.go` and `host/workbench/render_test.go`.
- **AC6 — the mutation drill.** Apply every §7 row on its own to the post-sprint tree, using
  the V3 instrument (exact single-occurrence replace, landed count must be 1). For each row
  record: landed y/n, **`go build ./...` rc (must be 0; a build failure is never a kill)**, the
  rc of `go test ./host/workbench ./host/daemon ./host/boundary -count=1`, and the red set.
  Restore by `cp` from a backup, and check that `git status --porcelain` is empty on the committed
  sprint tree before moving to the next row. Pass condition:
  - all **56** non-equivalent rows: build rc=0, suite rc=1, and the named killer from §7 in the
    red set;
  - **Q17**: build rc=0 and suite rc=0. This is the equivalence control; the row must survive;
  - H3–H7, Q18–Q20, Q23, Q26, Q27, Q30, W5, W8, V2, V4, V6, V7, V9 and V10 (the 20 survivors): also
    re-run each against `origin/dev`'s *tests* and confirm it **SURVIVES** there. That shows the
    new tests are what kill it.
- **AC7 — gofmt and size.** `gofmt -l host/daemon/workbench_test.go host/workbench/render_test.go`
  → empty. `git diff --shortstat origin/dev -- host/` → 0 deletions and ≤ 120 insertions.
  (Three other files under `host/` were already gofmt-dirty at `fd99840` [V18]. They are out of
  scope and must not be touched.)

---

## §7 Mutation table

**Every row was executed at design time**, not reasoned. The instrument was V3, the suite was
`./host/workbench ./host/daemon ./host/boundary -count=1`, and each file was restored by `cp`
before the next row. "Before" is the measurement at `fd99840`. "After" is the same mutant with
the §3 prototype in place (V15). The prototype was then deleted. "Named killer" is the test that
must be in the red set (AC6). For rows killed before this sprint, the named killer is one
measured member of today's red set, marked *(existing)*. Red counts include parent tests. `⏎`
marks a newline inside a multi-line match; indentation tabs are elided.

| ID | File | Old → new (exact, single occurrence) | Compiles | Before (fd99840 suite) | Named killer | After (prototype) |
|---|---|---|---|---|---|---|
| H1 | render.go | `{{if .World.Available}}` → `{{if false}}` | y | killed (1 red) | `TestWorkbenchRendersSeededWorldAndTimeline (existing)` | killed (1 red) |
| H2 | render.go | `{{if .PayloadTruncated}}<p>truncated</p>{{end}}` → `{{if false}}<p>truncated</p>{{end}}` | y | killed (2 red) | `TestWorkbenchPayloadPreviewBound/oversize (existing)` | killed (2 red) |
| H3 | render.go | `aria-label="test verdict PASS"` → `aria-label="x"` | y | **SURVIVED** | `TestGradeViewRequiresTestVerdict/pass` | killed (2 red) |
| H4 | render.go | `if query != "" && query[0] == '?' {` → `if query != "" {` | y | **SURVIVED** | `TestWorkbenchHrefAppendsOnlyQueryStrings` | killed (1 red) |
| H5 | workbench.go | `if len(query) != 2 {` → `if false && len(query) != 2 {` | y | **SURVIVED** | `TestSupportedWorkbenchQueryTruthTable/object+payload+world` | killed (15 red) |
| H6 | workbench.go | `if query["from"] != nil && query["entry"] != nil {` → `if query["from"] != nil \|\| query["entry"] != nil {` | y | **SURVIVED** | `TestSupportedWorkbenchQueryTruthTable/from+world` | killed (9 red) |
| H7 | workbench.go | `return query["object"] != nil && query["payload"] != nil` → `return query["object"] != nil \|\| query["payload"] != nil` | y | **SURVIVED** | `TestSupportedWorkbenchQueryTruthTable/payload+world` | killed (9 red) |
| Q1 | workbench.go | `if len(query) == 0 {` → `if false && len(query) == 0 {` | y | killed (21 red) | `TestWorkbenchPagingProbeStoreError/next-probe (existing)` | killed (23 red) |
| Q2 | workbench.go | `if len(query) == 0 {` → `if true \|\| len(query) == 0 {` | y | killed (2 red) | `TestWorkbenchRefusalBranches/unsupported-combination (existing)` | killed (33 red) |
| Q3 | workbench.go | `if len(query) == 0 {` → `if len(query) != 0 {` | y | killed (22 red) | `TestWorkbenchPagingProbeStoreError/next-probe (existing)` | killed (54 red) |
| Q4 | workbench.go | `if len(query) == 0 {⏎return true` → `if len(query) == 0 {⏎return false` | y | killed (21 red) | `TestWorkbenchPagingProbeStoreError/next-probe (existing)` | killed (23 red) |
| Q5 | workbench.go | `if len(query) == 1 {` → `if false && len(query) == 1 {` | y | killed (10 red) | `TestWorkbenchPayloadPreviewBound/default-off (existing)` | killed (13 red) |
| Q6 | workbench.go | `if len(query) == 1 {` → `if true \|\| len(query) == 1 {` | y | killed (17 red) | `TestWorkbenchPagingProbeStoreError/prev-probe (existing)` | killed (43 red) |
| Q7 | workbench.go | `if len(query) == 1 {` → `if len(query) != 1 {` | y | killed (23 red) | `TestWorkbenchPagingProbeStoreError/prev-probe (existing)` | killed (51 red) |
| Q8 | workbench.go | `return query["world"] != nil \|\| query["object"] != nil` → `return query["world"] != nil && query["object"] != nil` | y | killed (10 red) | `TestWorkbenchPayloadPreviewBound/default-off (existing)` | killed (13 red) |
| Q9 | workbench.go | `return query["world"] != nil \|\| query["object"] != nil` → `return query["object"] != nil` | y | killed (4 red) | `TestWorkbenchRefusalBranches/absent-world (existing)` | killed (6 red) |
| Q10 | workbench.go | `return query["world"] != nil \|\| query["object"] != nil` → `return query["world"] != nil` | y | killed (8 red) | `TestWorkbenchPayloadPreviewBound/default-off (existing)` | killed (10 red) |
| Q11 | workbench.go | `return query["world"] != nil \|\| query["object"] != nil` → `return query["world"] == nil \|\| query["object"] != nil` | y | killed (5 red) | `TestWorkbenchRefusalBranches/absent-world (existing)` | killed (10 red) |
| Q12 | workbench.go | `return query["world"] != nil \|\| query["object"] != nil` → `return query["world"] != nil \|\| query["object"] == nil` | y | killed (9 red) | `TestWorkbenchPayloadPreviewBound/default-off (existing)` | killed (14 red) |
| Q13 | workbench.go | `return query["world"] != nil \|\| query["object"] != nil` → `return true` | y | killed (2 red) | `TestWorkbenchRefusalBranches/unsupported-combination (existing)` | killed (6 red) |
| Q14 | workbench.go | `return query["world"] != nil \|\| query["object"] != nil` → `return false` | y | killed (10 red) | `TestWorkbenchPayloadPreviewBound/default-off (existing)` | killed (13 red) |
| Q15 | workbench.go | `if len(query) != 2 {` → `if true \|\| len(query) != 2 {` | y | killed (22 red) | `TestWorkbenchPagingProbeStoreError/prev-probe (existing)` | killed (25 red) |
| Q16 | workbench.go | `if len(query) != 2 {` → `if len(query) == 2 {` | y | killed (22 red) | `TestWorkbenchPagingProbeStoreError/prev-probe (existing)` | killed (38 red) |
| Q17 | workbench.go | `if len(query) != 2 {` → `if len(query) > 2 {` | y | **SURVIVED** | `— (equivalent, §2e)` | SURVIVED |
| Q18 | workbench.go | `if len(query) != 2 {⏎return false` → `if len(query) != 2 {⏎return true` | y | **SURVIVED** | `TestSupportedWorkbenchQueryTruthTable/entry+object+world` | killed (19 red) |
| Q19 | workbench.go | `if query["from"] != nil && query["entry"] != nil {` → `if query["entry"] != nil {` | y | **SURVIVED** | `TestSupportedWorkbenchQueryTruthTable/entry+world` | killed (4 red) |
| Q20 | workbench.go | `if query["from"] != nil && query["entry"] != nil {` → `if query["from"] != nil {` | y | **SURVIVED** | `TestSupportedWorkbenchQueryTruthTable/from+world` | killed (6 red) |
| Q21 | workbench.go | `if query["from"] != nil && query["entry"] != nil {` → `if query["from"] == nil && query["entry"] != nil {` | y | killed (17 red) | `TestWorkbenchPagingProbeStoreError/prev-probe (existing)` | killed (22 red) |
| Q22 | workbench.go | `if query["from"] != nil && query["entry"] != nil {` → `if query["from"] != nil && query["entry"] == nil {` | y | killed (17 red) | `TestWorkbenchPagingProbeStoreError/prev-probe (existing)` | killed (23 red) |
| Q23 | workbench.go | `if query["from"] != nil && query["entry"] != nil {` → `if true {` | y | **SURVIVED** | `TestSupportedWorkbenchQueryTruthTable/object+world` | killed (12 red) |
| Q24 | workbench.go | `if query["from"] != nil && query["entry"] != nil {` → `if false {` | y | killed (17 red) | `TestWorkbenchPagingProbeStoreError/prev-probe (existing)` | killed (19 red) |
| Q25 | workbench.go | `query["entry"] != nil {⏎return true` → `query["entry"] != nil {⏎return false` | y | killed (17 red) | `TestWorkbenchPagingProbeStoreError/prev-probe (existing)` | killed (19 red) |
| Q26 | workbench.go | `return query["object"] != nil && query["payload"] != nil` → `return query["payload"] != nil` | y | **SURVIVED** | `TestSupportedWorkbenchQueryTruthTable/payload+world` | killed (6 red) |
| Q27 | workbench.go | `return query["object"] != nil && query["payload"] != nil` → `return query["object"] != nil` | y | **SURVIVED** | `TestSupportedWorkbenchQueryTruthTable/object+world` | killed (4 red) |
| Q28 | workbench.go | `return query["object"] != nil && query["payload"] != nil` → `return query["object"] == nil && query["payload"] != nil` | y | killed (7 red) | `TestWorkbenchPayloadPreviewBound/opt-in (existing)` | killed (13 red) |
| Q29 | workbench.go | `return query["object"] != nil && query["payload"] != nil` → `return query["object"] != nil && query["payload"] == nil` | y | killed (7 red) | `TestWorkbenchPayloadPreviewBound/opt-in (existing)` | killed (12 red) |
| Q30 | workbench.go | `return query["object"] != nil && query["payload"] != nil` → `return true` | y | **SURVIVED** | `TestSupportedWorkbenchQueryTruthTable/object+world` | killed (12 red) |
| Q31 | workbench.go | `return query["object"] != nil && query["payload"] != nil` → `return false` | y | killed (7 red) | `TestWorkbenchPayloadPreviewBound/opt-in (existing)` | killed (9 red) |
| W1 | render.go | `if query != "" && query[0] == '?' {` → `if query[0] == '?' {` | y | killed (27 red) | `TestGradeViewRequiresTestVerdict/fail (existing)` | killed (23 red) |
| W2 | render.go | `if query != "" && query[0] == '?' {` → `if query != "" \|\| query[0] == '?' {` | y | killed (27 red) | `TestGradeViewRequiresTestVerdict/fail (existing)` | killed (23 red) |
| W3 | render.go | `if query != "" && query[0] == '?' {` → `if query != "" && query[0] != '?' {` | y | killed (15 red) | `TestRenderTimelinePagingLinks/next (existing)` | killed (16 red) |
| W4 | render.go | `if query != "" && query[0] == '?' {` → `if query == "" && query[0] == '?' {` | y | killed (27 red) | `TestGradeViewRequiresTestVerdict/fail (existing)` | killed (23 red) |
| W5 | render.go | `if query != "" && query[0] == '?' {` → `if true {` | y | **SURVIVED** | `TestWorkbenchHrefAppendsOnlyQueryStrings` | killed (1 red) |
| W6 | render.go | `if query != "" && query[0] == '?' {` → `if false {` | y | killed (15 red) | `TestRenderTimelinePagingLinks/next (existing)` | killed (16 red) |
| W7 | render.go | `return "/workbench" + query⏎` → `return "/workbench"⏎` | y | killed (15 red) | `TestRenderTimelinePagingLinks/next (existing)` | killed (16 red) |
| W8 | render.go | `}⏎return "/workbench"⏎}` → `}⏎return "/workbench" + query⏎}` | y | **SURVIVED** | `TestWorkbenchHrefAppendsOnlyQueryStrings` | killed (1 red) |
| W9 | render.go | `}⏎return "/workbench"⏎}` → `}⏎return ""⏎}` | y | killed (1 red) | `TestRenderEmitsOnlyLocalLinks (existing)` | killed (2 red) |
| V1 | render.go | `{{if .Grade.HasVerdict}}` → `{{if false}}` | y | killed (2 red) | `TestGradeViewRequiresTestVerdict/fail (existing)` | killed (3 red) |
| V2 | render.go | `{{if .Grade.HasVerdict}}` → `{{if true}}` | y | **SURVIVED** | `TestRenderGradeWithoutVerdictClaim/no-verdict` | killed (2 red) |
| V3 | render.go | `{{if eq .Grade.Verdict "FAIL"}}` → `{{if false}}` | y | killed (2 red) | `TestGradeViewRequiresTestVerdict/fail (existing)` | killed (2 red) |
| V4 | render.go | `{{if eq .Grade.Verdict "FAIL"}}` → `{{if true}}` | y | **SURVIVED** | `TestGradeViewRequiresTestVerdict/pass` | killed (2 red) |
| V5 | render.go | `{{if eq .Grade.Verdict "FAIL"}}` → `{{if ne .Grade.Verdict "FAIL"}}` | y | killed (2 red) | `TestGradeViewRequiresTestVerdict/fail (existing)` | killed (3 red) |
| V6 | render.go | `class="verdict-pass"` → `class="verdict-fail"` | y | **SURVIVED** | `TestGradeViewRequiresTestVerdict/pass` | killed (2 red) |
| V7 | render.go | `✓ verdict: {{.Grade.Verdict}}` → `✓ verdict: ` | y | **SURVIVED** | `TestGradeViewRequiresTestVerdict/pass` | killed (2 red) |
| V8 | render.go | `{{if .Grade.Available}}` → `{{if false}}` | y | killed (2 red) | `TestGradeViewRequiresTestVerdict/fail (existing)` | killed (5 red) |
| V9 | render.go | `{{if .Grade.Available}}` → `{{if true}}` | y | **SURVIVED** | `TestRenderGradeWithoutVerdictClaim/unavailable` | killed (2 red) |
| V10 | render.go | `{{.Grade.Unavailable}}</p>{{end}}` → `{{.Grade.Unavailable}}</p>{{end}}<span class="verdict-pass" aria-label="test verdict PASS">✓ verdict: PASS</span>` (r2 mutation control, astra) | y | **SURVIVED** (also SURVIVES the r1 prototype) | `TestRenderGradeWithoutVerdictClaim/no-verdict` and `/unavailable` (both required) | killed (3 red) |

**Tally.** 56 mutants, all of which compile. The count by prefix is H7 + Q31 + W9 + V9 = 56, measured with `grep -cE '^[HQWV][0-9]+ ' drill_before.txt` → `56` (the same command on `drill_after.txt` → `56`; per prefix, `grep -cE "^H[0-9]+ "` etc. → 7/31/9/9). **Before:** 36 killed and 20 survive (H3–H7, Q17–Q20,
Q23, Q26, Q27, Q30, W5, W8, V2, V4, V6, V7, V9). **After:** 55 killed and 1 survives (Q17,
equivalent, §2e). Every row that was killed before is still killed after. Each of the 19 rows that
the sprint turns from SURVIVED into killed has its named killer in the measured after-state red
set. The generator asserted this membership when it built the table (V15).

**Tally after revision 2 (supersedes the counts above for AC6).** **57** mutants: the 56 above plus
V10. Before (pristine `fd99840` tests): 36 killed, **21** survive (the 20 above plus V10, measured by
`v10_control.sh` against the pristine suite: build rc=0, suite rc=0). V10 also **survives the r1
prototype** (build rc=0, suite rc=0), so the r2 negative assertion is what kills it. After (the r2
prototype, `prototype_r2_controller.diff`): the whole drill was re-run by the controller with V10
added (`r2/mutdrive_r2.py` → `r2/drill_after_r2.txt`): **56 KILLED, 1 SURVIVED (Q17), 0
`compiles=n`, 0 NOT-LANDED** over 57 rows; V10's red set is `TestRenderGradeWithoutVerdictClaim`,
`/no-verdict`, `/unavailable` (V21).

---

## §8 Milestones

Each milestone ends green on AC1 and AC2.

| M | Content | ~LOC | Exit |
|---|---|---|---|
| **M1 — the grammar** | `workbench_test.go`: `TestSupportedWorkbenchQueryTruthTable` and the three `TestWorkbenchRefusalBranches` rows (§3.1) | ~37 | AC1, AC2; H5–H7, Q17–Q20, Q23, Q26, Q27, Q30 drilled |
| **M2 — href and grade line** | `render_test.go`: the `/pass` render arm, `TestRenderGradeWithoutVerdictClaim`, `TestWorkbenchHrefAppendsOnlyQueryStrings` (§3.2) | ~53 (measured, r2) | AC1–AC7; all 57 rows drilled |

---

## §9 Out of scope

These are recorded so that nobody widens this sprint to cover them. Each "row candidate" was
measured, and it is for the controller to file or drop it. This doc does not file them.

- **R-a (row candidate) — production never supplies a grade, so every object page renders
  `GRADE UNAVAILABLE — ` with an empty reason.** `host/daemon` has **0** non-test references to
  `Grade`. The same-file control is `PayloadTruncated` = **1** in `workbench.go`. The daemon
  probe's `GET /workbench?object=<stored>` returned 200, and its grade text was
  `GRADE UNAVAILABLE — </p>`, with no reason after the dash [V12]. This is the view-model-field
  class again (rows 35/38), and it is the honest-UNAVAILABLE half of clause 5: the page never
  says *why* the grade is unavailable. This sprint pins the template's behaviour for a
  *supplied* reason (V9's killer). It does not change the handler.
- **Observation, not a row — H1 and H2 are protected only across packages.** `go test
  ./host/workbench` on its own lets both survive [V5]. They are killed through `host/daemon`
  (§1.1), exactly like the `{{if .PayloadShown}}` non-gap that row 34 already records. The row's
  classification arm spans both packages, so this is not a gap under its own definition. If the
  controller wants per-package self-sufficiency as a rule, that is a protocol decision, not a
  test for this sprint.
- **Observation — three `host/` files are gofmt-dirty at `fd99840`**
  (`host/daemon/session_middleware_test.go`, `host/store/schema_version_test.go`,
  `host/store/store.go`) [V18]. None of them is touched here. A CI gofmt gate would be a separate
  item.
- The other mutable hunks in `handleWorkbench` (`:146-332`) and in the rest of `pageHTML`. The
  brief limits the rule-3n sweep to the three surfaces row 34 names. Rows 35/38's drill
  (`w-workbench-timeline-seam` §7) covered the paging and selection hunks.
- Any `.ail` change, any `/v1` route change, and any change to production workbench code.

---

## §10 Verification Log

All commands were run in this worktree at `fd99840` on 2026-09-24, with `AILANG_BIN=~/.pinned-ailang/ailang`.
Evidence files are under `~/.ailang/state/mission-world-iter185-evidence/designer/`: the
instrument `mutdrive.py`, the drill transcripts `drill_before.txt` / `drill_after.txt` /
`drill_h1h2_renderonly.txt`, the probe transcripts `probe_pristine.txt` / `probe_mutants.txt`,
the prototype `prototype.diff`, and the generated `table7.md`. Every zero below has a
same-command positive control.

| # | Claim | Command | Result |
|---|---|---|---|
| V0 | Base is `origin/dev` `fd99840` | `git rev-parse HEAD origin/dev` | both `fd998400a1cc…` |
| V1 | Pinned binary | `$AILANG_BIN --version` | `AILANG v0.41.0` (commit `24ee108`) |
| V2 | The base is green on the drill arm | `go build ./... && go test ./host/workbench ./host/daemon ./host/boundary -count=1` | build rc=0; `ok` ×3 |
| V3 | Before-drill: 56 mutants, all compile, 36 killed, 20 survive | `python3 mutdrive.py` (for each row: exact single-occurrence replace → `go build ./...` → the 3-package suite → `cp` restore) → `drill_before.txt` | every row `compiles=y`; SURVIVED: H3 H4 H5 H6 H7 Q17 Q18 Q19 Q20 Q23 Q26 Q27 Q30 W5 W8 V2 V4 V6 V7 V9; every other row `rc=1 KILLED`; final `git status --porcelain` = `''` |
| V4 | Red sets for H1/H2 | same run | H1 `rc=1 reds=1 TestWorkbenchRendersSeededWorldAndTimeline`; H2 `rc=1 reds=2 TestWorkbenchPayloadPreviewBound TestWorkbenchPayloadPreviewBound/oversize` |
| V5 | H1/H2 survive the render package on its own | `TESTCMD="go test ./host/workbench -count=1" python3 mutdrive.py H1 H2` | `H1 compiles=y rc=0 SURVIVED`, `H2 compiles=y rc=0 SURVIVED`. Control: the same mutants are killed by the 3-package arm (V4) |
| V6 | Where H1/H2's killers come from | `grep -n 'head.String()' host/daemon/workbench_test.go`; `git log -S'selected head world ref'` / `-S'oversize payload not marked truncated' -- host/daemon/workbench_test.go` | `:129-130`; the oversize check at `:289-290` in `t.Run("oversize")` (`:279`); introduced by `5fd6fb3` (WB.C, #85) and `e563339` (#87), both 2026-08-23 |
| V7 | Pristine grammar over all 32 subsets | throwaway `host/daemon/zz_probe185_test.go` (seeds genesis + 1 `testCommit`, sends GET for each subset with valid values, logs `supportedWorkbenchQuery` + status + message), run with `-run TestZZProbe185 -v`, then **deleted** | `supported=true status=200` for exactly `∅`, `world`, `object`, `entry+from`, `object+payload`; the other 27 give `supported=false status=400 unsupportedMsg=true` (`probe_pristine.txt`) |
| V8 | Subsets each survivor opens | the V7 probe under each mutant (`SHOWPROBE=1 TESTCMD=… mutdrive.py H3 … V9`) | H5: +12 (all ≥3-key, incl. `object+payload+world`, `entry+from+world`); H6: +`from+world`, `from+object`, `entry+world`, `entry+object`, `from+payload`, `entry+payload`; H7: +`object+world`, `from+object`, `entry+object`, `payload+world`, `from+payload`, `entry+payload`; Q18: +16 (every ≥3-key subset); Q19: +`entry+world`, `entry+object`, `entry+payload`; Q20: +`from+world`, `from+object`, `from+payload`; Q23 and Q30: +all 8 refused pairs; Q26: +`payload+world`, `from+payload`, `entry+payload`; Q27: +`object+world`, `from+object`, `entry+object` (`probe_mutants.txt`) |
| V9 | `workbenchHref` behaviour, pristine vs H4/W5/W8 | throwaway `host/workbench/zz_probe185_test.go` (deleted) | pristine: `""`→`/workbench`, `?object=abc`→`/workbench?object=abc`, `object=abc`/`//evil.example/x`/`javascript:alert(1)`→`/workbench`; `Render` of `SelectHref:"object=abc"` → `<a href="/workbench"`. Under H4, W5 and W8: `/workbenchobject=abc`, `/workbench//evil.example/x`, `/workbenchjavascript:alert(1)`, `<a href="/workbenchobject=abc"` |
| V10 | Grade-line rendering, pristine vs mutants | the V9 probe renders `NewGradeView(TESTED, TestReport, PASS/FAIL)`, `NewGradeView(PROVEN,"",nil)` and `NewGradeUnavailable("no canonical host projection")` | pristine PASS `<p><span>TESTED</span> <span class="verdict-pass" aria-label="test verdict PASS">✓ verdict: PASS</span></p>`; PROVEN `<p><span>PROVEN</span></p>`; unavailable `<p>GRADE UNAVAILABLE — no canonical host projection</p>`. H3: `aria-label="x"`; V2: PROVEN → `… <span class="verdict-pass" aria-label="test verdict PASS">✓ verdict: </span>`; V4: PASS → `class="verdict-fail" aria-label="test verdict FAIL">✗ verdict: PASS`; V6: `class="verdict-fail" aria-label="test verdict PASS"`; V7: `✓ verdict: </span>`; V9: unavailable → `<p><span></span></p>`. Control: the FAIL render is identical to pristine under every one of these mutants. `render_test.go:86-95` (`/pass`) calls no `Render` |
| V11 | No production writer passes a non-`?` non-empty href | `grep -rn 'workbenchHref\|Href:\|Href =\|SelectHref =' host --include='*.go' \| grep -v _test.go` | writers `workbench.go:141` `"?object="+target`, `:233` `"?object="+…`, `:302`/`:315`/`:329` `pageHref(…)` (`"?from="…`, `:114-116`); template literal arg only `workbenchHref ""` at `render.go:129`. The positive control is these 5 writer hits plus 8 `render.go` hits (the definition, 6 template calls and the FuncMap entry) |
| V12 | The daemon never sets `Grade`; the production grade text has no reason | `grep -rn Grade host/daemon --include='*.go' \| grep -v _test \| wc -l`; control `grep -c PayloadTruncated host/daemon/workbench.go`; V7 probe `GET /workbench?object=<stored>` | `0`; control `1`; status 200, text `"GRADE UNAVAILABLE — </p>\n…"` |
| V13 | No test calls the two functions directly | `grep -rn supportedWorkbenchQuery host --include='*_test.go' \| wc -l`; same with `--include='*.go'`; the same pair for `workbenchHref(` / `workbenchHref` | `0` / `4`; `0` / `8` |
| V14 | The prototype compiles, is clean and green; LOC | apply §3 → `gofmt -l host/` (only V18's three pre-existing files listed), `go vet ./...`, `go test ./... -count=1`, the AC4 `-run` command, `git diff --numstat`; saved `prototype.diff`; restored by `cp` | vet rc=0; test rc=0 (no non-`ok` lines); **59** `--- PASS` lines; numstat `37 0 host/daemon/workbench_test.go`, `47 0 host/workbench/render_test.go`; `git status --porcelain` empty after restore |
| V15 | After-drill with the prototype in place | `python3 mutdrive.py` → `drill_after.txt`; `table7.md` generated from both drills, with an `assert` that each named killer is in the after red set | 55 `KILLED`, 1 `SURVIVED` (Q17), 0 `compiles=n` over 56 rows; the asserts passed (table generated) |
| V16 | Q17 is equivalent | V8 probe under Q17 | accepted set identical to pristine (`∅`, `world`, `object`, `entry+from`, `object+payload`) on all 32 subsets; `:64-69` return for len 0/1 (`sed -n 63,77p host/daemon/workbench.go`) |
| V17 | `.ail` gate baseline | `./scripts/verify_ail.sh` | rc=0; `✓ verify gate PASSED: 11 required identities verified, 40 named tests pass` |
| V18 | Pre-existing gofmt debt | `gofmt -l host/` on the pristine tree | `host/daemon/session_middleware_test.go`, `host/store/schema_version_test.go`, `host/store/store.go` (none of them in this sprint's file set) |
| V19 | Anchor drift | `sed -n 63,77p host/daemon/workbench.go`; `cat -n host/workbench/render.go \| sed -n '97,102p;130p;154,155p'` | `if len(query) != 2 {` `:70`, from/entry pair `:73`, object/payload return `:76`; `workbenchHref` guard `render.go:98`; `{{if .World.Available}}` `:130`; grade line `:154`; payload line `:155` |
| V20 | (r1, glm) `9574d08` has landed in the base, and nothing in flight overlaps | `git merge-base --is-ancestor 9574d08 fd99840 && echo yes`; control `git merge-base --is-ancestor fd99840 9574d08; echo $?`; `gh pr list --repo sunholo-data/ailang-world --state open --json number \| jq length`; `git worktree list` (main checkout) | `yes`; control rc=`1`, so the check discriminates; `0` open PRs; the worktrees are only the main checkout (`fd99840 [dev]`), `.wt-world-iter182` (`sprint/w-prove-1-0-phase-a`, row 92, landed `157d6f9`) and this sprint's `.wt-world-iter185`. Measured by the controller at about 17:55Z on 2026-09-24 and reproduced by the designer |
| V21 | (r2, astra) the no-verdict/unavailable subtests reject verdict claims; V10 control | evidence dir `~/.ailang/state/mission-world-iter185-evidence/`: `v10_control.sh` run against (a) the r2 prototype, (b) the r1 prototype, (c) pristine `fd99840`; then `r2/mutdrive_r2.py` (57 rows) against the r2 prototype; `grep -c 'verdict-' host/workbench/render.go` | (a) build rc=0, suite rc=1, red = `TestRenderGradeWithoutVerdictClaim` + `/no-verdict` + `/unavailable`; (b) build rc=0, suite rc=0 (SURVIVES); (c) build rc=0, suite rc=0 (SURVIVES); full drill 56 KILLED / 1 SURVIVED (Q17) / 0 compiles=n; `verdict-` occurs on 2 lines: the CSS selectors at `:122` (`.verdict-fail{`, `.verdict-pass`) and the spans at `:154` — the CSS carries no `class="` prefix, so the negative assertion cannot self-match. Measured by the controller, 2026-09-24 |

---

## §11 Quorum log

| Round | Reviewer | Verdict | Objection (one line) | Disposition |
|---|---|---|---|---|
| r1 | gpt6-astra | reject | §7 has 56 rows (H7+Q31+W9+V9), yet V3/V15 claim 55 mutants, and AC6 requires only 54 non-equivalent kills plus Q17 | **Upheld.** The controller measured the designer's own transcripts: `drill_before.txt` has 56 rows (36 KILLED, 20 SURVIVED, 0 `compiles=n`), and `drill_after.txt` has 55 KILLED and 1 SURVIVED (Q17). The table and the drill were right; the prose totals were off by one. The header, §1.3, AC6, the §7 Tally (with the counting command), §8 and V3/V15 are corrected. The 20-survivor list is unchanged: 5 charter hunks + 14 new + Q17 |
| r1 | oc-glm-5-2 | reject | No verification row shows that `9574d08` is an ancestor of `fd99840`, and overlap with in-flight work is not checked | **Refuted by measurement (V20).** `9574d08` is an ancestor of `fd99840` (the reverse check has rc=1, so the check discriminates). There are 0 open PRs, and the only other worktree is row 92, already landed as `157d6f9`. §5 now cites V20. No design change |
| r1 | gemini-3-1-pro | pass | — | — |
| r2 | gpt6-astra | reject | `TestRenderGradeWithoutVerdictClaim` only checks the expected paragraph occurs somewhere; a renderer that keeps it and appends a false PASS span passes both subtests, so killing V2 does not establish I3 | **Upheld and applied under the narrow-refinement carve-out** (the objection carries a concrete `proposed_fix` and does not dispute the design direction). Applied verbatim: both subtests now reject `class="verdict-` and `aria-label="test verdict` anywhere in the page; mutation-control row V10 added and measured (survives pristine and the r1 prototype, killed by both subtests); drill totals updated (57 / 36 / 21 / 56); I3 scoped to the constructor-produced grade states. Controller-applied, V21 |
| r2 | gemini-3-1-pro | pass | (non-blocking) §3.2's `pass` arm uses `view` without showing its definition | Noted for the executor: `view` is the variable the existing `t.Run("pass", …)` body already binds (the prototype compiled, V14) |
| r2 | oc-glm-5-2 | pass | (non-blocking) §4 cites `acceptedWorkbenchKeys` at `workbench.go:34-40` without a V-row | Noted for the executor: re-anchor at execution time |

---

## §12 Execution record

Executed 2026-09-24 by mission-control iteration 185 (sprint-executor) in
`/Users/voightkampff/dev/sunholo-data/.wt-world-iter185`, branch
`sprint/w-workbench-unpinned-render-hunks`, from `8d99105` (the plan's `aa5df54` plus the
planner-findings doc commit; `git diff fd99840 8d99105 -- host cmd scripts` is empty). `origin/dev`
= `fd99840` throughout (re-fetched before the survival check). `AILANG_BIN=~/.pinned-ailang/ailang`
(`AILANG v0.41.0`, commit `24ee108`), `go1.26.6 darwin/arm64`. The test code was applied from
`prototype_r2_controller.diff` per file (`git apply --include=<file>`), so it is byte-identical to
the plan's §3.2/§4.2 blocks. Evidence: `~/.ailang/state/mission-world-iter185-evidence/executor/`
(`drill.py` = the planner's runner, sha256 `a4c00ade…fcfa4c5`, unedited; `m1_gates.txt`,
`m2_gates.txt`, `m1_ac.txt`, `m2_ac4.txt`, `drill_m1.txt`, `drill_m2.txt`, `drill_final.txt`,
`drill_origin_dev.txt`, `verify_ail.txt`, `ci_list.txt` and the per-step `ci_*.txt`).

| Milestone | Commit | Files |
|---|---|---|
| M1 — the grammar | `133a83f` | `host/daemon/workbench_test.go` (+37/−0) |
| M2 — href and grade line | `5ee393b` | `host/workbench/render_test.go` (+53/−0) |

**Gates.**

| Gate | M1 | M2 (final tree) |
|---|---|---|
| AC1 `go vet ./...` | rc=0 | rc=0 |
| AC2 `go test ./... -count=1` | rc=0; 20 `ok`, 0 `FAIL` | rc=0; 20 `ok`, 0 `FAIL` |
| `gofmt -l` (the two test files) | empty | empty |
| `go test -race ./host/workbench ./host/daemon -count=1` | both `ok` | both `ok` |
| M1 exit (`-run 'TestSupportedWorkbenchQueryTruthTable\|TestWorkbenchRefusalBranches'`) | rc=0; 51 PASS / 32 / 17 / 0 FAIL / 3 new refusal rows | — |
| AC3 `./scripts/verify_ail.sh` | — | rc=0; `✓ world package gate PASSED: 9/9 steps performed non-zero work`; `✓ verify gate PASSED: 11 required identities verified, 40 named tests pass`; `.ail` files changed vs `origin/dev`: 0 |
| AC4 (the §6 `-run` command) | — | rc=0; **59** `--- PASS`, 0 `--- FAIL`; 32 truth-table subtests; one PASS line each for the three new refusal rows, `/pass`, `/no-verdict`, `/unavailable`, `TestWorkbenchHrefAppendsOnlyQueryStrings` |
| AC5 `git diff --quiet origin/dev -- host/workbench/render.go host/daemon/workbench.go` | rc=0 | rc=0 (also after every drill); `--name-only origin/dev -- host/` = the two test files |
| AC7 size | `1 file changed, 37 insertions(+)` | `git diff --shortstat fd99840 -- host/` → **`2 files changed, 90 insertions(+)`** (numstat 37/0, 53/0) |
| `bash scripts/check_no_personal_email.sh` | — | ✓ |

**Mutation drills** (plan §2.3 runner; every row: exact single-occurrence replace, `go build ./...`
before the suite, `go test ./host/workbench ./host/daemon ./host/boundary -count=1`, `cp` restore).

| Drill | Tree | Rows | KILLED (named killer in reds) | SURVIVED | `compiles=n` / NOT-LANDED / `killer_in_reds=n` | Final `STATUS` |
|---|---|---|---|---|---|---|
| M1 | `133a83f` | H5 H6 H7 Q17 Q18 Q19 Q20 Q23 Q26 Q27 Q30 | 10 (reds 15/9/9/19/4/6/12/6/4/12) | 1 (Q17, rc=0) | 0 | `''` |
| M2 | `5ee393b` | H3 H4 W5 W8 V2 V4 V6 V7 V9 V10 | 10 (reds 2/1/1/1/2/2/2/2/2/3) | 0 | 0 | `''` |
| full | `5ee393b` | all 57 | **56** | **1** (Q17) | 0 | `''` |
| survival on `origin/dev`'s tests | detached sibling worktree `.exec-base-iter185` at `fd99840` (removed afterwards) | the 20 (H3–H7, Q18–Q20, Q23, Q26, Q27, Q30, W5, W8, V2, V4, V6, V7, V9, V10) | 0 | **20** (`compiles=y rc=0 SURVIVED`) | 0 | `''` |

The full drill's per-row red counts equal the §7 "After" column on all 57 rows (compared by
script: 57 parsed from each side, 0 differences). The full drill took about 3 min (20:18:02 →
20:21:05, load average ≈ 5.5).

**CI command list (plan §1.3), run locally on the final tree.** `ailang --version` → v0.41.0;
`verify_ail.sh` rc=0 (above); `go version` → go1.26.6; `go build ./...` rc=0;
`w-race-gate-blindspot/run.sh` rc=0 (affected toolchain BUG, known-good OK, pinned OK);
`bench_worldd.sh --smoke` rc=0 (PASSED); `--check-claims` rc=0 (PASSED);
`test_gate1_range_check.sh` 18/0; `test_queue_census.sh` 68/0; `queue_census.sh --doc
design_docs/world-mission.md --control-closed 1 --control-open 79` rc=0 (both controls ok);
`test_gate0_self_notices.sh` 88/0; `check_no_personal_email.sh` ✓; `test_check_no_personal_email.sh`
10/0. `verify_go.sh` was not run locally (standing rig rule, plan §1.3); its build/test/race
content is covered by the vet, whole-repo test and touched-package race legs above.

**Deviations.**

1. *Recorder-refusal arm not reproducible on this rig.* `bench_worldd.sh --record-pair --variant .
   --control .` did fail as required (rc=1), but its stderr is `✗ record-pair REFUSED: control commit
   is not the variant parent`, not `probe FAILED: sysctl -n hw.ncpu`. Why: the arm asserts the
   *off-rig* (Linux runner) refusal; on this macOS host `sysctl -n hw.ncpu` succeeds, so the recorder
   refuses at a later check. Evidence: `ci_rec_err.txt`; the step's own CI text calls it "the
   off-rig runner"; `git diff --quiet fd99840 -- scripts .github` → rc=0 (this sprint changes neither
   the script nor the workflow). CI is the binding gate for that step.
2. *Evidence directory name.* Transcripts are under `…/executor/` (the controller's brief), not the
   plan's `…/exec/`. No content difference.

No production file, `.ail` file, or other test was changed.

## §13 Evaluation

- **r1** — evaluator `sonnet` (Agent tool, own detached worktree `.eval-world-iter185` @ `2c3b997`): **PASS 95/100, zero blocking.** It re-derived every gate, AC4 and AC7; spot-checked H6, H4, V6, V10 and Q17 exactly; confirmed the milestone split (H5 killed at M1 `133a83f`, V6 survives there); and proved Q17 equivalent algebraically. Of its 3 self-chosen mutants, 1 was killed, 1 was equivalent (pair-check reorder), and **1 SURVIVED**: the FAIL span's glyph `✗ verdict:` → `✓ verdict:`. Report: `~/.ailang/state/mission-world-iter185-evidence/EVAL_REPORT_iter185.md`.
- **Controller closure** — the survivor was reproduced first-party, along with a second same-class survivor (FAIL span `class="verdict-fail"` → `"verdict-pass"`). `TestGradeViewRequiresTestVerdict/fail` now asserts the exact FAIL span, as `/pass` does, and is the sole killer of both; the G3 control (drop the verdict text) was already killed. Committed `121682d` (+4 lines; `render.go` byte-identical). Drill: `controller_g/{before,after}.txt`.
- **r2** — judge on that delta (own worktree @ `121682d`): **PASS 98/100, zero blocking.** It reproduced the survivor at `2c3b997` and the kill at `121682d`, and ran 3 further FAIL-span mutants, all killed. The class swap is killed only by the new lines. Report: `…/EVAL_REPORT_iter185_r2.md`.
- Final size: `git diff --shortstat fd99840 -- host/` → 2 files, **94** insertions, 0 deletions.
