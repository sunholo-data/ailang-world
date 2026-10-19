package daemon

// Row 189 M2 (D-189-2 A): profileFor resolves an episode's exec profile at
// request time — the credential's minted binding first, then the
// --exec-episode-project map for unbound (legacy) credentials, then the only
// profile — and a bound credential whose minted digest is not the loaded
// profile's refuses every exec, naming both digests and "re-mint".

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

const (
	digestA  = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	digestB  = "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	digestB2 = "sha256:b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2"
)

func twoProfiles() (*broker.ExecProfile, *broker.ExecProfile, *workspaceExec) {
	a := &broker.ExecProfile{Project: "proj-a", Digest: digestA}
	b := &broker.ExecProfile{Project: "proj-b", Digest: digestB}
	return a, b, &workspaceExec{
		byProject:      map[string]*broker.ExecProfile{"proj-a": a, "proj-b": b},
		episodeProject: map[string]string{"legacy": "proj-a", "mapped-broken": "broken"},
		sources:        map[string]string{"proj-a": "/p/proj-a.json", "proj-b": "/p/proj-b.json"},
		excluded:       []broker.ExecProfileExclusion{{File: "/p/broken.json", Project: "broken", Reason: "exec profile: path: is empty"}},
		excludedWhy:    map[string]string{"broken": "exec profile: path: is empty"},
	}
}

func TestProfileForResolutionOrder(t *testing.T) {
	a, b, e := twoProfiles()

	// 1. The binding wins, even over a flag mapping of the same episode.
	if p, why := e.profileFor("legacy", execBinding{Project: "proj-b", Digest: digestB}); p != b || why != "" {
		t.Fatalf("bound legacy -> %v %q, want proj-b (the binding beats the flag map)", p, why)
	}
	// 2. An unbound legacy credential resolves through the flag map: the
	// rig's live credentials, minted before the binding, keep working.
	if p, why := e.profileFor("legacy", execBinding{}); p != a || why != "" {
		t.Fatalf("unbound legacy -> %v %q, want proj-a via --exec-episode-project", p, why)
	}
	// Control: unbound and unmapped with several profiles is refused, and
	// the refusal now points at --profile at mint.
	if p, why := e.profileFor("fresh", execBinding{}); p != nil || !strings.Contains(why, "--profile PROJECT") || !strings.Contains(why, "--exec-episode-project fresh=PROJECT") {
		t.Fatalf("unbound unmapped -> %v %q", p, why)
	}
	// 3. The only profile, for an unbound unmapped episode.
	only := &workspaceExec{profile: a, byProject: map[string]*broker.ExecProfile{"proj-a": a}}
	if p, _ := only.profileFor("fresh", execBinding{}); p != a {
		t.Fatalf("the only profile was not the default: %v", p)
	}
	// A binding names its project even when one profile is the default.
	if p, why := only.profileFor("fresh", execBinding{Project: "proj-b", Digest: digestB}); p != nil || !strings.Contains(why, "has not loaded") {
		t.Fatalf("bound to an unloaded project beside the only profile -> %v %q", p, why)
	}
	// A left-out directory profile refuses with its reason, bound or mapped.
	for _, tc := range []struct {
		ep string
		b  execBinding
	}{{"ep", execBinding{Project: "broken", Digest: digestA}}, {"mapped-broken", execBinding{}}} {
		if p, why := e.profileFor(tc.ep, tc.b); p != nil || !strings.Contains(why, "left out") || !strings.Contains(why, "path: is empty") {
			t.Fatalf("%s -> %v %q, want the left-out reason", tc.ep, p, why)
		}
	}
}

// The digest rule. Mutant: drop the p.Digest != b.Digest arm and the
// changed-profile episode runs under bytes nobody confirmed.
func TestProfileForRefusesAMintedDigestTheDaemonDidNotLoad(t *testing.T) {
	_, b, e := twoProfiles()
	p, why := e.profileFor("w200", execBinding{Project: "proj-b", Digest: digestB2})
	if p != nil {
		t.Fatal("a credential bound to other profile bytes got the loaded profile")
	}
	for _, want := range []string{digestB2, digestB, "re-mint", "proj-b", "w200"} {
		if !strings.Contains(why, want) {
			t.Errorf("refusal lacks %q: %s", want, why)
		}
	}
	// Control: the minted digest itself resolves.
	if p, _ := e.profileFor("w200", execBinding{Project: "proj-b", Digest: digestB}); p != b {
		t.Fatalf("the minted digest did not resolve: %v", p)
	}
}

