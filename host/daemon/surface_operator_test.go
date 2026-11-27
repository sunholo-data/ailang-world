package daemon

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

const opTask = "# Task\n<script>alert(1)</script> ; $(literal) &"
const opOrigin = "http://localhost:7644"

type operatorLog struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (l *operatorLog) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.Write(p)
}

func (l *operatorLog) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.String()
}

type operatorHTTP struct {
	approvalRig
	handler http.Handler
	server  *httptest.Server
	now     atomic.Int64
	h       store.OperatorHeartbeat
	fake    string
	log     operatorLog
}

func newOperatorHTTP(t *testing.T, budget ...time.Duration) *operatorHTTP {
	t.Helper()
	f := newOperatorSurface(t, nil, budget...)
	f.server = httptest.NewServer(f.handler)
	t.Cleanup(f.server.Close)
	return f
}

func newOperatorSurface(t *testing.T, edit func(*store.OperatorStartReview), budget ...time.Duration) *operatorHTTP {
	t.Helper()
	f := &operatorHTTP{approvalRig: newApprovalRig(t)}
	if len(budget) != 0 {
		f.approvalRig.writeBudget = budget[0]
		f.d.commitBudget, f.d.readDeadline = budget[0], budget[0]
	}
	f.now.Store(apT0)
	f.d.approveNow = f.now.Load
	f.d.errLog = &f.log
	f.d.cfg.SurfaceOrigins = []string{opOrigin}
	root := t.TempDir()
	f.fake = filepath.Join(root, "ran")
	if err := os.WriteFile(filepath.Join(root, "fake-codex"), []byte("#!/bin/sh\nprintf x > \"${0%/*}/ran\"\n"), 0700); err != nil {
		t.Fatal(err)
	}
	physical, err := store.CanonicalPath(f.d.credentialsPath)
	if err != nil {
		t.Fatal(err)
	}
	build := CurrentBuild().ExecutableSHA256
	entry := store.OperatorStartReview{Store: physical, Workspace: root, Sessions: root, Binary: filepath.Join(root, "fake-codex"), Wall: "1m",
		Start: store.EpisodeStart{Schema: "world/episode-start/v1", Agent: "codex", Model: "fixture", TTLS: 60,
			Preset: store.StartPreset{Name: "sol", SHA256: "sha256:fixture"}, Overrides: []string{}}}
	if edit != nil {
		edit(&entry)
	}
	f.h = store.OperatorHeartbeat{Owner: store.OperatorOwner{Nonce: "fixture"}, Protocol: store.OperatorRunnerProtocol, Build: build,
		Readiness: "start", SeenAt: apT0, Catalogue: opJSON(map[string]store.OperatorStartReview{"sol": entry})}
	f.publish()
	f.handler = f.d.Handler()
	return f
}

func opJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func (f *operatorHTTP) publish() {
	f.write(func(ctx context.Context, c *store.Credentials) error {
		gen, err := c.PutOperatorHeartbeat(ctx, f.h)
		f.h.Owner.Generation = gen
		return err
	})
}

func (f *operatorHTTP) do(method, path string, body []byte, edit func(*http.Request)) *httptest.ResponseRecorder {
	f.t.Helper()
	req, err := http.NewRequestWithContext(boundedTestContext(f.t), method, f.server.URL+path, bytes.NewReader(body))
	if err != nil {
		f.t.Fatal(err)
	}
	if method == "POST" {
		req.Header.Set("Origin", opOrigin)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
	}
	if edit != nil {
		edit(req)
	}
	client := *f.server.Client()
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	res, err := client.Do(req)
	if err != nil {
		f.t.Fatal(err)
	}
	defer res.Body.Close()
	out := httptest.NewRecorder()
	out.Code, out.HeaderMap = res.StatusCode, res.Header
	if _, err := io.Copy(out.Body, res.Body); err != nil {
		f.t.Fatal(err)
	}
	return out
}

