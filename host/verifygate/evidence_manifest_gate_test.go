package verifygate

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const syntheticEvidencePackage = "github.com/sunholo-data/ailang-world/host/evidence"

var syntheticEvidenceTests = []string{
	"TestAttackerChosenValidatorCannotMintForHostAuthority",
	"TestConstructorNamesActuallyUsedUnorderedTimeouts",
	"TestConstructorPinsBusyTimeoutBelowObjectReadTimeout",
	"TestConstructorRefusesEmptyCompilerVersion",
	"TestConstructorRefusesEmptyRequiredIdentities",
	"TestConstructorRefusesNilReader",
	"TestConstructorRefusesNonPositiveObjectReadTimeout",
	"TestConstructorRefusesUnknownBusyTimeout",
	"TestConstructorRefusesUnsetCompilerIdentity",
	"TestDecodeProposalCapsBeforeParse",
	"TestEnvelopeCanonicalRoundTripAndMACDeferral",
	"TestEnvelopeCarriesTheReportItAlreadyDecoded",
	"TestEnvelopeStrictRefusals",
	"TestFailedProofReportIsRefused",
	"TestIncompleteProofReportIsRefused",
	"TestInvalidProofRefIsRefused",
	"TestMalformedProofReportIsRefused",
	"TestMismatchedProofSubjectIsRefused",
	"TestMismatchedProofToolIsRefused",
	"TestMissingProofReportIsRefused",
	"TestNestingDepthBombWithinByteCapIsRefused",
	"TestOtherwisePerfectReportWithWrongMACIsUnauthenticated",
	"TestOtherwisePerfectReportWithoutMACIsUnauthenticated",
	"TestOversizeProofReportIsRefused",
	"TestPayloadHashMismatchIsRefused",
	"TestProofReportCanonicalRoundTrip",
	"TestProofReportCaps",
	"TestProofReportStrictRefusals",
	"TestProposalStrictRefusals",
	"TestPublicAuthoritySurfaceIsFrozen",
	"TestReaderWaitBoundsCannotBeLostThroughWrapper",
	"TestRealStoreBlockedObjectReadReturnsWithinObjectReadTimeout",
	"TestTruncatedTailReportIsRefusedNotPanicked",
	"TestValidatorMintIdentitiesAreDistinct",
	"TestWrongInterfaceIsRefused",
	"TestWrongSemanticIDIsRefused",
	"TestZeroValueForgeryCannotResolve",
}

func newIsolatedEvidenceGateRoot(t *testing.T) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), "iso")
	if source := os.Getenv("EVIDENCE_MANIFEST_GATE_SOURCE"); source != "" {
		in, err := os.ReadFile(source)
		if err != nil {
			t.Fatalf("read overridden evidence gate source: %v", err)
		}
		dst := filepath.Join(root, "scripts", "verify_go.sh")
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(dst, in, 0o755); err != nil {
			t.Fatal(err)
		}
	} else {
		copyGateFile(t, root, "scripts/verify_go.sh", 0o755)
	}
	files := 0
	err := filepath.Walk(root, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files++
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if files != 1 {
		t.Fatalf("isolated evidence-gate copy landed %d files, want 1", files)
	}
	return root
}

