package main

// Row 243 V2a (plan §3.1, VC1): the why paths no other test drove, pinned to the bytes the CLI
// printed before the walk moved to host/why (the golden: WORLD_UPDATE_WHY_GOLDEN=1 on the base tree).

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const whyGolden = "testdata/why-characterise.golden.jsonl"

func TestWhyCharacterisedPaths(t *testing.T) {
	r := newWalkRig(t)
	inv := r.invoke("ep1", "ailang-read", "task-1", "hello\n")
	h, err := (worldReader{ctx: context.Background(), c: newClient(r.url)}).health()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	result := filepath.Join(dir, "result.json")
	writeFileT(t, result, inv.outputText+"\n")
	type run struct {
		Name, Stdout, Stderr string
		Code                 int
	}
	var got []run
	for _, c := range []struct {
		name, stdin string
		args        []string
	}{
		{"no target", "", nil},
		{"malformed sha256", "", []string{"sha256:zz"}},
		{"absent sha256", "", []string{"sha256:" + strings.Repeat("0", 64)}},
		{"flag before target", "", []string{"--json", "1"}},
		{"two targets", "", []string{"1", "2"}},
		{"target and --result", "", []string{"1", "--result", result}},
		{"--result absent file", "", []string{"--result", filepath.Join(dir, "absent")}},
		{"rest: unknown receipt", "", []string{"rest:nobody"}},
		{"--expect-store this store", "", []string{"1", "--expect-store", h.DBPath}},
		{"- not a result", "x", []string{"-", "--scan", "3"}},
		{"index past head", "", []string{"9"}},
	} {
		code, out, errOut := whyRun(t, r.url, c.stdin, c.args...)
		got = append(got, run{c.name, strings.ReplaceAll(out, dir, "$TMP"), strings.ReplaceAll(errOut, dir, "$TMP"), code})
	}
	var b []byte
	for _, g := range got { // one JSON line per run
		line, err := json.Marshal(g)
		if err != nil {
			t.Fatal(err)
		}
		b = append(append(b, line...), '\n')
	}
	if os.Getenv("WORLD_UPDATE_WHY_GOLDEN") == "1" {
		writeFileT(t, whyGolden, string(b))
	}
	want, err := os.ReadFile(whyGolden)
	if err != nil || string(want) != string(b) {
		t.Fatalf("why's output moved from the base's (%v):\n%s", err, b)
	}
}
