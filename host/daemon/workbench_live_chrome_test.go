package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// Instrument-health control only: the request trace below proves browser behavior.
func requireWorkbenchLiveQuickstart(t *testing.T) {
	t.Helper()
	b, err := os.ReadFile("../../docs/QUICKSTART.md")
	if err != nil {
		t.Fatal(err)
	}
	text := string(b)
	heading := "## 11. Open the live workbench"
	start := strings.Index(text, heading)
	if start < 0 {
		t.Fatal("missing QUICKSTART heading: Open the live workbench")
	}
	section := text[start:]
	for _, want := range []string{"http://127.0.0.1:7644/workbench", "Live list", "Footer", "Graph", "Decisions", "JavaScript disabled", "16", "hidden", "without reload", "observedHead", "nextWorld", "entryIndex", "base64", "--session", "genesis", "second"} {
		if !strings.Contains(section, want) {
			t.Errorf("live quickstart missing %q", want)
		}
	}
	b, err = os.ReadFile("daemon.go")
	if err != nil {
		t.Fatal(err)
	}
	text = string(b)
	end := strings.Index(text, "func (d *Daemon) Handler()")
	start = strings.LastIndex(text[:end], "// Handler")
	if start < 0 || !strings.Contains(text[start:end], "GET /workbench/live.js") {
		t.Error("Handler route documentation missing GET /workbench/live.js")
	}
}

func TestWorkbenchLiveQuickstartControl(t *testing.T) { requireWorkbenchLiveQuickstart(t) }

type workbenchChromeTrace struct {
	method string
	cursor int64
	err    error
}

