package daemon

// Row 201 A3 sprint 3 M5d (design §2.4, §6 RT6, RT13; the sprint-2 plan's Appendix C): the POST's
// SSHSIG arm judges an approvaltest signature with the A1b law before it records.

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/store"
)

// apSigner is the SSH approver (mark-ssh) of the SSHSIG tests.
var apSigner = approvaltest.NewSigner(73)

// enrolSSH enrols s as name at apT0-10: after the tests' requests (apT0-30), as enrolPasskey.
func (a approvalRig) enrolSSH(name string, s *approvaltest.Signer) {
	a.t.Helper()
	a.write(func(ctx context.Context, c *store.Credentials) error {
		return c.EnrolApprover(ctx, store.ApproverRow{Name: name, Format: approval.FormatSSH, PublicKey: s.PublicKey().Blob, EnrolledAt: apT0 - 10})
	})
}

// signatureFor is s's armored SSHSIG over row's 90-byte message.
func signatureFor(t *testing.T, s *approvaltest.Signer, row store.MintRequestRow) []byte {
	t.Helper()
	d, _ := hex.DecodeString(row.Digest)
	c, _ := hex.DecodeString(row.Challenge)
	msg, err := approval.SignedMessage([32]byte(d), [32]byte(c), row.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	return s.Sign(approval.Namespace, msg)
}

// sshBody is the page's POST body for digest with an armored SSHSIG.
func sshBody(digest string, armored []byte) []byte {
	b, _ := json.Marshal(map[string]any{"v": 1, "digest": digest, "sshsig": string(armored)})
	return b
}

// recorded is request id's surface_approvals row.
func (a approvalRig) recorded(id string) (store.SurfaceApprovalRow, bool) {
	a.t.Helper()
	ro, err := store.OpenCredentialsReadOnly(a.d.credentialsPath)
	if err != nil {
		a.t.Fatal(err)
	}
	defer func() { _ = ro.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 5e9)
	defer cancel()
	sa, ok, err := ro.SurfaceApproval(ctx, id)
	if err != nil {
		a.t.Fatal(err)
	}
	return sa, ok
}

// RT6: an enrolled key's SSHSIG over the request's message is recorded (ssh, ssh:mark-ssh, the
// armored bytes, the request's digest, received at the daemon's now); nothing mints.
func TestApprovePOSTRecordsAVerifiedSSHSIG(t *testing.T) {
	a, row := newPostRig(t)
	a.enrolSSH("mark-ssh", apSigner)
	sig := signatureFor(t, apSigner, row)
	a.d.approveNow = func() int64 { return apT0 + 4 }
	rec := a.post(row.RequestID, sshBody(row.Digest, sig))
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), `"approvedBy":"ssh:mark-ssh"`) {
		t.Fatalf("POST = %d %s", rec.Code, rec.Body.String())
	}
	sa, ok := a.recorded(row.RequestID)
	if !ok || sa.Format != "ssh" || sa.Approver != "ssh:mark-ssh" || !bytes.Equal(sa.Artifact, sig) || sa.Digest != row.Digest || sa.ReceivedAt != apT0+4 {
		t.Fatalf("recorded = %+v %v", sa, ok)
	}
	if n, s := a.recordedAndSessions(); n != 1 || s != 0 {
		t.Fatalf("%d recorded, %d sessions; want 1 and 0", n, s)
	}
	if body := a.serve("GET", "/surface/requests.json", apLoopback, nil).Body.String(); !strings.Contains(body, `"recorded":"ssh:mark-ssh"`) {
		t.Fatalf("requests.json does not show the approver: %s", body)
	}
}

// RT13: one SSH approver and no passkey; with --surface-origin https://studio.tail0.ts.net the
// card and the SSHSIG POST work through that origin under both its own Host and the loopback one
// (a proxy); without it they are refused naming the flag (403 via the loopback Host, 421 via the
// ts.net Host), and another ts.net origin is refused even when one is configured.
func TestApprovePOSTSSHSIGViaTailnetOriginWithNoPasskey(t *testing.T) {
	const ts, tsHost = "https://studio.tail0.ts.net", "studio.tail0.ts.net"
	rig := func(configured bool) (approvalRig, store.MintRequestRow, []byte) {
		a := newApprovalRig(t)
		a.enrolSSH("mark-ssh", apSigner)
		if configured {
			a.d.cfg.SurfaceOrigins = []string{ts}
		}
		row := a.putRequest("mr-00000000000000d1", "w201d", apT0-30)
		return a, row, sshBody(row.Digest, signatureFor(t, apSigner, row))
	}
	a, row, body := rig(true)
	for _, host := range []string{tsHost, apLoopback} {
		if rec := a.serve("GET", "/surface/requests/"+row.RequestID, host, nil); rec.Code != 200 {
			t.Fatalf("card via Host %s = %d %s", host, rec.Code, rec.Body.String())
		}
		rec := a.serve("POST", "/surface/approve/"+row.RequestID, host, body, "Origin", ts, "Sec-Fetch-Site", "same-origin", "Content-Type", "application/json")
		if sa, ok := a.recorded(row.RequestID); rec.Code != 200 || !ok || sa.Approver != "ssh:mark-ssh" {
			t.Fatalf("POST via Host %s = %d %s; recorded %+v", host, rec.Code, rec.Body.String(), sa)
		}
	}
	if n, _ := a.recordedAndSessions(); n != 1 {
		t.Fatalf("%d rows after two POSTs of one request, want 1", n)
	}
	other := a.serve("POST", "/surface/approve/"+row.RequestID, apLoopback, body, "Origin", "https://other.tail0.ts.net", "Sec-Fetch-Site", "same-origin", "Content-Type", "application/json")
	requireRefusal(t, "another ts.net origin", other, 403, "CrossOrigin")

	a, row, body = rig(false)
	r := requireRefusal(t, "unconfigured", a.serve("POST", "/surface/approve/"+row.RequestID, apLoopback, body, "Origin", ts, "Sec-Fetch-Site", "same-origin", "Content-Type", "application/json"), 403, "CrossOrigin")
	if !strings.Contains(r.Reason, "--surface-origin "+ts) {
		t.Fatalf("the refusal %q does not name the flag", r.Reason)
	}
	requireRefusal(t, "unconfigured ts.net Host", a.serve("POST", "/surface/approve/"+row.RequestID, tsHost, body, "Origin", ts, "Sec-Fetch-Site", "same-origin", "Content-Type", "application/json"), 421, "MisdirectedRequest")
	requireRefusal(t, "unconfigured card", a.serve("GET", "/surface/requests/"+row.RequestID, tsHost, nil), 421, "MisdirectedRequest")
	if n, s := a.recordedAndSessions(); n != 0 || s != 0 {
		t.Fatalf("%d recorded, %d sessions unconfigured; want 0 and 0", n, s)
	}
}
