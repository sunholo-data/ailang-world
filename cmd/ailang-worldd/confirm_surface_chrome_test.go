package main

// Row 243 V3 sprint 6 M16 (design §6 XC2, XC3): a finish confirmation given by a passkey in a live Chrome,
// on the A2/A3 driver (testdata/passkey_cdp.mjs, unchanged; its virtual authenticator holds approvaltest's
// key, enrolled by virtualPasskey), gated as XB1 (WORLD_CHROME, WORLD_REQUIRE_PASSKEY_BROWSER=1; skipped
// and counted without them). The card is S3/S4's /surface/confirm/<id> served by a real daemon.New on the
// loopback; the request is stored by the CLI's production put (XC2) or by the CLI itself (XC3).

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

// confirmView is the confirmation card after its passkey ceremony settles.
type confirmView struct {
	State, Status, Answer string
	Gets                  []cardGet
	Violations            []string
}

// confirmInChrome opens card, waits until it offers the passkey, clicks it and waits (each step bounded by
// the driver) until the card reads its request again, closed. It never calls t.Fatal (XC3 calls it from
// the CLI's wait).
func confirmInChrome(b *passkeyBrowser, card string) (confirmView, error) {
	var v confirmView
	if err := b.call(nil, "openCard", map[string]any{"url": card, "raw": true}); err != nil {
		return v, err
	}
	if err := b.call(nil, "waitFor", map[string]any{"expression": `document.body.dataset.state === 'ready' && !document.getElementById('go').hidden`}); err != nil {
		return v, err
	}
	if err := b.call(nil, "tap", map[string]any{"id": "go"}); err != nil {
		return v, err
	}
	err := b.call(&v, "waitFor", map[string]any{"ms": 30000, "expression": `document.body.dataset.state === 'closed' && { state: document.body.dataset.state,
		status: document.getElementById('status').textContent, answer: document.getElementById('answer').textContent, gets: window.__gets || [], violations: window.__cspv || [] }`})
	return v, err
}

// XC2: the virtual authenticator confirms a request stored by the CLI's own put on the card a real daemon
// serves; the daemon records a world decision decided by webauthn:mark-key and consumes the row naming it;
// the CLI's own re-check (a) and read-back (b) accept both. Control: a driver step is bounded.
func TestConfirmPageConfirmsWithAPasskeyInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	b.do(t, nil, "tab", map[string]any{"authenticator": "uv", "inject": "getRecorder"})
	if s := time.Now(); b.call(nil, "waitFor", map[string]any{"expression": "false", "ms": 300}) == nil || time.Since(s) > 10*time.Second {
		t.Fatalf("a waitFor that never holds did not fail within its bound (%s)", time.Since(s))
	}
	now := time.Now().Unix()
	creds, in := surfaceUnit(t)
	in.scope.ExpiresAt = now + 900 // a finish scope expires 900 s after the CLI builds it
	in.scopeText = broker.FormatFinishApprovalScope(in.scope)
	origin := cardOrigin(t, creds, "", nil)
	virtualPasskey(t, b, creds, origin, "mark-key", 93, now-100)
	own := storeOwn(t, in, "vr-00000000000000c7", 0xc7, now)
	v, err := confirmInChrome(b, origin+"/surface/confirm/"+own.id)
	if err != nil {
		t.Fatal(err)
	}
	ro, err := store.OpenCredentialsReadOnly(creds)
	if err != nil {
		t.Fatal(err)
	}
	row, ok, err := ro.VerbRequest(boundedTestContext(t), own.id)
	_ = ro.Close()
	if err != nil || !ok || row.Outcome != store.VerbOutcomeApproved || row.ApprovedBy != "webauthn:mark-key" || !strings.HasPrefix(row.DecisionRef, "sha256:") ||
		row.ConsumedAt < now || row.ConsumedAt > time.Now().Unix() {
		t.Fatalf("row %+v (%v %v)", row, ok, err)
	}
	h, _ := approval.VerbMessageHash(own.digest, own.challenge, own.expires)
	if v.Answer != "Confirmed by webauthn:mark-key: decision "+row.DecisionRef+"; the waiting episode finish opens the pull request." ||
		v.Status != "Already confirmed by webauthn:mark-key: decision "+row.DecisionRef+"; nothing to confirm here." || len(v.Violations) != 0 ||
		len(v.Gets) != 1 || v.Gets[0].Challenge != base64.RawURLEncoding.EncodeToString(h[:]) || v.Gets[0].UserVerification != "required" {
		t.Fatalf("the card shows %+v; want the confirmation of %s with challenge %s", v, row.DecisionRef, hex.EncodeToString(h[:]))
	}
	if refusal := readBackFinishDecision(worldReader{ctx: boundedTestContext(t), c: newClient(origin)}, row.DecisionRef, "webauthn:mark-key", in.scopeText); refusal != "" {
		t.Fatal(refusal)
	}
	if _, refusal, err := pollVerbRequest(creds, own, time.Now().Unix()); refusal != "" || err != nil {
		t.Fatalf("the CLI's re-check of the passkey's confirmation: %q %v", refusal, err)
	}
}

