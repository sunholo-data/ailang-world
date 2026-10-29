package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/daemon"
)

func TestServeSurfaceFollowFlag(t *testing.T) {
	got := captureServeConfig(t)
	var out, errOut bytes.Buffer
	if code := run([]string{"serve", "--db", "/x/world.db", "--surface-follow", "/x/follow.json"}, &out, &errOut); code != exitOK || got.SurfaceFollow != "/x/follow.json" {
		t.Fatalf("flag not forwarded: %d %+v %s", code, got, errOut.String())
	}
	*got = daemon.Config{}
	errOut.Reset()
	if code := run([]string{"serve", "--db", "/x/world.db", "--surface-follow", "/x/follow.json", "--surface-index", "/x/index.json"}, &out, &errOut); code != exitUsage ||
		!strings.Contains(errOut.String(), "surface follow refused") || !strings.Contains(errOut.String(), "mutually exclusive") ||
		!strings.Contains(errOut.String(), "next: ailang-worldd serve --help") || got.DBPath != "" {
		t.Fatalf("pair reached daemon or bad refusal: %d %+v %s", code, got, errOut.String())
	}
	out.Reset()
	errOut.Reset()
	if code := run([]string{"serve", "--help"}, &out, &errOut); code != exitOK {
		t.Fatalf("help exit %d", code)
	}
	_, block, ok := strings.Cut(out.String(), "serve flags:")
	if !ok || !strings.Contains(block, "\n  --surface-follow <file>\n") || !strings.Contains(block, "mutually exclusive with --surface-index") {
		t.Fatal(out.String())
	}
}

func TestHomeServiceSurfaceFollow(t *testing.T) {
	for _, kind := range []string{"forwarded", "relative", "both"} {
		t.Run(kind, func(t *testing.T) {
			f := newServiceFixture(t)
			root, err := filepath.EvalSymlinks(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(root, "follow.json")
			if kind == "relative" {
				path = "follow.json"
			}
			args := append(append([]string(nil), f.args...), "--surface-follow", path)
			if kind == "both" {
				args = append(args, "--surface-index", filepath.Join(root, "index.json"))
			}
			code, msg := f.run(args...)
			if kind != "forwarded" {
				if code != 1 || len(f.writes) != 0 || !strings.Contains(msg, "  why:") || !strings.Contains(msg, "  next:") {
					t.Fatalf("%d %s writes=%v", code, msg, f.writes)
				}
				if kind == "both" && (!strings.Contains(msg, "mutually exclusive") || !strings.Contains(msg, "ailang-worldd serve --help")) {
					t.Fatal(msg)
				}
				return
			}
			if code != 0 {
				t.Fatalf("%d %s", code, msg)
			}
			plist, err := os.ReadFile(f.env.paths.Plist)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(plist), "<string>--surface-follow</string>") || !strings.Contains(string(plist), "<string>"+path+"</string>") {
				t.Fatal(string(plist))
			}
		})
	}
}
