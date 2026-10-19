//go:build darwin || linux

package main

// Row 189 M3: episode start over a real loopback socket and a real child
// process with the production /dev/tty (realSessionEnv). These need loopback
// networking and setsid, which the agent sandbox does not give; the operator
// runs them unsandboxed.

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/ptytest"
)

func healthServer(t *testing.T, h daemon.HealthResponse) string {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/health" {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(h)
	}))
	t.Cleanup(srv.Close)
	return strings.TrimPrefix(srv.URL, "http://")
}

// The real fetch: a health body over loopback decodes; a closed port is
// "no daemon", which episode start answers with the serve line.
func TestGetDaemonHealthOverLoopback(t *testing.T) {
	want := daemon.HealthResponse{Status: "ok", Build: testBuild, ToolRef: "sha256:" + strings.Repeat("1", 64), CredentialsDB: "/c/credentials.db",
		ExecProfiles: []daemon.HealthExecProfile{{Project: "proj", Digest: "sha256:" + strings.Repeat("2", 64), File: "/p/proj.json"}}}
	got, err := getDaemonHealth(healthServer(t, want))
	if err != nil || got.Build != want.Build || got.ToolRef != want.ToolRef || len(got.ExecProfiles) != 1 || got.ExecProfiles[0] != want.ExecProfiles[0] {
		t.Fatalf("getDaemonHealth = %+v, %v", got, err)
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("loopback listen: %v", err)
	}
	closed := ln.Addr().String()
	_ = ln.Close()
	var down *daemonUnreachableError
	if _, err := getDaemonHealth(closed); !errors.As(err, &down) {
		t.Fatalf("a closed port: %v, want daemonUnreachableError", err)
	}
}

const episodePtyChildArgs = "AILANG_WORLDD_EPISODE_PTY_CHILD_ARGS"

// TestEpisodePtyChildProcess is the child the end-to-end test starts: the
// real run(), with this test binary given the stub daemon's build identity (a
// go test binary carries no vcs.revision of its own).
func TestEpisodePtyChildProcess(t *testing.T) {
	raw := os.Getenv(episodePtyChildArgs)
	if raw == "" {
		t.Skip("the child process of the episode pty test; it runs only when that test starts it")
	}
	var args []string
	if err := json.Unmarshal([]byte(raw), &args); err != nil {
		t.Fatal(err)
	}
	episodeSelfBuild = func() daemon.BuildIdentity { return testBuild }
	os.Exit(run(args, os.Stdout, os.Stderr))
}

// AC2 end to end: a matching daemon on loopback, no controlling terminal
// (setsid, the agent-harness shape), "y" piped to stdin — refused, with no
// worktree, no credential and no session file.
func TestEpisodeStartWithoutAControllingTerminalRefusesEndToEnd(t *testing.T) {
	r := newEpisodeRun(t)
	p := writeTestProfile(t, r.profiles, "proj", "")
	addr := healthServer(t, matchingHealth(t, r.f, r.profiles, p))
	encoded, err := json.Marshal(append([]string{"episode", "start"}, r.args("ep1", "--profile", "proj", "--addr", addr)...))
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestEpisodePtyChildProcess$", "-test.count=1")
	cmd.Env = append(os.Environ(), episodePtyChildArgs+"="+string(encoded))
	var stdout, stderr ptytest.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	cmd.Stdin = strings.NewReader("y\n")
	ptytest.WithoutControllingTerminal(cmd)
	code := childExitCode(t, cmd.Run())
	if code != exitUsage || !strings.Contains(stderr.String(), "refusing: no controlling terminal") {
		t.Fatalf("exit %d, want %d refusing for no controlling terminal\nstdout %q\nstderr %q", code, exitUsage, stdout.String(), stderr.String())
	}
	r.f.assertNothingProvisioned(t, "ep1")
	r.assertNoSessionFiles(t)
	if _, err := os.Stat(filepath.Join(r.f.ws, "ep1")); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("a worktree exists")
	}
}
