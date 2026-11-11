package authority

// Row 201 A3 sprint 1 M2 (design_docs/planned/w-world-approvals-by-key-a3.md §2.5): one judge for
// the CLI and the daemon. approve --webauthn-assertion and --signature (cmd/ailang-worldd), the
// minter approve --await, and from sprint 2 the surface's POST all judge an artifact here: the
// lookups on a read-only credentials handle, then the A1b/A2 law, unchanged.

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
)

// The steps a JudgeError names, in the order the judges run them.
const (
	JudgeStepRequest    = "request"     // the row's digest or challenge is not 64 hex
	JudgeStepParse      = "parse"       // the artifact is not an assertion or an SSHSIG
	JudgeStepStore      = "store"       // the lookup failed (below the schema its format needs)
	JudgeStepUnknownKey = "unknown-key" // no enrolled approver holds the artifact's key or credential
	JudgeStepStoredKey  = "stored-key"  // the enrolled passkey's stored key does not parse
	JudgeStepLaw        = "law"         // the law refused it
)

// ErrNotEnrolled is an unknown-key JudgeError's Err.
var ErrNotEnrolled = errors.New("authority: no enrolled approver holds that key")

// JudgeError is a refused judgement: the step that refused, what the CLI names for it, and the
// cause. The CLI prints its own text per step (cmd/ailang-worldd passkeyRefusal, sshRefusal).
type JudgeError struct {
	Step        string
	RawID       []byte // unknown-key, assertion: the credential ID the response names
	Fingerprint string // unknown-key, SSHSIG: the signature key's fingerprint
	Approver    string // stored-key: the approver whose stored key is invalid
	Err         error
}

func (e *JudgeError) Error() string { return fmt.Sprintf("authority: judge: %s: %v", e.Step, e.Err) }

func (e *JudgeError) Unwrap() error { return e.Err }

// judgeAdmit is the law's admission input and the challenge for row at now. The CLI has admitted
// row and checked its challenge before (F4); the daemon checks them too (R8, R10). A row that is
// not 64 hex in either is refused here, never converted.
func judgeAdmit(row store.MintRequestRow, now int64) (approval.AdmitInput, [32]byte, error) {
	stored, err := hex.DecodeString(row.Digest)
	if err != nil || len(stored) != 32 {
		return approval.AdmitInput{}, [32]byte{}, &JudgeError{Step: JudgeStepRequest, Err: fmt.Errorf("request %s's stored digest is not 64 hex", row.RequestID)}
	}
	challenge, err := hex.DecodeString(row.Challenge)
	if err != nil || len(challenge) != 32 {
		return approval.AdmitInput{}, [32]byte{}, &JudgeError{Step: JudgeStepRequest, Err: fmt.Errorf("request %s's stored challenge is not 64 hex", row.RequestID)}
	}
	return approval.AdmitInput{Now: now, CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt, Consumed: row.Consumed,
		StoredDigest: [32]byte(stored), BytesDigest: sha256.Sum256(row.Request)}, [32]byte(challenge), nil
}

// JudgeAssertion judges a passkey assertion response for row at now: ParseAssertionResponse, the
// passkey its credential ID names (PasskeyByCredentialID on ro), ParseStoredCOSEKey, then
// approval.ApproveByAssertion. It returns the verified approval or a *JudgeError.
func JudgeAssertion(ctx context.Context, ro *store.Credentials, row store.MintRequestRow, response []byte, now int64) (*approval.SignedApproval, error) {
	admit, challenge, err := judgeAdmit(row, now)
	if err != nil {
		return nil, err
	}
	resp, err := approval.ParseAssertionResponse(response)
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepParse, Err: err}
	}
	a, p, ok, err := ro.PasskeyByCredentialID(ctx, resp.RawID)
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepStore, Err: err}
	}
	if !ok {
		return nil, &JudgeError{Step: JudgeStepUnknownKey, RawID: resp.RawID, Err: ErrNotEnrolled}
	}
	key, err := approval.ParseStoredCOSEKey(a.PublicKey)
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepStoredKey, Approver: a.Name, Err: err}
	}
	signed, err := approval.ApproveByAssertion(approval.AssertionInput{
		Admit:     admit,
		Challenge: challenge,
		Response:  response,
		Approver: approval.Approver{Name: a.Name, Format: a.Format, EnrolledAt: a.EnrolledAt, LastUsedAt: a.LastUsedAt, RevokedAt: a.RevokedAt,
			Passkey: &approval.Passkey{RPID: p.RPID, Origin: p.Origin, CredentialID: p.CredentialID, Key: key,
				BackupEligible: p.BackupEligible, SignCount: uint32(p.SignCount)}},
	})
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepLaw, Err: err}
	}
	return signed, nil
}

// JudgeSSHSignature judges an armored SSHSIG for row at now: ParseSSHSIG, the approver its key
// names (ApproverByKey on ro), then approval.ApproveBySignature. It returns the verified approval
// or a *JudgeError.
func JudgeSSHSignature(ctx context.Context, ro *store.Credentials, row store.MintRequestRow, armored []byte, now int64) (*approval.SignedApproval, error) {
	admit, challenge, err := judgeAdmit(row, now)
	if err != nil {
		return nil, err
	}
	sig, err := approval.ParseSSHSIG(armored)
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepParse, Err: err}
	}
	enrolled, ok, err := ro.ApproverByKey(ctx, sig.PublicKey.Blob)
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepStore, Err: err}
	}
	if !ok {
		return nil, &JudgeError{Step: JudgeStepUnknownKey, Fingerprint: sig.PublicKey.Fingerprint(), Err: ErrNotEnrolled}
	}
	signed, err := approval.ApproveBySignature(approval.SignatureInput{
		Admit:     admit,
		Challenge: challenge,
		Signature: armored,
		Approver: approval.Approver{Name: enrolled.Name, Format: enrolled.Format, Key: approval.SSHPublicKey{Blob: enrolled.PublicKey},
			EnrolledAt: enrolled.EnrolledAt, LastUsedAt: enrolled.LastUsedAt, RevokedAt: enrolled.RevokedAt},
	})
	if err != nil {
		return nil, &JudgeError{Step: JudgeStepLaw, Err: err}
	}
	return signed, nil
}