// The same rule through the handler the registry binds: the refusal is an
// {ok:false, refused} answer that runs nothing, and a new mint (a new
// binding) is seen on the next call with no restart.
func TestExecHandlerReadsTheBindingPerRequest(t *testing.T) {
	_, _, e := twoProfiles()
	current := execBinding{Project: "proj-b", Digest: digestB2}
	w := &workspaceTools{exec: e, errLog: io.Discard, root: t.TempDir(),
		bindingFor: func(string) (execBinding, error) { return current, nil }}
	h := w.execHandler("w200", t.TempDir())
	r, ok := h.(broker.ExecRefusalHandler)
	if !ok || !strings.Contains(r.Why, "re-mint") {
		t.Fatalf("handler %T %+v, want the digest refusal", h, h)
	}
	out, err := h.Execute(context.Background(), broker.EffectRequest{Effect: broker.EffectWorkspaceExec, Scope: broker.WorkspaceScope}, nil)
	var body map[string]any
	if err != nil || json.Unmarshal(out, &body) != nil || body["ok"] != false || !strings.Contains(fmt.Sprint(body["refused"]), digestB2) {
		t.Fatalf("refusal answer %s %v", out, err)
	}
	// Re-minted against the loaded bytes: the next call is past the binding
	// gate (the handler build itself then needs a sandbox, absent here).
	current = execBinding{Project: "proj-b", Digest: digestB}
	if _, still := w.execHandler("w200", t.TempDir()).(broker.ExecRefusalHandler); still {
		t.Fatal("the re-minted binding was not read on the next request")
	}
	// A conflict between live credentials is a refusal, a store error a failure.
	w.bindingFor = func(string) (execBinding, error) {
		return execBinding{}, &bindingConflictError{episode: "w200", bindings: []string{"a", "b"}}
	}
	if _, ok := w.execHandler("w200", t.TempDir()).(broker.ExecRefusalHandler); !ok {
		t.Fatal("a binding conflict was not a refusal")
	}
	w.bindingFor = func(string) (execBinding, error) { return execBinding{}, errors.New("disk on fire") }
	if _, ok := w.execHandler("w200", t.TempDir()).(execFailedHandler); !ok {
		t.Fatal("a binding read error was not a recorded failure")
	}
}

func TestBindingOfAnEpisodesLiveCredentials(t *testing.T) {
	row := func(id int, ep, project, digest string, expires int64) store.SessionRow {
		return store.SessionRow{CredentialID: fmt.Sprintf("%064x", id), EpisodeID: ep, ExpiresAt: expires, Project: project, ProfileDigest: digest}
	}
	now := int64(1000)
	for name, tc := range map[string]struct {
		rows     []store.SessionRow
		want     execBinding
		conflict bool
	}{
		"none":              {nil, execBinding{}, false},
		"unbound only":      {[]store.SessionRow{row(1, "ep", "", "", 2000)}, execBinding{}, false},
		"bound and unbound": {[]store.SessionRow{row(1, "ep", "", "", 2000), row(2, "ep", "p", digestA, 2000)}, execBinding{"p", digestA}, false},
		"expired binding":   {[]store.SessionRow{row(1, "ep", "p", digestA, 999)}, execBinding{}, false},
		"same twice":        {[]store.SessionRow{row(1, "ep", "p", digestA, 2000), row(2, "ep", "p", digestA, 2000)}, execBinding{"p", digestA}, false},
		"two digests":       {[]store.SessionRow{row(1, "ep", "p", digestA, 2000), row(2, "ep", "p", digestB, 2000)}, execBinding{}, true},
		"two projects":      {[]store.SessionRow{row(1, "ep", "p", digestA, 2000), row(2, "ep", "q", digestA, 2000)}, execBinding{}, true},
		"other episode":     {[]store.SessionRow{row(1, "other", "p", digestA, 2000)}, execBinding{}, false},
	} {
		t.Run(name, func(t *testing.T) {
			got, err := bindingOf("ep", tc.rows, now)
			var conflict *bindingConflictError
			if errors.As(err, &conflict) != tc.conflict || (err != nil && !tc.conflict) || got != tc.want {
				t.Fatalf("bindingOf = %+v, %v; want %+v conflict=%v", got, err, tc.want, tc.conflict)
			}
		})
	}
}

// The daemon reads the binding from the credentials store, read-only, per
// request: a row minted after the "daemon" started is seen at once.
func TestEpisodeBindingReadsTheCredentialsStore(t *testing.T) {
	path := filepath.Join(t.TempDir(), store.CredentialsFileName)
	if err := ensureCredentialsStore(path, io.Discard); err != nil {
		t.Fatal(err)
	}
	d := &Daemon{credentialsPath: path, credentialBudget: 3 * time.Second}
	if b, err := d.episodeBinding("w200"); err != nil || b.bound() {
		t.Fatalf("before the mint: %+v %v", b, err)
	}
	c, err := store.OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := c.MintSession(ctx, store.SessionRow{CredentialID: fmt.Sprintf("%064x", 1), EpisodeID: "w200", GrantsJSON: "[]",
		ExpiresAt: time.Now().Unix() + 600, CreatedAt: time.Now().Unix(), Project: "proj-b", ProfileDigest: digestB}); err != nil {
		t.Fatal(err)
	}
	_ = c.Close()
	if b, err := d.episodeBinding("w200"); err != nil || b != (execBinding{"proj-b", digestB}) {
		t.Fatalf("after the mint: %+v %v", b, err)
	}
}

