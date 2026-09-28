package daemon

import (
	"bytes"
	"context"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/workbench"
)

func TestWorkbenchRouteIsReadOnly(t *testing.T) {
	d := newHandlerDaemon(t)
	seedGenesisEmbedded(t, d, "workbench-methods")
	for _, test := range []struct {
		method string
		want   int
	}{
		{http.MethodGet, http.StatusOK},
		{http.MethodPost, http.StatusMethodNotAllowed},
		{http.MethodPut, http.StatusMethodNotAllowed},
		{http.MethodDelete, http.StatusMethodNotAllowed},
		{http.MethodPatch, http.StatusMethodNotAllowed},
	} {
		t.Run(test.method, func(t *testing.T) {
			rec := requestRecorder(t, d, test.method, "/workbench", nil)
			if rec.Code != test.want {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, test.want, rec.Body)
			}
		})
	}
}

func TestWorkbenchSecurityHeaders(t *testing.T) {
	d := newHandlerDaemon(t)
	seedGenesisEmbedded(t, d, "workbench-headers")
	wants := map[string]string{
		"Content-Type":  "text/html; charset=utf-8",
		"Cache-Control": "no-store",
		// LITERAL, never the production `workbenchCSP` symbol. Asserting against the
		// constant the handler itself sets makes expected and actual move together, so
		// M22 (delete the `default-src 'none'; ` token) is invisible: measured, that
		// mutant landed, built rc=0 and left the whole package rc=0 with an empty FAIL
		// set. A tautological oracle cannot fail at any point in the sprint.
		"Content-Security-Policy": "default-src 'none'; style-src 'unsafe-inline'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'",
		"X-Content-Type-Options":  "nosniff",
		"Referrer-Policy":         "no-referrer",
	}
	for _, test := range []struct {
		name   string
		target string
		status int
	}{
		{"success", "/workbench", http.StatusOK},
		{"bad request", "/workbench?paylod=1", http.StatusBadRequest},
	} {
		t.Run(test.name, func(t *testing.T) {
			rec := requestRecorder(t, d, http.MethodGet, test.target, nil)
			if rec.Code != test.status {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, test.status, rec.Body)
			}
			for name, want := range wants {
				if got := rec.Header().Get(name); got != want {
					t.Errorf("%s = %q, want %q", name, got, want)
				}
			}
		})
	}
}

func TestWorkbenchRendersSeededWorldAndTimeline(t *testing.T) {
	d := newHandlerDaemon(t)
	genesis := seedGenesisEmbedded(t, d, "workbench-render")
	first := testCommit(genesis, 0, "workbench-first")
	if err := d.store.Commit(first); err != nil {
		t.Fatalf("first Commit: %v", err)
	}
	second := testCommit(first.NextWorld, 1, "workbench-second")
	if err := d.store.Commit(second); err != nil {
		t.Fatalf("second Commit: %v", err)
	}

	// The observable is the log entry hash, NOT the committed object hash. On
	// `/workbench` (no selection) the timeline rows render no edges, so the
	// object hash reaches the page only through the selected entry's
	// transitionRef edge -- pinned by TestWorkbenchSelectedEntry, not here. The
	// entry hash is written by the timeline loop and by nothing else on this
	// page, which is what makes it a pin on the mechanism rather than on a
	// sibling channel.
	entryHashes := make([]string, 0, 2)
	for index := int64(0); index < 2; index++ {
		entry, ok, err := d.store.GetLogEntry(context.Background(), index)
		if err != nil || !ok {
			t.Fatalf("GetLogEntry(%d): ok=%v err=%v", index, ok, err)
		}
		entryHashes = append(entryHashes, entry.EntryHash.String())
	}
	// POSITIVE CONTROL: two entries must genuinely differ before "both are
	// present" means anything -- one hash rendered twice would otherwise pass.
	if entryHashes[0] == entryHashes[1] {
		t.Fatalf("positive control: distinct commits produced equal entry hashes %q", entryHashes[0])
	}

	rec := requestRecorder(t, d, http.MethodGet, "/workbench", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body)
	}
	body := rec.Body.String()
	// Assert the hash inside the element ONLY the {{.EntryHash}} action produces. A
	// bare Contains is not sole-sourced for entry 0: store.Commit chains the log, so
	// entry 1's PrevEntryHash field renders entry 0's hash verbatim, and dropping
	// entry 0 from the page entirely still left the bare assertion green.
	for _, hash := range entryHashes {
		want := `<dt>entry hash</dt><dd><span class="hash" title="` + hash + `"`
		if !strings.Contains(body, want) {
			t.Errorf("rendered body does not contain timeline entry hash %q in its own entry-hash element", hash)
		}
	}
	// The world section is the other half of "seeded world AND timeline".
	head, ok, err := d.store.SelectedHead(context.Background())
	if err != nil || !ok {
		t.Fatalf("SelectedHead: ok=%v err=%v", ok, err)
	}
	if !strings.Contains(body, head.String()) {
		t.Errorf("rendered body does not contain selected head world ref %q", head.String())
	}
}

