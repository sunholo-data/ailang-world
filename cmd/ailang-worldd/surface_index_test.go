package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/surface"
)

// Row 188 M0a: `surface-index`. Building needs git, so these arms cover the
// flags, the refusals made before git runs, and --validate on the golden.

const surfaceGolden = "../../host/surface/testdata/index.json"

func runSurface(args ...string) (int, string, string) {
	var out, errw bytes.Buffer
	code := run(append([]string{"surface-index"}, args...), &out, &errw)
	return code, out.String(), errw.String()
}

func TestSurfaceIndexHelp(t *testing.T) {
	for _, args := range [][]string{{"help", "surface-index"}, {"surface-index", "--help"}} {
		var out, errw bytes.Buffer
		if code := run(args, &out, &errw); code != exitOK || out.String() != surfaceIndexHelp {
			t.Errorf("%v: exit %d, stdout %q, stderr %q", args, code, out.String(), errw.String())
		}
	}
	for _, flag := range []string{"--repo", "--charter", "--out", "--rev", "--limit", "--validate", "--rig", "--project", "--addr"} {
		if !strings.Contains(surfaceIndexHelp, "  "+flag+" ") {
			t.Errorf("help does not describe %s", flag)
		}
	}
	if !strings.Contains(usage, "\n  ailang-worldd surface-index --repo <dir> --charter <file> --out <file>") ||
		!strings.Contains(usage, "\n  ailang-worldd surface-index --validate <file>") ||
		!strings.Contains(usage, "[--rig <dir> --addr <url> [--project ailang-world]]") {
		t.Error("the top-level usage does not list surface-index and its flags")
	}
	for _, line := range []string{
		"ailang-worldd surface-index --repo . --charter c.md --out /tmp/i.json --rev dev --limit 50",
		"ailang-worldd surface-index --repo . --charter c.md --out /tmp/i.json --rig /rig --addr 127.0.0.1:7644 --project ailang-world",
	} {
		if err := checkWorlddLine(line); err != nil {
			t.Errorf("help check %q: %v", line, err)
		}
	}
}

// Row 188 M0c: the rig flags, refused before anything is read.
func TestSurfaceIndexRigFlagRefusals(t *testing.T) {
	repo := t.TempDir()
	outside := filepath.Join(t.TempDir(), "index.json")
	build := []string{"--repo", repo, "--charter", "c.md", "--out", outside}
	for _, tc := range []struct {
		args []string
		want string
	}{
		{[]string{"--rig", surfaceRig}, "--rig and --addr go together"},
		{[]string{"--addr", "127.0.0.1:7644"}, "--rig and --addr go together"},
		{[]string{"--project", "ailang-world"}, "--project needs --rig"},
		{[]string{"--rig", surfaceRig, "--addr", "127.0.0.1:7644", "--project", "../x"}, "is not a project name"},
		{[]string{"--rig", surfaceRig, "--addr", "example.com:7644"}, "is not loopback"},
		{[]string{"--rig", surfaceRig, "--addr", "http://10.0.0.1:7644"}, "is not loopback"},
		{[]string{"--rig", surfaceRig, "--addr", "http://127.0.0.2:7644"}, "is not loopback"},
		{[]string{"--rig", surfaceRig, "--addr", "http://127.0.0.1.example.com:7644"}, "is not loopback"},
		{[]string{"--rig", surfaceRig, "--addr", "http://u@127.0.0.1:7644"}, "is not a daemon base URL"},
		{[]string{"--rig", surfaceRig, "--addr", "http://127.0.0.1:7644/v1"}, "is not a daemon base URL"},
		{[]string{"--rig", surfaceRig, "--addr", "file:///tmp/x"}, "is not a daemon base URL"},
		{[]string{"--rig", filepath.Join(repo, "no-rig"), "--addr", "127.0.0.1:7644"}, "--charter"}, // the charter is read first
	} {
		code, stdout, stderr := runSurface(append(append([]string(nil), build...), tc.args...)...)
		if code != exitUsage || stdout != "" || !strings.Contains(stderr, tc.want) {
			t.Errorf("%v: exit %d, stdout %q, stderr %q; want %q", tc.args, code, stdout, stderr, tc.want)
		}
	}
	if code, _, stderr := runSurface("--validate", surfaceGolden, "--rig", surfaceRig, "--addr", "127.0.0.1:1"); code != exitUsage ||
		!strings.Contains(stderr, "--validate takes no other flag") {
		t.Errorf("--validate --rig: exit %d, %q", code, stderr)
	}
	for addr, want := range map[string]string{
		"127.0.0.1:7644": "http://127.0.0.1:7644", "localhost:1": "http://localhost:1",
		"http://[::1]:7644/": "http://[::1]:7644", "https://127.0.0.1": "https://127.0.0.1",
	} {
		if got, err := surfaceLoopbackBase(addr); err != nil || got != want {
			t.Errorf("surfaceLoopbackBase(%q) = %q, %v; want %q", addr, got, err, want)
		}
	}
}

