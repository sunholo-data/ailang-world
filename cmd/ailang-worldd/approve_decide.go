package main

// Row 201 A4 sprint 2 (design_docs/planned/w-world-approvals-by-key-a4.md §2.10): `approve
// <request-ref> [--deny]` decides a pending Human.Approve request whose decision packet the
// coordinator committed. The CLI never opens the world store: it reads the subject from the
// daemon serving --db (GET /surface/decisions/{ref}), checks both byte fields against their refs
// and the packet's requestRef here, shows every field, and POSTs the decision to
// /surface/decide/{ref}, where the daemon runs the decision law under humanMu and writes. The typed
// arm writes the row-236 operator proof beside the served credentials store only after a y typed
// at /dev/tty.

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/sshsign"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// decisionRefPattern is a Human.Approve request ref: approve's decision mode (design §2.10).
var decisionRefPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

// approveDecideBudget bounds the decision mode's reads before the y/N, and its POST after it.
const approveDecideBudget = 10 * time.Second

// decideOptions are approve's flags in the decision mode.
type decideOptions struct {
	ref, verb, db, credsDB, daemon, signWith string
	sigPath, msgPath, assertion              string // the other key flags (CD6, CD7, CD9)
	passkeyLink                              bool   // CD8
	expiry                                   int64  // with sigPath or assertion: the signed expiry
}

// decisionSubject is a subject the CLI checked: its request and packet, decoded from bytes that
// hash to their refs.
type decisionSubject struct {
	surfacepage.DecisionSubject
	request broker.ApprovalRequest
	packet  coordinator.DecisionPacket
}

// checkDecisionSubject decodes the daemon's GET answer for ref and checks it here (design §2.8:
// the daemon sends bytes, never a message or a hash): the request bytes hash to ref, the packet
// bytes hash to the packet ref, both decode strictly, and the packet carries ref. It returns the
// refusal text, or "".
func checkDecisionSubject(ref string, body []byte) (decisionSubject, string) {
	var s decisionSubject
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&s.DecisionSubject); err != nil {
		return s, fmt.Sprintf("the daemon's answer for request %s is not a decision subject (%v); nothing was decided", ref, err)
	}
	mismatch := fmt.Sprintf("the daemon's bytes for request %s do not match their refs; refusing before any prompt (nothing was decided)", ref)
	reqBytes, err1 := base64.RawURLEncoding.DecodeString(s.RequestBytes)
	pktBytes, err2 := base64.RawURLEncoding.DecodeString(s.PacketBytes)
	if err1 != nil || err2 != nil || sha256Ref(reqBytes) != ref || sha256Ref(pktBytes) != s.Packet {
		return s, mismatch
	}
	var err error
	if s.request, err = broker.DecodeApprovalRequest(reqBytes); err != nil {
		return s, mismatch
	}
	if s.packet, err = coordinator.DecodeDecisionPacket(pktBytes); err != nil || s.packet.RequestRef.String() != ref {
		return s, mismatch
	}
	return s, ""
}

// refDigest is the 32 bytes of a sha256: ref the subject check has held to its bytes.
func refDigest(ref string) [32]byte {
	d, _ := hex.DecodeString(strings.TrimPrefix(ref, "sha256:"))
	return [32]byte(d)
}

// sha256Ref is the sha256: ref of b.
func sha256Ref(b []byte) string { return "sha256:" + finishSHA(b) }

// renderDecision is the typed y/N: every field of both signed objects (design §2.8's field
// table), the entry and the proposal's page, then the verb in capitals and the deadline.
func renderDecision(ref, verb string, s decisionSubject) string {
	r, p := s.request, s.packet
	return fmt.Sprintf("request %s\n  effect               %s\n  scope                %s\n  cost                 %d\n  requester            %s\n  now                  %d\n"+
		"packet %s (log entry %d)\n  packetHash           %s\n  proposalHash         %s (/workbench?object=%s)\n  requestRef           %s\n"+
		"  createdAt            %d\n  deadlineAt           %d (%s)\n  escalationsRemaining %d\n  policy               %s\n"+
		"%s this request? A request is decided once; the decision is final. [y/N] ",
		ref, r.Effect, r.Scope, r.Cost, r.Requester, r.Now, s.Packet, s.Entry, p.PacketHash, p.ProposalHash, p.ProposalHash, p.RequestRef,
		p.CreatedAt, p.DeadlineAt, approverTime(p.DeadlineAt), p.EscalationsRemaining, p.Policy, strings.ToUpper(verb))
}

