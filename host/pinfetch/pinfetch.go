// Package pinfetch fetches and verifies a pinned AILANG release binary
// (row 138 M3, design_docs/planned/w-worldd-developer-cli.md §3.5, D-CLI-1).
//
// It is the ONLY non-loopback network code in this repository outside
// host/broker, and it has one job: put a byte-verified `ailang` binary at a
// destination path. The boundary gate (host/boundary) protects it as its own
// group and asserts that no store, replay or world-publish closure contains it.
//
// Trust (R-CLI-3): an install is accepted only when three tarball digests
// agree — the release's own `.sha256`, the hash computed while streaming, and
// the digest COMPILED INTO the caller's pin table — and the extracted binary
// hashes to the compiled-in binary digest. The release can therefore not vouch
// for itself (MUT-HASH-RELEASE-ONLY), and nothing is installed before every
// hash has matched (MUT-INSTALL-BEFORE-HASH). Cosign signatures are out of
// scope (D-CLI-5).
//
// Network bounds: https only, on the first URL and on every redirect; the
// host must be on a fixed allowlist; GET only; at most ten redirects; the
// URL is built from a compiled template (there is no URL input); every
// download is size-capped. pinfetch never executes a byte it downloaded.
package pinfetch

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

// Pin is one compiled-in pin-table row: a release asset for one platform and
// the two digests it must hash to.
type Pin struct {
	Release       string // e.g. "v0.41.0"
	Platform      string // GOOS/GOARCH, e.g. "darwin/arm64"
	Asset         string // the release asset name, e.g. "darwin.arm64.ailang.tar.gz"
	TarballSHA256 string // 64 lowercase hex
	BinarySHA256  string // 64 lowercase hex
}

// Caps. Variables (not constants) only so the over-cap arm can be driven with
// small fixtures; production never changes them.
var (
	// MaxChecksumBytes bounds the release `.sha256` asset.
	MaxChecksumBytes int64 = 1 << 10
	// MaxTarballBytes bounds the streamed tarball (128 MiB).
	MaxTarballBytes int64 = 128 << 20
	// MaxBinaryBytes bounds the one extracted `ailang` entry (256 MiB).
	MaxBinaryBytes int64 = 256 << 20
)

// Request budgets (§3.5): 30 s for a small asset, 5 min for the tarball.
const (
	checksumTimeout = 30 * time.Second
	tarballTimeout  = 5 * time.Minute
	maxRedirects    = 10
)

// releaseBase is the compiled URL template. There is no flag or environment
// variable that changes it.
const releaseBase = "https://github.com/sunholo-data/ailang/releases/download"

// gitHubHosts is the https host allowlist: the release page host and the two
// hosts GitHub redirects asset downloads to (measured V18: 302 to
// release-assets.githubusercontent.com).
var gitHubHosts = []string{"github.com", "release-assets.githubusercontent.com", "objects.githubusercontent.com"}

// Source yields one release asset's bytes, at most max of them.
type Source interface {
	// Open returns a reader for asset of release. The caller closes it.
	Open(ctx context.Context, release, asset string, timeout time.Duration) (io.ReadCloser, error)
	// Describe names the source for messages and pin.json.
	Describe(release, asset string) string
	// Requests reports how many network requests (or file opens) were made.
	Requests() int64
}

// httpSource is the network Source: the compiled release template, the host
// allowlist, https-only on every hop.
type httpSource struct {
	base   string
	hosts  []string
	client *http.Client
	count  atomic.Int64
}

// GitHub is the production Source: the compiled release URL template and the
// GitHub host allowlist.
func GitHub() Source { return newHTTPSource(releaseBase, gitHubHosts, &http.Client{}) }

// newHTTPSource wraps client with the redirect policy. It is unexported: the
// only exported constructor is GitHub, whose URL is compiled in.
func newHTTPSource(base string, hosts []string, client *http.Client) *httpSource {
	s := &httpSource{base: strings.TrimSuffix(base, "/"), hosts: hosts}
	c := *client
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= maxRedirects {
			return fmt.Errorf("pinfetch: refusing more than %d redirects", maxRedirects)
		}
		return s.checkURL(req.URL)
	}
	s.client = &c
	return s
}

