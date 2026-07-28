package broker

import (
	"context"
	"errors"
	"testing"
)

// CF-I-2 — a COMMITTED REPRODUCTION of a known, unresolved gap, following the
// `host/store/durability_repro_test.go` precedent: it asserts the CURRENT
// behaviour on purpose, so the gap is CI-enforced rather than a note in a log,
// and so whoever closes it must deliberately rewrite this file.
//
// THE GAP. Decision 3 freezes a pipeline with exactly two arms — denied and
// allowed. There is no HANDLER-ERROR arm. When a handler fails, the ledger has
// already been debited (broker.go debits before dispatch) and NO record is
// written, so:
//
//   - Decision 3's "the ledger is reconstructible ... from the record stream
//     alone" is FALSE on this path: the ledger says 2, the records say 5.
//   - A failed effect is invisible to audit and to replay, while a merely
//     DENIED effect is fully recorded — the weaker outcome is better recorded
//     than the stronger one.
//   - The handler may have partially executed before failing (bytes written,
//     tokens spent, a git object created), so this is not merely an accounting
//     nit.
//
// This is NOT the M3.D crash window (process death, blocked on ratification).
// It is an ordinary in-process error path and M3.D's question does not cover it.
//
// WHY IT IS NOT FIXED HERE. The two candidate fixes have opposite semantics —
// (a) roll the debit back, which refunds an effect that may have spent real
// money; (b) keep the debit and write a failure record, which adds a third arm
// to a frozen Decision. Choosing is a design call, not a sprint fix, so it is
// recorded rather than force-applied.
type cfi2FailingHandler struct{ calls int }

func (h *cfi2FailingHandler) Execute(ctx context.Context, req EffectRequest, payload []byte) ([]byte, error) {
	h.calls++
	return nil, errors.New("handler failed after partial execution")
}

func TestCFI2HandlerErrorDebitsLedgerAndWritesNoRecord(t *testing.T) {
	st := openTestStore(t)
	handler := &cfi2FailingHandler{}
	s := NewSession(st,
		[]Capability{{Effect: "FS.Write", Scope: "/p", ExpiresAt: 100, Budget: 5}},
		Registry{"FS.Write": handler})

	before := s.grants[0].Budget
	_, ref, err := s.Invoke(context.Background(),
		EffectRequest{Effect: "FS.Write", Scope: "/p", Cost: 3, Now: 1}, []byte("x"))

	if err == nil {
		t.Fatal("want a handler error, got nil")
	}
	if handler.calls != 1 {
		t.Fatalf("handler dispatched %d times, want exactly 1", handler.calls)
	}

	// CURRENT behaviour, asserted so a change is visible. When CF-I-2 is
	// resolved, EVERY assertion below is expected to fail and be rewritten.
	if got, want := s.grants[0].Budget, before-3; got != want {
		t.Fatalf("ledger budget = %d, want %d (the debit currently survives the error)", got, want)
	}
	if !ref.IsZero() {
		t.Fatalf("record ref = %s, want zero (no record is currently written on handler error)", ref.String())
	}
}
