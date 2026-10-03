package pinfetch

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func testCtx(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func hexSum(b []byte) string {
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:])
}

type entry struct {
	name     string
	typeflag byte
	body     []byte
	link     string
}

func makeTarball(t *testing.T, entries ...entry) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for _, e := range entries {
		hdr := &tar.Header{Name: e.name, Typeflag: e.typeflag, Mode: 0o755, Size: int64(len(e.body)), Linkname: e.link}
		if e.typeflag != tar.TypeReg {
			hdr.Size = 0
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if e.typeflag == tar.TypeReg {
			if _, err := tw.Write(e.body); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// release is one fake release served by a TLS httptest server.
type release struct {
	tarball  []byte
	checksum string // the .sha256 body
	binary   []byte
}

func goodRelease(t *testing.T, binary []byte) release {
	tb := makeTarball(t, entry{name: "ailang", typeflag: tar.TypeReg, body: binary})
	return release{tarball: tb, checksum: hexSum(tb) + "  darwin.arm64.ailang.tar.gz\n", binary: binary}
}

func pinFor(r release) Pin {
	return Pin{Release: "v9.9.9", Platform: "darwin/arm64", Asset: "darwin.arm64.ailang.tar.gz",
		TarballSHA256: hexSum(r.tarball), BinarySHA256: hexSum(r.binary)}
}

type server struct {
	srv  *httptest.Server
	hits atomic.Int64
}

func serveRelease(t *testing.T, r release) *server {
	t.Helper()
	s := &server{}
	s.srv = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		s.hits.Add(1)
		if req.Method != http.MethodGet {
			http.Error(w, "GET only", http.StatusMethodNotAllowed)
			return
		}
		switch req.URL.Path {
		case "/dl/v9.9.9/darwin.arm64.ailang.tar.gz":
			_, _ = w.Write(r.tarball)
		case "/dl/v9.9.9/darwin.arm64.ailang.tar.gz.sha256":
			_, _ = w.Write([]byte(r.checksum))
		default:
			http.NotFound(w, req)
		}
	}))
	t.Cleanup(s.srv.Close)
	return s
}

func (s *server) source() *httpSource {
	return newHTTPSource(s.srv.URL+"/dl", []string{"127.0.0.1"}, s.srv.Client())
}

// dirNames lists dir, sorted.
func dirNames(t *testing.T, dir string) []string {
	t.Helper()
	ents, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		t.Fatal(err)
	}
	var names []string
	for _, e := range ents {
		names = append(names, e.Name())
	}
	sort.Strings(names)
	return names
}

// AC3.1: a good release installs, leaves no tarball and writes pin.json.
func TestInstallGoodRelease(t *testing.T) {
	r := goodRelease(t, []byte("#!/bin/sh\necho fake ailang\n"))
	s := serveRelease(t, r)
	dest := filepath.Join(t.TempDir(), "pins", "ailang")
	fixed := time.Date(2026, 10, 3, 12, 0, 0, 0, time.UTC)
	res, err := Install(testCtx(t), s.source(), pinFor(r), dest, Options{Now: func() time.Time { return fixed }})
	if err != nil {
		t.Fatalf("Install: %v", err)
	}
	if res.Status != StatusInstalled {
		t.Fatalf("status %q, want installed", res.Status)
	}
	got, err := os.ReadFile(dest)
	if err != nil || !bytes.Equal(got, r.binary) {
		t.Fatalf("installed bytes differ (err %v)", err)
	}
	info, _ := os.Stat(dest)
	if info.Mode().Perm() != 0o755 {
		t.Fatalf("installed mode %v, want 0755", info.Mode().Perm())
	}
	if names := dirNames(t, filepath.Dir(dest)); strings.Join(names, ",") != "ailang,pin.json" {
		t.Fatalf("destination dir holds %v, want exactly [ailang pin.json] (no tarball, no temp file)", names)
	}
	var rec Record
	data, err := os.ReadFile(filepath.Join(filepath.Dir(dest), RecordName))
	if err != nil || json.Unmarshal(data, &rec) != nil {
		t.Fatalf("pin.json unreadable: %v %s", err, data)
	}
	want := Record{Release: "v9.9.9", Platform: "darwin/arm64", Asset: "darwin.arm64.ailang.tar.gz",
		TarballSHA256: hexSum(r.tarball), BinarySHA256: hexSum(r.binary),
		URL: s.srv.URL + "/dl/v9.9.9/darwin.arm64.ailang.tar.gz", FetchedAt: "2026-10-03T12:00:00Z"}
	if rec != want {
		t.Fatalf("pin.json = %+v, want %+v", rec, want)
	}
	if hits := s.hits.Load(); hits != 2 {
		t.Fatalf("server saw %d requests, want 2 (.sha256 + tarball)", hits)
	}
}

