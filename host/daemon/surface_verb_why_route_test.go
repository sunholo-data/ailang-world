package daemon

// Row 243 V2b (design §2.5, §2.7, §6 VR1–VR4, VR6, VR7): why in the bar, end to end through the route.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/why"
)

// whyRig is verbRig plus a committed REST call (a resolved rest: receipt) and a REST intent that never commits.
func whyRig(t *testing.T) (approvalRig, *reportSeed) {
	t.Helper()
	a, s := verbRig(t)
	s.call("rest", "w1", "3", "ailang-read", read(50, map[string]any{"ok": true, "content": "y"}), nil)
	s.pending("rest", "w1", "9", "ailang-read")
	return a, s
}

// whyCLI is the CLI line Lookup gives a why line.
func whyCLI(q string) string { return "ailang-worldd " + q }

// VR1: every outcome of design §2.7 end to end through approvalRoute; a result carries no reason
// (Lookup's V1 "cannot run it yet" text must not ride on an answer); a refusal carries the CLI line.
func TestSurfaceVerbWhyOutcomes(t *testing.T) {
	a, s := whyRig(t)
	zero := "sha256:" + strings.Repeat("0", 64)
	for _, c := range []struct {
		q, member, prefix, suffix string
		exit                      int
	}{
		{"why 1", "text", "why: entry 1 (matched by entry index)\n  ✗ world ", " link(s) BROKEN — integrity refusal (exit 3)\n", 3},
		{"why head --json", "json", "{\n  \"entry\": 3,\n  \"matchedBy\": \"head\",\n", "\n}\n", 3},
		{"why rest:w1:3", "text", "why: entry 3 (matched by rest receipt → world)\n", "(exit 3)\n", 3},
		{"why " + s.head.Ref.String(), "text", "why: entry 3 (matched by world ref)\n", "(exit 3)\n", 3},
		{"why 99", "stderr", "ailang-worldd why: log entry 99 not found\n", "", 1},
		{"why a2a:nobody:x --scan 1", "stderr", "ailang-worldd why: not found in the last 1 entries (invocationId)\n", "", 1},
		{"why rest:w1:9", "stderr", "ailang-worldd why: receipt rest:w1:9 is indeterminate: nothing committed\n", "", 1},
		{"why rest:nobody", "stderr", "ailang-worldd why: receipt rest:nobody is not-started: nothing committed\n", "", 1},
		{"why " + zero, "stderr", "ailang-worldd why: " + zero + " is neither an object nor a world in this store\n", "", 1},
	} {
		v, raw := verbGet(a, c.q)
		got := map[string]*string{"text": v.Text, "json": v.JSON, "stderr": v.Stderr}
		if v.Outcome != "result" || v.Class != "read" || v.Verb != "why" || v.CLI != whyCLI(c.q) || v.Reason != "" || v.Exit != c.exit ||
			got[c.member] == nil || !strings.HasPrefix(*got[c.member], c.prefix) || !strings.HasSuffix(*got[c.member], c.suffix) {
			t.Errorf("%q = %s", c.q, raw)
		}
		if c.member == "stderr" && *got[c.member] != c.prefix {
			t.Errorf("%q stderr = %q", c.q, *got[c.member])
		}
		for k, p := range got {
			if k != c.member && p != nil {
				t.Errorf("%q carries %s too: %s", c.q, k, raw)
			}
		}
	}
	for _, c := range []struct{ q, reason string }{
		{"why", "ailang-worldd why: a target is required"},
		{"why x", `ailang-worldd why: "x" is not an entry index, head, sha256:<ref>, a2a:/mcp:/rest:<id>, or -`},
		{"why sha256:zz", `ailang-worldd why: "sha256:zz" is not a content address: hashref: digest is not lowercase hexadecimal: "zz"`},
		{"why --scan 0 head", "ailang-worldd why: --scan must be 1..5000"},
		{"why 1 2", `ailang-worldd why: unexpected argument "2"`},
		{"why --bogus", "flag provided but not defined: -bogus"},
		{"why --help", "see ailang-worldd why --help at a terminal"},
		{"why -", "why - is not run from the app: the app has no stdin; at a terminal: ailang-worldd why -"},
		{"why --result f", "--result is not run from the app: the app reads no file on the daemon's host; at a terminal: ailang-worldd why --result f"},
		{"why --expect-store p head", "--expect-store is not run from the app: the app reads its own daemon's store; at a terminal: ailang-worldd why --expect-store p head"},
	} {
		if v, raw := verbGet(a, c.q); v.Outcome != "refused" || v.Reason != c.reason || v.CLI != whyCLI(c.q) || v.Text != nil || v.JSON != nil || v.Stderr != nil || v.Exit != 0 {
			t.Errorf("%q = %s, want refused %q", c.q, raw, c.reason)
		}
	}
	if v, _ := verbGet(a, "world why"); v.Outcome != "address" { // RK1: brane why still navigates
		t.Errorf("world why = %+v", *v)
	}
}

