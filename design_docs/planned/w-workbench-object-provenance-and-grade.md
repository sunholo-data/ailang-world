# w-workbench-object-provenance-and-grade — the object page's provenance walk and grade are never supplied

- Status: **planned** · Date: **2026-09-24** · Rows **98** and **102** (both clause-5)
- Iteration: 186, designer role · Measurement base: `origin/dev` = `6127ff3` (V0)
- Scope: ONE doc for BOTH rows. They are the same defect on the same page, the object inspector at
  `GET /workbench?object=<hash>`: the handler builds `ObjectView` without writing two fields the
  template reads (`Edges`, row 98; `Grade`, row 102). Their zero values then render as a silent
  blank and as an empty reason.
- Query grammar: **unchanged** (§4).
- Estimate: ~0.5 day, **289 changed lines** (277 insertions + 12 deletions) measured on the design-time
  prototype (V15): 62 production lines, 227 test lines. Two milestones (§8). Go host only; no `.ail`
  change.
- S3 ("why is this not a package?"): this is the host HTML adapter's projection of store rows. It
  adds no policy. The grade policy stays in AILANG (`gradeOf`, `world/types.ail`), and this design
  deliberately does **not** transcribe it (§2b).

---

## §1 Problem

Every fact below comes from a command logged in §10. The live-page facts come from a throwaway
probe test that called the daemon handler at `6127ff3` (V9). The probe was deleted afterwards.

**F1 — the provenance walk is blank on every page (row 98).** `handleWorkbench` builds
`page.Object` at `host/daemon/workbench.go:264-268` with seven fields. `Edges` is not one of them
(V1). The only non-test `Edges` write in `host/daemon` is `:287 selected.Edges = edges`, which is
the selected **log entry** (`EntryView`, row 35+38) and not the object (V2). The template at
`host/workbench/render.go:160` is `{{with .Object}}{{range .Edges}}{{template "edge" .}}{{end}}{{end}}`,
and neither action has an `{{else}}`. The probe measured the section as
`<h2>Provenance walk</h2>\n\n</section>`, a heading followed by nothing, on all four page shapes
it tried: `?object=X`, `?object=X&payload=1`, `/workbench` and `?from=0&entry=0` (V9). So the blank
is not limited to object pages. With no object selected, the section is blank as well.

**F2 — every object page shows `GRADE UNAVAILABLE — ` with an empty reason (row 102).** No
non-test file in `host/daemon` mentions `Grade`. The same-file control `PayloadTruncated` appears 1
time (V3). The zero `GradeView` has `Available=false` and `Unavailable=""`, so the `{{else}}` arm at
`render.go:154` renders `<p>GRADE UNAVAILABLE — </p>`. The probe measured exactly that on both
object-page shapes (V9). `NewGradeUnavailable` has **0** production callers. Its four callers are
all in `render_test.go` (V14). The defect is therefore a zero value that never passes through a
constructor. A guard placed on the constructor would not see it. §2b is built on this fact.

**F3 — the store records one typed reference per object, and nothing points back to it.**
`store.Object` has five fields (`store.go:92-98`). One of them, `InterfaceHash`, is a
`hashref.HashRef`. `Provenance` and `SemanticID` are UTF-8 labels. The schema comment says so in
those words: "semantic_id and provenance are UTF-8 labels, not digest fields" (V4). The schema
declares **0** indexes over 9 tables. The only `WHERE` clauses on a ref column are the two primary
keys (`objects.hash_ref`, `log_entries.entry_index`) and the `verification_cache` pair (V5). The
daemon's `readStore` seam has five methods, all point reads by primary key (V6). No read can answer
"which log entry or world references object X" without walking the log.

**F4 — in production, an interface target is never a stored object.** Every production producer
sets `InterfaceHash = SumSHA256([]byte(semanticID))`, a hash of a label string (broker, registry,
replay, journal, evidence; V8). On the probe store, the interface of the real epoch-registry
bootstrap object is not stored, and neither is the test object's interface (`GetObject` ok=false;
controls ok=true, V9). A checked interface edge therefore renders `UNAVAILABLE` in production today.
This is stated up front and chosen deliberately in §2a.

**F5 — the daemon holds nothing that could compute a grade.** `gradeOf` exists only in
`world/types.ail`. **0** Go files name `gradeOf` or `EvidenceGrade` (V10). **0** non-test Go files
import `host/evidence` (control: 3 daemon files import `host/store`, V10). The only path to
`PROVEN` is `evidence.Validator.ValidateProof` followed by `Resolve`. It needs a 32-byte HMAC key, a
compiler identity and version, a required-identity list and the **subject** ref the proof is
about (`validator.go:63-95, 120-183, 211-221`, V12). Outside `host/workbench`, no Go code names
`TestReport` (V18). `gradeOf` never yields `PROVEN` (V11).

---

## §2 Decisions

### (a) Provenance edges: **project the one exact edge, and name a stop for each relation the store cannot answer. Never render a blank.**

Each candidate was classified by measurement, not by reading:

