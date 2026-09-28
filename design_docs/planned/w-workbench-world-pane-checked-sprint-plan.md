# Sprint plan: w-workbench-world-pane-checked (iteration 205, rows 99 + 100)

**Base:** detached `2904ebc`. **Authority:** [approved design](w-workbench-world-pane-checked.md), especially D1–D5, §§4–7, 9–11. Quorum r1 and r2 blocked 3/3 on completeness only; r3 applied the r2 fixes verbatim under the narrow-refinement carve-out. This plan implements the approved design, including V13's two existing-test adaptations and V20–V21's command/anchor corrections. The controller owns any commits; no git write is part of this plan.

## Gates and measured base

Use this environment before **every** Go or AILANG command, including focused and mutation runs:

```sh
export PATH=/opt/homebrew/bin:$PATH AILANG_BIN=$HOME/.pinned-ailang/ailang
```

| Untouched base gate | Controller measurement at `aa3e36a` outside sandbox | Applicability at `2904ebc` |
|---|---|---|
| `go vet ./...` | rc 0 | V18: `host/` identical. |
| `go test ./... -run '^$'` | rc 0 | Compiles production and `_test.go`; V18 applies. |
| `go test -race -count=1 ./host/daemon ./host/workbench` | rc 0; daemon 17.2 s, workbench 1.4 s | V18 applies; CI repeats race. |
| `./scripts/verify_ail.sh` | rc 0 | Run again at final closure. |

The controller measurements are inherited evidence, not planner reruns. A local socket `bind: operation not permitted` is **UNINFORMATIVE UNDER SANDBOX**, neither pass nor fail. The relevant daemon tests use `requestRecorder`, not a listener. Every landing must pass its named tests, `go vet ./...`, `go test ./... -run '^$'`, and `go test -race -count=1 ./host/daemon ./host/workbench` at that boundary. Final closure also runs `go test -count=1 ./...` and `./scripts/verify_ail.sh`. Acceptance tests observe responses and use Go temporary storage; they write no acceptance artifact outside the worktree or into `$HOME`. If a server is needed, keep it foreground and bounded; never kill a PID read from a file, use `kill -1/0`, or `pkill`.

## Landing boundaries

| Landing | Added production LOC estimate | Compiling, bisectable boundary | Tests introduced or adapted here |
|---|---:|---|---|
| M1 | ~25 | Checked stateRoot read, error path, explicit no-world stop, world identity/revision/log-head markup. Keep `TimelineView.From/Limit` and old ratchet for now. | `TestRenderWorldPane`, `TestWorkbenchWorldPane`, widened `TestWorkbenchTimelinePaging/emitted-links-resolve`; adapt `TestWorkbenchSelectedEntry/object-store-error` and `TestWorkbenchReferenceWalk` **with** the production change. |
| M2 | 0 | Remove the two unused timeline fields/assignments and exemptions; switch to the typed, fail-closed nine-type census. M1 behavior remains green. | Replace `TestWorkbenchViewFieldsAllRender`; add `TestWorkbenchFieldCensusMutations` and its positive controls. All M2 analyzer work is test code. |

The M1 production estimate counts net-new Go/template lines; M2 deletes production lines. Both are far below ~150 added production LOC. Measure actual `git diff --numstat` at each landing. No M2 test or removed field may enter M1. If the M2 test helper proves too large to review in one landing, split M2a (typed analyzer and self-tests, retaining fields/exemptions) and M2b (field removal and zero-exemption ratchet switch), with all three gates and existing tests green at each boundary; the planned two-landings path remains preferred because M2 adds no production code.

## Implementation contracts for exact mutations

Use the §4.1 successful-world branch exactly enough to preserve the named anchors:

```go
stateRoot, err := d.checkedEdge(ctx, "stateRoot", world.StateRoot)
if err != nil {
    d.writeWorkbenchStoreError(w, r, ctx, err)
    return
}
page.World = workbench.WorldView{
    Ref: world.Ref.String(), Revision: world.Revision,
    StateRoot: stateRoot, LogHead: world.LogHead.String(), Available: true,
}
```

Gofmt may put each assignment on its own line; the mutation literals below use those individual field substrings. Initialize `page := workbench.Page{Title: "AILANG World Workbench", World: workbench.WorldView{Unavailable: "no world selected"}}`. Keep GetWorld-not-found as 404. Keep the local `from`/`limit` loops; M2 removes only `page.Timeline = workbench.TimelineView{From: from, Limit: limit}` and the struct fields. A root read error returns before `Render` and uses `writeWorkbenchStoreError`; a missing root is a 200 with metadata and an unavailable edge. There is one additional GetObject for a selected world per request; no timeline-row check.

