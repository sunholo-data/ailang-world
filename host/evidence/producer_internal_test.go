package evidence

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/procbound"
	"github.com/sunholo-data/ailang-world/host/proctest"
	"github.com/sunholo-data/ailang-world/host/store"
)

var internalTestKey = [32]byte{9, 8, 7}

// nopSeams is an inert ObjectReader and ObjectWriter for stage tests that never
// touch an object: it reports no busy window, so no ordering pin applies.
type nopSeams struct{}

func (nopSeams) ReadObject(context.Context, hashref.HashRef, int64) (store.ObjectMeta, []byte, error) {
	return store.ObjectMeta{}, nil, nil
}
func (nopSeams) PutObject(context.Context, store.Object) error { return nil }
func (nopSeams) BusyTimeout() time.Duration                    { return 0 }

// scriptedProducer writes an executable #!/bin/sh checker into its own temp dir,
// pins it by its own bytes, and constructs a producer through NewProducer. The
// script answers --version with a fixed line and otherwise runs body.
func scriptedProducer(t *testing.T, body string, execTimeout time.Duration, maxOut int64) (*Producer, string) {
	t.Helper()
	dir := t.TempDir()
	script := "#!/bin/sh\nif [ \"$1\" = \"--version\" ]; then echo \"AILANG scripted\"; exit 0; fi\n" + body + "\n"
	path := filepath.Join(dir, "ailang")
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	p, err := NewProducer(internalTestKey, nopSeams{}, nopSeams{}, CheckerConfig{
		ExecPath:        path,
		Compiler:        hashref.SumSHA256([]byte(script)),
		CompilerVersion: "AILANG scripted",
	}, ProducerConfig{
		ObjectReadTimeout:  time.Second,
		ObjectWriteTimeout: time.Second,
		ExecTimeout:        execTimeout,
		MaxOutputBytes:     maxOut,
	})
	if err != nil {
		t.Fatal(err)
	}
	return p, dir
}

func runEnv(dir string) []string {
	return []string{"PATH=" + os.Getenv("PATH"), "PIDDIR=" + dir}
}

// killRecorder observes every process-group kill, forwarding to the real kill.
type killRecorder struct {
	mu    sync.Mutex
	pgids []int
	times []time.Time
}

func (k *killRecorder) record(pgid int, at time.Time) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.pgids = append(k.pgids, pgid)
	k.times = append(k.times, at)
}

func (k *killRecorder) snapshot() ([]int, []time.Time) {
	k.mu.Lock()
	defer k.mu.Unlock()
	return append([]int(nil), k.pgids...), append([]time.Time(nil), k.times...)
}

func swapKillGroup(t *testing.T) *killRecorder {
	t.Helper()
	rec := &killRecorder{}
	orig := killGroup
	killGroup = func(pgid int) error {
		rec.record(pgid, time.Now())
		return orig(pgid)
	}
	t.Cleanup(func() { killGroup = orig })
	return rec
}

func TestProducerRefusesStderrOverflowIndependently(t *testing.T) {
	t.Run("overflow", func(t *testing.T) {
		rec := swapKillGroup(t)
		p, dir := scriptedProducer(t, `printf '{}'; head -c 64 /dev/zero | tr '\0' e >&2; exec sleep 30`, 10*time.Second, 16)
		start := time.Now()
		out, err := p.runBounded(context.Background(), dir, runEnv(dir), 16)
		elapsed := time.Since(start)
		if !errors.Is(err, ErrOutputLimit) {
			t.Fatalf("64 stderr bytes over a 16-byte cap: want ErrOutputLimit; got %v", err)
		}
		if errors.Is(err, ErrExecTimeout) {
			t.Fatalf("an overflow must not be reported as a timeout: %v", err)
		}
		if !strings.Contains(err.Error(), "stderr") {
			t.Fatalf("detail must name the stderr stream: %v", err)
		}
		if out != nil {
			t.Fatalf("overflow must return nil bytes; got %q", out)
		}
		if elapsed >= 5*time.Second {
			t.Fatalf("the overflow kill must end the child; it ran %s against a 10s timeout", elapsed)
		}
		if pgids, _ := rec.snapshot(); len(pgids) != 1 {
			t.Fatalf("overflow must kill the group exactly once; killGroup saw %v", pgids)
		}
	})

	t.Run("independence-control", func(t *testing.T) {
		p, dir := scriptedProducer(t, `printf '{}'; head -c 64 /dev/zero | tr '\0' e >&2`, 10*time.Second, 1024)
		out, err := p.runBounded(context.Background(), dir, runEnv(dir), 1024)
		if err != nil || string(out) != "{}" {
			t.Fatalf("64 stderr bytes under a 1024-byte cap must pass: out=%q err=%v", out, err)
		}
	})
}