| Candidate edge | What the store records (V) | Class | Decision |
|---|---|---|---|
| `interface`: the envelope's `InterfaceHash` | `objects.interface_hash_ref`, a typed `HashRef` on the object itself (V4) | **(a) exact** | **Project it**, existence-checked exactly like `entryEdges`: stored → link, unstored → `UNAVAILABLE: object <ref> is not stored`, store error → 5xx |
| log entries whose `transition_ref`/`transition_fn_ref`/`interpreter_ref` is X | columns of `log_entries`, **no index**; the seam reads entries by index only (V5, V6) | **(b) unbounded scan** | named stop `referencedBy` |
| worlds whose `state_root` is X | `worlds.state_root`, **no index**; `GetWorld` by world ref only (V5, V6) | **(b)** | folded into `referencedBy` |
| the commit that inserted X | **not recorded**: `objects` has 5 columns, none of them an entry or commit ref (V4) | **(c) absent** | named stop `committedBy` |
| the `provenance` string | a free-text label. Production writes `host/broker`, `store/journal`, `epoch-registry-bootstrap`, `host/transitionreg` or a replay caller label, and `POST /v1/commit` stores whatever the client sends (V7) | **(c) label, not a ref** | still rendered as text in the `<dl>` (unchanged); **not parsed** (original design §2.7, V13) |
| "is X the selected head's `StateRoot`" | the `?object=` page already reads the head world (`:209-237`) | exact, but **head-only** | **rejected**: a head-only reverse edge cannot distinguish "not the head's root" from "not a root anywhere", so its absence would carry two meanings. That world-pane edge belongs to row 99 |
| hits among the timeline entries the page already reads | an `?object=` page always reads entries 0–99 (`from := int64(0)`, `:168`, loop `:292`) | exact for hits, silent for misses | **rejected**: whether a link appears would depend on how long the log is, and a miss would look like an absence |

The walk's edge list is fixed and ordered. It holds `interface` (checked), then `committedBy` and
`referencedBy` (named stops). Their reasons name the missing store fact:

- `committedBy`: `the store records no commit-to-object relation, and the provenance field is a free-text label, not a reference`
- `referencedBy`: `no store index maps an object to the log entries or worlds that reference it`

The `interface` edge renders `UNAVAILABLE` in production today (F4). It is kept anyway for two
reasons. It is the envelope's only typed reference, so it is the one edge the store answers
exactly. And any producer that stores its schema bytes turns it into a working link with no further
change. Presenting `interface` as a provenance relation stretches the word "provenance". That is
accepted: the original design defines the section as a walk over *typed edges*, where "every stored
edge is clickable; every absent typed edge is a named `UNAVAILABLE` stop" (original design §1).

**The blank is closed at two layers.** (i) The template gets two `{{else}}` arms. A `range` with no
edges renders `UNAVAILABLE: no provenance edges were supplied for this object`. A page with no
object renders `UNAVAILABLE: no object selected`. After this change, no `Page` value can render the
heading followed by nothing. (ii) The daemon always supplies its three edges, and a daemon test
forbids the layer-(i) fallback text on daemon object pages. If the daemon regresses, CI goes red,
while production still shows a named gap and never a blank.

### (b) Grade: **every object gets `GRADE UNAVAILABLE — <one named reason>`. No label is ever derived.**

| Object kind | What a grade would need | Derivable from stored data? |
|---|---|---|
| labelled as a proof envelope (`SemanticID` `world/proof-report/v1`) | `Validator` + HMAC key + compiler config + required identities + the **subject** ref. The resulting grade belongs to the subject, not to the envelope. The daemon imports `host/evidence` 0 times (V10, V12). Its semantic ID is also a caller-supplied label (V7), so a grade keyed on it would be grade laundering | **no** |
| `TestReport` / any `Evidence` | `Evidence` exists only as AILANG values. No Go decoder exists, and no stored `TestReport` kind exists (V10, V18) | **no** |
| every other object (registries, broker, journal, replay sources, client objects) | no evidence relation at all | **no** |

No row is derivable, so the daemon has **one** constant reason and **no** per-kind branch. A branch
on `SemanticID` would trust a label (§2a). The reason reuses the original design's wording (§2.6,
V13) and adds the measured cause:

`no canonical host projection: evidence grades are computed only by gradeOf in world/types.ail, and the store records no decoded evidence for an object`

**The empty-reason guard: the render layer fails closed. `NewGradeUnavailable("")` is not
refused.** The template's `{{else}}` arm becomes
`GRADE UNAVAILABLE — {{with .Grade.Unavailable}}{{.}}{{else}}no grade reason was supplied{{end}}`.
Three measured reasons decide it:

1. The defect is a **zero value that never reaches a constructor** (F2: 0 production callers, 0
   `Grade` mentions in the daemon). A refusing constructor would leave today's defect exactly as it
   is. Only a guard at the point of rendering covers every future caller that forgets the field.
2. The daemon passes a **constant**, so a `(GradeView, error)` constructor would give the daemon an
   `if err != nil` branch that can never run. The repo has already learned that a branch no
   mutation can pin is not a guard (`host/evidence/validator.go:150-154`, iteration 106).
3. The fallback text is honest. It says the reason is missing and neither invents a reason nor
   implies a grade. `TestWorkbenchObjectGrade` forbids that text on daemon pages, so the daemon
   must supply the real reason (mutants H2, H3 and H16 prove the test fails without it).

---

## §3 Design

### 3.1 Template (`host/workbench/render.go`), two lines

| Anchor at `6127ff3` | Before | After |
|---|---|---|
| `:154` (grade `{{else}}` arm) | `{{else}}<p>GRADE UNAVAILABLE — {{.Grade.Unavailable}}</p>{{end}}` | `{{else}}<p>GRADE UNAVAILABLE — {{with .Grade.Unavailable}}{{.}}{{else}}no grade reason was supplied{{end}}</p>{{end}}` |
| `:160` (provenance walk) | `{{with .Object}}{{range .Edges}}{{template "edge" .}}{{end}}{{end}}` | `{{with .Object}}{{range .Edges}}{{template "edge" .}}{{else}}<p><span class="unavailable" role="note">UNAVAILABLE: no provenance edges were supplied for this object</span></p>{{end}}{{else}}<p><span class="unavailable" role="note">UNAVAILABLE: no object selected</span></p>{{end}}` |