// refusalReason is the reason of a daemon refusal body, or the body itself.
func refusalReason(body []byte) string {
	var r surfacepage.ApproveRefusal
	if json.Unmarshal(body, &r) == nil && r.Reason != "" {
		return r.Reason
	}
	return strings.TrimSpace(string(body))
}

// runApproveDecide is approve's decision mode (design §2.10). The daemon is --daemon, or the
// home service's address; it must serve --db and the credentials store beside it (the proof is
// written there). Every refusal before the POST decides nothing.
func runApproveDecide(o decideOptions, stdout, stderr io.Writer, env sessionEnv, refuse func(string, ...any) int) int {
	credPath, err := credentialsFor(o.db, o.credsDB)
	if err != nil {
		return refuse("%v", err)
	}
	if o.daemon == "" {
		o.daemon = defaultServiceAddr()
	}
	base := clientBaseURL(o.daemon)
	hostPort := strings.TrimSuffix(strings.TrimPrefix(base, "http://"), "/")
	h, err := fetchDaemonHealth(hostPort)
	if err != nil {
		return refuse("%v; nothing was decided; fix: start the daemon serving %s, or name it with --daemon <host:port>", err, o.db)
	}
	if same, _ := store.SameDatabase(h.DBPath, o.db); !same {
		return refuse("the daemon on %s serves store %s, not %s; nothing was decided; fix: --daemon <host:port> of the daemon serving it", hostPort, h.DBPath, o.db)
	}
	if same, _ := store.SameDatabase(h.CredentialsDB, credPath); !same {
		return refuse("the daemon on %s serves credentials store %s, not %s; nothing was decided (the operator proof is written beside the served one)", hostPort, h.CredentialsDB, credPath)
	}
	c := newClient(base)
	ctx, cancel := context.WithTimeout(context.Background(), approveDecideBudget)
	defer cancel()
	status, body, err := c.do(ctx, http.MethodGet, "/surface/decisions/"+o.ref, nil)
	if err != nil {
		return refuse("%v; nothing was decided", err)
	}
	if status != http.StatusOK {
		return refuse("the daemon will not show request %s (HTTP %d): %s", o.ref, status, refusalReason(body))
	}
	s, refusal := checkDecisionSubject(o.ref, body)
	if refusal != "" {
		return refuse("%s", refusal)
	}
	if s.Refusal != "" {
		return refuse("request %s cannot be decided now: %s", o.ref, s.Refusal)
	}
	if o.signWith != "" {
		return decideBySignWith(o, s, c, base, env, stdout, stderr, refuse)
	}
	if o.sigPath != "" || o.msgPath != "" || o.assertion != "" || o.passkeyLink {
		return decideByArtifact(ctx, o, credPath, hostPort, s, c, base, stdout, refuse)
	}
	const nothing = "nothing was decided; decide with an enrolled key instead (--sign-with <key> or --signature <file>)"
	term, code := openApproverTerminal(env, stderr, "ailang-worldd approve: ", nothing)
	if code != exitOK {
		return code
	}
	code = confirmOnTerminal(term, env, stderr, "ailang-worldd approve: ", renderDecision(o.ref, o.verb, s), "nothing was decided")
	_ = term.Close()
	if code != exitOK {
		return code
	}
	proof, err := writeOperatorProof(credPath, "decision-"+strings.TrimPrefix(o.ref, "sha256:"), "decide-"+o.verb)
	if err != nil {
		return refuse("cannot write the operator proof beside %s: %v; nothing was decided", credPath, err)
	}
	return postDecision(c, base, o.ref, surfacepage.DecideBody{V: 1, Verb: o.verb, Packet: strings.TrimPrefix(s.Packet, "sha256:"),
		Request: strings.TrimPrefix(o.ref, "sha256:"), Expiry: s.Window.Max, Typed: true}, proof, stdout, refuse)
}

