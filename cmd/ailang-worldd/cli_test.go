package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

type cliCommit struct {
	ObservedHead string      `json:"observedHead"`
	Objects      []cliObject `json:"objects"`
	NextWorld    cliWorld    `json:"nextWorld"`
	Entry        cliEntry    `json:"entry"`
}

type cliObject struct {
	Hash, InterfaceHash, SemanticID, Provenance string
	Payload                                     []byte
}

type cliWorld struct {
	Ref, StateRoot, LogHead string
	Revision                int64
}

type cliEntry struct {
	Header        cliHeader `json:"header"`
	EntryHash     string    `json:"entryHash"`
	TransitionRef string    `json:"transitionRef"`
}

type cliHeader struct {
	EntryIndex, SemanticsEpoch               int64
	TransitionFn, Interpreter, PrevEntryHash string
	WrittenBy                                string
}

func makeCLICommit(observed cliWorld, index int64, label string) cliCommit {
	payload := []byte("payload-" + label)
	objectHash := hashref.SumSHA256(payload).String()
	entryHash := hashref.SumSHA256([]byte("entry-" + label)).String()
	return cliCommit{
		ObservedHead: observed.Ref,
		Objects: []cliObject{{
			Hash: objectHash, InterfaceHash: hashref.SumSHA256([]byte("interface-" + label)).String(),
			SemanticID: "test/" + label, Provenance: "cli-e2e", Payload: payload,
		}},
		NextWorld: cliWorld{
			Ref: hashref.SumSHA256([]byte("world-" + label)).String(), Revision: index,
			StateRoot: hashref.SumSHA256([]byte("state-" + label)).String(), LogHead: entryHash,
		},
		Entry: cliEntry{
			Header: cliHeader{
				EntryIndex: index, SemanticsEpoch: 1,
				TransitionFn: hashref.SumSHA256([]byte("fn-" + label)).String(),
				Interpreter:  hashref.SumSHA256([]byte("interpreter")).String(),
				PrevEntryHash: func() string {
					if observed.LogHead != "" {
						return observed.LogHead
					}
					return hashref.SumSHA256([]byte("genesis-prev-" + label)).String()
				}(),
				WrittenBy: "cli-e2e",
			},
			EntryHash: entryHash, TransitionRef: objectHash,
		},
	}
}

