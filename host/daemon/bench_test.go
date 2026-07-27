package daemon

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"path/filepath"
	"sort"
	"testing"
	"time"

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

	var observed hashref.HashRef
	var previousLog hashref.HashRef
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
