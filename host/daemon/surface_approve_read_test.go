package daemon

// Row 201 A3 sprint 2 M4e (design §2.3, §6 RT2, RT3): what the read routes say.

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// RT2: requests.json lists the pending requests newest first with their recorded approver or
// null, a row that does not decode as refused, and the decided ones (at most 20, newest first);
// it carries no request bytes, challenge, artifact or token; a store it cannot read is 200 with
// unavailable and empty lists.
func TestRequestsJSON(t *testing.T) {
	a := newApprovalRig(t)
	a.putRequest("mr-00000000000000a1", "w201a", apT0-700) // expired at apT0-100: not listed
	r1 := a.putRequest("mr-00000000000000b1", "w201b", apT0-30)
	r2 := a.putRequest("mr-00000000000000c1", "w201c", apT0-20)
	bad := a.putRequest("mr-00000000000000d1", "w201d", apT0-10)
	artifact := []byte("ARTIFACT-NEVER-LISTED")
	a.write(func(ctx context.Context, c *store.Credentials) error {
		return c.RecordSurfaceApproval(ctx, store.SurfaceApprovalRow{RequestID: r2.RequestID, Format: "webauthn", Approver: "webauthn:mark", Artifact: artifact, Digest: r2.Digest}, apT0-5)
	})
	setRequestColumn(t, a.d.credentialsPath, bad.RequestID, "request", append(bad.Request, ' ')) // no longer hashes to its digest
	for i := 0; i < 21; i++ {
		row := a.putRequest("mr-0000000000000e"+hex.EncodeToString([]byte{byte(i)}), "w201e", apT0-1000+int64(i))
		a.consume(row, apT0-900+int64(i), strings.Repeat(hex.EncodeToString([]byte{byte(i)}), 32))
	}
	rec := a.serve("GET", "/surface/requests.json", apLoopback, nil)
	if err := approvalHeaders(rec); rec.Code != 200 || err != nil {
		t.Fatalf("requests.json = %d, %v", rec.Code, err)
	}
	var got surfacepage.RequestsJSON
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.V != 1 || got.Now != apT0 || got.Unavailable != "" || len(got.Pending) != 3 || len(got.Decided) != surfacepage.MaxDecidedListed {
		t.Fatalf("requests.json = %+v", got)
	}
	p := got.Pending
	if p[0].ID != bad.RequestID || p[0].Refused == "" || p[0].Agent != "" || p[1].ID != r2.RequestID || p[1].Recorded == nil || *p[1].Recorded != "webauthn:mark" ||
		p[2].ID != r1.RequestID || p[2].Recorded != nil || p[2].Grants != 2 || p[2].Agent != "claude" || p[2].Hash16 == "" || p[2].Expires != apT0+570 {
		t.Fatalf("pending = %+v", p)
	}
	if d := got.Decided; d[0].ID != "mr-0000000000000e14" || d[0].ApprovedBy != "tty" || d[0].At != apT0-880 || d[0].Credential16 != "1414141414141414" || d[19].ID != "mr-0000000000000e01" {
		t.Fatalf("decided = %+v … %+v", d[0], d[19])
	}
	for _, leak := range []string{`"request"`, `"challenge"`, `"artifact"`, `"token"`, string(artifact), r1.Challenge, r1.Digest, "/ws/", "session"} {
		if strings.Contains(rec.Body.String(), leak) {
			t.Errorf("requests.json carries %q", leak)
		}
	}
	setVersion(t, a.d.credentialsPath, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `PRAGMA user_version = 6`) // below 7: listed, nothing recorded
	if body := a.serve("GET", "/surface/requests.json", apLoopback, nil).Body.String(); strings.Count(body, `"recorded":null`) != 3 || !strings.Contains(body, `"unavailable":""`) {
		t.Fatalf("a schema-6 store: %s", body)
	}
	if err := os.Remove(a.d.credentialsPath); err != nil {
		t.Fatal(err)
	}
	rec = a.serve("GET", "/surface/requests.json", apLoopback, nil)
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), `"pending":[],"decided":[],"unavailable":"the mint requests could not be read"`) {
		t.Fatalf("a missing store: %d %s", rec.Code, rec.Body.String())
	}
}

