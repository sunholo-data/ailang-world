package main

import (
	"strconv"
	"strings"
	"testing"
)

// Row 114 F-7: `why 198` on an ailang-edit showed `{"new_text":"// Row 185:
// IsLinked accepts…… (1048 bytes)` — the edited file's path sat past the
// cut, and a walk concluded the file had not been edited.
func TestWhyInputLineShowsThePathOfALongEdit(t *testing.T) {
	payload := []byte(`{"new_text":"` + strings.Repeat(`// Row 185: IsLinked accepts the worktree Add makes\n`, 20) +
		`","old_text":"}\n","path":"host/worktree/worktree_test.go"}`)
	got := whyInputLine(payload)
	if !strings.HasPrefix(got, `old_text="}\n" path="host/worktree/worktree_test.go" · {"new_text":`) {
		t.Fatalf("key fields not first:\n%s", got)
	}
	if !strings.HasSuffix(got, "("+strconv.Itoa(len(payload))+" bytes)") {
		t.Fatalf("the elided payload and its size must still follow:\n%s", got)
	}
	if strings.Contains(got, "new_text=") {
		t.Fatalf("a long field must not be repeated as a key field:\n%s", got)
	}
}

// A payload that fits is printed whole and byte-identical to before (no key
// fields), so short inputs — most of them — read exactly as they did.
func TestWhyInputLineLeavesAShortInputUnchanged(t *testing.T) {
	for _, p := range []string{`{"path":"cmd/ailang-worldd/setup.go"}`, `{"args":["-run","TestX","./cmd"],"command":"go-test"}`, `"x"`} {
		if got := whyInputLine([]byte(p)); got != elide(p) || got != p {
			t.Fatalf("short input changed: %q -> %q", p, got)
		}
	}
}

// Small arrays count as key fields (an exec call's args); nested objects and
// over-long values do not; at most whyKeyFieldCount are shown.
func TestWhyInputLineKeyFieldBounds(t *testing.T) {
	long := strings.Repeat("a", whyKeyFieldMax)
	payload := []byte(`{"a1":1,"a2":2,"a3":3,"a4":4,"a5":5,"a6":6,"a7":7,"args":["x"],"big":"` + long +
		`","nested":{"k":"v"}}`)
	got := whyInputLine(payload)
	head, _, ok := strings.Cut(got, " · ")
	if !ok || head != "a1=1 a2=2 a3=3 a4=4 a5=5 a6=6" {
		t.Fatalf("key fields = %q (want the first %d in key order)", head, whyKeyFieldCount)
	}
	small := []byte(`{"args":["-run","TestSessionNew","./cmd/ailang-worldd"],"command":"go-test","note":"` + long + `"}`)
	head, _, _ = strings.Cut(whyInputLine(small), " · ")
	if head != `args=["-run","TestSessionNew","./cmd/ailang-worldd"] command="go-test"` {
		t.Fatalf("array/scalar key fields = %q", head)
	}
}

// Not an object, or an object with no short field: plain elide, as before.
func TestWhyInputLineFallsBackToElide(t *testing.T) {
	for _, p := range []string{`"` + strings.Repeat("z", 300) + `"`, `{"blob":"` + strings.Repeat("z", 300) + `"}`, strings.Repeat("not json ", 40)} {
		if got := whyInputLine([]byte(p)); got != elide(p) {
			t.Fatalf("fallback changed: %q", got)
		}
	}
}
