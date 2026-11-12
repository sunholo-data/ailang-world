package daemon

// Row 201 A4 sprint 1: the decision routes' test rig. Every store is under t.TempDir(); every
// packet is a forged, validated commit (workbench_packet_test.go's forgedChain); every operator
// proof is written as the CLI's writeOperatorProof writes it (D-WORLD-109: no real store or tty).

import (
	"crypto/sha256"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// decNow is the decision routes' pinned clock: after the packets' createdAt (packetNow), before
// their deadline (packetNow + 60).
const decNow = packetNow + 10

type decideRig struct{ approvalRig }

func newDecideRig(t *testing.T) decideRig {
	t.Helper()
	a := newApprovalRig(t)
	a.d.approveNow = func() int64 { return decNow }
	return decideRig{a}
}

// packet files a pending Human.Approve request (scope, requester) in the approvals chain and
// commits a validated packet carrying it, whose proposal is label's.
func (r decideRig) packet(label, scope, requester string) (hashref.HashRef, store.Object) {
	r.t.Helper()
	q := approvalRequestInChain(r.t, r.d, scope, requester, packetNow)
	p, _ := commitForgedChain(r.t, r.d, newForgedChain(label, q, scope))
	return q, p
}

// proof writes the operator proof for a decision of q with phase, beside the credentials store,
// as writeOperatorProof does (0700 directory, 0600 file, expiry within 300 s); it returns the secret.
func (r decideRig) proof(q hashref.HashRef, phase string) string {
	r.t.Helper()
	secret := hashref.SumSHA256([]byte(phase + q.String())).Digest()
	dir := filepath.Join(filepath.Dir(r.d.credentialsPath), "operator-proofs")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		r.t.Fatal(err)
	}
	data := mustMarshal(r.t, operatorProof{SHA256: hashref.SumSHA256([]byte(secret)).Digest(), Phase: phase, ExpiresAt: time.Now().Unix() + 200})
	if err := os.WriteFile(r.proofPath(q), data, 0o600); err != nil {
		r.t.Fatal(err)
	}
	return secret
}

func (r decideRig) proofPath(q hashref.HashRef) string {
	return filepath.Join(filepath.Dir(r.d.credentialsPath), "operator-proofs", "decision-"+q.Digest()+".json")
}

// proofExists reports whether q's proof file is still there (nothing consumed it).
func (r decideRig) proofExists(q hashref.HashRef) bool {
	_, err := os.Lstat(r.proofPath(q))
	return err == nil
}

// decideBody is a typed POST body for (verb, p, q, expiry).
func decideBody(verb string, p store.Object, q hashref.HashRef, expiry int64) []byte {
	return []byte(fmt.Sprintf(`{"v":1,"verb":%q,"packet":%q,"request":%q,"expiry":%d,"typed":true}`, verb, p.Hash.Digest(), q.Digest(), expiry))
}

// post is POST /surface/decide/<ref> from the loopback with A3a's headers and the proof; hdr adds
// or replaces headers.
func (r decideRig) post(ref string, body []byte, proof string, hdr ...string) *httptest.ResponseRecorder {
	r.t.Helper()
	h := []string{"Origin", "http://127.0.0.1:7644", "Content-Type", "application/json"}
	if proof != "" {
		h = append(h, "X-World-Operator-Proof", proof)
	}
	return r.serve("POST", "/surface/decide/"+ref, apLoopback, body, append(h, hdr...)...)
}

// decisions is the number of stored approval decisions and the approvals head.
func (r decideRig) decisions() (int, hashref.HashRef) {
	r.t.Helper()
	ctx := boundedTestContext(r.t)
	objs, err := r.d.store.ObjectsBySemanticID(ctx, broker.ApprovalDecisionV1, "", store.MaxSemanticIDPage)
	if err != nil {
		r.t.Fatal(err)
	}
	head, _, err := r.d.store.GetRegistryHead(ctx, broker.ApprovalsV1)
	if err != nil {
		r.t.Fatal(err)
	}
	return len(objs), head
}

// fillPackets stores n decodable decision packets that carry other requests (no commit: the
// discovery reads every stored packet) and returns their hashes, sorted.
func (r decideRig) fillPackets(label string, n int) []string {
	r.t.Helper()
	ctx := boundedTestContext(r.t)
	var hashes []string
	for i := 0; i < n; i++ {
		payload := coordinator.EncodeDecisionPacket(coordinator.DecisionPacket{ProposalHash: hashref.SumSHA256([]byte(fmt.Sprintf("%s-plan-%d", label, i))),
			RequestRef: hashref.SumSHA256([]byte(fmt.Sprintf("%s-request-%d", label, i))), CreatedAt: packetNow, DeadlineAt: packetNow + 60, Policy: "Cancel"})
		o := forgeObject(coordinator.DecisionPacketV1, payload)
		if err := r.d.store.PutObject(ctx, o); err != nil {
			r.t.Fatal(err)
		}
		hashes = append(hashes, o.Hash.String())
	}
	sort.Strings(hashes)
	return hashes
}

// packetAfter commits a validated packet carrying q whose hash sorts after position pos among
// fillers (sorted), trying proposal labels until one does; it returns the packet.
func (r decideRig) packetAfter(label string, q hashref.HashRef, scope string, fillers []string, pos int) store.Object {
	r.t.Helper()
	for i := 0; i < 200; i++ {
		f := newForgedChain(fmt.Sprintf("%s-%d", label, i), q, scope)
		h := hashref.SumSHA256(coordinator.EncodeDecisionPacket(f.packet)).String()
		if sort.SearchStrings(fillers, h) >= pos {
			p, _ := commitForgedChain(r.t, r.d, f)
			return p
		}
	}
	r.t.Fatal("no proposal label sorts the packet late enough")
	return store.Object{}
}

// packetPosition is p's 0-based position among every stored packet in hash_ref order.
func (r decideRig) packetPosition(p store.Object) int {
	r.t.Helper()
	ctx := boundedTestContext(r.t)
	pos, after := 0, ""
	for {
		page, err := r.d.store.ObjectsBySemanticID(ctx, coordinator.DecisionPacketV1, after, store.MaxSemanticIDPage)
		if err != nil {
			r.t.Fatal(err)
		}
		for _, o := range page {
			if o.Hash == p.Hash {
				return pos
			}
			pos++
		}
		if len(page) < store.MaxSemanticIDPage {
			r.t.Fatalf("packet %s is not stored", p.Hash)
		}
		after = page[len(page)-1].Hash.String()
	}
}

// sha256Of is the ref of payload.
func sha256Of(payload []byte) string {
	sum := sha256.Sum256(payload)
	return "sha256:" + fmt.Sprintf("%x", sum[:])
}
