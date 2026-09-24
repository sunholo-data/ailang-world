# w-workbench-unpinned-render-hunks — sprint plan (row 34, 2 milestones, tests only)

- Sprint: `w-workbench-unpinned-render-hunks` · queue row **34** (clause 5) · design doc
  [`w-workbench-unpinned-render-hunks.md`](w-workbench-unpinned-render-hunks.md) (quorum-closed:
  r1 BLOCKED → revision 1; r2 gemini PASS, glm PASS, astra REJECT closed by the narrow-refinement
  carve-out at `aa5df54`, doc §11).
- Base: `aa5df54` (`git rev-parse HEAD` in the sprint worktree) · branch
  `sprint/w-workbench-unpinned-render-hunks` · worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter185`.
- `origin/dev` = `fd99840` (the doc's measurement base). `git diff --quiet fd99840 aa5df54 -- host cmd scripts`
  → rc=0 (measured), so every `host/` anchor the doc took at `fd99840` is also an anchor at base. The
  three commits between them are the design doc alone.
- Planned 2026-09-24 by mission-control iteration 185 (sprint-planner). **The whole design was
  prototyped by the planner** in three detached scratch worktrees at `aa5df54`, siblings of the repo
  (`.planner-wt-iter185` = full r2 prototype, `-m1` = M1 half only, `-base` = pristine tests), all
  removed afterwards with `git worktree remove --force`. The prototype is the controller's
  `~/.ailang/state/mission-world-iter185-evidence/prototype_r2_controller.diff` applied unchanged.
  Every code block below is copied from that diff, which applied cleanly, was `gofmt`-clean, vetted
  and passed. It is not a paraphrase of the doc.
- Files changed by the sprint (doc §5, confirmed by `git diff --numstat` on the prototype):
  `host/daemon/workbench_test.go` (+37/−0) and `host/workbench/render_test.go` (+53/−0). **Nothing
  else.** No production file, no `.ail` file.
- Planner evidence: `~/.ailang/state/mission-world-iter185-evidence/planner/` — the drill runner
  `mutdrive_planner.py` (sha256 `a4c00ade…fcfa4c5`), the three drill transcripts `drill_base.txt`,
  `drill_m1.txt`, `drill_full.txt`, the AC4 transcripts `ac4.txt` / `m1_ac.txt`, the whole-repo test
  transcripts `full_test.txt` / `m1_full_test.txt`, and `verify_ail.txt`.

---

## §0 Planner findings (measured; read first)

The prototype **confirms every §6 acceptance number and every §7 red count** in the doc (§1.2
below: 57/57 rows match, before and after). No mutation row fails to compile, and no named killer
is missing from its row's red set. What remains are stale numbers in the doc's prose. None changes
what the executor builds. Each is recorded here with its measurement. The doc is **not** edited by
this plan.

| # | Doc says | Measured | Resolution in this plan |
|---|---|---|---|
| F1 | Header (line 9): "+84 insertions across 2 test files". §3.2 heading: "`render_test.go` (+47)". V14: numstat `47 0 host/workbench/render_test.go`. §8 M2: "~47". | The r2 carve-out added the six-line negative verdict-claim loop to `TestRenderGradeWithoutVerdictClaim`. `git diff --numstat` on the r2 prototype: `37 0 host/daemon/workbench_test.go`, **`53 0 host/workbench/render_test.go`**. `git diff --shortstat`: **`2 files changed, 90 insertions(+)`**. | The numbers in this plan are 37 / 53 / 90. AC7 (≤ 120 insertions, 0 deletions) still holds with a margin of 30. The header's "~90 LOC" figure is right; only its "+84" breakdown is stale. |
| F2 | §8 M1 exit: "H5–H7, Q17–Q20, Q23, Q26, Q27, Q30 drilled". §8 M2 exit: "all **56** rows drilled". | The table has **57** rows after revision 2 (V10 added); the drill runner's table has 57 entries. The M1 killable set is 10 rows (H5, H6, H7, Q18, Q19, Q20, Q23, Q26, Q27, Q30). Q17 is the equivalence control, not a kill. | M1 drills the 10 rows plus the Q17 control (§3.5). M2 drills its 10 rows, then the **full 57-row** drill (§4.5). |
| F3 | Header line 5, §1.3 and the first §7 tally: "56 mutants … 20 survive … 19 killable". | Superseded by revision 2 (header line 7 and the second §7 tally: 57 / 36 / 21 / 56 + Q17). The planner's pristine-tests drill measured **36 KILLED, 21 SURVIVED** (the 20 listed in AC6 plus Q17), 0 `compiles=n`, 0 NOT-LANDED. | The plan uses the revision-2 totals. The executor should not quote the line-5 totals in the PR body. |
| F4 | §2f's list of rows whose only killers are the new tests omits **V10**. | V10's red set after M2 is exactly `TestRenderGradeWithoutVerdictClaim`, `/no-verdict`, `/unavailable` (3 reds). It **survives** the pristine suite and the M1 state (measured in both drills). | V10 is an M2 row, named killer = both subtests (§4.5). |
| F5 | AC6 uses `go build ./...` as the compile fence. | Correct for this sprint, and it is not the known `go build` blind spot: every mutant edits a **production** file (`render.go`, `workbench.go`), and the test files are compiled by the `go test` step, whose failure would show up as a `[build failed]` red, not a kill with a named test. The per-milestone gate still uses `go vet ./...` (AC1) for the test files the executor writes. | No change. Stated so an evaluator does not flag it. |

Anchors re-verified at base (`sed -n`/`cat -n` on the files at `aa5df54`, all exact):
`workbench.go:34-40` (`acceptedWorkbenchKeys`, the r2 glm note), `:63-77`
(`supportedWorkbenchQuery`; `:70`, `:73`, `:76`); `render.go:97-102` (`workbenchHref`; guard `:98`),
`:122` (CSS selectors `.verdict-fail{`/`.verdict-pass`, no `class="` prefix), `:129`
(`workbenchHref ""`), `:130`, `:154` (grade line), `:155`; `workbench_test.go:134-216`
(`TestWorkbenchRefusalBranches`; `unsupported-combination` at `:157`; `commitWorkbenchPayload` at
`:218`); `render_test.go:60-96` (`TestGradeViewRequiresTestVerdict`; `/pass` at `:86-95`, binding
`view` at `:88`, the gemini r2 note). Both test files already import `bytes` and `strings`; no
import change is needed.

---

## §1 Baseline and prototype measurements

### 1.1 Commands and results (2026-09-24, load average 1.64 at start)

All commands ran with `export AILANG_BIN=$HOME/.pinned-ailang/ailang`
(`$AILANG_BIN --version` → `AILANG v0.41.0`; `go version` → `go1.26.6 darwin/arm64`).

| Gate | Command | Full r2 prototype (M1+M2) | M1 half only |
|---|---|---|---|
| vet | `go vet ./... ; echo rc=$?` | **rc=0** | **rc=0** |
| test | `go test ./... -count=1` | **rc=0**; **20 `ok`**, 0 `FAIL` (`go list ./... \| wc -l` = 20) | **rc=0**; 20 `ok`, 0 `FAIL` |
| gofmt (scoped) | `gofmt -l host/daemon/workbench_test.go host/workbench/render_test.go` | empty | empty |
| gofmt (whole) | `gofmt -l host cmd` | only the 3 pre-existing files (`host/daemon/session_middleware_test.go`, `host/store/schema_version_test.go`, `host/store/store.go`) — **red at base, out of scope, do not touch** | same |
| .ail | `./scripts/verify_ail.sh ; echo rc=$?` | **rc=0**; `✓ world package gate PASSED: 9/9 steps performed non-zero work`; `✓ verify gate PASSED: 11 required identities verified, 40 named tests pass`; `git diff --name-only aa5df54 -- '*.ail' \| wc -l` → `0` | — |
| AC4 | the §6 AC4 `-run` command | rc=0; **59** `--- PASS`, 0 `--- FAIL`; **32** `--- PASS: TestSupportedWorkbenchQueryTruthTable/`; 17 `--- PASS: TestWorkbenchRefusalBranches/` | M1 subset (§3.4): rc=0; **51** `--- PASS`; 32; 17; 0 FAIL |
| race (touched pkgs) | `go test -race ./host/workbench ./host/daemon -count=1` | both `ok` (workbench 1.4 s, daemon 12.7 s) | — |
| email | `bash scripts/check_no_personal_email.sh` | `✓ check-no-personal-email: no personal addresses in the loop-written surface` | — |
| size | `git diff --shortstat` / `--numstat` | `2 files changed, 90 insertions(+)`; 37 / 53 | `1 file changed, 37 insertions(+)` |

### 1.2 The three mutation drills (57 rows each, `mutdrive_planner.py`, §2.3)

| Drill | Tree | KILLED | SURVIVED | `compiles=n` | NOT-LANDED | named killer in reds (every KILLED row) | Survivors |
|---|---|---|---|---|---|---|---|
| base | pristine tests | 36 | **21** | 0 | 0 | y | H3 H4 H5 H6 H7 Q17 Q18 Q19 Q20 Q23 Q26 Q27 Q30 W5 W8 V2 V4 V6 V7 V9 V10 |
| m1 | + M1 only | 46 | **11** | 0 | 0 | y | H3 H4 Q17 W5 W8 V2 V4 V6 V7 V9 V10 |
| full | + M1 + M2 | **56** | **1** | 0 | 0 | y | Q17 |

Per-row red counts: the base drill's counts equal the doc's §7 "Before" column on all 57 rows, and
the full drill's counts equal the §7 "After" column on all 57 rows (compared by script, 0
differences). The full drill independently reproduces the controller's `r2/drill_after_r2.txt`
(56 / 1 / 0). Each run ended with production files restored (`STATUS` showed only the prototype's
test-file edits).