func TestProducerExecTimeoutKillsTheWholeProcessGroup(t *testing.T) {
	rec := swapKillGroup(t)
	p, dir := scriptedProducer(t, `sleep 30 & echo $! > "$PIDDIR/grandchild.pid"; echo $$ > "$PIDDIR/child.pid"; wait`, 300*time.Millisecond, 1024)
	start := time.Now()
	done := make(chan error, 1)
	go func() {
		_, err := p.runBounded(context.Background(), dir, runEnv(dir), 1024)
		done <- err
	}()
	var err error
	select {
	case err = <-done:
	case <-time.After(20 * time.Second):
		t.Fatalf("runBounded did not return within 20s; want ErrExecTimeout after ExecTimeout=300ms")
	}
	grandchild := proctest.ReadPid(t, filepath.Join(dir, "grandchild.pid"))
	child := proctest.ReadPid(t, filepath.Join(dir, "child.pid"))
	proctest.ReapOnCleanup(t, grandchild)
	proctest.ReapOnCleanup(t, child)

	if !errors.Is(err, ErrExecTimeout) {
		t.Fatalf("want ErrExecTimeout; got %v", err)
	}
	if errors.Is(err, procbound.ErrCleanupIncomplete) {
		t.Fatalf("the group kill must complete the cleanup; got %v", err)
	}
	pgids, times := rec.snapshot()
	if len(pgids) != 1 || pgids[0] != child {
		t.Fatalf("killGroup must see exactly the child's group %d; saw %v", child, pgids)
	}
	if lag := times[0].Sub(start); lag < 300*time.Millisecond {
		t.Fatalf("first kill at %s after the call; want at least the 300ms ExecTimeout", lag)
	}
	if !proctest.DeadWithin(t, grandchild, 2*time.Second) {
		t.Fatalf("grandchild %d survived the process-group kill", grandchild)
	}
}

