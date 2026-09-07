# Sprint plan — `w-flow-key-quote-trim-is-uncovered` (queue row 66)

- **Design doc (APPROVED, do not re-litigate):** `design_docs/planned/w-flow-key-quote-trim-is-uncovered.md`
- **Base commit:** `eff74891141567d63d9ba3bf1cba6ca54b624a54` (== `origin/dev`)
- **Worktree (work ONLY here):** `/Users/voightkampff/dev/sunholo-data/.wt-world-iter169`
- **Branch:** `mission/world-iter169-row66`
- **Executor lane:** `pi:ollama/deepseek-v4-flash:0731-cloud`
- **Size:** 0.07d (M0 0.02d + M1 0.05d)
- **Milestones:** 2 (M0 verify-and-commit-docs; M1 three table rows + mutation drill)
- **Acceptance criteria:** 6 (AC0–AC5)
- **Planner rehearsal:** every command below was executed by the planner on a scratch copy of the
  target file at this base. Where a command's output is quoted as "expected", the planner observed
  exactly that. Deviations are a REAL failure, not a transcription artifact — stop and report.

---

## 0. Environment preamble — prepend to EVERY bash call

The rig has a PATH gap and a mandatory test env var. Without these, `go`/`gh` are rc=127 and
`go test` is rc=1 **by design** — that is an instrument failure, not a result.

```bash
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN="$HOME/.pinned-ailang/ailang"
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter169
```

Two shell rules that have burned this rig before, both of which apply to commands in this plan:

- **Do not run under `set -e` unguarded.** `grep -c` exits rc=1 on a zero count, and AC2 EXPECTS a
  zero count. Read the printed NUMBER, not `$?`. Also, `set -e` resets `$?`.
- **Quote glob-ish flags.** `--include=*.go` unquoted is a zsh error (`no matches found`). Write
  `--include='*.go'`.

Timeouts: `go test ./host/verifygate/ -count=1` takes ~50s. `go test ./... -count=1` covers 19
packages and takes minutes. Allow at least 600000 ms for the full-suite run.

**Shorthands used below (expand them literally; do not rely on shell variables surviving between
calls):**

- `F` = `host/verifygate/dispatch_lever_gate_test.go`
- `CONSUMERS` = `'TestOnBlockTriggerParserShapes|TestOnBlockControlNeedle|TestOnBlockFailureMessagesUnchanged|TestEveryWorkflowDeclaresDispatchLever'`

---

## 1. STEP B — Baseline (run FIRST, before touching anything)

Purpose: a pre-existing red must be attributed to the repo, not to this change. Record all four
results in the evidence template §7.1 before making any edit.

```bash
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN="$HOME/.pinned-ailang/ailang"
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter169

# B1 — tree identity
git rev-parse HEAD
git status --porcelain
shasum -a 256 host/verifygate/dispatch_lever_gate_test.go

# B2 — compile fence (NEVER `go build`: it is rc=0 with a type error in a _test.go)
go vet ./host/verifygate/ ; echo "B2 rc=$?"

# B3 — package baseline
go test ./host/verifygate/ -count=1 ; echo "B3 rc=$?"

# B4 — full-suite baseline
go test ./... -count=1 ; echo "B4 rc=$?"
```

**Expected at B1:** `HEAD` = `eff74891141567d63d9ba3bf1cba6ca54b624a54`; `git status --porcelain`
shows exactly two entries — ` M design_docs/world-mission.md` (the controller's M0 row insertion,
uncommitted) and `?? design_docs/planned/w-flow-key-quote-trim-is-uncovered.md`, plus `??` for this
sprint-plan file; file sha256 =
`2f4d0efa13ad9dab015b85ce5eb7482d0888e3a338c8a7f4f24777dbd60bef2e`.

**Expected at B2/B3:** rc=0, `ok  ...ailang-world/host/verifygate  <N>s`.

**B4 is a RECORD, not a gate.** A green full suite prints `ok` for 19 packages. If B4 is rc≠0,
record the exact failing package and test name and CONTINUE — AC4's full-suite check is then
judged as "no NEW failure relative to B4", and you must say so explicitly in the evidence.
The known flake on this rig is `TestCLIRealSubprocessEpisode` in `cmd/ailang-worldd`. Do not
attempt to fix any pre-existing failure; it is out of scope.

**STOP CONDITION:** if B1's sha256 is not `2f4d0efa…`, the tree is not the one this plan was
rehearsed against. Stop and report; do not improvise.

---

## 2. MILESTONE M0 (0.02d) — verify the R0 queue row, then commit the docs

