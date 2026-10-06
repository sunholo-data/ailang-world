package main

import (
	"bytes"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/store"
)

// doctorFixture is a scratch HOME with both fixture pins installed by setup,
// a TTY observation that admits both fences, and nothing at --addr.
func doctorFixture(t *testing.T) *setupFixture {
	t.Helper()
	f := newSetupFixture(t)
	if code, out, errOut := runSetupCLI(t); code != exitOK {
		t.Fatalf("fixture setup rc %d\n%s%s", code, out, errOut)
	}
	old := observeTTY
	t.Cleanup(func() { observeTTY = old })
	observeTTY = func() ttyObservation { return ttyObservation{} }
	return f
}

func runDoctorCLI(t *testing.T, addr string, args ...string) (int, string, string) {
	t.Helper()
	var out, errw bytes.Buffer
	code := run(append([]string{"--addr", addr, "doctor"}, args...), &out, &errw)
	return code, out.String(), errw.String()
}

// deadAddr is a loopback URL nothing listens on.
func deadAddr(t *testing.T) string {
	srv := httptest.NewServer(http.NotFoundHandler())
	u := srv.URL
	srv.Close()
	return u
}

type fileState struct {
	names []string
	sizes map[string]int64
}

func snapshot(t *testing.T, dir string) fileState {
	t.Helper()
	st := fileState{sizes: map[string]int64{}}
	_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || p == dir {
			return nil
		}
		rel, _ := filepath.Rel(dir, p)
		st.names = append(st.names, rel)
		if info, err := d.Info(); err == nil && !d.IsDir() {
			st.sizes[rel] = info.Size()
		}
		return nil
	})
	sort.Strings(st.names)
	return st
}

func (a fileState) equal(b fileState) bool {
	if strings.Join(a.names, "\n") != strings.Join(b.names, "\n") {
		return false
	}
	for k, v := range a.sizes {
		if b.sizes[k] != v {
			return false
		}
	}
	return true
}

func storeLine(out string) string {
	var lines []string
	for _, l := range strings.Split(out, "\n") {
		if strings.Contains(l, " store ") || strings.HasPrefix(strings.TrimLeft(l, "✓!✗ "), "store") {
			lines = append(lines, l)
		}
	}
	return strings.Join(lines, "\n")
}

// AC4.1: lock held vs free is reported correctly, and doctor creates no
// file — not on a missing store, not beside a free one, not beside a held one.
// Kills MUT-DOCTOR-OPENS-STORE (store.Open in place of WriterLockHeld or of
// OpenReadOnly).
func TestDoctorLockProbeCreatesNothing(t *testing.T) {
	doctorFixture(t)
	addr := deadAddr(t)

	t.Run("missing-store", func(t *testing.T) {
		dir := t.TempDir()
		db := filepath.Join(dir, "world.db")
		before := snapshot(t, dir)
		_, out, _ := runDoctorCLI(t, addr, "--db", db)
		if !strings.Contains(out, "does not exist yet (doctor never creates it)") {
			t.Fatalf("missing store not reported:\n%s", out)
		}
		if after := snapshot(t, dir); !before.equal(after) {
			t.Fatalf("doctor created files beside a missing store: before %v after %v", before.names, after.names)
		}
	})

	dir := t.TempDir()
	db := filepath.Join(dir, "world.db")
	st, err := store.Open(db)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	t.Run("free", func(t *testing.T) {
		before := snapshot(t, dir)
		_, out, _ := runDoctorCLI(t, addr, "--db", db)
		if !strings.Contains(out, "✓ store      writer lock free") {
			t.Fatalf("free lock not reported free:\n%s", storeLine(out))
		}
		if !strings.Contains(out, "no world head yet") || !strings.Contains(out, "no transition registry") {
			t.Fatalf("an empty store's head/registry not reported:\n%s", storeLine(out))
		}
		if after := snapshot(t, dir); !before.equal(after) {
			t.Fatalf("doctor changed the store dir: before %v after %v", before, after)
		}
	})

	t.Run("held", func(t *testing.T) {
		w, err := store.Open(db)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = w.Close() }()
		before := snapshot(t, dir)
		_, out, _ := runDoctorCLI(t, addr, "--db", db)
		if !strings.Contains(out, "! store      writer lock held by another process") {
			t.Fatalf("held lock not reported held:\n%s", storeLine(out))
		}
		if strings.Contains(storeLine(out), "✗") {
			t.Fatalf("a held lock made a store check fail (doctor must read beside a live writer):\n%s", storeLine(out))
		}
		if after := snapshot(t, dir); !before.equal(after) {
			t.Fatalf("doctor changed the store dir: before %v after %v", before, after)
		}
	})
}

