package daemon

// Row 201 A4 sprint 4 (design §2.8): the decision card page and its script. Static embedded bytes
// under surfaceCSP, behind approvalRoute's Host check like every approval route; the page reads the
// subject itself from GET /surface/decisions/{ref} and POSTs to POST /surface/decide/{ref}, which
// runs the law again.

import (
	"net/http"

	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// handleSurfaceDecidePage is GET /surface/decide/{ref}: the same bytes for every well-formed ref.
func (d *Daemon) handleSurfaceDecidePage(w http.ResponseWriter, r *http.Request, _ []string) {
	if !decisionRefPattern.MatchString(r.PathValue("ref")) {
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "a request ref is sha256: and 64 lowercase hex")
		return
	}
	setSurfaceHeaders(w)
	_, _ = w.Write(surfacepage.DecidePage)
}

// handleSurfaceDecideScript is the card's script, under the same headers.
func (d *Daemon) handleSurfaceDecideScript(w http.ResponseWriter, _ *http.Request, _ []string) {
	setSurfaceHeaders(w)
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	_, _ = w.Write(surfacepage.DecideScript)
}