// XC3: episode finish --confirm-on-surface (the real CLI path on S5's rig, fake gh) stores its request and
// prints the card's link for the passkey's origin; the page confirms with the virtual authenticator while
// the CLI waits; the CLI re-checks, reads the decision back and prints the fake gh's PR URL. The terminal
// is never read; the rig's clock starts at the wall clock the daemon reads.
func TestFinishConfirmOnSurfaceInChrome(t *testing.T) {
	b := startPasskeyBrowser(t)
	b.do(t, nil, "tab", map[string]any{"authenticator": "uv"})
	r := newSurfaceRig(t)
	srv := httptest.NewServer(r.d.Handler())
	t.Cleanup(srv.Close)
	origin := strings.Replace(srv.URL, "127.0.0.1", "localhost", 1)
	virtualPasskey(t, b, r.creds, origin, "mark-key", 94, r.now-100)
	id, card := drawnID(1), ""
	var v confirmView
	var chromeErr error
	r.onWait = func(n int64) {
		if n == 1 {
			card = origin + "/surface/confirm/" + id
			v, chromeErr = confirmInChrome(b, card)
		}
	}
	rc, out, errs := r.run()
	if chromeErr != nil || rc != exitOK {
		t.Fatalf("chrome: %v; rc %d\nstdout %s\nstderr %s", chromeErr, rc, out, errs)
	}
	row := r.row(id)
	pr := r.f.finishRecords[0].PR.URL // the fake gh's own answer
	if !strings.Contains(out, "\n  open "+card+"\n") || !strings.Contains(out, fmt.Sprintf("\nconfirmed by webauthn:mark-key (decision %s)\n", row.DecisionRef)) ||
		!strings.HasSuffix(out, "\n"+pr+"\n") || pr == "" || !strings.HasPrefix(v.Answer, "Confirmed by webauthn:mark-key: decision "+row.DecisionRef) {
		t.Fatalf("stdout\n%s\ncard %+v; want the link %s, the confirmation and the PR %q last", out, v, card, pr)
	}
	if rec := r.f.finishRecords; len(rec) != 1 || rec[0].Approval != row.DecisionRef || rec[0].Outcome != "done" || row.ApprovedBy != "webauthn:mark-key" ||
		row.CreatedAt != r.start || r.tty.reads != 0 || len(r.f.scopes) != 0 {
		t.Fatalf("records %+v, row %+v, terminal reads %d, typed approve phases %d", rec, row, r.tty.reads, len(r.f.scopes))
	}
	var decision struct{ Decision, DecidedBy string }
	r.object(row.DecisionRef, &decision)
	if decision.Decision != "approve" || decision.DecidedBy != "webauthn:mark-key" {
		t.Fatalf("decision %+v", decision)
	}
}
