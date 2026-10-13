package coordinator

import (
	"context"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// emptyPacketReader answers every read with absence. The arms below refuse at
// the bytes check (step 1), before any read, so no store state is needed.
type emptyPacketReader struct{}

func (emptyPacketReader) GetObject(context.Context, hashref.HashRef) (store.Object, bool, error) {
	return store.Object{}, false, nil
}
func (emptyPacketReader) GetLogEntry(context.Context, int64) (store.LogEntry, bool, error) {
	return store.LogEntry{}, false, nil
}
func (emptyPacketReader) GetRegistryHead(context.Context, string) (hashref.HashRef, bool, error) {
	return hashref.HashRef{}, false, nil
}
func (emptyPacketReader) ObjectsBySemanticID(context.Context, string, string, int) ([]store.Object, error) {
	return nil, nil
}
func (emptyPacketReader) ObjectCommits(context.Context, hashref.HashRef, int64, int) ([]int64, error) {
	return nil, nil
}

// TestValidatePacketBytesCheckNeedsThePacketLabel (M12): a stored object whose
// payload decodes as a packet is refused at step 1 unless it carries the packet
// label and the packet interface hash.
func TestValidatePacketBytesCheckNeedsThePacketLabel(t *testing.T) {
	p := DecisionPacket{ProposalHash: hashref.SumSHA256([]byte("plan")), RequestRef: hashref.SumSHA256([]byte("req")),
		CreatedAt: 1700, DeadlineAt: 1760, Policy: "Cancel"}
	payload := EncodeDecisionPacket(p)
	packetIface := hashref.SumSHA256([]byte(DecisionPacketV1))
	arms := []struct {
		name  string
		label string
		iface hashref.HashRef
	}{
		{"wrong-semantic-id", "test/not-a-packet", packetIface},
		{"wrong-interface-hash", DecisionPacketV1, hashref.SumSHA256([]byte("test/not-a-packet"))},
	}
	for _, arm := range arms {
		t.Run(arm.name, func(t *testing.T) {
			obj := store.Object{Hash: hashref.SumSHA256(payload), InterfaceHash: arm.iface, SemanticID: arm.label,
				Provenance: "validate-test", Payload: payload}
			res, err := ValidatePacket(boundedTestContext(t), emptyPacketReader{}, obj)
			if err != nil {
				t.Fatal(err)
			}
			if res.Code != 1 {
				t.Fatalf("code = %d, want 1 (bytes) for %s", res.Code, arm.name)
			}
		})
	}
}
