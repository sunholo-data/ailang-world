package daemon

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/workbench"
)

const WorkbenchPageLimit = workbench.WorkbenchPageLimit

const workbenchCSP = "default-src 'none'; style-src 'unsafe-inline'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'"

const (
	unknownWorkbenchKeyMessage           = "unsupported workbench query parameter"
	duplicateWorkbenchKeyMessage         = "duplicate workbench query parameter"
	unsupportedWorkbenchQueryMessage     = "unsupported workbench parameter combination"
	malformedPayloadFlagMessage          = "malformed payload flag"
	malformedWorkbenchWorldMessage       = "malformed world reference"
	absentWorkbenchWorldMessage          = "world reference not found"
	malformedWorkbenchObjectMessage      = "malformed object reference"
	absentWorkbenchObjectMessage         = "object reference not found"
	negativeWorkbenchFromMessage         = "from index must be non-negative"
	malformedWorkbenchEntryMessage       = "malformed entry index"
	absentWorkbenchEntryMessage          = "log entry not found"
	workbenchFromOverflowMessage         = "from index overflows"
	workbenchInternalStoreFailureMessage = internalErrorMessage
	referenceIndexUnavailableMessage     = "object-reference indexes are absent or incompatible; open this store writable once to provision them, then reopen this read-only handle"
)

// Named reasons for what the object inspector cannot show. Each names the
// missing store fact; none is a guess at the value it stands in for.
const (
	objectGradeUnavailableReason = "no canonical host projection: this workbench handler does not resolve subject-bound evidence into an object grade"
)

var acceptedWorkbenchKeys = map[string]bool{
	"world":        true,
	"object":       true,
	"from":         true,
	"entry":        true,
	"payload":      true,
	"refsAfter":    true,
	"commitsAfter": true,
}

func setWorkbenchHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Security-Policy", workbenchCSP)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "no-referrer")
}

var workbenchErrorTemplate = template.Must(template.New("workbench-error").Parse(
	`<!doctype html><html lang="en"><head><meta charset="utf-8"><title>{{.Class}}</title></head><body><h1>{{.Class}}</h1><p>{{.Message}}</p></body></html>`,
))

func writeWorkbenchError(w http.ResponseWriter, status int, class, message string) {
	setWorkbenchHeaders(w)
	w.WriteHeader(status)
	_ = workbenchErrorTemplate.Execute(w, struct {
		Class   string
		Message string
	}{Class: class, Message: message})
}

func supportedWorkbenchQuery(query map[string][]string) bool {
	if len(query) == 0 {
		return true
	}
	if query["object"] != nil {
		for key := range query {
			if key != "object" && key != "payload" && key != "refsAfter" && key != "commitsAfter" {
				return false
			}
		}
		return true
	}
	return (len(query) == 1 && query["world"] != nil) || (len(query) == 2 && query["from"] != nil && query["entry"] != nil)
}

func objectPageHref(ref hashref.HashRef, payload string, refsAfter *store.ObjectReferenceCursor, commitsAfter *int64) string {
	q := url.Values{"object": {ref.String()}}
	if payload != "" {
		q.Set("payload", payload)
	}
	if refsAfter != nil {
		q.Set("refsAfter", encodeReferenceCursor(*refsAfter))
	}
	if commitsAfter != nil {
		q.Set("commitsAfter", strconv.FormatInt(*commitsAfter, 10))
	}
	return "?" + q.Encode()
}

func commitCursor(query url.Values, parsed int64) *int64 {
	if query["commitsAfter"] == nil {
		return nil
	}
	return &parsed
}

func (d *Daemon) writeWorkbenchStoreError(w http.ResponseWriter, r *http.Request, ctx context.Context, err error) {
	if timedOut(ctx, err) {
		writeWorkbenchError(w, http.StatusServiceUnavailable, "Timeout", "workbench read deadline exceeded")
		return
	}
	var unavailable *store.ReferenceIndexUnavailableError
	if errors.As(err, &unavailable) {
		writeWorkbenchError(w, http.StatusServiceUnavailable, "ReferenceIndexUnavailable", referenceIndexUnavailableMessage)
		return
	}
	d.writeWorkbenchInternalError(w, r, err)
}

func (d *Daemon) writeWorkbenchInternalError(w http.ResponseWriter, r *http.Request, err error) {
	// Keep operator detail in the daemon log and constant text on the wire.
	// The query string is deliberately excluded, matching writeInternalError.
	if err != nil {
		d.writeInternalErrorLog(r, err)
	}
	writeWorkbenchError(w, http.StatusInternalServerError, "Internal", workbenchInternalStoreFailureMessage)
}

