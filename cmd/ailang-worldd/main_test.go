package main

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
)

// TestMain removes the registry credential from this package's test process.
//
// It is not a convenience. run() now REFUSES to start when the credential is
// ambient (Decision 4), and the credential IS ambient in the mission
// operator's shell — so without this every test below would pass in CI and
// fail locally, which is the worst possible property for a gate. Unsetting it
// here makes the suite host-independent and leaves the ambient behaviour to
// the one test that sets it deliberately,
// TestRunRefusesAnAmbientRegistryCredential.
func TestMain(m *testing.M) {
	if err := os.Unsetenv(broker.RegistryCredentialVariable); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

// TestRunRefusesAnAmbientRegistryCredential is MUT-SM-AMBIENT-KEY's gate and
// its mandatory negative control: the SAME invocation, differing only in the
// environment. Identical outcomes in both arms would mean the check never read
// the environment.
func TestRunRefusesAnAmbientRegistryCredential(t *testing.T) {
	const sentinel = "worldd-ambient-sentinel-not-a-real-key"

	t.Setenv(broker.RegistryCredentialVariable, sentinel)
	var stdout, stderr bytes.Buffer
	if got := run([]string{"help"}, &stdout, &stderr); got != exitFatal {
		t.Fatalf("ambient arm exit = %d, want %d (fatal)", got, exitFatal)
	}
	if !strings.Contains(stderr.String(), broker.RegistryCredentialVariable) {
		t.Errorf("refusal %q does not name the variable", stderr.String())
	}
	if strings.Contains(stderr.String(), sentinel) || strings.Contains(stdout.String(), sentinel) {
		t.Fatal("the startup refusal printed the credential value")
	}

	// NEGATIVE CONTROL: same args, same streams, variable absent.
	t.Setenv(broker.RegistryCredentialVariable, "")
	stdout.Reset()
	stderr.Reset()
	if got := run([]string{"help"}, &stdout, &stderr); got != exitOK {
		t.Fatalf("NEGATIVE CONTROL failed: absent-variable arm exit = %d, want %d; "+
			"the same outcome in both arms means the check never read the environment",
			got, exitOK)
	}
}

// TestRunExitCodes pins Decision 5's exit-code contract and the flag-placement
// rule. Each case asserts the CODE, not merely "an error happened", because the
// exit code is the whole machine-readable surface of a CLI.
func TestRunExitCodes(t *testing.T) {
	db := filepath.Join(t.TempDir(), "world.db")

	cases := []struct {
		name string
		args []string
		want int
	}{
		// --addr is a GLOBAL CLIENT flag. Accepting it on serve would silently
		// ignore an operator's intent, so it is a usage error.
		{"addr is not a serve flag", []string{"--addr", "http://127.0.0.1:1", "serve", "--db", db}, exitUsage},
		{"serve rejects --addr as a serve flag", []string{"serve", "--addr", "http://127.0.0.1:1", "--db", db}, exitUsage},
		{"serve requires --db", []string{"serve"}, exitUsage},
		{"serve rejects a malformed --bind", []string{"serve", "--db", db, "--bind", "not-a-hostport"}, exitUsage},
		{"no arguments", nil, exitUsage},
		{"unknown verb", []string{"nope"}, exitUsage},
		{"help", []string{"help"}, exitOK},
		// A non-loopback bind is a FATAL startup refusal, not a usage error:
		// the flags were well-formed, the policy refused.
		{"non-loopback bind is fatal", []string{"serve", "--db", db, "--bind", "0.0.0.0:0"}, exitFatal},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if got := run(c.args, &stdout, &stderr); got != c.want {
				t.Fatalf("run(%q) = %d, want %d (stderr: %s)", c.args, got, c.want, stderr.String())
			}
		})
	}
}

// TestClientVerbsAgainstLiveDaemon drives the health and head verbs against a
// real in-process daemon on an ephemeral loopback port, asserting the printed
// CONTENT and the exit codes — a verb that printed nothing and returned 0 would
// fail. It also exercises the --addr global flag end to end.
func TestClientVerbsAgainstLiveDaemon(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "world.db")
	cfg := daemon.Config{DBPath: dbPath, BindHost: daemon.DefaultBindHost, BindPort: 0}

	pr, pw := io.Pipe()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ran := make(chan error, 1)
	go func() {
		err := daemon.Run(ctx, cfg, pw)
		_ = pw.Close()
		ran <- err
	}()

	line, err := bufio.NewReader(pr).ReadString('\n')
	if err != nil {
		t.Fatalf("read listen announcement: %v", err)
	}
	url := strings.TrimSpace(strings.TrimPrefix(line, daemon.ListenAnnouncePrefix))

	// health: exit 0 and the daemon's own facts on stdout.
	var stdout, stderr bytes.Buffer
	if code := run([]string{"--addr", url, "health"}, &stdout, &stderr); code != exitOK {
		t.Fatalf("health exit = %d, want %d (stderr: %s)", code, exitOK, stderr.String())
	}
	if !strings.Contains(stdout.String(), dbPath) {
		t.Fatalf("health stdout = %q, want it to carry the db path %q", stdout.String(), dbPath)
	}
	if !strings.Contains(stdout.String(), `"daemon_version"`) {
		t.Fatalf("health stdout = %q, want a daemon_version field", stdout.String())
	}

	// head: no head has been selected yet, so the daemon answers 404 and the
	// client reports a client error (exit 1) rather than a false success.
	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"--addr", url, "head"}, &stdout, &stderr); code != exitUsage {
		t.Fatalf("head exit = %d, want %d (stdout %q stderr %q)", code, exitUsage, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "404") {
		t.Fatalf("head stderr = %q, want it to report the 404 status", stderr.String())
	}

	cancel()
	select {
	case err := <-ran:
		if err != nil {
			t.Fatalf("daemon.Run returned %v, want a clean shutdown", err)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("daemon.Run did not return after cancellation")
	}
}

// TestClientCallIsBounded proves the D7 client deadline is real and injectable:
// the client dials a port nothing is listening on with a tiny timeout and must
// return an error quickly rather than hang. The injectable field is what lets
// M2.C's deadline test run in milliseconds instead of 30 s.
func TestClientCallIsBounded(t *testing.T) {
	c := newClient("http://127.0.0.1:1")
	if c.timeout != daemon.DefaultClientTimeout {
		t.Fatalf("newClient timeout = %s, want the D7 default %s", c.timeout, daemon.DefaultClientTimeout)
	}
	c.timeout = 100 * time.Millisecond

	start := time.Now()
	_, _, err := c.get("/v1/health")
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("client call to a dead port returned no error")
	}
	if elapsed > 10*time.Second {
		t.Fatalf("client call took %s — it is not bounded by the injected deadline", elapsed)
	}
}