M0 is **verify-only for content**: the queue row 87 is ALREADY PRESENT in
`design_docs/world-mission.md` in this worktree (inserted by the controller, currently
uncommitted). **DO NOT re-add it, do not edit it, do not renumber anything.** Your job is to
prove it is there and commit it, because the design doc's §4a makes M1 ungated until M0 is
committed.

### M0.1 — verify AC0

```bash
export PATH=/opt/homebrew/bin:$PATH
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter169
grep -nE '^[0-9]+\. \*\*w-lever-gate-quoted-key-refusal-contract\*\*' design_docs/world-mission.md
grep -cE '^[0-9]+\. \*\*w-lever-gate-quoted-key-refusal-contract\*\*' design_docs/world-mission.md
```

**Pass condition:** the first command prints exactly ONE line and it begins `5370:87. `; the second
prints `1`. Record the full printed line in evidence §7.2 (V17).

### M0.2 — capture the row's verbatim text for the record

```bash
sed -n '5370,5395p' design_docs/world-mission.md
```

Paste the output into evidence §7.2 (V17) verbatim. **Do not edit the approved design doc**
to insert V17 — the sprint-plan evidence IS the record; the controller mirrors it at landing.

### M0.3 — commit the docs (named files only)

```bash
export PATH=/opt/homebrew/bin:$PATH
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter169
git add design_docs/world-mission.md \
        design_docs/planned/w-flow-key-quote-trim-is-uncovered.md \
        design_docs/planned/w-flow-key-quote-trim-is-uncovered-sprint-plan.md
git status --porcelain
git commit -m "docs(mission): row 87 R0 split-out + row 66 design and sprint plan (M0)"
git log --oneline -1
```

**Pass condition:** `git status --porcelain` after `git add` shows exactly three staged entries and
NOTHING else staged; the commit succeeds; `git log --oneline -1` shows the new commit.

**PROHIBITED:** `git add -A`, `git add .`, `git commit -a`. Stage named files only.

---

## 3. MILESTONE M1 (0.05d) — insert the three table rows

### M1.1 — write the three rows to a temp file (quoted heredoc; preserves tabs and backslashes)

The three lines are copied verbatim from the design doc §4. Indentation is TWO TABS. Use a
**quoted** heredoc delimiter (`<<'ROWS'`) so nothing is expanded, and **not** `<<-` (which would
strip the tabs).

```bash
export PATH=/opt/homebrew/bin:$PATH
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter169
cat > /tmp/iter169_row66_rows.txt <<'ROWS'
		{"flow_quoted_keys", "on: {\"push\": {branches: [dev]}, 'workflow_dispatch': }\n", []string{"push", "workflow_dispatch"}, nil, nil},
		{"block_quoted_keys", "on:\n  \"push\":\n  'workflow_dispatch':\n", []string{"push", "workflow_dispatch"}, nil, nil},
		{"flow_unterminated_quote_key", "on: {\"workflow_dispatch: }\n", nil, nil, errUnhandledOnForm},
ROWS
wc -l /tmp/iter169_row66_rows.txt
head -c 4 /tmp/iter169_row66_rows.txt | od -c | head -1
```

**Pass condition:** `wc -l` prints `3`; the `od -c` line begins `0000000  \t  \t   {   "` — i.e.
the leading whitespace is two REAL TAB bytes, not spaces.

### M1.2 — confirm the insertion anchor, then insert after line 299

Line 299 is the `inline_comment_value` row — the last row of
`TestOnBlockTriggerParserShapes`'s table.

```bash
sed -n '299p' host/verifygate/dispatch_lever_gate_test.go
sed -n '300p' host/verifygate/dispatch_lever_gate_test.go
```

**Pass condition:** line 299 contains `{"inline_comment_value",` and line 300 is `	}` (a tab and a
closing brace). If not, STOP — the anchor moved and the plan's line numbers are invalid.

```bash
sed -i '' '299r /tmp/iter169_row66_rows.txt' host/verifygate/dispatch_lever_gate_test.go
awk 'NR>=298 && NR<=303 {printf "%d\t%s\n", NR, $0}' host/verifygate/dispatch_lever_gate_test.go
```

**Expected:** lines 300, 301, 302 are the three new rows in the order
`flow_quoted_keys`, `block_quoted_keys`, `flow_unterminated_quote_key`; line 303 is `	}`.

### M1.3 — format and compile fence

```bash
gofmt -l host/verifygate/dispatch_lever_gate_test.go ; echo "gofmt-listed rc=$?"
go vet ./host/verifygate/ ; echo "vet rc=$?"
```

**Pass condition:** `gofmt -l` prints NO filename (empty output). `go vet` rc=0.

