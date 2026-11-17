package daemon

// Row 243 V3 sprint 4 M10 (design §2.6, §6 CX14-CX21; S3 audit note 3's daemon half): the
// confirmation's writes. C14 records the world store's finish approval decided by the key; C15
// consumes the row and names it in one credentials transaction; the writer and humanMu are held
// across both. Each crash window of the design's table is driven by its instrument: an objects
// trigger refusing ApprovalDecisionV1 (A4's RD13 (a)), a verb_requests UPDATE trigger, and the C15
// seam for a supersession or an approver change between C14 and C15.

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// worldDecision is the finish approval decision ref names: its decision and decidedBy, and its
// request's effect, requester and scope, read back from the world store.
func (a approvalRig) worldDecision(ref string) (verb, decidedBy, effect, requester, scope string) {
	a.t.Helper()
	ctx := boundedTestContext(a.t)
	obj, ok, err := a.d.store.GetObject(ctx, hashref.MustParse(ref))
	var wire struct {
		RequestRef string `json:"requestRef"`
		Decision   string `json:"decision"`
		DecidedBy  string `json:"decidedBy"`
	}
	if err != nil || !ok || obj.SemanticID != broker.ApprovalDecisionV1 || json.Unmarshal(obj.Payload, &wire) != nil {
		a.t.Fatalf("decision %s: %v %v %s", ref, ok, err, obj.SemanticID)
	}
	q := hashref.MustParse(wire.RequestRef)
	reqObj, ok, err := a.d.store.GetObject(ctx, q)
	req, derr := broker.DecodeApprovalRequest(reqObj.Payload)
	if err != nil || !ok || derr != nil {
		a.t.Fatalf("request %s: %v %v %v", q, ok, err, derr)
	}
	if onChain, decided, chainVerb, chainBy, err := broker.ApprovalChainState(ctx, a.d.store, q); err != nil || !onChain || !decided || chainVerb != wire.Decision || chainBy != wire.DecidedBy {
		a.t.Fatalf("chain of %s: %v %v %q %q %v", q, onChain, decided, chainVerb, chainBy, err)
	}
	return wire.Decision, wire.DecidedBy, req.Effect, req.Requester, req.Scope
}

// CX17: an SSHSIG and a passkey assertion each confirm (the passkey one second before expiry): the
// answer is surfacepage.Confirmed of the row's decision; the world decision approves the request's
// own scope for operator:episode-finish, decided by the key; the row is consumed with five columns;
// the approver is used and the passkey's counter moved; a replay is 409 and writes nothing.
func TestConfirmPOSTConfirms(t *testing.T) {
	a, row := newConfirmRig(t)
	pk := a.putVerb("vr-00000000000000d1", "w-d1", apT0-30, false)
	m, _ := verbMessageOf(t, row)
	armored := apSigner.Sign(approval.VerbNamespace, m)
	_, h := verbMessageOf(t, pk)
	assertion := apPasskey.Assert(apOrigin, "localhost", h, 2)
	for i, tc := range []struct {
		row       store.VerbRequestRow
		body      []byte
		artifact  []byte
		now       int64
		by        string
		ep, draft string
	}{
		{row, sshBody(row.Digest, armored), armored, apT0, "ssh:mark-ssh", cfEp, "draft"},
		{pk, approveBody(pk.Digest, assertion), assertion, pk.ExpiresAt - 1, "webauthn:mark", "w-d1", "ready"},
	} {
		a.d.approveNow = func() int64 { return tc.now }
		rec := a.confirm(tc.row.RequestID, tc.body)
		got := a.readVerb(tc.row.RequestID)
		want, _ := json.Marshal(surfacepage.Confirmed(tc.row.RequestID, tc.by, got.DecisionRef))
		if rec.Code != 200 || !bytes.Equal(bytes.TrimSuffix(rec.Body.Bytes(), []byte("\n")), want) || approvalHeaders(rec) != nil {
			t.Fatalf("%s: %d %s, want %s", tc.by, rec.Code, rec.Body.String(), want)
		}
		if got.Outcome != store.VerbOutcomeApproved || got.ApprovedBy != tc.by || !bytes.Equal(got.Artifact, tc.artifact) || got.ConsumedAt != tc.now || !strings.HasPrefix(got.DecisionRef, "sha256:") {
			t.Fatalf("%s: row %+v", tc.by, got)
		}
		verb, decidedBy, effect, requester, scope := a.worldDecision(got.DecisionRef)
		if wantScope := verbFixture(apT0-30, tc.ep, tc.draft == "draft").Scope; verb != "approve" || decidedBy != tc.by || effect != broker.EffectHumanApprove || requester != "operator:episode-finish" || scope != wantScope {
			t.Fatalf("%s: world decision %q %q %q %q %q", tc.by, verb, decidedBy, effect, requester, scope)
		}
		if n, _ := (decideRig{a}).decisions(); n != i+1 {
			t.Fatalf("%s: %d decisions", tc.by, n)
		}
		before := a.confirmState()
		if r := requireRefusal(t, tc.by+" replay", a.confirm(tc.row.RequestID, tc.body), 409, "Conflict"); !strings.Contains(r.Reason, "already confirmed by "+tc.by) || a.confirmState() != before {
			t.Fatalf("%s replay: %q", tc.by, r.Reason)
		}
	}
	if s := a.credState(); !strings.Contains(s, "[mark-ssh 1800000000 0") || !strings.Contains(s, "[mark 2 ") {
		t.Fatalf("approver use or counter not moved: %s", s)
	}
}

