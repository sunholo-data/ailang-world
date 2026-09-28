# w-workbench-world-pane-checked — rows 99 + 100

- Status: planned; designer iteration 205; revision pass round 2; 2026-09-28.
- Measurement base: `aa3e36a` (V0). One bundled Go-host change; no AILANG kernel change.
- Estimate: about 0.6 day; two independently testable milestones. The controller owns commits.
- This document is the only deliverable. Round-1 probe source is reproduced below; those probes ran in `/tmp/world205-probe`. Controller measurements V11–V13 ran at base `aa3e36a` in `/Users/voightkampff/dev/sunholo-data/.probe-world-iter205`; round-2 source checks ran here (V14–V16).

## 1. Problem and clause mapping

The world pane offers a state-object URL under the world reference's label, without checking that the object exists. At base, both `/workbench` and `/workbench?from=0&entry=0` return 200 but their world-pane object link returns 404 (V1, V4). The pane also hides revision and log head. Timeline From/Limit are exempt from the field census, and WorldView/ObjectView are outside its checked types (V2, V5).

Rows 99 and 100 share the world constructor and template block. Together they support clause 5's requirement that a provenance walk answer a real “why did X happen?” question with verified evidence within five minutes (task mandate). The proposed observable is a trustworthy entry point: an identified world, its revision and log head, and either a working stateRoot hop or an explicit missing-object stop. This is a prerequisite improvement, **not evidence that the complete clause-5 bar is met**. No synthetic timing is presented as a real incident diagnosis.

## 2. Measured current state

Every base-code claim in this section is tied to the verification log.

1. `checkedEdge(ctx, relation, ref)` already performs one GetObject: found yields Available/Target/Href; missing yields `object <ref> is not stored`; errors propagate. The world constructor bypasses it and hardcodes Available true (V1).
2. With `seedGenesisEmbedded(...,"probe")` and `testCommit(genesis,0,"probe")`, both genesis and committed state roots are absent. The committed root is `sha256:ad312da471c6714e01bc3014671b0145a25d55e8b360470df959eed65283f170`. The rendered anchor's text is the different world ref `sha256:dde9258fd6b815e89c8bae1633c19aa91a72b376474d0b63591aacb6c9c757e3`. Both requested pages emit the dead link (V4). The committed transition object is a positive stored-object control.
3. Storing payload `state-0-probe` under its hash makes that same root resolvable (GetObject true and object page 200). Existing `TestWorkbenchCommitWalk` and `TestWorkbenchReferenceWalk` fixtures also set StateRoot to `state.Hash` and carry that object; their state walks return 200 (V4, V7). Both branches are testable without any schema change.
4. Extending the actual lexical ratchet to all nine view structs fails on **WorldView.Revision and WorldView.LogHead only**, with From/Limit still exempt. A separate census without exemptions reports exactly four omissions: TimelineView.From, TimelineView.Limit, WorldView.LogHead, WorldView.Revision. ObjectView contributes no additional omission (V5).
5. That result is insufficient assurance: replacing every `{{.World.Ref}}` with literal text leaves the lexical missing set unchanged, whereas a typed template-tree traversal adds **WorldView.Ref**. `.References` still satisfies the weak `.Ref` substring check. The typed traversal also distinguishes duplicate names such as Edges, Available and Truncated by owning struct (V5).
6. `readStore` has GetLogEntry by numeric index, but no entry-hash lookup. The store's GetLogEntry SQL filters `entry_index`; the schema does have UNIQUE `entry_hash_ref`. Thus an indexed SQL lookup could be added, but it is not a read exposed through this seam today. Revision must not be assumed to be an entry index. This design chooses no seam extension (V3).
7. From/Limit are copied into TimelineView but not read back by the handler; its loop and link logic use local `from` and `limit`. The positive search control finds NextHref assignment (V3).
8. `emitted-links-resolve` extracts timeline and selected-entry regions only, classifies every matched anchor, checks anchor counts, and requires paging/select/stored-edge categories to fire. It omits nav. A selected-entry error test uses an all-object-failing seam and expects the unselected page to succeed (V2). That control will need correction.
9. The existing store-error writer maps deadline errors to HTML 503 and ordinary internal errors to HTML 500, logging operator detail. It also maps ReferenceIndexUnavailableError to 503. A temporary world-check insertion measured 200 for the missing-root non-error, 500 for the sentinel storage error, and 503 for DeadlineExceeded (V1, V6).

10. `WorldView.Unavailable` already exists at base and the nav else-branch already renders `{{.World.Unavailable}}` (V11). This design adds no field; only the default reason `no world selected` is new.

## 3. Decisions, with alternatives

### D1 — Check StateRoot with checkedEdge

Use `d.checkedEdge(ctx, "stateRoot", world.StateRoot)` after successful GetWorld and before assigning Page.World. On error call `writeWorkbenchStoreError` and return before rendering. Missing root is a valid 200 world page with a named stop; it does not make the world itself unavailable.

Alternatives: hardcoded availability repeats the measured 404; checking in the template cannot access the store; duplicating GetObject handling would duplicate the already measured contract (V1, V4, V6). Cost: one additional object read per request with a world, including object pages. No per-timeline-row world checks.

### D2 — Separate world identity from the stateRoot link; render revision and log head

Render world ref as text, revision as a decimal integer, log head as text, then the existing edge partial with relation `stateRoot`. The link text/title/aria-label must all identify the **state object**, not the world. Retain the workbench home link.

Alternatives: keeping the world-labelled anchor misidentifies its target (V4); making the world ref a self-link adds navigation with no new information. Linking LogHead requires its entry index, which the current seam cannot obtain directly by hash (V3). A new store method could exploit the unique hash column, but is unnecessary for this small bundle. Do not guess from Revision or scan the log. Log head remains a full copyable hash with the existing visual hash treatment, explicitly labelled `log head`; genesis sentinel values are text too.

### D3 — Delete TimelineView.From and Limit; keep the handler locals

Delete the two fields, their constructor assignments, and both exemptions. Retain local `from`, `limit`, bounds, row iteration and paging behavior. These duplicated view values add no information needed by the chosen UI (V3, V5).

Alternative: render “from N, limit 100.” That exposes pagination bookkeeping rather than provenance and would retain otherwise unused copies. WorldView.Revision and LogHead, in contrast, identify the selected world's recorded state and must render (D2). No other ObjectView/WorldView field deletion or new exemption is needed by the measured census (V5).

### D4 — Replace lexical matching with a type-qualified template-tree census

Cover Page, TimelineView, EntryView, WorldView, ObjectView, GradeView, EdgeView, CommitView and ReferenceView. Traverse parsed template actions with the concrete Go type of dot. Record `(owning struct, exact field)` on every field-path segment. No substring matching and no global-name satisfaction. The base and mutation measurements in V5 justify this extra strength.

Alternative: merely adding WorldView/ObjectView to the lexical list discovers the two known missing actions but misses the measured Ref deletion. A handwritten expected-path list alone could drift when a field is added. Reflection must enumerate all exported fields and compare with actual traversal results.