HTTP-witness co-kills, measured in the m1 drill (doc §2b claims H5–H7, Q18, Q20, Q23, Q26, Q30):
H5 → `unsupported-triple-world-object-payload`; H6 → `unsupported-pair-world-from`; H7 →
`unsupported-pair-world-payload`; Q18 → triple; Q20 → pair-world-from; Q23 and Q30 → both pairs;
Q26 → pair-world-payload. Q19 and Q27 are killed by truth-table subtests only. **Confirmed.**

### 1.3 CI's command list (the gate list)

`.github/workflows/ci.yml` has two jobs. Setup steps (checkout, setup-go 1.26.6, binary and Z3
installs, the version asserts) are omitted. The commands each job **runs as its gate** are:

**Job `ailang-verify`** (`ailang-code verify gate`):
1. `ailang --version`
2. `./scripts/verify_ail.sh`

**Job `go-verify`** (`go host build + test gate`):
1. `go version` / `go env GOVERSION`
2. `./scripts/verify_go.sh` (timeout 25 min). Its internals: `AILANG_BIN` v0.41.0 assert, the
   tracked-binary hygiene gate, the World decision-ledger check, the toolchain miscompile guard,
   `go build ./...`, `go test -json ./host/evidence -count=1` + the 37-test named manifest,
   `go test ./... -count=1`, and `go test ./... -count=1 -race -timeout 8m`.
