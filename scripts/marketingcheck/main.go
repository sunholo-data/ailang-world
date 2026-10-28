package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"
)

const usageText = "usage: marketingcheck check --pass-a DIR --pass-b DIR --manifest DEMO_MANIFEST --rectangles RECTS [--out CAPTURE_MANIFEST]\n" +
	"       marketingcheck routes --manifest DEMO\n" +
	"       marketingcheck observe --root DIR [--out FILE]\n" +
	"       marketingcheck compare --before FILE --after FILE\n" +
	"       marketingcheck inventory --addr URL [--out FILE]\n" +
	"       marketingcheck inventory-compare --before FILE --after FILE\n" +
	"       marketingcheck docs --root DIR [--assets DIR [--capture-manifest FILE --require-assets]]"

func main() { os.Exit(run(os.Args[1:])) }
func run(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, usageText)
		return 2
	}
	switch args[0] {
	case "check":
		return runCheck(args[1:])
	case "routes":
		return runRoutes(args[1:])
	case "observe":
		return runObserve(args[1:])
	case "compare":
		return runCompare(args[1:])
	case "inventory", "inventory-compare":
		return runInventory(args[0], args[1:])
	case "docs":
		return runDocs(args[1:])
	}
	fmt.Fprintln(os.Stderr, usageText)
	return 2
}

// Inventory records the public head, ordered log and distinct transition bodies.
type Inventory struct {
	Schema  string            `json:"schema"`
	Head    json.RawMessage   `json:"head"`
	Entries []inventoryEntry  `json:"entries"`
	Objects []inventoryObject `json:"objects"`
}
type inventoryEntry struct {
	EntryHash     string `json:"entryHash"`
	TransitionRef string `json:"transitionRef"`
}
type inventoryObject struct {
	Ref    string `json:"ref"`
	SHA256 string `json:"sha256"`
}

func ReadInventory(base string) (Inventory, error) {
	inv := Inventory{Schema: "world/demo/capture-inventory/v1"}
	client := &http.Client{Timeout: 10 * time.Second}
	get := func(path string) ([]byte, error) {
		url := strings.TrimRight(base, "/") + path
		resp, err := client.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%s: %s", url, resp.Status)
		}
		return io.ReadAll(resp.Body)
	}
	raw, err := get("/v1/head")
	if err != nil {
		return inv, err
	}
	// The real daemon serves canonical hash text; JSON responses remain JSON.
	inv.Head = raw
	if !json.Valid(raw) {
		inv.Head, _ = json.Marshal(string(raw))
	}
	for from := 0; ; {
		raw, err := get(fmt.Sprintf("/v1/log?from=%d&limit=500", from))
		if err != nil {
			return inv, err
		}
		page := struct {
			Items []inventoryEntry `json:"items"`
		}{}
		if err := json.Unmarshal(raw, &page); err != nil {
			return inv, err
		}
		if len(page.Items) == 0 {
			break
		}
		inv.Entries = append(inv.Entries, page.Items...)
		from += len(page.Items)
	}
	refs := map[string]bool{}
	for _, entry := range inv.Entries {
		refs[entry.TransitionRef] = true
	}
	for ref := range refs {
		raw, err := get("/v1/objects/" + ref)
		if err != nil {
			return inv, err
		}
		sum := sha256.Sum256(raw)
		inv.Objects = append(inv.Objects, inventoryObject{ref, hex.EncodeToString(sum[:])})
	}
	slices.SortFunc(inv.Objects, func(a, b inventoryObject) int { return strings.Compare(a.Ref, b.Ref) })
	return inv, nil
}

func CompareInventory(before, after Inventory) error {
	for i, inv := range []Inventory{before, after} {
		if head := bytes.TrimSpace(inv.Head); len(head) == 0 || bytes.Equal(head, []byte("null")) || bytes.Equal(head, []byte(`""`)) || len(inv.Entries) == 0 {
			return refuse("inventory-empty", []string{"before", "after"}[i], "head or entries empty")
		}
	}
	if !bytes.Equal(before.Head, after.Head) {
		return refuse("inventory-mutated", "inventory", "head differs")
	}
	if len(before.Entries) != len(after.Entries) {
		return refuse("inventory-mutated", "inventory", fmt.Sprintf("entry count %d != %d", len(before.Entries), len(after.Entries)))
	}
	for i, entry := range before.Entries {
		if entry != after.Entries[i] {
			return refuse("inventory-mutated", "inventory", fmt.Sprintf("entry %d differs", i))
		}
	}
	objects := map[string]string{}
	for _, object := range after.Objects {
		objects[object.Ref] = object.SHA256
	}
	if len(before.Objects) != len(objects) {
		return refuse("inventory-mutated", "inventory", "object set differs")
	}
	for _, object := range before.Objects {
		if hash, ok := objects[object.Ref]; !ok || hash != object.SHA256 {
			return refuse("inventory-mutated", "inventory", "object set differs")
		}
	}
	return nil
}

