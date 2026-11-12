package daemon

// Row 201 A4 sprint 2 M5a (design §2.5 D11, §2.7, §6 RD10, RD11 and the judge's refusals; the
// passkey halves are sprint 3's): POST /surface/decide/{ref} with an SSHSIG by an enrolled
// approvaltest key (D-WORLD-109), through d.Handler().

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

var kdSigner = approvaltest.NewSigner(71)

// enrolDecider enrols s as name at packetNow, before the rig's clock (decNow).
func (r decideRig) enrolDecider(name string, s *approvaltest.Signer) {
	r.t.Helper()
	r.write(func(ctx context.Context, c *store.Credentials) error {
		return c.EnrolApprover(ctx, store.ApproverRow{Name: name, Format: approval.FormatSSH, PublicKey: s.PublicKey().Blob, EnrolledAt: packetNow})
	})
}

// decideSSHBody is a POST body deciding verb on (p, q, expiry) with s's SSHSIG under namespace.
func decideSSHBody(t *testing.T, s *approvaltest.Signer, namespace, verb string, p store.Object, q hashref.HashRef, expiry int64) []byte {
	t.Helper()
	v, err := approval.DecisionVerb(verb)
	if err != nil {
		v = approval.VerbApprove
	}
	msg, err := approval.DecisionMessage(v, [32]byte(mustDigest(t, p.Hash)), [32]byte(mustDigest(t, q)), expiry)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(map[string]any{"v": 1, "verb": verb, "packet": p.Hash.Digest(), "request": q.Digest(), "expiry": expiry, "sshsig": string(s.Sign(namespace, msg))})
	return b
}

func mustDigest(t *testing.T, h hashref.HashRef) []byte {
	t.Helper()
	d, err := hex.DecodeString(h.Digest())
	if err != nil || len(d) != 32 {
		t.Fatal(h, err)
	}
	return d
}

// keyRow is q's key_decisions row, read on a fresh read-only handle.
func (r decideRig) keyRow(q hashref.HashRef) (store.KeyDecisionRow, bool) {
	r.t.Helper()
	ro, err := store.OpenCredentialsReadOnly(r.d.credentialsPath)
	if err != nil {
		r.t.Fatal(err)
	}
	defer func() { _ = ro.Close() }()
	row, ok, err := ro.KeyDecision(boundedTestContext(r.t), q.String())
	if err != nil {
		r.t.Fatal(err)
	}
	return row, ok
}

// postKey is a key POST: the loopback headers and no proof.
func (r decideRig) postKey(q hashref.HashRef, body []byte) *httptest.ResponseRecorder {
	return r.post(q.String(), body, "")
}

// RD10, RD11: an enrolled key's SSHSIG decides approve and deny as ssh:<name>, derived from the
// verified key (a body naming decidedBy is 400); the decision is on the chain and the evidence row
// is landed with the decision's ref.
func TestDecidePOSTSSHSIG(t *testing.T) {
	r := newDecideRig(t)
	r.enrolDecider("mark-ssh", kdSigner)
	for _, verb := range []string{"approve", "deny"} {
		q, p := r.packet("ssh-"+verb, broker.WorkspaceScope, "agent-"+verb)
		forged := strings.Replace(string(decideSSHBody(t, kdSigner, approval.DecisionNamespace, verb, p, q, decNow+50)), `"v":1`, `"v":1,"decidedBy":"ssh:root"`, 1)
		requireRefusal(t, "RD11 decidedBy", r.postKey(q, []byte(forged)), 400, "BadRequest")
		rec := r.postKey(q, decideSSHBody(t, kdSigner, approval.DecisionNamespace, verb, p, q, decNow+50))
		var made surfacepage.DecisionMade
		if err := json.Unmarshal(rec.Body.Bytes(), &made); err != nil || rec.Code != 200 || made.Decided != verb || made.DecidedBy != "ssh:mark-ssh" {
			t.Fatalf("%s: %d %s", verb, rec.Code, rec.Body.String())
		}
		_, decided, gotVerb, by, err := broker.ApprovalChainState(boundedTestContext(t), r.d.store, q)
		row, ok := r.keyRow(q)
		if err != nil || !decided || gotVerb != verb || by != "ssh:mark-ssh" || !ok || row.DecisionRef != made.Decision || row.Verb != verb || !strings.HasPrefix(string(row.Artifact), "-----BEGIN SSH SIGNATURE-----") ||
			row.PacketRef != p.Hash.String() || row.Approver != "ssh:mark-ssh" || row.Expiry != decNow+50 || row.DecidedAt != decNow {
			t.Fatalf("%s: chain %v %q %q %v; row %+v %v; made %+v", verb, decided, gotVerb, by, err, row, ok, made)
		}
	}
}

