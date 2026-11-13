package store

// Row 201 A4 sprint 3 M6c (design §2.7 [r2], §6 KD2, KD7 passkey half, K2's passkey rule):
// RecordKeyDecision's counter (K4: usePasskey, the assertion law's W15 against the stored counter)
// off the exact match, and only the counter skipped on it. The fixtures are passkey_consume_test.go's
// (approver x, credential credX); every store is under t.TempDir() (D-WORLD-109).

import (
	"errors"
	"reflect"
	"testing"
)

// pkdUse is x's passkey use with counter count; pkdRow its decision row with artifact a.
func pkdUse(count int64) ApproverUse {
	return ApproverUse{Name: "x", Format: "webauthn", PublicKey: []byte("cose-x"), Passkey: &PasskeyUse{CredentialID: credX, SignCount: count}}
}

func pkdRow(a string) KeyDecisionRow {
	return KeyDecisionRow{RequestRef: kdQ, PacketRef: kdP, Verb: "approve", Format: "webauthn", Approver: "webauthn:x", Artifact: []byte(a), Expiry: 1760}
}

// KD2: a counter that does not advance past the stored one is ErrSignCountRegression with nothing
// written (no row, the counter and last_used_at unchanged); an advancing one is recorded and moves
// it; stored 0 and assertion 0 (an authenticator without a counter) is recorded, counter still 0.
func TestRecordKeyDecisionCounter(t *testing.T) {
	c, _ := passkeyAt(t, credX, 7)
	for _, n := range []int64{7, 6} {
		if err := c.RecordKeyDecision(testCtx(t), pkdRow("asr-1"), pkdUse(n), 1000); !errors.Is(err, ErrSignCountRegression) {
			t.Fatalf("counter %d over stored 7: %v", n, err)
		}
	}
	if count, used := passkeyState(t, c); count != 7 || used != 0 {
		t.Fatalf("refused: counter %d, last used %d", count, used)
	}
	if _, ok := kdRead(t, c); ok {
		t.Fatal("refused: a row was written")
	}
	if err := c.RecordKeyDecision(testCtx(t), pkdRow("asr-1"), pkdUse(8), 1000); err != nil {
		t.Fatal(err)
	}
	if count, used := passkeyState(t, c); count != 8 || used != 1000 {
		t.Fatalf("recorded: counter %d, last used %d", count, used)
	}
	z, _ := passkeyAt(t, credX, 0)
	if err := z.RecordKeyDecision(testCtx(t), pkdRow("asr-0"), pkdUse(0), 1000); err != nil {
		t.Fatalf("counter 0 over stored 0: %v", err)
	}
	if count, _ := passkeyState(t, z); count != 0 {
		t.Fatalf("counter 0: stored %d", count)
	}
}

// KD7 (passkey half): an unlanded row from a counter-7 assertion (stored now 7); the same artifact
// retried → the exact match: recorded without the counter (still 7), the row kept as first
// recorded; a second, independently made assertion that also carries counter 7 → ErrSignCountRegression,
// the row unchanged; one carrying 8 replaces it (K5) and moves the counter.
func TestRecordKeyDecisionIdempotentPasskeyRetry(t *testing.T) {
	c, _ := passkeyAt(t, credX, 6)
	if err := c.RecordKeyDecision(testCtx(t), pkdRow("asr-7"), pkdUse(7), 1000); err != nil {
		t.Fatal(err)
	}
	first, _ := kdRead(t, c)
	if err := c.RecordKeyDecision(testCtx(t), pkdRow("asr-7"), pkdUse(7), 1010); err != nil {
		t.Fatalf("the exact retry: %v", err)
	}
	if got, _ := kdRead(t, c); !reflect.DeepEqual(got, first) {
		t.Fatalf("the exact retry rewrote the row: %+v, was %+v", got, first)
	}
	if count, used := passkeyState(t, c); count != 7 || used != 1010 {
		t.Fatalf("the exact retry: counter %d, last used %d (useApprover runs, usePasskey does not)", count, used)
	}
	if err := c.RecordKeyDecision(testCtx(t), pkdRow("asr-7b"), pkdUse(7), 1020); !errors.Is(err, ErrSignCountRegression) {
		t.Fatalf("another assertion at 7: %v", err)
	}
	if got, _ := kdRead(t, c); !reflect.DeepEqual(got, first) {
		t.Fatalf("a refused assertion changed the row: %+v", got)
	}
	if err := c.RecordKeyDecision(testCtx(t), pkdRow("asr-8"), pkdUse(8), 1030); err != nil {
		t.Fatal(err)
	}
	if got, _ := kdRead(t, c); string(got.Artifact) != "asr-8" || got.DecidedAt != 1030 {
		t.Fatalf("a fresh assertion: %+v", got)
	}
	if count, _ := passkeyState(t, c); count != 8 {
		t.Fatalf("a fresh assertion: counter %d", count)
	}
}

// K2 (passkey rule): a webauthn decision without a passkey use, and an ssh one with one, are
// ErrBadKeyDecision; nothing is written and the counter does not move.
func TestRecordKeyDecisionPasskeyUseMatchesTheFormat(t *testing.T) {
	c, _ := passkeyAt(t, credX, 1)
	noPasskey := pkdUse(2)
	noPasskey.Passkey = nil
	if err := c.RecordKeyDecision(testCtx(t), pkdRow("asr"), noPasskey, 1000); !errors.Is(err, ErrBadKeyDecision) {
		t.Errorf("webauthn without a passkey use: %v", err)
	}
	if err := c.EnrolApprover(testCtx(t), ApproverRow{Name: "mark-ssh", Format: "ssh", PublicKey: []byte("k1"), EnrolledAt: 900}); err != nil {
		t.Fatal(err)
	}
	withPasskey := kdUse
	withPasskey.Passkey = &PasskeyUse{CredentialID: credX, SignCount: 2}
	if err := c.RecordKeyDecision(testCtx(t), kdRow(), withPasskey, 1000); !errors.Is(err, ErrBadKeyDecision) {
		t.Errorf("ssh with a passkey use: %v", err)
	}
	if _, ok := kdRead(t, c); ok {
		t.Fatal("a refused row was written")
	}
	if count, _ := passkeyState(t, c); count != 1 {
		t.Fatalf("the counter moved: %d", count)
	}
}
