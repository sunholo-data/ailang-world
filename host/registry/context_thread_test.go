package registry

import (
	"context"
	"errors"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"testing"
)

func TestBootstrapCallerCancellation(t *testing.T) {
	s := openMem(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, _, err := Bootstrap(ctx, s, m1Release); !errors.Is(err, context.Canceled) {
		t.Fatalf("bootstrap ignored caller cancellation: %v", err)
	}
	if _, _, err := Bootstrap(context.Background(), s, m1Release); err != nil {
		t.Fatalf("live control: %v", err)
	}
}

type bootstrapWitness struct {
	*store.Store
	t              *testing.T
	want           context.Context
	cancel         context.CancelFunc
	objects, heads int
}

func (w *bootstrapWitness) GetRegistryHead(ctx context.Context, name string) (hashref.HashRef, bool, error) {
	if ctx != w.want {
		w.t.Fatal("registry head lost caller context")
	}
	w.heads++
	ref, ok, err := w.Store.GetRegistryHead(ctx, name)
	w.cancel()
	return ref, ok, err
}
func (w *bootstrapWitness) GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error) {
	if ctx != w.want {
		w.t.Fatal("registry object lost caller context")
	}
	w.objects++
	return w.Store.GetObject(ctx, ref)
}
func TestBootstrapSecondReadCancellation(t *testing.T) {
	s := openMem(t)
	if _, _, err := Bootstrap(context.Background(), s, m1Release); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w := &bootstrapWitness{Store: s, t: t, want: ctx, cancel: cancel}
	if _, _, err := bootstrap(ctx, w, m1Release); !errors.Is(err, context.Canceled) {
		t.Fatalf("object read ignored caller: %v", err)
	}
	if w.heads != 1 || w.objects != 1 {
		t.Fatalf("read census %d/%d", w.heads, w.objects)
	}
}
