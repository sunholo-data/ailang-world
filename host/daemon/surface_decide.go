package daemon

// Row 201 A4 sprint 1 (design_docs/planned/w-world-approvals-by-key-a4.md §2.4, §2.5, §2.8): an
// operator decides a Human.Approve request whose decision packet the coordinator committed.
// GET /surface/decisions/{ref} serves the subject (the request's and the packet's stored bytes, the
// expiry window, the law's display verdict); POST /surface/decide/{ref} runs A3a's guards, then,
// under humanMu, the decision law (D1-D10) over facts read from the world store, then the format
// arm (D11), then broker.DecideApproval with a decidedBy the daemon derives. Sprint 1's only arm is
// typed: the row-236 operator proof the CLI writes after a typed y. Neither route is under /v1.

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// decisionRefPattern is D1's grammar: a request ref is sha256: and 64 lowercase hex.
var decisionRefPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

// typedDecider is the decidedBy of a decision made by the typed arm.
const typedDecider = "operator:terminal"

// decisionRead is what readDecision found: the law's facts and, when exactly one packet carries
// the request, the subject's objects.
type decisionRead struct {
	facts   broker.DecisionFacts
	request store.Object
	packet  store.Object
	entry   int64
}

// readDecision reads every fact the law needs for request ref at now from the world store (design
// §2.5): the request object, its place on the approvals chain, the complete packet discovery
// (§2.4 [r1]) and, for exactly one carrying packet, its ValidatePacket verdict. A ref that is not
// D1's grammar reads nothing. Store errors are returned, never a fact.
func (d *Daemon) readDecision(ctx context.Context, ref string, now int64) (decisionRead, error) {
	rd := decisionRead{facts: broker.DecisionFacts{Ref: ref, Now: now}}
	f := &rd.facts
	if !decisionRefPattern.MatchString(ref) { // D1
		return rd, nil
	}
	f.RefOK = true
	q := hashref.MustParse(ref)
	digest, _ := hex.DecodeString(q.Digest())
	f.Request = [32]byte(digest)
	obj, ok, err := d.store.GetObject(ctx, q)
	if err != nil {
		return rd, err
	}
	if req, err := broker.DecodeApprovalRequest(obj.Payload); ok && obj.SemanticID == broker.ApprovalRequestV1 && err == nil && req.Effect == broker.EffectHumanApprove { // D2
		f.RequestOK, f.Scope, rd.request = true, req.Scope, obj
	}
	if f.OnChain, f.Decided, f.DecidedVerb, f.DecidedBy, err = broker.ApprovalChainState(ctx, d.store, q); err != nil { // D3, D6
		return rd, err
	}
	var carrying store.Object
	if f.Discovery, carrying, err = coordinator.DiscoverRequestPackets(ctx, d.store, q); err != nil { // D4a
		return rd, err
	}
	if f.Discovery.Complete && f.Discovery.Undecodable == 0 && f.Discovery.Carrying == 1 { // D4b, D7, D8
		res, err := coordinator.ValidatePacket(ctx, d.store, carrying)
		if err != nil {
			return rd, err
		}
		pd, _ := hex.DecodeString(carrying.Hash.Digest())
		f.PacketCode, f.PacketSentence, f.Deadline, f.Packet = res.Code, packetRefusalSentence(res), res.Packet.DeadlineAt, [32]byte(pd)
		rd.packet, rd.entry = carrying, res.Entry
	}
	return rd, nil
}

// lawClass is the refusal class of a law status.
func lawClass(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "BadRequest"
	case http.StatusForbidden:
		return "Forbidden"
	case http.StatusNotFound:
		return "NotFound"
	}
	return "Conflict"
}

// handleSurfaceDecision is GET /surface/decisions/{ref}: the subject of a decision, as bytes. A
// request that cannot have a subject (D1-D4: no such request, not on the chain, not exactly one
// validated packet) is the law's refusal; otherwise 200 with the law's D5-D7 verdict at now as
// display (refusal, reserved, the decision when one landed). The POST runs the law again.
func (d *Daemon) handleSurfaceDecision(w http.ResponseWriter, r *http.Request, _ []string) {
	ref := r.PathValue("ref")
	now := d.approveNow()
	ctx, cancel := d.readCtx(r)
	defer cancel()
	rd, err := d.readDecision(ctx, ref, now)
	if err != nil {
		d.approvalInternal(w, r, err)
		return
	}
	law := broker.SubjectRefusal(rd.facts)
	switch law.Step {
	case "D1", "D2", "D3", "D4a", "D4b":
		writeApprovalRefusal(w, law.Status, lawClass(law.Status), law.Reason)
		return
	}
	p, _ := coordinator.DecodeDecisionPacket(rd.packet.Payload) // validated: code 0
	reserved, _ := broker.ReservedApprovalScope(rd.facts.Scope)
	passkeys, err := d.decisionPasskeys(ctx)
	if err != nil {
		d.approvalInternal(w, r, err)
		return
	}
	writeApprovalJSON(w, http.StatusOK, surfacepage.DecisionSubject{V: 1,
		Request: ref, RequestBytes: base64.RawURLEncoding.EncodeToString(rd.request.Payload),
		Packet: rd.packet.Hash.String(), PacketBytes: base64.RawURLEncoding.EncodeToString(rd.packet.Payload),
		Entry: rd.entry, Proposal: p.ProposalHash.String(), Now: now,
		Window:  surfacepage.DecisionWindow{Min: now + 1, Max: broker.DecisionWindowMax(now, p.DeadlineAt)},
		Refusal: law.Reason, Reserved: reserved, DecidedBy: rd.facts.DecidedBy, Decision: rd.facts.DecidedVerb,
		Passkeys: passkeys, SSH: surfacepage.DecisionSSH{Namespace: approval.DecisionNamespace}})
}

