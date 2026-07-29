package daemon

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

func percentileMS(samples []time.Duration, percentile float64) float64 {
	if len(samples) == 0 {
		return 0
	}
	ordered := append([]time.Duration(nil), samples...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i] < ordered[j] })
	index := int(math.Ceil(percentile*float64(len(ordered)))) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(ordered) {
		index = len(ordered) - 1
	}
	return float64(ordered[index]) / float64(time.Millisecond)
}

func reportPercentiles(b *testing.B, samples []time.Duration) {
	b.Helper()
	b.ReportMetric(percentileMS(samples, 0.50), "p50_ms")
	b.ReportMetric(percentileMS(samples, 0.95), "p95_ms")
}

func BenchmarkStoreCommit(b *testing.B) {
	dbPath := filepath.Join(b.TempDir(), "world.db")
	s, err := store.Open(dbPath)
	if err != nil {
		b.Fatalf("store.Open: %v", err)
	}
	b.Cleanup(func() {
		if err := s.Close(); err != nil {
			b.Errorf("store.Close: %v", err)
		}
	})

	// ObservedHead stays the zero value: it is the kernel's single zero-legal
	// COMPARED field, which is exactly how a genesis commit is expressed.
	var observed hashref.HashRef

	// previousLog must be a REAL content address from the first iteration on.
	// It was previously the zero HashRef, which made entry 0 carry a zero
	// PrevEntryHash — precisely the CF-B-2 poison this milestone now refuses at
	// the write path. That commit was never legal: M1's genesis convention seeds
	// entry 0's PrevEntryHash from the genesis world's LogHead, a real content
	// address (store_test.go:103), and a zero prev is unreadable by GetLogEntry.
	// The benchmark only ever "passed" because Commit validated nothing, so this
	// is a benchmark that was relying on the defect, not a regression in the fix.
	previousLog := hashref.SumSHA256([]byte("bench-genesis-prev"))
	samples := make([]time.Duration, 0, b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nextWorld := hashref.SumSHA256([]byte(fmt.Sprintf("bench-world-%d", i)))
		nextLog := hashref.SumSHA256([]byte(fmt.Sprintf("bench-log-%d", i)))
		commit := store.Commit{
			ObservedHead: observed,
			NextWorld: store.World{
				Ref:       nextWorld,
				Revision:  int64(i),
				StateRoot: hashref.SumSHA256([]byte(fmt.Sprintf("bench-state-%d", i))),
				LogHead:   nextLog,
			},
			Entry: store.LogEntry{
				Header: store.LogHeader{
					EntryIndex:     int64(i),
					SemanticsEpoch: 1,
					TransitionFn:   hashref.SumSHA256([]byte("bench-transition")),
					Interpreter:    hashref.SumSHA256([]byte("bench-interpreter")),
					PrevEntryHash:  previousLog,
					WrittenBy:      "benchmark",
				},
				EntryHash:     nextLog,
				TransitionRef: hashref.SumSHA256([]byte(fmt.Sprintf("bench-transition-body-%d", i))),
			},
		}
		start := time.Now()
		if err := s.Commit(commit); err != nil {
			b.Fatalf("Commit #%d: %v", i, err)
		}
		samples = append(samples, time.Since(start))
		observed = nextWorld
		previousLog = nextLog
	}
	b.StopTimer()
	reportPercentiles(b, samples)
}

func BenchmarkJournalAppend(b *testing.B) {
	dbPath := filepath.Join(b.TempDir(), "world.db")
	s, err := store.Open(dbPath)
	if err != nil {
		b.Fatalf("store.Open: %v", err)
	}
	b.Cleanup(func() {
		if err := s.Close(); err != nil {
			b.Errorf("store.Close: %v", err)
		}
	})

	samples := make([]time.Duration, 0, b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := fmt.Sprintf("bench-journal-%d", i)
		intent := store.JournalIntent{
			InvocationID:  id,
			WorldRef:      hashref.SumSHA256([]byte(fmt.Sprintf("bench-journal-world-%d", i))),
			EntryHash:     hashref.SumSHA256([]byte(fmt.Sprintf("bench-journal-entry-%d", i))),
			PrevEntryHash: hashref.SumSHA256([]byte(fmt.Sprintf("bench-journal-prev-%d", i))),
			TransitionFn:  hashref.SumSHA256([]byte("bench-journal-transition")),
			TransitionRef: hashref.SumSHA256([]byte(fmt.Sprintf("bench-journal-body-%d", i))),
			Interpreter:   hashref.SumSHA256([]byte("bench-journal-interpreter")),
			LogicalTime:   int64(i),
		}
		start := time.Now()
		if _, _, err := s.AppendIntent(id, intent); err != nil {
			b.Fatalf("AppendIntent #%d: %v", i, err)
		}
		samples = append(samples, time.Since(start))
	}
	b.StopTimer()
	reportPercentiles(b, samples)
}

