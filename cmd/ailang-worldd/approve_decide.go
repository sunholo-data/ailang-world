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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// decisionRefPattern is a Human.Approve request ref: approve's decision mode (design §2.10).
var decisionRefPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

// approveDecideBudget bounds the decision mode's reads before the y/N, and its POST after it.
const approveDecideBudget = 10 * time.Second

// decideOptions are approve's flags in the decision mode.
type decideOptions struct {
	ref, verb, db, credsDB, daemon string
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
