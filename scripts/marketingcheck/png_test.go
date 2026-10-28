package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func wantCode(t *testing.T, err error, code, assertion string) {
	t.Helper()
	var r *Refusal
	if !errors.As(err, &r) || r.Code != code {
		t.Fatalf("%s: want Code=%s, got %v", assertion, code, err)
	}
}
func pngOK(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
func encodePNG(t *testing.T, im image.Image) []byte {
	t.Helper()
	var b bytes.Buffer
	pngOK(t, (&png.Encoder{}).Encode(&b, im))
	return b.Bytes()
}
func putPNG(t *testing.T, dir, name string, b []byte) {
	t.Helper()
	pngOK(t, os.WriteFile(filepath.Join(dir, name), b, 0600))
}
func putJSON(t *testing.T, dir, name string, v any) string {
	t.Helper()
	b, e := json.Marshal(v)
	pngOK(t, e)
	p := filepath.Join(dir, name)
	pngOK(t, os.WriteFile(p, b, 0600))
	return p
}

type testChunk struct {
	typ  string
	data []byte
}

func chunks(t *testing.T, b []byte) []testChunk {
	t.Helper()
	var cs []testChunk
	for i := 8; i < len(b); {
		n := int(binary.BigEndian.Uint32(b[i : i+4]))
		cs = append(cs, testChunk{string(b[i+4 : i+8]), append([]byte(nil), b[i+8:i+8+n]...)})
		i += 12 + n
	}
	return cs
}
func chunkStream(cs []testChunk) []byte {
	b := []byte(PNGSignature)
	for _, c := range cs {
		start := len(b)
		b = binary.BigEndian.AppendUint32(b, uint32(len(c.data)))
		b = append(b, []byte(c.typ)...)
		b = append(b, c.data...)
		b = binary.BigEndian.AppendUint32(b, crc32.ChecksumIEEE(b[start+4:]))
	}
	return b
}
func addChunk(t *testing.T, b []byte, typ string, data []byte) []byte {
	cs := chunks(t, b)
	return chunkStream(append([]testChunk{cs[0], {typ, data}}, cs[1:]...))
}

// Legal optional chunks cover every size residue modulo the 12-byte chunk overhead.
// Contiguous empty IDAT chunks then reach the exact cap without invalid metadata.
func exactSizedPNG(t *testing.T, b []byte, size int) []byte {
	t.Helper()
	cs := chunks(t, b)
	var chroma []byte
	for _, v := range []uint32{31270, 32900, 64000, 33000, 30000, 60000, 15000, 6000} {
		chroma = binary.BigEndian.AppendUint32(chroma, v)
	}
	options := []testChunk{{"sRGB", []byte{0}}, {"gAMA", []byte{0, 0, 177, 143}}, {"pHYs", []byte{0, 0, 14, 196, 0, 0, 14, 196, 1}}, {"cHRM", chroma}, {"PLTE", []byte{255, 255, 255}}}
	for mask := 0; mask < 1<<len(options); mask++ {
		extra := []testChunk{}
		cost := 0
		for i, c := range options {
			if mask&(1<<i) != 0 {
				extra = append(extra, c)
				cost += 12 + len(c.data)
			}
		}
		gap := size - len(b) - cost
		if gap < 0 || gap%12 != 0 {
			continue
		}
		out := append([]testChunk{cs[0]}, extra...)
		for i := 0; i < gap/12; i++ {
			out = append(out, testChunk{"IDAT", nil})
		}
		out = append(out, cs[1:]...)
		result := chunkStream(out)
		if len(result) != size {
			t.Fatal("exact-size fixture length")
		}
		return result
	}
	t.Fatal("no legal padding combination")
	return nil
}
func noiseFixture(t *testing.T) ([]byte, ShotRects) {
	t.Helper()
	_, r := fixture(t, Shots[5])
	im := image.NewNRGBA(image.Rect(0, 0, 390, 844))
	state := uint32(0x1a2b3c4d)
	for y := 0; y < 844; y++ {
		for x := 0; x < 390; x++ {
			state ^= state << 13
			state ^= state >> 17
			state ^= state << 5
			col := color.NRGBA{uint8(state), uint8(state >> 8), uint8(state >> 16), 255}
			if x < 64 && y < 16 {
				col = color.NRGBA{240, 240, 240, 255}
			}
			im.SetNRGBA(x, y, col)
		}
	}
	b := encodePNG(t, im)
	pngOK(t, CheckContent("phone.png", b, r))
	return b, r
}

func fixture(t *testing.T, s Shot) ([]byte, ShotRects) {
	t.Helper()
	im := image.NewNRGBA(image.Rect(0, 0, s.Width, s.Height))
	bg := uint8(240)
	if s.Scheme == "dark" {
		bg = 24
	}
	for i := 0; i < len(im.Pix); i += 4 {
		im.Pix[i] = bg
		im.Pix[i+1] = bg
		im.Pix[i+2] = bg
		im.Pix[i+3] = 255
	}
	r := ShotRects{Background: Rect{X: 0, Y: 0, W: 64, H: 16}}
	names := []string{"header"}
	if s.Name == "packet.png" {
		names = []string{"validation", "roster-proven", "roster-tested", "roster-attested", "roster-claimed"}
	}
	for i, name := range names {
		c := Rect{Name: name, X: 32, Y: 64 + 160*i, W: min(s.Width-64, 900), H: 128}
		r.Content = append(r.Content, c)
		for y := c.Y; y < c.Y+c.H; y++ {
			for x := c.X; x < c.X+c.W; x++ {
				v := uint8((x+y)%200 + 20)
				col := color.NRGBA{v, 30, uint8(255 - v/2), 255}
				if y%8 < 3 {
					col = color.NRGBA{240, v, 20, 255}
				}
				im.SetNRGBA(x, y, col)
			}
		}
	}
	return encodePNG(t, im), r
}
func fixtureSet(t *testing.T) (map[string][]byte, Rectangles) {
	t.Helper()
	bs := map[string][]byte{}
	r := Rectangles{Schema: "world/demo/capture-rectangles/v1", Shots: map[string]ShotRects{}}
	for _, s := range Shots {
		b, sr := fixture(t, s)
		bs[s.Name] = b
		r.Shots[s.Name] = sr
	}
	return bs, r
}
func writeSet(t *testing.T, bs map[string][]byte) string {
	t.Helper()
	dir := t.TempDir()
	for n, b := range bs {
		putPNG(t, dir, n, b)
	}
	return dir
}
func demoFixture() DemoManifest {
	return DemoManifest{Schema: "world/demo/public-manifest/v1", Version: 1, References: []DemoReference{{Role: "frozen-head", Ref: "head-ref"}, {Role: "packet", Ref: "packet-ref"}, {Role: "invocation", Ref: "other-invocation", Index: 8}, {Role: "invocation", Ref: "approval-invocation", Index: 19}}, Why: []DemoWhy{{Role: "read", Entry: 8}, {Role: "approval", Entry: 19}}}
}

func TestMarketingPNGStructure(t *testing.T) {
	const fixtureFileCap = 1_048_576
	im := image.NewNRGBA(image.Rect(0, 0, 32, 32))
	b := encodePNG(t, im)
	t.Run("positive", func(t *testing.T) {
		pngOK(t, CheckStructure("control.png", b))
		for _, c := range chunks(t, b) {
			if c.typ != "IHDR" && c.typ != "IDAT" && c.typ != "IEND" {
				t.Fatalf("encoder emitted ancillary %s", c.typ)
			}
		}
	})
	for _, n := range []string{"../a.png", "A.png", "a_b.png", ".png", "a.png\n"} {
		t.Run("bad-name-"+n, func(t *testing.T) { wantCode(t, CheckStructure(n, b), "bad-name", "bad name rejected") })
	}
	t.Run("bad-signature", func(t *testing.T) {
		bad := append([]byte(nil), b...)
		bad[0] = 0
		wantCode(t, CheckStructure("a.png", bad), "bad-signature", "signature rejected")
	})
	t.Run("ihdr-not-first", func(t *testing.T) {
		cs := chunks(t, b)
		cs[0], cs[1] = cs[1], cs[0]
		wantCode(t, CheckStructure("a.png", chunkStream(cs)), "ihdr-not-first", "IHDR order rejected")
	})
	t.Run("ihdr-length", func(t *testing.T) {
		cs := chunks(t, b)
		cs[0].data = cs[0].data[:12]
		wantCode(t, CheckStructure("a.png", chunkStream(cs)), "ihdr-length", "IHDR length rejected")
	})
	t.Run("bad-crc", func(t *testing.T) {
		bad := append([]byte(nil), b...)
		bad[29] ^= 1
		wantCode(t, CheckStructure("a.png", bad), "bad-crc", "bad-CRC typed refusal")
	})
	for _, typ := range []string{"tEXt", "zTXt", "iTXt", "eXIf"} {
		t.Run("forbidden-"+typ, func(t *testing.T) {
			wantCode(t, CheckStructure("a.png", addChunk(t, b, typ, []byte("note"))), "forbidden-chunk", "forbidden "+typ+" rejected")
		})
	}
	t.Run("unknown", func(t *testing.T) {
		wantCode(t, CheckStructure("a.png", addChunk(t, b, "vpAg", nil)), "unknown-chunk", "unknown chunk rejected")
	})
	t.Run("allowed-ancillaries", func(t *testing.T) {
		for _, c := range []testChunk{{"sRGB", []byte{0}}, {"gAMA", []byte{0, 0, 177, 143}}, {"pHYs", make([]byte, 9)}} {
			b = addChunk(t, b, c.typ, c.data)
		}
		pngOK(t, CheckStructure("a.png", b))
	})
	t.Run("multiple-iend", func(t *testing.T) {
		cs := chunks(t, b)
		cs = append(cs, testChunk{"IEND", nil})
		wantCode(t, CheckStructure("a.png", chunkStream(cs)), "multiple-iend", "two IEND rejected")
	})
	t.Run("trailing-byte", func(t *testing.T) {
		wantCode(t, CheckStructure("a.png", append(append([]byte(nil), b...), 0)), "trailing-bytes", "one trailing byte rejected")
	})
	t.Run("missing-iend", func(t *testing.T) {
		wantCode(t, CheckStructure("a.png", b[:len(b)-12]), "missing-iend", "missing IEND rejected")
	})
	t.Run("truncated-header", func(t *testing.T) {
		wantCode(t, CheckStructure("a.png", b[:10]), "truncated-chunk", "truncated header rejected")
	})
	t.Run("truncated-data", func(t *testing.T) {
		wantCode(t, CheckStructure("a.png", b[:25]), "truncated-chunk", "truncated data rejected")
	})
	// Noise plus legal chunk boundaries exercises exact byte caps with valid decoded content.
	capBase, capRects := noiseFixture(t)
	capBytes := exactSizedPNG(t, capBase, fixtureFileCap)
	t.Run("at-cap", func(t *testing.T) {
		if len(capBytes) != fixtureFileCap {
			t.Fatal("fixture size")
		}
		pngOK(t, CheckStructure("phone.png", capBytes))
		pngOK(t, CheckContent("phone.png", capBytes, capRects))
		_, e := png.Decode(bytes.NewReader(capBytes))
		pngOK(t, e)
	})
	t.Run("cap-plus-one", func(t *testing.T) {
		over := exactSizedPNG(t, capBase, fixtureFileCap+1)
		if len(over) != fixtureFileCap+1 {
			t.Fatal("fixture size")
		}
		_, decodeErr := png.Decode(bytes.NewReader(over))
		pngOK(t, decodeErr)
		pngOK(t, CheckContent("phone.png", over, capRects))
		wantCode(t, CheckStructure("phone.png", over), "file-cap-exceeded", "cap+1 rejected")
	})
	t.Run("directory-unreadable", func(t *testing.T) {
		wantCode(t, CheckDirTotal(filepath.Join(t.TempDir(), "absent")), "dir-cap-exceeded", "unreadable directory rejected")
	})
	t.Run("directory-cap", func(t *testing.T) {
		dir := t.TempDir()
		under := exactSizedPNG(t, capBase, fixtureFileCap-1)
		pngOK(t, CheckStructure("a.png", under))
		for i := 0; i < 8; i++ {
			putPNG(t, dir, fmt.Sprintf("n%d.png", i), under)
		}
		putPNG(t, dir, "ignored.txt", capBytes)
		pngOK(t, CheckDirTotal(dir))
		putPNG(t, dir, "n8.png", under)
		wantCode(t, CheckDirTotal(dir), "dir-cap-exceeded", "aggregate > DirCap rejected")
	})
}

func TestMarketingPNGContent(t *testing.T) {
	bs, r := fixtureSet(t)
	name := "workbench-light.png"
	sr := r.Shots[name]
	for _, s := range Shots {
		t.Run("positive-"+s.Name, func(t *testing.T) { pngOK(t, CheckContent(s.Name, bs[s.Name], r.Shots[s.Name])) })
	}
	t.Run("decode", func(t *testing.T) {
		wantCode(t, CheckContent(name, []byte("bad"), sr), "decode-failed", "decoder failure typed")
	})
	t.Run("unknown-name", func(t *testing.T) {
		wantCode(t, CheckContent("other.png", bs[name], sr), "bad-name", "unknown content name rejected")
	})
	t.Run("dimensions", func(t *testing.T) {
		b := encodePNG(t, image.NewNRGBA(image.Rect(0, 0, 1439, 1100)))
		wantCode(t, CheckContent(name, b, sr), "wrong-dimensions", "one pixel dimension error rejected")
	})
	t.Run("uniform", func(t *testing.T) {
		im := image.NewNRGBA(image.Rect(0, 0, 1440, 1100))
		for y := 0; y < 1100; y++ {
			for x := 0; x < 1440; x++ {
				im.SetNRGBA(x, y, color.NRGBA{240, 240, 210, 255})
			}
		}
		b := encodePNG(t, im)
		wantCode(t, CheckColors(name, im), "too-few-colors", "uniform PNG rejected by color guard")
		wantCode(t, CheckContent(name, b, sr), "too-few-colors", "uniform PNG rejected")
	})
	t.Run("mostly-gray", func(t *testing.T) {
		im := image.NewNRGBA(image.Rect(0, 0, 1440, 1100))
		for y := 0; y < 1100; y++ {
			for x := 0; x < 1440; x++ {
				v := uint8(160 + x%96)
				im.SetNRGBA(x, y, color.NRGBA{v, v, v, 255})
			}
		}
		im.SetNRGBA(0, 0, color.NRGBA{255, 0, 0, 255})
		wantCode(t, CheckContent(name, encodePNG(t, im), sr), "not-chromatic", "mostly-gray PNG rejected independently")
	})
	t.Run("light-scheme", func(t *testing.T) {
		wantCode(t, CheckContent(name, bs["workbench-dark.png"], sr), "wrong-scheme", "light rejects dark background")
	})
	t.Run("dark-scheme", func(t *testing.T) {
		wantCode(t, CheckContent("workbench-dark.png", bs[name], sr), "wrong-scheme", "dark rejects light background")
	})
	for _, tc := range []struct {
		name, code string
		rect       Rect
	}{
		{"blank", "blank-region", Rect{Name: "flat", X: 0, Y: 0, W: 64, H: 16}},
		{"bounds", "rect-out-of-bounds", Rect{Name: "outside", X: 1400, Y: 64, W: 64, H: 16}},
		{"negative", "rect-out-of-bounds", Rect{Name: "outside", X: -1, Y: 64, W: 64, H: 16}},
		{"tiny", "rect-too-small", Rect{Name: "tiny", X: 32, Y: 64, W: 63, H: 16}},
		{"short", "rect-too-small", Rect{Name: "short", X: 32, Y: 64, W: 64, H: 15}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rr := sr
			rr.Content = []Rect{tc.rect}
			if tc.code == "rect-out-of-bounds" || tc.code == "rect-too-small" {
				wantCode(t, CheckRectBounds(name, tc.rect, image.Rect(0, 0, 1440, 1100), true), tc.code, tc.name+" rectangle guard rejected")
			}
			wantCode(t, CheckContent(name, bs[name], rr), tc.code, tc.name+" rectangle rejected")
		})
	}
	t.Run("low-variance", func(t *testing.T) {
		im := image.NewNRGBA(image.Rect(0, 0, 128, 32))
		for y := 0; y < 32; y++ {
			for x := 0; x < 128; x++ {
				v := uint8(240 + x%2)
				im.SetNRGBA(x, y, color.NRGBA{v, v, v, 255})
			}
		}
		wantCode(t, CheckRegion(name, im, Rect{Name: "subtle", W: 128, H: 32}), "blank-region", "two-color low variance region rejected")
	})
	t.Run("bounds-branches", func(t *testing.T) {
		for _, rr := range []Rect{{X: 0, Y: -1, W: 64, H: 16}, {X: 0, Y: 0, W: 0, H: 16}, {X: 0, Y: 0, W: 64, H: 0}, {X: 1441, Y: 0, W: 64, H: 16}, {X: 0, Y: 1101, W: 64, H: 16}, {X: 0, Y: 1090, W: 64, H: 16}, {X: 32, Y: 64, W: int(^uint(0) >> 1), H: 16}} {
			wantCode(t, CheckRectBounds(name, rr, image.Rect(0, 0, 1440, 1100), true), "rect-out-of-bounds", "rectangle bounds branch rejected")
		}
	})
	t.Run("color-boundaries", func(t *testing.T) {
		for _, count := range []int{32, 33} {
			im := image.NewNRGBA(image.Rect(0, 0, count*2, 16))
			for y := 0; y < 16; y++ {
				for x := 0; x < count*2; x++ {
					im.SetNRGBA(x, y, color.NRGBA{uint8(100 + x%count), 20, 200, 255})
				}
			}
			err := CheckColors(name, im)
			if count == 32 {
				wantCode(t, err, "too-few-colors", "32 RGB colors rejected")
			} else {
				pngOK(t, err)
			}
		}
	})
	t.Run("chromatic-boundaries", func(t *testing.T) {
		for _, n := range []int{99, 100} {
			im := image.NewNRGBA(image.Rect(0, 0, 100, 100))
			for y := 0; y < 100; y++ {
				for x := 0; x < 100; x++ {
					v := uint8(x)
					im.SetNRGBA(x, y, color.NRGBA{v, v, v, 255})
				}
			}
			for x := 0; x < n; x++ {
				im.SetNRGBA(x, 0, color.NRGBA{255, 0, 0, 255})
			}
			err := CheckColors(name, im)
			if n == 99 {
				wantCode(t, err, "not-chromatic", "less than one percent rejected")
			} else {
				pngOK(t, err)
			}
		}
		for _, delta := range []int{12, 13} {
			im := image.NewNRGBA(image.Rect(0, 0, 100, 16))
			for y := 0; y < 16; y++ {
				for x := 0; x < 100; x++ {
					v := uint8(x)
					im.SetNRGBA(x, y, color.NRGBA{v, v, v + uint8(delta), 255})
				}
			}
			err := CheckColors(name, im)
			if delta == 12 {
				wantCode(t, err, "not-chromatic", "delta exactly 12 is not chromatic")
			} else {
				pngOK(t, err)
			}
		}
	})
	t.Run("scheme-boundaries", func(t *testing.T) {
		pngOK(t, CheckScheme(name, "light", 160))
		pngOK(t, CheckScheme(name, "dark", 96))
		wantCode(t, CheckScheme(name, "light", 159), "wrong-scheme", "light below 160 rejected")
		wantCode(t, CheckScheme(name, "dark", 97), "wrong-scheme", "dark above 96 rejected")
	})
	t.Run("variance-boundaries", func(t *testing.T) {
		for _, delta := range []int{7, 8} {
			im := image.NewNRGBA(image.Rect(0, 0, 64, 16))
			for y := 0; y < 16; y++ {
				for x := 0; x < 64; x++ {
					v := uint8(120 + (x%2)*delta)
					im.SetNRGBA(x, y, color.NRGBA{v, v, v, 255})
				}
			}
			err := CheckRegion(name, im, Rect{Name: "threshold", W: 64, H: 16})
			if delta == 7 {
				wantCode(t, err, "blank-region", "standard deviation below four rejected")
			} else {
				pngOK(t, err)
			}
		}
	})

	t.Run("canonical-shots", func(t *testing.T) {
		expected := []Shot{{"workbench-light.png", 1440, 1100, "light"}, {"workbench-dark.png", 1440, 1100, "dark"}, {"packet.png", 1440, 3700, "light"}, {"graph.png", 1440, 2100, "light"}, {"provenance.png", 1440, 3000, "light"}, {"phone.png", 390, 844, "light"}}
		if len(Shots) != len(expected) {
			t.Fatal("six canonical shots required")
		}
		for _, want := range expected {
			got, ok := shotByName(want.Name)
			if !ok || got != want {
				t.Fatalf("canonical shot %s: got %+v want %+v", want.Name, got, want)
			}
		}
	})

	t.Run("background-bounds", func(t *testing.T) {
		rr := sr
		rr.Background.X = 2000
		wantCode(t, CheckContent(name, bs[name], rr), "rect-out-of-bounds", "background bounds rejected")
	})
	for _, missing := range []string{"validation", "roster-proven", "roster-tested", "roster-attested", "roster-claimed"} {
		t.Run("packet-missing-"+missing, func(t *testing.T) {
			rr := r.Shots["packet.png"]
			content := []Rect{}
			for _, c := range rr.Content {
				if c.Name != missing {
					content = append(content, c)
				}
			}
			rr.Content = content
			wantCode(t, CheckShotRects("packet.png", rr), "rect-missing", "packet required rectangle "+missing+" rejected")
		})
	}

	t.Run("packet-missing-row", func(t *testing.T) {
		rr := r.Shots["packet.png"]
		rr.Content = rr.Content[:4]
		wantCode(t, CheckContent("packet.png", bs["packet.png"], rr), "rect-missing", "packet roster-claimed required")
	})
	t.Run("rectangles-load", func(t *testing.T) {
		p := putJSON(t, t.TempDir(), "rects.json", r)
		_, e := LoadRectangles(p)
		pngOK(t, e)
	})
	for _, kind := range []string{"schema", "unknown", "missing", "empty", "duplicate", "unnamed", "packet-row", "json", "unreadable"} {
		t.Run("rectangles-"+kind, func(t *testing.T) {
			rr := Rectangles{Schema: r.Schema, Shots: map[string]ShotRects{}}
			for n, v := range r.Shots {
				v.Content = append([]Rect(nil), v.Content...)
				rr.Shots[n] = v
			}
			code := "bad-rectangles"
			switch kind {
			case "schema":
				rr.Schema = "wrong"
			case "unknown":
				rr.Shots["unknown.png"] = sr
			case "missing":
				delete(rr.Shots, "phone.png")
			case "empty":
				rr.Shots[name] = ShotRects{}
			case "duplicate":
				v := rr.Shots[name]
				v.Content = append(v.Content, v.Content[0])
				rr.Shots[name] = v
			case "unnamed":
				v := rr.Shots[name]
				v.Content[0].Name = ""
				rr.Shots[name] = v
			case "packet-row":
				v := rr.Shots["packet.png"]
				v.Content = v.Content[:4]
				rr.Shots["packet.png"] = v
				code = "rect-missing"
			}
			dir := t.TempDir()
			p := putJSON(t, dir, "rects.json", rr)
			if kind == "json" {
				pngOK(t, os.WriteFile(p, []byte("{"), 0600))
			}
			if kind == "unreadable" {
				p = filepath.Join(dir, "absent")
			}
			_, e := LoadRectangles(p)
			wantCode(t, e, code, "rectangle configuration "+kind+" rejected")
		})
	}
}

