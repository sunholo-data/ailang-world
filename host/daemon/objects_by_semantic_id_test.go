package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

func putSemanticObjects(t *testing.T, d *Daemon, id string, n int) []string {
	t.Helper()
	hashes := make([]string, 0, n)
	for i := 0; i < n; i++ {
		payload := []byte(fmt.Sprintf("route-sid-%s-%d", id, i))
		o := store.Object{
			Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte("route-sid-iface")),
			SemanticID: id, Provenance: "route-sid", Payload: payload,
		}
		if err := d.store.PutObject(o); err != nil {
			t.Fatalf("PutObject: %v", err)
		}
		hashes = append(hashes, o.Hash.String())
	}
	sort.Strings(hashes)
	return hashes
}

func decodePage(t *testing.T, rec interface{ Result() *http.Response }) objectPageResponse {
	t.Helper()
	var page objectPageResponse
	if err := json.NewDecoder(rec.Result().Body).Decode(&page); err != nil {
		t.Fatalf("decode page: %v", err)
	}
	return page
}

// TestObjectsBySemanticIDRoute pins AC-4 and AC-6: a match answers 200 with
// envelopes and no payload key, an unused name answers 200 with an empty list,
// and each refusal branch answers 400 without reaching the lookup.
func TestObjectsBySemanticIDRoute(t *testing.T) {
	d := newHandlerDaemon(t)
	want := putSemanticObjects(t, d, "world/mission/incident/route", 3)

	rec := requestRecorder(t, d, http.MethodGet, "/v1/objects/by-semantic-id/world/mission/incident/route", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("match status = %d body=%s", rec.Code, rec.Body)
	}
	var raw map[string][]map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	for _, item := range raw["items"] {
		if _, has := item["payload"]; has {
			t.Fatalf("item carries a payload key: %v", item)
		}
	}
	payloadRec := requestRecorder(t, d, http.MethodGet, "/v1/objects/by-semantic-id/world/mission/incident/route?payload=true", nil)
	var payloadRaw map[string][]map[string]any
	if err := json.Unmarshal(payloadRec.Body.Bytes(), &payloadRaw); err != nil {
		t.Fatal(err)
	}
	for _, item := range payloadRaw["items"] {
		if _, has := item["payload"]; has {
			t.Fatalf("payload=true returned payload: %v", item)
		}
	}
	postRec := requestRecorder(t, d, http.MethodPost, "/v1/objects/by-semantic-id/world/mission/incident/route", nil)
	if postRec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST lookup status=%d, want 405", postRec.Code)
	}
	page := decodePage(t, rec)
	if len(page.Items) != 3 {
		t.Fatalf("match page = %+v, want 3 items", page)
	}
	for i, item := range page.Items {
		if item.Hash != want[i] || item.SemanticID != "world/mission/incident/route" {
			t.Fatalf("item %d = %+v, want hash %s in ascending order", i, item, want[i])
		}
	}

	rec = requestRecorder(t, d, http.MethodGet, "/v1/objects/by-semantic-id/world/mission/incident/never", nil)
	if rec.Code != http.StatusOK || rec.Body.String() != "{\"items\":[]}\n" {
		t.Fatalf("unknown name = %d %q, want 200 {\"items\":[]}", rec.Code, rec.Body)
	}

	for _, tc := range []struct{ name, target string }{
		{"empty name", "/v1/objects/by-semantic-id/"},
		{"malformed after", "/v1/objects/by-semantic-id/world/mission/incident/route?after=not-a-ref"},
		{"non-integer limit", "/v1/objects/by-semantic-id/world/mission/incident/route?limit=ten"},
		{"no name segment", "/v1/objects/by-semantic-id"},
	} {
		rec := requestRecorder(t, d, http.MethodGet, tc.target, nil)
		assertErrorClass(t, rec, http.StatusBadRequest, "BadRequest")
	}

	rec = requestRecorder(t, d, http.MethodGet,
		"/v1/objects/by-semantic-id/world/mission/incident/route?after="+want[0], nil)
	if page := decodePage(t, rec); len(page.Items) != 2 || page.Items[0].Hash != want[1] {
		t.Fatalf("after=%s page = %+v, want the last two", want[0], page)
	}
}

