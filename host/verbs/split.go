// Package verbs is the operator verbs' one grammar (row 243 V1, design §2.1): the address bar's
// tokenizer (Split, Quote), the verb table (Lookup), the CLI's interleaved flag parse
// (ParseInterleaved) and episode report's flag set with the bar's admission. It is pure: no I/O, no
// clock, no randomness. The CLI and the daemon's GET /surface/verb.json both call it, so the bar
// and a terminal read one grammar.
package verbs

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

// MaxLine and MaxWords bound one line of the bar (design §2.1(a)).
const (
	MaxLine  = 1024
	MaxWords = 32
)

// The Split refusals, each named; the name leads the message.
var (
	ErrTooLong  = errors.New("TooLong: a line is at most 1024 bytes")
	ErrTooMany  = errors.New("TooMany: a line is at most 32 words")
	ErrNotText  = errors.New("NotText: a line is UTF-8 text with no control character")
	ErrUnclosed = errors.New("Unclosed: a quote is not closed")
	ErrShell    = errors.New("Shell: the bar runs no shell; quote $ ` ; | & < > ( ) to use them as text")
)

// shellChars are refused unquoted: the bar runs no shell, so what it shows is what a shell would run.
const shellChars = "$`;|&<>()"

// Split cuts one line into words: spaces separate them, '…' is literal, "…" is literal but for \"
// and \\, and adjacent segments join into one word.
func Split(line string) ([]string, error) {
	if len(line) > MaxLine {
		return nil, ErrTooLong
	}
	if !utf8.ValidString(line) {
		return nil, ErrNotText
	}
	for i := 0; i < len(line); i++ {
		if line[i] < 0x20 || line[i] == 0x7f {
			return nil, ErrNotText
		}
	}
	words := []string{}
	var w strings.Builder
	in := false
	for i := 0; i < len(line); i++ {
		c := line[i]
		switch {
		case c == ' ':
			if in {
				words = append(words, w.String())
				w.Reset()
				in = false
			}
			continue
		case c == '\'':
			j := strings.IndexByte(line[i+1:], '\'')
			if j < 0 {
				return nil, ErrUnclosed
			}
			w.WriteString(line[i+1 : i+1+j])
			i += j + 1
		case c == '"':
			j := i + 1
			for ; j < len(line) && line[j] != '"'; j++ {
				if line[j] == '\\' && j+1 < len(line) && (line[j+1] == '"' || line[j+1] == '\\') {
					j++
				}
				w.WriteByte(line[j])
			}
			if j == len(line) {
				return nil, ErrUnclosed
			}
			i = j
		case strings.IndexByte(shellChars, c) >= 0:
			return nil, ErrShell
		default:
			w.WriteByte(c)
		}
		in = true
	}
	if in {
		words = append(words, w.String())
	}
	if len(words) > MaxWords {
		return nil, ErrTooMany
	}
	return words, nil
}

// safeWord is a word Quote prints bare: nothing a shell or Split reads specially.
var safeWord = regexp.MustCompile(`^[A-Za-z0-9@%+=:,./_-]+$`)

// Quote prints argv as one line a shell and Split both read back as argv: a safe word bare, any
// other in single quotes, a ' as '"'"'.
func Quote(argv []string) string {
	out := make([]string, len(argv))
	for i, a := range argv {
		if safeWord.MatchString(a) {
			out[i] = a
			continue
		}
		out[i] = "'" + strings.ReplaceAll(a, "'", `'"'"'`) + "'"
	}
	return strings.Join(out, " ")
}
