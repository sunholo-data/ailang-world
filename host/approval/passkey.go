package approval

// Row 201 A2 (design §2.8): the registration law. A passkey approver is enrolled from one browser
// registration (fmt "none": nothing is attested), checked against the enrol process's own
// challenge and origin, and tied to the operator by the enrolment code he types (a 128-bit
// fingerprint of the exact registration bytes) and his typed y. No I/O.

import (
	"bytes"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

// Passkey is an enrolled passkey: what an assertion must match.
type Passkey struct {
	RPID, Origin   string
	CredentialID   []byte
	Key            COSEKey
	AAGUID         [16]byte
	BackupEligible bool   // stored; an assertion must agree (W13)
	BackupState    bool   // shown at enrolment only (synced or device-bound); not stored
	SignCount      uint32 // the registration's counter
}

// RegistrationInput is what NewPasskeyApprover judges.
type RegistrationInput struct {
	Name, Origin string
	Challenge    [32]byte // the enrol process's fresh challenge
	Response     []byte   // the registration file
	Now          int64
}

var (
	ErrClientDataType       = errors.New("approval: clientDataJSON has the wrong type for this ceremony")
	ErrChallenge            = errors.New("approval: the response was made for another challenge")
	ErrOriginMismatch       = errors.New("approval: the response was made at another origin")
	ErrCrossOrigin          = errors.New("approval: the response was made in a cross-origin frame")
	ErrAttestationObject    = errors.New("approval: the attestationObject is not a map of exactly fmt, attStmt and authData")
	ErrAttestationFormat    = errors.New(`approval: the registration must use attestation "none" with an empty attStmt`)
	ErrRPIDHash             = errors.New("approval: the authenticator signed for another relying party ID")
	ErrUserPresence         = errors.New("approval: the authenticator did not report user presence (UP)")
	ErrUserVerification     = errors.New("approval: the authenticator did not verify the user (UV: Face ID, Touch ID, a PIN)")
	ErrCredentialIDMismatch = errors.New("approval: the response's rawId is not the passkey's credential ID")
)

func b64Challenge(c []byte) string { return base64.RawURLEncoding.EncodeToString(c) }

// NewPasskeyApprover is the registration law (R1-R15, in this order; the first that fails is
// returned). R16 (the algorithm) is ParseCOSEKey's Q3, reached through R11.
func NewPasskeyApprover(in RegistrationInput) (Approver, error) {
	if !ApproverNamePattern.MatchString(in.Name) { // R1
		return Approver{}, ErrApproverName
	}
	rpID, err := ParseCeremonyOrigin(in.Origin) // R2
	if err != nil {
		return Approver{}, err
	}
	resp, err := ParseRegistrationResponse(in.Response) // R3
	if err != nil {
		return Approver{}, err
	}
	cd, err := ParseClientData(resp.ClientDataJSON) // R4
	if err != nil {
		return Approver{}, err
	}
	if cd.Type != "webauthn.create" { // R5
		return Approver{}, ErrClientDataType
	}
	if cd.Challenge != b64Challenge(in.Challenge[:]) { // R6
		return Approver{}, ErrChallenge
	}
	if cd.Origin != in.Origin { // R7
		return Approver{}, ErrOriginMismatch
	}
	if cd.CrossOrigin || cd.HasTopOrigin { // R8
		return Approver{}, ErrCrossOrigin
	}
	raw, err := attestedAuthData(resp.AttestationObject) // R9, R10
	if err != nil {
		return Approver{}, err
	}
	ad, err := ParseAuthData(raw, true) // R11
	if err != nil {
		return Approver{}, err
	}
	if ad.RPIDHash != sha256.Sum256([]byte(rpID)) { // R12
		return Approver{}, ErrRPIDHash
	}
	if !ad.UP { // R13
		return Approver{}, ErrUserPresence
	}
	if !ad.UV { // R14
		return Approver{}, ErrUserVerification
	}
	if !bytes.Equal(resp.RawID, ad.CredentialID) { // R15
		return Approver{}, ErrCredentialIDMismatch
	}
	return Approver{Name: in.Name, Format: FormatWebAuthn, EnrolledAt: in.Now, Passkey: &Passkey{
		RPID: rpID, Origin: in.Origin, CredentialID: ad.CredentialID, Key: ad.Key, AAGUID: ad.AAGUID,
		BackupEligible: ad.BE, BackupState: ad.BS, SignCount: ad.SignCount}}, nil
}

// attestedAuthData checks the attestationObject's shape (R9) and format (R10) and returns its
// authData.
func attestedAuthData(b []byte) ([]byte, error) {
	v, err := decodeCBORExact(b, attestationLimits)
	if err != nil { // R9
		return nil, fmt.Errorf("%w: %w", ErrAttestationObject, err)
	}
	f, okF := v.get(cborValue{kind: cborText, b: []byte("fmt")})
	st, okS := v.get(cborValue{kind: cborText, b: []byte("attStmt")})
	ad, okA := v.get(cborValue{kind: cborText, b: []byte("authData")})
	if len(v.items) != 6 || !okF || !okS || !okA || f.kind != cborText || st.kind != cborMap || ad.kind != cborBytes { // R9
		return nil, ErrAttestationObject
	}
	if string(f.b) != "none" || len(st.items) != 0 { // R10
		return nil, ErrAttestationFormat
	}
	return ad.b, nil
}

// EnrolmentCode is the 26-character, 128-bit fingerprint of a registration (design §2.8 [r1]):
// base32 (A-Z, 2-7, unpadded) of the first 16 bytes of SHA-256("world-enrolment-code/v1" || 0x00
// || SHA-256(clientDataJSON) || SHA-256(attestationObject)). It is public, not a secret.
func EnrolmentCode(clientDataJSON, attestationObject []byte) string {
	c, a := sha256.Sum256(clientDataJSON), sha256.Sum256(attestationObject)
	h := sha256.Sum256(append(append([]byte("world-enrolment-code/v1\x00"), c[:]...), a[:]...))
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(h[:16])
}

// GroupEnrolmentCode shows an EnrolmentCode as XXXXX-XXXXX-XXXXX-XXXXX-XXXXXX.
func GroupEnrolmentCode(code string) string {
	return code[:5] + "-" + code[5:10] + "-" + code[10:15] + "-" + code[15:20] + "-" + code[20:]
}

// PasskeyEnrolmentWarning is what every passkey enrolment y/N shows, verbatim (design §2.8).
const PasskeyEnrolmentWarning = "This passkey will approve mint requests with no terminal, each time after your device verifies\n" +
	"you (Face ID, Touch ID, a PIN or a security-key touch). World cannot check which device holds\n" +
	"it; the 26-character enrolment code you typed, which matched this exact registration, is what ties\n" +
	"it to the browser you just used. Enrol only a passkey you created yourself just now.\n"

// RenderPasskeyEnrolment is the body of a passkey enrolment y/N; a is NewPasskeyApprover's.
func RenderPasskeyEnrolment(a Approver) string {
	p := a.Passkey
	backup := "device-bound"
	switch {
	case p.BackupEligible && p.BackupState:
		backup = "synced passkey"
	case p.BackupEligible:
		backup = "backup-eligible"
	}
	cred := sha256.Sum256(p.CredentialID)
	return fmt.Sprintf("approver  %s\nformat    %s\npasskey   %s for %s at %s\ncred      %s\naaguid    %s\nbackup    %s\ncounter   %d\nrecords   approved_by=%s:%s\n\n%s",
		a.Name, a.Format, p.Key.AlgName(), p.RPID, p.Origin, hex.EncodeToString(cred[:8]), hex.EncodeToString(p.AAGUID[:]),
		backup, p.SignCount, a.Format, a.Name, PasskeyEnrolmentWarning)
}
