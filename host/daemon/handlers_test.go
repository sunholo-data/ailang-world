package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

func newHandlerDaemon(t *testing.T) *Daemon {
	t.Helper()
	d, err := New(Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: DefaultBindHost})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() {
		if err := d.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	})
	return d
}

func requestRecorder(t *testing.T, d *Daemon, method, target string, body io.Reader) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, target, body)
	rec := httptest.NewRecorder()
	d.Handler().ServeHTTP(rec, req)
	return rec
}

func testObject(label string) store.Object {
	payload := []byte("payload-" + label)
	return store.Object{
		Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte("interface-" + label)),
		SemanticID: "test/" + label, Provenance: "handlers-test", Payload: payload,
	}
}

func testCommit(observed store.World, index int64, label string) store.Commit {
	object := testObject(label)
	entryHash := hashref.SumSHA256([]byte(fmt.Sprintf("entry-%d-%s", index, label)))
	return store.Commit{
		ObservedHead: observed.Ref,
		Objects:      []store.Object{object},
		NextWorld: store.World{
			Ref: hashref.SumSHA256([]byte(fmt.Sprintf("world-%d-%s", index, label))), Revision: index,
			StateRoot: hashref.SumSHA256([]byte(fmt.Sprintf("state-%d-%s", index, label))), LogHead: entryHash,
		},
		Entry: store.LogEntry{
			Header: store.LogHeader{
				EntryIndex: index, SemanticsEpoch: 1,
				TransitionFn:  hashref.SumSHA256([]byte("fn-" + label)),
				Interpreter:   hashref.SumSHA256([]byte("interpreter")),
				PrevEntryHash: observed.LogHead, WrittenBy: "handlers-test",
			},
			EntryHash: entryHash, TransitionRef: object.Hash,
		},
	}
}

func encodeCommit(c store.Commit) []byte {
	objects := make([]objectRequest, len(c.Objects))
	for i, object := range c.Objects {
		objects[i] = objectRequest{
			Hash: object.Hash.String(), InterfaceHash: object.InterfaceHash.String(),
			SemanticID: object.SemanticID, Provenance: object.Provenance, Payload: object.Payload,
		}
	}
	body, err := json.Marshal(commitRequest{
		ObservedHead: c.ObservedHead.String(), Objects: objects,
		NextWorld: worldRequest{
			Ref: c.NextWorld.Ref.String(), Revision: c.NextWorld.Revision,
			StateRoot: c.NextWorld.StateRoot.String(), LogHead: c.NextWorld.LogHead.String(),
		},
		Entry: logEntryRequest{
			Header: logHeaderRequest{
				EntryIndex: c.Entry.Header.EntryIndex, SemanticsEpoch: c.Entry.Header.SemanticsEpoch,
				TransitionFn:  c.Entry.Header.TransitionFn.String(),
				Interpreter:   c.Entry.Header.Interpreter.String(),
				PrevEntryHash: c.Entry.Header.PrevEntryHash.String(), WrittenBy: c.Entry.Header.WrittenBy,
			},
			EntryHash: c.Entry.EntryHash.String(), TransitionRef: c.Entry.TransitionRef.String(),
		},
	})
	if err != nil {
		panic(err)
	}
	return body
}

// seedGenesisEmbedded plants a selected genesis world through the EMBEDDED
// store API. The name says "embedded" because it is not REST: it exists so the
// read-route, conflict and body-cap tests can start from a populated store
// without re-driving genesis over HTTP every time.
//
// It was called seedRESTGenesis until the M2.B evaluator pointed out the name
// was doing real damage — it made the equivalence test LOOK like it drove
// genesis over REST when it did not, which is how the genesis-over-REST gap
// stayed invisible through the executor's own review. The test that genuinely
// drives genesis over the wire is TestRESTGenesisAndCommitAreByteEquivalent.
func seedGenesisEmbedded(t *testing.T, d *Daemon, label string) store.World {
	t.Helper()
	genesis := store.World{
		Ref: hashref.SumSHA256([]byte("genesis-world-" + label)), Revision: 0,
		StateRoot: hashref.SumSHA256([]byte("genesis-state-" + label)),
		LogHead:   hashref.SumSHA256([]byte("genesis-log-" + label)),
	}
	if err := d.store.PutWorld(genesis); err != nil {
		t.Fatalf("PutWorld genesis: %v", err)
	}
	if err := d.store.SelectHead(genesis.Ref); err != nil {
		t.Fatalf("SelectHead genesis: %v", err)
	}
	return genesis
}