The census is an **action-use ratchet**, not proof that each value becomes visible text: condition fields are intentionally consumed as conditions. Focused rendering tests must pin visible output, branch behavior, escaping and labels. Do not advertise this static check as semantic rendering verification.

### D5 — Widen the link walk with a separately controlled world category

Extract `<nav aria-label="world browser">` through `</nav>` on each existing paging page and the additional fixtures below. Classify its home anchor as `world-home` and its checked state-object anchor as `world`. Require both categories independently to have positive counts. Require a stored-root page with exactly one world object link. A home link alone must not satisfy the world-edge control.

Alternative: scanning nav without category controls admits a dead extractor; scanning only absent-root fixtures exercises no world edge. The existing region/classification machinery supplies the walk controls (V2). The paging fixture’s selected-head root is absent, so all three original paging pages are missing-root cases (V12); add a stored-root fixture to exercise the world edge.

## 4. Design

### 4.1 Daemon construction

The intended replacement at the successful GetWorld branch is:

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

Initialize the existing Page.World.Unavailable field (already rendered at base, V11) to the new default `no world selected` when constructing Page, so an absent selected head has a named stop. The successful assignment replaces that default. A missing requested world still follows the existing GetWorld-not-found 404 path (V1). Remove `page.Timeline = workbench.TimelineView{From: from, Limit: limit}`; the Page zero value supplies the empty timeline before appends and link assignments.

The existence guarantee is scoped to the read that built this response, and link-resolution tests use an unchanged store. It is not a cross-request transaction or promise against subsequent storage failure.

### 4.2 View model and exact world markup

Remove only TimelineView.From/Limit. Keep all WorldView and ObjectView fields. Within the existing nav, replace its world conditional with:

```html
{{if .World.Available}}
<dl>
<dt>world</dt><dd><span class="hash" title="{{.World.Ref}}" aria-label="{{.World.Ref}}">{{.World.Ref}}</span></dd>
<dt>revision</dt><dd>{{.World.Revision}}</dd>
<dt>log head</dt><dd><span class="hash" title="{{.World.LogHead}}" aria-label="{{.World.LogHead}}">{{.World.LogHead}}</span></dd>
</dl>
{{template "edge" .World.StateRoot}}
{{else}}<span class="unavailable" role="note">UNAVAILABLE: {{.World.Unavailable}}</span>{{end}}
```

The nav retains `<a href="{{workbenchHref ""}}">workbench</a>`. This is a route link, not a stored-object edge. Every data-target link in the world pane must come from checkedEdge. The `edge` partial already renders the relation plus either missing text or an anchor whose identity is Target (V1). No CSS change is required by this design.

### 4.3 Type-aware ratchet

Implement the test helper in render_test.go, using the template parse tree before HTML execution. Start at Page. Resolve exact `FieldNode.Ident` chains through reflection, dereferencing pointers. A range changes dot to its collection element; a with changes dot to its value; if keeps dot; all else arms retain the outer dot. A template invocation traverses that named template with the argument's type. Inspect every command argument and pipeline, not just the first field. Visit both branches statically, including nil/empty cases. Key recursion protection by template name and argument type. In the reflected field loop, name the field `field` and its owner `typ`, and use `seen[typ.Name()+"."+field.Name]` as the exact membership expression; this fixes the A4 mutation anchor.

The present `edgeUnavailable .` helper reads EdgeView.Available indirectly (V1). Use an explicit function-read contract: only an EdgeView argument to `edgeUnavailable` records EdgeView.Available. Test that helper independently with true and false values. Do not exempt Available globally or infer every field as consumed when a whole struct is passed. All other fields must be visible in actual field actions.

Fail closed on unsupported type-affecting nodes, unknown function contracts or unresolved paths. If future templates introduce variables, assignments, chained expressions or new helpers, teach the analyzer their types or fail with an actionable error; silently skipping them is forbidden. Text and comments must not count as field use. Register the known presentation functions and builtins with sufficient signatures for traversal; they grant no arbitrary field credit.

Reflect over all nine types; require each type to have been visited and every exported field to be accounted for with no exemptions. Add self-tests that remove World.Ref, Object.Edges while leaving Entry.Edges, and CommitView.Truncated while leaving ReferenceView.Truncated. Also remove the `{{.World.Unavailable}}` action and require the precise omission `WorldView.Unavailable` (the field and action already exist, V11). Include a text/comment `.World.Ref` decoy. Each must report the precise missing owner/field. A known-good template is the positive control in the same test.

### 4.4 Handler and render tests

Keep the test packages and helpers separate (V15). `host/workbench/render_test.go` owns `TestRenderWorldPane`: use `renderPage(t, Page{...})`, which calls `Render(&body, p)` directly. Extract nav with a test-local string-region helper; use neither requestRecorder nor newHandlerDaemon. `host/daemon/workbench_test.go` owns `TestWorkbenchWorldPane` and the paging walk: use newHandlerDaemon, requestRecorder, workbenchRegion and test-local stores. No listener or live server.

Declare these exact names; mutation rows in §7 refer to them:

- `TestRenderWorldPane/metadata`, `TestRenderWorldPane/stored-root`, `TestRenderWorldPane/missing-root`, `TestRenderWorldPane/no-selected-world`, `TestRenderWorldPane/escaping`, `TestRenderWorldPane/revision-zero`.
- `TestWorkbenchWorldPane/missing-root`, `TestWorkbenchWorldPane/stored-root`, `TestWorkbenchWorldPane/explicit-world`, `TestWorkbenchWorldPane/no-selected-world`, `TestWorkbenchWorldPane/store-error/internal`, `TestWorkbenchWorldPane/store-error/timeout`.
- Existing test targets retained: `TestWorkbenchTimelinePaging/emitted-links-resolve` (V2), `TestWorkbenchViewFieldsAllRender` (V5). The new census self-test inventory is `TestWorkbenchFieldCensusMutations/world-ref-decoy` (WorldView.Ref), `TestWorkbenchFieldCensusMutations/world-unavailable` (WorldView.Unavailable), `TestWorkbenchFieldCensusMutations/object-edges` (ObjectView.Edges), and `TestWorkbenchFieldCensusMutations/commit-truncated` (CommitView.Truncated), each paired with a known-good control.

