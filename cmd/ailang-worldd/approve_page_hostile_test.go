package main

// Row 201 A3 sprint 4 (iteration 386 judge survivors C03/C04): a hostile signed value on the request
// card, in a live Chrome, is shown as text.

import (
	"testing"
	"time"
)

// A request whose model is markup and an attribute breaker: the card shows it verbatim as the
// model field's text, ready to approve; nothing it carries becomes an element or runs.
func TestApprovePageShowsAHostileValueAsText(t *testing.T) {
	const hostile = `<img src=x onerror="window.__xss=1"><script>window.__xss=2</script>"&amp;`
	b := startPasskeyBrowser(t)
	b.do(t, nil, "tab", nil)
	a := newApproveRig(t, "ep1", time.Now().Unix(), "--model", hostile)
	origin := cardOrigin(t, a.credPath, "", nil)
	var view cardView
	if b.do(t, &view, "openCard", map[string]any{"url": origin + "/surface/approve/" + a.id}); view.State != "ready" || view.Fields["model"] != hostile ||
		len(view.Errors) != 0 || len(view.Violations) != 0 {
		t.Fatalf("the card shows model %q, state %q, errors %q, CSP violations %q", view.Fields["model"], view.State, view.Errors, view.Violations)
	}
	var page struct {
		XSS             bool
		Images, Scripts int
	}
	b.do(t, &page, "waitFor", map[string]any{"expression": `({ xss: window.__xss !== undefined, images: document.querySelectorAll('img').length, scripts: document.scripts.length })`})
	if page.XSS || page.Images != 0 || page.Scripts != 1 {
		t.Fatalf("the hostile value became markup: %+v", page)
	}
}
