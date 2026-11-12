package daemon

// Row 201 A4 sprint 1 (design §2.4, §2.5, §2.8, §6 RD2, RD4-RD8, R4d): POST /surface/decide/{ref}
// with the typed arm: the guards, the law's arms, the proof, replay after a restart, R4 by L1.

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// requireNothingDecided fails unless no decision was stored, the head did not move and the proof
// of q is still unconsumed.
func (r decideRig) requireNothingDecided(what string, q hashref.HashRef, head hashref.HashRef) {
	r.t.Helper()
	if n, h := r.decisions(); n != 0 || h != head {
		r.t.Fatalf("%s: %d decisions, head %v (was %v)", what, n, h, head)
	}
	if !r.proofExists(q) {
		r.t.Fatalf("%s: the proof was consumed", what)
	}
}

// RD2: R0 Authorization, R0 Cookie, R2 a foreign Host, R3 no/null/other Origin, R4 same-site, R5
// text/plain, R6 an unknown field, decidedBy, a key format, typed false, v 2, a bad digest, two
// values, too large: each refused before the law, 0 decisions, the proof unconsumed; then the
// same proof decides (the control).
func TestDecidePOSTGuards(t *testing.T) {
	r := newDecideRig(t)
	q, p := r.packet("guards", broker.WorkspaceScope, "agent")
	proof := r.proof(q, "decide-approve")
	good := decideBody("approve", p, q, decNow+50)
	_, head := r.decisions()
	field := func(extra string) []byte { return []byte(strings.Replace(string(good), `"typed":true`, extra, 1)) }
	for _, tc := range []struct {
		name   string
		host   string
		body   []byte
		hdr    []string
		status int
		class  string
	}{
		{"R0 Authorization", apLoopback, good, []string{"Authorization", "Bearer agent-session"}, 403, "Forbidden"},
		{"R0 empty Authorization", apLoopback, good, []string{"Authorization", ""}, 403, "Forbidden"},
		{"R0 Cookie", apLoopback, good, []string{"Cookie", "s=1"}, 403, "Forbidden"},
		{"R2 foreign Host", "rebind.example:7644", good, nil, 421, "MisdirectedRequest"},
		{"R3 no Origin", apLoopback, good, []string{"Origin", ""}, 403, "CrossOrigin"},
		{"R3 null Origin", apLoopback, good, []string{"Origin", "null"}, 403, "CrossOrigin"},
		{"R3 other Origin", apLoopback, good, []string{"Origin", "http://localhost:7645"}, 403, "CrossOrigin"},
		{"R4 same-site", apLoopback, good, []string{"Sec-Fetch-Site", "same-site"}, 403, "CrossOrigin"},
		{"R5 text/plain", apLoopback, good, []string{"Content-Type", "text/plain"}, 415, "UnsupportedMediaType"},
		{"R6 unknown field", apLoopback, field(`"typed":true,"x":1`), nil, 400, "BadRequest"},
		{"R6 decidedBy", apLoopback, field(`"typed":true,"decidedBy":"ssh:mark"`), nil, 400, "BadRequest"},
		{"R6 a key format", apLoopback, field(`"typed":true,"sshsig":"-----BEGIN SSH SIGNATURE-----"`), nil, 400, "BadRequest"},
		{"R6 typed false", apLoopback, field(`"typed":false`), nil, 400, "BadRequest"},
		{"R6 no format", apLoopback, []byte(strings.Replace(string(good), `,"typed":true`, "", 1)), nil, 400, "BadRequest"},
		{"R6 v 2", apLoopback, []byte(strings.Replace(string(good), `"v":1`, `"v":2`, 1)), nil, 400, "BadRequest"},
		{"R6 upper-case packet", apLoopback, []byte(strings.Replace(string(good), p.Hash.Digest(), strings.ToUpper(p.Hash.Digest()), 1)), nil, 400, "BadRequest"},
		{"R6 short request", apLoopback, []byte(strings.Replace(string(good), q.Digest(), q.Digest()[:62], 1)), nil, 400, "BadRequest"},
		{"R6 two values", apLoopback, append(append([]byte{}, good...), good...), nil, 400, "BadRequest"},
		{"R6 too large", apLoopback, []byte(`{"v":1,"verb":"` + strings.Repeat("a", surfacepage.MaxApproveBodyBytes) + `"}`), nil, 413, "TooLarge"},
	} {
		hdr := append([]string{"Origin", "http://127.0.0.1:7644", "Content-Type", "application/json", "X-World-Operator-Proof", proof}, tc.hdr...)
		rec := r.serve("POST", "/surface/decide/"+q.String(), tc.host, tc.body, hdr...)
		requireRefusal(t, tc.name, rec, tc.status, tc.class)
		r.requireNothingDecided(tc.name, q, head)
	}
	if rec := r.post(q.String(), good, proof); rec.Code != 200 {
		t.Fatalf("control: %d %s", rec.Code, rec.Body.String())
	}
}