// /v1/health carries the build identity, the tool ref, the credentials
// store and the loaded and left-out profiles; the pre-row-189 fields stay.
func TestHealthReportsBuildToolCredentialsAndProfiles(t *testing.T) {
	_, _, e := twoProfiles()
	ref, err := hashref.Sum(hashref.AlgoSHA256, []byte("tool bytes"))
	if err != nil {
		t.Fatal(err)
	}
	d := &Daemon{cfg: Config{DBPath: "/s/world.db"}, credentialsPath: "/s/credentials.db",
		workspace: &workspaceTools{bin: "/archive/tool", binRef: ref, exec: e}}
	rec := httptest.NewRecorder()
	d.handleHealth(rec, httptest.NewRequest(http.MethodGet, "/v1/health", nil))
	var body map[string]json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"status", "daemon_version", "db_path", "interpreter_ref", "interpreter_version",
		"build", "tool_ref", "credentials_db", "exec_profiles", "exec_profiles_excluded"} {
		if _, ok := body[k]; !ok {
			t.Errorf("health lacks %q: %s", k, rec.Body.String())
		}
	}
	var h HealthResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &h); err != nil {
		t.Fatal(err)
	}
	if h.ToolRef != ref.String() || h.CredentialsDB != "/s/credentials.db" || h.Build != CurrentBuild() {
		t.Fatalf("health = %+v", h)
	}
	if len(h.ExecProfiles) != 2 || h.ExecProfiles[0] != (HealthExecProfile{"proj-a", digestA, "/p/proj-a.json"}) ||
		len(h.ExecProfilesExcluded) != 1 || h.ExecProfilesExcluded[0].Project != "broken" {
		t.Fatalf("profiles = %+v excluded %+v", h.ExecProfiles, h.ExecProfilesExcluded)
	}
	// Control: no tool binary, no exec: an empty tool ref and an empty list.
	d.workspace = &workspaceTools{}
	rec = httptest.NewRecorder()
	d.handleHealth(rec, httptest.NewRequest(http.MethodGet, "/v1/health", nil))
	if !strings.Contains(rec.Body.String(), `"tool_ref":""`) || !strings.Contains(rec.Body.String(), `"exec_profiles":[]`) {
		t.Fatalf("health without tools: %s", rec.Body.String())
	}
}

func TestSameBuild(t *testing.T) {
	clean := BuildIdentity{VCSRevision: "abc", ExecutableSHA256: "sha256:1"}
	dirty := BuildIdentity{VCSRevision: "abc", VCSModified: true, ExecutableSHA256: "sha256:1"}
	unstamped := BuildIdentity{ExecutableSHA256: "sha256:7"}
	for name, tc := range map[string]struct {
		a, b BuildIdentity
		same bool
	}{
		// The grader's case: one worktree build, so no vcs.* stamp, both sides.
		"same bytes, no stamp":             {unstamped, unstamped, true},
		"same bytes, one stamped":          {clean, BuildIdentity{ExecutableSHA256: "sha256:1"}, true},
		"same clean revision, other bytes": {clean, BuildIdentity{VCSRevision: "abc", ExecutableSHA256: "sha256:2"}, true},
		"same dirty bytes":                 {dirty, dirty, true},
		"other bytes, no stamps":           {unstamped, BuildIdentity{ExecutableSHA256: "sha256:8"}, false},
		"other bytes, dirty revision":      {dirty, BuildIdentity{VCSRevision: "abc", VCSModified: true, ExecutableSHA256: "sha256:9"}, false},
		"other bytes, clean vs dirty":      {clean, BuildIdentity{VCSRevision: "abc", VCSModified: true, ExecutableSHA256: "sha256:9"}, false},
		"other revision":                   {clean, BuildIdentity{VCSRevision: "abd"}, false},
		"no identity (old)":                {BuildIdentity{}, clean, false},
		"neither has identity":             {BuildIdentity{}, BuildIdentity{}, false},
		"other bytes, one unstamped":       {clean, unstamped, false},
	} {
		if got := SameBuild(tc.a, tc.b); got != tc.same {
			t.Errorf("%s: SameBuild = %v, want %v", name, got, tc.same)
		}
	}
}