func TestWorkbenchRefusalBranches(t *testing.T) {
	d := newHandlerDaemon(t)
	genesis := seedGenesisEmbedded(t, d, "workbench-refusals")
	commit := testCommit(genesis, 0, "workbench-refusals-entry")
	if err := d.store.Commit(commit); err != nil {
		t.Fatalf("Commit: %v", err)
	}

	absentWorld := hashref.SumSHA256([]byte("workbench-absent-world")).String()
	absentObject := hashref.SumSHA256([]byte("workbench-absent-object")).String()
	world := commit.NextWorld.Ref.String()
	object := commit.Objects[0].Hash.String()

	tests := []struct {
		name    string
		target  string
		status  int
		class   string
		message string
		setup   func()
	}{
		{"unknown-parameter", "/workbench?paylod=1", http.StatusBadRequest, "BadRequest", "unsupported workbench query parameter", nil},
		{"duplicate-parameter", "/workbench?world=" + world + "&world=" + world, http.StatusBadRequest, "BadRequest", "duplicate workbench query parameter", nil},
		{"unsupported-combination", "/workbench?from=0", http.StatusBadRequest, "BadRequest", "unsupported workbench parameter combination", nil},
		{"unsupported-pair-world-from", "/workbench?world=" + world + "&from=0", http.StatusBadRequest, "BadRequest", "unsupported workbench parameter combination", nil},
		{"unsupported-pair-world-payload", "/workbench?world=" + world + "&payload=1", http.StatusBadRequest, "BadRequest", "unsupported workbench parameter combination", nil},
		{"unsupported-triple-world-object-payload", "/workbench?world=" + world + "&object=" + object + "&payload=1", http.StatusBadRequest, "BadRequest", "unsupported workbench parameter combination", nil},
		{"malformed-payload", "/workbench?object=" + object + "&payload=true", http.StatusBadRequest, "BadRequest", "malformed payload flag", nil},
		{"malformed-world", "/workbench?world=not-a-hash", http.StatusBadRequest, "BadRequest", "malformed world reference", nil},
		{"absent-world", "/workbench?world=" + absentWorld, http.StatusNotFound, "NotFound", "world reference not found", nil},
		{"malformed-object", "/workbench?object=not-a-hash", http.StatusBadRequest, "BadRequest", "malformed object reference", nil},
		{"absent-object", "/workbench?object=" + absentObject, http.StatusNotFound, "NotFound", "object reference not found", nil},
		{"negative-from", "/workbench?from=-1&entry=0", http.StatusBadRequest, "BadRequest", "from index must be non-negative", nil},
		{"malformed-entry", "/workbench?from=0&entry=not-an-index", http.StatusBadRequest, "BadRequest", "malformed entry index", nil},
		{"absent-entry", "/workbench?from=0&entry=99", http.StatusNotFound, "NotFound", "log entry not found", nil},
		{"from-overflow", "/workbench?from=9223372036854775807&entry=0", http.StatusBadRequest, "BadRequest", "from index overflows", nil},
		{"store-error", "/workbench", http.StatusInternalServerError, "Internal", "internal store failure", func() {
			d.reads = failingStore{Store: d.store}
			d.errLog = &bytes.Buffer{}
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			oldReads, oldErrLog := d.reads, d.errLog
			defer func() {
				d.reads, d.errLog = oldReads, oldErrLog
			}()
			if test.setup != nil {
				test.setup()
			}
			rec := requestRecorder(t, d, http.MethodGet, test.target, nil)
			if rec.Code != test.status {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, test.status, rec.Body)
			}
			body := rec.Body.String()
			if !strings.Contains(body, ">"+test.class+"<") {
				t.Errorf("body does not contain class token %q: %s", test.class, body)
			}
			if !strings.Contains(body, ">"+test.message+"<") {
				t.Errorf("body does not contain branch message %q: %s", test.message, body)
			}
			assertWorkbenchSecurityHeaders(t, rec.Header())
		})
	}

	t.Run("accepted-keys-return-200", func(t *testing.T) {
		// Exercise every member of the closed key vocabulary across exactly the
		// four supported non-empty states; one all-keys query is intentionally
		// invalid because the grammar is a set of states, not a key allowlist.
		targets := []string{
			"/workbench?world=" + world,
			"/workbench?from=0&entry=0",
			"/workbench?object=" + object,
			"/workbench?object=" + object + "&payload=0",
			"/workbench?object=" + object + "&payload=1",
		}
		for _, target := range targets {
			rec := requestRecorder(t, d, http.MethodGet, target, nil)
			if rec.Code != http.StatusOK {
				t.Fatalf("%s: status = %d, want 200; body=%s", target, rec.Code, rec.Body)
			}
			assertWorkbenchSecurityHeaders(t, rec.Header())
		}
	})
}

// TestSupportedWorkbenchQueryTruthTable walks every subset of the five-key
// vocabulary. The function reads only key presence and count, so these 32
// subsets are its whole reachable domain: exactly five are accepted states.
func TestSupportedWorkbenchQueryTruthTable(t *testing.T) {
	keys := []string{"entry", "from", "object", "payload", "world"}
	accepted := map[string]bool{"none": true, "world": true, "object": true, "entry+from": true, "object+payload": true}
	seen := 0
	for mask := 0; mask < 1<<len(keys); mask++ {
		query := map[string][]string{}
		var names []string
		for i, key := range keys {
			if mask&(1<<i) != 0 {
				query[key] = []string{"x"}
				names = append(names, key)
			}
		}
		name := strings.Join(names, "+")
		if name == "" {
			name = "none"
		}
		if accepted[name] {
			seen++
		}
		t.Run(name, func(t *testing.T) {
			if got := supportedWorkbenchQuery(query); got != accepted[name] {
				t.Fatalf("supportedWorkbenchQuery(%s) = %v, want %v", name, got, accepted[name])
			}
		})
	}
	if seen != len(accepted) {
		t.Fatalf("accepted states enumerated = %d, want %d", seen, len(accepted))
	}
}

func commitWorkbenchPayload(t *testing.T, d *Daemon, payload []byte, label string) hashref.HashRef {
	t.Helper()
	genesis := seedGenesisEmbedded(t, d, label+"-genesis")
	commit := testCommit(genesis, 0, label)
	object := store.Object{
		Hash:          hashref.SumSHA256(payload),
		InterfaceHash: hashref.SumSHA256([]byte("interface-" + label)),
		SemanticID:    "test/" + label,
		Provenance:    "workbench-test",
		Payload:       payload,
	}
	commit.Objects = []store.Object{object}
	commit.Entry.TransitionRef = object.Hash
	if err := d.store.Commit(commit); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	return object.Hash
}

