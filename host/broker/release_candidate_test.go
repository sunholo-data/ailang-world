package broker

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/pkgproj"
)

// The committed ready-packet golden names the world/core@0.1.1 CANDIDATE
// (D-WORLD-41). world/core@0.1.0 is published and immutable; its served record
// is testdata/metadata_world_core_0.1.0.json, which TestReconcilePublished010-
// OverFourDigests pins. These two arms state what reconciliation must say when
// the candidate is compared with that preserved record: never success.

func candidateReconcileConfig(t *testing.T) (ReconcileConfig, pkgproj.ReadyPacket) {
	t.Helper()
	golden, err := pkgproj.LoadReadyPacket(filepath.Join("..", "..", "scripts", "world_package_ready_packet.golden.json"))
	if err != nil {
		t.Fatal(err)
	}
	if golden.Version == "0.1.0" {
		t.Fatal("the golden names the PUBLISHED version 0.1.0; a candidate must name a new version")
	}
	vendor, name, _ := strings.Cut(golden.Package, "/")
	return ReconcileConfig{Vendor: vendor, Name: name, Version: golden.Version,
		Expected: PublishHashes{TarballSHA256: golden.TarballSHA256, ContentHash: golden.ContentHash,
			InterfaceHash: golden.InterfaceHash},
		ExpectedInterfaceV2: golden.InterfaceHashV2}, golden
}

// Arm 1: a served 0.1.0 document is not evidence that the candidate landed.
func TestReconcileCandidateAgainstPublished010IsConflict(t *testing.T) {
	c, golden := candidateReconcileConfig(t)
	r, _ := resolvePresent(ReconcileReceipt{}, served01(t), c)
	if r.State != ReconcileConflict || !strings.Contains(r.Detail, "world/core@0.1.0, want world/core@"+golden.Version) {
		t.Fatalf("%s %s", r.State, r.Detail)
	}
}

// Arm 2: the candidate's bytes reconciled AS 0.1.0 mismatch the served digests.
// Control in the same test: the preserved 0.1.0 expectations still succeed, so
// the conflict is the digests, not a broken fixture.
func TestReconcileCandidateBytesAsVersion010IsDigestConflict(t *testing.T) {
	if r, _ := resolvePresent(ReconcileReceipt{}, served01(t), cfg01()); r.State != ReconcileSucceededReconciled {
		t.Fatalf("control: preserved 0.1.0 record no longer reconciles: %s %s", r.State, r.Detail)
	}
	c, _ := candidateReconcileConfig(t)
	c.Version = "0.1.0"
	r, _ := resolvePresent(ReconcileReceipt{}, served01(t), c)
	if r.State != ReconcileConflict || !strings.Contains(r.Detail, "public metadata") {
		t.Fatalf("%s %s", r.State, r.Detail)
	}
}
