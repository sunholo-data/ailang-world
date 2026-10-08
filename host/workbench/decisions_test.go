package workbench

import (
	"reflect"
	"strings"
	"testing"
)

func TestWorkbenchDecisionsReadOnly(t *testing.T) {
	p := Page{}
	v := reflect.ValueOf(&p).Elem().FieldByName("Decisions")
	if !v.IsValid() {
		t.Fatal("Page.Decisions missing")
	}
	body := renderPage(t, p)
	start := strings.Index(body, `<section aria-label="decisions">`)
	if start < 0 {
		t.Fatal("decisions section missing")
	}
	end := strings.Index(body[start:], "</section>")
	if end < 0 {
		t.Fatal("decisions section unclosed")
	}
	section := body[start : start+end]
	if !strings.Contains(section, "No approval requests recorded") {
		t.Fatal("positive empty state missing")
	}
	for _, bad := range []string{"<form", "<button", "method="} {
		if strings.Contains(section, bad) {
			t.Fatalf("action in decisions: %s", bad)
		}
	}
}
