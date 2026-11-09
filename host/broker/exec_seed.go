package broker

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

// Row 257 (design_docs/planned/w-exec-cache-seed-at-bind.md r3): the exec
// cache seed is copied on its own clock, started at bind (StartSeed) or by a
// call that finds no copy, never under a call's deadline. A copy is published
// only by an atomic rename of a staging tree that holds World's completion
// marker; readiness is the marker, never existence.

const (
	execSeedBindBudget  = 60 * time.Second // §3.4
	execSeedMaxAttempts = 3                // §3.3
	seedMarkerName      = ".world-seed-complete"
	seedMarkerMax       = 64 // bytes
	seedCleanupMax      = 8  // abandoned staging dirs removed per attempt (§3.1 step 3)
	seedGiveUpSuffix    = "; 3 of 3 attempts failed: fix the seed, then restart serve"
)

// The live seeding states (§3.6).
const (
	ExecSeedCopying = "copying"
	ExecSeedReady   = "ready"
	ExecSeedFailed  = "failed"
	ExecSeedGaveUp  = "gave-up"
)

// seedIncarnation names this daemon process in staging paths: pid and start
// time in nanoseconds (§3.1 step 3).
var seedIncarnation = strconv.Itoa(os.Getpid()) + "-" + strconv.FormatInt(time.Now().UnixNano(), 10)

type execSeed struct {
	mu       sync.Mutex
	attempts int             // copies started, ≤ execSeedMaxAttempts
	running  chan struct{}   // non-nil while a copy runs; closed when it ends
	started  time.Time       // of the running or last attempt
	took     time.Duration   // of the last finished attempt
	by       string          // "bind" or "call": who started the last attempt
	last     *HandlerFailure // nil after a success
	ok       bool            // the last attempt succeeded (or StartSeed found every cache verified)
}

type seedAction int

const (
	seedReady  seedAction = iota // every seeded cache is verified: return nil
	seedJoin                     // a copy is running: wait for it or the call's ctx
	seedStart                    // start attempt attempts+1 in the background, then join
	seedGiveUp                   // the cap is reached: refuse, spawn nothing
)

type seedView struct {
	attempts int
	running  bool
}

// ExecSeedStatus is one seeded cache's live seeding state (row 257 §3.6).
type ExecSeedStatus struct {
	Project  string
	Cache    string
	State    string // ExecSeedCopying | ExecSeedReady | ExecSeedFailed | ExecSeedGaveUp
	Attempts int
	By       string
	Took     time.Duration
	Failure  string // the HandlerFailure line; "" when none
}

// seedMarkerBytes is the completion marker World writes into a staged copy.
func seedMarkerBytes(name string) []byte {
	return []byte("world-exec-seed v1 " + name + "\n")
}

// seedMarkerValid reports whether data is exactly the marker for name.
func seedMarkerValid(name string, data []byte) bool {
	return len(data) <= seedMarkerMax && bytes.Equal(data, seedMarkerBytes(name))
}

// execSeedNext decides what a call does next (§2.3).
func execSeedNext(s seedView, allVerified bool) seedAction {
	switch {
	case allVerified:
		return seedReady
	case s.running:
		return seedJoin
	case s.attempts < execSeedMaxAttempts:
		return seedStart
	default:
		return seedGiveUp
	}
}

// seedStateName maps the seed bookkeeping to the live state (§2.4).
func seedStateName(attempts int, running, ok, failed bool) string {
	switch {
	case running:
		return ExecSeedCopying
	case ok:
		return ExecSeedReady
	case failed && attempts >= execSeedMaxAttempts:
		return ExecSeedGaveUp
	case failed:
		return ExecSeedFailed
	}
	return "" // nothing started
}

// seedAttemptDetail is the Detail of a failure of attempt k.
func seedAttemptDetail(name, what string, k int) string {
	return name + ": " + what + fmt.Sprintf(" (attempt %d of %d)", k, execSeedMaxAttempts)
}

// seedStillCopyingDetail is the Detail of a call that stopped waiting on a copy.
func seedStillCopyingDetail(name, by string, elapsed time.Duration) string {
	return name + ": still copying since " + by + " (" + elapsed.Round(100*time.Millisecond).String() + "); retry"
}

// seedStagingPrefix is the staging-directory name prefix for one cache.
func seedStagingPrefix(project, name string) string {
	return "." + project + "." + name + ".seeding-"
}

// seedStagingName is the staging-directory name of attempt k.
func seedStagingName(project, name, incarnation string, k int) string {
	return seedStagingPrefix(project, name) + incarnation + "-" + strconv.Itoa(k)
}

// defaultSeedCommand is the copy the seed uses: clonefile on darwin, reflink
// where the filesystem has it elsewhere (V21).
func defaultSeedCommand(seed, dst string) (string, []string) {
	if hostGOOS() == "darwin" {
		return "/bin/cp", []string{"-cR", seed, dst}
	}
	return "/bin/cp", []string{"-a", "--reflink=auto", seed, dst}
}
