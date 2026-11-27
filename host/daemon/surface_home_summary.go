package daemon

import (
	"net/http"

	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// handleHomeSummary serves GET /surface/home-summary.json, the Studio start page's card values
// (surfacepage.HomeSummarySchema). Needs you is counted from the same reads as requests.json; the
// other values come from the follow model's Now edge, and are left out without follow mode or
// before the first accepted graph. It reads only, never writes, and carries no names or paths.
func (d *Daemon) handleHomeSummary(w http.ResponseWriter, r *http.Request) {
	schema := surfacepage.HomeSummarySchema
	if r.URL.RawQuery != "" {
		d.writeSurfaceJSON(w, r, http.StatusBadRequest, schema, surfaceReason{Schema: schema, Reason: unknownSurfaceKeyMessage})
		return
	}
	ctx, cancel := d.readCtx(r)
	defer cancel()
	now := d.approveNow()
	var reqs surfacepage.RequestsJSON
	needs := -1
	if err := d.readRequests(ctx, now, &reqs); err != nil {
		d.writeInternalErrorLog(r, err)
	} else {
		needs = len(reqs.Pending)
		if err := d.readDecisionList(ctx, now, &reqs); err != nil {
			d.writeInternalErrorLog(r, err)
		}
		needs += len(reqs.Decisions) + len(reqs.OperatorActions)
		if reqs.OperatorActionsUnavailable != "" {
			needs = -1
		}
	}
	var edge *surfacepage.NowEdge
	if d.follow != nil {
		if s := d.followSnapshot(); len(s.Envelope.Repos) > 0 {
			m := surfacepage.BuildFollowModel(s.Envelope, d.followSources(s), d.followNow(r))
			edge = &m.NowEdge
		}
	}
	d.writeSurfaceJSON(w, r, http.StatusOK, schema, surfacepage.BuildHomeSummary(needs, edge))
}