func TestWorkbenchPayloadPreviewBound(t *testing.T) {
	t.Run("default-off", func(t *testing.T) {
		d := newHandlerDaemon(t)
		ref := commitWorkbenchPayload(t, d, []byte("PAYLOAD-MARKER-7ac"), "payload-default-off")
		rec := requestRecorder(t, d, http.MethodGet, "/workbench?object="+ref.String(), nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body)
		}
		if strings.Contains(rec.Body.String(), "PAYLOAD-MARKER-7ac") {
			t.Fatal("payload rendered without payload=1")
		}
	})

	t.Run("opt-in", func(t *testing.T) {
		d := newHandlerDaemon(t)
		ref := commitWorkbenchPayload(t, d, []byte("PAYLOAD-MARKER-7ac"), "payload-opt-in")
		rec := requestRecorder(t, d, http.MethodGet, "/workbench?object="+ref.String()+"&payload=1", nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body)
		}
		if !strings.Contains(rec.Body.String(), "PAYLOAD-MARKER-7ac") {
			t.Fatal("payload=1 did not render payload")
		}
	})

	t.Run("small-renders-in-full", func(t *testing.T) {
		d := newHandlerDaemon(t)
		payload := []byte("0123456789abcdefghijklmnopqrstuv")
		ref := commitWorkbenchPayload(t, d, payload, "payload-small")
		rec := requestRecorder(t, d, http.MethodGet, "/workbench?object="+ref.String()+"&payload=1", nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body)
		}
		body := rec.Body.String()
		if !strings.Contains(body, string(payload)) {
			t.Fatal("small payload did not render in full")
		}
		if strings.Contains(body, "truncated") {
			t.Fatal("small payload incorrectly marked truncated")
		}
	})

	t.Run("oversize", func(t *testing.T) {
		d := newHandlerDaemon(t)
		payload := bytes.Repeat([]byte("x"), workbench.MaxPayloadPreview+4096)
		copy(payload[len(payload)-len("TAIL-MARKER-3d1"):], "TAIL-MARKER-3d1")
		ref := commitWorkbenchPayload(t, d, payload, "payload-oversize")
		rec := requestRecorder(t, d, http.MethodGet, "/workbench?object="+ref.String()+"&payload=1", nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body)
		}
		body := rec.Body.String()
		if !strings.Contains(body, "truncated") {
			t.Fatal("oversize payload not marked truncated")
		}
		if strings.Contains(body, "TAIL-MARKER-3d1") {
			t.Fatal("oversize payload rendered bytes beyond preview cap")
		}
	})
}

// seedWorkbenchLog commits n chained testCommit entries (indices 0..n-1) on a
// fresh genesis. Every entry's TransitionRef object is stored; its TransitionFn
// and Interpreter targets are not (testCommit).
func seedWorkbenchLog(t *testing.T, d *Daemon, n int) {
	t.Helper()
	world := seedGenesisEmbedded(t, d, "workbench-timeline-bound")
	started := time.Now()
	for index := int64(0); index < int64(n); index++ {
		commit := testCommit(world, index, "workbench-timeline-bound")
		if err := d.store.Commit(commit); err != nil {
			t.Fatalf("Commit(%d): %v", index, err)
		}
		world = commit.NextWorld
		if time.Since(started) > 30*time.Second {
			t.Fatalf("seeding exceeded 30 seconds after %d commits", index+1)
		}
	}
}

func TestWorkbenchTimelineBound(t *testing.T) {
	d := newHandlerDaemon(t)
	seedWorkbenchLog(t, d, workbench.WorkbenchPageLimit+5)
	if _, ok, err := d.store.GetLogEntry(context.Background(), 104); err != nil || !ok {
		t.Fatalf("positive control GetLogEntry(104): ok=%v err=%v", ok, err)
	}

	// The empty query is state 1 of the closed grammar of §2.2: it defaults from=0 and leaves
	// Page.Selected nil, so only timeline rows emit entry headings. `?from=0` alone is NOT in
	// the enumeration and is refused 400 by design (see the unsupported-combination arm above).
	rec := requestRecorder(t, d, http.MethodGet, "/workbench", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body)
	}
	body := rec.Body.String()
	if got := strings.Count(body, "<h3>entry "); got != workbench.WorkbenchPageLimit {
		t.Errorf("timeline entry count = %d, want %d", got, workbench.WorkbenchPageLimit)
	}
	if !strings.Contains(body, "<h3>entry 99</h3>") {
		t.Error("timeline does not contain entry 99")
	}
	if strings.Contains(body, "<h3>entry 100</h3>") {
		t.Error("timeline contains entry 100 beyond page bound")
	}
}

func assertWorkbenchSecurityHeaders(t *testing.T, header http.Header) {
	t.Helper()
	wants := map[string]string{
		"Content-Type":            "text/html; charset=utf-8",
		"Cache-Control":           "no-store",
		"Content-Security-Policy": "default-src 'none'; style-src 'unsafe-inline'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'",
		"X-Content-Type-Options":  "nosniff",
		"Referrer-Policy":         "no-referrer",
	}
	for name, want := range wants {
		if got := header.Get(name); got != want {
			t.Errorf("%s = %q, want %q", name, got, want)
		}
	}
}

// workbenchRegion returns the substring from start to the first end after it.
func workbenchRegion(body, start, end string) (string, bool) {
	i := strings.Index(body, start)
	if i < 0 {
		return "", false
	}
	j := strings.Index(body[i:], end)
	if j < 0 {
		return "", false
	}
	return body[i : i+j], true
}

const (
	selectedEntryStart = `<article aria-label="selected entry">`
	timelineStart      = `<section aria-label="timeline">`
	worldStart         = `<nav aria-label="world browser">`
)

type rootReadStore struct {
	readStore
	root  hashref.HashRef
	err   error
	reads int
}

func (s *rootReadStore) GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error) {
	if ref == s.root {
		s.reads++
		if s.err != nil {
			return store.Object{}, false, s.err
		}
	}
	return s.readStore.GetObject(ctx, ref)
}

func assertWorldNav(t *testing.T, body string, world store.World, stored bool) string {
	t.Helper()
	nav, ok := workbenchRegion(body, worldStart, "</nav>")
	if !ok {
		t.Fatal("world nav missing")
	}
	for _, want := range []string{
		`<a href="/workbench">workbench</a>`,
		`<dt>world</dt><dd><span class="hash" title="` + world.Ref.String() + `" aria-label="` + world.Ref.String() + `">` + world.Ref.String() + `</span></dd>`,
		`<dt>revision</dt><dd>` + strconv.FormatInt(world.Revision, 10) + `</dd>`,
		`<dt>log head</dt><dd><span class="hash" title="` + world.LogHead.String() + `" aria-label="` + world.LogHead.String() + `">` + world.LogHead.String() + `</span></dd>`,
	} {
		if !strings.Contains(nav, want) {
			t.Errorf("nav missing %q: %s", want, nav)
		}
	}
	if strings.Contains(nav, `aria-label="`+world.Ref.String()+`">`+world.Ref.String()+`</a>`) || strings.Contains(nav, `aria-label="`+world.LogHead.String()+`">`+world.LogHead.String()+`</a>`) {
		t.Error("world or log head is an anchor")
	}
	if stored {
		want := `stateRoot: <a href="/workbench?object=` + world.StateRoot.String() + `" class="hash" title="` + world.StateRoot.String() + `" aria-label="` + world.StateRoot.String() + `">` + world.StateRoot.String() + `</a>`
		if !strings.Contains(nav, want) || strings.Count(nav, `?object=`) != 1 {
			t.Errorf("stored root nav: %s", nav)
		}
	} else {
		want := `stateRoot: <span class="unavailable" role="note">UNAVAILABLE: object ` + world.StateRoot.String() + ` is not stored</span>`
		if !strings.Contains(nav, want) || strings.Contains(nav, `?object=`) {
			t.Errorf("missing root nav: %s", nav)
		}
	}
	return nav
}