// AC3.6: a second run makes zero requests.
func TestInstallSecondRunMakesNoRequest(t *testing.T) {
	r := goodRelease(t, []byte("fake ailang v9.9.9"))
	s := serveRelease(t, r)
	dest := filepath.Join(t.TempDir(), "ailang")
	src := s.source()
	if _, err := Install(testCtx(t), src, pinFor(r), dest, Options{}); err != nil {
		t.Fatal(err)
	}
	before := s.hits.Load()
	res, err := Install(testCtx(t), src, pinFor(r), dest, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != StatusPresent {
		t.Fatalf("second run status %q, want present", res.Status)
	}
	if after := s.hits.Load(); after != before {
		t.Fatalf("second run made %d request(s), want 0", after-before)
	}
}

// AC3.2: a CONSISTENTLY tampered release — the tarball and its own .sha256
// agree with each other — is refused with exit-3 class IntegrityError, all
// three digests named, and the destination unchanged (absent, and present).
// Kills MUT-HASH-RELEASE-ONLY (trusting the release's own .sha256) and
// MUT-INSTALL-BEFORE-HASH (renaming into place before the digests agree).
func TestInstallRefusesConsistentlyTamperedRelease(t *testing.T) {
	good := goodRelease(t, []byte("the real binary"))
	pin := pinFor(good)
	evil := goodRelease(t, []byte("a trojan binary"))
	s := serveRelease(t, evil)

	for _, tc := range []struct {
		name     string
		existing []byte
		replace  bool
	}{
		{"dest-absent", nil, false},
		{"dest-present-replace", []byte("an older binary"), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			dest := filepath.Join(dir, "ailang")
			if tc.existing != nil {
				if err := os.WriteFile(dest, tc.existing, 0o755); err != nil {
					t.Fatal(err)
				}
			}
			before := dirNames(t, dir)
			_, err := Install(testCtx(t), s.source(), pin, dest, Options{Replace: tc.replace})
			var ie *IntegrityError
			if !errors.As(err, &ie) {
				t.Fatalf("Install err = %v, want *IntegrityError", err)
			}
			msg := err.Error()
			for _, d := range []string{hexSum(evil.tarball), pin.TarballSHA256} {
				if !strings.Contains(msg, d) {
					t.Errorf("refusal %q does not name digest %s", msg, d)
				}
			}
			if ie.Release != hexSum(evil.tarball) || ie.Computed != hexSum(evil.tarball) || ie.Pinned != pin.TarballSHA256 {
				t.Errorf("IntegrityError = %+v", ie)
			}
			after := dirNames(t, dir)
			if strings.Join(after, ",") != strings.Join(before, ",") {
				t.Fatalf("destination dir changed: before %v, after %v", before, after)
			}
			if tc.existing != nil {
				got, _ := os.ReadFile(dest)
				if !bytes.Equal(got, tc.existing) {
					t.Fatal("the existing destination was overwritten by a refused install")
				}
			}
		})
	}
}

// A good tarball whose binary does not match the compiled binary digest is
// refused too (the second compiled-in check).
func TestInstallRefusesBinaryDigestMismatch(t *testing.T) {
	r := goodRelease(t, []byte("binary"))
	pin := pinFor(r)
	pin.BinarySHA256 = strings.Repeat("0", 64)
	s := serveRelease(t, r)
	dir := t.TempDir()
	_, err := Install(testCtx(t), s.source(), pin, filepath.Join(dir, "ailang"), Options{})
	if !IsIntegrity(err) {
		t.Fatalf("err = %v, want integrity", err)
	}
	if names := dirNames(t, dir); len(names) != 0 {
		t.Fatalf("dir holds %v after a refused install", names)
	}
}

func TestInstallRefusesMismatchingExistingWithoutReplace(t *testing.T) {
	r := goodRelease(t, []byte("binary"))
	s := serveRelease(t, r)
	dir := t.TempDir()
	dest := filepath.Join(dir, "ailang")
	if err := os.WriteFile(dest, []byte("something else"), 0o755); err != nil {
		t.Fatal(err)
	}
	_, err := Install(testCtx(t), s.source(), pinFor(r), dest, Options{})
	var em *ExistsMismatch
	if !errors.As(err, &em) {
		t.Fatalf("err = %v, want ExistsMismatch", err)
	}
	if s.hits.Load() != 0 {
		t.Fatalf("a refused existing file still made %d request(s)", s.hits.Load())
	}
	res, err := Install(testCtx(t), s.source(), pinFor(r), dest, Options{Replace: true})
	if err != nil || res.Status != StatusReplaced {
		t.Fatalf("replace: %+v %v", res, err)
	}
	kept, _ := os.ReadFile(res.Kept)
	if string(kept) != "something else" || !strings.HasPrefix(filepath.Base(res.Kept), "ailang.prev-") {
		t.Fatalf("kept file %s = %q", res.Kept, kept)
	}
}