// D11 (sshsig): the judge's refusals are 422 with nothing decided and nothing recorded: an
// unenrolled key, a revoked approver, a mint-namespace signature, an approve signature posted as
// deny, a signature over another packet, another expiry, and garbage; and the law runs first: a
// valid signature with an expiry past the window is D10's 409.
func TestDecidePOSTSSHSIGRefusals(t *testing.T) {
	r := newDecideRig(t)
	r.enrolDecider("mark-ssh", kdSigner)
	r.enrolDecider("gone", approvaltest.NewSigner(72))
	r.write(func(ctx context.Context, c *store.Credentials) error {
		return c.RevokeApprover(ctx, "gone", packetNow+1)
	})
	q, p := r.packet("ssh-refusals", broker.WorkspaceScope, "agent")
	_, other := r.packet("ssh-other", broker.WorkspaceScope, "other")
	_, head := r.decisions()
	good := decideSSHBody(t, kdSigner, approval.DecisionNamespace, "approve", p, q, decNow+50)
	for _, tc := range []struct {
		name   string
		body   []byte
		status int
		reason string
	}{
		{"unenrolled", decideSSHBody(t, approvaltest.NewSigner(73), approval.DecisionNamespace, "approve", p, q, decNow+50), 422, "no enrolled approver holds that key"},
		{"revoked", decideSSHBody(t, approvaltest.NewSigner(72), approval.DecisionNamespace, "approve", p, q, decNow+50), 422, "the approver is revoked"},
		{"mint namespace", decideSSHBody(t, kdSigner, approval.Namespace, "approve", p, q, decNow+50), 422, "wrong namespace"},
		{"verb", []byte(strings.Replace(string(good), `"verb":"approve"`, `"verb":"deny"`, 1)), 422, "signature does not verify"},
		{"another packet", decideSSHBody(t, kdSigner, approval.DecisionNamespace, "approve", other, q, decNow+50), 409, "reload it"},
		{"another expiry", []byte(strings.Replace(string(good), fmt.Sprintf(`"expiry":%d`, decNow+50), fmt.Sprintf(`"expiry":%d`, decNow+49), 1)), 422, "signature does not verify"},
		{"garbage", []byte(strings.Replace(string(good), `"sshsig":"-----BEGIN`, `"sshsig":"x-----BEGIN`, 1)), 422, "bad armor"},
		{"D10 first", decideSSHBody(t, kdSigner, approval.DecisionNamespace, "approve", p, q, packetNow+61), 409, "stale or too far ahead"},
	} {
		got := requireRefusal(t, tc.name, r.postKey(q, tc.body), tc.status, map[int]string{422: "NotApproved", 409: "Conflict"}[tc.status])
		if !strings.Contains(got.Reason, tc.reason) {
			t.Errorf("%s: %q", tc.name, got.Reason)
		}
		if n, h := r.decisions(); n != 0 || h != head {
			t.Fatalf("%s: %d decisions", tc.name, n)
		}
		if _, ok := r.keyRow(q); ok {
			t.Fatalf("%s: evidence was recorded", tc.name)
		}
	}
	for _, drop := range [][]string{{`DROP TABLE key_decisions`, `PRAGMA user_version = 7`}, // v7: the writer's K1, then v4: the judge's store step
		{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`}} {
		setVersion(t, r.d.credentialsPath, drop...)
		requireRefusal(t, drop[len(drop)-1], r.postKey(q, good), 503, "Restart")
	}
}