If `gofmt -l` prints the filename, the indentation or escaping is wrong. Do NOT run `gofmt -w` and
proceed — a reflow would break AC2's `3 insertions(+)` assertion. Instead run
`git checkout -- host/verifygate/dispatch_lever_gate_test.go` and redo M1.1/M1.2 exactly.

### M1.4 — AC1 and AC2 (see §6 for the exact commands and pass conditions)

Run AC1, then AC2. Both must pass BEFORE committing.

### M1.5 — commit M1 (named file only) and record the post-M1 digest

**This commit must happen BEFORE the mutation drill.** The drill restores with
`git checkout -- <file>`, which restores to **HEAD**, not to an uncommitted state — if the rows are
not yet committed, the restore silently DELETES them and every later result is garbage.

```bash
export PATH=/opt/homebrew/bin:$PATH
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter169
git add host/verifygate/dispatch_lever_gate_test.go
git status --porcelain
git commit -m "test(verifygate): pin both quoted-trigger-key trims of the lever gate line scan (row 66 M1)"
git log --oneline -1
shasum -a 256 host/verifygate/dispatch_lever_gate_test.go
```

**Expected post-M1 sha256 (planner-rehearsed):**
`2ff462227c77dd564f27a0c9062a849688dc1781a2cab4ca54f682fdbe28ed5b`

Call this **SHA_M1**. It is the restore target for every mutant in §4. The base sha `2f4d0efa…` is
the PRE-M1 digest and is NOT the drill's restore target — do not compare against it after M1.

**PROHIBITED:** `git add -A`. **PROHIBITED:** `git push`. **PROHIBITED:** `gh pr create`.

---

## 4. The mutation drill (part of M1) — AC3

Run this AFTER M1.5's commit. One mutant at a time, on a clean tree. Both mutants sit at lines
162 and 237, which are ABOVE the insertion point, so their line numbers are unchanged by M1.

The two claims being re-tested are the ones nobody ever checks: **M-162 kills `flow_quoted_keys`
and ONLY that; M-237 kills `block_quoted_keys` and ONLY that; the control
`flow_unterminated_quote_key` stays green under BOTH.** Re-running this drill is the point of the
milestone — do not tick it from the design doc's table.

### 4.1 — Mutant M-162 (the flow-site trim, row 66's own claim)

```bash
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN="$HOME/.pinned-ailang/ailang"
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter169

# 0) clean-tree precondition
git status --porcelain host/verifygate/dispatch_lever_gate_test.go   # expect EMPTY
shasum -a 256 host/verifygate/dispatch_lever_gate_test.go            # expect SHA_M1

# 1) show the line about to be mutated
sed -n '162p' host/verifygate/dispatch_lever_gate_test.go

# 2) apply the mutant
sed -i '' '162s/.*/\t\t\t\tkey := strings.TrimSpace(item[:i])/' host/verifygate/dispatch_lever_gate_test.go

# 3) PROVE THE MUTANT LANDED (a mutant that silently fails to apply gives a vacuous rc=0)
git diff --stat -- host/verifygate/dispatch_lever_gate_test.go
sed -n '162p' host/verifygate/dispatch_lever_gate_test.go

# 4) READ VET BEFORE ANY VERDICT
go vet ./host/verifygate/ ; echo "vet rc=$?"

# 5) run the consumer set
go test ./host/verifygate/ -run 'TestOnBlockTriggerParserShapes|TestOnBlockControlNeedle|TestOnBlockFailureMessagesUnchanged|TestEveryWorkflowDeclaresDispatchLever' -count=1 -v 2>&1 | grep -E '^\s+--- FAIL:' | sort
echo "package rc (pipeline note: rc below is grep's; rerun without the pipe if you need the test rc)"
go test ./host/verifygate/ -run 'TestOnBlockTriggerParserShapes|TestOnBlockControlNeedle|TestOnBlockFailureMessagesUnchanged|TestEveryWorkflowDeclaresDispatchLever' -count=1 > /dev/null 2>&1 ; echo "M-162 package rc=$?"

# 6) restore and verify
git checkout -- host/verifygate/dispatch_lever_gate_test.go
git status --porcelain host/verifygate/dispatch_lever_gate_test.go   # expect EMPTY
shasum -a 256 host/verifygate/dispatch_lever_gate_test.go            # expect SHA_M1
```

**Pass conditions for M-162:**

- step 3: `git diff --stat` prints `1 file changed, 1 insertion(+), 1 deletion(-)` — **non-empty is
  mandatory**; an empty diff means the mutant did NOT land and any green is vacuous. Step 3's
  `sed -n '162p'` must print `				key := strings.TrimSpace(item[:i])` (four leading tabs, no
  `strings.Trim(` wrapper).