func runInventory(command string, args []string) int {
	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	var addr, out, before, after string
	if command == "inventory" {
		fs.StringVar(&addr, "addr", "", "daemon base URL")
		fs.StringVar(&out, "out", "", "inventory JSON (default stdout)")
	} else {
		fs.StringVar(&before, "before", "", "earlier inventory JSON")
		fs.StringVar(&after, "after", "", "later inventory JSON")
	}
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if fs.NArg() != 0 || (command == "inventory" && addr == "") || (command == "inventory-compare" && (before == "" || after == "")) {
		fs.Usage()
		return 2
	}
	var err error
	if command == "inventory" {
		var inv Inventory
		inv, err = ReadInventory(addr)
		if err == nil {
			err = writeJSON(inv, out)
		}
		if err != nil {
			err = refuse("inventory-unreadable", addr, err.Error())
		}
	} else {
		inputs := [2]Inventory{}
		for i, path := range []string{before, after} {
			var raw []byte
			raw, err = os.ReadFile(path)
			if err == nil {
				err = json.Unmarshal(raw, &inputs[i])
			}
			if err != nil {
				err = refuse("inventory-json", path, err.Error())
				break
			}
		}
		if err == nil {
			err = CompareInventory(inputs[0], inputs[1])
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "marketingcheck: %s\n", err)
		return 1
	}
	return 0
}

