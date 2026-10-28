package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMarketingCaptureInventoryMutation(t *testing.T) {
	d := newHandlerDaemon(t)
	c := workbenchLiveREST(t, d, []int64{0, 1, 2})
	srv := httptest.NewServer(d.Handler())
	defer srv.Close()
	dir := t.TempDir()
	bin := filepath.Join(dir, "marketingcheck")
	run := func(name string, argv ...string) (int, string) {
		t.Helper()
		ctx, cancel := context.WithTimeout(boundedTestContext(t), 60*time.Second)
		defer cancel()
		stdout, err := os.Create(filepath.Join(dir, name+".stdout"))
		if err != nil {
			t.Fatal(err)
		}
		defer stdout.Close()
		stderr, err := os.Create(filepath.Join(dir, name+".stderr"))
		if err != nil {
			t.Fatal(err)
		}
		defer stderr.Close()
		cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
		cmd.Stdout, cmd.Stderr = stdout, stderr
		rc := 0
		if err := cmd.Run(); err != nil {
			var exit *exec.ExitError
			if !errors.As(err, &exit) || ctx.Err() != nil {
				t.Fatalf("%s: %v (context: %v)", name, err, ctx.Err())
			}
			rc = exit.ExitCode()
		}
		raw, err := os.ReadFile(stderr.Name())
		if err != nil {
			t.Fatal(err)
		}
		return rc, string(raw)
	}
	if rc, detail := run("build", "go", "build", "-o", bin, "../../scripts/marketingcheck"); rc != 0 {
		t.Fatalf("build: rc %d: %s", rc, detail)
	}
	before, same, after := filepath.Join(dir, "before.json"), filepath.Join(dir, "same.json"), filepath.Join(dir, "after.json")
	read := func(name, path string) {
		t.Helper()
		if rc, detail := run(name, bin, "inventory", "--addr", srv.URL, "--out", path); rc != 0 {
			t.Fatalf("%s inventory: rc %d: %s", name, rc, detail)
		}
	}
	read("before", before)
	read("same", same)
	if rc, detail := run("equal", bin, "inventory-compare", "--before", before, "--after", same); rc != 0 {
		t.Fatalf("no-write capture inventory must compare equal: rc %d: %s", rc, detail)
	}
	extra := testCommit(c.NextWorld, 3, "workbench-live")
	r := requestRecorderAuth(t, d, authHeader(t, d), "POST", "/v1/commit", bytes.NewReader(encodeCommit(extra)))
	if r.Code != 200 {
		t.Fatalf("extra REST commit: %d %s", r.Code, r.Body)
	}
	read("after", after)
	if rc, detail := run("changed", bin, "inventory-compare", "--before", before, "--after", after); rc != 1 || !strings.Contains(detail, "marketingcheck: inventory-mutated") {
		t.Errorf("extra commit invalidates capture: rc %d: %s", rc, detail)
	}
	for path, count := range map[string]int{after: 4, before: 3} {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var inv struct{ Entries []json.RawMessage }
		if err := json.Unmarshal(raw, &inv); err != nil {
			t.Fatal(err)
		}
		if len(inv.Entries) != count {
			t.Errorf("%s must list all %d committed entries: got %d", filepath.Base(path), count, len(inv.Entries))
		}
	}
}
