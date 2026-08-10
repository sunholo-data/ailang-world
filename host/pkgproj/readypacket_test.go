package pkgproj

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// goldenPath is the ONE committed artifact this file is stated against. Every
// assertion below reads it; none re-derives it.
const goldenRelPath = "scripts/world_package_ready_packet.golden.json"

func readyPacketRepoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate repository root")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func goldenBytes(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(readyPacketRepoRoot(t), goldenRelPath))
	if err != nil {
		t.Fatalf("committed ready-packet golden is unreadable at %s: %v", goldenRelPath, err)
	}
	// ANTI-VACUITY: an empty golden would make every round-trip below trivially
	// consistent with itself.
	if len(data) == 0 {
		t.Fatalf("instrument failure: %s is zero bytes", goldenRelPath)
	}
	return data
}

// worldCoreGolden is the decoded committed packet, with a same-call assertion
// that it is populated. A zero-valued packet would satisfy several assertions
// below for the wrong reason.
func worldCoreGolden(t *testing.T) ReadyPacket {
	t.Helper()
	packet, err := DecodeReadyPacket(goldenBytes(t))
	if err != nil {
		t.Fatalf("decode committed golden: %v", err)
	}
	for _, arm := range []struct{ name, value string }{
		{"package", packet.Package},
		{"version", packet.Version},
		{"contentHash", packet.ContentHash},
		{"interfaceHash", packet.InterfaceHash},
		{"tarballSHA256", packet.TarballSHA256},
		{"compilerVersion", packet.CompilerVersion},
	} {
		if arm.value == "" {
			t.Fatalf("instrument failure: the committed golden decoded with an empty %s", arm.name)
		}
	}
	if len(packet.Exports) == 0 {
		t.Fatal("instrument failure: the committed golden decoded with zero exports")
	}
	if packet.TarballBytes <= 0 {
		t.Fatalf("instrument failure: the committed golden decoded tarballBytes=%d", packet.TarballBytes)
	}
	return packet
}

// TestEncodeReadyPacketReproducesTheCommittedGoldenBytes is the binding between
// this Go codec and the shell gate's python heredoc. They are two independent
// implementations of the same canonical document; if they ever disagree, the
// R-PACKET-DRIFT fence would refuse a packet the readiness gate calls correct
// (or, far worse, accept one it does not).
func TestEncodeReadyPacketReproducesTheCommittedGoldenBytes(t *testing.T) {
	want := goldenBytes(t)
	got := EncodeReadyPacket(worldCoreGolden(t))
	if string(got) != string(want) {
		t.Fatalf("re-encoded golden is not byte-identical:\n got %q\nwant %q", got, want)
	}
	// The trailing newline is part of the shell gate's `cmp -s` comparison, so
	// it is part of the contract rather than cosmetic.
	if got[len(got)-1] != '\n' {
		t.Fatal("canonical encoding does not end in the newline the gate writes")
	}
	t.Logf("canonical encoding reproduced %d bytes of %s byte-for-byte", len(got), goldenRelPath)
}

// TestReadyPacketFieldOrderIsFrozen reads the key order back OUT of the encoded
// bytes. A struct field added without a place in ReadyPacketFields would be
// invisible to Equal — i.e. a drifted field that no fence could see.
func TestReadyPacketFieldOrderIsFrozen(t *testing.T) {
	encoded := EncodeReadyPacket(worldCoreGolden(t))
	dec := json.NewDecoder(strings.NewReader(string(encoded)))
	if tok, err := dec.Token(); err != nil || tok != json.Delim('{') {
		t.Fatalf("encoded packet does not begin with an object: %v %v", tok, err)
	}
	var keys []string
	for dec.More() {
		key, err := dec.Token()
		if err != nil {
			t.Fatal(err)
		}
		name, ok := key.(string)
		if !ok {
			t.Fatalf("object key %v is not a string", key)
		}
		keys = append(keys, name)
		var discard json.RawMessage
		if err := dec.Decode(&discard); err != nil {
			t.Fatal(err)
		}
	}
	if len(keys) == 0 {
		t.Fatal("instrument failure: read ZERO keys out of the encoded packet")
	}
	if strings.Join(keys, ",") != strings.Join(ReadyPacketFields, ",") {
		t.Fatalf("wire key order = %v, want the frozen order %v", keys, ReadyPacketFields)
	}
}