func (d *Daemon) writeInternalErrorLog(r *http.Request, err error) {
	// Kept local to the HTML adapter so its response remains HTML rather than
	// passing through the JSON-only writeInternalError helper.
	fmt.Fprintf(d.errLog, "ailang-worldd: internal error: %s %s: %v\n", r.Method, r.URL.Path, err)
}

func entryView(entry store.LogEntry) workbench.EntryView {
	return workbench.EntryView{
		EntryIndex:     entry.Header.EntryIndex,
		EntryHash:      entry.EntryHash.String(),
		PrevEntryHash:  entry.Header.PrevEntryHash.String(),
		SemanticsEpoch: entry.Header.SemanticsEpoch,
		WrittenBy:      entry.Header.WrittenBy,
	}
}

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
		edge, err := d.checkedEdge(ctx, item.relation, item.ref)
		if err != nil {
			return nil, err
		}
		edges = append(edges, edge)
	}
	return edges, nil
}

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

// objectEdges checks the object's one typed interface edge.
func (d *Daemon) objectEdges(ctx context.Context, object store.Object) ([]workbench.EdgeView, error) {
	iface, err := d.checkedEdge(ctx, "interface", object.InterfaceHash)
	if err != nil {
		return nil, err
	}
	return []workbench.EdgeView{iface}, nil
}

// committedByEdges checks the existence of every displayed carrying entry.
func (d *Daemon) committedByEdges(ctx context.Context, indexes []int64) ([]workbench.EdgeView, error) {
	edges := make([]workbench.EdgeView, 0, len(indexes))
	for _, index := range indexes {
		_, commitEntryOK, commitReadErr := d.reads.GetLogEntry(ctx, index)
		if commitReadErr != nil {
			return nil, commitReadErr
		}
		if !commitEntryOK {
			edges = append(edges, workbench.EdgeView{Relation: "committedBy", Target: fmt.Sprintf("entry %d", index), Missing: fmt.Sprintf("log entry %d is not stored", index)})
			continue
		}
		from := index
		if from > math.MaxInt64-WorkbenchPageLimit {
			from = math.MaxInt64 - WorkbenchPageLimit
		}
		edges = append(edges, workbench.EdgeView{Relation: "committedBy", Target: fmt.Sprintf("entry %d", index), Href: pageHref(from, index), Available: true})
	}
	return edges, nil
}

func (d *Daemon) checkedReferenceEdge(ctx context.Context, ref hashref.HashRef, item store.ObjectReference) (workbench.EdgeView, error) {
	c := item.Cursor
	var actual hashref.HashRef
	var edge workbench.EdgeView
	switch c.Kind {
	case store.ReferenceTransitionRef, store.ReferenceTransitionFn, store.ReferenceInterpreter:
		entry, entryOK, err := d.reads.GetLogEntry(ctx, c.EntryIndex)
		if err != nil {
			return edge, err
		}
		if !entryOK {
			return workbench.EdgeView{Relation: referenceRole(c.Kind), Target: fmt.Sprintf("entry %d", c.EntryIndex), Missing: "source entry is no longer stored"}, nil
		}
		switch c.Kind {
		case store.ReferenceTransitionRef:
			actual = entry.TransitionRef
		case store.ReferenceTransitionFn:
			actual = entry.Header.TransitionFn
		case store.ReferenceInterpreter:
			actual = entry.Header.Interpreter
		}
		from := c.EntryIndex
		if from > math.MaxInt64-workbench.WorkbenchPageLimit {
			from = math.MaxInt64 - workbench.WorkbenchPageLimit
		}
		edge = workbench.EdgeView{Relation: referenceRole(c.Kind), Target: fmt.Sprintf("entry %d", c.EntryIndex), Href: pageHref(from, c.EntryIndex), Available: true}
	case store.ReferenceStateRoot:
		world, worldOK, err := d.reads.GetWorld(ctx, c.WorldRef)
		if err != nil {
			return edge, err
		}
		if !worldOK {
			return workbench.EdgeView{Relation: "stateRoot", Target: c.WorldRef.String(), Missing: "source world is no longer stored"}, nil
		}
		actual = world.StateRoot
		edge = workbench.EdgeView{Relation: "stateRoot", Target: c.WorldRef.String(), Href: "?world=" + c.WorldRef.String(), Available: true}
	default:
		return edge, fmt.Errorf("invalid reference kind %d", c.Kind)
	}
	if actual != ref {
		return workbench.EdgeView{}, fmt.Errorf("reference source relation mismatch")
	}
	return edge, nil
}

