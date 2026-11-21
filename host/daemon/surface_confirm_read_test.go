package daemon

// Row 243 V3 sprint 3 M9 (design §2.5, §6 CR1, CR2): the confirmation card's read routes. Every
// request is approval.CanonicalVerb of a request whose prompt is broker.FinishPromptText and whose
// scope is broker.FormatFinishApprovalScope; every expected card is surfacepage.ConfirmationCard of
// the same row (sprint 2's golden is the shape contract). Stores and keys are tests' (D-WORLD-109).

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// putVerb stores a pending verb request id for ep, created at created (expiring 600 s later), for a
// draft (or ready) pull request of head world/<ep>.
func (a approvalRig) putVerb(id, ep string, created int64, draft bool) store.VerbRequestRow {
	a.t.Helper()
	return a.putVerbBytes(id, ep, created, canonVerb(a.t, verbFixture(created, ep, draft)))
}

// verbFixture is a valid request: its prompt FinishPromptText's, its scope FormatFinishApprovalScope's.
func verbFixture(created int64, ep string, draft bool) approval.VerbRequest {
	title, body := "Fix the <bar> & "+ep, "## Summary\n"+ep+"\n"
	sum := func(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }
	s := broker.FinishApprovalScope{Repository: "sunholo-data/ailang-world", Remote: "origin", Head: "world/" + ep, SHA: strings.Repeat("0a", 20), Base: "dev",
		TitleSHA256: sum(title), BodySHA256: sum(body), Draft: draft, ExpiresAt: created + 900}
	r := approval.VerbRequest{Schema: approval.VerbRequestSchema, Verb: approval.VerbFinish, Argv: []string{"ailang-worldd", "episode", "finish", ep, "--preset", "p", "--confirm-on-surface"},
		Episode: ep, Prompt: broker.FinishPromptText(s, title, body, []string{"go.mod"}, ""), Scope: broker.FormatFinishApprovalScope(s), Title: title, Body: body, Protected: []string{"go.mod"}}
	return r
}