func TestWorkbenchWorldPane(t *testing.T) {
	for _, stored := range []bool{false, true} {
		name := "missing-root"
		if stored {
			name = "stored-root"
		}
		t.Run(name, func(t *testing.T) {
			d := newHandlerDaemon(t)
			genesis := seedGenesisEmbedded(t, d, name)
			commit := testCommit(genesis, 0, name)
			commit.NextWorld.Revision = 17
			if stored {
				state := workbenchTestObject("state-"+name, "world/state", hashref.SumSHA256([]byte("state-interface")))
				commit.NextWorld.StateRoot = state.Hash
				commit.Objects = append(commit.Objects, state)
			}
			if err := d.store.Commit(commit); err != nil {
				t.Fatal(err)
			}
			_, ok, err := d.store.GetObject(context.Background(), commit.NextWorld.StateRoot)
			if err != nil || ok != stored {
				t.Fatalf("root control stored=%v err=%v, want %v", ok, err, stored)
			}
			for _, path := range []string{"/workbench", "/workbench?from=0&entry=0", "/workbench?object=" + commit.Entry.TransitionRef.String()} {
				rec := requestRecorder(t, d, http.MethodGet, path, nil)
				if rec.Code != 200 {
					t.Fatalf("%s status %d: %s", path, rec.Code, rec.Body)
				}
				assertWorldNav(t, rec.Body.String(), commit.NextWorld, stored)
			}
			if stored {
				counter := &rootReadStore{readStore: d.store, root: commit.NextWorld.StateRoot}
				d.reads = counter
				rec := requestRecorder(t, d, http.MethodGet, "/workbench", nil)
				if rec.Code != 200 || counter.reads != 1 {
					t.Fatalf("root reads=%d status=%d", counter.reads, rec.Code)
				}
				d.reads = d.store
				link := "/workbench?object=" + commit.NextWorld.StateRoot.String()
				if rec := requestRecorder(t, d, http.MethodGet, link, nil); rec.Code != 200 {
					t.Fatalf("root link %s status %d", link, rec.Code)
				}
			}
		})
	}
	t.Run("explicit-world", func(t *testing.T) {
		d := newHandlerDaemon(t)
		genesis := seedGenesisEmbedded(t, d, "explicit-world")
		old := testCommit(genesis, 0, "old-world")
		old.NextWorld.Revision = 17
		state := workbenchTestObject("old-state", "world/state", hashref.SumSHA256([]byte("state-interface")))
		old.NextWorld.StateRoot = state.Hash
		old.Objects = append(old.Objects, state)
		if err := d.store.Commit(old); err != nil {
			t.Fatal(err)
		}
		newer := testCommit(old.NextWorld, 1, "new-world")
		if err := d.store.Commit(newer); err != nil {
			t.Fatal(err)
		}
		if old.NextWorld.Ref == newer.NextWorld.Ref || old.NextWorld.StateRoot == newer.NextWorld.StateRoot || old.NextWorld.LogHead == newer.NextWorld.LogHead || old.NextWorld.Revision == newer.NextWorld.Revision {
			t.Fatal("world metadata must differ")
		}
		rec := requestRecorder(t, d, http.MethodGet, "/workbench?world="+old.NextWorld.Ref.String(), nil)
		if rec.Code != 200 {
			t.Fatalf("status %d: %s", rec.Code, rec.Body)
		}
		assertWorldNav(t, rec.Body.String(), old.NextWorld, true)
	})
	t.Run("no-selected-world", func(t *testing.T) {
		d := newHandlerDaemon(t)
		rec := requestRecorder(t, d, http.MethodGet, "/workbench", nil)
		if rec.Code != 200 {
			t.Fatalf("status %d", rec.Code)
		}
		nav, ok := workbenchRegion(rec.Body.String(), worldStart, "</nav>")
		if !ok || !strings.Contains(nav, `UNAVAILABLE: no world selected`) || strings.Count(nav, "<a ") != 1 || strings.Contains(nav, `?object=`) {
			t.Fatalf("no-world nav: %s", nav)
		}
	})
	for _, tc := range []struct {
		name   string
		err    error
		status int
		class  string
	}{{"internal", errSentinelInternal, 500, "Internal"}, {"timeout", context.DeadlineExceeded, 503, "Timeout"}} {
		t.Run("store-error/"+tc.name, func(t *testing.T) {
			d := newHandlerDaemon(t)
			genesis := seedGenesisEmbedded(t, d, tc.name)
			commit := testCommit(genesis, 0, tc.name)
			if err := d.store.Commit(commit); err != nil {
				t.Fatal(err)
			}
			d.errLog = &bytes.Buffer{}
			wrapper := &rootReadStore{readStore: d.store, root: commit.NextWorld.StateRoot, err: tc.err}
			d.reads = wrapper
			if rec := requestRecorder(t, d, http.MethodGet, "/workbench?world="+genesis.Ref.String(), nil); rec.Code != 200 {
				t.Fatalf("untargeted status %d", rec.Code)
			}
			rec := requestRecorder(t, d, http.MethodGet, "/workbench", nil)
			if rec.Code != tc.status || !strings.Contains(rec.Body.String(), ">"+tc.class+"<") || strings.Contains(rec.Body.String(), tc.err.Error()) || strings.Contains(rec.Body.String(), worldStart) {
				t.Fatalf("status=%d body=%s", rec.Code, rec.Body)
			}
		})
	}
}