const surfaceRig = "../../host/surface/testdata/rig"

func TestSurfaceIndexLoadRig(t *testing.T) {
	rig, err := loadSurfaceRig(surfaceRig, "ailang-world")
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, ep := range rig.Episodes {
		ids = append(ids, ep.ID)
	}
	if strings.Join(ids, ",") != "e1,e2,e3,e4,e5" || len(rig.Episodes[0].Rounds) != 2 ||
		rig.PRs["e1"] != 327 || rig.PRs["e2"] != 328 || len(rig.PRs) != 3 {
		t.Errorf("rig: episodes %v, e1 rounds %d, prs %v", ids, len(rig.Episodes[0].Rounds), rig.PRs)
	}
	other, err := loadSurfaceRig(surfaceRig, "platform-source")
	if err != nil || len(other.Episodes) != 1 || other.Episodes[0].ID != "p1" || other.PRs["p1"] != 71 {
		t.Errorf("platform-source: %+v, %v", other, err)
	}
	ix, err := surface.Decode([]byte(mustRead(t, surfaceGolden)))
	if err != nil {
		t.Fatal(err)
	}
	// Only the episodes whose PR is a step in the index have a report fetched.
	if got := surfaceReportEpisodes(ix, rig); strings.Join(got, ",") != "e1,e2" {
		t.Errorf("report episodes %v, want e1,e2", got)
	}
}

func TestSurfaceIndexRigBounds(t *testing.T) {
	writeRig := func(t *testing.T, n int, size int) string {
		t.Helper()
		rig := t.TempDir()
		dir := filepath.Join(rig, "evidence", "ailang-world", "episodes")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(rig, "trial"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(rig, "trial", "prs.json"), []byte("{}"), 0o644); err != nil {
			t.Fatal(err)
		}
		for i := 0; i < n; i++ {
			rec := fmt.Sprintf(`{"schema":"trial-episode/1","episode":"x%d","started_at":"2026-10-20T07:00:00Z","finished_at":"2026-10-20T07:00:00Z","pad":%q}`,
				i, strings.Repeat("p", size))
			if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("x%d.json", i)), []byte(rec), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return rig
	}
	if rig, err := loadSurfaceRig(writeRig(t, 1, 0), "ailang-world"); err != nil || len(rig.Episodes) != 1 {
		t.Errorf("one small record: %+v, %v", rig, err)
	}
	if _, err := loadSurfaceRig(writeRig(t, 1, surfaceRecordMaxBytes), "ailang-world"); err == nil || !strings.Contains(err.Error(), "over the bound of 1048576 bytes") {
		t.Errorf("an oversized record: %v", err)
	}
	if _, err := loadSurfaceRig(writeRig(t, surfaceRecordMaxFiles+1, 0), "ailang-world"); err == nil || !strings.Contains(err.Error(), "more than 2000 episode records") {
		t.Errorf("too many records: %v", err)
	}
	if _, err := loadSurfaceRig(writeRig(t, surfaceRecordMaxFiles, 0), "ailang-world"); err != nil {
		t.Errorf("exactly %d records: %v", surfaceRecordMaxFiles, err)
	}
	noPRs := writeRig(t, 1, 0)
	_ = os.Remove(filepath.Join(noPRs, "trial", "prs.json"))
	if _, err := loadSurfaceRig(noPRs, "ailang-world"); err == nil {
		t.Error("a rig with no prs.json loaded")
	}
}

