package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

func healthyDoctorAddr(t *testing.T, db string) string {
	t.Helper()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"status": "ok", "daemon_version": "0.1.0", "db_path": db})
	}))
	t.Cleanup(s.Close)
	return s.URL
}

func TestDoctorDaemonDownSaysNothingListening(t *testing.T) {
	addr := deadAddr(t)
	db := filepath.Join(t.TempDir(), "named-world.db")
	code, out, errw := runDoctorCLI(t, addr, "--quick", "--db", db)
	for _, want := range []string{"✗ daemon     nothing is listening at " + addr, "start World: ailang-worldd serve --db " + shellQuote(db)} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q: %s", want, out)
		}
	}
	if code != 10 || strings.Contains(out+errw, "HTTP 0") {
		t.Fatalf("rc %d: %s%s", code, out, errw)
	}
}

func doctorCredentialsFixture(t *testing.T) (string, map[string]string, time.Time) {
	t.Helper()
	dir := t.TempDir()
	db := filepath.Join(dir, "world.db")
	st, err := store.Open(db)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	creds, err := store.OpenCredentials(filepath.Join(dir, "credentials.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer creds.Close()
	now := time.Now().Truncate(time.Second)
	files := map[string]string{}
	ctx, cancel := budgetContext(5 * time.Second)
	defer cancel()
	for _, name := range []string{"valid", "expired", "revoked"} {
		start := now.Unix()
		if name == "expired" {
			start -= 24 * 3600
		}
		grants := []broker.Capability{{Effect: broker.EffectWorkspaceRead, Scope: "worktree", Budget: 17, ExpiresAt: start + 8*3600 + 12*60}}
		token, id, _, err := authority.Mint(ctx, creds, "w188l-"+name, grants, 8*3600+12*60, start, nil)
		if err != nil {
			t.Fatal(err)
		}
		files[name] = filepath.Join(dir, name+".session")
		if err := os.WriteFile(files[name], []byte(token), 0600); err != nil {
			t.Fatal(err)
		}
		if name == "revoked" {
			if err := creds.RevokeSession(ctx, id); err != nil {
				t.Fatal(err)
			}
		}
	}
	return db, files, now
}

func TestDoctorSessionValidExpiredRevoked(t *testing.T) {
	db, files, now := doctorCredentialsFixture(t)
	for _, name := range []string{"valid", "expired", "revoked"} {
		var out bytes.Buffer
		r := &doctorReport{w: &out}
		doctorSession(r, db, "", files[name], now)
		want := map[string]string{
			"valid":   "✓ session    w188l-valid: 8 h 12 m left",
			"expired": "✗ session    w188l-expired's session expired at ",
			"revoked": "✗ session    session is revoked or unknown",
		}[name]
		if !strings.Contains(out.String(), want) {
			t.Errorf("%s: %s", name, out.String())
		}
		if name != "valid" && !strings.Contains(out.String(), "start a new session: ailang-worldd session new") {
			t.Fatal("no renewal command")
		}
		token, _ := os.ReadFile(files[name])
		if strings.Contains(out.String(), string(token)) || strings.Contains(out.String(), fmtHash(token)) {
			t.Fatal("credential leaked")
		}
		if name == "valid" && !strings.Contains(out.String(), "17") {
			t.Fatal("grant budget absent")
		}
	}
	var out bytes.Buffer
	r := &doctorReport{w: &out}
	doctorSession(r, db, "", filepath.Join(t.TempDir(), "missing.session"), now)
	if r.fails != 1 {
		t.Fatal("missing session accepted")
	}
	if err := os.Chmod(files["valid"], 0644); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	r = &doctorReport{w: &out}
	doctorSession(r, db, "", files["valid"], now)
	if r.fails != 1 || !strings.Contains(out.String(), "chmod 0600") {
		t.Fatal(out.String())
	}
}

func TestDoctorExecToolMoved(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	old := filepath.Join(dir, "Cellar", "node", "26.0.0", "bin", "node")
	current := filepath.Join(dir, "Cellar", "node", "26.12.1", "bin", "node")
	for _, path := range []string{old, current} {
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
			t.Fatal(err)
		}
	}
	bin := filepath.Join(dir, "bin")
	if err := os.Mkdir(bin, 0700); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(bin, "node")
	if err := os.Symlink(old, link); err != nil {
		t.Fatal(err)
	}
	profile := filepath.Join(dir, "profile.json")
	b, _ := json.Marshal(map[string]any{"path": []string{bin}, "commands": map[string]any{"test": map[string]any{"argv": []string{"node"}}}})
	if err := os.WriteFile(profile, b, 0600); err != nil {
		t.Fatal(err)
	}
	for _, remove := range []bool{false, true} {
		if err := os.Remove(link); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(current, link); err != nil {
			t.Fatal(err)
		}
		if remove {
			if err := os.Remove(old); err != nil {
				t.Fatal(err)
			}
		}
		var out bytes.Buffer
		r := &doctorReport{w: &out}
		doctorExecTools(r, profile, map[string]string{"node": old})
		for _, want := range []string{"✗ exec tools node: the daemon resolved " + old + " at startup;", "now " + current, "restart World so it resolves the tools again"} {
			if !strings.Contains(out.String(), want) {
				t.Errorf("missing %q: %s", want, out.String())
			}
		}
		if remove && !strings.Contains(out.String(), "that file is gone") {
			t.Fatal(out.String())
		}
	}
}

func TestDoctorQuickIsFast(t *testing.T) {
	db, files, _ := doctorCredentialsFixture(t)
	dir, bin := hostsLauncherFixture(t)
	t.Setenv("PATH", bin+":/usr/bin:/bin")
	addr := healthyDoctorAddr(t, db)
	start := time.Now()
	code, out, errw := runDoctorCLI(t, addr, "--quick", "--json", "--db", db, "--session", files["valid"], "--host", "claude", "--hosts-dir", dir)
	if elapsed := time.Since(start); elapsed >= 2*time.Second {
		t.Fatalf("quick took %s", elapsed)
	}
	if code != 0 {
		t.Fatalf("%d: %s%s", code, out, errw)
	}
	seen := map[string]bool{}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		var obj map[string]json.RawMessage
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			t.Fatal(err)
		}
		check, ok := obj["check"]
		if !ok {
			continue // The final failed-count summary has no check name.
		}
		var name string
		if err := json.Unmarshal(check, &name); err != nil {
			t.Fatal(err)
		}
		switch name {
		case "daemon", "session", "host":
			seen[name] = true
		default:
			t.Fatalf("quick ran %s: %s", name, out)
		}
	}
	for _, name := range []string{"daemon", "session", "host"} {
		if !seen[name] {
			t.Fatalf("quick omitted %s: %s", name, out)
		}
	}
}