type opEvidence struct {
	db, world                              [32]byte
	actions, operations, effects, sessions int
}

func (f *operatorHTTP) evidence() (out opEvidence) {
	f.t.Helper()
	fileHash := func(p string) [32]byte {
		b, err := os.ReadFile(p)
		if err != nil {
			f.t.Fatal(err)
		}
		return sha256.Sum256(b)
	}
	out.db, out.world = fileHash(f.d.credentialsPath), fileHash(f.d.cfg.DBPath)
	u := url.URL{Scheme: "file", Path: f.d.credentialsPath, RawQuery: "mode=ro"}
	db, err := sql.Open("sqlite", u.String())
	if err != nil {
		f.t.Fatal(err)
	}
	defer db.Close()
	err = db.QueryRowContext(boundedTestContext(f.t), `SELECT (SELECT count(*) FROM operator_actions),
		(SELECT count(*) FROM operator_action_steps WHERE kind='operation'),
		(SELECT count(*) FROM operator_action_steps WHERE kind<>'operation'), (SELECT count(*) FROM session_credentials)`).
		Scan(&out.actions, &out.operations, &out.effects, &out.sessions)
	if err != nil || out.effects != 0 || out.sessions != 0 {
		f.t.Fatalf("a runner or effect crossed the browser boundary: %+v, %v", out, err)
	}
	if _, err := os.Stat(f.fake); !os.IsNotExist(err) {
		f.t.Fatalf("unexpected fake child marker: %v", err)
	}
	return out
}

func opRequest(n int) operatorRequestBody {
	o := approval.OperatorOptions{Preset: "sol", Headless: true}
	return operatorRequestBody{Argv: approval.OperatorArgv("start", fmt.Sprintf("ep-%d", n), o),
		Task: opTask, OperationID: fmt.Sprintf("%032x", n), CatalogueGeneration: 1}
}

func (f *operatorHTTP) request(body operatorRequestBody) operatorStatus {
	f.t.Helper()
	r := f.do("POST", operatorPath, opJSON(body), nil)
	var s operatorStatus
	if err := json.Unmarshal(r.Body.Bytes(), &s); r.Code != 200 || err != nil || s.ID == "" || s.Href == "" {
		f.t.Fatalf("request: %d %s (%v)", r.Code, r.Body, err)
	}
	return s
}

func (f *operatorHTTP) refuse(before opEvidence, r *httptest.ResponseRecorder, status int, class string) {
	f.t.Helper()
	requireRefusal(f.t, class, r, status, class)
	if r.Header().Get("Location") != "" || before != f.evidence() {
		f.t.Fatal("refusal changed durable evidence or built a redirect")
	}
}

func (f *operatorHTTP) postOperator(body operatorRequestBody) *httptest.ResponseRecorder {
	return f.serve("POST", operatorPath, apLoopback, opJSON(body),
		"Origin", opOrigin, "Content-Type", "application/json", "Sec-Fetch-Site", "same-origin")
}

func TestOperatorReviewRefusesForgedCatalogue(t *testing.T) {
	for _, tc := range []struct {
		name string
		edit func(*store.OperatorStartReview)
	}{
		{"Store", func(r *store.OperatorStartReview) { r.Store += ".foreign" }},
		{"Agent", func(r *store.OperatorStartReview) { r.Start.Agent = "claude" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newOperatorSurface(t, tc.edit)
			before := f.evidence()
			line := "episode start ep-1 --preset sol --launch --headless"
			bar := f.serve("GET", "/surface/verb.json?q="+url.QueryEscape(line), apLoopback, nil)
			post := f.postOperator(opRequest(1))
			var a verbAnswer
			if err := json.Unmarshal(bar.Body.Bytes(), &a); bar.Code != 200 || err != nil ||
				a.Outcome != "terminal" || a.Reason != operatorUnavailable || a.Text != "" || a.CatalogueGeneration != 0 {
				t.Errorf("forged catalogue reached review: %d %s (%v)", bar.Code, bar.Body, err)
			}
			f.refuse(before, post, 503, "RunnerUnavailable")
		})
	}
}

