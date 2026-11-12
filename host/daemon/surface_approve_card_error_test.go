package daemon

// Row 201 A3 sprint 3 M6a (design §2.3; iteration 385's judge survivor J36): the card route's
// passkey read failing is the sanitized 500, never a card that lists no passkeys.

import (
	"strings"
	"testing"
)

// TestRequestCardPasskeyReadError: a schema-7 store whose webauthn_credentials table lost a column
// (the table set still opens, the request row still reads, and ActivePasskeys fails, not as an
// older store): the card GET is 500 Internal with no card in the body. Control: 200 before.
func TestRequestCardPasskeyReadError(t *testing.T) {
	a := newApprovalRig(t)
	row := a.putRequest("mr-00000000000000e1", "w201e", apT0-30)
	target := "/surface/requests/" + row.RequestID
	if rec := a.serve("GET", target, apLoopback, nil); rec.Code != 200 || !strings.Contains(rec.Body.String(), `"passkeys":[]`) {
		t.Fatalf("control: the card = %d %s", rec.Code, rec.Body.String())
	}
	setVersion(t, a.d.credentialsPath, `ALTER TABLE webauthn_credentials RENAME COLUMN origin TO origin_gone`) // still version 7
	rec := a.serve("GET", target, apLoopback, nil)
	requireRefusal(t, "a passkey read error", rec, 500, "Internal")
	for _, key := range []string{`"passkeys"`, `"request"`, `"challenge"`, row.Digest} {
		if strings.Contains(rec.Body.String(), key) {
			t.Fatalf("the 500 carries %s: %s", key, rec.Body.String())
		}
	}
}