// objectFailingStore fails GetObject only, so the entry-edge check is the one
// read that errors while every log and world read still reaches the real store.
type objectFailingStore struct {
	readStore
}

func (objectFailingStore) GetObject(context.Context, hashref.HashRef) (store.Object, bool, error) {
	return store.Object{}, false, errSentinelInternal
}

func TestWorkbenchSelectedEntry(t *testing.T) {
	d := newHandlerDaemon(t)
	genesis := seedGenesisEmbedded(t, d, "workbench-selected")
	commit := testCommit(genesis, 0, "workbench-selected")
	if err := d.store.Commit(commit); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	header := commit.Entry.Header
	transitionRef := commit.Entry.TransitionRef.String()
	transitionFn := header.TransitionFn.String()
	interpreter := header.Interpreter.String()
	// POSITIVE CONTROLS: the fixture must hold one stored and two unstored targets,
	// or the stored/unstored arms below assert nothing.
	for _, control := range []struct {
		name string
		ref  hashref.HashRef
		want bool
	}{
		{"TransitionRef", commit.Entry.TransitionRef, true},
		{"TransitionFn", header.TransitionFn, false},
		{"Interpreter", header.Interpreter, false},
	} {
		if _, ok, err := d.store.GetObject(context.Background(), control.ref); err != nil || ok != control.want {
			t.Fatalf("control GetObject(%s): ok=%v err=%v, want ok=%v", control.name, ok, err, control.want)
		}
	}

	selectedBody := func(t *testing.T) string {
		t.Helper()
		rec := requestRecorder(t, d, http.MethodGet, "/workbench?from=0&entry=0", nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body)
		}
		article, ok := workbenchRegion(rec.Body.String(), selectedEntryStart, "</article>")
		if !ok {
			t.Fatalf("no selected-entry article in %s", rec.Body)
		}
		return article
	}

	t.Run("stored-edge-links", func(t *testing.T) {
		article := selectedBody(t)
		if want := `transitionRef: <a href="/workbench?object=` + transitionRef + `"`; !strings.Contains(article, want) {
			t.Errorf("selected entry missing stored edge link %q: %s", want, article)
		}
	})

	t.Run("unstored-edge-unavailable", func(t *testing.T) {
		article := selectedBody(t)
		for _, edge := range []struct{ relation, ref string }{{"transitionFn", transitionFn}, {"interpreter", interpreter}} {
			want := edge.relation + `: <span class="unavailable" role="note">UNAVAILABLE: object ` + edge.ref + ` is not stored</span>`
			if !strings.Contains(article, want) {
				t.Errorf("selected entry missing %q: %s", want, article)
			}
		}
		if unwanted := `href="/workbench?object=` + transitionFn + `"`; strings.Contains(article, unwanted) {
			t.Errorf("unstored transitionFn rendered as a link %q", unwanted)
		}
	})

	t.Run("row-select-link", func(t *testing.T) {
		rec := requestRecorder(t, d, http.MethodGet, "/workbench", nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body)
		}
		body := rec.Body.String()
		if want := `href="/workbench?from=0&amp;entry=0">select entry 0</a>`; !strings.Contains(body, want) {
			t.Errorf("timeline row missing select link %q", want)
		}
		if strings.Contains(body, `aria-label="selected entry"`) {
			t.Error("/workbench with no entry= rendered a selected-entry article")
		}
	})

	t.Run("object-store-error", func(t *testing.T) {
		oldReads, oldErrLog := d.reads, d.errLog
		defer func() { d.reads, d.errLog = oldReads, oldErrLog }()
		if commit.NextWorld.StateRoot == commit.Entry.TransitionRef {
			t.Fatal("root and transition refs must differ")
		}
		d.reads = refFailingStore{readStore: d.store, fail: commit.Entry.TransitionRef}
		d.errLog = &bytes.Buffer{}
		// CONTROL: the same store serves the unselected page, so only the edge
		// check can be what fails below.
		if rec := requestRecorder(t, d, http.MethodGet, "/workbench", nil); rec.Code != http.StatusOK {
			t.Fatalf("control: /workbench status = %d, want 200; body=%s", rec.Code, rec.Body)
		}
		rec := requestRecorder(t, d, http.MethodGet, "/workbench?from=0&entry=0", nil)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500; body=%s", rec.Code, rec.Body)
		}
		if !strings.Contains(rec.Body.String(), ">Internal<") {
			t.Errorf("body does not contain >Internal<: %s", rec.Body)
		}
	})
}

var workbenchLinkPattern = regexp.MustCompile(`<a href="([^"]*)"[^>]*>([^<]*)</a>`)
var selectEntryText = regexp.MustCompile(`^select entry [0-9]+$`)

