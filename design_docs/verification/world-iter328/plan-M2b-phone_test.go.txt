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
	"sort"
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

func phoneCloneNRGBA(src image.Image) *image.NRGBA {
	w := image.NewNRGBA(src.Bounds())
	draw.Draw(w, src.Bounds(), src, image.Point{}, draw.Src)
	return w
}

func phoneCloneRGBA(src image.Image) *image.RGBA {
	o := image.NewRGBA(src.Bounds())
	draw.Draw(o, o.Bounds(), src, src.Bounds().Min, draw.Src)
	return o
}

func phoneNoScrollbar(w *image.NRGBA) {
	b := w.Bounds()
	for x := 0; x < b.Dx(); x++ {
		c := w.At(x, b.Dy()-21)
		for y := b.Dy() - 20; y < b.Dy(); y++ {
			w.Set(x, y, c)
		}
	}
}

func phoneUnclip(w *image.NRGBA) {
	for x := 372; x < 390; x++ {
		c := w.At(x, 44)
		for y := 45; y < 68; y++ {
			w.Set(x, y, c)
		}
	}
}

func phoneFill(im *image.RGBA, x0, y0, x1, y1 int, l uint8) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			im.Set(x, y, color.RGBA{l, l, l, 255})
		}
	}
}

func phoneNoThumb(im *image.RGBA) { phoneFill(im, 0, 824, 390, 844, 250) }

type phoneArm struct {
	Name string
	Make func(src image.Image) image.Image
}

// phoneJudgeArm wraps one v67b judge arm: clone to RGBA, then f.
func phoneJudgeArm(name string, f func(im *image.RGBA)) phoneArm {
	return phoneArm{"judge/" + name, func(src image.Image) image.Image { im := phoneCloneRGBA(src); f(im); return im }}
}