// AC4.2, restated 2026-10-06: the attended steps read their confirmation from
// /dev/tty itself, so ONE terminal fact is left to report — whether /dev/tty
// opens. The publish-fence reasons this used to report (stdin-not-a-terminal,
// stdin-is-not-the-controlling-terminal, and their `</dev/tty` fix) no longer
// exist in world-publish, so doctor must not send an operator after them.
// Each row is the other's control.
func TestDoctorTTYReasons(t *testing.T) {
	doctorFixture(t)
	addr := deadAddr(t)
	for _, tc := range []struct {
		name   string
		obs    ttyObservation
		want   []string
		absent []string
	}{
		{"no-ctty", ttyObservation{cttyErr: os.ErrNotExist},
			[]string{"! tty", "no controlling terminal", "world-publish and session new/mint will refuse here",
				"fix: run the attended steps yourself from a terminal window (an IDE terminal pane works); an agent cannot run them"},
			[]string{"✓ tty"}},
		{"ctty-opens", ttyObservation{},
			[]string{"✓ tty        /dev/tty opens: world-publish and session new/mint read your confirmation from it"},
			[]string{"! tty", "</dev/tty", "publish fence", "stdin-is-not-the-controlling-terminal"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			observeTTY = func() ttyObservation { return tc.obs }
			_, out, _ := runDoctorCLI(t, addr)
			for _, w := range tc.want {
				if !strings.Contains(out, w) {
					t.Errorf("doctor output lacks %q:\n%s", w, out)
				}
			}
			for _, a := range tc.absent {
				if strings.Contains(out, a) {
					t.Errorf("doctor output carries %q:\n%s", a, out)
				}
			}
		})
	}
}

// AC4.3: a flipped pin byte is ✗ with both digests, and doctor exits 1.
func TestDoctorFlippedPinByte(t *testing.T) {
	f := doctorFixture(t)
	addr := deadAddr(t)
	code, out, _ := runDoctorCLI(t, addr)
	if code != exitOK || strings.Contains(out, "✗") {
		t.Fatalf("control: a fresh setup should leave no ✗ (rc %d):\n%s", code, out)
	}
	tool := filepath.Join(f.home, ".pinned-ailang-tools", toolRelease(), "ailang")
	b, _ := os.ReadFile(tool)
	b[0] ^= 0x01
	if err := os.WriteFile(tool, b, 0o755); err != nil {
		t.Fatal(err)
	}
	code, out, _ = runDoctorCLI(t, addr)
	pin, _ := lookupPin(toolRelease(), setupPlatform)
	if code != exitUsage {
		t.Fatalf("rc %d, want 1:\n%s", code, out)
	}
	for _, want := range []string{"✗ pins", sha256Hex(b), pin.BinarySHA256, "setup --replace"} {
		if !strings.Contains(out, want) {
			t.Errorf("doctor output lacks %q:\n%s", want, out)
		}
	}
}

func TestDoctorNotesStaleFilesBesideAPin(t *testing.T) {
	f := doctorFixture(t)
	_ = os.WriteFile(filepath.Join(f.home, ".pinned-ailang", "darwin.arm64.ailang.tar.gz"), []byte("old"), 0o644)
	_, out, _ := runDoctorCLI(t, deadAddr(t))
	if !strings.Contains(out, "! pins       interpreter: stale files beside the pin: darwin.arm64.ailang.tar.gz") {
		t.Fatalf("stale tarball not noted:\n%s", out)
	}
}

