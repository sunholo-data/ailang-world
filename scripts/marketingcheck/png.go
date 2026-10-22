package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"regexp"
)

type Refusal struct{ Code, Name, Detail string }

func (r *Refusal) Error() string             { return fmt.Sprintf("%s %s: %s", r.Code, r.Name, r.Detail) }
func refuse(code, name, detail string) error { return &Refusal{code, name, detail} }

const (
	FileCap = 1_048_576
	DirCap  = 8_388_608
	// Content thresholds are provisional until M7c measures real renders.
	MinColors         = 32
	ChromaticDelta    = 12
	ChromaticFraction = 0.01
	LightMinY         = 160
	DarkMaxY          = 96
	MinRectW          = 64
	MinRectH          = 16
	MinRectStdDev     = 4.0
	PNGSignature      = "\x89PNG\r\n\x1a\n"
)

var pngName = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*\.png$`)
var allowedChunks = map[string]bool{"IHDR": true, "PLTE": true, "IDAT": true, "IEND": true, "sRGB": true, "gAMA": true, "cHRM": true, "iCCP": true, "pHYs": true}
var forbiddenChunks = map[string]bool{"tEXt": true, "zTXt": true, "iTXt": true, "eXIf": true}

type Shot struct {
	Name          string
	Width, Height int
	Scheme        string
}

var Shots = []Shot{
	{"workbench-light.png", 1440, 1100, "light"}, {"workbench-dark.png", 1440, 1100, "dark"},
	{"packet.png", 1440, 1400, "light"}, {"graph.png", 1440, 1100, "light"},
	{"provenance.png", 1440, 1400, "light"}, {"phone.png", 390, 844, "light"},
}

func shotByName(name string) (Shot, bool) {
	for _, s := range Shots {
		if s.Name == name {
			return s, true
		}
	}
	return Shot{}, false
}

// CheckStructure deliberately has no decoder dependency: CRC and policy stand alone.
func CheckStructure(name string, b []byte) error {
	if !pngName.MatchString(name) {
		return refuse("bad-name", name, "expected a lowercase PNG basename")
	}
	if len(b) > FileCap {
		return refuse("file-cap-exceeded", name, fmt.Sprintf("%d > %d", len(b), FileCap))
	}
	if !bytes.HasPrefix(b, []byte(PNGSignature)) {
		return refuse("bad-signature", name, "invalid PNG signature")
	}
	first, ended := true, false
	for i := len(PNGSignature); i < len(b); {
		if ended {
			if len(b)-i >= 12 && string(b[i+4:i+8]) == "IEND" && binary.BigEndian.Uint32(b[i:i+4]) == 0 {
				return refuse("multiple-iend", name, "second IEND")
			}
			return refuse("trailing-bytes", name, "bytes after IEND")
		}
		if len(b)-i < 12 {
			return refuse("truncated-chunk", name, "incomplete chunk header or CRC")
		}
		n := uint64(binary.BigEndian.Uint32(b[i : i+4]))
		end := uint64(i) + 12 + n
		if end > uint64(len(b)) {
			return refuse("truncated-chunk", name, "incomplete chunk data")
		}
		typ := string(b[i+4 : i+8])
		dataEnd := int(end) - 4
		if first && typ != "IHDR" {
			return refuse("ihdr-not-first", name, typ)
		}
		if typ == "IHDR" && n != 13 {
			return refuse("ihdr-length", name, "IHDR must have 13 bytes")
		}
		first = false
		if crc32.ChecksumIEEE(b[i+4:dataEnd]) != binary.BigEndian.Uint32(b[dataEnd:int(end)]) {
			return refuse("bad-crc", name, typ)
		}
		if forbiddenChunks[typ] {
			return refuse("forbidden-chunk", name, typ)
		}
		if !allowedChunks[typ] {
			return refuse("unknown-chunk", name, typ)
		}
		ended = typ == "IEND"
		i = int(end)
	}
	if !ended {
		return refuse("missing-iend", name, "no terminal IEND")
	}
	return nil
}
func CheckDirTotal(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return refuse("dir-cap-exceeded", dir, err.Error())
	}
	var total int64
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".png" {
			info, err := e.Info()
			if err != nil {
				return refuse("dir-cap-exceeded", e.Name(), err.Error())
			}
			total += info.Size()
		}
	}
	if total > DirCap {
		return refuse("dir-cap-exceeded", dir, fmt.Sprintf("%d > %d", total, DirCap))
	}
	return nil
}

type Rect struct {
	Name string `json:"name,omitempty"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	W    int    `json:"w"`
	H    int    `json:"h"`
}
type ShotRects struct {
	Background Rect   `json:"background"`
	Content    []Rect `json:"content"`
}
type Rectangles struct {
	Schema string               `json:"schema"`
	Shots  map[string]ShotRects `json:"shots"`
}

