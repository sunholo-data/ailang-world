package pkgproj

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

// ---------------------------------------------------------------------------
// The ready packet — the projected package's IDENTITY, as a Go value
//
// scripts/verify_world_package.sh step 9 already writes this document and
// compares it byte-for-byte against scripts/world_package_ready_packet.golden.
// json. That gate is a shell script and a python heredoc: nothing in Go could
// read the packet, so no Go caller could refuse to publish bytes that had
// drifted away from the reviewed artifact.
//
// This file is the Go half. It exists so cmd/world-publish's R-PACKET-DRIFT
// fence can RECOMPUTE the packet from the projection directory and compare it,
// field by field, with the committed golden — before any credential is loaded
// and before any request can leave the process.
//
// Everything here is pure and local: ContentHash, InterfaceHash, CreateTarball
// and TarballHash are all in-process, so recomputation needs no subprocess and
// no pinned binary. The one field that cannot be recomputed locally is
// compilerVersion, which is therefore a PARAMETER rather than something this
// package invents.
// ---------------------------------------------------------------------------

// ReadyPacket is the canonical ready-packet document.
//
// THE FIELD ORDER IS FROZEN and is exactly the key order the shell gate emits
// (python json.dumps(..., sort_keys=True), i.e. alphabetical). encoding/json
// emits struct fields in declaration order, so THIS DECLARATION IS THE WIRE
// ORDER, and EncodeReadyPacket reproduces the golden byte-for-byte.
// TestReadyPacketFieldOrderIsFrozen reads the key order back out of the encoded
// bytes rather than trusting this comment.
//
// compilerSHA256 is deliberately absent: it is provenance about the MACHINE,
// not about the package, and the shell gate keeps it out of the byte-compared
// golden for the same reason (a golden carrying it would pass on exactly one
// platform).
type ReadyPacket struct {
	CompilerVersion string   `json:"compilerVersion"`
	ContentHash     string   `json:"contentHash"`
	Effects         []string `json:"effects"`
	Exports         []string `json:"exports"`
	InterfaceHash   string   `json:"interfaceHash"`
	Package         string   `json:"package"`
	TarballBytes    int      `json:"tarballBytes"`
	TarballSHA256   string   `json:"tarballSHA256"`
	Version         string   `json:"version"`
}

// ReadyPacketFields is the frozen field order stated ONCE, independently of the
// struct, so the two can be compared. A single source would make the ordering
// test a tautology — the same argument publishPayloadFields is stated with in
// host/broker/registry_publish.go.
var ReadyPacketFields = []string{
	"compilerVersion", "contentHash", "effects", "exports", "interfaceHash",
	"package", "tarballBytes", "tarballSHA256", "version",
}

// EncodeReadyPacket renders the canonical document, INCLUDING the trailing
// newline the shell gate writes. It is byte-for-byte comparable with the
// committed golden; TestEncodeReadyPacketReproducesTheCommittedGoldenBytes is
// the measurement, not this sentence.
func EncodeReadyPacket(p ReadyPacket) []byte {
	p.Effects = normalizedList(p.Effects)
	p.Exports = normalizedList(p.Exports)
	encoded, err := json.Marshal(p)
	if err != nil {
		panic("pkgproj: fixed ready packet cannot fail JSON encoding: " + err.Error())
	}
	return append(encoded, '\n')
}

// normalizedList maps an empty list to a non-nil empty list so an absent and an
// empty set encode identically, and the document never carries JSON null where
// the gate emits [].
//
// The length test is deliberate and was MEASURED: the obvious `if items == nil`
// form is wrong, because `append([]string(nil), items...)` with zero elements
// returns a NIL slice again — so an already-decoded `"effects":[]` re-encoded as
// `"effects":null` and the byte-for-byte golden comparison failed. The empty
// case must be constructed, not appended.
func normalizedList(items []string) []string {
	if len(items) == 0 {
		return []string{}
	}
	return append([]string(nil), items...)
}

