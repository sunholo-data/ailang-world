package approval

// Row 201 A1b (design_docs/planned/w-world-approvals-by-key-a1b.md §2.2): the pure enrolment
// law. An approver is a name and one SSH public key, enrolled at the daemon's terminal (a typed
// y); its signature may later approve a mint request (sprint 2). No I/O here (coding standard S2).

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
)

const (
	// FormatSSH is an SSH key approver (A1b).
	FormatSSH = "ssh"
	// FormatWebAuthn is a passkey approver (row 201 A2): it enrols only through NewPasskeyApprover.
	FormatWebAuthn = "webauthn"
	// MaxKeyFileBytes caps a --pubkey file (one authorized_keys line).
	MaxKeyFileBytes = 8 << 10
)

// ApproverNamePattern is an approver's name: lowercase, digits, '.', '_' and '-', 1 to 48
// characters, starting with a letter or digit. It is printed on prompts and stored as ssh:<name>.
var ApproverNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,47}$`)

// Approver is one enrolled approver as the pure law sees it.
type Approver struct {
	Name       string
	Format     string       // FormatSSH or FormatWebAuthn
	Key        SSHPublicKey // Type + wire Blob (ssh approvers)
	EnrolledAt int64
	LastUsedAt int64    // 0: never used
	RevokedAt  int64    // 0: active
	Passkey    *Passkey // webauthn approvers only (row 201 A2); nil for ssh
}

var (
	ErrApproverName        = errors.New("approval: approver name must match ^[a-z0-9][a-z0-9._-]{0,47}$")
	ErrUseRegistration     = errors.New("approval: webauthn approvers enrol from a browser registration (approver enroll --format webauthn --origin … --registration …)")
	ErrApproverFormat      = errors.New("approval: approver format must be ssh")
	ErrApproverKeyFile     = errors.New("approval: the public key file must hold exactly one key line of at most 8 KiB")
	ErrPrivateKeyFile      = errors.New("approval: that is a PRIVATE key file; pass the .pub file (its contents were not shown or kept)")
	ErrApproverRevoked     = errors.New("approval: the approver is revoked")
	ErrApproverClockBehind = errors.New("approval: the clock is behind the approver's enrolment time")
)

// NewApprover is the enrolment law: it parses one public key line into an Approver, or refuses.
// The checks run in a fixed order (E1-E7); the first that fails is returned.
func NewApprover(name, format string, keyFile []byte, now int64) (Approver, error) {
	if !ApproverNamePattern.MatchString(name) { // E1
		return Approver{}, ErrApproverName
	}
	if format == FormatWebAuthn { // E2 (before E3, so webauthn names its own path)
		return Approver{}, ErrUseRegistration
	}
	if format != FormatSSH { // E3
		return Approver{}, ErrApproverFormat
	}
	if len(keyFile) == 0 || len(keyFile) > MaxKeyFileBytes { // E4
		return Approver{}, ErrApproverKeyFile
	}
	if bytes.HasPrefix(bytes.TrimLeft(keyFile, " \t\r\n"), []byte("-----BEGIN")) { // E5: never echoed
		return Approver{}, ErrPrivateKeyFile
	}
	line := bytes.TrimSuffix(keyFile, []byte("\n"))
	if bytes.ContainsAny(line, "\n") { // E6
		return Approver{}, ErrApproverKeyFile
	}
	key, err := ParseAuthorizedKey(line)
	if err != nil { // E7: ErrKeyType (sk-*, rsa, p384), ErrAuthorizedKey, ErrBase64
		return Approver{}, err
	}
	return Approver{Name: name, Format: FormatSSH, Key: key, EnrolledAt: now}, nil
}

// CheckRevoke is the revocation law, run before the y/N (the store re-checks it in SQL).
func CheckRevoke(a Approver, now int64) error {
	if a.RevokedAt != 0 { // R1
		return ErrApproverRevoked
	}
	if now < a.EnrolledAt { // R2
		return ErrApproverClockBehind
	}
	return nil
}

// ---- M1c ----

// EnrolmentWarning is the R6 text every enrolment y/N shows, verbatim (design §2.2).
const EnrolmentWarning = "This key will approve mint requests with no terminal. A signature proves that something holding\n" +
	"the private key signed. It does not prove that you were there. Enrol only a key that cannot sign\n" +
	"without you: one behind a passphrase, in an agent that asks before each use (`ssh-add -c`), or in\n" +
	"a Secure Enclave or hardware holder (for example Secretive). Compare the fingerprint above with\n" +
	"`ssh-keygen -l -f <your key>.pub` or `ssh-add -l` before you type y.\n"

// RenderEnrolment is the body of the enrolment y/N: what will be stored, and what it proves. It
// has no branches (AP5 pins it; RC4 checks it shows every persisted field).
func RenderEnrolment(a Approver) string {
	return fmt.Sprintf("approver  %s\nformat    %s\nkey       %s %s\nrecords   approved_by=%s:%s\n\n%s",
		a.Name, a.Format, a.Key.Type, a.Key.Fingerprint(), a.Format, a.Name, EnrolmentWarning)
}

// ---- sprint 2, M3b ----

// KeyFromBlob is the key a stored SSH wire blob holds, with ParseAuthorizedKey's refusals and no
// line: `approver list` and `approver revoke` show its type (sprint-2 plan D3).
func KeyFromBlob(blob []byte) (SSHPublicKey, error) {
	km, err := parseKeyBlob(blob)
	if err != nil {
		return SSHPublicKey{}, err
	}
	return SSHPublicKey{Type: km.typ, Blob: bytes.Clone(blob)}, nil
}
