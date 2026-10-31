package main

// Row 231 M2 (design_docs/planned/w-marketing-position-checker.md §2.2–§2.4, §3 M2, §4 F1–F11):
// `marketingcheck phone` on skeleton renders of the 10 real phone captures and 39 arms derived
// from the seed-5 render, its parity with the banked oracle (ledger: header-15px), its exit codes,
// and the env-gated real-corpus check.

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// ---- skeletons (design §2.3) ----

const (
	phoneTestdata        = "testdata/phone"
	phoneSkeletonSchema  = "world/marketingcheck/phone-skeleton/v1"
	phoneSkeletonCommand = "WORLD_PHONE_CORPUS=$HOME/.ailang/state/world-iter328-phone-corpus WORLD_PHONE_SKELETON_WRITE=1 go test -count=1 -run '^TestMarketingPhoneRealCorpus$' ./scripts/marketingcheck"
)

// phoneSkeleton is the V58 scan's geometry of one real capture, enough to draw it back.
type phoneSkeleton struct {
	Schema          string      `json:"schema"`
	GeneratedBy     string      `json:"generated_by"`
	Source          string      `json:"source"`
	SourceSHA256    string      `json:"source_sha256"`
	W               int         `json:"w"`
	H               int         `json:"h"`
	Background      [3]uint8    `json:"background"`
	Thumb           [3]uint8    `json:"thumb"`
	Bands           []phoneRect `json:"bands"`
	TitleRight      int         `json:"title_right"`
	Scrollbar       phoneRect   `json:"scrollbar"`
	RenderPixSHA256 string      `json:"render_pix_sha256"`
}

func phoneRGB(im image.Image, x, y int) [3]uint8 {
	c := color.RGBAModel.Convert(im.At(x, y)).(color.RGBA)
	return [3]uint8{c.R, c.G, c.B}
}

// renderPhoneSkeleton draws a skeleton: page colour; each band as ink on even columns from its
// X0 with its end columns inked; the header inked through title_right; the thumb filled.
func renderPhoneSkeleton(s phoneSkeleton) *image.RGBA {
	im := image.NewRGBA(image.Rect(0, 0, s.W, s.H))
	bg := color.RGBA{s.Background[0], s.Background[1], s.Background[2], 255}
	ink := color.RGBA{0, 0, 0, 255}
	draw.Draw(im, im.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)
	for i, b := range s.Bands {
		x1 := b.X1
		if i == 0 {
			x1 = s.TitleRight + 1
		}
		for y := b.Y0; y < b.Y1; y++ {
			for x := b.X0; x < x1; x++ {
				if (x-b.X0)%2 == 0 || x == b.X1-1 || x == x1-1 {
					im.Set(x, y, ink)
				}
			}
		}
	}
	t := color.RGBA{s.Thumb[0], s.Thumb[1], s.Thumb[2], 255}
	draw.Draw(im, image.Rect(s.Scrollbar.X0, s.Scrollbar.Y0, s.Scrollbar.X1, s.Scrollbar.Y1), &image.Uniform{t}, image.Point{}, draw.Src)
	return im
}

// derivePhoneSkeleton returns the skeleton JSON bytes of one real capture and its decoded image.
func derivePhoneSkeleton(raw []byte, source string) ([]byte, image.Image, error) {
	im, err := png.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, nil, err
	}
	b := phoneBands(im)
	if len(b) == 0 {
		return nil, nil, fmt.Errorf("%s: zero bands", source)
	}
	s := phoneSkeleton{Schema: phoneSkeletonSchema, GeneratedBy: phoneSkeletonCommand, Source: source, SourceSHA256: sha256Hex(raw),
		W: im.Bounds().Dx(), H: im.Bounds().Dy(), Background: phoneRGB(im, 389, 30), Bands: b, TitleRight: phoneTitleRight(im, b[0]), Scrollbar: phoneScrollbar(im)}
	s.Thumb = phoneRGB(im, (s.Scrollbar.X0+s.Scrollbar.X1)/2, s.Scrollbar.Y0)
	s.RenderPixSHA256 = sha256Hex(renderPhoneSkeleton(s).Pix)
	j, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return nil, nil, err
	}
	return append(j, '\n'), im, nil
}

func loadPhoneSkeleton(t *testing.T, path string) phoneSkeleton {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var s phoneSkeleton
	if err := json.Unmarshal(b, &s); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	if s.Schema != phoneSkeletonSchema {
		t.Fatalf("%s: schema %q", path, s.Schema)
	}
	return s
}

func writePhonePNG(t *testing.T, path string, im image.Image) {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, im); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
}

// ---- arms (Go ports of the banked generators) ----

// judgePhoneFile decodes one PNG as the subcommand does and returns its r3 reasons and scan.
func judgePhoneFile(t *testing.T, path string) ([]string, PhoneScan) {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	im, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	s := ScanPhone(im)
	return PhoneJudge(im, s), s
}

