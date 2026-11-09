package approval

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func fixtureRequest() MintRequest {
	return MintRequest{
		Schema: RequestSchema, Episode: "fixture-ep",
		WorldDB: "/ws/state/world.db", WorkspaceRoot: "/ws/workspace", Worktree: "/ws/workspace/fixture-ep",
		Repo: "/src/ailang-world", Git: "/usr/bin/git", Branch: "world/fixture-ep",
		Grants:     []Grant{{Effect: "Workspace.Read", Scope: "fixture-ep", Budget: 0}, {Effect: "Workspace.Write", Scope: "fixture-ep", Budget: 25}},
		TTLSeconds: 3600, Profile: ProfileRef{Project: "ailang-world", Digest: strings.Repeat("ab", 32)},
		ProfilesDir: "/ws/profiles", CredentialsDB: "/ws/state/credentials.db",
		SessionFile: "/ws/sessions/fixture-ep/session.token", MCPConfig: "/ws/sessions/fixture-ep/mcp.json",
		Agent: "claude", Model: "claude-opus<5>&x", Task: "/ws/tasks/fixture-ep.md",
		Daemon: DaemonRef{Addr: "127.0.0.1:7377", VCSRevision: "0123456789abcdef0123456789abcdef01234567",
			VCSModified: false, ExecutableSHA256: strings.Repeat("cd", 32), GoVersion: "go1.26.6"},
	}
}

const fixtureCanonicalJSON = `{"schema":"world/mint-request/v1","episode":"fixture-ep","world_db":"/ws/state/world.db","workspace_root":"/ws/workspace","worktree":"/ws/workspace/fixture-ep","reuse":false,"degraded":false,"repo":"/src/ailang-world","git":"/usr/bin/git","branch":"world/fixture-ep","grants":[{"effect":"Workspace.Read","scope":"fixture-ep","budget":0},{"effect":"Workspace.Write","scope":"fixture-ep","budget":25}],"ttl_s":3600,"profile":{"project":"ailang-world","digest":"abababababababababababababababababababababababababababababababab"},"profiles_dir":"/ws/profiles","credentials_db":"/ws/state/credentials.db","session_file":"/ws/sessions/fixture-ep/session.token","mcp_config":"/ws/sessions/fixture-ep/mcp.json","agent":"claude","model":"claude-opus<5>&x","task":"/ws/tasks/fixture-ep.md","daemon":{"addr":"127.0.0.1:7377","vcs_revision":"0123456789abcdef0123456789abcdef01234567","vcs_modified":false,"executable_sha256":"cdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcd","go_version":"go1.26.6"}}`

const fixtureDigestHex = "0ace927ef952f541b1e454b066f6d6b1941fb5f3d710c784c22dd80a53d771eb"

const fixtureMessageHashHex = "29181d224d982f659d1bd4703238d3f8489e26968ffe2a22bd786fa3bf15fef5"

func fixtureDigest(t *testing.T) [32]byte {
	t.Helper()
	b, err := hex.DecodeString(fixtureDigestHex)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 32 {
		t.Fatalf("fixtureDigestHex is %d bytes, want 32", len(b))
	}
	var d [32]byte
	copy(d[:], b)
	return d
}

// eachLeaf calls fn for every string, bool and int64 leaf of v, in field
// order, slice elements included. Paths are JSON-style.
func eachLeaf(v reflect.Value, path string, fn func(v reflect.Value, path string)) {
	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			p := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]
			if path != "" {
				p = path + "." + p
			}
			eachLeaf(v.Field(i), p, fn)
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			eachLeaf(v.Index(i), fmt.Sprintf("%s[%d]", path, i), fn)
		}
	case reflect.String, reflect.Bool, reflect.Int64:
		fn(v, path)
	}
}