func TestSurfaceIndexFetchReports(t *testing.T) {
	report := mustRead(t, filepath.Join(surfaceRig, "reports", "e1.json"))
	const base = "http://127.0.0.1:7644"
	var hosts []string
	prev := surfaceReportTransport
	surfaceReportTransport = handlerTransport(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hosts = append(hosts, r.URL.Host)
		switch r.URL.Path {
		case "/v1/episodes/e1/report":
			_, _ = io.WriteString(w, report)
		case "/v1/episodes/big/report":
			_, _ = w.Write(bytes.Repeat([]byte(" "), surfaceReportMaxBytes+1))
		case "/v1/episodes/broken/report":
			w.WriteHeader(http.StatusInternalServerError)
		case "/v1/episodes/moved/report":
			http.Redirect(w, r, "http://example.com/", http.StatusFound)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer func() { surfaceReportTransport = prev }()
	ctx := t.Context()
	got, err := fetchSurfaceReports(ctx, base, []string{"e1", "e2"})
	if err != nil {
		t.Fatal(err)
	}
	if rep, ok := got["e1"]; !ok || rep.TotalCalls != 42 || len(rep.Reads.Files) != 6 || rep.PlanRefusals.Count != 2 {
		t.Errorf("e1's report: %+v", rep)
	}
	if _, ok := got["e2"]; ok || len(got) != 1 {
		t.Errorf("a 404 gave a report: %v", got)
	}
	for ep, want := range map[string]string{"big": "over 4194304 bytes", "broken": "HTTP 500", "moved": "redirect refused"} {
		if _, err := fetchSurfaceReports(ctx, base, []string{ep}); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: err %v, want %q", ep, err, want)
		}
	}
	// The redirect to example.com was never followed.
	for _, h := range hosts {
		if h != "127.0.0.1:7644" {
			t.Errorf("a request went to %s", h)
		}
	}
}

// handlerTransport serves requests from a handler, in process.
func handlerTransport(h http.Handler) http.RoundTripper { return handlerRoundTripper{h} }

type handlerRoundTripper struct{ h http.Handler }

func (rt handlerRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	rt.h.ServeHTTP(rec, r)
	return rec.Result(), nil
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

func TestSurfaceIndexFlagRefusals(t *testing.T) {
	repo := t.TempDir()
	outside := filepath.Join(t.TempDir(), "index.json")
	for _, tc := range []struct {
		args []string
		want string
	}{
		{nil, "--repo, --charter and --out are required"},
		{[]string{"--repo", repo, "--charter", "c.md"}, "--repo, --charter and --out are required"},
		{[]string{"--validate", surfaceGolden, "--repo", repo}, "--validate takes no other flag"},
		{[]string{"--validate", surfaceGolden, "--rev", "dev"}, "--validate takes no other flag"},
		{[]string{"--repo", repo, "--charter", "c.md", "--out", outside, "--limit", "0"}, "--limit 0 is outside 1..2500"},
		{[]string{"--repo", repo, "--charter", "c.md", "--out", outside, "--limit", "2501"}, "--limit 2501 is outside 1..2500"},
		{[]string{"--repo", repo, "--charter", "c.md", "--out", outside, "--rev", "--upload-pack=x"}, "is not a revision"},
		{[]string{"--repo", repo, "--charter", "c.md", "--out", outside, "extra"}, `unexpected argument "extra"`},
		{[]string{"--repo", repo, "--charter", "c.md", "--out", filepath.Join(repo, "sub", "index.json")}, "is inside --repo"},
		{[]string{"--repo", repo, "--charter", "c.md", "--out", repo}, "is inside --repo"},
		{[]string{"--repo", filepath.Join(repo, "nope"), "--charter", "c.md", "--out", outside}, "--repo"},
		{[]string{"--repo", repo, "--charter", filepath.Join(repo, "nope.md"), "--out", outside}, "--charter"},
	} {
		code, stdout, stderr := runSurface(tc.args...)
		if code != exitUsage || stdout != "" || !strings.Contains(stderr, "refused: ") || !strings.Contains(stderr, tc.want) {
			t.Errorf("%v: exit %d, stdout %q, stderr %q; want %q", tc.args, code, stdout, stderr, tc.want)
		}
	}
	if _, err := os.Stat(outside); err == nil {
		t.Error("a refused build wrote --out")
	}
	var out, errw bytes.Buffer
	if code := run([]string{"--addr", "127.0.0.1:1", "surface-index", "--validate", surfaceGolden}, &out, &errw); code != exitUsage {
		t.Errorf("--addr surface-index: exit %d", code)
	}
}

func TestSurfaceIndexValidate(t *testing.T) {
	code, stdout, stderr := runSurface("--validate", surfaceGolden)
	if code != exitOK || !strings.Contains(stdout, "24 shared states, 24 steps, 6 branes") {
		t.Fatalf("golden: exit %d, stdout %q, stderr %q", code, stdout, stderr)
	}
	raw, err := os.ReadFile(surfaceGolden)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	for name, tc := range map[string]struct{ from, to, want string }{
		"duplicate id": {`"id": "f36428abd9a72e5093b61cba14e716c37019f1ef"`, `"id": "90035557b2012895551867c0ebe2f36b853aac09"`, "duplicates"},
		"bad kind":     {`"kind": "loop"`, `"kind": "nightly"`, `kind "nightly"`},
		"schema":       {surfaceSchemaLine, `"schema": "world/surface-index/v0"`, "schema"},
		"unknown key":  {surfaceSchemaLine, surfaceSchemaLine + `, "path": "x"`, "does not decode"},
	} {
		if !strings.Contains(string(raw), tc.from) {
			t.Fatalf("%s: the golden has no %q", name, tc.from)
		}
		p := filepath.Join(dir, strings.ReplaceAll(name, " ", "-")+".json")
		if err := os.WriteFile(p, []byte(strings.Replace(string(raw), tc.from, tc.to, 1)), 0o644); err != nil {
			t.Fatal(err)
		}
		code, stdout, stderr := runSurface("--validate", p)
		if code == exitOK || stdout != "" || !strings.Contains(stderr, tc.want) {
			t.Errorf("%s: exit %d, stdout %q, stderr %q; want %q", name, code, stdout, stderr, tc.want)
		}
	}
	if code, _, stderr := runSurface("--validate", filepath.Join(dir, "missing.json")); code == exitOK || !strings.Contains(stderr, "missing.json") {
		t.Errorf("missing file: exit %d, %q", code, stderr)
	}
}

const surfaceSchemaLine = `"schema": "world/surface-index/v1"`

func TestWriteFileAtomic(t *testing.T) {
	p := filepath.Join(t.TempDir(), "d", "index.json")
	for _, body := range []string{"one\n", "two\n"} {
		if err := writeFileAtomic(p, []byte(body)); err != nil {
			t.Fatal(err)
		}
		if got, _ := os.ReadFile(p); string(got) != body {
			t.Errorf("read %q, want %q", got, body)
		}
	}
	ents, _ := os.ReadDir(filepath.Dir(p))
	if len(ents) != 1 {
		t.Errorf("temporary files left: %v", ents)
	}
}