func TestMarketingTwoPassInventory(t *testing.T) {
	bs, r := fixtureSet(t)
	a, b := writeSet(t, bs), writeSet(t, bs)
	d := demoFixture()
	hashes, e := CheckTwoPass(a, b, r)
	pngOK(t, e)
	t.Run("equal-copy-manifest", func(t *testing.T) {
		m, e := BuildCaptureManifest(d, hashes)
		pngOK(t, e)
		raw, e := json.Marshal(m)
		pngOK(t, e)
		var round CaptureManifest
		pngOK(t, json.Unmarshal(raw, &round))
		pngOK(t, VerifyCaptureManifest(round, hashes))
		if m.FrozenHead != "head-ref" {
			t.Fatal("frozen head derived from demo")
		}
		routes := map[string]string{}
		for _, s := range m.Shots {
			routes[s.Name] = s.Route
			shot, _ := shotByName(s.Name)
			if s.Width != shot.Width || s.Height != shot.Height {
				t.Fatal("manifest dimensions")
			}
		}
		for n, want := range map[string]string{"packet.png": "/workbench?object=packet-ref", "graph.png": "/workbench?from=19&entry=19", "provenance.png": "/workbench?object=approval-invocation&payload=1", "phone.png": "/workbench", "workbench-light.png": "/workbench", "workbench-dark.png": "/workbench"} {
			if routes[n] != want {
				t.Errorf("%s route got %q want %q", n, routes[n], want)
			}
		}
	})
	t.Run("missing-shot", func(t *testing.T) {
		dir := writeSet(t, bs)
		pngOK(t, os.Remove(filepath.Join(dir, "phone.png")))
		_, e := CheckPass(dir, r)
		wantCode(t, e, "missing-shot", "missing shot rejected")
		_, e = CheckTwoPass(a, dir, r)
		wantCode(t, e, "missing-shot", "pass B missing shot rejected")
	})
	t.Run("extra-shot", func(t *testing.T) {
		dir := writeSet(t, bs)
		putPNG(t, dir, "extra.png", bs["phone.png"])
		_, e := CheckPass(dir, r)
		wantCode(t, e, "extra-shot", "extra shot rejected")
	})
	t.Run("nonregular", func(t *testing.T) {
		dir := writeSet(t, bs)
		p := filepath.Join(dir, "phone.png")
		pngOK(t, os.Remove(p))
		pngOK(t, os.Symlink(filepath.Join(a, "phone.png"), p))
		_, e := CheckPass(dir, r)
		wantCode(t, e, "missing-shot", "symlink is not regular shot")
	})
	variants := map[string][]byte{"ancillary": addChunk(t, bs["phone.png"], "sRGB", []byte{0})}
	im, e := png.Decode(bytes.NewReader(bs["phone.png"]))
	pngOK(t, e)
	changed := image.NewNRGBA(im.Bounds())
	for y := 0; y < im.Bounds().Dy(); y++ {
		for x := 0; x < im.Bounds().Dx(); x++ {
			changed.Set(x, y, im.At(x, y))
		}
	}
	changed.SetNRGBA(389, 843, color.NRGBA{210, 70, 30, 255})
	variants["pixel"] = encodePNG(t, changed)
	alteredDir := writeSet(t, bs)
	putPNG(t, alteredDir, "phone.png", variants["ancillary"])
	for _, kind := range []string{"ancillary", "pixel"} {
		data := variants[kind]
		t.Run("altered-"+kind, func(t *testing.T) {
			pngOK(t, CheckStructure("phone.png", data))
			pngOK(t, CheckContent("phone.png", data, r.Shots["phone.png"]))
			dir := writeSet(t, bs)
			putPNG(t, dir, "phone.png", data)
			_, e := CheckTwoPass(a, dir, r)
			wantCode(t, e, "pass-mismatch", "altered bytes rejected")
			var refusal *Refusal
			errors.As(e, &refusal)
			h := fmt.Sprintf("%x", sha256.Sum256(data))
			if refusal.Name != "phone.png" || !strings.Contains(refusal.Detail, hashes["phone.png"]) || !strings.Contains(refusal.Detail, h) {
				t.Fatal("pass mismatch must name shot and carry both raw hashes")
			}
		})
	}
	for _, kind := range []string{"hash", "schema", "missing", "extra", "duplicate", "unknown", "hash-set", "unknown-pair", "hash-key", "hash-extra"} {
		t.Run("manifest-"+kind, func(t *testing.T) {
			m, e := BuildCaptureManifest(d, hashes)
			pngOK(t, e)
			hs := hashes
			switch kind {
			case "hash":
				m.Shots[0].SHA256 = strings.Repeat("0", 64)
			case "schema":
				m.Schema = "wrong"
			case "missing":
				m.Shots = m.Shots[:5]
			case "extra":
				m.Shots = append(m.Shots, m.Shots[0])
			case "duplicate":
				m.Shots[1] = m.Shots[0]
			case "unknown":
				m.Shots[0].Name = "unknown.png"
			case "hash-set":
				hs = map[string]string{}
			case "hash-extra":
				hs = map[string]string{}
				for n, v := range hashes {
					hs[n] = v
				}
				hs["unexpected.png"] = "extra"
			case "unknown-pair":
				hs = map[string]string{}
				for n, v := range hashes {
					hs[n] = v
				}
				old := m.Shots[0].Name
				m.Shots[0].Name = "unknown.png"
				delete(hs, old)
				hs["unknown.png"] = m.Shots[0].SHA256
			case "hash-key":
				hs = map[string]string{}
				for n, v := range hashes {
					hs[n] = v
				}
				delete(hs, m.Shots[0].Name)
				hs["unknown.png"] = ""
				m.Shots[0].SHA256 = ""

			}
			wantCode(t, VerifyCaptureManifest(m, hs), "manifest-mismatch", "manifest "+kind+" rejected")
		})
	}
	for _, kind := range []string{"no-head", "duplicate-head", "empty-head", "no-packet", "duplicate-packet", "no-approval", "duplicate-approval", "zero-approval", "no-invocation", "duplicate-invocation", "schema", "version"} {
		t.Run("demo-"+kind, func(t *testing.T) {
			dd := demoFixture()
			switch kind {
			case "no-head":
				dd.References = dd.References[1:]
			case "duplicate-head":
				dd.References = append(dd.References, dd.References[0])
			case "empty-head":
				dd.References[0].Ref = ""
			case "no-packet":
				dd.References = append(dd.References[:1], dd.References[2:]...)
			case "duplicate-packet":
				dd.References = append(dd.References, dd.References[1])
			case "no-approval":
				dd.Why = dd.Why[:1]
			case "duplicate-approval":
				dd.Why = append(dd.Why, dd.Why[1])
			case "zero-approval":
				dd.Why[1].Entry = 0
			case "no-invocation":
				dd.References = dd.References[:3]
			case "duplicate-invocation":
				dd.References = append(dd.References, dd.References[3])
			case "schema":
				dd.Schema = "wrong"
			case "version":
				dd.Version = 2
			}
			if strings.Contains(kind, "approval") {
				_, err := approvalEntry(dd)
				wantCode(t, err, "bad-manifest", "approval guard "+kind+" rejected")
			}
			_, e := BuildCaptureManifest(dd, hashes)
			wantCode(t, e, "bad-manifest", "demo route "+kind+" rejected")
			if kind == "schema" || kind == "version" {
				p := putJSON(t, t.TempDir(), "demo.json", dd)
				_, e = LoadDemoManifest(p)
				wantCode(t, e, "bad-manifest", "demo load "+kind+" rejected")
			}
		})
	}
	t.Run("demo-load", func(t *testing.T) {
		dir := t.TempDir()
		p := putJSON(t, dir, "demo.json", d)
		_, e := LoadDemoManifest(p)
		pngOK(t, e)
		pngOK(t, os.WriteFile(p, []byte("{"), 0600))
		_, e = LoadDemoManifest(p)
		wantCode(t, e, "bad-manifest", "invalid demo JSON rejected")
		_, e = LoadDemoManifest(filepath.Join(dir, "absent"))
		wantCode(t, e, "bad-manifest", "missing demo file rejected")
	})
	t.Run("cli", func(t *testing.T) {
		dir := t.TempDir()
		bin := filepath.Join(dir, "marketingcheck")
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "go", "build", "-o", bin, ".")
		out, e := cmd.CombinedOutput()
		if e != nil {
			t.Fatalf("CLI build: %v %s", e, out)
		}
		dp := putJSON(t, dir, "demo.json", d)
		rp := putJSON(t, dir, "rects.json", r)
		args := []string{"check", "--pass-a", a, "--pass-b", b, "--manifest", dp, "--rectangles", rp}
		invoke := func(t *testing.T, args []string, want int, text string) []byte {
			t.Helper()
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, bin, args...)
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			out, e := cmd.Output()
			rc := 0
			if e != nil {
				var ee *exec.ExitError
				if !errors.As(e, &ee) {
					t.Fatal(e)
				}
				rc = ee.ExitCode()
			}
			if rc != want || !strings.Contains(stderr.String(), text) {
				t.Fatalf("CLI exit: want %d stderr containing %q, got %d: %s", want, text, rc, stderr.String())
			}
			return out
		}
		t.Run("success", func(t *testing.T) {
			raw := invoke(t, args, 0, "")
			var m CaptureManifest
			pngOK(t, json.Unmarshal(raw, &m))
			pngOK(t, VerifyCaptureManifest(m, hashes))
		})
		t.Run("mismatch", func(t *testing.T) {
			bad := append([]string(nil), args...)
			bad[4] = alteredDir
			invoke(t, bad, 1, "marketingcheck: pass-mismatch phone.png:")
		})
		t.Run("new-out", func(t *testing.T) {
			p := filepath.Join(dir, "capture.json")
			invoke(t, append(append([]string(nil), args...), "--out", p), 0, "")
			var m CaptureManifest
			raw, e := os.ReadFile(p)
			pngOK(t, e)
			pngOK(t, json.Unmarshal(raw, &m))
			pngOK(t, VerifyCaptureManifest(m, hashes))
		})
		t.Run("existing-out", func(t *testing.T) {
			p := filepath.Join(dir, "existing.json")
			pngOK(t, os.WriteFile(p, []byte("sentinel"), 0600))
			invoke(t, append(append([]string(nil), args...), "--out", p), 1, "marketingcheck:")
			raw, e := os.ReadFile(p)
			pngOK(t, e)
			if string(raw) != "sentinel" {
				t.Fatal("existing --out must remain unchanged")
			}
		})
		for _, tc := range []struct {
			name string
			args []string
			rc   int
		}{{"unknown", []string{"bogus"}, 2}, {"no-command", nil, 2}, {"missing-flags", []string{"check"}, 2}, {"unknown-flag", []string{"check", "--unknown"}, 2}, {"positional", append(append([]string(nil), args...), "extra"), 2}, {"help", []string{"check", "--help"}, 0}} {
			t.Run(tc.name, func(t *testing.T) { invoke(t, tc.args, tc.rc, "") })
		}
	})
}