// AC3.3: http://, a redirect to http, and a redirect to a host off the
// allowlist are each refused, with zero bytes written. Kills MUT-HTTPS-OFF
// and MUT-REDIRECT-OPEN.
func TestNetworkBounds(t *testing.T) {
	r := goodRelease(t, []byte("binary"))
	pin := pinFor(r)

	t.Run("plain-http-base", func(t *testing.T) {
		var hits atomic.Int64
		plain := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			hits.Add(1)
			_, _ = w.Write(r.tarball)
		}))
		defer plain.Close()
		src := newHTTPSource(plain.URL+"/dl", []string{"127.0.0.1"}, plain.Client())
		assertRefused(t, src, pin, func() int64 { return hits.Load() })
	})

	// A TLS server that redirects every request to target.
	redirector := func(target string) (*httptest.Server, *atomic.Int64) {
		var hits atomic.Int64
		srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			hits.Add(1)
			http.Redirect(w, req, target+req.URL.Path, http.StatusFound)
		}))
		t.Cleanup(srv.Close)
		return srv, &hits
	}

	t.Run("redirect-to-http", func(t *testing.T) {
		var landed atomic.Int64
		plain := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			landed.Add(1)
			_, _ = w.Write(r.tarball)
		}))
		defer plain.Close()
		srv, _ := redirector(plain.URL)
		src := newHTTPSource(srv.URL+"/dl", []string{"127.0.0.1"}, srv.Client())
		assertRefused(t, src, pin, func() int64 { return landed.Load() })
	})

	t.Run("redirect-to-unlisted-host", func(t *testing.T) {
		good := serveRelease(t, r)
		// Same listener, but named "localhost", which is not on the allowlist.
		target := strings.Replace(good.srv.URL, "127.0.0.1", "localhost", 1)
		srv, _ := redirector(target)
		src := newHTTPSource(srv.URL+"/dl", []string{"127.0.0.1"}, srv.Client())
		assertRefused(t, src, pin, func() int64 { return good.hits.Load() })
	})

	t.Run("too-many-redirects", func(t *testing.T) {
		var srv *httptest.Server
		srv = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, srv.URL+req.URL.Path, http.StatusFound)
		}))
		defer srv.Close()
		src := newHTTPSource(srv.URL+"/dl", []string{"127.0.0.1"}, srv.Client())
		dir := t.TempDir()
		_, err := Install(testCtx(t), src, pin, filepath.Join(dir, "ailang"), Options{})
		if err == nil || !strings.Contains(err.Error(), "redirects") {
			t.Fatalf("err = %v, want a redirect-count refusal", err)
		}
	})

	// Control: the same allowlisted TLS server, no redirect, succeeds — so
	// each refusal above is the bound, not a broken fixture.
	t.Run("control-allowlisted-https", func(t *testing.T) {
		good := serveRelease(t, r)
		if _, err := Install(testCtx(t), good.source(), pin, filepath.Join(t.TempDir(), "ailang"), Options{}); err != nil {
			t.Fatalf("control install failed: %v", err)
		}
	})

	t.Run("production-template", func(t *testing.T) {
		src := GitHub().(*httpSource)
		u := src.Describe("v0.52.1", "darwin.arm64.ailang.tar.gz")
		if u != "https://github.com/sunholo-data/ailang/releases/download/v0.52.1/darwin.arm64.ailang.tar.gz" {
			t.Fatalf("production URL = %s", u)
		}
		if fmt.Sprint(src.hosts) != "[github.com release-assets.githubusercontent.com objects.githubusercontent.com]" {
			t.Fatalf("production hosts = %v", src.hosts)
		}
	})
}

func assertRefused(t *testing.T, src Source, pin Pin, landed func() int64) {
	t.Helper()
	dir := t.TempDir()
	_, err := Install(testCtx(t), src, pin, filepath.Join(dir, "ailang"), Options{})
	var refusal *NetworkRefusal
	if !errors.As(err, &refusal) {
		t.Fatalf("err = %v, want *NetworkRefusal", err)
	}
	if n := landed(); n != 0 {
		t.Fatalf("the refused target was reached %d time(s)", n)
	}
	if names := dirNames(t, dir); len(names) != 0 {
		t.Fatalf("a refused fetch left %v in the destination dir", names)
	}
}

