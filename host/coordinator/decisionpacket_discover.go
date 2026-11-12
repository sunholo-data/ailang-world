package coordinator

// Row 201 A4 (design_docs/planned/w-world-approvals-by-key-a4.md §2.4 [r1], D4a): the complete,
// request-wide decision-packet discovery a decision runs under humanMu. Packets have no request
// index and ValidatePacket's duplicate scan is a 100-packet prefix, so a decision never trusts a
// bounded view: it pages through every stored packet, up to DecisionDiscoveryLimit plus one
// sentinel read, and reports whether the set it counted is complete.

import (
	"context"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// DiscoveryReader is the read surface the discovery needs; PacketReader and the daemon's
// readStore satisfy it.
type DiscoveryReader interface {
	GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error)
	ObjectsBySemanticID(ctx context.Context, id, after string, limit int) ([]store.Object, error)
}

// DiscoverRequestPackets counts the stored decision packets that name request, in keyset pages of
// store.MaxSemanticIDPage following the cursor. Complete is true only when a page came back short
// before the scan reached broker.DecisionDiscoveryLimit packets; a packet past the limit is the
// overflow sentinel, and the count is then not established. carrying is the last carrying packet
// read (the one, when Carrying is 1). Store errors are returned, never a count.
func DiscoverRequestPackets(ctx context.Context, r DiscoveryReader, request hashref.HashRef) (broker.Discovery, store.Object, error) {
	var d broker.Discovery
	var carrying store.Object
	after := ""
	for {
		page, err := r.ObjectsBySemanticID(ctx, DecisionPacketV1, after, store.MaxSemanticIDPage)
		if err != nil {
			return broker.Discovery{}, store.Object{}, err
		}
		for _, o := range page {
			if d.Scanned >= broker.DecisionDiscoveryLimit {
				return d, carrying, nil // the sentinel: more packets than the limit
			}
			d.Scanned++
			full, ok, err := r.GetObject(ctx, o.Hash)
			if err != nil {
				return broker.Discovery{}, store.Object{}, err
			}
			p, decodeErr := DecodeDecisionPacket(full.Payload)
			if !ok || decodeErr != nil {
				d.Undecodable++
				continue
			}
			if p.RequestRef == request {
				d.Carrying++
				carrying = full
			}
		}
		if len(page) < store.MaxSemanticIDPage {
			d.Complete = true
			return d, carrying, nil
		}
		after = page[len(page)-1].Hash.String()
	}
}
