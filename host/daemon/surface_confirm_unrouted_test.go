package daemon

// Row 243 V3 sprint 2 (design §2.5 [r1], §8): the confirmation card's page and its data exist only in
// host/surfacepage and tests; no daemon route serves them before sprint 3, which replaces this test
// with the routes' own (CR1).

import (
	"os"
	"strings"
	"testing"
)

// The Handler registers nothing under /surface/confirm, and the daemon answers 404 for the card, its
// script and its data. Control: the decision card's page, registered beside them, answers 200.
func TestConfirmRoutesNotRegistered(t *testing.T) {
	src, err := os.ReadFile("daemon.go")
	if err != nil {
		t.Fatal(err)
	}
	body := strings.SplitN(strings.SplitN(string(src), "func (d *Daemon) Handler() http.Handler {", 2)[1], "\n}\n", 2)[0]
	if strings.Contains(body, "/surface/confirm") || !strings.Contains(body, `"GET /surface/decide/{ref}"`) {
		t.Fatal("Handler registers a /surface/confirm route (or lost the decide control)")
	}
	a := newApprovalRig(t)
	for _, target := range []string{"/surface/confirm/vr-0123456789abcdef", "/surface/confirm.js", "/surface/confirmations/vr-0123456789abcdef"} {
		if rec := a.serve("GET", target, apLoopback, nil); rec.Code != 404 {
			t.Errorf("GET %s = %d, want 404 (no route before sprint 3)", target, rec.Code)
		}
	}
	if rec := a.serve("GET", "/surface/decide/sha256:"+strings.Repeat("ab", 32), apLoopback, nil); rec.Code != 200 {
		t.Fatalf("control: GET /surface/decide/… = %d", rec.Code)
	}
}
