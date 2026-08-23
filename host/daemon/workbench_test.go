package daemon

import (
	"context"
	"net/http"
	"strings"
	"testing"
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