// decisionPasskeys is the active passkeys by origin on a fresh read-only credentials handle, for
// the card's ceremony (sprint 4); a store below schema 6 has none.
func (d *Daemon) decisionPasskeys(ctx context.Context) ([]surfacepage.CardPasskeys, error) {
	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = ro.Close() }()
	active, err := ro.ActivePasskeys(ctx)
	if err != nil && !isOlderStore(err) {
		return nil, err
	}
	var passkeys []surfacepage.CardPasskey
	for _, e := range active {
		passkeys = append(passkeys, surfacepage.CardPasskey{Origin: e.Passkey.Origin, RPID: e.Passkey.RPID, CredentialID: e.Passkey.CredentialID})
	}
	return surfacepage.DecisionPasskeys(passkeys), nil
}

// handleSurfaceDecide is POST /surface/decide/{ref} (design §2.8). After R2 (approvalRoute): R0 any
// Authorization or Cookie is refused (the agent holds a session, so a session-bearing request is
// one it could make), then A3a's R3-R6 verbatim. Then, under humanMu, in this order (sprint 2's
// key arms take the arm's slot and nothing else moves; plan §1.3 S1-D9):
//  1. the law, D1-D10, over facts read now (a refusal writes nothing and consumes no proof);
//  2. the arm (D11): typed consumes the operator proof for this request and verb, unlinking it;
//  3. broker.DecideApproval, which lands the first decision only, with decidedBy derived here.
func (d *Daemon) handleSurfaceDecide(w http.ResponseWriter, r *http.Request, o []string) {
	if len(r.Header.Values("Authorization")) != 0 || len(r.Header.Values("Cookie")) != 0 { // R0
		writeApprovalRefusal(w, http.StatusForbidden, "Forbidden", "a decision is never made with a session or a cookie (the agent holds a session); decide at a terminal with ailang-worldd approve <request-ref>")
		return
	}
	if origin := r.Header.Get("Origin"); !allowedApprovalOrigin(origin, o) { // R3
		writeApprovalRefusal(w, http.StatusForbidden, "CrossOrigin", crossOriginReason(origin))
		return
	}
	if !sameOriginFetch(r.Header.Get("Sec-Fetch-Site")) { // R4
		writeApprovalRefusal(w, http.StatusForbidden, "CrossOrigin", "Sec-Fetch-Site is not same-origin")
		return
	}
	if !jsonMediaType(r.Header.Get("Content-Type")) { // R5
		writeApprovalRefusal(w, http.StatusUnsupportedMediaType, "UnsupportedMediaType", "the body must be application/json")
		return
	}
	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, surfacepage.MaxApproveBodyBytes)) // R6
	if err != nil {
		writeApprovalRefusal(w, http.StatusRequestEntityTooLarge, "TooLarge", fmt.Sprintf("the body is over %d bytes", surfacepage.MaxApproveBodyBytes))
		return
	}
	body, packet, request, err := surfacepage.ParseDecideBody(raw)
	if err != nil {
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", err.Error())
		return
	}
	ref := r.PathValue("ref")
	d.humanMu.Lock()
	defer d.humanMu.Unlock()
	ctx, cancel := context.WithTimeout(r.Context(), d.commitBudget)
	defer cancel()
	now := d.approveNow()
	rd, err := d.readDecision(ctx, ref, now)
	if err != nil {
		d.approvalInternal(w, r, err)
		return
	}
	f := rd.facts
	f.PostedPacket, f.PostedRequest, f.Verb, f.Expiry = packet, request, body.Verb, body.Expiry
	if law := broker.DecisionRefusal(f); law.Step != "" { // 1. the law
		writeApprovalRefusal(w, law.Status, lawClass(law.Status), law.Reason)
		return
	}
	q := hashref.MustParse(ref) // D1 held
	// 2. the arm (D11): a key's judge and evidence (an SSHSIG, sprint 2; a passkey, sprint 3;
	// surface_decide_key.go), or the typed proof, whose phase carries the verb.
	decidedBy, evidence := typedDecider, (*store.Credentials)(nil)
	if body.SSHSIG != nil || body.WebAuthn != nil {
		var ok bool
		if evidence, decidedBy, ok = d.keyArm(ctx, w, r, f, body, now); !ok {
			return
		}
		defer func() { _ = evidence.Close() }()
	} else if err := d.consumeOperatorProof("decision-"+q.Digest(), "decide-"+body.Verb, r.Header.Get("X-World-Operator-Proof")); err != nil {
		writeApprovalRefusal(w, http.StatusForbidden, "Forbidden", "the typed decision is not proven: "+err.Error()+"; decide at a terminal with ailang-worldd approve "+ref)
		return
	}
	decision, err := broker.DecideApproval(ctx, d.store, q, body.Verb, decidedBy, now) // 3.
	if err != nil {
		d.approvalInternal(w, r, err)
		return
	}
	if evidence != nil { // 4. the key decision's evidence names the decision that landed
		if err := evidence.LandKeyDecision(ctx, ref, decision.String()); err != nil {
			d.approvalInternal(w, r, err)
			return
		}
	}
	writeApprovalJSON(w, http.StatusOK, surfacepage.DecisionMade{V: 1, Decided: body.Verb, Request: ref, Decision: decision.String(), DecidedBy: decidedBy})
}