// AC3.4: an over-cap tarball is refused and its temp file removed. Kills
// MUT-NO-CAP.
func TestInstallRefusesOverCapDownload(t *testing.T) {
	r := goodRelease(t, bytes.Repeat([]byte("x"), 4096))
	s := serveRelease(t, r)
	old := MaxTarballBytes
	MaxTarballBytes = int64(len(r.tarball)) - 1
	defer func() { MaxTarballBytes = old }()
	dir := t.TempDir()
	_, err := Install(testCtx(t), s.source(), pinFor(r), filepath.Join(dir, "ailang"), Options{})
	if err == nil || !strings.Contains(err.Error(), "cap") {
		t.Fatalf("err = %v, want an over-cap refusal", err)
	}
	if names := dirNames(t, dir); len(names) != 0 {
		t.Fatalf("over-cap refusal left %v behind", names)
	}
	MaxTarballBytes = int64(len(r.tarball))
	if _, err := Install(testCtx(t), s.source(), pinFor(r), filepath.Join(dir, "ailang"), Options{}); err != nil {
		t.Fatalf("control at exactly the cap failed: %v", err)
	}
}

// AC3.5: tar tricks are refused. Each trick tarball is pinned correctly (its
// digests are the compiled-in ones), so the refusal is extraction's own.
// Kills MUT-TAR-NAME.
func TestInstallRefusesTarTricks(t *testing.T) {
	bin := []byte("binary")
	for _, tc := range []struct {
		name    string
		entries []entry
	}{
		{"symlink-named-ailang", []entry{{name: "ailang", typeflag: tar.TypeSymlink, link: "/bin/sh"}}},
		{"hardlink", []entry{{name: "ailang", typeflag: tar.TypeLink, link: "x"}}},
		{"dotdot", []entry{{name: "../ailang", typeflag: tar.TypeReg, body: bin}}},
		{"absolute", []entry{{name: "/ailang", typeflag: tar.TypeReg, body: bin}}},
		{"nested-name", []entry{{name: "bin/ailang", typeflag: tar.TypeReg, body: bin}}},
		{"other-name", []entry{{name: "ailang2", typeflag: tar.TypeReg, body: bin}}},
		{"extra-entry", []entry{{name: "ailang", typeflag: tar.TypeReg, body: bin}, {name: "README", typeflag: tar.TypeReg, body: []byte("hi")}}},
		{"duplicate", []entry{{name: "ailang", typeflag: tar.TypeReg, body: bin}, {name: "ailang", typeflag: tar.TypeReg, body: []byte("second")}}},
		{"empty", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tb := makeTarball(t, tc.entries...)
			r := release{tarball: tb, checksum: hexSum(tb) + "\n", binary: bin}
			s := serveRelease(t, r)
			dir := t.TempDir()
			_, err := Install(testCtx(t), s.source(), pinFor(r), filepath.Join(dir, "ailang"), Options{})
			if err == nil || IsIntegrity(err) {
				t.Fatalf("err = %v, want an extraction refusal", err)
			}
			if names := dirNames(t, dir); len(names) != 0 {
				t.Fatalf("a refused tarball left %v behind", names)
			}
		})
	}
}

func TestFromDirInstallsOffline(t *testing.T) {
	r := goodRelease(t, []byte("offline binary"))
	base := t.TempDir()
	src := filepath.Join(base, "v9.9.9")
	_ = os.MkdirAll(src, 0o755)
	if err := os.WriteFile(filepath.Join(src, "darwin.arm64.ailang.tar.gz"), r.tarball, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "darwin.arm64.ailang.tar.gz.sha256"), []byte(r.checksum), 0o644); err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(t.TempDir(), "ailang")
	if _, err := Install(testCtx(t), FromDir(base), pinFor(r), dest, Options{}); err != nil {
		t.Fatal(err)
	}
	// The same verification applies offline: a tampered local pair refuses.
	evil := goodRelease(t, []byte("trojan"))
	_ = os.WriteFile(filepath.Join(src, "darwin.arm64.ailang.tar.gz"), evil.tarball, 0o644)
	_ = os.WriteFile(filepath.Join(src, "darwin.arm64.ailang.tar.gz.sha256"), []byte(evil.checksum), 0o644)
	if _, err := Install(testCtx(t), FromDir(base), pinFor(r), filepath.Join(t.TempDir(), "ailang"), Options{}); !IsIntegrity(err) {
		t.Fatalf("tampered offline pair: err = %v, want integrity", err)
	}
}