func TestDoctorJSONClassesAndExitCodes(t *testing.T) {
	// Report order and priority are independent of warnings and legacy checks.
	for _, tc := range []struct {
		name  string
		code  int
		class string
	}{
		{"daemon", 10, "daemon"}, {"session", 11, "session"}, {"host", 12, "host"}, {"exec tools", 13, "exec"}, {"pins", 1, "other"},
	} {
		var out bytes.Buffer
		r := &doctorReport{w: &out, asJSON: true}
		r.add(markWarn, tc.name, "warning", "repair")
		r.add(markFail, tc.name, "named cause", "exact action")
		if r.classCode != tc.code || r.fails != 1 {
			t.Fatalf("%s: %+v", tc.name, r)
		}
		dec := json.NewDecoder(&out)
		for i := 0; i < 2; i++ {
			var obj map[string]any
			if err := dec.Decode(&obj); err != nil {
				t.Fatal(err)
			}
			if len(obj) != 5 || obj["class"] != tc.class || obj["check"] != tc.name || obj["fix"] == "" {
				t.Fatalf("%+v", obj)
			}
		}
	}
	var out bytes.Buffer
	r := &doctorReport{w: &out, asJSON: true}
	for _, name := range []string{"pins", "exec tools", "host", "session", "daemon"} {
		r.add(markFail, name, "cause", "repair")
	}
	if r.classCode != 10 {
		t.Fatal("first dependency class must win")
	}
	db, files, _ := doctorCredentialsFixture(t)
	addr := healthyDoctorAddr(t, db)
	code, output, errw := runDoctorCLI(t, addr, "--quick", "--json", "--db", db, "--session", files["expired"])
	if code != 11 {
		t.Fatalf("%d: %s%s", code, output, errw)
	}
	for _, tc := range []struct {
		name  string
		addr  string
		flags []string
		code  int
	}{
		{"daemon", deadAddr(t), []string{"--quick", "--session", files["expired"], "--host", "invalid"}, 10},
		{"host", addr, []string{"--quick", "--host", "invalid"}, 12},
		{"exec", addr, []string{"--profile", filepath.Join(t.TempDir(), "absent.json")}, 13},
	} {
		got, text, errors := runDoctorCLI(t, tc.addr, append([]string{"--json", "--db", db}, tc.flags...)...)
		if got != tc.code {
			t.Fatalf("%s: exit %d, want %d: %s%s", tc.name, got, tc.code, text, errors)
		}
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if lines[len(lines)-1] != "{\"failed\":1}" {
		t.Fatal(output)
	}
	for _, line := range lines {
		if !json.Valid([]byte(line)) {
			t.Fatalf("not JSON: %s", line)
		}
	}
}

func TestDoctorStaysReadOnly(t *testing.T) {
	db, files, _ := doctorCredentialsFixture(t)
	dir, bin := hostsLauncherFixture(t)
	t.Setenv("PATH", bin+":/usr/bin:/bin")
	addr := healthyDoctorAddr(t, db)
	profile := filepath.Join(t.TempDir(), "missing-profile.json")
	// Extend TestDoctorLockProbeCreatesNothing to every new read path, with
	// present and absent inputs. The fake host's --version performs no writes.
	root := filepath.Dir(db)
	before := snapshot(t, root)
	hostsBefore := snapshot(t, dir)
	_, _, _ = runDoctorCLI(t, addr, "--db", db, "--session", files["valid"], "--host", "claude", "--hosts-dir", dir, "--profile", profile)
	if !before.equal(snapshot(t, root)) || !hostsBefore.equal(snapshot(t, dir)) {
		t.Fatal("doctor changed files")
	}
	missing := filepath.Join(t.TempDir(), "absent")
	prior := snapshot(t, filepath.Dir(missing))
	_, _, _ = runDoctorCLI(t, addr, "--quick", "--db", filepath.Join(missing, "world.db"), "--session", filepath.Join(missing, "token"), "--host", "pi", "--hosts-dir", filepath.Join(missing, "hosts"))
	if !prior.equal(snapshot(t, filepath.Dir(missing))) {
		t.Fatal("doctor created missing paths")
	}
}

func TestDoctorHostScrubsEnvironmentAndDetectsDrift(t *testing.T) {
	dir, bin := hostsLauncherFixture(t)
	path := filepath.Join(bin, "claude")
	script := "#!/bin/sh\n[ -z \"$WORLD_SESSION$UNRELATED_SECRET$AILANG_REGISTRY_API_KEY\" ] || exit 99\nprintf 'host 1.0\\n'\n"
	if err := os.WriteFile(path, []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+":/usr/bin:/bin")
	t.Setenv("WORLD_SESSION", "never-inherit")
	t.Setenv("UNRELATED_SECRET", "never-inherit")
	t.Setenv("AILANG_REGISTRY_API_KEY", "never-inherit")
	var out bytes.Buffer
	r := &doctorReport{w: &out}
	doctorHost(r, "claude", dir)
	if r.fails != 0 || !strings.Contains(out.String(), "host 1.0") {
		t.Fatal(out.String())
	}
	if err := os.WriteFile(filepath.Join(dir, "claude-world.settings.json"), []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	doctorHost(r, "claude", dir)
	if !strings.Contains(out.String(), "observed, not confined") || !strings.Contains(out.String(), "--replace") {
		t.Fatal(out.String())
	}
}
