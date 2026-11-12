package store

// Row 201 A4 sprint 2 M4f (design §2.7, §6 KD3, KD4, KD7 and KD9 in their SSH halves, K2):
// RecordKeyDecision's transaction and grammar, LandKeyDecision. Every store is under t.TempDir()
// (D-WORLD-109); the fixtures are key_decisions_test.go's.

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

// KD3 (and KD9's store half): an approver revoked after the read-only judge refuses inside the
// transaction on the first record and on an exact-match retry alike; nothing is written.
func TestRecordKeyDecisionRevokedSince(t *testing.T) {
	c := kdStore(t)
	if err := c.RecordKeyDecision(testCtx(t), kdRow(), kdUse, 1000); err != nil {
		t.Fatal(err)
	}
	first, _ := kdRead(t, c)
	if err := c.RevokeApprover(testCtx(t), "mark-ssh", 1001); err != nil {
		t.Fatal(err)
	}
	for _, row := range []KeyDecisionRow{kdRow(), func() KeyDecisionRow { r := kdRow(); r.Artifact = []byte("sig-2"); return r }()} {
		if err := c.RecordKeyDecision(testCtx(t), row, kdUse, 1002); !errors.Is(err, ErrApproverRevoked) {
			t.Errorf("artifact %s after the revoke: %v, want ErrApproverRevoked", row.Artifact, err)
		}
	}
	if got, _ := kdRead(t, c); !reflect.DeepEqual(got, first) {
		t.Fatalf("the row changed: %+v, was %+v", got, first)
	}
	other := kdUse
	other.PublicKey = []byte("k2")
	if err := c.RecordKeyDecision(testCtx(t), func() KeyDecisionRow { r := kdRow(); r.RequestRef = kdD; return r }(), other, 1003); !errors.Is(err, ErrApproverRevoked) {
		t.Fatalf("another key for the name: %v", err)
	}
}

// KD4, KD7 (SSH half): the exact artifact retried keeps the unlanded row as it was (decided_at
// of the first record); a row differing in the artifact or any signed field replaces it; once
// landed, no row replaces it (ErrKeyDecisionLanded) and it cannot land twice.
func TestRecordKeyDecisionReplaceUnlanded(t *testing.T) {
	c := kdStore(t)
	if err := c.RecordKeyDecision(testCtx(t), kdRow(), kdUse, 1000); err != nil {
		t.Fatal(err)
	}
	if err := c.RecordKeyDecision(testCtx(t), kdRow(), kdUse, 1010); err != nil {
		t.Fatal(err)
	}
	if got, ok := kdRead(t, c); !ok || got.DecidedAt != 1000 || got.DecisionRef != "" {
		t.Fatalf("exact retry: %+v %v, want the first record kept", got, ok)
	}
	at := int64(1020)
	for name, edit := range map[string]func(r *KeyDecisionRow){
		"artifact": func(r *KeyDecisionRow) { r.Artifact = []byte("sig-2") },
		"packet":   func(r *KeyDecisionRow) { r.PacketRef = kdD },
		"verb":     func(r *KeyDecisionRow) { r.Verb = "deny" },
		"expiry":   func(r *KeyDecisionRow) { r.Expiry++ },
	} {
		row := kdRow()
		edit(&row)
		if err := c.RecordKeyDecision(testCtx(t), row, kdUse, at); err != nil {
			t.Fatal(err)
		}
		row.DecidedAt = at
		if got, _ := kdRead(t, c); !reflect.DeepEqual(got, row) {
			t.Errorf("%s: %+v, want the row replaced by %+v", name, got, row)
		}
		if err := c.RecordKeyDecision(testCtx(t), kdRow(), kdUse, at+1); err != nil { // back to the first, for the next arm
			t.Fatal(err)
		}
		at += 10
	}
	if err := c.LandKeyDecision(testCtx(t), kdQ, kdD); err != nil {
		t.Fatal(err)
	}
	if got, _ := kdRead(t, c); got.DecisionRef != kdD {
		t.Fatalf("landed: %+v", got)
	}
	for _, row := range []KeyDecisionRow{kdRow(), func() KeyDecisionRow { r := kdRow(); r.Artifact = []byte("sig-3"); return r }()} {
		if err := c.RecordKeyDecision(testCtx(t), row, kdUse, 2000); !errors.Is(err, ErrKeyDecisionLanded) {
			t.Errorf("after landing, %s: %v, want ErrKeyDecisionLanded", row.Artifact, err)
		}
	}
	if err := c.LandKeyDecision(testCtx(t), kdQ, kdD); !errors.Is(err, ErrKeyDecisionNotPending) {
		t.Fatalf("second land: %v", err)
	}
}

// K2: the grammar refuses a bad ref, an upper-case ref, a verb, an approver that is not the use's,
// a format that is not the use's, and an empty or oversized artifact; nothing is written.
func TestRecordKeyDecisionGrammar(t *testing.T) {
	c := kdStore(t)
	for name, edit := range map[string]func(r *KeyDecisionRow){
		"request ref":    func(r *KeyDecisionRow) { r.RequestRef = kdQ[:70] },
		"upper-case ref": func(r *KeyDecisionRow) { r.PacketRef = "sha256:" + strings.ToUpper(kdP[7:]) },
		"verb":           func(r *KeyDecisionRow) { r.Verb = "maybe" },
		"approver":       func(r *KeyDecisionRow) { r.Approver = "ssh:mallory" },
		"format":         func(r *KeyDecisionRow) { r.Format = "webauthn" },
		"empty artifact": func(r *KeyDecisionRow) { r.Artifact = nil },
		"large artifact": func(r *KeyDecisionRow) { r.Artifact = make([]byte, 16385) },
	} {
		row := kdRow()
		edit(&row)
		if err := c.RecordKeyDecision(testCtx(t), row, kdUse, 1000); !errors.Is(err, ErrBadKeyDecision) {
			t.Errorf("%s: %v, want ErrBadKeyDecision", name, err)
		}
	}
	if err := c.LandKeyDecision(testCtx(t), kdQ, "sha256:x"); !errors.Is(err, ErrBadKeyDecision) {
		t.Errorf("land a bad decision ref: %v", err)
	}
	if _, ok := kdRead(t, c); ok {
		t.Fatal("a refused row was written")
	}
}
