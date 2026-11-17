package daemon

// Row 243 V3 sprint 4 M10 (design_docs/planned/w-world-address-bar-verbs-v3.md §2.6): POST
// /surface/confirm/{id}, a verified enrolled key's confirmation of one stored verb request. After R2
// (approvalRoute) the guards C0-C6 (A4's R0, A3a's R1, R3-R6); then, under humanMu and within
// commitBudget, C7-C12 on a read-only handle (each refusal writes nothing), C13 a non-waiting writer,
// C14 the world store's finish approval decided by the key, and C15 one credentials transaction that
// consumes the row and names the decision. The writer is held across C14 and C15, as
// handleSurfaceDecide holds keyArm's across DecideApproval and LandKeyDecision. A failure after C14
// leaves an approval no row names (inert, design W62) and the row pending: the same body retries.

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// confirmRequester is the finish approval's requester, the typed approve phase's (W5).
const confirmRequester = "operator:episode-finish"

// confirmIncomplete is the 500 of a failed C14 or C15: nothing the CLI uses was written.
const confirmIncomplete = "the confirmation did not complete; approve again"

// Test seams: C11's judges (C12's refusal needs a verified confirmation of other bytes) and C15's
// consume (a request superseded, or its approver changed, between C14 and C15).
var (
	judgeVerbSSH       = authority.JudgeVerbSSH
	judgeVerbAssertion = authority.JudgeVerbAssertion
	consumeVerbRequest = (*store.Credentials).ConsumeVerbRequest
)

// confirmAdmission is C8: approval.Admit on row at now, "" when admitted.
func confirmAdmission(row store.VerbRequestRow, now int64) string {
	stored, _ := hex.DecodeString(row.Digest) // the CHECK holds it to 64 characters; not hex: no bytes hash to it
	var digest [32]byte
	copy(digest[:], stored)
	switch approval.Admit(approval.AdmitInput{Now: now, CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt, Consumed: row.Outcome != "",
		StoredDigest: digest, BytesDigest: sha256.Sum256(row.Request)}) {
	case approval.Admitted:
		return ""
	case approval.RefuseConsumed:
		if row.Outcome == store.VerbOutcomeSuperseded {
			return "request " + row.RequestID + " was superseded by a newer finish of " + row.EpisodeID
		}
		return "request " + row.RequestID + " was already confirmed by " + row.ApprovedBy
	case approval.RefuseClockBehind:
		return "the clock is behind request " + row.RequestID + "'s creation time; refusing"
	case approval.RefuseExpired:
		return "request " + row.RequestID + " expired; run episode finish again for a new one"
	}
	return "request " + row.RequestID + "'s stored bytes do not match its digest; refusing"
}

// confirmStoreError answers a credentials-store failure: below schema 9 503 Restart, a busy writer
// 503 Busy with Retry-After (never a wait), anything else the 500 that asks for a retry.
func (d *Daemon) confirmStoreError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case isOlderStore(err):
		writeApprovalRefusal(w, http.StatusServiceUnavailable, "Restart",
			"the credentials store predates verb requests: restart `ailang-worldd serve` with this build once (it upgrades the store), then approve again")
	case store.IsWriterAlreadyActive(err):
		w.Header().Set("Retry-After", "1")
		writeApprovalRefusal(w, http.StatusServiceUnavailable, "Busy", "the credentials store is busy (a mint or a session verb is writing it); try again")
	default:
		d.writeInternalErrorLog(r, err)
		writeApprovalRefusal(w, http.StatusInternalServerError, "Internal", confirmIncomplete)
	}
}

// consumeRefusal is C15's answer: the row raced (409), the approver or the counter changed since
// C11 (422), else confirmStoreError.
func (d *Daemon) consumeRefusal(w http.ResponseWriter, r *http.Request, row store.VerbRequestRow, err error) {
	switch {
	case errors.Is(err, store.ErrVerbRequestSuperseded):
		writeApprovalRefusal(w, http.StatusConflict, "Conflict", "request "+row.RequestID+" was superseded by a newer finish of "+row.EpisodeID)
	case errors.Is(err, store.ErrVerbRequestConsumed), errors.Is(err, store.ErrVerbRequestUnknown), errors.Is(err, store.ErrVerbRequestExpired),
		errors.Is(err, store.ErrVerbRequestClockBehind), errors.Is(err, store.ErrVerbRequestDigestChanged):
		writeApprovalRefusal(w, http.StatusConflict, "Conflict", "request "+row.RequestID+" changed while it was confirmed: "+err.Error())
	case isApproverRace(err), isCounterRefusal(err):
		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the key changed since it was judged: "+err.Error())
	default:
		d.confirmStoreError(w, r, err)
	}
}

