package main

// Row 201 A4 sprint 4 (design §2.8, §6 XD1): the decision card decides with a passkey in a live Chrome, gated as
// XB1 (WORLD_CHROME, WORLD_REQUIRE_PASSKEY_BROWSER=1; skipped and counted without them).

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// cardPasskeys is the subject's passkeys for approvaltest's ES256 seed at origin (rpID localhost).
func cardPasskeys(origin string, seed byte) []surfacepage.CardPasskeys {
	id := approvaltest.NewPasskey(approval.AlgES256, seed).CredentialID()
	return []surfacepage.CardPasskeys{{Origin: origin, RP: "localhost", Allow: []string{base64.RawURLEncoding.EncodeToString(id)}}}
}

// XD1: mark-mac's passkey (Chrome's virtual authenticator, approvaltest's key, enrolled in the
// credentials store) decides two packets from their cards: Approve on the first is one decision
// approve by webauthn:mark-mac, Deny on the second one deny. Each get's challenge is the decision
// message's hash for its own verb and the window's max, with the subject's rpID, its credential and
// user verification; each POST carries the verb, P, Q, the expiry and the assertion, with A3a's
// browser headers and no cookie; the page says the verb it did and then reads the request again,
// decided, with no decide control. With less than a minute of the window left, the page reads the
// subject again before it signs. No CSP violation, no page error.
func TestDecidePageApproveWithPasskey(t *testing.T) {
	b := startPasskeyBrowser(t)
	b.do(t, nil, "tab", map[string]any{"authenticator": "uv", "inject": "getRecorder"})
	r := newDecideCardRig(t)
	virtualPasskey(t, b, r.creds, r.origin, "mark-mac", 91, 1000)
	for i, tc := range []struct {
		verb, button, said string
		v                  byte
		now, deadline      int64
		gets               int
	}{
		{"approve", "go", "Approved by webauthn:mark-mac: decision sha256:", approval.VerbApprove, 1710, 2400, 1},
		{"deny", "deny", "Denied by webauthn:mark-mac: decision sha256:", approval.VerbDeny, 1710, 2400, 1},
		{"approve", "go", "Approved by webauthn:mark-mac: decision sha256:", approval.VerbApprove, 2350, 2400, 2}, // 50 s left: read again
	} {
		req, pkt := cardObjects(broker.WorkspaceScope, fmt.Sprintf("agent-%d", i), tc.deadline)
		s := cardSubject(req, pkt, tc.now, tc.deadline)
		s.Passkeys = cardPasskeys(r.origin, 91)
		url := r.add(s, broker.WorkspaceScope)
		var view cardView
		b.do(t, &view, "openCard", map[string]any{"url": url})
		ref16 := strings.TrimPrefix(s.Request, "sha256:")[:16]
		if view.State != "ready" || view.ButtonHidden || view.SaveHidden || !strings.HasPrefix(view.What, "Decide request "+s.Request+": Human.Approve for "+broker.WorkspaceScope+", requested by agent-") ||
			view.Fields["request.scope"] != broker.WorkspaceScope || view.Fields["packet.requestRef"] != s.Request || view.Fields["packet.policy"] != "Cancel" ||
			view.Fields["packet.deadlineAt"] != fmt.Sprintf("1970-01-01T00:%02d:%02d.000Z (%d)", tc.deadline/60, tc.deadline%60, tc.deadline) || len(view.Fields) != 13 ||
			view.SSHLine != "ssh-keygen -Y sign -n world-decision@sunholo -f <key> world-decision-approve-"+ref16+".msg\nssh-keygen -Y sign -n world-decision@sunholo -f <key> world-decision-deny-"+ref16+".msg" {
			t.Fatalf("%d: the card shows %+v", i, view)
		}
		var label, href string
		b.do(t, &label, "waitFor", map[string]any{"expression": "document.getElementById('" + tc.button + "').textContent"})
		b.do(t, &href, "waitFor", map[string]any{"expression": "document.getElementById('proposal').getAttribute('href')"})
		if label != map[string]string{"approve": "Approve with passkey", "deny": "Deny with passkey"}[tc.verb] || href != "/workbench?object=sha256%3A"+hashref.SumSHA256([]byte("plan")).Digest() {
			t.Fatalf("%d: the button reads %q, the proposal link %q", i, label, href)
		}
		b.do(t, &view, "press", map[string]any{"id": tc.button, "states": []string{"decided", "rejected", "failed", "unknown", "refused"}})
		if view.State != "decided" || !strings.HasPrefix(view.Status, tc.said) || !view.ButtonHidden || !view.SaveHidden || len(view.Violations) != 0 || len(view.Errors) != 0 {
			t.Fatalf("%d: after %s the page shows %+v", i, tc.verb, view)
		}
		g := view.Gets[len(view.Gets)-1]
		if len(view.Gets) != 1 || g.Challenge != decisionHash64(t, tc.v, s, s.Window.Max) || g.RPID != "localhost" || g.UserVerification != "required" ||
			len(g.Allow) != 1 || g.Allow[0] != s.Passkeys[0].Allow[0] || g.Timeout != 240000 || !g.Signal {
			t.Fatalf("%d: get %+v (of %d), want the %s hash", i, g, len(view.Gets), tc.verb)
		}
		p := r.posts[len(r.posts)-1]
		if len(r.posts) != i+1 || p.status != 200 || p.body.Verb != tc.verb || p.body.Packet != strings.TrimPrefix(s.Packet, "sha256:") ||
			p.body.Request != strings.TrimPrefix(s.Request, "sha256:") || p.body.Expiry != s.Window.Max || p.body.WebAuthn == nil || p.body.SSHSIG != nil || p.body.Typed ||
			s.Decision != tc.verb || s.DecidedBy != "webauthn:mark-mac" {
			t.Fatalf("%d: posts %d, last %+v; subject %+v", i, len(r.posts), p, s)
		}
		if gets := r.gets[s.Request]; gets != tc.gets+1 { // the reads before signing, and one after the decision
			t.Fatalf("%d: the subject was read %d times, want %d", i, gets, tc.gets+1)
		}
	}
}
