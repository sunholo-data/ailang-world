package main

// Row 201 A1b (design_docs/planned/w-world-approvals-by-key-a1b.md §2.7; sprint-3 plan §3 M5a):
// `approve <request-id> --signature <file>` approves a stored mint request with an enrolled
// approver's SSHSIG over the request's exact message, with no terminal.

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/sshsign"
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
	return approveArmored(ctx, credPath, "--signature "+path, armored, row, now)
}

// approveArmored judges an armored signature from source (--signature's file or --sign-with's
// signer) on a read-only handle by authority.JudgeSSHSignature (row 201 A3 §2.5: the parse F11, the
// lookup F12, F13, the law F14), the judge the minter and the daemon share.
func approveArmored(ctx context.Context, credPath, source string, armored []byte, row store.MintRequestRow, now int64) (*approval.SignedApproval, string) {
	ro, err := store.OpenCredentialsReadOnly(credPath)
	if err != nil {
		return nil, fmt.Sprintf("cannot read the credentials store %s: %v", credPath, err)
	}
	signed, err := authority.JudgeSSHSignature(ctx, ro, row, armored, now)
	_ = ro.Close()
	if err != nil {
		return nil, sshRefusal(err, source, row)
	}
	return signed, ""
}

// sshRefusal is the text of a JudgeSSHSignature refusal of source's signature (F11-F14).
func sshRefusal(err error, source string, row store.MintRequestRow) string {
	var je *authority.JudgeError
	if !errors.As(err, &je) {
		return err.Error()
	}
	switch je.Step {
	case authority.JudgeStepParse: // F11
		return fmt.Sprintf("%s is not an SSH signature: %v", source, je.Err)
	case authority.JudgeStepUnknownKey: // F13
		return fmt.Sprintf("the signature's key %s is not an enrolled approver; see ailang-worldd approver list", je.Fingerprint)
	case authority.JudgeStepLaw: // F14
		return fmt.Sprintf("the signature does not approve request %s: %v", row.RequestID, je.Err)
	}
	return je.Err.Error() // F12: below schema 5 names the serve restart
}

// ---- M5c ----

// writeApprovalMessage is approve --write-message (design §2.7 F5, F6): it writes the admitted
// request's 90-byte approval message to a NEW file (O_EXCL, mode 0644: the message is not secret)
// and prints the out-of-band signing line and the next step. It opens no terminal and writes no
// store row, so the request stays pending. It returns the refusal text, or "".
func writeApprovalMessage(stdout io.Writer, path string, row store.MintRequestRow, mr approval.MintRequest) string {
	msg, err := approvalMessage(row)
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

// approvalMessage is the admitted row's 90-byte approval message (design §2.3).
func approvalMessage(row store.MintRequestRow) ([]byte, error) {
	stored, _ := hex.DecodeString(row.Digest)       // admission has shown it is 32 bytes
	challenge, _ := hex.DecodeString(row.Challenge) // F4 has shown it is 32 bytes
	return approval.SignedMessage([32]byte(stored), [32]byte(challenge), row.ExpiresAt)
}

// ---- S1 (sprint 4) ----

// signForApproval is approve --sign-with's signing (design §2.7 step 3, §2.9): the admitted row's
// message, signed by env.sign under a new root bounded by the request's remaining lifetime. A
// request that expired since admission is refused before the signer runs (F7a), so the root is
// never made with a non-positive duration; a signer error is refused (F8). It returns the armored
// signature or the refusal text. The caller has checked env.sign (F7).
func signForApproval(env sessionEnv, key string, row store.MintRequestRow, stderr io.Writer) ([]byte, string) {
	rem := row.ExpiresAt - env.now()
	if rem <= 0 { // F7a
		return nil, expiredRefusal(row)
	}
	msg, err := approvalMessage(row)
	if err != nil {
		return nil, fmt.Sprintf("request %s: %v", row.RequestID, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(rem)*time.Second)
	defer cancel()
	armored, err := env.sign(ctx, key, msg, stderr)
	if err != nil { // F8
		return nil, fmt.Sprintf("--sign-with %s: %v; nothing was approved", key, err)
	}
	return armored, ""
}

// ---- S3 (sprint 4) ----

// signWithSSHKeygen is the production signer (design §2.9): ssh-keygen from PATH through
// host/sshsign's driven, environment-scrubbed site (AC10). PATH is not a trust boundary:
// whatever runs, approveArmored verifies its output against the enrolled key over World's own
// message.
func signWithSSHKeygen(ctx context.Context, key string, msg []byte, stderr io.Writer) ([]byte, error) {
	return sshsign.Sign(ctx, "ssh-keygen", key, msg, stderr)
}