- `TestRenderWorldPane`: `metadata` uses distinct world/root/log-head sentinels and nonzero revision 17. Assert labelled values inside nav, no world-ref or log-head anchor, root anchor identity in text/title/aria-label, unavailable-root text with no object anchor, and unavailable-world reason. Use `stored-root`, `missing-root` and `no-selected-world` for their respective branches; cover zero and malicious labels in `revision-zero` and `escaping`. The home anchor is expected in every case.
- `TestWorkbenchWorldPane`: default, explicit `?world=<ref>`, selected entry, and object-page requests. Use one testCommit root absent from GetObject, plus a world carrying a stored root (set NextWorld.StateRoot to an object included in Commit.Objects, as in V7). Run default, selected-entry and object-page requests within both `missing-root` and `stored-root`. Check GetObject true/false controls before rendering. On the absent-root page assert exact `stateRoot: ... UNAVAILABLE: object <ref> is not stored`, and **zero state-object anchors**; nav still contains home. Follow the stored-root anchor and observe 200. Assert identity, revision and log head against the fixture, not against values parsed back from HTML.
- `TestWorkbenchWorldPane/explicit-world`: include an explicit older-world request while a different head is selected, with different root, revision and log head. V14 shows the handler resolves the explicit query ref through GetWorld(ref), bypassing SelectedHead. This case kills accidental use of the selected head's metadata for `?world=`.
- `TestWorkbenchWorldPane/no-selected-world`: fresh store, no selected head; observe `UNAVAILABLE: no world selected`, one home anchor, no root link.
- `TestWorkbenchWorldPane/store-error/internal` and `TestWorkbenchWorldPane/store-error/timeout`: a ref-specific wrapper fails only the root. Untargeted ref control returns 200; root sentinel error returns 500/Internal without sentinel detail in the response; root DeadlineExceeded returns 503/Timeout. Check error responses have no partial nav success body. In `stored-root`, a counting wrapper pins one root GetObject on `/workbench` (avoid an object or selected page that legitimately reads the same ref elsewhere).
- Adapt `TestWorkbenchSelectedEntry/object-store-error`: fail only `commit.Entry.TransitionRef`, using the existing refFailingStore (V15). The unselected request must still be 200, and selected request 500. Distinct root/transition refs are a fixture control. Do not weaken this test to allow either status.

Extend `TestWorkbenchTimelinePaging/emitted-links-resolve` rather than replacing it. Keep its original three pages and category controls (V2). For each, extract with `world, ok := workbenchRegion(body, worldStart, "</nav>")`, fail on `!ok`, and add `regions = append(regions, struct{ name, text string }{"world", world})`. Require a nonempty world region, classify every anchor, and enforce `strings.Count(region,"<a ") == len(matches)`. Preserve errors for unclassified links and decode `&amp;` before fetching.

Add an explicit stored-root world page in the **same daemon/store**: PutObject a state fixture, PutWorld a distinct world referencing it, then request `?world=<ref>` without changing SelectedHead. The original three paging pages all render the same selected head and remain missing-root cases (V12). V14 verifies that the explicit query loads the distinct world without selecting it. Assert GetObject false for their root and true for the new root; assert missing text and absence of object anchors on the former, exactly one root anchor on the latter. Require `world-home >= 1` and `world >= 1` independently of paging/select/stored-edge. Fetch every classified link and require 200. A world extraction returning empty, a forgotten region append, and a missing stored fixture must each fail.

Adapt `TestWorkbenchReferenceWalk` in `host/daemon/workbench_references_test.go`: after requiring `worldStatus == 200`, extract ``refs, ok := workbenchRegion(stateBody, `<section aria-label="referencedBy">`, "</section>")``, fail on `!ok`, and require `strings.Count(refs, "stateRoot: <a") == 3`. The base assertion counts the whole body (V13, V15); the new nav adds a fourth edge. Keep the three-edge referencedBy guarantee. Bumping the body-wide literal to 4 would couple the reference walk to nav markup. `TestWorkbenchCommitWalk` needs no change under the measured prototype (V13).

## 5. Failure modes

| Failure | Required observation / containment |
|---|---|
| Root absent | World metadata remains visible; 200; named unavailable root; no root anchor. |
| Root read fails | 500/Internal, or 503/Timeout for deadline; no successful page partially emitted. |
| Wrong edge label or world hash substituted | Distinct sentinels make nav assertion fail even if the URL resolves. |
| Log-head hash treated as object/index | No log-head anchor is permitted; no extra hash-resolution read. |
| Shared field names hide omissions | Typed census reports the owning struct; mutation self-tests exercise collisions. |
| World extractor sees nothing | Region presence, exact stored-root anchor count and independent category controls fail. |
| Root check/nav breaks existing tests | V13 measured exactly two: `TestWorkbenchSelectedEntry/object-store-error` loses its 200 control because all object reads fail; narrow failure to the transition. `TestWorkbenchReferenceWalk` counts four stateRoot links over the whole body; scope its count to referencedBy and retain exactly three. |
| Concurrent store changes between requests | Checked at render time only; no snapshot guarantee introduced. |
| Empty world | Explicit reason; no speculative edge. |

## 6. Milestones and acceptance observations

Use `export AILANG_BIN="$HOME/.pinned-ailang/ailang"`. Commands below are execution gates, not claims that the proposed tests already exist. Tests use temporary storage and the normal external Go cache; no acceptance artifact is written into `$HOME` or outside the worktree. Test temp/cache activity is not an acceptance artifact. Do not set GOCACHE/TMPDIR to the worktree. The controller handles all commits.

**M1 — checked world pane and visible world metadata.** Implement daemon check, nav markup, missing-world reason and focused world tests; widen emitted-links-resolve; narrow the selected-entry failure wrapper; scope the reference-walk stateRoot count to referencedBy. Retain the old field ratchet and From/Limit until M2. This milestone fixes row 99 and the two world metadata omissions without depending on M2.

```sh
AILANG_BIN="$HOME/.pinned-ailang/ailang" go vet ./host/...
AILANG_BIN="$HOME/.pinned-ailang/ailang" go test -race ./host/workbench ./host/daemon -run 'Workbench|Render|Grade' -count=1
```

Observe green packages; stored world-root targets 200, missing roots rendered without object anchors, targeted storage failure 500, timeout 503, and all four world metadata values correctly labelled. Drill M1 mutations below. No test may skip because a socket is unavailable: these tests do not need one.

**M2 — remove dead view fields and strengthen the ratchet.** Delete From/Limit and exemptions, replace the ratchet with typed traversal and its mutation self-tests. Add no new UI. M1 focused tests remain green.

```sh
AILANG_BIN="$HOME/.pinned-ailang/ailang" go vet ./host/...
AILANG_BIN="$HOME/.pinned-ailang/ailang" go test -race ./host/workbench ./host/daemon -count=1
```

Observe all nine types accounted for, no omitted fields or exemptions, and named analyzer self-test failures for each deliberately damaged template. Run M2 mutations below. Broader CI uses race detection (task mandate); execution must report its actual result, not infer it from the design-time focused baseline.

## 7. Test plan and MUTATION TABLE

Every mutation applies singly to the proposed implementation, then is reverted before the next. These are required future drill outcomes, **not measured kills at design time**, except the explicitly measured lexical-vs-typed experiment V5. A mutation must land exactly once, compile under `go vet ./host/...`, and fail its named test; compile failures are not kills. The table uses the exact intended code/markup above as anchors.

