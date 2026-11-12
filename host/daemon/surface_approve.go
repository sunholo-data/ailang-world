package daemon

// Row 201 A3 sprint 2 (design_docs/planned/w-world-approvals-by-key-a3.md §2.3, §2.4): the
// approval routes under /surface, none under /v1. GET /surface/requests.json (Needs you and the
// decision list) and GET /surface/requests/{id} (the card's data) read the credentials store on a
// fresh read-only handle; POST /surface/approve/{id} records a passkey or SSH approval only after
// the A1b/A2 law verified it. Each route answers only a Host in the surface origin set O (R2); every
// existing route stays exactly as row 199 leaves it.

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// storedRequest is row as the approval model reads it.
func storedRequest(row store.MintRequestRow) surfacepage.StoredRequest {
	return surfacepage.StoredRequest{ID: row.RequestID, Episode: row.EpisodeID, Digest: row.Digest, Challenge: row.Challenge, Request: row.Request,
		Created: row.CreatedAt, Expires: row.ExpiresAt, ConsumedAt: row.ConsumedAt, Consumed: row.Consumed, ApprovedBy: row.ApprovedBy, CredentialID: row.CredentialID}
}

// writeApprovalJSON writes an approval route's JSON: no-store, nosniff, and never a CORS header.
func writeApprovalJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	writeJSON(w, status, v)
}

// writeApprovalRefusal is every refusal of the approval routes: inert JSON, nothing written.
func writeApprovalRefusal(w http.ResponseWriter, status int, class, reason string) {
	writeApprovalJSON(w, status, surfacepage.ApproveRefusal{V: 1, Error: class, Reason: reason})
}

// approvalInternal logs err for the operator and answers the sanitized 500.
func (d *Daemon) approvalInternal(w http.ResponseWriter, r *http.Request, err error) {
	d.writeInternalErrorLog(r, err)
	writeApprovalRefusal(w, http.StatusInternalServerError, "Internal", internalErrorMessage)
}

// isOlderStore reports a store below the schema a read needs (which reads as nothing there).
func isOlderStore(err error) bool {
	var older *store.OlderCredentialsStoreError
	return errors.As(err, &older)
}

// approvalOrigins is O for r (design §2.3): the loopback origins at the port r arrived on (its
// connection's local address, so an ephemeral port too), every --surface-origin, and the origin of
// every active passkey on a fresh read-only handle. A store it cannot read adds no passkey origin.
func (d *Daemon) approvalOrigins(ctx context.Context, r *http.Request) []string {
	var port string
	if addr, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr); ok {
		_, port, _ = net.SplitHostPort(addr.String())
	}
	var enrolled []string
	if ro, err := store.OpenCredentialsReadOnly(d.credentialsPath); err == nil {
		active, _ := ro.ActivePasskeys(ctx)
		_ = ro.Close()
		for _, e := range active {
			enrolled = append(enrolled, e.Passkey.Origin)
		}
	}
	return surfaceOrigins(port, d.cfg.SurfaceOrigins, enrolled)
}

// approvalRoute is an approval route behind R2: a Host outside O is 421 before anything else, and
// the handler gets O for its own checks.
func (d *Daemon) approvalRoute(h func(http.ResponseWriter, *http.Request, []string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := d.readCtx(r)
		o := d.approvalOrigins(ctx, r)
		cancel()
		if !allowedApprovalHost(r.Host, o) {
			writeApprovalRefusal(w, http.StatusMisdirectedRequest, "MisdirectedRequest", misdirectedReason(r.Host))
			return
		}
		h(w, r, o)
	}
}

// handleSurfaceRequests is GET /surface/requests.json: the pending requests, each with its
// recorded approver or null, and the last MaxDecidedListed decided ones. A store that cannot be
// read answers 200 with unavailable and empty lists, as /surface/model.json reports its sessions.
func (d *Daemon) handleSurfaceRequests(w http.ResponseWriter, r *http.Request, _ []string) {
	now := d.approveNow()
	out := surfacepage.RequestsJSON{V: 1, Now: now}
	ctx, cancel := d.readCtx(r)
	defer cancel()
	if err := d.readRequests(ctx, now, &out); err != nil {
		d.writeInternalErrorLog(r, err)
		out.Pending, out.Decided, out.Unavailable = nil, nil, "the mint requests could not be read"
	}
	if out.Pending == nil {
		out.Pending = []surfacepage.PendingRequest{}
	}
	if out.Decided == nil {
		out.Decided = []surfacepage.DecidedRequest{}
	}
	writeApprovalJSON(w, http.StatusOK, out)
}

