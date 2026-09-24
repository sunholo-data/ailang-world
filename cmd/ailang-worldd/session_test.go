package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/store"
)

// fakeTerm is an io.ReadWriteCloser whose Read yields the y/N confirmation
// answer and whose Write is discarded. It stands in for /dev/tty so the mint
// fence + confirmation can be driven from a test (the driven tty-fencing seam).
type fakeTerm struct{ r io.Reader }

func (f *fakeTerm) Read(p []byte) (int, error)  { return f.r.Read(p) }
func (f *fakeTerm) Write(p []byte) (int, error) { return len(p), nil }
func (f *fakeTerm) Close() error                { return nil }

// testMintEnv returns a sessionEnv that reports a controlling terminal whose
// confirmation reader yields `answer`, with a fixed clock.
func testMintEnv(answer string, now int64) sessionEnv {
	return sessionEnv{
		openTerminal: func() (io.ReadWriteCloser, error) { return &fakeTerm{r: strings.NewReader(answer)}, nil },
		now:          func() int64 { return now },
	}
}

// noTerminalEnv is the fence refusal seam: opening /dev/tty fails, so any mint
// under it must refuse.
func noTerminalEnv(now int64) sessionEnv {
	return sessionEnv{
		openTerminal: func() (io.ReadWriteCloser, error) { return nil, os.ErrNotExist },
		now:          func() int64 { return now },
	}
}

const hex64Re = `^[0-9a-f]{64}$`

var hex64TokenRegex = regexp.MustCompile(`\b[0-9a-f]{64}\b`)

func dbPathOf(t *testing.T) string { return filepath.Join(t.TempDir(), "world.db") }

