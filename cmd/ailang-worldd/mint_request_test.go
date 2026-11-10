package main

// Row 201 A0 (design §2.1a, §2.5): the field map between mintOptions/mintRequest and the stored
// approval.MintRequest, and the round trip that approve relies on.

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

// fieldMapProblems checks m against the mintOptions/mintRequest fields, broker.Capability's fields
// and approval.MintRequest's leaves. It returns nothing when the map is complete.
func fieldMapProblems(m map[string]string) []string {
	var problems []string
	fields := map[string]bool{}
	var walk func(t reflect.Type)
	walk = func(t reflect.Type) {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.Anonymous && f.Type.Kind() == reflect.Struct {
				walk(f.Type)
				continue
			}
			fields[f.Name] = true
		}
	}
	walk(reflect.TypeOf(mintRequest{}))
	capFields := map[string]bool{}
	ct := reflect.TypeOf(broker.Capability{})
	for i := 0; i < ct.NumField(); i++ {
		capFields["grants[]."+ct.Field(i).Name] = true
	}
	for name := range fields {
		if _, ok := m[name]; !ok {
			problems = append(problems, "unmapped field "+name)
		}
	}
	for key := range capFields {
		if _, ok := m[key]; !ok {
			problems = append(problems, "unmapped field "+key)
		}
	}
	for key := range m {
		if !fields[key] && !capFields[key] {
			problems = append(problems, "stale key "+key)
		}
	}
	targets := map[string]bool{}
	reqType := reflect.TypeOf(approval.MintRequest{})
	for key, v := range m {
		switch {
		case strings.HasPrefix(v, "MintRequest."):
			path := strings.TrimPrefix(v, "MintRequest.")
			if !resolvesOn(reqType, path) {
				problems = append(problems, "unresolved "+v+" for "+key)
			}
			targets[path] = true
		case strings.HasPrefix(v, "excluded: "), strings.HasPrefix(v, "fixed: "), strings.HasPrefix(v, "nil: "):
		default:
			problems = append(problems, "bad value for "+key+": "+v)
		}
	}
	leaves := []string{}
	leavesOf(reqType, "", &leaves)
	for _, leaf := range leaves {
		if !targets[leaf] && !inList(requestOnlyLeaves, leaf) {
			problems = append(problems, "unmapped leaf MintRequest."+leaf)
		}
	}
	for _, r := range requestOnlyLeaves {
		if targets[r] {
			problems = append(problems, "request-only leaf "+r+" is also mapped")
		}
	}
	// Non-vacuity: the counts this sprint's design fixes.
	if len(fields) < 29 {
		problems = append(problems, "vacuous: too few mint fields")
	}
	if len(capFields) != 4 {
		problems = append(problems, "vacuous: broker.Capability has a field count other than 4")
	}
	if len(leaves) != 28 {
		problems = append(problems, "vacuous: MintRequest has a leaf count other than 28")
	}
	return problems
}

// resolvesOn reports whether a dotted path names a field of t; a segment ending in [] is a slice
// field whose element type the path continues in.
func resolvesOn(t reflect.Type, path string) bool {
	for _, seg := range strings.Split(path, ".") {
		if t.Kind() != reflect.Struct {
			return false
		}
		if name, isSlice := strings.CutSuffix(seg, "[]"); isSlice {
			f, ok := t.FieldByName(name)
			if !ok || f.Type.Kind() != reflect.Slice {
				return false
			}
			t = f.Type.Elem()
			continue
		}
		f, ok := t.FieldByName(seg)
		if !ok {
			return false
		}
		t = f.Type
	}
	return true
}

// leavesOf lists the leaves of t: structs recurse as Parent.Child, slices of structs as
// Name[].Child, and everything else is a leaf.
func leavesOf(t reflect.Type, prefix string, out *[]string) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := prefix + f.Name
		switch {
		case f.Type.Kind() == reflect.Struct:
			leavesOf(f.Type, name+".", out)
		case f.Type.Kind() == reflect.Slice && f.Type.Elem().Kind() == reflect.Struct:
			leavesOf(f.Type.Elem(), name+"[].", out)
		default:
			*out = append(*out, name)
		}
	}
}

