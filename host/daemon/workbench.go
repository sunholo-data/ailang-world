package daemon

import (
	"html/template"
	"net/http"

	"github.com/sunholo-data/ailang-world/host/workbench"
)

const workbenchCSP = "default-src 'none'; style-src 'unsafe-inline'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'"

const unknownWorkbenchKeyMessage = "unknown workbench query parameter"

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

func (d *Daemon) handleWorkbench(w http.ResponseWriter, r *http.Request) {
	setWorkbenchHeaders(w)
	for key := range r.URL.Query() {
		if !acceptedWorkbenchKeys[key] {
			writeWorkbenchError(w, http.StatusBadRequest, "BadRequest", unknownWorkbenchKeyMessage)
			return
		}
	}

	ctx, cancel := d.readCtx(r)
	defer cancel()

	head, headOK, headErr := d.reads.SelectedHead(ctx)
	page := workbench.Page{Title: "AILANG World Workbench"}
	if headErr == nil && headOK {
		world, worldOK, worldErr := d.reads.GetWorld(ctx, head)
		if worldErr == nil && worldOK {
			page.World = workbench.WorldView{
				Ref:       world.Ref.String(),
				Revision:  world.Revision,
				StateRoot: workbench.EdgeView{Available: true, Target: world.StateRoot.String(), Href: "?object=" + world.StateRoot.String()},
				LogHead:   world.LogHead.String(),
				Available: true,
			}
		}
	}

	page.Timeline = workbench.TimelineView{From: 0, Limit: workbench.WorkbenchPageLimit}
	for index := int64(0); index < int64(workbench.WorkbenchPageLimit); index++ {
		entry, ok, err := d.reads.GetLogEntry(ctx, index)
		if err != nil || !ok {
			break
		}
		page.Timeline.Entries = append(page.Timeline.Entries, workbench.EntryView{
			EntryIndex:     entry.Header.EntryIndex,
			EntryHash:      entry.EntryHash.String(),
			PrevEntryHash:  entry.Header.PrevEntryHash.String(),
			SemanticsEpoch: entry.Header.SemanticsEpoch,
			TransitionFn:   workbench.EdgeView{Available: true, Target: entry.Header.TransitionFn.String(), Href: "?object=" + entry.Header.TransitionFn.String()},
			Interpreter:    workbench.EdgeView{Available: true, Target: entry.Header.Interpreter.String(), Href: "?object=" + entry.Header.Interpreter.String()},
			TransitionRef:  workbench.EdgeView{Available: true, Target: entry.TransitionRef.String(), Href: "?object=" + entry.TransitionRef.String()},
			WrittenBy:      entry.Header.WrittenBy,
		})
	}

	_ = workbench.Render(w, page)
}
