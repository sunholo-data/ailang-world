package main

// Row 201 A2 sprint 4 M7 (design_docs/planned/w-world-approvals-by-key-a2.md §2.14, §6 XB1, XB2):
// the harness of the passkey page's browser tests. testdata/passkey_cdp.mjs drives Chrome headless
// with a temporary profile and a CDP virtual authenticator, a test key (D-WORLD-109): never a real
// profile, keychain, passkey or platform authenticator. The page is served by the production daemon
// handler. Gated: WORLD_CHROME and node on PATH, else a counted skip, which
// WORLD_REQUIRE_PASSKEY_BROWSER=1 turns into a failure.

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/daemon"
)

// passkeyBrowserBudget bounds one driver (Chrome's start, every ceremony and the reap). Measured
// (plan §1.2 C19): one XB1 or XB2 test takes 4-19 s on the planner's machine at load 20-80.
const passkeyBrowserBudget = 120 * time.Second

// pageView is what the page shows (the driver reads it after each step).
type pageView struct {
	State, What, Status, Button, Code, Result, SaveHref, SaveName string
	StatusError, ButtonHidden, ButtonDisabled, SaveHidden         bool
	Violations                                                    []string
	Op                                                            *struct {
		Method                                                           string
		Signal, Aborted                                                  bool
		AbortedAfter, AbortedAfterDate, ResolvedAfter, ResolvedAfterDate *float64
	}
}

// pageCeremony is one click on the ceremony button: the page while pending, when settled, and
// after the Save click; the download events; the page's errors.
type pageCeremony struct {
	Pending, Settled, After pageView
	SaveClicked             bool
	Downloads               []struct{ Event, State, Name, URL string }
	Errors                  []string
	RealMs                  int64
}

type passkeyBrowser struct {
	in     io.WriteCloser
	out    *bufio.Reader
	chrome string
}

// startPasskeyBrowser starts the driver, or skips (counted) when Chrome or node is missing.
func startPasskeyBrowser(t *testing.T) *passkeyBrowser {
	t.Helper()
	chrome := os.Getenv("WORLD_CHROME")
	node, err := exec.LookPath("node")
	if chrome == "" || err != nil {
		msg := fmt.Sprintf("SKIP COUNTED: the passkey page browser test needs WORLD_CHROME (%q) and node on PATH (%v)", chrome, err)
		if os.Getenv("WORLD_REQUIRE_PASSKEY_BROWSER") == "1" {
			t.Fatal(msg)
		}
		t.Skip(msg)
	}
	ctx, cancel := context.WithTimeout(context.Background(), passkeyBrowserBudget)
	cmd := exec.CommandContext(ctx, node, filepath.Join("testdata", "passkey_cdp.mjs"), t.TempDir())
	var stderr syncBuffer // cli_test.go's: os/exec's copier writes until Wait returns (TestSubprocessSinkGate)
	cmd.Stderr = &stderr
	in, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		cancel()
		t.Fatal(err)
	}
	b := &passkeyBrowser{in: in, out: bufio.NewReader(out)}
	var ready struct {
		Ready  bool
		Chrome string
		Pid    int
	}
	t.Cleanup(func() {
		_ = in.Close()
		werr := cmd.Wait()
		cancel()
		if werr != nil && ready.Pid > 0 { // the driver did not reap Chrome itself
			_ = syscall.Kill(ready.Pid, syscall.SIGKILL)
		}
		if t.Failed() {
			t.Logf("driver: %v\n%s", werr, stderr.String())
		}
	})
	if err := b.read(&ready); err != nil || !ready.Ready {
		t.Fatalf("the driver did not start: %v", err)
	}
	b.chrome = ready.Chrome
	return b
}

func (b *passkeyBrowser) read(v any) error {
	line, err := b.out.ReadBytes('\n')
	if err != nil {
		return err
	}
	return json.Unmarshal(line, v)
}

// call sends one driver command and decodes its answer into out. It never calls t.Fatal, so a
// terminal's read goroutine may use it.
func (b *passkeyBrowser) call(out any, cmd string, args map[string]any) error {
	if args == nil {
		args = map[string]any{}
	}
	args["cmd"] = cmd
	line, _ := json.Marshal(args)
	if _, err := b.in.Write(append(line, '\n')); err != nil {
		return err
	}
	var reply struct {
		Ok    json.RawMessage
		Error string
	}
	if err := b.read(&reply); err != nil {
		return err
	}
	if reply.Error != "" {
		return fmt.Errorf("driver %s: %s", cmd, reply.Error)
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(reply.Ok, out)
}

func (b *passkeyBrowser) do(t *testing.T, out any, cmd string, args map[string]any) {
	t.Helper()
	if err := b.call(out, cmd, args); err != nil {
		t.Fatal(err)
	}
}

// passkeyPageOrigin serves the production daemon handler on localhost and returns its origin.
func passkeyPageOrigin(t *testing.T) string {
	t.Helper()
	d, err := daemon.New(boundedTestContext(t), daemon.Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: daemon.DefaultBindHost, ErrorLog: io.Discard})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = d.Close() })
	srv := httptest.NewServer(d.Handler())
	t.Cleanup(srv.Close)
	_, port, _ := net.SplitHostPort(srv.Listener.Addr().String())
	return "http://localhost:" + port
}

var printedLinkRE = regexp.MustCompile(`(?m)^  (http\S+/surface/passkey#\S+)$`)

// printedLink is the one ceremony link a CLI run printed.
func printedLink(out string) (string, error) {
	m := printedLinkRE.FindAllStringSubmatch(out, -1)
	if len(m) != 1 {
		return "", fmt.Errorf("%d ceremony links in %q, want 1", len(m), out)
	}
	return m[0][1], nil
}

// saved checks that c saved exactly the page's JSON plus a newline as dir/name through a blob:
// download of origin, with no CSP violation or page error, and returns the file.
func saved(t *testing.T, c pageCeremony, origin, dir, name string) []byte {
	t.Helper()
	if c.After.State != "done" || c.After.StatusError || !c.SaveClicked || c.After.SaveName != name {
		t.Fatalf("the page did not finish: state %q, status %q, save %q clicked %v", c.After.State, c.After.Status, c.After.SaveName, c.SaveClicked)
	}
	var began, completed bool
	for _, d := range c.Downloads {
		began = began || (d.Event == "Browser.downloadWillBegin" && d.Name == name && strings.HasPrefix(d.URL, "blob:"+origin+"/"))
		completed = completed || (d.Event == "Browser.downloadProgress" && d.State == "completed")
	}
	if !began || !completed {
		t.Fatalf("no completed blob: download of %s: %+v", name, c.Downloads)
	}
	if len(c.After.Violations) != 0 || len(c.Errors) != 0 {
		t.Fatalf("CSP violations %q, page errors %q", c.After.Violations, c.Errors)
	}
	file, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		t.Fatal(err)
	}
	if string(file) != c.After.Result+"\n" {
		t.Fatalf("the saved file is not the page's JSON plus a newline:\n%s\n%s", file, c.After.Result)
	}
	return file
}

// opString shows the injected operation's record.
func opString(o any) string { b, _ := json.Marshal(o); return string(b) }

// worldFiles is the world-*.json files in dir: what the page's Save link can write. Anything else
// Chrome leaves there is logged, not counted (plan §1.2 C21: a stray "downloads.html" download,
// outside every ceremony's events, seen under load).
func worldFiles(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "world-") && strings.HasSuffix(e.Name(), ".json") {
			names = append(names, e.Name())
		} else {
			t.Logf("not a page save, not counted: %s", e.Name())
		}
	}
	return names
}
