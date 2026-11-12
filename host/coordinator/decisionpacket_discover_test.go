package coordinator

// Row 201 A4 sprint 1 (design §2.4 [r1], §6 DS1): the discovery's completeness, its sentinel, its
// cursor and its undecodable count, over a fake keyset reader with the store's page contract.

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// keysetPackets is a fake DiscoveryReader: packets sorted by hash_ref, served in keyset pages of at
// most store.MaxSemanticIDPage; it records every after it was asked for.
type keysetPackets struct {
	objs   []store.Object
	afters []string
	gets   int
}

func (k *keysetPackets) ObjectsBySemanticID(_ context.Context, id, after string, limit int) ([]store.Object, error) {
	if id != DecisionPacketV1 || limit < 1 || limit > store.MaxSemanticIDPage {
		return nil, fmt.Errorf("unexpected page %q %d", id, limit)
	}
	k.afters = append(k.afters, after)
	var page []store.Object
	for _, o := range k.objs {
		if o.Hash.String() > after && len(page) < limit {
			page = append(page, store.Object{Hash: o.Hash, SemanticID: o.SemanticID})
		}
	}
	return page, nil
}

func (k *keysetPackets) GetObject(_ context.Context, ref hashref.HashRef) (store.Object, bool, error) {
	k.gets++
	for _, o := range k.objs {
		if o.Hash == ref {
			return o, true, nil
		}
	}
	return store.Object{}, false, nil
}

// discoveryRig stores n packets, the ones whose index is in carrying naming request, the rest
// naming their own request; undecodable payloads replace the packets whose index is in bad.
func discoveryRig(n int, request hashref.HashRef, carrying, bad map[int]bool) *keysetPackets {
	k := &keysetPackets{}
	for i := 0; i < n; i++ {
		req := hashref.SumSHA256([]byte(fmt.Sprintf("other-request-%d", i)))
		if carrying[i] {
			req = request
		}
		payload := EncodeDecisionPacket(DecisionPacket{ProposalHash: hashref.SumSHA256([]byte(fmt.Sprintf("plan-%d", i))), RequestRef: req,
			CreatedAt: 1, DeadlineAt: 2, Policy: "Cancel"})
		if bad[i] {
			payload = []byte(fmt.Sprintf(`{"not":"a packet %d"}`, i))
		}
		k.objs = append(k.objs, store.Object{Hash: hashref.SumSHA256(payload), SemanticID: DecisionPacketV1, Payload: payload})
	}
	sort.Slice(k.objs, func(i, j int) bool { return k.objs[i].Hash.String() < k.objs[j].Hash.String() })
	return k
}

// DS1: 0, 499, 500, 1,999 and exactly 2,000 packets are complete (the sentinel read returns none);
// 2,001 are not; the cursor follows after across pages; an undecodable payload is counted.
func TestDiscoverRequestPacketsCompleteness(t *testing.T) {
	ctx := context.Background()
	request := hashref.SumSHA256([]byte("the-request"))
	for _, tc := range []struct {
		n, carryAt, wantScanned, wantPages int
		complete                           bool
	}{
		{0, -1, 0, 1, true}, {499, 0, 499, 1, true}, {500, 499, 500, 2, true}, {1999, 1998, 1999, 4, true},
		{2000, 1999, 2000, 5, true}, {2001, 1999, 2000, 5, false},
	} {
		k := discoveryRig(tc.n, request, map[int]bool{tc.carryAt: true}, nil)
		d, carrying, err := DiscoverRequestPackets(ctx, k, request)
		if err != nil {
			t.Fatal(err)
		}
		wantCarrying := 0
		if tc.carryAt >= 0 {
			wantCarrying = 1
		}
		if d.Complete != tc.complete || d.Scanned != tc.wantScanned || d.Carrying != wantCarrying || d.Undecodable != 0 || len(k.afters) != tc.wantPages {
			t.Errorf("%d packets: %+v after %d pages, want complete=%v scanned=%d carrying=%d pages=%d", tc.n, d, len(k.afters), tc.complete, tc.wantScanned, wantCarrying, tc.wantPages)
		}
		if wantCarrying == 1 && d.Complete {
			p, err := DecodeDecisionPacket(carrying.Payload)
			if err != nil || p.RequestRef != request {
				t.Errorf("%d packets: the carrying packet is %v, %v", tc.n, p.RequestRef, err)
			}
		}
		// The cursor: each page starts strictly after the last hash of the one before.
		for i := 1; i < len(k.afters); i++ {
			if want := k.objs[i*store.MaxSemanticIDPage-1].Hash.String(); k.afters[i] != want {
				t.Errorf("%d packets: page %d after %q, want %q", tc.n, i, k.afters[i], want)
			}
		}
		if k.gets != d.Scanned {
			t.Errorf("%d packets: %d payload reads for %d scanned", tc.n, k.gets, d.Scanned)
		}
	}
	// Two carrying packets, one past the first page, are both counted.
	k := discoveryRig(1200, request, map[int]bool{3: true, 1100: true}, nil)
	if d, _, err := DiscoverRequestPackets(ctx, k, request); err != nil || d != (broker.Discovery{Complete: true, Scanned: 1200, Carrying: 2}) {
		t.Errorf("two carrying: %+v, %v", d, err)
	}
	// An undecodable payload is counted, not skipped; so is a listed packet whose payload is gone.
	k = discoveryRig(600, request, map[int]bool{7: true}, map[int]bool{550: true})
	if d, _, err := DiscoverRequestPackets(ctx, k, request); err != nil || !d.Complete || d.Undecodable != 1 || d.Carrying != 1 || d.Scanned != 600 {
		t.Errorf("undecodable: %+v, %v", d, err)
	}
	k = discoveryRig(3, request, map[int]bool{0: true}, nil)
	gone := k.objs[2].Hash
	if d, _, err := DiscoverRequestPackets(ctx, &missingPayload{k, gone}, request); err != nil || d.Undecodable != 1 || d.Scanned != 3 {
		t.Errorf("payload gone: %+v, %v", d, err)
	}
	// A store error is returned, never a count.
	if _, _, err := DiscoverRequestPackets(ctx, failingPages{}, request); err == nil || !strings.Contains(err.Error(), "page failed") {
		t.Errorf("store error: %v", err)
	}
	if _, _, err := DiscoverRequestPackets(ctx, failingGets{discoveryRig(3, request, nil, nil)}, request); err == nil || !strings.Contains(err.Error(), "read failed") {
		t.Errorf("payload read error: %v", err)
	}
}

// failingGets lists packets but fails every payload read.
type failingGets struct{ *keysetPackets }

func (failingGets) GetObject(context.Context, hashref.HashRef) (store.Object, bool, error) {
	return store.Object{}, false, fmt.Errorf("read failed")
}

// missingPayload lists every packet but answers not-found for one of them.
type missingPayload struct {
	*keysetPackets
	gone hashref.HashRef
}

func (m *missingPayload) GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error) {
	if ref == m.gone {
		return store.Object{}, false, nil
	}
	return m.keysetPackets.GetObject(ctx, ref)
}

type failingPages struct{}

func (failingPages) ObjectsBySemanticID(context.Context, string, string, int) ([]store.Object, error) {
	return nil, fmt.Errorf("page failed")
}

func (failingPages) GetObject(context.Context, hashref.HashRef) (store.Object, bool, error) {
	return store.Object{}, false, nil
}