func TestProducerAdmissionRefusesBacklog(t *testing.T) {
	p, dir := scriptedProducer(t, `touch "$PIDDIR/ran"`, 10*time.Second, 1024)
	before := procbound.Outstanding()
	var releases []func()
	for {
		release, err := procbound.Admit()
		if err != nil {
			if !errors.Is(err, procbound.ErrCleanupBacklog) {
				t.Fatalf("filling the procbound slots: want ErrCleanupBacklog at the cap; got %v", err)
			}
			break
		}
		releases = append(releases, release)
	}
	defer func() {
		for _, release := range releases {
			release()
		}
	}()
	_, err := p.runBounded(context.Background(), dir, runEnv(dir), 1024)
	if !errors.Is(err, procbound.ErrCleanupBacklog) {
		t.Fatalf("a full backlog must refuse the spawn: want ErrCleanupBacklog; got %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "ran")); statErr == nil {
		t.Fatalf("the refused checker ran anyway")
	}
	for _, release := range releases {
		release()
	}
	releases = nil
	if got := procbound.Outstanding(); got != before {
		t.Fatalf("after releasing, outstanding = %d; want the pre-test %d", got, before)
	}
}

// pinnedChecker pins the real AILANG_BIN by its bytes and its version line. A
// binary-dependent arm fails loudly when the variable is unset and never skips.
func pinnedChecker(t *testing.T) CheckerConfig {
	t.Helper()
	bin := os.Getenv("AILANG_BIN")
	if bin == "" {
		t.Fatal("AILANG_BIN unset: the producer's binary arms need the pinned v0.41.0 release; never skip")
	}
	out, err := exec.Command(bin, "--version").Output()
	if err != nil {
		t.Fatalf("pinned binary --version: %v", err)
	}
	if line := strings.SplitN(string(out), "\n", 2)[0]; line != "AILANG v0.41.0" {
		t.Fatalf("AILANG_BIN reports %q; want AILANG v0.41.0", line)
	}
	data, err := os.ReadFile(bin)
	if err != nil {
		t.Fatal(err)
	}
	return CheckerConfig{ExecPath: bin, Compiler: hashref.SumSHA256(data), CompilerVersion: "AILANG v0.41.0"}
}

func pinnedProducer(t *testing.T, execTimeout time.Duration, maxOut int64) *Producer {
	t.Helper()
	p, err := NewProducer(internalTestKey, nopSeams{}, nopSeams{}, pinnedChecker(t), ProducerConfig{
		ObjectReadTimeout:  time.Second,
		ObjectWriteTimeout: time.Second,
		ExecTimeout:        execTimeout,
		MaxOutputBytes:     maxOut,
	})
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestProducerRefusesStdoutOverflowIndependently(t *testing.T) {
	source, err := os.ReadFile("testdata/aicheck/verifiedpass.ail")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	t.Run("real-binary-overflow", func(t *testing.T) {
		t.Setenv("TMPDIR", t.TempDir())
		out, err := pinnedProducer(t, 2*time.Minute, 16).runCheck(ctx, source)
		if !errors.Is(err, ErrOutputLimit) || out != nil {
			t.Fatalf("real ai-check output over a 16-byte cap: want ErrOutputLimit and nil bytes; got %d bytes, err=%v", len(out), err)
		}
	})
	t.Run("real-binary-control", func(t *testing.T) {
		tmp := t.TempDir()
		t.Setenv("TMPDIR", tmp)
		out, err := pinnedProducer(t, 2*time.Minute, 1<<20).runCheck(ctx, source)
		if err != nil {
			t.Fatalf("under the 1 MiB cap the real run must return its JSON: %v", err)
		}
		var doc struct {
			Verify struct {
				Results []struct {
					Function string `json:"function"`
					Status   string `json:"status"`
				} `json:"results"`
			} `json:"verify"`
		}
		if err := json.Unmarshal(out, &doc); err != nil {
			t.Fatalf("control output is not ai-check JSON: %v", err)
		}
		verified := false
		for _, r := range doc.Verify.Results {
			verified = verified || (r.Function == "twicePlusOne" && r.Status == "verified")
		}
		if !verified {
			t.Fatalf("control output lacks twicePlusOne verified: %s", out)
		}
		entries, err := os.ReadDir(tmp)
		if err != nil {
			t.Fatal(err)
		}
		if len(entries) != 0 {
			t.Fatalf("the producer left %d entries in its TMPDIR; scratch must be removed", len(entries))
		}
	})
	t.Run("scripted-stdout-overflow", func(t *testing.T) {
		p, _ := scriptedProducer(t, `head -c 64 /dev/zero | tr '\0' o`, 10*time.Second, 16)
		out, err := p.runCheck(ctx, source)
		if !errors.Is(err, ErrOutputLimit) || out != nil {
			t.Fatalf("64 stdout bytes over a 16-byte cap: want ErrOutputLimit and nil bytes; got %d bytes, err=%v", len(out), err)
		}
		if !strings.Contains(err.Error(), "stdout") {
			t.Fatalf("detail must name the stdout stream: %v", err)
		}
	})
}
