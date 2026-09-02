# w-dispatch-lever-parser-false-reds-on-valid-yaml — harden the row-47 lever gate's line scan against three valid-YAML shapes

- **Status:** planned (design only; nothing implemented)
- **Date:** 2026-09-02 (iter-149 designer: `claude:claude-fable-5`)
- **Revision:** round 2 of 2 (2026-09-02). The quorum BLOCKED at round 1 on three premise
  objections; the controller measured each rather than forwarding it, and every measurement
  was RE-RUN first-party by this designer in this worktree at the same base commit before
  being folded in (V13 rewritten; V16–V20 new; §1(e), §3.1, AC6/AC7 repaired; needle audit
  added after AC8).
- **Base commit measured at:** `234d9da3ff6c76a0ab8b32b36f1fdfb4ab383c23` (== `origin/dev`), worktree `.wt-iter149-row55`
- **Owning queue row:** `design_docs/world-mission.md` row 55 — `w-dispatch-lever-parser-false-reds-on-valid-yaml` (clause-2, ~0.2d, gated on nothing)
- **Scope class: BRITTLENESS, not soundness.** The row-47 evaluator's attacks aimed at a
  SILENT FALSE GREEN all failed (hidden dotfile, nested subdirectory, case/extension
  variants, all nine canonical mutations, and the anti-vacuity floor's own precondition —
  charter row 55). Every defect this doc fixes fails LOUD, which is the accepted direction.
  This doc does not claim to close a soundness hole, because it measured none.

---

## 1. Problem

`host/verifygate/dispatch_lever_gate_test.go`'s helper `onBlockTriggerKeys` (L14–84),
consumed by `TestEveryWorkflowDeclaresDispatchLever` (L112–136), finds the trigger block by
EXACT byte equality `l == "on:"` at column 0 (L25) and computes indentation with
`strings.TrimLeft(l, " ")` — spaces only (L43, L50). Three shapes of **valid, lever-declaring
YAML** therefore red the gate:

- **(a) Quoted key `"on":`** → `instrument failure: <f> has no top-level `on:` trigger block`
  (V1). Quoting `on` is the STANDARD remedy for YAML 1.1 parsing bare `on` as boolean `true`
  — this gate breaks CI the day anyone applies the recommended fix for a famous GitHub
  Actions footgun. (In the full gate path it is even earlier: the known-positive control at
  L128, `strings.Contains(src, "on:")`, does not match `"on":` — the bytes `o`,`n`,`:` are
  not consecutive in `"on":` — so a fully-quoted file Fatals at the CONTROL before reaching
  the helper. Both sites must move together; see Decision.)
- **(b) Flow style `on: {push: {branches: [dev]}, workflow_dispatch: }`** → the identical
  Fatalf (V2), because the line is not byte-equal to `on:`.
- **(c) A TAB-indented first trigger line** → `TrimLeft(l, " ")` strips spaces only, so
  `lead=0 == onLead=0`, the loop reads the block as already exited, and the result is
  `keys=[]` (V3) — the caller then reports `on-block triggers=[] lack workflow_dispatch`,
  which **misreports total absence rather than a parse limitation**. For (c) the defect is
  the DISHONEST message, not the red: the fix is a loud, typed parse-limitation failure, not
  a green (see Decision §3.3).

Two documentation defects travel with the row (charter row 55, "carries with it"):

- **(d)** The scalar arm CASCADES: `workflow_dispatch: garbage` emits TWO messages (V6) —
  the scalar-value Errorf AND the caller's `triggers=[push] lack workflow_dispatch`. The
  row-47 doc's line 291–292 claims "EXACTLY ONE attributed message per defect … it never
  cascades" — measured FALSE. The behavior is kept (it is defense-in-depth: neutering either
  message alone still reds); the SENTENCE is corrected.
- **(e)** The row-47 doc's Declared Residual 3 is wrong in the safe direction in two claims:
  Go's `filepath.Glob(dir+"/*")` does NOT skip dotfiles (V7 — `.hidden.yml` IS enumerated;
  Go's Glob is not POSIX shell globbing), and a nested subdirectory is not "invisible" but a
  LOUD `t.Fatal` from `os.ReadFile` — `read <dir>: is a directory` (V7). A case-mismatched
  `CI.YML` is also enumerated by `*` (V7). The code comment at
  `dispatch_lever_gate_test.go:110-111` carries the SAME wrong claims IN DIFFERENT WORDS —
  NOT verbatim (round 1 said "duplicated"; V18 refutes that wording: the doc's phrases
  `skips dotfiles` and `invisible` both grep 0 in the Go file, with `nested subdirectory` → 1
  as the same-file known-positive control). L110–111 exactly:

  ```
  // root-level .yaml, a hidden file), a case-mismatched filename (the Glob is case-sensitive),
  // or a nested subdirectory (which GitHub itself does not scan either).
  ```

  The comment's "cannot see … a hidden file" is refuted by V7 (Glob DOES return
  `.hidden.yml`); its case-sensitivity aside is irrelevant to a bare `*` pattern (V7:
  `CI.YML` IS enumerated); and it describes the nested-subdirectory case as unseen when it
  is a loud Fatal. M3 repairs the comment in the comment's own words, not the doc's.

**Provenance:** defects (a)–(e) were measured by the iter-149 controller this session
(detached worktree at `234d9da`, throwaway test calling `onBlockTriggerKeys` directly,
worktree porcelain restored to 0). **Every measurement below was RE-RUN first-party by this
designer in the same worktree at the same base commit** — the Verification Log rows are this
designer's observed outputs, not copies. The round-2 rows (V16–V20, plus the V13 rewrite)
originate from the quorum's premise objections: the controller measured each objection, and
this designer re-ran every command first-party in this worktree before folding it in.

---

## 2. Verification Log

Base commit for every row: `234d9da3ff6c76a0ab8b32b36f1fdfb4ab383c23`. Method for V1–V7: a
throwaway test file `host/verifygate/iter149_row55_throwaway_test.go` calling
`onBlockTriggerKeys(t, path, src)` directly (and `filepath.Glob`/`os.ReadFile` for V7), run
as `go test ./host/verifygate/ -run 'TestIter149Row55' -v`; file deleted afterwards,
porcelain confirmed 0 (V12).

