package daemon

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sunholo-data/ailang-world/host/agui"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/workbench"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

type workbenchLatestReader interface {
	LogEntriesLatest(context.Context, int) ([]store.LogEntry, error)
}
type workbenchLatestReads struct {
	readStore
	calls  int
	limits []int
	fail   bool
}

func (p *workbenchLatestReads) LogEntriesLatest(ctx context.Context, n int) ([]store.LogEntry, error) {
	p.calls++
	p.limits = append(p.limits, n)
	if p.fail {
		return nil, errors.New("SECRET latest read detail")
	}
	return p.readStore.(workbenchLatestReader).LogEntriesLatest(ctx, n)
}
func workbenchLiveBody(t *testing.T, d *Daemon, path string) string {
	t.Helper()
	r := requestRecorder(t, d, "GET", path, nil)
	if r.Code != 200 {
		t.Fatalf("page %s: %d %s", path, r.Code, r.Body)
	}
	return r.Body.String()
}
func workbenchLiveCursor(t *testing.T, body string) int64 {
	t.Helper()
	m := regexp.MustCompile(`data-live-cursor="(-?[0-9]+)"`).FindStringSubmatch(body)
	if len(m) != 2 {
		t.Fatal("rendered live cursor missing")
	}
	n, err := strconv.ParseInt(m[1], 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	return n
}
func workbenchLiveRegion(t *testing.T, body string) string {
	t.Helper()
	start := strings.Index(body, `<section aria-label="live"`)
	if start < 0 {
		t.Fatal("live region missing")
	}
	end := strings.Index(body[start:], "</section>")
	if end < 0 {
		t.Fatal("live region not closed")
	}
	return body[start : start+end]
}
func workbenchLiveREST(t *testing.T, d *Daemon, indexes []int64) store.Commit {
	t.Helper()
	auth := authHeader(t, d)
	c := genesisCommit("workbench-live")
	for j, i := range indexes {
		if j > 0 {
			c = testCommit(c.NextWorld, i, "workbench-live")
		}
		r := requestRecorderAuth(t, d, auth, "POST", "/v1/commit", bytes.NewReader(encodeCommit(c)))
		if r.Code != 200 {
			t.Fatalf("REST commit %d: %d %s", i, r.Code, r.Body)
		}
		t.Logf("REST commit %d: 200", i)
	}
	return c
}
func TestWorkbenchLiveCursorCrossesGap(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		d := newHandlerDaemon(t)
		if got := workbenchLiveCursor(t, workbenchLiveBody(t, d, "/workbench")); got != -1 {
			t.Fatalf("empty cursor=%d want -1", got)
		}
	})
	t.Run("gap", func(t *testing.T) {
		d := newHandlerDaemon(t)
		c := workbenchLiveREST(t, d, []int64{0, 1, 5})
		spy := &workbenchLatestReads{readStore: d.reads}
		d.reads = spy
		for _, path := range []string{"/workbench", "/workbench?from=0&entry=5", "/workbench?object=" + c.Objects[0].Hash.String()} {
			body := workbenchLiveBody(t, d, path)
			if got := workbenchLiveCursor(t, body); got != 5 {
				t.Fatalf("gap cursor=%d want 5", got)
			}
			if !strings.Contains(workbenchLiveRegion(t, body), c.Entry.EntryHash.String()) {
				t.Error("live head missing")
			}
		}
		if spy.calls != 3 || !reflect.DeepEqual(spy.limits, []int{10, 10, 10}) {
			t.Fatalf("latest read calls=%d limits=%v", spy.calls, spy.limits)
		}
	})
	t.Run("tail-and-timeline", func(t *testing.T) {
		d := newHandlerDaemon(t)
		var indexes []int64
		for i := int64(0); i < 112; i++ {
			indexes = append(indexes, i)
		}
		workbenchLiveREST(t, d, indexes)
		body := workbenchLiveBody(t, d, "/workbench")
		live := workbenchLiveRegion(t, body)
		rows := regexp.MustCompile(`>select entry ([0-9]+)</a>`).FindAllStringSubmatch(live, -1)
		var got []int
		for _, r := range rows {
			n, _ := strconv.Atoi(r[1])
			got = append(got, n)
		}
		if !reflect.DeepEqual(got, []int{111, 110, 109, 108, 107, 106, 105, 104, 103, 102}) {
			t.Fatalf("newest 10=%v", got)
		}
		if workbenchLiveCursor(t, body) != 111 {
			t.Fatal("tail cursor incorrect")
		}
		if strings.Count(body, "<h3>entry ") != 100 || !strings.Contains(body, "<h3>entry 0</h3>") || !strings.Contains(body, "<h3>entry 99</h3>") || strings.Contains(body, "<h3>entry 100</h3>") {
			t.Fatal("oldest-first timeline changed")
		}
		next := workbenchLiveBody(t, d, "/workbench?from=100&entry=100")
		if strings.Count(next, "<h3>entry ") != 12 || !strings.Contains(next, "<h3>entry 100</h3>") {
			t.Fatal("timeline pagination changed")
		}
	})
}
func TestWorkbenchLiveCursorIsResumable(t *testing.T) {
	for _, race := range []bool{false, true} {
		t.Run(fmt.Sprintf("nonempty-race-%t", race), func(t *testing.T) {
			d := newHandlerDaemon(t)
			d.aguiBudget = 300 * time.Millisecond
			d.aguiTick = 5 * time.Millisecond
			c := workbenchLiveREST(t, d, []int64{0, 1, 5})
			cursor := workbenchLiveCursor(t, workbenchLiveBody(t, d, "/workbench"))
			if race {
				c = testCommit(c.NextWorld, 6, "render-post-race")
				if err := d.store.Commit(boundedTestContext(t), c); err != nil {
					t.Fatal(err)
				}
			}
			server := httptest.NewServer(d.Handler())
			defer server.Close()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			body := fmt.Sprintf(`{"threadId":"workbench","runId":"resume","messages":[],"state":{"schema":"world/agui-state/v1","lastIndex":%d}}`, cursor)
			req, _ := http.NewRequestWithContext(ctx, "POST", server.URL+"/agui/", strings.NewReader(body))
			resp, err := server.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				t.Fatalf("rendered cursor %d resume status=%d want 200", cursor, resp.StatusCode)
			}
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			got := aguiIndexes(t, data)
			want := []int64{}
			if race {
				want = []int64{6}
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("resume entries=%v want %v", got, want)
			}
			if !bytes.Contains(data, []byte(`"type":"RUN_FINISHED"`)) {
				t.Fatal("bounded run lacks terminal cursor")
			}
		})
	}
	t.Run("empty-live-entry-zero", func(t *testing.T) {
		d := newHandlerDaemon(t)
		d.aguiBudget = 5 * time.Second
		d.aguiTick = 5 * time.Millisecond
		cursor := workbenchLiveCursor(t, workbenchLiveBody(t, d, "/workbench"))
		polls := aguiPollReads{d.reads, make(chan struct{}, 1)}
		d.reads = polls
		server := httptest.NewServer(d.Handler())
		defer server.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		body := fmt.Sprintf(`{"threadId":"workbench","runId":"empty","messages":[],"state":{"schema":"world/agui-state/v1","lastIndex":%d}}`, cursor)
		req, _ := http.NewRequestWithContext(ctx, "POST", server.URL+"/agui/", strings.NewReader(body))
		resp, err := server.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Fatalf("empty rendered cursor %d resume status=%d want 200", cursor, resp.StatusCode)
		}
		frames := make(chan string, 16)
		done := make(chan struct{})
		go func() {
			defer close(done)
			scan := bufio.NewScanner(resp.Body)
			var frame strings.Builder
			for scan.Scan() {
				line := scan.Text()
				if line == "" {
					select {
					case frames <- frame.String():
					case <-ctx.Done():
						return
					}
					frame.Reset()
				} else {
					frame.WriteString(line)
					frame.WriteByte('\n')
				}
			}
		}()
		defer func() {
			cancel()
			resp.Body.Close()
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				t.Error("stream reader failed to join")
			}
		}()
		select {
		case <-polls.polls:
		case <-time.After(5 * time.Second):
			t.Fatal("initial poll not observed")
		}
		if err := d.store.Commit(boundedTestContext(t), genesisCommit("live-zero")); err != nil {
			t.Fatal(err)
		}
		timer := time.NewTimer(5 * time.Second)
		defer timer.Stop()
		custom, delta := false, false
		for !custom || !delta {
			select {
			case frame := <-frames:
				if strings.Contains(frame, `"type":"CUSTOM"`) && strings.Contains(frame, `"entryIndex":0`) {
					custom = true
				}
				if strings.Contains(frame, `"type":"STATE_DELTA"`) && strings.Contains(frame, `"path":"/lastIndex","value":0`) {
					delta = true
				}
			case <-done:
				t.Fatal("stream ended before entry zero and delta")
			case <-timer.C:
				t.Fatal("entry zero not delivered on genesis stream")
			}
		}
	})
}
func TestWorkbenchLiveStoreError(t *testing.T) {
	d := newHandlerDaemon(t)
	spy := &workbenchLatestReads{readStore: d.reads, fail: true}
	d.reads = spy
	r := requestRecorder(t, d, "GET", "/workbench", nil)
	if r.Code != 500 {
		t.Fatalf("latest read failure status=%d want 500", r.Code)
	}
	want := httptest.NewRecorder()
	writeWorkbenchError(want, 500, "Internal", workbenchInternalStoreFailureMessage)
	if r.Body.String() != want.Body.String() || strings.Contains(r.Body.String(), "SECRET") || strings.Contains(r.Body.String(), "<script") {
		t.Fatalf("error HTML differs from constant sanitized page: %s", r.Body)
	}
}

