package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sunholo-data/ailang-world/host/appmanifest"
	"github.com/sunholo-data/ailang-world/host/appregistry"
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
	Name string `json:"name"`
	// State is the listing law's verdict (published, published-local-differs, local), and
	// Launchable is D-WORLD-129 = A: only published apps launch (row 274 M1b-3).
	State       string                `json:"state"`
	Launchable  bool                  `json:"launchable"`
	Version     int64                 `json:"version,omitempty"`
	Manifest    string                `json:"manifest,omitempty"`
	PublishedAt string                `json:"published_at,omitempty"`
	Valid       bool                  `json:"valid"`
	Title       string                `json:"title,omitempty"`
	Summary     string                `json:"summary,omitempty"`
	RepoKind    string                `json:"repo_kind,omitempty"`
	Agent       string                `json:"agent,omitempty"`
	Mode        string                `json:"mode,omitempty"`
	Grants      []surfaceAppGrant     `json:"grants,omitempty"`
	Checks      []string              `json:"checks,omitempty"`
	Outputs     []string              `json:"outputs,omitempty"`
	Views       []string              `json:"views,omitempty"`
	Digest      string                `json:"digest,omitempty"`
	Problems    []appmanifest.Problem `json:"problems,omitempty"`
	Error       string                `json:"error,omitempty"`
}

type surfaceAppsResponse struct {
	Schema string            `json:"schema"`
	Apps   []surfaceAppEntry `json:"apps"`
}

// readSurfaceApps merges the app registry (the daemon's own store; absent when the daemon
// has none) with the local catalog, one row per app, through host/appregistry's listing. It
// reads on every call, as readSurfacePresets does, and never shows a machine path.
func (d *Daemon) readSurfaceApps(ctx context.Context) ([]surfaceAppEntry, error) {
	local, err := appmanifest.ReadCatalog(d.appsDir)
	if err != nil {
		return nil, err
	}
	var snap appregistry.Snapshot
	if d.store != nil {
		if snap, err = appregistry.Read(ctx, d.store); err != nil {
			return nil, err
		}
	}
	listing := appregistry.List(ctx, d.store, snap, local)
	rows := make([]surfaceAppEntry, 0, len(listing))
	for _, l := range listing {
		row := surfaceAppEntry{Name: l.Name, State: l.State, Launchable: l.Launchable(), Version: l.Version,
			Manifest: l.ManifestRef, PublishedAt: l.PublishedAt, Valid: l.Valid(), Problems: l.Problems, Error: l.Error}
		if l.Local != nil {
			row.Digest = l.Local.Digest
		}
		if m := l.Manifest; m != nil {
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
	ctx, cancel := d.readCtx(r)
	defer cancel()
	rows, err := d.readSurfaceApps(ctx)
	if err != nil {
		writeApprovalRefusal(w, http.StatusInternalServerError, "AppCatalogUnreadable", "the app catalog or registry could not be read")
		return
	}
	writeApprovalJSON(w, http.StatusOK, surfaceAppsResponse{Schema: SurfaceAppsSchema, Apps: rows})
}

func (d *Daemon) verbAppList(r *http.Request, args []string, a *verbAnswer) {
	asJSON, reason := parseListJSONFlag("app list", args)
	if reason != "" {
		a.Reason = reason
		return
	}
	ctx, cancel := d.readCtx(r)
	defer cancel()
	rows, err := d.readSurfaceApps(ctx)
	if err != nil {
		a.Reason = "the app catalog or registry could not be read; at a terminal: ailang-worldd app list"
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
	for _, row := range rows {
		r := row
		switch {
		case r.Error != "":
			fmt.Fprintf(&sb, "%-24s ERROR: %s\n", r.Name, r.Error)
		case !r.Valid:
			fmt.Fprintf(&sb, "%-24s INVALID %d problem(s): ailang-worldd app check %s\n", r.Name, len(r.Problems), r.Name)
		default:
			fmt.Fprintf(&sb, "%-24s %-24s %-6s %-11s grants=%d outputs=%d  %s\n", r.Name, appStateWords(r), r.Agent, r.Mode, len(r.Grants), len(r.Outputs), r.Title)
		}
	}
	a.Text = sb.String()
}

// appStateWords is a row's state in words, as the CLI prints it.
func appStateWords(r surfaceAppEntry) string {
	switch r.State {
	case appregistry.StatePublished:
		return fmt.Sprintf("published v%d", r.Version)
	case appregistry.StateLocalDiffers:
		return fmt.Sprintf("published v%d (draft differs)", r.Version)
	}
	return "local (not published)"
}