func TestCanonicalRequestGolden(t *testing.T) {
	got, err := Canonical(fixtureRequest())
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != fixtureCanonicalJSON {
		t.Fatalf("canonical bytes differ:\n got %s\nwant %s", got, fixtureCanonicalJSON)
	}
	if bytes.HasSuffix(got, []byte("\n")) {
		t.Fatal("canonical bytes end in a newline")
	}
	if !bytes.Contains(got, []byte("claude-opus<5>&x")) {
		t.Fatal("canonical bytes escaped the HTML characters in model")
	}
	d, err := Digest(fixtureRequest())
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(d[:]) != fixtureDigestHex {
		t.Fatalf("digest is %x, want %s", d, fixtureDigestHex)
	}
}

func TestDecodeIsStrict(t *testing.T) {
	canon := []byte(fixtureCanonicalJSON)
	got, err := Decode(canon)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, fixtureRequest()) {
		t.Fatalf("decoded request differs from fixtureRequest:\n%s", Diff(got, fixtureRequest()))
	}
	re, err := Canonical(got)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(re, canon) {
		t.Fatal("re-encoding the decoded request changed the bytes")
	}

	indented := new(bytes.Buffer)
	if err := json.Indent(indented, canon, "", "  "); err != nil {
		t.Fatal(err)
	}
	swapped := strings.Replace(fixtureCanonicalJSON, `"reuse":false,"degraded":false`, `"degraded":false,"reuse":false`, 1)
	if swapped == fixtureCanonicalJSON {
		t.Fatal("the reorder replacement changed nothing")
	}
	upperSchema := strings.Replace(fixtureCanonicalJSON, `"schema":`, `"Schema":`, 1)
	if upperSchema == fixtureCanonicalJSON {
		t.Fatal("the case replacement changed nothing")
	}

	arms := []struct {
		name  string
		in    []byte
		exact bool // true: must be exactly ErrNotCanonical
	}{
		{"unknown-field", []byte(`{"extra":1,` + fixtureCanonicalJSON[1:]), false},
		{"trailing-newline", append(append([]byte{}, canon...), '\n'), true},
		{"indented", indented.Bytes(), true},
		{"reordered-keys", []byte(swapped), true},
		{"twice", append(append([]byte{}, canon...), canon...), true},
		{"case-folded-key", []byte(upperSchema), true},
		{"empty", []byte{}, false},
	}
	for _, arm := range arms {
		r, err := Decode(arm.in)
		if err == nil {
			t.Errorf("%s: Decode accepted the input", arm.name)
			continue
		}
		if !reflect.DeepEqual(r, MintRequest{}) {
			t.Errorf("%s: Decode returned a non-zero request with an error", arm.name)
		}
		if arm.exact {
			if !errors.Is(err, ErrNotCanonical) {
				t.Errorf("%s: err = %v, want ErrNotCanonical", arm.name, err)
			}
		} else if !(errors.Is(err, ErrNotCanonical) || errors.Is(err, ErrDecode)) {
			t.Errorf("%s: err = %v, want ErrNotCanonical or ErrDecode", arm.name, err)
		}
	}
}

func TestDigestCoversEveryField(t *testing.T) {
	base, err := Digest(fixtureRequest())
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	eachLeaf(reflect.ValueOf(fixtureRequest()), "", func(_ reflect.Value, p string) {
		paths = append(paths, p)
	})
	if len(paths) < 20 {
		t.Fatalf("only %d leaves found, want >= 20", len(paths))
	}
	t.Logf("digest covers %d leaves", len(paths))
	for i, p := range paths {
		r := fixtureRequest() // a fresh copy per leaf, so Grants is never shared
		idx := 0
		eachLeaf(reflect.ValueOf(&r).Elem(), "", func(v reflect.Value, _ string) {
			if idx == i {
				switch v.Kind() {
				case reflect.String:
					v.SetString(v.String() + "x")
				case reflect.Bool:
					v.SetBool(!v.Bool())
				case reflect.Int64:
					v.SetInt(v.Int() + 1)
				}
			}
			idx++
		})
		d, err := Digest(r)
		if err != nil {
			t.Fatal(err)
		}
		if d == base {
			t.Errorf("changing %s did not change the digest", p)
		}
	}
}