// CX19: the confirmation holds humanMu, the mutex the typed finish approve phase takes
// (host/daemon/episode_finish.go), from its clock read through C11 to C15 (each probe's TryLock fails
// while the handler holds it); eight concurrent POSTs of one body confirm once, never a busy writer.
func TestConfirmPOSTHoldsHumanMu(t *testing.T) {
	a, row := newConfirmRig(t)
	body := confirmSSH(t, apSigner, approval.VerbNamespace, row)
	var mu sync.Mutex
	var probes []string
	probe := func(at string) {
		held := !a.d.humanMu.TryLock()
		if !held {
			a.d.humanMu.Unlock()
		}
		mu.Lock()
		probes = append(probes, fmt.Sprintf("%s:%v", at, held))
		mu.Unlock()
	}
	a.d.approveNow = func() int64 { probe("clock"); return apT0 }
	judgeVerbSSH = func(ctx context.Context, ro *store.Credentials, in authority.VerbJudgeInput, armored []byte) (*approval.SignedVerb, error) {
		probe("C11")
		return authority.JudgeVerbSSH(ctx, ro, in, armored)
	}
	consumeVerbRequest = func(c *store.Credentials, ctx context.Context, id, digest string, use store.ApproverUse, by string, artifact []byte, ref string, now int64) error {
		probe("C15")
		return c.ConsumeVerbRequest(ctx, id, digest, use, by, artifact, ref, now)
	}
	t.Cleanup(func() {
		judgeVerbSSH, consumeVerbRequest = authority.JudgeVerbSSH, (*store.Credentials).ConsumeVerbRequest
	})
	if rec := a.confirm(cfID, body); rec.Code != 200 || fmt.Sprint(probes) != "[clock:true C11:true C15:true]" {
		t.Fatalf("%d %s; probes %v", rec.Code, rec.Body.String(), probes)
	}
	next := a.putVerb("vr-00000000000000c6", "w-c6", apT0-30, true)
	body = confirmSSH(t, apSigner, approval.VerbNamespace, next)
	codes := make([]int, 8)
	var wg sync.WaitGroup
	for i := range codes {
		wg.Add(1)
		go func(i int) { defer wg.Done(); codes[i] = a.confirm(next.RequestID, body).Code }(i)
	}
	wg.Wait()
	ok := 0
	for _, c := range codes {
		if c == 200 {
			ok++
		} else if c != 409 {
			t.Errorf("a concurrent POST answered %d (codes %v)", c, codes)
		}
	}
	if n, _ := (decideRig{a}).decisions(); ok != 1 || n != 2 {
		t.Fatalf("%d confirmed, %d decisions (codes %v)", ok, n, codes)
	}
}