// VR2: an exhausted read deadline is (b), never a result, never "not found" — whether the walk's own
// context expires before any read (the scan) or a handler answers 503 while the walk's context lives.
func TestSurfaceVerbWhyDeadline(t *testing.T) {
	a, s := whyRig(t)
	s.d.errLog = io.Discard
	s.d.readDeadline = expiredReadDeadline
	for _, q := range []string{"why 1", "why a2a:nobody:x", "why head --json"} {
		if v, raw := verbGet(a, q); v.Outcome != "refused" || v.Reason != "read deadline (-1ns) exceeded: the walk did not finish in the app; at a terminal it has 60 s: "+whyCLI(q) || strings.Contains(raw, "not found") {
			t.Errorf("expired: %q = %s", q, raw)
		}
	}
	base := s.d.reads
	count := &whyFailReads{readStore: base}
	s.d.reads = count
	if v, _ := verbGet(a, "why sha256:"+strings.Repeat("0", 64)); v.Outcome != "refused" || count.n != 0 { // the walk's own deadline: no read starts
		t.Errorf("expired: %d object reads, %+v", count.n, *v)
	}
	s.d.readDeadline = readDeadline
	for _, reads := range []readStore{&whyFailReads{readStore: base, nth: 1, err: context.DeadlineExceeded}, // a 503 mid-walk
		&whyFailReads{readStore: &whyFailReads{readStore: base, nth: 2, err: errors.New("x")}, nth: 3, err: context.DeadlineExceeded}} { // a 500, then a 503: (b) first
		s.d.reads = reads
		if v, raw := verbGet(a, "why 1"); v.Outcome != "refused" || !strings.HasPrefix(v.Reason, "read deadline (10s) exceeded: the walk did not finish in the app") {
			t.Errorf("a 503 mid-walk = %s", raw)
		}
	}
	// (a): the page went away; the answer says so, never a deadline or a result (nobody reads it).
	s.d.reads = base
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), http.LocalAddrContextKey, net.Addr(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 7644})))
	cancel()
	req := httptest.NewRequest("GET", "/surface/verb.json?q=why+1", nil).WithContext(ctx)
	req.Host = apLoopback
	rec := httptest.NewRecorder()
	s.d.Handler().ServeHTTP(rec, req)
	if !strings.Contains(rec.Body.String(), `"reason":"the request ended before the walk did"`) {
		t.Errorf("a cancelled request = %d %s", rec.Code, rec.Body)
	}
}

// whyFailReads fails the nth GetObject (counting from 1) with err; n counts every call.
type whyFailReads struct {
	readStore
	nth, n int
	err    error
}

func (f *whyFailReads) GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error) {
	if f.n++; f.n == f.nth {
		return store.Object{}, false, f.err
	}
	return f.readStore.GetObject(ctx, ref)
}

// VR3: a read failure anywhere in the walk — every GetObject in turn, then a log read — is (c): no
// chain, no ✗, no "not found", no store error text; the operator's log has it.
func TestSurfaceVerbWhyReadFailure(t *testing.T) {
	a, s := whyRig(t)
	base := s.d.reads
	count := &whyFailReads{readStore: base}
	s.d.reads = count
	verbGet(a, "why 1")
	if count.n < 5 {
		t.Fatalf("why 1 read %d objects", count.n)
	}
	const want = "the daemon could not read World during the walk, so no chain is shown; at a terminal: ailang-worldd why 1"
	for nth := 1; nth <= count.n; nth++ {
		var log bytes.Buffer
		s.d.errLog = &log
		s.d.reads = &whyFailReads{readStore: base, nth: nth, err: errors.New("SECRET-store-path")}
		if v, raw := verbGet(a, "why 1"); v.Outcome != "refused" || v.Reason != want || strings.Contains(raw, "SECRET") || strings.Contains(raw, "✗") {
			t.Errorf("GetObject #%d fails: %s", nth, raw)
		}
		if !strings.Contains(log.String(), "SECRET-store-path") {
			t.Errorf("GetObject #%d: not logged: %q", nth, log.String())
		}
	}
	s.d.errLog = io.Discard
	s.d.reads = aguiFailReads{base, errors.New("SECRET-store-path")}
	if v, raw := verbGet(a, "why 1"); v.Outcome != "refused" || v.Reason != want || strings.Contains(raw, "not found") {
		t.Errorf("log read fails: %s", raw)
	}
}

// recordingGetter is why's Getter over HTTP to url, keeping every path's (status, body).
func recordingGetter(t *testing.T, url string, rec map[string][2]string) why.Getter {
	return func(ctx context.Context, path string) (int, []byte, error) {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url+path, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return 0, nil, err
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		rec[path] = [2]string{strconv.Itoa(resp.StatusCode), string(b)}
		return resp.StatusCode, b, err
	}
}

