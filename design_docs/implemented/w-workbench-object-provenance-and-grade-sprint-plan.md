# w-workbench-object-provenance-and-grade — sprint plan (rows 98+102, 2 milestones)

- Sprint: `w-workbench-object-provenance-and-grade` · queue rows **98** and **102** (clause 5) ·
  design doc [`w-workbench-object-provenance-and-grade.md`](w-workbench-object-provenance-and-grade.md)
  (quorum-closed: r1 BLOCKED on astra's upheld objection → designer revision `581c8cf`; r2 gemini and
  astra rejects closed by the controller's narrow-refinement carve-out `61c6e33`, doc §11).
- Base: `61c6e33` plus this plan's commit, on branch `sprint/w-workbench-object-provenance-and-grade`,
  worktree `/Users/voightkampff/dev/sunholo-data/.wt-world-iter186`. `origin/dev` = `6127ff3` (the
  doc's measurement base). The three commits between `6127ff3` and `61c6e33` touch only the design
  doc (`git diff --stat origin/dev 61c6e33` → 1 file, `design_docs/planned/…`), so every `host/`
  anchor the doc took at `6127ff3` is an anchor at base.
- Planned 2026-09-25 by mission-control iteration 186 (sprint-planner). **The whole design was
  prototyped by the planner in this worktree** (production + tests, §3/§6/§7 of the doc), gated,
  drilled in four states, measured, and then removed; `git status --porcelain` was empty before this
  plan was written. Every code block below is a copy of the banked prototype diff, which was
  `gofmt`-clean, vetted and green. It is not a paraphrase of the doc.
- Files changed by the sprint (doc §5, confirmed by `git diff --numstat` on the prototype):
  `host/workbench/render.go` (+2/−2), `host/workbench/render_test.go` (+61/−0),
  `host/daemon/workbench.go` (+46/−7), `host/daemon/workbench_test.go` (+171/−0). **Nothing else.**
  No `.ail` file, no view-model type change, no grammar change.
- Planner evidence: `~/.ailang/state/mission-world-iter186-evidence/planner/` —
  `prototype.diff` (sha256 `a5c83702ae651dc2c371e715d90c95335d91a54495e777cf4d5d67372daa46bc`),
  its two halves `m1.diff` (sha256 `3e0644425c76c0aef2f7adddb5df66a2418d4c23dbbca48c4c0e866411394eb3`)
  and `m2.diff` (sha256 `694db0457b38f9c51c8e2edcd888f50c029fe18d4c7001b7fd9a56054d5da101`), the
  drill runner `drill.py` (sha256 `7c0d37189c53ec6a9606dd51117c61b85835a3bdcff020d2b7ace0a8f83b820e`)
  and its wrapper `run_drills.sh`, the four drill transcripts `drill_full.txt`, `drill_m1.txt`,
  `drill_m1tests.txt`, `drill_basesuite.txt`, the gate transcripts `full_test.txt`,
  `m1_full_test.txt`, `verify_ail.txt`, `ac4.txt`, `m1_ac.txt`, `ac_checks.txt`,
  `v16_render_tests_at_base.txt`, and the pristine/prototype file copies under `base/` and `proto/`.

---

## §0 Planner findings (measured; read first)

The prototype **confirms the design**: all 22 §7 rows compile, all 22 turn the suite red, all 22
have their named killer in the red set, 14 are sole-killed, 17 survive the `6127ff3` tests and 5 are
already killed by a pre-existing test — every one of the doc's §7 per-row outcomes (compiles, rc,
named killer, sole, base-suite) was reproduced exactly, including every "Leaves" list. AC4–AC9 hold.
What follows are defects in how §7 *states* the mutants, plus stale numbers. None changes what the
executor builds. The design doc is **not** edited by this plan.

| # | Doc says | Measured | Resolution in this plan |
|---|---|---|---|
| P1 | §7 H5, H6, H8, H13: old text `if err != nil {`. | That text occurs **16 times** in the prototype `workbench.go`; a single-occurrence replace cannot land it. | §5 gives each row a context-bearing old text (the preceding assignment line + `\n\tif err != nil {`) that occurs exactly once. Measured: each lands once, compiles, and reproduces the doc's red set. |
| P2 | §7 H7: old text `if !ok {`. | Occurs **5 times** in the prototype `workbench.go`. | §5 anchors H7 as `if !ok {\n\t\treturn workbench.EdgeView{Relation: relation, Target: target,` (1 occurrence). The doc's first form `if false {` was re-run as control row **H7x**: `compiles=n`, `declared and not used: ok` — confirmed, and never counted. |
| P3 | §7 H17: "relation `"interface"` → `"transitionRef"`"; H14/H15/H16: `= "…"` (elided). | `"interface"` alone is not a stable anchor; the elided constants cannot be string-matched. | §5 gives H17 as `"interface", object.InterfaceHash` → `"transitionRef", object.InterfaceHash`, and H14–H16 with the full constant text. |
| P4 | §7 H10/H11 "delete `{Relation: …},`". | Deleting only the literal leaves a whitespace-only line (still compiles). | §5 deletes the whole line including its two leading tabs and newline, so the mutated file stays `gofmt`-shaped. Same outcome. |
| P5 | Header line 10, V15, §8: "62 production lines, 227 test lines"; numstat `48 10 workbench.go`, `161 0 workbench_test.go`, `66 0 render_test.go`; §8 M1 "70", M2 "219". | Planner prototype: `4 files changed, 280 insertions(+), 9 deletions(-)` = **289** (the same total). Numstat `46 7 host/daemon/workbench.go`, `171 0 host/daemon/workbench_test.go`, `2 2 host/workbench/render.go`, `61 0 host/workbench/render_test.go`. Production **57** changed lines (53 + 4), tests **232**. M1 = **65** (`2 files changed, 63 insertions(+), 2 deletions(-)`), M2 = **224**. | The split differs because the designer's deleted prototype is not recoverable and its line layout differed; this plan's numbers come from the banked diff. AC8 (≤ 330) holds with a margin of 41. |
| P6 | AC4 names 16 subtests but not `TestWorkbenchObjectGrade`'s six. | — | Fixed here: `test-object`, `test-object-payload`, `proof-labelled`, `proof-labelled-payload`, `registry`, `registry-payload` (suffix `-payload` = the `&payload=1` request). |
| P7 | AC7 fences with `go vet ./host/workbench ./host/daemon` and runs `go test ./host/workbench ./host/daemon`. | The planner drill fenced with `go vet ./...` and ran `./host/workbench ./host/daemon ./host/boundary`; results identical to the doc on every row. | The plan uses the wider forms (§2.3). Stricter than AC7, never looser. |
| P8 | §8 "M1 … Exit: R1–R5 drilled"; "M2 … H1–H17 drilled". | Split **confirmed** by two extra drills. At the M1-only tree (render files only) R1–R5 are all KILLED with their named killer; every H row is `NOT-LANDED count=0` (its code does not exist yet). With M2's *production* in place but M2's *tests* absent (`drill_m1tests.txt`), H1–H6, H9–H12, H14–H17 **survive** (rc=0) — M2's tests own them — and H7, H8, H13 are killed by pre-existing `TestWorkbenchSelectedEntry`/`TestWorkbenchTimelinePaging` subtests only (without the M2 named killer for H7/H8). | Ownership in §5 follows this: R1–R5 → M1; all H rows → M2. |
| P9 | V16: the new render tests are red at base. | Confirmed: prototype `render_test.go` against the `6127ff3` `render.go` → FAIL `/no-object`, `/object-without-edges`, `/zero-grade`; PASS `/supplied-edges`, `/supplied-reason`. | No change. |

Anchors re-verified at `61c6e33` (`cat -n`): `workbench.go:18-32` (message constants, closing `)` at
`:32`), `:34-40` (`acceptedWorkbenchKeys`), `:63-77` (`supportedWorkbenchQuery`), `:118-144`
(`entryEdges`, loop body `:132-141` exactly as V21 logs it), `:264-268` (object literal, 7 fields),
`:287` (`selected.Edges = edges`); `render.go:154` (grade line), `:160` (walk line),
`:166` (`"edge"` define), `:187` (`NewGradeUnavailable`); `workbench_test.go:396-407`
(`workbenchRegion`), `:414-422` (`objectFailingStore`), last line `:729`; `render_test.go:241-248`
(`renderPage`), last line `:359`; `handlers_test.go:71-77` (`testObject`), `:79` (`testCommit`),
`:141` (`seedGenesisEmbedded`); `read_deadline_test.go:681` (`errSentinelInternal`). Both test files
already import everything the new code uses (`bytes`, `context`, `net/http`, `strings`, `hashref`,
`store` in the daemon file; `strings` in the render file); **no import change**.

---

## §1 Baseline and prototype measurements

### 1.1 Commands and results (2026-09-24/25, load average 1.33 at start)

All commands ran with `export AILANG_BIN=$HOME/.pinned-ailang/ailang`
(`$AILANG_BIN --version` → `AILANG v0.41.0`, `Commit: 24ee108`; `go version` → `go1.26.6 darwin/arm64`).

| Gate | Command | Full prototype (M1+M2) | M1 half only |
|---|---|---|---|
| vet | `go vet ./... ; echo rc=$?` | **rc=0** | **rc=0** |
| test | `go test ./... -count=1` | **rc=0**; **20 `ok`**, 0 `FAIL` | **rc=0**; 20 `ok`, 0 `FAIL` |
| gofmt (touched) | `gofmt -l host/daemon/workbench.go host/daemon/workbench_test.go host/workbench/render.go host/workbench/render_test.go` | empty | empty |
| .ail | `bash ./scripts/verify_ail.sh ; echo rc=$?` | **rc=0**; `✓ world package gate PASSED: 9/9 steps performed non-zero work`; `✓ verify gate PASSED: 11 required identities verified, 40 named tests pass`; `git diff --name-only origin/dev -- '*.ail' \| wc -l` → `0` | — |
| AC4 | the §6 AC4 `-run` command | rc=0; **20** `--- PASS` (4 top-level + 16 subtests), 0 `--- FAIL` | render subset: rc=0, **7** `--- PASS` (2 + 5) |
| race (touched pkgs) | `go test -race ./host/workbench ./host/daemon -count=1` | both `ok` (workbench 1.3 s, daemon 12.8 s) | — |
| AC5 | the three `diff <(git show origin/dev:…) <(…)` commands | rc=0, rc=0, rc=0 | — |
| AC6 | the six `grep -c` controls | `1`, `1`, `0`, `1`, `0`, control `1` | — |
| AC9 / read budget | `git diff origin/dev -- host/daemon/workbench.go \| grep -c StateRoot`; `grep -c GetObject` / `context.Background` on `workbench.go` | `0`; `2` (base `2`); `0` | — |
| size | `git diff --shortstat origin/dev -- host/` | `4 files changed, 280 insertions(+), 9 deletions(-)` (289) | `2 files changed, 63 insertions(+), 2 deletions(-)` (65) |

`gofmt -l host cmd` over the whole tree lists pre-existing dirty files (for example
`host/daemon/session_middleware_test.go`); they are red at base, out of scope — **do not touch**.

### 1.2 The four drills (22 rows + control H7x + pristine P0 each; `drill.py`, §2.3)

| Drill | Tree | P0 | KILLED w/ named killer | SURVIVED | `compiles=n` | NOT-LANDED |
|---|---|---|---|---|---|---|
| full | M1 + M2 (prototype) | rc=0 | **22/22** (14 sole) | 0 | 1 (H7x, the control) | 0 |
| m1 | M1 only (render files) | rc=0 | R1–R5: **5/5** (R1, R3, R4 sole) | 0 | 0 | 18 (all H rows + H7x: code not present — expected) |
| m1tests | M2 production + M1 tests, `origin/dev` daemon tests | rc=0 | R1–R5 5/5; H13 (sole) | **14** (H1–H6, H9–H12, H14–H17) | 1 (H7x) | 0 |
| basesuite | M1+M2 production, both test files at `origin/dev` | rc=0 | H13 only (named killer is pre-existing) | **17** (R1, R3, R4, H1–H6, H9–H12, H14–H17) | 1 (H7x) | 0 |

In the m1tests and basesuite drills, R2, R5, H7 and H8 are killed (rc=1) but by pre-existing tests
only (`TestRenderGradeWithoutVerdictClaim/unavailable`, `TestRenderUnavailableProvenanceEdge`,
`TestWorkbenchSelectedEntry/unstored-edge-unavailable` + `TestWorkbenchTimelinePaging/emitted-links-resolve`,
`TestWorkbenchSelectedEntry/object-store-error`), exactly the doc's §7 base-suite column. Wall time:
about 80 s per 23-row drill on this rig (all four in 4 min 7 s).

### 1.3 CI's command list (the gate list)

`.github/workflows/ci.yml` has two jobs; CI green on the PR head **and** on the merge commit is the
binding gate. Job `ailang-verify` runs `./scripts/verify_ail.sh`. Job `go-verify` runs
`./scripts/verify_go.sh` (build, the evidence manifest, `go test ./... -count=1`, and the race leg),
the race-gate blindspot script, the bench smoke/claims/recorder-refusal arms, the gate-1 / queue-census
/ gate-0 script tests, and `scripts/check_no_personal_email.sh` + its test. Locally the executor
runs `verify_ail.sh`, `go vet ./...`, `go test ./... -count=1`, the race leg on the two touched
packages, and `check_no_personal_email.sh`. Per the standing rig rule, **do not run `verify_go.sh`
locally**. None of the other CI steps reads `host/workbench` or the workbench handler.

---

## §2 Execution protocol (binding, all milestones)

### 2.1 Environment (every shell call)

```bash
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter186 && export AILANG_BIN=$HOME/.pinned-ailang/ailang && <command>
```

The shell is zsh and does not keep `cwd` between calls. Run scripts with `bash`. Use `/usr/bin/grep`
(the default `grep` is ugrep). There is no `timeout(1)`. Never put a worktree or evidence under
`/tmp`; bank executor evidence under `~/.ailang/state/mission-world-iter186-evidence/exec/`.

### 2.2 Per-milestone order

1. Apply the milestone's edits (§3.2 / §4.2). **Preferred:** verify the banked half-diff's sha256
   and `git apply` it. **Fallback** (if `git apply --check` fails): make the edits by hand from the
   code blocks, which are copies of the same diff.
2. `gofmt -l <the milestone's files>` must print nothing.
3. Run the milestone exit gate (G-common plus the milestone's own). It must be green.
4. **Commit the milestone** (`git add <the milestone's two files> && git commit`). The drill restores
   each production file by `cp` from a backup it took just before mutating it, and its final
   `STATUS:` line must be `''`; that check only means something on a committed tree.
5. Run the milestone's drill rows (§2.3) and append each result line to the sprint Verification Log.
6. Confirm `git status --porcelain` is empty.

**G-common (every milestone exit):**

```bash
go vet ./... ; echo vet_rc=$?                                                          # AC1: vet_rc=0
go test ./... -count=1 2>&1 | /usr/bin/grep -E '^(ok|FAIL)' | awk '{print $1}' | sort | uniq -c   # AC2: 20 ok, 0 FAIL
gofmt -l host/daemon/workbench.go host/daemon/workbench_test.go host/workbench/render.go host/workbench/render_test.go   # prints nothing
go test -race ./host/workbench ./host/daemon -count=1                                  # both ok
```

`go build ./...` is **not** a compile fence for `_test.go` files; `go vet ./...` is. If AC2 reds only
on `TestHandlerTimeoutKillsTheWholeProcessGroup` or `TestCLIRealSubprocessEpisode` (known load
flakers, not in this sprint's packages), re-run that test alone once before calling it red; record
both runs.

### 2.3 The drill runner (banked outside the repo; do not re-derive)

```bash
E=~/.ailang/state/mission-world-iter186-evidence; X=$E/exec; mkdir -p $X
cp $E/planner/drill.py $X/drill.py && shasum -a 256 $X/drill.py   # must print 7c0d3718…3b820e
# usage: WT=<worktree> BK=<backup dir> python3 $X/drill.py [ID ...]   (no IDs = P0 + all rows)
```

Per row it does what AC7 requires: exact single-occurrence string replace (any other count prints
`NOT-LANDED count=N`); `go vet ./...` (rc≠0 prints `compiles=n … BUILD-FAIL` and is **never** a kill);
`go test ./host/workbench ./host/daemon ./host/boundary -count=1 -v`; the `--- FAIL:` set and its
leaves; `killer_in_reds=y|n`; `sole=y|n` (a subtest killer is sole when it is the only leaf; a
top-level killer is sole when every leaf is it or its subtest); restore by `cp` and assert the file's
sha256 equals the backup's. Last line: `STATUS: '<git status --porcelain>'`. `P0` runs the pristine
suite first. Row `H7x` is the doc's non-compiling first form of H7, kept as a compile-fence control:
it must print `compiles=n`.

A counted row **passes** only as `compiles=y rc=1 KILLED killer_in_reds=y`. Output shape (measured):

```
H5 compiles=y rc=1 KILLED killer_in_reds=y sole=y killer=TestWorkbenchObjectProvenanceWalk/interface-store-error reds=1 TestWorkbenchObjectProvenanceWalk/interface-store-error
H7x compiles=n vet_rc=1 BUILD-FAIL host/daemon/workbench.go:154:5: declared and not used: ok
```

---

## §3 M1 — the render guards (`host/workbench/render.go` +2/−2, `render_test.go` +61)

### 3.1 What M1 kills and why it is first

M1 closes the blank at the render layer (doc §2a layer (i), §2b guard). It can land alone: with
M1 only, a live object page shows `no grade reason was supplied` and
`no provenance edges were supplied for this object` — honest stops, not blanks. The whole repo is
green at M1 (measured, §1.1). R1–R5 are M1's; every H row's code does not exist yet.

### 3.2 Code

Apply:

```bash
E=~/.ailang/state/mission-world-iter186-evidence
shasum -a 256 $E/planner/m1.diff        # must print 3e064442…394eb3
git apply --check $E/planner/m1.diff && git apply $E/planner/m1.diff
git diff --numstat                      # 2 2 host/workbench/render.go / 61 0 host/workbench/render_test.go
```

What it contains. `render.go:154`, the grade line's `{{else}}` arm only:

```
before: {{else}}<p>GRADE UNAVAILABLE — {{.Grade.Unavailable}}</p>{{end}}
after:  {{else}}<p>GRADE UNAVAILABLE — {{with .Grade.Unavailable}}{{.}}{{else}}no grade reason was supplied{{end}}</p>{{end}}
```

`render.go:160`, the whole line:

```
before: {{with .Object}}{{range .Edges}}{{template "edge" .}}{{end}}{{end}}
after:  {{with .Object}}{{range .Edges}}{{template "edge" .}}{{else}}<p><span class="unavailable" role="note">UNAVAILABLE: no provenance edges were supplied for this object</span></p>{{end}}{{else}}<p><span class="unavailable" role="note">UNAVAILABLE: no object selected</span></p>{{end}}
```

Appended to the end of `render_test.go` (after `TestWorkbenchViewFieldsAllRender`, one blank line
between):

```go
// provenanceWalkBody returns the trimmed text between the provenance-walk
// section's </h2> and its </section>, failing if either marker is missing.
func provenanceWalkBody(t *testing.T, body string) string {
	t.Helper()
	const start = `<h2>Provenance walk</h2>`
	i := strings.Index(body, start)
	if i < 0 {
		t.Fatalf("no provenance-walk heading in %q", body)
	}
	rest := body[i+len(start):]
	j := strings.Index(rest, "</section>")
	if j < 0 {
		t.Fatalf("provenance-walk section is not closed in %q", body)
	}
	return strings.TrimSpace(rest[:j])
}

func TestRenderProvenanceWalkNeverBlank(t *testing.T) {
	for _, tc := range []struct {
		name, want, unwanted string
		page                 Page
	}{
		{"no-object", `<p><span class="unavailable" role="note">UNAVAILABLE: no object selected</span></p>`, "no provenance edges were supplied", Page{}},
		{"object-without-edges", `<p><span class="unavailable" role="note">UNAVAILABLE: no provenance edges were supplied for this object</span></p>`, "no object selected", Page{Object: &ObjectView{}}},
		{"supplied-edges", `<p>interface: <a href="/workbench?object=abc"`, "UNAVAILABLE", Page{Object: &ObjectView{Edges: []EdgeView{{Relation: "interface", Available: true, Target: "abc", Href: "?object=abc"}}}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			walk := provenanceWalkBody(t, renderPage(t, tc.page))
			if walk == "" {
				t.Fatal("provenance walk rendered a heading followed by nothing")
			}
			if !strings.Contains(walk, tc.want) {
				t.Errorf("provenance walk missing %q: %q", tc.want, walk)
			}
			if strings.Contains(walk, tc.unwanted) {
				t.Errorf("provenance walk contains %q: %q", tc.unwanted, walk)
			}
		})
	}
}

func TestRenderGradeReasonNeverEmpty(t *testing.T) {
	for _, tc := range []struct {
		name, want, unwanted string
		grade                GradeView
	}{
		{"zero-grade", `<p>GRADE UNAVAILABLE — no grade reason was supplied</p>`, `GRADE UNAVAILABLE — </p>`, GradeView{}},
		{"supplied-reason", `<p>GRADE UNAVAILABLE — named reason</p>`, "no grade reason was supplied", NewGradeUnavailable("named reason")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := renderPage(t, Page{Title: "grade", Object: &ObjectView{Grade: tc.grade}})
			if !strings.Contains(body, tc.want) {
				t.Errorf("grade line missing %q: %q", tc.want, body)
			}
			if strings.Contains(body, tc.unwanted) {
				t.Errorf("grade line contains %q: %q", tc.unwanted, body)
			}
		})
	}
}
```

The `—` is U+2014 (EM DASH), byte-identical to `render.go:154`. Copy, do not retype: a hyphen makes
`/zero-grade` and `/supplied-reason` red on the pristine tree (AC2 catches it at once).

### 3.3 M1 exit gate

G-common (§2.2), then:

```bash
go test ./host/workbench -count=1 -v -run 'TestRenderProvenanceWalkNeverBlank|TestRenderGradeReasonNeverEmpty' > $X/m1_ac.txt 2>&1; echo rc=$?   # rc=0
/usr/bin/grep -c -- '--- PASS' $X/m1_ac.txt           # 7
/usr/bin/grep -c -- '--- FAIL' $X/m1_ac.txt           # 0
git diff --quiet HEAD~1 -- host/daemon; echo daemon_rc=$?   # after the M1 commit: 0 (M1 touches no daemon file)
git diff --shortstat origin/dev -- host/              # 2 files changed, 63 insertions(+), 2 deletions(-)
```

### 3.4 M1 mutation drill (commit M1 first)

```bash
WT=/Users/voightkampff/dev/sunholo-data/.wt-world-iter186 BK=$X/bak-m1 python3 $X/drill.py P0 R1 R2 R3 R4 R5 > $X/drill_m1.txt 2>&1
```

| ID | Expected at M1 (measured, `drill_m1.txt`) |
|---|---|
| P0 | `pristine rc=0 reds=0` |
| R1 | KILLED, sole=y, reds=1: `TestRenderGradeReasonNeverEmpty/zero-grade` |
| R2 | KILLED, sole=n, reds=2: `…/supplied-reason`, `TestRenderGradeWithoutVerdictClaim/unavailable` |
| R3 | KILLED, sole=y, reds=1: `TestRenderProvenanceWalkNeverBlank/object-without-edges` |
| R4 | KILLED, sole=y, reds=1: `TestRenderProvenanceWalkNeverBlank/no-object` |
| R5 | KILLED, sole=n, reds=2: `…/supplied-edges`, `TestRenderUnavailableProvenanceEdge` |

Pass: 5 × `compiles=y rc=1 KILLED killer_in_reds=y`, last line `STATUS: ''`.

---

## §4 M2 — the daemon supplies both (`workbench.go` +46/−7, `workbench_test.go` +171)

### 4.1 What M2 kills

H1–H17. Measured (`drill_m1tests.txt`, M2 production with pre-M2 daemon tests): H1–H6, H9–H12 and
H14–H17 **survive** without M2's tests, so M2's tests are their only killers; H7, H8, H13 are
already killed by pre-existing selected-entry/paging subtests, which proves the `checkedEdge`
extraction keeps the row-35+38 pins (doc §5). After M2, every row is re-drilled (§4.5).

### 4.2 Code

Apply:

```bash
shasum -a 256 $E/planner/m2.diff        # must print 694db045…a5da101
git apply --check $E/planner/m2.diff && git apply $E/planner/m2.diff
git diff --numstat                      # 46 7 host/daemon/workbench.go / 171 0 host/daemon/workbench_test.go
```

What it contains, in `host/daemon/workbench.go`:

(a) A new `const` block **directly after** the message-constant block's closing `)` (`:32`) and
before `var acceptedWorkbenchKeys`, one blank line on each side:

```go
// Named reasons for what the object inspector cannot show. Each names the
// missing store fact; none is a guess at the value it stands in for.
const (
	objectGradeUnavailableReason = "no canonical host projection: this workbench handler does not resolve subject-bound evidence into an object grade"
	objectCommittedByMissing     = "the store records no commit-to-object relation, and the provenance field is a free-text label, not a reference"
	objectReferencedByMissing    = "no store index maps an object to the log entries or worlds that reference it"
)
```

(b) `entryEdges`'s loop body (`:132-141`) becomes the three lines below; its `refs` table, `make`,
comment and `return edges, nil` are unchanged:

```go
	for _, item := range refs {
		edge, err := d.checkedEdge(ctx, item.relation, item.ref)
		if err != nil {
			return nil, err
		}
		edges = append(edges, edge)
	}
```

(c) Two new methods **directly after** `entryEdges`'s closing `}` and before `func (d *Daemon) handleWorkbench`:

```go
// checkedEdge checks one edge target once. A stored target is a link; an
// unstored one is UNAVAILABLE with its ref visible; a store error is returned so
// the caller answers 5xx rather than "not stored".
func (d *Daemon) checkedEdge(ctx context.Context, relation string, ref hashref.HashRef) (workbench.EdgeView, error) {
	target := ref.String()
	_, ok, err := d.reads.GetObject(ctx, ref)
	if err != nil {
		return workbench.EdgeView{}, err
	}
	if !ok {
		return workbench.EdgeView{Relation: relation, Target: target, Missing: "object " + target + " is not stored"}, nil
	}
	return workbench.EdgeView{Relation: relation, Available: true, Target: target, Href: "?object=" + target}, nil
}

// objectEdges is the object's provenance walk from what the store records: the
// envelope's one typed reference, its interface, is existence-checked; the two
// relations the store cannot answer exactly are named stops, never a blank.
func (d *Daemon) objectEdges(ctx context.Context, object store.Object) ([]workbench.EdgeView, error) {
	iface, err := d.checkedEdge(ctx, "interface", object.InterfaceHash)
	if err != nil {
		return nil, err
	}
	return []workbench.EdgeView{
		iface,
		{Relation: "committedBy", Missing: objectCommittedByMissing},
		{Relation: "referencedBy", Missing: objectReferencedByMissing},
	}, nil
}
```

(d) The object block (`:258-268`): insert the `objectEdges` call between the truncation `if` and the
literal, and add **one** line to the literal (keep `Grade` and `Edges` on the same line — AC6 greps
and §5's H1/H2 anchors are measured against it):

```go
		edges, err := d.objectEdges(ctx, object)
		if err != nil {
			d.writeWorkbenchStoreError(w, r, ctx, err)
			return
		}
		page.Object = &workbench.ObjectView{
			Hash: object.Hash.String(), InterfaceHash: object.InterfaceHash.String(),
			SemanticID: object.SemanticID, Provenance: object.Provenance,
			PayloadShown: showPayload, PayloadPreview: string(preview), PayloadTruncated: truncated,
			Grade: workbench.NewGradeUnavailable(objectGradeUnavailableReason), Edges: edges,
		}
```

Appended to the end of `host/daemon/workbench_test.go` (after `TestWorkbenchPagingProbeStoreError`,
one blank line between):

```go
const provenanceWalkStart = `<section aria-label="provenance walk">`

// provenanceWalkSection returns the provenance-walk section of a workbench page
// up to its </section>, failing if the section is missing.
func provenanceWalkSection(t *testing.T, body string) string {
	t.Helper()
	section, ok := workbenchRegion(body, provenanceWalkStart, "</section>")
	if !ok {
		t.Fatalf("no provenance-walk section in %s", body)
	}
	return section
}

// refFailingStore fails GetObject for one ref only, so the object read itself
// succeeds while the walk's existence check on that ref errors.
type refFailingStore struct {
	readStore
	fail hashref.HashRef
}

func (s refFailingStore) GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error) {
	if ref == s.fail {
		return store.Object{}, false, errSentinelInternal
	}
	return s.readStore.GetObject(ctx, ref)
}

func workbenchTestObject(label, semanticID string, iface hashref.HashRef) store.Object {
	payload := []byte("payload-" + label)
	return store.Object{Hash: hashref.SumSHA256(payload), InterfaceHash: iface, SemanticID: semanticID, Provenance: "workbench-test", Payload: payload}
}

func TestWorkbenchObjectGrade(t *testing.T) {
	d := newHandlerDaemon(t)
	genesis := seedGenesisEmbedded(t, d, "workbench-grade")
	commit := testCommit(genesis, 0, "workbench-grade")
	proof := workbenchTestObject("workbench-grade-proof", "world/proof-report/v1", hashref.SumSHA256([]byte("world/authenticated-proof-envelope/v1")))
	commit.Objects = append(commit.Objects, proof)
	if err := d.store.Commit(commit); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	// CONTROL: the registry object is the real bootstrap-written one, not a fixture.
	registry, ok, err := d.store.GetRegistryHead(context.Background(), store.EpochRegistryV1)
	if err != nil || !ok {
		t.Fatalf("control GetRegistryHead(%s): ok=%v err=%v", store.EpochRegistryV1, ok, err)
	}
	want := `<p>GRADE UNAVAILABLE — ` + objectGradeUnavailableReason + `</p>`
	for _, object := range []struct {
		name string
		ref  hashref.HashRef
	}{{"test-object", commit.Objects[0].Hash}, {"proof-labelled", proof.Hash}, {"registry", registry}} {
		for _, suffix := range []struct{ name, query string }{{"", ""}, {"-payload", "&payload=1"}} {
			t.Run(object.name+suffix.name, func(t *testing.T) {
				rec := requestRecorder(t, d, http.MethodGet, "/workbench?object="+object.ref.String()+suffix.query, nil)
				if rec.Code != http.StatusOK {
					t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body)
				}
				body := rec.Body.String()
				if got := strings.Count(body, want); got != 1 {
					t.Errorf("grade line %q appears %d times, want 1", want, got)
				}
				for _, unwanted := range []string{`GRADE UNAVAILABLE — </p>`, "no grade reason was supplied", `<span>PROVEN</span>`, `<span>TESTED</span>`, `<span>ATTESTED</span>`, `<span>CLAIMED</span>`} {
					if strings.Contains(body, unwanted) {
						t.Errorf("object page contains %q", unwanted)
					}
				}
			})
		}
	}
}

func TestWorkbenchObjectProvenanceWalk(t *testing.T) {
	d := newHandlerDaemon(t)
	genesis := seedGenesisEmbedded(t, d, "workbench-walk")
	commit := testCommit(genesis, 0, "workbench-walk")
	plain := commit.Objects[0]
	schema := workbenchTestObject("workbench-walk-schema", "test/schema", hashref.SumSHA256([]byte("interface-workbench-walk-schema")))
	typed := workbenchTestObject("workbench-walk-typed", "test/typed", schema.Hash)
	commit.Objects = append(commit.Objects, schema, typed)
	if err := d.store.Commit(commit); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	// CONTROLS: one interface target is stored and one is not, or the link and
	// UNAVAILABLE arms below assert nothing.
	if _, ok, err := d.store.GetObject(context.Background(), schema.Hash); err != nil || !ok {
		t.Fatalf("control GetObject(schema): ok=%v err=%v, want ok=true", ok, err)
	}
	if _, ok, err := d.store.GetObject(context.Background(), plain.InterfaceHash); err != nil || ok {
		t.Fatalf("control GetObject(plain.InterfaceHash): ok=%v err=%v, want ok=false", ok, err)
	}
	get := func(t *testing.T, target string) string {
		t.Helper()
		rec := requestRecorder(t, d, http.MethodGet, target, nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status = %d, want 200; body=%s", target, rec.Code, rec.Body)
		}
		return rec.Body.String()
	}
	plainTarget := "/workbench?object=" + plain.Hash.String()
	typedTarget := "/workbench?object=" + typed.Hash.String()

	t.Run("interface-stored-link", func(t *testing.T) {
		section := provenanceWalkSection(t, get(t, typedTarget))
		href := "/workbench?object=" + schema.Hash.String()
		if want := `<p>interface: <a href="` + href + `"`; !strings.Contains(section, want) {
			t.Fatalf("walk missing stored interface link %q: %s", want, section)
		}
		get(t, href)
	})

	t.Run("interface-unstored-unavailable", func(t *testing.T) {
		body := get(t, plainTarget)
		iface := plain.InterfaceHash.String()
		want := `<p>interface: <span class="unavailable" role="note">UNAVAILABLE: object ` + iface + ` is not stored</span></p>`
		if section := provenanceWalkSection(t, body); !strings.Contains(section, want) {
			t.Errorf("walk missing %q: %s", want, section)
		}
		if unwanted := `href="/workbench?object=` + iface + `"`; strings.Contains(body, unwanted) {
			t.Errorf("unstored interface rendered as a link %q", unwanted)
		}
	})

	t.Run("named-stops", func(t *testing.T) {
		for _, target := range []string{plainTarget, typedTarget + "&payload=1"} {
			section := provenanceWalkSection(t, get(t, target))
			for _, want := range []string{
				`<p>committedBy: <span class="unavailable" role="note">UNAVAILABLE: ` + objectCommittedByMissing + `</span></p>`,
				`<p>referencedBy: <span class="unavailable" role="note">UNAVAILABLE: ` + objectReferencedByMissing + `</span></p>`,
			} {
				if !strings.Contains(section, want) {
					t.Errorf("%s: walk missing named stop %q: %s", target, want, section)
				}
			}
			if strings.Contains(section, "no provenance edges were supplied") {
				t.Errorf("%s: daemon object page fell back to the render-layer stop: %s", target, section)
			}
		}
	})

	t.Run("never-blank", func(t *testing.T) {
		for _, target := range []string{"/workbench", "/workbench?from=0&entry=0", plainTarget, typedTarget + "&payload=1"} {
			section := provenanceWalkSection(t, get(t, target))
			_, after, ok := strings.Cut(section, "</h2>")
			if !ok || strings.TrimSpace(after) == "" {
				t.Errorf("%s: provenance walk is a heading followed by nothing: %q", target, section)
			}
			if strings.Contains(section, "UNAVAILABLE: </span>") {
				t.Errorf("%s: walk rendered a stop with an empty reason: %s", target, section)
			}
		}
	})

	t.Run("interface-store-error", func(t *testing.T) {
		oldReads, oldErrLog := d.reads, d.errLog
		defer func() { d.reads, d.errLog = oldReads, oldErrLog }()
		d.reads = refFailingStore{readStore: oldReads, fail: schema.Hash}
		d.errLog = &bytes.Buffer{}
		// CONTROL: the same store serves a page whose walk does not read schema.
		if rec := requestRecorder(t, d, http.MethodGet, plainTarget, nil); rec.Code != http.StatusOK {
			t.Fatalf("control: %s status = %d, want 200; body=%s", plainTarget, rec.Code, rec.Body)
		}
		rec := requestRecorder(t, d, http.MethodGet, typedTarget, nil)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500; body=%s", rec.Code, rec.Body)
		}
		if !strings.Contains(rec.Body.String(), ">Internal<") {
			t.Errorf("body does not contain >Internal<: %s", rec.Body)
		}
	})
}
```

Keep every assertion string byte-for-byte: the §5 red sets were measured against them. The
absence assertions (`GRADE UNAVAILABLE — </p>` in `TestWorkbenchObjectGrade`, `UNAVAILABLE: </span>`
in `/never-blank`) are load-bearing for H3/H16 and H14/H15 respectively (doc §7 closing note) and
must not be dropped to save lines.

### 4.3 Size check

`git diff --shortstat origin/dev -- host/` → `4 files changed, 280 insertions(+), 9 deletions(-)`;
numstat `46 7 host/daemon/workbench.go`, `171 0 host/daemon/workbench_test.go`,
`2 2 host/workbench/render.go`, `61 0 host/workbench/render_test.go` (measured). AC8: 289 ≤ 330.

### 4.4 M2 exit gate = the full AC set (§6)

G-common, `bash ./scripts/verify_ail.sh`, `bash scripts/check_no_personal_email.sh`, and every AC in §6.

### 4.5 M2 mutation drill (commit M2 first)

Step 1 — every row on the committed M2 tree (AC7):

```bash
WT=/Users/voightkampff/dev/sunholo-data/.wt-world-iter186 BK=$X/bak-full python3 $X/drill.py > $X/drill_full.txt 2>&1
/usr/bin/grep -c ' KILLED killer_in_reds=y ' $X/drill_full.txt     # 22
/usr/bin/grep -c ' sole=y ' $X/drill_full.txt                       # 14
/usr/bin/grep -E '^H7x compiles=n ' $X/drill_full.txt | wc -l      # 1 (the control; never a kill)
/usr/bin/grep -cE 'SURVIVED|NOT-LANDED|killer_in_reds=n' $X/drill_full.txt   # 0
head -1 $X/drill_full.txt                                           # P0 pristine rc=0 reds=0
tail -1 $X/drill_full.txt                                           # STATUS: ''
```

Step 2 — the base-suite column (the doc's "Base suite at `6127ff3`"). Put `origin/dev`'s two test
files in place **by cp from a backup**, drill, restore:

```bash
mkdir -p $X/testbak && cp host/daemon/workbench_test.go host/workbench/render_test.go $X/testbak/
git show origin/dev:host/daemon/workbench_test.go > host/daemon/workbench_test.go
git show origin/dev:host/workbench/render_test.go > host/workbench/render_test.go
WT=/Users/voightkampff/dev/sunholo-data/.wt-world-iter186 BK=$X/bak-base python3 $X/drill.py > $X/drill_basesuite.txt 2>&1
cp $X/testbak/workbench_test.go host/daemon/workbench_test.go && cp $X/testbak/render_test.go host/workbench/render_test.go
git status --porcelain                                              # empty
/usr/bin/grep -c 'compiles=y rc=0 SURVIVED' $X/drill_basesuite.txt  # 17
/usr/bin/grep 'compiles=y rc=1' $X/drill_basesuite.txt | cut -d' ' -f1 | tr '\n' ' '   # R2 R5 H7 H8 H13
```

(`drill_basesuite.txt`'s own `STATUS:` line shows the two test files as modified — expected, the swap
is uncommitted; the `git status --porcelain` after the restore is the check.) If `origin/dev` has
moved past `6127ff3` and the counts differ, stop and report; do not absorb.

---

## §5 The mutation table (measured, 2026-09-25; `drill.py` rows)

Old text is exact and occurs **once** in the landed file; `⏎` is a newline and `⇥` a tab in the
file. Owner = the milestone whose tests kill it. "At M1" = result at the M1-only tree; "M2 prod +
pre-M2 tests" = `drill_m1tests.txt`; "Base" = production M1+M2 with both test files at `6127ff3`.

| ID | Owner | File | Old → new | Compiles | Suite rc (full) | Named killer (in red set) | Sole | Full red leaves | At M1 | M2 prod + pre-M2 tests | Base |
|---|---|---|---|---|---|---|---|---|---|---|---|
| R1 | M1 | render.go | `{{with .Grade.Unavailable}}{{.}}{{else}}no grade reason was supplied{{end}}` → `{{.Grade.Unavailable}}` | y | 1 | `TestRenderGradeReasonNeverEmpty/zero-grade` | y | 1 | KILLED (sole) | KILLED | 0 survives |
| R2 | M1 | render.go | same old → `no grade reason was supplied` | y | 1 | `TestRenderGradeReasonNeverEmpty/supplied-reason` | n | 8: + `TestRenderGradeWithoutVerdictClaim/unavailable`, 6× `TestWorkbenchObjectGrade/*` | KILLED (2) | KILLED (2) | 1 (`…WithoutVerdictClaim/unavailable`) |
| R3 | M1 | render.go | `{{else}}<p><span class="unavailable" role="note">UNAVAILABLE: no provenance edges were supplied for this object</span></p>` → `` | y | 1 | `TestRenderProvenanceWalkNeverBlank/object-without-edges` | y | 1 | KILLED (sole) | KILLED | 0 survives |
| R4 | M1 | render.go | `{{else}}<p><span class="unavailable" role="note">UNAVAILABLE: no object selected</span></p>` → `` | y | 1 | `TestRenderProvenanceWalkNeverBlank/no-object` | n | 2: + `TestWorkbenchObjectProvenanceWalk/never-blank` | KILLED (sole) | KILLED (sole) | 0 survives |
| R5 | M1 | render.go | `{{range .Edges}}{{template "edge" .}}{{else}}` → `{{range .Edges}}{{else}}` | y | 1 | `TestRenderProvenanceWalkNeverBlank/supplied-edges` | n | 6: + `TestRenderUnavailableProvenanceEdge`, walk `/interface-stored-link`, `/interface-unstored-unavailable`, `/named-stops`, `/never-blank` | KILLED (2) | KILLED (2) | 1 (`TestRenderUnavailableProvenanceEdge`) |
| H1 | M2 | workbench.go | `Edges: edges,` → `Edges: edges[:0],` | y | 1 | `TestWorkbenchObjectProvenanceWalk/named-stops` | n | 3: + `/interface-stored-link`, `/interface-unstored-unavailable` | n/a | survives | 0 survives |
| H2 | M2 | workbench.go | `Grade: workbench.NewGradeUnavailable(objectGradeUnavailableReason), ` (trailing space) → `` | y | 1 | `TestWorkbenchObjectGrade` | y | 6 (its subtests) | n/a | survives | 0 survives |
| H3 | M2 | workbench.go | `NewGradeUnavailable(objectGradeUnavailableReason)` → `NewGradeUnavailable("")` | y | 1 | `TestWorkbenchObjectGrade` | y | 6 | n/a | survives | 0 survives |
| H4 | M2 | workbench.go | `workbench.NewGradeUnavailable(objectGradeUnavailableReason)` → `workbench.GradeView{Available: true, Label: workbench.GradeCLAIMED}` | y | 1 | `TestWorkbenchObjectGrade` | y | 6 | n/a | survives | 0 survives |
| H5 | M2 | workbench.go | `iface, err := d.checkedEdge(ctx, "interface", object.InterfaceHash)⏎⇥if err != nil {` → `…⏎⇥if false && err != nil {` | y | 1 | `TestWorkbenchObjectProvenanceWalk/interface-store-error` | y | 1 | n/a | survives | 0 survives |
| H6 | M2 | workbench.go | `edges, err := d.objectEdges(ctx, object)⏎⇥⇥if err != nil {` → `…⏎⇥⇥if false && err != nil {` | y | 1 | `…/interface-store-error` | y | 1 | n/a | survives | 0 survives |
| H7 | M2 | workbench.go | `if !ok {⏎⇥⇥return workbench.EdgeView{Relation: relation, Target: target,` → `if false && !ok {⏎⇥⇥return workbench.EdgeView{Relation: relation, Target: target,` | y | 1 | `…/interface-unstored-unavailable` | n | 3: + `TestWorkbenchSelectedEntry/unstored-edge-unavailable`, `TestWorkbenchTimelinePaging/emitted-links-resolve` | n/a | KILLED (2, pre-existing only) | 1 |
| H7x | control | workbench.go | same old → `if false {⏎⇥⇥return …` (doc's first form) | **n** (`declared and not used: ok`) | — | — (never a kill) | — | — | n/a | n | n |
| H8 | M2 | workbench.go | `_, ok, err := d.reads.GetObject(ctx, ref)⏎⇥if err != nil {` → `…⏎⇥if false && err != nil {` | y | 1 | `…/interface-store-error` | n | 2: + `TestWorkbenchSelectedEntry/object-store-error` | n/a | KILLED (1, pre-existing only) | 1 |
| H9 | M2 | workbench.go | `d.checkedEdge(ctx, "interface", object.InterfaceHash)` → `d.checkedEdge(ctx, "interface", object.Hash)` | y | 1 | `…/interface-stored-link` | n | 3: + `/interface-unstored-unavailable`, `/interface-store-error` | n/a | survives | 0 survives |
| H10 | M2 | workbench.go | `⇥⇥{Relation: "committedBy", Missing: objectCommittedByMissing},⏎` → `` | y | 1 | `…/named-stops` | y | 1 | n/a | survives | 0 survives |
| H11 | M2 | workbench.go | `⇥⇥{Relation: "referencedBy", Missing: objectReferencedByMissing},⏎` → `` | y | 1 | `…/named-stops` | y | 1 | n/a | survives | 0 survives |
| H12 | M2 | workbench.go | `{Relation: "referencedBy", Missing: objectReferencedByMissing}` → `{Relation: "referencedBy", Available: true, Target: object.Hash.String(), Href: "?object=" + object.Hash.String()}` | y | 1 | `…/named-stops` | y | 1 | n/a | survives | 0 survives |
| H13 | M2 | workbench.go | `edge, err := d.checkedEdge(ctx, item.relation, item.ref)⏎⇥⇥if err != nil {` → `…⏎⇥⇥if false && err != nil {` | y | 1 | `TestWorkbenchSelectedEntry/object-store-error` | y | 1 | n/a | KILLED (sole) | 1 (named killer is pre-existing) |
| H14 | M2 | workbench.go | `objectCommittedByMissing     = "the store records no commit-to-object relation, and the provenance field is a free-text label, not a reference"` → `objectCommittedByMissing     = ""` | y | 1 | `…/never-blank` | y | 1 | n/a | survives | 0 survives |
| H15 | M2 | workbench.go | `objectReferencedByMissing    = "no store index maps an object to the log entries or worlds that reference it"` → `objectReferencedByMissing    = ""` | y | 1 | `…/never-blank` | y | 1 | n/a | survives | 0 survives |
| H16 | M2 | workbench.go | `objectGradeUnavailableReason = "no canonical host projection: this workbench handler does not resolve subject-bound evidence into an object grade"` → `objectGradeUnavailableReason = ""` | y | 1 | `TestWorkbenchObjectGrade` | y | 6 | n/a | survives | 0 survives |
| H17 | M2 | workbench.go | `"interface", object.InterfaceHash` → `"transitionRef", object.InterfaceHash` | y | 1 | `…/interface-stored-link` | n | 2: + `/interface-unstored-unavailable` | n/a | survives | 0 survives |

`…/` = `TestWorkbenchObjectProvenanceWalk/`. **Totals (counted rows, H7x excluded): 22 mutants;
22 compile; 22 rc=1; 22 with the named killer in the red set; 14 sole (R1, R3, H2, H3, H4, H5, H6,
H10, H11, H12, H13, H14, H15, H16).** Per milestone: **M1 owns 5** (R1–R5; 5/5 killed at M1, 3 sole
at M1: R1, R3, R4; at the full tree R4 gains `/never-blank`, so 2 sole: R1, R3). **M2 owns 17**
(H1–H17; 17/17 killed by the named killer at M2, 12 sole). Base suite: 17 survive `6127ff3`'s tests,
5 (R2, R5, H7, H8, H13) are killed by a pre-existing test — identical to the doc's §7.

Load-bearing claims and their mutations (coding-standards S6): the render `{{with}}` guard ← R1
(`/zero-grade`), its pass-through arm ← R2 (`/supplied-reason`); the `range`-else ← R3; the
`with`-else ← R4; the `range` body ← R5; the literal's `Edges` ← H1; the literal's `Grade` ← H2,
H3, H4, H16 (the grade-line count of exactly 1, and the absence of `GRADE UNAVAILABLE — </p>` and of
any `<span>LABEL</span>`); the walk's store-error propagation ← H5, H6, H8; unstored ⇒ UNAVAILABLE ←
H7; the interface target ← H9, H17; the named stops ← H10, H11, H12; their non-empty reasons ← H14,
H15 (the `UNAVAILABLE: </span>` absence); the `entryEdges` pin through the refactor ← H13. The AC6
greps and the AC4 PASS counts are instrument-health controls, never load-bearing claims.

---

## §6 Acceptance criteria → command → expected (the controller checklist)

| AC | Command | Expected (planner-measured on the prototype) |
|---|---|---|
| AC1 | `go vet ./... ; echo rc=$?` | rc=0 |
| AC2 | `go test ./... -count=1` | rc=0; 20 `ok`, 0 `FAIL` |
| AC3 | `bash ./scripts/verify_ail.sh ; echo rc=$?` · `git diff --name-only origin/dev -- '*.ail' \| wc -l` | rc=0; `✓ verify gate PASSED: 11 required identities verified, 40 named tests pass` (and `world package gate PASSED: 9/9`) · `0` |
| AC4 | `go test ./host/workbench ./host/daemon -count=1 -v -run 'TestRenderProvenanceWalkNeverBlank\|TestRenderGradeReasonNeverEmpty\|TestWorkbenchObjectGrade\|TestWorkbenchObjectProvenanceWalk'` | rc=0; **20** `--- PASS` (4 at column 0, 16 indented), 0 `--- FAIL`; the 16 are `TestRenderProvenanceWalkNeverBlank/{no-object,object-without-edges,supplied-edges}`, `TestRenderGradeReasonNeverEmpty/{zero-grade,supplied-reason}`, `TestWorkbenchObjectGrade/{test-object,test-object-payload,proof-labelled,proof-labelled-payload,registry,registry-payload}`, `TestWorkbenchObjectProvenanceWalk/{interface-stored-link,interface-unstored-unavailable,named-stops,never-blank,interface-store-error}` |
| AC5 | `bash -c 'diff <(git show origin/dev:host/daemon/workbench.go \| sed -n "/^func supportedWorkbenchQuery/,/^}/p") <(sed -n "/^func supportedWorkbenchQuery/,/^}/p" host/daemon/workbench.go)'`, the same for `/^var acceptedWorkbenchKeys/,/^}/`, and for `/^func TestWorkbenchRefusalBranches/,/^}/` on `workbench_test.go` | rc=0 ×3 |
| AC6 | on `host/daemon/workbench.go`: `grep -c 'Grade: workbench.NewGradeUnavailable(objectGradeUnavailableReason)'` · `grep -c 'objectGradeUnavailableReason = "no canonical host projection: this workbench handler does not resolve subject-bound evidence into an object grade"'` · `grep -c 'computed only by gradeOf'` · `grep -c 'Edges: edges,'` · `grep -c NewGradeView` · control `grep -c PayloadTruncated` | `1` · `1` · `0` · `1` · `0` · `1` (instrument-health only) |
| AC7 | §3.4 at M1 and §4.5 steps 1–2 at M2 | M1: R1–R5 KILLED w/ killer. M2: 22 KILLED w/ killer, 14 sole, H7x `compiles=n`, 0 SURVIVED/NOT-LANDED, `STATUS: ''`; base suite 17 survive, 5 killed (R2 R5 H7 H8 H13) |
| AC8 | `git diff --shortstat origin/dev -- host/` | `4 files changed, 280 insertions(+), 9 deletions(-)` (289 ≤ 330) |
| AC9 | `git diff origin/dev -- host/daemon/workbench.go \| /usr/bin/grep -c StateRoot` | `0` |
| gofmt | `gofmt -l host/daemon/workbench.go host/daemon/workbench_test.go host/workbench/render.go host/workbench/render_test.go` | empty |
| read budget | `/usr/bin/grep -c GetObject host/daemon/workbench.go` · `/usr/bin/grep -c context.Background host/daemon/workbench.go` | `2` (base `2`) · `0` |

Also run: `go test -race ./host/workbench ./host/daemon -count=1` (both `ok`) and
`bash scripts/check_no_personal_email.sh` (✓). Do **not** reformat pre-existing gofmt-dirty files.

---

## §7 Estimate and risks

- **Size:** 289 changed lines (57 production, 232 test), 4 files (measured). **Time:** ~0.5 day
  including the drills (about 80 s per full drill on this rig).
- **Risk 1 — character drift.** The em dash (U+2014) in the grade strings must be byte-identical to
  `render.go:154`. Mitigation: `git apply` the banked diffs; if hand-editing, copy the blocks.
- **Risk 2 — drilling an uncommitted milestone.** `STATUS: ''` is the only proof the tree is clean,
  and it only reads `''` on a committed milestone. Commit first (§2.2 step 4). Never restore with
  `git checkout --`.
- **Risk 3 — layout drift breaks the drill anchors.** Every §5 old text is whitespace-exact (tabs
  and newlines). If the executor reflows a line (for example splits `Grade: …, Edges: edges,` onto
  two lines, or lets `gofmt` realign the const block), rows go `NOT-LANDED`. A NOT-LANDED row is a
  finding, never a pass: restore the §4.2 layout rather than editing the anchor.
- **Risk 4 — load flakes** in `TestHandlerTimeoutKillsTheWholeProcessGroup` /
  `TestCLIRealSubprocessEpisode` (other packages). Re-run alone once; record both.
- **Risk 5 — `origin/dev` moves** before §4.5 step 2. A changed base-suite result is reported, not
  absorbed.
- **Risk 6 — M1 alone on a live page.** Between the M1 and M2 commits, object pages show the render
  fallbacks (`no grade reason was supplied`, `no provenance edges were supplied for this object`).
  That is by design (doc §8) and the whole suite is green there (measured); do not "fix" it in M1.
- **CI-only steps** (§1.3) are not reproduced locally. The binding gate is CI green on the PR head
  and on the merge commit.
- **Out of scope** (doc §9, unchanged): row 99 (world `StateRoot` unchecked — `checkedEdge` is the
  helper it will reuse), row 100 (ratchet over `ObjectView`), row 96, and the proposed P-a
  (object-reference index), P-b (commit membership), P-c (grade projection).
