package daemon

// Row 243 V3 sprint 4 M10 (design §2.6, §6 CX0-CX13): POST /surface/confirm/{id}'s guards and its
// refusals before any write. Requests are S3's putVerb rows (approval.CanonicalVerb of verbFixture);
// signatures are approvaltest's over approval.VerbMessage of the row; every refusal is compared with
// the stores before it (nothing written). Stores and keys are tests' (D-WORLD-109).

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/store"
)

const cfEp, cfID = "w243", "vr-00000000000000c1"

// newConfirmRig enrols apSigner as mark-ssh and apPasskey as mark (counter 1) at apT0-10 and stores
// one pending draft request cfID for cfEp, created at apT0-30.
func newConfirmRig(t *testing.T) (approvalRig, store.VerbRequestRow) {
	t.Helper()
	a := newApprovalRig(t)
	a.enrolSSH("mark-ssh", apSigner)
	a.enrolPasskey("mark", apPasskey, apOrigin)
	return a, a.putVerb(cfID, cfEp, apT0-30, true)
}

// verbMessageOf is row's 86-byte world-verb/v1 message and its hash.
func verbMessageOf(t *testing.T, row store.VerbRequestRow) ([]byte, [32]byte) {
	t.Helper()
	d, _ := hex.DecodeString(row.Digest)
	c, _ := hex.DecodeString(row.Challenge)
	m, err := approval.VerbMessage([32]byte(d), [32]byte(c), row.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	h, err := approval.VerbMessageHash([32]byte(d), [32]byte(c), row.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	return m, h
}

// confirmSSH is the card's body for row with s's SSHSIG over its message under namespace.
func confirmSSH(t *testing.T, s *approvaltest.Signer, namespace string, row store.VerbRequestRow) []byte {
	t.Helper()
	m, _ := verbMessageOf(t, row)
	return sshBody(row.Digest, s.Sign(namespace, m))
}

// confirmPasskey is the card's body for row with pk's assertion (counter n) over its message hash.
func confirmPasskey(t *testing.T, pk *approvaltest.Passkey, row store.VerbRequestRow, n uint32) []byte {
	t.Helper()
	_, h := verbMessageOf(t, row)
	return approveBody(row.Digest, pk.Assert(apOrigin, "localhost", h, n))
}

// confirmAt POSTs body to /surface/confirm/<id> with Host host, as the card would (loopback Origin,
// same-origin, JSON); hdr pairs override those headers. Until M10d registers the route, the request
// reaches the handler through approvalRoute, as the mux will route it.
func (a approvalRig) confirmAt(host, id string, body []byte, hdr ...string) *httptest.ResponseRecorder {
	a.t.Helper()
	req := httptest.NewRequest("POST", "/surface/confirm/"+id, bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), http.LocalAddrContextKey, net.Addr(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 7644})))
	req.Host = host
	req.SetPathValue("id", id)
	for h, i := append([]string{"Origin", apOrigin, "Sec-Fetch-Site", "same-origin", "Content-Type", "application/json"}, hdr...), 0; i+1 < len(h); i += 2 {
		req.Header.Set(h[i], h[i+1])
	}
	rec := httptest.NewRecorder()
	a.d.approvalRoute(a.d.handleSurfaceConfirm).ServeHTTP(rec, req)
	return rec
}

// confirm is confirmAt on the loopback Host.
func (a approvalRig) confirm(id string, body []byte, hdr ...string) *httptest.ResponseRecorder {
	a.t.Helper()
	return a.confirmAt(apLoopback, id, body, hdr...)
}

// confirmState is the "nothing written" witness: every verb_requests row, each approver's last use
// and revocation, each passkey's counter, and the world store's approval decisions and head.
func (a approvalRig) confirmState() string {
	a.t.Helper()
	n, head := decideRig{a}.decisions()
	return fmt.Sprintf("%s decisions=%d head=%s", a.credState(), n, head)
}

