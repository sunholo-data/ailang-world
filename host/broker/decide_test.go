package broker

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

func TestSketchRows(t *testing.T) {
	capability := func(effect, scope string, expiresAt, budget int64) Capability {
		return Capability{Effect: effect, Scope: scope, ExpiresAt: expiresAt, Budget: budget}
	}
	request := func(effect, scope string, cost, now int64) EffectRequest {
		return EffectRequest{Effect: effect, Scope: scope, Cost: cost, Now: now}
	}
	boolRows := []struct {
		name string
		line int
		got  bool
		want bool
	}{
		{"effectNameMatches/1", 47, effectNameMatches(capability("FS.Read", "/p", 10, 1), "FS.Read"), true},
		{"effectNameMatches/2", 48, effectNameMatches(capability("FS.Read", "/p", 10, 1), "FS.Write"), false},
		{"scopeMatches/1", 58, scopeMatches(capability("e", "/project/src", 10, 1), "/project/src"), true},
		{"scopeMatches/2", 59, scopeMatches(capability("e", "/project/src", 10, 1), "/project"), false},
		{"capabilityLive/1", 70, capabilityLive(capability("e", "s", 10, 1), 9), true},
		{"capabilityLive/2", 71, capabilityLive(capability("e", "s", 10, 1), 10), false},
		{"capabilityLive/3", 72, capabilityLive(capability("e", "s", 10, 1), -1), false},
		{"withinEffectBudget/1", 83, withinEffectBudget(5, 5), true},
		{"withinEffectBudget/2", 84, withinEffectBudget(5, 0), true},
		{"withinEffectBudget/3", 85, withinEffectBudget(5, 6), false},
		{"withinEffectBudget/4", 86, withinEffectBudget(-1, 0), false},
		{"withinEffectBudget/5", 87, withinEffectBudget(5, -1), false},
		{"effectAllowed/1", 122, effectAllowed(capability("FS.Read", "/p", 10, 5), request("FS.Read", "/p", 2, 3)), true},
		{"effectAllowed/2", 124, effectAllowed(capability("FS.Read", "/p", 10, 5), request("FS.Read", "/p", 2, 10)), false},
	}
	for _, row := range boolRows {
		t.Run(fmt.Sprintf("line_%d_%s", row.line, row.name), func(t *testing.T) {
			if row.got != row.want {
				t.Fatalf("sketch line %d: got %v, want %v", row.line, row.got, row.want)
			}
		})
	}

	debitRows := []struct {
		line         int
		budget, cost int64
		want         int64
	}{{102, 5, 2, 3}, {103, 5, 5, 0}, {104, 5, 0, 5}}
	for _, row := range debitRows {
		t.Run(fmt.Sprintf("line_%d_debit", row.line), func(t *testing.T) {
			if got := debit(row.budget, row.cost); got != row.want {
				t.Fatalf("sketch line %d: got %d, want %d", row.line, got, row.want)
			}
		})
	}

	decisionRows := []struct {
		line int
		c    Capability
		r    EffectRequest
		want string
	}{
		{225, capability("FS.Read", "/p", 10, 5), request("FS.Read", "/p", 2, 3), "allowed:3"},
		{227, capability("FS.Read", "/p", 10, 5), request("Git.Commit", "/p", 2, 3), "denied:effect-name"},
		{229, capability("FS.Read", "/p", 10, 5), request("FS.Read", "/q", 2, 3), "denied:scope"},
		{231, capability("FS.Read", "/p", 10, 5), request("FS.Read", "/p", 2, 10), "denied:expired"},
		{233, capability("FS.Read", "/p", 10, 5), request("FS.Read", "/p", 6, 3), "denied:budget"},
	}
	for _, row := range decisionRows {
		t.Run(fmt.Sprintf("line_%d_decideLabel", row.line), func(t *testing.T) {
			if got := Decide(row.c, row.r).Label; got != row.want {
				t.Fatalf("sketch line %d: got %q, want %q", row.line, got, row.want)
			}
		})
	}

	recordRows := []struct {
		line int
		rec  EffectRecord
		want bool
	}{
		{175, EffectRecord{Cost: 2, BudgetBefore: 5, BudgetAfter: 3, Allowed: true, ResultRef: hashref.SumSHA256([]byte("bb"))}, true},
		{179, EffectRecord{Cost: 2, BudgetBefore: 5, BudgetAfter: 5, Allowed: true, ResultRef: hashref.SumSHA256([]byte("bb"))}, false},
		{183, EffectRecord{Cost: 2, BudgetBefore: 5, BudgetAfter: 3, Allowed: true}, false},
		{187, EffectRecord{Cost: 2, BudgetBefore: 5, BudgetAfter: 3, Allowed: true, Failed: true}, true},
		{191, EffectRecord{Cost: 2, BudgetBefore: 5, BudgetAfter: 5, Allowed: true, Failed: true}, false},
		{195, EffectRecord{Cost: 2, BudgetBefore: 5, BudgetAfter: 3, Allowed: true, Failed: true, ResultRef: hashref.SumSHA256([]byte("bb"))}, false},
		{199, EffectRecord{Cost: 2, BudgetBefore: 5, BudgetAfter: 5, Denial: "budget"}, true},
		{203, EffectRecord{Cost: 2, BudgetBefore: 5, BudgetAfter: 3, Denial: "budget"}, false},
		{207, EffectRecord{Cost: 2, BudgetBefore: 5, BudgetAfter: 5, Failed: true, Denial: "budget"}, false},
	}
	for _, row := range recordRows {
		t.Run(fmt.Sprintf("line_%d_recordConsistent", row.line), func(t *testing.T) {
			if got := RecordConsistent(row.rec); got != row.want {
				t.Fatalf("sketch line %d: got %v, want %v", row.line, got, row.want)
			}
		})
	}
}

func TestCanonicalLabelSet(t *testing.T) {
	got := []string{
		allowedLabel(0),
		LabelDeniedEffectName,
		LabelDeniedScope,
		LabelDeniedExpired,
		LabelDeniedBudget,
	}
	want := []string{"allowed:0", "denied:effect-name", "denied:scope", "denied:expired", "denied:budget"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("canonical label set = %#v, want %#v", got, want)
	}
}