func writeCommitFile(t *testing.T, dir, name string, commit cliCommit) string {
	t.Helper()
	data, err := json.Marshal(commit)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, name+".json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func runCLI(t *testing.T, addr string, args ...string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	full := append([]string{"--addr", addr}, args...)
	code := run(full, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func requireCLIOK(t *testing.T, addr string, args ...string) string {
	t.Helper()
	code, stdout, stderr := runCLI(t, addr, args...)
	if code != exitOK {
		t.Fatalf("%v exit=%d stderr=%s stdout=%s", args, code, stderr, stdout)
	}
	if strings.TrimSpace(stdout) == "" {
		t.Fatalf("%v returned empty content", args)
	}
	return stdout
}

// TestCLIRealSubprocessEpisode is the M2.C non-vacuity gate. The server is the
// built binary in a separate process; every request is issued through run and
// cli.go, never through a test-side HTTP call.
func TestCLIRealSubprocessEpisode(t *testing.T) {
	temp := t.TempDir()
	binary := filepath.Join(temp, "ailang-worldd")
	buildCtx, cancelBuild := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelBuild()
	build := exec.CommandContext(buildCtx, "go", "build", "-o", binary, "./cmd/ailang-worldd")
	build.Dir = filepath.Join("..", "..")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build subprocess binary: %v\n%s", err, output)
	}

	daemonCtx, cancelDaemon := context.WithTimeout(context.Background(), 30*time.Second)
	cmd := exec.CommandContext(daemonCtx, binary, "serve", "--db", filepath.Join(temp, "world.db"), "--bind", "127.0.0.1:0")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	var daemonErr bytes.Buffer
	cmd.Stderr = &daemonErr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start daemon: %v", err)
	}
	waited := make(chan error, 1)
	go func() { waited <- cmd.Wait() }()
	t.Cleanup(func() {
		cancelDaemon()
		select {
		case <-waited:
		case <-time.After(2 * time.Second):
			_ = cmd.Process.Kill()
			select {
			case <-waited:
			case <-time.After(2 * time.Second):
				t.Errorf("daemon process did not exit after kill")
			}
		}
	})

	announced := make(chan string, 1)
	readErr := make(chan error, 1)
	go func() {
		scanner := bufio.NewScanner(stdout)
		if scanner.Scan() {
			announced <- scanner.Text()
			return
		}
		readErr <- scanner.Err()
	}()
	var line string
	select {
	case line = <-announced:
	case err := <-readErr:
		t.Fatalf("read daemon announcement: %v; stderr=%s", err, daemonErr.String())
	case err := <-waited:
		t.Fatalf("daemon exited before announcement: %v; stderr=%s", err, daemonErr.String())
	case <-time.After(5 * time.Second):
		t.Fatalf("daemon announcement timed out; stderr=%s", daemonErr.String())
	}
	if !strings.HasPrefix(line, daemon.ListenAnnouncePrefix) {
		t.Fatalf("announcement=%q, want prefix %q", line, daemon.ListenAnnouncePrefix)
	}
	addr := strings.TrimPrefix(line, daemon.ListenAnnouncePrefix)

	health := requireCLIOK(t, addr, "health")
	if !strings.Contains(health, `"status":"ok"`) || !strings.Contains(health, filepath.Join(temp, "world.db")) {
		t.Fatalf("health content=%s", health)
	}

	genesis := makeCLICommit(cliWorld{}, 0, "genesis")
	genesisFile := writeCommitFile(t, temp, "genesis", genesis)
	commitBody := requireCLIOK(t, addr, "commit", "--file", genesisFile)
	if !strings.Contains(commitBody, genesis.NextWorld.Ref) {
		t.Fatalf("genesis commit content=%s", commitBody)
	}
	if head := requireCLIOK(t, addr, "head"); strings.TrimSpace(head) != genesis.NextWorld.Ref {
		t.Fatalf("head=%q want=%q", head, genesis.NextWorld.Ref)
	}
	if world := requireCLIOK(t, addr, "world", "get", genesis.NextWorld.Ref); !strings.Contains(world, genesis.NextWorld.StateRoot) {
		t.Fatalf("world content=%s", world)
	}
	object := requireCLIOK(t, addr, "object", "get", genesis.Objects[0].Hash, "--payload")
	if !strings.Contains(object, "cGF5bG9hZC1nZW5lc2lz") {
		t.Fatalf("object payload content=%s", object)
	}
	if log := requireCLIOK(t, addr, "log", "get", "0"); !strings.Contains(log, genesis.Entry.EntryHash) {
		t.Fatalf("log content=%s", log)
	}
	if page := requireCLIOK(t, addr, "log", "range", "--from", "0", "--limit", "1"); !strings.Contains(page, `"entryIndex":0`) {
		t.Fatalf("log range content=%s", page)
	}
	registry := requireCLIOK(t, addr, "registry", "get", "world/epoch-registry/v1")
	if !strings.Contains(registry, `"name":"world/epoch-registry/v1"`) || !strings.Contains(registry, `"head":"sha256:`) {
		t.Fatalf("registry content=%s", registry)
	}

	winner := makeCLICommit(genesis.NextWorld, 1, "winner")
	requireCLIOK(t, addr, "commit", "--file", writeCommitFile(t, temp, "winner", winner))
	stale := makeCLICommit(genesis.NextWorld, 1, "stale")
	code, _, conflict := runCLI(t, addr, "commit", "--file", writeCommitFile(t, temp, "stale", stale))
	if code != exitUsage || !strings.Contains(conflict, "HTTP 409 HeadConflict") ||
		!strings.Contains(conflict, "observedHead="+genesis.NextWorld.Ref) ||
		!strings.Contains(conflict, "selectedHead="+winner.NextWorld.Ref) {
		t.Fatalf("conflict exit=%d stderr=%s", code, conflict)
	}
	selected := regexp.MustCompile(`selectedHead=([^)]*)`).FindStringSubmatch(conflict)
	if len(selected) != 2 {
		t.Fatalf("conflict does not expose selected head: %s", conflict)
	}
	replannedBase := winner.NextWorld
	replannedBase.Ref = selected[1]
	replanned := makeCLICommit(replannedBase, 2, "replanned")
	replannedBody := requireCLIOK(t, addr, "commit", "--file", writeCommitFile(t, temp, "replanned", replanned))
	if !strings.Contains(replannedBody, replanned.NextWorld.Ref) {
		t.Fatalf("replanned commit content=%s", replannedBody)
	}
}

func TestClientDeadlineAgainstAcceptingServer(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ln.Close() }()
	accepted := make(chan net.Conn, 1)
	go func() {
		conn, err := ln.Accept()
		if err == nil {
			accepted <- conn
		}
	}()

	c := newClient("http://" + ln.Addr().String())
	c.timeout = 150 * time.Millisecond
	start := time.Now()
	_, _, err = c.get("/v1/health")
	elapsed := time.Since(start)
	if err == nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error=%v, want context deadline exceeded", err)
	}
	if elapsed > time.Second {
		t.Fatalf("deadline call took %s, want <=1s", elapsed)
	}
	select {
	case conn := <-accepted:
		_ = conn.Close()
	case <-time.After(time.Second):
		t.Fatal("server never accepted client connection")
	}
}

func TestCommitFileBound(t *testing.T) {
	path := filepath.Join(t.TempDir(), "large.json")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Truncate(maxClientCommitBytes + 1); err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	var stdout, stderr bytes.Buffer
	if code := runCommit("http://127.0.0.1:1", []string{"--file", path}, &stdout, &stderr); code != exitUsage {
		t.Fatalf("exit=%d want=%d", code, exitUsage)
	}
	if !strings.Contains(stderr.String(), fmt.Sprintf("exceeds %d bytes", maxClientCommitBytes)) {
		t.Fatalf("stderr=%q", stderr.String())
	}
}