func BenchmarkCommitWithReceipt(b *testing.B) {
	dbPath := filepath.Join(b.TempDir(), "world.db")
	s, err := store.Open(dbPath)
	if err != nil {
		b.Fatalf("store.Open: %v", err)
	}
	b.Cleanup(func() {
		if err := s.Close(); err != nil {
			b.Errorf("store.Close: %v", err)
		}
	})

	var observed hashref.HashRef
	previousLog := hashref.SumSHA256([]byte("bench-receipt-genesis-prev"))
	commits := make([]store.Commit, b.N)
	for i := 0; i < b.N; i++ {
		id := fmt.Sprintf("bench-receipt-%d", i)
		nextWorld := hashref.SumSHA256([]byte(fmt.Sprintf("bench-receipt-world-%d", i)))
		nextLog := hashref.SumSHA256([]byte(fmt.Sprintf("bench-receipt-log-%d", i)))
		body := hashref.SumSHA256([]byte(fmt.Sprintf("bench-receipt-body-%d", i)))
		commits[i] = store.Commit{
			InvocationID: id,
			ObservedHead: observed,
			NextWorld: store.World{
				Ref: nextWorld, Revision: int64(i),
				StateRoot: hashref.SumSHA256([]byte(fmt.Sprintf("bench-receipt-state-%d", i))),
				LogHead:   nextLog,
			},
			Entry: store.LogEntry{
				Header: store.LogHeader{
					EntryIndex: int64(i), SemanticsEpoch: 1,
					TransitionFn:  hashref.SumSHA256([]byte("bench-receipt-transition")),
					Interpreter:   hashref.SumSHA256([]byte("bench-receipt-interpreter")),
					PrevEntryHash: previousLog, WrittenBy: "benchmark",
				},
				EntryHash: nextLog, TransitionRef: body,
			},
		}
		intent := store.JournalIntent{
			InvocationID: id, WorldRef: nextWorld, EntryHash: nextLog,
			ObservedHead: observed, PrevEntryHash: previousLog,
			TransitionFn: commits[i].Entry.Header.TransitionFn, TransitionRef: body,
			Interpreter: commits[i].Entry.Header.Interpreter, LogicalTime: int64(i),
		}
		if _, _, err := s.AppendIntent(id, intent); err != nil {
			b.Fatalf("stage AppendIntent #%d: %v", i, err)
		}
		observed, previousLog = nextWorld, nextLog
	}

	samples := make([]time.Duration, 0, b.N)
	b.ResetTimer()
	for i := range commits {
		start := time.Now()
		if err := s.Commit(commits[i]); err != nil {
			b.Fatalf("Commit with receipt #%d: %v", i, err)
		}
		samples = append(samples, time.Since(start))
	}
	b.StopTimer()
	reportPercentiles(b, samples)
}

func benchmarkDaemonGET(b *testing.B, route string, seed bool) {
	b.Helper()
	dbPath := filepath.Join(b.TempDir(), "world.db")
	d, err := New(Config{DBPath: dbPath, BindHost: DefaultBindHost, BindPort: 0})
	if err != nil {
		b.Fatalf("New: %v", err)
	}
	if seed {
		w := store.World{
			Ref:       hashref.SumSHA256([]byte("daemon-benchmark-world-genesis")),
			Revision:  0,
			StateRoot: hashref.SumSHA256([]byte("daemon-benchmark-state-genesis")),
			LogHead:   hashref.SumSHA256([]byte("daemon-benchmark-log-genesis")),
		}
		if err := d.store.PutWorld(w); err != nil {
			_ = d.Close()
			b.Fatalf("seed PutWorld: %v", err)
		}
		if err := d.store.SelectHead(w.Ref); err != nil {
			_ = d.Close()
			b.Fatalf("seed SelectHead: %v", err)
		}
	}
	if err := d.Listen(); err != nil {
		_ = d.Close()
		b.Fatalf("Listen: %v", err)
	}
	serveDone := make(chan error, 1)
	go func() { serveDone <- d.Serve() }()
	b.Cleanup(func() {
		if err := d.Shutdown(); err != nil {
			b.Errorf("Shutdown: %v", err)
		}
		if err := <-serveDone; err != nil && err != http.ErrServerClosed {
			b.Errorf("Serve: %v", err)
		}
		if err := d.Close(); err != nil {
			b.Errorf("Close: %v", err)
		}
	})

	client := &http.Client{Timeout: 10 * time.Second}
	doGET := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, d.URL()+route, nil)
		if err != nil {
			b.Fatalf("build GET %s: %v", route, err)
		}
		resp, err := client.Do(req)
		if err != nil {
			b.Fatalf("GET %s: %v", route, err)
		}
		if _, err := io.Copy(io.Discard, resp.Body); err != nil {
			_ = resp.Body.Close()
			b.Fatalf("drain GET %s body: %v", route, err)
		}
		if err := resp.Body.Close(); err != nil {
			b.Fatalf("close GET %s body: %v", route, err)
		}
		if resp.StatusCode != http.StatusOK {
			b.Fatalf("GET %s status = %d, want %d", route, resp.StatusCode, http.StatusOK)
		}
	}

	// Warm the listener and the client's keep-alive connection outside the
	// measured region.
	doGET()
	samples := make([]time.Duration, 0, b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := time.Now()
		doGET()
		samples = append(samples, time.Since(start))
	}
	b.StopTimer()
	reportPercentiles(b, samples)
}

