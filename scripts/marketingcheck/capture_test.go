package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// File sinks keep descendants from holding output pipes open beyond the deadline.
func runCaptureCommand(t *testing.T, bin string, args ...string) (int, string, string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	dir := t.TempDir()
	stdout, err := os.Create(filepath.Join(dir, "stdout"))
	pngOK(t, err)
	defer stdout.Close()
	stderr, err := os.Create(filepath.Join(dir, "stderr"))
	pngOK(t, err)
	defer stderr.Close()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Stdout, cmd.Stderr = stdout, stderr
	rc := 0
	if err := cmd.Run(); err != nil {
		var ee *exec.ExitError
		if !errors.As(err, &ee) {
			t.Fatal(err)
		}
		rc = ee.ExitCode()
	}
	return rc, readFile(t, stdout.Name()), readFile(t, stderr.Name())
}

type captureVerdict struct {
	Verdict string `json:"verdict"`
	Detail  string `json:"detail"`
	Version string `json:"version"`
	Shots   []struct {
		Name    string `json:"name"`
		Verdict string `json:"verdict"`
	} `json:"shots"`
}

func runCapture(t *testing.T, want, assertion string, args ...string) (captureVerdict, string) {
	t.Helper()
	py, err := exec.LookPath("python3")
	pngOK(t, err)
	rc, out, se := runCaptureCommand(t, py, append([]string{"../marketing_capture.py"}, args...)...)
	var got captureVerdict
	if err := json.Unmarshal([]byte(out), &got); err != nil || strings.Count(out, "\n") != 1 {
		t.Fatalf("one JSON verdict line required: %v stdout=%q stderr=%q", err, out, se)
	}
	wantRC := 1
	if want == "ok" {
		wantRC = 0
	}
	if got.Verdict != want || rc != wantRC {
		t.Errorf("%s: want %s rc=%d, got %+v rc=%d stderr=%q", assertion, want, wantRC, got, rc, se)
	}
	return got, se
}

