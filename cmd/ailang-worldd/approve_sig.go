package main

// Row 201 A1b (design_docs/planned/w-world-approvals-by-key-a1b.md §2.7; sprint-3 plan §3 M5a):
// `approve <request-id> --signature <file>` approves a stored mint request with an enrolled
// approver's SSHSIG over the request's exact message, with no terminal.

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
)

// approveMaxSignatureBytes bounds the --signature file (design §2.7 F9): an Ed25519 or P-256
// SSHSIG is under 1 KiB.
const approveMaxSignatureBytes = 16 << 10

// approveSignature is approve --signature's check, before any network call (design §2.7 F9-F14):
// the file is read once through a descriptor check (F9, F10), parsed (F11), its key looked up
// among the enrolled approvers on a read-only handle under approve's 3 s root (F12, F13), and
// judged by the signature law against the stored row (F14). It returns the verified approval, or
// the refusal text. The caller has already admitted row at now and decoded its challenge (F4).
func approveSignature(ctx context.Context, credPath, path string, row store.MintRequestRow, now int64) (*approval.SignedApproval, string) {
	// F9: one read, through a check of the open descriptor (a FIFO or a device is refused at once).
	armored, err := readRegularFileCapped(path, approveMaxSignatureBytes)
	if err != nil { // F10
		return nil, fmt.Sprintf("--signature: %v", err)
	}
	sig, err := approval.ParseSSHSIG(armored)
	if err != nil { // F11
		return nil, fmt.Sprintf("--signature %s is not an SSH signature: %v", path, err)
	}
	ro, err := store.OpenCredentialsReadOnly(credPath)
	if err != nil {
		return nil, fmt.Sprintf("cannot read the credentials store %s: %v", credPath, err)
	}
	enrolled, ok, err := ro.ApproverByKey(ctx, sig.PublicKey.Blob)
	_ = ro.Close()
	if err != nil { // F12: below schema 5 names the serve restart
		return nil, err.Error()
	}
	if !ok { // F13
		return nil, fmt.Sprintf("the signature's key %s is not an enrolled approver; see ailang-worldd approver list", sig.PublicKey.Fingerprint())
	}
	stored, _ := hex.DecodeString(row.Digest)       // admission has shown it is 32 bytes
	challenge, _ := hex.DecodeString(row.Challenge) // F4 has shown it is 32 bytes
	signed, err := approval.ApproveBySignature(approval.SignatureInput{
		Admit: approval.AdmitInput{Now: now, CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt, Consumed: row.Consumed,
			StoredDigest: [32]byte(stored), BytesDigest: sha256.Sum256(row.Request)},
		Challenge: [32]byte(challenge),
		Signature: armored,
		Approver: approval.Approver{Name: enrolled.Name, Format: enrolled.Format, Key: approval.SSHPublicKey{Blob: enrolled.PublicKey},
			EnrolledAt: enrolled.EnrolledAt, LastUsedAt: enrolled.LastUsedAt, RevokedAt: enrolled.RevokedAt},
	})
	if err != nil { // F14
		return nil, fmt.Sprintf("the signature does not approve request %s: %v", row.RequestID, err)
	}
	return signed, ""
}

// ---- M5c ----

// writeApprovalMessage is approve --write-message (design §2.7 F5, F6): it writes the admitted
// request's 90-byte approval message to a NEW file (O_EXCL, mode 0644: the message is not secret)
// and prints the out-of-band signing line and the next step. It opens no terminal and writes no
// store row, so the request stays pending. It returns the refusal text, or "".
func writeApprovalMessage(stdout io.Writer, path string, row store.MintRequestRow, mr approval.MintRequest) string {
	stored, _ := hex.DecodeString(row.Digest)       // admission has shown it is 32 bytes
	challenge, _ := hex.DecodeString(row.Challenge) // F4 has shown it is 32 bytes
	msg, err := approval.SignedMessage([32]byte(stored), [32]byte(challenge), row.ExpiresAt)
	if err != nil {
		return fmt.Sprintf("request %s: %v", row.RequestID, err)
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil { // F5: an existing file is never overwritten
		return fmt.Sprintf("--write-message: %v; nothing was written", err)
	}
	_, err = f.Write(msg)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		_ = os.Remove(path)
		return fmt.Sprintf("--write-message: %v; nothing was written", err)
	}
	fmt.Fprintf(stdout, "wrote request %s's %d-byte approval message to %s\n", row.RequestID, len(msg), path)
	fmt.Fprintf(stdout, "sign it with an enrolled key (the request expires at %s):\n  ssh-keygen -Y sign -n %s -f <your key> < %s > %s.sig\n",
		approverTime(row.ExpiresAt), approval.Namespace, shellQuote(path), shellQuote(path))
	fmt.Fprintf(stdout, "then approve it:\n  ailang-worldd approve %s --signature %s.sig --db %s --credentials-db %s\n",
		row.RequestID, shellQuote(path), shellQuote(mr.WorldDB), shellQuote(mr.CredentialsDB))
	return ""
}
