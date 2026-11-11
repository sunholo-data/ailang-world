package main

// Row 201 A2 sprint 3 (design_docs/planned/w-world-approvals-by-key-a2.md §2.13): `approve
// <request-id> --webauthn-link` prints the ceremony page's links for the request. It runs after
// admission and F4, under runApprove's 3 s read root.

import (
	"context"
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