| ID / milestone | Exact mutation | Test that must kill it |
|---|---|---|
| H1 / M1 | Replace `stateRoot, err := d.checkedEdge(ctx, "stateRoot", world.StateRoot)` with `stateRoot, err := workbench.EdgeView{Relation: "stateRoot", Available: true, Target: world.StateRoot.String(), Href: "?object=" + world.StateRoot.String()}, error(nil)` | `TestWorkbenchWorldPane/missing-root`; `TestWorkbenchTimelinePaging/emitted-links-resolve` missing-page assertion |
| H2 / M1 | At that call only, replace `world.StateRoot` with `world.Ref` | `TestWorkbenchWorldPane/stored-root` |
| H3 / M1 | At that call only, replace `"stateRoot"` with `"interface"` | `TestWorkbenchWorldPane/stored-root` exact relation assertion |
| H4 / M1 | In the new error arm only, replace `d.writeWorkbenchStoreError(w, r, ctx, err)` with `writeWorkbenchError(w, 404, "NotFound", "missing root")` | `TestWorkbenchWorldPane/store-error/internal`; `TestWorkbenchWorldPane/store-error/timeout` |
| H5 / M1 | In the new error arm only, replace that call with `d.writeWorkbenchInternalError(w, r, err)` | `TestWorkbenchWorldPane/store-error/timeout` |
| H6 / M1 | In WorldView assignment replace `Revision: world.Revision` with `Revision: 0` | `TestWorkbenchWorldPane/explicit-world` nonzero revision |
| H7 / M1 | In WorldView assignment replace `LogHead: world.LogHead.String()` with `LogHead: world.StateRoot.String()` | `TestWorkbenchWorldPane/explicit-world` |
| R1 / M1 | Replace `{{template "edge" .World.StateRoot}}` with empty text | `TestRenderWorldPane/stored-root`; `TestWorkbenchTimelinePaging/emitted-links-resolve` world category CONTROL |
| R2 / M1 | Replace `{{.World.Revision}}` with literal `0` | `TestRenderWorldPane/metadata` |
| R3 / M1 | Replace every `{{.World.LogHead}}` in nav with `{{.World.Ref}}` | `TestRenderWorldPane/metadata` |
| R4 / M1 | Replace every `{{.World.Ref}}` in nav with `literal-world` | `TestRenderWorldPane/metadata`; `TestWorkbenchViewFieldsAllRender` at M2 (weak census misses, measured V5) |
| R5 / M1 | In edge partial replace `{{if edgeUnavailable .}}` with `{{if false}}` | `TestRenderWorldPane/missing-root` |
| R6 / M1 | Replace `<dt>world</dt>` with `<dt>stateRoot</dt>` | `TestRenderWorldPane/metadata` label assertion |
| R7 / M1 | Delete the `{{.World.Unavailable}}` action from the nav else-branch | `TestRenderWorldPane/no-selected-world`; `TestWorkbenchWorldPane/no-selected-world`; `TestWorkbenchViewFieldsAllRender` at M2 |
| C1 / M1 | Delete `regions = append(regions, struct{ name, text string }{"world", world})` | `TestWorkbenchTimelinePaging/emitted-links-resolve` independent world/world-home controls |
| C2 / M1 | At the world extraction only, replace `workbenchRegion(body, worldStart, "</nav>")` with `"", true` | `TestWorkbenchTimelinePaging/emitted-links-resolve` missing-text / exact stored-root anchor assertions and category controls |
| A1 / M2 | In `{{with .Object}}{{range .Edges}}` replace `.Edges` with `.Commits.Edges` (only this range) | `TestWorkbenchViewFieldsAllRender`: ObjectView.Edges absent despite other Edges uses |
| A2 / M2 | Inside `{{with .Commits}}`, replace `{{if .Truncated}}` with `{{if false}}`, retaining ReferenceView.Truncated | `TestWorkbenchViewFieldsAllRender`: CommitView.Truncated absent |
| A3 / M2 | Add `Forgotten string` to ObjectView with no template action | `TestWorkbenchViewFieldsAllRender`: ObjectView.Forgotten absent |
| A4 / M2 | Replace `seen[typ.Name()+"."+field.Name]` in the census field loop with `strings.Contains(pageHTML+partialsHTML, "."+field.Name)` | `TestWorkbenchFieldCensusMutations/world-ref-decoy` must fail because the mutant analyzer no longer detects the missing field |

C1/C2 and A4 use the helper statements prescribed in §4.3–4.4; preserve those anchors or record an equivalent exact replacement in the sprint plan. Do not claim a landed mutant merely from a substring replacement count of zero. Self-tests operate on in-memory template strings and do not rewrite repository files.

## 8. Non-goals

No schema change, new `/v1` route, new store index, entry-hash read method, revision-to-index assumption, log scan, query grammar change, grade derivation, kernel policy, or payload interpretation. No redesign of other provenance panes. No new performance or five-minute mission claim. No separate sprint-plan artifact in this iteration.

## 9. Conflict Surface

Future implementation touches only:

| File | Changes |
|---|---|
| `host/daemon/workbench.go` | world root check/error path, missing-world reason, removal of timeline initializer |
| `host/workbench/render.go` | world nav, TimelineView fields |
| `host/daemon/workbench_test.go` | world fixtures/assertions, world link region/control, selected-entry error seam |
| `host/daemon/workbench_references_test.go` | scope reference-walk stateRoot count to referencedBy |
| `host/workbench/render_test.go` | world render cases, typed census and mutation self-tests |

V13 supersedes the round-1 claim that reference/commit walks need no changes: the combined handler/nav prototype broke exactly `TestWorkbenchReferenceWalk` and `TestWorkbenchSelectedEntry/object-store-error`. Both adaptations are specified in §4.4. `TestWorkbenchCommitWalk` and every other daemon/workbench test passed that prototype. This measured blast radius changes the test scope and M1 work; decisions D1–D5 remain unchanged. The reference documents used for structure and prior findings are `design_docs/implemented/w-workbench-object-provenance-and-grade.md`, its sprint plan, and `design_docs/implemented/w-workbench-timeline-seam.md` §9 (V8). This design-time diff contains only this document.

## 10. VERIFICATION LOG

V0–V10 record round-1 work; their status and line counts are historical. V11–V13 are controller-verified measurements at `aa3e36a` in `/Users/voightkampff/dev/sunholo-data/.probe-world-iter205`, outside the sandbox, supplied for this revision; V11 was also rechecked here. V14–V16 are round-2 read-only checks in this worktree. Other commands ran here unless a `/tmp/world205-probe` working directory is specified. Output below is trimmed. No skill was read, no git write command was run, and no server was started. The probes use the existing requestRecorder helpers; they opened no listening socket. There was no sandbox bind result to classify.

