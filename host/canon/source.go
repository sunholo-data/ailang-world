// Package canon implements the canonical transition-source-byte procedure from
// Decision 2 of the w-world-library-m1 design. A transition function's identity
// (transitionFn: HashRef) is the hash of these canonical bytes, so the
// procedure must be exact, editor-independent, and idempotent.
//
// The eight steps (Decision 2):
//
//  1. Decode input as UTF-8; malformed input produces CanonicalizationError.
//  2. Reject a UTF-8 BOM and the NUL code point.
//  3. Convert CRLF and lone CR line endings to LF (0x0a).
//  4. Remove trailing ASCII space (0x20) and tab (0x09) from every line.
//  5. Remove empty lines at the end of the file.
//  6. Join remaining lines with LF.
//  7. A populated source ends with exactly one LF; an empty source yields zero bytes.
//  8. Every other Unicode code point retains its original UTF-8 byte sequence.
package canon

import (
	"fmt"
	"unicode/utf8"
)

// utf8BOM is the three-byte UTF-8 byte-order mark (U+FEFF), rejected by step 2.
const (
	utf8BOM0 = 0xEF
	utf8BOM1 = 0xBB
	utf8BOM2 = 0xBF
)

// CanonicalizationError is the structured error returned when the input cannot
// be canonicalized: invalid UTF-8, a leading BOM, or an embedded NUL.
type CanonicalizationError struct {
	// Reason is a stable, human-readable explanation.
	Reason string
	// Offset is the byte offset at which the fault was detected, when known
	// (-1 when not applicable).
	Offset int
}

func (e *CanonicalizationError) Error() string {
	if e.Offset >= 0 {
		return fmt.Sprintf("canon: %s at byte offset %d", e.Reason, e.Offset)
	}
	return "canon: " + e.Reason
}

// Source applies the eight-step canonicalization to src and returns the
// canonical byte stream, or a *CanonicalizationError if src is not admissible.
//
// The result is idempotent: Source(Source(x)) == Source(x).
func Source(src []byte) ([]byte, error) {
	// Step 2 (BOM): reject a leading UTF-8 BOM before any other processing.
	if len(src) >= 3 && src[0] == utf8BOM0 && src[1] == utf8BOM1 && src[2] == utf8BOM2 {
		return nil, &CanonicalizationError{Reason: "input begins with a UTF-8 BOM", Offset: 0}
	}

	// Step 1 (UTF-8) and step 2 (NUL): validate the whole stream up front.
	if !utf8.Valid(src) {
		off := firstInvalidUTF8Offset(src)
		return nil, &CanonicalizationError{Reason: "input is not valid UTF-8", Offset: off}
	}
	for i := 0; i < len(src); i++ {
		if src[i] == 0x00 {
			return nil, &CanonicalizationError{Reason: "input contains a NUL code point", Offset: i}
		}
	}

	// Empty input yields zero bytes (step 7).
	if len(src) == 0 {
		return []byte{}, nil
	}

	// Step 3: split into lines on LF, CRLF, and lone CR, normalizing all
	// three to logical line boundaries. Steps 4 and 8 are applied per line:
	// trailing ASCII space/tab are trimmed; all other bytes (which are valid
	// UTF-8 by the check above) are preserved verbatim.
	lines := splitLines(src)
	for i := range lines {
		lines[i] = trimTrailingSpaceTab(lines[i])
	}

	// Step 5: drop empty lines at the end of the file.
	end := len(lines)
	for end > 0 && len(lines[end-1]) == 0 {
		end--
	}
	lines = lines[:end]

	// Step 7: an empty (or all-blank) source yields zero bytes.
	if len(lines) == 0 {
		return []byte{}, nil
	}

	// Steps 6 and 7: join remaining lines with LF and terminate with exactly
	// one trailing LF.
	out := make([]byte, 0, len(src)+1)
	for _, line := range lines {
		out = append(out, line...)
		out = append(out, '\n')
	}
	return out, nil
}

// SourceString is a string-typed convenience wrapper over Source.
func SourceString(src string) ([]byte, error) {
	return Source([]byte(src))
}

// splitLines splits src into line contents, treating CRLF, lone CR, and LF as
// equivalent terminators. The terminators themselves are not included in the
// returned slices. A trailing terminator produces a final empty line (which the
// caller strips in step 5); no trailing terminator means the last segment is
// the final line.
func splitLines(src []byte) [][]byte {
	var lines [][]byte
	start := 0
	i := 0
	for i < len(src) {
		c := src[i]
		switch c {
		case '\n':
			lines = append(lines, src[start:i])
			i++
			start = i
		case '\r':
			lines = append(lines, src[start:i])
			// Consume a following LF so CRLF counts as one boundary.
			if i+1 < len(src) && src[i+1] == '\n' {
				i += 2
			} else {
				i++
			}
			start = i
		default:
			i++
		}
	}
	// The remaining tail after the last terminator is the final line.
	lines = append(lines, src[start:])
	return lines
}

// trimTrailingSpaceTab removes trailing ASCII space (0x20) and tab (0x09) bytes
// from a single line (step 4). Because these are single-byte ASCII values, they
// never occur inside a multi-byte UTF-8 sequence, so byte-level trimming is
// safe and preserves step 8's byte-sequence invariant for all other runes.
func trimTrailingSpaceTab(line []byte) []byte {
	end := len(line)
	for end > 0 && (line[end-1] == 0x20 || line[end-1] == 0x09) {
		end--
	}
	return line[:end]
}

// firstInvalidUTF8Offset returns the byte offset of the first invalid UTF-8
// sequence in b, or -1 if b is valid.
func firstInvalidUTF8Offset(b []byte) int {
	for i := 0; i < len(b); {
		r, size := utf8.DecodeRune(b[i:])
		if r == utf8.RuneError && size == 1 {
			return i
		}
		i += size
	}
	return -1
}
