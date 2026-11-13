package daemon

// Row 201 A4 sprint 3 M6d (design §2.7 [r2], §2.8, §6 RD10 and RD11 passkey halves, KD2 and KD7 end
// to end, KD8, RD13 (b) (c), RD14): POST /surface/decide/{ref} with a passkey assertion by an
// enrolled approvaltest authenticator (D-WORLD-109), through d.Handler(). The judge defers the
// counter (W15); the writer judges it against the stored one, off the exact-match retry.

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

var kdPasskey = approvaltest.NewPasskey(approval.AlgES256, 72)

const kdOrigin, kdRP = "http://localhost:7644", "localhost"

// enrolPasskeyDecider enrols pk as name at packetNow (before the rig's clock) with counter stored.
func (r decideRig) enrolPasskeyDecider(name string, pk *approvaltest.Passkey, stored int64) {
	r.t.Helper()
	k := pk.Key()
	r.write(func(ctx context.Context, c *store.Credentials) error {
		return c.EnrolPasskey(ctx, store.ApproverRow{Name: name, Format: approval.FormatWebAuthn, PublicKey: k.Canonical(), EnrolledAt: packetNow},
			store.PasskeyRow{Name: name, RPID: kdRP, Origin: kdOrigin, CredentialID: pk.CredentialID(), AAGUID: make([]byte, 16), Alg: k.Alg, SignCount: stored})
	})
}

// decidePasskeyBody is a POST body deciding verb on (p, q, expiry) with pk's assertion of the
// decision message hash with counter count (and edits e).
func decidePasskeyBody(t *testing.T, pk *approvaltest.Passkey, verb string, p store.Object, q hashref.HashRef, expiry int64, count uint32, e ...approvaltest.Edit) []byte {
	t.Helper()
	v, _ := approval.DecisionVerb(verb)
	h, err := approval.DecisionMessageHash(v, [32]byte(mustDigest(t, p.Hash)), [32]byte(mustDigest(t, q)), expiry)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(map[string]any{"v": 1, "verb": verb, "packet": p.Hash.Digest(), "request": q.Digest(), "expiry": expiry,
		"webauthn": json.RawMessage(pk.Assert(kdOrigin, kdRP, h, count, e...))})
	return b
}

// signCount is name's stored passkey counter.
func (r decideRig) signCount(name string) int64 {
	r.t.Helper()
	ro, err := store.OpenCredentialsReadOnly(r.d.credentialsPath)
	if err != nil {
		r.t.Fatal(err)
	}
	defer func() { _ = ro.Close() }()
	p, ok, err := ro.Passkey(boundedTestContext(r.t), name)
	if err != nil || !ok {
		r.t.Fatal(name, ok, err)
	}
	return p.SignCount
}

