package store

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

func TestQuarantineRejectsOperationsUntilWorkerSettles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "world.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	release := make(chan struct{})
	entered := make(chan struct{})
	old := durableAfterCommitHook
	durableAfterCommitHook = func() { close(entered); <-release }
	defer func() { releaseOnce(release); durableAfterCommitHook = old; _ = s.Close() }()
	payload := []byte("quarantine worker")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	err = s.PutObject(ctx, Object{Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte("interface")), SemanticID: "test", Payload: payload})
	if !IsUncertain(err) {
		t.Fatalf("PutObject = %v, want uncertain", err)
	}
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("commit worker did not reach stall")
	}
	s.Quarantine()
	guardCtx, guardCancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer guardCancel()
	readErr, _ := boundedCall(t, func() error {
		_, _, err := s.GetObject(guardCtx, hashref.SumSHA256(payload))
		return err
	}, func() { releaseOnce(release) })
	if !errors.Is(readErr, ErrQuarantined) {
		t.Fatalf("GetObject(deadline) = %v", readErr)
	}
	if _, _, err := s.GetObject(nil, hashref.SumSHA256(payload)); !errors.Is(err, ErrQuarantined) {
		t.Fatalf("GetObject(nil) = %v", err)
	}
	if _, err := s.beginDurable(nil, "test"); !errors.Is(err, ErrQuarantined) {
		t.Fatalf("beginDurable(nil) = %v", err)
	}
	settled := make(chan struct{})
	go func() { _ = s.Close(); close(settled) }()
	if other, err := Open(path); err == nil {
		_ = other.Close()
		t.Fatal("second writer opened before worker settled")
	} else if !IsWriterAlreadyActive(err) {
		t.Fatalf("second Open = %v", err)
	}
	select {
	case <-settled:
		t.Fatal("cleanup settled before worker")
	default:
	}
	close(release)
	select {
	case <-settled:
	case <-time.After(2 * time.Second):
		t.Fatal("cleanup did not settle")
	}
	other, err := Open(path)
	if err != nil {
		t.Fatalf("Open after settlement: %v", err)
	}
	_ = other.Close()
}