// checkURL is the per-hop network bound: https and an allowlisted host.
func (s *httpSource) checkURL(u *url.URL) error {
	if u.Scheme != "https" {
		return &NetworkRefusal{URL: u.Redacted(), Reason: "not https"}
	}
	host := u.Hostname()
	for _, h := range s.hosts {
		if host == h {
			return nil
		}
	}
	return &NetworkRefusal{URL: u.Redacted(), Reason: fmt.Sprintf("host %q is not on the allowlist %v", host, s.hosts)}
}

func (s *httpSource) Describe(release, asset string) string {
	return s.base + "/" + release + "/" + asset
}

func (s *httpSource) Requests() int64 { return s.count.Load() }

func (s *httpSource) Open(ctx context.Context, release, asset string, timeout time.Duration) (io.ReadCloser, error) {
	u, err := url.Parse(s.Describe(release, asset))
	if err != nil {
		return nil, fmt.Errorf("pinfetch: parse release URL: %w", err)
	}
	if err := s.checkURL(u); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("pinfetch: build request: %w", err)
	}
	s.count.Add(1)
	resp, err := s.client.Do(req)
	if err != nil {
		cancel()
		var refusal *NetworkRefusal
		if errors.As(err, &refusal) {
			return nil, refusal
		}
		return nil, fmt.Errorf("pinfetch: GET %s: %w", u.Redacted(), err)
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		cancel()
		return nil, fmt.Errorf("pinfetch: GET %s: HTTP %d", u.Redacted(), resp.StatusCode)
	}
	return cancelOnClose{ReadCloser: resp.Body, cancel: cancel}, nil
}

type cancelOnClose struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (c cancelOnClose) Close() error {
	err := c.ReadCloser.Close()
	c.cancel()
	return err
}

// dirSource is the offline Source (`setup --from-dir`): the asset and its
// `.sha256` read from <dir>/<release>/, verified exactly like a download.
// The release subdirectory keeps the two pins apart: both releases name
// their tarball the same.
type dirSource struct {
	dir   string
	count atomic.Int64
}

// FromDir is the offline Source rooted at dir.
func FromDir(dir string) Source { return &dirSource{dir: dir} }