// readRequests fills out's lists on one read-only handle. Below schema 7 no approval is recorded.
func (d *Daemon) readRequests(ctx context.Context, now int64, out *surfacepage.RequestsJSON) error {
	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
	if err != nil {
		return err
	}
	defer func() { _ = ro.Close() }()
	pending, err := ro.PendingMintRequests(ctx, now, store.MaxPendingMintRequests)
	if err != nil {
		return err
	}
	for _, row := range pending {
		rec, _, err := ro.SurfaceApproval(ctx, row.RequestID) // none: rec.Approver is ""
		if err != nil && !isOlderStore(err) {
			return err
		}
		out.Pending = append(out.Pending, surfacepage.Pending(storedRequest(row), rec.Approver))
	}
	decided, err := ro.DecidedMintRequests(ctx, "", surfacepage.MaxDecidedListed)
	if err != nil {
		return err
	}
	for _, row := range decided {
		out.Decided = append(out.Decided, surfacepage.Decided(storedRequest(row)))
	}
	return nil
}

// handleSurfaceRequestCard is GET /surface/requests/{id}: the stored bytes, the digest, the raw
// nonce and the times, and the active passkeys by origin (none below schema 6). Never a message hash.
func (d *Daemon) handleSurfaceRequestCard(w http.ResponseWriter, r *http.Request, _ []string) {
	id := r.PathValue("id")
	if !store.MintRequestIDPattern.MatchString(id) {
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "a request id is mr- and 16 lowercase hex")
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
	row, ok, err := ro.MintRequest(ctx, id)
	if err != nil {
		d.approvalInternal(w, r, err)
		return
	}
	if !ok {
		writeApprovalRefusal(w, http.StatusNotFound, "NotFound", "no mint request "+id)
		return
	}
	active, err := ro.ActivePasskeys(ctx)
	if err != nil && !isOlderStore(err) {
		d.approvalInternal(w, r, err)
		return
	}
	var passkeys []surfacepage.CardPasskey
	for _, e := range active {
		passkeys = append(passkeys, surfacepage.CardPasskey{Origin: e.Passkey.Origin, RPID: e.Passkey.RPID, CredentialID: e.Passkey.CredentialID})
	}
	writeApprovalJSON(w, http.StatusOK, surfacepage.RequestCard(storedRequest(row), passkeys))
}

// approvalStoreError answers a store failure of the POST: a store below the schema a step needs is
// 503 naming the restart, a busy writer 503 Busy with Retry-After (never a wait), anything else
// the sanitized 500.
func (d *Daemon) approvalStoreError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case isOlderStore(err):
		writeApprovalRefusal(w, http.StatusServiceUnavailable, "Restart",
			"the credentials store predates surface approvals: restart `ailang-worldd serve` with this build once (it upgrades the store), then approve again")
	case store.IsWriterAlreadyActive(err):
		w.Header().Set("Retry-After", "1")
		writeApprovalRefusal(w, http.StatusServiceUnavailable, "Busy", "the credentials store is busy (a mint or a session verb is writing it); try again")
	default:
		d.approvalInternal(w, r, err)
	}
}

