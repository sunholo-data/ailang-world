package broker

import (
	"strings"
	"testing"
	"time"
)

// S7: the seed decisions are a pure table (row 257 §2.3, §2.4).
func TestSeedNextIsPure(t *testing.T) {
	next := []struct {
		attempts int
		running  bool
		verified bool
		want     seedAction
	}{
		{0, false, true, seedReady},
		{3, false, true, seedReady},
		{1, true, true, seedReady},
		{1, true, false, seedJoin},
		{3, true, false, seedJoin},
		{0, false, false, seedStart},
		{2, false, false, seedStart},
		{3, false, false, seedGiveUp},
		{4, false, false, seedGiveUp},
	}
	for _, c := range next {
		if got := execSeedNext(seedView{attempts: c.attempts, running: c.running}, c.verified); got != c.want {
			t.Errorf("execSeedNext(%d,%v,%v) = %v, want %v", c.attempts, c.running, c.verified, got, c.want)
		}
	}

	state := []struct {
		attempts            int
		running, ok, failed bool
		want                string
	}{
		{1, true, false, false, ExecSeedCopying},
		{0, false, true, false, ExecSeedReady},
		{1, false, true, false, ExecSeedReady},
		{1, false, false, true, ExecSeedFailed},
		{2, false, false, true, ExecSeedFailed},
		{3, false, false, true, ExecSeedGaveUp},
		{0, false, false, false, ""},
	}
	for _, c := range state {
		if got := seedStateName(c.attempts, c.running, c.ok, c.failed); got != c.want {
			t.Errorf("seedStateName(%d,%v,%v,%v) = %q, want %q", c.attempts, c.running, c.ok, c.failed, got, c.want)
		}
	}

	if got := seedAttemptDetail("GOCACHE", "the copy exited 1", 2); got != "GOCACHE: the copy exited 1 (attempt 2 of 3)" {
		t.Errorf("seedAttemptDetail = %q", got)
	}
	if got := seedStillCopyingDetail("GOCACHE", "bind", 4210*time.Millisecond); got != "GOCACHE: still copying since bind (4.2s); retry" {
		t.Errorf("seedStillCopyingDetail = %q", got)
	}
	if got := seedStagingName("proj", "GOCACHE", "1-2", 3); got != ".proj.GOCACHE.seeding-1-2-3" {
		t.Errorf("seedStagingName = %q", got)
	}
}

// S12: the marker is exact bytes, versioned, named, and size-capped (§2.1).
func TestSeedMarkerValidIsPure(t *testing.T) {
	if got := string(seedMarkerBytes("GOCACHE")); got != "world-exec-seed v1 GOCACHE\n" {
		t.Fatalf("seedMarkerBytes = %q", got)
	}
	if !seedMarkerValid("GOCACHE", []byte("world-exec-seed v1 GOCACHE\n")) {
		t.Fatal("the exact marker must be valid")
	}
	bad := []string{
		"world-exec-seed v1 GOCACHE\nx",
		"world-exec-seed v1 GOCACHE",
		"world-exec-seed v0 GOCACHE\n",
		"world-exec-seed v2 GOCACHE\n",
		"world-exec-seed v1 OTHER\n",
		"",
		strings.Repeat("x", 65),
	}
	for _, s := range bad {
		if seedMarkerValid("GOCACHE", []byte(s)) {
			t.Errorf("seedMarkerValid accepted %q", s)
		}
	}
	if seedMarkerValid("GOCACHE", nil) {
		t.Error("seedMarkerValid accepted nil")
	}
	n := strings.Repeat("n", 45)
	if got := len(seedMarkerBytes(n)); got != 65 {
		t.Fatalf("len(seedMarkerBytes(45-byte name)) = %d, want 65", got)
	}
	if seedMarkerValid(n, seedMarkerBytes(n)) {
		t.Error("the size cap must reject a 65-byte marker")
	}
}
