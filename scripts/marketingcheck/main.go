package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() { os.Exit(run(os.Args[1:])) }
func run(args []string) int {
	if len(args) == 0 || args[0] != "check" {
		fmt.Fprintln(os.Stderr, "usage: marketingcheck check --pass-a DIR --pass-b DIR --manifest DEMO_MANIFEST --rectangles RECTS [--out CAPTURE_MANIFEST]")
		return 2
	}
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	a := fs.String("pass-a", "", "first six-shot directory")
	b := fs.String("pass-b", "", "second six-shot directory")
	demo := fs.String("manifest", "", "public demo manifest JSON")
	rectangles := fs.String("rectangles", "", "controller-measured rectangles JSON")
	out := fs.String("out", "", "new absent output path (default stdout)")
	if err := fs.Parse(args[1:]); err != nil {
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