| ID | Exact command / reproducible procedure | Observed output |
|---|---|---|
| V0 | `git status --short; git rev-parse --short HEAD` | status empty; `aa3e36a`. |
| V1 | `sed -n '100,200p' host/daemon/workbench.go; sed -n '350,410p' host/daemon/workbench.go; cat host/workbench/render.go` | Error writer: timeout/reference-index unavailable 503, internal 500. checkedEdge has GetObject found/missing/error branches. World assignment unconditional Available true and no check. Nav world-ref anchor targets StateRoot.Href; partial uses Target and Missing. Full view structs and template actions visible. |
| V2 | `sed -n '320,365p' host/workbench/render_test.go; sed -n '490,648p' host/daemon/workbench_test.go` | Exempt From/Limit; lexical Contains; types EntryView, TimelineView, Page. Error control expects unselected 200 with objectFailingStore. Three-page walk regions timeline/selected, anchor census and three category controls. |
| V3 | Run source-instrument command in §10.1, plus `sed -n '485,555p' host/daemon/workbench.go; sed -n '609,632p' host/store/store.go; sed -n '25,47p' host/store/schema.sql` | GetLogEntry method count 1; GetLogEntryByHash / LogEntryIndex 0 with full eight-method seam printed as control. page.Timeline.From/Limit reads 0/0, NextHref 1. Loop uses locals. SQL filters entry_index; entry_hash_ref UNIQUE. |
| V4 | Copy tracked files and run daemon probe in §10.2: `AILANG_BIN="$HOME/.pinned-ailang/ailang" go test ./host/daemon -run '^TestProbe205$' -v -count=1` in copy | PASS as a defect-measuring test. Genesis root false; committed root false; committed-object control true. Both pages 200, root target 404, anchor text world ref. After exact preimage PutObject: GetObject true, target 200. checkedEdge returns both injected errors; base `/workbench` incorrectly remains 200 for both. |
| V5 | Extend actual ratchet as §10.3 and run `AILANG_BIN="$HOME/.pinned-ailang/ailang" go test ./host/workbench -run '^TestWorkbenchViewFieldsAllRender$' -v -count=1`; run §10.4 probe with `-run '^TestProbe205Census$'` | Extended actual test FAIL: WorldView.Revision and WorldView.LogHead. Probe PASS: both base censuses omit the four named fields; lexical mutated set unchanged; typed mutated set additionally omits WorldView.Ref. Positive base/field controls in same run. |
| V6 | Apply §10.5 temporary world-construction patch and run `AILANG_BIN="$HOME/.pinned-ailang/ailang" go test ./host/daemon -run '^TestProbe205Checked$' -v -count=1` | PASS: nil error 200; sentinel error 500; DeadlineExceeded 503. This is a handler-path measurement, not a full implementation prototype. |
| V7 | `sed -n '235,259p' host/daemon/workbench_commits_test.go; sed -n '280,312p' host/daemon/workbench_references_test.go; AILANG_BIN="$HOME/.pinned-ailang/ailang" go test ./host/daemon -run 'TestWorkbench(CommitWalk\|ReferenceWalk)$' -v -count=1` | Both fixtures set StateRoot=state.Hash and carry state. Both PASS; identical-state carried-by answer [0 1 2], object/followed statuses all 200; reference walk reports state 200. |
| V8 | `rg -n '^##' design_docs/implemented/w-workbench-object-provenance-and-grade.md; sed -n '1,60p' design_docs/implemented/w-workbench-object-provenance-and-grade-sprint-plan.md; sed -n '445,469p' design_docs/implemented/w-workbench-timeline-seam.md` | House doc sections include decisions/design/conflict/acceptance/mutations/milestones/log. Plan has two milestones and compile/drill gates. Seam §9 explicitly records R-a and R-c. Historical claims were not substituted for current probes. |
| V9 | `~/.pinned-ailang/ailang --version; AILANG_BIN="$HOME/.pinned-ailang/ailang" go vet ./host/...; AILANG_BIN="$HOME/.pinned-ailang/ailang" go test -race ./host/workbench ./host/daemon -run 'Workbench\|Render\|Grade' -count=1` | AILANG v0.41.0, commit 24ee108; vet rc=0; workbench ok 1.277s, daemon ok 11.822s. Baseline only, not proposed implementation. |
| V10 | `git status --short; git diff --name-only; wc -l design_docs/planned/w-workbench-world-pane-checked.md` | Only `?? design_docs/planned/w-workbench-world-pane-checked.md`; tracked diff empty; final document 395 lines. |


