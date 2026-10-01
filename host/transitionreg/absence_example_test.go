package transitionreg

import (
	"context"
	"errors"
	"fmt"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"testing"
	"time"
)

type absenceStore struct {
	*fakeObjectStore
	after func()
}

func (s absenceStore) GetRegistryHead(ctx context.Context, name string) (hashref.HashRef, bool, error) {
	h, ok, err := s.fakeObjectStore.GetRegistryHead(ctx, name)
	if s.after != nil {
		s.after()
	}
	return h, ok, err
}

type absenceCaps struct{}

func (absenceCaps) CapabilitySnapshot(int64) broker.CapabilitySnapshot {
	return broker.CapabilitySnapshot{}
}
func TestRegistryHeadAbsentErrorClassification(t *testing.T) {
	for _, name := range []string{"direct", "new_request"} {
		t.Run(name, func(t *testing.T) {
			s := &fakeObjectStore{}
			r := NewReader(s)
			var err error
			if name == "direct" {
				_, err = r.ReadSnapshot(boundedTestContext(t))
			} else {
				_, err = NewRequest(boundedTestContext(t), r, absenceCaps{}, 1)
			}
			var target RegistryHeadAbsentError
			if !errors.As(err, &target) {
				t.Fatalf("typed absence=false: %v", err)
			}
			if target.Error() != "read transition registry: head is absent" || s.headReads != 1 || s.objectReads != 0 {
				t.Fatalf("message/counters: %v %d/%d", err, s.headReads, s.objectReads)
			}
		})
	}
	for _, tc := range []struct {
		name string
		s    *fakeObjectStore
	}{
		{"store", &fakeObjectStore{headErr: errors.New("store failure")}},
		{"same_text", &fakeObjectStore{headErr: errors.New("read transition registry: head is absent")}},
		{"object", &fakeObjectStore{hasHead: true, head: testFn, objectErr: errors.New("object failure")}},
		{"missing_object", &fakeObjectStore{hasHead: true, head: testFn}},
		{"hash", &fakeObjectStore{hasHead: true, head: testFn, hasObject: true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewReader(tc.s).ReadSnapshot(boundedTestContext(t))
			var target RegistryHeadAbsentError
			if err == nil || errors.As(err, &target) {
				t.Fatalf("genuine error classified as absence: %v", err)
			}
		})
	}
}
func TestRegistryHeadAbsentAfterLookupCancellation(t *testing.T) {
	for _, name := range []string{"absent", "present", "deadline"} {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(boundedTestContext(t))
			defer cancel()
			want := error(context.Canceled)
			after := cancel
			if name == "deadline" {
				var stop context.CancelFunc
				ctx, stop = context.WithTimeout(boundedTestContext(t), time.Millisecond)
				defer stop()
				want = context.DeadlineExceeded
				after = func() { <-ctx.Done() }
			}
			s := &fakeObjectStore{hasHead: name == "present", head: testFn}
			_, err := NewReader(absenceStore{s, after}).ReadSnapshot(ctx)
			var target RegistryHeadAbsentError
			t.Logf("context=%t typedabsence=%t error=%v", errors.Is(err, want), errors.As(err, &target), err)
			if !errors.Is(err, want) || errors.As(err, &target) || s.headReads != 1 || s.objectReads != 0 {
				t.Fatalf("postlookup context/counters: %v %d/%d", err, s.headReads, s.objectReads)
			}
		})
	}
}
func TestRegistryHeadAbsentIntegrityClassification(t *testing.T) {
	for _, name := range []string{"semantic", "interface", "codec"} {
		t.Run(name, func(t *testing.T) {
			payload := []byte("not a revision")
			obj := store.Object{Hash: hashref.SumSHA256(payload), Payload: payload, SemanticID: SemanticIDV1, InterfaceHash: InterfaceHashV1}
			if name == "semantic" {
				obj.SemanticID = "wrong"
			}
			if name == "interface" {
				obj.InterfaceHash = testFn
			}
			s := &fakeObjectStore{head: obj.Hash, hasHead: true, hasObject: true, object: obj}
			_, err := NewRequest(boundedTestContext(t), NewReader(s), absenceCaps{}, 1)
			var target RegistryHeadAbsentError
			if err == nil || errors.As(err, &target) || s.headReads != 1 || s.objectReads != 1 {
				t.Fatalf("integrity classified as absence: %v counts %d/%d", err, s.headReads, s.objectReads)
			}
		})
	}
}
func ExampleRegistryHeadAbsentError() {
	err := fmt.Errorf("request: %w", RegistryHeadAbsentError{})
	var target RegistryHeadAbsentError
	fmt.Println(errors.As(err, &target))
	fmt.Println(target.Error())
	// Output:
	// true
	// read transition registry: head is absent
}
