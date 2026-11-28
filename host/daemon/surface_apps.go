package daemon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sunholo-data/ailang-world/host/appmanifest"
)

// Row 274 M1 (design_docs/planned/w-world-applications.md §2, §4): the read-only app
// catalog on the surface. GET /surface/apps.json lists the local catalog
// (~/.ailang-world/apps/<name>.json), each row validated by host/appmanifest, and the
// address-bar verb `app list` answers from the same rows. Nothing here launches an app:
// M2 resolves a manifest into a preset that row 243 V4's start path reviews.
const SurfaceAppsSchema = "world/surface-apps/v1"

type surfaceAppGrant struct {
	Effect string `json:"effect"`
	Budget int64  `json:"budget"`
}

type surfaceAppEntry struct {
	Name     string                `json:"name"`
	Valid    bool                  `json:"valid"`
	Title    string                `json:"title,omitempty"`
	Summary  string                `json:"summary,omitempty"`
	RepoKind string                `json:"repo_kind,omitempty"`
	Agent    string                `json:"agent,omitempty"`
	Mode     string                `json:"mode,omitempty"`
	Grants   []surfaceAppGrant     `json:"grants,omitempty"`
	Checks   []string              `json:"checks,omitempty"`
	Outputs  []string              `json:"outputs,omitempty"`
	Views    []string              `json:"views,omitempty"`
	Digest   string                `json:"digest,omitempty"`
	Problems []appmanifest.Problem `json:"problems,omitempty"`
	Error    string                `json:"error,omitempty"`
}

type surfaceAppsResponse struct {
	Schema string            `json:"schema"`
	Apps   []surfaceAppEntry `json:"apps"`
}

// readSurfaceApps reads the catalog on every call, as readSurfacePresets does. The file
// path is left out: the surface never shows a machine path.
func (d *Daemon) readSurfaceApps() ([]surfaceAppEntry, error) {
	entries, err := appmanifest.ReadCatalog(d.appsDir)
	if err != nil {
		return nil, err
	}
	rows := make([]surfaceAppEntry, 0, len(entries))
	for _, e := range entries {
		row := surfaceAppEntry{Name: e.Name, Valid: e.Valid(), Digest: e.Digest, Problems: e.Problems, Error: e.Error}
		if m := e.Manifest; m != nil {
			row.Title, row.Summary, row.RepoKind = m.Title, m.Summary, m.RepoKind
			row.Agent, row.Mode, row.Views = m.Preset.Agent, m.Preset.Launch.Mode, m.Views
			for _, g := range m.Preset.Grants {
				row.Grants = append(row.Grants, surfaceAppGrant{Effect: g.Effect, Budget: g.Budget})
			}
			for _, c := range m.Preset.Checks {
				row.Checks = append(row.Checks, c.Name)
			}
			for _, o := range m.Outputs {
				row.Outputs = append(row.Outputs, o.Glob)
			}
		}
		rows = append(rows, row)
	}
	return rows, nil
}

// handleSurfaceApps serves GET /surface/apps.json behind approvalRoute.
func (d *Daemon) handleSurfaceApps(w http.ResponseWriter, r *http.Request, _ []string) {
	if r.URL.RawQuery != "" {
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", unknownSurfaceKeyMessage)
		return
	}
	rows, err := d.readSurfaceApps()
	if err != nil {
		writeApprovalRefusal(w, http.StatusInternalServerError, "AppCatalogUnreadable", "the app catalog could not be read")
		return
	}
	writeApprovalJSON(w, http.StatusOK, surfaceAppsResponse{Schema: SurfaceAppsSchema, Apps: rows})
}

func (d *Daemon) verbAppList(args []string, a *verbAnswer) {
	asJSON, reason := parseListJSONFlag("app list", args)
	if reason != "" {
		a.Reason = reason
		return
	}
	rows, err := d.readSurfaceApps()
	if err != nil {
		a.Reason = "the app catalog could not be read; at a terminal: ailang-worldd app list"
		return
	}
	a.Outcome = "result"
	if asJSON {
		raw, _ := json.MarshalIndent(rows, "", "  ")
		a.JSON = string(raw) + "\n"
		return
	}
	if len(rows) == 0 {
		a.Text = "no apps installed\n"
		return
	}
	var sb strings.Builder
	for _, r := range rows {
		switch {
		case r.Error != "":
			fmt.Fprintf(&sb, "%-24s ERROR: %s\n", r.Name, r.Error)
		case !r.Valid:
			fmt.Fprintf(&sb, "%-24s INVALID %d problem(s): ailang-worldd app check %s\n", r.Name, len(r.Problems), r.Name)
		default:
			fmt.Fprintf(&sb, "%-24s %-6s %-11s grants=%d outputs=%d  %s\n", r.Name, r.Agent, r.Mode, len(r.Grants), len(r.Outputs), r.Title)
		}
	}
	a.Text = sb.String()
}