func (d *dirSource) Describe(release, asset string) string {
	return filepath.Join(d.dir, release, asset)
}
func (d *dirSource) Requests() int64 { return d.count.Load() }
func (d *dirSource) Open(_ context.Context, release, asset string, _ time.Duration) (io.ReadCloser, error) {
	d.count.Add(1)
	if release == "" || strings.ContainsAny(release, `/\`) || release == "." || release == ".." {
		return nil, fmt.Errorf("pinfetch: bad release name %q", release)
	}
	return os.Open(filepath.Join(d.dir, release, asset))
}

// NetworkRefusal is a URL refused by the network bounds before any byte of
// its body was read.
type NetworkRefusal struct {
	URL, Reason string
}

func (e *NetworkRefusal) Error() string {
	return fmt.Sprintf("pinfetch: refusing %s: %s", e.URL, e.Reason)
}

// IntegrityError is a digest that did not match. Nothing was installed.
type IntegrityError struct {
	What     string // "tarball" or "binary"
	Release  string // the release .sha256 (tarball only)
	Computed string
	Pinned   string
}

func (e *IntegrityError) Error() string {
	if e.What == "tarball" {
		return fmt.Sprintf("pinfetch: tarball digest mismatch: release .sha256 %s, computed %s, compiled-in pin %s",
			orNone(e.Release), e.Computed, e.Pinned)
	}
	return fmt.Sprintf("pinfetch: %s digest mismatch: computed %s, compiled-in pin %s", e.What, e.Computed, e.Pinned)
}

func orNone(s string) string {
	if s == "" {
		return "(none)"
	}
	return s
}

// IsIntegrity reports whether err is (or wraps) an *IntegrityError.
func IsIntegrity(err error) bool {
	var ie *IntegrityError
	return errors.As(err, &ie)
}

// ExistsMismatch is an existing destination whose hash is not the pin's;
// it is replaced only with Options.Replace.
type ExistsMismatch struct {
	Path, Got, Want string
}

func (e *ExistsMismatch) Error() string {
	return fmt.Sprintf("pinfetch: %s exists with sha256 %s, not the pinned %s; pass --replace to keep it as %s and install the pin",
		e.Path, e.Got, e.Want, prevName(e.Path, e.Got))
}

func prevName(dest, sum string) string {
	return dest + ".prev-" + sum[:8]
}

// Options tunes one Install.
type Options struct {
	// Replace keeps a mismatching existing destination as
	// <dest>.prev-<sha8> and installs the pin in its place.
	Replace bool
	// Now stamps pin.json; nil means time.Now.
	Now func() time.Time
}

// Status is what Install did.
type Status string

const (
	StatusPresent   Status = "present"   // the destination already hashed to the pin; no request made
	StatusInstalled Status = "installed" // downloaded, verified and installed
	StatusReplaced  Status = "replaced"  // installed over a kept mismatching file
)

// Result reports one Install.
type Result struct {
	Status Status
	Dest   string
	Kept   string // the .prev-<sha8> path, StatusReplaced only
}

// Record is pin.json, written next to an installed binary.
type Record struct {
	Release       string `json:"release"`
	Platform      string `json:"platform"`
	Asset         string `json:"asset"`
	TarballSHA256 string `json:"tarballSHA256"`
	BinarySHA256  string `json:"binarySHA256"`
	URL           string `json:"url"`
	FetchedAt     string `json:"fetchedAt"`
}

// RecordName is the provenance file written beside an installed binary.
const RecordName = "pin.json"

// HashFile returns the lowercase hex sha256 of the file at path.
func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Install puts pin's verified binary at dest (dest's directory is created
// 0755 when absent). It never executes a downloaded byte.
func Install(ctx context.Context, src Source, pin Pin, dest string, opts Options) (Result, error) {
	if !isHex64(pin.TarballSHA256) || !isHex64(pin.BinarySHA256) {
		return Result{}, fmt.Errorf("pinfetch: pin %s %s carries a malformed digest", pin.Release, pin.Platform)
	}
	now := time.Now
	if opts.Now != nil {
		now = opts.Now
	}

	// 1. Already there: hash it, make no request.
	existing := ""
	if info, err := os.Lstat(dest); err == nil {
		if !info.Mode().IsRegular() {
			return Result{}, fmt.Errorf("pinfetch: %s exists and is not a regular file", dest)
		}
		sum, err := HashFile(dest)
		if err != nil {
			return Result{}, fmt.Errorf("pinfetch: hash existing %s: %w", dest, err)
		}
		if sum == pin.BinarySHA256 {
			return Result{Status: StatusPresent, Dest: dest}, nil
		}
		if !opts.Replace {
			return Result{}, &ExistsMismatch{Path: dest, Got: sum, Want: pin.BinarySHA256}
		}
		existing = sum
	} else if !errors.Is(err, os.ErrNotExist) {
		return Result{}, fmt.Errorf("pinfetch: stat %s: %w", dest, err)
	}

	dir := filepath.Dir(dest)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Result{}, fmt.Errorf("pinfetch: create %s: %w", dir, err)
	}

	// 2. The release's own .sha256.
	releaseSum, err := fetchChecksum(ctx, src, pin)
	if err != nil {
		return Result{}, err
	}

	// 2b. Stream the tarball into an O_EXCL 0600 temp file, hashing it.
	tarPath, tarSum, err := download(ctx, src, pin, dir)
	if tarPath != "" {
		defer func() { _ = os.Remove(tarPath) }()
	}
	if err != nil {
		return Result{}, err
	}

	// 3. All three tarball digests must agree.
	if releaseSum != tarSum || tarSum != pin.TarballSHA256 {
		return Result{}, &IntegrityError{What: "tarball", Release: releaseSum, Computed: tarSum, Pinned: pin.TarballSHA256}
	}

	// 4. Extract exactly one regular `ailang` entry.
	binPath, binSum, err := extract(tarPath, dir)
	if binPath != "" {
		defer func() { _ = os.Remove(binPath) }() // a no-op once renamed into place
	}
	if err != nil {
		return Result{}, err
	}

	// 5. The binary must hash to the compiled-in digest.
	if binSum != pin.BinarySHA256 {
		return Result{}, &IntegrityError{What: "binary", Computed: binSum, Pinned: pin.BinarySHA256}
	}

	// 6. Only now: mode, keep a replaced file, atomic rename.
	if err := os.Chmod(binPath, 0o755); err != nil {
		return Result{}, fmt.Errorf("pinfetch: chmod %s: %w", binPath, err)
	}
	res := Result{Status: StatusInstalled, Dest: dest}
	if existing != "" {
		kept := prevName(dest, existing)
		if err := os.Rename(dest, kept); err != nil {
			return Result{}, fmt.Errorf("pinfetch: keep %s as %s: %w", dest, kept, err)
		}
		res.Status, res.Kept = StatusReplaced, kept
	}
	if err := os.Rename(binPath, dest); err != nil {
		return Result{}, fmt.Errorf("pinfetch: install %s: %w", dest, err)
	}

	// 7. pin.json; the tarball temp file is removed by the defer.
	rec := Record{Release: pin.Release, Platform: pin.Platform, Asset: pin.Asset,
		TarballSHA256: pin.TarballSHA256, BinarySHA256: pin.BinarySHA256,
		URL: src.Describe(pin.Release, pin.Asset), FetchedAt: now().UTC().Format(time.RFC3339)}
	data, _ := json.MarshalIndent(rec, "", "  ")
	if err := os.WriteFile(filepath.Join(dir, RecordName), append(data, '\n'), 0o644); err != nil {
		return res, fmt.Errorf("pinfetch: write %s: %w", RecordName, err)
	}
	return res, nil
}

// fetchChecksum reads the release `.sha256` (≤ MaxChecksumBytes) and returns
// its first field, which must be 64 hex.
func fetchChecksum(ctx context.Context, src Source, pin Pin) (string, error) {
	rc, err := src.Open(ctx, pin.Release, pin.Asset+".sha256", checksumTimeout)
	if err != nil {
		return "", err
	}
	defer func() { _ = rc.Close() }()
	data, err := io.ReadAll(io.LimitReader(rc, MaxChecksumBytes+1))
	if err != nil {
		return "", fmt.Errorf("pinfetch: read %s.sha256: %w", pin.Asset, err)
	}
	if int64(len(data)) > MaxChecksumBytes {
		return "", fmt.Errorf("pinfetch: %s.sha256 exceeds %d bytes", pin.Asset, MaxChecksumBytes)
	}
	fields := strings.Fields(string(data))
	if len(fields) == 0 || !isHex64(fields[0]) {
		return "", fmt.Errorf("pinfetch: %s.sha256 does not begin with a 64-hex digest", pin.Asset)
	}
	return strings.ToLower(fields[0]), nil
}

// download streams the tarball into an O_EXCL 0600 temp file in dir, hashing
// it as it arrives, capped at MaxTarballBytes. On an error after the temp file
// exists it still returns the path, so the caller's cleanup removes it.
func download(ctx context.Context, src Source, pin Pin, dir string) (string, string, error) {
	rc, err := src.Open(ctx, pin.Release, pin.Asset, tarballTimeout)
	if err != nil {
		return "", "", err
	}
	defer func() { _ = rc.Close() }()
	f, path, err := createExclusive(dir, ".ailang-download-")
	if err != nil {
		return "", "", err
	}
	h := sha256.New()
	n, copyErr := io.Copy(io.MultiWriter(f, h), io.LimitReader(rc, MaxTarballBytes+1))
	closeErr := f.Close()
	if copyErr != nil {
		return path, "", fmt.Errorf("pinfetch: download %s: %w", pin.Asset, copyErr)
	}
	if closeErr != nil {
		return path, "", fmt.Errorf("pinfetch: write %s: %w", path, closeErr)
	}
	if n > MaxTarballBytes {
		return path, "", fmt.Errorf("pinfetch: %s exceeds the %d-byte cap; refused", pin.Asset, MaxTarballBytes)
	}
	return path, hex.EncodeToString(h.Sum(nil)), nil
}

// extract takes exactly one regular entry named `ailang` from the gzip tar at
// tarPath into an O_EXCL temp file in dir. Any link, any other entry, a `..`
// segment or an absolute name refuses the whole archive.
func extract(tarPath, dir string) (string, string, error) {
	f, err := os.Open(tarPath)
	if err != nil {
		return "", "", fmt.Errorf("pinfetch: open tarball: %w", err)
	}
	defer func() { _ = f.Close() }()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return "", "", fmt.Errorf("pinfetch: tarball is not gzip: %w", err)
	}
	defer func() { _ = gz.Close() }()
	tr := tar.NewReader(gz)
	outPath, sum := "", ""
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return outPath, "", fmt.Errorf("pinfetch: read tarball: %w", err)
		}
		if err := checkEntry(hdr); err != nil {
			return outPath, "", err
		}
		if outPath != "" {
			return outPath, "", fmt.Errorf("pinfetch: tarball holds a second %q entry; refused", hdr.Name)
		}
		out, p, err := createExclusive(dir, ".ailang-extract-")
		if err != nil {
			return "", "", err
		}
		outPath = p
		h := sha256.New()
		n, copyErr := io.Copy(io.MultiWriter(out, h), io.LimitReader(tr, MaxBinaryBytes+1))
		closeErr := out.Close()
		if copyErr != nil {
			return outPath, "", fmt.Errorf("pinfetch: extract ailang: %w", copyErr)
		}
		if closeErr != nil {
			return outPath, "", fmt.Errorf("pinfetch: write %s: %w", outPath, closeErr)
		}
		if n > MaxBinaryBytes {
			return outPath, "", fmt.Errorf("pinfetch: the ailang entry exceeds the %d-byte cap; refused", MaxBinaryBytes)
		}
		sum = hex.EncodeToString(h.Sum(nil))
	}
	if outPath == "" {
		return "", "", errors.New("pinfetch: tarball holds no `ailang` entry")
	}
	return outPath, sum, nil
}

// checkEntry admits exactly a regular file named `ailang`.
func checkEntry(hdr *tar.Header) error {
	name := hdr.Name
	if strings.HasPrefix(name, "/") || filepath.IsAbs(name) {
		return fmt.Errorf("pinfetch: tarball entry %q is an absolute path; refused", name)
	}
	for _, seg := range strings.Split(name, "/") {
		if seg == ".." {
			return fmt.Errorf("pinfetch: tarball entry %q has a .. segment; refused", name)
		}
	}
	if hdr.Typeflag != tar.TypeReg {
		return fmt.Errorf("pinfetch: tarball entry %q is not a regular file (type %q); refused", name, hdr.Typeflag)
	}
	if name != "ailang" {
		return fmt.Errorf("pinfetch: tarball entry %q is not the one `ailang` binary; refused", name)
	}
	return nil
}

// createExclusive makes a new 0600 file in dir with O_EXCL.
func createExclusive(dir, prefix string) (*os.File, string, error) {
	var buf [8]byte
	for i := 0; i < 16; i++ {
		if _, err := rand.Read(buf[:]); err != nil {
			return nil, "", err
		}
		p := filepath.Join(dir, prefix+hex.EncodeToString(buf[:]))
		f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o600)
		if err == nil {
			return f, p, nil
		}
		if !errors.Is(err, os.ErrExist) {
			return nil, "", fmt.Errorf("pinfetch: create temp file in %s: %w", dir, err)
		}
	}
	return nil, "", fmt.Errorf("pinfetch: could not create a unique temp file in %s", dir)
}

func isHex64(s string) bool {
	if len(s) != 64 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}

// ReleaseChecksum fetches only pin's release `.sha256` (≤ MaxChecksumBytes)
// and returns its digest, for `doctor --online`: a reachability check that
// also compares the release's claim with the compiled-in pin. It installs
// nothing and downloads no tarball.
func ReleaseChecksum(ctx context.Context, src Source, pin Pin) (string, error) {
	return fetchChecksum(ctx, src, pin)
}