| ID | Claim | Command / input | Observed output |
|----|-------|-----------------|-----------------|
| V1 | Quoted `"on":` false-reds | helper on `"on":\n  push:\n    branches: [dev]\n  workflow_dispatch:\njobs:\n` | `--- FAIL … instrument failure: quoted.yml has no top-level `on:` trigger block` |
| V2 | Flow style false-reds | helper on `on: {push: {branches: [dev]}, workflow_dispatch: }\njobs:\n` | `--- FAIL … instrument failure: flow.yml has no top-level `on:` trigger block` |
| V3 | Tab-indented triggers silently empty the set | helper on `on:\n\tpush:\n\tworkflow_dispatch:\njobs:\n` | helper RETURNS `keys=[]`, no Fatal — the caller would then report `triggers=[] lack workflow_dispatch` (total absence, not a parse limit) |
| V4 | KNOWN-POSITIVE CONTROL, same call: canonical space-indent parses | helper on `on:\n  push:\n    branches: [dev]\n  workflow_dispatch:\njobs:\n` | `keys=[push workflow_dispatch]`, PASS |
| V5 | KNOWN-POSITIVE CONTROL, same path the real gate reads: the real workflow parses | helper on `os.ReadFile(repoRoot/.github/workflows/ci.yml)` | `keys=[push pull_request workflow_dispatch]`, PASS — so V1–V3 are the parser, not a broken instrument. Control scope: the same helper, same call form, on the one file the real gate enumerates. |
| V6 | The scalar arm CASCADES (two messages) | helper on `on:\n  push:\n  workflow_dispatch: garbage\njobs:\n` | ``scalar.yml: `workflow_dispatch:` has scalar value "garbage"; want an empty key or a mapping`` AND returned `keys=[push]` — the caller's `!slices.Contains` then emits the second message. Two messages per one defect. |
| V7 | Glob sees dotfiles AND case-mismatches; a nested dir is a LOUD read error | tempdir with `.hidden.yml`, `CI.YML`, `visible.yml`, `nested/`; `filepath.Glob(dir+"/*")` then `os.ReadFile(dir+"/nested")` | `Glob(dir/*) = [.hidden.yml CI.YML nested visible.yml]`; `os.ReadFile(nested dir) err = read …/nested: is a directory` |
| V8 | go.mod has EXACTLY ONE direct dependency | `cat go.mod` | `require modernc.org/sqlite v1.54.0`; every other require line is `// indirect` (9 lines, all `// indirect`) |
| V9 | The gate-path control at `dispatch_lever_gate_test.go:128` is BLIND to a quoted key only in the DEGENERATE case, and is satisfied INCIDENTALLY in the realistic one | `go run` over three fixtures printing `strings.Contains(src, "on:")` — quoted `"on":` WITH a `runs-on:` line; the same quoted fixture WITHOUT one; and a canonical unquoted file as known-positive control | `quoted+runs-on = true` · `quoted+NO-runs-on = false` · `CONTROL canonical = true`. **CORRECTED BY THE CONTROLLER, iter-149** — this row's first draft read "a fully-quoted workflow Fatals at the control before the helper runs" and was labelled *derived, not separately executed*. Executed, it is FALSE for any workflow carrying `runs-on:`, which is every real job: the substring `on:` inside `runs-on:` satisfies the control for a reason that has nothing to do with the trigger block. So the quoted-key defect surfaces at the HELPER (V1's Fatal) in the common case and at the CONTROL only in the degenerate case. **Two consequences, both load-bearing:** the L128 control is independently weak — it can pass on a file whose trigger block it cannot see — which is a second reason to feed it the shared needle set; and MUT-E's fixture must be the `runs-on:`-free form or the mutant cannot red at the control (see §Mutation table). |
| V10 | The wrong residual prose, doc site | `grep -n "skips dotfiles\|invisible" design_docs/planned/w-ci-recovery-lever-absent.md` | `424:   file (Glob `*` skips dotfiles), a case-mismatched filename, or a nested subdirectory is` / `425:   invisible — …` (block starts at 422) |
| V11 | The wrong cascade prose, doc site | `grep -n "EXACTLY ONE attributed" design_docs/planned/w-ci-recovery-lever-absent.md` | `291:derived set, and EXACTLY ONE attributed message per defect (the `t.Errorf` names the file` (sentence spans 291–293; line 503 separately RECORDS the over-claim as flagged and stays) |
| V12 | Worktree clean after measurement; build green at base | `rm …throwaway_test.go; git status --porcelain \| wc -l; go build ./...` | `0` / `BUILD_OK` |
| V13 | **REPO-WIDE** enumeration of every reader of `.github/workflows/*` — round 1 scoped this row to `host/verifygate/*.go` and so excluded `host/runbook` BY CONSTRUCTION (quorum OBJ-1, scope half, confirmed; re-measured first-party) | `grep -rn 'workflows"' --include='*.go' .` (repo-wide, unbounded) | 7 read sites across THREE files: `host/verifygate/toolchain_pin_gate_test.go:111,201,501` (`TestGoToolchainPinsAgreeAndMatchJobList` reads ci.yml at :111 AND globs `*` at :201 asserting the set is EXACTLY `[ci.yml]`; `TestMiscompileInstrumentStepIsGatedInCI` reads ci.yml at :501); `host/verifygate/ail_binary_gate_test.go:669` (`TestZ3PinDeclaredOnceAndInstalledInBothJobs` reads ci.yml); `host/verifygate/dispatch_lever_gate_test.go:113` (this gate's Glob); `host/runbook/runbook_stageb_test.go:339,361` (OUTSIDE host/verifygate — see §5). Round 1's verifygate-scoped pattern additionally matched non-reader sites — comment lines `dispatch_lever_gate_test.go:109,119`, message strings `ail_binary_gate_test.go:160,289` — which remain non-readers. All 7 read sites appear in the Conflict Surface (§5). |
| V14 | No TestMain / package-level AILANG_BIN gate — a targeted `-run` works without it | `grep -rn "TestMain\|AILANG_BIN is unset" host/verifygate/*.go` | one hit, inside a per-test Fatal at `ail_binary_gate_test.go:42`; the V1–V7 runs above executed without `AILANG_BIN` set |
| V15 | Pinned binary present and correct | `/tmp/ailang-v0300/ailang --version` | `AILANG v0.30.0` / `Commit: e37b370` |
| V16 | `onBlockTriggerKeys` has EXACTLY ONE caller, repo-wide (quorum OBJ-1, reuse half) | `grep -rn "onBlockTriggerKeys" --include='*.go' .` — CONTROL, same command shape and scope: `grep -rn "repoRoot" --include='*.go' . \| wc -l` | 3 hits, ALL in `host/verifygate/dispatch_lever_gate_test.go`: L11 (doc comment), L14 (definition), L131 (the ONE call site). CONTROL `repoRoot` → 86 hits — the grep sees positives at this scope |
| V17 | NO YAML machinery exists in this repo to reuse — the quorum's reuse hypothesis, measured FALSE | `grep -ci yaml go.mod` (CONTROL same file: `grep -ci modernc go.mod`); `grep -rln 'yaml' --include='*.go' .` | `yaml` in go.mod → **0** (control `modernc` → 4); yaml-mentioning Go files repo-wide → **1 file, and it is the gate itself** (`dispatch_lever_gate_test.go`). Derived, not asserted: there is no YAML dependency, no flow-collection scanner, and no parsing helper anywhere in this repo for §3.3(b)'s depth-counting flow parser to overlap or reuse |
| V18 | The Go comment's wrong claims are NOT verbatim copies of the doc's (quorum OBJ-2a), and the base counts for AC6's removal needles | in `host/verifygate/dispatch_lever_gate_test.go`: `grep -c 'skips dotfiles'`; `grep -c 'invisible'`; CONTROL same file: `grep -c 'nested subdirectory'`; then `grep -c 'a hidden file'`; `grep -c 'the Glob is case-sensitive'` | `skips dotfiles` → 0; `invisible` → 0; control `nested subdirectory` → **1** (instrument fires on this file). The comment's own phrasings, base counts for AC6: `a hidden file` → **1**; `the Glob is case-sensitive` → **1** (both at L110) |
| V19 | Casing of "anti-vacuity" is PER-AUTHOR across this repo — a case-sensitive presence needle reads zero on a healthy file (quorum OBJ-2b) | `grep -c 'anti-vacuity floor'` vs `grep -ci 'anti-vacuity'` on `dispatch_lever_gate_test.go`; repo-wide `--include='*.go'` both ways; `grep -ci 'anti-vacuity' design_docs/planned/w-ci-recovery-lever-absent.md` | Go file: case-sensitive → **0**, case-insensitive → **1** (the file writes it UPPERCASE: `// ANTI-VACUITY FLOOR:`, L117). Repo-wide `*.go`: case-sensitive `anti-vacuity floor` → 2, case-insensitive `anti-vacuity` → 22. Row-47 doc: → 9. Round-1 AC6 would have FAILED AT ITS OWN CONTROL |
| V20 | What D-WORLD-30 actually says (quorum OBJ-3 — round 1 cited it with no V-row) | `grep -n 'D-WORLD-30 \| RESOLVED' design_docs/world-mission.md` | 2 lines: **L740, the ledger row itself** (ANSWERED text quoted verbatim in §3.1; attended, Mark Edmondson, 2026-09-01) and L750, an iter-147 STATUS stamp referencing it. The phrase "such an actor can simply delete the test" appears twice in L740 — once in the recommendation, once in the answer |
| V21 | Scalar and sequence `on:` forms Fatal with the MISLEADING "no top-level `on:` block" message (quorum OBJ-3's second half — round 2 asked that this measurement enter the V-row set rather than sit only in the dispositions table) | throwaway test calling `onBlockTriggerKeys` on `on: push` and on `on: [push, workflow_dispatch]`, with the block form as known-positive control in the SAME call and the SAME scope; `go test ./host/verifygate/ -run TestIter149ScalarSeqOnForms -v` at `234d9da`; file deleted after, worktree porcelain restored | `scalar_on_push` -> `instrument failure: … has no top-level `on:` trigger block`; `sequence_on_list` -> the same; `CONTROL_block_form` -> `keys=[push workflow_dispatch]` **PASS**. Controller-measured, iter-149. |
| V22 | The helper ALREADY handles an inline comment on the lever line at base — so MUT-D's control premise is PROVEN rather than assumed (quorum OBJ-2, `gemini-3-1-pro`'s round-2 `proposed_fix`: *"if it greens, the V-row serves as the missing proof for MUT-D's control premise"*) | throwaway test on three fixtures — inline comment, own-line comment, and a no-comment known-positive control — `go test ./host/verifygate/ -run TestIter149InlineComment -v` at `234d9da`; file deleted after, porcelain restored | ALL THREE GREEN: `inline_comment_on_lever` -> `keys=[push workflow_dispatch] contains_lever=true`; `own_line_comment` -> same; `CONTROL_no_comment` -> same. `--- PASS` on all three. **It greens**, so per the reviewer's own disposition this row IS the missing proof, and the `eb215c3` inline-comment fix is confirmed live at base. |
| V23 | "Nothing to reuse" is now COMPLETE BY CONSTRUCTION, not a grep over names (quorum OBJ-1, `gpt5-6-sol`'s round-2 `proposed_fix`, run with ITS OWN named instruments) | (a) `go list -deps ./...` — the full transitive package graph, which cannot miss a parser whose name lacks "yaml"; (b) `git grep` over ALL TRACKED FILES (not just `*.go`) for `.github/workflows`; (c) `git grep -nIE 'func [A-Za-z_]*[Pp]arse' -- '*.go'` enumerating every parse-shaped routine repo-wide | (a) **257** packages in the graph, **0** matching `yaml` case-insensitively; controls in the same call: `modernc` **30**, `encoding/json` **1**, a fresh never-published absent literal **0**. (b) **43** tracked files reference `.github/workflows`, and all but the code sites are design-doc prose; the CODE sites remain the **7** in V13 across three files. (c) **19** parse-shaped funcs across 10 files — `parseStopLine`, `mustParseRefOrZero`, `parseCommandFile`, `ParsePublishApprovalScope`, `parseAndScan`, `parseRequiredRef`, `parseGenesisRef`, `parseRef`, `MustParse`, `Parse` (hashref), `parseDryRun`, `parseJSON`, `parseValue`, plus test funcs. **NONE parses YAML or a flow collection.** The closest is `host/transitionreg/codec.go`'s `parseJSON`/`parseValue`, and it does NOT transfer: JSON has no significant indentation and no block-vs-flow duality, which are the two properties this scan turns on. |

**Controller-verified baselines, rule 3e(a) — cited with provenance; the two full-suite rows
were NOT re-run this session (cost ~175s each) and MUST be re-run in M2 before merge:**

| Baseline | Result at base | Provenance |
|----------|----------------|------------|
| `go build ./...` | rc=0 | re-run this session (V12) |
| `AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/` | rc=0 (~175s) | controller, this iteration, pristine tree |
| `go test ./host/verifygate/` WITHOUT `AILANG_BIN` | **rc=1 AT BASE** — 17 tests fail "AILANG_BIN is unset", a deliberate fail-closed guard | controller, this iteration. Consequence: an AC of the bare form is vacuously red at base; **every test AC below sets `AILANG_BIN=/tmp/ailang-v0300/ailang`** |
| `./scripts/verify_go.sh` | **rc=1 AT BASE on this rig**, for two reasons owned elsewhere: (i) the known FLAKE `TestHandlerTimeoutKillsTheWholeProcessGroup` (row 58, recorded FLAKY not deterministic); (ii) the row-54 fleet arm reds while the driver copy is genuinely stale (759 diff-lines at last controller measure — "the fleet must commit", never "absorb it") | controller, this iteration. Consequence: **no AC in this doc has the form "verify_go.sh rc=0"**; whole-suite claims must be set comparisons against a recorded base failing set per row 58's amendment |

---

## 3. Decision

### 3.1 LINE SCAN, hardened — not a YAML parser

The row was filed because closing it "needs a decision this item does not own — whether the
gate adopts a structural YAML parser instead of a line scan" (charter row 55).
**`D-WORLD-30` (attended ruling, Mark Edmondson, 2026-09-01) answered exactly that question
for a SIBLING gate — row 52's CI step-scoping gate — and chose A: LINE SCAN, hardened.** The
ledger row (V20 — `design_docs/world-mission.md:740`) answers, verbatim:

> "ANSWERED — A — LINE SCAN, hardened, deriving the anchor from the SHALLOWEST enclosing
> `steps:` rather than the nearest one (measured to catch the round-2 attack that the doc's
> nearest-anchor derivation is blind to). This gate is a repo-internal tripwire against our
> own accidental regressions, not an adversarial boundary: the threat model that would
> justify B is an actor who can already edit `ci.yml`, and such an actor can simply delete
> the test. Not worth adding the second direct dependency to a `go.mod` that has exactly one.
> The residual is accepted and named: the scan stays text-level, and a future re-indent of
> `ci.yml` makes the gate fatal until the constant is updated."

Round 1 cited this ruling with no V-row and treated its rationale as transferring whole
(quorum OBJ-3 — the reviewer's point is right). Split explicitly:

- **TRANSFERABLE — properties of this repo, not of row 52's gate:** the repo-internal-
  tripwire threat model; "an actor who can already edit `ci.yml` can simply delete the test";
  and the one-direct-dependency `go.mod` constraint (V8, re-verified: exactly ONE direct
  require, `modernc.org/sqlite v1.54.0`; every other require line is `// indirect`).
- **NOT TRANSFERABLE — row-52 mechanics with no analogue in row 55:** "deriving the anchor
  from the SHALLOWEST enclosing `steps:`", and the accepted residual that a future re-indent
  of `ci.yml` makes the gate fatal until a constant is updated. Row 55's gate has no step
  anchor and no indent constant.

**Scope of the "nothing to reuse" claim, narrowed to what was measured (this is `gpt5-6-sol`'s own round-2 alternative — *"either reuse identified machinery or narrow the claim to the measured scope"* — and the measurements are V16/V17/V23):** the claim is that the FULL transitive package graph of this module contains no YAML package (0 of 257, controls firing), and that no parse-shaped routine anywhere in the tracked Go sources parses YAML or a flow collection (19 enumerated by name). It is NOT a claim about machinery reachable only through build tags, generated code, or a package this module does not import. Within that stated scope the enumeration is complete by construction rather than by inspection, which is what the objection asked for.

**Scope honesty: D-WORLD-30 does NOT resolve row 55** — it names row 52's gate. What carries
the line-scan choice here is the ruling's TRANSFERABLE half plus this doc's own measurements:
all three measured shapes are closable within a line scan (§3.2–§3.4); there is NO YAML
machinery in this repo to reuse (V17 — `yaml` greps 0 in go.mod, and the only yaml-mentioning
Go file repo-wide IS this gate); and the helper has exactly ONE caller repo-wide (V16), so
hardening it in place overlaps nothing. **Line scan, hardened, no new dependency, no go.mod
change.** No new decision needs raising: nothing here requires the structural parser that
would trigger one.

### 3.2 Split the helper: pure parse core, test-verdict wrapper

`onBlockTriggerKeys` currently interleaves parsing with `t.Fatalf`/`t.Errorf`, which is why
its failure arms were only ever measurable by throwaway harnesses (V1–V3). Refactor into:

```go
// pure: no *testing.T, no I/O — unit-testable on every arm, including failures
func parseOnBlockTriggers(src string) (keys []string, scalarValued map[string]string, err error)

// wrapper: preserves TODAY'S messages byte-for-byte for the existing arms
func onBlockTriggerKeys(t *testing.T, path, src string) []string
```

Typed sentinel errors: `errNoOnBlock`, `errDuplicateOnBlock`, `errTabIndent` (new),
`errUnhandledOnForm` (new). The wrapper maps `errNoOnBlock`/`errDuplicateOnBlock` to the
EXISTING Fatalf strings unchanged (so row-47's mutation evidence for those arms stays valid),
maps the two new errors to new honest messages, and emits the existing scalar-value `Errorf`
per `scalarValued` entry. This matches S1's spirit (pure core, effects at the boundary)
applied to Go test instruments, and is what makes M2's mutation table executable in-process
instead of via subprocess capture.

A single shared matcher feeds BOTH the block finder and the duplicate counter, and a shared
accepted-forms needle set feeds the gate's L128 known-positive control, so the three sites
cannot drift (V9 is what happens when they can):

```go
// matchTopLevelOn reports whether a line declares the top-level trigger key, and returns
// the remainder after the colon. Column-0-anchored, exact accepted forms only — preserves
// P14's rationale: a trimmed comparison would anchor on ANY nested `on:`.
func matchTopLevelOn(line string) (rest string, ok bool)   // accepts on: / "on": / 'on':
```

### 3.3 Per-shape dispositions

| Shape | Today | After | Why |
|-------|-------|-------|-----|
| (a) `"on":` / `'on':` quoted key | Fatal "no top-level on: block" (V1), or earlier at the L128 control (V9) | **GREEN** — quoted forms accepted at column 0 by `matchTopLevelOn`; L128 control broadened to the same needle set | valid YAML declaring the lever; the standard footgun remedy must not red CI |
| (b) `on: {…}` flow mapping | identical Fatal (V2) | **GREEN** — when `rest` (comment-stripped) is `{…}`, split on top-level commas with brace/bracket depth counting; key = text before first `:` per item, quotes trimmed; the scalar-value rule applies inside flow (`workflow_dispatch: garbage` in flow is a violation, `workflow_dispatch: ` and `workflow_dispatch: {}` are not). **GUARD SPECIFICATION — `oc-glm-5-2`'s round-2 `proposed_fix`, applied VERBATIM under the narrow-refinement carve-out:** *"The flow-mapping arm activates ONLY when `rest` (comment-stripped, whitespace-trimmed) both starts with `{` and ends with `}`. Any other non-empty `rest` — including a multi-line flow mapping whose first line yields an unclosed `{` — falls through to `errUnhandledOnForm` and Fatals loudly. This closes the silent-fallback path: a partial or unclosed flow never produces a key set."* This also discharges the second half of `gpt5-6-sol`'s round-2 `proposed_fix` (*"requiring an explicit typed error rather than partial interpretation"*) — the two reviewers converge on the same disposition, so quoted commas/braces/colons, escapes, nested collections and malformed or unbalanced input all reach `errUnhandledOnForm` rather than a partial key set. Multi-line flow mappings are therefore a DECLARED REFUSAL, not a silent miss (Residual R1). | valid YAML declaring the lever |
| (c) tab-indented trigger line | SILENT `keys=[]` → caller misreports total absence (V3) | **LOUD, TYPED** — any `\t` in a scanned line's indentation prefix returns `errTabIndent`; wrapper Fatals `instrument failure: <path> line N: tab in indentation — this line scan computes depth in spaces only`. NOT a green. | YAML 1.2 §6.1 defines block indentation as spaces only (spec citation, NOT a measurement — no YAML implementation was run this session; pyyaml is absent from this rig's python3). The fix does not depend on the citation: honest-loud beats misreporting absence either way, and if the file IS invalid YAML, green would be wrong. |
| `on: push` / `on: [push, …]` scalar & sequence forms | Fatal with the MISLEADING "no top-level on: block" message — same `l == "on:"` guard branch as V2. **MEASURED BY THE CONTROLLER, iter-149** (the row shipped as *inference … not separately executed*; executed, the inference HOLDS): a throwaway test calling `onBlockTriggerKeys` on `on: push` and on `on: [push, workflow_dispatch]` returns `instrument failure: <f> has no top-level `on:` trigger block` for BOTH, while the known-positive control in the same call — the block form, same helper, same scope — returns `keys=[push workflow_dispatch]` and PASSES | **LOUD, TYPED** — non-empty `rest` that is not a flow mapping returns `errUnhandledOnForm`; wrapper Fatals naming the form and this scan's limit, no longer claiming the block is absent | out of the row's measured set and unused in this repo (V5: ci.yml is block-form); an honest refusal is the 0.2d-proportionate disposition. Declared Residual R1. |
| (d) scalar-arm cascade | two messages (V6) | **KEPT, deliberately** — the cascade is defense-in-depth: neutering the scalar Errorf alone leaves the caller's absent-lever red; "fixing" it by appending the key after the Errorf would make a neutered Errorf a SILENT FALSE GREEN, i.e. would manufacture the soundness hole this row is fenced against | fence 1; only the row-47 doc's SENTENCE is wrong, so only the sentence changes (M3) |

### 3.4 Prose corrections (measured-wrong claims only)

Three sites, one milestone (M3):

1. `design_docs/planned/w-ci-recovery-lever-absent.md:291-293` — replace "EXACTLY ONE
   attributed message per defect (… it never cascades — P5's precedent)" with the measured
   truth: the scalar arm emits TWO messages (V6), deliberately, as defense-in-depth; the
   absent-lever arm alone is one message. Line 503 (which RECORDS that the over-claim was
   flagged) stays — it is a true historical record.
2. `design_docs/planned/w-ci-recovery-lever-absent.md:422-427` (Residual 3) — rewrite per V7:
   Glob `*` DOES enumerate dotfiles and case-mismatched names; a nested subdirectory is a
   LOUD `t.Fatal` ("is a directory"), not invisible. What remains genuinely invisible: a
   workflow outside `.github/workflows/` entirely (singular `.github/workflow` dir, a
   root-level `.yaml`) — inference from the Glob path literal at
   `dispatch_lever_gate_test.go:113`, stated as such.
3. `host/verifygate/dispatch_lever_gate_test.go:110-111` — the code comment carrying the
   same wrong claims in its OWN words, not the doc's (V18; quoted in §1(e)): "cannot see …
   a hidden file" is false (V7 — Glob returns `.hidden.yml`); "the Glob is case-sensitive"
   is irrelevant to a bare `*` pattern (V7 — `CI.YML` IS enumerated); and the
   nested-subdirectory case is a loud "is a directory" Fatal, not unseen. Rewrite those
   clauses in the comment's own words; same commit.

---

## 4. Milestones

Total **~0.2d**. All test ACs set `AILANG_BIN=/tmp/ailang-v0300/ailang` (base rc=1 without it
— Verification Log baselines). No AC has the form "verify_go.sh rc=0" (red at base for two
reasons owned by rows 58 and 54).

### M1 (0.10d) — pure-core refactor + the three hardenings

Deliverables: `parseOnBlockTriggers` + `matchTopLevelOn` + shared needle set for the L128
control; wrapper preserving existing messages byte-for-byte on existing arms; a table-driven
`TestOnBlockTriggerParserShapes` in `dispatch_lever_gate_test.go` covering, minimum: quoted
key → `[push workflow_dispatch]`; flow style → `[push workflow_dispatch]`; flow with scalar
violation → violation entry present; tab indent → `errTabIndent`; scalar/sequence `on:` forms
→ `errUnhandledOnForm`; duplicate mixed `on:`+`"on":` → `errDuplicateOnBlock`; canonical
block form → `[push workflow_dispatch]` (V4's control, now permanent); no `on:` at all →
`errNoOnBlock`; inline-comment value (the eb215c3 regression) → NOT a violation.

- **AC1:** `AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -run
  'TestOnBlockTriggerParserShapes|TestEveryWorkflowDeclaresDispatchLever' -v` → rc=0, and the
  `-v` output lists BOTH test names as `--- PASS` (guards against `-run` matching nothing —
  rc=0 on zero matched tests would otherwise pass this AC vacuously; the explicit name check
  is what can fail).
  *What would still pass if the claim were false:* a refactor that handles none of the three
  shapes passes the GATE test (ci.yml is canonical, V5) — it is the table cases in
  `TestOnBlockTriggerParserShapes` that red; a missing table case for a shape passes both.
  Hence AC2.
- **AC2 (case-presence floor):** `grep -c 'errTabIndent\|errUnhandledOnForm\|"on":'
  host/verifygate/dispatch_lever_gate_test.go` ≥ 3 hits AND the M2 mutation run (each mutant
  reds a NAMED table case — presence of a case that can fire, not just presence of code).
  *What would still pass if false:* grep proves text, not behavior — which is exactly why AC2
  is completed by M2's mutants rather than standing alone.
- **AC3 (no-dependency fence):** `git diff --stat -- go.mod go.sum` after M1 → empty.
  *What would still pass if false:* nothing — any dependency addition touches go.mod.

### M2 (0.05d) — mutation run + full-package regression

Run the mutation table below; then re-run the controller's full-package baseline:
`AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/` → rc=0 (this is the re-run
the Verification Log defers; ~175s).

- **AC4:** every mutant builds (`go build ./...` rc=0 on the mutated tree — a non-building
  mutant proves nothing) and reds EXACTLY its expected red set, with the canonical-control
  case staying green under every mutant (proves each mutant broke only its target arm, not
  the instrument).
  *What would still pass if false:* a mutant redding EVERYTHING (including the control) would
  pass a naive "mutant reds" check — the control-stays-green clause is what fails it.
- **AC5:** `AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/` rc=0 on the final
  (unmutated) tree.
  *What would still pass if false:* a hang under 10m would not (go test's default timeout
  reds); a test DELETED by the diff would pass — AC1's explicit `--- PASS` name listing covers
  the two tests this doc owns.

**Mutation table** (neutering via `if false && …` preferred over deletion so imports stay
used; each row names the mutant, the edit, and the expected red set):

| Mutant | Edit (neutered, must still build) | Expected RED set | Expected GREEN control |
|--------|-----------------------------------|------------------|------------------------|
| MUT-A quoted-form | `matchTopLevelOn` accepts only bare `on:` (quoted needles behind `if false &&`) | table cases: quoted key (gets `errNoOnBlock`, wants keys); duplicate mixed (gets 1 count, wants `errDuplicateOnBlock`) | canonical block form |
| MUT-B flow-parse | flow-mapping branch behind `if false &&` (falls through to `errUnhandledOnForm`) | table cases: flow style; flow-with-scalar-violation | canonical block form |
| MUT-C tab-detect | indentation `\t` check behind `if false &&` | table case: tab indent (gets `keys=[]` nil err, wants `errTabIndent`) | canonical block form |
| MUT-D scalar-value | the value-judging branch behind `if false &&` (key appended, no violation recorded) | table cases: flow-with-scalar-violation; block-form scalar (`workflow_dispatch: garbage` wants a violation entry) | inline-comment value case (must STAY non-violating — proves MUT-D's red is the check vanishing, not the comment-strip) |
| MUT-E control-drift | the gate's L128 needle set narrowed back to `on:` only (shared-slice mutant) | table case: quoted key via the GATE-path wrapper arm, using the **`runs-on:`-FREE fixture** — the mutant reds at the control only for that form, because a fixture carrying `runs-on:` satisfies `strings.Contains(src, "on:")` incidentally and the mutant then walks straight through (V9, measured). This is the mutant that proves the shared needle set actually feeds L128, so its fixture choice is load-bearing, not cosmetic. | real-ci.yml control case (V5's, now permanent), which MUST stay green — it carries `runs-on:` |

### M3 (0.05d) — the three prose corrections

Edit the three sites of §3.4 in one commit.

- **AC6 (absence of the measured-wrong claims, per site — every removal needle has a
  measured base count ≥1 in ITS OWN file, so each zero can fail):** doc site:
  `grep -c 'skips dotfiles' design_docs/planned/w-ci-recovery-lever-absent.md` → 0 (base 1,
  V10). Go site: `grep -c 'a hidden file' host/verifygate/dispatch_lever_gate_test.go` → 0
  AND `grep -c 'the Glob is case-sensitive'` (same file) → 0 (base 1 each, V18). Round 1
  aimed the DOC's phrases at the Go file, where they grep 0 at BASE (V18) — a vacuous check;
  the comment's wrong claims are in its own words, so the needles must be too (quorum
  OBJ-2a). Known-positive control in the same call, scoped to the same file-set:
  `grep -ci 'anti-vacuity' <file>` ≥ 1 in EACH of the two files, read as PRESENCE (≥1),
  never as an exact count (doc: 9 case-insensitive hits; Go file: 1, written
  `// ANTI-VACUITY FLOOR:` at L117 — V19). Why `-ci`: casing of this jargon is per-author
  across the repo (2 case-sensitive vs 22 case-insensitive hits repo-wide in `*.go`, V19),
  so a case-sensitive presence needle is an instrument that reads zero on a healthy file —
  the same defect class this loop already recorded for its own STATUS-stamp tell. Round 1's
  control was case-sensitive lowercase `anti-vacuity floor`, which greps 0 in the Go file:
  as written, AC6 would have FAILED AT ITS OWN CONTROL (quorum OBJ-2b, confirmed live).
  *What would still pass if false:* a rewording that keeps the wrong claim without those
  literal phrases ("Glob omits hidden files") — hence AC7.
- **AC7 (human-readable assertion, evaluator-checkable):** the rewritten Residual-3 block and
  code comment each contain BOTH corrected facts with their V7 citation: dotfiles/case
  mismatches ARE enumerated; a nested subdirectory is a LOUD "is a directory" Fatal.
  Checkable: `grep -ci "is a directory" <both files>` ≥1 each (`-ci` per the needle audit
  below — presence needles are case-insensitive so a future re-caser cannot make this AC
  read zero on a healthy file; the phrase itself quotes the OS error V7 measured).
  *What would still pass if false:* the phrase present but negated — residual accepted; the
  evaluator reads the sentence (this is a prose AC; the narrowest mechanical gate is the
  phrase floor, and the sentence-level truth is scored, not grepped).
- **AC8:** `grep -c "EXACTLY ONE attributed message"
  design_docs/planned/w-ci-recovery-lever-absent.md` → 0, control `grep -c "never cascades"`
  → exactly 1 (line 503's historical record survives; line 292's claim is gone — V11
  establishes the base counts as 1 and 2 respectively).
  *What would still pass if false:* deleting line 503 too would give 0/0 — the `exactly 1`
  on the control is what fails that.

**Needle audit — quorum OBJ-2b names the CLASS (a case-sensitive or verbatim-string needle
aimed at source text), not just the AC6 instance. Every grep needle in this doc's ACs and
mutation table, audited:**

- **AC6's presence control** — WAS the defect; fixed above (`-ci`, read as presence ≥1).
- **AC7's presence needle** — targets text M3 itself writes, but made `-ci` anyway so a
  future re-casing cannot turn the AC vacuous.
- **AC6's and AC8's ZERO-needles** — verbatim and case-sensitive ON PURPOSE, and sound:
  each targets specific existing bytes whose base count is measured ≥1 in the SAME file it
  is aimed at (V10, V11, V18), so the grep can fail, and casing drift cannot produce a false
  zero on a healthy file because a zero IS the pass condition. The failure mode of a
  case-sensitive needle is a false ZERO — fatal for a presence check, harmless for a
  measured-base absence check.
- **AC2's needles** (`errTabIndent`, `errUnhandledOnForm`, `"on":`) — Go identifiers and a
  literal this sprint itself introduces; identifier casing is enforced by the compiler, not
  by an author's prose style, so case-sensitive is correct there.
- **AC1's `--- PASS`** — go test's own machine output format, not source text.
- **The mutation table** contains no grep needles; expected-red sets are named table cases.

---

## 5. Conflict Surface

This helper is a parser over `.github/workflows/*`; every reader of that path **REPO-WIDE**
(complete set, V13 — repo-wide unbounded grep: 7 read sites, THREE files; round 1 scoped this
table to `host/verifygate` and missed `host/runbook` by construction — quorum OBJ-1):

| Consumer | What it does | Effect of this change |
|----------|--------------|-----------------------|
| `dispatch_lever_gate_test.go:113` (`TestEveryWorkflowDeclaresDispatchLever`) | the gate this doc hardens | behavior on the repo's actual workflow set is UNCHANGED (ci.yml is canonical block-form, V5); only non-canonical inputs change disposition |
| `toolchain_pin_gate_test.go:201` (`TestGoToolchainPinsAgreeAndMatchJobList`) | Globs `*` and asserts the workflow set is EXACTLY `[ci.yml]` | untouched. This doc adds NO workflow files (M1's table cases are inline strings and `t.TempDir()` fixtures, never files under `.github/workflows/`) — so its `[ci.yml]`-exactly assertion cannot fire from this diff |
| `toolchain_pin_gate_test.go:111` (same test) + `:501` (`TestMiscompileInstrumentStepIsGatedInCI`) + `ail_binary_gate_test.go:669` (`TestZ3PinDeclaredOnceAndInstalledInBothJobs`) | read `ci.yml` directly with their OWN line scans | untouched — they do not call `onBlockTriggerKeys` and this doc does not modify `ci.yml`. Their own line-scan limits (flow-style `with:`, block-scalar smuggling) are theirs, disclosed in their own comments, and NOT in this row's scope |
| `host/runbook/runbook_stageb_test.go:339,361` — OUTSIDE `host/verifygate` | appends `ci.yml` to a scan-target list and substring-scans its raw lines for `world-publish` (:339), then re-reads it counting `verify_go.sh` as its own known-positive control (:361); a byte-level scan, no trigger-block parsing | untouched — it never calls `onBlockTriggerKeys` (V16: the helper's ONE call site repo-wide is `dispatch_lever_gate_test.go:131`, and Go test-file helpers are not importable across packages anyway), and this doc does not modify `ci.yml`, so both its scan and its control read the same bytes before and after this change |
| `ail_binary_gate_test.go:160,289` | error-message strings naming the path | not readers; unaffected |

Frozen core: `tools/launchd/*` untouched. `go.mod`/`go.sum` untouched (AC3). No `.ail`
changes, so `scripts/verify_ail.sh` and the pinned-binary surface are untouched (V15 records
the pin for the test-AC env var only).

## 5b. Refer to the queue

The revision measurements surfaced NO pre-existing defect outside this row's scope that
needs a queue row: the `host/runbook` read sites are healthy byte-scans carrying their own
known-positive controls (§5), and the reuse question resolved to "nothing exists to reuse"
(V17), which removes work rather than adding it. Nothing is filed.

## 6. Declared Residuals

Per the hard gate: each residual below is measured or explicitly labeled as inference —
none is asserted on nobody's measurement (the row-47 mistake this row corrects).

- **R1 — scalar/sequence `on:` forms still red, now honestly.** `on: push` and
  `on: [push, workflow_dispatch]` are valid lever-adjacent YAML this scan still refuses —
  loudly, via `errUnhandledOnForm`, with a message that names the scan's limit instead of
  falsely claiming the block is absent. Basis: same guard branch as V2 (measured), form-level
  behavior inferred from source, labeled as such in §3.3. If a repo workflow ever legitimately
  adopts these forms, THAT sprint extends the parser or raises the parser question anew.
- **R2 — a line scan cannot see semantics.** Unchanged from row 47 and still true: DECLARED
  is not RUNNABLE (no dispatch run is created or proven green); the lever re-verifies a named
  ref's TIP, not an arbitrary SHA; a step-level `if:` can disable everything at runtime.
  These sentences already exist in the gate's comment (L92–106) and are not weakened here.
- **R3 — enumeration blind spots, now stated correctly.** After M3: a workflow outside
  `.github/workflows/` (singular dir, root-level `.yaml`) is invisible — inference from the
  Glob path literal, labeled as such. Dotfiles and case-mismatched names ARE seen (V7,
  measured); a nested subdirectory is a LOUD Fatal (V7, measured).
- **R4 — the tab disposition rests partly on a spec citation.** YAML 1.2 §6.1
  (spaces-only indentation) was cited, not executed — no YAML implementation ran this
  session. The design is chosen so nothing breaks if the citation were wrong: a loud typed
  refusal is honest regardless, and a future measurement can upgrade it.
- **R5 — the cascade stays.** Two messages per scalar-valued lever (V6), on purpose
  (§3.3 row d). Anyone reading the gate's output should expect both.

## Quorum verification log

| Round | Artifact | Verdict | Reviewers | Cost |
|---|---|---|---|---|
| 1 | `.ailang/state/mission-quorum/w-dispatch-lever-parser-false-reds-on-valid-yaml-2026-09-02T06-55-28Z.json` | **BLOCKED** | `gpt5-6-sol` reject · `gemini-3-1-pro` reject · `oc-glm-5-2` reject · controller pass | $0.0936 |
| 2 | `…-2026-09-02T07-04-17Z.json` | **BLOCKED** | `gpt5-6-sol` reject · `gemini-3-1-pro` reject · `oc-glm-5-2` reject · controller pass | $0.1275 |

Both rounds ran at **FULL STRENGTH**: `.synthesis.absent_reviewers` = `[]` in both, cross-checked
against `[.reviewers[]|select(.present==false)]` = `[]`, with `has("synthesis")` = `true` as the
control that the path resolves at all. No reviewer was waived and no verdict was read through a
`null`.

**Round 1 → revision.** All three objections were *premise* objections, so per rule 3f the
controller RAN each rather than forwarding it, and handed the designer measurements:
`gpt5-6-sol`'s reuse hypothesis was **refuted** (V16/V17) while its *scope* complaint was
**correct** and produced V13's repo-wide rewrite plus `host/runbook/runbook_stageb_test.go` in the
Conflict Surface; `gemini-3-1-pro` was **right on both halves**, and its (b) was a live defect in
this doc's own AC6, whose known-positive control would have failed on a casing mismatch (V19);
`oc-glm-5-2` was **right** that §3.1 cited `D-WORLD-30` with no V-row (now V20, which also splits
the transferable rationale from the row-52-specific mechanics).

**Round 2 → NARROW-REFINEMENT CARVE-OUT (ratified for this mission at iter-13), applied by the
controller.** The gate's conditions were checked before use and both hold: every remaining
blocking objection carries a **concrete reviewer-authored `proposed_fix`**, and **none disputes the
design DIRECTION** — line scan / hardened / no new dependency was unchallenged by all three, and
what they asked for is completeness (an audit), a missing baseline measurement, and determinism (a
flow-parser guard). The fixes were applied as the reviewers wrote them, not as the controller would
have resolved them:

- `oc-glm-5-2` → its guard specification pasted **VERBATIM** into §3.3 shape (b); plus V21, the
  V-row it asked for.
- `gemini-3-1-pro` → V22, run exactly as its `proposed_fix` prescribes. It **greens**, so by the
  reviewer's own stated disposition the row IS the missing proof for MUT-D's control premise.
- `gpt5-6-sol` → V23, run with **its own named instruments** (`go list -deps ./...`, `git grep`
  over all tracked files, a repo-wide parse-routine enumeration), and §3.1's claim **narrowed to
  the measured scope**, which is the alternative arm the reviewer itself offered. Its second half
  (a typed error rather than partial interpretation) is discharged by the same guard specification
  — the two reviewers converge on one disposition.

This SATISFIES the objections; it is not a force-pass. No objection was overridden, and no new
decision was raised — the carve-out exists precisely so a doc whose reviewers have written down
the fix is not parked for a human who has nothing to decide.
