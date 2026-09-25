package daemon

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/registry"
	"github.com/sunholo-data/ailang-world/host/store"
)

// TestProductionBusyTimeoutConfiguredBelowReadDeadline is the reorder test: New
// on a plain path must succeed, and the store it opened must report a configured
// lock-retry window numerically below the configured read deadline. Retune
// either constant into the wrong order and New refuses here. It is a
// configuration check, not a runtime bound on a lock-blocked read.
func TestProductionBusyTimeoutConfiguredBelowReadDeadline(t *testing.T) {
	d, err := New(Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: DefaultBindHost})
	if err != nil {
		t.Fatalf("New on the production defaults: %v", err)
	}
	defer d.Close()
	window := d.store.BusyTimeout()
	if window <= 0 {
		t.Fatalf("store BusyTimeout = %s; the ordering below would pass vacuously", window)
	}
	if window >= d.readDeadline {
		t.Fatalf("store busy_timeout %s is not below the read deadline %s", window, d.readDeadline)
	}
}

// TestNewRefusesBusyTimeoutAtOrAboveReadDeadline drives the window through the
// DSN, so the check must read the OPENED store's value, not a constant.
func TestNewRefusesBusyTimeoutAtOrAboveReadDeadline(t *testing.T) {
	for _, tc := range []struct {
		name    string
		window  time.Duration
		refused bool
	}{
		{"disabled (negative)", -time.Millisecond, false},
		{"disabled (zero)", 0, false},
		{"just below", readDeadline - time.Millisecond, false},
		{"equal", readDeadline, true},
		{"above", readDeadline + time.Millisecond, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "world.db")
			dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(%d)", path, tc.window.Milliseconds())
			d, err := New(Config{DBPath: dsn, BindHost: DefaultBindHost})
			if !tc.refused {
				if err != nil {
					t.Fatalf("New refused an ordered window %s: %v", tc.window, err)
				}
				if got := d.store.BusyTimeout(); got != tc.window {
					t.Fatalf("store BusyTimeout = %s, want the DSN's %s — the arm is not driving the window", got, tc.window)
				}
				// Control for the no-side-effect assertion below: an accepted
				// startup DOES bootstrap the registry head.
				if _, ok, err := d.store.GetRegistryHead(context.Background(), registry.SemanticID); err != nil || !ok {
					t.Fatalf("control: accepted startup has no registry head (ok=%v err=%v)", ok, err)
				}
				_ = d.Close()
				return
			}
			if err == nil {
				_ = d.Close()
				t.Fatalf("New accepted busy_timeout %s against read deadline %s", tc.window, readDeadline)
			}
			var se *StartupError
			if !errors.As(err, &se) || se.Stage != StageStoreOpen {
				t.Fatalf("New error = %v, want a StartupError at %q", err, StageStoreOpen)
			}
			if !errors.Is(err, ErrUnorderedTimeouts) {
				t.Fatalf("New error = %v, want ErrUnorderedTimeouts", err)
			}
			// The refusal names both configured durations, in order: the
			// operator must see which window collided with which deadline.
			if want := fmt.Sprintf("busy_timeout %s must be below read deadline %s", tc.window, readDeadline); !strings.Contains(err.Error(), want) {
				t.Fatalf("New error = %q, want it to name %q", err, want)
			}
			// A refused startup must not strand writer authority, and must refuse
			// BEFORE any lifecycle step writes to the store.
			s, err := store.Open(dsn)
			if err != nil {
				t.Fatalf("store.Open after the refusal: %v — New stranded writer authority", err)
			}
			defer s.Close()
			if _, ok, err := s.GetRegistryHead(context.Background(), registry.SemanticID); err != nil || ok {
				t.Fatalf("refused startup left a registry head (ok=%v err=%v) — the check ran after a write", ok, err)
			}
		})
	}
}

// TestCheckReadOrderingAcceptsDisabledWindow: a window <= 0 disables SQLite's
// busy handler, so a lock conflict fails immediately; it is accepted.
func TestCheckReadOrderingAcceptsDisabledWindow(t *testing.T) {
	for _, w := range []time.Duration{-time.Millisecond, 0} {
		if err := checkReadOrdering(w, readDeadline); err != nil {
			t.Fatalf("checkReadOrdering(%s) = %v; a disabled busy handler is ordered", w, err)
		}
	}
}