3. `./design_docs/verification/w-race-gate-blindspot/run.sh` (timeout 15 min)
4. `./scripts/bench_worldd.sh --smoke`
5. `./scripts/bench_worldd.sh --check-claims`
6. the recorder-refusal arm: `./scripts/bench_worldd.sh --record-pair --variant . --control .` must
   **fail** and its stderr must contain `probe FAILED: sysctl -n hw.ncpu`
7. `bash scripts/test_gate1_range_check.sh`
8. `bash scripts/test_queue_census.sh`
9. `bash scripts/queue_census.sh --doc design_docs/world-mission.md --control-closed 1 --control-open 79`
10. `bash scripts/test_gate0_self_notices.sh`
11. `bash scripts/check_no_personal_email.sh`
12. `bash scripts/test_check_no_personal_email.sh`

Local coverage: `verify_ail.sh`, `go vet ./...`, `go test ./... -count=1`, the race leg on the two
touched packages, and `check_no_personal_email.sh`. Per the standing rig rule the executor does
**not** run `verify_go.sh` locally. None of the other CI steps reads `host/workbench` or the
workbench handler; the whole-repo `go build`/`go test -race` inside `verify_go.sh` is covered for
the touched packages by the local race leg. **CI green on the PR head AND on the merge commit is
the binding gate. The local sweep does not replace it.**