// phoneArms: 12 phoneneg (world-iter315/v67-phoneneg-main.go.txt), 23 judge
// (world-iter315/v67b-judge-phone-arms-main.go.txt) and 4 new (world-iter328/v14-newarms-main.go.txt,
// v11-h15-main.go.txt).
var phoneArms = func() []phoneArm {
	shift := func(n int) phoneArm {
		return phoneArm{fmt.Sprintf("v67/shift-%d", n), func(src image.Image) image.Image {
			b := src.Bounds()
			s := image.NewNRGBA(b)
			draw.Draw(s, b, &image.Uniform{src.At(5, 5)}, image.Point{}, draw.Src)
			draw.Draw(s, image.Rect(0, n, b.Dx(), b.Dy()), src, image.Pt(0, 0), draw.Src)
			return s
		}}
	}
	arms := []phoneArm{
		{"v67/no-scrollbar", func(src image.Image) image.Image { w := phoneCloneNRGBA(src); phoneNoScrollbar(w); return w }},
		{"v67/title-unclipped", func(src image.Image) image.Image { w := phoneCloneNRGBA(src); phoneUnclip(w); return w }},
		{"v67/no-overflow", func(src image.Image) image.Image {
			w := phoneCloneNRGBA(src)
			phoneNoScrollbar(w)
			phoneUnclip(w)
			return w
		}},
		{"v67/crop-370", func(src image.Image) image.Image {
			c := image.NewNRGBA(image.Rect(0, 0, 370, src.Bounds().Dy()))
			draw.Draw(c, c.Bounds(), src, image.Point{}, draw.Src)
			return c
		}},
		{"v67/pad-370", func(src image.Image) image.Image {
			w := phoneCloneNRGBA(src)
			b := src.Bounds()
			draw.Draw(w, image.Rect(370, 0, b.Dx(), b.Dy()), &image.Uniform{src.At(389, 30)}, image.Point{}, draw.Src)
			return w
		}},
		{"v67/full-thumb", func(src image.Image) image.Image {
			w := phoneCloneNRGBA(src)
			tc := src.At(100, 836)
			for y := 833; y <= 840; y++ {
				for x := 3; x <= 387; x++ {
					w.Set(x, y, tc)
				}
			}
			return w
		}},
		shift(2), shift(4), shift(5), shift(10),
		{"v67/timeline-only", func(src image.Image) image.Image {
			b := src.Bounds()
			t := image.NewNRGBA(b)
			draw.Draw(t, b, &image.Uniform{src.At(5, 330)}, image.Point{}, draw.Src)
			draw.Draw(t, image.Rect(0, 0, b.Dx(), b.Dy()-330), src, image.Pt(0, 330), draw.Src)
			return t
		}},
		{"v67/blank", func(src image.Image) image.Image {
			bl := image.NewNRGBA(src.Bounds())
			draw.Draw(bl, bl.Bounds(), &image.Uniform{color.NRGBA{250, 250, 250, 255}}, image.Point{}, draw.Src)
			return bl
		}},
		phoneJudgeArm("title-right-386", func(im *image.RGBA) { phoneFill(im, 387, 45, 390, 68, 250) }),
		phoneJudgeArm("title-right-387", func(im *image.RGBA) { phoneFill(im, 388, 45, 390, 68, 250) }),
		phoneJudgeArm("B1-thumb-split-px150", func(im *image.RGBA) { phoneFill(im, 150, 833, 151, 841, 250) }),
		phoneJudgeArm("B1b-thumb-split-px250", func(im *image.RGBA) { phoneFill(im, 250, 833, 251, 841, 250) }),
		phoneJudgeArm("B1c-thumb-split-px8", func(im *image.RGBA) { phoneFill(im, 8, 833, 9, 841, 250) }),
		phoneJudgeArm("B3-grey-bar-x0-300-no-thumb", func(im *image.RGBA) { phoneNoThumb(im); phoneFill(im, 0, 830, 300, 835, 160) }),
		phoneJudgeArm("B3b-grey-bar-full-width", func(im *image.RGBA) { phoneNoThumb(im); phoneFill(im, 0, 830, 390, 835, 160) }),
		phoneJudgeArm("B4-run-100px", func(im *image.RGBA) { phoneNoThumb(im); phoneFill(im, 6, 833, 106, 841, 190) }),
		phoneJudgeArm("B4b-run-101px", func(im *image.RGBA) { phoneNoThumb(im); phoneFill(im, 6, 833, 107, 841, 190) }),
		phoneJudgeArm("B5-two-pieces-rows-differ", func(im *image.RGBA) {
			phoneNoThumb(im)
			phoneFill(im, 6, 833, 60, 837, 190)
			phoneFill(im, 200, 837, 330, 841, 190)
		}),
		phoneJudgeArm("title-right-385", func(im *image.RGBA) { phoneFill(im, 386, 45, 390, 68, 250) }),
		phoneJudgeArm("short-thumb-x20", func(im *image.RGBA) { phoneNoThumb(im); phoneFill(im, 20, 833, 300, 841, 190) }),
		phoneJudgeArm("short-thumb-x6-130", func(im *image.RGBA) { phoneNoThumb(im); phoneFill(im, 6, 833, 130, 841, 190) }),
		phoneJudgeArm("short-thumb-x6-90", func(im *image.RGBA) { phoneNoThumb(im); phoneFill(im, 6, 833, 90, 841, 190) }),
		phoneJudgeArm("thumb-x9", func(im *image.RGBA) { phoneNoThumb(im); phoneFill(im, 9, 833, 330, 841, 190) }),
		phoneJudgeArm("thumb-x8", func(im *image.RGBA) { phoneNoThumb(im); phoneFill(im, 8, 833, 330, 841, 190) }),
		phoneJudgeArm("thumb-end384", func(im *image.RGBA) { phoneNoThumb(im); phoneFill(im, 6, 833, 384, 841, 190) }),
		phoneJudgeArm("thumb-end385", func(im *image.RGBA) { phoneNoThumb(im); phoneFill(im, 6, 833, 385, 841, 190) }),
		phoneJudgeArm("scattered-pixels-no-thumb", func(im *image.RGBA) {
			phoneNoThumb(im)
			phoneFill(im, 3, 836, 4, 837, 160)
			phoneFill(im, 150, 836, 151, 837, 160)
		}),
		phoneJudgeArm("thumb-above-bottom20", func(im *image.RGBA) { phoneNoThumb(im); phoneFill(im, 6, 815, 330, 823, 190) }),
		phoneJudgeArm("thumb-dark-luma100", func(im *image.RGBA) { phoneNoThumb(im); phoneFill(im, 6, 833, 330, 841, 100) }),
		phoneJudgeArm("tall-image-844-ok-but-thumb-1row", func(im *image.RGBA) { phoneNoThumb(im); phoneFill(im, 6, 840, 330, 841, 190) }),
		phoneJudgeArm("title-ink-spoof-pixel-388", func(im *image.RGBA) { phoneFill(im, 360, 45, 390, 68, 250); phoneFill(im, 388, 50, 389, 51, 0) }),
	}
	solid := func(name string, r image.Rectangle) phoneArm {
		return phoneArm{"new/" + name, func(src image.Image) image.Image {
			w := phoneCloneRGBA(src)
			draw.Draw(w, r, &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.Point{}, draw.Src)
			return w
		}}
	}
	return append(arms,
		solid("header-solid-ink", image.Rect(25, 45, 370, 68)),
		solid("entry-solid-ink", image.Rect(41, 428, 364, 610)),
		phoneArm{"new/tall-845", func(src image.Image) image.Image {
			b := src.Bounds()
			t := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()+1))
			draw.Draw(t, t.Bounds(), &image.Uniform{src.At(5, 5)}, image.Point{}, draw.Src)
			draw.Draw(t, b, src, image.Point{}, draw.Src)
			return t
		}},
		phoneArm{"new/header-15px", func(src image.Image) image.Image {
			w := phoneCloneRGBA(src)
			bg := &image.Uniform{src.At(5, 5)}
			draw.Draw(w, image.Rect(0, 45, 390, 49), bg, image.Point{}, draw.Src)
			draw.Draw(w, image.Rect(0, 64, 390, 68), bg, image.Point{}, draw.Src)
			return w
		}},
	)
}()

