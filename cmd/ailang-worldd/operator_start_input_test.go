package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
)

type operatorStartRig struct {
	f      presetFixture
	creds  string
	hosts  string
	bin    string
	action approval.OperatorAction
	at     store.OperatorMintAttempt
	in     store.OperatorMintInput
}

const operatorStartTask = "# Fix fixture\nExact signed task bytes.\n"

// newOperatorStartRig reviews one headless start: the signed resolution is the
// projection of a first live read, as the catalogue would sign it.
func newOperatorStartRig(t *testing.T) operatorStartRig {
	t.Helper()
	f := newPresetFixture(t)
	r := operatorStartRig{f: f, bin: startReadFakeAgent(t, f), hosts: t.TempDir(), at: store.OperatorMintAttempt{Catalogue: 3}}
	r.creds = matchingHealth(t, f.r.f, f.r.profiles, f.profile).CredentialsDB
	canonical, err := store.CanonicalPath(r.creds)
	if err != nil {
		t.Fatal(err)
	}
	r.action = approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: fixtureActionID, Kind: "start", Episode: "ep",
		Store: canonical, Options: approval.OperatorOptions{Preset: "p", Headless: true}, Task: operatorStartTask, Review: "fixture"}
	r.action.Argv = approval.OperatorArgv(r.action.Kind, r.action.Episode, r.action.Options)
	in, facts, err := operatorStartFacts(presetTestContext(t), r.action, r.at, f.dir, f.r.profiles, f.r.f.ws, r.hosts, presetClock)
	if err != nil {
		t.Fatal(err)
	}
	r.action.Resolution, r.in = facts, in
	if err := r.action.Validate(); err != nil {
		t.Fatalf("projection is not a signable action: %v", err)
	}
	return r
}

func (r operatorStartRig) reconstruct(t *testing.T) (store.OperatorMintInput, error) {
	t.Helper()
	return reconstructOperatorStart(presetTestContext(t), r.action, r.at, r.f.dir, r.f.r.profiles, r.f.r.f.ws, r.hosts, presetClock)
}

func (r operatorStartRig) editPreset(t *testing.T, edit func(*launchPreset)) {
	t.Helper()
	var p launchPreset
	if err := json.Unmarshal(readStartFixtureFile(t, r.f.file), &p); err != nil {
		t.Fatal(err)
	}
	edit(&p)
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	seedStartFixtureFile(t, r.f.file, data)
}

// TestOperatorStartFacts pins every reviewed fact against fixture values
// derived independently of the projection.
func TestOperatorStartFacts(t *testing.T) {
	r := newOperatorStartRig(t)
	f := r.f
	names := []string{}
	value := map[string]string{}
	for _, fact := range r.action.Resolution {
		names = append(names, fact.Name)
		value[fact.Name] = fact.Value
	}
	want := []string{"launch_binary", "launch_wall", "mcp_path", "mint_catalogue", "mint_grants", "mint_profile", "mint_project",
		"mint_start", "mint_token_path", "served_build", "served_interpreter", "worktree_path"}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("reviewed fact names %v, want %v", names, want)
	}
	binary, err := filepath.EvalSymlinks(filepath.Join(r.bin, "codex"))
	if err != nil {
		t.Fatal(err)
	}
	build, _ := json.Marshal(testBuild)
	h := matchingHealth(t, f.r.f, f.r.profiles, f.profile)
	for name, v := range map[string]string{"launch_binary": binary, "launch_wall": "1h0m0s", "mcp_path": "",
		"mint_catalogue": "3", "mint_profile": f.profile.Digest, "mint_project": "proj",
		"mint_token_path": filepath.Join(f.r.sessions, "ep.session"), "served_build": string(build),
		"served_interpreter": h.ToolRef, "worktree_path": filepath.Join(f.r.f.ws, "ep")} {
		if value[name] != v {
			t.Fatalf("fact %s = %q, want %q", name, value[name], v)
		}
	}
	in := r.in
	if in.Start.Episode != "ep" || in.Start.Base.SHA != f.base || in.Start.Preset.SHA256 != serviceHash(f.data) ||
		in.Start.Branch != "world/ep" || in.Start.BranchExisted || in.Start.Agent != "codex" || in.Start.TTLS != 43200 ||
		in.Start.Task == nil || in.Start.Task.Path != filepath.Join(f.r.sessions, "ep.task.md") ||
		in.Start.Task.SHA256 != serviceHash([]byte(operatorStartTask)) || in.TokenPath != filepath.Join(f.r.sessions, "ep.session") ||
		in.Profile != (approval.ProfileRef{Project: "proj", Digest: f.profile.Digest}) || len(in.Grants) == 0 || in.Attempt != r.at {
		t.Fatalf("effective input differs from the fixture: %+v", in)
	}
	caps, err := sessionNewGrants("se-tools", nil, 50)
	if err != nil {
		t.Fatal(err)
	}
	grants := []approval.Grant{}
	for _, c := range caps {
		grants = append(grants, approval.Grant{Effect: c.Effect, Scope: c.Scope, Budget: c.Budget})
	}
	if !reflect.DeepEqual(in.Grants, grants) {
		t.Fatalf("effective grants %+v, want the preset's %+v", in.Grants, grants)
	}
	mint, err := in.Resolution()
	if err != nil || !reflect.DeepEqual(mint, r.action.Resolution[3:9]) {
		t.Fatalf("mint facts differ from the effective input: %v", err)
	}
}