func inList(xs []string, x string) bool {
	for _, s := range xs {
		if s == x {
			return true
		}
	}
	return false
}

// CL14: every field of the mint structs and every Capability field is mapped, and every leaf of
// the stored request is either the target of a mapping or request-only.
func TestMintRequestFieldMapIsComplete(t *testing.T) {
	if problems := fieldMapProblems(mintRequestFieldMap); len(problems) > 0 {
		t.Fatalf("mintRequestFieldMap is incomplete: %v", problems)
	}
	without := map[string]string{}
	for k, v := range mintRequestFieldMap {
		if k != "git" {
			without[k] = v
		}
	}
	if !inList(fieldMapProblems(without), "unmapped field git") {
		t.Errorf("dropping git reports no 'unmapped field git'")
	}
	extra := map[string]string{"nope": "fixed: x"}
	for k, v := range mintRequestFieldMap {
		extra[k] = v
	}
	if !inList(fieldMapProblems(extra), "stale key nope") {
		t.Errorf("an extra key reports no 'stale key nope'")
	}
}

// CL15: a request built from the validated mint, decoded and read back by approve, rebuilds the
// same mint: every field of the mint options is equal, and the MCP config writes the same bytes.
func TestApproveRebuildsTheSameMintRequest(t *testing.T) {
	for _, arm := range []struct{ name, agent, profile string }{{"claude", "claude", "proj"}, {"codex", "codex", ""}} {
		t.Run(arm.name, func(t *testing.T) {
			r := newEpisodeRun(t)
			f := r.f
			db, err := store.CanonicalPath(f.db)
			if err != nil {
				t.Fatal(err)
			}
			cred, err := credentialsFor(db, "")
			if err != nil {
				t.Fatal(err)
			}
			bound := arm.profile != ""
			if bound {
				writeTestProfile(t, r.profiles, arm.profile, "")
			}
			if err := os.MkdirAll(r.sessions, 0o700); err != nil {
				t.Fatal(err)
			}
			grants, err := episodeStartGrants("", nil, 50, bound)
			if err != nil {
				t.Fatal(err)
			}
			mcpFile := ""
			if arm.agent == "claude" {
				mcpFile = filepath.Join(r.sessions, "ep1.mcp.json")
			}
			task := filepath.Join(r.sessions, "task.md")
			A, err := prepareMint(mintOptions{verb: "episode start", episode: "ep1", db: db, credsDB: cred,
				wsRoot: f.ws, repo: f.repo, branch: "world/ep1", grants: grants, ttl: 3600,
				outPath: filepath.Join(r.sessions, "ep1.session"), outDirCreate: true, profile: arm.profile,
				profilesDir: r.profiles})
			if err != nil {
				t.Fatal(err)
			}
			if mcpFile != "" {
				A.afterToken = func(token string) error { return writeMCPConfig(mcpFile, "127.0.0.1:7644", token) }
			}
			mr, err := buildMintRequest(A, testBuild, "127.0.0.1:7644", arm.agent, "m1", task, mcpFile)
			if err != nil {
				t.Fatal(err)
			}
			canon, err := approval.Canonical(mr)
			if err != nil {
				t.Fatal(err)
			}
			dec, err := approval.Decode(canon)
			if err != nil {
				t.Fatal(err)
			}
			B, err := prepareMint(fromRequest(dec))
			if err != nil {
				t.Fatal(err)
			}
			B.afterToken = mcpAfterToken(dec.MCPConfig, dec.Daemon.Addr)
			if (A.afterToken == nil) != (B.afterToken == nil) || (A.beforeToken == nil) != (B.beforeToken == nil) ||
				(A.resumeRollback == nil) != (B.resumeRollback == nil) {
				t.Fatalf("callback nil-ness differs: A %v/%v/%v, B %v/%v/%v", A.afterToken == nil, A.beforeToken == nil,
					A.resumeRollback == nil, B.afterToken == nil, B.beforeToken == nil, B.resumeRollback == nil)
			}
			a, b := *A, *B
			a.verb, b.verb = "", ""
			a.afterToken, a.beforeToken, a.resumeRollback = nil, nil, nil
			b.afterToken, b.beforeToken, b.resumeRollback = nil, nil, nil
			if !reflect.DeepEqual(a, b) {
				t.Fatalf("rebuilt mint differs:\n%+v\nvs\n%+v", a, b)
			}
			if mcpFile != "" {
				A.afterToken("tok")
				first, err := os.ReadFile(mcpFile)
				if err != nil {
					t.Fatal(err)
				}
				os.Remove(mcpFile)
				B.afterToken("tok")
				second, err := os.ReadFile(mcpFile)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(first, second) {
					t.Errorf("MCP config bytes differ:\n%s\nvs\n%s", first, second)
				}
			}
		})
	}
}