| V11 | Controller: `grep -n 'Unavailable' host/workbench/render.go` at base; same command re-run here. | Line 35: `Unavailable string` (GradeView); line 100: `Unavailable string` (WorldView); line 150: nav else-branch already uses `{{.World.Unavailable}}`. The field exists and is rendered; this design adds no field. Other matches include the edge helper and grade rendering. |
| V12 | Controller daemon-package probe: `newHandlerDaemon(t); seedWorkbenchLog(t, d, WorkbenchPageLimit+5)`; read SelectedHead → GetWorld → GetObject(world.StateRoot), then GetObject(entry0.TransitionRef) as positive control. | `PAGING head=sha256:a584866e… revision=104 stateRoot=sha256:7b247566… stored=false err=<nil>`; `CONTROL entry0 transitionRef stored=true`. The paging head root is not stored. All three paging pages render that same selected head, so all three are missing-root pages. Controller supplied the probe procedure and output, not its shell invocation; it was not re-run in round 2. |
| V13 | Controller applied BOTH §10.5 handler patch and §4.2 nav markup verbatim to the base probe worktree; ran `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/daemon ./host/workbench -count=1`. | host/workbench ok; host/daemon FAIL with exactly TWO failing tests: (a) `TestWorkbenchReferenceWalk`, workbench_references_test.go:347–348, whole-body `strings.Count(stateBody, "stateRoot: <a") != 3` sees 4 because the stored selected-head root (V7) adds a nav edge; (b) `TestWorkbenchSelectedEntry/object-store-error`, workbench_test.go:506, `control: /workbench status = 500, want 200` because its all-object-failing seam now fails the root read. `TestWorkbenchCommitWalk` and every other daemon/workbench test passed. Not re-run in round 2. |
| V14 | `sed -n '345,398p' host/daemon/workbench.go` | `if values := query["world"]; values != nil` parses the explicit ref and assigns `worldRefText = ref.String()`; only the else-branch calls SelectedHead. The common branch parses worldRefText and calls `d.reads.GetWorld(ctx, ref)`. Thus `?world=<ref>` resolves that world, not the selected head. |
| V15 | `rg -n 'func (requestRecorder\|newHandlerDaemon\|workbenchRegion)\|refFailingStore\|object-store-error' host/daemon/*test.go`; `rg -n -A 9 'func renderPage' host/workbench/render_test.go; sed -n '390,412p' host/daemon/workbench_test.go; sed -n '490,515p' host/daemon/workbench_test.go; sed -n '20,48p' host/daemon/handlers_test.go; sed -n '335,355p' host/daemon/workbench_references_test.go; rg -n 'section aria-label="referencedBy"' host/workbench/render.go` | renderPage:241 uses bytes.Buffer and Render(&body, p). newHandlerDaemon uses a temp DB; requestRecorder directly calls Handler().ServeHTTP with httptest, without a listener. workbenchRegion returns from the start marker to the first end marker. Reference walk counts `stateRoot: <a` over stateBody; template line 178 supplies the referencedBy section marker. Selected-entry seam uses objectFailingStore and expects unselected 200; refFailingStore is declared at workbench_test.go:750 with GetObject at :755. |
| V16 | Before revision: `git status --short; wc -l design_docs/planned/w-workbench-world-pane-checked.md; git rev-parse --short HEAD` | Status empty; document 395 lines; current worktree HEAD `95b3b45`. Historical measurement base remains `aa3e36a`; controller V11–V13 explicitly used that base. |
| V17 | After revision: `git diff --check; git diff --stat; git status --short; wc -l design_docs/planned/w-workbench-world-pane-checked.md` | Diff check emits no diagnostics; positive controls show exactly one modified file, this document, and 427 lines. No code edits or git writes. |
| V18 | Controller, round 3 (glm/kimi r2 fix): `git log --oneline aa3e36a..61ac489 -- host/ \| wc -l`; `git diff --stat aa3e36a..61ac489 -- host/ \| wc -l`; control `git log --oneline aa3e36a..61ac489 \| wc -l` | `0`; `0`; control `2` (the two design-doc commits). The host tree at the implementation target is byte-identical to measurement base `aa3e36a`, so V1–V7, V11–V15 carry over unchanged. |
| V19 | Controller, round 3 (glm r2 fix): `rg -n 'func writeWorkbenchError\|func \(d \*Daemon\) writeWorkbenchInternalError\|func \(d \*Daemon\) writeWorkbenchStoreError' host/daemon/` | `workbench.go:67 func writeWorkbenchError(w http.ResponseWriter, status int, class, message string)`; `:112 func (d *Daemon) writeWorkbenchStoreError(w, r, ctx, err)`; `:125 func (d *Daemon) writeWorkbenchInternalError(w, r, err)`. H4/H5 mutants name real helpers with matching arity; H4's `404` is an untyped int constant accepted as `status int`. |
| V20 | Controller, round 3 (glm r2 fix): `rg -n -o '\{\{with \.Object\}\}\{\{range \.Edges\}\}\|\{\{with \.Commits\}\}\|\{\{if \.Truncated\}\}' host/workbench/render.go \| sort \| uniq -c` | line 178: `{{with .Object}}{{range .Edges}}` ×1, `{{with .Commits}}` ×1, `{{if .Truncated}}` ×2 (CommitView then ReferenceView). A1's anchor is unique. A2 must replace only the FIRST `{{if .Truncated}}` (inside `{{with .Commits}}`); the sprint plan must spell that anchor with its preceding context so it matches once. |
| V21 | Controller, round 3 (kimi r2 fix, supersedes V9's test run): `go test -race ./host/workbench ./host/daemon -run 'Workbench\|Render\|Grade' -v -count=1 \| grep -c '^=== RUN'` (backslash-pipe, as V9 logged it) vs the same with an unescaped pipe `-run 'Workbench|Render|Grade'` | escaped form: **0** tests run (RE2 reads `\|` as a literal pipe — V9's "ok" was vacuous); unescaped form: **176** `=== RUN` lines, both packages ok. Every acceptance command in §6 uses the unescaped form and must stay that way; a `\|` inside a Markdown table cell is table escaping, never shell text. |

### 10.1 Source instrument (negative results have positive controls)

```sh
python3 - <<'PY'
from pathlib import Path
s=Path('host/daemon/daemon.go').read_text().split('type readStore interface {')[1].split('\n}')[0]
print('readStore:',s)
for term in ['GetLogEntry(', 'GetLogEntryByHash(', 'LogEntryIndex(']: print(term,s.count(term))
s=Path('host/daemon/workbench.go').read_text()
print('Timeline fields reads From/Limit:',s.count('page.Timeline.From'),s.count('page.Timeline.Limit'),'CONTROL NextHref:',s.count('page.Timeline.NextHref'))
PY
```

### 10.2 Reproduce the daemon measurement safely

Create a fresh temporary copy (the measured directory was `/tmp/world205-probe`). The following copies tracked source without copying git metadata or writing git state. Execute each command from the indicated directory. All appended probe files belong only in the copy.

```sh
python3 - <<'PY'
import pathlib,subprocess,shutil
p=pathlib.Path('/tmp/world205-probe'); p.mkdir(exist_ok=True)
for f in subprocess.check_output(['git','ls-files'],text=True).splitlines():
 s=pathlib.Path(f)
 if s.is_file():
  d=p/f; d.parent.mkdir(parents=True,exist_ok=True); shutil.copy2(s,d)
print(p)
PY
```

Write this exact file as `/tmp/world205-probe/host/daemon/zz205_test.go`, then run V4 there:

```go
package daemon
import (
 "context"
 "testing"
 "strings"
 "bytes"
 "github.com/sunholo-data/ailang-world/host/hashref"
 "github.com/sunholo-data/ailang-world/host/store"
)
type root205 struct {readStore; fail hashref.HashRef; err error}
func(s root205)GetObject(ctx context.Context,r hashref.HashRef)(store.Object,bool,error){if r==s.fail{return store.Object{},false,s.err};return s.readStore.GetObject(ctx,r)}
func TestProbe205(t *testing.T){
 d:=newHandlerDaemon(t); ctx:=context.Background(); g:=seedGenesisEmbedded(t,d,"probe")
 c:=testCommit(g,0,"probe"); if err:=d.store.Commit(c);err!=nil{t.Fatal(err)}
 for _,w:=range []store.World{g,c.NextWorld}{_,ok,err:=d.store.GetObject(ctx,w.StateRoot); t.Logf("root revision=%d ref=%s stored=%v err=%v",w.Revision,w.StateRoot,ok,err)}
 _,ok,err:=d.store.GetObject(ctx,c.Objects[0].Hash); if !ok||err!=nil{t.Fatal("control",ok,err)}; t.Log("CONTROL committed object stored=true")
 for _,path:=range []string{"/workbench","/workbench?from=0&entry=0"}{
  r:=requestRecorder(t,d,"GET",path,nil);nav,ok:=workbenchRegion(r.Body.String(),`<nav aria-label="world browser">`,"</nav>");if !ok{t.Fatal("nav absent")}
  for _,m:=range workbenchLinkPattern.FindAllStringSubmatch(nav,-1){if strings.Contains(m[1],"?object="){got:=requestRecorder(t,d,"GET",m[1],nil);t.Logf("%s status=%d href=%s text=%s targetStatus=%d",path,r.Code,m[1],m[2],got.Code);if got.Code!=404{t.Fatal(got.Code)}}}
 }
 // Store the exact preimage used by testCommit; same world now has a stored root.
 obj:=testObject("unused");obj.Payload=[]byte("state-0-probe");obj.Hash=hashref.SumSHA256(obj.Payload)
 if obj.Hash!=c.NextWorld.StateRoot{t.Fatal("preimage")};if err:=d.store.PutObject(obj);err!=nil{t.Fatal(err)}
 _,ok,err=d.store.GetObject(ctx,obj.Hash);if !ok||err!=nil{t.Fatal(ok,err)}
 r:=requestRecorder(t,d,"GET","/workbench?object="+obj.Hash.String(),nil);t.Logf("stored-root CONTROL GetObject=%v targetStatus=%d",ok,r.Code);if r.Code!=200{t.Fatal(r.Code)}
 d.errLog=&bytes.Buffer{}
 for _,e:=range []error{errSentinelInternal,context.DeadlineExceeded}{d.reads=root205{d.store,c.NextWorld.StateRoot,e};edge,err:=d.checkedEdge(ctx,"stateRoot",c.NextWorld.StateRoot);if err==nil{t.Fatal("error lost")};t.Logf("checkedEdge error=%v available=%v",err,edge.Available)
 rec:=requestRecorder(t,d,"GET","/workbench",nil);t.Logf("base world check bypass error=%v status=%d",e,rec.Code)
 }
}
```

