package authority

import (
	"context"
	"errors"
	"io"
	"path/filepath"
	"testing"

	"github.com/sunholo-data/ailang-world/host/store"
)

// TestReadOnlyAtSeesMintAndRevokeWithoutRestart is row 189's resolver proof:
// ONE resolver over ReadOnlyAt (what the daemon builds once at startup) sees a
// credential minted after it was built, and stops seeing it on the request
// after a revoke — while the terminal's writer still holds the credentials
// lock. Mutant: cache one read-only handle across requests that snapshots the
// table, and the late mint resolves Unknown.
func TestReadOnlyAtSeesMintAndRevokeWithoutRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), store.CredentialsFileName)
	writer, err := store.OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = writer.Close() }()
	res := New(ReadOnlyAt(path))
	now := int64(1000)

	if out := resolveNow(t, res, "Bearer "+hex64, now); out.Denied == nil || *out.Denied != DenialUnknown {
		t.Fatalf("before any mint = %#v, want Unknown", out)
	}
	tok, id, _, err := Mint(boundedTestContext(t), writer, "ep-late", testGrants(), 3600, now, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if out := resolveNow(t, res, "Bearer "+tok, now); out.Success == nil || out.Success.EpisodeID != "ep-late" {
		t.Fatalf("after the mint = %#v, want ep-late", out)
	}
	if err := Revoke(boundedTestContext(t), writer, id); err != nil {
		t.Fatal(err)
	}
	if out := resolveNow(t, res, "Bearer "+tok, now); out.Denied == nil || *out.Denied != DenialUnknown {
		t.Fatalf("after the revoke = %#v, want Unknown on the next request", out)
	}
}

// A read failure is an error, never a denial (D-WORLD-40), and a
// deadline-free caller is refused before any file is opened.
func TestReadOnlyAtFailuresAreErrors(t *testing.T) {
	missing := ReadOnlyAt(filepath.Join(t.TempDir(), "absent.db"))
	ctx, cancel := context.WithTimeout(context.Background(), 5e9)
	defer cancel()
	if _, err := New(missing).ResolveContext(ctx, "Bearer "+hex64, 1000); err == nil {
		t.Fatal("a missing credentials store resolved without an error")
	}
	if _, _, err := missing.ResolveSession(context.Background(), "x"); !errors.Is(err, store.ErrNoDeadline) {
		t.Fatalf("deadline-free lookup = %v, want ErrNoDeadline", err)
	}
}