// TestOperatorStartChanged changes one live input after review at a time.
// Each change is a stale review with no lock, file, store or launch effect.
func TestOperatorStartChanged(t *testing.T) {
	for _, tc := range []struct {
		name, field string
		class       error
		change      func(*testing.T, *operatorStartRig)
	}{
		{name: "unchanged"},
		{name: "preset_ttl", field: "mint_start", change: func(t *testing.T, r *operatorStartRig) {
			r.editPreset(t, func(p *launchPreset) { p.TTLS-- })
		}},
		{name: "preset_grants", field: "mint_grants", change: func(t *testing.T, r *operatorStartRig) {
			r.editPreset(t, func(p *launchPreset) { p.Grants[0].Budget-- })
		}},
		{name: "base", field: "mint_start", change: func(t *testing.T, r *operatorStartRig) {
			presetGit(t, r.f.r.f, "branch", "-f", "dev", "feature")
		}},
		{name: "profile", field: "mint_profile", change: func(t *testing.T, r *operatorStartRig) {
			p := writeTestProfile(t, r.f.r.profiles, "proj", `, "env": {"FIXTURE": "changed"}`)
			h := matchingHealth(t, r.f.r.f, r.f.r.profiles, p)
			h.DBPath = r.f.r.f.db
			stubDaemon(t, h, nil)
		}},
		{name: "task", field: "mint_start", change: func(t *testing.T, r *operatorStartRig) { r.action.Task += "changed\n" }},
		{name: "catalogue", field: "mint_catalogue", change: func(t *testing.T, r *operatorStartRig) { r.at.Catalogue++ }},
		{name: "build", field: "served_build", change: func(t *testing.T, r *operatorStartRig) {
			other := daemon.BuildIdentity{VCSRevision: "fedcba9876543210fedcba9876543210fedcba98"}
			h := matchingHealth(t, r.f.r.f, r.f.r.profiles, r.f.profile)
			h.DBPath, h.Build = r.f.r.f.db, other
			stubDaemon(t, h, nil)
			episodeSelfBuild = func() daemon.BuildIdentity { return other }
		}},
		{name: "interpreter", field: "served_interpreter", change: func(t *testing.T, r *operatorStartRig) {
			r.action.Resolution[10].Value = "sha256:" + serviceHash([]byte("other"))
		}},
		{name: "launch_binary", field: "launch_binary", change: func(t *testing.T, r *operatorStartRig) {
			other := startReadFakeAgent(t, r.f)
			if other == r.bin {
				t.Fatal("fixture reused the agent directory")
			}
		}},
		{name: "mcp_path", field: "mcp_path", change: func(t *testing.T, r *operatorStartRig) {
			r.action.Resolution[2].Value = filepath.Join(r.f.r.sessions, "ep.mcp.json")
		}},
		{name: "worktree_path", field: "worktree_path", change: func(t *testing.T, r *operatorStartRig) {
			r.action.Resolution[11].Value += "-other"
		}},
		{name: "extra_fact", field: "resolution", change: func(t *testing.T, r *operatorStartRig) {
			r.action.Resolution = append(r.action.Resolution, approval.OperatorResolved{Name: "zz_extra", Value: "x"})
		}},
		{name: "missing_fact", field: "resolution", change: func(t *testing.T, r *operatorStartRig) {
			r.action.Resolution = r.action.Resolution[:len(r.action.Resolution)-1]
		}},
		{name: "store", field: "store", change: func(t *testing.T, r *operatorStartRig) {
			r.action.Store = filepath.Join(filepath.Dir(r.action.Store), "other.db")
		}},
		{name: "finish", class: store.ErrOperatorConflict, change: func(t *testing.T, r *operatorStartRig) {
			r.action.Kind = "finish"
		}},
		{name: "interactive", class: store.ErrOperatorConflict, change: func(t *testing.T, r *operatorStartRig) {
			r.action.Options.Headless = false
		}},
		{name: "no_profile", class: store.ErrOperatorInput, change: func(t *testing.T, r *operatorStartRig) {
			r.editPreset(t, func(p *launchPreset) {
				p.Profile, p.Checks = "", nil
				p.Grants = []presetGrant{{Effect: "Workspace.Read", Scope: "worktree", Budget: 1}}
			})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := newOperatorStartRig(t)
			if tc.change != nil {
				tc.change(t, &r)
			}
			before := observeProvision(t, r.f.r, r.creds)
			got, err := r.reconstruct(t)
			switch {
			case tc.field != "":
				var mismatch *approval.BindingError
				if !errors.As(err, &mismatch) || mismatch.Field != tc.field {
					t.Fatalf("expected stale review %s, got %v", tc.field, err)
				}
			case tc.class != nil:
				if !errors.Is(err, tc.class) {
					t.Fatalf("expected %v, got %v", tc.class, err)
				}
			default:
				if err != nil || !reflect.DeepEqual(got, r.in) {
					t.Fatalf("unchanged review did not reconstruct its input: %v", err)
				}
			}
			if !reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
				t.Fatal("reconstruction changed durable state")
			}
			for _, p := range []string{operatorEpisodeLockPath(r.creds, "ep"), filepath.Join(r.bin, "ran"), filepath.Join(r.f.r.f.ws, "ep")} {
				if _, err := os.Lstat(p); !os.IsNotExist(err) {
					t.Fatalf("reconstruction created %s: %v", p, err)
				}
			}
		})
	}
}
