# w-ci-recovery-lever-absent — a dropped webhook to `dev` is permanently unverifiable because this repo declares no API-driven trigger; a one-line `workflow_dispatch:` lever plus a static gate that pins EVERY workflow file to it

**Status**: PLANNED
**Date**: 2026-08-28
**Queue item**: 47, `w-ci-recovery-lever-absent` (clause-2). The mission has paid for this
defect twice (iterations 128 and 129); this is the closing row.
**Estimated**: ~0.10 day — ONE keyed line in `.github/workflows/ci.yml` (the `on:` block) and
ONE new test function (and one small parser helper) in `host/verifygate`. No `.ail`, no
`tools/launchd/*`, no other workflow file touched. Single-gate (Go) change.
**Designer**: `design-doc-creator` (DESIGN pass, iteration 135)
**Toolchain boundary**: every `VERIFIED-BY-ME` command below ran first-party in this checkout
(`/Users/voightkampff/dev/sunholo-data/.wt-world-iter136`, zsh with
`export PATH=/opt/homebrew/bin:$PATH`, darwin/arm64) at `3a98c79`, tree clean
(`git status --porcelain` = 0 lines, P11), 2026-08-28, with `go version` = `go1.26.6
darwin/arm64`. Rows marked INHERITED were measured by the controller this session and are
cited with attribution; the load-bearing non-mutation arms are re-derived here first-party.
The pinned released AILANG binary used to run the gate package is
`/tmp/ailang-v0300/ailang` (`AILANG v0.30.0`, P6).

> **Thesis:** a `push` to `dev` whose webhook GitHub silently drops is **permanently
> unverifiable** — the delivery is never replayed and no workflow run is ever created for
> that commit (measured twice by this mission: iteration 128's record push `a0b3162` still
> read `checks=0` / `actions/runs?head_sha=<40-char>` `total=0` a day later, with a
> rev-parsed control commit returning `checks=2`/`runs=1` in the same call). On a PR there
> is a workaround (a tree-identical empty commit through the git API fires a genuine
> `pull_request: synchronize`), but **on `dev` itself there is no lever at all**: advancing
> `dev` changes the very commit being verified, and `workflow_dispatch` — the one
> API-driven trigger that would re-run verification on the tip of a named ref (this repo's
> `dev`, P17) without moving `dev` — is not declared. This doc adds that lever (P2, P3) and pins it with a static
> gate in the existing `host/verifygate` family so it cannot be deleted silently (the
> removal-shaped and addition-shaped mutations, M1–M9). The claim is scoped **exactly** to
> what the lever buys: `workflow_dispatch` buys **A VERDICT ON THE TIP OF A NAMED REF**
> (`--ref dev` — P16, P17), never on an arbitrary commit SHA (residual 7); it does **NOT**
> buy **A MERGEABLE PR** — a dispatch run's checks do not satisfy PR branch protection
> (measured by this mission previously), a limitation stated in the code comment and in
> Declared residual 1, because an over-broad claim is the exact defect the sibling
> pin-tests were scoped to avoid.

## The finding in one paragraph

`.github/workflows/` holds exactly ONE file, `ci.yml` (P1), whose `on:` block declares only
`push: branches: [dev]` and `pull_request:` (P3) — there is no API-driven trigger
(P2: `grep -c 'workflow_dispatch'` → 0, with the same-call known-positive controls
`pull_request` → 1 and `^  push:` → 1 proving the read is live, plus a never-used absent
literal → 0 in the same breath). So when GitHub drops a webhook, the affected `dev` commit
can never be re-verified: there is no `workflow_dispatch`, the only workaround (an empty
commit through the git API) only works on a PR, and advancing `dev` moves the very commit at
stake. The repair is two parts, both small: add `workflow_dispatch:` as a trigger key at the
same depth as `push:`/`pull_request:` (P3), and add a static consistency test
`TestEveryWorkflowDeclaresDispatchLever` to `host/verifygate/dispatch_lever_gate_test.go` in the idiom of its siblings
`TestZ3PinDeclaredOnceAndInstalledInBothJobs` (`ail_binary_gate_test.go:668`) and
`TestGoToolchainPinsAgreeAndMatchJobList` (`toolchain_pin_gate_test.go:106`), asserting that
EVERY enumerated workflow file declares the lever in its `on:` block — with an anti-vacuity
floor (an empty enumeration FAILS LOUDLY), known-positive controls, and a single attributed
message per defect. The lever's limit is declared, not hidden: a dispatch run is not
equivalent to the event it replaces, so the gate's claim is "the lever is declared", never
"CI can always be re-triggered" or "this unblocks merges" (residual 1).

## Premises

`VERIFIED-BY-ME` = run by the designer this session, output observed first-hand. `INHERITED`
= measured by the controller this session (cited with attribution).