func TestValidateRefusals(t *testing.T) {
	codex := fixtureRequest()
	codex.Agent = "codex"
	codex.MCPConfig = ""
	if err := codex.Validate(); err != nil {
		t.Fatalf("codex variant refused: %v", err)
	}
	reuse := fixtureRequest()
	reuse.Reuse = true
	reuse.Git = ""
	if err := reuse.Validate(); err != nil {
		t.Fatalf("reuse variant refused: %v", err)
	}
	if err := fixtureRequest().Validate(); err != nil {
		t.Fatalf("fixture refused: %v", err)
	}

	arms := []struct {
		name, path string
		mutate     func(*MintRequest)
		want       error
	}{
		{"schema", "schema", func(r *MintRequest) { r.Schema = "world/mint-request/v2" }, ErrSchema},
		{"utf8", "model", func(r *MintRequest) { r.Model = "\xff" }, ErrString},
		{"nul-branch", "branch", func(r *MintRequest) { r.Branch = "a\x00b" }, ErrString},
		{"nul-grant", "grants[0].scope", func(r *MintRequest) { r.Grants[0].Scope = "a\x00" }, ErrString},
		{"nul-daemon", "daemon.addr", func(r *MintRequest) { r.Daemon.Addr = "1\x00" }, ErrString},
		{"nul-path", "world_db", func(r *MintRequest) { r.WorldDB = "/ws/\x00" }, ErrString},
		{"episode-case", "episode", func(r *MintRequest) { r.Episode = "Bad_Ep" }, ErrEpisode},
		{"episode-empty", "episode", func(r *MintRequest) { r.Episode = "" }, ErrEpisode},
		{"episode-long", "episode", func(r *MintRequest) { r.Episode = strings.Repeat("a", 65) }, ErrEpisode},
		{"world-db-rel", "world_db", func(r *MintRequest) { r.WorldDB = "ws/state/world.db" }, ErrPath},
		{"world-db-empty", "world_db", func(r *MintRequest) { r.WorldDB = "" }, ErrPath},
		{"ws-root-unclean", "workspace_root", func(r *MintRequest) { r.WorkspaceRoot = "/ws/workspace/" }, ErrPath},
		{"worktree-dotdot", "worktree", func(r *MintRequest) { r.Worktree = "/ws/workspace/../x" }, ErrPath},
		{"profiles-dir-rel", "profiles_dir", func(r *MintRequest) { r.ProfilesDir = "profiles" }, ErrPath},
		{"creds-rel", "credentials_db", func(r *MintRequest) { r.CredentialsDB = "credentials.db" }, ErrPath},
		{"session-rel", "session_file", func(r *MintRequest) { r.SessionFile = "s.token" }, ErrPath},
		{"repo-rel", "repo", func(r *MintRequest) { r.Repo = "src/ailang-world" }, ErrPath},
		{"git-rel", "git", func(r *MintRequest) { r.Git = "git" }, ErrPath},
		{"mcp-rel", "mcp_config", func(r *MintRequest) { r.MCPConfig = "mcp.json" }, ErrPath},
		{"task-rel", "task", func(r *MintRequest) { r.Task = "task.md" }, ErrPath},
		{"git-missing", "git", func(r *MintRequest) { r.Git = "" }, ErrGitReuse},
		{"git-on-reuse", "git", func(r *MintRequest) { r.Reuse = true }, ErrGitReuse},
		{"no-grants-nil", "grants", func(r *MintRequest) { r.Grants = nil }, ErrNoGrants},
		{"no-grants-empty", "grants", func(r *MintRequest) { r.Grants = []Grant{} }, ErrNoGrants},
		{"effect-empty", "grants[1]", func(r *MintRequest) { r.Grants[1].Effect = "" }, ErrGrantEmpty},
		{"scope-empty", "grants[0]", func(r *MintRequest) { r.Grants[0].Scope = "" }, ErrGrantEmpty},
		{"budget-neg", "grants[1].budget", func(r *MintRequest) { r.Grants[1].Budget = -1 }, ErrNegative},
		{"ttl-neg", "ttl_s", func(r *MintRequest) { r.TTLSeconds = -1 }, ErrNegative},
		{"half-profile-a", "profile", func(r *MintRequest) { r.Profile.Project = "" }, ErrHalfProfile},
		{"half-profile-b", "profile", func(r *MintRequest) { r.Profile.Digest = "" }, ErrHalfProfile},
		{"agent", "agent", func(r *MintRequest) { r.Agent = "gemini" }, ErrAgent},
		{"mcp-for-codex", "mcp_config", func(r *MintRequest) { r.Agent = "codex" }, ErrMCPConfig},
		{"mcp-missing", "mcp_config", func(r *MintRequest) { r.MCPConfig = "" }, ErrMCPConfig},
	}
	for _, arm := range arms {
		r := fixtureRequest()
		arm.mutate(&r)
		err := r.Validate()
		if !errors.Is(err, arm.want) {
			t.Errorf("%s: err = %v, want %v", arm.name, err, arm.want)
			continue
		}
		if !strings.Contains(err.Error(), arm.path) {
			t.Errorf("%s: err = %q, does not name path %q", arm.name, err, arm.path)
		}
	}
}

