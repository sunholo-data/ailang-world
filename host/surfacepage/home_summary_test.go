package surfacepage

import (
	"reflect"
	"testing"
)

// The start page's card shows Needs you first, and only what the Now card shows.
func TestBuildHomeSummary(t *testing.T) {
	healthy := &NowEdge{Health: "World healthy", Sessions: []ModelSession{{}, {}}, Today: Today{Landed: 3}}
	follow := &NowEdge{Health: "World healthy", Open: &ModelOpen{Live: []ModelRunning{{}}}, Today: Today{Landed: 5},
		Service: &ModelService{Warn: false}}
	for _, c := range []struct {
		name  string
		needs int
		edge  *NowEdge
		want  []HomeItem
	}{
		{"no model, unreadable requests", -1, nil, []HomeItem{{"Needs you", "unknown", "attention"}}},
		{"no model, nothing pending", 0, nil, []HomeItem{{"Needs you", "0", ""}}},
		{"index mode counts sessions", 2, healthy, []HomeItem{{"Needs you", "2", "attention"}, {"Running", "2", ""},
			{"Landed today", "3", ""}, {"Health", "healthy", "ok"}}},
		{"follow mode counts open runs", 0, follow, []HomeItem{{"Needs you", "0", ""}, {"Running", "1", ""},
			{"Landed today", "5", ""}, {"Health", "healthy", "ok"}}},
		{"unreadable sessions", 0, &NowEdge{Health: "World healthy", SessionsUnavailable: "x"}, []HomeItem{{"Needs you", "0", ""},
			{"Running", "unknown", ""}, {"Landed today", "0", ""}, {"Health", "healthy", "ok"}}},
		{"stale build warns", 0, &NowEdge{Health: "World healthy", Service: &ModelService{Warn: true}}, []HomeItem{{"Needs you", "0", ""},
			{"Running", "0", ""}, {"Landed today", "0", ""}, {"Health", "check", "attention"}}},
		{"unhealthy", 0, &NowEdge{Health: "World not serving"}, []HomeItem{{"Needs you", "0", ""},
			{"Running", "0", ""}, {"Landed today", "0", ""}, {"Health", "check", "attention"}}},
	} {
		got := BuildHomeSummary(c.needs, c.edge)
		if got.Schema != HomeSummarySchema || !reflect.DeepEqual(got.Items, c.want) {
			t.Errorf("%s: got %+v, want %+v", c.name, got, c.want)
		}
		if len(got.Items) > 4 {
			t.Errorf("%s: %d items; the start page shows at most 4", c.name, len(got.Items))
		}
	}
}