func phoneArmByName(t *testing.T, name string) phoneArm {
	t.Helper()
	for _, a := range phoneArms {
		if a.Name == name {
			return a
		}
	}
	t.Fatalf("no phone arm %s", name)
	return phoneArm{}
}

// phoneArmInputs writes every arm of base into dir and returns the labelled PNG paths.
func phoneArmInputs(t *testing.T, dir string, base image.Image) []labelledPath {
	t.Helper()
	var in []labelledPath
	for _, a := range phoneArms {
		p := filepath.Join(dir, strings.ReplaceAll(a.Name, "/", "_")+".png")
		writePhonePNG(t, p, a.Make(base))
		in = append(in, labelledPath{a.Name, p})
	}
	return in
}

// phoneRenderInputs is the CI arm set: the seed-5 and seed-2 base renders plus the 39 arms of
// the seed-5 render (41 PNGs in dir).
func phoneRenderInputs(t *testing.T, dir string) []labelledPath {
	t.Helper()
	seed5 := renderPhoneSkeleton(loadPhoneSkeleton(t, filepath.Join(phoneTestdata, "skeleton-r1c-seed5-pass-a-phone.json")))
	seed2 := renderPhoneSkeleton(loadPhoneSkeleton(t, filepath.Join(phoneTestdata, "skeleton-r1-seed2-pass-a-phone.json")))
	in := []labelledPath{{"seed5-render", filepath.Join(dir, "seed5-render.png")}, {"seed2-render", filepath.Join(dir, "seed2-render.png")}}
	writePhonePNG(t, in[0].Path, seed5)
	writePhonePNG(t, in[1].Path, seed2)
	in = append(in, phoneArmInputs(t, dir, seed5)...)
	if len(in) != 41 {
		t.Fatalf("phone inputs: %d, want 41", len(in))
	}
	return in
}

// phoneTag is a reason's clause tag: zero-bands, (i), (ii), dims, thumb or title-right.
func phoneTag(reason string) string {
	switch {
	case strings.HasPrefix(reason, "zero-bands"):
		return "zero-bands"
	case strings.HasPrefix(reason, "(i) "):
		return "(i)"
	case strings.HasPrefix(reason, "(ii) "):
		return "(ii)"
	case strings.HasPrefix(reason, "(iii′) dims"):
		return "dims"
	case strings.HasPrefix(reason, "(iii′) thumb"):
		return "thumb"
	case strings.HasPrefix(reason, "(iii′) title-right"):
		return "title-right"
	}
	return "?" + reason
}

func phoneVerdict(reasons []string) string {
	if len(reasons) == 0 {
		return "ACCEPT"
	}
	set := map[string]bool{}
	for _, r := range reasons {
		set[phoneTag(r)] = true
	}
	var tags []string
	for k := range set {
		tags = append(tags, k)
	}
	sort.Strings(tags)
	return "REFUSE " + strings.Join(tags, ",")
}

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