// CX20 (window "C14 fails"): the world write fails (A4's objects trigger) → 500, no decision, the
// row pending and the counter unmoved; the identical passkey body after the trigger is dropped → 200,
// one decision, the row names it; a replay → 409.
func TestConfirmPOSTRetriesAFailedWorldWrite(t *testing.T) {
	a, row := newConfirmRig(t)
	body := confirmPasskey(t, apPasskey, row, 2)
	cred := a.credState()
	drop := (decideRig{a}).failWorldWrites()
	if r := requireRefusal(t, "C14 fails", a.confirm(cfID, body), 500, "Internal"); r.Reason != confirmIncomplete {
		t.Fatalf("reason %q", r.Reason)
	}
	if n, _ := (decideRig{a}).decisions(); n != 0 || a.credState() != cred {
		t.Fatalf("after the failure: %d decisions; credentials %s (was %s)", n, a.credState(), cred)
	}
	drop()
	rec := a.confirm(cfID, body)
	var c surfacepage.ConfirmedJSON
	_ = json.Unmarshal(rec.Body.Bytes(), &c)
	if n, _ := (decideRig{a}).decisions(); rec.Code != 200 || n != 1 || a.readVerb(cfID).DecisionRef != c.Decision {
		t.Fatalf("retry: %d %s; %d decisions; row %+v", rec.Code, rec.Body.String(), n, a.readVerb(cfID))
	}
	requireRefusal(t, "replay", a.confirm(cfID, body), 409, "Conflict")
}

