package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const usageText = "usage: marketingcheck check --pass-a DIR --pass-b DIR --manifest DEMO_MANIFEST --rectangles RECTS [--out CAPTURE_MANIFEST]\n" +
	"       marketingcheck routes --manifest DEMO\n" +
	"       marketingcheck observe --root DIR [--out FILE]\n" +
	"       marketingcheck compare --before FILE --after FILE"

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
	}
	fmt.Fprintln(os.Stderr, usageText)
	return 2
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
		err = writeObservation(records, *out)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "marketingcheck: observation-io %s: %s\n", *root, err)
		return 1
	}
	return 0
}

func writeObservation(records []ObservationRecord, out string) error {
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
