package daemon

import (
	"bytes"
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
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

	// The observable is the log entry hash, NOT the committed object hash. The
	// object hash reaches an EntryView only through TransitionRef, and the WB.B
	// template renders no TransitionRef action at all (see the sprint plan's
	// s7d) -- so asserting on the object hash here would be asserting on a
	// value this route cannot emit. The entry hash is written by the timeline
	// loop and by nothing else on the page, which is what makes it a pin on the
	// mechanism rather than on a sibling channel.
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
