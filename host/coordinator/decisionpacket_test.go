package coordinator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// TestDecisionPacketCodec pins the frozen world/decision-packet/v1 byte form
// (row 172 D2): the golden bytes, the round trip, every strict refusal, and the
// mint schedule (D3).
func TestDecisionPacketCodec(t *testing.T) {
	P := hashref.SumSHA256([]byte("plan"))
	Q := hashref.SumSHA256([]byte("req"))
	p := DecisionPacket{ProposalHash: P, RequestRef: Q, CreatedAt: 1700, DeadlineAt: 1760, Policy: "Cancel"}

	body := fmt.Sprintf(`{"proposalHash":"%s","requestRef":"%s","createdAt":1700,"deadlineAt":1760,"escalationsRemaining":0,"policy":"Cancel"}`, P, Q)
	want := fmt.Sprintf(`{"packetHash":"%s",`, hashref.SumSHA256([]byte(body))) + body[1:]
	enc := EncodeDecisionPacket(p)
	if string(enc) != want {
		t.Fatalf("golden bytes differ:\n got %s\nwant %s", enc, want)
	}

	t.Run("round-trip", func(t *testing.T) {
		got, err := DecodeDecisionPacket(enc)
		if err != nil {
			t.Fatalf("decode canonical: %v", err)
		}
		wantP := p
		wantP.PacketHash = hashref.SumSHA256([]byte(body))
		if got != wantP {
			t.Fatalf("round trip = %+v, want %+v", got, wantP)
		}
	})

	refusals := map[string]string{
		"whitespace":       strings.Replace(string(enc), `,"createdAt"`, `, "createdAt"`, 1),
		"unknown member":   strings.TrimSuffix(string(enc), "}") + `,"x":1}`,
		"trailing value":   string(enc) + "{}",
		"missing member":   mustMissingMember(t, enc),
		"malformed ref":    strings.Replace(string(enc), P.String(), "sha256:zz", 1),
		"unknown policy":   strings.Replace(string(enc), `"policy":"Cancel"`, `"policy":"cancel"`, 1),
		"wrong packetHash": strings.Replace(string(enc), hashref.SumSHA256([]byte(body)).String(), hashref.SumSHA256([]byte("wrong")).String(), 1),
	}
	for name, bad := range refusals {
		t.Run("refuses/"+strings.ReplaceAll(name, " ", "-"), func(t *testing.T) {
			if _, err := DecodeDecisionPacket([]byte(bad)); !errors.Is(err, errPacketEncoding) {
				t.Fatalf("decode %q: err = %v, want errPacketEncoding", bad, err)
			}
		})
	}

	t.Run("mintPacket", func(t *testing.T) {
		obj, err := mintPacket(1700, 1760, P, Q, "coordinator:a2a")
		if err != nil {
			t.Fatalf("mint: %v", err)
		}
		if obj.SemanticID != DecisionPacketV1 || obj.Provenance != "coordinator:a2a" {
			t.Fatalf("semantic id %q provenance %q", obj.SemanticID, obj.Provenance)
		}
		if obj.InterfaceHash != hashref.SumSHA256([]byte(DecisionPacketV1)) || obj.Hash != hashref.SumSHA256(obj.Payload) {
			t.Fatalf("interface hash or content hash does not match the payload")
		}
		got, err := DecodeDecisionPacket(obj.Payload)
		if err != nil {
			t.Fatalf("minted payload does not decode: %v", err)
		}
		if got.CreatedAt != 1700 || got.DeadlineAt != 1760 || got.EscalationsRemaining != 0 || got.Policy != "Cancel" {
			t.Fatalf("minted schedule = %+v", got)
		}
	})
	for _, c := range []struct{ now, expires int64 }{{1760, 1760}, {-1, 5}} {
		t.Run(fmt.Sprintf("mintPacket/refuses-%d-%d", c.now, c.expires), func(t *testing.T) {
			if _, err := mintPacket(c.now, c.expires, P, Q, "coordinator:a2a"); !errors.Is(err, errPacketSchedule) {
				t.Fatalf("err = %v, want errPacketSchedule", err)
			}
		})
	}
}

// mustMissingMember encodes the canonical members minus escalationsRemaining.
func mustMissingMember(t *testing.T, enc []byte) string {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(enc, &m); err != nil {
		t.Fatal(err)
	}
	delete(m, "escalationsRemaining")
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
