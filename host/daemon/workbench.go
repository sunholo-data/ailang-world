package daemon

import (
	"context"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"

	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/workbench"
)

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
)

var acceptedWorkbenchKeys = map[string]bool{
	"world":   true,
	"object":  true,
	"from":    true,
	"entry":   true,
	"payload": true,
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
	if len(query) == 1 {
		return query["world"] != nil || query["object"] != nil
	}
	if len(query) != 2 {
		return false
	}
	if query["from"] != nil && query["entry"] != nil {
		return true
	}
	return query["object"] != nil && query["payload"] != nil
}

func (d *Daemon) writeWorkbenchStoreError(w http.ResponseWriter, r *http.Request, ctx context.Context, err error) {
	if timedOut(ctx, err) {
		writeWorkbenchError(w, http.StatusServiceUnavailable, "Timeout", "workbench read deadline exceeded")
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
		TransitionFn:   workbench.EdgeView{Available: true, Target: entry.Header.TransitionFn.String(), Href: "?object=" + entry.Header.TransitionFn.String()},
		Interpreter:    workbench.EdgeView{Available: true, Target: entry.Header.Interpreter.String(), Href: "?object=" + entry.Header.Interpreter.String()},
		TransitionRef:  workbench.EdgeView{Available: true, Target: entry.TransitionRef.String(), Href: "?object=" + entry.TransitionRef.String()},
		WrittenBy:      entry.Header.WrittenBy,
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
	page := workbench.Page{Title: "AILANG World Workbench"}

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
		page.World = workbench.WorldView{
			Ref:       world.Ref.String(),
			Revision:  world.Revision,
			StateRoot: workbench.EdgeView{Available: true, Target: world.StateRoot.String(), Href: "?object=" + world.StateRoot.String()},
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
		page.Object = &workbench.ObjectView{
			Hash: object.Hash.String(), InterfaceHash: object.InterfaceHash.String(),
			SemanticID: object.SemanticID, Provenance: object.Provenance,
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
		page.Timeline.Entries = append(page.Timeline.Entries, entryView(entry))
	}

	_ = workbench.Render(w, page)
}
