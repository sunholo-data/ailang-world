package surfacepage

// Row 201 A3 sprint 2 (design_docs/planned/w-world-approvals-by-key-a3.md §2.3, §2.4): the JSON
// the approval routes speak, built here without a daemon or a store. GET /surface/requests.json
// lists the pending requests (Needs you) and the decided ones; GET /surface/requests/{id} is the
// card's data (the stored bytes and the raw nonce, never a message hash); POST
// /surface/approve/{id} takes ApproveBody and answers ApproveRecorded or ApproveRefusal.

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
)

// MaxDecidedListed is the decision list's length (design §2.3, Q3).
const MaxDecidedListed = 20

// MaxApproveBodyBytes bounds a POST /surface/approve/{id} body (design §2.4 R6).
const MaxApproveBodyBytes = 20 << 10

// StoredRequest is one mint_requests row as the routes read it (store.MintRequestRow's fields).
type StoredRequest struct {
	ID, Episode, Digest, Challenge string
	Request                        []byte
	Created, Expires, ConsumedAt   int64
	Consumed                       bool
	ApprovedBy, CredentialID       string
}

// PendingRequest is one Needs-you line: a summary decoded from the stored bytes, or, when they do
// not decode or do not hash to the digest, Refused and no summary (never hidden).
type PendingRequest struct {
	ID       string  `json:"id"`
	Episode  string  `json:"episode"`
	Agent    string  `json:"agent,omitempty"`
	Model    string  `json:"model,omitempty"`
	Grants   int     `json:"grants,omitempty"`
	TTL      *int64  `json:"ttl_s,omitempty"`
	Expires  int64   `json:"expires"`
	Hash16   string  `json:"hash16,omitempty"`
	Recorded *string `json:"recorded"` // the verified approver of a recorded approval, or null
	Refused  string  `json:"refused,omitempty"`
}

// DecidedRequest is one line of the decision list.
type DecidedRequest struct {
	ID           string `json:"id"`
	Episode      string `json:"episode"`
	ApprovedBy   string `json:"approvedBy"`
	At           int64  `json:"at"`
	Credential16 string `json:"credential16"`
}

// RequestsJSON is GET /surface/requests.json. It carries no request bytes, challenge, artifact,
// token or path.
type RequestsJSON struct {
	V           int              `json:"v"`
	Now         int64            `json:"now"`
	Pending     []PendingRequest `json:"pending"`
	Decided     []DecidedRequest `json:"decided"`
	Unavailable string           `json:"unavailable"`
	// Decisions are the pending Human.Approve requests a decision card decides (row 201 A4 §2.8).
	Decisions          []PendingDecision `json:"decisions"`
	DecisionsTruncated bool              `json:"decisionsTruncated"`
}

// hex32 decodes 64 hex characters.
func hex32(s string) ([32]byte, bool) {
	b, err := hex.DecodeString(s)
	if err != nil || len(b) != 32 {
		return [32]byte{}, false
	}
	return [32]byte(b), true
}

// decodeStored is a stored request's bytes as a canonical request whose SHA-256 is its digest,
// and its message hash (approval.MessageHash, what a signer signs).
func decodeStored(r StoredRequest) (approval.MintRequest, [32]byte, error) {
	digest, ok := hex32(r.Digest)
	if !ok || sha256.Sum256(r.Request) != digest {
		return approval.MintRequest{}, [32]byte{}, errors.New("the stored bytes do not match the request's digest")
	}
	mr, err := approval.Decode(r.Request)
	if err != nil {
		return approval.MintRequest{}, [32]byte{}, fmt.Errorf("the stored bytes are not a canonical mint request (%v)", err)
	}
	challenge, ok := hex32(r.Challenge)
	if !ok {
		return approval.MintRequest{}, [32]byte{}, errors.New("the stored challenge is not 64 hex")
	}
	h, err := approval.MessageHash(digest, challenge, r.Expires)
	if err != nil {
		return approval.MintRequest{}, [32]byte{}, err
	}
	return mr, h, nil
}

