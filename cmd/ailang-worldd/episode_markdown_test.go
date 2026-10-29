package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
)

func TestFinishSharedScrubVectors(t *testing.T) {
	data, err := os.ReadFile("../../scripts/trial/testdata/scrub_vectors.json")
	if err != nil {
		t.Fatal(err)
	}
	var vectors []struct {
		Name   string
		Text   string
		Refuse bool
	}
	if err := json.Unmarshal(data, &vectors); err != nil || len(vectors) == 0 {
		t.Fatalf("vectors: %v", err)
	}
	for _, v := range vectors {
		t.Run(v.Name, func(t *testing.T) {
			if got := scrubMarkdown(v.Text) != nil; got != v.Refuse {
				t.Errorf("refused = %v, want %v", got, v.Refuse)
			}
		})
	}
}

func ac5Report() daemon.EpisodeReport {
	r := sampleReport()
	r.Episode = "ep1"
	r.Entries = [2]int64{1, 6000}
	r.Start = &store.EpisodeStart{Agent: "codex", Model: "fixture-model"}
	r.Writes.Written = []daemon.ReportWrite{{Path: "a.go", Op: "write", Count: 1, FirstEntry: 12}}
	r.Writes.Produced = []daemon.ReportRead{{Path: "gen.txt", Count: 1, FirstEntry: 14}}
	return r
}

func TestFinishMarkdownOrderCapAndAC5(t *testing.T) {
	r := ac5Report()
	body, err := renderReportMarkdown(r, "/fixture/world.db", reportAdditions{})
	if err != nil {
		t.Fatal(err)
	}
	last := -1
	for _, want := range []string{"Episode ep1: codex fixture-model", "## What the agent did", "## Files", "a.go (write)", "gen.txt", "World entries 1-6000", "World-Provenance:"} {
		i := strings.Index(body, want)
		if i <= last {
			t.Fatalf("missing/out-of-order %q: %s", want, body)
		}
		last = i
	}
	trailer := "World-Provenance: store=/fixture/world.db episode=ep1 entries=1-6000"
	if !strings.HasSuffix(body, trailer+"\n") {
		t.Fatal(body)
	}
	t.Log("AC5 report --markdown:\n" + body)
	r.CallsByTool = map[string]int{strings.Repeat("large ", 15000): 1}
	body, err = renderReportMarkdown(r, "/fixture/world.db", reportAdditions{})
	if err != nil || len(body) > 60000 || !strings.Contains(body, "truncated; run ailang-worldd episode report ep1") || !strings.HasSuffix(body, trailer+"\n") {
		t.Fatalf("cap: len %d, err %v", len(body), err)
	}
	// Scrub the full source even when the unsafe value would be truncated.
	r.Writes.Produced = append(r.Writes.Produced, daemon.ReportRead{Path: strings.Repeat("a", 64)})
	if _, err := renderReportMarkdown(r, "/fixture/world.db", reportAdditions{}); err == nil {
		t.Fatal("truncated secret escaped scrub")
	}
}

func TestFinishMarkdownCLIRefusalsLoopback(t *testing.T) {
	for _, unsafe := range []string{strings.Repeat("a", 64), "Bearer x", "sha256:" + strings.Repeat("a", 64)} {
		r := ac5Report()
		r.PlanRefusals.First[0].Reason = unsafe
		c := memoryDaemon(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/v1/health" {
				json.NewEncoder(w).Encode(map[string]any{"db_path": "/fixture/world.db"})
				return
			}
			json.NewEncoder(w).Encode(r)
		})
		var out, errw bytes.Buffer
		ctx, cancel := budgetContext(whyBudget)
		code := reportOnce(ctx, c, reportOpts{episode: "ep1", markdown: true}, &out, &errw)
		cancel()
		if strings.HasPrefix(unsafe, "sha256:") {
			if code != 0 || !strings.Contains(out.String(), unsafe) {
				t.Fatalf("prefixed: %d %s", code, errw.String())
			}
		} else if code != 1 || out.Len() != 0 || !strings.Contains(errw.String(), "refusing markdown") {
			t.Fatalf("scrub: %d stdout %q stderr %q", code, out.String(), errw.String())
		}
	}
	for _, extra := range []string{"--json", "--follow"} {
		var out, errw bytes.Buffer
		if rc := runEpisodeReport("127.0.0.1:1", []string{"ep1", "--markdown", extra}, &out, &errw); rc != 1 || out.Len() != 0 || !strings.Contains(errw.String(), "cannot combine") {
			t.Fatalf("flags: %d %s", rc, errw.String())
		}
	}
}
