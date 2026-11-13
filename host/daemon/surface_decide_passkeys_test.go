package daemon

// Row 201 A4 sprint 4 (design §2.8; sprint-1 plan P7): the subject's passkeys.

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// The subject's passkeys (sprint-1 P7, sprint 4): the active passkeys grouped by origin, as the
// request card's; none enrolled is an empty list; a credentials store that cannot be opened is the
// sanitized 500, never a subject without its passkeys.
func TestDecisionSubjectPasskeys(t *testing.T) {
	r := newDecideRig(t)
	q, _ := r.packet("passkeys", broker.WorkspaceScope, "agent")
	subject := func() (int, surfacepage.DecisionSubject, string) {
		rec := r.serve("GET", "/surface/decisions/"+q.String(), apLoopback, nil)
		var s surfacepage.DecisionSubject
		_ = json.Unmarshal(rec.Body.Bytes(), &s)
		return rec.Code, s, rec.Body.String()
	}
	if code, s, body := subject(); code != 200 || s.Passkeys == nil || len(s.Passkeys) != 0 || !strings.Contains(body, `"passkeys":[]`) {
		t.Fatalf("no passkey: %d %s", code, body)
	}
	a, b, c := approvaltest.NewPasskey(-7, 1), approvaltest.NewPasskey(-7, 2), approvaltest.NewPasskey(-7, 3)
	r.enrolPasskey("mark-mac", a, "https://studio.tail.ts.net")
	r.enrolPasskey("mark-phone", b, "http://localhost:7644")
	r.enrolPasskey("mark-key", c, "https://studio.tail.ts.net")
	r.write(func(ctx context.Context, cr *store.Credentials) error {
		return cr.RevokeApprover(ctx, "mark-key", apT0)
	})
	code, s, body := subject()
	got := map[string][]string{}
	for _, g := range s.Passkeys {
		sort.Strings(g.Allow)
		got[g.Origin+" "+g.RP] = g.Allow
	}
	enc := func(p *approvaltest.Passkey) string { return base64.RawURLEncoding.EncodeToString(p.CredentialID()) }
	want := map[string][]string{"https://studio.tail.ts.net studio.tail.ts.net": {enc(a)}, "http://localhost:7644 localhost": {enc(b)}}
	if code != 200 || len(s.Passkeys) != 2 || !reflect.DeepEqual(got, want) {
		t.Fatalf("passkeys %d %v, want %v\n%s", code, got, want, body)
	}
	r.d.credentialsPath = filepath.Join(t.TempDir(), "missing", "credentials.db")
	requireRefusal(t, "no credentials store", r.serve("GET", "/surface/decisions/"+q.String(), apLoopback, nil), 500, "Internal")
}