// credState is confirmState's credentials half.
func (a approvalRig) credState() string {
	a.t.Helper()
	db, err := sql.Open("sqlite", a.d.credentialsPath)
	if err != nil {
		a.t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	var b strings.Builder
	for _, q := range []string{`SELECT request_id, coalesce(consumed_at, 0), coalesce(outcome, ''), coalesce(approved_by, ''), hex(artifact), coalesce(decision_ref, '') FROM verb_requests ORDER BY request_id`,
		`SELECT name, coalesce(last_used_at, 0), coalesce(revoked_at, 0), '', '', '' FROM approvers ORDER BY name`, `SELECT name, sign_count, '', '', '', '' FROM webauthn_credentials ORDER BY name`} {
		rows, err := db.Query(q)
		if err != nil { // a schema-8 store has no verb_requests
			fmt.Fprintf(&b, "%v;", err)
			continue
		}
		for rows.Next() {
			var c [6]string
			if err := rows.Scan(&c[0], &c[1], &c[2], &c[3], &c[4], &c[5]); err != nil {
				a.t.Fatal(err)
			}
			fmt.Fprintf(&b, "%v;", c)
		}
		_ = rows.Close()
	}
	return b.String()
}

// CX0-CX6: C0 (a session or a cookie), C1 (the id), R2 (the Host, approvalRoute), C3-C5 (A3a's
// guards), C6 (size, ParseApproveBody: no typed arm) each refuse with nothing written; the same
// signature then confirms (control).
func TestConfirmPOSTGuards(t *testing.T) {
	a, row := newConfirmRig(t)
	m, _ := verbMessageOf(t, row)
	sig := string(apSigner.Sign(approval.VerbNamespace, m))
	body := sshBody(row.Digest, []byte(sig))
	before := a.confirmState()
	for _, tc := range []struct {
		name, id      string
		body          []byte
		hdr           []string
		status        int
		class, reason string
	}{
		{"C0 authorization", cfID, body, []string{"Authorization", "Bearer x"}, 403, "Forbidden", "never made with a session"},
		{"C0 cookie", cfID, body, []string{"Cookie", "s=1"}, 403, "Forbidden", "never made with a session"},
		{"C0 before C1", "vr-0123", body, []string{"Cookie", "s=1"}, 403, "Forbidden", ""},
		{"C1 short", "vr-0123", body, nil, 400, "BadRequest", "vr- and 16 lowercase hex"},
		{"C1 upper", "vr-00000000000000C1", body, nil, 400, "BadRequest", ""},
		{"C1 mint id", "mr-00000000000000c1", body, nil, 400, "BadRequest", ""},
		{"C1 before C3", "vr-0123", body, []string{"Origin", "http://evil.example"}, 400, "BadRequest", ""},
		{"C3 missing", cfID, body, []string{"Origin", ""}, 403, "CrossOrigin", ""},
		{"C3 null", cfID, body, []string{"Origin", "null"}, 403, "CrossOrigin", ""},
		{"C3 foreign", cfID, body, []string{"Origin", "http://localhost:7645"}, 403, "CrossOrigin", ""},
		{"C4 same-site", cfID, body, []string{"Sec-Fetch-Site", "same-site"}, 403, "CrossOrigin", "Sec-Fetch-Site"},
		{"C4 cross-site", cfID, body, []string{"Sec-Fetch-Site", "cross-site"}, 403, "CrossOrigin", "Sec-Fetch-Site"},
		{"C5 text", cfID, body, []string{"Content-Type", "text/plain"}, 415, "UnsupportedMediaType", ""},
		{"C6 oversize", cfID, append(body[:len(body)-1:len(body)-1], []byte(`,"x":"`+strings.Repeat("a", 20<<10)+`"}`)...), nil, 413, "TooLarge", ""},
		{"C6 typed", cfID, []byte(`{"v":1,"digest":"` + row.Digest + `","typed":"y"}`), nil, 400, "BadRequest", `unknown field "typed"`},
		{"C6 y", cfID, []byte(`{"v":1,"digest":"` + row.Digest + `","sshsig":` + fmt.Sprintf("%q", sig) + `,"y":true}`), nil, 400, "BadRequest", "unknown field"},
		{"C6 both formats", cfID, []byte(`{"v":1,"digest":"` + row.Digest + `","sshsig":` + fmt.Sprintf("%q", sig) + `,"webauthn":{}}`), nil, 400, "BadRequest", "both"},
		{"C6 neither", cfID, []byte(`{"v":1,"digest":"` + row.Digest + `"}`), nil, 400, "BadRequest", "neither"},
		{"C6 v2", cfID, []byte(`{"v":2,"digest":"` + row.Digest + `","sshsig":` + fmt.Sprintf("%q", sig) + `}`), nil, 400, "BadRequest", "v is 2"},
		{"C6 two values", cfID, append(append([]byte{}, body...), body...), nil, 400, "BadRequest", "more than one"},
	} {
		got := requireRefusal(t, tc.name, a.confirm(tc.id, tc.body, tc.hdr...), tc.status, tc.class)
		if !strings.Contains(got.Reason, tc.reason) {
			t.Errorf("%s: reason %q, want %q", tc.name, got.Reason, tc.reason)
		}
	}
	requireRefusal(t, "R2 foreign host", a.confirmAt("rebind.example:7644", cfID, body), 421, "MisdirectedRequest")
	if after := a.confirmState(); after != before {
		t.Fatalf("a guard wrote:\n%s\n%s", before, after)
	}
	if rec := a.confirm(cfID, body); rec.Code != 200 {
		t.Fatalf("control: %d %s", rec.Code, rec.Body.String())
	}
}