| # | Claim | Command (verbatim) | Observed | Status |
|---|---|---|---|---|
| P1 | `.github/workflows/` holds exactly ONE workflow file, named `ci.yml`, and that scope exists | `ls .github/workflows/ \| wc -l`; `ls .github/workflows/`; `test -f .github/workflows/ci.yml` | `1`; `ci.yml`; `true` (rc=0) — scope asserted, not read empty | VERIFIED-BY-ME |
| P2 | `ci.yml` declares NO `workflow_dispatch` — and the read is live | `grep -c 'workflow_dispatch' .github/workflows/ci.yml`; same-call controls `grep -c 'pull_request' …` and `grep -cE '^  push:' …`; absent-literal control `grep -c 'schedule:' …` | `workflow_dispatch` → **0, rc=1**; `pull_request` → **1, rc=0**; `^  push:` → **1, rc=0**; `schedule:` → **0, rc=1** — the 1/0 rc pair proves grep's rc=1 is a legitimate zero, and the absent-literal control proves the scanner sees absences | VERIFIED-BY-ME |
| P3 | `ci.yml`'s `on:` block is exactly `push: branches:[dev]` + `pull_request:`, with `on:` at lead 0 and triggers at lead 2 | `sed -n '/^on:/,/^env:/p' .github/workflows/ci.yml`; `grep -nE 'on:\|push:\|pull_request\|workflow_dispatch\|branches' .github/workflows/ci.yml`; awk leading-space count on lines 3–8 | block reads `on:` / `  push:` / `    branches: [dev]` / `  pull_request:`; `on:` at `:3`, `push:` at `:4`, `pull_request:` at `:6`, `jobs:` at `:16`; lead counts `on:`=0, `push:`=2, `pull_request:`=2 | VERIFIED-BY-ME — `workflow_dispatch:` belongs at lead 2, between `:6` and the `env:` comment |
| P4 | The sibling pin-tests that the new gate must match exist at the cited locations | `grep -n "func TestZ3PinDeclaredOnceAndInstalledInBothJobs\|func TestGoToolchainPinsAgreeAndMatchJobList" host/verifygate/*_test.go` | `ail_binary_gate_test.go:668` and `toolchain_pin_gate_test.go:106` | VERIFIED-BY-ME |
| P5 | The sibling pin-test ALREADY enumerates workflow files with the `filepath.Glob(…, "workflows", "*")` idiom at `toolchain_pin_gate_test.go:197–207` and asserts the set is EXACTLY `[ci.yml]` — my new gate's enumerator and anti-vacuity floor copy this precedent | `grep -n "filepath.Glob\|workflowFiles\|want exactly \[ci.yml\]" host/verifygate/toolchain_pin_gate_test.go` | `:197` `filepath.Glob(filepath.Join(repoRoot, ".github", "workflows", "*"))`; `:201–205` build + `slices.Sort` the base names; `:206` `slices.Equal(workflowFiles, []string{"ci.yml"})`; `:207` message | VERIFIED-BY-ME — precedent enumerator; its `want exactly [ci.yml]` is itself a Conflict Surface (see §Conflict Surface) |
| P6 | The gate package is RED at base without `AILANG_BIN` and GREEN with it — the two arms differ, so the red is the environment | `go test ./host/verifygate/ -count=1` (no env); then `AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -count=1`; `test -f /tmp/ailang-v0300/ailang && /tmp/ailang-v0300/ailang --version` | no-env: **rc=1**, **17 `--- FAIL`** lines, **17** `AILANG_BIN is unset`; with env: **rc=0**, **0 FAIL**, `ok … 46.882s`; binary exists, `AILANG v0.30.0` | VERIFIED-BY-ME (re-derives controller measurement, identical) — therefore EVERY acceptance criterion running this package carries the `AILANG_BIN=/tmp/ailang-v0300/ailang` prefix, or it is red at base and vacuous |
| P7 | Base hygiene of the repo toolchain gates (unchanged by this design) | `gofmt -l .`; `go vet ./...`; `go build ./...` | `gofmt -l .` rc=0, **0 bytes**; `go vet ./...` rc=0; `go build ./...` rc=0 | VERIFIED-BY-ME |
| P8 | `go build ./...` does NOT compile `_test.go` files; `go vet` on the affected package (here `./host/verifygate/`) is the typecheck for test-file mutants | controller, this session | a test-file undefined symbol gives `go build ./...` rc=0 and `go vet` rc=1 | INHERITED (controller) — binding for any mutation table row that edits a `_test.go` file |
| P9 | No `actionlint` (or any external YAML linter) runs anywhere in this repo, so nothing else validates the `on:`-block shape | `grep -rln 'actionlint' --include='*.go' --include='*.sh' --include='*.yml' --include='*.yaml' . --exclude-dir=.git` | 2 hits, both COMMENT lines in `host/verifygate/*_test.go` declaring "no actionlint runs anywhere in this repo" — the repo's own word; no linter backs the lever's shape | VERIFIED-BY-ME (corroborates the siblings' own V18 record) |
| P10 | `verify_go.sh` (the gated CI `go test ./...`) refuses to run without `AILANG_BIN`, so the whole-repo gate is not silently false-green | `grep -n 'AILANG_BIN' scripts/verify_go.sh` | `:120` `if [ -z "${AILANG_BIN:-}" ]`, `:121` `✗ AILANG_BIN is unset … false-green`, `:132/:144` version asserts; the gated job runs `./scripts/verify_go.sh` (ci.yml:165) | VERIFIED-BY-ME — an AC that uses the whole-repo gate inherits the same prefix requirement |
| P11 | HEAD and pristine tree at the time of measurement | `git log -1 --format='%h %s'`; `git status --porcelain \| wc -l` | `3a98c79 record(world) iter 135: …`; `0` | VERIFIED-BY-ME |
| P12 | No CODE surface writes `workflow_dispatch` anywhere in the repo; the only occurrences are design-doc prose, so the lever's observable is owned by `ci.yml` alone | `grep -rn 'workflow_dispatch' --include='*.go' --include='*.sh' --include='*.yml' --include='*.yaml' . --exclude-dir=.git` (code surfaces); `grep -rn 'workflow_dispatch' . --exclude-dir=.git \| wc -l` (whole repo) | code surfaces → **0**; whole repo → **15**, all in `design_docs/world-mission*.md` prose (P2's own record of the defect) — no competing mechanism writes the observable | VERIFIED-BY-ME |
| P13 | A second trigger line `workflow_dispatch:` at lead 2 cannot corrupt the sibling's job enumeration: the `on:` block sits at `:3–7` BEFORE `jobs:` at `:16`, and the sibling's `jobLine` regexp `^  ([a-z0-9-]+):$` does not match the underscore in `workflow_dispatch` | `grep -n '^jobs:' .github/workflows/ci.yml`; `printf '  workflow_dispatch:\n' \| grep -E '^  ([a-z0-9-]+):$'` | `jobs:` at `:16` (after the `on:` block); the trigger line does NOT match `[a-z0-9-]+` (underscore excluded), `rc=1` — the added line is inert to `TestGoToolchainPinsAgreeAndMatchJobList`'s job set | VERIFIED-BY-ME (code read + regex probe) |
| P14 | `on:` appears exactly once at top level, so the new on-block parser is unambiguous | `grep -c '^on:$' .github/workflows/ci.yml` (COLUMN-0 exact) with the same-call control `grep -c '^  pull_request:$'`; **the original `awk '{$1=$1}; $0=="on:"'` was REPLACED by the controller after measuring that `{$1=$1}` rebuilds the record and strips leading whitespace, so it matched an indented `  on:` too** (two arms: awk MATCHes both a column-0 `on:` and an indented `  on:`; `grep -c '^on:$'` returns 0 on the indented file and 1 on the top-level one — control fires) | **1** (one column-0 `on:`), control `^  pull_request:$` -> **1**. The repaired command is now STRICTER than, and identical in kind to, the parser it certifies — closing the verify/ship gap oc-glm-5-2 named | VERIFIED-BY-ME |
| P15 | The `host/verifygate` test package already imports `os`, `path/filepath`, `strings`, `slices`, `testing` — the new test adds zero imports | `sed -n '1,12p' host/verifygate/toolchain_pin_gate_test.go` | import block contains `"go/version"`, `"os"`, `"path/filepath"`, `"regexp"`, `"slices"`, `"strings"`, `"testing"` | VERIFIED-BY-ME |
| P16 | `gh workflow run` dispatches against a branch/tag NAME, not an arbitrary commit SHA: `-r, --ref string` is documented as the name of the ref containing the workflow file to run | `gh workflow run --help 2>&1 \| grep -n -- '\-\-ref'` | line 22 `-r, --ref string  Branch or tag name which contains the version of the workflow file you'd like to run`; line 36 `$ gh workflow run triage.yml --ref my-branch` — a NAME, never a SHA | VERIFIED-BY-ME |
| P17 | This repo's default branch is `dev` — the very branch the lever lands on — so `workflow_dispatch` (availability tied to the default branch) is live exactly where the lever ships | `gh api /repos/sunholo-data/ailang-world --jq .default_branch` | `dev`, rc=0 | VERIFIED-BY-ME |

