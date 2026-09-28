package daemon

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/registry"
	"github.com/sunholo-data/ailang-world/host/store"
)

func TestStartupBudgetQuarantinesPostCutoffBootstrap(t *testing.T) {
	path := filepath.Join(t.TempDir(), "world.db")
	s, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	reader, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	held, err := reader.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer held.Rollback()
	var count int
	if err := held.QueryRow("SELECT count(*) FROM objects").Scan(&count); err != nil {
		t.Fatal(err)
	}
	d := &Daemon{store: s}
	payload := []byte("bootstrap post-cutoff")
	object := store.Object{Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte("interface")), SemanticID: "bootstrap", Payload: payload}
	writeResult := make(chan error, 1)
	d.bootstrap = func(ctx context.Context, actual *store.Store, _ string) (registry.Registry, hashref.HashRef, error) {
		if actual != s {
			t.Error("bootstrap did not receive real store")
		}
		err := actual.PutObject(ctx, object)
		writeResult <- err
		return registry.Registry{}, hashref.HashRef{}, err
	}
	outer, stop := context.WithCancel(boundedTestContext(t))
	defer stop()
	begin := time.Now()
	_, err = startBounded(outer, Config{}, 100*time.Millisecond, func(ctx context.Context, _ Config) (*Daemon, error) {
		return nil, d.bootstrapRegistry(ctx, unpinnedRelease)
	})
	var startup *StartupError
	if !errors.As(err, &startup) || startup.Stage != StageRegistry || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("startup = %v", err)
	}
	if time.Since(begin) > time.Second {
		t.Fatal("startup exceeded budget")
	}
	if writeErr := <-writeResult; !store.IsUncertain(writeErr) {
		t.Fatalf("bootstrap write = %v, want uncertain", writeErr)
	}
	if startup.Settled == nil {
		t.Fatal("no settlement signal")
	}
	if err := outer.Err(); err != nil {
		t.Fatalf("serve context canceled by startup timer: %v", err)
	}
	if _, _, err := s.GetRegistryHead(nil, registry.SemanticID); !errors.Is(err, store.ErrQuarantined) {
		t.Fatalf("quarantined read = %v", err)
	}
	if other, err := store.Open(path); err == nil {
		_ = other.Close()
		t.Fatal("second writer opened during cleanup")
	} else if !store.IsWriterAlreadyActive(err) {
		t.Fatalf("Open while worker held = %v", err)
	}
	select {
	case <-startup.Settled:
		t.Fatal("cleanup settled before worker")
	default:
	}
	if err := held.Rollback(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-startup.Settled:
	case <-time.After(3 * time.Second):
		t.Fatal("cleanup did not settle")
	}
	reopened, err := store.Open(path)
	if err != nil {
		t.Fatalf("Open after cleanup: %v", err)
	}
	_ = reopened.Close()
}

func TestStartupBudgetHeldConnectionSettles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "world.db")
	s, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	blocker, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer blocker.Close()
	conn, err := blocker.Conn(boundedTestContext(t))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if _, err := conn.ExecContext(boundedTestContext(t), "BEGIN EXCLUSIVE"); err != nil {
		t.Fatal(err)
	}
	defer conn.ExecContext(boundedTestContext(t), "ROLLBACK")
	readDone := make(chan struct{})
	go func() {
		_, _, _ = s.GetObject(boundedTestContext(t), hashref.SumSHA256([]byte("held")))
		close(readDone)
	}()
	time.Sleep(30 * time.Millisecond)
	d := &Daemon{store: s, bootstrap: registry.Bootstrap}
	begin := time.Now()
	_, err = startBounded(boundedTestContext(t), Config{}, 100*time.Millisecond, func(ctx context.Context, _ Config) (*Daemon, error) {
		return nil, d.bootstrapRegistry(ctx, unpinnedRelease)
	})
	var startup *StartupError
	if !errors.As(err, &startup) || startup.Stage != StageRegistry || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("startup=%v", err)
	}
	if time.Since(begin) > time.Second {
		t.Fatal("held-connection startup exceeded budget")
	}
	select {
	case <-startup.Settled:
	case <-time.After(time.Second):
		t.Fatal("no-worker cleanup stalled behind held connection")
	}
	if _, err := conn.ExecContext(boundedTestContext(t), "ROLLBACK"); err != nil {
		t.Fatal(err)
	}
	select {
	case <-readDone:
	case <-time.After(3 * time.Second):
		t.Fatal("held read did not settle")
	}
	reopened, err := store.Open(path)
	if err != nil {
		t.Fatalf("Open after cleanup: %v", err)
	}
	_ = reopened.Close()
}