func referenceRole(kind store.ReferenceKind) string {
	switch kind {
	case store.ReferenceTransitionRef:
		return "transitionRef"
	case store.ReferenceTransitionFn:
		return "transitionFn"
	case store.ReferenceInterpreter:
		return "interpreter"
	case store.ReferenceStateRoot:
		return "stateRoot"
	default:
		return "unknown"
	}
}

func (d *Daemon) handleWorkbench(w http.ResponseWriter, r *http.Request) {
	setWorkbenchHeaders(w)
	query := r.URL.Query()
	for key, values := range query {
		if !acceptedWorkbenchKeys[key] {
			writeWorkbenchError(w, http.StatusBadRequest, "BadRequest", unknownWorkbenchKeyMessage)
			return
		}
		if len(values) > 1 {
			writeWorkbenchError(w, http.StatusBadRequest, "BadRequest", duplicateWorkbenchKeyMessage)
			return
		}
	}
	if !supportedWorkbenchQuery(query) {
		writeWorkbenchError(w, http.StatusBadRequest, "BadRequest", unsupportedWorkbenchQueryMessage)
		return
	}
	if payload := query.Get("payload"); payload != "" && payload != "0" && payload != "1" {
		writeWorkbenchError(w, http.StatusBadRequest, "BadRequest", malformedPayloadFlagMessage)
		return
	}
	var refsAfter *store.ObjectReferenceCursor
	if values := query["refsAfter"]; values != nil {
		cursor, err := decodeReferenceCursor(values[0])
		if err != nil {
			writeWorkbenchError(w, http.StatusBadRequest, "BadRequest", "malformed refsAfter cursor")
			return
		}
		refsAfter = &cursor
	}

	afterEntry := int64(-1)
	var parsedCommitsAfter int64
	if values := query["commitsAfter"]; values != nil {
		var err error
		parsedCommitsAfter, err = strconv.ParseInt(values[0], 10, 64)
		if err != nil || parsedCommitsAfter < 0 {
			writeWorkbenchError(w, http.StatusBadRequest, "BadRequest", malformedWorkbenchEntryMessage)
			return
		}
		afterEntry = parsedCommitsAfter
	}

	from := int64(0)
	if text := query.Get("from"); text != "" {
		parsed, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			writeWorkbenchError(w, http.StatusBadRequest, "BadRequest", negativeWorkbenchFromMessage)
			return
		}
		from = parsed
		if from < 0 {
			writeWorkbenchError(w, http.StatusBadRequest, "BadRequest", negativeWorkbenchFromMessage)
			return
		}
	}
	limit := workbench.WorkbenchPageLimit
	if from > math.MaxInt64-int64(limit) {
		writeWorkbenchError(w, http.StatusBadRequest, "BadRequest", workbenchFromOverflowMessage)
		return
	}

	var selectedIndex *int64
	if text := query.Get("entry"); text != "" {
		parsed, err := strconv.ParseInt(text, 10, 64)
		if err != nil || parsed < 0 {
			writeWorkbenchError(w, http.StatusBadRequest, "BadRequest", malformedWorkbenchEntryMessage)
			return
		}
		selectedIndex = &parsed
	}

	ctx, cancel := d.readCtx(r)
	defer cancel()
	page := workbench.Page{Title: "AILANG World Workbench", World: workbench.WorldView{Unavailable: "no world selected"}}

	var worldRefText string
	if values := query["world"]; values != nil {
		ref, err := parseRef(values[0], "world ref")
		if err != nil {
			writeWorkbenchError(w, http.StatusBadRequest, "BadRequest", malformedWorkbenchWorldMessage)
			return
		}
		worldRefText = ref.String()
	} else {
		head, ok, err := d.reads.SelectedHead(ctx)
		if err != nil {
			d.writeWorkbenchStoreError(w, r, ctx, err)
			return
		}
		if ok {
			worldRefText = head.String()
		}
	}
	if worldRefText != "" {
		ref, _ := parseRef(worldRefText, "world ref")
		world, ok, err := d.reads.GetWorld(ctx, ref)
		if err != nil {
			d.writeWorkbenchStoreError(w, r, ctx, err)
			return
		}
		if !ok {
			writeWorkbenchError(w, http.StatusNotFound, "NotFound", absentWorkbenchWorldMessage)
			return
		}
		stateRoot, err := d.checkedEdge(ctx, "stateRoot", world.StateRoot)
		if err != nil {
			d.writeWorkbenchStoreError(w, r, ctx, err)
			return
		}
		page.World = workbench.WorldView{
			Ref:       world.Ref.String(),
			Revision:  world.Revision,
			StateRoot: stateRoot,
			LogHead:   world.LogHead.String(),
			Available: true,
		}
	}

	if values := query["object"]; values != nil {
		ref, err := parseRef(values[0], "object ref")
		if err != nil {
			writeWorkbenchError(w, http.StatusBadRequest, "BadRequest", malformedWorkbenchObjectMessage)
			return
		}
		object, ok, err := d.reads.GetObject(ctx, ref)
		if err != nil {
			d.writeWorkbenchStoreError(w, r, ctx, err)
			return
		}
		if !ok {
			writeWorkbenchError(w, http.StatusNotFound, "NotFound", absentWorkbenchObjectMessage)
			return
		}
		showPayload := false
		if query.Get("payload") == "1" {
			showPayload = true
		}
		preview := object.Payload
		truncated := false
		if len(preview) > workbench.MaxPayloadPreview {
			preview = preview[:workbench.MaxPayloadPreview]
			truncated = true
		}
		edges, err := d.objectEdges(ctx, object)
		if err != nil {
			d.writeWorkbenchStoreError(w, r, ctx, err)
			return
		}
		commits, err := d.reads.ObjectCommits(ctx, ref, afterEntry, WorkbenchPageLimit+1)
		if err != nil {
			d.writeWorkbenchStoreError(w, r, ctx, err)
			return
		}
		commitView := &workbench.CommitView{Truncated: len(commits) > WorkbenchPageLimit, Continued: query["commitsAfter"] != nil}
		if commitView.Continued {
			commitView.FirstHref = objectPageHref(ref, query.Get("payload"), refsAfter, nil)
		}
		if commitView.Truncated {
			nextCommit := commits[WorkbenchPageLimit-1]
			commitView.NextHref = objectPageHref(ref, query.Get("payload"), refsAfter, &nextCommit)
		}
		commitView.Edges, err = d.committedByEdges(ctx, commits[:min(len(commits), WorkbenchPageLimit)])
		if err != nil {
			d.writeWorkbenchStoreError(w, r, ctx, err)
			return
		}
		refs, err := d.reads.ObjectReferences(ctx, ref, refsAfter, WorkbenchPageLimit+1)
		if err != nil {
			d.writeWorkbenchStoreError(w, r, ctx, err)
			return
		}
		references := &workbench.ReferenceView{Truncated: len(refs) > WorkbenchPageLimit, Continued: refsAfter != nil}
		if refsAfter != nil {
			references.FirstHref = objectPageHref(ref, query.Get("payload"), nil, commitCursor(query, parsedCommitsAfter))
		}
		if references.Truncated {
			next := refs[WorkbenchPageLimit-1].Cursor
			references.NextHref = objectPageHref(ref, query.Get("payload"), &next, commitCursor(query, parsedCommitsAfter))
		}
		for _, item := range refs[:min(len(refs), workbench.WorkbenchPageLimit)] {
			edge, err := d.checkedReferenceEdge(ctx, ref, item)
			if err != nil {
				d.writeWorkbenchStoreError(w, r, ctx, err)
				return
			}
			references.Edges = append(references.Edges, edge)
		}
		page.Object = &workbench.ObjectView{
			Hash: object.Hash.String(), InterfaceHash: object.InterfaceHash.String(),
			SemanticID: object.SemanticID, Provenance: object.Provenance,
			PayloadShown: showPayload, PayloadPreview: string(preview), PayloadTruncated: truncated,
			Grade: workbench.NewGradeUnavailable(objectGradeUnavailableReason), Edges: edges, Commits: commitView, References: references,
		}
	}

	if selectedIndex != nil {
		entry, ok, err := d.reads.GetLogEntry(ctx, *selectedIndex)
		if err != nil {
			d.writeWorkbenchStoreError(w, r, ctx, err)
			return
		}
		if !ok {
			writeWorkbenchError(w, http.StatusNotFound, "NotFound", absentWorkbenchEntryMessage)
			return
		}
		selected := entryView(entry)
		edges, err := d.entryEdges(ctx, entry)
		if err != nil {
			d.writeWorkbenchStoreError(w, r, ctx, err)
			return
		}
		selected.Edges = edges
		page.Selected = &selected
	}

	page.Timeline = workbench.TimelineView{From: from, Limit: limit}
	for offset := int64(0); offset < int64(limit); offset++ {
		entry, ok, err := d.reads.GetLogEntry(ctx, from+offset)
		if err != nil {
			d.writeWorkbenchStoreError(w, r, ctx, err)
			return
		}
		if !ok {
			break
		}
		view := entryView(entry)
		view.SelectHref = pageHref(from, view.EntryIndex)
		page.Timeline.Entries = append(page.Timeline.Entries, view)
	}
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

	_ = workbench.Render(w, page)
}
