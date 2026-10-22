package daemon

import (
	"fmt"
	"html"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
)

func TestWorkbenchStillAllRegions(t *testing.T) {
	d := newHandlerDaemon(t)
	publishHumanTool(t, d)
	res, packet := producePacket(t, d, humanPlan("still"), "still-task")
	p, err := coordinator.DecodeDecisionPacket(packet.Payload)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := broker.DecideApproval(boundedTestContext(t), d.store, p.RequestRef, "deny", "still-operator", packetNow+1); err != nil {
		t.Fatal(err)
	}
	ctx := boundedTestContext(t)
	headBefore, ok, err := d.store.SelectedHead(ctx)
	if err != nil || !ok {
		t.Fatalf("selected head: ok=%v err=%v", ok, err)
	}
	logBefore, err := d.store.LogEntriesLatest(ctx, 1)
	if err != nil || len(logBefore) != 1 {
		t.Fatalf("log head: %v, err=%v", logBefore, err)
	}
	paths := []string{"/workbench", fmt.Sprintf("/workbench?from=0&entry=%d", res.EntryIndex), "/workbench?object=" + packet.Hash.String()}
	defaults := make([]*httptest.ResponseRecorder, len(paths))
	for i, path := range paths {
		defaults[i] = requestRecorder(t, d, "GET", path, nil)
		if defaults[i].Code != 200 {
			t.Fatalf("default GET %s: %d", path, defaults[i].Code)
		}
	}
	d.cfg.WorkbenchStill = true
	first := make([]string, len(paths))
	for i, path := range paths {
		r := requestRecorder(t, d, "GET", path, nil)
		if r.Code != 200 {
			t.Fatalf("still GET %s: %d", path, r.Code)
		}
		body := r.Body.String()
		first[i] = body
		if got := strings.Count(body, "<script"); got != 0 {
			t.Fatalf("still page must have zero script tags: got %d on %s", got, path)
		}
		if !reflect.DeepEqual(r.Header(), defaults[i].Header()) {
			t.Fatalf("still response headers differ from default on %s: %v != %v", path, r.Header(), defaults[i].Header())
		}
		for _, header := range []string{"Content-Security-Policy", "Content-Type", "Cache-Control", "X-Content-Type-Options", "Referrer-Policy"} {
			if r.Header().Get(header) == "" {
				t.Fatalf("missing response header %s", header)
			}
		}
		for _, want := range []string{"<header><h1>AILANG World Workbench</h1>", `aria-label="world browser"`, headBefore.String(), "Live updates: off. Reload to refresh."} {
			if !strings.Contains(body, want) {
				t.Fatalf("still page missing %q on %s", want, path)
			}
		}
		if s := workbenchNamedSection(t, body, "timeline"); !strings.Contains(s, fmt.Sprintf("<h3>entry %d</h3>", res.EntryIndex)) {
			t.Fatal("timeline missing recorded entry")
		}
		if s := workbenchLiveRegion(t, body); !strings.Contains(s, `<ul class="live-list">`) || !strings.Contains(s, logBefore[0].EntryHash.String()) {
			t.Fatal("live list missing recorded log head")
		}
		graph := workbenchNamedSection(t, body, "world graph")
		if i < 2 && (!strings.Contains(graph, fmt.Sprintf(">entry %d</text>", res.EntryIndex)) || !strings.Contains(graph, "<line ") || !strings.Contains(graph, "<title>transitionRef</title>")) {
			t.Fatal("world graph missing recorded entry edges")
		}
		if s := workbenchNamedSection(t, body, "decisions"); !strings.Contains(s, p.RequestRef.String()) || !strings.Contains(s, "deny by still-operator") {
			t.Fatal("decisions pane missing recorded decision")
		}
		if want := strings.Replace(defaults[i].Body.String(), `<script src="/workbench/live.js" defer></script>`, "", 1); body != want {
			t.Fatalf("still mode changed HTML beyond the script on %s", path)
		}
	}
	sec := packetSection(t, first[2])
	for _, want := range []string{html.EscapeString(fmt.Sprintf(packetValidatedFormat, res.EntryIndex)), `href="/workbench?object=` + p.RequestRef.String() + `"`, `href="/workbench?object=` + p.ProposalHash.String() + `"`} {
		if !strings.Contains(sec, want) {
			t.Fatalf("decision packet missing validated content/ref %q", want)
		}
	}
	if s := workbenchNamedSection(t, first[2], "inspector"); !strings.Contains(s, packet.Hash.String()) || !strings.Contains(s, coordinator.DecisionPacketV1) {
		t.Fatal("object pane missing stored packet")
	}
	// The daemon has no injectable request clock. Sleep after all first GETs,
	// so every repeated request is separated from its first by at least 2 s.
	start := time.Now()
	time.Sleep(2 * time.Second)
	for i, path := range paths {
		body := workbenchLiveBody(t, d, path)
		if body != first[i] {
			at := 0
			for at < len(body) && at < len(first[i]) && body[at] == first[i][at] {
				at++
			}
			t.Fatalf("repeat GET %s varies at byte %d: first %q, repeat %q", path, at, first[i][max(0, at-80):min(len(first[i]), at+160)], body[max(0, at-80):min(len(body), at+160)])
		}
	}
	t.Logf("repeat GETs byte-identical on all three routes after %s (real sleep)", time.Since(start))
	headAfter, ok, err := d.store.SelectedHead(ctx)
	if err != nil || !ok || headAfter != headBefore {
		t.Fatalf("GET changed selected head: %v -> %v, ok=%v err=%v", headBefore, headAfter, ok, err)
	}
	logAfter, err := d.store.LogEntriesLatest(ctx, 1)
	if err != nil || !reflect.DeepEqual(logBefore, logAfter) {
		t.Fatalf("GET changed log head: %v -> %v, err=%v", logBefore, logAfter, err)
	}
}