// RouteRow carries the capture route and the authoritative shot geometry and scheme.
type RouteRow struct {
	Name   string `json:"name"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Scheme string `json:"scheme"`
	Route  string `json:"route"`
}

func ShotRoutes(d DemoManifest) ([]RouteRow, error) {
	m, err := BuildCaptureManifest(d, map[string]string{})
	if err != nil {
		return nil, err
	}
	rows := make([]RouteRow, 0, len(Shots))
	for i, s := range Shots {
		rows = append(rows, RouteRow{s.Name, s.Width, s.Height, s.Scheme, m.Shots[i].Route})
	}
	return rows, nil
}

func runRoutes(args []string) int {
	fs := flag.NewFlagSet("routes", flag.ContinueOnError)
	demo := fs.String("manifest", "", "public demo manifest JSON")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if fs.NArg() != 0 || *demo == "" {
		fs.Usage()
		return 2
	}
	d, err := LoadDemoManifest(*demo)
	var rows []RouteRow
	if err == nil {
		rows, err = ShotRoutes(d)
	}
	if err == nil {
		err = json.NewEncoder(os.Stdout).Encode(rows)
	}
	if err != nil {
		var r *Refusal
		if errors.As(err, &r) {
			fmt.Fprintf(os.Stderr, "marketingcheck: %s %s: %s\n", r.Code, r.Name, r.Detail)
		} else {
			fmt.Fprintf(os.Stderr, "marketingcheck: bad-manifest output: %s\n", err)
		}
		return 1
	}
	return 0
}

func runCheck(args []string) int {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	a := fs.String("pass-a", "", "first six-shot directory")
	b := fs.String("pass-b", "", "second six-shot directory")
	demo := fs.String("manifest", "", "public demo manifest JSON")
	rectangles := fs.String("rectangles", "", "controller-measured rectangles JSON")
	out := fs.String("out", "", "new absent output path (default stdout)")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if fs.NArg() != 0 || *a == "" || *b == "" || *demo == "" || *rectangles == "" {
		fs.Usage()
		return 2
	}
	if err := check(*a, *b, *demo, *rectangles, *out); err != nil {
		var r *Refusal
		if errors.As(err, &r) {
			fmt.Fprintf(os.Stderr, "marketingcheck: %s %s: %s\n", r.Code, r.Name, r.Detail)
		} else {
			fmt.Fprintf(os.Stderr, "marketingcheck: bad-manifest output: %s\n", err)
		}
		return 1
	}
	return 0
}
func check(a, b, demo, rectangles, out string) error {
	rects, err := LoadRectangles(rectangles)
	if err != nil {
		return err
	}
	hashes, err := CheckTwoPass(a, b, rects)
	if err != nil {
		return err
	}
	d, err := LoadDemoManifest(demo)
	if err != nil {
		return err
	}
	m, err := BuildCaptureManifest(d, hashes)
	if err != nil {
		return err
	}
	if err = VerifyCaptureManifest(m, hashes); err != nil {
		return err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if out == "" {
		_, err = os.Stdout.Write(data)
		return err
	}
	f, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	closeErr := f.Close()
	if err != nil {
		return err
	}
	return closeErr
}

// ObservationRecord is one C11 protected-path observation.
type ObservationRecord struct {
	Path   string `json:"path"`   // slash-separated, relative to the observed root
	SHA256 string `json:"sha256"` // lowercase hex of the raw file bytes
	Lines  int    `json:"lines"`  // bytes.Count(raw, []byte{'\n'})
}

// Observe reads (never writes) the C11 protected set under root.
func Observe(root string) ([]ObservationRecord, error) {
	var paths []string
	world, err := os.ReadDir(filepath.Join(root, "world"))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	for _, e := range world {
		if ok, _ := filepath.Match("*.ail", e.Name()); ok {
			paths = append(paths, "world/"+e.Name())
		}
	}
	paths = append(paths, "go.mod", "go.sum", "design_docs/world-mission.md")
	err = filepath.WalkDir(filepath.Join(root, "tools", "launchd"), func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil
			}
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		paths = append(paths, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		return nil, err
	}
	records := []ObservationRecord{}
	for _, rel := range paths {
		full := filepath.Join(root, filepath.FromSlash(rel))
		info, err := os.Stat(full)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, err
		}
		if !info.Mode().IsRegular() {
			continue
		}
		raw, err := os.ReadFile(full)
		if err != nil {
			return nil, err
		}
		sum := sha256.Sum256(raw)
		records = append(records, ObservationRecord{Path: rel, SHA256: hex.EncodeToString(sum[:]), Lines: bytes.Count(raw, []byte{'\n'})})
	}
	slices.SortFunc(records, func(a, b ObservationRecord) int { return strings.Compare(a.Path, b.Path) })
	return records, nil
}

func observationPaths(records []ObservationRecord) []string {
	paths := make([]string, 0, len(records))
	for _, r := range records {
		paths = append(paths, r.Path)
	}
	slices.Sort(paths)
	return paths
}

// CompareObservations is the C11 equality comparator.
func CompareObservations(before, after []ObservationRecord) error {
	if len(before) == 0 {
		return refuse("observation-empty", "before", "an empty observation proves nothing")
	}
	beforePaths, afterPaths := observationPaths(before), observationPaths(after)
	if !slices.Equal(beforePaths, afterPaths) || len(slices.Compact(slices.Clone(beforePaths))) != len(beforePaths) {
		return refuse("observation-set", "paths", fmt.Sprintf("before has %d paths, after has %d; the protected set changed or repeats a path", len(beforePaths), len(afterPaths)))
	}
	afterByPath := make(map[string]ObservationRecord, len(after))
	for _, a := range after {
		afterByPath[a.Path] = a
	}
	for _, b := range before {
		a := afterByPath[b.Path] // zero value when absent; never skip a miss
		if b.SHA256 != a.SHA256 {
			return refuse("observation-digest", b.Path, fmt.Sprintf("sha256 %s became %s", b.SHA256, a.SHA256))
		}
		if b.Lines != a.Lines {
			return refuse("observation-lines", b.Path, fmt.Sprintf("lines %d became %d", b.Lines, a.Lines))
		}
	}
	return nil
}

func runObserve(args []string) int {
	fs := flag.NewFlagSet("observe", flag.ContinueOnError)
	root := fs.String("root", "", "repository root to observe (read-only)")
	out := fs.String("out", "", "new absent output path (default stdout)")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if fs.NArg() != 0 || *root == "" {
		fs.Usage()
		return 2
	}
	records, err := Observe(*root)
	if err == nil {
		err = writeJSON(records, *out)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "marketingcheck: observation-io %s: %s\n", *root, err)
		return 1
	}
	return 0
}

func writeJSON(records any, out string) error {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if out == "" {
		_, err = os.Stdout.Write(data)
		return err
	}
	f, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	closeErr := f.Close()
	if err != nil {
		return err
	}
	return closeErr
}

func readObservations(path string) ([]ObservationRecord, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, refuse("observation-json", path, err.Error())
	}
	var records []ObservationRecord
	if err := json.Unmarshal(raw, &records); err != nil {
		return nil, refuse("observation-json", path, err.Error())
	}
	return records, nil
}

func runCompare(args []string) int {
	fs := flag.NewFlagSet("compare", flag.ContinueOnError)
	beforePath := fs.String("before", "", "earlier observation JSON")
	afterPath := fs.String("after", "", "later observation JSON")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if fs.NArg() != 0 || *beforePath == "" || *afterPath == "" {
		fs.Usage()
		return 2
	}
	before, err := readObservations(*beforePath)
	if err == nil {
		var after []ObservationRecord
		after, err = readObservations(*afterPath)
		if err == nil {
			err = CompareObservations(before, after)
		}
	}
	if err != nil {
		var r *Refusal
		if errors.As(err, &r) {
			fmt.Fprintf(os.Stderr, "marketingcheck: %s %s: %s\n", r.Code, r.Name, r.Detail)
		} else {
			fmt.Fprintf(os.Stderr, "marketingcheck: observation-json compare: %s\n", err)
		}
		return 1
	}
	return 0
}

// Repository files the docs gate reads (row 150 M8a).
const (
	docsPagePath       = "website/docs/getting-started/marketing-capture.md"
	docsReadmePath     = "README.md"
	docsCaptionsPath   = "website/static/img/v1/captions.md"
	docsSchemaPath     = "docs/marketing-capture-manifest.schema.json"
	docsQuickstartPath = "docs/QUICKSTART.md"
	docsRunbookHeading = "## 13. Capture the release screenshots (attended)"
	docsSchemaID       = "world/demo/capture-manifest/v1"
)

// docsDisclosures are the strings captions.md must carry verbatim (parent §2/§3,
// orchestration §4.1 and §8).
var docsDisclosures = []string{
	"static capture of recorded history",
	"demo-authored",
	"1440x3700",
	"1440x2100",
	"1440x3000",
	"the stored invocation record, opened from the approval entry, with its leading fields legible; the complete reference set is the accompanying `why` JSON",
	"v0.41.0",
	"v0.54.15",
	"scrollbar",
	"overflow",
	"clipped",
	"truncated",
}

// docsAltForbidden are the route and secret markers no alt text may carry.
var docsAltForbidden = []string{"/workbench", "?object=", "/v1/", "session", "token", "proof.key", "Bearer", "sha256:"}

var (
	docsPageImage    = regexp.MustCompile(`<img src=\{useBaseUrl\('/img/v1/([^']*)'\)\} alt="([^"]*)" />`)
	docsPageEmbed    = regexp.MustCompile(`!\[[^\]]*\]\([^)]*img/v1/`)
	docsReadmeEmbed  = regexp.MustCompile(`!\[[^\]]*\]\([^)]*website/static/img/v1/`)
	docsMarkdownLink = regexp.MustCompile(`\[([^\]]*)\]\(([^)\s]+)\)`)
)

