package daemon

// Row 243 V3 sprint 3 M9 (design_docs/planned/w-world-address-bar-verbs-v3.md §2.5): the
// confirmation card's read routes, each behind approvalRoute's Host check, served as the decision
// card's are; the data is surfacepage.ConfirmationCard of the row. Nothing here writes (POST: sprint 4).

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/broker"
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

// readConfirmations fills out.Confirmations with the live, pending verb requests on one read-only
// handle: a summary decoded from the stored bytes, or Refused and no summary.
func (d *Daemon) readConfirmations(ctx context.Context, now int64, out *surfacepage.RequestsJSON) error {
	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
	if err != nil {
		return err
	}
	defer func() { _ = ro.Close() }()
	rows, err := ro.PendingVerbRequests(ctx, now, store.MaxPendingVerbRequests)
	if err != nil {
		return err
	}
	for _, row := range rows {
		out.Confirmations = append(out.Confirmations, pendingConfirmation(row))
	}
	return nil
}

// pendingConfirmation summarises row, or refuses it when its bytes do not hash to its digest, do
// not decode as a valid world/verb-request/v1, or carry no finish approval scope.
func pendingConfirmation(row store.VerbRequestRow) surfacepage.PendingConfirmation {
	p := surfacepage.PendingConfirmation{ID: row.RequestID, Expires: row.ExpiresAt}
	if sum := sha256.Sum256(row.Request); hex.EncodeToString(sum[:]) != row.Digest {
		p.Refused = "the stored request does not match its digest"
		return p
	}
	r, err := approval.DecodeVerb(row.Request)
	if err == nil {
		err = r.Validate()
	}
	if err != nil {
		p.Refused = "the stored request is not a valid verb request: " + err.Error()
		return p
	}
	s, err := broker.ParseFinishApprovalScope(r.Scope)
	if err != nil {
		p.Refused = "the stored request's scope is not a finish approval"
		return p
	}
	p.Episode, p.Repository, p.Head, p.SHA, p.Draft = r.Episode, s.Repository, s.Head, s.SHA, &s.Draft
	return p
}
