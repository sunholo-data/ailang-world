package approval

// Row 201 A4 (design_docs/planned/w-world-approvals-by-key-a4.md §2.3): the message a decision of
// a Human.Approve request signs. It is not the mint message (SignedMessage): a decision is approve
// or deny, so the verb is signed, and a distinct prefix and SSHSIG namespace keep a mint signature
// from deciding and a decision signature from minting. Pure: no I/O, no clock.

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
)

// DecisionPrefix starts every signed decision message (design §2.3).
const DecisionPrefix = "world-decision/v1"

// DecisionNamespace is the SSHSIG namespace every decision is signed under.
const DecisionNamespace = "world-decision@sunholo"

// DecisionMessageLen is the length of DecisionMessage: 17 + 1 + 1 + 32 + 32 + 8.
const DecisionMessageLen = len(DecisionPrefix) + 1 + 1 + 32 + 32 + 8

// The signed verb bytes; anything else is refused.
const (
	VerbApprove byte = 0x01
	VerbDeny    byte = 0x02
)

// ErrDecisionVerb is a verb byte other than VerbApprove and VerbDeny.
var ErrDecisionVerb = errors.New("approval: a decision verb is approve (0x01) or deny (0x02)")

// DecisionVerb is the signed byte of a decision word ("approve" or "deny").
func DecisionVerb(word string) (byte, error) {
	switch word {
	case "approve":
		return VerbApprove, nil
	case "deny":
		return VerbDeny, nil
	}
	return 0, ErrDecisionVerb
}

// DecisionMessage is the 91-byte message a decision signs (design §2.3):
// DecisionPrefix ‖ 0x00 ‖ verb ‖ packet digest P ‖ request digest Q ‖ expiry (uint64, big-endian).
func DecisionMessage(verb byte, packet, request [32]byte, expiry int64) ([]byte, error) {
	if verb != VerbApprove && verb != VerbDeny {
		return nil, ErrDecisionVerb
	}
	if expiry < 0 {
		return nil, ErrNegativeExpiry
	}
	m := make([]byte, 0, DecisionMessageLen)
	m = append(m, DecisionPrefix...)
	m = append(m, 0)
	m = append(m, verb)
	m = append(m, packet[:]...)
	m = append(m, request[:]...)
	m = binary.BigEndian.AppendUint64(m, uint64(expiry))
	return m, nil
}

// DecisionMessageHash is the SHA-256 of DecisionMessage: a passkey's WebAuthn challenge.
func DecisionMessageHash(verb byte, packet, request [32]byte, expiry int64) ([32]byte, error) {
	m, err := DecisionMessage(verb, packet, request, expiry)
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256(m), nil
}