// RD4: a store built per D1-D10 → the law's status and class, 0 decisions, the proof unconsumed
// (the law runs before the arm).
func TestDecidePOSTLawArms(t *testing.T) {
	r := newDecideRig(t)
	q, p := r.packet("law", broker.WorkspaceScope, "agent")
	proof := r.proof(q, "decide-approve")
	_, head := r.decisions()
	ctx := boundedTestContext(t)
	offChain := forgeObject(broker.ApprovalRequestV1, []byte(`{"effect":"Human.Approve","scope":"worktree","cost":1,"requester":"off","now":1700}`))
	notApprove := forgeObject(broker.ApprovalRequestV1, []byte(`{"effect":"Human.PollApproval","scope":"worktree","cost":1,"requester":"x","now":1700}`))
	for _, o := range []store.Object{offChain, notApprove} {
		if err := r.d.store.PutObject(ctx, o); err != nil {
			t.Fatal(err)
		}
	}
	pending := approvalRequestInChain(t, r.d, broker.WorkspaceScope, "no-packet", packetNow)
	failed := approvalRequestInChain(t, r.d, broker.WorkspaceScope, "failed", packetNow)
	fc := newForgedChain("law-failed", failed, broker.WorkspaceScope)
	fc.failed = true
	fp, _ := commitForgedChain(t, r.d, fc)
	_, head = r.decisions()
	for _, tc := range []struct {
		name, ref string
		body      []byte
		now       int64
		status    int
		class     string
		reason    string
	}{
		{"D1", "sha256:XYZ", decideBody("approve", p, q, decNow+50), decNow, 400, "BadRequest", "is not a request ref"},
		{"D2 absent", "sha256:" + strings.Repeat("0", 64), decideBody("approve", p, q, decNow+50), decNow, 404, "NotFound", "no such approval request"},
		{"D2 a packet", p.Hash.String(), decideBody("approve", p, q, decNow+50), decNow, 404, "NotFound", "no such approval request"},
		{"D2 another effect", notApprove.Hash.String(), decideBody("approve", p, q, decNow+50), decNow, 404, "NotFound", "no such approval request"},
		{"D3", offChain.Hash.String(), decideBody("approve", p, q, decNow+50), decNow, 404, "NotFound", "not on the approvals chain"},
		{"D4b zero", pending.String(), decideBody("approve", p, pending, decNow+50), decNow, 409, "Conflict", "no committed decision packet carries request"},
		{"D4b refused packet", failed.String(), decideBody("approve", fp, failed, decNow+50), decNow, 409, "Conflict", "is refused: not a validated decision packet: requestRef"},
		{"D7", q.String(), decideBody("approve", p, q, packetNow+61), packetNow + 60, 409, "Conflict", "deadline passed"},
		{"D8 packet", q.String(), decideBody("approve", fp, q, decNow+50), decNow, 409, "Conflict", "reload it"},
		{"D8 request", q.String(), []byte(strings.Replace(string(decideBody("approve", p, q, decNow+50)), q.Digest(), pending.Digest(), 1)), decNow, 409, "Conflict", "reload it"},
		{"D9", q.String(), decideBody("maybe", p, q, decNow+50), decNow, 400, "BadRequest", "a decision is approve or deny"},
		{"D10 stale", q.String(), decideBody("approve", p, q, decNow), decNow, 409, "Conflict", "at most 10 minutes"},
		{"D10 past the deadline", q.String(), decideBody("approve", p, q, packetNow+61), decNow, 409, "Conflict", "at most 10 minutes"},
	} {
		r.d.approveNow = func() int64 { return tc.now }
		got := requireRefusal(t, tc.name, r.post(tc.ref, tc.body, proof), tc.status, tc.class)
		if !strings.Contains(got.Reason, tc.reason) {
			t.Errorf("%s: %q lacks %q", tc.name, got.Reason, tc.reason)
		}
		r.requireNothingDecided(tc.name, q, head)
	}
	// D6 is RD8's; D5 is R4d's; D4a and the counts are RD9*'s.
}