func TestMarketingPhoneSkeleton(t *testing.T) {
	files, err := filepath.Glob(filepath.Join(phoneTestdata, "skeleton-*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 10 {
		t.Fatalf("%d skeletons, want 10", len(files))
	}
	for _, f := range files {
		t.Run(strings.TrimSuffix(filepath.Base(f), ".json"), func(t *testing.T) {
			s := loadPhoneSkeleton(t, f)
			if want := "skeleton-" + strings.TrimSuffix(s.Source, ".png") + ".json"; filepath.Base(f) != want {
				t.Errorf("source %s belongs in %s", s.Source, want)
			}
			r := renderPhoneSkeleton(s)
			if got := sha256Hex(r.Pix); got != s.RenderPixSHA256 {
				t.Errorf("render_pix_sha256: got %s, recorded %s", got, s.RenderPixSHA256)
			}
			var buf bytes.Buffer
			if err := png.Encode(&buf, r); err != nil {
				t.Fatal(err)
			}
			rt, err := png.Decode(bytes.NewReader(buf.Bytes()))
			if err != nil {
				t.Fatal(err)
			}
			if rgba, ok := rt.(*image.RGBA); !ok || !bytes.Equal(rgba.Pix, r.Pix) {
				t.Errorf("PNG round trip is not lossless (%T)", rt)
			}
			bands := phoneBands(rt)
			if fmt.Sprint(bands) != fmt.Sprint(s.Bands) {
				t.Errorf("render bands %v, recorded %v", bands, s.Bands)
			}
			sc := ScanPhone(rt)
			if sc.W != s.W || sc.H != s.H || sc.Scrollbar != s.Scrollbar || sc.TitleRight != s.TitleRight {
				t.Errorf("render scan %dx%d thumb %v titleRight %d, recorded %dx%d %v %d", sc.W, sc.H, sc.Scrollbar, sc.TitleRight, s.W, s.H, s.Scrollbar, s.TitleRight)
			}
			if sc.Header != phoneR1Header || sc.FirstEntry != phoneR1Entry {
				t.Errorf("render header %v entry %v, want R1 %v %v", sc.Header, sc.FirstEntry, phoneR1Header, phoneR1Entry)
			}
			if why := PhoneJudge(rt, sc); len(why) != 0 {
				t.Errorf("render REFUSE: %q", why)
			}
		})
	}
}

func TestMarketingPhoneCLI(t *testing.T) {
	dir := t.TempDir()
	seed5 := renderPhoneSkeleton(loadPhoneSkeleton(t, filepath.Join(phoneTestdata, "skeleton-r1c-seed5-pass-a-phone.json")))
	accept := filepath.Join(dir, "accept.png")
	writePhonePNG(t, accept, seed5)
	// blank is the v67/blank arm: a 390x844 page of luma 250 with no ink.
	blankIm := image.NewNRGBA(image.Rect(0, 0, 390, 844))
	draw.Draw(blankIm, blankIm.Bounds(), &image.Uniform{color.NRGBA{250, 250, 250, 255}}, image.Point{}, draw.Src)
	blank := filepath.Join(dir, "blank.png")
	writePhonePNG(t, blank, blankIm)
	notPNG := filepath.Join(dir, "not.png")
	if err := os.WriteFile(notPNG, []byte("not a png\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		name   string
		args   []string
		rc     int
		stdout []string
		stderr string
	}{
		{"accept-rc", []string{"accept=" + accept}, 0, []string{`^accept +ACCEPT 390x844 bands=21 `, `^TOTAL candidates=1 accept=1 refuse=0$`}, ""},
		{"refuse-rc", []string{"accept=" + accept, "refuse=" + blank}, 1, []string{`^accept +ACCEPT `, `^refuse +REFUSE 390x844 bands=0 `, `^TOTAL candidates=2 accept=1 refuse=1$`}, ""},
		{"blank", []string{"blank=" + blank}, 1, []string{`^blank +REFUSE 390x844 bands=0 `, `^    zero-bands: no ink rows read$`}, ""},
		{"decode-failed", []string{"x=" + notPNG}, 1, nil, "marketingcheck: decode-failed " + notPNG},
		{"missing-file", []string{"x=" + filepath.Join(dir, "absent.png")}, 1, nil, "marketingcheck: decode-failed"},
		{"no-candidates", nil, 2, nil, "no candidates"},
		{"bad-token", []string{accept}, 2, nil, "want LABEL=PATH"},
		{"duplicate-label", []string{"a=" + accept, "a=" + blank}, 2, nil, "want LABEL=PATH"},
		{"empty-label", []string{"=" + accept}, 2, nil, "want LABEL=PATH"},
		{"bad-flag", []string{"--bogus"}, 2, nil, ""},
		{"help", []string{"--help"}, 0, nil, phoneUsage},
	} {
		t.Run(c.name, func(t *testing.T) {
			var out, errb bytes.Buffer
			rc := runPhone(c.args, &out, &errb)
			if rc != c.rc {
				t.Fatalf("rc %d, want %d\nstdout: %s\nstderr: %s", rc, c.rc, out.String(), errb.String())
			}
			for _, re := range c.stdout {
				if !regexp.MustCompile(`(?m)` + re).MatchString(out.String()) {
					t.Errorf("stdout lacks %s:\n%s", re, out.String())
				}
			}
			if !strings.Contains(errb.String(), c.stderr) {
				t.Errorf("stderr lacks %q:\n%s", c.stderr, errb.String())
			}
		})
	}
	t.Run("dispatch", func(t *testing.T) {
		if rc := run([]string{"phone"}); rc != 2 {
			t.Fatalf("marketingcheck phone with no arguments: rc %d, want 2", rc)
		}
		if !strings.Contains(usageText, "marketingcheck phone LABEL=PNG...") {
			t.Fatal("usageText lacks the phone line")
		}
	})
}

// phoneRealThumbX1 is each real image's thumb X1 (orchestration V66: 324–340; V50 329).
var phoneRealThumbX1 = map[string]int{
	"eval-p1-pass-a-phone.png":   338,
	"eval-p2-pass-a-phone.png":   324,
	"eval-p3-pass-a-phone.png":   326,
	"r1-seed1-pass-a-phone.png":  324,
	"r1-seed2-pass-a-phone.png":  340,
	"r1b-seed3-pass-a-phone.png": 333,
	"r1b-seed4-pass-a-phone.png": 332,
	"r1c-seed5-pass-a-phone.png": 333,
	"r1c-seed6-pass-a-phone.png": 334,
	"v50-reference-phone.png":    329,
}

// TestMarketingPhoneRealCorpus is the one env-gated test (design §2.3, §2.5 step 1). Without
// WORLD_PHONE_CORPUS it skips; with it, a missing corpus or any hash difference is BLOCKED.
// WORLD_PHONE_SKELETON_WRITE=1 (re)writes the 10 committed skeletons from the corpus first.
func TestMarketingPhoneRealCorpus(t *testing.T) {
	corpus := os.Getenv("WORLD_PHONE_CORPUS")
	if corpus == "" {
		t.Skip("WORLD_PHONE_CORPUS unset: the real phone corpus is host-local (design §2.3); capture-day step 1 runs this test")
	}
	sumsFile, err := os.Open(filepath.Join(phoneTestdata, "corpus-SHA256SUMS"))
	if err != nil {
		t.Fatal(err)
	}
	defer sumsFile.Close()
	var names []string
	sc := bufio.NewScanner(sumsFile)
	for sc.Scan() {
		f := strings.Fields(sc.Text())
		if len(f) != 2 {
			t.Fatalf("corpus-SHA256SUMS: malformed line %q", sc.Text())
		}
		raw, err := os.ReadFile(filepath.Join(corpus, f[1]))
		if err != nil {
			t.Fatalf("BLOCKED: corpus %s: %v", f[1], err)
		}
		if got := sha256Hex(raw); got != f[0] {
			t.Fatalf("BLOCKED: corpus %s hash=false: sha256 %s, committed %s", f[1], got, f[0])
		}
		names = append(names, f[1])
	}
	if len(names) != 11 {
		t.Fatalf("BLOCKED: corpus-SHA256SUMS lists %d images, want 11", len(names))
	}
	write := os.Getenv("WORLD_PHONE_SKELETON_WRITE") == "1"
	byteEqual, desktop := 0, ""
	for _, name := range names {
		path := filepath.Join(corpus, name)
		if !strings.HasSuffix(name, "-phone.png") {
			desktop = path
			continue
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		fresh, _, err := derivePhoneSkeleton(raw, name)
		if err != nil {
			t.Fatal(err)
		}
		skel := filepath.Join(phoneTestdata, "skeleton-"+strings.TrimSuffix(name, ".png")+".json")
		if write {
			if err := os.WriteFile(skel, fresh, 0o644); err != nil {
				t.Fatal(err)
			}
		}
		committed, err := os.ReadFile(skel)
		if err != nil || !bytes.Equal(fresh, committed) {
			t.Errorf("%s: fresh skeleton is not byte-equal to %s (%v)", name, skel, err)
		} else {
			byteEqual++
		}
		reasons, s := judgePhoneFile(t, path)
		if len(reasons) != 0 || s.Header != phoneR1Header || s.FirstEntry != phoneR1Entry {
			t.Errorf("%s: want ACCEPT with header and entry == R1; got %q header %v entry %v", name, reasons, s.Header, s.FirstEntry)
		}
		if want, ok := phoneRealThumbX1[name]; !ok || s.Scrollbar.X1 != want {
			t.Errorf("%s: thumb X1 %d, want %d", name, s.Scrollbar.X1, want)
		}
	}
	if byteEqual != 10 {
		t.Errorf("CHECK images=10 byte-equal=%d", byteEqual)
	}
	if desktop == "" {
		t.Fatal("BLOCKED: the corpus lacks the seed-5 desktop shot")
	}
	if reasons, _ := judgePhoneFile(t, desktop); len(reasons) == 0 {
		t.Errorf("desktop %s: want REFUSE, got ACCEPT", desktop)
	}
	t.Logf("CHECK images=10 byte-equal=%d", byteEqual)
}