## Conflict Surface

- **The sibling's `[ci.yml]`-exact list at `toolchain_pin_gate_test.go:206`.** The sibling
  `TestGoToolchainPinsAgreeAndMatchJobList` asserts
  `if !slices.Equal(workflowFiles, []string{"ci.yml"})` at `:206` (message
  `workflow files=%v, want exactly [ci.yml]; a second workflow may carry unscanned toolchain
  pins`), measured at exactly **1** occurrence. Therefore adding ANY second workflow file
  reds that sibling **by construction** — every second-workflow arm in this design (M4, M5,
  AC4) is a **multi-owner red**: it moves BOTH MY gate and the sibling's `[ci.yml]` list.
  Consequently M4 and M5 assert on MY test's OWN message text (the `deploy.yml` name / MY
  test's `--- PASS`), **never on a bare `rc`** — `rc` cannot distinguish which gate reds —
  and their expected result is an **ENUMERATED red set** whose sibling member is explained,
  not `rc=0`/`rc=1`.
- **The added `workflow_dispatch:` line is inert to the sibling's job set.** The sibling's
  job-enumeration regexp `^  ([a-z0-9-]+):$` excludes underscore, so the trigger line does
  NOT match (controller re-measured: the trigger line rc=1 while a real `go-verify:` job
  line rc=0 — the control fires, P13). The added line cannot corrupt the sibling's job list.
- **`host/verifygate` base-red without `AILANG_BIN`.** The package reds at base without the
  pinned binary (P6: 17 `--- FAIL`, all `AILANG_BIN is unset`). Every acceptance criterion
  and mutation arm touching this package carries the
  `AILANG_BIN=/tmp/ailang-v0300/ailang` prefix; a red there may be the environment, not the
  mutation, so a mutation's result is read only in a prefixed, LANDED-proven run.
- **If a second workflow EVER legitimately lands.** The sibling's `:206` list must change in
  the SAME edit that adds the workflow, or the sibling stays red forever (residual 4,
  Deferred Scope). That is a policy decision for a future row, not this one — today there is
  exactly one workflow and both gates agree on it.
- **Thesis reconciliation.** The headline "pins EVERY workflow file to it" is MY gate's
  DESIGNED coverage: the enumerator derives the workflow set at runtime and requires each
  file to declare the lever (option C rejects a hardcoded single-file read for exactly this
  reason). The sibling's `[ci.yml]` list is a SEPARATE, currently-narrower policy (one
  workflow) that a second workflow would force BOTH gates to revisit together. Read that
  way the two claims are not in tension: they describe different gates with different,
  currently-coincident scopes.

## Options considered

**A — Status quo / "advance `dev`" (no lever).** Rejected. On a dropped delivery there is
literally no action that verifies the offending commit: advancing `dev` replaces the commit
under test, and the git-API empty-commit workaround only fires on a PR. This is the exact
dead end that cost iterations 128 and 129; it is the defect, not an option.

**B — `workflow_dispatch` with `inputs:` (parameterized re-run).** Considered and rejected.
The lever's only job is to re-run verification on the tip of a named ref (`--ref dev` —
P16, P17). A dispatch run is created against the tip of whatever branch or tag name the
caller passes to `gh workflow run --ref`; the workflow already reads its inputs from the
checkout, not from dispatch parameters. `inputs:` would invite divergent parameterization
(which job, which leg, a marker nobody writes) that this repo does not need and that a static
gate would then have to police. The default — a bare `workflow_dispatch:` with no inputs —
is the minimal, valid, GitHub-supported form: it enables "Run workflow" with no parameters.
Adding inputs is scope creep, not capability.

**C — Pin the lever to `ci.yml` only (hardcode the single file, no enumeration).**
Considered and rejected on the "every workflow file" requirement. A gate that reads only a
hardcoded `ci.yml` passes silently the day a second workflow file appears without the lever
(the addition-shaped hole this item explicitly names). The enumeration must be derived at
runtime — the `filepath.Glob` idiom the sibling already uses (P5) — so the gate's coverage
grows with the workflow set, and the anti-vacuity floor keeps an empty enumeration from
printing a checkmark.

**D — The shipped design: add `workflow_dispatch:` (no inputs) + `TestEveryWorkflowDeclaresDispatchLever`.**
Chosen. One keyed line in `ci.yml`'s `on:` block (P3), plus a static consistency gate in
`host/verifygate` in the exact idiom of the siblings (P4): iterate the enumerated workflow
files, require each to declare the lever as a trigger in its `on:` block, fail loudly on an
empty enumeration, and carry one attributed message per defect. The lever's limitation is
declared in the code comment and in residual 1. Cost ~0.10 day.

## Decision

Add to `.github/workflows/ci.yml`'s `on:` block, at trigger depth (lead 2, matching `push:`
and `pull_request:` — P3), a single bare trigger:

```yaml
on:
  push:
    branches: [dev]
  pull_request:
  workflow_dispatch:
```

