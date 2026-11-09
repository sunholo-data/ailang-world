package broker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

// seedCacheVerified reports whether dst is a published copy: a directory
// holding World's marker, exactly (§2.5).
func seedCacheVerified(dst, name string) bool {
	fi, err := os.Lstat(dst)
	if err != nil || !fi.IsDir() {
		return false
	}
	m := filepath.Join(dst, seedMarkerName)
	mi, err := os.Lstat(m)
	if err != nil || !mi.Mode().IsRegular() || mi.Size() > seedMarkerMax {
		return false
	}
	data, err := os.ReadFile(m)
	if err != nil {
		return false
	}
	return seedMarkerValid(name, data)
}

// seededCaches lists the caches that have a seed, in byte order.
func (h *ExecHandler) seededCaches() []string {
	var out []string
	for _, name := range sortedKeys(h.profile.Caches) {
		if h.profile.Caches[name].Seed != "" {
			out = append(out, name)
		}
	}
	return out
}

// seedsVerified returns the first seeded cache that is not verified, or
// ("", true) when every seeded cache is.
func (h *ExecHandler) seedsVerified() (string, bool) {
	for _, name := range h.seededCaches() {
		if !seedCacheVerified(filepath.Join(h.cache, name), name) {
			return name, false
		}
	}
	return "", true
}

// StartSeed starts the seed copy at bind, on the copy's own clock (§2.6). It
// is idempotent and spawns nothing when every seeded cache is verified.
func (h *ExecHandler) StartSeed() {
	if len(h.seededCaches()) == 0 {
		return
	}
	_, verified := h.seedsVerified() // I/O outside the lock
	h.seed.mu.Lock()
	defer h.seed.mu.Unlock()
	if h.seed.running != nil || h.seed.attempts > 0 || h.seed.ok {
		return // idempotent
	}
	if verified {
		h.seed.ok, h.seed.by = true, "bind" // a restart over a verified cache: no copy
		return
	}
	h.startAttemptLocked("bind")
}

// startAttemptLocked starts attempt attempts+1; the caller holds h.seed.mu.
func (h *ExecHandler) startAttemptLocked(by string) chan struct{} {
	h.seed.attempts++
	k := h.seed.attempts
	done := make(chan struct{})
	h.seed.running, h.seed.started, h.seed.by = done, time.Now(), by
	go h.runSeedAttempt(k, by, done)
	return done
}

// runSeedAttempt is the goroutine that runs attempt k. It is the ONE context
// root of the seed path: its budget is execSeedBindBudget, never a call's.
func (h *ExecHandler) runSeedAttempt(k int, by string, done chan struct{}) {
	ctx, cancel := context.WithTimeout(context.Background(), h.seedBudget)
	start := time.Now()
	f := h.seedAttempt(ctx, k)
	cancel()
	took := time.Since(start)
	if h.log != nil {
		if f == nil {
			fmt.Fprintf(h.log, "ailang-worldd: exec cache for episode %q (%s): %s seeded at %s in %.2fs\n",
				h.episode, h.profile.Project, strings.Join(h.seededCaches(), ","), by, took.Seconds())
		} else {
			fmt.Fprintf(h.log, "ailang-worldd: exec cache for episode %q (%s): seeding at %s failed: %s\n",
				h.episode, h.profile.Project, by, f.Line())
		}
	}
	h.seed.mu.Lock()
	h.seed.took, h.seed.last, h.seed.ok, h.seed.running = took, f, f == nil, nil
	h.seed.mu.Unlock()
	close(done)
}

// seedAttempt seeds every seeded cache in order; the first failure ends it.
func (h *ExecHandler) seedAttempt(ctx context.Context, k int) *HandlerFailure {
	for _, name := range h.seededCaches() {
		if f := h.seedOne(ctx, k, name, h.profile.Caches[name].Seed); f != nil {
			return f
		}
	}
	return nil
}

// seedOne is steps 0-6 of one attempt (§3.1): copy into unique staging,
// write the marker, then publish by an atomic rename.
func (h *ExecHandler) seedOne(ctx context.Context, k int, name, seed string) *HandlerFailure {
	epDir := filepath.Dir(h.cache)
	dst := filepath.Join(h.cache, name)
	own := seedStagingName(h.profile.Project, name, h.incarnation, k)
	staging := filepath.Join(epDir, own)
	fail := func(what string, err error) *HandlerFailure {
		return &HandlerFailure{Class: FailSeedCache, Detail: seedAttemptDetail(name, what, k), Err: err}
	}
	wrap := func(err error) error { return fmt.Errorf("broker: exec handler: seed cache %s: %w", name, err) }
	// 0
	if err := os.MkdirAll(h.cache, 0o700); err != nil {
		return fail("cache dir cannot be created", wrap(err))
	}
	// 1
	if _, err := os.Lstat(filepath.Join(seed, seedMarkerName)); err == nil {
		return fail("the seed holds World's marker; remove it from the seed",
			fmt.Errorf("broker: exec handler: seed %s holds %s", seed, seedMarkerName))
	}
	// 2
	if seedCacheVerified(dst, name) {
		return nil
	}
	if _, err := os.Lstat(dst); err == nil {
		if err := os.RemoveAll(dst); err != nil {
			return fail("an unverified copy cannot be removed", wrap(err))
		}
	}
	// 3
	h.cleanupStaging(epDir, name, own)
	if err := os.Mkdir(staging, 0o700); err != nil {
		_ = os.RemoveAll(staging)
		return fail("a staging directory cannot be created", wrap(err))
	}
	// 4
	bin, args := h.seedCommand(seed, filepath.Join(staging, "copy"))
	h.spawns.Add(1)
	if _, err := runBounded(ctx, handlerBounds{execTimeout: h.seedBudget}, handlerCommand{
		path: bin, args: args, dir: epDir, env: execHostEnv(),
	}); err != nil {
		_ = os.RemoveAll(staging)
		return &HandlerFailure{Class: FailSeedCache, Detail: seedAttemptDetail(name, seedWhat(err), k),
			Err: fmt.Errorf("broker: exec handler: seed cache %s from %s: %w", name, seed, err)}
	}
	// 5
	if err := os.WriteFile(filepath.Join(staging, "copy", seedMarkerName), seedMarkerBytes(name), 0o600); err != nil {
		_ = os.RemoveAll(staging)
		return fail("the completion marker cannot be written", wrap(err))
	}
	// 6
	if err := os.Rename(filepath.Join(staging, "copy"), dst); err != nil {
		_ = os.RemoveAll(staging)
		return fail("the copy cannot be published", wrap(err))
	}
	_ = os.RemoveAll(staging)
	return nil
}

