package main

import (
	"errors"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
)

func testOperatorCataloguePublication(t *testing.T, f *operatorLaunchRig) {
	t.Helper()
	r, c, h := f.worktreeStartRig, f.c, f.h
	r.env.operatorCatalogue = func(c runnerConfig) []byte { return r.env.readOperatorCatalogue(c, r.src) }
	if err := r.env.runnerTick(c, &h); err != nil {
		t.Fatal(err)
	}
	seen, err := runnerHeartbeat(c.Credentials, activationTestBudget)
	if err != nil || seen.CatalogueGeneration != 2 || !seen.Available(c.Executables[r.env.paths.Binary], "start", r.src.clock()) {
		t.Fatalf("published review generation: %d %v", seen.CatalogueGeneration, err)
	}
	preview, err := seen.StartReview(r.action.Options.Preset)
	if err != nil {
		t.Fatalf("missing signed start inputs: %v; catalogue %q", err, seen.Catalogue)
	}
	a, err := preview.Action(r.action.Episode, r.action.Task, seen.CatalogueGeneration)
	if err != nil {
		t.Fatal(err)
	}
	for i := range r.action.Resolution {
		if r.action.Resolution[i].Name == "mint_catalogue" {
			r.action.Resolution[i].Value = "2"
		}
	}
	if !reflect.DeepEqual(a.Resolution, r.action.Resolution) || !reflect.DeepEqual(a.Argv, r.action.Argv) {
		t.Fatalf("published review does not round-trip to exact launch authority:\n%+v\n%+v", a.Resolution, r.action.Resolution)
	}

	assertUnpublished := func(t *testing.T, c runnerConfig) {
		t.Helper()
		if err := r.env.runnerTick(c, &h); err != nil {
			t.Fatal(err)
		}
		seen, err := runnerHeartbeat(c.Credentials, activationTestBudget)
		if err != nil || !seen.Available(c.Executables[r.env.paths.Binary], "start", r.src.clock()) {
			t.Fatalf("expected a live heartbeat: %+v (%v)", seen, err)
		}
		_, reviewErr := seen.StartReview(r.action.Options.Preset)
		var actions int
		err = r.db.QueryRowContext(activationTestContext(t), "SELECT count(*) FROM operator_actions").Scan(&actions)
		if err != nil || actions != 1 || !errors.Is(reviewErr, store.ErrOperatorRunner) || string(seen.Catalogue) != "{}" {
			t.Fatalf("refused preset readback: catalogue %q, review error %v, actions %d (%v)", seen.Catalogue, reviewErr, actions, err)
		}
	}
	t.Run("BuildDigest", func(t *testing.T) {
		publish, self := r.env.operatorCatalogue, episodeSelfBuild
		t.Cleanup(func() { r.env.operatorCatalogue, episodeSelfBuild = publish, self })
		installed := r.health.Build
		installed.VCSRevision = "fixture-revision"
		episodeSelfBuild = func() daemon.BuildIdentity { return installed }
		r.env.operatorCatalogue = func(c runnerConfig) []byte {
			// The pair check sees the installed digest. The catalogue read sees
			// a different executable from the same clean revision.
			original := r.health.Build
			defer func() { r.health.Build = original }()
			r.health.Build = installed
			r.health.Build.ExecutableSHA256 = "sha256:other-executable"
			return publish(c)
		}
		assertUnpublished(t, c)
	})
	t.Run("Agent", func(t *testing.T) {
		original := readStartFixtureFile(t, r.f.file)
		t.Cleanup(func() { seedStartFixtureFile(t, r.f.file, original) })
		r.editPreset(t, func(p *launchPreset) { p.Agent = "claude" })
		binary := filepath.Join(r.bin, "claude")
		if err := os.WriteFile(binary, []byte(fakeLaunchMarker), 0700); err != nil {
			t.Fatal(err)
		}
		c := c
		c.Executables = maps.Clone(c.Executables)
		c.Executables[binary] = serviceFileHash(binary)
		assertUnpublished(t, c)
	})
	t.Run("RemoteCredentials", func(t *testing.T) {
		presetGit(t, r.f.r.f, "remote", "set-url", "origin", "https://fixture:fake-password@example.invalid/repo")
		assertUnpublished(t, c)
	})
}