func CheckShotRects(name string, r ShotRects) error {
	if len(r.Content) == 0 {
		return refuse("bad-rectangles", name, "content rectangles required")
	}
	seen := map[string]bool{}
	for _, c := range r.Content {
		if c.Name == "" || seen[c.Name] {
			return refuse("bad-rectangles", name, "empty or duplicate content name")
		}
		seen[c.Name] = true
	}
	if name == "packet.png" {
		for _, s := range []string{"validation", "roster-proven", "roster-tested", "roster-attested", "roster-claimed"} {
			if !seen[s] {
				return refuse("rect-missing", name, s)
			}
		}
	}
	return nil
}
func LoadRectangles(path string) (Rectangles, error) {
	var r Rectangles
	b, err := os.ReadFile(path)
	if err != nil {
		return r, refuse("bad-rectangles", path, err.Error())
	}
	if err = json.Unmarshal(b, &r); err != nil {
		return r, refuse("bad-rectangles", path, err.Error())
	}
	if r.Schema != "world/demo/capture-rectangles/v1" {
		return r, refuse("bad-rectangles", path, "wrong schema")
	}
	for name := range r.Shots {
		if _, ok := shotByName(name); !ok {
			return r, refuse("bad-rectangles", name, "unknown shot")
		}
	}
	for _, s := range Shots {
		sr, ok := r.Shots[s.Name]
		if !ok {
			return r, refuse("bad-rectangles", s.Name, "missing shot")
		}
		if err = CheckShotRects(s.Name, sr); err != nil {
			return r, err
		}
	}
	return r, nil
}
func CheckRectBounds(name string, r Rect, bounds image.Rectangle, content bool) error {
	if r.X < bounds.Min.X || r.Y < bounds.Min.Y || r.W <= 0 || r.H <= 0 || r.X > bounds.Max.X || r.Y > bounds.Max.Y || r.W > bounds.Max.X-r.X || r.H > bounds.Max.Y-r.Y {
		return refuse("rect-out-of-bounds", name, r.Name)
	}
	if content && (r.W < MinRectW || r.H < MinRectH) {
		return refuse("rect-too-small", name, r.Name)
	}
	return nil
}
func rgbAt(im image.Image, x, y int) (uint8, uint8, uint8) {
	switch p := im.(type) {
	case *image.RGBA:
		i := p.PixOffset(x, y)
		if p.Pix[i+3] == 255 {
			return p.Pix[i], p.Pix[i+1], p.Pix[i+2]
		}
	case *image.NRGBA:
		i := p.PixOffset(x, y)
		return p.Pix[i], p.Pix[i+1], p.Pix[i+2]
	}
	c := color.NRGBAModel.Convert(im.At(x, y)).(color.NRGBA)
	return c.R, c.G, c.B
}
func luminance(r, g, b uint8) float64 { return .299*float64(r) + .587*float64(g) + .114*float64(b) }
func CheckColors(name string, im image.Image) error {
	colors := map[uint32]bool{}
	chromatic := 0
	bounds := im.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b := rgbAt(im, x, y)
			if len(colors) <= MinColors {
				colors[uint32(r)<<16|uint32(g)<<8|uint32(b)] = true
			}
			hi, lo := max(r, g, b), min(r, g, b)
			if int(hi)-int(lo) > ChromaticDelta {
				chromatic++
			}
			// Counts only increase; once both whole-image minima are proved, stop scanning.
			if len(colors) > MinColors && float64(chromatic) >= ChromaticFraction*float64(bounds.Dx()*bounds.Dy()) {
				return nil
			}
		}
	}
	if len(colors) <= MinColors {
		return refuse("too-few-colors", name, fmt.Sprintf("%d distinct RGB colors", len(colors)))
	}
	if float64(chromatic)/float64(bounds.Dx()*bounds.Dy()) < ChromaticFraction {
		return refuse("not-chromatic", name, "less than 1% chromatic pixels")
	}
	return nil
}
func rectStats(im image.Image, r Rect) (float64, float64, int) {
	var sum, sum2 float64
	colors := map[uint32]bool{}
	for y := r.Y; y < r.Y+r.H; y++ {
		for x := r.X; x < r.X+r.W; x++ {
			red, g, b := rgbAt(im, x, y)
			v := luminance(red, g, b)
			sum += v
			sum2 += v * v
			if len(colors) < 2 {
				colors[uint32(red)<<16|uint32(g)<<8|uint32(b)] = true
			}
		}
	}
	n := float64(r.W * r.H)
	mean := sum / n
	return mean, math.Sqrt(math.Max(0, sum2/n-mean*mean)), len(colors)
}
func CheckScheme(name, scheme string, mean float64) error {
	if scheme == "light" && mean < LightMinY {
		return refuse("wrong-scheme", name, fmt.Sprintf("light mean Y=%g", mean))
	}
	if scheme == "dark" && mean > DarkMaxY {
		return refuse("wrong-scheme", name, fmt.Sprintf("dark mean Y=%g", mean))
	}
	return nil
}
func CheckRegion(name string, im image.Image, r Rect) error {
	if err := CheckRectBounds(name, r, im.Bounds(), true); err != nil {
		return err
	}
	_, std, colors := rectStats(im, r)
	if std < MinRectStdDev || colors < 2 {
		return refuse("blank-region", name, r.Name)
	}
	return nil
}
func CheckContent(name string, b []byte, rects ShotRects) error {
	im, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		return refuse("decode-failed", name, err.Error())
	}
	shot, ok := shotByName(name)
	if !ok {
		return refuse("bad-name", name, "unknown shot")
	}
	if im.Bounds().Dx() != shot.Width || im.Bounds().Dy() != shot.Height {
		return refuse("wrong-dimensions", name, fmt.Sprintf("got %dx%d, need %dx%d", im.Bounds().Dx(), im.Bounds().Dy(), shot.Width, shot.Height))
	}
	if err = CheckColors(name, im); err != nil {
		return err
	}
	if err = CheckShotRects(name, rects); err != nil {
		return err
	}
	if err = CheckRectBounds(name, rects.Background, im.Bounds(), false); err != nil {
		return err
	}
	mean, _, _ := rectStats(im, rects.Background)
	if err = CheckScheme(name, shot.Scheme, mean); err != nil {
		return err
	}
	for _, r := range rects.Content {
		if err = CheckRegion(name, im, r); err != nil {
			return err
		}
	}
	return nil
}
func CheckPass(dir string, rects Rectangles) (map[string]string, error) {
	if err := CheckDirTotal(dir); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, refuse("missing-shot", dir, err.Error())
	}
	present := map[string]bool{}
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".png" {
			continue
		}
		if _, ok := shotByName(e.Name()); !ok {
			return nil, refuse("extra-shot", e.Name(), "unexpected PNG")
		}
		if !e.Type().IsRegular() {
			return nil, refuse("missing-shot", e.Name(), "must be a regular file")
		}
		present[e.Name()] = true
	}
	for _, s := range Shots {
		if !present[s.Name] {
			return nil, refuse("missing-shot", s.Name, "required PNG absent")
		}
	}
	hashes := map[string]string{}
	for _, s := range Shots {
		if !present[s.Name] {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, s.Name))
		if err != nil {
			return nil, refuse("missing-shot", s.Name, err.Error())
		}
		if err = CheckStructure(s.Name, b); err != nil {
			return nil, err
		}
		if err = CheckContent(s.Name, b, rects.Shots[s.Name]); err != nil {
			return nil, err
		}
		hashes[s.Name] = fmt.Sprintf("%x", sha256.Sum256(b))
	}
	return hashes, nil
}
func CheckTwoPass(passA, passB string, rects Rectangles) (map[string]string, error) {
	a, err := CheckPass(passA, rects)
	if err != nil {
		return nil, err
	}
	b, err := CheckPass(passB, rects)
	if err != nil {
		return nil, err
	}
	for _, s := range Shots {
		if a[s.Name] != b[s.Name] {
			return nil, refuse("pass-mismatch", s.Name, fmt.Sprintf("%s != %s", a[s.Name], b[s.Name]))
		}
	}
	return a, nil
}