No `inputs:` (Option B's argument above). This is purely additive: GitHub triggers are
orthogonal — declaring a new event does not change how `push`/`pull_request` fire, so an
ordinary push or PR runs exactly as before (reasoned from GitHub's documented event model,
NOT measured live against a hosted repo — the repo has no external trigger this session; the
claim's strength is recorded accordingly and its observable is the unchanged `push:`
`pull_request:` bytes, P3).

Add a static gate to `host/verifygate`. The design sketch the implementing sprint owns the
final bytes of:

```go
package verifygate

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// onBlockTriggerKeys returns the trigger keys declared DIRECTLY under the on: block — the
// same depth as push:/pull_request: — so a workflow_dispatch folded into nested inputs: or
// buried in a comment does not count.
func onBlockTriggerKeys(t *testing.T, path, src string) []string {
	t.Helper()
	lines := strings.Split(src, "\n")
	start := -1
	for i, l := range lines {
		// EXACT byte equality at column 0 -- deliberately NOT strings.TrimSpace. A trimmed
		// comparison anchors on ANY nested `on:` (under a job step's `with:`/`inputs:`),
		// so the gate could scan a non-trigger block and pass vacuously. This matches
		// P14's verification method (see P14 -- its command was itself repaired to a
		// column-0 test, because `awk '{$1=$1}' ` strips leading whitespace and matched
		// an indented `  on:` too; measured, arms differ).
		if l == "on:" {
			start = i
			break
		}
	}
	if start == -1 {
		t.Fatalf("instrument failure: %s has no top-level `on:` trigger block", path)
	}
	// gpt5-6-sol R2: fail if the top-level `on:` is DUPLICATED, not merely absent.
	dupes := 0
	for _, l := range lines {
		if l == "on:" {
			dupes++
		}
	}
	if dupes != 1 {
		t.Fatalf("instrument failure: %s has %d top-level `on:` blocks, want exactly 1", path, dupes)
	}
	onLead := len(lines[start]) - len(strings.TrimLeft(lines[start], " "))
	triggerLead, keys := -1, []string{}
	for _, l := range lines[start+1:] {
		tok := strings.TrimSpace(l)
		if tok == "" || strings.HasPrefix(tok, "#") {
			continue
		}
		lead := len(l) - len(strings.TrimLeft(l, " "))
		if lead <= onLead {
			break // left the on: block
		}
		if triggerLead == -1 {
			triggerLead = lead
		}
		if lead == triggerLead {
			if kv := strings.SplitN(tok, ":", 2); len(kv) == 2 {
				// gpt5-6-sol R2: validate the VALUE. `workflow_dispatch` is legal as an
				// empty/null mapping key or as a mapping-valued trigger config -- never as a
				// bare scalar. `workflow_dispatch: garbage` must NOT count as declared.
				if kv[0] == "workflow_dispatch" && strings.TrimSpace(kv[1]) != "" {
					t.Errorf("%s: `workflow_dispatch:` has scalar value %q; want an empty key or a mapping",
						path, strings.TrimSpace(kv[1]))
					continue
				}
				keys = append(keys, kv[0])
			}
		}
	}
	return keys
}

// TestEveryWorkflowDeclaresDispatchLever pins queue-item 47's lever so it cannot be deleted
// silently. A push to dev whose webhook GitHub drops is PERMANENTLY unverifiable (never
// replayed, no run created), and the only API-driven lever is workflow_dispatch, which no
// workflow currently declares (P2). This gate asserts EVERY enumerated workflow file declares
// the lever as a trigger in its on: block.
//
// DECLARED RESIDUAL: this is a STATIC text scan over YAML. It proves the lever is DECLARED,
// never that a dispatch RUN is created or is green. And a workflow_dispatch run is NOT
// equivalent to the event it replaces: its checks do not satisfy PR branch protection
// (measured by this mission — required contexts can read success on the head SHA while
// gh pr checks --required still lists only the pull_request-event context with
// mergeStateStatus=BLOCKED). The lever buys A VERDICT ON A COMMIT, not A MERGEABLE PR. It
// also cannot see a workflow added OUTSIDE .github/workflows/ (a singular `workflow` dir, a
// root-level .yaml, a hidden file), a case-mismatched filename (the Glob is case-sensitive),
// or a nested subdirectory (which GitHub itself does not scan either).
func TestEveryWorkflowDeclaresDispatchLever(t *testing.T) {
	matches, err := filepath.Glob(filepath.Join(repoRoot, ".github", "workflows", "*"))
	if err != nil {
		t.Fatal(err)
	}
	// ANTI-VACUITY FLOOR: an empty enumeration FAILS LOUDLY rather than printing a checkmark.
	if len(matches) == 0 {
		t.Fatal("instrument failure: no workflow files enumerated under .github/workflows/ — an empty set proves nothing about the lever")
	}
	for _, m := range matches {
		raw, err := os.ReadFile(m)
		if err != nil {
			t.Fatal(err)
		}
		src := string(raw)
		// known-positive control: a scan that cannot see the on: block is reading the wrong file.
		if !strings.Contains(src, "on:") {
			t.Fatalf("instrument failure: %s does not contain known-positive control %q", m, "on:")
		}
		triggers := onBlockTriggerKeys(t, m, src)
		if !slices.Contains(triggers, "workflow_dispatch") {
			t.Errorf("%s: on-block triggers=%v lack workflow_dispatch — a dropped push to dev is permanently unverifiable; every workflow file must declare the lever", filepath.Base(m), triggers)
		}
	}
}
```

This matches the sibling idiom: `repoRoot`/`findRepoRoot` (`ail_binary_gate_test.go:27,31`),
known-positive control floors, line-based parsing over YAML, the anti-vacuity floor on a
derived set, and EXACTLY ONE attributed message per defect (the `t.Errorf` names the file
and the absent lever; it never cascades — P5's precedent). It reuses the sibling's import set
(P15), so no imports change.

## Milestones

Total **~0.10 day** — deliberately small; this is a one-line lever plus one gate function.

- **M1 (0.05d)** — add `  workflow_dispatch:` to `ci.yml`'s `on:` block (P3); the bare form,
  no `inputs:`. Verify: `grep -cE '^  workflow_dispatch:' .github/workflows/ci.yml` → 1 with
  same-call controls `pull_request` → 1 and `^  push:` → 1.
- **M2 (0.05d)** — add `TestEveryWorkflowDeclaresDispatchLever` (+ `onBlockTriggerKeys`) to
  `host/verifygate/dispatch_lever_gate_test.go`; `gofmt`; run green under the prefixed
  binary; rehearse mutations M1–M9 in probe worktrees (restore byte-identical, porcelain
  0). Verify: AC1–AC6.

No `.ail`, no `tools/launchd/*`, no other workflow file, no `go.mod`/dependency change.

## Acceptance criteria

Every criterion runs with `export PATH=/opt/homebrew/bin:$PATH`. Any command touching the
gate package carries the `AILANG_BIN=/tmp/ailang-v0300/ailang` prefix (P6). "Probe worktree"
= detached worktree at the landing commit, mutations proven LANDED before results are read,
restored byte-identical, porcelain 0 (house recipe).

- **AC1 — the lever is declared as a trigger in `ci.yml`'s `on:` block.**
  `grep -cE '^  workflow_dispatch:' .github/workflows/ci.yml` → **1**, with same-call
  controls `grep -c 'pull_request'` → **1** and `grep -cE '^  push:'` → **1**.
  **Base:** **0** (P2 — this is the finding the sprint reverses; red at base is the point,
  not a broken criterion).
- **AC2 — the new gate exists and RUNS green on the unmutated tree.** The file
  `host/verifygate/dispatch_lever_gate_test.go` exists and is gofmt-clean.
  `AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -run TestEveryWorkflowDeclaresDispatchLever -count=1 -v`
  → rc=0 with one `=== RUN` and one `--- PASS`; paired nonsense control
  `-run '^TestNoSuchGateZZZ$'` prints `[no tests to run]`. **Base:** N/A (the test does not
  exist at HEAD; the whole-package prefixed run is P6 rc=0/0 FAIL). AC2's teeth are M1–M9.
- **AC3 — removal-shaped: the gate REDS when the lever is deleted.** Probe worktree: remove
  the `workflow_dispatch:` trigger line from `ci.yml` (LANDED: `grep -cE '^  workflow_dispatch:'` → 0
  with same-call control `grep -c 'pull_request'` still 1); prefixed scoped run → rc≠0 whose
  output names `ci.yml` and `workflow_dispatch` exactly once (single attributed message).
  **Base:** N/A (test absent).
- **AC4 — addition-shaped: the gate REDS when a SECOND workflow file lacks the lever.**
  Probe worktree: add `.github/workflows/deploy.yml` with `on:` / `  push:` and no
  `workflow_dispatch` (LANDED: `ls .github/workflows/ \| wc -l` 1→2); prefixed scoped run
  → rc≠0 whose output names `deploy.yml` — assert on MY test's message text (the sibling
  `TestGoToolchainPinsAgreeAndMatchJobList` ALSO reds its `[ci.yml]` list here, so rc alone
  cannot distinguish the gates; see Conflict Surface). This proves the gate LOOKS beyond the
  file it was born from. **Base:** N/A.
- **AC5 — anti-vacuity: the gate REDS loudly on an empty enumeration.** Probe worktree:
  `mv .github/workflows /tmp/wf_backup_$$` so the Glob returns empty (LANDED:
  `ls .github/workflows/ \| wc -l` 1→0); prefixed scoped run → rc≠0 whose output contains the
  `instrument failure: no workflow files enumerated` floor message. Assert on THAT message
  (single-owned) because the sibling tests also red on the missing `ci.yml` read.
  **Base:** N/A.
- **AC6 — hygiene, base-green preservation, and the typecheck discipline.**
  `gofmt -l host/verifygate/` → empty; `go vet ./host/verifygate/` → rc=0 (also the
  typecheck for any test-file mutation, per P8);
  `AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -count=1` → rc=0 on the
  unmutated tree (the well-formed `push:`/`pull_request:`/`workflow_dispatch:` triggers pass
  the new parser AND the sibling pin-tests stay green — with one workflow there is no
  sibling collision); `git status --porcelain` → 0 lines after all rehearsals. **Base:** P6
  (prefixed rc=0/0 FAIL) + P7 (gofmt/vet clean) — this criterion stays green; AC3–AC5 are
  its teeth.

## Named RED mutations

Venue: `./scripts/verify_go.sh` runs `go test ./... -count=1` inside the gated `go-verify`
job (ci.yml:165) and refuses to run without `AILANG_BIN` (P10), so every arm below is a
CI-gate red; the local rehearsal command is the prefixed scoped test. House recipe per arm:
prove LANDED by a count that MOVES before reading the result; restore byte-identical;
porcelain 0. Cross-checked against the Decision: every arm targets a construct the design
SHIPS (the lever line, the enumerator, the floor, the trigger-parser). The single test-file arm
(M6) uses `go vet ./host/verifygate/` as their typecheck, because `go build ./...` does not
compile `_test.go` (P8).

| # | Mutation | File | What it neuters | Predicted result | Landed-proof before reading the result |
|---|---|---|---|---|---|
| M1 | Delete the `  workflow_dispatch:` trigger line from `ci.yml`'s `on:` block | `ci.yml` | the lever itself — the removal shape | prefixed scoped run RED: output names `ci.yml` and `workflow_dispatch` exactly once | `grep -cE '^  workflow_dispatch:'` 1→0 with same-call control `grep -c 'pull_request'` still 1 |
| M2 | Rename the trigger `workflow_dispatch` → `workflow_dispatcht` (trailing `t`) | `ci.yml` | the exact trigger-key name — a plain substring scan could pass on the typo | prefixed scoped run RED naming the typo'd file/absent exact key | `grep -cE '^  workflow_dispatch:'` 1→0 AND `grep -cE '^  workflow_dispatcht:'` 0→1 (both in one call) |
| M3 | Replace the trigger line with `  # workflow_dispatch: recovery lever` (comment) | `ci.yml` | the on-block depth/comment handling — proves the parser requires a REAL trigger, not a substring (a `strings.Count(src,"workflow_dispatch")>=1` gate would pass this) | prefixed scoped run RED: trigger-set membership fails | keyed-trigger `grep -cE '^  workflow_dispatch:'` 1→0 while raw `grep -c 'workflow_dispatch'` stays **1** (the comment form) — proves it is a comment, not a trigger |
| M4 | ADD `.github/workflows/deploy.yml`: `on:` + `  push:` only, no `workflow_dispatch` | `.github/workflows/deploy.yml` (new) | the "EVERY workflow file" requirement — the addition shape that proves the gate LOOKS | ENUMERATED RED SET: (1) MY test `TestEveryWorkflowDeclaresDispatchLever` reds naming `deploy.yml` — the single-owned assertion, keyed on that message text; (2) the sibling `TestGoToolchainPinsAgreeAndMatchJobList` reds its `[ci.yml]` list (`:206`) — multi-owned by construction (Conflict Surface). Assert on MY test's `deploy.yml` message, never on bare `rc` | `ls .github/workflows/ \| wc -l` 1→2; `grep -c 'workflow_dispatch' deploy.yml` → 0 |
| M5 | (green control, must-stay-green) ADD `.github/workflows/deploy.yml` declaring `workflow_dispatch:` in its `on:` block — every file has the lever | `.github/workflows/deploy.yml` (new) | over-firing: the gate must not refuse a second lever-bearing workflow | MY gate must STAY GREEN: prefixed scoped `-run TestEveryWorkflowDeclaresDispatchLever` → `--- PASS` (assert MY test's PASS, NOT package `rc`); the package as a whole still reds ONLY on the sibling's `[ci.yml]` list (`:206`) — that single red member is ENUMERATED and explained (Conflict Surface), not a failure of MY gate | `ls .github/workflows/ \| wc -l` 1→2; `grep -cE '^  workflow_dispatch:' deploy.yml` 0→1 |
| M6 | Neuter the anti-vacuity floor (test-file mutation): delete the `if len(matches)==0 { t.Fatal(…) }` guard | `host/verifygate/dispatch_lever_gate_test.go` | the anti-vacuity floor — combined with the empty-enumeration scenario it is what reds, not the (multi-owned) sibling reads | typecheck: `go vet ./host/verifygate/` rc=0 (removal compiles); behavior: with the floor removed AND the empty-dir scenario (M6 + M6b) the scoped run prints a checkmark / rc=0 instead of the floor message — proving the floor is load-bearing. Restore both | `grep -c 'instrument failure: no workflow files enumerated' host/verifygate/dispatch_lever_gate_test.go` 1→0 (floor gone); and `ls .github/workflows/ \| wc -l` →0 for the empty-dir arm (M6b), restored after |

| M7 | Remove the top-level `on:` and add a NESTED decoy: a job step gains `        on:` with `          workflow_dispatch:` beneath it | `ci.yml` | the parser's column-0 ANCHORING — a trimmed match would anchor on the decoy and pass | prefixed scoped run RED with the `no top-level \`on:\` trigger block` instrument-failure message; it must NOT report the lever as declared | `grep -c '^on:$'` 1->0 AND `grep -c 'workflow_dispatch'` stays >=1 (the decoy is still in the file) — the pair is what proves the decoy was planted rather than the lever deleted |
| M8 | Duplicate the top-level `on:` block (a second column-0 `on:` later in the file) | `ci.yml` | the duplicate-detection floor gpt5-6-sol required | prefixed scoped run RED naming `has 2 top-level \`on:\` blocks, want exactly 1` | `grep -c '^on:$'` 1->2 |
| M9 | Replace the bare trigger with the scalar form `  workflow_dispatch: garbage` | `ci.yml` | the trigger-VALUE validation — a key-only scan would accept a malformed value | prefixed scoped run RED naming `has scalar value "garbage"` | `grep -cE '^  workflow_dispatch:$'` 1->0 AND `grep -c '^  workflow_dispatch: garbage$'` 0->1 (both in one call) |

Green control for all arms: the unmutated post-sprint tree passes AC1–AC6; M5 is the named
must-stay-green arm (addition WITH lever passes — only an addition lacking the lever reds,
M4). Each assertion above is keyed on a SINGLE-OWNED observable (the specific trigger count,
the specific file name in the message, the specific floor message) so a mutation cannot pass
for an adjacent reason; where the observable is genuinely multi-owned (M4/M5 also move the
sibling's `[ci.yml]` list, M6b also breaks sibling `os.ReadFile`), the arm asserts on MY
test's message text instead.

## Deferred Scope

- **Live end-to-end proof that a dispatch run is created and green against the tip of the
  named ref (`dev`).** This doc adds the static lever and the static gate. Creating an
  actual `workflow_dispatch` run against a live GitHub repo is a runtime action gated to
  the sprint's CI, not a unit test; a live dispatch rehearsal belongs to the implementation
  sprint's Verification Log, not to a design-doc acceptance criterion.
- **The static gate proves the lever is DECLARED, not that a run is created.** A green gate
  means every enumerated workflow file TEXT declares `workflow_dispatch:` in its `on:`
  block. It does NOT prove any dispatch RUN is ever created, nor its result. Nobody may
  read the gate's green as the stronger claim; only a live dispatch against `dev`'s tip
  proves a run (residuals 2, 7).
- **Relaxing the sibling's `[ci.yml]`-only workflow list (P5, residual 4).** If this repo
  ever legitimately ships a second workflow (even one WITH the lever), the sibling's
  `TestGoToolchainPinsAgreeAndMatchJobList` `:206` assertion must change in the SAME edit.
  That is a policy decision for a future row, not this one — today there is exactly one
  workflow and both gates agree on it.
- **Structural YAML parsing / `actionlint` adoption.** The line-scan residual class is
  already declared by the sibling tests and inherited here; no linter runs in this repo
  (P9). Not a row-47 repair.

## Declared residuals

1. **The lever buys a VERDICT ON A COMMIT, not a MERGEABLE PR.** A `workflow_dispatch` run's
   checks do NOT satisfy branch protection on a PR: measured by this mission previously, all
   four required contexts can read `success` on the head SHA while `gh pr checks --required`
   still lists only the context from a real `pull_request` event, with
   `mergeStateStatus=BLOCKED`. The gate's claim is therefore scoped: it pins that the lever
   is DECLARED so a dropped webhook to `dev` can be re-verified; it does NOT claim "CI can
   always be re-triggered", "a dropped event is always recoverable", or "this unblocks
   merges". Stated in the code comment and here.
2. **Static text scan over YAML.** Like the siblings, the gate sees the lever TEXT, never
   that a dispatch RUN is created, never its result, and never a step-level `if:` that
   disables a job at runtime. A workflow_dispatch is valid syntax whether or not its jobs
   run; only a live dispatch rehearses the run.
3. **The enumerator sees only `.github/workflows/*` at one level.** A workflow added as
   `.github/workflow` (singular — a documented footgun), a root-level `.yaml`, a hidden
   file (Glob `*` skips dotfiles), a case-mismatched filename, or a nested subdirectory is
   invisible — though GitHub itself also does not scan nested subdirectories, so the
   top-level Glob is aligned with the platform's actual behavior. The anti-vacuity floor
   (AC5/M6) makes an EMPTY enumeration loud; it cannot make an UNSEEN enumeration loud.
4. **The sibling's `[ci.yml]`-only list conflicts with ANY second workflow file (P5).**
   Adding a second workflow — even one WITH the lever (M5) — reds
   `TestGoToolchainPinsAgreeAndMatchJobList:206`. The two gates currently agree (one
   workflow); the day a second legitimately appears, the sibling's list must be updated in
   the same edit (Deferred Scope). This doc does not change that gate.
5. **`inputs:` deliberately omitted.** A bare `workflow_dispatch:` gives the caller the
   "Run workflow" button with no parameters; there is nothing to parameterize for a
   re-verify-the-tip-of-a-named-ref use case. If a future row wants selective legs,
   `inputs:` plus a parallel gate update is its own change.
6. **Adding the trigger does not change ordinary push/PR runs — stated, not live-measured.**
   GitHub event triggers are orthogonal and purely additive; the claim rests on the
   platform's documented event model and on the unchanged `push:`/`pull_request:` bytes
   (P3), not on a hosted-repo observation (this session had none). An over-strong
   "measured live" phrasing is deliberately avoided.
7. **The lever re-verifies the TIP OF A NAMED REF, never an arbitrary commit SHA.**
   `gh workflow run` dispatches against a branch/tag NAME (`-r, --ref string` — P16), and
   `workflow_dispatch` availability is tied to the workflow existing on the DEFAULT branch,
   which here is `dev` (P17) — the branch the lever lands on. So the row-47 scenario the
   lever recovers is precisely a dropped push that left `dev`'s HEAD with no run:
   dispatching `--ref dev` runs against `dev`'s TIP, which IS that unverified commit only
   while `dev` has not advanced since. Once `dev` advances, the dropped commit is no longer
   any ref's tip and remains permanently unverifiable. This is distinct from residual 1
   (branch protection): residual 1 bounds the PR case; this residual bounds the recoverable
   window even on `dev` itself. The code comment and Deferred Scope state it plainly.

## In-PR hardening after the evaluator (iteration 136)

The `sonnet` evaluator scored **82/100 PASS, zero blocking** (generator≠judge: the executor was
OpenAI's codex) in its own worktree, and — the part that makes the verdict a measurement — it
**named the attacks that FAILED**: a hidden dotfile, a nested subdirectory, case/extension variants,
all nine canonical mutations, and the anti-vacuity floor's own precondition test. None produced a
false green. **The gate's core promise — never silently certify a workflow that lacks the lever —
held under every attack aimed at it.** What it found instead was the opposite failure mode: the
line-scan parser fails LOUD on several forms of *valid, lever-declaring* YAML.

**Two of its five findings were reproduced first-party and FIXED IN THIS PR** (`gofmt`/`vet` clean,
full drill re-run below):

1. **An inline `#` comment on the trigger line was misread as a scalar value.**
   `  workflow_dispatch: # manual re-run lever` is valid YAML that declares the lever, and it is the
   single most natural edit a maintainer makes to explain the line — yet it tripped the
   scalar-value branch and redded CI with **two** cascading messages. **The discriminating control
   is what makes this a mechanism rather than a guess:** the same comment on its OWN line was
   rc=0/PASS=1, so the defect was specific to the inline form and the parser was not merely
   comment-hostile. Fixed by stripping an inline comment before judging the value.
2. **The doc asserted something false about its own shipped code.** Residual 7 claimed *"the code
   comment and Deferred Scope state it plainly"* about the named-ref window; grepping the shipped
   comment returned **0** hits for named-ref/tip/SHA, with the control firing at **2** for the
   branch-protection sentence that IS there. The comment now states the named-ref limitation, so
   a reader auditing only the code meets it. (This is the same shape as the residue the carve-out
   left in the mutation ranges: a claim about a sibling artifact that nobody re-derived.)

**The hardening was proven NON-WEAKENING before it was believed** — the risk of stripping comments
is over-approximation, so the drill includes the arm that would catch it. Seven arms, one tree,
each restored byte-identical (`sha256` equal to base):

| Arm | Expectation | Observed |
|---|---|---|
| green control (pristine) | PASS | rc=0 RUN=1 PASS=1 |
| inline comment (the finding) | now PASS | rc=0 RUN=1 PASS=1 |
| M9 `workflow_dispatch: garbage` | still RED | rc=1 FAIL=1, `has scalar value "garbage"` |
| **`garbage # note`** (over-approximation control) | **still RED** | rc=1 FAIL=1, `has scalar value "garbage"` — the value BEFORE the `#` is still judged |
| M1 delete the lever | still RED | rc=1 FAIL=1, `lack workflow_dispatch` |
| M8 duplicate top-level `on:` | still RED | rc=1 FAIL=1, `want exactly 1` |
| M7 nested decoy, no top-level `on:` | still RED | rc=1 FAIL=1, `instrument failure` |
| final green control | PASS | rc=0 RUN=1 PASS=1 |

**Three findings are NOT fixed here and become a queue row rather than a silent scope widening**
(they need a decision this item does not own — whether the gate adopts a structural YAML parser):
the parser false-reds on the **quoted `"on":`** form (which is the standard remedy for YAML 1.1's
`on` → boolean footgun, so this would break CI the day anyone applies it), on **flow-style**
`on: {push: …, workflow_dispatch: }`, and on a **tab-indented** first trigger line (`TrimLeft(l, " ")`
strips spaces only, so the block reads as already exited and the trigger set silently empties). All
three fail LOUD, which is the accepted direction, and none is a silent pass. Also unfixed and
declared: the scalar-value arm emits **two** messages rather than one, so the Decision's
"never cascades" phrasing is an over-claim the planner's D5 already flagged.

**A correction to D0, so the record is not stale:** the doc's sketch WAS shipped byte-verbatim by
the executor and was measured green that way; the two hardening edits above were applied afterwards
by the controller in response to the evaluator, so the shipped bytes are no longer identical to the
sketch in `## Decision`.

## Quorum verification log

Designer: `pi:ollama/deepseek-v4-flash:0731-cloud` (designer rotation; FIRST authoring run on this
lane since the 2026-08-28 rotation amendment). Authoring verdict `ok` (304 s, 29 tool executions);
revision verdict `ok` (203 s, 20 tool executions). Metered cost of the lane: **$0** (flat-rate).

### Round 1 — BLOCKED at full strength (`absent_reviewers` EMPTY, 3/3 external reviewers present)

Artifact `.ailang/state/mission-quorum/w-ci-recovery-lever-absent-2026-08-28T09-29-30Z.json`.
Metered $0.0901. 3/3 reject + controller reject. Objections landed on **three different surfaces**,
so this is doc immaturity, not a decomposition signal.

| Reviewer | Surface | Objection | Disposition |
|---|---|---|---|
| `gpt5-6-sol` | capability claim | the doc claimed the lever re-runs verification on an "arbitrary head SHA" | **MEASURED FIRST-PARTY BY THE CONTROLLER, CONFIRMED, and sharpened**: `gh workflow run --help` documents `-r, --ref` as *"Branch or tag name"* — not a SHA. Separately measured: this repo's `default_branch` is **`dev`**, which is what makes the lever *available* and what makes it *sufficient* for the row's actual scenario. Fix applied: claim rescoped to THE TIP OF A NAMED REF; new residual added. |
| `gemini-3-1-pro` | rehearsability | mutation row M6 carried the literal placeholder `<new_test>`, a shell redirect from a nonexistent file | **MEASURED**: 1 occurrence (control: a present literal at 7); the doc never named the new test file anywhere. Fix applied: `host/verifygate/dispatch_lever_gate_test.go` named at every site; placeholder-shaped tokens re-derived as a SET and swept to 0. |
| `oc-glm-5-2` | sibling-gate collision | the sibling `TestGoToolchainPinsAgreeAndMatchJobList` asserts the workflow set is EXACTLY `[ci.yml]`, so every second-workflow arm is a multi-owner red | **CONFIRMED first-party** at `toolchain_pin_gate_test.go:206` (assertion occurs exactly 1x). Reached independently by the controller, whose own blocking objection was that the required `## Conflict Surface` section did not exist (`grep -ciE '^#{1,6} .*conflict'` -> **0**) while premise P5 carried a *dangling forward reference to it*. Fix applied: the section now exists, consolidating the collision, the multi-owner red sets, the inert-to-job-regexp finding, and the `AILANG_BIN` base-red. |

### Round 2 — BLOCKED at full strength (`absent_reviewers` EMPTY, 3/3 external present)

Artifact `.ailang/state/mission-quorum/w-ci-recovery-lever-absent-2026-08-28T09-37-31Z.json`.
Metered $0.0999. Objections **LOCALISED onto ONE surface** — the `onBlockTriggerKeys` parser's
anchoring — with `gemini-3-1-pro`'s a separate mechanical completeness gap. No reviewer flipped to
pass, so this is not the SPLIT signal; it is one helper needing three concrete repairs.

Closed under the **ratified narrow-refinement carve-out** (ratified iteration 44, so no first-use
gate): all three objections carry concrete reviewer-authored `proposed_fix` text and **none disputes
the design DIRECTION** (bare `workflow_dispatch:` + a static `host/verifygate` gate is accepted by
all three). Fixes applied in the reviewers' own words:

- **`oc-glm-5-2`, verbatim** — replace `strings.TrimSpace(l) == "on:"` with exact byte equality
  `l == "on:"`, and comment that this matches P14's verification method.
- **`gpt5-6-sol`, verbatim** — require an exact zero-indentation `on:`, **fail if it is absent or
  duplicated**, and **validate `workflow_dispatch` as an empty/null mapping key or a mapping-valued
  trigger rather than accepting any scalar**; plus RED mutations for (1) a nested `on:`/
  `workflow_dispatch:` decoy with no top-level `on:`, (2) duplicate top-level `on:` blocks, and
  (3) `workflow_dispatch: garbage` — landed as **M7**, **M8**, **M9**.
- **`gemini-3-1-pro`, verbatim** — prepend the `package verifygate` declaration and the `os` /
  `path/filepath` / `slices` / `strings` / `testing` import block to the sketch.

**AND A REVIEWER'S CITED CONTROL WAS REFUTED BY MEASUREMENT WHILE ITS FIX WAS APPLIED ANYWAY — the
distinction is this row.** `oc-glm-5-2` justified the exact-equality fix by saying it *"matches
P14's verification method (awk `$0=="on:"` exact match at column 0)"*. Measured before acting, two
arms: `awk '{$1=$1}; $0=="on:"'` prints MATCH for a column-0 `on:` **and** for an indented `  on:`,
because `{$1=$1}` rebuilds the record and strips leading whitespace; a genuine column-0 test
(`grep -c '^on:$'`) returns **0** on the indented file and **1** on the top-level one (control
fires). So P14 never established top-level-ness either. **That makes the objection STRONGER than
filed, not weaker** — neither the parser nor the premise row that certified it was column-0 aware —
and the fix was therefore applied AND P14's own command was repaired to the column-0 form, so the
verification method is now identical in kind to the code it certifies. Applying a reviewer's fix
while correcting the rationale it cited is not overriding the objection; it is satisfying it.

Every carve-out edit was proven LANDED by a count that MOVED before the result was read:
`TrimSpace(l) == "on:"` 1->0; `if l == "on:"` 0->2; `package verifygate` 0->1;
`want exactly 1` (duplicate floor) 0->2; scalar-value validation 0->1; mutation rows `M7|M8|M9`
0->3; the old trimming awk in P14 1->0; control `onBlockTriggerKeys` present throughout at 4.

## Related Documents

- Sibling gate to imitate: `design_docs/implemented/w-toolchain-pin-normalizer-accepts-malformed-gotoolchain.md`
  (section order, Premises/Acceptance/Mutations/residuals discipline; the example this doc
  mirrors).
- Sibling tests: `host/verifygate/ail_binary_gate_test.go:668`
  (`TestZ3PinDeclaredOnceAndInstalledInBothJobs`), `host/verifygate/toolchain_pin_gate_test.go:106`
  (`TestGoToolchainPinsAgreeAndMatchJobList`, whose `:197–207` enumerator this doc copies).
- Charter: `design_docs/world-mission.md` (queue row 47 record; iterations 128/129 account
  of the `a0b3162` dropped-delivery incident and the branch-protection limitation cited in
  residual 1).
- Coding standards: `design_docs/coding-standards.md` (non-vacuous gates, effects at the
  boundary, one attributed message per defect).