func TestSignedMessageGolden(t *testing.T) {
	m, err := SignedMessage(fixtureDigest(t), fixtureChallenge(), fixtureExpiry)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(m) != fixtureMessageHex {
		t.Fatalf("signed message is %x, want %s", m, fixtureMessageHex)
	}
	if len(m) != 90 {
		t.Fatalf("signed message is %d bytes, want 90", len(m))
	}
	if string(m[:17]) != MessagePrefix {
		t.Fatalf("prefix is %q, want %q", m[:17], MessagePrefix)
	}
	if m[17] != 0 {
		t.Fatal("byte 17 is not NUL")
	}
	d := fixtureDigest(t)
	if !bytes.Equal(m[18:50], d[:]) {
		t.Fatal("bytes 18:50 are not the digest")
	}
	c := fixtureChallenge()
	if !bytes.Equal(m[50:82], c[:]) {
		t.Fatal("bytes 50:82 are not the challenge")
	}
	if binary.BigEndian.Uint64(m[82:90]) != uint64(fixtureExpiry) {
		t.Fatal("bytes 82:90 are not the big-endian expiry")
	}
	mh, err := MessageHash(fixtureDigest(t), fixtureChallenge(), fixtureExpiry)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(mh[:]) != fixtureMessageHashHex {
		t.Fatalf("message hash is %x, want %s", mh, fixtureMessageHashHex)
	}
	if want := sha256.Sum256(m); mh != want {
		t.Fatal("message hash is not sha256 of the signed message")
	}
	if _, err := SignedMessage(fixtureDigest(t), fixtureChallenge(), -1); !errors.Is(err, ErrNegativeExpiry) {
		t.Fatalf("negative expiry: err = %v, want ErrNegativeExpiry", err)
	}
	if _, err := SignedMessage(fixtureDigest(t), fixtureChallenge(), 0); err != nil {
		t.Fatalf("expiry 0 refused: %v", err)
	}
}

func TestDiffNamesTheField(t *testing.T) {
	f := fixtureRequest()
	if d := Diff(f, f); d != "" {
		t.Fatalf("Diff(f, f) = %q, want empty", d)
	}
	cases := []struct {
		name, want string
		mutate     func(*MintRequest)
	}{
		{"profile-digest", "profile.digest", func(r *MintRequest) { r.Profile.Digest = strings.Repeat("ef", 32) }},
		{"grant-budget", "grants[1].budget", func(r *MintRequest) { r.Grants[1].Budget++ }},
		{"daemon-exe", "daemon.executable_sha256", func(r *MintRequest) { r.Daemon.ExecutableSHA256 = strings.Repeat("ef", 32) }},
		{"grant-appended", "grants", func(r *MintRequest) {
			r.Grants = append(r.Grants, Grant{Effect: "Workspace.Write", Scope: "fixture-ep", Budget: 1})
		}},
		{"episode-and-model", "episode", func(r *MintRequest) {
			r.Episode = "other-ep"
			r.Model = "other-model"
		}},
	}
	for _, c := range cases {
		g := fixtureRequest()
		c.mutate(&g)
		if d := Diff(f, g); d != c.want {
			t.Errorf("%s: Diff = %q, want %q", c.name, d, c.want)
		}
	}
}
