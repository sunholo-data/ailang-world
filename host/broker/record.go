package broker

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

const (
	EffectRequestV1 = "world/effect-request/v1"
	EffectRecordV1  = "world/effect-record/v1"
	EffectResultV1  = "world/effect-result/v1"
)

// EffectRecord is the immutable accounting and replay input for one request.
type EffectRecord struct {
	Effect       string
	Scope        string
	Cost         int64
	BudgetBefore int64
	BudgetAfter  int64
	Allowed      bool
	Failed       bool
	Denial       string
	RequestRef   hashref.HashRef
	ResultRef    hashref.HashRef
	// Failure is row 255's bounded failure line (FailureLine), set only on a
	// failed record whose handler returned a *HandlerFailure. It is omitted
	// from the wire when empty, so every other record keeps its bytes.
	Failure string
}

type recordWire struct {
	Effect       string `json:"effect"`
	Scope        string `json:"scope"`
	Cost         int64  `json:"cost"`
	BudgetBefore int64  `json:"budgetBefore"`
	BudgetAfter  int64  `json:"budgetAfter"`
	Allowed      bool   `json:"allowed"`
	Failed       bool   `json:"failed"`
	Denial       string `json:"denial"`
	RequestRef   string `json:"requestRef"`
	ResultRef    string `json:"resultRef"`
	Failure      string `json:"failure,omitempty"`
}

// EncodeRecord is the single fixed-field-order record codec.
func EncodeRecord(rec EffectRecord) []byte {
	payload, err := json.Marshal(recordWire{
		rec.Effect, rec.Scope, rec.Cost, rec.BudgetBefore, rec.BudgetAfter,
		rec.Allowed, rec.Failed, rec.Denial, rec.RequestRef.String(), rec.ResultRef.String(), rec.Failure,
	})
	if err != nil {
		panic("broker: fixed effect record cannot fail JSON encoding: " + err.Error())
	}
	return payload
}

// DecodeRecord decodes and validates the reference syntax in a record.
func DecodeRecord(payload []byte) (EffectRecord, error) {
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.DisallowUnknownFields()
	var wire recordWire
	if err := dec.Decode(&wire); err != nil {
		return EffectRecord{}, fmt.Errorf("broker: decode effect record: %w", err)
	}
	requestRef, err := parseRequiredRef("requestRef", wire.RequestRef)
	if err != nil {
		return EffectRecord{}, err
	}
	var resultRef hashref.HashRef
	if wire.ResultRef != "" {
		resultRef, err = parseRequiredRef("resultRef", wire.ResultRef)
		if err != nil {
			return EffectRecord{}, err
		}
	}
	rec := EffectRecord{
		wire.Effect, wire.Scope, wire.Cost, wire.BudgetBefore, wire.BudgetAfter,
		wire.Allowed, wire.Failed, wire.Denial, requestRef, resultRef, wire.Failure,
	}
	if !FailureConsistent(rec) {
		return EffectRecord{}, fmt.Errorf("broker: decode effect record failure: %q is not a bounded failure line of a failed record", wire.Failure)
	}
	return rec, nil
}

func parseRequiredRef(field, text string) (hashref.HashRef, error) {
	ref, err := hashref.Parse(text)
	if err != nil {
		return hashref.HashRef{}, fmt.Errorf("broker: decode effect record %s: %w", field, err)
	}
	return ref, nil
}

// RecordConsistent mirrors design_docs/sketches/effectbroker.ail. The arm table
// is (allowed, failed): denied=(false,false), success=(true,false), and
// failed=(true,true); (false,true) is illegal. The sketch's
// rec.resultRef.digest == "" is equivalent to ResultRef.IsZero() and to
// "resultRef":"" on the wire: DecodeRecord makes a zero wire value the full Go
// zero ref, while non-empty wire refs must parse and therefore cannot be half-zero.
func RecordConsistent(rec EffectRecord) bool {
	switch {
	case rec.Allowed && !rec.Failed:
		return rec.Denial == "" && rec.BudgetAfter == rec.BudgetBefore-rec.Cost && !rec.ResultRef.IsZero()
	case rec.Allowed && rec.Failed:
		return rec.Denial == "" && rec.BudgetAfter == rec.BudgetBefore-rec.Cost && rec.ResultRef.IsZero()
	default: // !Allowed
		return !rec.Failed && rec.Denial != "" && rec.BudgetAfter == rec.BudgetBefore && rec.ResultRef.IsZero()
	}
}

// FailureConsistent is row 255's rule for the failure line: only an allowed,
// failed record carries one, and it is a fixed point of FailureLine. It sits
// beside RecordConsistent, not inside it: RecordConsistent mirrors the
// accounting law of sketches/effectbroker.ail, which the line does not touch.
func FailureConsistent(rec EffectRecord) bool {
	return rec.Failure == "" || rec.Allowed && rec.Failed && ValidFailureLine(rec.Failure)
}

func recordObject(rec EffectRecord) store.Object {
	return brokerObject(EffectRecordV1, EncodeRecord(rec))
}

func resultObject(payload []byte) store.Object {
	return brokerObject(EffectResultV1, payload)
}

func brokerObject(semanticID string, payload []byte) store.Object {
	return store.Object{
		Hash:          hashref.SumSHA256(payload),
		InterfaceHash: hashref.SumSHA256([]byte(semanticID)),
		SemanticID:    semanticID,
		Provenance:    "host/broker",
		Payload:       payload,
	}
}