### 10.3 Run the extended actual lexical test

In the temporary copy, run this patch, then V5's first command. The two reported failures are expected measurements of the base defect.

```sh
python3 - <<'PY'
from pathlib import Path
p=Path('host/workbench/render_test.go');s=p.read_text()
s=s.replace('reflect.TypeOf(Page{})}', 'reflect.TypeOf(Page{}), reflect.TypeOf(WorldView{}), reflect.TypeOf(ObjectView{}), reflect.TypeOf(CommitView{}), reflect.TypeOf(ReferenceView{}), reflect.TypeOf(GradeView{}), reflect.TypeOf(EdgeView{})}')
p.write_text(s)
PY
```

### 10.4 Typed-vs-lexical experiment

Write this as `host/workbench/zz205_test.go` in the copy. It is a bounded proof of the collision on the current template, not the full fail-closed analyzer specified for production tests in §4.3. Run `AILANG_BIN="$HOME/.pinned-ailang/ailang" go test ./host/workbench -run '^TestProbe205Census$' -v -count=1`.

```go
package workbench
import("testing";"reflect";"strings";"sort";"text/template/parse";"html/template")
var types205=[]reflect.Type{reflect.TypeOf(Page{}),reflect.TypeOf(WorldView{}),reflect.TypeOf(TimelineView{}),reflect.TypeOf(EntryView{}),reflect.TypeOf(ObjectView{}),reflect.TypeOf(CommitView{}),reflect.TypeOf(ReferenceView{}),reflect.TypeOf(GradeView{}),reflect.TypeOf(EdgeView{})}
func census205(src string, strong bool)[]string{
 seen:=map[string]bool{}
 deref:=func(t reflect.Type)reflect.Type{for t!=nil&&(t.Kind()==reflect.Pointer||t.Kind()==reflect.Slice){t=t.Elem()};return t}
 tm:=template.Must(template.New("workbench").Funcs(template.FuncMap{"edgeUnavailable":edgeUnavailable,"workbenchHref":workbenchHref}).Parse(src))
 var eval func(parse.Node,reflect.Type)reflect.Type
 var walk func(parse.Node,reflect.Type)
 eval=func(n parse.Node,d reflect.Type)reflect.Type{
 switch x:=n.(type){
 case *parse.DotNode:return d
 case *parse.FieldNode:
  for _,name:=range x.Ident{d=deref(d);if d==nil||d.Kind()!=reflect.Struct{panic("unknown field "+name)};f,ok:=d.FieldByName(name);if !ok{panic(name)};seen[d.Name()+"."+name]=true;d=f.Type};return d
 case *parse.PipeNode:var r reflect.Type;for _,c:=range x.Cmds{r=eval(c,d)};return r
 case *parse.CommandNode:
  for _,a:=range x.Args{eval(a,d)}
  if f,ok:=x.Args[0].(*parse.IdentifierNode);ok&&f.Ident=="edgeUnavailable"{arg:=deref(eval(x.Args[1],d));if arg!=reflect.TypeOf(EdgeView{}){panic("bad edge argument")};seen["EdgeView.Available"]=true;return reflect.TypeOf(true)}
  return evalFirst205(x.Args[0],d,eval)
 };return nil
 }
 walk=func(n parse.Node,d reflect.Type){
 switch x:=n.(type){
 case *parse.ListNode:if x!=nil{for _,c:=range x.Nodes{walk(c,d)}}
 case *parse.ActionNode:eval(x.Pipe,d)
 case *parse.IfNode:eval(x.Pipe,d);walk(x.List,d);walk(x.ElseList,d)
 case *parse.WithNode:next:=eval(x.Pipe,d);walk(x.List,deref(next));walk(x.ElseList,d)
 case *parse.RangeNode:next:=eval(x.Pipe,d);walk(x.List,deref(next));walk(x.ElseList,d)
 case *parse.TemplateNode:next:=eval(x.Pipe,d);walk(tm.Lookup(x.Name).Tree.Root,deref(next))
 }
 }
 if strong{walk(tm.Tree.Root,reflect.TypeOf(Page{}))}
 var missing []string
 for _,typ:=range types205{for i:=0;i<typ.NumField();i++{name:=typ.Field(i).Name;key:=typ.Name()+"."+name;hit:=seen[key];if !strong{hit=strings.Contains(src,"."+name)};if !hit{missing=append(missing,key)}}};sort.Strings(missing);return missing
}
func evalFirst205(n parse.Node,d reflect.Type,f func(parse.Node,reflect.Type)reflect.Type)reflect.Type{switch n.(type){case *parse.FieldNode,*parse.DotNode:return f(n,d)};return nil}
func TestProbe205Census(t *testing.T){
 src:=pageHTML+partialsHTML
 for _,s:=range []bool{false,true}{t.Logf("strong=%v base missing=%v",s,census205(src,s))}
 mutated:=strings.ReplaceAll(src,"{{.World.Ref}}","literal-world")
 if mutated==src{t.Fatal("mutation absent")}
 for _,s:=range []bool{false,true}{t.Logf("strong=%v World.Ref mutation missing=%v",s,census205(mutated,s))}
 if strings.Contains(strings.Join(census205(mutated,false),","),"WorldView.Ref"){t.Fatal("weak unexpectedly caught")}
 if !strings.Contains(strings.Join(census205(mutated,true),","),"WorldView.Ref"){t.Fatal("strong missed")}
}
```

The sorted outputs were:

```text
strong=false base missing=[TimelineView.From TimelineView.Limit WorldView.LogHead WorldView.Revision]
strong=true base missing=[TimelineView.From TimelineView.Limit WorldView.LogHead WorldView.Revision]
strong=false World.Ref mutation missing=[TimelineView.From TimelineView.Limit WorldView.LogHead WorldView.Revision]
strong=true World.Ref mutation missing=[TimelineView.From TimelineView.Limit WorldView.LogHead WorldView.Ref WorldView.Revision]
```

### 10.5 Error-path prototype

In the copy only, replace the unconditional world-edge construction, rename the defect logger so it is not run against the changed handler, and append the error-path probe:

