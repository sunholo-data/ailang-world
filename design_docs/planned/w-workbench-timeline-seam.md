# w-workbench-timeline-seam — make the timeline a paged, selectable surface whose entries lead to objects

- Status: **planned** (design only; nothing implemented) · Date: **2026-09-24** · Rows **35** and **38** (both clause-5), plus the unfiled same-class `Page.Selected` defect (§1 F3)
- Iteration: 184, designer role · Measurement base: `origin/dev` = `82e3630` (V0)
- Scope: ONE doc for BOTH rows, as row 38 asks ("both touch the same `TimelineView`/`EntryView` seam in `render.go`, so designing them together is likely cheaper than twice"). The two rows are also one defect class seen from both ends: view-model fields the handler writes and the template never reads (35, `Truncated`, `Selected`), and template actions that read fields the handler never writes (`NextHref`/`PrevHref`).
- Query grammar: **unchanged** (§4). `supportedWorkbenchQuery` stays byte-identical, so charter row 34's open hunks there are neither fixed nor moved.
- Estimate: ~1 day, ~270 LOC including tests, three milestones (§8). Go-host only; no `.ail` changes.
- Revision 1 (quorum r1, BLOCKED 3/3): I1 narrowed to the timeline and selected-entry regions, word for word as gpt6-astra asked. `emitted-links-resolve` now isolates those two DOM regions and requires ≥1 link in each category. The world-pane exclusion is measured (V15: the seeded world's `StateRoot` link returns 404, and R-c stays open). The gemini and glm objections were measured and refuted (V16, V17a/b). The design and code are unchanged. See §11.

---

## §1 Problem

The workbench (`design_docs/implemented/w-workbench-read-only.md`, "the original design" below) exists
so that a human can walk provenance in a browser: clause 5. Its timeline is the way into that walk.
At `82e3630` the timeline cannot do that job. Every fact below comes from the command in brackets, and
each command is listed again in the Verification Log (§10).

**F1 — no entry's transition edges render (row 35).** `EntryView` has three `EdgeView` fields,
`TransitionFn`, `Interpreter` and `TransitionRef` (`host/workbench/render.go:63-65`). `entryView()`
fills each one with a `Target` and an `?object=` `Href` (`host/daemon/workbench.go:107-109`). In
`pageHTML` each is used by **0** template actions, while the five scalar fields next to them are
used by 1–3 actions each (EntryIndex 1, EntryHash 3, PrevEntryHash 3, SemanticsEpoch 1, WrittenBy 1)
[V1: field census over the `pageHTML` literal]. `TransitionRef` is the only field that could put a
committed object's hash on the timeline, and the comment at `host/daemon/workbench_test.go:87-93`
records that the route cannot show one.

**F2 — pagination is broken in both directions (row 38).** `TimelineView.Truncated` is written once
(`host/daemon/workbench.go:265`) and read **0** times anywhere under `host/`, tests included. As a
control, `PayloadTruncated` has **3** references [V3]. `NextHref`/`PrevHref` are read by two template
actions (`render.go:140-141`, each inside `{{if}}`), and **0** files in `host/daemon` mention either
field. The control is **3** `page.Timeline` lines in `workbench.go` [V2]. So the timeline stops at
100 entries and never offers a next link.

**F3 — `Page.Selected` renders nowhere (not yet a queue row; controller-found this iteration,
re-measured here).** `grep -c Selected host/workbench/render.go` returns **1**, the struct field at
`:92`. As a same-file control, `.Timeline.` appears on **3** template lines [V1, V1c]. The handler
still performs `?entry=N` in full. It reads the entry, returns 404 when the entry is absent, and
builds `page.Selected` (`workbench.go:239-251`). Then it discards the result. This is F1's class
again, and it is where the original design put the edges: "Selecting entry N shows the measured
fields below … transitionFn ──> object inspector if object exists" (original design §2.7).

**F4 — the `Truncated` predicate is wrong as a "there is more" signal.** The loop at
`workbench.go:254-264` reads indices `from … from+limit-1` and stops early only when an entry is
missing. It never reads `from+limit`. So `len(Entries) == limit` is true for a log that holds exactly
`from+limit` entries, where nothing comes after. A next link built from that predicate would lead to
an empty page. The only `GetLogEntry` calls are at `:240` (the selection) and `:255` (the loop)
[V5b]. The has-more decision therefore needs a new read (§2a).

**F5 — the default page shows the *oldest* history, not the newest.** When the query is empty,
`from = 0` (`:136`), and the timeline counts upward from entry 0. On a log with more than 100 entries,
the default page shows entries 0–99, and every later entry needs a hand-built
`?from=N&entry=N` URL. §2.2 of the original design refuses `?from=N` on its own, and row 38 already
records this cost. A timeline limited to one page therefore does not show the head. It shows only
the oldest page.

**F6 — the store does not guarantee that an edge's target object exists.** `Store.Commit` checks
`TransitionFn`, `Interpreter` and `TransitionRef` only for ref syntax (`validateRef` →
`hashref.Parse`, `host/store/store.go:437-446, 874-890`). None of them has to be among the stored
objects. The repo's own `testCommit` stores only the `TransitionRef` object. `TransitionFn` and
`Interpreter` are hashes of strings that are never stored (`host/daemon/handlers_test.go:79-99`).
`entryView()` still sets `Available: true` on all three. If those fields were rendered as they
stand, two of the three links on every test entry would lead to a 404 page.

**F7 — `GetObject` reads the full payload.** Its SELECT includes `payload`
(`store.go:482-483`) [V6], and a commit body can be up to `maxCommitBytes = 8388608`
(`host/daemon/daemon.go:108`). Checking every edge target on a 100-row page could mean up to
3 × 100 full-payload reads. Checking only the selected entry means 3.

---

## §2 Decisions

### (a) Paged or head-only? **PAGED.**

Head-only would mean deleting `Truncated`, `NextHref`, `PrevHref` and the two template actions. It
fails clause 5 because of F5. The timeline counts up from entry 0, so "head-only" really means
"entries 0–99 only". On a log that grows with each commit, the entries a human most needs for "why
did X happen" questions become unreachable from the page after the 100th commit. Turning the
timeline into a real head view (newest first) would reverse the original design's §2.7 ordering and
give up the "entry N" addressing. That is a larger redesign than wiring two links that already
exist. Paging costs one extra point read per direction per request and needs no grammar change (b).

How "there is more" is decided (F4): **probe, don't infer.** The handler emits `NextHref` only when
`GetLogEntry(ctx, from+limit)` returns `ok`. It emits `PrevHref` only when `from > 0` and
`GetLogEntry(ctx, prev)` returns `ok`, where `prev = max(0, from−limit)`. Each link is emitted
because the handler read the entry it selects during this request, and the log is append-only. So
**every paging link the page emits resolves to 200** (up to a store error or deadline on the next
request). Reading `limit+1` entries would also work, but the probe is simpler to state and pin: it is
the exact read that the target page's `entry=` will repeat. The probe is not skipped when
`len(Entries) < limit`. The store does not enforce contiguous indices (`entry_index` is only the
primary key; `Commit` does not check it against the previous index, V7), so a gap earlier in the
page does not prove that `from+limit` is absent.

`Truncated` is **deleted**. A non-empty `NextHref` is now the "there is more" signal. Keeping the
bool as a second copy with no readers would keep the defect class this doc closes.

### (b) What query state may a paging link emit? **`?from=N&entry=N` — the existing pair. No grammar change.**

Among the grammar's four non-empty states (original design §2.2), `from`+`entry` is the only one
that carries `from`. A paging link therefore emits `?from=T&entry=T`. It pages to `T` and selects
entry `T`, the first entry of the target page. Because of (d), that selection is meaningful: the
page opens with its first entry's edges available to walk. One builder produces the one link shape
the page can emit:

```go
func pageHref(from, entry int64) string {
	return "?from=" + strconv.FormatInt(from, 10) + "&entry=" + strconv.FormatInt(entry, 10)
}
```

Rejected: accepting `?from=N` alone. It would change `supportedWorkbenchQuery`'s `len(query) == 1`
branch (`workbench.go:66-67`). That function holds three of row 34's open hunks (`:69`, `:72`,
`:75`), so changing it would drag this sprint into row 34's territory. It would also add a newly
accepted form that needs its own still-refused neighbour, and it would give nothing the existing
pair lacks once selection renders. `?from=N` alone stays refused 400 (the existing
`TestWorkbenchRefusalBranches/unsupported-combination` arm), and §6 AC7 checks that no emitted link
takes that form.

