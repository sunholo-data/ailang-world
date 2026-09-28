package store

import (
	"context"
	"errors"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"testing"
)

func TestVerifyResultCallerCancellation(t *testing.T) {
	s, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	ref := hashref.SumSHA256([]byte("verify"))
	ctx, cancel := context.WithCancel(boundedTestContext(t))
	cancel()
	if _, _, err := s.GetVerifyResult(ctx, ref, ref); !errors.Is(err, context.Canceled) {
		t.Fatalf("cache ignored cancellation: %v", err)
	}
	if _, ok, err := s.GetVerifyResult(boundedTestContext(t), ref, ref); ok || err != nil {
		t.Fatalf("live cache miss control: %v %v", ok, err)
	}
}