// AC4.4: a worldd listener is told apart from a foreign one and from none.
func TestDoctorDaemonVsForeignListener(t *testing.T) {
	f := doctorFixture(t)
	db := filepath.Join(f.home, ".ailang", "world", "world.db")
	worldd := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/health" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok","daemon_version":"0.1.0","db_path":"` + db + `","interpreter_version":"AILANG v0.41.0"}`))
	}))
	defer worldd.Close()
	foreign := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("<html>some other service</html>"))
	}))
	defer foreign.Close()

	code, out, _ := runDoctorCLI(t, worldd.URL)
	if !strings.Contains(out, "✓ daemon     ailang-worldd 0.1.0 at "+worldd.URL+", store "+db) {
		t.Fatalf("worldd not recognised (rc %d):\n%s", code, out)
	}
	code, out, _ = runDoctorCLI(t, foreign.URL)
	if code != exitUsage || !strings.Contains(out, "✗ daemon     something that is not ailang-worldd answers at "+foreign.URL) {
		t.Fatalf("foreign listener not flagged (rc %d):\n%s", code, out)
	}
	_, out, _ = runDoctorCLI(t, deadAddr(t))
	if !strings.Contains(out, "! daemon     nothing answers at") {
		t.Fatalf("an empty port not reported:\n%s", out)
	}
}

func TestDoctorWorkspaceAndExamples(t *testing.T) {
	f := doctorFixture(t)
	ws := filepath.Join(f.home, ".ailang", "world-ws")
	// one real-shaped worktree, one plain dir, one bad name
	_ = os.MkdirAll(filepath.Join(ws, "ep1"), 0o755)
	_ = os.WriteFile(filepath.Join(ws, "ep1", ".git"), []byte("gitdir: /x/.git/worktrees/ep1\n"), 0o644)
	_ = os.MkdirAll(filepath.Join(ws, "plain"), 0o755)
	_ = os.MkdirAll(filepath.Join(ws, "Bad_Name"), 0o755)
	ex := filepath.Join(f.home, "examples")
	_ = os.MkdirAll(filepath.Join(ex, "sub"), 0o755)
	_ = os.WriteFile(filepath.Join(ex, "a.ail"), nil, 0o644)
	_ = os.WriteFile(filepath.Join(ex, "sub", "b.ail"), nil, 0o644)
	_, out, _ := runDoctorCLI(t, deadAddr(t), "--examples-dir", ex)
	for _, want := range []string{
		"1 episode worktree(s): ep1",
		"Bad_Name (not a valid episode name)",
		"plain (not a git worktree)",
		"✓ examples   2 .ail files in " + ex,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("doctor output lacks %q:\n%s", want, out)
		}
	}
}

// The doctor's episode grammar is the daemon's.
func TestDoctorEpisodeGrammarMatchesTheDaemon(t *testing.T) {
	src, err := os.ReadFile(filepath.Join("..", "..", "host", "daemon", "workspace.go"))
	if err != nil {
		t.Fatal(err)
	}
	want := "regexp.MustCompile(`" + workspaceEpisodePattern.String() + "`)"
	if !strings.Contains(string(src), "episodeIDPattern = "+want) {
		t.Fatalf("host/daemon/workspace.go's episodeIDPattern is not %s", want)
	}
}

func TestDoctorHelpAndAPIKeyGuard(t *testing.T) {
	var out, errw bytes.Buffer
	if code := run([]string{"doctor", "--help"}, &out, &errw); code != exitOK || !strings.Contains(out.String(), "usage: ailang-worldd [--addr <url>] doctor") {
		t.Fatalf("doctor --help rc %d: %q", code, out.String())
	}
	// The API-key finding is the guard itself: doctor is not exempt.
	t.Setenv("AILANG_REGISTRY_API_KEY", "sentinel-not-a-key")
	out.Reset()
	errw.Reset()
	code := run([]string{"doctor"}, &out, &errw)
	if code != exitFatal || !strings.Contains(errw.String(), "AILANG_REGISTRY_API_KEY") || strings.Contains(errw.String(), "sentinel-not-a-key") {
		t.Fatalf("doctor with the key set: rc %d stderr %q", code, errw.String())
	}
}