func TestOperatorReplayConflicts(t *testing.T) {
	f := newOperatorSurface(t, nil)
	if r := f.postOperator(opRequest(1)); r.Code != 200 {
		t.Fatalf("initial request: %d %s", r.Code, r.Body)
	}
	before := f.evidence()
	if before.actions != 1 || before.operations != 1 {
		t.Fatalf("initial request evidence: %+v", before)
	}
	for _, tc := range []struct {
		name string
		edit func(*operatorRequestBody)
	}{
		{"Argv", func(b *operatorRequestBody) { b.Argv = opRequest(2).Argv }},
		{"Generation", func(b *operatorRequestBody) { b.CatalogueGeneration++ }},
	} {
		f.sub(tc.name, func(t *testing.T) {
			b := opRequest(1)
			tc.edit(&b)
			f.refuse(before, f.postOperator(b), 409, "Conflict")
		})
	}
}

type operatorBodySpy struct{ reads int }

func (b *operatorBodySpy) Read([]byte) (int, error) { b.reads++; return 0, io.EOF }
func (b *operatorBodySpy) Close() error             { return nil }

func (f *operatorHTTP) sub(name string, run func(*testing.T)) {
	f.t.Run(name, func(t *testing.T) {
		old := f.t
		f.t = t
		defer func() { f.t = old }()
		run(t)
	})
}

