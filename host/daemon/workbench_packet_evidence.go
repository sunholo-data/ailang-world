package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

const (
	packetRosterInputLimit = 64 << 10
	packetRosterReadLimit  = 1 << 20
	packetRosterClaimLimit = 4 << 10
	// TESTED consumes these in M5b; no test stdout is decoded in M5a.
	packetRosterStdoutLimit    = 1 << 20
	packetRosterNamedTestLimit = 256
	packetRosterReadTimeout    = 10 * time.Second
)

// These transport projections carry no object-grade authority. The package's
// evidence.ail owns kind admission and constructor grades (executable drift gate).
// M5b will adapt these internal fields to the workbench renderer.
type packetRoster struct {
	rows        []packetEvidenceRow
	unavailable string
}
type packetEvidenceRow struct {
	kind, ref, label, producer, unavailable string
}
type packetRosterReader interface {
	GetObject(context.Context, hashref.HashRef) (store.Object, bool, error)
	GetLogEntry(context.Context, int64) (store.LogEntry, bool, error)
}
type packetRosterPair struct {
	Kind string `json:"kind"`
	Ref  string `json:"ref"`
}
type packetRosterInput struct {
	Goal       string             `json:"goal"`
	Candidates []string           `json:"candidates"`
	Roster     []packetRosterPair `json:"roster"`
}

