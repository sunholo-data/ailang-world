package surfacepage

// Row 243 V3 sprint 2 M4 (design_docs/planned/w-world-address-bar-verbs-v3.md §2.5, §6 CP0): the
// confirmation card's data, the answer of GET /surface/confirmations/{id} (sprint 3 registers that
// route and reads the row), built here from primitives without a daemon or a store: the stored
// canonical bytes and the raw challenge, from which the page computes the digest and the 86-byte
// message itself. There is no message-hash field and no artifact field.

import (
	"encoding/base64"

	"github.com/sunholo-data/ailang-world/host/approval"
)

// The states a confirmation card shows (design §2.5); only ConfirmationPending offers controls.
const (
	ConfirmationPending    = "pending"
	ConfirmationApproved   = "approved"
	ConfirmationSuperseded = "superseded"
	ConfirmationExpired    = "expired"
)

// ConfirmationInput is one stored verb request as a read route reads it: primitives only.
type ConfirmationInput struct {
	ID                   string // vr-<16 hex>
	Request              []byte // the stored canonical world/verb-request/v1 bytes
	Digest, Challenge    string // 64 hex each, as stored
	Created, Expires     int64
	Now                  int64
	Outcome              string // "" while unconsumed, else ConfirmationApproved or ConfirmationSuperseded
	ApprovedBy, Decision string // the approver and the world decision's ref, for an approved row
}

// Confirmation is GET /surface/confirmations/{id} (design §2.5).
type Confirmation struct {
	V          int            `json:"v"`
	ID         string         `json:"id"`
	Request    string         `json:"request"` // base64url of the stored canonical bytes
	Digest     string         `json:"digest"`
	Challenge  string         `json:"challenge"`
	Created    int64          `json:"created"`
	Expires    int64          `json:"expires"`
	Now        int64          `json:"now"`
	State      string         `json:"state"`
	ApprovedBy string         `json:"approvedBy"`
	Decision   string         `json:"decision"`
	Passkeys   []CardPasskeys `json:"passkeys"`
	SSH        CardSSH        `json:"ssh"`
}

// ConfirmationCard is in's card data, its active passkeys grouped by origin as RequestCard groups them.
func ConfirmationCard(in ConfirmationInput, passkeys []CardPasskey) Confirmation {
	c := Confirmation{V: 1, ID: in.ID, Request: base64.RawURLEncoding.EncodeToString(in.Request), Digest: in.Digest, Challenge: in.Challenge,
		Created: in.Created, Expires: in.Expires, Now: in.Now, State: confirmationState(in),
		Passkeys: RequestCard(StoredRequest{}, passkeys).Passkeys, SSH: CardSSH{Namespace: approval.VerbNamespace}}
	if c.State == ConfirmationApproved {
		c.ApprovedBy, c.Decision = in.ApprovedBy, in.Decision
	}
	return c
}

// confirmationState is the row's outcome if it has one, else expired at or after its expiry (as
// approval.Admit), else pending.
func confirmationState(in ConfirmationInput) string {
	if in.Outcome != "" {
		return in.Outcome
	}
	if in.Now >= in.Expires {
		return ConfirmationExpired
	}
	return ConfirmationPending
}

// ConfirmedJSON is POST /surface/confirm/{id}'s success (design §2.6): the world store's finish
// approval, decided by the key, is named by the request.
type ConfirmedJSON struct {
	V          int    `json:"v"`
	Confirmed  bool   `json:"confirmed"`
	Request    string `json:"request"`
	ApprovedBy string `json:"approvedBy"`
	Decision   string `json:"decision"`
	Next       string `json:"next"`
}

// Confirmed is the success body for request id confirmed by approvedBy with decision.
func Confirmed(id, approvedBy, decision string) ConfirmedJSON {
	return ConfirmedJSON{V: 1, Confirmed: true, Request: id, ApprovedBy: approvedBy, Decision: decision,
		Next: "the waiting episode finish opens the pull request"}
}