type DemoReference struct {
	Role  string `json:"role"`
	Kind  string `json:"kind"`
	Ref   string `json:"ref"`
	Index int64  `json:"index,omitempty"`
}
type DemoWhy struct {
	Role   string `json:"role"`
	Entry  int64  `json:"entry"`
	SHA256 string `json:"sha256"`
}
type DemoManifest struct {
	Schema     string          `json:"schema"`
	Version    int             `json:"version"`
	References []DemoReference `json:"references"`
	Why        []DemoWhy       `json:"why"`
}

func CheckDemoSchema(d DemoManifest) error {
	if d.Schema != "world/demo/public-manifest/v1" || d.Version != 1 {
		return refuse("bad-manifest", "demo", "wrong schema or version")
	}
	return nil
}
func LoadDemoManifest(path string) (DemoManifest, error) {
	var d DemoManifest
	b, err := os.ReadFile(path)
	if err != nil {
		return d, refuse("bad-manifest", path, err.Error())
	}
	if err = json.Unmarshal(b, &d); err != nil {
		return d, refuse("bad-manifest", path, err.Error())
	}
	return d, CheckDemoSchema(d)
}

type CaptureRow struct {
	Name   string `json:"name"`
	SHA256 string `json:"sha256"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Route  string `json:"route"`
}
type CaptureManifest struct {
	Schema     string       `json:"schema"`
	FrozenHead string       `json:"frozen-head"`
	Shots      []CaptureRow `json:"shots"`
}

func manifestRef(d DemoManifest, role string, index int64, indexed bool) (string, error) {
	value := ""
	count := 0
	for _, r := range d.References {
		if r.Role == role && (!indexed || r.Index == index) {
			value = r.Ref
			count++
		}
	}
	if count != 1 || value == "" {
		return "", refuse("bad-manifest", role, "expected exactly one nonempty reference")
	}
	return value, nil
}
func approvalEntry(d DemoManifest) (int64, error) {
	var n int64
	count := 0
	for _, w := range d.Why {
		if w.Role == "approval" {
			n = w.Entry
			count++
		}
	}
	if count != 1 || n <= 0 {
		return 0, refuse("bad-manifest", "approval", "expected exactly one positive approval entry")
	}
	return n, nil
}
func BuildCaptureManifest(d DemoManifest, hashes map[string]string) (CaptureManifest, error) {
	m := CaptureManifest{Schema: "world/demo/capture-manifest/v1"}
	if err := CheckDemoSchema(d); err != nil {
		return m, err
	}
	head, err := manifestRef(d, "frozen-head", 0, false)
	if err != nil {
		return m, err
	}
	m.FrozenHead = head
	packet, err := manifestRef(d, "packet", 0, false)
	if err != nil {
		return m, err
	}
	n, err := approvalEntry(d)
	if err != nil {
		return m, err
	}
	invocation, err := manifestRef(d, "invocation", n, true)
	if err != nil {
		return m, err
	}
	for _, s := range Shots {
		route := "/workbench"
		switch s.Name {
		case "packet.png":
			route += "?object=" + packet
		case "graph.png":
			route += fmt.Sprintf("?from=%d&entry=%d", n, n)
		case "provenance.png":
			route += "?object=" + invocation + "&payload=1"
		}
		m.Shots = append(m.Shots, CaptureRow{s.Name, hashes[s.Name], s.Width, s.Height, route})
	}
	return m, nil
}
func VerifyCaptureManifest(m CaptureManifest, hashes map[string]string) error {
	if m.Schema != "world/demo/capture-manifest/v1" || len(m.Shots) != len(Shots) || len(hashes) != len(Shots) {
		return refuse("manifest-mismatch", "capture", "schema or inventory mismatch")
	}
	seen := map[string]bool{}
	for _, r := range m.Shots {
		hash, ok := hashes[r.Name]
		_, known := shotByName(r.Name)
		if !known || !ok || seen[r.Name] || r.SHA256 != hash {
			return refuse("manifest-mismatch", r.Name, "name or raw SHA256 mismatch")
		}
		seen[r.Name] = true
	}
	return nil
}
