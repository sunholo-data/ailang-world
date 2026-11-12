package daemon

// Row 201 A3 sprint 3 (design §2.7): the request card page and its script. Static embedded bytes
// under surfaceCSP, behind approvalRoute's Host check like every approval route; the page reads the
// request itself from GET /surface/requests/{id}.

import (
	"net/http"

	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// handleSurfaceApprovePage is GET /surface/approve/{id}: the same bytes for every well-formed id.
func (d *Daemon) handleSurfaceApprovePage(w http.ResponseWriter, r *http.Request, _ []string) {
	if !store.MintRequestIDPattern.MatchString(r.PathValue("id")) {
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "a request id is mr- and 16 lowercase hex")
		return
	}
	setSurfaceHeaders(w)
	_, _ = w.Write(surfacepage.ApprovePage)
}

// handleSurfaceApproveScript is the card's script, under the same headers.
func (d *Daemon) handleSurfaceApproveScript(w http.ResponseWriter, _ *http.Request, _ []string) {
	setSurfaceHeaders(w)
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	_, _ = w.Write(surfacepage.ApproveScript)
}
