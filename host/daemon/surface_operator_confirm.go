package daemon

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

func operatorCeremony(ctx context.Context, db *store.Credentials, row store.OperatorActionRow, href string) (string, error) {
	active, err := db.ActivePasskeys(ctx)
	if err != nil {
		return "", err
	}
	var keys []surfacepage.CardPasskey
	for _, e := range active {
		keys = append(keys, surfacepage.CardPasskey{Origin: e.Passkey.Origin, RPID: e.Passkey.RPID, CredentialID: e.Passkey.CredentialID})
	}
	m, err := approval.OperatorActionMessage([32]byte(row.Digest), [32]byte(row.Challenge), row.ExpiresAt)
	if err != nil {
		return "", err
	}
	passkeys := surfacepage.DecisionPasskeys(keys)
	wire, err := json.MarshalIndent(struct {
		Post      string                     `json:"post"`
		Digest    string                     `json:"digest"`
		Message   string                     `json:"messageHex"`
		Namespace string                     `json:"namespace"`
		Passkeys  []surfacepage.CardPasskeys `json:"passkeys"`
	}{href, hex.EncodeToString(row.Digest), hex.EncodeToString(m), approval.OperatorActionNamespace, passkeys}, "", "  ")
	return string(wire), err
}

type operatorConfirmationAnswer struct {
	operatorStatus
	Result     string `json:"result"`
	ApprovedBy string `json:"approvedBy"`
}

func (d *Daemon) handleOperatorConfirm(w http.ResponseWriter, r *http.Request, id string, origins []string) {
	raw, ok := readOperatorPOST(w, r, origins, surfacepage.MaxApproveBodyBytes)
	if !ok {
		return
	}
	var fields map[string]json.RawMessage
	if approval.DecodeOperatorObject(raw, &fields, "v", "digest", "webauthn", "sshsig") != nil || len(fields) != 3 {
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "send a v1 digest and one key artifact")
		return
	}
	body, err := surfacepage.ParseApproveBody(raw)
	if err != nil || (body.WebAuthn != nil && body.WebAuthn[0] != '{') {
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "send a v1 digest and one key artifact")
		return
	}
	artifact := []byte(body.WebAuthn)
	if body.SSHSIG != nil {
		if strings.HasPrefix(strings.TrimSpace(*body.SSHSIG), "{") {
			writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "send an armored SSH signature in sshsig")
			return
		}
		artifact = []byte(*body.SSHSIG)
	}

	d.humanMu.Lock()
	defer d.humanMu.Unlock()
	ctx, cancel := context.WithTimeout(r.Context(), d.commitBudget)
	defer cancel()
	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
	if err != nil {
		d.operatorError(w, err)
		return
	}
	defer ro.Close()
	now := d.approveNow()
	row, err := ro.OperatorAction(ctx, id, now)
	if err != nil {
		d.operatorError(w, err)
		return
	}
	a, err := approval.DecodeOperatorAction(row.Action)
	digest := sha256.Sum256(row.Action)
	served, pathErr := store.CanonicalPath(d.credentialsPath)
	if pathErr != nil {
		d.operatorError(w, pathErr)
		return
	}
	supported := a.Kind == "start" && a.Options.Headless || a.Kind == "finish" && a.Options.ConfirmOnSurface
	if err != nil || !supported || a.ID != id || a.Episode != row.Episode || a.Kind != row.Kind || a.Store != served ||
		!bytes.Equal(row.Digest, digest[:]) || body.Digest != hex.EncodeToString(row.Digest) || len(row.Challenge) != 32 {
		d.operatorError(w, store.ErrOperatorConflict)
		return
	}
	if row.ApprovedBy != "" {
		if !bytes.Equal(row.Artifact, artifact) {
			d.operatorError(w, store.ErrOperatorConflict)
			return
		}
		writeApprovalJSON(w, http.StatusOK, operatorConfirmationAnswer{operatorStatusFor(row, now), "already-confirmed", row.ApprovedBy})
		return
	}
	if row.State == "expired" || now < row.CreatedAt {
		operatorConfirmExpired(w)
		return
	}
	if row.State != "requested" {
		d.operatorError(w, store.ErrOperatorConflict)
		return
	}
	proof, err := authority.JudgeOperatorAction(ctx, ro, authority.OperatorJudgeInput{
		Action: a, Digest: [32]byte(row.Digest), BytesDigest: digest, Challenge: [32]byte(row.Challenge),
		CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt, Now: now,
	}, artifact)
	var judge *authority.JudgeError
	switch {
	case errors.Is(err, approval.ErrOperatorLifetime):
		operatorConfirmExpired(w)
		return
	case errors.As(err, &judge) && judge.Step == authority.JudgeStepStore:
		d.operatorError(w, err)
		return
	case err != nil:
		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the enrolled key did not confirm this exact action")
		return
	}

	wr, err := d.openConfirmationWriter()
	if err != nil {
		d.operatorError(w, err)
		return
	}
	defer wr.Close()
	now = d.approveNow()
	row, err = wr.ConfirmOperatorAction(ctx, store.OperatorCAS{ID: id, Digest: digest, Generation: row.Generation}, artifact, proof, now)
	switch {
	case errors.Is(err, store.ErrOperatorExpired):
		operatorConfirmExpired(w)
	case isApproverRace(err), isCounterRefusal(err):
		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the enrolled key changed during confirmation")
	case err != nil:
		d.operatorError(w, err)
	default:
		writeApprovalJSON(w, http.StatusOK, operatorConfirmationAnswer{operatorStatusFor(row, now), "confirmed", row.ApprovedBy})
	}
}

func operatorConfirmExpired(w http.ResponseWriter) {
	writeApprovalRefusal(w, http.StatusConflict, "Expired", "the action is outside its signed lifetime; review again")
}