func TestMarketingCapturePassFakeChrome(t *testing.T) {
	t.Run("routes", func(t *testing.T) {
		dir := t.TempDir()
		bin := filepath.Join(dir, "marketingcheck")
		if rc, out, se := runCaptureCommand(t, "go", "build", "-o", bin, "."); rc != 0 {
			t.Fatalf("CLI build rc=%d stdout=%s stderr=%s", rc, out, se)
		}
		demo := demoFixture()
		manifest := putJSON(t, dir, "demo.json", demo)
		rc, out, se := runCaptureCommand(t, bin, "routes", "--manifest", manifest)
		if rc != 0 || se != "" {
			t.Fatalf("routes CLI rc=%d stderr=%s", rc, se)
		}
		var rows []RouteRow
		pngOK(t, json.Unmarshal([]byte(out), &rows))
		if len(rows) != 6 {
			t.Fatalf("routes must have exactly six rows: %+v", rows)
		}
		m, err := BuildCaptureManifest(demo, map[string]string{})
		pngOK(t, err)
		literalHeights := map[string]int{"packet.png": 3700, "graph.png": 2100, "provenance.png": 3000}
		for i, row := range rows {
			s := Shots[i]
			if row.Name != s.Name || row.Width != s.Width || row.Height != s.Height || row.Scheme != s.Scheme ||
				(literalHeights[row.Name] != 0 && row.Height != literalHeights[row.Name]) {
				t.Errorf("routes must carry the M7c0 shot dimensions: got %+v want %+v", row, s)
			}
			if row.Route != m.Shots[i].Route {
				t.Errorf("routes must equal the capture-manifest routes: %s got %q want %q", row.Name, row.Route, m.Shots[i].Route)
			}
		}
		for _, args := range [][]string{{"routes"}, {"routes", "--unknown"}, {"routes", "--manifest", manifest, "extra"}} {
			if rc, _, _ := runCaptureCommand(t, bin, args...); rc != 2 {
				t.Errorf("routes usage errors exit 2: %q rc=%d", args, rc)
			}
		}
		bad := putJSON(t, dir, "bad.json", DemoManifest{})
		if rc, _, se := runCaptureCommand(t, bin, "routes", "--manifest", bad); rc != 1 ||
			!strings.HasPrefix(se, "marketingcheck: bad-manifest demo: ") {
			t.Errorf("routes bad-manifest refusal: rc=%d stderr=%q", rc, se)
		}
	})

	t.Run("pass-six", func(t *testing.T) {
		dir := t.TempDir()
		bs, rects := fixtureSet(t)
		fixtures := writeSet(t, bs)
		t.Setenv("CAPTURE_FIXTURES", fixtures)
		t.Setenv("CAPTURE_RECORD", dir)
		chrome := fakeScript(t, dir, "chrome", `scheme=light
for arg do
 case "$arg" in
 --screenshot=*) png=${arg#--screenshot=} ;;
 --window-size=*) size=${arg#--window-size=} ;;
 --blink-settings=preferredColorScheme=0) scheme=dark ;;
 esac
done
printf '%s\n' "$@" > "$CAPTURE_RECORD/$(basename "$png").args"
case "$size" in
 1440,1100) fixture=workbench-$scheme.png ;;
 1440,3700) fixture=packet.png ;;
 1440,2100) fixture=graph.png ;;
 1440,3000) fixture=provenance.png ;;
 390,844) fixture=phone.png ;;
 *) exit 4 ;;
esac
cp "$CAPTURE_FIXTURES/$fixture" "$png"
`)
		rows, err := ShotRoutes(demoFixture())
		pngOK(t, err)
		routes := putJSON(t, dir, "routes.json", rows)
		outDir := filepath.Join(dir, "pass")
		got, se := runCapture(t, "ok", "six-shot pass succeeds", "pass", "--routes", routes,
			"--base-url", "http://example.invalid", "--chrome", chrome, "--out-dir", outDir)
		if len(got.Shots) != 6 {
			t.Fatalf("pass must report six shots: %+v", got)
		}
		entries, err := os.ReadDir(outDir)
		pngOK(t, err)
		if len(entries) != 6 {
			t.Errorf("pass must create exactly six files: %v", entries)
		}
		for i, row := range rows {
			if got.Shots[i].Name != row.Name || got.Shots[i].Verdict != "ok" {
				t.Errorf("pass reports each shot in routes order: %+v", got.Shots[i])
			}
			argv := strings.Split(strings.TrimSpace(readFile(t, filepath.Join(dir, row.Name+".args"))), "\n")
			scheme := "1"
			if row.Scheme == "dark" {
				scheme = "0"
			}
			if !slices.Contains(argv, "--blink-settings=preferredColorScheme="+scheme) {
				t.Errorf("every shot explicitly pins correct scheme: %s argv=%q", row.Name, argv)
			}
			if !slices.Contains(argv, fmt.Sprintf("--window-size=%d,%d", row.Width, row.Height)) ||
				argv[len(argv)-1] != "http://example.invalid"+row.Route {
				t.Errorf("pass must forward the row viewport and route: %s argv=%q", row.Name, argv)
			}
		}
		if _, err := CheckPass(outDir, rects); err != nil {
			t.Errorf("CheckPass must accept captured fixtures: %v", err)
		}
		var progress []string
		for i, row := range rows {
			progress = append(progress, fmt.Sprintf("progress: %d/6 %s ok\n", i+1, row.Name))
		}
		if se != strings.Join(progress, "") {
			t.Errorf("progress goes only to stderr in shot order: %q", se)
		}
	})

	t.Run("version", func(t *testing.T) {
		for _, tc := range []struct{ version, want, assertion string }{
			{"Google Chrome 156.0.8078.47", "ok", "Chrome 156 is accepted"},
			{"Google Chrome 155.0.8059.41", "ok", "Chrome 155 is accepted"},
			{"Google Chrome 154.0.1.2", "ChromeVersion", "Chrome 154 must be refused"},
			{"garbage", "ChromeVersion", "unparsable Chrome version is refused"},
			{"exit-3", "NonzeroExit", "nonzero version command is refused"},
		} {
			t.Setenv("CAPTURE_VERSION", tc.version)
			chrome := fakeScript(t, t.TempDir(), "chrome", `[ "$#" -eq 1 ] && [ "$1" = --version ] || exit 4
[ "$CAPTURE_VERSION" != exit-3 ] || exit 3
printf '%s\n' "$CAPTURE_VERSION"
`)
			got, _ := runCapture(t, tc.want, tc.assertion, "chrome-version", "--chrome", chrome)
			if tc.want == "ok" && got.Version != tc.version {
				t.Errorf("full Chrome version must be recorded: got %q want %q", got.Version, tc.version)
			}
		}
	})

	t.Run("ready", func(t *testing.T) {
		var healthCalls, headCalls atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/v1/health":
				if healthCalls.Add(1) <= 2 {
					w.WriteHeader(http.StatusServiceUnavailable)
				}
			case "/v1/head":
				headCalls.Add(1)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		defer server.Close()
		runCapture(t, "ok", "readiness waits for health and head", "ready", "--base-url", server.URL)
		if healthCalls.Load() < 3 || headCalls.Load() == 0 {
			t.Errorf("readiness must poll health and then head: health=%d head=%d", healthCalls.Load(), headCalls.Load())
		}
		unready := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer unready.Close()
		runCapture(t, "Timeout", "readiness deadline refuses unavailable server", "ready", "--base-url", unready.URL, "--startup-timeout", "1")
	})

	t.Run("fresh-dir", func(t *testing.T) {
		dir := t.TempDir()
		rows, err := ShotRoutes(demoFixture())
		pngOK(t, err)
		routes := putJSON(t, dir, "routes.json", rows)
		chrome := fakeScript(t, dir, "chrome", "exit 9\n")
		runCapture(t, "DestinationExists", "pass destination must be fresh", "pass", "--routes", routes,
			"--base-url", "http://example.invalid", "--chrome", chrome, "--out-dir", dir)
	})
}