func liveAsset(t *testing.T) []byte {
	t.Helper()
	b := workbench.LiveScript
	if len(b) == 0 {
		t.Fatal("live asset empty")
	}
	return b
}
func TestWorkbenchLiveScriptTag(t *testing.T) {
	d := newHandlerDaemon(t)
	body := workbenchLiveBody(t, d, "/workbench")
	if strings.Count(body, "<script") != 1 || !strings.Contains(body, `<script src="/workbench/live.js" defer></script>`) {
		t.Fatal("success page needs exactly one external deferred script")
	}
	if regexp.MustCompile(`\bon[a-z]+\s*=`).MatchString(body) {
		t.Fatal("inline event handler")
	}
}
func TestWorkbenchErrorPageInert(t *testing.T) {
	d := newHandlerDaemon(t)
	if !strings.Contains(workbenchLiveBody(t, d, "/workbench"), `<script src="/workbench/live.js" defer></script>`) {
		t.Fatal("success-script positive control missing")
	}
	for _, tc := range []struct {
		path   string
		status int
	}{{"/workbench?bogus=1", 400}, {"/workbench?from=0&entry=99", 404}, {"/workbench", 500}} {
		t.Run(strconv.Itoa(tc.status), func(t *testing.T) {
			if tc.status == 500 {
				d.reads = &workbenchLatestReads{readStore: d.reads, fail: true}
			}
			r := requestRecorder(t, d, "GET", tc.path, nil)
			if r.Code != tc.status {
				t.Fatalf("status=%d want %d", r.Code, tc.status)
			}
			if strings.Contains(r.Body.String(), "<script") {
				t.Fatal("error page contains script")
			}
		})
	}
}
func TestWorkbenchLiveScriptRoute(t *testing.T) {
	d := newHandlerDaemon(t)
	if !d.isProtected(httptest.NewRequest("POST", "/v1/commit", nil)) {
		t.Fatal("protected positive control")
	}
	for _, path := range []string{"/workbench/live.js", "/workbench"} {
		if d.isProtected(httptest.NewRequest("GET", path, nil)) {
			t.Fatalf("GET %s protected", path)
		}
	}
	r := requestRecorder(t, d, "GET", "/workbench/live.js", nil)
	if r.Code != 200 {
		t.Fatalf("script status=%d want 200", r.Code)
	}
	for k, v := range map[string]string{"Content-Type": "text/javascript; charset=utf-8", "X-Content-Type-Options": "nosniff", "Cache-Control": "no-store", "Content-Security-Policy": "default-src 'none'"} {
		if r.Header().Get(k) != v {
			t.Errorf("%s=%q want %q", k, r.Header().Get(k), v)
		}
	}
	if !bytes.Equal(r.Body.Bytes(), liveAsset(t)) {
		t.Fatal("script differs from asset")
	}
	if r := requestRecorder(t, d, "POST", "/workbench/live.js", nil); r.Code != 405 {
		t.Fatalf("POST script=%d want 405", r.Code)
	}
}
func TestLiveScriptContract(t *testing.T) {
	script := string(liveAsset(t))
	extract := func(pattern string) string {
		t.Helper()
		m := regexp.MustCompile(pattern).FindStringSubmatch(script)
		if len(m) != 2 || m[1] == "" {
			t.Fatalf("missing extraction: %s", pattern)
		}
		return m[1]
	}
	if got := extract(`const STATE_SCHEMA = "([^"]+)";`); got != agui.StateSchema {
		t.Fatalf("schema=%s", got)
	}
	var regions []string
	if err := json.Unmarshal([]byte(extract(`const REGIONS = (\[[^\n]+\]);`)), &regions); err != nil {
		t.Fatal(err)
	}
	want := []string{`nav[aria-label="world browser"]`, `section[aria-label="timeline"]`, `section[aria-label="live"]`, `section[aria-label="world graph"]`, `section[aria-label="decisions"]`}
	if !reflect.DeepEqual(regions, want) {
		t.Fatalf("REGIONS=%v want exact five without footer %v", regions, want)
	}
	if got := extract(`fetch\("([^"]+)", \{method: "POST"`); got != "/agui/" {
		t.Fatalf("POST path=%s", got)
	}
	empty, seeded := newHandlerDaemon(t), newHandlerDaemon(t)
	c := workbenchLiveREST(t, seeded, []int64{0})
	variants := []struct {
		name string
		d    *Daemon
		path string
	}{{"empty", empty, "/workbench"}, {"home", seeded, "/workbench"}, {"entry", seeded, "/workbench?from=0&entry=0"}, {"object", seeded, "/workbench?object=" + c.Objects[0].Hash.String()}}
	if len(variants) == 0 {
		t.Fatal("no variants")
	}
	for _, v := range variants {
		t.Run(v.name, func(t *testing.T) {
			body := workbenchLiveBody(t, v.d, v.path)
			for _, selector := range regions {
				label := regexp.MustCompile(`aria-label="([^"]+)"`).FindStringSubmatch(selector)
				if len(label) != 2 {
					t.Fatal("invalid selector")
				}
				if n := strings.Count(body, `aria-label="`+label[1]+`"`); n != 1 {
					t.Fatalf("%s count=%d", selector, n)
				}
			}
		})
	}
}
