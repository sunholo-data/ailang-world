package store

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// This file proves Decision 6 across real process death. In particular, the
// mid-commit stop is reached from commitBeforeOutcomeHook after the world, log,
// and head writes have executed but before the outcome is inserted. The hook
// prints READY and blocks; no sleep determines where SIGKILL lands.
const (
	crashHelperDBEnv     = "WORLDD_CRASH_HELPER_DB"
	crashHelperStopEnv   = "WORLDD_CRASH_HELPER_STOP"
	crashHelperEffectEnv = "WORLDD_CRASH_HELPER_EFFECT"
	crashHelperReady     = "READY"
	crashInvocationID    = "crash-proof-invocation"
	crashDeadline        = 30 * time.Second
)

func crashFixture(t *testing.T, s *Store) Commit {
	t.Helper()
	genesis := seedGenesis(t, s)
	body := obj("crash-proof-transition-body", "transition/body")
	entryHash := hashref.SumSHA256([]byte("crash-proof-entry"))
	return Commit{
		InvocationID: crashInvocationID,
		ObservedHead: genesis.Ref,
		Objects:      []Object{body},
		NextWorld: World{
			Ref:       hashref.SumSHA256([]byte("crash-proof-world")),
			Revision:  1,
			StateRoot: hashref.SumSHA256([]byte("crash-proof-state")),
			LogHead:   entryHash,
		},
		Entry: LogEntry{
			Header: LogHeader{
				EntryIndex:     1,
				SemanticsEpoch: 1,
				TransitionFn:   hashref.SumSHA256([]byte("crash-proof-fn")),
				Interpreter:    hashref.SumSHA256([]byte("crash-proof-interpreter")),
				PrevEntryHash:  genesis.LogHead,
				WrittenBy:      "crash-proof",
			},
			EntryHash:     entryHash,
			TransitionRef: body.Hash,
		},
	}
}

// TestCrashHelperProcess is the re-exec body. READY is emitted only after the
// requested stop point has genuinely been reached. Closing stdin lets the
// after-outcome negative control exit cleanly; crash cases use real SIGKILL.
func TestCrashHelperProcess(t *testing.T) {
	dbPath := os.Getenv(crashHelperDBEnv)
	if dbPath == "" {
		t.Skip("subprocess helper; runs only when re-exec'd with " + crashHelperDBEnv)
	}
	stop := os.Getenv(crashHelperStopEnv)
	s, err := Open(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "crash helper Open: %v\n", err)
		os.Exit(3)
	}
	c := crashFixture(t, s)
	if _, _, err := s.AppendIntent(crashInvocationID, testCommitIntent(crashInvocationID, c)); err != nil {
		fmt.Fprintf(os.Stderr, "crash helper AppendIntent: %v\n", err)
		os.Exit(4)
	}
	readyAndBlock := func() {
		fmt.Println(crashHelperReady)
		_, _ = io.Copy(io.Discard, os.Stdin)
	}
	switch stop {
	case "after-intent":
		readyAndBlock()
	case "after-external-effect":
		if err := os.WriteFile(os.Getenv(crashHelperEffectEnv), []byte("dispatch\n"), 0o600); err != nil {
			fmt.Fprintf(os.Stderr, "crash helper probe effect: %v\n", err)
			os.Exit(5)
		}
		readyAndBlock()
	case "mid-commit-before-outcome":
		commitBeforeOutcomeHook = readyAndBlock
		if err := s.Commit(c); err != nil {
			fmt.Fprintf(os.Stderr, "crash helper Commit: %v\n", err)
			os.Exit(6)
		}
	case "after-outcome":
		if err := s.Commit(c); err != nil {
			fmt.Fprintf(os.Stderr, "crash helper Commit: %v\n", err)
			os.Exit(7)
		}
		readyAndBlock()
	default:
		fmt.Fprintf(os.Stderr, "unknown crash stop %q\n", stop)
		os.Exit(8)
	}
	_ = s.Close()
	os.Exit(0)
}

type crashProcess struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stderr *syncBuffer
	waited bool
}

func startCrashProcess(t *testing.T, dbPath, effectPath, stop string) *crashProcess {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestCrashHelperProcess$")
	cmd.Env = append(os.Environ(),
		crashHelperDBEnv+"="+dbPath,
		crashHelperStopEnv+"="+stop,
		crashHelperEffectEnv+"="+effectPath,
	)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	stderr := &syncBuffer{}
	cmd.Stderr = stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start crash helper: %v", err)
	}
	p := &crashProcess{cmd: cmd, stdin: stdin, stderr: stderr}
	t.Cleanup(func() {
		if !p.waited {
			p.killAndWait(t)
		}
	})

	lineCh := make(chan string, 1)
	errCh := make(chan error, 1)
	go func() {
		line, err := bufio.NewReader(stdout).ReadString('\n')
		if err != nil {
			errCh <- err
			return
		}
		lineCh <- strings.TrimSpace(line)
	}()
	deadline := time.Now().Add(crashDeadline)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case line := <-lineCh:
			if line != crashHelperReady {
				t.Fatalf("crash helper first line=%q, want READY; stderr:\n%s", line, stderr.String())
			}
			return p
		case err := <-errCh:
			t.Fatalf("crash helper never signalled READY: %v; stderr:\n%s", err, stderr.String())
		case <-ticker.C:
			// Poll the captured process, never a process-name pattern.
			if err := cmd.Process.Signal(syscall.Signal(0)); err != nil {
				t.Fatalf("captured crash helper pid %d exited before READY: %v; stderr:\n%s",
					cmd.Process.Pid, err, stderr.String())
			}
			if time.Now().After(deadline) {
				t.Fatalf("captured crash helper pid %d missed READY deadline; stderr:\n%s",
					cmd.Process.Pid, stderr.String())
			}
		}
	}
}