// DecodeReadyPacket parses a ready packet STRICTLY. Unknown fields are
// rejected and trailing bytes are rejected: a document this function cannot
// fully account for must not be used to authorize an irreversible write.
func DecodeReadyPacket(data []byte) (ReadyPacket, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var packet ReadyPacket
	if err := dec.Decode(&packet); err != nil {
		return ReadyPacket{}, fmt.Errorf("pkgproj: decode ready packet: %w", err)
	}
	// A second JSON value after the first would otherwise be silently ignored,
	// so a golden with a smuggled second document would decode as its first.
	if _, err := dec.Token(); err != io.EOF {
		return ReadyPacket{}, fmt.Errorf(
			"pkgproj: ready packet carries trailing content after the first JSON value")
	}
	return packet, nil
}

// LoadReadyPacket reads and strictly decodes the committed golden.
func LoadReadyPacket(path string) (ReadyPacket, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ReadyPacket{}, fmt.Errorf("pkgproj: read ready packet %s: %w", path, err)
	}
	return DecodeReadyPacket(data)
}

// RecomputeReadyPacket rebuilds the packet from the bytes on disk RIGHT NOW.
//
// Nothing carried in a manifest or a golden is trusted as a hash: content,
// interface and tarball digests are all computed here from the projection
// directory, exactly as RegistryPublishHandler.Execute recomputes them
// immediately before dispatch.
func RecomputeReadyPacket(dir string, manifest Manifest, compilerVersion string) (ReadyPacket, error) {
	content, err := ContentHash(dir)
	if err != nil {
		return ReadyPacket{}, fmt.Errorf("pkgproj: recompute content hash: %w", err)
	}
	tarball, err := CreateTarball(dir)
	if err != nil {
		return ReadyPacket{}, fmt.Errorf("pkgproj: recompute tarball: %w", err)
	}
	return ReadyPacket{
		CompilerVersion: compilerVersion,
		ContentHash:     content,
		Effects:         normalizedList(manifest.Effects.Max),
		Exports:         normalizedList(manifest.Exports.Modules),
		InterfaceHash:   InterfaceHash(manifest),
		Package:         manifest.Package.Name,
		TarballBytes:    len(tarball),
		TarballSHA256:   TarballHash(tarball),
		Version:         manifest.Package.Version,
	}, nil
}

// Equal reports whether two packets are identical, and NAMES THE FIRST FIELD
// that differs when they are not.
//
// The field name is the whole point. "the packet drifted" sends an attended
// operator to diff two 64-hex documents by eye at the one moment improvisation
// is most expensive; "contentHash differs" tells them a source file changed,
// and "tarballSHA256 differs" alone tells them the packaging changed while the
// sources did not — opposite remedies (design doc DD-1).
//
// It walks ReadyPacketFields, so a field added to the struct without being
// added to the frozen list is INVISIBLE here — which is why
// TestReadyPacketFieldOrderIsFrozen asserts the two agree.
func (p ReadyPacket) Equal(other ReadyPacket) (field string, equal bool) {
	for _, name := range ReadyPacketFields {
		mine, theirs := p.Field(name), other.Field(name)
		if mine != theirs {
			return name, false
		}
	}
	return "", true
}

// Field renders one field as canonical text. Comparing rendered text rather
// than reflecting over the struct keeps the comparison total: a slice compared
// with == would not compile, and a slice compared by length alone would call
// ["a","b"] equal to ["b","a"].
func (p ReadyPacket) Field(name string) string {
	switch name {
	case "compilerVersion":
		return p.CompilerVersion
	case "contentHash":
		return p.ContentHash
	case "effects":
		return renderList(p.Effects)
	case "exports":
		return renderList(p.Exports)
	case "interfaceHash":
		return p.InterfaceHash
	case "package":
		return p.Package
	case "tarballBytes":
		return strconv.Itoa(p.TarballBytes)
	case "tarballSHA256":
		return p.TarballSHA256
	case "version":
		return p.Version
	default:
		// An unknown name can only arrive from ReadyPacketFields, so this is a
		// programming error in THIS file, not a bad input. Panicking is right:
		// silently returning "" would make two different packets compare equal
		// on the field nobody can read.
		panic("pkgproj: ready packet has no field named " + strconv.Quote(name))
	}
}

// renderList is order-sensitive and length-sensitive by construction: it
// encodes the slice as JSON rather than joining it, so neither reordering nor a
// separator smuggled into an element can make two different lists render the
// same.
func renderList(items []string) string {
	encoded, err := json.Marshal(normalizedList(items))
	if err != nil {
		panic("pkgproj: fixed string list cannot fail JSON encoding: " + err.Error())
	}
	return string(encoded)
}
