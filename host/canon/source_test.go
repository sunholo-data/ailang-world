package canon

import (
	"bytes"
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

func mustCanon(t *testing.T, src string) []byte {
	t.Helper()
	out, err := SourceString(src)
	if err != nil {
		t.Fatalf("SourceString(%q) unexpected error: %v", src, err)
	}
	return out
}

// TestCanonNormalization exercises steps 3-7 across line endings, trailing
// whitespace, and terminal empty lines.
func TestCanonNormalization(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"already canonical", "line1\nline2\n", "line1\nline2\n"},
		{"crlf to lf", "line1\r\nline2\r\n", "line1\nline2\n"},
		{"lone cr to lf", "line1\rline2\r", "line1\nline2\n"},
		{"mixed cr crlf lf", "a\r\nb\rc\nd", "a\nb\nc\nd\n"},
		{"trailing spaces removed", "code   \nmore  \n", "code\nmore\n"},
		{"trailing tabs removed", "code\t\t\nmore\t\n", "code\nmore\n"},
		{"trailing mixed space tab", "code \t \t\n", "code\n"},
		{"terminal empty lines removed", "a\nb\n\n\n\n", "a\nb\n"},
		{"terminal blank-with-space lines removed", "a\n   \n\t\n", "a\n"},
		{"no trailing newline gets one", "solo", "solo\n"},
		{"interior blank lines preserved", "a\n\nb\n", "a\n\nb\n"},
		{"empty input yields empty", "", ""},
		{"only newlines yields empty", "\n\n\n", ""},
		{"only blank lines yields empty", "   \n\t\n  ", ""},
		{"leading blank lines preserved", "\n\na\n", "\n\na\n"},
		{"crlf with trailing space", "x  \r\n", "x\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := mustCanon(t, tc.in)
			if string(got) != tc.want {
				t.Fatalf("Source(%q) = %q, want %q", tc.in, string(got), tc.want)
			}
		})
	}
}

// TestCanonPreservesInteriorWhitespaceAndUnicode confirms step 8: only trailing
// ASCII space/tab are stripped; interior whitespace and multi-byte runes are
// preserved byte-for-byte.
func TestCanonPreservesInteriorWhitespaceAndUnicode(t *testing.T) {
	in := "  indented\tcode  \nrésumé — café ☕\n"
	want := "  indented\tcode\nrésumé — café ☕\n"
	got := mustCanon(t, in)
	if string(got) != want {
		t.Fatalf("Source(%q) = %q, want %q", in, string(got), want)
	}
	// The multi-byte runes must survive verbatim.
	if !bytes.Contains(got, []byte("résumé — café ☕")) {
		t.Errorf("multi-byte runes were altered: %q", string(got))
	}
}

func TestCanonRejectsInvalidUTF8(t *testing.T) {
	// 0x80 is a lone continuation byte: invalid UTF-8.
	_, err := Source([]byte{'a', 0x80, 'b'})
	assertCanonError(t, err, "invalid UTF-8")
}

func TestCanonRejectsBOM(t *testing.T) {
	in := append([]byte{utf8BOM0, utf8BOM1, utf8BOM2}, []byte("hello\n")...)
	_, err := Source(in)
	assertCanonError(t, err, "BOM")
}

func TestCanonRejectsNUL(t *testing.T) {
	_, err := Source([]byte("valid\x00text\n"))
	assertCanonError(t, err, "NUL")
}

func assertCanonError(t *testing.T, err error, what string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected error, got nil", what)
	}
	var ce *CanonicalizationError
	if !errors.As(err, &ce) {
		t.Fatalf("%s: expected *CanonicalizationError, got %T (%v)", what, err, err)
	}
}

// TestCanonIdempotence verifies canon(canon(x)) == canon(x) across a range of
// inputs, including ones that need normalization.
func TestCanonIdempotence(t *testing.T) {
	inputs := []string{
		"",
		"\n\n\n",
		"solo",
		"line1\r\nline2\r\n",
		"a\rb\rc",
		"trailing   \nspaces\t\t\n\n\n",
		"  interior  kept \nrésumé ☕  \n",
		"\n\nleading blanks\n",
	}
	for _, in := range inputs {
		once := mustCanon(t, in)
		twice, err := Source(once)
		if err != nil {
			t.Fatalf("second canonicalization of %q failed: %v", in, err)
		}
		if !bytes.Equal(once, twice) {
			t.Errorf("not idempotent for %q: canon=%q, canon(canon)=%q", in, string(once), string(twice))
		}
	}
}

// TestCanonGoldenHashRef pins a golden sha256 HashRef over canonicalized bytes,
// tying Decision 2 (canon) to Decision 3 (hashref). The messy input below
// canonicalizes to exactly "line1\nline2\n", whose sha256 is verified with
// shasum -a 256.
func TestCanonGoldenHashRef(t *testing.T) {
	const messy = "line1  \r\nline2\t\n\n\n"
	const wantCanon = "line1\nline2\n"
	// Verified: printf 'line1\nline2\n' | shasum -a 256
	const wantHashRef = "sha256:2751a3a2f303ad21752038085e2b8c5f98ecff61a2e4ebbd43506a941725be80"

	canon := mustCanon(t, messy)
	if string(canon) != wantCanon {
		t.Fatalf("canon = %q, want %q", string(canon), wantCanon)
	}
	ref := hashref.SumSHA256(canon)
	if ref.String() != wantHashRef {
		t.Fatalf("HashRef over canon bytes = %q, want %q", ref.String(), wantHashRef)
	}
}
