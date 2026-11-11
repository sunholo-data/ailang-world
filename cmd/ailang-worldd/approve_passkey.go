package main

// Row 201 A2 sprint 3 (design_docs/planned/w-world-approvals-by-key-a2.md §2.13): `approve
// <request-id> --webauthn-link` prints the ceremony page's links for the request, and
// `--webauthn-assertion <file>` approves it with an enrolled passkey's assertion, with no terminal.
// Both run after admission and F4, under runApprove's 3 s read root.

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
)

// writePasskeyLinks is approve --webauthn-link (design §2.13): one get link per distinct origin
// among the active passkeys, each allowing that origin's credentials, then the next step. It opens
// no terminal and writes no row: the request stays pending. It returns the refusal text, or "".
func writePasskeyLinks(ctx context.Context, credPath string, stdout io.Writer, row store.MintRequestRow, mr approval.MintRequest) string {
	ro, err := store.OpenCredentialsReadOnly(credPath)
	if err != nil {
		return fmt.Sprintf("cannot read the credentials store %s: %v", credPath, err)
	}
	active, err := ro.ActivePasskeys(ctx)
	_ = ro.Close()
	if err != nil { // below schema 6 names the serve restart
		return err.Error()
	}
	if len(active) == 0 { // WF1
		return "no passkey approver is enrolled; see ailang-worldd approver list"
	}
	h, err := requestHash(row)
	if err != nil {
		return fmt.Sprintf("request %s: %v", row.RequestID, err)
	}
	hash16 := hex.EncodeToString(h[:8])
	fmt.Fprintf(stdout, "approve request %s (expires %s, message hash %s) with a passkey: open the link for its origin, check the page shows the same request, expiry and hash, and save the file it offers:\n",
		row.RequestID, approverTime(row.ExpiresAt), hash16)
	links := map[string]*ceremonyLink{}
	var order []string
	for _, e := range active { // ActivePasskeys is ordered by name
		l := links[e.Passkey.Origin]
		if l == nil {
			l = &ceremonyLink{V: 1, Op: "get", RP: e.Passkey.RPID, Origin: e.Passkey.Origin, Request: row.RequestID, Expires: row.ExpiresAt,
				Hash16: hash16, Challenge: base64.RawURLEncoding.EncodeToString(h[:])}
			links[e.Passkey.Origin] = l
			order = append(order, e.Passkey.Origin)
		}
		l.Allow = append(l.Allow, base64.RawURLEncoding.EncodeToString(e.Passkey.CredentialID))
	}
	for _, o := range order {
		fmt.Fprintf(stdout, "  %s\n", links[o].link())
	}
	fmt.Fprintf(stdout, "then approve it:\n  ailang-worldd approve %s --webauthn-assertion <the saved file> --db %s --credentials-db %s\n",
		row.RequestID, shellQuote(mr.WorldDB), shellQuote(mr.CredentialsDB))
	return ""
}

// requestHash is the admitted row's message hash, the challenge a passkey signs (design §2.1).
func requestHash(row store.MintRequestRow) ([32]byte, error) {
	stored, _ := hex.DecodeString(row.Digest)       // admission has shown it is 32 bytes
	challenge, _ := hex.DecodeString(row.Challenge) // F4 has shown it is 32 bytes
	return approval.MessageHash([32]byte(stored), [32]byte(challenge), row.ExpiresAt)
}

// approveAssertion is approve --webauthn-assertion's check, before any network call (design §2.13
// WF2-WF6): the file is read once through a descriptor check (WF2) and parsed (WF3), the passkey
// its credential ID names is looked up on a read-only handle (WF4, WF5), and the assertion law
// judges it against the stored row and that passkey (WF6). It returns the verified approval, or
// the refusal text. The caller has already admitted row at now and decoded its challenge (F4).
func approveAssertion(ctx context.Context, credPath, path string, row store.MintRequestRow, now int64) (*approval.SignedApproval, string) {
	file, err := readRegularFileCapped(path, approval.MaxPasskeyResponseBytes)
	if err != nil { // WF2
		return nil, fmt.Sprintf("--webauthn-assertion: %v", err)
	}
	resp, err := approval.ParseAssertionResponse(file)
	if err != nil { // WF3
		return nil, fmt.Sprintf("--webauthn-assertion %s is not a passkey assertion: %v", path, err)
	}
	ro, err := store.OpenCredentialsReadOnly(credPath)
	if err != nil {
		return nil, fmt.Sprintf("cannot read the credentials store %s: %v", credPath, err)
	}
	a, p, ok, err := ro.PasskeyByCredentialID(ctx, resp.RawID)
	_ = ro.Close()
	if err != nil { // WF4: below schema 6 names the serve restart
		return nil, err.Error()
	}
	if !ok { // WF5
		return nil, fmt.Sprintf("the assertion's credential cred %s is not an enrolled passkey; see ailang-worldd approver list", passkeyCred(resp.RawID))
	}
	key, err := approval.ParseStoredCOSEKey(a.PublicKey)
	if err != nil { // WF5: only a hand-edited row
		return nil, fmt.Sprintf("approver %q's stored passkey key is not valid (%v); see ailang-worldd approver list", a.Name, err)
	}
	stored, _ := hex.DecodeString(row.Digest)       // admission has shown it is 32 bytes
	challenge, _ := hex.DecodeString(row.Challenge) // F4 has shown it is 32 bytes
	signed, err := approval.ApproveByAssertion(approval.AssertionInput{
		Admit: approval.AdmitInput{Now: now, CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt, Consumed: row.Consumed,
			StoredDigest: [32]byte(stored), BytesDigest: sha256.Sum256(row.Request)},
		Challenge: [32]byte(challenge),
		Response:  file,
		Approver: approval.Approver{Name: a.Name, Format: a.Format, EnrolledAt: a.EnrolledAt, LastUsedAt: a.LastUsedAt, RevokedAt: a.RevokedAt,
			Passkey: &approval.Passkey{RPID: p.RPID, Origin: p.Origin, CredentialID: p.CredentialID, Key: key,
				BackupEligible: p.BackupEligible, SignCount: uint32(p.SignCount)}},
	})
	if err != nil { // WF6
		return nil, fmt.Sprintf("the passkey assertion does not approve request %s: %v", row.RequestID, err)
	}
	return signed, ""
}