// TestObjectsBySemanticIDRouteCapAndDefault pins the default and capped pages,
// then uses the last hash as the cursor. The envelope never has a next key.
func TestObjectsBySemanticIDRouteCapAndDefault(t *testing.T) {
	d := newHandlerDaemon(t)
	want := putSemanticObjects(t, d, "world/journal-intent/v1", 520)
	base := "/v1/objects/by-semantic-id/world/journal-intent/v1"
	check := func(target string, n int) objectPageResponse {
		t.Helper()
		rec := requestRecorder(t, d, http.MethodGet, target, nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("status=%d body=%s", rec.Code, rec.Body)
		}
		var raw map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &raw); err != nil {
			t.Fatal(err)
		}
		if _, has := raw["next"]; has {
			t.Fatalf("unexpected next key in %s", rec.Body)
		}
		page := decodePage(t, rec)
		if len(page.Items) != n {
			t.Fatalf("%s: %d items, want %d", target, len(page.Items), n)
		}
		return page
	}
	first := check(base, 100)
	if first.Items[99].Hash != want[99] {
		t.Fatalf("default page ends %s, want %s", first.Items[99].Hash, want[99])
	}
	big := check(base+"?limit=100000", 500)
	if big.Items[499].Hash != want[499] {
		t.Fatalf("capped page ends %s, want %s", big.Items[499].Hash, want[499])
	}
	last := check(base+"?limit=500&after="+big.Items[499].Hash, 20)
	got := make([]string, 0, 520)
	for _, item := range big.Items {
		got = append(got, item.Hash)
	}
	for _, item := range last.Items {
		got = append(got, item.Hash)
	}
	for i, h := range got {
		if h != want[i] {
			t.Fatalf("union[%d]=%s, want %s", i, h, want[i])
		}
	}
}

// walkIncidentID is the one semanticId the re-timed walk looks for.
const walkIncidentID = "world/mission/incident/walk-target"

// seedWalkLog commits n log entries shaped like row 92's: each entry's
// transitionRef is an object whose semanticId names what it is, and each commit
// also carries one journal object under the shared, non-unique id
// "world/journal-intent/v1". The LAST entry is the incident (worst case for a
// scan from 0); its payload lists three source objects committed with it.
func seedWalkLog(tb testing.TB, d *Daemon, n int) {
	tb.Helper()
	current := store.World{
		Ref: hashref.SumSHA256([]byte("walk-genesis")), Revision: 0,
		StateRoot: hashref.SumSHA256([]byte("walk-genesis-state")),
		LogHead:   hashref.SumSHA256([]byte("walk-genesis-log")),
	}
	if err := d.store.PutWorld(current); err != nil {
		tb.Fatal(err)
	}
	if err := d.store.SelectHead(current.Ref); err != nil {
		tb.Fatal(err)
	}
	obj := func(id, body string) store.Object {
		p := []byte(body)
		return store.Object{Hash: hashref.SumSHA256(p), InterfaceHash: hashref.SumSHA256([]byte("walk-iface/" + id)),
			SemanticID: id, Provenance: "walk-seed", Payload: p}
	}
	for i := int64(0); i < int64(n); i++ {
		transition := obj(fmt.Sprintf("world/mission/transition/%d", i), fmt.Sprintf("transition-%d", i))
		objects := []store.Object{transition, obj("world/journal-intent/v1", fmt.Sprintf("intent-%d", i))}
		if i == int64(n)-1 {
			var sources []string
			for s := 0; s < 3; s++ {
				src := obj("world/mission/source/"+strconv.Itoa(s), fmt.Sprintf("source-%d-evidence", s))
				objects = append(objects, src)
				sources = append(sources, src.Hash.String())
			}
			body, _ := json.Marshal(map[string]any{"answer": "why X happened", "sources": sources})
			transition = obj(walkIncidentID, string(body))
			objects[0] = transition
		}
		entryHash := hashref.SumSHA256([]byte(fmt.Sprintf("walk-entry-%d", i)))
		c := store.Commit{
			ObservedHead: current.Ref, Objects: objects,
			NextWorld: store.World{Ref: hashref.SumSHA256([]byte(fmt.Sprintf("walk-world-%d", i))), Revision: i + 1,
				StateRoot: hashref.SumSHA256([]byte(fmt.Sprintf("walk-state-%d", i))), LogHead: entryHash},
			Entry: store.LogEntry{
				Header: store.LogHeader{EntryIndex: i, SemanticsEpoch: 1,
					TransitionFn: hashref.SumSHA256([]byte("walk-fn")), Interpreter: hashref.SumSHA256([]byte("interpreter")),
					PrevEntryHash: current.LogHead, WrittenBy: "walk-seed"},
				EntryHash: entryHash, TransitionRef: transition.Hash,
			},
		}
		if err := d.store.Commit(c); err != nil {
			tb.Fatalf("seed Commit(%d): %v", i, err)
		}
		current = c.NextWorld
	}
}