// RT3: the card is the stored bytes (base64url), the row's digest, nonce and times, and the active
// passkeys by origin (a revoked one left out); it has no message-hash field. A bad id is 400, an
// unknown one 404, a consumed one says so.
func TestRequestCardJSON(t *testing.T) {
	a := newApprovalRig(t)
	row := a.putRequest("mr-00000000000000b1", "w201b", apT0-30)
	a.enrolPasskey("a-mac", approvaltest.NewPasskey(approval.AlgEdDSA, 61), "http://localhost:7644")
	a.enrolPasskey("b-phone", approvaltest.NewPasskey(approval.AlgES256, 62), "https://studio.tail0.ts.net")
	a.enrolPasskey("c-key", approvaltest.NewPasskey(approval.AlgEdDSA, 63), "http://localhost:7644")
	a.enrolPasskey("d-gone", approvaltest.NewPasskey(approval.AlgEdDSA, 64), "https://gone.tail0.ts.net")
	a.write(func(ctx context.Context, c *store.Credentials) error { return c.RevokeApprover(ctx, "d-gone", apT0-5) })
	rec := a.serve("GET", "/surface/requests/"+row.RequestID, apLoopback, nil)
	if err := approvalHeaders(rec); rec.Code != 200 || err != nil {
		t.Fatalf("card = %d, %v", rec.Code, err)
	}
	var keys map[string]json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &keys); err != nil {
		t.Fatal(err)
	}
	var names []string
	for k := range keys {
		names = append(names, k)
	}
	sort.Strings(names)
	if want := []string{"approvedBy", "challenge", "consumed", "created", "digest", "expires", "id", "passkeys", "request", "ssh", "v"}; !reflect.DeepEqual(names, want) {
		t.Fatalf("card keys %v, want %v", names, want)
	}
	var c surfacepage.Card
	if err := json.Unmarshal(rec.Body.Bytes(), &c); err != nil {
		t.Fatal(err)
	}
	b, err := base64.RawURLEncoding.DecodeString(c.Request)
	b64 := func(pk *approvaltest.Passkey) string { return base64.RawURLEncoding.EncodeToString(pk.CredentialID()) }
	want := []surfacepage.CardPasskeys{{Origin: "http://localhost:7644", RP: "localhost", Allow: []string{b64(approvaltest.NewPasskey(approval.AlgEdDSA, 61)), b64(approvaltest.NewPasskey(approval.AlgEdDSA, 63))}},
		{Origin: "https://studio.tail0.ts.net", RP: "studio.tail0.ts.net", Allow: []string{b64(approvaltest.NewPasskey(approval.AlgES256, 62))}}}
	if err != nil || !bytes.Equal(b, row.Request) || c.Digest != row.Digest || c.Challenge != row.Challenge || c.Created != row.CreatedAt ||
		c.Expires != row.ExpiresAt || c.Consumed || !reflect.DeepEqual(c.Passkeys, want) || c.SSH.Namespace != approval.Namespace {
		t.Fatalf("card = %+v", c)
	}
	d, _ := hex.DecodeString(row.Digest)
	n, _ := hex.DecodeString(row.Challenge)
	h, err := approval.MessageHash([32]byte(d), [32]byte(n), row.ExpiresAt)
	if err != nil || strings.Contains(rec.Body.String(), hex.EncodeToString(h[:])) || strings.Contains(rec.Body.String(), base64.RawURLEncoding.EncodeToString(h[:])) {
		t.Fatalf("the card carries the message hash (%v)", err)
	}
	requireRefusal(t, "bad id", a.serve("GET", "/surface/requests/MR-00000000000000B1", apLoopback, nil), 400, "BadRequest")
	requireRefusal(t, "short id", a.serve("GET", "/surface/requests/mr-00", apLoopback, nil), 400, "BadRequest")
	requireRefusal(t, "unknown", a.serve("GET", "/surface/requests/mr-00000000000000f1", apLoopback, nil), 404, "NotFound")
	a.consume(row, apT0-10, strings.Repeat("ab", 32))
	if err := json.Unmarshal(a.serve("GET", "/surface/requests/"+row.RequestID, apLoopback, nil).Body.Bytes(), &c); err != nil || !c.Consumed || c.ApprovedBy != "tty" {
		t.Fatalf("consumed card = %+v, %v", c, err)
	}
	setVersion(t, a.d.credentialsPath, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`)
	if rec := a.serve("GET", "/surface/requests/"+row.RequestID, apLoopback, nil); rec.Code != 200 || !strings.Contains(rec.Body.String(), `"passkeys":[]`) {
		t.Fatalf("a schema-5 store's card = %d %s, want 200 with no passkeys", rec.Code, rec.Body.String())
	}
	if err := os.Remove(a.d.credentialsPath); err != nil {
		t.Fatal(err)
	}
	requireRefusal(t, "no store", a.serve("GET", "/surface/requests/"+row.RequestID, apLoopback, nil), 500, "Internal")
}
