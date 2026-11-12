package daemon

// Row 201 A3 sprint 2 M5a (design §2.4, §6 RT4, RT5): POST /surface/approve/{id} judges with the
// A2 law before it records, and each refusal writes nothing.

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// apPasskey is the enrolled passkey (mark) of the POST tests, at the loopback origin apOrigin.
var apPasskey = approvaltest.NewPasskey(approval.AlgEdDSA, 71)

const apOrigin = "http://localhost:7644"

// newPostRig is a rig with apPasskey enrolled as mark and one pending request.
func newPostRig(t *testing.T) (approvalRig, store.MintRequestRow) {
	t.Helper()
	a := newApprovalRig(t)
	a.enrolPasskey("mark", apPasskey, apOrigin)
	return a, a.putRequest("mr-00000000000000b1", "w201b", apT0-30)
}

// assertionFor is pk's assertion (counter n) at apOrigin over row's message.
func assertionFor(t *testing.T, pk *approvaltest.Passkey, row store.MintRequestRow, n uint32) []byte {
	t.Helper()
	d, _ := hex.DecodeString(row.Digest)
	c, _ := hex.DecodeString(row.Challenge)
	h, err := approval.MessageHash([32]byte(d), [32]byte(c), row.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	return pk.Assert(apOrigin, "localhost", h, n)
}

// approveBody is the page's POST body for digest with a webauthn assertion.
func approveBody(digest string, assertion []byte) []byte {
	b, _ := json.Marshal(map[string]any{"v": 1, "digest": digest, "webauthn": json.RawMessage(assertion)})
	return b
}

// post POSTs body to /surface/approve/<id> as the card page would (loopback Host and Origin,
// same-origin, JSON); hdr pairs override those headers ("" removes one).
func (a approvalRig) post(id string, body []byte, hdr ...string) *httptest.ResponseRecorder {
	a.t.Helper()
	return a.serve("POST", "/surface/approve/"+id, apLoopback, body,
		append([]string{"Origin", apOrigin, "Sec-Fetch-Site", "same-origin", "Content-Type", "application/json"}, hdr...)...)
}

// recordedAndSessions counts surface_approvals and session_credentials rows (straight from the file).
func (a approvalRig) recordedAndSessions() (recorded, sessions int) {
	a.t.Helper()
	db, err := sql.Open("sqlite", a.d.credentialsPath)
	if err != nil {
		a.t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	_ = db.QueryRow(`SELECT count(*) FROM surface_approvals`).Scan(&recorded) // below schema 7: 0
	if err := db.QueryRow(`SELECT count(*) FROM session_credentials`).Scan(&sessions); err != nil {
		a.t.Fatal(err)
	}
	return recorded, sessions
}

// RT4: R1-R14, each refusal with its status and class, and 0 surface_approvals rows and 0
// sessions after it. R2's 421 and R12 are RT7's and the law's (named equivalent, plan §4).
func TestApprovePOSTRefusalArms(t *testing.T) {
	other := approvaltest.NewPasskey(approval.AlgES256, 72)
	for _, arm := range []struct {
		name, id, class, reason string
		status                  int
		setup                   func(a approvalRig, row store.MintRequestRow) (body []byte, hdr []string)
	}{
		{"R1-id", "MR-00000000000000B1", "BadRequest", "mr- and 16 lowercase hex", 400, nil},
		{"R3-origin", "", "CrossOrigin", "http://localhost:7645", 403, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), []string{"Origin", "http://localhost:7645"}
		}},
		{"R4-site", "", "CrossOrigin", "Sec-Fetch-Site", 403, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), []string{"Sec-Fetch-Site", "same-site"}
		}},
		{"R5-type", "", "UnsupportedMediaType", "application/json", 415, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), []string{"Content-Type", "text/plain"}
		}},
		{"R6-size", "", "TooLarge", "20480", 413, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			return approveBody(row.Digest, []byte(`"`+strings.Repeat("x", surfacepage.MaxApproveBodyBytes)+`"`)), nil
		}},
		{"R6-grammar", "", "BadRequest", "both", 400, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			return []byte(`{"v":1,"digest":"` + row.Digest + `","webauthn":{},"sshsig":"x"}`), nil
		}},
		{"R7-unknown", "mr-00000000000000f1", "NotFound", "no mint request mr-00000000000000f1", 404, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), nil
		}},
		{"R8-expired", "", "Conflict", "expired at", 409, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			a.d.approveNow = func() int64 { return row.ExpiresAt }
			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), nil
		}},
		{"R8-bytes", "", "Conflict", "do not match its digest", 409, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			setRequestColumn(t, a.d.credentialsPath, row.RequestID, "digest", strings.Repeat("0", 64))
			return approveBody(strings.Repeat("0", 64), assertionFor(t, apPasskey, row, 2)), nil
		}},
		{"R9-digest", "", "Conflict", "changed since the card was shown", 409, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			return approveBody(strings.Repeat("ab", 32), assertionFor(t, apPasskey, row, 2)), nil
		}},
		{"R10-challenge", "", "Conflict", "challenge is not 64 hex", 409, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			setRequestColumn(t, a.d.credentialsPath, row.RequestID, "challenge", strings.Repeat("zz", 32))
			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), nil
		}},
		{"R11-other-request", "", "NotApproved", "does not approve request", 422, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			return approveBody(row.Digest, assertionFor(t, apPasskey, a.putRequest("mr-00000000000000c1", "w201c", apT0-20), 2)), nil
		}},
		{"R11-unenrolled", "", "NotApproved", "no enrolled approver", 422, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			return approveBody(row.Digest, assertionFor(t, other, row, 2)), nil
		}},
		{"R11-counter", "", "NotApproved", "does not approve request", 422, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 1)), nil
		}},
		{"R11-garbage", "", "NotApproved", "does not approve request", 422, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			return approveBody(row.Digest, []byte(`{"id":"x"}`)), nil
		}},
		{"R11-sshsig-garbage", "", "NotApproved", "the SSH signature does not approve request", 422, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			return sshBody(row.Digest, []byte("-----BEGIN SSH SIGNATURE-----\n")), nil
		}},
		{"R11-sshsig-unenrolled", "", "NotApproved", "no enrolled approver", 422, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			return sshBody(row.Digest, signatureFor(t, approvaltest.NewSigner(74), row)), nil
		}},
		{"R11-sshsig-store", "", "Internal", "internal store failure", 500, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			setVersion(t, a.d.credentialsPath, `ALTER TABLE approvers RENAME COLUMN public_key TO public_key_gone`)
			return sshBody(row.Digest, signatureFor(t, apSigner, row)), nil
		}},
		{"R11-store-v5", "", "Restart", "restart `ailang-worldd serve`", 503, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			body := approveBody(row.Digest, assertionFor(t, apPasskey, row, 2))
			setVersion(t, a.d.credentialsPath, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`)
			return body, nil
		}},
		{"R13-busy", "", "Busy", "busy", 503, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			held, err := store.OpenCredentials(a.d.credentialsPath)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = held.Close() })
			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), nil
		}},
		{"R14-v6", "", "Restart", "restart `ailang-worldd serve`", 503, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			setVersion(t, a.d.credentialsPath, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `PRAGMA user_version = 6`)
			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), nil
		}},
		{"R14-race", "", "Conflict", "changed while it was approved", 409, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
			reads := 0 // the record reads the clock again: the request expired meanwhile
			a.d.approveNow = func() int64 { reads++; return map[bool]int64{true: apT0, false: row.ExpiresAt}[reads == 1] }
			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), nil
		}},
	} {
		t.Run(arm.name, func(t *testing.T) {
			a, row := newPostRig(t)
			var body []byte
			var hdr []string
			if arm.setup != nil {
				body, hdr = arm.setup(a, row)
			}
			id := map[bool]string{true: row.RequestID, false: arm.id}[arm.id == ""]
			rec := a.post(id, body, hdr...)
			if r := requireRefusal(t, arm.name, rec, arm.status, arm.class); !strings.Contains(r.Reason, arm.reason) {
				t.Fatalf("reason %q lacks %q", r.Reason, arm.reason)
			}
			if arm.name == "R13-busy" && rec.Header().Get("Retry-After") != "1" {
				t.Fatalf("Retry-After = %q, want 1", rec.Header().Get("Retry-After"))
			}
			if n, s := a.recordedAndSessions(); n != 0 || s != 0 {
				t.Fatalf("%d recorded, %d sessions after a refusal; want 0 and 0", n, s)
			}
		})
	}
}

// RT5: a passkey assertion the law verifies is recorded (one row: webauthn, webauthn:mark, the
// posted envelope's bytes, the request's digest, received at the daemon's now) and nothing mints;
// requests.json then shows the approver. A refused one before it recorded nothing.
func TestApprovePOSTRecordsAVerifiedPasskey(t *testing.T) {
	a, row := newPostRig(t)
	if rec := a.post(row.RequestID, approveBody(row.Digest, assertionFor(t, apPasskey, row, 1))); rec.Code != 422 {
		t.Fatalf("a stale counter = %d, want 422", rec.Code)
	}
	if n, _ := a.recordedAndSessions(); n != 0 {
		t.Fatalf("a refused assertion was recorded (%d rows)", n)
	}
	assertion := assertionFor(t, apPasskey, row, 2)
	a.d.approveNow = func() int64 { return apT0 + 3 }
	rec := a.post(row.RequestID, approveBody(row.Digest, assertion))
	var got surfacepage.ApproveRecorded
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil || rec.Code != 200 || got != surfacepage.Recorded(row.RequestID, "webauthn:mark") {
		t.Fatalf("POST = %d %s", rec.Code, rec.Body.String())
	}
	if err := approvalHeaders(rec); err != nil {
		t.Fatal(err)
	}
	ro, err := store.OpenCredentialsReadOnly(a.d.credentialsPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ro.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 5e9)
	defer cancel()
	sa, ok, err := ro.SurfaceApproval(ctx, row.RequestID)
	if err != nil || !ok || sa.Format != "webauthn" || sa.Approver != "webauthn:mark" || !bytes.Equal(sa.Artifact, assertion) || sa.Digest != row.Digest || sa.ReceivedAt != apT0+3 {
		t.Fatalf("recorded = %+v %v %v", sa, ok, err)
	}
	if again, _, _ := ro.MintRequest(ctx, row.RequestID); again.Consumed {
		t.Fatal("the POST consumed the request")
	}
	if n, s := a.recordedAndSessions(); n != 1 || s != 0 {
		t.Fatalf("%d recorded, %d sessions; want 1 and 0", n, s)
	}
	if body := a.serve("GET", "/surface/requests.json", apLoopback, nil).Body.String(); !strings.Contains(body, `"recorded":"webauthn:mark"`) {
		t.Fatalf("requests.json does not show the approver: %s", body)
	}
}
