package approvaltest_test

import (
	"bytes"
	"errors"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
)

// TS1: the test signer's key parses as an authorized key, and its signatures verify through
// host/approval's own SSHSIG code. Deterministic: the same seed signs the same bytes the same way.
func TestTestSignerVerifies(t *testing.T) {
	s1, s2 := approvaltest.NewSigner(1), approvaltest.NewSigner(2)
	msg, err := approval.SignedMessage([32]byte{1}, [32]byte{2}, 1800000600)
	if err != nil {
		t.Fatal(err)
	}

	// (a) the .pub line parses to the same key, and ends with a newline.
	k, err := approval.ParseAuthorizedKey(s1.AuthorizedKey())
	if err != nil {
		t.Fatal(err)
	}
	if k.Type != "ssh-ed25519" || !bytes.Equal(k.Blob, s1.PublicKey().Blob) {
		t.Errorf("parsed key %s, want ssh-ed25519 with the signer's blob", k.Type)
	}
	if !bytes.HasSuffix(s1.AuthorizedKey(), []byte("\n")) {
		t.Errorf("AuthorizedKey does not end with a newline")
	}

	// (b) two seeds make two keys.
	if bytes.Equal(s1.PublicKey().Blob, s2.PublicKey().Blob) {
		t.Errorf("seeds 1 and 2 made the same key")
	}

	// (c) deterministic.
	sig := s1.Sign(approval.Namespace, msg)
	if !bytes.Equal(sig, approvaltest.NewSigner(1).Sign(approval.Namespace, msg)) {
		t.Errorf("the same seed signed the same message differently")
	}

	// (d) the signature parses and verifies under the signer's key, sha512, the Namespace.
	p, err := approval.ParseSSHSIG(sig)
	if err != nil {
		t.Fatal(err)
	}
	if p.HashAlg != "sha512" || p.Namespace != approval.Namespace || !bytes.Equal(p.PublicKey.Blob, s1.PublicKey().Blob) {
		t.Errorf("parsed sig: hash %q namespace %q, want sha512 and %q", p.HashAlg, p.Namespace, approval.Namespace)
	}
	if err := approval.VerifySSHSIG(p, approval.Namespace, msg, s1.PublicKey()); err != nil {
		t.Fatalf("verify: %v", err)
	}

	// (e) one flipped byte of the message does not verify.
	tampered := bytes.Clone(msg)
	tampered[0] ^= 1
	if err := approval.VerifySSHSIG(p, approval.Namespace, tampered, s1.PublicKey()); !errors.Is(err, approval.ErrBadSignature) {
		t.Errorf("tampered message: %v, want ErrBadSignature", err)
	}

	// (f) against another key, the key-mismatch refusal.
	if err := approval.VerifySSHSIG(p, approval.Namespace, msg, s2.PublicKey()); !errors.Is(err, approval.ErrKeyMismatch) {
		t.Errorf("other key: %v, want ErrKeyMismatch", err)
	}

	// (g) the approval entry point: the same digest, challenge and expiry verify; a later expiry does not.
	if err := approval.VerifyApproval(sig, [32]byte{1}, [32]byte{2}, 1800000600, s1.PublicKey()); err != nil {
		t.Errorf("VerifyApproval: %v", err)
	}
	if err := approval.VerifyApproval(sig, [32]byte{1}, [32]byte{2}, 1800000601, s1.PublicKey()); !errors.Is(err, approval.ErrBadSignature) {
		t.Errorf("later expiry: %v, want ErrBadSignature", err)
	}
}

// TS3 (a control): no production file imports approvaltest. The walk covers the repository from
// its go.mod down, skipping dot-directories, node_modules and vendor, and parses imports only.
// Its own controls: at least 300 .go files were parsed, and at least one _test.go file imports it.
func TestApprovalTestIsTestOnly(t *testing.T) {
	const self = "github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(root)
		if parent == root {
			t.Fatal("no go.mod above the working directory")
		}
		root = parent
	}
	fset := token.NewFileSet()
	parsed, testImporters := 0, 0
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if d.IsDir() {
			if path != root && (strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(name, ".go") {
			return nil
		}
		f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		parsed++
		for _, imp := range f.Imports {
			if strings.Trim(imp.Path.Value, `"`) != self {
				continue
			}
			if strings.HasSuffix(name, "_test.go") {
				testImporters++
			} else {
				t.Errorf("%s imports approvaltest, which is test support only", strings.TrimPrefix(path, root+"/"))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if parsed < 300 {
		t.Errorf("parsed %d .go files, want at least 300 (the walk is not seeing the repository)", parsed)
	}
	if testImporters < 1 {
		t.Errorf("no _test.go file imports approvaltest; the control is vacuous")
	}
}