// docsReadmeLink matches an ordinary (non-embedded) markdown link to target.
func docsReadmeLink(target string) *regexp.Regexp {
	return regexp.MustCompile(`(^|[^!])\[[^\]]*\]\(` + regexp.QuoteMeta(target) + `\)`)
}

// checkCaptions requires one `## <name>` section per shot, each with one Alt line and
// a caption, no route or secret in any alt, and every disclosure string. It returns the
// alt text per shot.
func checkCaptions(text string) (map[string]string, error) {
	type section struct {
		alts []string
		body int
	}
	sections := map[string]*section{}
	var current *section
	for _, line := range strings.Split(text, "\n") {
		if name, ok := strings.CutPrefix(line, "## "); ok {
			name = strings.TrimSpace(name)
			if _, dup := sections[name]; dup {
				return nil, refuse("docs-captions", name, "duplicate section")
			}
			current = &section{}
			sections[name] = current
			continue
		}
		if current == nil {
			continue
		}
		trimmed := strings.TrimSpace(line)
		if alt, ok := strings.CutPrefix(trimmed, "Alt: "); ok {
			current.alts = append(current.alts, strings.TrimSpace(alt))
		} else if trimmed != "" {
			current.body++
		}
	}
	alts := map[string]string{}
	for name := range sections {
		if filepath.Ext(name) == ".png" {
			if _, ok := shotByName(name); !ok {
				return nil, refuse("docs-captions", name, "section names a PNG that is not in the shot table")
			}
		}
	}
	for _, s := range Shots {
		sec, ok := sections[s.Name]
		if !ok {
			return nil, refuse("docs-captions", s.Name, "no section for this shot")
		}
		if len(sec.alts) != 1 || sec.alts[0] == "" || sec.body == 0 {
			return nil, refuse("docs-captions", s.Name, "a section needs exactly one Alt line and a caption paragraph")
		}
		for _, marker := range docsAltForbidden {
			if strings.Contains(sec.alts[0], marker) {
				return nil, refuse("docs-alt", s.Name, "alt text carries "+marker)
			}
		}
		alts[s.Name] = sec.alts[0]
	}
	for _, d := range docsDisclosures {
		if !strings.Contains(text, d) {
			return nil, refuse("docs-captions", d, "required disclosure is absent")
		}
	}
	return alts, nil
}

