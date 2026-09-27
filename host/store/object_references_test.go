package store

import (
	"context"
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

func TestObjectReferencesValidation(t *testing.T) {
	ref := hashref.SumSHA256([]byte("target"))
	world := hashref.SumSHA256([]byte("world"))
	t.Run("limits", func(t *testing.T) {
		for _, limit := range []int{-1, 0, 501, 1000} {
			err := validateObjectReferences(ref, nil, limit)
			var bad *InvalidLimitError
			if !errors.As(err, &bad) {
				t.Fatalf("limit %d: %v", limit, err)
			}
		}
		for _, limit := range []int{1, 500} {
			if err := validateObjectReferences(ref, nil, limit); err != nil {
				t.Fatal(err)
			}
		}
	})
	t.Run("ref", func(t *testing.T) {
		err := validateObjectReferences(hashref.HashRef{}, nil, 1)
		var bad *InvalidRefError
		if !errors.As(err, &bad) {
			t.Fatalf("ref error: %v", err)
		}
	})
	t.Run("cursor", func(t *testing.T) {
		good := []ObjectReferenceCursor{{Kind: ReferenceTransitionRef, EntryIndex: 0}, {Kind: ReferenceTransitionFn, EntryIndex: 2}, {Kind: ReferenceInterpreter, EntryIndex: 10}, {Kind: ReferenceStateRoot, WorldRef: world}}
		for _, c := range good {
			if err := validateObjectReferences(ref, &c, 1); err != nil {
				t.Fatalf("good %+v: %v", c, err)
			}
		}
		bad := []ObjectReferenceCursor{{Kind: -1}, {Kind: 4}, {Kind: ReferenceTransitionRef, EntryIndex: -1}, {Kind: ReferenceTransitionFn, WorldRef: world}, {Kind: ReferenceStateRoot}, {Kind: ReferenceStateRoot, EntryIndex: 1, WorldRef: world}}
		for _, c := range bad {
			err := validateObjectReferences(ref, &c, 1)
			var invalid *InvalidObjectReferenceCursorError
			if !errors.As(err, &invalid) {
				t.Fatalf("bad %+v: %v", c, err)
			}
		}
	})
	s := openMem(t)
	got, err := s.ObjectReferences(context.Background(), ref, nil, 1)
	if err != nil || got == nil || len(got) != 0 {
		t.Fatalf("empty page=%v %v", got, err)
	}
}