func writeSyntheticEvidenceEvents(t *testing.T, root string, tests []string) string {
	t.Helper()
	path := filepath.Join(root, "observed.json")
	var raw bytes.Buffer
	enc := json.NewEncoder(&raw)
	for _, name := range tests {
		if err := enc.Encode(map[string]string{"Action": "pass", "Package": syntheticEvidencePackage, "Test": name}); err != nil {
			t.Fatal(err)
		}
	}
	if err := enc.Encode(map[string]string{"Action": "pass", "Package": syntheticEvidencePackage}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, raw.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func runSyntheticEvidenceGate(t *testing.T, root string, tests []string, exact int) (int, string) {
	t.Helper()
	observed := writeSyntheticEvidenceEvents(t, root, tests)
	cmd := exec.Command(filepath.Join(root, "scripts", "verify_go.sh"), "--evidence-manifest-check", observed, strconv.Itoa(exact))
	cmd.Dir = root
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	err := cmd.Run()
	if err == nil {
		return 0, output.String()
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), output.String()
	}
	t.Fatalf("start isolated evidence gate: %v", err)
	return -1, output.String()
}

func requireSyntheticEvidenceControl(t *testing.T, root string) {
	t.Helper()
	rc, out := runSyntheticEvidenceGate(t, root, syntheticEvidenceTests, len(syntheticEvidenceTests))
	if rc != 0 || !strings.Contains(out, "all 37 required top-level evidence tests passed exactly once") {
		t.Fatalf("pristine synthetic control did not reach success: rc=%d\n%s", rc, out)
	}
}

func emptyCopiedEvidenceRequiredSet(t *testing.T, root string) {
	t.Helper()
	path := filepath.Join(root, "scripts", "verify_go.sh")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	start := strings.Index(string(raw), "REQUIRED_EVIDENCE_TESTS = {")
	endMarker := "}\nEXACT_EVIDENCE_TESTS = int(sys.argv[2])"
	end := strings.Index(string(raw), endMarker)
	if start < 0 || end < start {
		t.Fatalf("required-set mutation anchors missing")
	}
	mutant := string(raw[:start]) + "REQUIRED_EVIDENCE_TESTS = {}\nEXACT_EVIDENCE_TESTS = int(sys.argv[2])" + string(raw[end+len(endMarker):])
	if err := os.WriteFile(path, []byte(mutant), 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestEvidenceNamedManifestRejectsUnpinnedTest(t *testing.T) {
	t.Run("extra observed test", func(t *testing.T) {
		root := newIsolatedEvidenceGateRoot(t)
		requireSyntheticEvidenceControl(t, root)
		mutant := append(append([]string{}, syntheticEvidenceTests...), "TestUnpinnedEvidenceProbe")
		rc, out := runSyntheticEvidenceGate(t, root, mutant, len(syntheticEvidenceTests))
		if rc != 1 || !strings.Contains(out, "evidence test set differs from REQUIRED_EVIDENCE_TESTS") || !strings.Contains(out, "TestUnpinnedEvidenceProbe") {
			t.Fatalf("extra test was not refused: rc=%d\n%s", rc, out)
		}
	})
	t.Run("missing required test", func(t *testing.T) {
		root := newIsolatedEvidenceGateRoot(t)
		requireSyntheticEvidenceControl(t, root)
		rc, out := runSyntheticEvidenceGate(t, root, syntheticEvidenceTests[1:], len(syntheticEvidenceTests))
		if rc != 1 || !strings.Contains(out, "missing=['TestAttackerChosenValidatorCannotMintForHostAuthority']") {
			t.Fatalf("missing test was not refused: rc=%d\n%s", rc, out)
		}
	})
	t.Run("duplicate observed identity", func(t *testing.T) {
		root := newIsolatedEvidenceGateRoot(t)
		requireSyntheticEvidenceControl(t, root)
		mutant := append(append([]string{}, syntheticEvidenceTests...), syntheticEvidenceTests[0])
		rc, out := runSyntheticEvidenceGate(t, root, mutant, len(syntheticEvidenceTests))
		if rc != 1 || !strings.Contains(out, "duplicate=['TestAttackerChosenValidatorCannotMintForHostAuthority']") {
			t.Fatalf("duplicate test was not refused: rc=%d\n%s", rc, out)
		}
	})
	t.Run("empty required set", func(t *testing.T) {
		root := newIsolatedEvidenceGateRoot(t)
		requireSyntheticEvidenceControl(t, root)
		emptyCopiedEvidenceRequiredSet(t, root)
		rc, out := runSyntheticEvidenceGate(t, root, syntheticEvidenceTests, 0)
		if rc != 1 || !strings.Contains(out, "FATAL INSTRUMENT FAILURE: REQUIRED_EVIDENCE_TESTS is empty") {
			t.Fatalf("empty required set did not fail loudly: rc=%d\n%s", rc, out)
		}
	})
	t.Run("empty observed enumeration", func(t *testing.T) {
		root := newIsolatedEvidenceGateRoot(t)
		requireSyntheticEvidenceControl(t, root)
		rc, out := runSyntheticEvidenceGate(t, root, nil, len(syntheticEvidenceTests))
		if rc != 1 || !strings.Contains(out, "FATAL INSTRUMENT FAILURE: observed terminal named-test pass set is empty") {
			t.Fatalf("empty observed enumeration did not fail loudly: rc=%d\n%s", rc, out)
		}
	})
}
