package main

// Row 201 A3 sprint 3 M8c (design §2.7 [r2], §6 XA5): the card's network deadlines in a live Chrome.

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
)

// cardNet stands between the browser and the daemon: in a stall mode the matching request gets no
// header, or a header and part of its body, and then nothing; post-recorded runs the daemon's POST
// (it records) and withholds its answer; restart answers a POST 503 Restart, get-500 the card's GET a
// 500, and list-header stalls GET /surface/requests.json (the status check). It counts POSTs.
type cardNet struct {
	mu          sync.Mutex
	mode        string
	posts, ssh  int
	release     chan struct{}
	recordedRun int
}

func (n *cardNet) set(mode string) { n.mu.Lock(); n.mode = mode; n.mu.Unlock() }
func (n *cardNet) count() int      { n.mu.Lock(); defer n.mu.Unlock(); return n.posts }

func (n *cardNet) wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		post := r.Method == "POST" && strings.HasPrefix(r.URL.Path, "/surface/approve/")
		get := r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/surface/requests/")
		list := r.Method == "GET" && r.URL.Path == "/surface/requests.json"
		if post {
			body, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewReader(body))
			n.mu.Lock()
			n.posts++
			if bytes.Contains(body, []byte(`"sshsig"`)) {
				n.ssh++
			}
			n.mu.Unlock()
		}
		n.mu.Lock()
		mode := n.mode
		n.mu.Unlock()
		stall := func() {
			select {
			case <-r.Context().Done():
			case <-n.release:
			}
		}
		switch {
		case (get && mode == "get-header") || (post && mode == "post-header") || (list && mode == "list-header"):
			stall()
		case (get && mode == "get-body") || (post && mode == "post-body"):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{"v":1,`)
			w.(http.Flusher).Flush()
			stall()
		case post && mode == "post-recorded":
			h.ServeHTTP(newDiscardWriter(), r)
			n.mu.Lock()
			n.recordedRun++
			n.mu.Unlock()
			stall()
		case get && mode == "get-500":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, `{"v":1,"error":"Internal","reason":"internal error (test)"}`)
		case post && mode == "restart":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = io.WriteString(w, `{"v":1,"error":"Restart","reason":"restart serve (test)"}`)
		default:
			h.ServeHTTP(w, r)
		}
	})
}

// discardWriter is a ResponseWriter whose answer nobody reads.
type discardWriter struct{ h http.Header }

func newDiscardWriter() *discardWriter               { return &discardWriter{h: http.Header{}} }
func (d *discardWriter) Header() http.Header         { return d.h }
func (d *discardWriter) Write(b []byte) (int, error) { return len(b), nil }
func (d *discardWriter) WriteHeader(int)             {}

// XA5: every card GET and POST, body included, ends at its 10 s deadline. A stalled GET (no header;
// a header and part of the body) shows a retryable error, and the retry reads the card. A stalled
// POST (no header; part of the body; recorded by the daemon but unanswered) shows "outcome
// unknown", never refused or minted, and its status check tells the truth (not recorded; recorded;
// minted meanwhile: no longer pending), itself under the deadline (a stalled check offers another).
// A card GET answered 500 shows the daemon's reason and the same retry.
// Only an explicit 503 Busy is retried: 1 s apart, at most 3 times (a writer held 1.5 s: 3 POSTs,
// recorded; held throughout: 4 POSTs, then refused); a 503 Restart is not retried. No timeout or
// refusal ever sends the other format.
func TestApproveCardNetworkDeadlines(t *testing.T) {
	setup := func(t *testing.T) (*passkeyBrowser, approveRig, *cardNet, string) {
		b := startPasskeyBrowser(t)
		b.do(t, nil, "tab", map[string]any{"authenticator": "uv"})
		now := time.Now().Unix()
		a := newApproveRig(t, "ep1", now)
		n := &cardNet{release: make(chan struct{})}
		origin := cardOrigin(t, a.credPath, "", n.wrap)
		t.Cleanup(func() { close(n.release) }) // before the server's Close, which waits for its handlers
		virtualPasskey(t, b, a.credPath, origin, "mark-mac", 82, now-100)
		return b, a, n, origin + "/surface/approve/" + a.id
	}
	timed := func(t *testing.T, b *passkeyBrowser, cmd string, args map[string]any) (cardView, time.Duration) {
		t.Helper()
		var v cardView
		t0 := time.Now()
		b.do(t, &v, cmd, args)
		return v, time.Since(t0)
	}
	deadline := func(t *testing.T, what string, took time.Duration) {
		t.Helper()
		if took < 9500*time.Millisecond || took > 25*time.Second {
			t.Fatalf("%s ended after %v, want the 10 s deadline", what, took)
		}
	}

	t.Run("get", func(t *testing.T) {
		b, _, n, card := setup(t)
		n.set("get-500")
		if v, _ := timed(t, b, "openCard", map[string]any{"url": card}); v.State != "failed" || v.Status != "The daemon did not give the request (HTTP 500: internal error (test))." || v.Again != "Read it again" || v.AgainHidden {
			t.Fatalf("get-500: the card shows %+v", v)
		}
		for _, mode := range []string{"get-header", "get-body"} {
			n.set(mode)
			v, took := timed(t, b, "openCard", map[string]any{"url": card, "settleMs": 25000})
			deadline(t, mode, took)
			if v.State != "failed" || !v.StatusError || v.Status != "The request could not be read (no answer within 10 s). Nothing to approve yet." ||
				v.AgainHidden || v.Again != "Read it again" || !v.ButtonHidden || len(v.Fields) != 0 {
				t.Fatalf("%s: the card shows %+v", mode, v)
			}
			n.set("")
			if v, _ = timed(t, b, "press", map[string]any{"id": "again", "states": []string{"ready", "failed", "refused"}}); v.State != "ready" {
				t.Fatalf("%s: the retry shows %+v", mode, v)
			}
		}
	})

	t.Run("post", func(t *testing.T) {
		b, a, n, card := setup(t)
		if v, _ := timed(t, b, "openCard", map[string]any{"url": card}); v.State != "ready" {
			t.Fatalf("the card shows %+v", v)
		}
		for _, tc := range []struct{ mode, after, status string }{
			{"post-header", "ready", "Not recorded: the request is still pending with no approval. You may approve it again."},
			{"post-body", "ready", "Not recorded: the request is still pending with no approval. You may approve it again."},
			{"post-recorded", "recorded", "Recorded: approved by webauthn:mark-mac, waiting for the minter (approve " + a.id + " --await)."},
		} {
			n.set(tc.mode)
			v, took := timed(t, b, "press", map[string]any{"id": "go", "states": []string{"unknown", "recorded", "rejected", "failed"}, "settleMs": 25000})
			deadline(t, tc.mode, took)
			if v.State != "unknown" || !strings.HasPrefix(v.Status, "Approval recording outcome unknown: no answer within 10 s.") || v.AgainHidden || v.Again != "Check status" {
				t.Fatalf("%s: the card shows %+v", tc.mode, v)
			}
			if _, ok := surfaceRecorded(t, a.credPath, a.id); ok != (tc.mode == "post-recorded") {
				t.Fatalf("%s: recorded %v", tc.mode, ok)
			}
			if tc.mode == "post-header" { // the status check has its own deadline, and offers another
				n.set("list-header")
				b.do(t, nil, "press", map[string]any{"id": "again"})
				var said string
				b.do(t, &said, "waitFor", map[string]any{"expression": `(() => { const s = document.getElementById('status').textContent; return s.startsWith('The status could not be read') && s; })()`, "ms": 25000})
				if said != "The status could not be read (no answer within 10 s)." {
					t.Fatalf("a stalled status check says %q", said)
				}
			}
			n.set("")
			if v, _ = timed(t, b, "press", map[string]any{"id": "again", "states": []string{"ready", "recorded", "used"}}); v.State != tc.after || v.Status != tc.status {
				t.Fatalf("%s: the status check shows %+v", tc.mode, v)
			}
		}

		// Recorded with no answer, then minted before the check: the check says it is no longer pending.
		if v, _ := timed(t, b, "openCard", map[string]any{"url": card}); v.State != "ready" {
			t.Fatalf("the card shows %+v", v)
		}
		n.set("post-recorded")
		if v, _ := timed(t, b, "press", map[string]any{"id": "go", "states": []string{"unknown", "recorded", "rejected"}, "settleMs": 25000}); v.State != "unknown" {
			t.Fatalf("post-recorded again: the card shows %+v", v)
		}
		n.set("")
		var opened bool
		aenv, _ := awaitEnv(time.Now().Unix(), &opened, nil)
		var stdout, stderr bytes.Buffer
		if code := runApprove(a.args("--await", "--await-for", "5s"), &stdout, &stderr, aenv); code != exitOK || opened {
			t.Fatalf("approve --await: rc %d\nstderr %s", code, stderr.String())
		}
		b.do(t, nil, "press", map[string]any{"id": "again"})
		var said string
		b.do(t, &said, "waitFor", map[string]any{"expression": `(() => { const s = document.getElementById('status').textContent; return s.startsWith('The request is no longer') && s; })()`})
		if said != "The request is no longer pending (used or expired): read the card again." {
			t.Fatalf("after the mint the status check says %q", said)
		}

		a2 := newApproveRigIn(t, a.r, "ep2", time.Now().Unix()) // a pending request for the Busy arms
		card = strings.TrimSuffix(card, a.id) + a2.id
		busy := func(hold time.Duration) (cardView, int) {
			t.Helper()
			if v, _ := timed(t, b, "openCard", map[string]any{"url": card}); v.State != "ready" {
				t.Fatalf("the card shows %+v", v)
			}
			w, err := store.OpenCredentials(a.credPath) // the daemon's non-waiting writer is refused: 503 Busy
			if err != nil {
				t.Fatal(err)
			}
			var once sync.Once
			free := func() { once.Do(func() { _ = w.Close() }) }
			defer free()
			if hold > 0 {
				time.AfterFunc(hold, free)
			}
			before := n.count()
			v, _ := timed(t, b, "press", map[string]any{"id": "go", "states": []string{"recorded", "rejected", "unknown", "failed"}, "settleMs": 20000})
			return v, n.count() - before
		}
		if v, posts := busy(1500 * time.Millisecond); v.State != "recorded" || posts != 3 {
			t.Fatalf("a writer held 1.5 s: %d POSTs, the card shows %+v", posts, v)
		}
		if v, posts := busy(0); v.State != "rejected" || posts != 4 || !strings.Contains(v.Status, "(HTTP 503 Busy): the credentials store is busy") {
			t.Fatalf("a writer held throughout: %d POSTs, the card shows %+v", posts, v)
		}
		if v, _ := timed(t, b, "openCard", map[string]any{"url": card}); v.State != "ready" {
			t.Fatalf("the card shows %+v", v)
		}
		n.set("restart")
		before := n.count()
		if v, _ := timed(t, b, "press", map[string]any{"id": "go", "states": []string{"recorded", "rejected", "unknown"}}); v.State != "rejected" || n.count()-before != 1 || !strings.Contains(v.Status, "HTTP 503 Restart") {
			t.Fatalf("a 503 Restart: %d POSTs, the card shows %+v", n.count()-before, v)
		}
		n.mu.Lock()
		defer n.mu.Unlock()
		if n.ssh != 0 || n.recordedRun != 2 {
			t.Fatalf("%d POSTs sent an SSH signature; the withheld POST ran %d times", n.ssh, n.recordedRun)
		}
	})
}

// XA5, the passkey's own bound (A2 §2.14 [r2], design §2.7 step 4): a get that never answers is
// aborted at 240 s of page time and nothing is sent; a result after the 240 s is ignored and
// nothing is sent; an authenticator without user verification makes no response and nothing is
// sent. Each leaves the button usable.
func TestApproveCardPasskeyDeadline(t *testing.T) {
	b := startPasskeyBrowser(t)
	now := time.Now().Unix()
	a := newApproveRig(t, "ep1", now)
	n := &cardNet{release: make(chan struct{})}
	origin := cardOrigin(t, a.credPath, "", n.wrap)
	t.Cleanup(func() { close(n.release) })
	card := origin + "/surface/approve/" + a.id
	for i, tc := range []struct{ inject, authenticator, state, status string }{
		{"abort", "uv", "expired", "Timed out: no passkey answer within 4 minutes. Nothing was sent; try again."},
		{"late", "uv", "expired", "Timed out: no passkey answer within 4 minutes. Nothing was sent; try again."},
		{"", "no-uv", "failed", "The browser did not make a passkey response (NotAllowedError: "},
	} {
		b.do(t, nil, "tab", map[string]any{"authenticator": tc.authenticator, "inject": tc.inject})
		if i == 0 {
			virtualPasskey(t, b, a.credPath, origin, "mark-mac", 84, now-100)
		}
		var v cardView
		if b.do(t, &v, "openCard", map[string]any{"url": card}); v.State != "ready" {
			t.Fatalf("%s: the card shows %+v", tc.inject, v)
		}
		args := map[string]any{"id": "go", "states": []string{"expired", "failed", "rejected", "recorded"}}
		if tc.inject != "" {
			args = map[string]any{"id": "go", "virtualBudget": 250000}
		}
		var after cardView
		b.do(t, &after, "press", args)
		if after.State != tc.state || !strings.HasPrefix(after.Status, tc.status) || after.ButtonDisabled || n.count() != 0 {
			t.Fatalf("%s/%s: %d POSTs, the card shows %+v", tc.inject, tc.authenticator, n.count(), after)
		}
		if tc.inject != "" && (after.Op == nil || !after.Op.Signal || after.Op.AbortedAfter == nil || *after.Op.AbortedAfter < 239000 || *after.Op.AbortedAfter > 241000) {
			t.Fatalf("%s: the get was not aborted at 240 s: %s", tc.inject, opString(after.Op))
		}
	}
}
