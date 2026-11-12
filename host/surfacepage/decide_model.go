package surfacepage

// Row 201 A4 sprint 1 (design_docs/planned/w-world-approvals-by-key-a4.md §2.8): the decision
// routes' JSON shapes and the POST body's grammar, pure. Sprint 2 admits exactly one of typed and
// sshsig; webauthn arrives with sprint 3's passkey arm, so until then it is an unknown field.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// DecisionSubject is GET /surface/decisions/{ref}: the stored bytes of the request and of the one
// packet carrying it (never a precomputed message or hash), the window an expiry must fall in, and
// the law's display verdict at now (the POST runs the law again).
type DecisionSubject struct {
	V            int            `json:"v"`
	Request      string         `json:"request"`
	RequestBytes string         `json:"requestBytes"` // base64url, no padding
	Packet       string         `json:"packet"`
	PacketBytes  string         `json:"packetBytes"` // base64url, no padding
	Entry        int64          `json:"entry"`
	Proposal     string         `json:"proposal"`
	Now          int64          `json:"now"`
	Window       DecisionWindow `json:"window"`
	Refusal      string         `json:"refusal"`
	Reserved     string         `json:"reserved"`
	DecidedBy    string         `json:"decidedBy"`
	Decision     string         `json:"decision"`
	SSH          DecisionSSH    `json:"ssh"`
}

// DecisionWindow is the expiry window of a decision made now: min = now + 1, max = min(now + 600, deadline).
type DecisionWindow struct {
	Min int64 `json:"min"`
	Max int64 `json:"max"`
}

// DecisionSSH names the SSHSIG namespace a decision is signed under.
type DecisionSSH struct {
	Namespace string `json:"namespace"`
}

// DecisionMade is the POST's success: the verb, the request, the decision object and who decided.
type DecisionMade struct {
	V         int    `json:"v"`
	Decided   string `json:"decided"`
	Request   string `json:"request"`
	Decision  string `json:"decision"`
	DecidedBy string `json:"decidedBy"`
}

// DecideBody is the POST body (design §2.8): the verb, the packet digest P and request digest Q
// the caller saw (64 lowercase hex), the expiry, and exactly one format: typed: true (sprint 1) or
// sshsig, an armored SSH signature over the decision message (sprint 2).
type DecideBody struct {
	V       int     `json:"v"`
	Verb    string  `json:"verb"`
	Packet  string  `json:"packet"`
	Request string  `json:"request"`
	Expiry  int64   `json:"expiry"`
	Typed   bool    `json:"typed,omitempty"`
	SSHSIG  *string `json:"sshsig,omitempty"`
}

// ErrDecideBody refuses a body that is not DecideBody's grammar.
var ErrDecideBody = errors.New(`the body is not {"v":1,"verb":…,"packet":<64 hex>,"request":<64 hex>,"expiry":…} with exactly one of "typed":true and "sshsig":<armored SSHSIG>`)

// ParseDecideBody is the POST's grammar: strict JSON (no unknown field, so no decidedBy; exactly
// one value), v 1, exactly one format, both digests 64 lowercase hex. The verb and the expiry are
// the law's (D9, D10), not the grammar's.
func ParseDecideBody(b []byte) (body DecideBody, packet, request [32]byte, err error) {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		return DecideBody{}, packet, request, fmt.Errorf("%w: %v", ErrDecideBody, err)
	}
	if _, err := dec.Token(); err != io.EOF {
		return DecideBody{}, packet, request, fmt.Errorf("%w: more than one JSON value", ErrDecideBody)
	}
	switch {
	case body.V != 1:
		return DecideBody{}, packet, request, fmt.Errorf("%w: v is %d", ErrDecideBody, body.V)
	case !body.Typed && body.SSHSIG == nil:
		return DecideBody{}, packet, request, fmt.Errorf("%w: it has no format", ErrDecideBody)
	case body.Typed && body.SSHSIG != nil:
		return DecideBody{}, packet, request, fmt.Errorf("%w: it has two formats", ErrDecideBody)
	case !lowerHex64(body.Packet) || !lowerHex64(body.Request):
		return DecideBody{}, packet, request, fmt.Errorf("%w: packet and request are each 64 lowercase hex", ErrDecideBody)
	}
	packet, _ = hex32(body.Packet)
	request, _ = hex32(body.Request)
	return body, packet, request, nil
}