func TestWorkbenchTimelinePaging(t *testing.T) {
	d := newHandlerDaemon(t)
	seedWorkbenchLog(t, d, workbench.WorkbenchPageLimit+5)
	if _, ok, err := d.store.GetLogEntry(context.Background(), 104); err != nil || !ok {
		t.Fatalf("positive control GetLogEntry(104): ok=%v err=%v", ok, err)
	}
	get := func(t *testing.T, target string) string {
		t.Helper()
		rec := requestRecorder(t, d, http.MethodGet, target, nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status = %d, want 200; body=%s", target, rec.Code, rec.Body)
		}
		return rec.Body.String()
	}

	t.Run("head-page-has-next", func(t *testing.T) {
		body := get(t, "/workbench")
		if want := `href="/workbench?from=100&amp;entry=100">next</a>`; !strings.Contains(body, want) {
			t.Errorf("head page missing next link %q", want)
		}
		if strings.Contains(body, ">previous</a>") {
			t.Error("head page (from=0) rendered a previous link")
		}
	})

	t.Run("last-page-has-prev-no-next", func(t *testing.T) {
		body := get(t, "/workbench?from=100&entry=100")
		if got := strings.Count(body, "<h3>entry "); got != 5 {
			t.Errorf("last page entry count = %d, want 5", got)
		}
		for _, want := range []string{
			`href="/workbench?from=0&amp;entry=0">previous</a>`,
			`href="/workbench?from=100&amp;entry=104">select entry 104</a>`,
		} {
			if !strings.Contains(body, want) {
				t.Errorf("last page missing %q", want)
			}
		}
		if strings.Contains(body, ">next</a>") {
			t.Error("last page rendered a next link past the log end")
		}
	})

	t.Run("exactly-limit-no-next", func(t *testing.T) {
		body := get(t, "/workbench?from=5&entry=5")
		// CONTROL: the page is exactly full, so the old len==limit predicate holds here.
		if got := strings.Count(body, "<h3>entry "); got != workbench.WorkbenchPageLimit {
			t.Fatalf("control: entry count = %d, want %d", got, workbench.WorkbenchPageLimit)
		}
		if strings.Contains(body, ">next</a>") {
			t.Error("exactly-full last page rendered a next link to an empty page (F4)")
		}
		if want := `href="/workbench?from=0&amp;entry=0">previous</a>`; !strings.Contains(body, want) {
			t.Errorf("from=5 page missing clamped previous link %q", want)
		}
	})

	t.Run("emitted-links-resolve", func(t *testing.T) {
		counts := map[string]int{}
		ctx := context.Background()
		head, ok, err := d.store.SelectedHead(ctx)
		if err != nil || !ok {
			t.Fatalf("selected head ok=%v err=%v", ok, err)
		}
		headWorld, ok, err := d.store.GetWorld(ctx, head)
		if err != nil || !ok {
			t.Fatalf("head world ok=%v err=%v", ok, err)
		}
		if _, ok, err := d.store.GetObject(ctx, headWorld.StateRoot); err != nil || ok {
			t.Fatalf("head root stored=%v err=%v, want false", ok, err)
		}
		state := workbenchTestObject("paging-stored-state", "world/state", hashref.SumSHA256([]byte("paging-state-interface")))
		if err := d.store.PutObject(context.Background(), state); err != nil {
			t.Fatal(err)
		}
		storedWorld := store.World{Ref: hashref.SumSHA256([]byte("paging-stored-world")), Revision: 17, StateRoot: state.Hash, LogHead: hashref.SumSHA256([]byte("paging-stored-log"))}
		if err := d.store.PutWorld(storedWorld); err != nil {
			t.Fatal(err)
		}
		if _, ok, err := d.store.GetObject(ctx, state.Hash); err != nil || !ok {
			t.Fatalf("stored root control ok=%v err=%v", ok, err)
		}
		for _, page := range []struct {
			target      string
			hasSelected bool
			storedWorld bool
		}{
			{"/workbench", false, false},
			{"/workbench?from=100&entry=100", true, false},
			{"/workbench?from=5&entry=5", true, false},
			{"/workbench?world=" + storedWorld.Ref.String(), false, true},
		} {
			body := get(t, page.target)
			timeline, ok := workbenchRegion(body, timelineStart, "</section>")
			if !ok {
				t.Fatalf("%s: no timeline region", page.target)
			}
			regions := []struct{ name, text string }{{"timeline", timeline}}
			world, ok := workbenchRegion(body, worldStart, "</nav>")
			if !ok || world == "" {
				t.Fatalf("%s: no world region", page.target)
			}
			regions = append(regions, struct{ name, text string }{"world", world})
			root := headWorld.StateRoot
			if page.storedWorld {
				root = storedWorld.StateRoot
			}
			if page.storedWorld {
				if strings.Count(world, `?object=`) != 1 || !strings.Contains(world, `stateRoot: <a href="/workbench?object=`+root.String()+`"`) {
					t.Fatalf("%s: stored world root: %s", page.target, world)
				}
			} else if strings.Contains(world, `?object=`) || !strings.Contains(world, `UNAVAILABLE: object `+root.String()+` is not stored`) {
				t.Fatalf("%s: missing world root: %s", page.target, world)
			}
			selected, hasSelected := workbenchRegion(body, selectedEntryStart, "</article>")
			if hasSelected != page.hasSelected {
				t.Fatalf("%s: selected-entry region present=%v, want %v", page.target, hasSelected, page.hasSelected)
			}
			if hasSelected {
				regions = append(regions, struct{ name, text string }{"selected", selected})
			}
			for _, region := range regions {
				matches := workbenchLinkPattern.FindAllStringSubmatch(region.text, -1)
				// Nothing escapes classification: every anchor must be a pattern match.
				if anchors := strings.Count(region.text, "<a "); anchors != len(matches) {
					t.Fatalf("%s %s region: %d anchors but %d pattern matches", page.target, region.name, anchors, len(matches))
				}
				for _, match := range matches {
					href, text := match[1], match[2]
					category := ""
					switch {
					case region.name == "world" && href == "/workbench" && text == "workbench":
						category = "world-home"
					case region.name == "world" && strings.HasPrefix(href, "/workbench?object=") && text == storedWorld.StateRoot.String():
						category = "world"
					case region.name == "timeline" && (text == "previous" || text == "next"):
						category = "paging"
					case region.name == "timeline" && selectEntryText.MatchString(text):
						category = "select"
					case region.name == "selected" && strings.HasPrefix(href, "/workbench?object="):
						category = "stored-edge"
					default:
						t.Errorf("%s %s region: unclassified link href=%q text=%q", page.target, region.name, href, text)
						continue
					}
					counts[category]++
					target := strings.ReplaceAll(href, "&amp;", "&")
					if rec := requestRecorder(t, d, http.MethodGet, target, nil); rec.Code != http.StatusOK {
						t.Errorf("%s: %s link %q returned %d, want 200", page.target, category, target, rec.Code)
					}
				}
			}
		}
		// CONTROL: every category must be exercised, or a dead extractor passes vacuously.
		for _, category := range []string{"paging", "select", "stored-edge", "world-home", "world"} {
			if counts[category] == 0 {
				t.Errorf("category %s matched 0 links across the three pages", category)
			}
		}
	})
}

// denseLogStore answers every non-negative index with a copy of stored entry 0,
// so the from bound can be exercised without writing 2^63 entries.
type denseLogStore struct {
	readStore
}

func (s denseLogStore) GetLogEntry(ctx context.Context, index int64) (store.LogEntry, bool, error) {
	if index < 0 {
		return store.LogEntry{}, false, nil
	}
	entry, ok, err := s.readStore.GetLogEntry(ctx, 0)
	if err != nil || !ok {
		return entry, ok, err
	}
	entry.Header.EntryIndex = index
	return entry, true, nil
}