func TestOperatorAdmissionHTTP(t *testing.T) {
	f := newOperatorHTTP(t)
	valid := opJSON(opRequest(1))
	before := f.evidence()
	refuse := func(path, method string, raw []byte, code int, class string, edit func(*http.Request)) {
		f.refuse(before, f.do(method, path, raw, edit), code, class)
	}
	hdr := func(k, v string) func(*http.Request) {
		return func(r *http.Request) {
			if v == "" {
				r.Header.Del(k)
			} else {
				r.Header.Set(k, v)
			}
		}
	}
	f.sub("Misdirected", func(t *testing.T) {
		refuse(operatorPath, "POST", valid, 421, "MisdirectedRequest", func(r *http.Request) { r.Host = "localhost:7644.evil" })
		spy := &operatorBodySpy{}
		r := httptest.NewRequest("POST", operatorPath, spy)
		r = r.WithContext(context.WithValue(r.Context(), http.LocalAddrContextKey, &net.TCPAddr{Port: 7644}))
		r.Host = "localhost:7644.evil"
		w := httptest.NewRecorder()
		f.handler.ServeHTTP(w, r)
		f.refuse(before, w, 421, "MisdirectedRequest")
		if spy.reads != 0 {
			t.Fatal("foreign Host reached body read")
		}
	})
	f.sub("CrossOrigin", func(t *testing.T) {
		for _, origin := range []string{"", "null", opOrigin + ".evil", "https://localhost:7644", "http://localhost:7645"} {
			refuse(operatorPath, "POST", valid, 403, "CrossOrigin", hdr("Origin", origin))
		}
		refuse(operatorPath, "POST", valid, 403, "CrossOrigin", func(r *http.Request) { r.Header.Add("Origin", opOrigin) })
	})
	f.sub("FetchSite", func(t *testing.T) {
		for _, site := range []string{"same-site", "cross-site", "same-origin, cross-site", "Same-Origin"} {
			refuse(operatorPath, "POST", valid, 403, "CrossOrigin", hdr("Sec-Fetch-Site", site))
		}
		refuse(operatorPath, "POST", valid, 403, "CrossOrigin", func(r *http.Request) { r.Header.Add("Sec-Fetch-Site", "cross-site") })
	})
	f.sub("MediaAndBody", func(t *testing.T) {
		refuse(operatorPath, "POST", valid, 415, "UnsupportedMediaType", hdr("Content-Type", "text/plain"))
		refuse(operatorPath, "POST", valid, 415, "UnsupportedMediaType", func(r *http.Request) { r.Header.Add("Content-Type", "text/plain") })
		refuse(operatorPath, "POST", bytes.Repeat([]byte(" "), (64<<10)+1), 413, "TooLarge", nil)
	})
	f.sub("UnsupportedMethodAndPath", func(t *testing.T) {
		id := "/oa-" + strings.Repeat("a", 32)
		for _, path := range []string{operatorPath, operatorPath + id, operatorPath + id + "/status.json"} {
			refuse(path, "OPTIONS", nil, 405, "MethodNotAllowed", nil)
		}
		for _, path := range []string{operatorPath + id + "/", operatorPath + id + "?next=https://elsewhere.invalid", operatorPath + "/" + id,
			operatorPath + "/x/.." + id, "/surface/x/../operator-actions" + id, operatorPath + "/%6fa-" + strings.Repeat("a", 32), operatorPath + strings.ToUpper(id),
			operatorPath + id + "x", operatorPath + id[:len(id)-1], operatorPath + id + "%0a"} {
			refuse(path, "GET", nil, 400, "BadRequest", nil)
		}
		refuse(operatorPath+id+"/status.json/x", "GET", nil, 404, "NotFound", nil)
		refuse(operatorPath+id+"/confirm/x", "POST", valid, 404, "NotFound", nil)
	})
	f.sub("StrictJSONAndLimits", func(t *testing.T) {
		invalid := [][]byte{nil, []byte("[]"), []byte(`{"argv":"start"}`), []byte(`{"argv":[42]}`), append(bytes.Clone(valid), []byte(" {}")...),
			append([]byte(`{"task":"duplicate",`), valid[1:]...), append([]byte(`{"Task":"alias",`), valid[1:]...), append([]byte(`{"other":true,`), valid[1:]...),
			bytes.Replace(valid, []byte("Task"), []byte{0xff}, 1), bytes.Replace(valid, []byte("Task"), []byte{0}, 1),
			bytes.Replace(valid, []byte("Task"), []byte(`\u0000`), 1)}
		for _, path := range invalid {
			refuse(operatorPath, "POST", path, 400, "BadRequest", nil)
		}
		for _, text := range []string{"", strings.Repeat("x", (8<<10)+1)} {
			b := opRequest(1)
			b.Task = text
			refuse(operatorPath, "POST", opJSON(b), 400, "BadRequest", nil)
		}
		for _, id := range []string{"", strings.Repeat("A", 32), "x" + opRequest(1).OperationID} {
			b := opRequest(1)
			b.OperationID = id
			refuse(operatorPath, "POST", opJSON(b), 400, "BadRequest", nil)
		}
	})
	f.sub("ExplicitHold", func(t *testing.T) {
		refuse(operatorPath, "POST", valid, 403, "Forbidden", hdr("Authorization", "Bearer session"))
		refuse(operatorPath, "POST", valid, 403, "Forbidden", hdr("Cookie", "session=x"))
	})
	checkBar := func(line, want string) {
		t.Helper()
		result := f.do("GET", "/surface/verb.json?q="+url.QueryEscape(line), nil, nil)
		var a verbAnswer
		if err := json.Unmarshal(result.Body.Bytes(), &a); err != nil || result.Code != 200 || a.Outcome != want ||
			(want == "review" && (a.CatalogueGeneration != 1 || !reflect.DeepEqual(a.Argv, opRequest(1).Argv))) || before != f.evidence() {
			t.Fatalf("bar %q: %d %s (%v)", line, result.Code, result.Body, err)
		}
	}
	checkBar("episode start ep-1 --preset=sol --headless --launch", "review")
	for _, line := range []string{"world-publish", "approver enroll", "approver revoke", "session mint", "session revoke", "service runner start"} {
		checkBar(line, "refused")
	}
	checkBar("approve mr-0123456789abcdef", "navigate")
	checkBar("episode report ep-1", "result")
	for _, v := range []struct {
		generation int64
		preset     string
	}{{2, "sol"}, {1, "absent"}} {
		b := opRequest(1)
		b.CatalogueGeneration, b.Argv[5] = v.generation, v.preset
		refuse(operatorPath, "POST", opJSON(b), 503, "RunnerUnavailable", nil)
	}
	s := f.request(opRequest(1))
	if after := f.evidence(); after.world != before.world || after.actions != 1 || after.operations != 1 || s.State != "requested" {
		t.Fatalf("inert request evidence: %+v status=%+v", after, s)
	}
	f.sub("Status", func(t *testing.T) { testOperatorStatus(f, s) })
	if duplicate := f.request(opRequest(1)); duplicate.ID != s.ID || f.evidence().actions != 1 {
		t.Fatal("same operation did not recover the existing row")
	}
	b := opRequest(1)
	b.OperationID = opRequest(90).OperationID
	if identical := f.request(b); identical.ID != s.ID || f.evidence().actions != 1 {
		t.Fatal("identical intent from a different operation created a row")
	}
	before = f.evidence()
	b.Task += " changed"
	refuse(operatorPath, "POST", opJSON(b), 409, "Conflict", nil)
	f.now.Store(apT0 + 10)
	checkBar("episode start ep-1 --preset sol --headless --launch", "terminal")
	refuse(operatorPath, "POST", opJSON(opRequest(2)), 503, "RunnerUnavailable", nil)
	if recovered := f.request(opRequest(1)); recovered.ID != s.ID || before != f.evidence() {
		t.Fatal("lost request response did not survive catalogue expiry")
	}
	f.now.Store(apT0)
	entry, err := f.h.StartReview("sol")
	if err != nil {
		t.Fatal(err)
	}
	for i := 2; i <= 16; i++ {
		b := opRequest(i)
		a, err := entry.Action(fmt.Sprintf("ep-%d", i), b.Task, 1)
		if err != nil {
			t.Fatal(err)
		}
		a.ID = "oa-" + b.OperationID
		f.write(func(ctx context.Context, db *store.Credentials) error {
			_, err := db.RequestOperatorAction(ctx, store.OperatorRequest{Action: a, OperationID: b.OperationID,
				Challenge: sha256.Sum256([]byte(a.ID)), CreatedAt: apT0, ExpiresAt: apT0 + 600}, apT0)
			return err
		})
	}
	before = f.evidence()
	if before.actions != 16 {
		t.Fatalf("pending cap positive control: %d", before.actions)
	}
	refuse(operatorPath, "POST", opJSON(opRequest(17)), 409, "QueueFull", nil)
	f.write(func(ctx context.Context, c *store.Credentials) error {
		return c.CancelUnstartedOperatorAction(ctx, s.ID, apT0)
	})
	if r := f.do("GET", "/surface/requests.json", nil, nil); r.Code != 200 || strings.Contains(r.Body.String(), s.ID) {
		t.Fatalf("cancelled row in Needs you: %d %s", r.Code, r.Body)
	}
}