// RD5: a typed proof decides: one decision by operator:terminal, the proof unlinked; deny too.
func TestDecideTypedLands(t *testing.T) {
	r := newDecideRig(t)
	for _, verb := range []string{"approve", "deny"} {
		q, p := r.packet("typed-"+verb, broker.WorkspaceScope, "agent-"+verb)
		before, _ := r.decisions()
		rec := r.post(q.String(), decideBody(verb, p, q, decNow+50), r.proof(q, "decide-"+verb))
		var got surfacepage.DecisionMade
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil || rec.Code != 200 || got.V != 1 || got.Decided != verb || got.Request != q.String() || got.DecidedBy != "operator:terminal" {
			t.Fatalf("%s: %d %s", verb, rec.Code, rec.Body.String())
		}
		if err := approvalHeaders(rec); err != nil {
			t.Fatal(err)
		}
		on, decided, v, by, err := broker.ApprovalChainState(boundedTestContext(t), r.d.store, q)
		if err != nil || !on || !decided || v != verb || by != "operator:terminal" {
			t.Fatalf("%s: chain %v %v %q %q %v", verb, on, decided, v, by, err)
		}
		obj, ok, err := r.d.store.GetObject(boundedTestContext(t), hashref.MustParse(got.Decision))
		if err != nil || !ok || obj.SemanticID != broker.ApprovalDecisionV1 || !strings.Contains(string(obj.Payload), `"now":1710`) {
			t.Fatalf("%s: decision object %s %v %v", verb, obj.Payload, ok, err)
		}
		if after, _ := r.decisions(); after != before+1 || r.proofExists(q) {
			t.Fatalf("%s: decisions %d -> %d, proof left %v", verb, before, after, r.proofExists(q))
		}
	}
}

// RD6: the proof binds the verb and the request: a decide-deny proof cannot approve; another
// request's proof, a wrong secret and no proof are refused; 0 decisions.
func TestDecideTypedProofBindsTheVerb(t *testing.T) {
	r := newDecideRig(t)
	q, p := r.packet("verb", broker.WorkspaceScope, "agent")
	other, _ := r.packet("verb-other", broker.WorkspaceScope, "other")
	_, head := r.decisions()
	body := decideBody("approve", p, q, decNow+50)
	for _, tc := range []struct{ name, phase, proof string }{
		{"deny proof", "decide-deny", ""}, {"finish phase", "approve", ""}, {"wrong secret", "decide-approve", strings.Repeat("e", 64)}, {"no header", "decide-approve", "-"},
	} {
		secret := r.proof(q, tc.phase)
		switch tc.proof {
		case "":
		case "-":
			secret = ""
		default:
			secret = tc.proof
		}
		got := requireRefusal(t, tc.name, r.post(q.String(), body, secret), 403, "Forbidden")
		if !strings.Contains(got.Reason, "the typed decision is not proven: operator proof") {
			t.Errorf("%s: %q", tc.name, got.Reason)
		}
		if n, h := r.decisions(); n != 0 || h != head {
			t.Fatalf("%s: %d decisions", tc.name, n)
		}
	}
	// Another request's proof is not this request's.
	_ = os.Remove(r.proofPath(q))
	got := requireRefusal(t, "other request", r.post(q.String(), body, r.proof(other, "decide-approve")), 403, "Forbidden")
	if !strings.Contains(got.Reason, "operator proof is missing") || !r.proofExists(other) {
		t.Errorf("other request: %q", got.Reason)
	}
}

