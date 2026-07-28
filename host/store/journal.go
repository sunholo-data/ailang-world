package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/canon"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

const (
	JournalIntentV1  = "world/journal-intent/v1"
	JournalOutcomeV1 = "world/journal-outcome/v1"

	// MaxPendingIntentsPage is the kernel-owned upper bound for one pending scan.
	MaxPendingIntentsPage = 1000
)

// JournalIntent is the canonical statement of a planned commit. LogicalTime is
// supplied by the caller; journal payloads never read a wall clock.
type JournalIntent struct {
	InvocationID  string          `json:"invocationId"`
	WorldRef      hashref.HashRef `json:"-"`
	EntryHash     hashref.HashRef `json:"-"`
	ObservedHead  hashref.HashRef `json:"-"`
	PrevEntryHash hashref.HashRef `json:"-"`
	TransitionFn  hashref.HashRef `json:"-"`
	TransitionRef hashref.HashRef `json:"-"`
	Interpreter   hashref.HashRef `json:"-"`
	LogicalTime   int64           `json:"logicalTime"`
}

// JournalOutcome is the canonical statement of a known invocation result.
type JournalOutcome struct {
	InvocationID string          `json:"invocationId"`
	Status       string          `json:"status"`
	ResultRef    hashref.HashRef `json:"-"`
	LogicalTime  int64           `json:"logicalTime"`
}

// ReceiptState is the frozen receiptState projection.
type ReceiptState string

const (
	ReceiptNotStarted    ReceiptState = "not-started"
	ReceiptIndeterminate ReceiptState = "indeterminate"
	ReceiptResolved      ReceiptState = "resolved"
)

// Receipt reports the durable state and, when present, its payload references.
type Receipt struct {
	InvocationID string
	State        ReceiptState
	IntentSeq    int64
	IntentRef    hashref.HashRef
	OutcomeSeq   int64
	OutcomeRef   hashref.HashRef
	Intent       *JournalIntent
	Outcome      *JournalOutcome
}

// PendingIntent is one intent lacking an outcome. Seq is its keyset cursor.
type PendingIntent struct {
	Seq          int64
	InvocationID string
	ObjectRef    hashref.HashRef
	Intent       JournalIntent
}

// DuplicateInvocationError reports reuse of an invocation ID with different
// canonical bytes, or an attempt to append a second outcome.
type DuplicateInvocationError struct {
	ID   string
	Kind string
}

func (e *DuplicateInvocationError) Error() string {
	return fmt.Sprintf("store: duplicate %s for invocation %q", e.Kind, e.ID)
}

func IsDuplicateInvocation(err error) bool {
	var duplicate *DuplicateInvocationError
	return errors.As(err, &duplicate)
}

// InvocationMismatchError reports the first field that differs from the
// durable intent. Want and Got are canonical text.
type InvocationMismatchError struct {
	ID, Field, Want, Got string
}

func (e *InvocationMismatchError) Error() string {
	return fmt.Sprintf("store: invocation %q field %s mismatch: want %q, got %q",
		e.ID, e.Field, e.Want, e.Got)
}

func IsInvocationMismatch(err error) bool {
	var mismatch *InvocationMismatchError
	return errors.As(err, &mismatch)
}

type intentWire struct {
	InvocationID  string `json:"invocationId"`
	WorldRef      string `json:"worldRef"`
	EntryHash     string `json:"entryHash"`
	ObservedHead  string `json:"observedHead"`
	PrevEntryHash string `json:"prevEntryHash"`
	TransitionFn  string `json:"transitionFn"`
	TransitionRef string `json:"transitionRef"`
	Interpreter   string `json:"interpreter"`
	LogicalTime   int64  `json:"logicalTime"`
}

type outcomeWire struct {
	InvocationID string `json:"invocationId"`
	Status       string `json:"status"`
	ResultRef    string `json:"resultRef"`
	LogicalTime  int64  `json:"logicalTime"`
}

func canonicalJSON(v any) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return canon.Source(raw)
}

func encodeJournalIntent(v JournalIntent) ([]byte, error) {
	return canonicalJSON(intentWire{
		v.InvocationID, v.WorldRef.String(), v.EntryHash.String(), v.ObservedHead.String(),
		v.PrevEntryHash.String(), v.TransitionFn.String(), v.TransitionRef.String(),
		v.Interpreter.String(), v.LogicalTime,
	})
}

