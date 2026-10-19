package main

// Row 198: `episode report` and friction #16 (--addr). The *Loopback tests
// serve a canned daemon over httptest, so they need a loopback bind (not
// available inside an episode's exec sandbox); the others need no network.
// The report's own counts are tested on a seeded store in
// host/daemon/episode_report_test.go.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/daemon"
)

func TestClientBaseURLAcceptsABareHostPort(t *testing.T) {
	for in, want := range map[string]string{
		"127.0.0.1:47330":        "http://127.0.0.1:47330",
		"localhost:7644":         "http://localhost:7644",
		"[::1]:7644":             "http://[::1]:7644",
		"http://127.0.0.1:7644":  "http://127.0.0.1:7644",
		"https://127.0.0.1:7644": "https://127.0.0.1:7644",
	} {
		if got := clientBaseURL(in); got != want {
			t.Errorf("clientBaseURL(%q) = %q, want %q", in, got, want)
		}
		if got := newClient(in).base; got != want {
			t.Errorf("newClient(%q).base = %q, want %q", in, got, want)
		}
	}
}

// The global --addr reaches the verb normalised: the error names the
// http:// URL, never Go's "first path segment in URL cannot contain colon".
func TestGlobalBareAddrIsNormalisedForEveryClientVerb(t *testing.T) {
	for _, argv := range [][]string{{"health"}, {"provenance"}, {"episode", "report", "ep1"}} {
		var out, errw bytes.Buffer
		code := run(append([]string{"--addr", "127.0.0.1:1"}, argv...), &out, &errw)
		if code != exitUsage {
			t.Errorf("%v: exit %d, want 1 (nothing listens on port 1)", argv, code)
		}
		if strings.Contains(errw.String(), "first path segment") || !strings.Contains(errw.String(), "http://127.0.0.1:1") {
			t.Errorf("%v: stderr %q, want a dial error naming http://127.0.0.1:1", argv, errw.String())
		}
	}
}

func TestEpisodeGlobalAddrRefusedForStartOnly(t *testing.T) {
	var out, errw bytes.Buffer
	if code := run([]string{"--addr", "127.0.0.1:1", "episode", "start", "ep1"}, &out, &errw); code != exitUsage ||
		!strings.Contains(errw.String(), "the global --addr is not valid for 'episode start'") {
		t.Errorf("episode start with the global --addr: exit %d stderr %q, want the refusal", code, errw.String())
	}
	out.Reset()
	errw.Reset()
	run([]string{"--addr", "127.0.0.1:1", "episode", "report", "ep1"}, &out, &errw)
	if strings.Contains(errw.String(), "not valid") || !strings.Contains(errw.String(), "episode report:") {
		t.Errorf("episode report with the global --addr: stderr %q, want it to reach the verb", errw.String())
	}
}

func TestEpisodeReportUsage(t *testing.T) {
	for _, argv := range [][]string{{"episode", "report"}, {"episode", "report", "a", "b"}, {"episode", "report", "ep", "--interval", "0s"}} {
		var out, errw bytes.Buffer
		if code := run(argv, &out, &errw); code != exitUsage {
			t.Errorf("%v: exit %d, want 1", argv, code)
		}
	}
	for _, argv := range [][]string{{"episode", "report", "--help"}, {"help", "episode", "report"}} {
		var out, errw bytes.Buffer
		if code := run(argv, &out, &errw); code != exitOK || !strings.Contains(out.String(), "episode report <episode>") {
			t.Errorf("%v: exit %d, out %q", argv, code, out.String())
		}
	}
}

func sampleReport() daemon.EpisodeReport {
	return daemon.EpisodeReport{
		Episode: "w1", Entries: [2]int64{10, 20}, TotalCalls: 4,
		CallsByTool: map[string]int{"ailang-read": 2, "ailang-cli": 1, "workspace-exec": 1},
		Effects: []daemon.ReportEffect{{Effect: "Ailang.CLI", Decision: "allowed", Count: 1},
			{Effect: "Workspace.Exec", Decision: "denied:budget", Count: 1}, {Effect: "Workspace.Read", Decision: "allowed", Count: 1}},
		Budgets: []daemon.ReportBudget{{Effect: "Workspace.Exec", Scope: "worktree", Minted: 600, Used: 91, Left: 509}},
		PlanRefusals: daemon.ReportRefusals{Count: 1, ByTool: map[string]int{"ailang-read": 1},
			First: []daemon.ReportRefusal{{Entry: 12, Tool: "ailang-read", Reason: "missing path"}}},
		HandlerRefusals: daemon.ReportRefusals{Count: 1, ByTool: map[string]int{"ailang-cli": 1},
			First: []daemon.ReportRefusal{{Entry: 15, Tool: "ailang-cli", Reason: "fmt write refuses"}}},
		FirstAt: "2026-10-19T09:00:00Z", LastAt: "2026-10-19T09:30:00Z",
	}
}