No view-model type changes. The stop markup matches the existing `"edge"` define's unavailable
span byte-for-byte, so the page has one visual idiom for "unavailable".

### 3.2 Handler (`host/daemon/workbench.go`)

```go
// Named reasons for what the object inspector cannot show. Each names the
// missing store fact; none is a guess at the value it stands in for.
const (
	objectGradeUnavailableReason = "no canonical host projection: evidence grades are computed only by gradeOf in world/types.ail, and the store records no decoded evidence for an object"
	objectCommittedByMissing     = "the store records no commit-to-object relation, and the provenance field is a free-text label, not a reference"
	objectReferencedByMissing    = "no store index maps an object to the log entries or worlds that reference it"
)

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

1. **`checkedEdge` is extracted from `entryEdges` (`:121-144`).** It is the loop body's three
   branches, verbatim. `entryEdges` keeps its `refs` table and its loop, and the loop body becomes
   `edge, err := d.checkedEdge(ctx, item.relation, item.ref); if err != nil { return nil, err }; edges = append(edges, edge)`.
   The object walk and the selected entry then share one existence-check path. Row 35+38's
   killers still kill through it (§5; mutants H7, H8, H13).
2. **Object block (`:264-268`)**: before the literal, `edges, err := d.objectEdges(ctx, object)`,
   and on error `d.writeWorkbenchStoreError(w, r, ctx, err); return`. The literal gains
   `Grade: workbench.NewGradeUnavailable(objectGradeUnavailableReason), Edges: edges,`.

The read budget grows by exactly **one** `GetObject` per object page, on the request's `readCtx`
context. The count of `GetObject` mentions in `workbench.go` stays at 2 (V15), and no
`context.Background` is added. `GetObject` reads the payload, as the row-35+38 doc's F7 recorded,
so this one read can fetch up to one commit body. That is the same cost the selected entry already
pays three times.

### 3.3 Invariants this sprint makes persistent

- **I1** The provenance-walk section is never a heading followed by nothing, for any `Page` value.
  Pinned at the render layer by `TestRenderProvenanceWalkNeverBlank` and on live pages by
  `TestWorkbenchObjectProvenanceWalk/never-blank`.
- **I2** No object page renders an empty grade reason or any grade label. Pinned by
  `TestRenderGradeReasonNeverEmpty` and `TestWorkbenchObjectGrade`. Both make **absence**
  assertions (`GRADE UNAVAILABLE — </p>`, the fallback text, `<span>PROVEN|TESTED|ATTESTED|CLAIMED</span>`),
  not only presence assertions.
- **I3** An edge link is emitted only for a target the request itself found stored. A missing
  target is `UNAVAILABLE`, never a link, and a store error is 5xx, never `UNAVAILABLE`. This is the
  row-35+38 I2 invariant, now on the object walk too.

---

## §4 Query grammar delta

**None.** `acceptedWorkbenchKeys` (`:34-40`) and `supportedWorkbenchQuery` (`:63-77`) stay
byte-identical (AC5). The only link this design adds is the `interface` edge's `?object=<ref>`,
which uses the existing `object` state. That link is emitted only after `GetObject` returned `ok`
in the same request (I3), and `/interface-stored-link` fetches it and requires 200.

---

## §5 Conflict Surface

| Touched | How | Consequence |
|---|---|---|
| **Row 35+38** `entryEdges` (`:121-144`) | Its loop body moves into `checkedEdge` | That doc's §7 H9/H10 anchors move into `checkedEdge`. Their killers are unchanged and still fire on the pre-existing suite: H7 → `TestWorkbenchSelectedEntry/unstored-edge-unavailable`, H8/H13 → `TestWorkbenchSelectedEntry/object-store-error` (§7 base-suite column) |
| **Row 34** grade line `render.go:154`, pinned by `TestRenderGradeWithoutVerdictClaim/unavailable` | The `{{else}}` arm gains `{{with}}…{{else}}…{{end}}` | The pristine prototype passes that test unchanged (V15), and it still kills mutant R2 (§7). `supportedWorkbenchQuery`/`workbenchHref` are untouched |
| **Row 99** (world `StateRoot` unchecked) | **Not touched.** The world-pane block `:219-237` is byte-identical | `StateRoot` appears 0 times in the prototype diff (V15, AC9). `checkedEdge` is the helper row 99 will want, which makes row 99 cheaper. This sprint does not absorb it |
| **Row 100** (unrendered fields; extend the ratchet to `ObjectView`) | `ObjectView.Grade`/`.Edges` are now written, and both already had template reads | Row 100's `ObjectView` extension has two fewer fields to reconcile. The ratchet itself is not extended here |
| **Row 96** (semantic-ID lookup) | Not touched | — |
| `TestRenderUnavailableProvenanceEdge`, `TestRenderEscapesAllObjectText`, `TestRenderEmitsOnlyLocalLinks` | Unchanged | Pass on the prototype (V15) |
| `TestWorkbenchTimelinePaging/emitted-links-resolve` | Scans only the timeline section and the selected-entry article | The new walk links are outside its regions. `/interface-stored-link` covers them |
| Read-deadline tests (`read_deadline_test.go`) | One extra `GetObject` on object pages, on `readCtx` | Pass on the prototype (V15) |

Files the implementation changes: `host/workbench/render.go`, `host/workbench/render_test.go`,
`host/daemon/workbench.go`, `host/daemon/workbench_test.go`. Nothing else.

---

## §6 Acceptance criteria

`export AILANG_BIN=$HOME/.pinned-ailang/ailang` (v0.41.0, V17) applies to every command.

- **AC1** `go vet ./...` → rc=0.
- **AC2** `go test ./... -count=1` → rc=0.
- **AC3** `bash ./scripts/verify_ail.sh` → rc=0 ending `verify gate PASSED: 11 required identities verified, 40 named tests pass`,
  and `git diff --name-only origin/dev -- '*.ail' | wc -l` → `0`.
- **AC4 — the new tests exist and pass.** `go test ./host/workbench ./host/daemon -count=1 -v -run 'TestRenderProvenanceWalkNeverBlank|TestRenderGradeReasonNeverEmpty|TestWorkbenchObjectGrade|TestWorkbenchObjectProvenanceWalk'`
  → rc=0, with a `--- PASS` line for each of the 4 top-level names and 16 subtests, and 0 `--- FAIL`. Assertions:
  - `TestRenderProvenanceWalkNeverBlank` (render). The helper `provenanceWalkBody` returns the
    trimmed text between the section's `</h2>` and `</section>`, and fails if the markers are
    missing. Each arm first asserts that text is **non-empty**. `/no-object` (`Page{}`): contains
    `<p><span class="unavailable" role="note">UNAVAILABLE: no object selected</span></p>`, and does
    **not** contain `no provenance edges were supplied`. `/object-without-edges` (`Object: &ObjectView{}`):
    contains the `…no provenance edges were supplied for this object…` stop, and does **not** contain
    `no object selected`. `/supplied-edges` (one available edge): contains
    `<p>interface: <a href="/workbench?object=abc"`, and does **not** contain `UNAVAILABLE`.
  - `TestRenderGradeReasonNeverEmpty` (render). `/zero-grade` (`Object: &ObjectView{}`): contains
    `<p>GRADE UNAVAILABLE — no grade reason was supplied</p>`, and does **not** contain
    `GRADE UNAVAILABLE — </p>`. `/supplied-reason` (`NewGradeUnavailable("named reason")`):
    contains `<p>GRADE UNAVAILABLE — named reason</p>`, and does **not** contain `no grade reason was supplied`.
  - `TestWorkbenchObjectGrade` (daemon). Seed genesis plus one `testCommit` that also stores a
    proof-report-labelled object (`SemanticID "world/proof-report/v1"`, interface
    `SumSHA256("world/authenticated-proof-envelope/v1")`). **Control**: `GetRegistryHead(EpochRegistryV1)`
    returns ok (the real bootstrap-written object). Six subtests: {`test-object`, `proof-labelled`,
    `registry`} × {``, `&payload=1`}. Each does `GET /workbench?object=<ref><suffix>` → 200. The
    count of `<p>GRADE UNAVAILABLE — ` + `objectGradeUnavailableReason` + `</p>` must be **exactly 1**.
    **Absence**: none of `GRADE UNAVAILABLE — </p>`, `no grade reason was supplied`,
    `<span>PROVEN</span>`, `<span>TESTED</span>`, `<span>ATTESTED</span>`, `<span>CLAIMED</span>`.
  - `TestWorkbenchObjectProvenanceWalk` (daemon). One commit stores `plain` (the `testCommit`
    object, whose interface is unstored), `schema`, and `typed` (with `InterfaceHash = schema.Hash`).
    **Controls**: `GetObject(schema.Hash)` ok=true and `GetObject(plain.InterfaceHash)` ok=false. Subtests:
    `/interface-stored-link`: the `?object=<typed>` section contains `<p>interface: <a href="/workbench?object=<schema>"`, and fetching that href → 200.
    `/interface-unstored-unavailable`: the `?object=<plain>` section contains
    `<p>interface: <span class="unavailable" role="note">UNAVAILABLE: object <plain.InterfaceHash> is not stored</span></p>`,
    and the **whole body** contains no `href="/workbench?object=<plain.InterfaceHash>"`.
    `/named-stops`: on `?object=<plain>` and `?object=<typed>&payload=1`, the section contains the exact
    `committedBy` and `referencedBy` stop paragraphs built from the constants, and does **not**
    contain `no provenance edges were supplied`.
    `/never-blank`: on `/workbench`, `?from=0&entry=0`, `?object=<plain>` and `?object=<typed>&payload=1`,
    the text after `</h2>` is non-empty, and the section does **not** contain `UNAVAILABLE: </span>`
    (a stop with an empty reason).
    `/interface-store-error`: `d.reads = refFailingStore{readStore: orig, fail: schema.Hash}`, which
    fails `GetObject` for that one ref. **Control**: `?object=<plain>` → 200 on the same store.
    Then `?object=<typed>` → 500 and the body contains `>Internal<`.
- **AC5 — the grammar is byte-identical.** `diff <(git show origin/dev:host/daemon/workbench.go | sed -n '/^func supportedWorkbenchQuery/,/^}/p') <(sed -n '/^func supportedWorkbenchQuery/,/^}/p' host/daemon/workbench.go)`
  → rc=0. The same holds for `/^var acceptedWorkbenchKeys/,/^}/`, and for `/^func TestWorkbenchRefusalBranches/,/^}/` on `workbench_test.go`.
- **AC6 — the dead reads are fed. These are instrument-health controls only (S6); AC7 carries the
  load-bearing claims.** `grep -c 'Grade: workbench.NewGradeUnavailable(objectGradeUnavailableReason)' host/daemon/workbench.go`
  → `1`. `grep -c 'Edges: edges,' host/daemon/workbench.go` → `1`. `grep -c 'NewGradeView' host/daemon/workbench.go`
  → `0` (the daemon never constructs a grade label). Control: `grep -c PayloadTruncated host/daemon/workbench.go` → `1`.
- **AC7 — mutation drill.** Apply each §7 row on its own to the landed code. For each row: `go vet ./host/workbench ./host/daemon`
  → rc=0 (it **compiles**), `go test ./host/workbench ./host/daemon -count=1 -v` → rc=1, and the
  named killer appears in the `--- FAIL:` set. A build failure never counts as a kill. Restore by
  copying a backup back (`cp`) and confirm `git diff --quiet -- host/` against the landed commit
  before the next row. The pristine run → rc=0. Record the per-row table (compiled y/n, rc, red
  leaves, restored) in the sprint's Verification Log.
- **AC8 — size.** `git diff --shortstat origin/dev -- host/` → insertions + deletions ≤ 330. The
  prototype measured 289 (V15).
- **AC9 — row 99 is not absorbed.** `git diff origin/dev -- host/daemon/workbench.go | grep -c StateRoot` → `0`.

---

## §7 Mutation table

The rows are enumerated from the prototype's diff (V15): every new or changed template action,
every new conditional, the edge constructor, each reason string, and the empty-reason guard. **Every
row was executed at design time** by `.proto186/drill.py`, a throwaway driver deleted before
commit. For each row, the driver copied the prototype into place, applied the single replacement
(asserting it matched exactly once), ran `go vet ./host/workbench ./host/daemon` and then
`go test ./host/workbench ./host/daemon -count=1 -v`, and parsed the `--- FAIL:` lines. It then
re-ran the suite with the **`6127ff3` test files** against the same mutant ("base suite"), and
restored everything by `cp`. A base-suite rc of 0 means **the mutant survives the suite at `6127ff3`**,
so the new tests are what kill it. A base-suite rc of 1 means a pre-existing test already kills it,
and the table names that test. The pristine control P0 was vet 0, suite 0, base suite 0.

"Leaves" lists the failing tests at their deepest level (a parent is omitted when one of its
subtests failed). "Sole" is `y` when the named killer is the only leaf. H7's first form,
`if !ok {` → `if false {`, **did not compile** (`ok` declared and not used). Per the rule, it was
replaced by the compiling form shown.

| ID | File · anchor | Mutation (old → new) | Compiles | Suite rc | Named killer (in red set) | Sole | Base suite at `6127ff3` |
|---|---|---|---|---|---|---|---|
| R1 | render.go · grade `{{else}}` | `{{with .Grade.Unavailable}}{{.}}{{else}}no grade reason was supplied{{end}}` → `{{.Grade.Unavailable}}` | y | 1 | `TestRenderGradeReasonNeverEmpty/zero-grade` | y | 0 (survives) |
| R2 | render.go · grade `{{else}}` | same old → `no grade reason was supplied` | y | 1 | `TestRenderGradeReasonNeverEmpty/supplied-reason` | n (+`TestRenderGradeWithoutVerdictClaim/unavailable`, 6× `TestWorkbenchObjectGrade/*`) | 1 (`TestRenderGradeWithoutVerdictClaim/unavailable`) |
| R3 | render.go · walk `range`-else | delete `{{else}}<p><span …>UNAVAILABLE: no provenance edges were supplied for this object</span></p>` | y | 1 | `TestRenderProvenanceWalkNeverBlank/object-without-edges` | y | 0 (survives) |
| R4 | render.go · walk `with`-else | delete `{{else}}<p><span …>UNAVAILABLE: no object selected</span></p>` | y | 1 | `TestRenderProvenanceWalkNeverBlank/no-object` | n (+`TestWorkbenchObjectProvenanceWalk/never-blank`) | 0 (survives) |
| R5 | render.go · walk `range` body | `{{range .Edges}}{{template "edge" .}}{{else}}…` → `{{range .Edges}}{{else}}…` | y | 1 | `TestRenderProvenanceWalkNeverBlank/supplied-edges` | n (+`TestRenderUnavailableProvenanceEdge`, 4 walk subtests) | 1 (`TestRenderUnavailableProvenanceEdge`) |
| H1 | workbench.go · object literal | `Edges: edges,` → `Edges: edges[:0],` | y | 1 | `TestWorkbenchObjectProvenanceWalk/named-stops` | n (+`/interface-stored-link`, `/interface-unstored-unavailable`) | 0 (survives) |
| H2 | workbench.go · object literal | delete `Grade: workbench.NewGradeUnavailable(objectGradeUnavailableReason),` | y | 1 | `TestWorkbenchObjectGrade` | y (its 6 subtests only) | 0 (survives) |
| H3 | workbench.go · object literal | `NewGradeUnavailable(objectGradeUnavailableReason)` → `NewGradeUnavailable("")` | y | 1 | `TestWorkbenchObjectGrade` | y (its 6 subtests only) | 0 (survives) |
| H4 | workbench.go · object literal | `workbench.NewGradeUnavailable(objectGradeUnavailableReason)` → `workbench.GradeView{Available: true, Label: workbench.GradeCLAIMED}` | y | 1 | `TestWorkbenchObjectGrade` | y (its 6 subtests only) | 0 (survives) |
| H5 | workbench.go · `objectEdges` | `if err != nil {` → `if false && err != nil {` | y | 1 | `TestWorkbenchObjectProvenanceWalk/interface-store-error` | y | 0 (survives) |
| H6 | workbench.go · object block | after `objectEdges`: `if err != nil {` → `if false && err != nil {` | y | 1 | `TestWorkbenchObjectProvenanceWalk/interface-store-error` | y | 0 (survives) |
| H7 | workbench.go · `checkedEdge` | `if !ok {` → `if false && !ok {` | y | 1 | `TestWorkbenchObjectProvenanceWalk/interface-unstored-unavailable` | n (+`TestWorkbenchSelectedEntry/unstored-edge-unavailable`, `TestWorkbenchTimelinePaging/emitted-links-resolve`) | 1 (`TestWorkbenchSelectedEntry/unstored-edge-unavailable`) |
| H8 | workbench.go · `checkedEdge` | `if err != nil {` → `if false && err != nil {` | y | 1 | `TestWorkbenchObjectProvenanceWalk/interface-store-error` | n (+`TestWorkbenchSelectedEntry/object-store-error`) | 1 (`TestWorkbenchSelectedEntry/object-store-error`) |
| H9 | workbench.go · `objectEdges` | `d.checkedEdge(ctx, "interface", object.InterfaceHash)` → `…, object.Hash)` | y | 1 | `TestWorkbenchObjectProvenanceWalk/interface-stored-link` | n (+`/interface-unstored-unavailable`, `/interface-store-error`) | 0 (survives) |
| H10 | workbench.go · `objectEdges` | delete `{Relation: "committedBy", Missing: objectCommittedByMissing},` | y | 1 | `TestWorkbenchObjectProvenanceWalk/named-stops` | y | 0 (survives) |
| H11 | workbench.go · `objectEdges` | delete `{Relation: "referencedBy", Missing: objectReferencedByMissing},` | y | 1 | `TestWorkbenchObjectProvenanceWalk/named-stops` | y | 0 (survives) |
| H12 | workbench.go · `objectEdges` | `{Relation: "referencedBy", Missing: objectReferencedByMissing}` → `{Relation: "referencedBy", Available: true, Target: object.Hash.String(), Href: "?object=" + object.Hash.String()}` (an edge the store does not back) | y | 1 | `TestWorkbenchObjectProvenanceWalk/named-stops` | y | 0 (survives) |
| H13 | workbench.go · `entryEdges` loop | `if err != nil {` → `if false && err != nil {` | y | 1 | `TestWorkbenchSelectedEntry/object-store-error` | y | 1 (`TestWorkbenchSelectedEntry/object-store-error`) |
| H14 | workbench.go · const | `objectCommittedByMissing = "…"` → `= ""` | y | 1 | `TestWorkbenchObjectProvenanceWalk/never-blank` | y | 0 (survives) |
| H15 | workbench.go · const | `objectReferencedByMissing = "…"` → `= ""` | y | 1 | `TestWorkbenchObjectProvenanceWalk/never-blank` | y | 0 (survives) |
| H16 | workbench.go · const | `objectGradeUnavailableReason = "…"` → `= ""` | y | 1 | `TestWorkbenchObjectGrade` | y (its 6 subtests only) | 0 (survives) |
| H17 | workbench.go · `objectEdges` | relation `"interface"` → `"transitionRef"` | y | 1 | `TestWorkbenchObjectProvenanceWalk/interface-stored-link` | n (+`/interface-unstored-unavailable`) | 0 (survives) |

**Totals (recounted by command from the rows above, V19):** 22 mutants. All 22 compile, all 22
turn the suite red, and all 22 have the named killer in the red set. For 14 the named killer is
the sole leaf, or its own subtests are the only leaves. 17 survive the suite at `6127ff3`. The other
5 (R2, R5, H7, H8, H13) are already killed by a named pre-existing test, which shows the row-34
and row-35+38 pins still hold through the refactor.

H14 and H15 show why `/never-blank` asserts **absence** of `UNAVAILABLE: </span>`. `/named-stops`
builds its expected text from the same constants, so emptying a constant would still match a
presence-only check. The absence assertion is the only one that catches it. H3 and H16 are killed
the same way, by the absence of `GRADE UNAVAILABLE — </p>`.

---

## §8 Milestones

| M | Content | ~LOC (prototype) | Exit |
|---|---|---|---|
| **M1 — the render guards** | `render.go` `:154` and `:160` (§3.1); `render_test.go`: `provenanceWalkBody`, `TestRenderProvenanceWalkNeverBlank`, `TestRenderGradeReasonNeverEmpty` | 70 (2+2 render.go, 66 tests) | AC1, AC2; R1–R5 drilled |
| **M2 — the daemon supplies both** | `workbench.go`: the three constants, the `checkedEdge` extraction, `objectEdges`, the object-block wiring; `workbench_test.go`: `provenanceWalkSection`, `refFailingStore`, `workbenchTestObject`, `TestWorkbenchObjectGrade`, `TestWorkbenchObjectProvenanceWalk` | 219 (48+10 workbench.go, 161 tests) | AC1–AC9; H1–H17 drilled |

Each milestone ends green on AC1 and AC2. M1 can land alone: with the render guards in place and
before M2, a live object page shows `no grade reason was supplied` and
`no provenance edges were supplied for this object`. Those are honest stops, not blanks.

---

## §9 Out of scope, and proposed queue rows

Nothing here widens this sprint. The proposed rows are measured and are for the controller to file
or drop.

- **Row 99** (world `StateRoot` unchecked), **row 100** (the ratchet over `WorldView`/`ObjectView`),
  **row 96** (semantic-ID lookup): not touched (§5). Row 99 can reuse `checkedEdge` directly.
- **P-a (proposed row) — `w-store-object-reference-index`, clause-5.** Makes `referencedBy` exact.
  Evidence: the schema declares 0 indexes over 9 tables (V5). `log_entries.transition_ref`/
  `transition_fn_ref`/`interpreter_ref` and `worlds.state_root` have no index. The `readStore` seam
  has no reverse read (V6). THE ITEM: add indexes plus a bounded `readStore` read ("entries/worlds
  referencing X, limit N"), and replace the `referencedBy` stop with checked edges. It needs a
  schema change and a seam change, so it is excluded here by the brief.
- **P-b (proposed row) — `w-store-commit-object-membership`, clause-5.** Makes `committedBy` exact.
  Evidence: `objects` has 5 columns and none of them relates an object to the commit or entry that
  inserted it (V4). `Store.Commit` inserts `Objects` in the same transaction as the entry but
  records no link between them. THE ITEM: record commit membership (a schema change), then project
  it.
- **P-c (proposed row, likely PARKED behind item 13's validated-boundary follow-on) —
  `w-workbench-grade-projection`.** A real grade requires stored, decoded evidence bound to a
  subject, and a daemon-held validator (F5, V10, V12). Until then, §2b's constant reason is the
  correct rendering. Filing it only records that the constant is a waypoint and not the destination.
- Parsing `provenance` text as a reference, searching payload bytes for hash-shaped strings, and
  promoting `writtenBy` to an agent link all stay forbidden (original design §2.7, V13).

---

## §10 Verification Log

All commands were run in this worktree at `6127ff3` on 2026-09-24 unless marked *prototype*.
Every zero is paired with a positive control in the same scope.

| # | Claim | Command | Result |
|---|---|---|---|
| V0 | Base is `origin/dev` `6127ff3` | `git rev-parse origin/dev HEAD` | both `6127ff32491e37067e05e6b80f7b839c1ee6df24` |
| V1 | The object literal writes no `Grade` or `Edges` | `sed -n 264,268p host/daemon/workbench.go` | `page.Object = &workbench.ObjectView{ Hash, InterfaceHash, SemanticID, Provenance, PayloadShown, PayloadPreview, PayloadTruncated }`, 7 fields |
| V2 | The only daemon `Edges` write is the selected entry's | `/usr/bin/grep -rn --include='*.go' -E '\bEdges\b' host/daemon \| /usr/bin/grep -v _test.go` | exactly 1 hit: `host/daemon/workbench.go:287: selected.Edges = edges` |
| V3 | No daemon source mentions `Grade` | `/usr/bin/grep -ln Grade host/daemon/*.go \| /usr/bin/grep -v _test \| wc -l`; control `/usr/bin/grep -c PayloadTruncated host/daemon/workbench.go` | `0`; control `1` |
| V4 | `store.Object` and `objects` carry one typed ref; provenance is a label | `sed -n 90,98p host/store/store.go`; `sed -n 9,20p host/store/schema.sql` | fields `Hash, InterfaceHash hashref.HashRef; SemanticID, Provenance string; Payload []byte`; columns `hash_ref, interface_hash_ref, semantic_id, provenance, payload`; comment `semantic_id and provenance are UTF-8 labels, not digest fields` |
| V5 | No index; no reverse lookup by ref | `/usr/bin/grep -c 'CREATE INDEX\|CREATE UNIQUE INDEX' host/store/schema.sql host/store/store.go`; `/usr/bin/grep -c 'CREATE TABLE' host/store/schema.sql`; `/usr/bin/grep -rn --include='*.go' -E 'WHERE (transition_ref\|transition_fn_ref\|interpreter_ref\|state_root\|object_ref)' host/store \| grep -v _test`; control `… -E 'WHERE (entry_index\|hash_ref) '` | indexes `0` / `0`; tables `9`; ref-column WHERE: only `store.go:790` (`verification_cache` pair); control: `store.go:484` objects by `hash_ref`, `:572` log by `entry_index`, `:699`, `scan.go:60`, `read_object.go:31,33` |
| V6 | The daemon's read seam is 5 point reads | `sed -n 331,337p host/daemon/daemon.go` | `GetObject(ref)`, `GetWorld(ref)`, `GetLogEntry(index)`, `GetRegistryHead(name)`, `SelectedHead()` |
| V7 | Provenance values are producer labels or client text | `/usr/bin/grep -rn 'Provenance:' --include='*.go' host cmd \| grep -v _test` | `broker/record.go:120 "host/broker"`, `transitionreg.go:256 "host/transitionreg"`, `registry.go:96 "epoch-registry-bootstrap"`, `replay.go:393 provenance` (caller arg), `journal.go:313 "store/journal"`, `handlers.go:660 object.Provenance` (the `POST /v1/commit` request body) |
| V8 | Production interface hashes are hashes of label strings | `/usr/bin/grep -rn 'InterfaceHash' --include='*.go' host cmd \| grep -v _test` (read at each producer) | `record.go:118`, `registry.go:95`, `replay.go:392`, `journal.go:313`: `SumSHA256([]byte(semanticID))`; `evidence/validator.go:19`: `SumSHA256([]byte("world/authenticated-proof-envelope/v1"))`; `transitionreg/codec.go:30` a fixed constant |
| V9 | **Live page at base** (probe) | throwaway `host/daemon/zz_probe186_test.go`: `newHandlerDaemon` + `seedGenesisEmbedded` + one `testCommit`; GET four targets; log the `<p>GRADE UNAVAILABLE[^<]*</p>` match and the provenance-walk section; `GetObject` on the interface refs. `go test ./host/daemon -run TestZZProbe186 -count=1 -v`. **Deleted afterwards** | `?object=X` and `?object=X&payload=1`: 200, grade `"<p>GRADE UNAVAILABLE — </p>"`, section `"<section aria-label=\"provenance walk\">\n<h2>Provenance walk</h2>\n\n</section>"`; `/workbench` and `?from=0&entry=0`: 200, 0 `GRADE`, the same blank section. `GetObject(testObject.InterfaceHash)` ok=false (control `GetObject(obj.Hash)` ok=true); epoch-registry head ok=true, object `sid="world/epoch-registry/v1" prov="epoch-registry-bootstrap"`, its interface stored=**false**; transition-registry head ok=false |
| V10 | No Go grade code; the daemon cannot reach the validator | `/usr/bin/grep -rlE 'gradeOf\|EvidenceGrade' --include='*.go' . \| wc -l`, control `/usr/bin/grep -lE gradeOf world/*.ail`; `/usr/bin/grep -rl --include='*.go' 'ailang-world/host/evidence"' . \| grep -v _test.go \| wc -l`, control `/usr/bin/grep -l 'ailang-world/host/store"' host/daemon/*.go \| grep -v _test \| wc -l` | `0`, control `world/types.ail`; `0`, control `3` |
| V11 | `gradeOf` never yields PROVEN | `sed -n 42,59p world/types.ail \| /usr/bin/grep -c PROVEN`; control `sed -n 15,40p world/types.ail` | `0`; the type declares `PROVEN \| TESTED \| ATTESTED \| CLAIMED`; `ProofReceipt(_) => CLAIMED` |
| V12 | PROVEN needs key, config, identities and subject | `sed -n 63,95p;120,183p;211,221p host/evidence/validator.go` | `NewValidator(key [32]byte, reader, cfg CompilerConfig, requiredIdentities)` refuses nil reader, zero compiler, empty version, empty identities; `ValidateProof(ctx, reportRef, expectedSubject)` checks semantic ID, interface, MAC, subject, tool, proof, identities; `Resolve` → `ResolvedGradeProven` only for its own seal |
| V13 | The original design's grade and walk rules | `sed -n 206,214p;255,262p design_docs/implemented/w-workbench-read-only.md` | `GRADE UNAVAILABLE — no canonical host projection`; "does not downgrade the object to `CLAIMED`"; "object `provenance` text is not parsed as a reference, and payload bytes are not searched for hash-shaped substrings" |
| V14 | `NewGradeUnavailable` has no production caller | `/usr/bin/grep -rn 'NewGradeUnavailable' --include='*.go' .` | definition `render.go:187`; callers only `render_test.go:43, 53, 121, 168` |
| V15 | *prototype*: the design compiles, the suite is green, size | prototype applied to the four files (§3), then `go vet ./...`; `go test ./... -count=1`; `git diff --shortstat -- host/`; `git diff --numstat -- host/`; `grep -c GetObject` and `grep -c context.Background` on the prototype `workbench.go`; `grep -c StateRoot` on the prototype diff; `gofmt -l` on the prototype files | vet rc=0; test rc=0 (no non-`ok` line); `4 files changed, 277 insertions(+), 12 deletions(-)`; numstat `48 10 workbench.go`, `161 0 workbench_test.go`, `2 2 render.go`, `66 0 render_test.go`; `GetObject` 2 (base 2); `context.Background` 0; `StateRoot` 0; gofmt clean |
| V16 | *prototype*: the new render tests are red at base | the prototype `render_test.go` run against the `6127ff3` `render.go`, `go test ./host/workbench -run 'TestRenderProvenanceWalkNeverBlank\|TestRenderGradeReasonNeverEmpty' -v` | FAIL `/no-object`, FAIL `/object-without-edges`, FAIL `/zero-grade`; PASS `/supplied-edges`, PASS `/supplied-reason` (controls). The daemon tests do not compile at base (`undefined: objectGradeUnavailableReason`), so their base behaviour is carried by mutants H1/H2, which restore it and survive the base suite |
| V17 | Pinned binary; `.ail` gate is unaffected | `$AILANG_BIN --version`; *prototype* `bash ./scripts/verify_ail.sh` | `AILANG v0.41.0` (`24ee108`); rc=0, `verify gate PASSED: 11 required identities verified, 40 named tests pass` |
| V18 | No stored `TestReport` kind | `/usr/bin/grep -rn --include='*.go' TestReport host cmd \| grep -v _test \| grep -v host/workbench \| wc -l`; control `/usr/bin/grep -c KindTestReport host/workbench/render.go` | `0`; control `2` |
| V19 | §7 totals, recounted from the table | `D=design_docs/planned/w-workbench-object-provenance-and-grade.md; T=$(/usr/bin/grep -E '^\| (R\|H)[0-9]+ \|' $D)`; rows `echo "$T" \| wc -l`; compiles `echo "$T" \| awk -F'\|' '$5 ~ /^ y $/' \| wc -l`; red `… '$6 ~ /^ 1 $/'`; killer `… '$7 ~ /Test/'`; sole `… '$8 ~ /^ y/'`; survive base `… '$9 ~ /^ 0/'` | 22; 22; 22; 22; 14; 17 |

---

## §11 Quorum log

| Round | Reviewer | Verdict | Objection (one line) | Disposition |
|---|---|---|---|---|
| — | — | — | — | — |