Render the §4.2 nav exactly: retain the `workbench` home anchor; world ref as `<dt>world</dt><dd><span class="hash" title="{{.World.Ref}}" aria-label="{{.World.Ref}}">{{.World.Ref}}</span></dd>`; revision as `<dt>revision</dt><dd>{{.World.Revision}}</dd>`; log head as `<dt>log head</dt><dd><span class="hash" title="{{.World.LogHead}}" aria-label="{{.World.LogHead}}">{{.World.LogHead}}</span></dd>`; then `{{template "edge" .World.StateRoot}}`. The existing else arm includes `{{.World.Unavailable}}`. The edge partial alone decides whether a stateRoot object anchor exists and labels it with its Target in text/title/aria-label. World ref and log head stay full, copyable, non-anchor hashes. No CSS or route change.

In M1, `TestWorkbenchWorldPane` must exercise absent and stored roots on default, selected-entry, and object pages; distinct sentinels pin ref, root, revision (17 and zero), and log head. Its `explicit-world` case requests an older ref while another head is selected. A targeted wrapper fails only the root for 500/Internal and 503/Timeout, with an untargeted 200 control, sanitized body and no partial nav. A counting wrapper pins one root GetObject on `/workbench`. `TestRenderWorldPane` uses `renderPage(t, Page{...})`, not daemon helpers; extract nav and assert metadata, root identity, missing and no-world stops, escaping, and zero revision. Use no listener.

Adapt `TestWorkbenchSelectedEntry/object-store-error` to `refFailingStore{readStore: d.store, fail: commit.Entry.TransitionRef}`; assert root differs from transition, unselected 200 and selected 500. Adapt `TestWorkbenchReferenceWalk` after `worldStatus == 200` by extracting `refs, ok := workbenchRegion(stateBody, `<section aria-label="referencedBy">`, "</section>")`; require `ok` and `strings.Count(refs, "stateRoot: <a") == 3`. V13 measured that the old whole-body count becomes four after the nav edge. These adaptations are part of M1, never separate precursor landings.

In `TestWorkbenchTimelinePaging/emitted-links-resolve`, keep all three paging pages, timeline/selected regions, anchor census, classified-link fetches, and existing category controls. For each page extract `world, ok := workbenchRegion(body, worldStart, "</nav>")`; require `ok` and nonempty region; execute `regions = append(regions, struct{ name, text string }{"world", world})`. Assert the three selected-head roots are absent and have zero world object anchors. In the same store add a stored state object and a distinct world pointing to it, request `?world=<ref>` without changing SelectedHead, append its nav, and require exactly one world object anchor. Classify `world-home` and `world` separately, require both counts positive in addition to paging/select/stored-edge, require `strings.Count(region.text, "<a ") == len(matches)` for every region, decode `&amp;`, and fetch every anchor for 200. Empty extraction, omitted append, or omitted stored fixture must red.

M2's census parses `pageHTML` and `partialsHTML`, begins at `Page`, and traverses actions and template invocations with reflected dot types. Resolve every FieldNode segment exactly, dereference pointers, use element dot for `range`, value dot for `with`, outer dot for `if` and all `else` arms, inspect all pipeline commands/arguments, and guard recursion by template name plus argument type. Visit Page, TimelineView, EntryView, WorldView, ObjectView, GradeView, EdgeView, CommitView and ReferenceView. Require every type visited and every exported field consumed, with **no exemptions**. The live-template test and mutation self-tests must call the **same** reflected-field completeness function. Its field loop must name `typ` and `field` and use `seen[typ.Name()+"."+field.Name]` as its exact membership expression. `edgeUnavailable` receives only EdgeView and explicitly consumes EdgeView.Available; its true/false helper behavior gets direct tests. Unknown functions, unresolved paths, variables/assignments/chains or type-affecting nodes fail closed with a useful error until supported. Text/comments grant no credit. Self-tests remove World.Ref (with text/comment decoy), World.Unavailable, Object.Edges while Entry.Edges remains, and Commit.Truncated while Reference.Truncated remains; each includes a known-good template control and requires the precise missing owner/field.

No schema change, `/v1` route, store index or method, log-head lookup, guessed revision-to-entry index, log scan, AILANG kernel edit, or five-minute incident claim. The world check is evidence at response construction time, not a cross-request snapshot.

## M1 — checked world pane and recorder walk

**Files:** `host/daemon/workbench.go`, `host/workbench/render.go`, `host/daemon/workbench_test.go`, `host/daemon/workbench_references_test.go`, `host/workbench/render_test.go`.