// TestWalkRetimedAtScale is AC-11: row 92's walk, re-timed at N log entries,
// once with the linear scan (log range pages of 100 + one object get per entry
// until the semanticId matches) and once with the new route. Steps 2-3 (read
// the incident payload, then each source) are identical in both. It is a
// MEASUREMENT, gated on WORLD_WALK_N so the default suite stays fast; it
// asserts only that both walks reach the same incident and verified sources.
// WORLD_WALK_DB keeps the seeded store for the CLI re-timing.
func TestWalkRetimedAtScale(t *testing.T) {
	n, _ := strconv.Atoi(os.Getenv("WORLD_WALK_N"))
	if n <= 0 {
		t.Skip("set WORLD_WALK_N to run the walk re-timing")
	}
	runs := 5
	dbPath := os.Getenv("WORLD_WALK_DB")
	if dbPath == "" {
		dbPath = filepath.Join(t.TempDir(), "walk.db")
	}
	d, err := New(context.Background(), Config{DBPath: dbPath, BindHost: DefaultBindHost})
	if err != nil {
		t.Fatal(err)
	}
	seedStart := time.Now()
	seedWalkLog(t, d, n)
	t.Logf("seeded N=%d entries (%d objects) in %s", n, 2*n+3, time.Since(seedStart))
	if err := d.Listen(); err != nil {
		t.Fatal(err)
	}
	serveDone := make(chan error, 1)
	go func() { serveDone <- d.Serve() }()
	t.Cleanup(func() {
		_ = d.Shutdown()
		<-serveDone
		_ = d.Close()
	})
	client := &http.Client{Timeout: 30 * time.Second}
	calls := 0
	get := func(path string, into any) {
		calls++
		resp, err := client.Get(d.URL() + path)
		if err != nil {
			t.Fatalf("GET %s: %v", path, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("GET %s: %d %s", path, resp.StatusCode, body)
		}
		if err := json.NewDecoder(resp.Body).Decode(into); err != nil {
			t.Fatalf("decode %s: %v", path, err)
		}
	}
	readIncident := func(ref string) []string {
		var incident objectResponse
		get("/v1/objects/"+ref+"?payload=true", &incident)
		var body struct{ Sources []string }
		if err := json.Unmarshal(*incident.Payload, &body); err != nil || len(body.Sources) != 3 {
			t.Fatalf("incident payload %q: %v", *incident.Payload, err)
		}
		for _, src := range body.Sources {
			var o objectResponse
			get("/v1/objects/"+src+"?payload=true", &o)
			if got := hashref.SumSHA256(*o.Payload).String(); got != src {
				t.Fatalf("source %s payload hashes to %s", src, got)
			}
		}
		return body.Sources
	}
	scan := func() string {
		for from := 0; ; from += 100 {
			var page logRangeResponse
			get("/v1/log?from="+strconv.Itoa(from), &page)
			if len(page.Items) == 0 {
				t.Fatalf("scan reached the end without finding %s", walkIncidentID)
			}
			for _, e := range page.Items {
				var o objectResponse
				get("/v1/objects/"+e.TransitionRef, &o)
				if o.SemanticID == walkIncidentID {
					return o.Hash
				}
			}
		}
	}
	lookup := func() string {
		var page objectPageResponse
		get("/v1/objects/by-semantic-id/"+walkIncidentID, &page)
		if len(page.Items) != 1 {
			t.Fatalf("lookup = %+v, want exactly the incident", page)
		}
		return page.Items[0].Hash
	}
	measure := func(name string, locate func() string) []time.Duration {
		var samples []time.Duration
		for r := 0; r < runs; r++ {
			calls = 0
			start := time.Now()
			ref := locate()
			sources := readIncident(ref)
			samples = append(samples, time.Since(start))
			if r == 0 {
				t.Logf("%s: incident %s, %d sources verified, %d HTTP calls", name, ref, len(sources), calls)
			}
		}
		sort.Slice(samples, func(i, j int) bool { return samples[i] < samples[j] })
		t.Logf("WALK N=%d method=%s runs=%d p50=%s max=%s", n, name, runs, samples[runs/2], samples[runs-1])
		return samples
	}
	before := measure("scan", scan)
	after := measure("route", lookup)
	if scanRef, routeRef := scan(), lookup(); scanRef != routeRef {
		t.Fatalf("scan found %s but the route found %s", scanRef, routeRef)
	}
	_ = before
	_ = after
}

type semanticLookupFailStore struct {
	readStore
	err error
}

func (s semanticLookupFailStore) ObjectsBySemanticID(context.Context, string, string, int) ([]store.Object, error) {
	return nil, s.err
}

func TestObjectsBySemanticIDRouteStoreErrors(t *testing.T) {
	d := newHandlerDaemon(t)
	d.reads = semanticLookupFailStore{readStore: d.reads, err: errors.New("host-secret-detail")}
	rec := requestRecorder(t, d, http.MethodGet, "/v1/objects/by-semantic-id/x", nil)
	assertErrorClass(t, rec, http.StatusInternalServerError, "Internal")
	if strings.Contains(rec.Body.String(), "host-secret-detail") {
		t.Fatalf("internal detail leaked: %s", rec.Body)
	}
}

func TestLookupIndexUnavailableIs503(t *testing.T) {
	vectors := sketchHTTPStatusVectors(t)
	if vectors["LookupIndexUnavailable"] != 503 {
		t.Fatalf("sketch status=%d, want 503", vectors["LookupIndexUnavailable"])
	}
	d := newHandlerDaemon(t)
	d.reads = semanticLookupFailStore{readStore: d.reads, err: &store.LookupIndexUnavailableError{}}
	rec := requestRecorder(t, d, http.MethodGet, "/v1/objects/by-semantic-id/x", nil)
	assertErrorClass(t, rec, http.StatusServiceUnavailable, "LookupIndexUnavailable")
	var body APIError
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Error.Message != lookupIndexUnavailableMessage {
		t.Fatalf("message=%q, want fixed remediation", body.Error.Message)
	}
}