- step 4: `go vet` rc=0. If vet is non-zero, the mutant did not compile — the drill is void; restore
  and report. **Never read a test verdict before vet.**
- step 5: the FAIL list is EXACTLY one line:
  `    --- FAIL: TestOnBlockTriggerParserShapes/flow_quoted_keys (0.00s)`
  (duration may differ). Package rc=1.
- `flow_unterminated_quote_key` and `block_quoted_keys` must NOT appear in the FAIL list.
- step 6: sha256 == SHA_M1 and porcelain empty.

### 4.2 — Mutant M-237 (the block-site trim, the site row 66 did not name)

Identical protocol; only steps 1–3 change:

```bash
sed -n '237p' host/verifygate/dispatch_lever_gate_test.go
sed -i '' '237s/.*/\t\tkey := strings.TrimSpace(kv[0])/' host/verifygate/dispatch_lever_gate_test.go
git diff --stat -- host/verifygate/dispatch_lever_gate_test.go
sed -n '237p' host/verifygate/dispatch_lever_gate_test.go
```

then steps 4, 5, 6 verbatim as in 4.1.

**Pass conditions for M-237:**

- step 3: `git diff --stat` prints `1 file changed, 1 insertion(+), 1 deletion(-)`; line 237 prints
  `		key := strings.TrimSpace(kv[0])` (two leading tabs).
- step 4: `go vet` rc=0.
- step 5: FAIL list is EXACTLY one line:
  `    --- FAIL: TestOnBlockTriggerParserShapes/block_quoted_keys (0.00s)`. Package rc=1.
- `flow_unterminated_quote_key` and `flow_quoted_keys` must NOT appear.
- step 6: sha256 == SHA_M1, porcelain empty.

### 4.3 — the control, stated explicitly

`flow_unterminated_quote_key` must be **PASS at M1 (AC1)** and **absent from both FAIL lists**. If
it ever reds, the flow site's unclosed-quote refusal has moved into the trim — a change this
milestone does not make. That is a STOP: restore, do not "fix", report.

---

## 5. Final gates — AC4 and AC5

**These have NOT been executed by anyone before you.** The design doc says so explicitly
(§4, `gpt6-astra` round-3 fix): "AC4 and AC5 remain prospective acceptance checks unless separately
recorded." You must run them on the real M1 tree — after M1.5's commit and after the drill has
restored the tree — and record exact tree (HEAD sha + file sha), exact command, exit status and
observed output. **They may not be ticked from the baseline or from the design doc.**

Precondition before running: `git status --porcelain` is EMPTY (both commits made, no mutant
residue) and `shasum -a 256` of the test file == SHA_M1.

---

## 6. Acceptance criteria — exact commands and pass conditions

