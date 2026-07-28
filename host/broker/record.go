package broker

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

const (
	EffectRecordV1 = "world/effect-record/v1"
	EffectResultV1 = "world/effect-result/v1"
)

// EffectRecord is the immutable accounting and replay input for one request.
type EffectRecord struct {
	Effect       string
	Scope        string
	Cost         int64
	BudgetBefore int64
	BudgetAfter  int64
	Allowed      bool
	Denial       string
	RequestRef   hashref.HashRef
	ResultRef    hashref.HashRef
}

type recordWire struct {
	Effect       string `json:"effect"`
	Scope        string `json:"scope"`
	Cost         int64  `json:"cost"`
	BudgetBefore int64  `json:"budgetBefore"`
	BudgetAfter  int64  `json:"budgetAfter"`
	Allowed      bool   `json:"allowed"`
	Denial       string `json:"denial"`
	RequestRef   string `json:"requestRef"`
	ResultRef    string `json:"resultRef"`
}

// EncodeRecord is the single fixed-field-order record codec.
func EncodeRecord(rec EffectRecord) []byte {
	payload, err := json.Marshal(recordWire{
		rec.Effect, rec.Scope, rec.Cost, rec.BudgetBefore, rec.BudgetAfter,
		rec.Allowed, rec.Denial, rec.RequestRef.String(), rec.ResultRef.String(),
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
	return EffectRecord{
		wire.Effect, wire.Scope, wire.Cost, wire.BudgetBefore, wire.BudgetAfter,
		wire.Allowed, wire.Denial, requestRef, resultRef,
	}, nil
}

func parseRequiredRef(field, text string) (hashref.HashRef, error) {
	ref, err := hashref.Parse(text)
	if err != nil {
		return hashref.HashRef{}, fmt.Errorf("broker: decode effect record %s: %w", field, err)
	}
	return ref, nil
}

// RecordConsistent mirrors the sketch's recordConsistent law.
func RecordConsistent(rec EffectRecord) bool {
	return rec.Allowed && rec.BudgetAfter == rec.BudgetBefore-rec.Cost && rec.Denial == "" ||
		!rec.Allowed && rec.BudgetAfter == rec.BudgetBefore
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