// RD7: P or Q of another packet → 409 D8.
func TestDecidePOSTStaleDigests(t *testing.T) {
	r := newDecideRig(t)
	q, p := r.packet("stale", broker.WorkspaceScope, "agent")
	q2, p2 := r.packet("stale-2", broker.WorkspaceScope, "agent-2")
	proof := r.proof(q, "decide-approve")
	_, head := r.decisions()
	for name, body := range map[string][]byte{"P": decideBody("approve", p2, q, decNow+50), "Q": decideBody("approve", p, q2, decNow+50)} {
		got := requireRefusal(t, name, r.post(q.String(), body, proof), 409, "Conflict")
		if !strings.Contains(got.Reason, "the request changed since it was shown; reload it") {
			t.Errorf("%s: %q", name, got.Reason)
		}
		r.requireNothingDecided(name, q, head)
	}
}

// RD8: decide, close the daemon, open a new one on the same stores, POST again with a fresh proof
// → 409 D6, still one decision.
func TestDecideReplayAfterRestart(t *testing.T) {
	r := newDecideRig(t)
	q, p := r.packet("replay", broker.WorkspaceScope, "agent")
	body := decideBody("deny", p, q, decNow+50)
	if rec := r.post(q.String(), body, r.proof(q, "decide-deny")); rec.Code != 200 {
		t.Fatalf("first: %d %s", rec.Code, rec.Body.String())
	}
	n, head := r.decisions()
	db := r.d.cfg.DBPath
	if err := r.d.Close(); err != nil {
		t.Fatal(err)
	}
	d2, err := New(boundedTestContext(t), Config{DBPath: db, BindHost: DefaultBindHost, ErrorLog: io.Discard})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = d2.Close() })
	d2.approveNow = func() int64 { return decNow + 1 }
	r2 := decideRig{approvalRig{t: t, d: d2}}
	if filepath.Dir(d2.credentialsPath) != filepath.Dir(r.d.credentialsPath) {
		t.Fatal("the restarted daemon has another credentials directory")
	}
	for _, verb := range []string{"deny", "approve"} {
		got := requireRefusal(t, "replay "+verb, r2.post(q.String(), decideBody(verb, p, q, decNow+50), r2.proof(q, "decide-"+verb)), 409, "Conflict")
		if !strings.Contains(got.Reason, "was already decided deny by operator:terminal") {
			t.Errorf("replay %s: %q", verb, got.Reason)
		}
		r2.requireNothingDecidedSince(n, head)
	}
}

func (r decideRig) requireNothingDecidedSince(n int, head hashref.HashRef) {
	r.t.Helper()
	if m, h := r.decisions(); m != n || h != head {
		r.t.Fatalf("%d decisions (was %d), head %v (was %v)", m, n, h, head)
	}
}

// R4d: publish- and finish-shaped requests are refused for the typed arm too (L1, D5), with the
// owning command; 0 decisions, the proof unconsumed.
func TestDecideRefusesReservedScopesTyped(t *testing.T) {
	r := newDecideRig(t)
	publish := broker.PublishApprovalScopeFor(broker.PublishIdentity{RegistryOrigin: "https://registry.example", Vendor: "v", Name: "n", Version: "1.0.0",
		ManifestRef: hashref.SumSHA256([]byte("m"))}, broker.PublishHashes{TarballSHA256: "t", ContentHash: "c", InterfaceHash: "i"}, packetNow+60)
	finish := broker.FormatFinishApprovalScope(broker.FinishApprovalScope{Repository: "o/r", Remote: "https://github.com/o/r.git", Head: "world/e", SHA: strings.Repeat("a", 40),
		Base: "dev", TitleSHA256: strings.Repeat("b", 64), BodySHA256: strings.Repeat("c", 64), ExpiresAt: packetNow + 60})
	for _, tc := range []struct{ name, scope, owner string }{
		{"publish", publish, "`world-publish approve` at a terminal (R4)"}, {"publish near miss", "registry:x", "`world-publish approve`"},
		{"finish", finish, "`ailang-worldd episode finish`"}, {"finish near miss", "github-pr:x", "`ailang-worldd episode finish`"},
	} {
		q, p := r.packet("reserved-"+tc.name, tc.scope, "agent")
		_, head := r.decisions()
		got := requireRefusal(t, tc.name, r.post(q.String(), decideBody("approve", p, q, decNow+50), r.proof(q, "decide-approve")), 403, "Forbidden")
		if !strings.Contains(got.Reason, tc.owner) {
			t.Errorf("%s: %q lacks %q", tc.name, got.Reason, tc.owner)
		}
		r.requireNothingDecided(tc.name, q, head)
	}
}
