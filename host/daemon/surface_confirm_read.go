package daemon

// Row 243 V3 sprint 3 M9 (design_docs/planned/w-world-address-bar-verbs-v3.md §2.5): the
// confirmation card's read routes, each behind approvalRoute's Host check, served as the decision
// card's are; the data is surfacepage.ConfirmationCard of the row. Nothing here writes (POST: sprint 4).

import (
	"net/http"

	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

const confirmIDReason = "a verb request id is vr- and 16 lowercase hex"

// handleSurfaceConfirmPage is GET /surface/confirm/{id}: the same bytes for every well-formed id.
func (d *Daemon) handleSurfaceConfirmPage(w http.ResponseWriter, r *http.Request, _ []string) {
	if !store.VerbRequestIDPattern.MatchString(r.PathValue("id")) {
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", confirmIDReason)
		return
	}
	setSurfaceHeaders(w)
	_, _ = w.Write(surfacepage.ConfirmPage)
}

// handleSurfaceConfirmScript is the card's script, under the same headers.
func (d *Daemon) handleSurfaceConfirmScript(w http.ResponseWriter, _ *http.Request, _ []string) {
	setSurfaceHeaders(w)
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	_, _ = w.Write(surfacepage.ConfirmScript)
}

// handleSurfaceConfirmation is GET /surface/confirmations/{id}: 400 for a bad id, 503 naming the
// serve restart below schema 9, 404 for an unknown one, else the card's data with the active
// passkeys by origin.
func (d *Daemon) handleSurfaceConfirmation(w http.ResponseWriter, r *http.Request, _ []string) {
	id := r.PathValue("id")
	if !store.VerbRequestIDPattern.MatchString(id) {
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", confirmIDReason)
		return
	}
	ctx, cancel := d.readCtx(r)
	defer cancel()
	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
	if err != nil {
		d.approvalInternal(w, r, err)
		return
	}
	defer func() { _ = ro.Close() }()
	row, ok, err := ro.VerbRequest(ctx, id)
	switch {
	case isOlderStore(err):
		writeApprovalRefusal(w, http.StatusServiceUnavailable, "Restart",
			"the credentials store predates verb requests: restart `ailang-worldd serve` with this build once (it upgrades the store)")
		return
	case err != nil:
		d.approvalInternal(w, r, err)
		return
	case !ok:
		writeApprovalRefusal(w, http.StatusNotFound, "NotFound", "no verb request "+id)
		return
	}
	active, err := ro.ActivePasskeys(ctx)
	if err != nil {
		d.approvalInternal(w, r, err)
		return
	}
	var passkeys []surfacepage.CardPasskey
	for _, e := range active {
		passkeys = append(passkeys, surfacepage.CardPasskey{Origin: e.Passkey.Origin, RPID: e.Passkey.RPID, CredentialID: e.Passkey.CredentialID})
	}
	writeApprovalJSON(w, http.StatusOK, surfacepage.ConfirmationCard(confirmationInput(row, d.approveNow()), passkeys))
}

// confirmationInput is row as the card reads it at now.
func confirmationInput(row store.VerbRequestRow, now int64) surfacepage.ConfirmationInput {
	return surfacepage.ConfirmationInput{ID: row.RequestID, Request: row.Request, Digest: row.Digest, Challenge: row.Challenge, Created: row.CreatedAt,
		Expires: row.ExpiresAt, Now: now, Outcome: row.Outcome, ApprovedBy: row.ApprovedBy, Decision: row.DecisionRef}
}