func BenchmarkHeadRead(b *testing.B) {
	benchmarkDaemonGET(b, "/v1/head", true)
}

func BenchmarkHealth(b *testing.B) {
	benchmarkDaemonGET(b, "/v1/health", false)
}

func BenchmarkRESTCommit(b *testing.B) {
	dbPath := filepath.Join(b.TempDir(), "world.db")
	d, err := New(Config{DBPath: dbPath, BindHost: DefaultBindHost, BindPort: 0})
	if err != nil {
		b.Fatalf("New: %v", err)
	}
	genesis := store.World{
		Ref: hashref.SumSHA256([]byte("rest-bench-genesis")), Revision: 0,
		StateRoot: hashref.SumSHA256([]byte("rest-bench-genesis-state")),
		LogHead:   hashref.SumSHA256([]byte("rest-bench-genesis-log")),
	}
	if err := d.store.PutWorld(genesis); err != nil {
		b.Fatalf("PutWorld genesis: %v", err)
	}
	if err := d.store.SelectHead(genesis.Ref); err != nil {
		b.Fatalf("SelectHead genesis: %v", err)
	}
	if err := d.Listen(); err != nil {
		_ = d.Close()
		b.Fatalf("Listen: %v", err)
	}
	serveDone := make(chan error, 1)
	go func() { serveDone <- d.Serve() }()
	b.Cleanup(func() {
		if err := d.Shutdown(); err != nil {
			b.Errorf("Shutdown: %v", err)
		}
		if err := <-serveDone; err != nil && err != http.ErrServerClosed {
			b.Errorf("Serve: %v", err)
		}
		if err := d.Close(); err != nil {
			b.Errorf("Close: %v", err)
		}
	})

	client := &http.Client{Timeout: 10 * time.Second}
	current := genesis
	samples := make([]time.Duration, 0, b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		commit := testCommit(current, int64(i), fmt.Sprintf("rest-bench-%d", i))
		body := encodeCommit(commit)
		req, err := http.NewRequest(http.MethodPost, d.URL()+"/v1/commit", bytes.NewReader(body))
		if err != nil {
			b.Fatalf("build POST: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		start := time.Now()
		resp, err := client.Do(req)
		if err != nil {
			b.Fatalf("POST commit #%d: %v", i, err)
		}
		if _, err := io.Copy(io.Discard, resp.Body); err != nil {
			_ = resp.Body.Close()
			b.Fatalf("drain POST body: %v", err)
		}
		if err := resp.Body.Close(); err != nil {
			b.Fatalf("close POST body: %v", err)
		}
		samples = append(samples, time.Since(start))
		if resp.StatusCode != http.StatusOK {
			b.Fatalf("POST commit #%d status=%d", i, resp.StatusCode)
		}
		current = commit.NextWorld
	}
	b.StopTimer()
	reportPercentiles(b, samples)
}

func BenchmarkLogRange(b *testing.B) {
	for _, limit := range []int{100, 500} {
		b.Run(fmt.Sprintf("limit_%d", limit), func(b *testing.B) {
			dbPath := filepath.Join(b.TempDir(), "world.db")
			d, err := New(Config{DBPath: dbPath, BindHost: DefaultBindHost, BindPort: 0})
			if err != nil {
				b.Fatalf("New: %v", err)
			}
			current := store.World{
				Ref: hashref.SumSHA256([]byte("range-bench-genesis")), Revision: 0,
				StateRoot: hashref.SumSHA256([]byte("range-bench-genesis-state")),
				LogHead:   hashref.SumSHA256([]byte("range-bench-genesis-log")),
			}
			if err := d.store.PutWorld(current); err != nil {
				b.Fatalf("PutWorld genesis: %v", err)
			}
			if err := d.store.SelectHead(current.Ref); err != nil {
				b.Fatalf("SelectHead genesis: %v", err)
			}
			for i := int64(0); i < 500; i++ {
				commit := testCommit(current, i, fmt.Sprintf("range-bench-%d", i))
				if err := d.store.Commit(commit); err != nil {
					b.Fatalf("seed Commit(%d): %v", i, err)
				}
				current = commit.NextWorld
			}
			if err := d.Listen(); err != nil {
				_ = d.Close()
				b.Fatalf("Listen: %v", err)
			}
			serveDone := make(chan error, 1)
			go func() { serveDone <- d.Serve() }()
			b.Cleanup(func() {
				if err := d.Shutdown(); err != nil {
					b.Errorf("Shutdown: %v", err)
				}
				if err := <-serveDone; err != nil && err != http.ErrServerClosed {
					b.Errorf("Serve: %v", err)
				}
				if err := d.Close(); err != nil {
					b.Errorf("Close: %v", err)
				}
			})

			client := &http.Client{Timeout: 10 * time.Second}
			route := fmt.Sprintf("/v1/log?from=0&limit=%d", limit)
			doGET := func() {
				resp, err := client.Get(d.URL() + route)
				if err != nil {
					b.Fatalf("GET %s: %v", route, err)
				}
				if _, err := io.Copy(io.Discard, resp.Body); err != nil {
					_ = resp.Body.Close()
					b.Fatalf("drain GET body: %v", err)
				}
				if err := resp.Body.Close(); err != nil {
					b.Fatalf("close GET body: %v", err)
				}
				if resp.StatusCode != http.StatusOK {
					b.Fatalf("GET %s status=%d", route, resp.StatusCode)
				}
			}
			doGET() // warm listener and keep-alive outside the measured window
			samples := make([]time.Duration, 0, b.N)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				start := time.Now()
				doGET()
				samples = append(samples, time.Since(start))
			}
			b.StopTimer()
			reportPercentiles(b, samples)
		})
	}
}

func BenchmarkBrokerDecide(b *testing.B) {
	capability := broker.Capability{
		Effect: broker.EffectFSRead, Scope: "/bench/input", ExpiresAt: int64(b.N) + 1, Budget: int64(b.N),
	}
	samples := make([]time.Duration, 0, b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := broker.EffectRequest{
			Effect: broker.EffectFSRead, Scope: "/bench/input", Cost: 1, Now: int64(i),
		}
		start := time.Now()
		decision := broker.Decide(capability, req)
		samples = append(samples, time.Since(start))
		if !decision.Allowed {
			b.Fatalf("Decide #%d = %#v, want allowed", i, decision)
		}
	}
	b.StopTimer()
	reportPercentiles(b, samples)
}

func BenchmarkBrokerFSRead(b *testing.B) {
	dbPath := filepath.Join(b.TempDir(), "world.db")
	s, err := store.Open(dbPath)
	if err != nil {
		b.Fatalf("store.Open: %v", err)
	}
	b.Cleanup(func() {
		if err := s.Close(); err != nil {
			b.Errorf("store.Close: %v", err)
		}
	})

	inputPath := filepath.Join(b.TempDir(), "input.txt")
	if err := os.WriteFile(inputPath, []byte("broker-fs-read-seed"), 0o600); err != nil {
		b.Fatalf("seed input: %v", err)
	}
	session := broker.NewSession(s, []broker.Capability{{
		Effect: broker.EffectFSRead, Scope: inputPath,
		ExpiresAt: int64(b.N) + 1, Budget: int64(b.N),
	}}, broker.Registry{broker.EffectFSRead: broker.FSHandler{}})

	samples := make([]time.Duration, 0, b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		want := []byte(fmt.Sprintf("broker-fs-read-%d", i))
		if err := os.WriteFile(inputPath, want, 0o600); err != nil {
			b.Fatalf("write input #%d: %v", i, err)
		}
		req := broker.EffectRequest{
			Effect: broker.EffectFSRead, Scope: inputPath, Cost: 1, Now: int64(i),
		}
		b.StartTimer()
		start := time.Now()
		got, recordRef, err := session.Invoke(context.Background(), req, nil)
		samples = append(samples, time.Since(start))
		if err != nil {
			b.Fatalf("Invoke #%d: %v", i, err)
		}
		if !bytes.Equal(got, want) || recordRef.IsZero() {
			b.Fatalf("Invoke #%d returned result %q, record %s", i, got, recordRef.String())
		}
	}
	b.StopTimer()
	reportPercentiles(b, samples)
}