// phoneWant is the 41-row verdict table (design §3 M2; V14, V15 column 3): 12 ACCEPT.
var phoneWant = map[string]string{
	"seed5-render":                           "ACCEPT",
	"seed2-render":                           "ACCEPT",
	"v67/no-scrollbar":                       "REFUSE thumb",
	"v67/title-unclipped":                    "REFUSE title-right",
	"v67/no-overflow":                        "REFUSE thumb,title-right",
	"v67/crop-370":                           "REFUSE dims,title-right",
	"v67/pad-370":                            "REFUSE title-right",
	"v67/full-thumb":                         "REFUSE thumb",
	"v67/shift-2":                            "ACCEPT",
	"v67/shift-4":                            "ACCEPT",
	"v67/shift-5":                            "REFUSE (i)",
	"v67/shift-10":                           "REFUSE (i)",
	"v67/timeline-only":                      "REFUSE (i),(ii),thumb,title-right",
	"v67/blank":                              "REFUSE (i),(ii),thumb,title-right,zero-bands",
	"judge/title-right-386":                  "REFUSE title-right",
	"judge/title-right-387":                  "ACCEPT",
	"judge/B1-thumb-split-px150":             "REFUSE thumb",
	"judge/B1b-thumb-split-px250":            "ACCEPT",
	"judge/B1c-thumb-split-px8":              "REFUSE thumb",
	"judge/B3-grey-bar-x0-300-no-thumb":      "ACCEPT",
	"judge/B3b-grey-bar-full-width":          "REFUSE thumb",
	"judge/B4-run-100px":                     "REFUSE thumb",
	"judge/B4b-run-101px":                    "ACCEPT",
	"judge/B5-two-pieces-rows-differ":        "REFUSE thumb",
	"judge/title-right-385":                  "REFUSE title-right",
	"judge/short-thumb-x20":                  "REFUSE thumb",
	"judge/short-thumb-x6-130":               "ACCEPT",
	"judge/short-thumb-x6-90":                "REFUSE thumb",
	"judge/thumb-x9":                         "REFUSE thumb",
	"judge/thumb-x8":                         "ACCEPT",
	"judge/thumb-end384":                     "ACCEPT",
	"judge/thumb-end385":                     "REFUSE thumb",
	"judge/scattered-pixels-no-thumb":        "REFUSE thumb",
	"judge/thumb-above-bottom20":             "REFUSE thumb",
	"judge/thumb-dark-luma100":               "REFUSE thumb",
	"judge/tall-image-844-ok-but-thumb-1row": "ACCEPT",
	"judge/title-ink-spoof-pixel-388":        "REFUSE (i)",
	"new/header-solid-ink":                   "REFUSE (ii)",
	"new/entry-solid-ink":                    "REFUSE (i)",
	"new/tall-845":                           "REFUSE dims",
	"new/header-15px":                        "REFUSE (ii)",
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

func TestMarketingPhoneArms(t *testing.T) {
	accept := 0
	for _, in := range phoneRenderInputs(t, t.TempDir()) {
		t.Run(in.Label, func(t *testing.T) {
			reasons, _ := judgePhoneFile(t, in.Path)
			got := phoneVerdict(reasons)
			if want, ok := phoneWant[in.Label]; !ok || got != want {
				t.Errorf("want %q, got %q (reasons %q)\n\t%q: %q,", want, got, reasons, in.Label, got)
			}
		})
		if phoneWant[in.Label] == "ACCEPT" {
			accept++
		}
	}
	if len(phoneWant) != 41 || accept != 12 {
		t.Errorf("verdict table: %d rows, %d ACCEPT; want 41 and 12", len(phoneWant), accept)
	}
}

var phoneOracleLine = regexp.MustCompile(`^(\S+)\s.*\| \(i\)=(\w+) \(ii\)=(\w+) \|.* dims=(\w+) thumb=(\w+) clip=(\w+) -> r3 (ACCEPT|REFUSE)$`)

// phoneParity runs the oracle over inputs (ref only feeds its r2 field) and compares each
// label's r3 verdict and clause pass/fail with the subcommand's reasons. The only allowed
// difference is the ledger's: header-15px, oracle ACCEPT, phone REFUSE (design §2.2).
func phoneParity(t *testing.T, oracle, ref string, inputs []labelledPath) {
	t.Helper()
	args := []string{ref}
	for _, in := range inputs {
		args = append(args, in.Label+"="+in.Path)
	}
	res := mustRun(t, subprocRunLimit, "", nil, oracle, args...)
	if res.Code != 0 {
		t.Fatalf("oracle rc %d: %s%s", res.Code, res.Stdout, res.Stderr)
	}
	type row struct {
		verdict string
		clauses [5]bool // (i), (ii), dims, thumb, clip
	}
	oracleOut := map[string]row{}
	for _, line := range strings.Split(res.Stdout, "\n") {
		if m := phoneOracleLine.FindStringSubmatch(line); m != nil {
			oracleOut[m[1]] = row{m[7], [5]bool{m[2] == "true", m[3] == "true", m[4] == "true", m[5] == "true", m[6] == "true"}}
		}
	}
	ledger := map[string]bool{"new/header-15px": true}
	compared, diffs := 0, 0
	for _, in := range inputs {
		o, ok := oracleOut[in.Label]
		if !ok {
			t.Errorf("%s: no oracle line", in.Label)
			continue
		}
		reasons, _ := judgePhoneFile(t, in.Path)
		v := "ACCEPT"
		if len(reasons) > 0 {
			v = "REFUSE"
		}
		joined := strings.Join(reasons, "\n")
		clauses := [5]bool{!strings.Contains(joined, "(i) "), !strings.Contains(joined, "(ii) "), !strings.Contains(joined, "(iii′) dims"), !strings.Contains(joined, "(iii′) thumb"), !strings.Contains(joined, "(iii′) title-right")}
		compared++
		if ledger[in.Label] {
			if o.verdict != "ACCEPT" || v != "REFUSE" || !strings.Contains(joined, "(ii) rect-too-small phone.png: header") {
				t.Errorf("ledger %s: want oracle ACCEPT, phone REFUSE (ii) rect-too-small phone.png: header; got oracle %s, phone %s %q", in.Label, o.verdict, v, reasons)
			}
			continue
		}
		if v != o.verdict || clauses != o.clauses {
			diffs++
			t.Errorf("%s: oracle %s %v, phone %s %v %q", in.Label, o.verdict, o.clauses, v, clauses, reasons)
		}
	}
	if compared != len(inputs) || len(oracleOut) != len(inputs) {
		t.Fatalf("compared %d labels (oracle printed %d), want %d", compared, len(oracleOut), len(inputs))
	}
	t.Logf("PARITY phone: inputs=%d diffs=%d ledger=%d", compared, diffs, len(ledger))
}

func TestMarketingPhoneParity(t *testing.T) {
	t.Run("pin-registry", func(t *testing.T) { checkPinRegistry(t, phoneOracleSource, phoneOraclePin) })
	oracle := buildOracle(t, phoneOracleSource, phoneOraclePin)
	inputs := phoneRenderInputs(t, t.TempDir())
	phoneParity(t, oracle, inputs[0].Path, inputs)
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
	var realInputs []labelledPath
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
		label := strings.TrimSuffix(name, "-pass-a-phone.png")
		if label == "v50-reference-phone.png" {
			realInputs = append(realInputs, labelledPath{"v50-ref", path})
			continue
		}
		// Pass B is byte-equal to pass A for every store (V08), so pass A stands for both.
		realInputs = append(realInputs, labelledPath{label + "/A", path}, labelledPath{label + "/B", path})
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
	realInputs = append(realInputs, labelledPath{"desktop-light", desktop})
	seed5Raw, err := os.ReadFile(filepath.Join(corpus, "r1c-seed5-pass-a-phone.png"))
	if err != nil {
		t.Fatal(err)
	}
	seed5, err := png.Decode(bytes.NewReader(seed5Raw))
	if err != nil {
		t.Fatal(err)
	}
	arms := phoneArmInputs(t, t.TempDir(), seed5)
	for _, in := range arms {
		reasons, _ := judgePhoneFile(t, in.Path)
		got := "ACCEPT"
		if len(reasons) > 0 {
			got = "REFUSE"
		}
		if want, _, _ := strings.Cut(phoneWant[in.Label], " "); got != want {
			t.Errorf("real seed-5 arm %s: want %s, got %s %q", in.Label, want, got, reasons)
		}
	}
	realInputs = append(realInputs, arms...)
	if len(realInputs) != 59 {
		t.Fatalf("real parity inputs %d, want 59", len(realInputs))
	}
	oracle := buildOracle(t, phoneOracleSource, phoneOraclePin)
	phoneParity(t, oracle, filepath.Join(corpus, "v50-reference-phone.png"), realInputs)
	t.Logf("CHECK images=10 byte-equal=%d; real arms=%d; parity inputs=%d", byteEqual, len(arms), len(realInputs))
}