// TestReadyPacketEqualNamesTheFirstDifferingField drives ONE mutation per
// frozen field and requires Equal to name that exact field. A comparator that
// returned "not equal" without the field name would pass a weaker test; a
// comparator that only looked at, say, the three digests would pass a test that
// mutated only digests.
func TestReadyPacketEqualNamesTheFirstDifferingField(t *testing.T) {
	base := worldCoreGolden(t)

	// POSITIVE CONTROL, first: the comparator must be able to say EQUAL. Without
	// it, a comparator that reported every pair as different would pass every
	// arm below.
	if field, equal := base.Equal(base); !equal {
		t.Fatalf("instrument failure: a packet does not equal itself (first difference %q)", field)
	}

	mutations := map[string]func(ReadyPacket) ReadyPacket{
		"compilerVersion": func(p ReadyPacket) ReadyPacket { p.CompilerVersion += "-x"; return p },
		"contentHash":     func(p ReadyPacket) ReadyPacket { p.ContentHash = flipHexNibble(p.ContentHash); return p },
		"effects":         func(p ReadyPacket) ReadyPacket { p.Effects = []string{"IO"}; return p },
		"exports": func(p ReadyPacket) ReadyPacket {
			// REORDER only. Same length, same elements: a comparator that
			// compared sorted sets or lengths would miss this, and the interface
			// hash genuinely depends on the export set.
			p.Exports = append([]string(nil), p.Exports...)
			p.Exports[0], p.Exports[len(p.Exports)-1] = p.Exports[len(p.Exports)-1], p.Exports[0]
			return p
		},
		"interfaceHash": func(p ReadyPacket) ReadyPacket { p.InterfaceHash = flipHexNibble(p.InterfaceHash); return p },
		"package":       func(p ReadyPacket) ReadyPacket { p.Package = "world/other"; return p },
		"tarballBytes":  func(p ReadyPacket) ReadyPacket { p.TarballBytes++; return p },
		"tarballSHA256": func(p ReadyPacket) ReadyPacket { p.TarballSHA256 = flipHexNibble(p.TarballSHA256); return p },
		"version":       func(p ReadyPacket) ReadyPacket { p.Version = "0.1.1"; return p },
	}
	if len(mutations) != len(ReadyPacketFields) {
		t.Fatalf("this test drives %d fields but the frozen list has %d: a field would go untested",
			len(mutations), len(ReadyPacketFields))
	}

	for _, name := range ReadyPacketFields {
		mutate, ok := mutations[name]
		if !ok {
			t.Fatalf("frozen field %q has no mutation arm", name)
		}
		t.Run(name, func(t *testing.T) {
			mutant := mutate(base)
			field, equal := base.Equal(mutant)
			if equal {
				t.Fatalf("Equal reported EQUAL after mutating %s", name)
			}
			if field != name {
				t.Fatalf("Equal named %q as the first difference, want %q", field, name)
			}
			// Symmetry: an asymmetric comparator would let drift in one
			// direction pass depending on which side the fence put first.
			if reverse, equal := mutant.Equal(base); equal || reverse != name {
				t.Fatalf("reversed Equal named %q (equal=%v), want %q", reverse, equal, name)
			}
		})
	}
}

// TestDecodeReadyPacketIsStrict pins the two refusals that keep a golden from
// being partially believed.
func TestDecodeReadyPacketIsStrict(t *testing.T) {
	valid := EncodeReadyPacket(worldCoreGolden(t))

	// CONTROL: the decoder accepts the real document, so the refusals below are
	// refusals of the mutation and not of everything.
	if _, err := DecodeReadyPacket(valid); err != nil {
		t.Fatalf("instrument failure: the decoder rejects the committed golden: %v", err)
	}

	for _, arm := range []struct{ name, body string }{
		{"unknown field", strings.Replace(string(valid), `{"compilerVersion"`, `{"compilerSHA256":"x","compilerVersion"`, 1)},
		{"trailing document", strings.TrimRight(string(valid), "\n") + "\n{}\n"},
		{"truncated", string(valid[:len(valid)/2])},
	} {
		if _, err := DecodeReadyPacket([]byte(arm.body)); err == nil {
			t.Errorf("decoder accepted a %s packet", arm.name)
		}
	}
}

// TestLoadReadyPacketReportsAMissingFile keeps the fence's "cannot read the
// golden" path distinguishable from "the golden says something else".
func TestLoadReadyPacketReportsAMissingFile(t *testing.T) {
	absent := filepath.Join(t.TempDir(), "no-such-packet.json")
	if _, err := os.Stat(absent); err == nil {
		t.Fatal("negative control failed: the 'absent' path exists")
	}
	if _, err := LoadReadyPacket(absent); err == nil {
		t.Fatal("LoadReadyPacket accepted a path that does not exist")
	}
	// And the positive half in the same test: the real golden loads.
	loaded, err := LoadReadyPacket(filepath.Join(readyPacketRepoRoot(t), goldenRelPath))
	if err != nil {
		t.Fatalf("LoadReadyPacket(%s): %v", goldenRelPath, err)
	}
	if field, equal := loaded.Equal(worldCoreGolden(t)); !equal {
		t.Fatalf("LoadReadyPacket disagrees with DecodeReadyPacket at %q", field)
	}
}

// TestInterfaceHashIgnoresTheVersion measures the claim made in Package's doc
// comment. If a future edit fed Version into InterfaceHash, every approval
// minted against the current interface hash would silently stop matching.
func TestInterfaceHashIgnoresTheVersion(t *testing.T) {
	base := Manifest{
		Package: Package{Name: "world/core", Edition: "1", AILANG: ">=0.30.0", Version: "0.1.0"},
		Exports: Exports{Modules: []string{"world/types"}},
		Effects: Effects{Max: []string{}},
	}
	bumped := base
	bumped.Package.Version = "9.9.9"
	if InterfaceHash(base) != InterfaceHash(bumped) {
		t.Fatal("InterfaceHash changed when only the version changed")
	}
	// CONTROL in the same test: the hash IS sensitive to something, so the
	// equality above is not the equality of a constant function.
	changed := base
	changed.Exports.Modules = []string{"world/types", "world/other"}
	if InterfaceHash(base) == InterfaceHash(changed) {
		t.Fatal("instrument failure: InterfaceHash is insensitive to the export set")
	}
}

// flipHexNibble changes exactly one hex digit of a digest, leaving its length
// and its "sha256:" prefix alone.
func flipHexNibble(hash string) string {
	if hash == "" {
		panic("flipHexNibble on an empty hash")
	}
	last := hash[len(hash)-1]
	replacement := byte('0')
	if last == '0' {
		replacement = '1'
	}
	return hash[:len(hash)-1] + string(replacement)
}
