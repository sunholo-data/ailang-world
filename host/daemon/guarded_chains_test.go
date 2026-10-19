package daemon

import (
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/registry"
	"github.com/sunholo-data/ailang-world/host/store"
)

func TestGuardedStartupChain(t *testing.T) {
	for _, bounded := range []bool{false, true} {
		s, err := store.Open(filepath.Join(t.TempDir(), "world.db"))
		if err != nil {
			t.Fatal(err)
		}
		d := &Daemon{store: s}
		ref := hashref.SumSHA256([]byte("bootstrap"))
		d.bootstrap = func(ctx context.Context, actual *store.Store, _ string) (registry.Registry, hashref.HashRef, error) {
			err := actual.SetRegistryHead(ctx, "bootstrap", ref)
			return registry.Registry{}, ref, err
		}
		var ctx context.Context = context.Background()
		if bounded {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, time.Second)
			defer cancel()
		}
		err = d.bootstrapRegistry(ctx, unpinnedRelease)
		if bounded {
			if err != nil {
				t.Fatalf("bounded bootstrap = %v", err)
			}
			_ = s.Close()
		} else if !errors.Is(err, store.ErrNoDeadline) {
			t.Fatalf("unbounded bootstrap = %v", err)
		}
	}
}

func TestGuardedCredentialLookupChain(t *testing.T) {
	dir := t.TempDir()
	s, err := store.Open(filepath.Join(dir, "world.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	credPath := filepath.Join(dir, store.CredentialsFileName)
	if err := ensureCredentialsStore(credPath, io.Discard); err != nil {
		t.Fatal(err)
	}
	// The production resolver shape (row 189): a read-only open per lookup.
	d := &Daemon{store: s, credentialsPath: credPath}
	d.resolver = authority.New(authority.ReadOnlyAt(d.credentialsPath))
	header := "Bearer " + strings.Repeat("a", 64)
	if _, err := d.resolver.ResolveContext(context.Background(), header, 0); !errors.Is(err, store.ErrNoDeadline) {
		t.Fatalf("unbounded lookup = %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if _, err := d.resolver.ResolveContext(ctx, header, 0); err != nil {
		t.Fatalf("bounded lookup = %v", err)
	}
}