func TestEpisodeReportTablePrintsEveryClass(t *testing.T) {
	var out bytes.Buffer
	printReport(&out, sampleReport())
	for _, want := range []string{
		"episode w1: 4 calls, entries 10-20, 2026-10-19T09:00:00Z .. 2026-10-19T09:30:00Z",
		"  ailang-read             2", "Workspace.Exec         denied:budget",
		"Workspace.Exec@worktree    600     91    509",
		"plan refusals (refused before any effect): 1 (ailang-read 1)", "  #12 ailang-read: missing path",
		"handler refusals (effect allowed, handler ok:false): 1 (ailang-cli 1)", "  #15 ailang-cli: fmt write refuses",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("table lacks %q:\n%s", want, out.String())
		}
	}
}

// fakeReportDaemon serves /v1/head and the report route; bump moves the head
// and the report on.
type fakeReportDaemon struct {
	mu    sync.Mutex
	head  int
	calls int
	paths []string
}

func (f *fakeReportDaemon) bump() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.head++
	f.calls++
}

func (f *fakeReportDaemon) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.paths = append(f.paths, r.URL.EscapedPath())
	switch {
	case r.URL.Path == "/v1/head":
		fmt.Fprintf(w, "sha256:%064d", f.head)
	case r.URL.Path == "/v1/episodes/w1/report" && f.calls > 0:
		rep := sampleReport()
		rep.TotalCalls = f.calls
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(rep)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"error":{"class":"NotFound","message":"episode has no committed calls"}}`)
	}
}

func TestEpisodeReportJSONAndTableLoopback(t *testing.T) {
	f := &fakeReportDaemon{}
	f.bump()
	srv := httptest.NewServer(f)
	defer srv.Close()
	bare := strings.TrimPrefix(srv.URL, "http://")

	// --json is the route's body; the global --addr is a bare host:port.
	code, stdout, stderr := runCLI(t, bare, "episode", "report", "w1", "--json")
	if code != exitOK {
		t.Fatalf("--json: exit %d stderr %s", code, stderr)
	}
	var got daemon.EpisodeReport
	if err := json.Unmarshal([]byte(stdout), &got); err != nil || got.TotalCalls != 1 || got.Episode != "w1" {
		t.Errorf("--json printed %q (%v)", stdout, err)
	}
	// The verb's own --addr wins over the global one.
	code, stdout, _ = runCLI(t, "127.0.0.1:1", "episode", "report", "w1", "--addr", bare)
	if code != exitOK || !strings.Contains(stdout, "episode w1: 1 calls") {
		t.Errorf("table: exit %d stdout %q", code, stdout)
	}
	code, _, stderr = runCLI(t, srv.URL, "episode", "report", "nobody")
	if code != exitUsage || !strings.Contains(stderr, "episode nobody has no committed calls") {
		t.Errorf("no calls: exit %d stderr %q", code, stderr)
	}
}

func TestEpisodeReportFollowReprintsOnChangeLoopback(t *testing.T) {
	f := &fakeReportDaemon{}
	srv := httptest.NewServer(f)
	defer srv.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var out, errw syncBuffer
	done := make(chan int, 1)
	go func() {
		done <- reportFollow(ctx, newClient(srv.URL), reportOpts{episode: "w1", json: true, interval: 20 * time.Millisecond}, &out, &errw)
	}()
	lines := func() []string { return strings.Split(strings.TrimSpace(out.String()), "\n") }
	waitFor := func(n int) {
		t.Helper()
		for deadline := time.Now().Add(3 * time.Second); time.Now().Before(deadline); time.Sleep(10 * time.Millisecond) {
			if strings.TrimSpace(out.String()) != "" && len(lines()) >= n {
				return
			}
		}
		t.Fatalf("follow printed %d line(s), want %d:\n%s\nstderr %s", len(lines()), n, out.String(), errw.String())
	}
	time.Sleep(60 * time.Millisecond)
	if !strings.Contains(errw.String(), "waiting for episode w1") {
		t.Errorf("no waiting note before the first call: %q", errw.String())
	}
	f.bump()
	waitFor(1)
	time.Sleep(80 * time.Millisecond) // unchanged head: no reprint
	f.bump()
	waitFor(2)
	cancel()
	if code := <-done; code != exitOK {
		t.Errorf("follow exit %d after Ctrl-C, want 0", code)
	}
	if n := len(lines()); n != 2 {
		t.Errorf("follow printed %d lines for 2 changes:\n%s", n, out.String())
	}
	for i, l := range lines() {
		var rep daemon.EpisodeReport
		if err := json.Unmarshal([]byte(l), &rep); err != nil || rep.TotalCalls != i+1 {
			t.Errorf("line %d = %q (%v), want one compact report with %d calls", i, l, err, i+1)
		}
	}
}
