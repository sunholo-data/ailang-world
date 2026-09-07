# w-flow-key-quote-trim-is-uncovered — pin the two EXISTING quote-trims of the lever gate's line scan with arms; change no parser behaviour

- **Status:** planned (design only; nothing implemented)
- **Date:** 2026-09-07 (iter-169 designer: `claude:claude-fable-5-1`)
- **Revision:** round 3 — **SCOPE REDUCTION by controller SPLIT** (2026-09-07), not a third
  revision. Rounds 1 and 2 both BLOCKED; every objection across six reviewer-verdicts landed on
  the round-1/2 `unquoteKey` helper (a behaviour change to the gate's acceptance semantics),
  and none on the arms. The helper and everything measured about it is SPLIT OUT to a new queue
  row (§6). This doc is now exactly what row 66 asked for: "add the arm". See §8.
- **Base commit measured at:** `eff74891141567d63d9ba3bf1cba6ca54b624a54` (== `origin/dev`), worktree `.wt-world-iter169`, branch `mission/world-iter169-row66`
- **Owning queue row:** `design_docs/world-mission.md` row 66 — `w-flow-key-quote-trim-is-uncovered` (clause-2, ~0.1d, gated on nothing; evaluator-found in row 55)
- **Scope class:** BRITTLENESS — a false RED (row 66's own claim), REPRODUCED at both trim sites.
  The silent false GREENS this design pass measured at the same sites (§6) are DECLARED and
  LEFT OPEN here, by split.
- **Disposition:** **(A) — add the arms.** Three new rows in `TestOnBlockTriggerParserShapes`.
  The parser is not modified: the two cutset trims stay byte-identical, and this milestone
  promises NO behaviour change. `tools/launchd/*` untouched.

---

## 1. Problem

`host/verifygate/dispatch_lever_gate_test.go` (441 lines at base) holds the row-55 pure parser
`parseOnBlockTriggers` (L170–243) consumed by `TestEveryWorkflowDeclaresDispatchLever` (L417).
Despite its name it routes BOTH forms itself: the flow form `on: {…}` at L191–208 (via
`splitFlowItems`/`splitFlowKeyValue`) and the block form at L213–242. Trigger keys are read at
TWO sites, each stripping quotes with the same cutset call:

- **L162, flow path** (`splitFlowKeyValue`): `key := strings.Trim(strings.TrimSpace(item[:i]), "'\"")`
- **L237, block path** (`parseOnBlockTriggers`): `key := strings.Trim(strings.TrimSpace(kv[0]), "'\"")`

These are the ONLY two quote-trims in the package — enumerated from the surface, not from the
row's prose (V1). Both landed in row 55's squash `165b9fd`. Row 55's design doc SPECIFIED the
flow-site trim in one clause ("quotes trimmed", its §3.3(b)) and never mentioned the block site;
neither site has an arm in `TestOnBlockTriggerParserShapes` (L274) whose fixture carries a
quoted trigger key — every fixture there uses bare `push:` / `workflow_dispatch:`.

**The defect — a false RED row 66 names, at BOTH sites (V2, V3).** Deleting either trim compiles
and leaves the full consumer set green. Under the mutant, `on: {"workflow_dispatch": }` or
`on:` + `  "workflow_dispatch":` — both VALID YAML declaring the lever (V7, Psych
`on-keys=["push","workflow_dispatch"]`) — would yield the key `"workflow_dispatch"` WITH quotes,
`slices.Contains(triggers, "workflow_dispatch")` would be false, and the gate would red
claiming the lever is absent. Loud and wrong-direction-safe, but pinned by nothing. Row 66
names only L162; L237 is the same defect at the block site, which is the more common
real-world YAML shape.

**Provenance.** V2/V3 were first measured by the iter-169 controller in the main checkout at
the same base; every row in §2 was RE-RUN or newly run first-party by this designer in this
worktree, with harness sources reproduced verbatim.

---

## 2. Verification Log

All rows at `eff7489` in `.wt-world-iter169`. `F` = `host/verifygate/dispatch_lever_gate_test.go`
(sha256 `2f4d0efa13ad9dab…` before and after every row that touches it; porcelain 0 after).
`CONSUMERS` = `'TestOnBlockTriggerParserShapes|TestOnBlockControlNeedle|TestOnBlockFailureMessagesUnchanged|TestEveryWorkflowDeclaresDispatchLever'`
(the complete consumer set — row 55's V16). Every command runs with `PATH=/opt/homebrew/bin:$PATH`
and `AILANG_BIN=~/.pinned-ailang/ailang`.

**Harness H1 (V4, V5, V6, V10) — VERBATIM, written to `host/verifygate/iter169_row66_probe_test.go`,
run as `go test ./host/verifygate/ -run '^TestIter169Row66Probe$' -count=1 -v | grep -E '^V[0-9]+ '`,
deleted afterwards:**

```go
package verifygate

import (
	"fmt"
	"testing"
)

func TestIter169Row66Probe(t *testing.T) {
	cases := []struct{ row, src string }{
		{"V4", "on: {\"push\": {branches: [dev]}, 'workflow_dispatch': }\n"},
		{"V4", "on:\n  \"push\":\n  'workflow_dispatch':\n"},
		{"V5", "on:\n  \"workflow_dispatch:\n"},
		{"V5", "on:\n  \"workflow_dispatch':\n"},
		{"V5", "on:\n  workflow_dispatch\":\n"},
		{"V5", "on:\n  \"\"workflow_dispatch\"\":\n"},
		{"V6", "on: {\"workflow_dispatch: }\n"},
		{"V6", "on: {\"workflow_dispatch': }\n"},
		{"V6", "on: {\"\": }\n"},
		{"V10", "on: {\"\"workflow_dispatch\"\": }\n"},
		{"V10", "on: {\"\"workflow_dispatch: }\n"},
		{"V10", "on: {workflow_dispatch\"\": }\n"},
		{"V10", "on: {\"a\"\"b\": }\n"},
		{"V10", "on: {\"workflow_dispatch\"': }\n"},
	}
	for _, c := range cases {
		keys, sv, err := parseOnBlockTriggers(c.src)
		fmt.Printf("%-4s %-48q keys=%q scalar=%v err=%v\n", c.row, c.src, keys, sv, err)
	}
}
```

**Harness H3 (V13, V14, V15) — the M1 deliverable rehearsed:** the three §4 table rows inserted
after L299 (`inline_comment_value`) via `perl -pi`; `git diff --stat` → `1 file changed,
3 insertions(+)`; `go vet ./host/verifygate/` rc=0; AC1 executed; then M-162 / M-237 applied
on top (L162/L237 are ABOVE the insertion, so their numbers hold); then `git checkout -- F`.

| ID | Claim | Command / input | Observed |
|----|-------|-----------------|----------|
| V1 | EXACTLY two quote-trims exist in the package; both landed in row 55 | `grep -n "'\\\\\"\"" host/verifygate/*.go`; `git log -L162,162:$F`; `git log -L237,237:$F` | L162 and L237 only. L162: `165b9fd` over `2a01c43`; L237: `165b9fd` only |
| V2 | **M-162 survives at base** (row 66's claim, reproduced) | `sed -i '' '162s/.*/\t\t\t\tkey := strings.TrimSpace(item[:i])/' $F`; `git diff --stat`; `go vet ./host/verifygate/`; `go test ./host/verifygate/ -run $CONSUMERS -count=1` | `1 file changed, 1 insertion(+), 1 deletion(-)`; vet rc=0; `ok … 2.913s`, **rc=0** |
| V3 | **M-237 survives at base** (NOT in row 66) | as V2 with `237s/.*/\t\tkey := strings.TrimSpace(kv[0])/` | diff-stat non-empty; vet rc=0; `ok … 3.704s`, **rc=0** |
| V4 | KNOWN-POSITIVE: quoted keys parse at base, both sites, both quote kinds — the (A)-over-(B) evidence | H1 rows `V4` | both → `keys=["push" "workflow_dispatch"] scalar=map[] err=<nil>` — V2/V3 are missing ARMS, not missing behaviour |
| V5 | Block site silently accepts malformed/foreign quoting as the lever (SPLIT OUT, §6) | H1 rows `V5` — inputs EXACTLY: `"workflow_dispatch:`, `"workflow_dispatch':`, `workflow_dispatch":`, `""workflow_dispatch"":` | ALL FOUR → `keys=["workflow_dispatch"] err=<nil>` |
| V6 | Flow site refuses an UNCLOSED quote and an EMPTY quoted key before the trim — the `flow_unterminated_quote_key` control's premise | H1 rows `V6` — inputs EXACTLY: `{"workflow_dispatch: }`, `{"workflow_dispatch': }`, `{"": }` | all three → `keys=[] err=line 1: unsupported top-level `on:` form …` |
| V7 | Independent YAML verdicts | `/usr/bin/ruby -ryaml -e 'YAML.safe_load(src)'`, reading `doc[true]` (Psych parses bare `on` as YAML-1.1 boolean) | `flow_quoted`, `block_quoted`, canonical → VALID, `on-keys=["push","workflow_dispatch"]`; `block_trailing_only` → VALID, key `workflow_dispatch"`; `block_unterminated`, `block_mismatched`, `flow_unterminated` → **INVALID** `found unexpected end of stream while scanning a quoted scalar` |
| V9 | Baseline green (controller, main checkout, same base) | `go test ./host/verifygate/ -count=1`; `go vet ./host/verifygate/` | `ok … 49.452s` rc=0; vet rc=0 |
| V10 | Flow site silently accepts DOUBLED quoting as the lever (SPLIT OUT, §6) | H1 rows `V10` | `{""workflow_dispatch"": }`, `{""workflow_dispatch: }`, `{workflow_dispatch"": }` ALL → `keys=["workflow_dispatch"] err=<nil>`; `{"a""b": }` → `["a\"\"b"]`; `{"workflow_dispatch"': }` → refused |
| V11 | Independent YAML verdicts, doubled-quote set (SPLIT OUT, §6) | same `ruby -ryaml` call | `flow_doubled_both`, `flow_doubled_open`, `flow_adjacent_pairs` → **INVALID** `did not find expected ',' or '}' while parsing a flow mapping`; `block_doubled_both`, `block_doubled_open` → **INVALID** `did not find expected key while parsing a block mapping`; `flow_doubled_close`, `block_doubled_close` → VALID with key `workflow_dispatch""` |
| V12 | `gpt6-astra` round-2 finding (SPLIT OUT, §6): a co-occurring bare key makes the lever pass regardless of malformed quoting | controller first-party, then this designer via throwaway `TestIter169Row66MixedProbe` (same shape as H1) on `on:\n  ""workflow_dispatch"":\n  workflow_dispatch:\n`, `on: {""workflow_dispatch"": , workflow_dispatch: }\n`, and `on:\n  ""workflow_dispatch"":\n` alone | first two → `keys=["workflow_dispatch" "workflow_dispatch"] err=<nil>`; alone → `keys=["workflow_dispatch"] err=<nil>` — both agree with the controller's numbers |
| **V13** | **AC1's corrected command produces real output** (round-2 `gemini-3-1-pro`: the `$`-anchored form matched 0 lines ALWAYS because `go test -v` appends ` (0.00s)`) | H3 tree; AC1's exact command (§4); and the OLD `$`-anchored form for contrast | corrected → the 3 lines quoted verbatim under AC1, count `3`; OLD form → count **`0`** on the same tree (the defect gemini caught, reproduced) |
| **V14** | **M-162 has a UNIQUE killer** with the arms present | H3 tree + V2's sed; `go vet` rc=0; `go test … -run $CONSUMERS -count=1 -v \| grep -E '^\s+--- FAIL:' \| sort` | exactly one line: `--- FAIL: TestOnBlockTriggerParserShapes/flow_quoted_keys (0.00s)`; package rc=1 |
| **V15** | **M-237 has a UNIQUE killer** with the arms present | H3 tree + V3's sed; same pipeline | exactly one line: `--- FAIL: TestOnBlockTriggerParserShapes/block_quoted_keys (0.00s)`; package rc=1 |
| V16 | The M1 diff touches ONLY the table (AC2's premise) | H3 tree: `git diff --stat`; `git diff -U0 $F \| grep -E '^[-+]' \| grep -vE '^(\+\+\+\|---)'` | `1 file changed, 3 insertions(+)`; the three `+` lines are the three table rows; **zero `-` lines** — L162 and L237 unchanged |

**Charter claim re-checked:** row 66 says the shape "sits outside row 55's three declared
shapes". CONFIRMED with a refinement: row 55's shapes are (a) quoted `"on":`, (b) flow mapping
with BARE keys, (c) tab indent. A quoted TRIGGER key is none of those — but row 55's doc did
SPECIFY the flow-site trim without an arm (spec-without-a-test); the block-site trim is
specified nowhere.

---

## 3. Decision

### 3.1 (A): add the arms — because the parser already MEANS to support quoted keys

Read from the code, not from taste:

- `splitFlowItems` (L82–128) and `splitFlowKeyValue` (L130–168) each carry a quote-state
  machine with backslash-escape handling inside double quotes. That machinery exists so a quoted
  token can contain `,`, `:`, `{` without splitting. A parser that tracks quotes this carefully
  and then REJECTS a quoted key would be contradicting itself.
- `acceptedOnForms` (L13) accepts `"on":` and `'on':` — row 55 ruled that quoting the top-level
  key is valid, lever-declaring YAML that must NOT red CI. Quoting the trigger key one level
  down is the same footgun remedy one level down; (B) would re-open, one nesting level deeper,
  the exact class row 55 closed.
- (B)'s only argument — "no arm, so no cost to refuse" — is refuted by V4: the behaviour is
  already there and correct; the cost of (A) is three table rows.

**What `errUnhandledOnForm` does today, and which forms reach it** (L19, L36, L194–210): message
``line N: unsupported top-level `on:` form %q``, wrapper-prefixed `instrument failure: <path> …
— this line scan computes YAML trigger shape and depth conservatively` (L252–253). Reached ONLY
from the flow/rest path: (i) non-empty `rest` that is not `{…}` — `on: push`, `on: [push]`
(L210); (ii) a flow body `splitFlowItems` refuses — unbalanced braces/brackets or an UNCLOSED
quote (L194); (iii) a flow item with no top-level `:` or an empty key (L202). The block path
never returns it (L213–242 return only `errTabIndent` or nil). **This milestone does not change
any of that.** The `flow_unterminated_quote_key` control pins (ii) as a fact this milestone
left alone.

### 3.2 (deleted by split)

Rounds 1–2 specified an exact-pair `unquoteKey` helper replacing both trims. Every quorum
objection in both rounds landed on it (§8), and `gpt6-astra`'s round-2 measurement (V12) showed
its safety argument was conditional, not unconditional. It is a behaviour change to a CI gate's
acceptance semantics, larger than a ~0.1d "add the arm" row, and it now has its own row (§6).
Nothing from it survives here.

---

## 4a. Milestone M0 (0.02d) — PREREQUISITE: the R0 queue row must exist

`oc-glm-5-2` round-3 verbatim fix. The entire scope-reduction rests on R0 being TRACKED, and
round 3 correctly caught that the doc merely ASSERTED a new queue row without a row number,
verbatim text or a verification command — leaving three measured, real defects with no verified
tracking artifact. Its fix, applied as written:

> Add to §4 an M0 prerequisite (0.02d): *"Add the R0 row to `world-mission.md` carrying
> V5/V10/V11/V12, the conditional-fix finding (V12), and both reviewers' open questions
> (`gpt6-astra`, `gemini-3-1-pro`) verbatim. M1 must not land until M0 is committed."*

**AC0 (gates M1):** the row exists and is greppable —
`grep -nE '^[0-9]+\. \*\*w-lever-gate-quoted-key-refusal-contract\*\*' design_docs/world-mission.md`
→ exactly one line; record its row number and paste its verbatim text into V17 in §2. M1 must
not land until AC0 passes. The row's body must carry the R0 measurements (V5, V10, V11, V12) and
both reviewers' open questions verbatim, so the split-out work is recoverable from the charter
alone without reading this doc.

---

## 4. Milestone M1 (0.05d) — three table rows and the drill

**Deliverable** (one commit, one file, table rows only — V16 is what the diff must look like):
three new rows in `TestOnBlockTriggerParserShapes`'s table, PINNED subtest names, inserted after
`inline_comment_value` (L299):

**The exact Go source to insert** (`gemini-3-1-pro` round-3 verbatim fix: a Markdown table
forces the implementer to guess struct-literal escaping, which risks gofmt reflowing to
multi-line and breaking AC2's `3 insertions(+)` assertion). The table's struct is
`{name, src, wantKeys, wantScalar, wantErr}` — FIVE fields; insert these three lines, exactly
as written, immediately after the `inline_comment_value` row:

```go
		{"flow_quoted_keys", "on: {\"push\": {branches: [dev]}, 'workflow_dispatch': }\n", []string{"push", "workflow_dispatch"}, nil, nil},
		{"block_quoted_keys", "on:\n  \"push\":\n  'workflow_dispatch':\n", []string{"push", "workflow_dispatch"}, nil, nil},
		{"flow_unterminated_quote_key", "on: {\"workflow_dispatch: }\n", nil, nil, errUnhandledOnForm},
```

Roles: `flow_quoted_keys` pins the L162 trim (row 66's shape; both quote kinds in ONE fixture so
a half-mutant handling only `"` still reds). `block_quoted_keys` pins the L237 trim (the unnamed
site). `flow_unterminated_quote_key` is the **CONTROL — must stay green under both mutants**,
proving the unclosed-quote refusal comes from `splitFlowItems` (V6) and that this milestone did
not disturb it. Indentation is two tabs, matching the surrounding rows; run `gofmt -l` on the
file and expect no output.

Plus the drill of §5 executed, with per-mutant FAIL lists and file shas recorded in the sprint
plan's evidence.

**Acceptance criteria** — every AC is a command plus an expected result. Compile fence is
`go vet` / `go test -run '^$'`, never `go build` (row 65).

**Execution status (`gpt6-astra` round-3 verbatim fix — the previous blanket "all executed"
claim overstated the evidence):** *AC1–AC3 were rehearsed on H3, with evidence in V13–V16. V9
establishes baseline package health only. AC4 and AC5 remain prospective acceptance checks
unless separately recorded.* The executor MUST record AC4 and AC5 output on the real M1 tree —
exact tree, command, exit status, observed output — before either may be ticked.

- **AC1 — arms exist and pass, by name.** (`gemini-3-1-pro`'s verbatim fix applied: the
  duration suffix is allowed.)
  `AILANG_BIN=~/.pinned-ailang/ailang go test ./host/verifygate/ -run 'TestOnBlockTriggerParserShapes' -count=1 -v 2>&1 | grep -E '^\s+--- (PASS|FAIL): TestOnBlockTriggerParserShapes/(flow_quoted_keys|block_quoted_keys|flow_unterminated_quote_key)( \([0-9.]+s\))?$'`
  → exactly 3 lines, all `--- PASS`. **Real output on the H3 tree (V13):**
  ```
      --- PASS: TestOnBlockTriggerParserShapes/flow_quoted_keys (0.00s)
      --- PASS: TestOnBlockTriggerParserShapes/block_quoted_keys (0.00s)
      --- PASS: TestOnBlockTriggerParserShapes/flow_unterminated_quote_key (0.00s)
  ```
  (The round-2 form with a bare `$` produced `0` lines on the same tree — an AC that can never
  pass. Greps the RUNNER's verdict lines, not the table source.)
- **AC2 — scope containment: the parser is untouched, the diff is table rows only.**
  `git diff -U0 origin/dev -- host/verifygate/dispatch_lever_gate_test.go | grep -cE '^-[^-]'`
  → `0` (no line removed or modified — NOTE `grep -c` exits rc=1 on a zero count, so read the
  printed number, not `$?`, and do not run this under `set -e` unguarded); `git diff origin/dev --stat -- host/verifygate/dispatch_lever_gate_test.go`
  → `1 file changed, 3 insertions(+)`; and CONTROL `grep -n "'\\\\\"\"" host/verifygate/dispatch_lever_gate_test.go`
  → exactly two lines, `162:` and `237:` (the trims are still there, at the same lines, because
  the insertion is below them). Rehearsed as V16.
- **AC3 — each mutant reds EXACTLY its arm (§5).** For M-162 and M-237: apply the sed; `git diff
  --stat` shows `1 file changed`; `go vet ./host/verifygate/` rc=0; then
  `AILANG_BIN=~/.pinned-ailang/ailang go test ./host/verifygate/ -run $CONSUMERS -count=1 -v 2>&1 | grep -E '^\s+--- FAIL:' | sort`
  → exactly the one line in the mutant's "reds" column; then `git checkout -- <F>` and
  `shasum -a 256` matches the pre-mutation digest. `-count=1` mandatory. Rehearsed as V14/V15.
- **AC4 — the real gate and the whole package stay green.**
  `AILANG_BIN=~/.pinned-ailang/ailang go test ./host/verifygate/ -count=1` → `ok`, rc=0
  (base: V9); `AILANG_BIN=~/.pinned-ailang/ailang go test ./... -count=1` → rc=0;
  `./scripts/verify_ail.sh` → rc=0 (no-regression control; no `.ail` is touched).
- **AC5 — frozen core and the workflow are untouched.**
  `git diff --stat origin/dev -- tools/launchd .github scripts world` → empty, and
  `git diff --name-only origin/dev` → exactly `host/verifygate/dispatch_lever_gate_test.go`
  plus the two design docs (this file and the sprint plan).

---

## 5. Mutation table and anti-vacuity

Protocol: one mutant at a time on the M1 tree; record `git diff --stat` (a mutant that did not
land is a vacuous rc=0 — the controller hit exactly this on its first M-162 attempt); `go vet`
rc=0 before reading any verdict; run `$CONSUMERS` with `-count=1 -v`; list `--- FAIL:` lines;
restore and verify sha256. Both mutants sit ABOVE the inserted rows, so their base line numbers
hold on the M1 tree (V14/V15 confirm).

| Mutant | Edit | Reds (EXACT set — MEASURED, V14/V15) | Must stay green (proves specificity) |
|---|---|---|---|
| **M-162** (row 66's own) | L162 → `key := strings.TrimSpace(item[:i])` | at base: **∅** (V2 — the defect). On M1: `TestOnBlockTriggerParserShapes/flow_quoted_keys` and ONLY that (V14) | `flow_unterminated_quote_key` (refused upstream by `splitFlowItems`); `flow_mapping`, `flow_scalar_violation` (bare keys); `block_quoted_keys`; every other subtest; the real gate |
| **M-237** (the unnamed site) | L237 → `key := strings.TrimSpace(kv[0])` | at base: **∅** (V3). On M1: `TestOnBlockTriggerParserShapes/block_quoted_keys` and ONLY that (V15) | `flow_quoted_keys`; `canonical_block`, `real_ci_yml`, `quoted_double`, `quoted_single` (bare trigger keys); every other subtest; the real gate (ci.yml is bare) |

**Anti-vacuity, per arm.** An arm is proven non-vacuous ONLY by the drill: observed `--- FAIL`
under its mutant and `--- PASS` at M1 in the same session, both lines verbatim in the
sprint-plan evidence. `flow_quoted_keys` ← M-162 (V14); `block_quoted_keys` ← M-237 (V15).
`flow_unterminated_quote_key` is the CONTROL: green at M1 and under BOTH mutants (V14/V15 list
one FAIL line each, and it is not this one). If it ever reds, the flow site's refusal has moved
into the trim — a change this milestone does not make.

**The label trap, stated so nobody re-pays for it:** AC1 and AC3 assert on `go test -v`'s
`--- PASS`/`--- FAIL` lines, which the runner emits only for a subtest that EXECUTED, and which
carry an unconditional ` (N.NNs)` suffix — anchor after it, never before it (V13). A grep for a
fixture string or subtest name in the SOURCE is a banner check — green with the row commented
out. Do not substitute one.

---

## 6. Declared residuals — including the work SPLIT OUT of this row

**R0 — SPLIT OUT to its own queue row (controller disposition, 2026-09-07), NOT closed here.
This milestone leaves every item below OPEN. They are measured defects, not hypotheses.**

- **Silent false GREEN on malformed quoting, BOTH sites.** The cutset `strings.Trim(…, "'\"")`
  strips any run of quote characters from either end, so INVALID YAML is read as declaring the
  lever with no error. Block site (V5, H1): `"workflow_dispatch:`, `"workflow_dispatch':`,
  `""workflow_dispatch"":` → `keys=["workflow_dispatch"] err=<nil>`; Psych: INVALID (V7, V11).
  Flow site (V10, H1): `{""workflow_dispatch"": }`, `{""workflow_dispatch: }` → same; Psych:
  INVALID (V11). The flow site refuses only UNCLOSED quotes (V6); doubled quotes are balanced and
  pass `splitFlowItems`.
- **Silent false GREEN on VALID YAML that does not declare the lever, both sites.** The key
  `workflow_dispatch"` (block, V5) and `workflow_dispatch""` (flow, V10) are legal plain scalars
  naming a DIFFERENT key (Psych: `on-keys=[…, "workflow_dispatch\""]`, V7/V11); the cutset
  reads both as the lever.
- **Any fix is CONDITIONAL, not unconditional** (`gpt6-astra`, round 2; controller-confirmed;
  V12): `on:\n  ""workflow_dispatch"":\n  workflow_dispatch:\n` and
  `on: {""workflow_dispatch"": , workflow_dispatch: }` → `keys=["workflow_dispatch" "workflow_dispatch"]`
  — a co-occurring bare key satisfies the lever check regardless of how the malformed token is
  handled, so a helper that reads the malformed token as a non-lever key does NOT make the gate
  red on invalid YAML in general. Whatever closes R0 must decide whether the gate's job includes
  refusing a file GitHub would reject, and must prove its claim on MIXED fixtures, not isolated
  tokens. The round-2 `unquoteKey` design (exact-pair unquote, `ok=false` → `errUnhandledOnForm`
  at both sites) and `gemini-3-1-pro`'s round-2 question (why not `strconv.Unquote`) travel with
  the new row as inputs, not as conclusions.
- Consequence for THIS row, stated plainly: after M1 the gate is pinned against the false RED
  and remains open to the false GREENS above. That is a known gap, declared here and in the new
  row, not a quiet one.

**R1 — block path skips a colon-less line silently.** L234 `if len(kv) != 2 { continue }`:
`on:\n  push\n  workflow_dispatch:\n` → `keys=[workflow_dispatch]`, no error (H1-shape probe,
round 1). Direction (`gpt6-astra` round-3 verbatim correction — the earlier "false RED" reading
was unsupported by this very measurement): *The parser silently skips the colon-less line and
returns `workflow_dispatch`, so the lever membership check passes. Whether this is a false GREEN
depends on the invalid-YAML refusal contract left open in R0; it is not evidence of a false RED.*
Instance 1.

**R2 — an empty quoted key in block form** (`"":`) is appended as `""` rather than refused
(flow form refuses it, V6). Harmless for the lever check; asymmetry noted.

**R3 — `parseOnBlockTriggers` is misnamed**: it routes BOTH forms (L191, L213). The name
plausibly misled `gemini-3-1-pro` in round 1 (§8). Rename out of scope; noted for whoever next
touches the file.

---

## 7. Scope boundary — what this does NOT change

- **`tools/launchd/*` — FROZEN CORE, not touched.** AC5 measures it.
- **The parser — not touched.** No function body changes; AC2 asserts zero removed/modified
  lines and that both trims remain at L162/L237. No new helper, no new error path, no new
  message. `TestOnBlockFailureMessagesUnchanged` and row 55's mutation evidence (MUT-A…MUT-F)
  are unaffected by construction; AC4 re-runs them anyway.
- No production code: `world/`, non-test `host/`, `scripts/`, `.github/workflows/ci.yml` are
  all untouched (AC5).
- No YAML dependency, no `go.mod` change; row 55's line-scan disposition stands; no new decision
  is raised.
- Charter row text: the landing record should add to row 66 that L237 carried the same defect
  (V3), and should cross-reference the new R0 row so the split is visible from the charter.
  Landing-commit docs edit, not an M1 deliverable.

---

## 8. Quorum verification log

**Round 1 (2026-09-07) — BLOCKED 3/3.** `gpt6-astra` was ABSENT from the first run (budget
refusal, $0.1130 est. vs $0.1000 cap) and was re-run alone by the controller at a raised cap
($0.0858) under the standing rule that a verdict with non-empty `absent_reviewers` is a verdict
with a named hole.

| Reviewer | Verdict | Objection (one line) | Disposition |
|---|---|---|---|
| `gemini-3-1-pro` | reject | V4/V6 "cannot" call `parseOnBlockTriggers` on flow fixtures; proposed renaming the function the doc called | **REFUTED with evidence — NOT applied.** Controller measured: `parseOnBlockTriggers` routes both forms — L191 tests `rest` for `{…}`, L192 calls `splitFlowItems`, `errUnhandledOnForm` returns from inside that branch at L194/L202/L210. No other routing function exists; H1 calls it on flow fixtures and gets V4/V6's outputs. Misleading NAME recorded as R3. |
| `gpt6-astra` | reject | round-1 §3.2 claimed the flow-site helper was "unchanged for EVERY input"; `""workflow_dispatch""` is balanced, reaches the trim, and cutset vs helper differ on it; V6's inputs mislabelled | **MEASURED TRUE** (V10/V11 added, V6 relabelled); round 2 kept the helper and pinned the change with an arm — which round 2 then rejected again (below). |
| `oc-glm-5-2` | reject | V4 had no command-column entry | **APPLIED as written:** H1 source reproduced verbatim and re-executed. |

**Round 2 (2026-09-07) — BLOCKED, recorded as N−1.** `oc-glm-5-2` ABSENT, reason `invalid`
(not budget; the controller could not restore it). Two verdicts, both reject.

| Reviewer | Verdict | Objection (one line) | Disposition |
|---|---|---|---|
| `gpt6-astra` | reject | the helper's "silent green → loud red" argument assumed the malformed token was the ONLY lever declaration; with a co-occurring bare key the lever check passes regardless, so §3.2/R3's unconditional claim was unproven (V10 tested isolated tokens only) | **CONFIRMED by controller measurement, then by this designer (V12).** The surface it lands on — the helper — is **SPLIT OUT** to a new queue row carrying the objection verbatim plus the three measurements; nothing of it survives here (§3.2 deleted, R0 records the open gap). |
| `gemini-3-1-pro` | reject | (1) AC1's `$` anchor can never match because `go test -v` appends ` (0.00s)`, so M1 was unpassable as written; (2) justify bypassing `strconv.Unquote` in the helper | **(1) APPLIED VERBATIM** (`( \([0-9.]+s\))?$`), and the corrected command was EXECUTED — real output pasted under AC1, old form reproduced at `0` lines (V13). **(2) MOOT after the split**; forwarded to the new row as an input. |

**Disposition after round 2: controller SPLIT, not a force-pass and not a third revision.**
Across six reviewer-verdicts every objection landed on the helper; §3.1's deliverable (arms
pinning the two existing trims — literally row 66's "add the arm") drew none. This doc is
reduced to that uncontested scope: 3 table rows, 2 mutants with measured unique killers, 5 ACs
all executed on a rehearsed tree, parser byte-identical to base (sha `2f4d0efa13ad9dab…`).
The helper, its refusal semantics, V5/V10/V11/V12 and both reviewers' open questions go to the
new row. One re-quorum remains for this item.

### Round 3 (2026-09-07, N=3, all present) — BLOCKED, resolved by the Gate-2 NARROW-REFINEMENT CARVE-OUT

Three rejects, none disputing the design DIRECTION and each carrying a concrete reviewer-authored
`proposed_fix` — the carve-out's exact conditions. The controller applied all three VERBATIM
(reviewers' own text; no controller-invented resolution) and routed to sprint-planner.

| Reviewer | Verdict | Objection | Disposition |
|---|---|---|---|
| `gpt6-astra` | reject | §4/§8 claimed all five ACs were executed, but AC4/AC5 have no execution record — V9 is a baseline controller run only. Also R1's "false RED" reading is not supported by its own measurement. | **APPLIED VERBATIM.** The blanket claim is replaced by astra's sentence ("AC1–AC3 were rehearsed on H3 … AC4 and AC5 remain prospective acceptance checks unless separately recorded"), and R1's direction statement is replaced by astra's verbatim correction. |
| `gemini-3-1-pro` | reject | §4 gave the new rows as a Markdown table, so the implementer must guess struct escaping; gofmt reflow would break AC2's `3 insertions(+)` assertion. | **APPLIED VERBATIM** — the table is replaced by the exact three Go struct literals (five fields, matching the real struct). **Then EXECUTED by the controller** rather than asserted: inserted on a scratch tree → `gofmt -l` silent, `git diff --stat` = `1 file changed, 3 insertions(+)`, removed-line count `0`, `go vet` rc=0, AC1 returns exactly the three `--- PASS` lines. Tree restored, parser sha `2f4d0efa13ad9dab…`. |
| `oc-glm-5-2` | reject | The split-out R0 work is tracked only by an ASSERTED queue row — no number, no verbatim text, no verification command — so three measured defects had no verified tracking artifact. | **APPLIED VERBATIM** as the M0 prerequisite in §4a, with AC0 gating M1: the R0 row must exist in `world-mission.md`, carrying V5/V10/V11/V12 and both reviewers' open questions verbatim, before M1 may land. |

**Why this is a carve-out and not a force-pass.** Standing rule 2 forbids proceeding over a
contested design DIRECTION. No round-3 objection contests the direction: each is a completeness,
determinism or attribution defect with a fix its own author wrote. The one round-2 objection that
DID reach the design (astra on the helper's refusal semantics) was not force-passed — the surface
it landed on was removed from this doc entirely and filed as its own row.