func startWorkbenchChrome(t *testing.T, chrome, url string, jsOff bool) {
	t.Helper()
	root := t.TempDir()
	for _, dir := range []string{"cache", "tmp", "profile"} {
		if err := os.Mkdir(filepath.Join(root, dir), 0700); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	args := []string{"--headless=new", "--no-first-run", "--no-default-browser-check", "--disable-background-networking", "--user-data-dir=" + filepath.Join(root, "profile")}
	if jsOff {
		args = append(args, "--blink-settings=scriptEnabled=false")
	}
	args = append(args, url)
	cmd := exec.CommandContext(ctx, chrome, args...)
	for _, e := range os.Environ() {
		key := strings.SplitN(e, "=", 2)[0]
		// HOME is deliberately NOT overridden: on macOS headless Chrome with a
		// relocated HOME never loads its URL (measured iteration 245: zero
		// requests with HOME set, one without). --user-data-dir isolates the
		// profile instead.
		if key != "XDG_CACHE_HOME" && key != "TMPDIR" && key != "TMP" && key != "TEMP" {
			cmd.Env = append(cmd.Env, e)
		}
	}
	cmd.Env = append(cmd.Env, "XDG_CACHE_HOME="+filepath.Join(root, "cache"), "TMPDIR="+filepath.Join(root, "tmp"), "TMP="+filepath.Join(root, "tmp"), "TEMP="+filepath.Join(root, "tmp"))
	logPath := filepath.Join(root, "chrome.log")
	log, err := os.Create(logPath)
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	cmd.Stdout, cmd.Stderr = log, log
	if err := cmd.Start(); err != nil {
		cancel()
		log.Close()
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	t.Cleanup(func() {
		cancel()
		<-done
		log.Close()
		if t.Failed() {
			b, _ := os.ReadFile(logPath)
			t.Logf("Chrome log: %s", b)
		}
	})
}

func waitWorkbenchChromeTrace(t *testing.T, events <-chan workbenchChromeTrace, method string, cursor int64) {
	t.Helper()
	guard := time.NewTimer(20 * time.Second)
	defer guard.Stop()
	for {
		select {
		case e := <-events:
			if e.err != nil {
				t.Fatal(e.err)
			}
			t.Logf("trace %s /%s lastIndex=%d", e.method, map[string]string{"GET": "workbench", "POST": "agui/"}[e.method], e.cursor)
			if e.method == method && (method != "POST" || e.cursor == cursor) {
				return
			}
		case <-guard.C:
			t.Fatalf("waiting for %s cursor %d", method, cursor)
		}
	}
}

func TestWorkbenchLiveChromeDrill(t *testing.T) {
	// PD5: fail on missing docs before the environment gate, never launch for red-first.
	requireWorkbenchLiveQuickstart(t)
	if t.Failed() {
		return
	}
	chrome := os.Getenv("WORLD_CHROME")
	if chrome == "" {
		t.Skip("WORLD_CHROME unset; controller must run the local Chrome drill")
	}
	info, err := os.Stat(chrome)
	if err != nil {
		t.Fatal(err)
	}
	if !info.Mode().IsRegular() || info.Mode().Perm()&0111 == 0 {
		t.Fatal("WORLD_CHROME must name an executable file")
	}
	for _, jsOff := range []bool{false, true} {
		name := "js-on"
		if jsOff {
			name = "js-off"
		}
		t.Run(name, func(t *testing.T) {
			d := newHandlerDaemon(t)
			c := genesisCommit("chrome-drill")
			if err := d.store.Commit(boundedTestContext(t), c); err != nil {
				t.Fatal(err)
			}
			k := c.Entry.Header.EntryIndex
			d.aguiBudget, d.aguiTick = time.Second, 10*time.Millisecond
			events := make(chan workbenchChromeTrace, 128)
			handler := d.Handler()
			var mu sync.Mutex
			posts := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "POST" && r.URL.Path == "/agui/" {
					var copy bytes.Buffer
					b, err := io.ReadAll(io.TeeReader(r.Body, &copy))
					r.Body.Close()
					r.Body = io.NopCloser(bytes.NewReader(b))
					var input struct {
						State struct {
							LastIndex *int64 `json:"lastIndex"`
							Schema    string `json:"schema"`
						} `json:"state"`
					}
					if err == nil {
						err = json.Unmarshal(copy.Bytes(), &input)
					}
					if err == nil && (input.State.LastIndex == nil || input.State.Schema != "world/agui-state/v1") {
						err = fmt.Errorf("POST missing recognized schema/lastIndex: %s", copy.Bytes())
					}
					cursor := int64(-2)
					if input.State.LastIndex != nil {
						cursor = *input.State.LastIndex
					}
					mu.Lock()
					posts++
					mu.Unlock()
					events <- workbenchChromeTrace{method: "POST", cursor: cursor, err: err}
				}
				handler.ServeHTTP(w, r)
				if r.Method == "GET" && r.URL.Path == "/workbench" {
					events <- workbenchChromeTrace{method: "GET"}
				}
			}))
			t.Cleanup(server.Close)
			startWorkbenchChrome(t, chrome, server.URL+"/workbench", jsOff)
			waitWorkbenchChromeTrace(t, events, "GET", 0)
			if jsOff {
				window := time.NewTimer(time.Second)
				defer window.Stop()
				select {
				case e := <-events:
					t.Fatalf("JS-off unexpected request: %+v", e)
				case <-window.C:
				}
				mu.Lock()
				n := posts
				mu.Unlock()
				if n != 0 {
					t.Fatalf("JS-off POST count=%d want 0", n)
				}
				t.Log("JS-off: page GET observed; zero POSTs for 1 s (instrument control)")
				return
			}
			waitWorkbenchChromeTrace(t, events, "POST", k)
			next := testCommit(c.NextWorld, k+1, "chrome-drill")
			if err := d.store.Commit(boundedTestContext(t), next); err != nil {
				t.Fatal(err)
			}
			t.Logf("committed entry %d after POST cursor %d", k+1, k)
			waitWorkbenchChromeTrace(t, events, "GET", 0)
			waitWorkbenchChromeTrace(t, events, "POST", k+1)
		})
	}
}