func TestWorkbenchNextLinkOverflowGuard(t *testing.T) {
	d := newHandlerDaemon(t)
	genesis := seedGenesisEmbedded(t, d, "workbench-overflow")
	if err := d.store.Commit(testCommit(genesis, 0, "workbench-overflow")); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	d.reads = denseLogStore{readStore: d.store}

	// MaxInt64-100 is the largest from the handler's from bound accepts; its next
	// link would carry from=MaxInt64, which that bound refuses, so none may render.
	rec := requestRecorder(t, d, http.MethodGet, "/workbench?from=9223372036854775707&entry=9223372036854775707", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body)
	}
	if strings.Contains(rec.Body.String(), ">next</a>") {
		t.Error("from=MaxInt64-100 rendered a next link whose own from overflows")
	}
	// CONTROL: one page earlier, the dense store does produce a next link.
	rec = requestRecorder(t, d, http.MethodGet, "/workbench?from=9223372036854775607&entry=9223372036854775607", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("control status = %d, want 200; body=%s", rec.Code, rec.Body)
	}
	if want := `href="/workbench?from=9223372036854775707&amp;entry=9223372036854775707">next</a>`; !strings.Contains(rec.Body.String(), want) {
		t.Errorf("control: missing next link %q", want)
	}
}

// probeFailingStore fails GetLogEntry at exactly one index, so a paging probe
// can be made to error while every timeline row read still succeeds.
type probeFailingStore struct {
	readStore
	failAt int64
}

func (s probeFailingStore) GetLogEntry(ctx context.Context, index int64) (store.LogEntry, bool, error) {
	if index == s.failAt {
		return store.LogEntry{}, false, errSentinelInternal
	}
	return s.readStore.GetLogEntry(ctx, index)
}

func TestWorkbenchPagingProbeStoreError(t *testing.T) {
	d := newHandlerDaemon(t)
	seedWorkbenchLog(t, d, 105)
	for _, tc := range []struct {
		name   string
		path   string
		failAt int64
	}{
		// /workbench reads rows 0-99; only the next probe reads entry 100.
		{"next-probe", "/workbench", 100},
		// from=5 reads rows 5-104 and probes 105 (absent); only the previous probe reads entry 0.
		{"prev-probe", "/workbench?from=5&entry=5", 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			oldReads, oldErrLog := d.reads, d.errLog
			defer func() { d.reads, d.errLog = oldReads, oldErrLog }()
			d.errLog = &bytes.Buffer{}
			// CONTROL: failing an index no read touches leaves the page at 200.
			d.reads = probeFailingStore{readStore: d.store, failAt: 1 << 40}
			if rec := requestRecorder(t, d, http.MethodGet, tc.path, nil); rec.Code != http.StatusOK {
				t.Fatalf("control: %s status = %d, want 200; body=%s", tc.path, rec.Code, rec.Body)
			}
			d.reads = probeFailingStore{readStore: d.store, failAt: tc.failAt}
			rec := requestRecorder(t, d, http.MethodGet, tc.path, nil)
			if rec.Code != http.StatusInternalServerError {
				t.Fatalf("%s with entry %d failing: status = %d, want 500; body=%s", tc.path, tc.failAt, rec.Code, rec.Body)
			}
			if !strings.Contains(rec.Body.String(), ">Internal<") {
				t.Errorf("body does not contain >Internal<: %s", rec.Body)
			}
		})
	}
}

const provenanceWalkStart = `<section aria-label="provenance walk">`

// provenanceWalkSection returns the provenance-walk section of a workbench page
// up to its </section>, failing if the section is missing.
func provenanceWalkSection(t *testing.T, body string) string {
	t.Helper()
	start := strings.Index(body, provenanceWalkStart)
	if start < 0 {
		t.Fatalf("no provenance-walk section in %s", body)
	}
	end := strings.Index(body[start:], "</section>\n</main>")
	if end < 0 {
		t.Fatalf("unterminated provenance-walk section in %s", body)
	}
	return body[start : start+end+len("</section>")]
}

// refFailingStore fails GetObject for one ref only, so the object read itself
// succeeds while the walk's existence check on that ref errors.
type refFailingStore struct {
	readStore
	fail hashref.HashRef
}

func (s refFailingStore) GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error) {
	if ref == s.fail {
		return store.Object{}, false, errSentinelInternal
	}
	return s.readStore.GetObject(ctx, ref)
}

func workbenchTestObject(label, semanticID string, iface hashref.HashRef) store.Object {
	payload := []byte("payload-" + label)
	return store.Object{Hash: hashref.SumSHA256(payload), InterfaceHash: iface, SemanticID: semanticID, Provenance: "workbench-test", Payload: payload}
}

func TestWorkbenchObjectGrade(t *testing.T) {
	d := newHandlerDaemon(t)
	genesis := seedGenesisEmbedded(t, d, "workbench-grade")
	commit := testCommit(genesis, 0, "workbench-grade")
	proof := workbenchTestObject("workbench-grade-proof", "world/proof-report/v1", hashref.SumSHA256([]byte("world/authenticated-proof-envelope/v1")))
	commit.Objects = append(commit.Objects, proof)
	if err := d.store.Commit(commit); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	// CONTROL: the registry object is the real bootstrap-written one, not a fixture.
	registry, ok, err := d.store.GetRegistryHead(context.Background(), store.EpochRegistryV1)
	if err != nil || !ok {
		t.Fatalf("control GetRegistryHead(%s): ok=%v err=%v", store.EpochRegistryV1, ok, err)
	}
	want := `<p>GRADE UNAVAILABLE — ` + objectGradeUnavailableReason + `</p>`
	for _, object := range []struct {
		name string
		ref  hashref.HashRef
	}{{"test-object", commit.Objects[0].Hash}, {"proof-labelled", proof.Hash}, {"registry", registry}} {
		for _, suffix := range []struct{ name, query string }{{"", ""}, {"-payload", "&payload=1"}} {
			t.Run(object.name+suffix.name, func(t *testing.T) {
				rec := requestRecorder(t, d, http.MethodGet, "/workbench?object="+object.ref.String()+suffix.query, nil)
				if rec.Code != http.StatusOK {
					t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body)
				}
				body := rec.Body.String()
				if got := strings.Count(body, want); got != 1 {
					t.Errorf("grade line %q appears %d times, want 1", want, got)
				}
				for _, unwanted := range []string{`GRADE UNAVAILABLE — </p>`, "no grade reason was supplied", `<span>PROVEN</span>`, `<span>TESTED</span>`, `<span>ATTESTED</span>`, `<span>CLAIMED</span>`} {
					if strings.Contains(body, unwanted) {
						t.Errorf("object page contains %q", unwanted)
					}
				}
			})
		}
	}
}