```sh
python3 - <<'PY'
from pathlib import Path
p=Path('host/daemon/workbench.go');s=p.read_text()
s=s.replace('page.World = workbench.WorldView{','stateRoot, err := d.checkedEdge(ctx, "stateRoot", world.StateRoot)\n\t\tif err != nil {\n\t\t\td.writeWorkbenchStoreError(w, r, ctx, err)\n\t\t\treturn\n\t\t}\n\t\tpage.World = workbench.WorldView{')
s=s.replace('workbench.EdgeView{Available: true, Target: world.StateRoot.String(), Href: "?object=" + world.StateRoot.String()}','stateRoot')
p.write_text(s)
p=Path('host/daemon/zz205_test.go');s=p.read_text().replace('func TestProbe205(', 'func SkipProbe205(');p.write_text(s)
PY
```

Append this to `host/daemon/zz205_test.go` in the copy, then run V6:

```go
func TestProbe205Checked(t *testing.T){
 d:=newHandlerDaemon(t);g:=seedGenesisEmbedded(t,d,"probe");c:=testCommit(g,0,"probe");if err:=d.store.Commit(c);err!=nil{t.Fatal(err)}
 d.errLog=&bytes.Buffer{}
 for _,x:=range []struct{err error;want int}{{nil,200},{errSentinelInternal,500},{context.DeadlineExceeded,503}}{
  d.reads=root205{d.store,c.NextWorld.StateRoot,x.err};if x.err==nil{d.reads=d.store}
  r:=requestRecorder(t,d,"GET","/workbench",nil);t.Logf("error=%v status=%d",x.err,r.Code);if r.Code!=x.want{t.Fatal(r.Code)}
 }
}
```

### 10.6 Measurement limits and final artifact check

All five requested design measurements were completed. A full final implementation, the future mutation table, and a new real-incident five-minute provenance walk were **not** executed: this iteration authorizes a design document only. V6 prototypes only the new daemon read/error path; V5 prototypes only the census experiment. Controller V13 additionally tested the combined handler/nav prototype across both packages and measured the two failures documented above. The corrected implementation and its mutation drills remain future acceptance gates. The focused baseline race suite and host vet passed (V9). No browser-based visual assessment was performed.

Round-1 final check is preserved in V10. Its temporary probe copy was removed with `shutil.rmtree("/tmp/world205-probe")`. Round 2 creates no throwaway probes and edits only this now-tracked design document. V12–V13 are attributed controller measurements, not newly executed tests. The controller owns commits.


### 10.7 V12 reproducible probe (sonnet r2 fix)

Controller-run at base `aa3e36a` in the detached probe worktree `.probe-world-iter205` (outside any sandbox), then deleted. Write as `host/daemon/zz205probe_test.go` in a scratch copy and run `AILANG_BIN="$HOME/.pinned-ailang/ailang" go test ./host/daemon -run '^TestProbe205PagingRoot$' -v -count=1`:

```go
package daemon

import (
	"context"
	"testing"

	"github.com/sunholo-data/ailang-world/host/workbench"
)

func TestProbe205PagingRoot(t *testing.T) {
	d := newHandlerDaemon(t)
	seedWorkbenchLog(t, d, workbench.WorkbenchPageLimit+5)
	ctx := context.Background()
	head, ok, err := d.store.SelectedHead(ctx)
	if err != nil || !ok {
		t.Fatalf("head ok=%v err=%v", ok, err)
	}
	w, ok, err := d.store.GetWorld(ctx, head)
	if err != nil || !ok {
		t.Fatalf("world ok=%v err=%v", ok, err)
	}
	_, stored, err := d.store.GetObject(ctx, w.StateRoot)
	t.Logf("PAGING head=%s revision=%d stateRoot=%s stored=%v err=%v", head, w.Revision, w.StateRoot, stored, err)
	e, ok, err := d.store.GetLogEntry(ctx, 0)
	if err != nil || !ok {
		t.Fatalf("control entry0 ok=%v err=%v", ok, err)
	}
	_, tstored, _ := d.store.GetObject(ctx, e.TransitionRef)
	t.Logf("CONTROL entry0 transitionRef stored=%v", tstored)
}
```

Literal output:

```text
zz205probe_test.go:23: PAGING head=sha256:a584866ed27514a33c00d6d96a31b49ffa22fb2ee8f05381851c464f5389ab2e revision=104 stateRoot=sha256:7b2475660d4b6f7093df919899cfcf15e5f281004d7b8cd9efe674868bed04f8 stored=false err=<nil>
zz205probe_test.go:29: CONTROL entry0 transitionRef stored=true
ok  	github.com/sunholo-data/ailang-world/host/daemon	0.492s
```


## 11. Quorum verification log

Round 1: **REJECT / BLOCKED 3/3** — oc-glm-5-3, oc-kimi-k3, claude-sonnet-5. All accepted the design direction; objections concerned completeness or measured premises.

- **oc-glm-5-3 — blast radius:** V13 measured two failures; §4.4 scopes the reference count to referencedBy and narrows the selected-entry error seam; §5, M1 and §9 now include both. This replaces the incorrect no-reference-walk-change claim.
- **oc-glm-5-3 — test placement and names:** V15 confirms package helpers; §4.4 separates direct render tests from daemon tests and declares exact test/subtest names used by every mutation row.
- **claude-sonnet-5 — Unavailable premise and coverage:** V11 confirms the existing field/action; §2 item 10 and §4.1 say only its default is new; R7 deletes the action and the census self-test inventory names WorldView.Unavailable.
- **oc-kimi-k3 — paging-root premise:** V12 measures the actual paging head root absent with a stored transition control; D5 and §4.4 cite V12 and require the additional stored-root world fixture.
- **oc-kimi-k3 — explicit world resolution:** V14 records the query branch and GetWorld(ref); §4.4 requires an older explicit world whose metadata differs from SelectedHead.

Round 2: **REJECT / BLOCKED 3/3** — same seats. Again no reviewer disputed the direction; every objection was a narrow completeness gap carrying a concrete `proposed_fix`. Per the narrow-refinement carve-out, the **controller** measured each and applied the reviewers' fixes verbatim (round 3 is this edit; no designer re-run, no further quorum):

- **oc-glm-5-3 / oc-kimi-k3 — base drift unbounded:** V18 shows `host/` unchanged between `aa3e36a` and `61ac489` (0 commits, 0 files; whole-range control 2 commits). Premises carry over.
- **oc-glm-5-3 — H4/H5 helpers and A1/A2 anchors unmeasured:** V19 quotes all three error-writer signatures; V20 quotes the anchors and records that `{{if .Truncated}}` occurs twice, so A2's anchor must include its `{{with .Commits}}` context.
- **oc-kimi-k3 — V9 possibly vacuous:** confirmed. V21 measured **0** tests for the `\|` form V9 logged and **176** for the unescaped form; V21 supersedes V9's test claim. §6 already uses the unescaped form.
- **claude-sonnet-5 — V12 not reproducible:** §10.7 now carries the exact probe source and its literal output.