**Acceptance** (run environment block above first; each command is copy-pasteable):

```sh
go test -race -count=1 ./host/daemon ./host/workbench -run '^(TestRenderWorldPane|TestWorkbenchWorldPane|TestWorkbenchTimelinePaging|TestWorkbenchSelectedEntry|TestWorkbenchReferenceWalk|TestWorkbenchCommitWalk)$' -v
go vet ./...
go test ./... -run '^$'
go test -race -count=1 ./host/daemon ./host/workbench
```

All names exist at this boundary. Observe stored-link follow 200, missing-root 200 with named stop and no root anchor, root error 500/503 without partial output, correct metadata and labels, four independently exercised link categories, and the retained three referencedBy edges. V13's `TestWorkbenchCommitWalk` stays green. Each mutation below is a **single, exactly-one-match literal replacement in the named file**, applied alone, then restored. Run its named test; compile failures are not kills. New-code literals above are binding spellings for the executor.

| Mutation | Named file: exact `old` → `new` (one match) | Must red |
|---|---|---|
| H1 | `host/daemon/workbench.go`: `stateRoot, err := d.checkedEdge(ctx, "stateRoot", world.StateRoot)` → `stateRoot, err := workbench.EdgeView{Relation: "stateRoot", Available: true, Target: world.StateRoot.String(), Href: "?object=" + world.StateRoot.String()}, error(nil)` | `TestWorkbenchWorldPane/missing-root`; paging link walk. |
| H2 | `host/daemon/workbench.go`: `d.checkedEdge(ctx, "stateRoot", world.StateRoot)` → `d.checkedEdge(ctx, "stateRoot", world.Ref)` | `TestWorkbenchWorldPane/stored-root`. |
| H3 | `host/daemon/workbench.go`: `d.checkedEdge(ctx, "stateRoot", world.StateRoot)` → `d.checkedEdge(ctx, "interface", world.StateRoot)` | `TestWorkbenchWorldPane/stored-root` relation. |
| H4 | `host/daemon/workbench.go`: `stateRoot, err := d.checkedEdge(ctx, "stateRoot", world.StateRoot)\n\t\tif err != nil {\n\t\t\td.writeWorkbenchStoreError(w, r, ctx, err)` → `stateRoot, err := d.checkedEdge(ctx, "stateRoot", world.StateRoot)\n\t\tif err != nil {\n\t\t\twriteWorkbenchError(w, 404, "NotFound", "missing root")` | `TestWorkbenchWorldPane/store-error/internal` and `/timeout`. |
| H5 | `host/daemon/workbench.go`: `stateRoot, err := d.checkedEdge(ctx, "stateRoot", world.StateRoot)\n\t\tif err != nil {\n\t\t\td.writeWorkbenchStoreError(w, r, ctx, err)` → `stateRoot, err := d.checkedEdge(ctx, "stateRoot", world.StateRoot)\n\t\tif err != nil {\n\t\t\td.writeWorkbenchInternalError(w, r, err)` | `TestWorkbenchWorldPane/store-error/timeout`. |
| H6 | `host/daemon/workbench.go`: `Revision:  world.Revision,` → `Revision:  0,` | `TestWorkbenchWorldPane/explicit-world`. |
| H7 | `host/daemon/workbench.go`: `LogHead:   world.LogHead.String(),` → `LogHead:   world.StateRoot.String(),` | `TestWorkbenchWorldPane/explicit-world`. |
| R1 | `host/workbench/render.go`: `{{template "edge" .World.StateRoot}}` → `` (empty string) | `TestRenderWorldPane/stored-root`; paging world control. |
| R2 | `host/workbench/render.go`: `<dt>revision</dt><dd>{{.World.Revision}}</dd>` → `<dt>revision</dt><dd>0</dd>` | `TestRenderWorldPane/metadata`. |
| R3 | `host/workbench/render.go`: `<dt>log head</dt><dd><span class="hash" title="{{.World.LogHead}}" aria-label="{{.World.LogHead}}">{{.World.LogHead}}</span></dd>` → `<dt>log head</dt><dd><span class="hash" title="{{.World.Ref}}" aria-label="{{.World.Ref}}">{{.World.Ref}}</span></dd>` | `TestRenderWorldPane/metadata`. |
| R4 | `host/workbench/render.go`: `<dt>world</dt><dd><span class="hash" title="{{.World.Ref}}" aria-label="{{.World.Ref}}">{{.World.Ref}}</span></dd>` → `<dt>world</dt><dd><span class="hash" title="literal-world" aria-label="literal-world">literal-world</span></dd>` | `TestRenderWorldPane/metadata`; M2 census also detects absent Ref. |
| R5 | `host/workbench/render.go`: `{{if edgeUnavailable .}}` → `{{if false}}` | `TestRenderWorldPane/missing-root`. |
| R6 | `host/workbench/render.go`: `<dt>world</dt>` → `<dt>stateRoot</dt>` | `TestRenderWorldPane/metadata`. |
| R7 | `host/workbench/render.go`: `UNAVAILABLE: {{.World.Unavailable}}</span>` → `UNAVAILABLE: </span>` | `TestRenderWorldPane/no-selected-world`; M2 census. |
| C1 | `host/daemon/workbench_test.go`: `regions = append(regions, struct{ name, text string }{"world", world})` → `` (empty string) | `TestWorkbenchTimelinePaging/emitted-links-resolve` independent controls. |
| C2 | `host/daemon/workbench_test.go`: `world, ok := workbenchRegion(body, worldStart, "</nav>")` → `world, ok := "", true` | `TestWorkbenchTimelinePaging/emitted-links-resolve` nonempty/anchor controls. |