// docsLinkTargets returns the markdown link targets of text.
func docsLinkTargets(text string) []string {
	var out []string
	for _, m := range docsMarkdownLink.FindAllStringSubmatch(text, -1) {
		out = append(out, m[2])
	}
	return out
}

// checkDocLinks resolves each relative link target against the real root on disk.
// dir is the directory the links are written relative to, slash-separated and
// relative to root. Absolute URLs, pure anchors and the promoted image targets
// are not repository links and are skipped.
func checkDocLinks(root, dir, text, code string) error {
	for _, target := range docsLinkTargets(text) {
		if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") || strings.HasPrefix(target, "mailto:") || strings.HasPrefix(target, "#") {
			continue
		}
		if strings.Contains(target, "img/v1/") {
			continue
		}
		path, _, _ := strings.Cut(target, "#")
		if path == "" {
			continue
		}
		info, err := os.Stat(filepath.Join(root, filepath.FromSlash(dir), filepath.FromSlash(path)))
		if err != nil || info.IsDir() {
			return refuse(code, target, "relative link does not resolve to a file from "+dirOrRoot(dir))
		}
	}
	return nil
}

func dirOrRoot(dir string) string {
	if dir == "" {
		return "the repository root"
	}
	return dir
}

// subsection returns the body of the first heading named heading, up to the next heading.
func subsection(text, heading string) (string, bool) {
	var out []string
	in := false
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "#") {
			if in {
				break
			}
			if strings.TrimSpace(strings.TrimLeft(line, "#")) == heading {
				in = true
			}
			continue
		}
		if in {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n"), in
}

// checkPage requires the six documented images, each through useBaseUrl with the
// exact inventory name and the captions alt text, and resolves the page's links.
func checkPage(root, page string, alts map[string]string) error {
	const detail = "documented image target must match promoted inventory"
	if !strings.Contains(page, "import useBaseUrl from '@docusaurus/useBaseUrl';") {
		return refuse("docs-image", docsPagePath, "page must import useBaseUrl")
	}
	if docsPageEmbed.MatchString(page) {
		return refuse("docs-image", docsPagePath, "markdown image syntax points at /img/v1; use the useBaseUrl img form")
	}
	matches := docsPageImage.FindAllStringSubmatch(page, -1)
	if len(matches) != len(Shots) || strings.Count(page, "/img/v1/") != len(matches) {
		return refuse("docs-image", docsPagePath, detail)
	}
	seen := map[string]bool{}
	for _, m := range matches {
		name, alt := m[1], m[2]
		if _, ok := shotByName(name); !ok || seen[name] {
			return refuse("docs-image", name, detail)
		}
		seen[name] = true
		for _, marker := range docsAltForbidden {
			if strings.Contains(alt, marker) {
				return refuse("docs-alt", name, "alt text carries "+marker)
			}
		}
		if alts[name] != alt {
			return refuse("docs-image", name, "documented alt must equal captions.md")
		}
	}
	// The page is the public face of the captions, so it carries the same disclosures.
	for _, d := range docsDisclosures {
		if !strings.Contains(page, d) {
			return refuse("docs-disclosure", d, "required disclosure is absent from the website page")
		}
	}
	return checkDocLinks(root, "website/docs/getting-started", page, "docs-link")
}