// strictRosterJSON rejects duplicate keys as well as unknown keys and trailing
// values. encoding/json's default last-key-wins behavior is unsuitable here.
func strictRosterJSON(raw []byte, dst any) error {
	dec := json.NewDecoder(bytes.NewReader(raw))
	var value func() error
	value = func() error {
		token, err := dec.Token()
		if err != nil {
			return err
		}
		delim, ok := token.(json.Delim)
		if !ok {
			return nil
		}
		switch delim {
		case '{':
			seen := map[string]bool{}
			for dec.More() {
				key, err := dec.Token()
				if err != nil {
					return err
				}
				s, ok := key.(string)
				if !ok || seen[s] {
					return errors.New("duplicate JSON key")
				}
				// encoding/json accepts case-insensitive field aliases. The
				// producer's strict shape admits only the exact JSON keys.
				for _, field := range []string{"goal", "candidates", "roster", "kind", "ref", "author", "text"} {
					if s != field && strings.EqualFold(s, field) {
						return errors.New("noncanonical JSON key")
					}
				}
				seen[s] = true
				if err = value(); err != nil {
					return err
				}
			}
		case '[':
			for dec.More() {
				if err = value(); err != nil {
					return err
				}
			}
		default:
			return errors.New("unexpected JSON delimiter")
		}
		_, err = dec.Token()
		return err
	}
	if err := value(); err != nil {
		return err
	}
	if _, err := dec.Token(); err != io.EOF {
		return errors.New("trailing JSON")
	}
	dec = json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

// packetRosterFromValidated follows only the carrying entry of an already
// validated packet. No semantic-ID lookup can discover or substitute evidence.
// All refusals discard the entire roster before any label can escape.
func packetRosterFromValidated(parent context.Context, r packetRosterReader, packet coordinator.PacketResult) packetRoster {
	refuse := func(reason string) packetRoster { return packetRoster{unavailable: reason} }
	if packet.Code != 0 {
		return refuse("packet not validated")
	}
	ctx, cancel := context.WithTimeout(parent, packetRosterReadTimeout)
	defer cancel()
	read := func(ref hashref.HashRef, role, id string, limit int) (store.Object, string) {
		if ctx.Err() != nil {
			return store.Object{}, "read deadline"
		}
		o, ok, err := r.GetObject(ctx, ref)
		if ctx.Err() != nil || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return store.Object{}, "read deadline"
		}
		if err != nil || !ok {
			return store.Object{}, role + " missing"
		}
		if len(o.Payload) > limit {
			return store.Object{}, role + " size"
		}
		if o.Hash != ref || hashref.SumSHA256(o.Payload) != ref {
			return store.Object{}, role + " hash"
		}
		if o.SemanticID != id || o.InterfaceHash != hashref.SumSHA256([]byte(id)) {
			return store.Object{}, role + " kind"
		}
		return o, ""
	}
	entry, ok, err := r.GetLogEntry(ctx, packet.Entry)
	if ctx.Err() != nil || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return refuse("read deadline")
	}
	if err != nil || !ok {
		return refuse("carrying entry missing")
	}
	record, reason := read(entry.TransitionRef, "record", coordinator.RecordV2, packetRosterReadLimit)
	if reason != "" {
		return refuse(reason)
	}
	var rec struct {
		Input string `json:"input"`
	}
	if json.Unmarshal(record.Payload, &rec) != nil {
		return refuse("record schema")
	}
	ref, err := hashref.Parse(rec.Input)
	if err != nil {
		return refuse("input ref")
	}
	// The smaller input cap is checked on decoded object bytes, before parsing.
	input, reason := read(ref, "input", coordinator.InputV1, packetRosterReadLimit)
	if reason != "" {
		return refuse(reason)
	}
	if len(input.Payload) > packetRosterInputLimit {
		return refuse("input size")
	}
	var in packetRosterInput
	if strictRosterJSON(input.Payload, &in) != nil {
		return refuse("input shape")
	}
	if in.Goal == "" || len(in.Candidates) != 2 || in.Candidates[0] == "" || in.Candidates[1] == "" || in.Candidates[0] == in.Candidates[1] {
		return refuse("input schema")
	}
	if len(in.Roster) != 4 {
		return refuse("roster count")
	}
	seen := map[string]bool{}
	refs := map[string]hashref.HashRef{}
	for _, pair := range in.Roster {
		if seen[pair.Kind] {
			return refuse("duplicate kind")
		}
		seen[pair.Kind] = true
		// Transport transcription of world-demo/evidence.rosterKindAllowed;
		// TestWorldDemoLegacyGradeDrift runs the package's law against the kernel.
		switch pair.Kind {
		case "test-invocation", "recorded-effect", "authored-claim", "proof-envelope":
		default:
			return refuse("unknown kind")
		}
		ref, err := hashref.Parse(pair.Ref)
		if err != nil {
			return refuse("ref syntax")
		}
		refs[pair.Kind] = ref
	}
	effect, reason := read(refs["recorded-effect"], "effect", broker.EffectRecordV1, packetRosterReadLimit)
	if reason != "" {
		return refuse(reason)
	}
	er, err := broker.DecodeRecord(effect.Payload)
	if err != nil {
		return refuse("effect decode")
	}
	if !broker.RecordConsistent(er) {
		return refuse("effect inconsistent")
	}
	if _, reason = read(er.RequestRef, "request", broker.EffectRequestV1, packetRosterReadLimit); reason != "" {
		return refuse(reason)
	}
	if !er.ResultRef.IsZero() {
		if _, reason = read(er.ResultRef, "result", broker.EffectResultV1, packetRosterReadLimit); reason != "" {
			return refuse(reason)
		}
	}
	claim, reason := read(refs["authored-claim"], "claim", "world/demo/authored-claim/v1", packetRosterClaimLimit)
	if reason != "" {
		return refuse(reason)
	}
	var authored struct {
		Kind   string `json:"kind"`
		Author string `json:"author"`
		Text   string `json:"text"`
	}
	if strictRosterJSON(claim.Payload, &authored) != nil || authored.Text == "" {
		return refuse("claim shape")
	}
	if authored.Kind != "authored-claim" {
		return refuse("claim kind")
	}
	if authored.Author != "demo-authored" || claim.Provenance != "demo-authored" {
		return refuse("claim author")
	}
	if ctx.Err() != nil {
		return refuse("read deadline")
	}
	rows := make([]packetEvidenceRow, 0, 4)
	for _, pair := range in.Roster {
		row := packetEvidenceRow{kind: pair.Kind, ref: pair.Ref}
		switch pair.Kind {
		case "recorded-effect":
			row.label = "ATTESTED"
			row.producer = effect.Provenance
		case "authored-claim":
			row.label = "CLAIMED"
			row.producer = authored.Author
		default:
			row.unavailable = "resolution pending M5b"
		}
		rows = append(rows, row)
	}
	return packetRoster{rows: rows}
}
