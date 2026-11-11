package main

// Row 201 A2 sprint 3 (design_docs/planned/w-world-approvals-by-key-a2.md §2.12): `approver enroll
// --format webauthn`. The browser's registration reaches the CLI as a file; the operator types, at
// /dev/tty, the 26-character enrolment code the page showed (a fingerprint of that file, not a
// secret), then y. No store handle is open while either prompt waits.

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
)

// passkeyCeremonyBudget bounds the enrolment code's read, counted from its prompt (WN6): the
// operator opens a browser, passes Face ID or Touch ID (or uses a phone), and saves a file. A
// constant, not a flag, so a script cannot widen it.
const passkeyCeremonyBudget = 300 * time.Second

// approverEnrolPasskey is the passkey enrol write, a test seam (CW6 replaces it to inject a race);
// production never reassigns it.
var approverEnrolPasskey = (*store.Credentials).EnrolPasskey

// ceremonyLink is the ceremony page's fragment (design §2.12 step 6, §2.13): what one create or
// get asks. The page reads it from location.hash; it never reaches a server.
type ceremonyLink struct {
	V         int      `json:"v"`
	Op        string   `json:"op"`
	RP        string   `json:"rp"`
	Origin    string   `json:"origin"`
	Name      string   `json:"name,omitempty"`
	Request   string   `json:"request,omitempty"`
	Expires   int64    `json:"expires,omitempty"`
	Hash16    string   `json:"hash16,omitempty"`
	Challenge string   `json:"challenge"`
	Allow     []string `json:"allow,omitempty"`
}

// link is <origin>/surface/passkey#<base64url(JSON)>.
func (c ceremonyLink) link() string {
	b, _ := json.Marshal(c) // strings, integers and a string slice: Marshal cannot fail
	return c.Origin + "/surface/passkey#" + base64.RawURLEncoding.EncodeToString(b)
}