// BM1: relative inputs become absolute, a grant carrying an expiry is refused, and the rendered
// request shows what approve will act on.
func TestBuildMintRequestPathsAndRefusals(t *testing.T) {
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", base)
	t.Chdir(base)
	db := filepath.Join(base, "world.db")
	if st, err := store.Open(db); err != nil {
		t.Fatal(err)
	} else if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"ws", "proj", "profiles", "sessions"} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	p := writeTestProfile(t, "profiles", "proj", "")
	grants, err := episodeStartGrants("", nil, 50, true)
	if err != nil {
		t.Fatal(err)
	}
	r, err := prepareMint(mintOptions{verb: "episode start", episode: "ep1", db: "world.db", wsRoot: "ws",
		repo: "proj", branch: "world/ep1", grants: grants, ttl: 3600,
		outPath: filepath.Join(base, "sessions", "ep1.session"), outDirCreate: true, profile: "proj",
		profilesDir: "profiles"})
	if err != nil {
		t.Fatal(err)
	}
	mcpFile := filepath.Join(base, "sessions", "ep1.mcp.json")
	mr, err := buildMintRequest(r, testBuild, "127.0.0.1:7644", "claude", "m1", "", mcpFile)
	if err != nil {
		t.Fatalf("buildMintRequest: %v", err)
	}
	canonDB, _ := store.CanonicalPath("world.db")
	if mr.WorldDB != canonDB {
		t.Errorf("WorldDB %q, want %q", mr.WorldDB, canonDB)
	}
	absWs, _ := filepath.Abs("ws")
	absRepo, _ := filepath.Abs("proj")
	absProfiles, _ := filepath.Abs("profiles")
	if mr.WorkspaceRoot != absWs || mr.Repo != absRepo || mr.ProfilesDir != absProfiles {
		t.Errorf("paths %q %q %q, want %q %q %q", mr.WorkspaceRoot, mr.Repo, mr.ProfilesDir, absWs, absRepo, absProfiles)
	}
	if mr.Daemon.VCSRevision != testBuild.VCSRevision || mr.Daemon.Addr != "127.0.0.1:7644" {
		t.Errorf("daemon ref %+v", mr.Daemon)
	}

	r.grants[0].ExpiresAt = 5
	if _, err := buildMintRequest(r, testBuild, "127.0.0.1:7644", "claude", "m1", "", mcpFile); err == nil ||
		!strings.Contains(err.Error(), "carries an expiry") {
		t.Errorf("a grant with an expiry: err %v, want 'carries an expiry'", err)
	}
	r.grants[0].ExpiresAt = 0

	const now = int64(1800000000)
	digest := strings.Repeat("ab", 32)
	out := renderRequest("mr-0102030405060708", mr, digest, now+600, now)
	for _, want := range []string{"mr-0102030405060708", digest, grantList(toCapabilities(mr.Grants)),
		mr.SessionFile, mr.MCPConfig, p.Digest, mr.Daemon.Addr, "(in 10m0s)"} {
		if !strings.Contains(out, want) {
			t.Errorf("renderRequest lacks %q:\n%s", want, out)
		}
	}
}