// CX21 (window "after C14, C15 fails"): a verb_requests UPDATE trigger fails the consume after the
// world write → 500, one decision that no row names, the row pending, the counter unmoved. The
// identical body after the trigger is dropped: in the same second the identical decision answers its
// existing ref (W59: still one decision); a second later a second decision lands and the row names
// it, the first left inert.
func TestConfirmPOSTRetriesAFailedConsume(t *testing.T) {
	a, row := newConfirmRig(t)
	later := a.putVerb("vr-00000000000000c7", "w-c7", apT0-30, true)
	db, err := sql.Open("sqlite", a.d.credentialsPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	trigger := func(sql string) {
		if _, err := db.Exec(sql); err != nil {
			t.Fatal(err)
		}
	}
	for i, tc := range []struct {
		row   store.VerbRequestRow
		retry int64
		n     int
	}{{row, apT0, 1}, {later, apT0 + 1, 3}} {
		body := confirmPasskey(t, apPasskey, tc.row, uint32(2+i))
		a.d.approveNow = func() int64 { return apT0 }
		trigger(`CREATE TRIGGER fail_consume BEFORE UPDATE ON verb_requests WHEN NEW.outcome = 'approved' BEGIN SELECT RAISE(ABORT, 'injected consume failure'); END`)
		cred := a.credState()
		n0, _ := (decideRig{a}).decisions()
		if r := requireRefusal(t, "C15 fails", a.confirm(tc.row.RequestID, body), 500, "Internal"); r.Reason != confirmIncomplete {
			t.Fatalf("reason %q", r.Reason)
		}
		if n, _ := (decideRig{a}).decisions(); n != n0+1 || a.credState() != cred {
			t.Fatalf("after the failure: %d decisions (was %d); credentials %s (was %s)", n, n0, a.credState(), cred)
		}
		trigger(`DROP TRIGGER fail_consume`)
		a.d.approveNow = func() int64 { return tc.retry }
		rec := a.confirm(tc.row.RequestID, body)
		var c surfacepage.ConfirmedJSON
		_ = json.Unmarshal(rec.Body.Bytes(), &c)
		if n, _ := (decideRig{a}).decisions(); rec.Code != 200 || n != tc.n || a.readVerb(tc.row.RequestID).DecisionRef != c.Decision {
			t.Fatalf("retry at %d: %d %s; %d decisions, want %d", tc.retry, rec.Code, rec.Body.String(), n, tc.n)
		}
	}
}

// Audit note 3 and C15's refusals (window "after C14, C15 refuses"), through the C15 seam: the
// request superseded by a newer finish of its episode between C14 and C15 → 409 naming it; a
// concurrent consume that won → 409; the approver revoked, or the passkey's counter moved past the
// assertion's → 422; any other store error → 500. Each leaves one decision no row names (inert, W62)
// and, but for the won consume, the row unconsumed and the key unused.
func TestConfirmPOSTRefusedAtC15(t *testing.T) {
	a, _ := newConfirmRig(t)
	t.Cleanup(func() { consumeVerbRequest = (*store.Credentials).ConsumeVerbRequest })
	for i, tc := range []struct {
		name          string
		before        func(ctx context.Context, c *store.Credentials, row store.VerbRequestRow) error
		passkey       bool
		status        int
		class, reason string
		outcome       string
	}{
		{"superseded", func(ctx context.Context, c *store.Credentials, row store.VerbRequestRow) error {
			n := row
			n.RequestID = "vr-00000000000000e9" // the same episode
			return c.PutVerbRequest(ctx, n, apT0)
		}, false, 409, "Conflict", "superseded by a newer finish of w-e1", store.VerbOutcomeSuperseded},
		{"consume lost", nil, false, 409, "Conflict", "changed while it was confirmed", store.VerbOutcomeApproved},
		{"store error", func(context.Context, *store.Credentials, store.VerbRequestRow) error { return errors.New("injected") }, false, 500, "Internal", confirmIncomplete, ""},
		{"counter", func(context.Context, *store.Credentials, store.VerbRequestRow) error {
			db, err := sql.Open("sqlite", a.d.credentialsPath)
			if err == nil {
				_, err = db.Exec(`UPDATE webauthn_credentials SET sign_count = 9 WHERE name = 'mark'`)
				_ = db.Close()
			}
			return err
		}, true, 422, "NotApproved", "counter", ""},
		{"approver revoked", func(ctx context.Context, c *store.Credentials, _ store.VerbRequestRow) error {
			return c.RevokeApprover(ctx, "mark-ssh", apT0)
		}, false, 422, "NotApproved", "revoked", ""},
	} {
		row := a.putVerb(fmt.Sprintf("vr-00000000000000e%d", i+1), fmt.Sprintf("w-e%d", i+1), apT0-30, true)
		body := confirmSSH(t, apSigner, approval.VerbNamespace, row)
		if tc.passkey {
			body = confirmPasskey(t, apPasskey, row, 2)
		}
		consumeVerbRequest = func(c *store.Credentials, ctx context.Context, id, digest string, use store.ApproverUse, by string, artifact []byte, ref string, now int64) error {
			if tc.before == nil { // another consume of the same row wins first
				if err := c.ConsumeVerbRequest(ctx, id, digest, use, by, artifact, ref, now); err != nil {
					return err
				}
			} else if err := tc.before(ctx, c, row); err != nil {
				return err
			}
			return c.ConsumeVerbRequest(ctx, id, digest, use, by, artifact, ref, now)
		}
		n0, _ := (decideRig{a}).decisions()
		if r := requireRefusal(t, tc.name, a.confirm(row.RequestID, body), tc.status, tc.class); !strings.Contains(r.Reason, tc.reason) {
			t.Errorf("%s: reason %q", tc.name, r.Reason)
		}
		got := a.readVerb(row.RequestID)
		if n, _ := (decideRig{a}).decisions(); n != n0+1 || got.Outcome != tc.outcome || (tc.outcome != store.VerbOutcomeApproved && got.DecisionRef != "") {
			t.Errorf("%s: %d decisions (was %d), row %+v", tc.name, n, n0, got)
		}
	}
	if s := a.credState(); !strings.Contains(s, "[mark 9 ") || strings.Count(s, " approved ") != 1 { // 9: the counter row's own edit
		t.Errorf("a refused C15 used a key or consumed a row: %s", s)
	}
	consumeVerbRequest = (*store.Credentials).ConsumeVerbRequest
}

// C15's sentinel → status mapping, operand by operand, as ConsumeVerbRequest wraps them.
func TestConsumeRefusalOperands(t *testing.T) {
	a, row := newConfirmRig(t)
	for _, tc := range []struct {
		err    error
		status int
	}{
		{store.ErrVerbRequestSuperseded, 409}, {store.ErrVerbRequestConsumed, 409}, {store.ErrVerbRequestUnknown, 409},
		{store.ErrVerbRequestExpired, 409}, {store.ErrVerbRequestClockBehind, 409}, {store.ErrVerbRequestDigestChanged, 409},
		{store.ErrApproverRevoked, 422}, {store.ErrApproverKeyChanged, 422}, {store.ErrApproverUnknown, 422}, {store.ErrApproverClockBehind, 422},
		{store.ErrSignCountRegression, 422}, {store.ErrPasskeyChanged, 422},
		{&store.OlderCredentialsStoreError{Version: 8, Need: 9}, 503}, {store.ErrBadVerbRequest, 500}, {errors.New("disk I/O error"), 500},
	} {
		rec := httptest.NewRecorder()
		a.d.consumeRefusal(rec, httptest.NewRequest("POST", "/surface/confirm/"+cfID, nil), row, fmt.Errorf("%w: %s", tc.err, cfID))
		if rec.Code != tc.status {
			t.Errorf("%v → %d, want %d", tc.err, rec.Code, tc.status)
		}
	}
}