// checkReadme requires the two screenshots as ordinary links, the runbook link, and
// resolvable relative links in the Screenshots section.
func checkReadme(root, readme string) error {
	if docsReadmeEmbed.MatchString(readme) {
		return refuse("docs-readme", docsReadmePath, "a screenshot is embedded as an image; it must be an ordinary link until promotion")
	}
	for _, target := range []string{"website/static/img/v1/workbench-light.png", "website/static/img/v1/packet.png"} {
		if !docsReadmeLink(target).MatchString(readme) {
			return refuse("docs-readme", target, "screenshot must be linked as an ordinary markdown link")
		}
	}
	if !strings.Contains(readme, "](docs/QUICKSTART.md") {
		return refuse("docs-readme", docsReadmePath, "the reproduction runbook docs/QUICKSTART.md is not linked")
	}
	section, ok := subsection(readme, "Screenshots")
	if !ok {
		return refuse("docs-readme", docsReadmePath, "no Screenshots section")
	}
	if !strings.Contains(section, "static captures of recorded history") {
		return refuse("docs-readme", docsReadmePath, "the Screenshots section must say the images are static captures of recorded history")
	}
	return checkDocLinks(root, "", section, "docs-link")
}

// checkSchema requires the capture-manifest schema shape and that it carries no hash.
func checkSchema(text string) error {
	var doc map[string]any
	if err := json.Unmarshal([]byte(text), &doc); err != nil {
		return refuse("docs-schema", docsSchemaPath, "not valid JSON: "+err.Error())
	}
	props, _ := doc["properties"].(map[string]any)
	schemaProp, _ := props["schema"].(map[string]any)
	if schemaProp["const"] != docsSchemaID {
		return refuse("docs-schema", docsSchemaPath, "properties.schema.const must be "+docsSchemaID)
	}
	required := jsonStrings(doc["required"])
	for _, k := range []string{"schema", "frozen-head", "shots"} {
		if !slices.Contains(required, k) {
			return refuse("docs-schema", docsSchemaPath, "top-level required lacks "+k)
		}
	}
	shots, _ := props["shots"].(map[string]any)
	if shots["minItems"] != float64(len(Shots)) || shots["maxItems"] != float64(len(Shots)) {
		return refuse("docs-schema", docsSchemaPath, "shots must have exactly six items")
	}
	items, _ := shots["items"].(map[string]any)
	itemRequired := jsonStrings(items["required"])
	for _, k := range []string{"name", "sha256", "width", "height", "route"} {
		if !slices.Contains(itemRequired, k) {
			return refuse("docs-schema", docsSchemaPath, "shot items require lacks "+k)
		}
	}
	if docsHex64.MatchString(text) {
		return refuse("docs-schema", docsSchemaPath, "a schema must not carry a 64-hex string")
	}
	return nil
}

var docsHex64 = regexp.MustCompile(`[0-9a-fA-F]{64}`)

