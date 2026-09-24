# w-workbench-timeline-seam — sprint plan (rows 35 + 38, 3 milestones)

- Sprint: `w-workbench-timeline-seam` · queue rows **35** and **38** (clause 5), plus the unfiled
  `Page.Selected` defect (doc §1 F3) · design doc
  [`w-workbench-timeline-seam.md`](w-workbench-timeline-seam.md) (quorum-closed, revision 2).
- Base: `3d77c61647ab8177f91ff37cd1b1e2b20f02c4d6` (`git rev-parse HEAD`, measured) · branch
  `sprint/w-workbench-timeline-seam` · worktree `/Users/voightkampff/dev/sunholo-data/.wt-world-iter184`.
- `origin/dev` = `82e36306d32f…` (the doc's measurement base). `git diff --stat 82e3630 3d77c61 -- host cmd`
  is **empty** (measured), so every `host/` anchor the doc took at `82e3630` is also an anchor at base.
  The three commits between them are the design doc alone.
- Planned 2026-09-24 by mission-control iteration 184 (sprint-planner). **The whole design was
  prototyped by the planner on a throwaway branch in this worktree.** All 7 AC4 tests passed, and
  all 26 §7 mutation rows were drilled with the runner in §2.3. The prototype was then deleted
  (`git branch -D scratch-iter184-proto`; `git diff --quiet 3d77c61 -- host/` → rc=0 afterwards).
  Every code block below is the prototype's code, copied from the file after `gofmt -w` and
  shown to compile and pass. It is not a paraphrase of the doc.
- Files changed by the sprint (doc §5, confirmed by the prototype's `git diff --stat`):
  `host/workbench/render.go`, `host/workbench/render_test.go`, `host/daemon/workbench.go`,
  `host/daemon/workbench_test.go`. **Nothing else.** No `.ail` file changes.

---

## §0 Where this plan departs from the design doc (measured; read first)

The prototype turned up five places where the doc cannot be followed exactly as written. Each
item below gives the measurement and the prescribed resolution. The executor follows the
resolution and records it in the PR body. It is not a judgement call left to the executor.

| # | Doc says | Measured | Resolution in this plan |
|---|---|---|---|
| P1 | §7 H9: "delete the `if !ok { …; continue }` block". §7 preamble: "Each one compiles." | Deleting the block leaves `ok` unused. `go vet` → `host/daemon/workbench.go:133:6: declared and not used: ok`; **compiled=n**. This is the same defect class as r2's H11. | H9 becomes `if !ok {` → `if false && !ok {`. The mutant compiles, the unstored target renders as an `Available` link, and `TestWorkbenchSelectedEntry/unstored-edge-unavailable` reds (**measured kill**). |
| P2 | §7 H3: "replace the check with `if len(page.Timeline.Entries) == limit {`" | Replacing `if ok {` leaves the probe's `ok` unused. `declared and not used: ok` at `:309`; **compiled=n**. | H3 becomes `if ok {` → `if _ = ok; len(page.Timeline.Entries) == limit {`. That is exactly the old predicate with `ok` kept in use. It compiles, and `/exactly-limit-no-next` reds (**measured kill**). |
| P3 | §8: M2 drills H9–H14 | H14's killer is `TestWorkbenchTimelinePaging/last-page-has-prev-no-next`, an **M3** test. At M2 the mutant compiled and the red set was **empty** (measured: with one entry and `from=0`, `pageHref(0,0)` equals `pageHref(from,0)`). | H14 is drilled in **M3**. M2 drills H9–H13. |
| P4 | §6 AC9: `git diff --shortstat origin/dev -- host/` insertions+deletions ≤ **320**. §8 says ~270 LOC. | The prototype, which carries exactly AC4's assertions and controls, measured **517 insertions + 31 deletions = 548**. Per file (`--numstat`): workbench.go 71/5, workbench_test.go 306/10, render.go 19/16, render_test.go 121/0. The non-test share is **111** lines. | AC9 as written is **expected RED**, and the executor must **NOT** delete assertions or controls to meet it. The executor records both numbers (total and `git diff --shortstat origin/dev -- host/ ':!*_test.go'`). The controller decides the disposition, for example re-baselining AC9 as non-test ≤ 150 and total ≤ 600. This is flagged for the controller, not silently amended (§6 AC9 row). |
| P5 | §5: "It does not move row 34's line anchors either: every change in `workbench.go` lands at `:101` and below" | The `host/hashref` import that §3.3 itself requires is inserted at `:11`, which moves row 34's hunks **+1**. Measured on the prototype: `if len(query) != 2` goes `:69` → `:70`, `from`+`entry` goes `:72` → `:73`, `object`+`payload` goes `:75` → `:76`. The function bodies stay byte-identical (AC5 rc=0 on all three). | Nothing is changed in code (the import is required). The executor states the +1 shift in the PR body so that row 34's next drill re-anchors by content. |

Doc anchor drift, also corrected (all other cited anchors re-verified exact at base; see §1.3):
§5 cites `TestWorkbenchRefusalBranches` as `workbench_test.go:133-234`. At base it is **`:133-215`**
(`:217` is `commitWorkbenchPayload`). §5 cites `TestWorkbenchTimelineBound` as `:297-330`. At base
it is **`:297-332`**. F6 cites `validateRef` at `store.go:437-446`. At base `:437` is
`validateRefText` and `validateRef` is at **`:444`**. The claim itself holds (`hashref.Parse` at `:438`).

One strengthening beyond the doc's text, not a change of intent. `emitted-links-resolve` in §5.2
also asserts that the number of `<a ` openings in each region equals the number of
classification-regex matches. The doc's "none is dropped" otherwise depends on the regex, and a
link that the regex failed to match would escape classification silently.

---

## §1 Baseline (pristine tree, measured by the planner before any edit)

### 1.1 Commands and results (base `3d77c61`, 2026-09-24, load average 2.18 at start)

All commands run from the worktree with `export AILANG_BIN=$HOME/.pinned-ailang/ailang`
(`$AILANG_BIN --version` → `AILANG v0.41.0`, commit `24ee108`; `go version` → `go1.26.6 darwin/arm64`).

| Gate | Command | Result at base |
|---|---|---|
| vet | `go vet ./... ; echo rc=$?` | **rc=0** |
| test | `go test ./... -count=1 2>&1 \| tail -30` | **rc=0**; **20 `ok`**, **0 `FAIL`**. Control: `go list ./... \| wc -l` = **20** packages, all 20 with test files. No flake occurred, so no re-run was needed (the known flakers `TestHandlerTimeoutKillsTheWholeProcessGroup` and `TestCLIRealSubprocessEpisode` passed in the first run). |
| .ail | `./scripts/verify_ail.sh ; echo rc=$?` | **rc=0**. `✓ world package gate PASSED: 9/9 steps performed non-zero work`; `✓ verify gate PASSED: 11 required identities verified, 40 named tests pass` |
| gofmt | `gofmt -l host cmd` | **NOT EMPTY AT BASE: 3 files.** `host/daemon/session_middleware_test.go` (comment alignment), `host/store/schema_version_test.go` (one mis-indented line `:345`) and `host/store/store.go` (trailing blank line). **This gate is red at base, so it measures the repo, not the change.** None of the three is a file this sprint touches. The sprint's gofmt gate is therefore scoped: `gofmt -l host/workbench host/daemon/workbench.go host/daemon/workbench_test.go` must print nothing. It printed nothing at base (measured), and the pre-existing `session_middleware_test.go` in `host/daemon` is excluded by naming files. Do **not** reformat the three base files, because they are out of scope. |
| race (touched pkgs) | `go test -race ./host/workbench ./host/daemon -count=1` | both `ok` (workbench 1.39 s, daemon 4.86 s) |
| email | `bash scripts/check_no_personal_email.sh` | rc=0, `✓ … no personal addresses in the loop-written surface` |

### 1.2 CI's command list vs this local list

`.github/workflows/ci.yml` has two jobs. Job `ailang-verify` runs `./scripts/verify_ail.sh`, which is
covered locally. Job `go-verify` runs `./scripts/verify_go.sh` followed by nine more steps. The
local sweep above is a subset. **CI steps that are not in the local list:**

- `verify_go.sh` internals: the tracked-binary hygiene gate, the World decision-ledger check
  (`scripts/mission_decisions.sh --check`), the race-detector known-positive control, `go build ./...`,
  the `host/evidence` 37-test named manifest, and **`go test ./... -count=1 -race -timeout 8m`**.
  Locally this plan runs the race leg on the two touched packages only (the table above). Per the
  standing rig rule, the executor does **not** run `verify_go.sh` locally.
- `./design_docs/verification/w-race-gate-blindspot/run.sh` (compiler reproducer, platform-gated).
- `./scripts/bench_worldd.sh --smoke`, `--check-claims`, and the `--record-pair` off-rig refusal arm.
- `scripts/test_gate1_range_check.sh`, `scripts/test_queue_census.sh`,
  `scripts/queue_census.sh --doc design_docs/world-mission.md --control-closed 1 --control-open 79`,
  `scripts/test_gate0_self_notices.sh`, `scripts/check_no_personal_email.sh`,
  `scripts/test_check_no_personal_email.sh`.

None of these reads `host/workbench` or the workbench handler, apart from the whole-repo `go
build`/`go test -race`, which the local race leg covers for the touched packages. **CI green on the PR
head AND on the merge commit is the binding gate. The local sweep does not replace it.**

### 1.3 Anchor re-verification (doc cites → base `3d77c61`)

Each anchor below was checked with `git show 3d77c61:<file> | sed -n '<range>p'`, and each was exact
unless marked. `render.go`: `:63-65` (three `EdgeView` fields), `:67` (`Edges []EdgeView`), `:74`
(`Truncated`), `:92` (`Selected`), `:101` (`workbenchHref` body), `:114` (`const pageHTML`),
`:140-141` (prev/next), `:143-147` (row article), `:151` (`<h2>Inspector</h2>`), `:162` (inline
edge), `:168-171` (`pageTemplate`). `workbench.go`: `:7` (`"math"`), `:33-39`
(`acceptedWorkbenchKeys`), `:62-76` (`supportedWorkbenchQuery`; `:69/:72/:75`), `:101-112`
(`entryView`; `:107-109` the three edge lines), `:136-147` (`from` parse), `:150` (overflow
bound), `:201` (world `StateRoot`), `:213` (`GetObject`), `:239-251` (selection; `:240`
`GetLogEntry`; `:249-250`), `:253` (`TimelineView{From, Limit}`), `:254-264` (loop; `:255`
`GetLogEntry`; `:263` append), `:265` (`Truncated` write). `workbench_test.go`: `:75-131`,
`:87-93` (stale comment), `:156` (`unsupported-combination`), **`:133-215`** (doc said `-234`),
**`:297-332`** (doc said `-330`). `handlers_test.go:79-99` (`testCommit`). `read_deadline_test.go`:
`:502-505` (workbench route), `:693-715` (`failingStore`, 5 methods). `store.go`: `:115-116`,
`:128` (`hashref.HashRef` fields), `:482-483` (`GetObject` SELECT includes `payload`), **`:444`**
`validateRef` (doc said `:437`). `daemon.go`: `:53` (`hashref` import), `:108`
(`maxCommitBytes = 8388608`), `:331-337` (`readStore`, 5 methods).

---

## §2 Execution protocol (binding, all milestones)

### 2.1 Environment (every shell call)

```bash
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter184 && export AILANG_BIN=$HOME/.pinned-ailang/ailang && <command>
```

The shell does not keep `cwd` between calls. Every command below assumes this prefix.

### 2.2 Per-milestone order

1. Make the milestone's code and test edits exactly as given.
2. `gofmt -w host/workbench/render.go host/workbench/render_test.go host/daemon/workbench.go host/daemon/workbench_test.go`
   (gofmt realigns the trailing comments on the `next` lines in M3; the drill patterns below
   already allow for this).
3. Run the milestone exit gate (G-common plus the milestone's own). It must be green.
4. **Commit the milestone** (`git add host/ && git commit`). The mutation drill restores with
   `git checkout -- <file>`, which is only correct against a committed milestone. Drilling on an
   uncommitted tree destroys the milestone.
5. Run the milestone's drill rows (§2.3) one at a time and append each result line to the Verification Log.
6. Confirm `git status --short -- host/` is empty and `git diff --quiet -- host/` rc=0.

**G-common (every milestone exit):**

```bash
go vet ./... ; echo vet_rc=$?                                              # AC1: rc=0
go test ./... -count=1 2>&1 | grep -E '^(ok|FAIL|---)' | sort | uniq -c | sort -rn | head   # AC2: 20 ok, 0 FAIL
gofmt -l host/workbench host/daemon/workbench.go host/daemon/workbench_test.go   # must print nothing
go test -race ./host/workbench ./host/daemon -count=1                      # both ok
```

If AC2 reds only on `TestHandlerTimeoutKillsTheWholeProcessGroup` or `TestCLIRealSubprocessEpisode`,
re-run that package once with the machine unloaded before calling it red. Record both runs.

### 2.3 The drill runner (create once, outside the repo)

Evidence goes under `~/.ailang/state/world-iter184/` because `/tmp` is wiped. This is a
scratch input, not a repo file:

```bash
mkdir -p ~/.ailang/state/world-iter184 && cat > ~/.ailang/state/world-iter184/drill.sh <<'EOF'
#!/bin/bash
# usage: drill.sh ID FILE 'PERL_MATCH' 'PERL_REPLACEMENT' KILLER
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter184 || exit 9
export AILANG_BIN=$HOME/.pinned-ailang/ailang
L=~/.ailang/state/world-iter184
id=$1 f=$2 m=$3 r=$4 killer=$5
n=$(M="$m" perl -0ne 'my $c = () = /$ENV{M}/g; print $c' "$f")
if [ "$n" != 1 ]; then echo "$id | count=$n (want 1) | NOT APPLIED"; exit 1; fi
M="$m" R="$r" perl -0pi -e 's/$ENV{M}/$ENV{R}/' "$f"
landed=n; git diff --quiet -- "$f" || landed=y
compiled=n; go vet ./host/workbench ./host/daemon >"$L/vet.$id" 2>&1 && compiled=y
out=$(go test ./host/workbench ./host/daemon -count=1 2>&1); printf '%s\n' "$out" >"$L/test.$id"
reds=$(printf '%s\n' "$out" | grep -oE -- '--- FAIL: [^ ]+' | sed 's/--- FAIL: //' | sort -u | tr '\n' ' ')
killed=n
[ "$compiled" = y ] && printf '%s\n' "$out" | grep -qE -- "--- FAIL: ${killer}( |$)" && killed=y
git checkout -- "$f"; restored=n; git diff --quiet -- host/ && restored=y
echo "$id | count=1 | landed=$landed | compiled=$compiled | killer=$killer killed=$killed | restored=$restored | reds: $reds"
EOF
chmod +x ~/.ailang/state/world-iter184/drill.sh
```

The runner performs each of the steps AC8 requires. It asserts an occurrence count of exactly 1
before editing, and a row with any other count is **NOT APPLIED** and is a finding. It records
`landed` and `compiled` separately, and **a build failure is never a kill** (`killed` needs
`compiled=y`). It records the red set, restores with `git checkout -- <file>`, and checks
`git diff --quiet -- host/`. A row passes only with `landed=y compiled=y killed=y restored=y`.
Each regex is a perl pattern applied with `-0` to the whole file. The replacement is inserted
literally. Below, `D=$HOME/.ailang/state/world-iter184/drill.sh`.

---

## §3 M1 — the seam (render.go view model + template; handler's mechanical edits)

**Files:** `host/workbench/render.go`, `host/workbench/render_test.go`, `host/daemon/workbench.go`.

### 3.1 `host/workbench/render.go`

1. **`EntryView` (`:58-68`)**: delete `:63-65` (`TransitionFn`/`Interpreter`/`TransitionRef EdgeView`)
   and add `SelectHref string` after `WrittenBy`. Result:

   ```go
   type EntryView struct {
   	EntryIndex     int64
   	EntryHash      string
   	PrevEntryHash  string
   	SemanticsEpoch int64
   	WrittenBy      string
   	SelectHref     string
   	Edges          []EdgeView
   }
   ```

2. **`TimelineView` (`:70-77`)**: delete `:74` `Truncated bool`. gofmt re-aligns the rest:

   ```go
   type TimelineView struct {
   	From     int64
   	Limit    int
   	Entries  []EntryView
   	NextHref string
   	PrevHref string
   }
   ```

3. **Row article (`:143-147`)**: replace the three inline `<dt>` lines and `</dl></article>` so
   that the five lines read:

   ```
   <article><h3>entry {{.EntryIndex}}</h3><dl>
   {{template "entryFields" .}}
   </dl><p><a href="{{workbenchHref .SelectHref}}">select entry {{.EntryIndex}}</a></p></article>
   ```

4. **Selected block**: directly after `:151` `<h2>Inspector</h2>` (and before `{{with .Object}}`), insert:

   ```
   {{with .Selected}}<article aria-label="selected entry"><h3>selected entry {{.EntryIndex}}</h3><dl>
   {{template "entryFields" .}}
   </dl>{{range .Edges}}{{template "edge" .}}{{end}}</article>{{end}}
   ```

5. **Provenance walk (`:162`)** becomes `{{with .Object}}{{range .Edges}}{{template "edge" .}}{{end}}{{end}}`.
6. **`pageTemplate` (`:168-171`)**: replace it with the `partialsHTML` constant and the double
   `Parse`. The `"edge"` body is `:162`'s inner body byte-for-byte, and the `"entryFields"` body is
   `:144-146` byte-for-byte (doc §3.2):

   ```go
   const partialsHTML = `{{define "edge"}}<p>{{.Relation}}: {{if edgeUnavailable .}}<span class="unavailable" role="note">UNAVAILABLE: {{.Missing}}</span>{{else}}<a href="{{workbenchHref .Href}}" class="hash" title="{{.Target}}" aria-label="{{.Target}}">{{.Target}}</a>{{end}}</p>{{end}}
   {{define "entryFields"}}<dt>entry hash</dt><dd><span class="hash" title="{{.EntryHash}}" aria-label="{{.EntryHash}}">{{.EntryHash}}</span></dd>
   <dt>previous entry</dt><dd><span class="hash" title="{{.PrevEntryHash}}" aria-label="{{.PrevEntryHash}}">{{.PrevEntryHash}}</span></dd>
   <dt>semantics epoch</dt><dd>{{.SemanticsEpoch}}</dd><dt>written by</dt><dd>{{.WrittenBy}}</dd>{{end}}`

   var pageTemplate = template.Must(template.Must(template.New("workbench").Funcs(template.FuncMap{
   	"edgeUnavailable": edgeUnavailable,
   	"workbenchHref":   workbenchHref,
   }).Parse(pageHTML)).Parse(partialsHTML))
   ```

   `workbenchHref` (`:100-105`) and `edgeUnavailable` stay **byte-identical**, because row 34 owns `:101`.

### 3.2 `host/daemon/workbench.go` (mechanical, so that it compiles)

- `entryView` (`:101-112`): delete `:107-109`. The other five fields stay, and it still has no store access.
- Delete `:265` `page.Timeline.Truncated = len(page.Timeline.Entries) == limit`.

At M1 the rows' `SelectHref` is empty, so the select link renders `href="/workbench"`. That is a
transitional state, and M2 fills it.

### 3.3 New tests (append to `host/workbench/render_test.go`; add `"reflect"` to its imports)

These are the prototype's code, verbatim and gofmt-clean. They carry AC4's four render-package
assertions, including the `Selected: nil` control arm and the self-checking exemption arm:

```go
func renderPage(t *testing.T, p Page) string {
	t.Helper()
	var body bytes.Buffer
	if err := Render(&body, p); err != nil {
		t.Fatal(err)
	}
	return body.String()
}

// selectedArticle returns the substring from the selected-entry marker to the
// first </article> after it, or "" when the marker is absent.
func selectedArticle(body string) string {
	const marker = `<article aria-label="selected entry">`
	start := strings.Index(body, marker)
	if start < 0 {
		return ""
	}
	end := strings.Index(body[start:], "</article>")
	if end < 0 {
		return body[start:]
	}
	return body[start : start+end]
}

func TestRenderSelectedEntry(t *testing.T) {
	p := Page{Title: "selected", Selected: &EntryView{
		EntryIndex: 7,
		EntryHash:  "SEL-HASH-7",
		Edges: []EdgeView{
			{Relation: "transitionRef", Available: true, Target: "abc", Href: "?object=abc"},
			{Relation: "interpreter", Target: "def", Missing: "object def is not stored"},
		},
	}}
	article := selectedArticle(renderPage(t, p))
	if article == "" {
		t.Fatal(`rendered page has no <article aria-label="selected entry">`)
	}
	for _, want := range []string{
		`<h3>selected entry 7</h3>`,
		`title="SEL-HASH-7"`,
		`transitionRef: <a href="/workbench?object=abc"`,
		`interpreter: <span class="unavailable" role="note">UNAVAILABLE: object def is not stored</span>`,
	} {
		if !strings.Contains(article, want) {
			t.Errorf("selected-entry article missing %q: %q", want, article)
		}
	}
	// CONTROL: the article is Selected's and nothing else's.
	if body := renderPage(t, Page{Title: "unselected"}); strings.Contains(body, `aria-label="selected entry"`) {
		t.Errorf("Selected=nil still rendered a selected-entry article: %q", body)
	}
}

func TestRenderTimelineRowSelectLink(t *testing.T) {
	body := renderPage(t, Page{Title: "rows", Timeline: TimelineView{Entries: []EntryView{{EntryIndex: 3, SelectHref: "?from=0&entry=3"}}}})
	if want := `<a href="/workbench?from=0&amp;entry=3">select entry 3</a>`; !strings.Contains(body, want) {
		t.Errorf("row select link %q missing: %q", want, body)
	}
}

func TestRenderTimelinePagingLinks(t *testing.T) {
	t.Run("prev", func(t *testing.T) {
		body := renderPage(t, Page{Title: "prev", Timeline: TimelineView{PrevHref: "?from=0&entry=0"}})
		if want := `<a href="/workbench?from=0&amp;entry=0">previous</a>`; !strings.Contains(body, want) {
			t.Errorf("previous link %q missing: %q", want, body)
		}
	})
	t.Run("next", func(t *testing.T) {
		body := renderPage(t, Page{Title: "next", Timeline: TimelineView{NextHref: "?from=100&entry=100"}})
		if want := `<a href="/workbench?from=100&amp;entry=100">next</a>`; !strings.Contains(body, want) {
			t.Errorf("next link %q missing: %q", want, body)
		}
	})
	t.Run("neither", func(t *testing.T) {
		body := renderPage(t, Page{Title: "neither"})
		for _, unwanted := range []string{">previous</a>", ">next</a>"} {
			if strings.Contains(body, unwanted) {
				t.Errorf("zero TimelineView rendered %q", unwanted)
			}
		}
	})
}

// unrenderedExempt names the view-model fields that are deliberately not
// rendered (row candidate R-a). The self-checking arm below fails if one of
// them starts rendering, so an exemption cannot outlive its reason.
var unrenderedExempt = map[string]bool{
	"TimelineView.From":  true,
	"TimelineView.Limit": true,
}

// TestWorkbenchViewFieldsAllRender is the I3 class ratchet: a view-model field
// that the handler writes and the template never reads fails here. The check is
// lexical and name-based by design; specific mutations name specific killers.
func TestWorkbenchViewFieldsAllRender(t *testing.T) {
	text := pageHTML + partialsHTML
	checked := 0
	for _, typ := range []reflect.Type{reflect.TypeOf(EntryView{}), reflect.TypeOf(TimelineView{}), reflect.TypeOf(Page{})} {
		for i := 0; i < typ.NumField(); i++ {
			name := typ.Field(i).Name
			key := typ.Name() + "." + name
			rendered := strings.Contains(text, "."+name)
			if unrenderedExempt[key] {
				if rendered {
					t.Errorf("%s is exempt as unrendered but the template now renders .%s; remove the exemption", key, name)
				}
				continue
			}
			checked++
			if !rendered {
				t.Errorf("%s is never rendered by the workbench template (no .%s action)", key, name)
			}
		}
	}
	// CONTROL: the census must have walked real fields, not an empty set.
	if checked < 10 {
		t.Fatalf("field census checked %d fields, want >= 10", checked)
	}
}
```

`TestWorkbenchViewFieldsAllRender` walks 16 non-exempt fields (EntryView 7 + TimelineView 3 +
Page 6). The `checked < 10` control guards against a census that walks nothing. The existing
render tests (`TestRenderUnavailableProvenanceEdge`, `TestRenderEscapesAllObjectText`,
`TestRenderEmitsOnlyLocalLinks`) are **not modified** and must pass as they are.

### 3.4 M1 exit gate

G-common, plus:

```bash
go test ./host/workbench -count=1 -v -run 'TestRenderSelectedEntry|TestRenderTimelineRowSelectLink|TestRenderTimelinePagingLinks|TestWorkbenchViewFieldsAllRender' 2>&1 | grep -E '^--- '   # 4 × --- PASS, 0 FAIL
grep -rnw --include='*.go' Truncated host/workbench host/daemon | wc -l          # 0   (control: grep -rn --include='*.go' PayloadTruncated host | wc -l → 3)
grep -cE '^\s+(TransitionFn|Interpreter|TransitionRef)\s+EdgeView' host/workbench/render.go   # 0 (control: grep -cE '^\s+Edges\s+\[\]EdgeView' … → 2)
```

### 3.5 M1 mutation drill (R1–R11; all measured on the prototype: count=1, compiled=y, killed=y, restored=y)

```bash
F=host/workbench/render.go
$D R1  $F '(?s)\{\{with \.Selected\}\}.*?</article>\{\{end\}\}\n' '' TestRenderSelectedEntry
$D R2  $F '</dl>\{\{range \.Edges\}\}\{\{template "edge" \.\}\}\{\{end\}\}</article>' '</dl></article>' TestRenderSelectedEntry
$D R3  $F '(?<=<h3>selected entry \{\{\.EntryIndex\}\}</h3><dl>\n)\{\{template "entryFields" \.\}\}\n' '' TestRenderSelectedEntry
$D R4  $F '\{\{if edgeUnavailable \.\}\}' '{{if false}}' TestRenderUnavailableProvenanceEdge
$D R5  $F '(?<=\{\{with \.Object\}\}\{\{range \.Edges\}\})\{\{template "edge" \.\}\}' '' TestRenderUnavailableProvenanceEdge
$D R6  $F '<p><a href="\{\{workbenchHref \.SelectHref\}\}">select entry \{\{\.EntryIndex\}\}</a></p>' '' TestRenderTimelineRowSelectLink
$D R7  $F '(?<=<article><h3>entry \{\{\.EntryIndex\}\}</h3><dl>\n)\{\{template "entryFields" \.\}\}\n' '' TestWorkbenchRendersSeededWorldAndTimeline
$D R8  $F '\{\{if \.Timeline\.NextHref\}\}[^\n]*\n' '' TestRenderTimelinePagingLinks/next
$D R9  $F '\{\{if \.Timeline\.PrevHref\}\}[^\n]*\n' '' TestRenderTimelinePagingLinks/prev
$D R10 $F '(?<=\tSelectHref     string)(?=\n)' '; Extra string' TestWorkbenchViewFieldsAllRender
$D R11 $F '<h2>Timeline</h2>' '<h2>Timeline {{.Timeline.From}}</h2>' TestWorkbenchViewFieldsAllRender
```

R10 adds the field on the same line (`SelectHref     string; Extra string`), which is valid Go
field-list syntax. This keeps the replacement free of newlines. (`$(printf …)` would strip a
trailing newline and join two field lines. The planner measured that pitfall, and this form avoids it.)

Red sets that the planner measured (for comparison, since killer membership is the criterion):
R1 {SelectedEntry, ViewFieldsAllRender}; R2, R3 {SelectedEntry}; R4 {SelectedEntry,
UnavailableProvenanceEdge}; R5 {UnavailableProvenanceEdge}; R6 {RowSelectLink, ViewFieldsAllRender};
R7 {RendersSeededWorldAndTimeline}; R8 {PagingLinks/next, ViewFieldsAllRender}; R9
{PagingLinks/prev, ViewFieldsAllRender}; R10, R11 {ViewFieldsAllRender}.

**M1 ~LOC:** prototype render.go 19+/16−, render_test.go 121+, workbench.go 0+/4−. About 160 lines, against the doc's ~110.

---

## §4 M2 — selection edges

**Files:** `host/daemon/workbench.go`, `host/daemon/workbench_test.go`.

### 4.1 `host/daemon/workbench.go`

1. **Import**: add `"github.com/sunholo-data/ailang-world/host/hashref"` above the `store` import
   (base `:11`). This shifts row 34's anchors by +1 (§0 P5). Record it in the PR body.
2. **After `entryView` (base `:112`)**, add `pageHref` and `entryEdges`:

   ```go
   // pageHref is the only builder of a from/entry workbench query: every paging
   // and row-selection link uses the existing from+entry grammar state.
   func pageHref(from, entry int64) string {
   	return "?from=" + strconv.FormatInt(from, 10) + "&entry=" + strconv.FormatInt(entry, 10)
   }

   // entryEdges checks each of the entry's three edge targets once. A stored target
   // is a link; an unstored one is UNAVAILABLE with its ref visible; a store error
   // is returned so the caller answers 5xx rather than "not stored".
   func (d *Daemon) entryEdges(ctx context.Context, entry store.LogEntry) ([]workbench.EdgeView, error) {
   	refs := []struct {
   		relation string
   		ref      hashref.HashRef
   	}{
   		{"transitionFn", entry.Header.TransitionFn},
   		{"interpreter", entry.Header.Interpreter},
   		{"transitionRef", entry.TransitionRef},
   	}
   	edges := make([]workbench.EdgeView, 0, len(refs))
   	for _, item := range refs {
   		target := item.ref.String()
   		_, ok, err := d.reads.GetObject(ctx, item.ref)
   		if err != nil {
   			return nil, err
   		}
   		if !ok {
   			edges = append(edges, workbench.EdgeView{Relation: item.relation, Target: target, Missing: "object " + target + " is not stored"})
   			continue
   		}
   		edges = append(edges, workbench.EdgeView{Relation: item.relation, Available: true, Target: target, Href: "?object=" + target})
   	}
   	return edges, nil
   }
   ```

3. **Selected block (base `:249-250`)**:

   ```go
   		selected := entryView(entry)
   		edges, err := d.entryEdges(ctx, entry)
   		if err != nil {
   			d.writeWorkbenchStoreError(w, r, ctx, err)
   			return
   		}
   		selected.Edges = edges
   		page.Selected = &selected
   ```

4. **Timeline loop (base `:263`)**:

   ```go
   		view := entryView(entry)
   		view.SelectHref = pageHref(from, view.EntryIndex)
   		page.Timeline.Entries = append(page.Timeline.Entries, view)
   ```

### 4.2 `host/daemon/workbench_test.go`

1. **Rewrite the stale comment at `:87-93`** (its "renders no TransitionRef action at all" is false
   after this sprint). The assertions below the comment do not change:

   ```go
   	// The observable is the log entry hash, NOT the committed object hash. On
   	// `/workbench` (no selection) the timeline rows render no edges, so the
   	// object hash reaches the page only through the selected entry's
   	// transitionRef edge -- pinned by TestWorkbenchSelectedEntry, not here. The
   	// entry hash is written by the timeline loop and by nothing else on this
   	// page, which is what makes it a pin on the mechanism rather than on a
   	// sibling channel.
   ```

2. **Append** the region helper (M3 reuses it), `objectFailingStore`, and `TestWorkbenchSelectedEntry`.
   `errSentinelInternal` already exists in the package (used by `failingStore`,
   `read_deadline_test.go:697-715`). `objectFailingStore` embeds the `readStore` interface, so
   only `GetObject` is overridden (doc V14: `failingStore` overrides all 5 methods and cannot
   isolate the edge check):

```go
// workbenchRegion returns the substring from start to the first end after it.
func workbenchRegion(body, start, end string) (string, bool) {
	i := strings.Index(body, start)
	if i < 0 {
		return "", false
	}
	j := strings.Index(body[i:], end)
	if j < 0 {
		return "", false
	}
	return body[i : i+j], true
}

const (
	selectedEntryStart = `<article aria-label="selected entry">`
	timelineStart      = `<section aria-label="timeline">`
)

// objectFailingStore fails GetObject only, so the entry-edge check is the one
// read that errors while every log and world read still reaches the real store.
type objectFailingStore struct {
	readStore
}

func (objectFailingStore) GetObject(context.Context, hashref.HashRef) (store.Object, bool, error) {
	return store.Object{}, false, errSentinelInternal
}

func TestWorkbenchSelectedEntry(t *testing.T) {
	d := newHandlerDaemon(t)
	genesis := seedGenesisEmbedded(t, d, "workbench-selected")
	commit := testCommit(genesis, 0, "workbench-selected")
	if err := d.store.Commit(commit); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	header := commit.Entry.Header
	transitionRef := commit.Entry.TransitionRef.String()
	transitionFn := header.TransitionFn.String()
	interpreter := header.Interpreter.String()
	// POSITIVE CONTROLS: the fixture must hold one stored and two unstored targets,
	// or the stored/unstored arms below assert nothing.
	for _, control := range []struct {
		name string
		ref  hashref.HashRef
		want bool
	}{
		{"TransitionRef", commit.Entry.TransitionRef, true},
		{"TransitionFn", header.TransitionFn, false},
		{"Interpreter", header.Interpreter, false},
	} {
		if _, ok, err := d.store.GetObject(context.Background(), control.ref); err != nil || ok != control.want {
			t.Fatalf("control GetObject(%s): ok=%v err=%v, want ok=%v", control.name, ok, err, control.want)
		}
	}

	selectedBody := func(t *testing.T) string {
		t.Helper()
		rec := requestRecorder(t, d, http.MethodGet, "/workbench?from=0&entry=0", nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body)
		}
		article, ok := workbenchRegion(rec.Body.String(), selectedEntryStart, "</article>")
		if !ok {
			t.Fatalf("no selected-entry article in %s", rec.Body)
		}
		return article
	}

	t.Run("stored-edge-links", func(t *testing.T) {
		article := selectedBody(t)
		if want := `transitionRef: <a href="/workbench?object=` + transitionRef + `"`; !strings.Contains(article, want) {
			t.Errorf("selected entry missing stored edge link %q: %s", want, article)
		}
	})

	t.Run("unstored-edge-unavailable", func(t *testing.T) {
		article := selectedBody(t)
		for _, edge := range []struct{ relation, ref string }{{"transitionFn", transitionFn}, {"interpreter", interpreter}} {
			want := edge.relation + `: <span class="unavailable" role="note">UNAVAILABLE: object ` + edge.ref + ` is not stored</span>`
			if !strings.Contains(article, want) {
				t.Errorf("selected entry missing %q: %s", want, article)
			}
		}
		if unwanted := `href="/workbench?object=` + transitionFn + `"`; strings.Contains(article, unwanted) {
			t.Errorf("unstored transitionFn rendered as a link %q", unwanted)
		}
	})

	t.Run("row-select-link", func(t *testing.T) {
		rec := requestRecorder(t, d, http.MethodGet, "/workbench", nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body)
		}
		body := rec.Body.String()
		if want := `href="/workbench?from=0&amp;entry=0">select entry 0</a>`; !strings.Contains(body, want) {
			t.Errorf("timeline row missing select link %q", want)
		}
		if strings.Contains(body, `aria-label="selected entry"`) {
			t.Error("/workbench with no entry= rendered a selected-entry article")
		}
	})

	t.Run("object-store-error", func(t *testing.T) {
		oldReads, oldErrLog := d.reads, d.errLog
		defer func() { d.reads, d.errLog = oldReads, oldErrLog }()
		d.reads = objectFailingStore{readStore: d.store}
		d.errLog = &bytes.Buffer{}
		// CONTROL: the same store serves the unselected page, so only the edge
		// check can be what fails below.
		if rec := requestRecorder(t, d, http.MethodGet, "/workbench", nil); rec.Code != http.StatusOK {
			t.Fatalf("control: /workbench status = %d, want 200; body=%s", rec.Code, rec.Body)
		}
		rec := requestRecorder(t, d, http.MethodGet, "/workbench?from=0&entry=0", nil)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500; body=%s", rec.Code, rec.Body)
		}
		if !strings.Contains(rec.Body.String(), ">Internal<") {
			t.Errorf("body does not contain >Internal<: %s", rec.Body)
		}
	})
}
```

### 4.3 M2 exit gate

G-common, plus:

```bash
go test ./host/daemon -count=1 -v -run 'TestWorkbenchSelectedEntry' 2>&1 | grep -E '^\s*--- '   # 1 + 4 subtests PASS
go test ./host/daemon -count=1 -run 'TestWorkbenchRefusalBranches|TestReadCtx|TestWorkbench' 2>&1 | tail -1   # ok (the accepted-keys arm now exercises entryEdges on /workbench?from=0&entry=0)
grep -c 'context.Background' host/daemon/workbench.go       # 0   (control: grep -c 'readCtx' host/daemon/workbench.go → 1)
```

### 4.4 M2 mutation drill (H9–H13; H14 moves to M3 per §0 P3)

```bash
F=host/daemon/workbench.go
$D H9  $F 'if !ok \{(?=\n\t\t\tedges = append)' 'if false && !ok {' TestWorkbenchSelectedEntry/unstored-edge-unavailable     # §0 P1 correction
$D H10 $F 'if err != nil \{(?=\n\t\t\treturn nil, err\n)' 'if false && err != nil {' TestWorkbenchSelectedEntry/object-store-error
$D H11 $F 'selected\.Edges = edges(?=\n)' 'selected.Edges = edges[:0]' TestWorkbenchSelectedEntry/stored-edge-links
$D H12 $F '(?<=edges, err := d\.entryEdges\(ctx, entry\)\n\t\t)if err != nil \{' 'if false && err != nil {' TestWorkbenchSelectedEntry/object-store-error
$D H13 $F '\t\tview\.SelectHref = pageHref\(from, view\.EntryIndex\)\n' '' TestWorkbenchSelectedEntry/row-select-link
```

H11 is the doc's r2 form (`edges[:0]` keeps `edges` in use so the mutant compiles). The lookahead
anchors it to the one assignment line.

Planner-measured red sets: H9 {SelectedEntry/unstored-edge-unavailable}; H10, H12
{SelectedEntry/object-store-error}; H11 {…/stored-edge-links, …/unstored-edge-unavailable};
H13 {…/row-select-link}. **Literal doc H9 was measured compiled=n**, and it must not be recorded as a kill.

**M2 ~LOC:** workbench.go +45 production; workbench_test.go about +130 (comment 7/7, helper,
store, test). About 180 lines, against the doc's ~80.

---

## §5 M3 — paging

**Files:** `host/daemon/workbench.go`, `host/daemon/workbench_test.go`.

### 5.1 `host/daemon/workbench.go`: the two probes, where `Truncated` was (base `:265`, after the loop's closing `}`)

```go
	// Probe, don't infer: a paging link is emitted only after this request has
	// read the entry it selects, so every emitted paging link resolves.
	next := from + int64(limit)             // cannot overflow: the from bound above refused from > MaxInt64-limit
	if next <= math.MaxInt64-int64(limit) { // the link's own from must pass that bound on the next request
		_, ok, err := d.reads.GetLogEntry(ctx, next)
		if err != nil {
			d.writeWorkbenchStoreError(w, r, ctx, err)
			return
		}
		if ok {
			page.Timeline.NextHref = pageHref(next, next)
		}
	}
	if from > 0 {
		prev := from - int64(limit)
		if prev < 0 {
			prev = 0
		}
		_, ok, err := d.reads.GetLogEntry(ctx, prev)
		if err != nil {
			d.writeWorkbenchStoreError(w, r, ctx, err)
			return
		}
		if ok {
			page.Timeline.PrevHref = pageHref(prev, prev)
		}
	}
```

`"math"` is already imported (`:7`) and already used (`:150`). No new import is needed.

### 5.2 `host/daemon/workbench_test.go`

1. **Extract `seedWorkbenchLog`** from `TestWorkbenchTimelineBound` (base `:297-313`). The test
   keeps its positive control and all of its assertions **unchanged**:

```go
// seedWorkbenchLog commits n chained testCommit entries (indices 0..n-1) on a
// fresh genesis. Every entry's TransitionRef object is stored; its TransitionFn
// and Interpreter targets are not (testCommit).
func seedWorkbenchLog(t *testing.T, d *Daemon, n int) {
	t.Helper()
	world := seedGenesisEmbedded(t, d, "workbench-timeline-bound")
	started := time.Now()
	for index := int64(0); index < int64(n); index++ {
		commit := testCommit(world, index, "workbench-timeline-bound")
		if err := d.store.Commit(commit); err != nil {
			t.Fatalf("Commit(%d): %v", index, err)
		}
		world = commit.NextWorld
		if time.Since(started) > 30*time.Second {
			t.Fatalf("seeding exceeded 30 seconds after %d commits", index+1)
		}
	}
}

func TestWorkbenchTimelineBound(t *testing.T) {
	d := newHandlerDaemon(t)
	seedWorkbenchLog(t, d, workbench.WorkbenchPageLimit+5)
	if _, ok, err := d.store.GetLogEntry(context.Background(), 104); err != nil || !ok {
	// ... unchanged from here (positive control body, the /workbench request, and the three assertions)
```

2. **Add `"regexp"`** to the imports and **append** `TestWorkbenchTimelinePaging`, `denseLogStore`
   and `TestWorkbenchNextLinkOverflowGuard`:

```go
var workbenchLinkPattern = regexp.MustCompile(`<a href="([^"]*)"[^>]*>([^<]*)</a>`)
var selectEntryText = regexp.MustCompile(`^select entry [0-9]+$`)

func TestWorkbenchTimelinePaging(t *testing.T) {
	d := newHandlerDaemon(t)
	seedWorkbenchLog(t, d, workbench.WorkbenchPageLimit+5)
	if _, ok, err := d.store.GetLogEntry(context.Background(), 104); err != nil || !ok {
		t.Fatalf("positive control GetLogEntry(104): ok=%v err=%v", ok, err)
	}
	get := func(t *testing.T, target string) string {
		t.Helper()
		rec := requestRecorder(t, d, http.MethodGet, target, nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status = %d, want 200; body=%s", target, rec.Code, rec.Body)
		}
		return rec.Body.String()
	}

	t.Run("head-page-has-next", func(t *testing.T) {
		body := get(t, "/workbench")
		if want := `href="/workbench?from=100&amp;entry=100">next</a>`; !strings.Contains(body, want) {
			t.Errorf("head page missing next link %q", want)
		}
		if strings.Contains(body, ">previous</a>") {
			t.Error("head page (from=0) rendered a previous link")
		}
	})

	t.Run("last-page-has-prev-no-next", func(t *testing.T) {
		body := get(t, "/workbench?from=100&entry=100")
		if got := strings.Count(body, "<h3>entry "); got != 5 {
			t.Errorf("last page entry count = %d, want 5", got)
		}
		for _, want := range []string{
			`href="/workbench?from=0&amp;entry=0">previous</a>`,
			`href="/workbench?from=100&amp;entry=104">select entry 104</a>`,
		} {
			if !strings.Contains(body, want) {
				t.Errorf("last page missing %q", want)
			}
		}
		if strings.Contains(body, ">next</a>") {
			t.Error("last page rendered a next link past the log end")
		}
	})

	t.Run("exactly-limit-no-next", func(t *testing.T) {
		body := get(t, "/workbench?from=5&entry=5")
		// CONTROL: the page is exactly full, so the old len==limit predicate holds here.
		if got := strings.Count(body, "<h3>entry "); got != workbench.WorkbenchPageLimit {
			t.Fatalf("control: entry count = %d, want %d", got, workbench.WorkbenchPageLimit)
		}
		if strings.Contains(body, ">next</a>") {
			t.Error("exactly-full last page rendered a next link to an empty page (F4)")
		}
		if want := `href="/workbench?from=0&amp;entry=0">previous</a>`; !strings.Contains(body, want) {
			t.Errorf("from=5 page missing clamped previous link %q", want)
		}
	})

	t.Run("emitted-links-resolve", func(t *testing.T) {
		counts := map[string]int{}
		for _, page := range []struct {
			target      string
			hasSelected bool
		}{
			{"/workbench", false},
			{"/workbench?from=100&entry=100", true},
			{"/workbench?from=5&entry=5", true},
		} {
			body := get(t, page.target)
			timeline, ok := workbenchRegion(body, timelineStart, "</section>")
			if !ok {
				t.Fatalf("%s: no timeline region", page.target)
			}
			regions := []struct{ name, text string }{{"timeline", timeline}}
			selected, hasSelected := workbenchRegion(body, selectedEntryStart, "</article>")
			if hasSelected != page.hasSelected {
				t.Fatalf("%s: selected-entry region present=%v, want %v", page.target, hasSelected, page.hasSelected)
			}
			if hasSelected {
				regions = append(regions, struct{ name, text string }{"selected", selected})
			}
			for _, region := range regions {
				matches := workbenchLinkPattern.FindAllStringSubmatch(region.text, -1)
				// Nothing escapes classification: every anchor must be a pattern match.
				if anchors := strings.Count(region.text, "<a "); anchors != len(matches) {
					t.Fatalf("%s %s region: %d anchors but %d pattern matches", page.target, region.name, anchors, len(matches))
				}
				for _, match := range matches {
					href, text := match[1], match[2]
					category := ""
					switch {
					case region.name == "timeline" && (text == "previous" || text == "next"):
						category = "paging"
					case region.name == "timeline" && selectEntryText.MatchString(text):
						category = "select"
					case region.name == "selected" && strings.HasPrefix(href, "/workbench?object="):
						category = "stored-edge"
					default:
						t.Errorf("%s %s region: unclassified link href=%q text=%q", page.target, region.name, href, text)
						continue
					}
					counts[category]++
					target := strings.ReplaceAll(href, "&amp;", "&")
					if rec := requestRecorder(t, d, http.MethodGet, target, nil); rec.Code != http.StatusOK {
						t.Errorf("%s: %s link %q returned %d, want 200", page.target, category, target, rec.Code)
					}
				}
			}
		}
		// CONTROL: every category must be exercised, or a dead extractor passes vacuously.
		for _, category := range []string{"paging", "select", "stored-edge"} {
			if counts[category] == 0 {
				t.Errorf("category %s matched 0 links across the three pages", category)
			}
		}
	})
}

// denseLogStore answers every non-negative index with a copy of stored entry 0,
// so the from bound can be exercised without writing 2^63 entries.
type denseLogStore struct {
	readStore
}

func (s denseLogStore) GetLogEntry(ctx context.Context, index int64) (store.LogEntry, bool, error) {
	if index < 0 {
		return store.LogEntry{}, false, nil
	}
	entry, ok, err := s.readStore.GetLogEntry(ctx, 0)
	if err != nil || !ok {
		return entry, ok, err
	}
	entry.Header.EntryIndex = index
	return entry, true, nil
}

func TestWorkbenchNextLinkOverflowGuard(t *testing.T) {
	d := newHandlerDaemon(t)
	genesis := seedGenesisEmbedded(t, d, "workbench-overflow")
	if err := d.store.Commit(testCommit(genesis, 0, "workbench-overflow")); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	d.reads = denseLogStore{readStore: d.store}

	// MaxInt64-100 is the largest from the handler's from bound accepts; its next
	// link would carry from=MaxInt64, which that bound refuses, so none may render.
	rec := requestRecorder(t, d, http.MethodGet, "/workbench?from=9223372036854775707&entry=9223372036854775707", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body)
	}
	if strings.Contains(rec.Body.String(), ">next</a>") {
		t.Error("from=MaxInt64-100 rendered a next link whose own from overflows")
	}
	// CONTROL: one page earlier, the dense store does produce a next link.
	rec = requestRecorder(t, d, http.MethodGet, "/workbench?from=9223372036854775607&entry=9223372036854775607", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("control status = %d, want 200; body=%s", rec.Code, rec.Body)
	}
	if want := `href="/workbench?from=9223372036854775707&amp;entry=9223372036854775707">next</a>`; !strings.Contains(rec.Body.String(), want) {
		t.Errorf("control: missing next link %q", want)
	}
}
```

### 5.3 M3 exit gate = the full AC set (§6)

### 5.4 M3 mutation drill (H1–H8, H14, H15)

```bash
F=host/daemon/workbench.go
$D H1  $F '\t\t\tpage\.Timeline\.NextHref = pageHref\(next, next\)\n' '' TestWorkbenchTimelinePaging/head-page-has-next
$D H2  $F '(?<=\t\t\})\n\t\tif ok \{(?=\n\t\t\tpage\.Timeline\.NextHref)' "$(printf '\n\t\tif !ok {')" TestWorkbenchTimelinePaging/exactly-limit-no-next
$D H3  $F '(?<=\t\t\})\n\t\tif ok \{(?=\n\t\t\tpage\.Timeline\.NextHref)' "$(printf '\n\t\tif _ = ok; len(page.Timeline.Entries) == limit {')" TestWorkbenchTimelinePaging/exactly-limit-no-next   # §0 P2 correction
$D H4  $F 'if next <= math\.MaxInt64-int64\(limit\) \{' 'if true {' TestWorkbenchNextLinkOverflowGuard
$D H5  $F '\t\t\tpage\.Timeline\.PrevHref = pageHref\(prev, prev\)\n' '' TestWorkbenchTimelinePaging/last-page-has-prev-no-next
$D H6  $F '\t\tif prev < 0 \{\n\t\t\tprev = 0\n\t\t\}\n' '' TestWorkbenchTimelinePaging/exactly-limit-no-next
$D H7  $F 'if from > 0 \{' 'if true {' TestWorkbenchTimelinePaging/head-page-has-next
$D H8  $F '"\?from=" \+ strconv\.FormatInt\(from, 10\) \+ "&entry=" \+ strconv\.FormatInt\(entry, 10\)' '"?from=" + strconv.FormatInt(from, 10)' TestWorkbenchTimelinePaging/emitted-links-resolve
$D H14 $F 'pageHref\(from, view\.EntryIndex\)' 'pageHref(view.EntryIndex, view.EntryIndex)' TestWorkbenchTimelinePaging/last-page-has-prev-no-next   # moved from M2, §0 P3
$D H15 $F 'pageHref\(next, next\)' 'pageHref(next, from)' TestWorkbenchTimelinePaging/head-page-has-next
```

Leading `\n` in a replacement is kept by `$(printf …)`, because only trailing newlines are stripped.
The `if from > 0 {` count is 1 (the M2/M3 code has no other occurrence; the runner asserts it).
Before H7, the executor confirms `grep -c 'if from > 0 {' host/daemon/workbench.go` → 1.

Planner-measured red sets (each includes its named killer): H1 {…/head-page-has-next,
OverflowGuard}; H2 {head, last-page, exactly-limit, emitted-links, OverflowGuard}; H3 (corrected)
{exactly-limit, emitted-links}; H4 {OverflowGuard}; H5 {last-page, exactly-limit}; H6
{exactly-limit}; H7 {head}; H8 {emitted-links, head, last-page, exactly-limit, OverflowGuard,
SelectedEntry/row-select-link}; H14 {last-page}; H15 {head, OverflowGuard}. **Literal doc H3 was
measured compiled=n.**

**M3 ~LOC:** workbench.go +26 production; workbench_test.go about +185. About 210 lines, against the doc's ~90.

---

### 5.5 Verbatim re-run of every drill line in this plan

After this plan was written, the planner took the three drill blocks above and the §2.3 runner
out of this file's text. It ran them **as written** against the re-applied final prototype
(scratch branch, since deleted), and **26/26 rows returned `count=1 landed=y compiled=y killed=y
restored=y`**. Red sets on the final tree are supersets of the per-milestone sets listed above,
because later milestones' tests also red (for example, R1 additionally reds
`TestWorkbenchSelectedEntry/*` and `emitted-links-resolve`). The criterion is membership of the
killer, not equality of the set.

## §6 Controller checklist (AC1–AC9 → command)

Every line runs with the §2.1 prefix, on the final sprint head. The planner's prototype results are shown in brackets.

| AC | Command | Pass condition [prototype] |
|---|---|---|
| AC1 | `go vet ./... ; echo rc=$?` | rc=0 [rc=0] |
| AC2 | `go test ./... -count=1 2>&1 \| grep -cE '^ok'` and `… \| grep -cE '^(FAIL\|--- FAIL)'` | 20 and 0; rc=0 [20, 0] |
| AC3 | `./scripts/verify_ail.sh ; echo rc=$?` and `git diff --name-only origin/dev -- '*.ail' \| wc -l` | rc=0 and `0` [0] |
| AC4 | `go test ./host/workbench ./host/daemon -count=1 -v -run 'TestRenderSelectedEntry\|TestRenderTimelineRowSelectLink\|TestRenderTimelinePagingLinks\|TestWorkbenchViewFieldsAllRender\|TestWorkbenchSelectedEntry\|TestWorkbenchTimelinePaging\|TestWorkbenchNextLinkOverflowGuard' 2>&1 \| grep -c '^--- PASS'` | 7, and `grep -c -- '--- FAIL'` → 0 [7, 0] |
| AC5 | `for e in '/^func supportedWorkbenchQuery/,/^}/p' '/^var acceptedWorkbenchKeys/,/^}/p'; do diff <(git show origin/dev:host/daemon/workbench.go \| sed -n "$e") <(sed -n "$e" host/daemon/workbench.go); echo rc=$?; done; e='/^func TestWorkbenchRefusalBranches/,/^}/p'; diff <(git show origin/dev:host/daemon/workbench_test.go \| sed -n "$e") <(sed -n "$e" host/daemon/workbench_test.go); echo rc=$?` | three × rc=0. Control: each extraction is non-empty (15 / 7 / 83 lines) [rc=0 ×3; 15/7/83] |
| AC6 | `grep -rnw --include='*.go' Truncated host/workbench host/daemon \| wc -l` · `grep -rn --include='*.go' PayloadTruncated host \| wc -l` · `grep -cE '^\s+(TransitionFn\|Interpreter\|TransitionRef)\s+EdgeView' host/workbench/render.go` · `grep -cE '^\s+Edges\s+\[\]EdgeView' host/workbench/render.go` · `grep -cE 'Timeline\.(NextHref\|PrevHref) = ' host/daemon/workbench.go` · `grep -c 'Selected' host/workbench/render.go` | 0 · ≥3 · 0 · 2 · 2 · ≥2 [0·3·0·2·2·2]. These are instrument-health controls. The load-bearing claims are discharged by the drill (S6). |
| AC7 | `grep -c 'context.Background' host/daemon/workbench.go` · `grep -c 'readCtx' host/daemon/workbench.go` · `grep -c '"?from=" + strconv' host/daemon/workbench.go` | 0 · 1 · 1 [0·1·1] |
| AC8 | The §3.5, §4.4 and §5.4 drill lines, all 26 rows (R1–R11, H1–H15) | every row `count=1 landed=y compiled=y killed=y restored=y` [26/26 with the §0 P1/P2 corrections; literal H3 and H9 measured compiled=n] |
| AC9 | `git diff --shortstat origin/dev -- host/` and `git diff --shortstat origin/dev -- host/ ':!*_test.go'` | **RE-BASELINED by the controller (iteration 184, before execution; rule 3h — adjudicated by the planner's measured prototype, 548 total / 111 non-test):** non-test ≤ **150** AND total ≤ **600**. The doc's ≤ 320 was an estimate, not a correctness bound. Record both numbers; tests are never cut to meet it. |

Also: gofmt scoped gate (§1.1) prints nothing; `go test -race ./host/workbench ./host/daemon -count=1`
ok [ok, daemon 12.95 s]; and after merge, **CI on the PR head and on the merge commit** (both jobs) is green.

---

## §7 Estimate and risks

| M | Doc ~LOC | Prototype-measured | Time |
|---|---|---|---|
| M1 | 110 | about 160 (render.go 35, render_test.go 121, workbench.go 4) | 0.3 d |
| M2 | 80 | about 180 | 0.3 d |
| M3 | 90 | about 210 | 0.4 d |
| **Total** | **~270** | **548** (517+/31−; 111 non-test) | **~1 d**, plus about 25 min of drill (26 rows at about 50 s each on this rig) |

Risks, in honest order:

1. **AC9 is red as written (P4).** The size was underestimated by about 2× because the doc's
   own AC4 assertions and controls cost that much test code. If the controller will not
   re-baseline it, the sprint cannot close green without weakening the tests, and weakening the
   tests is forbidden. This is the top risk and needs a decision **before** the executor starts,
   or it must be accepted as a recorded red.
2. **The drill runner's perl quoting.** Replacements that need a newline (R10, H2, H3, H11) are the
   fragile ones. The runner's `landed` and `compiled` columns make a botched edit visible (a
   botched edit shows as `compiled=n` or count≠1, and never as a false kill). §4.4 gives the
   newline-free H11 form.
3. **Test-time growth in `host/daemon`.** `TestWorkbenchTimelinePaging` seeds 105 commits and
   fetches about 110 links (measured 0.45 s, versus 0.8 s for the whole AC4 daemon set). Under
   `-race` the daemon package measured 4.86 s at base and 12.95 s on the prototype. The two runs
   were under different machine load, so the delta is not attributed to the sprint alone. Both are
   far inside CI's 25-minute `timeout-minutes` and the 8-minute `-race -timeout`.
4. **Known flakers** (`TestHandlerTimeoutKillsTheWholeProcessGroup`, `TestCLIRealSubprocessEpisode`)
   can red AC2 under load. Re-run them once unloaded, and record both runs.
5. **Row 34 re-anchoring (P5).** Row 34's hunk line numbers move +1. Its content is unchanged
   (AC5), but a drill that is anchored by line number would miss. This is PR-body-only; no code
   changes.
6. **R-c stays open.** The world pane's `StateRoot` link still 404s on the seeded world (doc V15).
   `emitted-links-resolve` excludes the world pane by region, and this is intended (I1). The
   sprint does **not** fix it (doc §9).
7. **CI-only steps (§1.2)** are not reproduced locally. The binding gate is CI on the PR head and
   on the merge commit.

---

## Execution record (iteration 184)

Executed by the sprint executor on 2026-09-24 in this worktree, branch `sprint/w-workbench-timeline-seam`,
starting from `0a69f1a` (the plan with the controller's AC9 re-baseline). Every command ran with
`export AILANG_BIN=$HOME/.pinned-ailang/ailang` (AILANG v0.41.0). `origin/dev` = `82e36306d32f…` (measured).
Evidence (per-row `vet.*`/`test.*` outputs, drill logs, full test outputs) is under
`~/.ailang/state/world-iter184/`. The drill runner is §2.3's text, verbatim.

### Commits

| Milestone | SHA | Files | numstat vs parent |
|---|---|---|---|
| M1 the seam | `750f8f6` | render.go, render_test.go, daemon/workbench.go | 19/16, 121/0, 0/4 (= prototype) |
| M2 selection edges | `01acc46` | daemon/workbench.go, daemon/workbench_test.go | 45+ / 136+ (173+/8− total) |
| M3 paging | `db94908` | daemon/workbench.go, daemon/workbench_test.go | 204+/3− total |

Final `git diff --numstat origin/dev -- host/` is byte-for-byte the prototype's: workbench.go 71/5,
workbench_test.go 306/10, render.go 19/16, render_test.go 121/0.

### Per-milestone exit gates

| Gate | M1 (`750f8f6`) | M2 (`01acc46`) | M3 (`db94908`) |
|---|---|---|---|
| `go vet ./...` | rc=0 | rc=0 | rc=0 |
| `go test ./... -count=1` | rc=0; 20 `ok`, 0 FAIL | rc=0; 20 `ok`, 0 FAIL | rc=0; 20 `ok`, 0 FAIL |
| scoped `gofmt -l host/workbench host/daemon/workbench.go host/daemon/workbench_test.go` | empty | empty | empty |
| `go test -race ./host/workbench ./host/daemon -count=1` | ok / ok (daemon 4.96 s) | ok / ok (4.96 s) | ok / ok (12.59 s) |
| milestone-specific | 4 × `--- PASS` (§3.4); `Truncated` 0 (control `PayloadTruncated` 3); dead EdgeView fields 0 (control `Edges []EdgeView` 2) | `TestWorkbenchSelectedEntry` + 4 subtests PASS; `-run 'TestWorkbenchRefusalBranches\|TestReadCtx\|TestWorkbench'` ok; `context.Background` 0 (control `readCtx` 1) | full AC set below |

No known flaker (`TestHandlerTimeoutKillsTheWholeProcessGroup`, `TestCLIRealSubprocessEpisode`) went red in
any run, so no re-run was needed.

### Mutation drill (AC8) — 26/26 killed by the named test

Every row: `count=1`, run one at a time, restored with `git checkout -- <file>`, then `git diff --quiet -- host/` rc=0.
`git status --short -- host/` was empty after each block.

| Row | Landed | Compiled | Red set | Named killer red? | Restored |
|---|---|---|---|---|---|
| R1 | y | y | TestRenderSelectedEntry, TestWorkbenchViewFieldsAllRender | y (TestRenderSelectedEntry) | y |
| R2 | y | y | TestRenderSelectedEntry | y | y |
| R3 | y | y | TestRenderSelectedEntry | y | y |
| R4 | y | y | TestRenderSelectedEntry, TestRenderUnavailableProvenanceEdge | y (TestRenderUnavailableProvenanceEdge) | y |
| R5 | y | y | TestRenderUnavailableProvenanceEdge | y | y |
| R6 | y | y | TestRenderTimelineRowSelectLink, TestWorkbenchViewFieldsAllRender | y (TestRenderTimelineRowSelectLink) | y |
| R7 | y | y | TestWorkbenchRendersSeededWorldAndTimeline | y | y |
| R8 | y | y | TestRenderTimelinePagingLinks{,/next}, TestWorkbenchViewFieldsAllRender | y (…/next) | y |
| R9 | y | y | TestRenderTimelinePagingLinks{,/prev}, TestWorkbenchViewFieldsAllRender | y (…/prev) | y |
| R10 | y | y | TestWorkbenchViewFieldsAllRender | y | y |
| R11 | y | y | TestWorkbenchViewFieldsAllRender | y | y |
| H9 (§0 P1 form) | y | y | TestWorkbenchSelectedEntry{,/unstored-edge-unavailable} | y | y |
| H10 | y | y | TestWorkbenchSelectedEntry{,/object-store-error} | y | y |
| H11 | y | y | TestWorkbenchSelectedEntry{,/stored-edge-links,/unstored-edge-unavailable} | y (…/stored-edge-links) | y |
| H12 | y | y | TestWorkbenchSelectedEntry{,/object-store-error} | y | y |
| H13 | y | y | TestWorkbenchSelectedEntry{,/row-select-link} | y | y |
| H1 | y | y | TestWorkbenchTimelinePaging{,/head-page-has-next}, TestWorkbenchNextLinkOverflowGuard | y (…/head-page-has-next) | y |
| H2 | y | y | TestWorkbenchTimelinePaging{,/head-page-has-next,/last-page-has-prev-no-next,/exactly-limit-no-next,/emitted-links-resolve}, TestWorkbenchNextLinkOverflowGuard | y (…/exactly-limit-no-next) | y |
| H3 (§0 P2 form) | y | y | TestWorkbenchTimelinePaging{,/exactly-limit-no-next,/emitted-links-resolve} | y (…/exactly-limit-no-next) | y |
| H4 | y | y | TestWorkbenchNextLinkOverflowGuard | y | y |
| H5 | y | y | TestWorkbenchTimelinePaging{,/last-page-has-prev-no-next,/exactly-limit-no-next} | y (…/last-page-has-prev-no-next) | y |
| H6 | y | y | TestWorkbenchTimelinePaging{,/exactly-limit-no-next} | y | y |
| H7 (pre-check `grep -c 'if from > 0 {'` = 1) | y | y | TestWorkbenchTimelinePaging{,/head-page-has-next} | y | y |
| H8 | y | y | TestWorkbenchTimelinePaging{,/head-page-has-next,/last-page-has-prev-no-next,/exactly-limit-no-next,/emitted-links-resolve}, TestWorkbenchSelectedEntry{,/row-select-link}, TestWorkbenchNextLinkOverflowGuard | y (…/emitted-links-resolve) | y |
| H14 (moved to M3, §0 P3) | y | y | TestWorkbenchTimelinePaging{,/last-page-has-prev-no-next} | y | y |
| H15 | y | y | TestWorkbenchTimelinePaging{,/head-page-has-next}, TestWorkbenchNextLinkOverflowGuard | y (…/head-page-has-next) | y |

No row failed to compile, no row survived, no row had count≠1. Every red set equals the planner's measured set.

### Final gates on `db94908` (load average 2.72 at start)

| Gate | Command | Result |
|---|---|---|
| AC1 | `go vet ./...` | rc=0 |
| AC2 | `go test ./... -count=1` | rc=0; 20 `ok`, 0 `FAIL`/`--- FAIL` (no flake, no re-run) |
| race | `go test -race ./host/workbench ./host/daemon -count=1` | rc=0; workbench 1.29 s, daemon 13.20 s |
| AC3 | `./scripts/verify_ail.sh` · `git diff --name-only origin/dev -- '*.ail' \| wc -l` | rc=0 (`world package gate PASSED: 9/9`, `verify gate PASSED: 11 required identities verified, 40 named tests pass`) · 0 |
| email | `bash scripts/check_no_personal_email.sh` | rc=0 |
| gofmt (scoped) | `gofmt -l host/workbench host/daemon/workbench.go host/daemon/workbench_test.go` | empty |
| AC4 | the 7-name `-run` | rc=0; 7 `--- PASS`, 0 `--- FAIL` (Paging 0.45 s) |
| AC5 | the three `diff` extractions vs `origin/dev` | rc=0 ×3; extraction sizes 15 / 7 / 83 lines (non-empty controls) |
| AC6 | the six greps | 0 · 3 · 0 · 2 · 2 · 2 — PASS |
| AC7 | the three greps | 0 · 1 · 1 — PASS |
| AC8 | the 26 drill rows | 26/26 `count=1 landed=y compiled=y killed=y restored=y` — PASS |
| AC9 | `git diff --shortstat origin/dev -- host/` · `… ':!*_test.go'` | **total 548** (517+/31−) ≤ 600 · **non-test 111** (90+/21−) ≤ 150 — PASS under the controller's re-baseline (the doc's original ≤ 320 would be red, as §0 P4 predicted) |

AC1–AC9: **all PASS**. CI on the PR head and the merge commit remains the binding gate (§1.2); it was not
run here (no push, no PR from this role).

§0 P5 confirmed: `if len(query) != 2` moved `:69` → `:70`, and `origin/dev:69-76` diffs clean (rc=0)
against the sprint's `:70-77`. Row 34's next drill must re-anchor by content (+1 line).

### Deviations from the plan

1. **Commit staging.** §2.2 step 4 says `git add host/`. Per the controller's brief, each milestone staged
   its exact files by name instead. The staged set is identical (only those files were modified).
2. **M1 AC2 tally re-run.** §2.2's G-common AC2 line ends in `| head`, which shows only 10 of the 20 `ok`
   lines, so it cannot show "20 ok" by itself. For M1 the full suite was re-run with explicit
   `grep -cE '^ok'` / `grep -cE '^(FAIL|--- FAIL)'` counts (20 / 0). For M2, M3 and the final gate the
   counted form was used from the start. The instrument, not the code, was changed.
3. **The M3 probes' position.** They were inserted directly after the timeline loop's closing `}`, where
   M1 removed the `Truncated` line. The blank line before `_ = workbench.Render(w, page)` is kept. The
   final numstat equals the prototype's, which is consistent with the same placement.

No code deviation. The §0 P1/P2/P3 corrections (H9, H3, H14 moved to M3) were applied as the plan
prescribes. Nothing else was changed: no `.ail`, no `tools/launchd/*`, no `.claude/`, and the three
un-gofmt'd base files were not touched.

### Controller addendum (iteration 184, after the judge)

The evaluator (`sonnet`, own worktree @ `88855dd`, **PASS 95/100, zero blocking**) found two survivors of its own
choosing: `if err != nil` → `if false && err != nil` on the next-probe and on the prev-probe `GetLogEntry` (M3). The
shipped code was correct; nothing pinned it. The controller added `TestWorkbenchPagingProbeStoreError`
(`probeFailingStore` fails exactly one index; per-subtest control fails an unread index → 200). Measured: with
the test present, each mutant lands (1 occurrence), compiles, and reds **only** `TestWorkbenchPagingProbeStoreError`
in `./host/workbench ./host/daemon`; restored byte-identical by sha256. Test-only delta; no production change.
