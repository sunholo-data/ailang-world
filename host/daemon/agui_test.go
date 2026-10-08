package daemon

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sunholo-data/ailang-world/host/authority"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

var _ readStore = (*store.Store)(nil)

func TestAGUIReadStoreSeam(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "seam.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	h := hashref.SumSHA256([]byte("seam"))
	if err := s.Commit(ctx, store.Commit{
		NextWorld: store.World{Ref: h, Revision: 0, StateRoot: h, LogHead: h},
		Entry:     store.LogEntry{Header: store.LogHeader{EntryIndex: 0, SemanticsEpoch: 1, TransitionFn: h, Interpreter: h, PrevEntryHash: h, WrittenBy: "seam"}, EntryHash: h, TransitionRef: h},
	}); err != nil {
		t.Fatal(err)
	}
	var reads readStore = s
	entries, err := reads.LogEntriesAfter(ctx, -1, 100)
	if err != nil || len(entries) != 1 || entries[0].Header.EntryIndex != 0 || entries[0].EntryHash != h {
		t.Fatalf("readStore keyset: %+v err=%v", entries, err)
	}
}

// AC3.14 is intentionally the first M3 oracle: REST preconditions precede AGUI.
func TestAGUIRESTGapCrossed(t *testing.T) {
	d := newHandlerDaemon(t)
	d.aguiBudget = 30 * time.Millisecond
	d.aguiTick = 5 * time.Millisecond
	auth := authHeader(t, d)
	c := genesisCommit("agui-gap")
	for _, i := range []int64{0, 1, 5} {
		if i != 0 {
			c = testCommit(c.NextWorld, i, "agui-gap")
		}
		rec := requestRecorderAuth(t, d, auth, "POST", "/v1/commit", bytes.NewReader(encodeCommit(c)))
		if rec.Code != 200 {
			t.Fatalf("REST commit %d: %d %s", i, rec.Code, rec.Body)
		}
		t.Logf("REST commit %d: 200", i)
	}
	rec := requestRecorder(t, d, "GET", "/v1/log?from=0", nil)
	var rangeBody struct {
		Items []struct {
			Header struct {
				EntryIndex int64 `json:"entryIndex"`
			} `json:"header"`
		} `json:"items"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &rangeBody); err != nil {
		t.Fatal(err)
	}
	var indexes []int64
	for _, item := range rangeBody.Items {
		indexes = append(indexes, item.Header.EntryIndex)
	}
	if !reflect.DeepEqual(indexes, []int64{0, 1}) {
		t.Fatalf("frozen range: %s indexes=%v", rec.Body, indexes)
	}
	if rec = requestRecorder(t, d, "GET", "/v1/log/5", nil); rec.Code != 200 {
		t.Fatalf("entry5: %d", rec.Code)
	}
	t.Log("frozen range [0 1], entry5 200")
	// No new production seam is needed for this initial red.
	req := httptest.NewRequest("POST", "/agui/", bytes.NewBufferString(`{"threadId":"t","runId":"r","messages":[]}`))
	rec = httptest.NewRecorder()
	d.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("AGUI genesis status=%d, want200: %s", rec.Code, rec.Body)
	}
	var got []int64
	for _, frame := range bytes.Split(rec.Body.Bytes(), []byte("\n\n")) {
		for _, line := range bytes.Split(frame, []byte("\n")) {
			if !bytes.HasPrefix(line, []byte("data: ")) {
				continue
			}
			var e struct {
				Type  string `json:"type"`
				Value struct {
					Header struct {
						EntryIndex int64 `json:"entryIndex"`
					} `json:"header"`
				} `json:"value"`
			}
			if err := json.Unmarshal(line[6:], &e); err != nil {
				t.Fatal(err)
			}
			if e.Type == "CUSTOM" {
				got = append(got, e.Value.Header.EntryIndex)
			}
		}
	}
	if !reflect.DeepEqual(got, []int64{0, 1, 5}) {
		t.Fatalf("AGUI crosses gap: %v", got)
	}
}

const aguiPayload = `{"threadId":"t","runId":"r","messages":[]}`

func aguiTestDaemon(t *testing.T) *Daemon {
	d := newHandlerDaemon(t)
	d.aguiBudget = 30 * time.Millisecond
	d.aguiTick = 5 * time.Millisecond
	return d
}
func aguiFixture(t *testing.T, d *Daemon) []store.LogEntry {
	t.Helper()
	var entries []store.LogEntry
	for _, i := range []int64{0, 1, 2, 5, 6} {
		ref := func(n int64) hashref.HashRef {
			h, err := hashref.Parse(fmt.Sprintf("sha256:%064x", n))
			if err != nil {
				t.Fatal(err)
			}
			return h
		}
		e := store.LogEntry{Header: store.LogHeader{EntryIndex: i, SemanticsEpoch: 1, TransitionFn: ref(100), Interpreter: ref(101), PrevEntryHash: ref(i + 100), WrittenBy: "fixture <>&"}, EntryHash: ref(i + 1), TransitionRef: ref(i + 200)}
		c := store.Commit{Entry: e, NextWorld: store.World{Ref: ref(i + 500), Revision: i, StateRoot: ref(100), LogHead: e.EntryHash}}
		if len(entries) > 0 {
			c.ObservedHead = ref(entries[len(entries)-1].Header.EntryIndex + 500)
		}
		if err := d.store.Commit(boundedTestContext(t), c); err != nil {
			t.Fatal(err)
		}
		entries = append(entries, e)
	}
	return entries
}
func aguiRequest(t *testing.T, d *Daemon, body, header string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest("POST", "/agui/", strings.NewReader(body))
	if header != "" {
		req.Header.Set("Last-Event-ID", header)
	}
	rec := httptest.NewRecorder()
	d.Handler().ServeHTTP(rec, req)
	return rec
}
func aguiEvents(t *testing.T, b []byte) []map[string]json.RawMessage {
	t.Helper()
	var out []map[string]json.RawMessage
	frames := bytes.Split(b, []byte("\n\n"))
	for _, f := range frames[:len(frames)-1] {
		for _, line := range bytes.Split(f, []byte("\n")) {
			if bytes.HasPrefix(line, []byte("data: ")) {
				var e map[string]json.RawMessage
				if err := json.Unmarshal(line[6:], &e); err != nil {
					t.Fatal(err)
				}
				out = append(out, e)
			}
		}
	}
	return out
}
func aguiIndexes(t *testing.T, b []byte) []int64 {
	t.Helper()
	out := []int64{}
	for _, e := range aguiEvents(t, b) {
		if string(e["type"]) == `"CUSTOM"` {
			var v struct {
				Header struct {
					EntryIndex int64 `json:"entryIndex"`
				} `json:"header"`
			}
			if err := json.Unmarshal(e["value"], &v); err != nil {
				t.Fatal(err)
			}
			out = append(out, v.Header.EntryIndex)
		}
	}
	return out
}
func requireIndexes(t *testing.T, b []byte, want []int64) {
	t.Helper()
	if got := aguiIndexes(t, b); !reflect.DeepEqual(got, want) {
		t.Fatalf("indexes=%v want=%v", got, want)
	}
}
func TestAGUIGoldenRun(t *testing.T) {
	d := aguiTestDaemon(t)
	aguiFixture(t, d)
	d.aguiBudget = 300 * time.Millisecond
	golden, err := os.ReadFile("../agui/testdata/stream_fixture.golden")
	if err != nil {
		t.Fatal(err)
	}
	want := append(golden, []byte("data: {\"type\":\"RUN_FINISHED\",\"threadId\":\"t\",\"runId\":\"r\",\"result\":{\"lastIndex\":6}}\n\n")...)
	for n := 0; n < 2; n++ {
		rec := aguiRequest(t, d, aguiPayload, "")
		if rec.Code != 200 || !bytes.Equal(rec.Body.Bytes(), want) {
			t.Fatalf("golden run %d status=%d bytes differ: %s", n, rec.Code, rec.Body)
		}
	}
	t.Run("full-page", func(t *testing.T) {
		d := aguiTestDaemon(t)
		w := store.World{LogHead: hashref.SumSHA256([]byte("initial"))}
		want := []int64{}
		for i := int64(0); i < 205; i++ {
			c := testCommit(w, i, fmt.Sprintf("page-%d", i))
			if err := d.store.Commit(boundedTestContext(t), c); err != nil {
				t.Fatal(err)
			}
			w = c.NextWorld
			want = append(want, i)
		}
		requireIndexes(t, aguiRequest(t, d, aguiPayload, "").Body.Bytes(), want)
	})
}
func TestAGUIResumeExact(t *testing.T) {
	d := aguiTestDaemon(t)
	aguiFixture(t, d)
	for _, k := range []int64{-1, 0, 1, 2, 5, 6} {
		t.Run(fmt.Sprint(k), func(t *testing.T) {
			want := []int64{}
			for _, i := range []int64{0, 1, 2, 5, 6} {
				if i > k {
					want = append(want, i)
				}
			}
			h := aguiRequest(t, d, aguiPayload, fmt.Sprint(k))
			s := aguiRequest(t, d, fmt.Sprintf(`{"threadId":"t","runId":"r","messages":[],"state":{"schema":"world/agui-state/v1","lastIndex":%d,"logHead":"forged"}}`, k), "")
			if h.Code != 200 || s.Code != 200 || !bytes.Equal(h.Body.Bytes(), s.Body.Bytes()) {
				t.Fatalf("header/state resume differs: %d/%d", h.Code, s.Code)
			}
			requireIndexes(t, h.Body.Bytes(), want)
			match := aguiRequest(t, d, fmt.Sprintf(`{"threadId":"t","runId":"r","messages":[],"state":{"schema":"world/agui-state/v1","lastIndex":%d}}`, k), fmt.Sprint(k))
			if !bytes.Equal(h.Body.Bytes(), match.Body.Bytes()) {
				t.Fatal("matching cursors differ")
			}
		})
	}
	for _, tc := range []struct {
		state, header string
		status        int
	}{{`{}`, "", 200}, {`[]`, "", 200}, {`5`, "", 200}, {`"initial"`, "", 200}, {`null`, "", 200}, {`{}`, "1", 200}, {`{"schema":"other","lastIndex":5}`, "", 200}, {`{"schema":"world/agui-state/v1","lastIndex":"2"}`, "", 400}, {`{"schema":"world/agui-state/v1","lastIndex":-2}`, "", 400}, {`{"schema":"world/agui-state/v1","lastIndex":1.5}`, "", 400}, {`{"schema":"world/agui-state/v1","lastIndex":2}`, "1", 400}, {`{"schema":"world/agui-state/v1","lastIndex":3}`, "", 404}, {`{}`, "999", 404}, {`{}`, "oops", 400}, {`{}`, "9223372036854775808", 400}, {`{}`, "-2", 400}} {
		rec := aguiRequest(t, d, `{"threadId":"t","runId":"r","messages":[],"state":`+tc.state+`}`, tc.header)
		if rec.Code != tc.status {
			t.Fatalf("state=%s header=%s status=%d want%d", tc.state, tc.header, rec.Code, tc.status)
		}
		if tc.status != 200 && bytes.Contains(rec.Body.Bytes(), []byte("data:")) {
			t.Fatal("refusal contains SSE")
		}
	}
	rec := aguiRequest(t, d, `{"threadId":"t","runId":"r","messages":[],"forwardedProps":{"world":{"after":6}}}`, "")
	requireIndexes(t, rec.Body.Bytes(), []int64{0, 1, 2, 5, 6})
}

// Test-local stock state application: unknown operations are a hard failure.
func aguiState(t *testing.T, b []byte) map[string]any {
	t.Helper()
	s := map[string]any{}
	for _, e := range aguiEvents(t, b) {
		switch string(e["type"]) {
		case `"STATE_SNAPSHOT"`:
			if err := json.Unmarshal(e["snapshot"], &s); err != nil {
				t.Fatal(err)
			}
		case `"STATE_DELTA"`:
			var ops []struct {
				Op, Path string
				Value    any
			}
			if err := json.Unmarshal(e["delta"], &ops); err != nil {
				t.Fatal(err)
			}
			for _, op := range ops {
				if op.Op != "replace" && op.Op != "add" {
					t.Fatalf("unsupported op %s", op.Op)
				}
				if op.Path != "/lastIndex" && op.Path != "/logHead" {
					t.Fatalf("unsupported path %s", op.Path)
				}
				s[op.Path[1:]] = op.Value
			}
		}
	}
	return s
}
func TestAGUIStockClientStateResume(t *testing.T) {
	d := aguiTestDaemon(t)
	aguiFixture(t, d)
	state := aguiState(t, aguiRequest(t, d, aguiPayload, "").Body.Bytes())
	h, _, err := d.store.SelectedHead(boundedTestContext(t))
	if err != nil {
		t.Fatal(err)
	}
	w, _, err := d.store.GetWorld(boundedTestContext(t), h)
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range []int64{7, 8} {
		c := testCommit(w, i, fmt.Sprint(i))
		if err := d.store.Commit(boundedTestContext(t), c); err != nil {
			t.Fatal(err)
		}
		w = c.NextWorld
	}
	b, _ := json.Marshal(state)
	rec := aguiRequest(t, d, `{"threadId":"t","runId":"r","messages":[],"state":`+string(b)+`}`, "")
	requireIndexes(t, rec.Body.Bytes(), []int64{7, 8})
}
func TestAGUIEntryValueEqualsLogRoute(t *testing.T) {
	d := aguiTestDaemon(t)
	aguiFixture(t, d)
	rec := aguiRequest(t, d, aguiPayload, "")
	n := 0
	for _, e := range aguiEvents(t, rec.Body.Bytes()) {
		if string(e["type"]) != `"CUSTOM"` {
			continue
		}
		var v struct {
			Header struct {
				EntryIndex int64 `json:"entryIndex"`
			} `json:"header"`
		}
		json.Unmarshal(e["value"], &v)
		log := requestRecorder(t, d, "GET", fmt.Sprintf("/v1/log/%d", v.Header.EntryIndex), nil)
		if !bytes.Equal(append(e["value"], '\n'), log.Body.Bytes()) {
			t.Fatalf("raw value differs for %d", v.Header.EntryIndex)
		}
		if !bytes.Contains(e["value"], []byte(`\u003c\u003e\u0026`)) {
			t.Fatal("escaping control absent")
		}
		n++
	}
	if n != 5 {
		t.Fatalf("checked %d entries", n)
	}
}
func TestAGUIBudgetDerivation(t *testing.T) {
	d := newHandlerDaemon(t)
	if aguiRunBudget != writeTimeout-readDeadline-aguiWriteMargin {
		t.Fatal("budget derivation violated")
	}
	if aguiRunBudget != 18*time.Second || d.aguiBudget != aguiRunBudget || aguiRunBudget <= 0 || aguiWriteMargin != 2*time.Second || d.aguiBodyBound != 2*time.Second || d.aguiTick != 250*time.Millisecond || aguiPage != 100 || cap(d.aguiSlots) != 16 {
		t.Fatal("AGUI budget/default derivation violated")
	}
}
func TestAGUIReadPostureMatchesLog(t *testing.T) {
	d := aguiTestDaemon(t)
	aguiFixture(t, d)
	if d.isProtected(httptest.NewRequest("POST", "/agui/", nil)) != d.isProtected(httptest.NewRequest("GET", "/v1/log", nil)) || d.isProtected(httptest.NewRequest("POST", "/agui/", nil)) {
		t.Fatal("read posture differs")
	}
	live := authHeader(t, d)
	expired, _, _, err := authority.Mint(boundedTestContext(t), d.store, "expired", grantForAuth(), 1, time.Now().Unix()-10, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, auth := range []string{"", "Bearer invalid", live, "Bearer " + expired} {
		log := requestRecorderAuth(t, d, auth, "GET", "/v1/log", nil)
		rec := requestRecorderAuth(t, d, auth, "POST", "/agui/", strings.NewReader(aguiPayload))
		if log.Code != 200 || rec.Code != 200 {
			t.Fatalf("read auth posture %q: %d %d", auth, log.Code, rec.Code)
		}
		requireIndexes(t, rec.Body.Bytes(), []int64{0, 1, 2, 5, 6})
	}
}
func TestAGUIRejectsGET(t *testing.T) {
	d := aguiTestDaemon(t)
	if rec := requestRecorder(t, d, "GET", "/agui/", nil); rec.Code != 405 {
		t.Fatalf("GET=%d want405", rec.Code)
	}
}
func TestAGUIInputErrors(t *testing.T) {
	d := aguiTestDaemon(t)
	emptyIDs := aguiRequest(t, d, `{"threadId":"","runId":"","messages":[]}`, "")
	if emptyIDs.Code != 200 {
		t.Fatalf("opaque string IDs may be empty: %d", emptyIDs.Code)
	}
	for _, body := range []string{`{`, aguiPayload + `{}`, `{"threadId":"t","messages":[]}`, `{"runId":"r","messages":[]}`, `{"threadId":2,"runId":"r","messages":[]}`, `{"threadId":"t","runId":false,"messages":[]}`, `{"threadId":"t","runId":"r"}`, `{"threadId":"t","runId":"r","messages":null}`, `{"threadId":"t","runId":"r","messages":{}}`} {
		r := aguiRequest(t, d, body, "")
		if r.Code != 400 || !json.Valid(r.Body.Bytes()) || bytes.Contains(r.Body.Bytes(), []byte("data:")) {
			t.Fatalf("bad input %s: %d %s", body, r.Code, r.Body)
		}
	}
	r := aguiRequest(t, d, `{"threadId":"`+strings.Repeat("x", 65536)+`","runId":"r","messages":[]}`, "")
	if r.Code != 413 || !strings.Contains(r.Body.String(), "PayloadTooLarge") {
		t.Fatalf("large body: %d %s", r.Code, r.Body)
	}
	r = aguiRequest(t, d, `{"threadId":"t","runId":"r","messages":[null,7],"tools":{},"context":false,"forwardedProps":5}`, "")
	if r.Code != 200 {
		t.Fatalf("optional values: %d %s", r.Code, r.Body)
	}
	for k, v := range map[string]string{"Content-Type": "text/event-stream", "Cache-Control": "no-store", "X-Content-Type-Options": "nosniff"} {
		if r.Header().Get(k) != v {
			t.Fatalf("header %s=%s", k, r.Header().Get(k))
		}
	}
	if r.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("CORS added")
	}
}

// Blocking reads expose cancellation rather than merely observing handler exit.
type aguiBlockReads struct {
	readStore
	entered  chan context.Context
	release  chan struct{}
	canceled chan struct{}
}

func (b *aguiBlockReads) LogEntriesAfter(ctx context.Context, after int64, limit int) ([]store.LogEntry, error) {
	select {
	case b.entered <- ctx:
	default:
	}
	select {
	case <-ctx.Done():
		select {
		case b.canceled <- struct{}{}:
		default:
		}
		return nil, ctx.Err()
	case <-b.release:
		return nil, nil
	}
}
func aguiBlock(t *testing.T, d *Daemon) *aguiBlockReads {
	t.Helper()
	b := &aguiBlockReads{readStore: d.reads, entered: make(chan context.Context, 1), release: make(chan struct{}), canceled: make(chan struct{}, 1)}
	d.reads = b
	t.Cleanup(func() { close(b.release) })
	return b
}
func aguiWait(t *testing.T, ch <-chan struct{}, what string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(150 * time.Millisecond):
		t.Fatalf("%s did not finish promptly", what)
	}
}
func TestAGUIClientDisconnectReturns(t *testing.T) {
	d := aguiTestDaemon(t)
	d.aguiBudget = time.Second
	b := aguiBlock(t, d)
	base := runtime.NumGoroutine()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := httptest.NewRequest("POST", "/agui/", strings.NewReader(aguiPayload)).WithContext(ctx)
	done := make(chan struct{})
	go func() { d.Handler().ServeHTTP(httptest.NewRecorder(), req); close(done) }()
	select {
	case <-b.entered:
	case <-time.After(time.Second):
		t.Fatal("read never entered")
	}
	cancel()
	aguiWait(t, b.canceled, "read cancellation")
	aguiWait(t, done, "handler")
	deadline := time.Now().Add(150 * time.Millisecond)
	for runtime.NumGoroutine() > base && time.Now().Before(deadline) {
		runtime.Gosched()
	}
	if runtime.NumGoroutine() > base {
		t.Fatalf("goroutines baseline=%d got=%d", base, runtime.NumGoroutine())
	}
}
func TestAGUIShutdownEndsRuns(t *testing.T) {
	for _, blocked := range []bool{false, true} {
		t.Run(fmt.Sprint(blocked), func(t *testing.T) {
			d := aguiTestDaemon(t)
			d.aguiBudget = time.Second
			var b *aguiBlockReads
			if blocked {
				b = aguiBlock(t, d)
			}
			if err := d.Listen(); err != nil {
				t.Fatal(err)
			}
			serveDone := make(chan struct{})
			go func() { d.Serve(); close(serveDone) }()
			t.Cleanup(func() { d.srv.Close(); <-serveDone })
			req, _ := http.NewRequest("POST", d.URL()+"/agui/", strings.NewReader(aguiPayload))
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			if blocked {
				select {
				case <-b.entered:
				case <-time.After(time.Second):
					t.Fatal("read never entered")
				}
			}
			start := time.Now()
			done := make(chan error, 1)
			go func() { done <- d.Shutdown() }()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if time.Since(start) > 150*time.Millisecond || !bytes.Contains(body, []byte(`"RUN_FINISHED"`)) {
				t.Fatalf("shutdown terminal late/missing: %s", body)
			}
			select {
			case err := <-done:
				if err != nil {
					t.Fatal(err)
				}
			case <-time.After(2 * time.Second):
				t.Fatal("shutdown held drain")
			}
			if blocked {
				aguiWait(t, b.canceled, "shutdown read cancellation")
			}
		})
	}
}

type aguiFailReads struct {
	readStore
	err error
}

func (f aguiFailReads) LogEntriesAfter(context.Context, int64, int) ([]store.LogEntry, error) {
	return nil, f.err
}
func (f aguiFailReads) GetLogEntry(context.Context, int64) (store.LogEntry, bool, error) {
	return store.LogEntry{}, false, f.err
}
func TestAGUIStoreErrorIsRunError(t *testing.T) {
	for _, tc := range []struct {
		err  error
		code string
	}{{errors.New("SECRET-store-path"), "Internal"}, {context.DeadlineExceeded, "Timeout"}} {
		d := aguiTestDaemon(t)
		var log bytes.Buffer
		d.errLog = &log
		d.reads = aguiFailReads{d.reads, tc.err}
		r := aguiRequest(t, d, aguiPayload, "")
		events := aguiEvents(t, r.Body.Bytes())
		if len(events) != 3 {
			t.Fatalf("terminal count: %s", r.Body)
		}
		e := events[2]
		if string(e["type"]) != `"RUN_ERROR"` || string(e["code"]) != strconv.Quote(tc.code) || string(e["message"]) != strconv.Quote(internalErrorMessage) || e["threadId"] != nil || bytes.Contains(r.Body.Bytes(), []byte("SECRET")) {
			t.Fatalf("unsanitized/wrong error: %s", r.Body)
		}
		if !strings.Contains(log.String(), tc.err.Error()) {
			t.Fatal("detail absent from internal log")
		}
		pre := aguiRequest(t, d, aguiPayload, "1")
		if pre.Code != 500 || bytes.Contains(pre.Body.Bytes(), []byte("data:")) {
			t.Fatalf("prestream error: %d %s", pre.Code, pre.Body)
		}
	}
}

type aguiPollReads struct {
	readStore
	polls chan struct{}
}

func (p aguiPollReads) LogEntriesAfter(ctx context.Context, a int64, n int) ([]store.LogEntry, error) {
	rows, err := p.readStore.LogEntriesAfter(ctx, a, n)
	select {
	case p.polls <- struct{}{}:
	default:
	}
	return rows, err
}
func TestAGUISeesDirectStoreCommit(t *testing.T) {
	d := aguiTestDaemon(t)
	d.aguiBudget = time.Second
	d.aguiTick = 20 * time.Millisecond
	p := aguiPollReads{d.reads, make(chan struct{}, 100)}
	d.reads = p
	s := httptest.NewServer(d.Handler())
	defer s.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "POST", s.URL+"/agui/", strings.NewReader(aguiPayload))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	lines := make(chan string, 100)
	done := make(chan struct{})
	go func() {
		defer close(done)
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
	}()
	defer func() { cancel(); resp.Body.Close(); <-done }()
	aguiWait(t, p.polls, "initial poll")
	c := genesisCommit("direct")
	start := time.Now()
	if err := d.store.Commit(boundedTestContext(t), c); err != nil {
		t.Fatal(err)
	}
	waitEntry := func(i int64) {
		t.Helper()
		deadline := time.NewTimer(2*d.aguiTick + 50*time.Millisecond)
		defer deadline.Stop()
		for {
			select {
			case line := <-lines:
				if strings.Contains(line, fmt.Sprintf(`"entryIndex":%d`, i)) {
					return
				}
			case <-deadline.C:
				t.Fatalf("entry%d not polled within two ticks", i)
			}
		}
	}
	waitEntry(0)
	if time.Since(start) > 2*d.aguiTick+50*time.Millisecond {
		t.Fatal("direct commit late")
	}
	auth := authHeader(t, d)
	c = testCommit(c.NextWorld, 1, "rest-live")
	postCommit(t, s.URL, c, auth)
	waitEntry(1)
}
func TestAGUIGlobalCap(t *testing.T) {
	d := aguiTestDaemon(t)
	d.aguiBudget = time.Second
	s := httptest.NewServer(d.Handler())
	defer s.Close()
	type held struct {
		cancel context.CancelFunc
		body   io.ReadCloser
	}
	runs := []held{}
	open := func() held {
		t.Helper()
		ctx, cancel := context.WithCancel(context.Background())
		req, _ := http.NewRequestWithContext(ctx, "POST", s.URL+"/agui/", strings.NewReader(aguiPayload))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			cancel()
			t.Fatal(err)
		}
		if resp.StatusCode != 200 {
			cancel()
			resp.Body.Close()
			t.Fatalf("admission=%d", resp.StatusCode)
		}
		return held{cancel, resp.Body}
	}
	defer func() {
		for _, r := range runs {
			r.cancel()
			r.body.Close()
		}
	}()
	for i := 0; i < 16; i++ {
		runs = append(runs, open())
	}
	// Recorder avoids a blocked 17th request when the cap mutant is applied.
	req := httptest.NewRequest("POST", "/agui/", strings.NewReader(aguiPayload))
	ctx, cancel := context.WithTimeout(req.Context(), 50*time.Millisecond)
	defer cancel()
	rec := httptest.NewRecorder()
	d.Handler().ServeHTTP(rec, req.WithContext(ctx))
	if rec.Code != 503 || rec.Header().Get("Retry-After") != "1" || !strings.Contains(rec.Body.String(), "StreamLimit") || strings.Contains(rec.Body.String(), "data:") {
		t.Fatalf("17th: %d %s", rec.Code, rec.Body)
	}
	runs[0].cancel()
	runs[0].body.Close()
	deadline := time.Now().Add(150 * time.Millisecond)
	for len(d.aguiSlots) == 16 && time.Now().Before(deadline) {
		runtime.Gosched()
	}
	if len(d.aguiSlots) == 16 {
		t.Fatal("disconnect slot not released")
	}
	runs = append(runs, open())
	// Refusals and store errors must also release their slot.
	for _, r := range runs {
		r.cancel()
		r.body.Close()
	}
	deadline = time.Now().Add(150 * time.Millisecond)
	for len(d.aguiSlots) > 0 && time.Now().Before(deadline) {
		runtime.Gosched()
	}
	if len(d.aguiSlots) != 0 {
		t.Fatal("slots stranded")
	}
	aguiRequest(t, d, aguiPayload, "999")
	d.reads = aguiFailReads{d.reads, errors.New("failure")}
	aguiRequest(t, d, aguiPayload, "")
	if len(d.aguiSlots) != 0 {
		t.Fatal("refusal/error slot stranded")
	}
}

type aguiDelayedBody struct {
	data  []byte
	delay time.Duration
	once  bool
}

func (b *aguiDelayedBody) Read(p []byte) (int, error) {
	if !b.once {
		time.Sleep(b.delay)
		b.once = true
	}
	if len(b.data) == 0 {
		return 0, io.EOF
	}
	n := copy(p, b.data)
	b.data = b.data[n:]
	return n, nil
}
func (b *aguiDelayedBody) Close() error { return nil }

type aguiDripBody struct{ n int }

func (b *aguiDripBody) Read(p []byte) (int, error) {
	time.Sleep(30 * time.Millisecond)
	if b.n == 0 {
		return 0, io.EOF
	}
	b.n--
	p[0] = ' '
	return 1, nil
}
func (b *aguiDripBody) Close() error { return nil }
func TestAGUISlowBodyNeverTruncates(t *testing.T) {
	d := aguiTestDaemon(t)
	d.aguiBodyBound = 100 * time.Millisecond
	d.aguiBudget = 300 * time.Millisecond
	d.aguiTick = 10 * time.Millisecond
	req := httptest.NewRequest("POST", "/agui/", nil)
	req.Body = &aguiDripBody{10}
	r := httptest.NewRecorder()
	d.Handler().ServeHTTP(r, req)
	if r.Code != 408 || !json.Valid(r.Body.Bytes()) || !strings.Contains(r.Body.String(), "SlowBody") || strings.Contains(r.Body.String(), "data:") {
		t.Fatalf("slow body: %d %s", r.Code, r.Body)
	}
	req = httptest.NewRequest("POST", "/agui/", nil)
	req.Body = &aguiDelayedBody{data: []byte(aguiPayload), delay: 80 * time.Millisecond}
	r = httptest.NewRecorder()
	start := time.Now()
	d.Handler().ServeHTTP(r, req)
	if elapsed := time.Since(start); elapsed > 300*time.Millisecond+d.aguiTick+50*time.Millisecond || !strings.Contains(r.Body.String(), `"RUN_FINISHED"`) {
		t.Fatalf("entry anchored terminal elapsed=%s body=%s", elapsed, r.Body)
	}
	// A real stalled body exhausts equal read/write deadlines, yielding clean EOF.
	srv := httptest.NewUnstartedServer(d.Handler())
	srv.Config.ReadTimeout = 150 * time.Millisecond
	srv.Config.WriteTimeout = 150 * time.Millisecond
	srv.Start()
	defer srv.Close()
	conn, err := net.DialTimeout("tcp", srv.Listener.Addr().String(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Second))
	fmt.Fprintf(conn, "POST /agui/ HTTP/1.1\r\nHost: localhost\r\nContent-Length: 500\r\n\r\n")
	b, err := io.ReadAll(conn)
	if err != nil || len(b) != 0 {
		t.Fatalf("stalled body must be zero bytes/EOF: %q %v", b, err)
	}
}

// The id oracle is independent of state: an id acknowledges its own complete delta.
func aguiCompleteCursor(t *testing.T, b []byte) int64 {
	t.Helper()
	cursor := int64(-1)
	id := int64(-1)
	frames := bytes.Split(b, []byte("\n\n"))
	for _, f := range frames[:len(frames)-1] {
		var ev map[string]json.RawMessage
		frameID := ""
		for _, l := range bytes.Split(f, []byte("\n")) {
			if bytes.HasPrefix(l, []byte("id: ")) {
				frameID = string(l[4:])
			}
			if bytes.HasPrefix(l, []byte("data: ")) {
				if err := json.Unmarshal(l[6:], &ev); err != nil {
					t.Fatal(err)
				}
			}
		}
		if string(ev["type"]) == `"STATE_SNAPSHOT"` {
			var s struct {
				LastIndex int64 `json:"lastIndex"`
			}
			json.Unmarshal(ev["snapshot"], &s)
			cursor = s.LastIndex
			id = cursor
		}
		if string(ev["type"]) == `"STATE_DELTA"` {
			var ops []struct {
				Path  string
				Value json.RawMessage
			}
			json.Unmarshal(ev["delta"], &ops)
			for _, o := range ops {
				if o.Path == "/lastIndex" {
					json.Unmarshal(o.Value, &cursor)
				}
			}
			if frameID == "" {
				t.Fatal("delta lacks final-frame id")
			}
		}
		if frameID != "" {
			if string(ev["type"]) != `"STATE_DELTA"` {
				t.Fatal("id acknowledges incomplete entry before delta")
			}
			n, err := strconv.ParseInt(frameID, 10, 64)
			if err != nil || n != cursor {
				t.Fatal("id/delta mismatch")
			}
			id = n
		}
		if id != cursor {
			t.Fatal("last complete id differs from delta cursor")
		}
	}
	return cursor
}
func TestAGUISeveranceResumable(t *testing.T) {
	d := aguiTestDaemon(t)
	aguiFixture(t, d)
	full := aguiRequest(t, d, aguiPayload, "").Body.Bytes()
	requireIndexes(t, full, []int64{0, 1, 2, 5, 6})
	for cut := 0; cut <= len(full); cut++ {
		cursor := aguiCompleteCursor(t, full[:cut])
		seen := []int64{}
		for _, i := range []int64{0, 1, 2, 5, 6} {
			if i <= cursor {
				seen = append(seen, i)
			}
		}
		for _, e := range aguiEvents(t, full) {
			if string(e["type"]) != `"CUSTOM"` {
				continue
			}
			var v struct {
				Header struct {
					EntryIndex int64 `json:"entryIndex"`
				} `json:"header"`
			}
			json.Unmarshal(e["value"], &v)
			if v.Header.EntryIndex > cursor {
				seen = append(seen, v.Header.EntryIndex)
			}
		}
		if !reflect.DeepEqual(seen, []int64{0, 1, 2, 5, 6}) {
			t.Fatalf("cut%d loses entry: %v", cut, seen)
		}
	}
	// Both cursor paths are exercised over the actual handler; enumeration stays pure.
	for _, k := range []int64{-1, 0, 1, 2, 5, 6} {
		want := []int64{}
		for _, i := range []int64{0, 1, 2, 5, 6} {
			if i > k {
				want = append(want, i)
			}
		}
		requireIndexes(t, aguiRequest(t, d, aguiPayload, fmt.Sprint(k)).Body.Bytes(), want)
	}
}

// End replay as soon as its final delta arrives, rather than burning an idle run.
type aguiReplayRecorder struct {
	*httptest.ResponseRecorder
	cancel context.CancelFunc
}

func (r aguiReplayRecorder) Write(b []byte) (int, error) {
	n, err := r.ResponseRecorder.Write(b)
	if bytes.Contains(b, []byte("id: 4999\n")) {
		r.cancel()
	}
	return n, err
}

func TestAGUISlowReaderCutIsResumable(t *testing.T) {
	d := aguiTestDaemon(t)
	d.aguiBudget = 2 * time.Second
	w := store.World{LogHead: hashref.SumSHA256([]byte("initial"))}
	for i := int64(0); i < 5000; i++ {
		c := testCommit(w, i, fmt.Sprint(i))
		c.Objects = nil
		c.Entry.Header.WrittenBy = strings.Repeat("large-entry", 400)
		if err := d.store.Commit(boundedTestContext(t), c); err != nil {
			t.Fatal(err)
		}
		w = c.NextWorld
	}
	srv := httptest.NewUnstartedServer(d.Handler())
	srv.Config.WriteTimeout = 200 * time.Millisecond
	srv.Start()
	defer srv.Close()
	conn, err := net.DialTimeout("tcp", srv.Listener.Addr().String(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	fmt.Fprintf(conn, "POST /agui/ HTTP/1.1\r\nHost: localhost\r\nContent-Length: %d\r\nContent-Type: application/json\r\n\r\n%s", len(aguiPayload), aguiPayload)
	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	prefix := make([]byte, 1024)
	n, err := io.ReadFull(resp.Body, prefix)
	if err != nil {
		t.Fatal(err)
	}
	prefix = prefix[:n]
	time.Sleep(250 * time.Millisecond)
	rest, err := io.ReadAll(resp.Body)
	prefix = append(prefix, rest...)
	if err == nil || bytes.Contains(prefix, []byte(`"RUN_FINISHED"`)) || bytes.Contains(prefix, []byte(`"RUN_ERROR"`)) {
		t.Fatalf("expected write-timeout cut, got err=%v bytes=%d", err, len(prefix))
	}
	if ne, ok := err.(net.Error); ok && ne.Timeout() {
		t.Fatal("client deadline, not server severance")
	}
	cursor := aguiCompleteCursor(t, prefix)
	if cursor < 0 || cursor >= 4999 {
		t.Fatalf("cut cursor=%d", cursor)
	}
	// Recorder consumes the remaining large response immediately, with no socket pause.
	for _, body := range []string{aguiPayload, fmt.Sprintf(`{"threadId":"t","runId":"r","messages":[],"state":{"schema":"world/agui-state/v1","lastIndex":%d}}`, cursor)} {
		header := ""
		if body == aguiPayload {
			header = fmt.Sprint(cursor)
		}
		ctx, cancel := context.WithCancel(context.Background())
		req := httptest.NewRequest("POST", "/agui/", strings.NewReader(body)).WithContext(ctx)
		req.Header.Set("Last-Event-ID", header)
		r := httptest.NewRecorder()
		d.Handler().ServeHTTP(aguiReplayRecorder{r, cancel}, req)
		cancel()
		want := []int64{}
		for i := cursor + 1; i < 5000; i++ {
			want = append(want, i)
		}
		requireIndexes(t, r.Body.Bytes(), want)
		aguiCompleteCursor(t, r.Body.Bytes())
	}
}
