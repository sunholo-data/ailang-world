package main

import (
	"reflect"
	"testing"
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
}