func testOperatorStatus(f *operatorHTTP, s operatorStatus) {
	t := f.t
	defer f.now.Store(apT0)
	before := f.evidence()
	get := func(p string) *httptest.ResponseRecorder { return f.do("GET", p, nil, nil) }
	card := get(s.Href)
	if card.Code != 200 || card.Header().Get("Content-Security-Policy") != surfaceCSP ||
		card.Header().Get("Cache-Control") != "no-store" || !strings.Contains(card.Body.String(), "&lt;script&gt;alert(1)&lt;/script&gt;") ||
		strings.Contains(card.Body.String(), "<script>") || !strings.Contains(card.Body.String(), s.StatusHref) {
		t.Fatalf("stable protected card: %d %s", card.Code, card.Body)
	}
	for _, target := range []string{s.Href, s.StatusHref} {
		f.refuse(before, f.do("GET", target, nil, func(r *http.Request) { r.Host = "other.invalid:7644" }), 421, "MisdirectedRequest")
	}
	out := get(s.StatusHref)
	var current operatorStatus
	if err := json.Unmarshal(out.Body.Bytes(), &current); out.Code != 200 || err != nil || current.ID != s.ID ||
		current.State != "requested" || current.ExpiresAt != apT0+600 || current.RunnerSeenAt != apT0 {
		t.Fatalf("status: %d %s (%v)", out.Code, out.Body, err)
	}
	checkNeeds := func(want int) {
		t.Helper()
		r := get("/surface/requests.json")
		sum := get("/surface/home-summary.json")
		var needs surfacepage.RequestsJSON
		var home surfacepage.HomeSummary
		if json.Unmarshal(r.Body.Bytes(), &needs) != nil || json.Unmarshal(sum.Body.Bytes(), &home) != nil || len(needs.OperatorActions) != want ||
			len(home.Items) == 0 || home.Items[0].Value != fmt.Sprint(want) ||
			want == 1 && needs.OperatorActions[0].Href != s.Href {
			t.Fatalf("Needs you diverged from summary: %s / %s", r.Body, sum.Body)
		}
		for _, text := range []string{r.Body.String(), sum.Body.String(), out.Body.String(), f.log.String()} {
			if strings.Contains(text, "alert(1)") || strings.Contains(text, "challenge") {
				t.Fatalf("private review or ceremony leaked: %q", text)
			}
		}
	}
	checkNeeds(1)
	f.now.Store(apT0 - 1)
	checkNeeds(0)
	f.now.Store(apT0 + 10)
	if r := get(s.StatusHref); r.Code != 200 || !strings.Contains(r.Body.String(), "runner disconnected") || !strings.Contains(r.Body.String(), `"state":"requested"`) {
		t.Fatalf("last durable status during disconnect: %d %s", r.Code, r.Body)
	}
	f.now.Store(apT0 + 600)
	checkNeeds(0)
	if r := get(s.StatusHref); r.Code != 200 || !strings.Contains(r.Body.String(), `"state":"expired"`) || f.evidence() != before {
		t.Fatalf("expiry projection altered durable evidence: %d %s", r.Code, r.Body)
	}
}

