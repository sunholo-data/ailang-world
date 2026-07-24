package hashref

import (
	"errors"
	"testing"
)

// Golden SHA-256 digests (lowercase hex), verified with shasum -a 256.
const (
	sha256Empty = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	sha256ABC   = "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
)

func TestSumSHA256Golden(t *testing.T) {
	cases := []struct {
		name    string
		payload []byte
		want    string
	}{
		{"empty", []byte(""), "sha256:" + sha256Empty},
		{"abc", []byte("abc"), "sha256:" + sha256ABC},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SumSHA256(tc.payload)
			if got.String() != tc.want {
				t.Fatalf("SumSHA256(%q) = %q, want %q", tc.payload, got.String(), tc.want)
			}
			if got.Algo() != AlgoSHA256 {
				t.Errorf("Algo() = %q, want %q", got.Algo(), AlgoSHA256)
			}
			if len(got.Digest()) != 64 {
				t.Errorf("Digest() length = %d, want 64", len(got.Digest()))
			}
		})
	}
}

func TestSumUnsupportedAlgo(t *testing.T) {
	_, err := Sum("md5", []byte("x"))
	if err == nil {
		t.Fatal("Sum with unsupported algo: expected error, got nil")
	}
	var he *HashError
	if !errors.As(err, &he) {
		t.Fatalf("expected *HashError, got %T", err)
	}
}

func TestParseRoundTrip(t *testing.T) {
	text := "sha256:" + sha256Empty
	h, err := Parse(text)
	if err != nil {
		t.Fatalf("Parse(%q) unexpected error: %v", text, err)
	}
	if h.String() != text {
		t.Errorf("round trip: got %q, want %q", h.String(), text)
	}
	if h.Algo() != AlgoSHA256 {
		t.Errorf("Algo() = %q, want sha256", h.Algo())
	}
	if h.Digest() != sha256Empty {
		t.Errorf("Digest() = %q, want %q", h.Digest(), sha256Empty)
	}
}

// TestParseRejects covers every rejection required by Decision 3: malformed
// text, unsupported tags, uppercase hex, and bare digests.
func TestParseRejects(t *testing.T) {
	upperDigest := "E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855"
	cases := []struct {
		name string
		text string
	}{
		{"empty string", ""},
		{"bare digest (no tag)", sha256Empty},
		{"uppercase hex digest", "sha256:" + upperDigest},
		{"unsupported tag", "md5:" + "d41d8cd98f00b204e9800998ecf8427e"},
		{"empty algo", ":" + sha256Empty},
		{"empty digest", "sha256:"},
		{"multiple colons", "sha256:" + sha256Empty + ":extra"},
		{"digest too short", "sha256:abc123"},
		{"digest too long", "sha256:" + sha256Empty + "00"},
		{"non-hex character g", "sha256:" + "g3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"leading whitespace", " sha256:" + sha256Empty},
		{"trailing whitespace", "sha256:" + sha256Empty + " "},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.text)
			if err == nil {
				t.Fatalf("Parse(%q): expected error, got nil", tc.text)
			}
			var he *HashError
			if !errors.As(err, &he) {
				t.Fatalf("Parse(%q): expected *HashError, got %T (%v)", tc.text, err, err)
			}
		})
	}
}

func TestNewValidation(t *testing.T) {
	if _, err := New(AlgoSHA256, sha256Empty); err != nil {
		t.Errorf("New with valid inputs: unexpected error: %v", err)
	}
	if _, err := New("sha512", sha256Empty); err == nil {
		t.Error("New with unsupported algo: expected error")
	}
	if _, err := New(AlgoSHA256, "ABC"); err == nil {
		t.Error("New with uppercase hex: expected error")
	}
	if _, err := New(AlgoSHA256, "abcd"); err == nil {
		t.Error("New with wrong-width digest: expected error")
	}
}

func TestZeroValue(t *testing.T) {
	var h HashRef
	if !h.IsZero() {
		t.Error("zero value IsZero() = false, want true")
	}
	if h.String() != "" {
		t.Errorf("zero value String() = %q, want empty", h.String())
	}
}

func TestSupported(t *testing.T) {
	if !Supported(AlgoSHA256) {
		t.Error("sha256 should be supported")
	}
	if Supported("md5") {
		t.Error("md5 should not be supported")
	}
}

func TestMustParsePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("MustParse on malformed input: expected panic")
		}
	}()
	MustParse("not-a-hashref")
}