// TestSessionMint_PrintsOnceExactly is AC-M1-1: under a driven tty-fencing seam
// the mint prints the raw 64-hex token to stdout EXACTLY once (never logged
// twice, never echoed to stderr) and persists one row. The credential_id hash
// goes to stderr, so stdout carries the single token and nothing else hex.
func TestSessionMint_PrintsOnceExactly(t *testing.T) {
	db := dbPathOf(t)
	var stdout, stderr bytes.Buffer
	code := runSessionMint([]string{
		"--db", db, "--episode", "e1", "--grant", "fs.read=/tmp/b:1", "--ttl", "3600",
	}, &stdout, &stderr, testMintEnv("y\n", 1000))
	if code != exitOK {
		t.Fatalf("mint exit=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	tokens := hex64TokenRegex.FindAllString(stdout.String(), -1)
	if len(tokens) != 1 {
		t.Fatalf("raw token appears %d times in stdout %q, want exactly 1", len(tokens), stdout.String())
	}
	tok := tokens[0]
	if !regexp.MustCompile(hex64Re).MatchString(tok) {
		t.Fatalf("extracted token %q is not 64 lowercase hex", tok)
	}
	// stderr carries the credential_id diagnostic (the hash) but never the raw
	// token.
	if strings.Contains(stderr.String(), tok) {
		t.Fatalf("stderr leaked the raw token: %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "credential_id=") {
		t.Fatalf("stderr does not name credential_id for revocation: %q", stderr.String())
	}

	// Exactly one row persisted, and it resolves.
	st, err := store.Open(db)
	if err != nil {
		t.Fatalf("reopen store: %v", err)
	}
	defer func() { _ = st.Close() }()
	out := authority.New(st).Resolve("Bearer "+tok, 1001)
	if out.Success == nil || out.Success.EpisodeID != "e1" {
		t.Fatalf("minted token did not resolve to episode e1: %#v", out)
	}
}

// TestSessionMint_DistinctCredentials is AC-M1-2: two mints for the same
// episode produce distinct raw tokens.
func TestSessionMint_DistinctCredentials(t *testing.T) {
	db := dbPathOf(t)
	mint := func() string {
		var stdout, stderr bytes.Buffer
		code := runSessionMint([]string{
			"--db", db, "--episode", "e1", "--grant", "fs.read=/tmp/b:1", "--ttl", "3600",
		}, &stdout, &stderr, testMintEnv("y\n", 1000))
		if code != exitOK {
			t.Fatalf("mint exit=%d stderr=%s", code, stderr.String())
		}
		toks := hex64TokenRegex.FindAllString(stdout.String(), -1)
		if len(toks) != 1 {
			t.Fatalf("expected exactly one token, got %d: %q", len(toks), stdout.String())
		}
		return toks[0]
	}
	a, b := mint(), mint()
	if a == b {
		t.Fatalf("two mints for the same episode produced the same credential: %s", a)
	}
}

// TestSessionMint_RequiresControllingTerminal is the D1 tty fence: a mint with
// no controlling terminal refuses (usage) and prints nothing.
func TestSessionMint_RequiresControllingTerminal(t *testing.T) {
	db := dbPathOf(t)
	var stdout, stderr bytes.Buffer
	code := runSessionMint([]string{
		"--db", db, "--episode", "e1", "--grant", "fs.read=/tmp/b:1", "--ttl", "3600",
	}, &stdout, &stderr, noTerminalEnv(1000))
	if code != exitUsage {
		t.Fatalf("no-terminal mint exit=%d, want %d (usage refusal)", code, exitUsage)
	}
	if stdout.String() != "" {
		t.Fatalf("no-terminal mint printed a credential to stdout: %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "controlling terminal") {
		t.Fatalf("no-terminal refusal does not name the controlling terminal: %q", stderr.String())
	}
}

// TestSessionMint_ConfirmationRejected: a "n" answer aborts the mint — no
// credential is minted or printed.
func TestSessionMint_ConfirmationRejected(t *testing.T) {
	db := dbPathOf(t)
	var stdout, stderr bytes.Buffer
	code := runSessionMint([]string{
		"--db", db, "--episode", "e1", "--grant", "fs.read=/tmp/b:1", "--ttl", "3600",
	}, &stdout, &stderr, testMintEnv("n\n", 1000))
	if code != exitUsage {
		t.Fatalf("declined mint exit=%d, want %d", code, exitUsage)
	}
	if stdout.String() != "" {
		t.Fatalf("declined mint printed a credential: %q", stdout.String())
	}
}

// TestSessionRevokeCLI is AC-M2-5's CLI half: `session revoke <credential_id>`.
func TestSessionRevokeCLI(t *testing.T) {
	db := dbPathOf(t)
	var mintOut, mintErr bytes.Buffer
	if code := runSessionMint([]string{
		"--db", db, "--episode", "e1", "--grant", "fs.read=/tmp/b:1", "--ttl", "3600",
	}, &mintOut, &mintErr, testMintEnv("y\n", 1000)); code != exitOK {
		t.Fatalf("mint exit=%d stderr=%s", code, mintErr.String())
	}
	// Extract the credential_id from the stderr diagnostic.
	credLine := ""
	for _, line := range strings.Split(mintErr.String(), "\n") {
		if strings.HasPrefix(line, "ailang-worldd session mint: credential_id=") {
			credLine = strings.TrimSuffix(line, "\n")
			break
		}
	}
	if credLine == "" {
		t.Fatalf("mint stderr lacks credential_id: %q", mintErr.String())
	}
	// The credential_id is exactly the 64-hex hash on that diagnostic line
	// (the rest of the line is the human summary).
	credID := hex64TokenRegex.FindString(credLine)
	if credID == "" {
		t.Fatalf("credential_id line %q has no 64-hex hash", credLine)
	}
	var out, errOut bytes.Buffer
	code := runSessionRevoke([]string{"--db", db, credID}, &out, &errOut)
	if code != exitOK {
		t.Fatalf("revoke exit=%d stderr=%s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "revoked") {
		t.Fatalf("revoke stdout %q does not confirm", out.String())
	}

	// The token now resolves as UNKNOWN (delete-from-mapping = unknown, D4).
	tok := hex64TokenRegex.FindString(mintOut.String())
	st, err := store.Open(db)
	if err != nil {
		t.Fatalf("reopen store: %v", err)
	}
	defer func() { _ = st.Close() }()
	outR := authority.New(st).Resolve("Bearer "+tok, 1001)
	if outR.Success != nil || outR.Denied == nil || *outR.Denied != authority.DenialUnknown {
		t.Fatalf("post-revoke resolve = %#v, want DenialUnknown", outR)
	}
}

// TestSessionMint_OutFileModeAndOnce: --out writes the credential to the file
// exactly once at mode 0600 and prints no credential to stdout.
func TestSessionMint_OutFileModeAndOnce(t *testing.T) {
	dir := t.TempDir()
	db := filepath.Join(dir, "world.db")
	outPath := filepath.Join(dir, "credential.txt")
	var stdout, stderr bytes.Buffer
	code := runSessionMint([]string{
		"--db", db, "--episode", "e1", "--grant", "fs.read=/tmp/b:1", "--ttl", "3600", "--out", outPath,
	}, &stdout, &stderr, testMintEnv("y\n", 1000))
	if code != exitOK {
		t.Fatalf("mint exit=%d stderr=%s", code, stderr.String())
	}
	if stdout.String() == "" || !strings.Contains(stdout.String(), outPath) {
		t.Fatalf("--out confirmation stdout %q does not name the path", stdout.String())
	}
	if hex64TokenRegex.MatchString(stdout.String()) {
		t.Fatalf("--out mode leaked the credential to stdout: %q", stdout.String())
	}
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read out file: %v", err)
	}
	tok := strings.TrimSpace(string(data))
	if !regexp.MustCompile(hex64Re).MatchString(tok) {
		t.Fatalf("out file does not hold exactly one 64-hex credential: %q", data)
	}
	if fi, err := os.Stat(outPath); err != nil || fi.Mode()&0o777 != 0o600 {
		t.Fatalf("out file mode = %v, want 0600 (err=%v)", fi.Mode(), err)
	}
	if strings.Contains(stderr.String(), tok) {
		t.Fatalf("stderr leaked the raw token: %q", stderr.String())
	}
}
