package main

// Row 201 A4 sprint 4 (design §2.8, §6 XD4): the decision card's SSH path, in a live Chrome, gated as
// XB1 (WORLD_CHROME, WORLD_REQUIRE_PASSKEY_BROWSER=1; skipped and counted without them); its
// real-tool sibling needs ssh-keygen (WORLD_REQUIRE_SSH_KEYGEN=1 makes a missing one a failure).

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
)

// XD4: the card's SSH path. Save the approve message and Save the deny message download exactly Go's
// approval.DecisionMessage (91 bytes) for their verbs and the window's max, under the names the
// signing lines give; approvaltest's SSHSIG over the approve message, pasted, decides approve as
// ssh:mark-mac (200); a signature over the approve message posted as a deny is refused by the law
// (422) and decides nothing. Real-tool sibling (directive retro 11): ssh-keygen -Y sign over the
// saved deny file with a key made in a temp dir decides deny as ssh:mark-key.
func TestDecidePageSSHSIGPaste(t *testing.T) {
	b := startPasskeyBrowser(t)
	b.do(t, nil, "tab", map[string]any{"authenticator": "uv"})
	r := newDecideCardRig(t)
	seedApprover(t, r.creds, "mark-mac", sigSigner.PublicKey(), 1000)
	for i, scope := range []string{"worktree:a", "worktree:b"} {
		req, pkt := cardObjects(scope, "agent", 2400)
		s := cardSubject(req, pkt, 1710, 2400)
		note := "No passkey is enrolled; decide with an SSH key below."
		if i == 0 { // a passkey enrolled only at another origin: the card names it and offers no passkey
			s.Passkeys = cardPasskeys("https://studio.tail0.ts.net", 93)
			note = "No passkey is enrolled for " + r.origin + ": open this card at https://studio.tail0.ts.net."
		}
		url := r.add(s, scope)
		var view cardView
		b.do(t, &view, "openCard", map[string]any{"url": url})
		ref16 := strings.TrimPrefix(s.Request, "sha256:")[:16]
		if view.State != "ready" || !view.ButtonHidden || view.Note != note || view.SaveHidden {
			t.Fatalf("%d: the card shows %+v", i, view)
		}
		dl := t.TempDir()
		msgs := map[string][]byte{}
		for id, verb := range map[string]byte{"save": approval.VerbApprove, "save-deny": approval.VerbDeny} {
			var downloads []struct{ Event, State, Name, URL string }
			b.do(t, &downloads, "save", map[string]any{"id": id, "dir": dl})
			name := map[byte]string{approval.VerbApprove: "approve", approval.VerbDeny: "deny"}[verb]
			got, err := os.ReadFile(filepath.Join(dl, "world-decision-"+name+"-"+ref16+".msg"))
			want, _ := approval.DecisionMessage(verb, refDigest(s.Packet), refDigest(s.Request), s.Window.Max)
			if err != nil || len(want) != 91 || !bytes.Equal(got, want) || !strings.HasPrefix(downloads[0].URL, "blob:"+r.origin+"/") {
				t.Fatalf("%d: %s saved %x (%v), want %x; %+v", i, id, got, err, want, downloads)
			}
			msgs[name] = got
		}
		if i == 1 {
			ssh, err := exec.LookPath("ssh-keygen")
			if err != nil {
				if os.Getenv("WORLD_REQUIRE_SSH_KEYGEN") == "1" {
					t.Fatal("WORLD_REQUIRE_SSH_KEYGEN=1 but ssh-keygen was not found")
				}
				t.Log("SKIP COUNTED: ssh-keygen not found; the card's real-tool sibling did NOT run")
				return
			}
			_ = ssh
			t.Setenv("HOME", t.TempDir())
			t.Setenv("SSH_AUTH_SOCK", "")
			key := filepath.Join(t.TempDir(), "k")
			sshKeygen(t, "-q", "-t", "ed25519", "-N", "", "-C", "xd4", "-f", key)
			pub, err := os.ReadFile(key + ".pub")
			if err != nil {
				t.Fatal(err)
			}
			k, err := approval.ParseAuthorizedKey(pub)
			if err != nil {
				t.Fatal(err)
			}
			seedApprover(t, r.creds, "mark-key", k, 1000)
			file := filepath.Join(dl, "world-decision-deny-"+ref16+".msg")
			sshKeygen(t, "-Y", "sign", "-n", approval.DecisionNamespace, "-f", key, file)
			sig, err := os.ReadFile(file + ".sig")
			if err != nil {
				t.Fatal(err)
			}
			b.do(t, nil, "fill", map[string]any{"id": "sig", "value": string(sig)})
			b.do(t, &view, "press", map[string]any{"id": "post-sig-deny", "states": []string{"decided", "rejected", "failed", "unknown"}})
			if view.State != "decided" || !strings.HasPrefix(view.Status, "Denied by ssh:mark-key: decision sha256:") || s.Decision != "deny" || s.DecidedBy != "ssh:mark-key" {
				t.Fatalf("ssh-keygen: the page shows %+v; subject %+v", view, s)
			}
			continue
		}
		sig := sigSigner.Sign(approval.DecisionNamespace, msgs["approve"])
		b.do(t, nil, "fill", map[string]any{"id": "sig", "value": string(sig)})
		b.do(t, &view, "press", map[string]any{"id": "post-sig-deny", "states": []string{"decided", "rejected", "failed", "unknown"}})
		if view.State != "rejected" || !strings.Contains(view.Status, "HTTP 422 NotApproved") || s.Decision != "" {
			t.Fatalf("the approve signature posted as a deny: %+v", view)
		}
		b.do(t, &view, "press", map[string]any{"id": "post-sig", "states": []string{"decided", "rejected", "failed", "unknown"}})
		p := r.posts[len(r.posts)-1]
		if view.State != "decided" || !strings.HasPrefix(view.Status, "Approved by ssh:mark-mac: decision sha256:") || s.Decision != "approve" ||
			p.body.SSHSIG == nil || *p.body.SSHSIG != string(sig) || p.body.Expiry != s.Window.Max || p.body.WebAuthn != nil {
			t.Fatalf("the approve signature: %+v; post %+v", view, p)
		}
	}
}
