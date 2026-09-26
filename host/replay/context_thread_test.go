package replay

import (
	"context"
	"errors"
	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"os"
	"path/filepath"
	"testing"
)

func TestReplayCallerCancellation(t *testing.T) {
	s, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	e := NewEngine(s, nil)
	entry := EpisodeEntry{TransitionFn: hashref.SumSHA256([]byte("missing"))}
	ep := Episode{Entries: []EpisodeEntry{entry}}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err = e.ReplayEntry(ctx, ep, 0, entry); !errors.Is(err, context.Canceled) {
		t.Fatalf("entry ignored cancellation: %v", err)
	}
	if _, err = e.ReplayEpisode(ctx, ep); !errors.Is(err, context.Canceled) {
		t.Fatalf("episode ignored cancellation: %v", err)
	}
	_, err = e.ReplayEpisode(context.Background(), ep)
	var absent *archive.ReplayError
	if !errors.As(err, &absent) || absent.Kind != archive.KindAbsentArtifact {
		t.Fatalf("live control: %v", err)
	}
}

type replayWitness struct {
	*store.Store
	t               *testing.T
	want            context.Context
	cancel          context.CancelFunc
	objects, caches int
}

func (w *replayWitness) GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error) {
	if ctx != w.want {
		w.t.Fatal("source read lost caller context")
	}
	w.objects++
	o, ok, err := w.Store.GetObject(ctx, ref)
	w.cancel()
	return o, ok, err
}
func (w *replayWitness) GetVerifyResult(ctx context.Context, fn, interp hashref.HashRef) (store.VerifyResult, bool, error) {
	if ctx != w.want {
		w.t.Fatal("verify cache lost caller context")
	}
	w.caches++
	return w.Store.GetVerifyResult(ctx, fn, interp)
}
func TestReplayCacheCallerCancellation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fake-ailang")
	if err := os.WriteFile(path, []byte("#!/bin/sh\necho 'AILANG v0.30.0'\n"), 0700); err != nil {
		t.Fatal(err)
	}
	a := archive.New(filepath.Join(t.TempDir(), "world.db"))
	interp, err := a.Archive(path)
	if err != nil {
		t.Fatal(err)
	}
	s, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	payload := []byte("source")
	ref := hashref.SumSHA256(payload)
	if err := s.PutObject(store.Object{Hash: ref, InterfaceHash: hashref.SumSHA256([]byte("iface")), SemanticID: "test/source", Payload: payload}); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w := &replayWitness{Store: s, t: t, want: ctx, cancel: cancel}
	e := &Engine{store: w, archive: a}
	ep := Episode{Entries: []EpisodeEntry{{TransitionFn: ref, Interpreter: interp}}}
	if _, err := e.ReplayEpisode(ctx, ep); !errors.Is(err, context.Canceled) {
		t.Fatalf("cache ignored caller cancellation: %v", err)
	}
	if w.objects != 1 || w.caches != 1 {
		t.Fatalf("read census %d/%d", w.objects, w.caches)
	}
}
