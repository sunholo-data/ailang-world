package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
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
		chrome := fakeCaptureChrome(t, dir)
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
			argv := strings.Split(strings.TrimSpace(readFile(t, filepath.Join(dir, "pass-"+row.Name+".args"))), "\n")
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

	t.Run("stop-first", func(t *testing.T) {
		dir := t.TempDir()
		calls := filepath.Join(dir, "calls")
		t.Setenv("CAPTURE_CALLS", calls)
		// graph.png (the fourth row) fails; nothing after it may be attempted.
		chrome := fakeScript(t, dir, "chrome", `echo x >> "$CAPTURE_CALLS"
case " $* " in *" --window-size=1440,2100 "*) exit 5 ;; esac
`+screenshotPath+writePNG)
		rows, err := ShotRoutes(demoFixture())
		pngOK(t, err)
		routes := putJSON(t, dir, "routes.json", rows)
		outDir := filepath.Join(dir, "pass")
		got, se := runCapture(t, "NonzeroExit", "pass stops at the first non-ok shot", "pass", "--routes", routes,
			"--base-url", "http://example.invalid", "--chrome", chrome, "--out-dir", outDir)
		entries, err := os.ReadDir(outDir)
		pngOK(t, err)
		if len(got.Shots) != 4 || got.Shots[3].Name != "graph.png" || got.Shots[3].Verdict != "NonzeroExit" ||
			strings.Count(readFile(t, calls), "\n") != 4 || len(entries) != 3 ||
			!strings.HasSuffix(se, "progress: 4/6 graph.png NonzeroExit\n") {
			t.Errorf("pass stops at the first non-ok shot: shots=%+v calls=%q files=%d stderr=%q",
				got.Shots, readFile(t, calls), len(entries), se)
		}
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

func fakeCaptureChrome(t *testing.T, dir string) string {
	return fakeScript(t, dir, "chrome", `if [ "${1:-}" = --version ]; then
 if [ "${CAPTURE_CHANGE:-}" = yes ] && [ -f "$CAPTURE_RECORD/version" ]; then
  echo "Google Chrome 156.0.8078.47"
 else
  echo "${CAPTURE_VERSION:-Google Chrome 155.0.8059.41}"
 fi
 touch "$CAPTURE_RECORD/version"
 exit 0
fi
scheme=light
for arg do
 case "$arg" in
 --screenshot=*) png=${arg#--screenshot=} ;;
 --window-size=*) size=${arg#--window-size=} ;;
 --blink-settings=preferredColorScheme=0) scheme=dark ;;
 esac
done
printf '%s\n' "$@" > "$CAPTURE_RECORD/$(basename "$(dirname "$png")")-$(basename "$png").args"
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
}

func TestMarketingCaptureFakeChromeSixShots(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "marketingcheck")
	if rc, _, se := runCaptureCommand(t, "go", "build", "-o", bin, "."); rc != 0 {
		t.Fatal(se)
	}
	bs, rects := fixtureSet(t)
	fixtures := writeSet(t, bs)
	manifest := putJSON(t, dir, "demo.json", demoFixture())
	rectangles := putJSON(t, dir, "rectangles.json", rects)
	for _, name := range []string{"six-shots", "serve-argv", "version-changed", "old-chrome", "cleanup", "scheme", "bind-refused"} {
		t.Run(name, func(t *testing.T) {
			d := t.TempDir()
			t.Setenv("CAPTURE_FIXTURES", fixtures)
			t.Setenv("CAPTURE_RECORD", d)
			if name == "version-changed" {
				t.Setenv("CAPTURE_CHANGE", "yes")
			}
			if name == "old-chrome" {
				t.Setenv("CAPTURE_VERSION", "Google Chrome 154.0.1.2")
			}
			bind := "127.0.0.1:1"
			if name == "bind-refused" {
				bind = "0.0.0.0:1"
			} else if name != "old-chrome" {
				listener, err := net.Listen("tcp", "127.0.0.1:0")
				pngOK(t, err)
				bind = listener.Addr().String()
				pngOK(t, listener.Close())
			}
			worldd := filepath.Join(d, "worldd")
			pngOK(t, os.WriteFile(worldd, []byte(`#!/usr/bin/env python3
import http.server, json, os, pathlib, sys
root = pathlib.Path(__file__).parent
(root / 'argv').write_text(json.dumps(sys.argv[1:]) + '\n')
(root / 'pid').write_text(str(os.getpid()))
host, port = sys.argv[sys.argv.index('--bind') + 1].split(':')
class Handler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200 if self.path in ('/v1/health', '/v1/head') else 404)
        self.end_headers()
        self.wfile.write(b'{}')
http.server.HTTPServer((host, int(port)), Handler).serve_forever()
`), 0700))
			out := filepath.Join(d, "out")
			tool := filepath.Join(d, "operator tool")
			args := []string{"../capture_marketing.sh", "--db", filepath.Join(d, "db"), "--manifest", manifest,
				"--proof-key-file", filepath.Join(d, "key"), "--ailang-bin", filepath.Join(d, "ailang"),
				"--workspace-root", d, "--tool-ailang-bin", tool, "--chrome", fakeCaptureChrome(t, d),
				"--out", out, "--rectangles", rectangles, "--worldd-bin", worldd, "--bind", bind, "--marketingcheck-bin", bin, "--verify-twice"}
			rc, _, se := runCaptureCommand(t, "bash", args...)
			if name == "bind-refused" {
				_, err := os.Stat(filepath.Join(d, "argv"))
				if rc == 0 || !os.IsNotExist(err) || !strings.Contains(se, bind) {
					t.Fatalf("non-loopback bind must refuse before daemon starts: rc=%d stderr=%q argv stat=%v", rc, se, err)
				}
				return
			}
			if name == "old-chrome" {
				shots, _ := filepath.Glob(filepath.Join(d, "*.args"))
				_, err := os.Stat(filepath.Join(d, "argv"))
				if rc != 1 || len(shots) != 0 || !os.IsNotExist(err) || !strings.Contains(se, "chrome-version A") {
					t.Fatalf("old Chrome must refuse before any shot: rc=%d stderr=%q", rc, se)
				}
				return
			}
			if name == "version-changed" {
				if rc != 1 || !strings.Contains(se, "chrome version changed during capture") {
					t.Fatalf("a Chrome version change during capture must refuse: rc=%d stderr=%q", rc, se)
				}
				return
			}
			if name == "serve-argv" {
				var argv []string
				pngOK(t, json.Unmarshal([]byte(readFile(t, filepath.Join(d, "argv"))), &argv))
				paired := func(flag, value string) bool {
					for i := 0; i+1 < len(argv); i++ {
						if argv[i] == flag && argv[i+1] == value {
							return true
						}
					}
					return false
				}
				if !paired("--proof-required-identity", "sameRef") || !paired("--proof-required-identity", "servesEntry") {
					t.Errorf("both seed identities must be configured: %q", argv)
				}
				if len(argv) == 0 || argv[0] != "serve" || !slices.Contains(argv, "--workbench-still") || !paired("--tool-ailang-bin", tool) {
					t.Errorf("serve must forward frozen configuration and operator tool path: %q", argv)
				}
				for _, forbidden := range []string{"commit", "session", "mint", "approve"} {
					if slices.Contains(argv, forbidden) {
						t.Errorf("capture must not mutate: %q", argv)
					}
				}
			}
			if name == "scheme" {
				for _, pass := range []string{"pass-a", "pass-b"} {
					for _, shot := range Shots {
						argv := strings.Split(strings.TrimSpace(readFile(t, filepath.Join(d, pass+"-"+shot.Name+".args"))), "\n")
						scheme := "1"
						if shot.Scheme == "dark" {
							scheme = "0"
						}
						if !slices.Contains(argv, "--blink-settings=preferredColorScheme="+scheme) {
							t.Errorf("every shot explicitly pins correct scheme: %s/%s argv=%q", pass, shot.Name, argv)
						}
					}
				}
			}
			if name == "cleanup" {
				var pid int
				_, err := fmt.Sscan(readFile(t, filepath.Join(d, "pid")), &pid)
				pngOK(t, err)
				assertDead(t, pid, "fake worldd after capture")
			}
			if rc != 0 {
				t.Fatalf("capture must succeed: rc=%d stderr=%q", rc, se)
			}
			if name == "six-shots" {
				for _, pass := range []string{"pass-a", "pass-b"} {
					entries, err := os.ReadDir(filepath.Join(out, pass))
					pngOK(t, err)
					if len(entries) != 6 {
						t.Fatalf("each pass must contain six names: %v", entries)
					}
				}
				for _, shot := range Shots {
					if !bytes.Equal([]byte(readFile(t, filepath.Join(out, "pass-a", shot.Name))), []byte(readFile(t, filepath.Join(out, "pass-b", shot.Name)))) {
						t.Errorf("passes must be byte-equal: %s", shot.Name)
					}
				}
				hashes, err := CheckPass(filepath.Join(out, "pass-a"), rects)
				pngOK(t, err)
				var capture CaptureManifest
				pngOK(t, json.Unmarshal([]byte(readFile(t, filepath.Join(out, "capture-manifest.json"))), &capture))
				pngOK(t, VerifyCaptureManifest(capture, hashes))
				diagnostics := readFile(t, filepath.Join(out, "diagnostics.txt"))
				if !strings.Contains(diagnostics, "chrome-version: Google Chrome 155.0.8059.41\n") || !strings.Contains(diagnostics, "bind: "+bind+"\n") {
					t.Errorf("diagnostics must record version and bind: %q", diagnostics)
				}
			}
		})
	}
	t.Run("daemon-readiness-config", func(t *testing.T) {
		for _, flags := range [][]string{nil, {"--ready-file", filepath.Join(dir, "ready"), "--ready-url", "http://127.0.0.1:1"}} {
			args := append([]string{"../marketing_capture.py", "daemon", "--started-file", filepath.Join(dir, "started")}, flags...)
			args = append(args, "--", "does-not-exist")
			rc, _, se := runCaptureCommand(t, "python3", args...)
			if rc != 2 || !strings.Contains(se, "--ready") {
				t.Errorf("exactly one readiness source required: rc=%d stderr=%q", rc, se)
			}
		}
	})
}
