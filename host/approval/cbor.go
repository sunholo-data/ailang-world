package approval

// Row 201 A2 (design_docs/planned/w-world-approvals-by-key-a2.md §2.3): the bounded CBOR subset a
// WebAuthn registration needs (RFC 8949 major types 0-5 and the simple values false/true/null),
// read from a byte slice with fixed bounds. Every refusal is its own branch (K1-K16). No I/O.

import (
	"encoding/binary"
	"errors"
	"math"
	"unicode/utf8"
)

var (
	ErrCBORTruncated    = errors.New("approval: cbor: truncated")
	ErrCBORMalformed    = errors.New("approval: cbor: malformed")
	ErrCBORIndefinite   = errors.New("approval: cbor: indefinite length is not accepted")
	ErrCBORNonCanonical = errors.New("approval: cbor: an argument is not in its shortest form")
	ErrCBORTooLarge     = errors.New("approval: cbor: over a bound")
	ErrCBORTooDeep      = errors.New("approval: cbor: nested too deep")
	ErrCBORUnsupported  = errors.New("approval: cbor: a tag, float, simple value or map key type outside the subset")
	ErrCBORDuplicateKey = errors.New("approval: cbor: a duplicate map key")
	ErrCBORTrailing     = errors.New("approval: cbor: trailing bytes")
)

type cborKind uint8

const (
	cborInt cborKind = iota + 1
	cborBytes
	cborText
	cborArray
	cborMap
	cborBool
	cborNull
)

// cborValue is one decoded item. A map keeps its entries in the order they were read.
type cborValue struct {
	kind  cborKind
	i     int64
	b     []byte // bytes, or text's UTF-8
	t     bool
	items []cborValue // array elements, or map keys and values alternating
}

// get returns the value a map holds under an integer or text key.
func (v cborValue) get(key cborValue) (cborValue, bool) {
	for j := 0; v.kind == cborMap && j < len(v.items); j += 2 {
		if sameKey(v.items[j], key) {
			return v.items[j+1], true
		}
	}
	return cborValue{}, false
}

func sameKey(a, b cborValue) bool {
	return a.kind == b.kind && a.i == b.i && string(a.b) == string(b.b)
}

type cborLimits struct{ maxDepth, maxItemsPerContainer, maxTotalItems, maxStringBytes int }

var attestationLimits = cborLimits{maxDepth: 4, maxItemsPerContainer: 16, maxTotalItems: 64, maxStringBytes: 2048}

// decodeCBOR reads one item and returns it and the bytes after it.
func decodeCBOR(b []byte, lim cborLimits) (cborValue, []byte, error) {
	d := &cborDecoder{b: b, lim: lim}
	v, err := d.item(1)
	if err != nil {
		return cborValue{}, nil, err
	}
	return v, d.b, nil
}

// decodeCBORExact reads one item and refuses trailing bytes.
func decodeCBORExact(b []byte, lim cborLimits) (cborValue, error) {
	v, rest, err := decodeCBOR(b, lim)
	if err != nil {
		return cborValue{}, err
	}
	if len(rest) != 0 { // K16
		return cborValue{}, ErrCBORTrailing
	}
	return v, nil
}

type cborDecoder struct {
	b     []byte
	lim   cborLimits
	total int
}

// head reads an initial byte and its argument.
func (d *cborDecoder) head() (major byte, arg uint64, err error) {
	if len(d.b) == 0 { // K1
		return 0, 0, ErrCBORTruncated
	}
	major, ai := d.b[0]>>5, d.b[0]&0x1f
	switch {
	case major == 6: // K9
		return 0, 0, ErrCBORUnsupported
	case major == 7:
		if ai < 20 || ai > 22 { // K10
			return 0, 0, ErrCBORUnsupported
		}
		d.b = d.b[1:]
		return major, uint64(ai), nil
	case ai < 24:
		d.b = d.b[1:]
		return major, uint64(ai), nil
	case ai == 31: // K3
		return 0, 0, ErrCBORIndefinite
	case ai > 27: // K2
		return 0, 0, ErrCBORMalformed
	}
	n := 1 << (ai - 24)
	if len(d.b) < 1+n { // K1
		return 0, 0, ErrCBORTruncated
	}
	var buf [8]byte
	copy(buf[8-n:], d.b[1:1+n])
	arg = binary.BigEndian.Uint64(buf[:])
	if arg < [4]uint64{24, 1 << 8, 1 << 16, 1 << 32}[ai-24] { // K4
		return 0, 0, ErrCBORNonCanonical
	}
	d.b = d.b[1+n:]
	return major, arg, nil
}

func (d *cborDecoder) item(depth int) (cborValue, error) {
	if depth > d.lim.maxDepth { // K8
		return cborValue{}, ErrCBORTooDeep
	}
	if d.total++; d.total > d.lim.maxTotalItems { // K15
		return cborValue{}, ErrCBORTooLarge
	}
	major, arg, err := d.head()
	if err != nil {
		return cborValue{}, err
	}
	switch major {
	case 0, 1:
		if arg > math.MaxInt64 { // K14
			return cborValue{}, ErrCBORTooLarge
		}
		if major == 1 {
			return cborValue{kind: cborInt, i: -1 - int64(arg)}, nil
		}
		return cborValue{kind: cborInt, i: int64(arg)}, nil
	case 2, 3:
		if arg > uint64(d.lim.maxStringBytes) { // K6
			return cborValue{}, ErrCBORTooLarge
		}
		if arg > uint64(len(d.b)) { // K5
			return cborValue{}, ErrCBORTruncated
		}
		s := d.b[:arg:arg]
		d.b = d.b[arg:]
		if major == 3 && !utf8.Valid(s) { // K11
			return cborValue{}, ErrCBORMalformed
		}
		return cborValue{kind: cborBytes + cborKind(major-2), b: s}, nil
	case 4, 5:
		if arg > uint64(d.lim.maxItemsPerContainer) { // K7
			return cborValue{}, ErrCBORTooLarge
		}
		v := cborValue{kind: cborArray + cborKind(major-4)}
		for range int(arg) {
			e, err := d.item(depth + 1)
			if err != nil {
				return cborValue{}, err
			}
			if major == 5 {
				if e.kind != cborInt && e.kind != cborText { // K12
					return cborValue{}, ErrCBORUnsupported
				}
				if _, dup := v.get(e); dup { // K13
					return cborValue{}, ErrCBORDuplicateKey
				}
				val, err := d.item(depth + 1)
				if err != nil {
					return cborValue{}, err
				}
				v.items = append(v.items, e, val)
				continue
			}
			v.items = append(v.items, e)
		}
		return v, nil
	}
	if arg == 22 {
		return cborValue{kind: cborNull}, nil
	}
	return cborValue{kind: cborBool, t: arg == 21}, nil
}