// Measured bottom edge (CSS px, device scale 1, getBoundingClientRect().bottom at
// scroll 0) of each tall shot's last required region, max across two independently
// seeded stores. Source: w-marketing-capture-orchestration §10 V51 (iteration 305)
// and the R0 re-confirm V-row V59 (iteration 306, two fresh seeds on Chrome 156.0.8078.48, delta 0.00). A Go test cannot read a
// V-row; these constants are the transcription.
const (
	measuredPacketBottom     = 3564.56 // section[aria-label="decision packet"]
	measuredGraphBottom      = 2042.14 // section[aria-label="world graph"]
	measuredProvenanceBottom = 2881.17 // pre.payload
)

func shotRuleHeight(bottom float64) int { return int(math.Ceil((bottom+40)/100)) * 100 }

func TestMarketingShotGeometry(t *testing.T) {
	for _, tc := range []struct {
		name      string
		bottom    float64
		assertion string
	}{
		{"packet.png", measuredPacketBottom, "packet height must equal the rule value from measuredPacketBottom"},
		{"graph.png", measuredGraphBottom, "graph height must equal the rule value from measuredGraphBottom"},
		{"provenance.png", measuredProvenanceBottom, "provenance height must equal the rule value from measuredProvenanceBottom"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s, ok := shotByName(tc.name)
			want := shotRuleHeight(tc.bottom)
			if !ok || s.Width != 1440 || s.Height != want {
				t.Fatalf("%s: got %+v, want 1440x%d", tc.assertion, s, want)
			}
			if want > 4400 {
				t.Fatalf("%s needs %d px > 4,400: BLOCKED, escalate; never clamp", tc.name, want)
			}
			if tc.bottom > float64(s.Height-40) {
				t.Fatalf("%s: required bottom %.2f must be <= H-40 = %d", tc.name, tc.bottom, s.Height-40)
			}
		})
	}
	t.Run("unchanged", func(t *testing.T) {
		for _, want := range []Shot{{"workbench-light.png", 1440, 1100, "light"}, {"workbench-dark.png", 1440, 1100, "dark"}, {"phone.png", 390, 844, "light"}} {
			if got, ok := shotByName(want.Name); !ok || got != want {
				t.Fatalf("light, dark and phone geometry must stay unchanged: got %+v want %+v", got, want)
			}
		}
	})
}
