package procbound

import (
	"errors"
	"sync"
	"testing"
	"time"
)

// drained fails the test unless every reservation is back within 5 s, and
// returns the (zero) base the exactness assertions below rely on.
func drained(t *testing.T) int64 {
	t.Helper()
	for start := time.Now(); Outstanding() != 0 && time.Since(start) < 5*time.Second; time.Sleep(time.Millisecond) {
	}
	if got := Outstanding(); got != 0 {
		t.Fatalf("reservations held at test start = %d, want 0", got)
	}
	return 0
}

func TestWaitReturnsResultWithinBoundAndReleases(t *testing.T) {
	drained(t)
	release, err := Admit()
	if err != nil {
		t.Fatal(err)
	}
	want := errors.New("exit 3")
	if err := Wait(func() error { return want }, time.Minute, release); err != want {
		t.Fatalf("Wait = %v, want %v", err, want)
	}
	if got := Outstanding(); got != 0 {
		t.Fatalf("Outstanding after a reaped wait = %d, want 0", got)
	}
}

// A wait that outlives the bound is abandoned; its reservation is RETAINED
// until the wait finishes, a full set of reservations refuses admission, and
// each is released exactly once.
func TestAbandonedWaitersRetainReservationsUntilReaped(t *testing.T) {
	drained(t)
	unblock := make(chan struct{})
	for i := int64(0); i < MaxOutstanding; i++ {
		release, err := Admit()
		if err != nil {
			t.Fatalf("Admit %d = %v", i, err)
		}
		if err := Wait(func() error { <-unblock; return nil }, time.Millisecond, release); !errors.Is(err, ErrCleanupIncomplete) {
			t.Fatalf("Wait = %v, want ErrCleanupIncomplete", err)
		}
	}
	if got := Outstanding(); got != MaxOutstanding {
		t.Fatalf("Outstanding = %d, want %d (abandoned waiters keep their reservations)", got, MaxOutstanding)
	}
	if release, err := Admit(); !errors.Is(err, ErrCleanupBacklog) || release != nil {
		t.Fatalf("Admit at cap = (%v, %v), want (nil, ErrCleanupBacklog)", release != nil, err)
	}
	close(unblock)
	drained(t)
	release, err := Admit()
	if err != nil {
		t.Fatalf("Admit after drain = %v", err)
	}
	release()
}

// release is idempotent: a second call must not free a slot it does not own.
func TestReleaseIsIdempotent(t *testing.T) {
	drained(t)
	r1, err := Admit()
	if err != nil {
		t.Fatal(err)
	}
	r2, err := Admit()
	if err != nil {
		t.Fatal(err)
	}
	r1()
	r1()
	if got := Outstanding(); got != 1 {
		t.Fatalf("Outstanding after releasing one reservation twice = %d, want 1", got)
	}
	r2()
	if got := Outstanding(); got != 0 {
		t.Fatalf("Outstanding after both released = %d, want 0", got)
	}
}

// AC10 (unit). Many barrier-released callers race Admit: exactly
// MaxOutstanding succeed, the rest are refused without waiting, and the count
// never exceeds the limit.
func TestAdmitIsExactUnderConcurrency(t *testing.T) {
	drained(t)
	for round := 0; round < 200; round++ {
		const callers = 64
		var wg sync.WaitGroup
		barrier := make(chan struct{})
		releases := make(chan func(), callers)
		var refused, peakViolations int64
		var mu sync.Mutex
		for i := 0; i < callers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-barrier
				release, err := Admit()
				if n := Outstanding(); n > MaxOutstanding {
					mu.Lock()
					peakViolations++
					mu.Unlock()
				}
				if err != nil {
					mu.Lock()
					refused++
					mu.Unlock()
					return
				}
				releases <- release
			}()
		}
		close(barrier)
		wg.Wait()
		close(releases)
		admitted := int64(0)
		for release := range releases {
			admitted++
			release()
		}
		if admitted != MaxOutstanding || refused != callers-MaxOutstanding || peakViolations != 0 {
			t.Fatalf("round %d: admitted=%d refused=%d peakViolations=%d, want %d/%d/0",
				round, admitted, refused, peakViolations, MaxOutstanding, callers-MaxOutstanding)
		}
		if got := Outstanding(); got != 0 {
			t.Fatalf("round %d: Outstanding after release = %d, want 0", round, got)
		}
	}
}