// Pending is r's Needs-you line; recorded is the verified approver of its recorded approval, or "".
func Pending(r StoredRequest, recorded string) PendingRequest {
	p := PendingRequest{ID: r.ID, Episode: r.Episode, Expires: r.Expires}
	if recorded != "" {
		p.Recorded = &recorded
	}
	mr, h, err := decodeStored(r)
	if err != nil {
		p.Refused = err.Error()
		return p
	}
	ttl := mr.TTLSeconds
	p.Agent, p.Model, p.Grants, p.TTL, p.Hash16 = mr.Agent, mr.Model, len(mr.Grants), &ttl, hex.EncodeToString(h[:8])
	return p
}

// Decided is r's line in the decision list.
func Decided(r StoredRequest) DecidedRequest {
	c := r.CredentialID
	if len(c) > 16 {
		c = c[:16]
	}
	return DecidedRequest{ID: r.ID, Episode: r.Episode, ApprovedBy: r.ApprovedBy, At: r.ConsumedAt, Credential16: c}
}

// CardPasskey is one active passkey as the card offers it.
type CardPasskey struct {
	Origin, RPID string
	CredentialID []byte
}

// CardPasskeys is the active passkeys of one origin: its rpID and their credential IDs (base64url).
type CardPasskeys struct {
	Origin string   `json:"origin"`
	RP     string   `json:"rp"`
	Allow  []string `json:"allow"`
}

// CardSSH is what an SSH signer needs besides the message: the namespace.
type CardSSH struct {
	Namespace string `json:"namespace"`
}

// Card is GET /surface/requests/{id}: the stored bytes and the raw nonce, from which the page
// computes the digest and the message hash itself. There is no message-hash field.
type Card struct {
	V          int            `json:"v"`
	ID         string         `json:"id"`
	Request    string         `json:"request"` // base64url of the stored canonical bytes
	Digest     string         `json:"digest"`
	Challenge  string         `json:"challenge"`
	Created    int64          `json:"created"`
	Expires    int64          `json:"expires"`
	Consumed   bool           `json:"consumed"`
	ApprovedBy string         `json:"approvedBy"`
	Passkeys   []CardPasskeys `json:"passkeys"`
	SSH        CardSSH        `json:"ssh"`
}

// RequestCard is r's card, with the active passkeys grouped by origin in the order given.
func RequestCard(r StoredRequest, passkeys []CardPasskey) Card {
	c := Card{V: 1, ID: r.ID, Request: base64.RawURLEncoding.EncodeToString(r.Request), Digest: r.Digest, Challenge: r.Challenge,
		Created: r.Created, Expires: r.Expires, Consumed: r.Consumed, ApprovedBy: r.ApprovedBy,
		Passkeys: []CardPasskeys{}, SSH: CardSSH{Namespace: approval.Namespace}}
	at := map[string]int{}
	for _, p := range passkeys {
		i, ok := at[p.Origin]
		if !ok {
			i = len(c.Passkeys)
			at[p.Origin] = i
			c.Passkeys = append(c.Passkeys, CardPasskeys{Origin: p.Origin, RP: p.RPID, Allow: []string{}})
		}
		c.Passkeys[i].Allow = append(c.Passkeys[i].Allow, base64.RawURLEncoding.EncodeToString(p.CredentialID))
	}
	return c
}

// ApproveBody is a POST /surface/approve/{id} body: v 1, the digest the page displayed, and
// exactly one of a WebAuthn assertion envelope (A2 §2.7) and an armored SSHSIG.
type ApproveBody struct {
	V        int             `json:"v"`
	Digest   string          `json:"digest"`
	WebAuthn json.RawMessage `json:"webauthn"`
	SSHSIG   *string         `json:"sshsig"`
}

// ErrApproveBody refuses a body that is not ApproveBody's grammar (R6).
var ErrApproveBody = errors.New("the body is not {\"v\":1,\"digest\":<64 hex>} with exactly one of webauthn and sshsig")