func TestWorkbenchObjectProvenanceWalk(t *testing.T) {
	d := newHandlerDaemon(t)
	genesis := seedGenesisEmbedded(t, d, "workbench-walk")
	commit := testCommit(genesis, 0, "workbench-walk")
	plain := workbenchTestObject("workbench-walk-plain", "test/plain", hashref.SumSHA256([]byte("interface-workbench-walk-plain")))
	schema := workbenchTestObject("workbench-walk-schema", "test/schema", hashref.SumSHA256([]byte("interface-workbench-walk-schema")))
	typed := workbenchTestObject("workbench-walk-typed", "test/typed", schema.Hash)
	for _, object := range []store.Object{plain, schema, typed} {
		if err := d.store.PutObject(context.Background(), object); err != nil {
			t.Fatal(err)
		}
	}
	if err := d.store.Commit(commit); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	// CONTROLS: one interface target is stored and one is not, or the link and
	// UNAVAILABLE arms below assert nothing.
	if _, ok, err := d.store.GetObject(context.Background(), schema.Hash); err != nil || !ok {
		t.Fatalf("control GetObject(schema): ok=%v err=%v, want ok=true", ok, err)
	}
	if _, ok, err := d.store.GetObject(context.Background(), plain.InterfaceHash); err != nil || ok {
		t.Fatalf("control GetObject(plain.InterfaceHash): ok=%v err=%v, want ok=false", ok, err)
	}
	get := func(t *testing.T, target string) string {
		t.Helper()
		rec := requestRecorder(t, d, http.MethodGet, target, nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status = %d, want 200; body=%s", target, rec.Code, rec.Body)
		}
		return rec.Body.String()
	}
	plainTarget := "/workbench?object=" + plain.Hash.String()
	typedTarget := "/workbench?object=" + typed.Hash.String()

	t.Run("interface-stored-link", func(t *testing.T) {
		section := provenanceWalkSection(t, get(t, typedTarget))
		href := "/workbench?object=" + schema.Hash.String()
		if want := `<p>interface: <a href="` + href + `"`; !strings.Contains(section, want) {
			t.Fatalf("walk missing stored interface link %q: %s", want, section)
		}
		get(t, href)
	})

	t.Run("interface-unstored-unavailable", func(t *testing.T) {
		body := get(t, plainTarget)
		iface := plain.InterfaceHash.String()
		want := `<p>interface: <span class="unavailable" role="note">UNAVAILABLE: object ` + iface + ` is not stored</span></p>`
		if section := provenanceWalkSection(t, body); !strings.Contains(section, want) {
			t.Errorf("walk missing %q: %s", want, section)
		}
		if unwanted := `href="/workbench?object=` + iface + `"`; strings.Contains(body, unwanted) {
			t.Errorf("unstored interface rendered as a link %q", unwanted)
		}
	})

	t.Run("named-stops", func(t *testing.T) {
		for _, target := range []string{plainTarget, typedTarget + "&payload=1"} {
			section := provenanceWalkSection(t, get(t, target))
			for _, want := range []string{
				`<h3>committedBy</h3>`,
				`no commit carried this object: it was stored outside any commit (PutObject or journal). Entries that only reference it are listed under referencedBy.`,
				`<h3>referencedBy</h3>`,
			} {
				if !strings.Contains(section, want) {
					t.Errorf("%s: walk missing named stop %q: %s", target, want, section)
				}
			}
			if strings.Contains(section, "no provenance edges were supplied") {
				t.Errorf("%s: daemon object page fell back to the render-layer stop: %s", target, section)
			}
		}
	})

	t.Run("edge-order", func(t *testing.T) {
		// The walk orders interface, committedBy, then referencedBy.
		for _, target := range []string{plainTarget, typedTarget} {
			section := provenanceWalkSection(t, get(t, target))
			last := -1
			for _, relation := range []string{"<p>interface: ", "<h3>committedBy</h3>", "<h3>referencedBy</h3>"} {
				if n := strings.Count(section, relation); n != 1 {
					t.Fatalf("%s: %q occurs %d times, want 1: %s", target, relation, n, section)
				}
				at := strings.Index(section, relation)
				if at <= last {
					t.Errorf("%s: %q is out of order: %s", target, relation, section)
				}
				last = at
			}
		}
	})

	t.Run("never-blank", func(t *testing.T) {
		for _, target := range []string{"/workbench", "/workbench?from=0&entry=0", plainTarget, typedTarget + "&payload=1"} {
			section := provenanceWalkSection(t, get(t, target))
			_, after, ok := strings.Cut(section, "</h2>")
			if !ok || strings.TrimSpace(after) == "" {
				t.Errorf("%s: provenance walk is a heading followed by nothing: %q", target, section)
			}
			if strings.Contains(section, "UNAVAILABLE: </span>") {
				t.Errorf("%s: walk rendered a stop with an empty reason: %s", target, section)
			}
		}
	})

	t.Run("interface-store-error", func(t *testing.T) {
		oldReads, oldErrLog := d.reads, d.errLog
		defer func() { d.reads, d.errLog = oldReads, oldErrLog }()
		d.reads = refFailingStore{readStore: oldReads, fail: schema.Hash}
		d.errLog = &bytes.Buffer{}
		// CONTROL: the same store serves a page whose walk does not read schema.
		if rec := requestRecorder(t, d, http.MethodGet, plainTarget, nil); rec.Code != http.StatusOK {
			t.Fatalf("control: %s status = %d, want 200; body=%s", plainTarget, rec.Code, rec.Body)
		}
		rec := requestRecorder(t, d, http.MethodGet, typedTarget, nil)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500; body=%s", rec.Code, rec.Body)
		}
		if !strings.Contains(rec.Body.String(), ">Internal<") {
			t.Errorf("body does not contain >Internal<: %s", rec.Body)
		}
	})
}