---

## §2 Execution protocol (binding, all milestones)

### 2.1 Environment (every shell call)

```bash
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter185 && export AILANG_BIN=$HOME/.pinned-ailang/ailang && <command>
```

The shell is zsh and does not keep `cwd` between calls. Run scripts with `bash`. Use
`/usr/bin/grep` (the default `grep` is ugrep). There is no `timeout(1)`.

### 2.2 Per-milestone order

1. Make the milestone's test edits exactly as given (§3.2 / §4.2).
2. `gofmt -l host/daemon/workbench_test.go host/workbench/render_test.go` must print nothing (the
   code below is already gofmt-clean; if it prints a file, run `gofmt -w` on **that file only**).
3. Run the milestone exit gate (G-common plus the milestone's own). It must be green.
4. **Commit the milestone** (`git add host/daemon/workbench_test.go host/workbench/render_test.go && git commit`).
   The drill restores each production file by `cp` from a backup it took just before mutating it,
   and its final `STATUS` line must be `''`. That check only means something on a committed tree.
5. Run the milestone's drill rows (§2.3) and append each result line to the Verification Log.
6. Confirm `git status --porcelain` is empty.

**G-common (every milestone exit):**

```bash
go vet ./... ; echo vet_rc=$?                                                        # AC1: vet_rc=0
go test ./... -count=1 2>&1 | /usr/bin/grep -E '^(ok|FAIL|---)' | sort | uniq -c | sort -rn | head   # AC2: 20 ok, 0 FAIL
gofmt -l host/daemon/workbench_test.go host/workbench/render_test.go                 # prints nothing
go test -race ./host/workbench ./host/daemon -count=1                                # both ok
```

`go build ./...` is **not** a compile fence for `_test.go` files; `go vet ./...` is. If AC2 reds only
on `TestHandlerTimeoutKillsTheWholeProcessGroup` or `TestCLIRealSubprocessEpisode` (known load
flakers, not in this sprint's packages), re-run that package once with the machine unloaded before
calling it red. Record both runs.

### 2.3 The drill runner (banked outside the repo; do not re-derive)

The runner is `~/.ailang/state/mission-world-iter185-evidence/planner/mutdrive_planner.py`
(103 lines, sha256 `a4c00adeac68af6298407c5f69e1fbe393b42c1ff06e1503e41f85358fcfa4c5`). It is the
controller's `r2/mutdrive_r2.py` with two changes only: `WT` and `BK` come from the environment, and
each row's named killer(s) from doc §7 are embedded (`KILLERS`) and checked. Copy it; do not edit
the banked copy:

```bash
E=~/.ailang/state/mission-world-iter185-evidence; X=$E/exec; mkdir -p $X/bak
cp $E/planner/mutdrive_planner.py $X/drill.py && shasum -a 256 $X/drill.py   # must print a4c00ade…fcfa4c5
# usage: WT=<worktree> BK=$X/bak python3 $X/drill.py [ID ...]      (no IDs = all 57 rows)
```

What it does per row, which is what AC6 requires: exact single-occurrence string replace (any other
count prints `NOT-LANDED count=N` — a finding, never a pass); `go build ./...` (a failure prints
`compiles=n BUILD-FAIL` and **is never a kill**); `go test ./host/workbench ./host/daemon ./host/boundary -count=1`;
the red set; `killer_in_reds=y|n` (all named killers must be in the red set; `-` for Q17); restore
by `cp` from the backup. The last line is `STATUS: '<git status --porcelain>'`.

A row **passes** only as `compiles=y rc=1 KILLED killer_in_reds=y`. Q17 passes only as
`compiles=y rc=0 SURVIVED`. Output line shape (measured):

```
H5 compiles=y rc=1 KILLED killer_in_reds=y killer=TestSupportedWorkbenchQueryTruthTable/object+payload+world reds=15 TestSupportedWorkbenchQueryTruthTable ...
Q17 compiles=y rc=0 SURVIVED killer_in_reds=- killer=- reds=0
```

Wall time: three full 57-row drills run in parallel in separate worktrees finished together in about 4 minutes (last backup-dir write 20:03:54, all three `.done` markers 20:08:05).

---

## §3 M1 — the grammar (`host/daemon/workbench_test.go`, +37)

### 3.1 What M1 kills and why it is first

The truth table and the three HTTP witnesses are M1's only content. Their killers exist at the end
of M1, so every grammar row is drilled here. The rows whose killers are M2 tests (H3, H4, W5, W8,
V2, V4, V6, V7, V9, V10) **survive** the M1 state (measured, §1.2 m1 drill) and belong to M2.

### 3.2 Code

In `TestWorkbenchRefusalBranches`'s table, insert these three rows **directly after** the
`unsupported-combination` row (`:157`) and before `malformed-payload`:

```go
		{"unsupported-pair-world-from", "/workbench?world=" + world + "&from=0", http.StatusBadRequest, "BadRequest", "unsupported workbench parameter combination", nil},
		{"unsupported-pair-world-payload", "/workbench?world=" + world + "&payload=1", http.StatusBadRequest, "BadRequest", "unsupported workbench parameter combination", nil},
		{"unsupported-triple-world-object-payload", "/workbench?world=" + world + "&object=" + object + "&payload=1", http.StatusBadRequest, "BadRequest", "unsupported workbench parameter combination", nil},
```

Insert this new top-level test **between** the closing `}` of `TestWorkbenchRefusalBranches`
(`:216`) and `func commitWorkbenchPayload` (`:218`), separated by one blank line on each side:

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

Keep every assertion string byte-for-byte (doc §3): §7's red sets were measured against them.
Subtest names are the sorted key set joined by `+` because `keys` is in alphabetical order; the §7
named killers (for example `/object+payload+world`) depend on that order.

### 3.3 Size check

`git diff --shortstat origin/dev -- host/` → `1 file changed, 37 insertions(+)` (measured on the
M1 half).

### 3.4 M1 exit gate

G-common (§2.2), then:

```bash
go test ./host/daemon -count=1 -v -run 'TestSupportedWorkbenchQueryTruthTable|TestWorkbenchRefusalBranches' > /tmp/m1.txt 2>&1; echo rc=$?
/usr/bin/grep -c -- '--- PASS' /tmp/m1.txt                                            # 51
/usr/bin/grep -c -- '--- PASS: TestSupportedWorkbenchQueryTruthTable/' /tmp/m1.txt    # 32
/usr/bin/grep -c -- '--- PASS: TestWorkbenchRefusalBranches/' /tmp/m1.txt             # 17
/usr/bin/grep -c -- '--- FAIL' /tmp/m1.txt                                            # 0
/usr/bin/grep -E -- '--- PASS: TestWorkbenchRefusalBranches/unsupported-(pair|triple)' /tmp/m1.txt | wc -l   # 3
git diff --quiet origin/dev -- host/workbench/render.go host/daemon/workbench.go; echo prod_rc=$?   # 0 (AC5)
```

(`/tmp` is fine for a transient grep input; copy anything you want to keep into
`~/.ailang/state/mission-world-iter185-evidence/exec/`.) Measured: rc=0, 51, 32, 17, 0, 3.
51 = 1 truth-table parent + 32 subtests + 1 refusal parent + 17 refusal subtests (14 existing + 3 new).

### 3.5 M1 mutation drill (commit M1 first)

```bash
WT=/Users/voightkampff/dev/sunholo-data/.wt-world-iter185 BK=$X/bak python3 $X/drill.py H5 H6 H7 Q17 Q18 Q19 Q20 Q23 Q26 Q27 Q30
```

| ID | Mutant (doc §7) | Named killer | Expected (measured on the m1 drill) |
|---|---|---|---|
| H5 | `if len(query) != 2 {` → `if false && len(query) != 2 {` | `TestSupportedWorkbenchQueryTruthTable/object+payload+world` | KILLED, reds=15, killer_in_reds=y |
| H6 | from/entry `&&` → `\|\|` | `…TruthTable/from+world` | KILLED, reds=9 |
| H7 | object/payload `&&` → `\|\|` | `…TruthTable/payload+world` | KILLED, reds=9 |
| Q18 | `!= 2 {⏎return false` → `return true` | `…TruthTable/entry+object+world` | KILLED, reds=19 |
| Q19 | from/entry pair → `if query["entry"] != nil {` | `…TruthTable/entry+world` | KILLED, reds=4 (truth table only) |
| Q20 | from/entry pair → `if query["from"] != nil {` | `…TruthTable/from+world` | KILLED, reds=6 |
| Q23 | from/entry pair → `if true {` | `…TruthTable/object+world` | KILLED, reds=12 |
| Q26 | object/payload return → `return query["payload"] != nil` | `…TruthTable/payload+world` | KILLED, reds=6 |
| Q27 | object/payload return → `return query["object"] != nil` | `…TruthTable/object+world` | KILLED, reds=4 (truth table only) |
| Q30 | object/payload return → `return true` | `…TruthTable/object+world` | KILLED, reds=12 |
| Q17 | `if len(query) != 2 {` → `if len(query) > 2 {` | — (equivalent, doc §2e) | **SURVIVED**, rc=0 (the control) |

Pass: 10 × `compiles=y rc=1 KILLED killer_in_reds=y`, Q17 `compiles=y rc=0 SURVIVED`, `STATUS: ''`.

---

## §4 M2 — href and grade line (`host/workbench/render_test.go`, +53)

### 4.1 What M2 kills

H3, H4, W5, W8, V2, V4, V6, V7, V9, V10. Each survived the M1 state (measured) and each is killed
only by a test this milestone adds. After M2 the full 57-row drill runs (§4.5).

### 4.2 Code

Append to `t.Run("pass", …)` inside `TestGradeViewRequiresTestVerdict`, after the existing
`if view.Label != GradeTESTED || …` check and before the subtest's closing `})` (`view` is the
variable the existing body already binds at `:88`):

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

Insert these two new top-level tests directly after the closing `}` of
`TestGradeViewRequiresTestVerdict` and before `func TestRenderEscapesAllObjectText`, one blank line
between functions:

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

The `—` in the unavailable string is U+2014 (EM DASH), byte-identical to `render.go:154`, and `✓`
is U+2713. Copy from this file; do not retype. The negative loop is load-bearing for V10 and must
not be dropped to save lines (§4.5).

### 4.3 Size check

`git diff --shortstat origin/dev -- host/` → `2 files changed, 90 insertions(+)`;
`git diff --numstat origin/dev -- host/` → `37 0 host/daemon/workbench_test.go`,
`53 0 host/workbench/render_test.go` (measured).

### 4.4 M2 exit gate = the full AC set (§5)

G-common, `./scripts/verify_ail.sh`, and every AC in §5.

### 4.5 M2 mutation drill (commit M2 first)

Step 1 — the M2 rows:

```bash
WT=/Users/voightkampff/dev/sunholo-data/.wt-world-iter185 BK=$X/bak python3 $X/drill.py H3 H4 W5 W8 V2 V4 V6 V7 V9 V10
```

| ID | Mutant (doc §7) | Named killer | Expected (measured on the full drill) |
|---|---|---|---|
| H3 | `aria-label="test verdict PASS"` → `aria-label="x"` | `TestGradeViewRequiresTestVerdict/pass` | KILLED, reds=2 (`…Verdict`, `…/pass`) |
| H4 | `if query != "" && query[0] == '?' {` → `if query != "" {` | `TestWorkbenchHrefAppendsOnlyQueryStrings` | KILLED, reds=1 |
| W5 | the href guard → `if true {` | `TestWorkbenchHrefAppendsOnlyQueryStrings` | KILLED, reds=1 |
| W8 | fallback `return "/workbench"` → `return "/workbench" + query` | `TestWorkbenchHrefAppendsOnlyQueryStrings` | KILLED, reds=1 |
| V2 | `{{if .Grade.HasVerdict}}` → `{{if true}}` | `TestRenderGradeWithoutVerdictClaim/no-verdict` | KILLED, reds=2 |
| V4 | `{{if eq .Grade.Verdict "FAIL"}}` → `{{if true}}` | `TestGradeViewRequiresTestVerdict/pass` | KILLED, reds=2 |
| V6 | `class="verdict-pass"` → `class="verdict-fail"` | `TestGradeViewRequiresTestVerdict/pass` | KILLED, reds=2 |
| V7 | `✓ verdict: {{.Grade.Verdict}}` → `✓ verdict: ` | `TestGradeViewRequiresTestVerdict/pass` | KILLED, reds=2 |
| V9 | `{{if .Grade.Available}}` → `{{if true}}` | `TestRenderGradeWithoutVerdictClaim/unavailable` | KILLED, reds=2 |
| V10 | append an unconditional PASS span after `{{.Grade.Unavailable}}</p>{{end}}` | `TestRenderGradeWithoutVerdictClaim/no-verdict` **and** `/unavailable` | KILLED, reds=3 |

Step 2 — the full 57-row drill on the committed M2 tree (AC6 bullets 1–2):

```bash
WT=/Users/voightkampff/dev/sunholo-data/.wt-world-iter185 BK=$X/bak python3 $X/drill.py > $X/drill_final.txt 2>&1
/usr/bin/grep -c ' KILLED killer_in_reds=y ' $X/drill_final.txt     # 56
/usr/bin/grep -c ' SURVIVED ' $X/drill_final.txt                     # 1
/usr/bin/grep ' SURVIVED ' $X/drill_final.txt | cut -d' ' -f1        # Q17
/usr/bin/grep -cE 'compiles=n|NOT-LANDED|killer_in_reds=n' $X/drill_final.txt   # 0
tail -1 $X/drill_final.txt                                           # STATUS: ''
```

Step 3 — the 20 survivors against `origin/dev`'s tests (AC6 bullet 3). Use a detached sibling
worktree, never `/tmp`, and remove it afterwards:

```bash
R=/Users/voightkampff/dev/sunholo-data/ailang-world; B=/Users/voightkampff/dev/sunholo-data/.exec-base-wt-iter185
git -C $R fetch origin dev && git -C $R worktree add --detach $B origin/dev
WT=$B BK=$X/bak python3 $X/drill.py H3 H4 H5 H6 H7 Q18 Q19 Q20 Q23 Q26 Q27 Q30 W5 W8 V2 V4 V6 V7 V9 V10 > $X/drill_origin_dev.txt 2>&1
/usr/bin/grep -c 'compiles=y rc=0 SURVIVED' $X/drill_origin_dev.txt   # 20
tail -1 $X/drill_origin_dev.txt                                       # STATUS: ''
git -C $R worktree remove --force $B
```

Measured on the planner's pristine-tests drill: all 20 SURVIVED (plus Q17), the other 36 KILLED
with their doc §7 *(existing)* killer in the red set. If `origin/dev` has moved past `fd99840` and a
row no longer survives there, stop and report it: the row's pin then comes from someone else's
change, and the controller decides.

---

## §5 Acceptance criteria → command → expected (the controller checklist)

| AC | Command | Expected (planner-measured on the prototype) |
|---|---|---|
| AC1 | `go vet ./... ; echo rc=$?` | rc=0 |
| AC2 | `go test ./... -count=1` | rc=0; 20 `ok`, 0 `FAIL` |
| AC3 | `./scripts/verify_ail.sh ; echo rc=$?` · `git diff --name-only origin/dev -- '*.ail' \| wc -l` | rc=0; `✓ verify gate PASSED: 11 required identities verified, 40 named tests pass` (and `world package gate PASSED: 9/9`) · `0` |
| AC4 | `go test ./host/workbench ./host/daemon -count=1 -v -run 'TestSupportedWorkbenchQueryTruthTable\|TestWorkbenchRefusalBranches\|TestGradeViewRequiresTestVerdict\|TestRenderGradeWithoutVerdictClaim\|TestWorkbenchHrefAppendsOnlyQueryStrings'` | rc=0; **59** `--- PASS`; 0 `--- FAIL`; `grep -c -- '--- PASS: TestSupportedWorkbenchQueryTruthTable/'` → **32**; PASS lines present for `TestWorkbenchRefusalBranches/unsupported-pair-world-from`, `/unsupported-pair-world-payload`, `/unsupported-triple-world-object-payload`, `TestGradeViewRequiresTestVerdict/pass`, `TestRenderGradeWithoutVerdictClaim/no-verdict`, `/unavailable`, `TestWorkbenchHrefAppendsOnlyQueryStrings` |
| AC5 | `git diff --quiet origin/dev -- host/workbench/render.go host/daemon/workbench.go; echo rc=$?` · `git diff --name-only origin/dev -- host/` | rc=0 · exactly `host/daemon/workbench_test.go`, `host/workbench/render_test.go` |
| AC6 | §4.5 steps 1–3 (and §3.5 at M1) | 56 KILLED with named killer in reds; Q17 SURVIVED; 0 `compiles=n` / NOT-LANDED; the 20 survive on `origin/dev`'s tests; `STATUS: ''` after each run |
| AC7 | `gofmt -l host/daemon/workbench_test.go host/workbench/render_test.go` · `git diff --shortstat origin/dev -- host/` | empty · `2 files changed, 90 insertions(+)` (≤ 120, 0 deletions) |

Also run: `go test -race ./host/workbench ./host/daemon -count=1` (both `ok`) and
`bash scripts/check_no_personal_email.sh` (✓). Do **not** reformat the three pre-existing
gofmt-dirty files.

Load-bearing claims and their mutations (coding-standards S6): every new assertion is discharged by
a named §7 row, never by a count of its own text. The truth table ← H5–H7, Q18–Q20, Q23, Q26, Q27,
Q30 (subtest-named); the HTTP witnesses ← H5/H6/H7 (one per witness, §1.2); the `/pass` render arm
← H3, V4, V6, V7; the `no-verdict` subtest ← V2 and V10; the `unavailable` subtest ← V9 and V10;
the negative verdict-claim loop ← V10 (it survives without the loop: the r1 prototype, doc V21);
the href test ← H4, W5, W8. The `seen` self-check and the AC4 PASS counts are instrument-health
controls, not load-bearing claims.

---

## §6 Estimate and risks

- **Size:** 90 test lines, 0 production lines, 2 files (measured). **Time:** ~0.25 day including
  both drills (about 4 min wall for a full 57-row pass on this rig).
- **Risk 1 — character drift.** The em dash (U+2014) and check mark (U+2713) in the expected strings
  must be byte-identical to `render.go:154`. A retyped hyphen makes `/unavailable` red on the
  pristine tree, which AC2 catches at once. Mitigation: copy the code blocks.
- **Risk 2 — drilling an uncommitted milestone.** The drill's `cp` restore returns the production
  file to its pre-mutation bytes, but `STATUS: ''` is the only proof the tree is clean, and it only
  reads `''` on a committed milestone. Commit first (§2.2 step 4). Never restore with
  `git checkout --`.
- **Risk 3 — load flakes** in `TestHandlerTimeoutKillsTheWholeProcessGroup` /
  `TestCLIRealSubprocessEpisode` (other packages). Re-run once unloaded; record both.
- **Risk 4 — `origin/dev` moves** before §4.5 step 3. The step is measured against the moving
  `origin/dev` on purpose; a changed result is reported, not absorbed (§4.5).
- **CI-only steps** (§1.3) are not reproduced locally. The binding gate is CI green on the PR head
  and on the merge commit.
- **Out of scope** (doc §9, unchanged): R-a (production never supplies a grade reason); H1/H2's
  cross-package-only protection; the three gofmt-dirty base files; the rest of `handleWorkbench` /
  `pageHTML`.