// VR4, VR6: every path the walk reads over HTTP (every target kind) answers byte for byte the same
// in-process, /v1/health on the members an answer embeds; the walks' answers are equal; the private
// mux's patterns are exactly the ones those paths reach, each bound to the method Handler registers.
func TestWhyGetterTransportParity(t *testing.T) {
	_, s := whyRig(t)
	srv := httptest.NewServer(s.d.Handler())
	defer srv.Close()
	rec := map[string][2]string{}
	ctx := boundedTestContext(t)
	for _, c := range []struct {
		target string
		o      why.Options
	}{{"1", why.Options{Scan: 500}}, {"head", why.Options{Scan: 500, JSON: true}}, {"rest:w1:3", why.Options{Scan: 500}},
		{"rest:nobody", why.Options{Scan: 500}}, {s.head.Ref.String(), why.Options{Scan: 500}}, {"a2a:nobody:x", why.Options{Scan: 500}}, {"99", why.Options{Scan: 500}}} {
		tg, err := why.ParseTarget(c.target)
		if err != nil {
			t.Fatal(err)
		}
		var tally whyTally
		overHTTP := why.Run(why.Reader{Ctx: ctx, Get: recordingGetter(t, srv.URL, rec)}, tg, c.o)
		inProc := why.Run(why.Reader{Ctx: ctx, Get: s.d.whyGetter(httptest.NewRequest("GET", "/surface/verb.json", nil), s.d.whyMux(), &tally)}, tg, c.o)
		if overHTTP != inProc || tally != (whyTally{}) {
			t.Errorf("%s: in-process %+v %+v, over HTTP %+v", c.target, inProc, tally, overHTTP)
		}
	}
	mux, reached := s.d.whyMux(), map[string]bool{}
	for path, want := range rec {
		var tally whyTally
		status, body, err := s.d.whyGetter(httptest.NewRequest("GET", "/surface/verb.json", nil), mux, &tally)(ctx, path)
		got := [2]string{strconv.Itoa(status), string(body)}
		if strings.HasPrefix(path, "/v1/health") {
			got[1], want[1] = healthPin(t, got[1]), healthPin(t, want[1])
		}
		if err != nil || got != want {
			t.Errorf("%s: in-process %v %v, over HTTP %v", path, got, err, want)
		}
		req, _ := http.NewRequest("GET", path, nil)
		_, pattern := mux.Handler(req)
		reached[pattern] = true
	}
	var patterns []string
	for p := range s.d.whyRoutes() {
		patterns = append(patterns, p)
	}
	if want := mapKeys(reached); !reflect.DeepEqual(whySorted(patterns), want) {
		t.Errorf("whyRoutes %v, but the recorded walks reach %v", whySorted(patterns), want)
	}
}

func healthPin(t *testing.T, body string) string {
	var h struct {
		Ref     string `json:"interpreter_ref"`
		Version string `json:"interpreter_version"`
	}
	if err := json.Unmarshal([]byte(body), &h); err != nil {
		t.Fatalf("health %q: %v", body, err)
	}
	return h.Ref + "|" + strings.SplitN(h.Version, "\n", 2)[0]
}

func mapKeys(m map[string]bool) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	return whySorted(out)
}

func whySorted(s []string) []string { sort.Strings(s); return s }

// VR7 (end to end): the cap runs in the route's deferred writer, for why (text and json) and for
// episode report alike; the refusal carries the size and the CLI line and none of the members.
func TestVerbAnswerCapEndToEnd(t *testing.T) {
	a, _ := whyRig(t)
	sizes := map[string]int{}
	for _, q := range []string{"why 1", "why head --json", "episode report w1", "why 99"} {
		v, _ := verbGet(a, q)
		for _, m := range []*string{v.Text, v.JSON, v.Stderr} {
			if m != nil {
				sizes[q] += len(*m)
			}
		}
	}
	verbAnswerCap = sizes["why 99"]
	t.Cleanup(func() { verbAnswerCap = MaxVerbAnswerBytes })
	for _, q := range []string{"why 1", "why head --json", "episode report w1"} {
		v, raw := verbGet(a, q)
		if v.Outcome != "refused" || v.Reason != fmt.Sprintf("the answer is %d bytes, over the app's 1 MiB cap; at a terminal: ailang-worldd %s", sizes[q], q) ||
			strings.Contains(raw, `"text"`) || strings.Contains(raw, `"json"`) || strings.Contains(raw, `"stderr"`) || strings.Contains(raw, `"exit"`) {
			t.Errorf("%q over the cap = %s", q, raw)
		}
	}
	if v, raw := verbGet(a, "why 99"); v.Outcome != "result" || v.Stderr == nil { // at the cap exactly
		t.Errorf("at the cap = %s", raw)
	}
}