func encodeJournalOutcome(v JournalOutcome) ([]byte, error) {
	return canonicalJSON(outcomeWire{
		v.InvocationID, v.Status, v.ResultRef.String(), v.LogicalTime,
	})
}

func decodeJournalIntent(payload []byte) (JournalIntent, error) {
	var w intentWire
	if err := json.Unmarshal(payload, &w); err != nil {
		return JournalIntent{}, err
	}
	parse := func(field, text string) (hashref.HashRef, error) {
		if text == "" && field == "ObservedHead" {
			return hashref.HashRef{}, nil
		}
		ref, err := hashref.Parse(text)
		if err != nil {
			return hashref.HashRef{}, fmt.Errorf("%s: %w", field, err)
		}
		return ref, nil
	}
	world, err := parse("WorldRef", w.WorldRef)
	if err != nil {
		return JournalIntent{}, err
	}
	entry, err := parse("EntryHash", w.EntryHash)
	if err != nil {
		return JournalIntent{}, err
	}
	observed, err := parse("ObservedHead", w.ObservedHead)
	if err != nil {
		return JournalIntent{}, err
	}
	prev, err := parse("PrevEntryHash", w.PrevEntryHash)
	if err != nil {
		return JournalIntent{}, err
	}
	fn, err := parse("TransitionFn", w.TransitionFn)
	if err != nil {
		return JournalIntent{}, err
	}
	transition, err := parse("TransitionRef", w.TransitionRef)
	if err != nil {
		return JournalIntent{}, err
	}
	interpreter, err := parse("Interpreter", w.Interpreter)
	if err != nil {
		return JournalIntent{}, err
	}
	return JournalIntent{w.InvocationID, world, entry, observed, prev, fn, transition, interpreter, w.LogicalTime}, nil
}

func decodeJournalOutcome(payload []byte) (JournalOutcome, error) {
	var w outcomeWire
	if err := json.Unmarshal(payload, &w); err != nil {
		return JournalOutcome{}, err
	}
	ref, err := hashref.Parse(w.ResultRef)
	if err != nil {
		return JournalOutcome{}, fmt.Errorf("ResultRef: %w", err)
	}
	return JournalOutcome{w.InvocationID, w.Status, ref, w.LogicalTime}, nil
}

func journalObject(semanticID string, payload []byte) Object {
	return Object{
		Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte(semanticID)),
		SemanticID: semanticID, Provenance: "store/journal", Payload: payload,
	}
}

func validateIntent(id string, v JournalIntent) error {
	if id == "" {
		return &InvocationMismatchError{ID: id, Field: "InvocationID", Want: "non-empty", Got: id}
	}
	if v.InvocationID != id {
		return &InvocationMismatchError{ID: id, Field: "InvocationID", Want: id, Got: v.InvocationID}
	}
	refs := []struct {
		field string
		ref   hashref.HashRef
	}{
		{"WorldRef", v.WorldRef}, {"EntryHash", v.EntryHash},
		{"PrevEntryHash", v.PrevEntryHash}, {"TransitionFn", v.TransitionFn},
		{"TransitionRef", v.TransitionRef}, {"Interpreter", v.Interpreter},
	}
	for _, item := range refs {
		if err := validateRef("AppendIntent", item.field, item.ref); err != nil {
			return err
		}
	}
	if !v.ObservedHead.IsZero() {
		if err := validateRef("AppendIntent", "ObservedHead", v.ObservedHead); err != nil {
			return err
		}
	}
	return nil
}

// nextJournalSeqTx assigns max(seq)+1 inside the transaction. This is gapless
// because the ratified cross-process single-writer lock makes the writer unique.
func nextJournalSeqTx(tx *sql.Tx) (int64, error) {
	var seq int64
	if err := tx.QueryRow(`SELECT COALESCE(MAX(seq), 0) + 1 FROM journal`).Scan(&seq); err != nil {
		return 0, fmt.Errorf("store: next journal seq: %w", err)
	}
	return seq, nil
}

func insertJournalObjectTx(tx *sql.Tx, o Object) error {
	if _, err := tx.Exec(`INSERT OR IGNORE INTO objects
		(hash_ref, interface_hash_ref, semantic_id, provenance, payload)
		VALUES (?, ?, ?, ?, ?)`, o.Hash.String(), o.InterfaceHash.String(),
		o.SemanticID, o.Provenance, o.Payload); err != nil {
		return fmt.Errorf("store: append journal object: %w", err)
	}
	return nil
}

