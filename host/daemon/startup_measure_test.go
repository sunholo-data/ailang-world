package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"
)

// TestMeasureStartup is the startup instrument for the bound table in
// design_docs/planned/w-store-bounded-durable-operations.md. It is skipped
// unless WORLD_BOUND_MEASURE=1; it asserts nothing about speed.
func TestMeasureStartup(t *testing.T) {
	if os.Getenv("WORLD_BOUND_MEASURE") != "1" {
		t.Skip("set WORLD_BOUND_MEASURE=1 to measure")
	}
	bin := os.Getenv("AILANG_BIN")
	samples := map[string][]time.Duration{}
	dir := t.TempDir()
	for i := 0; i < 50; i++ {
		path := filepath.Join(dir, fmt.Sprintf("w-%d.db", i))
		for _, arm := range []string{"New(fresh, archive)", "New(existing, archive)"} {
			start := time.Now()
			d, err := New(boundedTestContext(t), Config{DBPath: path, BindHost: DefaultBindHost, AilangBin: bin})
			if err != nil {
				t.Fatalf("%s: %v", arm, err)
			}
			samples[arm] = append(samples[arm], time.Since(start))
			_ = d.Close()
		}
	}
	for arm, d := range samples {
		sort.Slice(d, func(i, j int) bool { return d[i] < d[j] })
		t.Logf("%-25s n=%d p50=%v p99=%v max=%v", arm, len(d), d[len(d)/2], d[int(0.99*float64(len(d)-1))], d[len(d)-1])
	}
}