| AC | Command (run from the worktree root, with §0's preamble) | Pass condition |
|---|---|---|
| **AC0** | `grep -nE '^[0-9]+\. \*\*w-lever-gate-quoted-key-refusal-contract\*\*' design_docs/world-mission.md` | Exactly ONE line, beginning `5370:87. `. Also `grep -c …` prints `1`. Row 87 is already present — VERIFY, do not re-add. Gates M1. |
| **AC1** | `go test ./host/verifygate/ -run 'TestOnBlockTriggerParserShapes' -count=1 -v 2>&1 \| grep -E '^\s+--- (PASS\|FAIL): TestOnBlockTriggerParserShapes/(flow_quoted_keys\|block_quoted_keys\|flow_unterminated_quote_key)( \([0-9.]+s\))?$'` | Exactly 3 lines, ALL `--- PASS`: `flow_quoted_keys`, `block_quoted_keys`, `flow_unterminated_quote_key`. Pipe to `\| wc -l` to record the count as `3`. **Do not drop the `( \([0-9.]+s\))?` suffix group** — `go test -v` always appends ` (0.00s)`; the bare-`$` form matched 0 lines and is a permanently-unpassable AC (V13). This greps the RUNNER's verdict lines, not the source; a `grep` for the fixture string in the file is a banner check and is NOT a substitute. |
| **AC2a** | `git diff -U0 origin/dev -- host/verifygate/dispatch_lever_gate_test.go \| grep -cE '^-[^-]'` | Prints `0`. **Read the printed number, not `$?`** — `grep -c` exits rc=1 on a zero count. Do not run under `set -e`. |
| **AC2b** | `git diff origin/dev --stat -- host/verifygate/dispatch_lever_gate_test.go` | `1 file changed, 3 insertions(+)` — and NO `deletions(-)` term. |
| **AC2c** | `grep -n "'\\\\\"\"" host/verifygate/dispatch_lever_gate_test.go` | Exactly TWO lines, prefixed `162:` and `237:` — both original trims still present, unchanged, at the same line numbers. |
| **AC3** | §4 in full: for M-162 and M-237 — diff-stat non-empty, `go vet` rc=0, then `go test … -run <CONSUMERS> -count=1 -v 2>&1 \| grep -E '^\s+--- FAIL:' \| sort` | M-162 → exactly `--- FAIL: TestOnBlockTriggerParserShapes/flow_quoted_keys`; M-237 → exactly `--- FAIL: TestOnBlockTriggerParserShapes/block_quoted_keys`; `flow_unterminated_quote_key` in NEITHER list; after each, `git checkout --` restores and sha256 == SHA_M1. `-count=1` is mandatory. |
| **AC4a** | `go test ./host/verifygate/ -count=1` | `ok  …/host/verifygate` and rc=0. Record rc explicitly. |
| **AC4b** | `go test ./... -count=1` | rc=0 (19 packages `ok`). If B4 was already red, the pass condition degrades to "no failure that was not in B4" — say so explicitly and name the pre-existing failure. |
| **AC4c** | `./scripts/verify_ail.sh ; echo "rc=$?"` | rc=0. No-regression control; no `.ail` is touched by this sprint. |
| **AC5a** | `git diff --stat origin/dev -- tools/launchd .github scripts world` | EMPTY output. Frozen core and the workflow untouched. |
| **AC5b** | `git diff --name-only origin/dev \| sort` | Exactly these FOUR paths (see §8 note 1 — the design doc's AC5 listed only three, omitting `world-mission.md`, which M0 necessarily modifies): `design_docs/planned/w-flow-key-quote-trim-is-uncovered-sprint-plan.md`, `design_docs/planned/w-flow-key-quote-trim-is-uncovered.md`, `design_docs/world-mission.md`, `host/verifygate/dispatch_lever_gate_test.go`. Nothing else. |

---

## 7. Evidence template — the executor FILLS THIS IN

Fill every cell. "As expected" is not evidence; paste the observed text. An empty cell fails the
sprint.

### 7.1 Baseline (before any edit)

| ID | Command | Exit status | Observed output |
|---|---|---|---|
| B1a | `git rev-parse HEAD` | 0 | `eff74891141567d63d9ba3bf1cba6ca54b624a54` |
| B1b | `git status --porcelain` | 0 | ` M design_docs/world-mission.md`; `?? design_docs/planned/w-flow-key-quote-trim-is-uncovered-sprint-plan.md`; `?? design_docs/planned/w-flow-key-quote-trim-is-uncovered.md` |
| B1c | `shasum -a 256 host/verifygate/dispatch_lever_gate_test.go` | 0 | `2f4d0efa13ad9dab015b85ce5eb7482d0888e3a338c8a7f4f24777dbd60bef2e` |
| B2 | `go vet ./host/verifygate/` | 0 | (no output; rc=0) |
| B3 | `go test ./host/verifygate/ -count=1` | 0 | `ok  github.com/sunholo-data/ailang-world/host/verifygate  47.743s` |
| B4 | `go test ./... -count=1` | 0 | all 19 packages `ok`; rc=0. No pre-existing failure. |

### 7.2 M0

| ID | Command | Exit status | Observed output |
|---|---|---|---|
| AC0 | (§6 AC0 grep) | 0 | `5370:87. **w-lever-gate-quoted-key-refusal-contract** · clause-2 · **THE DISPATCH-LEVER GATE READS` (exactly ONE line, begins `5370:87. `) |
| AC0-count | (§6 AC0 `grep -c`) | 0 | `1` |
| V17 | `sed -n '5370,5395p' design_docs/world-mission.md` | 0 | `87. **w-lever-gate-quoted-key-refusal-contract** · clause-2 · **THE DISPATCH-LEVER GATE READS` / `    INVALID YAML — AND VALID YAML NAMING A DIFFERENT KEY — AS DECLARING THE \`workflow_dispatch\`` / `    LEVER, AT BOTH TRIM SITES, WITH NO ERROR.** SPLIT OUT of row 66 at iteration 169 by controller` / `    disposition: across three quorum rounds every objection landed on this surface while row 66's` / `    own deliverable (arms pinning the two existing trims) drew none, so row 66 was reduced to the` / `    uncontested part and this is the remainder. \`strings.Trim\` with a cutset strips ANY run of` / `    quote characters from EITHER end, so at \`dispatch_lever_gate_test.go:237\` (block path) all of` / `    \`"workflow_dispatch:\`, \`"workflow_dispatch':\` and \`workflow_dispatch":\` parse to` / `    \`keys=[workflow_dispatch] err=nil\` — the first two are INVALID YAML (Psych: *found unexpected` / `    end of stream while scanning a quoted scalar*) and the third is VALID YAML whose key is` / `    literally \`workflow_dispatch"\`, which is NOT the lever (designer-measured, V5/V7/V11). The flow` / `    path at \`:162\` refuses the unterminated forms upstream in \`splitFlowItems\` but shares the` / `    defect for doubled quotes: \`on: {""workflow_dispatch"": }\` → \`keys=["workflow_dispatch"]\` and` / `    \`on: {workflow_dispatch"": }\` → the non-lever key (V10/V11). **The asymmetry between the two` / `    sites is itself undeclared.** **AND ANY FIX IS CONDITIONAL, NOT UNCONDITIONAL** —` / `    \`gpt6-astra\`'s round-2 objection, which the controller then reproduced first-party at base:` / `    \`on:\n  ""workflow_dispatch"":\n  workflow_dispatch:\n\` and its flow equivalent both return` / `    \`keys=[workflow_dispatch workflow_dispatch] err=<nil>\`, so a co-occurring bare key satisfies` / `    the lever check regardless of how the malformed token is handled. A design that reads the` / `    malformed token as a non-lever key therefore does NOT make the gate red on invalid YAML in` / `    general, and any proposal here must prove its claim on MIXED fixtures, not isolated tokens.` / `    **THE ITEM:** decide the gate's refusal contract — does its job include refusing a workflow` / `    file GitHub itself would reject? — then implement it and pin it. Two reviewer inputs travel` / `    with this row as INPUTS, not conclusions: the round-2 \`unquoteKey\` design (strip exactly one` / `    matching pair; \`ok=false\` routes to \`errUnhandledOnForm\` at both sites), and` / `    \`gemini-3-1-pro\`'s round-3 question of why \`strconv.Unquote\` is not reused (it rejects` (row 87 verbatim, lines 5370–5395) |
| M0-commit | `git commit -m "docs(mission): …"` + `git log --oneline -1` | 0 | `[mission/world-iter169-row66 88c8105] docs(mission): row 87 R0 split-out + row 66 design and sprint plan (M0)` / `88c8105 docs(mission): row 87 R0 split-out + row 66 design and sprint plan (M0)` |

### 7.3 M1 construction

| ID | Command | Exit status | Observed output |
|---|---|---|---|
| M1.1 | `wc -l /tmp/iter169_row66_rows.txt`; `od -c` first line | 0 | `3 /tmp/iter169_row66_rows.txt`; `0000000   \t  \t   {   "` (two real TAB bytes) |
| M1.2a | `sed -n '299p' …` (anchor) | 0 | line 299 = `		{"inline_comment_value", "on:\n  push:\n  workflow_dispatch: # manual re-run lever\n", []string{"push", "workflow_dispatch"}, nil, nil},`; line 300 = `	}` |
| M1.2b | `awk 'NR>=298 && NR<=303 …'` after insert | 0 | 300=`flow_quoted_keys`, 301=`block_quoted_keys`, 302=`flow_unterminated_quote_key`, 303=`	}` |
| M1.3a | `gofmt -l host/verifygate/dispatch_lever_gate_test.go` | 0 | (no filename printed; empty output) |
| M1.3b | `go vet ./host/verifygate/` | 0 | (no output; rc=0) |
| AC1 | (§6 AC1) | 0 | `    --- PASS: TestOnBlockTriggerParserShapes/flow_quoted_keys (0.00s)` / `    --- PASS: TestOnBlockTriggerParserShapes/block_quoted_keys (0.00s)` / `    --- PASS: TestOnBlockTriggerParserShapes/flow_unterminated_quote_key (0.00s)`; `wc -l` = `3` |
| AC2a | (§6 AC2a) | 0 | `0` (printed number) |
| AC2b | (§6 AC2b) | 0 | ` host/verifygate/dispatch_lever_gate_test.go | 3 +++` / ` 1 file changed, 3 insertions(+)` (no deletions) |
| AC2c | (§6 AC2c) | 0 | `162:				key := strings.Trim(strings.TrimSpace(item[:i]), "'\"")` / `237:		key := strings.Trim(strings.TrimSpace(kv[0]), "'\"")` |
| M1.5 | `git commit` + `shasum -a 256` | 0 | `[mission/world-iter169-row66 089aed9] test(verifygate): pin both quoted-trigger-key trims of the lever gate line scan (row 66 M1)` / `089aed9 test(verifygate): pin both quoted-trigger-key trims of the lever gate line scan (row 66 M1)` / **SHA_M1 = `2ff462227c77dd564f27a0c9062a849688dc1781a2cab4ca54f682fdbe28ed5b`** |

### 7.4 Mutation drill — one block per mutant (AC3)

**M-162**

| Field | Value |
|---|---|
| Pre-mutation sha256 | `2ff462227c77dd564f27a0c9062a849688dc1781a2cab4ca54f682fdbe28ed5b` |
| Pre-mutation porcelain | (empty) |
| `sed -n '162p'` BEFORE | `				key := strings.Trim(strings.TrimSpace(item[:i]), "'\"")` |
| `git diff --stat` (mutant landed?) | `1 file changed, 1 insertion(+), 1 deletion(-)` — NON-EMPTY, landed |
| `sed -n '162p'` AFTER | `				key := strings.TrimSpace(item[:i])` |
| `go vet` rc | 0 |
| `--- FAIL:` list (sorted, verbatim) | `    --- FAIL: TestOnBlockTriggerParserShapes/flow_quoted_keys (0.00s)` |
| package rc | 1 |
| Control `flow_unterminated_quote_key` in FAIL list? | **no** (observed) |
| Post-restore sha256 | `2ff462227c77dd564f27a0c9062a849688dc1781a2cab4ca54f682fdbe28ed5b` |
| Post-restore porcelain | (empty) |

**M-237** — same table, with line 237 and `strings.TrimSpace(kv[0])`.

| Field | Value |
|---|---|
| Pre-mutation sha256 | `2ff462227c77dd564f27a0c9062a849688dc1781a2cab4ca54f682fdbe28ed5b` |
| Pre-mutation porcelain | (empty) |
| `sed -n '237p'` BEFORE | `		key := strings.Trim(strings.TrimSpace(kv[0]), "'\"")` |
| `git diff --stat` (mutant landed?) | `1 file changed, 1 insertion(+), 1 deletion(-)` — NON-EMPTY, landed |
| `sed -n '237p'` AFTER | `		key := strings.TrimSpace(kv[0])` |
| `go vet` rc | 0 |
| `--- FAIL:` list (sorted, verbatim) | `    --- FAIL: TestOnBlockTriggerParserShapes/block_quoted_keys (0.00s)` |
| package rc | 1 |
| Control `flow_unterminated_quote_key` in FAIL list? | **no** (observed) |
| Post-restore sha256 | `2ff462227c77dd564f27a0c9062a849688dc1781a2cab4ca54f682fdbe28ed5b` |
| Post-restore porcelain | (empty) |

### 7.5 Final gates

| ID | Command | Exit status | Observed output | Tree at time of run (HEAD sha / file sha256) |
|---|---|---|---|---|
| AC4a | `go test ./host/verifygate/ -count=1` | 0 | `ok  github.com/sunholo-data/ailang-world/host/verifygate  46.854s` | `089aed904b0f52ff333fbe9c5ab7e7200734c03a` / `2ff462227c77dd564f27a0c9062a849688dc1781a2cab4ca54f682fdbe28ed5b` |
| AC4b | `go test ./... -count=1` | 0 | all 19 packages `ok`; rc=0 | `089aed904b0f52ff333fbe9c5ab7e7200734c03a` / `2ff462227c77dd564f27a0c9062a849688dc1781a2cab4ca54f682fdbe28ed5b` |
| AC4c | `./scripts/verify_ail.sh` | 0 | `✓ verify gate PASSED: 11 required identities verified, 40 named tests pass`; `✓ world package gate PASSED: 9/9 steps performed non-zero work` | `089aed904b0f52ff333fbe9c5ab7e7200734c03a` / `2ff462227c77dd564f27a0c9062a849688dc1781a2cab4ca54f682fdbe28ed5b` |
| AC5a | `git diff --stat origin/dev -- tools/launchd .github scripts world` | 0 | (empty output — frozen core and workflow untouched) | `089aed904b0f52ff333fbe9c5ab7e7200734c03a` / `2ff462227c77dd564f27a0c9062a849688dc1781a2cab4ca54f682fdbe28ed5b` |
| AC5b | `git diff --name-only origin/dev \| sort` | 0 | `design_docs/planned/w-flow-key-quote-trim-is-uncovered-sprint-plan.md` / `design_docs/planned/w-flow-key-quote-trim-is-uncovered.md` / `design_docs/world-mission.md` / `host/verifygate/dispatch_lever_gate_test.go` (exactly these four) | `089aed904b0f52ff333fbe9c5ab7e7200734c03a` / `2ff462227c77dd564f27a0c9062a849688dc1781a2cab4ca54f682fdbe28ed5b` |

### 7.6 Summary

- ACs passed: 6 / 6
- Commits made (sha + subject):
  - `88c8105` — `docs(mission): row 87 R0 split-out + row 66 design and sprint plan (M0)`
  - `089aed9` — `test(verifygate): pin both quoted-trigger-key trims of the lever gate line scan (row 66 M1)`
- Anything that deviated from this plan, and what you did: **None.** Every command was executed verbatim from the plan; all expected outputs matched (B1 sha, SHA_M1, AC0 line `5370:87. `, AC1 3 PASS lines, AC2a `0`, AC2b `3 insertions(+)`, AC2c lines 162/237, M-162 reds exactly `flow_quoted_keys`, M-237 reds exactly `block_quoted_keys`, control green under both, AC4a/b/c rc=0, AC5a empty, AC5b exactly four paths).

---

## 8. Under-specifications the planner resolved (read these; they are the traps)

1. **AC5's file list was incomplete in the design doc.** §4 AC5 says the changed set is the test
   file "plus the two design docs". But M0 necessarily modifies `design_docs/world-mission.md`
   (that is where row 87 lives), so the real set is FOUR paths. AC5b above lists all four. A
   three-path expectation would fail a correct sprint.
2. **AC5 is only meaningful after M0's commit.** `git diff --name-only origin/dev` does not list
   untracked files; both design docs are untracked at the start. Hence AC5 runs at §5, last.
3. **The restore target after M1 is SHA_M1, not the base sha.** §2 of the design doc quotes
   `2f4d0efa…` as the file's digest "before and after every row" — that is the PRE-M1 digest.
   After M1 the file legitimately differs. Comparing to `2f4d0efa…` after M1 would look like a
   failure and tempt a wrong "fix".
4. **M1 must be committed BEFORE the drill.** `git checkout -- <file>` restores to HEAD. If M1 is
   uncommitted, the first mutant's restore deletes the three new rows silently and every subsequent
   result is vacuous. The design doc's drill protocol does not say this because it was rehearsed on
   a scratch tree where base == HEAD.
5. **The insertion mechanism was unspecified.** The doc gives the three Go lines but not how to get
   them into the file (its rehearsal used `perl -pi`). This plan prescribes a quoted heredoc plus
   `sed -i '' '299r …'`, which the planner executed on a scratch copy of the exact base file:
   `gofmt -l` silent, 3 insertions, 0 deletions, lines 162/237 untouched, resulting sha256
   `2ff46222…`.
6. **AC0's "paste into V17 in §2" would edit the approved design doc.** This plan routes V17 into
   the sprint-plan evidence (§7.2) instead; the executor must NOT edit
   `w-flow-key-quote-trim-is-uncovered.md`.
7. **`$CONSUMERS` is never expanded in the design doc's AC text.** Every command in this plan
   spells the regex out literally.
8. **BSD `sed` vs `\t`.** Verified on this rig: macOS `sed -i ''` DOES emit real tab bytes for `\t`
   in the replacement, so the mutation seds produce compiling Go. (Checked with `od -c`, not
   assumed.)

---

## 9. Hard prohibitions for the executor

- **DO NOT push.** No `git push`, no `git push --force`, no branch publishing. The controller lands.
- **DO NOT open a PR.** No `gh pr create`.
- **DO NOT `git add -A` / `git add .` / `git commit -a`.** Stage named files only.
- **DO NOT modify `tools/launchd/*`** (FROZEN CORE), `.github/`, `scripts/`, or `world/`.
- **DO NOT modify the parser** — `host/verifygate/dispatch_lever_gate_test.go` lines 1–299 stay
  byte-identical. The only permitted edit is the three inserted rows.
- **DO NOT edit** `design_docs/planned/w-flow-key-quote-trim-is-uncovered.md` (approved).
- **DO NOT touch** `/Users/voightkampff/dev/sunholo-data/ailang` (a different mission's repo) or
  `~/.ailang/state/mission-v1*`.
- **DO NOT use `go build` as a compile fence** — it is rc=0 with a type error in a `_test.go`. Use
  `go vet` (or `go test -run '^$'`).
- **DO NOT substitute a source `grep` for AC1's runner-output grep.** A source grep is green with
  the row commented out.
- **DO NOT "fix" a pre-existing failure found at B4.** Record it and continue.

## 10. Stop conditions — report, do not improvise

- B1's file sha256 ≠ `2f4d0efa13ad9dab015b85ce5eb7482d0888e3a338c8a7f4f24777dbd60bef2e`.
- AC0 returns 0 lines or more than 1 line.
- M1.2's anchor check does not show `inline_comment_value` at line 299.
- `gofmt -l` names the file after insertion and a clean redo does not fix it.
- `go vet` is non-zero at any point.
- A mutant's `git diff --stat` is EMPTY (mutant did not land).
- A mutant's FAIL list is empty, has more than one line, or contains
  `flow_unterminated_quote_key`.
- AC2b shows any `deletions(-)`.
