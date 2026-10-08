package daemon

import (
	"context"
	"fmt"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"strings"
	"testing"
)

func workbenchNamedSection(t *testing.T, body, label string) string {
	t.Helper()
	start := strings.Index(body, `<section aria-label="`+label+`">`)
	if start < 0 {
		t.Fatalf("%s region missing", label)
	}
	end := strings.Index(body[start:], "</section>")
	if end < 0 {
		t.Fatal("section unclosed")
	}
	return body[start : start+end+10]
}

type graphReads struct {
	readStore
	objects int
}

func (r *graphReads) GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error) {
	r.objects++
	return r.readStore.GetObject(ctx, ref)
}
func TestWorkbenchGraphWiring(t *testing.T) {
	d := newHandlerDaemon(t)
	c := workbenchLiveREST(t, d, []int64{0, 1, 2, 3, 4, 5, 6})
	spy := &graphReads{readStore: d.reads}
	d.reads = spy
	s := workbenchNamedSection(t, workbenchLiveBody(t, d, "/workbench"), "world graph")
	for i := 2; i <= 6; i++ {
		if !strings.Contains(s, fmt.Sprintf(">entry %d</text>", i)) {
			t.Fatalf("newest entry %d missing", i)
		}
	}
	if strings.Contains(s, ">entry 1</text>") || strings.Count(s, "<line ") != 15 || spy.objects > 16 {
		t.Fatalf("home graph edges/reads incorrect: %d", spy.objects)
	}
	spy.objects = 0
	s = workbenchNamedSection(t, workbenchLiveBody(t, d, "/workbench?from=0&entry=1"), "world graph")
	if strings.Count(s, "<line ") != 3 || !strings.Contains(s, ">entry 1</text>") || strings.Contains(s, ">entry 6</text>") || spy.objects != 4 {
		t.Fatalf("selected graph or reused reads incorrect: %d", spy.objects)
	}
	for _, relation := range []string{"transitionFn", "interpreter", "transitionRef"} {
		if !strings.Contains(s, "<title>"+relation+"</title>") {
			t.Fatal("relation missing")
		}
	}
	for _, relation := range []string{"stateRoot", "committedBy", "evidence"} {
		if strings.Contains(s, relation) {
			t.Fatal("excluded graph relation")
		}
	}
	s = workbenchNamedSection(t, workbenchLiveBody(t, d, "/workbench?object="+c.Objects[0].Hash.String()), "world graph")
	if !strings.Contains(s, "graph: select an entry") || strings.Contains(s, "<line ") {
		t.Fatal("object graph placeholder missing")
	}
}

func TestWorkbenchGraphStoreError(t *testing.T) {
	d := newHandlerDaemon(t)
	c := workbenchLiveREST(t, d, []int64{0})
	d.reads = refFailingStore{readStore: d.reads, fail: c.Entry.TransitionRef}
	s := workbenchNamedSection(t, workbenchLiveBody(t, d, "/workbench"), "world graph")
	if !strings.Contains(s, "UNAVAILABLE: graph checked edges could not be read") || strings.Contains(s, "<line ") || strings.Contains(s, "private") {
		t.Fatal("home graph failure must be explicit and sanitized")
	}
	if r := requestRecorder(t, d, "GET", "/workbench?from=0&entry=0", nil); r.Code != 500 {
		t.Fatal("selected-entry error posture changed")
	}
}