// handleSurfaceConfirm is POST /surface/confirm/{id} (design §2.6, steps C0-C15 in order).
func (d *Daemon) handleSurfaceConfirm(w http.ResponseWriter, r *http.Request, o []string) {
	if len(r.Header.Values("Authorization")) != 0 || len(r.Header.Values("Cookie")) != 0 { // C0
		writeApprovalRefusal(w, http.StatusForbidden, "Forbidden", "a confirmation is never made with a session or a cookie (the agent holds a session); confirm with your key on the card")
		return
	}
	id := r.PathValue("id")
	if !store.VerbRequestIDPattern.MatchString(id) { // C1
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", confirmIDReason)
		return
	}
	if origin := r.Header.Get("Origin"); !allowedApprovalOrigin(origin, o) { // C3
		writeApprovalRefusal(w, http.StatusForbidden, "CrossOrigin", crossOriginReason(origin))
		return
	}
	if !sameOriginFetch(r.Header.Get("Sec-Fetch-Site")) { // C4
		writeApprovalRefusal(w, http.StatusForbidden, "CrossOrigin", "Sec-Fetch-Site is not same-origin")
		return
	}
	if !jsonMediaType(r.Header.Get("Content-Type")) { // C5
		writeApprovalRefusal(w, http.StatusUnsupportedMediaType, "UnsupportedMediaType", "the body must be application/json")
		return
	}
	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, surfacepage.MaxApproveBodyBytes)) // C6
	if err != nil {
		writeApprovalRefusal(w, http.StatusRequestEntityTooLarge, "TooLarge", fmt.Sprintf("the body is over %d bytes", surfacepage.MaxApproveBodyBytes))
		return
	}
	body, err := surfacepage.ParseApproveBody(raw) // no typed arm: a typed or y key is unknown
	if err != nil {
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", err.Error())
		return
	}
	d.humanMu.Lock()
	defer d.humanMu.Unlock()
	ctx, cancel := context.WithTimeout(r.Context(), d.commitBudget)
	defer cancel()
	now := d.approveNow()
	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
	if err != nil {
		d.approvalInternal(w, r, err)
		return
	}
	defer func() { _ = ro.Close() }()
	row, ok, err := ro.VerbRequest(ctx, id) // C7
	switch {
	case err != nil:
		d.confirmStoreError(w, r, err)
		return
	case !ok:
		writeApprovalRefusal(w, http.StatusNotFound, "NotFound", "no verb request "+id)
		return
	}
	if reason := confirmAdmission(row, now); reason != "" { // C8
		writeApprovalRefusal(w, http.StatusConflict, "Conflict", reason)
		return
	}
	if body.Digest != row.Digest { // C9
		writeApprovalRefusal(w, http.StatusConflict, "Conflict", "the request changed since the card was shown; reload it")
		return
	}
	req, err := approval.DecodeVerb(row.Request) // C10
	if err == nil {
		err = req.Validate()
	}
	if err == nil {
		if reason := broker.FinishConfirmRefusal(broker.FinishConfirmFacts{Scope: req.Scope, Title: req.Title, Body: req.Body, Prompt: req.Prompt,
			Committer: req.Committer, Protected: req.Protected, RequestExpires: row.ExpiresAt, Now: now}); reason != "" {
			err = errors.New(reason)
		}
	}
	if err != nil {
		writeApprovalRefusal(w, http.StatusConflict, "NotConfirmable", "request "+id+" cannot be confirmed: "+err.Error()+"; confirm at the terminal")
		return
	}
	digest, challenge := sha256.Sum256(row.Request), [32]byte{} // C8 held: the bytes hash to the digest
	c, _ := hex.DecodeString(row.Challenge)                     // PutVerbRequest's grammar: 64 lowercase hex
	copy(challenge[:], c)
	in := authority.VerbJudgeInput{Digest: digest, Challenge: challenge, BytesDigest: digest, CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt, Now: now}
	judge, artifact, what := judgeVerbAssertion, []byte(body.WebAuthn), "passkey assertion"
	if body.SSHSIG != nil {
		judge, artifact, what = judgeVerbSSH, []byte(*body.SSHSIG), "SSH signature"
	}
	signed, err := judge(ctx, ro, in, artifact) // C11
	if err == nil {
		err = signed.Binds(digest, challenge, row.ExpiresAt) // C12
	}
	var je *authority.JudgeError
	switch {
	case errors.As(err, &je) && je.Step == authority.JudgeStepStore:
		d.confirmStoreError(w, r, err)
		return
	case err != nil:
		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the "+what+" does not confirm request "+id+": "+err.Error())
		return
	}
	wr, err := store.OpenCredentials(d.credentialsPath) // C13: non-waiting, before any write
	if err != nil {
		d.confirmStoreError(w, r, err)
		return
	}
	defer func() { _ = wr.Close() }() // held across C14 and C15, as handleSurfaceDecide holds keyArm's

	_, decision, err := broker.RecordOperatorApproval(ctx, d.store, req.Scope, confirmRequester, signed.ApprovedBy(), now) // C14
	if err != nil {
		d.writeInternalErrorLog(r, err)
		writeApprovalRefusal(w, http.StatusInternalServerError, "Internal", confirmIncomplete)
		return
	}
	name, format, key, _ := signed.Use()
	use := store.ApproverUse{Name: name, Format: format, PublicKey: key}
	if cid, count, ok := signed.PasskeyUse(); ok {
		use.Passkey = &store.PasskeyUse{CredentialID: cid, SignCount: int64(count)}
	}
	if err := consumeVerbRequest(wr, ctx, id, row.Digest, use, signed.ApprovedBy(), artifact, decision.String(), now); err != nil { // C15
		d.consumeRefusal(w, r, row, err)
		return
	}
	writeApprovalJSON(w, http.StatusOK, surfacepage.Confirmed(id, signed.ApprovedBy(), decision.String()))
}