// AppendIntent durably appends a canonical intent and its index atomically.
func (s *Store) AppendIntent(id string, intent JournalIntent) (int64, hashref.HashRef, error) {
	if err := validateIntent(id, intent); err != nil {
		return 0, hashref.HashRef{}, err
	}
	payload, err := encodeJournalIntent(intent)
	if err != nil {
		return 0, hashref.HashRef{}, fmt.Errorf("store: encode intent: %w", err)
	}
	object := journalObject(JournalIntentV1, payload)
	tx, err := s.db.Begin()
	if err != nil {
		return 0, hashref.HashRef{}, fmt.Errorf("store: begin append intent: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var seq int64
	var existingText string
	err = tx.QueryRow(`SELECT seq, object_ref FROM journal
		WHERE invocation_id = ? AND kind = 'intent'`, id).Scan(&seq, &existingText)
	if err == nil {
		if existingText == object.Hash.String() {
			return seq, object.Hash, nil
		}
		return 0, hashref.HashRef{}, &DuplicateInvocationError{ID: id, Kind: "intent"}
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, hashref.HashRef{}, fmt.Errorf("store: find intent %q: %w", id, err)
	}
	seq, err = nextJournalSeqTx(tx)
	if err != nil {
		return 0, hashref.HashRef{}, err
	}
	if err := insertJournalObjectTx(tx, object); err != nil {
		return 0, hashref.HashRef{}, err
	}
	if _, err := tx.Exec(`INSERT INTO journal(seq, kind, invocation_id, object_ref)
		VALUES (?, 'intent', ?, ?)`, seq, id, object.Hash.String()); err != nil {
		return 0, hashref.HashRef{}, fmt.Errorf("store: append intent %q: %w", id, err)
	}
	if err := tx.Commit(); err != nil {
		return 0, hashref.HashRef{}, fmt.Errorf("store: commit intent: %w", err)
	}
	return seq, object.Hash, nil
}

// AppendOutcome requires a durable intent and atomically appends one outcome.
func (s *Store) AppendOutcome(id string, outcome JournalOutcome) (int64, hashref.HashRef, error) {
	if id == "" || outcome.InvocationID != id {
		return 0, hashref.HashRef{}, &InvocationMismatchError{
			ID: id, Field: "InvocationID", Want: id, Got: outcome.InvocationID,
		}
	}
	if err := validateRef("AppendOutcome", "ResultRef", outcome.ResultRef); err != nil {
		return 0, hashref.HashRef{}, err
	}
	payload, err := encodeJournalOutcome(outcome)
	if err != nil {
		return 0, hashref.HashRef{}, err
	}
	object := journalObject(JournalOutcomeV1, payload)
	tx, err := s.db.Begin()
	if err != nil {
		return 0, hashref.HashRef{}, err
	}
	defer func() { _ = tx.Rollback() }()
	var exists int
	if err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM journal
		WHERE invocation_id = ? AND kind = 'intent')`, id).Scan(&exists); err != nil {
		return 0, hashref.HashRef{}, fmt.Errorf("store: find outcome intent: %w", err)
	}
	if exists == 0 {
		return 0, hashref.HashRef{}, &InvocationMismatchError{
			ID: id, Field: "Intent", Want: "durable", Got: "missing",
		}
	}
	var existing int
	if err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM journal
		WHERE invocation_id = ? AND kind = 'outcome')`, id).Scan(&existing); err != nil {
		return 0, hashref.HashRef{}, err
	}
	if existing != 0 {
		return 0, hashref.HashRef{}, &DuplicateInvocationError{ID: id, Kind: "outcome"}
	}
	seq, err := nextJournalSeqTx(tx)
	if err != nil {
		return 0, hashref.HashRef{}, err
	}
	if err := insertJournalObjectTx(tx, object); err != nil {
		return 0, hashref.HashRef{}, err
	}
	if _, err := tx.Exec(`INSERT INTO journal(seq, kind, invocation_id, object_ref)
		VALUES (?, 'outcome', ?, ?)`, seq, id, object.Hash.String()); err != nil {
		return 0, hashref.HashRef{}, fmt.Errorf("store: append outcome: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return 0, hashref.HashRef{}, err
	}
	return seq, object.Hash, nil
}

type journalRow struct {
	seq int64
	ref hashref.HashRef
	obj Object
}

func journalRowFor(q interface{ QueryRow(string, ...any) *sql.Row }, id, kind string) (journalRow, bool, error) {
	var seq int64
	var refText string
	var ifaceText, semantic, provenance string
	var payload []byte
	err := q.QueryRow(`SELECT j.seq, j.object_ref, o.interface_hash_ref,
		o.semantic_id, o.provenance, o.payload FROM journal j
		JOIN objects o ON o.hash_ref = j.object_ref
		WHERE j.invocation_id = ? AND j.kind = ?`, id, kind).
		Scan(&seq, &refText, &ifaceText, &semantic, &provenance, &payload)
	if errors.Is(err, sql.ErrNoRows) {
		return journalRow{}, false, nil
	}
	if err != nil {
		return journalRow{}, false, fmt.Errorf("store: load journal %s: %w", kind, err)
	}
	ref, err := hashref.Parse(refText)
	if err != nil {
		return journalRow{}, false, err
	}
	iface, err := hashref.Parse(ifaceText)
	if err != nil {
		return journalRow{}, false, err
	}
	return journalRow{seq, ref, Object{ref, iface, semantic, provenance, payload}}, true, nil
}

// GetReceipt mirrors receiptState and never reports not-started with an intent.
func (s *Store) GetReceipt(id string) (Receipt, bool, error) {
	intentRow, hasIntent, err := journalRowFor(s.db, id, "intent")
	if err != nil {
		return Receipt{}, false, err
	}
	outcomeRow, hasOutcome, err := journalRowFor(s.db, id, "outcome")
	if err != nil {
		return Receipt{}, false, err
	}
	if !hasIntent && hasOutcome {
		return Receipt{}, false, &InvocationMismatchError{ID: id, Field: "Intent", Want: "durable", Got: "missing"}
	}
	if !hasIntent {
		return Receipt{InvocationID: id, State: ReceiptNotStarted}, false, nil
	}
	intent, err := decodeJournalIntent(intentRow.obj.Payload)
	if err != nil {
		return Receipt{}, false, fmt.Errorf("store: decode intent: %w", err)
	}
	receipt := Receipt{InvocationID: id, State: ReceiptIndeterminate, IntentSeq: intentRow.seq,
		IntentRef: intentRow.ref, Intent: &intent}
	if !hasOutcome {
		return receipt, true, nil
	}
	outcome, err := decodeJournalOutcome(outcomeRow.obj.Payload)
	if err != nil {
		return Receipt{}, false, fmt.Errorf("store: decode outcome: %w", err)
	}
	receipt.State, receipt.OutcomeSeq, receipt.OutcomeRef, receipt.Outcome =
		ReceiptResolved, outcomeRow.seq, outcomeRow.ref, &outcome
	return receipt, true, nil
}

// PendingIntents returns oldest pending intents using seq keyset pagination.
// The optional fromIndex cursor is exclusive; omitting it starts before seq 1.
func (s *Store) PendingIntents(limit int, fromIndex ...int64) ([]PendingIntent, error) {
	if limit < 1 || limit > MaxPendingIntentsPage {
		return nil, &InvalidLimitError{Op: "PendingIntents", Limit: limit, Max: MaxPendingIntentsPage}
	}
	var after int64
	if len(fromIndex) > 1 {
		return nil, fmt.Errorf("store: PendingIntents: at most one fromIndex cursor")
	}
	if len(fromIndex) == 1 {
		after = fromIndex[0]
	}
	rows, err := s.db.Query(`SELECT j.seq, j.invocation_id, j.object_ref, o.payload
		FROM journal j JOIN objects o ON o.hash_ref = j.object_ref
		WHERE j.kind = 'intent' AND j.seq > ?
		AND NOT EXISTS (SELECT 1 FROM journal x WHERE x.invocation_id = j.invocation_id
			AND x.kind = 'outcome')
		ORDER BY j.seq LIMIT ?`, after, limit)
	if err != nil {
		return nil, fmt.Errorf("store: pending intents: %w", err)
	}
	defer rows.Close()
	var pending []PendingIntent
	for rows.Next() {
		var item PendingIntent
		var refText string
		var payload []byte
		if err := rows.Scan(&item.Seq, &item.InvocationID, &refText, &payload); err != nil {
			return nil, err
		}
		item.ObjectRef, err = hashref.Parse(refText)
		if err != nil {
			return nil, err
		}
		item.Intent, err = decodeJournalIntent(payload)
		if err != nil {
			return nil, err
		}
		pending = append(pending, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return pending, nil
}