func (p *crashProcess) wait(t *testing.T) error {
	t.Helper()
	if p.waited {
		return nil
	}
	done := make(chan error, 1)
	go func() { done <- p.cmd.Wait() }()
	deadline := time.Now().Add(crashDeadline)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case err := <-done:
			p.waited = true
			return err
		case <-ticker.C:
			// Deadline discipline is tied to the captured os.Process.
			_ = p.cmd.Process.Signal(syscall.Signal(0))
			if time.Now().After(deadline) {
				t.Fatalf("captured crash helper pid %d missed exit deadline", p.cmd.Process.Pid)
			}
		}
	}
}

func (p *crashProcess) killAndWait(t *testing.T) {
	t.Helper()
	if p.waited {
		return
	}
	if err := p.cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
		t.Fatalf("SIGKILL captured crash helper pid %d: %v", p.cmd.Process.Pid, err)
	}
	_ = p.stdin.Close()
	if err := p.wait(t); err == nil {
		t.Fatalf("SIGKILLed helper pid %d exited cleanly; real process death was not observed",
			p.cmd.Process.Pid)
	}
}

func assertCrashStore(t *testing.T, dbPath, effectPath, stop string) {
	t.Helper()
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("reopen after %s: %v", stop, err)
	}
	defer func() { _ = s.Close() }()
	c := crashFixtureForRead()
	receipt, ok, err := s.GetReceipt(crashInvocationID)
	if err != nil {
		t.Fatal(err)
	}
	wantState := ReceiptIndeterminate
	wantCommitted := false
	if stop == "after-outcome" {
		wantState, wantCommitted = ReceiptResolved, true
	}
	if !ok || receipt.State != wantState {
		t.Fatalf("%s receipt=(ok=%v,state=%s), want (true,%s)", stop, ok, receipt.State, wantState)
	}
	pending, err := s.PendingIntents(MaxPendingIntentsPage)
	if err != nil {
		t.Fatal(err)
	}
	wantPending := 1
	if wantCommitted {
		wantPending = 0
	}
	if len(pending) != wantPending {
		t.Fatalf("%s PendingIntents=%d, want %d", stop, len(pending), wantPending)
	}
	_, worldOK, err := s.GetWorld(c.NextWorld.Ref)
	if err != nil {
		t.Fatal(err)
	}
	_, entryOK, err := s.GetLogEntry(c.Entry.Header.EntryIndex)
	if err != nil {
		t.Fatal(err)
	}
	if worldOK != wantCommitted || entryOK != wantCommitted {
		t.Fatalf("%s world=%v entry=%v, want both %v", stop, worldOK, entryOK, wantCommitted)
	}
	effect, err := os.ReadFile(effectPath)
	hasEffect := err == nil
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
	wantEffect := stop == "after-external-effect"
	if hasEffect != wantEffect {
		t.Fatalf("%s probe effect=%v, want %v", stop, hasEffect, wantEffect)
	}
	if hasEffect && string(effect) != "dispatch\n" {
		t.Fatalf("%s probe effect content=%q: recovery re-dispatched it", stop, effect)
	}
}

// This read-only projection avoids trying to seed genesis after reopening.
func crashFixtureForRead() Commit {
	entryHash := hashref.SumSHA256([]byte("crash-proof-entry"))
	return Commit{
		NextWorld: World{Ref: hashref.SumSHA256([]byte("crash-proof-world"))},
		Entry:     LogEntry{Header: LogHeader{EntryIndex: 1}, EntryHash: entryHash},
	}
}

func TestCrashReceiptLawAtNamedStopPoints(t *testing.T) {
	for _, stop := range []string{
		"after-intent", "after-external-effect", "mid-commit-before-outcome", "after-outcome",
	} {
		t.Run(stop, func(t *testing.T) {
			dir := t.TempDir()
			dbPath := filepath.Join(dir, "world.db")
			effectPath := filepath.Join(dir, "probe.dispatches")
			p := startCrashProcess(t, dbPath, effectPath, stop)
			p.killAndWait(t)
			assertCrashStore(t, dbPath, effectPath, stop)
		})
	}
}

func TestCrashNegativeControlCleanRunProducesReceipt(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "world.db")
	effectPath := filepath.Join(dir, "probe.dispatches")
	p := startCrashProcess(t, dbPath, effectPath, "after-outcome")
	_ = p.stdin.Close()
	if err := p.wait(t); err != nil {
		t.Fatalf("clean helper: %v; stderr:\n%s", err, p.stderr.String())
	}
	assertCrashStore(t, dbPath, effectPath, "after-outcome")
}