// cleanupStaging removes abandoned staging dirs of this cache, lexical order,
// at most seedCleanupMax attempts; a failure is logged, never fatal (§1.2 item 6).
func (h *ExecHandler) cleanupStaging(epDir, name, own string) {
	ents, err := os.ReadDir(epDir)
	if err != nil {
		return // ENOENT is normal
	}
	prefix := seedStagingPrefix(h.profile.Project, name)
	removed := 0
	for _, e := range ents {
		if !strings.HasPrefix(e.Name(), prefix) || e.Name() == own {
			continue
		}
		if removed == seedCleanupMax {
			break
		}
		removed++
		if err := os.RemoveAll(filepath.Join(epDir, e.Name())); err != nil && h.log != nil {
			fmt.Fprintf(h.log, "ailang-worldd: exec cache for episode %q (%s): abandoned staging %s not removed: %v\n",
				h.episode, h.profile.Project, e.Name(), err)
		}
	}
}

// awaitSeed is the per-call path (§2.7): it returns nil once every seeded
// cache is verified, and otherwise joins, starts or refuses per execSeedNext.
func (h *ExecHandler) awaitSeed(ctx context.Context) error {
	for {
		name, verified := h.seedsVerified()
		h.seed.mu.Lock()
		act := execSeedNext(seedView{attempts: h.seed.attempts, running: h.seed.running != nil}, verified)
		var done chan struct{}
		switch act {
		case seedReady:
			h.seed.mu.Unlock()
			return nil
		case seedGiveUp:
			f := h.giveUpLocked(name)
			h.seed.mu.Unlock()
			return f
		case seedStart:
			done = h.startAttemptLocked("call")
		case seedJoin:
			done = h.seed.running
		}
		by, started := h.seed.by, h.seed.started
		h.seed.mu.Unlock()
		select {
		case <-done:
		case <-ctx.Done():
			return &HandlerFailure{Class: FailSeedCache, Detail: seedStillCopyingDetail(name, by, time.Since(started)),
				Err: fmt.Errorf("broker: exec handler: seed cache %s: %w", name, ctx.Err())}
		}
		if _, ok := h.seedsVerified(); ok {
			return nil
		}
		h.seed.mu.Lock()
		last := h.seed.last
		h.seed.mu.Unlock()
		if last != nil {
			c := *last
			return &c // the joined attempt's own failure (§1.2 item 3)
		}
		// The joined attempt succeeded but the cache is unverified again: decide again (bounded by the cap).
	}
}

// giveUpLocked is the refusal once the attempt cap is reached (§2.2); the
// caller holds h.seed.mu.
func (h *ExecHandler) giveUpLocked(name string) *HandlerFailure {
	if last := h.seed.last; last != nil {
		return &HandlerFailure{Class: FailSeedCache, Detail: last.Detail + seedGiveUpSuffix, Err: last.Err}
	}
	return &HandlerFailure{Class: FailSeedCache,
		Detail: seedAttemptDetail(name, "the published copy is missing or unverified", execSeedMaxAttempts) + seedGiveUpSuffix,
		Err:    errors.New("broker: exec handler: seed cache " + name + " is unverified after 3 attempts")}
}

// SeedStatus is the live seeding state, one entry per seeded cache (§2.8);
// nil when nothing has started.
func (h *ExecHandler) SeedStatus() []ExecSeedStatus {
	h.seed.mu.Lock()
	attempts, running, ok, last := h.seed.attempts, h.seed.running != nil, h.seed.ok, h.seed.last
	by, took := h.seed.by, h.seed.took
	h.seed.mu.Unlock()
	state := seedStateName(attempts, running, ok, last != nil)
	if state == "" {
		return nil
	}
	failure := ""
	if last != nil {
		failure = last.Line()
	}
	var out []ExecSeedStatus
	for _, name := range h.seededCaches() {
		out = append(out, ExecSeedStatus{Project: h.profile.Project, Cache: name, State: state,
			Attempts: attempts, By: by, Took: took, Failure: failure})
	}
	return out
}