// decideBySignWith is approve <request-ref> --sign-with (design §2.10, CD5): the 91-byte decision
// message for the window's latest expiry, signed by env.signDecision (ssh-keygen -Y sign under
// DecisionNamespace) under a root bounded by that expiry, and posted as sshsig. The daemon judges
// it (design §2.7); no terminal opens and no proof is written.
func decideBySignWith(o decideOptions, s decisionSubject, c *client, base string, env sessionEnv, stdout, stderr io.Writer, refuse func(string, ...any) int) int {
	if env.signDecision == nil {
		fmt.Fprintln(stderr, "ailang-worldd approve: internal error: no signer (sessionEnv.signDecision)")
		return exitFatal
	}
	verb, _ := approval.DecisionVerb(o.verb)
	body := surfacepage.DecideBody{V: 1, Verb: o.verb, Packet: strings.TrimPrefix(s.Packet, "sha256:"), Request: strings.TrimPrefix(o.ref, "sha256:"), Expiry: s.Window.Max}
	msg, err := approval.DecisionMessage(verb, refDigest(s.Packet), refDigest(o.ref), body.Expiry)
	if err != nil {
		return refuse("%v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(max(body.Expiry-env.now(), 1))*time.Second)
	defer cancel()
	armored, err := env.signDecision(ctx, o.signWith, msg, stderr)
	if err != nil {
		return refuse("--sign-with %s: %v; nothing was decided", o.signWith, err)
	}
	sig := string(armored)
	body.SSHSIG = &sig
	return postDecision(c, base, o.ref, body, "", stdout, refuse)
}

// signDecisionWithSSHKeygen is the production decision signer: ssh-keygen from PATH through
// host/sshsign's one driven, environment-scrubbed site, under DecisionNamespace. PATH is not a
// trust boundary: the daemon verifies the output against the enrolled key over World's message.
func signDecisionWithSSHKeygen(ctx context.Context, key string, msg []byte, stderr io.Writer) ([]byte, error) {
	return sshsign.SignIn(ctx, "ssh-keygen", approval.DecisionNamespace, key, msg, stderr)
}

// decideByArtifact is the decision mode's artifact flags (design §2.10): --write-message (CD6)
// writes the 91-byte decision message for the window's latest expiry to a NEW file and prints the
// signing line and the --signature step; --webauthn-link (CD8) prints, per origin of the active
// passkeys, a link to the A2 passkey page (unchanged) over the message's hash, and the
// --webauthn-assertion step; --signature (CD7) and --webauthn-assertion (CD9) post a file made for
// --expiry. The daemon judges the artifact (design §2.7); no terminal opens and no proof is written.
func decideByArtifact(ctx context.Context, o decideOptions, credPath, hostPort string, s decisionSubject, c *client, base string, stdout io.Writer, refuse func(string, ...any) int) int {
	verb, _ := approval.DecisionVerb(o.verb)
	body := surfacepage.DecideBody{V: 1, Verb: o.verb, Packet: strings.TrimPrefix(s.Packet, "sha256:"), Request: strings.TrimPrefix(o.ref, "sha256:"), Expiry: s.Window.Max}
	if o.expiry != 0 {
		body.Expiry = o.expiry
	}
	msg, err := approval.DecisionMessage(verb, refDigest(s.Packet), refDigest(o.ref), body.Expiry)
	if err != nil {
		return refuse("%v", err)
	}
	next := fmt.Sprintf("ailang-worldd approve %s%s %%s %%s --expiry %d --daemon %s --db %s --credentials-db %s", o.ref,
		map[string]string{"approve": "", "deny": " --deny"}[o.verb], body.Expiry, hostPort, shellQuote(o.db), shellQuote(credPath))
	switch {
	case o.msgPath != "": // CD6
		f, err := os.OpenFile(o.msgPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644) // never overwrites
		if err == nil {
			_, err = f.Write(msg)
			if cerr := f.Close(); err == nil {
				err = cerr
			}
		}
		if err != nil {
			return refuse("--write-message: %v; nothing was written", err)
		}
		fmt.Fprintf(stdout, "wrote the %d-byte message deciding %s request %s to %s\nsign it with an enrolled key before %s:\n  ssh-keygen -Y sign -n %s -f <your key> < %s > %s.sig\nthen decide it:\n  %s\n",
			len(msg), strings.ToUpper(o.verb), o.ref, o.msgPath, approverTime(body.Expiry), approval.DecisionNamespace, shellQuote(o.msgPath), shellQuote(o.msgPath),
			fmt.Sprintf(next, "--signature", shellQuote(o.msgPath+".sig")))
		return exitOK
	case o.passkeyLink: // CD8
		return decisionPasskeyLinks(ctx, o, credPath, body.Expiry, msg, fmt.Sprintf(next, "--webauthn-assertion", "<the saved file>"), stdout, refuse)
	case o.sigPath != "": // CD7
		armored, err := readRegularFileCapped(o.sigPath, approveMaxSignatureBytes)
		if err != nil {
			return refuse("--signature: %v; nothing was decided", err)
		}
		sig := string(armored)
		body.SSHSIG = &sig
	default: // CD9
		file, err := readRegularFileCapped(o.assertion, approval.MaxPasskeyResponseBytes)
		if err != nil {
			return refuse("--webauthn-assertion: %v; nothing was decided", err)
		}
		if _, err := approval.ParseAssertionResponse(file); err != nil {
			return refuse("--webauthn-assertion %s is not a passkey assertion: %v; nothing was decided", o.assertion, err)
		}
		asr := json.RawMessage(file)
		body.WebAuthn = &asr
	}
	return postDecision(c, base, o.ref, body, "", stdout, refuse)
}

// decisionPasskeyLinks is --webauthn-link on a request ref (CD8): one get link per distinct origin
// among the active passkeys (read-only), each allowing that origin's credentials, whose challenge is
// SHA-256 of the decision message. The A2 page signs the challenge its link carries and reads
// "Approve request <request>": the link's request names the verb (decide-<verb>-<16 hex of Q>), and
// the page shows the message hash printed here. It writes nothing and posts nothing.
func decisionPasskeyLinks(ctx context.Context, o decideOptions, credPath string, expiry int64, msg []byte, next string, stdout io.Writer, refuse func(string, ...any) int) int {
	ro, err := store.OpenCredentialsReadOnly(credPath)
	if err != nil {
		return refuse("cannot read the credentials store %s: %v; nothing was decided", credPath, err)
	}
	active, err := ro.ActivePasskeys(ctx)
	_ = ro.Close()
	if err != nil {
		return refuse("%v; nothing was decided", err)
	}
	if len(active) == 0 {
		return refuse("no passkey approver is enrolled; see ailang-worldd approver list")
	}
	h := sha256.Sum256(msg)
	hash16, request := hex.EncodeToString(h[:8]), "decide-"+o.verb+"-"+strings.TrimPrefix(o.ref, "sha256:")[:16]
	fmt.Fprintf(stdout, "%s request %s with a passkey (expires %s, message hash %s): open the link for its origin; the page reads \"Approve request %s\" for either verb: check it shows message hash %s, and save the file it offers:\n",
		strings.ToUpper(o.verb), o.ref, approverTime(expiry), hash16, request, hash16)
	links := map[string]*ceremonyLink{}
	var order []string
	for _, e := range active { // ActivePasskeys is ordered by name
		l := links[e.Passkey.Origin]
		if l == nil {
			l = &ceremonyLink{V: 1, Op: "get", RP: e.Passkey.RPID, Origin: e.Passkey.Origin, Request: request, Expires: expiry, Hash16: hash16,
				Challenge: base64.RawURLEncoding.EncodeToString(h[:])}
			links[e.Passkey.Origin] = l
			order = append(order, e.Passkey.Origin)
		}
		l.Allow = append(l.Allow, base64.RawURLEncoding.EncodeToString(e.Passkey.CredentialID))
	}
	for _, origin := range order {
		fmt.Fprintf(stdout, "  %s\n", links[origin].link())
	}
	fmt.Fprintf(stdout, "then decide it:\n  %s\n", next)
	return exitOK
}

// postDecision POSTs body to /surface/decide/{ref} with the daemon's own loopback Origin (in O by
// construction; design §2.8) and, for the typed arm, the operator proof; it prints the decision.
func postDecision(c *client, origin, ref string, body surfacepage.DecideBody, proof string, stdout io.Writer, refuse func(string, ...any) int) int {
	data, err := json.Marshal(body)
	if err != nil {
		return refuse("%v", err)
	}
	hdr := map[string]string{"Origin": origin}
	if proof != "" {
		hdr["X-World-Operator-Proof"] = proof
	}
	ctx, cancel := context.WithTimeout(context.Background(), approveDecideBudget)
	defer cancel()
	status, _, resp, err := c.doRequest(ctx, http.MethodPost, "/surface/decide/"+ref, bytes.NewReader(data), hdr)
	if err != nil {
		return refuse("%v; run approve %s again to see whether it was decided", err, ref)
	}
	if status != http.StatusOK {
		return refuse("the daemon refused to decide request %s (HTTP %d): %s", ref, status, refusalReason(resp))
	}
	var made surfacepage.DecisionMade
	if err := json.Unmarshal(resp, &made); err != nil || made.Request != ref {
		return refuse("the daemon's answer for request %s is not a decision (%s)", ref, strings.TrimSpace(string(resp)))
	}
	fmt.Fprintf(stdout, "decided %s: request %s\ndecision %s by %s\n", made.Decided, made.Request, made.Decision, made.DecidedBy)
	return exitOK
}