A paging link does **not** carry the current selection (the link always selects its target page's
first entry). Carrying it would give paging two rules, "keep the current entry if there is one,
otherwise the first", which means one more branch and one more mutation for a small UX gain. It is
listed out of scope (§9).

### (c) Render the three `EntryView` edges, or stop populating them? **Render them on the selected entry, checked for existence; delete the three named fields.**

- Timeline rows render **no** edges. Each row gets a `select entry N` link to `pageHref(from, N)`,
  so any entry's edges are one click away. Rendering edges on every row would force a choice
  between an unchecked `Available: true`, which leads to 404s (F6), and up to 300 full-payload
  existence reads per page (F7).
- The selected entry renders its three edges. Each edge's target is checked once with
  `d.reads.GetObject` (3 reads, all on the request's `readCtx` context). A stored target renders as
  a link to `?object=<ref>`. A missing target renders through the **existing** `edgeUnavailable`
  → `UNAVAILABLE: {{.Missing}}` path with `Missing = "object <ref> is not stored"`. The ref stays
  visible, and the human stays on a page they can navigate instead of landing on a 404. A store
  error on the check returns 500/503 through the existing `writeWorkbenchStoreError`. It is never
  shown as "not stored".
- The three named fields `TransitionFn`/`Interpreter`/`TransitionRef` are **deleted**. The edges go
  into `EntryView.Edges []EdgeView` (`render.go:67`), which today has **0** writers and **0** readers
  of its own [V4b], with `Relation` = `transitionFn` / `interpreter` / `transitionRef` (the original
  design's §2.7 labels). This reuses `EdgeView.Relation` and makes the entry edges and the object
  provenance walk one rendering path instead of two.
- The edge markup at `render.go:162` moves into a shared `{{define "edge"}}`. Both the provenance
  walk and the selected entry call it, so there is one `edgeUnavailable` path and
  `TestRenderUnavailableProvenanceEdge` pins it for both callers.

### (d) `Page.Selected`: **IN.**

`?entry=N` is required by the only grammar state that pages (b), and it already costs a read, a
possible 404 and a view-model build. Leaving it unrendered would ship (a)/(b) with every paging link
selecting an entry the page then hides. It renders as an `<article aria-label="selected entry">` at
the top of the Inspector section:

- heading `<h3>selected entry N</h3>`. The heading deliberately does not begin with `entry `, so
  `TestWorkbenchTimelineBound`'s `<h3>entry ` count is unaffected (§5);
- the same `<dl>` of scalar fields the rows render (shared `{{define "entryFields"}}`);
- its three edges (c).

---

## §3 Design

### 3.1 View model (`host/workbench/render.go`)

```go
type EntryView struct {
	EntryIndex     int64
	EntryHash      string
	PrevEntryHash  string
	SemanticsEpoch int64
	WrittenBy      string
	SelectHref     string     // NEW: "?from=F&entry=N" — rows only
	Edges          []EdgeView // now WRITTEN (selected entry only) and READ
}                             // DELETED: TransitionFn, Interpreter, TransitionRef (:63-65)

type TimelineView struct {
	From     int64
	Limit    int
	Entries  []EntryView
	NextHref string
	PrevHref string
}                             // DELETED: Truncated (:74)
```

`Page` is unchanged (`Selected *EntryView` stays, now rendered).

### 3.2 Template (`host/workbench/render.go`)

Add a second constant that contains only definitions, and parse it into the same template (text/template
allows only one of several `Parse` calls on a receiver to contain non-definition text):

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

The `"edge"` body is `:162`'s inner body, byte-for-byte. The `"entryFields"` body is `:144-146`,
byte-for-byte. The rendered row markup is therefore unchanged, and so are the substrings existing
tests assert on.

`pageHTML` changes (line anchors are at `82e3630`):

| Anchor | Before | After |
|---|---|---|
| `:140-141` | unchanged (`{{if .Timeline.PrevHref}}…previous…`, `{{if .Timeline.NextHref}}…next…`) | unchanged; the handler now writes them |
| `:143-147` | `<article><h3>entry {{.EntryIndex}}</h3><dl>` + the three inline `<dt>` lines + `</dl></article>` | `<article><h3>entry {{.EntryIndex}}</h3><dl>`⏎`{{template "entryFields" .}}`⏎`</dl><p><a href="{{workbenchHref .SelectHref}}">select entry {{.EntryIndex}}</a></p></article>` |
| after `:151` (`<h2>Inspector</h2>`) | — | `{{with .Selected}}<article aria-label="selected entry"><h3>selected entry {{.EntryIndex}}</h3><dl>`⏎`{{template "entryFields" .}}`⏎`</dl>{{range .Edges}}{{template "edge" .}}{{end}}</article>{{end}}` |
| `:162` | `{{with .Object}}{{range .Edges}}<p>…inline edge…</p>{{end}}{{end}}` | `{{with .Object}}{{range .Edges}}{{template "edge" .}}{{end}}{{end}}` |

The select link sits **outside** the `<h3>`, so the literal `<h3>entry 99</h3>` that
`TestWorkbenchTimelineBound` asserts stays byte-identical. html/template renders `&` in an `href`
as `&amp;`. This was measured on the host Go toolchain with a throwaway program outside the repo:
`?from=100&entry=100` renders `href="/workbench?from=100&amp;entry=100"` [V8]. Tests assert that
escaped form.

### 3.3 Handler (`host/daemon/workbench.go`)

1. **`entryView` (`:101-112`)** drops the three `EdgeView` lines (`:107-109`). It keeps the five
   scalar fields and has no store access.
2. **New `pageHref(from, entry int64) string`** (§2b), the only builder of `from`/`entry` query
   strings.
3. **New method** that checks each edge target (§2c):

   ```go
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

   (It adds a `host/hashref` import to `workbench.go`. `daemon.go:53` already imports that package,
   so the package's import closure does not change.)
4. **Selected block (`:249-250`)**: after `selected := entryView(entry)`, call
   `edges, err := d.entryEdges(ctx, entry)`. On error, `d.writeWorkbenchStoreError(w, r, ctx, err); return`.
   Then set `selected.Edges = edges` and `page.Selected = &selected`.
5. **Timeline loop (`:263`)**: `view := entryView(entry); view.SelectHref = pageHref(from, view.EntryIndex)`, then append `view`.
6. **Replace `:265`** (the `Truncated` write) with the two checks:

   ```go
   next := from + int64(limit) // cannot overflow: :150 already refused from > MaxInt64-limit
   if next <= math.MaxInt64-int64(limit) { // the link's own from must pass :150 on the next request
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

Read budget per request: previously at most 100 log reads, 1 selection read and the world/object
reads. Now at most 2 more log reads (the probes) and, only when `entry` is present, 3 object reads.
All of them go through `d.reads` on the `readCtx` context, so the elapsed-time contract the original
design inherited (§3.3 there) covers them unchanged, and `TestNoNewDeadlineFreeStoreReads` sees no
new `context.Background()`.

### 3.4 Invariants this sprint makes persistent

- **I1** Every timeline paging link, timeline row-selection link, and stored-target link in the
  selected-entry article resolves to 200 against an unchanged store, absent store errors or deadline
  expiry. Existing world-pane and object-pane links are outside this invariant; R-c remains open.
  Pinned by `TestWorkbenchTimelinePaging/emitted-links-resolve`. (The exclusion is measured, not
  assumed: when no `?world` is given, `handleWorkbench` resolves the world from the selected head, so the
  world pane renders on the paging pages too. On the seeded test world, its `StateRoot` link returns
  **404** [V15]. A whole-page href sweep would therefore red at base because of R-c, not because of
  this sprint.)
- **I2** A missing edge target is `UNAVAILABLE`, never a link, and a store failure on the check is
  5xx, never `UNAVAILABLE`.
- **I3** No `EntryView`, `TimelineView` or `Page` field goes unrendered unless it is on a named,
  self-checking exemption list. Pinned by `TestWorkbenchViewFieldsAllRender` (§6). This is the
  class ratchet, so row 35's class cannot quietly return.

---

## §4 Query grammar delta

**None.** `acceptedWorkbenchKeys` (`workbench.go:33-39`) and `supportedWorkbenchQuery` (`:62-76`)
stay byte-identical (§6 AC5). The four accepted non-empty states stay exactly the ones in §2.2 of
the original design. Each emitted link uses the existing `from`+`entry` state:

| Link | Emitted form | Grammar state |
|---|---|---|
| next | `?from=F+L&entry=F+L` | `from`+`entry` (existing) |
| previous | `?from=P&entry=P`, `P = max(0, F−L)` | `from`+`entry` (existing) |
| select entry N | `?from=F&entry=N` | `from`+`entry` (existing) |
| entry edge | `?object=<ref>` | `object` (existing) |

These forms stay refused, and so do their existing tests: `?from=0` alone (400
`unsupported-combination`), unknown keys, duplicate keys, a negative `from`, an overflowing `from`,
and an absent `entry` (404). No form is newly accepted, so no newly accepted form needs a refused
neighbour. The neighbour this sprint could plausibly produce by mistake is a link that emits `from`
without `entry`. `emitted-links-resolve` rules that out by fetching each link (a `?from=N` link
would return 400), and mutation H8 checks that it does.

---

## §5 Conflict Surface

| Touched | How | Consequence |
|---|---|---|
| **Row 34** (`w-workbench-unpinned-render-hunks`): hunks at `supportedWorkbenchQuery` `:69`, `:72`, `:75` and `workbenchHref` `render.go:101` | **Not touched.** Both functions stay byte-identical (AC5). The new tests exercise only the *accepting* `from`+`entry` path and only `?`-prefixed hrefs | Row 34 stays open with all seven hunks still unpinned. This sprint **does not absorb it** (standing rule 1) and does not add row 34's recommended `world`+`from` / `world`+`payload` negative tests. It **does not move row 34's line anchors either**: every change in `workbench.go` lands at `:101` and below |
| **Row 35**, **Row 38** | Discharged by M1–M3 | Close both rows when the sprint lands. F3 (`Selected`) is fixed in the same sprint and never needs its own row |
| `TestWorkbenchTimelineBound` (`workbench_test.go:297-330`) | Counts `<h3>entry ` on `/workbench` (= 100) and asserts `<h3>entry 99</h3>` | Unaffected: `/workbench` has no selection, the selected heading is `selected entry`, and the select link sits outside the `<h3>`. M3 reuses its 105-entry seeding through an extracted helper |
| `TestWorkbenchRendersSeededWorldAndTimeline` (`:75-131`) | Sole-source assertion on `<dt>entry hash</dt><dd><span class="hash" title="HASH"` | Still one source per entry on `/workbench` (no selection). Its comment at `:87-93` ("the WB.B template renders no TransitionRef action at all") becomes false. M2 rewrites it to point at `TestWorkbenchSelectedEntry` |
| `TestWorkbenchRefusalBranches` (`:133-234`) | Byte-identical (AC5) | The grammar's tests are unchanged |
| `TestRenderUnavailableProvenanceEdge`, `TestRenderEscapesAllObjectText`, `TestRenderEmitsOnlyLocalLinks` (`render_test.go`) | Now run through the shared `"edge"` define | Must pass unmodified. `TestRenderEscapesAllObjectText` builds `EntryView{WrittenBy: hostile}`, which still compiles (no removed field is set) |
| Read-deadline tests (`read_deadline_test.go`, route table entry `workbench` `:502`; `blockingStore`, `recordingStore`) | 1 extra probe read on `/workbench` (`from=0` skips the prev probe) | The blocking store blocks at the first `GetLogEntry` anyway. `recordingStore` only notes the context. Both are expected to pass unchanged, and AC2 covers them |
| Original design §6 M12 (page cap 100→101 killed by `TimelineBound`) and M16 (`edgeUnavailable` neuter) | M16's anchor moves from `:162` to the `"edge"` define | Same function and same killer. The anchor move is recorded here so the next drill finds it |
| `bench_test.go` | Uses `store.LogHeader.TransitionFn` etc., **not** `EntryView` (V9) | Unaffected |

Files the implementation changes: `host/workbench/render.go`, `host/workbench/render_test.go`,
`host/daemon/workbench.go`, `host/daemon/workbench_test.go`. Nothing else.

---

## §6 Acceptance criteria

Each item is a command or a named test with its exact assertion. `export AILANG_BIN=~/.pinned-ailang/ailang` (v0.41.0, V10) applies to every command.

- **AC1** `go vet ./...` → rc=0.
- **AC2** `go test ./... -count=1` → rc=0.
- **AC3** `./scripts/verify_ail.sh` → rc=0, and `git diff --name-only origin/dev -- '*.ail' | wc -l` → `0`.
- **AC4 — new tests exist and pass.** `go test ./host/workbench ./host/daemon -count=1 -v -run 'TestRenderSelectedEntry|TestRenderTimelineRowSelectLink|TestRenderTimelinePagingLinks|TestWorkbenchViewFieldsAllRender|TestWorkbenchSelectedEntry|TestWorkbenchTimelinePaging|TestWorkbenchNextLinkOverflowGuard'` → rc=0, a `--- PASS` line for each of the 7 top-level names, and 0 `--- FAIL`. Assertions:
  - `TestRenderSelectedEntry` (render pkg): `Page{Selected: &EntryView{EntryIndex: 7, EntryHash: "SEL-HASH-7", Edges: []EdgeView{{Relation: "transitionRef", Available: true, Target: "abc", Href: "?object=abc"}, {Relation: "interpreter", Target: "def", Missing: "object def is not stored"}}}}`. The substring from `aria-label="selected entry"` to the next `</article>` must contain `<h3>selected entry 7</h3>`, `title="SEL-HASH-7"`, `transitionRef: <a href="/workbench?object=abc"` and `interpreter: <span class="unavailable" role="note">UNAVAILABLE: object def is not stored</span>`. Control arm: with `Selected: nil` the body contains no `aria-label="selected entry"`.
  - `TestRenderTimelineRowSelectLink`: a row `{EntryIndex: 3, SelectHref: "?from=0&entry=3"}` renders `<a href="/workbench?from=0&amp;entry=3">select entry 3</a>`.
  - `TestRenderTimelinePagingLinks`: `/prev` — `PrevHref: "?from=0&entry=0"` renders `<a href="/workbench?from=0&amp;entry=0">previous</a>`; `/next` — `NextHref: "?from=100&entry=100"` renders `<a href="/workbench?from=100&amp;entry=100">next</a>`; `/neither` — the zero `TimelineView` renders neither `>previous</a>` nor `>next</a>`.
  - `TestWorkbenchViewFieldsAllRender` (render pkg): for every field of `EntryView`, `TimelineView` and `Page` (by `reflect`), the text `pageHTML+partialsHTML` contains `"."+field.Name`, except the fields listed in `unrenderedExempt = {"TimelineView.From", "TimelineView.Limit"}` (row candidate R-a, §9). **Self-checking arm:** each exempt field must *not* appear, so rendering one fails until its exemption is removed. The check is lexical and name-based: `.Edges` is satisfied by either `range`, and one field name can satisfy another struct's field. It is the class ratchet, not the killer of any specific mutation. The mutation table names specific killers.
  - `TestWorkbenchSelectedEntry` (daemon): seed genesis plus one `testCommit`. Positive controls measured inside the test: `GetObject(TransitionRef)` ok=true, `GetObject(TransitionFn)` ok=false, `GetObject(Interpreter)` ok=false. Subtests:
    `/stored-edge-links` — in `GET /workbench?from=0&entry=0` → 200, the selected-article substring contains `transitionRef: <a href="/workbench?object=<TransitionRef>"`;
    `/unstored-edge-unavailable` — the same substring contains `transitionFn: <span class="unavailable" role="note">UNAVAILABLE: object <TransitionFn> is not stored</span>` and the same line for `interpreter`, and contains no `href="/workbench?object=<TransitionFn>"`;
    `/row-select-link` — `GET /workbench` contains `href="/workbench?from=0&amp;entry=0">select entry 0</a>` and no `aria-label="selected entry"`;
    `/object-store-error` — with `d.reads = objectFailingStore{readStore: d.store}` (overrides `GetObject` only, returns an error), `GET /workbench?from=0&entry=0` → 500 and the body contains `>Internal<`. **Control:** the same store serves `GET /workbench` with 200, which proves only the edge check fails.
  - `TestWorkbenchTimelinePaging` (daemon): seed 105 entries with the helper extracted from `TestWorkbenchTimelineBound` (`seedWorkbenchLog(t, d, n)`). Subtests:
    `/head-page-has-next` — `GET /workbench` contains `href="/workbench?from=100&amp;entry=100">next</a>` and no `>previous</a>`;
    `/last-page-has-prev-no-next` — `GET /workbench?from=100&entry=100`: `strings.Count(body, "<h3>entry ")` = 5, contains `href="/workbench?from=0&amp;entry=0">previous</a>`, contains `href="/workbench?from=100&amp;entry=104">select entry 104</a>`, and no `>next</a>`;
    `/exactly-limit-no-next` — `GET /workbench?from=5&entry=5`: **control** `strings.Count(body, "<h3>entry ")` = 100 (so the old predicate `len == limit` holds), no `>next</a>` (F4), and contains `href="/workbench?from=0&amp;entry=0">previous</a>` (clamp 5−100 → 0);
    `/emitted-links-resolve` (pins I1) — run on each of the three pages above. **Region isolation:** (i) the *timeline region* is the substring from `<section aria-label="timeline">` to the first `</section>` after it (the timeline section has nested `<article>`s but no nested `<section>`); (ii) the *selected-entry region* is the substring from `<article aria-label="selected entry">` to the first `</article>` after it (the article has nothing nested inside it). If the timeline start marker is missing, the test fails. On `/workbench` the selected-entry marker must be absent, and on the two `entry=`-bearing pages it must be present. Nothing outside these two regions is scanned, so the `<nav aria-label="world browser">` world link and the inspector/provenance-walk sections are excluded (I1; V15). **Classification, not filtering:** every `<a href="(/workbench\?[^"]*)"[^>]*>([^<]*)</a>` match inside a region is put in exactly one category: *paging* (timeline region, text `previous` or `next`), *select* (timeline region, text `select entry N`), or *stored-edge* (selected-entry region, href `/workbench?object=…`). A match that fits no category fails the test, and none is dropped. Each classified href has `&amp;` replaced with `&` and is fetched against the same unchanged store. Every status must be 200, and the failure message names the category and href. **Control (per category, over the three pages):** paging ≥1 (`/workbench` has `next`, the other two have `previous`), select ≥1 (every page has rows), stored-edge ≥1 (`?from=100&entry=100` and `?from=5&entry=5` each select a `testCommit` entry whose `TransitionRef` object is stored, V7b). A category with 0 links fails the test, so a region extractor that matches nothing cannot pass vacuously.
  - `TestWorkbenchNextLinkOverflowGuard` (daemon): `d.reads = denseLogStore{readStore: d.store}`, whose `GetLogEntry` returns a copy of stored entry 0 with `EntryIndex = index` for every `index ≥ 0`. `GET /workbench?from=9223372036854775707&entry=9223372036854775707` (= MaxInt64−100, the largest `from` that `:150` accepts) → 200 and no `>next</a>`. **Control:** `?from=9223372036854775607&entry=9223372036854775607` → 200 and contains `href="/workbench?from=9223372036854775707&amp;entry=9223372036854775707">next</a>`.
- **AC5 — the grammar is byte-identical.** For each of `supportedWorkbenchQuery`, `acceptedWorkbenchKeys` and `TestWorkbenchRefusalBranches`: `diff <(git show origin/dev:host/daemon/workbench.go | sed -n '/^func supportedWorkbenchQuery/,/^}/p') <(sed -n '/^func supportedWorkbenchQuery/,/^}/p' host/daemon/workbench.go)` → rc=0, and likewise with `/^var acceptedWorkbenchKeys/,/^}/` on `workbench.go` and `/^func TestWorkbenchRefusalBranches/,/^}/` on `workbench_test.go`.
- **AC6 — the dead fields are gone and the dead reads are fed.** `grep -rnw --include='*.go' Truncated host/workbench host/daemon | wc -l` → `0` (control: `grep -rn --include='*.go' PayloadTruncated host | wc -l` ≥ `3`); `grep -cE '^\s+(TransitionFn|Interpreter|TransitionRef)\s+EdgeView' host/workbench/render.go` → `0` (control: `grep -cE '^\s+Edges\s+\[\]EdgeView' host/workbench/render.go` → `2`); `grep -cE 'Timeline\.(NextHref|PrevHref) = ' host/daemon/workbench.go` → `2`; `grep -c 'Selected' host/workbench/render.go` ≥ `2`.
- **AC7 — no deadline-free read and no bare-`from` link.** `grep -c 'context.Background' host/daemon/workbench.go` → `0` (control: `grep -c 'readCtx' host/daemon/workbench.go` → `1`); `grep -c '"?from=" + strconv' host/daemon/workbench.go` → `1` (only `pageHref` builds a `from` query).
- **AC8 — mutation drill.** Apply each row of §7 on its own. `go test ./host/workbench ./host/daemon -count=1` reds, and the named test is among the reds. Restore with `git diff --quiet -- host/` rc=0 before the next row. Record the per-row table (mutant landed y/n, red set, restored) in the sprint's Verification Log.
- **AC9 — size.** `git diff --shortstat origin/dev -- host/` shows insertions + deletions ≤ 320.

---

## §7 Mutation table

Every mutation is anchored to the diff this sprint ships, not to the rows' defects. Each one
compiles. "Killer" names the one test that must red. Other tests may also red.

| ID | File · anchor (post-sprint) | Mutation | Killer |
|---|---|---|---|
| R1 | render.go · Selected block | delete `{{with .Selected}}…{{end}}` | `TestRenderSelectedEntry` |
| R2 | render.go · Selected block | delete `{{range .Edges}}{{template "edge" .}}{{end}}` | `TestRenderSelectedEntry` |
| R3 | render.go · Selected block | delete `{{template "entryFields" .}}` inside Selected | `TestRenderSelectedEntry` (asserts `title="SEL-HASH-7"` inside the article) |
| R4 | render.go · `"edge"` define | `{{if edgeUnavailable .}}` → `{{if false}}` | `TestRenderUnavailableProvenanceEdge` |
| R5 | render.go · provenance section | `{{template "edge" .}}` → empty in `{{with .Object}}{{range .Edges}}` | `TestRenderUnavailableProvenanceEdge` |
| R6 | render.go · row | delete the `<p><a …>select entry …</a></p>` | `TestRenderTimelineRowSelectLink` |
| R7 | render.go · row | delete `{{template "entryFields" .}}` inside the row | `TestWorkbenchRendersSeededWorldAndTimeline` (existing entry-hash sole-source assertion) |
| R8 | render.go · `:141` | delete the `{{if .Timeline.NextHref}}…{{end}}` line | `TestRenderTimelinePagingLinks/next` |
| R9 | render.go · `:140` | delete the `{{if .Timeline.PrevHref}}…{{end}}` line | `TestRenderTimelinePagingLinks/prev` |
| R10 | render.go · `EntryView` | add a field `Extra string` that nothing renders | `TestWorkbenchViewFieldsAllRender` |
| R11 | render.go · `TimelineView` | render `{{.Timeline.From}}` in the h2 (exempt field now rendered) | `TestWorkbenchViewFieldsAllRender` (self-checking arm) |
| H1 | workbench.go · next check | delete `page.Timeline.NextHref = pageHref(next, next)` | `TestWorkbenchTimelinePaging/head-page-has-next` |
| H2 | workbench.go · next check | `if ok {` → `if !ok {` | `TestWorkbenchTimelinePaging/exactly-limit-no-next` |
| H3 | workbench.go · next check | replace the check with `if len(page.Timeline.Entries) == limit {` (the old predicate) | `TestWorkbenchTimelinePaging/exactly-limit-no-next` |
| H4 | workbench.go · next check | `if next <= math.MaxInt64-int64(limit) {` → `if true {` | `TestWorkbenchNextLinkOverflowGuard` |
| H5 | workbench.go · prev check | delete `page.Timeline.PrevHref = pageHref(prev, prev)` | `TestWorkbenchTimelinePaging/last-page-has-prev-no-next` |
| H6 | workbench.go · prev check | delete `if prev < 0 { prev = 0 }` | `TestWorkbenchTimelinePaging/exactly-limit-no-next`. The killer input is `from=5` < `limit=100`, a state the grammar accepts because `from` is any non-negative integer parsed at `:136-147` [V17b]. Unclamped, `prev = −95` and `GetLogEntry(−95)` returns `ok=false`, so no previous link renders. The test expects `href="/workbench?from=0&amp;entry=0">previous</a>` |
| H7 | workbench.go · prev check | `if from > 0 {` → `if true {` | `TestWorkbenchTimelinePaging/head-page-has-next` (asserts no previous) |
| H8 | workbench.go · `pageHref` | return `"?from=" + strconv.FormatInt(from, 10)` (drop `&entry=`) | `TestWorkbenchTimelinePaging/emitted-links-resolve` (400 on fetch) |
| H9 | workbench.go · `entryEdges` | delete the `if !ok { …; continue }` block | `TestWorkbenchSelectedEntry/unstored-edge-unavailable` |
| H10 | workbench.go · `entryEdges` | `if err != nil {` → `if false && err != nil {` | `TestWorkbenchSelectedEntry/object-store-error` |
| H11 | workbench.go · Selected block | delete `selected.Edges = edges` | `TestWorkbenchSelectedEntry/stored-edge-links` |
| H12 | workbench.go · Selected block | `if err != nil { d.writeWorkbenchStoreError(…); return }` → `if false && err != nil {` | `TestWorkbenchSelectedEntry/object-store-error` |
| H13 | workbench.go · timeline loop | delete `view.SelectHref = pageHref(from, view.EntryIndex)` | `TestWorkbenchSelectedEntry/row-select-link` |
| H14 | workbench.go · timeline loop | `pageHref(from, …)` → `pageHref(view.EntryIndex, …)` | `TestWorkbenchTimelinePaging/last-page-has-prev-no-next` (expects `from=100&amp;entry=104`) |
| H15 | workbench.go · next check | `pageHref(next, next)` → `pageHref(next, from)` | `TestWorkbenchTimelinePaging/head-page-has-next` (exact href) |

H10 and H12 have the same killer by design. H10 turns a store error into `UNAVAILABLE`, and H12
makes the handler ignore the returned error and render a nil `Edges`. Each produces a 200 where
the test requires 500.

---

## §8 Milestones

Each milestone ends green on AC1 and AC2. Deleting a field breaks compilation in both packages, so
M1 includes the handler's mechanical edits.

| M | Content | ~LOC | Exit |
|---|---|---|---|
| **M1 — the seam** | `render.go`: delete the three `EdgeView` fields and `Truncated`; add `SelectHref`; add `partialsHTML` (`"edge"`, `"entryFields"`); row → `entryFields` + select link; Selected block; provenance → `"edge"`. `workbench.go`: `entryView` drops `:107-109`; delete `:265`. Tests: `TestRenderSelectedEntry`, `TestRenderTimelineRowSelectLink`, `TestRenderTimelinePagingLinks`, `TestWorkbenchViewFieldsAllRender` | ~110 | AC1, AC2; R1–R11 drilled |
| **M2 — selection edges** | `workbench.go`: `entryEdges`, the Selected wiring, `pageHref`, `SelectHref` in the loop; rewrite the stale `:87-93` comment. Tests: `TestWorkbenchSelectedEntry` + `objectFailingStore` | ~80 | AC1, AC2; H9–H14 drilled |
| **M3 — paging** | `workbench.go`: the next/prev checks. Tests: `seedWorkbenchLog` extracted from `TestWorkbenchTimelineBound`, `TestWorkbenchTimelinePaging`, `TestWorkbenchNextLinkOverflowGuard` + `denseLogStore` | ~90 | AC1–AC9; H1–H8, H15 drilled |

---

## §9 Out of scope

These are recorded so that nobody widens this sprint to cover them. Each "row candidate" is
measured and is for the controller to file or drop. This doc does not file them.

- **Row 34**: all seven unpinned hunks and its recommended negative tests (§5).
- **R-a (row candidate) — `TimelineView.From`/`Limit`, `WorldView.Revision`/`LogHead` render
  nowhere.** Template counts are **0** for each (V1). They are the same class as row 35. The first
  two are the explicit exemptions in `TestWorkbenchViewFieldsAllRender`. `WorldView` and
  `ObjectView` are outside that test's struct set in this sprint.
- **R-b (row candidate) — `ObjectView.Edges` has 0 writers in `host/daemon`** [V4]. The
  "Provenance walk" section's `{{range .Edges}}` therefore renders nothing in production. This is
  row 38's other half of the class: a read with no writer. The original design's §2.7 explains why
  the data does not exist yet (there is no typed provenance projection). The defect is that the
  section shows *blank* instead of an explicit `UNAVAILABLE`. After this sprint the selected
  entry's edges are the page's only working provenance hop.
- **R-c (row candidate) — the world's `StateRoot` edge is `Available: true` without checking the
  target** (`workbench.go:201`). This is the F6 pattern on the world pane. `testCommit`'s state
  roots are unstored hashes, and the link returns 404 on the paging pages themselves [V15]. That is
  why I1 and `emitted-links-resolve` exclude the world pane explicitly (revision 1).
- Carrying the current selection across paging (§2b). A newest-first/head-anchored timeline (§2a).
  The timeline stopping at the first index gap in a non-contiguous log (existing behaviour, original
  design §3.2). Combining `world=` with timeline paging (paging links carry no `world`, and the
  timeline reads the store's single log whatever the `world`).
- Any `.ail` change, any `/v1` route change, any `readStore` seam change (the check uses the
  existing `GetObject`/`GetLogEntry`).

---

## §10 Verification Log

All commands were run in this worktree at `82e3630` on 2026-09-24. Template counts are over the
`pageHTML` literal, extracted with `sed -n '/^const pageHTML/,/^<\/html>`/p' host/workbench/render.go`
(written `$T` below). Every zero has a same-scope positive control in the same call.

| # | Claim | Command | Result |
|---|---|---|---|
| V0 | Base is `origin/dev` `82e3630` | `git rev-parse origin/dev HEAD` | both `82e36306d32f…` |
| V1 | Field census in the template | `for f in …; printf '%s' "$T" \| grep -o "\.$f\b" \| wc -l` | TransitionFn 0, Interpreter 0, TransitionRef 0, Timeline.Truncated 0, Selected 0, Timeline.From 0, Timeline.Limit 0, World.Revision 0, World.LogHead 0; controls EntryIndex 1, EntryHash 3, PrevEntryHash 3, SemanticsEpoch 1, WrittenBy 1, Timeline.NextHref 2, Timeline.PrevHref 2, PayloadTruncated 1 |
| V1b | `Selected` appears once in render.go (the struct field) | `grep -c Selected host/workbench/render.go` | `1` |
| V1c | Same-file control: template lines using `.Timeline.` | `sed -n '/^const pageHTML/,/^<\/html>`/p' … \| grep -c '\.Timeline\.'` | `3` |
| V2 | No daemon file writes NextHref/PrevHref | `grep -rl --include='*.go' -E 'NextHref\|PrevHref' host/daemon \| wc -l`; control `grep -c 'page\.Timeline' host/daemon/workbench.go` | `0` files; control `3` |
| V3 | `Truncated` read by nothing, tests included | `grep -rn --include='*_test.go' '\.Truncated\b' host \| wc -l`; `grep -rn --include='*.go' -E 'Truncated' host \| grep -v PayloadTruncated` | `0`; the only non-evidence hits are `render.go:74` (decl) and `workbench.go:265` (write) |
| V3c | Control for V3 | `grep -rn --include='*.go' PayloadTruncated host \| wc -l` | `3` |
| V4 | `ObjectView.Edges` has no daemon writer | `grep -rn --include='*.go' Edges host/daemon \| wc -l`; control `grep -c 'StateRoot: workbench.EdgeView' host/daemon/workbench.go` | `0`; control `1` |
| V4b | `EntryView.Edges` has no writer, and its only template read is the object `range` | `grep -rn --include='*.go' -E '\.Edges\|Edges:' host` | writes only in `render_test.go` `ObjectView` literals; template `.Edges` count `1` (V1) is `{{with .Object}}{{range .Edges}}` at `:162` |
| V5 | `entryView` is built at two sites (selection, loop) | `grep -c 'entryView(entry)' host/daemon/workbench.go` | `2` |
| V5b | The loop never reads `from+limit` | `grep -n 'GetLogEntry\|GetObject' host/daemon/workbench.go`; `sed -n 254,265p` | calls at `:213` (object), `:240` (selection), `:255` (loop `offset < limit`); none at `from+limit` |
| V6 | `GetObject` selects the payload | `grep -c 'SELECT interface_hash_ref, semantic_id, provenance, payload' host/store/store.go` | `1` (`:483`); `maxCommitBytes = 8388608` at `daemon.go:108` |
| V7 | Commit checks edge refs only for syntax and does not enforce index contiguity | `sed -n 437,446p;873,945p host/store/store.go`; `grep -n 'entry_index' host/store/*.go \| grep -i primary` | `validateRef` → `hashref.Parse` only; `entry_index INTEGER PRIMARY KEY` in schema fixtures; the append guard compares `ObservedHead` only |
| V7b | `testCommit` stores only the `TransitionRef` object | `sed -n 79,99p host/daemon/handlers_test.go` | `Objects: []store.Object{object}`, `TransitionRef: object.Hash`; `TransitionFn`/`Interpreter` are `SumSHA256("fn-"+label)` / `("interpreter")`, never stored |
| V8 | html/template renders `&` in href as `&amp;` | throwaway `go run` in `/tmp` (deleted), template `<a href="{{h .}}">` on `?from=100&entry=100` | `<a href="/workbench?from=100&amp;entry=100">next</a>` |
| V9 | No test builds `EntryView`'s deleted fields | `grep -rn --include='*.go' EntryView host cmd \| grep -v render.go:` | only `render_test.go:102` (`EntryView{WrittenBy: …}`) and a comment at `workbench_test.go:88`; `bench_test.go` hits are `store.LogHeader` fields |
| V10 | Pinned binary | `~/.pinned-ailang/ailang --version` | `AILANG v0.41.0` |
| V11 | Base is green on the touched packages | `go vet ./host/workbench ./host/daemon`; `AILANG_BIN=~/.pinned-ailang/ailang go test ./host/workbench ./host/daemon -run 'Workbench\|Render\|Grade' -count=1` | vet rc=0; both `ok` (47 RUN lines with `-v`) |
| V12 | `from` alone is refused today and tested | `grep -n 'unsupported-combination' host/daemon/workbench_test.go` | `:156` `"/workbench?from=0"` → 400 |
| V13 | Row 34's grammar hunks sit at `:69`/`:72`/`:75`, above every line this sprint edits | `sed -n 62,76p host/daemon/workbench.go` | confirmed; the first edited line is `:101` (`entryView`) |
| V14 | `failingStore` fails every method, so it cannot isolate the edge check | `grep -n 'func (failingStore)' host/daemon/read_deadline_test.go` | 5 methods overridden (`:697-713`), hence the new `objectFailingStore` |
| V15 | (r1) The world pane's `StateRoot` link on the paging pages returns 404, so I1 excludes it (R-c) | throwaway `host/daemon/zz_probe_r1_test.go` (run, then **deleted**; not in the diff): `newHandlerDaemon` + `seedGenesisEmbedded(…, "probe")` + one `testCommit(genesis, 0, "probe")`. For `GET /workbench` and `GET /workbench?from=0&entry=0`, extract the `<a href="/workbench?object=…"` inside `<nav aria-label="world browser">`, call `d.store.GetObject` on its ref, then GET the href. `AILANG_BIN=~/.pinned-ailang/ailang go test ./host/daemon -run TestZZProbeStateRoot -count=1 -v` | on both pages: page status `200`, world `sha256:dde9258f…57e3` (the head, resolved with no `?world`), href `/workbench?object=sha256:ad312da4…f170` (= `commit.NextWorld.StateRoot`, `isCommitStateRoot=true`), `GetObject` `ok=false err=<nil>`, GET status **`404`**. The probe test PASSed only as a logger, then was removed (`git status --short` clean). Run at `919a10b`, whose `host/` tree is identical to `82e3630` (`git diff --stat 82e3630 919a10b -- host` is empty) |
| V16 | (r1, gemini) `entryEdges`' ref fields are `hashref.HashRef` on the store types it reads. The `string` fields cited are a JSON wire struct | `sed -n 112,128p host/store/store.go`; `sed -n 149,157p host/store/journal.go` | `LogHeader`: `TransitionFn   hashref.HashRef` (`:115`), `Interpreter    hashref.HashRef` (`:116`); `LogEntry`: `TransitionRef hashref.HashRef` (`:128`). The `string`-typed `TransitionFn`/`TransitionRef`/`Interpreter` at `journal.go:155-157` are fields of `type intentWire struct` (`:149`, JournalIntent's JSON wire form), not `store.LogEntry` |
| V17a | (r1, glm) `"math"` is already imported and used in `workbench.go` | `sed -n 3,13p host/daemon/workbench.go`; `grep -n math host/daemon/workbench.go` | import block line `:7` `"math"`; use at `:150` `if from > math.MaxInt64-int64(limit) {` |
| V17b | (r1, glm) The prev clamp is reachable: `0 < from < limit` is an accepted state | `sed -n 136,150p host/daemon/workbench.go` | `from` is `strconv.ParseInt` of any `?from=` text, refused only if `< 0` (`:136-147`) or `> MaxInt64−limit` (`:150`). With `limit = 100`, `?from=5&entry=5` passes both, so `prev = 5−100 < 0` is reached. It is AC4's `/exactly-limit-no-next` input and H6's killer |

---

## §11 Quorum log

| Round | Reviewer | Verdict | Objection (one line) | Disposition |
|---|---|---|---|---|
| r1 | gpt6-astra | reject | I1 ("every href on a 200 page resolves") and `emitted-links-resolve` contradict the deferred R-c. The world-pane `StateRoot` link is emitted `Available: true` without a check and renders on the paging pages | **Upheld; applied verbatim.** I1 replaced with the reviewer's text (§3.4). `emitted-links-resolve` now scans only the timeline section and the selected-entry article, and requires ≥1 link per category (AC4). The exclusion is substantiated by V15 (404) |
| r1 | gemini-3-1-pro | reject | `entryEdges` assigns `string` fields to a `hashref.HashRef` | **Refuted, V16.** gemini r1 objection refuted: `store.LogHeader.TransitionFn`/`.Interpreter` and `store.LogEntry.TransitionRef` are `hashref.HashRef` (`store.go:115,116,128`). The `string` fields cited are `intentWire`, the journal's JSON wire struct (`journal.go:149-157`). No change |
| r1 | oc-glm-5-2 | reject | (a) `math` is not imported for the next-probe guard; (b) the prev clamp `prev < 0 → 0` is unreachable | **Refuted on both halves, V17a/V17b.** glm r1 objection refuted: (a) `"math"` is imported at `workbench.go:7` and used at `:150`; (b) `from` is any non-negative parsed integer (`:136-147`), so `?from=5&entry=5` (from=5 < limit=100) is accepted, reaches the clamp, and is H6's killer (`/exactly-limit-no-next`). No change |
