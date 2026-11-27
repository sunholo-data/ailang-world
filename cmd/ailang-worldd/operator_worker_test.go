package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
)

func TestOperatorLaunchWorker(t *testing.T) {
	f := newOperatorLaunchRig(t, "#!/bin/sh\n: >&3 || exit 3\n# The locks are open read/write. Inspect their presence without writing.\n: >&4 || exit 4\n"+strings.TrimPrefix(fakeLaunchMarker, "#!/bin/sh\n"))
	e := &f.env
	config, err := os.ReadFile(e.runnerPaths().Record)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("WORLD_OPERATOR_CONFIG_SHA256", serviceHash(config))
	e.operatorSupervisor = func(job operatorLaunchJob, log io.Writer) error {
		input := operatorWorkerInput{build: f.h.Build}
		run := func(other *os.File, handoff []byte, out io.Writer) int {
			input.handoff = bytes.NewReader(handoff)
			input.file = func(n uintptr, name string) *os.File {
				fd := job.locks[n-3]
				if n == 3 && other != nil {
					fd = other
				}
				dup, err := syscall.Dup(int(fd.Fd()))
				if err != nil {
					t.Fatal(err)
				}
				return os.NewFile(uintptr(dup), name)
			}
			return e.runOperatorLaunchWorker([]string{job.row.ID}, out, input)
		}
		before := activationEvidence(t, f.worktreeStartRig)
		foreign, err := os.CreateTemp(t.TempDir(), "unrelated")
		if err != nil {
			t.Fatal(err)
		}
		defer foreign.Close()
		occupied, err := os.OpenFile(f.creds+".operator-runner.lock", os.O_RDWR, 0)
		if err != nil {
			t.Fatal(err)
		}
		defer occupied.Close()
		for _, wrong := range []*os.File{foreign, occupied} {
			if rc := run(wrong, job.secret[:], io.Discard); rc != exitUsage ||
				!reflect.DeepEqual(before, activationEvidence(t, f.worktreeStartRig)) {
				t.Fatalf("bad inherited descriptor crossed reservation, rc=%d", rc)
			}
		}
		rc := run(nil, job.secret[:], log)

		if rc != 0 {
			t.Fatalf("worker exited %d", rc)
		}
		return nil
	}
	f.tick(t)
	marker, err := os.ReadFile(filepath.Join(f.bin, "ran"))
	s := f.receipt(t)
	if err != nil || string(marker) != "x" || s.Receipt.Status != "stopped" || s.Receipt.Birth == "" {
		t.Fatalf("worker child or fds 3/4: marker=%q receipt=%+v err=%v", marker, s.Receipt, err)
	}
	t.Run("catalogue", func(t *testing.T) { testOperatorCataloguePublication(t, f) })
}

func TestOperatorStartedRequiresBirth(t *testing.T) {
	var r store.OperatorLaunchReceipt
	if err := requireOperatorProcessBirth(&r, 42, func(int) (string, bool) { return "", false }); err == nil || r.Birth != "" {
		t.Fatalf("Started can advance without a kernel birth: %+v %v", r, err)
	}
	if err := requireOperatorProcessBirth(&r, 42, func(pid int) (string, bool) {
		if pid != 42 {
			t.Errorf("wrong PID: %d", pid)
		}
		return "kernel-start", true
	}); err != nil || r.PID != 42 || r.Birth != "kernel-start" {
		t.Fatalf("Started did not retain the birth: %+v %v", r, err)
	}
}

func TestOperatorSupervisorCap(t *testing.T) {
	dir := t.TempDir()
	hold, release, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer hold.Close()
	defer release.Close()
	events, report, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer events.Close()
	defer report.Close()
	if err := events.SetReadDeadline(time.Now().Add(activationTestBudget)); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "worker")
	script := `#!/bin/sh
trap 'printf "TERM\n" >&4; exit' TERM
printf 'ready\n' >&4
read -r unused <&3
printf 'released\n' >&4
`
	if err := os.WriteFile(path, []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	e := serviceEnv{}
	e.paths.Binary = path
	fired, stopped := make(chan time.Time, 1), false
	event := bufio.NewReader(events)
	wall := 2 * time.Minute
	e.operatorSupervisorCap = func(got time.Duration) (<-chan time.Time, func()) {
		if got != wall+time.Minute {
			t.Errorf("supervisor cap duration %s", got)
		}
		if line, err := event.ReadString('\n'); err != nil || line != "ready\n" {
			t.Fatalf("worker did not install its signal trap: %q (%v)", line, err)
		}
		fired <- time.Time{}
		return fired, func() { stopped = true }
	}
	job := operatorLaunchJob{wall: wall, locks: [3]*os.File{hold, report, report}}
	err = e.superviseOperatorLaunch(job, nil, io.Discard, store.OperatorHeartbeat{}, "")
	if err == nil || !strings.Contains(err.Error(), "did not report completion") || !stopped || errors.Is(err, os.ErrClosed) {
		t.Fatalf("supervisor did not honor its injected cap: %v", err)
	}
	// Unblock the worker even if the cap failed to signal it.
	if err := release.Close(); err != nil {
		t.Fatal(err)
	}
	if line, err := event.ReadString('\n'); err != nil || line != "TERM\n" {
		t.Fatalf("worker did not record SIGTERM: %q (%v)", line, err)
	}
}
