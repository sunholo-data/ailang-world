package broker

import (
	"reflect"
	"strings"
	"testing"
)

func TestFinishApprovalScopeStrict(t *testing.T) {
	scope := FinishApprovalScope{
		Repository: "o/r", Remote: "git@github.com:o/r.git", Head: "world/ep1",
		SHA: strings.Repeat("a", 40), Base: "dev", TitleSHA256: strings.Repeat("b", 64),
		BodySHA256: strings.Repeat("c", 64), Draft: true, CreateLabel: true, ExpiresAt: 1700000900,
	}
	raw := FormatFinishApprovalScope(scope)
	parsed, err := ParseFinishApprovalScope(raw)
	if err != nil || !reflect.DeepEqual(parsed, scope) {
		t.Fatalf("roundtrip: %+v %v", parsed, err)
	}
	for _, bad := range []string{
		strings.Replace(raw, "&head=world/ep1&sha="+scope.SHA, "&sha="+scope.SHA+"&head=world/ep1", 1),
		strings.Replace(raw, "&base=dev", "", 1),
		strings.Replace(raw, "#remote=", "##remote=", 1),
		strings.Replace(raw, "expires=1700000900", "expires=tomorrow", 1),
		strings.Replace(raw, "expires=1700000900", "expires=-1", 1),
		raw + "&unknown=true",
		strings.Replace(raw, "draft=true", "draft=1", 1),
	} {
		if _, err := ParseFinishApprovalScope(bad); err == nil {
			t.Errorf("accepted %s", bad)
		}
	}
}

func TestFinishMatchProtected(t *testing.T) {
	patterns, err := ParseProtectedPatterns([]string{"*_test.go"})
	if err != nil {
		t.Fatal(err)
	}
	if !MatchProtected(patterns, "a_test.go") || MatchProtected(patterns, "a.go") {
		t.Fatal("protected match changed")
	}
}