// runApproverEnrollPasskey is approver enroll --format webauthn (design §2.12 WN2-WN11 with N3,
// N6, N7 and N9-N14), after runApproverEnroll's WN1.
func runApproverEnrollPasskey(name, origin, regPath, db, creds string, stdout, stderr io.Writer, env sessionEnv) int {
	const prefix, nothing = "ailang-worldd approver enroll: ", "nothing was enrolled"
	if name == "" || origin == "" || regPath == "" { // WN2
		return approverRefuse(stderr, prefix, "--name, --origin and --registration are required with --format webauthn")
	}
	credPath, refusal := approverStores(db, creds) // N3
	if refusal != "" {
		return approverRefuse(stderr, prefix, "%s", refusal)
	}
	rpID, err := approval.ParseCeremonyOrigin(origin)
	if err != nil { // WN3
		return approverRefuse(stderr, prefix, "--origin %q: %v", origin, err)
	}
	if _, err := os.Lstat(regPath); !errors.Is(err, fs.ErrNotExist) { // WN4
		return approverRefuse(stderr, prefix, "--registration %s already exists or cannot be checked; pass a new path: the browser saves the new registration there", regPath)
	}
	pre, cancelPre := context.WithTimeout(context.Background(), approveReadBudget)
	refusal = passkeyEnrolPrecheck(pre, credPath, name, nil) // N6, N7, closed before any terminal
	cancelPre()
	if refusal != "" {
		return approverRefuse(stderr, prefix, "%s", refusal)
	}
	var challenge [32]byte
	if env.random == nil || env.random(challenge[:]) != nil || challenge == [32]byte{} { // WN5
		fmt.Fprintf(stderr, "%sinternal error: no fresh challenge (sessionEnv.random is missing, failed, or gave 32 zero bytes); %s\n", prefix, nothing)
		return exitFatal
	}
	term, code := openApproverTerminal(env, stderr, prefix, nothing) // N9a, N9: before the link, which only a person can use
	if code != exitOK {
		return code
	}
	defer func() { _ = term.Close() }()
	fmt.Fprintf(stdout, "open this link in a browser (the daemon serving %s shows the page; nothing is sent to it), create the passkey, and save the file it offers as %s:\n  %s\n",
		origin, regPath, ceremonyLink{V: 1, Op: "create", RP: rpID, Origin: origin, Name: name, Challenge: base64.RawURLEncoding.EncodeToString(challenge[:])}.link())
	line, ok, complete := askApproverLine(term, "Enrolment code shown by the browser (XXXXX-XXXXX-XXXXX-XXXXX-XXXXXX): ", passkeyCeremonyBudget, env.after)
	if !ok { // WN6
		return approverRefuse(stderr, prefix, "no code within %s; %s", passkeyCeremonyBudget, nothing)
	}
	if !complete { // WN7: end of input, a read error or a partial line
		return approverRefuse(stderr, prefix, "aborted (no complete code line); %s. fix: run it again and type the code the page shows", nothing)
	}
	reg, err := readRegularFileCapped(regPath, approval.MaxPasskeyResponseBytes)
	if err != nil { // WN8
		return approverRefuse(stderr, prefix, "--registration: %v; %s", err, nothing)
	}
	a, err := approval.NewPasskeyApprover(approval.RegistrationInput{Name: name, Origin: origin, Challenge: challenge, Response: reg, Now: env.now()})
	if err != nil { // WN9
		return approverRefuse(stderr, prefix, "%s does not enrol: %v; %s", regPath, err, nothing)
	}
	resp, _ := approval.ParseRegistrationResponse(reg) // R3 has parsed it
	typed := strings.NewReplacer("-", "", " ", "", "\t", "", "\n", "").Replace(strings.ToUpper(line))
	if subtle.ConstantTimeCompare([]byte(typed), []byte(approval.EnrolmentCode(resp.ClientDataJSON, resp.AttestationObject))) != 1 { // WN10: one attempt
		return approverRefuse(stderr, prefix, "the code you typed is not the code of %s; %s", regPath, nothing)
	}
	pre, cancelPre = context.WithTimeout(context.Background(), approveReadBudget)
	refusal = passkeyEnrolPrecheck(pre, credPath, name, &a) // WN11
	cancelPre()
	if refusal != "" {
		return approverRefuse(stderr, prefix, "%s; %s", refusal, nothing)
	}
	if code := confirmOnTerminal(term, env, stderr, prefix, approval.RenderPasskeyEnrolment(a)+fmt.Sprintf("Enrol approver %s? [y/N] ", name), nothing); code != exitOK {
		return code // N10a, N10
	}

	// The write, under the credentials writer lock taken only now (N11-N13).
	ctx, cancel := context.WithTimeout(context.Background(), approveReadBudget)
	defer cancel()
	st, err := openCredentialsForSession(ctx, db, credPath)
	if err != nil { // N11
		fmt.Fprintf(stderr, "%scannot open the credentials store for writing: %v\n", prefix, err)
		return exitFatal
	}
	defer func() { _ = st.Close() }()
	p := a.Passkey
	err = approverEnrolPasskey(st, ctx, store.ApproverRow{Name: a.Name, Format: a.Format, PublicKey: p.Key.Canonical(), EnrolledAt: a.EnrolledAt},
		store.PasskeyRow{Name: a.Name, RPID: p.RPID, Origin: p.Origin, CredentialID: p.CredentialID, AAGUID: p.AAGUID[:], Alg: p.Key.Alg,
			BackupEligible: p.BackupEligible, SignCount: int64(p.SignCount)})
	switch {
	case errors.Is(err, store.ErrApproverNameTaken), errors.Is(err, store.ErrApproverKeyEnrolled), errors.Is(err, store.ErrPasskeyCredentialEnrolled): // N12
		return approverRefuse(stderr, prefix, "%v (another enrol ran first); %s", err, nothing)
	case err != nil: // N13
		fmt.Fprintf(stderr, "%scannot enrol approver %q: %v\n", prefix, a.Name, err)
		return exitFatal
	}
	fmt.Fprintf(stdout, "enrolled approver %s: %s; approvals by it record approved_by=webauthn:%s\n", a.Name, passkeyLabel(p.Key.Alg, p.RPID, p.CredentialID), a.Name) // N14
	return exitOK
}

// passkeyEnrolPrecheck is enrol's read-only refusal on its own handle: "" when the store is at
// schema 6 (N6) and no approver holds the name (N7) and, given the verified registration a,
// neither its key nor its credential ID (WN11). Revoked rows count: none is ever reused.
func passkeyEnrolPrecheck(ctx context.Context, credPath, name string, a *approval.Approver) string {
	ro, err := store.OpenCredentialsReadOnly(credPath)
	if err != nil {
		return fmt.Sprintf("cannot read the credentials store %s: %v", credPath, err)
	}
	defer func() { _ = ro.Close() }()
	if _, _, err := ro.Passkey(ctx, name); err != nil { // N6: below schema 6 names the serve restart
		return err.Error()
	}
	_, taken, err := ro.Approver(ctx, name)
	if err != nil {
		return err.Error()
	}
	if taken { // N7
		return fmt.Sprintf("an approver named %q already exists (a name is never reused, even after revoke); fix: choose another --name", name)
	}
	if a == nil {
		return ""
	}
	holder, enrolled, err := ro.ApproverByKey(ctx, a.Passkey.Key.Canonical())
	if err != nil {
		return err.Error()
	}
	if enrolled { // WN11: the key
		return fmt.Sprintf("this passkey's key is already enrolled as approver %q (a key is never enrolled twice, even after revoke)", holder.Name)
	}
	holder, _, enrolled, err = ro.PasskeyByCredentialID(ctx, a.Passkey.CredentialID)
	if err != nil {
		return err.Error()
	}
	if enrolled { // WN11: the credential
		return fmt.Sprintf("credential cred %s is already enrolled as approver %q (a credential is never enrolled twice, even after revoke)", passkeyCred(a.Passkey.CredentialID), holder.Name)
	}
	return ""
}