// failWorldWrites makes every approval decision write fail until the returned func drops the trigger.
func (r decideRig) failWorldWrites() func() {
	r.t.Helper()
	db, err := sql.Open("sqlite", r.d.cfg.DBPath)
	if err != nil {
		r.t.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TRIGGER fail_decision BEFORE INSERT ON objects WHEN NEW.semantic_id = '` + broker.ApprovalDecisionV1 + `' BEGIN SELECT RAISE(ABORT, 'injected world-write failure'); END`); err != nil {
		r.t.Fatal(err)
	}
	return func() {
		r.t.Helper()
		defer func() { _ = db.Close() }()
		if _, err := db.Exec(`DROP TRIGGER fail_decision`); err != nil {
			r.t.Fatal(err)
		}
	}
}

// RD10, RD11 (passkey halves), KD2 end to end: an enrolled passkey's assertion decides approve and
// deny as webauthn:<name> (derived; a body naming decidedBy is 400), the evidence row is landed and
// the counter moved; a webauthn member beside another format is 400; a mint assertion and a
// counter that did not advance are 422 with nothing decided, recorded or moved.
func TestDecidePOSTPasskey(t *testing.T) {
	r := newDecideRig(t)
	r.enrolPasskeyDecider("mark-mac", kdPasskey, 1)
	q, p := r.packet("passkey", broker.WorkspaceScope, "agent")
	body := decidePasskeyBody(t, kdPasskey, "approve", p, q, decNow+50, 2)
	for name, b := range map[string][]byte{
		"decidedBy":          []byte(strings.Replace(string(body), `"v":1`, `"v":1,"decidedBy":"webauthn:mallory"`, 1)),
		"webauthn and ssh":   []byte(strings.Replace(string(body), `"v":1`, `"v":1,"sshsig":"-----BEGIN SSH SIGNATURE-----"`, 1)),
		"webauthn and typed": []byte(strings.Replace(string(body), `"v":1`, `"v":1,"typed":true`, 1)),
	} {
		requireRefusal(t, name, r.postKey(q, b), 400, "BadRequest")
	}
	mint, _ := approval.MessageHash([32]byte(mustDigest(t, q)), [32]byte(mustDigest(t, p.Hash)), decNow+50)
	mintBody, _ := json.Marshal(map[string]any{"v": 1, "verb": "approve", "packet": p.Hash.Digest(), "request": q.Digest(), "expiry": decNow + 50,
		"webauthn": json.RawMessage(kdPasskey.Assert(kdOrigin, kdRP, mint, 2))})
	if got := requireRefusal(t, "a mint assertion", r.postKey(q, mintBody), 422, "NotApproved"); !strings.Contains(got.Reason, "passkey assertion does not decide") {
		t.Errorf("a mint assertion: %q", got.Reason)
	}
	if n, _ := r.decisions(); n != 0 || r.signCount("mark-mac") != 1 {
		t.Fatalf("refused: %d decisions, counter %d", n, r.signCount("mark-mac"))
	}
	rec := r.postKey(q, body)
	var made surfacepage.DecisionMade
	if err := json.Unmarshal(rec.Body.Bytes(), &made); rec.Code != 200 || err != nil || made.DecidedBy != "webauthn:mark-mac" || made.Decided != "approve" {
		t.Fatalf("approve: %d %s", rec.Code, rec.Body.String())
	}
	row, ok := r.keyRow(q)
	if !ok || row.DecisionRef != made.Decision || row.Format != "webauthn" || row.Approver != "webauthn:mark-mac" || r.signCount("mark-mac") != 2 {
		t.Fatalf("evidence %+v %v, counter %d", row, ok, r.signCount("mark-mac"))
	}
	q2, p2 := r.packet("passkey-deny", broker.WorkspaceScope, "agent-2")
	if rec := r.postKey(q2, decidePasskeyBody(t, kdPasskey, "deny", p2, q2, decNow+50, 3)); rec.Code != 200 || !strings.Contains(rec.Body.String(), `"decided":"deny"`) {
		t.Fatalf("deny: %d %s", rec.Code, rec.Body.String())
	}
	q3, p3 := r.packet("passkey-regress", broker.WorkspaceScope, "agent-3")
	n, head := r.decisions()
	if got := requireRefusal(t, "counter 3 over stored 3", r.postKey(q3, decidePasskeyBody(t, kdPasskey, "approve", p3, q3, decNow+50, 3)), 422, "NotApproved"); !strings.Contains(got.Reason, "sign counter did not advance") {
		t.Errorf("counter: %q", got.Reason)
	}
	r.requireNothingDecidedSince(n, head)
	if _, ok := r.keyRow(q3); ok || r.signCount("mark-mac") != 3 {
		t.Fatalf("a refused counter recorded evidence or moved the counter (%d)", r.signCount("mark-mac"))
	}
}

// KD8 and RD13 (c): an authenticator without a counter (stored 0, assertion 0) decides once: the
// identical POST again is D6's 409 with one decision and one row; after a failed world write the
// identical POST lands that one decision (K0's exact match) and the next is 409.
func TestDecidePOSTCounterZeroAssertionLandsOnce(t *testing.T) {
	r := newDecideRig(t)
	r.enrolPasskeyDecider("mark-mac", kdPasskey, 0)
	q, p := r.packet("zero", broker.WorkspaceScope, "agent")
	body := decidePasskeyBody(t, kdPasskey, "approve", p, q, decNow+50, 0)
	if rec := r.postKey(q, body); rec.Code != 200 {
		t.Fatalf("first: %d %s", rec.Code, rec.Body.String())
	}
	requireRefusal(t, "the identical assertion again", r.postKey(q, body), 409, "Conflict")
	if n, _ := r.decisions(); n != 1 {
		t.Fatalf("%d decisions", n)
	}
	q2, p2 := r.packet("zero-crash", broker.WorkspaceScope, "agent-2")
	body2 := decidePasskeyBody(t, kdPasskey, "deny", p2, q2, decNow+50, 0)
	release := r.failWorldWrites()
	requireRefusal(t, "failed world write", r.postKey(q2, body2), 500, "Internal")
	if row, ok := r.keyRow(q2); !ok || row.DecisionRef != "" {
		t.Fatalf("after the failure: row %+v %v", row, ok)
	}
	release()
	if rec := r.postKey(q2, body2); rec.Code != 200 {
		t.Fatalf("retry: %d %s", rec.Code, rec.Body.String())
	}
	requireRefusal(t, "after the retry", r.postKey(q2, body2), 409, "Conflict")
	if n, _ := r.decisions(); n != 2 || r.signCount("mark-mac") != 0 {
		t.Fatalf("%d decisions, counter %d", n, r.signCount("mark-mac"))
	}
}

// RD13 (b), KD7 end to end: a counter-7 assertion (stored 6) whose world write fails → 500, one
// unlanded row, the counter 7, 0 decisions; a second, independently made assertion that also
// carries 7 → 422 (the writer's W15), nothing changed; the identical body → 200, the row landed,
// the counter still 7. Then, for a second request, a fresh assertion at N+1 lands by replacement.
func TestDecidePOSTPasskeyRetryAfterAFailedWorldWrite(t *testing.T) {
	r := newDecideRig(t)
	r.enrolPasskeyDecider("mark-mac", kdPasskey, 6)
	q, p := r.packet("pk-retry", broker.WorkspaceScope, "agent")
	body := decidePasskeyBody(t, kdPasskey, "approve", p, q, decNow+50, 7)
	release := r.failWorldWrites()
	requireRefusal(t, "failed world write", r.postKey(q, body), 500, "Internal")
	first, ok := r.keyRow(q)
	if n, _ := r.decisions(); n != 0 || !ok || first.DecisionRef != "" || r.signCount("mark-mac") != 7 {
		t.Fatalf("after the failure: %d decisions, row %+v %v, counter %d", n, first, ok, r.signCount("mark-mac"))
	}
	release()
	other := decidePasskeyBody(t, kdPasskey, "approve", p, q, decNow+50, 7, approvaltest.Edit{UserHandle: []byte("mark-mac")})
	requireRefusal(t, "another assertion at 7", r.postKey(q, other), 422, "NotApproved")
	if row, _ := r.keyRow(q); string(row.Artifact) != string(first.Artifact) {
		t.Fatal("a refused assertion replaced the row")
	}
	rec := r.postKey(q, body)
	row, _ := r.keyRow(q)
	if n, _ := r.decisions(); rec.Code != 200 || n != 1 || row.DecisionRef == "" || row.DecidedAt != first.DecidedAt || r.signCount("mark-mac") != 7 {
		t.Fatalf("retry: %d %s; %d decisions; row %+v; counter %d", rec.Code, rec.Body.String(), n, row, r.signCount("mark-mac"))
	}
	q2, p2 := r.packet("pk-fresh", broker.WorkspaceScope, "agent-2")
	release = r.failWorldWrites()
	requireRefusal(t, "failed world write 2", r.postKey(q2, decidePasskeyBody(t, kdPasskey, "deny", p2, q2, decNow+50, 8)), 500, "Internal")
	release()
	if rec := r.postKey(q2, decidePasskeyBody(t, kdPasskey, "deny", p2, q2, decNow+50, 9)); rec.Code != 200 || r.signCount("mark-mac") != 9 {
		t.Fatalf("a fresh assertion: %d %s, counter %d", rec.Code, rec.Body.String(), r.signCount("mark-mac"))
	}
	if row, _ := r.keyRow(q2); row.DecisionRef == "" {
		t.Fatalf("a fresh assertion: row %+v", row)
	}
}

// RD14: the counter-7 retry across a daemon restart: stored 7 in SQLite after the failed world
// write; the restarted daemon refuses another assertion carrying 7 and lands the identical one
// through the full HTTP/judge path without moving the counter.
func TestDecidePOSTRetryAfterAFailedWorldWriteAcrossRestart(t *testing.T) {
	r := newDecideRig(t)
	r.enrolPasskeyDecider("mark-mac", kdPasskey, 6)
	q, p := r.packet("pk-restart", broker.WorkspaceScope, "agent")
	body := decidePasskeyBody(t, kdPasskey, "approve", p, q, decNow+50, 7)
	release := r.failWorldWrites()
	requireRefusal(t, "failed world write", r.postKey(q, body), 500, "Internal")
	release()
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
	if r2.signCount("mark-mac") != 7 {
		t.Fatalf("stored counter after the restart: %d", r2.signCount("mark-mac"))
	}
	requireRefusal(t, "another assertion at 7", r2.postKey(q, decidePasskeyBody(t, kdPasskey, "approve", p, q, decNow+50, 7, approvaltest.Edit{UserHandle: []byte("mark-mac")})), 422, "NotApproved")
	if rec := r2.postKey(q, body); rec.Code != 200 || r2.signCount("mark-mac") != 7 {
		t.Fatalf("the identical retry after the restart: %d %s, counter %d", rec.Code, rec.Body.String(), r2.signCount("mark-mac"))
	}
	if row, _ := r2.keyRow(q); row.DecisionRef == "" {
		t.Fatalf("not landed: %+v", row)
	}
}

// The writer's counter refusals answer 422 (both operands); an approver race is not one of them.
// ErrPasskeyChanged is unit-only: the judge finds the passkey by its credential ID first (plan §3.1).
func TestIsCounterRefusalOperands(t *testing.T) {
	for _, base := range []error{store.ErrSignCountRegression, store.ErrPasskeyChanged} {
		if err := fmt.Errorf("%w: %q", base, "mark-mac"); !isCounterRefusal(err) {
			t.Errorf("%v: not a counter refusal", err)
		}
	}
	if isCounterRefusal(store.ErrApproverRevoked) {
		t.Error("a revoked approver is a counter refusal")
	}
}