H4/H5's `\n` and `\t` in this table denote literal newline and tab bytes in the gofmt file; assert a count of one on the fully expanded old string. Preserve `err` and the return so both mutants compile. R3/R4 replace a whole labelled element to mutate all its repeated field actions in **one** replacement. C1/C2's statement spellings are test implementation contracts; if a helper changes them, update the plan before claiming a drill result.

## M2 — remove dead timeline values and install typed census

**Files:** `host/daemon/workbench.go`, `host/workbench/render.go`, `host/workbench/render_test.go`.

**Acceptance** (run environment block above first):

```sh
go test -race -count=1 ./host/workbench ./host/daemon -run '^(TestWorkbenchViewFieldsAllRender|TestWorkbenchFieldCensusMutations|TestRenderWorldPane|TestWorkbenchWorldPane|TestWorkbenchTimelinePaging|TestWorkbenchSelectedEntry|TestWorkbenchReferenceWalk|TestWorkbenchCommitWalk)$' -v
go vet ./...
go test ./... -run '^$'
go test -race -count=1 ./host/daemon ./host/workbench
```

The M2 test file contains all named self-tests and the good-template controls. Count all nine visited types and zero exempt fields. Run A1–A4 singly, restoring between rows; require compile and the named behavioral red. A2's old string includes the **commit-specific prose** because `{{if .Truncated}}` occurs twice on line 178 (V20).

| Mutation | Named file: exact `old` → `new` (one match) | Must red |
|---|---|---|
| A1 | `host/workbench/render.go`: `{{with .Object}}{{range .Edges}}` → `{{with .Object}}{{range .Commits.Edges}}` | `TestWorkbenchViewFieldsAllRender` reports ObjectView.Edges. |
| A2 | `host/workbench/render.go`: `{{if .Truncated}}<p>Showing 100 commits; more recorded</p>` → `{{if false}}<p>Showing 100 commits; more recorded</p>` | `TestWorkbenchViewFieldsAllRender` reports CommitView.Truncated. |
| A3 | `host/workbench/render.go`: `type ObjectView struct {` → `type ObjectView struct {\n\tForgotten string` | `TestWorkbenchViewFieldsAllRender` reports ObjectView.Forgotten. |
| A4 | `host/workbench/render_test.go`: `seen[typ.Name()+"."+field.Name]` → `strings.Contains(pageHTML+partialsHTML, "."+field.Name)` | `TestWorkbenchFieldCensusMutations/world-ref-decoy` detects the analyzer's false credit. |

A3's `\n`/`\t` represent newline/tab bytes. A4 requires the exact expression in the field loop; it must compile because `strings` remains imported. The self-test's decoy must make the weak substring result true while its typed traversal still reports missing `WorldView.Ref`. An analyzer test that only checks the live template would not kill A4.

## Closure and verification log

After M2: rerun its focused command and all three boundary gates, then `go test -count=1 ./...` and `./scripts/verify_ail.sh` under the stated environment. Record actual rc, test counts, race result, added production LOC, one-match counts, compile and named red for each of **20** mutations, restoration green, and CI result. Do not describe design-time mutation rows as measured kills. V21 established that `-run 'Workbench\|Render\|Grade'` runs **zero** tests; every command here uses an unescaped `|`. Keep the conflict surface to the five files in design §9 and recheck it before landing. A final successful walk provides verified links and named stops; it does not itself establish the mission's five-minute diagnosis bar.