func jsonStrings(v any) []string {
	raw, _ := v.([]any)
	var out []string
	for _, x := range raw {
		if s, ok := x.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// CheckDocs proves the M8a wiring against the Shots inventory: the website page's six
// image references, the README's links, the captions and the schema, every relative
// link on disk, and the runbook heading. read returns a repository file by its
// slash-separated path relative to root; nil reads under root.
func CheckDocs(root string, read func(rel string) ([]byte, error)) error {
	if read == nil {
		read = func(rel string) ([]byte, error) {
			return os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		}
	}
	load := func(rel string) (string, error) {
		b, err := read(rel)
		if err != nil {
			return "", refuse("docs-io", rel, err.Error())
		}
		return string(b), nil
	}
	captions, err := load(docsCaptionsPath)
	if err != nil {
		return err
	}
	alts, err := checkCaptions(captions)
	if err != nil {
		return err
	}
	page, err := load(docsPagePath)
	if err != nil {
		return err
	}
	if err := checkPage(root, page, alts); err != nil {
		return err
	}
	readme, err := load(docsReadmePath)
	if err != nil {
		return err
	}
	if err := checkReadme(root, readme); err != nil {
		return err
	}
	schema, err := load(docsSchemaPath)
	if err != nil {
		return err
	}
	if err := checkSchema(schema); err != nil {
		return err
	}
	quick, err := load(docsQuickstartPath)
	if err != nil {
		return err
	}
	if !slices.Contains(strings.Split(quick, "\n"), docsRunbookHeading) {
		return refuse("docs-runbook", docsQuickstartPath, "heading absent: "+docsRunbookHeading)
	}
	return nil
}

// CheckPromotedAssets checks a promoted directory of the six shots: each a regular
// PNG that passes CheckStructure, within the directory cap. In release mode
// (requireAssets) it also reads the attended capture manifest and requires every
// row's SHA-256 to match the bytes and its dimensions to match the shot table.
func CheckPromotedAssets(dir, captureManifest string, requireAssets bool) error {
	if requireAssets && captureManifest == "" {
		return refuse("docs-usage", "capture-manifest", "--require-assets needs --capture-manifest")
	}
	if dir == "" {
		return nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return refuse("missing-shot", dir, err.Error())
	}
	present := map[string]bool{}
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".png" {
			continue
		}
		if _, ok := shotByName(e.Name()); !ok {
			return refuse("extra-shot", e.Name(), "unexpected PNG")
		}
		if !e.Type().IsRegular() {
			return refuse("missing-shot", e.Name(), "must be a regular file")
		}
		present[e.Name()] = true
	}
	for _, s := range Shots {
		if !present[s.Name] {
			return refuse("missing-shot", s.Name, "required PNG absent")
		}
	}
	if err := CheckDirTotal(dir); err != nil {
		return err
	}
	hashes := map[string]string{}
	for _, s := range Shots {
		b, err := os.ReadFile(filepath.Join(dir, s.Name))
		if err != nil {
			return refuse("missing-shot", s.Name, err.Error())
		}
		if err := CheckStructure(s.Name, b); err != nil {
			return err
		}
		sum := sha256.Sum256(b)
		hashes[s.Name] = hex.EncodeToString(sum[:])
	}
	if !requireAssets {
		return nil
	}
	raw, err := os.ReadFile(captureManifest)
	if err != nil {
		return refuse("bad-manifest", captureManifest, err.Error())
	}
	var m CaptureManifest
	if err := json.Unmarshal(raw, &m); err != nil {
		return refuse("bad-manifest", captureManifest, err.Error())
	}
	if err := VerifyCaptureManifest(m, hashes); err != nil {
		return err
	}
	for _, row := range m.Shots {
		shot, _ := shotByName(row.Name)
		if row.Width != shot.Width || row.Height != shot.Height {
			return refuse("manifest-mismatch", row.Name, "dimensions differ from the shot table")
		}
	}
	return nil
}

func runDocs(args []string) int {
	fs := flag.NewFlagSet("docs", flag.ContinueOnError)
	root := fs.String("root", "", "repository root (read-only)")
	assets := fs.String("assets", "", "promoted six-PNG directory to check")
	captureManifest := fs.String("capture-manifest", "", "attended capture-manifest.json (requires --require-assets)")
	requireAssets := fs.Bool("require-assets", false, "release mode: the PNG bytes must match the capture manifest")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if fs.NArg() != 0 || *root == "" || (*requireAssets && (*assets == "" || *captureManifest == "")) || (*captureManifest != "" && !*requireAssets) {
		fs.Usage()
		return 2
	}
	err := CheckDocs(*root, nil)
	if err == nil && *assets != "" {
		err = CheckPromotedAssets(*assets, *captureManifest, *requireAssets)
	}
	if err != nil {
		var r *Refusal
		if errors.As(err, &r) {
			fmt.Fprintf(os.Stderr, "marketingcheck: %s %s: %s\n", r.Code, r.Name, r.Detail)
		} else {
			fmt.Fprintf(os.Stderr, "marketingcheck: docs-io: %s\n", err)
		}
		return 1
	}
	return 0
}