func assertErrorClass(t *testing.T, rec *httptest.ResponseRecorder, status int, class string) APIError {
	t.Helper()
	if rec.Code != status {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, status, rec.Body)
	}
	var body APIError
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode API error: %v; body=%s", err, rec.Body)
	}
	if body.Error.Class != class {
		t.Fatalf("error class = %q, want %q; body=%s", body.Error.Class, class, rec.Body)
	}
	return body
}

func TestClampLimitMirrorsSketchVectors(t *testing.T) {
	for _, test := range []struct{ in, want int }{{0, 100}, {25, 25}, {1000, 500}, {-2, 100}} {
		if got := clampLimit(test.in); got != test.want {
			t.Errorf("clampLimit(%d) = %d, want %d", test.in, got, test.want)
		}
	}
}

func TestReadRoutesAndPayloadGate(t *testing.T) {
	d := newHandlerDaemon(t)
	genesis := seedGenesisEmbedded(t, d, "reads")
	commit := testCommit(genesis, 1, "reads")
	if err := d.store.Commit(commit); err != nil {
		t.Fatalf("seed Commit: %v", err)
	}

	t.Run("world", func(t *testing.T) {
		rec := requestRecorder(t, d, http.MethodGet, "/v1/worlds/"+commit.NextWorld.Ref.String(), nil)
		if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), commit.NextWorld.StateRoot.String()) {
			t.Fatalf("world response: status=%d body=%s", rec.Code, rec.Body)
		}
	})
	t.Run("object payload is opt-in base64", func(t *testing.T) {
		path := "/v1/objects/" + commit.Objects[0].Hash.String()
		without := requestRecorder(t, d, http.MethodGet, path, nil)
		if without.Code != http.StatusOK {
			t.Fatalf("default object status=%d body=%s", without.Code, without.Body)
		}
		var defaultBody map[string]json.RawMessage
		if err := json.Unmarshal(without.Body.Bytes(), &defaultBody); err != nil {
			t.Fatal(err)
		}
		if _, present := defaultBody["payload"]; present {
			t.Fatalf("default object leaked payload: %s", without.Body)
		}
		with := requestRecorder(t, d, http.MethodGet, path+"?payload=true", nil)
		var body objectResponse
		if err := json.Unmarshal(with.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if body.Payload == nil || !bytes.Equal(*body.Payload, commit.Objects[0].Payload) {
			t.Fatalf("payload round trip = %v, want %q", body.Payload, commit.Objects[0].Payload)
		}
	})
	t.Run("log", func(t *testing.T) {
		rec := requestRecorder(t, d, http.MethodGet, "/v1/log/1", nil)
		if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), commit.Entry.EntryHash.String()) {
			t.Fatalf("log response: status=%d body=%s", rec.Code, rec.Body)
		}
	})
	t.Run("multi-segment registry name", func(t *testing.T) {
		rec := requestRecorder(t, d, http.MethodGet, "/v1/registry/world/epoch-registry/v1", nil)
		if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"name":"world/epoch-registry/v1"`) {
			t.Fatalf("registry response: status=%d body=%s", rec.Code, rec.Body)
		}
	})
}

func TestBadRequestAndNotFoundAreDistinct(t *testing.T) {
	d := newHandlerDaemon(t)
	absent := hashref.SumSHA256([]byte("well-formed-but-absent")).String()
	for _, test := range []struct {
		name, target string
		status       int
		class        string
	}{
		{"malformed world", "/v1/worlds/not-a-ref", 400, "BadRequest"},
		{"malformed object", "/v1/objects/not-a-ref", 400, "BadRequest"},
		{"absent world", "/v1/worlds/" + absent, 404, "NotFound"},
		{"absent object", "/v1/objects/" + absent, 404, "NotFound"},
		{"absent log", "/v1/log/987654", 404, "NotFound"},
		{"absent registry", "/v1/registry/world/absent/v1", 404, "NotFound"},
	} {
		t.Run(test.name, func(t *testing.T) {
			assertErrorClass(t, requestRecorder(t, d, http.MethodGet, test.target, nil), test.status, test.class)
		})
	}
}

func TestHeadErrorsUseAPIEnvelope(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		d := newHandlerDaemon(t)
		assertErrorClass(t, requestRecorder(t, d, http.MethodGet, "/v1/head", nil),
			http.StatusNotFound, "NotFound")
	})
	t.Run("internal", func(t *testing.T) {
		d := newHandlerDaemon(t)
		if err := d.store.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
		assertErrorClass(t, requestRecorder(t, d, http.MethodGet, "/v1/head", nil),
			http.StatusInternalServerError, "Internal")
	})
}

func TestGETRoutesRejectOtherMethods(t *testing.T) {
	d := newHandlerDaemon(t)
	for _, test := range []struct {
		method string
		target string
	}{
		{http.MethodPost, "/v1/health"},
		{http.MethodDelete, "/v1/head"},
	} {
		rec := requestRecorder(t, d, test.method, test.target, nil)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s %s status=%d, want 405; body=%s", test.method, test.target, rec.Code, rec.Body)
		}
	}
}

func TestStaleHeadConflictBodySupportsReplan(t *testing.T) {
	d := newHandlerDaemon(t)
	genesis := seedGenesisEmbedded(t, d, "conflict")
	first := testCommit(genesis, 1, "winner")
	if rec := requestRecorder(t, d, http.MethodPost, "/v1/commit", bytes.NewReader(encodeCommit(first))); rec.Code != 200 {
		t.Fatalf("first commit: status=%d body=%s", rec.Code, rec.Body)
	}

	stale := testCommit(genesis, 2, "stale")
	conflictRec := requestRecorder(t, d, http.MethodPost, "/v1/commit", bytes.NewReader(encodeCommit(stale)))
	body := assertErrorClass(t, conflictRec, http.StatusConflict, "HeadConflict")
	observed, err := hashref.Parse(body.Error.ObservedHead)
	if err != nil {
		t.Fatalf("parse conflict observedHead: %v", err)
	}
	selected, err := hashref.Parse(body.Error.SelectedHead)
	if err != nil {
		t.Fatalf("parse conflict selectedHead: %v", err)
	}
	if observed.String() != genesis.Ref.String() || selected.String() != first.NextWorld.Ref.String() {
		t.Fatalf("conflict heads observed=%s selected=%s", observed, selected)
	}

	// Genuine re-plan: load the selected world named by the body, construct the
	// successor against it, and successfully commit that successor.
	selectedWorld, ok, err := d.store.GetWorld(context.Background(), selected)
	if err != nil || !ok {
		t.Fatalf("GetWorld(selected from 409): ok=%v err=%v", ok, err)
	}
	replanned := testCommit(selectedWorld, 2, "replanned")
	rec := requestRecorder(t, d, http.MethodPost, "/v1/commit", bytes.NewReader(encodeCommit(replanned)))
	if rec.Code != http.StatusOK {
		t.Fatalf("replanned commit: status=%d body=%s", rec.Code, rec.Body)
	}
}

func TestLogRangeClampAndDefaultAreNonVacuous(t *testing.T) {
	d := newHandlerDaemon(t)
	current := seedGenesisEmbedded(t, d, "range")
	const entries = 510
	for i := int64(0); i < entries; i++ {
		commit := testCommit(current, i, fmt.Sprintf("range-%d", i))
		if err := d.store.Commit(commit); err != nil {
			t.Fatalf("Commit(%d): %v", i, err)
		}
		current = commit.NextWorld
	}
	for _, test := range []struct {
		name, query string
		want        int
	}{
		{"default", "", 100},
		{"zero defaults", "?limit=0", 100},
		{"clamped", "?limit=1000", 500},
	} {
		t.Run(test.name, func(t *testing.T) {
			rec := requestRecorder(t, d, http.MethodGet, "/v1/log"+test.query, nil)
			if rec.Code != http.StatusOK {
				t.Fatalf("status=%d body=%s", rec.Code, rec.Body)
			}
			var body logRangeResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if len(body.Items) != test.want {
				t.Fatalf("returned %d items, want %d (store has %d)", len(body.Items), test.want, entries)
			}
		})
	}

	t.Run("non-zero from starts at requested index", func(t *testing.T) {
		const from = 37
		rec := requestRecorder(t, d, http.MethodGet, "/v1/log?from=37&limit=3", nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("status=%d body=%s", rec.Code, rec.Body)
		}
		var body logRangeResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if len(body.Items) != 3 || body.Items[0].Header.EntryIndex != from {
			t.Fatalf("indices begin at %v, want %d; body=%s", func() any {
				if len(body.Items) == 0 {
					return "<empty>"
				}
				return body.Items[0].Header.EntryIndex
			}(), from, rec.Body)
		}
	})
}

// genesisCommit builds a TRUE genesis commit: the ZERO observed head, which is
// the only value store.Commit accepts against a nil selected head, and which
// renders as the empty string in canonical HashRef text — exactly what
// parseGenesisRef reads back.
//
// PrevEntryHash is a REAL hash, not the zero value, matching M1's own genesis
// convention (store_test.go:103 seeds it from the genesis world's LogHead). A
// zero PrevEntryHash is writable by store.Commit but NOT readable by
// store.GetLogEntry, so using one here would build an episode whose own log
// cannot be read back. See parseGenesisRef.
func genesisCommit(label string) store.Commit {
	object := testObject("genesis-" + label)
	entryHash := hashref.SumSHA256([]byte("genesis-entry-" + label))
	return store.Commit{
		ObservedHead: hashref.HashRef{},
		Objects:      []store.Object{object},
		NextWorld: store.World{
			Ref: hashref.SumSHA256([]byte("genesis-world-" + label)), Revision: 0,
			StateRoot: hashref.SumSHA256([]byte("genesis-state-" + label)), LogHead: entryHash,
		},
		Entry: store.LogEntry{
			Header: store.LogHeader{
				EntryIndex: 0, SemanticsEpoch: 1,
				TransitionFn:  hashref.SumSHA256([]byte("genesis-fn-" + label)),
				Interpreter:   hashref.SumSHA256([]byte("interpreter")),
				PrevEntryHash: hashref.SumSHA256([]byte("genesis-prev-" + label)),
				WrittenBy:     "handlers-test",
			},
			EntryHash: entryHash, TransitionRef: object.Hash,
		},
	}
}

func postCommit(t *testing.T, baseURL string, c store.Commit) {
	t.Helper()
	resp, err := http.Post(baseURL+"/v1/commit", "application/json", bytes.NewReader(encodeCommit(c)))
	if err != nil {
		t.Fatalf("REST POST /v1/commit: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("REST commit status=%d body=%s", resp.StatusCode, body)
	}
}

// TestRESTGenesisAndCommitAreByteEquivalent is the M2.B acceptance check "a
// genesis+commit episode driven ENTIRELY over REST reproduces the store-level
// result byte-for-byte".
//
// Entirely means entirely: the REST daemon receives BOTH the genesis commit and
// its successor over a real loopback socket and never has an embedded store call
// made against it. The comparison is on BYTES — the selected head, the stored
// object payloads, and the encoded log entries — not on response shapes, because
// two surfaces can agree on JSON and still write different things to the kernel.
//
// Round 1 of this milestone shipped a weaker version of this test that seeded
// genesis through the embedded API on both sides, and the surface could not in
// fact express a genesis commit at all (POST /v1/commit answered 400). The
// evaluator blocked on it; parseGenesisRef is the fix and this test is its gate.
func TestRESTGenesisAndCommitAreByteEquivalent(t *testing.T) {
	embedded := newHandlerDaemon(t)
	rest := newHandlerDaemon(t)

	genesis := genesisCommit("equivalence")
	successor := testCommit(genesis.NextWorld, 1, "equivalence")

	// Arm 1: the embedded kernel API.
	if err := embedded.store.Commit(genesis); err != nil {
		t.Fatalf("embedded genesis Commit: %v", err)
	}
	if err := embedded.store.Commit(successor); err != nil {
		t.Fatalf("embedded successor Commit: %v", err)
	}

	// Arm 2: the same episode, entirely over REST.
	server := httptest.NewServer(rest.Handler())
	defer server.Close()
	postCommit(t, server.URL, genesis)
	postCommit(t, server.URL, successor)

	embeddedHead, okE, errE := embedded.store.SelectedHead(context.Background())
	restHead, okR, errR := rest.store.SelectedHead(context.Background())
	if errE != nil || errR != nil || !okE || !okR {
		t.Fatalf("SelectedHead: embedded(ok=%v err=%v) REST(ok=%v err=%v)", okE, errE, okR, errR)
	}
	if embeddedHead.String() != restHead.String() {
		t.Fatalf("selected heads differ: embedded=%s REST=%s", embeddedHead, restHead)
	}
	if restHead.String() != successor.NextWorld.Ref.String() {
		t.Fatalf("REST head = %s, want the successor world %s", restHead, successor.NextWorld.Ref)
	}

	// Both objects, byte-for-byte.
	for _, c := range []store.Commit{genesis, successor} {
		ref := c.Objects[0].Hash
		embeddedObject, ok, err := embedded.store.GetObject(context.Background(), ref)
		if err != nil || !ok {
			t.Fatalf("embedded GetObject(%s): ok=%v err=%v", ref, ok, err)
		}
		restObject, ok, err := rest.store.GetObject(context.Background(), ref)
		if err != nil || !ok {
			t.Fatalf("REST GetObject(%s): ok=%v err=%v", ref, ok, err)
		}
		if !bytes.Equal(embeddedObject.Payload, restObject.Payload) {
			t.Fatalf("payload bytes differ for %s: embedded=%x REST=%x", ref, embeddedObject.Payload, restObject.Payload)
		}
	}

	// Both log entries, compared through their canonical encodings.
	for index := int64(0); index <= 1; index++ {
		embeddedEntry, ok, err := embedded.store.GetLogEntry(context.Background(), index)
		if err != nil || !ok {
			t.Fatalf("embedded GetLogEntry(%d): ok=%v err=%v", index, ok, err)
		}
		restEntry, ok, err := rest.store.GetLogEntry(context.Background(), index)
		if err != nil || !ok {
			t.Fatalf("REST GetLogEntry(%d): ok=%v err=%v", index, ok, err)
		}
		embeddedJSON, _ := json.Marshal(logJSON(embeddedEntry))
		restJSON, _ := json.Marshal(logJSON(restEntry))
		if !bytes.Equal(embeddedJSON, restJSON) {
			t.Fatalf("log entry %d differs:\n embedded=%s\n REST    =%s", index, embeddedJSON, restJSON)
		}
	}
}

// TestGenesisRefLenienceIsExactlyOneField guards the blast radius of
// parseGenesisRef. The empty string is accepted for observedHead — the one field
// the kernel defines a zero value for — and for NOTHING else. Without this test
// the fix for the genesis gap could quietly decay into "a missing ref is a zero
// ref", which is a far worse defect than the one it closed.
//
// prevEntryHash is in this list on purpose: it is the field that LOOKS like it
// should be lenient at genesis and must not be, because store.Commit will write
// a zero there that store.GetLogEntry cannot read back.
func TestGenesisRefLenienceIsExactlyOneField(t *testing.T) {
	d := newHandlerDaemon(t)
	for _, test := range []struct{ name, field string }{
		{"nextWorld.ref", "nextWorld.ref"},
		{"nextWorld.stateRoot", "nextWorld.stateRoot"},
		{"nextWorld.logHead", "nextWorld.logHead"},
		{"entry.entryHash", "entry.entryHash"},
		{"entry.transitionRef", "entry.transitionRef"},
		{"entry.header.transitionFn", "entry.header.transitionFn"},
		{"entry.header.interpreter", "entry.header.interpreter"},
		{"entry.header.prevEntryHash", "entry.header.prevEntryHash"},
		{"objects[0].hash", "objects[0].hash"},
		{"objects[0].interfaceHash", "objects[0].interfaceHash"},
	} {
		t.Run(test.name, func(t *testing.T) {
			var request commitRequest
			if err := json.Unmarshal(encodeCommit(genesisCommit("lenience")), &request); err != nil {
				t.Fatal(err)
			}
			blankField(t, &request, test.field)
			body, err := json.Marshal(request)
			if err != nil {
				t.Fatal(err)
			}
			got := assertErrorClass(t, requestRecorder(t, d, http.MethodPost, "/v1/commit", bytes.NewReader(body)),
				http.StatusBadRequest, "BadRequest")
			if !strings.Contains(got.Error.Message, test.field) {
				t.Fatalf("400 message %q does not name the offending field %q", got.Error.Message, test.field)
			}
		})
	}
}

func blankField(t *testing.T, r *commitRequest, field string) {
	t.Helper()
	switch field {
	case "nextWorld.ref":
		r.NextWorld.Ref = ""
	case "nextWorld.stateRoot":
		r.NextWorld.StateRoot = ""
	case "nextWorld.logHead":
		r.NextWorld.LogHead = ""
	case "entry.entryHash":
		r.Entry.EntryHash = ""
	case "entry.transitionRef":
		r.Entry.TransitionRef = ""
	case "entry.header.transitionFn":
		r.Entry.Header.TransitionFn = ""
	case "entry.header.interpreter":
		r.Entry.Header.Interpreter = ""
	case "entry.header.prevEntryHash":
		r.Entry.Header.PrevEntryHash = ""
	case "objects[0].hash":
		r.Objects[0].Hash = ""
	case "objects[0].interfaceHash":
		r.Objects[0].InterfaceHash = ""
	default:
		t.Fatalf("blankField: unhandled field %q", field)
	}
}