func TestOperatorLinksProjection(t *testing.T) {
	r := store.OperatorActionRow{ID: "oa-" + strings.Repeat("a", 32), Episode: "episode-1", Kind: "start", State: "succeeded"}
	if s := operatorStatusFor(r, 1); s.ReportHref != "/v1/episodes/episode-1/report" || s.StatusHref != operatorPath+"/"+r.ID+"/status.json" {
		t.Fatalf("exact link control: %+v", s)
	}
	for _, id := range []string{r.ID + "\n", r.ID + "/x", "x" + r.ID, strings.ToUpper(r.ID), "https://github.com/other/repo/pull/1?x=y", "//elsewhere.invalid/"} {
		bad := r
		bad.ID = id
		if s := operatorStatusFor(bad, 1); s.Href != "" || s.StatusHref != "" || s.ReportHref != "" || s.State != "unavailable" {
			t.Fatalf("untrusted action link: %+v", s)
		}
	}
	for _, bad := range []store.OperatorActionRow{
		{ID: r.ID, Episode: "../x", Kind: r.Kind, State: r.State},
		{ID: r.ID, Episode: r.Episode, Kind: r.Kind, State: "other"},
		{ID: r.ID, Episode: r.Episode, Kind: "other", State: r.State},
	} {
		if s := operatorStatusFor(bad, 1); s.Href != "" || s.StatusHref != "" || s.ReportHref != "" || s.State != "unavailable" {
			t.Fatalf("unsafe or unclassified state projection: %+v", s)
		}
	}
}