func canonVerb(t *testing.T, r approval.VerbRequest) []byte {
	b, err := approval.CanonicalVerb(r)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func (a approvalRig) putVerbBytes(id, ep string, created int64, b []byte) store.VerbRequestRow {
	a.t.Helper()
	d, c := sha256.Sum256(b), sha256.Sum256([]byte(id))
	row := store.VerbRequestRow{RequestID: id, Request: b, Digest: hex.EncodeToString(d[:]), Challenge: hex.EncodeToString(c[:]), EpisodeID: ep, CreatedAt: created, ExpiresAt: created + 600}
	a.write(func(ctx context.Context, c *store.Credentials) error { return c.PutVerbRequest(ctx, row, created) })
	return row
}

// readVerb is the stored row id.
func (a approvalRig) readVerb(id string) (row store.VerbRequestRow) {
	a.t.Helper()
	a.write(func(ctx context.Context, c *store.Credentials) error { row, _, _ = c.VerbRequest(ctx, id); return nil })
	return row
}

// The route table: exactly three GETs and the one POST (sprint 4) under /surface/confirm, each behind
// approvalRoute; other methods 405 with no CORS header; a foreign Host 421; other names 404.
func TestConfirmRoutesRegistered(t *testing.T) {
	src, err := os.ReadFile("daemon.go")
	if err != nil {
		t.Fatal(err)
	}
	body := strings.SplitN(strings.SplitN(string(src), "func (d *Daemon) Handler() http.Handler {", 2)[1], "\n}\n", 2)[0]
	var got []string
	for _, m := range regexp.MustCompile(`mux\.HandleFunc\("([A-Z]+ /surface/confirm[^"]*)", (d\.approvalRoute\()?`).FindAllStringSubmatch(body, -1) {
		if m[2] == "" {
			t.Errorf("%s is not behind approvalRoute", m[1])
		}
		got = append(got, m[1])
	}
	sort.Strings(got)
	if want := []string{"GET /surface/confirm.js", "GET /surface/confirm/{id}", "GET /surface/confirmations/{id}", "POST /surface/confirm/{id}"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("confirmation routes %v, want %v", got, want)
	}
	a := newApprovalRig(t)
	id := "vr-0123456789abcdef"
	for _, target := range []string{"/surface/confirm/" + id, "/surface/confirm.js", "/surface/confirmations/" + id} {
		requireRefusal(t, "foreign "+target, a.serve("GET", target, "rebind.example:7644", nil), 421, "MisdirectedRequest")
		for _, method := range []string{"POST", "PUT", "DELETE", "OPTIONS"} {
			if method == "POST" && target == "/surface/confirm/"+id {
				continue // the confirmation (surface_confirm_test.go)
			}
			rec := a.serve(method, target, apLoopback, nil, "Origin", "http://localhost:7644", "Access-Control-Request-Method", "POST")
			if rec.Code != 405 || rec.Header().Get("Access-Control-Allow-Origin") != "" {
				t.Errorf("%s %s = %d, want 405 with no CORS header", method, target, rec.Code)
			}
		}
	}
	for _, target := range []string{"/surface/confirm", "/surface/confirm/", "/surface/confirm/" + id + "/x", "/surface/confirmations", "/surface/confirmations/" + id + "/x", "/surface/confirm.jsx"} {
		if rec := a.serve("GET", target, apLoopback, nil); rec.Code != 404 {
			t.Errorf("GET %s = %d, want 404", target, rec.Code)
		}
	}
}

// PD2 for the confirmation card: the page and the script are host/surfacepage/confirm's bytes,
// under surfaceCSP and the surface headers; a bad id is 400.
func TestConfirmPageAssetsServed(t *testing.T) {
	a := newApprovalRig(t)
	for _, tc := range []struct{ target, file, contentType string }{
		{"/surface/confirm/vr-0123456789abcdef", "index.html", "text/html; charset=utf-8"},
		{"/surface/confirm/vr-0000000000000000", "index.html", "text/html; charset=utf-8"},
		{"/surface/confirm.js", "confirm.js", "text/javascript; charset=utf-8"},
	} {
		want, err := os.ReadFile(filepath.Join("..", "surfacepage", "confirm", tc.file))
		if err != nil {
			t.Fatal(err)
		}
		rec := a.serve("GET", tc.target, apLoopback, nil)
		if rec.Code != 200 || len(want) == 0 || !bytes.Equal(rec.Body.Bytes(), want) {
			t.Errorf("GET %s: %d, or not the bytes of host/surfacepage/confirm/%s", tc.target, rec.Code, tc.file)
		}
		for k, v := range map[string]string{"Content-Security-Policy": surfaceCSP, "Content-Type": tc.contentType,
			"Cache-Control": "no-store", "X-Content-Type-Options": "nosniff", "Referrer-Policy": "no-referrer"} {
			if got := rec.Header().Get(k); got != v {
				t.Errorf("GET %s: %s = %q, want %q", tc.target, k, got, v)
			}
		}
		if rec.Header().Get("Set-Cookie") != "" {
			t.Errorf("GET %s sets a cookie", tc.target)
		}
		requireRefusal(t, "foreign host", a.serve("GET", tc.target, "evil.example", nil), 421, "MisdirectedRequest")
	}
	for _, bad := range []string{"vr-0123456789ABCDEF", "vr-0123", "mr-0123456789abcdef", "sha256:" + strings.Repeat("ab", 32)} {
		requireRefusal(t, bad, a.serve("GET", "/surface/confirm/"+bad, apLoopback, nil), 400, "BadRequest")
	}
}

// CR1: GET /surface/confirmations/{id} is surfacepage.ConfirmationCard of the stored row, byte for
// byte, for a pending, an approved and a superseded row: the bytes and the raw challenge, never the
// message hash, never the artifact. 400 bad id, 404 unknown, 503 below schema 9, 500 no store.
func TestConfirmationCardServed(t *testing.T) {
	a := newApprovalRig(t)
	pk := approvaltest.NewPasskey(approval.AlgEdDSA, 61)
	a.enrolPasskey("a-mac", pk, "http://localhost:7644")
	row := a.putVerb("vr-00000000000000a1", "w243", apT0-30, true)
	other := a.putVerb("vr-00000000000000b1", "w244", apT0-20, false)
	artifact := []byte("ARTIFACT-NEVER-SERVED")
	decision := "sha256:" + strings.Repeat("d4", 32)
	passkeys := []surfacepage.CardPasskey{{Origin: "http://localhost:7644", RPID: "localhost", CredentialID: pk.CredentialID()}}
	check := func(what string, r store.VerbRequestRow, state string) {
		t.Helper()
		rec := a.serve("GET", "/surface/confirmations/"+r.RequestID, apLoopback, nil)
		want, err := json.Marshal(surfacepage.ConfirmationCard(surfacepage.ConfirmationInput{ID: r.RequestID, Request: r.Request, Digest: r.Digest, Challenge: r.Challenge,
			Created: r.CreatedAt, Expires: r.ExpiresAt, Now: apT0, Outcome: r.Outcome, ApprovedBy: r.ApprovedBy, Decision: r.DecisionRef}, passkeys))
		if err := approvalHeaders(rec); err != nil || rec.Code != 200 {
			t.Fatalf("%s: %d %v", what, rec.Code, err)
		}
		if err != nil || !bytes.Equal(bytes.TrimSuffix(rec.Body.Bytes(), []byte("\n")), want) {
			t.Fatalf("%s:\n got %s\nwant %s", what, rec.Body.Bytes(), want)
		}
		var c surfacepage.Confirmation
		if err := json.Unmarshal(rec.Body.Bytes(), &c); err != nil || c.State != state || c.Request != base64.RawURLEncoding.EncodeToString(r.Request) || c.Challenge != r.Challenge {
			t.Fatalf("%s: %+v %v, want state %s", what, c, err, state)
		}
		d, _ := hex.DecodeString(r.Digest)
		n, _ := hex.DecodeString(r.Challenge)
		h, err := approval.VerbMessageHash([32]byte(d), [32]byte(n), r.ExpiresAt)
		if err != nil {
			t.Fatal(err)
		}
		for _, leak := range []string{hex.EncodeToString(h[:]), base64.RawURLEncoding.EncodeToString(h[:]), base64.StdEncoding.EncodeToString(h[:]), string(artifact), base64.StdEncoding.EncodeToString(artifact), `"artifact"`, `"hash"`} {
			if strings.Contains(rec.Body.String(), leak) {
				t.Fatalf("%s: the card carries %q", what, leak)
			}
		}
	}
	check("pending", row, surfacepage.ConfirmationPending)
	a.write(func(ctx context.Context, c *store.Credentials) error {
		return c.ConsumeVerbRequest(ctx, other.RequestID, other.Digest, store.ApproverUse{Name: "a-mac", Format: approval.FormatWebAuthn, PublicKey: pk.Key().Canonical(),
			Passkey: &store.PasskeyUse{CredentialID: pk.CredentialID(), SignCount: 2}}, "webauthn:a-mac", artifact, decision, apT0-10)
	})
	if got := a.readVerb(other.RequestID); got.Outcome != store.VerbOutcomeApproved || got.DecisionRef != decision {
		t.Fatalf("approved row = %+v", got)
	}
	check("approved", a.readVerb(other.RequestID), surfacepage.ConfirmationApproved)
	a.putVerb("vr-00000000000000a2", "w243", apT0-5, true) // supersedes row
	check("superseded", a.readVerb(row.RequestID), surfacepage.ConfirmationSuperseded)
	requireRefusal(t, "bad id", a.serve("GET", "/surface/confirmations/VR-00000000000000A1", apLoopback, nil), 400, "BadRequest")
	requireRefusal(t, "mint id", a.serve("GET", "/surface/confirmations/mr-00000000000000a1", apLoopback, nil), 400, "BadRequest")
	requireRefusal(t, "unknown", a.serve("GET", "/surface/confirmations/vr-00000000000000f1", apLoopback, nil), 404, "NotFound")
	requireRefusal(t, "foreign host", a.serve("GET", "/surface/confirmations/"+row.RequestID, "evil.example", nil), 421, "MisdirectedRequest")
	setVersion(t, a.d.credentialsPath, `DROP TABLE verb_requests`, `PRAGMA user_version = 8`)
	if r := requireRefusal(t, "v8", a.serve("GET", "/surface/confirmations/"+row.RequestID, apLoopback, nil), 503, "Restart"); !strings.Contains(r.Reason, "restart `ailang-worldd serve`") {
		t.Fatalf("v8 reason %q", r.Reason)
	}
	if err := os.Remove(a.d.credentialsPath); err != nil {
		t.Fatal(err)
	}
	requireRefusal(t, "no store", a.serve("GET", "/surface/confirmations/"+row.RequestID, apLoopback, nil), 500, "Internal")
}

// CR2: requests.json's confirmations are the live pending verb requests, newest first, each
// summarised from its stored bytes; a row whose bytes do not hash to its digest, or do not decode,
// is listed refused with no summary; approved, superseded and expired rows are not listed; below
// schema 9 the list is empty with confirmationsUnavailable and the other lists intact.
func TestRequestsJSONConfirmations(t *testing.T) {
	a := newApprovalRig(t)
	mint := a.putRequest("mr-00000000000000b1", "w201b", apT0-30)
	a.putVerb("vr-00000000000000a0", "wold", apT0-700, true) // expired at apT0-100
	sup := a.putVerb("vr-00000000000000a1", "w1", apT0-60, true)
	ready := a.putVerb("vr-00000000000000a2", "w2", apT0-50, false)
	done := a.putVerb("vr-00000000000000a3", "w3", apT0-45, true)
	draft := a.putVerb("vr-00000000000000a4", "w1", apT0-40, true) // supersedes sup
	undecodable := a.putVerbBytes("vr-00000000000000a5", "w5", apT0-30, []byte(`{"schema":"world/verb-request/v1"}`))
	edited := a.putVerb("vr-00000000000000a6", "w6", apT0-20, true)
	setVerbColumn(t, a.d.credentialsPath, edited.RequestID, "request", append(edited.Request, ' '))
	start, noScope := verbFixture(apT0-15, "w7", true), verbFixture(apT0-12, "w8", true)
	start.Verb, noScope.Scope = "episode start", "github-pr:o/r#x" // each decodes; Validate, then the scope parse, refuses
	invalid, unscoped := a.putVerbBytes("vr-00000000000000a7", "w7", apT0-15, canonVerb(t, start)), a.putVerbBytes("vr-00000000000000a8", "w8", apT0-12, canonVerb(t, noScope))
	a.write(func(ctx context.Context, c *store.Credentials) error {
		if err := c.EnrolApprover(ctx, store.ApproverRow{Name: "m", Format: "ssh", PublicKey: []byte("k"), EnrolledAt: apT0 - 100}); err != nil {
			return err
		}
		return c.ConsumeVerbRequest(ctx, done.RequestID, done.Digest, store.ApproverUse{Name: "m", Format: "ssh", PublicKey: []byte("k")}, "ssh:m", []byte("sig"), "sha256:"+strings.Repeat("d4", 32), apT0-10)
	})
	get := func() (surfacepage.RequestsJSON, string) {
		rec := a.serve("GET", "/surface/requests.json", apLoopback, nil)
		var out surfacepage.RequestsJSON
		if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil || rec.Code != 200 {
			t.Fatalf("requests.json = %d %v", rec.Code, err)
		}
		return out, rec.Body.String()
	}
	out, body := get()
	yes, no := true, false
	want := []surfacepage.PendingConfirmation{
		{ID: unscoped.RequestID, Expires: apT0 + 588, Refused: "the stored request's scope is not a finish approval"},
		{ID: invalid.RequestID, Expires: apT0 + 585, Refused: `the stored request is not a valid verb request: approval: bad verb request: verb is "episode start", want "episode finish"`},
		{ID: edited.RequestID, Expires: apT0 + 580, Refused: "the stored request does not match its digest"},
		{ID: undecodable.RequestID, Expires: apT0 + 570, Refused: out.Confirmations[3].Refused},
		{ID: draft.RequestID, Episode: "w1", Repository: "sunholo-data/ailang-world", Head: "world/w1", SHA: strings.Repeat("0a", 20), Draft: &yes, Expires: apT0 + 560},
		{ID: ready.RequestID, Episode: "w2", Repository: "sunholo-data/ailang-world", Head: "world/w2", SHA: strings.Repeat("0a", 20), Draft: &no, Expires: apT0 + 550},
	}
	if !reflect.DeepEqual(out.Confirmations, want) || out.ConfirmationsUnavailable != "" || !strings.HasPrefix(want[3].Refused, "the stored request is not a valid verb request: ") {
		t.Fatalf("confirmations = %s", body)
	}
	if len(out.Pending) != 1 || out.Pending[0].ID != mint.RequestID || out.Unavailable != "" {
		t.Fatalf("pending = %+v", out.Pending)
	}
	for _, leak := range []string{sup.RequestID, done.RequestID, draft.Challenge, draft.Digest, `"request"`, `"challenge"`, "sig"} {
		if strings.Contains(body, leak) {
			t.Errorf("requests.json carries %q", leak)
		}
	}
	setVersion(t, a.d.credentialsPath, `DROP TABLE verb_requests`, `PRAGMA user_version = 8`)
	out, body = get()
	if len(out.Pending) != 1 || out.Unavailable != "" || !strings.Contains(body, `"confirmations":[],"confirmationsUnavailable":"the credentials store predates verb requests: restart`) {
		t.Fatalf("a schema-8 store: %s", body)
	}
	if err := os.Remove(a.d.credentialsPath); err != nil {
		t.Fatal(err)
	}
	if _, body = get(); !strings.Contains(body, `"confirmations":[],"confirmationsUnavailable":"the verb requests could not be read"`) {
		t.Fatalf("a missing store: %s", body)
	}
}

// setVerbColumn sets one column of verb request id straight in the credentials file (a hand edit).
func setVerbColumn(t *testing.T, path, id, column string, v any) {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	if res, err := db.Exec(`UPDATE verb_requests SET `+column+` = ? WHERE request_id = ?`, v, id); err != nil {
		t.Fatal(err)
	} else if n, _ := res.RowsAffected(); n != 1 {
		t.Fatalf("%d rows updated", n)
	}
}

func TestConfirmationLinksRequireVerifiedEpisodeAndRecordedFinish(t *testing.T) {
	cfg, _ := followFixture(t, 2)
	d := newFollowDaemon(t, cfg)
	d.approveNow = func() int64 { return apT0 }
	a := approvalRig{t: t, d: d}
	snapshot := followSnapshotFor(t, "confirm-links", []string{"repo0", "repo1"}, nil)
	injectFollow(d, snapshot)
	const episode = "wa1"
	const remote = "git@github.com:sunholo-data/ailang-world.git"
	request := verbFixture(apT0-30, episode, true)
	scope, err := broker.ParseFinishApprovalScope(request.Scope)
	if err != nil {
		t.Fatal(err)
	}
	scope.Remote = remote
	request.Scope = broker.FormatFinishApprovalScope(scope)
	request.Prompt = broker.FinishPromptText(scope, request.Title, request.Body, request.Protected, request.Committer)
	row := a.putVerbBytes("vr-0000000000000041", episode, apT0-30, canonVerb(t, request))
	readCard := func() surfacepage.Confirmation {
		t.Helper()
		res := a.serve("GET", "/surface/confirmations/"+row.RequestID, apLoopback, nil)
		var card surfacepage.Confirmation
		if res.Code != 200 || json.Unmarshal(res.Body.Bytes(), &card) != nil {
			t.Fatalf("confirmation: %d %s", res.Code, res.Body)
		}
		return card
	}
	if card := readCard(); card.EpisodeHref != "" || card.PRHref != "" {
		t.Fatalf("unverified episode has links: %+v", card)
	}
	start := overlayStart(episode, d.follow.config.Repos[0].Repo, time.Unix(apT0, 0))
	start.Branch, start.Remote.URL = scope.Head, remote
	mintSurfaceTaskStart(t, d, start)
	card := readCard()
	href, err := url.Parse(card.EpisodeHref)
	if err != nil || href.Path != "/surface" || href.Host != "" || card.PRHref != "" ||
		href.Query().Get("at") != "repo0:episode:"+episode ||
		href.Query().Get("brane") != "repo0:brane:row-189" {
		t.Fatalf("pending confirmation links: episode %q, pull request %q", card.EpisodeHref, card.PRHref)
	}

	pk := approvaltest.NewPasskey(approval.AlgEdDSA, 62)
	a.enrolPasskey("links-key", pk, apOrigin)
	decision := "sha256:" + strings.Repeat("da", 32)
	a.write(func(ctx context.Context, c *store.Credentials) error {
		return c.ConsumeVerbRequest(ctx, row.RequestID, row.Digest, store.ApproverUse{Name: "links-key", Format: approval.FormatWebAuthn,
			PublicKey: pk.Key().Canonical(), Passkey: &store.PasskeyUse{CredentialID: pk.CredentialID(), SignCount: 2}},
			"webauthn:links-key", []byte("signed artifact"), decision, apT0-10)
	})
	if card := readCard(); card.State != surfacepage.ConfirmationApproved || card.EpisodeHref == "" || card.PRHref != "" {
		t.Fatalf("approval alone should only link to the episode: %+v", card)
	}
	const pr = "https://github.com/sunholo-data/ailang-world/pull/631"
	finish := EpisodeFinish{Episode: episode, Branch: scope.Head, Remote: scope.Remote, Scope: request.Scope,
		Commit: scope.SHA, Approval: decision, Outcome: "done", PR: FinishPR{URL: pr, Number: 631, Draft: true}}
	for _, tc := range []struct {
		name, provenance string
		edit             func(*EpisodeFinish)
	}{
		{"wrong provenance", "agent:fixture", func(f *EpisodeFinish) { f.FinishedAt = "2026-10-21T08:00:00Z" }},
		{"other commit", finishProvenance, func(f *EpisodeFinish) { f.Commit = strings.Repeat("a", 40) }},
		{"other repository", finishProvenance, func(f *EpisodeFinish) { f.PR.URL = "https://github.com/other/ailang-world/pull/631" }},
		{"different approval", finishProvenance, func(f *EpisodeFinish) { f.Approval = "sha256:" + strings.Repeat("ff", 32) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			wrong := finish
			tc.edit(&wrong)
			putFinish(t, d, wrong, tc.provenance)
			if card := readCard(); card.EpisodeHref == "" || card.PRHref != "" {
				t.Fatalf("incorrect finish exposed a PR: %+v", card)
			}
		})
	}
	putFinish(t, d, finish, finishProvenance)
	if card := readCard(); card.PRHref != pr {
		t.Fatalf("verified successful finish PR = %q", card.PRHref)
	}
	var step surface.State
	for _, s := range snapshot.Envelope.Repos[0].Index.States {
		if s.Kind == surface.KindStep && s.Brane != "" {
			step = s
			break
		}
	}
	snapshot.Envelope.Open.Episodes = snapshot.Envelope.Open.Episodes[1:]
	snapshot.Envelope.Open.Aliases = []surface.EpisodeAlias{{ID: "repo0:episode:" + episode, Target: "repo0:commit:" + step.ID}}
	if err := surface.ValidateFollow(snapshot.Envelope); err != nil {
		t.Fatal(err)
	}
	card = readCard()
	href, err = url.Parse(card.EpisodeHref)
	if err != nil || href.Query().Get("at") != "repo0:episode:"+episode ||
		href.Query().Get("brane") != "repo0:brane:"+step.Brane || card.PRHref != pr {
		t.Fatalf("landed confirmation lost its qualified episode alias: episode %q, PR %q", card.EpisodeHref, card.PRHref)
	}
	scope.Repository = "another-owner/ailang-world"
	request.Scope = broker.FormatFinishApprovalScope(scope)
	request.Prompt = broker.FinishPromptText(scope, request.Title, request.Body, request.Protected, request.Committer)
	row = a.putVerbBytes("vr-0000000000000042", episode, apT0-2, canonVerb(t, request))
	if card := readCard(); card.EpisodeHref != "" || card.PRHref != "" {
		t.Fatalf("signed repository mismatch carried links: %+v", card)
	}
}