// handleSurfaceApprove is POST /surface/approve/{id} (design §2.4): after R2 (approvalRoute), the
// checks R1, R3-R12 in order on a read-only handle (R11 judges the passkey assertion or the SSH
// signature with the CLI's own judge), each refusal writing nothing; then R13 a
// non-waiting writer, opened for R14's RecordSurfaceApproval and closed at once. It reads no
// Cookie and no Authorization, answers no CORS, and never mints: the minter (approve --await)
// judges the recorded artifact again.
func (d *Daemon) handleSurfaceApprove(w http.ResponseWriter, r *http.Request, o []string) {
	id := r.PathValue("id")
	if !store.MintRequestIDPattern.MatchString(id) { // R1
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "a request id is mr- and 16 lowercase hex")
		return
	}
	if origin := r.Header.Get("Origin"); !allowedApprovalOrigin(origin, o) { // R3
		writeApprovalRefusal(w, http.StatusForbidden, "CrossOrigin", crossOriginReason(origin))
		return
	}
	if !sameOriginFetch(r.Header.Get("Sec-Fetch-Site")) { // R4
		writeApprovalRefusal(w, http.StatusForbidden, "CrossOrigin", "Sec-Fetch-Site is not same-origin")
		return
	}
	if !jsonMediaType(r.Header.Get("Content-Type")) { // R5
		writeApprovalRefusal(w, http.StatusUnsupportedMediaType, "UnsupportedMediaType", "the body must be application/json")
		return
	}
	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, surfacepage.MaxApproveBodyBytes)) // R6
	if err != nil {
		writeApprovalRefusal(w, http.StatusRequestEntityTooLarge, "TooLarge", fmt.Sprintf("the body is over %d bytes", surfacepage.MaxApproveBodyBytes))
		return
	}
	body, err := surfacepage.ParseApproveBody(raw)
	if err != nil {
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", err.Error())
		return
	}
	now := d.approveNow()
	ctx, cancel := d.readCtx(r)
	defer cancel()
	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
	if err != nil {
		d.approvalInternal(w, r, err)
		return
	}
	defer func() { _ = ro.Close() }()
	row, ok, err := ro.MintRequest(ctx, id) // R7 (serve upgraded the store to 7 at start)
	if err != nil {
		d.approvalInternal(w, r, err)
		return
	}
	if !ok {
		writeApprovalRefusal(w, http.StatusNotFound, "NotFound", "no mint request "+id)
		return
	}
	if reason := surfacepage.AdmissionRefusal(storedRequest(row), now); reason != "" { // R8
		writeApprovalRefusal(w, http.StatusConflict, "Conflict", reason)
		return
	}
	if body.Digest != row.Digest { // R9
		writeApprovalRefusal(w, http.StatusConflict, "Conflict", "the request changed since the card was shown; reload it")
		return
	}
	challenge, _ := hex.DecodeString(row.Challenge) // R10 (the table's CHECK holds it to 64 characters)
	if len(challenge) != 32 {
		writeApprovalRefusal(w, http.StatusConflict, "Conflict", "request "+id+"'s stored challenge is not 64 hex; make a new request")
		return
	}
	format, artifact, what := approval.FormatWebAuthn, []byte(body.WebAuthn), "the passkey's assertion"
	if body.SSHSIG != nil {
		format, artifact, what = approval.FormatSSH, []byte(*body.SSHSIG), "the SSH signature"
	}
	var signed *approval.SignedApproval
	if format == approval.FormatSSH { // R11: the law, before anything is recorded
		signed, err = authority.JudgeSSHSignature(ctx, ro, row, artifact, now)
	} else {
		signed, err = authority.JudgeAssertion(ctx, ro, row, artifact, now)
	}
	if err != nil {
		var je *authority.JudgeError
		if errors.As(err, &je) && je.Step == authority.JudgeStepStore {
			d.approvalStoreError(w, r, err)
			return
		}
		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", what+" does not approve request "+id+": "+err.Error())
		return
	}
	digest, _ := hex.DecodeString(row.Digest)                                                  // R8 has shown it is 64 hex
	if err := signed.Binds([32]byte(digest), [32]byte(challenge), row.ExpiresAt); err != nil { // R12
		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the verified approval is not bound to request "+id+": "+err.Error())
		return
	}
	wr, err := store.OpenCredentials(d.credentialsPath) // R13: non-waiting, closed at once
	if err != nil {
		d.approvalStoreError(w, r, err)
		return
	}
	err = wr.RecordSurfaceApproval(ctx, store.SurfaceApprovalRow{RequestID: id, Format: format, Approver: signed.ApprovedBy(),
		Artifact: artifact, Digest: row.Digest}, d.approveNow()) // R14, at the time it is received
	_ = wr.Close()
	switch {
	case isRequestRace(err):
		writeApprovalRefusal(w, http.StatusConflict, "Conflict", "request "+id+" changed while it was approved: "+err.Error())
	case err != nil:
		d.approvalStoreError(w, r, err)
	default:
		writeApprovalJSON(w, http.StatusOK, surfacepage.Recorded(id, signed.ApprovedBy()))
	}
}