// ParseApproveBody is R6's grammar: strict JSON (no unknown field, exactly one value), v 1,
// exactly one format, the digest 64 lowercase hex.
func ParseApproveBody(b []byte) (ApproveBody, error) {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	var body ApproveBody
	if err := dec.Decode(&body); err != nil {
		return ApproveBody{}, fmt.Errorf("%w: %v", ErrApproveBody, err)
	}
	if _, err := dec.Token(); err != io.EOF {
		return ApproveBody{}, fmt.Errorf("%w: more than one JSON value", ErrApproveBody)
	}
	switch {
	case body.V != 1:
		return ApproveBody{}, fmt.Errorf("%w: v is %d", ErrApproveBody, body.V)
	case (body.WebAuthn != nil) == (body.SSHSIG != nil):
		return ApproveBody{}, fmt.Errorf("%w: it has %s", ErrApproveBody, map[bool]string{true: "both", false: "neither"}[body.SSHSIG != nil])
	case !lowerHex64(body.Digest):
		return ApproveBody{}, fmt.Errorf("%w: digest is not 64 lowercase hex", ErrApproveBody)
	}
	return body, nil
}

// lowerHex64 is exactly 64 lowercase hex characters.
func lowerHex64(s string) bool {
	b, ok := hex32(s)
	return ok && hex.EncodeToString(b[:]) == s
}

// AdmissionRefusal is approve's admission (approval.Admit) of r at now, in the CLI's words
// (cmd/ailang-worldd approveAdmission); "" when r is admitted (R8).
func AdmissionRefusal(r StoredRequest, now int64) string {
	if _, err := approval.Decode(r.Request); err != nil {
		return fmt.Sprintf("request %s is not a canonical mint request (%v); fix: make a new request with episode start --request", r.ID, err)
	}
	digestRefusal := fmt.Sprintf("request %s's stored bytes do not match its digest; refusing; fix: make a new request with episode start --request", r.ID)
	stored, _ := hex32(r.Digest) // not 64 hex: the zero array, which no bytes hash to (the digest refusal)
	switch approval.Admit(approval.AdmitInput{Now: now, CreatedAt: r.Created, ExpiresAt: r.Expires, Consumed: r.Consumed,
		StoredDigest: stored, BytesDigest: sha256.Sum256(r.Request)}) {
	case approval.Admitted:
		return ""
	case approval.RefuseConsumed:
		return fmt.Sprintf("request %s was already used at %s by %s (credential %s); fix: make a new request with episode start --request",
			r.ID, rfc3339(r.ConsumedAt), r.ApprovedBy, r.CredentialID)
	case approval.RefuseClockBehind:
		return fmt.Sprintf("the clock is behind request %s's creation time (%s); refusing", r.ID, rfc3339(r.Created))
	case approval.RefuseExpired:
		return fmt.Sprintf("request %s expired at %s; fix: make a new request with episode start --request", r.ID, rfc3339(r.Expires))
	default:
		return digestRefusal
	}
}

func rfc3339(t int64) string { return time.Unix(t, 0).UTC().Format(time.RFC3339) }

// ApproveRefusal is every refusal of the approval routes: inert JSON, nothing written.
type ApproveRefusal struct {
	V      int    `json:"v"`
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

// ApproveRecorded is a POST's success: the verified approval is recorded for the minter.
type ApproveRecorded struct {
	V          int    `json:"v"`
	Recorded   bool   `json:"recorded"`
	Request    string `json:"request"`
	ApprovedBy string `json:"approvedBy"`
	Next       string `json:"next"`
}

// Recorded is the success body for request id approved by approvedBy.
func Recorded(id, approvedBy string) ApproveRecorded {
	return ApproveRecorded{V: 1, Recorded: true, Request: id, ApprovedBy: approvedBy,
		Next: "the minter (approve " + id + " --await) mints it; it then appears under Decided"}
}
