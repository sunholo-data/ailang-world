package daemon

import (
	"bytes"
	"context"
	"net/http"
	"regexp"
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
)

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
		d.reads = objectFailingStore{readStore: d.store}
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
		for _, page := range []struct {
			target      string
			hasSelected bool
		}{
			{"/workbench", false},
			{"/workbench?from=100&entry=100", true},
			{"/workbench?from=5&entry=5", true},
		} {
			body := get(t, page.target)
			timeline, ok := workbenchRegion(body, timelineStart, "</section>")
			if !ok {
				t.Fatalf("%s: no timeline region", page.target)
			}
			regions := []struct{ name, text string }{{"timeline", timeline}}
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
		for _, category := range []string{"paging", "select", "stored-edge"} {
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
