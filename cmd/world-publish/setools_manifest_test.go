package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/projection"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// TestSeToolsManifestPublishes is row 134 AC5.2: the CHECKED-IN
// packages/se-tools/transitions.json, published by this verb (in-process, the
// attended fence satisfied by the test probe) with the pinned released
// interpreter into a TEST store, yields exactly the 8 descriptors of the
// ailang_only lane — IDs that are already their MCP names
// (EncodeMCPName(id) == id), each tool's effect at cost 0 as access and at
// cost 1 as its one declared effect in scope "worktree". Run from the repo
// root, because the manifest's transitionFnFile paths are repo-relative (the
// attended publish is run from the checkout). Every source passes the real
// publish loadability check (standalone `check` in an empty root), which is
// what proves the modules are self-contained. A second run is UNCHANGED.
func TestSeToolsManifestPublishes(t *testing.T) {
	bin := os.Getenv("AILANG_BIN")
	if bin == "" {
		t.Fatal("AILANG_BIN unset: publishing the se-tools manifest needs the pinned released interpreter; never skip")
	}
	if !filepath.IsAbs(bin) {
		abs, err := filepath.Abs(bin)
		if err != nil {
			t.Fatal(err)
		}
		bin = abs
	}
	storePath := filepath.Join(t.TempDir(), "world.db")
	bootstrapEpochRegistry(t, storePath, firstVersionLine(t, bin))
	t.Chdir(filepath.Join("..", ".."))

	flags := map[string]string{"store": storePath, "manifest": "packages/se-tools/transitions.json", "ailang-bin": bin}
	res := driveTransitions(t, flags, transitionsOK, noEnv)
	if res.code != exitOK || !strings.Contains(res.stdout, "published transition registry revision 1") {
		t.Fatalf("publish = (%d, %q, %q), want exit 0 and revision 1", res.code, res.stdout, res.stderr)
	}
	again := driveTransitions(t, flags, transitionsOK, noEnv)
	if again.code != exitOK || !strings.Contains(again.stdout, "UNCHANGED at revision 1") {
		t.Fatalf("republish = (%d, %q, %q), want the idempotent no-op", again.code, again.stdout, again.stderr)
	}

	db, err := store.Open(storePath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	_, rev, ok, err := transitionreg.NewReader(db).CurrentRevision(boundedTestContext(t))
	if err != nil || !ok {
		t.Fatalf("read published revision: ok=%v err=%v", ok, err)
	}
	want := map[string]string{
		"ailang-read": "Workspace.Read", "ailang-write": "Workspace.Write", "ailang-edit": "Workspace.Write",
		"ailang-check": "Ailang.Check", "ailang-run": "Ailang.Run", "builtins-search": "Ailang.Discover",
		"examples-search": "Ailang.Discover", "ailang-cli": "Ailang.CLI",
	}
	if len(rev.Entries) != len(want) {
		t.Fatalf("published %d descriptors, want %d", len(rev.Entries), len(want))
	}
	for _, d := range rev.Entries {
		effect, ok := want[d.ID]
		if !ok {
			t.Fatalf("unexpected descriptor %q", d.ID)
		}
		delete(want, d.ID)
		name, err := projection.EncodeMCPName(d.ID)
		if err != nil || name != d.ID {
			t.Fatalf("EncodeMCPName(%q) = (%q, %v), want the ID itself", d.ID, name, err)
		}
		access := transitionreg.EffectRequirement{Effect: effect, Scope: "worktree", Cost: 0}
		declared := []transitionreg.EffectRequirement{{Effect: effect, Scope: "worktree", Cost: 1}}
		if d.ID == "ailang-run" {
			// Row 135 (§4.2 L3): the Env and Net runs are their own effects.
			declared = append(declared,
				transitionreg.EffectRequirement{Effect: "Ailang.RunEnv", Scope: "worktree", Cost: 1},
				transitionreg.EffectRequirement{Effect: "Ailang.RunNet", Scope: "worktree", Cost: 1})
		}
		if d.Access != access || !reflect.DeepEqual(d.DeclaredEffects, declared) {
			t.Fatalf("%s: access %+v declared %+v; want %+v and %+v", d.ID, d.Access, d.DeclaredEffects, access, declared)
		}
		if d.Title == "" || d.Description == "" || d.SemanticsEpoch != 1 {
			t.Fatalf("%s: title %q description %d bytes epoch %d", d.ID, d.Title, len(d.Description), d.SemanticsEpoch)
		}
	}
}
